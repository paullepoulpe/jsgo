// Port of com/google/javascript/jscomp/Compiler.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"weak"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/resources"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// Compiler (and the other classes in this package) does the following:
//
//   - parses JS code
//   - checks for undefined variables
//   - performs optimizations such as constant folding and constants inlining
//   - renames variables (to short names)
//   - outputs compact JavaScript code
//
// External variables are declared in 'externs' files. For instance, the
// file may include definitions for global javascript/browser objects such
// as window, document.
//
// Go-only: language-adapter — Java's Compiler implements ErrorHandler and
// SourceFileMapping, whose methods collide in Go with AbstractCompiler's
// report(JSError) and SourceExcerptProvider's getSourceMapping(String, int,
// int) (a *string parameter there): the Compiler seen as one of those
// interfaces is AsErrorHandler() / AsSourceFileMapping().
//
// TODO(tbreisacher): Rename Compiler to JsCompiler and remove this suppression.
type Compiler struct {
	options *CompilerOptions

	passes PassConfig

	// The externs inputs
	externs []*CompilerInput

	// The source chunk graph, denoting dependencies between chunks.
	chunkGraph *JSChunkGraph

	// The module loader for resolving paths into module URIs.
	moduleLoader *ModuleLoader

	// Map of module names to module types - used for module rewriting
	moduleTypesByName *jutil.LinkedHashMap[string, CompilerInput_ModuleType]

	// error manager to which error management is delegated
	errorManager ErrorManager

	// Warnings guard for filtering warnings.
	warningsGuard WarningsGuard

	// Parse tree root nodes
	externsRoot     rhino.Node
	jsRoot          rhino.Node
	externAndJsRoot rhino.Node

	// Used for debugging; to see the compiled code between passes
	lastJsSource *string

	// Used to track the invariant about what language features may be in the AST at the given point
	// in compilation.
	allowableFeatures *parser.FeatureSet

	// Keyed by InputId.GetIdName() (InputId equality is id-name string equality).
	inputsById map[string]*CompilerInput

	scriptNodeByFilename *jutil.LinkedHashMap[string, rhino.Node]

	// Go-only: AST compaction state (compiler_compact.go).
	astCompact astCompaction

	inputPathByWebpackId *jutil.LinkedHashMap[string, string]

	transpilationNamespace rhino.StaticScope

	runtimeJsLibManager *RuntimeJsLibManager

	// Original sources referenced by the source maps.
	sourceMapOriginalSources *jutil.LinkedHashMap[string, *SourceFile]

	// Configured SourceMapInputs, plus any source maps discovered in source
	// files.
	//
	// Go-only: language-adapter — Java HashMap order. A ConcurrentHashMap
	// keyed by String whose values() are iterated (to add each input source
	// map's sourcesContent). JavaHashMap gives HashMap's iteration order,
	// which is not ConcurrentHashMap's: black-box on JDK 21, the two differ
	// for most putAll calls of 2 or more entries into an empty map (as
	// initBasedOnOptions does with the --source_map_input maps) and once 12
	// entries are put one by one. The order is observable only when
	// applyInputSourceMaps and sourceMapIncludeSourcesContent are on: which
	// sourcesContent wins when two input source maps carry different contents
	// for the same source, and the raw report order of SOURCEMAP_PARSE_FAILED /
	// SOURCEMAP_RESOLVE_FAILED (the CLI's SortingErrorManager sorts it away).
	// Matching it needs a jutil model of ConcurrentHashMap.
	inputSourceMaps *jutil.JavaHashMap[string, *SourceMapInput]

	// Map from filenames to lists of all the comments in each file.
	commentsPerFile *jutil.LinkedHashMap[string, []parser.Comment]

	// The source code map
	sourceMap *SourceMap

	// Extra context appended to the "print source after each pass" output.
	debugMessage *string

	// The externs created from the exports.
	externExports *string

	uniqueIdSupplier *UniqueIdSupplier

	// Ids for function inlining so that each declared name remains unique.
	uniqueNameId         int
	uniqueNameIdSupplier func() string

	// Whether to assume there are references to the RegExp Global object properties.
	hasRegExpGlobalReferences bool

	runJ2clPasses bool

	// Detects Google-specific coding conventions.
	defaultCodingConvention CodingConvention

	typeRegistry        *rhino.JSTypeRegistry
	colorRegistry       *colors.ColorRegistry
	parserConfig        *parsing.Config
	externsParserConfig *parsing.Config

	abstractInterpreter ReverseAbstractInterpreter
	typeValidator       *TypeValidator
	// The compiler can ask phaseOptimizer for things like which pass is currently
	// running, or which functions have been changed by optimizations
	phaseOptimizer *PhaseOptimizer

	tracker *PerformanceTracker

	// Runtime-library files deserailized from a TypedAST JAR resource; indexed by filename
	runtimeLibraryTypedAsts *jutil.LinkedHashMap[string, func() rhino.Node]

	// Types that have been forward declared
	forwardDeclaredTypes compiler_forwardDeclaredTypes

	typeCheckingHasRun bool

	// Go-only: the registry that ClearJSTypeRegistry dropped, for
	// typeRegistryReleased.
	clearedTypeRegistry weak.Pointer[rhino.JSTypeRegistry]

	// Go-only: the store this compilation's nodes are allocated in
	// (rhino/node_store.go); compaction replaces it.
	nodeStore *rhino.NodeStore

	// This error reporter gets the messages from the current Rhino parser or TypeRegistry.
	oldErrorReporter rhino.ErrorReporter

	compilerExecutor *CompilerExecutor

	outStream *jutil.PrintStream

	lastPassName *string

	externProperties *jutil.LinkedHashSet[string]
	accessorSummary  *AccessorSummary

	changeTracker *ChangeTracker

	// There is exactly one cached resolvedSourceMap, which you can think of
	// as a cache of size one.
	resolvedSourceMap compiler_ResolvedSourceMap
	sourceMapMu       sync.RWMutex

	typedAstFilesystem          *jutil.LinkedHashMap[*SourceFile, func() rhino.Node]
	typedAstScriptDeserializers map[*SourceFile]*ScriptNodeDeserializer

	validityCheck *PassFactory

	currentPassIndex int

	stage AbstractCompiler_LifeCycleStage

	typedScopeCreator *TypedScopeCreator

	topScope *TypedScope

	preferRegexParser bool

	// Non-static because this file represents different content for every
	// Compiler (and TSAN complains if one instance is shared by all
	// threads).
	SYNTHETIC_EXTERNS_FILE *SourceFile

	syntheticExternsInput *CompilerInput // matches SYNTHETIC_EXTERNS_FILE

	// Non-static because this file represents different content for every
	// Compiler (and TSAN complains if one instance is shared by all
	// threads).
	syntheticTypeSummaryFile *SourceFile

	syntheticTypeSummaryInput *CompilerInput // matches syntheticTypeSummaryFile

	// Ids for cross-chunk method stubbing, so that each method has a unique id.
	crossChunkIdGenerator *IdGenerator

	// Arguments passed to getCssName() found during compilation.
	cssNames *jutil.LinkedHashSet[string]

	// The variable renaming map
	variableMap *VariableMap

	// The property renaming map
	propertyMap *VariableMap

	// String replacement map
	stringMap *VariableMap

	// Mapping for Instrumentation parameter encoding
	instrumentationMapping *VariableMap

	// Id generator map
	idGeneratorMap *string

	// Conformance violations persisted across compilation stages.
	conformanceViolations             []*ConformanceViolationsProto
	unattributedConformanceViolations []*ViolationProto

	// Did any source file need any feature transpiled
	transpiledFiles bool

	// Names exported by goog.exportSymbol.
	exportedNames *jutil.LinkedHashSet[string]

	// Names defined by goog.define.
	defineNames *jutil.LinkedHashSet[string]

	moduleMetadataMap *ModuleMetadataMap

	moduleMap *ModuleMap

	// Go-only: language-adapter — the Compiler as an ErrorHandler and as a
	// SourceFileMapping (see the type comment); one instance each, so that
	// the same object is handed out every time, as `this` is in Java.
	asErrorHandler      *compiler_errorHandler
	asSourceFileMapping *compiler_sourceFileMapping
	cachedAstAnalyzer   *AstAnalyzer
	printScratch        codePrinterScratch
	printScratchBusy    atomic.Bool

	// Go-only: language-adapter — the object's monitor, which its
	// synchronized methods and blocks hold.
	monitor jutil.Monitor
}

// Java: Compiler.CHUNK_DEPENDENCY_ERROR
var Compiler_CHUNK_DEPENDENCY_ERROR = DiagnosticType_Error(
	"JSC_CHUNK_DEPENDENCY_ERROR",
	"Bad dependency: {0} -> {1}. "+"Chunks must be listed in dependency order.")

// Java: Compiler.MISSING_ENTRY_ERROR
var Compiler_MISSING_ENTRY_ERROR = DiagnosticType_Error(
	"JSC_MISSING_ENTRY_ERROR", "required entry point \"{0}\" never provided")

// Java: Compiler.MISSING_CHUNK_ERROR
var Compiler_MISSING_CHUNK_ERROR = DiagnosticType_Error(
	"JSC_MISSING_CHUNK_ERROR", "unknown chunk \"{0}\" specified in entry point spec")

// Compiler_ExternalSourceLoader: subclasses are responsible for loading
// sources that were not provided as explicit inputs to the compiler. For
// example, looking up sources referenced within sourcemaps.
//
// Java: Compiler.ExternalSourceLoader
type Compiler_ExternalSourceLoader struct {
	// Go-only: language-adapter — distinct identity
	_ byte
}

// Java: Compiler.ExternalSourceLoader()
func NewCompiler_ExternalSourceLoader() *Compiler_ExternalSourceLoader {
	return &Compiler_ExternalSourceLoader{}
}

// Java: Compiler.ExternalSourceLoader.loadSource(String)
func (l *Compiler_ExternalSourceLoader) LoadSource(filename string) *SourceFile {
	panic(jutil.NewRuntimeException("Cannot load without a valid loader."))
}

// Compiler_logger is the logger for the whole com.google.javascript.jscomp
// domain - setting configuration for this logger affects all loggers in
// other classes within the compiler.
//
// Java: Compiler.logger
var Compiler_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp")

// Java: Compiler.pathJoiner (Joiner.on(Platform.getFileSeperator()))
func compiler_pathJoiner(parts []string) string {
	return strings.Join(parts, Platform_GetFileSeperator())
}

// compiler_ResolvedSourceMap: when mapping symbols from a source map, we
// must repeatedly combine the path of the original file with the path from
// the source map to compute the SourceFile of the underlying code. When
// iterating through adjacent symbols this computation ends up computing the
// same thing.
//
// This object tracks the last source map we resolved to reduce that work.
// It caches: for a given originalPath+sourceMapPath, what sourceFile did we
// find?
//
// Java: Compiler.ResolvedSourceMap
type compiler_ResolvedSourceMap struct {
	hasCached     bool
	originalPath  string
	sourceMapPath string
	relativePath  string
}

// compiler_forwardDeclaredTypes is the Set<String> of forward-declared
// types: a LinkedHashSet, or the "everything is forward declared" set.
//
// Go-only: language-adapter — the Java field's static type Set<String>.
type compiler_forwardDeclaredTypes interface {
	jutil.Set[string]
	Add(e string) bool
}

// compiler_allForwardDeclared is the anonymous AbstractSet<String> of
// initOptions: it reports all types as forward declared types.
type compiler_allForwardDeclared struct {
	// Go-only: language-adapter — distinct identity
	_ byte
}

func (*compiler_allForwardDeclared) Contains(o string) bool {
	return true // Report all types as forward declared types.
}

func (*compiler_allForwardDeclared) Add(e string) bool {
	return false
}

func (*compiler_allForwardDeclared) All() iter.Seq[string] {
	return func(yield func(string) bool) {}
}

func (*compiler_allForwardDeclared) Size() int {
	return 0
}

func (*compiler_allForwardDeclared) IsEmpty() bool {
	return true
}

// initFields runs the Java field initializers, in declaration order.
func (c *Compiler) initFields() {
	c.externs = nil
	c.inputsById = make(map[string]*CompilerInput)
	c.scriptNodeByFilename = jutil.NewLinkedHashMap[string, rhino.Node]()
	c.sourceMapOriginalSources = jutil.NewLinkedHashMap[string, *SourceFile]()
	c.inputSourceMaps = jutil.NewJavaHashMap[string, *SourceMapInput]()
	c.commentsPerFile = jutil.NewLinkedHashMap[string, []parser.Comment]()
	c.uniqueIdSupplier = NewUniqueIdSupplier()
	c.uniqueNameId = 0
	c.hasRegExpGlobalReferences = true
	c.runJ2clPasses = false
	c.defaultCodingConvention = NewClosureCodingConvention()
	c.forwardDeclaredTypes = jutil.NewLinkedHashSet[string]()
	c.typeCheckingHasRun = false
	c.asErrorHandler = &compiler_errorHandler{c: c}
	c.asSourceFileMapping = &compiler_sourceFileMapping{c: c}
	c.oldErrorReporter = RhinoErrorReporter_ForOldRhino(c.asErrorHandler)
	c.compilerExecutor = c.CreateCompilerExecutor()
	c.changeTracker = NewChangeTracker()
	c.resolvedSourceMap = compiler_ResolvedSourceMap{}
	c.validityCheck = PassFactory_Builder_().
		SetName("validityCheck").
		SetRunInFixedPointLoop(true).
		SetInternalFactory(func(compiler AbstractCompiler) CompilerPass { return NewValidityCheck(compiler) }).
		// TODO(b/280906684): Cannot use a lambda nor a simple anonymous function (which might be
		// converted to lambdas) directly due to this inference problem. Needs the cast to
		// make jscompiler happy.
		SetCondition(
			func(o *CompilerOptions) bool {
				panic(jutil.NewIllegalStateExceptionMsg("Unexpected"))
			}).
		Build()
	c.currentPassIndex = -1
	c.stage = AbstractCompiler_LifeCycleStage_RAW
	c.preferRegexParser = false
	c.SYNTHETIC_EXTERNS_FILE = SourceFile_FromCodeWithKind(compiler_SYNTHETIC_FILE_NAME_PREFIX+"externs] ", "", rhino.StaticSourceFile_SourceKind_EXTERN)
	c.syntheticTypeSummaryFile = SourceFile_FromCodeWithKind(compiler_SYNTHETIC_FILE_NAME_PREFIX+"typeSummary] ", "", rhino.StaticSourceFile_SourceKind_EXTERN)
	c.crossChunkIdGenerator = NewIdGenerator()
	c.exportedNames = jutil.NewLinkedHashSet[string]()
	c.defineNames = jutil.NewLinkedHashSet[string]()
}

// NewCompiler creates a Compiler that reports errors and warnings to its
// logger.
//
// Java: Compiler()
func NewCompiler() *Compiler {
	return NewCompilerWithOutStream(nil)
}

// NewCompilerWithOutStream creates a Compiler that reports errors and
// warnings to an output stream.
//
// Java: Compiler(PrintStream)
func NewCompilerWithOutStream(outStream *jutil.PrintStream) *Compiler {
	// Go-only: the first compiler of a run (or of a test) starts the node
	// store; compilers created after it share it (rhino/node_store.go).
	c := &Compiler{nodeStore: rhino.EnsureCompilationStore()}
	c.initFields()
	c.changeTracker.AddChangeHandler(c.changeTracker.GetRecentChange())
	c.outStream = outStream
	c.moduleTypesByName = jutil.NewLinkedHashMap[string, CompilerInput_ModuleType]()
	return c
}

// NewCompilerWithErrorManager creates a Compiler that uses a custom error
// manager.
//
// Java: Compiler(ErrorManager)
func NewCompilerWithErrorManager(errorManager ErrorManager) *Compiler {
	c := NewCompiler()
	c.SetErrorManager(errorManager)
	return c
}

// compiler_errorHandler is the Compiler as an ErrorHandler.
type compiler_errorHandler struct {
	c *Compiler
}

// Report is Compiler.report(CheckLevel, JSError).
func (h *compiler_errorHandler) Report(level *CheckLevel, error *JSError) {
	h.c.ReportWithLevel(level, error)
}

// compiler_sourceFileMapping is the Compiler as a SourceFileMapping.
type compiler_sourceFileMapping struct {
	c *Compiler
}

// GetSourceMapping is Compiler.getSourceMapping(String, int, int).
func (m *compiler_sourceFileMapping) GetSourceMapping(fileName string, lineNo int, columnNo int) *sourcemap.Mapping_OriginalMapping {
	return m.c.GetSourceMapping(&fileName, lineNo, columnNo)
}

// AsErrorHandler is the Compiler as an ErrorHandler (Java's `this`).
func (c *Compiler) AsErrorHandler() ErrorHandler { return c.asErrorHandler }

// AsSourceFileMapping is the Compiler as a SourceFileMapping (Java's
// `this`).
func (c *Compiler) AsSourceFileMapping() SourceFileMapping { return c.asSourceFileMapping }

// SetErrorManager sets the error manager.
//
// errorManager: the error manager, it cannot be null
//
// Java: Compiler.setErrorManager(ErrorManager)
func (c *Compiler) SetErrorManager(errorManager ErrorManager) {
	jutil.CheckNotNullMsg(errorManager, "the error manager cannot be null")
	c.errorManager = NewThreadSafeDelegatingErrorManager(errorManager)
}

// createMessageFormatter creates a message formatter instance
// corresponding to the value of CompilerOptions.
//
// Java: Compiler.createMessageFormatter()
func (c *Compiler) createMessageFormatter() MessageFormatter {
	colorize := c.options.ShouldColorizeErrorOutput()
	return c.options.GetErrorFormat().ToFormatter(c, colorize)
}

// Java: Compiler.initExperimentalOutputFeatureSetOptions(CompilerOptions)
func (c *Compiler) initExperimentalOutputFeatureSetOptions(options *CompilerOptions) {
	featureSet := options.GetOutputFeatureSet()
	if jutil.ObjectsEquals(featureSet, parser.FeatureSet_ES5) || jutil.ObjectsEquals(featureSet, parser.FeatureSet_ES3) {
		// do not change anything if language_out is configured to be ES5 or ES3.
		return
	}

	if options.GetExperimentalOutputFeatureSet() != nil {
		switch *options.GetExperimentalOutputFeatureSet() {
		case CompilerOptions_ExperimentalOutputFeatureSet_ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS:
			options.SetOutputFeatureSet(
				parser.FeatureSet_ES5.WithFeatures(
					parser.FeatureSet_Feature_ASYNC_FUNCTIONS,
					parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION,
					parser.FeatureSet_Feature_CONST_DECLARATIONS,
					parser.FeatureSet_Feature_GENERATORS,
					parser.FeatureSet_Feature_LET_DECLARATIONS))
		case CompilerOptions_ExperimentalOutputFeatureSet_BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD:
			options.SetOutputFeatureSet(
				parser.FeatureSet_BROWSER_2019.WithoutFeatures(
					parser.FeatureSet_Feature_CLASSES,
					parser.FeatureSet_Feature_CLASS_GETTER_SETTER,
					parser.FeatureSet_Feature_NEW_TARGET,
					parser.FeatureSet_Feature_SUPER,
					parser.FeatureSet_Feature_SPREAD_EXPRESSIONS))
		}
	}
}

// InitOptions initializes the compiler options. It's called as part of a
// normal compile() job. Public for the callers that are not doing a normal
// compile() job.
//
// Java: Compiler.initOptions(CompilerOptions)
func (c *Compiler) InitOptions(options *CompilerOptions) {
	c.options = options
	c.initASTCompaction()
	c.allowableFeatures = options.GetLanguageIn().ToFeatureSet()
	c.initExperimentalOutputFeatureSetOptions(options)
	if c.errorManager == nil {
		if c.outStream == nil {
			c.SetErrorManager(NewLoggerErrorManagerWithFormatter(c.createMessageFormatter(), Compiler_logger))
		} else {
			// ImmutableSet.builder(): duplicates (by identity) are dropped.
			var builder []SortingErrorManager_ErrorReportGenerator
			add := func(g SortingErrorManager_ErrorReportGenerator) {
				for _, e := range builder {
					if jutil.ObjectsEquals(e, g) {
						return
					}
				}
				builder = append(builder, g)
			}
			add(NewPrintStreamErrorReportGenerator(
				c.createMessageFormatter(), c.outStream, options.GetSummaryDetailLevel()))
			for _, g := range options.GetExtraReportGenerators() {
				add(g)
			}
			c.SetErrorManager(NewSortingErrorManager(builder))
		}
	}

	if options.GetSkipNonTranspilationPasses() {
		options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_NO_OP)
	}

	c.moduleLoader = moduleLoader_NewEmpty()

	c.ReconcileOptionsWithGuards()
	switch options.GetShouldValidateRequiredInlinings() {
	case base.Tri_TRUE:
		sufficientOptimizations := InlineFunctions_CheckOptimizationLevelSupportsRequireInlining(options)
		jutil.CheckStateMsg1(
			sufficientOptimizations.ShouldValidate(),
			"shouldValidateRequiredInlinings cannot be set unless CompilerOptions enables"+
				" sufficient optimization, but found: %s",
			sufficientOptimizations.OptMessage())
	case base.Tri_FALSE:
		// do nothing.
	case base.Tri_UNKNOWN:
		options.SetValidateRequiredInlinings(
			InlineFunctions_CheckOptimizationLevelSupportsRequireInlining(options).
				ShouldValidate())
	}

	// TODO(johnlenz): generally, the compiler should not be changing the options object
	// provided by the user.  This should be handled a different way.

	// Turn off type-based optimizations when type checking is off
	if !options.IsTypecheckingEnabled() {
		options.SetUseTypesForLocalOptimization(false)
		options.SetUseTypesForOptimization(false)
	}

	if options.AssumeForwardDeclaredForMissingTypes() {
		c.forwardDeclaredTypes = &compiler_allForwardDeclared{}
	}

	c.initWarningsGuard(options.GetWarningsGuard())
	if c.IsDebugLoggingEnabled() {
		// If debug logs are requested, then we'll always log the configuration for convenience.
		options.SetPrintConfig(true)
	}

	// If pruning many unused inputs, it's possible that it's faster to regex-parse everything and
	// AST-parse only the non-pruned inputs, rather than just AST-parse everything once.
	c.preferRegexParser = c.preferRegexParser || options.GetDependencyOptions().ShouldPrune()

	if options.GetMergedPrecompiledLibraries() &&
		options.GetDependencyOptions().NeedsManagement() {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Using precompiled libraries (i.e. TypedAST) is incompatible with flags that " +
				"automatically order/prune dependencies: " +
				jutil.ToJavaString(options.GetDependencyOptions()) +
				""))
	}

	if options.GetMergedPrecompiledLibraries() && len(options.GetConformanceConfigs()) != 0 {
		conformanceConfigFiles := jutil.NewLinkedHashSet[string]()
		for _, conformanceConfig := range options.GetConformanceConfigs() {
			for _, requirement := range conformanceConfig.GetRequirementList() {
				conformanceConfigFiles.AddAll(requirement.GetConfigFileList()...)
			}
		}
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"use_precompiled_libraries mode does not support passing JS conformance configs: this" +
				" mode skips all check passes including conformance\n" +
				"Unsupported conformance configs: " +
				conformanceConfigFiles.String()))
	}

	if options.ShouldRunDeadAssignmentElimination() &&
		!(options.ShouldRemoveUnusedVariables() || options.ShouldRemoveUnusedLocalVariables()) {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Invalid flag combination: enabling dead assignment elimination also requires unused" +
				" variable removal. Either disable dead assignment elimination or enable unused" +
				" variable removal."))
	}
	if options.ShouldRunDeadPropertyAssignmentElimination() && options.IsPolymerPassEnabled() {
		// The Polymer source is usually not included in the compilation, but it creates
		// getters/setters for many properties in compiled code. Dead property property assignment
		// elimination is only safe when it knows about getters/setters.
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Invalid flag combination: dead assignment elimination is incompatible with" +
				" setPolymerPass(true)."))
	}
}

// Java: Compiler.printConfig()
func (c *Compiler) PrintConfig() {
	if c.IsDebugLoggingEnabled() {
		// Log to separate files for convenience.
		c.logToFile("externs.log", func() string { return jutil.CollectionToString(c.externs) })
		c.logToFile(
			"inputs.json",
			func() string {
				return jutil.NewGsonBuilder().SetPrettyPrinting().Create().ToJson(c.chunkGraph.ToJson())
			})
		c.logToFile("options.log", func() string { return c.options.String() })
		c.logToFile("warningsGuard.log", func() string { return jutil.ToJavaString(c.warningsGuard) })
	} else {
		// Debug logging is not enabled, so use stderr
		err := jutil.System_Err()
		err.Println("==== Externs ====")
		err.Println(jutil.CollectionToString(c.externs))
		err.Println("==== Inputs ====")
		err.Println(jutil.NewGsonBuilder().SetPrettyPrinting().Create().ToJson(c.chunkGraph.ToJson()))
		err.Println("==== CompilerOptions ====")
		err.Println(c.options.String())
		err.Println("==== WarningsGuard ====")
		err.Println(jutil.ToJavaString(c.warningsGuard))
	}
}

// Java: Compiler.maybeLogPrunedInputs()
func (c *Compiler) maybeLogPrunedInputs() {
	if c.IsDebugLoggingEnabled() {
		c.logToFile(
			"inputs_pruned.json",
			func() string {
				return jutil.NewGsonBuilder().SetPrettyPrinting().Create().ToJson(c.chunkGraph.ToJson())
			})
	}
}

// Java: Compiler.logToFile(String, Supplier<String>)
func (c *Compiler) logToFile(logFileName string, logStringSupplier func() string) {
	log := c.CreateOrReopenLog("Compiler", logFileName)
	defer log.Close()
	log.LogSupplier(logStringSupplier)
}

// Java: Compiler.initWarningsGuard(WarningsGuard)
func (c *Compiler) initWarningsGuard(warningsGuard WarningsGuard) {
	guards := []WarningsGuard{}
	guards = append(guards,
		NewJ2clSuppressWarningsGuard(),
		NewClosureUnawareCodeWarningsGuard(),
		NewSuppressDocWarningsGuard(c, DiagnosticGroups_GetRegisteredGroups()),
		jutil.CheckNotNull(warningsGuard))

	c.warningsGuard = NewComposeWarningsGuardFromList(guards)
}

// ReconcileOptionsWithGuards: when the CompilerOptions and its
// WarningsGuard overlap, reconcile any discrepancies.
//
// Java: Compiler.reconcileOptionsWithGuards()
func (c *Compiler) ReconcileOptionsWithGuards() {
	options := c.options
	// DiagnosticGroups override the plain checkTypes option.
	if options.Enables(DiagnosticGroups_CHECK_TYPES) {
		options.SetCheckTypes(true)
	} else if options.Disables(DiagnosticGroups_CHECK_TYPES) {
		options.SetCheckTypes(false)
	} else if !options.GetCheckTypes() {
		// If DiagnosticGroups did not override the plain checkTypes
		// option, and checkTypes is disabled, then turn off the
		// parser type warnings.
		options.SetWarningLevel(
			DiagnosticGroup_ForType(RhinoErrorReporter_TYPE_PARSE_ERROR), CheckLevel_OFF)
	}

	if !options.GetCheckTypes() {
		options.SetWarningLevel(DiagnosticGroups_BOUNDED_GENERICS, CheckLevel_OFF)
	}

	// All passes must run the variable check. This synthesizes
	// variables later so that the compiler doesn't crash. It also
	// checks the externs file for validity. If you don't want to warn
	// about missing variable declarations, we shut that specific
	// error off.
	if !options.GetCheckSymbols() && !options.Enables(DiagnosticGroups_CHECK_VARIABLES) {
		options.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_OFF)
	}

	// If we're in transpile-only mode, we don't need to do checks for suspicious var usages.
	// Since we still have to run VariableReferenceCheck before transpilation to check block-scoped
	// variables, though, we disable the unnecessary warnings it produces relating to vars here.
	if options.GetSkipNonTranspilationPasses() &&
		!options.Enables(DiagnosticGroups_CHECK_VARIABLES) {
		options.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_OFF)
	}

	// If we're in transpile-only mode, we don't need to check for missing requires unless the user
	// explicitly enables missing-provide checks.
	if options.GetSkipNonTranspilationPasses() &&
		!options.Enables(DiagnosticGroups_MISSING_PROVIDE) {
		options.SetWarningLevel(DiagnosticGroups_MISSING_PROVIDE, CheckLevel_OFF)
	}

	// The default warning level of ARTIFICIAL_FUNCTION_PURITY_VALIDATION varies with whether or not
	// property disambiguation is enabled: if yes, it defaults to ERROR; if no, it defaults to OFF.
	// This is because without property disambiguation, it's often impossible for the compiler to
	// enforce function purity; and we expect projects that are code-size-sensitive enough to care
	// to have disambiguation already enabled.
	// We still want to allow individual projects to explicitly disable or enable this error.
	if options.ShouldDisambiguateProperties() &&
		options.
			GetWarningsGuard().
			MustRunChecks(DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION) == base.Tri_UNKNOWN {
		options.SetWarningLevel(
			DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION, CheckLevel_ERROR)
	}
}

// Java: Compiler.getTypedAstDeserializer(SourceFile)
func (c *Compiler) GetTypedAstDeserializer(file *SourceFile) func() rhino.Node {
	if c.typedAstFilesystem == nil {
		return nil
	}

	sd := c.typedAstScriptDeserializers[file]
	delete(c.typedAstScriptDeserializers, file)
	ast := c.typedAstFilesystem.Remove(file)
	if sd != nil {
		return sd.DeserializeNew
	}
	jutil.CheckStateMsg1(
		ast != nil ||
			strings.HasPrefix(file.GetName(), compiler_SYNTHETIC_FILE_NAME_PREFIX) ||
			AbstractCompiler_IsFillFileName(file.GetName()),
		"TypedAST filesystem initialized, but missing requested file: %s",
		file)
	return ast
}

func (c *Compiler) takeTypedAstDeserializer(file *SourceFile) (*ScriptNodeDeserializer, func() rhino.Node) {
	if c.typedAstFilesystem == nil {
		return nil, nil
	}
	sd := c.typedAstScriptDeserializers[file]
	delete(c.typedAstScriptDeserializers, file)
	ast := c.typedAstFilesystem.Remove(file)
	jutil.CheckStateMsg1(
		ast != nil ||
			strings.HasPrefix(file.GetName(), compiler_SYNTHETIC_FILE_NAME_PREFIX) ||
			AbstractCompiler_IsFillFileName(file.GetName()),
		"TypedAST filesystem initialized, but missing requested file: %s",
		file)
	return sd, ast
}

// InitWithTypedAstFilesystem initializes a compiler with deserialized state
// from the given TypedAst.List
//
// This method initializes all the state needed to run `.stage2Passes()` or
// do any similar sort of optimization work.
//
// typedAstListStream: a gzipped, binary-serialized TypedAst.List proto
//
// Java: Compiler.initWithTypedAstFilesystem(List<SourceFile>, List<SourceFile>, CompilerOptions, InputStream)
func (c *Compiler) InitWithTypedAstFilesystem(
	externs []*SourceFile,
	sources []*SourceFile,
	options *CompilerOptions,
	typedAstListStream io.Reader) {

	files := jutil.NewLinkedHashSet[*SourceFile]()
	files.AddAll(externs...)
	files.AddAll(sources...)

	options.SetMergedPrecompiledLibraries(true)

	c.InitOptions(options)
	c.Init(externs, sources, options)
	c.mergeAndDeserializeTypedAsts(files, typedAstListStream, options)
}

// InitChunksWithTypedAstFilesystem initializes a compiler with deserialized
// state from the given TypedAst.List
//
// This method initializes all the state needed to run `.stage2Passes()` or
// do any similar sort of optimization work.
//
// typedAstListStream: a gzipped, binary-serialized TypedAst.List proto
//
// Java: Compiler.initChunksWithTypedAstFilesystem(List<SourceFile>, List<JSChunk>, CompilerOptions, InputStream)
func (c *Compiler) InitChunksWithTypedAstFilesystem(
	externs []*SourceFile,
	chunks []*JSChunk,
	options *CompilerOptions,
	typedAstListStream io.Reader) {

	filesBuilder := jutil.NewLinkedHashSet[*SourceFile]()
	filesBuilder.AddAll(externs...)
	for _, chunk := range chunks {
		for _, input := range chunk.GetInputs() {
			filesBuilder.Add(input.GetSourceFile())
		}
	}
	files := filesBuilder

	options.SetMergedPrecompiledLibraries(true)

	c.InitOptions(options)
	c.InitChunks(externs, chunks, options)
	c.mergeAndDeserializeTypedAsts(files, typedAstListStream, options)
}

// Java: Compiler.mergeAndDeserializeTypedAsts(ImmutableSet<SourceFile>, InputStream, CompilerOptions)
func (c *Compiler) mergeAndDeserializeTypedAsts(
	requiredInputFiles *jutil.LinkedHashSet[*SourceFile],
	typedAstListStream io.Reader,
	options *CompilerOptions) {
	jutil.CheckState(c.typedAstFilesystem == nil)
	c.MaybeSetTracker()

	c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC)
	// To speed up builds that don't run type-based optimizations, skip type deserialization
	deserializeTypes := options.RequiresTypesForOptimization()

	astData := Compiler_RunInCompilerThread(c,
		func() *TypedAstDeserializer_DeserializedAst {
			tracer := c.NewTracer("deserializeTypedAst")
			defer c.StopTracer(tracer, "deserializeTypedAst")
			return TypedAstDeserializer_DeserializeFullAst(
				c,
				c.SYNTHETIC_EXTERNS_FILE,
				requiredInputFiles,
				typedAstListStream,
				deserializeTypes,
				options.GetResolveSourceMapAnnotations(),
				options.GetParseInlineSourceMaps())
		})

	c.typedAstFilesystem = astData.rawFilesystem()
	c.typedAstScriptDeserializers = astData.scriptDeserializers
	c.externProperties = astData.GetExternProperties()
	c.colorRegistry = astData.GetColorRegistry()
	c.SetTypeCheckingHasRun(deserializeTypes)

	for library := range astData.GetRuntimeLibraries().All() {
		c.runtimeJsLibManager.EnsureLibraryInjected(library, false)
	}

	c.GetSynthesizedExternsInput() // Force lazy creation.

	// TODO(lharker): refactor things to avoid this misleading 'parse' call. It's not actually
	// parsing JS source code, just doing more deserialization & building the Rhino AST.
	c.runInCompilerThread(c.parseForCompilationInternal)
}

// compiler_getResourceAsStream is Compiler.class.getResourceAsStream(path):
// nil (Java null) when there is no such resource.
//
// Go-only: language-adapter — classpath resources are embedded in package
// resources.
func compiler_getResourceAsStream(path string) io.Reader {
	data := resources.ResourceLoader_GetResourceAsStream("com.google.javascript.jscomp.Compiler", path)
	if data == nil {
		return nil
	}
	return bytes.NewReader(data)
}

// Java: Compiler.initRuntimeLibraryTypedAsts(Optional<ColorPool.Builder>)
func (c *Compiler) InitRuntimeLibraryTypedAsts(colorPoolBuilder *ColorPool_Builder) {
	jutil.CheckState(c.runtimeLibraryTypedAsts == nil)

	path := strings.Join([]string{
		"/runtime_libs.typedast"}, "")

	var stream io.Reader
	if data := resources.ResourceLoader_GetResourceAsStream("com.google.javascript.jscomp.Compiler", path); data != nil {
		stream = &runtimeLibsTypedAstStream{data: data}
	} else {
		stream = compiler_getResourceAsStream(path)
	}

	astData := TypedAstDeserializer_DeserializeRuntimeLibraries(
		c,
		c.SYNTHETIC_EXTERNS_FILE,
		colorPoolBuilder,
		stream,
		c.GetOptions().GetResolveSourceMapAnnotations(),
		c.GetOptions().GetParseInlineSourceMaps())

	// Re-index the runtime libraries by file name rather than SourceFile object
	runtimeLibraryTypedAsts := jutil.NewLinkedHashMap[string, func() rhino.Node]()
	for key, value := range astData.GetFilesystem().All() {
		runtimeLibraryTypedAsts.ComputeIfAbsent(
			key.GetName(), func(f string) func() rhino.Node { return value })
	}

	c.runtimeLibraryTypedAsts = jutil.NewLinkedHashMapFrom(runtimeLibraryTypedAsts)
}

// Init initializes the instance state needed for a compile job.
//
// Java: Compiler.init(List<SourceFile>, List<SourceFile>, CompilerOptions)
func (c *Compiler) Init(externs []*SourceFile, sources []*SourceFile, options *CompilerOptions) {
	chunk := NewJSChunk(JSChunk_STRONG_CHUNK_NAME)
	for _, source := range sources {
		chunk.Add(NewCompilerInput(source))
	}

	chunks := make([]*JSChunk, 0, 1)
	chunks = append(chunks, chunk)
	c.InitChunks(externs, chunks, options)
	c.addFilesToSourceMap(sources)
}

// InitChunks initializes the instance state needed for a compile job if
// the sources are in chunks.
//
// Java: Compiler.initChunks(List<SourceFile>, List<JSChunk>, CompilerOptions)
func (c *Compiler) InitChunks(externs []*SourceFile, chunks []*JSChunk, options *CompilerOptions) {
	c.InitOptions(options)

	c.checkFirstChunk(chunks)

	clear(c.externs)
	c.externs = c.externs[:0]
	for _, file := range externs {
		c.externs = append(c.externs, NewCompilerInputWithIsExtern(file /* isExtern= */, true))
	}

	// Generate the chunk graph, and report any errors in the chunk specification as errors.
	failed := false
	jutil.Try(func() {
		c.chunkGraph = NewJSChunkGraph(chunks)
	}, func(e *JSChunkGraph_ChunkDependenceException) {
		// problems with the chunk format.  Report as an error.  The
		// message gives all details.
		c.Report(
			JSError_Make(
				Compiler_CHUNK_DEPENDENCY_ERROR, e.GetChunk().GetName(), e.GetDependentChunk().GetName()))
		failed = true
	})
	if failed {
		return
	}

	// Creating the chunk graph can move weak source around, and end up with empty chunks.
	c.fillEmptyChunks(c.GetChunks())

	c.commentsPerFile = jutil.NewLinkedHashMap[string, []parser.Comment]()
	c.InitBasedOnOptions()

	c.InitInputsByIdMap()

	c.initAST()

	if c.IsDebugLoggingEnabled() || options.ShouldPrintConfig() {
		// Always print configuration logs when debug logging is enabled.
		// NOTE: initChunks() is called by every execution path through the compiler code that
		// actually performs a compilation action. That's what makes this a good place to put this
		// logging.
		c.PrintConfig()
	}

	// If debug logging is enabled, write out the module / chunk graph in GraphViz format.
	// This graph is often too big to reasonably render.
	// Using gvpr(1) is recommended to extract the parts of the graph that are of interest.
	// `dot -Tpng graph_file.dot > graph_file.png` will render an image.
	func() {
		chunkGraphLog := c.CreateOrReopenLog("Compiler", "chunk_graph.dot")
		defer chunkGraphLog.Close()
		if chunkGraphLog.IsLogging() {
			chunkGraphLog.LogString(DotFormatter_ToDotOfGraphvizGraph(c.chunkGraph.ToGraphvizGraph()))
		}
	}()

	var nodeForCodeInsertion func() rhino.Node
	if options.GetRuntimeLibraryMode() == RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES {
		nodeForCodeInsertion = c.getSynthesizedExternsRoot
	} else {
		nodeForCodeInsertion = c.getNodeForRuntimeCodeInsertion
	}
	c.runtimeJsLibManager = RuntimeJsLibManager_Create(
		options.GetRuntimeLibraryMode(),
		c.loadResourceContents,
		c.changeTracker,
		nodeForCodeInsertion)
}

// Java: Compiler.getRuntimeJsLibManager()
func (c *Compiler) GetRuntimeJsLibManager() *RuntimeJsLibManager {
	return c.runtimeJsLibManager
}

// InitBasedOnOptions does any initialization that is dependent on the
// compiler options.
//
// Java: Compiler.initBasedOnOptions()
func (c *Compiler) InitBasedOnOptions() {
	for k, v := range c.options.GetInputSourceMaps().All() {
		c.inputSourceMaps.Put(k, v)
	}
	// Create the source map if necessary.
	if c.options.ShouldGatherSourceMapInfo() {
		c.sourceMap = c.options.GetSourceMapFormat().GetInstance()
		c.sourceMap.SetPrefixMappings(c.options.GetSourceMapLocationMappings())
		// Go-only: fast mode encodes the mappings on workers (byte-identical).
		if c.options.isFastMode() && c.options.GetNumParallelThreads() > 1 && c.sourceMap.generatorV3 != nil {
			c.sourceMap.generatorV3.SetParallelWorkers(c.options.GetNumParallelThreads())
		}
		if c.options.GetApplyInputSourceMaps() {
			c.sourceMap.SetSourceFileMapping(c.AsSourceFileMapping())
			if c.options.GetSourceMapIncludeSourcesContent() {
				for inputSourceMap := range c.inputSourceMaps.Values() {
					c.addSourceMapSourceFiles(inputSourceMap)
				}
			}
		}
	}
}

// Java: Compiler.EMPTY_CHUNK_LIST_ERROR
var compiler_EMPTY_CHUNK_LIST_ERROR = DiagnosticType_Error("JSC_EMPTY_CHUNK_LIST_ERROR", "At least one chunk must be provided")

// Java: Compiler.EMPTY_ROOT_CHUNK_ERROR
var compiler_EMPTY_ROOT_CHUNK_ERROR = DiagnosticType_Error(
	"JSC_EMPTY_ROOT_CHUNK_ERROR",
	"Root chunk ''{0}'' must contain at least one source code input")

// checkFirstChunk verifies that at least one chunk has been provided and
// that the first one has at least one source code input.
//
// Java: Compiler.checkFirstChunk(List<JSChunk>)
func (c *Compiler) checkFirstChunk(chunks []*JSChunk) {
	if len(chunks) == 0 {
		c.Report(JSError_Make(compiler_EMPTY_CHUNK_LIST_ERROR))
	} else if len(chunks[0].GetInputs()) == 0 && len(chunks) > 1 {
		// The root chunk may only be empty if there is exactly 1 chunk.
		c.Report(JSError_Make(compiler_EMPTY_ROOT_CHUNK_ERROR, chunks[0].GetName()))
	}
}

// Compiler_JoinPathParts creates an OS specific path string from parts
//
// Java: Compiler.joinPathParts(String...)
func Compiler_JoinPathParts(pathParts ...string) string {
	return compiler_pathJoiner(pathParts)
}

// fillEmptyChunks fills any empty chunks with a place holder file. It makes
// any cross chunk motion easier.
//
// Java: Compiler.fillEmptyChunks(Iterable<JSChunk>)
func (c *Compiler) fillEmptyChunks(chunks []*JSChunk) {
	for _, chunk := range chunks {
		if chunk.GetName() != JSChunk_WEAK_CHUNK_NAME && len(chunk.GetInputs()) == 0 {
			input := NewCompilerInput(SourceFile_FromCode(AbstractCompiler_CreateFillFileName(chunk.GetName()), ""))
			input.SetCompiler(c)
			chunk.Add(input)
		}
	}
}

// RebuildInputsFromChunks rebuilds the internal input map by iterating
// over all chunks. This is necessary if inputs have been added to or
// removed from a chunk after an init or initChunks call.
//
// Java: Compiler.rebuildInputsFromChunks()
func (c *Compiler) RebuildInputsFromChunks() {
	c.InitInputsByIdMap()
}

// Java: Compiler.DUPLICATE_INPUT
var Compiler_DUPLICATE_INPUT = DiagnosticType_Error("JSC_DUPLICATE_INPUT", "Duplicate input: {0}")

// Java: Compiler.DUPLICATE_EXTERN_INPUT
var Compiler_DUPLICATE_EXTERN_INPUT = DiagnosticType_Error("JSC_DUPLICATE_EXTERN_INPUT", "Duplicate extern input: {0}")

// InitInputsByIdMap creates a map to make looking up an input by name fast.
// Also checks for duplicate inputs.
//
// Java: Compiler.initInputsByIdMap()
func (c *Compiler) InitInputsByIdMap() {
	clear(c.inputsById)
	for _, input := range c.externs {
		previous := c.putCompilerInput(input)
		if previous != nil {
			c.Report(JSError_Make(Compiler_DUPLICATE_EXTERN_INPUT, input.GetName()))
		}
	}

	for _, input := range c.chunkGraph.GetAllInputs() {
		previous := c.putCompilerInput(input)
		if previous != nil {
			c.Report(JSError_Make(Compiler_DUPLICATE_INPUT, input.GetName()))
		}
	}
}

// initAST sets up the skeleton of the AST (the externs and root).
//
// Java: Compiler.initAST()
func (c *Compiler) initAST() {
	c.jsRoot = rhino.IR_Root()
	c.externsRoot = rhino.IR_Root()
	c.externAndJsRoot = rhino.IR_Root(c.externsRoot, c.jsRoot)
}

// Java: Compiler.createAstFactory()
func (c *Compiler) CreateAstFactory() *AstFactory {
	if c.HasTypeCheckingRun() {
		if c.HasOptimizationColors() {
			return AstFactory_CreateFactoryWithColors(c.stage, c.GetColorRegistry(), c.runtimeJsLibManager)
		}
		return AstFactory_CreateFactoryWithTypes(c.stage, c.GetTypeRegistry(), c.runtimeJsLibManager)
	}
	return AstFactory_CreateFactoryWithoutTypes(c.stage, c.runtimeJsLibManager)
}

// Java: Compiler.createAstFactoryWithoutTypes()
func (c *Compiler) CreateAstFactoryWithoutTypes() *AstFactory {
	return AstFactory_CreateFactoryWithoutTypes(c.stage, c.runtimeJsLibManager)
}

// Java: Compiler.getAstAnalyzer()
func (c *Compiler) GetAstAnalyzer() *AstAnalyzer {
	options := c.GetOptions()
	useTypes := options.ShouldUseTypesForLocalOptimization()
	assumeGetters := options.GetAssumeGettersArePure()
	hasRegexp := c.HasRegExpGlobalReferences()
	assumeBuiltins := options.GetAssumePropertiesAreStaticallyAnalyzable()
	accessorSummary := c.GetAccessorSummary()
	if a := c.cachedAstAnalyzer; a != nil &&
		a.typeRegistry == c.typeRegistry &&
		a.accessorSummary == accessorSummary &&
		a.useTypesForLocalOptimization == useTypes &&
		a.assumeGettersArePure == assumeGetters &&
		a.hasRegexpGlobalReferences == hasRegexp &&
		a.assumeKnownBuiltinsArePure == assumeBuiltins {
		return a
	}
	analysisSettings := AstAnalyzer_Options_Builder_().
		SetUseTypesForLocalOptimization(useTypes).
		SetAssumeGettersArePure(assumeGetters).
		SetHasRegexpGlobalReferences(hasRegexp).
		SetAssumeKnownBuiltinsArePure(assumeBuiltins).
		Build()
	a := NewAstAnalyzer(analysisSettings, c.typeRegistry, accessorSummary)
	c.cachedAstAnalyzer = a
	return a
}

// Java: Compiler.createDefaultExpressionDecomposer()
func (c *Compiler) CreateDefaultExpressionDecomposer() *ExpressionDecomposer {
	return c.CreateExpressionDecomposer(
		c.GetUniqueNameIdSupplier(),
		nil,
		Scope_CreateGlobalScope(rhino.NewNode(rhino.Token_SCRIPT)))
}

// Java: Compiler.createExpressionDecomposer(Supplier<String>, ImmutableSet<String>, Scope)
func (c *Compiler) CreateExpressionDecomposer(
	uniqueNameIdSupplier func() string,
	knownConstantFunctions map[string]bool,
	scope *Scope) *ExpressionDecomposer {
	// If the output is ES5, then it may end up running on IE11, so enable a workaround
	// for one of its bugs.
	brokenIE11LocationAssign := parser.FeatureSet_ES5.Contains(c.GetOptions().GetOutputFeatureSet())
	return NewExpressionDecomposer(
		c, uniqueNameIdSupplier, knownConstantFunctions, scope, brokenIE11LocationAssign)
}

// Java: Compiler.isDebugLoggingEnabled()
func (c *Compiler) IsDebugLoggingEnabled() bool {
	return c.GetOptions().GetDebugLogDirectory() != nil
}

// Java: Compiler.getDebugLogFilterList()
func (c *Compiler) GetDebugLogFilterList() []string {
	if c.GetOptions().GetDebugLogFilter() == nil {
		return []string{}
	}
	var out []string
	for _, s := range strings.Split(*c.GetOptions().GetDebugLogFilter(), ",") {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// Java: Compiler.createOrReopenLog(Class<?>, String, String...)
func (c *Compiler) CreateOrReopenLog(owner string, firstNamePart string, restNameParts ...string) diagnostic.LogFile {
	if !c.IsDebugLoggingEnabled() {
		return diagnostic.LogFile_CreateNoOp()
	}

	dir := c.GetOptions().GetDebugLogDirectory()
	relativeParts := jutil.Path_Of(firstNamePart, restNameParts...)
	file := dir.ResolveString(owner).Resolve(relativeParts)

	// If a filter list for log file names was provided, only create a log file if any
	// of the filter strings matches.
	filters := c.GetDebugLogFilterList()
	if len(filters) == 0 {
		return diagnostic.LogFile_CreateOrReopen(file)
	}

	for _, filter := range filters {
		if strings.Contains(file.String(), filter) {
			return diagnostic.LogFile_CreateOrReopen(file)
		}
	}
	return diagnostic.LogFile_CreateNoOp()
}

// Java: Compiler.createOrReopenIndexedLog(Class<?>, String, String...)
func (c *Compiler) CreateOrReopenIndexedLog(owner string, firstNamePart string, restNameParts ...string) diagnostic.LogFile {
	jutil.CheckStateMsg1(c.currentPassIndex >= 0, "%s", c.currentPassIndex)

	index := strconv.Itoa(c.currentPassIndex)
	for len(index) < 3 {
		index = "0" + index
	}
	length := len(restNameParts)
	if length == 0 {
		firstNamePart = index + "_" + firstNamePart
	} else {
		restNameParts[length-1] = index + "_" + restNameParts[length-1]
	}

	return c.CreateOrReopenLog(owner, firstNamePart, restNameParts...)
}

// Java: Compiler.maybePrintSourceAfterEachPass(String)
func (c *Compiler) maybePrintSourceAfterEachPass(passName string) {
	if !c.options.ShouldPrintSourceAfterEachPass() {
		return
	}

	currentJsSource := c.GetCurrentJsSource()
	if c.lastJsSource != nil && currentJsSource == *c.lastJsSource {
		return
	}

	if c.IsDebugLoggingEnabled() {
		func() {
			log := c.CreateOrReopenIndexedLog("Compiler", "source_after_pass", passName)
			defer log.Close()
			log.LogString(currentJsSource)
		}()
	} else {
		err := jutil.System_Err()
		err.PrintlnEmpty()
		if c.debugMessage != nil {
			err.Println("// DEBUG: " + *c.debugMessage)
		}
		err.Println("// " + passName + " yields:")
		err.Println("// ************************************")
		err.Println(currentJsSource)
		err.Println("// ************************************")
	}

	c.lastJsSource = &currentJsSource
}

// Java: Compiler.getCurrentJsSource()
func (c *Compiler) GetCurrentJsSource() string {
	c.ResetAndIntitializeSourceMap()

	fileNameRegexList := c.options.GetFilesToPrintAfterEachPassRegexList()
	chunkNameRegexList := c.options.GetChunksToPrintAfterEachPassRegexList()
	qnameSet := jutil.NewLinkedHashSet[string](c.options.GetQnameUsesToPrintAfterEachPassList()...)
	builder := jutil.NewStringBuilder()

	if len(fileNameRegexList) == 0 && len(chunkNameRegexList) == 0 && qnameSet.IsEmpty() {
		return c.ToSource()
	}
	if len(fileNameRegexList) != 0 {
		rhino.CheckNotNull(c.externsRoot)
		rhino.CheckNotNull(c.jsRoot)
		for _, r := range []rhino.Node{c.externsRoot, c.jsRoot} {
			for fileNode := r.GetFirstChild(); !fileNode.IsNil(); fileNode = fileNode.GetNext() {
				fileName := fileNode.GetSourceFileName()
				for _, regex := range fileNameRegexList {
					if jutil.Pattern_Compile(regex).Matches(*fileName) {
						source := "// " + *fileName + "\n" + c.ToSourceOfNode(fileNode)
						builder.Append(source)
						break
					}
				}
			}
		}
		if builder.Length() == 0 {
			builder.Append("// No files matched any of: ").Append(jutil.CollectionToString(fileNameRegexList))
		}
	}
	if len(chunkNameRegexList) != 0 {
		chunks := jutil.CheckNotNull(c.GetChunks())
		unmatchedChunkNames := []string{}
		// As we are emitting all chunks at once in dependency order, we can be smarter about
		// deciding which licenses to emit where in the source.
		lt := NewCompiler_ChunkGraphAwareLicenseTracker(c)
		for _, jsChunk := range chunks {
			chunkName := jsChunk.GetName()
			lt.SetCurrentChunkContext(jsChunk)
			for _, regex := range chunkNameRegexList {
				if jutil.Pattern_Compile(regex).Matches(chunkName) {
					source := "// module '" + chunkName + "'\n" + c.ToSourceOfLicenseTrackerAndJSChunk(lt, jsChunk)
					builder.Append(source)
					break
				}
			}
			unmatchedChunkNames = append(unmatchedChunkNames, chunkName)
		}
		if builder.Length() == 0 {
			builder.Append("// No chunks were matched:\n")
			for _, chunkName := range unmatchedChunkNames {
				builder.Append("// ").Append(chunkName).Append("\n")
			}
		}
	}
	if !qnameSet.IsEmpty() {
		// Keep track of how a qualified name we're interested in gets renamed.
		// Note that a qname can actually be renamed multiple times during compilation,
		// so we need a multimap here.
		originalToNewQNameMap := jutil.NewSetMultimap[string, string]()
		// Get a stream of top level statement nodes that contain at least one reference to one of the
		// qnames we are interested in. Take note of relevant renamings while we're at it.
		//
		// Go-only: language-adapter — Java's lazy Stream.filter runs when
		// the stream is consumed (after the header below is written), so
		// the filter is applied in the loop at the end.
		statements := c.getTopLevelStatements(c.jsRoot)
		statementFilter := func(statement rhino.Node) bool {
			// filter to just those statements containing interesting qnames
			return NodeUtil_Has(
				statement,
				func(n rhino.Node) bool {
					rhino.CheckNotNull(n)
					if !n.IsQualifiedName() {
						return false
					}
					qname := n.GetQualifiedName()
					if qnameSet.Contains(*qname) {
						return true
					}
					originalQname := n.GetOriginalQualifiedName()
					if originalQname != nil && qnameSet.Contains(*originalQname) {
						originalToNewQNameMap.Put(*originalQname, *qname)
						return true
					} else {
						return false
					}
				},
				func(node rhino.Node) bool { return true })
		}
		builder.Append("//\n")
		builder.Append("// closure-compiler: Printing all of the top-level statements\n")
		builder.Append("// that contain references to these qualified names.\n")
		builder.Append("//\n")
		for qname := range qnameSet.All() {
			builder.Append(fmt.Sprintf("// '%s'\n", qname))
			for newName := range originalToNewQNameMap.Get(qname).All() {
				builder.Append(fmt.Sprintf("// '%s' (originally '%s')\n", newName, qname))
			}
		}
		builder.Append("//\n")
		for _, statement := range statements {
			if !statementFilter(statement) {
				continue
			}
			builder.
				Append(fmt.Sprintf("// %s\n", statement.GetLocation())).
				Append(c.ToSourceOfNode(statement)).
				Append("\n")
		}
	}
	return builder.String()
}

// Java: Compiler.getTopLevelStatements(Node)
func (c *Compiler) getTopLevelStatements(root rhino.Node) []rhino.Node {
	cb := &compiler_topLevelStatementsCallback{}
	NodeTraversal_Traverse(c, root, cb)
	return cb.builder
}

// compiler_topLevelStatementsCallback is the anonymous
// AbstractPreOrderCallback of getTopLevelStatements: it collects the
// outermost statements.
type compiler_topLevelStatementsCallback struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	builder []rhino.Node
}

// Java: Compiler.getTopLevelStatements(...).new AbstractPreOrderCallback() {...}.shouldTraverse(NodeTraversal, Node, Node)
func (cb *compiler_topLevelStatementsCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if NodeUtil_IsStatement(n) {
		cb.builder = append(cb.builder, n)
		return false
	} else {
		return true
	}
}

// Java: Compiler.getScriptNode(String)
func (c *Compiler) GetScriptNode(filename string) rhino.Node {
	return c.scriptNodeByFilename.Get(jutil.CheckNotNull(filename))
}

// NewTracer returns a new tracer for the given pass name.
//
// Java: Compiler.newTracer(String)
func (c *Compiler) NewTracer(passName string) *Tracer {
	comment := passName
	if c.changeTracker.GetRecentChange().HasCodeChanged() {
		comment += " on recently changed AST"
	}
	if c.options.GetTracerMode().IsOn() && c.tracker != nil {
		c.tracker.RecordPassStart(passName, true)
	}
	if DebugPassLogHook != nil {
		DebugPassLogHook(DebugPassEvent{Name: passName, Tracer: true, Start: true})
	}
	typ := "Compiler"
	return NewTracerWithType(&typ, &comment)
}

// Java: Compiler.stopTracer(Tracer, String)
func (c *Compiler) StopTracer(t *Tracer, passName string) {
	result := t.Stop()
	if DebugPassLogHook != nil {
		DebugPassLogHook(DebugPassEvent{Name: passName, Tracer: true})
	}
	if c.options.GetTracerMode().IsOn() && c.tracker != nil {
		c.tracker.RecordPassStop(passName, result)
	}
}

// GetResult returns the result of the compilation.
//
// Java: Compiler.getResult()
func (c *Compiler) GetResult() *Result {
	return NewResultWithAllMaps(
		c.GetErrors(),
		c.GetWarnings(),
		c.variableMap,
		c.propertyMap,
		nil,
		c.stringMap,
		c.instrumentationMapping,
		c.sourceMap,
		c.externExports,
		c.cssNames,
		c.idGeneratorMap,
		c.transpiledFiles)
}

// SetDefaultCodingConvention overrides the compiler's fallback CodingConvention
// used when CompilerOptions does not specify one.
func (c *Compiler) SetDefaultCodingConvention(cc CodingConvention) {
	c.defaultCodingConvention = jutil.CheckNotNull(cc)
}

// GetDefaultCodingConvention returns the compiler's fallback CodingConvention.
func (c *Compiler) GetDefaultCodingConvention() CodingConvention {
	return c.defaultCodingConvention
}

// GetIdGeneratorMap returns the serialized ID generator mapping, or "" when
// none has been produced, without allocating a full Result.
func (c *Compiler) GetIdGeneratorMap() string {
	if c.idGeneratorMap == nil {
		return ""
	}
	return *c.idGeneratorMap
}

// GetConformanceViolations returns the attributed conformance violation protos
// persisted on this Compiler across compilation stages.
func (c *Compiler) GetConformanceViolations() []*ConformanceViolationsProto {
	return c.conformanceViolations
}

// SetConformanceViolations replaces the attributed conformance violation protos
// stored on this Compiler.
func (c *Compiler) SetConformanceViolations(violations []*ConformanceViolationsProto) {
	c.conformanceViolations = slices.Clone(violations)
}

// AddConformanceViolations appends an attributed conformance violations proto
// to this Compiler.
func (c *Compiler) AddConformanceViolations(violations *ConformanceViolationsProto) {
	c.conformanceViolations = append(c.conformanceViolations, jutil.CheckNotNull(violations))
}

// GetUnattributedConformanceViolations returns the unattributed conformance
// violation protos persisted on this Compiler across compilation stages.
func (c *Compiler) GetUnattributedConformanceViolations() []*ViolationProto {
	return c.unattributedConformanceViolations
}

// SetUnattributedConformanceViolations replaces the unattributed conformance
// violation protos stored on this Compiler.
func (c *Compiler) SetUnattributedConformanceViolations(violations []*ViolationProto) {
	c.unattributedConformanceViolations = slices.Clone(violations)
}

// AddUnattributedConformanceViolation appends an unattributed conformance
// violation proto to this Compiler.
func (c *Compiler) AddUnattributedConformanceViolation(violation *ViolationProto) {
	c.unattributedConformanceViolations = append(c.unattributedConformanceViolations, jutil.CheckNotNull(violation))
}

// GetErrors returns the list of errors (never null).
//
// Java: Compiler.getErrors()
func (c *Compiler) GetErrors() []*JSError {
	if c.errorManager == nil {
		return []*JSError{}
	}
	return c.errorManager.GetErrors()
}

// GetWarnings returns the list of warnings (never null).
//
// Java: Compiler.getWarnings()
func (c *Compiler) GetWarnings() []*JSError {
	if c.errorManager == nil {
		return []*JSError{}
	}
	return c.errorManager.GetWarnings()
}

// Java: Compiler.getRoot()
func (c *Compiler) GetRoot() rhino.Node {
	return c.externAndJsRoot
}

// Java: Compiler.getAllowableFeatures()
func (c *Compiler) GetAllowableFeatures() *parser.FeatureSet {
	return c.allowableFeatures
}

// Java: Compiler.setAllowableFeatures(FeatureSet)
func (c *Compiler) SetAllowableFeatures(allowableFeatures *parser.FeatureSet) {
	c.allowableFeatures = allowableFeatures
}

// Java: Compiler.markFeatureNotAllowed(Feature)
func (c *Compiler) MarkFeatureNotAllowed(feature parser.FeatureSet_Feature) {
	c.allowableFeatures = c.allowableFeatures.WithoutFeatures(feature)
}

// Java: Compiler.markFeatureSetNotAllowed(FeatureSet)
func (c *Compiler) MarkFeatureSetNotAllowed(featureSet *parser.FeatureSet) {
	c.allowableFeatures = c.allowableFeatures.Without(featureSet)
}

// Java: Compiler.getUniqueIdSupplier()
func (c *Compiler) GetUniqueIdSupplier() *UniqueIdSupplier {
	return c.uniqueIdSupplier
}

// nextUniqueNameId creates a new id for making unique names.
//
// Deprecated.
//
// Java: Compiler.nextUniqueNameId()
func (c *Compiler) nextUniqueNameId() int {
	id := c.uniqueNameId
	c.uniqueNameId++
	return id
}

// GetUniqueNameIdSupplier is the legacy supplier for getting unique names.
//
// Deprecated: This is deprecated because this supplier fails to generate
// unique Ids across input files when used for tagged template literal
// transpilation. Kindly use the new UniqueIdSupplier instead.
//
// Java: Compiler.getUniqueNameIdSupplier()
func (c *Compiler) GetUniqueNameIdSupplier() func() string {
	if c.uniqueNameIdSupplier == nil {
		c.uniqueNameIdSupplier = func() string { return nodeUtilFormatInt(c.nextUniqueNameId()) }
	}
	return c.uniqueNameIdSupplier
}

// Java: Compiler.areNodesEqualForInlining(Node, Node)
func (c *Compiler) AreNodesEqualForInlining(n1 rhino.Node, n2 rhino.Node) bool {
	if c.options.ShouldAmbiguateProperties() || c.options.ShouldDisambiguateProperties() {
		// The type based optimizations require that type information is preserved
		// during other optimizations.
		return n1.IsEquivalentToTyped(n2)
	} else {
		return n1.IsEquivalentTo(n2)
	}
}

// ------------------------------------------------------------------------
// Inputs
// ------------------------------------------------------------------------

// TODO(nicksantos): Decide which parts of these belong in an AbstractCompiler
// interface, and which ones should always be injected.

// Java: Compiler.getInput(InputId)
func (c *Compiler) GetInput(id *rhino.InputId) *CompilerInput {
	// TODO(bradfordcsmith): Allowing null id is less ideal. Add checkNotNull(id) here and fix
	// call sites that break.
	if id == nil {
		return nil
	}
	return c.inputsById[id.GetIdName()]
}

// PutCompilerInput registers input by its InputId on c.
//
// Java: Compiler.putCompilerInput(CompilerInput)
func (c *Compiler) PutCompilerInput(input *CompilerInput) *CompilerInput {
	return c.putCompilerInput(input)
}

func (c *Compiler) putCompilerInput(input *CompilerInput) *CompilerInput {
	input.SetCompiler(c)
	idName := input.GetInputId().GetIdName()
	prev := c.inputsById[idName]
	c.inputsById[idName] = input
	return prev
}

// GetChunkGraph gets the graph of JS source chunks.
//
// Returns null if #init or #initChunks hasn't been called yet. Otherwise,
// the result is always a chunk graph, even in the degenerate case where
// there's only one chunk.
//
// Java: Compiler.getChunkGraph()
func (c *Compiler) GetChunkGraph() *JSChunkGraph {
	return c.chunkGraph
}

// GetChunks gets the JS source chunks in dependency order.
//
// Returns null if #init or #initChunks hasn't been called yet. Otherwise,
// the result is always non-empty, even in the degenerate case where
// there's only one chunk.
//
// Java: Compiler.getChunks()
func (c *Compiler) GetChunks() []*JSChunk {
	if c.chunkGraph != nil {
		return c.chunkGraph.GetAllChunks()
	}
	return nil
}

// Java: Compiler.clearJSTypeRegistry()
func (c *Compiler) ClearJSTypeRegistry() {
	if c.typeRegistry != nil {
		c.clearedTypeRegistry = weak.Make(c.typeRegistry)
		// Drop the type column: every type a Node referred to is released
		// here, whether or not the node is still reachable.
		c.typeRegistry.ReleaseNodeTypeHandles()
	}
	c.typeRegistry = nil
	c.typeValidator = nil
	c.abstractInterpreter = nil
	// Go-only: the cached analyzer holds the registry it was built with.
	c.cachedAstAnalyzer = nil
	if AfterTypesToColorsHook != nil {
		AfterTypesToColorsHook()
	}
}

// AfterTypesToColorsHook, when set, runs at the end of ClearJSTypeRegistry.
//
// Go-only: measurement aid (cmd/jsgo records the live heap there).
var AfterTypesToColorsHook func()

// typeRegistryReleased reports whether the registry that ClearJSTypeRegistry
// dropped has been garbage collected. Every JSType points to its registry,
// so this is a strict check: it fails as long as any JSType at all is still
// reachable (for example from a Node, live or detached, that still carries
// it), even if the registry emptied its tables and the rest of the type
// graph was freed.
//
// Go-only: debugging aid for tests.
func (c *Compiler) typeRegistryReleased() bool {
	runtime.GC()
	runtime.GC()
	return c.clearedTypeRegistry.Value() == nil
}

// TypeRegistryReleased is typeRegistryReleased for cmd/jsgo's debugging
// hook (Go-only).
func (c *Compiler) TypeRegistryReleased() bool { return c.typeRegistryReleased() }

// Java: Compiler.isTypeRegistryCleared()
func (c *Compiler) IsTypeRegistryCleared() bool {
	return c.typeCheckingHasRun && c.typeRegistry == nil
}

// Java: Compiler.getTypeRegistry()
func (c *Compiler) GetTypeRegistry() *rhino.JSTypeRegistry {
	if c.typeRegistry == nil {
		jutil.CheckStateMsg(
			!c.HasTypeCheckingRun(),
			// This could throw when calling "getTypeRegistry()" during the optimizations phase when
			// JSTypes have been converted to optimization colors
			"Attempted to re-initialize JSTypeRegistry after it had been cleared")
		c.typeRegistry = rhino.NewJSTypeRegistryWithForwardDeclaredTypes(c.oldErrorReporter, c.forwardDeclaredTypes)
	}
	return c.typeRegistry
}

// Java: Compiler.getColorRegistry()
func (c *Compiler) GetColorRegistry() *colors.ColorRegistry {
	return jutil.CheckNotNullMsg(c.colorRegistry, "Color registry has not been initialized yet")
}

// Java: Compiler.setColorRegistry(ColorRegistry)
func (c *Compiler) SetColorRegistry(colorRegistry *colors.ColorRegistry) {
	jutil.CheckState(c.runtimeLibraryTypedAsts != nil)
	c.colorRegistry = colorRegistry
}

// Java: Compiler.forwardDeclareType(String)
func (c *Compiler) ForwardDeclareType(typeName string) {
	c.forwardDeclaredTypes.Add(typeName)
}

// Java: Compiler.setTypeCheckingHasRun(boolean)
func (c *Compiler) SetTypeCheckingHasRun(hasRun bool) {
	c.typeCheckingHasRun = hasRun
}

// Java: Compiler.hasTypeCheckingRun()
func (c *Compiler) HasTypeCheckingRun() bool {
	return c.typeCheckingHasRun
}

// Java: Compiler.hasOptimizationColors()
func (c *Compiler) HasOptimizationColors() bool {
	return c.colorRegistry != nil
}

// Java: Compiler.getTypedScopeCreator()
func (c *Compiler) GetTypedScopeCreator() ScopeCreator {
	if c.typedScopeCreator == nil {
		jutil.CheckStateMsg(
			!c.HasTypeCheckingRun(),
			// This could throw when calling "getTypedScopeCreator()" during the optimizations phase
			// when JSTypes have been converted to optimization colors
			"Attempted to re-initialize TypedScopeCreator after it had been cleared")
		c.typedScopeCreator = NewTypedScopeCreator(c)
	}
	return c.typedScopeCreator
}

// Java: Compiler.clearTypedScopeCreator()
func (c *Compiler) ClearTypedScopeCreator() {
	if c.typedScopeCreator != nil {
		c.typedScopeCreator.memoized = nil
		c.typedScopeCreator = nil
	}
}

// Java: Compiler.ensureDefaultPassConfig()
func (c *Compiler) EnsureDefaultPassConfig() *DefaultPassConfig {
	passes := c.GetPassConfig().GetBasePassConfig()
	defaultPassConfig, ok := passes.(*DefaultPassConfig)
	jutil.CheckStateMsg(
		ok,
		"PassConfigs must eventually delegate to the DefaultPassConfig")
	return defaultPassConfig
}

// Java: Compiler.buildKnownSymbolTable()
func (c *Compiler) BuildKnownSymbolTable() *SymbolTable {
	symbolTable := NewSymbolTable(c, c.GetTypeRegistry())

	if c.typedScopeCreator != nil {
		SymbolTable_AddScopes(symbolTable, c.typedScopeCreator.GetAllMemoizedScopes())
		SymbolTable_AddSymbolsFrom[*TypedVar, *TypedVar](symbolTable, c.typedScopeCreator)
	} else {
		symbolTable.FindScopes(c.externsRoot, c.jsRoot)
	}

	globalNamespace := NewGlobalNamespaceWithRoots(c, c.externsRoot, c.jsRoot)
	SymbolTable_AddSymbolsFrom[*GlobalNamespace_Name, *GlobalNamespace_Ref](symbolTable, globalNamespace)

	refCollector := NewReferenceCollector(
		c, ReferenceCollector_DO_NOTHING_BEHAVIOR, NewSyntacticScopeCreator(c))
	refCollector.ProcessRoot(c.GetRoot())
	SymbolTable_AddSymbolsFrom(symbolTable, refCollector.AsStaticSymbolTable())

	preprocessorSymbolTable := c.EnsureDefaultPassConfig().GetPreprocessorSymbolTable()
	if preprocessorSymbolTable != nil {
		SymbolTable_AddSymbolsFrom[rhino.SimpleSlot, *PreprocessorSymbolTable_Reference](symbolTable, preprocessorSymbolTable)
	}

	symbolTable.FlattenGoogModuleExports()
	symbolTable.FillNamespaceReferences()
	symbolTable.FillPropertyScopes()
	symbolTable.FillThisReferences(c.externsRoot, c.jsRoot)
	symbolTable.FillPropertySymbols(c.externsRoot, c.jsRoot)
	symbolTable.FillSuperReferences(c.externsRoot, c.jsRoot)
	symbolTable.FillJSDocInfo(c.externsRoot, c.jsRoot)
	symbolTable.FillSymbolVisibility(c.externsRoot, c.jsRoot)
	symbolTable.FillGoogProvideModuleRequires(c.externsRoot, c.jsRoot)
	symbolTable.RemoveGeneratedSymbols()

	return symbolTable
}

// Java: Compiler.getTopScope()
func (c *Compiler) GetTopScope() *TypedScope {
	return c.topScope
}

// Java: Compiler.setTopScope(TypedScope)
func (c *Compiler) SetTopScope(x *TypedScope) {
	jutil.CheckStateMsg1(x == nil || x.GetParent() == nil, "%s", x)
	c.topScope = x
}

// Java: Compiler.getReverseAbstractInterpreter()
func (c *Compiler) GetReverseAbstractInterpreter() ReverseAbstractInterpreter {
	if c.abstractInterpreter == nil {
		var interpreter ChainableReverseAbstractInterpreter = NewSemanticReverseAbstractInterpreter(c.GetTypeRegistry())
		if c.options.GetClosurePass() {
			interpreter = NewClosureReverseAbstractInterpreter(c.GetTypeRegistry()).Append(interpreter).GetFirst()
		}
		c.abstractInterpreter = interpreter
	}
	return c.abstractInterpreter
}

// Java: Compiler.getTypeValidator()
func (c *Compiler) GetTypeValidator() *TypeValidator {
	if c.typeValidator == nil {
		jutil.CheckStateMsg(
			!c.HasTypeCheckingRun(),
			// This could throw when calling "getTypeValidator()" during the optimizations phase when
			// JSTypes have been converted to optimization colors
			"Attempted to re-initialize TypeValidator after it had been cleared")
		c.typeValidator = NewTypeValidator(c)
	}
	return c.typeValidator
}

// Java: Compiler.getTypeMismatches()
func (c *Compiler) GetTypeMismatches() []*TypeMismatch {
	if c.typeCheckingHasRun {
		return c.GetTypeValidator().GetMismatches()
	}
	panic(jutil.NewRuntimeException("Can't ask for type mismatches before type checking."))
}

// Java: Compiler.getTranspilationNamespace()
func (c *Compiler) GetTranspilationNamespace() rhino.StaticScope {
	if c.transpilationNamespace == nil {
		// Note: this also needs to happen /after/ the runtime libraries are injected, or we will not
		// be able to find the runtime library code.
		// TODO(bradfordcsmith): It would be good to add a check to ensure that library injection has
		// already happened. We used to check for normalization, but that's going too far and
		// prevents us moving transpilation after normalization.
		gn := NewGlobalNamespaceWithRoots(c, c.GetExternsRoot(), c.GetJsRoot())
		// Exclude user-defined code names, to save memory and help prevent misuse.
		gn.SetShouldTraverseScriptPredicate(
			func(script rhino.Node) bool { return script.IsFromExterns() || script.IsFirstChildOf(c.GetJsRoot()) })
		gn.EnsureGenerated() // ensure namespace is generated
		// wrap GlobalNamespace to prevent looking up extra information in GlobalNamespace
		c.transpilationNamespace = gn
	}

	return c.transpilationNamespace
}

// Java: Compiler.maybeSetTracker()
func (c *Compiler) MaybeSetTracker() {
	if !c.options.GetTracerMode().IsOn() || c.tracker != nil {
		return
	}

	c.tracker = NewPerformanceTracker(c.externsRoot, c.jsRoot, c.options.GetTracerMode())
	c.changeTracker.AddChangeHandler(c.tracker.GetCodeChangeHandler())
}

// Java: Compiler.initializeModuleLoader()
func (c *Compiler) InitializeModuleLoader() {

	var moduleResolverFactory ModuleLoader_ModuleResolverFactory
	switch c.options.GetModuleResolutionMode() {
	case ModuleLoader_ResolutionMode_BROWSER:
		moduleResolverFactory = BrowserModuleResolver_FACTORY
	case ModuleLoader_ResolutionMode_NODE:
		// processJsonInputs requires a module loader to already be defined
		// so we redefine it afterwards with the package.json inputs
		moduleResolverFactory = NewNodeModuleResolver_FactoryWithEntries(c.ProcessJsonInputs(c.chunkGraph.GetAllInputs()))
	case ModuleLoader_ResolutionMode_WEBPACK:
		moduleResolverFactory = NewWebpackModuleResolver_Factory(c.inputPathByWebpackId)
	case ModuleLoader_ResolutionMode_BROWSER_WITH_TRANSFORMED_PREFIXES:
		moduleResolverFactory = NewBrowserWithTransformedPrefixesModuleResolver_Factory(
			c.options.GetBrowserResolverPrefixReplacements())
	}

	inputs := c.chunkGraph.GetAllInputs()
	dependencyInfos := make([]DependencyInfo, len(inputs))
	for i, in := range inputs {
		dependencyInfos[i] = in
	}
	c.moduleLoader = ModuleLoader_Builder_().
		SetModuleRoots(c.options.GetModuleRoots()).
		SetInputs(dependencyInfos).
		SetFactory(moduleResolverFactory).
		SetPathResolver(ModuleLoader_PathResolver_RELATIVE).
		SetPathEscaper(c.options.GetPathEscaper()).
		Build()
}

var compiler_emptyLicenseSet = jutil.NewLinkedHashSet[string]()

// ComputeCFG is Control Flow Analysis.
//
// Java: Compiler.computeCFG()
func (c *Compiler) ComputeCFG() ControlFlowGraph[rhino.Node] {
	Compiler_logger.Fine("Computing Control Flow Graph")
	tracer := c.NewTracer("computeCFG")
	cfg := ControlFlowAnalysis_Builder_().
		SetCompiler(c).
		SetCfgRoot(c.jsRoot).
		SetTraverseFunctions(true).
		ComputeCfg()
	c.StopTracer(tracer, "computeCFG")
	return cfg
}

// Java: Compiler.SYNTHETIC_FILE_NAME_PREFIX
const compiler_SYNTHETIC_FILE_NAME_PREFIX = " [synthetic:"

// Java: Compiler.SYNTHETIC_CODE_INPUT_ID
var compiler_SYNTHETIC_CODE_INPUT_ID = rhino.NewInputId(compiler_SYNTHETIC_FILE_NAME_PREFIX + "input] ")

// Port of com/google/javascript/jscomp/Compiler.java @ closure-compiler v20261015

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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"iter"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unsafe"

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
	AbstractCompilerBase

	options *CompilerOptions

	passes PassConfig

	// The externs inputs
	externs *jutil.ArrayList[*CompilerInput]

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
	externsRoot     *rhino.Node
	jsRoot          *rhino.Node
	externAndJsRoot *rhino.Node

	// Used for debugging; to see the compiled code between passes
	lastJsSource *string

	// Used to track the invariant about what language features may be in the AST at the given point
	// in compilation.
	allowableFeatures *parser.FeatureSet

	// Go-only: language-adapter — a ConcurrentHashMap keyed by InputId
	// (value equality); never iterated.
	inputsById *jutil.JavaHashMap[*rhino.InputId, *CompilerInput]

	scriptNodeByFilename *jutil.LinkedHashMap[string, *rhino.Node]

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
	commentsPerFile *jutil.LinkedHashMap[string, []*parser.Comment]

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
	runtimeLibraryTypedAsts *jutil.LinkedHashMap[string, func() *rhino.Node]

	// Types that have been forward declared
	forwardDeclaredTypes compiler_forwardDeclaredTypes

	typeCheckingHasRun bool

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

	typedAstFilesystem          *jutil.LinkedHashMap[*SourceFile, func() *rhino.Node]
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

	indexProvidersByType *jutil.LinkedHashMap[reflect.Type, IndexProvider[any]]

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
	c.externs = jutil.NewArrayList[*CompilerInput]()
	c.inputsById = jutil.NewJavaHashMap[*rhino.InputId, *CompilerInput]()
	c.scriptNodeByFilename = jutil.NewLinkedHashMap[string, *rhino.Node]()
	c.sourceMapOriginalSources = jutil.NewLinkedHashMap[string, *SourceFile]()
	c.inputSourceMaps = jutil.NewJavaHashMap[string, *SourceMapInput]()
	c.commentsPerFile = jutil.NewLinkedHashMap[string, []*parser.Comment]()
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
	c.indexProvidersByType = jutil.NewLinkedHashMap[reflect.Type, IndexProvider[any]]()
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
	c := &Compiler{}
	c.InitAbstractCompilerBase(c)
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
			for g := range options.GetExtraReportGenerators().All() {
				add(g)
			}
			c.SetErrorManager(NewSortingErrorManager(builder))
		}
	}

	if options.GetSkipNonTranspilationPasses() {
		options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_NO_OP)
	}

	c.moduleLoader = ModuleLoader_EMPTY

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
		c.logToFile("externs.log", c.externs.String)
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
		err.Println(c.externs.String())
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
	log := c.CreateOrReopenLog(reflect.TypeOf(c).Elem(), logFileName)
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
func (c *Compiler) GetTypedAstDeserializer(file *SourceFile) func() *rhino.Node {
	if c.typedAstFilesystem == nil {
		return nil
	}

	delete(c.typedAstScriptDeserializers, file)
	ast := c.typedAstFilesystem.Remove(file)
	jutil.CheckStateMsg1(
		ast != nil ||
			strings.HasPrefix(file.GetName(), compiler_SYNTHETIC_FILE_NAME_PREFIX) ||
			AbstractCompiler_IsFillFileName(file.GetName()),
		"TypedAST filesystem initialized, but missing requested file: %s",
		file)
	return ast
}

func (c *Compiler) takeTypedAstDeserializer(file *SourceFile) (*ScriptNodeDeserializer, func() *rhino.Node) {
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

	c.typedAstFilesystem = astData.GetFilesystem()
	c.typedAstScriptDeserializers = astData.scriptDeserializers
	c.externProperties = astData.GetExternProperties()
	c.colorRegistry = astData.GetColorRegistry()
	c.SetTypeCheckingHasRun(deserializeTypes)

	for library := range astData.GetRuntimeLibraries().All() {
		c.runtimeJsLibManager.EnsureLibraryInjected(library, false)
	}

	c.GetSynthesizedExternsInput() // Force lazy creation.

	Compiler_RunInCompilerThread(c,
		func() any {
			// TODO(lharker): refactor things to avoid this misleading 'parse' call. It's not actually
			// parsing JS source code, just doing more deserialization & building the Rhino AST.
			c.parseForCompilationInternal()
			return nil
		})
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
	runtimeLibraryTypedAsts := jutil.NewLinkedHashMap[string, func() *rhino.Node]()
	for key, value := range astData.GetFilesystem().All() {
		runtimeLibraryTypedAsts.ComputeIfAbsent(
			key.GetName(), func(f string) func() *rhino.Node { return value })
	}

	c.runtimeLibraryTypedAsts = jutil.NewLinkedHashMapFrom(runtimeLibraryTypedAsts)
}

// Init initializes the instance state needed for a compile job.
//
// Java: Compiler.init(List<SourceFile>, List<SourceFile>, CompilerOptions)
func (c *Compiler) Init(externs []*SourceFile, sources []*SourceFile, options *CompilerOptions) {
	chunk := NewJSChunk(JSChunk_STRONG_CHUNK_NAME)
	for _, source := range sources {
		chunk.Add(NewCompilerInputWithIsExtern(source /* isExtern= */, false))
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

	c.externs.Clear()
	for _, file := range externs {
		c.externs.Add(NewCompilerInputWithIsExtern(file /* isExtern= */, true))
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

	c.commentsPerFile = jutil.NewLinkedHashMap[string, []*parser.Comment]()
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
		chunkGraphLog := c.CreateOrReopenLog(reflect.TypeOf(c).Elem(), "chunk_graph.dot")
		defer chunkGraphLog.Close()
		if chunkGraphLog.IsLogging() {
			chunkGraphLog.LogString(DotFormatter_ToDotOfGraphvizGraph(c.chunkGraph.ToGraphvizGraph()))
		}
	}()

	var nodeForCodeInsertion func() *rhino.Node
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
	c.inputsById.Clear()
	for input := range c.externs.All() {
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

// CompileOfSourceFiles compiles a single source file and a single externs
// file.
//
// Java: Compiler.compile(SourceFile, SourceFile, CompilerOptions)
func (c *Compiler) CompileOfSourceFiles(extern *SourceFile, input *SourceFile, options *CompilerOptions) *Result {
	return c.Compile([]*SourceFile{extern}, []*SourceFile{input}, options)
}

// Compile compiles a list of inputs.
//
// This is a convenience method to wrap up all the work of compilation,
// including generating the error and warning report.
//
// NOTE: All methods called here must be public, because client code must
// be able to replicate and customize this.
//
// Java: Compiler.compile(List<SourceFile>, List<SourceFile>, CompilerOptions)
func (c *Compiler) Compile(externs []*SourceFile, inputs []*SourceFile, options *CompilerOptions) *Result {
	// The compile method should only be called once.
	jutil.CheckState(c.jsRoot == nil)

	func() {
		defer c.GenerateReport()
		c.Init(externs, inputs, options)
		if !c.HasErrors() {
			c.ParseForCompilation()
		}
		if !c.HasErrors() {
			if options.GetInstrumentForCoverageOnly() {
				// TODO(bradfordcsmith): The option to instrument for coverage only should belong to the
				//     runner, not the compiler.
				c.InstrumentForCoverage()
			} else {
				c.Stage1Passes()
				if !c.HasErrors() {
					c.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
					if !c.HasErrors() {
						c.Stage3Passes()
					}
				}
			}
			c.PerformPostCompilationTasks()
		}
	}()
	return c.GetResult()
}

// GenerateReport generates a report of all warnings and errors found during
// compilation to stderr.
//
// Client code must call this method explicitly if it doesn't use one of
// the convenience methods that do so automatically.
//
// Always call this method, even if the compiler throws an exception. The
// report will include information about the exception.
//
// Java: Compiler.generateReport()
func (c *Compiler) GenerateReport() {
	t := c.NewTracer("generateReport")
	c.errorManager.GenerateReport()
	c.StopTracer(t, "generateReport")
}

// CompileChunks compiles a list of chunks.
//
// This is a convenience method to wrap up all the work of compilation,
// including generating the error and warning report.
//
// NOTE: All methods called here must be public, because client code must
// be able to replicate and customize this.
//
// Java: Compiler.compileChunks(List<SourceFile>, List<JSChunk>, CompilerOptions)
func (c *Compiler) CompileChunks(externs []*SourceFile, chunks []*JSChunk, options *CompilerOptions) *Result {
	// The compile method should only be called once.
	jutil.CheckState(c.jsRoot == nil)

	func() {
		defer c.GenerateReport()
		c.InitChunks(externs, chunks, options)
		if !c.HasErrors() {
			c.ParseForCompilation()
		}
		if !c.HasErrors() {
			// TODO(bradfordcsmith): The option to instrument for coverage only should belong to the
			//     runner, not the compiler.
			if options.GetInstrumentForCoverageOnly() {
				c.InstrumentForCoverage()
			} else {
				c.Stage1Passes()
				if !c.HasErrors() {
					c.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
					if !c.HasErrors() {
						c.Stage3Passes()
					}
				}
			}
			c.PerformPostCompilationTasks()
		}
	}()
	return c.GetResult()
}

// Stage1Passes performs compiler passes for stage 1 of compilation.
//
// Stage 1 consists primarily of error and type checking passes.
//
// parseForCompilation() must be called before this method is called.
//
// The caller is responsible for also calling generateReport() to generate
// a report of warnings and errors to stderr. See the invocation in compile
// for a good example.
//
// Java: Compiler.stage1Passes()
func (c *Compiler) Stage1Passes() {
	jutil.CheckStateMsg(c.chunkGraph != nil, "No inputs. Did you call init() or initChunks()?")
	jutil.CheckState(!c.HasErrors())
	jutil.CheckState(!c.options.GetInstrumentForCoverageOnly())
	Compiler_RunInCompilerThread(c,
		func() any {
			c.performChecks()
			c.evictPassArenaCaches()
			return nil
		})
}

// Stage2Passes performs compiler passes for stage 2 of compilation.
//
// Stage 2 consists primarily of transpilation and optimization passes.
//
// stage1Passes() must be called before this method is called.
//
// The caller is responsible for also calling generateReport() to generate
// a report of warnings and errors to stderr. See the invocation in compile
// for a good example.
//
// Java: Compiler.stage2Passes(SegmentOfCompilationToRun)
func (c *Compiler) Stage2Passes(segmentOfCompilationToRun CompilerOptions_SegmentOfCompilationToRun) {
	jutil.CheckStateMsg1(
		segmentOfCompilationToRun == CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS ||
			segmentOfCompilationToRun == CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF ||
			segmentOfCompilationToRun == CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF,
		"Unsupported segment of compilation to run in stage 2: %s",
		segmentOfCompilationToRun)
	jutil.CheckStateMsg(c.chunkGraph != nil, "No inputs. Did you call init() or initChunks()?")
	jutil.CheckState(!c.HasErrors())
	jutil.CheckState(!c.options.GetInstrumentForCoverageOnly())
	weakChunk := c.chunkGraph.GetChunkByName(JSChunk_WEAK_CHUNK_NAME)
	if weakChunk != nil {
		for _, i := range c.chunkGraph.GetAllInputs() {
			if i.GetSourceFile().IsWeak() {
				jutil.CheckStateMsg(i.GetChunk() == weakChunk, "Expected all weak files to be in the weak chunk.")
			}
		}
	}
	Compiler_RunInCompilerThread(c,
		func() any {
			if c.options.ShouldOptimize() {
				c.PerformTranspilationAndOptimizations(segmentOfCompilationToRun)
			}
			c.evictPassArenaCaches()
			return nil
		})
}

// Stage3Passes performs compiler passes for stage 3 of compilation.
//
// Stage 3 consists primarily of localization passes.
//
// stage2Passes() must be called before this method is called.
//
// The caller is responsible for also calling generateReport() to generate
// a report of warnings and errors to stderr. See the invocation in compile
// for a good example.
//
// Java: Compiler.stage3Passes()
func (c *Compiler) Stage3Passes() {
	jutil.CheckStateMsg(c.chunkGraph != nil, "No inputs. Did you call init() or initChunks()?")
	jutil.CheckStateMsg1(!c.HasErrors(), "Cannot run stage3Passes with errors present: %s", c.GetErrors())
	jutil.CheckState(!c.options.GetInstrumentForCoverageOnly())
	Compiler_RunInCompilerThread(c,
		func() any {
			if c.options.ShouldOptimize() {
				c.PerformFinalizations()
			}
			c.evictPassArenaCaches()
			return nil
		})
}

// DisableThreads disables threads. This is for clients that run on
// AppEngine and don't have threads.
//
// Java: Compiler.disableThreads()
func (c *Compiler) DisableThreads() {
	c.compilerExecutor.DisableThreads()
}

// SetTimeout sets the timeout when Compiler is run in a thread
//
// timeout: seconds to wait before timeout
//
// Java: Compiler.setTimeout(int)
func (c *Compiler) SetTimeout(timeout int) {
	c.compilerExecutor.SetTimeout(timeout)
}

// Compiler_RunInCompilerThread: the primary purpose of this method is to
// run the provided code with a larger than standard stack.
//
// Go-only: language-adapter — Java's generic instance method
// <T> T runInCompilerThread(Callable<T>); Go methods cannot have type
// parameters.
//
// Java: Compiler.runInCompilerThread(Callable<T>)
func Compiler_RunInCompilerThread[T any](c *Compiler, callable func() T) T {
	return CompilerExecutor_RunInCompilerThread(
		c.compilerExecutor, callable, c.options != nil && c.options.GetTracerMode().IsOn())
}

// Java: Compiler.performChecks()
func (c *Compiler) performChecks() {
	if c.options.GetSkipNonTranspilationPasses() {
		// i.e. whitespace-only mode, which will not work with goog.module without:
		c.WhitespaceOnlyPasses()
		if c.options.NeedsTranspilationFrom(c.options.GetLanguageIn().ToFeatureSet()) {
			c.TranspileAndDontCheck()
		}
	} else {
		c.Check()
	}
}

// PerformPostCompilationTasks performs all the bookkeeping required at the
// end of a compilation.
//
// This method must be called if the compilation makes it as far as doing
// checks.
//
// DON'T call it if the compiler threw an exception.
//
// DO call it even when hasErrors() returns true.
//
// Java: Compiler.performPostCompilationTasks()
func (c *Compiler) PerformPostCompilationTasks() {
	Compiler_RunInCompilerThread(c,
		func() any {
			c.performPostCompilationTasksInternal()
			return nil
		})
}

// performPostCompilationTasksInternal performs all the bookkeeping
// required at the end of a compilation.
//
// Java: Compiler.performPostCompilationTasksInternal()
func (c *Compiler) performPostCompilationTasksInternal() {
	if c.options.GetDevMode() == CompilerOptions_DevMode_START_AND_END {
		c.runValidityCheck()
	}

	// TODO: b/501060517 - consider deleting this. The code printer will already print out
	// @closureUnaware contents as if they were unwrapped, but the explicit wrapping might be
	// more maintainable.
	ManageClosureUnawareCode_Unwrap(c).Process(c.externsRoot, c.jsRoot)

	if c.tracker != nil {
		if c.options.GetTracerOutput() == nil {
			c.tracker.OutputTracerReport(c.outStream)
		} else {
			jutil.Try(func() {
				stream := jutil.Files_NewOutputStream(c.options.GetTracerOutput())
				out := jutil.NewPrintStream(stream)
				defer func() {
					// PrintStream.close: flush, then close the stream.
					out.Flush()
					stream.Close()
				}()
				c.tracker.OutputTracerReport(out)
			}, func(e jutil.Throwable) {
				// Java: catch (Exception e), which does not see java.lang.Errors.
				if _, isError := e.(jutil.IError); isError {
					panic(e)
				}
				panic(jutil.NewRuntimeExceptionWithCause(e))
			})
		}
	}
}

// InstrumentForCoverage instruments code for coverage.
//
// parseForCompilation() must be called before this method is called.
//
// The caller is responsible for also calling generateReport() to generate
// a report of warnings and errors to stderr. See the invocation in compile
// for a good example.
//
// This method is mutually exclusive with stage1Passes() and
// stage2Passes(). Either call those two methods or this one, but not both.
//
// Java: Compiler.instrumentForCoverage()
func (c *Compiler) InstrumentForCoverage() {
	jutil.CheckStateMsg(c.chunkGraph != nil, "No inputs. Did you call init() or initChunks()?")
	jutil.CheckState(!c.HasErrors())
	Compiler_RunInCompilerThread(c,
		func() any {
			jutil.CheckState(c.options.GetInstrumentForCoverageOnly())
			jutil.CheckState(!c.HasErrors())
			if c.options.GetInstrumentForCoverageOption() != CompilerOptions_InstrumentOption_NONE {
				c.instrumentForCoverageInternal(c.options.GetInstrumentForCoverageOption())
			}
			return nil
		})
}

// Java: Compiler.instrumentForCoverageInternal(InstrumentOption)
func (c *Compiler) instrumentForCoverageInternal(instrumentForCoverageOption CompilerOptions_InstrumentOption) {
	tracer := c.NewTracer("instrumentationPass")
	c.Process(
		NewCoverageInstrumentationPassWithInstrumentOption(
			c,
			CoverageInstrumentationPass_CoverageReach_ALL,
			instrumentForCoverageOption,
			c.options.GetProductionInstrumentationArrayName()))
	c.StopTracer(tracer, "instrumentationPass")
}

// ParseForCompilation parses input files in preparation for compilation.
//
// Either init() or initChunks() must be called first to set up the input
// files to be read.
//
// TODO(bradfordcsmith): Rename this to parse()
//
// Java: Compiler.parseForCompilation()
func (c *Compiler) ParseForCompilation() {
	jutil.CheckStateMsg(
		c.typedAstFilesystem == nil, "Unnecessary if initWithTypedAstFilesystem was called")
	Compiler_RunInCompilerThread(c,
		func() any {
			c.parseForCompilationInternal()
			return nil
		})
}

// parseForCompilationInternal parses input files in preparation for
// compilation.
//
// Either init() or initChunks() must be called first to set up the input
// files to be read.
//
// TODO(bradfordcsmith): Rename this to parse()
//
// Java: Compiler.parseForCompilationInternal()
func (c *Compiler) parseForCompilationInternal() {
	CompilerOptionsPreprocessor_Preprocess(c.options)
	c.MaybeSetTracker()
	c.ParseInputs()
}

// Parse parses input files without doing progress tracking that is part of
// a full compile.
//
// Either init() or initChunks() must be called first to set up the input
// files to be read.
//
// TODO(bradfordcsmith): Rename this to parseIndependentOfCompilation() or
// similar.
//
// Java: Compiler.parse()
func (c *Compiler) Parse() {
	c.ParseInputs()
}

// Java: Compiler.parse(SourceFile)
func (c *Compiler) ParseOfSourceFile(file *SourceFile) *rhino.Node {
	c.InitCompilerOptionsIfTesting()
	Compiler_logger.Finest("Parsing: " + file.GetName())
	return NewCompilerInput(file).GetAstRoot(c)
}

// Java: Compiler.getPassConfig()
func (c *Compiler) GetPassConfig() PassConfig {
	if c.passes == nil {
		c.passes = c.CreatePassConfigInternal()
	}
	return c.passes
}

// CreatePassConfigInternal creates the passes object. Clients should use
// setPassConfig instead of overriding this.
//
// Java: Compiler.createPassConfigInternal()
func (c *Compiler) CreatePassConfigInternal() PassConfig {
	return NewDefaultPassConfig(c.options)
}

// SetPassConfig sets the PassConfig to use with this Compiler.
//
// Panics with NullPointerException if passes is null, and with
// IllegalStateException if this.passes has already been assigned.
//
// Java: Compiler.setPassConfig(PassConfig)
func (c *Compiler) SetPassConfig(passes PassConfig) {
	// Important to check for null because if setPassConfig(null) is
	// called before this.passes is set, getPassConfig() will create a
	// new PassConfig object and use that, which is probably not what
	// the client wanted since they probably meant to use their
	// own PassConfig object.
	jutil.CheckNotNull(passes)
	jutil.CheckStateMsg(c.passes == nil, "setPassConfig was already called")
	c.passes = passes
}

// Java: Compiler.whitespaceOnlyPasses()
func (c *Compiler) WhitespaceOnlyPasses() {
	c.runCustomPasses(CustomPassExecutionTime_BEFORE_CHECKS)

	t := c.NewTracer("runWhitespaceOnlyPasses")
	defer c.StopTracer(t, "runWhitespaceOnlyPasses")
	for _, pf := range c.GetPassConfig().GetWhitespaceOnlyPasses().Build() {
		pf.Create(c).Process(c.externsRoot, c.jsRoot)
	}
}

// Java: Compiler.markTranspiledFiles()
func (c *Compiler) markTranspiledFiles() {
	// note whether any script will get transpiled
	for child := c.jsRoot.GetFirstChild(); child != nil; child = child.GetNext() {
		jutil.CheckState(child.IsScript())
		featureSet := NodeUtil_GetFeatureSetOfScript(child)
		if !c.options.GetOutputFeatureSet().Contains(featureSet) {
			c.SetTranspiledFiles(true)
			return // marked atleast 1 script that needs tranpile
		}
	}
}

// Java: Compiler.transpileAndDontCheck()
func (c *Compiler) TranspileAndDontCheck() {
	t := c.NewTracer("runTranspileOnlyPasses")
	c.markTranspiledFiles()
	defer c.StopTracer(t, "runTranspileOnlyPasses")
	for _, pf := range c.GetPassConfig().GetTranspileOnlyPasses().Build() {
		if c.HasErrors() {
			return
		}
		c.BeforePass(pf.GetName())
		pf.Create(c).Process(c.externsRoot, c.jsRoot)
		c.AfterPass(pf.GetName())
	}
}

// Java: Compiler.createPhaseOptimizer()
func (c *Compiler) createPhaseOptimizer() *PhaseOptimizer {
	phaseOptimizer := NewPhaseOptimizer(c, c.tracker)
	if c.options.GetDevMode() == CompilerOptions_DevMode_EVERY_PASS {
		phaseOptimizer.SetValidityCheck(c.validityCheck)
	}
	if c.options.GetCheckDeterminism() {
		phaseOptimizer.SetPrintAstHashcodes(true)
	}
	return phaseOptimizer
}

// Java: Compiler.check()
func (c *Compiler) Check() {
	c.runCustomPasses(CustomPassExecutionTime_BEFORE_CHECKS)

	c.phaseOptimizer = c.createPhaseOptimizer()
	c.phaseOptimizer.Consume(c.GetPassConfig().GetChecks().Build())
	c.phaseOptimizer.Process(c.externsRoot, c.jsRoot)
	if c.HasErrors() {
		return
	}

	c.runCustomPasses(CustomPassExecutionTime_BEFORE_OPTIMIZATIONS)
	c.phaseOptimizer = nil
}

// Java: Compiler.setExternExports(String)
func (c *Compiler) SetExternExports(externExports string) {
	c.externExports = &externExports
}

// Java: Compiler.process(CompilerPass)
func (c *Compiler) Process(p CompilerPass) {
	p.Process(c.externsRoot, c.jsRoot)
}

// Java: Compiler.runValidityCheck()
func (c *Compiler) runValidityCheck() {
	c.validityCheck.Create(c).Process(c.externsRoot, c.jsRoot)
}

// runCustomPasses runs custom passes that are designated to run at a
// particular time.
//
// Java: Compiler.runCustomPasses(CustomPassExecutionTime)
func (c *Compiler) runCustomPasses(executionTime CustomPassExecutionTime) {
	customPasses := c.options.GetCustomPassesAt(executionTime)
	t := c.NewTracer("runCustomPasses")
	defer c.StopTracer(t, "runCustomPasses")
	for _, p := range customPasses {
		c.Process(p)
	}
}

// Java: Compiler.beforePass(String)
func (c *Compiler) BeforePass(passName string) {
	c.currentPassIndex++
}

// Java: Compiler.afterPass(String)
func (c *Compiler) AfterPass(passName string) {
	c.maybePrintSourceAfterEachPass(passName)
}

// GetLifeCycleStage returns the current life-cycle stage of the AST we're
// working on.
//
// Java: Compiler.getLifeCycleStage()
func (c *Compiler) GetLifeCycleStage() AbstractCompiler_LifeCycleStage {
	return c.stage
}

// SetLifeCycleStage sets the current life-cycle state.
//
// Java: Compiler.setLifeCycleStage(LifeCycleStage)
func (c *Compiler) SetLifeCycleStage(stage AbstractCompiler_LifeCycleStage) {
	c.stage = stage
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
		jutil.NewLinkedHashSet[string](),
		Scope_CreateGlobalScope(rhino.NewNode(rhino.Token_SCRIPT)))
}

// Java: Compiler.createExpressionDecomposer(Supplier<String>, ImmutableSet<String>, Scope)
func (c *Compiler) CreateExpressionDecomposer(
	uniqueNameIdSupplier func() string,
	knownConstantFunctions *jutil.LinkedHashSet[string],
	scope *Scope) *ExpressionDecomposer {
	// If the output is ES5, then it may end up running on IE11, so enable a workaround
	// for one of its bugs.
	var enabledWorkarounds *jutil.EnumSet[ExpressionDecomposer_Workaround]
	if parser.FeatureSet_ES5.Contains(c.GetOptions().GetOutputFeatureSet()) {
		enabledWorkarounds = jutil.NewEnumSet(ExpressionDecomposer_Workaround_BROKEN_IE11_LOCATION_ASSIGN)
	} else {
		enabledWorkarounds = jutil.NewEnumSet[ExpressionDecomposer_Workaround]()
	}
	return NewExpressionDecomposer(
		c, uniqueNameIdSupplier, knownConstantFunctions, scope, enabledWorkarounds)
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
func (c *Compiler) CreateOrReopenLog(owner reflect.Type, firstNamePart string, restNameParts ...string) diagnostic.LogFile {
	if !c.IsDebugLoggingEnabled() {
		return diagnostic.LogFile_CreateNoOp()
	}

	dir := c.GetOptions().GetDebugLogDirectory()
	relativeParts := jutil.Path_Of(firstNamePart, restNameParts...)
	file := dir.ResolveString(owner.Name()).Resolve(relativeParts)

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
func (c *Compiler) CreateOrReopenIndexedLog(owner reflect.Type, firstNamePart string, restNameParts ...string) diagnostic.LogFile {
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
			log := c.CreateOrReopenIndexedLog(reflect.TypeOf(c).Elem(), "source_after_pass", passName)
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
		jutil.CheckNotNull(c.externsRoot)
		jutil.CheckNotNull(c.jsRoot)
		for _, r := range []*rhino.Node{c.externsRoot, c.jsRoot} {
			for fileNode := r.GetFirstChild(); fileNode != nil; fileNode = fileNode.GetNext() {
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
		statementFilter := func(statement *rhino.Node) bool {
			// filter to just those statements containing interesting qnames
			return NodeUtil_Has(
				statement,
				func(n *rhino.Node) bool {
					jutil.CheckNotNull(n)
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
				func(node *rhino.Node) bool { return true })
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
func (c *Compiler) getTopLevelStatements(root *rhino.Node) []*rhino.Node {
	cb := &compiler_topLevelStatementsCallback{}
	NodeTraversal_Traverse(c, root, cb)
	return cb.builder
}

// compiler_topLevelStatementsCallback is the anonymous
// AbstractPreOrderCallback of getTopLevelStatements: it collects the
// outermost statements.
type compiler_topLevelStatementsCallback struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	builder []*rhino.Node
}

// Java: Compiler.getTopLevelStatements(...).new AbstractPreOrderCallback() {...}.shouldTraverse(NodeTraversal, Node, Node)
func (cb *compiler_topLevelStatementsCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if NodeUtil_IsStatement(n) {
		cb.builder = append(cb.builder, n)
		return false
	} else {
		return true
	}
}

// Java: Compiler.getScriptNode(String)
func (c *Compiler) GetScriptNode(filename string) *rhino.Node {
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
	typ := "Compiler"
	return NewTracerWithType(&typ, &comment)
}

// Java: Compiler.stopTracer(Tracer, String)
func (c *Compiler) StopTracer(t *Tracer, passName string) {
	result := t.Stop()
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
func (c *Compiler) GetRoot() *rhino.Node {
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
		c.uniqueNameIdSupplier = func() string { return strconv.Itoa(c.nextUniqueNameId()) }
	}
	return c.uniqueNameIdSupplier
}

// Java: Compiler.areNodesEqualForInlining(Node, Node)
func (c *Compiler) AreNodesEqualForInlining(n1 *rhino.Node, n2 *rhino.Node) bool {
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
	return c.inputsById.Get(id)
}

// Java: Compiler.putCompilerInput(CompilerInput)
func (c *Compiler) putCompilerInput(input *CompilerInput) *CompilerInput {
	input.SetCompiler(c)
	return c.inputsById.Put(input.GetInputId(), input)
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
	c.typeRegistry = nil
	c.typeValidator = nil
	c.abstractInterpreter = nil
}

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
			func(script *rhino.Node) bool { return script.IsFromExterns() || script.IsFirstChildOf(c.GetJsRoot()) })
		gn.GetNameForest() // ensure namespace is generated
		gn.ReleaseScopeArena()
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

// ------------------------------------------------------------------------
// Parsing
// ------------------------------------------------------------------------

// ParseInputs parses the externs and main inputs.
//
// Returns a synthetic root node whose two children are the externs root
// and the main root
//
// Java: Compiler.parseInputs()
func (c *Compiler) ParseInputs() (result *rhino.Node) {
	defer parser.DrainParseArenaPool()
	devMode := c.options.GetDevMode() != CompilerOptions_DevMode_OFF

	// If old roots exist (we are parsing a second time), detach each of the
	// individual file parse trees.
	c.externsRoot.DetachChildren()
	c.jsRoot.DetachChildren()
	c.scriptNodeByFilename.Clear()

	tracer := c.NewTracer(PassNames_PARSE_INPUTS)
	c.BeforePass(PassNames_PARSE_INPUTS)

	defer func() {
		c.AfterPass(PassNames_PARSE_INPUTS)
		c.StopTracer(tracer, PassNames_PARSE_INPUTS)
	}()
	// Parse externs sources.
	if c.options.GetNumParallelThreads() > 1 {
		NewPrebuildAst(c, c.options.GetNumParallelThreads()).Prebuild(c.externs.All())
	}
	for input := range c.externs.All() {
		n := jutil.CheckNotNull(input.GetAstRoot(c))
		if c.HasErrors() {
			return nil
		}
		c.externsRoot.AddChildToBack(n)
		c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), n)
	}

	if c.options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) ||
		c.options.GetProcessCommonJSModules() {
		c.InitializeModuleLoader()
	} else {
		// Use an empty module loader if we're not actually dealing with modules.
		c.moduleLoader = ModuleLoader_EMPTY
	}

	if c.options.GetDependencyOptions().NeedsManagement() {
		c.findModulesFromEntryPoints(
			c.options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES),
			c.options.GetProcessCommonJSModules())
	} else if c.options.NeedsTranspilationFrom(parser.FeatureSet_ES2015_MODULES) ||
		c.options.GetProcessCommonJSModules() {
		if c.options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
			c.parsePotentialModules(c.chunkGraph.GetAllInputs())
		}

		// Build a map of module identifiers for any input which provides no namespace.
		// These files could be imported modules which have no exports, but do have side effects.
		inputModuleIdentifiers := jutil.NewLinkedHashMap[string, *CompilerInput]()
		for _, input := range c.chunkGraph.GetAllInputs() {
			if len(input.GetKnownProvides()) == 0 {
				modPath := c.moduleLoader.Resolve(input.GetSourceFile().GetName())
				inputModuleIdentifiers.Put(modPath.ToModuleName(), input)
			}
		}

		// Find out if any input attempted to import a module that had no exports.
		// In this case we must force module rewriting to occur on the imported file
		inputsToRewrite := jutil.NewLinkedHashMap[string, *CompilerInput]()
		for _, input := range c.chunkGraph.GetAllInputs() {
			for _, require := range input.GetKnownRequiredSymbols() {
				if inputModuleIdentifiers.ContainsKey(require) &&
					!inputsToRewrite.ContainsKey(require) {
					inputsToRewrite.Put(require, inputModuleIdentifiers.Get(require))
				}
			}
		}

		for input := range inputsToRewrite.Values() {
			input.SetJsModuleType(CompilerInput_ModuleType_IMPORTED_SCRIPT)
		}
	}

	if c.moduleLoader != nil {
		c.moduleLoader.SetErrorHandler(c.AsErrorHandler())
	}

	c.OrderInputs()
	c.markClosureUnawareCode()

	// If in IDE mode, we ignore the error and keep going.
	if c.HasErrors() {
		return nil
	}

	// Build the AST.
	allInputs := c.chunkGraph.GetAllInputs()
	if c.options.GetNumParallelThreads() > 1 {
		NewPrebuildAst(c, c.options.GetNumParallelThreads()).Prebuild(slices.Values(allInputs))
	}

	needSourceInfo := c.options.ShouldGatherSourceMapInfo() ||
		c.options.GetExternExportsPath() != nil ||
		len(c.options.GetReplaceStringsFunctionDescriptions()) != 0
	annotatedInParallel := false
	if needSourceInfo && !devMode && c.options.GetDevMode() == CompilerOptions_DevMode_OFF &&
		c.options.GetNumParallelThreads() > 1 && len(allInputs) > 1 && !c.HasErrors() {
		sia := SourceInformationAnnotator_Create()
		failures := parallelPrebuild(len(allInputs), c.options.GetNumParallelThreads(), func(i int) {
			if root := allInputs[i].GetAstRoot(c); root != nil {
				sia.AnnotateTree(root)
			}
		})
		for _, f := range failures {
			if f != nil {
				panic(f)
			}
		}
		annotatedInParallel = true
	}

	for _, input := range allInputs {
		n := jutil.CheckNotNull(input.GetAstRoot(c))
		if devMode {
			c.runValidityCheck()
			if c.HasErrors() {
				return nil
			}
		}

		if needSourceInfo && !annotatedInParallel {
			// Annotate the nodes in the tree with information from the
			// input file. This information is used to construct the SourceMap.
			if c.options.GetDevMode() == CompilerOptions_DevMode_OFF {
				SourceInformationAnnotator_Create().AnnotateTree(n)
			} else {
				sia := SourceInformationAnnotator_CreateWithAnnotationChecks(input.GetName())
				NodeTraversal_Traverse(c, n, sia)
			}
		}

		if NodeUtil_IsFromTypeSummary(n) {
			input.SetIsExtern()
			c.externsRoot.AddChildToBack(n)
		} else {
			c.jsRoot.AddChildToBack(n)
		}
		c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), n)
	}

	if c.HasErrors() {
		return nil
	}
	// Save on memory. Any future calls to "createInputConsideringTypedAstFilesystem" will throw.
	// TODO(lharker): do we actually need the synthetic externs file in stage 2?
	if c.typedAstFilesystem != nil {
		// If anything (besides the synthetic externs) is in the typedAstFilesystem, that indicates
		// it never should have been added in the first place: every CompilerInput should now be
		// initialized, a process which removes it from the typedAstFilesystem
		jutil.CheckStateMsg1(
			c.typedAstFilesystem.Size() == 1 || c.typedAstFilesystem.IsEmpty(),
			"%s",
			c.typedAstFilesystem.Size())
	}
	return c.externAndJsRoot
}

// PreferRegexParser hints to other parts of the compiler as to whether
// every input will eventually require a full AST parse.
//
// NOTE: this is just intended as a hint towards improving performance.
// Returning true or false shouldn't affect any non-performance behavior.
//
// Java: Compiler.preferRegexParser()
func (c *Compiler) PreferRegexParser() bool {
	return c.preferRegexParser
}

// SetPreferRegexParser artificially sets the return value of
// `preferRegexParser`
//
// Non-public and not exposed on AbstractCompiler, as for now we want to
// maintain control of when the regex vs. AST parser is used.
//
// Java: Compiler.setPreferRegexParser(boolean)
func (c *Compiler) SetPreferRegexParser(preferRegexParser bool) {
	c.preferRegexParser = preferRegexParser
}

// Java: Compiler.orderInputsWithLargeStack()
func (c *Compiler) OrderInputsWithLargeStack() {
	Compiler_RunInCompilerThread(c,
		func() any {
			tracer := c.NewTracer("orderInputsWithLargeStack")
			defer c.StopTracer(tracer, "orderInputsWithLargeStack")
			c.OrderInputs()
			return nil
		})
}

// Java: Compiler.orderInputs()
func (c *Compiler) OrderInputs() {
	c.maybeDoThreadedParsing()

	// Before dependency pruning, save a copy of the original inputs to use for externs hoisting.
	originalInputs := append([]*CompilerInput(nil), c.chunkGraph.GetAllInputs()...)

	if c.tracker != nil {
		c.tracker.RecordPrePruningInputCount(slices.Values(originalInputs), c.externs.All())
	}

	// Externs must be marked before dependency management since it needs to know what is an extern.
	c.markExterns(originalInputs)

	// Check if the sources need to be re-ordered.
	staleInputs := false
	if c.options.GetDependencyOptions().NeedsManagement() {
		jutil.Try(func() {
			jutil.Try(func() {
				result :=
					c.chunkGraph.ManageDependencies(c, c.options.GetDependencyOptions())
				staleInputs = true
				c.maybeLogPruningAnalysis(result)
			}, func(e *SortedDependencies_MissingProvideException) {
				c.Report(JSError_Make(Compiler_MISSING_ENTRY_ERROR, e.GetMessage()))
			})
		}, func(e *JSChunkGraph_MissingChunkException) {
			c.Report(JSError_Make(Compiler_MISSING_CHUNK_ERROR, e.GetMessage()))
		})
	}
	c.HoistExterns(originalInputs)

	// Manage dependencies may move weak sources around, and end up with empty chunks.
	c.fillEmptyChunks(c.GetChunks())
	c.HoistNoCompileFiles()

	if staleInputs {
		c.repartitionInputs()
	}
	c.maybeLogPrunedInputs()
}

// findModulesFromEntryPoints finds modules by recursively traversing
// dependencies starting with the entry points.
//
// Triggers either a regex parse or full AST-based parse of each file in
// order to gather dependency information.
//
// If the dependency mode is set to PRUNE_LEGACY, inputs which the regex
// parse does not identify as ES modules and which do not contain any
// provide statements are considered to be additional entry points.
//
// Java: Compiler.findModulesFromEntryPoints(boolean, boolean)
func (c *Compiler) findModulesFromEntryPoints(
	supportEs6Modules bool, supportCommonJSModules bool) {
	c.maybeDoThreadedParsing()
	entryPoints := []*CompilerInput{}
	inputsByProvide := jutil.NewLinkedHashMap[string, *CompilerInput]()
	inputsByIdentifier := jutil.NewLinkedHashMap[string, *CompilerInput]()
	for _, input := range c.chunkGraph.GetAllInputs() {
		var provides []string
		for _, p := range input.GetProvides() {
			if !strings.HasPrefix(p, "module$") {
				provides = append(provides, p)
			}
		}
		if !c.options.GetDependencyOptions().ShouldDropMoochers() && len(provides) == 0 {
			entryPoints = append(entryPoints, input)
		}
		inputsByIdentifier.Put(
			ModuleIdentifier_ForFile(input.GetPath().String()).String(), input)
		for _, provide := range provides {
			inputsByProvide.Put(provide, input)
		}
	}
	for _, moduleIdentifier := range c.options.GetDependencyOptions().EntryPoints() {
		input := inputsByProvide.Get(moduleIdentifier.String())
		if input == nil {
			input = inputsByIdentifier.Get(moduleIdentifier.String())
		}
		if input != nil {
			entryPoints = append(entryPoints, input)
		}
	}

	workingInputSet := jutil.NewLinkedHashSet[*CompilerInput]()
	for _, input := range c.chunkGraph.GetAllInputs() {
		workingInputSet.Add(input)
	}
	for _, entryPoint := range entryPoints {
		c.findModulesFromInput(
			entryPoint,
			/* wasImportedByModule= */ false,
			workingInputSet,
			inputsByIdentifier,
			inputsByProvide,
			supportEs6Modules,
			supportCommonJSModules)
	}
}

// findModulesFromInput traverses an input's dependencies to find
// additional modules.
//
// Java: Compiler.findModulesFromInput(CompilerInput, boolean, Set<CompilerInput>, Map<String, CompilerInput>, Map<String, CompilerInput>, boolean, boolean)
func (c *Compiler) findModulesFromInput(
	input *CompilerInput,
	wasImportedByModule bool,
	inputs *jutil.LinkedHashSet[*CompilerInput],
	inputsByIdentifier *jutil.LinkedHashMap[string, *CompilerInput],
	inputsByProvide *jutil.LinkedHashMap[string, *CompilerInput],
	supportEs6Modules bool,
	supportCommonJSModules bool) {
	if !inputs.Remove(input) {
		// It's possible for a module to be included as both a script
		// and a module in the same compilation. In these cases, it should
		// be forced to be a module.
		if wasImportedByModule && input.GetJsModuleType() == CompilerInput_ModuleType_NONE {
			input.SetJsModuleType(CompilerInput_ModuleType_IMPORTED_SCRIPT)
		}
		return
	}

	findDeps := NewFindModuleDependencies(
		c, supportEs6Modules, supportCommonJSModules, c.inputPathByWebpackId)
	findDeps.Process(jutil.CheckNotNull(input.GetAstRoot(c)))

	// If this input was imported by another module, it is itself a module
	// so we force it to be detected as such.
	if wasImportedByModule && input.GetJsModuleType() == CompilerInput_ModuleType_NONE {
		input.SetJsModuleType(CompilerInput_ModuleType_IMPORTED_SCRIPT)
	}
	c.moduleTypesByName.Put(input.GetPath().ToModuleName(), input.GetJsModuleType())

	var allDeps []string
	allDeps = append(allDeps, input.GetRequiredSymbols()...)
	allDeps = append(allDeps, input.GetDynamicRequires()...)
	allDeps = append(allDeps, input.GetTypeRequires()...)
	for _, requiredNamespace := range allDeps {
		var requiredInput *CompilerInput
		requiredByModuleImport := false
		if inputsByProvide.ContainsKey(requiredNamespace) {
			requiredInput = inputsByProvide.Get(requiredNamespace)
		} else if inputsByIdentifier.ContainsKey(requiredNamespace) {
			requiredByModuleImport = true
			requiredInput = inputsByIdentifier.Get(requiredNamespace)
		}

		if requiredInput != nil {
			c.findModulesFromInput(
				requiredInput,
				requiredByModuleImport,
				inputs,
				inputsByIdentifier,
				inputsByProvide,
				supportEs6Modules,
				supportCommonJSModules)
		}
	}
}

// HoistExterns hoists inputs with the @externs annotation into the externs
// list.
//
// Java: Compiler.hoistExterns(ImmutableList<CompilerInput>)
func (c *Compiler) HoistExterns(originalInputs []*CompilerInput) {
	staleInputs := false

	for _, input := range originalInputs {
		if c.hoistIfExtern(input) {
			staleInputs = true
		}
	}
	if staleInputs {
		c.repartitionInputs()
	}
}

// hoistIfExtern hoists a compiler input to externs if it contains the
// @externs annotation. Return whether or not the given input was hoisted.
//
// Java: Compiler.hoistIfExtern(CompilerInput)
func (c *Compiler) hoistIfExtern(input *CompilerInput) bool {
	if input.GetHasExternsAnnotation() {
		// If the input file is explicitly marked as an externs file, then move it out of the main
		// JS root and put it with the other externs.
		root := input.GetAstRoot(c)
		c.externsRoot.AddChildToBack(root)
		c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)

		chunk := input.GetChunk()
		if chunk != nil {
			chunk.Remove(input)
		}

		c.externs.Add(input)
		return true
	}
	return false
}

// markClosureUnawareCode marks source files annotated with
// `@closureUnaware` with the closureUnawareCode bit on SourceFile.
//
// Java: Compiler.markClosureUnawareCode()
func (c *Compiler) markClosureUnawareCode() {
	for _, input := range c.chunkGraph.GetAllInputs() {
		root := input.GetAstRoot(c)
		info := root.GetJSDocInfo()
		if info != nil && info.IsClosureUnawareCode() {
			input.GetSourceFile().MarkAsClosureUnawareCode()
		}
	}
}

// markExterns marks inputs with the @externs annotation as an Extern
// source file type. This is so that externs marking can be done before
// dependency management, and externs hoisting done after dependency
// management.
//
// Java: Compiler.markExterns(ImmutableList<CompilerInput>)
func (c *Compiler) markExterns(originalInputs []*CompilerInput) {
	for _, input := range originalInputs {
		if input.GetHasExternsAnnotation() {
			input.SetIsExtern()
		}
	}
}

// HoistNoCompileFiles hoists inputs with the @nocompile annotation out of
// the inputs.
//
// Java: Compiler.hoistNoCompileFiles()
func (c *Compiler) HoistNoCompileFiles() {
	staleInputs := false
	c.maybeDoThreadedParsing()
	// Iterate a copy because hoisting modifies what we're iterating over.
	for _, input := range append([]*CompilerInput(nil), c.chunkGraph.GetAllInputs()...) {
		if input.GetHasNoCompileAnnotation() {
			input.GetChunk().Remove(input)
			staleInputs = true
		}
	}

	if staleInputs {
		c.repartitionInputs()
	}
}

// Java: Compiler.maybeDoThreadedParsing()
func (c *Compiler) maybeDoThreadedParsing() {
	if c.options.GetNumParallelThreads() > 1 {
		NewPrebuildDependencyInfo(c.options.GetNumParallelThreads()).
			Prebuild(slices.Values(c.chunkGraph.GetAllInputs()))
	}
}

// Java: Compiler.repartitionInputs()
func (c *Compiler) repartitionInputs() {
	c.fillEmptyChunks(c.GetChunks())
	c.RebuildInputsFromChunks()
}

// ProcessJsonInputs transforms JSON files to a module export that closure
// compiler can process and keeps track of any "main" entries in
// package.json files.
//
// Java: Compiler.processJsonInputs(Iterable<CompilerInput>)
func (c *Compiler) ProcessJsonInputs(inputsToProcess []*CompilerInput) *jutil.LinkedHashMap[string, string] {
	rewriteJson := NewRewriteJsonToModule(c)
	for _, input := range inputsToProcess {
		if !strings.HasSuffix(input.GetSourceFile().GetName(), ".json") {
			continue
		}

		input.SetCompiler(c)
		readFailed := false
		jutil.Try(func() {
			// JSON objects need wrapped in parens to parse properly
			input.GetSourceFile().SetCodeDeprecated("(" + input.GetSourceFile().GetCode() + ")")
		}, func(e jutil.IIOException) {
			readFailed = true
		})
		if readFailed {
			continue
		}

		root := jutil.CheckNotNull(input.GetAstRoot(c))
		input.SetJsModuleType(CompilerInput_ModuleType_JSON)
		rewriteJson.Process(nil, root)
	}
	return rewriteJson.GetPackageJsonMainEntries()
}

// Java: Compiler.parsePotentialModules(Iterable<CompilerInput>)
func (c *Compiler) parsePotentialModules(inputsToProcess []*CompilerInput) {
	filteredInputs := []*CompilerInput{}
	for _, input := range inputsToProcess {
		// Only process files that are detected as ES6 modules
		if !c.options.GetDependencyOptions().ShouldPrune() ||
			!JsFileRegexParser_IsSupported() ||
			input.IsEs6Module() {
			filteredInputs = append(filteredInputs, input)
		}
	}
	if c.options.GetNumParallelThreads() > 1 {
		NewPrebuildAst(c, c.options.GetNumParallelThreads()).Prebuild(slices.Values(filteredInputs))
	}
	for _, input := range filteredInputs {
		input.SetCompiler(c)
		// Call getRequires to force regex-based dependency parsing to happen.
		input.GetRequires()
		input.SetJsModuleType(CompilerInput_ModuleType_ES6)
	}
}

// NewCompilerOptions allows subclasses to override the default
// CompileOptions object.
//
// Java: Compiler.newCompilerOptions()
func (c *Compiler) NewCompilerOptions() *CompilerOptions {
	return NewCompilerOptions()
}

// Java: Compiler.initCompilerOptionsIfTesting()
func (c *Compiler) InitCompilerOptionsIfTesting() {
	if c.options == nil {
		// initialization for tests that don't initialize the compiler
		// by the normal mechanisms.
		c.InitOptions(c.NewCompilerOptions())
		// allow tests to use unsupported features.
		c.options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	}
}

// Java: Compiler.parseSyntheticCode(String, String)
func (c *Compiler) ParseSyntheticCode(fileName string, js string) *rhino.Node {
	c.InitCompilerOptionsIfTesting()
	source := SourceFile_FromCode(compiler_SYNTHETIC_FILE_NAME_PREFIX+fileName+"] ", js)
	c.addFilesToSourceMap([]*SourceFile{source})
	return c.parseCodeHelper(source)
}

// Java: Compiler.parseTestCode(String)
func (c *Compiler) ParseTestCode(js string) *rhino.Node {
	c.InitCompilerOptionsIfTesting()
	c.InitBasedOnOptions()
	// NOTE: The file name used here is important, because at least Es6RewriteClassTest expects
	// it to match the default name used by the CompilerTestCase class for the source file
	// created when its `test("source string", "expected string")` method is called.
	// Also, some test cases rely on generated variable names that are based partially on the file
	// name or error messages that include the file name.
	return c.parseCodeHelper(SourceFile_FromCode("testcode", js))
}

// Java: Compiler.parseTestCode(ImmutableList<String>)
func (c *Compiler) ParseTestCodeOfList(code []string) *rhino.Node {
	c.InitCompilerOptionsIfTesting()
	c.InitBasedOnOptions()
	srcs := make([]*SourceFile, len(code))
	for index, value := range code {
		srcs[index] = SourceFile_FromCode("testcode"+strconv.Itoa(index), value)
	}
	return c.parseCodeHelperOfList(srcs)
}

// Java: Compiler.parseCodeHelper(SourceFile)
func (c *Compiler) parseCodeHelper(src *SourceFile) *rhino.Node {
	input := NewCompilerInput(src)
	c.putCompilerInput(input)
	root := input.GetAstRoot(c)
	c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)
	return jutil.CheckNotNull(root)
}

// Java: Compiler.parseCodeHelper(List<SourceFile>)
func (c *Compiler) parseCodeHelperOfList(srcs []*SourceFile) *rhino.Node {
	root := rhino.IR_Root()
	for _, src := range srcs {
		root.AddChildToBack(c.parseCodeHelper(src))
	}
	return root
}

// Java: Compiler.getDefaultErrorReporter()
func (c *Compiler) GetDefaultErrorReporter() rhino.ErrorReporter {
	return c.oldErrorReporter
}

// ------------------------------------------------------------------------
// Convert back to source code
// ------------------------------------------------------------------------

func (c *Compiler) evictPassArenaCaches() {
	evictCachedGlobalNamespaceArenas()
	if a := cachedRUCArena.Swap(nil); a != nil {
		removeUnusedCodeArenaPool.Put(a)
	}
	for i := range cachedRCArenas {
		if a := cachedRCArenas[i].Swap(nil); a != nil {
			referenceCollectorArenaPool.Put(a)
		}
	}
	if a := cachedCCAArena.Swap(nil); a != nil {
		crossChunkArenaPool.Put(a)
	}
	for i := range cachedScopeArenas {
		if a := cachedScopeArenas[i].Swap(nil); a != nil {
			a.trimSmall()
			scopeArenaPool.Put(a)
		}
	}
	if ic := cachedInferConstsCollector.Swap(nil); ic != nil {
		inferConstsCollectorPool.Put(ic)
	}
	if p := cachedIfPoolState.Swap(nil); p != nil && len(p.fnStateChunks) <= 64 {
		ifPool.Put(p)
	}
}

// ToSource converts the main parse tree back to JS code.
//
// Java: Compiler.toSource()
func (c *Compiler) ToSource() string {
	return Compiler_RunInCompilerThread(c,
		func() string {
			c.evictPassArenaCaches()
			tracer := c.NewTracer("toSource")
			defer c.StopTracer(tracer, "toSource")
			cb := NewCompiler_CodeBuilder()
			// We are emitting all the sources at once, so use the SingleBinaryLicenseTracker
			// to de-dupe seen licenses across all inputs.
			lt := newCompiler_SingleBinaryLicenseTracker(c)
			if c.jsRoot != nil {
				var scripts []*rhino.Node
				if c.options.ShouldPrintExterns() {
					for scriptNode := c.externsRoot.GetFirstChild(); scriptNode != nil; scriptNode = scriptNode.GetNext() {
						scripts = append(scripts, scriptNode)
					}
				}
				for scriptNode := c.jsRoot.GetFirstChild(); scriptNode != nil; scriptNode = scriptNode.GetNext() {
					scripts = append(scripts, scriptNode)
				}
				if c.options.GetNumParallelThreads() > 1 && len(scripts) > 1 {
					batchSize := 64
					if len(scripts) <= 64 {
						batchSize = max(c.options.GetNumParallelThreads(), 4)
					}
					type scriptPrinted struct {
						sm       CodePrinter_SourceAndMappings
						licenses []string
						lt       *Compiler_SingleBinaryLicenseTracker
					}
					printed := make([]scriptPrinted, min(len(scripts), batchSize))
					numThreads := c.options.GetNumParallelThreads()
					for batchStart := 0; batchStart < len(scripts); batchStart += batchSize {
						batchEnd := min(batchStart+batchSize, len(scripts))
						batchLen := batchEnd - batchStart
						failures := parallelPrebuild(batchLen, numThreads, func(j int) {
							idx := batchStart + j
							localLT := printed[j].lt
							if localLT == nil {
								localLT = newCompiler_SingleBinaryLicenseTracker(c)
								printed[j].lt = localLT
							} else {
								localLT.currentlySeenLicenses.Clear()
								localLT.lastSeenFile = ""
							}
							c.toSourceAndMappingsInto(scripts[idx], idx == 0, localLT, &printed[j].sm)
							if localLT.currentlySeenLicenses.IsEmpty() {
								printed[j].licenses = nil
							} else {
								printed[j].licenses = localLT.currentlySeenLicenses.ToSlice()
							}
						})
						for _, f := range failures {
							if f != nil {
								for j := 0; j < batchLen; j++ {
									printed[j].sm.Release()
								}
								panic(f)
							}
						}
						for j := 0; j < batchLen; j++ {
							idx := batchStart + j
							for _, lic := range printed[j].licenses {
								if !lt.globallyUniqueLicenses.Contains(lic) {
									lt.currentlySeenLicenses.Add(lic)
								}
							}
							printed[j].licenses = nil
							c.appendSourceAndMappingsToCodeBuilder(cb, lt, idx, scripts[idx], &printed[j].sm)
						}
					}
				} else {
					for i, scriptNode := range scripts {
						c.ToSourceOfCodeBuilder(cb, lt, i, scriptNode)
					}
				}
			}
			return cb.String()
		})
}

// ToSourceOfJSChunk converts the parse tree for a chunk back to JS code.
//
// Consider using toSource(JSChunk, LicenseTracker) for better license
// handling. This call will emit all license text attached to all direct
// inputs to the chunk, which can be very inefficient.
//
// Java: Compiler.toSource(JSChunk)
func (c *Compiler) ToSourceOfJSChunk(chunk *JSChunk) string {
	return c.ToSourceOfLicenseTrackerAndJSChunk(
		// Backwards-compatible implementation of license tracking.
		NewCompiler_ScriptNodeLicensesOnlyTracker(c), chunk)
}

// ToSourceOfLicenseTrackerAndJSChunk converts the parse tree for a chunk
// back to JS code, using the given License Tracker to determine which
// licenses should be emitted before the source code.
//
// licenseTracker: The license tracker implementation to use.
// ChunkGraphAwareLicenseTracker is a suitable implementation when this
// method is being called on each chunk in the chunk graph in dependency
// order.
//
// Java: Compiler.toSource(LicenseTracker, JSChunk)
func (c *Compiler) ToSourceOfLicenseTrackerAndJSChunk(licenseTracker CodePrinter_LicenseTracker, chunk *JSChunk) string {
	return Compiler_RunInCompilerThread(c,
		func() string {
			inputs := chunk.GetInputs()
			numInputs := len(inputs)
			if numInputs == 0 {
				return ""
			}
			cb := NewCompiler_CodeBuilder()
			if numInputs > 1 && c.options != nil && c.options.GetNumParallelThreads() > 1 {
				replay, onlyScript, canParallel := compiler_licenseTrackerReplay(licenseTracker)
				if canParallel {
					scripts := make([]*rhino.Node, numInputs)
					for i := 0; i < numInputs; i++ {
						scriptNode := inputs[i].GetAstRoot(c)
						if scriptNode == nil {
							panic(jutil.NewIllegalArgumentExceptionMsg("Bad chunk: " + chunk.GetName()))
						}
						scripts[i] = scriptNode
					}
					batchSize := 64
					if numInputs <= 64 {
						batchSize = max(c.options.GetNumParallelThreads(), 4)
					}
					type scriptPrinted struct {
						sm  CodePrinter_SourceAndMappings
						rec compiler_ParallelScriptLicenseRecorder
					}
					printed := make([]scriptPrinted, min(numInputs, batchSize))
					numThreads := c.options.GetNumParallelThreads()
					for batchStart := 0; batchStart < numInputs; batchStart += batchSize {
						batchEnd := min(batchStart+batchSize, numInputs)
						batchLen := batchEnd - batchStart
						failures := parallelPrebuild(batchLen, numThreads, func(j int) {
							idx := batchStart + j
							printed[j].rec.Reset(onlyScript)
							c.toSourceAndMappingsInto(scripts[idx], idx == 0, &printed[j].rec, &printed[j].sm)
						})
						for _, f := range failures {
							if f != nil {
								for j := 0; j < batchLen; j++ {
									printed[j].sm.Release()
								}
								panic(f)
							}
						}
						for j := 0; j < batchLen; j++ {
							idx := batchStart + j
							for _, file := range printed[j].rec.seenFiles {
								replay(file)
							}
							c.appendSourceAndMappingsToCodeBuilder(cb, licenseTracker, idx, scripts[idx], &printed[j].sm)
						}
					}
					return cb.String()
				}
			}
			for i := 0; i < numInputs; i++ {
				scriptNode := inputs[i].GetAstRoot(c)
				if scriptNode == nil {
					panic(jutil.NewIllegalArgumentExceptionMsg("Bad chunk: " + chunk.GetName()))
				}
				c.ToSourceOfCodeBuilder(cb, licenseTracker, i, scriptNode)
			}
			return cb.String()
		})
}

// ToSourceOfCodeBuilder writes out JS code from a root node. If printing
// input delimiters, this method will attach a comment to the start of the
// text indicating which input the output derived from. If there were any
// preserve annotations within the root's source, they will also be printed
// in a block comment at the beginning of the output.
//
// The LicenseTracker provided determines which licenses attached to the
// nodes visited should be included in the output. When building a single
// JS bundle, consider using the SingleBinaryLicenseTracker implementation.
//
// Java: Compiler.toSource(CodeBuilder, LicenseTracker, int, Node)
func (c *Compiler) ToSourceOfCodeBuilder(
	cb *Compiler_CodeBuilder,
	licenseTracker CodePrinter_LicenseTracker,
	inputSeqNum int,
	root *rhino.Node) {
	Compiler_RunInCompilerThread(c,
		func() any {
			var sourceAndMappings CodePrinter_SourceAndMappings
			c.toSourceAndMappingsInto(root, inputSeqNum == 0, licenseTracker, &sourceAndMappings)
			c.appendSourceAndMappingsToCodeBuilder(cb, licenseTracker, inputSeqNum, root, &sourceAndMappings)
			return nil
		})
}

func (c *Compiler) appendSourceAndMappingsToCodeBuilder(
	cb *Compiler_CodeBuilder,
	licenseTracker CodePrinter_LicenseTracker,
	inputSeqNum int,
	root *rhino.Node,
	sourceAndMappings *CodePrinter_SourceAndMappings) {
	defer sourceAndMappings.Release()
	if c.options.ShouldPrintInputDelimiter() {
		if (cb.GetLength() > 0) && !cb.EndsWith("\n") {
			cb.Append("\n") // Make sure that the label starts on a new line
		}
		jutil.CheckState(root.IsScript())

		delimiter := c.options.GetInputDelimiter()

		inputName := root.GetInputId().GetIdName()
		sourceName := root.GetSourceFileName()
		jutil.CheckState(sourceName != nil)
		jutil.CheckState(*sourceName != "")

		delimiter = strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(delimiter, "%name%", inputName),
				"%num%", strconv.Itoa(inputSeqNum)),
			"%n%", "\n")

		cb.Append(delimiter).Append("\n")
	}

	code := sourceAndMappings.source

	// Check whether there is any license information that should be emitted.
	if emitted := licenseTracker.EmitLicenses(); emitted != nil && !emitted.IsEmpty() {
		for license := range emitted.All() {
			cb.Append("/*\n").Append(license).Append("*/\n")
		}
	}

	// Check whether there's any actual code to emit.
	if code == "" {
		// Nothing to do.
		return
	}

	// If there is a valid source map, then indicate to it that the current
	// root node's mappings are offset by the given string builder buffer.
	// This offset is a result of licenses being added to the output buffer.
	if c.options.ShouldGatherSourceMapInfo() {
		c.sourceMap.SetStartingPosition(cb.GetLineIndex(), cb.GetColumnIndex())
	}

	cb.Append(code)

	// In order to avoid parse ambiguity when files are concatenated
	// together, all files should end in a semi-colon. Do a quick
	// heuristic check if there's an obvious semi-colon already there.
	// ';' and '\n' are ASCII and never appear as trailing bytes of multi-byte UTF-8/WTF-8.
	hasSemiColon := strings.HasSuffix(code, ";") || strings.HasSuffix(code, ";\n")
	if !hasSemiColon {
		cb.Append(";")
	}

	if c.options.ShouldGatherSourceMapInfo() {
		if len(sourceAndMappings.mappingChunks) > 0 || len(sourceAndMappings.mappingSlab) > 0 {
			for _, chunk := range sourceAndMappings.mappingChunks {
				for i := range chunk {
					m := &chunk[i]
					c.sourceMap.AddMappingOfNodePositions(m.node, m.start, m.end)
				}
			}
			for i := range sourceAndMappings.mappingSlab {
				m := &sourceAndMappings.mappingSlab[i]
				c.sourceMap.AddMappingOfNodePositions(m.node, m.start, m.end)
			}
			sourceAndMappings.Release()
		} else if sourceAndMappings.mappings != nil {
			for _, mapping := range sourceAndMappings.mappings {
				c.sourceMap.AddMapping(mapping)
			}
		}
	}
}

// ToSourceOfNode generates JavaScript source code for an AST, doesn't
// generate source map info.
//
// Java: Compiler.toSource(Node)
func (c *Compiler) ToSourceOfNode(n *rhino.Node) string {
	c.InitCompilerOptionsIfTesting()
	if n != nil && !n.IsRoot() && !n.IsScript() {
		return c.toSourceWithoutMappings(n, false, nil)
	}
	sb := jutil.NewStringBuilder()
	// We don't know whether this node is part of a single JS binary, or part of a module, or
	// something else. Use the legacy license tracking implementation to match existing expectations
	// about what licenses are included in the resulting JS.
	lt := NewCompiler_ScriptNodeLicensesOnlyTracker(c)
	code := c.toSourceWithoutMappings(n, false, lt)
	if emitted := lt.EmitLicenses(); emitted != nil && !emitted.IsEmpty() {
		for license := range emitted.All() {
			sb.Append("/*\n").Append(license).Append("*/\n")
		}
	}
	sb.Append(code)
	return sb.String()
}

// compiler_SeenSetLicenseTracker is a license tracker implementation that
// maintains two sets of licenses: the "currently seen set" which will be
// emitted on the next call to emitLicenses, and the "globally seen set"
// which includes every license ever seen. Licenses in the global set are
// never added to the currently seen set. On every call to emitLicenses, the
// currently seen set is added to the globally seen set, and the currently
// seen set is cleared. There is no way to clear the globally seen set.
//
// Java: Compiler.SeenSetLicenseTracker
type compiler_SeenSetLicenseTracker interface {
	CodePrinter_LicenseTracker
	// Java: Compiler.SeenSetLicenseTracker.shouldUseLicenseInfo(Node)
	ShouldUseLicenseInfo(node *rhino.Node) bool
}

// compiler_SeenSetLicenseTrackerBase holds the fields and base
// implementations of SeenSetLicenseTracker.
type compiler_SeenSetLicenseTrackerBase struct {
	self                   compiler_SeenSetLicenseTracker
	compiler               AbstractCompiler
	globallyUniqueLicenses *jutil.LinkedHashSet[string]
	currentlySeenLicenses  *jutil.LinkedHashSet[string]
	lastSeenFile           string
}

// Java: Compiler.SeenSetLicenseTracker(AbstractCompiler)
func (t *compiler_SeenSetLicenseTrackerBase) initCompiler_SeenSetLicenseTrackerBase(self compiler_SeenSetLicenseTracker, compiler AbstractCompiler) {
	t.self = self
	t.compiler = compiler
	t.globallyUniqueLicenses = jutil.NewLinkedHashSet[string]()
	t.currentlySeenLicenses = jutil.NewLinkedHashSet[string]()
	t.lastSeenFile = ""
}

// Java: Compiler.SeenSetLicenseTracker.shouldUseLicenseInfo(Node)
func (t *compiler_SeenSetLicenseTrackerBase) ShouldUseLicenseInfo(node *rhino.Node) bool {
	return !node.IsRoot() && !node.IsScript()
}

// Java: Compiler.SeenSetLicenseTracker.trackLicensesForNode(Node)
func (t *compiler_SeenSetLicenseTrackerBase) TrackLicensesForNode(node *rhino.Node) {
	if !t.self.ShouldUseLicenseInfo(node) {
		return
	}
	file := node.GetSourceFileName()
	if file == nil {
		return
	}
	t.trackLicensesForFileName(*file)
}

func (t *compiler_SeenSetLicenseTrackerBase) trackLicensesForFileName(file string) {
	if file == t.lastSeenFile {
		// Skip files we just saw.
		return
	} else {
		t.lastSeenFile = file
	}
	license := compiler_getLicenseForFileName(t.compiler, file)
	if license == nil {
		return
	}
	if t.globallyUniqueLicenses.Contains(*license) {
		return
	}
	t.currentlySeenLicenses.Add(*license)
}

var compiler_emptyLicenseSet = jutil.NewLinkedHashSet[string]()

// Java: Compiler.SeenSetLicenseTracker.emitLicenses()
func (t *compiler_SeenSetLicenseTrackerBase) EmitLicenses() *jutil.LinkedHashSet[string] {
	if t.currentlySeenLicenses.IsEmpty() {
		return compiler_emptyLicenseSet
	}
	slice := t.currentlySeenLicenses.ToSlice()
	rv := jutil.NewLinkedHashSet[string](slice...)
	t.globallyUniqueLicenses.AddAll(slice...)
	t.currentlySeenLicenses.Clear()
	return rv
}

// Compiler_ScriptNodeLicensesOnlyTracker is a license tracker
// implementation that emits all licenses attached to script nodes,
// regardless of if the content of that script is used anywhere in the
// compilation.
//
// This license tracker encapsulates the logic that the compiler used
// before the concept of license tracking was implemented, and should only
// used in places where a better implementation cannot be used. See
// SingleBinaryLicenseTracker and ChunkGraphAwareLicenseTracker for smarter
// implementations of a license tracker for the single-binary and
// multiple-chunk emit use-cases.
//
// Java: Compiler.ScriptNodeLicensesOnlyTracker
type Compiler_ScriptNodeLicensesOnlyTracker struct {
	compiler_SeenSetLicenseTrackerBase
}

// Java: Compiler.ScriptNodeLicensesOnlyTracker(AbstractCompiler)
func NewCompiler_ScriptNodeLicensesOnlyTracker(compiler AbstractCompiler) *Compiler_ScriptNodeLicensesOnlyTracker {
	t := &Compiler_ScriptNodeLicensesOnlyTracker{}
	t.initCompiler_SeenSetLicenseTrackerBase(t, compiler)
	return t
}

// Java: Compiler.ScriptNodeLicensesOnlyTracker.shouldUseLicenseInfo(Node)
func (t *Compiler_ScriptNodeLicensesOnlyTracker) ShouldUseLicenseInfo(node *rhino.Node) bool {
	// Only use license information that is found on the script or root nodes of the tree.
	return node.IsRoot() || node.IsScript()
}

// Compiler_SingleBinaryLicenseTracker is a license tracker implementation
// that only retains licenses for which there are useful nodes in a binary.
// This will ignore licenses attached to script or root nodes when there
// are no other nodes from that script in use in the binary (likely removed
// due to compiler optimizations).
//
// Java: Compiler.SingleBinaryLicenseTracker
type Compiler_SingleBinaryLicenseTracker struct {
	compiler_SeenSetLicenseTrackerBase
}

// Java: Compiler.SingleBinaryLicenseTracker(AbstractCompiler)
func newCompiler_SingleBinaryLicenseTracker(compiler AbstractCompiler) *Compiler_SingleBinaryLicenseTracker {
	t := &Compiler_SingleBinaryLicenseTracker{}
	t.initCompiler_SeenSetLicenseTrackerBase(t, compiler)
	return t
}

// Java: Compiler.SingleBinaryLicenseTracker.shouldUseLicenseInfo(Node)
func (t *Compiler_SingleBinaryLicenseTracker) ShouldUseLicenseInfo(node *rhino.Node) bool {
	// Ignore the licensing information attached to root or script nodes, and only carry license
	// info for nodes that emit visible JS.
	return !node.IsRoot() && !node.IsScript()
}

// Compiler_ChunkGraphAwareLicenseTracker is an implementation of a license
// tracker to be used when licenses can be de-duped among Chunks (because
// licenses might be loaded by a chunk's transitive dependencies. This
// tracker has a specific lifecycle to function correctly:
//
//   - all chunks must be processed in dependency-order
//   - setCurrentChunkContext must be called before any nodes in a Chunk are
//     processed
//   - emitLicenses can either be called once per chunk, or once per input
//     file in the chunk. Each call to emitLicenses will push all
//     newly-seen licenses into tracker-internal storage, and those licenses
//     won't be emitted again until setCurrentChunkContext resets this
//     internal state.
//
// Java: Compiler.ChunkGraphAwareLicenseTracker
type Compiler_ChunkGraphAwareLicenseTracker struct {
	compiler           AbstractCompiler
	licensesFromChunks *jutil.LinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]]
	cumulativeLicenses map[*JSChunk]*jutil.LinkedHashSet[string]

	haveInitializedCurrentChunkLicenses bool
	currentChunkLicencesInTDeps         *jutil.LinkedHashSet[string]
	currentChunkTDepsRef                *jutil.LinkedHashSet[string]

	lastSeenFile             string
	licensesNewInCurrentFile *jutil.LinkedHashSet[string]

	// Not a final set as this is re-initialized in setCurrentChunkContext, and then stored into
	// the licensesFromChunks map.
	licensesNewInCurrentChunk *jutil.LinkedHashSet[string]
	log                       diagnostic.LogFile
	currentChunk              *JSChunk
}

// Java: Compiler.ChunkGraphAwareLicenseTracker(AbstractCompiler)
func NewCompiler_ChunkGraphAwareLicenseTracker(compiler AbstractCompiler) *Compiler_ChunkGraphAwareLicenseTracker {
	tdeps := jutil.NewLinkedHashSet[string]()
	return &Compiler_ChunkGraphAwareLicenseTracker{
		licensesFromChunks:                  jutil.NewLinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]](),
		cumulativeLicenses:                  make(map[*JSChunk]*jutil.LinkedHashSet[string], 16),
		haveInitializedCurrentChunkLicenses: false,
		currentChunkLicencesInTDeps:         tdeps,
		currentChunkTDepsRef:                tdeps,
		lastSeenFile:                        "",
		licensesNewInCurrentFile:            jutil.NewLinkedHashSet[string](),
		licensesNewInCurrentChunk:           compiler_emptyLicenseSet,
		log:                                 diagnostic.LogFile_CreateNoOp(),
		compiler:                            compiler,
	}
}

// Java: Compiler.ChunkGraphAwareLicenseTracker.setLogFile(LogFile)
func (t *Compiler_ChunkGraphAwareLicenseTracker) SetLogFile(file diagnostic.LogFile) {
	t.log = file
}

// Java: Compiler.ChunkGraphAwareLicenseTracker.setCurrentChunkContext(JSChunk)
func (t *Compiler_ChunkGraphAwareLicenseTracker) SetCurrentChunkContext(chunk *JSChunk) {
	if t.log.IsLogging() {
		t.log.LogFormat("Initializing licenses from deps for chunk %s", chunk)
	}
	t.currentChunk = chunk
	if !t.currentChunkLicencesInTDeps.IsEmpty() {
		t.currentChunkLicencesInTDeps.Clear()
	}
	t.currentChunkTDepsRef = compiler_emptyLicenseSet
	if !t.licensesNewInCurrentFile.IsEmpty() {
		t.licensesNewInCurrentFile.Clear()
	}
	t.haveInitializedCurrentChunkLicenses = false
	if t.licensesFromChunks.ContainsKey(chunk) {
		panic(jutil.NewIllegalStateExceptionMsg("Visiting a chunk more than once is not allowed."))
	}
	t.licensesNewInCurrentChunk = compiler_emptyLicenseSet
	t.licensesFromChunks.Put(chunk, compiler_emptyLicenseSet)
}

func (t *Compiler_ChunkGraphAwareLicenseTracker) ensureCumulativeLicenses(chunk *JSChunk) *jutil.LinkedHashSet[string] {
	if cum, ok := t.cumulativeLicenses[chunk]; ok {
		return cum
	}
	chunkNew, visited := t.licensesFromChunks.Lookup(chunk)
	if !visited {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Chunk aware license analysis error - analysis out of order."))
	}
	deps := chunk.deps
	if len(deps) == 0 {
		t.cumulativeLicenses[chunk] = chunkNew
		return chunkNew
	}
	if len(deps) == 1 {
		cum0 := t.ensureCumulativeLicenses(deps[0])
		if chunkNew.IsEmpty() {
			t.cumulativeLicenses[chunk] = cum0
			return cum0
		}
		if cum0.IsEmpty() {
			t.cumulativeLicenses[chunk] = chunkNew
			return chunkNew
		}
		merged := jutil.NewLinkedHashSet[string](cum0.ToSlice()...)
		merged.AddAll(chunkNew.ToSlice()...)
		t.cumulativeLicenses[chunk] = merged
		return merged
	}
	var shared *jutil.LinkedHashSet[string]
	needMerge := false
	for _, dep := range deps {
		cumDep := t.ensureCumulativeLicenses(dep)
		if cumDep.IsEmpty() {
			continue
		}
		if shared == nil {
			shared = cumDep
		} else if shared != cumDep {
			needMerge = true
		}
	}
	if !needMerge && chunkNew.IsEmpty() {
		if shared == nil {
			shared = compiler_emptyLicenseSet
		}
		t.cumulativeLicenses[chunk] = shared
		return shared
	}
	merged := jutil.NewLinkedHashSet[string]()
	for _, dep := range deps {
		cumDep := t.cumulativeLicenses[dep]
		if !cumDep.IsEmpty() {
			merged.AddAll(cumDep.ToSlice()...)
		}
	}
	if !chunkNew.IsEmpty() {
		merged.AddAll(chunkNew.ToSlice()...)
	}
	t.cumulativeLicenses[chunk] = merged
	return merged
}

// Java: Compiler.ChunkGraphAwareLicenseTracker.trackLicensesForNode(Node)
func (t *Compiler_ChunkGraphAwareLicenseTracker) TrackLicensesForNode(node *rhino.Node) {
	if node.IsRoot() || node.IsScript() {
		// Neither of these node kinds on their own emit code, so skip their analysis.
		return
	}
	sourceFile := node.GetSourceFileName()
	if sourceFile == nil {
		return
	}
	t.trackLicensesForFileName(*sourceFile)
}

func (t *Compiler_ChunkGraphAwareLicenseTracker) trackLicensesForFileName(sourceFile string) {
	if t.lastSeenFile == sourceFile {
		// We just saw this file during the previous call to trackLicensesForNode
		return
	} else {
		// Record the file whose license info we are going to analyze so that if we see this file
		// again in consecutive calls to trackLicensesForNode we can skip redundant analysis.
		t.lastSeenFile = sourceFile
	}
	// Lazily initialize the set of licenses provided by the chunk's transitive dependencies
	if !t.haveInitializedCurrentChunkLicenses {
		if t.log.IsLogging() {
			for depChunk := range t.currentChunk.GetAllDependencies().All() {
				t.log.LogFormat("Chunk %s depends on chunk %s", t.currentChunk, depChunk)
				if !t.licensesFromChunks.ContainsKey(depChunk) {
					panic(jutil.NewIllegalStateExceptionMsg(
						"Chunk aware license analysis error - analysis out of order."))
				}
				depSet := t.licensesFromChunks.Get(depChunk)
				if !depSet.IsEmpty() {
					t.currentChunkLicencesInTDeps.AddAll(depSet.ToSlice()...)
				}
			}
			t.currentChunkTDepsRef = t.currentChunkLicencesInTDeps
		} else {
			deps := t.currentChunk.deps
			if len(deps) == 0 {
				t.currentChunkTDepsRef = compiler_emptyLicenseSet
			} else if len(deps) == 1 {
				t.currentChunkTDepsRef = t.ensureCumulativeLicenses(deps[0])
			} else {
				var shared *jutil.LinkedHashSet[string]
				needMerge := false
				for _, dep := range deps {
					cumDep := t.ensureCumulativeLicenses(dep)
					if cumDep.IsEmpty() {
						continue
					}
					if shared == nil {
						shared = cumDep
					} else if shared != cumDep {
						needMerge = true
					}
				}
				if !needMerge {
					if shared == nil {
						shared = compiler_emptyLicenseSet
					}
					t.currentChunkTDepsRef = shared
				} else {
					for _, dep := range deps {
						cumDep := t.cumulativeLicenses[dep]
						if !cumDep.IsEmpty() {
							t.currentChunkLicencesInTDeps.AddAll(cumDep.ToSlice()...)
						}
					}
					t.currentChunkTDepsRef = t.currentChunkLicencesInTDeps
				}
			}
		}
		t.haveInitializedCurrentChunkLicenses = true
	}
	license := compiler_getLicenseForFileName(t.compiler, sourceFile)
	if license == nil {
		return
	}
	if t.currentChunkTDepsRef.Contains(*license) {
		// This node's license is covered by one of its transitive dependencies - no need to
		// duplicate it.
		return
	}
	if t.licensesNewInCurrentChunk.Contains(*license) {
		// This license is new as of this chunk, but another file has already seen it.
		return
	}
	newlyAddedLicense := t.licensesNewInCurrentFile.Add(*license)
	if t.log.IsLogging() {
		if !newlyAddedLicense {
			t.log.LogFormat(
				"Chunk %s already depends on license\n==================\n%s\n==================\n\n",
				t.currentChunk, *license)
		} else {
			t.log.LogFormat(
				"Chunk %s has new license:\n==================\n%s\n==================\n\n",
				t.currentChunk, *license)
		}
	}
}

type compiler_ParallelScriptLicenseRecorder struct {
	onlyScript      bool
	lastSeenFilePtr *string
	lastSeenFile    string
	hasLastSeen     bool
	seenFiles       []string
}

func (r *compiler_ParallelScriptLicenseRecorder) Reset(onlyScript bool) {
	r.onlyScript = onlyScript
	r.lastSeenFilePtr = nil
	r.lastSeenFile = ""
	r.hasLastSeen = false
	r.seenFiles = r.seenFiles[:0]
}

func (r *compiler_ParallelScriptLicenseRecorder) TrackLicensesForNode(node *rhino.Node) {
	isRootOrScript := node.IsRoot() || node.IsScript()
	if r.onlyScript != isRootOrScript {
		return
	}
	file := node.GetSourceFileName()
	if file == nil {
		return
	}
	if r.hasLastSeen && (file == r.lastSeenFilePtr || *file == r.lastSeenFile) {
		return
	}
	r.lastSeenFilePtr = file
	r.lastSeenFile = *file
	r.hasLastSeen = true
	r.seenFiles = append(r.seenFiles, *file)
}

func (r *compiler_ParallelScriptLicenseRecorder) EmitLicenses() *jutil.LinkedHashSet[string] {
	return nil
}

func compiler_licenseTrackerReplay(lt CodePrinter_LicenseTracker) (replay func(string), onlyScript bool, ok bool) {
	switch t := lt.(type) {
	case *Compiler_ChunkGraphAwareLicenseTracker:
		return t.trackLicensesForFileName, false, true
	case *Compiler_SingleBinaryLicenseTracker:
		return t.trackLicensesForFileName, false, true
	case *Compiler_ScriptNodeLicensesOnlyTracker:
		return t.trackLicensesForFileName, true, true
	default:
		return nil, false, false
	}
}

// Java: Compiler.ChunkGraphAwareLicenseTracker.emitLicenses()
func (t *Compiler_ChunkGraphAwareLicenseTracker) EmitLicenses() *jutil.LinkedHashSet[string] {
	if t.licensesNewInCurrentFile.IsEmpty() {
		return compiler_emptyLicenseSet
	}
	slice := t.licensesNewInCurrentFile.ToSlice()
	licensesNewInFile := jutil.NewLinkedHashSet[string](slice...)
	if t.licensesNewInCurrentChunk == compiler_emptyLicenseSet {
		t.licensesNewInCurrentChunk = jutil.NewLinkedHashSet[string](slice...)
		t.licensesFromChunks.Put(t.currentChunk, t.licensesNewInCurrentChunk)
	} else {
		t.licensesNewInCurrentChunk.AddAll(slice...)
	}
	t.licensesNewInCurrentFile.Clear()
	return licensesNewInFile
}

// toSourceAndMappings generates JavaScript source code for an AST.
//
// Java: Compiler.toSourceAndMappings(Node, boolean, LicenseTracker)
func (c *Compiler) toSourceAndMappings(
	n *rhino.Node, firstOutput bool, licenseTracker CodePrinter_LicenseTracker) *CodePrinter_SourceAndMappings {
	sm := &CodePrinter_SourceAndMappings{}
	c.toSourceAndMappingsInto(n, firstOutput, licenseTracker, sm)
	sm.mappings = sm.iterMappings
	return sm
}

func (c *Compiler) toSourceAndMappingsInto(
	n *rhino.Node, firstOutput bool, licenseTracker CodePrinter_LicenseTracker, sm *CodePrinter_SourceAndMappings) {
	builder := CodePrinter_Builder{root: n}
	builder.SetCompilerOptions(c.options)
	builder.SetTagAsTypeSummary(c.options.ShouldGenerateTypedExterns())
	builder.SetTagAsStrict(firstOutput && c.options.ShouldEmitUseStrict())
	builder.SetLicenseTracker(licenseTracker)
	builder.buildSourceAndMappingsInto(sm, true)
}

func (c *Compiler) toSourceWithoutMappings(
	n *rhino.Node, firstOutput bool, licenseTracker CodePrinter_LicenseTracker) string {
	builder := CodePrinter_Builder{root: n}
	builder.SetCompilerOptions(c.options)
	builder.SetTagAsTypeSummary(c.options.ShouldGenerateTypedExterns())
	builder.SetTagAsStrict(firstOutput && c.options.ShouldEmitUseStrict())
	builder.SetLicenseTracker(licenseTracker)
	return builder.Build()
}

// Java: Compiler.getLicenseForFile(AbstractCompiler, String)
func Compiler_GetLicenseForFile(compiler AbstractCompiler, fileName *string) *string {
	if fileName == nil {
		return nil
	}
	return compiler_getLicenseForFileName(compiler, *fileName)
}

func compiler_getLicenseForFileName(compiler AbstractCompiler, fileName string) *string {
	licenseRoot := compiler.GetScriptNode(fileName)
	if licenseRoot == nil {
		return nil
	}
	docInfo := licenseRoot.GetJSDocInfo()
	if docInfo == nil {
		return nil
	}
	return docInfo.GetLicense()
}

// ToSourceArray converts the parse tree for each input in a chunk back to
// JS code, using the given License Tracker implementation to decide which
// licenses should be emitted before each input file.
//
// licenseTracker: The license tracker implementation to use.
// ChunkGraphAwareLicenseTracker is a suitable implementation when this
// method is being called on each chunk in the chunk graph in dependency
// order - see its javadoc on how to use it correctly. chunk: the chunk
// being converted to source.
//
// Java: Compiler.toSourceArray(LicenseTracker, JSChunk)
func (c *Compiler) ToSourceArray(licenseTracker CodePrinter_LicenseTracker, chunk *JSChunk) []string {
	return Compiler_RunInCompilerThread(c,
		func() []string {
			c.evictPassArenaCaches()
			inputs := chunk.GetInputs()
			numInputs := len(inputs)
			if numInputs == 0 {
				return []string{}
			}

			sources := make([]string, numInputs)
			cb := NewCompiler_CodeBuilder()
			for i := 0; i < numInputs; i++ {
				scriptNode := inputs[i].GetAstRoot(c)
				if scriptNode == nil {
					panic(jutil.NewIllegalArgumentExceptionMsg("Bad chunk input: " + inputs[i].GetName()))
				}

				cb.Reset()
				c.ToSourceOfCodeBuilder(cb, licenseTracker, i, scriptNode)
				sources[i] = cb.String()
			}
			return sources
		})
}

// Compiler_CodeBuilder stores a buffer of text to which more can be
// appended. This is just like a StringBuilder except that we also track the
// number of lines.
//
// Java: Compiler.CodeBuilder
type Compiler_CodeBuilder struct {
	buf         []byte
	utf16Len    int
	lineCount   int
	colCount    int
	stringTaken bool
}

// Java: Compiler.CodeBuilder()
func NewCompiler_CodeBuilder() *Compiler_CodeBuilder {
	return &Compiler_CodeBuilder{}
}

// Reset removes all text, but leaves the line count unchanged.
//
// Java: Compiler.CodeBuilder.reset()
func (b *Compiler_CodeBuilder) Reset() {
	if b.stringTaken {
		b.buf = nil
		b.stringTaken = false
	} else {
		b.buf = b.buf[:0]
	}
	b.utf16Len = 0
}

// Append appends the given string to the text buffer.
//
// Java: Compiler.CodeBuilder.append(String)
func (b *Compiler_CodeBuilder) Append(str string) *Compiler_CodeBuilder {
	if len(str) == 0 {
		return b
	}
	if b.stringTaken {
		newCap := max(len(b.buf)+len(str), max(cap(b.buf)*2, 4096))
		newBuf := make([]byte, len(b.buf), newCap)
		copy(newBuf, b.buf)
		b.buf = newBuf
		b.stringTaken = false
	} else if len(b.buf)+len(str) > cap(b.buf) {
		newCap := max(len(b.buf)+len(str), max(cap(b.buf)*2, 4096))
		newBuf := make([]byte, len(b.buf), newCap)
		copy(newBuf, b.buf)
		b.buf = newBuf
	}
	b.buf = append(b.buf, str...)

	// Adjust the line and column information for the new text.
	//
	// Go-only: performance — Java loops over str.indexOf('\n', index + 1),
	// which here would convert all of str to UTF-16 once per line (quadratic
	// in the size of a script). '\n' never occurs inside a multi-byte WTF-8
	// sequence, so counting and finding it on the bytes gives the same line
	// count, and the UTF-16 length after the last '\n' is
	// str.length() - (lastIndex + 1).
	lastIndex := strings.LastIndexByte(str, '\n')
	if lastIndex == -1 {
		// No new lines, append the new characters added.
		n := jutil.StringLength(str)
		b.utf16Len += n
		b.colCount += n
	} else {
		tailLen := jutil.StringLength(str[lastIndex+1:])
		b.utf16Len += jutil.StringLength(str[:lastIndex]) + 1 + tailLen
		b.lineCount += 1 + strings.Count(str[:lastIndex], "\n")
		b.colCount = tailLen
	}

	return b
}

// String returns all text in the text buffer.
//
// Java: Compiler.CodeBuilder.toString()
func (b *Compiler_CodeBuilder) String() string {
	if len(b.buf) == 0 {
		return ""
	}
	b.stringTaken = true
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}

// GetLength returns the length of the text buffer.
//
// Java: Compiler.CodeBuilder.getLength()
func (b *Compiler_CodeBuilder) GetLength() int {
	return b.utf16Len
}

// GetLineIndex returns the (zero-based) index of the last line in the text
// buffer.
//
// Java: Compiler.CodeBuilder.getLineIndex()
func (b *Compiler_CodeBuilder) GetLineIndex() int {
	return b.lineCount
}

// GetColumnIndex returns the (zero-based) index of the last column in the
// text buffer.
//
// Java: Compiler.CodeBuilder.getColumnIndex()
func (b *Compiler_CodeBuilder) GetColumnIndex() int {
	return b.colCount
}

// EndsWith determines whether the text ends with the given suffix.
//
// Java: Compiler.CodeBuilder.endsWith(String)
func (b *Compiler_CodeBuilder) EndsWith(suffix string) bool {
	return b.utf16Len > jutil.StringLength(suffix) &&
		len(b.buf) >= len(suffix) &&
		unsafe.String(unsafe.SliceData(b.buf[len(b.buf)-len(suffix):]), len(suffix)) == suffix
}

// ------------------------------------------------------------------------
// Optimizations
// ------------------------------------------------------------------------

// Java: Compiler.performTranspilationAndOptimizations(SegmentOfCompilationToRun)
func (c *Compiler) PerformTranspilationAndOptimizations(segmentOfCompilationToRun CompilerOptions_SegmentOfCompilationToRun) {
	jutil.CheckState(c.options.ShouldOptimize())
	// getOptimizations() also includes transpilation passes
	allOptimizationPassesToRun := c.GetPassConfig().GetOptimizations().Build()
	if len(allOptimizationPassesToRun) == 0 {
		return
	}

	// Get the subset of optimization passes to run in the current segment of optimizations.
	optimizationPassesToRunInCurrentSegment := c.getOptimizationPassesToRunInCurrentSegment(
		segmentOfCompilationToRun, allOptimizationPassesToRun)

	// note whether any script gets transpiled
	c.markTranspiledFiles()

	c.phaseOptimizer = c.createPhaseOptimizer()
	c.phaseOptimizer.Consume(optimizationPassesToRunInCurrentSegment)
	c.phaseOptimizer.Process(c.externsRoot, c.jsRoot)
	c.phaseOptimizer = nil
}

// Java: Compiler.performFinalizations()
func (c *Compiler) PerformFinalizations() {
	finalizations := c.GetPassConfig().GetFinalizations().Build()
	if len(finalizations) == 0 {
		return
	}

	c.phaseOptimizer = c.createPhaseOptimizer()
	c.phaseOptimizer.Consume(finalizations)
	c.phaseOptimizer.Process(c.externsRoot, c.jsRoot)
	c.phaseOptimizer = nil
}

// ComputeCFG is Control Flow Analysis.
//
// Java: Compiler.computeCFG()
func (c *Compiler) ComputeCFG() ControlFlowGraph[*rhino.Node] {
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

// getOptimizationPassesToRunInCurrentSegment takes the entire list of
// optimization passes to run in the current build, and returns the subset
// of passes that should be run in this current segment of optimizations.
//
// The 3 cases are:
//
//  1. If we are running `SegmentOfCompilationToRun.OPTIMIZATIONS`, this
//     will return the entire list of optimization passes.
//  2. If we are running `SegmentOfCompilationToRun.OPTIMIZATIONS_FIRST_HALF`,
//     this will return all passes up to the half-way point marker pass
//     `OPTIMIZATIONS_HALFWAY_POINT`.
//  3. If we are running `SegmentOfCompilationToRun.OPTIMIZATIONS_SECOND_HALF`,
//     this will return all passes after the half-way point marker pass
//     `OPTIMIZATIONS_HALFWAY_POINT`.
//
// Java: Compiler.getOptimizationPassesToRunInCurrentSegment(SegmentOfCompilationToRun, List<PassFactory>)
func (c *Compiler) getOptimizationPassesToRunInCurrentSegment(
	segmentOfCompilationToRun CompilerOptions_SegmentOfCompilationToRun,
	allOptimizationPassesToRun []*PassFactory) []*PassFactory {

	// Create a list of passes based on which segment of optimizations we are running.
	switch segmentOfCompilationToRun {
	case CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS:
		return allOptimizationPassesToRun
	case CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF:
		return c.passListUntil(allOptimizationPassesToRun, PassNames_OPTIMIZATIONS_HALFWAY_POINT)
	case CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF:
		return c.passListAfter(allOptimizationPassesToRun, PassNames_OPTIMIZATIONS_HALFWAY_POINT)
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Unsupported segment of optimizations to run: " + segmentOfCompilationToRun.String()))
	}
}

// passListUntil returns a subset of passes from the given list of passes,
// up to but not including the given pass name.
//
// Java: Compiler.passListUntil(List<PassFactory>, String)
func (c *Compiler) passListUntil(passes []*PassFactory, passName string) []*PassFactory {
	passList := []*PassFactory{}
	for _, pass := range passes {
		if pass.GetName() == passName {
			return passList
		}
		passList = append(passList, pass)
	}
	return passList
}

// passListAfter returns a subset of passes from the given list of passes,
// starting right after the given pass name.
//
// Java: Compiler.passListAfter(List<PassFactory>, String)
func (c *Compiler) passListAfter(passes []*PassFactory, passName string) []*PassFactory {
	passList := []*PassFactory{}
	for _, pass := range passes {
		if pass.GetName() == passName {
			passList = passList[:0]
			continue
		}
		passList = append(passList, pass)
	}
	return passList
}

// Java: Compiler.SYNTHETIC_FILE_NAME_PREFIX
const compiler_SYNTHETIC_FILE_NAME_PREFIX = " [synthetic:"

// Java: Compiler.SYNTHETIC_CODE_INPUT_ID
var compiler_SYNTHETIC_CODE_INPUT_ID = rhino.NewInputId(compiler_SYNTHETIC_FILE_NAME_PREFIX + "input] ")

// Java: Compiler.getChangeTracker()
func (c *Compiler) GetChangeTracker() *ChangeTracker {
	return c.changeTracker
}

// Java: Compiler.addIndexProvider(IndexProvider<?>)
func (c *Compiler) AddIndexProvider(indexProvider IndexProvider[any]) {
	type_ := indexProvider.GetType()
	if c.indexProvidersByType.Put(type_, indexProvider) != nil {
		panic(jutil.NewIllegalStateExceptionMsg(
			"A provider is already registered for index of type " + type_.Name()))
	}
}

// Java: Compiler.getIndex(Class<T>)
func (c *Compiler) GetIndex(key reflect.Type) any {
	indexProvider := c.indexProvidersByType.Get(key)
	if indexProvider == nil {
		return nil
	}
	return indexProvider.Get()
}

// Java: Compiler.getExternsRoot()
func (c *Compiler) GetExternsRoot() *rhino.Node {
	return c.externsRoot
}

// Java: Compiler.getJsRoot()
func (c *Compiler) GetJsRoot() *rhino.Node {
	return c.jsRoot
}

// Java: Compiler.hasScopeChanged(Node)
func (c *Compiler) HasScopeChanged(n *rhino.Node) bool {
	if c.phaseOptimizer == nil {
		return true
	}
	return c.phaseOptimizer.HasScopeChanged(n)
}

// Java: Compiler.reportChangeToChangeScope(Node)
func (c *Compiler) ReportChangeToChangeScope(changeScopeRoot *rhino.Node) {
	c.changeTracker.ReportChangeToChangeScope(changeScopeRoot)
}

// Java: Compiler.reportFunctionDeleted(Node)
func (c *Compiler) ReportFunctionDeleted(n *rhino.Node) {
	c.changeTracker.ReportFunctionDeleted(n)
}

// Java: Compiler.reportDisambiguatePropertiesSummary(Supplier<String>)
func (c *Compiler) ReportDisambiguatePropertiesSummary(summarySupplier func() string) {
	if c.tracker != nil {
		c.tracker.SetDisambiguatePropertiesSummary(summarySupplier())
	}
}

// Java: Compiler.reportAmbiguatePropertiesSummary(Supplier<String>)
func (c *Compiler) ReportAmbiguatePropertiesSummary(summarySupplier func() string) {
	if c.tracker != nil {
		c.tracker.SetAmbiguatePropertiesSummary(summarySupplier())
	}
}

// Java: Compiler.reportChangeToEnclosingScope(Node)
func (c *Compiler) ReportChangeToEnclosingScope(n *rhino.Node) {
	c.changeTracker.ReportChangeToEnclosingScope(n)
}

// Java: Compiler.getCodingConvention()
func (c *Compiler) GetCodingConvention() CodingConvention {
	convention := c.options.GetCodingConvention()
	if convention == nil {
		convention = c.defaultCodingConvention
	}
	return convention
}

// Java: Compiler.getParserConfigLanguageMode(CompilerOptions.LanguageMode)
func (c *Compiler) getParserConfigLanguageMode(languageMode CompilerOptions_LanguageMode) parsing.Config_LanguageMode {
	switch languageMode {
	case CompilerOptions_LanguageMode_ECMASCRIPT3:
		return parsing.Config_LanguageMode_ECMASCRIPT3
	case CompilerOptions_LanguageMode_ECMASCRIPT5, CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT:
		return parsing.Config_LanguageMode_ECMASCRIPT5
	case CompilerOptions_LanguageMode_ECMASCRIPT_2015:
		return parsing.Config_LanguageMode_ECMASCRIPT_2015
	case CompilerOptions_LanguageMode_ECMASCRIPT_2016:
		return parsing.Config_LanguageMode_ECMASCRIPT_2016
	case CompilerOptions_LanguageMode_ECMASCRIPT_2017:
		return parsing.Config_LanguageMode_ECMASCRIPT_2017
	case CompilerOptions_LanguageMode_ECMASCRIPT_2018:
		return parsing.Config_LanguageMode_ECMASCRIPT_2018
	case CompilerOptions_LanguageMode_ECMASCRIPT_2019:
		return parsing.Config_LanguageMode_ECMASCRIPT_2019
	case CompilerOptions_LanguageMode_ECMASCRIPT_2020:
		return parsing.Config_LanguageMode_ECMASCRIPT_2020
	case CompilerOptions_LanguageMode_ECMASCRIPT_2021:
		return parsing.Config_LanguageMode_ECMASCRIPT_2021
	case CompilerOptions_LanguageMode_ECMASCRIPT_NEXT:
		return parsing.Config_LanguageMode_ES_NEXT
	case CompilerOptions_LanguageMode_UNSTABLE:
		return parsing.Config_LanguageMode_UNSTABLE
	case CompilerOptions_LanguageMode_UNSUPPORTED:
		return parsing.Config_LanguageMode_UNSUPPORTED
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected language mode: " + c.options.GetLanguageIn().String()))
	}
}

// Java: Compiler.getParserConfig(ConfigContext)
func (c *Compiler) GetParserConfig(context AbstractCompiler_ConfigContext) *parsing.Config {
	if c.parserConfig == nil || c.externsParserConfig == nil {
		func() {
			c.monitor.Lock()
			defer c.monitor.Unlock()
			if c.parserConfig == nil {
				configLanguageMode := c.getParserConfigLanguageMode(c.options.GetLanguageIn())
				strictMode := parsing.Config_StrictMode_SLOPPY
				if c.options.ExpectStrictModeInput() {
					strictMode = parsing.Config_StrictMode_STRICT
				}
				c.parserConfig = c.CreateConfig(configLanguageMode, strictMode)
				// Externs must always be parsed with at least ES5 language mode.
				if configLanguageMode == parsing.Config_LanguageMode_ECMASCRIPT3 {
					c.externsParserConfig = c.CreateConfig(parsing.Config_LanguageMode_ECMASCRIPT5, strictMode)
				} else {
					c.externsParserConfig = c.parserConfig
				}
			}
		}()
	}
	if context == AbstractCompiler_ConfigContext_EXTERNS {
		return c.externsParserConfig
	}
	return c.parserConfig
}

// Java: Compiler.createConfig(LanguageMode, StrictMode)
func (c *Compiler) CreateConfig(mode parsing.Config_LanguageMode, strictMode parsing.Config_StrictMode) *parsing.Config {
	runMode := parsing.Config_RunMode_STOP_AFTER_ERROR
	if c.options.CanContinueAfterErrors() {
		runMode = parsing.Config_RunMode_KEEP_GOING
	}
	return parsing.ParserRunner_CreateConfigWithModes(
		mode,
		c.options.IsParseJsDocDocumentation(),
		runMode,
		c.options.GetExtraAnnotationNames(),
		c.options.GetParseInlineSourceMaps(),
		strictMode)
}

// ------------------------------------------------------------------------
// Error reporting
// ------------------------------------------------------------------------

// GetDiagnosticGroups: the warning classes that are available from the
// command-line, and are suppressible by the @suppress annotation.
//
// Java: Compiler.getDiagnosticGroups()
func (c *Compiler) GetDiagnosticGroups() *DiagnosticGroups {
	return NewDiagnosticGroups()
}

// Java: Compiler.report(JSError)
func (c *Compiler) Report(error *JSError) {
	level := error.DefaultLevel().Ptr()
	if c.warningsGuard != nil {
		newLevel := c.warningsGuard.Level(error)
		if newLevel != nil {
			level = newLevel
		}
	}

	if level.IsOn() {
		c.InitCompilerOptionsIfTesting()
		if c.GetOptions().GetErrorHandler() != nil {
			c.GetOptions().GetErrorHandler().Report(level, error)
		}
		c.errorManager.Report(level, error)
	}
}

// ReportWithLevel is report(CheckLevel, JSError), Compiler's ErrorHandler
// method (see AsErrorHandler).
//
// Java: Compiler.report(CheckLevel, JSError)
func (c *Compiler) ReportWithLevel(ignoredLevel *CheckLevel, error *JSError) {
	c.Report(error)
}

// Java: Compiler.getErrorLevel(JSError)
func (c *Compiler) GetErrorLevel(error *JSError) *CheckLevel {
	jutil.CheckNotNull(c.options)
	return c.warningsGuard.Level(error)
}

// ThrowInternalError reports an internal error.
//
// Java: Compiler.throwInternalError(String, Throwable)
func (c *Compiler) ThrowInternalError(message string, cause jutil.Throwable) {
	panic(jutil.NewRuntimeExceptionWithMessageAndCause(
		"INTERNAL COMPILER ERROR.\nPlease report this problem.\n\n"+message, cause))
}

// GetErrorCount gets the number of errors.
//
// Java: Compiler.getErrorCount()
func (c *Compiler) GetErrorCount() int {
	return c.errorManager.GetErrorCount()
}

// GetWarningCount gets the number of warnings.
//
// Java: Compiler.getWarningCount()
func (c *Compiler) GetWarningCount() int {
	return c.errorManager.GetWarningCount()
}

// Java: Compiler.hasHaltingErrors()
func (c *Compiler) HasHaltingErrors() bool {
	return !c.GetOptions().CanContinueAfterErrors() && c.errorManager.HasHaltingErrors()
}

// HasErrors consults the ErrorManager to see if we've encountered errors
// that should halt compilation.
//
// If CompilerOptions#canContinueAfterErrors is true, this function always
// returns false without consulting the error manager. The error manager
// will continue to be told about new errors and warnings, but the compiler
// will complete compilation of all inputs.
//
// Java: Compiler.hasErrors()
func (c *Compiler) HasErrors() bool {
	return c.HasHaltingErrors()
}

// Java: Compiler.getSourceFileByName(String)
func (c *Compiler) GetSourceFileByName(sourceName string) *SourceFile {
	// Here we assume that the source name is the input name, this
	// is true of JavaScript parsed from source.
	input := c.inputsById.Get(rhino.NewInputId(sourceName))
	if input != nil {
		return input.GetSourceFile()
	}
	// Alternatively, the sourceName might have been reverse-mapped by
	// an input source-map, so let's look in our sourcemap original sources.
	return c.sourceMapOriginalSources.Get(sourceName)
}

// getSourceFileByNullableName is getSourceFileByName for a nullable name.
//
// Go-only: language-adapter — the Java null check of getSourceFileByName,
// for the callers whose name may be null.
func (c *Compiler) getSourceFileByNullableName(sourceName *string) *SourceFile {
	if sourceName != nil {
		return c.GetSourceFileByName(*sourceName)
	}
	return nil
}

// Java: Compiler.getSourceFileContentByName(String)
func (c *Compiler) GetSourceFileContentByName(sourceName string) *string {
	file := c.GetSourceFileByName(sourceName)
	if file.IsStubSourceFileForAlreadyProvidedInput() {
		return nil
	}
	jutil.CheckNotNull(file)
	var result *string
	jutil.Try(func() {
		code := file.GetCode()
		result = &code
	}, func(e jutil.IIOException) {
		result = nil
	})
	return result
}

// Java: Compiler.addInputSourceMap(String, SourceMapInput)
func (c *Compiler) AddInputSourceMap(sourceFileName string, inputSourceMap *SourceMapInput) {
	c.inputSourceMaps.Put(sourceFileName, inputSourceMap)
	if c.options.GetSourceMapIncludeSourcesContent() && c.sourceMap != nil {
		c.addSourceMapSourceFiles(inputSourceMap)
	}
}

// GetBase64SourceMapContents returns the <encoded_source_map> from a `//#
// sourceMappingURL=data:application/json;base64,<encoded_source_map>`
// comment. This sourceMappingURL comment is a JS file's inline source map,
// embedded into the input JS file with a base64-encoded "data url" stored
// in `//# sourceMappingURL=` comment.
//
// Java: Compiler.getBase64SourceMapContents(String)
func (c *Compiler) GetBase64SourceMapContents(sourceFileName string) *string {
	sourceMapInput := c.inputSourceMaps.GetOrDefault(sourceFileName, nil)
	// This sourceMappingURL ends with a ".inline.map", and we want to get its base4-encoded
	// contents (i.e. "data:application/json;base64,<encoded_source_map>").
	var sourceMappingURL *string
	if sourceMapInput != nil {
		originalPath := sourceMapInput.GetOriginalPath()
		sourceMappingURL = &originalPath
	}
	if sourceMappingURL != nil && strings.HasSuffix(*sourceMappingURL, ".inline.map") {
		unencoded := sourceMapInput.GetRawSourceMapContents() // This is the sourcemap.
		encoded := jutil.BaseEncoding_Base64().Encode([]byte(*unencoded))
		return &encoded
	}
	return nil
}

// addSourceMapSourceFiles adds file name to content mappings for all
// sources found in a source map. This is used to populate sourcesContent
// array in the output source map even for sources embedded in the input
// source map.
//
// Java: Compiler.addSourceMapSourceFiles(SourceMapInput)
func (c *Compiler) addSourceMapSourceFiles(inputSourceMap *SourceMapInput) {
	// synchronized annotation guards concurrent access to sourceMap during parsing.
	c.monitor.Lock()
	defer c.monitor.Unlock()
	consumer := inputSourceMap.GetSourceMap(c.errorManager)
	if consumer == nil {
		return
	}
	sourcesContent := consumer.GetOriginalSourcesContent()
	if sourcesContent == nil {
		return
	}
	content := sourcesContent
	sources := consumer.GetOriginalSources()
	i := 0
	for i < len(sources) && i < len(content) {
		code := content[i]
		source := SourceMapResolver_GetRelativePath(inputSourceMap.GetOriginalPath(), *sources[i])
		i++
		if source != nil {
			c.sourceMap.AddSourceFile(source.GetName(), code)
		}
	}
	if i < len(sources) || i < len(content) {
		panic(jutil.NewRuntimeException(
			"Source map's \"sources\" and \"sourcesContent\" lengths do not match."))
	}
}

// Java: Compiler.getSourceMapping(String, int, int)
func (c *Compiler) GetSourceMapping(
	sourceName *string, lineNumber int, columnNumber int) *sourcemap.Mapping_OriginalMapping {
	if sourceName == nil {
		return nil
	}
	result, _, hasInputMap, isDefaultEmpty := c.getSourceMappingValue(*sourceName, lineNumber, columnNumber)
	if !hasInputMap {
		return nil
	}
	if isDefaultEmpty {
		return sourcemap.Mapping_OriginalMapping_GetDefaultInstance()
	}
	return &result
}

func (c *Compiler) getSourceMappingValue(
	sourceName string, lineNumber int, columnNumber int) (sourcemap.Mapping_OriginalMapping, *string, bool, bool) {
	if c.inputSourceMaps.IsEmpty() {
		return sourcemap.Mapping_OriginalMapping{}, nil, false, false
	}
	sourceMap := c.inputSourceMaps.Get(sourceName)
	if sourceMap == nil {
		return sourcemap.Mapping_OriginalMapping{}, nil, false, false
	}

	// JSCompiler uses 1-indexing for lineNumber and 0-indexing for columnNumber.
	// Sourcemaps use 1-indexing for both.
	consumer := sourceMap.GetSourceMap(c.errorManager)
	if consumer == nil {
		return sourcemap.Mapping_OriginalMapping{}, nil, false, false
	}
	result, identifierPtr, ok := consumer.GetMappingForLineValue(lineNumber, columnNumber+1)
	if !ok {
		return sourcemap.Mapping_OriginalMapping{}, nil, true, true
	}

	// First check to see if the original file was loaded from an input source map.
	sourceMapOriginalPath := sourceMap.GetOriginalPath()
	resultOriginalPath := result.GetOriginalFile()
	var relativePath string

	// Resolving the paths to a source file is expensive, so check the cache first.
	if c.resolvedSourceMap.hasCached &&
		sourceMapOriginalPath == c.resolvedSourceMap.originalPath &&
		resultOriginalPath == c.resolvedSourceMap.sourceMapPath {
		relativePath = c.resolvedSourceMap.relativePath
	} else {
		relativePath = compiler_resolveSibling(sourceMapOriginalPath, resultOriginalPath)
		source := c.GetSourceFileByName(relativePath)
		if source == nil && resultOriginalPath != "" {
			source = SourceMapResolver_GetRelativePath(
				sourceMap.GetOriginalPath(), result.GetOriginalFile())
			if source != nil {
				c.sourceMapOriginalSources.PutIfAbsent(relativePath, source)
			}
		}

		// Cache this resolved source for the next caller.
		c.resolvedSourceMap.hasCached = true
		c.resolvedSourceMap.originalPath = sourceMapOriginalPath
		c.resolvedSourceMap.sourceMapPath = resultOriginalPath
		c.resolvedSourceMap.relativePath = relativePath
	}

	result.SetOriginalFileAndColumnInPlace(relativePath, result.GetColumnPosition()-1)
	return result, identifierPtr, true, false
}

// Java: Compiler.getSourceLine(String, int)
func (c *Compiler) GetSourceLine(sourceName *string, lineNumber int) *string {
	if lineNumber < 1 {
		return nil
	}
	input := c.getSourceFileByNullableName(sourceName)
	if input != nil && !input.IsStubSourceFileForAlreadyProvidedInput() {
		return input.GetLine(lineNumber)
	}
	return nil
}

// Java: Compiler.getSourceLines(String, int, int)
func (c *Compiler) GetSourceLines(sourceName *string, lineNumber int, length int) Region {
	if lineNumber < 1 {
		return nil
	}
	input := c.getSourceFileByNullableName(sourceName)
	if input != nil && !input.IsStubSourceFileForAlreadyProvidedInput() {
		return input.GetLines(lineNumber, length)
	}
	return nil
}

// Java: Compiler.getSourceRegion(String, int)
func (c *Compiler) GetSourceRegion(sourceName *string, lineNumber int) Region {
	if lineNumber < 1 {
		return nil
	}
	input := c.getSourceFileByNullableName(sourceName)
	if input != nil && !input.IsStubSourceFileForAlreadyProvidedInput() {
		return input.GetRegion(lineNumber)
	}
	return nil
}

// ------------------------------------------------------------------------
// Package-private and Protected helpers
// ------------------------------------------------------------------------

// Java: Compiler.getNodeForCodeInsertion(JSChunk)
func (c *Compiler) GetNodeForCodeInsertion(chunk *JSChunk) *rhino.Node {
	if c.inputsById.ContainsKey(compiler_SYNTHETIC_CODE_INPUT_ID) {
		return c.inputsById.Get(compiler_SYNTHETIC_CODE_INPUT_ID).GetAstRoot(c)
	}
	if chunk == nil {
		allInputs := c.chunkGraph.GetAllInputs()
		if c.chunkGraph == nil || len(allInputs) == 0 {
			panic(jutil.NewIllegalStateExceptionMsg("No inputs"))
		}
		firstInput := allInputs[0]
		return compiler_checkNotModule(firstInput.GetAstRoot(c), "Cannot insert code into a module")
	}

	chunkInputs := chunk.GetInputs()
	if len(chunkInputs) != 0 {
		return compiler_checkNotModule(
			chunkInputs[0].GetAstRoot(c), "Cannot insert code into a module")
	}
	panic(jutil.NewIllegalStateExceptionMsg("Root chunk has no inputs"))
}

// Java: Compiler.getSourceMap()
func (c *Compiler) GetSourceMap() *SourceMap {
	return c.sourceMap
}

// Java: Compiler.setVariableMap(VariableMap)
func (c *Compiler) SetVariableMap(variableMap *VariableMap) {
	c.variableMap = variableMap
}

// Java: Compiler.getVariableMap()
func (c *Compiler) GetVariableMap() *VariableMap {
	return c.variableMap
}

// Java: Compiler.setPropertyMap(VariableMap)
func (c *Compiler) SetPropertyMap(propertyMap *VariableMap) {
	c.propertyMap = propertyMap
}

// Java: Compiler.getPropertyMap()
func (c *Compiler) GetPropertyMap() *VariableMap {
	return c.propertyMap
}

// Java: Compiler.setStringMap(VariableMap)
func (c *Compiler) SetStringMap(stringMap *VariableMap) {
	c.stringMap = stringMap
}

// Java: Compiler.setCssNames(Set<String>)
func (c *Compiler) SetCssNames(cssNames *jutil.LinkedHashSet[string]) {
	c.cssNames = cssNames
}

// Java: Compiler.setIdGeneratorMap(String)
func (c *Compiler) SetIdGeneratorMap(serializedIdMappings string) {
	c.idGeneratorMap = &serializedIdMappings
}

// Java: Compiler.getTranspiledFiles()
func (c *Compiler) GetTranspiledFiles() bool {
	return c.transpiledFiles
}

// Java: Compiler.setTranspiledFiles(boolean)
func (c *Compiler) SetTranspiledFiles(transpiledFiles bool) {
	c.transpiledFiles = transpiledFiles
}

// Java: Compiler.getCrossChunkIdGenerator()
func (c *Compiler) GetCrossChunkIdGenerator() *IdGenerator {
	return c.crossChunkIdGenerator
}

// Java: Compiler.setAnonymousFunctionNameMap(VariableMap)
func (c *Compiler) SetAnonymousFunctionNameMap(functionMap *VariableMap) {
	// NOTE: remove this method
}

// Java: Compiler.getStringMap()
func (c *Compiler) GetStringMap() *VariableMap {
	return c.stringMap
}

// Java: Compiler.setInstrumentationMapping(VariableMap)
func (c *Compiler) SetInstrumentationMapping(instrumentationMapping *VariableMap) {
	c.instrumentationMapping = instrumentationMapping
}

// Java: Compiler.getInstrumentationMapping()
func (c *Compiler) GetInstrumentationMapping() *VariableMap {
	return c.instrumentationMapping
}

// Java: Compiler.addExportedNames(Set<String>)
func (c *Compiler) AddExportedNames(exportedNames *jutil.LinkedHashSet[string]) {
	c.exportedNames.AddAll(exportedNames.ToSlice()...)
}

// Java: Compiler.getExportedNames()
func (c *Compiler) GetExportedNames() *jutil.LinkedHashSet[string] {
	return c.exportedNames
}

// Java: Compiler.setDefineNames(Collection<String>)
func (c *Compiler) SetDefineNames(defineNames []string) {
	c.defineNames = jutil.NewLinkedHashSet[string](defineNames...)
}

// Java: Compiler.getDefineNames()
func (c *Compiler) GetDefineNames() *jutil.LinkedHashSet[string] {
	return c.defineNames
}

// Java: Compiler.getOptions()
func (c *Compiler) GetOptions() *CompilerOptions {
	return c.options
}

// Compiler_SetLoggingLevel sets the logging level for the
// com.google.javascript.jscomp package.
//
// Java: Compiler.setLoggingLevel(Level)
func Compiler_SetLoggingLevel(level *jutil.Level) {
	Compiler_logger.SetLevel(level)
}

// GetAstDotGraph gets the DOT graph of the AST generated at the end of
// compilation.
//
// Java: Compiler.getAstDotGraph()
func (c *Compiler) GetAstDotGraph() string {
	if c.jsRoot != nil {
		cfg := ControlFlowAnalysis_Builder_().
			SetCompiler(c).
			SetCfgRoot(c.jsRoot).
			SetTraverseFunctions(true).
			ComputeCfg()
		return DotFormatter_ToDotWithCfg(c.jsRoot, cfg)
	} else {
		return ""
	}
}

// Java: Compiler.getErrorManager()
func (c *Compiler) GetErrorManager() ErrorManager {
	if c.options == nil {
		c.InitOptions(NewCompilerOptions())
	}
	return c.errorManager
}

// Java: Compiler.getInputsInOrder()
func (c *Compiler) GetInputsInOrder() []*CompilerInput {
	return c.chunkGraph.GetAllInputs()
}

// Java: Compiler.getNumberOfInputs()
func (c *Compiler) GetNumberOfInputs() int {
	// In some testing cases inputs will be null, but obviously there must be at least one input.
	// The intended use of this method is to allow passes to estimate how much memory they will
	// need for data structures, so it's not necessary that the returned value be exactly right
	// in the corner cases where inputs ends up being null.
	if c.chunkGraph != nil {
		return c.chunkGraph.GetInputCount()
	}
	return 1
}

// GetInputsById returns an unmodifiable view of the compiler inputs
// indexed by id.
//
// Java: Compiler.getInputsById()
func (c *Compiler) GetInputsById() *jutil.JavaHashMap[*rhino.InputId, *CompilerInput] {
	return c.inputsById
}

// GetExternsInOrder gets the externs in the order in which they are being
// processed.
//
// Java: Compiler.getExternsInOrder()
func (c *Compiler) GetExternsInOrder() []*CompilerInput {
	return c.externs.Slice()
}

// Java: Compiler.getInputsForTesting()
func (c *Compiler) GetInputsForTesting() []*CompilerInput {
	if c.chunkGraph != nil {
		return append([]*CompilerInput(nil), c.chunkGraph.GetAllInputs()...)
	}
	return nil
}

// Java: Compiler.getExternsForTesting()
func (c *Compiler) GetExternsForTesting() *jutil.ArrayList[*CompilerInput] {
	return c.externs
}

// Java: Compiler.hasRegExpGlobalReferences()
func (c *Compiler) HasRegExpGlobalReferences() bool {
	return c.hasRegExpGlobalReferences
}

// Java: Compiler.setHasRegExpGlobalReferences(boolean)
func (c *Compiler) SetHasRegExpGlobalReferences(references bool) {
	c.hasRegExpGlobalReferences = references
}

// SetRunJ2clPasses sets whether J2CL passes should run
//
// Java: Compiler.setRunJ2clPasses(boolean)
func (c *Compiler) SetRunJ2clPasses(runJ2clPasses bool) {
	c.runJ2clPasses = runJ2clPasses
}

// RunJ2clPasses reports whether J2CL passes should run
//
// Java: Compiler.runJ2clPasses()
func (c *Compiler) RunJ2clPasses() bool {
	return c.runJ2clPasses
}

// Java: Compiler.getSynthesizedExternsInput()
func (c *Compiler) GetSynthesizedExternsInput() *CompilerInput {
	if c.syntheticExternsInput != nil {
		return c.syntheticExternsInput
	}

	input := NewCompilerInputWithIsExtern(c.SYNTHETIC_EXTERNS_FILE /* isExtern= */, true)
	root := jutil.CheckNotNull(input.GetAstRoot(c))
	c.putCompilerInput(input)
	c.syntheticExternsInput = input
	c.externsRoot.AddChildToFront(root)
	c.externs.AddAt(0, input)
	c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)

	return input
}

// Java: Compiler.getSynthesizedExternsRoot()
func (c *Compiler) getSynthesizedExternsRoot() *rhino.Node {
	externsInput := c.GetSynthesizedExternsInput()
	return externsInput.GetAstRoot(c)
}

// Java: Compiler.getSynthesizedTypeSummaryInput()
func (c *Compiler) GetSynthesizedTypeSummaryInput() *CompilerInput {
	if c.syntheticTypeSummaryInput != nil {
		return c.syntheticTypeSummaryInput
	}

	input := NewCompilerInputWithIsExtern(c.syntheticTypeSummaryFile /* isExtern= */, true)
	root := jutil.CheckNotNull(input.GetAstRoot(c))
	c.putCompilerInput(input)
	c.syntheticTypeSummaryInput = input
	c.externsRoot.AddChildToFront(root)
	c.externs.AddAt(0, input)
	c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordTypeSummary()
	root.SetJSDocInfo(builder.Build())

	return input
}

// Java: Compiler.getSyntheticCodeInputId()
func (c *Compiler) GetSyntheticCodeInputId() *rhino.InputId {
	return compiler_SYNTHETIC_CODE_INPUT_ID
}

// Java: Compiler.initializeSyntheticCodeInput()
func (c *Compiler) InitializeSyntheticCodeInput() {
	jutil.CheckStateMsg(
		!c.inputsById.ContainsKey(compiler_SYNTHETIC_CODE_INPUT_ID),
		"Already initialized synthetic input")
	ast := NewCompilerInput(SourceFile_FromCode(compiler_SYNTHETIC_CODE_INPUT_ID.GetIdName(), ""))
	if c.inputsById.ContainsKey(ast.GetInputId()) {
		panic(jutil.NewIllegalStateExceptionMsg("Conflicting synthetic id name"))
	}
	c.jsRoot.AddChildToFront(jutil.CheckNotNull(ast.GetAstRoot(c)))

	var firstChunk *JSChunk
	if chunks := c.GetChunks(); len(chunks) > 0 {
		firstChunk = chunks[0]
	}
	if firstChunk.GetName() == JSChunk_STRONG_CHUNK_NAME {
		firstChunk.Add(ast)
	}
	ast.SetChunk(firstChunk)
	c.putCompilerInput(ast)

	c.commentsPerFile.Put(compiler_SYNTHETIC_CODE_INPUT_ID.GetIdName(), []*parser.Comment{})
	c.ReportChangeToChangeScope(ast.GetAstRoot(c))
}

// Java: Compiler.removeSyntheticCodeInput()
func (c *Compiler) RemoveSyntheticCodeInput() {
	c.removeSyntheticCodeInputMerging( /* mergeContentIntoFirstInput= */ false)
}

// Java: Compiler.mergeSyntheticCodeInput()
func (c *Compiler) MergeSyntheticCodeInput() {
	c.removeSyntheticCodeInputMerging( /* mergeContentIntoFirstInput= */ true)
}

// removeSyntheticCodeInputMerging deletes the synthesized code input and
// optionally moves all its subtree contents into the first non-synthetic
// input
//
// Java: Compiler.removeSyntheticCodeInput(boolean)
func (c *Compiler) removeSyntheticCodeInputMerging(mergeContentIntoFirstInput bool) {
	jutil.CheckStateMsg(
		c.inputsById.ContainsKey(compiler_SYNTHETIC_CODE_INPUT_ID),
		"Never initialized the synthetic input")
	input := c.inputsById.Get(compiler_SYNTHETIC_CODE_INPUT_ID)
	astRoot := input.GetAstRoot(c)
	jutil.CheckState(astRoot.IsFirstChildOf(c.jsRoot))
	jutil.CheckState(compiler_SYNTHETIC_CODE_INPUT_ID.Equals(input.GetInputId()))

	if mergeContentIntoFirstInput && astRoot.HasChildren() {
		// If we've inserted anything into the synthetic AST, move it into the top of the next script.
		next := astRoot.GetNext()
		jutil.CheckNotNullMsg(next, "Must provide at least one source")
		compiler_checkNotModule(
			next, "Cannot remove synthetic code input until modules are rewritten: %s", next)
		next.AddChildrenToFront(astRoot.RemoveChildren())
		c.ReportChangeToChangeScope(next)
	}

	astRoot.Detach()

	// bookkeeping to mark scopes and nodes as deleted
	c.ReportChangeToChangeScope(astRoot)
	astRoot.SetDeleted(true)
	NodeUtil_MarkFunctionsDeleted(astRoot, c)

	input.GetChunk().Remove(input)
	c.inputsById.Remove(input.GetInputId())
}

// Java: Compiler.getLastPassName()
func (c *Compiler) GetLastPassName() *string {
	return c.lastPassName
}

// Java: Compiler.setExternProperties(ImmutableSet<String>)
func (c *Compiler) SetExternProperties(externProperties *jutil.LinkedHashSet[string]) {
	c.externProperties = externProperties
}

// Java: Compiler.getExternProperties()
func (c *Compiler) GetExternProperties() *jutil.LinkedHashSet[string] {
	return c.externProperties
}

// Java: Compiler.getAccessorSummary()
func (c *Compiler) GetAccessorSummary() *AccessorSummary {
	return c.accessorSummary
}

// Java: Compiler.setAccessorSummary(AccessorSummary)
func (c *Compiler) SetAccessorSummary(summary *AccessorSummary) {
	c.accessorSummary = jutil.CheckNotNull(summary)
}

// loadResourceContents reads runtime library .js files into memory in AST
// format
//
// Java: Compiler.loadResourceContents(String, String)
func (c *Compiler) loadResourceContents(resourceName string, path string) *rhino.Node {
	if c.GetLifeCycleStage().HasColorAndSimplifiedJSDoc() {
		jutil.CheckNotNullMsg(
			c.runtimeLibraryTypedAsts,
			"Must call initRuntimeLibraryTypedAsts before calling ensureLibraryInjected during"+
				" optimizations")

		typedAstSupplier := jutil.CheckNotNullMsg1(
			c.runtimeLibraryTypedAsts.Get(path),
			strings.Join([]string{
				"Missing precompiled .typedast for '%s'. If this file is newly added, ",
				"you may need to regenerate runtime_libs.typedast.textproto",
				""}, ""),
			path)
		return typedAstSupplier()
	}
	jutil.CheckStateMsg(
		!c.HasTypeCheckingRun(),
		"runtime library injected after type checking but before optimization colors")
	originalCode := resources.ResourceLoader_LoadTextResource("com.google.javascript.jscomp.Compiler", "js/"+resourceName+".js")

	source := SourceFile_FromCode(path, originalCode)
	c.addFilesToSourceMap([]*SourceFile{source})
	return c.parseCodeHelper(source)
}

// Java: Compiler.getNodeForRuntimeCodeInsertion()
func (c *Compiler) getNodeForRuntimeCodeInsertion() *rhino.Node {
	// For stage 2 optimizing builds, insert runtime libraries as actual code in the AST.
	// For library builds that outputs a TypedAST, the runtime libraries are injected to a synthetic
	// @typeSummary (.i.js) node so that the types are available.
	if c.options.GetTypedAstOutputFile() == nil {
		return c.GetNodeForCodeInsertion(nil)
	}
	return c.GetSynthesizedTypeSummaryInput().GetAstRoot(c)
}

// Java: Compiler.addComments(String, List<Comment>)
func (c *Compiler) AddComments(filename string, comments []*parser.Comment) {
	if !c.GetOptions().PreservesDetailedSourceInfo() {
		panic(jutil.NewUnsupportedOperationExceptionMsg("addComments may only be called in IDE mode."))
	}
	c.commentsPerFile.Put(filename, comments)
}

// Java: Compiler.getComments(String)
func (c *Compiler) GetComments(filename string) []*parser.Comment {
	if !c.GetOptions().PreservesDetailedSourceInfo() {
		panic(jutil.NewUnsupportedOperationExceptionMsg("getComments may only be called in IDE mode."))
	}
	return c.commentsPerFile.Get(filename)
}

// Java: Compiler.getModuleLoader()
func (c *Compiler) GetModuleLoader() *ModuleLoader {
	return c.moduleLoader
}

// Java: Compiler.addFilesToSourceMap(Iterable<SourceFile>)
func (c *Compiler) addFilesToSourceMap(files []*SourceFile) {
	// synchronized annotation guards concurrent access to sourceMap during parsing.
	c.monitor.Lock()
	defer c.monitor.Unlock()
	if c.GetOptions().GetSourceMapIncludeSourcesContent() && c.GetSourceMap() != nil {
		for _, file := range files {
			jutil.Try(func() {
				c.GetSourceMap().AddSourceFile(file.GetName(), jutil.StringPtr(file.GetCode()))
			}, func(e jutil.IIOException) {
				panic(jutil.NewRuntimeExceptionWithMessageAndCause("Cannot read code of a source map's source file.", e))
			})
		}
	}
}

// Java: Compiler.initWebpackMap(ImmutableMap<String, String>)
func (c *Compiler) InitWebpackMap(inputPathByWebpackId *jutil.LinkedHashMap[string, string]) {
	c.inputPathByWebpackId = inputPathByWebpackId
}

// SetDebugMessage sets an additional human-readable message appended to
// debugging contexts, like crash reports & source-after-each-pass logging.
//
// This is used by NestedCompilerRunner to distinguish between the main
// compiler & nested compiler.
//
// Java: Compiler.setDebugMessage(String)
func (c *Compiler) SetDebugMessage(debugMessage *string) {
	c.compilerExecutor.SetDebugMessage(debugMessage)
	c.debugMessage = debugMessage
}

// Java: Compiler.createCompilerExecutor()
func (c *Compiler) CreateCompilerExecutor() *CompilerExecutor {
	return NewCompilerExecutor()
}

// Java: Compiler.getCompilerExecutor()
func (c *Compiler) GetCompilerExecutor() *CompilerExecutor {
	return c.compilerExecutor
}

// Serializable state of the compiler specific to multistage binary builds
//
// Only contains state that does not make sense in 'multilevel' binary builds
// (where library-level TypedASTs are the input). Such state belongs in the
// jscomp.TypedAst proto.

// Java: Compiler.toProto(Feature)
func compiler_toProto(feature parser.FeatureSet_Feature) FeatureProto {
	var result FeatureProto
	jutil.Try(func() {
		result = FeatureProto_ValueOfName("FEATURE_" + feature.Name())
	}, func(e jutil.IIllegalArgumentException) {
		result = FeatureProto_FEATURE_UNKNOWN
	})
	return result
}

// Java: Compiler.toProto(LifeCycleStage)
func compiler_toProtoOfLifeCycleStage(stage AbstractCompiler_LifeCycleStage) LifeCycleStageProto {
	switch stage {
	case AbstractCompiler_LifeCycleStage_RAW:
		return LifeCycleStageProto_LIFE_CYCLE_STAGE_RAW
	case AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC:
		return LifeCycleStageProto_LIFE_CYCLE_STAGE_COLORS_AND_SIMPLIFIED_JSDOC
	case AbstractCompiler_LifeCycleStage_NORMALIZED:
		return LifeCycleStageProto_LIFE_CYCLE_STAGE_NORMALIZED
	case AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED:
		return LifeCycleStageProto_LIFE_CYCLE_STAGE_NORMALIZED_OBFUSCATED
	}
	// Go-only: language-adapter — Java's exhaustive enum switch expression
	// throws MatchException for an unknown constant.
	panic(jutil.NewIllegalStateExceptionMsg("unexpected LifeCycleStage " + stage.String()))
}

// Java: Compiler.toVariableMapEntries(VariableMap)
func compiler_toVariableMapEntries(m *VariableMap) []*VariableMapEntryProto {
	entries := []*VariableMapEntryProto{}
	if m != nil {
		for key, value := range m.GetOriginalNameToNewNameMap().All() {
			entries = append(entries,
				VariableMapEntryProto_NewBuilder().
					SetOriginalName(key).
					SetNewName(value).
					Build())
		}
	}
	return entries
}

// Java: Compiler.fromProto(List<FeatureProto>)
func compiler_fromProto(protos []FeatureProto) *parser.FeatureSet {
	features := jutil.NewEnumSet[parser.FeatureSet_Feature]()
	for _, p := range protos {
		if p != FeatureProto_FEATURE_UNKNOWN {
			jutil.Try(func() {
				features.Add(parser.FeatureSet_Feature_ValueOf(p.String()[len("FEATURE_"):]))
			}, func(e jutil.IIllegalArgumentException) {
				// ignore
			})
		}
	}
	return parser.FeatureSet_BARE_MINIMUM.WithEnumSet(features)
}

// Java: Compiler.fromProto(LifeCycleStageProto)
func compiler_fromProtoOfLifeCycleStage(stage LifeCycleStageProto) AbstractCompiler_LifeCycleStage {
	switch stage {
	case LifeCycleStageProto_LIFE_CYCLE_STAGE_RAW:
		return AbstractCompiler_LifeCycleStage_RAW
	case LifeCycleStageProto_LIFE_CYCLE_STAGE_COLORS_AND_SIMPLIFIED_JSDOC:
		return AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC
	case LifeCycleStageProto_LIFE_CYCLE_STAGE_NORMALIZED:
		return AbstractCompiler_LifeCycleStage_NORMALIZED
	case LifeCycleStageProto_LIFE_CYCLE_STAGE_NORMALIZED_OBFUSCATED:
		return AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED
	default:
		return AbstractCompiler_LifeCycleStage_RAW
	}
}

// Java: Compiler.fromVariableMapEntries(List<VariableMapEntryProto>)
func compiler_fromVariableMapEntries(entries []*VariableMapEntryProto) *VariableMap {
	m := jutil.NewLinkedHashMap[string, string]()
	for _, entry := range entries {
		m.Put(entry.GetOriginalName(), entry.GetNewName())
	}
	return NewVariableMap(m)
}

// GetCompilerStateProtoBuilder returns a builder for the serializable state
// of the compiler.
//
// Java: Compiler.getCompilerStateProtoBuilder()
func (c *Compiler) GetCompilerStateProtoBuilder() *JSCompilerStateProto_Builder {
	builder := JSCompilerStateProto_NewBuilder()

	for feature := range c.allowableFeatures.GetFeatures().All() {
		builder.AddAllowableFeatures(compiler_toProto(feature))
	}
	builder.
		SetTypeCheckingHasRun(c.typeCheckingHasRun).
		SetHasRegExpGlobalReferences(c.hasRegExpGlobalReferences).
		SetLifeCycleStage(compiler_toProtoOfLifeCycleStage(c.GetLifeCycleStage())).
		SetMergedPrecompiledLibraries(c.options.GetMergedPrecompiledLibraries())

	chunkList := c.chunkGraph.GetIntegralChunkArrayForSerialization()
	chunkNameMap := jutil.NewLinkedHashMap[string, int]()
	for i := 0; i < len(chunkList); i++ {
		chunkNameMap.Put(chunkList[i].GetName(), i)
	}

	for _, chunk := range chunkList {
		chunkBuilder := ChunkProto_NewBuilder()
		chunkBuilder.SetName(chunk.GetName())
		for _, dep := range chunk.GetDependencies() {
			depIdx, ok := chunkNameMap.Lookup(dep.GetName())
			jutil.CheckStateMsg2(
				ok, "Missing dependency %s for chunk %s", dep.GetName(), chunk.GetName())
			chunkBuilder.AddDependencies(int32(depIdx))
		}
		for _, input := range chunk.GetInputs() {
			chunkBuilder.AddInputIds(input.GetInputId().GetIdName())
		}
		builder.AddChunks(chunkBuilder.Build())
	}

	builder.
		SetUniqueNameId(int32(c.uniqueNameId)).
		AddAllUniqueIdSupplier(c.uniqueIdSupplier.ToProto()).
		AddAllExportedNames(c.exportedNames.ToSlice())

	builder.SetHasCssNames(c.cssNames != nil)
	if c.cssNames != nil {
		builder.AddAllCssNames(c.cssNames.ToSlice())
	}

	if c.idGeneratorMap != nil {
		builder.SetIdGeneratorMap(*c.idGeneratorMap)
	}

	builder.
		SetTranspiledFiles(c.transpiledFiles).
		SetIdGeneratorCurrentId(int32(c.crossChunkIdGenerator.GetCurrentId())).
		SetRunJ2ClPasses(c.runJ2clPasses)

	for extern := range c.externs.All() {
		builder.AddExterns(extern.GetInputId().GetIdName())
	}

	builder.AddAllInjectedLibraries(c.runtimeJsLibManager.GetInjectedLibraries())
	lastInjectedLibrary := c.runtimeJsLibManager.GetLastInjectedLibrary()
	if lastInjectedLibrary != nil && c.jsRoot != nil && c.jsRoot.HasChildren() {
		builder.SetLastInjectedLibraryIndexInFirstScript(
			int32(c.jsRoot.GetFirstChild().GetIndexOfChild(lastInjectedLibrary)))
	} else {
		builder.SetLastInjectedLibraryIndexInFirstScript(-1)
	}

	if c.accessorSummary != nil {
		builder.SetAccessorSummary(c.accessorSummary.ToProto())
	}

	builder.SetHasStringMap(c.stringMap != nil)
	if c.stringMap != nil {
		builder.AddAllStringMap(compiler_toVariableMapEntries(c.stringMap))
	}

	builder.SetHasInstrumentationMapping(c.instrumentationMapping != nil)
	if c.instrumentationMapping != nil {
		builder.AddAllInstrumentationMapping(compiler_toVariableMapEntries(c.instrumentationMapping))
	}

	return builder
}

// RestoreFromState restores the portions of the compiler state that don't
// require access to the serialized AST.
//
// Java: Compiler.restoreFromState(JSCompilerStateProto)
func (c *Compiler) RestoreFromState(compilerState *JSCompilerStateProto) {
	c.allowableFeatures = compiler_fromProto(compilerState.GetAllowableFeaturesList())
	c.typeCheckingHasRun = compilerState.GetTypeCheckingHasRun()
	c.scriptNodeByFilename.Clear()
	c.hasRegExpGlobalReferences = compilerState.GetHasRegExpGlobalReferences()

	for _, library := range compilerState.GetInjectedLibrariesList() {
		c.GetRuntimeJsLibManager().RecordLibraryInjected(library)
	}

	stage := AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC
	if compilerState.GetLifeCycleStage() != LifeCycleStageProto_LIFE_CYCLE_STAGE_RAW {
		stage = compiler_fromProtoOfLifeCycleStage(compilerState.GetLifeCycleStage())
	}
	c.SetLifeCycleStage(stage)
	c.GetOptions().SetMergedPrecompiledLibraries(compilerState.GetMergedPrecompiledLibraries())

	c.uniqueNameId = int(compilerState.GetUniqueNameId())
	c.uniqueIdSupplier = UniqueIdSupplier_FromProto(compilerState.GetUniqueIdSupplierList())

	c.exportedNames.Clear()
	c.exportedNames.AddAll(compilerState.GetExportedNamesList()...)

	if compilerState.GetHasCssNames() {
		c.cssNames = jutil.NewLinkedHashSet(compilerState.GetCssNamesList()...)
	} else {
		c.cssNames = nil
	}

	c.variableMap = nil
	c.propertyMap = nil

	if compilerState.HasIdGeneratorMap() {
		idGeneratorMap := compilerState.GetIdGeneratorMap()
		c.idGeneratorMap = &idGeneratorMap
	} else {
		c.idGeneratorMap = nil
	}
	c.transpiledFiles = compilerState.GetTranspiledFiles()

	c.crossChunkIdGenerator = NewIdGenerator()
	c.crossChunkIdGenerator.SetCurrentId(int(compilerState.GetIdGeneratorCurrentId()))

	c.runJ2clPasses = compilerState.GetRunJ2ClPasses()

	if compilerState.HasAccessorSummary() {
		c.accessorSummary = AccessorSummary_FromProto(compilerState.GetAccessorSummary())
	} else {
		c.accessorSummary = nil
	}

	if compilerState.GetHasStringMap() {
		c.stringMap = compiler_fromVariableMapEntries(compilerState.GetStringMapList())
	} else {
		c.stringMap = nil
	}

	if compilerState.GetHasInstrumentationMapping() {
		c.instrumentationMapping =
			compiler_fromVariableMapEntries(compilerState.GetInstrumentationMappingList())
	} else {
		c.instrumentationMapping = nil
	}

	// We don't save changeStamp, because its value turned out to be non-deterministic
	// (Reasons for this are unknown.), and there's no benefit to saving it anyway.
	// We don't save the change stamps that are stored on the AST nodes,
	// so all the AST Nodes we read in effectively have a change stamp of 0,
	// and we can just start the compiler's counter over at 1.
	c.changeTracker.ResetChangeStamp()
}

// Java: Compiler.saveState(OutputStream)
func (c *Compiler) SaveState(outputStream io.Writer) {
	// Do not close the outputstream, caller is responsible for closing it.
	Compiler_RunInCompilerThread(c,
		func() any {
			c.evictPassArenaCaches()
			tracer := c.NewTracer("serializeCompilerState")
			nonClosingOut := NewCompiler_NonClosingOutputStream(outputStream)
			var compressionStream io.Writer
			wrapper := c.options.GetStateCompressionWrapper()
			if wrapper != nil {
				compressionStream = wrapper(nonClosingOut)
			} else {
				compressionStream = NewFastGzipOutputStream(nonClosingOut)
			}

			bufferedCompressionStream := compiler_newBufferedOutputStream(compressionStream, 64*1024)

			builder := c.GetCompilerStateProtoBuilder()
			compiler_writeDelimitedTo(builder.Build().ToByteArray(), bufferedCompressionStream)
			c.StopTracer(tracer, "serializeCompilerState")
			tracer = c.NewTracer("serializeTypedAst")
			SerializeTypedAstPass_CreateFromOutputStream(
				c,
				bufferedCompressionStream,
				SerializationOptions_Builder_().
					SetRuntimeLibraries([]string{}).
					SetIncludeDebugInfo(c.GetOptions().ShouldSerializeExtraDebugInfo()).
					Build()).
				Process(c.externsRoot, c.jsRoot)
			c.StopTracer(tracer, "serializeTypedAst")
			// Close the buffered stream to flush uncompressed bytes and close the compressor
			bufferedCompressionStream.Close()
			return nil
		})
}

// Java: Compiler.restoreState(InputStream)
func (c *Compiler) RestoreState(inputStream io.Reader) {
	c.initWarningsGuard(c.options.GetWarningsGuard())
	c.MaybeSetTracker()

	Compiler_RunInCompilerThread(c,
		func() any {
			tracer := c.NewTracer(PassNames_DESERIALIZE_COMPILER_STATE)
			Compiler_logger.Fine("Deserializing the CompilerState")
			nonClosingIn := NewCompiler_NonClosingInputStream(inputStream)
			var decompressStream io.Reader
			wrapper := c.options.GetStateDecompressionWrapper()
			if wrapper != nil {
				decompressStream = wrapper(nonClosingIn)
			} else {
				decompressStream = jutil.NewGZIPInputStreamWithSize(nonClosingIn, 64*1024)
			}
			bufferedDecompressStream := compiler_newBufferedInputStream(decompressStream, 64*1024)
			defer func() {
				bufferedDecompressStream.Close()
				Compiler_logger.Fine("Finished deserializing CompilerState")
				c.StopTracer(tracer, PassNames_DESERIALIZE_COMPILER_STATE)
			}()
			c.deserializeCompilerState(bufferedDecompressStream)
			return nil
		})

	if c.tracker != nil {
		c.tracker.UpdateAfterDeserialize(c.jsRoot)
	}
}

// deserializeCompilerState must be called from within a "compiler thread"
// with a larger stack
//
// Java: Compiler.deserializeCompilerState(InputStream)
func (c *Compiler) deserializeCompilerState(inputStream io.Reader) {
	// Do not close the input stream, caller is responsible for closing it.
	stateProto := compiler_parseDelimitedFrom(inputStream)

	jutil.CheckNotNullMsg(
		c.chunkGraph, "Did you forget to call .init or .initChunks before restoreState?")
	externFilesBuilder := jutil.NewLinkedHashMap[string, *SourceFile]()
	codeFilesBuilder := jutil.NewLinkedHashMap[string, *SourceFile]()
	allInputFiles := jutil.NewLinkedHashSet[*SourceFile]()
	compiler_putOrThrow := func(m *jutil.LinkedHashMap[string, *SourceFile], k string, v *SourceFile) {
		// ImmutableMap.Builder.buildOrThrow: duplicate keys are an error.
		if m.ContainsKey(k) {
			panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " + k))
		}
		m.Put(k, v)
	}
	for _, input := range c.chunkGraph.GetAllInputs() {
		allInputFiles.Add(input.GetSourceFile())
		compiler_putOrThrow(codeFilesBuilder, input.GetInputId().GetIdName(), input.GetSourceFile())
	}
	for extern := range c.externs.All() {
		allInputFiles.Add(extern.GetSourceFile())
		compiler_putOrThrow(externFilesBuilder, extern.GetInputId().GetIdName(), extern.GetSourceFile())
	}

	deserializedAst := TypedAstDeserializer_DeserializeFullAst(
		c,
		c.SYNTHETIC_EXTERNS_FILE,
		allInputFiles,
		inputStream,
		stateProto.GetTypeCheckingHasRun(),
		c.GetOptions().GetResolveSourceMapAnnotations(),
		c.GetOptions().GetParseInlineSourceMaps())

	c.RestoreFromState(stateProto)

	// Restore the chunk graph
	chunks := []*JSChunk{}
	for _, chunkProto := range stateProto.GetChunksList() {
		chunk := NewJSChunk(chunkProto.GetName())
		chunks = append(chunks, chunk)
	}
	for i := 0; i < len(chunks); i++ {
		chunk := chunks[i]
		chunkProto := stateProto.GetChunks(i)
		for _, depIdx := range chunkProto.GetDependenciesList() {
			chunk.AddDependency(chunks[depIdx])
		}
	}
	c.chunkGraph = NewJSChunkGraph(chunks)

	// Restore TypedAST and related fields
	c.externProperties = deserializedAst.GetExternProperties()
	c.externAndJsRoot = rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root())
	c.externsRoot = c.externAndJsRoot.GetFirstChild()
	c.jsRoot = c.externAndJsRoot.GetLastChild()
	c.inputsById.Clear()
	c.externs.Clear()

	c.colorRegistry = deserializedAst.GetColorRegistry()
	runtimeLibraries := deserializedAst.GetRuntimeLibraries()

	// Tells CompilerInput::getRoot to deserialize an AST rather than re-parsing the file
	c.typedAstFilesystem = deserializedAst.GetFilesystem()
	c.typedAstScriptDeserializers = deserializedAst.scriptDeserializers
	deserializedAst = nil

	// overwrite any existing CompilerInput instances. Reuse the SourceFiles created by
	// TypedAstDeserializer; otherwise lookups in this.typedAstFilesystem will fail.
	externFiles := externFilesBuilder
	codeFiles := codeFilesBuilder
	externFilesBuilder = nil
	codeFilesBuilder = nil
	for _, externName := range stateProto.GetExternsList() {
		if externName == c.SYNTHETIC_EXTERNS_FILE.GetName() {
			c.GetSynthesizedExternsInput()
			continue
		}
		externFile := externFiles.Get(externName)
		if externFile == nil {
			// Extern files may passed in with regular source files, and later 'hoisted' to externs
			// because they are annotated `@externs`.
			externFile = jutil.CheckNotNullMsg1(codeFiles.Get(externName), "Missing %s", externName)
		}
		input := NewCompilerInputWithIsExtern(externFile /* isExtern= */, true)
		script := input.GetAstRoot(c) // accesses this.typedAstFilesystem

		c.externsRoot.AddChildToBack(script)
		c.putCompilerInput(input)
		c.scriptNodeByFilename.Put(externFile.GetName(), script)
		c.externs.Add(input)
	}
	c.typedAstFilesystem.Remove(c.SYNTHETIC_EXTERNS_FILE)
	externFiles = nil

	// Restore the inputs for each chunk.
	type restoredChunkInput struct {
		chunkIdx int
		src      *SourceFile
		input    *CompilerInput
	}
	var restoredInputs []restoredChunkInput
	var prebuildInputs []*CompilerInput
	parallelRestore := c.options != nil && c.options.GetNumParallelThreads() > 1
	for i := 0; i < len(chunks); i++ {
		chunkProto := stateProto.GetChunks(i)
		for _, inputId := range chunkProto.GetInputIdsList() {
			src := codeFiles.Get(inputId)
			if src == nil {
				// The auto-generated empty fill files used to facilitate CCCM for
				// empty chunks may not have gotten serialized, but all the others
				// should have.
				jutil.CheckStateMsg1(AbstractCompiler_IsFillFileName(inputId), "Missing %s", inputId)
				continue
			}
			input := NewCompilerInput(src)
			if parallelRestore {
				restoredInputs = append(restoredInputs, restoredChunkInput{chunkIdx: i, src: src, input: input})
				prebuildInputs = append(prebuildInputs, input)
			} else {
				script := input.GetAstRoot(c) // accesses this.typedAstFilesystem

				c.jsRoot.AddChildToBack(script)
				c.scriptNodeByFilename.Put(src.GetName(), script)

				c.putCompilerInput(input) // overwrite the old input
				chunks[i].Add(input)
			}
		}
	}
	codeFiles = nil
	if parallelRestore && len(restoredInputs) > 0 {
		NewPrebuildAst(c, c.options.GetNumParallelThreads()).Prebuild(slices.Values(prebuildInputs))
		prebuildInputs = nil
		for _, item := range restoredInputs {
			script := item.input.GetAstRoot(c)

			c.jsRoot.AddChildToBack(script)
			c.scriptNodeByFilename.Put(item.src.GetName(), script)

			c.putCompilerInput(item.input) // overwrite the old input
			chunks[item.chunkIdx].Add(item.input)
		}
		restoredInputs = nil
	}

	for library := range runtimeLibraries.All() {
		c.runtimeJsLibManager.EnsureLibraryInjected(library, false)
	}

	c.typedAstFilesystem = nil // allow garbage collection
	c.typedAstScriptDeserializers = nil

	if stateProto.GetLastInjectedLibraryIndexInFirstScript() != -1 {
		lastInjectedLibrary := c.jsRoot.
			GetFirstChild().
			GetChildAtIndex(int(stateProto.GetLastInjectedLibraryIndexInFirstScript()))
		c.runtimeJsLibManager.SetLastInjectedLibrary(lastInjectedLibrary)
	}
}

// compiler_writeDelimitedTo is message.writeDelimitedTo(out): the message's
// size as a varint, then the message.
//
// Go-only: language-adapter — MessageLite.writeDelimitedTo(OutputStream) of
// the protoc-generated message.
func compiler_writeDelimitedTo(message []byte, out io.Writer) {
	buf := protoAppendVarint(make([]byte, 0, 10+len(message)), uint64(len(message)))
	buf = append(buf, message...)
	if _, err := out.Write(buf); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
}

// compiler_parseDelimitedFrom is JSCompilerStateProto.parseDelimitedFrom(in):
// it reads a varint size, then parses at most that many bytes of in (fewer if
// in ends first, as protobuf's LimitedInputStream does), and leaves in right
// after them. It returns nil (Java null) if in is at its end.
//
// Go-only: language-adapter — Parser.parseDelimitedFrom(InputStream) of the
// protoc-generated message.
func compiler_parseDelimitedFrom(in io.Reader) *JSCompilerStateProto {
	readByte := func() (byte, bool) {
		var b [1]byte
		for {
			n, err := in.Read(b[:])
			if n == 1 {
				return b[0], true
			}
			if err == io.EOF {
				return 0, false
			}
			if err != nil {
				panic(jutil.NewIOException(err.Error()))
			}
		}
	}
	firstByte, ok := readByte()
	if !ok {
		return nil
	}
	// CodedInputStream.readRawVarint32(int, InputStream)
	size := uint32(firstByte & 0x7f)
	if firstByte&0x80 != 0 {
		offset := 7
		for ; offset < 32; offset += 7 {
			b, ok := readByte()
			if !ok {
				panic(protoTruncatedMessage())
			}
			size |= uint32(b&0x7f) << offset
			if b&0x80 == 0 {
				break
			}
		}
		if offset >= 32 {
			// Keep reading up to 64 bits.
			for ; offset < 64; offset += 7 {
				b, ok := readByte()
				if !ok {
					panic(protoTruncatedMessage())
				}
				if b&0x80 == 0 {
					break
				}
			}
			if offset >= 64 {
				panic(protoMalformedVarint())
			}
		}
	}
	if int32(size) < 0 {
		panic(protoNegativeSize())
	}
	data, err := io.ReadAll(io.LimitReader(in, int64(size)))
	if err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	return JSCompilerStateProto_ParseFrom(data)
}

// compiler_bufferedOutputStream is new BufferedOutputStream(out, size).
//
// Go-only: language-adapter — java.io.BufferedOutputStream over bufio:
// close() flushes the buffer, then closes out (FilterOutputStream.close).
type compiler_bufferedOutputStream struct {
	w   *bufio.Writer
	out io.Writer
}

func compiler_newBufferedOutputStream(out io.Writer, size int) *compiler_bufferedOutputStream {
	return &compiler_bufferedOutputStream{w: bufio.NewWriterSize(out, size), out: out}
}

func (s *compiler_bufferedOutputStream) Write(b []byte) (int, error) {
	return s.w.Write(b)
}

func (s *compiler_bufferedOutputStream) Close() {
	if err := s.w.Flush(); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	if closer, ok := s.out.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
	}
}

// compiler_bufferedInputStream is new BufferedInputStream(in, size).
//
// Go-only: language-adapter — java.io.BufferedInputStream over bufio:
// close() closes in.
type compiler_bufferedInputStream struct {
	r  *bufio.Reader
	in io.Reader
}

func compiler_newBufferedInputStream(in io.Reader, size int) *compiler_bufferedInputStream {
	return &compiler_bufferedInputStream{r: bufio.NewReaderSize(in, size), in: in}
}

func (s *compiler_bufferedInputStream) Read(b []byte) (int, error) {
	return s.r.Read(b)
}

func (s *compiler_bufferedInputStream) Close() {
	if closer, ok := s.in.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
	}
}

// GetModuleTypeByName returns the module type for the provided namespace.
//
// Java: Compiler.getModuleTypeByName(String)
func (c *Compiler) GetModuleTypeByName(moduleName string) *CompilerInput_ModuleType {
	// Map.get: null when absent.
	if moduleType, ok := c.moduleTypesByName.Lookup(moduleName); ok {
		return &moduleType
	}
	return nil
}

// Java: Compiler.getModuleMetadataMap()
func (c *Compiler) GetModuleMetadataMap() *ModuleMetadataMap {
	return c.moduleMetadataMap
}

// Java: Compiler.setModuleMetadataMap(ModuleMetadataMap)
func (c *Compiler) SetModuleMetadataMap(moduleMetadataMap *ModuleMetadataMap) {
	c.moduleMetadataMap = moduleMetadataMap
}

// Java: Compiler.getModuleMap()
func (c *Compiler) GetModuleMap() *ModuleMap {
	return c.moduleMap
}

// Java: Compiler.setModuleMap(ModuleMap)
func (c *Compiler) SetModuleMap(moduleMap *ModuleMap) {
	c.moduleMap = moduleMap
}

// compiler_resolveSibling is a simplistic implementation of the
// java.nio.file.Path resolveSibling method that always uses "/" regardless
// of platform. This is required for consistency of path handling when
// using input source maps, which are expected to always use "/".
//
// If `toPath` is already absolute, just return it.
//
// Otherwise, strip the filename off of `fromPath` and logically replaces
// it with `toPath`. Then normalize the path by resolving and removing any
// "." or ".." elements.
//
// fromPath - must be a file (not directory); toPath - must be a file (not
// directory). Returns `toPath` adjusted so that it has the same parent
// directory as `fromPath`.
//
// Java: Compiler.resolveSibling(String, String)
func compiler_resolveSibling(fromPath string, toPath string) string {
	// If the destination is an absolute path, nothing to do.
	if strings.HasPrefix(toPath, "/") {
		return toPath
	}

	fromPathParts := jutil.StringSplitLiteral(fromPath, "/")
	toPathParts := jutil.StringSplitLiteral(toPath, "/")
	if len(fromPathParts) != 0 {
		fromPathParts = fromPathParts[:len(fromPathParts)-1]
	}

	for len(fromPathParts) != 0 && len(toPathParts) != 0 {
		if toPathParts[0] == "." {
			toPathParts = toPathParts[1:]
		} else if toPathParts[0] == ".." {
			toPathParts = toPathParts[1:]
			fromPathParts = fromPathParts[:len(fromPathParts)-1]
		} else {
			break
		}
	}

	fromPathParts = append(fromPathParts, toPathParts...)
	return strings.Join(fromPathParts, "/")
}

// Java: Compiler.resetAndIntitializeSourceMap()
func (c *Compiler) ResetAndIntitializeSourceMap() {
	if c.sourceMap == nil {
		return
	}
	c.sourceMap.Reset()
	if c.options.GetSourceMapIncludeSourcesContent() {
		if c.options.GetApplyInputSourceMaps() {
			// Add any input source map content files to the source map as potential sources
			for inputSourceMap := range c.inputSourceMaps.Values() {
				c.addSourceMapSourceFiles(inputSourceMap)
			}
		}

		// Add all the compilation sources to the source map as potential sources
		allChunks := c.GetChunks()
		if allChunks != nil {
			sourceFiles := []*SourceFile{}
			for _, chunk := range allChunks {
				for _, input := range chunk.GetInputs() {
					sourceFile := input.GetSourceFile()
					if sourceFile.IsStubSourceFileForAlreadyProvidedInput() {
						// do not add stub source files to the source map
						continue
					}
					sourceFiles = append(sourceFiles, input.GetSourceFile())
				}
			}
			c.addFilesToSourceMap(sourceFiles)
		}
	}
}

// Java: Compiler.checkNotModule(Node, String, Object...)
func compiler_checkNotModule(script *rhino.Node, msg string, args ...any) *rhino.Node {
	jutil.CheckArgumentMsg1(script.IsScript(), "%s", script)
	if !script.HasOneChild() {
		return script
	}
	jutil.CheckStateMsg(!script.GetFirstChild().IsModuleBody(), msg, args...)
	return script
}

// Java: Compiler.maybeLogPruningAnalysis(DependencyManagementResult)
func (c *Compiler) maybeLogPruningAnalysis(result *JSChunkGraph_DependencyManagementResult) {
	if c.tracker == nil || !c.options.GetTracerMode().DoPruningAnalysis() {
		return
	}
	analysis := PruningAnalysis_Create(result.Sorter(), result.EntryPointInputs().ToSlice())
	pruningResult := analysis.Analyze()
	c.tracker.SetPruningAnalysisSummary(
		"\nTransitive dependency count per entry point:\n" +
			jutil.ToJavaString(pruningResult.EntryPointDependencyCount()) +
			"\n\nTop bottleneck dependencies:\n" +
			jutil.ToJavaString(pruningResult.BottleneckBlame()))
}

// Compiler_NonClosingOutputStream is a FilterOutputStream whose close() only
// flushes.
//
// Java: Compiler.NonClosingOutputStream
type Compiler_NonClosingOutputStream struct {
	out io.Writer
}

// Java: Compiler.NonClosingOutputStream(OutputStream)
func NewCompiler_NonClosingOutputStream(out io.Writer) *Compiler_NonClosingOutputStream {
	return &Compiler_NonClosingOutputStream{out: out}
}

// Write is FilterOutputStream.write(byte[], int, int).
//
// Go-only: language-adapter — io.Writer for OutputStream.
func (s *Compiler_NonClosingOutputStream) Write(b []byte) (int, error) {
	return s.out.Write(b)
}

// Flush is FilterOutputStream.flush(): it flushes out.
//
// Go-only: language-adapter — an io.Writer flushes if it has a Flush method.
func (s *Compiler_NonClosingOutputStream) Flush() error {
	if f, ok := s.out.(interface{ Flush() error }); ok {
		return f.Flush()
	}
	return nil
}

// Java: Compiler.NonClosingOutputStream.close()
func (s *Compiler_NonClosingOutputStream) Close() error {
	return s.Flush()
}

// Compiler_NonClosingInputStream is a FilterInputStream whose close() does not
// close the underlying stream.
//
// Java: Compiler.NonClosingInputStream
type Compiler_NonClosingInputStream struct {
	in io.Reader
}

// Java: Compiler.NonClosingInputStream(InputStream)
func NewCompiler_NonClosingInputStream(in io.Reader) *Compiler_NonClosingInputStream {
	return &Compiler_NonClosingInputStream{in: in}
}

// Read is FilterInputStream.read(byte[], int, int).
//
// Go-only: language-adapter — io.Reader for InputStream.
func (s *Compiler_NonClosingInputStream) Read(b []byte) (int, error) {
	return s.in.Read(b)
}

// Java: Compiler.NonClosingInputStream.close()
func (s *Compiler_NonClosingInputStream) Close() error {
	// Do not close the underlying stream
	return nil
}

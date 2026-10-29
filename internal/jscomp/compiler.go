// Port of com/google/javascript/jscomp/Compiler.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/Result.java @ closure-compiler v20261021

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

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"iter"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"weak"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
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
	// files, with deterministic insertion-order tracking.
	inputSourceMaps     map[string]*SourceMapInput
	inputSourceMapOrder []string

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
	forwardDeclaredTypes compilerForwardDeclaredSet

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

// compilerForwardDeclaredSet implements jutil.Set[string] for
// Compiler.forwardDeclaredTypes without LinkedHashSet or interface dispatch
// on Compiler itself. When all is true (AssumeForwardDeclaredForMissingTypes),
// Contains reports true for every type name.
type compilerForwardDeclaredSet struct {
	all   bool
	seen  map[string]struct{}
	order []string
}

func (s *compilerForwardDeclaredSet) Add(e string) bool {
	if s.all {
		return false
	}
	if s.seen == nil {
		s.seen = make(map[string]struct{})
	}
	if _, ok := s.seen[e]; ok {
		return false
	}
	s.seen[e] = struct{}{}
	s.order = append(s.order, e)
	return true
}

func (s *compilerForwardDeclaredSet) Contains(e string) bool {
	if s.all {
		return true
	}
	_, ok := s.seen[e]
	return ok
}

func (s *compilerForwardDeclaredSet) Size() int {
	if s.all {
		return 0
	}
	return len(s.order)
}

func (s *compilerForwardDeclaredSet) IsEmpty() bool {
	if s.all {
		return true
	}
	return len(s.order) == 0
}

func (s *compilerForwardDeclaredSet) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		if s.all {
			return
		}
		for _, e := range s.order {
			if !yield(e) {
				return
			}
		}
	}
}

// initFields runs the Java field initializers, in declaration order.
func (c *Compiler) initFields() {
	c.externs = nil
	c.inputsById = make(map[string]*CompilerInput)
	c.scriptNodeByFilename = jutil.NewLinkedHashMap[string, rhino.Node]()
	c.sourceMapOriginalSources = jutil.NewLinkedHashMap[string, *SourceFile]()
	c.inputSourceMaps = make(map[string]*SourceMapInput)
	c.inputSourceMapOrder = nil
	c.commentsPerFile = jutil.NewLinkedHashMap[string, []parser.Comment]()
	c.uniqueIdSupplier = NewUniqueIdSupplier()
	c.uniqueNameId = 0
	c.hasRegExpGlobalReferences = true
	c.runJ2clPasses = false
	c.defaultCodingConvention = NewClosureCodingConvention()
	c.forwardDeclaredTypes = compilerForwardDeclaredSet{}
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
		c.typeRegistry = rhino.NewJSTypeRegistryWithForwardDeclaredTypes(c.oldErrorReporter, &c.forwardDeclaredTypes)
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
func (c *Compiler) GetTypeMismatches() []TypeMismatch {
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

// ComputeCFG is Control Flow Analysis.
//
// Java: Compiler.computeCFG()
func (c *Compiler) ComputeCFG() *ControlFlowGraph[rhino.Node] {
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

// Result holds compilation results.
type Result struct {
	// Java's public final fields; the ones other Go packages read are
	// exported (PORTING.md §3).
	Success  bool
	Errors   []*JSError
	Warnings []*JSError

	variableMap             *VariableMap
	propertyMap             *VariableMap
	namedAnonFunctionMap    *VariableMap
	stringMap               *VariableMap                 // nil is null
	instrumentationMappings *VariableMap                 // nil is null
	SourceMap               *SourceMap                   // nil is null
	cssNames                *jutil.LinkedHashSet[string] // nil is null
	ExternExport            *string                      // nil is null
	idGeneratorMap          *string                      // nil is null
	TranspiledFiles         bool
}

// Java: Result(ImmutableList<JSError>, ImmutableList<JSError>, VariableMap, VariableMap, VariableMap, VariableMap, VariableMap, SourceMap, String, Set<String>, String, boolean)
func NewResultWithAllMaps(
	errors []*JSError,
	warnings []*JSError,
	variableMap *VariableMap,
	propertyMap *VariableMap,
	namedAnonFunctionMap *VariableMap,
	stringMap *VariableMap,
	instrumentationMappings *VariableMap,
	sourceMap *SourceMap,
	externExport *string,
	cssNames *jutil.LinkedHashSet[string],
	idGeneratorMap *string,
	transpiledFiles bool) *Result {
	r := &Result{}
	r.Success = len(errors) == 0
	r.Errors = errors
	r.Warnings = warnings
	r.variableMap = variableMap
	r.propertyMap = propertyMap
	r.namedAnonFunctionMap = namedAnonFunctionMap
	r.stringMap = stringMap
	r.instrumentationMappings = instrumentationMappings
	r.SourceMap = sourceMap
	r.ExternExport = externExport
	r.cssNames = cssNames
	r.idGeneratorMap = idGeneratorMap
	r.TranspiledFiles = transpiledFiles
	return r
}

// NewResult is @VisibleForTesting and deprecated in Java.
//
// Java: Result(ImmutableList<JSError>, ImmutableList<JSError>, VariableMap, VariableMap, VariableMap, SourceMap, String)
func NewResult(
	errors []*JSError,
	warnings []*JSError,
	variableMap *VariableMap,
	propertyMap *VariableMap,
	namedAnonFunctionMap *VariableMap,
	sourceMap *SourceMap,
	externExport *string) *Result {
	return NewResultWithAllMaps(
		errors,
		warnings,
		variableMap,
		propertyMap,
		namedAnonFunctionMap,
		/* stringMap= */ nil,
		/* instrumentationMappings= */ nil,
		sourceMap,
		externExport,
		/* cssNames= */ nil,
		/* idGeneratorMap= */ nil,
		/* transpiledFiles= */ false)
}

// Result_PruneResultForPartialCompilation returns an almost empty result
// that is more appropriate for a partial compilation.
//
// For a partial compilation we only care about errors and warnings. It is
// unnecessary to examine all of the other results.
//
// result is the full `Result` object provided by the compiler.
//
// Java: Result.pruneResultForPartialCompilation(Result)
func Result_PruneResultForPartialCompilation(result *Result) *Result {
	emptyVariableMap := NewVariableMap(jutil.NewLinkedHashMap[string, string]())
	return NewResultWithAllMaps(
		/* errors= */ result.Errors,
		/* warnings= */ result.Warnings,
		/* variableMap= */ emptyVariableMap,
		/* propertyMap= */ emptyVariableMap,
		/* namedAnonFunctionMap= */ emptyVariableMap,
		/* stringMap= */ nil,
		/* instrumentationMappings= */ emptyVariableMap,
		/* sourceMap= */ nil,
		/* externExport= */ jutil.StringPtr(""),
		/* cssNames= */ nil,
		/* idGeneratorMap= */ nil,
		/* transpiledFiles= */ false)
}

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
	"iter"
	"strings"
	"sync"
	"sync/atomic"
	"weak"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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
// Go-only: language-adapter — Java's Compiler implements ErrorHandler,
// whose method collides in Go with AbstractCompiler's report(JSError):
// the Compiler seen as ErrorHandler is AsErrorHandler().
//
// compilerMap is a native Go map with thin method helpers matching the
// lookup/update operations used across Compiler files.
type compilerMap[K comparable, V any] map[K]V

func (m compilerMap[K, V]) Get(k K) V {
	return m[k]
}

func (m compilerMap[K, V]) Lookup(k K) (V, bool) {
	v, ok := m[k]
	return v, ok
}

func (m compilerMap[K, V]) Put(k K, v V) {
	m[k] = v
}

func (m compilerMap[K, V]) PutIfAbsent(k K, v V) {
	if _, ok := m[k]; !ok {
		m[k] = v
	}
}

func (m compilerMap[K, V]) Clear() {
	clear(m)
}

func (m compilerMap[K, V]) KeySlice() []K {
	if len(m) == 0 {
		return nil
	}
	keys := make([]K, 0, len(m))
	for k := range m { // order-independent: callers only update each independent map key in place
		keys = append(keys, k)
	}
	return keys
}

func (m compilerMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m { // order-independent: callers only clear independent cached state on each value
			if !yield(v) {
				return
			}
		}
	}
}

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
	moduleTypesByName compilerMap[string, CompilerInput_ModuleType]

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

	scriptNodeByFilename compilerMap[string, rhino.Node]

	// Go-only: AST compaction state (compiler_compact.go).
	astCompact astCompaction

	inputPathByWebpackId *jutil.LinkedHashMap[string, string]

	transpilationNamespace rhino.StaticScope

	runtimeJsLibManager *RuntimeJsLibManager

	// Original sources referenced by the source maps.
	sourceMapOriginalSources compilerMap[string, *SourceFile]

	// Configured SourceMapInputs, plus any source maps discovered in source
	// files, with deterministic insertion-order tracking.
	inputSourceMaps     map[string]*SourceMapInput
	inputSourceMapOrder []string

	// Map from filenames to lists of all the comments in each file.
	commentsPerFile map[string][]parser.Comment

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
	runtimeLibraryTypedAsts map[string]*ScriptNodeDeserializer

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

	externProperties *jutil.LinkedHashSet[string]
	accessorSummary  *AccessorSummary

	changeTracker *ChangeTracker

	// There is exactly one cached resolvedSourceMap, which you can think of
	// as a cache of size one.
	resolvedSourceMap compiler_ResolvedSourceMap

	typedAstFilesystem typedAstFilesystemMap

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

	// Go-only: language-adapter — the Compiler as an ErrorHandler (see the
	// type comment); one instance so that the same object is handed out every
	// time, as `this` is in Java.
	asErrorHandler    *compiler_errorHandler
	cachedAstAnalyzer *AstAnalyzer
	sync              *compilerSync
	*parallelRecorder
}

// compilerSync holds non-copyable synchronization and scratch state shared
// between the main Compiler and shallow-copied parallel worker Compilers.
type compilerSync struct {
	mu               sync.Mutex
	sourceMapMu      sync.RWMutex
	printScratchBusy atomic.Bool
	printScratch     codePrinterScratch
}

func (c *Compiler) ensureSync() *compilerSync {
	if c.sync == nil {
		c.sync = &compilerSync{}
	}
	return c.sync
}

func (c *Compiler) isParallelWorker() bool {
	return c.parallelRecorder != nil
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

var compiler_validityCheck = &PassFactory{
	name:                "validityCheck",
	runInFixedPointLoop: true,
	preconditionCheck:   passFactoryAlwaysValid,
	internalFactory:     func(compiler AbstractCompiler) CompilerPass { return NewValidityCheck(compiler) },
	// TODO(b/280906684): Cannot use a lambda nor a simple anonymous function (which might be
	// converted to lambdas) directly due to this inference problem. Needs the cast to
	// make jscompiler happy.
	condition: func(o *CompilerOptions) bool {
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected"))
	},
}

// Java: Compiler.SYNTHETIC_FILE_NAME_PREFIX
const compiler_SYNTHETIC_FILE_NAME_PREFIX = " [synthetic:"

// Java: Compiler.SYNTHETIC_CODE_INPUT_ID
var compiler_SYNTHETIC_CODE_INPUT_ID = rhino.NewInputId(compiler_SYNTHETIC_FILE_NAME_PREFIX + "input] ")

// initFields runs the Java field initializers, in declaration order.
func (c *Compiler) initFields() {
	c.externs = nil
	c.inputsById = make(map[string]*CompilerInput)
	c.scriptNodeByFilename = make(compilerMap[string, rhino.Node])
	c.sourceMapOriginalSources = make(compilerMap[string, *SourceFile])
	c.inputSourceMaps = make(map[string]*SourceMapInput)
	c.inputSourceMapOrder = nil
	c.commentsPerFile = make(map[string][]parser.Comment)
	c.uniqueIdSupplier = NewUniqueIdSupplier()
	c.uniqueNameId = 0
	c.hasRegExpGlobalReferences = true
	c.runJ2clPasses = false
	c.defaultCodingConvention = NewClosureCodingConvention()
	c.forwardDeclaredTypes = compilerForwardDeclaredSet{}
	c.typeCheckingHasRun = false
	c.asErrorHandler = &compiler_errorHandler{c: c}
	c.oldErrorReporter = RhinoErrorReporter_ForOldRhino(c.asErrorHandler)
	c.compilerExecutor = c.CreateCompilerExecutor()
	c.changeTracker = NewChangeTracker()
	c.resolvedSourceMap = compiler_ResolvedSourceMap{}
	c.validityCheck = compiler_validityCheck
	c.currentPassIndex = -1
	c.stage = AbstractCompiler_LifeCycleStage_RAW
	c.preferRegexParser = false
	c.SYNTHETIC_EXTERNS_FILE = SourceFile_FromCodeWithKind(compiler_SYNTHETIC_FILE_NAME_PREFIX+"externs] ", "", rhino.StaticSourceFile_SourceKind_EXTERN)
	c.syntheticTypeSummaryFile = SourceFile_FromCodeWithKind(compiler_SYNTHETIC_FILE_NAME_PREFIX+"typeSummary] ", "", rhino.StaticSourceFile_SourceKind_EXTERN)
	c.crossChunkIdGenerator = NewIdGenerator()
	c.exportedNames = jutil.NewLinkedHashSet[string]()
	c.defineNames = jutil.NewLinkedHashSet[string]()
	c.sync = &compilerSync{}
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
	c.moduleTypesByName = make(compilerMap[string, CompilerInput_ModuleType])
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

// AsErrorHandler is the Compiler as an ErrorHandler (Java's `this`).
func (c *Compiler) AsErrorHandler() ErrorHandler { return c.asErrorHandler }

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

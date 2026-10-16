package jscomp

// Go-only: test-helper — the AbstractCompiler test double the core tests and
// oracles use until com/google/javascript/jscomp/Compiler.java is ported.
//
// newTestCompiler(options) is the state of the Java `new Compiler()` followed
// by `initOptions(options)`: every AbstractCompiler method below is the body
// of the corresponding Compiler method, translated line by line. Where that
// body needs a class that is not ported yet (it lives above the core), the
// method panics with "test compiler: <method> needs <class>", so a test that
// reaches it fails visibly instead of silently diverging from Java.
//
// What initOptions does and what this double does:
//   - allowableFeatures, initExperimentalForceTranspileOptions, the error
//     manager (LoggerErrorManager over the options' ErrorFormat), the
//     NO_OP runtime-library mode for skipNonTranspilationPasses,
//     ModuleLoader.EMPTY, turning off type-based optimizations without type
//     checking, preferRegexParser, and the flag-combination checks: as in
//     Java.
//   - reconcileOptionsWithGuards needs DiagnosticGroups; the
//     validateRequiredInlinings step needs InlineFunctions; the
//     assumeForwardDeclaredForMissingTypes set needs a "contains everything"
//     Set. These are not done (the options keep their values), and
//     options.assumeForwardDeclaredForMissingTypes() makes newTestCompiler
//     panic.
//   - initWarningsGuard composes J2clSuppressWarningsGuard,
//     ClosureUnawareCodeWarningsGuard, SuppressDocWarningsGuard and the
//     options' guard. J2clSuppressWarningsGuard and SuppressDocWarningsGuard
//     are not ported; the double composes the other two, so report() does
//     not apply @suppress annotations or the J2CL suppressions. Tests that
//     depend on those report a difference, they do not hide one.

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// Java: Compiler.SYNTHETIC_FILE_NAME_PREFIX
const testCompiler_SYNTHETIC_FILE_NAME_PREFIX = " [synthetic:"

// Java: Compiler.SYNTHETIC_CODE_INPUT_ID
var testCompiler_SYNTHETIC_CODE_INPUT_ID = rhino.NewInputId(testCompiler_SYNTHETIC_FILE_NAME_PREFIX + "input] ")

// Java: Compiler.logger
var testCompiler_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp")

// testCompiler is the state of com.google.javascript.jscomp.Compiler that
// the AbstractCompiler methods read and write.
type testCompiler struct {
	AbstractCompilerBase

	options *CompilerOptions

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
	externsRoot     *rhino.Node
	jsRoot          *rhino.Node
	externAndJsRoot *rhino.Node

	allowableFeatures *parser.FeatureSet

	inputsById           *jutil.LinkedHashMap[string, *CompilerInput] // keyed by InputId.getIdName()
	scriptNodeByFilename *jutil.LinkedHashMap[string, *rhino.Node]

	transpilationNamespace rhino.StaticScope
	runtimeJsLibManager    *RuntimeJsLibManager

	sourceMapOriginalSources *jutil.LinkedHashMap[string, *SourceFile]
	inputSourceMaps          *jutil.LinkedHashMap[string, *SourceMapInput]
	commentsPerFile          *jutil.LinkedHashMap[string, []*parser.Comment]

	sourceMap     *SourceMap
	externExports *string

	uniqueIdSupplier *UniqueIdSupplier
	uniqueNameId     int

	hasRegExpGlobalReferences bool
	runJ2clPasses             bool

	typeRegistry            *rhino.JSTypeRegistry
	colorRegistry           *colors.ColorRegistry
	parserConfig            *parsing.Config
	externsParserConfig     *parsing.Config
	abstractInterpreter     ReverseAbstractInterpreter
	typeValidator           *TypeValidator
	hasPhaseOptimizer       bool // Compiler.phaseOptimizer != null
	runtimeLibraryTypedAsts *jutil.LinkedHashMap[string, func() *rhino.Node]

	forwardDeclaredTypes *jutil.LinkedHashSet[string]
	typeCheckingHasRun   bool

	// This error reporter gets the messages from the current Rhino parser or TypeRegistry.
	oldErrorReporter rhino.ErrorReporter

	lastPassName     *string
	externProperties *jutil.LinkedHashSet[string]
	accessorSummary  *AccessorSummary

	changeTracker *ChangeTracker

	currentPassIndex  int
	stage             AbstractCompiler_LifeCycleStage
	preferRegexParser bool

	typedScopeCreator ScopeCreator
	topScope          *TypedScope

	syntheticExternsFile      *SourceFile
	syntheticExternsInput     *CompilerInput
	syntheticTypeSummaryFile  *SourceFile
	syntheticTypeSummaryInput *CompilerInput
	indexProvidersByType      *jutil.LinkedHashMap[reflect.Type, IndexProvider[any]]

	crossChunkIdGenerator  *IdGenerator
	cssNames               *jutil.LinkedHashSet[string]
	variableMap            *VariableMap
	propertyMap            *VariableMap
	stringMap              *VariableMap
	instrumentationMapping *VariableMap
	idGeneratorMap         *string
	transpiledFiles        bool
	exportedNames          *jutil.LinkedHashSet[string]
	defineNames            *jutil.LinkedHashSet[string]

	moduleMetadataMap *ModuleMetadataMap
	moduleMap         *ModuleMap
}

var _ AbstractCompiler = (*testCompiler)(nil)

// testCompiler_errorHandler is the ErrorHandler the Compiler passes to
// RhinoErrorReporter.forOldRhino(this): Go cannot give testCompiler both
// report(JSError) and report(CheckLevel, JSError) under one name.
type testCompiler_errorHandler struct{ c *testCompiler }

func (h testCompiler_errorHandler) Report(level *CheckLevel, error *JSError) {
	h.c.ReportWithLevel(level, error)
}

// newTestCompiler is `new Compiler()` followed by `initOptions(options)`.
func newTestCompiler(options *CompilerOptions) AbstractCompiler {
	c := newTestCompilerUninitialized()
	c.initOptions(options)
	return c
}

// newTestCompilerUninitialized is `new Compiler()` (no options yet; like
// Compiler, it initializes default options when getErrorManager(), report()
// or the parse methods are called first).
func newTestCompilerUninitialized() *testCompiler {
	c := &testCompiler{
		externs:                   []*CompilerInput{},
		inputsById:                jutil.NewLinkedHashMap[string, *CompilerInput](),
		scriptNodeByFilename:      jutil.NewLinkedHashMap[string, *rhino.Node](),
		sourceMapOriginalSources:  jutil.NewLinkedHashMap[string, *SourceFile](),
		inputSourceMaps:           jutil.NewLinkedHashMap[string, *SourceMapInput](),
		commentsPerFile:           jutil.NewLinkedHashMap[string, []*parser.Comment](),
		uniqueIdSupplier:          NewUniqueIdSupplier(),
		hasRegExpGlobalReferences: true,
		forwardDeclaredTypes:      jutil.NewLinkedHashSet[string](),
		changeTracker:             NewChangeTracker(),
		currentPassIndex:          -1,
		stage:                     AbstractCompiler_LifeCycleStage_RAW,
		indexProvidersByType:      jutil.NewLinkedHashMap[reflect.Type, IndexProvider[any]](),
		crossChunkIdGenerator:     NewIdGenerator(),
		exportedNames:             jutil.NewLinkedHashSet[string](),
		defineNames:               jutil.NewLinkedHashSet[string](),
	}
	c.InitAbstractCompilerBase(c)
	c.oldErrorReporter = RhinoErrorReporter_ForOldRhino(testCompiler_errorHandler{c})
	c.syntheticExternsFile = SourceFile_FromCodeWithKind(
		testCompiler_SYNTHETIC_FILE_NAME_PREFIX+"externs] ", "", rhino.StaticSourceFile_SourceKind_EXTERN)
	c.syntheticTypeSummaryFile = SourceFile_FromCodeWithKind(
		testCompiler_SYNTHETIC_FILE_NAME_PREFIX+"typeSummary] ", "", rhino.StaticSourceFile_SourceKind_EXTERN)
	// Compiler(PrintStream)
	c.changeTracker.AddChangeHandler(c.changeTracker.GetRecentChange())
	c.moduleTypesByName = jutil.NewLinkedHashMap[string, CompilerInput_ModuleType]()
	return c
}

// SetErrorManager is Compiler.setErrorManager(ErrorManager).
func (c *testCompiler) SetErrorManager(errorManager ErrorManager) {
	jutil.CheckNotNullMsg(errorManager, "the error manager cannot be null")
	c.errorManager = NewThreadSafeDelegatingErrorManager(errorManager)
}

// Java: Compiler.createMessageFormatter()
func (c *testCompiler) createMessageFormatter() MessageFormatter {
	colorize := c.options.ShouldColorizeErrorOutput()
	return c.options.GetErrorFormat().ToFormatter(c, colorize)
}

// Java: Compiler.initExperimentalForceTranspileOptions(CompilerOptions)
func (c *testCompiler) initExperimentalForceTranspileOptions(options *CompilerOptions) {
	featureSet := options.GetOutputFeatureSet()
	if jutil.ObjectsEquals(featureSet, parser.FeatureSet_ES5) || jutil.ObjectsEquals(featureSet, parser.FeatureSet_ES3) {
		// do not change anything if language_out is configured to be ES5 or ES3.
		return
	}

	for _, experimentalForceTranspile := range options.GetExperimentalForceTranspiles() {
		switch experimentalForceTranspile {
		case CompilerOptions_ExperimentalForceTranspile_PERFORMANT_WITH_ASYNC_STACKS:
			options.SetOutputFeatureSet(
				parser.FeatureSet_ES5.WithFeatures(
					parser.FeatureSet_Feature_ASYNC_FUNCTIONS,
					parser.FeatureSet_Feature_ASYNC_GENERATORS,
					parser.FeatureSet_Feature_GENERATORS,
					parser.FeatureSet_Feature_LET_DECLARATIONS,
					parser.FeatureSet_Feature_CONST_DECLARATIONS,
					parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION))
		}
	}
}

// Java: Compiler.initOptions(CompilerOptions) (see the file comment for the
// steps that need unported classes).
func (c *testCompiler) initOptions(options *CompilerOptions) {
	c.options = options
	c.allowableFeatures = options.GetLanguageIn().ToFeatureSet()
	c.initExperimentalForceTranspileOptions(options)
	if c.errorManager == nil {
		// outStream == null
		c.SetErrorManager(NewLoggerErrorManagerWithFormatter(c.createMessageFormatter(), testCompiler_logger))
	}

	if options.GetSkipNonTranspilationPasses() {
		options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_NO_OP)
	}

	c.moduleLoader = ModuleLoader_EMPTY

	// reconcileOptionsWithGuards() and the validateRequiredInlinings switch:
	// not done (DiagnosticGroups, InlineFunctions).

	// Turn off type-based optimizations when type checking is off
	if !options.IsTypecheckingEnabled() {
		options.SetUseTypesForLocalOptimization(false)
		options.SetUseTypesForOptimization(false)
	}

	if options.AssumeForwardDeclaredForMissingTypes() {
		panic("test compiler: initOptions with assumeForwardDeclaredForMissingTypes needs " +
			"Compiler's report-everything forwardDeclaredTypes set")
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
				options.GetDependencyOptions().String() +
				""))
	}

	if options.GetMergedPrecompiledLibraries() && len(options.GetConformanceConfigs()) != 0 {
		panic("test compiler: initOptions with precompiled libraries and conformance configs " +
			"needs the conformance config file list")
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

// Java: Compiler.initWarningsGuard(WarningsGuard), without the guards that
// are not ported (J2clSuppressWarningsGuard, SuppressDocWarningsGuard).
func (c *testCompiler) initWarningsGuard(warningsGuard WarningsGuard) {
	c.warningsGuard = NewComposeWarningsGuard(NewClosureUnawareCodeWarningsGuard(), warningsGuard)
}

// Java: Compiler.initCompilerOptionsIfTesting()
func (c *testCompiler) initCompilerOptionsIfTesting() {
	if c.options == nil {
		// initialization for tests that don't initialize the compiler
		// by the normal mechanisms.
		c.initOptions(NewCompilerOptions())
		// allow tests to use unsupported features.
		c.options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	}
}

// Java: Compiler.initBasedOnOptions()
func (c *testCompiler) initBasedOnOptions() {
	for k, v := range c.options.GetInputSourceMaps().All() {
		c.inputSourceMaps.Put(k, v)
	}
	// Create the source map if necessary.
	if c.options.ShouldGatherSourceMapInfo() {
		c.sourceMap = c.options.GetSourceMapFormat().GetInstance()
		c.sourceMap.SetPrefixMappings(c.options.GetSourceMapLocationMappings())
		if c.options.GetApplyInputSourceMaps() {
			panic("test compiler: initBasedOnOptions with applyInputSourceMaps needs " +
				"Compiler as a SourceFileMapping")
		}
	}
}

// GetErrors is Compiler.getErrors() (Go-only: test-helper accessor).
func (c *testCompiler) GetErrors() []*JSError {
	if c.errorManager == nil {
		return []*JSError{}
	}
	return c.errorManager.GetErrors()
}

// GetWarnings is Compiler.getWarnings() (Go-only: test-helper accessor).
func (c *testCompiler) GetWarnings() []*JSError {
	if c.errorManager == nil {
		return []*JSError{}
	}
	return c.errorManager.GetWarnings()
}

// SetPreferRegexParser is Compiler.setPreferRegexParser(boolean).
func (c *testCompiler) SetPreferRegexParser(preferRegexParser bool) {
	c.preferRegexParser = preferRegexParser
}

// --- AbstractCompiler -------------------------------------------------------

func (c *testCompiler) BeforePass(passName string) {
	c.currentPassIndex++
}

func (c *testCompiler) AfterPass(passName string) {
	// maybePrintSourceAfterEachPass(passName)
	if !c.options.ShouldPrintSourceAfterEachPass() {
		return
	}
	panic("test compiler: afterPass with printSourceAfterEachPass needs CodePrinter")
}

func (c *testCompiler) GetTypedAstDeserializer(file *SourceFile) func() *rhino.Node {
	// typedAstFilesystem == null
	return nil
}

func (c *testCompiler) GetInput(id *rhino.InputId) *CompilerInput {
	// TODO(bradfordcsmith): Allowing null id is less ideal. Add checkNotNull(id) here and fix
	// call sites that break.
	if id == nil {
		return nil
	}
	return c.inputsById.Get(id.GetIdName())
}

// Java: Compiler.putCompilerInput(CompilerInput)
func (c *testCompiler) putCompilerInput(input *CompilerInput) *CompilerInput {
	input.SetCompiler(c)
	return c.inputsById.Put(input.GetInputId().GetIdName(), input)
}

func (c *testCompiler) GetSourceFileByName(sourceName string) *SourceFile {
	// Here we assume that the source name is the input name, this
	// is true of JavaScript parsed from source.
	input := c.inputsById.Get(sourceName)
	if input != nil {
		return input.GetSourceFile()
	}
	// Alternatively, the sourceName might have been reverse-mapped by
	// an input source-map, so let's look in our sourcemap original sources.
	return c.sourceMapOriginalSources.Get(sourceName)
}

// getSourceFileByNullableName is getSourceFileByName(String) with a null
// check, for the SourceExcerptProvider methods.
func (c *testCompiler) getSourceFileByNullableName(sourceName *string) *SourceFile {
	if sourceName == nil {
		return nil
	}
	return c.GetSourceFileByName(*sourceName)
}

func (c *testCompiler) GetScriptNode(filename string) *rhino.Node {
	return c.scriptNodeByFilename.Get(filename)
}

func (c *testCompiler) GetChunkGraph() *JSChunkGraph { return c.chunkGraph }

func (c *testCompiler) GetInputsInOrder() []*CompilerInput {
	return c.chunkGraph.GetAllInputs()
}

func (c *testCompiler) GetNumberOfInputs() int {
	// In some testing cases inputs will be null, but obviously there must be at least one input.
	if c.chunkGraph != nil {
		return c.chunkGraph.GetInputCount()
	}
	return 1
}

func (c *testCompiler) AddExportedNames(exportedNames *jutil.LinkedHashSet[string]) {
	c.exportedNames.AddAll(exportedNames.ToSlice()...)
}

func (c *testCompiler) GetExportedNames() *jutil.LinkedHashSet[string] { return c.exportedNames }

func (c *testCompiler) SetDefineNames(defineNames []string) {
	c.defineNames = jutil.NewLinkedHashSet(defineNames...)
}

func (c *testCompiler) GetDefineNames() *jutil.LinkedHashSet[string] { return c.defineNames }

func (c *testCompiler) SetVariableMap(variableMap *VariableMap) { c.variableMap = variableMap }

func (c *testCompiler) SetPropertyMap(propertyMap *VariableMap) { c.propertyMap = propertyMap }

func (c *testCompiler) SetStringMap(stringMap *VariableMap) { c.stringMap = stringMap }

func (c *testCompiler) SetCssNames(cssNames *jutil.LinkedHashSet[string]) { c.cssNames = cssNames }

func (c *testCompiler) SetInstrumentationMapping(instrumentationMapping *VariableMap) {
	c.instrumentationMapping = instrumentationMapping
}

func (c *testCompiler) SetIdGeneratorMap(serializedIdMappings string) {
	c.idGeneratorMap = &serializedIdMappings
}

func (c *testCompiler) GetTranspiledFiles() bool { return c.transpiledFiles }

func (c *testCompiler) GetCrossChunkIdGenerator() *IdGenerator { return c.crossChunkIdGenerator }

func (c *testCompiler) SetAnonymousFunctionNameMap(functionMap *VariableMap) {
	// NOTE: remove this method
}

func (c *testCompiler) SetTypeCheckingHasRun(hasRun bool) { c.typeCheckingHasRun = hasRun }

func (c *testCompiler) HasTypeCheckingRun() bool { return c.typeCheckingHasRun }

func (c *testCompiler) HasOptimizationColors() bool { return c.colorRegistry != nil }

func (c *testCompiler) IsTypeRegistryCleared() bool {
	return c.typeCheckingHasRun && c.typeRegistry == nil
}

func (c *testCompiler) GetTypeRegistry() *rhino.JSTypeRegistry {
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

func (c *testCompiler) ClearJSTypeRegistry() {
	c.typeRegistry = nil
	c.typeValidator = nil
	c.abstractInterpreter = nil
}

func (c *testCompiler) GetColorRegistry() *colors.ColorRegistry {
	return jutil.CheckNotNullMsg(c.colorRegistry, "Color registry has not been initialized yet")
}

func (c *testCompiler) SetColorRegistry(colorRegistry *colors.ColorRegistry) {
	jutil.CheckState(c.runtimeLibraryTypedAsts != nil)
	c.colorRegistry = colorRegistry
}

func (c *testCompiler) ForwardDeclareType(typeName string) { c.forwardDeclaredTypes.Add(typeName) }

func (c *testCompiler) GetTypedScopeCreator() ScopeCreator {
	if c.typedScopeCreator == nil {
		jutil.CheckStateMsg(
			!c.HasTypeCheckingRun(),
			"Attempted to re-initialize TypedScopeCreator after it had been cleared")
		panic("test compiler: getTypedScopeCreator needs TypedScopeCreator")
	}
	return c.typedScopeCreator
}

func (c *testCompiler) ClearTypedScopeCreator() { c.typedScopeCreator = nil }

func (c *testCompiler) GetTopScope() *TypedScope { return c.topScope }

func (c *testCompiler) SetTopScope(x *TypedScope) {
	jutil.CheckStateMsg1(x == nil || x.GetParent() == nil, "%s", x)
	c.topScope = x
}

func (c *testCompiler) GetTranspilationNamespace() rhino.StaticScope {
	if c.transpilationNamespace == nil {
		panic("test compiler: getTranspilationNamespace needs GlobalNamespace")
	}
	return c.transpilationNamespace
}

func (c *testCompiler) Report(error *JSError) {
	level := error.DefaultLevel()
	if c.warningsGuard != nil {
		newLevel := c.warningsGuard.Level(error)
		if newLevel != nil {
			level = *newLevel
		}
	}

	if level.IsOn() {
		c.initCompilerOptionsIfTesting()
		if c.GetOptions().GetErrorHandler() != nil {
			c.GetOptions().GetErrorHandler().Report(level.Ptr(), error)
		}
		c.errorManager.Report(level.Ptr(), error)
	}
}

// ReportWithLevel is Compiler.report(CheckLevel, JSError) (ErrorHandler).
func (c *testCompiler) ReportWithLevel(ignoredLevel *CheckLevel, error *JSError) {
	c.Report(error)
}

func (c *testCompiler) ThrowInternalError(message string, cause jutil.Throwable) {
	panic(jutil.NewRuntimeExceptionWithMessageAndCause(
		"INTERNAL COMPILER ERROR.\nPlease report this problem.\n\n"+message, cause))
}

func (c *testCompiler) GetCodingConvention() CodingConvention {
	convention := c.options.GetCodingConvention()
	if convention != nil {
		return convention
	}
	panic("test compiler: getCodingConvention (default) needs ClosureCodingConvention")
}

func (c *testCompiler) ReportChangeToEnclosingScope(n *rhino.Node) {
	c.changeTracker.ReportChangeToEnclosingScope(n)
}

func (c *testCompiler) ReportChangeToChangeScope(changeScopeRoot *rhino.Node) {
	c.changeTracker.ReportChangeToChangeScope(changeScopeRoot)
}

func (c *testCompiler) ReportFunctionDeleted(n *rhino.Node) {
	c.changeTracker.ReportFunctionDeleted(n)
}

func (c *testCompiler) ReportDisambiguatePropertiesSummary(summarySupplier func() string) {
	// tracker == null
}

func (c *testCompiler) ReportAmbiguatePropertiesSummary(summarySupplier func() string) {
	// tracker == null
}

func (c *testCompiler) GetNodeForCodeInsertion(chunk *JSChunk) *rhino.Node {
	if c.inputsById.ContainsKey(testCompiler_SYNTHETIC_CODE_INPUT_ID.GetIdName()) {
		return c.inputsById.Get(testCompiler_SYNTHETIC_CODE_INPUT_ID.GetIdName()).GetAstRoot(c)
	}
	if chunk == nil {
		allInputs := c.chunkGraph.GetAllInputs()
		if c.chunkGraph == nil || len(allInputs) == 0 {
			panic(jutil.NewIllegalStateExceptionMsg("No inputs"))
		}
		firstInput := allInputs[0]
		return testCompiler_checkNotModule(firstInput.GetAstRoot(c), "Cannot insert code into a module")
	}

	chunkInputs := chunk.GetInputs()
	if len(chunkInputs) != 0 {
		return testCompiler_checkNotModule(
			chunkInputs[0].GetAstRoot(c), "Cannot insert code into a module")
	}
	panic(jutil.NewIllegalStateExceptionMsg("Root chunk has no inputs"))
}

// Java: Compiler.checkNotModule(Node, String, Object...)
func testCompiler_checkNotModule(script *rhino.Node, msg string, args ...any) *rhino.Node {
	jutil.CheckArgumentMsg1(script.IsScript(), "%s", script)
	if !script.HasOneChild() {
		return script
	}
	jutil.CheckStateMsg(!script.GetFirstChild().IsModuleBody(), msg, args...)
	return script
}

func (c *testCompiler) GetTypeValidator() *TypeValidator {
	if c.typeValidator == nil {
		jutil.CheckStateMsg(
			!c.HasTypeCheckingRun(),
			"Attempted to re-initialize TypeValidator after it had been cleared")
		c.typeValidator = NewTypeValidator(c)
	}
	return c.typeValidator
}

func (c *testCompiler) GetTypeMismatches() []*TypeMismatch {
	if c.typeCheckingHasRun {
		return c.GetTypeValidator().GetMismatches()
	}
	panic(jutil.NewRuntimeException("Can't ask for type mismatches before type checking."))
}

func (c *testCompiler) SetExternExports(externExports string) {
	c.externExports = &externExports
}

func (c *testCompiler) ParseSyntheticCode(fileName string, js string) *rhino.Node {
	c.initCompilerOptionsIfTesting()
	source := SourceFile_FromCode(testCompiler_SYNTHETIC_FILE_NAME_PREFIX+fileName+"] ", js)
	c.addFilesToSourceMap([]*SourceFile{source})
	return c.parseCodeHelper(source)
}

func (c *testCompiler) ParseTestCode(js string) *rhino.Node {
	c.initCompilerOptionsIfTesting()
	c.initBasedOnOptions()
	// NOTE: The file name used here is important, because at least Es6RewriteClassTest expects
	// it to match the default name used by the CompilerTestCase class for the source file
	// created when its `test("source string", "expected string")` method is called.
	return c.parseCodeHelper(SourceFile_FromCode("testcode", js))
}

func (c *testCompiler) ParseTestCodeOfList(code []string) *rhino.Node {
	c.initCompilerOptionsIfTesting()
	c.initBasedOnOptions()
	srcs := make([]*SourceFile, len(code))
	for index, value := range code {
		srcs[index] = SourceFile_FromCode("testcode"+strconv.Itoa(index), value)
	}
	return c.parseCodeHelperOfList(srcs)
}

// Java: Compiler.parseCodeHelper(SourceFile)
func (c *testCompiler) parseCodeHelper(src *SourceFile) *rhino.Node {
	input := NewCompilerInput(src)
	c.putCompilerInput(input)
	root := input.GetAstRoot(c)
	c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)
	return jutil.CheckNotNull(root)
}

// Java: Compiler.parseCodeHelper(List<SourceFile>)
func (c *testCompiler) parseCodeHelperOfList(srcs []*SourceFile) *rhino.Node {
	root := rhino.IR_Root()
	for _, src := range srcs {
		root.AddChildToBack(c.parseCodeHelper(src))
	}
	return root
}

// Java: Compiler.addFilesToSourceMap(Iterable<SourceFile>)
func (c *testCompiler) addFilesToSourceMap(files []*SourceFile) {
	if c.GetOptions().GetSourceMapIncludeSourcesContent() && c.sourceMap != nil {
		for _, file := range files {
			c.sourceMap.AddSourceFile(file.GetName(), jutil.StringPtr(file.GetCode()))
		}
	}
}

func (c *testCompiler) ToSource() string {
	panic("test compiler: toSource needs CodePrinter")
}

func (c *testCompiler) ToSourceOfNode(n *rhino.Node) string {
	panic("test compiler: toSource(Node) needs CodePrinter")
}

func (c *testCompiler) PreferRegexParser() bool { return c.preferRegexParser }

func (c *testCompiler) GetDefaultErrorReporter() rhino.ErrorReporter { return c.oldErrorReporter }

func (c *testCompiler) GetReverseAbstractInterpreter() ReverseAbstractInterpreter {
	if c.abstractInterpreter == nil {
		var interpreter ChainableReverseAbstractInterpreter = NewSemanticReverseAbstractInterpreter(c.GetTypeRegistry())
		if c.options.GetClosurePass() {
			interpreter = NewClosureReverseAbstractInterpreter(c.GetTypeRegistry()).Append(interpreter).GetFirst()
		}
		c.abstractInterpreter = interpreter
	}
	return c.abstractInterpreter
}

func (c *testCompiler) GetLifeCycleStage() AbstractCompiler_LifeCycleStage { return c.stage }

func (c *testCompiler) InitRuntimeLibraryTypedAsts(colorPoolBuilder *ColorPool_Builder) {
	panic("test compiler: initRuntimeLibraryTypedAsts needs TypedAstDeserializer")
}

func (c *testCompiler) GetUniqueIdSupplier() *UniqueIdSupplier { return c.uniqueIdSupplier }

func (c *testCompiler) GetUniqueNameIdSupplier() func() string {
	return func() string {
		id := c.uniqueNameId
		c.uniqueNameId++
		return strconv.Itoa(id)
	}
}

func (c *testCompiler) HasHaltingErrors() bool {
	return !c.GetOptions().CanContinueAfterErrors() && c.errorManager.HasHaltingErrors()
}

func (c *testCompiler) GetChangeTracker() *ChangeTracker { return c.changeTracker }

func (c *testCompiler) AddIndexProvider(indexProvider IndexProvider[any]) {
	type_ := indexProvider.GetType()
	if c.indexProvidersByType.Put(type_, indexProvider) != nil {
		panic(jutil.NewIllegalStateExceptionMsg(
			"A provider is already registered for index of type " + type_.Name()))
	}
}

func (c *testCompiler) GetIndex(key reflect.Type) any {
	indexProvider := c.indexProvidersByType.Get(key)
	if indexProvider == nil {
		return nil
	}
	return indexProvider.Get()
}

func (c *testCompiler) GetJsRoot() *rhino.Node { return c.jsRoot }

func (c *testCompiler) HasScopeChanged(n *rhino.Node) bool {
	if !c.hasPhaseOptimizer {
		return true
	}
	panic("test compiler: hasScopeChanged needs PhaseOptimizer")
}

// Java: Compiler.getParserConfigLanguageMode(CompilerOptions.LanguageMode)
func (c *testCompiler) getParserConfigLanguageMode(languageMode CompilerOptions_LanguageMode) parsing.Config_LanguageMode {
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

func (c *testCompiler) GetParserConfig(context AbstractCompiler_ConfigContext) *parsing.Config {
	if c.parserConfig == nil || c.externsParserConfig == nil {
		if c.parserConfig == nil {
			configLanguageMode := c.getParserConfigLanguageMode(c.options.GetLanguageIn())
			strictMode := parsing.Config_StrictMode_SLOPPY
			if c.options.ExpectStrictModeInput() {
				strictMode = parsing.Config_StrictMode_STRICT
			}
			c.parserConfig = c.createConfig(configLanguageMode, strictMode)
			// Externs must always be parsed with at least ES5 language mode.
			if configLanguageMode == parsing.Config_LanguageMode_ECMASCRIPT3 {
				c.externsParserConfig = c.createConfig(parsing.Config_LanguageMode_ECMASCRIPT5, strictMode)
			} else {
				c.externsParserConfig = c.parserConfig
			}
		}
	}
	if context == AbstractCompiler_ConfigContext_EXTERNS {
		return c.externsParserConfig
	}
	return c.parserConfig
}

// Java: Compiler.createConfig(LanguageMode, StrictMode)
func (c *testCompiler) createConfig(mode parsing.Config_LanguageMode, strictMode parsing.Config_StrictMode) *parsing.Config {
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

func (c *testCompiler) GetErrorManager() ErrorManager {
	if c.options == nil {
		c.initOptions(NewCompilerOptions())
	}
	return c.errorManager
}

func (c *testCompiler) SetLifeCycleStage(stage AbstractCompiler_LifeCycleStage) { c.stage = stage }

func (c *testCompiler) AreNodesEqualForInlining(n1 *rhino.Node, n2 *rhino.Node) bool {
	if c.options.ShouldAmbiguateProperties() || c.options.ShouldDisambiguateProperties() {
		// The type based optimizations require that type information is preserved
		// during other optimizations.
		return n1.IsEquivalentToTyped(n2)
	}
	return n1.IsEquivalentTo(n2)
}

func (c *testCompiler) SetHasRegExpGlobalReferences(references bool) {
	c.hasRegExpGlobalReferences = references
}

func (c *testCompiler) HasRegExpGlobalReferences() bool { return c.hasRegExpGlobalReferences }

func (c *testCompiler) GetErrorLevel(error *JSError) *CheckLevel {
	jutil.CheckNotNull(c.options)
	return c.warningsGuard.Level(error)
}

func (c *testCompiler) Process(p CompilerPass) {
	p.Process(c.externsRoot, c.jsRoot)
}

func (c *testCompiler) GetRoot() *rhino.Node { return c.externAndJsRoot }

func (c *testCompiler) GetOptions() *CompilerOptions { return c.options }

func (c *testCompiler) GetAllowableFeatures() *parser.FeatureSet { return c.allowableFeatures }

func (c *testCompiler) SetAllowableFeatures(allowableFeatures *parser.FeatureSet) {
	c.allowableFeatures = allowableFeatures
}

func (c *testCompiler) MarkFeatureNotAllowed(feature parser.FeatureSet_Feature) {
	c.allowableFeatures = c.allowableFeatures.WithoutFeatures(feature)
}

func (c *testCompiler) MarkFeatureSetNotAllowed(featureSet *parser.FeatureSet) {
	c.allowableFeatures = c.allowableFeatures.Without(featureSet)
}

func (c *testCompiler) GetSynthesizedExternsInput() *CompilerInput {
	if c.syntheticExternsInput != nil {
		return c.syntheticExternsInput
	}

	input := NewCompilerInputWithIsExtern(c.syntheticExternsFile, true /* isExtern */)
	root := jutil.CheckNotNull(input.GetAstRoot(c))
	c.putCompilerInput(input)
	c.syntheticExternsInput = input
	c.externsRoot.AddChildToFront(root)
	c.externs = append([]*CompilerInput{input}, c.externs...)
	c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)

	return input
}

func (c *testCompiler) GetSynthesizedTypeSummaryInput() *CompilerInput {
	if c.syntheticTypeSummaryInput != nil {
		return c.syntheticTypeSummaryInput
	}

	input := NewCompilerInputWithIsExtern(c.syntheticTypeSummaryFile, true /* isExtern */)
	root := jutil.CheckNotNull(input.GetAstRoot(c))
	c.putCompilerInput(input)
	c.syntheticTypeSummaryInput = input
	c.externsRoot.AddChildToFront(root)
	c.externs = append([]*CompilerInput{input}, c.externs...)
	c.scriptNodeByFilename.Put(input.GetSourceFile().GetName(), root)
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordTypeSummary()
	root.SetJSDocInfo(builder.Build())

	return input
}

func (c *testCompiler) GetLastPassName() *string { return c.lastPassName }

func (c *testCompiler) GetRuntimeJsLibManager() *RuntimeJsLibManager { return c.runtimeJsLibManager }

func (c *testCompiler) SetExternProperties(externProperties *jutil.LinkedHashSet[string]) {
	c.externProperties = externProperties
}

func (c *testCompiler) GetExternProperties() *jutil.LinkedHashSet[string] { return c.externProperties }

func (c *testCompiler) AddInputSourceMap(sourceFileName string, inputSourceMap *SourceMapInput) {
	c.inputSourceMaps.Put(sourceFileName, inputSourceMap)
	if c.options.GetSourceMapIncludeSourcesContent() && c.sourceMap != nil {
		panic("test compiler: addInputSourceMap with sourcesContent needs Compiler.addSourceMapSourceFiles")
	}
}

func (c *testCompiler) GetBase64SourceMapContents(sourceFileName string) *string {
	sourceMapInput := c.inputSourceMaps.Get(sourceFileName)
	// This sourceMappingURL ends with a ".inline.map", and we want to get its base4-encoded
	// contents (i.e. "data:application/json;base64,<encoded_source_map>").
	if sourceMapInput != nil && strings.HasSuffix(sourceMapInput.GetOriginalPath(), ".inline.map") {
		unencoded := sourceMapInput.GetRawSourceMapContents() // This is the sourcemap.
		encoded := jutil.BaseEncoding_Base64().Encode(jutil.StringGetBytesUTF8(*unencoded))
		return &encoded
	}
	return nil
}

func (c *testCompiler) AddComments(filename string, comments []*parser.Comment) {
	if !c.GetOptions().PreservesDetailedSourceInfo() {
		panic(jutil.NewUnsupportedOperationExceptionMsg("addComments may only be called in IDE mode."))
	}
	c.commentsPerFile.Put(filename, comments)
}

func (c *testCompiler) GetAccessorSummary() *AccessorSummary { return c.accessorSummary }

func (c *testCompiler) SetAccessorSummary(summary *AccessorSummary) {
	c.accessorSummary = jutil.CheckNotNull(summary)
}

func (c *testCompiler) GetComments(filename string) []*parser.Comment {
	if !c.GetOptions().PreservesDetailedSourceInfo() {
		panic(jutil.NewUnsupportedOperationExceptionMsg("getComments may only be called in IDE mode."))
	}
	return c.commentsPerFile.Get(filename)
}

func (c *testCompiler) GetModuleLoader() *ModuleLoader { return c.moduleLoader }

func (c *testCompiler) GetModuleTypeByName(moduleName string) *CompilerInput_ModuleType {
	if t, ok := c.moduleTypesByName.Lookup(moduleName); ok {
		return t.Ptr()
	}
	return nil
}

func (c *testCompiler) SetRunJ2clPasses(runJ2clPasses bool) { c.runJ2clPasses = runJ2clPasses }

func (c *testCompiler) RunJ2clPasses() bool { return c.runJ2clPasses }

func (c *testCompiler) CreateAstFactory() *AstFactory {
	if c.HasTypeCheckingRun() {
		if c.HasOptimizationColors() {
			return AstFactory_CreateFactoryWithColors(c.stage, c.GetColorRegistry(), c.runtimeJsLibManager)
		}
		return AstFactory_CreateFactoryWithTypes(c.stage, c.GetTypeRegistry(), c.runtimeJsLibManager)
	}
	return AstFactory_CreateFactoryWithoutTypes(c.stage, c.runtimeJsLibManager)
}

func (c *testCompiler) CreateAstFactoryWithoutTypes() *AstFactory {
	return AstFactory_CreateFactoryWithoutTypes(c.stage, c.runtimeJsLibManager)
}

func (c *testCompiler) GetAstAnalyzer() *AstAnalyzer {
	options := c.GetOptions()
	analysisSettings := AstAnalyzer_Options_Builder_().
		SetUseTypesForLocalOptimization(options.ShouldUseTypesForLocalOptimization()).
		SetAssumeGettersArePure(options.GetAssumeGettersArePure()).
		SetHasRegexpGlobalReferences(c.HasRegExpGlobalReferences()).
		SetAssumeKnownBuiltinsArePure(options.GetAssumePropertiesAreStaticallyAnalyzable()).
		Build()
	return NewAstAnalyzer(analysisSettings, c.typeRegistry, c.GetAccessorSummary())
}

func (c *testCompiler) CreateDefaultExpressionDecomposer() *ExpressionDecomposer {
	return c.CreateExpressionDecomposer(
		c.GetUniqueNameIdSupplier(),
		jutil.NewLinkedHashSet[string](),
		Scope_CreateGlobalScope(rhino.NewNode(rhino.Token_SCRIPT)))
}

func (c *testCompiler) CreateExpressionDecomposer(
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
	return NewExpressionDecomposer(c, uniqueNameIdSupplier, knownConstantFunctions, scope, enabledWorkarounds)
}

func (c *testCompiler) GetModuleMetadataMap() *ModuleMetadataMap { return c.moduleMetadataMap }

func (c *testCompiler) SetModuleMetadataMap(moduleMetadataMap *ModuleMetadataMap) {
	c.moduleMetadataMap = moduleMetadataMap
}

func (c *testCompiler) GetModuleMap() *ModuleMap { return c.moduleMap }

func (c *testCompiler) SetModuleMap(moduleMap *ModuleMap) { c.moduleMap = moduleMap }

func (c *testCompiler) IsDebugLoggingEnabled() bool {
	return c.GetOptions().GetDebugLogDirectory() != nil
}

func (c *testCompiler) GetDebugLogFilterList() []string {
	if c.GetOptions().GetDebugLogFilter() == nil {
		return []string{}
	}
	// Splitter.on(',').omitEmptyStrings().splitToList(...)
	out := []string{}
	for _, part := range strings.Split(*c.GetOptions().GetDebugLogFilter(), ",") {
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func (c *testCompiler) CreateOrReopenLog(owner reflect.Type, firstNamePart string, restNameParts ...string) diagnostic.LogFile {
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

func (c *testCompiler) CreateOrReopenIndexedLog(owner reflect.Type, firstNamePart string, restNameParts ...string) diagnostic.LogFile {
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

func (c *testCompiler) GetSyntheticCodeInputId() *rhino.InputId {
	return testCompiler_SYNTHETIC_CODE_INPUT_ID
}

func (c *testCompiler) InitializeSyntheticCodeInput() {
	jutil.CheckStateMsg(
		!c.inputsById.ContainsKey(testCompiler_SYNTHETIC_CODE_INPUT_ID.GetIdName()),
		"Already initialized synthetic input")
	ast := NewCompilerInput(SourceFile_FromCode(testCompiler_SYNTHETIC_CODE_INPUT_ID.GetIdName(), ""))
	if c.inputsById.ContainsKey(ast.GetInputId().GetIdName()) {
		panic(jutil.NewIllegalStateExceptionMsg("Conflicting synthetic id name"))
	}
	c.jsRoot.AddChildToFront(jutil.CheckNotNull(ast.GetAstRoot(c)))

	var firstChunk *JSChunk
	if c.chunkGraph != nil && len(c.chunkGraph.GetAllChunks()) != 0 {
		firstChunk = c.chunkGraph.GetAllChunks()[0]
	}
	if firstChunk.GetName() == JSChunk_STRONG_CHUNK_NAME {
		firstChunk.Add(ast)
	}
	ast.SetChunk(firstChunk)
	c.putCompilerInput(ast)

	c.commentsPerFile.Put(testCompiler_SYNTHETIC_CODE_INPUT_ID.GetIdName(), []*parser.Comment{})
	c.ReportChangeToChangeScope(ast.GetAstRoot(c))
}

func (c *testCompiler) RemoveSyntheticCodeInput() {
	c.removeSyntheticCodeInput( /* mergeContentIntoFirstInput= */ false)
}

func (c *testCompiler) MergeSyntheticCodeInput() {
	c.removeSyntheticCodeInput( /* mergeContentIntoFirstInput= */ true)
}

// Java: Compiler.removeSyntheticCodeInput(boolean)
func (c *testCompiler) removeSyntheticCodeInput(mergeContentIntoFirstInput bool) {
	jutil.CheckStateMsg(
		c.inputsById.ContainsKey(testCompiler_SYNTHETIC_CODE_INPUT_ID.GetIdName()),
		"Never initialized the synthetic input")
	input := c.inputsById.Get(testCompiler_SYNTHETIC_CODE_INPUT_ID.GetIdName())
	astRoot := input.GetAstRoot(c)
	jutil.CheckState(astRoot.IsFirstChildOf(c.jsRoot))
	jutil.CheckState(testCompiler_SYNTHETIC_CODE_INPUT_ID.Equals(input.GetInputId()))

	if mergeContentIntoFirstInput && astRoot.HasChildren() {
		// If we've inserted anything into the synthetic AST, move it into the top of the next script.
		next := astRoot.GetNext()
		jutil.CheckNotNullMsg(next, "Must provide at least one source")
		testCompiler_checkNotModule(
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
	c.inputsById.Remove(input.GetInputId().GetIdName())
}

// --- SourceExcerptProvider -------------------------------------------------

func (c *testCompiler) GetSourceLine(sourceName *string, lineNumber int) *string {
	if lineNumber < 1 {
		return nil
	}
	input := c.getSourceFileByNullableName(sourceName)
	if input != nil && !input.IsStubSourceFileForAlreadyProvidedInput() {
		return input.GetLine(lineNumber)
	}
	return nil
}

func (c *testCompiler) GetSourceLines(sourceName *string, lineNumber int, length int) Region {
	if lineNumber < 1 {
		return nil
	}
	input := c.getSourceFileByNullableName(sourceName)
	if input != nil && !input.IsStubSourceFileForAlreadyProvidedInput() {
		return input.GetLines(lineNumber, length)
	}
	return nil
}

func (c *testCompiler) GetSourceRegion(sourceName *string, lineNumber int) Region {
	if lineNumber < 1 {
		return nil
	}
	input := c.getSourceFileByNullableName(sourceName)
	if input != nil && !input.IsStubSourceFileForAlreadyProvidedInput() {
		return input.GetRegion(lineNumber)
	}
	return nil
}

func (c *testCompiler) GetSourceMapping(sourceName *string, lineNumber int, columnNumber int) *sourcemap.Mapping_OriginalMapping {
	if sourceName == nil {
		return nil
	}
	sourceMap := c.inputSourceMaps.Get(*sourceName)
	if sourceMap == nil {
		return nil
	}
	panic("test compiler: getSourceMapping with input source maps needs Compiler.resolveSibling")
}

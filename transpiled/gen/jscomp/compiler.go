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
// Port of com/google/javascript/jscomp/Compiler.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_jscomptype "github.com/paullepoulpe/jsgo/transpiled/gen/jscomptype"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_resources "github.com/paullepoulpe/jsgo/transpiled/gen/resources"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	pkg_sourcemap "github.com/paullepoulpe/jsgo/transpiled/gen/sourcemap"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_compiler_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_compiler_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_compiler_3 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_compiler_4 j2gort.Ext_com_google_common_collect_Streams_FunctionWithIndex

func init() {
	lambdaSingleton_compiler_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
		o := j2gort.CastPtr[CompilerOptions](lp_1)
		_ = o
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected")))
	})
	lambdaSingleton_compiler_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewValidityCheck(j2gort.AsIface[AbstractCompiler](mrp0_)) })
	lambdaSingleton_compiler_3 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_73 uint32) bool {
		node := pkg_rhino.Node(rawNode_73)
		_ = node
		return true
	})
	lambdaSingleton_compiler_4 = j2gort.NewLambda_Ext_com_google_common_collect_Streams_FunctionWithIndex(func(lp_129 any, index int64) any {
		value := j2gort.CastString(lp_129)
		_ = value
		_ = index
		return j2gort.PtrToAny(SourceFile_FromCode_2(("testcode" + j2gort.StrValInt64(index)), value))
	})
}

// Java class: com.google.javascript.jscomp.Compiler
type Compiler struct {
	AbstractCompilerBase
	options                   *CompilerOptions                                    `gson:"options"`
	passes                    PassConfig                                          `gson:"passes"`
	externs                   j2gort.Ext_java_util_ArrayList                      `gson:"externs"`
	chunkGraph                *JSChunkGraph                                       `gson:"chunkGraph"`
	moduleLoader              *ModuleLoader                                       `gson:"moduleLoader"`
	moduleTypesByName         j2gort.OrderedMap[string, CompilerInput_ModuleType] `gson:"moduleTypesByName"`
	errorManager              ErrorManager                                        `gson:"errorManager"`
	warningsGuard             WarningsGuard                                       `gson:"warningsGuard"`
	lastJsSource              string                                              `gson:"lastJsSource"`
	allowableFeatures         *pkg_parser.FeatureSet                              `gson:"allowableFeatures"`
	inputsById                j2gort.Ext_java_util_Map                            `gson:"inputsById"`
	scriptNodeByFilename      j2gort.Ext_java_util_Map                            `gson:"scriptNodeByFilename"`
	hoistedExternsByChunk     j2gort.Ext_com_google_common_collect_SetMultimap    `gson:"hoistedExternsByChunk"`
	inputPathByWebpackId      j2gort.Ext_com_google_common_collect_ImmutableMap   `gson:"inputPathByWebpackId"`
	transpilationNamespace    pkg_rhino.StaticScope                               `gson:"transpilationNamespace"`
	runtimeJsLibManager       *RuntimeJsLibManager                                `gson:"runtimeJsLibManager"`
	sourceMapOriginalSources  *j2gort.Ext_java_util_concurrent_ConcurrentHashMap  `gson:"sourceMapOriginalSources"`
	inputSourceMaps           *j2gort.Ext_java_util_concurrent_ConcurrentHashMap  `gson:"inputSourceMaps"`
	commentsPerFile           j2gort.Ext_java_util_Map                            `gson:"commentsPerFile"`
	sourceMap                 *SourceMap                                          `gson:"sourceMap"`
	debugMessage              string                                              `gson:"debugMessage"`
	externExports             string                                              `gson:"externExports"`
	uniqueIdSupplier          *UniqueIdSupplier                                   `gson:"uniqueIdSupplier"`
	defaultCodingConvention   CodingConvention                                    `gson:"defaultCodingConvention"`
	typeRegistry              *pkg_rhino.JSTypeRegistry                           `gson:"typeRegistry"`
	colorRegistry             *pkg_colors.ColorRegistry                           `gson:"colorRegistry"`
	parserConfig              pkg_parsing.Config                                  `gson:"parserConfig"`
	externsParserConfig       pkg_parsing.Config                                  `gson:"externsParserConfig"`
	abstractInterpreter       pkg_jscomptype.ReverseAbstractInterpreter           `gson:"abstractInterpreter"`
	typeValidator             *TypeValidator                                      `gson:"typeValidator"`
	phaseOptimizer            *PhaseOptimizer                                     `gson:"phaseOptimizer"`
	tracker                   *PerformanceTracker                                 `gson:"tracker"`
	runtimeLibraryTypedAsts   j2gort.Ext_com_google_common_collect_ImmutableMap   `gson:"runtimeLibraryTypedAsts"`
	forwardDeclaredTypes      j2gort.Ext_java_util_Set                            `gson:"forwardDeclaredTypes"`
	oldErrorReporter          pkg_rhino.ErrorReporter                             `gson:"oldErrorReporter"`
	compilerExecutor          *CompilerExecutor                                   `gson:"compilerExecutor"`
	outStream                 *j2gort.Ext_java_io_PrintStream                     `gson:"outStream"`
	lastPassName              string                                              `gson:"lastPassName"`
	externProperties          j2gort.Ext_com_google_common_collect_ImmutableSet   `gson:"externProperties"`
	accessorSummary           *AccessorSummary                                    `gson:"accessorSummary"`
	cachedAstAnalyzer         *AstAnalyzer                                        `gson:"-"`
	cachedAstAnalyzerRegistry *pkg_rhino.JSTypeRegistry                           `gson:"-"`
	cachedAstAnalyzerSummary  *AccessorSummary                                    `gson:"-"`
	changeTracker             *ChangeTracker                                      `gson:"changeTracker"`
	resolvedSourceMap         *Compiler_ResolvedSourceMap                         `gson:"resolvedSourceMap"`
	typedAstFilesystem        j2gort.Ext_java_util_concurrent_ConcurrentMap       `gson:"typedAstFilesystem"`
	validityCheck             PassFactory                                         `gson:"validityCheck"`
	stage                     *AbstractCompiler_LifeCycleStage                    `gson:"stage"`
	typedScopeCreator         *TypedScopeCreator                                  `gson:"typedScopeCreator"`
	topScope                  *TypedScope                                         `gson:"topScope"`
	u_SYNTHETIC_EXTERNS_FILE  *SourceFile                                         `gson:"SYNTHETIC_EXTERNS_FILE"`
	syntheticExternsInput     *CompilerInput                                      `gson:"syntheticExternsInput"`
	syntheticTypeSummaryFile  *SourceFile                                         `gson:"syntheticTypeSummaryFile"`
	syntheticTypeSummaryInput *CompilerInput                                      `gson:"syntheticTypeSummaryInput"`
	indexProvidersByType      j2gort.Ext_java_util_Map                            `gson:"indexProvidersByType"`
	crossChunkIdGenerator     *IdGenerator                                        `gson:"crossChunkIdGenerator"`
	cssNames                  j2gort.Ext_java_util_Set                            `gson:"cssNames"`
	variableMap               *VariableMap                                        `gson:"variableMap"`
	propertyMap               *VariableMap                                        `gson:"propertyMap"`
	stringMap                 *VariableMap                                        `gson:"stringMap"`
	instrumentationMapping    *VariableMap                                        `gson:"instrumentationMapping"`
	idGeneratorMap            string                                              `gson:"idGeneratorMap"`
	exportedNames             *j2gort.Ext_java_util_LinkedHashSet                 `gson:"exportedNames"`
	defineNames               j2gort.Ext_com_google_common_collect_ImmutableSet   `gson:"defineNames"`
	moduleMetadataMap         *ModuleMetadataMap                                  `gson:"moduleMetadataMap"`
	moduleMap                 *ModuleMap                                          `gson:"moduleMap"`
	_syncMu                   sync.Mutex                                          `gson:"-"`
	externsRoot               pkg_rhino.Node                                      `gson:"externsRoot"`
	jsRoot                    pkg_rhino.Node                                      `gson:"jsRoot"`
	externAndJsRoot           pkg_rhino.Node                                      `gson:"externAndJsRoot"`
	uniqueNameId              int32                                               `gson:"uniqueNameId"`
	cachedAstAnalyzerMask     int32                                               `gson:"-"`
	currentPassIndex          int32                                               `gson:"currentPassIndex"`
	hasRegExpGlobalReferences bool                                                `gson:"hasRegExpGlobalReferences"`
	runJ2clPasses             bool                                                `gson:"runJ2clPasses"`
	typeCheckingHasRun        bool                                                `gson:"typeCheckingHasRun"`
	preferRegexParser         bool                                                `gson:"preferRegexParser"`
	transpiledFiles           bool                                                `gson:"transpiledFiles"`
}

func (this_ *Compiler) BindSelf_(s any) {
	_ = s
}
func (this_ *Compiler) IsJavaNil_() bool { return this_ == nil }
func (*Compiler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler", "Compiler"
}
func (this_ *Compiler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractCompilerBase.IdentityHash_()
}
func (*Compiler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Compiler", "com.google.javascript.jscomp.Compiler", func(v_ any) bool {
		return j2gort.InstanceOf[*Compiler](v_)
	})
}
func (*Compiler) Implements_jscomp_AbstractCompiler()      {}
func (*Compiler) Implements_jscomp_ErrorHandler()          {}
func (*Compiler) Implements_jscomp_SourceFileMapping()     {}
func (*Compiler) Implements_jscomp_SourceExcerptProvider() {}
func (*Compiler) Implements_jscomp_CompilerInputProvider() {}

var static_Compiler_CHUNK_DEPENDENCY_ERROR *DiagnosticType

func StaticGet_Compiler_CHUNK_DEPENDENCY_ERROR() *DiagnosticType {
	clinit_Compiler()
	return static_Compiler_CHUNK_DEPENDENCY_ERROR
}
func StaticSet_Compiler_CHUNK_DEPENDENCY_ERROR(v *DiagnosticType) {
	clinit_Compiler()
	static_Compiler_CHUNK_DEPENDENCY_ERROR = v
}

var static_Compiler_MISSING_ENTRY_ERROR *DiagnosticType

func StaticGet_Compiler_MISSING_ENTRY_ERROR() *DiagnosticType {
	clinit_Compiler()
	return static_Compiler_MISSING_ENTRY_ERROR
}
func StaticSet_Compiler_MISSING_ENTRY_ERROR(v *DiagnosticType) {
	clinit_Compiler()
	static_Compiler_MISSING_ENTRY_ERROR = v
}

var static_Compiler_MISSING_CHUNK_ERROR *DiagnosticType

func StaticGet_Compiler_MISSING_CHUNK_ERROR() *DiagnosticType {
	clinit_Compiler()
	return static_Compiler_MISSING_CHUNK_ERROR
}
func StaticSet_Compiler_MISSING_CHUNK_ERROR(v *DiagnosticType) {
	clinit_Compiler()
	static_Compiler_MISSING_CHUNK_ERROR = v
}

var static_Compiler_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_Compiler_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_Compiler()
	return static_Compiler_logger
}
func StaticSet_Compiler_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_Compiler()
	static_Compiler_logger = v
}

var static_Compiler_pathJoiner *j2gort.Ext_com_google_common_base_Joiner

func StaticGet_Compiler_PathJoiner() *j2gort.Ext_com_google_common_base_Joiner {
	clinit_Compiler()
	return static_Compiler_pathJoiner
}
func StaticSet_Compiler_PathJoiner(v *j2gort.Ext_com_google_common_base_Joiner) {
	clinit_Compiler()
	static_Compiler_pathJoiner = v
}

var static_Compiler_EMPTY_CHUNK_LIST_ERROR *DiagnosticType

func StaticGet_Compiler_EMPTY_CHUNK_LIST_ERROR() *DiagnosticType {
	clinit_Compiler()
	return static_Compiler_EMPTY_CHUNK_LIST_ERROR
}
func StaticSet_Compiler_EMPTY_CHUNK_LIST_ERROR(v *DiagnosticType) {
	clinit_Compiler()
	static_Compiler_EMPTY_CHUNK_LIST_ERROR = v
}

var static_Compiler_EMPTY_ROOT_CHUNK_ERROR *DiagnosticType

func StaticGet_Compiler_EMPTY_ROOT_CHUNK_ERROR() *DiagnosticType {
	clinit_Compiler()
	return static_Compiler_EMPTY_ROOT_CHUNK_ERROR
}
func StaticSet_Compiler_EMPTY_ROOT_CHUNK_ERROR(v *DiagnosticType) {
	clinit_Compiler()
	static_Compiler_EMPTY_ROOT_CHUNK_ERROR = v
}

var static_Compiler_DUPLICATE_INPUT *DiagnosticType

func StaticGet_Compiler_DUPLICATE_INPUT() *DiagnosticType {
	clinit_Compiler()
	return static_Compiler_DUPLICATE_INPUT
}
func StaticSet_Compiler_DUPLICATE_INPUT(v *DiagnosticType) {
	clinit_Compiler()
	static_Compiler_DUPLICATE_INPUT = v
}

var static_Compiler_DUPLICATE_EXTERN_INPUT *DiagnosticType

func StaticGet_Compiler_DUPLICATE_EXTERN_INPUT() *DiagnosticType {
	clinit_Compiler()
	return static_Compiler_DUPLICATE_EXTERN_INPUT
}
func StaticSet_Compiler_DUPLICATE_EXTERN_INPUT(v *DiagnosticType) {
	clinit_Compiler()
	static_Compiler_DUPLICATE_EXTERN_INPUT = v
}
func StaticGet_Compiler_SYNTHETIC_FILE_NAME_PREFIX() string { return " [synthetic:" }

var static_Compiler_SYNTHETIC_CODE_INPUT_ID *pkg_rhino.InputId

func StaticGet_Compiler_SYNTHETIC_CODE_INPUT_ID() *pkg_rhino.InputId {
	clinit_Compiler()
	return static_Compiler_SYNTHETIC_CODE_INPUT_ID
}
func StaticSet_Compiler_SYNTHETIC_CODE_INPUT_ID(v *pkg_rhino.InputId) {
	clinit_Compiler()
	static_Compiler_SYNTHETIC_CODE_INPUT_ID = v
}

var clinitDone_Compiler uint32
var clinitState_Compiler j2gort.ClinitState

func clinit_Compiler() {
	if atomic.LoadUint32(&clinitDone_Compiler) != 0 {
		return
	}
	clinitSlow_Compiler()
}

func clinitSlow_Compiler() {
	if !j2gort.EnterClinit(&clinitState_Compiler, &clinitDone_Compiler) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Compiler, &clinitDone_Compiler)
	static_Compiler_CHUNK_DEPENDENCY_ERROR = DiagnosticType_Error("JSC_CHUNK_DEPENDENCY_ERROR", "Bad dependency: {0} -> {1}. Chunks must be listed in dependency order.")
	static_Compiler_MISSING_ENTRY_ERROR = DiagnosticType_Error("JSC_MISSING_ENTRY_ERROR", "required entry point \"{0}\" never provided")
	static_Compiler_MISSING_CHUNK_ERROR = DiagnosticType_Error("JSC_MISSING_CHUNK_ERROR", "unknown chunk \"{0}\" specified in entry point spec")
	static_Compiler_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1("com.google.javascript.jscomp")
	static_Compiler_pathJoiner = j2gort.Ext_com_google_common_base_Joiner_On_1_string(Platform_GetFileSeperator())
	static_Compiler_EMPTY_CHUNK_LIST_ERROR = DiagnosticType_Error("JSC_EMPTY_CHUNK_LIST_ERROR", "At least one chunk must be provided")
	static_Compiler_EMPTY_ROOT_CHUNK_ERROR = DiagnosticType_Error("JSC_EMPTY_ROOT_CHUNK_ERROR", "Root chunk ''{0}'' must contain at least one source code input")
	static_Compiler_DUPLICATE_INPUT = DiagnosticType_Error("JSC_DUPLICATE_INPUT", "Duplicate input: {0}")
	static_Compiler_DUPLICATE_EXTERN_INPUT = DiagnosticType_Error("JSC_DUPLICATE_EXTERN_INPUT", "Duplicate extern input: {0}")
	static_Compiler_SYNTHETIC_CODE_INPUT_ID = pkg_rhino.NewInputId((" [synthetic:" + "input] "))
}

func (this_ *Compiler) initFields_Compiler() {
	_ = this_

	this_.externs = j2gort.NewExt_java_util_ArrayList()
	this_.lastJsSource = j2gort.NullString
	this_.inputsById = j2gort.NewExt_java_util_concurrent_ConcurrentHashMap()
	this_.scriptNodeByFilename = j2gort.NewExt_java_util_concurrent_ConcurrentHashMap()
	this_.hoistedExternsByChunk = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0())
	this_.sourceMapOriginalSources = j2gort.NewExt_java_util_concurrent_ConcurrentHashMap()
	this_.inputSourceMaps = j2gort.NewExt_java_util_concurrent_ConcurrentHashMap()
	this_.commentsPerFile = j2gort.NewExt_java_util_concurrent_ConcurrentHashMap()
	this_.debugMessage = j2gort.NullString
	this_.externExports = j2gort.NullString
	this_.uniqueIdSupplier = NewUniqueIdSupplier()

	this_.hasRegExpGlobalReferences = true

	this_.defaultCodingConvention = NewClosureCodingConvention()

	this_.forwardDeclaredTypes = j2gort.NewExt_java_util_LinkedHashSet()

	this_.oldErrorReporter = RhinoErrorReporter_ForOldRhino(this_)
	this_.compilerExecutor = this_.CreateCompilerExecutor()
	this_.lastPassName = j2gort.NullString

	this_.cachedAstAnalyzerMask = int32(-1)
	this_.changeTracker = NewChangeTracker()
	this_.resolvedSourceMap = NewCompiler_ResolvedSourceMap()
	this_.validityCheck = PassFactory_S_Builder().SetName("validityCheck").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_compiler_2).SetCondition((lambdaSingleton_compiler_1)).Build()
	this_.currentPassIndex = int32(-1)
	this_.stage = StaticGet_AbstractCompiler_LifeCycleStage_RAW()

	this_.u_SYNTHETIC_EXTERNS_FILE = SourceFile_FromCode_3((" [synthetic:" + "externs] "), "", pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN())
	this_.syntheticTypeSummaryFile = SourceFile_FromCode_3((" [synthetic:" + "typeSummary] "), "", pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN())
	this_.indexProvidersByType = j2gort.NewExt_java_util_LinkedHashMap()
	this_.crossChunkIdGenerator = NewIdGenerator()

	this_.idGeneratorMap = j2gort.NullString

	this_.exportedNames = j2gort.NewExt_java_util_LinkedHashSet()
	this_.defineNames = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
}

// Java constructor: Compiler()
func NewCompiler() *Compiler {
	clinit_Compiler()
	this_ := &Compiler{}
	InitCompiler(this_)
	return this_
}

func InitCompiler(this_ *Compiler) {
	InitCompiler_1_Ext_java_io_PrintStream(this_, nil)
}

// Java constructor: Compiler(java.io.@org.jspecify.annotations.Nullable PrintStream outStream)
func NewCompiler_1_Ext_java_io_PrintStream(outStream *j2gort.Ext_java_io_PrintStream) *Compiler {
	clinit_Compiler()
	this_ := &Compiler{}
	InitCompiler_1_Ext_java_io_PrintStream(this_, outStream)
	return this_
}

func InitCompiler_1_Ext_java_io_PrintStream(this_ *Compiler, outStream *j2gort.Ext_java_io_PrintStream) {
	_ = outStream
	InitAbstractCompiler(&this_.AbstractCompilerBase)
	this_.initFields_Compiler()
	this_.changeTracker.AddChangeHandler(ToIface_CodeChangeHandler(this_.changeTracker.GetRecentChange()))
	this_.outStream = outStream
	this_.moduleTypesByName = j2gort.NewOrderedMap[string, CompilerInput_ModuleType](

	// Java constructor: Compiler(com.google.javascript.jscomp.ErrorManager errorManager)
	)
}

func NewCompiler_1_ErrorManager(errorManager ErrorManager) *Compiler {
	clinit_Compiler()
	this_ := &Compiler{}
	InitCompiler_1_ErrorManager(this_, errorManager)
	return this_
}

func InitCompiler_1_ErrorManager(this_ *Compiler, errorManager ErrorManager) {
	_ = errorManager
	InitCompiler(this_)
	this_.SetErrorManager(errorManager)
}

// Java method: Compiler.setErrorManager(com.google.javascript.jscomp.ErrorManager errorManager)
func (this_ *Compiler) SetErrorManager(errorManager ErrorManager) {
	_ = errorManager
	j2gort.CheckNotNullIface(errorManager)
	this_.errorManager = j2gort.AsIface[ErrorManager](NewThreadSafeDelegatingErrorManager(errorManager))
}

// Java method: Compiler.createMessageFormatter()
func (this_ *Compiler) priv_Compiler_createMessageFormatter_1437() MessageFormatter {
	var colorize bool = this_.options.ShouldColorizeErrorOutput()
	_ = colorize
	return this_.options.GetErrorFormat().ToFormatter(this_, colorize)
}

// Java method: Compiler.initExperimentalOutputFeatureSetOptions(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *Compiler) priv_Compiler_initExperimentalOutputFeatureSetOptions_1438(options *CompilerOptions) {
	_ = options
	var featureSet *pkg_parser.FeatureSet = options.GetOutputFeatureSet()
	_ = featureSet
	if j2gort.ObjEqCustomPtr(featureSet, pkg_parser.StaticGet_FeatureSet_ES5()) || j2gort.ObjEqCustomPtr(featureSet, pkg_parser.StaticGet_FeatureSet_ES3()) {
		return
	}
	if options.GetExperimentalOutputFeatureSet().IsPresent() {
		switch (j2gort.Cast[CompilerOptions_ExperimentalOutputFeatureSet](options.GetExperimentalOutputFeatureSet().Get())).Ordinal() {
		case int32(0) /* ExperimentalOutputFeatureSet.ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS */ :
			options.SetOutputFeatureSet(pkg_parser.StaticGet_FeatureSet_ES5().With_1_FeatureSet_Feature_7181([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS(), pkg_parser.StaticGet_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS()}))
		case int32(1) /* ExperimentalOutputFeatureSet.BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD */ :
			options.SetOutputFeatureSet(pkg_parser.StaticGet_FeatureSet_BROWSER_2019().Without_2(pkg_parser.StaticGet_FeatureSet_Feature_CLASSES(), []*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER(), pkg_parser.StaticGet_FeatureSet_Feature_NEW_TARGET(), pkg_parser.StaticGet_FeatureSet_Feature_SUPER(), pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS()}))
		}
	}
}

// Java method: Compiler.initOptions(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *Compiler) InitOptions(options *CompilerOptions) {
	_ = options
	this_.options = options
	this_.allowableFeatures = options.GetLanguageIn().ToFeatureSet()
	this_.priv_Compiler_initExperimentalOutputFeatureSetOptions_1438(options)
	if this_.errorManager == nil || this_.errorManager.IsJavaNil_() {
		if this_.outStream == nil {
			this_.SetErrorManager(NewLoggerErrorManager_2(this_.priv_Compiler_createMessageFormatter_1437(), StaticGet_Compiler_Logger()))
		} else {
			var builder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
			_ = builder
			j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](builder.Add_1_any(NewPrintStreamErrorReportGenerator(this_.priv_Compiler_createMessageFormatter_1437(), this_.outStream, options.GetSummaryDetailLevel()))).AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](options.GetExtraReportGenerators()))
			this_.SetErrorManager(j2gort.AsIface[ErrorManager](NewSortingErrorManager(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](builder.Build()))))
		}
	}
	if options.GetSkipNonTranspilationPasses() {
		options.SetRuntimeLibraryMode(StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_NO_OP())
	}
	this_.moduleLoader = StaticGet_ModuleLoader_EMPTY()
	this_.ReconcileOptionsWithGuards()
	switch (options.GetShouldValidateRequiredInlinings()).Ordinal() {
	case int32(2) /* Tri.TRUE */ :
		var sufficientOptimizations *InlineFunctions_ShouldRequireInlining = InlineFunctions_CheckOptimizationLevelSupportsRequireInlining(options)
		_ = sufficientOptimizations
		if !(sufficientOptimizations.ShouldValidate()) {
			j2gort.CheckStateFmt(false, "shouldValidateRequiredInlinings cannot be set unless CompilerOptions enables sufficient optimization, but found: %s", j2gort.BoxString(sufficientOptimizations.OptMessage()))
		}
	case int32(0) /* Tri.FALSE */ :
	case int32(1) /* Tri.UNKNOWN */ :
		options.SetValidateRequiredInlinings(InlineFunctions_CheckOptimizationLevelSupportsRequireInlining(options).ShouldValidate())
	}
	if !(options.IsTypecheckingEnabled()) {
		options.SetUseTypesForLocalOptimization(false)
		options.SetUseTypesForOptimization(false)
	}
	if options.AssumeForwardDeclaredForMissingTypes() {
		this_.forwardDeclaredTypes = NewCompiler__Anon_12(this_)
	}
	this_.priv_Compiler_initWarningsGuard_1443(options.GetWarningsGuard())
	if this_.IsDebugLoggingEnabled() {
		options.SetPrintConfig(true)
	}
	this_.preferRegexParser = (this_.preferRegexParser || options.GetDependencyOptions().ShouldPrune())
	if options.GetMergedPrecompiledLibraries() && options.GetDependencyOptions().NeedsManagement() {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(("Using precompiled libraries (i.e. TypedAST) is incompatible with flags that automatically order/prune dependencies: " + j2gort.StrVal(options.GetDependencyOptions()) + "")))))
	}
	if options.GetMergedPrecompiledLibraries() && len(j2gort.CheckNotNullPtr(options).Slice_GetConformanceConfigs()) != 0 {
		var conformanceConfigFiles j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
		_ = conformanceConfigFiles
		{

			for _, conformanceConfig := range j2gort.CheckNotNullPtr(options).
				Slice_GetConformanceConfigs() {

				_ = conformanceConfig
				{
					cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(conformanceConfig.GetRequirementList()))
					for cur_4.HasNext() {
						rawElem_5 := cur_4.Next()
						requirement := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_Requirement](rawElem_5)
						_ = requirement
						conformanceConfigFiles.AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](requirement.GetConfigFileList()))
					}
				}
			}
		}
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "use_precompiled_libraries mode does not support passing JS conformance configs: this mode skips all check passes including conformance\nUnsupported conformance configs: "+j2gort.StrVal(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](conformanceConfigFiles.Build())))))
	}
	if options.ShouldRunDeadAssignmentElimination() && !(options.ShouldRemoveUnusedVariables() || options.ShouldRemoveUnusedLocalVariables()) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Invalid flag combination: enabling dead assignment elimination also requires unused variable removal. Either disable dead assignment elimination or enable unused variable removal.")))
	}
	if options.ShouldRunDeadPropertyAssignmentElimination() && options.IsPolymerPassEnabled() {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Invalid flag combination: dead assignment elimination is incompatible with setPolymerPass(true).")))
	}
}

// Java method: Compiler.printConfig()
func (this_ *Compiler) PrintConfig() {
	if this_.IsDebugLoggingEnabled() {
		this_.priv_Compiler_logToFile_1442("externs.log", func() j2gort.Ext_java_util_function_Supplier {
			bound_6 := this_.externs
			return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.BoxString(bound_6.String()) })
		}())
		this_.priv_Compiler_logToFile_1442("inputs.json", func() j2gort.Ext_java_util_function_Supplier {
			bound_7 := j2gort.NewExt_com_google_gson_GsonBuilder().SetPrettyPrinting().Create().ToJson_1_Ext_com_google_gson_JsonElement(j2gort.AsIface[j2gort.Ext_com_google_gson_JsonElement](this_.chunkGraph.ToJson()))
			return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.CheckNotNullString(bound_7) })
		}())
		this_.priv_Compiler_logToFile_1442("options.log", j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return j2gort.BoxString(this_.options.String())
		}))
		this_.priv_Compiler_logToFile_1442("warningsGuard.log", j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return javalang.ToJavaString(j2gort.CheckNotNull(this_.warningsGuard))
		}))
	} else {
		var err *j2gort.Ext_java_io_PrintStream = j2gort.Ext_java_lang_System_F_Err
		_ = err
		err.Println_1_string("==== Externs ====")
		err.Println_1_any(this_.externs)
		err.Println_1_string("==== Inputs ====")
		err.Println_1_string(j2gort.NewExt_com_google_gson_GsonBuilder().SetPrettyPrinting().Create().ToJson_1_Ext_com_google_gson_JsonElement(j2gort.AsIface[j2gort.Ext_com_google_gson_JsonElement](this_.chunkGraph.ToJson())))
		err.Println_1_string("==== CompilerOptions ====")
		err.Println_1_any(j2gort.PtrToAny(this_.options))
		err.Println_1_string("==== WarningsGuard ====")
		err.Println_1_any(this_.warningsGuard)
	}
}

// Java method: Compiler.maybeLogPrunedInputs()
func (this_ *Compiler) priv_Compiler_maybeLogPrunedInputs_1441() {
	if this_.IsDebugLoggingEnabled() {
		this_.priv_Compiler_logToFile_1442("inputs_pruned.json", func() j2gort.Ext_java_util_function_Supplier {
			bound_8 := j2gort.NewExt_com_google_gson_GsonBuilder().SetPrettyPrinting().Create().ToJson_1_Ext_com_google_gson_JsonElement(j2gort.AsIface[j2gort.Ext_com_google_gson_JsonElement](this_.chunkGraph.ToJson()))
			return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.CheckNotNullString(bound_8) })
		}())
	}
}

// Java method: Compiler.logToFile(java.lang.String logFileName, java.util.function.Supplier<java.lang.String> logStringSupplier)
func (this_ *Compiler) priv_Compiler_logToFile_1442(logFileName string, logStringSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = logFileName
	_ = logStringSupplier
	func() {
		var log pkg_diagnostic.LogFile = this_.CreateOrReopenLog(j2gort.GetObjectClass(this_), logFileName, []string{})
		_ = log
		defer j2gort.CloseResource(log)
		if log.IsLogging() {
			log.Log_1_Ext_java_util_function_Supplier(logStringSupplier)
		}
	}()
}

// Java method: Compiler.initWarningsGuard(com.google.javascript.jscomp.WarningsGuard warningsGuard)
func (this_ *Compiler) priv_Compiler_initWarningsGuard_1443(warningsGuard WarningsGuard) {
	_ = warningsGuard
	var guards *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = guards
	guards.AddSelf_(NewJ2clSuppressWarningsGuard()).AddSelf_(NewClosureUnawareCodeWarningsGuard()).AddSelf_(NewSuppressDocWarningsGuard(this_, j2gort.AsIface[j2gort.Ext_java_util_Map](DiagnosticGroups_GetRegisteredGroups()))).AddVoid_(warningsGuard)
	this_.warningsGuard = NewComposeWarningsGuard_1_Ext_java_util_List(j2gort.AsIface[j2gort.Ext_java_util_List](guards.BuildImmList_()))
}

// Java method: Compiler.reconcileOptionsWithGuards()
func (this_ *Compiler) ReconcileOptionsWithGuards() {
	if this_.options.Enables(StaticGet_DiagnosticGroups_CHECK_TYPES()) {
		this_.options.SetCheckTypes(true)
	} else if this_.options.Disables(StaticGet_DiagnosticGroups_CHECK_TYPES()) {
		this_.options.SetCheckTypes(false)
	} else if !(this_.options.GetCheckTypes()) {
		this_.options.SetWarningLevel(DiagnosticGroup_ForType(StaticGet_RhinoErrorReporter_TYPE_PARSE_ERROR()), StaticGet_CheckLevel_OFF())
	}
	if !(this_.options.GetCheckTypes()) {
		this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_BOUNDED_GENERICS(), StaticGet_CheckLevel_OFF())
	}
	if !(this_.options.GetCheckSymbols()) && !(this_.options.Enables(StaticGet_DiagnosticGroups_CHECK_VARIABLES())) {
		this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_CHECK_VARIABLES(), StaticGet_CheckLevel_OFF())
	}
	if this_.options.GetSkipNonTranspilationPasses() && !(this_.options.Enables(StaticGet_DiagnosticGroups_CHECK_VARIABLES())) {
		this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_CHECK_VARIABLES(), StaticGet_CheckLevel_OFF())
	}
	if this_.options.GetSkipNonTranspilationPasses() && !(this_.options.Enables(StaticGet_DiagnosticGroups_MISSING_PROVIDE())) {
		this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_MISSING_PROVIDE(), StaticGet_CheckLevel_OFF())
	}
	if this_.options.ShouldDisambiguateProperties() && func() bool {
		re_9 := this_.options.GetWarningsGuard().MustRunChecks(StaticGet_DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION())
		ae_10 := pkg_base.StaticGet_Tri_UNKNOWN()
		if re_9 == 0 {
			j2gort.ThrowNPE()
		}
		return re_9 == ae_10
	}() {
		this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION(), StaticGet_CheckLevel_ERROR())
	}
}

// Java method: Compiler.getTypedAstDeserializer(com.google.javascript.jscomp.SourceFile file)
func (this_ *Compiler) GetTypedAstDeserializer(file *SourceFile) j2gort.Ext_java_util_function_Supplier {
	_ = file
	if this_.typedAstFilesystem == nil {
		return nil
	}
	var ast j2gort.Ext_java_util_function_Supplier = j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](this_.typedAstFilesystem.Remove_1(j2gort.PtrToAny(file)))
	_ = ast
	if !(((ast != nil) || javalang.StringStartsWith(j2gort.CheckNotNullString(file.GetName()), " [synthetic:")) || AbstractCompiler_IsFillFileName(file.GetName())) {
		j2gort.CheckStateFmt(false, "TypedAST filesystem initialized, but missing requested file: %s", j2gort.PtrToAny(file))
	}
	return ast
}

// Java method: Compiler.initWithTypedAstFilesystem(java.util.List<com.google.javascript.jscomp.SourceFile> externs, java.util.List<com.google.javascript.jscomp.SourceFile> sources, com.google.javascript.jscomp.CompilerOptions options, java.io.InputStream typedAstListStream)
func (this_ *Compiler) InitWithTypedAstFilesystem(externs j2gort.Ext_java_util_List, sources j2gort.Ext_java_util_List, options *CompilerOptions, typedAstListStream j2gort.Ext_java_io_InputStream) {
	_ = externs
	_ = sources
	_ = options
	_ = typedAstListStream
	var files j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](externs))).AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](sources))).Build())
	_ = files
	options.SetMergedPrecompiledLibraries(true)
	this_.InitOptions(options)
	this_.Init(externs, sources, options)
	this_.priv_Compiler_mergeAndDeserializeTypedAsts_1447(files, typedAstListStream, options)
}

// Java method: Compiler.initChunksWithTypedAstFilesystem(java.util.List<com.google.javascript.jscomp.SourceFile> externs, java.util.List<com.google.javascript.jscomp.JSChunk> chunks, com.google.javascript.jscomp.CompilerOptions options, java.io.InputStream typedAstListStream)
func (this_ *Compiler) InitChunksWithTypedAstFilesystem(externs j2gort.Ext_java_util_List, chunks j2gort.Ext_java_util_List, options *CompilerOptions, typedAstListStream j2gort.Ext_java_io_InputStream) {
	_ = externs
	_ = chunks
	_ = options
	_ = typedAstListStream
	var filesBuilder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = filesBuilder
	filesBuilder.AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](externs))
	{
		cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunks))
		for cur_11.HasNext() {
			rawElem_12 := cur_11.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_12)
			_ = chunk
			{
				cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetInputs()))
				for cur_13.HasNext() {
					rawElem_14 := cur_13.Next()
					input := j2gort.CastPtr[CompilerInput](rawElem_14)
					_ = input
					filesBuilder.Add_1_any(j2gort.PtrToAny(input.GetSourceFile()))
				}
			}
		}
	}
	var files j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](filesBuilder.Build())
	_ = files
	options.SetMergedPrecompiledLibraries(true)
	this_.InitOptions(options)
	this_.InitChunks(externs, chunks, options)
	this_.priv_Compiler_mergeAndDeserializeTypedAsts_1447(files, typedAstListStream, options)
}

// Java method: Compiler.mergeAndDeserializeTypedAsts(com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.SourceFile> requiredInputFiles, java.io.InputStream typedAstListStream, com.google.javascript.jscomp.CompilerOptions options)
func (this_ *Compiler) priv_Compiler_mergeAndDeserializeTypedAsts_1447(requiredInputFiles j2gort.Ext_com_google_common_collect_ImmutableSet, typedAstListStream j2gort.Ext_java_io_InputStream, options *CompilerOptions) {
	_ = requiredInputFiles
	_ = typedAstListStream
	_ = options
	if !(this_.typedAstFilesystem == nil) {
		j2gort.CheckState(false)
	}
	this_.MaybeSetTracker()
	this_.SetLifeCycleStage(StaticGet_AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC())
	var deserializeTypes bool = options.RequiresTypesForOptimization()
	_ = deserializeTypes
	var astData TypedAstDeserializer_DeserializedAst = j2gort.AsIface[TypedAstDeserializer_DeserializedAst](this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		var tracer *Tracer = this_.NewTracer("deserializeTypedAst")
		_ = tracer
		defer func() {
			this_.StopTracer(tracer, "deserializeTypedAst")
		}()
		{
			return TypedAstDeserializer_DeserializeFullAst(this_, this_.u_SYNTHETIC_EXTERNS_FILE, requiredInputFiles, typedAstListStream, deserializeTypes, options.GetResolveSourceMapAnnotations(), options.GetParseInlineSourceMaps())
		}
		panic("unreachable")
	})))
	_ = astData
	this_.typedAstFilesystem = astData.GetFilesystem()
	this_.externProperties = astData.GetExternProperties()
	if astData.GetColorRegistry() != nil {
		this_.colorRegistry = j2gort.CastPtr[pkg_colors.ColorRegistry](astData.GetColorRegistry().OrNull())
	} else {
		this_.colorRegistry = nil
	}
	this_.SetTypeCheckingHasRun(deserializeTypes)
	{
		cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(astData.GetRuntimeLibraries()))
		for cur_15.HasNext() {
			library := cur_15.NextStr()
			_ = library
			this_.runtimeJsLibManager.EnsureLibraryInjected(library, false)
		}
	}
	this_.GetSynthesizedExternsInput()
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		this_.priv_Compiler_parseForCompilationInternal_1472()
		return nil
	}))
}

// Java method: Compiler.initRuntimeLibraryTypedAsts(com.google.common.base.Optional<com.google.javascript.jscomp.serialization.ColorPool.Builder> colorPoolBuilder)
func (this_ *Compiler) InitRuntimeLibraryTypedAsts(colorPoolBuilder j2gort.Ext_com_google_common_base_Optional) {
	_ = colorPoolBuilder
	if !(this_.runtimeLibraryTypedAsts == nil) {
		j2gort.CheckState(false)
	}
	var path string = j2gort.Ext_java_lang_String_Join_2_any_slice_any("", []any{"/runtime_libs.typedast"})
	_ = path
	var astData TypedAstDeserializer_DeserializedAst = TypedAstDeserializer_DeserializeRuntimeLibraries(this_, this_.u_SYNTHETIC_EXTERNS_FILE, colorPoolBuilder, func() j2gort.Ext_java_io_InputStream {
		if b_, ok_ := j2gort.ClassResourceBytes(j2gort.ClassLiteral("com.google.javascript.jscomp.Compiler", "Compiler"), path); ok_ {
			s_ := &j2gort.Ext_java_io_InputStreamBase{Core_: b_}
			s_.Self_ = s_
			return s_
		}
		return nil
	}(), this_.GetOptions().GetResolveSourceMapAnnotations(), this_.GetOptions().GetParseInlineSourceMaps())
	_ = astData
	var runtimeLibraryTypedAsts j2gort.OrderedMap[string, j2gort.Ext_java_util_function_Supplier] = j2gort.NewOrderedMap[string, j2gort.Ext_java_util_function_Supplier]()
	_ = runtimeLibraryTypedAsts
	{
		cur_16 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](astData.GetFilesystem().EntrySet())))
		for cur_16.HasNext() {
			rawElem_17 := cur_16.Next()
			library := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](rawElem_17)
			_ = library
			j2gort.OrderedMapComputeIfAbsent(&runtimeLibraryTypedAsts, j2gort.CastPtr[SourceFile](library.GetKey()).GetName(), nil, func(lp_22 string) j2gort.Ext_java_util_function_Supplier {
				f := lp_22
				_ = f
				return j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](library.GetValue())
			})

		}
	}
	this_.runtimeLibraryTypedAsts = runtimeLibraryTypedAsts.ToImmutableMap()
}

// Java method: Compiler.init(java.util.List<com.google.javascript.jscomp.SourceFile> externs, java.util.List<com.google.javascript.jscomp.SourceFile> sources, com.google.javascript.jscomp.CompilerOptions options)
func (this_ *Compiler) Init(externs j2gort.Ext_java_util_List, sources j2gort.Ext_java_util_List, options *CompilerOptions) {
	_ = externs
	_ = sources
	_ = options
	var val_23 JSChunk
	InitJSChunk(&val_23, "$strong$")
	var chunk *JSChunk = &val_23
	_ = chunk
	{
		cur_24 := j2gort.NewIterCursor(j2gort.CheckNotNull(sources))
		for cur_24.HasNext() {
			rawElem_25 := cur_24.Next()
			source := j2gort.CastPtr[SourceFile](rawElem_25)
			_ = source
			chunk.Add_1_CompilerInput(NewCompilerInput_2_SourceFile_bool(source, false))
		}
	}
	var chunks j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_int32(int32(1))
	_ = chunks
	chunks.Add_1(chunk)
	this_.InitChunks(externs, chunks, options)
	this_.priv_Compiler_addFilesToSourceMap_1558(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](sources))
}

// Java method: Compiler.initChunks(java.util.List<com.google.javascript.jscomp.SourceFile> externs, java.util.List<com.google.javascript.jscomp.JSChunk> chunks, com.google.javascript.jscomp.CompilerOptions options)
func (this_ *Compiler) InitChunks(externs j2gort.Ext_java_util_List, chunks j2gort.Ext_java_util_List, options *CompilerOptions) {
	_ = externs
	_ = chunks
	_ = options
	this_.InitOptions(options)
	this_.priv_Compiler_checkFirstChunk_1451(chunks)
	this_.externs.Clear()
	{
		cur_26 := j2gort.NewIterCursor(j2gort.CheckNotNull(externs))
		for cur_26.HasNext() {
			rawElem_27 := cur_26.Next()
			file := j2gort.CastPtr[SourceFile](rawElem_27)
			_ = file
			this_.externs.Add_1(NewCompilerInput_2_SourceFile_bool(file, true))
		}
	}
	var didRet_28 bool
	func() {
		defer func() {
			if rec_29 := recover(); rec_29 != nil {
				jex_30 := javalang.RecoverException(rec_29)
				_ = jex_30
				if j2gort.InstanceOf[*JSChunkGraph_ChunkDependenceException](jex_30) {
					e := j2gort.AsIface[*JSChunkGraph_ChunkDependenceException](jex_30)
					_ = e
					this_.Report_1(JSError_Make_2(StaticGet_Compiler_CHUNK_DEPENDENCY_ERROR(), []string{e.GetChunk().GetName(), e.GetDependentChunk().GetName()}))
					didRet_28 = true
					return
				} else {
					panic(rec_29)
				}
			}
		}()
		this_.chunkGraph = NewJSChunkGraph_1_Ext_java_util_List(chunks)
	}()
	if didRet_28 {
		return
	}
	this_.priv_Compiler_fillEmptyChunks_1452(this_.GetChunks())
	this_.commentsPerFile = j2gort.NewExt_java_util_concurrent_ConcurrentHashMap_1_int32(this_.chunkGraph.GetInputCount())
	this_.InitBasedOnOptions()
	this_.InitInputsByIdMap()
	this_.priv_Compiler_initAST_1455()
	if this_.IsDebugLoggingEnabled() || options.ShouldPrintConfig() {
		this_.PrintConfig()
	}
	func() {
		var chunkGraphLog pkg_diagnostic.LogFile = this_.CreateOrReopenLog(j2gort.GetObjectClass(this_), "chunk_graph.dot", []string{})
		_ = chunkGraphLog
		defer j2gort.CloseResource(chunkGraphLog)
		if chunkGraphLog.IsLogging() {
			chunkGraphLog.Log_1_string(DotFormatter_ToDot_1_GraphvizGraph(j2gort.AsIface[pkg_graph.GraphvizGraph](this_.chunkGraph.ToGraphvizGraph())))
		}
	}()
	this_.runtimeJsLibManager = RuntimeJsLibManager_Create(options.GetRuntimeLibraryMode(), func() RuntimeJsLibManager_ResourceProvider {
		bound_31 := this_
		return NewLambda_RuntimeJsLibManager_ResourceProvider(func(mrp0_ string, mrp1_ string) pkg_rhino.Node {
			return bound_31.priv_Compiler_loadResourceContents_1556(mrp0_, mrp1_)
		})
	}(), this_.changeTracker, func() j2gort.Ext_java_util_function_Supplier {
		if options.GetRuntimeLibraryMode() == StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES() {
			return func() j2gort.Ext_java_util_function_Supplier {
				bound_32 := this_
				return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.BoxNode(bound_32.priv_Compiler_getSynthesizedExternsRoot_1554()) })
			}()
		}
		return func() j2gort.Ext_java_util_function_Supplier {
			bound_33 := this_
			return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.BoxNode(bound_33.priv_Compiler_getNodeForRuntimeCodeInsertion_1557()) })
		}()
	}())
}

// Java method: Compiler.getRuntimeJsLibManager()
func (this_ *Compiler) GetRuntimeJsLibManager() *RuntimeJsLibManager {
	return this_.runtimeJsLibManager
}

// Java method: Compiler.initBasedOnOptions()
func (this_ *Compiler) InitBasedOnOptions() {
	this_.inputSourceMaps.PutAll(j2gort.AsIface[j2gort.Ext_java_util_Map](this_.options.GetInputSourceMaps()))
	if this_.options.ShouldGatherSourceMapInfo() {
		this_.sourceMap = this_.options.GetSourceMapFormat().GetInstance()
		this_.sourceMap.SetPrefixMappings(j2gort.AsIface[j2gort.Ext_java_util_List](this_.options.GetSourceMapLocationMappings()))
		if this_.options.GetApplyInputSourceMaps() {
			this_.sourceMap.SetSourceFileMapping(this_)
			if this_.options.GetSourceMapIncludeSourcesContent() {
				{
					cur_34 := j2gort.NewMapValCursor(this_.inputSourceMaps)
					for cur_34.HasNext() {
						rawElem_35 := cur_34.Next()
						inputSourceMap := j2gort.CastPtr[SourceMapInput](rawElem_35)
						_ = inputSourceMap
						this_.priv_Compiler_addSourceMapSourceFiles_1538(inputSourceMap)
					}
				}
			}
		}
	}
}

// Java method: Compiler.checkFirstChunk(java.util.List<com.google.javascript.jscomp.JSChunk> chunks)
func (this_ *Compiler) priv_Compiler_checkFirstChunk_1451(chunks j2gort.Ext_java_util_List) {
	_ = chunks
	if chunks.IsEmpty() {
		this_.Report_1(JSError_Make_2(StaticGet_Compiler_EMPTY_CHUNK_LIST_ERROR(), []string{}))
	} else if j2gort.CastPtr[JSChunk](chunks.Get(int32(0))).GetInputs().IsEmpty() && (chunks.Size() > int32(1)) {
		this_.Report_1(JSError_Make_2(StaticGet_Compiler_EMPTY_ROOT_CHUNK_ERROR(), []string{j2gort.CastPtr[JSChunk](chunks.Get(int32(0))).GetName()}))
	}
}

// Java static method: Compiler.joinPathParts(java.lang.String[] pathParts)
func Compiler_JoinPathParts(pathParts []string) string {
	clinit_Compiler()
	_ = pathParts
	return StaticGet_Compiler_PathJoiner().Join_1_slice_any(func() []any {
		arr_36 := pathParts
		if arr_36 == nil {
			return nil
		}
		out_37 := make([]any, len(arr_36))
		for i_, el_38 := range arr_36 {
			out_37[i_] = j2gort.BoxString(el_38)
		}
		return out_37
	}())
}

// Java method: Compiler.fillEmptyChunks(java.lang.Iterable<com.google.javascript.jscomp.JSChunk> chunks)
func (this_ *Compiler) priv_Compiler_fillEmptyChunks_1452(chunks j2gort.Ext_java_lang_Iterable) {
	_ = chunks
	{
		cur_39 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunks))
		for cur_39.HasNext() {
			rawElem_40 := cur_39.Next()
			chunk := j2gort.CastPtr[JSChunk](rawElem_40)
			_ = chunk
			{
				var iifeCond_1 bool
				{

					sr_41 := chunk.GetName()
					se_42 := "$weak$"
					iifeCond_1 = j2gort.CheckNotNullString(sr_41) == se_42
				}
				if !(iifeCond_1) && chunk.GetInputs().IsEmpty() {
					var input *CompilerInput = NewCompilerInput_1(SourceFile_FromCode_2(AbstractCompiler_CreateFillFileName(chunk.GetName()), ""))
					_ = input
					input.SetCompiler(this_)
					chunk.Add_1_CompilerInput(input)
				}
			}
		}
	}
}

// Java method: Compiler.rebuildInputsFromChunks()
func (this_ *Compiler) RebuildInputsFromChunks() {
	this_.InitInputsByIdMap()
}

// Java method: Compiler.initInputsByIdMap()
func (this_ *Compiler) InitInputsByIdMap() {
	this_.inputsById.Clear()
	{
		cur_43 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.externs))
		for cur_43.HasNext() {
			rawElem_44 := cur_43.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_44)
			_ = input
			var previous *CompilerInput = this_.priv_Compiler_putCompilerInput_1494(input)
			_ = previous
			if previous != nil {
				this_.Report_1(JSError_Make_2(StaticGet_Compiler_DUPLICATE_EXTERN_INPUT(), []string{input.GetName()}))
			}
		}
	}
	{
		cur_45 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.chunkGraph.GetAllInputs()))
		for cur_45.HasNext() {
			rawElem_46 := cur_45.Next()
			input_2 := j2gort.CastPtr[CompilerInput](rawElem_46)
			_ = input_2
			var previous_2 *CompilerInput = this_.priv_Compiler_putCompilerInput_1494(input_2)
			_ = previous_2
			if previous_2 != nil {
				this_.Report_1(JSError_Make_2(StaticGet_Compiler_DUPLICATE_INPUT(), []string{input_2.GetName()}))
			}
		}
	}
}

// Java method: Compiler.initAST()
func (this_ *Compiler) priv_Compiler_initAST_1455() {
	this_.jsRoot = pkg_rhino.IR_Root([]pkg_rhino.Node{})
	this_.externsRoot = pkg_rhino.IR_Root([]pkg_rhino.Node{})
	this_.externAndJsRoot = pkg_rhino.IR_Root([]pkg_rhino.Node{this_.externsRoot, this_.jsRoot})
}

// Java method: Compiler.compile(com.google.javascript.jscomp.SourceFile extern, com.google.javascript.jscomp.SourceFile input, com.google.javascript.jscomp.CompilerOptions options)
func (this_ *Compiler) Compile_3_SourceFile_SourceFile_CompilerOptions(extern *SourceFile, input *SourceFile, options *CompilerOptions) *Result {
	_ = extern
	_ = input
	_ = options
	return this_.Compile_3_Ext_java_util_List_Ext_java_util_List_CompilerOptions(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(extern)), j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(input)), options)
}

// Java method: Compiler.compile(java.util.List<com.google.javascript.jscomp.SourceFile> externs, java.util.List<com.google.javascript.jscomp.SourceFile> inputs, com.google.javascript.jscomp.CompilerOptions options)
func (this_ *Compiler) Compile_3_Ext_java_util_List_Ext_java_util_List_CompilerOptions(externs j2gort.Ext_java_util_List, inputs j2gort.Ext_java_util_List, options *CompilerOptions) *Result {
	_ = externs
	_ = inputs
	_ = options
	if !(this_.jsRoot == 0) {
		j2gort.CheckState(false)
	}
	func() {
		defer func() {
			this_.GenerateReport()
		}()
		this_.Init(externs, inputs, options)
		if !(this_.HasErrors()) {
			this_.ParseForCompilation()
		}
		if !(this_.HasErrors()) {
			if options.GetInstrumentForCoverageOnly() {
				this_.InstrumentForCoverage()
			} else {
				this_.Stage1Passes()
				if !(this_.HasErrors()) {
					this_.Stage2Passes(StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS())
					if !(this_.HasErrors()) {
						this_.Stage3Passes()
					}
				}
			}
			this_.PerformPostCompilationTasks()
		}
	}()
	return this_.GetResult()
}

// Java method: Compiler.generateReport()
func (this_ *Compiler) GenerateReport() {
	var t *Tracer = this_.NewTracer("generateReport")
	_ = t
	this_.errorManager.GenerateReport()
	this_.StopTracer(t, "generateReport")
}

// Java method: Compiler.compileChunks(java.util.List<com.google.javascript.jscomp.SourceFile> externs, java.util.List<com.google.javascript.jscomp.JSChunk> chunks, com.google.javascript.jscomp.CompilerOptions options)
func (this_ *Compiler) CompileChunks(externs j2gort.Ext_java_util_List, chunks j2gort.Ext_java_util_List, options *CompilerOptions) *Result {
	_ = externs
	_ = chunks
	_ = options
	if !(this_.jsRoot == 0) {
		j2gort.CheckState(false)
	}
	func() {
		defer func() {
			this_.GenerateReport()
		}()
		this_.InitChunks(externs, chunks, options)
		if !(this_.HasErrors()) {
			this_.ParseForCompilation()
		}
		if !(this_.HasErrors()) {
			if options.GetInstrumentForCoverageOnly() {
				this_.InstrumentForCoverage()
			} else {
				this_.Stage1Passes()
				if !(this_.HasErrors()) {
					this_.Stage2Passes(StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS())
					if !(this_.HasErrors()) {
						this_.Stage3Passes()
					}
				}
			}
			this_.PerformPostCompilationTasks()
		}
	}()
	return this_.GetResult()
}

// Java method: Compiler.stage1Passes()
func (this_ *Compiler) Stage1Passes() {
	if !(this_.chunkGraph != nil) {
		j2gort.CheckStateStr(false, "No inputs. Did you call init() or initChunks()?")
	}
	if !(!(this_.HasErrors())) {
		j2gort.CheckState(false)
	}
	if !(!(this_.options.GetInstrumentForCoverageOnly())) {
		j2gort.CheckState(false)
	}
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		this_.priv_Compiler_performChecks_1466()
		return nil
	}))
}

// Java method: Compiler.stage2Passes(com.google.javascript.jscomp.CompilerOptions.SegmentOfCompilationToRun segmentOfCompilationToRun)
func (this_ *Compiler) Stage2Passes(segmentOfCompilationToRun CompilerOptions_SegmentOfCompilationToRun) {
	_ = segmentOfCompilationToRun
	if !(((segmentOfCompilationToRun == StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS()) || (segmentOfCompilationToRun == StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF())) || (segmentOfCompilationToRun == StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF())) {
		j2gort.CheckStateFmt(false, "Unsupported segment of compilation to run in stage 2: %s", j2gort.BoxNullableVal(segmentOfCompilationToRun))
	}
	if !(this_.chunkGraph != nil) {
		j2gort.CheckStateStr(false, "No inputs. Did you call init() or initChunks()?")
	}
	if !(!(this_.HasErrors())) {
		j2gort.CheckState(false)
	}
	if !(!(this_.options.GetInstrumentForCoverageOnly())) {
		j2gort.CheckState(false)
	}
	var weakChunk *JSChunk = this_.chunkGraph.GetChunkByName("$weak$")
	_ = weakChunk
	if weakChunk != nil {
		{
			cur_47 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.chunkGraph.GetAllInputs()))
			for cur_47.HasNext() {
				rawElem_48 := cur_47.Next()
				i := j2gort.CastPtr[CompilerInput](rawElem_48)
				_ = i
				if i.GetSourceFile().IsWeak() {
					if !(i.GetChunk() == weakChunk) {
						j2gort.CheckStateStr(false, "Expected all weak files to be in the weak chunk.")
					}
				}
			}
		}
	}
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		if this_.options.ShouldOptimize() {
			this_.PerformTranspilationAndOptimizations(segmentOfCompilationToRun)
		}
		return nil
	}))
}

// Java method: Compiler.stage3Passes()
func (this_ *Compiler) Stage3Passes() {
	if !(this_.chunkGraph != nil) {
		j2gort.CheckStateStr(false, "No inputs. Did you call init() or initChunks()?")
	}
	if !(!(this_.HasErrors())) {
		j2gort.CheckStateFmt(false, "Cannot run stage3Passes with errors present: %s", this_.GetErrors())
	}
	if !(!(this_.options.GetInstrumentForCoverageOnly())) {
		j2gort.CheckState(false)
	}
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		if this_.options.ShouldOptimize() {
			this_.PerformFinalizations()
		}
		return nil
	}))
}

// Java method: Compiler.disableThreads()
func (this_ *Compiler) DisableThreads() {
	this_.compilerExecutor.DisableThreads()
}

// Java method: Compiler.setTimeout(int timeout)
func (this_ *Compiler) SetTimeout(timeout int32) {
	_ = timeout
	this_.compilerExecutor.SetTimeout(timeout)
}

// Java method: Compiler.runInCompilerThread(java.util.concurrent.Callable<T> callable)
func (this_ *Compiler) RunInCompilerThread(callable j2gort.Ext_java_util_concurrent_Callable) any {
	_ = callable
	return this_.compilerExecutor.RunInCompilerThread(callable, ((this_.options != nil) && this_.options.GetTracerMode().IsOn()))
}

// Java method: Compiler.performChecks()
func (this_ *Compiler) priv_Compiler_performChecks_1466() {
	if this_.options.GetSkipNonTranspilationPasses() {
		this_.WhitespaceOnlyPasses()
		if this_.options.NeedsTranspilationFrom(this_.options.GetLanguageIn().ToFeatureSet()) {
			this_.TranspileAndDontCheck()
		}
	} else {
		this_.Check()
	}
}

// Java method: Compiler.performPostCompilationTasks()
func (this_ *Compiler) PerformPostCompilationTasks() {
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		this_.priv_Compiler_performPostCompilationTasksInternal_1468()
		return nil
	}))
}

// Java method: Compiler.performPostCompilationTasksInternal()
func (this_ *Compiler) priv_Compiler_performPostCompilationTasksInternal_1468() {
	if this_.options.GetDevMode() == StaticGet_CompilerOptions_DevMode_START_AND_END() {
		this_.priv_Compiler_runValidityCheck_1483()
	}
	ManageClosureUnawareCode_Unwrap(this_).Process_2(this_.externsRoot, this_.jsRoot)
	if this_.tracker != nil {
		if this_.options.GetTracerOutput() == nil {
			this_.tracker.OutputTracerReport(this_.outStream)
		} else {
			func() {
				defer func() {
					if rec_49 := recover(); rec_49 != nil {
						jex_50 := javalang.RecoverException(rec_49)
						_ = jex_50
						if j2gort.InstanceOf[javalang.IException](jex_50) {
							e := j2gort.AsIface[javalang.IException](jex_50)
							_ = e
							panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
						} else {
							panic(rec_49)
						}
					}
				}()
				var out *j2gort.Ext_java_io_PrintStream = j2gort.NewExt_java_io_PrintStream_1_Ext_java_io_OutputStream(j2gort.Ext_java_nio_file_Files_NewOutputStream(this_.options.GetTracerOutput(), []j2gort.Ext_java_nio_file_OpenOption{}))
				_ = out
				defer j2gort.CloseResource(out)
				this_.tracker.OutputTracerReport(out)
			}()
		}
	}
}

// Java method: Compiler.instrumentForCoverage()
func (this_ *Compiler) InstrumentForCoverage() {
	if !(this_.chunkGraph != nil) {
		j2gort.CheckStateStr(false, "No inputs. Did you call init() or initChunks()?")
	}
	if !(!(this_.HasErrors())) {
		j2gort.CheckState(false)
	}
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		if !(this_.options.GetInstrumentForCoverageOnly()) {
			j2gort.CheckState(false)
		}
		if !(!(this_.HasErrors())) {
			j2gort.CheckState(false)
		}
		if this_.options.GetInstrumentForCoverageOption() != StaticGet_CompilerOptions_InstrumentOption_NONE() {
			this_.priv_Compiler_instrumentForCoverageInternal_1470(this_.options.GetInstrumentForCoverageOption())
		}
		return nil
	}))
}

// Java method: Compiler.instrumentForCoverageInternal(com.google.javascript.jscomp.CompilerOptions.InstrumentOption instrumentForCoverageOption)
func (this_ *Compiler) priv_Compiler_instrumentForCoverageInternal_1470(instrumentForCoverageOption CompilerOptions_InstrumentOption) {
	_ = instrumentForCoverageOption
	var tracer *Tracer = this_.NewTracer("instrumentationPass")
	_ = tracer
	this_.Process(NewCoverageInstrumentationPass_4(this_, StaticGet_CoverageInstrumentationPass_CoverageReach_ALL(), instrumentForCoverageOption, this_.options.GetProductionInstrumentationArrayName()))
	this_.StopTracer(tracer, "instrumentationPass")
}

// Java method: Compiler.parseForCompilation()
func (this_ *Compiler) ParseForCompilation() {
	if !(this_.typedAstFilesystem == nil) {
		j2gort.CheckStateStr(false, "Unnecessary if initWithTypedAstFilesystem was called")
	}
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		this_.priv_Compiler_parseForCompilationInternal_1472()
		return nil
	}))
}

// Java method: Compiler.parseForCompilationInternal()
func (this_ *Compiler) priv_Compiler_parseForCompilationInternal_1472() {
	CompilerOptionsPreprocessor_Preprocess(this_.options)
	this_.MaybeSetTracker()
	this_.ParseInputs()
}

// Java method: Compiler.parse()
func (this_ *Compiler) Parse_0() {
	this_.ParseInputs()
}

// Java method: Compiler.parse(com.google.javascript.jscomp.SourceFile file)
func (this_ *Compiler) Parse_1(file *SourceFile) pkg_rhino.Node {
	_ = file
	this_.InitCompilerOptionsIfTesting()
	StaticGet_Compiler_Logger().Finest_1_string(("Parsing: " + j2gort.StrValString(file.GetName())))
	return NewCompilerInput_1(file).GetAstRoot(this_)
}

// Java method: Compiler.getPassConfig()
func (this_ *Compiler) GetPassConfig() PassConfig {
	if this_.passes == nil || this_.passes.IsJavaNil_() {
		this_.passes = this_.CreatePassConfigInternal()
	}
	return this_.passes
}

// Java method: Compiler.createPassConfigInternal()
func (this_ *Compiler) CreatePassConfigInternal() PassConfig {
	return NewDefaultPassConfig(this_.options)
}

// Java method: Compiler.setPassConfig(com.google.javascript.jscomp.PassConfig passes)
func (this_ *Compiler) SetPassConfig(passes PassConfig) {
	_ = passes
	j2gort.CheckNotNullIface(passes)
	if this_.passes != nil && !this_.passes.IsJavaNil_() {
		j2gort.CheckStateStr(false, "setPassConfig was already called")
	}
	this_.passes = passes
}

// Java method: Compiler.whitespaceOnlyPasses()
func (this_ *Compiler) WhitespaceOnlyPasses() {
	this_.priv_Compiler_runCustomPasses_1484(StaticGet_CustomPassExecutionTime_BEFORE_CHECKS())
	var t *Tracer = this_.NewTracer("runWhitespaceOnlyPasses")
	_ = t
	defer func() {
		this_.StopTracer(t, "runWhitespaceOnlyPasses")
	}()
	{
		{
			cur_51 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetPassConfig().GetWhitespaceOnlyPasses().Build()))
			for cur_51.HasNext() {
				rawElem_52 := cur_51.Next()
				pf := j2gort.AsIface[PassFactory](rawElem_52)
				_ = pf
				pf.Create(this_).Process_2(this_.externsRoot, this_.jsRoot)
			}
		}
	}
}

// Java method: Compiler.markTranspiledFiles()
func (this_ *Compiler) priv_Compiler_markTranspiledFiles_1479() {
	{
		var child pkg_rhino.Node = this_.jsRoot.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if !(child.IsScript()) {
				j2gort.CheckState(false)
			}
			var featureSet *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(child)
			_ = featureSet
			if !(this_.options.GetOutputFeatureSet().Contains_1_FeatureSet(featureSet)) {
				this_.SetTranspiledFiles(true)
				return
			}
		}
	}
}

// Java method: Compiler.transpileAndDontCheck()
func (this_ *Compiler) TranspileAndDontCheck() {
	var t *Tracer = this_.NewTracer("runTranspileOnlyPasses")
	_ = t
	this_.priv_Compiler_markTranspiledFiles_1479()
	defer func() {
		this_.StopTracer(t, "runTranspileOnlyPasses")
	}()
	{
		{
			cur_53 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetPassConfig().GetTranspileOnlyPasses().Build()))
			for cur_53.HasNext() {
				rawElem_54 := cur_53.Next()
				pf := j2gort.AsIface[PassFactory](rawElem_54)
				_ = pf
				if this_.HasErrors() {
					return
				}
				this_.BeforePass(pf.GetName())
				pf.Create(this_).Process_2(this_.externsRoot, this_.jsRoot)
				this_.AfterPass(pf.GetName())
			}
		}
	}
}

// Java method: Compiler.createPhaseOptimizer()
func (this_ *Compiler) priv_Compiler_createPhaseOptimizer_1481() *PhaseOptimizer {
	clinit_PhaseOptimizer()
	var val_55 PhaseOptimizer
	InitPhaseOptimizer(&val_55, this_, this_.tracker)
	var phaseOptimizer *PhaseOptimizer = &val_55
	_ = phaseOptimizer
	if this_.options.GetDevMode() == StaticGet_CompilerOptions_DevMode_EVERY_PASS() {
		phaseOptimizer.SetValidityCheck(this_.validityCheck)
	}
	if this_.options.GetCheckDeterminism() {
		phaseOptimizer.SetPrintAstHashcodes(true)
	}
	return phaseOptimizer
}

// Java method: Compiler.check()
func (this_ *Compiler) Check() {
	this_.priv_Compiler_runCustomPasses_1484(StaticGet_CustomPassExecutionTime_BEFORE_CHECKS())
	this_.phaseOptimizer = this_.priv_Compiler_createPhaseOptimizer_1481()
	this_.phaseOptimizer.Consume(j2gort.AsIface[j2gort.Ext_java_util_List](this_.GetPassConfig().GetChecks().Build()))
	this_.phaseOptimizer.Process_2(this_.externsRoot, this_.jsRoot)
	if this_.HasErrors() {
		return
	}
	this_.priv_Compiler_runCustomPasses_1484(StaticGet_CustomPassExecutionTime_BEFORE_OPTIMIZATIONS())
	this_.phaseOptimizer = nil
}

// Java method: Compiler.setExternExports(java.lang.String externExports)
func (this_ *Compiler) SetExternExports(externExports string) {
	_ = externExports
	this_.externExports = externExports
}

// Java method: Compiler.process(com.google.javascript.jscomp.CompilerPass p)
func (this_ *Compiler) Process(p CompilerPass) {
	_ = p
	p.Process_2(this_.externsRoot, this_.jsRoot)
}

// Java method: Compiler.runValidityCheck()
func (this_ *Compiler) priv_Compiler_runValidityCheck_1483() {
	this_.validityCheck.Create(this_).Process_2(this_.externsRoot, this_.jsRoot)
}

// Java method: Compiler.runCustomPasses(com.google.javascript.jscomp.CustomPassExecutionTime executionTime)
func (this_ *Compiler) priv_Compiler_runCustomPasses_1484(executionTime CustomPassExecutionTime) {
	_ = executionTime
	var customPasses []CompilerPass = j2gort.ImmListToSliceNullable[CompilerPass](this_.options.GetCustomPassesAt(executionTime))
	_ = customPasses
	var t *Tracer = this_.NewTracer("runCustomPasses")
	_ = t
	defer func() {
		this_.StopTracer(t, "runCustomPasses")
	}()
	{
		{

			for _, p := range customPasses {

				_ = p
				this_.Process(p)
			}
		}
	}
}

// Java method: Compiler.beforePass(java.lang.String passName)
func (this_ *Compiler) BeforePass(passName string) {
	_ = passName
	this_.currentPassIndex = this_.currentPassIndex + 1
}

// Java method: Compiler.afterPass(java.lang.String passName)
func (this_ *Compiler) AfterPass(passName string) {
	_ = passName
	this_.priv_Compiler_maybePrintSourceAfterEachPass_1485(passName)
}

// Java method: Compiler.getLifeCycleStage()
func (this_ *Compiler) GetLifeCycleStage() *AbstractCompiler_LifeCycleStage {
	return this_.stage
}

// Java method: Compiler.setLifeCycleStage(com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage stage)
func (this_ *Compiler) SetLifeCycleStage(stage *AbstractCompiler_LifeCycleStage) {
	_ = stage
	this_.stage = stage
}

// Java method: Compiler.createAstFactory()
func (this_ *Compiler) CreateAstFactory() *AstFactory {
	if this_.HasTypeCheckingRun() {
		if this_.HasOptimizationColors() {
			return AstFactory_CreateFactoryWithColors(this_.stage, this_.GetColorRegistry(), this_.runtimeJsLibManager)
		} else {
			return AstFactory_CreateFactoryWithTypes(this_.stage, this_.GetTypeRegistry(), this_.runtimeJsLibManager)
		}
	} else {
		return AstFactory_CreateFactoryWithoutTypes(this_.stage, this_.runtimeJsLibManager)
	}
}

// Java method: Compiler.createAstFactoryWithoutTypes()
func (this_ *Compiler) CreateAstFactoryWithoutTypes() *AstFactory {
	return AstFactory_CreateFactoryWithoutTypes(this_.stage, this_.runtimeJsLibManager)
}

// Java method: Compiler.getAstAnalyzer()
func (this_ *Compiler) GetAstAnalyzer() *AstAnalyzer {
	var options *CompilerOptions = this_.GetOptions()
	_ = options
	var useTypes bool = options.ShouldUseTypesForLocalOptimization()
	_ = useTypes
	var pureGetters bool = options.GetAssumeGettersArePure()
	_ = pureGetters
	var regexpRefs bool = this_.HasRegExpGlobalReferences()
	_ = regexpRefs
	var pureBuiltins bool = options.GetAssumePropertiesAreStaticallyAnalyzable()
	_ = pureBuiltins
	var mask int32 = ((((func() int32 {
		if useTypes {
			return int32(1)
		}
		return int32(0)
	}()) | (func() int32 {
		if pureGetters {
			return int32(2)
		}
		return int32(0)
	}())) | (func() int32 {
		if regexpRefs {
			return int32(4)
		}
		return int32(0)
	}())) | (func() int32 {
		if pureBuiltins {
			return int32(8)
		}
		return int32(0)
	}()))
	_ = mask
	var reg *pkg_rhino.JSTypeRegistry = this_.typeRegistry
	_ = reg
	var summary *AccessorSummary = this_.GetAccessorSummary()
	_ = summary
	var cached *AstAnalyzer = this_.cachedAstAnalyzer
	_ = cached
	if (((cached != nil) && (this_.cachedAstAnalyzerMask == mask)) && (this_.cachedAstAnalyzerRegistry == reg)) && (this_.cachedAstAnalyzerSummary == summary) {
		return cached
	}
	var analysisSettings *AstAnalyzer_Options = AstAnalyzer_Options_S_Builder().SetUseTypesForLocalOptimization(useTypes).SetAssumeGettersArePure(pureGetters).SetHasRegexpGlobalReferences(regexpRefs).SetAssumeKnownBuiltinsArePure(pureBuiltins).Build()
	_ = analysisSettings
	var val_58 AstAnalyzer
	InitAstAnalyzer(&val_58, analysisSettings, reg, summary)
	var analyzer *AstAnalyzer = &val_58
	_ = analyzer
	this_.cachedAstAnalyzer = analyzer
	this_.cachedAstAnalyzerMask = mask
	this_.cachedAstAnalyzerRegistry = reg
	this_.cachedAstAnalyzerSummary = summary
	return analyzer
}

// Java method: Compiler.createDefaultExpressionDecomposer()
func (this_ *Compiler) CreateDefaultExpressionDecomposer() *ExpressionDecomposer {
	return this_.CreateExpressionDecomposer(j2gort.AsIface[j2gort.Ext_com_google_common_base_Supplier](this_.GetUniqueNameIdSupplier()), j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0(), Scope_CreateGlobalScope(pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT())))
}

// Java method: Compiler.createExpressionDecomposer(java.util.function.Supplier<java.lang.String> uniqueNameIdSupplier, com.google.common.collect.ImmutableSet<java.lang.String> knownConstantFunctions, com.google.javascript.jscomp.Scope scope)
func (this_ *Compiler) CreateExpressionDecomposer(uniqueNameIdSupplier j2gort.Ext_java_util_function_Supplier, knownConstantFunctions j2gort.Ext_com_google_common_collect_ImmutableSet, scope *Scope) *ExpressionDecomposer {
	_ = uniqueNameIdSupplier
	_ = knownConstantFunctions
	_ = scope
	var enableBrokenIe11LocationAssign bool = pkg_parser.StaticGet_FeatureSet_ES5().Contains_1_FeatureSet(this_.GetOptions().GetOutputFeatureSet())
	_ = enableBrokenIe11LocationAssign
	return NewExpressionDecomposer_5_AbstractCompiler_Ext_java_util_function_Supplier_Ext_com_google_common_collect_ImmutableSet_Scope_bool(this_, uniqueNameIdSupplier, knownConstantFunctions, scope, enableBrokenIe11LocationAssign)
}

// Java method: Compiler.isDebugLoggingEnabled()
func (this_ *Compiler) IsDebugLoggingEnabled() bool {
	return (this_.GetOptions().GetDebugLogDirectory() != nil)
}

// Java method: Compiler.getDebugLogFilterList()
func (this_ *Compiler) GetDebugLogFilterList() j2gort.Ext_java_util_List {
	if this_.GetOptions().GetDebugLogFilter() == j2gort.NullString {
		return j2gort.NewExt_java_util_ArrayList()
	}
	return j2gort.SplitterSplitToListStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(44)).OmitEmptyStrings(), this_.GetOptions().GetDebugLogFilter())
}

// Java method: Compiler.createOrReopenLog(java.lang.Class<?> owner, java.lang.String firstNamePart, java.lang.String[] restNameParts)
func (this_ *Compiler) CreateOrReopenLog(owner *javalang.Class, firstNamePart string, restNameParts []string) pkg_diagnostic.LogFile {
	_ = owner
	_ = firstNamePart
	_ = restNameParts
	if !(this_.IsDebugLoggingEnabled()) {
		return pkg_diagnostic.LogFile_CreateNoOp()
	}
	var dir j2gort.Ext_java_nio_file_Path = this_.GetOptions().GetDebugLogDirectory()
	_ = dir
	var relativeParts j2gort.Ext_java_nio_file_Path = j2gort.Ext_java_nio_file_Path_Of_2(firstNamePart, restNameParts)
	_ = relativeParts
	var file j2gort.Ext_java_nio_file_Path = dir.Resolve_1_string(j2gort.CheckNotNullPtr(owner).GetSimpleName()).Resolve_1_Ext_java_nio_file_Path(relativeParts)
	_ = file
	var filters j2gort.Ext_java_util_List = this_.GetDebugLogFilterList()
	_ = filters
	if filters.IsEmpty() {
		return pkg_diagnostic.LogFile_CreateOrReopen(file)
	}
	{
		cur_59 := j2gort.NewIterCursor(j2gort.CheckNotNull(filters))
		for cur_59.HasNext() {
			filter := cur_59.NextStr()
			_ = filter
			if j2gort.JStrIndexOf(file.String(), j2gort.CheckNotNullString(filter)) >= 0 {
				return pkg_diagnostic.LogFile_CreateOrReopen(file)
			}
		}
	}
	return pkg_diagnostic.LogFile_CreateNoOp()
}

// Java method: Compiler.createOrReopenIndexedLog(java.lang.Class<?> owner, java.lang.String firstNamePart, java.lang.String[] restNameParts)
func (this_ *Compiler) CreateOrReopenIndexedLog(owner *javalang.Class, firstNamePart string, restNameParts []string) pkg_diagnostic.LogFile {
	_ = owner
	_ = firstNamePart
	_ = restNameParts
	if !(this_.currentPassIndex >= int32(0)) {
		j2gort.CheckStateMsg(false, j2gort.ValueOfInt(this_.currentPassIndex))
	}
	var index string = j2gort.Ext_com_google_common_base_Strings_PadStart(strconv.FormatInt(int64(this_.currentPassIndex), 10), int32(3), uint16(48))
	_ = index
	var length int32 = int32(len(restNameParts))
	_ = length
	if length == int32(0) {
		firstNamePart = (j2gort.StrValString(index) + "_" + j2gort.StrValString(firstNamePart))
	} else {
		{
			a_60 := restNameParts
			i_61 := (length - int32(1))
			a_60[i_61] = (j2gort.StrValString(index) + "_" + j2gort.StrValString(restNameParts[(length-int32(1))]))
		}
	}
	return this_.CreateOrReopenLog(owner, firstNamePart, restNameParts)
}

// Java method: Compiler.maybePrintSourceAfterEachPass(java.lang.String passName)
func (this_ *Compiler) priv_Compiler_maybePrintSourceAfterEachPass_1485(passName string) {
	_ = passName
	if !(this_.options.ShouldPrintSourceAfterEachPass()) {
		return
	}
	var currentJsSource string = this_.GetCurrentJsSource()
	_ = currentJsSource
	{
		var iifeCond_1 bool
		{

			sr_62 := currentJsSource
			se_63 := this_.lastJsSource
			iifeCond_1 = j2gort.CheckNotNullString(sr_62) == se_63
		}
		if iifeCond_1 {
			return
		}
	}
	if this_.IsDebugLoggingEnabled() {
		func() {
			var log pkg_diagnostic.LogFile = this_.CreateOrReopenIndexedLog(j2gort.GetObjectClass(this_), "source_after_pass", []string{passName})
			_ = log
			defer j2gort.CloseResource(log)
			if log.IsLogging() {
				log.Log_1_string(currentJsSource)
			}
		}()
	} else {
		j2gort.Ext_java_lang_System_F_Err.Println_0()
		if this_.debugMessage != j2gort.NullString {
			j2gort.Ext_java_lang_System_F_Err.Println_1_string(("// DEBUG: " + j2gort.StrValString(this_.debugMessage)))
		}
		j2gort.Ext_java_lang_System_F_Err.Println_1_string(("// " + j2gort.StrValString(passName) + " yields:"))
		j2gort.Ext_java_lang_System_F_Err.Println_1_string("// ************************************")
		j2gort.Ext_java_lang_System_F_Err.Println_1_string(currentJsSource)
		j2gort.Ext_java_lang_System_F_Err.Println_1_string("// ************************************")
	}
	this_.lastJsSource = currentJsSource
}

// Java method: Compiler.getCurrentJsSource()
func (this_ *Compiler) GetCurrentJsSource() string {
	this_.ResetAndIntitializeSourceMap()
	var fileNameRegexList j2gort.Ext_java_util_List = this_.options.GetFilesToPrintAfterEachPassRegexList()
	_ = fileNameRegexList
	var chunkNameRegexList j2gort.Ext_java_util_List = this_.options.GetChunksToPrintAfterEachPassRegexList()
	_ = chunkNameRegexList
	var qnameSet j2gort.LocalOrderedSet[string] = j2gort.LocalOrderedSetFromStrCol(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.options.GetQnameUsesToPrintAfterEachPassList()))
	_ = qnameSet
	var builder *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = builder
	if (fileNameRegexList.IsEmpty() && chunkNameRegexList.IsEmpty()) && qnameSet.IsEmpty() {
		return this_.ToSource_0()
	}
	if !(fileNameRegexList.IsEmpty()) {
		j2gort.CheckNotNullNode(this_.externsRoot)
		j2gort.CheckNotNullNode(this_.jsRoot)
		{
			cur_64 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.BoxNode(this_.externsRoot), j2gort.BoxNode(this_.jsRoot))))
			for cur_64.HasNext() {
				r := pkg_rhino.Node(cur_64.NextNode())
				_ = r
				{
					var fileNode pkg_rhino.Node = r.GetFirstChild()
					_ = fileNode
					for ; fileNode != 0; fileNode = fileNode.GetNext() {
						var fileName string = fileNode.GetSourceFileName()
						_ = fileName
						{
							cur_65 := j2gort.NewIterCursor(j2gort.CheckNotNull(fileNameRegexList))
							for cur_65.HasNext() {
								regex := cur_65.NextStr()
								_ = regex
								if j2gort.JStrMatches(fileName, regex) {
									var source string = ("// " + j2gort.StrValString(fileName) + "\n" + j2gort.StrValString(this_.ToSource_1_Node(fileNode)))
									_ = source
									builder.Append_1_string(source)
									break
								}
							}
						}
					}
				}
			}
		}
		if builder.Length() == int32(0) {
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](builder.Append_1_string("// No files matched any of: ")).Append_1_any_9815(fileNameRegexList)
		}
	}
	if !(chunkNameRegexList.IsEmpty()) {
		var chunks j2gort.Ext_java_lang_Iterable = j2gort.CheckNotNull(this_.GetChunks())
		_ = chunks
		var unmatchedChunkNames []string = []string{}
		_ = unmatchedChunkNames
		var val_66 Compiler_ChunkGraphAwareLicenseTracker
		InitCompiler_ChunkGraphAwareLicenseTracker(&val_66, this_)
		var lt *Compiler_ChunkGraphAwareLicenseTracker = &val_66
		_ = lt
		{
			cur_67 := j2gort.NewIterCursor(chunks)
			for cur_67.HasNext() {
				rawElem_68 := cur_67.Next()
				jsChunk := j2gort.CastPtr[JSChunk](rawElem_68)
				_ = jsChunk
				var chunkName string = jsChunk.GetName()
				_ = chunkName
				lt.SetCurrentChunkContext(jsChunk)
				{
					cur_69 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunkNameRegexList))
					for cur_69.HasNext() {
						regex_2 := cur_69.NextStr()
						_ = regex_2
						if j2gort.JStrMatches(chunkName, regex_2) {
							var source_2 string = ("// module '" + j2gort.StrValString(chunkName) + "'\n" + j2gort.StrValString(this_.ToSource_2(lt, jsChunk)))
							_ = source_2
							builder.Append_1_string(source_2)
							break
						}
					}
				}
				unmatchedChunkNames = append(unmatchedChunkNames, chunkName)
			}
		}
		if builder.Length() == int32(0) {
			builder.Append_1_string("// No chunks were matched:\n")
			for _, chunkName_2 := range unmatchedChunkNames {
				_ = chunkName_2
				j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](builder.Append_1_string("// ")).Append_1_string(chunkName_2)).Append_1_string("\n")
			}

		}
	}
	if !(qnameSet.IsEmpty()) {
		var originalToNewQNameMap j2gort.Ext_com_google_common_collect_Multimap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0())
		_ = originalToNewQNameMap
		var statementStream j2gort.Ext_java_util_stream_Stream = this_.priv_Compiler_getTopLevelStatements_1487(this_.jsRoot).Filter(j2gort.NewNodeLambda_Ext_java_util_function_Predicate(func(rawNode_71 uint32) bool {
			statement := pkg_rhino.Node(rawNode_71)
			_ = statement
			return nodeUtilHasPredFn_(statement, func(rawNode_72 uint32) bool {
				n := pkg_rhino.Node(rawNode_72)
				_ = n
				j2gort.CheckNotNullNode(n)
				if !(n.IsQualifiedName()) {
					return false
				}
				var qname string = n.GetQualifiedName()
				_ = qname
				if qnameSet.Contains(qname) {
					return true
				}
				var originalQname string = n.GetOriginalQualifiedName()
				_ = originalQname
				if (originalQname != j2gort.NullString) && qnameSet.Contains(originalQname) {
					j2gort.MultimapPutStr(originalToNewQNameMap, originalQname, j2gort.BoxString(qname))
					return true
				} else {
					return false
				}
			}, lambdaSingleton_compiler_3)
		}))
		_ = statementStream
		builder.Append_1_string("//\n")
		builder.Append_1_string("// closure-compiler: Printing all of the top-level statements\n")
		builder.Append_1_string("// that contain references to these qualified names.\n")
		builder.Append_1_string("//\n")
		{
			cur_74 := qnameSet.Iter()
			for cur_74.HasNext() {
				qname := cur_74.
					Next()
				_ = qname
				builder.Append_1_string(j2gort.Ext_java_lang_String_Format_2("// '%s'\n", []any{j2gort.BoxString(qname)}))
				{
					cur_75 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.MultimapGetStr(originalToNewQNameMap, qname))))
					for cur_75.HasNext() {
						newName := cur_75.NextStr()
						_ = newName
						builder.Append_1_string(j2gort.Ext_java_lang_String_Format_2("// '%s' (originally '%s')\n", []any{j2gort.BoxString(newName), j2gort.BoxString(qname)}))
					}
				}
			}
		}
		builder.Append_1_string("//\n")
		statementStream.ForEach(j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_76 any) {
			statement := j2gort.CastNode[pkg_rhino.Node](lp_76)
			_ = statement
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](builder.Append_1_string(j2gort.Ext_java_lang_String_Format_2("// %s\n", []any{j2gort.BoxString(statement.GetLocation())}))).Append_1_string(this_.ToSource_1_Node(statement))).Append_1_string("\n")
		}))
	}
	return builder.String()
}

// Java method: Compiler.getTopLevelStatements(com.google.javascript.rhino.Node root)
func (this_ *Compiler) priv_Compiler_getTopLevelStatements_1487(root pkg_rhino.Node) j2gort.Ext_java_util_stream_Stream {
	_ = root
	var builder j2gort.Ext_java_util_stream_Stream_Builder = j2gort.Ext_java_util_stream_Stream_S_Builder()
	_ = builder
	NodeTraversal_Traverse(this_, root, NewCompiler__Anon_13(this_, builder))
	return builder.Build()
}

// Java method: Compiler.getScriptNode(java.lang.String filename)
func (this_ *Compiler) GetScriptNode(filename string) pkg_rhino.Node {
	_ = filename
	return j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetStr(this_.scriptNodeByFilename, j2gort.CheckNotNullString(filename)))
}

// Java method: Compiler.newTracer(java.lang.String passName)
func (this_ *Compiler) NewTracer(passName string) *Tracer {
	_ = passName
	var comment string = (j2gort.StrValString(passName) + j2gort.StrValString(func() string {
		if this_.changeTracker.GetRecentChange().HasCodeChanged() {
			return " on recently changed AST"
		}
		return ""
	}()))
	_ = comment
	if this_.options.GetTracerMode().IsOn() && (this_.tracker != nil) {
		this_.tracker.RecordPassStart(passName, true)
	}
	return NewTracer_2("Compiler", comment)
}

// Java method: Compiler.stopTracer(com.google.javascript.jscomp.Tracer t, java.lang.String passName)
func (this_ *Compiler) StopTracer(t *Tracer, passName string) {
	_ = t
	_ = passName
	var result int64 = t.Stop_0()
	_ = result
	if this_.options.GetTracerMode().IsOn() && (this_.tracker != nil) {
		this_.tracker.RecordPassStop(passName, result)
	}
}

// Java method: Compiler.getResult()
func (this_ *Compiler) GetResult() *Result {
	return NewResult_12(this_.GetErrors(), this_.GetWarnings(), this_.variableMap, this_.propertyMap, nil, this_.stringMap, this_.instrumentationMapping, this_.sourceMap, this_.externExports, this_.cssNames, this_.idGeneratorMap, this_.transpiledFiles)
}

// Java method: Compiler.getErrors()
func (this_ *Compiler) GetErrors() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.errorManager == nil || this_.errorManager.IsJavaNil_() {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	} else {
		return this_.errorManager.GetErrors()
	}
}

// Java method: Compiler.getWarnings()
func (this_ *Compiler) GetWarnings() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.errorManager == nil || this_.errorManager.IsJavaNil_() {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	} else {
		return this_.errorManager.GetWarnings()
	}
}

// Java method: Compiler.getRoot()
func (this_ *Compiler) GetRoot() pkg_rhino.Node {
	return this_.externAndJsRoot
}

// Java method: Compiler.getAllowableFeatures()
func (this_ *Compiler) GetAllowableFeatures() *pkg_parser.FeatureSet {
	return this_.allowableFeatures
}

// Java method: Compiler.setAllowableFeatures(com.google.javascript.jscomp.parsing.parser.FeatureSet allowableFeatures)
func (this_ *Compiler) SetAllowableFeatures(allowableFeatures *pkg_parser.FeatureSet) {
	_ = allowableFeatures
	this_.allowableFeatures = allowableFeatures
}

// Java method: Compiler.markFeatureNotAllowed(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *Compiler) MarkFeatureNotAllowed(feature *pkg_parser.FeatureSet_Feature) {
	_ = feature
	this_.allowableFeatures = this_.allowableFeatures.Without_2(feature, []*pkg_parser.FeatureSet_Feature{})
}

// Java method: Compiler.markFeatureSetNotAllowed(com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet)
func (this_ *Compiler) MarkFeatureSetNotAllowed(featureSet *pkg_parser.FeatureSet) {
	_ = featureSet
	this_.allowableFeatures = this_.allowableFeatures.Without_1(featureSet)
}

// Java method: Compiler.getUniqueIdSupplier()
func (this_ *Compiler) GetUniqueIdSupplier() *UniqueIdSupplier {
	return this_.uniqueIdSupplier
}

// Java method: Compiler.nextUniqueNameId()
func (this_ *Compiler) priv_Compiler_nextUniqueNameId_1493() int32 {
	return func() int32 {
		old_77 := this_.uniqueNameId
		nw_78 := old_77 + 1
		this_.uniqueNameId = nw_78
		return old_77
	}()
}

// Java method: Compiler.getUniqueNameIdSupplier()
func (this_ *Compiler) GetUniqueNameIdSupplier() any {
	return j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
		return strconv.FormatInt(int64(this_.priv_Compiler_nextUniqueNameId_1493()), 10)
	})
}

// Java method: Compiler.areNodesEqualForInlining(com.google.javascript.rhino.Node n1, com.google.javascript.rhino.Node n2)
func (this_ *Compiler) AreNodesEqualForInlining(n1 pkg_rhino.Node, n2 pkg_rhino.Node) bool {
	_ = n1
	_ = n2
	if this_.options.ShouldAmbiguateProperties() || this_.options.ShouldDisambiguateProperties() {
		return n1.IsEquivalentToTyped(n2)
	} else {
		return n1.IsEquivalentTo_1(n2)
	}
}

// Java method: Compiler.getInput(com.google.javascript.rhino.InputId id)
func (this_ *Compiler) GetInput(id *pkg_rhino.InputId) *CompilerInput {
	_ = id
	if id == nil {
		return nil
	}
	return j2gort.CastPtr[CompilerInput](this_.inputsById.Get(j2gort.PtrToAny(id)))
}

// Java method: Compiler.putCompilerInput(com.google.javascript.jscomp.CompilerInput input)
func (this_ *Compiler) priv_Compiler_putCompilerInput_1494(input *CompilerInput) *CompilerInput {
	_ = input
	input.SetCompiler(this_)
	return j2gort.CastPtr[CompilerInput](this_.inputsById.Put_2(j2gort.PtrToAny(input.GetInputId()), j2gort.PtrToAny(input)))
}

// Java method: Compiler.getChunkGraph()
func (this_ *Compiler) GetChunkGraph() *JSChunkGraph {
	return this_.chunkGraph
}

// Java method: Compiler.getChunks()
func (this_ *Compiler) GetChunks() j2gort.Ext_java_lang_Iterable {
	if this_.chunkGraph != nil {
		return this_.chunkGraph.GetAllChunks()
	} else {
		return nil
	}
}

// Java method: Compiler.clearJSTypeRegistry()
func (this_ *Compiler) ClearJSTypeRegistry() {
	this_.typeRegistry = nil
	this_.typeValidator = nil
	this_.abstractInterpreter = nil
}

// Java method: Compiler.isTypeRegistryCleared()
func (this_ *Compiler) IsTypeRegistryCleared() bool {
	return (this_.typeCheckingHasRun && (this_.typeRegistry == nil))
}

// Java method: Compiler.getTypeRegistry()
func (this_ *Compiler) GetTypeRegistry() *pkg_rhino.JSTypeRegistry {
	if this_.typeRegistry == nil {
		if !(!(this_.HasTypeCheckingRun())) {
			j2gort.CheckStateStr(false, "Attempted to re-initialize JSTypeRegistry after it had been cleared")
		}
		this_.typeRegistry = pkg_rhino.NewJSTypeRegistry_2(this_.oldErrorReporter, this_.forwardDeclaredTypes)
	}
	return this_.typeRegistry
}

// Java method: Compiler.getColorRegistry()
func (this_ *Compiler) GetColorRegistry() *pkg_colors.ColorRegistry {
	return j2gort.CheckNotNullPtr(this_.colorRegistry)
}

// Java method: Compiler.setColorRegistry(com.google.javascript.jscomp.colors.ColorRegistry colorRegistry)
func (this_ *Compiler) SetColorRegistry(colorRegistry *pkg_colors.ColorRegistry) {
	_ = colorRegistry
	if !(this_.runtimeLibraryTypedAsts != nil) {
		j2gort.CheckState(false)
	}
	this_.colorRegistry = colorRegistry
}

// Java method: Compiler.forwardDeclareType(java.lang.String typeName)
func (this_ *Compiler) ForwardDeclareType(typeName string) {
	_ = typeName
	j2gort.SetAddStr(this_.forwardDeclaredTypes, typeName)
}

// Java method: Compiler.setTypeCheckingHasRun(boolean hasRun)
func (this_ *Compiler) SetTypeCheckingHasRun(hasRun bool) {
	_ = hasRun
	this_.typeCheckingHasRun = hasRun
}

// Java method: Compiler.hasTypeCheckingRun()
func (this_ *Compiler) HasTypeCheckingRun() bool {
	return this_.typeCheckingHasRun
}

// Java method: Compiler.hasOptimizationColors()
func (this_ *Compiler) HasOptimizationColors() bool {
	return (this_.colorRegistry != nil)
}

// Java method: Compiler.getTypedScopeCreator()
func (this_ *Compiler) GetTypedScopeCreator() ScopeCreator {
	if this_.typedScopeCreator == nil {
		if !(!(this_.HasTypeCheckingRun())) {
			j2gort.CheckStateStr(false, "Attempted to re-initialize TypedScopeCreator after it had been cleared")
		}
		this_.typedScopeCreator = NewTypedScopeCreator_1(this_)
	}
	return ToIface_ScopeCreator(this_.typedScopeCreator)
}

// Java method: Compiler.clearTypedScopeCreator()
func (this_ *Compiler) ClearTypedScopeCreator() {
	this_.typedScopeCreator = nil
}

// Java method: Compiler.ensureDefaultPassConfig()
func (this_ *Compiler) EnsureDefaultPassConfig() *DefaultPassConfig {
	var passes PassConfig = this_.GetPassConfig().GetBasePassConfig()
	_ = passes
	if !(j2gort.InstanceOf[*DefaultPassConfig](passes)) {
		j2gort.CheckStateStr(false, "PassConfigs must eventually delegate to the DefaultPassConfig")
	}
	return j2gort.CastPtr[DefaultPassConfig](passes)
}

// Java method: Compiler.buildKnownSymbolTable()
func (this_ *Compiler) BuildKnownSymbolTable() *SymbolTable {
	clinit_SymbolTable()
	var val_79 SymbolTable
	InitSymbolTable(&val_79, this_, this_.GetTypeRegistry())
	var symbolTable *SymbolTable = &val_79
	_ = symbolTable
	if this_.typedScopeCreator != nil {
		symbolTable.AddScopes(this_.typedScopeCreator.GetAllMemoizedScopes())
		symbolTable.AddSymbolsFrom(pkg_rhino.ToIface_StaticSymbolTable(this_.typedScopeCreator))
	} else {
		symbolTable.FindScopes(this_.externsRoot, this_.jsRoot)
	}
	var val_80 GlobalNamespace
	InitGlobalNamespace_3_AbstractCompiler_Node_Node(&val_80, this_, this_.externsRoot, this_.jsRoot)
	var globalNamespace *GlobalNamespace = &val_80
	_ = globalNamespace
	symbolTable.AddSymbolsFrom(globalNamespace)
	clinit_ReferenceCollector()
	var val_81 ReferenceCollector
	InitReferenceCollector_3(&val_81, this_, StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR(), NewSyntacticScopeCreator_1(this_))
	var refCollector *ReferenceCollector = &val_81
	_ = refCollector
	refCollector.Process_1(this_.GetRoot())
	symbolTable.AddSymbolsFrom(refCollector)
	var preprocessorSymbolTable *PreprocessorSymbolTable = this_.EnsureDefaultPassConfig().GetPreprocessorSymbolTable()
	_ = preprocessorSymbolTable
	if preprocessorSymbolTable != nil {
		symbolTable.AddSymbolsFrom(pkg_rhino.ToIface_StaticSymbolTable(preprocessorSymbolTable))
	}
	symbolTable.FlattenGoogModuleExports()
	symbolTable.FillNamespaceReferences()
	symbolTable.FillPropertyScopes()
	symbolTable.FillThisReferences(this_.externsRoot, this_.jsRoot)
	symbolTable.FillPropertySymbols(this_.externsRoot, this_.jsRoot)
	symbolTable.FillSuperReferences(this_.externsRoot, this_.jsRoot)
	symbolTable.FillJSDocInfo(this_.externsRoot, this_.jsRoot)
	symbolTable.FillSymbolVisibility(this_.externsRoot, this_.jsRoot)
	symbolTable.FillGoogProvideModuleRequires(this_.externsRoot, this_.jsRoot)
	symbolTable.RemoveGeneratedSymbols()
	return symbolTable
}

// Java method: Compiler.getTopScope()
func (this_ *Compiler) GetTopScope() *TypedScope {
	return this_.topScope
}

// Java method: Compiler.setTopScope(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope x)
func (this_ *Compiler) SetTopScope(x *TypedScope) {
	_ = x
	if !((x == nil) || (j2gort.CastPtr[TypedScope](x.GetParent()) == nil)) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(x))
	}
	this_.topScope = x
}

// Java method: Compiler.getReverseAbstractInterpreter()
func (this_ *Compiler) GetReverseAbstractInterpreter() pkg_jscomptype.ReverseAbstractInterpreter {
	if this_.abstractInterpreter == nil || this_.abstractInterpreter.IsJavaNil_() {
		var interpreter pkg_jscomptype.ChainableReverseAbstractInterpreter = pkg_jscomptype.NewSemanticReverseAbstractInterpreter(this_.GetTypeRegistry())
		_ = interpreter
		if this_.options.GetClosurePass() {
			interpreter = pkg_jscomptype.NewClosureReverseAbstractInterpreter(this_.GetTypeRegistry()).Append(interpreter).GetFirst()
		}
		this_.abstractInterpreter = j2gort.AsIface[pkg_jscomptype.ReverseAbstractInterpreter](interpreter)
	}
	return this_.abstractInterpreter
}

// Java method: Compiler.getTypeValidator()
func (this_ *Compiler) GetTypeValidator() *TypeValidator {
	if this_.typeValidator == nil {
		if !(!(this_.HasTypeCheckingRun())) {
			j2gort.CheckStateStr(false, "Attempted to re-initialize TypeValidator after it had been cleared")
		}
		this_.typeValidator = NewTypeValidator(this_)
	}
	return this_.typeValidator
}

// Java method: Compiler.getTypeMismatches()
func (this_ *Compiler) GetTypeMismatches() j2gort.Ext_java_lang_Iterable {
	if this_.typeCheckingHasRun {
		return this_.GetTypeValidator().GetMismatches()
	}
	panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Can't ask for type mismatches before type checking.")))
}

// Java method: Compiler.getTranspilationNamespace()
func (this_ *Compiler) GetTranspilationNamespace() pkg_rhino.StaticScope {
	if this_.transpilationNamespace == nil || this_.transpilationNamespace.IsJavaNil_() {
		var val_82 GlobalNamespace
		InitGlobalNamespace_3_AbstractCompiler_Node_Node(&val_82, this_, this_.GetExternsRoot(), this_.GetJsRoot())
		var gn *GlobalNamespace = &val_82
		_ = gn
		gn.SetShouldTraverseScriptPredicate(j2gort.NewNodeLambda_Ext_java_util_function_Predicate(func(rawNode_83 uint32) bool {
			script := pkg_rhino.Node(rawNode_83)
			_ = script
			return (script.IsFromExterns() || script.IsFirstChildOf(this_.GetJsRoot()))
		}))
		gn.EnsureGenerated()
		this_.transpilationNamespace = gn
	}
	return this_.transpilationNamespace
}

// Java method: Compiler.maybeSetTracker()
func (this_ *Compiler) MaybeSetTracker() {
	if !(this_.options.GetTracerMode().IsOn()) || (this_.tracker != nil) {
		return
	}
	this_.tracker = NewPerformanceTracker(this_.externsRoot, this_.jsRoot, this_.options.GetTracerMode())
	this_.changeTracker.AddChangeHandler(this_.tracker.GetCodeChangeHandler())
}

// Java method: Compiler.initializeModuleLoader()
func (this_ *Compiler) InitializeModuleLoader() {
	var moduleResolverFactory ModuleLoader_ModuleResolverFactory = nil
	_ = moduleResolverFactory
	switch (this_.options.GetModuleResolutionMode()).Ordinal() {
	case int32(0) /* ResolutionMode.BROWSER */ :
		moduleResolverFactory = StaticGet_BrowserModuleResolver_FACTORY()
	case int32(2) /* ResolutionMode.NODE */ :
		moduleResolverFactory = NewNodeModuleResolver_Factory_1(this_.ProcessJsonInputs(this_.chunkGraph.GetAllInputs()))
	case int32(3) /* ResolutionMode.WEBPACK */ :
		moduleResolverFactory = NewWebpackModuleResolver_Factory(j2gort.AsIface[j2gort.Ext_java_util_Map](this_.inputPathByWebpackId))
	case int32(1) /* ResolutionMode.BROWSER_WITH_TRANSFORMED_PREFIXES */ :
		moduleResolverFactory = NewBrowserWithTransformedPrefixesModuleResolver_Factory(this_.options.GetBrowserResolverPrefixReplacements())
	}
	this_.moduleLoader = ModuleLoader_S_Builder().SetModuleRoots(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.options.GetModuleRoots())).SetInputs(this_.chunkGraph.GetAllInputs()).SetFactory(moduleResolverFactory).SetPathResolver(StaticGet_ModuleLoader_PathResolver_RELATIVE()).SetPathEscaper(this_.options.GetPathEscaper()).Build()
}

// Java method: Compiler.parseInputs()
func (this_ *Compiler) ParseInputs() pkg_rhino.Node {
	var devMode bool = (this_.options.GetDevMode() != StaticGet_CompilerOptions_DevMode_OFF())
	_ = devMode
	this_.externsRoot.DetachChildren()
	this_.jsRoot.DetachChildren()
	this_.scriptNodeByFilename.Clear()
	var tracer *Tracer = this_.NewTracer("parseInputs")
	_ = tracer
	this_.BeforePass("parseInputs")
	defer func() {
		this_.AfterPass("parseInputs")
		this_.StopTracer(tracer, "parseInputs")
	}()
	{
		if this_.options.GetNumParallelThreads() > int32(1) {
			{
				var val_84 PrebuildAst
				InitPrebuildAst(&val_84, this_, this_.options.GetNumParallelThreads())
				(&val_84).Prebuild(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.externs))
			}
		}
		{
			cur_85 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.externs))
			for cur_85.HasNext() {
				rawElem_86 := cur_85.Next()
				input := j2gort.CastPtr[CompilerInput](rawElem_86)
				_ = input
				var n pkg_rhino.Node = j2gort.CheckNotNullNode(input.GetAstRoot(this_))
				_ = n
				if this_.HasErrors() {
					return pkg_rhino.Node(0)
				}
				this_.externsRoot.AddChildToBack(n)
				j2gort.MapPutStr(this_.scriptNodeByFilename, input.GetSourceFile().GetName(), j2gort.BoxNode(n))
			}
		}
		if this_.options.GetLanguageIn().ToFeatureSet().Has(pkg_parser.StaticGet_FeatureSet_Feature_MODULES()) || this_.options.GetProcessCommonJSModules() {
			this_.InitializeModuleLoader()
		} else {
			this_.moduleLoader = StaticGet_ModuleLoader_EMPTY()
		}
		if this_.options.GetDependencyOptions().NeedsManagement() {
			this_.priv_Compiler_findModulesFromEntryPoints_1504(this_.options.GetLanguageIn().ToFeatureSet().Has(pkg_parser.StaticGet_FeatureSet_Feature_MODULES()), this_.options.GetProcessCommonJSModules())
		} else if this_.options.NeedsTranspilationFrom(pkg_parser.StaticGet_FeatureSet_ES2015_MODULES()) || this_.options.GetProcessCommonJSModules() {
			if this_.options.GetLanguageIn().ToFeatureSet().Has(pkg_parser.StaticGet_FeatureSet_Feature_MODULES()) {
				this_.priv_Compiler_parsePotentialModules_1514(this_.chunkGraph.GetAllInputs())
			}
			var inputModuleIdentifiers j2gort.OrderedMap[string, *CompilerInput] = j2gort.NewOrderedMap[string, *CompilerInput]()
			_ = inputModuleIdentifiers
			{
				cur_87 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.chunkGraph.GetAllInputs()))
				for cur_87.HasNext() {
					rawElem_88 := cur_87.Next()
					input_2 := j2gort.CastPtr[CompilerInput](rawElem_88)
					_ = input_2
					if input_2.GetKnownProvides().IsEmpty() {
						var modPath *ModuleLoader_ModulePath = this_.moduleLoader.Resolve(input_2.GetSourceFile().GetName())
						_ = modPath
						inputModuleIdentifiers.Put(modPath.ToModuleName(), input_2)
					}
				}
			}
			var inputsToRewrite j2gort.OrderedMap[string, *CompilerInput] = j2gort.NewOrderedMap[string, *CompilerInput]()
			_ = inputsToRewrite
			{
				cur_89 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.chunkGraph.GetAllInputs()))
				for cur_89.HasNext() {
					rawElem_90 := cur_89.Next()
					input_3 := j2gort.CastPtr[CompilerInput](rawElem_90)
					_ = input_3
					{
						cur_91 := j2gort.NewIterCursor(j2gort.CheckNotNull(input_3.GetKnownRequiredSymbols()))
						for cur_91.HasNext() {
							require := cur_91.NextStr()
							_ = require
							if inputModuleIdentifiers.ContainsKey(require) && !(inputsToRewrite.ContainsKey(require)) {
								inputsToRewrite.Put(require, inputModuleIdentifiers.Get(require))
							}
						}
					}
				}
			}
			{
				cur_92 := inputsToRewrite.IterValues()
				for cur_92.HasNext() {

					input_4 := cur_92.
						Next()
					_ = input_4
					input_4.SetJsModuleType(StaticGet_CompilerInput_ModuleType_IMPORTED_SCRIPT())
				}
			}
		}
		if this_.moduleLoader != nil {
			this_.moduleLoader.SetErrorHandler(this_)
		}
		this_.OrderInputs()
		this_.priv_Compiler_markClosureUnawareCode_1508()
		if this_.HasErrors() {
			return pkg_rhino.Node(0)
		}
		if this_.options.GetNumParallelThreads() > int32(1) {
			{
				var val_94 PrebuildAst
				InitPrebuildAst(&val_94, this_, this_.options.GetNumParallelThreads())
				(&val_94).Prebuild(this_.chunkGraph.GetAllInputs())
			}
		}
		{
			cur_95 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.chunkGraph.GetAllInputs()))
			for cur_95.HasNext() {
				rawElem_96 := cur_95.Next()
				input_5 := j2gort.CastPtr[CompilerInput](rawElem_96)
				_ = input_5
				var n_2 pkg_rhino.Node = j2gort.CheckNotNullNode(input_5.GetAstRoot(this_))
				_ = n_2
				if devMode {
					this_.priv_Compiler_runValidityCheck_1483()
					if this_.HasErrors() {
						return pkg_rhino.Node(0)
					}
				}
				if (this_.options.ShouldGatherSourceMapInfo() || (this_.options.GetExternExportsPath() != j2gort.NullString)) || !(this_.options.GetReplaceStringsFunctionDescriptions().IsEmpty()) {
					var sia *SourceInformationAnnotator
					_ = sia
					if this_.options.GetDevMode() == StaticGet_CompilerOptions_DevMode_OFF() {
						sia = SourceInformationAnnotator_Create()
					} else {
						sia = SourceInformationAnnotator_CreateWithAnnotationChecks(input_5.GetName())
					}
					NodeTraversal_Traverse(this_, n_2, ToIface_NodeTraversal_Callback(sia))
				}
				if NodeUtil_IsFromTypeSummary(n_2) {
					input_5.SetIsExtern()
					this_.externsRoot.AddChildToBack(n_2)
				} else {
					this_.jsRoot.AddChildToBack(n_2)
				}
				j2gort.MapPutStr(this_.scriptNodeByFilename, input_5.GetSourceFile().GetName(), j2gort.BoxNode(n_2))
			}
		}
		if this_.HasErrors() {
			return pkg_rhino.Node(0)
		}
		if this_.typedAstFilesystem != nil {
			j2gort.CheckStateMsg(((this_.typedAstFilesystem.Size() == int32(1)) || this_.typedAstFilesystem.IsEmpty()), j2gort.ValueOfInt(this_.typedAstFilesystem.Size()))
		}
		return this_.externAndJsRoot
	}
	panic("unreachable")
}

// Java method: Compiler.preferRegexParser()
func (this_ *Compiler) PreferRegexParser() bool {
	return this_.preferRegexParser
}

// Java method: Compiler.setPreferRegexParser(boolean preferRegexParser)
func (this_ *Compiler) SetPreferRegexParser(preferRegexParser bool) {
	_ = preferRegexParser
	this_.preferRegexParser = preferRegexParser
}

// Java method: Compiler.orderInputsWithLargeStack()
func (this_ *Compiler) OrderInputsWithLargeStack() {
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		var tracer *Tracer = this_.NewTracer("orderInputsWithLargeStack")
		_ = tracer
		func() {
			defer func() {
				this_.StopTracer(tracer, "orderInputsWithLargeStack")
			}()
			this_.OrderInputs()
		}()
		return nil
	}))
}

// Java method: Compiler.orderInputs()
func (this_ *Compiler) OrderInputs() {
	this_.priv_Compiler_maybeDoThreadedParsing_1511()
	var originalInputs j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(this_.chunkGraph.GetAllInputs())
	_ = originalInputs
	if this_.tracker != nil {
		this_.tracker.RecordPrePruningInputCount(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](originalInputs), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.externs))
	}
	this_.priv_Compiler_markExterns_1509(originalInputs)
	var staleInputs bool = false
	_ = staleInputs
	if this_.options.GetDependencyOptions().NeedsManagement() {
		func() {
			defer func() {
				if rec_97 := recover(); rec_97 != nil {
					jex_98 := javalang.RecoverException(rec_97)
					_ = jex_98
					if j2gort.InstanceOf[*SortedDependencies_MissingProvideException](jex_98) {
						e := j2gort.AsIface[*SortedDependencies_MissingProvideException](jex_98)
						_ = e
						this_.Report_1(JSError_Make_2(StaticGet_Compiler_MISSING_ENTRY_ERROR(), []string{j2gort.ThrowableGetMessage(e)}))
					} else if j2gort.InstanceOf[*JSChunkGraph_MissingChunkException](jex_98) {
						e_2 := j2gort.AsIface[*JSChunkGraph_MissingChunkException](jex_98)
						_ = e_2
						this_.Report_1(JSError_Make_2(StaticGet_Compiler_MISSING_CHUNK_ERROR(), []string{j2gort.ThrowableGetMessage(e_2)}))
					} else {
						panic(rec_97)
					}
				}
			}()
			var result *JSChunkGraph_DependencyManagementResult = this_.chunkGraph.ManageDependencies(this_, this_.options.GetDependencyOptions())
			_ = result
			staleInputs = true
			this_.priv_Compiler_maybeLogPruningAnalysis_1569(result)
		}()
	}
	this_.HoistExterns(originalInputs)
	this_.priv_Compiler_fillEmptyChunks_1452(this_.GetChunks())
	this_.HoistNoCompileFiles()
	if staleInputs {
		this_.priv_Compiler_repartitionInputs_1512()
	}
	this_.priv_Compiler_maybeLogPrunedInputs_1441()
}

// Java method: Compiler.findModulesFromEntryPoints(boolean supportEs6Modules, boolean supportCommonJSModules)
func (this_ *Compiler) priv_Compiler_findModulesFromEntryPoints_1504(supportEs6Modules bool, supportCommonJSModules bool) {
	_ = supportEs6Modules
	_ = supportCommonJSModules
	this_.priv_Compiler_maybeDoThreadedParsing_1511()
	var entryPoints []*CompilerInput = []*CompilerInput{}
	_ = entryPoints
	var inputsByProvide j2gort.OrderedMap[string, *CompilerInput] = j2gort.NewOrderedMap[string, *CompilerInput]()
	_ = inputsByProvide
	var inputsByIdentifier j2gort.OrderedMap[string, *CompilerInput] = j2gort.NewOrderedMap[string, *CompilerInput]()
	_ = inputsByIdentifier
	var keepMoochers bool = !(this_.options.GetDependencyOptions().ShouldDropMoochers())
	_ = keepMoochers
	{
		cur_99 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.chunkGraph.GetAllInputs()))
		for cur_99.HasNext() {
			rawElem_100 := cur_99.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_100)
			_ = input
			var hasNonModuleProvide bool = false
			_ = hasNonModuleProvide
			{
				cur_101 := j2gort.NewIterCursor(j2gort.CheckNotNull(input.GetProvides()))
				for cur_101.HasNext() {
					provide := cur_101.NextStr()
					_ = provide
					if !(javalang.StringStartsWith(j2gort.CheckNotNullString(provide), "module$")) {
						hasNonModuleProvide = true
						inputsByProvide.Put(provide, input)
					}
				}
			}
			if keepMoochers && !(hasNonModuleProvide) {
				entryPoints = append(entryPoints, input)
			}
			inputsByIdentifier.Put(ModuleIdentifier_ForFile(input.GetPath().String()).String(), input)
		}
	}
	{

		for _, moduleIdentifier := range j2gort.CheckNotNullSlice(j2gort.CheckNotNullPtr(this_.options.GetDependencyOptions()).
			Slice_EntryPoints()) {

			_ = moduleIdentifier
			var input_2 *CompilerInput = inputsByProvide.Get(moduleIdentifier.String())
			_ = input_2
			if input_2 == nil {
				input_2 = inputsByIdentifier.Get(moduleIdentifier.String())
			}
			if input_2 != nil {
				entryPoints = append(entryPoints, input_2)
			}
		}
	}
	var workingInputSet j2gort.OrderedSet[*CompilerInput] = j2gort.NewOrderedSet[*CompilerInput]()
	_ = workingInputSet
	{
		cur_104 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.chunkGraph.GetAllInputs()))
		for cur_104.HasNext() {
			rawElem_105 := cur_104.Next()
			input_3 := j2gort.CastPtr[CompilerInput](rawElem_105)
			_ = input_3
			workingInputSet.Add(input_3)
		}
	}
	for _, entryPoint := range entryPoints {
		_ = entryPoint
		this_.priv_Compiler_findModulesFromInput_1505(entryPoint, false, &workingInputSet, &inputsByIdentifier, &inputsByProvide, supportEs6Modules, supportCommonJSModules)
	}

}

// Java method: Compiler.findModulesFromInput(com.google.javascript.jscomp.CompilerInput input, boolean wasImportedByModule, java.util.Set<com.google.javascript.jscomp.CompilerInput> inputs, java.util.Map<java.lang.String,com.google.javascript.jscomp.CompilerInput> inputsByIdentifier, java.util.Map<java.lang.String,com.google.javascript.jscomp.CompilerInput> inputsByProvide, boolean supportEs6Modules, boolean supportCommonJSModules)
func (this_ *Compiler) priv_Compiler_findModulesFromInput_1505(input *CompilerInput, wasImportedByModule bool, inputs *j2gort.OrderedSet[*CompilerInput], inputsByIdentifier *j2gort.OrderedMap[string, *CompilerInput], inputsByProvide *j2gort.OrderedMap[string, *CompilerInput], supportEs6Modules bool, supportCommonJSModules bool) {
	_ = input
	_ = wasImportedByModule
	_ = inputs
	_ = inputsByIdentifier
	_ = inputsByProvide
	_ = supportEs6Modules
	_ = supportCommonJSModules
	if !(inputs.Remove(input)) {
		if wasImportedByModule && (input.GetJsModuleType() == StaticGet_CompilerInput_ModuleType_NONE()) {
			input.SetJsModuleType(StaticGet_CompilerInput_ModuleType_IMPORTED_SCRIPT())
		}
		return
	}
	var val_108 FindModuleDependencies
	InitFindModuleDependencies(&val_108, this_, supportEs6Modules, supportCommonJSModules, this_.inputPathByWebpackId)
	var findDeps *FindModuleDependencies = &val_108
	_ = findDeps
	findDeps.Process(j2gort.CheckNotNullNode(input.GetAstRoot(this_)))
	if wasImportedByModule && (input.GetJsModuleType() == StaticGet_CompilerInput_ModuleType_NONE()) {
		input.SetJsModuleType(StaticGet_CompilerInput_ModuleType_IMPORTED_SCRIPT())
	}
	this_.moduleTypesByName.Put(input.GetPath().ToModuleName(), input.GetJsModuleType())
	var allDeps j2gort.Ext_java_lang_Iterable = j2gort.Ext_com_google_common_collect_Iterables_Concat_3(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](input.GetRequiredSymbols()), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](input.GetDynamicRequires()), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](input.GetTypeRequires()))
	_ = allDeps
	{
		cur_109 := j2gort.NewIterCursor(j2gort.CheckNotNull(allDeps))
		for cur_109.HasNext() {
			requiredNamespace := cur_109.NextStr()
			_ = requiredNamespace
			var requiredInput *CompilerInput = nil
			_ = requiredInput
			var requiredByModuleImport bool = false
			_ = requiredByModuleImport
			if inputsByProvide.ContainsKey(requiredNamespace) {
				requiredInput = inputsByProvide.Get(requiredNamespace)
			} else if inputsByIdentifier.ContainsKey(requiredNamespace) {
				requiredByModuleImport = true
				requiredInput = inputsByIdentifier.Get(requiredNamespace)
			}
			if requiredInput != nil {
				this_.priv_Compiler_findModulesFromInput_1505(requiredInput, requiredByModuleImport, inputs, inputsByIdentifier, inputsByProvide, supportEs6Modules, supportCommonJSModules)
			}
		}
	}
}

// Java method: Compiler.hoistExterns(com.google.common.collect.ImmutableList<com.google.javascript.jscomp.CompilerInput> originalInputs)
func (this_ *Compiler) HoistExterns(originalInputs j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = originalInputs
	var staleInputs bool = false
	_ = staleInputs
	{
		cur_110 := j2gort.NewIterCursor(j2gort.CheckNotNull(originalInputs))
		for cur_110.HasNext() {
			rawElem_111 := cur_110.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_111)
			_ = input
			if this_.priv_Compiler_hoistIfExtern_1507(input) {
				staleInputs = true
			}
		}
	}
	if staleInputs {
		this_.priv_Compiler_repartitionInputs_1512()
	}
}

// Java method: Compiler.hoistIfExtern(com.google.javascript.jscomp.CompilerInput input)
func (this_ *Compiler) priv_Compiler_hoistIfExtern_1507(input *CompilerInput) bool {
	_ = input
	if input.GetHasExternsAnnotation() {
		var root pkg_rhino.Node = input.GetAstRoot(this_)
		_ = root
		this_.externsRoot.AddChildToBack(root)
		j2gort.MapPutStr(this_.scriptNodeByFilename, input.GetSourceFile().GetName(), j2gort.BoxNode(root))
		var chunk *JSChunk = input.GetChunk()
		_ = chunk
		if chunk != nil {
			this_.hoistedExternsByChunk.Put(j2gort.PtrToAny(chunk), j2gort.PtrToAny(input))
			chunk.Remove(input)
		}
		this_.externs.Add_1(j2gort.PtrToAny(input))
		return true
	}
	return false
}

// Java method: Compiler.markClosureUnawareCode()
func (this_ *Compiler) priv_Compiler_markClosureUnawareCode_1508() {
	{
		cur_112 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.chunkGraph.GetAllInputs()))
		for cur_112.HasNext() {
			rawElem_113 := cur_112.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_113)
			_ = input
			var root pkg_rhino.Node = input.GetAstRoot(this_)
			_ = root
			var info *pkg_rhino.JSDocInfo = root.GetJSDocInfo()
			_ = info
			if (info != nil) && info.IsClosureUnawareCode() {
				input.GetSourceFile().MarkAsClosureUnawareCode()
			}
		}
	}
}

// Java method: Compiler.markExterns(com.google.common.collect.ImmutableList<com.google.javascript.jscomp.CompilerInput> originalInputs)
func (this_ *Compiler) priv_Compiler_markExterns_1509(originalInputs j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = originalInputs
	{
		cur_114 := j2gort.NewIterCursor(j2gort.CheckNotNull(originalInputs))
		for cur_114.HasNext() {
			rawElem_115 := cur_114.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_115)
			_ = input
			if input.GetHasExternsAnnotation() {
				input.SetIsExtern()
			}
		}
	}
}

// Java method: Compiler.hoistNoCompileFiles()
func (this_ *Compiler) HoistNoCompileFiles() {
	var staleInputs bool = false
	_ = staleInputs
	this_.priv_Compiler_maybeDoThreadedParsing_1511()
	{
		cur_116 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(this_.chunkGraph.GetAllInputs())))
		for cur_116.HasNext() {
			rawElem_117 := cur_116.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_117)
			_ = input
			if input.GetHasNoCompileAnnotation() {
				input.GetChunk().Remove(input)
				staleInputs = true
			}
		}
	}
	if staleInputs {
		this_.priv_Compiler_repartitionInputs_1512()
	}
}

// Java method: Compiler.maybeDoThreadedParsing()
func (this_ *Compiler) priv_Compiler_maybeDoThreadedParsing_1511() {
	if this_.options.GetNumParallelThreads() > int32(1) {
		{
			var val_118 PrebuildDependencyInfo
			InitPrebuildDependencyInfo(&val_118, this_.options.GetNumParallelThreads())
			(&val_118).Prebuild(this_.chunkGraph.GetAllInputs())
		}
	}
}

// Java method: Compiler.repartitionInputs()
func (this_ *Compiler) priv_Compiler_repartitionInputs_1512() {
	this_.priv_Compiler_fillEmptyChunks_1452(this_.GetChunks())
	this_.RebuildInputsFromChunks()
}

// Java method: Compiler.processJsonInputs(java.lang.Iterable<com.google.javascript.jscomp.CompilerInput> inputsToProcess)
func (this_ *Compiler) ProcessJsonInputs(inputsToProcess j2gort.Ext_java_lang_Iterable) j2gort.Ext_java_util_Map {
	_ = inputsToProcess
	var rewriteJson *RewriteJsonToModule = NewRewriteJsonToModule(this_)
	_ = rewriteJson
	{
		cur_119 := j2gort.NewIterCursor(j2gort.CheckNotNull(inputsToProcess))
		for cur_119.HasNext() {
			rawElem_120 := cur_119.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_120)
			_ = input
			if !(javalang.StringEndsWith(j2gort.CheckNotNullString(input.GetSourceFile().GetName()), ".json")) {
				continue
			}
			input.SetCompiler(this_)
			var tryFlow_121 int
			func() {
				defer func() {
					if rec_122 := recover(); rec_122 != nil {
						jex_123 := javalang.RecoverException(rec_122)
						_ = jex_123
						if j2gort.InstanceOf[javaio.IIOException](jex_123) {
							e := j2gort.AsIface[javaio.IIOException](jex_123)
							_ = e
							tryFlow_121 = 2
							return
						} else {
							panic(rec_122)
						}
					}
				}()
				input.GetSourceFile().SetCodeDeprecated(("(" + j2gort.StrValString(input.GetSourceFile().GetCode()) + ")"))
			}()
			if tryFlow_121 == 2 {
				continue
			}
			var root pkg_rhino.Node = j2gort.CheckNotNullNode(input.GetAstRoot(this_))
			_ = root
			input.SetJsModuleType(StaticGet_CompilerInput_ModuleType_JSON())
			rewriteJson.Process_2(pkg_rhino.Node(0), root)
		}
	}
	return j2gort.AsIface[j2gort.Ext_java_util_Map](rewriteJson.GetPackageJsonMainEntries())
}

// Java method: Compiler.parsePotentialModules(java.lang.Iterable<com.google.javascript.jscomp.CompilerInput> inputsToProcess)
func (this_ *Compiler) priv_Compiler_parsePotentialModules_1514(inputsToProcess j2gort.Ext_java_lang_Iterable) {
	_ = inputsToProcess
	var filteredInputs j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = filteredInputs
	{
		cur_124 := j2gort.NewIterCursor(j2gort.CheckNotNull(inputsToProcess))
		for cur_124.HasNext() {
			rawElem_125 := cur_124.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_125)
			_ = input
			if (!(this_.options.GetDependencyOptions().ShouldPrune()) || !(JsFileRegexParser_IsSupported())) || input.IsEs6Module() {
				filteredInputs.Add_1(j2gort.PtrToAny(input))
			}
		}
	}
	if this_.options.GetNumParallelThreads() > int32(1) {
		{
			var val_126 PrebuildAst
			InitPrebuildAst(&val_126, this_, this_.options.GetNumParallelThreads())
			(&val_126).Prebuild(filteredInputs)
		}
	}
	{
		cur_127 := j2gort.NewIterCursor(filteredInputs)
		for cur_127.HasNext() {
			rawElem_128 := cur_127.Next()
			input_2 := j2gort.CastPtr[CompilerInput](rawElem_128)
			_ = input_2
			input_2.SetCompiler(this_)
			input_2.GetRequires()
			input_2.SetJsModuleType(StaticGet_CompilerInput_ModuleType_ES6())
		}
	}
}

// Java method: Compiler.newCompilerOptions()
func (this_ *Compiler) NewCompilerOptions() *CompilerOptions {
	return NewCompilerOptions()
}

// Java method: Compiler.initCompilerOptionsIfTesting()
func (this_ *Compiler) InitCompilerOptionsIfTesting() {
	if this_.options == nil {
		this_.InitOptions(this_.NewCompilerOptions())
		this_.options.SetLanguage(StaticGet_CompilerOptions_LanguageMode_UNSUPPORTED())
	}
}

// Java method: Compiler.parseSyntheticCode(java.lang.String fileName, java.lang.String js)
func (this_ *Compiler) ParseSyntheticCode(fileName string, js string) pkg_rhino.Node {
	_ = fileName
	_ = js
	this_.InitCompilerOptionsIfTesting()
	var source *SourceFile = SourceFile_FromCode_2((" [synthetic:" + j2gort.StrValString(fileName) + "] "), js)
	_ = source
	this_.priv_Compiler_addFilesToSourceMap_1558(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.ImmListOfPtr(source)))
	return this_.priv_Compiler_parseCodeHelper_1517(source)
}

// Java method: Compiler.parseTestCode(java.lang.String js)
func (this_ *Compiler) ParseTestCode_1_string(js string) pkg_rhino.Node {
	_ = js
	this_.InitCompilerOptionsIfTesting()
	this_.InitBasedOnOptions()
	return this_.priv_Compiler_parseCodeHelper_1517(SourceFile_FromCode_2("testcode", js))
}

// Java method: Compiler.parseTestCode(com.google.common.collect.ImmutableList<java.lang.String> code)
func (this_ *Compiler) ParseTestCode_1_Ext_com_google_common_collect_ImmutableList(code j2gort.Ext_com_google_common_collect_ImmutableList) pkg_rhino.Node {
	_ = code
	this_.InitCompilerOptionsIfTesting()
	this_.InitBasedOnOptions()
	return this_.priv_Compiler_parseCodeHelper_1518(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_Streams_MapWithIndex_2_Ext_java_util_stream_Stream_Ext_com_google_common_collect_Streams_FunctionWithIndex(code.Stream(), lambdaSingleton_compiler_4).Collect_1(j2gort.Ext_java_util_stream_Collectors_ToList())))
}

// Java method: Compiler.parseCodeHelper(com.google.javascript.jscomp.SourceFile src)
func (this_ *Compiler) priv_Compiler_parseCodeHelper_1517(src *SourceFile) pkg_rhino.Node {
	_ = src
	var val_130 CompilerInput
	InitCompilerInput_1(&val_130, src)
	var input *CompilerInput = &val_130
	_ = input
	this_.priv_Compiler_putCompilerInput_1494(input)
	var root pkg_rhino.Node = input.GetAstRoot(this_)
	_ = root
	j2gort.MapPutStr(this_.scriptNodeByFilename, input.GetSourceFile().GetName(), j2gort.BoxNode(root))
	return j2gort.CheckNotNullNode(root)
}

// Java method: Compiler.parseCodeHelper(java.util.List<com.google.javascript.jscomp.SourceFile> srcs)
func (this_ *Compiler) priv_Compiler_parseCodeHelper_1518(srcs j2gort.Ext_java_util_List) pkg_rhino.Node {
	_ = srcs
	var root pkg_rhino.Node = pkg_rhino.IR_Root([]pkg_rhino.Node{})
	_ = root
	{
		cur_131 := j2gort.NewIterCursor(j2gort.CheckNotNull(srcs))
		for cur_131.HasNext() {
			rawElem_132 := cur_131.Next()
			src := j2gort.CastPtr[SourceFile](rawElem_132)
			_ = src
			root.AddChildToBack(this_.priv_Compiler_parseCodeHelper_1517(src))
		}
	}
	return root
}

// Java method: Compiler.getDefaultErrorReporter()
func (this_ *Compiler) GetDefaultErrorReporter() pkg_rhino.ErrorReporter {
	return this_.oldErrorReporter
}

// Java method: Compiler.toSource()
func (this_ *Compiler) ToSource_0() string {
	return j2gort.CastString(this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		var tracer *Tracer = this_.NewTracer("toSource")
		_ = tracer
		defer func() {
			this_.StopTracer(tracer, "toSource")
		}()
		{
			var val_133 Compiler_CodeBuilder
			InitCompiler_CodeBuilder(&val_133)
			var cb *Compiler_CodeBuilder = &val_133
			_ = cb
			var lt *Compiler_SingleBinaryLicenseTracker = NewCompiler_SingleBinaryLicenseTracker(this_)
			_ = lt
			if this_.jsRoot != 0 {
				var i int32 = int32(0)
				_ = i
				if this_.options.ShouldPrintExterns() {
					{
						var scriptNode pkg_rhino.Node = this_.externsRoot.GetFirstChild()
						_ = scriptNode
						for ; scriptNode != 0; scriptNode = scriptNode.GetNext() {
							this_.ToSource_4(cb, lt, func() int32 { old_134 := i; nw_135 := old_134 + 1; i = nw_135; return old_134 }(), scriptNode)
						}
					}
				}
				{
					var scriptNode_2 pkg_rhino.Node = this_.jsRoot.GetFirstChild()
					_ = scriptNode_2
					for ; scriptNode_2 != 0; scriptNode_2 = scriptNode_2.GetNext() {
						this_.ToSource_4(cb, lt, func() int32 { old_136 := i; nw_137 := old_136 + 1; i = nw_137; return old_136 }(), scriptNode_2)
					}
				}
			}
			return j2gort.BoxString(cb.String())
		}
		panic("unreachable")
	})))
}

// Java method: Compiler.toSource(com.google.javascript.jscomp.JSChunk chunk)
func (this_ *Compiler) ToSource_1_JSChunk(chunk *JSChunk) string {
	_ = chunk
	return this_.ToSource_2(NewCompiler_ScriptNodeLicensesOnlyTracker(this_), chunk)
}

// Java method: Compiler.toSource(com.google.javascript.jscomp.CodePrinter.LicenseTracker licenseTracker, com.google.javascript.jscomp.JSChunk chunk)
func (this_ *Compiler) ToSource_2(licenseTracker CodePrinter_LicenseTracker, chunk *JSChunk) string {
	_ = licenseTracker
	_ = chunk
	return j2gort.CastString(this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		var inputs []*CompilerInput = j2gort.ImmListToSliceNullable[*CompilerInput](chunk.GetInputs())
		_ = inputs
		var chunkExterns j2gort.Ext_java_util_Set
		_ = chunkExterns
		if this_.options.ShouldPrintExterns() {
			chunkExterns = j2gort.AsIface[j2gort.Ext_java_util_Set](this_.hoistedExternsByChunk.Get(j2gort.PtrToAny(chunk)))
		} else {
			chunkExterns = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
		}
		if len(inputs) == 0 && chunkExterns.IsEmpty() {
			return ""
		}
		var val_138 Compiler_CodeBuilder
		InitCompiler_CodeBuilder(&val_138)
		var cb *Compiler_CodeBuilder = &val_138
		_ = cb
		var seq int32 = int32(0)
		_ = seq
		{
			cur_139 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunkExterns))
			for cur_139.HasNext() {
				rawElem_140 := cur_139.Next()
				externInput := j2gort.CastPtr[CompilerInput](rawElem_140)
				_ = externInput
				var scriptNode pkg_rhino.Node = externInput.GetAstRoot(this_)
				_ = scriptNode
				if scriptNode == 0 {
					panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Bad chunk: "+j2gort.StrValString(chunk.GetName()))))
				}
				this_.ToSource_4(cb, licenseTracker, func() int32 { old_141 := seq; nw_142 := old_141 + 1; seq = nw_142; return old_141 }(), scriptNode)
			}
		}
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(inputs)); i = i + 1 {
				if !(chunkExterns.IsEmpty()) && AbstractCompiler_IsFillFileName(inputs[i].GetName()) {
					continue
				}
				var scriptNode_2 pkg_rhino.Node = inputs[i].GetAstRoot(this_)
				_ = scriptNode_2
				if scriptNode_2 == 0 {
					panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Bad chunk: "+j2gort.StrValString(chunk.GetName()))))
				}
				this_.ToSource_4(cb, licenseTracker, func() int32 { old_143 := seq; nw_144 := old_143 + 1; seq = nw_144; return old_143 }(), scriptNode_2)
			}
		}
		return j2gort.BoxString(cb.String())
	})))
}

// Java method: Compiler.toSource(com.google.javascript.jscomp.Compiler.CodeBuilder cb, com.google.javascript.jscomp.CodePrinter.LicenseTracker licenseTracker, int inputSeqNum, com.google.javascript.rhino.Node root)
func (this_ *Compiler) ToSource_4(cb *Compiler_CodeBuilder, licenseTracker CodePrinter_LicenseTracker, inputSeqNum int32, root pkg_rhino.Node) {
	_ = cb
	_ = licenseTracker
	_ = inputSeqNum
	_ = root
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		if this_.options.ShouldPrintInputDelimiter() {
			if (cb.GetLength() > int32(0)) && !(cb.EndsWith("\n")) {
				cb.Append("\n")
			}
			if !(root.IsScript()) {
				j2gort.CheckState(false)
			}
			var delimiter string = this_.options.GetInputDelimiter()
			_ = delimiter
			var inputName string = root.GetInputId().GetIdName()
			_ = inputName
			var sourceName string = root.GetSourceFileName()
			_ = sourceName
			if !(sourceName != j2gort.NullString) {
				j2gort.CheckState(false)
			}
			if !(!(j2gort.CheckNotNullString(sourceName) == "")) {
				j2gort.CheckState(false)
			}
			delimiter = javalang.StringReplace(javalang.StringReplace(javalang.StringReplace(delimiter, "%name%", j2gort.CheckNotNullString(inputName)), "%num%", strconv.FormatInt(int64(inputSeqNum), 10)), "%n%", "\n")
			cb.Append(delimiter).Append("\n")
		}
		var sourceAndMappings *CodePrinter_SourceAndMappings = this_.priv_Compiler_toSourceAndMappings_1522(root, (inputSeqNum == int32(0)), licenseTracker)
		_ = sourceAndMappings
		var code string = sourceAndMappings.source
		_ = code
		{
			cur_145 := j2gort.NewIterCursor(j2gort.CheckNotNull(licenseTracker.EmitLicenses()))
			for cur_145.HasNext() {
				license := cur_145.NextStr()
				_ = license
				cb.Append("/*\n").Append(license).Append("*/\n")
			}
		}
		if j2gort.CheckNotNullString(code) == "" {
			return nil
		}
		if this_.options.ShouldGatherSourceMapInfo() {
			this_.sourceMap.SetStartingPosition(cb.GetLineIndex(), cb.GetColumnIndex())
		}
		cb.Append(code)
		var length int32 = j2gort.JStrLen(code)
		_ = length
		var lastChar uint16 = j2gort.JStrCharAt(code, (length - int32(1)))
		_ = lastChar
		var secondLastChar uint16
		_ = secondLastChar
		if length >= int32(2) {
			secondLastChar = j2gort.JStrCharAt(code, (length - int32(2)))
		} else {
			secondLastChar = uint16(0)
		}
		var hasSemiColon bool = ((int32(lastChar) == int32(uint16(59))) || ((int32(lastChar) == int32(uint16(10))) && (int32(secondLastChar) == int32(uint16(59)))))
		_ = hasSemiColon
		if !(hasSemiColon) {
			cb.Append(";")
		}
		if this_.options.ShouldGatherSourceMapInfo() {
			{
				cur_146 := j2gort.NewIterCursor(j2gort.CheckNotNull(sourceAndMappings.mappings))
				for cur_146.HasNext() {
					rawElem_147 := cur_146.Next()
					mapping := j2gort.CastPtr[SourceMap_Mapping](rawElem_147)
					_ = mapping
					this_.sourceMap.AddMapping_1(mapping)
				}
			}
		}
		return nil
	}))
}

// Java method: Compiler.toSource(com.google.javascript.rhino.Node n)
func (this_ *Compiler) ToSource_1_Node(n pkg_rhino.Node) string {
	_ = n
	this_.InitCompilerOptionsIfTesting()
	var sb j2gort.LocalStringBuilder

	var lt *Compiler_ScriptNodeLicensesOnlyTracker = NewCompiler_ScriptNodeLicensesOnlyTracker(this_)
	_ = lt
	var code string = this_.priv_Compiler_toSourceAndMappings_1522(n, false, lt).source
	_ = code
	{
		cur_148 := j2gort.NewIterCursor(j2gort.CheckNotNull(lt.EmitLicenses()))
		for cur_148.HasNext() {
			license := cur_148.NextStr()
			_ = license
			sb.AppendStr("/*\n").AppendStr(license).AppendStr("*/\n")
		}
	}
	sb.AppendStr(code)
	return sb.String()
}

// Java method: Compiler.toSourceAndMappings(com.google.javascript.rhino.Node n, boolean firstOutput, com.google.javascript.jscomp.CodePrinter.LicenseTracker licenseTracker)
func (this_ *Compiler) priv_Compiler_toSourceAndMappings_1522(n pkg_rhino.Node, firstOutput bool, licenseTracker CodePrinter_LicenseTracker) *CodePrinter_SourceAndMappings {
	_ = n
	_ = firstOutput
	_ = licenseTracker
	var val_149 CodePrinter_Builder
	InitCodePrinter_Builder(&val_149, n)
	var builder *CodePrinter_Builder = &val_149
	_ = builder
	builder.SetCompilerOptions(this_.options)
	builder.SetTagAsTypeSummary(this_.options.ShouldGenerateTypedExterns())
	builder.SetTagAsStrict((firstOutput && this_.options.ShouldEmitUseStrict()))
	builder.SetLicenseTracker(licenseTracker)
	return builder.BuildWithSourceMappings()
}

// Java static method: Compiler.getLicenseForFile(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.@org.jspecify.annotations.Nullable String fileName)
func Compiler_GetLicenseForFile(compiler AbstractCompiler, fileName string) string {
	clinit_Compiler()
	_ = compiler
	_ = fileName
	if fileName == j2gort.NullString {
		return j2gort.NullString
	}
	var licenseRoot pkg_rhino.Node = compiler.GetScriptNode(fileName)
	_ = licenseRoot
	if licenseRoot == 0 {
		return j2gort.NullString
	}
	var docInfo *pkg_rhino.JSDocInfo = licenseRoot.GetJSDocInfo()
	_ = docInfo
	if docInfo == nil {
		return j2gort.NullString
	}
	var license string = docInfo.GetLicense()
	_ = license
	return license
}

// Java method: Compiler.toSourceArray(com.google.javascript.jscomp.CodePrinter.LicenseTracker licenseTracker, com.google.javascript.jscomp.JSChunk chunk)
func (this_ *Compiler) ToSourceArray(licenseTracker CodePrinter_LicenseTracker, chunk *JSChunk) []string {
	_ = licenseTracker
	_ = chunk
	return j2gort.Cast[[]string](this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		var inputs []*CompilerInput = j2gort.ImmListToSliceNullable[*CompilerInput](chunk.GetInputs())
		_ = inputs
		var numInputs int32 = int32(len(inputs))
		_ = numInputs
		if numInputs == int32(0) {
			return j2gort.SliceToAny(j2gort.NewStringArray(int32(0)))
		}
		var sources []string = j2gort.NewStringArray(numInputs)
		_ = sources
		var val_150 Compiler_CodeBuilder
		InitCompiler_CodeBuilder(&val_150)
		var cb *Compiler_CodeBuilder = &val_150
		_ = cb
		{
			var i int32 = int32(0)
			_ = i
			for ; i < numInputs; i = i + 1 {
				var scriptNode pkg_rhino.Node = inputs[i].GetAstRoot(this_)
				_ = scriptNode
				if scriptNode == 0 {
					panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Bad chunk input: "+j2gort.StrValString(inputs[i].GetName()))))
				}
				cb.Reset()
				this_.ToSource_4(cb, licenseTracker, i, scriptNode)
				{
					a_151 := sources
					i_152 := i
					a_151[i_152] = cb.String()
				}
			}
		}
		return j2gort.SliceToAny(sources)
	})))
}

// Java method: Compiler.performTranspilationAndOptimizations(com.google.javascript.jscomp.CompilerOptions.SegmentOfCompilationToRun segmentOfCompilationToRun)
func (this_ *Compiler) PerformTranspilationAndOptimizations(segmentOfCompilationToRun CompilerOptions_SegmentOfCompilationToRun) {
	_ = segmentOfCompilationToRun
	defer j2gort.PurgeAllRuntimePools()

	if !(this_.options.ShouldOptimize()) {
		j2gort.CheckState(false)
	}
	var allOptimizationPassesToRun j2gort.Ext_com_google_common_collect_ImmutableList = this_.GetPassConfig().GetOptimizations().Build()
	_ = allOptimizationPassesToRun
	if allOptimizationPassesToRun.IsEmpty() {
		return
	}
	var optimizationPassesToRunInCurrentSegment j2gort.Ext_java_util_List = this_.priv_Compiler_getOptimizationPassesToRunInCurrentSegment_1527(segmentOfCompilationToRun, j2gort.AsIface[j2gort.Ext_java_util_List](allOptimizationPassesToRun))
	_ = optimizationPassesToRunInCurrentSegment
	this_.priv_Compiler_markTranspiledFiles_1479()
	this_.phaseOptimizer = this_.priv_Compiler_createPhaseOptimizer_1481()
	this_.phaseOptimizer.Consume(optimizationPassesToRunInCurrentSegment)
	this_.phaseOptimizer.Process_2(this_.externsRoot, this_.jsRoot)
	this_.phaseOptimizer = nil
}

// Java method: Compiler.performFinalizations()
func (this_ *Compiler) PerformFinalizations() {
	var finalizations j2gort.Ext_com_google_common_collect_ImmutableList = this_.GetPassConfig().GetFinalizations().Build()
	_ = finalizations
	if finalizations.IsEmpty() {
		return
	}
	this_.phaseOptimizer = this_.priv_Compiler_createPhaseOptimizer_1481()
	this_.phaseOptimizer.Consume(j2gort.AsIface[j2gort.Ext_java_util_List](finalizations))
	this_.phaseOptimizer.Process_2(this_.externsRoot, this_.jsRoot)
	this_.phaseOptimizer = nil
}

// Java method: Compiler.computeCFG()
func (this_ *Compiler) ComputeCFG() ControlFlowGraph[pkg_rhino.Node] {
	StaticGet_Compiler_Logger().Fine_1_string("Computing Control Flow Graph")
	var tracer *Tracer = this_.NewTracer("computeCFG")
	_ = tracer
	var cfg ControlFlowGraph[pkg_rhino.Node] = ControlFlowAnalysis_S_Builder().SetCompiler(this_).SetCfgRoot(this_.jsRoot).SetTraverseFunctions(true).ComputeCfg()
	_ = cfg
	this_.StopTracer(tracer, "computeCFG")
	return cfg
}

// Java method: Compiler.getOptimizationPassesToRunInCurrentSegment(com.google.javascript.jscomp.CompilerOptions.SegmentOfCompilationToRun segmentOfCompilationToRun, java.util.List<com.google.javascript.jscomp.PassFactory> allOptimizationPassesToRun)
func (this_ *Compiler) priv_Compiler_getOptimizationPassesToRunInCurrentSegment_1527(segmentOfCompilationToRun CompilerOptions_SegmentOfCompilationToRun, allOptimizationPassesToRun j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	_ = segmentOfCompilationToRun
	_ = allOptimizationPassesToRun
	switch (segmentOfCompilationToRun).Ordinal() {
	case int32(4) /* SegmentOfCompilationToRun.OPTIMIZATIONS */ :
		return allOptimizationPassesToRun
	case int32(2) /* SegmentOfCompilationToRun.OPTIMIZATIONS_FIRST_HALF */ :
		return this_.priv_Compiler_passListUntil_1528(allOptimizationPassesToRun, "optimizationsHalfwayPoint")
	case int32(3) /* SegmentOfCompilationToRun.OPTIMIZATIONS_SECOND_HALF */ :
		return this_.priv_Compiler_passListAfter_1529(allOptimizationPassesToRun, "optimizationsHalfwayPoint")
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unsupported segment of optimizations to run: "+j2gort.StrVal(segmentOfCompilationToRun))))
	}
	panic("unreachable switch expression")
}

// Java method: Compiler.passListUntil(java.util.List<com.google.javascript.jscomp.PassFactory> passes, java.lang.String passName)
func (this_ *Compiler) priv_Compiler_passListUntil_1528(passes j2gort.Ext_java_util_List, passName string) j2gort.Ext_java_util_List {
	_ = passes
	_ = passName
	var passList []PassFactory = []PassFactory{}
	_ = passList
	{
		cur_153 := j2gort.NewIterCursor(j2gort.CheckNotNull(passes))
		for cur_153.HasNext() {
			rawElem_154 := cur_153.Next()
			pass := j2gort.AsIface[PassFactory](rawElem_154)
			_ = pass
			if j2gort.CheckNotNullString(pass.GetName()) == passName {
				return j2gort.WrapTypedArrayList(passList)
			}
			passList = append(passList, pass)
		}
	}
	return j2gort.WrapTypedArrayList(passList)
}

// Java method: Compiler.passListAfter(java.util.List<com.google.javascript.jscomp.PassFactory> passes, java.lang.String passName)
func (this_ *Compiler) priv_Compiler_passListAfter_1529(passes j2gort.Ext_java_util_List, passName string) j2gort.Ext_java_util_List {
	_ = passes
	_ = passName
	var passList []PassFactory = []PassFactory{}
	_ = passList
	{
		cur_155 := j2gort.NewIterCursor(j2gort.CheckNotNull(passes))
		for cur_155.HasNext() {
			rawElem_156 := cur_155.Next()
			pass := j2gort.AsIface[PassFactory](rawElem_156)
			_ = pass
			if j2gort.CheckNotNullString(pass.GetName()) == passName {
				clear(passList)
				passList = passList[:0]
				continue
			}
			passList = append(passList, pass)
		}
	}
	return j2gort.WrapTypedArrayList(passList)
}

// Java method: Compiler.getChangeTracker()
func (this_ *Compiler) GetChangeTracker() *ChangeTracker {
	return this_.changeTracker
}

// Java method: Compiler.addIndexProvider(com.google.javascript.jscomp.IndexProvider<?> indexProvider)
func (this_ *Compiler) AddIndexProvider(indexProvider IndexProvider) {
	_ = indexProvider
	var type_ *javalang.Class = indexProvider.GetType()
	_ = type_
	if !(j2gort.IfaceEquals[IndexProvider](j2gort.AsIface[IndexProvider](this_.indexProvidersByType.Put_2(j2gort.PtrToAny(type_), indexProvider)), nil)) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "A provider is already registered for index of type "+j2gort.StrValString(j2gort.CheckNotNullPtr(type_).GetSimpleName()))))
	}
}

// Java method: Compiler.getIndex(java.lang.Class<T> key)
func (this_ *Compiler) GetIndex(key *javalang.Class) any {
	_ = key
	var indexProvider IndexProvider = j2gort.AsIface[IndexProvider](this_.indexProvidersByType.Get(j2gort.PtrToAny(key)))
	_ = indexProvider
	if indexProvider == nil || indexProvider.IsJavaNil_() {
		return nil
	}
	return indexProvider.Get()
}

// Java method: Compiler.getExternsRoot()
func (this_ *Compiler) GetExternsRoot() pkg_rhino.Node {
	return this_.externsRoot
}

// Java method: Compiler.getJsRoot()
func (this_ *Compiler) GetJsRoot() pkg_rhino.Node {
	return this_.jsRoot
}

// Java method: Compiler.hasScopeChanged(com.google.javascript.rhino.Node n)
func (this_ *Compiler) HasScopeChanged(n pkg_rhino.Node) bool {
	_ = n
	if this_.phaseOptimizer == nil {
		return true
	}
	return this_.phaseOptimizer.HasScopeChanged(n)
}

// Java method: Compiler.reportChangeToChangeScope(com.google.javascript.rhino.Node changeScopeRoot)
func (this_ *Compiler) ReportChangeToChangeScope(changeScopeRoot pkg_rhino.Node) {
	_ = changeScopeRoot
	this_.changeTracker.ReportChangeToChangeScope(changeScopeRoot)
}

// Java method: Compiler.reportFunctionDeleted(com.google.javascript.rhino.Node n)
func (this_ *Compiler) ReportFunctionDeleted(n pkg_rhino.Node) {
	_ = n
	this_.changeTracker.ReportFunctionDeleted(n)
}

// Java method: Compiler.reportDisambiguatePropertiesSummary(java.util.function.Supplier<java.lang.String> summarySupplier)
func (this_ *Compiler) ReportDisambiguatePropertiesSummary(summarySupplier j2gort.Ext_java_util_function_Supplier) {
	_ = summarySupplier
	if this_.tracker != nil {
		this_.tracker.SetDisambiguatePropertiesSummary(j2gort.CastString(summarySupplier.Get()))
	}
}

// Java method: Compiler.reportAmbiguatePropertiesSummary(java.util.function.Supplier<java.lang.String> summarySupplier)
func (this_ *Compiler) ReportAmbiguatePropertiesSummary(summarySupplier j2gort.Ext_java_util_function_Supplier) {
	_ = summarySupplier
	if this_.tracker != nil {
		this_.tracker.SetAmbiguatePropertiesSummary(j2gort.CastString(summarySupplier.Get()))
	}
}

// Java method: Compiler.reportChangeToEnclosingScope(com.google.javascript.rhino.Node n)
func (this_ *Compiler) ReportChangeToEnclosingScope(n pkg_rhino.Node) {
	_ = n
	this_.changeTracker.ReportChangeToEnclosingScope(n)
}

// Java method: Compiler.getCodingConvention()
func (this_ *Compiler) GetCodingConvention() CodingConvention {
	var convention CodingConvention = this_.options.GetCodingConvention()
	_ = convention
	if convention != nil && !convention.IsJavaNil_() {
		convention = convention
	} else {
		convention = this_.defaultCodingConvention
	}
	return convention
}

// Java method: Compiler.getParserConfigLanguageMode(com.google.javascript.jscomp.CompilerOptions.LanguageMode languageMode)
func (this_ *Compiler) priv_Compiler_getParserConfigLanguageMode_1531(languageMode CompilerOptions_LanguageMode) *pkg_parsing.Config_LanguageMode {
	_ = languageMode
	switch (languageMode).Ordinal() {
	case int32(0) /* LanguageMode.ECMASCRIPT3 */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT3()
	case int32(1) /* LanguageMode.ECMASCRIPT5 */, int32(2) /* LanguageMode.ECMASCRIPT5_STRICT */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT5()
	case int32(3) /* LanguageMode.ECMASCRIPT_2015 */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT_2015()
	case int32(4) /* LanguageMode.ECMASCRIPT_2016 */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT_2016()
	case int32(5) /* LanguageMode.ECMASCRIPT_2017 */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT_2017()
	case int32(6) /* LanguageMode.ECMASCRIPT_2018 */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT_2018()
	case int32(7) /* LanguageMode.ECMASCRIPT_2019 */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT_2019()
	case int32(8) /* LanguageMode.ECMASCRIPT_2020 */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT_2020()
	case int32(9) /* LanguageMode.ECMASCRIPT_2021 */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT_2021()
	case int32(11) /* LanguageMode.ECMASCRIPT_NEXT */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_ES_NEXT()
	case int32(14) /* LanguageMode.UNSTABLE */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_UNSTABLE()
	case int32(15) /* LanguageMode.UNSUPPORTED */ :
		return pkg_parsing.StaticGet_Config_LanguageMode_UNSUPPORTED()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected language mode: "+j2gort.StrVal(this_.options.GetLanguageIn()))))
	}
	panic("unreachable switch expression")
}

// Java method: Compiler.getParserConfig(com.google.javascript.jscomp.AbstractCompiler.ConfigContext context)
func (this_ *Compiler) GetParserConfig(context AbstractCompiler_ConfigContext) pkg_parsing.Config {
	_ = context
	if (this_.parserConfig == nil || this_.parserConfig.IsJavaNil_()) || (this_.externsParserConfig == nil || this_.externsParserConfig.IsJavaNil_()) {
		func() {
			if this_ == nil {
				j2gort.ThrowNPE()
			}
			this_._syncMu.Lock()
			defer this_._syncMu.Unlock()
			if this_.parserConfig == nil || this_.parserConfig.IsJavaNil_() {
				var configLanguageMode *pkg_parsing.Config_LanguageMode = this_.priv_Compiler_getParserConfigLanguageMode_1531(this_.options.GetLanguageIn())
				_ = configLanguageMode
				var strictMode pkg_parsing.Config_StrictMode
				_ = strictMode
				if this_.options.ExpectStrictModeInput() {
					strictMode = pkg_parsing.StaticGet_Config_StrictMode_STRICT()
				} else {
					strictMode = pkg_parsing.StaticGet_Config_StrictMode_SLOPPY()
				}
				this_.parserConfig = this_.CreateConfig(configLanguageMode, strictMode)
				if j2gort.RefEquals(configLanguageMode, pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT3()) {
					this_.externsParserConfig = this_.CreateConfig(pkg_parsing.StaticGet_Config_LanguageMode_ECMASCRIPT5(), strictMode)
				} else {
					this_.externsParserConfig = this_.parserConfig
				}
			}
		}()
	}
	if context == StaticGet_AbstractCompiler_ConfigContext_EXTERNS() {
		return this_.externsParserConfig
	}
	return this_.parserConfig
}

// Java method: Compiler.createConfig(com.google.javascript.jscomp.parsing.Config.LanguageMode mode, com.google.javascript.jscomp.parsing.Config.StrictMode strictMode)
func (this_ *Compiler) CreateConfig(mode *pkg_parsing.Config_LanguageMode, strictMode pkg_parsing.Config_StrictMode) pkg_parsing.Config {
	_ = mode
	_ = strictMode
	return pkg_parsing.ParserRunner_CreateConfig_6(mode, this_.options.IsParseJsDocDocumentation(), func() pkg_parsing.Config_RunMode {
		if this_.options.CanContinueAfterErrors() {
			return pkg_parsing.StaticGet_Config_RunMode_KEEP_GOING()
		}
		return pkg_parsing.StaticGet_Config_RunMode_STOP_AFTER_ERROR()
	}(), j2gort.AsIface[j2gort.Ext_java_util_Set](this_.options.GetExtraAnnotationNames()), this_.options.GetParseInlineSourceMaps(), strictMode)
}

// Java method: Compiler.getDiagnosticGroups()
func (this_ *Compiler) GetDiagnosticGroups() *DiagnosticGroups {
	return NewDiagnosticGroups()
}

// Java method: Compiler.report(com.google.javascript.jscomp.JSError error)
func (this_ *Compiler) Report_1(error_ *JSError) {
	_ = error_
	var level *CheckLevel = error_.DefaultLevel()
	_ = level
	if this_.warningsGuard != nil && !this_.warningsGuard.IsJavaNil_() {
		var newLevel *CheckLevel = this_.warningsGuard.Level(error_)
		_ = newLevel
		if newLevel != nil {
			level = newLevel
		}
	}
	if level.IsOn() {
		this_.InitCompilerOptionsIfTesting()
		if !(j2gort.IfaceEquals[ErrorHandler](this_.GetOptions().GetErrorHandler(), nil)) {
			this_.GetOptions().GetErrorHandler().Report_2(level, error_)
		}
		this_.errorManager.Report_2(level, error_)
	}
}

// Java method: Compiler.report(com.google.javascript.jscomp.CheckLevel ignoredLevel, com.google.javascript.jscomp.JSError error)
func (this_ *Compiler) Report_2(ignoredLevel *CheckLevel, error_ *JSError) {
	_ = ignoredLevel
	_ = error_
	this_.Report_1(error_)
}

// Java method: Compiler.getErrorLevel(com.google.javascript.jscomp.JSError error)
func (this_ *Compiler) GetErrorLevel(error_ *JSError) *CheckLevel {
	_ = error_
	j2gort.CheckNotNullPtr(this_.options)
	return this_.warningsGuard.Level(error_)
}

// Java method: Compiler.throwInternalError(java.lang.String message, java.lang.Throwable cause)
func (this_ *Compiler) ThrowInternalError(message string, cause javalang.Throwable) {
	_ = message
	_ = cause
	panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "INTERNAL COMPILER ERROR.\nPlease report this problem.\n\n"+j2gort.StrValString(message), cause)))
}

// Java method: Compiler.getErrorCount()
func (this_ *Compiler) GetErrorCount() int32 {
	return this_.errorManager.GetErrorCount()
}

// Java method: Compiler.getWarningCount()
func (this_ *Compiler) GetWarningCount() int32 {
	return this_.errorManager.GetWarningCount()
}

// Java method: Compiler.hasHaltingErrors()
func (this_ *Compiler) HasHaltingErrors() bool {
	return (!(this_.GetOptions().CanContinueAfterErrors()) && this_.errorManager.HasHaltingErrors())
}

// Java method: Compiler.hasErrors()
func (this_ *Compiler) HasErrors() bool {
	return this_.HasHaltingErrors()
}

// Java method: Compiler.getSourceFileByName(java.lang.String sourceName)
func (this_ *Compiler) GetSourceFileByName(sourceName string) *SourceFile {
	_ = sourceName
	if sourceName != j2gort.NullString {
		var input *CompilerInput = j2gort.CastPtr[CompilerInput](this_.inputsById.Get(pkg_rhino.NewInputId(sourceName)))
		_ = input
		if input != nil {
			return input.GetSourceFile()
		}
		return j2gort.CastPtr[SourceFile](j2gort.MapGetStr(this_.sourceMapOriginalSources, sourceName))
	}
	return nil
}

// Java method: Compiler.getSourceFileContentByName(java.lang.String sourceName)
func (this_ *Compiler) GetSourceFileContentByName(sourceName string) any {
	_ = sourceName
	var file *SourceFile = this_.GetSourceFileByName(sourceName)
	_ = file
	if file.IsStubSourceFileForAlreadyProvidedInput() {
		return nil
	}
	j2gort.CheckNotNullPtr(file)
	var didRet_157 bool
	var retVal_158 any
	func() {
		defer func() {
			if rec_159 := recover(); rec_159 != nil {
				jex_160 := javalang.RecoverException(rec_159)
				_ = jex_160
				if j2gort.InstanceOf[javaio.IIOException](jex_160) {
					e := j2gort.AsIface[javaio.IIOException](jex_160)
					_ = e
					retVal_158 = nil
					didRet_157 = true
					return
				} else {
					panic(rec_159)
				}
			}
		}()
		retVal_158 = j2gort.BoxString(file.GetCode())
		didRet_157 = true
		return
	}()
	if didRet_157 {
		return retVal_158
	}
	panic("unreachable")
}

// Java method: Compiler.addInputSourceMap(java.lang.String sourceFileName, com.google.javascript.jscomp.SourceMapInput inputSourceMap)
func (this_ *Compiler) AddInputSourceMap(sourceFileName string, inputSourceMap *SourceMapInput) {
	_ = sourceFileName
	_ = inputSourceMap
	j2gort.MapPutStr(this_.inputSourceMaps, sourceFileName, j2gort.PtrToAny(inputSourceMap))
	if this_.options.GetSourceMapIncludeSourcesContent() && (this_.sourceMap != nil) {
		this_.priv_Compiler_addSourceMapSourceFiles_1538(inputSourceMap)
	}
}

// Java method: Compiler.getBase64SourceMapContents(java.lang.String sourceFileName)
func (this_ *Compiler) GetBase64SourceMapContents(sourceFileName string) string {
	_ = sourceFileName
	var sourceMapInput *SourceMapInput = j2gort.CastPtr[SourceMapInput](j2gort.MapGetOrDefaultStr(this_.inputSourceMaps, sourceFileName, nil))
	_ = sourceMapInput
	var sourceMappingURL string
	_ = sourceMappingURL
	if sourceMapInput != nil {
		sourceMappingURL = sourceMapInput.GetOriginalPath()
	} else {
		sourceMappingURL = j2gort.NullString
	}
	if (sourceMappingURL != j2gort.NullString) && javalang.StringEndsWith(j2gort.CheckNotNullString(sourceMappingURL), ".inline.map") {
		var unencoded string = sourceMapInput.GetRawSourceMapContents()
		_ = unencoded
		return j2gort.Ext_com_google_common_io_BaseEncoding_Base64().Encode_1(func() []int8 {
			_ = j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8
			return j2gort.JStrGetBytes(unencoded)
		}())
	}
	return j2gort.NullString
}

// Java method: Compiler.addSourceMapSourceFiles(com.google.javascript.jscomp.SourceMapInput inputSourceMap)
func (this_ *Compiler) priv_Compiler_addSourceMapSourceFiles_1538(inputSourceMap *SourceMapInput) {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	this_._syncMu.Lock()
	defer this_._syncMu.Unlock()
	_ = inputSourceMap
	var consumer *pkg_sourcemap.SourceMapConsumerV3 = inputSourceMap.GetSourceMap(this_.errorManager)
	_ = consumer
	if consumer == nil {
		return
	}
	var sourcesContent j2gort.Ext_java_util_Collection = consumer.GetOriginalSourcesContent()
	_ = sourcesContent
	if sourcesContent == nil {
		return
	}
	var content j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](sourcesContent.Iterator())
	_ = content
	var sources j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](consumer.GetOriginalSources().Iterator())
	_ = sources
	for sources.HasNext() && content.HasNext() {
		var code string = j2gort.CastString(content.Next())
		_ = code
		var source *SourceFile = SourceMapResolver_GetRelativePath(inputSourceMap.GetOriginalPath(), j2gort.CastString(sources.Next()))
		_ = source
		if source != nil {
			this_.sourceMap.AddSourceFile(source.GetName(), code)
		}
	}
	if sources.HasNext() || content.HasNext() {
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Source map's \"sources\" and \"sourcesContent\" lengths do not match.")))
	}
}

// Java method: Compiler.getSourceMapping(java.lang.String sourceName, int lineNumber, int columnNumber)
func (this_ *Compiler) GetSourceMapping(sourceName string, lineNumber int32, columnNumber int32) *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping {
	_ = sourceName
	_ = lineNumber
	_ = columnNumber
	if sourceName == j2gort.NullString {
		return nil
	}
	var sourceMap *SourceMapInput = j2gort.CastPtr[SourceMapInput](j2gort.MapGetStr(this_.inputSourceMaps, sourceName))
	_ = sourceMap
	if sourceMap == nil {
		return nil
	}
	var consumer *pkg_sourcemap.SourceMapConsumerV3 = sourceMap.GetSourceMap(this_.errorManager)
	_ = consumer
	if consumer == nil {
		return nil
	}
	var result *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping = consumer.GetMappingForLine(lineNumber, (columnNumber + int32(1)))
	_ = result
	if result == nil {
		return j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_GetDefaultInstance()
	}
	var sourceMapOriginalPath string = sourceMap.GetOriginalPath()
	_ = sourceMapOriginalPath
	var resultOriginalPath string = result.GetOriginalFile()
	_ = resultOriginalPath
	var relativePath string = j2gort.NullString
	_ = relativePath
	{
		var iifeCond_1 bool
		{

			sr_161 := sourceMapOriginalPath
			se_162 := this_.resolvedSourceMap.originalPath
			iifeCond_1 = j2gort.CheckNotNullString(sr_161) == se_162
		}
		if iifeCond_1 && func() bool {
			sr_163 := resultOriginalPath
			se_164 := this_.resolvedSourceMap.sourceMapPath
			return j2gort.CheckNotNullString(sr_163) == se_164
		}() {
			relativePath = this_.resolvedSourceMap.relativePath
		} else {
			relativePath = Compiler_ResolveSibling(sourceMapOriginalPath, resultOriginalPath)
			var source *SourceFile = this_.GetSourceFileByName(relativePath)
			_ = source
			if (source == nil) && !(j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(resultOriginalPath)) {
				source = SourceMapResolver_GetRelativePath(sourceMap.GetOriginalPath(), result.GetOriginalFile())
				if source != nil {
					this_.sourceMapOriginalSources.PutIfAbsent(j2gort.BoxString(relativePath), j2gort.PtrToAny(source))
				}
			}
			{
				r_165 := this_.resolvedSourceMap
				r_165.originalPath = sourceMapOriginalPath
			}
			{
				r_166 := this_.resolvedSourceMap
				r_166.sourceMapPath = resultOriginalPath
			}
			{
				r_167 := this_.resolvedSourceMap
				r_167.relativePath = relativePath
			}
		}
	}
	return j2gort.CastPtr[j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping](j2gort.CastPtr[j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Builder](result.ToBuilder()).SetOriginalFile(relativePath).SetColumnPosition((result.GetColumnPosition() - int32(1))).Build())
}

// Java method: Compiler.getSourceLine(java.lang.String sourceName, int lineNumber)
func (this_ *Compiler) GetSourceLine(sourceName string, lineNumber int32) string {
	_ = sourceName
	_ = lineNumber
	if lineNumber < int32(1) {
		return j2gort.NullString
	}
	var input *SourceFile = this_.GetSourceFileByName(sourceName)
	_ = input
	if (input != nil) && !(input.IsStubSourceFileForAlreadyProvidedInput()) {
		return input.GetLine(lineNumber)
	}
	return j2gort.NullString
}

// Java method: Compiler.getSourceLines(java.lang.String sourceName, int lineNumber, int length)
func (this_ *Compiler) GetSourceLines(sourceName string, lineNumber int32, length int32) Region {
	_ = sourceName
	_ = lineNumber
	_ = length
	if lineNumber < int32(1) {
		return nil
	}
	var input *SourceFile = this_.GetSourceFileByName(sourceName)
	_ = input
	if (input != nil) && !(input.IsStubSourceFileForAlreadyProvidedInput()) {
		return input.GetLines(lineNumber, length)
	}
	return nil
}

// Java method: Compiler.getSourceRegion(java.lang.String sourceName, int lineNumber)
func (this_ *Compiler) GetSourceRegion(sourceName string, lineNumber int32) Region {
	_ = sourceName
	_ = lineNumber
	if lineNumber < int32(1) {
		return nil
	}
	var input *SourceFile = this_.GetSourceFileByName(sourceName)
	_ = input
	if (input != nil) && !(input.IsStubSourceFileForAlreadyProvidedInput()) {
		return input.GetRegion(lineNumber)
	}
	return nil
}

// Java method: Compiler.getNodeForCodeInsertion(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable JSChunk chunk)
func (this_ *Compiler) GetNodeForCodeInsertion(chunk *JSChunk) pkg_rhino.Node {
	_ = chunk
	if this_.inputsById.ContainsKey(j2gort.PtrToAny(StaticGet_Compiler_SYNTHETIC_CODE_INPUT_ID())) {
		return j2gort.CastPtr[CompilerInput](this_.inputsById.Get(j2gort.PtrToAny(StaticGet_Compiler_SYNTHETIC_CODE_INPUT_ID()))).GetAstRoot(this_)
	}
	if chunk == nil {
		var allInputs j2gort.Ext_java_lang_Iterable = this_.chunkGraph.GetAllInputs()
		_ = allInputs
		if (this_.chunkGraph == nil) || j2gort.Ext_com_google_common_collect_Iterables_IsEmpty(allInputs) {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "No inputs")))
		}
		var firstInput *CompilerInput = j2gort.CastPtr[CompilerInput](j2gort.Ext_com_google_common_collect_Iterables_GetFirst(allInputs, nil))
		_ = firstInput
		return Compiler_CheckNotModule(firstInput.GetAstRoot(this_), "Cannot insert code into a module", []any{})
	}
	var chunkInputs []*CompilerInput = j2gort.ImmListToSliceNullable[*CompilerInput](chunk.GetInputs())
	_ = chunkInputs
	if len(chunkInputs) != 0 {
		return Compiler_CheckNotModule(chunkInputs[int32(0)].GetAstRoot(this_), "Cannot insert code into a module", []any{})
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Root chunk has no inputs")))
}

// Java method: Compiler.getSourceMap()
func (this_ *Compiler) GetSourceMap() *SourceMap {
	return this_.sourceMap
}

// Java method: Compiler.setVariableMap(com.google.javascript.jscomp.VariableMap variableMap)
func (this_ *Compiler) SetVariableMap(variableMap *VariableMap) {
	_ = variableMap
	this_.variableMap = variableMap
}

// Java method: Compiler.getVariableMap()
func (this_ *Compiler) GetVariableMap() *VariableMap {
	return this_.variableMap
}

// Java method: Compiler.setPropertyMap(com.google.javascript.jscomp.VariableMap propertyMap)
func (this_ *Compiler) SetPropertyMap(propertyMap *VariableMap) {
	_ = propertyMap
	this_.propertyMap = propertyMap
}

// Java method: Compiler.getPropertyMap()
func (this_ *Compiler) GetPropertyMap() *VariableMap {
	return this_.propertyMap
}

// Java method: Compiler.setStringMap(com.google.javascript.jscomp.VariableMap stringMap)
func (this_ *Compiler) SetStringMap(stringMap *VariableMap) {
	_ = stringMap
	this_.stringMap = stringMap
}

// Java method: Compiler.setCssNames(java.util.Set<java.lang.String> cssNames)
func (this_ *Compiler) SetCssNames(cssNames j2gort.Ext_java_util_Set) {
	_ = cssNames
	this_.cssNames = cssNames
}

// Java method: Compiler.setIdGeneratorMap(java.lang.String serializedIdMappings)
func (this_ *Compiler) SetIdGeneratorMap(serializedIdMappings string) {
	_ = serializedIdMappings
	this_.idGeneratorMap = serializedIdMappings
}

// Java method: Compiler.getTranspiledFiles()
func (this_ *Compiler) GetTranspiledFiles() bool {
	return this_.transpiledFiles
}

// Java method: Compiler.setTranspiledFiles(boolean transpiledFiles)
func (this_ *Compiler) SetTranspiledFiles(transpiledFiles bool) {
	_ = transpiledFiles
	this_.transpiledFiles = transpiledFiles
}

// Java method: Compiler.getCrossChunkIdGenerator()
func (this_ *Compiler) GetCrossChunkIdGenerator() *IdGenerator {
	return this_.crossChunkIdGenerator
}

// Java method: Compiler.setAnonymousFunctionNameMap(com.google.javascript.jscomp.VariableMap functionMap)
func (this_ *Compiler) SetAnonymousFunctionNameMap(functionMap *VariableMap) {
	_ = functionMap
}

// Java method: Compiler.getStringMap()
func (this_ *Compiler) GetStringMap() *VariableMap {
	return this_.stringMap
}

// Java method: Compiler.setInstrumentationMapping(com.google.javascript.jscomp.VariableMap instrumentationMapping)
func (this_ *Compiler) SetInstrumentationMapping(instrumentationMapping *VariableMap) {
	_ = instrumentationMapping
	this_.instrumentationMapping = instrumentationMapping
}

// Java method: Compiler.getInstrumentationMapping()
func (this_ *Compiler) GetInstrumentationMapping() *VariableMap {
	return this_.instrumentationMapping
}

// Java method: Compiler.addExportedNames(java.util.Set<java.lang.String> exportedNames)
func (this_ *Compiler) AddExportedNames(exportedNames j2gort.Ext_java_util_Set) {
	_ = exportedNames
	this_.exportedNames.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](exportedNames))
}

// Java method: Compiler.getExportedNames()
func (this_ *Compiler) GetExportedNames() j2gort.Ext_java_util_Set {
	return j2gort.AsIface[j2gort.Ext_java_util_Set](this_.exportedNames)
}

// Java method: Compiler.setDefineNames(java.util.Collection<java.lang.String> defineNames)
func (this_ *Compiler) SetDefineNames(defineNames j2gort.Ext_java_util_Collection) {
	_ = defineNames
	this_.defineNames = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(defineNames)
}

// Java method: Compiler.getDefineNames()
func (this_ *Compiler) GetDefineNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.defineNames
}

// Java method: Compiler.getOptions()
func (this_ *Compiler) GetOptions() *CompilerOptions {
	return this_.options
}

// Java static method: Compiler.setLoggingLevel(java.util.logging.Level level)
func Compiler_SetLoggingLevel(level *j2gort.Ext_java_util_logging_Level) {
	clinit_Compiler()
	_ = level
	StaticGet_Compiler_Logger().SetLevel(level)
}

// Java method: Compiler.getAstDotGraph()
func (this_ *Compiler) GetAstDotGraph() string {
	if this_.jsRoot != 0 {
		var cfg ControlFlowGraph[pkg_rhino.Node] = ControlFlowAnalysis_S_Builder().SetCompiler(this_).SetCfgRoot(this_.jsRoot).SetTraverseFunctions(true).ComputeCfg()
		_ = cfg
		if cfg != nil && !cfg.IsJavaNil_() {
			if r_, ok_ := any(cfg).(interface{ ReleaseSlabs_() }); ok_ {
				defer r_.ReleaseSlabs_()
			}
		}

		return DotFormatter_ToDot_2(this_.jsRoot, cfg)
	} else {
		return ""
	}
}

// Java method: Compiler.getErrorManager()
func (this_ *Compiler) GetErrorManager() ErrorManager {
	if this_.options == nil {
		this_.InitOptions(NewCompilerOptions())
	}
	return this_.errorManager
}

// Java method: Compiler.getInputsInOrder()
func (this_ *Compiler) GetInputsInOrder() j2gort.Ext_java_lang_Iterable {
	return this_.chunkGraph.GetAllInputs()
}

// Java method: Compiler.getNumberOfInputs()
func (this_ *Compiler) GetNumberOfInputs() int32 {
	if this_.chunkGraph != nil {
		return this_.chunkGraph.GetInputCount()
	} else {
		return int32(1)
	}
}

// Java method: Compiler.getInputsById()
func (this_ *Compiler) GetInputsById() j2gort.Ext_java_util_Map {
	return j2gort.Ext_java_util_Collections_UnmodifiableMap(this_.inputsById)
}

// Java method: Compiler.getExternsInOrder()
func (this_ *Compiler) GetExternsInOrder() j2gort.Ext_java_util_List {
	return j2gort.Ext_java_util_Collections_UnmodifiableList(j2gort.AsIface[j2gort.Ext_java_util_List](this_.externs))
}

// Java method: Compiler.getInputsForTesting()
func (this_ *Compiler) GetInputsForTesting() j2gort.Ext_java_util_List {
	if this_.chunkGraph != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(this_.chunkGraph.GetAllInputs()))
	} else {
		return nil
	}
}

// Java method: Compiler.getExternsForTesting()
func (this_ *Compiler) GetExternsForTesting() j2gort.Ext_java_util_List {
	return j2gort.AsIface[j2gort.Ext_java_util_List](this_.externs)
}

// Java method: Compiler.hasRegExpGlobalReferences()
func (this_ *Compiler) HasRegExpGlobalReferences() bool {
	return this_.hasRegExpGlobalReferences
}

// Java method: Compiler.setHasRegExpGlobalReferences(boolean references)
func (this_ *Compiler) SetHasRegExpGlobalReferences(references bool) {
	_ = references
	this_.hasRegExpGlobalReferences = references
}

// Java method: Compiler.setRunJ2clPasses(boolean runJ2clPasses)
func (this_ *Compiler) SetRunJ2clPasses(runJ2clPasses bool) {
	_ = runJ2clPasses
	this_.runJ2clPasses = runJ2clPasses
}

// Java method: Compiler.runJ2clPasses()
func (this_ *Compiler) RunJ2clPasses() bool {
	return this_.runJ2clPasses
}

// Java method: Compiler.getSynthesizedExternsInput()
func (this_ *Compiler) GetSynthesizedExternsInput() *CompilerInput {
	if this_.syntheticExternsInput != nil {
		return this_.syntheticExternsInput
	}
	var val_168 CompilerInput
	InitCompilerInput_2_SourceFile_bool(&val_168, this_.u_SYNTHETIC_EXTERNS_FILE, true)
	var input *CompilerInput = &val_168
	_ = input
	var root pkg_rhino.Node = j2gort.CheckNotNullNode(input.GetAstRoot(this_))
	_ = root
	this_.priv_Compiler_putCompilerInput_1494(input)
	this_.syntheticExternsInput = input
	this_.externsRoot.AddChildToFront(root)
	this_.externs.Add_2(int32(0), input)
	j2gort.MapPutStr(this_.scriptNodeByFilename, input.GetSourceFile().GetName(), j2gort.BoxNode(root))
	return input
}

// Java method: Compiler.getSynthesizedExternsRoot()
func (this_ *Compiler) priv_Compiler_getSynthesizedExternsRoot_1554() pkg_rhino.Node {
	var externsInput *CompilerInput = this_.GetSynthesizedExternsInput()
	_ = externsInput
	return externsInput.GetAstRoot(this_)
}

// Java method: Compiler.getSynthesizedTypeSummaryInput()
func (this_ *Compiler) GetSynthesizedTypeSummaryInput() *CompilerInput {
	if this_.syntheticTypeSummaryInput != nil {
		return this_.syntheticTypeSummaryInput
	}
	var val_169 CompilerInput
	InitCompilerInput_2_SourceFile_bool(&val_169, this_.syntheticTypeSummaryFile, true)
	var input *CompilerInput = &val_169
	_ = input
	var root pkg_rhino.Node = j2gort.CheckNotNullNode(input.GetAstRoot(this_))
	_ = root
	this_.priv_Compiler_putCompilerInput_1494(input)
	this_.syntheticTypeSummaryInput = input
	this_.externsRoot.AddChildToFront(root)
	this_.externs.Add_2(int32(0), input)
	j2gort.MapPutStr(this_.scriptNodeByFilename, input.GetSourceFile().GetName(), j2gort.BoxNode(root))
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
	_ = builder
	builder.RecordTypeSummary()
	root.SetJSDocInfo(builder.Build_0())
	return input
}

// Java method: Compiler.getSyntheticCodeInputId()
func (this_ *Compiler) GetSyntheticCodeInputId() *pkg_rhino.InputId {
	return StaticGet_Compiler_SYNTHETIC_CODE_INPUT_ID()
}

// Java method: Compiler.initializeSyntheticCodeInput()
func (this_ *Compiler) InitializeSyntheticCodeInput() {
	if !(!(this_.inputsById.ContainsKey(j2gort.PtrToAny(StaticGet_Compiler_SYNTHETIC_CODE_INPUT_ID())))) {
		j2gort.CheckStateStr(false, "Already initialized synthetic input")
	}
	var val_170 CompilerInput
	InitCompilerInput_1(&val_170, SourceFile_FromCode_2(StaticGet_Compiler_SYNTHETIC_CODE_INPUT_ID().GetIdName(), ""))
	var ast *CompilerInput = &val_170
	_ = ast
	if this_.inputsById.ContainsKey(j2gort.PtrToAny(ast.GetInputId())) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Conflicting synthetic id name")))
	}
	this_.jsRoot.AddChildToFront(j2gort.CheckNotNullNode(ast.GetAstRoot(this_)))
	var firstChunk *JSChunk = j2gort.CastPtr[JSChunk](j2gort.Ext_com_google_common_collect_Iterables_GetFirst(this_.GetChunks(), nil))
	_ = firstChunk
	{
		var iifeCond_1 bool
		{

			sr_171 := firstChunk.GetName()
			se_172 := "$strong$"
			iifeCond_1 = j2gort.CheckNotNullString(sr_171) == se_172
		}
		if iifeCond_1 {
			firstChunk.Add_1_CompilerInput(ast)
		}
	}
	ast.SetChunk(firstChunk)
	this_.priv_Compiler_putCompilerInput_1494(ast)
	j2gort.MapPutStr(this_.commentsPerFile, StaticGet_Compiler_SYNTHETIC_CODE_INPUT_ID().GetIdName(), j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	this_.ReportChangeToChangeScope(ast.GetAstRoot(this_))
}

// Java method: Compiler.removeSyntheticCodeInput()
func (this_ *Compiler) RemoveSyntheticCodeInput_0() {
	this_.priv_Compiler_removeSyntheticCodeInput_1555(false)
}

// Java method: Compiler.mergeSyntheticCodeInput()
func (this_ *Compiler) MergeSyntheticCodeInput() {
	this_.priv_Compiler_removeSyntheticCodeInput_1555(true)
}

// Java method: Compiler.removeSyntheticCodeInput(boolean mergeContentIntoFirstInput)
func (this_ *Compiler) priv_Compiler_removeSyntheticCodeInput_1555(mergeContentIntoFirstInput bool) {
	_ = mergeContentIntoFirstInput
	if !(this_.inputsById.ContainsKey(j2gort.PtrToAny(StaticGet_Compiler_SYNTHETIC_CODE_INPUT_ID()))) {
		j2gort.CheckStateStr(false, "Never initialized the synthetic input")
	}
	var input *CompilerInput = j2gort.CastPtr[CompilerInput](this_.inputsById.Get(j2gort.PtrToAny(StaticGet_Compiler_SYNTHETIC_CODE_INPUT_ID())))
	_ = input
	var astRoot pkg_rhino.Node = input.GetAstRoot(this_)
	_ = astRoot
	if !(astRoot.IsFirstChildOf(this_.jsRoot)) {
		j2gort.CheckState(false)
	}
	if !(StaticGet_Compiler_SYNTHETIC_CODE_INPUT_ID().Equals(j2gort.PtrToAny(input.GetInputId()))) {
		j2gort.CheckState(false)
	}
	if mergeContentIntoFirstInput && astRoot.HasChildren() {
		var next pkg_rhino.Node = astRoot.GetNext()
		_ = next
		j2gort.CheckNotNullNode(next)
		Compiler_CheckNotModule(next, "Cannot remove synthetic code input until modules are rewritten: %s", []any{j2gort.BoxNode(next)})
		next.AddChildrenToFront(astRoot.RemoveChildren())
		this_.ReportChangeToChangeScope(next)
	}
	astRoot.Detach()
	this_.ReportChangeToChangeScope(astRoot)
	astRoot.SetDeleted(true)
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(astRoot, this_)
	input.GetChunk().Remove(input)
	this_.inputsById.Remove_1(j2gort.PtrToAny(input.GetInputId()))
}

// Java method: Compiler.getLastPassName()
func (this_ *Compiler) GetLastPassName() string {
	return this_.lastPassName
}

// Java method: Compiler.setExternProperties(com.google.common.collect.ImmutableSet<java.lang.String> externProperties)
func (this_ *Compiler) SetExternProperties(externProperties j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = externProperties
	this_.externProperties = externProperties
}

// Java method: Compiler.getExternProperties()
func (this_ *Compiler) GetExternProperties() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.externProperties
}

// Java method: Compiler.getAccessorSummary()
func (this_ *Compiler) GetAccessorSummary() *AccessorSummary {
	return this_.accessorSummary
}

// Java method: Compiler.setAccessorSummary(com.google.javascript.jscomp.AccessorSummary summary)
func (this_ *Compiler) SetAccessorSummary(summary *AccessorSummary) {
	_ = summary
	this_.accessorSummary = j2gort.CheckNotNullPtr(summary)
}

// Java method: Compiler.loadResourceContents(java.lang.String resourceName, java.lang.String path)
func (this_ *Compiler) priv_Compiler_loadResourceContents_1556(resourceName string, path string) pkg_rhino.Node {
	_ = resourceName
	_ = path
	if this_.GetLifeCycleStage().HasColorAndSimplifiedJSDoc() {
		j2gort.CheckNotNull(this_.runtimeLibraryTypedAsts)
		var typedAstSupplier j2gort.Ext_java_util_function_Supplier = j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](j2gort.MapGetStr(this_.runtimeLibraryTypedAsts, path)))
		_ = typedAstSupplier
		return j2gort.CastNode[pkg_rhino.Node](typedAstSupplier.Get())
	}
	if !(!(this_.HasTypeCheckingRun())) {
		j2gort.CheckStateStr(false, "runtime library injected after type checking but before optimization colors")
	}
	var originalCode string = pkg_resources.ResourceLoader_LoadTextResource(j2gort.ClassLiteral("com.google.javascript.jscomp.Compiler", "Compiler"), ("js/" + j2gort.StrValString(resourceName) + ".js"))
	_ = originalCode
	var source *SourceFile = SourceFile_FromCode_2(path, originalCode)
	_ = source
	this_.priv_Compiler_addFilesToSourceMap_1558(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.ImmListOfPtr(source)))
	return this_.priv_Compiler_parseCodeHelper_1517(source)
}

// Java method: Compiler.getNodeForRuntimeCodeInsertion()
func (this_ *Compiler) priv_Compiler_getNodeForRuntimeCodeInsertion_1557() pkg_rhino.Node {
	if this_.options.GetTypedAstOutputFile() == nil {
		return this_.GetNodeForCodeInsertion(nil)
	}
	return this_.GetSynthesizedTypeSummaryInput().GetAstRoot(this_)
}

// Java method: Compiler.addComments(java.lang.String filename, java.util.List<com.google.javascript.jscomp.parsing.parser.trees.Comment> comments)
func (this_ *Compiler) AddComments(filename string, comments j2gort.Ext_java_util_List) {
	_ = filename
	_ = comments
	if !(this_.GetOptions().PreservesDetailedSourceInfo()) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "addComments may only be called in IDE mode.")))
	}
	j2gort.MapPutStr(this_.commentsPerFile, filename, comments)
}

// Java method: Compiler.getComments(java.lang.String filename)
func (this_ *Compiler) GetComments(filename string) j2gort.Ext_java_util_List {
	_ = filename
	if !(this_.GetOptions().PreservesDetailedSourceInfo()) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "getComments may only be called in IDE mode.")))
	}
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapGetStr(this_.commentsPerFile, filename))
}

// Java method: Compiler.getModuleLoader()
func (this_ *Compiler) GetModuleLoader() *ModuleLoader {
	return this_.moduleLoader
}

// Java method: Compiler.addFilesToSourceMap(java.lang.Iterable<com.google.javascript.jscomp.SourceFile> files)
func (this_ *Compiler) priv_Compiler_addFilesToSourceMap_1558(files j2gort.Ext_java_lang_Iterable) {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	this_._syncMu.Lock()
	defer this_._syncMu.Unlock()
	_ = files
	if this_.GetOptions().GetSourceMapIncludeSourcesContent() && (this_.GetSourceMap() != nil) {
		{
			cur_173 := j2gort.NewIterCursor(j2gort.CheckNotNull(files))
			for cur_173.HasNext() {
				rawElem_174 := cur_173.Next()
				file := j2gort.CastPtr[SourceFile](rawElem_174)
				_ = file
				func() {
					defer func() {
						if rec_175 := recover(); rec_175 != nil {
							jex_176 := javalang.RecoverException(rec_175)
							_ = jex_176
							if j2gort.InstanceOf[javaio.IIOException](jex_176) {
								e := j2gort.AsIface[javaio.IIOException](jex_176)
								_ = e
								panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Cannot read code of a source map's source file.", e)))
							} else {
								panic(rec_175)
							}
						}
					}()
					this_.GetSourceMap().AddSourceFile(file.GetName(), file.GetCode())
				}()
			}
		}
	}
}

// Java method: Compiler.initWebpackMap(com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> inputPathByWebpackId)
func (this_ *Compiler) InitWebpackMap(inputPathByWebpackId j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = inputPathByWebpackId
	this_.inputPathByWebpackId = inputPathByWebpackId
}

// Java method: Compiler.setDebugMessage(java.lang.String debugMessage)
func (this_ *Compiler) SetDebugMessage(debugMessage string) {
	_ = debugMessage
	this_.compilerExecutor.SetDebugMessage(debugMessage)
	this_.debugMessage = debugMessage
}

// Java method: Compiler.createCompilerExecutor()
func (this_ *Compiler) CreateCompilerExecutor() *CompilerExecutor {
	return NewCompilerExecutor()
}

// Java method: Compiler.getCompilerExecutor()
func (this_ *Compiler) GetCompilerExecutor() *CompilerExecutor {
	return this_.compilerExecutor
}

// Java static method: Compiler.toProto(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func Compiler_ToProto_1_FeatureSet_Feature(feature *pkg_parser.FeatureSet_Feature) *j2gort.Ext_com_google_javascript_jscomp_FeatureProto {
	clinit_Compiler()
	_ = feature
	var didRet_177 bool
	var retVal_178 *j2gort.Ext_com_google_javascript_jscomp_FeatureProto
	func() {
		defer func() {
			if rec_179 := recover(); rec_179 != nil {
				jex_180 := javalang.RecoverException(rec_179)
				_ = jex_180
				if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_180) {
					e := j2gort.AsIface[javalang.IIllegalArgumentException](jex_180)
					_ = e
					retVal_178 = j2gort.Ext_com_google_javascript_jscomp_FeatureProto_F_FEATURE_UNKNOWN
					didRet_177 = true
					return
				} else {
					panic(rec_179)
				}
			}
		}()
		retVal_178 = j2gort.Ext_com_google_javascript_jscomp_FeatureProto_ValueOf_1_string(("FEATURE_" + j2gort.StrValString(j2gort.EnumName(feature))))
		didRet_177 = true
		return
	}()
	if didRet_177 {
		return retVal_178
	}
	panic("unreachable")
}

// Java static method: Compiler.toProto(com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage stage)
func Compiler_ToProto_1_AbstractCompiler_LifeCycleStage(stage *AbstractCompiler_LifeCycleStage) *j2gort.Ext_com_google_javascript_jscomp_LifeCycleStageProto {
	clinit_Compiler()
	_ = stage
	switch j2gort.CheckNotNullPtr(stage).Ordinal() {
	case int32(0) /* LifeCycleStage.RAW */ :
		return j2gort.Ext_com_google_javascript_jscomp_LifeCycleStageProto_F_LIFE_CYCLE_STAGE_RAW
	case int32(1) /* LifeCycleStage.COLORS_AND_SIMPLIFIED_JSDOC */ :
		return j2gort.Ext_com_google_javascript_jscomp_LifeCycleStageProto_F_LIFE_CYCLE_STAGE_COLORS_AND_SIMPLIFIED_JSDOC
	case int32(2) /* LifeCycleStage.NORMALIZED */ :
		return j2gort.Ext_com_google_javascript_jscomp_LifeCycleStageProto_F_LIFE_CYCLE_STAGE_NORMALIZED
	case int32(3) /* LifeCycleStage.NORMALIZED_OBFUSCATED */ :
		return j2gort.Ext_com_google_javascript_jscomp_LifeCycleStageProto_F_LIFE_CYCLE_STAGE_NORMALIZED_OBFUSCATED
	}
	panic("unreachable switch expression")
}

// Java static method: Compiler.toVariableMapEntries(com.google.javascript.jscomp.VariableMap map)
func Compiler_ToVariableMapEntries(map_ *VariableMap) j2gort.Ext_java_util_List {
	clinit_Compiler()
	_ = map_
	var entries []*j2gort.Ext_com_google_javascript_jscomp_VariableMapEntryProto = []*j2gort.Ext_com_google_javascript_jscomp_VariableMapEntryProto{}
	_ = entries
	if map_ != nil {
		{
			cur_181 := j2gort.NewMapEntryCursor(map_.GetOriginalNameToNewNameMap())
			for cur_181.NextEntry() {
				entries = append(entries, j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_VariableMapEntryProto](j2gort.Ext_com_google_javascript_jscomp_VariableMapEntryProto_NewBuilder_0().SetOriginalName(cur_181.KeyStr()).SetNewName(j2gort.CastString(cur_181.GetValue())).Build()))
			}
		}
	}
	return j2gort.WrapTypedArrayList(entries)
}

// Java static method: Compiler.fromProto(java.util.List<com.google.javascript.jscomp.FeatureProto> protos)
func Compiler_FromProto_1_Ext_java_util_List(protos j2gort.Ext_java_util_List) *pkg_parser.FeatureSet {
	clinit_Compiler()
	_ = protos
	var features j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = features
	{
		cur_182 := j2gort.NewIterCursor(j2gort.CheckNotNull(protos))
		for cur_182.HasNext() {
			rawElem_183 := cur_182.Next()
			p := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_FeatureProto](rawElem_183)
			_ = p
			if p != j2gort.Ext_com_google_javascript_jscomp_FeatureProto_F_FEATURE_UNKNOWN {
				func() {
					defer func() {
						if rec_184 := recover(); rec_184 != nil {
							jex_185 := javalang.RecoverException(rec_184)
							_ = jex_185
							if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_185) {
								e := j2gort.AsIface[javalang.IIllegalArgumentException](jex_185)
								_ = e
							} else {
								panic(rec_184)
							}
						}
					}()
					features.Add_1(j2gort.PtrToAny(pkg_parser.FeatureSet_Feature_ValueOf(j2gort.JStrSubFrom(j2gort.EnumName(p), j2gort.JStrLen("FEATURE_")))))
				}()
			}
		}
	}
	return pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_Ext_java_util_Set(features)
}

// Java static method: Compiler.fromProto(com.google.javascript.jscomp.LifeCycleStageProto stage)
func Compiler_FromProto_1_Ext_com_google_javascript_jscomp_LifeCycleStageProto(stage *j2gort.Ext_com_google_javascript_jscomp_LifeCycleStageProto) *AbstractCompiler_LifeCycleStage {
	clinit_Compiler()
	_ = stage
	switch j2gort.CheckNotNullPtr(stage).Ordinal() {
	case int32(1) /* LifeCycleStageProto.LIFE_CYCLE_STAGE_RAW */ :
		return StaticGet_AbstractCompiler_LifeCycleStage_RAW()
	case int32(2) /* LifeCycleStageProto.LIFE_CYCLE_STAGE_COLORS_AND_SIMPLIFIED_JSDOC */ :
		return StaticGet_AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC()
	case int32(3) /* LifeCycleStageProto.LIFE_CYCLE_STAGE_NORMALIZED */ :
		return StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED()
	case int32(4) /* LifeCycleStageProto.LIFE_CYCLE_STAGE_NORMALIZED_OBFUSCATED */ :
		return StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED()
	default:
		return StaticGet_AbstractCompiler_LifeCycleStage_RAW()
	}
	panic("unreachable switch expression")
}

// Java static method: Compiler.fromVariableMapEntries(java.util.List<com.google.javascript.jscomp.VariableMapEntryProto> entries)
func Compiler_FromVariableMapEntries(entries j2gort.Ext_java_util_List) *VariableMap {
	clinit_Compiler()
	_ = entries
	var map_ j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = map_
	{
		cur_186 := j2gort.NewIterCursor(j2gort.CheckNotNull(entries))
		for cur_186.HasNext() {
			rawElem_187 := cur_186.Next()
			entry := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_VariableMapEntryProto](rawElem_187)
			_ = entry
			j2gort.MapPutStr(map_, entry.GetOriginalName(), j2gort.BoxString(entry.GetNewName()))
		}
	}
	return NewVariableMap(map_)
}

// Java method: Compiler.getCompilerStateProtoBuilder()
func (this_ *Compiler) GetCompilerStateProtoBuilder() *j2gort.Ext_com_google_javascript_jscomp_JSCompilerStateProto_Builder {
	var builder *j2gort.Ext_com_google_javascript_jscomp_JSCompilerStateProto_Builder = j2gort.Ext_com_google_javascript_jscomp_JSCompilerStateProto_NewBuilder_0()
	_ = builder
	{
		cur_188 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.allowableFeatures.GetFeatures()))
		for cur_188.HasNext() {
			rawElem_189 := cur_188.Next()
			feature := j2gort.CastPtr[pkg_parser.FeatureSet_Feature](rawElem_189)
			_ = feature
			builder.AddAllowableFeatures(Compiler_ToProto_1_FeatureSet_Feature(feature))
		}
	}
	builder.SetTypeCheckingHasRun(this_.typeCheckingHasRun).SetHasRegExpGlobalReferences(this_.hasRegExpGlobalReferences).SetLifeCycleStage(Compiler_ToProto_1_AbstractCompiler_LifeCycleStage(this_.GetLifeCycleStage())).SetMergedPrecompiledLibraries(this_.options.GetMergedPrecompiledLibraries())
	var chunkList []*JSChunk = this_.chunkGraph.GetIntegralChunkArrayForSerialization()
	_ = chunkList
	var chunkNameMap j2gort.OrderedMap[string, *j2gort.BoxInt] = j2gort.NewOrderedMap[string, *j2gort.BoxInt]()
	_ = chunkNameMap
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(chunkList)); i = i + 1 {
			chunkNameMap.Put(chunkList[i].GetName(), j2gort.ValueOfInt(i))
		}
	}
	for _, el_190 := range j2gort.CheckNotNullSlice(chunkList) {
		chunk := el_190
		_ = chunk
		var chunkBuilder *j2gort.Ext_com_google_javascript_jscomp_ChunkProto_Builder = j2gort.Ext_com_google_javascript_jscomp_ChunkProto_NewBuilder_0()
		_ = chunkBuilder
		chunkBuilder.SetName(chunk.GetName())
		{
			cur_191 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetDependencies()))
			for cur_191.HasNext() {
				rawElem_192 := cur_191.Next()
				dep := j2gort.CastPtr[JSChunk](rawElem_192)
				_ = dep
				var depIdx *j2gort.BoxInt = chunkNameMap.Get(dep.GetName())
				_ = depIdx
				if !(depIdx != nil) {
					j2gort.CheckStateFmt(false, "Missing dependency %s for chunk %s", j2gort.BoxString(dep.GetName()), j2gort.BoxString(chunk.GetName()))
				}
				chunkBuilder.AddDependencies(j2gort.UnboxInt(depIdx))
			}
		}
		{
			cur_193 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetInputs()))
			for cur_193.HasNext() {
				rawElem_194 := cur_193.Next()
				input := j2gort.CastPtr[CompilerInput](rawElem_194)
				_ = input
				chunkBuilder.AddInputIds(input.GetInputId().GetIdName())
			}
		}
		builder.AddChunks_1_Ext_com_google_javascript_jscomp_ChunkProto(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_ChunkProto](chunkBuilder.Build()))
	}
	builder.SetUniqueNameId(this_.uniqueNameId).AddAllUniqueIdSupplier(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.uniqueIdSupplier.ToProto())).AddAllExportedNames(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.exportedNames))
	builder.SetHasCssNames((this_.cssNames != nil))
	if this_.cssNames != nil {
		builder.AddAllCssNames(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.cssNames))
	}
	if this_.idGeneratorMap != j2gort.NullString {
		builder.SetIdGeneratorMap(this_.idGeneratorMap)
	}
	builder.SetTranspiledFiles(this_.transpiledFiles).SetIdGeneratorCurrentId(this_.crossChunkIdGenerator.GetCurrentId()).SetRunJ2ClPasses(this_.runJ2clPasses)
	{
		cur_195 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.externs))
		for cur_195.HasNext() {
			rawElem_196 := cur_195.Next()
			extern := j2gort.CastPtr[CompilerInput](rawElem_196)
			_ = extern
			builder.AddExterns(extern.GetInputId().GetIdName())
		}
	}
	builder.AddAllInjectedLibraries(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.runtimeJsLibManager.GetInjectedLibraries()))
	var lastInjectedLibrary pkg_rhino.Node = this_.runtimeJsLibManager.GetLastInjectedLibrary()
	_ = lastInjectedLibrary
	builder.SetLastInjectedLibraryIndexInFirstScript(func() int32 {
		if ((lastInjectedLibrary != 0) && (this_.jsRoot != 0)) && this_.jsRoot.HasChildren() {
			return this_.jsRoot.GetFirstChild().GetIndexOfChild(lastInjectedLibrary)
		}
		return int32(-1)
	}())
	if this_.accessorSummary != nil {
		builder.SetAccessorSummary_1_Ext_com_google_javascript_jscomp_AccessorSummaryProto(this_.accessorSummary.ToProto())
	}
	builder.SetHasStringMap((this_.stringMap != nil))
	if this_.stringMap != nil {
		builder.AddAllStringMap(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](Compiler_ToVariableMapEntries(this_.stringMap)))
	}
	builder.SetHasInstrumentationMapping((this_.instrumentationMapping != nil))
	if this_.instrumentationMapping != nil {
		builder.AddAllInstrumentationMapping(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](Compiler_ToVariableMapEntries(this_.instrumentationMapping)))
	}
	return builder
}

// Java method: Compiler.restoreFromState(com.google.javascript.jscomp.JSCompilerStateProto compilerState)
func (this_ *Compiler) RestoreFromState(compilerState *j2gort.Ext_com_google_javascript_jscomp_JSCompilerStateProto) {
	_ = compilerState
	this_.allowableFeatures = Compiler_FromProto_1_Ext_java_util_List(compilerState.GetAllowableFeaturesList())
	this_.typeCheckingHasRun = compilerState.GetTypeCheckingHasRun()
	this_.scriptNodeByFilename.Clear()
	this_.hasRegExpGlobalReferences = compilerState.GetHasRegExpGlobalReferences()
	{
		cur_197 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](compilerState.GetInjectedLibrariesList())))
		for cur_197.HasNext() {
			library := cur_197.NextStr()
			_ = library
			this_.GetRuntimeJsLibManager().RecordLibraryInjected(library)
		}
	}
	var stage *AbstractCompiler_LifeCycleStage
	_ = stage
	if compilerState.GetLifeCycleStage() == j2gort.Ext_com_google_javascript_jscomp_LifeCycleStageProto_F_LIFE_CYCLE_STAGE_RAW {
		stage = StaticGet_AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC()
	} else {
		stage = Compiler_FromProto_1_Ext_com_google_javascript_jscomp_LifeCycleStageProto(compilerState.GetLifeCycleStage())
	}
	this_.SetLifeCycleStage(stage)
	this_.GetOptions().SetMergedPrecompiledLibraries(compilerState.GetMergedPrecompiledLibraries())
	this_.uniqueNameId = compilerState.GetUniqueNameId()
	this_.uniqueIdSupplier = UniqueIdSupplier_FromProto(compilerState.GetUniqueIdSupplierList())
	this_.exportedNames.Clear()
	this_.exportedNames.AddAll_1(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](compilerState.GetExportedNamesList()))
	if compilerState.GetHasCssNames() {
		this_.cssNames = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](compilerState.GetCssNamesList()))
	} else {
		this_.cssNames = nil
	}
	this_.variableMap = nil
	this_.propertyMap = nil
	if compilerState.HasIdGeneratorMap() {
		this_.idGeneratorMap = compilerState.GetIdGeneratorMap()
	} else {
		this_.idGeneratorMap = j2gort.NullString
	}
	this_.transpiledFiles = compilerState.GetTranspiledFiles()
	this_.crossChunkIdGenerator = NewIdGenerator()
	this_.crossChunkIdGenerator.SetCurrentId(compilerState.GetIdGeneratorCurrentId())
	this_.runJ2clPasses = compilerState.GetRunJ2ClPasses()
	if compilerState.HasAccessorSummary() {
		this_.accessorSummary = AccessorSummary_FromProto_1_Ext_com_google_javascript_jscomp_AccessorSummaryProto(compilerState.GetAccessorSummary())
	} else {
		this_.accessorSummary = nil
	}
	if compilerState.GetHasStringMap() {
		this_.stringMap = Compiler_FromVariableMapEntries(compilerState.GetStringMapList())
	} else {
		this_.stringMap = nil
	}
	if compilerState.GetHasInstrumentationMapping() {
		this_.instrumentationMapping = Compiler_FromVariableMapEntries(compilerState.GetInstrumentationMappingList())
	} else {
		this_.instrumentationMapping = nil
	}
	this_.changeTracker.ResetChangeStamp()
}

// Java method: Compiler.saveState(java.io.OutputStream outputStream)
func (this_ *Compiler) SaveState(outputStream j2gort.Ext_java_io_OutputStream) {
	_ = outputStream
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		var tracer *Tracer = this_.NewTracer("serializeCompilerState")
		_ = tracer
		var nonClosingOut j2gort.Ext_java_io_OutputStream = NewCompiler_NonClosingOutputStream(outputStream)
		_ = nonClosingOut
		var compressionStream j2gort.Ext_java_io_OutputStream = nil
		_ = compressionStream
		var wrapper j2gort.Ext_java_util_function_Function = this_.options.GetStateCompressionWrapper()
		_ = wrapper
		if wrapper != nil {
			compressionStream = j2gort.AsIface[j2gort.Ext_java_io_OutputStream](wrapper.Apply(nonClosingOut))
		} else {
			compressionStream = NewFastGzipOutputStream_1(nonClosingOut)
		}
		var bufferedCompressionStream j2gort.Ext_java_io_OutputStream = j2gort.NewExt_java_io_BufferedOutputStream_2(compressionStream, (int32(64) * int32(1024)))
		_ = bufferedCompressionStream
		var builder *j2gort.Ext_com_google_javascript_jscomp_JSCompilerStateProto_Builder = this_.GetCompilerStateProtoBuilder()
		_ = builder
		j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_JSCompilerStateProto](builder.Build()).WriteDelimitedTo(bufferedCompressionStream)
		this_.StopTracer(tracer, "serializeCompilerState")
		tracer = this_.NewTracer("serializeTypedAst")
		SerializeTypedAstPass_CreateFromOutputStream(this_, bufferedCompressionStream, SerializationOptions_S_Builder().SetRuntimeLibraries(j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()).SetIncludeDebugInfo(this_.GetOptions().ShouldSerializeExtraDebugInfo()).Build()).Process_2(this_.externsRoot, this_.jsRoot)
		this_.StopTracer(tracer, "serializeTypedAst")
		bufferedCompressionStream.Close_0()
		return nil
	}))
}

// Java method: Compiler.restoreState(java.io.InputStream inputStream)
func (this_ *Compiler) RestoreState(inputStream j2gort.Ext_java_io_InputStream) {
	_ = inputStream
	this_.priv_Compiler_initWarningsGuard_1443(this_.options.GetWarningsGuard())
	this_.MaybeSetTracker()
	this_.RunInCompilerThread(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
		var tracer *Tracer = this_.NewTracer("deserializeCompilerState")
		_ = tracer
		StaticGet_Compiler_Logger().Fine_1_string("Deserializing the CompilerState")
		var nonClosingIn j2gort.Ext_java_io_InputStream = NewCompiler_NonClosingInputStream(inputStream)
		_ = nonClosingIn
		var decompressStream j2gort.Ext_java_io_InputStream = nil
		_ = decompressStream
		var wrapper j2gort.Ext_java_util_function_Function = this_.options.GetStateDecompressionWrapper()
		_ = wrapper
		if wrapper != nil {
			decompressStream = j2gort.AsIface[j2gort.Ext_java_io_InputStream](wrapper.Apply(nonClosingIn))
		} else {
			decompressStream = j2gort.NewExt_java_util_zip_GZIPInputStream_2(nonClosingIn, (int32(64) * int32(1024)))
		}
		var bufferedDecompressStream j2gort.Ext_java_io_InputStream = j2gort.NewExt_java_io_BufferedInputStream_2(decompressStream, (int32(64) * int32(1024)))
		_ = bufferedDecompressStream
		defer func() {
			bufferedDecompressStream.Close_0()
			StaticGet_Compiler_Logger().Fine_1_string("Finished deserializing CompilerState")
			this_.StopTracer(tracer, "deserializeCompilerState")
		}()
		{
			this_.priv_Compiler_deserializeCompilerState_1567(bufferedDecompressStream)
			return nil
		}
		panic("unreachable")
	}))
	if this_.tracker != nil {
		this_.tracker.UpdateAfterDeserialize(this_.jsRoot)
	}
}

// Java method: Compiler.deserializeCompilerState(java.io.InputStream inputStream)
func (this_ *Compiler) priv_Compiler_deserializeCompilerState_1567(inputStream j2gort.Ext_java_io_InputStream) {
	_ = inputStream
	var stateProto *j2gort.Ext_com_google_javascript_jscomp_JSCompilerStateProto = j2gort.Ext_com_google_javascript_jscomp_JSCompilerStateProto_ParseDelimitedFrom_1(inputStream)
	_ = stateProto
	j2gort.CheckNotNullPtr(this_.chunkGraph)
	var externFilesBuilder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = externFilesBuilder
	var codeFilesBuilder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = codeFilesBuilder
	var allInputFiles j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = allInputFiles
	{
		cur_198 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.chunkGraph.GetAllInputs()))
		for cur_198.HasNext() {
			rawElem_199 := cur_198.Next()
			input := j2gort.CastPtr[CompilerInput](rawElem_199)
			_ = input
			allInputFiles.Add_1_any(j2gort.PtrToAny(input.GetSourceFile()))
			codeFilesBuilder.Put_2(j2gort.BoxString(input.GetInputId().GetIdName()), j2gort.PtrToAny(input.GetSourceFile()))
		}
	}
	{
		cur_200 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.externs))
		for cur_200.HasNext() {
			rawElem_201 := cur_200.Next()
			extern := j2gort.CastPtr[CompilerInput](rawElem_201)
			_ = extern
			allInputFiles.Add_1_any(j2gort.PtrToAny(extern.GetSourceFile()))
			externFilesBuilder.Put_2(j2gort.BoxString(extern.GetInputId().GetIdName()), j2gort.PtrToAny(extern.GetSourceFile()))
		}
	}
	var deserializedAst TypedAstDeserializer_DeserializedAst = TypedAstDeserializer_DeserializeFullAst(this_, this_.u_SYNTHETIC_EXTERNS_FILE, j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](allInputFiles.Build()), inputStream, stateProto.GetTypeCheckingHasRun(), this_.GetOptions().GetResolveSourceMapAnnotations(), this_.GetOptions().GetParseInlineSourceMaps())
	_ = deserializedAst
	this_.RestoreFromState(stateProto)
	var chunks j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = chunks
	{
		cur_202 := j2gort.NewIterCursor(j2gort.CheckNotNull(stateProto.GetChunksList()))
		for cur_202.HasNext() {
			rawElem_203 := cur_202.Next()
			chunkProto := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_ChunkProto](rawElem_203)
			_ = chunkProto
			var chunk *JSChunk = NewJSChunk(chunkProto.GetName())
			_ = chunk
			chunks.Add_1(chunk)
		}
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < chunks.Size(); i = i + 1 {
			var chunk_2 *JSChunk = j2gort.CastPtr[JSChunk](chunks.Get(i))
			_ = chunk_2
			var chunkProto_2 *j2gort.Ext_com_google_javascript_jscomp_ChunkProto = stateProto.GetChunks(i)
			_ = chunkProto_2
			{
				cur_204 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunkProto_2.GetDependenciesList()))
				for cur_204.HasNext() {
					rawElem_205 := cur_204.Next()
					depIdx := javalang.D2I(j2gort.UnboxToDouble(rawElem_205))
					_ = depIdx
					chunk_2.AddDependency(j2gort.CastPtr[JSChunk](chunks.Get(depIdx)))
				}
			}
		}
	}
	this_.chunkGraph = NewJSChunkGraph_1_Ext_java_util_List(chunks)
	this_.externProperties = deserializedAst.GetExternProperties()
	this_.externAndJsRoot = pkg_rhino.IR_Root([]pkg_rhino.Node{pkg_rhino.IR_Root([]pkg_rhino.Node{}), pkg_rhino.IR_Root([]pkg_rhino.Node{})})
	this_.externsRoot = this_.externAndJsRoot.GetFirstChild()
	this_.jsRoot = this_.externAndJsRoot.GetLastChild()
	this_.inputsById.Clear()
	this_.externs.Clear()
	this_.colorRegistry = j2gort.CastPtr[pkg_colors.ColorRegistry](deserializedAst.GetColorRegistry().OrNull())
	this_.typedAstFilesystem = deserializedAst.GetFilesystem()
	var externFiles j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](externFilesBuilder.BuildOrThrow())
	_ = externFiles
	var codeFiles j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](codeFilesBuilder.BuildOrThrow())
	_ = codeFiles
	{
		cur_206 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](stateProto.GetExternsList())))
		for cur_206.HasNext() {
			externName := cur_206.NextStr()
			_ = externName
			{
				var iifeCond_1 bool
				{

					sr_207 := externName
					se_208 := this_.u_SYNTHETIC_EXTERNS_FILE.GetName()
					iifeCond_1 = j2gort.CheckNotNullString(sr_207) == se_208
				}
				if iifeCond_1 {
					this_.GetSynthesizedExternsInput()
					continue
				}
			}
			var externFile *SourceFile = j2gort.CastPtr[SourceFile](j2gort.MapGetStr(externFiles, externName))
			_ = externFile
			if externFile == nil {
				externFile = j2gort.CheckNotNullPtr(j2gort.CastPtr[SourceFile](j2gort.MapGetStr(codeFiles, externName)))
			}
			var input_2 *CompilerInput = NewCompilerInput_2_SourceFile_bool(externFile, true)
			_ = input_2
			var script pkg_rhino.Node = input_2.GetAstRoot(this_)
			_ = script
			this_.externsRoot.AddChildToBack(script)
			this_.priv_Compiler_putCompilerInput_1494(input_2)
			j2gort.MapPutStr(this_.scriptNodeByFilename, externFile.GetName(), j2gort.BoxNode(script))
			this_.externs.Add_1(input_2)
		}
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < chunks.Size(); i_2 = i_2 + 1 {
			var chunkProto_3 *j2gort.Ext_com_google_javascript_jscomp_ChunkProto = stateProto.GetChunks(i_2)
			_ = chunkProto_3
			{
				cur_209 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](chunkProto_3.GetInputIdsList())))
				for cur_209.HasNext() {
					inputId := cur_209.NextStr()
					_ = inputId
					var src *SourceFile = j2gort.CastPtr[SourceFile](j2gort.MapGetStr(codeFiles, inputId))
					_ = src
					if src == nil {
						if !(AbstractCompiler_IsFillFileName(inputId)) {
							j2gort.CheckStateFmt(false, "Missing %s", j2gort.BoxString(inputId))
						}
						continue
					}
					var input_3 *CompilerInput = NewCompilerInput_1(src)
					_ = input_3
					var script_2 pkg_rhino.Node = input_3.GetAstRoot(this_)
					_ = script_2
					this_.jsRoot.AddChildToBack(script_2)
					j2gort.MapPutStr(this_.scriptNodeByFilename, src.GetName(), j2gort.BoxNode(script_2))
					this_.priv_Compiler_putCompilerInput_1494(input_3)
					j2gort.CastPtr[JSChunk](chunks.Get(i_2)).Add_1_CompilerInput(input_3)
				}
			}
		}
	}
	{
		cur_210 := j2gort.NewIterCursor(j2gort.CheckNotNull(deserializedAst.GetRuntimeLibraries()))
		for cur_210.HasNext() {
			library := cur_210.NextStr()
			_ = library
			this_.runtimeJsLibManager.EnsureLibraryInjected(library, false)
		}
	}
	this_.typedAstFilesystem = nil
	if stateProto.GetLastInjectedLibraryIndexInFirstScript() != int32(-1) {
		var lastInjectedLibrary pkg_rhino.Node = this_.jsRoot.GetFirstChild().GetChildAtIndex(stateProto.GetLastInjectedLibraryIndexInFirstScript())
		_ = lastInjectedLibrary
		this_.runtimeJsLibManager.SetLastInjectedLibrary(lastInjectedLibrary)
	}
}

// Java method: Compiler.getModuleTypeByName(java.lang.String moduleName)
func (this_ *Compiler) GetModuleTypeByName(moduleName string) CompilerInput_ModuleType {
	_ = moduleName
	return this_.moduleTypesByName.Get(moduleName)
}

// Java method: Compiler.getModuleMetadataMap()
func (this_ *Compiler) GetModuleMetadataMap() *ModuleMetadataMap {
	return this_.moduleMetadataMap
}

// Java method: Compiler.setModuleMetadataMap(com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap)
func (this_ *Compiler) SetModuleMetadataMap(moduleMetadataMap *ModuleMetadataMap) {
	_ = moduleMetadataMap
	this_.moduleMetadataMap = moduleMetadataMap
}

// Java method: Compiler.getModuleMap()
func (this_ *Compiler) GetModuleMap() *ModuleMap {
	return this_.moduleMap
}

// Java method: Compiler.setModuleMap(com.google.javascript.jscomp.modules.ModuleMap moduleMap)
func (this_ *Compiler) SetModuleMap(moduleMap *ModuleMap) {
	_ = moduleMap
	this_.moduleMap = moduleMap
}

// Java static method: Compiler.resolveSibling(java.lang.String fromPath, java.lang.String toPath)
func Compiler_ResolveSibling(fromPath string, toPath string) string {
	clinit_Compiler()
	_ = fromPath
	_ = toPath
	if javalang.StringStartsWith(j2gort.CheckNotNullString(toPath), "/") {
		return toPath
	}
	var fromPathParts j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Arrays_AsList(func() []any {
		arr_211 := j2gort.JStrSplit(fromPath, "/")
		if arr_211 == nil {
			return nil
		}
		out_212 := make([]any, len(arr_211))
		for i_, el_213 := range arr_211 {
			out_212[i_] = j2gort.BoxString(el_213)
		}
		return out_212
	}())))
	_ = fromPathParts
	var toPathParts j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Arrays_AsList(func() []any {
		arr_214 := j2gort.JStrSplit(toPath, "/")
		if arr_214 == nil {
			return nil
		}
		out_215 := make([]any, len(arr_214))
		for i_, el_216 := range arr_214 {
			out_215[i_] = j2gort.BoxString(el_216)
		}
		return out_215
	}())))
	_ = toPathParts
	if !(fromPathParts.IsEmpty()) {
		fromPathParts.Remove_1_int32((fromPathParts.Size() - int32(1)))
	}
	for !(fromPathParts.IsEmpty()) && !(toPathParts.IsEmpty()) {
		if j2gort.CheckNotNullString(j2gort.ListGetStr(toPathParts, int32(0))) == "." {
			toPathParts.Remove_1_int32(int32(0))
		} else if j2gort.CheckNotNullString(j2gort.ListGetStr(toPathParts, int32(0))) == ".." {
			toPathParts.Remove_1_int32(int32(0))
			fromPathParts.Remove_1_int32((fromPathParts.Size() - int32(1)))
		} else {
			break
		}
	}
	fromPathParts.AddAll_1(toPathParts)
	return j2gort.Ext_java_lang_String_Join_2_any_Ext_java_lang_Iterable("/", fromPathParts)
}

// Java method: Compiler.resetAndIntitializeSourceMap()
func (this_ *Compiler) ResetAndIntitializeSourceMap() {
	if this_.sourceMap == nil {
		return
	}
	this_.sourceMap.Reset()
	if this_.options.GetSourceMapIncludeSourcesContent() {
		if this_.options.GetApplyInputSourceMaps() {
			{
				cur_217 := j2gort.NewMapValCursor(this_.inputSourceMaps)
				for cur_217.HasNext() {
					rawElem_218 := cur_217.Next()
					inputSourceMap := j2gort.CastPtr[SourceMapInput](rawElem_218)
					_ = inputSourceMap
					this_.priv_Compiler_addSourceMapSourceFiles_1538(inputSourceMap)
				}
			}
		}
		var allChunks j2gort.Ext_java_lang_Iterable = this_.GetChunks()
		_ = allChunks
		if allChunks != nil {
			var sourceFiles j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
			_ = sourceFiles
			{
				cur_219 := j2gort.NewIterCursor(j2gort.CheckNotNull(allChunks))
				for cur_219.HasNext() {
					rawElem_220 := cur_219.Next()
					chunk := j2gort.CastPtr[JSChunk](rawElem_220)
					_ = chunk
					{
						cur_221 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunk.GetInputs()))
						for cur_221.HasNext() {
							rawElem_222 := cur_221.Next()
							input := j2gort.CastPtr[CompilerInput](rawElem_222)
							_ = input
							var sourceFile *SourceFile = input.GetSourceFile()
							_ = sourceFile
							if sourceFile.IsStubSourceFileForAlreadyProvidedInput() {
								continue
							}
							sourceFiles.Add_1(j2gort.PtrToAny(input.GetSourceFile()))
						}
					}
				}
			}
			this_.priv_Compiler_addFilesToSourceMap_1558(sourceFiles)
		}
	}
}

// Java static method: Compiler.checkNotModule(com.google.javascript.rhino.Node script, java.lang.String msg, java.lang.Object[] args)
func Compiler_CheckNotModule(script pkg_rhino.Node, msg string, args []any) pkg_rhino.Node {
	clinit_Compiler()
	_ = script
	_ = msg
	_ = args
	if !(script.IsScript()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(script))
	}
	if !(script.HasOneChild()) {
		return script
	}
	if !(!(script.GetFirstChild().IsModuleBody())) {
		j2gort.CheckStateFmt(false, j2gort.StrValString(msg), j2gort.SliceToAny(args))
	}
	return script
}

// Java method: Compiler.maybeLogPruningAnalysis(com.google.javascript.jscomp.JSChunkGraph.DependencyManagementResult result)
func (this_ *Compiler) priv_Compiler_maybeLogPruningAnalysis_1569(result *JSChunkGraph_DependencyManagementResult) {
	_ = result
	if (this_.tracker == nil) || !(this_.options.GetTracerMode().DoPruningAnalysis()) {
		return
	}
	var analysis *PruningAnalysis = PruningAnalysis_Create(result.Sorter(), j2gort.AsIface[j2gort.Ext_java_util_Collection](result.EntryPointInputs()))
	_ = analysis
	var pruningResult *PruningAnalysis_Result = analysis.Analyze()
	_ = pruningResult
	this_.tracker.SetPruningAnalysisSummary(("\nTransitive dependency count per entry point:\n" + j2gort.StrVal(pruningResult.EntryPointDependencyCount()) + "\n\nTop bottleneck dependencies:\n" + j2gort.StrVal(pruningResult.BottleneckBlame())))
}

// Java class: com.google.javascript.jscomp.Compiler.ExternalSourceLoader
type Compiler_ExternalSourceLoader struct {
	idHash_ atomic.Int32
}

func (this_ *Compiler_ExternalSourceLoader) BindSelf_(s any) {
	_ = s
}
func (this_ *Compiler_ExternalSourceLoader) IsJavaNil_() bool { return this_ == nil }
func (*Compiler_ExternalSourceLoader) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$ExternalSourceLoader", "ExternalSourceLoader"
}
func (this_ *Compiler_ExternalSourceLoader) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Compiler_ExternalSourceLoader) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Compiler.ExternalSourceLoader", "com.google.javascript.jscomp.Compiler$ExternalSourceLoader", func(v_ any) bool {
		return j2gort.InstanceOf[*Compiler_ExternalSourceLoader](v_)
	})
}

func (this_ *Compiler_ExternalSourceLoader) initFields_Compiler_ExternalSourceLoader() {
	_ = this_
}

// Java constructor: ExternalSourceLoader()
func NewCompiler_ExternalSourceLoader() *Compiler_ExternalSourceLoader {
	this_ := &Compiler_ExternalSourceLoader{}
	InitCompiler_ExternalSourceLoader(this_)
	return this_
}

func InitCompiler_ExternalSourceLoader(this_ *Compiler_ExternalSourceLoader) {
	this_.initFields_Compiler_ExternalSourceLoader()
}

// Java method: ExternalSourceLoader.loadSource(java.lang.String filename)
func (this_ *Compiler_ExternalSourceLoader) LoadSource(filename string) *SourceFile {
	_ = filename
	panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Cannot load without a valid loader.")))
}

// Java class: com.google.javascript.jscomp.Compiler.ResolvedSourceMap
type Compiler_ResolvedSourceMap struct {
	originalPath  string `gson:"originalPath"`
	sourceMapPath string `gson:"sourceMapPath"`
	relativePath  string `gson:"relativePath"`
	idHash_       atomic.Int32
}

func (this_ *Compiler_ResolvedSourceMap) BindSelf_(s any) {
	_ = s
}
func (this_ *Compiler_ResolvedSourceMap) IsJavaNil_() bool { return this_ == nil }
func (*Compiler_ResolvedSourceMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$ResolvedSourceMap", "ResolvedSourceMap"
}
func (this_ *Compiler_ResolvedSourceMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Compiler_ResolvedSourceMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Compiler.ResolvedSourceMap", "com.google.javascript.jscomp.Compiler$ResolvedSourceMap", func(v_ any) bool {
		return j2gort.InstanceOf[*Compiler_ResolvedSourceMap](v_)
	})
}

func (this_ *Compiler_ResolvedSourceMap) initFields_Compiler_ResolvedSourceMap() {
	_ = this_
	this_.originalPath = j2gort.NullString
	this_.sourceMapPath = j2gort.NullString
	this_.relativePath = j2gort.NullString
}

// Java constructor: ResolvedSourceMap()
func NewCompiler_ResolvedSourceMap() *Compiler_ResolvedSourceMap {
	this_ := &Compiler_ResolvedSourceMap{}
	InitCompiler_ResolvedSourceMap(this_)
	return this_
}

func InitCompiler_ResolvedSourceMap(this_ *Compiler_ResolvedSourceMap) {
	this_.initFields_Compiler_ResolvedSourceMap()
}

// Java polymorphic class interface: com.google.javascript.jscomp.Compiler.SeenSetLicenseTracker
type Compiler_SeenSetLicenseTracker interface {
	Implements_jscomp_Compiler_SeenSetLicenseTracker()
	IsJavaNil_() bool
	GetCompiler_SeenSetLicenseTrackerBase() *Compiler_SeenSetLicenseTrackerBase
	CodePrinter_LicenseTracker
	ShouldUseLicenseInfo(pkg_rhino.Node) bool
}

func ToIface_Compiler_SeenSetLicenseTracker[T interface {
	Compiler_SeenSetLicenseTracker
	comparable
}](v T) Compiler_SeenSetLicenseTracker {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.Compiler.SeenSetLicenseTracker
type Compiler_SeenSetLicenseTrackerBase struct {
	self_                  Compiler_SeenSetLicenseTracker
	compiler               AbstractCompiler          `gson:"compiler"`
	globallyUniqueLicenses j2gort.OrderedSet[string] `gson:"globallyUniqueLicenses"`
	currentlySeenLicenses  j2gort.OrderedSet[string] `gson:"currentlySeenLicenses"`
	lastSeenFile           string                    `gson:"lastSeenFile"`
	idHash_                atomic.Int32
}

func (this_ *Compiler_SeenSetLicenseTrackerBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(Compiler_SeenSetLicenseTracker)
}
func (this_ *Compiler_SeenSetLicenseTrackerBase) IsJavaNil_() bool { return this_ == nil }
func (*Compiler_SeenSetLicenseTrackerBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$SeenSetLicenseTracker", "SeenSetLicenseTracker"
}
func (this_ *Compiler_SeenSetLicenseTrackerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Compiler.SeenSetLicenseTracker", "com.google.javascript.jscomp.Compiler$SeenSetLicenseTracker", func(v_ any) bool {
		return j2gort.InstanceOf[Compiler_SeenSetLicenseTracker](v_)
	})
}
func (b *Compiler_SeenSetLicenseTrackerBase) GetCompiler_SeenSetLicenseTrackerBase() *Compiler_SeenSetLicenseTrackerBase {
	return b
}
func (*Compiler_SeenSetLicenseTrackerBase) Implements_jscomp_Compiler_SeenSetLicenseTracker() {}
func (*Compiler_SeenSetLicenseTrackerBase) Implements_jscomp_CodePrinter_LicenseTracker()     {}

func (this_ *Compiler_SeenSetLicenseTrackerBase) initFields_Compiler_SeenSetLicenseTracker() {
	_ = this_
	this_.globallyUniqueLicenses = j2gort.NewOrderedSet[string]()
	this_.currentlySeenLicenses = j2gort.NewOrderedSet[string]()

}

func InitCompiler_SeenSetLicenseTracker(this_ *Compiler_SeenSetLicenseTrackerBase, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Compiler_SeenSetLicenseTracker()
	this_.compiler = compiler
}

// Java method: SeenSetLicenseTracker.shouldUseLicenseInfo(com.google.javascript.rhino.Node node)
func (this_ *Compiler_SeenSetLicenseTrackerBase) ShouldUseLicenseInfo(node pkg_rhino.Node) bool {
	_ = node
	return (!(node.IsRoot()) && !(node.IsScript()))
}

// Java method: SeenSetLicenseTracker.trackLicensesForNode(com.google.javascript.rhino.Node node)
func (this_ *Compiler_SeenSetLicenseTrackerBase) TrackLicensesForNode(node pkg_rhino.Node) {
	_ = node
	if !(this_.self_.ShouldUseLicenseInfo(node)) {
		return
	}
	var file string = node.GetSourceFileName()
	_ = file
	if file == j2gort.NullString {
		return
	}
	{
		var iifeCond_1 bool
		{

			sr_223 := file
			se_224 := this_.lastSeenFile
			iifeCond_1 = j2gort.CheckNotNullString(sr_223) == se_224
		}
		if iifeCond_1 {
			return
		} else {
			this_.lastSeenFile = file
		}
	}
	var license string = Compiler_GetLicenseForFile(this_.compiler, file)
	_ = license
	if license == j2gort.NullString {
		return
	}
	if this_.globallyUniqueLicenses.Contains(license) {
		return
	}
	this_.currentlySeenLicenses.Add(license)
}

// Java method: SeenSetLicenseTracker.emitLicenses()
func (this_ *Compiler_SeenSetLicenseTrackerBase) EmitLicenses() j2gort.Ext_com_google_common_collect_ImmutableSet {
	var rv j2gort.Ext_com_google_common_collect_ImmutableSet = this_.currentlySeenLicenses.ToImmutableSet()
	_ = rv
	this_.globallyUniqueLicenses.AddAllSet(this_.currentlySeenLicenses)
	this_.currentlySeenLicenses.Clear()
	return rv
}

// Java class: com.google.javascript.jscomp.Compiler.ScriptNodeLicensesOnlyTracker
type Compiler_ScriptNodeLicensesOnlyTracker struct {
	Compiler_SeenSetLicenseTrackerBase
}

func (this_ *Compiler_ScriptNodeLicensesOnlyTracker) BindSelf_(s any) {
	_ = s
	this_.Compiler_SeenSetLicenseTrackerBase.BindSelf_(s)
}
func (this_ *Compiler_ScriptNodeLicensesOnlyTracker) IsJavaNil_() bool { return this_ == nil }
func (*Compiler_ScriptNodeLicensesOnlyTracker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$ScriptNodeLicensesOnlyTracker", "ScriptNodeLicensesOnlyTracker"
}
func (this_ *Compiler_ScriptNodeLicensesOnlyTracker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Compiler_SeenSetLicenseTrackerBase.IdentityHash_()
}
func (*Compiler_ScriptNodeLicensesOnlyTracker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Compiler.ScriptNodeLicensesOnlyTracker", "com.google.javascript.jscomp.Compiler$ScriptNodeLicensesOnlyTracker", func(v_ any) bool {
		return j2gort.InstanceOf[*Compiler_ScriptNodeLicensesOnlyTracker](v_)
	})
}
func (*Compiler_ScriptNodeLicensesOnlyTracker) Implements_jscomp_Compiler_SeenSetLicenseTracker() {}
func (*Compiler_ScriptNodeLicensesOnlyTracker) Implements_jscomp_CodePrinter_LicenseTracker()     {}

func (this_ *Compiler_ScriptNodeLicensesOnlyTracker) initFields_Compiler_ScriptNodeLicensesOnlyTracker() {
	_ = this_
}

// Java constructor: ScriptNodeLicensesOnlyTracker(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCompiler_ScriptNodeLicensesOnlyTracker(compiler AbstractCompiler) *Compiler_ScriptNodeLicensesOnlyTracker {
	this_ := &Compiler_ScriptNodeLicensesOnlyTracker{}
	this_.BindSelf_(this_)
	InitCompiler_ScriptNodeLicensesOnlyTracker(this_, compiler)
	return this_
}

func InitCompiler_ScriptNodeLicensesOnlyTracker(this_ *Compiler_ScriptNodeLicensesOnlyTracker, compiler AbstractCompiler) {
	_ = compiler
	InitCompiler_SeenSetLicenseTracker(&this_.Compiler_SeenSetLicenseTrackerBase, compiler)
	this_.initFields_Compiler_ScriptNodeLicensesOnlyTracker()
}

// Java method: ScriptNodeLicensesOnlyTracker.shouldUseLicenseInfo(com.google.javascript.rhino.Node node)
func (this_ *Compiler_ScriptNodeLicensesOnlyTracker) ShouldUseLicenseInfo(node pkg_rhino.Node) bool {
	_ = node
	return (node.IsRoot() || node.IsScript())
}

// Java class: com.google.javascript.jscomp.Compiler.SingleBinaryLicenseTracker
type Compiler_SingleBinaryLicenseTracker struct {
	Compiler_SeenSetLicenseTrackerBase
}

func (this_ *Compiler_SingleBinaryLicenseTracker) BindSelf_(s any) {
	_ = s
	this_.Compiler_SeenSetLicenseTrackerBase.BindSelf_(s)
}
func (this_ *Compiler_SingleBinaryLicenseTracker) IsJavaNil_() bool { return this_ == nil }
func (*Compiler_SingleBinaryLicenseTracker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$SingleBinaryLicenseTracker", "SingleBinaryLicenseTracker"
}
func (this_ *Compiler_SingleBinaryLicenseTracker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Compiler_SeenSetLicenseTrackerBase.IdentityHash_()
}
func (*Compiler_SingleBinaryLicenseTracker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Compiler.SingleBinaryLicenseTracker", "com.google.javascript.jscomp.Compiler$SingleBinaryLicenseTracker", func(v_ any) bool {
		return j2gort.InstanceOf[*Compiler_SingleBinaryLicenseTracker](v_)
	})
}
func (*Compiler_SingleBinaryLicenseTracker) Implements_jscomp_Compiler_SeenSetLicenseTracker() {}
func (*Compiler_SingleBinaryLicenseTracker) Implements_jscomp_CodePrinter_LicenseTracker()     {}

func (this_ *Compiler_SingleBinaryLicenseTracker) initFields_Compiler_SingleBinaryLicenseTracker() {
	_ = this_
}

// Java constructor: SingleBinaryLicenseTracker(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCompiler_SingleBinaryLicenseTracker(compiler AbstractCompiler) *Compiler_SingleBinaryLicenseTracker {
	this_ := &Compiler_SingleBinaryLicenseTracker{}
	this_.BindSelf_(this_)
	InitCompiler_SingleBinaryLicenseTracker(this_, compiler)
	return this_
}

func InitCompiler_SingleBinaryLicenseTracker(this_ *Compiler_SingleBinaryLicenseTracker, compiler AbstractCompiler) {
	_ = compiler
	InitCompiler_SeenSetLicenseTracker(&this_.Compiler_SeenSetLicenseTrackerBase, compiler)
	this_.initFields_Compiler_SingleBinaryLicenseTracker()
}

// Java method: SingleBinaryLicenseTracker.shouldUseLicenseInfo(com.google.javascript.rhino.Node node)
func (this_ *Compiler_SingleBinaryLicenseTracker) ShouldUseLicenseInfo(node pkg_rhino.Node) bool {
	_ = node
	return (!(node.IsRoot()) && !(node.IsScript()))
}

// Java class: com.google.javascript.jscomp.Compiler.ChunkGraphAwareLicenseTracker
type Compiler_ChunkGraphAwareLicenseTracker struct {
	compiler                            AbstractCompiler                                      `gson:"compiler"`
	licensesFromChunks                  j2gort.OrderedMap[*JSChunk, j2gort.Ext_java_util_Set] `gson:"licensesFromChunks"`
	currentChunkLicencesInTDeps         j2gort.OrderedSet[string]                             `gson:"currentChunkLicencesInTDeps"`
	lastSeenFile                        string                                                `gson:"lastSeenFile"`
	licensesNewInCurrentFile            j2gort.OrderedSet[string]                             `gson:"licensesNewInCurrentFile"`
	licensesNewInCurrentChunk           j2gort.Ext_java_util_Set                              `gson:"licensesNewInCurrentChunk"`
	log                                 pkg_diagnostic.LogFile                                `gson:"log"`
	currentChunk                        *JSChunk                                              `gson:"currentChunk"`
	idHash_                             atomic.Int32
	haveInitializedCurrentChunkLicenses bool `gson:"haveInitializedCurrentChunkLicenses"`
}

func (this_ *Compiler_ChunkGraphAwareLicenseTracker) BindSelf_(s any) {
	_ = s
}
func (this_ *Compiler_ChunkGraphAwareLicenseTracker) IsJavaNil_() bool { return this_ == nil }
func (*Compiler_ChunkGraphAwareLicenseTracker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$ChunkGraphAwareLicenseTracker", "ChunkGraphAwareLicenseTracker"
}
func (this_ *Compiler_ChunkGraphAwareLicenseTracker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Compiler_ChunkGraphAwareLicenseTracker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Compiler.ChunkGraphAwareLicenseTracker", "com.google.javascript.jscomp.Compiler$ChunkGraphAwareLicenseTracker", func(v_ any) bool {
		return j2gort.InstanceOf[*Compiler_ChunkGraphAwareLicenseTracker](v_)
	})
}
func (*Compiler_ChunkGraphAwareLicenseTracker) Implements_jscomp_CodePrinter_LicenseTracker() {}

func (this_ *Compiler_ChunkGraphAwareLicenseTracker) initFields_Compiler_ChunkGraphAwareLicenseTracker() {
	_ = this_
	this_.licensesFromChunks = j2gort.NewOrderedMap[*JSChunk, j2gort.Ext_java_util_Set]()

	this_.currentChunkLicencesInTDeps = j2gort.NewOrderedSet[string]()

	this_.licensesNewInCurrentFile = j2gort.NewOrderedSet[string]()
	this_.licensesNewInCurrentChunk = j2gort.NewExt_java_util_LinkedHashSet()
	this_.log = pkg_diagnostic.LogFile_CreateNoOp()
}

// Java constructor: ChunkGraphAwareLicenseTracker(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCompiler_ChunkGraphAwareLicenseTracker(compiler AbstractCompiler) *Compiler_ChunkGraphAwareLicenseTracker {
	this_ := &Compiler_ChunkGraphAwareLicenseTracker{}
	InitCompiler_ChunkGraphAwareLicenseTracker(this_, compiler)
	return this_
}

func InitCompiler_ChunkGraphAwareLicenseTracker(this_ *Compiler_ChunkGraphAwareLicenseTracker, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Compiler_ChunkGraphAwareLicenseTracker()
	this_.compiler = compiler
}

// Java method: ChunkGraphAwareLicenseTracker.setLogFile(com.google.javascript.jscomp.diagnostic.LogFile file)
func (this_ *Compiler_ChunkGraphAwareLicenseTracker) SetLogFile(file pkg_diagnostic.LogFile) {
	_ = file
	this_.log = file
}

// Java method: ChunkGraphAwareLicenseTracker.setCurrentChunkContext(com.google.javascript.jscomp.JSChunk chunk)
func (this_ *Compiler_ChunkGraphAwareLicenseTracker) SetCurrentChunkContext(chunk *JSChunk) {
	_ = chunk
	if this_.log.IsLogging() {
		this_.log.Log_2("Initializing licenses from deps for chunk %s", []any{j2gort.PtrToAny(chunk)})
	}
	this_.currentChunk = chunk
	this_.currentChunkLicencesInTDeps.Clear()
	this_.licensesNewInCurrentFile.Clear()
	this_.haveInitializedCurrentChunkLicenses = false
	if this_.licensesFromChunks.ContainsKey(chunk) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Visiting a chunk more than once is not allowed.")))
	}
	this_.licensesNewInCurrentChunk = j2gort.NewExt_java_util_LinkedHashSet()
	this_.licensesFromChunks.Put(chunk, this_.licensesNewInCurrentChunk)
}

// Java method: ChunkGraphAwareLicenseTracker.trackLicensesForNode(com.google.javascript.rhino.Node node)
func (this_ *Compiler_ChunkGraphAwareLicenseTracker) TrackLicensesForNode(node pkg_rhino.Node) {
	_ = node
	if node.IsRoot() || node.IsScript() {
		return
	}
	var sourceFile string = node.GetSourceFileName()
	_ = sourceFile
	if sourceFile == j2gort.NullString {
		return
	}
	if j2gort.CheckNotNullString(this_.lastSeenFile) == sourceFile {
		return
	} else {
		this_.lastSeenFile = sourceFile
	}
	if !(this_.haveInitializedCurrentChunkLicenses) {
		{
			cur_225 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.currentChunk.GetAllDependencies()))
			for cur_225.HasNext() {
				rawElem_226 := cur_225.Next()
				depChunk := j2gort.CastPtr[JSChunk](rawElem_226)
				_ = depChunk
				if this_.log.IsLogging() {
					this_.log.Log_2("Chunk %s depends on chunk %s", []any{j2gort.PtrToAny(this_.currentChunk), j2gort.PtrToAny(depChunk)})
				}
				if !(this_.licensesFromChunks.ContainsKey(depChunk)) {
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Chunk aware license analysis error - analysis out of order.")))
				}
				j2gort.OrderedSetAddAllStr(&this_.currentChunkLicencesInTDeps, this_.licensesFromChunks.Get(depChunk))
			}
		}
		this_.haveInitializedCurrentChunkLicenses = true
	}
	var license string = Compiler_GetLicenseForFile(this_.compiler, node.GetSourceFileName())
	_ = license
	if license == j2gort.NullString {
		return
	}
	if this_.currentChunkLicencesInTDeps.Contains(license) {
		return
	}
	if j2gort.SetContainsStr(this_.licensesNewInCurrentChunk, license) {
		return
	}
	var newlyAddedLicense bool = this_.licensesNewInCurrentFile.Add(license)
	_ = newlyAddedLicense
	if !(newlyAddedLicense) {
		if this_.log.IsLogging() {
			this_.log.Log_2("Chunk %s already depends on license\n==================\n%s\n==================\n\n", []any{j2gort.PtrToAny(this_.currentChunk), j2gort.BoxString(license)})
		}
	} else {
		if this_.log.IsLogging() {
			this_.log.Log_2("Chunk %s has new license:\n==================\n%s\n==================\n\n", []any{j2gort.PtrToAny(this_.currentChunk), j2gort.BoxString(license)})
		}
	}
}

// Java method: ChunkGraphAwareLicenseTracker.emitLicenses()
func (this_ *Compiler_ChunkGraphAwareLicenseTracker) EmitLicenses() j2gort.Ext_com_google_common_collect_ImmutableSet {
	var licensesNewInFile j2gort.Ext_com_google_common_collect_ImmutableSet = this_.licensesNewInCurrentFile.ToImmutableSet()
	_ = licensesNewInFile
	this_.licensesNewInCurrentChunk.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](licensesNewInFile))
	this_.licensesNewInCurrentFile.Clear()
	return licensesNewInFile
}

// Java class: com.google.javascript.jscomp.Compiler.CodeBuilder
type Compiler_CodeBuilder struct {
	sb        *j2gort.Ext_java_lang_StringBuilder `gson:"sb"`
	lineCount int32                               `gson:"lineCount"`
	colCount  int32                               `gson:"colCount"`
	idHash_   atomic.Int32
}

func (this_ *Compiler_CodeBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *Compiler_CodeBuilder) IsJavaNil_() bool { return this_ == nil }
func (*Compiler_CodeBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$CodeBuilder", "CodeBuilder"
}
func (this_ *Compiler_CodeBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Compiler_CodeBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Compiler.CodeBuilder", "com.google.javascript.jscomp.Compiler$CodeBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*Compiler_CodeBuilder](v_)
	})
}

func (this_ *Compiler_CodeBuilder) initFields_Compiler_CodeBuilder() {
	_ = this_
	this_.sb = j2gort.NewExt_java_lang_StringBuilder()

}

// Java constructor: CodeBuilder()
func NewCompiler_CodeBuilder() *Compiler_CodeBuilder {
	this_ := &Compiler_CodeBuilder{}
	InitCompiler_CodeBuilder(this_)
	return this_
}

func InitCompiler_CodeBuilder(this_ *Compiler_CodeBuilder) {
	this_.initFields_Compiler_CodeBuilder()
}

// Java method: CodeBuilder.reset()
func (this_ *Compiler_CodeBuilder) Reset() {
	this_.sb.SetLength(int32(0))
}

// Java method: CodeBuilder.append(java.lang.String str)
func (this_ *Compiler_CodeBuilder) Append(str string) *Compiler_CodeBuilder {
	_ = str
	this_.sb.Append_1_string(str)
	var index int32 = int32(-1)
	_ = index
	var lastIndex int32 = index
	_ = lastIndex
	for (func() int32 {
		var asgn_227 int32 = j2gort.JStrIndexOfCpFrom(str, int32(uint16(10)), (index + int32(1)))
		index = asgn_227
		return asgn_227
	}()) >= int32(0) {
		this_.lineCount = this_.lineCount + 1
		lastIndex = index
	}
	if lastIndex == int32(-1) {
		this_.colCount = (this_.colCount + j2gort.JStrLen(str))
	} else {
		this_.colCount = (j2gort.JStrLen(str) - (lastIndex + int32(1)))
	}
	return this_
}

// Java method: CodeBuilder.toString()
func (this_ *Compiler_CodeBuilder) String() string {
	return this_.sb.String()
}

// Java method: CodeBuilder.getLength()
func (this_ *Compiler_CodeBuilder) GetLength() int32 {
	return this_.sb.Length()
}

// Java method: CodeBuilder.getLineIndex()
func (this_ *Compiler_CodeBuilder) GetLineIndex() int32 {
	return this_.lineCount
}

// Java method: CodeBuilder.getColumnIndex()
func (this_ *Compiler_CodeBuilder) GetColumnIndex() int32 {
	return this_.colCount
}

// Java method: CodeBuilder.endsWith(java.lang.String suffix)
func (this_ *Compiler_CodeBuilder) EndsWith(suffix string) bool {
	_ = suffix
	return ((this_.sb.Length() > j2gort.JStrLen(suffix)) && func() bool {
		sr_228 := suffix
		se_229 := this_.sb.Substring_1((this_.sb.Length() - j2gort.JStrLen(suffix)))
		return j2gort.CheckNotNullString(sr_228) == se_229
	}())
}

// Java class: com.google.javascript.jscomp.Compiler.NonClosingOutputStream
type Compiler_NonClosingOutputStream struct {
	j2gort.Ext_java_io_FilterOutputStreamBase
	idHash_ atomic.Int32
}

func (this_ *Compiler_NonClosingOutputStream) BindSelf_(s any) {
	_ = s
	this_.Ext_java_io_FilterOutputStreamBase.BindExtSelf_(s)
}
func (this_ *Compiler_NonClosingOutputStream) IsJavaNil_() bool { return this_ == nil }
func (*Compiler_NonClosingOutputStream) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$NonClosingOutputStream", "NonClosingOutputStream"
}
func (this_ *Compiler_NonClosingOutputStream) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Compiler.NonClosingOutputStream", "com.google.javascript.jscomp.Compiler$NonClosingOutputStream", func(v_ any) bool {
		return j2gort.InstanceOf[*Compiler_NonClosingOutputStream](v_)
	})
}
func (*Compiler_NonClosingOutputStream) Implements_j2gort_Ext_java_io_FilterOutputStream() {}
func (*Compiler_NonClosingOutputStream) Implements_j2gort_Ext_java_io_OutputStream()       {}
func (*Compiler_NonClosingOutputStream) Implements_j2gort_Ext_java_io_Closeable()          {}
func (*Compiler_NonClosingOutputStream) Implements_j2gort_Ext_java_io_Flushable()          {}
func (*Compiler_NonClosingOutputStream) Implements_j2gort_Ext_java_lang_AutoCloseable()    {}

func (this_ *Compiler_NonClosingOutputStream) initFields_Compiler_NonClosingOutputStream() {
	_ = this_
}

// Java constructor: NonClosingOutputStream(java.io.OutputStream out)
func NewCompiler_NonClosingOutputStream(out j2gort.Ext_java_io_OutputStream) *Compiler_NonClosingOutputStream {
	this_ := &Compiler_NonClosingOutputStream{}
	this_.BindSelf_(this_)
	InitCompiler_NonClosingOutputStream(this_, out)
	return this_
}

func InitCompiler_NonClosingOutputStream(this_ *Compiler_NonClosingOutputStream, out j2gort.Ext_java_io_OutputStream) {
	_ = out
	j2gort.InitExt_java_io_FilterOutputStream(&this_.Ext_java_io_FilterOutputStreamBase, out)
	this_.initFields_Compiler_NonClosingOutputStream()
}

// Java method: NonClosingOutputStream.close()
func (this_ *Compiler_NonClosingOutputStream) Close_0() {
	this_.Flush()
}

// Java class: com.google.javascript.jscomp.Compiler.NonClosingInputStream
type Compiler_NonClosingInputStream struct {
	j2gort.Ext_java_io_FilterInputStreamBase
	idHash_ atomic.Int32
}

func (this_ *Compiler_NonClosingInputStream) BindSelf_(s any) {
	_ = s
	this_.Ext_java_io_FilterInputStreamBase.BindExtSelf_(s)
}
func (this_ *Compiler_NonClosingInputStream) IsJavaNil_() bool { return this_ == nil }
func (*Compiler_NonClosingInputStream) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$NonClosingInputStream", "NonClosingInputStream"
}
func (this_ *Compiler_NonClosingInputStream) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Compiler.NonClosingInputStream", "com.google.javascript.jscomp.Compiler$NonClosingInputStream", func(v_ any) bool {
		return j2gort.InstanceOf[*Compiler_NonClosingInputStream](v_)
	})
}
func (*Compiler_NonClosingInputStream) Implements_j2gort_Ext_java_io_FilterInputStream() {}
func (*Compiler_NonClosingInputStream) Implements_j2gort_Ext_java_io_InputStream()       {}
func (*Compiler_NonClosingInputStream) Implements_j2gort_Ext_java_io_Closeable()         {}
func (*Compiler_NonClosingInputStream) Implements_j2gort_Ext_java_lang_AutoCloseable()   {}

func (this_ *Compiler_NonClosingInputStream) initFields_Compiler_NonClosingInputStream() {
	_ = this_
}

// Java constructor: NonClosingInputStream(java.io.InputStream in)
func NewCompiler_NonClosingInputStream(in j2gort.Ext_java_io_InputStream) *Compiler_NonClosingInputStream {
	this_ := &Compiler_NonClosingInputStream{}
	this_.BindSelf_(this_)
	InitCompiler_NonClosingInputStream(this_, in)
	return this_
}

func InitCompiler_NonClosingInputStream(this_ *Compiler_NonClosingInputStream, in j2gort.Ext_java_io_InputStream) {
	_ = in
	j2gort.InitExt_java_io_FilterInputStream(&this_.Ext_java_io_FilterInputStreamBase, in)
	this_.initFields_Compiler_NonClosingInputStream()
}

// Java method: NonClosingInputStream.close()
func (this_ *Compiler_NonClosingInputStream) Close_0() {
}

// Java class:
type Compiler__Anon_12 struct {
	j2gort.Ext_java_util_AbstractSetBase
	outer_ *Compiler
}

func (this_ *Compiler__Anon_12) BindSelf_(s any) {
	_ = s
	this_.Ext_java_util_AbstractSetBase.BindExtSelf_(s)
}
func (this_ *Compiler__Anon_12) IsJavaNil_() bool { return this_ == nil }
func (*Compiler__Anon_12) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$1", ""
}
func (this_ *Compiler__Anon_12) SetOuter_(o *Compiler)                         { this_.outer_ = o }
func (this_ *Compiler__Anon_12) GetOuter_() *Compiler                          { return this_.outer_ }
func (*Compiler__Anon_12) Implements_j2gort_Ext_java_util_AbstractSet()        {}
func (*Compiler__Anon_12) Implements_j2gort_Ext_java_util_AbstractCollection() {}
func (*Compiler__Anon_12) Implements_j2gort_Ext_java_util_Set()                {}
func (*Compiler__Anon_12) Implements_j2gort_Ext_java_util_Collection()         {}
func (*Compiler__Anon_12) Implements_j2gort_Ext_java_lang_Iterable()           {}

func (this_ *Compiler__Anon_12) initFields_Compiler__Anon_12() {
	_ = this_
}

func NewCompiler__Anon_12(outer_ *Compiler) *Compiler__Anon_12 {
	this_ := &Compiler__Anon_12{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	j2gort.InitExt_java_util_AbstractSet(&this_.Ext_java_util_AbstractSetBase)
	this_.initFields_Compiler__Anon_12()
	return this_
}

// Java method: .contains(java.lang.Object o)
func (this_ *Compiler__Anon_12) Contains(o any) bool {
	_ = o
	return true
}

// Java method: .add(java.lang.String e)
func (this_ *Compiler__Anon_12) Add_1(arg_230 any) bool {
	e := j2gort.CastString(arg_230)
	_ = e
	return false
}

// Java method: .iterator()
func (this_ *Compiler__Anon_12) Iterator() any {
	return j2gort.Ext_java_util_Collections_EmptyIterator()
}

// Java method: .size()
func (this_ *Compiler__Anon_12) Size() int32 {
	return int32(0)
}

// Java class:
type Compiler__Anon_13 struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer_        *Compiler
	cap_0_builder j2gort.Ext_java_util_stream_Stream_Builder
}

func (this_ *Compiler__Anon_13) BindSelf_(s any) {
	_ = s
}
func (this_ *Compiler__Anon_13) IsJavaNil_() bool { return this_ == nil }
func (*Compiler__Anon_13) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Compiler$2", ""
}
func (this_ *Compiler__Anon_13) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*Compiler__Anon_13) HasIdentityEquals_()                                       {}
func (this_ *Compiler__Anon_13) SetOuter_(o *Compiler)                               { this_.outer_ = o }
func (this_ *Compiler__Anon_13) GetOuter_() *Compiler                                { return this_.outer_ }
func (*Compiler__Anon_13) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {}
func (*Compiler__Anon_13) Implements_jscomp_NodeTraversal_Callback()                 {}

func (this_ *Compiler__Anon_13) initFields_Compiler__Anon_13() {
	_ = this_
}

func NewCompiler__Anon_13(outer_ *Compiler, cap_p_builder j2gort.Ext_java_util_stream_Stream_Builder) *Compiler__Anon_13 {
	this_ := &Compiler__Anon_13{}
	this_.outer_ = outer_
	this_.cap_0_builder = cap_p_builder
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_Compiler__Anon_13()
	return this_
}

// Java method: .shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Compiler__Anon_13) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsStatement(n) {
		this_.cap_0_builder.Add(j2gort.BoxNode(n))
		return false
	} else {
		return true
	}
}
func (*Compiler__Anon_13) NoOpVisit_()                 {}
func (*Compiler__Anon_13) CanReleaseSyntacticScopes_() {}

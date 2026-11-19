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
// Port of com/google/javascript/jscomp/CompilerOptions.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CompilerOptions
type CompilerOptions struct {
	emitUseStrict                                j2gort.Ext_com_google_common_base_Optional                   `gson:"emitUseStrict"`
	outputFeatureSet                             j2gort.Ext_com_google_common_base_Optional                   `gson:"outputFeatureSet"`
	experimentalOutputFeatureSet                 j2gort.Ext_com_google_common_base_Optional                   `gson:"experimentalOutputFeatureSet"`
	languageOutIsDefaultStrict                   j2gort.Ext_com_google_common_base_Optional                   `gson:"languageOutIsDefaultStrict"`
	browserFeaturesetYear                        *CompilerOptions_BrowserFeaturesetYear                       `gson:"browserFeaturesetYear"`
	typedAstOutputFile                           j2gort.Ext_java_nio_file_Path                                `gson:"typedAstOutputFile"`
	dependencyOptions                            *DependencyOptions                                           `gson:"dependencyOptions"`
	messageBundle                                MessageBundle                                                `gson:"messageBundle"`
	extraAnnotationNames                         j2gort.Ext_com_google_common_collect_ImmutableSet            `gson:"extraAnnotationNames"`
	renamePrefix                                 string                                                       `gson:"renamePrefix"`
	renamePrefixNamespace                        string                                                       `gson:"renamePrefixNamespace"`
	inputSourceMaps                              j2gort.Ext_com_google_common_collect_ImmutableMap            `gson:"inputSourceMaps"`
	inputVariableMap                             *VariableMap                                                 `gson:"inputVariableMap"`
	inputPropertyMap                             *VariableMap                                                 `gson:"inputPropertyMap"`
	nameGenerator                                NameGenerator                                                `gson:"nameGenerator"`
	tcProjectId                                  string                                                       `gson:"tcProjectId"`
	codingConvention                             CodingConvention                                             `gson:"codingConvention"`
	syntheticBlockStartMarker                    string                                                       `gson:"syntheticBlockStartMarker"`
	syntheticBlockEndMarker                      string                                                       `gson:"syntheticBlockEndMarker"`
	locale                                       string                                                       `gson:"locale"`
	j2clMinifierPruningManifest                  string                                                       `gson:"j2clMinifierPruningManifest"`
	stripTypes                                   j2gort.Ext_com_google_common_collect_ImmutableSet            `gson:"stripTypes"`
	stripNameSuffixes                            j2gort.Ext_com_google_common_collect_ImmutableSet            `gson:"stripNameSuffixes"`
	stripNamePrefixes                            j2gort.Ext_com_google_common_collect_ImmutableSet            `gson:"stripNamePrefixes"`
	customPasses                                 j2gort.Ext_com_google_common_collect_SetMultimap             `gson:"-"`
	defineReplacements                           *j2gort.Ext_java_util_LinkedHashMap                          `gson:"defineReplacements"`
	cssRenamingMap                               CssRenamingMap                                               `gson:"cssRenamingMap"`
	cssRenamingSkiplist                          j2gort.Ext_com_google_common_collect_ImmutableSet            `gson:"cssRenamingSkiplist"`
	idGenerators                                 j2gort.Ext_com_google_common_collect_ImmutableMap            `gson:"idGenerators"`
	xidHashFunction                              Xid_HashFunction                                             `gson:"xidHashFunction"`
	chunkIdHashFunction                          Xid_HashFunction                                             `gson:"chunkIdHashFunction"`
	idGeneratorsMapSerialized                    string                                                       `gson:"idGeneratorsMapSerialized"`
	replaceStringsFunctionDescriptions           j2gort.Ext_com_google_common_collect_ImmutableList           `gson:"replaceStringsFunctionDescriptions"`
	replaceStringsPlaceholderToken               string                                                       `gson:"replaceStringsPlaceholderToken"`
	propertiesThatMustDisambiguate               j2gort.Ext_com_google_common_collect_ImmutableSet            `gson:"propertiesThatMustDisambiguate"`
	moduleRoots                                  j2gort.Ext_com_google_common_collect_ImmutableList           `gson:"moduleRoots"`
	forceLibraryInjection                        j2gort.Ext_com_google_common_collect_ImmutableList           `gson:"forceLibraryInjection"`
	unusedImportsToRemove                        j2gort.Ext_com_google_common_collect_ImmutableSet            `gson:"unusedImportsToRemove"`
	inputDelimiter                               string                                                       `gson:"inputDelimiter"`
	debugLogDirectory                            j2gort.Ext_java_nio_file_Path                                `gson:"debugLogDirectory"`
	debugLogFilter                               string                                                       `gson:"debugLogFilter"`
	stateCompressionWrapper                      j2gort.Ext_java_util_function_Function                       `gson:"stateCompressionWrapper"`
	stateDecompressionWrapper                    j2gort.Ext_java_util_function_Function                       `gson:"stateDecompressionWrapper"`
	filesToPrintAfterEachPassRegexList           j2gort.Ext_java_util_List                                    `gson:"filesToPrintAfterEachPassRegexList"`
	chunksToPrintAfterEachPassRegexList          j2gort.Ext_java_util_List                                    `gson:"chunksToPrintAfterEachPassRegexList"`
	qnameUsesToPrintAfterEachPassList            j2gort.Ext_java_util_List                                    `gson:"qnameUsesToPrintAfterEachPassList"`
	tracerOutput                                 j2gort.Ext_java_nio_file_Path                                `gson:"tracerOutput"`
	errorFormat                                  *ErrorFormat                                                 `gson:"errorFormat"`
	warningsGuard                                *ComposeWarningsGuard                                        `gson:"warningsGuard"`
	externExportsPath                            string                                                       `gson:"externExportsPath"`
	extraReportGenerators                        j2gort.Ext_java_util_List                                    `gson:"extraReportGenerators"`
	sourceMapOutputPath                          string                                                       `gson:"sourceMapOutputPath"`
	sourceMapDetailLevel                         *SourceMap_DetailLevel                                       `gson:"sourceMapDetailLevel"`
	sourceMapFormat                              *SourceMap_Format                                            `gson:"sourceMapFormat"`
	sourceMapLocationMappings                    j2gort.Ext_com_google_common_collect_ImmutableList           `gson:"sourceMapLocationMappings"`
	outputCharset                                j2gort.Ext_java_nio_charset_Charset                          `gson:"-"`
	errorHandler                                 ErrorHandler                                                 `gson:"-"`
	productionInstrumentationArrayName           string                                                       `gson:"productionInstrumentationArrayName"`
	conformanceConfigs                           []*j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig `gson:"conformanceConfigs"`
	conformanceRemoveRegexFromPath               j2gort.Ext_com_google_common_base_Optional                   `gson:"conformanceRemoveRegexFromPath"`
	isStrictModeInput                            j2gort.Ext_com_google_common_base_Optional                   `gson:"isStrictModeInput"`
	browserResolverPrefixReplacements            j2gort.Ext_com_google_common_collect_ImmutableMap            `gson:"browserResolverPrefixReplacements"`
	pathEscaper                                  *ModuleLoader_PathEscaper                                    `gson:"pathEscaper"`
	packageJsonEntryNames                        j2gort.Ext_java_util_List                                    `gson:"packageJsonEntryNames"`
	dynamicImportAlias                           string                                                       `gson:"dynamicImportAlias"`
	unknownDefinesToIgnore                       j2gort.Ext_com_google_common_collect_ImmutableList           `gson:"unknownDefinesToIgnore"`
	enableZonesDefineName                        string                                                       `gson:"enableZonesDefineName"`
	zoneInputPattern                             *j2gort.Ext_java_util_regex_Pattern                          `gson:"zoneInputPattern"`
	numParallelThreads                           int32                                                        `gson:"numParallelThreads"`
	maxFunctionSizeAfterInlining                 int32                                                        `gson:"maxFunctionSizeAfterInlining"`
	optimizationLoopMaxIterations                int32                                                        `gson:"optimizationLoopMaxIterations"`
	summaryDetailLevel                           int32                                                        `gson:"summaryDetailLevel"`
	lineLengthThreshold                          int32                                                        `gson:"lineLengthThreshold"`
	idHash_                                      atomic.Int32
	languageIn                                   CompilerOptions_LanguageMode                           `gson:"languageIn"`
	environment                                  CompilerOptions_Environment                            `gson:"environment"`
	instrumentForCoverageOnly                    bool                                                   `gson:"instrumentForCoverageOnly"`
	mergedPrecompiledLibraries                   bool                                                   `gson:"mergedPrecompiledLibraries"`
	inferConsts                                  bool                                                   `gson:"inferConsts"`
	assumeStrictThis                             bool                                                   `gson:"assumeStrictThis"`
	preserveDetailedSourceInfo                   bool                                                   `gson:"preserveDetailedSourceInfo"`
	preserveNonJSDocComments                     bool                                                   `gson:"preserveNonJSDocComments"`
	continueAfterErrors                          bool                                                   `gson:"continueAfterErrors"`
	incrementalCheckMode                         CompilerOptions_IncrementalCheckMode                   `gson:"incrementalCheckMode"`
	parseJsDocDocumentation                      pkg_parsing.Config_JsDocParsing                        `gson:"parseJsDocDocumentation"`
	printExterns                                 bool                                                   `gson:"printExterns"`
	inferTypes                                   bool                                                   `gson:"inferTypes"`
	skipNonTranspilationPasses                   bool                                                   `gson:"skipNonTranspilationPasses"`
	devMode                                      CompilerOptions_DevMode                                `gson:"devMode"`
	checkDeterminism                             bool                                                   `gson:"checkDeterminism"`
	strictMessageReplacement                     bool                                                   `gson:"strictMessageReplacement"`
	checkSymbols                                 bool                                                   `gson:"checkSymbols"`
	checkSuspiciousCode                          bool                                                   `gson:"checkSuspiciousCode"`
	checkTypes                                   bool                                                   `gson:"checkTypes"`
	foldConstants                                bool                                                   `gson:"foldConstants"`
	deadAssignmentElimination                    bool                                                   `gson:"deadAssignmentElimination"`
	deadPropertyAssignmentElimination            pkg_base.Tri                                           `gson:"deadPropertyAssignmentElimination"`
	inlineConstantVars                           bool                                                   `gson:"inlineConstantVars"`
	inlineProperties                             bool                                                   `gson:"inlineProperties"`
	crossChunkCodeMotion                         bool                                                   `gson:"crossChunkCodeMotion"`
	crossChunkCodeMotionNoStubMethods            bool                                                   `gson:"crossChunkCodeMotionNoStubMethods"`
	parentChunkCanSeeSymbolsDeclaredInChildren   bool                                                   `gson:"parentChunkCanSeeSymbolsDeclaredInChildren"`
	coalesceVariableNames                        bool                                                   `gson:"coalesceVariableNames"`
	optimizeLetAndConst                          bool                                                   `gson:"optimizeLetAndConst"`
	assumeGlobalScopeIsIsolated                  bool                                                   `gson:"assumeGlobalScopeIsIsolated"`
	crossChunkMethodMotion                       bool                                                   `gson:"crossChunkMethodMotion"`
	inlineGetters                                bool                                                   `gson:"inlineGetters"`
	inlineVariables                              bool                                                   `gson:"inlineVariables"`
	inlineLocalVariables                         bool                                                   `gson:"inlineLocalVariables"`
	flowSensitiveInlineVariables                 bool                                                   `gson:"flowSensitiveInlineVariables"`
	smartNameRemoval                             bool                                                   `gson:"smartNameRemoval"`
	extractPrototypeMemberDeclarations           CompilerOptions_ExtractPrototypeMemberDeclarationsMode `gson:"extractPrototypeMemberDeclarations"`
	removeUnusedPrototypeProperties              bool                                                   `gson:"removeUnusedPrototypeProperties"`
	removeUnusedClassProperties                  bool                                                   `gson:"removeUnusedClassProperties"`
	removeUnusedVars                             bool                                                   `gson:"removeUnusedVars"`
	removeUnusedLocalVars                        bool                                                   `gson:"removeUnusedLocalVars"`
	collapseVariableDeclarations                 bool                                                   `gson:"collapseVariableDeclarations"`
	collapseAnonymousFunctions                   bool                                                   `gson:"collapseAnonymousFunctions"`
	aliasStringsMode                             CompilerOptions_AliasStringsMode                       `gson:"aliasStringsMode"`
	outputJsStringUsage                          bool                                                   `gson:"outputJsStringUsage"`
	convertToDottedProperties                    bool                                                   `gson:"convertToDottedProperties"`
	rewriteFunctionExpressions                   bool                                                   `gson:"rewriteFunctionExpressions"`
	optimizeCalls                                bool                                                   `gson:"optimizeCalls"`
	optimizeESClassConstructors                  bool                                                   `gson:"optimizeESClassConstructors"`
	useTypesForLocalOptimization                 bool                                                   `gson:"useTypesForLocalOptimization"`
	useSizeHeuristicToStopOptimizationLoop       bool                                                   `gson:"useSizeHeuristicToStopOptimizationLoop"`
	variableRenaming                             VariableRenamingPolicy                                 `gson:"variableRenaming"`
	propertyRenaming                             PropertyRenamingPolicy                                 `gson:"propertyRenaming"`
	propertyRenamingOnlyCompilationMode          bool                                                   `gson:"propertyRenamingOnlyCompilationMode"`
	labelRenaming                                bool                                                   `gson:"labelRenaming"`
	reserveRawExports                            bool                                                   `gson:"reserveRawExports"`
	preferStableNames                            bool                                                   `gson:"preferStableNames"`
	generatePseudoNames                          bool                                                   `gson:"generatePseudoNames"`
	renamePrefixNamespaceAssumeCrossChunkNames   bool                                                   `gson:"renamePrefixNamespaceAssumeCrossChunkNames"`
	optimizeLocalAccessForGlobalSymbolNamespace  CompilerOptions_OptimizeLocalAccess                    `gson:"optimizeLocalAccessForGlobalSymbolNamespace"`
	collapsePropertiesLevel                      CompilerOptions_PropertyCollapseLevel                  `gson:"collapsePropertiesLevel"`
	collapseObjectLiterals                       bool                                                   `gson:"collapseObjectLiterals"`
	devirtualizeMethods                          bool                                                   `gson:"devirtualizeMethods"`
	computeFunctionSideEffects                   bool                                                   `gson:"computeFunctionSideEffects"`
	disambiguateProperties                       bool                                                   `gson:"disambiguateProperties"`
	ambiguateProperties                          bool                                                   `gson:"ambiguateProperties"`
	exportTestFunctions                          bool                                                   `gson:"exportTestFunctions"`
	replaceMessagesWithChromeI18n                bool                                                   `gson:"replaceMessagesWithChromeI18n"`
	doLateLocalization                           bool                                                   `gson:"doLateLocalization"`
	markAsCompiled                               bool                                                   `gson:"markAsCompiled"`
	closurePass                                  bool                                                   `gson:"closurePass"`
	preserveClosurePrimitives                    bool                                                   `gson:"preserveClosurePrimitives"`
	angularPass                                  bool                                                   `gson:"angularPass"`
	polymerPass                                  bool                                                   `gson:"polymerPass"`
	chromePass                                   bool                                                   `gson:"chromePass"`
	j2clPassMode                                 CompilerOptions_J2clPassMode                           `gson:"j2clPassMode"`
	j2clMinifierEnabled                          bool                                                   `gson:"j2clMinifierEnabled"`
	removeAbstractMethods                        bool                                                   `gson:"removeAbstractMethods"`
	removeClosureAsserts                         bool                                                   `gson:"removeClosureAsserts"`
	removeJ2clAsserts                            bool                                                   `gson:"removeJ2clAsserts"`
	gatherCssNames                               bool                                                   `gson:"gatherCssNames"`
	tweakProcessing                              CompilerOptions_TweakProcessing                        `gson:"tweakProcessing"`
	rewriteGlobalDeclarationsForTryCatchWrapping bool                                                   `gson:"rewriteGlobalDeclarationsForTryCatchWrapping"`
	checksOnly                                   bool                                                   `gson:"checksOnly"`
	outputJs                                     CompilerOptions_OutputJs                               `gson:"outputJs"`
	generateExports                              bool                                                   `gson:"generateExports"`
	exportLocalPropertyDefinitions               bool                                                   `gson:"exportLocalPropertyDefinitions"`
	replaceIdGenerators                          bool                                                   `gson:"replaceIdGenerators"`
	validateRequireInliningAnnotation            pkg_base.Tri                                           `gson:"validateRequireInliningAnnotation"`
	processCommonJSModules                       bool                                                   `gson:"processCommonJSModules"`
	rewritePolyfills                             bool                                                   `gson:"rewritePolyfills"`
	isolatePolyfills                             bool                                                   `gson:"isolatePolyfills"`
	injectPolyfillsNewerThan                     CompilerOptions_LanguageMode                           `gson:"injectPolyfillsNewerThan"`
	es6SubclassTranspilation                     CompilerOptions_Es6SubclassTranspilation               `gson:"es6SubclassTranspilation"`
	instrumentAsyncContext                       bool                                                   `gson:"instrumentAsyncContext"`
	runtimeLibraryMode                           RuntimeJsLibManager_RuntimeLibraryMode                 `gson:"runtimeLibraryMode"`
	assumeForwardDeclaredForMissingTypes         bool                                                   `gson:"assumeForwardDeclaredForMissingTypes"`
	preserveTypeAnnotations                      bool                                                   `gson:"preserveTypeAnnotations"`
	gentsMode                                    bool                                                   `gson:"gentsMode"`
	prettyPrint                                  bool                                                   `gson:"prettyPrint"`
	lineBreak                                    bool                                                   `gson:"lineBreak"`
	printInputDelimiter                          bool                                                   `gson:"printInputDelimiter"`
	serializeExtraDebugInfo                      bool                                                   `gson:"serializeExtraDebugInfo"`
	quoteKeywordProperties                       bool                                                   `gson:"quoteKeywordProperties"`
	preferSingleQuotes                           bool                                                   `gson:"preferSingleQuotes"`
	trustedStrings                               bool                                                   `gson:"trustedStrings"`
	printSourceAfterEachPass                     bool                                                   `gson:"printSourceAfterEachPass"`
	tracer                                       CompilerOptions_TracerMode                             `gson:"tracer"`
	colorizeErrorOutput                          bool                                                   `gson:"colorizeErrorOutput"`
	useOriginalNamesInOutput                     bool                                                   `gson:"useOriginalNamesInOutput"`
	shouldAlwaysGatherSourceMapInfo              bool                                                   `gson:"shouldAlwaysGatherSourceMapInfo"`
	parseInlineSourceMaps                        bool                                                   `gson:"parseInlineSourceMaps"`
	applyInputSourceMaps                         bool                                                   `gson:"applyInputSourceMaps"`
	resolveSourceMapAnnotations                  bool                                                   `gson:"resolveSourceMapAnnotations"`
	sourceMapIncludeSourcesContent               bool                                                   `gson:"sourceMapIncludeSourcesContent"`
	protectHiddenSideEffects                     bool                                                   `gson:"protectHiddenSideEffects"`
	assumeGettersArePure                         bool                                                   `gson:"assumeGettersArePure"`
	assumePropertiesAreStaticallyAnalyzable      bool                                                   `gson:"assumePropertiesAreStaticallyAnalyzable"`
	assumeStaticInheritanceIsNotUsed             bool                                                   `gson:"assumeStaticInheritanceIsNotUsed"`
	instrumentForCoverageOption                  CompilerOptions_InstrumentOption                       `gson:"instrumentForCoverageOption"`
	conformanceReportingMode                     CompilerOptions_ConformanceReportingMode               `gson:"conformanceReportingMode"`
	wrapGoogModulesForWhitespaceOnly             bool                                                   `gson:"wrapGoogModulesForWhitespaceOnly"`
	printConfig                                  bool                                                   `gson:"printConfig"`
	rewriteModulesBeforeTypechecking             bool                                                   `gson:"rewriteModulesBeforeTypechecking"`
	enableModuleRewriting                        bool                                                   `gson:"enableModuleRewriting"`
	moduleResolutionMode                         ModuleLoader_ResolutionMode                            `gson:"moduleResolutionMode"`
	allowDynamicImport                           bool                                                   `gson:"allowDynamicImport"`
	chunkOutputType                              CompilerOptions_ChunkOutputType                        `gson:"chunkOutputType"`
	inlineFunctionsLevel                         CompilerOptions_Reach                                  `gson:"inlineFunctionsLevel"`
	es6ModuleTranspilation                       CompilerOptions_Es6ModuleTranspilation                 `gson:"es6ModuleTranspilation"`
}

func (this_ *CompilerOptions) BindSelf_(s any) {
	_ = s
}
func (this_ *CompilerOptions) IsJavaNil_() bool { return this_ == nil }
func (*CompilerOptions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions", "CompilerOptions"
}
func (this_ *CompilerOptions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CompilerOptions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions", "com.google.javascript.jscomp.CompilerOptions", func(v_ any) bool {
		return j2gort.InstanceOf[*CompilerOptions](v_)
	})
}

func StaticGet_CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD() int32 { return int32(500) }

var static_CompilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CompilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CompilerOptions()
	return static_CompilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS
}
func StaticSet_CompilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CompilerOptions()
	static_CompilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS = v
}

var static_CompilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CompilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CompilerOptions()
	return static_CompilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS
}
func StaticSet_CompilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CompilerOptions()
	static_CompilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS = v
}

var static_CompilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_CompilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CompilerOptions()
	return static_CompilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS
}
func StaticSet_CompilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_CompilerOptions()
	static_CompilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS = v
}
func StaticGet_CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING() int32 { return int32(-1) }

var clinitDone_CompilerOptions uint32
var clinitState_CompilerOptions j2gort.ClinitState

func clinit_CompilerOptions() {
	if atomic.LoadUint32(&clinitDone_CompilerOptions) != 0 {
		return
	}
	clinitSlow_CompilerOptions()
}

func clinitSlow_CompilerOptions() {
	if !j2gort.EnterClinit(&clinitState_CompilerOptions, &clinitDone_CompilerOptions) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CompilerOptions, &clinitDone_CompilerOptions)
	static_CompilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_primitives_Chars_AsList(javalang.StringToUTF16("ABCDEFGHIJKLMNOPQRSTUVWXYZ$"))))
	static_CompilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_2(j2gort.ValueOfChar(uint16(95)), j2gort.ValueOfChar(uint16(36)))
	static_CompilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_1(j2gort.ValueOfChar(uint16(36)))
}

func (this_ *CompilerOptions) initFields_CompilerOptions() {
	_ = this_
	this_.emitUseStrict = j2gort.Ext_com_google_common_base_Optional_Absent()
	this_.outputFeatureSet = j2gort.Ext_com_google_common_base_Optional_Absent()
	this_.experimentalOutputFeatureSet = j2gort.Ext_com_google_common_base_Optional_Absent()
	this_.languageOutIsDefaultStrict = j2gort.Ext_com_google_common_base_Optional_Absent()

	this_.inferConsts = true

	this_.incrementalCheckMode = StaticGet_CompilerOptions_IncrementalCheckMode_OFF()
	this_.parseJsDocDocumentation = pkg_parsing.StaticGet_Config_JsDocParsing_TYPES_ONLY()
	this_.dependencyOptions = DependencyOptions_None()

	this_.numParallelThreads = int32(1)
	this_.useSizeHeuristicToStopOptimizationLoop = true
	this_.renamePrefix = j2gort.NullString
	this_.renamePrefixNamespace = j2gort.NullString

	this_.optimizeLocalAccessForGlobalSymbolNamespace = StaticGet_CompilerOptions_OptimizeLocalAccess_DISABLED()
	this_.tcProjectId = j2gort.NullString
	this_.syntheticBlockStartMarker = j2gort.NullString
	this_.syntheticBlockEndMarker = j2gort.NullString
	this_.locale = j2gort.NullString
	this_.j2clMinifierEnabled = true
	this_.j2clMinifierPruningManifest = j2gort.NullString
	this_.removeJ2clAsserts = true
	this_.replaceIdGenerators = true
	this_.idGeneratorsMapSerialized = j2gort.NullString
	this_.replaceStringsPlaceholderToken = j2gort.NullString

	this_.moduleRoots = j2gort.ImmListOfStr1("./")

	this_.es6SubclassTranspilation = StaticGet_CompilerOptions_Es6SubclassTranspilation_CONCISE_UNSAFE()

	this_.forceLibraryInjection = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	this_.runtimeLibraryMode = StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_INJECT()

	this_.inputDelimiter = "// Input %num%"
	this_.debugLogFilter = j2gort.NullString

	this_.filesToPrintAfterEachPassRegexList = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	this_.chunksToPrintAfterEachPassRegexList = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	this_.qnameUsesToPrintAfterEachPassList = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	this_.warningsGuard = NewComposeWarningsGuard_1_WarningsGuard([]WarningsGuard{})
	this_.summaryDetailLevel = int32(1)
	this_.lineLengthThreshold = int32(500)

	this_.externExportsPath = j2gort.NullString
	this_.extraReportGenerators = j2gort.NewExt_java_util_ArrayList()
	this_.sourceMapOutputPath = j2gort.NullString

	this_.sourceMapDetailLevel = StaticGet_SourceMap_DetailLevel_ALL()
	this_.sourceMapFormat = StaticGet_SourceMap_Format_DEFAULT()
	this_.parseInlineSourceMaps = true

	this_.resolveSourceMapAnnotations = true
	this_.sourceMapLocationMappings = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()

	this_.assumeGettersArePure = true
	this_.assumePropertiesAreStaticallyAnalyzable = true
	this_.assumeStaticInheritanceIsNotUsed = true
	this_.productionInstrumentationArrayName = j2gort.NullString
	this_.conformanceConfigs = []*j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig{}
	this_.conformanceReportingMode = StaticGet_CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG()
	this_.conformanceRemoveRegexFromPath = j2gort.Ext_com_google_common_base_Optional_Of(j2gort.PtrToAny(j2gort.Ext_java_util_regex_Pattern_Compile_1("^((.*/)?google3/)?(/?(blaze|bazel)-out/[^/]+/(bin|genfiles)/)?")))
	this_.wrapGoogModulesForWhitespaceOnly = true

	this_.isStrictModeInput = j2gort.Ext_com_google_common_base_Optional_Absent()

	this_.enableModuleRewriting = true

	this_.dynamicImportAlias = j2gort.NullString
	this_.unknownDefinesToIgnore = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	this_.enableZonesDefineName = j2gort.NullString

	this_.es6ModuleTranspilation = StaticGet_CompilerOptions_Es6ModuleTranspilation_COMPILE()
}

// Java constructor: CompilerOptions()
func NewCompilerOptions() *CompilerOptions {
	clinit_CompilerOptions()
	this_ := &CompilerOptions{}
	InitCompilerOptions(this_)
	return this_
}

func InitCompilerOptions(this_ *CompilerOptions) {
	this_.initFields_CompilerOptions()
	this_.languageIn = StaticGet_CompilerOptions_LanguageMode_STABLE_IN()
	this_.environment = StaticGet_CompilerOptions_Environment_BROWSER()
	this_.browserResolverPrefixReplacements = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	this_.moduleResolutionMode = StaticGet_ModuleLoader_ResolutionMode_BROWSER()
	this_.packageJsonEntryNames = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfStr3("browser", "module", "main"))
	this_.pathEscaper = StaticGet_ModuleLoader_PathEscaper_ESCAPE()
	this_.rewriteModulesBeforeTypechecking = false
	this_.enableModuleRewriting = true
	this_.skipNonTranspilationPasses = false
	this_.devMode = StaticGet_CompilerOptions_DevMode_OFF()
	this_.checkDeterminism = false
	this_.checkSymbols = false
	this_.checkSuspiciousCode = false
	this_.checkTypes = false
	this_.computeFunctionSideEffects = false
	this_.extraAnnotationNames = nil
	this_.foldConstants = false
	this_.coalesceVariableNames = false
	this_.optimizeLetAndConst = true
	this_.assumeGlobalScopeIsIsolated = false
	this_.deadAssignmentElimination = false
	this_.deadPropertyAssignmentElimination = pkg_base.StaticGet_Tri_UNKNOWN()
	this_.inlineConstantVars = false
	this_.inlineFunctionsLevel = StaticGet_CompilerOptions_Reach_NONE()
	this_.maxFunctionSizeAfterInlining = int32(-1)
	this_.assumeStrictThis = false
	this_.inlineProperties = false
	this_.crossChunkCodeMotion = false
	this_.parentChunkCanSeeSymbolsDeclaredInChildren = false
	this_.crossChunkMethodMotion = false
	this_.inlineGetters = false
	this_.inlineVariables = false
	this_.inlineLocalVariables = false
	this_.smartNameRemoval = false
	this_.extractPrototypeMemberDeclarations = StaticGet_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF()
	this_.removeUnusedPrototypeProperties = false
	this_.removeUnusedClassProperties = false
	this_.removeUnusedVars = false
	this_.removeUnusedLocalVars = false
	this_.collapseVariableDeclarations = false
	this_.collapseAnonymousFunctions = false
	this_.aliasStringsMode = StaticGet_CompilerOptions_AliasStringsMode_NONE()
	this_.outputJsStringUsage = false
	this_.convertToDottedProperties = false
	this_.rewriteFunctionExpressions = false
	this_.validateRequireInliningAnnotation = pkg_base.StaticGet_Tri_UNKNOWN()
	this_.variableRenaming = StaticGet_VariableRenamingPolicy_OFF()
	this_.propertyRenaming = StaticGet_PropertyRenamingPolicy_OFF()
	this_.propertyRenamingOnlyCompilationMode = false
	this_.labelRenaming = false
	this_.generatePseudoNames = false
	this_.preferStableNames = false
	this_.renamePrefix = j2gort.NullString
	this_.collapsePropertiesLevel = StaticGet_CompilerOptions_PropertyCollapseLevel_NONE()
	this_.collapseObjectLiterals = false
	this_.devirtualizeMethods = false
	this_.disambiguateProperties = false
	this_.ambiguateProperties = false
	this_.exportTestFunctions = false
	this_.nameGenerator = NewDefaultNameGenerator()
	this_.syntheticBlockStartMarker = j2gort.NullString
	this_.syntheticBlockEndMarker = j2gort.NullString
	this_.locale = j2gort.NullString
	this_.doLateLocalization = false
	this_.markAsCompiled = false
	this_.closurePass = false
	this_.preserveClosurePrimitives = false
	this_.angularPass = false
	this_.polymerPass = false
	this_.j2clPassMode = StaticGet_CompilerOptions_J2clPassMode_AUTO()
	this_.j2clMinifierEnabled = true
	this_.removeAbstractMethods = false
	this_.removeClosureAsserts = false
	this_.stripTypes = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	this_.stripNameSuffixes = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	this_.stripNamePrefixes = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	this_.customPasses = nil
	this_.defineReplacements = j2gort.NewExt_java_util_LinkedHashMap()
	this_.tweakProcessing = StaticGet_CompilerOptions_TweakProcessing_OFF()
	this_.rewriteGlobalDeclarationsForTryCatchWrapping = false
	this_.checksOnly = false
	this_.outputJs = StaticGet_CompilerOptions_OutputJs_NORMAL()
	this_.generateExports = true
	this_.exportLocalPropertyDefinitions = true
	this_.cssRenamingMap = nil
	this_.cssRenamingSkiplist = nil
	this_.idGenerators = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	this_.replaceStringsFunctionDescriptions = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	this_.replaceStringsPlaceholderToken = ""
	this_.propertiesThatMustDisambiguate = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	this_.inputSourceMaps = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	this_.instrumentForCoverageOption = StaticGet_CompilerOptions_InstrumentOption_NONE()
	this_.productionInstrumentationArrayName = ""
	this_.preserveTypeAnnotations = false
	this_.gentsMode = false
	this_.printInputDelimiter = false
	this_.prettyPrint = false
	this_.lineBreak = false
	this_.tracer = StaticGet_CompilerOptions_TracerMode_OFF()
	this_.colorizeErrorOutput = false
	this_.errorFormat = StaticGet_ErrorFormat_FULL()
	this_.chunkOutputType = StaticGet_CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE()
	this_.errorHandler = nil
	this_.printSourceAfterEachPass = false
	this_.strictMessageReplacement = false
}

// Java static method: CompilerOptions.getAngularPropertyReservedFirstChars()
func CompilerOptions_GetAngularPropertyReservedFirstChars() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CompilerOptions()
	return StaticGet_CompilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS()
}

// Java static method: CompilerOptions.getPolymerPropertyReservedFirstChars()
func CompilerOptions_GetPolymerPropertyReservedFirstChars() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_CompilerOptions()
	return StaticGet_CompilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS()
}

// Java method: CompilerOptions.shouldRunCrossChunkCodeMotion()
func (this_ *CompilerOptions) ShouldRunCrossChunkCodeMotion() bool {
	return this_.crossChunkCodeMotion
}

// Java method: CompilerOptions.shouldRunCrossChunkMethodMotion()
func (this_ *CompilerOptions) ShouldRunCrossChunkMethodMotion() bool {
	return this_.crossChunkMethodMotion
}

// Java method: CompilerOptions.getSourceMapOutputPath()
func (this_ *CompilerOptions) GetSourceMapOutputPath() string {
	return this_.sourceMapOutputPath
}

// Java method: CompilerOptions.shouldGatherSourceMapInfo()
func (this_ *CompilerOptions) ShouldGatherSourceMapInfo() bool {
	return (this_.shouldAlwaysGatherSourceMapInfo || !(j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(this_.sourceMapOutputPath)))
}

// Java method: CompilerOptions.setAlwaysGatherSourceMapInfo(boolean shouldAlwaysGatherSourceMapInfo)
func (this_ *CompilerOptions) SetAlwaysGatherSourceMapInfo(shouldAlwaysGatherSourceMapInfo bool) {
	_ = shouldAlwaysGatherSourceMapInfo
	this_.shouldAlwaysGatherSourceMapInfo = shouldAlwaysGatherSourceMapInfo
}

// Java method: CompilerOptions.getBrowserFeaturesetYearObject()
func (this_ *CompilerOptions) GetBrowserFeaturesetYearObject() *CompilerOptions_BrowserFeaturesetYear {
	return this_.browserFeaturesetYear
}

// Java method: CompilerOptions.setBrowserFeaturesetYear(int year)
func (this_ *CompilerOptions) SetBrowserFeaturesetYear_1_int32(year int32) {
	_ = year
	this_.SetBrowserFeaturesetYear_1_CompilerOptions_BrowserFeaturesetYear(CompilerOptions_BrowserFeaturesetYear_From(year))
}

// Java method: CompilerOptions.setBrowserFeaturesetYear(com.google.javascript.jscomp.CompilerOptions.BrowserFeaturesetYear year)
func (this_ *CompilerOptions) SetBrowserFeaturesetYear_1_CompilerOptions_BrowserFeaturesetYear(year *CompilerOptions_BrowserFeaturesetYear) {
	_ = year
	this_.browserFeaturesetYear = year
	this_.browserFeaturesetYear.SetDependentValuesFromYear(this_)
}

// Java method: CompilerOptions.setInstrumentForCoverageOnly(boolean instrumentForCoverageOnly)
func (this_ *CompilerOptions) SetInstrumentForCoverageOnly(instrumentForCoverageOnly bool) {
	_ = instrumentForCoverageOnly
	this_.instrumentForCoverageOnly = instrumentForCoverageOnly
}

// Java method: CompilerOptions.getInstrumentForCoverageOnly()
func (this_ *CompilerOptions) GetInstrumentForCoverageOnly() bool {
	return this_.instrumentForCoverageOnly
}

// Java method: CompilerOptions.setTypedAstOutputFile(java.nio.file.@org.jspecify.annotations.Nullable Path file)
func (this_ *CompilerOptions) SetTypedAstOutputFile(file j2gort.Ext_java_nio_file_Path) {
	_ = file
	this_.typedAstOutputFile = file
}

// Java method: CompilerOptions.getTypedAstOutputFile()
func (this_ *CompilerOptions) GetTypedAstOutputFile() j2gort.Ext_java_nio_file_Path {
	return this_.typedAstOutputFile
}

// Java method: CompilerOptions.setMergedPrecompiledLibraries(boolean mergedPrecompiledLibraries)
func (this_ *CompilerOptions) SetMergedPrecompiledLibraries(mergedPrecompiledLibraries bool) {
	_ = mergedPrecompiledLibraries
	this_.mergedPrecompiledLibraries = mergedPrecompiledLibraries
}

// Java method: CompilerOptions.getMergedPrecompiledLibraries()
func (this_ *CompilerOptions) GetMergedPrecompiledLibraries() bool {
	return this_.mergedPrecompiledLibraries
}

// Java method: CompilerOptions.setSkipTranspilationAndCrash(boolean value)
func (this_ *CompilerOptions) SetSkipTranspilationAndCrash(value bool) {
	_ = value
}

// Java method: CompilerOptions.setInputSourceMaps(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.SourceMapInput> inputSourceMaps)
func (this_ *CompilerOptions) SetInputSourceMaps(inputSourceMaps j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = inputSourceMaps
	this_.inputSourceMaps = inputSourceMaps
}

// Java method: CompilerOptions.getInputSourceMaps()
func (this_ *CompilerOptions) GetInputSourceMaps() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.inputSourceMaps
}

// Java method: CompilerOptions.setInferConst(boolean value)
func (this_ *CompilerOptions) SetInferConst(value bool) {
	_ = value
	this_.inferConsts = value
}

// Java method: CompilerOptions.shouldInferConsts()
func (this_ *CompilerOptions) ShouldInferConsts() bool {
	return this_.inferConsts
}

// Java method: CompilerOptions.setIncrementalChecks(com.google.javascript.jscomp.CompilerOptions.IncrementalCheckMode value)
func (this_ *CompilerOptions) SetIncrementalChecks(value CompilerOptions_IncrementalCheckMode) {
	_ = value
	this_.incrementalCheckMode = value
	switch (value).Ordinal() {
	case int32(0) /* IncrementalCheckMode.OFF */ :
	case int32(1) /* IncrementalCheckMode.GENERATE_IJS */ :
		this_.SetPreserveTypeAnnotations(true)
		this_.SetOutputJs(StaticGet_CompilerOptions_OutputJs_NORMAL())
	}
}

// Java method: CompilerOptions.shouldGenerateTypedExterns()
func (this_ *CompilerOptions) ShouldGenerateTypedExterns() bool {
	return (this_.incrementalCheckMode == StaticGet_CompilerOptions_IncrementalCheckMode_GENERATE_IJS())
}

// Java method: CompilerOptions.setPrintExterns(boolean printExterns)
func (this_ *CompilerOptions) SetPrintExterns(printExterns bool) {
	_ = printExterns
	this_.printExterns = printExterns
}

// Java method: CompilerOptions.shouldPrintExterns()
func (this_ *CompilerOptions) ShouldPrintExterns() bool {
	return (this_.printExterns || (this_.incrementalCheckMode == StaticGet_CompilerOptions_IncrementalCheckMode_GENERATE_IJS()))
}

// Java method: CompilerOptions.setCheckGlobalThisLevel(com.google.javascript.jscomp.CheckLevel level)
func (this_ *CompilerOptions) SetCheckGlobalThisLevel(level *CheckLevel) {
	_ = level
}

// Java method: CompilerOptions.setNumParallelThreads(int parallelism)
func (this_ *CompilerOptions) SetNumParallelThreads(parallelism int32) {
	_ = parallelism
	this_.numParallelThreads = parallelism
}

// Java method: CompilerOptions.getNumParallelThreads()
func (this_ *CompilerOptions) GetNumParallelThreads() int32 {
	return this_.numParallelThreads
}

// Java method: CompilerOptions.setRenamePrefixNamespaceAssumeCrossChunkNames(boolean assume)
func (this_ *CompilerOptions) SetRenamePrefixNamespaceAssumeCrossChunkNames(assume bool) {
	_ = assume
	this_.renamePrefixNamespaceAssumeCrossChunkNames = assume
}

// Java method: CompilerOptions.assumeCrossChunkNamesForRenamePrefixNamespace()
func (this_ *CompilerOptions) AssumeCrossChunkNamesForRenamePrefixNamespace() bool {
	return this_.renamePrefixNamespaceAssumeCrossChunkNames
}

// Java method: CompilerOptions.setOptimizeLocalAccessForGlobalSymbolNamespace(com.google.javascript.jscomp.CompilerOptions.OptimizeLocalAccess optimize)
func (this_ *CompilerOptions) SetOptimizeLocalAccessForGlobalSymbolNamespace(optimize CompilerOptions_OptimizeLocalAccess) {
	_ = optimize
	this_.optimizeLocalAccessForGlobalSymbolNamespace = optimize
}

// Java method: CompilerOptions.getOptimizeLocalAccessForGlobalSymbolNamespace()
func (this_ *CompilerOptions) GetOptimizeLocalAccessForGlobalSymbolNamespace() CompilerOptions_OptimizeLocalAccess {
	return this_.optimizeLocalAccessForGlobalSymbolNamespace
}

// Java method: CompilerOptions.shouldCollapseProperties()
func (this_ *CompilerOptions) ShouldCollapseProperties() bool {
	return (this_.collapsePropertiesLevel != StaticGet_CompilerOptions_PropertyCollapseLevel_NONE())
}

// Java method: CompilerOptions.getPropertyCollapseLevel()
func (this_ *CompilerOptions) GetPropertyCollapseLevel() CompilerOptions_PropertyCollapseLevel {
	return this_.collapsePropertiesLevel
}

// Java method: CompilerOptions.setCollapseObjectLiterals(boolean enabled)
func (this_ *CompilerOptions) SetCollapseObjectLiterals(enabled bool) {
	_ = enabled
	this_.collapseObjectLiterals = enabled
}

// Java method: CompilerOptions.getCollapseObjectLiterals()
func (this_ *CompilerOptions) GetCollapseObjectLiterals() bool {
	return this_.collapseObjectLiterals
}

// Java method: CompilerOptions.setNameGenerator(com.google.javascript.jscomp.NameGenerator nameGenerator)
func (this_ *CompilerOptions) SetNameGenerator(nameGenerator NameGenerator) {
	_ = nameGenerator
	this_.nameGenerator = nameGenerator
}

// Java method: CompilerOptions.getNameGenerator()
func (this_ *CompilerOptions) GetNameGenerator() NameGenerator {
	return this_.nameGenerator
}

// Java method: CompilerOptions.setReplaceMessagesWithChromeI18n(boolean replaceMessagesWithChromeI18n, java.lang.String tcProjectId)
func (this_ *CompilerOptions) SetReplaceMessagesWithChromeI18n(replaceMessagesWithChromeI18n bool, tcProjectId string) {
	_ = replaceMessagesWithChromeI18n
	_ = tcProjectId
	if (replaceMessagesWithChromeI18n && (this_.messageBundle != nil && !this_.messageBundle.IsJavaNil_())) && !(j2gort.InstanceOf[*EmptyMessageBundle](this_.messageBundle)) {
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "When replacing messages with chrome.i18n.getMessage, a message bundle should not be specified.")))
	}
	this_.replaceMessagesWithChromeI18n = replaceMessagesWithChromeI18n
	this_.tcProjectId = tcProjectId
}

// Java method: CompilerOptions.getTcProjectId()
func (this_ *CompilerOptions) GetTcProjectId() string {
	return this_.tcProjectId
}

// Java method: CompilerOptions.shouldRunReplaceMessagesForChrome()
func (this_ *CompilerOptions) ShouldRunReplaceMessagesForChrome() bool {
	if this_.replaceMessagesWithChromeI18n {
		if !((this_.messageBundle == nil || this_.messageBundle.IsJavaNil_()) || j2gort.InstanceOf[*EmptyMessageBundle](this_.messageBundle)) {
			j2gort.CheckStateStr(false, "When replacing messages with chrome.i18n.getMessage, a message bundle should not be specified.")
		}
		if !(!(this_.doLateLocalization)) {
			j2gort.CheckStateStr(false, "Late localization is not supported for chrome.i18n.getMessage")
		}
		return true
	} else {
		return false
	}
}

// Java method: CompilerOptions.setAssumeForwardDeclaredForMissingTypes(boolean assumeForwardDeclaredForMissingTypes)
func (this_ *CompilerOptions) SetAssumeForwardDeclaredForMissingTypes(assumeForwardDeclaredForMissingTypes bool) {
	_ = assumeForwardDeclaredForMissingTypes
	this_.assumeForwardDeclaredForMissingTypes = assumeForwardDeclaredForMissingTypes
}

// Java method: CompilerOptions.assumeForwardDeclaredForMissingTypes()
func (this_ *CompilerOptions) AssumeForwardDeclaredForMissingTypes() bool {
	return this_.assumeForwardDeclaredForMissingTypes
}

// Java method: CompilerOptions.setPreferSingleQuotes(boolean enabled)
func (this_ *CompilerOptions) SetPreferSingleQuotes(enabled bool) {
	_ = enabled
	this_.preferSingleQuotes = enabled
}

// Java method: CompilerOptions.shouldPreferSingleQuotes()
func (this_ *CompilerOptions) ShouldPreferSingleQuotes() bool {
	return this_.preferSingleQuotes
}

// Java method: CompilerOptions.setTrustedStrings(boolean yes)
func (this_ *CompilerOptions) SetTrustedStrings(yes bool) {
	_ = yes
	this_.trustedStrings = yes
}

// Java method: CompilerOptions.assumeTrustedStrings()
func (this_ *CompilerOptions) AssumeTrustedStrings() bool {
	return this_.trustedStrings
}

// Java method: CompilerOptions.setPrintSourceAfterEachPass(boolean printSource)
func (this_ *CompilerOptions) SetPrintSourceAfterEachPass(printSource bool) {
	_ = printSource
	this_.printSourceAfterEachPass = printSource
}

// Java method: CompilerOptions.shouldPrintSourceAfterEachPass()
func (this_ *CompilerOptions) ShouldPrintSourceAfterEachPass() bool {
	return this_.printSourceAfterEachPass
}

// Java method: CompilerOptions.setFilesToPrintAfterEachPassRegexList(java.util.List<java.lang.String> filePathRegexList)
func (this_ *CompilerOptions) SetFilesToPrintAfterEachPassRegexList(filePathRegexList j2gort.Ext_java_util_List) {
	_ = filePathRegexList
	this_.filesToPrintAfterEachPassRegexList = filePathRegexList
}

// Java method: CompilerOptions.setChunksToPrintAfterEachPassRegexList(java.util.List<java.lang.String> chunkPathRegexList)
func (this_ *CompilerOptions) SetChunksToPrintAfterEachPassRegexList(chunkPathRegexList j2gort.Ext_java_util_List) {
	_ = chunkPathRegexList
	this_.chunksToPrintAfterEachPassRegexList = chunkPathRegexList
}

// Java method: CompilerOptions.setQnameUsesToPrintAfterEachPassList(java.util.List<java.lang.String> qnameRegexList)
func (this_ *CompilerOptions) SetQnameUsesToPrintAfterEachPassList(qnameRegexList j2gort.Ext_java_util_List) {
	_ = qnameRegexList
	this_.qnameUsesToPrintAfterEachPassList = qnameRegexList
}

// Java method: CompilerOptions.getFilesToPrintAfterEachPassRegexList()
func (this_ *CompilerOptions) GetFilesToPrintAfterEachPassRegexList() j2gort.Ext_java_util_List {
	return this_.filesToPrintAfterEachPassRegexList
}

// Java method: CompilerOptions.getChunksToPrintAfterEachPassRegexList()
func (this_ *CompilerOptions) GetChunksToPrintAfterEachPassRegexList() j2gort.Ext_java_util_List {
	return this_.chunksToPrintAfterEachPassRegexList
}

// Java method: CompilerOptions.getQnameUsesToPrintAfterEachPassList()
func (this_ *CompilerOptions) GetQnameUsesToPrintAfterEachPassList() j2gort.Ext_java_util_List {
	return this_.qnameUsesToPrintAfterEachPassList
}

// Java method: CompilerOptions.getTracerMode()
func (this_ *CompilerOptions) GetTracerMode() CompilerOptions_TracerMode {
	return this_.tracer
}

// Java method: CompilerOptions.setTracerMode(com.google.javascript.jscomp.CompilerOptions.TracerMode mode)
func (this_ *CompilerOptions) SetTracerMode(mode CompilerOptions_TracerMode) {
	_ = mode
	this_.tracer = mode
}

// Java method: CompilerOptions.getTracerOutput()
func (this_ *CompilerOptions) GetTracerOutput() j2gort.Ext_java_nio_file_Path {
	return this_.tracerOutput
}

// Java method: CompilerOptions.setTracerOutput(java.nio.file.Path out)
func (this_ *CompilerOptions) SetTracerOutput(out j2gort.Ext_java_nio_file_Path) {
	_ = out
	this_.tracerOutput = out
}

// Java method: CompilerOptions.getExtraReportGenerators()
func (this_ *CompilerOptions) GetExtraReportGenerators() j2gort.Ext_java_util_List {
	return this_.extraReportGenerators
}

// Java method: CompilerOptions.addReportGenerator(com.google.javascript.jscomp.SortingErrorManager.ErrorReportGenerator generator)
func (this_ *CompilerOptions) AddReportGenerator(generator SortingErrorManager_ErrorReportGenerator) {
	_ = generator
	this_.extraReportGenerators.Add_1(generator)
}

// Java method: CompilerOptions.setProtectHiddenSideEffects(boolean enable)
func (this_ *CompilerOptions) SetProtectHiddenSideEffects(enable bool) {
	_ = enable
	this_.protectHiddenSideEffects = enable
}

// Java method: CompilerOptions.shouldProtectHiddenSideEffects()
func (this_ *CompilerOptions) ShouldProtectHiddenSideEffects() bool {
	return (this_.protectHiddenSideEffects && (!(this_.checksOnly) || (this_.GetTypedAstOutputFile() != nil)))
}

// Java method: CompilerOptions.setAssumeGettersArePure(boolean x)
func (this_ *CompilerOptions) SetAssumeGettersArePure(x bool) {
	_ = x
	this_.assumeGettersArePure = x
}

// Java method: CompilerOptions.getAssumeGettersArePure()
func (this_ *CompilerOptions) GetAssumeGettersArePure() bool {
	return this_.assumeGettersArePure
}

// Java method: CompilerOptions.setAssumePropertiesAreStaticallyAnalyzable(boolean x)
func (this_ *CompilerOptions) SetAssumePropertiesAreStaticallyAnalyzable(x bool) {
	_ = x
	this_.assumePropertiesAreStaticallyAnalyzable = x
}

// Java method: CompilerOptions.getAssumePropertiesAreStaticallyAnalyzable()
func (this_ *CompilerOptions) GetAssumePropertiesAreStaticallyAnalyzable() bool {
	return this_.assumePropertiesAreStaticallyAnalyzable
}

// Java method: CompilerOptions.setAssumeStaticInheritanceIsNotUsed(boolean x)
func (this_ *CompilerOptions) SetAssumeStaticInheritanceIsNotUsed(x bool) {
	_ = x
	this_.assumeStaticInheritanceIsNotUsed = x
}

// Java method: CompilerOptions.getAssumeStaticInheritanceIsNotUsed()
func (this_ *CompilerOptions) GetAssumeStaticInheritanceIsNotUsed() bool {
	return this_.assumeStaticInheritanceIsNotUsed
}

// Java method: CompilerOptions.getConformanceReportingMode()
func (this_ *CompilerOptions) GetConformanceReportingMode() CompilerOptions_ConformanceReportingMode {
	return this_.conformanceReportingMode
}

// Java method: CompilerOptions.setConformanceReportingMode(com.google.javascript.jscomp.CompilerOptions.ConformanceReportingMode mode)
func (this_ *CompilerOptions) SetConformanceReportingMode(mode CompilerOptions_ConformanceReportingMode) {
	_ = mode
	this_.conformanceReportingMode = mode
}

// Java method: CompilerOptions.setConformanceRemoveRegexFromPath(com.google.common.base.Optional<java.util.regex.Pattern> pattern)
func (this_ *CompilerOptions) SetConformanceRemoveRegexFromPath(pattern j2gort.Ext_com_google_common_base_Optional) {
	_ = pattern
	this_.conformanceRemoveRegexFromPath = pattern
}

// Java method: CompilerOptions.getConformanceRemoveRegexFromPath()
func (this_ *CompilerOptions) GetConformanceRemoveRegexFromPath() j2gort.Ext_com_google_common_base_Optional {
	return this_.conformanceRemoveRegexFromPath
}

// Java method: CompilerOptions.setWrapGoogModulesForWhitespaceOnly(boolean enable)
func (this_ *CompilerOptions) SetWrapGoogModulesForWhitespaceOnly(enable bool) {
	_ = enable
	this_.wrapGoogModulesForWhitespaceOnly = enable
}

// Java method: CompilerOptions.shouldWrapGoogModulesForWhitespaceOnly()
func (this_ *CompilerOptions) ShouldWrapGoogModulesForWhitespaceOnly() bool {
	return this_.wrapGoogModulesForWhitespaceOnly
}

// Java method: CompilerOptions.setBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf(boolean b)
func (this_ *CompilerOptions) SetBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf(b bool) {
	_ = b
	this_.rewriteModulesBeforeTypechecking = b
}

// Java method: CompilerOptions.shouldRewriteModulesBeforeTypechecking()
func (this_ *CompilerOptions) ShouldRewriteModulesBeforeTypechecking() bool {
	return (this_.enableModuleRewriting && (this_.rewriteModulesBeforeTypechecking || this_.processCommonJSModules))
}

// Java method: CompilerOptions.setEnableModuleRewriting(boolean enable)
func (this_ *CompilerOptions) SetEnableModuleRewriting(enable bool) {
	_ = enable
	this_.enableModuleRewriting = enable
}

// Java method: CompilerOptions.shouldRewriteModulesAfterTypechecking()
func (this_ *CompilerOptions) ShouldRewriteModulesAfterTypechecking() bool {
	return (this_.enableModuleRewriting && !(this_.rewriteModulesBeforeTypechecking))
}

// Java method: CompilerOptions.shouldRewriteModules()
func (this_ *CompilerOptions) ShouldRewriteModules() bool {
	return this_.enableModuleRewriting
}

// Java method: CompilerOptions.setPrintConfig(boolean printConfig)
func (this_ *CompilerOptions) SetPrintConfig(printConfig bool) {
	_ = printConfig
	this_.printConfig = printConfig
}

// Java method: CompilerOptions.shouldPrintConfig()
func (this_ *CompilerOptions) ShouldPrintConfig() bool {
	return this_.printConfig
}

// Java method: CompilerOptions.setAllowDynamicImport(boolean value)
func (this_ *CompilerOptions) SetAllowDynamicImport(value bool) {
	_ = value
	this_.allowDynamicImport = value
}

// Java method: CompilerOptions.shouldAllowDynamicImport()
func (this_ *CompilerOptions) ShouldAllowDynamicImport() bool {
	return this_.allowDynamicImport
}

// Java method: CompilerOptions.getDynamicImportAlias()
func (this_ *CompilerOptions) GetDynamicImportAlias() string {
	return this_.dynamicImportAlias
}

// Java method: CompilerOptions.setDynamicImportAlias(java.lang.String value)
func (this_ *CompilerOptions) SetDynamicImportAlias(value string) {
	_ = value
	this_.dynamicImportAlias = value
}

// Java method: CompilerOptions.shouldAliasDynamicImport()
func (this_ *CompilerOptions) ShouldAliasDynamicImport() bool {
	return (this_.dynamicImportAlias != j2gort.NullString)
}

// Java method: CompilerOptions.isRemoveUnusedClassProperties()
func (this_ *CompilerOptions) IsRemoveUnusedClassProperties() bool {
	return this_.removeUnusedClassProperties
}

// Java method: CompilerOptions.setRemoveUnusedClassProperties(boolean removeUnusedClassProperties)
func (this_ *CompilerOptions) SetRemoveUnusedClassProperties(removeUnusedClassProperties bool) {
	_ = removeUnusedClassProperties
	this_.removeUnusedClassProperties = removeUnusedClassProperties
}

// Java method: CompilerOptions.getDefineReplacements()
func (this_ *CompilerOptions) GetDefineReplacements() j2gort.Ext_com_google_common_collect_ImmutableMap {
	var map_ j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = map_
	{
		cur_1 := j2gort.NewMapEntryCursor(this_.defineReplacements)
		for cur_1.NextEntry() {
			var name string = cur_1.KeyStr()
			_ = name
			var value any = cur_1.GetValue()
			_ = value
			var b *j2gort.BoxBoolean
			_ = b
			var number j2gort.Number
			_ = number
			var string_ string
			_ = string_
			var ok_2 bool
			b, ok_2 = j2gort.TryCast[*j2gort.BoxBoolean](value)
			if ok_2 {
				map_.Put_2(j2gort.BoxString(name), j2gort.BoxNode(NodeUtil_BooleanNode(b.BooleanValue())))
			} else {
				var ok_3 bool
				number, ok_3 = j2gort.TryCast[j2gort.Number](value)
				if ok_3 {
					map_.Put_2(j2gort.BoxString(name), j2gort.BoxNode(NodeUtil_NumberNode(number.DoubleValue(), pkg_rhino.Node(0))))
				} else {
					var ok_4 bool
					string_, ok_4 = j2gort.TryCastString(value)
					if ok_4 {
						map_.Put_2(j2gort.BoxString(name), j2gort.BoxNode(pkg_rhino.IR_String(string_)))
					} else {
						panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.StrVal(value))))
					}
				}
			}
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](map_.BuildOrThrow())
}

// Java method: CompilerOptions.setDefineToBooleanLiteral(java.lang.String defineName, boolean value)
func (this_ *CompilerOptions) SetDefineToBooleanLiteral(defineName string, value bool) {
	_ = defineName
	_ = value
	j2gort.MapPutStr(this_.defineReplacements, defineName, j2gort.ValueOfBoolean(value))
}

// Java method: CompilerOptions.setDefineToStringLiteral(java.lang.String defineName, java.lang.String value)
func (this_ *CompilerOptions) SetDefineToStringLiteral(defineName string, value string) {
	_ = defineName
	_ = value
	j2gort.MapPutStr(this_.defineReplacements, defineName, j2gort.BoxString(value))
}

// Java method: CompilerOptions.setDefineToNumberLiteral(java.lang.String defineName, int value)
func (this_ *CompilerOptions) SetDefineToNumberLiteral(defineName string, value int32) {
	_ = defineName
	_ = value
	j2gort.MapPutStr(this_.defineReplacements, defineName, j2gort.ValueOfInt(value))
}

// Java method: CompilerOptions.setDefineToDoubleLiteral(java.lang.String defineName, double value)
func (this_ *CompilerOptions) SetDefineToDoubleLiteral(defineName string, value float64) {
	_ = defineName
	_ = value
	j2gort.MapPutStr(this_.defineReplacements, defineName, j2gort.ValueOfDouble(value))
}

// Java method: CompilerOptions.skipAllCompilerPasses()
func (this_ *CompilerOptions) SkipAllCompilerPasses() {
	this_.SetSkipNonTranspilationPasses(true)
}

// Java method: CompilerOptions.enables(com.google.javascript.jscomp.DiagnosticGroup group)
func (this_ *CompilerOptions) Enables(group *DiagnosticGroup) bool {
	_ = group
	return (this_.warningsGuard.MustRunChecks(group) == pkg_base.StaticGet_Tri_TRUE())
}

// Java method: CompilerOptions.disables(com.google.javascript.jscomp.DiagnosticGroup group)
func (this_ *CompilerOptions) Disables(group *DiagnosticGroup) bool {
	_ = group
	return (this_.warningsGuard.MustRunChecks(group) == pkg_base.StaticGet_Tri_FALSE())
}

// Java method: CompilerOptions.setUnknownDefinesToIgnore(com.google.common.collect.ImmutableList<java.lang.String> unknownDefinesToIgnore)
func (this_ *CompilerOptions) SetUnknownDefinesToIgnore(unknownDefinesToIgnore j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = unknownDefinesToIgnore
	this_.unknownDefinesToIgnore = unknownDefinesToIgnore
}

// Java method: CompilerOptions.addUnknownDefinesToIgnore(com.google.common.collect.ImmutableList<java.lang.String> unknownDefinesToIgnore)
func (this_ *CompilerOptions) AddUnknownDefinesToIgnore(unknownDefinesToIgnore j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = unknownDefinesToIgnore
	this_.unknownDefinesToIgnore = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}).AddAllSelf_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.unknownDefinesToIgnore)).AddAllSelf_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](unknownDefinesToIgnore)).BuildImmList_()
}

// Java method: CompilerOptions.getUnknownDefinesToIgnore()
func (this_ *CompilerOptions) GetUnknownDefinesToIgnore() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.unknownDefinesToIgnore
}

// Java method: CompilerOptions.setWarningLevel(com.google.javascript.jscomp.DiagnosticGroup type, com.google.javascript.jscomp.CheckLevel level)
func (this_ *CompilerOptions) SetWarningLevel(type_ *DiagnosticGroup, level *CheckLevel) {
	_ = type_
	_ = level
	this_.AddWarningsGuard(NewDiagnosticGroupWarningsGuard(type_, level))
}

// Java method: CompilerOptions.getWarningsGuard()
func (this_ *CompilerOptions) GetWarningsGuard() WarningsGuard {
	return ToIface_WarningsGuard(this_.warningsGuard)
}

// Java method: CompilerOptions.resetWarningsGuard()
func (this_ *CompilerOptions) ResetWarningsGuard() {
	this_.warningsGuard = NewComposeWarningsGuard_1_WarningsGuard([]WarningsGuard{})
}

// Java method: CompilerOptions.addWarningsGuard(com.google.javascript.jscomp.WarningsGuard guard)
func (this_ *CompilerOptions) AddWarningsGuard(guard WarningsGuard) {
	_ = guard
	this_.warningsGuard.AddGuard(guard)
}

// Java method: CompilerOptions.setRenamingPolicy(com.google.javascript.jscomp.VariableRenamingPolicy newVariablePolicy, com.google.javascript.jscomp.PropertyRenamingPolicy newPropertyPolicy)
func (this_ *CompilerOptions) SetRenamingPolicy(newVariablePolicy VariableRenamingPolicy, newPropertyPolicy PropertyRenamingPolicy) {
	_ = newVariablePolicy
	_ = newPropertyPolicy
	this_.variableRenaming = newVariablePolicy
	this_.propertyRenaming = newPropertyPolicy
}

// Java method: CompilerOptions.setReplaceIdGenerators(boolean replaceIdGenerators)
func (this_ *CompilerOptions) SetReplaceIdGenerators(replaceIdGenerators bool) {
	_ = replaceIdGenerators
	this_.replaceIdGenerators = replaceIdGenerators
}

// Java method: CompilerOptions.shouldReplaceIdGenerators()
func (this_ *CompilerOptions) ShouldReplaceIdGenerators() bool {
	return this_.replaceIdGenerators
}

// Java method: CompilerOptions.setIdGenerators(java.util.Set<java.lang.String> idGenerators)
func (this_ *CompilerOptions) SetIdGenerators_1_Ext_java_util_Set(idGenerators j2gort.Ext_java_util_Set) {
	_ = idGenerators
	var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = builder
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(idGenerators))
		for cur_5.HasNext() {
			name := cur_5.NextStr()
			_ = name
			builder.Put_2(j2gort.BoxString(name), j2gort.PtrToAny(StaticGet_RenamingToken_INCONSISTENT()))
		}
	}
	this_.idGenerators = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
}

// Java method: CompilerOptions.setIdGenerators(java.util.Map<java.lang.String,com.google.javascript.jscomp.RenamingMap> idGenerators)
func (this_ *CompilerOptions) SetIdGenerators_1_Ext_java_util_Map(idGenerators j2gort.Ext_java_util_Map) {
	_ = idGenerators
	this_.idGenerators = j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(idGenerators)
}

// Java method: CompilerOptions.getIdGenerators()
func (this_ *CompilerOptions) GetIdGenerators() j2gort.Ext_java_util_Map {
	return j2gort.AsIface[j2gort.Ext_java_util_Map](this_.idGenerators)
}

// Java method: CompilerOptions.setIdGeneratorsMap(java.lang.String previousMappings)
func (this_ *CompilerOptions) SetIdGeneratorsMap(previousMappings string) {
	_ = previousMappings
	this_.idGeneratorsMapSerialized = previousMappings
}

// Java method: CompilerOptions.getIdGeneratorsMapSerialized()
func (this_ *CompilerOptions) GetIdGeneratorsMapSerialized() string {
	return this_.idGeneratorsMapSerialized
}

// Java method: CompilerOptions.setXidHashFunction(com.google.javascript.jscomp.Xid.HashFunction xidHashFunction)
func (this_ *CompilerOptions) SetXidHashFunction(xidHashFunction Xid_HashFunction) {
	_ = xidHashFunction
	this_.xidHashFunction = xidHashFunction
}

// Java method: CompilerOptions.getXidHashFunction()
func (this_ *CompilerOptions) GetXidHashFunction() Xid_HashFunction {
	return this_.xidHashFunction
}

// Java method: CompilerOptions.setChunkIdHashFunction(com.google.javascript.jscomp.Xid.HashFunction chunkIdHashFunction)
func (this_ *CompilerOptions) SetChunkIdHashFunction(chunkIdHashFunction Xid_HashFunction) {
	_ = chunkIdHashFunction
	this_.chunkIdHashFunction = chunkIdHashFunction
}

// Java method: CompilerOptions.getChunkIdHashFunction()
func (this_ *CompilerOptions) GetChunkIdHashFunction() Xid_HashFunction {
	return this_.chunkIdHashFunction
}

// Java method: CompilerOptions.setInlineFunctions(com.google.javascript.jscomp.CompilerOptions.Reach reach)
func (this_ *CompilerOptions) SetInlineFunctions(reach CompilerOptions_Reach) {
	_ = reach
	this_.inlineFunctionsLevel = reach
}

// Java method: CompilerOptions.getInlineFunctionsLevel()
func (this_ *CompilerOptions) GetInlineFunctionsLevel() CompilerOptions_Reach {
	return this_.inlineFunctionsLevel
}

// Java method: CompilerOptions.setMaxFunctionSizeAfterInlining(int funAstSize)
func (this_ *CompilerOptions) SetMaxFunctionSizeAfterInlining(funAstSize int32) {
	_ = funAstSize
	if !(funAstSize > int32(0)) {
		j2gort.CheckArgument(false)
	}
	this_.maxFunctionSizeAfterInlining = funAstSize
}

// Java method: CompilerOptions.getMaxFunctionSizeAfterInlining()
func (this_ *CompilerOptions) GetMaxFunctionSizeAfterInlining() int32 {
	return this_.maxFunctionSizeAfterInlining
}

// Java method: CompilerOptions.setInlineVariables(boolean inlineVariables)
func (this_ *CompilerOptions) SetInlineVariables_1_bool(inlineVariables bool) {
	_ = inlineVariables
	this_.inlineVariables = inlineVariables
}

// Java method: CompilerOptions.shouldInlineVariables()
func (this_ *CompilerOptions) ShouldInlineVariables() bool {
	return this_.inlineVariables
}

// Java method: CompilerOptions.setInlineVariables(com.google.javascript.jscomp.CompilerOptions.Reach reach)
func (this_ *CompilerOptions) SetInlineVariables_1_CompilerOptions_Reach(reach CompilerOptions_Reach) {
	_ = reach
	switch (reach).Ordinal() {
	case int32(0) /* Reach.ALL */ :
		this_.inlineVariables = true
		this_.inlineLocalVariables = true
	case int32(1) /* Reach.LOCAL_ONLY */ :
		this_.inlineVariables = false
		this_.inlineLocalVariables = true
	case int32(2) /* Reach.NONE */ :
		this_.inlineVariables = false
		this_.inlineLocalVariables = false
	}
}

// Java method: CompilerOptions.setInlineProperties(boolean enable)
func (this_ *CompilerOptions) SetInlineProperties(enable bool) {
	_ = enable
	this_.inlineProperties = enable
}

// Java method: CompilerOptions.shouldInlineProperties()
func (this_ *CompilerOptions) ShouldInlineProperties() bool {
	return this_.inlineProperties
}

// Java method: CompilerOptions.setRemoveUnusedVariables(com.google.javascript.jscomp.CompilerOptions.Reach reach)
func (this_ *CompilerOptions) SetRemoveUnusedVariables(reach CompilerOptions_Reach) {
	_ = reach
	switch (reach).Ordinal() {
	case int32(0) /* Reach.ALL */ :
		this_.removeUnusedVars = true
		this_.removeUnusedLocalVars = true
	case int32(1) /* Reach.LOCAL_ONLY */ :
		this_.removeUnusedVars = false
		this_.removeUnusedLocalVars = true
	case int32(2) /* Reach.NONE */ :
		this_.removeUnusedVars = false
		this_.removeUnusedLocalVars = false
	}
}

// Java method: CompilerOptions.shouldRemoveUnusedVariables()
func (this_ *CompilerOptions) ShouldRemoveUnusedVariables() bool {
	return this_.removeUnusedVars
}

// Java method: CompilerOptions.shouldRemoveUnusedLocalVariables()
func (this_ *CompilerOptions) ShouldRemoveUnusedLocalVariables() bool {
	return this_.removeUnusedLocalVars
}

// Java method: CompilerOptions.setReplaceStringsConfiguration(java.lang.String placeholderToken, java.util.List<java.lang.String> functionDescriptors)
func (this_ *CompilerOptions) SetReplaceStringsConfiguration(placeholderToken string, functionDescriptors j2gort.Ext_java_util_List) {
	_ = placeholderToken
	_ = functionDescriptors
	this_.replaceStringsPlaceholderToken = placeholderToken
	this_.replaceStringsFunctionDescriptions = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](functionDescriptors))
}

// Java method: CompilerOptions.setRemoveAbstractMethods(boolean remove)
func (this_ *CompilerOptions) SetRemoveAbstractMethods(remove bool) {
	_ = remove
	this_.removeAbstractMethods = remove
}

// Java method: CompilerOptions.shouldRemoveAbstractMethods()
func (this_ *CompilerOptions) ShouldRemoveAbstractMethods() bool {
	return this_.removeAbstractMethods
}

// Java method: CompilerOptions.setRemoveClosureAsserts(boolean remove)
func (this_ *CompilerOptions) SetRemoveClosureAsserts(remove bool) {
	_ = remove
	this_.removeClosureAsserts = remove
}

// Java method: CompilerOptions.shouldRemoveClosureAsserts()
func (this_ *CompilerOptions) ShouldRemoveClosureAsserts() bool {
	return this_.removeClosureAsserts
}

// Java method: CompilerOptions.setRemoveJ2clAsserts(boolean remove)
func (this_ *CompilerOptions) SetRemoveJ2clAsserts(remove bool) {
	_ = remove
	this_.removeJ2clAsserts = remove
}

// Java method: CompilerOptions.shouldRemoveJ2clAsserts()
func (this_ *CompilerOptions) ShouldRemoveJ2clAsserts() bool {
	return this_.removeJ2clAsserts
}

// Java method: CompilerOptions.setColorizeErrorOutput(boolean colorizeErrorOutput)
func (this_ *CompilerOptions) SetColorizeErrorOutput(colorizeErrorOutput bool) {
	_ = colorizeErrorOutput
	this_.colorizeErrorOutput = colorizeErrorOutput
}

// Java method: CompilerOptions.shouldColorizeErrorOutput()
func (this_ *CompilerOptions) ShouldColorizeErrorOutput() bool {
	return this_.colorizeErrorOutput
}

// Java method: CompilerOptions.setChecksOnly(boolean checksOnly)
func (this_ *CompilerOptions) SetChecksOnly(checksOnly bool) {
	_ = checksOnly
	this_.checksOnly = checksOnly
}

// Java method: CompilerOptions.isChecksOnly()
func (this_ *CompilerOptions) IsChecksOnly() bool {
	return this_.checksOnly
}

// Java method: CompilerOptions.setOutputJs(com.google.javascript.jscomp.CompilerOptions.OutputJs outputJs)
func (this_ *CompilerOptions) SetOutputJs(outputJs CompilerOptions_OutputJs) {
	_ = outputJs
	this_.outputJs = outputJs
}

// Java method: CompilerOptions.getOutputJs()
func (this_ *CompilerOptions) GetOutputJs() CompilerOptions_OutputJs {
	return this_.outputJs
}

// Java method: CompilerOptions.setGenerateExports(boolean generateExports)
func (this_ *CompilerOptions) SetGenerateExports(generateExports bool) {
	_ = generateExports
	this_.generateExports = generateExports
}

// Java method: CompilerOptions.shouldGenerateExports()
func (this_ *CompilerOptions) ShouldGenerateExports() bool {
	return this_.generateExports
}

// Java method: CompilerOptions.setExportLocalPropertyDefinitions(boolean export)
func (this_ *CompilerOptions) SetExportLocalPropertyDefinitions(export bool) {
	_ = export
	this_.exportLocalPropertyDefinitions = export
}

// Java method: CompilerOptions.shouldExportLocalPropertyDefinitions()
func (this_ *CompilerOptions) ShouldExportLocalPropertyDefinitions() bool {
	return this_.exportLocalPropertyDefinitions
}

// Java method: CompilerOptions.setAngularPass(boolean angularPass)
func (this_ *CompilerOptions) SetAngularPass(angularPass bool) {
	_ = angularPass
	this_.angularPass = angularPass
}

// Java method: CompilerOptions.getAngularPass()
func (this_ *CompilerOptions) GetAngularPass() bool {
	return this_.angularPass
}

// Java method: CompilerOptions.setPolymerVersion(java.lang.Integer polymerVersion)
func (this_ *CompilerOptions) SetPolymerVersion(polymerVersion *j2gort.BoxInt) {
	_ = polymerVersion
	if !(((polymerVersion == nil) || (j2gort.UnboxInt(polymerVersion) == int32(1))) || (j2gort.UnboxInt(polymerVersion) == int32(2))) {
		j2gort.CheckArgumentFmt(false, "Invalid Polymer version: (%s)", j2gort.PtrToAny(polymerVersion))
	}
	this_.polymerPass = (polymerVersion != nil)
}

// Java method: CompilerOptions.isPolymerPassEnabled()
func (this_ *CompilerOptions) IsPolymerPassEnabled() bool {
	return this_.polymerPass
}

// Java method: CompilerOptions.setChromePass(boolean chromePass)
func (this_ *CompilerOptions) SetChromePass(chromePass bool) {
	_ = chromePass
	this_.chromePass = chromePass
}

// Java method: CompilerOptions.isChromePassEnabled()
func (this_ *CompilerOptions) IsChromePassEnabled() bool {
	return this_.chromePass
}

// Java method: CompilerOptions.setJ2clPass(com.google.javascript.jscomp.CompilerOptions.J2clPassMode j2clPassMode)
func (this_ *CompilerOptions) SetJ2clPass(j2clPassMode CompilerOptions_J2clPassMode) {
	_ = j2clPassMode
	this_.j2clPassMode = j2clPassMode
}

// Java method: CompilerOptions.getJ2clPass()
func (this_ *CompilerOptions) GetJ2clPass() CompilerOptions_J2clPassMode {
	return this_.j2clPassMode
}

// Java method: CompilerOptions.setJ2clMinifierEnabled(boolean enabled)
func (this_ *CompilerOptions) SetJ2clMinifierEnabled(enabled bool) {
	_ = enabled
	this_.j2clMinifierEnabled = enabled
}

// Java method: CompilerOptions.isJ2clMinifierEnabled()
func (this_ *CompilerOptions) IsJ2clMinifierEnabled() bool {
	return this_.j2clMinifierEnabled
}

// Java method: CompilerOptions.setJ2clMinifierPruningManifest(java.lang.String j2clMinifierPruningManifest)
func (this_ *CompilerOptions) SetJ2clMinifierPruningManifest(j2clMinifierPruningManifest string) {
	_ = j2clMinifierPruningManifest
	this_.j2clMinifierPruningManifest = j2clMinifierPruningManifest
}

// Java method: CompilerOptions.getJ2clMinifierPruningManifest()
func (this_ *CompilerOptions) GetJ2clMinifierPruningManifest() string {
	return this_.j2clMinifierPruningManifest
}

// Java method: CompilerOptions.setCodingConvention(com.google.javascript.jscomp.CodingConvention codingConvention)
func (this_ *CompilerOptions) SetCodingConvention(codingConvention CodingConvention) {
	_ = codingConvention
	this_.codingConvention = codingConvention
}

// Java method: CompilerOptions.getCodingConvention()
func (this_ *CompilerOptions) GetCodingConvention() CodingConvention {
	return this_.codingConvention
}

// Java method: CompilerOptions.setDependencyOptions(com.google.javascript.jscomp.DependencyOptions dependencyOptions)
func (this_ *CompilerOptions) SetDependencyOptions(dependencyOptions *DependencyOptions) {
	_ = dependencyOptions
	this_.dependencyOptions = dependencyOptions
}

// Java method: CompilerOptions.getDependencyOptions()
func (this_ *CompilerOptions) GetDependencyOptions() *DependencyOptions {
	return this_.dependencyOptions
}

// Java method: CompilerOptions.setSummaryDetailLevel(int summaryDetailLevel)
func (this_ *CompilerOptions) SetSummaryDetailLevel(summaryDetailLevel int32) {
	_ = summaryDetailLevel
	this_.summaryDetailLevel = summaryDetailLevel
}

// Java method: CompilerOptions.getSummaryDetailLevel()
func (this_ *CompilerOptions) GetSummaryDetailLevel() int32 {
	return this_.summaryDetailLevel
}

// Java method: CompilerOptions.setExtraAnnotationNames(java.lang.Iterable<java.lang.String> extraAnnotationNames)
func (this_ *CompilerOptions) SetExtraAnnotationNames(extraAnnotationNames j2gort.Ext_java_lang_Iterable) {
	_ = extraAnnotationNames
	this_.extraAnnotationNames = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_lang_Iterable(extraAnnotationNames)
}

// Java method: CompilerOptions.getExtraAnnotationNames()
func (this_ *CompilerOptions) GetExtraAnnotationNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.extraAnnotationNames
}

// Java method: CompilerOptions.setOutputCharset(java.nio.charset.Charset charset)
func (this_ *CompilerOptions) SetOutputCharset(charset j2gort.Ext_java_nio_charset_Charset) {
	_ = charset
	this_.outputCharset = charset
}

// Java method: CompilerOptions.getOutputCharset()
func (this_ *CompilerOptions) GetOutputCharset() j2gort.Ext_java_nio_charset_Charset {
	return this_.outputCharset
}

// Java method: CompilerOptions.setTweakProcessing(com.google.javascript.jscomp.CompilerOptions.TweakProcessing tweakProcessing)
func (this_ *CompilerOptions) SetTweakProcessing(tweakProcessing CompilerOptions_TweakProcessing) {
	_ = tweakProcessing
	this_.tweakProcessing = tweakProcessing
}

// Java method: CompilerOptions.getTweakProcessing()
func (this_ *CompilerOptions) GetTweakProcessing() CompilerOptions_TweakProcessing {
	return this_.tweakProcessing
}

// Java method: CompilerOptions.setLanguage(com.google.javascript.jscomp.CompilerOptions.LanguageMode language)
func (this_ *CompilerOptions) SetLanguage(language CompilerOptions_LanguageMode) {
	_ = language
	if !(language != StaticGet_CompilerOptions_LanguageMode_NO_TRANSPILE()) {
		j2gort.CheckState(false)
	}
	this_.SetLanguageIn(language)
	this_.SetLanguageOut(language)
}

// Java method: CompilerOptions.setLanguageIn(com.google.javascript.jscomp.CompilerOptions.LanguageMode languageIn)
func (this_ *CompilerOptions) SetLanguageIn(languageIn CompilerOptions_LanguageMode) {
	_ = languageIn
	if !(languageIn != StaticGet_CompilerOptions_LanguageMode_NO_TRANSPILE()) {
		j2gort.CheckState(false)
	}
	if languageIn == StaticGet_CompilerOptions_LanguageMode_STABLE() {
		this_.languageIn = StaticGet_CompilerOptions_LanguageMode_STABLE_IN()
	} else {
		this_.languageIn = languageIn
	}
}

// Java method: CompilerOptions.getLanguageIn()
func (this_ *CompilerOptions) GetLanguageIn() CompilerOptions_LanguageMode {
	return this_.languageIn
}

// Java method: CompilerOptions.setLanguageOut(com.google.javascript.jscomp.CompilerOptions.LanguageMode languageOut)
func (this_ *CompilerOptions) SetLanguageOut(languageOut CompilerOptions_LanguageMode) {
	_ = languageOut
	if languageOut == StaticGet_CompilerOptions_LanguageMode_NO_TRANSPILE() {
		this_.languageOutIsDefaultStrict = j2gort.Ext_com_google_common_base_Optional_Absent()
		this_.outputFeatureSet = j2gort.Ext_com_google_common_base_Optional_Absent()
	} else {
		if languageOut == StaticGet_CompilerOptions_LanguageMode_STABLE() {
			languageOut = StaticGet_CompilerOptions_LanguageMode_STABLE_OUT()
		} else {
			languageOut = languageOut
		}
		this_.languageOutIsDefaultStrict = j2gort.Ext_com_google_common_base_Optional_Of(j2gort.ValueOfBoolean(languageOut.IsDefaultStrict()))
		this_.SetOutputFeatureSet(languageOut.ToFeatureSet())
	}
}

// Java method: CompilerOptions.setOutputFeatureSet(com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet)
func (this_ *CompilerOptions) SetOutputFeatureSet(featureSet *pkg_parser.FeatureSet) {
	_ = featureSet
	this_.outputFeatureSet = j2gort.Ext_com_google_common_base_Optional_Of(j2gort.PtrToAny(featureSet))
}

// Java method: CompilerOptions.legacySetOutputFeatureSet(com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet)
func (this_ *CompilerOptions) LegacySetOutputFeatureSet(featureSet *pkg_parser.FeatureSet) {
	_ = featureSet
	this_.SetOutputFeatureSet(featureSet)
}

// Java method: CompilerOptions.getOutputFeatureSet()
func (this_ *CompilerOptions) GetOutputFeatureSet() *pkg_parser.FeatureSet {
	if this_.outputFeatureSet.IsPresent() {
		return j2gort.CastPtr[pkg_parser.FeatureSet](this_.outputFeatureSet.Get())
	}
	return this_.languageIn.ToFeatureSet()
}

// Java method: CompilerOptions.setTranspilePublicClassFields(boolean transpile)
func (this_ *CompilerOptions) SetTranspilePublicClassFields(transpile bool) {
	_ = transpile
	if transpile {
		this_.SetOutputFeatureSet(this_.GetOutputFeatureSet().Without_2(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS(), []*pkg_parser.FeatureSet_Feature{}))
	} else {
		this_.SetOutputFeatureSet(this_.GetOutputFeatureSet().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS()))
	}
}

// Java method: CompilerOptions.getTranspilePublicClassFields()
func (this_ *CompilerOptions) GetTranspilePublicClassFields() bool {
	return !(this_.GetOutputFeatureSet().Has(pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS()))
}

// Java method: CompilerOptions.setExperimentalOutputFeatureSet(com.google.javascript.jscomp.CompilerOptions.ExperimentalOutputFeatureSet experimentalOutputFeatureSet)
func (this_ *CompilerOptions) SetExperimentalOutputFeatureSet(experimentalOutputFeatureSet CompilerOptions_ExperimentalOutputFeatureSet) {
	_ = experimentalOutputFeatureSet
	this_.experimentalOutputFeatureSet = j2gort.Ext_com_google_common_base_Optional_FromNullable(j2gort.BoxNullableVal(experimentalOutputFeatureSet))
}

// Java method: CompilerOptions.getExperimentalOutputFeatureSet()
func (this_ *CompilerOptions) GetExperimentalOutputFeatureSet() j2gort.Ext_com_google_common_base_Optional {
	return this_.experimentalOutputFeatureSet
}

// Java method: CompilerOptions.needsTranspilationFrom(com.google.javascript.jscomp.parsing.parser.FeatureSet languageLevel)
func (this_ *CompilerOptions) NeedsTranspilationFrom(languageLevel *pkg_parser.FeatureSet) bool {
	_ = languageLevel
	return (this_.GetLanguageIn().ToFeatureSet().Contains_1_FeatureSet(languageLevel) && !(this_.GetOutputFeatureSet().Contains_1_FeatureSet(languageLevel)))
}

// Java method: CompilerOptions.needsTranspilationOf(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *CompilerOptions) NeedsTranspilationOf(feature *pkg_parser.FeatureSet_Feature) bool {
	_ = feature
	return (this_.GetLanguageIn().ToFeatureSet().Has(feature) && !(this_.GetOutputFeatureSet().Has(feature)))
}

// Java method: CompilerOptions.setEnvironment(com.google.javascript.jscomp.CompilerOptions.Environment environment)
func (this_ *CompilerOptions) SetEnvironment(environment CompilerOptions_Environment) {
	_ = environment
	this_.environment = environment
}

// Java method: CompilerOptions.getEnvironment()
func (this_ *CompilerOptions) GetEnvironment() CompilerOptions_Environment {
	return this_.environment
}

// Java method: CompilerOptions.setErrorHandler(com.google.javascript.jscomp.ErrorHandler handler)
func (this_ *CompilerOptions) SetErrorHandler(handler ErrorHandler) {
	_ = handler
	this_.errorHandler = handler
}

// Java method: CompilerOptions.getErrorHandler()
func (this_ *CompilerOptions) GetErrorHandler() ErrorHandler {
	return this_.errorHandler
}

// Java method: CompilerOptions.setInferTypes(boolean enable)
func (this_ *CompilerOptions) SetInferTypes(enable bool) {
	_ = enable
	this_.inferTypes = enable
}

// Java method: CompilerOptions.getInferTypes()
func (this_ *CompilerOptions) GetInferTypes() bool {
	return this_.inferTypes
}

// Java method: CompilerOptions.setNewTypeInference(boolean enable)
func (this_ *CompilerOptions) SetNewTypeInference(enable bool) {
	_ = enable
}

// Java method: CompilerOptions.isTypecheckingEnabled()
func (this_ *CompilerOptions) IsTypecheckingEnabled() bool {
	return this_.checkTypes
}

// Java method: CompilerOptions.assumeStrictThis()
func (this_ *CompilerOptions) AssumeStrictThis() bool {
	return this_.assumeStrictThis
}

// Java method: CompilerOptions.setAssumeStrictThis(boolean enable)
func (this_ *CompilerOptions) SetAssumeStrictThis(enable bool) {
	_ = enable
	this_.assumeStrictThis = enable
}

// Java method: CompilerOptions.setPropertiesThatMustDisambiguate(java.util.Set<java.lang.String> names)
func (this_ *CompilerOptions) SetPropertiesThatMustDisambiguate(names j2gort.Ext_java_util_Set) {
	_ = names
	this_.propertiesThatMustDisambiguate = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](names))
}

// Java method: CompilerOptions.getPropertiesThatMustDisambiguate()
func (this_ *CompilerOptions) GetPropertiesThatMustDisambiguate() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.propertiesThatMustDisambiguate
}

// Java method: CompilerOptions.setValidateRequiredInlinings(boolean validateRequireInliningAnnotation)
func (this_ *CompilerOptions) SetValidateRequiredInlinings(validateRequireInliningAnnotation bool) {
	_ = validateRequireInliningAnnotation
	this_.validateRequireInliningAnnotation = pkg_base.Tri_ForBoolean(validateRequireInliningAnnotation)
}

// Java method: CompilerOptions.getShouldValidateRequiredInlinings()
func (this_ *CompilerOptions) GetShouldValidateRequiredInlinings() pkg_base.Tri {
	return this_.validateRequireInliningAnnotation
}

// Java method: CompilerOptions.setPreserveDetailedSourceInfo(boolean preserveDetailedSourceInfo)
func (this_ *CompilerOptions) SetPreserveDetailedSourceInfo(preserveDetailedSourceInfo bool) {
	_ = preserveDetailedSourceInfo
	this_.preserveDetailedSourceInfo = preserveDetailedSourceInfo
}

// Java method: CompilerOptions.preservesDetailedSourceInfo()
func (this_ *CompilerOptions) PreservesDetailedSourceInfo() bool {
	return this_.preserveDetailedSourceInfo
}

// Java method: CompilerOptions.setPreserveNonJSDocComments(boolean preserveNonJSDocComments)
func (this_ *CompilerOptions) SetPreserveNonJSDocComments(preserveNonJSDocComments bool) {
	_ = preserveNonJSDocComments
	this_.preserveNonJSDocComments = preserveNonJSDocComments
}

// Java method: CompilerOptions.getPreserveNonJSDocComments()
func (this_ *CompilerOptions) GetPreserveNonJSDocComments() bool {
	return this_.preserveNonJSDocComments
}

// Java method: CompilerOptions.setContinueAfterErrors(boolean continueAfterErrors)
func (this_ *CompilerOptions) SetContinueAfterErrors(continueAfterErrors bool) {
	_ = continueAfterErrors
	this_.continueAfterErrors = continueAfterErrors
}

// Java method: CompilerOptions.canContinueAfterErrors()
func (this_ *CompilerOptions) CanContinueAfterErrors() bool {
	return this_.continueAfterErrors
}

// Java method: CompilerOptions.setParseJsDocDocumentation(com.google.javascript.jscomp.parsing.Config.JsDocParsing parseJsDocDocumentation)
func (this_ *CompilerOptions) SetParseJsDocDocumentation(parseJsDocDocumentation pkg_parsing.Config_JsDocParsing) {
	_ = parseJsDocDocumentation
	this_.parseJsDocDocumentation = parseJsDocDocumentation
}

// Java method: CompilerOptions.isParseJsDocDocumentation()
func (this_ *CompilerOptions) IsParseJsDocDocumentation() pkg_parsing.Config_JsDocParsing {
	return this_.parseJsDocDocumentation
}

// Java method: CompilerOptions.setSkipNonTranspilationPasses(boolean skipNonTranspilationPasses)
func (this_ *CompilerOptions) SetSkipNonTranspilationPasses(skipNonTranspilationPasses bool) {
	_ = skipNonTranspilationPasses
	this_.skipNonTranspilationPasses = skipNonTranspilationPasses
}

// Java method: CompilerOptions.getSkipNonTranspilationPasses()
func (this_ *CompilerOptions) GetSkipNonTranspilationPasses() bool {
	return this_.skipNonTranspilationPasses
}

// Java method: CompilerOptions.setDevMode(com.google.javascript.jscomp.CompilerOptions.DevMode devMode)
func (this_ *CompilerOptions) SetDevMode(devMode CompilerOptions_DevMode) {
	_ = devMode
	this_.devMode = devMode
}

// Java method: CompilerOptions.getDevMode()
func (this_ *CompilerOptions) GetDevMode() CompilerOptions_DevMode {
	return this_.devMode
}

// Java method: CompilerOptions.setCheckDeterminism(boolean checkDeterminism)
func (this_ *CompilerOptions) SetCheckDeterminism(checkDeterminism bool) {
	_ = checkDeterminism
	this_.checkDeterminism = checkDeterminism
}

// Java method: CompilerOptions.getCheckDeterminism()
func (this_ *CompilerOptions) GetCheckDeterminism() bool {
	return this_.checkDeterminism
}

// Java method: CompilerOptions.setMessageBundle(com.google.javascript.jscomp.MessageBundle messageBundle)
func (this_ *CompilerOptions) SetMessageBundle(messageBundle MessageBundle) {
	_ = messageBundle
	this_.messageBundle = messageBundle
}

// Java method: CompilerOptions.getMessageBundle()
func (this_ *CompilerOptions) GetMessageBundle() MessageBundle {
	return this_.messageBundle
}

// Java method: CompilerOptions.setCheckSymbols(boolean checkSymbols)
func (this_ *CompilerOptions) SetCheckSymbols(checkSymbols bool) {
	_ = checkSymbols
	this_.checkSymbols = checkSymbols
}

// Java method: CompilerOptions.getCheckSymbols()
func (this_ *CompilerOptions) GetCheckSymbols() bool {
	return this_.checkSymbols
}

// Java method: CompilerOptions.setCheckSuspiciousCode(boolean checkSuspiciousCode)
func (this_ *CompilerOptions) SetCheckSuspiciousCode(checkSuspiciousCode bool) {
	_ = checkSuspiciousCode
	this_.checkSuspiciousCode = checkSuspiciousCode
}

// Java method: CompilerOptions.getCheckSuspiciousCode()
func (this_ *CompilerOptions) GetCheckSuspiciousCode() bool {
	return this_.checkSuspiciousCode
}

// Java method: CompilerOptions.setCheckTypes(boolean checkTypes)
func (this_ *CompilerOptions) SetCheckTypes(checkTypes bool) {
	_ = checkTypes
	this_.checkTypes = checkTypes
}

// Java method: CompilerOptions.getCheckTypes()
func (this_ *CompilerOptions) GetCheckTypes() bool {
	return this_.checkTypes
}

// Java method: CompilerOptions.setFoldConstants(boolean foldConstants)
func (this_ *CompilerOptions) SetFoldConstants(foldConstants bool) {
	_ = foldConstants
	this_.foldConstants = foldConstants
}

// Java method: CompilerOptions.shouldFoldConstants()
func (this_ *CompilerOptions) ShouldFoldConstants() bool {
	return this_.foldConstants
}

// Java method: CompilerOptions.setDeadAssignmentElimination(boolean deadAssignmentElimination)
func (this_ *CompilerOptions) SetDeadAssignmentElimination(deadAssignmentElimination bool) {
	_ = deadAssignmentElimination
	this_.deadAssignmentElimination = deadAssignmentElimination
}

// Java method: CompilerOptions.shouldRunDeadAssignmentElimination()
func (this_ *CompilerOptions) ShouldRunDeadAssignmentElimination() bool {
	return this_.deadAssignmentElimination
}

// Java method: CompilerOptions.setDeadPropertyAssignmentElimination(boolean deadPropertyAssignmentElimination)
func (this_ *CompilerOptions) SetDeadPropertyAssignmentElimination(deadPropertyAssignmentElimination bool) {
	_ = deadPropertyAssignmentElimination
	this_.deadPropertyAssignmentElimination = pkg_base.Tri_ForBoolean(deadPropertyAssignmentElimination)
}

// Java method: CompilerOptions.shouldRunDeadPropertyAssignmentElimination()
func (this_ *CompilerOptions) ShouldRunDeadPropertyAssignmentElimination() bool {
	return this_.deadPropertyAssignmentElimination.ToBoolean((this_.deadAssignmentElimination && !(this_.polymerPass)))
}

// Java method: CompilerOptions.setInlineConstantVars(boolean inlineConstantVars)
func (this_ *CompilerOptions) SetInlineConstantVars(inlineConstantVars bool) {
	_ = inlineConstantVars
	this_.inlineConstantVars = inlineConstantVars
}

// Java method: CompilerOptions.shouldInlineConstantVars()
func (this_ *CompilerOptions) ShouldInlineConstantVars() bool {
	return this_.inlineConstantVars
}

// Java method: CompilerOptions.setCrossChunkCodeMotion(boolean crossChunkCodeMotion)
func (this_ *CompilerOptions) SetCrossChunkCodeMotion(crossChunkCodeMotion bool) {
	_ = crossChunkCodeMotion
	this_.crossChunkCodeMotion = crossChunkCodeMotion
}

// Java method: CompilerOptions.setCrossChunkCodeMotionNoStubMethods(boolean crossChunkCodeMotionNoStubMethods)
func (this_ *CompilerOptions) SetCrossChunkCodeMotionNoStubMethods(crossChunkCodeMotionNoStubMethods bool) {
	_ = crossChunkCodeMotionNoStubMethods
	this_.crossChunkCodeMotionNoStubMethods = crossChunkCodeMotionNoStubMethods
}

// Java method: CompilerOptions.getCrossChunkCodeMotionNoStubMethods()
func (this_ *CompilerOptions) GetCrossChunkCodeMotionNoStubMethods() bool {
	return this_.crossChunkCodeMotionNoStubMethods
}

// Java method: CompilerOptions.setParentChunkCanSeeSymbolsDeclaredInChildren(boolean parentChunkCanSeeSymbolsDeclaredInChildren)
func (this_ *CompilerOptions) SetParentChunkCanSeeSymbolsDeclaredInChildren(parentChunkCanSeeSymbolsDeclaredInChildren bool) {
	_ = parentChunkCanSeeSymbolsDeclaredInChildren
	this_.parentChunkCanSeeSymbolsDeclaredInChildren = parentChunkCanSeeSymbolsDeclaredInChildren
}

// Java method: CompilerOptions.getParentChunkCanSeeSymbolsDeclaredInChildren()
func (this_ *CompilerOptions) GetParentChunkCanSeeSymbolsDeclaredInChildren() bool {
	return this_.parentChunkCanSeeSymbolsDeclaredInChildren
}

// Java method: CompilerOptions.setCrossChunkMethodMotion(boolean crossChunkMethodMotion)
func (this_ *CompilerOptions) SetCrossChunkMethodMotion(crossChunkMethodMotion bool) {
	_ = crossChunkMethodMotion
	this_.crossChunkMethodMotion = crossChunkMethodMotion
}

// Java method: CompilerOptions.getCrossChunkMethodMotion()
func (this_ *CompilerOptions) GetCrossChunkMethodMotion() bool {
	return this_.crossChunkMethodMotion
}

// Java method: CompilerOptions.setCoalesceVariableNames(boolean coalesceVariableNames)
func (this_ *CompilerOptions) SetCoalesceVariableNames(coalesceVariableNames bool) {
	_ = coalesceVariableNames
	this_.coalesceVariableNames = coalesceVariableNames
}

// Java method: CompilerOptions.shouldCoalesceVariableNames()
func (this_ *CompilerOptions) ShouldCoalesceVariableNames() bool {
	return this_.coalesceVariableNames
}

// Java method: CompilerOptions.setOptimizeLetAndConst(boolean optimizeLetAndConst)
func (this_ *CompilerOptions) SetOptimizeLetAndConst(optimizeLetAndConst bool) {
	_ = optimizeLetAndConst
	this_.optimizeLetAndConst = optimizeLetAndConst
}

// Java method: CompilerOptions.shouldOptimizeLetAndConst()
func (this_ *CompilerOptions) ShouldOptimizeLetAndConst() bool {
	return this_.optimizeLetAndConst
}

// Java method: CompilerOptions.setAssumeGlobalScopeIsIsolated(boolean assumeGlobalScopeIsIsolated)
func (this_ *CompilerOptions) SetAssumeGlobalScopeIsIsolated(assumeGlobalScopeIsIsolated bool) {
	_ = assumeGlobalScopeIsIsolated
	this_.assumeGlobalScopeIsIsolated = assumeGlobalScopeIsIsolated
}

// Java method: CompilerOptions.shouldTreatGlobalScopeAsIsolated()
func (this_ *CompilerOptions) ShouldTreatGlobalScopeAsIsolated() bool {
	return (this_.assumeGlobalScopeIsIsolated || (this_.renamePrefixNamespace != j2gort.NullString))
}

// Java method: CompilerOptions.setInlineLocalVariables(boolean inlineLocalVariables)
func (this_ *CompilerOptions) SetInlineLocalVariables(inlineLocalVariables bool) {
	_ = inlineLocalVariables
	this_.inlineLocalVariables = inlineLocalVariables
}

// Java method: CompilerOptions.shouldInlineLocalVariables()
func (this_ *CompilerOptions) ShouldInlineLocalVariables() bool {
	return this_.inlineLocalVariables
}

// Java method: CompilerOptions.setFlowSensitiveInlineVariables(boolean enabled)
func (this_ *CompilerOptions) SetFlowSensitiveInlineVariables(enabled bool) {
	_ = enabled
	this_.flowSensitiveInlineVariables = enabled
}

// Java method: CompilerOptions.getFlowSensitiveInlineVariables()
func (this_ *CompilerOptions) GetFlowSensitiveInlineVariables() bool {
	return this_.flowSensitiveInlineVariables
}

// Java method: CompilerOptions.setSmartNameRemoval(boolean smartNameRemoval)
func (this_ *CompilerOptions) SetSmartNameRemoval(smartNameRemoval bool) {
	_ = smartNameRemoval
	this_.smartNameRemoval = smartNameRemoval
	if smartNameRemoval {
		this_.removeUnusedVars = true
		this_.removeUnusedPrototypeProperties = true
	}
}

// Java method: CompilerOptions.getSmartNameRemoval()
func (this_ *CompilerOptions) GetSmartNameRemoval() bool {
	return this_.smartNameRemoval
}

// Java method: CompilerOptions.setExtractPrototypeMemberDeclarations(boolean enabled)
func (this_ *CompilerOptions) SetExtractPrototypeMemberDeclarations_1_bool(enabled bool) {
	_ = enabled
	if enabled {
		this_.extractPrototypeMemberDeclarations = StaticGet_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP()
	} else {
		this_.extractPrototypeMemberDeclarations = StaticGet_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF()
	}
}

// Java method: CompilerOptions.setExtractPrototypeMemberDeclarations(com.google.javascript.jscomp.CompilerOptions.ExtractPrototypeMemberDeclarationsMode mode)
func (this_ *CompilerOptions) SetExtractPrototypeMemberDeclarations_1_CompilerOptions_ExtractPrototypeMemberDeclarationsMode(mode CompilerOptions_ExtractPrototypeMemberDeclarationsMode) {
	_ = mode
	this_.extractPrototypeMemberDeclarations = mode
}

// Java method: CompilerOptions.getExtractPrototypeMemberDeclarationsMode()
func (this_ *CompilerOptions) GetExtractPrototypeMemberDeclarationsMode() CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	return this_.extractPrototypeMemberDeclarations
}

// Java method: CompilerOptions.setRemoveUnusedPrototypeProperties(boolean enabled)
func (this_ *CompilerOptions) SetRemoveUnusedPrototypeProperties(enabled bool) {
	_ = enabled
	this_.removeUnusedPrototypeProperties = enabled
	this_.inlineGetters = enabled
}

// Java method: CompilerOptions.shouldRemoveUnusedPrototypeProperties()
func (this_ *CompilerOptions) ShouldRemoveUnusedPrototypeProperties() bool {
	return this_.removeUnusedPrototypeProperties
}

// Java method: CompilerOptions.shouldInlineGetters()
func (this_ *CompilerOptions) ShouldInlineGetters() bool {
	return this_.inlineGetters
}

// Java method: CompilerOptions.setCollapseVariableDeclarations(boolean enabled)
func (this_ *CompilerOptions) SetCollapseVariableDeclarations(enabled bool) {
	_ = enabled
	this_.collapseVariableDeclarations = enabled
}

// Java method: CompilerOptions.shouldCollapseVariableDeclarations()
func (this_ *CompilerOptions) ShouldCollapseVariableDeclarations() bool {
	return this_.collapseVariableDeclarations
}

// Java method: CompilerOptions.setCollapseAnonymousFunctions(boolean enabled)
func (this_ *CompilerOptions) SetCollapseAnonymousFunctions(enabled bool) {
	_ = enabled
	this_.collapseAnonymousFunctions = enabled
}

// Java method: CompilerOptions.shouldCollapseAnonymousFunctions()
func (this_ *CompilerOptions) ShouldCollapseAnonymousFunctions() bool {
	return this_.collapseAnonymousFunctions
}

// Java method: CompilerOptions.setAliasStringsMode(com.google.javascript.jscomp.CompilerOptions.AliasStringsMode aliasStringsMode)
func (this_ *CompilerOptions) SetAliasStringsMode(aliasStringsMode CompilerOptions_AliasStringsMode) {
	_ = aliasStringsMode
	this_.aliasStringsMode = aliasStringsMode
}

// Java method: CompilerOptions.getAliasStringsMode()
func (this_ *CompilerOptions) GetAliasStringsMode() CompilerOptions_AliasStringsMode {
	return this_.aliasStringsMode
}

// Java method: CompilerOptions.setOutputJsStringUsage(boolean outputJsStringUsage)
func (this_ *CompilerOptions) SetOutputJsStringUsage(outputJsStringUsage bool) {
	_ = outputJsStringUsage
	this_.outputJsStringUsage = outputJsStringUsage
}

// Java method: CompilerOptions.shouldOutputJsStringUsage()
func (this_ *CompilerOptions) ShouldOutputJsStringUsage() bool {
	return this_.outputJsStringUsage
}

// Java method: CompilerOptions.setConvertToDottedProperties(boolean convertToDottedProperties)
func (this_ *CompilerOptions) SetConvertToDottedProperties(convertToDottedProperties bool) {
	_ = convertToDottedProperties
	this_.convertToDottedProperties = convertToDottedProperties
}

// Java method: CompilerOptions.shouldConvertToDottedProperties()
func (this_ *CompilerOptions) ShouldConvertToDottedProperties() bool {
	return this_.convertToDottedProperties
}

// Java method: CompilerOptions.setUseTypesForLocalOptimization(boolean useTypesForLocalOptimization)
func (this_ *CompilerOptions) SetUseTypesForLocalOptimization(useTypesForLocalOptimization bool) {
	_ = useTypesForLocalOptimization
	this_.useTypesForLocalOptimization = useTypesForLocalOptimization
}

// Java method: CompilerOptions.shouldUseTypesForLocalOptimization()
func (this_ *CompilerOptions) ShouldUseTypesForLocalOptimization() bool {
	return this_.useTypesForLocalOptimization
}

// Java method: CompilerOptions.setUseTypesForOptimization(boolean useTypesForOptimization)
func (this_ *CompilerOptions) SetUseTypesForOptimization(useTypesForOptimization bool) {
	_ = useTypesForOptimization
	if useTypesForOptimization {
		this_.disambiguateProperties = true
		this_.ambiguateProperties = true
		this_.inlineProperties = true
		this_.useTypesForLocalOptimization = true
	}
}

// Java method: CompilerOptions.requiresTypesForOptimization()
func (this_ *CompilerOptions) RequiresTypesForOptimization() bool {
	return ((((this_.disambiguateProperties || this_.ambiguateProperties) || this_.inlineProperties) || this_.useTypesForLocalOptimization) || (this_.closurePass && (this_.removeAbstractMethods || this_.removeClosureAsserts)))
}

// Java method: CompilerOptions.setRewriteFunctionExpressions(boolean rewriteFunctionExpressions)
func (this_ *CompilerOptions) SetRewriteFunctionExpressions(rewriteFunctionExpressions bool) {
	_ = rewriteFunctionExpressions
	this_.rewriteFunctionExpressions = rewriteFunctionExpressions
}

// Java method: CompilerOptions.shouldRewriteFunctionExpressions()
func (this_ *CompilerOptions) ShouldRewriteFunctionExpressions() bool {
	return this_.rewriteFunctionExpressions
}

// Java method: CompilerOptions.setOptimizeCalls(boolean optimizeCalls)
func (this_ *CompilerOptions) SetOptimizeCalls(optimizeCalls bool) {
	_ = optimizeCalls
	this_.optimizeCalls = optimizeCalls
}

// Java method: CompilerOptions.shouldOptimizeCalls()
func (this_ *CompilerOptions) ShouldOptimizeCalls() bool {
	return this_.optimizeCalls
}

// Java method: CompilerOptions.getOptimizeESClassConstructors()
func (this_ *CompilerOptions) GetOptimizeESClassConstructors() bool {
	return this_.optimizeESClassConstructors
}

// Java method: CompilerOptions.setOptimizeESClassConstructors(boolean optimizeESClassConstructors)
func (this_ *CompilerOptions) SetOptimizeESClassConstructors(optimizeESClassConstructors bool) {
	_ = optimizeESClassConstructors
	this_.optimizeESClassConstructors = optimizeESClassConstructors
}

// Java method: CompilerOptions.setVariableRenaming(com.google.javascript.jscomp.VariableRenamingPolicy variableRenaming)
func (this_ *CompilerOptions) SetVariableRenaming(variableRenaming VariableRenamingPolicy) {
	_ = variableRenaming
	this_.variableRenaming = variableRenaming
}

// Java method: CompilerOptions.getVariableRenaming()
func (this_ *CompilerOptions) GetVariableRenaming() VariableRenamingPolicy {
	return this_.variableRenaming
}

// Java method: CompilerOptions.setPropertyRenaming(com.google.javascript.jscomp.PropertyRenamingPolicy propertyRenaming)
func (this_ *CompilerOptions) SetPropertyRenaming(propertyRenaming PropertyRenamingPolicy) {
	_ = propertyRenaming
	this_.propertyRenaming = propertyRenaming
}

// Java method: CompilerOptions.getPropertyRenaming()
func (this_ *CompilerOptions) GetPropertyRenaming() PropertyRenamingPolicy {
	return this_.propertyRenaming
}

// Java method: CompilerOptions.setLabelRenaming(boolean labelRenaming)
func (this_ *CompilerOptions) SetLabelRenaming(labelRenaming bool) {
	_ = labelRenaming
	this_.labelRenaming = labelRenaming
}

// Java method: CompilerOptions.shouldRenameLabels()
func (this_ *CompilerOptions) ShouldRenameLabels() bool {
	return this_.labelRenaming
}

// Java method: CompilerOptions.setReserveRawExports(boolean reserveRawExports)
func (this_ *CompilerOptions) SetReserveRawExports(reserveRawExports bool) {
	_ = reserveRawExports
	this_.reserveRawExports = reserveRawExports
}

// Java method: CompilerOptions.shouldReserveRawExports()
func (this_ *CompilerOptions) ShouldReserveRawExports() bool {
	return this_.reserveRawExports
}

// Java method: CompilerOptions.setPreferStableNames(boolean preferStableNames)
func (this_ *CompilerOptions) SetPreferStableNames(preferStableNames bool) {
	_ = preferStableNames
	this_.preferStableNames = preferStableNames
}

// Java method: CompilerOptions.shouldPreferStableNames()
func (this_ *CompilerOptions) ShouldPreferStableNames() bool {
	return this_.preferStableNames
}

// Java method: CompilerOptions.setGeneratePseudoNames(boolean generatePseudoNames)
func (this_ *CompilerOptions) SetGeneratePseudoNames(generatePseudoNames bool) {
	_ = generatePseudoNames
	this_.generatePseudoNames = generatePseudoNames
}

// Java method: CompilerOptions.shouldGeneratePseudoNames()
func (this_ *CompilerOptions) ShouldGeneratePseudoNames() bool {
	return this_.generatePseudoNames
}

// Java method: CompilerOptions.setPropertyRenamingOnlyCompilationMode(boolean propertyRenamingOnlyCompilationMode)
func (this_ *CompilerOptions) SetPropertyRenamingOnlyCompilationMode(propertyRenamingOnlyCompilationMode bool) {
	_ = propertyRenamingOnlyCompilationMode
	this_.propertyRenamingOnlyCompilationMode = propertyRenamingOnlyCompilationMode
}

// Java method: CompilerOptions.isPropertyRenamingOnlyCompilationMode()
func (this_ *CompilerOptions) IsPropertyRenamingOnlyCompilationMode() bool {
	return this_.propertyRenamingOnlyCompilationMode
}

// Java method: CompilerOptions.setRenamePrefix(java.lang.String renamePrefix)
func (this_ *CompilerOptions) SetRenamePrefix(renamePrefix string) {
	_ = renamePrefix
	this_.renamePrefix = renamePrefix
}

// Java method: CompilerOptions.getRenamePrefix()
func (this_ *CompilerOptions) GetRenamePrefix() string {
	return this_.renamePrefix
}

// Java method: CompilerOptions.getRenamePrefixNamespace()
func (this_ *CompilerOptions) GetRenamePrefixNamespace() string {
	return this_.renamePrefixNamespace
}

// Java method: CompilerOptions.setRenamePrefixNamespace(java.lang.String renamePrefixNamespace)
func (this_ *CompilerOptions) SetRenamePrefixNamespace(renamePrefixNamespace string) {
	_ = renamePrefixNamespace
	this_.renamePrefixNamespace = renamePrefixNamespace
}

// Java method: CompilerOptions.setCollapsePropertiesLevel(com.google.javascript.jscomp.CompilerOptions.PropertyCollapseLevel level)
func (this_ *CompilerOptions) SetCollapsePropertiesLevel(level CompilerOptions_PropertyCollapseLevel) {
	_ = level
	this_.collapsePropertiesLevel = level
}

// Java method: CompilerOptions.setCollapseProperties(boolean fullyCollapse)
func (this_ *CompilerOptions) SetCollapseProperties(fullyCollapse bool) {
	_ = fullyCollapse
	if fullyCollapse {
		this_.collapsePropertiesLevel = StaticGet_CompilerOptions_PropertyCollapseLevel_ALL()
	} else {
		this_.collapsePropertiesLevel = StaticGet_CompilerOptions_PropertyCollapseLevel_NONE()
	}
}

// Java method: CompilerOptions.getCollapsePropertiesLevel()
func (this_ *CompilerOptions) GetCollapsePropertiesLevel() CompilerOptions_PropertyCollapseLevel {
	return this_.collapsePropertiesLevel
}

// Java method: CompilerOptions.setDevirtualizeMethods(boolean devirtualizeMethods)
func (this_ *CompilerOptions) SetDevirtualizeMethods(devirtualizeMethods bool) {
	_ = devirtualizeMethods
	this_.devirtualizeMethods = devirtualizeMethods
}

// Java method: CompilerOptions.shouldDevirtualizeMethods()
func (this_ *CompilerOptions) ShouldDevirtualizeMethods() bool {
	return this_.devirtualizeMethods
}

// Java method: CompilerOptions.setComputeFunctionSideEffects(boolean computeFunctionSideEffects)
func (this_ *CompilerOptions) SetComputeFunctionSideEffects(computeFunctionSideEffects bool) {
	_ = computeFunctionSideEffects
	this_.computeFunctionSideEffects = computeFunctionSideEffects
}

// Java method: CompilerOptions.shouldComputeFunctionSideEffects()
func (this_ *CompilerOptions) ShouldComputeFunctionSideEffects() bool {
	return this_.computeFunctionSideEffects
}

// Java method: CompilerOptions.setDisambiguateProperties(boolean disambiguateProperties)
func (this_ *CompilerOptions) SetDisambiguateProperties(disambiguateProperties bool) {
	_ = disambiguateProperties
	this_.disambiguateProperties = disambiguateProperties
}

// Java method: CompilerOptions.shouldDisambiguateProperties()
func (this_ *CompilerOptions) ShouldDisambiguateProperties() bool {
	return this_.disambiguateProperties
}

// Java method: CompilerOptions.setAmbiguateProperties(boolean ambiguateProperties)
func (this_ *CompilerOptions) SetAmbiguateProperties(ambiguateProperties bool) {
	_ = ambiguateProperties
	this_.ambiguateProperties = ambiguateProperties
}

// Java method: CompilerOptions.shouldAmbiguateProperties()
func (this_ *CompilerOptions) ShouldAmbiguateProperties() bool {
	return this_.ambiguateProperties
}

// Java method: CompilerOptions.setInputVariableMap(com.google.javascript.jscomp.VariableMap inputVariableMap)
func (this_ *CompilerOptions) SetInputVariableMap(inputVariableMap *VariableMap) {
	_ = inputVariableMap
	this_.inputVariableMap = inputVariableMap
}

// Java method: CompilerOptions.getInputVariableMap()
func (this_ *CompilerOptions) GetInputVariableMap() *VariableMap {
	return this_.inputVariableMap
}

// Java method: CompilerOptions.setInputPropertyMap(com.google.javascript.jscomp.VariableMap inputPropertyMap)
func (this_ *CompilerOptions) SetInputPropertyMap(inputPropertyMap *VariableMap) {
	_ = inputPropertyMap
	this_.inputPropertyMap = inputPropertyMap
}

// Java method: CompilerOptions.getInputPropertyMap()
func (this_ *CompilerOptions) GetInputPropertyMap() *VariableMap {
	return this_.inputPropertyMap
}

// Java method: CompilerOptions.setExportTestFunctions(boolean exportTestFunctions)
func (this_ *CompilerOptions) SetExportTestFunctions(exportTestFunctions bool) {
	_ = exportTestFunctions
	this_.exportTestFunctions = exportTestFunctions
}

// Java method: CompilerOptions.shouldExportTestFunctions()
func (this_ *CompilerOptions) ShouldExportTestFunctions() bool {
	return this_.exportTestFunctions
}

// Java method: CompilerOptions.setSyntheticBlockStartMarker(java.lang.String syntheticBlockStartMarker)
func (this_ *CompilerOptions) SetSyntheticBlockStartMarker(syntheticBlockStartMarker string) {
	_ = syntheticBlockStartMarker
	this_.syntheticBlockStartMarker = syntheticBlockStartMarker
}

// Java method: CompilerOptions.setSyntheticBlockEndMarker(java.lang.String syntheticBlockEndMarker)
func (this_ *CompilerOptions) SetSyntheticBlockEndMarker(syntheticBlockEndMarker string) {
	_ = syntheticBlockEndMarker
	this_.syntheticBlockEndMarker = syntheticBlockEndMarker
}

// Java method: CompilerOptions.getSyntheticBlockStartMarker()
func (this_ *CompilerOptions) GetSyntheticBlockStartMarker() string {
	return this_.syntheticBlockStartMarker
}

// Java method: CompilerOptions.getSyntheticBlockEndMarker()
func (this_ *CompilerOptions) GetSyntheticBlockEndMarker() string {
	return this_.syntheticBlockEndMarker
}

// Java method: CompilerOptions.setLocale(java.lang.String locale)
func (this_ *CompilerOptions) SetLocale(locale string) {
	_ = locale
	this_.locale = locale
}

// Java method: CompilerOptions.getLocale()
func (this_ *CompilerOptions) GetLocale() string {
	return this_.locale
}

// Java method: CompilerOptions.setDoLateLocalization(boolean doLateLocalization)
func (this_ *CompilerOptions) SetDoLateLocalization(doLateLocalization bool) {
	_ = doLateLocalization
	this_.doLateLocalization = doLateLocalization
}

// Java method: CompilerOptions.doLateLocalization()
func (this_ *CompilerOptions) DoLateLocalization() bool {
	return this_.doLateLocalization
}

// Java method: CompilerOptions.shouldRunReplaceMessagesPass()
func (this_ *CompilerOptions) ShouldRunReplaceMessagesPass() bool {
	return (!(this_.ShouldRunReplaceMessagesForChrome()) && (this_.messageBundle != nil && !this_.messageBundle.IsJavaNil_()))
}

// Java method: CompilerOptions.setMarkAsCompiled(boolean markAsCompiled)
func (this_ *CompilerOptions) SetMarkAsCompiled(markAsCompiled bool) {
	_ = markAsCompiled
	this_.markAsCompiled = markAsCompiled
}

// Java method: CompilerOptions.shouldMarkAsCompiled()
func (this_ *CompilerOptions) ShouldMarkAsCompiled() bool {
	return this_.markAsCompiled
}

// Java method: CompilerOptions.setClosurePass(boolean closurePass)
func (this_ *CompilerOptions) SetClosurePass(closurePass bool) {
	_ = closurePass
	this_.closurePass = closurePass
}

// Java method: CompilerOptions.getClosurePass()
func (this_ *CompilerOptions) GetClosurePass() bool {
	return this_.closurePass
}

// Java method: CompilerOptions.setPreserveClosurePrimitives(boolean preserveClosurePrimitives)
func (this_ *CompilerOptions) SetPreserveClosurePrimitives(preserveClosurePrimitives bool) {
	_ = preserveClosurePrimitives
	this_.preserveClosurePrimitives = preserveClosurePrimitives
}

// Java method: CompilerOptions.shouldPreservesGoogProvidesAndRequires()
func (this_ *CompilerOptions) ShouldPreservesGoogProvidesAndRequires() bool {
	return this_.preserveClosurePrimitives
}

// Java method: CompilerOptions.shouldPreserveGoogModule()
func (this_ *CompilerOptions) ShouldPreserveGoogModule() bool {
	return this_.preserveClosurePrimitives
}

// Java method: CompilerOptions.shouldPreserveGoogLibraryPrimitives()
func (this_ *CompilerOptions) ShouldPreserveGoogLibraryPrimitives() bool {
	return this_.preserveClosurePrimitives
}

// Java method: CompilerOptions.setPreserveTypeAnnotations(boolean preserveTypeAnnotations)
func (this_ *CompilerOptions) SetPreserveTypeAnnotations(preserveTypeAnnotations bool) {
	_ = preserveTypeAnnotations
	this_.preserveTypeAnnotations = preserveTypeAnnotations
}

// Java method: CompilerOptions.shouldPreserveTypeAnnotations()
func (this_ *CompilerOptions) ShouldPreserveTypeAnnotations() bool {
	return this_.preserveTypeAnnotations
}

// Java method: CompilerOptions.setGentsMode(boolean gentsMode)
func (this_ *CompilerOptions) SetGentsMode(gentsMode bool) {
	_ = gentsMode
	this_.gentsMode = gentsMode
}

// Java method: CompilerOptions.getGentsMode()
func (this_ *CompilerOptions) GetGentsMode() bool {
	return this_.gentsMode
}

// Java method: CompilerOptions.setGatherCssNames(boolean gatherCssNames)
func (this_ *CompilerOptions) SetGatherCssNames(gatherCssNames bool) {
	_ = gatherCssNames
	this_.gatherCssNames = gatherCssNames
}

// Java method: CompilerOptions.shouldGatherCssNames()
func (this_ *CompilerOptions) ShouldGatherCssNames() bool {
	return this_.gatherCssNames
}

// Java method: CompilerOptions.setStripTypes(java.util.Set<java.lang.String> stripTypes)
func (this_ *CompilerOptions) SetStripTypes(stripTypes j2gort.Ext_java_util_Set) {
	_ = stripTypes
	this_.stripTypes = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](stripTypes))
}

// Java method: CompilerOptions.getStripTypes()
func (this_ *CompilerOptions) GetStripTypes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.stripTypes
}

// Java method: CompilerOptions.setStripNameSuffixes(java.util.Set<java.lang.String> stripNameSuffixes)
func (this_ *CompilerOptions) SetStripNameSuffixes(stripNameSuffixes j2gort.Ext_java_util_Set) {
	_ = stripNameSuffixes
	this_.stripNameSuffixes = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](stripNameSuffixes))
}

// Java method: CompilerOptions.getStripNameSuffixes()
func (this_ *CompilerOptions) GetStripNameSuffixes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.stripNameSuffixes
}

// Java method: CompilerOptions.setStripNamePrefixes(java.util.Set<java.lang.String> stripNamePrefixes)
func (this_ *CompilerOptions) SetStripNamePrefixes(stripNamePrefixes j2gort.Ext_java_util_Set) {
	_ = stripNamePrefixes
	this_.stripNamePrefixes = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](stripNamePrefixes))
}

// Java method: CompilerOptions.getStripNamePrefixes()
func (this_ *CompilerOptions) GetStripNamePrefixes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.stripNamePrefixes
}

// Java method: CompilerOptions.addCustomPass(com.google.javascript.jscomp.CustomPassExecutionTime time, com.google.javascript.jscomp.CompilerPass customPass)
func (this_ *CompilerOptions) AddCustomPass(time CustomPassExecutionTime, customPass CompilerPass) {
	_ = time
	_ = customPass
	if this_.customPasses == nil {
		this_.customPasses = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0())
	}
	this_.customPasses.Put(j2gort.BoxNullableVal(time), customPass)
}

// Java method: CompilerOptions.getCustomPassesAt(com.google.javascript.jscomp.CustomPassExecutionTime executionTime)
func (this_ *CompilerOptions) GetCustomPassesAt(executionTime CustomPassExecutionTime) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = executionTime
	if this_.customPasses == nil {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.customPasses.Get(j2gort.BoxNullableVal(executionTime))))
	}
}

// Java method: CompilerOptions.setDefineReplacements(java.util.Map<java.lang.String,java.lang.Object> defineReplacements)
func (this_ *CompilerOptions) SetDefineReplacements(defineReplacements j2gort.Ext_java_util_Map) {
	_ = defineReplacements
	this_.defineReplacements.Clear()
	this_.defineReplacements.PutAll(defineReplacements)
}

// Java method: CompilerOptions.getEnableZonesDefineName()
func (this_ *CompilerOptions) GetEnableZonesDefineName() string {
	return this_.enableZonesDefineName
}

// Java method: CompilerOptions.setEnableZonesDefineName(java.lang.@org.jspecify.annotations.Nullable String enableZonesDefineName)
func (this_ *CompilerOptions) SetEnableZonesDefineName(enableZonesDefineName string) {
	_ = enableZonesDefineName
	this_.enableZonesDefineName = enableZonesDefineName
}

// Java method: CompilerOptions.getZoneInputPattern()
func (this_ *CompilerOptions) GetZoneInputPattern() *j2gort.Ext_java_util_regex_Pattern {
	return this_.zoneInputPattern
}

// Java method: CompilerOptions.setZoneInputPattern(java.util.regex.@org.jspecify.annotations.Nullable Pattern zoneInputPattern)
func (this_ *CompilerOptions) SetZoneInputPattern(zoneInputPattern *j2gort.Ext_java_util_regex_Pattern) {
	_ = zoneInputPattern
	this_.zoneInputPattern = zoneInputPattern
}

// Java method: CompilerOptions.setRewriteGlobalDeclarationsForTryCatchWrapping(boolean rewrite)
func (this_ *CompilerOptions) SetRewriteGlobalDeclarationsForTryCatchWrapping(rewrite bool) {
	_ = rewrite
	this_.rewriteGlobalDeclarationsForTryCatchWrapping = rewrite
}

// Java method: CompilerOptions.shouldRewriteGlobalDeclarationsForTryCatchWrapping()
func (this_ *CompilerOptions) ShouldRewriteGlobalDeclarationsForTryCatchWrapping() bool {
	return this_.rewriteGlobalDeclarationsForTryCatchWrapping
}

// Java method: CompilerOptions.setCssRenamingMap(com.google.javascript.jscomp.CssRenamingMap cssRenamingMap)
func (this_ *CompilerOptions) SetCssRenamingMap(cssRenamingMap CssRenamingMap) {
	_ = cssRenamingMap
	this_.cssRenamingMap = cssRenamingMap
}

// Java method: CompilerOptions.getCssRenamingMap()
func (this_ *CompilerOptions) GetCssRenamingMap() CssRenamingMap {
	return this_.cssRenamingMap
}

// Java method: CompilerOptions.setCssRenamingSkiplist(java.util.Set<java.lang.String> skiplist)
func (this_ *CompilerOptions) SetCssRenamingSkiplist(skiplist j2gort.Ext_java_util_Set) {
	_ = skiplist
	this_.cssRenamingSkiplist = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](skiplist))
}

// Java method: CompilerOptions.getCssRenamingSkiplist()
func (this_ *CompilerOptions) GetCssRenamingSkiplist() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.cssRenamingSkiplist
}

// Java method: CompilerOptions.setReplaceStringsFunctionDescriptions(java.util.List<java.lang.String> replaceStringsFunctionDescriptions)
func (this_ *CompilerOptions) SetReplaceStringsFunctionDescriptions(replaceStringsFunctionDescriptions j2gort.Ext_java_util_List) {
	_ = replaceStringsFunctionDescriptions
	this_.replaceStringsFunctionDescriptions = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](replaceStringsFunctionDescriptions))
}

// Java method: CompilerOptions.getReplaceStringsFunctionDescriptions()
func (this_ *CompilerOptions) GetReplaceStringsFunctionDescriptions() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.replaceStringsFunctionDescriptions
}

// Java method: CompilerOptions.setReplaceStringsPlaceholderToken(java.lang.String replaceStringsPlaceholderToken)
func (this_ *CompilerOptions) SetReplaceStringsPlaceholderToken(replaceStringsPlaceholderToken string) {
	_ = replaceStringsPlaceholderToken
	this_.replaceStringsPlaceholderToken = replaceStringsPlaceholderToken
}

// Java method: CompilerOptions.getReplaceStringsPlaceholderToken()
func (this_ *CompilerOptions) GetReplaceStringsPlaceholderToken() string {
	return this_.replaceStringsPlaceholderToken
}

// Java method: CompilerOptions.setPrettyPrint(boolean prettyPrint)
func (this_ *CompilerOptions) SetPrettyPrint(prettyPrint bool) {
	_ = prettyPrint
	this_.prettyPrint = prettyPrint
}

// Java method: CompilerOptions.isPrettyPrint()
func (this_ *CompilerOptions) IsPrettyPrint() bool {
	return this_.prettyPrint
}

// Java method: CompilerOptions.setLineBreak(boolean lineBreak)
func (this_ *CompilerOptions) SetLineBreak(lineBreak bool) {
	_ = lineBreak
	this_.lineBreak = lineBreak
}

// Java method: CompilerOptions.shouldAddLineBreak()
func (this_ *CompilerOptions) ShouldAddLineBreak() bool {
	return this_.lineBreak
}

// Java method: CompilerOptions.setPrintInputDelimiter(boolean printInputDelimiter)
func (this_ *CompilerOptions) SetPrintInputDelimiter(printInputDelimiter bool) {
	_ = printInputDelimiter
	this_.printInputDelimiter = printInputDelimiter
}

// Java method: CompilerOptions.shouldPrintInputDelimiter()
func (this_ *CompilerOptions) ShouldPrintInputDelimiter() bool {
	return this_.printInputDelimiter
}

// Java method: CompilerOptions.setInputDelimiter(java.lang.String inputDelimiter)
func (this_ *CompilerOptions) SetInputDelimiter(inputDelimiter string) {
	_ = inputDelimiter
	this_.inputDelimiter = inputDelimiter
}

// Java method: CompilerOptions.getInputDelimiter()
func (this_ *CompilerOptions) GetInputDelimiter() string {
	return this_.inputDelimiter
}

// Java method: CompilerOptions.setDebugLogDirectory(java.nio.file.@org.jspecify.annotations.Nullable Path dir)
func (this_ *CompilerOptions) SetDebugLogDirectory(dir j2gort.Ext_java_nio_file_Path) {
	_ = dir
	this_.debugLogDirectory = dir
}

// Java method: CompilerOptions.getDebugLogDirectory()
func (this_ *CompilerOptions) GetDebugLogDirectory() j2gort.Ext_java_nio_file_Path {
	return this_.debugLogDirectory
}

// Java method: CompilerOptions.setDebugLogFilter(java.lang.String filter)
func (this_ *CompilerOptions) SetDebugLogFilter(filter string) {
	_ = filter
	this_.debugLogFilter = filter
}

// Java method: CompilerOptions.getDebugLogFilter()
func (this_ *CompilerOptions) GetDebugLogFilter() string {
	return this_.debugLogFilter
}

// Java method: CompilerOptions.shouldSerializeExtraDebugInfo()
func (this_ *CompilerOptions) ShouldSerializeExtraDebugInfo() bool {
	return (this_.serializeExtraDebugInfo || (this_.GetDebugLogDirectory() != nil))
}

// Java method: CompilerOptions.setSerializeExtraDebugInfo(boolean serializeExtraDebugInfo)
func (this_ *CompilerOptions) SetSerializeExtraDebugInfo(serializeExtraDebugInfo bool) {
	_ = serializeExtraDebugInfo
	this_.serializeExtraDebugInfo = serializeExtraDebugInfo
}

// Java method: CompilerOptions.setStateCompressionWrapper(java.util.function.@org.jspecify.annotations.Nullable Function<java.io.OutputStream,java.io.OutputStream> wrapper)
func (this_ *CompilerOptions) SetStateCompressionWrapper(wrapper j2gort.Ext_java_util_function_Function) {
	_ = wrapper
	this_.stateCompressionWrapper = wrapper
}

// Java method: CompilerOptions.getStateCompressionWrapper()
func (this_ *CompilerOptions) GetStateCompressionWrapper() j2gort.Ext_java_util_function_Function {
	return this_.stateCompressionWrapper
}

// Java method: CompilerOptions.setStateDecompressionWrapper(java.util.function.@org.jspecify.annotations.Nullable Function<java.io.InputStream,java.io.InputStream> wrapper)
func (this_ *CompilerOptions) SetStateDecompressionWrapper(wrapper j2gort.Ext_java_util_function_Function) {
	_ = wrapper
	this_.stateDecompressionWrapper = wrapper
}

// Java method: CompilerOptions.getStateDecompressionWrapper()
func (this_ *CompilerOptions) GetStateDecompressionWrapper() j2gort.Ext_java_util_function_Function {
	return this_.stateDecompressionWrapper
}

// Java method: CompilerOptions.setQuoteKeywordProperties(boolean quoteKeywordProperties)
func (this_ *CompilerOptions) SetQuoteKeywordProperties(quoteKeywordProperties bool) {
	_ = quoteKeywordProperties
	this_.quoteKeywordProperties = quoteKeywordProperties
}

// Java method: CompilerOptions.shouldQuoteKeywordProperties()
func (this_ *CompilerOptions) ShouldQuoteKeywordProperties() bool {
	if this_.incrementalCheckMode == StaticGet_CompilerOptions_IncrementalCheckMode_GENERATE_IJS() {
		return false
	}
	return (this_.quoteKeywordProperties || pkg_parser.StaticGet_FeatureSet_ES3().Contains_1_FeatureSet(this_.GetOutputFeatureSet()))
}

// Java method: CompilerOptions.setErrorFormat(com.google.javascript.jscomp.ErrorFormat errorFormat)
func (this_ *CompilerOptions) SetErrorFormat(errorFormat *ErrorFormat) {
	_ = errorFormat
	this_.errorFormat = errorFormat
}

// Java method: CompilerOptions.getErrorFormat()
func (this_ *CompilerOptions) GetErrorFormat() *ErrorFormat {
	return this_.errorFormat
}

// Java method: CompilerOptions.setWarningsGuard(com.google.javascript.jscomp.ComposeWarningsGuard warningsGuard)
func (this_ *CompilerOptions) SetWarningsGuard(warningsGuard *ComposeWarningsGuard) {
	_ = warningsGuard
	this_.warningsGuard = warningsGuard
}

// Java method: CompilerOptions.setLineLengthThreshold(int lineLengthThreshold)
func (this_ *CompilerOptions) SetLineLengthThreshold(lineLengthThreshold int32) {
	_ = lineLengthThreshold
	this_.lineLengthThreshold = lineLengthThreshold
}

// Java method: CompilerOptions.getLineLengthThreshold()
func (this_ *CompilerOptions) GetLineLengthThreshold() int32 {
	return this_.lineLengthThreshold
}

// Java method: CompilerOptions.setUseOriginalNamesInOutput(boolean useOriginalNamesInOutput)
func (this_ *CompilerOptions) SetUseOriginalNamesInOutput(useOriginalNamesInOutput bool) {
	_ = useOriginalNamesInOutput
	this_.useOriginalNamesInOutput = useOriginalNamesInOutput
}

// Java method: CompilerOptions.getUseOriginalNamesInOutput()
func (this_ *CompilerOptions) GetUseOriginalNamesInOutput() bool {
	return this_.useOriginalNamesInOutput
}

// Java method: CompilerOptions.setExternExportsPath(java.lang.@org.jspecify.annotations.Nullable String externExportsPath)
func (this_ *CompilerOptions) SetExternExportsPath(externExportsPath string) {
	_ = externExportsPath
	this_.externExportsPath = externExportsPath
}

// Java method: CompilerOptions.getExternExportsPath()
func (this_ *CompilerOptions) GetExternExportsPath() string {
	return this_.externExportsPath
}

// Java method: CompilerOptions.setSourceMapOutputPath(java.lang.String sourceMapOutputPath)
func (this_ *CompilerOptions) SetSourceMapOutputPath(sourceMapOutputPath string) {
	_ = sourceMapOutputPath
	this_.sourceMapOutputPath = sourceMapOutputPath
}

// Java method: CompilerOptions.setApplyInputSourceMaps(boolean applyInputSourceMaps)
func (this_ *CompilerOptions) SetApplyInputSourceMaps(applyInputSourceMaps bool) {
	_ = applyInputSourceMaps
	this_.applyInputSourceMaps = applyInputSourceMaps
}

// Java method: CompilerOptions.getApplyInputSourceMaps()
func (this_ *CompilerOptions) GetApplyInputSourceMaps() bool {
	return this_.applyInputSourceMaps
}

// Java method: CompilerOptions.setResolveSourceMapAnnotations(boolean resolveSourceMapAnnotations)
func (this_ *CompilerOptions) SetResolveSourceMapAnnotations(resolveSourceMapAnnotations bool) {
	_ = resolveSourceMapAnnotations
	this_.resolveSourceMapAnnotations = resolveSourceMapAnnotations
}

// Java method: CompilerOptions.getResolveSourceMapAnnotations()
func (this_ *CompilerOptions) GetResolveSourceMapAnnotations() bool {
	return this_.resolveSourceMapAnnotations
}

// Java method: CompilerOptions.setSourceMapIncludeSourcesContent(boolean sourceMapIncludeSourcesContent)
func (this_ *CompilerOptions) SetSourceMapIncludeSourcesContent(sourceMapIncludeSourcesContent bool) {
	_ = sourceMapIncludeSourcesContent
	this_.sourceMapIncludeSourcesContent = sourceMapIncludeSourcesContent
}

// Java method: CompilerOptions.getSourceMapIncludeSourcesContent()
func (this_ *CompilerOptions) GetSourceMapIncludeSourcesContent() bool {
	return this_.sourceMapIncludeSourcesContent
}

// Java method: CompilerOptions.setParseInlineSourceMaps(boolean parseInlineSourceMaps)
func (this_ *CompilerOptions) SetParseInlineSourceMaps(parseInlineSourceMaps bool) {
	_ = parseInlineSourceMaps
	this_.parseInlineSourceMaps = parseInlineSourceMaps
}

// Java method: CompilerOptions.getParseInlineSourceMaps()
func (this_ *CompilerOptions) GetParseInlineSourceMaps() bool {
	return this_.parseInlineSourceMaps
}

// Java method: CompilerOptions.setSourceMapDetailLevel(com.google.javascript.jscomp.SourceMap.DetailLevel sourceMapDetailLevel)
func (this_ *CompilerOptions) SetSourceMapDetailLevel(sourceMapDetailLevel *SourceMap_DetailLevel) {
	_ = sourceMapDetailLevel
	this_.sourceMapDetailLevel = sourceMapDetailLevel
}

// Java method: CompilerOptions.getSourceMapDetailLevel()
func (this_ *CompilerOptions) GetSourceMapDetailLevel() *SourceMap_DetailLevel {
	return this_.sourceMapDetailLevel
}

// Java method: CompilerOptions.setSourceMapFormat(com.google.javascript.jscomp.SourceMap.Format sourceMapFormat)
func (this_ *CompilerOptions) SetSourceMapFormat(sourceMapFormat *SourceMap_Format) {
	_ = sourceMapFormat
	this_.sourceMapFormat = sourceMapFormat
}

// Java method: CompilerOptions.getSourceMapFormat()
func (this_ *CompilerOptions) GetSourceMapFormat() *SourceMap_Format {
	return this_.sourceMapFormat
}

// Java method: CompilerOptions.setSourceMapLocationMappings(java.util.List<? extends com.google.javascript.jscomp.SourceMap.LocationMapping> sourceMapLocationMappings)
func (this_ *CompilerOptions) SetSourceMapLocationMappings(sourceMapLocationMappings j2gort.Ext_java_util_List) {
	_ = sourceMapLocationMappings
	this_.sourceMapLocationMappings = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](sourceMapLocationMappings))
}

// Java method: CompilerOptions.getSourceMapLocationMappings()
func (this_ *CompilerOptions) GetSourceMapLocationMappings() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.sourceMapLocationMappings
}

// Java method: CompilerOptions.setProcessCommonJSModules(boolean processCommonJSModules)
func (this_ *CompilerOptions) SetProcessCommonJSModules(processCommonJSModules bool) {
	_ = processCommonJSModules
	this_.processCommonJSModules = processCommonJSModules
}

// Java method: CompilerOptions.getProcessCommonJSModules()
func (this_ *CompilerOptions) GetProcessCommonJSModules() bool {
	return this_.processCommonJSModules
}

// Java method: CompilerOptions.setEs6ModuleTranspilation(com.google.javascript.jscomp.CompilerOptions.Es6ModuleTranspilation value)
func (this_ *CompilerOptions) SetEs6ModuleTranspilation(value CompilerOptions_Es6ModuleTranspilation) {
	_ = value
	this_.es6ModuleTranspilation = value
}

// Java method: CompilerOptions.getEs6ModuleTranspilation()
func (this_ *CompilerOptions) GetEs6ModuleTranspilation() CompilerOptions_Es6ModuleTranspilation {
	return this_.es6ModuleTranspilation
}

// Java method: CompilerOptions.setCommonJSModulePathPrefix(java.lang.String commonJSModulePathPrefix)
func (this_ *CompilerOptions) SetCommonJSModulePathPrefix(commonJSModulePathPrefix string) {
	_ = commonJSModulePathPrefix
	this_.SetModuleRoots(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfStr1(commonJSModulePathPrefix)))
}

// Java method: CompilerOptions.setModuleRoots(java.util.List<java.lang.String> moduleRoots)
func (this_ *CompilerOptions) SetModuleRoots(moduleRoots j2gort.Ext_java_util_List) {
	_ = moduleRoots
	this_.moduleRoots = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](moduleRoots))
}

// Java method: CompilerOptions.getModuleRoots()
func (this_ *CompilerOptions) GetModuleRoots() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.moduleRoots
}

// Java method: CompilerOptions.setRewritePolyfills(boolean rewritePolyfills)
func (this_ *CompilerOptions) SetRewritePolyfills(rewritePolyfills bool) {
	_ = rewritePolyfills
	this_.rewritePolyfills = rewritePolyfills
}

// Java method: CompilerOptions.getRewritePolyfills()
func (this_ *CompilerOptions) GetRewritePolyfills() bool {
	return this_.rewritePolyfills
}

// Java method: CompilerOptions.setIsolatePolyfills(boolean isolatePolyfills)
func (this_ *CompilerOptions) SetIsolatePolyfills(isolatePolyfills bool) {
	_ = isolatePolyfills
	this_.isolatePolyfills = isolatePolyfills
	if this_.isolatePolyfills {
		this_.SetDefineToBooleanLiteral("$jscomp.ISOLATE_POLYFILLS", isolatePolyfills)
	}
}

// Java method: CompilerOptions.getIsolatePolyfills()
func (this_ *CompilerOptions) GetIsolatePolyfills() bool {
	return this_.isolatePolyfills
}

// Java method: CompilerOptions.setInjectPolyfillsNewerThan(com.google.javascript.jscomp.CompilerOptions.LanguageMode injectPolyfillsNewerThan)
func (this_ *CompilerOptions) SetInjectPolyfillsNewerThan(injectPolyfillsNewerThan CompilerOptions_LanguageMode) {
	_ = injectPolyfillsNewerThan
	this_.injectPolyfillsNewerThan = injectPolyfillsNewerThan
}

// Java method: CompilerOptions.getInjectPolyfillsNewerThan()
func (this_ *CompilerOptions) GetInjectPolyfillsNewerThan() CompilerOptions_LanguageMode {
	return this_.injectPolyfillsNewerThan
}

// Java method: CompilerOptions.setInstrumentAsyncContext(boolean instrumentAsyncContext)
func (this_ *CompilerOptions) SetInstrumentAsyncContext(instrumentAsyncContext bool) {
	_ = instrumentAsyncContext
	this_.instrumentAsyncContext = instrumentAsyncContext
	this_.SetDefineToBooleanLiteral("$jscomp.INSTRUMENT_ASYNC_CONTEXT", instrumentAsyncContext)
}

// Java method: CompilerOptions.getInstrumentAsyncContext()
func (this_ *CompilerOptions) GetInstrumentAsyncContext() bool {
	return this_.instrumentAsyncContext
}

// Java method: CompilerOptions.setEs6SubclassTranspilation(com.google.javascript.jscomp.CompilerOptions.Es6SubclassTranspilation es6SubclassTranspilation)
func (this_ *CompilerOptions) SetEs6SubclassTranspilation(es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) {
	_ = es6SubclassTranspilation
	this_.es6SubclassTranspilation = es6SubclassTranspilation
}

// Java method: CompilerOptions.getEs6SubclassTranspilation()
func (this_ *CompilerOptions) GetEs6SubclassTranspilation() CompilerOptions_Es6SubclassTranspilation {
	return this_.es6SubclassTranspilation
}

// Java method: CompilerOptions.setForceLibraryInjection(java.lang.Iterable<java.lang.String> libraries)
func (this_ *CompilerOptions) SetForceLibraryInjection(libraries j2gort.Ext_java_lang_Iterable) {
	_ = libraries
	this_.forceLibraryInjection = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(libraries)
}

// Java method: CompilerOptions.getForceLibraryInjectionList()
func (this_ *CompilerOptions) GetForceLibraryInjectionList() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.forceLibraryInjection
}

// Java method: CompilerOptions.setPreventLibraryInjection(boolean preventLibraryInjection)
func (this_ *CompilerOptions) SetPreventLibraryInjection(preventLibraryInjection bool) {
	_ = preventLibraryInjection
	if preventLibraryInjection {
		this_.runtimeLibraryMode = StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_NO_OP()
	} else {
		this_.runtimeLibraryMode = StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_INJECT()
	}
}

// Java method: CompilerOptions.setRuntimeLibraryMode(com.google.javascript.jscomp.js.RuntimeJsLibManager.RuntimeLibraryMode runtimeLibraryMode)
func (this_ *CompilerOptions) SetRuntimeLibraryMode(runtimeLibraryMode RuntimeJsLibManager_RuntimeLibraryMode) {
	_ = runtimeLibraryMode
	this_.runtimeLibraryMode = runtimeLibraryMode
}

// Java method: CompilerOptions.getRuntimeLibraryMode()
func (this_ *CompilerOptions) GetRuntimeLibraryMode() RuntimeJsLibManager_RuntimeLibraryMode {
	return this_.runtimeLibraryMode
}

// Java method: CompilerOptions.setUnusedImportsToRemove(com.google.common.collect.@org.jspecify.annotations.Nullable ImmutableSet<java.lang.String> unusedImportsToRemove)
func (this_ *CompilerOptions) SetUnusedImportsToRemove(unusedImportsToRemove j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = unusedImportsToRemove
	this_.unusedImportsToRemove = unusedImportsToRemove
}

// Java method: CompilerOptions.getUnusedImportsToRemove()
func (this_ *CompilerOptions) GetUnusedImportsToRemove() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.unusedImportsToRemove
}

// Java method: CompilerOptions.setInstrumentForCoverageOption(com.google.javascript.jscomp.CompilerOptions.InstrumentOption instrumentForCoverageOption)
func (this_ *CompilerOptions) SetInstrumentForCoverageOption(instrumentForCoverageOption CompilerOptions_InstrumentOption) {
	_ = instrumentForCoverageOption
	this_.instrumentForCoverageOption = j2gort.CheckNotNullEnum(instrumentForCoverageOption)
}

// Java method: CompilerOptions.getInstrumentForCoverageOption()
func (this_ *CompilerOptions) GetInstrumentForCoverageOption() CompilerOptions_InstrumentOption {
	return this_.instrumentForCoverageOption
}

// Java method: CompilerOptions.setProductionInstrumentationArrayName(java.lang.String productionInstrumentationArrayName)
func (this_ *CompilerOptions) SetProductionInstrumentationArrayName(productionInstrumentationArrayName string) {
	_ = productionInstrumentationArrayName
	this_.productionInstrumentationArrayName = j2gort.CheckNotNullString(productionInstrumentationArrayName)
}

// Java method: CompilerOptions.getProductionInstrumentationArrayName()
func (this_ *CompilerOptions) GetProductionInstrumentationArrayName() string {
	return this_.productionInstrumentationArrayName
}

// Java method: CompilerOptions.getConformanceConfigs()
func (this_ *CompilerOptions) GetConformanceConfigs() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(this_.conformanceConfigs)
}

// Java method: CompilerOptions.setConformanceConfig(com.google.javascript.jscomp.ConformanceConfig conformanceConfig)
func (this_ *CompilerOptions) SetConformanceConfig(conformanceConfig *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig) {
	_ = conformanceConfig
	this_.SetConformanceConfigs(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(conformanceConfig)))
}

// Java method: CompilerOptions.setConformanceConfigs(java.util.List<com.google.javascript.jscomp.ConformanceConfig> configs)
func (this_ *CompilerOptions) SetConformanceConfigs(configs j2gort.Ext_java_util_List) {
	_ = configs
	this_.conformanceConfigs = j2gort.ImmListCopyToSlice[*j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig](j2gort.AsIface[j2gort.Ext_java_util_Collection](configs))
}

// Java method: CompilerOptions.clearConformanceConfigs()
func (this_ *CompilerOptions) ClearConformanceConfigs() {
	this_.conformanceConfigs = []*j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig{}
}

// Java method: CompilerOptions.shouldEmitUseStrict()
func (this_ *CompilerOptions) ShouldEmitUseStrict() bool {
	return j2gort.UnboxBoolean(j2gort.CastPtr[j2gort.BoxBoolean](this_.emitUseStrict.Or_1_Ext_com_google_common_base_Optional(this_.languageOutIsDefaultStrict).Or_1_any(j2gort.ValueOfBoolean(this_.languageIn.IsDefaultStrict()))))
}

// Java method: CompilerOptions.setEmitUseStrict(boolean emitUseStrict)
func (this_ *CompilerOptions) SetEmitUseStrict(emitUseStrict bool) *CompilerOptions {
	_ = emitUseStrict
	this_.emitUseStrict = j2gort.Ext_com_google_common_base_Optional_Of(j2gort.ValueOfBoolean(emitUseStrict))
	return this_
}

// Java method: CompilerOptions.getModuleResolutionMode()
func (this_ *CompilerOptions) GetModuleResolutionMode() ModuleLoader_ResolutionMode {
	return this_.moduleResolutionMode
}

// Java method: CompilerOptions.setModuleResolutionMode(com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode moduleResolutionMode)
func (this_ *CompilerOptions) SetModuleResolutionMode(moduleResolutionMode ModuleLoader_ResolutionMode) {
	_ = moduleResolutionMode
	this_.moduleResolutionMode = moduleResolutionMode
}

// Java method: CompilerOptions.getBrowserResolverPrefixReplacements()
func (this_ *CompilerOptions) GetBrowserResolverPrefixReplacements() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.browserResolverPrefixReplacements
}

// Java method: CompilerOptions.setBrowserResolverPrefixReplacements(com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> browserResolverPrefixReplacements)
func (this_ *CompilerOptions) SetBrowserResolverPrefixReplacements(browserResolverPrefixReplacements j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = browserResolverPrefixReplacements
	this_.browserResolverPrefixReplacements = browserResolverPrefixReplacements
}

// Java method: CompilerOptions.setPathEscaper(com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper pathEscaper)
func (this_ *CompilerOptions) SetPathEscaper(pathEscaper *ModuleLoader_PathEscaper) {
	_ = pathEscaper
	this_.pathEscaper = pathEscaper
}

// Java method: CompilerOptions.getPathEscaper()
func (this_ *CompilerOptions) GetPathEscaper() *ModuleLoader_PathEscaper {
	return this_.pathEscaper
}

// Java method: CompilerOptions.getPackageJsonEntryNames()
func (this_ *CompilerOptions) GetPackageJsonEntryNames() j2gort.Ext_java_util_List {
	return this_.packageJsonEntryNames
}

// Java method: CompilerOptions.setPackageJsonEntryNames(java.util.List<java.lang.String> names)
func (this_ *CompilerOptions) SetPackageJsonEntryNames(names j2gort.Ext_java_util_List) {
	_ = names
	this_.packageJsonEntryNames = names
}

// Java method: CompilerOptions.setUseSizeHeuristicToStopOptimizationLoop(boolean mayStopEarly)
func (this_ *CompilerOptions) SetUseSizeHeuristicToStopOptimizationLoop(mayStopEarly bool) {
	_ = mayStopEarly
	this_.useSizeHeuristicToStopOptimizationLoop = mayStopEarly
}

// Java method: CompilerOptions.shouldUseSizeHeuristicToStopOptimizationLoop()
func (this_ *CompilerOptions) ShouldUseSizeHeuristicToStopOptimizationLoop() bool {
	return this_.useSizeHeuristicToStopOptimizationLoop
}

// Java method: CompilerOptions.setMaxOptimizationLoopIterations(int maxIterations)
func (this_ *CompilerOptions) SetMaxOptimizationLoopIterations(maxIterations int32) {
	_ = maxIterations
	this_.optimizationLoopMaxIterations = maxIterations
}

// Java method: CompilerOptions.getMaxOptimizationLoopIterations()
func (this_ *CompilerOptions) GetMaxOptimizationLoopIterations() int32 {
	return this_.optimizationLoopMaxIterations
}

// Java method: CompilerOptions.getChunkOutputType()
func (this_ *CompilerOptions) GetChunkOutputType() CompilerOptions_ChunkOutputType {
	return this_.chunkOutputType
}

// Java method: CompilerOptions.setChunkOutputType(com.google.javascript.jscomp.CompilerOptions.ChunkOutputType chunkOutputType)
func (this_ *CompilerOptions) SetChunkOutputType(chunkOutputType CompilerOptions_ChunkOutputType) {
	_ = chunkOutputType
	this_.chunkOutputType = chunkOutputType
}

// Java method: CompilerOptions.setStrictMessageReplacement(boolean strictMessageReplacement)
func (this_ *CompilerOptions) SetStrictMessageReplacement(strictMessageReplacement bool) {
	_ = strictMessageReplacement
	this_.strictMessageReplacement = strictMessageReplacement
}

// Java method: CompilerOptions.getStrictMessageReplacement()
func (this_ *CompilerOptions) GetStrictMessageReplacement() bool {
	return this_.strictMessageReplacement
}

// Java method: CompilerOptions.toString()
func (this_ *CompilerOptions) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).OmitNullValues().Add_2_string_any("aliasStringsMode", j2gort.BoxNullableVal(this_.GetAliasStringsMode())).Add_2_string_bool("ambiguateProperties", this_.ambiguateProperties).Add_2_string_bool("angularPass", this_.angularPass).Add_2_string_bool("assumeGettersArePure", this_.assumeGettersArePure).Add_2_string_bool("assumePropertiesAreStaticallyAnalyzable", this_.assumePropertiesAreStaticallyAnalyzable).Add_2_string_bool("assumeStrictThis", this_.AssumeStrictThis()).Add_2_string_any("browserResolverPrefixReplacements", this_.browserResolverPrefixReplacements).Add_2_string_bool("checkDeterminism", this_.GetCheckDeterminism()).Add_2_string_bool("checkSuspiciousCode", this_.checkSuspiciousCode).Add_2_string_bool("checkSymbols", this_.checkSymbols).Add_2_string_bool("checkTypes", this_.checkTypes).Add_2_string_bool("checksOnly", this_.checksOnly).Add_2_string_any("chunksToPrintAfterEachPassRegexList", this_.chunksToPrintAfterEachPassRegexList).Add_2_string_bool("closurePass", this_.closurePass).Add_2_string_bool("coalesceVariableNames", this_.coalesceVariableNames).Add_2_string_any("codingConvention", this_.GetCodingConvention()).Add_2_string_bool("collapseAnonymousFunctions", this_.collapseAnonymousFunctions).Add_2_string_bool("collapseObjectLiterals", this_.collapseObjectLiterals).Add_2_string_any("collapseProperties", j2gort.BoxNullableVal(this_.collapsePropertiesLevel)).Add_2_string_bool("collapseVariableDeclarations", this_.collapseVariableDeclarations).Add_2_string_bool("colorizeErrorOutput", this_.ShouldColorizeErrorOutput()).Add_2_string_bool("computeFunctionSideEffects", this_.computeFunctionSideEffects).Add_2_string_any("conformanceConfigs", this_.GetConformanceConfigs()).Add_2_string_any("conformanceRemoveRegexFromPath", this_.conformanceRemoveRegexFromPath).Add_2_string_any("conformanceReportingMode", j2gort.BoxNullableVal(this_.conformanceReportingMode)).Add_2_string_bool("continueAfterErrors", this_.CanContinueAfterErrors()).Add_2_string_bool("convertToDottedProperties", this_.convertToDottedProperties).Add_2_string_bool("crossChunkCodeMotion", this_.crossChunkCodeMotion).Add_2_string_bool("crossChunkCodeMotionNoStubMethods", this_.crossChunkCodeMotionNoStubMethods).Add_2_string_bool("crossChunkMethodMotion", this_.crossChunkMethodMotion).Add_2_string_any("cssRenamingMap", this_.cssRenamingMap).Add_2_string_any("cssRenamingSkiplist", this_.cssRenamingSkiplist).Add_2_string_any("customPasses", this_.customPasses).Add_2_string_bool("deadAssignmentElimination", this_.deadAssignmentElimination).Add_2_string_any("debugLogDirectory", this_.debugLogDirectory).Add_2_string_any("defineReplacements", this_.GetDefineReplacements()).Add_2_string_any("dependencyOptions", j2gort.PtrToAny(this_.GetDependencyOptions())).Add_2_string_any("devMode", j2gort.BoxNullableVal(this_.devMode)).Add_2_string_bool("devirtualizeMethods", this_.devirtualizeMethods).Add_2_string_bool("disambiguateProperties", this_.disambiguateProperties).Add_2_string_any("emitUseStrict", this_.emitUseStrict).Add_2_string_bool("enableModuleRewriting", this_.enableModuleRewriting).Add_2_string_any("environment", j2gort.BoxNullableVal(this_.GetEnvironment())).Add_2_string_any("errorFormat", j2gort.PtrToAny(this_.errorFormat)).Add_2_string_any("errorHandler", this_.errorHandler).Add_2_string_any("es6ModuleTranspilation", j2gort.BoxNullableVal(this_.es6ModuleTranspilation)).Add_2_string_any("es6SubclassTranspilation", j2gort.BoxNullableVal(this_.es6SubclassTranspilation)).Add_2_string_bool("exportLocalPropertyDefinitions", this_.exportLocalPropertyDefinitions).Add_2_string_bool("exportTestFunctions", this_.exportTestFunctions).Add_2_string_any("externExportsPath", j2gort.BoxString(this_.externExportsPath)).Add_2_string_any("extraAnnotationNames", this_.extraAnnotationNames).Add_2_string_any("extractPrototypeMemberDeclarations", j2gort.BoxNullableVal(this_.extractPrototypeMemberDeclarations)).Add_2_string_any("filesToPrintAfterEachPassRegexList", this_.filesToPrintAfterEachPassRegexList).Add_2_string_bool("flowSensitiveInlineVariables", this_.flowSensitiveInlineVariables).Add_2_string_bool("foldConstants", this_.foldConstants).Add_2_string_any("forceLibraryInjection", this_.forceLibraryInjection).Add_2_string_bool("gatherCssNames", this_.gatherCssNames).Add_2_string_bool("generateExports", this_.generateExports).Add_2_string_bool("generatePseudoNames", this_.generatePseudoNames).Add_2_string_bool("generateTypedExterns", this_.ShouldGenerateTypedExterns()).Add_2_string_any("idGenerators", this_.idGenerators).Add_2_string_any("idGeneratorsMapSerialized", j2gort.BoxString(this_.idGeneratorsMapSerialized)).Add_2_string_any("incrementalCheckMode", j2gort.BoxNullableVal(this_.incrementalCheckMode)).Add_2_string_bool("inferConsts", this_.inferConsts).Add_2_string_bool("inferTypes", this_.inferTypes).Add_2_string_bool("inlineConstantVars", this_.inlineConstantVars).Add_2_string_any("inlineFunctionsLevel", j2gort.BoxNullableVal(this_.inlineFunctionsLevel)).Add_2_string_bool("inlineGetters", this_.inlineGetters).Add_2_string_bool("inlineLocalVariables", this_.inlineLocalVariables).Add_2_string_bool("inlineProperties", this_.inlineProperties).Add_2_string_bool("inlineVariables", this_.inlineVariables).Add_2_string_any("inputDelimiter", j2gort.BoxString(this_.inputDelimiter)).Add_2_string_any("inputPropertyMap", j2gort.PtrToAny(this_.inputPropertyMap)).Add_2_string_any("inputSourceMaps", this_.inputSourceMaps).Add_2_string_any("inputVariableMap", j2gort.PtrToAny(this_.inputVariableMap)).Add_2_string_bool("instrumentForCoverageOnly", this_.instrumentForCoverageOnly).Add_2_string_any("instrumentForCoverageOption", j2gort.BoxString(j2gort.CheckNotNullEnum(this_.instrumentForCoverageOption).String())).Add_2_string_bool("isolatePolyfills", this_.isolatePolyfills).Add_2_string_bool("j2clMinifierEnabled", this_.j2clMinifierEnabled).Add_2_string_any("j2clMinifierPruningManifest", j2gort.BoxString(this_.j2clMinifierPruningManifest)).Add_2_string_any("j2clPassMode", j2gort.BoxNullableVal(this_.j2clPassMode)).Add_2_string_bool("labelRenaming", this_.labelRenaming).Add_2_string_any("languageIn", j2gort.BoxNullableVal(this_.GetLanguageIn())).Add_2_string_any("languageOutIsDefaultStrict", this_.languageOutIsDefaultStrict).Add_2_string_bool("lineBreak", this_.lineBreak).Add_2_string_int32("lineLengthThreshold", this_.lineLengthThreshold).Add_2_string_any("locale", j2gort.BoxString(this_.locale)).Add_2_string_bool("markAsCompiled", this_.markAsCompiled).Add_2_string_int32("maxFunctionSizeAfterInlining", this_.maxFunctionSizeAfterInlining).Add_2_string_any("messageBundle", this_.messageBundle).Add_2_string_any("moduleRoots", this_.moduleRoots).Add_2_string_any("nameGenerator", func() any {
		if this_.nameGenerator == nil || this_.nameGenerator.IsJavaNil_() {
			return nil
		}
		return j2gort.BoxString(j2gort.CheckNotNullPtr(j2gort.GetObjectClass(this_.nameGenerator)).GetName())
	}()).Add_2_string_int32("numParallelThreads", this_.numParallelThreads).Add_2_string_bool("optimizeCalls", this_.optimizeCalls).Add_2_string_bool("optimizeESClassConstructors", this_.optimizeESClassConstructors).Add_2_string_any("optimizeLocalAccessForGlobalSymbolNamespace", j2gort.BoxNullableVal(this_.optimizeLocalAccessForGlobalSymbolNamespace)).Add_2_string_any("outputCharset", this_.outputCharset).Add_2_string_any("outputFeatureSet", this_.outputFeatureSet).Add_2_string_any("outputJs", j2gort.BoxNullableVal(this_.outputJs)).Add_2_string_bool("outputJsStringUsage", this_.outputJsStringUsage).Add_2_string_bool("parentChunkCanSeeSymbolsDeclaredInChildren", this_.parentChunkCanSeeSymbolsDeclaredInChildren).Add_2_string_any("parseJsDocDocumentation", j2gort.BoxNullableVal(this_.IsParseJsDocDocumentation())).Add_2_string_any("pathEscaper", j2gort.PtrToAny(this_.pathEscaper)).Add_2_string_bool("polymerPass", this_.polymerPass).Add_2_string_bool("preferSingleQuotes", this_.preferSingleQuotes).Add_2_string_bool("preferStableNames", this_.preferStableNames).Add_2_string_bool("preserveDetailedSourceInfo", this_.PreservesDetailedSourceInfo()).Add_2_string_bool("preserveGoogProvidesAndRequires", this_.preserveClosurePrimitives).Add_2_string_bool("preserveNonJSDocComments", this_.GetPreserveNonJSDocComments()).Add_2_string_bool("preserveTypeAnnotations", this_.preserveTypeAnnotations).Add_2_string_bool("prettyPrint", this_.prettyPrint).Add_2_string_bool("printConfig", this_.printConfig).Add_2_string_bool("printInputDelimiter", this_.printInputDelimiter).Add_2_string_bool("printSourceAfterEachPass", this_.printSourceAfterEachPass).Add_2_string_bool("processCommonJSModules", this_.processCommonJSModules).Add_2_string_any("productionInstrumentationArrayName", j2gort.BoxString(this_.productionInstrumentationArrayName)).Add_2_string_any("propertiesThatMustDisambiguate", this_.propertiesThatMustDisambiguate).Add_2_string_any("propertyRenaming", j2gort.BoxNullableVal(this_.propertyRenaming)).Add_2_string_bool("propertyRenamingOnlyCompilationMode", this_.propertyRenamingOnlyCompilationMode).Add_2_string_bool("protectHiddenSideEffects", this_.protectHiddenSideEffects).Add_2_string_any("qnameUsesToPrintAfterEachPassRegexList", this_.qnameUsesToPrintAfterEachPassList).Add_2_string_bool("quoteKeywordProperties", this_.quoteKeywordProperties).Add_2_string_bool("removeAbstractMethods", this_.removeAbstractMethods).Add_2_string_bool("removeClosureAsserts", this_.removeClosureAsserts).Add_2_string_bool("removeJ2clAsserts", this_.removeJ2clAsserts).Add_2_string_bool("removeUnusedClassProperties", this_.removeUnusedClassProperties).Add_2_string_bool("removeUnusedLocalVars", this_.removeUnusedLocalVars).Add_2_string_bool("removeUnusedPrototypeProperties", this_.removeUnusedPrototypeProperties).Add_2_string_bool("removeUnusedVars", this_.removeUnusedVars).Add_2_string_any("renamePrefix", j2gort.BoxString(this_.renamePrefix)).Add_2_string_any("renamePrefixNamespace", j2gort.BoxString(this_.renamePrefixNamespace)).Add_2_string_bool("renamePrefixNamespaceAssumeCrossChunkNames", this_.renamePrefixNamespaceAssumeCrossChunkNames).Add_2_string_bool("replaceIdGenerators", this_.replaceIdGenerators).Add_2_string_bool("replaceMessagesWithChromeI18n", this_.replaceMessagesWithChromeI18n).Add_2_string_any("replaceStringsFunctionDescriptions", this_.replaceStringsFunctionDescriptions).Add_2_string_any("replaceStringsPlaceholderToken", j2gort.BoxString(this_.replaceStringsPlaceholderToken)).Add_2_string_bool("reserveRawExports", this_.reserveRawExports).Add_2_string_bool("rewriteFunctionExpressions", this_.rewriteFunctionExpressions).Add_2_string_bool("rewriteGlobalDeclarationsForTryCatchWrapping", this_.rewriteGlobalDeclarationsForTryCatchWrapping).Add_2_string_bool("rewriteModulesBeforeTypechecking", this_.rewriteModulesBeforeTypechecking).Add_2_string_bool("rewritePolyfills", this_.rewritePolyfills).Add_2_string_any("runtimeLibraryMode", j2gort.BoxNullableVal(this_.runtimeLibraryMode)).Add_2_string_bool("skipNonTranspilationPasses", this_.skipNonTranspilationPasses).Add_2_string_bool("smartNameRemoval", this_.smartNameRemoval).Add_2_string_any("sourceMapDetailLevel", j2gort.PtrToAny(this_.sourceMapDetailLevel)).Add_2_string_any("sourceMapFormat", j2gort.PtrToAny(this_.sourceMapFormat)).Add_2_string_any("sourceMapLocationMappings", this_.sourceMapLocationMappings).Add_2_string_any("sourceMapOutputPath", j2gort.BoxString(this_.sourceMapOutputPath)).Add_2_string_bool("strictMessageReplacement", this_.strictMessageReplacement).Add_2_string_any("stripNamePrefixes", this_.stripNamePrefixes).Add_2_string_any("stripNameSuffixes", this_.stripNameSuffixes).Add_2_string_any("stripTypes", this_.stripTypes).Add_2_string_int32("summaryDetailLevel", this_.summaryDetailLevel).Add_2_string_any("syntheticBlockEndMarker", j2gort.BoxString(this_.syntheticBlockEndMarker)).Add_2_string_any("syntheticBlockStartMarker", j2gort.BoxString(this_.syntheticBlockStartMarker)).Add_2_string_any("tcProjectId", j2gort.BoxString(this_.tcProjectId)).Add_2_string_any("tracer", j2gort.BoxNullableVal(this_.tracer)).Add_2_string_bool("trustedStrings", this_.trustedStrings).Add_2_string_any("tweakProcessing", j2gort.BoxNullableVal(this_.GetTweakProcessing())).Add_2_string_any("unusedImportsToRemove", this_.unusedImportsToRemove).Add_2_string_bool("useTypesForLocalOptimization", this_.useTypesForLocalOptimization).Add_2_string_any("variableRenaming", j2gort.BoxNullableVal(this_.variableRenaming)).Add_2_string_any("warningsGuard", this_.GetWarningsGuard()).Add_2_string_bool("wrapGoogModulesForWhitespaceOnly", this_.wrapGoogModulesForWhitespaceOnly).String()
}

// Java method: CompilerOptions.expectStrictModeInput()
func (this_ *CompilerOptions) ExpectStrictModeInput() bool {
	return j2gort.UnboxBoolean(j2gort.CastPtr[j2gort.BoxBoolean](this_.isStrictModeInput.Or_1_any(j2gort.ValueOfBoolean(this_.GetLanguageIn().IsDefaultStrict()))))
}

// Java method: CompilerOptions.setStrictModeInput(boolean isStrictModeInput)
func (this_ *CompilerOptions) SetStrictModeInput(isStrictModeInput bool) *CompilerOptions {
	_ = isStrictModeInput
	this_.isStrictModeInput = j2gort.Ext_com_google_common_base_Optional_Of(j2gort.ValueOfBoolean(isStrictModeInput))
	return this_
}

// Java method: CompilerOptions.getPropertyReservedNamingFirstChars()
func (this_ *CompilerOptions) GetPropertyReservedNamingFirstChars() j2gort.Ext_com_google_common_collect_ImmutableSet {
	if this_.polymerPass {
		return StaticGet_CompilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS()
	} else if this_.angularPass {
		return StaticGet_CompilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS()
	}
	return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
}

// Java method: CompilerOptions.getPropertyReservedNamingNonFirstChars()
func (this_ *CompilerOptions) GetPropertyReservedNamingNonFirstChars() j2gort.Ext_com_google_common_collect_ImmutableSet {
	if this_.polymerPass {
		return StaticGet_CompilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS()
	}
	return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
}

// Java method: CompilerOptions.shouldOptimize()
func (this_ *CompilerOptions) ShouldOptimize() bool {
	return (((!(this_.skipNonTranspilationPasses) && !(this_.checksOnly)) && !(this_.ShouldGenerateTypedExterns())) && !(this_.instrumentForCoverageOnly))
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.Reach
type CompilerOptions_Reach uint8

const (
	enum_CompilerOptions_Reach_ALL        CompilerOptions_Reach = 1
	enum_CompilerOptions_Reach_LOCAL_ONLY CompilerOptions_Reach = 2
	enum_CompilerOptions_Reach_NONE       CompilerOptions_Reach = 3
)

var enumBaseTable_CompilerOptions_Reach = [...]j2gort.EnumBase{
	{Name_: "ALL", Ordinal_: 0},
	{Name_: "LOCAL_ONLY", Ordinal_: 1},
	{Name_: "NONE", Ordinal_: 2},
}
var enumValues_CompilerOptions_Reach = []CompilerOptions_Reach{1, 2, 3}

func StaticGet_CompilerOptions_Reach_ALL() CompilerOptions_Reach {
	return enum_CompilerOptions_Reach_ALL
}
func StaticGet_CompilerOptions_Reach_LOCAL_ONLY() CompilerOptions_Reach {
	return enum_CompilerOptions_Reach_LOCAL_ONLY
}
func StaticGet_CompilerOptions_Reach_NONE() CompilerOptions_Reach {
	return enum_CompilerOptions_Reach_NONE
}

func CompilerOptions_Reach_Values() []CompilerOptions_Reach {
	return j2gort.CloneSlice(enumValues_CompilerOptions_Reach)
}

func CompilerOptions_Reach_ValueOf(name string) CompilerOptions_Reach {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_Reach {
		if enumBaseTable_CompilerOptions_Reach[i_].Name_ == name {
			return CompilerOptions_Reach(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.Reach." + name))
}

// Java method: Reach.isOn()
func (this_ CompilerOptions_Reach) IsOn() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ != StaticGet_CompilerOptions_Reach_NONE())
}

// Java method: Reach.includesGlobals()
func (this_ CompilerOptions_Reach) IncludesGlobals() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ == StaticGet_CompilerOptions_Reach_ALL())
}

func (this_ CompilerOptions_Reach) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_Reach) HasIdentityEquals_()                                 {}
func (CompilerOptions_Reach) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_Reach) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_Reach) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_Reach) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$Reach", "Reach"
}
func (this_ CompilerOptions_Reach) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_Reach[this_-1]
}
func (this_ CompilerOptions_Reach) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_Reach) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_Reach) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_Reach[this_-1].Name_
}
func (this_ CompilerOptions_Reach) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_Reach[this_-1].Name_
}
func (this_ CompilerOptions_Reach) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_Reach)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_Reach) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_Reach[this_-1].HashCode()
}
func (this_ CompilerOptions_Reach) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_Reach[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_Reach) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_Reach[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_Reach(1), CompilerOptions_Reach(2), CompilerOptions_Reach(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.Reach", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$Reach", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.Reach", "com.google.javascript.jscomp.CompilerOptions$Reach", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_Reach](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.PropertyCollapseLevel
type CompilerOptions_PropertyCollapseLevel uint8

const (
	enum_CompilerOptions_PropertyCollapseLevel_ALL           CompilerOptions_PropertyCollapseLevel = 1
	enum_CompilerOptions_PropertyCollapseLevel_NONE          CompilerOptions_PropertyCollapseLevel = 2
	enum_CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT CompilerOptions_PropertyCollapseLevel = 3
)

var enumBaseTable_CompilerOptions_PropertyCollapseLevel = [...]j2gort.EnumBase{
	{Name_: "ALL", Ordinal_: 0},
	{Name_: "NONE", Ordinal_: 1},
	{Name_: "MODULE_EXPORT", Ordinal_: 2},
}
var enumValues_CompilerOptions_PropertyCollapseLevel = []CompilerOptions_PropertyCollapseLevel{1, 2, 3}

func StaticGet_CompilerOptions_PropertyCollapseLevel_ALL() CompilerOptions_PropertyCollapseLevel {
	return enum_CompilerOptions_PropertyCollapseLevel_ALL
}
func StaticGet_CompilerOptions_PropertyCollapseLevel_NONE() CompilerOptions_PropertyCollapseLevel {
	return enum_CompilerOptions_PropertyCollapseLevel_NONE
}
func StaticGet_CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT() CompilerOptions_PropertyCollapseLevel {
	return enum_CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT
}

func CompilerOptions_PropertyCollapseLevel_Values() []CompilerOptions_PropertyCollapseLevel {
	return j2gort.CloneSlice(enumValues_CompilerOptions_PropertyCollapseLevel)
}

func CompilerOptions_PropertyCollapseLevel_ValueOf(name string) CompilerOptions_PropertyCollapseLevel {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_PropertyCollapseLevel {
		if enumBaseTable_CompilerOptions_PropertyCollapseLevel[i_].Name_ == name {
			return CompilerOptions_PropertyCollapseLevel(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.PropertyCollapseLevel." + name))
}

func (this_ CompilerOptions_PropertyCollapseLevel) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_PropertyCollapseLevel) HasIdentityEquals_()                                 {}
func (CompilerOptions_PropertyCollapseLevel) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_PropertyCollapseLevel) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_PropertyCollapseLevel) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_PropertyCollapseLevel) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$PropertyCollapseLevel", "PropertyCollapseLevel"
}
func (this_ CompilerOptions_PropertyCollapseLevel) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_PropertyCollapseLevel[this_-1]
}
func (this_ CompilerOptions_PropertyCollapseLevel) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_PropertyCollapseLevel) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_PropertyCollapseLevel) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_PropertyCollapseLevel[this_-1].Name_
}
func (this_ CompilerOptions_PropertyCollapseLevel) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_PropertyCollapseLevel[this_-1].Name_
}
func (this_ CompilerOptions_PropertyCollapseLevel) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_PropertyCollapseLevel)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_PropertyCollapseLevel) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_PropertyCollapseLevel[this_-1].HashCode()
}
func (this_ CompilerOptions_PropertyCollapseLevel) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_PropertyCollapseLevel[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_PropertyCollapseLevel) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_PropertyCollapseLevel[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_PropertyCollapseLevel(1), CompilerOptions_PropertyCollapseLevel(2), CompilerOptions_PropertyCollapseLevel(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.PropertyCollapseLevel", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$PropertyCollapseLevel", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.PropertyCollapseLevel", "com.google.javascript.jscomp.CompilerOptions$PropertyCollapseLevel", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_PropertyCollapseLevel](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.BrowserFeaturesetYear
type CompilerOptions_BrowserFeaturesetYear struct {
	j2gort.EnumBase
	year       int32
	featureSet *pkg_parser.FeatureSet
}

func (this_ *CompilerOptions_BrowserFeaturesetYear) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *CompilerOptions_BrowserFeaturesetYear) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*CompilerOptions_BrowserFeaturesetYear) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*CompilerOptions_BrowserFeaturesetYear) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *CompilerOptions_BrowserFeaturesetYear) initFields_CompilerOptions_BrowserFeaturesetYear() {
	_ = this_
}

func initEnumCompilerOptions_BrowserFeaturesetYear(this_ *CompilerOptions_BrowserFeaturesetYear, enumName_ string, enumOrd_ int32, year int32, featureSet *pkg_parser.FeatureSet) *CompilerOptions_BrowserFeaturesetYear {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = year
	_ = featureSet
	this_.initFields_CompilerOptions_BrowserFeaturesetYear()
	this_.year = year
	this_.featureSet = featureSet
	return this_
}

var enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2012 *CompilerOptions_BrowserFeaturesetYear

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2012() *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2012
}

var enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2018 *CompilerOptions_BrowserFeaturesetYear

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2018() *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2018
}

var enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2019 *CompilerOptions_BrowserFeaturesetYear

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2019() *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2019
}

var enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2020 *CompilerOptions_BrowserFeaturesetYear

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2020() *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2020
}

var enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2021 *CompilerOptions_BrowserFeaturesetYear

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2021() *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2021
}

var enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2022 *CompilerOptions_BrowserFeaturesetYear

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2022() *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2022
}

var enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2023 *CompilerOptions_BrowserFeaturesetYear

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2023() *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2023
}

var enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2024 *CompilerOptions_BrowserFeaturesetYear

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2024() *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2024
}

var enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2025 *CompilerOptions_BrowserFeaturesetYear

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2025() *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2025
}

var enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2026 *CompilerOptions_BrowserFeaturesetYear

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2026() *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2026
}

var static_CompilerOptions_BrowserFeaturesetYear_YEAR_MAP j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_MAP() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return static_CompilerOptions_BrowserFeaturesetYear_YEAR_MAP
}
func StaticSet_CompilerOptions_BrowserFeaturesetYear_YEAR_MAP(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	static_CompilerOptions_BrowserFeaturesetYear_YEAR_MAP = v
}

var clinitDone_CompilerOptions_BrowserFeaturesetYear uint32
var clinitState_CompilerOptions_BrowserFeaturesetYear j2gort.ClinitState

func clinit_CompilerOptions_BrowserFeaturesetYear() {
	if atomic.LoadUint32(&clinitDone_CompilerOptions_BrowserFeaturesetYear) != 0 {
		return
	}
	clinitSlow_CompilerOptions_BrowserFeaturesetYear()
}

func clinitSlow_CompilerOptions_BrowserFeaturesetYear() {
	if !j2gort.EnterClinit(&clinitState_CompilerOptions_BrowserFeaturesetYear, &clinitDone_CompilerOptions_BrowserFeaturesetYear) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CompilerOptions_BrowserFeaturesetYear, &clinitDone_CompilerOptions_BrowserFeaturesetYear)
	enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2012 = &CompilerOptions_BrowserFeaturesetYear{}
	initEnumCompilerOptions_BrowserFeaturesetYear(enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2012, "YEAR_2012", int32(0), int32(2012), StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT().ToFeatureSet())
	enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2018 = &CompilerOptions_BrowserFeaturesetYear{}
	initEnumCompilerOptions_BrowserFeaturesetYear(enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2018, "YEAR_2018", int32(1), int32(2018), StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_2016().ToFeatureSet())
	enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2019 = &CompilerOptions_BrowserFeaturesetYear{}
	initEnumCompilerOptions_BrowserFeaturesetYear(enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2019, "YEAR_2019", int32(2), int32(2019), pkg_parser.StaticGet_FeatureSet_BROWSER_2019())
	enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2020 = &CompilerOptions_BrowserFeaturesetYear{}
	initEnumCompilerOptions_BrowserFeaturesetYear(enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2020, "YEAR_2020", int32(3), int32(2020), pkg_parser.StaticGet_FeatureSet_BROWSER_2020())
	enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2021 = &CompilerOptions_BrowserFeaturesetYear{}
	initEnumCompilerOptions_BrowserFeaturesetYear(enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2021, "YEAR_2021", int32(4), int32(2021), pkg_parser.StaticGet_FeatureSet_BROWSER_2021())
	enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2022 = &CompilerOptions_BrowserFeaturesetYear{}
	initEnumCompilerOptions_BrowserFeaturesetYear(enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2022, "YEAR_2022", int32(5), int32(2022), pkg_parser.StaticGet_FeatureSet_BROWSER_2022())
	enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2023 = &CompilerOptions_BrowserFeaturesetYear{}
	initEnumCompilerOptions_BrowserFeaturesetYear(enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2023, "YEAR_2023", int32(6), int32(2023), pkg_parser.StaticGet_FeatureSet_BROWSER_2023())
	enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2024 = &CompilerOptions_BrowserFeaturesetYear{}
	initEnumCompilerOptions_BrowserFeaturesetYear(enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2024, "YEAR_2024", int32(7), int32(2024), pkg_parser.StaticGet_FeatureSet_BROWSER_2024())
	enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2025 = &CompilerOptions_BrowserFeaturesetYear{}
	initEnumCompilerOptions_BrowserFeaturesetYear(enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2025, "YEAR_2025", int32(8), int32(2025), pkg_parser.StaticGet_FeatureSet_BROWSER_2025())
	enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2026 = &CompilerOptions_BrowserFeaturesetYear{}
	initEnumCompilerOptions_BrowserFeaturesetYear(enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2026, "YEAR_2026", int32(9), int32(2026), pkg_parser.StaticGet_FeatureSet_BROWSER_2026())
	enumValues_CompilerOptions_BrowserFeaturesetYear = []*CompilerOptions_BrowserFeaturesetYear{enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2012, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2018, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2019, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2020, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2021, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2022, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2023, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2024, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2025, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2026}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.BrowserFeaturesetYear", []any{enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2012, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2018, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2019, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2020, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2021, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2022, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2023, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2024, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2025, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2026})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$BrowserFeaturesetYear", []any{enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2012, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2018, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2019, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2020, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2021, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2022, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2023, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2024, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2025, enum_CompilerOptions_BrowserFeaturesetYear_YEAR_2026})
	static_CompilerOptions_BrowserFeaturesetYear_YEAR_MAP = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_20(j2gort.ValueOfInt(int32(2012)), j2gort.PtrToAny(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2012()), j2gort.ValueOfInt(int32(2018)), j2gort.PtrToAny(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2018()), j2gort.ValueOfInt(int32(2019)), j2gort.PtrToAny(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2019()), j2gort.ValueOfInt(int32(2020)), j2gort.PtrToAny(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2020()), j2gort.ValueOfInt(int32(2021)), j2gort.PtrToAny(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2021()), j2gort.ValueOfInt(int32(2022)), j2gort.PtrToAny(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2022()), j2gort.ValueOfInt(int32(2023)), j2gort.PtrToAny(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2023()), j2gort.ValueOfInt(int32(2024)), j2gort.PtrToAny(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2024()), j2gort.ValueOfInt(int32(2025)), j2gort.PtrToAny(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2025()), j2gort.ValueOfInt(int32(2026)), j2gort.PtrToAny(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_2026()))
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.CompilerOptions.BrowserFeaturesetYear", "com.google.javascript.jscomp.CompilerOptions$BrowserFeaturesetYear", func() {
		clinit_CompilerOptions_BrowserFeaturesetYear()
	})
}

var enumValues_CompilerOptions_BrowserFeaturesetYear []*CompilerOptions_BrowserFeaturesetYear

func CompilerOptions_BrowserFeaturesetYear_Values() []*CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	return j2gort.CloneSlice(enumValues_CompilerOptions_BrowserFeaturesetYear)
}

func CompilerOptions_BrowserFeaturesetYear_ValueOf(name string) *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	for _, c := range enumValues_CompilerOptions_BrowserFeaturesetYear {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.BrowserFeaturesetYear." + name))
}

// Java static method: BrowserFeaturesetYear.from(int year)
func CompilerOptions_BrowserFeaturesetYear_From(year int32) *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	_ = year
	if !(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_MAP().ContainsKey(j2gort.ValueOfInt(year))) {
		j2gort.CheckStateFmt(false, "Illegal browser_featureset_year=%s. We support values 2012, or 2018..2026 only", j2gort.ValueOfInt(year))
	}
	return j2gort.CastPtr[CompilerOptions_BrowserFeaturesetYear](StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_MAP().Get(j2gort.ValueOfInt(year)))
}

// Java method: BrowserFeaturesetYear.setDependentValuesFromYear(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *CompilerOptions_BrowserFeaturesetYear) SetDependentValuesFromYear(options *CompilerOptions) {
	_ = options
	options.SetOutputFeatureSet(this_.featureSet)
	{
		r_6 := options
		r_6.languageOutIsDefaultStrict = j2gort.Ext_com_google_common_base_Optional_Of(j2gort.ValueOfBoolean(true))
	}
	options.SetDefineToNumberLiteral("goog.FEATURESET_YEAR", this_.year)
	options.SetDefineToBooleanLiteral("$jscomp.ASSUME_ES5", (this_.year > int32(2012)))
	options.SetDefineToBooleanLiteral("$jscomp.ASSUME_ES6", (this_.year >= int32(2018)))
	options.SetDefineToBooleanLiteral("$jscomp.ASSUME_ES2020", (this_.year >= int32(2021)))
}

// Java method: BrowserFeaturesetYear.getFeatureSet()
func (this_ *CompilerOptions_BrowserFeaturesetYear) GetFeatureSet() *pkg_parser.FeatureSet {
	return this_.featureSet
}

// Java method: BrowserFeaturesetYear.getYear()
func (this_ *CompilerOptions_BrowserFeaturesetYear) GetYear() int32 {
	return this_.year
}

// Java static method: BrowserFeaturesetYear.minimumRequiredFor(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor(feature *pkg_parser.FeatureSet_Feature) *CompilerOptions_BrowserFeaturesetYear {
	clinit_CompilerOptions_BrowserFeaturesetYear()
	_ = feature
	{
		cur_7 := j2gort.NewMapKeyCursor(StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_MAP())
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			year := javalang.D2I(j2gort.UnboxToDouble(rawElem_8))
			_ = year
			var yearObject *CompilerOptions_BrowserFeaturesetYear = j2gort.CastPtr[CompilerOptions_BrowserFeaturesetYear](StaticGet_CompilerOptions_BrowserFeaturesetYear_YEAR_MAP().Get(j2gort.ValueOfInt(year)))
			_ = yearObject
			if yearObject.GetFeatureSet().Contains_1_FeatureSet_Feature(feature) {
				return yearObject
			}
		}
	}
	return nil
}

func (this_ *CompilerOptions_BrowserFeaturesetYear) IsJavaNil_() bool { return this_ == nil }
func (*CompilerOptions_BrowserFeaturesetYear) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$BrowserFeaturesetYear", "BrowserFeaturesetYear"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.BrowserFeaturesetYear", "com.google.javascript.jscomp.CompilerOptions$BrowserFeaturesetYear", func(v_ any) bool {
		return j2gort.InstanceOf[*CompilerOptions_BrowserFeaturesetYear](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.IncrementalCheckMode
type CompilerOptions_IncrementalCheckMode uint8

const (
	enum_CompilerOptions_IncrementalCheckMode_OFF          CompilerOptions_IncrementalCheckMode = 1
	enum_CompilerOptions_IncrementalCheckMode_GENERATE_IJS CompilerOptions_IncrementalCheckMode = 2
)

var enumBaseTable_CompilerOptions_IncrementalCheckMode = [...]j2gort.EnumBase{
	{Name_: "OFF", Ordinal_: 0},
	{Name_: "GENERATE_IJS", Ordinal_: 1},
}
var enumValues_CompilerOptions_IncrementalCheckMode = []CompilerOptions_IncrementalCheckMode{1, 2}

func StaticGet_CompilerOptions_IncrementalCheckMode_OFF() CompilerOptions_IncrementalCheckMode {
	return enum_CompilerOptions_IncrementalCheckMode_OFF
}
func StaticGet_CompilerOptions_IncrementalCheckMode_GENERATE_IJS() CompilerOptions_IncrementalCheckMode {
	return enum_CompilerOptions_IncrementalCheckMode_GENERATE_IJS
}

func CompilerOptions_IncrementalCheckMode_Values() []CompilerOptions_IncrementalCheckMode {
	return j2gort.CloneSlice(enumValues_CompilerOptions_IncrementalCheckMode)
}

func CompilerOptions_IncrementalCheckMode_ValueOf(name string) CompilerOptions_IncrementalCheckMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_IncrementalCheckMode {
		if enumBaseTable_CompilerOptions_IncrementalCheckMode[i_].Name_ == name {
			return CompilerOptions_IncrementalCheckMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.IncrementalCheckMode." + name))
}

func (this_ CompilerOptions_IncrementalCheckMode) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_IncrementalCheckMode) HasIdentityEquals_()                                 {}
func (CompilerOptions_IncrementalCheckMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_IncrementalCheckMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_IncrementalCheckMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_IncrementalCheckMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$IncrementalCheckMode", "IncrementalCheckMode"
}
func (this_ CompilerOptions_IncrementalCheckMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_IncrementalCheckMode[this_-1]
}
func (this_ CompilerOptions_IncrementalCheckMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_IncrementalCheckMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_IncrementalCheckMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_IncrementalCheckMode[this_-1].Name_
}
func (this_ CompilerOptions_IncrementalCheckMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_IncrementalCheckMode[this_-1].Name_
}
func (this_ CompilerOptions_IncrementalCheckMode) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_IncrementalCheckMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_IncrementalCheckMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_IncrementalCheckMode[this_-1].HashCode()
}
func (this_ CompilerOptions_IncrementalCheckMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_IncrementalCheckMode[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_IncrementalCheckMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_IncrementalCheckMode[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_IncrementalCheckMode(1), CompilerOptions_IncrementalCheckMode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.IncrementalCheckMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$IncrementalCheckMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.IncrementalCheckMode", "com.google.javascript.jscomp.CompilerOptions$IncrementalCheckMode", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_IncrementalCheckMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.ExtractPrototypeMemberDeclarationsMode
type CompilerOptions_ExtractPrototypeMemberDeclarationsMode uint8

const (
	enum_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF             CompilerOptions_ExtractPrototypeMemberDeclarationsMode = 1
	enum_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP CompilerOptions_ExtractPrototypeMemberDeclarationsMode = 2
	enum_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP  CompilerOptions_ExtractPrototypeMemberDeclarationsMode = 3
	enum_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_IIFE        CompilerOptions_ExtractPrototypeMemberDeclarationsMode = 4
)

var enumBaseTable_CompilerOptions_ExtractPrototypeMemberDeclarationsMode = [...]j2gort.EnumBase{
	{Name_: "OFF", Ordinal_: 0},
	{Name_: "USE_GLOBAL_TEMP", Ordinal_: 1},
	{Name_: "USE_CHUNK_TEMP", Ordinal_: 2},
	{Name_: "USE_IIFE", Ordinal_: 3},
}
var enumValues_CompilerOptions_ExtractPrototypeMemberDeclarationsMode = []CompilerOptions_ExtractPrototypeMemberDeclarationsMode{1, 2, 3, 4}

func StaticGet_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF() CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	return enum_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF
}
func StaticGet_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP() CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	return enum_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP
}
func StaticGet_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP() CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	return enum_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP
}
func StaticGet_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_IIFE() CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	return enum_CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_IIFE
}

func CompilerOptions_ExtractPrototypeMemberDeclarationsMode_Values() []CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	return j2gort.CloneSlice(enumValues_CompilerOptions_ExtractPrototypeMemberDeclarationsMode)
}

func CompilerOptions_ExtractPrototypeMemberDeclarationsMode_ValueOf(name string) CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
		if enumBaseTable_CompilerOptions_ExtractPrototypeMemberDeclarationsMode[i_].Name_ == name {
			return CompilerOptions_ExtractPrototypeMemberDeclarationsMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.ExtractPrototypeMemberDeclarationsMode." + name))
}

func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) IsJavaNil_() bool {
	return this_ == 0
}
func (CompilerOptions_ExtractPrototypeMemberDeclarationsMode) HasIdentityEquals_() {}
func (CompilerOptions_ExtractPrototypeMemberDeclarationsMode) Implements_j2gort_Ext_java_lang_Enum() {
}
func (CompilerOptions_ExtractPrototypeMemberDeclarationsMode) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_ExtractPrototypeMemberDeclarationsMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$ExtractPrototypeMemberDeclarationsMode", "ExtractPrototypeMemberDeclarationsMode"
}
func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_ExtractPrototypeMemberDeclarationsMode[this_-1]
}
func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_ExtractPrototypeMemberDeclarationsMode[this_-1].Name_
}
func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_ExtractPrototypeMemberDeclarationsMode[this_-1].Name_
}
func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_ExtractPrototypeMemberDeclarationsMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_ExtractPrototypeMemberDeclarationsMode[this_-1].HashCode()
}
func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_ExtractPrototypeMemberDeclarationsMode[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_ExtractPrototypeMemberDeclarationsMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_ExtractPrototypeMemberDeclarationsMode[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_ExtractPrototypeMemberDeclarationsMode(1), CompilerOptions_ExtractPrototypeMemberDeclarationsMode(2), CompilerOptions_ExtractPrototypeMemberDeclarationsMode(3), CompilerOptions_ExtractPrototypeMemberDeclarationsMode(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.ExtractPrototypeMemberDeclarationsMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$ExtractPrototypeMemberDeclarationsMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.ExtractPrototypeMemberDeclarationsMode", "com.google.javascript.jscomp.CompilerOptions$ExtractPrototypeMemberDeclarationsMode", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_ExtractPrototypeMemberDeclarationsMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.OutputJs
type CompilerOptions_OutputJs uint8

const (
	enum_CompilerOptions_OutputJs_NONE     CompilerOptions_OutputJs = 1
	enum_CompilerOptions_OutputJs_SENTINEL CompilerOptions_OutputJs = 2
	enum_CompilerOptions_OutputJs_NORMAL   CompilerOptions_OutputJs = 3
)

var enumBaseTable_CompilerOptions_OutputJs = [...]j2gort.EnumBase{
	{Name_: "NONE", Ordinal_: 0},
	{Name_: "SENTINEL", Ordinal_: 1},
	{Name_: "NORMAL", Ordinal_: 2},
}
var enumValues_CompilerOptions_OutputJs = []CompilerOptions_OutputJs{1, 2, 3}

func StaticGet_CompilerOptions_OutputJs_NONE() CompilerOptions_OutputJs {
	return enum_CompilerOptions_OutputJs_NONE
}
func StaticGet_CompilerOptions_OutputJs_SENTINEL() CompilerOptions_OutputJs {
	return enum_CompilerOptions_OutputJs_SENTINEL
}
func StaticGet_CompilerOptions_OutputJs_NORMAL() CompilerOptions_OutputJs {
	return enum_CompilerOptions_OutputJs_NORMAL
}

func CompilerOptions_OutputJs_Values() []CompilerOptions_OutputJs {
	return j2gort.CloneSlice(enumValues_CompilerOptions_OutputJs)
}

func CompilerOptions_OutputJs_ValueOf(name string) CompilerOptions_OutputJs {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_OutputJs {
		if enumBaseTable_CompilerOptions_OutputJs[i_].Name_ == name {
			return CompilerOptions_OutputJs(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.OutputJs." + name))
}

func (this_ CompilerOptions_OutputJs) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_OutputJs) HasIdentityEquals_()                                 {}
func (CompilerOptions_OutputJs) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_OutputJs) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_OutputJs) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_OutputJs) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$OutputJs", "OutputJs"
}
func (this_ CompilerOptions_OutputJs) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_OutputJs[this_-1]
}
func (this_ CompilerOptions_OutputJs) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_OutputJs) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_OutputJs) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_OutputJs[this_-1].Name_
}
func (this_ CompilerOptions_OutputJs) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_OutputJs[this_-1].Name_
}
func (this_ CompilerOptions_OutputJs) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_OutputJs)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_OutputJs) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_OutputJs[this_-1].HashCode()
}
func (this_ CompilerOptions_OutputJs) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_OutputJs[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_OutputJs) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_OutputJs[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_OutputJs(1), CompilerOptions_OutputJs(2), CompilerOptions_OutputJs(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.OutputJs", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$OutputJs", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.OutputJs", "com.google.javascript.jscomp.CompilerOptions$OutputJs", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_OutputJs](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.ConformanceReportingMode
type CompilerOptions_ConformanceReportingMode uint8

const (
	enum_CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG  CompilerOptions_ConformanceReportingMode = 1
	enum_CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG CompilerOptions_ConformanceReportingMode = 2
)

var enumBaseTable_CompilerOptions_ConformanceReportingMode = [...]j2gort.EnumBase{
	{Name_: "IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG", Ordinal_: 0},
	{Name_: "RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG", Ordinal_: 1},
}
var enumValues_CompilerOptions_ConformanceReportingMode = []CompilerOptions_ConformanceReportingMode{1, 2}

func StaticGet_CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG() CompilerOptions_ConformanceReportingMode {
	return enum_CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
}
func StaticGet_CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG() CompilerOptions_ConformanceReportingMode {
	return enum_CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
}

func CompilerOptions_ConformanceReportingMode_Values() []CompilerOptions_ConformanceReportingMode {
	return j2gort.CloneSlice(enumValues_CompilerOptions_ConformanceReportingMode)
}

func CompilerOptions_ConformanceReportingMode_ValueOf(name string) CompilerOptions_ConformanceReportingMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_ConformanceReportingMode {
		if enumBaseTable_CompilerOptions_ConformanceReportingMode[i_].Name_ == name {
			return CompilerOptions_ConformanceReportingMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.ConformanceReportingMode." + name))
}

func (this_ CompilerOptions_ConformanceReportingMode) IsJavaNil_() bool                { return this_ == 0 }
func (CompilerOptions_ConformanceReportingMode) HasIdentityEquals_()                   {}
func (CompilerOptions_ConformanceReportingMode) Implements_j2gort_Ext_java_lang_Enum() {}
func (CompilerOptions_ConformanceReportingMode) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ CompilerOptions_ConformanceReportingMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_ConformanceReportingMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$ConformanceReportingMode", "ConformanceReportingMode"
}
func (this_ CompilerOptions_ConformanceReportingMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_ConformanceReportingMode[this_-1]
}
func (this_ CompilerOptions_ConformanceReportingMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_ConformanceReportingMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_ConformanceReportingMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_ConformanceReportingMode[this_-1].Name_
}
func (this_ CompilerOptions_ConformanceReportingMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_ConformanceReportingMode[this_-1].Name_
}
func (this_ CompilerOptions_ConformanceReportingMode) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_ConformanceReportingMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_ConformanceReportingMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_ConformanceReportingMode[this_-1].HashCode()
}
func (this_ CompilerOptions_ConformanceReportingMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_ConformanceReportingMode[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_ConformanceReportingMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_ConformanceReportingMode[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_ConformanceReportingMode(1), CompilerOptions_ConformanceReportingMode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.ConformanceReportingMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$ConformanceReportingMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.ConformanceReportingMode", "com.google.javascript.jscomp.CompilerOptions$ConformanceReportingMode", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_ConformanceReportingMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.ExperimentalOutputFeatureSet
type CompilerOptions_ExperimentalOutputFeatureSet uint8

const (
	enum_CompilerOptions_ExperimentalOutputFeatureSet_ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS CompilerOptions_ExperimentalOutputFeatureSet = 1
	enum_CompilerOptions_ExperimentalOutputFeatureSet_BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD             CompilerOptions_ExperimentalOutputFeatureSet = 2
)

var enumBaseTable_CompilerOptions_ExperimentalOutputFeatureSet = [...]j2gort.EnumBase{
	{Name_: "ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS", Ordinal_: 0},
	{Name_: "BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD", Ordinal_: 1},
}
var enumValues_CompilerOptions_ExperimentalOutputFeatureSet = []CompilerOptions_ExperimentalOutputFeatureSet{1, 2}

func StaticGet_CompilerOptions_ExperimentalOutputFeatureSet_ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS() CompilerOptions_ExperimentalOutputFeatureSet {
	return enum_CompilerOptions_ExperimentalOutputFeatureSet_ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS
}
func StaticGet_CompilerOptions_ExperimentalOutputFeatureSet_BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD() CompilerOptions_ExperimentalOutputFeatureSet {
	return enum_CompilerOptions_ExperimentalOutputFeatureSet_BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD
}

func CompilerOptions_ExperimentalOutputFeatureSet_Values() []CompilerOptions_ExperimentalOutputFeatureSet {
	return j2gort.CloneSlice(enumValues_CompilerOptions_ExperimentalOutputFeatureSet)
}

func CompilerOptions_ExperimentalOutputFeatureSet_ValueOf(name string) CompilerOptions_ExperimentalOutputFeatureSet {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_ExperimentalOutputFeatureSet {
		if enumBaseTable_CompilerOptions_ExperimentalOutputFeatureSet[i_].Name_ == name {
			return CompilerOptions_ExperimentalOutputFeatureSet(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.ExperimentalOutputFeatureSet." + name))
}

func (this_ CompilerOptions_ExperimentalOutputFeatureSet) IsJavaNil_() bool                { return this_ == 0 }
func (CompilerOptions_ExperimentalOutputFeatureSet) HasIdentityEquals_()                   {}
func (CompilerOptions_ExperimentalOutputFeatureSet) Implements_j2gort_Ext_java_lang_Enum() {}
func (CompilerOptions_ExperimentalOutputFeatureSet) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ CompilerOptions_ExperimentalOutputFeatureSet) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_ExperimentalOutputFeatureSet) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$ExperimentalOutputFeatureSet", "ExperimentalOutputFeatureSet"
}
func (this_ CompilerOptions_ExperimentalOutputFeatureSet) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_ExperimentalOutputFeatureSet[this_-1]
}
func (this_ CompilerOptions_ExperimentalOutputFeatureSet) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_ExperimentalOutputFeatureSet) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ CompilerOptions_ExperimentalOutputFeatureSet) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_ExperimentalOutputFeatureSet[this_-1].Name_
}
func (this_ CompilerOptions_ExperimentalOutputFeatureSet) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_ExperimentalOutputFeatureSet[this_-1].Name_
}
func (this_ CompilerOptions_ExperimentalOutputFeatureSet) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_ExperimentalOutputFeatureSet)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_ExperimentalOutputFeatureSet) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_ExperimentalOutputFeatureSet[this_-1].HashCode()
}
func (this_ CompilerOptions_ExperimentalOutputFeatureSet) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_ExperimentalOutputFeatureSet[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_ExperimentalOutputFeatureSet) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_ExperimentalOutputFeatureSet[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_ExperimentalOutputFeatureSet(1), CompilerOptions_ExperimentalOutputFeatureSet(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.ExperimentalOutputFeatureSet", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$ExperimentalOutputFeatureSet", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.ExperimentalOutputFeatureSet", "com.google.javascript.jscomp.CompilerOptions$ExperimentalOutputFeatureSet", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_ExperimentalOutputFeatureSet](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.Es6ModuleTranspilation
type CompilerOptions_Es6ModuleTranspilation uint8

const (
	enum_CompilerOptions_Es6ModuleTranspilation_NONE                      CompilerOptions_Es6ModuleTranspilation = 1
	enum_CompilerOptions_Es6ModuleTranspilation_RELATIVIZE_IMPORT_PATHS   CompilerOptions_Es6ModuleTranspilation = 2
	enum_CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES CompilerOptions_Es6ModuleTranspilation = 3
	enum_CompilerOptions_Es6ModuleTranspilation_COMPILE                   CompilerOptions_Es6ModuleTranspilation = 4
)

var enumBaseTable_CompilerOptions_Es6ModuleTranspilation = [...]j2gort.EnumBase{
	{Name_: "NONE", Ordinal_: 0},
	{Name_: "RELATIVIZE_IMPORT_PATHS", Ordinal_: 1},
	{Name_: "TO_COMMON_JS_LIKE_MODULES", Ordinal_: 2},
	{Name_: "COMPILE", Ordinal_: 3},
}
var enumValues_CompilerOptions_Es6ModuleTranspilation = []CompilerOptions_Es6ModuleTranspilation{1, 2, 3, 4}

func StaticGet_CompilerOptions_Es6ModuleTranspilation_NONE() CompilerOptions_Es6ModuleTranspilation {
	return enum_CompilerOptions_Es6ModuleTranspilation_NONE
}
func StaticGet_CompilerOptions_Es6ModuleTranspilation_RELATIVIZE_IMPORT_PATHS() CompilerOptions_Es6ModuleTranspilation {
	return enum_CompilerOptions_Es6ModuleTranspilation_RELATIVIZE_IMPORT_PATHS
}
func StaticGet_CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES() CompilerOptions_Es6ModuleTranspilation {
	return enum_CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES
}
func StaticGet_CompilerOptions_Es6ModuleTranspilation_COMPILE() CompilerOptions_Es6ModuleTranspilation {
	return enum_CompilerOptions_Es6ModuleTranspilation_COMPILE
}

func CompilerOptions_Es6ModuleTranspilation_Values() []CompilerOptions_Es6ModuleTranspilation {
	return j2gort.CloneSlice(enumValues_CompilerOptions_Es6ModuleTranspilation)
}

func CompilerOptions_Es6ModuleTranspilation_ValueOf(name string) CompilerOptions_Es6ModuleTranspilation {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_Es6ModuleTranspilation {
		if enumBaseTable_CompilerOptions_Es6ModuleTranspilation[i_].Name_ == name {
			return CompilerOptions_Es6ModuleTranspilation(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.Es6ModuleTranspilation." + name))
}

func (this_ CompilerOptions_Es6ModuleTranspilation) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_Es6ModuleTranspilation) HasIdentityEquals_()                                 {}
func (CompilerOptions_Es6ModuleTranspilation) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_Es6ModuleTranspilation) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_Es6ModuleTranspilation) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_Es6ModuleTranspilation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$Es6ModuleTranspilation", "Es6ModuleTranspilation"
}
func (this_ CompilerOptions_Es6ModuleTranspilation) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_Es6ModuleTranspilation[this_-1]
}
func (this_ CompilerOptions_Es6ModuleTranspilation) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_Es6ModuleTranspilation) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_Es6ModuleTranspilation) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_Es6ModuleTranspilation[this_-1].Name_
}
func (this_ CompilerOptions_Es6ModuleTranspilation) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_Es6ModuleTranspilation[this_-1].Name_
}
func (this_ CompilerOptions_Es6ModuleTranspilation) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_Es6ModuleTranspilation)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_Es6ModuleTranspilation) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_Es6ModuleTranspilation[this_-1].HashCode()
}
func (this_ CompilerOptions_Es6ModuleTranspilation) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_Es6ModuleTranspilation[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_Es6ModuleTranspilation) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_Es6ModuleTranspilation[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_Es6ModuleTranspilation(1), CompilerOptions_Es6ModuleTranspilation(2), CompilerOptions_Es6ModuleTranspilation(3), CompilerOptions_Es6ModuleTranspilation(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.Es6ModuleTranspilation", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$Es6ModuleTranspilation", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.Es6ModuleTranspilation", "com.google.javascript.jscomp.CompilerOptions$Es6ModuleTranspilation", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_Es6ModuleTranspilation](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.Es6SubclassTranspilation
type CompilerOptions_Es6SubclassTranspilation uint8

const (
	enum_CompilerOptions_Es6SubclassTranspilation_CONCISE_UNSAFE         CompilerOptions_Es6SubclassTranspilation = 1
	enum_CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT CompilerOptions_Es6SubclassTranspilation = 2
)

var enumBaseTable_CompilerOptions_Es6SubclassTranspilation = [...]j2gort.EnumBase{
	{Name_: "CONCISE_UNSAFE", Ordinal_: 0},
	{Name_: "SAFE_REFLECT_CONSTRUCT", Ordinal_: 1},
}
var enumValues_CompilerOptions_Es6SubclassTranspilation = []CompilerOptions_Es6SubclassTranspilation{1, 2}

func StaticGet_CompilerOptions_Es6SubclassTranspilation_CONCISE_UNSAFE() CompilerOptions_Es6SubclassTranspilation {
	return enum_CompilerOptions_Es6SubclassTranspilation_CONCISE_UNSAFE
}
func StaticGet_CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT() CompilerOptions_Es6SubclassTranspilation {
	return enum_CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT
}

func CompilerOptions_Es6SubclassTranspilation_Values() []CompilerOptions_Es6SubclassTranspilation {
	return j2gort.CloneSlice(enumValues_CompilerOptions_Es6SubclassTranspilation)
}

func CompilerOptions_Es6SubclassTranspilation_ValueOf(name string) CompilerOptions_Es6SubclassTranspilation {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_Es6SubclassTranspilation {
		if enumBaseTable_CompilerOptions_Es6SubclassTranspilation[i_].Name_ == name {
			return CompilerOptions_Es6SubclassTranspilation(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.Es6SubclassTranspilation." + name))
}

func (this_ CompilerOptions_Es6SubclassTranspilation) IsJavaNil_() bool                { return this_ == 0 }
func (CompilerOptions_Es6SubclassTranspilation) HasIdentityEquals_()                   {}
func (CompilerOptions_Es6SubclassTranspilation) Implements_j2gort_Ext_java_lang_Enum() {}
func (CompilerOptions_Es6SubclassTranspilation) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ CompilerOptions_Es6SubclassTranspilation) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_Es6SubclassTranspilation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$Es6SubclassTranspilation", "Es6SubclassTranspilation"
}
func (this_ CompilerOptions_Es6SubclassTranspilation) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_Es6SubclassTranspilation[this_-1]
}
func (this_ CompilerOptions_Es6SubclassTranspilation) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_Es6SubclassTranspilation) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_Es6SubclassTranspilation) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_Es6SubclassTranspilation[this_-1].Name_
}
func (this_ CompilerOptions_Es6SubclassTranspilation) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_Es6SubclassTranspilation[this_-1].Name_
}
func (this_ CompilerOptions_Es6SubclassTranspilation) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_Es6SubclassTranspilation)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_Es6SubclassTranspilation) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_Es6SubclassTranspilation[this_-1].HashCode()
}
func (this_ CompilerOptions_Es6SubclassTranspilation) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_Es6SubclassTranspilation[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_Es6SubclassTranspilation) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_Es6SubclassTranspilation[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_Es6SubclassTranspilation(1), CompilerOptions_Es6SubclassTranspilation(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.Es6SubclassTranspilation", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$Es6SubclassTranspilation", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.Es6SubclassTranspilation", "com.google.javascript.jscomp.CompilerOptions$Es6SubclassTranspilation", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_Es6SubclassTranspilation](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.InstrumentOption
type CompilerOptions_InstrumentOption uint8

const (
	enum_CompilerOptions_InstrumentOption_NONE        CompilerOptions_InstrumentOption = 1
	enum_CompilerOptions_InstrumentOption_LINE_ONLY   CompilerOptions_InstrumentOption = 2
	enum_CompilerOptions_InstrumentOption_BRANCH_ONLY CompilerOptions_InstrumentOption = 3
	enum_CompilerOptions_InstrumentOption_PRODUCTION  CompilerOptions_InstrumentOption = 4
)

var enumBaseTable_CompilerOptions_InstrumentOption = [...]j2gort.EnumBase{
	{Name_: "NONE", Ordinal_: 0},
	{Name_: "LINE_ONLY", Ordinal_: 1},
	{Name_: "BRANCH_ONLY", Ordinal_: 2},
	{Name_: "PRODUCTION", Ordinal_: 3},
}
var enumValues_CompilerOptions_InstrumentOption = []CompilerOptions_InstrumentOption{1, 2, 3, 4}

func StaticGet_CompilerOptions_InstrumentOption_NONE() CompilerOptions_InstrumentOption {
	return enum_CompilerOptions_InstrumentOption_NONE
}
func StaticGet_CompilerOptions_InstrumentOption_LINE_ONLY() CompilerOptions_InstrumentOption {
	return enum_CompilerOptions_InstrumentOption_LINE_ONLY
}
func StaticGet_CompilerOptions_InstrumentOption_BRANCH_ONLY() CompilerOptions_InstrumentOption {
	return enum_CompilerOptions_InstrumentOption_BRANCH_ONLY
}
func StaticGet_CompilerOptions_InstrumentOption_PRODUCTION() CompilerOptions_InstrumentOption {
	return enum_CompilerOptions_InstrumentOption_PRODUCTION
}

func CompilerOptions_InstrumentOption_Values() []CompilerOptions_InstrumentOption {
	return j2gort.CloneSlice(enumValues_CompilerOptions_InstrumentOption)
}

func CompilerOptions_InstrumentOption_ValueOf(name string) CompilerOptions_InstrumentOption {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_InstrumentOption {
		if enumBaseTable_CompilerOptions_InstrumentOption[i_].Name_ == name {
			return CompilerOptions_InstrumentOption(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.InstrumentOption." + name))
}

// Java static method: InstrumentOption.fromString(java.lang.String value)
func CompilerOptions_InstrumentOption_FromString(value string) CompilerOptions_InstrumentOption {
	_ = value
	if value == j2gort.NullString {
		return CompilerOptions_InstrumentOption(0)
	}
	switch j2gort.CheckNotNullString(value) {
	case "NONE":
		return StaticGet_CompilerOptions_InstrumentOption_NONE()
	case "LINE":
		return StaticGet_CompilerOptions_InstrumentOption_LINE_ONLY()
	case "BRANCH":
		return StaticGet_CompilerOptions_InstrumentOption_BRANCH_ONLY()
	case "PRODUCTION":
		return StaticGet_CompilerOptions_InstrumentOption_PRODUCTION()
	default:
		return CompilerOptions_InstrumentOption(0)
	}
	panic("unreachable switch expression")
}

func (this_ CompilerOptions_InstrumentOption) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_InstrumentOption) HasIdentityEquals_()                                 {}
func (CompilerOptions_InstrumentOption) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_InstrumentOption) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_InstrumentOption) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_InstrumentOption) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$InstrumentOption", "InstrumentOption"
}
func (this_ CompilerOptions_InstrumentOption) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_InstrumentOption[this_-1]
}
func (this_ CompilerOptions_InstrumentOption) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_InstrumentOption) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_InstrumentOption) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_InstrumentOption[this_-1].Name_
}
func (this_ CompilerOptions_InstrumentOption) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_InstrumentOption[this_-1].Name_
}
func (this_ CompilerOptions_InstrumentOption) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_InstrumentOption)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_InstrumentOption) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_InstrumentOption[this_-1].HashCode()
}
func (this_ CompilerOptions_InstrumentOption) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_InstrumentOption[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_InstrumentOption) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_InstrumentOption[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_InstrumentOption(1), CompilerOptions_InstrumentOption(2), CompilerOptions_InstrumentOption(3), CompilerOptions_InstrumentOption(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.InstrumentOption", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$InstrumentOption", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.InstrumentOption", "com.google.javascript.jscomp.CompilerOptions$InstrumentOption", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_InstrumentOption](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.ChunkOutputType
type CompilerOptions_ChunkOutputType uint8

const (
	enum_CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE CompilerOptions_ChunkOutputType = 1
	enum_CompilerOptions_ChunkOutputType_ES_MODULES       CompilerOptions_ChunkOutputType = 2
)

var enumBaseTable_CompilerOptions_ChunkOutputType = [...]j2gort.EnumBase{
	{Name_: "GLOBAL_NAMESPACE", Ordinal_: 0},
	{Name_: "ES_MODULES", Ordinal_: 1},
}
var enumValues_CompilerOptions_ChunkOutputType = []CompilerOptions_ChunkOutputType{1, 2}

func StaticGet_CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE() CompilerOptions_ChunkOutputType {
	return enum_CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE
}
func StaticGet_CompilerOptions_ChunkOutputType_ES_MODULES() CompilerOptions_ChunkOutputType {
	return enum_CompilerOptions_ChunkOutputType_ES_MODULES
}

func CompilerOptions_ChunkOutputType_Values() []CompilerOptions_ChunkOutputType {
	return j2gort.CloneSlice(enumValues_CompilerOptions_ChunkOutputType)
}

func CompilerOptions_ChunkOutputType_ValueOf(name string) CompilerOptions_ChunkOutputType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_ChunkOutputType {
		if enumBaseTable_CompilerOptions_ChunkOutputType[i_].Name_ == name {
			return CompilerOptions_ChunkOutputType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.ChunkOutputType." + name))
}

func (this_ CompilerOptions_ChunkOutputType) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_ChunkOutputType) HasIdentityEquals_()                                 {}
func (CompilerOptions_ChunkOutputType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_ChunkOutputType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_ChunkOutputType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_ChunkOutputType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$ChunkOutputType", "ChunkOutputType"
}
func (this_ CompilerOptions_ChunkOutputType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_ChunkOutputType[this_-1]
}
func (this_ CompilerOptions_ChunkOutputType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_ChunkOutputType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_ChunkOutputType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_ChunkOutputType[this_-1].Name_
}
func (this_ CompilerOptions_ChunkOutputType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_ChunkOutputType[this_-1].Name_
}
func (this_ CompilerOptions_ChunkOutputType) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_ChunkOutputType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_ChunkOutputType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_ChunkOutputType[this_-1].HashCode()
}
func (this_ CompilerOptions_ChunkOutputType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_ChunkOutputType[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_ChunkOutputType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_ChunkOutputType[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_ChunkOutputType(1), CompilerOptions_ChunkOutputType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.ChunkOutputType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$ChunkOutputType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.ChunkOutputType", "com.google.javascript.jscomp.CompilerOptions$ChunkOutputType", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_ChunkOutputType](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.OptimizeLocalAccess
type CompilerOptions_OptimizeLocalAccess uint8

const (
	enum_CompilerOptions_OptimizeLocalAccess_DISABLED                                              CompilerOptions_OptimizeLocalAccess = 1
	enum_CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY                                   CompilerOptions_OptimizeLocalAccess = 2
	enum_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS                                            CompilerOptions_OptimizeLocalAccess = 3
	enum_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS          CompilerOptions_OptimizeLocalAccess = 4
	enum_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED CompilerOptions_OptimizeLocalAccess = 5
	enum_CompilerOptions_OptimizeLocalAccess_DELETE_GLOBALS_ON_LAST_ACCESS                         CompilerOptions_OptimizeLocalAccess = 6
)

var enumBaseTable_CompilerOptions_OptimizeLocalAccess = [...]j2gort.EnumBase{
	{Name_: "DISABLED", Ordinal_: 0},
	{Name_: "DEFINING_CHUNK_ONLY", Ordinal_: 1},
	{Name_: "ALL_CHUNKS", Ordinal_: 2},
	{Name_: "ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS", Ordinal_: 3},
	{Name_: "ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED", Ordinal_: 4},
	{Name_: "DELETE_GLOBALS_ON_LAST_ACCESS", Ordinal_: 5},
}
var enumValues_CompilerOptions_OptimizeLocalAccess = []CompilerOptions_OptimizeLocalAccess{1, 2, 3, 4, 5, 6}

func StaticGet_CompilerOptions_OptimizeLocalAccess_DISABLED() CompilerOptions_OptimizeLocalAccess {
	return enum_CompilerOptions_OptimizeLocalAccess_DISABLED
}
func StaticGet_CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY() CompilerOptions_OptimizeLocalAccess {
	return enum_CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
}
func StaticGet_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS() CompilerOptions_OptimizeLocalAccess {
	return enum_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS
}
func StaticGet_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS() CompilerOptions_OptimizeLocalAccess {
	return enum_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
}
func StaticGet_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED() CompilerOptions_OptimizeLocalAccess {
	return enum_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
}
func StaticGet_CompilerOptions_OptimizeLocalAccess_DELETE_GLOBALS_ON_LAST_ACCESS() CompilerOptions_OptimizeLocalAccess {
	return enum_CompilerOptions_OptimizeLocalAccess_DELETE_GLOBALS_ON_LAST_ACCESS
}

func CompilerOptions_OptimizeLocalAccess_Values() []CompilerOptions_OptimizeLocalAccess {
	return j2gort.CloneSlice(enumValues_CompilerOptions_OptimizeLocalAccess)
}

func CompilerOptions_OptimizeLocalAccess_ValueOf(name string) CompilerOptions_OptimizeLocalAccess {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_OptimizeLocalAccess {
		if enumBaseTable_CompilerOptions_OptimizeLocalAccess[i_].Name_ == name {
			return CompilerOptions_OptimizeLocalAccess(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.OptimizeLocalAccess." + name))
}

// Java method: OptimizeLocalAccess.shouldOptimize()
func (this_ CompilerOptions_OptimizeLocalAccess) ShouldOptimize() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ != StaticGet_CompilerOptions_OptimizeLocalAccess_DISABLED())
}

// Java method: OptimizeLocalAccess.shouldOptimizeAllChunks()
func (this_ CompilerOptions_OptimizeLocalAccess) ShouldOptimizeAllChunks() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return ((this_).Ordinal() >= (StaticGet_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS()).Ordinal())
}

// Java method: OptimizeLocalAccess.shouldWrapReassignableSymbols()
func (this_ CompilerOptions_OptimizeLocalAccess) ShouldWrapReassignableSymbols() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return ((this_).Ordinal() >= (StaticGet_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS()).Ordinal())
}

// Java method: OptimizeLocalAccess.shouldPreserveLocalVarDeclarations()
func (this_ CompilerOptions_OptimizeLocalAccess) ShouldPreserveLocalVarDeclarations() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return ((this_).Ordinal() >= (StaticGet_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED()).Ordinal())
}

// Java method: OptimizeLocalAccess.shouldRewriteCompoundAndIncrementWritesInPlace()
func (this_ CompilerOptions_OptimizeLocalAccess) ShouldRewriteCompoundAndIncrementWritesInPlace() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return ((this_).Ordinal() >= (StaticGet_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED()).Ordinal())
}

// Java method: OptimizeLocalAccess.shouldKeepLocalAliasForWrappedReassignables()
func (this_ CompilerOptions_OptimizeLocalAccess) ShouldKeepLocalAliasForWrappedReassignables() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return ((this_).Ordinal() >= (StaticGet_CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED()).Ordinal())
}

// Java method: OptimizeLocalAccess.shouldDeleteGlobalsOnLastAccess()
func (this_ CompilerOptions_OptimizeLocalAccess) ShouldDeleteGlobalsOnLastAccess() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ == StaticGet_CompilerOptions_OptimizeLocalAccess_DELETE_GLOBALS_ON_LAST_ACCESS())
}

func (this_ CompilerOptions_OptimizeLocalAccess) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_OptimizeLocalAccess) HasIdentityEquals_()                                 {}
func (CompilerOptions_OptimizeLocalAccess) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_OptimizeLocalAccess) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_OptimizeLocalAccess) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_OptimizeLocalAccess) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$OptimizeLocalAccess", "OptimizeLocalAccess"
}
func (this_ CompilerOptions_OptimizeLocalAccess) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_OptimizeLocalAccess[this_-1]
}
func (this_ CompilerOptions_OptimizeLocalAccess) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_OptimizeLocalAccess) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_OptimizeLocalAccess) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_OptimizeLocalAccess[this_-1].Name_
}
func (this_ CompilerOptions_OptimizeLocalAccess) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_OptimizeLocalAccess[this_-1].Name_
}
func (this_ CompilerOptions_OptimizeLocalAccess) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_OptimizeLocalAccess)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_OptimizeLocalAccess) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_OptimizeLocalAccess[this_-1].HashCode()
}
func (this_ CompilerOptions_OptimizeLocalAccess) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_OptimizeLocalAccess[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_OptimizeLocalAccess) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_OptimizeLocalAccess[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_OptimizeLocalAccess(1), CompilerOptions_OptimizeLocalAccess(2), CompilerOptions_OptimizeLocalAccess(3), CompilerOptions_OptimizeLocalAccess(4), CompilerOptions_OptimizeLocalAccess(5), CompilerOptions_OptimizeLocalAccess(6)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.OptimizeLocalAccess", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$OptimizeLocalAccess", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.OptimizeLocalAccess", "com.google.javascript.jscomp.CompilerOptions$OptimizeLocalAccess", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_OptimizeLocalAccess](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.LanguageMode
type CompilerOptions_LanguageMode uint8

const (
	enum_CompilerOptions_LanguageMode_ECMASCRIPT3        CompilerOptions_LanguageMode = 1
	enum_CompilerOptions_LanguageMode_ECMASCRIPT5        CompilerOptions_LanguageMode = 2
	enum_CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT CompilerOptions_LanguageMode = 3
	enum_CompilerOptions_LanguageMode_ECMASCRIPT_2015    CompilerOptions_LanguageMode = 4
	enum_CompilerOptions_LanguageMode_ECMASCRIPT_2016    CompilerOptions_LanguageMode = 5
	enum_CompilerOptions_LanguageMode_ECMASCRIPT_2017    CompilerOptions_LanguageMode = 6
	enum_CompilerOptions_LanguageMode_ECMASCRIPT_2018    CompilerOptions_LanguageMode = 7
	enum_CompilerOptions_LanguageMode_ECMASCRIPT_2019    CompilerOptions_LanguageMode = 8
	enum_CompilerOptions_LanguageMode_ECMASCRIPT_2020    CompilerOptions_LanguageMode = 9
	enum_CompilerOptions_LanguageMode_ECMASCRIPT_2021    CompilerOptions_LanguageMode = 10
	enum_CompilerOptions_LanguageMode_ECMASCRIPT_2022    CompilerOptions_LanguageMode = 11
	enum_CompilerOptions_LanguageMode_ECMASCRIPT_NEXT    CompilerOptions_LanguageMode = 12
	enum_CompilerOptions_LanguageMode_STABLE             CompilerOptions_LanguageMode = 13
	enum_CompilerOptions_LanguageMode_NO_TRANSPILE       CompilerOptions_LanguageMode = 14
	enum_CompilerOptions_LanguageMode_UNSTABLE           CompilerOptions_LanguageMode = 15
	enum_CompilerOptions_LanguageMode_UNSUPPORTED        CompilerOptions_LanguageMode = 16
)

var enumBaseTable_CompilerOptions_LanguageMode = [...]j2gort.EnumBase{
	{Name_: "ECMASCRIPT3", Ordinal_: 0},
	{Name_: "ECMASCRIPT5", Ordinal_: 1},
	{Name_: "ECMASCRIPT5_STRICT", Ordinal_: 2},
	{Name_: "ECMASCRIPT_2015", Ordinal_: 3},
	{Name_: "ECMASCRIPT_2016", Ordinal_: 4},
	{Name_: "ECMASCRIPT_2017", Ordinal_: 5},
	{Name_: "ECMASCRIPT_2018", Ordinal_: 6},
	{Name_: "ECMASCRIPT_2019", Ordinal_: 7},
	{Name_: "ECMASCRIPT_2020", Ordinal_: 8},
	{Name_: "ECMASCRIPT_2021", Ordinal_: 9},
	{Name_: "ECMASCRIPT_2022", Ordinal_: 10},
	{Name_: "ECMASCRIPT_NEXT", Ordinal_: 11},
	{Name_: "STABLE", Ordinal_: 12},
	{Name_: "NO_TRANSPILE", Ordinal_: 13},
	{Name_: "UNSTABLE", Ordinal_: 14},
	{Name_: "UNSUPPORTED", Ordinal_: 15},
}
var enumValues_CompilerOptions_LanguageMode = []CompilerOptions_LanguageMode{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT3() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT3
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT5() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT5
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_2015() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT_2015
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_2016() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT_2016
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_2017() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT_2017
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_2018() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT_2018
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_2019() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT_2019
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_2020() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT_2020
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_2021() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT_2021
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_2022() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT_2022
}
func StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_NEXT() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
}
func StaticGet_CompilerOptions_LanguageMode_STABLE() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_STABLE
}
func StaticGet_CompilerOptions_LanguageMode_NO_TRANSPILE() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_NO_TRANSPILE
}
func StaticGet_CompilerOptions_LanguageMode_UNSTABLE() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_UNSTABLE
}
func StaticGet_CompilerOptions_LanguageMode_UNSUPPORTED() CompilerOptions_LanguageMode {
	return enum_CompilerOptions_LanguageMode_UNSUPPORTED
}

var static_CompilerOptions_LanguageMode_STABLE_IN CompilerOptions_LanguageMode

func StaticGet_CompilerOptions_LanguageMode_STABLE_IN() CompilerOptions_LanguageMode {
	clinit_CompilerOptions_LanguageMode()
	return static_CompilerOptions_LanguageMode_STABLE_IN
}
func StaticSet_CompilerOptions_LanguageMode_STABLE_IN(v CompilerOptions_LanguageMode) {
	clinit_CompilerOptions_LanguageMode()
	static_CompilerOptions_LanguageMode_STABLE_IN = v
}

var static_CompilerOptions_LanguageMode_STABLE_OUT CompilerOptions_LanguageMode

func StaticGet_CompilerOptions_LanguageMode_STABLE_OUT() CompilerOptions_LanguageMode {
	clinit_CompilerOptions_LanguageMode()
	return static_CompilerOptions_LanguageMode_STABLE_OUT
}
func StaticSet_CompilerOptions_LanguageMode_STABLE_OUT(v CompilerOptions_LanguageMode) {
	clinit_CompilerOptions_LanguageMode()
	static_CompilerOptions_LanguageMode_STABLE_OUT = v
}

var clinitDone_CompilerOptions_LanguageMode uint32
var clinitState_CompilerOptions_LanguageMode j2gort.ClinitState

func clinit_CompilerOptions_LanguageMode() {
	if atomic.LoadUint32(&clinitDone_CompilerOptions_LanguageMode) != 0 {
		return
	}
	clinitSlow_CompilerOptions_LanguageMode()
}

func clinitSlow_CompilerOptions_LanguageMode() {
	if !j2gort.EnterClinit(&clinitState_CompilerOptions_LanguageMode, &clinitDone_CompilerOptions_LanguageMode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CompilerOptions_LanguageMode, &clinitDone_CompilerOptions_LanguageMode)
	static_CompilerOptions_LanguageMode_STABLE_IN = StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_NEXT()
	static_CompilerOptions_LanguageMode_STABLE_OUT = StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT5()
}

func CompilerOptions_LanguageMode_Values() []CompilerOptions_LanguageMode {
	clinit_CompilerOptions_LanguageMode()
	return j2gort.CloneSlice(enumValues_CompilerOptions_LanguageMode)
}

func CompilerOptions_LanguageMode_ValueOf(name string) CompilerOptions_LanguageMode {
	clinit_CompilerOptions_LanguageMode()
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_LanguageMode {
		if enumBaseTable_CompilerOptions_LanguageMode[i_].Name_ == name {
			return CompilerOptions_LanguageMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.LanguageMode." + name))
}

// Java method: LanguageMode.isDefaultStrict()
func (this_ CompilerOptions_LanguageMode) IsDefaultStrict() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	switch (this_).Ordinal() {
	case int32(0) /* LanguageMode.ECMASCRIPT3 */, int32(1) /* LanguageMode.ECMASCRIPT5 */ :
		return false
	default:
		return true
	}
	panic("unreachable switch expression")
}

// Java static method: LanguageMode.validCommandLineNames()
func CompilerOptions_LanguageMode_ValidCommandLineNames() j2gort.Ext_java_util_List {
	clinit_CompilerOptions_LanguageMode()
	var names []string = []string{}
	_ = names
	{
		cur_9 := j2gort.NewTypedIterCursor[CompilerOptions_LanguageMode](j2gort.CheckNotNull(j2gort.Ext_java_util_EnumSet_AllOf(j2gort.ClassLiteral("com.google.javascript.jscomp.CompilerOptions.LanguageMode", "LanguageMode"))))
		for cur_9.HasNext() {
			mode := cur_9.Next()
			_ = mode
			if mode != StaticGet_CompilerOptions_LanguageMode_UNSUPPORTED() {
				names = append(names, (mode).Name())
				if javalang.StringStartsWith(j2gort.CheckNotNullString((mode).Name()), "ECMASCRIPT") {
					names = append(names, j2gort.JStrReplaceFirst((mode).Name(), "ECMASCRIPT", "ES"))
				}
			}
		}
	}
	names = append(names, "ECMASCRIPT6")
	names = append(names, "ES6")
	names = append(names, "ECMASCRIPT6_STRICT")
	names = append(names, "ES6_STRICT")
	return j2gort.WrapTypedArrayList(names)
}

// Java static method: LanguageMode.fromString(java.lang.String value)
func CompilerOptions_LanguageMode_FromString(value string) CompilerOptions_LanguageMode {
	clinit_CompilerOptions_LanguageMode()
	_ = value
	if value == j2gort.NullString {
		return CompilerOptions_LanguageMode(0)
	}
	var canonicalizedName string = j2gort.JStrReplaceFirst(j2gort.Ext_com_google_common_base_Ascii_ToUpperCase_1_string(javalang.StringTrim(j2gort.CheckNotNullString(value))), "^ES", "ECMASCRIPT")
	_ = canonicalizedName
	if (j2gort.CheckNotNullString(canonicalizedName) == "ECMASCRIPT6") || (j2gort.CheckNotNullString(canonicalizedName) == "ECMASCRIPT6_STRICT") {
		return StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_2015()
	}
	var didRet_10 bool
	var retVal_11 CompilerOptions_LanguageMode
	func() {
		defer func() {
			if rec_12 := recover(); rec_12 != nil {
				jex_13 := javalang.RecoverException(rec_12)
				_ = jex_13
				if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_13) {
					e := j2gort.AsIface[javalang.IIllegalArgumentException](jex_13)
					_ = e
					retVal_11 = CompilerOptions_LanguageMode(0)
					didRet_10 = true
					return
				} else {
					panic(rec_12)
				}
			}
		}()
		retVal_11 = CompilerOptions_LanguageMode_ValueOf(canonicalizedName)
		didRet_10 = true
		return
	}()
	if didRet_10 {
		return retVal_11
	}
	panic("unreachable")
}

// Java method: LanguageMode.toFeatureSet()
func (this_ CompilerOptions_LanguageMode) ToFeatureSet() *pkg_parser.FeatureSet {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	switch (this_).Ordinal() {
	case int32(0) /* LanguageMode.ECMASCRIPT3 */ :
		return pkg_parser.StaticGet_FeatureSet_ES3()
	case int32(1) /* LanguageMode.ECMASCRIPT5 */, int32(2) /* LanguageMode.ECMASCRIPT5_STRICT */ :
		return pkg_parser.StaticGet_FeatureSet_ES5()
	case int32(3) /* LanguageMode.ECMASCRIPT_2015 */ :
		return pkg_parser.StaticGet_FeatureSet_ES2015_MODULES()
	case int32(4) /* LanguageMode.ECMASCRIPT_2016 */ :
		return pkg_parser.StaticGet_FeatureSet_ES2016_MODULES()
	case int32(5) /* LanguageMode.ECMASCRIPT_2017 */ :
		return pkg_parser.StaticGet_FeatureSet_ES2017_MODULES()
	case int32(6) /* LanguageMode.ECMASCRIPT_2018 */ :
		return pkg_parser.StaticGet_FeatureSet_ES2018_MODULES()
	case int32(7) /* LanguageMode.ECMASCRIPT_2019 */ :
		return pkg_parser.StaticGet_FeatureSet_ES2019_MODULES()
	case int32(8) /* LanguageMode.ECMASCRIPT_2020 */ :
		return pkg_parser.StaticGet_FeatureSet_ES2020_MODULES()
	case int32(9) /* LanguageMode.ECMASCRIPT_2021 */ :
		return pkg_parser.StaticGet_FeatureSet_ES2021_MODULES()
	case int32(10) /* LanguageMode.ECMASCRIPT_2022 */ :
		return pkg_parser.StaticGet_FeatureSet_ES2022_MODULES()
	case int32(11) /* LanguageMode.ECMASCRIPT_NEXT */ :
		return pkg_parser.StaticGet_FeatureSet_ES_NEXT()
	case int32(13) /* LanguageMode.NO_TRANSPILE */, int32(14) /* LanguageMode.UNSTABLE */ :
		return pkg_parser.StaticGet_FeatureSet_ES_UNSTABLE()
	case int32(15) /* LanguageMode.UNSUPPORTED */ :
		return pkg_parser.StaticGet_FeatureSet_ES_UNSUPPORTED()
	case int32(12) /* LanguageMode.STABLE */ :
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "STABLE has different feature sets for language in and out. Use STABLE_IN or STABLE_OUT.")))
	}
	panic("unreachable switch expression")
}

func (this_ CompilerOptions_LanguageMode) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_LanguageMode) HasIdentityEquals_()                                 {}
func (CompilerOptions_LanguageMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_LanguageMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_LanguageMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_LanguageMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$LanguageMode", "LanguageMode"
}
func (this_ CompilerOptions_LanguageMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_LanguageMode[this_-1]
}
func (this_ CompilerOptions_LanguageMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_LanguageMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_LanguageMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_LanguageMode[this_-1].Name_
}
func (this_ CompilerOptions_LanguageMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_LanguageMode[this_-1].Name_
}
func (this_ CompilerOptions_LanguageMode) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_LanguageMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_LanguageMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_LanguageMode[this_-1].HashCode()
}
func (this_ CompilerOptions_LanguageMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_LanguageMode[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_LanguageMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_LanguageMode[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_LanguageMode(1), CompilerOptions_LanguageMode(2), CompilerOptions_LanguageMode(3), CompilerOptions_LanguageMode(4), CompilerOptions_LanguageMode(5), CompilerOptions_LanguageMode(6), CompilerOptions_LanguageMode(7), CompilerOptions_LanguageMode(8), CompilerOptions_LanguageMode(9), CompilerOptions_LanguageMode(10), CompilerOptions_LanguageMode(11), CompilerOptions_LanguageMode(12), CompilerOptions_LanguageMode(13), CompilerOptions_LanguageMode(14), CompilerOptions_LanguageMode(15), CompilerOptions_LanguageMode(16)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.LanguageMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$LanguageMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.LanguageMode", "com.google.javascript.jscomp.CompilerOptions$LanguageMode", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_LanguageMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.DevMode
type CompilerOptions_DevMode uint8

const (
	enum_CompilerOptions_DevMode_OFF           CompilerOptions_DevMode = 1
	enum_CompilerOptions_DevMode_START         CompilerOptions_DevMode = 2
	enum_CompilerOptions_DevMode_START_AND_END CompilerOptions_DevMode = 3
	enum_CompilerOptions_DevMode_EVERY_PASS    CompilerOptions_DevMode = 4
)

var enumBaseTable_CompilerOptions_DevMode = [...]j2gort.EnumBase{
	{Name_: "OFF", Ordinal_: 0},
	{Name_: "START", Ordinal_: 1},
	{Name_: "START_AND_END", Ordinal_: 2},
	{Name_: "EVERY_PASS", Ordinal_: 3},
}
var enumValues_CompilerOptions_DevMode = []CompilerOptions_DevMode{1, 2, 3, 4}

func StaticGet_CompilerOptions_DevMode_OFF() CompilerOptions_DevMode {
	return enum_CompilerOptions_DevMode_OFF
}
func StaticGet_CompilerOptions_DevMode_START() CompilerOptions_DevMode {
	return enum_CompilerOptions_DevMode_START
}
func StaticGet_CompilerOptions_DevMode_START_AND_END() CompilerOptions_DevMode {
	return enum_CompilerOptions_DevMode_START_AND_END
}
func StaticGet_CompilerOptions_DevMode_EVERY_PASS() CompilerOptions_DevMode {
	return enum_CompilerOptions_DevMode_EVERY_PASS
}

func CompilerOptions_DevMode_Values() []CompilerOptions_DevMode {
	return j2gort.CloneSlice(enumValues_CompilerOptions_DevMode)
}

func CompilerOptions_DevMode_ValueOf(name string) CompilerOptions_DevMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_DevMode {
		if enumBaseTable_CompilerOptions_DevMode[i_].Name_ == name {
			return CompilerOptions_DevMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.DevMode." + name))
}

func (this_ CompilerOptions_DevMode) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_DevMode) HasIdentityEquals_()                                 {}
func (CompilerOptions_DevMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_DevMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_DevMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_DevMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$DevMode", "DevMode"
}
func (this_ CompilerOptions_DevMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_DevMode[this_-1]
}
func (this_ CompilerOptions_DevMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_DevMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_DevMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_DevMode[this_-1].Name_
}
func (this_ CompilerOptions_DevMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_DevMode[this_-1].Name_
}
func (this_ CompilerOptions_DevMode) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_DevMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_DevMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_DevMode[this_-1].HashCode()
}
func (this_ CompilerOptions_DevMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_DevMode[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_DevMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_DevMode[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_DevMode(1), CompilerOptions_DevMode(2), CompilerOptions_DevMode(3), CompilerOptions_DevMode(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.DevMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$DevMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.DevMode", "com.google.javascript.jscomp.CompilerOptions$DevMode", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_DevMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.TracerMode
type CompilerOptions_TracerMode uint8

const (
	enum_CompilerOptions_TracerMode_ALL                  CompilerOptions_TracerMode = 1
	enum_CompilerOptions_TracerMode_RAW_SIZE             CompilerOptions_TracerMode = 2
	enum_CompilerOptions_TracerMode_AST_SIZE_AND_PRUNING CompilerOptions_TracerMode = 3
	enum_CompilerOptions_TracerMode_AST_SIZE             CompilerOptions_TracerMode = 4
	enum_CompilerOptions_TracerMode_TIMING_ONLY          CompilerOptions_TracerMode = 5
	enum_CompilerOptions_TracerMode_OFF                  CompilerOptions_TracerMode = 6
)

var enumBaseTable_CompilerOptions_TracerMode = [...]j2gort.EnumBase{
	{Name_: "ALL", Ordinal_: 0},
	{Name_: "RAW_SIZE", Ordinal_: 1},
	{Name_: "AST_SIZE_AND_PRUNING", Ordinal_: 2},
	{Name_: "AST_SIZE", Ordinal_: 3},
	{Name_: "TIMING_ONLY", Ordinal_: 4},
	{Name_: "OFF", Ordinal_: 5},
}
var enumValues_CompilerOptions_TracerMode = []CompilerOptions_TracerMode{1, 2, 3, 4, 5, 6}

func StaticGet_CompilerOptions_TracerMode_ALL() CompilerOptions_TracerMode {
	return enum_CompilerOptions_TracerMode_ALL
}
func StaticGet_CompilerOptions_TracerMode_RAW_SIZE() CompilerOptions_TracerMode {
	return enum_CompilerOptions_TracerMode_RAW_SIZE
}
func StaticGet_CompilerOptions_TracerMode_AST_SIZE_AND_PRUNING() CompilerOptions_TracerMode {
	return enum_CompilerOptions_TracerMode_AST_SIZE_AND_PRUNING
}
func StaticGet_CompilerOptions_TracerMode_AST_SIZE() CompilerOptions_TracerMode {
	return enum_CompilerOptions_TracerMode_AST_SIZE
}
func StaticGet_CompilerOptions_TracerMode_TIMING_ONLY() CompilerOptions_TracerMode {
	return enum_CompilerOptions_TracerMode_TIMING_ONLY
}
func StaticGet_CompilerOptions_TracerMode_OFF() CompilerOptions_TracerMode {
	return enum_CompilerOptions_TracerMode_OFF
}

func CompilerOptions_TracerMode_Values() []CompilerOptions_TracerMode {
	return j2gort.CloneSlice(enumValues_CompilerOptions_TracerMode)
}

func CompilerOptions_TracerMode_ValueOf(name string) CompilerOptions_TracerMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_TracerMode {
		if enumBaseTable_CompilerOptions_TracerMode[i_].Name_ == name {
			return CompilerOptions_TracerMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.TracerMode." + name))
}

// Java method: TracerMode.isOn()
func (this_ CompilerOptions_TracerMode) IsOn() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ != StaticGet_CompilerOptions_TracerMode_OFF())
}

// Java method: TracerMode.doPruningAnalysis()
func (this_ CompilerOptions_TracerMode) DoPruningAnalysis() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (((this_ == StaticGet_CompilerOptions_TracerMode_AST_SIZE_AND_PRUNING()) || (this_ == StaticGet_CompilerOptions_TracerMode_ALL())) || (this_ == StaticGet_CompilerOptions_TracerMode_RAW_SIZE()))
}

func (this_ CompilerOptions_TracerMode) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_TracerMode) HasIdentityEquals_()                                 {}
func (CompilerOptions_TracerMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_TracerMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_TracerMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_TracerMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$TracerMode", "TracerMode"
}
func (this_ CompilerOptions_TracerMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_TracerMode[this_-1]
}
func (this_ CompilerOptions_TracerMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_TracerMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_TracerMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_TracerMode[this_-1].Name_
}
func (this_ CompilerOptions_TracerMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_TracerMode[this_-1].Name_
}
func (this_ CompilerOptions_TracerMode) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_TracerMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_TracerMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_TracerMode[this_-1].HashCode()
}
func (this_ CompilerOptions_TracerMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_TracerMode[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_TracerMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_TracerMode[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_TracerMode(1), CompilerOptions_TracerMode(2), CompilerOptions_TracerMode(3), CompilerOptions_TracerMode(4), CompilerOptions_TracerMode(5), CompilerOptions_TracerMode(6)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.TracerMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$TracerMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.TracerMode", "com.google.javascript.jscomp.CompilerOptions$TracerMode", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_TracerMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.TweakProcessing
type CompilerOptions_TweakProcessing uint8

const (
	enum_CompilerOptions_TweakProcessing_OFF   CompilerOptions_TweakProcessing = 1
	enum_CompilerOptions_TweakProcessing_CHECK CompilerOptions_TweakProcessing = 2
	enum_CompilerOptions_TweakProcessing_STRIP CompilerOptions_TweakProcessing = 3
)

var enumBaseTable_CompilerOptions_TweakProcessing = [...]j2gort.EnumBase{
	{Name_: "OFF", Ordinal_: 0},
	{Name_: "CHECK", Ordinal_: 1},
	{Name_: "STRIP", Ordinal_: 2},
}
var enumValues_CompilerOptions_TweakProcessing = []CompilerOptions_TweakProcessing{1, 2, 3}

func StaticGet_CompilerOptions_TweakProcessing_OFF() CompilerOptions_TweakProcessing {
	return enum_CompilerOptions_TweakProcessing_OFF
}
func StaticGet_CompilerOptions_TweakProcessing_CHECK() CompilerOptions_TweakProcessing {
	return enum_CompilerOptions_TweakProcessing_CHECK
}
func StaticGet_CompilerOptions_TweakProcessing_STRIP() CompilerOptions_TweakProcessing {
	return enum_CompilerOptions_TweakProcessing_STRIP
}

func CompilerOptions_TweakProcessing_Values() []CompilerOptions_TweakProcessing {
	return j2gort.CloneSlice(enumValues_CompilerOptions_TweakProcessing)
}

func CompilerOptions_TweakProcessing_ValueOf(name string) CompilerOptions_TweakProcessing {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_TweakProcessing {
		if enumBaseTable_CompilerOptions_TweakProcessing[i_].Name_ == name {
			return CompilerOptions_TweakProcessing(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.TweakProcessing." + name))
}

// Java method: TweakProcessing.isOn()
func (this_ CompilerOptions_TweakProcessing) IsOn() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ != StaticGet_CompilerOptions_TweakProcessing_OFF())
}

// Java method: TweakProcessing.shouldStrip()
func (this_ CompilerOptions_TweakProcessing) ShouldStrip() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ == StaticGet_CompilerOptions_TweakProcessing_STRIP())
}

func (this_ CompilerOptions_TweakProcessing) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_TweakProcessing) HasIdentityEquals_()                                 {}
func (CompilerOptions_TweakProcessing) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_TweakProcessing) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_TweakProcessing) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_TweakProcessing) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$TweakProcessing", "TweakProcessing"
}
func (this_ CompilerOptions_TweakProcessing) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_TweakProcessing[this_-1]
}
func (this_ CompilerOptions_TweakProcessing) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_TweakProcessing) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_TweakProcessing) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_TweakProcessing[this_-1].Name_
}
func (this_ CompilerOptions_TweakProcessing) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_TweakProcessing[this_-1].Name_
}
func (this_ CompilerOptions_TweakProcessing) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_TweakProcessing)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_TweakProcessing) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_TweakProcessing[this_-1].HashCode()
}
func (this_ CompilerOptions_TweakProcessing) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_TweakProcessing[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_TweakProcessing) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_TweakProcessing[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_TweakProcessing(1), CompilerOptions_TweakProcessing(2), CompilerOptions_TweakProcessing(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.TweakProcessing", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$TweakProcessing", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.TweakProcessing", "com.google.javascript.jscomp.CompilerOptions$TweakProcessing", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_TweakProcessing](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.IsolationMode
type CompilerOptions_IsolationMode uint8

const (
	enum_CompilerOptions_IsolationMode_NONE CompilerOptions_IsolationMode = 1
	enum_CompilerOptions_IsolationMode_IIFE CompilerOptions_IsolationMode = 2
)

var enumBaseTable_CompilerOptions_IsolationMode = [...]j2gort.EnumBase{
	{Name_: "NONE", Ordinal_: 0},
	{Name_: "IIFE", Ordinal_: 1},
}
var enumValues_CompilerOptions_IsolationMode = []CompilerOptions_IsolationMode{1, 2}

func StaticGet_CompilerOptions_IsolationMode_NONE() CompilerOptions_IsolationMode {
	return enum_CompilerOptions_IsolationMode_NONE
}
func StaticGet_CompilerOptions_IsolationMode_IIFE() CompilerOptions_IsolationMode {
	return enum_CompilerOptions_IsolationMode_IIFE
}

func CompilerOptions_IsolationMode_Values() []CompilerOptions_IsolationMode {
	return j2gort.CloneSlice(enumValues_CompilerOptions_IsolationMode)
}

func CompilerOptions_IsolationMode_ValueOf(name string) CompilerOptions_IsolationMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_IsolationMode {
		if enumBaseTable_CompilerOptions_IsolationMode[i_].Name_ == name {
			return CompilerOptions_IsolationMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.IsolationMode." + name))
}

func (this_ CompilerOptions_IsolationMode) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_IsolationMode) HasIdentityEquals_()                                 {}
func (CompilerOptions_IsolationMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_IsolationMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_IsolationMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_IsolationMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$IsolationMode", "IsolationMode"
}
func (this_ CompilerOptions_IsolationMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_IsolationMode[this_-1]
}
func (this_ CompilerOptions_IsolationMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_IsolationMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_IsolationMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_IsolationMode[this_-1].Name_
}
func (this_ CompilerOptions_IsolationMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_IsolationMode[this_-1].Name_
}
func (this_ CompilerOptions_IsolationMode) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_IsolationMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_IsolationMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_IsolationMode[this_-1].HashCode()
}
func (this_ CompilerOptions_IsolationMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_IsolationMode[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_IsolationMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_IsolationMode[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_IsolationMode(1), CompilerOptions_IsolationMode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.IsolationMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$IsolationMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.IsolationMode", "com.google.javascript.jscomp.CompilerOptions$IsolationMode", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_IsolationMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.SegmentOfCompilationToRun
type CompilerOptions_SegmentOfCompilationToRun uint8

const (
	enum_CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION              CompilerOptions_SegmentOfCompilationToRun = 1
	enum_CompilerOptions_SegmentOfCompilationToRun_CHECKS                          CompilerOptions_SegmentOfCompilationToRun = 2
	enum_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF        CompilerOptions_SegmentOfCompilationToRun = 3
	enum_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF       CompilerOptions_SegmentOfCompilationToRun = 4
	enum_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS                   CompilerOptions_SegmentOfCompilationToRun = 5
	enum_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_AND_FINALIZATIONS CompilerOptions_SegmentOfCompilationToRun = 6
	enum_CompilerOptions_SegmentOfCompilationToRun_FINALIZATIONS                   CompilerOptions_SegmentOfCompilationToRun = 7
)

var enumBaseTable_CompilerOptions_SegmentOfCompilationToRun = [...]j2gort.EnumBase{
	{Name_: "ENTIRE_COMPILATION", Ordinal_: 0},
	{Name_: "CHECKS", Ordinal_: 1},
	{Name_: "OPTIMIZATIONS_FIRST_HALF", Ordinal_: 2},
	{Name_: "OPTIMIZATIONS_SECOND_HALF", Ordinal_: 3},
	{Name_: "OPTIMIZATIONS", Ordinal_: 4},
	{Name_: "OPTIMIZATIONS_AND_FINALIZATIONS", Ordinal_: 5},
	{Name_: "FINALIZATIONS", Ordinal_: 6},
}
var enumValues_CompilerOptions_SegmentOfCompilationToRun = []CompilerOptions_SegmentOfCompilationToRun{1, 2, 3, 4, 5, 6, 7}

func StaticGet_CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION() CompilerOptions_SegmentOfCompilationToRun {
	return enum_CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION
}
func StaticGet_CompilerOptions_SegmentOfCompilationToRun_CHECKS() CompilerOptions_SegmentOfCompilationToRun {
	return enum_CompilerOptions_SegmentOfCompilationToRun_CHECKS
}
func StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF() CompilerOptions_SegmentOfCompilationToRun {
	return enum_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF
}
func StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF() CompilerOptions_SegmentOfCompilationToRun {
	return enum_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF
}
func StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS() CompilerOptions_SegmentOfCompilationToRun {
	return enum_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS
}
func StaticGet_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_AND_FINALIZATIONS() CompilerOptions_SegmentOfCompilationToRun {
	return enum_CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_AND_FINALIZATIONS
}
func StaticGet_CompilerOptions_SegmentOfCompilationToRun_FINALIZATIONS() CompilerOptions_SegmentOfCompilationToRun {
	return enum_CompilerOptions_SegmentOfCompilationToRun_FINALIZATIONS
}

func CompilerOptions_SegmentOfCompilationToRun_Values() []CompilerOptions_SegmentOfCompilationToRun {
	return j2gort.CloneSlice(enumValues_CompilerOptions_SegmentOfCompilationToRun)
}

func CompilerOptions_SegmentOfCompilationToRun_ValueOf(name string) CompilerOptions_SegmentOfCompilationToRun {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_SegmentOfCompilationToRun {
		if enumBaseTable_CompilerOptions_SegmentOfCompilationToRun[i_].Name_ == name {
			return CompilerOptions_SegmentOfCompilationToRun(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.SegmentOfCompilationToRun." + name))
}

func (this_ CompilerOptions_SegmentOfCompilationToRun) IsJavaNil_() bool                { return this_ == 0 }
func (CompilerOptions_SegmentOfCompilationToRun) HasIdentityEquals_()                   {}
func (CompilerOptions_SegmentOfCompilationToRun) Implements_j2gort_Ext_java_lang_Enum() {}
func (CompilerOptions_SegmentOfCompilationToRun) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ CompilerOptions_SegmentOfCompilationToRun) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_SegmentOfCompilationToRun) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$SegmentOfCompilationToRun", "SegmentOfCompilationToRun"
}
func (this_ CompilerOptions_SegmentOfCompilationToRun) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_SegmentOfCompilationToRun[this_-1]
}
func (this_ CompilerOptions_SegmentOfCompilationToRun) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_SegmentOfCompilationToRun) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_SegmentOfCompilationToRun) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_SegmentOfCompilationToRun[this_-1].Name_
}
func (this_ CompilerOptions_SegmentOfCompilationToRun) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_SegmentOfCompilationToRun[this_-1].Name_
}
func (this_ CompilerOptions_SegmentOfCompilationToRun) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_SegmentOfCompilationToRun)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_SegmentOfCompilationToRun) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_SegmentOfCompilationToRun[this_-1].HashCode()
}
func (this_ CompilerOptions_SegmentOfCompilationToRun) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_SegmentOfCompilationToRun[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_SegmentOfCompilationToRun) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_SegmentOfCompilationToRun[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_SegmentOfCompilationToRun(1), CompilerOptions_SegmentOfCompilationToRun(2), CompilerOptions_SegmentOfCompilationToRun(3), CompilerOptions_SegmentOfCompilationToRun(4), CompilerOptions_SegmentOfCompilationToRun(5), CompilerOptions_SegmentOfCompilationToRun(6), CompilerOptions_SegmentOfCompilationToRun(7)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.SegmentOfCompilationToRun", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$SegmentOfCompilationToRun", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.SegmentOfCompilationToRun", "com.google.javascript.jscomp.CompilerOptions$SegmentOfCompilationToRun", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_SegmentOfCompilationToRun](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.AliasStringsMode
type CompilerOptions_AliasStringsMode uint8

const (
	enum_CompilerOptions_AliasStringsMode_NONE           CompilerOptions_AliasStringsMode = 1
	enum_CompilerOptions_AliasStringsMode_LARGE          CompilerOptions_AliasStringsMode = 2
	enum_CompilerOptions_AliasStringsMode_ALL            CompilerOptions_AliasStringsMode = 3
	enum_CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE CompilerOptions_AliasStringsMode = 4
)

var enumBaseTable_CompilerOptions_AliasStringsMode = [...]j2gort.EnumBase{
	{Name_: "NONE", Ordinal_: 0},
	{Name_: "LARGE", Ordinal_: 1},
	{Name_: "ALL", Ordinal_: 2},
	{Name_: "ALL_AGGRESSIVE", Ordinal_: 3},
}
var enumValues_CompilerOptions_AliasStringsMode = []CompilerOptions_AliasStringsMode{1, 2, 3, 4}

func StaticGet_CompilerOptions_AliasStringsMode_NONE() CompilerOptions_AliasStringsMode {
	return enum_CompilerOptions_AliasStringsMode_NONE
}
func StaticGet_CompilerOptions_AliasStringsMode_LARGE() CompilerOptions_AliasStringsMode {
	return enum_CompilerOptions_AliasStringsMode_LARGE
}
func StaticGet_CompilerOptions_AliasStringsMode_ALL() CompilerOptions_AliasStringsMode {
	return enum_CompilerOptions_AliasStringsMode_ALL
}
func StaticGet_CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE() CompilerOptions_AliasStringsMode {
	return enum_CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE
}

func CompilerOptions_AliasStringsMode_Values() []CompilerOptions_AliasStringsMode {
	return j2gort.CloneSlice(enumValues_CompilerOptions_AliasStringsMode)
}

func CompilerOptions_AliasStringsMode_ValueOf(name string) CompilerOptions_AliasStringsMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_AliasStringsMode {
		if enumBaseTable_CompilerOptions_AliasStringsMode[i_].Name_ == name {
			return CompilerOptions_AliasStringsMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.AliasStringsMode." + name))
}

func (this_ CompilerOptions_AliasStringsMode) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_AliasStringsMode) HasIdentityEquals_()                                 {}
func (CompilerOptions_AliasStringsMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_AliasStringsMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_AliasStringsMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_AliasStringsMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$AliasStringsMode", "AliasStringsMode"
}
func (this_ CompilerOptions_AliasStringsMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_AliasStringsMode[this_-1]
}
func (this_ CompilerOptions_AliasStringsMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_AliasStringsMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_AliasStringsMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_AliasStringsMode[this_-1].Name_
}
func (this_ CompilerOptions_AliasStringsMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_AliasStringsMode[this_-1].Name_
}
func (this_ CompilerOptions_AliasStringsMode) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_AliasStringsMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_AliasStringsMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_AliasStringsMode[this_-1].HashCode()
}
func (this_ CompilerOptions_AliasStringsMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_AliasStringsMode[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_AliasStringsMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_AliasStringsMode[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_AliasStringsMode(1), CompilerOptions_AliasStringsMode(2), CompilerOptions_AliasStringsMode(3), CompilerOptions_AliasStringsMode(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.AliasStringsMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$AliasStringsMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.AliasStringsMode", "com.google.javascript.jscomp.CompilerOptions$AliasStringsMode", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_AliasStringsMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.Environment
type CompilerOptions_Environment uint8

const (
	enum_CompilerOptions_Environment_BROWSER CompilerOptions_Environment = 1
	enum_CompilerOptions_Environment_CUSTOM  CompilerOptions_Environment = 2
)

var enumBaseTable_CompilerOptions_Environment = [...]j2gort.EnumBase{
	{Name_: "BROWSER", Ordinal_: 0},
	{Name_: "CUSTOM", Ordinal_: 1},
}
var enumValues_CompilerOptions_Environment = []CompilerOptions_Environment{1, 2}

func StaticGet_CompilerOptions_Environment_BROWSER() CompilerOptions_Environment {
	return enum_CompilerOptions_Environment_BROWSER
}
func StaticGet_CompilerOptions_Environment_CUSTOM() CompilerOptions_Environment {
	return enum_CompilerOptions_Environment_CUSTOM
}

func CompilerOptions_Environment_Values() []CompilerOptions_Environment {
	return j2gort.CloneSlice(enumValues_CompilerOptions_Environment)
}

func CompilerOptions_Environment_ValueOf(name string) CompilerOptions_Environment {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_Environment {
		if enumBaseTable_CompilerOptions_Environment[i_].Name_ == name {
			return CompilerOptions_Environment(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.Environment." + name))
}

func (this_ CompilerOptions_Environment) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_Environment) HasIdentityEquals_()                                 {}
func (CompilerOptions_Environment) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_Environment) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_Environment) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_Environment) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$Environment", "Environment"
}
func (this_ CompilerOptions_Environment) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_Environment[this_-1]
}
func (this_ CompilerOptions_Environment) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_Environment) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_Environment) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_Environment[this_-1].Name_
}
func (this_ CompilerOptions_Environment) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_Environment[this_-1].Name_
}
func (this_ CompilerOptions_Environment) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_Environment)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_Environment) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_Environment[this_-1].HashCode()
}
func (this_ CompilerOptions_Environment) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_Environment[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_Environment) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_Environment[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_Environment(1), CompilerOptions_Environment(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.Environment", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$Environment", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.Environment", "com.google.javascript.jscomp.CompilerOptions$Environment", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_Environment](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.JsonStreamMode
type CompilerOptions_JsonStreamMode uint8

const (
	enum_CompilerOptions_JsonStreamMode_NONE CompilerOptions_JsonStreamMode = 1
	enum_CompilerOptions_JsonStreamMode_IN   CompilerOptions_JsonStreamMode = 2
	enum_CompilerOptions_JsonStreamMode_OUT  CompilerOptions_JsonStreamMode = 3
	enum_CompilerOptions_JsonStreamMode_BOTH CompilerOptions_JsonStreamMode = 4
)

var enumBaseTable_CompilerOptions_JsonStreamMode = [...]j2gort.EnumBase{
	{Name_: "NONE", Ordinal_: 0},
	{Name_: "IN", Ordinal_: 1},
	{Name_: "OUT", Ordinal_: 2},
	{Name_: "BOTH", Ordinal_: 3},
}
var enumValues_CompilerOptions_JsonStreamMode = []CompilerOptions_JsonStreamMode{1, 2, 3, 4}

func StaticGet_CompilerOptions_JsonStreamMode_NONE() CompilerOptions_JsonStreamMode {
	return enum_CompilerOptions_JsonStreamMode_NONE
}
func StaticGet_CompilerOptions_JsonStreamMode_IN() CompilerOptions_JsonStreamMode {
	return enum_CompilerOptions_JsonStreamMode_IN
}
func StaticGet_CompilerOptions_JsonStreamMode_OUT() CompilerOptions_JsonStreamMode {
	return enum_CompilerOptions_JsonStreamMode_OUT
}
func StaticGet_CompilerOptions_JsonStreamMode_BOTH() CompilerOptions_JsonStreamMode {
	return enum_CompilerOptions_JsonStreamMode_BOTH
}

func CompilerOptions_JsonStreamMode_Values() []CompilerOptions_JsonStreamMode {
	return j2gort.CloneSlice(enumValues_CompilerOptions_JsonStreamMode)
}

func CompilerOptions_JsonStreamMode_ValueOf(name string) CompilerOptions_JsonStreamMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_JsonStreamMode {
		if enumBaseTable_CompilerOptions_JsonStreamMode[i_].Name_ == name {
			return CompilerOptions_JsonStreamMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.JsonStreamMode." + name))
}

func (this_ CompilerOptions_JsonStreamMode) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_JsonStreamMode) HasIdentityEquals_()                                 {}
func (CompilerOptions_JsonStreamMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_JsonStreamMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_JsonStreamMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_JsonStreamMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$JsonStreamMode", "JsonStreamMode"
}
func (this_ CompilerOptions_JsonStreamMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_JsonStreamMode[this_-1]
}
func (this_ CompilerOptions_JsonStreamMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_JsonStreamMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_JsonStreamMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_JsonStreamMode[this_-1].Name_
}
func (this_ CompilerOptions_JsonStreamMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_JsonStreamMode[this_-1].Name_
}
func (this_ CompilerOptions_JsonStreamMode) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_JsonStreamMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_JsonStreamMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_JsonStreamMode[this_-1].HashCode()
}
func (this_ CompilerOptions_JsonStreamMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_JsonStreamMode[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_JsonStreamMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_JsonStreamMode[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_JsonStreamMode(1), CompilerOptions_JsonStreamMode(2), CompilerOptions_JsonStreamMode(3), CompilerOptions_JsonStreamMode(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.JsonStreamMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$JsonStreamMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.JsonStreamMode", "com.google.javascript.jscomp.CompilerOptions$JsonStreamMode", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_JsonStreamMode](v_)
	})
}

// Java enum: com.google.javascript.jscomp.CompilerOptions.J2clPassMode
type CompilerOptions_J2clPassMode uint8

const (
	enum_CompilerOptions_J2clPassMode_OFF  CompilerOptions_J2clPassMode = 1
	enum_CompilerOptions_J2clPassMode_AUTO CompilerOptions_J2clPassMode = 2
)

var enumBaseTable_CompilerOptions_J2clPassMode = [...]j2gort.EnumBase{
	{Name_: "OFF", Ordinal_: 0},
	{Name_: "AUTO", Ordinal_: 1},
}
var enumValues_CompilerOptions_J2clPassMode = []CompilerOptions_J2clPassMode{1, 2}

func StaticGet_CompilerOptions_J2clPassMode_OFF() CompilerOptions_J2clPassMode {
	return enum_CompilerOptions_J2clPassMode_OFF
}
func StaticGet_CompilerOptions_J2clPassMode_AUTO() CompilerOptions_J2clPassMode {
	return enum_CompilerOptions_J2clPassMode_AUTO
}

func CompilerOptions_J2clPassMode_Values() []CompilerOptions_J2clPassMode {
	return j2gort.CloneSlice(enumValues_CompilerOptions_J2clPassMode)
}

func CompilerOptions_J2clPassMode_ValueOf(name string) CompilerOptions_J2clPassMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerOptions_J2clPassMode {
		if enumBaseTable_CompilerOptions_J2clPassMode[i_].Name_ == name {
			return CompilerOptions_J2clPassMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerOptions.J2clPassMode." + name))
}

// Java method: J2clPassMode.shouldAddJ2clPasses()
func (this_ CompilerOptions_J2clPassMode) ShouldAddJ2clPasses() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ == StaticGet_CompilerOptions_J2clPassMode_AUTO())
}

func (this_ CompilerOptions_J2clPassMode) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerOptions_J2clPassMode) HasIdentityEquals_()                                 {}
func (CompilerOptions_J2clPassMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerOptions_J2clPassMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerOptions_J2clPassMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerOptions_J2clPassMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerOptions$J2clPassMode", "J2clPassMode"
}
func (this_ CompilerOptions_J2clPassMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerOptions_J2clPassMode[this_-1]
}
func (this_ CompilerOptions_J2clPassMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerOptions_J2clPassMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerOptions_J2clPassMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_J2clPassMode[this_-1].Name_
}
func (this_ CompilerOptions_J2clPassMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerOptions_J2clPassMode[this_-1].Name_
}
func (this_ CompilerOptions_J2clPassMode) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerOptions_J2clPassMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerOptions_J2clPassMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerOptions_J2clPassMode[this_-1].HashCode()
}
func (this_ CompilerOptions_J2clPassMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerOptions_J2clPassMode[this_-1].CompareTo(other_)
}
func (this_ CompilerOptions_J2clPassMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerOptions_J2clPassMode[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerOptions_J2clPassMode(1), CompilerOptions_J2clPassMode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions.J2clPassMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerOptions$J2clPassMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerOptions.J2clPassMode", "com.google.javascript.jscomp.CompilerOptions$J2clPassMode", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerOptions_J2clPassMode](v_)
	})
}
func (this_ *CompilerOptions) Slice_GetConformanceConfigs() []*j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig {
	return this_.conformanceConfigs
}

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
// Port of com/google/javascript/jscomp/AbstractCompiler.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_jscomptype "github.com/paullepoulpe/jsgo/transpiled/gen/jscomptype"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.AbstractCompiler
type AbstractCompiler interface {
	Implements_jscomp_AbstractCompiler()
	IsJavaNil_() bool
	GetAbstractCompilerBase() *AbstractCompilerBase
	SourceExcerptProvider
	CompilerInputProvider
	BeforePass(string)
	AfterPass(string)
	GetTypedAstDeserializer(*SourceFile) j2gort.Ext_java_util_function_Supplier
	GetSourceFileByName(string) *SourceFile
	GetScriptNode(string) pkg_rhino.Node
	GetChunkGraph() *JSChunkGraph
	GetInputsInOrder() j2gort.Ext_java_lang_Iterable
	GetNumberOfInputs() int32
	AddExportedNames(j2gort.Ext_java_util_Set)
	GetExportedNames() j2gort.Ext_java_util_Set
	SetDefineNames(j2gort.Ext_java_util_Collection)
	GetDefineNames() j2gort.Ext_com_google_common_collect_ImmutableSet
	SetVariableMap(*VariableMap)
	SetPropertyMap(*VariableMap)
	SetStringMap(*VariableMap)
	SetCssNames(j2gort.Ext_java_util_Set)
	SetInstrumentationMapping(*VariableMap)
	SetIdGeneratorMap(string)
	GetTranspiledFiles() bool
	GetCrossChunkIdGenerator() *IdGenerator
	SetAnonymousFunctionNameMap(*VariableMap)
	SetTypeCheckingHasRun(bool)
	HasTypeCheckingRun() bool
	HasOptimizationColors() bool
	IsTypeRegistryCleared() bool
	GetTypeRegistry() *pkg_rhino.JSTypeRegistry
	ClearJSTypeRegistry()
	GetColorRegistry() *pkg_colors.ColorRegistry
	SetColorRegistry(*pkg_colors.ColorRegistry)
	ForwardDeclareType(string)
	GetTypedScopeCreator() ScopeCreator
	ClearTypedScopeCreator()
	GetTopScope() *TypedScope
	SetTopScope(*TypedScope)
	GetTranspilationNamespace() pkg_rhino.StaticScope
	Report_1(*JSError)
	ThrowInternalError(string, javalang.Throwable)
	GetCodingConvention() CodingConvention
	ReportChangeToEnclosingScope(pkg_rhino.Node)
	ReportChangeToChangeScope(pkg_rhino.Node)
	ReportFunctionDeleted(pkg_rhino.Node)
	ReportDisambiguatePropertiesSummary(j2gort.Ext_java_util_function_Supplier)
	ReportAmbiguatePropertiesSummary(j2gort.Ext_java_util_function_Supplier)
	GetNodeForCodeInsertion(*JSChunk) pkg_rhino.Node
	GetTypeValidator() *TypeValidator
	GetTypeMismatches() j2gort.Ext_java_lang_Iterable
	SetExternExports(string)
	ParseSyntheticCode(string, string) pkg_rhino.Node
	ParseTestCode_1_string(string) pkg_rhino.Node
	ParseTestCode_1_Ext_com_google_common_collect_ImmutableList(j2gort.Ext_com_google_common_collect_ImmutableList) pkg_rhino.Node
	ToSource_0() string
	ToSource_1_Node(pkg_rhino.Node) string
	PreferRegexParser() bool
	GetDefaultErrorReporter() pkg_rhino.ErrorReporter
	GetReverseAbstractInterpreter() pkg_jscomptype.ReverseAbstractInterpreter
	GetLifeCycleStage() *AbstractCompiler_LifeCycleStage
	InitRuntimeLibraryTypedAsts(j2gort.Ext_com_google_common_base_Optional)
	GetUniqueIdSupplier() *UniqueIdSupplier
	GetUniqueNameIdSupplier() any
	HasHaltingErrors() bool
	GetChangeTracker() *ChangeTracker
	AddIndexProvider(IndexProvider)
	GetIndex(*javalang.Class) any
	GetJsRoot() pkg_rhino.Node
	HasScopeChanged(pkg_rhino.Node) bool
	GetParserConfig(AbstractCompiler_ConfigContext) pkg_parsing.Config
	GetErrorManager() ErrorManager
	SetLifeCycleStage(*AbstractCompiler_LifeCycleStage)
	AreNodesEqualForInlining(pkg_rhino.Node, pkg_rhino.Node) bool
	SetHasRegExpGlobalReferences(bool)
	HasRegExpGlobalReferences() bool
	GetErrorLevel(*JSError) *CheckLevel
	Process(CompilerPass)
	GetRoot() pkg_rhino.Node
	GetOptions() *CompilerOptions
	GetAllowableFeatures() *pkg_parser.FeatureSet
	SetAllowableFeatures(*pkg_parser.FeatureSet)
	MarkFeatureNotAllowed(*pkg_parser.FeatureSet_Feature)
	MarkFeatureSetNotAllowed(*pkg_parser.FeatureSet)
	GetSynthesizedExternsInput() *CompilerInput
	GetSynthesizedTypeSummaryInput() *CompilerInput
	GetLastPassName() string
	GetRuntimeJsLibManager() *RuntimeJsLibManager
	SetExternProperties(j2gort.Ext_com_google_common_collect_ImmutableSet)
	GetExternProperties() j2gort.Ext_com_google_common_collect_ImmutableSet
	AddInputSourceMap(string, *SourceMapInput)
	GetBase64SourceMapContents(string) string
	AddComments(string, j2gort.Ext_java_util_List)
	GetAccessorSummary() *AccessorSummary
	SetAccessorSummary(*AccessorSummary)
	GetComments(string) j2gort.Ext_java_util_List
	GetModuleLoader() *ModuleLoader
	GetModuleTypeByName(string) CompilerInput_ModuleType
	SetRunJ2clPasses(bool)
	RunJ2clPasses() bool
	CreateAstFactory() *AstFactory
	CreateAstFactoryWithoutTypes() *AstFactory
	GetAstAnalyzer() *AstAnalyzer
	CreateDefaultExpressionDecomposer() *ExpressionDecomposer
	CreateExpressionDecomposer(j2gort.Ext_java_util_function_Supplier, j2gort.Ext_com_google_common_collect_ImmutableSet, *Scope) *ExpressionDecomposer
	GetModuleMetadataMap() *ModuleMetadataMap
	SetModuleMetadataMap(*ModuleMetadataMap)
	GetModuleMap() *ModuleMap
	SetModuleMap(*ModuleMap)
	IsDebugLoggingEnabled() bool
	GetDebugLogFilterList() j2gort.Ext_java_util_List
	CreateOrReopenLog(*javalang.Class, string, []string) pkg_diagnostic.LogFile
	CreateOrReopenIndexedLog(*javalang.Class, string, []string) pkg_diagnostic.LogFile
	GetSyntheticCodeInputId() *pkg_rhino.InputId
	InitializeSyntheticCodeInput()
	RemoveSyntheticCodeInput_0()
	MergeSyntheticCodeInput()
}

func ToIface_AbstractCompiler[T interface {
	AbstractCompiler
	comparable
}](v T) AbstractCompiler {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_AbstractCompiler[T AbstractCompiler](v T) AbstractCompiler {
	var vi AbstractCompiler = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.AbstractCompiler
type AbstractCompilerBase struct {
	idHash_ atomic.Int32
}

func (this_ *AbstractCompilerBase) BindSelf_(s any) {
	_ = s
}
func (this_ *AbstractCompilerBase) IsJavaNil_() bool { return this_ == nil }
func (*AbstractCompilerBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCompiler", "AbstractCompiler"
}
func (this_ *AbstractCompilerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCompiler", "com.google.javascript.jscomp.AbstractCompiler", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractCompiler](v_)
	})
}
func (b *AbstractCompilerBase) GetAbstractCompilerBase() *AbstractCompilerBase { return b }
func (*AbstractCompilerBase) Implements_jscomp_AbstractCompiler()              {}
func (*AbstractCompilerBase) Implements_jscomp_SourceExcerptProvider()         {}
func (*AbstractCompilerBase) Implements_jscomp_CompilerInputProvider()         {}

var static_AbstractCompiler_READ_ERROR *DiagnosticType

func StaticGet_AbstractCompiler_READ_ERROR() *DiagnosticType {
	clinit_AbstractCompiler()
	return static_AbstractCompiler_READ_ERROR
}
func StaticSet_AbstractCompiler_READ_ERROR(v *DiagnosticType) {
	clinit_AbstractCompiler()
	static_AbstractCompiler_READ_ERROR = v
}
func StaticGet_AbstractCompiler_FILL_FILE_SUFFIX() string { return "$fillFile" }
func StaticGet_AbstractCompiler_RUNTIME_LIB_DIR() string {
	return "src/com/google/javascript/jscomp/js/"
}

var clinitDone_AbstractCompiler uint32
var clinitState_AbstractCompiler j2gort.ClinitState

func clinit_AbstractCompiler() {
	if atomic.LoadUint32(&clinitDone_AbstractCompiler) != 0 {
		return
	}
	clinitSlow_AbstractCompiler()
}

func clinitSlow_AbstractCompiler() {
	if !j2gort.EnterClinit(&clinitState_AbstractCompiler, &clinitDone_AbstractCompiler) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AbstractCompiler, &clinitDone_AbstractCompiler)
	static_AbstractCompiler_READ_ERROR = DiagnosticType_Error("JSC_READ_ERROR", "Cannot read file {0}: {1}")
}

func (this_ *AbstractCompilerBase) initFields_AbstractCompiler() {
	_ = this_
}

func InitAbstractCompiler(this_ *AbstractCompilerBase) {
	this_.initFields_AbstractCompiler()
}

// Java method: AbstractCompiler.beforePass(java.lang.String passName)
func (this_ *AbstractCompilerBase) BeforePass(passName string) {
	_ = passName
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.afterPass(java.lang.String passName)
func (this_ *AbstractCompilerBase) AfterPass(passName string) {
	_ = passName
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getTypedAstDeserializer(com.google.javascript.jscomp.SourceFile file)
func (this_ *AbstractCompilerBase) GetTypedAstDeserializer(file *SourceFile) j2gort.Ext_java_util_function_Supplier {
	_ = file
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getInput(com.google.javascript.rhino.InputId inputId)
func (this_ *AbstractCompilerBase) GetInput(inputId *pkg_rhino.InputId) *CompilerInput {
	_ = inputId
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getSourceFileByName(java.lang.String sourceName)
func (this_ *AbstractCompilerBase) GetSourceFileByName(sourceName string) *SourceFile {
	_ = sourceName
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getScriptNode(java.lang.String filename)
func (this_ *AbstractCompilerBase) GetScriptNode(filename string) pkg_rhino.Node {
	_ = filename
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getChunkGraph()
func (this_ *AbstractCompilerBase) GetChunkGraph() *JSChunkGraph {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getInputsInOrder()
func (this_ *AbstractCompilerBase) GetInputsInOrder() j2gort.Ext_java_lang_Iterable {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getNumberOfInputs()
func (this_ *AbstractCompilerBase) GetNumberOfInputs() int32 {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.addExportedNames(java.util.Set<java.lang.String> exportedVariableNames)
func (this_ *AbstractCompilerBase) AddExportedNames(exportedVariableNames j2gort.Ext_java_util_Set) {
	_ = exportedVariableNames
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getExportedNames()
func (this_ *AbstractCompilerBase) GetExportedNames() j2gort.Ext_java_util_Set {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setDefineNames(java.util.Collection<java.lang.String> defineNames)
func (this_ *AbstractCompilerBase) SetDefineNames(defineNames j2gort.Ext_java_util_Collection) {
	_ = defineNames
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getDefineNames()
func (this_ *AbstractCompilerBase) GetDefineNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setVariableMap(com.google.javascript.jscomp.VariableMap variableMap)
func (this_ *AbstractCompilerBase) SetVariableMap(variableMap *VariableMap) {
	_ = variableMap
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setPropertyMap(com.google.javascript.jscomp.VariableMap propertyMap)
func (this_ *AbstractCompilerBase) SetPropertyMap(propertyMap *VariableMap) {
	_ = propertyMap
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setStringMap(com.google.javascript.jscomp.VariableMap stringMap)
func (this_ *AbstractCompilerBase) SetStringMap(stringMap *VariableMap) {
	_ = stringMap
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setCssNames(java.util.Set<java.lang.String> newCssNames)
func (this_ *AbstractCompilerBase) SetCssNames(newCssNames j2gort.Ext_java_util_Set) {
	_ = newCssNames
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setInstrumentationMapping(com.google.javascript.jscomp.VariableMap instrumentationMapping)
func (this_ *AbstractCompilerBase) SetInstrumentationMapping(instrumentationMapping *VariableMap) {
	_ = instrumentationMapping
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setIdGeneratorMap(java.lang.String serializedIdMappings)
func (this_ *AbstractCompilerBase) SetIdGeneratorMap(serializedIdMappings string) {
	_ = serializedIdMappings
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getTranspiledFiles()
func (this_ *AbstractCompilerBase) GetTranspiledFiles() bool {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getCrossChunkIdGenerator()
func (this_ *AbstractCompilerBase) GetCrossChunkIdGenerator() *IdGenerator {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setAnonymousFunctionNameMap(com.google.javascript.jscomp.VariableMap functionMap)
func (this_ *AbstractCompilerBase) SetAnonymousFunctionNameMap(functionMap *VariableMap) {
	_ = functionMap
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setTypeCheckingHasRun(boolean hasRun)
func (this_ *AbstractCompilerBase) SetTypeCheckingHasRun(hasRun bool) {
	_ = hasRun
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.hasTypeCheckingRun()
func (this_ *AbstractCompilerBase) HasTypeCheckingRun() bool {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.hasOptimizationColors()
func (this_ *AbstractCompilerBase) HasOptimizationColors() bool {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.isTypeRegistryCleared()
func (this_ *AbstractCompilerBase) IsTypeRegistryCleared() bool {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getTypeRegistry()
func (this_ *AbstractCompilerBase) GetTypeRegistry() *pkg_rhino.JSTypeRegistry {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.clearJSTypeRegistry()
func (this_ *AbstractCompilerBase) ClearJSTypeRegistry() {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getColorRegistry()
func (this_ *AbstractCompilerBase) GetColorRegistry() *pkg_colors.ColorRegistry {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setColorRegistry(com.google.javascript.jscomp.colors.ColorRegistry registry)
func (this_ *AbstractCompilerBase) SetColorRegistry(registry *pkg_colors.ColorRegistry) {
	_ = registry
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.forwardDeclareType(java.lang.String typeName)
func (this_ *AbstractCompilerBase) ForwardDeclareType(typeName string) {
	_ = typeName
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getTypedScopeCreator()
func (this_ *AbstractCompilerBase) GetTypedScopeCreator() ScopeCreator {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.clearTypedScopeCreator()
func (this_ *AbstractCompilerBase) ClearTypedScopeCreator() {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getTopScope()
func (this_ *AbstractCompilerBase) GetTopScope() *TypedScope {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setTopScope(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope x)
func (this_ *AbstractCompilerBase) SetTopScope(x *TypedScope) {
	_ = x
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getTranspilationNamespace()
func (this_ *AbstractCompilerBase) GetTranspilationNamespace() pkg_rhino.StaticScope {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.report(com.google.javascript.jscomp.JSError error)
func (this_ *AbstractCompilerBase) Report_1(error_ *JSError) {
	_ = error_
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.throwInternalError(java.lang.String msg, java.lang.Throwable cause)
func (this_ *AbstractCompilerBase) ThrowInternalError(msg string, cause javalang.Throwable) {
	_ = msg
	_ = cause
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getCodingConvention()
func (this_ *AbstractCompilerBase) GetCodingConvention() CodingConvention {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.reportChangeToEnclosingScope(com.google.javascript.rhino.Node n)
func (this_ *AbstractCompilerBase) ReportChangeToEnclosingScope(n pkg_rhino.Node) {
	_ = n
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.reportChangeToChangeScope(com.google.javascript.rhino.Node changeScopeRoot)
func (this_ *AbstractCompilerBase) ReportChangeToChangeScope(changeScopeRoot pkg_rhino.Node) {
	_ = changeScopeRoot
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.reportFunctionDeleted(com.google.javascript.rhino.Node node)
func (this_ *AbstractCompilerBase) ReportFunctionDeleted(node pkg_rhino.Node) {
	_ = node
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.reportDisambiguatePropertiesSummary(java.util.function.Supplier<java.lang.String> summarySupplier)
func (this_ *AbstractCompilerBase) ReportDisambiguatePropertiesSummary(summarySupplier j2gort.Ext_java_util_function_Supplier) {
	_ = summarySupplier
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.reportAmbiguatePropertiesSummary(java.util.function.Supplier<java.lang.String> summarySupplier)
func (this_ *AbstractCompilerBase) ReportAmbiguatePropertiesSummary(summarySupplier j2gort.Ext_java_util_function_Supplier) {
	_ = summarySupplier
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getNodeForCodeInsertion(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable JSChunk chunk)
func (this_ *AbstractCompilerBase) GetNodeForCodeInsertion(chunk *JSChunk) pkg_rhino.Node {
	_ = chunk
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getTypeValidator()
func (this_ *AbstractCompilerBase) GetTypeValidator() *TypeValidator {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getTypeMismatches()
func (this_ *AbstractCompilerBase) GetTypeMismatches() j2gort.Ext_java_lang_Iterable {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setExternExports(java.lang.String externExports)
func (this_ *AbstractCompilerBase) SetExternExports(externExports string) {
	_ = externExports
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.parseSyntheticCode(java.lang.String filename, java.lang.String code)
func (this_ *AbstractCompilerBase) ParseSyntheticCode(filename string, code string) pkg_rhino.Node {
	_ = filename
	_ = code
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.parseTestCode(java.lang.String code)
func (this_ *AbstractCompilerBase) ParseTestCode_1_string(code string) pkg_rhino.Node {
	_ = code
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.parseTestCode(com.google.common.collect.ImmutableList<java.lang.String> code)
func (this_ *AbstractCompilerBase) ParseTestCode_1_Ext_com_google_common_collect_ImmutableList(code j2gort.Ext_com_google_common_collect_ImmutableList) pkg_rhino.Node {
	_ = code
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.toSource()
func (this_ *AbstractCompilerBase) ToSource_0() string {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.toSource(com.google.javascript.rhino.Node root)
func (this_ *AbstractCompilerBase) ToSource_1_Node(root pkg_rhino.Node) string {
	_ = root
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.preferRegexParser()
func (this_ *AbstractCompilerBase) PreferRegexParser() bool {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getDefaultErrorReporter()
func (this_ *AbstractCompilerBase) GetDefaultErrorReporter() pkg_rhino.ErrorReporter {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getReverseAbstractInterpreter()
func (this_ *AbstractCompilerBase) GetReverseAbstractInterpreter() pkg_jscomptype.ReverseAbstractInterpreter {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getLifeCycleStage()
func (this_ *AbstractCompilerBase) GetLifeCycleStage() *AbstractCompiler_LifeCycleStage {
	panic("abstract or native method invoked")
}

// Java static method: AbstractCompiler.createFillFileName(java.lang.String chunkName)
func AbstractCompiler_CreateFillFileName(chunkName string) string {
	clinit_AbstractCompiler()
	_ = chunkName
	return (j2gort.StrValString(chunkName) + "$fillFile")
}

// Java static method: AbstractCompiler.isFillFileName(java.lang.String fileName)
func AbstractCompiler_IsFillFileName(fileName string) bool {
	clinit_AbstractCompiler()
	_ = fileName
	return javalang.StringEndsWith(j2gort.CheckNotNullString(fileName), "$fillFile")
}

// Java method: AbstractCompiler.initRuntimeLibraryTypedAsts(com.google.common.base.Optional<com.google.javascript.jscomp.serialization.ColorPool.Builder> colorPoolBuilder)
func (this_ *AbstractCompilerBase) InitRuntimeLibraryTypedAsts(colorPoolBuilder j2gort.Ext_com_google_common_base_Optional) {
	_ = colorPoolBuilder
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Implementation in Compiler.java is not J2CL compatible.")))
}

// Java method: AbstractCompiler.getUniqueIdSupplier()
func (this_ *AbstractCompilerBase) GetUniqueIdSupplier() *UniqueIdSupplier {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getUniqueNameIdSupplier()
func (this_ *AbstractCompilerBase) GetUniqueNameIdSupplier() any {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.hasHaltingErrors()
func (this_ *AbstractCompilerBase) HasHaltingErrors() bool {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getChangeTracker()
func (this_ *AbstractCompilerBase) GetChangeTracker() *ChangeTracker {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.addIndexProvider(com.google.javascript.jscomp.IndexProvider<?> indexProvider)
func (this_ *AbstractCompilerBase) AddIndexProvider(indexProvider IndexProvider) {
	_ = indexProvider
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getIndex(java.lang.Class<T> type)
func (this_ *AbstractCompilerBase) GetIndex(type_ *javalang.Class) any {
	_ = type_
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getJsRoot()
func (this_ *AbstractCompilerBase) GetJsRoot() pkg_rhino.Node {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.hasScopeChanged(com.google.javascript.rhino.Node n)
func (this_ *AbstractCompilerBase) HasScopeChanged(n pkg_rhino.Node) bool {
	_ = n
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getParserConfig(com.google.javascript.jscomp.AbstractCompiler.ConfigContext context)
func (this_ *AbstractCompilerBase) GetParserConfig(context AbstractCompiler_ConfigContext) pkg_parsing.Config {
	_ = context
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getErrorManager()
func (this_ *AbstractCompilerBase) GetErrorManager() ErrorManager {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setLifeCycleStage(com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage stage)
func (this_ *AbstractCompilerBase) SetLifeCycleStage(stage *AbstractCompiler_LifeCycleStage) {
	_ = stage
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.areNodesEqualForInlining(com.google.javascript.rhino.Node n1, com.google.javascript.rhino.Node n2)
func (this_ *AbstractCompilerBase) AreNodesEqualForInlining(n1 pkg_rhino.Node, n2 pkg_rhino.Node) bool {
	_ = n1
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setHasRegExpGlobalReferences(boolean references)
func (this_ *AbstractCompilerBase) SetHasRegExpGlobalReferences(references bool) {
	_ = references
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.hasRegExpGlobalReferences()
func (this_ *AbstractCompilerBase) HasRegExpGlobalReferences() bool {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getErrorLevel(com.google.javascript.jscomp.JSError error)
func (this_ *AbstractCompilerBase) GetErrorLevel(error_ *JSError) *CheckLevel {
	_ = error_
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.process(com.google.javascript.jscomp.CompilerPass pass)
func (this_ *AbstractCompilerBase) Process(pass CompilerPass) {
	_ = pass
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getRoot()
func (this_ *AbstractCompilerBase) GetRoot() pkg_rhino.Node {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getOptions()
func (this_ *AbstractCompilerBase) GetOptions() *CompilerOptions {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getAllowableFeatures()
func (this_ *AbstractCompilerBase) GetAllowableFeatures() *pkg_parser.FeatureSet {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setAllowableFeatures(com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet)
func (this_ *AbstractCompilerBase) SetAllowableFeatures(featureSet *pkg_parser.FeatureSet) {
	_ = featureSet
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.markFeatureNotAllowed(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *AbstractCompilerBase) MarkFeatureNotAllowed(feature *pkg_parser.FeatureSet_Feature) {
	_ = feature
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.markFeatureSetNotAllowed(com.google.javascript.jscomp.parsing.parser.FeatureSet featureSet)
func (this_ *AbstractCompilerBase) MarkFeatureSetNotAllowed(featureSet *pkg_parser.FeatureSet) {
	_ = featureSet
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getSynthesizedExternsInput()
func (this_ *AbstractCompilerBase) GetSynthesizedExternsInput() *CompilerInput {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getSynthesizedTypeSummaryInput()
func (this_ *AbstractCompilerBase) GetSynthesizedTypeSummaryInput() *CompilerInput {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getLastPassName()
func (this_ *AbstractCompilerBase) GetLastPassName() string {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getRuntimeJsLibManager()
func (this_ *AbstractCompilerBase) GetRuntimeJsLibManager() *RuntimeJsLibManager {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setExternProperties(com.google.common.collect.ImmutableSet<java.lang.String> externProperties)
func (this_ *AbstractCompilerBase) SetExternProperties(externProperties j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = externProperties
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getExternProperties()
func (this_ *AbstractCompilerBase) GetExternProperties() j2gort.Ext_com_google_common_collect_ImmutableSet {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.addInputSourceMap(java.lang.String name, com.google.javascript.jscomp.SourceMapInput sourceMap)
func (this_ *AbstractCompilerBase) AddInputSourceMap(name string, sourceMap *SourceMapInput) {
	_ = name
	_ = sourceMap
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getBase64SourceMapContents(java.lang.String sourceFileName)
func (this_ *AbstractCompilerBase) GetBase64SourceMapContents(sourceFileName string) string {
	_ = sourceFileName
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.addComments(java.lang.String filename, java.util.List<com.google.javascript.jscomp.parsing.parser.trees.Comment> comments)
func (this_ *AbstractCompilerBase) AddComments(filename string, comments j2gort.Ext_java_util_List) {
	_ = filename
	_ = comments
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getAccessorSummary()
func (this_ *AbstractCompilerBase) GetAccessorSummary() *AccessorSummary {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setAccessorSummary(com.google.javascript.jscomp.AccessorSummary summary)
func (this_ *AbstractCompilerBase) SetAccessorSummary(summary *AccessorSummary) {
	_ = summary
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getComments(java.lang.String filename)
func (this_ *AbstractCompilerBase) GetComments(filename string) j2gort.Ext_java_util_List {
	_ = filename
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getModuleLoader()
func (this_ *AbstractCompilerBase) GetModuleLoader() *ModuleLoader {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getModuleTypeByName(java.lang.String moduleName)
func (this_ *AbstractCompilerBase) GetModuleTypeByName(moduleName string) CompilerInput_ModuleType {
	_ = moduleName
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setRunJ2clPasses(boolean runJ2clPasses)
func (this_ *AbstractCompilerBase) SetRunJ2clPasses(runJ2clPasses bool) {
	_ = runJ2clPasses
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.runJ2clPasses()
func (this_ *AbstractCompilerBase) RunJ2clPasses() bool {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.createAstFactory()
func (this_ *AbstractCompilerBase) CreateAstFactory() *AstFactory {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.createAstFactoryWithoutTypes()
func (this_ *AbstractCompilerBase) CreateAstFactoryWithoutTypes() *AstFactory {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getAstAnalyzer()
func (this_ *AbstractCompilerBase) GetAstAnalyzer() *AstAnalyzer {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.createDefaultExpressionDecomposer()
func (this_ *AbstractCompilerBase) CreateDefaultExpressionDecomposer() *ExpressionDecomposer {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.createExpressionDecomposer(java.util.function.Supplier<java.lang.String> uniqueNameIdSupplier, com.google.common.collect.ImmutableSet<java.lang.String> knownConstantFunctions, com.google.javascript.jscomp.Scope scope)
func (this_ *AbstractCompilerBase) CreateExpressionDecomposer(uniqueNameIdSupplier j2gort.Ext_java_util_function_Supplier, knownConstantFunctions j2gort.Ext_com_google_common_collect_ImmutableSet, scope *Scope) *ExpressionDecomposer {
	_ = uniqueNameIdSupplier
	_ = knownConstantFunctions
	_ = scope
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getModuleMetadataMap()
func (this_ *AbstractCompilerBase) GetModuleMetadataMap() *ModuleMetadataMap {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setModuleMetadataMap(com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap)
func (this_ *AbstractCompilerBase) SetModuleMetadataMap(moduleMetadataMap *ModuleMetadataMap) {
	_ = moduleMetadataMap
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getModuleMap()
func (this_ *AbstractCompilerBase) GetModuleMap() *ModuleMap {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.setModuleMap(com.google.javascript.jscomp.modules.ModuleMap moduleMap)
func (this_ *AbstractCompilerBase) SetModuleMap(moduleMap *ModuleMap) {
	_ = moduleMap
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.isDebugLoggingEnabled()
func (this_ *AbstractCompilerBase) IsDebugLoggingEnabled() bool {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getDebugLogFilterList()
func (this_ *AbstractCompilerBase) GetDebugLogFilterList() j2gort.Ext_java_util_List {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.createOrReopenLog(java.lang.Class<?> owner, java.lang.String firstNamePart, java.lang.String[] restNameParts)
func (this_ *AbstractCompilerBase) CreateOrReopenLog(owner *javalang.Class, firstNamePart string, restNameParts []string) pkg_diagnostic.LogFile {
	_ = owner
	_ = firstNamePart
	_ = restNameParts
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.createOrReopenIndexedLog(java.lang.Class<?> owner, java.lang.String firstNamePart, java.lang.String[] restNameParts)
func (this_ *AbstractCompilerBase) CreateOrReopenIndexedLog(owner *javalang.Class, firstNamePart string, restNameParts []string) pkg_diagnostic.LogFile {
	_ = owner
	_ = firstNamePart
	_ = restNameParts
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.getSyntheticCodeInputId()
func (this_ *AbstractCompilerBase) GetSyntheticCodeInputId() *pkg_rhino.InputId {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.initializeSyntheticCodeInput()
func (this_ *AbstractCompilerBase) InitializeSyntheticCodeInput() {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.removeSyntheticCodeInput()
func (this_ *AbstractCompilerBase) RemoveSyntheticCodeInput_0() {
	panic("abstract or native method invoked")
}

// Java method: AbstractCompiler.mergeSyntheticCodeInput()
func (this_ *AbstractCompilerBase) MergeSyntheticCodeInput() {
	panic("abstract or native method invoked")
}

// Java enum: com.google.javascript.jscomp.AbstractCompiler.ConfigContext
type AbstractCompiler_ConfigContext uint8

const (
	enum_AbstractCompiler_ConfigContext_DEFAULT AbstractCompiler_ConfigContext = 1
	enum_AbstractCompiler_ConfigContext_EXTERNS AbstractCompiler_ConfigContext = 2
)

var enumBaseTable_AbstractCompiler_ConfigContext = [...]j2gort.EnumBase{
	{Name_: "DEFAULT", Ordinal_: 0},
	{Name_: "EXTERNS", Ordinal_: 1},
}
var enumValues_AbstractCompiler_ConfigContext = []AbstractCompiler_ConfigContext{1, 2}

func StaticGet_AbstractCompiler_ConfigContext_DEFAULT() AbstractCompiler_ConfigContext {
	return enum_AbstractCompiler_ConfigContext_DEFAULT
}
func StaticGet_AbstractCompiler_ConfigContext_EXTERNS() AbstractCompiler_ConfigContext {
	return enum_AbstractCompiler_ConfigContext_EXTERNS
}

func AbstractCompiler_ConfigContext_Values() []AbstractCompiler_ConfigContext {
	return j2gort.CloneSlice(enumValues_AbstractCompiler_ConfigContext)
}

func AbstractCompiler_ConfigContext_ValueOf(name string) AbstractCompiler_ConfigContext {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_AbstractCompiler_ConfigContext {
		if enumBaseTable_AbstractCompiler_ConfigContext[i_].Name_ == name {
			return AbstractCompiler_ConfigContext(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.AbstractCompiler.ConfigContext." + name))
}

func (this_ AbstractCompiler_ConfigContext) IsJavaNil_() bool                              { return this_ == 0 }
func (AbstractCompiler_ConfigContext) HasIdentityEquals_()                                 {}
func (AbstractCompiler_ConfigContext) Implements_j2gort_Ext_java_lang_Enum()               {}
func (AbstractCompiler_ConfigContext) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ AbstractCompiler_ConfigContext) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (AbstractCompiler_ConfigContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCompiler$ConfigContext", "ConfigContext"
}
func (this_ AbstractCompiler_ConfigContext) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_AbstractCompiler_ConfigContext[this_-1]
}
func (this_ AbstractCompiler_ConfigContext) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ AbstractCompiler_ConfigContext) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ AbstractCompiler_ConfigContext) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AbstractCompiler_ConfigContext[this_-1].Name_
}
func (this_ AbstractCompiler_ConfigContext) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_AbstractCompiler_ConfigContext[this_-1].Name_
}
func (this_ AbstractCompiler_ConfigContext) Equals(other_ any) bool {
	o_, ok_ := other_.(AbstractCompiler_ConfigContext)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ AbstractCompiler_ConfigContext) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_AbstractCompiler_ConfigContext[this_-1].HashCode()
}
func (this_ AbstractCompiler_ConfigContext) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AbstractCompiler_ConfigContext[this_-1].CompareTo(other_)
}
func (this_ AbstractCompiler_ConfigContext) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_AbstractCompiler_ConfigContext[this_-1].Name_)
}
func init() {
	vals_ := []any{AbstractCompiler_ConfigContext(1), AbstractCompiler_ConfigContext(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractCompiler.ConfigContext", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractCompiler$ConfigContext", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCompiler.ConfigContext", "com.google.javascript.jscomp.AbstractCompiler$ConfigContext", func(v_ any) bool {
		return j2gort.InstanceOf[AbstractCompiler_ConfigContext](v_)
	})
}

// Java enum: com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage
type AbstractCompiler_LifeCycleStage struct {
	j2gort.EnumBase
}

func (this_ *AbstractCompiler_LifeCycleStage) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *AbstractCompiler_LifeCycleStage) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*AbstractCompiler_LifeCycleStage) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*AbstractCompiler_LifeCycleStage) Implements_j2gort_Ext_java_lang_constant_Constable() {}

var enumBacking_AbstractCompiler_LifeCycleStage = [4]AbstractCompiler_LifeCycleStage{
	{EnumBase: j2gort.EnumBase{Name_: "RAW", Ordinal_: 0}},
	{EnumBase: j2gort.EnumBase{Name_: "COLORS_AND_SIMPLIFIED_JSDOC", Ordinal_: 1}},
	{EnumBase: j2gort.EnumBase{Name_: "NORMALIZED", Ordinal_: 2}},
	{EnumBase: j2gort.EnumBase{Name_: "NORMALIZED_OBFUSCATED", Ordinal_: 3}},
}
var enum_AbstractCompiler_LifeCycleStage_RAW *AbstractCompiler_LifeCycleStage = &enumBacking_AbstractCompiler_LifeCycleStage[0]

func StaticGet_AbstractCompiler_LifeCycleStage_RAW() *AbstractCompiler_LifeCycleStage {
	return &enumBacking_AbstractCompiler_LifeCycleStage[0]
}

var enum_AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC *AbstractCompiler_LifeCycleStage = &enumBacking_AbstractCompiler_LifeCycleStage[1]

func StaticGet_AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC() *AbstractCompiler_LifeCycleStage {
	return &enumBacking_AbstractCompiler_LifeCycleStage[1]
}

var enum_AbstractCompiler_LifeCycleStage_NORMALIZED *AbstractCompiler_LifeCycleStage = &enumBacking_AbstractCompiler_LifeCycleStage[2]

func StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED() *AbstractCompiler_LifeCycleStage {
	return &enumBacking_AbstractCompiler_LifeCycleStage[2]
}

var enum_AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED *AbstractCompiler_LifeCycleStage = &enumBacking_AbstractCompiler_LifeCycleStage[3]

func StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED() *AbstractCompiler_LifeCycleStage {
	return &enumBacking_AbstractCompiler_LifeCycleStage[3]
}

var enumValues_AbstractCompiler_LifeCycleStage = []*AbstractCompiler_LifeCycleStage{&enumBacking_AbstractCompiler_LifeCycleStage[0], &enumBacking_AbstractCompiler_LifeCycleStage[1], &enumBacking_AbstractCompiler_LifeCycleStage[2], &enumBacking_AbstractCompiler_LifeCycleStage[3]}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage", "com.google.javascript.jscomp.AbstractCompiler$LifeCycleStage", func() {
		vals_ := make([]any, len(enumValues_AbstractCompiler_LifeCycleStage))
		for i_, v_ := range enumValues_AbstractCompiler_LifeCycleStage {
			vals_[i_] = v_
		}
		j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage", vals_)
		j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AbstractCompiler$LifeCycleStage", vals_)
	})
}

func AbstractCompiler_LifeCycleStage_Values() []*AbstractCompiler_LifeCycleStage {
	return j2gort.CloneSlice(enumValues_AbstractCompiler_LifeCycleStage)
}

func AbstractCompiler_LifeCycleStage_ValueOf(name string) *AbstractCompiler_LifeCycleStage {
	for _, c := range enumValues_AbstractCompiler_LifeCycleStage {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage." + name))
}

// Java method: LifeCycleStage.isNormalized()
func (this_ *AbstractCompiler_LifeCycleStage) IsNormalized() bool {
	return ((this_ == StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED()) || (this_ == StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED()))
}

// Java method: LifeCycleStage.isNormalizedUnobfuscated()
func (this_ *AbstractCompiler_LifeCycleStage) IsNormalizedUnobfuscated() bool {
	return (this_ == StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED())
}

// Java method: LifeCycleStage.isNormalizedObfuscated()
func (this_ *AbstractCompiler_LifeCycleStage) IsNormalizedObfuscated() bool {
	return (this_ == StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED())
}

// Java method: LifeCycleStage.hasColorAndSimplifiedJSDoc()
func (this_ *AbstractCompiler_LifeCycleStage) HasColorAndSimplifiedJSDoc() bool {
	return (this_ != StaticGet_AbstractCompiler_LifeCycleStage_RAW())
}

func (this_ *AbstractCompiler_LifeCycleStage) IsJavaNil_() bool { return this_ == nil }
func (*AbstractCompiler_LifeCycleStage) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AbstractCompiler$LifeCycleStage", "LifeCycleStage"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage", "com.google.javascript.jscomp.AbstractCompiler$LifeCycleStage", func(v_ any) bool {
		return j2gort.InstanceOf[*AbstractCompiler_LifeCycleStage](v_)
	})
}

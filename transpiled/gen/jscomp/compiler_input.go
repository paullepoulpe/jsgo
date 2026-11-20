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
// Port of com/google/javascript/jscomp/CompilerInput.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CompilerInput
type CompilerInput struct {
	chunk                        *JSChunk                                           `gson:"chunk"`
	id                           *pkg_rhino.InputId                                 `gson:"id"`
	sourceFile                   *SourceFile                                        `gson:"sourceFile"`
	ast                          *CompilerInput_JsAst                               `gson:"ast"`
	dependencyInfo               DependencyInfo                                     `gson:"dependencyInfo"`
	extraRequires                j2gort.Ext_java_util_List                          `gson:"extraRequires"`
	extraProvides                j2gort.Ext_java_util_List                          `gson:"extraProvides"`
	orderedRequires              j2gort.Ext_java_util_List                          `gson:"orderedRequires"`
	cachedOrderedRequires        j2gort.Ext_com_google_common_collect_ImmutableList `gson:"-"`
	cachedOrderedRequiredSymbols j2gort.Ext_com_google_common_collect_ImmutableList `gson:"-"`
	dynamicRequires              j2gort.Ext_java_util_List                          `gson:"dynamicRequires"`
	requireDynamicImports        j2gort.Ext_java_util_List                          `gson:"requireDynamicImports"`
	typedScope                   *TypedScope                                        `gson:"typedScope"`
	compiler                     AbstractCompiler                                   `gson:"compiler"`
	modulePath                   *ModuleLoader_ModulePath                           `gson:"modulePath"`
	idHash_                      atomic.Int32
	hasFullParseDependencyInfo   bool                     `gson:"hasFullParseDependencyInfo"`
	jsModuleType                 CompilerInput_ModuleType `gson:"jsModuleType"`
}

func (this_ *CompilerInput) BindSelf_(s any) {
	_ = s
}
func (this_ *CompilerInput) IsJavaNil_() bool { return this_ == nil }
func (*CompilerInput) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerInput", "CompilerInput"
}
func (this_ *CompilerInput) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CompilerInput) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerInput", "com.google.javascript.jscomp.CompilerInput", func(v_ any) bool {
		return j2gort.InstanceOf[*CompilerInput](v_)
	})
}
func (*CompilerInput) Implements_jscomp_DependencyInfo() {}

func StaticGet_CompilerInput_SerialVersionUID() int64 { return int64(2) }
func (this_ *CompilerInput) initFields_CompilerInput() {
	_ = this_
	this_.ast = NewCompilerInput_JsAst(this_)
	this_.extraRequires = j2gort.NewExt_java_util_ArrayList()
	this_.extraProvides = j2gort.NewExt_java_util_ArrayList()
	this_.orderedRequires = j2gort.NewExt_java_util_ArrayList()
	this_.dynamicRequires = j2gort.NewExt_java_util_ArrayList()
	this_.requireDynamicImports = j2gort.NewExt_java_util_ArrayList()

	this_.jsModuleType = StaticGet_CompilerInput_ModuleType_NONE()
}

// Java constructor: CompilerInput(com.google.javascript.jscomp.SourceFile file, com.google.javascript.rhino.InputId inputId)
func NewCompilerInput_2_SourceFile_InputId(file *SourceFile, inputId *pkg_rhino.InputId) *CompilerInput {
	this_ := &CompilerInput{}
	InitCompilerInput_2_SourceFile_InputId(this_, file, inputId)
	return this_
}

func InitCompilerInput_2_SourceFile_InputId(this_ *CompilerInput, file *SourceFile, inputId *pkg_rhino.InputId) {
	_ = file
	_ = inputId
	InitCompilerInput_3_SourceFile_InputId_bool(this_, file, inputId, false)
}

// Java constructor: CompilerInput(com.google.javascript.jscomp.SourceFile sourceFile, java.lang.String inputId, boolean isExtern)
func NewCompilerInput_3_SourceFile_string_bool(sourceFile *SourceFile, inputId string, isExtern bool) *CompilerInput {
	this_ := &CompilerInput{}
	InitCompilerInput_3_SourceFile_string_bool(this_, sourceFile, inputId, isExtern)
	return this_
}

func InitCompilerInput_3_SourceFile_string_bool(this_ *CompilerInput, sourceFile *SourceFile, inputId string, isExtern bool) {
	_ = sourceFile
	_ = inputId
	_ = isExtern
	InitCompilerInput_3_SourceFile_InputId_bool(this_, sourceFile, pkg_rhino.NewInputId(inputId), isExtern)
}

// Java constructor: CompilerInput(com.google.javascript.jscomp.SourceFile sourceFile, com.google.javascript.rhino.InputId inputId, boolean isExtern)
func NewCompilerInput_3_SourceFile_InputId_bool(sourceFile *SourceFile, inputId *pkg_rhino.InputId, isExtern bool) *CompilerInput {
	this_ := &CompilerInput{}
	InitCompilerInput_3_SourceFile_InputId_bool(this_, sourceFile, inputId, isExtern)
	return this_
}

func InitCompilerInput_3_SourceFile_InputId_bool(this_ *CompilerInput, sourceFile *SourceFile, inputId *pkg_rhino.InputId, isExtern bool) {
	_ = sourceFile
	_ = inputId
	_ = isExtern
	this_.initFields_CompilerInput()
	this_.sourceFile = j2gort.CheckNotNullPtr(sourceFile)
	this_.id = j2gort.CheckNotNullPtr(inputId)
	if isExtern {
		this_.SetIsExtern()
	}
}

// Java constructor: CompilerInput(com.google.javascript.jscomp.SourceFile file)
func NewCompilerInput_1(file *SourceFile) *CompilerInput {
	this_ := &CompilerInput{}
	InitCompilerInput_1(this_, file)
	return this_
}

func InitCompilerInput_1(this_ *CompilerInput, file *SourceFile) {
	_ = file
	InitCompilerInput_2_SourceFile_bool(this_, file, false)
}

// Java constructor: CompilerInput(com.google.javascript.jscomp.SourceFile file, boolean isExtern)
func NewCompilerInput_2_SourceFile_bool(file *SourceFile, isExtern bool) *CompilerInput {
	this_ := &CompilerInput{}
	InitCompilerInput_2_SourceFile_bool(this_, file, isExtern)
	return this_
}

func InitCompilerInput_2_SourceFile_bool(this_ *CompilerInput, file *SourceFile, isExtern bool) {
	_ = file
	_ = isExtern
	InitCompilerInput_3_SourceFile_InputId_bool(this_, file, pkg_rhino.NewInputId(file.GetName()), isExtern)
}

// Java method: CompilerInput.setTypedScope(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope typedScope)
func (this_ *CompilerInput) SetTypedScope(typedScope *TypedScope) {
	_ = typedScope
	this_.typedScope = typedScope
}

// Java method: CompilerInput.getTypedScope()
func (this_ *CompilerInput) GetTypedScope() *TypedScope {
	return this_.typedScope
}

// Java method: CompilerInput.getInputId()
func (this_ *CompilerInput) GetInputId() *pkg_rhino.InputId {
	return this_.id
}

// Java method: CompilerInput.getName()
func (this_ *CompilerInput) GetName() string {
	return this_.id.GetIdName()
}

// Java method: CompilerInput.getPathRelativeToClosureBase()
func (this_ *CompilerInput) GetPathRelativeToClosureBase() string {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: CompilerInput.getAstRoot(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *CompilerInput) GetAstRoot(compiler AbstractCompiler) pkg_rhino.Node {
	_ = compiler
	var root pkg_rhino.Node = j2gort.CheckNotNullNode(this_.ast.GetAstRoot(compiler))
	_ = root
	if !(root.IsScript()) {
		j2gort.CheckState(false)
	}
	j2gort.CheckNotNullPtr(root.GetInputId())
	return root
}

// Java method: CompilerInput.clearAst()
func (this_ *CompilerInput) ClearAst() {
	this_.ast.ClearAst()
}

// Java method: CompilerInput.initShadowAst(com.google.javascript.rhino.Node shadowScript)
func (this_ *CompilerInput) InitShadowAst(shadowScript pkg_rhino.Node) {
	_ = shadowScript
	if !(shadowScript.IsScript()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(shadowScript))
	}
	if !(shadowScript.GetIsInClosureUnawareSubtree()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(shadowScript))
	}
	{
		r_1 := this_.ast
		r_1.root = shadowScript
	}
}

// Java method: CompilerInput.getSourceFile()
func (this_ *CompilerInput) GetSourceFile() *SourceFile {
	return this_.sourceFile
}

// Java method: CompilerInput.setCompiler(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *CompilerInput) SetCompiler(compiler AbstractCompiler) {
	_ = compiler
	this_.compiler = compiler
}

// Java method: CompilerInput.getRequires()
func (this_ *CompilerInput) GetRequires() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.hasFullParseDependencyInfo {
		var reqs j2gort.Ext_com_google_common_collect_ImmutableList = this_.cachedOrderedRequires
		_ = reqs
		if reqs == nil {
			reqs = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.orderedRequires))
			this_.cachedOrderedRequires = reqs
		}
		return reqs
	}
	return this_.GetDependencyInfo().GetRequires()
}

// Java method: CompilerInput.getRequiredSymbols()
func (this_ *CompilerInput) GetRequiredSymbols() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.hasFullParseDependencyInfo {
		var syms j2gort.Ext_com_google_common_collect_ImmutableList = this_.cachedOrderedRequiredSymbols
		_ = syms
		if syms == nil {
			syms = DependencyInfo_Require_AsSymbolList(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.orderedRequires))
			this_.cachedOrderedRequiredSymbols = syms
		}
		return syms
	}
	return this_.GetDependencyInfo().GetRequiredSymbols()
}

// Java method: CompilerInput.getTypeRequires()
func (this_ *CompilerInput) GetTypeRequires() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.GetDependencyInfo().GetTypeRequires()
}

// Java method: CompilerInput.getKnownRequires()
func (this_ *CompilerInput) GetKnownRequires() j2gort.Ext_com_google_common_collect_ImmutableCollection {
	if this_.extraRequires.IsEmpty() {
		if this_.dependencyInfo != nil && !this_.dependencyInfo.IsJavaNil_() {
			return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](this_.dependencyInfo.GetRequires())
		} else {
			return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
		}
	}
	var iifeArg_1 j2gort.
		Ext_java_lang_Iterable
iifeSwitch_1:
	switch {
	default:

		if this_.dependencyInfo != nil && !this_.dependencyInfo.IsJavaNil_() {
			iifeArg_1 = j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.dependencyInfo.GetRequires())
			break iifeSwitch_1
		}
		iifeArg_1 = j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
		break iifeSwitch_1
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](CompilerInput_Concat(iifeArg_1, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.extraRequires)))
}

// Java method: CompilerInput.getKnownRequiredSymbols()
func (this_ *CompilerInput) GetKnownRequiredSymbols() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.extraRequires.IsEmpty() {
		if this_.dependencyInfo != nil && !this_.dependencyInfo.IsJavaNil_() {
			return this_.dependencyInfo.GetRequiredSymbols()
		} else {
			return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
		}
	}
	return DependencyInfo_Require_AsSymbolList(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.GetKnownRequires()))
}

// Java method: CompilerInput.getProvides()
func (this_ *CompilerInput) GetProvides() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.GetDependencyInfo().GetProvides()
}

// Java method: CompilerInput.getHasExternsAnnotation()
func (this_ *CompilerInput) GetHasExternsAnnotation() bool {
	return this_.GetDependencyInfo().GetHasExternsAnnotation()
}

// Java method: CompilerInput.getHasNoCompileAnnotation()
func (this_ *CompilerInput) GetHasNoCompileAnnotation() bool {
	return this_.GetDependencyInfo().GetHasNoCompileAnnotation()
}

// Java method: CompilerInput.getKnownProvides()
func (this_ *CompilerInput) GetKnownProvides() j2gort.Ext_com_google_common_collect_ImmutableCollection {
	if this_.extraProvides.IsEmpty() {
		if this_.dependencyInfo != nil && !this_.dependencyInfo.IsJavaNil_() {
			return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](this_.dependencyInfo.GetProvides())
		} else {
			return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
		}
	}
	var iifeArg_2 j2gort.
		Ext_java_lang_Iterable
iifeSwitch_2:
	switch {
	default:

		if this_.dependencyInfo != nil && !this_.dependencyInfo.IsJavaNil_() {
			iifeArg_2 = j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.dependencyInfo.GetProvides())
			break iifeSwitch_2
		}
		iifeArg_2 = j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
		break iifeSwitch_2
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](CompilerInput_Concat(iifeArg_2, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.extraProvides)))
}

// Java method: CompilerInput.addProvide(java.lang.String provide)
func (this_ *CompilerInput) AddProvide(provide string) {
	_ = provide
	this_.extraProvides.Add_1(j2gort.BoxString(provide))
}

// Java method: CompilerInput.addOrderedRequire(com.google.javascript.jscomp.deps.DependencyInfo.Require require)
func (this_ *CompilerInput) AddOrderedRequire(require *DependencyInfo_Require) bool {
	_ = require
	if !(this_.orderedRequires.Contains(j2gort.PtrToAny(require))) {
		this_.orderedRequires.Add_1(j2gort.PtrToAny(require))
		this_.cachedOrderedRequires = nil
		this_.cachedOrderedRequiredSymbols = nil
		return true
	}
	return false
}

// Java method: CompilerInput.getDynamicRequires()
func (this_ *CompilerInput) GetDynamicRequires() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.dynamicRequires.IsEmpty() {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.dynamicRequires))
	}
}

// Java method: CompilerInput.addDynamicRequire(java.lang.String require)
func (this_ *CompilerInput) AddDynamicRequire(require string) bool {
	_ = require
	if !(this_.dynamicRequires.Contains(j2gort.BoxString(require))) {
		this_.dynamicRequires.Add_1(j2gort.BoxString(require))
		return true
	}
	return false
}

// Java method: CompilerInput.getRequireDynamicImports()
func (this_ *CompilerInput) GetRequireDynamicImports() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.requireDynamicImports.IsEmpty() {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.requireDynamicImports))
	}
}

// Java method: CompilerInput.addRequireDynamicImports(java.lang.String require)
func (this_ *CompilerInput) AddRequireDynamicImports(require string) {
	_ = require
	if !(this_.requireDynamicImports.Contains(j2gort.BoxString(require))) {
		this_.requireDynamicImports.Add_1(j2gort.BoxString(require))
	}
}

// Java method: CompilerInput.setHasFullParseDependencyInfo(boolean hasFullParseDependencyInfo)
func (this_ *CompilerInput) SetHasFullParseDependencyInfo(hasFullParseDependencyInfo bool) {
	_ = hasFullParseDependencyInfo
	this_.hasFullParseDependencyInfo = hasFullParseDependencyInfo
}

// Java method: CompilerInput.getJsModuleType()
func (this_ *CompilerInput) GetJsModuleType() CompilerInput_ModuleType {
	return this_.jsModuleType
}

// Java method: CompilerInput.setJsModuleType(com.google.javascript.jscomp.CompilerInput.ModuleType moduleType)
func (this_ *CompilerInput) SetJsModuleType(moduleType CompilerInput_ModuleType) {
	_ = moduleType
	this_.jsModuleType = moduleType
}

// Java method: CompilerInput.addRequire(com.google.javascript.jscomp.deps.DependencyInfo.Require require)
func (this_ *CompilerInput) AddRequire(require *DependencyInfo_Require) {
	_ = require
	this_.extraRequires.Add_1(j2gort.PtrToAny(require))
}

// Java method: CompilerInput.getDependencyInfo()
func (this_ *CompilerInput) GetDependencyInfo() DependencyInfo {
	if this_.dependencyInfo == nil || this_.dependencyInfo.IsJavaNil_() {
		this_.dependencyInfo = this_.priv_CompilerInput_generateDependencyInfo_1626()
	}
	if !(this_.extraRequires.IsEmpty()) || !(this_.extraProvides.IsEmpty()) {
		this_.dependencyInfo = Upcast_DependencyInfo(SimpleDependencyInfo_S_Builder(this_.GetName(), this_.GetName()).SetProvides_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](CompilerInput_Concat(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.dependencyInfo.GetProvides()), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.extraProvides)))).SetRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](CompilerInput_Concat(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.dependencyInfo.GetRequires()), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.extraRequires)))).SetTypeRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.dependencyInfo.GetTypeRequires())).SetLoadFlags(j2gort.AsIface[j2gort.Ext_java_util_Map](this_.dependencyInfo.GetLoadFlags())).SetHasExternsAnnotation(this_.dependencyInfo.GetHasExternsAnnotation()).SetHasNoCompileAnnotation(this_.dependencyInfo.GetHasNoCompileAnnotation()).Build())
		this_.extraRequires.Clear()
		this_.extraProvides.Clear()
	}
	return this_.dependencyInfo
}

// Java method: CompilerInput.generateDependencyInfo()
func (this_ *CompilerInput) priv_CompilerInput_generateDependencyInfo_1626() DependencyInfo {
	j2gort.CheckNotNullIface(this_.compiler)
	j2gort.CheckNotNullIface(this_.compiler.GetErrorManager())
	if JsFileRegexParser_IsSupported() && this_.compiler.PreferRegexParser() {
		var didRet_2 bool
		var retVal_3 DependencyInfo
		func() {
			defer func() {
				if rec_4 := recover(); rec_4 != nil {
					jex_5 := javalang.RecoverException(rec_4)
					_ = jex_5
					if j2gort.InstanceOf[javaio.IIOException](jex_5) {
						e := j2gort.AsIface[javaio.IIOException](jex_5)
						_ = e
						this_.compiler.GetErrorManager().Report_2(StaticGet_CheckLevel_ERROR(), JSError_Make_2(StaticGet_AbstractCompiler_READ_ERROR(), []string{this_.GetName(), j2gort.ThrowableGetMessage(e)}))
						retVal_3 = Upcast_DependencyInfo(StaticGet_SimpleDependencyInfo_EMPTY())
						didRet_2 = true
						return
					} else {
						panic(rec_4)
					}
				}
			}()
			var info DependencyInfo = NewJsFileRegexParser(this_.compiler.GetErrorManager()).SetModuleLoader(this_.compiler.GetModuleLoader()).SetIncludeGoogBase(true).ParseFile(this_.GetName(), this_.GetName(), this_.GetCode())
			_ = info
			retVal_3 = NewLazyParsedDependencyInfo(info, this_, this_.compiler)
			didRet_2 = true
			return
		}()
		if didRet_2 {
			return retVal_3
		}
	} else {
		var val_6 CompilerInput_DepsFinder
		InitCompilerInput_DepsFinder(&val_6, this_.GetPath())
		var finder *CompilerInput_DepsFinder = &val_6
		_ = finder
		var root pkg_rhino.Node = this_.GetAstRoot(this_.compiler)
		_ = root
		if root == 0 {
			return Upcast_DependencyInfo(StaticGet_SimpleDependencyInfo_EMPTY())
		}
		finder.VisitTree(root)
		var info_2 *pkg_rhino.JSDocInfo = root.GetJSDocInfo()
		_ = info_2
		return Upcast_DependencyInfo(SimpleDependencyInfo_S_Builder("", "").SetProvides_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](finder.provides)).SetRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](finder.requires)).SetTypeRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](finder.typeRequires)).SetLoadFlags(finder.loadFlags).SetHasExternsAnnotation(((info_2 != nil) && info_2.IsExterns())).SetHasNoCompileAnnotation(((info_2 != nil) && info_2.IsNoCompile())).Build())
	}
	panic("unreachable")
}

// Java method: CompilerInput.getCode()
func (this_ *CompilerInput) GetCode() string {
	return this_.GetSourceFile().GetCode()
}

// Java method: CompilerInput.getChunk()
func (this_ *CompilerInput) GetChunk() *JSChunk {
	return this_.chunk
}

// Java method: CompilerInput.setChunk(com.google.javascript.jscomp.JSChunk chunk)
func (this_ *CompilerInput) SetChunk(chunk *JSChunk) {
	_ = chunk
	if !(((chunk == nil) || (this_.chunk == nil)) || (this_.chunk == chunk)) {
		j2gort.CheckArgument(false)
	}
	this_.chunk = chunk
}

// Java method: CompilerInput.overrideChunk(com.google.javascript.jscomp.JSChunk chunk)
func (this_ *CompilerInput) OverrideChunk(chunk *JSChunk) {
	_ = chunk
	this_.chunk = chunk
}

// Java method: CompilerInput.isExtern()
func (this_ *CompilerInput) IsExtern() bool {
	return this_.sourceFile.IsExtern()
}

// Java method: CompilerInput.setIsExtern()
func (this_ *CompilerInput) SetIsExtern() {
	this_.sourceFile.SetKind(pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN())
}

// Java method: CompilerInput.getLineOffset(int lineno)
func (this_ *CompilerInput) GetLineOffset(lineno int32) int32 {
	_ = lineno
	return this_.sourceFile.GetLineOffset(lineno)
}

// Java method: CompilerInput.toString()
func (this_ *CompilerInput) String() string {
	return this_.GetName()
}

// Java method: CompilerInput.isEs6Module()
func (this_ *CompilerInput) IsEs6Module() bool {
	return this_.GetDependencyInfo().IsEs6Module()
}

// Java method: CompilerInput.isGoogModule()
func (this_ *CompilerInput) IsGoogModule() bool {
	return this_.GetDependencyInfo().IsGoogModule()
}

// Java method: CompilerInput.getLoadFlags()
func (this_ *CompilerInput) GetLoadFlags() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.GetDependencyInfo().GetLoadFlags()
}

// Java static method: CompilerInput.concat(java.lang.Iterable<T> first, java.lang.Iterable<T> second)
func CompilerInput_Concat(first j2gort.Ext_java_lang_Iterable, second j2gort.Ext_java_lang_Iterable) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = first
	_ = second
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(first)).AddAll_1_Ext_java_lang_Iterable(second)).Build())
}

// Java method: CompilerInput.getPath()
func (this_ *CompilerInput) GetPath() *ModuleLoader_ModulePath {
	if this_.modulePath == nil {
		var moduleLoader *ModuleLoader = this_.compiler.GetModuleLoader()
		_ = moduleLoader
		this_.modulePath = moduleLoader.Resolve(this_.GetName())
	}
	return this_.modulePath
}

// Java method: CompilerInput.getFeatures(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *CompilerInput) GetFeatures(compiler AbstractCompiler) *pkg_parser.FeatureSet {
	_ = compiler
	var unused pkg_rhino.Node = this_.ast.GetAstRoot(compiler)
	_ = unused
	return this_.ast.features
}

// Java class: com.google.javascript.jscomp.CompilerInput.DepsFinder
type CompilerInput_DepsFinder struct {
	loadFlags    j2gort.Ext_java_util_Map  `gson:"loadFlags"`
	provides     j2gort.Ext_java_util_List `gson:"provides"`
	requires     j2gort.Ext_java_util_List `gson:"requires"`
	typeRequires j2gort.Ext_java_util_List `gson:"typeRequires"`
	modulePath   *ModuleLoader_ModulePath  `gson:"modulePath"`
	idHash_      atomic.Int32
}

func (this_ *CompilerInput_DepsFinder) BindSelf_(s any) {
	_ = s
}
func (this_ *CompilerInput_DepsFinder) IsJavaNil_() bool { return this_ == nil }
func (*CompilerInput_DepsFinder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerInput$DepsFinder", "DepsFinder"
}
func (this_ *CompilerInput_DepsFinder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CompilerInput_DepsFinder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerInput.DepsFinder", "com.google.javascript.jscomp.CompilerInput$DepsFinder", func(v_ any) bool {
		return j2gort.InstanceOf[*CompilerInput_DepsFinder](v_)
	})
}

func (this_ *CompilerInput_DepsFinder) initFields_CompilerInput_DepsFinder() {
	_ = this_
	this_.loadFlags = j2gort.NewExt_java_util_TreeMap()
	this_.provides = j2gort.NewExt_java_util_ArrayList()
	this_.requires = j2gort.NewExt_java_util_ArrayList()
	this_.typeRequires = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: DepsFinder(com.google.javascript.jscomp.deps.ModuleLoader.ModulePath modulePath)
func NewCompilerInput_DepsFinder(modulePath *ModuleLoader_ModulePath) *CompilerInput_DepsFinder {
	this_ := &CompilerInput_DepsFinder{}
	InitCompilerInput_DepsFinder(this_, modulePath)
	return this_
}

func InitCompilerInput_DepsFinder(this_ *CompilerInput_DepsFinder, modulePath *ModuleLoader_ModulePath) {
	_ = modulePath
	this_.initFields_CompilerInput_DepsFinder()
	this_.modulePath = modulePath
}

// Java method: DepsFinder.visitTree(com.google.javascript.rhino.Node n)
func (this_ *CompilerInput_DepsFinder) VisitTree(n pkg_rhino.Node) {
	_ = n
	this_.VisitSubtree(n, pkg_rhino.Node(0))
	if !(n.IsScript()) {
		j2gort.CheckArgument(false)
	}
	var features *pkg_parser.FeatureSet = j2gort.CastPtr[pkg_parser.FeatureSet](n.GetFeatureSet())
	_ = features
	if features != nil {
		var version string = features.Version()
		_ = version
		if !(j2gort.CheckNotNullString(version) == "es3") {
			j2gort.MapPutStr(this_.loadFlags, "lang", j2gort.BoxString(version))
		}
	}
}

// Java method: DepsFinder.visitSubtree(com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent)
func (this_ *CompilerInput_DepsFinder) VisitSubtree(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */ :
		if (n.HasTwoChildren() && n.GetFirstChild().IsGetProp()) && n.GetFirstFirstChild().MatchesName_1_string("goog") {
			if !(this_.requires.Contains(j2gort.PtrToAny(StaticGet_DependencyInfo_Require_BASE()))) {
				this_.requires.Add_1(j2gort.PtrToAny(StaticGet_DependencyInfo_Require_BASE()))
			}
			var callee pkg_rhino.Node = n.GetFirstChild()
			_ = callee
			var argument pkg_rhino.Node = n.GetLastChild()
			_ = argument
			switch j2gort.CheckNotNullString(callee.GetString()) {
			case "module":
				if parent.IsExprResult() && parent.GetParent().IsModuleBody() {
					j2gort.MapPutStr(this_.loadFlags, "module", "goog")
				}
				fallthrough
			case "provide":
				if !(argument.IsStringLit()) {
					return
				}
				this_.provides.Add_1(j2gort.BoxString(argument.GetString()))
				return
			case "require":
				if !(argument.IsStringLit()) {
					return
				}
				this_.requires.Add_1(j2gort.PtrToAny(DependencyInfo_Require_GoogRequireSymbol(argument.GetString())))
				return
			case "requireType", "forwardDeclare":
				if !(argument.IsStringLit()) {
					return
				}
				this_.typeRequires.Add_1(j2gort.BoxString(argument.GetString()))
				return
			case "loadModule":
				if argument.IsStringLit() {
					panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unsupported parse of goog.loadModule with string literal")))
				}
				n = argument.GetLastChild()
			case "declareModuleId":
				if !(argument.IsStringLit()) {
					return
				}
				this_.provides.Add_1(j2gort.BoxString(argument.GetString()))
			default:
				return
			}
		}
	case int32(123) /* Token.MODULE_BODY */ :
		if !(parent.GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE())) {
			this_.provides.Add_1(j2gort.BoxString(this_.modulePath.ToModuleName()))
			j2gort.MapPutStr(this_.loadFlags, "module", "es6")
		}
	case int32(116) /* Token.IMPORT */ :
		this_.VisitEs6ModuleName(n.GetLastChild(), n)
		return
	case int32(120) /* Token.EXPORT */ :
		if NodeUtil_IsExportFrom(n) {
			this_.VisitEs6ModuleName(n.GetLastChild(), n)
		}
		return
	case int32(93) /* Token.EXPR_RESULT */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */, int32(84) /* Token.VAR */, int32(90) /* Token.BLOCK */, int32(32) /* Token.NAME */, int32(104) /* Token.DESTRUCTURING_LHS */ :
	case int32(94) /* Token.SCRIPT */ :
		var jsdoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = jsdoc
		if (jsdoc != nil) && jsdoc.IsProvideGoog() {
			this_.provides.Add_1("goog")
		}
	default:
		return
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.VisitSubtree(child, n)
		}
	}
}

// Java method: DepsFinder.visitEs6ModuleName(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CompilerInput_DepsFinder) VisitEs6ModuleName(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	if !(n.IsStringLit()) {
		j2gort.CheckArgument(false)
	}
	if !(parent.IsExport() || parent.IsImport()) {
		j2gort.CheckArgument(false)
	}
	var moduleName string = n.GetString()
	_ = moduleName
	if javalang.StringStartsWith(j2gort.CheckNotNullString(moduleName), "goog:") {
		this_.requires.Add_1(j2gort.PtrToAny(DependencyInfo_Require_GoogRequireSymbol(j2gort.JStrSubFrom(moduleName, int32(5)))))
		return
	}
	var importedModule *ModuleLoader_ModulePath = this_.modulePath.ResolveJsModule(moduleName, this_.modulePath.String(), n.GetLineno(), n.GetCharno())
	_ = importedModule
	if importedModule == nil {
		importedModule = this_.modulePath.ResolveModuleAsPath(moduleName)
	}
	this_.requires.Add_1(j2gort.PtrToAny(DependencyInfo_Require_Es6Import(importedModule.ToModuleName(), n.GetString())))
}

// Java enum: com.google.javascript.jscomp.CompilerInput.ModuleType
type CompilerInput_ModuleType uint8

const (
	enum_CompilerInput_ModuleType_NONE            CompilerInput_ModuleType = 1
	enum_CompilerInput_ModuleType_GOOG            CompilerInput_ModuleType = 2
	enum_CompilerInput_ModuleType_ES6             CompilerInput_ModuleType = 3
	enum_CompilerInput_ModuleType_COMMONJS        CompilerInput_ModuleType = 4
	enum_CompilerInput_ModuleType_JSON            CompilerInput_ModuleType = 5
	enum_CompilerInput_ModuleType_IMPORTED_SCRIPT CompilerInput_ModuleType = 6
)

var enumBaseTable_CompilerInput_ModuleType = [...]j2gort.EnumBase{
	{Name_: "NONE", Ordinal_: 0},
	{Name_: "GOOG", Ordinal_: 1},
	{Name_: "ES6", Ordinal_: 2},
	{Name_: "COMMONJS", Ordinal_: 3},
	{Name_: "JSON", Ordinal_: 4},
	{Name_: "IMPORTED_SCRIPT", Ordinal_: 5},
}
var enumValues_CompilerInput_ModuleType = []CompilerInput_ModuleType{1, 2, 3, 4, 5, 6}

func StaticGet_CompilerInput_ModuleType_NONE() CompilerInput_ModuleType {
	return enum_CompilerInput_ModuleType_NONE
}
func StaticGet_CompilerInput_ModuleType_GOOG() CompilerInput_ModuleType {
	return enum_CompilerInput_ModuleType_GOOG
}
func StaticGet_CompilerInput_ModuleType_ES6() CompilerInput_ModuleType {
	return enum_CompilerInput_ModuleType_ES6
}
func StaticGet_CompilerInput_ModuleType_COMMONJS() CompilerInput_ModuleType {
	return enum_CompilerInput_ModuleType_COMMONJS
}
func StaticGet_CompilerInput_ModuleType_JSON() CompilerInput_ModuleType {
	return enum_CompilerInput_ModuleType_JSON
}
func StaticGet_CompilerInput_ModuleType_IMPORTED_SCRIPT() CompilerInput_ModuleType {
	return enum_CompilerInput_ModuleType_IMPORTED_SCRIPT
}

func CompilerInput_ModuleType_Values() []CompilerInput_ModuleType {
	return j2gort.CloneSlice(enumValues_CompilerInput_ModuleType)
}

func CompilerInput_ModuleType_ValueOf(name string) CompilerInput_ModuleType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilerInput_ModuleType {
		if enumBaseTable_CompilerInput_ModuleType[i_].Name_ == name {
			return CompilerInput_ModuleType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilerInput.ModuleType." + name))
}

func (this_ CompilerInput_ModuleType) IsJavaNil_() bool                              { return this_ == 0 }
func (CompilerInput_ModuleType) HasIdentityEquals_()                                 {}
func (CompilerInput_ModuleType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CompilerInput_ModuleType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilerInput_ModuleType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CompilerInput_ModuleType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerInput$ModuleType", "ModuleType"
}
func (this_ CompilerInput_ModuleType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilerInput_ModuleType[this_-1]
}
func (this_ CompilerInput_ModuleType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilerInput_ModuleType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilerInput_ModuleType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerInput_ModuleType[this_-1].Name_
}
func (this_ CompilerInput_ModuleType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilerInput_ModuleType[this_-1].Name_
}
func (this_ CompilerInput_ModuleType) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilerInput_ModuleType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilerInput_ModuleType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilerInput_ModuleType[this_-1].HashCode()
}
func (this_ CompilerInput_ModuleType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilerInput_ModuleType[this_-1].CompareTo(other_)
}
func (this_ CompilerInput_ModuleType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilerInput_ModuleType[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilerInput_ModuleType(1), CompilerInput_ModuleType(2), CompilerInput_ModuleType(3), CompilerInput_ModuleType(4), CompilerInput_ModuleType(5), CompilerInput_ModuleType(6)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerInput.ModuleType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilerInput$ModuleType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerInput.ModuleType", "com.google.javascript.jscomp.CompilerInput$ModuleType", func(v_ any) bool {
		return j2gort.InstanceOf[CompilerInput_ModuleType](v_)
	})
}

// Java class: com.google.javascript.jscomp.CompilerInput.JsAst
type CompilerInput_JsAst struct {
	outer_   *CompilerInput
	features *pkg_parser.FeatureSet `gson:"features"`
	root     pkg_rhino.Node         `gson:"root"`
	idHash_  atomic.Int32
}

func (this_ *CompilerInput_JsAst) BindSelf_(s any) {
	_ = s
}
func (this_ *CompilerInput_JsAst) IsJavaNil_() bool { return this_ == nil }
func (*CompilerInput_JsAst) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerInput$JsAst", "JsAst"
}
func (this_ *CompilerInput_JsAst) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CompilerInput_JsAst) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerInput.JsAst", "com.google.javascript.jscomp.CompilerInput$JsAst", func(v_ any) bool {
		return j2gort.InstanceOf[*CompilerInput_JsAst](v_)
	})
}
func (this_ *CompilerInput_JsAst) SetOuter_(o *CompilerInput) { this_.outer_ = o }
func (this_ *CompilerInput_JsAst) GetOuter_() *CompilerInput  { return this_.outer_ }

func (this_ *CompilerInput_JsAst) initFields_CompilerInput_JsAst() {
	_ = this_
}

// Java constructor: JsAst()
func NewCompilerInput_JsAst(outer_ *CompilerInput) *CompilerInput_JsAst {
	this_ := &CompilerInput_JsAst{}
	this_.outer_ = outer_
	InitCompilerInput_JsAst(this_)
	return this_
}

func InitCompilerInput_JsAst(this_ *CompilerInput_JsAst) {
	this_.initFields_CompilerInput_JsAst()
}

// Java method: JsAst.getAstRoot(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *CompilerInput_JsAst) GetAstRoot(compiler AbstractCompiler) pkg_rhino.Node {
	_ = compiler
	if this_.priv_JsAst_isParsed_1645() {
		return this_.root
	}
	var astRootSource j2gort.Ext_java_util_function_Supplier = compiler.GetTypedAstDeserializer(this_.outer_.sourceFile)
	_ = astRootSource
	if astRootSource != nil {
		this_.root = j2gort.CastNode[pkg_rhino.Node](astRootSource.Get())
		this_.features = j2gort.CastPtr[pkg_parser.FeatureSet](this_.root.GetFeatureSet())
	} else {
		this_.priv_JsAst_parse_1646(compiler)
	}
	if !(j2gort.IfaceEquals[pkg_rhino.StaticSourceFile](this_.root.GetStaticSourceFile(), this_.outer_.sourceFile)) {
		j2gort.CheckState(false)
	}
	this_.root.SetInputId(this_.outer_.id)
	this_.outer_.sourceFile.ClearCachedSource()
	return this_.root
}

// Java method: JsAst.clearAst()
func (this_ *CompilerInput_JsAst) ClearAst() {
	this_.root = pkg_rhino.Node(0)
	this_.outer_.sourceFile.ClearCachedSource()
}

// Java method: JsAst.isParsed()
func (this_ *CompilerInput_JsAst) priv_JsAst_isParsed_1645() bool {
	return (this_.root != 0)
}

// Java method: JsAst.parse(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *CompilerInput_JsAst) priv_JsAst_parse_1646(compiler AbstractCompiler) {
	_ = compiler
	func() {
		defer func() {
			if rec_7 := recover(); rec_7 != nil {
				jex_8 := javalang.RecoverException(rec_7)
				_ = jex_8
				if j2gort.InstanceOf[javaio.IIOException](jex_8) {
					e := j2gort.AsIface[javaio.IIOException](jex_8)
					_ = e
					compiler.Report_1(JSError_Make_2(StaticGet_AbstractCompiler_READ_ERROR(), []string{this_.outer_.sourceFile.GetName(), j2gort.ThrowableGetMessage(e)}))
				} else {
					panic(rec_7)
				}
			}
		}()
		var result *pkg_parsing.ParserRunner_ParseResult = pkg_parsing.ParserRunner_Parse(pkg_rhino.ToIface_StaticSourceFile(this_.outer_.sourceFile), this_.outer_.sourceFile.GetCode(), compiler.GetParserConfig(func() AbstractCompiler_ConfigContext {
			if this_.outer_.sourceFile.IsExtern() {
				return StaticGet_AbstractCompiler_ConfigContext_EXTERNS()
			}
			return StaticGet_AbstractCompiler_ConfigContext_DEFAULT()
		}()), compiler.GetDefaultErrorReporter())
		_ = result
		this_.root = result.FieldGet_ast()
		this_.features = result.FieldGet_features()
		if compiler.GetOptions().PreservesDetailedSourceInfo() {
			compiler.AddComments(this_.outer_.sourceFile.GetName(), result.FieldGet_comments())
		}
		if (result.FieldGet_sourceMapURL() != j2gort.NullString) && compiler.GetOptions().GetResolveSourceMapAnnotations() {
			var parseInline bool = compiler.GetOptions().GetParseInlineSourceMaps()
			_ = parseInline
			var sourceMapSourceFile *SourceFile = SourceMapResolver_ExtractSourceMap(this_.outer_.sourceFile, result.FieldGet_sourceMapURL(), parseInline)
			_ = sourceMapSourceFile
			if sourceMapSourceFile != nil {
				compiler.AddInputSourceMap(this_.outer_.sourceFile.GetName(), NewSourceMapInput(sourceMapSourceFile))
			}
		}
	}()
	if this_.root == 0 {
		this_.root = pkg_rhino.IR_Script_0()
	}
	this_.root.SetStaticSourceFile(pkg_rhino.ToIface_StaticSourceFile(this_.outer_.sourceFile))
}

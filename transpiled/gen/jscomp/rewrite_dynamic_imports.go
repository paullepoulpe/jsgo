/*
 * Copyright 2021 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RewriteDynamicImports.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.RewriteDynamicImports
type RewriteDynamicImports struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler                         AbstractCompiler `gson:"compiler"`
	astFactory                       *AstFactory      `gson:"astFactory"`
	alias                            string           `gson:"alias"`
	requiresAliasing                 bool             `gson:"requiresAliasing"`
	shouldWrapDynamicImportCallbacks bool             `gson:"shouldWrapDynamicImportCallbacks"`
	dynamicImportsRemoved            bool             `gson:"dynamicImportsRemoved"`
	wrappedDynamicImportCallback     bool             `gson:"wrappedDynamicImportCallback"`
}

func (this_ *RewriteDynamicImports) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteDynamicImports) IsJavaNil_() bool { return this_ == nil }
func (*RewriteDynamicImports) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteDynamicImports", "RewriteDynamicImports"
}
func (this_ *RewriteDynamicImports) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RewriteDynamicImports) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteDynamicImports", "com.google.javascript.jscomp.RewriteDynamicImports", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteDynamicImports](v_)
	})
}
func (*RewriteDynamicImports) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*RewriteDynamicImports) Implements_jscomp_CompilerPass()                            {}
func (*RewriteDynamicImports) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_RewriteDynamicImports_DYNAMIC_IMPORT_ALIASING_REQUIRED *DiagnosticType

func StaticGet_RewriteDynamicImports_DYNAMIC_IMPORT_ALIASING_REQUIRED() *DiagnosticType {
	clinit_RewriteDynamicImports()
	return static_RewriteDynamicImports_DYNAMIC_IMPORT_ALIASING_REQUIRED
}
func StaticSet_RewriteDynamicImports_DYNAMIC_IMPORT_ALIASING_REQUIRED(v *DiagnosticType) {
	clinit_RewriteDynamicImports()
	static_RewriteDynamicImports_DYNAMIC_IMPORT_ALIASING_REQUIRED = v
}

var static_RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS *DiagnosticType

func StaticGet_RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS() *DiagnosticType {
	clinit_RewriteDynamicImports()
	return static_RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS
}
func StaticSet_RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS(v *DiagnosticType) {
	clinit_RewriteDynamicImports()
	static_RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS = v
}

var clinitDone_RewriteDynamicImports uint32
var clinitState_RewriteDynamicImports j2gort.ClinitState

func clinit_RewriteDynamicImports() {
	if atomic.LoadUint32(&clinitDone_RewriteDynamicImports) != 0 {
		return
	}
	clinitSlow_RewriteDynamicImports()
}

func clinitSlow_RewriteDynamicImports() {
	if !j2gort.EnterClinit(&clinitState_RewriteDynamicImports, &clinitDone_RewriteDynamicImports) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RewriteDynamicImports, &clinitDone_RewriteDynamicImports)
	static_RewriteDynamicImports_DYNAMIC_IMPORT_ALIASING_REQUIRED = DiagnosticType_Warning("JSC_DYNAMIC_IMPORT_ALIASING_REQUIRED", "Dynamic import expressions should be aliased for language level. Use the --dynamic_import_alias flag.")
	static_RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS = DiagnosticType_Error("JSC_DYNAMIC_IMPORT_INVALID_ALIAS", "Dynamic import alias is not a valid name")
}

func (this_ *RewriteDynamicImports) initFields_RewriteDynamicImports() {
	_ = this_
	this_.alias = j2gort.NullString

}

// Java constructor: RewriteDynamicImports(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.@org.jspecify.annotations.Nullable String alias, com.google.javascript.jscomp.CompilerOptions.ChunkOutputType chunkOutputType)
func NewRewriteDynamicImports(compiler AbstractCompiler, alias string, chunkOutputType CompilerOptions_ChunkOutputType) *RewriteDynamicImports {
	clinit_RewriteDynamicImports()
	this_ := &RewriteDynamicImports{}
	InitRewriteDynamicImports(this_, compiler, alias, chunkOutputType)
	return this_
}

func InitRewriteDynamicImports(this_ *RewriteDynamicImports, compiler AbstractCompiler, alias string, chunkOutputType CompilerOptions_ChunkOutputType) {
	_ = compiler
	_ = alias
	_ = chunkOutputType
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RewriteDynamicImports()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.alias = alias
	this_.requiresAliasing = !(compiler.GetOptions().GetOutputFeatureSet().Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT()))
	this_.shouldWrapDynamicImportCallbacks = (chunkOutputType == StaticGet_CompilerOptions_ChunkOutputType_ES_MODULES())
}

// Java method: RewriteDynamicImports.aliasIsValid()
func (this_ *RewriteDynamicImports) priv_RewriteDynamicImports_aliasIsValid_4979() bool {
	return ((this_.alias != j2gort.NullString) && ((j2gort.CheckNotNullString(this_.alias) == "import") || NodeUtil_IsValidQualifiedName_2_FeatureSet_string(this_.compiler.GetAllowableFeatures(), this_.alias)))
}

// Java method: RewriteDynamicImports.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewriteDynamicImports) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.dynamicImportsRemoved = false
	if !(externs.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(externs))
	}
	if !(root.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(root))
	}
	NodeTraversal_Traverse(this_.compiler, root, this_)
	if this_.wrappedDynamicImportCallback {
		this_.priv_RewriteDynamicImports_injectWrappingFunctionExtern_4985()
	}
	if this_.dynamicImportsRemoved {
		this_.compiler.MarkFeatureNotAllowed(pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT())
		NodeUtil_RemoveFeatureFromAllScripts(root, pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT(), this_.compiler)
		if this_.requiresAliasing && this_.priv_RewriteDynamicImports_aliasIsValid_4979() {
			NodeTraversal_Traverse(this_.compiler, externs, NewRewriteDynamicImports_AliasInjectingTraversal(this_))
		}
	}
}

// Java method: RewriteDynamicImports.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteDynamicImports) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.RawToken_() != uint8(124) {
		return
	}
	var moduleMap *ModuleMap = this_.compiler.GetModuleMap()
	_ = moduleMap
	var importSpecifier pkg_rhino.Node = n.GetFirstChild()
	_ = importSpecifier
	if importSpecifier.IsStringLit() && (moduleMap != nil) {
		var targetPath *ModuleLoader_ModulePath = t.GetInput().GetPath().ResolveJsModule(importSpecifier.GetString(), n.GetSourceFileName(), n.GetLineno(), n.GetCharno())
		_ = targetPath
		var module *Module
		_ = module
		if targetPath == nil {
			module = nil
		} else {
			module = this_.compiler.GetModuleMap().GetModule_1_ModuleLoader_ModulePath(targetPath)
		}
		var targetModuleVarName string
		_ = targetModuleVarName
		if module == nil {
			targetModuleVarName = j2gort.NullString
		} else {
			targetModuleVarName = ModuleRenaming_GlobalizedModuleName_Create_3(module.Metadata(), j2gort.NullString, nil).AliasName().Join_0()
		}
		var targetModuleNS *Var
		_ = targetModuleNS
		if targetModuleVarName == j2gort.NullString {
			targetModuleNS = nil
		} else {
			targetModuleNS = j2gort.CastPtr[Var](t.GetScope().GetVar(targetModuleVarName))
		}
		if targetModuleNS != nil {
			var targetChunk *JSChunk = targetModuleNS.GetInput().GetChunk()
			_ = targetChunk
			if t.GetChunk() == targetChunk {
				this_.priv_RewriteDynamicImports_replaceDynamicImportWithPromise_4980(t, n, targetModuleNS)
				return
			} else {
				this_.priv_RewriteDynamicImports_retargetImportSpecifier_4981(t, n, targetChunk)
				if NodeUtil_IsExpressionResultUsed(n) {
					this_.priv_RewriteDynamicImports_addChainedThen_4982(t, n, targetModuleNS)
				}
			}
		}
	}
	if this_.priv_RewriteDynamicImports_aliasIsValid_4979() {
		this_.priv_RewriteDynamicImports_aliasDynamicImport_4983(t, n)
	} else if this_.alias != j2gort.NullString {
		t.Report_3(n, StaticGet_RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS(), []string{})
	} else if this_.requiresAliasing {
		t.Report_3(n, StaticGet_RewriteDynamicImports_DYNAMIC_IMPORT_ALIASING_REQUIRED(), []string{})
	}
}

// Java method: RewriteDynamicImports.replaceDynamicImportWithPromise(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node dynamicImport, com.google.javascript.jscomp.Var targetModuleNS)
func (this_ *RewriteDynamicImports) priv_RewriteDynamicImports_replaceDynamicImportWithPromise_4980(t *NodeTraversal, dynamicImport pkg_rhino.Node, targetModuleNS *Var) {
	_ = t
	_ = dynamicImport
	_ = targetModuleNS
	var registry *pkg_rhino.JSTypeRegistry = this_.compiler.GetTypeRegistry()
	_ = registry
	var promiseDotResolve pkg_rhino.Node = this_.astFactory.CreateQName_2_StaticScope_string(this_.compiler.GetTranspilationNamespace(), "Promise.resolve")
	_ = promiseDotResolve
	var promiseResolveCall pkg_rhino.Node = this_.astFactory.CreateCall(promiseDotResolve, AstFactory_S_Type_1_Node(dynamicImport), []pkg_rhino.Node{})
	_ = promiseResolveCall
	var isExpressionUsed bool = NodeUtil_IsExpressionResultUsed(dynamicImport)
	_ = isExpressionUsed
	if isExpressionUsed {
		var moduleNamespaceType pkg_rhino.JSType = nil
		_ = moduleNamespaceType
		if !(j2gort.IfaceEquals[pkg_rhino.JSType](dynamicImport.GetJSType(), nil)) && dynamicImport.GetJSType().IsTemplatizedType() {
			var templatizedType *pkg_rhino.TemplatizedType = j2gort.CastPtr[pkg_rhino.TemplatizedType](dynamicImport.GetJSType())
			_ = templatizedType
			moduleNamespaceType = j2gort.CheckNotNullSlice(j2gort.CheckNotNullPtr(templatizedType).Slice_GetTemplateTypes())[int32(0)]
		} else {
			moduleNamespaceType = registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}
		promiseDotResolve.SetJSType(pkg_rhino.ToIface_JSType(registry.CreateFunctionType_2_JSType_JSType(pkg_rhino.ToIface_JSType(registry.CreateTemplatizedType_2_ObjectType_JSType(registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE()), []pkg_rhino.JSType{moduleNamespaceType})), []pkg_rhino.JSType{registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ALL_TYPE())})))
		promiseResolveCall.AddChildToBack(this_.priv_RewriteDynamicImports_createModuleNamespaceNode_4984(targetModuleNS))
		promiseResolveCall.CopyTypeFrom(dynamicImport)
	} else {
		promiseDotResolve.SetJSType(pkg_rhino.ToIface_JSType(registry.CreateFunctionType_2_JSType_JSType(pkg_rhino.ToIface_JSType(registry.CreateTemplatizedType_2_ObjectType_JSType(registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE()), []pkg_rhino.JSType{registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())})), []pkg_rhino.JSType{registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ALL_TYPE())})))
		promiseResolveCall.SetJSType(pkg_rhino.ToIface_JSType(registry.CreateTemplatizedType_2_ObjectType_JSType(registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE()), []pkg_rhino.JSType{registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())})))
	}
	promiseResolveCall.SrcrefTree(dynamicImport)
	var parent pkg_rhino.Node = dynamicImport.GetParent()
	_ = parent
	dynamicImport.ReplaceWith(promiseResolveCall)
	t.ReportCodeChange_1(parent)
	this_.dynamicImportsRemoved = true
}

// Java static method: RewriteDynamicImports.getChunkFileName(com.google.javascript.jscomp.JSChunk chunk)
func RewriteDynamicImports_GetChunkFileName(chunk *JSChunk) string {
	clinit_RewriteDynamicImports()
	_ = chunk
	return (j2gort.StrValString(chunk.GetName()) + ".js")
}

// Java method: RewriteDynamicImports.retargetImportSpecifier(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node dynamicImport, com.google.javascript.jscomp.JSChunk targetChunk)
func (this_ *RewriteDynamicImports) priv_RewriteDynamicImports_retargetImportSpecifier_4981(t *NodeTraversal, dynamicImport pkg_rhino.Node, targetChunk *JSChunk) {
	_ = t
	_ = dynamicImport
	_ = targetChunk
	var retargetedSpecifier string = j2gort.NullString
	_ = retargetedSpecifier
	var importingChunkFilename string = RewriteDynamicImports_GetChunkFileName(t.GetInput().GetChunk())
	_ = importingChunkFilename
	var targetChunkFilename string = RewriteDynamicImports_GetChunkFileName(targetChunk)
	_ = targetChunkFilename
	var didRet_1 bool
	func() {
		defer func() {
			if rec_2 := recover(); rec_2 != nil {
				jex_3 := javalang.RecoverException(rec_2)
				_ = jex_3
				if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_3) {
					e := j2gort.AsIface[javalang.IIllegalArgumentException](jex_3)
					_ = e
					this_.compiler.Report_1(JSError_Make_3(dynamicImport, StaticGet_ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH(), []string{importingChunkFilename, targetChunkFilename}))
					didRet_1 = true
					return
				} else {
					panic(rec_2)
				}
			}
		}()
		retargetedSpecifier = ModuleLoader_RelativePathFrom(importingChunkFilename, targetChunkFilename)
	}()
	if didRet_1 {
		return
	}
	var originalSpecifierNode pkg_rhino.Node = dynamicImport.GetFirstChild()
	_ = originalSpecifierNode
	var newSpecifier pkg_rhino.Node = this_.astFactory.CreateString(retargetedSpecifier).Srcref(originalSpecifierNode)
	_ = newSpecifier
	originalSpecifierNode.ReplaceWith(newSpecifier)
	t.ReportCodeChange_1(newSpecifier)
}

// Java method: RewriteDynamicImports.addChainedThen(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node dynamicImport, com.google.javascript.jscomp.Var targetModuleNs)
func (this_ *RewriteDynamicImports) priv_RewriteDynamicImports_addChainedThen_4982(t *NodeTraversal, dynamicImport pkg_rhino.Node, targetModuleNs *Var) {
	_ = t
	_ = dynamicImport
	_ = targetModuleNs
	var registry *pkg_rhino.JSTypeRegistry = this_.compiler.GetTypeRegistry()
	_ = registry
	var importParent pkg_rhino.Node = dynamicImport.GetParent()
	_ = importParent
	var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
	_ = placeholder
	dynamicImport.ReplaceWith(placeholder)
	var moduleNamespaceNode pkg_rhino.Node = this_.priv_RewriteDynamicImports_createModuleNamespaceNode_4984(targetModuleNs)
	_ = moduleNamespaceNode
	var callbackFn pkg_rhino.Node = this_.astFactory.CreateZeroArgArrowFunctionForExpression(moduleNamespaceNode)
	_ = callbackFn
	callbackFn.SetJSType(pkg_rhino.ToIface_JSType(registry.CreateFunctionType_2_JSType_JSType(pkg_rhino.ToIface_JSType(registry.CreateTemplatizedType_2_ObjectType_JSType(registry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE()), []pkg_rhino.JSType{registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())})), []pkg_rhino.JSType{pkg_rhino.ToIface_JSType(registry.CreateFunctionType_2_JSType_JSType(registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ALL_TYPE()), []pkg_rhino.JSType{}))})))
	var thenArgument pkg_rhino.Node = callbackFn
	_ = thenArgument
	if this_.shouldWrapDynamicImportCallbacks {
		var wrappingFunction pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type("jscomp$DynamicImportCallback", AstFactory_S_Type_1_JSType(pkg_rhino.ToIface_JSType(registry.CreateFunctionType_2_JSType_JSType(callbackFn.GetJSType(), []pkg_rhino.JSType{callbackFn.GetJSType()}))))
		_ = wrappingFunction
		thenArgument = this_.astFactory.CreateCall(wrappingFunction, AstFactory_S_Type_1_Node(callbackFn), []pkg_rhino.Node{callbackFn})
		this_.wrappedDynamicImportCallback = true
	}
	var importThenCall pkg_rhino.Node = this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(dynamicImport, "then"), AstFactory_S_Type_1_Node(dynamicImport), []pkg_rhino.Node{thenArgument})
	_ = importThenCall
	importThenCall.SrcrefTreeIfMissing(dynamicImport)
	if !(j2gort.IfaceEquals[pkg_rhino.JSType](dynamicImport.GetJSType(), nil)) {
		importThenCall.CopyTypeFrom(dynamicImport)
	}
	placeholder.ReplaceWith(importThenCall)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS(), this_.compiler)
	this_.compiler.ReportChangeToChangeScope(callbackFn)
	this_.compiler.ReportChangeToEnclosingScope(importParent)
}

// Java method: RewriteDynamicImports.aliasDynamicImport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node dynamicImport)
func (this_ *RewriteDynamicImports) priv_RewriteDynamicImports_aliasDynamicImport_4983(t *NodeTraversal, dynamicImport pkg_rhino.Node) {
	_ = t
	_ = dynamicImport
	j2gort.CheckNotNullString(this_.alias)
	var aliasNode pkg_rhino.Node = this_.astFactory.CreateQNameWithUnknownType_1_string(this_.alias)
	_ = aliasNode
	aliasNode.SetOriginalName("import")
	var moduleSpecifier pkg_rhino.Node = dynamicImport.RemoveFirstChild()
	_ = moduleSpecifier
	var importAliasCall pkg_rhino.Node = this_.astFactory.CreateCall(aliasNode, AstFactory_S_Type_1_Node(dynamicImport), []pkg_rhino.Node{moduleSpecifier}).SrcrefTreeIfMissing(dynamicImport)
	_ = importAliasCall
	if !(j2gort.IfaceEquals[pkg_rhino.JSType](dynamicImport.GetJSType(), nil)) {
		importAliasCall.CopyTypeFrom(dynamicImport)
	}
	dynamicImport.ReplaceWith(importAliasCall)
	t.ReportCodeChange_1(importAliasCall)
	this_.dynamicImportsRemoved = true
}

// Java method: RewriteDynamicImports.createModuleNamespaceNode(com.google.javascript.jscomp.Var moduleVar)
func (this_ *RewriteDynamicImports) priv_RewriteDynamicImports_createModuleNamespaceNode_4984(moduleVar *Var) pkg_rhino.Node {
	_ = moduleVar
	var moduleVarNode pkg_rhino.Node = moduleVar.GetNode()
	_ = moduleVarNode
	var moduleNamespace pkg_rhino.Node = moduleVarNode.CloneNode_0()
	_ = moduleNamespace
	moduleNamespace.CopyTypeFrom(moduleVarNode)
	return moduleNamespace
}

// Java method: RewriteDynamicImports.injectWrappingFunctionExtern()
func (this_ *RewriteDynamicImports) priv_RewriteDynamicImports_injectWrappingFunctionExtern_4985() {
	var registry *pkg_rhino.JSTypeRegistry = this_.compiler.GetTypeRegistry()
	_ = registry
	var templateT *pkg_rhino.TemplateType = registry.CreateTemplateType_1("T")
	_ = templateT
	var wrappingFunctionDefinition pkg_rhino.Node = this_.astFactory.CreateFunction("jscomp$DynamicImportCallback", this_.astFactory.CreateParamList([]string{"importCallback"}), this_.astFactory.CreateBlock([]pkg_rhino.Node{}), AstFactory_S_Type_1_JSType(pkg_rhino.ToIface_JSType(registry.CreateFunctionType_2_JSType_JSType(pkg_rhino.ToIface_JSType(templateT), []pkg_rhino.JSType{pkg_rhino.ToIface_JSType(templateT)}))))
	_ = wrappingFunctionDefinition
	var externsRoot pkg_rhino.Node = this_.compiler.GetSynthesizedExternsInput().GetAstRoot(this_.compiler)
	_ = externsRoot
	wrappingFunctionDefinition.SrcrefTree(externsRoot)
	externsRoot.AddChildToBack(wrappingFunctionDefinition)
	this_.compiler.ReportChangeToEnclosingScope(wrappingFunctionDefinition)
}

// Java class: com.google.javascript.jscomp.RewriteDynamicImports.AliasInjectingTraversal
type RewriteDynamicImports_AliasInjectingTraversal struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer_ *RewriteDynamicImports
}

func (this_ *RewriteDynamicImports_AliasInjectingTraversal) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteDynamicImports_AliasInjectingTraversal) IsJavaNil_() bool { return this_ == nil }
func (*RewriteDynamicImports_AliasInjectingTraversal) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteDynamicImports$AliasInjectingTraversal", "AliasInjectingTraversal"
}
func (this_ *RewriteDynamicImports_AliasInjectingTraversal) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*RewriteDynamicImports_AliasInjectingTraversal) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteDynamicImports.AliasInjectingTraversal", "com.google.javascript.jscomp.RewriteDynamicImports$AliasInjectingTraversal", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteDynamicImports_AliasInjectingTraversal](v_)
	})
}
func (this_ *RewriteDynamicImports_AliasInjectingTraversal) SetOuter_(o *RewriteDynamicImports) {
	this_.outer_ = o
}
func (this_ *RewriteDynamicImports_AliasInjectingTraversal) GetOuter_() *RewriteDynamicImports {
	return this_.outer_
}
func (*RewriteDynamicImports_AliasInjectingTraversal) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*RewriteDynamicImports_AliasInjectingTraversal) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *RewriteDynamicImports_AliasInjectingTraversal) initFields_RewriteDynamicImports_AliasInjectingTraversal() {
	_ = this_
}

// Java constructor: AliasInjectingTraversal()
func NewRewriteDynamicImports_AliasInjectingTraversal(outer_ *RewriteDynamicImports) *RewriteDynamicImports_AliasInjectingTraversal {
	this_ := &RewriteDynamicImports_AliasInjectingTraversal{}
	this_.outer_ = outer_
	InitRewriteDynamicImports_AliasInjectingTraversal(this_)
	return this_
}

func InitRewriteDynamicImports_AliasInjectingTraversal(this_ *RewriteDynamicImports_AliasInjectingTraversal) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_RewriteDynamicImports_AliasInjectingTraversal()
}

// Java method: AliasInjectingTraversal.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteDynamicImports_AliasInjectingTraversal) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if !(this_.outer_.priv_RewriteDynamicImports_aliasIsValid_4979()) {
		j2gort.CheckState(false)
	}
	if NodeUtil_IsValidSimpleName(this_.outer_.alias) && (j2gort.CastPtr[Var](t.GetScope().GetVar(this_.outer_.alias)) != nil) {
		return false
	}
	var registry *pkg_rhino.JSTypeRegistry = this_.outer_.compiler.GetTypeRegistry()
	_ = registry
	var unknownType pkg_rhino.JSType = registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = unknownType
	var promiseType pkg_rhino.JSType = registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE())
	_ = promiseType
	var promiseTemplatizedType *pkg_rhino.TemplatizedType = registry.CreateTemplatizedType_2_ObjectType_Ext_com_google_common_collect_ImmutableList(promiseType.ToObjectType(), j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(unknownType))
	_ = promiseTemplatizedType
	var stringType pkg_rhino.JSType = registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
	_ = stringType
	var externsRoot pkg_rhino.Node = this_.outer_.compiler.GetSynthesizedExternsInput().GetAstRoot(this_.outer_.compiler)
	_ = externsRoot
	if (j2gort.CheckNotNullString(this_.outer_.alias) == "import") || NodeUtil_IsValidSimpleName(this_.outer_.alias) {
		var aliasNode pkg_rhino.Node = this_.outer_.astFactory.CreateFunction(this_.outer_.alias, this_.outer_.astFactory.CreateParamList([]string{"specifier"}), this_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{}), AstFactory_S_Type_1_JSType(pkg_rhino.ToIface_JSType(registry.CreateFunctionType_2_JSType_JSType(pkg_rhino.ToIface_JSType(promiseTemplatizedType), []pkg_rhino.JSType{stringType}))))
		_ = aliasNode
		aliasNode.SrcrefTree(externsRoot)
		externsRoot.AddChildToBack(aliasNode)
	} else {
		var aliasRootName string = NodeUtil_GetRootOfQualifiedName_1_string(this_.outer_.alias)
		_ = aliasRootName
		var aliasVar *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(aliasRootName))
		_ = aliasVar
		if aliasVar != nil {
			return false
		}
		var aliasRootNode pkg_rhino.Node = this_.outer_.astFactory.CreateSingleVarNameDeclaration_2(aliasRootName, this_.outer_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{}))
		_ = aliasRootNode
		aliasRootNode.SrcrefTree(externsRoot)
		externsRoot.AddChildToBack(aliasRootNode)
		var aliasNameParts j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.SplitterSplitStr(j2gort.Ext_com_google_common_base_Splitter_On_1_string("."), this_.outer_.alias).Iterator())
		_ = aliasNameParts
		var qName pkg_rhino.Node = aliasRootNode.GetFirstChild().CloneNode_0()
		_ = qName
		aliasNameParts.Next()
		for aliasNameParts.HasNext() {
			qName = this_.outer_.astFactory.CreateGetPropWithUnknownType(qName.CloneTree_0(), j2gort.CastString(aliasNameParts.Next()))
			var assignedValue pkg_rhino.Node = pkg_rhino.Node(0)
			_ = assignedValue
			if !(aliasNameParts.HasNext()) {
				assignedValue = this_.outer_.astFactory.CreateFunction("", this_.outer_.astFactory.CreateParamList([]string{"specifier"}), this_.outer_.astFactory.CreateBlock([]pkg_rhino.Node{}), AstFactory_S_Type_1_JSType(pkg_rhino.ToIface_JSType(registry.CreateFunctionType_2_JSType_JSType(pkg_rhino.ToIface_JSType(promiseTemplatizedType), []pkg_rhino.JSType{stringType}))))
			} else {
				assignedValue = this_.outer_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})
			}
			var expr pkg_rhino.Node = this_.outer_.astFactory.ExprResult(this_.outer_.astFactory.CreateAssign_2_Node_Node(qName, assignedValue))
			_ = expr
			expr.SrcrefTree(externsRoot)
			externsRoot.AddChildToBack(expr)
		}
	}
	this_.outer_.compiler.ReportChangeToChangeScope(externsRoot)
	return false
}
func (*RewriteDynamicImports) AlwaysShouldTraverse_()                              {}
func (*RewriteDynamicImports) CanReleaseSyntacticScopes_()                         {}
func (*RewriteDynamicImports_AliasInjectingTraversal) NoOpVisit_()                 {}
func (*RewriteDynamicImports_AliasInjectingTraversal) CanReleaseSyntacticScopes_() {}

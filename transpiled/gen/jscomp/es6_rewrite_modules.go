/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Es6RewriteModules.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.Es6RewriteModules
type Es6RewriteModules struct {
	compiler                AbstractCompiler                                                 `gson:"compiler"`
	astFactory              *AstFactory                                                      `gson:"astFactory"`
	unknownType             pkg_rhino.JSType                                                 `gson:"unknownType"`
	preprocessorSymbolTable *PreprocessorSymbolTable                                         `gson:"preprocessorSymbolTable"`
	namesToInlineByAlias    *j2gort.OrderedMap[string, *ModuleRenaming_GlobalizedModuleName] `gson:"namesToInlineByAlias"`
	typedefs                *j2gort.OrderedSet[string]                                       `gson:"typedefs"`
	moduleMetadataMap       *ModuleMetadataMap                                               `gson:"moduleMetadataMap"`
	moduleMap               *ModuleMap                                                       `gson:"moduleMap"`
	globalTypedScope        *TypedScope                                                      `gson:"globalTypedScope"`
	idHash_                 atomic.Int32
	chunkOutputType         CompilerOptions_ChunkOutputType `gson:"chunkOutputType"`
}

func (this_ *Es6RewriteModules) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteModules) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteModules) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteModules", "Es6RewriteModules"
}
func (this_ *Es6RewriteModules) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteModules) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteModules", "com.google.javascript.jscomp.Es6RewriteModules", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteModules](v_)
	})
}
func (*Es6RewriteModules) Implements_jscomp_CompilerPass()           {}
func (*Es6RewriteModules) Implements_jscomp_NodeTraversal_Callback() {}

var static_Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST *DiagnosticType

func StaticGet_Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST() *DiagnosticType {
	clinit_Es6RewriteModules()
	return static_Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST
}
func StaticSet_Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST(v *DiagnosticType) {
	clinit_Es6RewriteModules()
	static_Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST = v
}

var static_Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST *DiagnosticType

func StaticGet_Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST() *DiagnosticType {
	clinit_Es6RewriteModules()
	return static_Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST
}
func StaticSet_Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST(v *DiagnosticType) {
	clinit_Es6RewriteModules()
	static_Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST = v
}

var static_Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST *DiagnosticType

func StaticGet_Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST() *DiagnosticType {
	clinit_Es6RewriteModules()
	return static_Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST
}
func StaticSet_Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST(v *DiagnosticType) {
	clinit_Es6RewriteModules()
	static_Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST = v
}

var static_Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE *DiagnosticType

func StaticGet_Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE() *DiagnosticType {
	clinit_Es6RewriteModules()
	return static_Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE
}
func StaticSet_Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE(v *DiagnosticType) {
	clinit_Es6RewriteModules()
	static_Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE = v
}

var static_Es6RewriteModules_DOT_SPLITTER *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_Es6RewriteModules_DOT_SPLITTER() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_Es6RewriteModules()
	return static_Es6RewriteModules_DOT_SPLITTER
}
func StaticSet_Es6RewriteModules_DOT_SPLITTER(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_Es6RewriteModules()
	static_Es6RewriteModules_DOT_SPLITTER = v
}

var static_Es6RewriteModules_GOOG_REQUIRE pkg_rhino.QualifiedName

func StaticGet_Es6RewriteModules_GOOG_REQUIRE() pkg_rhino.QualifiedName {
	clinit_Es6RewriteModules()
	return static_Es6RewriteModules_GOOG_REQUIRE
}
func StaticSet_Es6RewriteModules_GOOG_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_Es6RewriteModules()
	static_Es6RewriteModules_GOOG_REQUIRE = v
}

var static_Es6RewriteModules_GOOG_REQUIRETYPE pkg_rhino.QualifiedName

func StaticGet_Es6RewriteModules_GOOG_REQUIRETYPE() pkg_rhino.QualifiedName {
	clinit_Es6RewriteModules()
	return static_Es6RewriteModules_GOOG_REQUIRETYPE
}
func StaticSet_Es6RewriteModules_GOOG_REQUIRETYPE(v pkg_rhino.QualifiedName) {
	clinit_Es6RewriteModules()
	static_Es6RewriteModules_GOOG_REQUIRETYPE = v
}

var static_Es6RewriteModules_GOOG_MODULE_GET pkg_rhino.QualifiedName

func StaticGet_Es6RewriteModules_GOOG_MODULE_GET() pkg_rhino.QualifiedName {
	clinit_Es6RewriteModules()
	return static_Es6RewriteModules_GOOG_MODULE_GET
}
func StaticSet_Es6RewriteModules_GOOG_MODULE_GET(v pkg_rhino.QualifiedName) {
	clinit_Es6RewriteModules()
	static_Es6RewriteModules_GOOG_MODULE_GET = v
}

var static_Es6RewriteModules_GOOG_FORWARDDECLARE pkg_rhino.QualifiedName

func StaticGet_Es6RewriteModules_GOOG_FORWARDDECLARE() pkg_rhino.QualifiedName {
	clinit_Es6RewriteModules()
	return static_Es6RewriteModules_GOOG_FORWARDDECLARE
}
func StaticSet_Es6RewriteModules_GOOG_FORWARDDECLARE(v pkg_rhino.QualifiedName) {
	clinit_Es6RewriteModules()
	static_Es6RewriteModules_GOOG_FORWARDDECLARE = v
}

var static_Es6RewriteModules_GOOG_DECLAREMODULEID pkg_rhino.QualifiedName

func StaticGet_Es6RewriteModules_GOOG_DECLAREMODULEID() pkg_rhino.QualifiedName {
	clinit_Es6RewriteModules()
	return static_Es6RewriteModules_GOOG_DECLAREMODULEID
}
func StaticSet_Es6RewriteModules_GOOG_DECLAREMODULEID(v pkg_rhino.QualifiedName) {
	clinit_Es6RewriteModules()
	static_Es6RewriteModules_GOOG_DECLAREMODULEID = v
}

var clinitDone_Es6RewriteModules uint32
var clinitState_Es6RewriteModules j2gort.ClinitState

func clinit_Es6RewriteModules() {
	if atomic.LoadUint32(&clinitDone_Es6RewriteModules) != 0 {
		return
	}
	clinitSlow_Es6RewriteModules()
}

func clinitSlow_Es6RewriteModules() {
	if !j2gort.EnterClinit(&clinitState_Es6RewriteModules, &clinitDone_Es6RewriteModules) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6RewriteModules, &clinitDone_Es6RewriteModules)
	static_Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST = DiagnosticType_Error("JSC_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST", "The left side of a goog.require() or goog.requireType() must use ''const'' (not ''let'' or ''var'')")
	static_Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST = DiagnosticType_Error("JSC_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST", "goog.requireType alias for ES6 module should be const.")
	static_Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST = DiagnosticType_Error("JSC_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST", "goog.forwardDeclare alias for ES6 module should be const.")
	static_Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE = DiagnosticType_Warning("JSC_SHOULD_IMPORT_ES6_MODULE", "ES6 modules should import other ES6 modules rather than goog.require them.")
	static_Es6RewriteModules_DOT_SPLITTER = j2gort.Ext_com_google_common_base_Splitter_On_1_string(".")
	static_Es6RewriteModules_GOOG_REQUIRE = pkg_rhino.QualifiedName_Of("goog.require")
	static_Es6RewriteModules_GOOG_REQUIRETYPE = pkg_rhino.QualifiedName_Of("goog.requireType")
	static_Es6RewriteModules_GOOG_MODULE_GET = pkg_rhino.QualifiedName_Of("goog.module.get")
	static_Es6RewriteModules_GOOG_FORWARDDECLARE = pkg_rhino.QualifiedName_Of("goog.forwardDeclare")
	static_Es6RewriteModules_GOOG_DECLAREMODULEID = pkg_rhino.QualifiedName_Of("goog.declareModuleId")
}

func (this_ *Es6RewriteModules) initFields_Es6RewriteModules() {
	_ = this_
}

// Java constructor: Es6RewriteModules(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap, com.google.javascript.jscomp.modules.ModuleMap moduleMap, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable PreprocessorSymbolTable preprocessorSymbolTable, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope globalTypedScope)
func NewEs6RewriteModules_5(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap, moduleMap *ModuleMap, preprocessorSymbolTable *PreprocessorSymbolTable, globalTypedScope *TypedScope) *Es6RewriteModules {
	clinit_Es6RewriteModules()
	this_ := &Es6RewriteModules{}
	InitEs6RewriteModules_5(this_, compiler, moduleMetadataMap, moduleMap, preprocessorSymbolTable, globalTypedScope)
	return this_
}

func InitEs6RewriteModules_5(this_ *Es6RewriteModules, compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap, moduleMap *ModuleMap, preprocessorSymbolTable *PreprocessorSymbolTable, globalTypedScope *TypedScope) {
	_ = compiler
	_ = moduleMetadataMap
	_ = moduleMap
	_ = preprocessorSymbolTable
	_ = globalTypedScope
	InitEs6RewriteModules_6(this_, compiler, moduleMetadataMap, moduleMap, preprocessorSymbolTable, globalTypedScope, StaticGet_CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE())
}

// Java constructor: Es6RewriteModules(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.modules.ModuleMetadataMap moduleMetadataMap, com.google.javascript.jscomp.modules.ModuleMap moduleMap, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable PreprocessorSymbolTable preprocessorSymbolTable, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope globalTypedScope, com.google.javascript.jscomp.CompilerOptions.ChunkOutputType chunkOutputType)
func NewEs6RewriteModules_6(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap, moduleMap *ModuleMap, preprocessorSymbolTable *PreprocessorSymbolTable, globalTypedScope *TypedScope, chunkOutputType CompilerOptions_ChunkOutputType) *Es6RewriteModules {
	clinit_Es6RewriteModules()
	this_ := &Es6RewriteModules{}
	InitEs6RewriteModules_6(this_, compiler, moduleMetadataMap, moduleMap, preprocessorSymbolTable, globalTypedScope, chunkOutputType)
	return this_
}

func InitEs6RewriteModules_6(this_ *Es6RewriteModules, compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap, moduleMap *ModuleMap, preprocessorSymbolTable *PreprocessorSymbolTable, globalTypedScope *TypedScope, chunkOutputType CompilerOptions_ChunkOutputType) {
	_ = compiler
	_ = moduleMetadataMap
	_ = moduleMap
	_ = preprocessorSymbolTable
	_ = globalTypedScope
	_ = chunkOutputType
	this_.initFields_Es6RewriteModules()
	j2gort.CheckNotNullPtr(moduleMetadataMap)
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.moduleMetadataMap = moduleMetadataMap
	this_.moduleMap = moduleMap
	this_.preprocessorSymbolTable = preprocessorSymbolTable
	this_.globalTypedScope = globalTypedScope
	this_.unknownType = compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	this_.chunkOutputType = chunkOutputType
}

// Java static method: Es6RewriteModules.isEs6ModuleRoot(com.google.javascript.rhino.Node scriptNode)
func Es6RewriteModules_IsEs6ModuleRoot(scriptNode pkg_rhino.Node) bool {
	clinit_Es6RewriteModules()
	_ = scriptNode
	if !(scriptNode.IsScript()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(scriptNode))
	}
	if scriptNode.GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE()) {
		return false
	}
	return (scriptNode.HasChildren() && scriptNode.GetFirstChild().IsModuleBody())
}

// Java method: Es6RewriteModules.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteModules) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(externs.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(externs))
	}
	if !(root.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(root))
	}
	NodeTraversal_TraverseRoots(this_.compiler, this_, externs, root)
	NodeUtil_RemoveFeatureFromAllScripts(root, pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), this_.compiler)
	GatherGetterAndSetterProperties_Update(this_.compiler, externs, root)
}

// Java method: Es6RewriteModules.clearPerFileState()
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_clearPerFileState_2741() {
	this_.typedefs = &j2gort.OrderedSet[string]{}
	this_.namesToInlineByAlias = &j2gort.OrderedMap[string, *ModuleRenaming_GlobalizedModuleName]{}
}

// Java method: Es6RewriteModules.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModules) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if n.IsScript() {
		nodeTraversal.GetScope()
		NewEs6RewriteModules_RewriteRequiresForEs6Modules(this_).Rewrite(n)
		if Es6RewriteModules_IsEs6ModuleRoot(n) {
			this_.priv_Es6RewriteModules_clearPerFileState_2741()
		} else {
			return false
		}
	}
	return true
}

// Java method: Es6RewriteModules.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModules) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsImport() {
		this_.priv_Es6RewriteModules_maybeWarnExternModule_2742(t, n, parent)
		this_.priv_Es6RewriteModules_visitImport_2743(t, n, parent)
	} else if n.IsExport() {
		this_.priv_Es6RewriteModules_maybeWarnExternModule_2742(t, n, parent)
		this_.priv_Es6RewriteModules_visitExport_2744(t, n, parent)
	} else if n.IsScript() {
		this_.priv_Es6RewriteModules_visitScript_2748(t, n)
	} else if n.IsCall() {
		if StaticGet_Es6RewriteModules_GOOG_DECLAREMODULEID().Matches(n.GetFirstChild()) {
			n.GetParent().Detach()
		}
	} else if n.IsImportMeta() {
		if this_.chunkOutputType != StaticGet_CompilerOptions_ChunkOutputType_ES_MODULES() {
			t.Report_3(n, StaticGet_TranspilationUtil_CANNOT_CONVERT(), []string{"import.meta. Use --chunk_output_type=ES_MODULES to allow passthrough support."})
		}
	}
}

// Java method: Es6RewriteModules.maybeWarnExternModule(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_maybeWarnExternModule_2742(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(parent.IsModuleBody()) {
		j2gort.CheckState(false)
	}
	if parent.IsFromExterns() && !(NodeUtil_IsFromTypeSummary(parent.GetParent())) {
		t.Report_3(n, StaticGet_TranspilationUtil_CANNOT_CONVERT_YET(), []string{"ES6 modules in externs"})
	}
}

// Java method: Es6RewriteModules.visitImport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node importDecl, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_visitImport_2743(t *NodeTraversal, importDecl pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = importDecl
	_ = parent
	if !(parent.IsModuleBody()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(parent))
	}
	var importName string = importDecl.GetLastChild().GetString()
	_ = importName
	var isNamespaceImport bool = javalang.StringStartsWith(j2gort.CheckNotNullString(importName), "goog:")
	_ = isNamespaceImport
	if isNamespaceImport {
		var namespace string = j2gort.JStrSubFrom(importName, j2gort.JStrLen("goog:"))
		_ = namespace
		var m ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleMetadataMap.GetModulesByGoogNamespace(), namespace))
		_ = m
		if m == nil || m.IsJavaNil_() {
			t.Report_3(importDecl, StaticGet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE(), []string{namespace})
		} else {
			if !((m.GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module() || m.GetModuleMetadataMap_ModuleMetadataBase().IsGoogModule()) || m.GetModuleMetadataMap_ModuleMetadataBase().IsGoogProvide()) {
				j2gort.CheckState(false)
			}
		}
	} else {
		var modulePath *ModuleLoader_ModulePath = t.GetInput().GetPath().ResolveJsModule(importName, importDecl.GetSourceFileName(), importDecl.GetLineno(), importDecl.GetCharno())
		_ = modulePath
		if modulePath == nil {
			modulePath = t.GetInput().GetPath().ResolveModuleAsPath(importName)
		}
		this_.priv_Es6RewriteModules_maybeAddImportedFileReferenceToSymbolTable_2758(importDecl.GetLastChild(), modulePath.String())
	}
	{
		var child pkg_rhino.Node = importDecl.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsImportSpecs() {
				{
					var grandChild pkg_rhino.Node = child.GetFirstChild()
					_ = grandChild
					for ; grandChild != 0; grandChild = grandChild.GetNext() {
						this_.priv_Es6RewriteModules_maybeAddAliasToSymbolTable_2757(grandChild.GetFirstChild(), t.GetSourceName())
						if !(grandChild.HasTwoChildren()) {
							j2gort.CheckState(false)
						}
					}
				}
			} else if child.IsImportStar() {
				this_.priv_Es6RewriteModules_maybeAddAliasToSymbolTable_2757(child, t.GetSourceName())
			}
		}
	}
	importDecl.Detach()
	t.ReportCodeChange_0()
}

// Java method: Es6RewriteModules.visitExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_visitExport_2744(t *NodeTraversal, export pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = export
	_ = parent
	if !(parent.IsModuleBody()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(parent))
	}
	if export.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
		var child pkg_rhino.Node = export.GetFirstChild()
		_ = child
		var name string = j2gort.NullString
		_ = name
		if child.IsFunction() || child.IsClass() {
			name = NodeUtil_GetName(child)
		}
		if name != j2gort.NullString {
			var decl pkg_rhino.Node = child.Detach()
			_ = decl
			export.ReplaceWith(decl)
		} else {
			var var_ pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_2("$jscompDefaultExport", export.RemoveFirstChild())
			_ = var_
			var_.SetJSDocInfo(child.GetJSDocInfo())
			child.SetJSDocInfo(nil)
			var_.SrcrefTreeIfMissing(export)
			export.ReplaceWith(var_)
		}
		t.ReportCodeChange_0()
	} else if (export.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_ALL_FROM()) || export.HasTwoChildren()) || export.GetFirstChild().IsExportSpecs() {
		export.Detach()
		t.ReportCodeChange_0()
	} else {
		this_.priv_Es6RewriteModules_visitExportDeclaration_2746(t, export)
	}
}

// Java method: Es6RewriteModules.visitExportNameDeclaration(com.google.javascript.rhino.Node declaration)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_visitExportNameDeclaration_2745(declaration pkg_rhino.Node) {
	_ = declaration
	if (declaration.GetJSDocInfo() != nil) && declaration.GetJSDocInfo().HasTypedefType() {
		NodeUtil_VisitLhsNodesInNode(declaration, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_1 any) {
			lhs := j2gort.CastNode[pkg_rhino.Node](lp_1)
			_ = lhs
			this_.typedefs.Add(lhs.GetString())
		}))
	}
}

// Java method: Es6RewriteModules.visitExportDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_visitExportDeclaration_2746(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	var declaration pkg_rhino.Node = export.GetFirstChild()
	_ = declaration
	if NodeUtil_IsNameDeclaration(declaration) {
		this_.priv_Es6RewriteModules_visitExportNameDeclaration_2745(declaration)
	}
	export.ReplaceWith(declaration.Detach())
	t.ReportCodeChange_0()
}

// Java method: Es6RewriteModules.inlineModuleToGlobalScope(com.google.javascript.rhino.Node moduleNode)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_inlineModuleToGlobalScope_2747(moduleNode pkg_rhino.Node) {
	_ = moduleNode
	if !(moduleNode.IsModuleBody()) {
		j2gort.CheckState(false)
	}
	var scriptNode pkg_rhino.Node = moduleNode.GetParent()
	_ = scriptNode
	moduleNode.Detach()
	scriptNode.AddChildrenToFront(moduleNode.RemoveChildren())
}

// Java method: Es6RewriteModules.visitScript(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node script)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_visitScript_2748(t *NodeTraversal, script pkg_rhino.Node) {
	_ = t
	_ = script
	var moduleBody pkg_rhino.Node = script.GetFirstChild()
	_ = moduleBody
	var moduleObjectType AstFactory_Type = AstFactory_S_Type_1_Node(moduleBody)
	_ = moduleObjectType
	this_.priv_Es6RewriteModules_inlineModuleToGlobalScope_2747(moduleBody)
	ClosureRewriteModule_CheckAndSetStrictModeDirective(t, script)
	var thisModule *Module = this_.moduleMap.GetModule_1_ModuleLoader_ModulePath(t.GetInput().GetPath())
	_ = thisModule
	var qualifiedName pkg_rhino.QualifiedName = ModuleRenaming_GetGlobalName_2(thisModule.Metadata(), j2gort.NullString)
	_ = qualifiedName
	if !(qualifiedName.IsSimple()) {
		j2gort.CheckStateFmt(false, "Unexpected qualified name %s", qualifiedName)
	}
	var moduleName string = qualifiedName.GetQualifiedNameBase().GetRoot()
	_ = moduleName
	var moduleVar pkg_rhino.Node = this_.priv_Es6RewriteModules_createExportsObject_2749(moduleName, t, script, moduleObjectType)
	_ = moduleVar
	NodeTraversal_Traverse(this_.compiler, script, NewEs6RewriteModules_RenameGlobalVars(this_, thisModule))
	moduleVar.GetFirstChild().SetString(moduleName)
	moduleVar.MakeNonIndexableRecursive()
	this_.priv_Es6RewriteModules_declareGlobalVariable_2759(moduleVar.GetFirstChild(), t)
	this_.priv_Es6RewriteModules_rewriteRequires_2751(script)
	t.ReportCodeChange_0()
}

// Java method: Es6RewriteModules.createExportsObject(java.lang.String moduleName, com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node script, com.google.javascript.jscomp.AstFactory.Type moduleObjectType)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_createExportsObject_2749(moduleName string, t *NodeTraversal, script pkg_rhino.Node, moduleObjectType AstFactory_Type) pkg_rhino.Node {
	_ = moduleName
	_ = t
	_ = script
	_ = moduleObjectType
	var moduleObject pkg_rhino.Node = this_.astFactory.CreateObjectLit_2(moduleObjectType, []pkg_rhino.Node{})
	_ = moduleObject
	var moduleVar pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_2("$jscomp$tmp$exports$module$name", moduleObject)
	_ = moduleVar
	moduleVar.GetFirstChild().PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
	var infoBuilder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
	_ = infoBuilder
	infoBuilder.RecordConstancy()
	moduleVar.SetJSDocInfo(infoBuilder.Build_0())
	moduleVar.GetFirstChild().SetDeclaredConstantVar(true)
	script.AddChildToBack(moduleVar.SrcrefTreeIfMissing(script))
	var thisModule *Module = this_.moduleMap.GetModule_1_ModuleLoader_ModulePath(t.GetInput().GetPath())
	_ = thisModule
	{
		cur_2 := j2gort.NewMapEntryCursor(thisModule.Namespace())
		for cur_2.NextEntry() {
			var exportedName string = cur_2.KeyStr()
			_ = exportedName
			var binding *Binding = j2gort.CastPtr[Binding](cur_2.GetValue())
			_ = binding
			var nodeForSourceInfo pkg_rhino.Node = binding.SourceNode()
			_ = nodeForSourceInfo
			var mutated bool = binding.IsMutated()
			_ = mutated
			var boundVariableQualifiedName pkg_rhino.QualifiedName = ModuleRenaming_GetGlobalName_1_Binding(binding)
			_ = boundVariableQualifiedName
			if !(boundVariableQualifiedName.IsSimple()) {
				j2gort.CheckStateFmt(false, "unexpected qualified name: %s", boundVariableQualifiedName)
			}
			var boundVariableName string = boundVariableQualifiedName.GetQualifiedNameBase().GetRoot()
			_ = boundVariableName
			var getProp pkg_rhino.Node = this_.astFactory.CreateGetPropWithoutColor(this_.astFactory.CreateName_2_string_AstFactory_Type(moduleName, moduleObjectType), exportedName)
			_ = getProp
			getProp.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
			if this_.typedefs.Contains(exportedName) {
				var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
				_ = builder
				var typeExpr pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(this_.astFactory.CreateString(exportedName).Srcref(nodeForSourceInfo), script.GetSourceFileName())
				_ = typeExpr
				builder.RecordTypedef(typeExpr)
				var info *pkg_rhino.JSDocInfo = builder.Build_0()
				_ = info
				getProp.SetJSDocInfo(info)
				var exprResult pkg_rhino.Node = this_.astFactory.ExprResult(getProp).SrcrefTreeIfMissing(nodeForSourceInfo)
				_ = exprResult
				script.AddChildToBack(exprResult)
			} else if mutated {
				var globalExportName pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(boundVariableName, AstFactory_S_Type_1_Node(getProp))
				_ = globalExportName
				this_.priv_Es6RewriteModules_addGetterExport_2750(script, nodeForSourceInfo, moduleObject, exportedName, globalExportName)
				NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_GETTER(), this_.compiler)
			} else {
				var assign pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(getProp, this_.astFactory.CreateName_2_string_AstFactory_Type(boundVariableName, AstFactory_S_Type_1_Node(getProp)))
				_ = assign
				var builder_2 *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
				_ = builder_2
				builder_2.RecordConstancy()
				var info_2 *pkg_rhino.JSDocInfo = builder_2.Build_0()
				_ = info_2
				assign.SetJSDocInfo(info_2)
				script.AddChildToBack(this_.astFactory.ExprResult(assign).SrcrefTreeIfMissing(nodeForSourceInfo))
			}
		}
	}
	return moduleVar
}

// Java method: Es6RewriteModules.addGetterExport(com.google.javascript.rhino.Node script, com.google.javascript.rhino.Node forSourceInfo, com.google.javascript.rhino.Node objLit, java.lang.String exportedName, com.google.javascript.rhino.Node value)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_addGetterExport_2750(script pkg_rhino.Node, forSourceInfo pkg_rhino.Node, objLit pkg_rhino.Node, exportedName string, value pkg_rhino.Node) {
	_ = script
	_ = forSourceInfo
	_ = objLit
	_ = exportedName
	_ = value
	var getter pkg_rhino.Node = this_.astFactory.CreateGetterDef(exportedName, value)
	_ = getter
	getter.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT(), true)
	objLit.AddChildToBack(getter)
	if !(this_.astFactory.IsAddingTypes()) {
		var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder().ParseDocumentation()
		_ = builder
		builder.RecordReturnType(pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_QMARK()).Srcref(forSourceInfo), script.GetSourceFileName()))
		getter.SetJSDocInfo(builder.Build_0())
	} else {
		getter.SetJSType(pkg_rhino.ToIface_JSType(this_.compiler.GetTypeRegistry().CreateFunctionType_2_JSType_JSType(value.GetJSType(), []pkg_rhino.JSType{})))
	}
	getter.SrcrefTreeIfMissing(forSourceInfo)
	this_.compiler.ReportChangeToEnclosingScope(getter.GetFirstChild().GetLastChild())
	this_.compiler.ReportChangeToEnclosingScope(getter)
}

// Java method: Es6RewriteModules.rewriteRequires(com.google.javascript.rhino.Node script)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_rewriteRequires_2751(script pkg_rhino.Node) {
	_ = script
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(newReleasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface(func(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
		_ = t
		_ = n
		_ = parent
		if n.IsCall() {
			var fn pkg_rhino.Node = n.GetFirstChild()
			_ = fn
			if StaticGet_Es6RewriteModules_GOOG_REQUIRE().Matches(fn) || StaticGet_Es6RewriteModules_GOOG_REQUIRETYPE().Matches(fn) {
				this_.priv_Es6RewriteModules_visitRequireOrGet_2755(t, n, parent, true)
			} else if StaticGet_Es6RewriteModules_GOOG_MODULE_GET().Matches(fn) {
				this_.priv_Es6RewriteModules_visitGoogModuleGet_2753(t, n, parent)
			}
		}
	})).Traverse(script)
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(newReleasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface(func(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
		_ = t
		_ = n
		_ = parent
		var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = info
		if info != nil {
			{
				var i int32 = int32(0)
				_ = i
				var count int32 = info.GetTypeNodeCount()
				_ = count
				for ; i < count; i = i + 1 {
					this_.priv_Es6RewriteModules_inlineAliasedTypes_2752(t, info.GetTypeNodeAt(i))
				}
			}
		}
		if n.IsName() && this_.namesToInlineByAlias.ContainsKey(n.GetString()) {
			var v *Var = t.GetScope().LookupVar(n.GetString())
			_ = v
			if (v == nil) || (v.GetNameNode() != n) {
				var replacementName *ModuleRenaming_GlobalizedModuleName = this_.namesToInlineByAlias.Get(n.GetString())
				_ = replacementName
				var replacement pkg_rhino.Node = replacementName.ToQname(this_.astFactory).SrcrefTree(n)
				_ = replacement
				n.ReplaceWith(replacement)
			}
		}
	})).Traverse(script)
}

// Java method: Es6RewriteModules.inlineAliasedTypes(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node typeNode)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_inlineAliasedTypes_2752(t *NodeTraversal, typeNode pkg_rhino.Node) {
	_ = t
	_ = typeNode
	if typeNode.IsStringLit() {
		var name string = typeNode.GetString()
		_ = name
		var dot int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
		_ = dot
		var baseName string
		_ = baseName
		if dot < int32(0) {
			baseName = name
		} else {
			baseName = j2gort.JStrSub(name, int32(0), dot)
		}
		if t.GetScope().LookupVar(baseName) == nil {
			var replacement *ModuleRenaming_GlobalizedModuleName = this_.namesToInlineByAlias.Get(baseName)
			_ = replacement
			if replacement != nil {
				var rest string
				_ = rest
				if dot < int32(0) {
					rest = ""
				} else {
					rest = j2gort.JStrSubFrom(name, dot)
				}
				typeNode.SetOriginalName(name)
				typeNode.SetString((j2gort.StrValString(replacement.AliasName().Join_0()) + j2gort.StrValString(rest)))
				t.ReportCodeChange_0()
			}
		}
	}
	{
		var child pkg_rhino.Node = typeNode.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_Es6RewriteModules_inlineAliasedTypes_2752(t, child)
		}
	}
}

// Java method: Es6RewriteModules.visitGoogModuleGet(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node getCall, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_visitGoogModuleGet_2753(t *NodeTraversal, getCall pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = getCall
	_ = parent
	if !(getCall.HasTwoChildren()) || !(getCall.GetLastChild().IsStringLit()) {
		t.Report_3(getCall, StaticGet_ClosurePrimitiveErrors_INVALID_GET_NAMESPACE(), []string{})
		return
	}
	if t.InGlobalHoistScope() {
		t.Report_3(getCall, StaticGet_ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET(), []string{})
		return
	}
	this_.priv_Es6RewriteModules_visitRequireOrGet_2755(t, getCall, parent, false)
}

// Java method: Es6RewriteModules.getFallbackMetadataForNamespace(java.lang.String namespace)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_getFallbackMetadataForNamespace_2754(namespace string) ModuleMetadataMap_ModuleMetadata {
	_ = namespace
	var builder ModuleMetadataMap_ModuleMetadata_Builder = ModuleMetadataMap_ModuleMetadata_S_Builder().ModuleType_1(StaticGet_ModuleMetadataMap_ModuleType_GOOG_PROVIDE()).UsesClosure(true).IsTestOnly(false)
	_ = builder
	builder.GoogNamespacesBuilder().Add_1_any(j2gort.BoxString(namespace))
	return builder.Build()
}

// Java method: Es6RewriteModules.visitRequireOrGet(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node requireCall, com.google.javascript.rhino.Node parent, boolean isRequire)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_visitRequireOrGet_2755(t *NodeTraversal, requireCall pkg_rhino.Node, parent pkg_rhino.Node, isRequire bool) {
	_ = t
	_ = requireCall
	_ = parent
	_ = isRequire
	if !(requireCall.HasTwoChildren()) || !(requireCall.GetLastChild().IsStringLit()) {
		t.Report_3(requireCall, StaticGet_ClosurePrimitiveErrors_INVALID_REQUIRE_NAMESPACE(), []string{})
		return
	}
	if isRequire && !(t.GetScope().IsGlobal()) {
		t.Report_3(requireCall, StaticGet_ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR(), []string{})
		return
	}
	var namespace string = requireCall.GetLastChild().GetString()
	_ = namespace
	var isStoredInDeclaration bool = NodeUtil_IsDeclaration(parent.GetParent())
	_ = isStoredInDeclaration
	if isStoredInDeclaration && !(parent.GetParent().IsConst()) {
		this_.compiler.Report_1(JSError_Make_3(parent.GetParent(), StaticGet_Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST(), []string{}))
	}
	var m ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.moduleMetadataMap.GetModulesByGoogNamespace(), namespace))
	_ = m
	var isFromFallbackMetadata bool = (m == nil || m.IsJavaNil_())
	_ = isFromFallbackMetadata
	if isFromFallbackMetadata {
		t.Report_3(requireCall, StaticGet_ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE(), []string{namespace})
		m = this_.priv_Es6RewriteModules_getFallbackMetadataForNamespace_2754(namespace)
	}
	if isStoredInDeclaration {
		if isRequire {
			var toDetach pkg_rhino.Node = pkg_rhino.Node(0)
			_ = toDetach
			if parent.IsDestructuringLhs() {
				if !(parent.GetFirstChild().IsObjectPattern()) {
					j2gort.CheckState(false)
				}
				toDetach = parent.GetParent()
				{
					var child pkg_rhino.Node = parent.GetFirstFirstChild()
					_ = child
					for ; child != 0; child = child.GetNext() {
						if !(child.IsStringKey() && child.GetFirstChild().IsName()) {
							j2gort.CheckStateMsg(false, j2gort.BoxNode(child))
						}
						var rep *ModuleRenaming_GlobalizedModuleName = this_.priv_Es6RewriteModules_getGlobalNameAndType_2756(m, namespace, isFromFallbackMetadata).Getprop(child.GetString())
						_ = rep
						this_.namesToInlineByAlias.Put(child.GetFirstChild().GetString(), rep)
					}
				}
			} else if parent.IsName() {
				var alias *ModuleRenaming_GlobalizedModuleName = this_.priv_Es6RewriteModules_getGlobalNameAndType_2756(m, namespace, isFromFallbackMetadata)
				_ = alias
				this_.namesToInlineByAlias.Put(parent.GetString(), alias)
				toDetach = parent.GetParent()
			} else {
				if !(parent.IsExprResult()) {
					j2gort.CheckState(false)
				}
				toDetach = parent
			}
			toDetach.Detach()
		} else {
			var name *ModuleRenaming_GlobalizedModuleName = this_.priv_Es6RewriteModules_getGlobalNameAndType_2756(m, namespace, isFromFallbackMetadata)
			_ = name
			var replacement pkg_rhino.Node = name.ToQname(this_.astFactory).SrcrefTree(requireCall)
			_ = replacement
			requireCall.ReplaceWith(replacement)
		}
	} else {
		if !(requireCall.GetParent().IsExprResult()) {
			j2gort.CheckState(false)
		}
		requireCall.GetParent().Detach()
	}
}

// Java method: Es6RewriteModules.getGlobalNameAndType(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, java.lang.@org.jspecify.annotations.Nullable String googNamespace, boolean isFromMissingModuleOrProvide)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_getGlobalNameAndType_2756(metadata ModuleMetadataMap_ModuleMetadata, googNamespace string, isFromMissingModuleOrProvide bool) *ModuleRenaming_GlobalizedModuleName {
	_ = metadata
	_ = googNamespace
	_ = isFromMissingModuleOrProvide
	if isFromMissingModuleOrProvide {
		var globalName pkg_rhino.QualifiedName = ModuleRenaming_GetGlobalName_2(metadata, googNamespace)
		_ = globalName
		return ModuleRenaming_GlobalizedModuleName_Create_2(globalName, this_.unknownType)
	}
	return ModuleRenaming_GlobalizedModuleName_Create_3(metadata, googNamespace, this_.globalTypedScope)
}

// Java method: Es6RewriteModules.maybeAddAliasToSymbolTable(com.google.javascript.rhino.Node n, java.lang.String module)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_maybeAddAliasToSymbolTable_2757(n pkg_rhino.Node, module string) {
	_ = n
	_ = module
	if this_.preprocessorSymbolTable == nil {
		return
	}
	n.PutBooleanProp(pkg_rhino.StaticGet_Node_MODULE_ALIAS(), true)
	var nodeName string
	_ = nodeName
	if n.IsStringLit() || n.IsImportStar() {
		nodeName = n.GetString()
	} else {
		nodeName = this_.preprocessorSymbolTable.GetQualifiedName(n)
	}
	var name string = ("alias_" + j2gort.StrValString(module) + "_" + j2gort.StrValString(nodeName))
	_ = name
	this_.preprocessorSymbolTable.AddReference_2(n, name)
}

// Java method: Es6RewriteModules.maybeAddImportedFileReferenceToSymbolTable(com.google.javascript.rhino.Node importStringNode, java.lang.String importedFilePath)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_maybeAddImportedFileReferenceToSymbolTable_2758(importStringNode pkg_rhino.Node, importedFilePath string) {
	_ = importStringNode
	_ = importedFilePath
	if this_.preprocessorSymbolTable == nil {
		return
	}
	if j2gort.IfaceEquals[pkg_rhino.SimpleSlot](j2gort.AsIface[pkg_rhino.SimpleSlot](this_.preprocessorSymbolTable.GetSlot(importedFilePath)), nil) {
		var scriptNode pkg_rhino.Node = this_.compiler.GetScriptNode(importedFilePath)
		_ = scriptNode
		if scriptNode != 0 {
			this_.preprocessorSymbolTable.AddReference_2(scriptNode, importedFilePath)
		}
	}
	this_.preprocessorSymbolTable.AddReference_2(importStringNode, importedFilePath)
}

// Java method: Es6RewriteModules.declareGlobalVariable(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6RewriteModules) priv_Es6RewriteModules_declareGlobalVariable_2759(n pkg_rhino.Node, t *NodeTraversal) {
	_ = n
	_ = t
	if !(n.IsName()) {
		j2gort.CheckState(false)
	}
	if !(this_.astFactory.IsAddingTypes()) {
		return
	}
	j2gort.CheckNotNullPtr(this_.globalTypedScope)
	var name string = n.GetString()
	_ = name
	if this_.globalTypedScope.HasOwnSlot(name) {
		t.Report_3(t.GetCurrentScript(), StaticGet_ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT(), []string{name})
	} else {
		var type_ pkg_rhino.JSType = j2gort.CheckNotNullIface(n.GetJSType())
		_ = type_
		this_.globalTypedScope.Declare(name, n, type_, t.GetInput(), false)
	}
}

// Java class: com.google.javascript.jscomp.Es6RewriteModules.RewriteRequiresForEs6Modules
type Es6RewriteModules_RewriteRequiresForEs6Modules struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_      *Es6RewriteModules
	renameTable j2gort.Ext_com_google_common_collect_Table `gson:"renameTable"`
}

func (this_ *Es6RewriteModules_RewriteRequiresForEs6Modules) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteModules_RewriteRequiresForEs6Modules) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteModules_RewriteRequiresForEs6Modules) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteModules$RewriteRequiresForEs6Modules", "RewriteRequiresForEs6Modules"
}
func (this_ *Es6RewriteModules_RewriteRequiresForEs6Modules) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteModules_RewriteRequiresForEs6Modules) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteModules.RewriteRequiresForEs6Modules", "com.google.javascript.jscomp.Es6RewriteModules$RewriteRequiresForEs6Modules", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteModules_RewriteRequiresForEs6Modules](v_)
	})
}
func (this_ *Es6RewriteModules_RewriteRequiresForEs6Modules) SetOuter_(o *Es6RewriteModules) {
	this_.outer_ = o
}
func (this_ *Es6RewriteModules_RewriteRequiresForEs6Modules) GetOuter_() *Es6RewriteModules {
	return this_.outer_
}
func (*Es6RewriteModules_RewriteRequiresForEs6Modules) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*Es6RewriteModules_RewriteRequiresForEs6Modules) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *Es6RewriteModules_RewriteRequiresForEs6Modules) initFields_Es6RewriteModules_RewriteRequiresForEs6Modules() {
	_ = this_
}

// Java constructor: RewriteRequiresForEs6Modules()
func NewEs6RewriteModules_RewriteRequiresForEs6Modules(outer_ *Es6RewriteModules) *Es6RewriteModules_RewriteRequiresForEs6Modules {
	this_ := &Es6RewriteModules_RewriteRequiresForEs6Modules{}
	this_.outer_ = outer_
	InitEs6RewriteModules_RewriteRequiresForEs6Modules(this_)
	return this_
}

func InitEs6RewriteModules_RewriteRequiresForEs6Modules(this_ *Es6RewriteModules_RewriteRequiresForEs6Modules) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6RewriteModules_RewriteRequiresForEs6Modules()
}

// Java method: RewriteRequiresForEs6Modules.rewrite(com.google.javascript.rhino.Node scriptNode)
func (this_ *Es6RewriteModules_RewriteRequiresForEs6Modules) Rewrite(scriptNode pkg_rhino.Node) {
	_ = scriptNode
	this_.renameTable = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Table](j2gort.Ext_com_google_common_collect_HashBasedTable_Create_0())
	NodeTraversal_Traverse(this_.outer_.compiler, scriptNode, this_)
	if !(this_.renameTable.IsEmpty()) {
		NodeTraversal_Traverse(this_.outer_.compiler, scriptNode, NewEs6RenameTypeReferences(this_.renameTable))
	}
}

// Java method: RewriteRequiresForEs6Modules.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModules_RewriteRequiresForEs6Modules) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsCall()) {
		return
	}
	var callTarget pkg_rhino.Node = n.GetFirstChild()
	_ = callTarget
	if !(callTarget.IsGetProp()) {
		return
	}
	var isRequire bool = StaticGet_Es6RewriteModules_GOOG_REQUIRE().Matches(callTarget)
	_ = isRequire
	var isRequireType bool = StaticGet_Es6RewriteModules_GOOG_REQUIRETYPE().Matches(callTarget)
	_ = isRequireType
	var isGet bool = StaticGet_Es6RewriteModules_GOOG_MODULE_GET().Matches(callTarget)
	_ = isGet
	var isForwardDeclare bool = StaticGet_Es6RewriteModules_GOOG_FORWARDDECLARE().Matches(callTarget)
	_ = isForwardDeclare
	if ((!(isRequire) && !(isRequireType)) && !(isGet)) && !(isForwardDeclare) {
		return
	}
	if !(n.HasTwoChildren()) || !(n.GetLastChild().IsStringLit()) {
		return
	}
	var name string = n.GetLastChild().GetString()
	_ = name
	var moduleMetadata ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.outer_.moduleMetadataMap.GetModulesByGoogNamespace(), name))
	_ = moduleMetadata
	if (moduleMetadata == nil || moduleMetadata.IsJavaNil_()) || !(moduleMetadata.GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module()) {
		return
	}
	if isRequire {
		var currentModuleMetadata ModuleMetadataMap_ModuleMetadata = j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.outer_.moduleMetadataMap.GetModulesByPath(), t.GetInput().GetPath().String()))
		_ = currentModuleMetadata
		if (currentModuleMetadata != nil && !currentModuleMetadata.IsJavaNil_()) && currentModuleMetadata.GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module() {
			t.Report_3(n, StaticGet_Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE(), []string{})
		}
	}
	if isGet && t.InGlobalHoistScope() {
		return
	}
	var statementNode pkg_rhino.Node = NodeUtil_GetEnclosingStatement(n)
	_ = statementNode
	var importHasAlias bool = NodeUtil_IsNameDeclaration(statementNode)
	_ = importHasAlias
	if importHasAlias {
		if statementNode.GetFirstChild().IsDestructuringLhs() {
			if isForwardDeclare {
				t.Report_3(n, StaticGet_ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE(), []string{})
				return
			}
			if isRequireType {
				if !(statementNode.IsConst()) {
					t.Report_3(statementNode, StaticGet_Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST(), []string{})
					return
				}
				{
					var child pkg_rhino.Node = statementNode.GetFirstFirstChild().GetFirstChild()
					_ = child
					for ; child != 0; child = child.GetNext() {
						if !(child.IsStringKey()) {
							j2gort.CheckState(false)
						}
						if !(child.GetFirstChild().IsName()) {
							j2gort.CheckState(false)
						}
						this_.renameTable.Put(j2gort.BoxNode(t.GetScopeRoot()), j2gort.BoxString(child.GetFirstChild().GetString()), j2gort.BoxString(ModuleRenaming_GetGlobalName_2(moduleMetadata, name).GetQualifiedNameBase().Getprop(child.GetString()).Join_0()))
					}
				}
			} else {
				{
					var child_2 pkg_rhino.Node = statementNode.GetFirstFirstChild().GetFirstChild()
					_ = child_2
					for ; child_2 != 0; child_2 = child_2.GetNext() {
						if !(child_2.IsStringKey()) {
							j2gort.CheckState(false)
						}
						if !(child_2.GetFirstChild().IsName()) {
							j2gort.CheckState(false)
						}
						var globalName *ModuleRenaming_GlobalizedModuleName = this_.outer_.priv_Es6RewriteModules_getGlobalNameAndType_2756(moduleMetadata, name, false).Getprop(child_2.GetString())
						_ = globalName
						var constNode pkg_rhino.Node = this_.outer_.astFactory.CreateSingleConstNameDeclaration(child_2.GetFirstChild().GetString(), globalName.ToQname(this_.outer_.astFactory))
						_ = constNode
						constNode.SrcrefTree(child_2)
						constNode.InsertBefore(statementNode)
					}
				}
			}
			statementNode.Detach()
			t.ReportCodeChange_0()
		} else {
			if isForwardDeclare || isRequireType {
				if !(statementNode.IsConst()) {
					var diagnostic *DiagnosticType
					_ = diagnostic
					if isForwardDeclare {
						diagnostic = StaticGet_Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST()
					} else {
						diagnostic = StaticGet_Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST()
					}
					t.Report_3(statementNode, diagnostic, []string{})
					return
				}
				this_.renameTable.Put(j2gort.BoxNode(t.GetScopeRoot()), j2gort.BoxString(statementNode.GetFirstChild().GetString()), j2gort.BoxString(ModuleRenaming_GetGlobalName_2(moduleMetadata, name).Join_0()))
				statementNode.Detach()
				t.ReportCodeChange_0()
			} else {
				n.ReplaceWith(this_.outer_.astFactory.CreateName_2_string_AstFactory_Type(ModuleRenaming_GetGlobalName_2(moduleMetadata, name).GetQualifiedNameBase().GetRoot(), AstFactory_S_Type_1_Node(n)).SrcrefTree(n))
				t.ReportCodeChange_0()
			}
		}
	} else {
		if isForwardDeclare || isRequireType {
			this_.renameTable.Put(j2gort.BoxNode(t.GetScopeRoot()), j2gort.BoxString(name), j2gort.BoxString(ModuleRenaming_GetGlobalName_2(moduleMetadata, name).Join_0()))
			statementNode.Detach()
		} else {
			if statementNode.IsExprResult() && (statementNode.GetFirstChild() == n) {
				statementNode.Detach()
			} else {
				n.ReplaceWith(this_.outer_.astFactory.CreateName_2_string_AstFactory_Type(ModuleRenaming_GetGlobalName_2(moduleMetadata, name).GetQualifiedNameBase().GetRoot(), AstFactory_S_Type_1_Node(n)).SrcrefTree(n))
			}
		}
		t.ReportCodeChange_0()
	}
}

// Java class: com.google.javascript.jscomp.Es6RewriteModules.RenameGlobalVars
type Es6RewriteModules_RenameGlobalVars struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_     *Es6RewriteModules
	thisModule *Module `gson:"thisModule"`
}

func (this_ *Es6RewriteModules_RenameGlobalVars) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteModules_RenameGlobalVars) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteModules_RenameGlobalVars) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteModules$RenameGlobalVars", "RenameGlobalVars"
}
func (this_ *Es6RewriteModules_RenameGlobalVars) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteModules_RenameGlobalVars) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteModules.RenameGlobalVars", "com.google.javascript.jscomp.Es6RewriteModules$RenameGlobalVars", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteModules_RenameGlobalVars](v_)
	})
}
func (this_ *Es6RewriteModules_RenameGlobalVars) SetOuter_(o *Es6RewriteModules) { this_.outer_ = o }
func (this_ *Es6RewriteModules_RenameGlobalVars) GetOuter_() *Es6RewriteModules  { return this_.outer_ }
func (*Es6RewriteModules_RenameGlobalVars) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*Es6RewriteModules_RenameGlobalVars) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *Es6RewriteModules_RenameGlobalVars) initFields_Es6RewriteModules_RenameGlobalVars() {
	_ = this_
}

// Java constructor: RenameGlobalVars(com.google.javascript.jscomp.modules.Module thisModule)
func NewEs6RewriteModules_RenameGlobalVars(outer_ *Es6RewriteModules, thisModule *Module) *Es6RewriteModules_RenameGlobalVars {
	this_ := &Es6RewriteModules_RenameGlobalVars{}
	this_.outer_ = outer_
	InitEs6RewriteModules_RenameGlobalVars(this_, thisModule)
	return this_
}

func InitEs6RewriteModules_RenameGlobalVars(this_ *Es6RewriteModules_RenameGlobalVars, thisModule *Module) {
	_ = thisModule
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6RewriteModules_RenameGlobalVars()
	this_.thisModule = thisModule
}

// Java method: RenameGlobalVars.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModules_RenameGlobalVars) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if info != nil {
		{
			var i int32 = int32(0)
			_ = i
			var count int32 = info.GetTypeNodeCount()
			_ = count
			for ; i < count; i = i + 1 {
				this_.priv_RenameGlobalVars_fixTypeNode_2762(t, info.GetTypeNodeAt(i))
			}
		}
	}
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		var var_ *Var = t.GetScope().LookupVar(name)
		_ = var_
		if (var_ != nil) && var_.IsGlobal() {
			var newName string = ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(this_.thisModule.Metadata(), name).Join_0()
			_ = newName
			n.SetString(newName)
			n.SetOriginalName(name)
			t.ReportCodeChange_1(n)
			if NodeUtil_IsDeclarationLValue(n) {
				this_.outer_.priv_Es6RewriteModules_declareGlobalVariable_2759(n, t)
			}
		} else if (var_ == nil) && j2gort.MapContainsKeyStr(this_.thisModule.BoundNames(), name) {
			this_.outer_.priv_Es6RewriteModules_maybeAddAliasToSymbolTable_2757(n, t.GetSourceName())
			var binding *Binding = j2gort.CastPtr[Binding](j2gort.MapGetStr(this_.thisModule.BoundNames(), name))
			_ = binding
			var replacement pkg_rhino.Node = this_.priv_RenameGlobalVars_replace_2761(n, binding)
			_ = replacement
			if ((replacement.IsName() && parent.IsCall()) && (parent.GetFirstChild() == n)) && parent.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL()) {
				parent.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
			}
			if NodeUtil_IsDeclarationLValue(n) {
				this_.outer_.priv_Es6RewriteModules_declareGlobalVariable_2759(n, t)
			}
			t.ReportCodeChange_0()
		}
	}
}

// Java method: RenameGlobalVars.replace(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.modules.Binding binding)
func (this_ *Es6RewriteModules_RenameGlobalVars) priv_RenameGlobalVars_replace_2761(n pkg_rhino.Node, binding *Binding) pkg_rhino.Node {
	_ = n
	_ = binding
	if !(n.IsName()) {
		j2gort.CheckState(false)
	}
	for (binding.IsModuleNamespace() && binding.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module()) && n.GetParent().IsGetProp() {
		var propertyName string = n.GetParent().GetString()
		_ = propertyName
		var m *Module = this_.outer_.moduleMap.GetModule_1_ModuleLoader_ModulePath(binding.Metadata().Path())
		_ = m
		if j2gort.MapContainsKeyStr(m.Namespace(), propertyName) {
			binding = j2gort.CastPtr[Binding](j2gort.MapGetStr(m.Namespace(), propertyName))
			n = n.GetParent()
		} else {
			break
		}
	}
	var globalName pkg_rhino.QualifiedName = ModuleRenaming_GetGlobalName_1_Binding(binding)
	_ = globalName
	var newNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newNode
	base_globalName_ := globalName.GetQualifiedNameBase()
	if !(globalName.IsSimple()) {
		var root string = base_globalName_.GetRoot()
		_ = root
		if (this_.outer_.globalTypedScope != nil) && this_.outer_.globalTypedScope.HasSlot(root) {
			newNode = this_.outer_.astFactory.CreateQNameUsingJSTypeInfo(this_.outer_.globalTypedScope, globalName.Join_0())
		} else {
			newNode = this_.outer_.astFactory.CreateQNameWithUnknownType_1_string(globalName.Join_0())
		}
	} else {
		newNode = this_.outer_.astFactory.CreateName_2_string_AstFactory_Type(base_globalName_.GetRoot(), AstFactory_S_Type_1_Node(n))
	}
	n.ReplaceWith(newNode)
	newNode.SrcrefTree(n)
	newNode.SetOriginalName(n.GetString())
	return newNode
}

// Java method: RenameGlobalVars.fixTypeNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node typeNode)
func (this_ *Es6RewriteModules_RenameGlobalVars) priv_RenameGlobalVars_fixTypeNode_2762(t *NodeTraversal, typeNode pkg_rhino.Node) {
	_ = t
	_ = typeNode
	if typeNode.IsStringLit() {
		var thisModule *Module = this_.outer_.moduleMap.GetModule_1_ModuleLoader_ModulePath(t.GetInput().GetPath())
		_ = thisModule
		var name string = typeNode.GetString()
		_ = name
		var dot int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
		_ = dot
		var baseName string
		_ = baseName
		if dot < int32(0) {
			baseName = name
		} else {
			baseName = j2gort.JStrSub(name, int32(0), dot)
		}
		var var_ *Var = t.GetScope().LookupVar(baseName)
		_ = var_
		if (var_ != nil) && var_.IsGlobal() {
			var rest string
			_ = rest
			if dot < int32(0) {
				rest = ""
			} else {
				rest = j2gort.JStrSubFrom(name, dot)
			}
			this_.priv_RenameGlobalVars_maybeSetNewName_2763(t, typeNode, name, (j2gort.StrValString(ModuleRenaming_GetGlobalNameOfEsModuleLocalVariable(thisModule.Metadata(), baseName).Join_0()) + j2gort.StrValString(rest)))
		} else if (var_ == nil) && j2gort.MapContainsKeyStr(thisModule.BoundNames(), baseName) {
			var binding *Binding = j2gort.CastPtr[Binding](j2gort.MapGetStr(thisModule.BoundNames(), baseName))
			_ = binding
			var propChain j2gort.Ext_java_util_List
			_ = propChain
			if dot < int32(0) {
				propChain = j2gort.Ext_java_util_List_Of_0()
			} else {
				propChain = j2gort.SplitterSplitToListStr(StaticGet_Es6RewriteModules_DOT_SPLITTER(), j2gort.JStrSubFrom(name, (dot+int32(1))))
			}
			var globalName string = ModuleRenaming_GetGlobalNameForJsDoc(this_.outer_.moduleMap, binding, propChain)
			_ = globalName
			this_.priv_RenameGlobalVars_maybeSetNewName_2763(t, typeNode, name, globalName)
			if this_.outer_.preprocessorSymbolTable != nil {
				var onlyBaseName pkg_rhino.Node = pkg_rhino.Node_NewString_1(baseName).Srcref(typeNode)
				_ = onlyBaseName
				onlyBaseName.SetLength(j2gort.JStrLen(baseName))
				this_.outer_.priv_Es6RewriteModules_maybeAddAliasToSymbolTable_2757(onlyBaseName, t.GetSourceName())
			}
		}
		typeNode.SetOriginalName(name)
	}
	{
		var child pkg_rhino.Node = typeNode.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_RenameGlobalVars_fixTypeNode_2762(t, child)
		}
	}
}

// Java method: RenameGlobalVars.maybeSetNewName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node, java.lang.String name, java.lang.String newName)
func (this_ *Es6RewriteModules_RenameGlobalVars) priv_RenameGlobalVars_maybeSetNewName_2763(t *NodeTraversal, node pkg_rhino.Node, name string, newName string) {
	_ = t
	_ = node
	_ = name
	_ = newName
	if !(j2gort.CheckNotNullString(name) == newName) {
		node.SetString(newName)
		node.SetOriginalName(name)
		t.ReportCodeChange_0()
	}
}
func (*Es6RewriteModules) CanReleaseSyntacticScopes_()                              {}
func (*Es6RewriteModules_RewriteRequiresForEs6Modules) AlwaysShouldTraverse_()      {}
func (*Es6RewriteModules_RewriteRequiresForEs6Modules) CanReleaseSyntacticScopes_() {}
func (*Es6RewriteModules_RenameGlobalVars) AlwaysShouldTraverse_()                  {}
func (*Es6RewriteModules_RenameGlobalVars) CanReleaseSyntacticScopes_()             {}

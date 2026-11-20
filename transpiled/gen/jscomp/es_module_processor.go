/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/modules/EsModuleProcessor.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.modules.EsModuleProcessor
type EsModuleProcessor struct {
	compiler             AbstractCompiler                           `gson:"compiler"`
	currentModuleBuilder *EsModuleProcessor_UnresolvedModuleBuilder `gson:"currentModuleBuilder"`
	metadata             ModuleMetadataMap_ModuleMetadata           `gson:"metadata"`
	idHash_              atomic.Int32
}

func (this_ *EsModuleProcessor) BindSelf_(s any) {
	_ = s
}
func (this_ *EsModuleProcessor) IsJavaNil_() bool { return this_ == nil }
func (*EsModuleProcessor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.EsModuleProcessor", "EsModuleProcessor"
}
func (this_ *EsModuleProcessor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*EsModuleProcessor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.EsModuleProcessor", "com.google.javascript.jscomp.modules.EsModuleProcessor", func(v_ any) bool {
		return j2gort.InstanceOf[*EsModuleProcessor](v_)
	})
}
func (*EsModuleProcessor) Implements_jscomp_NodeTraversal_Callback()           {}
func (*EsModuleProcessor) Implements_jscomp_ModuleMapCreator_ModuleProcessor() {}

var static_EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION *DiagnosticType

func StaticGet_EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION() *DiagnosticType {
	clinit_EsModuleProcessor()
	return static_EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION
}
func StaticSet_EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION(v *DiagnosticType) {
	clinit_EsModuleProcessor()
	static_EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION = v
}

var static_EsModuleProcessor_CYCLIC_EXPORT_DEFINITION *DiagnosticType

func StaticGet_EsModuleProcessor_CYCLIC_EXPORT_DEFINITION() *DiagnosticType {
	clinit_EsModuleProcessor()
	return static_EsModuleProcessor_CYCLIC_EXPORT_DEFINITION
}
func StaticSet_EsModuleProcessor_CYCLIC_EXPORT_DEFINITION(v *DiagnosticType) {
	clinit_EsModuleProcessor()
	static_EsModuleProcessor_CYCLIC_EXPORT_DEFINITION = v
}

var static_EsModuleProcessor_DUPLICATE_EXPORT *DiagnosticType

func StaticGet_EsModuleProcessor_DUPLICATE_EXPORT() *DiagnosticType {
	clinit_EsModuleProcessor()
	return static_EsModuleProcessor_DUPLICATE_EXPORT
}
func StaticSet_EsModuleProcessor_DUPLICATE_EXPORT(v *DiagnosticType) {
	clinit_EsModuleProcessor()
	static_EsModuleProcessor_DUPLICATE_EXPORT = v
}

var static_EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT *DiagnosticType

func StaticGet_EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT() *DiagnosticType {
	clinit_EsModuleProcessor()
	return static_EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT
}
func StaticSet_EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT(v *DiagnosticType) {
	clinit_EsModuleProcessor()
	static_EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT = v
}

var static_EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR *DiagnosticType

func StaticGet_EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR() *DiagnosticType {
	clinit_EsModuleProcessor()
	return static_EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR
}
func StaticSet_EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR(v *DiagnosticType) {
	clinit_EsModuleProcessor()
	static_EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR = v
}

var static_EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE *DiagnosticType

func StaticGet_EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE() *DiagnosticType {
	clinit_EsModuleProcessor()
	return static_EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE
}
func StaticSet_EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE(v *DiagnosticType) {
	clinit_EsModuleProcessor()
	static_EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE = v
}

var clinitDone_EsModuleProcessor uint32
var clinitState_EsModuleProcessor j2gort.ClinitState

func clinit_EsModuleProcessor() {
	if atomic.LoadUint32(&clinitDone_EsModuleProcessor) != 0 {
		return
	}
	clinitSlow_EsModuleProcessor()
}

func clinitSlow_EsModuleProcessor() {
	if !j2gort.EnterClinit(&clinitState_EsModuleProcessor, &clinitDone_EsModuleProcessor) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_EsModuleProcessor, &clinitDone_EsModuleProcessor)
	static_EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION = DiagnosticType_Warning("JSC_AMBIGUOUS_EXPORT_DEFINITION", "The export \"{0}\" is ambiguous.")
	static_EsModuleProcessor_CYCLIC_EXPORT_DEFINITION = DiagnosticType_Error("JSC_CYCLIC_EXPORT_DEFINITION", "Cyclic export detected while resolving name \"{0}\".")
	static_EsModuleProcessor_DUPLICATE_EXPORT = DiagnosticType_Error("JSC_DUPLICATE_EXPORT", "Duplicate export of \"{0}\".")
	static_EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT = DiagnosticType_Error("JSC_IMPORTED_AMBIGUOUS_EXPORT", "The requested name \"{0}\" is ambiguous.")
	static_EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR = DiagnosticType_Error("JSC_NAMESPACE_IMPORT_CANNOT_USE_STAR", "Namespace imports ('goog:some.Namespace') cannot use import * as. Did you mean to import {0} from ''{1}'';?")
	static_EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE = DiagnosticType_Error("JSC_CANNOT_PATH_IMPORT_CLOSURE_FILE", "Cannot import Closure files by path. Use either import 'goog:namespace' or goog.require('namespace')")
}

func (this_ *EsModuleProcessor) initFields_EsModuleProcessor() {
	_ = this_
}

// Java constructor: EsModuleProcessor(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEsModuleProcessor(compiler AbstractCompiler) *EsModuleProcessor {
	clinit_EsModuleProcessor()
	this_ := &EsModuleProcessor{}
	InitEsModuleProcessor(this_, compiler)
	return this_
}

func InitEsModuleProcessor(this_ *EsModuleProcessor, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_EsModuleProcessor()
	this_.compiler = compiler
}

// Java method: EsModuleProcessor.process(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, com.google.javascript.jscomp.deps.ModuleLoader.ModulePath path, com.google.javascript.rhino.Node script)
func (this_ *EsModuleProcessor) Process(metadata ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, script pkg_rhino.Node) UnresolvedModule {
	_ = metadata
	_ = path
	_ = script
	this_.metadata = metadata
	this_.currentModuleBuilder = NewEsModuleProcessor_UnresolvedModuleBuilder(this_, path, script)
	NodeTraversal_Traverse(this_.compiler, script, this_)
	var m UnresolvedModule = this_.currentModuleBuilder.Build()
	_ = m
	this_.metadata = nil
	this_.currentModuleBuilder = nil
	return m
}

// Java method: EsModuleProcessor.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *EsModuleProcessor) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(89) /* Token.ROOT */, int32(94) /* Token.SCRIPT */, int32(123) /* Token.MODULE_BODY */, int32(120) /* Token.EXPORT */, int32(116) /* Token.IMPORT */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: EsModuleProcessor.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *EsModuleProcessor) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(120) /* Token.EXPORT */ :
		this_.priv_EsModuleProcessor_visitExport_6734(t, n)
	case int32(116) /* Token.IMPORT */ :
		this_.priv_EsModuleProcessor_visitImport_6738(t, n)
	default:
	}
}

// Java method: EsModuleProcessor.visitExportAllFrom(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitExportAllFrom_6727(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	this_.currentModuleBuilder.Add_1_Export(Export_S_Builder().ExportName(j2gort.NullString).ModuleRequest(export.GetSecondChild().GetString()).ImportName("*").LocalName(j2gort.NullString).ModulePath(t.GetInput().GetPath()).ExportNode(export).ModuleMetadata(this_.metadata).GetExport_BuilderBase(

	// Java method: EsModuleProcessor.visitExportDefault(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
	).Build())
}

func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitExportDefault_6728(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	var child pkg_rhino.Node = export.GetFirstChild()
	_ = child
	var name string = "*default*"
	_ = name
	if child.IsFunction() || child.IsClass() {
		var maybeName string = NodeUtil_GetName(child)
		_ = maybeName
		if !(j2gort.Ext_com_google_common_base_Strings_IsNullOrEmpty(maybeName)) {
			name = maybeName
		}
	}
	if !(this_.currentModuleBuilder.Add_1_Export(Export_S_Builder().ExportName("default").ModuleRequest(j2gort.NullString).ImportName(j2gort.NullString).LocalName(name).ModulePath(t.GetInput().GetPath()).ExportNode(export).ModuleMetadata(this_.metadata).GetExport_BuilderBase().Build())) {
		t.Report_3(export, StaticGet_EsModuleProcessor_DUPLICATE_EXPORT(), []string{"default"})
	}
}

// Java method: EsModuleProcessor.visitExportFrom(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitExportFrom_6729(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	{
		var child pkg_rhino.Node = export.GetFirstFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			var importName string = child.GetFirstChild().GetString()
			_ = importName
			var exportedName string = child.GetLastChild().GetString()
			_ = exportedName
			if !(this_.currentModuleBuilder.Add_1_Export(Export_S_Builder().ExportName(exportedName).ModuleRequest(export.GetSecondChild().GetString()).ImportName(importName).LocalName(j2gort.NullString).ModulePath(t.GetInput().GetPath()).ExportNode(export).NameNode(child.GetFirstChild()).ModuleMetadata(this_.metadata).GetExport_BuilderBase().Build())) {
				t.Report_3(export, StaticGet_EsModuleProcessor_DUPLICATE_EXPORT(), []string{exportedName})
			}
		}
	}
}

// Java method: EsModuleProcessor.visitExportSpecs(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitExportSpecs_6730(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	{
		var child pkg_rhino.Node = export.GetFirstFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			var localName string = child.GetFirstChild().GetString()
			_ = localName
			var exportedName string = child.GetLastChild().GetString()
			_ = exportedName
			if !(this_.currentModuleBuilder.Add_1_Export(Export_S_Builder().ExportName(exportedName).ModuleRequest(j2gort.NullString).ImportName(j2gort.NullString).LocalName(localName).ModulePath(t.GetInput().GetPath()).ExportNode(export).NameNode(child.GetFirstChild()).ModuleMetadata(this_.metadata).GetExport_BuilderBase().Build())) {
				t.Report_3(export, StaticGet_EsModuleProcessor_DUPLICATE_EXPORT(), []string{exportedName})
			}
		}
	}
}

// Java method: EsModuleProcessor.visitExportNameDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export, com.google.javascript.rhino.Node declaration)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitExportNameDeclaration_6731(t *NodeTraversal, export pkg_rhino.Node, declaration pkg_rhino.Node) {
	_ = t
	_ = export
	_ = declaration
	NodeUtil_VisitLhsNodesInNode(declaration, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_1 any) {
		lhs := j2gort.CastNode[pkg_rhino.Node](lp_1)
		_ = lhs
		this_.priv_EsModuleProcessor_addExportNameDeclaration_6732(t, lhs, export)
	}))
}

// Java method: EsModuleProcessor.addExportNameDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node export)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_addExportNameDeclaration_6732(t *NodeTraversal, lhs pkg_rhino.Node, export pkg_rhino.Node) {
	_ = t
	_ = lhs
	_ = export
	if !(lhs.IsName()) {
		j2gort.CheckState(false)
	}
	var name string = lhs.GetString()
	_ = name
	if !(this_.currentModuleBuilder.Add_1_Export(Export_S_Builder().ExportName(name).ModuleRequest(j2gort.NullString).ImportName(j2gort.NullString).LocalName(name).ModulePath(t.GetInput().GetPath()).ExportNode(export).NameNode(lhs).ModuleMetadata(this_.metadata).GetExport_BuilderBase().Build())) {
		t.Report_3(export, StaticGet_EsModuleProcessor_DUPLICATE_EXPORT(), []string{name})
	}
}

// Java method: EsModuleProcessor.visitExportFunctionOrClass(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export, com.google.javascript.rhino.Node declaration)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitExportFunctionOrClass_6733(t *NodeTraversal, export pkg_rhino.Node, declaration pkg_rhino.Node) {
	_ = t
	_ = export
	_ = declaration
	if !(declaration.IsFunction() || declaration.IsClass()) {
		j2gort.CheckState(false)
	}
	var nameNode pkg_rhino.Node = declaration.GetFirstChild()
	_ = nameNode
	var name string = nameNode.GetString()
	_ = name
	if !(this_.currentModuleBuilder.Add_1_Export(Export_S_Builder().ExportName(name).ModuleRequest(j2gort.NullString).ImportName(j2gort.NullString).LocalName(name).ModulePath(t.GetInput().GetPath()).ExportNode(export).NameNode(nameNode).ModuleMetadata(this_.metadata).GetExport_BuilderBase().Build())) {
		t.Report_3(export, StaticGet_EsModuleProcessor_DUPLICATE_EXPORT(), []string{name})
	}
}

// Java method: EsModuleProcessor.visitExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitExport_6734(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	if export.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_ALL_FROM()) {
		this_.priv_EsModuleProcessor_visitExportAllFrom_6727(t, export)
	} else if export.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
		this_.priv_EsModuleProcessor_visitExportDefault_6728(t, export)
	} else if export.HasTwoChildren() {
		this_.priv_EsModuleProcessor_visitExportFrom_6729(t, export)
	} else if export.GetFirstChild().IsExportSpecs() {
		this_.priv_EsModuleProcessor_visitExportSpecs_6730(t, export)
	} else {
		var declaration pkg_rhino.Node = export.GetFirstChild()
		_ = declaration
		if NodeUtil_IsNameDeclaration(declaration) {
			this_.priv_EsModuleProcessor_visitExportNameDeclaration_6731(t, export, declaration)
		} else {
			this_.priv_EsModuleProcessor_visitExportFunctionOrClass_6733(t, export, declaration)
		}
	}
}

// Java method: EsModuleProcessor.visitImportDefault(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node importNode, java.lang.String moduleRequest)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitImportDefault_6735(t *NodeTraversal, importNode pkg_rhino.Node, moduleRequest string) {
	_ = t
	_ = importNode
	_ = moduleRequest
	this_.currentModuleBuilder.Add_1_Import(Import_S_Builder().ModuleRequest(moduleRequest).ImportName("default").LocalName(importNode.GetFirstChild().GetString()).ModulePath(t.GetInput().GetPath()).ImportNode(importNode).NameNode(importNode.GetFirstChild()).Build())
}

// Java method: EsModuleProcessor.visitImportSpecs(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node importNode, java.lang.String moduleRequest)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitImportSpecs_6736(t *NodeTraversal, importNode pkg_rhino.Node, moduleRequest string) {
	_ = t
	_ = importNode
	_ = moduleRequest
	{
		var child pkg_rhino.Node = importNode.GetSecondChild().GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			var importName string = child.GetFirstChild().GetString()
			_ = importName
			var localName string = child.GetLastChild().GetString()
			_ = localName
			this_.currentModuleBuilder.Add_1_Import(Import_S_Builder().ModuleRequest(moduleRequest).ImportName(importName).LocalName(localName).ModulePath(t.GetInput().GetPath()).ImportNode(importNode).NameNode(child.GetSecondChild()).Build())
		}
	}
}

// Java method: EsModuleProcessor.visitImportStar(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node importNode, java.lang.String moduleRequest)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitImportStar_6737(t *NodeTraversal, importNode pkg_rhino.Node, moduleRequest string) {
	_ = t
	_ = importNode
	_ = moduleRequest
	this_.currentModuleBuilder.Add_1_Import(Import_S_Builder().ModuleRequest(moduleRequest).ImportName("*").LocalName(importNode.GetSecondChild().GetString()).ImportNode(importNode).ModulePath(t.GetInput().GetPath()).NameNode(importNode.GetSecondChild()).Build())
}

// Java method: EsModuleProcessor.visitImport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node importNode)
func (this_ *EsModuleProcessor) priv_EsModuleProcessor_visitImport_6738(t *NodeTraversal, importNode pkg_rhino.Node) {
	_ = t
	_ = importNode
	var moduleRequest string = importNode.GetLastChild().GetString()
	_ = moduleRequest
	if importNode.GetFirstChild().IsName() {
		this_.priv_EsModuleProcessor_visitImportDefault_6735(t, importNode, moduleRequest)
	}
	if importNode.GetSecondChild().IsImportSpecs() {
		this_.priv_EsModuleProcessor_visitImportSpecs_6736(t, importNode, moduleRequest)
	} else if importNode.GetSecondChild().IsImportStar() {
		this_.priv_EsModuleProcessor_visitImportStar_6737(t, importNode, moduleRequest)
	}
}

// Java class: com.google.javascript.jscomp.modules.EsModuleProcessor.FindMutableExports
type EsModuleProcessor_FindMutableExports struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	exports            j2gort.Ext_java_util_List                         `gson:"exports"`
	exportsByLocalName j2gort.Ext_com_google_common_collect_ListMultimap `gson:"exportsByLocalName"`
}

func (this_ *EsModuleProcessor_FindMutableExports) BindSelf_(s any) {
	_ = s
}
func (this_ *EsModuleProcessor_FindMutableExports) IsJavaNil_() bool { return this_ == nil }
func (*EsModuleProcessor_FindMutableExports) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.EsModuleProcessor$FindMutableExports", "FindMutableExports"
}
func (this_ *EsModuleProcessor_FindMutableExports) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*EsModuleProcessor_FindMutableExports) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.EsModuleProcessor.FindMutableExports", "com.google.javascript.jscomp.modules.EsModuleProcessor$FindMutableExports", func(v_ any) bool {
		return j2gort.InstanceOf[*EsModuleProcessor_FindMutableExports](v_)
	})
}
func (*EsModuleProcessor_FindMutableExports) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*EsModuleProcessor_FindMutableExports) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *EsModuleProcessor_FindMutableExports) initFields_EsModuleProcessor_FindMutableExports() {
	_ = this_
}

// Java constructor: FindMutableExports(java.util.List<com.google.javascript.jscomp.modules.Export> exports)
func NewEsModuleProcessor_FindMutableExports(exports j2gort.Ext_java_util_List) *EsModuleProcessor_FindMutableExports {
	this_ := &EsModuleProcessor_FindMutableExports{}
	InitEsModuleProcessor_FindMutableExports(this_, exports)
	return this_
}

func InitEsModuleProcessor_FindMutableExports(this_ *EsModuleProcessor_FindMutableExports, exports j2gort.Ext_java_util_List) {
	_ = exports
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_EsModuleProcessor_FindMutableExports()
	this_.exports = exports
	this_.exportsByLocalName = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ListMultimap](j2gort.Ext_com_google_common_collect_ArrayListMultimap_Create_0())
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(exports))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			e := j2gort.CastPtr[Export](rawElem_3)
			_ = e
			j2gort.MultimapPutStr(this_.exportsByLocalName, e.LocalName(), j2gort.PtrToAny(e))
		}
	}
}

// Java method: FindMutableExports.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *EsModuleProcessor_FindMutableExports) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsName()) {
		return
	}
	var scope *Scope = t.GetScope()
	_ = scope
	if !(NodeUtil_IsLValue(n)) || j2gort.CastPtr[Scope](scope.GetClosestHoistScope()).IsModuleScope() {
		return
	}
	var exports j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MultimapGetStr(this_.exportsByLocalName, n.GetString()))
	_ = exports
	if exports.IsEmpty() {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(n.GetString()))
	_ = var_
	if (var_ == nil) || !(j2gort.CastPtr[Scope](var_.GetScope()).IsModuleScope()) {
		return
	}
	{
		cur_4 := j2gort.NewIterCursor(j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](exports)))
		for cur_4.HasNext() {
			rawElem_5 := cur_4.Next()
			e := j2gort.CastPtr[Export](rawElem_5)
			_ = e
			var i int32 = this_.exports.IndexOf_1(j2gort.PtrToAny(e))
			_ = i
			var mutated *Export = e.MutatedCopy()
			_ = mutated
			this_.exports.Set(i, j2gort.PtrToAny(mutated))
			j2gort.MultimapRemoveStr(this_.exportsByLocalName, e.LocalName(), j2gort.PtrToAny(e))
			j2gort.MultimapPutStr(this_.exportsByLocalName, e.LocalName(), j2gort.PtrToAny(mutated))
		}
	}
}

// Java class: com.google.javascript.jscomp.modules.EsModuleProcessor.UnresolvedModuleBuilder
type EsModuleProcessor_UnresolvedModuleBuilder struct {
	outer_             *EsModuleProcessor
	path               *ModuleLoader_ModulePath           `gson:"path"`
	importsByLocalName j2gort.OrderedMap[string, *Import] `gson:"importsByLocalName"`
	exports            []*Export                          `gson:"exports"`
	exportedNames      j2gort.OrderedSet[string]          `gson:"exportedNames"`
	root               pkg_rhino.Node                     `gson:"root"`
	idHash_            atomic.Int32
}

func (this_ *EsModuleProcessor_UnresolvedModuleBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *EsModuleProcessor_UnresolvedModuleBuilder) IsJavaNil_() bool { return this_ == nil }
func (*EsModuleProcessor_UnresolvedModuleBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.EsModuleProcessor$UnresolvedModuleBuilder", "UnresolvedModuleBuilder"
}
func (this_ *EsModuleProcessor_UnresolvedModuleBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*EsModuleProcessor_UnresolvedModuleBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.EsModuleProcessor.UnresolvedModuleBuilder", "com.google.javascript.jscomp.modules.EsModuleProcessor$UnresolvedModuleBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*EsModuleProcessor_UnresolvedModuleBuilder](v_)
	})
}
func (this_ *EsModuleProcessor_UnresolvedModuleBuilder) SetOuter_(o *EsModuleProcessor) {
	this_.outer_ = o
}
func (this_ *EsModuleProcessor_UnresolvedModuleBuilder) GetOuter_() *EsModuleProcessor {
	return this_.outer_
}

func (this_ *EsModuleProcessor_UnresolvedModuleBuilder) initFields_EsModuleProcessor_UnresolvedModuleBuilder() {
	_ = this_
}

// Java constructor: UnresolvedModuleBuilder(com.google.javascript.jscomp.deps.ModuleLoader.ModulePath path, com.google.javascript.rhino.Node root)
func NewEsModuleProcessor_UnresolvedModuleBuilder(outer_ *EsModuleProcessor, path *ModuleLoader_ModulePath, root pkg_rhino.Node) *EsModuleProcessor_UnresolvedModuleBuilder {
	this_ := &EsModuleProcessor_UnresolvedModuleBuilder{}
	this_.outer_ = outer_
	InitEsModuleProcessor_UnresolvedModuleBuilder(this_, path, root)
	return this_
}

func InitEsModuleProcessor_UnresolvedModuleBuilder(this_ *EsModuleProcessor_UnresolvedModuleBuilder, path *ModuleLoader_ModulePath, root pkg_rhino.Node) {
	_ = path
	_ = root
	this_.initFields_EsModuleProcessor_UnresolvedModuleBuilder()
	this_.path = path
	this_.root = root
	this_.importsByLocalName = j2gort.NewOrderedMap[string, *Import]()
	this_.exports = []*Export{}
	this_.exportedNames = j2gort.NewOrderedSet[string]()
}

// Java method: UnresolvedModuleBuilder.add(com.google.javascript.jscomp.modules.Import i)
func (this_ *EsModuleProcessor_UnresolvedModuleBuilder) Add_1_Import(i *Import) {
	_ = i
	this_.importsByLocalName.Put(i.LocalName(), i)
}

// Java method: UnresolvedModuleBuilder.add(com.google.javascript.jscomp.modules.Export e)
func (this_ *EsModuleProcessor_UnresolvedModuleBuilder) Add_1_Export(e *Export) bool {
	_ = e
	if e.ExportName() == j2gort.NullString {
		this_.exports = append(this_.exports, e)
		return true
	} else {
		if this_.exportedNames.Contains(e.ExportName()) {
			return false
		}
		this_.exports = append(this_.exports, e)
		this_.exportedNames.Add(e.ExportName())
		return true
	}
}

// Java method: UnresolvedModuleBuilder.build()
func (this_ *EsModuleProcessor_UnresolvedModuleBuilder) Build() UnresolvedModule {
	var localExports j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = localExports
	var indirectExportsBuilder []*Export = []*Export{}
	_ = indirectExportsBuilder
	var starExportsBuilder []*Export = []*Export{}
	_ = starExportsBuilder
	for _, ee := range this_.exports {
		_ = ee
		if ee.ModuleRequest() == j2gort.NullString {
			if this_.importsByLocalName.ContainsKey(ee.LocalName()) {
				indirectExportsBuilder = append(indirectExportsBuilder, j2gort.CheckNotNullElem(ee))
			} else {
				localExports.Add_1(j2gort.PtrToAny(ee))
			}
		} else {
			var iifeCond_1 bool
			{

				sr_8 := "*"
				se_9 := ee.ImportName()
				iifeCond_1 = j2gort.CheckNotNullString(sr_8) == se_9
			}
			if iifeCond_1 {
				starExportsBuilder = append(starExportsBuilder, j2gort.CheckNotNullElem(ee))
			} else {
				indirectExportsBuilder = append(indirectExportsBuilder, j2gort.CheckNotNullElem(ee))
			}
		}
	}

	NodeTraversal_Traverse(this_.outer_.compiler, this_.root, NewEsModuleProcessor_FindMutableExports(localExports))
	return NewEsModuleProcessor_UnresolvedEsModule(this_.outer_, this_.outer_.metadata, this_.path, this_.importsByLocalName.ToImmutableMap(), j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(localExports), j2gort.ImmListFromTypedSlice(indirectExportsBuilder), j2gort.ImmListFromTypedSlice(starExportsBuilder))
}

// Java class: com.google.javascript.jscomp.modules.EsModuleProcessor.UnresolvedEsModule
type EsModuleProcessor_UnresolvedEsModule struct {
	UnresolvedModuleBase
	outer_             *EsModuleProcessor
	metadata           ModuleMetadataMap_ModuleMetadata                  `gson:"metadata"`
	path               *ModuleLoader_ModulePath                          `gson:"path"`
	importsByLocalName j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"importsByLocalName"`
	localExports       []*Export                                         `gson:"localExports"`
	indirectExports    []*Export                                         `gson:"indirectExports"`
	starExports        []*Export                                         `gson:"starExports"`
	exportedNames      j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"exportedNames"`
	resolvedImports    j2gort.OrderedMap[string, *ResolveExportResult]   `gson:"resolvedImports"`
	resolvedExports    j2gort.OrderedMap[string, *ResolveExportResult]   `gson:"resolvedExports"`
	resolved           *Module                                           `gson:"resolved"`
}

func (this_ *EsModuleProcessor_UnresolvedEsModule) BindSelf_(s any) {
	_ = s
	this_.UnresolvedModuleBase.BindSelf_(s)
}
func (this_ *EsModuleProcessor_UnresolvedEsModule) IsJavaNil_() bool { return this_ == nil }
func (*EsModuleProcessor_UnresolvedEsModule) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.EsModuleProcessor$UnresolvedEsModule", "UnresolvedEsModule"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.EsModuleProcessor.UnresolvedEsModule", "com.google.javascript.jscomp.modules.EsModuleProcessor$UnresolvedEsModule", func(v_ any) bool {
		return j2gort.InstanceOf[*EsModuleProcessor_UnresolvedEsModule](v_)
	})
}
func (this_ *EsModuleProcessor_UnresolvedEsModule) SetOuter_(o *EsModuleProcessor) { this_.outer_ = o }
func (this_ *EsModuleProcessor_UnresolvedEsModule) GetOuter_() *EsModuleProcessor {
	return this_.outer_
}
func (*EsModuleProcessor_UnresolvedEsModule) Implements_jscomp_UnresolvedModule() {}

func (this_ *EsModuleProcessor_UnresolvedEsModule) initFields_EsModuleProcessor_UnresolvedEsModule() {
	_ = this_
}

// Java constructor: UnresolvedEsModule(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, com.google.javascript.jscomp.deps.ModuleLoader.ModulePath path, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.modules.Import> importsByLocalName, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.modules.Export> localExports, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.modules.Export> indirectExports, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.modules.Export> starExports)
func NewEsModuleProcessor_UnresolvedEsModule(outer_ *EsModuleProcessor, metadata ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, importsByLocalName j2gort.Ext_com_google_common_collect_ImmutableMap, localExports j2gort.Ext_com_google_common_collect_ImmutableList, indirectExports j2gort.Ext_com_google_common_collect_ImmutableList, starExports j2gort.Ext_com_google_common_collect_ImmutableList) *EsModuleProcessor_UnresolvedEsModule {
	this_ := &EsModuleProcessor_UnresolvedEsModule{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitEsModuleProcessor_UnresolvedEsModule(this_, metadata, path, importsByLocalName, localExports, indirectExports, starExports)
	return this_
}

func InitEsModuleProcessor_UnresolvedEsModule(this_ *EsModuleProcessor_UnresolvedEsModule, metadata ModuleMetadataMap_ModuleMetadata, path *ModuleLoader_ModulePath, importsByLocalName j2gort.Ext_com_google_common_collect_ImmutableMap, localExports j2gort.Ext_com_google_common_collect_ImmutableList, indirectExports j2gort.Ext_com_google_common_collect_ImmutableList, starExports j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = metadata
	_ = path
	_ = importsByLocalName
	_ = localExports
	_ = indirectExports
	_ = starExports
	InitUnresolvedModule(&this_.UnresolvedModuleBase)
	this_.initFields_EsModuleProcessor_UnresolvedEsModule()
	this_.metadata = metadata
	this_.path = path
	this_.importsByLocalName = importsByLocalName
	this_.localExports = j2gort.ImmListToSliceNullable[*Export](localExports)
	this_.indirectExports = j2gort.ImmListToSliceNullable[*Export](indirectExports)
	this_.starExports = j2gort.ImmListToSliceNullable[*Export](starExports)
	this_.exportedNames = nil
	this_.resolvedImports = j2gort.NewOrderedMap[string, *ResolveExportResult]()
	this_.resolvedExports = j2gort.NewOrderedMap[string, *ResolveExportResult]()
	this_.resolved = nil
}

// Java method: UnresolvedEsModule.reset()
func (this_ *EsModuleProcessor_UnresolvedEsModule) Reset() {
	this_.resolved = nil
	this_.exportedNames = nil
	this_.resolvedImports.Clear()
	this_.resolvedExports.Clear()
}

// Java method: UnresolvedEsModule.resolve(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.@org.jspecify.annotations.Nullable String moduleSpecifier)
func (this_ *EsModuleProcessor_UnresolvedEsModule) Resolve_2(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string) *Module {
	_ = moduleRequestResolver
	_ = moduleSpecifier
	if this_.resolved == nil {
		var boundNames j2gort.OrderedMap[string, *Binding] = j2gort.OrderedMapFromJavaMap[string, *Binding](this_.GetAllResolvedImports(moduleRequestResolver))
		_ = boundNames
		var localNameToLocalExport j2gort.OrderedMap[string, *Export] = j2gort.NewOrderedMap[string, *Export]()
		_ = localNameToLocalExport
		{

			for _, e := range this_.localExports {

				_ = e
				localNameToLocalExport.Put(e.LocalName(), e)
				{
					var iifeCond_1 bool
					{

						sr_12 := "*default*"
						se_13 := e.LocalName()
						iifeCond_1 = j2gort.CheckNotNullString(sr_12) == se_13
					}
					if !(iifeCond_1) {
						var b *ResolveExportResult = this_.ResolveExport_2(moduleRequestResolver, e.ExportName())
						_ = b
						if !(b.Resolved()) {
							j2gort.CheckStateStr(false, "Cannot have invalid missing own export!")
						}
						if !(b.IsAmbiguous()) {
							boundNames.Put(e.LocalName(), b.GetBinding())
						}
					}
				}
			}
		}
		this_.resolved = Module_S_Builder().BoundNames(boundNames.ToImmutableMap()).Namespace(j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(this_.GetAllResolvedExports(moduleRequestResolver))).Metadata(this_.metadata).Path(this_.path).LocalNameToLocalExport(localNameToLocalExport.ToImmutableMap()).Build()
	}
	return this_.resolved
}

// Java method: UnresolvedEsModule.metadata()
func (this_ *EsModuleProcessor_UnresolvedEsModule) Metadata() ModuleMetadataMap_ModuleMetadata {
	return this_.metadata
}

// Java method: UnresolvedEsModule.getAllResolvedImports(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver)
func (this_ *EsModuleProcessor_UnresolvedEsModule) GetAllResolvedImports(moduleRequestResolver ModuleRequestResolver) j2gort.Ext_java_util_Map {
	_ = moduleRequestResolver
	var imports j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = imports
	{
		cur_14 := j2gort.NewMapKeyStrCursor(this_.importsByLocalName)
		for cur_14.HasNext() {
			name := cur_14.NextStr()
			_ = name
			var b *ResolveExportResult = this_.ResolveImport_2(moduleRequestResolver, name)
			_ = b
			if b.Resolved() {
				j2gort.MapPutStr(imports, name, j2gort.PtrToAny(b.GetBinding()))
			}
		}
	}
	return imports
}

// Java method: UnresolvedEsModule.resolveImport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.String name, java.util.Set<com.google.javascript.jscomp.modules.ExportTrace> resolveSet, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> exportStarSet)
func (this_ *EsModuleProcessor_UnresolvedEsModule) ResolveImport_4(moduleRequestResolver ModuleRequestResolver, name string, resolveSet j2gort.Ext_java_util_Set, exportStarSet j2gort.Ext_java_util_Set) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = name
	_ = resolveSet
	_ = exportStarSet
	if this_.resolvedImports.ContainsKey(name) {
		return this_.resolvedImports.Get(name)
	}
	var b *ResolveExportResult = this_.priv_UnresolvedEsModule_resolveImportImpl_6745(moduleRequestResolver, name, resolveSet, exportStarSet)
	_ = b
	this_.resolvedImports.Put(name, b)
	return b
}

// Java method: UnresolvedEsModule.resolveImport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.String name)
func (this_ *EsModuleProcessor_UnresolvedEsModule) ResolveImport_2(moduleRequestResolver ModuleRequestResolver, name string) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = name
	return this_.ResolveImport_4(moduleRequestResolver, name, j2gort.NewExt_java_util_LinkedHashSet(), j2gort.NewExt_java_util_LinkedHashSet())
}

// Java method: UnresolvedEsModule.resolveImportImpl(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.String name, java.util.Set<com.google.javascript.jscomp.modules.ExportTrace> resolveSet, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> exportStarSet)
func (this_ *EsModuleProcessor_UnresolvedEsModule) priv_UnresolvedEsModule_resolveImportImpl_6745(moduleRequestResolver ModuleRequestResolver, name string, resolveSet j2gort.Ext_java_util_Set, exportStarSet j2gort.Ext_java_util_Set) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = name
	_ = resolveSet
	_ = exportStarSet
	var i *Import = j2gort.CastPtr[Import](j2gort.MapGetStr(this_.importsByLocalName, name))
	_ = i
	var requested UnresolvedModule = moduleRequestResolver.Resolve_1_Import(i)
	_ = requested
	if requested == nil || requested.IsJavaNil_() {
		return StaticGet_ResolveExportResult_ERROR()
	} else {
		var importStar bool = (j2gort.CheckNotNullString(i.ImportName()) == "*")
		_ = importStar
		if importStar || (func() bool {
			sr_15 := i.ImportName()
			se_16 := "default"
			return j2gort.CheckNotNullString(sr_15) == se_16
		}() && (requested.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsGoogProvide() || requested.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsGoogModule())) {
			if !(GoogEsImports_IsGoogImportSpecifier(i.ModuleRequest())) && (requested.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsGoogModule() || requested.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsGoogProvide()) {
				this_.outer_.compiler.Report_1(JSError_Make_5(this_.path.String(), i.ImportNode().GetLineno(), i.ImportNode().GetCharno(), StaticGet_EsModuleProcessor_CANNOT_PATH_IMPORT_CLOSURE_FILE(), []string{i.LocalName(), i.ModuleRequest()}))
				return StaticGet_ResolveExportResult_ERROR()
			}
			if importStar && GoogEsImports_IsGoogImportSpecifier(i.ModuleRequest()) {
				this_.outer_.compiler.Report_1(JSError_Make_5(this_.path.String(), i.ImportNode().GetLineno(), i.ImportNode().GetCharno(), StaticGet_EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR(), []string{i.LocalName(), i.ModuleRequest()}))
				return StaticGet_ResolveExportResult_ERROR()
			}
			var closureNamespace string
			_ = closureNamespace
			if GoogEsImports_IsGoogImportSpecifier(i.ModuleRequest()) {
				closureNamespace = GoogEsImports_GetClosureIdFromGoogImportSpecifier(i.ModuleRequest())
			} else {
				closureNamespace = j2gort.NullString
			}
			return ResolveExportResult_Of(Binding_From_3(requested.Metadata(), closureNamespace, i.NameNode()))
		} else {
			var result *ResolveExportResult = requested.ResolveExport_5(moduleRequestResolver, i.ModuleRequest(), i.ImportName(), resolveSet, exportStarSet)
			_ = result
			if !(result.Found()) && !(result.HadError()) {
				this_.outer_.compiler.Report_1(JSError_Make_5(this_.path.String(), i.ImportNode().GetLineno(), i.ImportNode().GetCharno(), StaticGet_ModuleMapCreator_DOES_NOT_HAVE_EXPORT(), []string{i.ImportName()}))
				return StaticGet_ResolveExportResult_ERROR()
			} else if result.IsAmbiguous() {
				this_.outer_.compiler.Report_1(JSError_Make_5(this_.path.String(), i.ImportNode().GetLineno(), i.ImportNode().GetCharno(), StaticGet_EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT(), []string{i.ImportName()}))
				return StaticGet_ResolveExportResult_ERROR()
			}
			var forSourceInfo pkg_rhino.Node
			_ = forSourceInfo
			if i.NameNode() == 0 {
				forSourceInfo = i.ImportNode()
			} else {
				forSourceInfo = i.NameNode()
			}
			return result.Copy(forSourceInfo, StaticGet_Binding_CreatedBy_IMPORT())
		}
	}
}

// Java method: UnresolvedEsModule.getExportedNames(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver)
func (this_ *EsModuleProcessor_UnresolvedEsModule) GetExportedNames_1(moduleRequestResolver ModuleRequestResolver) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = moduleRequestResolver
	if this_.exportedNames == nil {
		this_.exportedNames = this_.GetExportedNames_2(moduleRequestResolver, j2gort.NewExt_java_util_LinkedHashSet())
	}
	return this_.exportedNames
}

// Java method: UnresolvedEsModule.getExportedNames(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> visited)
func (this_ *EsModuleProcessor_UnresolvedEsModule) GetExportedNames_2(moduleRequestResolver ModuleRequestResolver, visited j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = moduleRequestResolver
	_ = visited
	if visited.Contains(this_) {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	}
	visited.Add_1(this_)
	var exportedNames j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet()
	_ = exportedNames
	{

		for _, e := range this_.localExports {

			_ = e
			j2gort.SetAddStr(exportedNames, e.ExportName())
		}
	}
	{

		for _, e_2 := range this_.indirectExports {

			_ = e_2
			j2gort.SetAddStr(exportedNames, e_2.ExportName())
		}
	}
	{

		for _, e_3 := range this_.starExports {

			_ = e_3
			var requested UnresolvedModule = moduleRequestResolver.Resolve_1_Export(e_3)
			_ = requested
			if requested != nil && !requested.IsJavaNil_() {
				if requested.Metadata().GetModuleMetadataMap_ModuleMetadataBase().IsEs6Module() {
					{
						cur_23 := j2gort.NewIterCursor(j2gort.CheckNotNull(requested.GetExportedNames_2(moduleRequestResolver, visited)))
						for cur_23.HasNext() {
							n := cur_23.NextStr()
							_ = n
							if !("default" == n) && !(j2gort.SetContainsStr(exportedNames, n)) {
								j2gort.SetAddStr(exportedNames, n)
							}
						}
					}
				} else {
					this_.outer_.compiler.Report_1(JSError_Make_3(e_3.ExportNode(), StaticGet_TranspilationUtil_CANNOT_CONVERT_YET(), []string{"Wildcard export for non-ES module"}))
				}
			}
		}
	}
	return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(exportedNames)
}

// Java method: UnresolvedEsModule.getAllResolvedExports(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver)
func (this_ *EsModuleProcessor_UnresolvedEsModule) GetAllResolvedExports(moduleRequestResolver ModuleRequestResolver) j2gort.Ext_java_util_Map {
	_ = moduleRequestResolver
	var exports j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = exports
	{
		cur_24 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetExportedNames_1(moduleRequestResolver)))
		for cur_24.HasNext() {
			name := cur_24.NextStr()
			_ = name
			var b *ResolveExportResult = this_.ResolveExport_2(moduleRequestResolver, name)
			_ = b
			if !(b.Found()) {
				j2gort.CheckStateStr(false, "Cannot have invalid own export.")
			}
			if b.Resolved() {
				j2gort.MapPutStr(exports, name, j2gort.PtrToAny(b.GetBinding()))
			} else if b.IsAmbiguous() {
				this_.outer_.compiler.Report_1(JSError_Make_5(this_.path.String(), int32(-1), int32(-1), StaticGet_EsModuleProcessor_AMBIGUOUS_EXPORT_DEFINITION(), []string{name}))
			}
		}
	}
	return exports
}

// Java method: UnresolvedEsModule.resolveExport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.String exportName, java.util.Set<com.google.javascript.jscomp.modules.ExportTrace> resolveSet, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> exportStarSet)
func (this_ *EsModuleProcessor_UnresolvedEsModule) priv_UnresolvedEsModule_resolveExport_6747(moduleRequestResolver ModuleRequestResolver, exportName string, resolveSet j2gort.Ext_java_util_Set, exportStarSet j2gort.Ext_java_util_Set) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = exportName
	_ = resolveSet
	_ = exportStarSet
	if !(j2gort.SetContainsStr(this_.GetExportedNames_1(moduleRequestResolver), exportName)) {
		return StaticGet_ResolveExportResult_NOT_FOUND()
	}
	if !(resolveSet.Add_1(j2gort.PtrToAny(ExportTrace_Create(this_, exportName)))) {
		this_.outer_.compiler.Report_1(JSError_Make_5(this_.path.String(), int32(0), int32(0), StaticGet_EsModuleProcessor_CYCLIC_EXPORT_DEFINITION(), []string{exportName}))
		return StaticGet_ResolveExportResult_ERROR()
	}
	{

		for _, e := range this_.localExports {

			_ = e
			{
				var iifeCond_1 bool
				{

					sr_27 := exportName
					se_28 := e.ExportName()
					iifeCond_1 = j2gort.CheckNotNullString(sr_27) == se_28
				}
				if iifeCond_1 {
					var forSourceInfo pkg_rhino.Node
					_ = forSourceInfo
					if e.NameNode() != 0 {
						forSourceInfo = e.NameNode()
					} else {
						forSourceInfo = e.ExportNode()
					}
					return ResolveExportResult_Of(Binding_From_2(e, forSourceInfo))
				}
			}
		}
	}
	{

		for _, e_2 := range this_.indirectExports {

			_ = e_2
			{
				var iifeCond_1 bool
				{

					sr_31 := exportName
					se_32 := e_2.ExportName()
					iifeCond_1 = j2gort.CheckNotNullString(sr_31) == se_32
				}
				if iifeCond_1 {
					if j2gort.MapContainsKeyStr(this_.importsByLocalName, e_2.LocalName()) {
						return this_.ResolveImport_4(moduleRequestResolver, e_2.LocalName(), resolveSet, exportStarSet).Copy(e_2.NameNode(), StaticGet_Binding_CreatedBy_EXPORT())
					} else {
						var requested UnresolvedModule = moduleRequestResolver.Resolve_1_Export(e_2)
						_ = requested
						if requested == nil || requested.IsJavaNil_() {
							return StaticGet_ResolveExportResult_ERROR()
						} else {
							var result *ResolveExportResult = requested.ResolveExport_5(moduleRequestResolver, e_2.ModuleRequest(), e_2.ImportName(), resolveSet, exportStarSet)
							_ = result
							if !(result.Found()) && !(result.HadError()) {
								this_.outer_.compiler.Report_1(JSError_Make_5(this_.path.String(), e_2.ExportNode().GetLineno(), e_2.ExportNode().GetCharno(), StaticGet_ModuleMapCreator_DOES_NOT_HAVE_EXPORT(), []string{e_2.ImportName()}))
								return StaticGet_ResolveExportResult_ERROR()
							} else if result.IsAmbiguous() {
								this_.outer_.compiler.Report_1(JSError_Make_5(this_.path.String(), e_2.ExportNode().GetLineno(), e_2.ExportNode().GetCharno(), StaticGet_EsModuleProcessor_IMPORTED_AMBIGUOUS_EXPORT(), []string{e_2.ImportName()}))
							}
							return result.Copy(e_2.NameNode(), StaticGet_Binding_CreatedBy_EXPORT())
						}
					}
				}
			}
		}
	}
	if !(!("default" == exportName)) {
		j2gort.CheckStateStr(false, "Default export cannot come from export *.")
	}
	if exportStarSet.Contains(this_) {
		this_.outer_.compiler.Report_1(JSError_Make_5(this_.path.String(), int32(-1), int32(-1), StaticGet_EsModuleProcessor_CYCLIC_EXPORT_DEFINITION(), []string{exportName}))
		return StaticGet_ResolveExportResult_ERROR()
	}
	exportStarSet.Add_1(this_)
	var starResolution *ResolveExportResult = nil
	_ = starResolution
	{

		for _, e_3 := range this_.starExports {

			_ = e_3
			var requested_2 UnresolvedModule = moduleRequestResolver.Resolve_1_Export(e_3)
			_ = requested_2
			if requested_2 == nil || requested_2.IsJavaNil_() {
				return StaticGet_ResolveExportResult_ERROR()
			} else if j2gort.SetContainsStr(requested_2.GetExportedNames_1(moduleRequestResolver), exportName) {
				var resolution *ResolveExportResult = requested_2.ResolveExport_5(moduleRequestResolver, e_3.ModuleRequest(), exportName, resolveSet, exportStarSet)
				_ = resolution
				if resolution.HadError() {
					return resolution
				} else if resolution.IsAmbiguous() {
					return resolution
				} else {
					if starResolution == nil {
						starResolution = resolution.Copy(e_3.ExportNode(), StaticGet_Binding_CreatedBy_EXPORT())
					} else {
						if starResolution != resolution {
							return StaticGet_ResolveExportResult_AMBIGUOUS()
						}
					}
				}
			}
		}
	}
	if starResolution == nil {
		return StaticGet_ResolveExportResult_ERROR()
	}
	return starResolution
}

// Java method: UnresolvedEsModule.resolveExport(com.google.javascript.jscomp.modules.ModuleRequestResolver moduleRequestResolver, java.lang.@org.jspecify.annotations.Nullable String moduleSpecifier, java.lang.String exportName, java.util.Set<com.google.javascript.jscomp.modules.ExportTrace> resolveSet, java.util.Set<com.google.javascript.jscomp.modules.UnresolvedModule> exportStarSet)
func (this_ *EsModuleProcessor_UnresolvedEsModule) ResolveExport_5(moduleRequestResolver ModuleRequestResolver, moduleSpecifier string, exportName string, resolveSet j2gort.Ext_java_util_Set, exportStarSet j2gort.Ext_java_util_Set) *ResolveExportResult {
	_ = moduleRequestResolver
	_ = moduleSpecifier
	_ = exportName
	_ = resolveSet
	_ = exportStarSet
	if this_.resolvedExports.ContainsKey(exportName) {
		return this_.resolvedExports.Get(exportName)
	}
	var b *ResolveExportResult = this_.priv_UnresolvedEsModule_resolveExport_6747(moduleRequestResolver, exportName, resolveSet, exportStarSet)
	_ = b
	this_.resolvedExports.Put(exportName, b)
	return b
}
func (*EsModuleProcessor) CanReleaseSyntacticScopes_()                    {}
func (*EsModuleProcessor_FindMutableExports) AlwaysShouldTraverse_()      {}
func (*EsModuleProcessor_FindMutableExports) CanReleaseSyntacticScopes_() {}

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
// Port of com/google/javascript/jscomp/Es6RewriteModulesToCommonJsModules.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules
type Es6RewriteModulesToCommonJsModules struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *Es6RewriteModulesToCommonJsModules) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteModulesToCommonJsModules) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteModulesToCommonJsModules) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules", "Es6RewriteModulesToCommonJsModules"
}
func (this_ *Es6RewriteModulesToCommonJsModules) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteModulesToCommonJsModules) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules", "com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteModulesToCommonJsModules](v_)
	})
}
func (*Es6RewriteModulesToCommonJsModules) Implements_jscomp_CompilerPass() {}

func StaticGet_Es6RewriteModulesToCommonJsModules_JSCOMP_DEFAULT_EXPORT() string { return "$$default" }
func StaticGet_Es6RewriteModulesToCommonJsModules_MODULE() string                { return "$$module" }
func StaticGet_Es6RewriteModulesToCommonJsModules_EXPORTS() string               { return "$$exports" }
func StaticGet_Es6RewriteModulesToCommonJsModules_REQUIRE() string               { return "$$require" }
func (this_ *Es6RewriteModulesToCommonJsModules) initFields_Es6RewriteModulesToCommonJsModules() {
	_ = this_
}

// Java constructor: Es6RewriteModulesToCommonJsModules(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6RewriteModulesToCommonJsModules(compiler AbstractCompiler) *Es6RewriteModulesToCommonJsModules {
	this_ := &Es6RewriteModulesToCommonJsModules{}
	InitEs6RewriteModulesToCommonJsModules(this_, compiler)
	return this_
}

func InitEs6RewriteModulesToCommonJsModules(this_ *Es6RewriteModulesToCommonJsModules, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Es6RewriteModulesToCommonJsModules()
	this_.compiler = compiler
}

// Java method: Es6RewriteModulesToCommonJsModules.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteModulesToCommonJsModules) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		var script pkg_rhino.Node = root.GetFirstChild()
		_ = script
		for ; script != 0; script = script.GetNext() {
			if Es6RewriteModules_IsEs6ModuleRoot(script) {
				NodeTraversal_Traverse(this_.compiler, script, NewEs6RewriteModulesToCommonJsModules_Rewriter(this_.compiler, script))
			}
		}
	}
	NodeUtil_RemoveFeatureFromAllScripts(root, pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), this_.compiler)
	GatherGetterAndSetterProperties_Update(this_.compiler, externs, root)
}

// Java static method: Es6RewriteModulesToCommonJsModules.normalizePath(java.lang.String path)
func Es6RewriteModulesToCommonJsModules_NormalizePath(path string) string {
	_ = path
	var indexOfProtocol int32 = j2gort.JStrIndexOf(path, "://")
	_ = indexOfProtocol
	if indexOfProtocol > int32(-1) {
		path = j2gort.JStrSubFrom(path, (indexOfProtocol + int32(3)))
		var indexOfSlash int32 = j2gort.JStrIndexOfCp(path, int32(uint16(47)))
		_ = indexOfSlash
		if indexOfSlash > int32(-1) {
			path = j2gort.JStrSubFrom(path, (indexOfSlash + int32(1)))
		}
	} else if javalang.StringStartsWith(j2gort.CheckNotNullString(path), "/") {
		path = j2gort.JStrSubFrom(path, int32(1))
	}
	return path
}

// Java class: com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules.LocalQName
type Es6RewriteModulesToCommonJsModules_LocalQName struct {
	qName             string         `gson:"qName"`
	nodeForSourceInfo pkg_rhino.Node `gson:"nodeForSourceInfo"`
	idHash_           atomic.Int32
}

func (this_ *Es6RewriteModulesToCommonJsModules_LocalQName) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteModulesToCommonJsModules_LocalQName) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteModulesToCommonJsModules_LocalQName) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules$LocalQName", "LocalQName"
}
func (this_ *Es6RewriteModulesToCommonJsModules_LocalQName) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteModulesToCommonJsModules_LocalQName) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules.LocalQName", "com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules$LocalQName", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteModulesToCommonJsModules_LocalQName](v_)
	})
}

func (this_ *Es6RewriteModulesToCommonJsModules_LocalQName) initFields_Es6RewriteModulesToCommonJsModules_LocalQName() {
	_ = this_
	this_.qName = j2gort.NullString
}

// Java constructor: LocalQName(java.lang.String qName, com.google.javascript.rhino.Node nodeForSourceInfo)
func NewEs6RewriteModulesToCommonJsModules_LocalQName(qName string, nodeForSourceInfo pkg_rhino.Node) *Es6RewriteModulesToCommonJsModules_LocalQName {
	this_ := &Es6RewriteModulesToCommonJsModules_LocalQName{}
	InitEs6RewriteModulesToCommonJsModules_LocalQName(this_, qName, nodeForSourceInfo)
	return this_
}

func InitEs6RewriteModulesToCommonJsModules_LocalQName(this_ *Es6RewriteModulesToCommonJsModules_LocalQName, qName string, nodeForSourceInfo pkg_rhino.Node) {
	_ = qName
	_ = nodeForSourceInfo
	this_.initFields_Es6RewriteModulesToCommonJsModules_LocalQName()
	this_.qName = qName
	this_.nodeForSourceInfo = nodeForSourceInfo
}

// Java record: com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules.ModuleRequest
type Es6RewriteModulesToCommonJsModules_ModuleRequest struct {
	specifier string `gson:"specifier"`
	varName   string `gson:"varName"`
}

func (this_ *Es6RewriteModulesToCommonJsModules_ModuleRequest) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteModulesToCommonJsModules_ModuleRequest) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteModulesToCommonJsModules_ModuleRequest) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules$ModuleRequest", "ModuleRequest"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules.ModuleRequest", "com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules$ModuleRequest", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteModulesToCommonJsModules_ModuleRequest](v_)
	})
}
func (*Es6RewriteModulesToCommonJsModules_ModuleRequest) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *Es6RewriteModulesToCommonJsModules_ModuleRequest) initFields_Es6RewriteModulesToCommonJsModules_ModuleRequest() {
	_ = this_
}

// Java constructor: ModuleRequest(java.lang.String specifier, java.lang.String varName)
func NewEs6RewriteModulesToCommonJsModules_ModuleRequest(specifier string, varName string) *Es6RewriteModulesToCommonJsModules_ModuleRequest {
	this_ := &Es6RewriteModulesToCommonJsModules_ModuleRequest{}
	InitEs6RewriteModulesToCommonJsModules_ModuleRequest(this_, specifier, varName)
	return this_
}

func InitEs6RewriteModulesToCommonJsModules_ModuleRequest(this_ *Es6RewriteModulesToCommonJsModules_ModuleRequest, specifier string, varName string) {
	_ = specifier
	_ = varName
	this_.initFields_Es6RewriteModulesToCommonJsModules_ModuleRequest()
	j2gort.CheckNotNullString(specifier)
	j2gort.CheckNotNullString(varName)
	this_.specifier = specifier
	this_.varName = varName
}

// Java static method: ModuleRequest.create(java.lang.String specifier, java.lang.String varName)
func Es6RewriteModulesToCommonJsModules_ModuleRequest_Create(specifier string, varName string) *Es6RewriteModulesToCommonJsModules_ModuleRequest {
	_ = specifier
	_ = varName
	return NewEs6RewriteModulesToCommonJsModules_ModuleRequest(specifier, varName)
}

func (this_ *Es6RewriteModulesToCommonJsModules_ModuleRequest) Specifier() string {
	return this_.specifier
}
func (this_ *Es6RewriteModulesToCommonJsModules_ModuleRequest) VarName() string { return this_.varName }
func (this_ *Es6RewriteModulesToCommonJsModules_ModuleRequest) Equals(other_ any) bool {
	o_, ok_ := other_.(*Es6RewriteModulesToCommonJsModules_ModuleRequest)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.specifier != o_.specifier {
		return false
	}
	if this_.varName != o_.varName {
		return false
	}
	return true
}
func (this_ *Es6RewriteModulesToCommonJsModules_ModuleRequest) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_1 := this_.specifier
		h_ = 31 * h_
		if hs_1 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_1)
		}
	}
	{
		hs_2 := this_.varName
		h_ = 31 * h_
		if hs_2 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_2)
		}
	}
	return h_
}
func (this_ *Es6RewriteModulesToCommonJsModules_ModuleRequest) String() string {
	_ = fmt.Sprint
	return "ModuleRequest[" + "specifier=" + j2gort.StrVal(j2gort.BoxString(this_.specifier)) + ", " + "varName=" + j2gort.StrVal(j2gort.BoxString(this_.varName)) + "]"
}

// Java class: com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules.Rewriter
type Es6RewriteModulesToCommonJsModules_Rewriter struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	exportedNameToLocalQName j2gort.Ext_java_util_Map                                                   `gson:"exportedNameToLocalQName"`
	imports                  j2gort.OrderedSet[pkg_rhino.Node]                                          `gson:"imports"`
	importRequests           j2gort.CustomOrderedSet[*Es6RewriteModulesToCommonJsModules_ModuleRequest] `gson:"importRequests"`
	compiler                 AbstractCompiler                                                           `gson:"compiler"`
	modulePath               *ModuleLoader_ModulePath                                                   `gson:"modulePath"`
	requireInsertSpot        pkg_rhino.Node                                                             `gson:"requireInsertSpot"`
	script                   pkg_rhino.Node                                                             `gson:"script"`
}

func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteModulesToCommonJsModules_Rewriter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules$Rewriter", "Rewriter"
}
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteModulesToCommonJsModules_Rewriter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules.Rewriter", "com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules$Rewriter", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteModulesToCommonJsModules_Rewriter](v_)
	})
}
func (*Es6RewriteModulesToCommonJsModules_Rewriter) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*Es6RewriteModulesToCommonJsModules_Rewriter) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) initFields_Es6RewriteModulesToCommonJsModules_Rewriter() {
	_ = this_
}

// Java constructor: Rewriter(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node script)
func NewEs6RewriteModulesToCommonJsModules_Rewriter(compiler AbstractCompiler, script pkg_rhino.Node) *Es6RewriteModulesToCommonJsModules_Rewriter {
	this_ := &Es6RewriteModulesToCommonJsModules_Rewriter{}
	InitEs6RewriteModulesToCommonJsModules_Rewriter(this_, compiler, script)
	return this_
}

func InitEs6RewriteModulesToCommonJsModules_Rewriter(this_ *Es6RewriteModulesToCommonJsModules_Rewriter, compiler AbstractCompiler, script pkg_rhino.Node) {
	_ = compiler
	_ = script
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6RewriteModulesToCommonJsModules_Rewriter()
	this_.compiler = compiler
	this_.script = script
	this_.requireInsertSpot = pkg_rhino.Node(0)
	this_.exportedNameToLocalQName = j2gort.NewExt_java_util_TreeMap()
	this_.importRequests = j2gort.NewCustomOrderedSet[*Es6RewriteModulesToCommonJsModules_ModuleRequest]()
	this_.imports = j2gort.NewOrderedSet[pkg_rhino.Node]()
	this_.modulePath = compiler.GetInput(script.GetInputId()).GetPath()
}

// Java method: Rewriter.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(116) /* Token.IMPORT */ :
		this_.priv_Rewriter_visitImport_2766(t.GetInput().GetPath(), n)
	case int32(120) /* Token.EXPORT */ :
		this_.priv_Rewriter_visitExport_2774(t, n, parent)
	case int32(94) /* Token.SCRIPT */ :
		this_.priv_Rewriter_visitScript_2761(t, n)
	case int32(32) /* Token.NAME */ :
		this_.priv_Rewriter_maybeRenameImportedValue_2760(t, n)
	default:
	}
}

// Java method: Rewriter.getVarNameOfImport(com.google.javascript.rhino.Node importDecl)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_getVarNameOfImport_2756(importDecl pkg_rhino.Node) string {
	_ = importDecl
	if !(importDecl.IsImport()) {
		j2gort.CheckState(false)
	}
	if importDecl.GetSecondChild().IsImportStar() {
		return importDecl.GetSecondChild().GetString()
	}
	return this_.priv_Rewriter_getVarNameOfImport_2757(importDecl.GetLastChild().GetString())
}

// Java method: Rewriter.getVarNameOfImport(java.lang.String importRequest)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_getVarNameOfImport_2757(importRequest string) string {
	_ = importRequest
	return this_.modulePath.ResolveModuleAsPath(importRequest).ToModuleName()
}

// Java method: Rewriter.getNameOfImportedValue(com.google.javascript.rhino.Node nameNode)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_getNameOfImportedValue_2758(nameNode pkg_rhino.Node) string {
	_ = nameNode
	var importDecl pkg_rhino.Node = nameNode
	_ = importDecl
	for !(importDecl.IsImport()) {
		importDecl = importDecl.GetParent()
	}
	var moduleName string = this_.priv_Rewriter_getVarNameOfImport_2756(importDecl)
	_ = moduleName
	if nameNode.GetParent().IsImportSpec() {
		return (j2gort.StrValString(moduleName) + "." + j2gort.StrValString(nameNode.GetParent().GetFirstChild().GetString()))
	} else if nameNode.IsImportStar() {
		return moduleName
	} else {
		if !(nameNode.GetParent().IsImport()) {
			j2gort.CheckState(false)
		}
		return (j2gort.StrValString(moduleName) + ".default")
	}
}

// Java method: Rewriter.maybeGetNameOfImportedValue(com.google.javascript.jscomp.Scope s, com.google.javascript.rhino.Node nameNode)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_maybeGetNameOfImportedValue_2759(s *Scope, nameNode pkg_rhino.Node) string {
	_ = s
	_ = nameNode
	if !(nameNode.IsName()) {
		j2gort.CheckState(false)
	}
	var var_ *Var = j2gort.CastPtr[Var](s.GetVar(nameNode.GetString()))
	_ = var_
	if (((var_ != nil) && (var_.GetNameNode() != 0)) && NodeUtil_IsImportedName(var_.GetNameNode())) && (nameNode != var_.GetNameNode()) {
		return this_.priv_Rewriter_getNameOfImportedValue_2758(var_.GetNameNode())
	}
	return j2gort.NullString
}

// Java method: Rewriter.maybeRenameImportedValue(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_maybeRenameImportedValue_2760(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if !(n.IsName()) {
		j2gort.CheckState(false)
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if ((parent.IsExport() || parent.IsExportSpec()) || parent.IsImport()) || parent.IsImportSpec() {
		return
	}
	var qName string = this_.priv_Rewriter_maybeGetNameOfImportedValue_2759(t.GetScope(), n)
	_ = qName
	if qName != j2gort.NullString {
		n.ReplaceWith(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, qName))
		t.ReportCodeChange_0()
	}
}

// Java method: Rewriter.visitScript(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node script)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_visitScript_2761(t *NodeTraversal, script pkg_rhino.Node) {
	_ = t
	_ = script
	if !(this_.script == script) {
		j2gort.CheckState(false)
	}
	var moduleNode pkg_rhino.Node = script.GetFirstChild()
	_ = moduleNode
	if !(moduleNode.IsModuleBody()) {
		j2gort.CheckState(false)
	}
	moduleNode.Detach()
	script.AddChildrenToFront(moduleNode.RemoveChildren())
	this_.priv_Rewriter_addRequireCalls_2762()
	this_.priv_Rewriter_addExportDef_2764()
	this_.priv_Rewriter_registerAndLoadModule_2763(t)
}

// Java method: Rewriter.addRequireCalls()
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_addRequireCalls_2762() {
	if !(this_.importRequests.IsEmpty()) {
		{
			cur_1 := this_.imports.Iter()
			for cur_1.HasNext() {
				importDecl := cur_1.
					Next()
				_ = importDecl
				importDecl.Detach()
			}
		}
		var importedNames j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
		_ = importedNames
		{
			cur_2 := this_.importRequests.Iter()
			for cur_2.HasNext() {

				request := cur_2.
					Next()
				_ = request
				var varName string = request.VarName()
				_ = varName
				if importedNames.Add(varName) {
					var requireCall pkg_rhino.Node = pkg_rhino.IR_Call(pkg_rhino.IR_Name("$$require"), []pkg_rhino.Node{pkg_rhino.IR_String(request.Specifier())})
					_ = requireCall
					requireCall.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
					var decl pkg_rhino.Node = pkg_rhino.IR_Var_2(pkg_rhino.IR_Name(varName), requireCall)
					_ = decl
					decl.SrcrefTreeIfMissing(this_.script)
					if this_.requireInsertSpot == 0 {
						this_.script.AddChildToFront(decl)
					} else {
						decl.InsertAfter(this_.requireInsertSpot)
					}
					this_.requireInsertSpot = decl
				}
			}
		}
	}
}

// Java method: Rewriter.registerAndLoadModule(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_registerAndLoadModule_2763(t *NodeTraversal) {
	_ = t
	var block pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = block
	block.AddChildrenToFront(this_.script.RemoveChildren())
	var moduleFunction pkg_rhino.Node = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{pkg_rhino.IR_Name("$$require"), pkg_rhino.IR_Name("$$exports"), pkg_rhino.IR_Name("$$module")}), block)
	_ = moduleFunction
	var shallowDeps pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_ARRAYLIT())
	_ = shallowDeps
	{
		cur_4 := this_.importRequests.Iter()
		for cur_4.HasNext() {

			request := cur_4.
				Next()
			_ = request
			shallowDeps.AddChildToBack(pkg_rhino.IR_String(request.Specifier()))
		}
	}
	var exprResult pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("$jscomp"), "registerAndLoadModule"), []pkg_rhino.Node{moduleFunction, pkg_rhino.IR_String(Es6RewriteModulesToCommonJsModules_NormalizePath(this_.compiler.GetModuleLoader().Resolve(t.GetInput().GetName()).String())), shallowDeps}))
	_ = exprResult
	this_.script.AddChildToBack(exprResult.SrcrefTreeIfMissing(this_.script))
	this_.compiler.ReportChangeToChangeScope(this_.script)
	this_.compiler.ReportChangeToChangeScope(moduleFunction)
	t.ReportCodeChange_0()
}

// Java method: Rewriter.addExportDef()
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_addExportDef_2764() {
	if !(this_.exportedNameToLocalQName.IsEmpty()) {
		var definePropertiesLit pkg_rhino.Node = pkg_rhino.IR_Objectlit([]pkg_rhino.Node{})
		_ = definePropertiesLit
		{
			cur_6 := j2gort.NewMapEntryCursor(this_.exportedNameToLocalQName)
			for cur_6.NextEntry() {
				this_.priv_Rewriter_addExport_2765(definePropertiesLit, cur_6.KeyStr(), j2gort.CastPtr[Es6RewriteModulesToCommonJsModules_LocalQName](cur_6.GetValue()))
			}
		}
		this_.script.AddChildToFront(pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, "Object.defineProperties"), []pkg_rhino.Node{pkg_rhino.IR_Name("$$exports"), definePropertiesLit})).SrcrefTreeIfMissing(this_.script))
	}
}

// Java method: Rewriter.addExport(com.google.javascript.rhino.Node definePropertiesLit, java.lang.String exportedName, com.google.javascript.jscomp.Es6RewriteModulesToCommonJsModules.LocalQName localQName)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_addExport_2765(definePropertiesLit pkg_rhino.Node, exportedName string, localQName *Es6RewriteModulesToCommonJsModules_LocalQName) {
	_ = definePropertiesLit
	_ = exportedName
	_ = localQName
	var exportedValue pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, localQName.qName)
	_ = exportedValue
	var getterFunction pkg_rhino.Node = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{}), pkg_rhino.IR_Block_1_Node(pkg_rhino.IR_ReturnNode_1(exportedValue)))
	_ = getterFunction
	getterFunction.SrcrefTree(localQName.nodeForSourceInfo)
	var objLit pkg_rhino.Node = pkg_rhino.IR_Objectlit([]pkg_rhino.Node{pkg_rhino.IR_StringKey_2("enumerable", pkg_rhino.IR_TrueNode()), pkg_rhino.IR_StringKey_2("get", getterFunction)})
	_ = objLit
	definePropertiesLit.AddChildToBack(pkg_rhino.IR_StringKey_2(exportedName, objLit))
	this_.compiler.ReportChangeToChangeScope(getterFunction)
}

// Java method: Rewriter.visitImport(com.google.javascript.jscomp.deps.ModuleLoader.ModulePath path, com.google.javascript.rhino.Node importDecl)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_visitImport_2766(path *ModuleLoader_ModulePath, importDecl pkg_rhino.Node) {
	_ = path
	_ = importDecl
	if j2gort.JStrIndexOf(importDecl.GetLastChild().GetString(), "://") >= 0 {
		this_.compiler.Report_1(JSError_Make_3(importDecl, StaticGet_TranspilationUtil_CANNOT_CONVERT(), []string{"Module requests with protocols."}))
	}
	this_.importRequests.Add(Es6RewriteModulesToCommonJsModules_ModuleRequest_Create(Es6RewriteModulesToCommonJsModules_NormalizePath(path.ResolveModuleAsPath(importDecl.GetLastChild().GetString()).String()), this_.priv_Rewriter_getVarNameOfImport_2756(importDecl)))
	this_.imports.Add(importDecl)
}

// Java method: Rewriter.visitExportDefault(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_visitExportDefault_2767(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
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
		name = "$$default"
		var var_ pkg_rhino.Node = pkg_rhino.IR_ConstNode(pkg_rhino.IR_Name(name), export.RemoveFirstChild())
		_ = var_
		export.ReplaceWith(var_.SrcrefTreeIfMissing(export))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
	}
	j2gort.MapPutStr(this_.exportedNameToLocalQName, "default", NewEs6RewriteModulesToCommonJsModules_LocalQName(name, export))
	t.ReportCodeChange_0()
}

// Java method: Rewriter.visitExportFrom(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_visitExportFrom_2768(t *NodeTraversal, export pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = export
	_ = parent
	var moduleIdentifier pkg_rhino.Node = export.GetLastChild()
	_ = moduleIdentifier
	var importNode pkg_rhino.Node = pkg_rhino.IR_ImportNode(pkg_rhino.IR_Empty(), pkg_rhino.IR_Empty(), moduleIdentifier.CloneNode_0())
	_ = importNode
	importNode.Srcref(export)
	importNode.InsertBefore(export)
	this_.Visit_3(t, importNode, parent)
	var moduleName string = this_.priv_Rewriter_getVarNameOfImport_2757(moduleIdentifier.GetString())
	_ = moduleName
	{
		var exportSpec pkg_rhino.Node = export.GetFirstFirstChild()
		_ = exportSpec
		for ; exportSpec != 0; exportSpec = exportSpec.GetNext() {
			j2gort.MapPutStr(this_.exportedNameToLocalQName, exportSpec.GetLastChild().GetString(), NewEs6RewriteModulesToCommonJsModules_LocalQName((j2gort.StrValString(moduleName)+"."+j2gort.StrValString(exportSpec.GetFirstChild().GetString())), exportSpec))
		}
	}
	export.Detach()
	t.ReportCodeChange_0()
}

// Java method: Rewriter.visitExportSpecs(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_visitExportSpecs_2769(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	{
		var exportSpec pkg_rhino.Node = export.GetFirstFirstChild()
		_ = exportSpec
		for ; exportSpec != 0; exportSpec = exportSpec.GetNext() {
			var localName string = exportSpec.GetFirstChild().GetString()
			_ = localName
			var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(localName))
			_ = var_
			if (var_ != nil) && NodeUtil_IsImportedName(var_.GetNameNode()) {
				localName = this_.priv_Rewriter_maybeGetNameOfImportedValue_2759(t.GetScope(), exportSpec.GetFirstChild())
				j2gort.CheckNotNullString(localName)
			}
			j2gort.MapPutStr(this_.exportedNameToLocalQName, exportSpec.GetLastChild().GetString(), NewEs6RewriteModulesToCommonJsModules_LocalQName(localName, exportSpec))
		}
	}
	export.Detach()
	t.ReportCodeChange_0()
}

// Java method: Rewriter.visitExportNameDeclaration(com.google.javascript.rhino.Node declaration)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_visitExportNameDeclaration_2770(declaration pkg_rhino.Node) {
	_ = declaration
	var iifeArg_1 j2gort.
		Ext_java_util_function_Consumer
	{

		bound_7 := this_
		iifeArg_1 = j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) { bound_7.priv_Rewriter_addExportedName_2771(j2gort.CastNode[pkg_rhino.Node](mrp0_)) })
	}
	NodeUtil_VisitLhsNodesInNode(declaration, iifeArg_1)
}

// Java method: Rewriter.addExportedName(com.google.javascript.rhino.Node lhs)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_addExportedName_2771(lhs pkg_rhino.Node) {
	_ = lhs
	if !(lhs.IsName()) {
		j2gort.CheckState(false)
	}
	var name string = lhs.GetString()
	_ = name
	j2gort.MapPutStr(this_.exportedNameToLocalQName, name, NewEs6RewriteModulesToCommonJsModules_LocalQName(name, lhs))
}

// Java method: Rewriter.visitExportDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_visitExportDeclaration_2772(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	var declaration pkg_rhino.Node = export.GetFirstChild()
	_ = declaration
	if NodeUtil_IsNameDeclaration(declaration) {
		this_.priv_Rewriter_visitExportNameDeclaration_2770(declaration)
	} else {
		if !(declaration.IsFunction() || declaration.IsClass()) {
			j2gort.CheckState(false)
		}
		var name string = declaration.GetFirstChild().GetString()
		_ = name
		j2gort.MapPutStr(this_.exportedNameToLocalQName, name, NewEs6RewriteModulesToCommonJsModules_LocalQName(name, export))
	}
	export.ReplaceWith(declaration.Detach())
	t.ReportCodeChange_0()
}

// Java method: Rewriter.visitExportStar(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_visitExportStar_2773(t *NodeTraversal, export pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = export
	_ = parent
	var moduleIdentifier pkg_rhino.Node = export.GetLastChild()
	_ = moduleIdentifier
	var importNode pkg_rhino.Node = pkg_rhino.IR_ImportNode(pkg_rhino.IR_Empty(), pkg_rhino.IR_Empty(), moduleIdentifier.CloneNode_0())
	_ = importNode
	importNode.Srcref(export)
	importNode.InsertBefore(export)
	this_.Visit_3(t, importNode, parent)
	var moduleName string = this_.priv_Rewriter_getVarNameOfImport_2757(moduleIdentifier.GetString())
	_ = moduleName
	export.ReplaceWith(pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("$$module"), "exportAllFrom"), []pkg_rhino.Node{pkg_rhino.IR_Name(moduleName)})).SrcrefTree(export))
	t.ReportCodeChange_0()
}

// Java method: Rewriter.visitExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteModulesToCommonJsModules_Rewriter) priv_Rewriter_visitExport_2774(t *NodeTraversal, export pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = export
	_ = parent
	if export.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
		this_.priv_Rewriter_visitExportDefault_2767(t, export)
	} else if export.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_ALL_FROM()) {
		this_.priv_Rewriter_visitExportStar_2773(t, export, parent)
	} else if export.HasTwoChildren() {
		this_.priv_Rewriter_visitExportFrom_2768(t, export, parent)
	} else {
		if export.GetFirstChild().IsExportSpecs() {
			this_.priv_Rewriter_visitExportSpecs_2769(t, export)
		} else {
			this_.priv_Rewriter_visitExportDeclaration_2772(t, export)
		}
	}
}
func (*Es6RewriteModulesToCommonJsModules_Rewriter) AlwaysShouldTraverse_()      {}
func (*Es6RewriteModulesToCommonJsModules_Rewriter) CanReleaseSyntacticScopes_() {}

/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PolymerPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.PolymerPass
type PolymerPass struct {
	NodeTraversal_ExternsSkippingCallbackBase
	compiler                     AbstractCompiler                                  `gson:"compiler"`
	tagNameMap                   j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"tagNameMap"`
	nativeExternsAdded           j2gort.OrderedSet[string]                         `gson:"nativeExternsAdded"`
	polymerElementProps          []pkg_rhino.Node                                  `gson:"polymerElementProps"`
	globalNames                  *GlobalNamespace                                  `gson:"globalNames"`
	behaviorExtractor            *PolymerBehaviorExtractor                         `gson:"behaviorExtractor"`
	polymerElementExterns        pkg_rhino.Node                                    `gson:"polymerElementExterns"`
	warnedPolymer1ExternsMissing bool                                              `gson:"warnedPolymer1ExternsMissing"`
	propertySinkExternInjected   bool                                              `gson:"propertySinkExternInjected"`
}

func (this_ *PolymerPass) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerPass) IsJavaNil_() bool { return this_ == nil }
func (*PolymerPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerPass", "PolymerPass"
}
func (this_ *PolymerPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_ExternsSkippingCallbackBase.IdentityHash_()
}
func (*PolymerPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerPass", "com.google.javascript.jscomp.PolymerPass", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerPass](v_)
	})
}
func (*PolymerPass) Implements_jscomp_NodeTraversal_ExternsSkippingCallback() {}
func (*PolymerPass) Implements_jscomp_CompilerPass()                          {}
func (*PolymerPass) Implements_jscomp_NodeTraversal_Callback()                {}

func StaticGet_PolymerPass_VIRTUAL_FILE() string { return "<PolymerPass.java>" }
func (this_ *PolymerPass) initFields_PolymerPass() {
	_ = this_

}

// Java constructor: PolymerPass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewPolymerPass(compiler AbstractCompiler) *PolymerPass {
	this_ := &PolymerPass{}
	InitPolymerPass(this_, compiler)
	return this_
}

func InitPolymerPass(this_ *PolymerPass, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_ExternsSkippingCallback(&this_.NodeTraversal_ExternsSkippingCallbackBase)
	this_.initFields_PolymerPass()
	this_.compiler = compiler
	this_.tagNameMap = TagNameToType_GetMap()
	this_.nativeExternsAdded = j2gort.NewOrderedSet[string]()
}

// Java method: PolymerPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *PolymerPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer func() {
		this_.behaviorExtractor = nil
		if this_.globalNames != nil {
			this_.globalNames.ReleaseSlabs_()
			this_.globalNames = nil
		}
	}()

	var val_1 PolymerPassFindExterns
	InitPolymerPassFindExterns(&val_1)
	var externsCallback *PolymerPassFindExterns = &val_1
	_ = externsCallback
	NodeTraversal_Traverse(this_.compiler, externs, externsCallback)
	this_.polymerElementExterns = externsCallback.GetPolymerElementExterns()
	this_.polymerElementProps = j2gort.ImmListToNodeSliceNullable[pkg_rhino.Node](externsCallback.GetPolymerElementProps())
	this_.globalNames = NewGlobalNamespace_3_AbstractCompiler_Node_Node(this_.compiler, externs, root)
	this_.behaviorExtractor = NewPolymerBehaviorExtractor(this_.compiler, this_.globalNames, this_.compiler.GetModuleMetadataMap(), this_.compiler.GetModuleMap())
	var externsAndJsRoot pkg_rhino.Node = root.GetParent()
	_ = externsAndJsRoot
	NodeTraversal_Traverse(this_.compiler, externsAndJsRoot, this_)
	var suppressBehaviorsCallback *PolymerPassSuppressBehaviorsAndProtectKeys = NewPolymerPassSuppressBehaviorsAndProtectKeys(this_.compiler)
	_ = suppressBehaviorsCallback
	NodeTraversal_Traverse(this_.compiler, externsAndJsRoot, suppressBehaviorsCallback)
}

// Java method: PolymerPass.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *PolymerPass) Visit_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = node
	_ = parent
	if PolymerPassStaticUtils_IsPolymerCall(node) {
		if this_.polymerElementExterns != 0 {
			this_.priv_PolymerPass_rewritePolymer1ClassDefinition_4401(node, parent, traversal)
		} else if !(this_.warnedPolymer1ExternsMissing) {
			this_.compiler.Report_1(JSError_Make_3(node, StaticGet_PolymerPassErrors_POLYMER_MISSING_EXTERNS(), []string{}))
			this_.warnedPolymer1ExternsMissing = true
		}
	} else if PolymerPassStaticUtils_IsPolymerClass(node) {
		this_.priv_PolymerPass_rewritePolymer2ClassDefinition_4402(node, traversal)
	}
}

// Java method: PolymerPass.getModuleMetadata(com.google.javascript.jscomp.NodeTraversal traversal)
func (this_ *PolymerPass) priv_PolymerPass_getModuleMetadata_4400(traversal *NodeTraversal) ModuleMetadataMap_ModuleMetadata {
	_ = traversal
	var script pkg_rhino.Node = traversal.GetCurrentScript()
	_ = script
	if (script != 0) && script.GetFirstChild().IsModuleBody() {
		return ModuleImportResolver_GetModuleFromScopeRoot(this_.compiler.GetModuleMap(), Upcast_CompilerInputProvider(this_.compiler), script.GetFirstChild()).Metadata()
	}
	var scopeRoot pkg_rhino.Node = traversal.GetScopeRoot()
	_ = scopeRoot
	if NodeUtil_IsBundledGoogModuleScopeRoot(scopeRoot) {
		return ModuleImportResolver_GetModuleFromScopeRoot(this_.compiler.GetModuleMap(), Upcast_CompilerInputProvider(this_.compiler), scopeRoot).Metadata()
	}
	return nil
}

// Java method: PolymerPass.rewritePolymer1ClassDefinition(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.NodeTraversal traversal)
func (this_ *PolymerPass) priv_PolymerPass_rewritePolymer1ClassDefinition_4401(node pkg_rhino.Node, parent pkg_rhino.Node, traversal *NodeTraversal) {
	_ = node
	_ = parent
	_ = traversal
	var grandparent pkg_rhino.Node = parent.GetParent()
	_ = grandparent
	if grandparent.IsConst() {
		grandparent.SetToken(pkg_rhino.StaticGet_Token_LET())
		var scriptNode pkg_rhino.Node = traversal.GetCurrentScript()
		_ = scriptNode
		if scriptNode != 0 {
			NodeUtil_AddFeatureToScript(scriptNode, pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
		}
		traversal.ReportCodeChange_0()
	}
	if NodeUtil_IsNameDeclaration(grandparent) && grandparent.GetParent().IsExport() {
		this_.priv_PolymerPass_normalizePolymerExport_4403(grandparent, grandparent.GetParent())
		traversal.ReportCodeChange_0()
	}
	var def *PolymerClassDefinition = PolymerClassDefinition_ExtractFromCallNode(node, this_.compiler, this_.priv_PolymerPass_getModuleMetadata_4400(traversal), this_.behaviorExtractor)
	_ = def
	if def != nil {
		if def.nativeBaseElement != j2gort.NullString {
			this_.priv_PolymerPass_appendPolymerElementExterns_4404(def)
		}
		clinit_PolymerClassRewriter()
		var val_2 PolymerClassRewriter
		InitPolymerClassRewriter(&val_2, this_.compiler)
		var rewriter *PolymerClassRewriter = &val_2
		_ = rewriter
		rewriter.RewritePolymerCall(def, traversal)
	}
}

// Java method: PolymerPass.rewritePolymer2ClassDefinition(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.NodeTraversal traversal)
func (this_ *PolymerPass) priv_PolymerPass_rewritePolymer2ClassDefinition_4402(node pkg_rhino.Node, traversal *NodeTraversal) {
	_ = node
	_ = traversal
	var def *PolymerClassDefinition = PolymerClassDefinition_ExtractFromClassNode(node, this_.compiler)
	_ = def
	if def != nil {
		clinit_PolymerClassRewriter()
		var val_3 PolymerClassRewriter
		InitPolymerClassRewriter(&val_3, this_.compiler)
		var rewriter *PolymerClassRewriter = &val_3
		_ = rewriter
		{
			r_4 := rewriter
			r_4.propertySinkExternInjected = this_.propertySinkExternInjected
		}
		rewriter.RewritePolymerClassDeclaration(node, traversal, def)
		this_.propertySinkExternInjected = rewriter.propertySinkExternInjected
		PolymerPass_ExportTsickleProps(def.definition, def.target.GetQualifiedName())
	}
}

// Java method: PolymerPass.normalizePolymerExport(com.google.javascript.rhino.Node nameDecl, com.google.javascript.rhino.Node export)
func (this_ *PolymerPass) priv_PolymerPass_normalizePolymerExport_4403(nameDecl pkg_rhino.Node, export pkg_rhino.Node) {
	_ = nameDecl
	_ = export
	var name pkg_rhino.Node = nameDecl.GetFirstChild()
	_ = name
	nameDecl.Detach().InsertBefore(export)
	var exportSpec pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_EXPORT_SPEC())
	_ = exportSpec
	exportSpec.AddChildToFront(name.CloneNode_0())
	exportSpec.AddChildToFront(name.CloneNode_0())
	export.AddChildToFront(pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_EXPORT_SPECS(), exportSpec).SrcrefTree(export))
}

// Java method: PolymerPass.appendPolymerElementExterns(com.google.javascript.jscomp.PolymerClassDefinition def)
func (this_ *PolymerPass) priv_PolymerPass_appendPolymerElementExterns_4404(def *PolymerClassDefinition) {
	_ = def
	if !(this_.nativeExternsAdded.Add(def.nativeBaseElement)) {
		return
	}
	var block pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = block
	var baseExterns pkg_rhino.Node = this_.polymerElementExterns.CloneTree_0()
	_ = baseExterns
	var polymerElementType string = PolymerPassStaticUtils_GetPolymerElementType(def)
	_ = polymerElementType
	baseExterns.GetFirstChild().SetString(polymerElementType)
	var elementType string = j2gort.CastString(j2gort.MapGetStr(this_.tagNameMap, def.nativeBaseElement))
	_ = elementType
	if elementType == j2gort.NullString {
		this_.compiler.Report_1(JSError_Make_3(def.descriptor, StaticGet_PolymerPassErrors_POLYMER_INVALID_EXTENDS(), []string{def.nativeBaseElement}))
		return
	}
	var elementBaseType pkg_rhino.JSTypeExpression = pkg_rhino.NewJSTypeExpression(pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BANG(), pkg_rhino.IR_String(elementType).SrcrefTree(this_.polymerElementExterns)), "<PolymerPass.java>")
	_ = elementBaseType
	var baseDocs *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_CopyFrom(baseExterns.GetJSDocInfo())
	_ = baseDocs
	baseDocs.ChangeBaseType(elementBaseType)
	baseExterns.SetJSDocInfo(baseDocs.Build_0())
	block.AddChildToBack(baseExterns)
	{

		for _, baseProp := range this_.polymerElementProps {

			_ = baseProp
			var newProp pkg_rhino.Node = baseProp.CloneTree_0()
			_ = newProp
			var newPropRootName pkg_rhino.Node = NodeUtil_GetRootOfQualifiedName_1_Node(newProp.GetFirstFirstChild())
			_ = newPropRootName
			newPropRootName.SetString(polymerElementType)
			block.AddChildToBack(newProp)
		}
	}
	block.SrcrefTreeIfMissing(this_.polymerElementExterns)
	var parent pkg_rhino.Node = this_.polymerElementExterns.GetParent()
	_ = parent
	var stmts pkg_rhino.Node = block.RemoveChildren()
	_ = stmts
	parent.AddChildrenAfter(stmts, this_.polymerElementExterns)
	this_.compiler.ReportChangeToEnclosingScope(stmts)
}

// Java static method: PolymerPass.exportTsickleProps(com.google.javascript.rhino.Node classNode, java.lang.String classQualifiedName)
func PolymerPass_ExportTsickleProps(classNode pkg_rhino.Node, classQualifiedName string) {
	_ = classNode
	_ = classQualifiedName
	var classStmt pkg_rhino.Node = j2gort.CheckNotNullNode(NodeUtil_GetEnclosingStatement(classNode))
	_ = classStmt
	{
		var sibling pkg_rhino.Node = classStmt.GetNext()
		_ = sibling
		for ; sibling != 0; sibling = sibling.GetNext() {
			var ifBlock pkg_rhino.Node = PolymerPass_MaybeGetTsicklePropTypesIfBlock(sibling)
			_ = ifBlock
			if ifBlock == 0 {
				continue
			}
			{
				var stmt pkg_rhino.Node = ifBlock.GetFirstChild()
				_ = stmt
				for ; stmt != 0; stmt = stmt.GetNext() {
					var prop pkg_rhino.Node = PolymerPass_GetTsickleClassPrototypeProp(stmt, classQualifiedName)
					_ = prop
					if prop == 0 {
						continue
					}
					var jsdoc *pkg_rhino.JSDocInfo = j2gort.CheckNotNullPtr(prop.GetJSDocInfo())
					_ = jsdoc
					if jsdoc.GetVisibility() == pkg_rhino.StaticGet_JSDocInfo_Visibility_PRIVATE() {
						continue
					}
					var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_CopyFrom(jsdoc)
					_ = builder
					builder.RecordExport()
					builder.OverwriteVisibility(pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED())
					prop.SetJSDocInfo(builder.Build_0())
				}
			}
		}
	}
}

// Java static method: PolymerPass.maybeGetTsicklePropTypesIfBlock(com.google.javascript.rhino.Node node)
func PolymerPass_MaybeGetTsicklePropTypesIfBlock(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if !(node.IsIf()) {
		return pkg_rhino.Node(0)
	}
	var condition pkg_rhino.Node = node.GetFirstChild()
	_ = condition
	if !(condition.IsFalse()) {
		return pkg_rhino.Node(0)
	}
	var block pkg_rhino.Node = condition.GetNext()
	_ = block
	if !((block != 0) && block.IsBlock()) {
		j2gort.CheckStateFmt(false, "Expected block but got %s", j2gort.BoxNode(block))
	}
	return block
}

// Java static method: PolymerPass.getTsickleClassPrototypeProp(com.google.javascript.rhino.Node stmt, java.lang.String classQualifiedName)
func PolymerPass_GetTsickleClassPrototypeProp(stmt pkg_rhino.Node, classQualifiedName string) pkg_rhino.Node {
	_ = stmt
	_ = classQualifiedName
	if !(stmt.IsExprResult()) {
		return pkg_rhino.Node(0)
	}
	var expr pkg_rhino.Node = stmt.GetFirstChild()
	_ = expr
	if !(expr.IsGetProp()) {
		return pkg_rhino.Node(0)
	}
	var owner pkg_rhino.Node = expr.GetFirstChild()
	_ = owner
	if !(owner.IsGetProp()) || !(j2gort.CheckNotNullString(owner.GetString()) == "prototype") {
		return pkg_rhino.Node(0)
	}
	var classRef pkg_rhino.Node = owner.GetFirstChild()
	_ = classRef
	if !(classRef.MatchesQualifiedName_1_string(classQualifiedName)) {
		return pkg_rhino.Node(0)
	}
	return expr
}

// Java class: com.google.javascript.jscomp.PolymerPass.MemberDefinition
type PolymerPass_MemberDefinition struct {
	info            *pkg_rhino.JSDocInfo `gson:"info"`
	name            pkg_rhino.Node       `gson:"name"`
	value           pkg_rhino.Node       `gson:"value"`
	enclosingModule pkg_rhino.Node       `gson:"enclosingModule"`
	idHash_         atomic.Int32
}

func (this_ *PolymerPass_MemberDefinition) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerPass_MemberDefinition) IsJavaNil_() bool { return this_ == nil }
func (*PolymerPass_MemberDefinition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerPass$MemberDefinition", "MemberDefinition"
}
func (this_ *PolymerPass_MemberDefinition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolymerPass_MemberDefinition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerPass.MemberDefinition", "com.google.javascript.jscomp.PolymerPass$MemberDefinition", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerPass_MemberDefinition](v_)
	})
}

func (this_ *PolymerPass_MemberDefinition) initFields_PolymerPass_MemberDefinition() {
	_ = this_
}

// Java constructor: MemberDefinition(com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.Node name, com.google.javascript.rhino.Node value, com.google.javascript.rhino.Node enclosingModule)
func NewPolymerPass_MemberDefinition(info *pkg_rhino.JSDocInfo, name pkg_rhino.Node, value pkg_rhino.Node, enclosingModule pkg_rhino.Node) *PolymerPass_MemberDefinition {
	this_ := &PolymerPass_MemberDefinition{}
	InitPolymerPass_MemberDefinition(this_, info, name, value, enclosingModule)
	return this_
}

func InitPolymerPass_MemberDefinition(this_ *PolymerPass_MemberDefinition, info *pkg_rhino.JSDocInfo, name pkg_rhino.Node, value pkg_rhino.Node, enclosingModule pkg_rhino.Node) {
	_ = info
	_ = name
	_ = value
	_ = enclosingModule
	this_.initFields_PolymerPass_MemberDefinition()
	this_.info = info
	this_.name = name
	this_.value = value
	this_.enclosingModule = enclosingModule
}

// Java method: MemberDefinition.toString()
func (this_ *PolymerPass_MemberDefinition) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_any("name", j2gort.BoxNode(this_.name)).Add_2_string_any("value", j2gort.BoxNode(this_.value)).String()
}

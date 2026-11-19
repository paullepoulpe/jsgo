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
// Port of com/google/javascript/jscomp/Es6ConvertSuper.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_es6_convert_super_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_es6_convert_super_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_1 uint32) bool {
		n := pkg_rhino.Node(rawNode_1)
		_ = n
		{

			switch n.TokenSwitchOrd_() {
			case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(129) /* Token.COMPUTED_PROP */ :
				return true
			default:
				return false
			}
			panic("unreachable switch expression")
		}

	})
}

// Java class: com.google.javascript.jscomp.Es6ConvertSuper
type Es6ConvertSuper struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler               AbstractCompiler                `gson:"compiler"`
	astFactory             *AstFactory                     `gson:"astFactory"`
	constructorSynthesizer *SynthesizeExplicitConstructors `gson:"constructorSynthesizer"`
}

func (this_ *Es6ConvertSuper) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6ConvertSuper) IsJavaNil_() bool { return this_ == nil }
func (*Es6ConvertSuper) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6ConvertSuper", "Es6ConvertSuper"
}
func (this_ *Es6ConvertSuper) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6ConvertSuper) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6ConvertSuper", "com.google.javascript.jscomp.Es6ConvertSuper", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6ConvertSuper](v_)
	})
}
func (*Es6ConvertSuper) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*Es6ConvertSuper) Implements_jscomp_CompilerPass()                            {}
func (*Es6ConvertSuper) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_Es6ConvertSuper_FEATURES_TO_RUN_FOR *pkg_parser.FeatureSet

func StaticGet_Es6ConvertSuper_FEATURES_TO_RUN_FOR() *pkg_parser.FeatureSet {
	clinit_Es6ConvertSuper()
	return static_Es6ConvertSuper_FEATURES_TO_RUN_FOR
}
func StaticSet_Es6ConvertSuper_FEATURES_TO_RUN_FOR(v *pkg_parser.FeatureSet) {
	clinit_Es6ConvertSuper()
	static_Es6ConvertSuper_FEATURES_TO_RUN_FOR = v
}

var clinitDone_Es6ConvertSuper uint32
var clinitState_Es6ConvertSuper j2gort.ClinitState

func clinit_Es6ConvertSuper() {
	if atomic.LoadUint32(&clinitDone_Es6ConvertSuper) != 0 {
		return
	}
	clinitSlow_Es6ConvertSuper()
}

func clinitSlow_Es6ConvertSuper() {
	if !j2gort.EnterClinit(&clinitState_Es6ConvertSuper, &clinitDone_Es6ConvertSuper) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6ConvertSuper, &clinitDone_Es6ConvertSuper)
	static_Es6ConvertSuper_FEATURES_TO_RUN_FOR = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7232([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_SUPER(), pkg_parser.StaticGet_FeatureSet_Feature_CLASSES()})
}

func (this_ *Es6ConvertSuper) initFields_Es6ConvertSuper() {
	_ = this_
}

// Java constructor: Es6ConvertSuper(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6ConvertSuper(compiler AbstractCompiler) *Es6ConvertSuper {
	clinit_Es6ConvertSuper()
	this_ := &Es6ConvertSuper{}
	InitEs6ConvertSuper(this_, compiler)
	return this_
}

func InitEs6ConvertSuper(this_ *Es6ConvertSuper, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6ConvertSuper()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.constructorSynthesizer = NewSynthesizeExplicitConstructors(compiler)
}

// Java method: Es6ConvertSuper.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6ConvertSuper) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsClass() {
		this_.constructorSynthesizer.SynthesizeClassConstructorIfMissing(t, n)
	} else if n.IsSuper() {
		this_.priv_Es6ConvertSuper_visitSuper_2467(n, parent)
	}
}

// Java method: Es6ConvertSuper.isInterface(com.google.javascript.rhino.Node classNode)
func (this_ *Es6ConvertSuper) priv_Es6ConvertSuper_isInterface_2466(classNode pkg_rhino.Node) bool {
	_ = classNode
	var classJsDocInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(classNode)
	_ = classJsDocInfo
	return ((classJsDocInfo != nil) && classJsDocInfo.IsInterface())
}

// Java method: Es6ConvertSuper.visitSuper(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *Es6ConvertSuper) priv_Es6ConvertSuper_visitSuper_2467(node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = node
	_ = parent
	if !(node.IsSuper()) {
		j2gort.CheckState(false)
	}
	var exprRoot pkg_rhino.Node = node
	_ = exprRoot
	if exprRoot.GetParent().IsGetElem() || exprRoot.GetParent().IsGetProp() {
		exprRoot = exprRoot.GetParent()
	}
	var enclosingMemberDef pkg_rhino.Node = NodeUtil_GetEnclosingNode(exprRoot, lambdaSingleton_es6_convert_super_1)
	_ = enclosingMemberDef
	if (enclosingMemberDef == 0) || enclosingMemberDef.GetParent().IsObjectLit() {
		this_.compiler.Report_1(JSError_Make_3(node, StaticGet_TranspilationUtil_CANNOT_CONVERT(), []string{"super is not supported in object literal methods."}))
		return
	}
	if parent.IsCall() {
		this_.priv_Es6ConvertSuper_visitSuperCall_2468(node, parent, enclosingMemberDef)
	} else if parent.IsGetProp() || parent.IsGetElem() {
		if parent.GetFirstChild() == node {
			if parent.GetParent().IsCall() && NodeUtil_IsInvocationTarget(parent) {
				this_.priv_Es6ConvertSuper_visitSuperPropertyCall_2469(node, parent, enclosingMemberDef)
			} else {
				this_.priv_Es6ConvertSuper_visitSuperPropertyAccess_2470(node, parent, enclosingMemberDef)
			}
		} else {
			this_.compiler.Report_1(JSError_Make_3(node, StaticGet_TranspilationUtil_CANNOT_CONVERT_YET(), []string{"Only calls to super or to a method of super are supported."}))
		}
	} else if parent.IsNew() {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "This should never happen. Did Es6SuperCheck fail to run?")))
	} else {
		this_.compiler.Report_1(JSError_Make_3(node, StaticGet_TranspilationUtil_CANNOT_CONVERT_YET(), []string{"Only calls to super or to a method of super are supported."}))
	}
}

// Java method: Es6ConvertSuper.visitSuperCall(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node enclosingMemberDef)
func (this_ *Es6ConvertSuper) priv_Es6ConvertSuper_visitSuperCall_2468(node pkg_rhino.Node, parent pkg_rhino.Node, enclosingMemberDef pkg_rhino.Node) {
	_ = node
	_ = parent
	_ = enclosingMemberDef
	if !(parent.IsCall()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
	}
	if !(node.IsSuper()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
	}
	var clazz pkg_rhino.Node = NodeUtil_GetEnclosingClass(node)
	_ = clazz
	var superName pkg_rhino.Node = clazz.GetSecondChild()
	_ = superName
	if !(superName.IsQualifiedName()) {
		return
	}
	if NodeUtil_IsEs6ConstructorMemberFunctionDef(enclosingMemberDef) {
		if node.IsFromExterns() || this_.priv_Es6ConvertSuper_isInterface_2466(clazz) {
			var enclosingStatement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(node)
			_ = enclosingStatement
			var enclosingStatementParent pkg_rhino.Node = enclosingStatement.GetParent()
			_ = enclosingStatementParent
			enclosingStatement.Detach()
			this_.compiler.ReportChangeToEnclosingScope(enclosingStatementParent)
		}
		return
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "This should never happen. Did Es6SuperCheck fail to run?")))
	}
}

// Java method: Es6ConvertSuper.visitSuperPropertyCall(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node enclosingMemberDef)
func (this_ *Es6ConvertSuper) priv_Es6ConvertSuper_visitSuperPropertyCall_2469(node pkg_rhino.Node, parent pkg_rhino.Node, enclosingMemberDef pkg_rhino.Node) {
	_ = node
	_ = parent
	_ = enclosingMemberDef
	if !(parent.IsGetProp() || parent.IsGetElem()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
	}
	if !(node.IsSuper()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
	}
	var grandparent pkg_rhino.Node = parent.GetParent()
	_ = grandparent
	if !(grandparent.IsCall()) {
		j2gort.CheckState(false)
	}
	var clazz pkg_rhino.Node = NodeUtil_GetEnclosingClass(node)
	_ = clazz
	var superName pkg_rhino.Node = clazz.GetSecondChild()
	_ = superName
	if !(superName.IsQualifiedName()) {
		return
	}
	var callTarget pkg_rhino.Node = parent
	_ = callTarget
	if enclosingMemberDef.IsStaticMember() {
		var expandedSuper pkg_rhino.Node = superName.CloneTree_0().SrcrefTree(node)
		_ = expandedSuper
		expandedSuper.SetOriginalName("super")
		node.ReplaceWith(expandedSuper)
		callTarget = this_.astFactory.CreateGetPropWithUnknownType(callTarget.Detach(), "call")
		grandparent.AddChildToFront(callTarget)
		var thisNode pkg_rhino.Node = this_.astFactory.CreateThis(AstFactory_S_Type_1_Node(clazz))
		_ = thisNode
		thisNode.MakeNonIndexable()
		thisNode.InsertAfter(callTarget)
		grandparent.SrcrefTreeIfMissing(parent)
		this_.priv_Es6ConvertSuper_ensureNoCollapseOnStaticMemberDef_2471(enclosingMemberDef)
	} else {
		var expandedSuper_2 pkg_rhino.Node = this_.astFactory.CreatePrototypeAccess(superName.CloneTree_0()).SrcrefTree(node)
		_ = expandedSuper_2
		expandedSuper_2.SetOriginalName("super")
		node.ReplaceWith(expandedSuper_2)
		callTarget = this_.astFactory.CreateGetPropWithUnknownType(callTarget.Detach(), "call")
		grandparent.AddChildToFront(callTarget)
		var thisNode_2 pkg_rhino.Node = this_.astFactory.CreateThisForEs6Class(clazz)
		_ = thisNode_2
		thisNode_2.MakeNonIndexable()
		thisNode_2.InsertAfter(callTarget)
		grandparent.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), false)
		grandparent.SrcrefTreeIfMissing(parent)
	}
	this_.compiler.ReportChangeToEnclosingScope(grandparent)
}

// Java method: Es6ConvertSuper.visitSuperPropertyAccess(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node enclosingMemberDef)
func (this_ *Es6ConvertSuper) priv_Es6ConvertSuper_visitSuperPropertyAccess_2470(node pkg_rhino.Node, parent pkg_rhino.Node, enclosingMemberDef pkg_rhino.Node) {
	_ = node
	_ = parent
	_ = enclosingMemberDef
	if !(parent.IsGetProp() || parent.IsGetElem()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
	}
	if !(node.IsSuper()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
	}
	var grandparent pkg_rhino.Node = parent.GetParent()
	_ = grandparent
	if NodeUtil_IsLValue(parent) {
		this_.compiler.Report_1(JSError_Make_3(parent, StaticGet_TranspilationUtil_CANNOT_CONVERT_YET(), []string{"assigning to a super property"}))
		return
	}
	var clazz pkg_rhino.Node = NodeUtil_GetEnclosingClass(node)
	_ = clazz
	var superName pkg_rhino.Node = clazz.GetSecondChild()
	_ = superName
	if !(superName.IsQualifiedName()) {
		return
	}
	var newProp pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newProp
	if enclosingMemberDef.IsStaticMember() {
		newProp = superName.CloneTree_0().SrcrefTree(node)
	} else {
		newProp = this_.astFactory.CreatePrototypeAccess(superName.CloneTree_0()).SrcrefTree(node)
	}
	if parent.IsGetProp() && this_.compiler.GetAccessorSummary().GetKind(parent.GetString()).HasGetter() {
		var superPlaceholder pkg_rhino.Node = pkg_rhino.IR_Name("sp")
		_ = superPlaceholder
		var reflectGetCall pkg_rhino.Node = this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateConstantName("Reflect", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), "get"), AstFactory_S_Type_1_Node(parent), []pkg_rhino.Node{superPlaceholder, func() pkg_rhino.Node {
			if parent.IsGetProp() {
				return this_.astFactory.CreateCall(this_.astFactory.CreateName_2_string_AstFactory_Type("JSCompiler_renameProperty", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_STRING()), []pkg_rhino.Node{this_.astFactory.CreateString(parent.GetString()), superName.CloneTree_0().SrcrefTree(node)})
			}
			return node.GetNext().Detach()
		}(), this_.astFactory.CreateThis(AstFactory_S_Type_1_Node(clazz))}).SrcrefTreeIfMissing(parent)
		_ = reflectGetCall
		parent.ReplaceWith(reflectGetCall)
		superPlaceholder.ReplaceWith(node.Detach())
		this_.priv_Es6ConvertSuper_ensureNoCollapseOnStaticMemberDef_2471(enclosingMemberDef)
	}
	node.ReplaceWith(newProp)
	newProp.SetOriginalName("super")
	this_.compiler.ReportChangeToEnclosingScope(grandparent)
}

// Java method: Es6ConvertSuper.ensureNoCollapseOnStaticMemberDef(com.google.javascript.rhino.Node memberDef)
func (this_ *Es6ConvertSuper) priv_Es6ConvertSuper_ensureNoCollapseOnStaticMemberDef_2471(memberDef pkg_rhino.Node) {
	_ = memberDef
	if !(memberDef.IsStaticMember()) {
		return
	}
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(memberDef.GetJSDocInfo())
	_ = builder
	if builder.RecordNoCollapse() {
		memberDef.SetJSDocInfo(builder.Build_0())
	}
}

// Java method: Es6ConvertSuper.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6ConvertSuper) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	TranspilationPasses_ProcessTranspile(this_.compiler, externs, StaticGet_Es6ConvertSuper_FEATURES_TO_RUN_FOR(), []NodeTraversal_Callback{this_})
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_Es6ConvertSuper_FEATURES_TO_RUN_FOR(), []NodeTraversal_Callback{this_})
}
func (*Es6ConvertSuper) AlwaysShouldTraverse_()      {}
func (*Es6ConvertSuper) CanReleaseSyntacticScopes_() {}

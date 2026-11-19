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
// Port of com/google/javascript/jscomp/Denormalize.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_denormalize_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_denormalize_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsIn() })
}

// Java class: com.google.javascript.jscomp.Denormalize
type Denormalize struct {
	compiler         AbstractCompiler       `gson:"compiler"`
	outputFeatureSet *pkg_parser.FeatureSet `gson:"outputFeatureSet"`
	idHash_          atomic.Int32
}

func (this_ *Denormalize) BindSelf_(s any) {
	_ = s
}
func (this_ *Denormalize) IsJavaNil_() bool { return this_ == nil }
func (*Denormalize) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Denormalize", "Denormalize"
}
func (this_ *Denormalize) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Denormalize) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Denormalize", "com.google.javascript.jscomp.Denormalize", func(v_ any) bool {
		return j2gort.InstanceOf[*Denormalize](v_)
	})
}
func (*Denormalize) Implements_jscomp_CompilerPass()                {}
func (*Denormalize) Implements_jscomp_NodeTraversal_Callback()      {}
func (*Denormalize) Implements_jscomp_ReferenceCollector_Behavior() {}

func (this_ *Denormalize) initFields_Denormalize() {
	_ = this_
}

// Java constructor: Denormalize(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.parsing.parser.FeatureSet outputFeatureSet)
func NewDenormalize(compiler AbstractCompiler, outputFeatureSet *pkg_parser.FeatureSet) *Denormalize {
	this_ := &Denormalize{}
	InitDenormalize(this_, compiler, outputFeatureSet)
	return this_
}

func InitDenormalize(this_ *Denormalize, compiler AbstractCompiler, outputFeatureSet *pkg_parser.FeatureSet) {
	_ = compiler
	_ = outputFeatureSet
	this_.initFields_Denormalize()
	this_.compiler = compiler
	this_.outputFeatureSet = outputFeatureSet
}

// Java method: Denormalize.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Denormalize) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	if this_.compiler.GetOptions().GetSyntheticBlockStartMarker() == j2gort.NullString {
		{
			clinit_ReferenceCollector()
			var val_1 ReferenceCollector
			InitReferenceCollector_3(&val_1, this_.compiler, this_, NewSyntacticScopeCreator_1(this_.compiler))
			(&val_1).ownsScopeCreator_ = true
			defer (&val_1).ReleaseSlabs_()

			(&val_1).Process_1(root)
		}
	}
}

// Java method: Denormalize.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Denormalize) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return true
}

// Java method: Denormalize.afterExitScope(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.ReferenceMap referenceMap)
func (this_ *Denormalize) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	_ = t
	_ = referenceMap
	var scopeRoot pkg_rhino.Node = t.GetScopeRoot()
	_ = scopeRoot
	if scopeRoot.IsBlock() && scopeRoot.GetParent().IsFunction() {
		var changed bool = false
		_ = changed
		{
			cur_2 := (t.GetScope()).iterVars_()
			for cur_2.HasNext() {

				v := cur_2.Next()

				_ = v
				var references *ReferenceCollection = j2gort.CastPtr[ReferenceCollection](referenceMap.GetReferences(j2gort.PtrToAny(v)))
				_ = references
				var declaration *Reference = nil
				_ = declaration
				var assign *Reference = nil
				_ = assign
				{
					cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(references))
					for cur_4.HasNext() {
						rawElem_5 := cur_4.Next()
						r := j2gort.CastPtr[Reference](rawElem_5)
						_ = r
						if (r.IsVarDeclaration() && NodeUtil_IsStatement(r.GetNode().GetParent())) && !(r.IsInitializingDeclaration()) {
							declaration = r
						} else if ((assign == nil) && r.IsSimpleAssignmentToName()) && func() bool {
							rp_6 := j2gort.CastPtr[Scope](r.GetScope().GetClosestHoistScope())
							ap_7 := t.GetScope()
							return j2gort.CheckNotNullPtr(rp_6) == ap_7
						}() {
							assign = r
						}
					}
				}
				if (declaration != nil) && (assign != nil) {
					var lhs pkg_rhino.Node = assign.GetNode()
					_ = lhs
					var assignNode pkg_rhino.Node = lhs.GetParent()
					_ = assignNode
					if assignNode.GetParent().IsExprResult() {
						var rhs pkg_rhino.Node = lhs.GetNext()
						_ = rhs
						assignNode.GetParent().ReplaceWith(pkg_rhino.IR_Var_2(lhs.Detach(), rhs.Detach()))
						var var_ pkg_rhino.Node = declaration.GetNode().GetParent()
						_ = var_
						if !(var_.IsVar()) {
							j2gort.CheckStateMsg(false, j2gort.BoxNode(var_))
						}
						NodeUtil_RemoveChild(var_, declaration.GetNode())
						changed = true
					}
				}
			}
		}
		if changed {
			t.ReportCodeChange_0()
		}
	}
}

// Java method: Denormalize.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Denormalize) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.priv_Denormalize_maybeCollapseIntoForStatements_2386(n, parent)
	this_.priv_Denormalize_maybeCollapseLogicalAssignShorthand_2388(t, n, parent)
	this_.priv_Denormalize_maybeCollapseAssignShorthand_2387(n, parent)
}

// Java method: Denormalize.maybeCollapseIntoForStatements(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Denormalize) priv_Denormalize_maybeCollapseIntoForStatements_2386(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	if (parent == 0) || !(NodeUtil_IsStatementBlock(parent)) {
		return
	}
	if !(n.IsExprResult()) && !(n.IsVar()) {
		return
	}
	var nextSibling pkg_rhino.Node = n.GetNext()
	_ = nextSibling
	if nextSibling == 0 {
		return
	} else if nextSibling.IsForIn() || nextSibling.IsForOf() {
		var forNode pkg_rhino.Node = nextSibling
		_ = forNode
		var forVar pkg_rhino.Node = forNode.GetFirstChild()
		_ = forVar
		if (forVar.IsName() && n.IsVar()) && n.HasOneChild() {
			var name pkg_rhino.Node = n.GetFirstChild()
			_ = name
			if !(name.HasChildren()) && func() bool {
				sr_8 := forVar.GetString()
				se_9 := name.GetString()
				return j2gort.CheckNotNullString(sr_8) == se_9
			}() {
				n.Detach()
				forVar.ReplaceWith(n)
				this_.compiler.ReportChangeToEnclosingScope(parent)
			}
		}
	} else if nextSibling.IsVanillaFor() && nextSibling.GetFirstChild().IsEmpty() {
		if NodeUtil_Has(n, lambdaSingleton_denormalize_1, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue()) {
			return
		}
		var forNode_2 pkg_rhino.Node = nextSibling
		_ = forNode_2
		var oldInitializer pkg_rhino.Node = forNode_2.GetFirstChild()
		_ = oldInitializer
		n.Detach()
		var newInitializer pkg_rhino.Node = pkg_rhino.Node(0)
		_ = newInitializer
		if n.IsVar() {
			newInitializer = n
		} else {
			if !(n.HasOneChild()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
			}
			newInitializer = n.GetFirstChild()
			newInitializer.Detach()
		}
		oldInitializer.ReplaceWith(newInitializer)
		this_.compiler.ReportChangeToEnclosingScope(forNode_2)
	}
}

// Java method: Denormalize.maybeCollapseAssignShorthand(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Denormalize) priv_Denormalize_maybeCollapseAssignShorthand_2387(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	if !(this_.priv_Denormalize_isCollapsableAssign_2389(n)) {
		return
	}
	var op pkg_rhino.Node = n.GetLastChild()
	_ = op
	var assignOp pkg_rhino.Token = this_.priv_Denormalize_getAssignOpFromOp_2391(op)
	_ = assignOp
	{
		var iifeCond_1 bool
		{

			sr_10 := n.GetFirstChild().GetString()
			se_11 := op.GetFirstChild().GetString()
			iifeCond_1 = j2gort.CheckNotNullString(sr_10) == se_11
		}
		if iifeCond_1 {
			op.SetToken(assignOp)
			var opDetached pkg_rhino.Node = op.Detach()
			_ = opDetached
			opDetached.SetJSDocInfo(n.GetJSDocInfo())
			n.ReplaceWith(opDetached)
			this_.compiler.ReportChangeToEnclosingScope(parent)
		}
	}
}

// Java method: Denormalize.maybeCollapseLogicalAssignShorthand(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Denormalize) priv_Denormalize_maybeCollapseLogicalAssignShorthand_2388(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(this_.priv_Denormalize_isCollapsableLogicalAssign_2390(n)) {
		return
	}
	var op pkg_rhino.Node = n.GetLastChild()
	_ = op
	var assignOp pkg_rhino.Token = this_.priv_Denormalize_getAssignOpFromOp_2391(n)
	_ = assignOp
	{
		var iifeCond_1 bool
		{

			sr_12 := n.GetFirstChild().GetString()
			se_13 := op.GetFirstChild().GetString()
			iifeCond_1 = j2gort.CheckNotNullString(sr_12) == se_13
		}
		if iifeCond_1 {
			op.SetToken(assignOp)
			var opDetached pkg_rhino.Node = op.Detach()
			_ = opDetached
			opDetached.SetJSDocInfo(n.GetJSDocInfo())
			n.ReplaceWith(opDetached)
			NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT(), this_.compiler)
			this_.compiler.ReportChangeToEnclosingScope(parent)
		}
	}
}

// Java method: Denormalize.isCollapsableAssign(com.google.javascript.rhino.Node n)
func (this_ *Denormalize) priv_Denormalize_isCollapsableAssign_2389(n pkg_rhino.Node) bool {
	_ = n
	return (((n.IsAssign() && n.GetFirstChild().IsName()) && this_.priv_Denormalize_hasCorrespondingAssignmentOp_2392(n.GetLastChild())) && n.GetLastChild().GetFirstChild().IsName())
}

// Java method: Denormalize.isCollapsableLogicalAssign(com.google.javascript.rhino.Node n)
func (this_ *Denormalize) priv_Denormalize_isCollapsableLogicalAssign_2390(n pkg_rhino.Node) bool {
	_ = n
	return ((((this_.outputFeatureSet.Has(pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT()) && ((n.IsOr() || n.IsAnd()) || n.IsNullishCoalesce())) && n.GetFirstChild().IsName()) && n.GetLastChild().IsAssign()) && n.GetLastChild().GetFirstChild().IsName())
}

// Java method: Denormalize.getAssignOpFromOp(com.google.javascript.rhino.Node n)
func (this_ *Denormalize) priv_Denormalize_getAssignOpFromOp_2391(n pkg_rhino.Node) pkg_rhino.Token {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(1) /* Token.BITOR */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_BITOR()
	case int32(2) /* Token.BITXOR */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_BITXOR()
	case int32(3) /* Token.BITAND */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_BITAND()
	case int32(10) /* Token.LSH */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_LSH()
	case int32(11) /* Token.RSH */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_RSH()
	case int32(12) /* Token.URSH */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_URSH()
	case int32(13) /* Token.ADD */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_ADD()
	case int32(14) /* Token.SUB */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_SUB()
	case int32(15) /* Token.MUL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_MUL()
	case int32(18) /* Token.EXPONENT */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_EXPONENT()
	case int32(16) /* Token.DIV */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_DIV()
	case int32(17) /* Token.MOD */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_MOD()
	case int32(68) /* Token.OR */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_OR()
	case int32(69) /* Token.AND */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_AND()
	case int32(70) /* Token.COALESCE */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_COALESCE()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected operator: "+j2gort.StrVal(n))))
	}
	panic("unreachable switch expression")
}

// Java method: Denormalize.hasCorrespondingAssignmentOp(com.google.javascript.rhino.Node n)
func (this_ *Denormalize) priv_Denormalize_hasCorrespondingAssignmentOp_2392(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(3) /* Token.BITAND */, int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(12) /* Token.URSH */, int32(13) /* Token.ADD */, int32(14) /* Token.SUB */, int32(15) /* Token.MUL */, int32(16) /* Token.DIV */, int32(17) /* Token.MOD */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

func (this_ *Denormalize) BeforeTraverseNode(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node) {
	IfaceDefault_ReferenceCollector_Behavior_BeforeTraverseNode(this_, p0_, p1_, p2_)
}

func (this_ *Denormalize) AfterVisitNode(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node, p3_ *Reference) {
	IfaceDefault_ReferenceCollector_Behavior_AfterVisitNode(this_, p0_, p1_, p2_, p3_)
}
func (*Denormalize) AlwaysShouldTraverse_()      {}
func (*Denormalize) CanReleaseSyntacticScopes_() {}

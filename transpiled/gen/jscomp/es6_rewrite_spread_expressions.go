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
// Port of com/google/javascript/jscomp/Es6RewriteSpreadExpressions.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.Es6RewriteSpreadExpressions
type Es6RewriteSpreadExpressions struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler   AbstractCompiler      `gson:"compiler"`
	astFactory *AstFactory           `gson:"astFactory"`
	namespace  pkg_rhino.StaticScope `gson:"namespace"`
}

func (this_ *Es6RewriteSpreadExpressions) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteSpreadExpressions) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteSpreadExpressions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteSpreadExpressions", "Es6RewriteSpreadExpressions"
}
func (this_ *Es6RewriteSpreadExpressions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteSpreadExpressions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteSpreadExpressions", "com.google.javascript.jscomp.Es6RewriteSpreadExpressions", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteSpreadExpressions](v_)
	})
}
func (*Es6RewriteSpreadExpressions) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*Es6RewriteSpreadExpressions) Implements_jscomp_CompilerPass()                            {}
func (*Es6RewriteSpreadExpressions) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_Es6RewriteSpreadExpressions_TRANSPILED_FEATURES *pkg_parser.FeatureSet

func StaticGet_Es6RewriteSpreadExpressions_TRANSPILED_FEATURES() *pkg_parser.FeatureSet {
	clinit_Es6RewriteSpreadExpressions()
	return static_Es6RewriteSpreadExpressions_TRANSPILED_FEATURES
}
func StaticSet_Es6RewriteSpreadExpressions_TRANSPILED_FEATURES(v *pkg_parser.FeatureSet) {
	clinit_Es6RewriteSpreadExpressions()
	static_Es6RewriteSpreadExpressions_TRANSPILED_FEATURES = v
}
func StaticGet_Es6RewriteSpreadExpressions_FRESH_SPREAD_VAR() string { return "$jscomp$spread$args" }

var static_Es6RewriteSpreadExpressions_arrayType AstFactory_Type

func StaticGet_Es6RewriteSpreadExpressions_ArrayType() AstFactory_Type {
	clinit_Es6RewriteSpreadExpressions()
	return static_Es6RewriteSpreadExpressions_arrayType
}
func StaticSet_Es6RewriteSpreadExpressions_ArrayType(v AstFactory_Type) {
	clinit_Es6RewriteSpreadExpressions()
	static_Es6RewriteSpreadExpressions_arrayType = v
}

var static_Es6RewriteSpreadExpressions_concatFnType AstFactory_Type

func StaticGet_Es6RewriteSpreadExpressions_ConcatFnType() AstFactory_Type {
	clinit_Es6RewriteSpreadExpressions()
	return static_Es6RewriteSpreadExpressions_concatFnType
}
func StaticSet_Es6RewriteSpreadExpressions_ConcatFnType(v AstFactory_Type) {
	clinit_Es6RewriteSpreadExpressions()
	static_Es6RewriteSpreadExpressions_concatFnType = v
}

var clinitDone_Es6RewriteSpreadExpressions uint32
var clinitState_Es6RewriteSpreadExpressions j2gort.ClinitState

func clinit_Es6RewriteSpreadExpressions() {
	if atomic.LoadUint32(&clinitDone_Es6RewriteSpreadExpressions) != 0 {
		return
	}
	clinitSlow_Es6RewriteSpreadExpressions()
}

func clinitSlow_Es6RewriteSpreadExpressions() {
	if !j2gort.EnterClinit(&clinitState_Es6RewriteSpreadExpressions, &clinitDone_Es6RewriteSpreadExpressions) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6RewriteSpreadExpressions, &clinitDone_Es6RewriteSpreadExpressions)
	static_Es6RewriteSpreadExpressions_TRANSPILED_FEATURES = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS())
	static_Es6RewriteSpreadExpressions_arrayType = AstFactory_S_Type_1_ColorId(pkg_colors.StaticGet_StandardColors_ARRAY_ID())
	static_Es6RewriteSpreadExpressions_concatFnType = AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT())
}

func (this_ *Es6RewriteSpreadExpressions) initFields_Es6RewriteSpreadExpressions() {
	_ = this_
}

// Java constructor: Es6RewriteSpreadExpressions(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6RewriteSpreadExpressions(compiler AbstractCompiler) *Es6RewriteSpreadExpressions {
	clinit_Es6RewriteSpreadExpressions()
	this_ := &Es6RewriteSpreadExpressions{}
	InitEs6RewriteSpreadExpressions(this_, compiler)
	return this_
}

func InitEs6RewriteSpreadExpressions(this_ *Es6RewriteSpreadExpressions, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6RewriteSpreadExpressions()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.namespace = compiler.GetTranspilationNamespace()
}

// Java method: Es6RewriteSpreadExpressions.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteSpreadExpressions) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_Es6RewriteSpreadExpressions_TRANSPILED_FEATURES(), []NodeTraversal_Callback{this_})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_Es6RewriteSpreadExpressions_TRANSPILED_FEATURES())
}

// Java method: Es6RewriteSpreadExpressions.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node current, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteSpreadExpressions) Visit_3(traversal *NodeTraversal, current pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = current
	_ = parent
	var token pkg_rhino.Token = current.GetToken()
	_ = token
	if ((token == pkg_rhino.StaticGet_Token_ARRAYLIT()) || (token == pkg_rhino.StaticGet_Token_NEW())) || (token == pkg_rhino.StaticGet_Token_CALL()) {
		{
			var child pkg_rhino.Node = current.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if child.IsSpread() {
					this_.priv_Es6RewriteSpreadExpressions_visitArrayLitOrCallWithSpread_2767(traversal, current)
					break
				}
			}
		}
	}
}

// Java method: Es6RewriteSpreadExpressions.visitArrayLitOrCallWithSpread(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node spreadParent)
func (this_ *Es6RewriteSpreadExpressions) priv_Es6RewriteSpreadExpressions_visitArrayLitOrCallWithSpread_2767(t *NodeTraversal, spreadParent pkg_rhino.Node) {
	_ = t
	_ = spreadParent
	if spreadParent.IsArrayLit() {
		this_.priv_Es6RewriteSpreadExpressions_visitArrayLitContainingSpread_2769(spreadParent)
	} else if spreadParent.IsCall() {
		this_.priv_Es6RewriteSpreadExpressions_visitCallContainingSpread_2770(t, spreadParent)
	} else {
		if !(spreadParent.IsNew()) {
			j2gort.CheckArgumentMsg(false, j2gort.BoxNode(spreadParent))
		}
		this_.priv_Es6RewriteSpreadExpressions_visitNewWithSpread_2772(spreadParent)
	}
}

// Java method: Es6RewriteSpreadExpressions.extractSpreadGroups(com.google.javascript.rhino.Node spreadParent)
func (this_ *Es6RewriteSpreadExpressions) priv_Es6RewriteSpreadExpressions_extractSpreadGroups_2768(spreadParent pkg_rhino.Node) j2gort.Ext_java_util_List {
	_ = spreadParent
	if !((spreadParent.IsCall() || spreadParent.IsArrayLit()) || spreadParent.IsNew()) {
		j2gort.CheckArgument(false)
	}
	var groups []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = groups
	var currGroup pkg_rhino.Node = pkg_rhino.Node(0)
	_ = currGroup
	{
		var currElement pkg_rhino.Node = spreadParent.RemoveFirstChild()
		_ = currElement
		for ; currElement != 0; currElement = spreadParent.RemoveFirstChild() {
			if currElement.IsSpread() {
				var spreadExpression pkg_rhino.Node = currElement.RemoveFirstChild()
				_ = spreadExpression
				if spreadExpression.IsArrayLit() {
					if currGroup == 0 {
						currGroup = spreadExpression
					} else {
						currGroup.AddChildrenToBack(spreadExpression.RemoveChildren())
					}
				} else {
					if currGroup != 0 {
						groups = append(groups, currGroup)
						currGroup = pkg_rhino.Node(0)
					}
					groups = append(groups, this_.astFactory.CreateJscompArrayFromIterableCall(spreadExpression, this_.namespace))
				}
			} else {
				if currGroup == 0 {
					currGroup = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{})
				}
				currGroup.AddChildToBack(currElement)
			}
		}
	}
	if currGroup != 0 {
		groups = append(groups, currGroup)
	}
	return j2gort.WrapTypedArrayList(groups)
}

// Java method: Es6RewriteSpreadExpressions.visitArrayLitContainingSpread(com.google.javascript.rhino.Node spreadParent)
func (this_ *Es6RewriteSpreadExpressions) priv_Es6RewriteSpreadExpressions_visitArrayLitContainingSpread_2769(spreadParent pkg_rhino.Node) {
	_ = spreadParent
	if !(spreadParent.IsArrayLit()) {
		j2gort.CheckArgument(false)
	}
	var groups j2gort.Ext_java_util_List = this_.priv_Es6RewriteSpreadExpressions_extractSpreadGroups_2768(spreadParent)
	_ = groups
	var baseArrayLit pkg_rhino.Node = pkg_rhino.Node(0)
	_ = baseArrayLit
	if j2gort.CastNode[pkg_rhino.Node](groups.Get(int32(0))).IsArrayLit() {
		baseArrayLit = j2gort.CastNode[pkg_rhino.Node](groups.Remove_1_int32(int32(0)))
	} else {
		baseArrayLit = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{})
	}
	var joinedGroups pkg_rhino.Node = pkg_rhino.Node(0)
	_ = joinedGroups
	if groups.IsEmpty() {
		joinedGroups = baseArrayLit
	} else {
		var concat pkg_rhino.Node = this_.astFactory.CreateGetProp(baseArrayLit, "concat", StaticGet_Es6RewriteSpreadExpressions_ConcatFnType())
		_ = concat
		joinedGroups = this_.astFactory.CreateCall(concat, AstFactory_S_Type_1_Node(spreadParent), func() []pkg_rhino.Node {
			arr_4 := groups.ToArray_1_slice_any(func() []any {
				arr_1 := make([]pkg_rhino.Node, int32(0))
				if arr_1 == nil {
					return nil
				}
				out_2 := make([]any, len(arr_1))
				for i_, el_3 := range arr_1 {
					out_2[i_] = j2gort.BoxNode(el_3)
				}
				return out_2
			}())
			if arr_4 == nil {
				return nil
			}
			out_5 := make([]pkg_rhino.Node, len(arr_4))
			for i_, el_6 := range arr_4 {
				out_5[i_] = j2gort.CastNode[pkg_rhino.Node](el_6)
			}
			return out_5
		}())
	}
	joinedGroups.SrcrefTreeIfMissing(spreadParent)
	spreadParent.ReplaceWith(joinedGroups)
	this_.compiler.ReportChangeToEnclosingScope(joinedGroups)
}

// Java method: Es6RewriteSpreadExpressions.visitCallContainingSpread(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node spreadParent)
func (this_ *Es6RewriteSpreadExpressions) priv_Es6RewriteSpreadExpressions_visitCallContainingSpread_2770(t *NodeTraversal, spreadParent pkg_rhino.Node) {
	_ = t
	_ = spreadParent
	if !(spreadParent.IsCall()) {
		j2gort.CheckArgument(false)
	}
	var callee pkg_rhino.Node = spreadParent.GetFirstChild()
	_ = callee
	if !(!(callee.IsSuper())) {
		j2gort.CheckStateStr(false, "Cannot spread into super calls")
	}
	var calleeMayHaveSideEffects bool = this_.compiler.GetAstAnalyzer().MayHaveSideEffects(callee)
	_ = calleeMayHaveSideEffects
	callee.Detach()
	for callee.IsCast() {
		callee = callee.RemoveFirstChild()
	}
	var joinedGroups pkg_rhino.Node = pkg_rhino.Node(0)
	_ = joinedGroups
	if spreadParent.HasOneChild() && this_.priv_Es6RewriteSpreadExpressions_isSpreadOfArguments_2771(spreadParent.GetOnlyChild()) {
		joinedGroups = spreadParent.RemoveFirstChild().RemoveFirstChild()
	} else {
		var groups j2gort.Ext_java_util_List = this_.priv_Es6RewriteSpreadExpressions_extractSpreadGroups_2768(spreadParent)
		_ = groups
		if !(!(groups.IsEmpty())) {
			j2gort.CheckState(false)
		}
		if groups.Size() == int32(1) {
			joinedGroups = j2gort.CastNode[pkg_rhino.Node](groups.Remove_1_int32(int32(0)))
		} else {
			var baseArrayLit pkg_rhino.Node
			_ = baseArrayLit
			if j2gort.CastNode[pkg_rhino.Node](groups.Get(int32(0))).IsArrayLit() {
				baseArrayLit = j2gort.CastNode[pkg_rhino.Node](groups.Remove_1_int32(int32(0)))
			} else {
				baseArrayLit = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{})
			}
			var concat pkg_rhino.Node = this_.astFactory.CreateGetProp(baseArrayLit, "concat", StaticGet_Es6RewriteSpreadExpressions_ConcatFnType())
			_ = concat
			joinedGroups = this_.astFactory.CreateCall(concat, StaticGet_Es6RewriteSpreadExpressions_ArrayType(), func() []pkg_rhino.Node {
				arr_10 := groups.ToArray_1_slice_any(func() []any {
					arr_7 := make([]pkg_rhino.Node, int32(0))
					if arr_7 == nil {
						return nil
					}
					out_8 := make([]any, len(arr_7))
					for i_, el_9 := range arr_7 {
						out_8[i_] = j2gort.BoxNode(el_9)
					}
					return out_8
				}())
				if arr_10 == nil {
					return nil
				}
				out_11 := make([]pkg_rhino.Node, len(arr_10))
				for i_, el_12 := range arr_10 {
					out_11[i_] = j2gort.CastNode[pkg_rhino.Node](el_12)
				}
				return out_11
			}())
		}
	}
	var isFreeCall bool = spreadParent.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL())
	_ = isFreeCall
	var callToApply pkg_rhino.Node = pkg_rhino.Node(0)
	_ = callToApply
	if (calleeMayHaveSideEffects && (callee.IsGetProp() || callee.IsGetElem())) && !(isFreeCall) {
		var freshVar pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$spread$args" + j2gort.StrValString(this_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput()))), AstFactory_S_Type_1_Node(callee.GetFirstChild()))
		_ = freshVar
		var freshVarDeclaration pkg_rhino.Node = pkg_rhino.IR_Var_1(freshVar.CloneTree_0())
		_ = freshVarDeclaration
		var statementContainingSpread pkg_rhino.Node = NodeUtil_GetEnclosingStatement(spreadParent)
		_ = statementContainingSpread
		freshVarDeclaration.SrcrefTreeIfMissing(statementContainingSpread)
		freshVarDeclaration.InsertBefore(statementContainingSpread)
		callee.AddChildToFront(this_.astFactory.CreateAssign_2_Node_Node(freshVar.CloneTree_0(), callee.RemoveFirstChild()))
		callToApply = this_.astFactory.CreateCallWithUnknownType(this_.astFactory.CreateGetPropWithUnknownType(callee, "apply"), []pkg_rhino.Node{freshVar.CloneTree_0(), joinedGroups})
	} else {
		var context pkg_rhino.Node
		_ = context
		if (callee.IsGetProp() || callee.IsGetElem()) && !(isFreeCall) {
			context = callee.GetFirstChild().CloneTree_0()
		} else {
			context = this_.astFactory.CreateNull()
		}
		callToApply = this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(callee, "apply"), AstFactory_S_Type_1_Node(spreadParent), []pkg_rhino.Node{context, joinedGroups})
	}
	callToApply.SetColor(spreadParent.GetColor())
	callToApply.SrcrefTreeIfMissing(spreadParent)
	spreadParent.ReplaceWith(callToApply)
	this_.compiler.ReportChangeToEnclosingScope(callToApply)
}

// Java method: Es6RewriteSpreadExpressions.isSpreadOfArguments(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteSpreadExpressions) priv_Es6RewriteSpreadExpressions_isSpreadOfArguments_2771(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsSpread() && n.GetOnlyChild().MatchesName_1_string("arguments"))
}

// Java method: Es6RewriteSpreadExpressions.visitNewWithSpread(com.google.javascript.rhino.Node spreadParent)
func (this_ *Es6RewriteSpreadExpressions) priv_Es6RewriteSpreadExpressions_visitNewWithSpread_2772(spreadParent pkg_rhino.Node) {
	_ = spreadParent
	if !(spreadParent.IsNew()) {
		j2gort.CheckArgument(false)
	}
	var callee pkg_rhino.Node = spreadParent.RemoveFirstChild()
	_ = callee
	var groups j2gort.Ext_java_util_List = this_.priv_Es6RewriteSpreadExpressions_extractSpreadGroups_2768(spreadParent)
	_ = groups
	var baseArrayLit pkg_rhino.Node = pkg_rhino.Node(0)
	_ = baseArrayLit
	if j2gort.CastNode[pkg_rhino.Node](groups.Get(int32(0))).IsArrayLit() {
		baseArrayLit = j2gort.CastNode[pkg_rhino.Node](groups.Remove_1_int32(int32(0)))
	} else {
		baseArrayLit = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{})
	}
	baseArrayLit.AddChildToFront(this_.astFactory.CreateNull())
	var joinedGroups pkg_rhino.Node
	_ = joinedGroups
	if groups.IsEmpty() {
		joinedGroups = baseArrayLit
	} else {
		joinedGroups = this_.astFactory.CreateCall(this_.astFactory.CreateGetProp(baseArrayLit, "concat", StaticGet_Es6RewriteSpreadExpressions_ConcatFnType()), StaticGet_Es6RewriteSpreadExpressions_ArrayType(), func() []pkg_rhino.Node {
			arr_16 := groups.ToArray_1_slice_any(func() []any {
				arr_13 := make([]pkg_rhino.Node, int32(0))
				if arr_13 == nil {
					return nil
				}
				out_14 := make([]any, len(arr_13))
				for i_, el_15 := range arr_13 {
					out_14[i_] = j2gort.BoxNode(el_15)
				}
				return out_14
			}())
			if arr_16 == nil {
				return nil
			}
			out_17 := make([]pkg_rhino.Node, len(arr_16))
			for i_, el_18 := range arr_16 {
				out_17[i_] = j2gort.CastNode[pkg_rhino.Node](el_18)
			}
			return out_17
		}())
	}
	if pkg_parser.StaticGet_FeatureSet_ES3().Contains_1_FeatureSet(this_.compiler.GetOptions().GetOutputFeatureSet()) {
		TranspilationUtil_CannotConvert(this_.compiler, spreadParent, "\"...\" passed to a constructor (consider using --language_out=ES5)")
	}
	var bindApply pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreatePrototypeAccess(this_.astFactory.CreateName_2_StaticScope_string(this_.namespace, "Function")), "bind"), "apply")
	_ = bindApply
	var result pkg_rhino.Node = pkg_rhino.IR_NewNode(this_.astFactory.CreateCallWithUnknownType(bindApply, []pkg_rhino.Node{callee, joinedGroups}), []pkg_rhino.Node{}).SetColor(spreadParent.GetColor())
	_ = result
	result.SrcrefTreeIfMissing(spreadParent)
	spreadParent.ReplaceWith(result)
	this_.compiler.ReportChangeToEnclosingScope(result)
}
func (*Es6RewriteSpreadExpressions) AlwaysShouldTraverse_()      {}
func (*Es6RewriteSpreadExpressions) CanReleaseSyntacticScopes_() {}

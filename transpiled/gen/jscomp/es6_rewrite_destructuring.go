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
// Port of com/google/javascript/jscomp/Es6RewriteDestructuring.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_es6_rewrite_destructuring_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_es6_rewrite_destructuring_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_10 uint32) bool {
		n := pkg_rhino.Node(rawNode_10)
		_ = n
		return (n.IsAwait() || n.IsYield())
	})
}

// Java class: com.google.javascript.jscomp.Es6RewriteDestructuring
type Es6RewriteDestructuring struct {
	compiler                     AbstractCompiler                                           `gson:"compiler"`
	astFactory                   *AstFactory                                                `gson:"astFactory"`
	namespace                    pkg_rhino.StaticScope                                      `gson:"namespace"`
	featuresToTriggerRunningPass *pkg_parser.FeatureSet                                     `gson:"featuresToTriggerRunningPass"`
	featuresToMarkAsRemoved      *pkg_parser.FeatureSet                                     `gson:"featuresToMarkAsRemoved"`
	patternNestingStack          j2gort.Deque[*Es6RewriteDestructuring_PatternNestingLevel] `gson:"patternNestingStack"`
	destructuringVarCounter      int32                                                      `gson:"destructuringVarCounter"`
	idHash_                      atomic.Int32
	rewriteMode                  Es6RewriteDestructuring_ObjectDestructuringRewriteMode `gson:"rewriteMode"`
	useJscompAssign              bool                                                   `gson:"useJscompAssign"`
}

func (this_ *Es6RewriteDestructuring) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteDestructuring) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteDestructuring) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteDestructuring", "Es6RewriteDestructuring"
}
func (this_ *Es6RewriteDestructuring) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteDestructuring) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteDestructuring", "com.google.javascript.jscomp.Es6RewriteDestructuring", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteDestructuring](v_)
	})
}
func (*Es6RewriteDestructuring) Implements_jscomp_NodeTraversal_Callback() {}
func (*Es6RewriteDestructuring) Implements_jscomp_CompilerPass()           {}

func StaticGet_Es6RewriteDestructuring_DESTRUCTURING_TEMP_VAR() string {
	return "$jscomp$destructuring$var"
}
func (this_ *Es6RewriteDestructuring) initFields_Es6RewriteDestructuring() {
	_ = this_
	this_.patternNestingStack = j2gort.Deque[*Es6RewriteDestructuring_PatternNestingLevel]{}

}

// Java constructor: Es6RewriteDestructuring(com.google.javascript.jscomp.Es6RewriteDestructuring.Builder builder)
func NewEs6RewriteDestructuring(builder *Es6RewriteDestructuring_Builder) *Es6RewriteDestructuring {
	this_ := &Es6RewriteDestructuring{}
	InitEs6RewriteDestructuring(this_, builder)
	return this_
}

func InitEs6RewriteDestructuring(this_ *Es6RewriteDestructuring, builder *Es6RewriteDestructuring_Builder) {
	_ = builder
	this_.initFields_Es6RewriteDestructuring()
	this_.compiler = builder.compiler
	this_.rewriteMode = builder.rewriteMode
	this_.astFactory = this_.compiler.CreateAstFactory()
	this_.namespace = this_.compiler.GetTranspilationNamespace()
	this_.useJscompAssign = !(this_.compiler.GetOptions().GetOutputFeatureSet().Contains_1_FeatureSet(pkg_parser.StaticGet_FeatureSet_ES2015()))
	switch (this_.rewriteMode).Ordinal() {
	case int32(0) /* ObjectDestructuringRewriteMode.REWRITE_ALL_OBJECT_PATTERNS */ :
		this_.featuresToTriggerRunningPass = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7237([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_DEFAULT_PARAMETERS(), pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_DESTRUCTURING(), pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_PATTERN_REST(), pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_DESTRUCTURING()})
		this_.featuresToMarkAsRemoved = this_.featuresToTriggerRunningPass.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST())
	case int32(1) /* ObjectDestructuringRewriteMode.REWRITE_OBJECT_REST */ :
		this_.featuresToTriggerRunningPass = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST())
		this_.featuresToMarkAsRemoved = this_.featuresToTriggerRunningPass
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Es6RewriteDestructuring cannot handle ObjectDestructuringRewriteMode "+j2gort.StrVal(this_.rewriteMode))))
	}
}

// Java method: Es6RewriteDestructuring.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteDestructuring) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.patternNestingStack.IsEmpty()) {
		j2gort.CheckState(false)
	}
	NodeTraversal_Traverse(this_.compiler, root, this_)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, this_.featuresToMarkAsRemoved)
	if !(this_.patternNestingStack.IsEmpty()) {
		j2gort.CheckState(false)
	}
}

// Java method: Es6RewriteDestructuring.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteDestructuring) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(n)
		_ = scriptFeatures
		return scriptFeatures.ContainsAtLeastOneOf(this_.featuresToTriggerRunningPass)
	}
	switch n.TokenSwitchOrd_() {
	case int32(49) /* Token.PARAM_LIST */ :
		this_.priv_Es6RewriteDestructuring_pullDestructuringOutOfParams_2629(n, parent)
	case int32(102) /* Token.ARRAY_PATTERN */, int32(103) /* Token.OBJECT_PATTERN */ :
		var hasRest bool = ((n.IsObjectPattern() && n.HasChildren()) && n.GetLastChild().IsRest())
		_ = hasRest
		if !(this_.patternNestingStack.IsEmpty()) && hasRest {
			{
				cur_1 := this_.patternNestingStack.Iter()
				for cur_1.HasNext() {

					level := cur_1.Next()
					_ = level
					if level.hasNestedObjectRest {
						break
					}
					{
						r_3 := level
						r_3.hasNestedObjectRest = true
					}
				}
			}
			{
				r_4 := this_.patternNestingStack.PeekLast()
				r_4.hasNestedObjectRest = true
			}
		}
		this_.patternNestingStack.AddLast(NewEs6RewriteDestructuring_PatternNestingLevel(n, hasRest))
	default:
	}
	return true
}

// Java method: Es6RewriteDestructuring.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteDestructuring) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(102) /* Token.ARRAY_PATTERN */, int32(103) /* Token.OBJECT_PATTERN */ :
		this_.priv_Es6RewriteDestructuring_visitPattern_2634(t, n)
		if n == this_.patternNestingStack.GetLast().pattern {
			this_.patternNestingStack.RemoveLast()
		}
	default:
	}
}

// Java method: Es6RewriteDestructuring.pullDestructuringOutOfParams(com.google.javascript.rhino.Node paramList, com.google.javascript.rhino.Node function)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_pullDestructuringOutOfParams_2629(paramList pkg_rhino.Node, function pkg_rhino.Node) {
	_ = paramList
	_ = function
	var insertSpot pkg_rhino.Node = pkg_rhino.Node(0)
	_ = insertSpot
	var body pkg_rhino.Node = function.GetLastChild()
	_ = body
	var next pkg_rhino.Node = pkg_rhino.Node(0)
	_ = next
	{
		var param pkg_rhino.Node = paramList.GetFirstChild()
		_ = param
		for ; param != 0; param = next {
			next = param.GetNext()
			if param.IsDefaultValue() {
				var nameOrPattern pkg_rhino.Node = param.RemoveFirstChild()
				_ = nameOrPattern
				var jsDoc *pkg_rhino.JSDocInfo = nameOrPattern.GetJSDocInfo()
				_ = jsDoc
				nameOrPattern.SetJSDocInfo(nil)
				var defaultValue pkg_rhino.Node = param.RemoveFirstChild()
				_ = defaultValue
				var newParam pkg_rhino.Node = pkg_rhino.Node(0)
				_ = newParam
				var isNoop bool = false
				_ = isNoop
				if !(nameOrPattern.IsName()) {
				} else if defaultValue.IsName() {
					{

						sr_5 := "undefined"
						se_6 := defaultValue.GetString()
						isNoop = j2gort.CheckNotNullString(sr_5) == se_6
					}

				} else if defaultValue.IsVoid() {
					isNoop = NodeUtil_IsImmutableValue(defaultValue.GetFirstChild())
				}
				if isNoop {
					newParam = nameOrPattern.CloneTree_0()
				} else {
					if nameOrPattern.IsName() {
						newParam = nameOrPattern
					} else {
						newParam = this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(this_.priv_Es6RewriteDestructuring_getTempVariableName_2633(), AstFactory_S_Type_1_Node(nameOrPattern))
					}
					var lhs pkg_rhino.Node = nameOrPattern.CloneTree_0()
					_ = lhs
					var rhs pkg_rhino.Node = this_.priv_Es6RewriteDestructuring_defaultValueHook_2646(newParam.CloneTree_0(), defaultValue)
					_ = rhs
					var newStatement pkg_rhino.Node
					_ = newStatement
					if nameOrPattern.IsName() {
						newStatement = pkg_rhino.IR_ExprResult(this_.astFactory.CreateAssign_2_Node_Node(lhs, rhs))
					} else {
						newStatement = pkg_rhino.IR_Var_2(lhs, rhs)
					}
					newStatement.SrcrefTreeIfMissing(param)
					if insertSpot == 0 {
						insertSpot = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(body)
						if insertSpot != 0 {
							newStatement.InsertBefore(insertSpot)
						} else {
							body.AddChildToBack(newStatement)
						}
					} else {
						newStatement.InsertAfter(insertSpot)
					}
					insertSpot = newStatement
					if newStatement.GetFirstChild().IsDestructuringLhs() {
						insertSpot = this_.priv_Es6RewriteDestructuring_rewriteVarDestructuringDeclaration_2631(newStatement.GetFirstChild())
					}
				}
				param.ReplaceWith(newParam)
				newParam.SetJSDocInfo(jsDoc)
				this_.compiler.ReportChangeToChangeScope(function)
			} else if param.IsDestructuringPattern() {
				insertSpot = this_.priv_Es6RewriteDestructuring_replacePatternParamWithTempVar_2630(function, insertSpot, param, this_.priv_Es6RewriteDestructuring_getTempVariableName_2633())
				this_.compiler.ReportChangeToChangeScope(function)
			} else if param.IsRest() && param.GetFirstChild().IsDestructuringPattern() {
				insertSpot = this_.priv_Es6RewriteDestructuring_replacePatternParamWithTempVar_2630(function, insertSpot, param.GetFirstChild(), this_.priv_Es6RewriteDestructuring_getTempVariableName_2633())
				this_.compiler.ReportChangeToChangeScope(function)
			}
		}
	}
}

// Java method: Es6RewriteDestructuring.replacePatternParamWithTempVar(com.google.javascript.rhino.Node function, com.google.javascript.rhino.Node insertSpot, com.google.javascript.rhino.Node patternParam, java.lang.String tempVarName)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_replacePatternParamWithTempVar_2630(function pkg_rhino.Node, insertSpot pkg_rhino.Node, patternParam pkg_rhino.Node, tempVarName string) pkg_rhino.Node {
	_ = function
	_ = insertSpot
	_ = patternParam
	_ = tempVarName
	var paramType AstFactory_Type = AstFactory_S_Type_1_Node(patternParam)
	_ = paramType
	var newParam pkg_rhino.Node = this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, paramType)
	_ = newParam
	newParam.SetJSDocInfo(patternParam.GetJSDocInfo())
	patternParam.ReplaceWith(newParam)
	var newDecl pkg_rhino.Node = pkg_rhino.IR_Var_2(patternParam, this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, paramType))
	_ = newDecl
	newDecl.SrcrefTreeIfMissing(patternParam)
	if insertSpot == 0 {
		insertSpot = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(function.GetLastChild())
		if insertSpot != 0 {
			newDecl.InsertBefore(insertSpot)
		} else {
			function.GetLastChild().AddChildToBack(newDecl)
		}
	} else {
		newDecl.InsertAfter(insertSpot)
	}
	return this_.priv_Es6RewriteDestructuring_rewriteVarDestructuringDeclaration_2631(newDecl.GetFirstChild())
}

// Java method: Es6RewriteDestructuring.rewriteVarDestructuringDeclaration(com.google.javascript.rhino.Node destructuringLhs)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_rewriteVarDestructuringDeclaration_2631(destructuringLhs pkg_rhino.Node) pkg_rhino.Node {
	_ = destructuringLhs
	if !(destructuringLhs.IsDestructuringLhs()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(destructuringLhs))
	}
	if !(!(destructuringLhs.GetGrandparent().IsExport())) {
		j2gort.CheckStateStr(false, "Export destructuring declarations not expected inside a function body.")
	}
	var var_ pkg_rhino.Node = destructuringLhs.GetParent()
	_ = var_
	NodeUtil_VisitLhsNodesInNode(destructuringLhs, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_7 any) {
		name := j2gort.CastNode[pkg_rhino.Node](lp_7)
		_ = name
		if !(name.IsName()) {
			j2gort.CheckStateFmt(false, "lhs in destructuring declaration should be a simple name. (%s)", j2gort.BoxNode(name))
		}
		var newName pkg_rhino.Node = pkg_rhino.IR_Name(name.GetString()).Srcref(name)
		_ = newName
		if name.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()) {
			newName.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
		}
		var newVar pkg_rhino.Node = pkg_rhino.IR_Var_1(newName).Srcref(name)
		_ = newVar
		newVar.InsertBefore(var_)
	}))
	var destructuringPattern pkg_rhino.Node = destructuringLhs.RemoveFirstChild()
	_ = destructuringPattern
	if !(destructuringPattern.IsDestructuringPattern()) {
		j2gort.CheckStateFmt(false, "Expected destructuring pattern but got %s", j2gort.BoxString(destructuringPattern.ToStringTree()))
	}
	var rhs pkg_rhino.Node = destructuringLhs.RemoveFirstChild()
	_ = rhs
	var assign pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(destructuringPattern, rhs)
	_ = assign
	assign.Srcref(var_)
	var expr pkg_rhino.Node = this_.astFactory.ExprResult(assign)
	_ = expr
	expr.Srcref(var_)
	var_.ReplaceWith(expr)
	return expr
}

// Java method: Es6RewriteDestructuring.createTempVarNameNode(java.lang.String name, com.google.javascript.jscomp.AstFactory.Type type)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_createTempVarNameNode_2632(name string, type_ AstFactory_Type) pkg_rhino.Node {
	_ = name
	_ = type_
	return this_.astFactory.CreateConstantName(name, type_)
}

// Java method: Es6RewriteDestructuring.getTempVariableName()
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_getTempVariableName_2633() string {
	return ("$jscomp$destructuring$var" + j2gort.StrValInt32(func() int32 {
		old_8 := this_.destructuringVarCounter
		nw_9 := old_8 + 1
		this_.destructuringVarCounter = nw_9
		return old_8
	}()))
}

// Java method: Es6RewriteDestructuring.visitPattern(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node pattern)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_visitPattern_2634(t *NodeTraversal, pattern pkg_rhino.Node) {
	_ = t
	_ = pattern
	var parent pkg_rhino.Node = pattern.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(104) /* Token.DESTRUCTURING_LHS */ :
		{
			var declaration pkg_rhino.Node = parent.GetParent()
			_ = declaration
			var declarationParent pkg_rhino.Node = declaration.GetParent()
			_ = declarationParent
			if declarationParent.IsVanillaFor() {
				this_.priv_Es6RewriteDestructuring_visitDestructuringPatternInVanillaForInnerVars_2642(t, pattern)
			} else if NodeUtil_IsEnhancedFor(declarationParent) {
				this_.priv_Es6RewriteDestructuring_visitDestructuringPatternInEnhancedForInnerVars_2643(pattern)
			} else {
				this_.priv_Es6RewriteDestructuring_replacePattern_2635(t, pattern, pattern.GetNext(), declaration, declaration)
			}
		}
	case int32(51) /* Token.ASSIGN */ :
		if parent.GetParent().IsExprResult() {
			this_.priv_Es6RewriteDestructuring_replacePattern_2635(t, pattern, pattern.GetNext(), parent, parent.GetParent())
		} else {
			this_.priv_Es6RewriteDestructuring_wrapAssignOrDestructuringInCallToArrow_2640(t, parent)
		}
	case int32(126) /* Token.OBJECT_REST */, int32(125) /* Token.ITER_REST */, int32(100) /* Token.STRING_KEY */, int32(102) /* Token.ARRAY_PATTERN */, int32(134) /* Token.DEFAULT_VALUE */, int32(129) /* Token.COMPUTED_PROP */ :
	case int32(112) /* Token.FOR_OF */, int32(81) /* Token.FOR_IN */, int32(113) /* Token.FOR_AWAIT_OF */ :
		this_.priv_Es6RewriteDestructuring_visitDestructuringPatternInEnhancedForWithOuterVars_2644(pattern)
	case int32(86) /* Token.CATCH */ :
		this_.priv_Es6RewriteDestructuring_visitDestructuringPatternInCatch_2645(t, pattern)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected parent")))
	}
}

// Java method: Es6RewriteDestructuring.replacePattern(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node pattern, com.google.javascript.rhino.Node rhs, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node nodeToDetach)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_replacePattern_2635(t *NodeTraversal, pattern pkg_rhino.Node, rhs pkg_rhino.Node, parent pkg_rhino.Node, nodeToDetach pkg_rhino.Node) {
	_ = t
	_ = pattern
	_ = rhs
	_ = parent
	_ = nodeToDetach
	if !(NodeUtil_IsStatement(nodeToDetach)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(nodeToDetach))
	}
	switch pattern.TokenSwitchOrd_() {
	case int32(102) /* Token.ARRAY_PATTERN */ :
		this_.priv_Es6RewriteDestructuring_replaceArrayPattern_2639(t, pattern, rhs, parent, nodeToDetach)
	case int32(103) /* Token.OBJECT_PATTERN */ :
		this_.priv_Es6RewriteDestructuring_replaceObjectPattern_2636(t, pattern, rhs, parent, nodeToDetach)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected")))
	}
}

// Java method: Es6RewriteDestructuring.replaceObjectPattern(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node objectPattern, com.google.javascript.rhino.Node rhs, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node nodeToDetach)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_replaceObjectPattern_2636(t *NodeTraversal, objectPattern pkg_rhino.Node, rhs pkg_rhino.Node, parent pkg_rhino.Node, nodeToDetach pkg_rhino.Node) {
	_ = t
	_ = objectPattern
	_ = rhs
	_ = parent
	_ = nodeToDetach
	var tempVarName string = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
	_ = tempVarName
	var tempVarType AstFactory_Type = AstFactory_S_Type_1_Node(objectPattern)
	_ = tempVarType
	var restTempVarName string = j2gort.NullString
	_ = restTempVarName
	var propsToDeleteForRest j2gort.Ext_java_util_ArrayList = nil
	_ = propsToDeleteForRest
	if objectPattern.HasChildren() && objectPattern.GetLastChild().IsRest() {
		propsToDeleteForRest = j2gort.NewExt_java_util_ArrayList()
		restTempVarName = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
	} else if this_.rewriteMode == StaticGet_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST() {
		if !(this_.patternNestingStack.PeekLast().hasNestedObjectRest) {
			this_.destructuringVarCounter = this_.destructuringVarCounter - 1
			return
		}
	}
	var tempDecl pkg_rhino.Node = pkg_rhino.IR_Var_2(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, tempVarType), rhs.Detach()).SrcrefTreeIfMissing(objectPattern)
	_ = tempDecl
	if parent.IsConst() {
		var jsDoc *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
		_ = jsDoc
		jsDoc.RecordConstancy()
		tempDecl.SetJSDocInfo(jsDoc.Build_0())
	}
	tempDecl.InsertBefore(nodeToDetach)
	{
		var child pkg_rhino.Node = objectPattern.GetFirstChild()
		_ = child
		var next pkg_rhino.Node = pkg_rhino.Node(0)
		_ = next
		for ; child != 0; child = next {
			next = child.GetNext()
			var newLHS pkg_rhino.Node = pkg_rhino.Node(0)
			_ = newLHS
			var newRHS pkg_rhino.Node = pkg_rhino.Node(0)
			_ = newRHS
			if child.IsStringKey() {
				var tempVarNameNode pkg_rhino.Node = this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, tempVarType)
				_ = tempVarNameNode
				var getprop pkg_rhino.Node
				_ = getprop
				if child.IsQuotedStringKey() {
					getprop = this_.astFactory.CreateGetElem(tempVarNameNode, this_.astFactory.CreateString(child.GetString()))
				} else {
					getprop = this_.astFactory.CreateGetProp(tempVarNameNode, child.GetString(), tempVarType)
				}
				var value pkg_rhino.Node = child.RemoveFirstChild()
				_ = value
				if !(value.IsDefaultValue()) {
					newLHS = value
					newRHS = getprop
				} else {
					newLHS = value.RemoveFirstChild()
					var defaultValue pkg_rhino.Node = value.RemoveFirstChild()
					_ = defaultValue
					var intermediateTempVarName string = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
					_ = intermediateTempVarName
					var intermediateDecl pkg_rhino.Node = pkg_rhino.IR_Var_2(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(intermediateTempVarName, AstFactory_S_Type_1_Node(getprop)), getprop)
					_ = intermediateDecl
					intermediateDecl.SrcrefTreeIfMissing(child)
					intermediateDecl.InsertBefore(nodeToDetach)
					newRHS = this_.priv_Es6RewriteDestructuring_defaultValueHook_2646(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(intermediateTempVarName, AstFactory_S_Type_1_Node(getprop)), defaultValue)
				}
				if propsToDeleteForRest != nil {
					propsToDeleteForRest.Add_1(j2gort.BoxNode(child))
				}
			} else if child.IsComputedProp() {
				var hasDefault bool = child.GetLastChild().IsDefaultValue()
				_ = hasDefault
				var defaultValue_2 pkg_rhino.Node = pkg_rhino.Node(0)
				_ = defaultValue_2
				var propExpr pkg_rhino.Node = child.RemoveFirstChild()
				_ = propExpr
				if hasDefault {
					var defaultNode pkg_rhino.Node = child.GetLastChild()
					_ = defaultNode
					newLHS = defaultNode.RemoveFirstChild()
					defaultValue_2 = defaultNode.RemoveFirstChild()
				} else {
					newLHS = child.RemoveFirstChild()
					defaultValue_2 = pkg_rhino.Node(0)
				}
				if propsToDeleteForRest != nil {
					var exprEvalTempVarName string = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
					_ = exprEvalTempVarName
					var exprEvalTempVarModel pkg_rhino.Node = this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(exprEvalTempVarName, AstFactory_S_Type_1_Node(propExpr))
					_ = exprEvalTempVarModel
					var exprEvalDecl pkg_rhino.Node = pkg_rhino.IR_Var_2(exprEvalTempVarModel.CloneNode_0(), propExpr)
					_ = exprEvalDecl
					exprEvalDecl.SrcrefTreeIfMissing(child)
					exprEvalDecl.InsertBefore(nodeToDetach)
					propExpr = exprEvalTempVarModel.CloneNode_0()
					propsToDeleteForRest.Add_1(j2gort.BoxNode(exprEvalTempVarModel.CloneNode_0()))
				}
				if hasDefault {
					var getelem pkg_rhino.Node = this_.astFactory.CreateGetElem(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, tempVarType), propExpr)
					_ = getelem
					var intermediateTempVarName_2 string = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
					_ = intermediateTempVarName_2
					var intermediateDecl_2 pkg_rhino.Node = pkg_rhino.IR_Var_2(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(intermediateTempVarName_2, AstFactory_S_Type_1_Node(getelem)), getelem)
					_ = intermediateDecl_2
					intermediateDecl_2.SrcrefTreeIfMissing(child)
					intermediateDecl_2.InsertBefore(nodeToDetach)
					newRHS = this_.priv_Es6RewriteDestructuring_defaultValueHook_2646(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(intermediateTempVarName_2, AstFactory_S_Type_1_Node(getelem)), defaultValue_2)
				} else {
					newRHS = this_.astFactory.CreateGetElem(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, AstFactory_S_Type_1_Node(newLHS)), propExpr)
				}
			} else if child.IsRest() {
				if next != 0 {
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "object rest may not be followed by any properties")))
				}
				var assignCall pkg_rhino.Node
				_ = assignCall
				if this_.useJscompAssign {
					assignCall = this_.astFactory.CreateJscompAssignCall(this_.namespace, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()), []pkg_rhino.Node{this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{}), this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, tempVarType)})
				} else {
					assignCall = this_.astFactory.CreateObjectDotAssignCall(this_.namespace, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()), []pkg_rhino.Node{this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{}), this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, tempVarType)})
				}
				var restTempDecl pkg_rhino.Node = pkg_rhino.IR_Var_2(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(restTempVarName, tempVarType), assignCall)
				_ = restTempDecl
				restTempDecl.SrcrefTreeIfMissing(objectPattern)
				restTempDecl.InsertBefore(nodeToDetach)
				var restName pkg_rhino.Node = child.GetOnlyChild()
				_ = restName
				if javalang.StringStartsWith(j2gort.CheckNotNullString(restName.GetString()), "$jscomp$destructuring$var") {
					if !(restName.IsName()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(restName))
					}
					newLHS = this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(restName.GetString(), AstFactory_S_Type_1_Node(restName))
				} else {
					newLHS = restName.Detach()
				}
				newRHS = this_.priv_Es6RewriteDestructuring_objectPatternRestRHS_2637(objectPattern, child, restTempVarName, propsToDeleteForRest)
			} else {
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected child")))
			}
			var newNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = newNode
			if NodeUtil_IsNameDeclaration(parent) {
				if parent.IsConst() {
					newLHS.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
				}
				newNode = pkg_rhino.IR_Declaration_3(newLHS, newRHS, parent.GetToken())
			} else if parent.IsAssign() {
				newNode = pkg_rhino.IR_ExprResult(this_.astFactory.CreateAssign_2_Node_Node(newLHS, newRHS))
			} else {
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "not reached")))
			}
			newNode.SrcrefTreeIfMissing(child)
			newNode.InsertBefore(nodeToDetach)
			this_.Visit_3(t, newLHS, newLHS.GetParent())
		}
	}
	nodeToDetach.Detach()
	t.ReportCodeChange_0()
}

// Java method: Es6RewriteDestructuring.objectPatternRestRHS(com.google.javascript.rhino.Node objectPattern, com.google.javascript.rhino.Node rest, java.lang.String restTempVarName, java.util.ArrayList<com.google.javascript.rhino.Node> statedProperties)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_objectPatternRestRHS_2637(objectPattern pkg_rhino.Node, rest pkg_rhino.Node, restTempVarName string, statedProperties j2gort.Ext_java_util_ArrayList) pkg_rhino.Node {
	_ = objectPattern
	_ = rest
	_ = restTempVarName
	_ = statedProperties
	if !(objectPattern.GetLastChild() == rest) {
		j2gort.CheckArgument(false)
	}
	var restTempVarModel pkg_rhino.Node = this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(restTempVarName, AstFactory_S_Type_1_Node(objectPattern))
	_ = restTempVarModel
	var result pkg_rhino.Node = restTempVarModel.CloneNode_0()
	_ = result
	if !(statedProperties.IsEmpty()) {
		var propItr j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](statedProperties.Iterator())
		_ = propItr
		var comma pkg_rhino.Node = this_.priv_Es6RewriteDestructuring_deletionNodeForRestProperty_2638(restTempVarModel.CloneNode_0(), j2gort.CastNode[pkg_rhino.Node](propItr.Next()))
		_ = comma
		for propItr.HasNext() {
			comma = this_.astFactory.CreateComma(comma, this_.priv_Es6RewriteDestructuring_deletionNodeForRestProperty_2638(restTempVarModel.CloneNode_0(), j2gort.CastNode[pkg_rhino.Node](propItr.Next())))
		}
		result = this_.astFactory.CreateComma(comma, result)
	}
	result.SrcrefTreeIfMissing(rest)
	return result
}

// Java method: Es6RewriteDestructuring.deletionNodeForRestProperty(com.google.javascript.rhino.Node restTempVarNameNode, com.google.javascript.rhino.Node property)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_deletionNodeForRestProperty_2638(restTempVarNameNode pkg_rhino.Node, property pkg_rhino.Node) pkg_rhino.Node {
	_ = restTempVarNameNode
	_ = property
	var get pkg_rhino.Node = pkg_rhino.Node(0)
	_ = get
	switch property.TokenSwitchOrd_() {
	case int32(100) /* Token.STRING_KEY */ :
	iifeSwitch_1:
		switch {
		default:

			if property.IsQuotedStringKey() {
				get = this_.astFactory.CreateGetElem(restTempVarNameNode, this_.astFactory.CreateString(property.GetString()))
				break iifeSwitch_1
			}
			get = this_.astFactory.CreateGetPropWithUnknownType(restTempVarNameNode, property.GetString())
			break iifeSwitch_1
		}

	case int32(32) /* Token.NAME */ :
		get = this_.astFactory.CreateGetElem(restTempVarNameNode, property)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected property to delete node: "+j2gort.StrValString(property.ToStringTree()))))
	}
	return this_.astFactory.CreateDelProp(get)
}

// Java method: Es6RewriteDestructuring.replaceArrayPattern(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node arrayPattern, com.google.javascript.rhino.Node rhs, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node nodeToDetach)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_replaceArrayPattern_2639(t *NodeTraversal, arrayPattern pkg_rhino.Node, rhs pkg_rhino.Node, parent pkg_rhino.Node, nodeToDetach pkg_rhino.Node) {
	_ = t
	_ = arrayPattern
	_ = rhs
	_ = parent
	_ = nodeToDetach
	if this_.rewriteMode == StaticGet_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST() {
		if this_.patternNestingStack.IsEmpty() || !(this_.patternNestingStack.PeekLast().hasNestedObjectRest) {
			return
		}
	}
	var tempVarName string = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
	_ = tempVarName
	var makeIteratorCall pkg_rhino.Node = this_.astFactory.CreateJSCompMakeIteratorCall(rhs.Detach(), this_.namespace)
	_ = makeIteratorCall
	var tempDecl pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_2(tempVarName, makeIteratorCall)
	_ = tempDecl
	var tempVarModel pkg_rhino.Node = tempDecl.GetFirstChild()
	_ = tempVarModel
	tempDecl.SrcrefTreeIfMissing(arrayPattern)
	tempDecl.InsertBefore(nodeToDetach)
	{
		var child pkg_rhino.Node = arrayPattern.GetFirstChild()
		_ = child
		var next pkg_rhino.Node = pkg_rhino.Node(0)
		_ = next
		for ; child != 0; child = next {
			next = child.GetNext()
			if child.IsEmpty() {
				var nextCall pkg_rhino.Node = pkg_rhino.IR_ExprResult(this_.astFactory.CreateCallWithUnknownType(this_.astFactory.CreateGetProp(tempVarModel.CloneNode_0(), "next", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT())), []pkg_rhino.Node{}))
				_ = nextCall
				nextCall.SrcrefTreeIfMissing(child)
				nextCall.InsertBefore(nodeToDetach)
				continue
			}
			var newLHS pkg_rhino.Node = pkg_rhino.Node(0)
			_ = newLHS
			var newRHS pkg_rhino.Node = pkg_rhino.Node(0)
			_ = newRHS
			if child.IsDefaultValue() {
				var nextVarName string = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
				_ = nextVarName
				var nextCallDotValue pkg_rhino.Node = this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateCallWithUnknownType(this_.astFactory.CreateGetPropWithUnknownType(tempVarModel.CloneNode_0(), "next"), []pkg_rhino.Node{}), "value")
				_ = nextCallDotValue
				var nextVarType AstFactory_Type = AstFactory_S_Type_1_Node(nextCallDotValue)
				_ = nextVarType
				var var_ pkg_rhino.Node = pkg_rhino.IR_Var_2(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(nextVarName, nextVarType), nextCallDotValue)
				_ = var_
				var_.SrcrefTreeIfMissing(child)
				var_.InsertBefore(nodeToDetach)
				newLHS = child.RemoveFirstChild()
				newRHS = this_.priv_Es6RewriteDestructuring_defaultValueHook_2646(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(nextVarName, nextVarType), child.GetLastChild().Detach())
			} else if child.IsRest() {
				newLHS = child.RemoveFirstChild()
				newRHS = this_.astFactory.CreateJscompArrayFromIteratorCall(tempVarModel.CloneNode_0(), this_.namespace)
			} else {
				newLHS = child.Detach()
				newRHS = this_.astFactory.CreateGetProp(this_.astFactory.CreateCallWithUnknownType(this_.astFactory.CreateGetPropWithUnknownType(tempVarModel.CloneNode_0(), "next"), []pkg_rhino.Node{}), "value", AstFactory_S_Type_1_Node(child))
			}
			var newNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = newNode
			if parent.IsAssign() {
				var assignment pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(newLHS, newRHS)
				_ = assignment
				newNode = pkg_rhino.IR_ExprResult(assignment)
			} else {
				newNode = pkg_rhino.IR_Declaration_3(newLHS, newRHS, parent.GetToken())
			}
			newNode.SrcrefTreeIfMissing(arrayPattern)
			newNode.InsertBefore(nodeToDetach)
			this_.Visit_3(t, newLHS, newLHS.GetParent())
		}
	}
	nodeToDetach.Detach()
	t.ReportCodeChange_0()
}

// Java method: Es6RewriteDestructuring.wrapAssignOrDestructuringInCallToArrow(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node assignment)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_wrapAssignOrDestructuringInCallToArrow_2640(t *NodeTraversal, assignment pkg_rhino.Node) {
	_ = t
	_ = assignment
	var lhs pkg_rhino.Node = assignment.GetFirstChild()
	_ = lhs
	var rhs pkg_rhino.Node = assignment.GetLastChild().Detach()
	_ = rhs
	this_.priv_Es6RewriteDestructuring_validateNoAwaitOrYieldInNestedAssignmentPattern_2641(lhs)
	var tempVarName string = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
	_ = tempVarName
	var tempVarModel pkg_rhino.Node = this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, AstFactory_S_Type_1_Node(rhs))
	_ = tempVarModel
	var paramList pkg_rhino.Node = pkg_rhino.IR_ParamList([]pkg_rhino.Node{tempVarModel.CloneNode_0()})
	_ = paramList
	var replacementExpr pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(assignment.RemoveFirstChild(), tempVarModel.CloneNode_0())
	_ = replacementExpr
	var exprResult pkg_rhino.Node = pkg_rhino.IR_ExprResult(replacementExpr)
	_ = exprResult
	var returnNode pkg_rhino.Node = pkg_rhino.IR_ReturnNode_1(tempVarModel.CloneNode_0())
	_ = returnNode
	var block pkg_rhino.Node = pkg_rhino.IR_Block_1_Node_2([]pkg_rhino.Node{exprResult, returnNode})
	_ = block
	var arrowFn pkg_rhino.Node = this_.astFactory.CreateFunction("", paramList, block, AstFactory_S_Type_2_JSTypeNative_Color(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE(), pkg_colors.StaticGet_StandardColors_UNKNOWN()))
	_ = arrowFn
	arrowFn.SetIsArrowFunction(true)
	var call pkg_rhino.Node = this_.astFactory.CreateCall(arrowFn, AstFactory_S_Type_1_Node(rhs), []pkg_rhino.Node{rhs})
	_ = call
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS(), this_.compiler)
	call.SrcrefTreeIfMissing(assignment)
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	assignment.ReplaceWith(call)
	NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(call, this_.compiler)
	this_.priv_Es6RewriteDestructuring_replacePattern_2635(t, replacementExpr.GetFirstChild(), replacementExpr.GetLastChild(), replacementExpr, exprResult)
}

// Java method: Es6RewriteDestructuring.validateNoAwaitOrYieldInNestedAssignmentPattern(com.google.javascript.rhino.Node pattern)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_validateNoAwaitOrYieldInNestedAssignmentPattern_2641(pattern pkg_rhino.Node) {
	_ = pattern
	var hasAwaitOrYieldInLhs bool = (NodeUtil_FindPreorder(pattern, lambdaSingleton_es6_rewrite_destructuring_1, StaticGet_NodeUtil_MATCH_NOT_FUNCTION()) != 0)
	_ = hasAwaitOrYieldInLhs
	if !(!(hasAwaitOrYieldInLhs)) {
		j2gort.CheckStateFmt(false, "Cannot transpile yet: destructuring assignment referencing await or yield in lhs, in nested sub-expression: %s", j2gort.BoxNode(pattern))
	}
}

// Java method: Es6RewriteDestructuring.visitDestructuringPatternInVanillaForInnerVars(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node pattern)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_visitDestructuringPatternInVanillaForInnerVars_2642(t *NodeTraversal, pattern pkg_rhino.Node) {
	_ = t
	_ = pattern
	if !(pattern.IsDestructuringPattern()) {
		j2gort.CheckArgument(false)
	}
	var destructuringLhs pkg_rhino.Node = pattern.GetParent()
	_ = destructuringLhs
	var declaration pkg_rhino.Node = destructuringLhs.GetParent()
	_ = declaration
	var insertionPoint pkg_rhino.Node = declaration.GetParent()
	_ = insertionPoint
	for insertionPoint.GetParent().IsLabel() {
		insertionPoint = insertionPoint.GetParent()
	}
	switch declaration.TokenSwitchOrd_() {
	case int32(97) /* Token.CONST */ :
		{
			var block pkg_rhino.Node = pkg_rhino.IR_Block_0().Srcref(insertionPoint)
			_ = block
			insertionPoint.ReplaceWith(block)
			block.AddChildToBack(insertionPoint)
		}
		fallthrough
	case int32(84) /* Token.VAR */ :
		{
			{
				var c pkg_rhino.Node = declaration.GetFirstChild()
				_ = c
				for ; c != destructuringLhs; c = declaration.GetFirstChild() {
					var newDeclaration pkg_rhino.Node = declaration.CloneNode_0()
					_ = newDeclaration
					newDeclaration.AddChildToBack(c.Detach())
					newDeclaration.InsertBefore(insertionPoint)
				}
			}
			var newDeclaration_2 pkg_rhino.Node = declaration.CloneNode_0()
			_ = newDeclaration_2
			newDeclaration_2.AddChildToBack(destructuringLhs.Detach())
			newDeclaration_2.InsertBefore(insertionPoint)
			this_.priv_Es6RewriteDestructuring_replacePattern_2635(t, pattern, pattern.GetNext(), newDeclaration_2, newDeclaration_2)
			if !(declaration.HasChildren()) {
				declaration.ReplaceWith(pkg_rhino.IR_Empty())
			}
		}
	case int32(111) /* Token.LET */ :
		{
			pkg_parsing.ParsingUtil_GetParamOrPatternNames(pattern, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_11 any) {
				name := j2gort.CastNode[pkg_rhino.Node](lp_11)
				_ = name
				name.CloneNode_0().InsertBefore(destructuringLhs)
			}))
			var unusedVar pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType((j2gort.StrValString(this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()) + "$unused")).Srcref(destructuringLhs)
			_ = unusedVar
			destructuringLhs.ReplaceWith(unusedVar)
			unusedVar.AddChildToBack(destructuringLhs)
			this_.priv_Es6RewriteDestructuring_wrapAssignOrDestructuringInCallToArrow_2640(t, destructuringLhs)
		}
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(declaration.String()))))
	}
	this_.compiler.ReportChangeToEnclosingScope(insertionPoint)
}

// Java method: Es6RewriteDestructuring.visitDestructuringPatternInEnhancedForInnerVars(com.google.javascript.rhino.Node pattern)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_visitDestructuringPatternInEnhancedForInnerVars_2643(pattern pkg_rhino.Node) {
	_ = pattern
	if !(pattern.IsDestructuringPattern()) {
		j2gort.CheckArgument(false)
	}
	var tempVarName string = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
	_ = tempVarName
	var destructuringLhs pkg_rhino.Node = pattern.GetParent()
	_ = destructuringLhs
	if !(destructuringLhs.IsDestructuringLhs()) {
		j2gort.CheckState(false)
	}
	var declarationNode pkg_rhino.Node = destructuringLhs.GetParent()
	_ = declarationNode
	var forNode pkg_rhino.Node = declarationNode.GetParent()
	_ = forNode
	if !(NodeUtil_IsEnhancedFor(forNode)) {
		j2gort.CheckState(false)
	}
	var block pkg_rhino.Node = forNode.GetLastChild()
	_ = block
	destructuringLhs.ReplaceWith(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, AstFactory_S_Type_1_Node(pattern)).Srcref(pattern))
	var declarationType pkg_rhino.Token = declarationNode.GetToken()
	_ = declarationType
	var decl pkg_rhino.Node = pkg_rhino.IR_Declaration_3(pattern.Detach(), this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, AstFactory_S_Type_1_Node(pattern)), declarationType)
	_ = decl
	decl.SrcrefTreeIfMissing(pattern)
	var newBlock pkg_rhino.Node = pkg_rhino.IR_Block_1_Node(decl)
	_ = newBlock
	block.ReplaceWith(newBlock)
	newBlock.AddChildToBack(block)
}

// Java method: Es6RewriteDestructuring.visitDestructuringPatternInEnhancedForWithOuterVars(com.google.javascript.rhino.Node pattern)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_visitDestructuringPatternInEnhancedForWithOuterVars_2644(pattern pkg_rhino.Node) {
	_ = pattern
	if !(pattern.IsDestructuringPattern()) {
		j2gort.CheckArgument(false)
	}
	var tempVarName string = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
	_ = tempVarName
	var forNode pkg_rhino.Node = pattern.GetParent()
	_ = forNode
	var block pkg_rhino.Node = forNode.GetLastChild()
	_ = block
	var name pkg_rhino.Node = this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, AstFactory_S_Type_1_Node(pattern))
	_ = name
	var decl pkg_rhino.Node = pkg_rhino.IR_Var_1(name)
	_ = decl
	decl.SrcrefTreeIfMissing(pattern)
	decl.InsertBefore(forNode)
	var clonedName pkg_rhino.Node = name.CloneNode_0()
	_ = clonedName
	clonedName.SrcrefTreeIfMissing(pattern)
	pattern.ReplaceWith(clonedName)
	var exprResult pkg_rhino.Node = pkg_rhino.IR_ExprResult(this_.astFactory.CreateAssign_2_Node_Node(pattern, this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, AstFactory_S_Type_1_Node(pattern))))
	_ = exprResult
	exprResult.SrcrefTreeIfMissing(pattern)
	block.AddChildToFront(exprResult)
}

// Java method: Es6RewriteDestructuring.visitDestructuringPatternInCatch(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node pattern)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_visitDestructuringPatternInCatch_2645(t *NodeTraversal, pattern pkg_rhino.Node) {
	_ = t
	_ = pattern
	var tempVarName string = this_.priv_Es6RewriteDestructuring_getTempVariableName_2633()
	_ = tempVarName
	var catchBlock pkg_rhino.Node = pattern.GetNext()
	_ = catchBlock
	var patternType AstFactory_Type = AstFactory_S_Type_1_Node(pattern)
	_ = patternType
	pattern.ReplaceWith(this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, patternType))
	catchBlock.AddChildToFront(pkg_rhino.IR_Declaration_3(pattern, this_.priv_Es6RewriteDestructuring_createTempVarNameNode_2632(tempVarName, patternType), pkg_rhino.StaticGet_Token_LET()))
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.compiler)
}

// Java method: Es6RewriteDestructuring.defaultValueHook(com.google.javascript.rhino.Node getprop, com.google.javascript.rhino.Node defaultValue)
func (this_ *Es6RewriteDestructuring) priv_Es6RewriteDestructuring_defaultValueHook_2646(getprop pkg_rhino.Node, defaultValue pkg_rhino.Node) pkg_rhino.Node {
	_ = getprop
	_ = defaultValue
	var undefined pkg_rhino.Node = this_.astFactory.CreateUndefinedValue()
	_ = undefined
	undefined.MakeNonIndexable()
	var getpropClone pkg_rhino.Node = getprop.CloneTree_0().SetColor(getprop.GetColor())
	_ = getpropClone
	return this_.astFactory.CreateHook(this_.astFactory.CreateSheq(getprop, undefined), defaultValue, getpropClone)
}

// Java enum: com.google.javascript.jscomp.Es6RewriteDestructuring.ObjectDestructuringRewriteMode
type Es6RewriteDestructuring_ObjectDestructuringRewriteMode uint8

const (
	enum_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS Es6RewriteDestructuring_ObjectDestructuringRewriteMode = 1
	enum_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST         Es6RewriteDestructuring_ObjectDestructuringRewriteMode = 2
)

var enumBaseTable_Es6RewriteDestructuring_ObjectDestructuringRewriteMode = [...]j2gort.EnumBase{
	{Name_: "REWRITE_ALL_OBJECT_PATTERNS", Ordinal_: 0},
	{Name_: "REWRITE_OBJECT_REST", Ordinal_: 1},
}
var enumValues_Es6RewriteDestructuring_ObjectDestructuringRewriteMode = []Es6RewriteDestructuring_ObjectDestructuringRewriteMode{1, 2}

func StaticGet_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS() Es6RewriteDestructuring_ObjectDestructuringRewriteMode {
	return enum_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS
}
func StaticGet_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST() Es6RewriteDestructuring_ObjectDestructuringRewriteMode {
	return enum_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST
}

func Es6RewriteDestructuring_ObjectDestructuringRewriteMode_Values() []Es6RewriteDestructuring_ObjectDestructuringRewriteMode {
	return j2gort.CloneSlice(enumValues_Es6RewriteDestructuring_ObjectDestructuringRewriteMode)
}

func Es6RewriteDestructuring_ObjectDestructuringRewriteMode_ValueOf(name string) Es6RewriteDestructuring_ObjectDestructuringRewriteMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Es6RewriteDestructuring_ObjectDestructuringRewriteMode {
		if enumBaseTable_Es6RewriteDestructuring_ObjectDestructuringRewriteMode[i_].Name_ == name {
			return Es6RewriteDestructuring_ObjectDestructuringRewriteMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.Es6RewriteDestructuring.ObjectDestructuringRewriteMode." + name))
}

func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) IsJavaNil_() bool {
	return this_ == 0
}
func (Es6RewriteDestructuring_ObjectDestructuringRewriteMode) HasIdentityEquals_() {}
func (Es6RewriteDestructuring_ObjectDestructuringRewriteMode) Implements_j2gort_Ext_java_lang_Enum() {
}
func (Es6RewriteDestructuring_ObjectDestructuringRewriteMode) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Es6RewriteDestructuring_ObjectDestructuringRewriteMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteDestructuring$ObjectDestructuringRewriteMode", "ObjectDestructuringRewriteMode"
}
func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Es6RewriteDestructuring_ObjectDestructuringRewriteMode[this_-1]
}
func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Es6RewriteDestructuring_ObjectDestructuringRewriteMode[this_-1].Name_
}
func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Es6RewriteDestructuring_ObjectDestructuringRewriteMode[this_-1].Name_
}
func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) Equals(other_ any) bool {
	o_, ok_ := other_.(Es6RewriteDestructuring_ObjectDestructuringRewriteMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Es6RewriteDestructuring_ObjectDestructuringRewriteMode[this_-1].HashCode()
}
func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Es6RewriteDestructuring_ObjectDestructuringRewriteMode[this_-1].CompareTo(other_)
}
func (this_ Es6RewriteDestructuring_ObjectDestructuringRewriteMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Es6RewriteDestructuring_ObjectDestructuringRewriteMode[this_-1].Name_)
}
func init() {
	vals_ := []any{Es6RewriteDestructuring_ObjectDestructuringRewriteMode(1), Es6RewriteDestructuring_ObjectDestructuringRewriteMode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.Es6RewriteDestructuring.ObjectDestructuringRewriteMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.Es6RewriteDestructuring$ObjectDestructuringRewriteMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteDestructuring.ObjectDestructuringRewriteMode", "com.google.javascript.jscomp.Es6RewriteDestructuring$ObjectDestructuringRewriteMode", func(v_ any) bool {
		return j2gort.InstanceOf[Es6RewriteDestructuring_ObjectDestructuringRewriteMode](v_)
	})
}

// Java class: com.google.javascript.jscomp.Es6RewriteDestructuring.Builder
type Es6RewriteDestructuring_Builder struct {
	compiler    AbstractCompiler `gson:"compiler"`
	idHash_     atomic.Int32
	rewriteMode Es6RewriteDestructuring_ObjectDestructuringRewriteMode `gson:"rewriteMode"`
}

func (this_ *Es6RewriteDestructuring_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteDestructuring_Builder) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteDestructuring_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteDestructuring$Builder", "Builder"
}
func (this_ *Es6RewriteDestructuring_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteDestructuring_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteDestructuring.Builder", "com.google.javascript.jscomp.Es6RewriteDestructuring$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteDestructuring_Builder](v_)
	})
}

func (this_ *Es6RewriteDestructuring_Builder) initFields_Es6RewriteDestructuring_Builder() {
	_ = this_
	this_.rewriteMode = StaticGet_Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS()
}

// Java constructor: Builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6RewriteDestructuring_Builder(compiler AbstractCompiler) *Es6RewriteDestructuring_Builder {
	this_ := &Es6RewriteDestructuring_Builder{}
	InitEs6RewriteDestructuring_Builder(this_, compiler)
	return this_
}

func InitEs6RewriteDestructuring_Builder(this_ *Es6RewriteDestructuring_Builder, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Es6RewriteDestructuring_Builder()
	this_.compiler = compiler
}

// Java method: Builder.setDestructuringRewriteMode(com.google.javascript.jscomp.Es6RewriteDestructuring.ObjectDestructuringRewriteMode rewriteMode)
func (this_ *Es6RewriteDestructuring_Builder) SetDestructuringRewriteMode(rewriteMode Es6RewriteDestructuring_ObjectDestructuringRewriteMode) *Es6RewriteDestructuring_Builder {
	_ = rewriteMode
	this_.rewriteMode = rewriteMode
	return this_
}

// Java method: Builder.build()
func (this_ *Es6RewriteDestructuring_Builder) Build() *Es6RewriteDestructuring {
	return NewEs6RewriteDestructuring(this_)
}

// Java class: com.google.javascript.jscomp.Es6RewriteDestructuring.PatternNestingLevel
type Es6RewriteDestructuring_PatternNestingLevel struct {
	pattern             pkg_rhino.Node `gson:"pattern"`
	idHash_             atomic.Int32
	hasNestedObjectRest bool `gson:"hasNestedObjectRest"`
}

func (this_ *Es6RewriteDestructuring_PatternNestingLevel) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteDestructuring_PatternNestingLevel) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteDestructuring_PatternNestingLevel) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteDestructuring$PatternNestingLevel", "PatternNestingLevel"
}
func (this_ *Es6RewriteDestructuring_PatternNestingLevel) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteDestructuring_PatternNestingLevel) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteDestructuring.PatternNestingLevel", "com.google.javascript.jscomp.Es6RewriteDestructuring$PatternNestingLevel", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteDestructuring_PatternNestingLevel](v_)
	})
}

func (this_ *Es6RewriteDestructuring_PatternNestingLevel) initFields_Es6RewriteDestructuring_PatternNestingLevel() {
	_ = this_
}

// Java constructor: PatternNestingLevel(com.google.javascript.rhino.Node pattern, boolean hasNestedRest)
func NewEs6RewriteDestructuring_PatternNestingLevel(pattern pkg_rhino.Node, hasNestedRest bool) *Es6RewriteDestructuring_PatternNestingLevel {
	this_ := &Es6RewriteDestructuring_PatternNestingLevel{}
	InitEs6RewriteDestructuring_PatternNestingLevel(this_, pattern, hasNestedRest)
	return this_
}

func InitEs6RewriteDestructuring_PatternNestingLevel(this_ *Es6RewriteDestructuring_PatternNestingLevel, pattern pkg_rhino.Node, hasNestedRest bool) {
	_ = pattern
	_ = hasNestedRest
	this_.initFields_Es6RewriteDestructuring_PatternNestingLevel()
	this_.pattern = pattern
	this_.hasNestedObjectRest = hasNestedRest
}
func (*Es6RewriteDestructuring) CanReleaseSyntacticScopes_() {}

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
// Port of com/google/javascript/jscomp/ExpressionDecomposer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ExpressionDecomposer
type ExpressionDecomposer struct {
	compiler                       AbstractCompiler                                  `gson:"compiler"`
	astAnalyzer                    *AstAnalyzer                                      `gson:"astAnalyzer"`
	astFactory                     *AstFactory                                       `gson:"astFactory"`
	safeNameIdSupplier             j2gort.Ext_java_util_function_Supplier            `gson:"safeNameIdSupplier"`
	knownConstantFunctions         j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"knownConstantFunctions"`
	knownConstantFunctionsSet      *j2gort.OrderedSet[string]                        `gson:"knownConstantFunctionsSet"`
	scope                          *Scope                                            `gson:"scope"`
	unknownType                    pkg_rhino.JSType                                  `gson:"unknownType"`
	tempNamePrefix                 string                                            `gson:"tempNamePrefix"`
	tempConstPrefix                string                                            `gson:"tempConstPrefix"`
	resultNamePrefix               string                                            `gson:"resultNamePrefix"`
	idHash_                        atomic.Int32
	enableBrokenIe11LocationAssign bool `gson:"enableBrokenIe11LocationAssign"`
}

func (this_ *ExpressionDecomposer) BindSelf_(s any) {
	_ = s
}
func (this_ *ExpressionDecomposer) IsJavaNil_() bool { return this_ == nil }
func (*ExpressionDecomposer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExpressionDecomposer", "ExpressionDecomposer"
}
func (this_ *ExpressionDecomposer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ExpressionDecomposer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExpressionDecomposer", "com.google.javascript.jscomp.ExpressionDecomposer", func(v_ any) bool {
		return j2gort.InstanceOf[*ExpressionDecomposer](v_)
	})
}

func StaticGet_ExpressionDecomposer_MAX_ITERATIONS() int32 { return int32(1000) }

var static_ExpressionDecomposer_WINDOW_LOCATION_ASSIGN pkg_rhino.QualifiedName

func StaticGet_ExpressionDecomposer_WINDOW_LOCATION_ASSIGN() pkg_rhino.QualifiedName {
	clinit_ExpressionDecomposer()
	return static_ExpressionDecomposer_WINDOW_LOCATION_ASSIGN
}
func StaticSet_ExpressionDecomposer_WINDOW_LOCATION_ASSIGN(v pkg_rhino.QualifiedName) {
	clinit_ExpressionDecomposer()
	static_ExpressionDecomposer_WINDOW_LOCATION_ASSIGN = v
}

var clinitDone_ExpressionDecomposer uint32
var clinitState_ExpressionDecomposer j2gort.ClinitState

func clinit_ExpressionDecomposer() {
	if atomic.LoadUint32(&clinitDone_ExpressionDecomposer) != 0 {
		return
	}
	clinitSlow_ExpressionDecomposer()
}

func clinitSlow_ExpressionDecomposer() {
	if !j2gort.EnterClinit(&clinitState_ExpressionDecomposer, &clinitDone_ExpressionDecomposer) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ExpressionDecomposer, &clinitDone_ExpressionDecomposer)
	static_ExpressionDecomposer_WINDOW_LOCATION_ASSIGN = pkg_rhino.QualifiedName_Of("window.location.assign")
}

func (this_ *ExpressionDecomposer) initFields_ExpressionDecomposer() {
	_ = this_
	this_.tempNamePrefix = "JSCompiler_temp"
	this_.tempConstPrefix = ("JSCompiler_temp_const" + "$jscomp$")
	this_.resultNamePrefix = "JSCompiler_inline_result"
}

// Java constructor: ExpressionDecomposer(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.function.Supplier<java.lang.String> safeNameIdSupplier, com.google.common.collect.ImmutableSet<java.lang.String> constFunctionNames, com.google.javascript.jscomp.Scope scope, java.util.EnumSet<com.google.javascript.jscomp.ExpressionDecomposer.Workaround> enabledWorkarounds)
func NewExpressionDecomposer_5_AbstractCompiler_Ext_java_util_function_Supplier_Ext_com_google_common_collect_ImmutableSet_Scope_Ext_java_util_EnumSet(compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier, constFunctionNames j2gort.Ext_com_google_common_collect_ImmutableSet, scope *Scope, enabledWorkarounds j2gort.Ext_java_util_EnumSet) *ExpressionDecomposer {
	this_ := &ExpressionDecomposer{}
	InitExpressionDecomposer_5_AbstractCompiler_Ext_java_util_function_Supplier_Ext_com_google_common_collect_ImmutableSet_Scope_Ext_java_util_EnumSet(this_, compiler, safeNameIdSupplier, constFunctionNames, scope, enabledWorkarounds)
	return this_
}

func InitExpressionDecomposer_5_AbstractCompiler_Ext_java_util_function_Supplier_Ext_com_google_common_collect_ImmutableSet_Scope_Ext_java_util_EnumSet(this_ *ExpressionDecomposer, compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier, constFunctionNames j2gort.Ext_com_google_common_collect_ImmutableSet, scope *Scope, enabledWorkarounds j2gort.Ext_java_util_EnumSet) {
	_ = compiler
	_ = safeNameIdSupplier
	_ = constFunctionNames
	_ = scope
	_ = enabledWorkarounds
	InitExpressionDecomposer_5_AbstractCompiler_Ext_java_util_function_Supplier_Ext_com_google_common_collect_ImmutableSet_Scope_bool(this_, compiler, safeNameIdSupplier, constFunctionNames, scope, j2gort.SetContainsTypedU8_(enabledWorkarounds, StaticGet_ExpressionDecomposer_Workaround_BROKEN_IE11_LOCATION_ASSIGN()))
}

// Java constructor: ExpressionDecomposer(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.function.Supplier<java.lang.String> safeNameIdSupplier, com.google.common.collect.ImmutableSet<java.lang.String> constFunctionNames, com.google.javascript.jscomp.Scope scope, boolean enableBrokenIe11LocationAssign)
func NewExpressionDecomposer_5_AbstractCompiler_Ext_java_util_function_Supplier_Ext_com_google_common_collect_ImmutableSet_Scope_bool(compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier, constFunctionNames j2gort.Ext_com_google_common_collect_ImmutableSet, scope *Scope, enableBrokenIe11LocationAssign bool) *ExpressionDecomposer {
	this_ := &ExpressionDecomposer{}
	InitExpressionDecomposer_5_AbstractCompiler_Ext_java_util_function_Supplier_Ext_com_google_common_collect_ImmutableSet_Scope_bool(this_, compiler, safeNameIdSupplier, constFunctionNames, scope, enableBrokenIe11LocationAssign)
	return this_
}

func InitExpressionDecomposer_5_AbstractCompiler_Ext_java_util_function_Supplier_Ext_com_google_common_collect_ImmutableSet_Scope_bool(this_ *ExpressionDecomposer, compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier, constFunctionNames j2gort.Ext_com_google_common_collect_ImmutableSet, scope *Scope, enableBrokenIe11LocationAssign bool) {
	_ = compiler
	_ = safeNameIdSupplier
	_ = constFunctionNames
	_ = scope
	_ = enableBrokenIe11LocationAssign
	this_.initFields_ExpressionDecomposer()
	j2gort.CheckNotNullIface(compiler)
	j2gort.CheckNotNull(safeNameIdSupplier)
	j2gort.CheckNotNull(constFunctionNames)
	this_.compiler = compiler
	this_.astAnalyzer = compiler.GetAstAnalyzer()
	this_.astFactory = compiler.CreateAstFactory()
	this_.safeNameIdSupplier = safeNameIdSupplier
	this_.knownConstantFunctions = constFunctionNames
	if constFunctionNames.IsEmpty() {
		this_.knownConstantFunctionsSet = nil
	} else {
		this_.knownConstantFunctionsSet = j2gort.OrderedSetPtrFromStrCol(j2gort.AsIface[j2gort.Ext_java_util_Collection](constFunctionNames))
	}
	this_.scope = scope
	this_.enableBrokenIe11LocationAssign = enableBrokenIe11LocationAssign
	if compiler.HasTypeCheckingRun() && !(compiler.HasOptimizationColors()) {
		this_.unknownType = compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	} else {
		this_.unknownType = nil
	}
}

// Java method: ExpressionDecomposer.reset(com.google.javascript.jscomp.Scope scope, com.google.common.collect.ImmutableSet<java.lang.String> constFunctionNames)
func (this_ *ExpressionDecomposer) Reset(scope *Scope, constFunctionNames j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = scope
	_ = constFunctionNames
	this_.scope = scope
	if !(j2gort.RefEquals(this_.knownConstantFunctions, constFunctionNames)) {
		j2gort.CheckNotNull(constFunctionNames)
		this_.knownConstantFunctions = constFunctionNames
		if constFunctionNames.IsEmpty() {
			this_.knownConstantFunctionsSet = nil
		} else {
			this_.knownConstantFunctionsSet = j2gort.OrderedSetPtrFromStrCol(j2gort.AsIface[j2gort.Ext_java_util_Collection](constFunctionNames))
		}
	}
}

// Java method: ExpressionDecomposer.maybeExposeExpression(com.google.javascript.rhino.Node expression)
func (this_ *ExpressionDecomposer) MaybeExposeExpression(expression pkg_rhino.Node) {
	_ = expression
	var i int32 = int32(0)
	_ = i
	for StaticGet_ExpressionDecomposer_DecompositionType_DECOMPOSABLE() == this_.CanExposeExpression(expression) {
		if !(this_.priv_ExpressionDecomposer_exposeExpression_2817(expression)) {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "exposeExpression exposed nothing for:\n"+j2gort.StrValString(expression.ToStringTree()))))
		}
		i = i + 1
		if i > int32(1000) {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "exposeExpression depth exceeded on:\n"+j2gort.StrValString(expression.ToStringTree()))))
		}
	}
}

// Java method: ExpressionDecomposer.exposeExpression(com.google.javascript.rhino.Node expression)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_exposeExpression_2817(expression pkg_rhino.Node) bool {
	_ = expression
	this_.priv_ExpressionDecomposer_rewriteAllContainingOptionalChains_2819(expression)
	var expressionRoot pkg_rhino.Node = ExpressionDecomposer_FindExpressionRoot(expression)
	_ = expressionRoot
	j2gort.CheckNotNullNode(expressionRoot)
	if !(NodeUtil_IsStatement(expressionRoot)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(expressionRoot))
	}
	return this_.priv_ExpressionDecomposer_exposeExpression_2818(expressionRoot, expression)
}

// Java method: ExpressionDecomposer.exposeExpression(com.google.javascript.rhino.Node expressionRoot, com.google.javascript.rhino.Node subExpression)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_exposeExpression_2818(expressionRoot pkg_rhino.Node, subExpression pkg_rhino.Node) bool {
	_ = expressionRoot
	_ = subExpression
	var exposedSomething bool = false
	_ = exposedSomething
	var nodeWithNonconditionalParent pkg_rhino.Node = ExpressionDecomposer_FindNonconditionalParent(subExpression, expressionRoot)
	_ = nodeWithNonconditionalParent
	var hasFollowingSideEffects bool = this_.astAnalyzer.MayHaveSideEffects(nodeWithNonconditionalParent)
	_ = hasFollowingSideEffects
	var exprInjectionPoint pkg_rhino.Node = ExpressionDecomposer_FindInjectionPoint(nodeWithNonconditionalParent)
	_ = exprInjectionPoint
	var val_1 ExpressionDecomposer_DecompositionState
	InitExpressionDecomposer_DecompositionState(&val_1)
	var state *ExpressionDecomposer_DecompositionState = &val_1
	_ = state
	{
		r_2 := state
		r_2.sideEffects = hasFollowingSideEffects
	}
	{
		r_3 := state
		r_3.extractBeforeStatement = exprInjectionPoint
	}
	var lastExposedSubexpression pkg_rhino.Node = pkg_rhino.Node(0)
	_ = lastExposedSubexpression
	var expressionToExpose pkg_rhino.Node = nodeWithNonconditionalParent
	_ = expressionToExpose
	var expressionParent pkg_rhino.Node = expressionToExpose.GetParent()
	_ = expressionParent
	for expressionParent != expressionRoot {
		if !(!(ExpressionDecomposer_IsConditionalOp(expressionParent)) || expressionToExpose.IsFirstChildOf(expressionParent)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(expressionParent))
		}
		if expressionParent.IsAssign() {
			if this_.priv_ExpressionDecomposer_isSafeAssign_2836(expressionParent, state.sideEffects) {
			} else if !(expressionToExpose.IsFirstChildOf(expressionParent)) {
				var left pkg_rhino.Node = expressionParent.GetFirstChild()
				_ = left
				switch left.TokenSwitchOrd_() {
				case int32(102) /* Token.ARRAY_PATTERN */, int32(103) /* Token.OBJECT_PATTERN */ :
				case int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */ :
					exposedSomething = (this_.priv_ExpressionDecomposer_decomposeSubExpressions_2821(left.GetFirstChild(), pkg_rhino.Node(0), state) || exposedSomething)
				default:
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Expected a property access or destructuring pattern: "+j2gort.StrValString(left.ToStringTree()))))
				}
			}
		} else if expressionParent.IsCall() && NodeUtil_IsNormalGet(expressionParent.GetFirstChild()) {
			var callee pkg_rhino.Node = expressionParent.GetFirstChild()
			_ = callee
			if callee != expressionToExpose {
				exposedSomething = (this_.priv_ExpressionDecomposer_decomposeSubExpressions_2821(callee.GetNext(), expressionToExpose, state) || exposedSomething)
			}
			if this_.priv_ExpressionDecomposer_isExpressionTreeUnsafe_2837(callee, state.sideEffects) && (lastExposedSubexpression != callee.GetFirstChild()) {
				{
					r_4 := state
					r_4.sideEffects = true
				}
				this_.priv_ExpressionDecomposer_rewriteCallExpression_2825(expressionParent, state)
				exposedSomething = true
			}
		} else {
			exposedSomething = (this_.priv_ExpressionDecomposer_decomposeSubExpressions_2821(expressionParent.GetFirstChild(), expressionToExpose, state) || exposedSomething)
		}
		lastExposedSubexpression = expressionToExpose
		expressionToExpose = expressionParent
		expressionParent = expressionToExpose.GetParent()
	}
	if nodeWithNonconditionalParent == subExpression {
	} else {
		var parent pkg_rhino.Node = nodeWithNonconditionalParent.GetParent()
		_ = parent
		var needResult bool = !(parent.IsExprResult())
		_ = needResult
		this_.priv_ExpressionDecomposer_extractConditional_2822(nodeWithNonconditionalParent, exprInjectionPoint, needResult)
		exposedSomething = true
	}
	return exposedSomething
}

// Java method: ExpressionDecomposer.rewriteAllContainingOptionalChains(com.google.javascript.rhino.Node subExpression)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_rewriteAllContainingOptionalChains_2819(subExpression pkg_rhino.Node) {
	_ = subExpression
	var optChainRewriterBuilder *OptionalChainRewriter_Builder = OptionalChainRewriter_S_Builder(this_.compiler).SetTmpVarNameCreator(func() OptionalChainRewriter_TmpVarNameCreator {
		bound_5 := this_
		return NewLambda_OptionalChainRewriter_TmpVarNameCreator(func() string { return bound_5.priv_ExpressionDecomposer_getTempConstantValueName_2830() })
	}()).SetScope(this_.scope)
	_ = optChainRewriterBuilder
	var rewriters j2gort.Deque[*OptionalChainRewriter]
	_ = rewriters
	{
		var exprParent pkg_rhino.Node = subExpression.GetParent()
		_ = exprParent
		for ; !(NodeUtil_IsStatement(exprParent)); exprParent = exprParent.GetParent() {
			if NodeUtil_IsEndOfFullOptChain(exprParent) {
				rewriters.AddFirst(optChainRewriterBuilder.Build(exprParent))
			} else if exprParent.IsCall() {
				var callee pkg_rhino.Node = exprParent.GetFirstChild()
				_ = callee
				if NodeUtil_IsOptChainGet(callee) {
					rewriters.AddFirst(optChainRewriterBuilder.Build(callee))
				}
			}
		}
	}
	{
		cur_6 := rewriters.Iter()
		for cur_6.HasNext() {

			rewriter := cur_6.Next()
			_ = rewriter
			rewriter.Rewrite()
		}
	}
}

// Java method: ExpressionDecomposer.moveExpression(com.google.javascript.rhino.Node expression)
func (this_ *ExpressionDecomposer) MoveExpression(expression pkg_rhino.Node) {
	_ = expression
	var resultName string = this_.priv_ExpressionDecomposer_getResultValueName_2829()
	_ = resultName
	var injectionPoint pkg_rhino.Node = ExpressionDecomposer_FindInjectionPoint(expression)
	_ = injectionPoint
	j2gort.CheckNotNullNode(injectionPoint)
	var injectionPointParent pkg_rhino.Node = injectionPoint.GetParent()
	_ = injectionPointParent
	j2gort.CheckNotNullNode(injectionPointParent)
	if !(NodeUtil_IsStatementBlock(injectionPointParent)) {
		j2gort.CheckState(false)
	}
	expression.ReplaceWith(pkg_rhino.IR_Name(resultName).CopyTypeFrom(expression))
	var newExpressionRoot pkg_rhino.Node = NodeUtil_NewVarNode_2_string_Node(resultName, expression)
	_ = newExpressionRoot
	newExpressionRoot.GetFirstChild().CopyTypeFrom(expression)
	newExpressionRoot.InsertBefore(injectionPoint)
	this_.compiler.ReportChangeToEnclosingScope(injectionPointParent)
}

// Java static method: ExpressionDecomposer.findNonconditionalParent(com.google.javascript.rhino.Node subExpression, com.google.javascript.rhino.Node expressionRoot)
func ExpressionDecomposer_FindNonconditionalParent(subExpression pkg_rhino.Node, expressionRoot pkg_rhino.Node) pkg_rhino.Node {
	_ = subExpression
	_ = expressionRoot
	var result pkg_rhino.Node = subExpression
	_ = result
	{
		var child pkg_rhino.Node = subExpression
		_ = child
		var parent pkg_rhino.Node = child.GetParent()
		_ = parent
		for ; parent != expressionRoot; func() { child = parent; parent = child.GetParent() }() {
			if ExpressionDecomposer_IsConditionalOp(parent) && !(child.IsFirstChildOf(parent)) {
				result = parent
			}
		}
	}
	if NodeUtil_IsOptChainNode(result) {
		result = NodeUtil_GetEndOfOptChainSegment(result)
	}
	return result
}

// Java method: ExpressionDecomposer.decomposeSubExpressions(com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node stopNode, com.google.javascript.jscomp.ExpressionDecomposer.DecompositionState state)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_decomposeSubExpressions_2821(n pkg_rhino.Node, stopNode pkg_rhino.Node, state *ExpressionDecomposer_DecompositionState) bool {
	_ = n
	_ = stopNode
	_ = state
	if (n == 0) || (n == stopNode) {
		return false
	}
	var decomposedSomething bool = false
	_ = decomposedSomething
	decomposedSomething = (this_.priv_ExpressionDecomposer_decomposeSubExpressions_2821(n.GetNext(), stopNode, state) || decomposedSomething)
	if NodeUtil_MayBeObjectLitKey(n) || n.IsComputedProp() {
		if n.IsComputedProp() {
			decomposedSomething = (this_.priv_ExpressionDecomposer_decomposeSubExpressions_2821(n.GetSecondChild(), stopNode, state) || decomposedSomething)
		}
		n = n.GetFirstChild()
	} else if n.IsTemplateLitSub() {
		n = n.GetFirstChild()
	} else if n.IsSpread() {
	} else if !(pkg_rhino.IR_MayBeExpression(n)) {
		return decomposedSomething
	}
	if this_.priv_ExpressionDecomposer_isExpressionTreeUnsafe_2837(n, state.sideEffects) {
		{
			r_8 := state
			r_8.sideEffects = true
		}
		{
			r_9 := state
			r_9.extractBeforeStatement = this_.priv_ExpressionDecomposer_extractExpression_2824(n, state.extractBeforeStatement)
		}
		decomposedSomething = true
	}
	return decomposedSomething
}

// Java method: ExpressionDecomposer.extractConditional(com.google.javascript.rhino.Node expr, com.google.javascript.rhino.Node injectionPoint, boolean needResult)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_extractConditional_2822(expr pkg_rhino.Node, injectionPoint pkg_rhino.Node, needResult bool) pkg_rhino.Node {
	_ = expr
	_ = injectionPoint
	_ = needResult
	var parent pkg_rhino.Node = expr.GetParent()
	_ = parent
	var tempName string = this_.priv_ExpressionDecomposer_getTempValueName_2827()
	_ = tempName
	var first pkg_rhino.Node = expr.GetFirstChild()
	_ = first
	var second pkg_rhino.Node = first.GetNext()
	_ = second
	var last pkg_rhino.Node = expr.GetLastChild()
	_ = last
	expr.DetachChildren()
	var cond pkg_rhino.Node = pkg_rhino.Node(0)
	_ = cond
	var trueExpr pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{}).Srcref(expr)
	_ = trueExpr
	var falseExpr pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{}).Srcref(expr)
	_ = falseExpr
	switch expr.TokenSwitchOrd_() {
	case int32(67) /* Token.HOOK */ :
		cond = first
		trueExpr.AddChildToFront(this_.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(second, needResult, tempName)))
		falseExpr.AddChildToFront(this_.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(last, needResult, tempName)))
	case int32(69) /* Token.AND */ :
		cond = ExpressionDecomposer_BuildResultExpression(first, needResult, tempName)
		trueExpr.AddChildToFront(this_.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(last, needResult, tempName)))
	case int32(68) /* Token.OR */ :
		cond = ExpressionDecomposer_BuildResultExpression(first, needResult, tempName)
		falseExpr.AddChildToFront(this_.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(last, needResult, tempName)))
	case int32(70) /* Token.COALESCE */ :
		var tempNameAssign string = this_.priv_ExpressionDecomposer_getTempValueName_2827()
		_ = tempNameAssign
		var tempVarNodeAssign pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_1(tempNameAssign).SrcrefTreeIfMissing(expr)
		_ = tempVarNodeAssign
		tempVarNodeAssign.InsertBefore(injectionPoint)
		var assignLhs pkg_rhino.Node = ExpressionDecomposer_BuildResultExpression(first, true, tempNameAssign)
		_ = assignLhs
		var nullNode pkg_rhino.Node = this_.astFactory.CreateNull().Srcref(expr)
		_ = nullNode
		cond = this_.astFactory.CreateNe(assignLhs, nullNode).Srcref(expr)
		trueExpr.AddChildToFront(this_.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(this_.astFactory.CreateName_2_string_AstFactory_Type(tempNameAssign, AstFactory_S_Type_1_Node(first)).Srcref(expr), needResult, tempName)))
		falseExpr.AddChildToFront(this_.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(last, needResult, tempName)))
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected expression: "+j2gort.StrVal(expr))))
	}
	var ifNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = ifNode
	if falseExpr.HasChildren() {
		ifNode = this_.astFactory.CreateIf_3(cond, trueExpr, falseExpr)
	} else {
		ifNode = this_.astFactory.CreateIf_2(cond, trueExpr)
	}
	ifNode.SrcrefIfMissing(expr)
	if needResult {
		var tempVarNode pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_1(tempName).SrcrefTreeIfMissing(expr)
		_ = tempVarNode
		tempVarNode.InsertBefore(injectionPoint)
		ifNode.InsertAfter(tempVarNode)
		var replacementValueNode pkg_rhino.Node = pkg_rhino.IR_Name(tempName).CopyTypeFrom(expr)
		_ = replacementValueNode
		expr.ReplaceWith(replacementValueNode)
	} else {
		if !(parent.IsExprResult()) {
			j2gort.CheckArgument(false)
		}
		parent.ReplaceWith(ifNode)
	}
	return ifNode
}

// Java static method: ExpressionDecomposer.buildResultExpression(com.google.javascript.rhino.Node expr, boolean needResult, java.lang.String tempName)
func ExpressionDecomposer_BuildResultExpression(expr pkg_rhino.Node, needResult bool, tempName string) pkg_rhino.Node {
	_ = expr
	_ = needResult
	_ = tempName
	if needResult {
		return pkg_rhino.IR_Assign(pkg_rhino.IR_Name(tempName).CopyTypeFrom(expr), expr).CopyTypeFrom(expr).SrcrefTree(expr)
	} else {
		return expr
	}
}

// Java method: ExpressionDecomposer.isConstantNameNode(com.google.javascript.rhino.Node n)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_isConstantNameNode_2823(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsName() && (NodeUtil_IsConstantVar(n, this_.scope) || ((this_.knownConstantFunctionsSet != nil) && this_.knownConstantFunctionsSet.Contains(n.GetString()))))
}

// Java method: ExpressionDecomposer.extractExpression(com.google.javascript.rhino.Node expr, com.google.javascript.rhino.Node injectionPoint)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_extractExpression_2824(expr pkg_rhino.Node, injectionPoint pkg_rhino.Node) pkg_rhino.Node {
	_ = expr
	_ = injectionPoint
	if !(!(NodeUtil_IsLogicalAssignmentOp(expr))) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(expr))
	}
	var parent pkg_rhino.Node = expr.GetParent()
	_ = parent
	var isLhsOfAssignOp bool = ((NodeUtil_IsAssignmentOp(parent) && !(parent.IsAssign())) && expr.IsFirstChildOf(parent))
	_ = isLhsOfAssignOp
	var firstExtractedNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = firstExtractedNode
	if isLhsOfAssignOp && NodeUtil_IsNormalGet(expr) {
		{
			var n pkg_rhino.Node = expr.GetFirstChild()
			_ = n
			var next pkg_rhino.Node = pkg_rhino.Node(0)
			_ = next
			for ; n != 0; n = next {
				next = n.GetNext()
				if !(n.IsStringLit()) && !(this_.priv_ExpressionDecomposer_isConstantNameNode_2823(n)) {
					var extractedNode pkg_rhino.Node = this_.priv_ExpressionDecomposer_extractExpression_2824(n, injectionPoint)
					_ = extractedNode
					if firstExtractedNode == 0 {
						firstExtractedNode = extractedNode
					}
				}
			}
		}
	}
	var tempName string = this_.priv_ExpressionDecomposer_getTempConstantValueName_2830()
	_ = tempName
	var replacementValueNode pkg_rhino.Node = pkg_rhino.IR_Name(tempName).CopyTypeFrom(expr).Srcref(expr)
	_ = replacementValueNode
	var tempNameValue pkg_rhino.Node = pkg_rhino.Node(0)
	_ = tempNameValue
	if isLhsOfAssignOp {
		if !(expr.IsName() || NodeUtil_IsNormalGet(expr)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(expr))
		}
		var opNode pkg_rhino.Node = pkg_rhino.NewNode_1(NodeUtil_GetOpFromAssignmentOp(parent)).CopyTypeFrom(parent).SrcrefIfMissing(parent)
		_ = opNode
		var rightOperand pkg_rhino.Node = parent.GetLastChild()
		_ = rightOperand
		parent.SetToken(pkg_rhino.StaticGet_Token_ASSIGN())
		rightOperand.ReplaceWith(opNode)
		opNode.AddChildToFront(replacementValueNode)
		opNode.AddChildToBack(rightOperand)
		tempNameValue = expr.CloneTree_0()
	} else if expr.IsSpread() {
		var spreadCopy pkg_rhino.Node = expr.CloneNode_0()
		_ = spreadCopy
		spreadCopy.AddChildToBack(replacementValueNode)
		expr.ReplaceWith(spreadCopy)
	iifeSwitch_1:
		switch {
		default:

			switch parent.TokenSwitchOrd_() {
			case int32(46) /* Token.ARRAYLIT */, int32(28) /* Token.CALL */, int32(23) /* Token.NEW */ :
				tempNameValue = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{expr}).Srcref(expr.GetOnlyChild())
				break iifeSwitch_1
			case int32(47) /* Token.OBJECTLIT */ :
				tempNameValue = this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{expr}).Srcref(expr.GetOnlyChild())
				break iifeSwitch_1
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected parent of SPREAD:"+j2gort.StrValString(parent.ToStringTree()))))
			}
			panic("unreachable switch expression")
		}

	} else {
		expr.ReplaceWith(replacementValueNode)
		tempNameValue = expr
	}
	var tempVarNode pkg_rhino.Node = NodeUtil_NewVarNode_2_string_Node(tempName, tempNameValue)
	_ = tempVarNode
	tempVarNode.GetFirstChild().CopyTypeFrom(tempNameValue)
	tempVarNode.GetFirstChild().SetInferredConstantVar(true)
	var containingHoistScope *Scope = j2gort.CastPtr[Scope](this_.scope.GetClosestHoistScope())
	_ = containingHoistScope
	containingHoistScope.Declare(tempName, tempVarNode.GetFirstChild(), nil)
	tempVarNode.InsertBefore(injectionPoint)
	if firstExtractedNode == 0 {
		firstExtractedNode = tempVarNode
	}
	if !(firstExtractedNode.IsVar()) {
		j2gort.CheckState(false)
	}
	return firstExtractedNode
}

// Java method: ExpressionDecomposer.rewriteCallExpression(com.google.javascript.rhino.Node call, com.google.javascript.jscomp.ExpressionDecomposer.DecompositionState state)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_rewriteCallExpression_2825(call pkg_rhino.Node, state *ExpressionDecomposer_DecompositionState) {
	_ = call
	_ = state
	if !(call.IsCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(call))
	}
	var first pkg_rhino.Node = call.GetFirstChild()
	_ = first
	if !(NodeUtil_IsNormalGet(first)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(first))
	}
	var fnCallType pkg_rhino.JSType = nil
	_ = fnCallType
	if this_.astFactory.IsAddingTypes() {
		var fnType pkg_rhino.JSType = first.GetJSType()
		_ = fnType
		if fnType.GetJSTypeBase().IsFunctionType() {
			fnCallType = fnType.ToMaybeFunctionType().GetPropertyType_1_string("call")
		} else {
			fnCallType = this_.unknownType
		}
	}
	var getVarNode pkg_rhino.Node = this_.priv_ExpressionDecomposer_extractExpression_2824(first, state.extractBeforeStatement)
	_ = getVarNode
	{
		r_10 := state
		r_10.extractBeforeStatement = getVarNode
	}
	var functionNameNode pkg_rhino.Node = getVarNode.GetFirstChild().CloneNode_0()
	_ = functionNameNode
	if call.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL()) {
		call.GetFirstChild().ReplaceWith(functionNameNode)
		return
	}
	var getExprNode pkg_rhino.Node = getVarNode.GetFirstFirstChild()
	_ = getExprNode
	if !(NodeUtil_IsNormalGet(getExprNode)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(getExprNode))
	}
	var origThisValue pkg_rhino.Node = getExprNode.GetFirstChild()
	_ = origThisValue
	var receiverNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = receiverNode
	if origThisValue.IsThis() {
		receiverNode = origThisValue.CloneNode_0()
	} else if origThisValue.IsSuper() {
		receiverNode = this_.astFactory.CreateThis(AstFactory_S_Type_1_Node(origThisValue)).Srcref(origThisValue)
	} else {
		var thisVarNode pkg_rhino.Node = this_.priv_ExpressionDecomposer_extractExpression_2824(origThisValue, state.extractBeforeStatement)
		_ = thisVarNode
		{
			r_11 := state
			r_11.extractBeforeStatement = thisVarNode
		}
		receiverNode = thisVarNode.GetFirstChild().CloneNode_0()
	}
	call.RemoveFirstChild()
	call.AddChildToFront(receiverNode)
	call.AddChildToFront(this_.astFactory.CreateGetProp(functionNameNode, "call", AstFactory_S_Type_2_JSType_Color(fnCallType, pkg_colors.StaticGet_StandardColors_TOP_OBJECT())).SrcrefTreeIfMissing(call))
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), false)
}

// Java method: ExpressionDecomposer.setTempNamePrefix(java.lang.String prefix)
func (this_ *ExpressionDecomposer) SetTempNamePrefix(prefix string) {
	_ = prefix
	this_.tempNamePrefix = prefix
	this_.tempConstPrefix = (j2gort.StrValString(prefix) + "_const" + "$jscomp$")
}

// Java method: ExpressionDecomposer.getTempValueName()
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_getTempValueName_2827() string {
	return (j2gort.StrValString(this_.tempNamePrefix) + "$jscomp$" + j2gort.StrValString(j2gort.CastString(this_.safeNameIdSupplier.Get())))
}

// Java method: ExpressionDecomposer.setResultNamePrefix(java.lang.String prefix)
func (this_ *ExpressionDecomposer) SetResultNamePrefix(prefix string) {
	_ = prefix
	this_.resultNamePrefix = prefix
}

// Java method: ExpressionDecomposer.getResultValueName()
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_getResultValueName_2829() string {
	return (j2gort.StrValString(this_.resultNamePrefix) + "$jscomp$" + j2gort.StrValString(j2gort.CastString(this_.safeNameIdSupplier.Get())))
}

// Java method: ExpressionDecomposer.getTempConstantValueName()
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_getTempConstantValueName_2830() string {
	return (j2gort.StrValString(this_.tempConstPrefix) + j2gort.StrValString(j2gort.CastString(this_.safeNameIdSupplier.Get())))
}

// Java method: ExpressionDecomposer.isTempConstantValueName(com.google.javascript.rhino.Node name)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_isTempConstantValueName_2831(name pkg_rhino.Node) bool {
	_ = name
	return (name.IsName() && javalang.StringStartsWith(j2gort.CheckNotNullString(name.GetString()), j2gort.CheckNotNullString(this_.tempConstPrefix)))
}

// Java static method: ExpressionDecomposer.findInjectionPoint(com.google.javascript.rhino.Node subExpression)
func ExpressionDecomposer_FindInjectionPoint(subExpression pkg_rhino.Node) pkg_rhino.Node {
	_ = subExpression
	var expressionRoot pkg_rhino.Node = ExpressionDecomposer_FindExpressionRoot(subExpression)
	_ = expressionRoot
	j2gort.CheckNotNullNode(expressionRoot)
	var injectionPoint pkg_rhino.Node = expressionRoot
	_ = injectionPoint
	var parent pkg_rhino.Node = injectionPoint.GetParent()
	_ = parent
	for parent.IsLabel() {
		injectionPoint = parent
		parent = injectionPoint.GetParent()
	}
	if !(NodeUtil_IsStatementBlock(parent)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
	}
	return injectionPoint
}

// Java static method: ExpressionDecomposer.isConditionalOp(com.google.javascript.rhino.Node n)
func ExpressionDecomposer_IsConditionalOp(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(67) /* Token.HOOK */, int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(31) /* Token.OPTCHAIN_CALL */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: ExpressionDecomposer.findExpressionRoot(com.google.javascript.rhino.Node subExpression)
func ExpressionDecomposer_FindExpressionRoot(subExpression pkg_rhino.Node) pkg_rhino.Node {
	_ = subExpression
	var child pkg_rhino.Node = subExpression
	_ = child
	{
		curAnc_ := child
		if curAnc_ == 0 {
			j2gort.ThrowNPE()
		}
		curAnc_ = curAnc_.GetParent()
		for curAnc_ != 0 {
			current := curAnc_
			curAnc_ = curAnc_.GetParent()

			_ = current
			var parent pkg_rhino.Node = current.GetParent()
			_ = parent
			switch current.TokenSwitchOrd_() {
			case int32(93) /* Token.EXPR_RESULT */, int32(74) /* Token.IF */, int32(75) /* Token.SWITCH */, int32(0) /* Token.RETURN */, int32(43) /* Token.THROW */ :
				if !(child.IsFirstChildOf(current)) {
					j2gort.CheckState(false)
				}
				return current
			case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
				if NodeUtil_IsAnyFor(parent) {
					break
				}
				return current
			case int32(80) /* Token.FOR */ :
				if child.IsFirstChildOf(current) {
					return current
				}
				fallthrough
			case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(79) /* Token.DO */, int32(78) /* Token.WHILE */, int32(94) /* Token.SCRIPT */, int32(90) /* Token.BLOCK */, int32(92) /* Token.LABEL */, int32(76) /* Token.CASE */, int32(91) /* Token.SWITCH_BODY */, int32(77) /* Token.DEFAULT_CASE */, int32(134) /* Token.DEFAULT_VALUE */, int32(49) /* Token.PARAM_LIST */, int32(105) /* Token.CLASS */ :
				return pkg_rhino.Node(0)
			default:
			}
			child = current
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected AST structure.")))
}

// Java method: ExpressionDecomposer.canExposeExpression(com.google.javascript.rhino.Node subExpression)
func (this_ *ExpressionDecomposer) CanExposeExpression(subExpression pkg_rhino.Node) ExpressionDecomposer_DecompositionType {
	_ = subExpression
	var expressionRoot pkg_rhino.Node = ExpressionDecomposer_FindExpressionRoot(subExpression)
	_ = expressionRoot
	if expressionRoot != 0 {
		return this_.priv_ExpressionDecomposer_isSubexpressionMovable_2833(expressionRoot, subExpression)
	}
	return StaticGet_ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE()
}

// Java method: ExpressionDecomposer.isSubexpressionMovable(com.google.javascript.rhino.Node expressionRoot, com.google.javascript.rhino.Node subExpression)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_isSubexpressionMovable_2833(expressionRoot pkg_rhino.Node, subExpression pkg_rhino.Node) ExpressionDecomposer_DecompositionType {
	_ = expressionRoot
	_ = subExpression
	var requiresDecomposition bool = false
	_ = requiresDecomposition
	var seenSideEffects bool = this_.astAnalyzer.MayHaveSideEffects(subExpression)
	_ = seenSideEffects
	if NodeUtil_IsOptChainNode(subExpression) && !(NodeUtil_IsEndOfFullOptChain(subExpression)) {
		requiresDecomposition = true
	}
	var child pkg_rhino.Node = subExpression
	_ = child
	{
		curAnc_ := child
		if curAnc_ == 0 {
			j2gort.ThrowNPE()
		}
		curAnc_ = curAnc_.GetParent()
		for curAnc_ != 0 {
			parent := curAnc_
			curAnc_ = curAnc_.GetParent()

			_ = parent
			if NodeUtil_IsNameDeclaration(parent) && !(child.IsFirstChildOf(parent)) {
				return StaticGet_ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE()
			}
			if parent.IsTaggedTemplateLit() && !(parent.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL())) {
				return StaticGet_ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE()
			}
			if parent == expressionRoot {
				if requiresDecomposition {
					return StaticGet_ExpressionDecomposer_DecompositionType_DECOMPOSABLE()
				} else {
					return StaticGet_ExpressionDecomposer_DecompositionType_MOVABLE()
				}
			}
			if ExpressionDecomposer_IsConditionalOp(parent) {
				if child != parent.GetFirstChild() {
					requiresDecomposition = true
				}
			} else {
				if this_.priv_ExpressionDecomposer_isSafeAssign_2836(parent, seenSideEffects) {
				} else {
					var direction ExpressionDecomposer_EvaluationDirection = ExpressionDecomposer_GetEvaluationDirection(parent)
					_ = direction
					{
						var n pkg_rhino.Node = this_.priv_ExpressionDecomposer_getFirstEvaluatedChild_2834(parent, direction)
						_ = n
						for ; n != 0; n = this_.priv_ExpressionDecomposer_getNextEvaluatedSibling_2835(n, direction) {
							if n == child {
								break
							}
							if this_.priv_ExpressionDecomposer_isExpressionTreeUnsafe_2837(n, seenSideEffects) {
								seenSideEffects = true
								requiresDecomposition = true
							}
						}
					}
					var first pkg_rhino.Node = parent.GetFirstChild()
					_ = first
					if (requiresDecomposition && parent.IsCall()) && NodeUtil_IsNormalGet(first) {
						return StaticGet_ExpressionDecomposer_DecompositionType_DECOMPOSABLE()
					}
				}
			}
			child = parent
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected.")))
}

// Java static method: ExpressionDecomposer.getEvaluationDirection(com.google.javascript.rhino.Node node)
func ExpressionDecomposer_GetEvaluationDirection(node pkg_rhino.Node) ExpressionDecomposer_EvaluationDirection {
	_ = node
	switch node.TokenSwitchOrd_() {
	case int32(104) /* Token.DESTRUCTURING_LHS */, int32(51) /* Token.ASSIGN */, int32(134) /* Token.DEFAULT_VALUE */ :
		if node.GetFirstChild().IsDestructuringPattern() {
			return StaticGet_ExpressionDecomposer_EvaluationDirection_REVERSE()
		}
		fallthrough
	default:
		return StaticGet_ExpressionDecomposer_EvaluationDirection_FORWARD()
	}
	panic("unreachable")
}

// Java method: ExpressionDecomposer.getFirstEvaluatedChild(com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.ExpressionDecomposer.EvaluationDirection direction)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_getFirstEvaluatedChild_2834(parent pkg_rhino.Node, direction ExpressionDecomposer_EvaluationDirection) pkg_rhino.Node {
	_ = parent
	_ = direction
	if direction == StaticGet_ExpressionDecomposer_EvaluationDirection_FORWARD() {
		return parent.GetFirstChild()
	} else {
		return parent.GetLastChild()
	}
}

// Java method: ExpressionDecomposer.getNextEvaluatedSibling(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.ExpressionDecomposer.EvaluationDirection direction)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_getNextEvaluatedSibling_2835(node pkg_rhino.Node, direction ExpressionDecomposer_EvaluationDirection) pkg_rhino.Node {
	_ = node
	_ = direction
	if direction == StaticGet_ExpressionDecomposer_EvaluationDirection_FORWARD() {
		return node.GetNext()
	} else {
		return node.GetPrevious()
	}
}

// Java method: ExpressionDecomposer.isSafeAssign(com.google.javascript.rhino.Node n, boolean seenSideEffects)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_isSafeAssign_2836(n pkg_rhino.Node, seenSideEffects bool) bool {
	_ = n
	_ = seenSideEffects
	if n.IsAssign() {
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		switch lhs.TokenSwitchOrd_() {
		case int32(32) /* Token.NAME */ :
			return true
		case int32(26) /* Token.GETPROP */ :
			return !(this_.priv_ExpressionDecomposer_isExpressionTreeUnsafe_2837(lhs.GetFirstChild(), seenSideEffects))
		case int32(27) /* Token.GETELEM */ :
			return (!(this_.priv_ExpressionDecomposer_isExpressionTreeUnsafe_2837(lhs.GetFirstChild(), seenSideEffects)) && !(this_.priv_ExpressionDecomposer_isExpressionTreeUnsafe_2837(lhs.GetLastChild(), seenSideEffects)))
		default:
		}
	}
	return false
}

// Java method: ExpressionDecomposer.isExpressionTreeUnsafe(com.google.javascript.rhino.Node tree, boolean followingSideEffectsExist)
func (this_ *ExpressionDecomposer) priv_ExpressionDecomposer_isExpressionTreeUnsafe_2837(tree pkg_rhino.Node, followingSideEffectsExist bool) bool {
	_ = tree
	_ = followingSideEffectsExist
	if tree.IsSpread() {
		if this_.priv_ExpressionDecomposer_isTempConstantValueName_2831(tree.GetOnlyChild()) {
			return false
		}
	}
	if followingSideEffectsExist {
		var parent pkg_rhino.Node = tree.GetParent()
		_ = parent
		if NodeUtil_IsObjectCallMethod(parent, "call") || parent.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL()) {
			var callee pkg_rhino.Node
			_ = callee
			if tree.IsGetProp() && (j2gort.CheckNotNullString(tree.GetString()) == "call") {
				callee = tree.GetFirstChild()
			} else {
				callee = tree
			}
			if tree.IsFirstChildOf(parent) && (this_.priv_ExpressionDecomposer_isTempConstantValueName_2831(callee) || j2gort.SetContainsStr(this_.knownConstantFunctions, callee.GetQualifiedName())) {
				return false
			}
		}
		if this_.enableBrokenIe11LocationAssign && StaticGet_ExpressionDecomposer_WINDOW_LOCATION_ASSIGN().Matches(tree) {
			return false
		}
		return NodeUtil_CanBeSideEffected_3(tree, j2gort.AsIface[j2gort.Ext_java_util_Set](this_.knownConstantFunctions), this_.scope)
	} else {
		return this_.astAnalyzer.MayHaveSideEffects(tree)
	}
}

// Java enum: com.google.javascript.jscomp.ExpressionDecomposer.DecompositionType
type ExpressionDecomposer_DecompositionType uint8

const (
	enum_ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE ExpressionDecomposer_DecompositionType = 1
	enum_ExpressionDecomposer_DecompositionType_MOVABLE        ExpressionDecomposer_DecompositionType = 2
	enum_ExpressionDecomposer_DecompositionType_DECOMPOSABLE   ExpressionDecomposer_DecompositionType = 3
)

var enumBaseTable_ExpressionDecomposer_DecompositionType = [...]j2gort.EnumBase{
	{Name_: "UNDECOMPOSABLE", Ordinal_: 0},
	{Name_: "MOVABLE", Ordinal_: 1},
	{Name_: "DECOMPOSABLE", Ordinal_: 2},
}
var enumValues_ExpressionDecomposer_DecompositionType = []ExpressionDecomposer_DecompositionType{1, 2, 3}

func StaticGet_ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE() ExpressionDecomposer_DecompositionType {
	return enum_ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE
}
func StaticGet_ExpressionDecomposer_DecompositionType_MOVABLE() ExpressionDecomposer_DecompositionType {
	return enum_ExpressionDecomposer_DecompositionType_MOVABLE
}
func StaticGet_ExpressionDecomposer_DecompositionType_DECOMPOSABLE() ExpressionDecomposer_DecompositionType {
	return enum_ExpressionDecomposer_DecompositionType_DECOMPOSABLE
}

func ExpressionDecomposer_DecompositionType_Values() []ExpressionDecomposer_DecompositionType {
	return j2gort.CloneSlice(enumValues_ExpressionDecomposer_DecompositionType)
}

func ExpressionDecomposer_DecompositionType_ValueOf(name string) ExpressionDecomposer_DecompositionType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ExpressionDecomposer_DecompositionType {
		if enumBaseTable_ExpressionDecomposer_DecompositionType[i_].Name_ == name {
			return ExpressionDecomposer_DecompositionType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ExpressionDecomposer.DecompositionType." + name))
}

func (this_ ExpressionDecomposer_DecompositionType) IsJavaNil_() bool                              { return this_ == 0 }
func (ExpressionDecomposer_DecompositionType) HasIdentityEquals_()                                 {}
func (ExpressionDecomposer_DecompositionType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ExpressionDecomposer_DecompositionType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ExpressionDecomposer_DecompositionType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ExpressionDecomposer_DecompositionType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExpressionDecomposer$DecompositionType", "DecompositionType"
}
func (this_ ExpressionDecomposer_DecompositionType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ExpressionDecomposer_DecompositionType[this_-1]
}
func (this_ ExpressionDecomposer_DecompositionType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ExpressionDecomposer_DecompositionType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ExpressionDecomposer_DecompositionType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ExpressionDecomposer_DecompositionType[this_-1].Name_
}
func (this_ ExpressionDecomposer_DecompositionType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ExpressionDecomposer_DecompositionType[this_-1].Name_
}
func (this_ ExpressionDecomposer_DecompositionType) Equals(other_ any) bool {
	o_, ok_ := other_.(ExpressionDecomposer_DecompositionType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ExpressionDecomposer_DecompositionType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ExpressionDecomposer_DecompositionType[this_-1].HashCode()
}
func (this_ ExpressionDecomposer_DecompositionType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ExpressionDecomposer_DecompositionType[this_-1].CompareTo(other_)
}
func (this_ ExpressionDecomposer_DecompositionType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ExpressionDecomposer_DecompositionType[this_-1].Name_)
}
func init() {
	vals_ := []any{ExpressionDecomposer_DecompositionType(1), ExpressionDecomposer_DecompositionType(2), ExpressionDecomposer_DecompositionType(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ExpressionDecomposer.DecompositionType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ExpressionDecomposer$DecompositionType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExpressionDecomposer.DecompositionType", "com.google.javascript.jscomp.ExpressionDecomposer$DecompositionType", func(v_ any) bool {
		return j2gort.InstanceOf[ExpressionDecomposer_DecompositionType](v_)
	})
}

// Java enum: com.google.javascript.jscomp.ExpressionDecomposer.Workaround
type ExpressionDecomposer_Workaround uint8

const (
	enum_ExpressionDecomposer_Workaround_BROKEN_IE11_LOCATION_ASSIGN ExpressionDecomposer_Workaround = 1
)

var enumBaseTable_ExpressionDecomposer_Workaround = [...]j2gort.EnumBase{
	{Name_: "BROKEN_IE11_LOCATION_ASSIGN", Ordinal_: 0},
}
var enumValues_ExpressionDecomposer_Workaround = []ExpressionDecomposer_Workaround{1}

func StaticGet_ExpressionDecomposer_Workaround_BROKEN_IE11_LOCATION_ASSIGN() ExpressionDecomposer_Workaround {
	return enum_ExpressionDecomposer_Workaround_BROKEN_IE11_LOCATION_ASSIGN
}

func ExpressionDecomposer_Workaround_Values() []ExpressionDecomposer_Workaround {
	return j2gort.CloneSlice(enumValues_ExpressionDecomposer_Workaround)
}

func ExpressionDecomposer_Workaround_ValueOf(name string) ExpressionDecomposer_Workaround {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ExpressionDecomposer_Workaround {
		if enumBaseTable_ExpressionDecomposer_Workaround[i_].Name_ == name {
			return ExpressionDecomposer_Workaround(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ExpressionDecomposer.Workaround." + name))
}

func (this_ ExpressionDecomposer_Workaround) IsJavaNil_() bool                              { return this_ == 0 }
func (ExpressionDecomposer_Workaround) HasIdentityEquals_()                                 {}
func (ExpressionDecomposer_Workaround) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ExpressionDecomposer_Workaround) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ExpressionDecomposer_Workaround) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ExpressionDecomposer_Workaround) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExpressionDecomposer$Workaround", "Workaround"
}
func (this_ ExpressionDecomposer_Workaround) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ExpressionDecomposer_Workaround[this_-1]
}
func (this_ ExpressionDecomposer_Workaround) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ExpressionDecomposer_Workaround) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ExpressionDecomposer_Workaround) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ExpressionDecomposer_Workaround[this_-1].Name_
}
func (this_ ExpressionDecomposer_Workaround) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ExpressionDecomposer_Workaround[this_-1].Name_
}
func (this_ ExpressionDecomposer_Workaround) Equals(other_ any) bool {
	o_, ok_ := other_.(ExpressionDecomposer_Workaround)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ExpressionDecomposer_Workaround) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ExpressionDecomposer_Workaround[this_-1].HashCode()
}
func (this_ ExpressionDecomposer_Workaround) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ExpressionDecomposer_Workaround[this_-1].CompareTo(other_)
}
func (this_ ExpressionDecomposer_Workaround) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ExpressionDecomposer_Workaround[this_-1].Name_)
}
func init() {
	vals_ := []any{ExpressionDecomposer_Workaround(1)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ExpressionDecomposer.Workaround", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ExpressionDecomposer$Workaround", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExpressionDecomposer.Workaround", "com.google.javascript.jscomp.ExpressionDecomposer$Workaround", func(v_ any) bool {
		return j2gort.InstanceOf[ExpressionDecomposer_Workaround](v_)
	})
}

// Java class: com.google.javascript.jscomp.ExpressionDecomposer.DecompositionState
type ExpressionDecomposer_DecompositionState struct {
	extractBeforeStatement pkg_rhino.Node `gson:"extractBeforeStatement"`
	idHash_                atomic.Int32
	sideEffects            bool `gson:"sideEffects"`
}

func (this_ *ExpressionDecomposer_DecompositionState) BindSelf_(s any) {
	_ = s
}
func (this_ *ExpressionDecomposer_DecompositionState) IsJavaNil_() bool { return this_ == nil }
func (*ExpressionDecomposer_DecompositionState) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExpressionDecomposer$DecompositionState", "DecompositionState"
}
func (this_ *ExpressionDecomposer_DecompositionState) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ExpressionDecomposer_DecompositionState) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExpressionDecomposer.DecompositionState", "com.google.javascript.jscomp.ExpressionDecomposer$DecompositionState", func(v_ any) bool {
		return j2gort.InstanceOf[*ExpressionDecomposer_DecompositionState](v_)
	})
}

func (this_ *ExpressionDecomposer_DecompositionState) initFields_ExpressionDecomposer_DecompositionState() {
	_ = this_
}

// Java constructor: DecompositionState()
func NewExpressionDecomposer_DecompositionState() *ExpressionDecomposer_DecompositionState {
	this_ := &ExpressionDecomposer_DecompositionState{}
	InitExpressionDecomposer_DecompositionState(this_)
	return this_
}

func InitExpressionDecomposer_DecompositionState(this_ *ExpressionDecomposer_DecompositionState) {
	this_.initFields_ExpressionDecomposer_DecompositionState()
}

// Java method: DecompositionState.toString()
func (this_ *ExpressionDecomposer_DecompositionState) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_bool("sideEffects", this_.sideEffects).Add_2_string_any("extractBeforeStatement", j2gort.BoxNode(this_.extractBeforeStatement)).String()
}

// Java enum: com.google.javascript.jscomp.ExpressionDecomposer.EvaluationDirection
type ExpressionDecomposer_EvaluationDirection uint8

const (
	enum_ExpressionDecomposer_EvaluationDirection_FORWARD ExpressionDecomposer_EvaluationDirection = 1
	enum_ExpressionDecomposer_EvaluationDirection_REVERSE ExpressionDecomposer_EvaluationDirection = 2
)

var enumBaseTable_ExpressionDecomposer_EvaluationDirection = [...]j2gort.EnumBase{
	{Name_: "FORWARD", Ordinal_: 0},
	{Name_: "REVERSE", Ordinal_: 1},
}
var enumValues_ExpressionDecomposer_EvaluationDirection = []ExpressionDecomposer_EvaluationDirection{1, 2}

func StaticGet_ExpressionDecomposer_EvaluationDirection_FORWARD() ExpressionDecomposer_EvaluationDirection {
	return enum_ExpressionDecomposer_EvaluationDirection_FORWARD
}
func StaticGet_ExpressionDecomposer_EvaluationDirection_REVERSE() ExpressionDecomposer_EvaluationDirection {
	return enum_ExpressionDecomposer_EvaluationDirection_REVERSE
}

func ExpressionDecomposer_EvaluationDirection_Values() []ExpressionDecomposer_EvaluationDirection {
	return j2gort.CloneSlice(enumValues_ExpressionDecomposer_EvaluationDirection)
}

func ExpressionDecomposer_EvaluationDirection_ValueOf(name string) ExpressionDecomposer_EvaluationDirection {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ExpressionDecomposer_EvaluationDirection {
		if enumBaseTable_ExpressionDecomposer_EvaluationDirection[i_].Name_ == name {
			return ExpressionDecomposer_EvaluationDirection(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ExpressionDecomposer.EvaluationDirection." + name))
}

func (this_ ExpressionDecomposer_EvaluationDirection) IsJavaNil_() bool                { return this_ == 0 }
func (ExpressionDecomposer_EvaluationDirection) HasIdentityEquals_()                   {}
func (ExpressionDecomposer_EvaluationDirection) Implements_j2gort_Ext_java_lang_Enum() {}
func (ExpressionDecomposer_EvaluationDirection) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ ExpressionDecomposer_EvaluationDirection) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ExpressionDecomposer_EvaluationDirection) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExpressionDecomposer$EvaluationDirection", "EvaluationDirection"
}
func (this_ ExpressionDecomposer_EvaluationDirection) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ExpressionDecomposer_EvaluationDirection[this_-1]
}
func (this_ ExpressionDecomposer_EvaluationDirection) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ExpressionDecomposer_EvaluationDirection) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ExpressionDecomposer_EvaluationDirection) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ExpressionDecomposer_EvaluationDirection[this_-1].Name_
}
func (this_ ExpressionDecomposer_EvaluationDirection) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ExpressionDecomposer_EvaluationDirection[this_-1].Name_
}
func (this_ ExpressionDecomposer_EvaluationDirection) Equals(other_ any) bool {
	o_, ok_ := other_.(ExpressionDecomposer_EvaluationDirection)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ExpressionDecomposer_EvaluationDirection) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ExpressionDecomposer_EvaluationDirection[this_-1].HashCode()
}
func (this_ ExpressionDecomposer_EvaluationDirection) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ExpressionDecomposer_EvaluationDirection[this_-1].CompareTo(other_)
}
func (this_ ExpressionDecomposer_EvaluationDirection) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ExpressionDecomposer_EvaluationDirection[this_-1].Name_)
}
func init() {
	vals_ := []any{ExpressionDecomposer_EvaluationDirection(1), ExpressionDecomposer_EvaluationDirection(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ExpressionDecomposer.EvaluationDirection", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ExpressionDecomposer$EvaluationDirection", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExpressionDecomposer.EvaluationDirection", "com.google.javascript.jscomp.ExpressionDecomposer$EvaluationDirection", func(v_ any) bool {
		return j2gort.InstanceOf[ExpressionDecomposer_EvaluationDirection](v_)
	})
}

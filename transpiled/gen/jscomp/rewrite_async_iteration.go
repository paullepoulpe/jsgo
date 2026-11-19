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
// Port of com/google/javascript/jscomp/RewriteAsyncIteration.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.RewriteAsyncIteration
type RewriteAsyncIteration struct {
	asyncGeneratorWrapper RuntimeJsLibManager_JsLibField                      `gson:"asyncGeneratorWrapper"`
	actionRecordName      RuntimeJsLibManager_JsLibField                      `gson:"actionRecordName"`
	actionEnumAwait       RuntimeJsLibManager_JsLibField                      `gson:"actionEnumAwait"`
	actionEnumYield       RuntimeJsLibManager_JsLibField                      `gson:"actionEnumYield"`
	actionEnumYieldStar   RuntimeJsLibManager_JsLibField                      `gson:"actionEnumYieldStar"`
	compiler              AbstractCompiler                                    `gson:"compiler"`
	contextStack          j2gort.Deque[*RewriteAsyncIteration_LexicalContext] `gson:"contextStack"`
	astFactory            *AstFactory                                         `gson:"astFactory"`
	namespace             pkg_rhino.StaticScope                               `gson:"namespace"`
	nextForAwaitId        int32                                               `gson:"nextForAwaitId"`
	idHash_               atomic.Int32
}

func (this_ *RewriteAsyncIteration) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteAsyncIteration) IsJavaNil_() bool { return this_ == nil }
func (*RewriteAsyncIteration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteAsyncIteration", "RewriteAsyncIteration"
}
func (this_ *RewriteAsyncIteration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteAsyncIteration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteAsyncIteration", "com.google.javascript.jscomp.RewriteAsyncIteration", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteAsyncIteration](v_)
	})
}
func (*RewriteAsyncIteration) Implements_jscomp_NodeTraversal_Callback() {}
func (*RewriteAsyncIteration) Implements_jscomp_CompilerPass()           {}

var static_RewriteAsyncIteration_transpiledFeatures *pkg_parser.FeatureSet

func StaticGet_RewriteAsyncIteration_TranspiledFeatures() *pkg_parser.FeatureSet {
	clinit_RewriteAsyncIteration()
	return static_RewriteAsyncIteration_transpiledFeatures
}
func StaticSet_RewriteAsyncIteration_TranspiledFeatures(v *pkg_parser.FeatureSet) {
	clinit_RewriteAsyncIteration()
	static_RewriteAsyncIteration_transpiledFeatures = v
}

var static_RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN *DiagnosticType

func StaticGet_RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN() *DiagnosticType {
	clinit_RewriteAsyncIteration()
	return static_RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN
}
func StaticSet_RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN(v *DiagnosticType) {
	clinit_RewriteAsyncIteration()
	static_RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN = v
}
func StaticGet_RewriteAsyncIteration_FOR_AWAIT_ITERATOR_TEMP_NAME() string {
	return "$jscomp$forAwait$tempIterator"
}
func StaticGet_RewriteAsyncIteration_FOR_AWAIT_RESULT_TEMP_NAME() string {
	return "$jscomp$forAwait$tempResult"
}
func StaticGet_RewriteAsyncIteration_FOR_AWAIT_ERROR_RESULT_TEMP_NAME() string {
	return "$jscomp$forAwait$errResult"
}
func StaticGet_RewriteAsyncIteration_FOR_AWAIT_CATCH_PARAM_TEMP_NAME() string {
	return "$jscomp$forAwait$catchErrParam"
}
func StaticGet_RewriteAsyncIteration_FOR_AWAIT_RETURN_FN_TEMP_NAME() string {
	return "$jscomp$forAwait$retFn"
}
func StaticGet_RewriteAsyncIteration_THIS_VAR_NAME() string { return "$jscomp$asyncIter$this$" }
func StaticGet_RewriteAsyncIteration_ARGUMENTS_VAR_NAME() string {
	return "$jscomp$asyncIter$arguments"
}
func StaticGet_RewriteAsyncIteration_SUPER_PROP_GETTER_PREFIX() string {
	return "$jscomp$asyncIter$super$get$"
}

var clinitDone_RewriteAsyncIteration uint32
var clinitState_RewriteAsyncIteration j2gort.ClinitState

func clinit_RewriteAsyncIteration() {
	if atomic.LoadUint32(&clinitDone_RewriteAsyncIteration) != 0 {
		return
	}
	clinitSlow_RewriteAsyncIteration()
}

func clinitSlow_RewriteAsyncIteration() {
	if !j2gort.EnterClinit(&clinitState_RewriteAsyncIteration, &clinitDone_RewriteAsyncIteration) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RewriteAsyncIteration, &clinitDone_RewriteAsyncIteration)
	static_RewriteAsyncIteration_transpiledFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7200([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_GENERATORS(), pkg_parser.StaticGet_FeatureSet_Feature_FOR_AWAIT_OF()})
	static_RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN = DiagnosticType_Error("JSC_CANNOT_CONVERT_ASYNCGEN", "Cannot convert async generator. {0}")
}

func (this_ *RewriteAsyncIteration) initFields_RewriteAsyncIteration() {
	_ = this_

}

// Java constructor: RewriteAsyncIteration(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.AstFactory astFactory, com.google.javascript.rhino.StaticScope namespace)
func NewRewriteAsyncIteration(compiler AbstractCompiler, astFactory *AstFactory, namespace pkg_rhino.StaticScope) *RewriteAsyncIteration {
	clinit_RewriteAsyncIteration()
	this_ := &RewriteAsyncIteration{}
	InitRewriteAsyncIteration(this_, compiler, astFactory, namespace)
	return this_
}

func InitRewriteAsyncIteration(this_ *RewriteAsyncIteration, compiler AbstractCompiler, astFactory *AstFactory, namespace pkg_rhino.StaticScope) {
	_ = compiler
	_ = astFactory
	_ = namespace
	this_.initFields_RewriteAsyncIteration()
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.astFactory = j2gort.CheckNotNullPtr(astFactory)
	this_.namespace = j2gort.CheckNotNullIface(namespace)
	this_.contextStack = j2gort.Deque[*RewriteAsyncIteration_LexicalContext]{}
	var runtimeLibManager *RuntimeJsLibManager = compiler.GetRuntimeJsLibManager()
	_ = runtimeLibManager
	this_.asyncGeneratorWrapper = runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper")
	this_.actionRecordName = runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionRecord")
	this_.actionEnumAwait = runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionEnum.AWAIT_VALUE")
	this_.actionEnumYield = runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE")
	this_.actionEnumYieldStar = runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_STAR")
}

// Java static method: RewriteAsyncIteration.create(com.google.javascript.jscomp.AbstractCompiler compiler)
func RewriteAsyncIteration_Create(compiler AbstractCompiler) *RewriteAsyncIteration {
	clinit_RewriteAsyncIteration()
	_ = compiler
	var astFactory *AstFactory = compiler.CreateAstFactory()
	_ = astFactory
	var namespace pkg_rhino.StaticScope = compiler.GetTranspilationNamespace()
	_ = namespace
	return NewRewriteAsyncIteration(compiler, astFactory, namespace)
}

// Java method: RewriteAsyncIteration.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewriteAsyncIteration) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.contextStack.IsEmpty()) {
		j2gort.CheckState(false)
	}
	this_.contextStack.Push(RewriteAsyncIteration_LexicalContext_NewGlobalContext(root))
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_RewriteAsyncIteration_TranspiledFeatures(), []NodeTraversal_Callback{this_})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_RewriteAsyncIteration_TranspiledFeatures())
	if !(this_.contextStack.GetFirst().function == 0) {
		j2gort.CheckState(false)
	}
	this_.contextStack.RemoveFirst()
	if !(this_.contextStack.IsEmpty()) {
		j2gort.CheckState(false)
	}
}

// Java method: RewriteAsyncIteration.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteAsyncIteration) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if n.IsFunction() {
		this_.contextStack.Push(RewriteAsyncIteration_LexicalContext_NewContextForFunction(this_.contextStack.GetFirst(), n, this_.compiler))
	} else if n.IsParamList() {
		this_.contextStack.Push(RewriteAsyncIteration_LexicalContext_NewContextForParamList(this_.contextStack.GetFirst(), n, this_.compiler))
	}
	return true
}

// Java method: RewriteAsyncIteration.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteAsyncIteration) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var ctx *RewriteAsyncIteration_LexicalContext = this_.contextStack.GetFirst()
	_ = ctx
	switch n.TokenSwitchOrd_() {
	case int32(49) /* Token.PARAM_LIST */ :
		{
			var iifeCond_1 bool
			{

				rn_1 := n
				an_2 := ctx.contextRoot
				iifeCond_1 = j2gort.CheckNotNullNode(rn_1) == an_2
			}
			if !(iifeCond_1) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
			}
		}
		this_.contextStack.Pop()
	case int32(73) /* Token.FUNCTION */ :
		{
			var iifeCond_1 bool
			{

				rn_3 := n
				an_4 := ctx.contextRoot
				iifeCond_1 = j2gort.CheckNotNullNode(rn_3) == an_4
			}
			if !(iifeCond_1) {
				j2gort.CheckState(false)
			}
		}
		if n.IsAsyncGeneratorFunction() {
			this_.priv_RewriteAsyncIteration_convertAsyncGenerator_4953(n)
			this_.priv_RewriteAsyncIteration_prependTempVarDeclarations_4966(ctx, t)
		}
		this_.contextStack.Pop()
	case int32(115) /* Token.AWAIT */ :
		j2gort.CheckNotNullNode(ctx.function)
		if ctx.function.IsAsyncGeneratorFunction() {
			this_.priv_RewriteAsyncIteration_convertAwaitOfAsyncGenerator_4954(ctx, n)
		}
	case int32(114) /* Token.YIELD */ :
		j2gort.CheckNotNullNode(ctx.function)
		if ctx.function.IsAsyncGeneratorFunction() {
			this_.priv_RewriteAsyncIteration_convertYieldOfAsyncGenerator_4955(ctx, n)
		}
	case int32(0) /* Token.RETURN */ :
		j2gort.CheckNotNullNode(ctx.function)
		if ctx.function.IsAsyncGeneratorFunction() {
			this_.priv_RewriteAsyncIteration_convertReturnOfAsyncGenerator_4956(ctx, n)
		}
	case int32(113) /* Token.FOR_AWAIT_OF */ :
		j2gort.CheckNotNullNode(ctx.function)
		if !(ctx.function.IsAsyncFunction()) {
			j2gort.CheckState(false)
		}
		this_.priv_RewriteAsyncIteration_replaceForAwaitOf_4957(ctx, n)
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
	case int32(37) /* Token.THIS */ :
		if ctx.MustReplaceThisSuperArgs() {
			this_.priv_RewriteAsyncIteration_replaceThis_4963(ctx, n)
		}
	case int32(32) /* Token.NAME */ :
		if ctx.MustReplaceThisSuperArgs() && n.MatchesName_1_string("arguments") {
			this_.priv_RewriteAsyncIteration_replaceArguments_4964(ctx, n)
		}
	case int32(110) /* Token.SUPER */ :
		if ctx.MustReplaceThisSuperArgs() {
			this_.priv_RewriteAsyncIteration_replaceSuper_4965(ctx, n, parent)
		}
	default:
	}
}

// Java method: RewriteAsyncIteration.convertAsyncGenerator(com.google.javascript.rhino.Node originalFunction)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_convertAsyncGenerator_4953(originalFunction pkg_rhino.Node) {
	_ = originalFunction
	j2gort.CheckNotNullNode(originalFunction)
	if !(originalFunction.IsAsyncGeneratorFunction()) {
		j2gort.CheckState(false)
	}
	var asyncGeneratorWrapperRef pkg_rhino.Node = this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.asyncGeneratorWrapper)
	_ = asyncGeneratorWrapperRef
	var innerFunction pkg_rhino.Node = this_.astFactory.CreateEmptyAsyncGeneratorWrapperArgument(nil)
	_ = innerFunction
	var innerBlock pkg_rhino.Node = originalFunction.GetLastChild()
	_ = innerBlock
	innerBlock.Detach()
	innerFunction.GetLastChild().ReplaceWith(innerBlock)
	var outerBlock pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{this_.astFactory.CreateReturn(this_.astFactory.CreateNewNode(asyncGeneratorWrapperRef, []pkg_rhino.Node{this_.astFactory.CreateCall(innerFunction, AstFactory_S_Type_1_ColorId(pkg_colors.StaticGet_StandardColors_GENERATOR_ID()), []pkg_rhino.Node{})}))})
	_ = outerBlock
	originalFunction.AddChildToBack(outerBlock)
	originalFunction.SetIsAsyncFunction(false)
	originalFunction.SetIsGeneratorFunction(false)
	originalFunction.SrcrefTreeIfMissing(originalFunction)
	this_.compiler.ReportChangeToChangeScope(originalFunction)
	this_.compiler.ReportChangeToChangeScope(innerFunction)
}

// Java method: RewriteAsyncIteration.convertAwaitOfAsyncGenerator(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, com.google.javascript.rhino.Node awaitNode)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_convertAwaitOfAsyncGenerator_4954(ctx *RewriteAsyncIteration_LexicalContext, awaitNode pkg_rhino.Node) {
	_ = ctx
	_ = awaitNode
	j2gort.CheckNotNullNode(awaitNode)
	if !(awaitNode.IsAwait()) {
		j2gort.CheckState(false)
	}
	if !((ctx != nil) && (ctx.function != 0)) {
		j2gort.CheckState(false)
	}
	if !(ctx.function.IsAsyncGeneratorFunction()) {
		j2gort.CheckState(false)
	}
	var expression pkg_rhino.Node = awaitNode.RemoveFirstChild()
	_ = expression
	j2gort.CheckNotNullNode(expression)
	var newActionRecord pkg_rhino.Node = this_.astFactory.CreateNewNode(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionRecordName), []pkg_rhino.Node{this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionEnumAwait), expression})
	_ = newActionRecord
	newActionRecord.SrcrefTreeIfMissing(awaitNode)
	awaitNode.AddChildToFront(newActionRecord)
	awaitNode.SetToken(pkg_rhino.StaticGet_Token_YIELD())
}

// Java method: RewriteAsyncIteration.convertYieldOfAsyncGenerator(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, com.google.javascript.rhino.Node yieldNode)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_convertYieldOfAsyncGenerator_4955(ctx *RewriteAsyncIteration_LexicalContext, yieldNode pkg_rhino.Node) {
	_ = ctx
	_ = yieldNode
	j2gort.CheckNotNullNode(yieldNode)
	if !(yieldNode.IsYield()) {
		j2gort.CheckState(false)
	}
	if !((ctx != nil) && (ctx.function != 0)) {
		j2gort.CheckState(false)
	}
	if !(ctx.function.IsAsyncGeneratorFunction()) {
		j2gort.CheckState(false)
	}
	var expression pkg_rhino.Node = yieldNode.RemoveFirstChild()
	_ = expression
	var newActionRecord pkg_rhino.Node = this_.astFactory.CreateNewNode(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionRecordName), []pkg_rhino.Node{})
	_ = newActionRecord
	if yieldNode.IsYieldAll() {
		j2gort.CheckNotNullNode(expression)
		newActionRecord.AddChildToBack(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionEnumYieldStar))
		newActionRecord.AddChildToBack(expression)
	} else {
		if expression == 0 {
			expression = NodeUtil_NewUndefinedNode(pkg_rhino.Node(0))
		}
		newActionRecord.AddChildToBack(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionEnumYield))
		newActionRecord.AddChildToBack(expression)
	}
	newActionRecord.SrcrefTreeIfMissing(yieldNode)
	yieldNode.AddChildToFront(newActionRecord)
	yieldNode.PutBooleanProp(pkg_rhino.StaticGet_Node_YIELD_ALL(), false)
}

// Java method: RewriteAsyncIteration.convertReturnOfAsyncGenerator(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, com.google.javascript.rhino.Node returnNode)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_convertReturnOfAsyncGenerator_4956(ctx *RewriteAsyncIteration_LexicalContext, returnNode pkg_rhino.Node) {
	_ = ctx
	_ = returnNode
	j2gort.CheckNotNullNode(returnNode)
	if !(returnNode.IsReturn()) {
		j2gort.CheckState(false)
	}
	if !((ctx != nil) && (ctx.function != 0)) {
		j2gort.CheckState(false)
	}
	if !(ctx.function.IsAsyncGeneratorFunction()) {
		j2gort.CheckState(false)
	}
	var expression pkg_rhino.Node = returnNode.RemoveFirstChild()
	_ = expression
	var newActionRecord pkg_rhino.Node = this_.astFactory.CreateNewNode(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionRecordName), []pkg_rhino.Node{})
	_ = newActionRecord
	if expression == 0 {
		expression = NodeUtil_NewUndefinedNode(pkg_rhino.Node(0))
	}
	newActionRecord.AddChildToBack(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionEnumYield))
	newActionRecord.AddChildToBack(expression)
	newActionRecord.SrcrefTreeIfMissing(returnNode)
	returnNode.AddChildToFront(newActionRecord)
}

// Java method: RewriteAsyncIteration.replaceForAwaitOf(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, com.google.javascript.rhino.Node forAwaitOf)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_replaceForAwaitOf_4957(ctx *RewriteAsyncIteration_LexicalContext, forAwaitOf pkg_rhino.Node) {
	_ = ctx
	_ = forAwaitOf
	var forAwaitId int32 = func() int32 {
		old_5 := this_.nextForAwaitId
		nw_6 := old_5 + 1
		this_.nextForAwaitId = nw_6
		return old_5
	}()
	_ = forAwaitId
	var iteratorTempName string = ("$jscomp$forAwait$tempIterator" + j2gort.StrValInt32(forAwaitId))
	_ = iteratorTempName
	var resultTempName string = ("$jscomp$forAwait$tempResult" + j2gort.StrValInt32(forAwaitId))
	_ = resultTempName
	var errorResultTempName string = ("$jscomp$forAwait$errResult" + j2gort.StrValInt32(forAwaitId))
	_ = errorResultTempName
	var catchErrorParamTempName string = ("$jscomp$forAwait$catchErrParam" + j2gort.StrValInt32(forAwaitId))
	_ = catchErrorParamTempName
	var returnFuncTempName string = ("$jscomp$forAwait$retFn" + j2gort.StrValInt32(forAwaitId))
	_ = returnFuncTempName
	if !(forAwaitOf.HasParent()) {
		j2gort.CheckStateStr(false, "Cannot replace parentless for-await-of")
	}
	var forAwaitOfParent pkg_rhino.Node = forAwaitOf.GetParent()
	_ = forAwaitOfParent
	var replacementPoint pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacementPoint
	if forAwaitOfParent.IsLabel() {
		if !(forAwaitOf.IsSecondChildOf(forAwaitOfParent)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(forAwaitOfParent))
		}
		replacementPoint = forAwaitOfParent
	} else {
		replacementPoint = forAwaitOf
	}
	var lhs pkg_rhino.Node = forAwaitOf.RemoveFirstChild()
	_ = lhs
	var rhs pkg_rhino.Node = forAwaitOf.RemoveFirstChild()
	_ = rhs
	var originalBody pkg_rhino.Node = forAwaitOf.RemoveFirstChild()
	_ = originalBody
	var initializer pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_2(iteratorTempName, this_.astFactory.CreateJSCompMakeAsyncIteratorCall(rhs, this_.namespace)).SrcrefTreeIfMissing(rhs)
	_ = initializer
	var iterableResultType AstFactory_Type = AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT())
	_ = iterableResultType
	var breakIfDone pkg_rhino.Node = this_.astFactory.CreateIf_2(this_.astFactory.CreateGetProp(this_.astFactory.CreateName_2_string_AstFactory_Type(resultTempName, iterableResultType), "done", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_BOOLEAN())), this_.astFactory.CreateBlock([]pkg_rhino.Node{this_.astFactory.CreateBreak()}))
	_ = breakIfDone
	var lhsAssignment pkg_rhino.Node = pkg_rhino.Node(0)
	_ = lhsAssignment
	var resultType AstFactory_Type = nil
	_ = resultType
	if lhs.IsValidAssignmentTarget() {
		resultType = AstFactory_S_Type_1_Node(lhs)
		lhsAssignment = this_.astFactory.ExprResult(this_.astFactory.CreateAssign_2_Node_Node(lhs, this_.astFactory.CreateGetProp(this_.astFactory.CreateName_2_string_AstFactory_Type(resultTempName, iterableResultType), "value", resultType)))
	} else if NodeUtil_IsNameDeclaration(lhs) {
		var declarationTarget pkg_rhino.Node = lhs.GetFirstChild()
		_ = declarationTarget
		if declarationTarget.IsName() {
			resultType = AstFactory_S_Type_1_Node(declarationTarget)
			declarationTarget.AddChildToBack(this_.astFactory.CreateGetProp(this_.astFactory.CreateName_2_string_AstFactory_Type(resultTempName, iterableResultType), "value", resultType))
		} else {
			if !(declarationTarget.IsDestructuringLhs()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(declarationTarget))
			}
			var destructuringPattern pkg_rhino.Node = declarationTarget.GetOnlyChild()
			_ = destructuringPattern
			resultType = AstFactory_S_Type_1_Node(destructuringPattern)
			declarationTarget.AddChildToBack(this_.astFactory.CreateGetProp(this_.astFactory.CreateName_2_string_AstFactory_Type(resultTempName, iterableResultType), "value", resultType))
		}
		lhsAssignment = lhs
	} else {
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "unexpected for-await-of lhs")))
	}
	lhsAssignment.SrcrefTreeIfMissing(lhs)
	var errorResDecl pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_1(errorResultTempName).SrcrefTreeIfMissing(forAwaitOf)
	_ = errorResDecl
	var tempResultDecl pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_1(resultTempName).SrcrefTreeIfMissing(forAwaitOf)
	_ = tempResultDecl
	var returnFuncDecl pkg_rhino.Node = this_.astFactory.CreateSingleVarNameDeclaration_1(returnFuncTempName).SrcrefTreeIfMissing(forAwaitOf)
	_ = returnFuncDecl
	var resetTempResult pkg_rhino.Node = this_.astFactory.ExprResult(this_.astFactory.CreateAssign_2_string_Node(resultTempName, this_.astFactory.CreateUndefinedValue()))
	_ = resetTempResult
	var resultDeclaration pkg_rhino.Node = this_.astFactory.ExprResult(this_.astFactory.CreateAssign_2_string_Node(resultTempName, this_.priv_RewriteAsyncIteration_constructAwaitNextResult_4962(ctx, iteratorTempName, resultType, iterableResultType)))
	_ = resultDeclaration
	var newForLoop pkg_rhino.Node = this_.astFactory.CreateFor(this_.astFactory.CreateEmpty(), this_.astFactory.CreateEmpty(), this_.astFactory.CreateEmpty(), this_.astFactory.CreateBlock([]pkg_rhino.Node{resetTempResult, resultDeclaration, breakIfDone, lhsAssignment, this_.priv_RewriteAsyncIteration_ensureBlock_4961(originalBody)}))
	_ = newForLoop
	if replacementPoint.IsLabel() {
		newForLoop = this_.astFactory.CreateLabel(replacementPoint.GetFirstChild().CloneNode_0(), newForLoop)
	}
	var tryNode pkg_rhino.Node = this_.priv_RewriteAsyncIteration_createOuterTry_4958(newForLoop)
	_ = tryNode
	initializer.InsertBefore(newForLoop)
	var catchNode pkg_rhino.Node = this_.priv_RewriteAsyncIteration_createOuterCatch_4959(catchErrorParamTempName, errorResultTempName)
	_ = catchNode
	var finallyNode pkg_rhino.Node = this_.priv_RewriteAsyncIteration_createOuterFinally_4960(ctx, iterableResultType, resultType, resultTempName, returnFuncTempName, iteratorTempName, errorResultTempName)
	_ = finallyNode
	var tryCatchFinally pkg_rhino.Node = this_.astFactory.CreateTryCatchFinally(tryNode, catchNode, finallyNode)
	_ = tryCatchFinally
	replacementPoint.ReplaceWith(tryCatchFinally)
	tryCatchFinally.SrcrefTreeIfMissing(replacementPoint)
	errorResDecl.InsertBefore(tryCatchFinally)
	tempResultDecl.InsertBefore(tryCatchFinally)
	returnFuncDecl.InsertBefore(tryCatchFinally)
	this_.compiler.ReportChangeToEnclosingScope(tryCatchFinally)
}

// Java method: RewriteAsyncIteration.createOuterTry(com.google.javascript.rhino.Node newForLoop)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_createOuterTry_4958(newForLoop pkg_rhino.Node) pkg_rhino.Node {
	_ = newForLoop
	var tryNode pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = tryNode
	tryNode.AddChildToBack(newForLoop)
	return tryNode
}

// Java method: RewriteAsyncIteration.createOuterCatch(java.lang.String catchErrorParamTempName, java.lang.String errorResultTempName)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_createOuterCatch_4959(catchErrorParamTempName string, errorResultTempName string) pkg_rhino.Node {
	_ = catchErrorParamTempName
	_ = errorResultTempName
	var catchBodyStmt pkg_rhino.Node = this_.astFactory.ExprResult(this_.astFactory.CreateAssign_2_string_Node(errorResultTempName, this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{this_.astFactory.CreateStringKey("error", this_.astFactory.CreateNameWithUnknownType(catchErrorParamTempName))})))
	_ = catchBodyStmt
	var wrapperCatchBlockNode pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = wrapperCatchBlockNode
	wrapperCatchBlockNode.AddChildToBack(catchBodyStmt)
	return this_.astFactory.CreateCatch(this_.astFactory.CreateNameWithUnknownType(catchErrorParamTempName), wrapperCatchBlockNode)
}

// Java method: RewriteAsyncIteration.createOuterFinally(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, com.google.javascript.jscomp.AstFactory.Type iterableResultType, com.google.javascript.jscomp.AstFactory.Type resultType, java.lang.String resultTempName, java.lang.String returnFuncTempName, java.lang.String iteratorTempName, java.lang.String errorResultTempName)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_createOuterFinally_4960(ctx *RewriteAsyncIteration_LexicalContext, iterableResultType AstFactory_Type, resultType AstFactory_Type, resultTempName string, returnFuncTempName string, iteratorTempName string, errorResultTempName string) pkg_rhino.Node {
	_ = ctx
	_ = iterableResultType
	_ = resultType
	_ = resultTempName
	_ = returnFuncTempName
	_ = iteratorTempName
	_ = errorResultTempName
	var finallyNode pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = finallyNode
	var tmpResNameNode pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(resultTempName)
	_ = tmpResNameNode
	var tmpResDoneGetProp pkg_rhino.Node = this_.astFactory.CreateGetProp(this_.astFactory.CreateName_2_string_AstFactory_Type(resultTempName, iterableResultType), "done", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_BOOLEAN()))
	_ = tmpResDoneGetProp
	var and pkg_rhino.Node = this_.astFactory.CreateAnd(tmpResNameNode, this_.astFactory.CreateNot(tmpResDoneGetProp))
	_ = and
	var assign pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(this_.astFactory.CreateNameWithUnknownType(returnFuncTempName), this_.astFactory.CreateGetProp(this_.astFactory.CreateName_2_string_AstFactory_Type(iteratorTempName, resultType), "return", AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())))
	_ = assign
	var ifCond pkg_rhino.Node = this_.astFactory.CreateAnd(and, assign)
	_ = ifCond
	var awaitOrYieldStmt pkg_rhino.Node = pkg_rhino.Node(0)
	_ = awaitOrYieldStmt
	if ctx.function.IsAsyncGeneratorFunction() {
		awaitOrYieldStmt = this_.astFactory.ExprResult(this_.astFactory.CreateYield(iterableResultType, this_.astFactory.CreateNewNode(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionRecordName), []pkg_rhino.Node{this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionEnumAwait), this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateName_2_string_AstFactory_Type(returnFuncTempName, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), "call"), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN()), []pkg_rhino.Node{this_.astFactory.CreateName_2_string_AstFactory_Type(iteratorTempName, resultType)})})))
	} else {
		awaitOrYieldStmt = this_.astFactory.ExprResult(this_.astFactory.CreateAwait(iterableResultType, this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateName_2_string_AstFactory_Type(returnFuncTempName, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_UNKNOWN())), "call"), AstFactory_S_Type_1_ColorId(pkg_colors.StaticGet_StandardColors_PROMISE_ID()), []pkg_rhino.Node{this_.astFactory.CreateName_2_string_AstFactory_Type(iteratorTempName, resultType)})))
	}
	var ifBody pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = ifBody
	ifBody.AddChildToBack(awaitOrYieldStmt)
	var ifBlock pkg_rhino.Node = this_.astFactory.CreateIf_2(ifCond, ifBody)
	_ = ifBlock
	var innerTryBlock pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = innerTryBlock
	innerTryBlock.AddChildToBack(ifBlock)
	var innerFinallyBlock pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = innerFinallyBlock
	var secondIfBody pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = secondIfBody
	var throwStmt pkg_rhino.Node = this_.astFactory.CreateThrow(this_.astFactory.CreateGetPropWithUnknownType(this_.astFactory.CreateNameWithUnknownType(errorResultTempName), "error"))
	_ = throwStmt
	secondIfBody.AddChildToBack(throwStmt)
	var secondIfCond pkg_rhino.Node = this_.astFactory.CreateNameWithUnknownType(errorResultTempName)
	_ = secondIfCond
	var secondIfBlock pkg_rhino.Node = this_.astFactory.CreateIf_2(secondIfCond, secondIfBody)
	_ = secondIfBlock
	innerFinallyBlock.AddChildToBack(secondIfBlock)
	var finallyBody pkg_rhino.Node = this_.astFactory.CreateTryFinally(innerTryBlock, innerFinallyBlock)
	_ = finallyBody
	finallyNode.AddChildToBack(finallyBody)
	return finallyNode
}

// Java method: RewriteAsyncIteration.ensureBlock(com.google.javascript.rhino.Node possiblyBlock)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_ensureBlock_4961(possiblyBlock pkg_rhino.Node) pkg_rhino.Node {
	_ = possiblyBlock
	if possiblyBlock.IsBlock() {
		return possiblyBlock
	} else {
		return this_.astFactory.CreateBlock([]pkg_rhino.Node{possiblyBlock}).Srcref(possiblyBlock)
	}
}

// Java method: RewriteAsyncIteration.constructAwaitNextResult(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, java.lang.String iteratorTempName, com.google.javascript.jscomp.AstFactory.Type iteratorType, com.google.javascript.jscomp.AstFactory.Type iterableResultType)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_constructAwaitNextResult_4962(ctx *RewriteAsyncIteration_LexicalContext, iteratorTempName string, iteratorType AstFactory_Type, iterableResultType AstFactory_Type) pkg_rhino.Node {
	_ = ctx
	_ = iteratorTempName
	_ = iteratorType
	_ = iterableResultType
	j2gort.CheckNotNullNode(ctx.function)
	var result pkg_rhino.Node = pkg_rhino.Node(0)
	_ = result
	var iteratorTemp pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(iteratorTempName, iteratorType)
	_ = iteratorTemp
	if ctx.function.IsAsyncGeneratorFunction() {
		result = this_.astFactory.CreateYield(iterableResultType, this_.astFactory.CreateNewNode(this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionRecordName), []pkg_rhino.Node{this_.astFactory.CreateQName_2_StaticScope_RuntimeJsLibManager_JsLibField(this_.namespace, this_.actionEnumAwait), this_.astFactory.CreateCallWithUnknownType(this_.astFactory.CreateGetPropWithUnknownType(iteratorTemp, "next"), []pkg_rhino.Node{})}))
	} else {
		result = this_.astFactory.CreateAwait(iterableResultType, this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(iteratorTemp, "next"), AstFactory_S_Type_1_ColorId(pkg_colors.StaticGet_StandardColors_PROMISE_ID()), []pkg_rhino.Node{}))
	}
	return result
}

// Java method: RewriteAsyncIteration.replaceThis(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, com.google.javascript.rhino.Node n)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_replaceThis_4963(ctx *RewriteAsyncIteration_LexicalContext, n pkg_rhino.Node) {
	_ = ctx
	_ = n
	if !(n.IsThis()) {
		j2gort.CheckArgument(false)
	}
	if !((ctx != nil) && ctx.MustReplaceThisSuperArgs()) {
		j2gort.CheckArgument(false)
	}
	if !(ctx.function != 0) {
		j2gort.CheckArgumentStr(false, "Cannot prepend declarations to root scope")
	}
	j2gort.CheckNotNullPtr(ctx.thisSuperArgsContext)
	n.ReplaceWith(this_.astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$asyncIter$this$" + j2gort.StrValString(ctx.thisSuperArgsContext.uniqueId)), AstFactory_S_Type_1_Node(n)).Srcref(n))
	{
		r_7 := ctx.thisSuperArgsContext
		r_7.thisNodeToAdd = this_.astFactory.CreateThis(AstFactory_S_Type_1_Node(n))
	}
	this_.compiler.ReportChangeToChangeScope(ctx.function)
}

// Java method: RewriteAsyncIteration.replaceArguments(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, com.google.javascript.rhino.Node n)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_replaceArguments_4964(ctx *RewriteAsyncIteration_LexicalContext, n pkg_rhino.Node) {
	_ = ctx
	_ = n
	if !(n.IsName() && func() bool {
		sr_8 := "arguments"
		se_9 := n.GetString()
		return j2gort.CheckNotNullString(sr_8) == se_9
	}()) {
		j2gort.CheckArgument(false)
	}
	if !((ctx != nil) && ctx.MustReplaceThisSuperArgs()) {
		j2gort.CheckArgument(false)
	}
	if !(ctx.function != 0) {
		j2gort.CheckArgumentStr(false, "Cannot prepend declarations to root scope")
	}
	j2gort.CheckNotNullPtr(ctx.thisSuperArgsContext)
	n.ReplaceWith(this_.astFactory.CreateName_2_string_AstFactory_Type("$jscomp$asyncIter$arguments", AstFactory_S_Type_1_Node(n)).Srcref(n))
	{
		r_10 := ctx.thisSuperArgsContext
		r_10.usedArguments = true
	}
	this_.compiler.ReportChangeToChangeScope(ctx.function)
}

// Java method: RewriteAsyncIteration.replaceSuper(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_replaceSuper_4965(ctx *RewriteAsyncIteration_LexicalContext, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = ctx
	_ = n
	_ = parent
	if !(parent.IsGetProp()) {
		this_.compiler.Report_1(JSError_Make_3(parent, StaticGet_RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN(), []string{"super only allowed with getprop (like super.foo(), not super['foo']())"}))
		return
	}
	if !(n.IsSuper()) {
		j2gort.CheckArgument(false)
	}
	if !((ctx != nil) && ctx.MustReplaceThisSuperArgs()) {
		j2gort.CheckArgument(false)
	}
	if !(ctx.function != 0) {
		j2gort.CheckArgumentStr(false, "Cannot prepend declarations to root scope")
	}
	j2gort.CheckNotNullPtr(ctx.thisSuperArgsContext)
	var propertyName string = parent.GetString()
	_ = propertyName
	var propertyReplacementNameText string = ("$jscomp$asyncIter$super$get$" + j2gort.StrValString(propertyName))
	_ = propertyReplacementNameText
	var getPropReplacement pkg_rhino.Node = this_.astFactory.CreateCall(this_.astFactory.CreateName_2_string_AstFactory_Type(propertyReplacementNameText, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT())), AstFactory_S_Type_1_Node(parent), []pkg_rhino.Node{})
	_ = getPropReplacement
	var grandparent pkg_rhino.Node = parent.GetParent()
	_ = grandparent
	if grandparent.IsCall() && (grandparent.GetFirstChild() == parent) {
		getPropReplacement = this_.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "call")
		{
			r_11 := ctx.thisSuperArgsContext
			r_11.thisNodeToAdd = this_.astFactory.CreateThisForEs6ClassMember(ctx.contextRoot.GetParent())
		}
		this_.astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$asyncIter$this$" + j2gort.StrValString(ctx.thisSuperArgsContext.uniqueId)), AstFactory_S_Type_1_Node(ctx.thisSuperArgsContext.thisNodeToAdd)).Srcref(parent).InsertAfter(parent)
	} else if grandparent.IsTaggedTemplateLit() && (grandparent.GetFirstChild() == parent) {
		{
			r_12 := ctx.thisSuperArgsContext
			r_12.thisNodeToAdd = this_.astFactory.CreateThisForEs6ClassMember(ctx.contextRoot.GetParent())
		}
		var thisAlias pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(("$jscomp$asyncIter$this$" + j2gort.StrValString(ctx.thisSuperArgsContext.uniqueId)), AstFactory_S_Type_1_Node(ctx.thisSuperArgsContext.thisNodeToAdd)).Srcref(parent)
		_ = thisAlias
		getPropReplacement = this_.astFactory.CreateCallWithUnknownType(this_.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "bind"), []pkg_rhino.Node{thisAlias})
		grandparent.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	}
	getPropReplacement.SrcrefTree(parent)
	parent.ReplaceWith(getPropReplacement)
	ctx.thisSuperArgsContext.usedSuperProperties.Add(parent)
	this_.compiler.ReportChangeToChangeScope(ctx.function)
}

// Java method: RewriteAsyncIteration.prependTempVarDeclarations(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_prependTempVarDeclarations_4966(ctx *RewriteAsyncIteration_LexicalContext, t *NodeTraversal) {
	_ = ctx
	_ = t
	if !(ctx != nil) {
		j2gort.CheckArgument(false)
	}
	if !(ctx.function != 0) {
		j2gort.CheckArgumentStr(false, "Cannot prepend declarations to root scope")
	}
	j2gort.CheckNotNullPtr(ctx.thisSuperArgsContext)
	var thisSuperArgsCtx *RewriteAsyncIteration_ThisSuperArgsContext = ctx.thisSuperArgsContext
	_ = thisSuperArgsCtx
	var function pkg_rhino.Node = ctx.function
	_ = function
	var block pkg_rhino.Node = function.GetLastChild()
	_ = block
	j2gort.CheckNotNullNode(block)
	var prefixBlock pkg_rhino.Node = this_.astFactory.CreateBlock([]pkg_rhino.Node{})
	_ = prefixBlock
	if thisSuperArgsCtx.thisNodeToAdd != 0 {
		prefixBlock.AddChildToBack(this_.astFactory.CreateSingleConstNameDeclaration(("$jscomp$asyncIter$this$" + j2gort.StrValString(thisSuperArgsCtx.uniqueId)), thisSuperArgsCtx.thisNodeToAdd).SrcrefTree(block))
	}
	if thisSuperArgsCtx.usedArguments {
		prefixBlock.AddChildToBack(this_.astFactory.CreateSingleConstNameDeclaration("$jscomp$asyncIter$arguments", this_.astFactory.CreateArgumentsReference()).SrcrefTree(block))
	}
	{
		cur_13 := thisSuperArgsCtx.usedSuperProperties.Iter()
		for cur_13.HasNext() {
			replacedMethodReference := cur_13.
				Next()
			_ = replacedMethodReference
			prefixBlock.AddChildToBack(this_.priv_RewriteAsyncIteration_createSuperMethodReferenceGetter_4967(replacedMethodReference, t))
		}
	}
	prefixBlock.SrcrefTreeIfMissing(block)
	block.AddChildrenToFront(prefixBlock.RemoveChildren())
	if ((thisSuperArgsCtx.thisNodeToAdd != 0) || thisSuperArgsCtx.usedArguments) || !(thisSuperArgsCtx.usedSuperProperties.IsEmpty()) {
		this_.compiler.ReportChangeToChangeScope(function)
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS(), this_.compiler)
	}
}

// Java method: RewriteAsyncIteration.createSuperMethodReferenceGetter(com.google.javascript.rhino.Node replacedMethodReference, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *RewriteAsyncIteration) priv_RewriteAsyncIteration_createSuperMethodReferenceGetter_4967(replacedMethodReference pkg_rhino.Node, t *NodeTraversal) pkg_rhino.Node {
	_ = replacedMethodReference
	_ = t
	var typeOfSuper AstFactory_Type = AstFactory_S_Type_1_Node(replacedMethodReference.GetFirstChild())
	_ = typeOfSuper
	var superReference pkg_rhino.Node = this_.astFactory.CreateSuper(typeOfSuper)
	_ = superReference
	var replacedMethodName string = replacedMethodReference.GetString()
	_ = replacedMethodName
	var arrowFunction pkg_rhino.Node = this_.astFactory.CreateZeroArgArrowFunctionForExpression(this_.astFactory.CreateBlock([]pkg_rhino.Node{this_.astFactory.CreateReturn(this_.astFactory.CreateGetProp(superReference, replacedMethodName, AstFactory_S_Type_1_Node(replacedMethodReference)))}))
	_ = arrowFunction
	this_.compiler.ReportChangeToChangeScope(arrowFunction)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS(), this_.compiler)
	var superReplacementName string = ("$jscomp$asyncIter$super$get$" + j2gort.StrValString(replacedMethodName))
	_ = superReplacementName
	return this_.astFactory.CreateSingleConstNameDeclaration(superReplacementName, arrowFunction)
}

// Java class: com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext
type RewriteAsyncIteration_LexicalContext struct {
	thisSuperArgsContext *RewriteAsyncIteration_ThisSuperArgsContext `gson:"thisSuperArgsContext"`
	contextRoot          pkg_rhino.Node                              `gson:"contextRoot"`
	function             pkg_rhino.Node                              `gson:"function"`
	idHash_              atomic.Int32
}

func (this_ *RewriteAsyncIteration_LexicalContext) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteAsyncIteration_LexicalContext) IsJavaNil_() bool { return this_ == nil }
func (*RewriteAsyncIteration_LexicalContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteAsyncIteration$LexicalContext", "LexicalContext"
}
func (this_ *RewriteAsyncIteration_LexicalContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteAsyncIteration_LexicalContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext", "com.google.javascript.jscomp.RewriteAsyncIteration$LexicalContext", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteAsyncIteration_LexicalContext](v_)
	})
}

func (this_ *RewriteAsyncIteration_LexicalContext) initFields_RewriteAsyncIteration_LexicalContext() {
	_ = this_
}

// Java constructor: LexicalContext(com.google.javascript.rhino.Node contextRoot)
func NewRewriteAsyncIteration_LexicalContext_1(contextRoot pkg_rhino.Node) *RewriteAsyncIteration_LexicalContext {
	this_ := &RewriteAsyncIteration_LexicalContext{}
	InitRewriteAsyncIteration_LexicalContext_1(this_, contextRoot)
	return this_
}

func InitRewriteAsyncIteration_LexicalContext_1(this_ *RewriteAsyncIteration_LexicalContext, contextRoot pkg_rhino.Node) {
	_ = contextRoot
	this_.initFields_RewriteAsyncIteration_LexicalContext()
	this_.contextRoot = j2gort.CheckNotNullNode(contextRoot)
	this_.function = pkg_rhino.Node(0)
	this_.thisSuperArgsContext = nil
}

// Java constructor: LexicalContext(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext parent, com.google.javascript.rhino.Node contextRoot, com.google.javascript.rhino.Node function, com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRewriteAsyncIteration_LexicalContext_4(parent *RewriteAsyncIteration_LexicalContext, contextRoot pkg_rhino.Node, function pkg_rhino.Node, compiler AbstractCompiler) *RewriteAsyncIteration_LexicalContext {
	this_ := &RewriteAsyncIteration_LexicalContext{}
	InitRewriteAsyncIteration_LexicalContext_4(this_, parent, contextRoot, function, compiler)
	return this_
}

func InitRewriteAsyncIteration_LexicalContext_4(this_ *RewriteAsyncIteration_LexicalContext, parent *RewriteAsyncIteration_LexicalContext, contextRoot pkg_rhino.Node, function pkg_rhino.Node, compiler AbstractCompiler) {
	_ = parent
	_ = contextRoot
	_ = function
	_ = compiler
	this_.initFields_RewriteAsyncIteration_LexicalContext()
	j2gort.CheckNotNullPtr(parent)
	j2gort.CheckNotNullNode(contextRoot)
	if !((contextRoot == function) || contextRoot.IsParamList()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(contextRoot))
	}
	j2gort.CheckNotNullNode(function)
	if !(function.IsFunction()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(function))
	}
	this_.contextRoot = contextRoot
	this_.function = function
	if function.IsArrowFunction() {
		this_.thisSuperArgsContext = parent.thisSuperArgsContext
	} else if contextRoot.IsFunction() {
		var newUniqueId string = compiler.GetUniqueIdSupplier().GetUniqueId(compiler.GetInput(NodeUtil_GetInputId(contextRoot)))
		_ = newUniqueId
		this_.thisSuperArgsContext = NewRewriteAsyncIteration_ThisSuperArgsContext(this_, newUniqueId)
	} else {
		this_.thisSuperArgsContext = nil
	}
}

// Java static method: LexicalContext.newGlobalContext(com.google.javascript.rhino.Node contextRoot)
func RewriteAsyncIteration_LexicalContext_NewGlobalContext(contextRoot pkg_rhino.Node) *RewriteAsyncIteration_LexicalContext {
	_ = contextRoot
	return NewRewriteAsyncIteration_LexicalContext_1(contextRoot)
}

// Java static method: LexicalContext.newContextForFunction(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext parent, com.google.javascript.rhino.Node function, com.google.javascript.jscomp.AbstractCompiler compiler)
func RewriteAsyncIteration_LexicalContext_NewContextForFunction(parent *RewriteAsyncIteration_LexicalContext, function pkg_rhino.Node, compiler AbstractCompiler) *RewriteAsyncIteration_LexicalContext {
	_ = parent
	_ = function
	_ = compiler
	return NewRewriteAsyncIteration_LexicalContext_4(parent, function, function, compiler)
}

// Java static method: LexicalContext.newContextForParamList(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext parent, com.google.javascript.rhino.Node paramList, com.google.javascript.jscomp.AbstractCompiler compiler)
func RewriteAsyncIteration_LexicalContext_NewContextForParamList(parent *RewriteAsyncIteration_LexicalContext, paramList pkg_rhino.Node, compiler AbstractCompiler) *RewriteAsyncIteration_LexicalContext {
	_ = parent
	_ = paramList
	_ = compiler
	return NewRewriteAsyncIteration_LexicalContext_4(parent, paramList, parent.function, compiler)
}

// Java method: LexicalContext.getFunctionDeclaringThisArgsSuper()
func (this_ *RewriteAsyncIteration_LexicalContext) GetFunctionDeclaringThisArgsSuper() pkg_rhino.Node {
	return this_.thisSuperArgsContext.ctx.function
}

// Java method: LexicalContext.mustReplaceThisSuperArgs()
func (this_ *RewriteAsyncIteration_LexicalContext) MustReplaceThisSuperArgs() bool {
	return ((this_.thisSuperArgsContext != nil) && this_.GetFunctionDeclaringThisArgsSuper().IsAsyncGeneratorFunction())
}

// Java class: com.google.javascript.jscomp.RewriteAsyncIteration.ThisSuperArgsContext
type RewriteAsyncIteration_ThisSuperArgsContext struct {
	ctx                 *RewriteAsyncIteration_LexicalContext `gson:"ctx"`
	usedSuperProperties j2gort.OrderedSet[pkg_rhino.Node]     `gson:"usedSuperProperties"`
	uniqueId            string                                `gson:"uniqueId"`
	thisNodeToAdd       pkg_rhino.Node                        `gson:"thisNodeToAdd"`
	idHash_             atomic.Int32
	usedArguments       bool `gson:"usedArguments"`
}

func (this_ *RewriteAsyncIteration_ThisSuperArgsContext) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteAsyncIteration_ThisSuperArgsContext) IsJavaNil_() bool { return this_ == nil }
func (*RewriteAsyncIteration_ThisSuperArgsContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteAsyncIteration$ThisSuperArgsContext", "ThisSuperArgsContext"
}
func (this_ *RewriteAsyncIteration_ThisSuperArgsContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteAsyncIteration_ThisSuperArgsContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteAsyncIteration.ThisSuperArgsContext", "com.google.javascript.jscomp.RewriteAsyncIteration$ThisSuperArgsContext", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteAsyncIteration_ThisSuperArgsContext](v_)
	})
}

func (this_ *RewriteAsyncIteration_ThisSuperArgsContext) initFields_RewriteAsyncIteration_ThisSuperArgsContext() {
	_ = this_
	this_.usedSuperProperties = j2gort.NewOrderedSet[pkg_rhino.Node]()

	this_.uniqueId = j2gort.NullString
}

// Java constructor: ThisSuperArgsContext(com.google.javascript.jscomp.RewriteAsyncIteration.LexicalContext ctx, java.lang.String uniqueId)
func NewRewriteAsyncIteration_ThisSuperArgsContext(ctx *RewriteAsyncIteration_LexicalContext, uniqueId string) *RewriteAsyncIteration_ThisSuperArgsContext {
	this_ := &RewriteAsyncIteration_ThisSuperArgsContext{}
	InitRewriteAsyncIteration_ThisSuperArgsContext(this_, ctx, uniqueId)
	return this_
}

func InitRewriteAsyncIteration_ThisSuperArgsContext(this_ *RewriteAsyncIteration_ThisSuperArgsContext, ctx *RewriteAsyncIteration_LexicalContext, uniqueId string) {
	_ = ctx
	_ = uniqueId
	this_.initFields_RewriteAsyncIteration_ThisSuperArgsContext()
	this_.ctx = ctx
	this_.uniqueId = uniqueId
}
func (*RewriteAsyncIteration) CanReleaseSyntacticScopes_() {}

/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/FunctionInjector.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.FunctionInjector
type FunctionInjector struct {
	compiler                 AbstractCompiler                                  `gson:"compiler"`
	knownConstantFunctions   j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"knownConstantFunctions"`
	safeNameIdSupplier       j2gort.Ext_java_util_function_Supplier            `gson:"safeNameIdSupplier"`
	throwawayNameSupplier    j2gort.Ext_java_util_function_Supplier            `gson:"throwawayNameSupplier"`
	functionArgumentInjector *FunctionArgumentInjector                         `gson:"functionArgumentInjector"`
	blockMutator             *FunctionToBlockMutator                           `gson:"blockMutator"`
	callerFnStateCache       j2gort.OrderedMap[pkg_rhino.Node, pkg_rhino.Node] `gson:"callerFnStateCache"`
	reusableDecomposer       *ExpressionDecomposer                             `gson:"reusableDecomposer"`
	foundInnerFn             pkg_rhino.Node                                    `gson:"foundInnerFn"`
	foundInnerFnCount        int32                                             `gson:"foundInnerFnCount"`
	idHash_                  atomic.Int32
	allowDecomposition       bool `gson:"allowDecomposition"`
	assumeStrictThis         bool `gson:"assumeStrictThis"`
}

func (this_ *FunctionInjector) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionInjector) IsJavaNil_() bool { return this_ == nil }
func (*FunctionInjector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionInjector", "FunctionInjector"
}
func (this_ *FunctionInjector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionInjector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionInjector", "com.google.javascript.jscomp.FunctionInjector", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionInjector](v_)
	})
}

var static_FunctionInjector_NO_FUNCTIONS pkg_rhino.Node

func StaticGet_FunctionInjector_NO_FUNCTIONS() pkg_rhino.Node {
	clinit_FunctionInjector()
	return static_FunctionInjector_NO_FUNCTIONS
}
func StaticSet_FunctionInjector_NO_FUNCTIONS(v pkg_rhino.Node) {
	clinit_FunctionInjector()
	static_FunctionInjector_NO_FUNCTIONS = v
}

var static_FunctionInjector_MULTIPLE_FUNCTIONS pkg_rhino.Node

func StaticGet_FunctionInjector_MULTIPLE_FUNCTIONS() pkg_rhino.Node {
	clinit_FunctionInjector()
	return static_FunctionInjector_MULTIPLE_FUNCTIONS
}
func StaticSet_FunctionInjector_MULTIPLE_FUNCTIONS(v pkg_rhino.Node) {
	clinit_FunctionInjector()
	static_FunctionInjector_MULTIPLE_FUNCTIONS = v
}

var static_FunctionInjector_REFERENCES_EVAL pkg_rhino.Node

func StaticGet_FunctionInjector_REFERENCES_EVAL() pkg_rhino.Node {
	clinit_FunctionInjector()
	return static_FunctionInjector_REFERENCES_EVAL
}
func StaticSet_FunctionInjector_REFERENCES_EVAL(v pkg_rhino.Node) {
	clinit_FunctionInjector()
	static_FunctionInjector_REFERENCES_EVAL = v
}

var static_FunctionInjector_MATCH_DECLARATION *NodeUtil_MatchDeclaration

func StaticGet_FunctionInjector_MATCH_DECLARATION() *NodeUtil_MatchDeclaration {
	clinit_FunctionInjector()
	return static_FunctionInjector_MATCH_DECLARATION
}
func StaticSet_FunctionInjector_MATCH_DECLARATION(v *NodeUtil_MatchDeclaration) {
	clinit_FunctionInjector()
	static_FunctionInjector_MATCH_DECLARATION = v
}

var static_FunctionInjector_MATCH_SHALLOW_STATEMENT *NodeUtil_MatchShallowStatement

func StaticGet_FunctionInjector_MATCH_SHALLOW_STATEMENT() *NodeUtil_MatchShallowStatement {
	clinit_FunctionInjector()
	return static_FunctionInjector_MATCH_SHALLOW_STATEMENT
}
func StaticSet_FunctionInjector_MATCH_SHALLOW_STATEMENT(v *NodeUtil_MatchShallowStatement) {
	clinit_FunctionInjector()
	static_FunctionInjector_MATCH_SHALLOW_STATEMENT = v
}

var static_FunctionInjector_NAME_COST_ESTIMATE int32

func StaticGet_FunctionInjector_NAME_COST_ESTIMATE() int32 {
	clinit_FunctionInjector()
	return static_FunctionInjector_NAME_COST_ESTIMATE
}
func StaticSet_FunctionInjector_NAME_COST_ESTIMATE(v int32) {
	clinit_FunctionInjector()
	static_FunctionInjector_NAME_COST_ESTIMATE = v
}
func StaticGet_FunctionInjector_COMMA_COST() int32 { return int32(1) }
func StaticGet_FunctionInjector_PAREN_COST() int32 { return int32(2) }

var clinitDone_FunctionInjector uint32
var clinitState_FunctionInjector j2gort.ClinitState

func clinit_FunctionInjector() {
	if atomic.LoadUint32(&clinitDone_FunctionInjector) != 0 {
		return
	}
	clinitSlow_FunctionInjector()
}

func clinitSlow_FunctionInjector() {
	if !j2gort.EnterClinit(&clinitState_FunctionInjector, &clinitDone_FunctionInjector) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FunctionInjector, &clinitDone_FunctionInjector)
	static_FunctionInjector_NO_FUNCTIONS = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_FUNCTION())
	static_FunctionInjector_MULTIPLE_FUNCTIONS = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_FUNCTION())
	static_FunctionInjector_REFERENCES_EVAL = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_FUNCTION())
	static_FunctionInjector_MATCH_DECLARATION = NewNodeUtil_MatchDeclaration()
	static_FunctionInjector_MATCH_SHALLOW_STATEMENT = NewNodeUtil_MatchShallowStatement()
	static_FunctionInjector_NAME_COST_ESTIMATE = StaticGet_InlineCostEstimator_ESTIMATED_IDENTIFIER_COST()
}

func (this_ *FunctionInjector) initFields_FunctionInjector() {
	_ = this_
	this_.knownConstantFunctions = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	this_.throwawayNameSupplier = NewFunctionInjector__Anon_45(this_)
	this_.callerFnStateCache = j2gort.NewOrderedMap[pkg_rhino.Node, pkg_rhino.Node](

	// Java constructor: FunctionInjector(com.google.javascript.jscomp.FunctionInjector.Builder builder)
	)
}

func NewFunctionInjector(builder *FunctionInjector_Builder) *FunctionInjector {
	clinit_FunctionInjector()
	this_ := &FunctionInjector{}
	InitFunctionInjector(this_, builder)
	return this_
}

func InitFunctionInjector(this_ *FunctionInjector, builder *FunctionInjector_Builder) {
	_ = builder
	this_.initFields_FunctionInjector()
	this_.compiler = j2gort.CheckNotNullIface(builder.compiler)
	this_.safeNameIdSupplier = j2gort.CheckNotNull(builder.safeNameIdSupplier)
	this_.assumeStrictThis = builder.assumeStrictThis
	this_.allowDecomposition = builder.allowDecomposition
	this_.functionArgumentInjector = j2gort.CheckNotNullPtr(builder.functionArgumentInjector)
	this_.blockMutator = NewFunctionToBlockMutator_3(this_.compiler, this_.safeNameIdSupplier, this_.functionArgumentInjector)
}

// Java method: FunctionInjector.doesFunctionMeetMinimumRequirements(java.lang.String fnName, com.google.javascript.rhino.Node fnNode)
func (this_ *FunctionInjector) DoesFunctionMeetMinimumRequirements(fnName string, fnNode pkg_rhino.Node) bool {
	_ = fnName
	_ = fnNode
	var block pkg_rhino.Node = NodeUtil_GetFunctionBody(fnNode)
	_ = block
	var fnRecursionName string = fnNode.GetFirstChild().GetString()
	_ = fnRecursionName
	if !(fnRecursionName != j2gort.NullString) {
		j2gort.CheckState(false)
	}
	if NodeUtil_IsNameReferenced_3(block, "arguments", StaticGet_NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION()) {
		return false
	}
	return !(FunctionInjector_ContainsInjectionBlocker(block, fnName, fnRecursionName))
}

// Java static method: FunctionInjector.containsInjectionBlocker(com.google.javascript.rhino.Node n, java.lang.String fnName, java.lang.String fnRecursionName)
func FunctionInjector_ContainsInjectionBlocker(n pkg_rhino.Node, fnName string, fnRecursionName string) bool {
	clinit_FunctionInjector()
	_ = n
	_ = fnName
	_ = fnRecursionName
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		if ((j2gort.CheckNotNullString(name) == "eval") || (!(j2gort.CheckNotNullString(fnName) == "") && (j2gort.CheckNotNullString(name) == fnName))) || (!(j2gort.CheckNotNullString(fnRecursionName) == "") && (j2gort.CheckNotNullString(name) == fnRecursionName)) {
			return true
		}
	} else if n.IsSuper() || (n.RawToken_() == uint8(135)) {
		return true
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if FunctionInjector_ContainsInjectionBlocker(c, fnName, fnRecursionName) {
				return true
			}
		}
	}
	return false
}

// Java method: FunctionInjector.canInlineReferenceToFunction(com.google.javascript.jscomp.FunctionInjector.Reference ref, com.google.javascript.rhino.Node fnNode, com.google.common.collect.ImmutableSet<java.lang.String> needAliases, boolean referencesThis, boolean containsClosures)
func (this_ *FunctionInjector) CanInlineReferenceToFunction(ref FunctionInjector_Reference, fnNode pkg_rhino.Node, needAliases j2gort.Ext_com_google_common_collect_ImmutableSet, referencesThis bool, containsClosures bool) FunctionInjector_CanInlineResult {
	_ = ref
	_ = fnNode
	_ = needAliases
	_ = referencesThis
	_ = containsClosures
	base_ref_ := ref.GetFunctionInjector_ReferenceBase()
	var callNode pkg_rhino.Node = base_ref_.callNode
	_ = callNode
	if !(this_.priv_FunctionInjector_isSupportedCallType_2938(callNode, fnNode)) {
		return StaticGet_FunctionInjector_CanInlineResult_NO()
	}
	if FunctionInjector_HasSpreadCallArgument(callNode) {
		return StaticGet_FunctionInjector_CanInlineResult_NO()
	}
	if containsClosures {
		if !(base_ref_.scope.IsGlobal()) {
			return StaticGet_FunctionInjector_CanInlineResult_NO()
		} else if NodeUtil_IsWithinLoop(callNode) {
			return StaticGet_FunctionInjector_CanInlineResult_NO()
		}
	}
	if referencesThis && !(NodeUtil_IsFunctionObjectCall(callNode)) {
		return StaticGet_FunctionInjector_CanInlineResult_NO()
	}
	if base_ref_.mode == StaticGet_FunctionInjector_InliningMode_DIRECT() {
		return this_.priv_FunctionInjector_canInlineReferenceDirectly_2951(ref, fnNode, needAliases)
	} else {
		return this_.priv_FunctionInjector_canInlineReferenceAsStatementBlock_2947(ref, fnNode, needAliases)
	}
}

// Java method: FunctionInjector.isSupportedCallType(com.google.javascript.rhino.Node callNode, com.google.javascript.rhino.Node fnNode)
func (this_ *FunctionInjector) priv_FunctionInjector_isSupportedCallType_2938(callNode pkg_rhino.Node, fnNode pkg_rhino.Node) bool {
	_ = callNode
	_ = fnNode
	if !(callNode.GetFirstChild().IsName()) {
		if NodeUtil_IsFunctionObjectCall(callNode) {
			if !(this_.assumeStrictThis) && !(fnNode.IsArrowFunction()) {
				var thisValue pkg_rhino.Node = callNode.GetSecondChild()
				_ = thisValue
				if (thisValue == 0) || !(thisValue.IsThis()) {
					return false
				}
			}
		} else if NodeUtil_IsFunctionObjectApply(callNode) {
			return false
		}
	}
	return true
}

// Java static method: FunctionInjector.hasSpreadCallArgument(com.google.javascript.rhino.Node callNode)
func FunctionInjector_HasSpreadCallArgument(callNode pkg_rhino.Node) bool {
	clinit_FunctionInjector()
	_ = callNode
	if !(NodeUtil_IsNormalOrOptChainCall(callNode)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(callNode))
	}
	{
		var arg pkg_rhino.Node = callNode.GetSecondChild()
		_ = arg
		for ; arg != 0; arg = arg.GetNext() {
			if arg.IsSpread() {
				return true
			}
		}
	}
	return false
}

// Java method: FunctionInjector.inline(com.google.javascript.jscomp.FunctionInjector.Reference ref, java.lang.String fnName, com.google.javascript.rhino.Node fnNode)
func (this_ *FunctionInjector) Inline(ref FunctionInjector_Reference, fnName string, fnNode pkg_rhino.Node) pkg_rhino.Node {
	_ = ref
	_ = fnName
	_ = fnNode
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	return this_.priv_FunctionInjector_internalInline_2941(ref, fnName, fnNode)
}

// Java method: FunctionInjector.unsafeInline(com.google.javascript.jscomp.FunctionInjector.Reference ref, java.lang.String fnName, com.google.javascript.rhino.Node fnNode)
func (this_ *FunctionInjector) UnsafeInline(ref FunctionInjector_Reference, fnName string, fnNode pkg_rhino.Node) pkg_rhino.Node {
	_ = ref
	_ = fnName
	_ = fnNode
	return this_.priv_FunctionInjector_internalInline_2941(ref, fnName, fnNode)
}

// Java method: FunctionInjector.internalInline(com.google.javascript.jscomp.FunctionInjector.Reference ref, java.lang.String fnName, com.google.javascript.rhino.Node fnNode)
func (this_ *FunctionInjector) priv_FunctionInjector_internalInline_2941(ref FunctionInjector_Reference, fnName string, fnNode pkg_rhino.Node) pkg_rhino.Node {
	_ = ref
	_ = fnName
	_ = fnNode
	var result pkg_rhino.Node = pkg_rhino.Node(0)
	_ = result
	if ref.GetFunctionInjector_ReferenceBase().mode == StaticGet_FunctionInjector_InliningMode_DIRECT() {
		result = this_.priv_FunctionInjector_inlineReturnValue_2942(ref, fnNode)
	} else {
		result = this_.priv_FunctionInjector_inlineFunction_2946(ref, fnNode, fnName)
	}
	this_.compiler.ReportChangeToEnclosingScope(result)
	return result
}

// Java method: FunctionInjector.inlineReturnValue(com.google.javascript.jscomp.FunctionInjector.Reference ref, com.google.javascript.rhino.Node fnNode)
func (this_ *FunctionInjector) priv_FunctionInjector_inlineReturnValue_2942(ref FunctionInjector_Reference, fnNode pkg_rhino.Node) pkg_rhino.Node {
	_ = ref
	_ = fnNode
	var callNode pkg_rhino.Node = ref.GetFunctionInjector_ReferenceBase().callNode
	_ = callNode
	var block pkg_rhino.Node = fnNode.GetLastChild()
	_ = block
	var newExpression pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newExpression
	var paramToArgMap *FunctionArgumentInjector_ParamArgList = this_.functionArgumentInjector.AcquireFunctionCallParameterMap(fnNode, callNode, this_.safeNameIdSupplier)
	_ = paramToArgMap
	func() {
		defer func() {
			this_.functionArgumentInjector.ReleaseParamArgList(paramToArgMap)
		}()
		if !(block.HasChildren()) {
			var srcLocation pkg_rhino.Node = block
			_ = srcLocation
			newExpression = NodeUtil_NewUndefinedNode(srcLocation)
		} else {
			var returnNode pkg_rhino.Node = block.GetFirstChild()
			_ = returnNode
			if !(returnNode.IsReturn()) {
				j2gort.CheckArgumentMsg(false, j2gort.BoxNode(returnNode))
			}
			var safeReturnNode pkg_rhino.Node = returnNode.CloneTree_0()
			_ = safeReturnNode
			var inlineResult pkg_rhino.Node = this_.functionArgumentInjector.Inject_4_AbstractCompiler_Node_Node_FunctionArgumentInjector_ParamArgList(nil, safeReturnNode, pkg_rhino.Node(0), paramToArgMap)
			_ = inlineResult
			if !(safeReturnNode == inlineResult) {
				j2gort.CheckArgument(false)
			}
			newExpression = safeReturnNode.RemoveFirstChild()
			NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(newExpression, this_.compiler)
		}
	}()
	var typeBeforeCast pkg_rhino.JSType = callNode.GetJSTypeBeforeCast()
	_ = typeBeforeCast
	if typeBeforeCast != nil && !typeBeforeCast.IsJavaNil_() {
		newExpression.SetJSTypeBeforeCast(typeBeforeCast)
		newExpression.SetJSType(callNode.GetJSType())
	}
	if (callNode.GetColor() != nil) && callNode.IsColorFromTypeCast() {
		newExpression.SetColor(callNode.GetColor())
		newExpression.SetColorFromTypeCast()
	}
	callNode.ReplaceWith(newExpression)
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(callNode, this_.compiler)
	return newExpression
}

// Java method: FunctionInjector.classifyCallSite(com.google.javascript.jscomp.FunctionInjector.Reference ref)
func (this_ *FunctionInjector) priv_FunctionInjector_classifyCallSite_2943(ref FunctionInjector_Reference) *FunctionInjector_CallSiteType {
	_ = ref
	base_ref_ := ref.GetFunctionInjector_ReferenceBase()
	var callNode pkg_rhino.Node = base_ref_.callNode
	_ = callNode
	var parent pkg_rhino.Node = callNode.GetParent()
	_ = parent
	var grandParent pkg_rhino.Node = parent.GetParent()
	_ = grandParent
	if NodeUtil_IsExprCall(parent) {
		return StaticGet_FunctionInjector_CallSiteType_SIMPLE_CALL()
	} else if ((NodeUtil_IsExprAssign(grandParent) && !(NodeUtil_IsNameDeclOrSimpleAssignLhs(callNode, parent))) && parent.GetFirstChild().IsName()) && !(NodeUtil_IsConstantName(parent.GetFirstChild())) {
		return StaticGet_FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT()
	} else if ((parent.IsName() && !(NodeUtil_IsConstantName(parent))) && grandParent.IsVar()) && grandParent.HasOneChild() {
		return StaticGet_FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT()
	} else {
		var decomposer *ExpressionDecomposer = this_.priv_FunctionInjector_getDecomposer_2944(base_ref_.scope)
		_ = decomposer
		switch (decomposer.CanExposeExpression(callNode)).Ordinal() {
		case int32(1) /* DecompositionType.MOVABLE */ :
			return StaticGet_FunctionInjector_CallSiteType_EXPRESSION()
		case int32(2) /* DecompositionType.DECOMPOSABLE */ :
			return StaticGet_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION()
		case int32(0) /* DecompositionType.UNDECOMPOSABLE */ :
		}
	}
	return StaticGet_FunctionInjector_CallSiteType_UNSUPPORTED()
}

// Java method: FunctionInjector.getDecomposer(com.google.javascript.jscomp.Scope scope)
func (this_ *FunctionInjector) priv_FunctionInjector_getDecomposer_2944(scope *Scope) *ExpressionDecomposer {
	_ = scope
	var decomposer *ExpressionDecomposer = this_.reusableDecomposer
	_ = decomposer
	if decomposer == nil {
		decomposer = this_.compiler.CreateExpressionDecomposer(this_.safeNameIdSupplier, this_.knownConstantFunctions, scope)
		this_.reusableDecomposer = decomposer
	} else {
		decomposer.Reset(scope, this_.knownConstantFunctions)
	}
	return decomposer
}

// Java method: FunctionInjector.maybePrepareCall(com.google.javascript.jscomp.FunctionInjector.Reference ref)
func (this_ *FunctionInjector) MaybePrepareCall(ref FunctionInjector_Reference) {
	_ = ref
	var callSiteType *FunctionInjector_CallSiteType = this_.priv_FunctionInjector_classifyCallSite_2943(ref)
	_ = callSiteType
	callSiteType.Prepare(this_, ref)
}

// Java method: FunctionInjector.inlineFunction(com.google.javascript.jscomp.FunctionInjector.Reference ref, com.google.javascript.rhino.Node fnNode, java.lang.String fnName)
func (this_ *FunctionInjector) priv_FunctionInjector_inlineFunction_2946(ref FunctionInjector_Reference, fnNode pkg_rhino.Node, fnName string) pkg_rhino.Node {
	_ = ref
	_ = fnNode
	_ = fnName
	base_ref_ := ref.GetFunctionInjector_ReferenceBase()
	var callNode pkg_rhino.Node = base_ref_.callNode
	_ = callNode
	var parent pkg_rhino.Node = callNode.GetParent()
	_ = parent
	var grandParent pkg_rhino.Node = parent.GetParent()
	_ = grandParent
	var callSiteType *FunctionInjector_CallSiteType = this_.priv_FunctionInjector_classifyCallSite_2943(ref)
	_ = callSiteType
	if !(callSiteType != StaticGet_FunctionInjector_CallSiteType_UNSUPPORTED()) {
		j2gort.CheckArgument(false)
	}
	var resultName string = j2gort.NullString
	_ = resultName
	var needsDefaultReturnResult bool = true
	_ = needsDefaultReturnResult
	switch j2gort.CheckNotNullPtr(callSiteType).Ordinal() {
	case int32(2) /* CallSiteType.SIMPLE_ASSIGNMENT */ :
		resultName = parent.GetFirstChild().GetString()
		FunctionInjector_RemoveConstantVarAnnotation(base_ref_.scope, resultName)
	case int32(3) /* CallSiteType.VAR_DECL_SIMPLE_ASSIGNMENT */ :
		resultName = parent.GetString()
		FunctionInjector_RemoveConstantVarAnnotation(base_ref_.scope, resultName)
	case int32(1) /* CallSiteType.SIMPLE_CALL */ :
		resultName = j2gort.NullString
		needsDefaultReturnResult = false
	case int32(4) /* CallSiteType.EXPRESSION */ :
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Movable expressions must be moved before inlining.")))
	case int32(5) /* CallSiteType.DECOMPOSABLE_EXPRESSION */ :
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Decomposable expressions must be decomposed before inlining.")))
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected call site type.")))
	}
	var isCallInLoop bool = NodeUtil_IsWithinLoop(callNode)
	_ = isCallInLoop
	var newBlock pkg_rhino.Node = this_.blockMutator.Mutate(fnName, fnNode, callNode, resultName, needsDefaultReturnResult, isCallInLoop)
	_ = newBlock
	NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(newBlock, this_.compiler)
	switch j2gort.CheckNotNullPtr(callSiteType).Ordinal() {
	case int32(3) /* CallSiteType.VAR_DECL_SIMPLE_ASSIGNMENT */ :
		var firstChild pkg_rhino.Node = parent.RemoveFirstChild()
		_ = firstChild
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(firstChild, this_.compiler)
		if !(!(parent.HasChildren())) {
			j2gort.CheckState(false)
		}
		newBlock.InsertAfter(grandParent)
	case int32(2) /* CallSiteType.SIMPLE_ASSIGNMENT */ :
		if !(grandParent.IsExprResult()) {
			j2gort.CheckState(false)
		}
		grandParent.ReplaceWith(newBlock)
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(grandParent, this_.compiler)
	case int32(1) /* CallSiteType.SIMPLE_CALL */ :
		if !(parent.IsExprResult()) {
			j2gort.CheckState(false)
		}
		parent.ReplaceWith(newBlock)
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(parent, this_.compiler)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected call site type.")))
	}
	return newBlock
}

// Java static method: FunctionInjector.removeConstantVarAnnotation(com.google.javascript.jscomp.Scope scope, java.lang.String name)
func FunctionInjector_RemoveConstantVarAnnotation(scope *Scope, name string) {
	clinit_FunctionInjector()
	_ = scope
	_ = name
	var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(name))
	_ = var_
	var nameNode pkg_rhino.Node
	_ = nameNode
	if var_ == nil {
		nameNode = pkg_rhino.Node(0)
	} else {
		nameNode = var_.GetNameNode()
	}
	if nameNode == 0 {
		return
	}
	if nameNode.IsDeclaredConstantVar() {
		nameNode.SetDeclaredConstantVar(false)
	}
}

// Java static method: FunctionInjector.isDirectCallNodeReplacementPossible(com.google.javascript.rhino.Node fnNode)
func FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode pkg_rhino.Node) bool {
	clinit_FunctionInjector()
	_ = fnNode
	var block pkg_rhino.Node = NodeUtil_GetFunctionBody(fnNode)
	_ = block
	if !(block.HasChildren()) {
		return true
	} else if block.HasOneChild() {
		if block.GetFirstChild().IsReturn() && (block.GetFirstFirstChild() != 0) {
			return true
		}
	}
	return false
}

// Java method: FunctionInjector.canInlineReferenceAsStatementBlock(com.google.javascript.jscomp.FunctionInjector.Reference ref, com.google.javascript.rhino.Node fnNode, com.google.common.collect.ImmutableSet<java.lang.String> namesToAlias)
func (this_ *FunctionInjector) priv_FunctionInjector_canInlineReferenceAsStatementBlock_2947(ref FunctionInjector_Reference, fnNode pkg_rhino.Node, namesToAlias j2gort.Ext_com_google_common_collect_ImmutableSet) FunctionInjector_CanInlineResult {
	_ = ref
	_ = fnNode
	_ = namesToAlias
	var callSiteType *FunctionInjector_CallSiteType = this_.priv_FunctionInjector_classifyCallSite_2943(ref)
	_ = callSiteType
	if callSiteType == StaticGet_FunctionInjector_CallSiteType_UNSUPPORTED() {
		return StaticGet_FunctionInjector_CanInlineResult_NO()
	}
	if !(this_.allowDecomposition) && ((callSiteType == StaticGet_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION()) || (callSiteType == StaticGet_FunctionInjector_CallSiteType_EXPRESSION())) {
		return StaticGet_FunctionInjector_CanInlineResult_NO()
	}
	if !(this_.priv_FunctionInjector_callMeetsBlockInliningRequirements_2950(ref, fnNode, namesToAlias)) {
		return StaticGet_FunctionInjector_CanInlineResult_NO()
	}
	if (callSiteType == StaticGet_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION()) || (callSiteType == StaticGet_FunctionInjector_CallSiteType_EXPRESSION()) {
		return StaticGet_FunctionInjector_CanInlineResult_AFTER_PREPARATION()
	} else {
		return StaticGet_FunctionInjector_CanInlineResult_YES()
	}
}

// Java static method: FunctionInjector.containsShallowEval(com.google.javascript.rhino.Node n)
func FunctionInjector_ContainsShallowEval(n pkg_rhino.Node) bool {
	clinit_FunctionInjector()
	_ = n
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if c.IsName() && (j2gort.CheckNotNullString(c.GetString()) == "eval") {
				return true
			}
			if !(c.IsFunction()) && FunctionInjector_ContainsShallowEval(c) {
				return true
			}
		}
	}
	return false
}

// Java method: FunctionInjector.scanInnerFunctions(com.google.javascript.rhino.Node n)
func (this_ *FunctionInjector) priv_FunctionInjector_scanInnerFunctions_2948(n pkg_rhino.Node) {
	_ = n
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(c) {
				if this_.foundInnerFnCount == int32(0) {
					this_.foundInnerFn = c
					this_.foundInnerFnCount = int32(1)
				} else {
					this_.foundInnerFnCount = int32(2)
					return
				}
			}
			this_.priv_FunctionInjector_scanInnerFunctions_2948(c)
			if this_.foundInnerFnCount > int32(1) {
				return
			}
		}
	}
}

// Java method: FunctionInjector.getCallerFnState(com.google.javascript.rhino.Node callerFn)
func (this_ *FunctionInjector) priv_FunctionInjector_getCallerFnState_2949(callerFn pkg_rhino.Node) pkg_rhino.Node {
	_ = callerFn
	if !(callerFn.IsFunction()) {
		j2gort.CheckState(false)
	}
	var cached pkg_rhino.Node = this_.callerFnStateCache.Get(callerFn)
	_ = cached
	if cached != 0 {
		return cached
	}
	if FunctionInjector_ContainsShallowEval(callerFn) {
		cached = StaticGet_FunctionInjector_REFERENCES_EVAL()
	} else {
		this_.foundInnerFn = pkg_rhino.Node(0)
		this_.foundInnerFnCount = int32(0)
		this_.priv_FunctionInjector_scanInnerFunctions_2948(callerFn)
	iifeSwitch_1:
		switch {
		default:

			switch this_.foundInnerFnCount {
			case int32(0):
				cached = StaticGet_FunctionInjector_NO_FUNCTIONS()
				break iifeSwitch_1
			case int32(1):
				cached = this_.foundInnerFn
				break iifeSwitch_1
			default:
				cached = StaticGet_FunctionInjector_MULTIPLE_FUNCTIONS()
				break iifeSwitch_1
			}
			panic("unreachable switch expression")
		}

		this_.foundInnerFn = pkg_rhino.Node(0)
	}
	if this_.callerFnStateCache.IsEmpty() {
		this_.callerFnStateCache = j2gort.AcquireOrderedMapCap[pkg_rhino.Node, pkg_rhino.Node](int32(2048))
	}

	this_.callerFnStateCache.Put(callerFn, cached)
	return cached
}

// Java method: FunctionInjector.callMeetsBlockInliningRequirements(com.google.javascript.jscomp.FunctionInjector.Reference callRef, com.google.javascript.rhino.Node calleeFn, com.google.common.collect.ImmutableSet<java.lang.String> namesToAlias)
func (this_ *FunctionInjector) priv_FunctionInjector_callMeetsBlockInliningRequirements_2950(callRef FunctionInjector_Reference, calleeFn pkg_rhino.Node, namesToAlias j2gort.Ext_com_google_common_collect_ImmutableSet) bool {
	_ = callRef
	_ = calleeFn
	_ = namesToAlias
	var calleeContainsVars bool = NodeUtil_Has(NodeUtil_GetFunctionBody(calleeFn), j2gort.ToIface_Ext_com_google_common_base_Predicate(StaticGet_FunctionInjector_MATCH_DECLARATION()), j2gort.ToIface_Ext_com_google_common_base_Predicate(StaticGet_FunctionInjector_MATCH_SHALLOW_STATEMENT()))
	_ = calleeContainsVars
	var forbidTemps bool = false
	_ = forbidTemps
	if !(j2gort.CastPtr[Scope](callRef.GetFunctionInjector_ReferenceBase().scope.GetClosestHoistScope()).IsGlobal()) {
		var callerFn pkg_rhino.Node = j2gort.CastPtr[Scope](callRef.GetFunctionInjector_ReferenceBase().scope.GetClosestHoistScope()).GetRootNode_0().GetParent()
		_ = callerFn
		var callerState pkg_rhino.Node = this_.priv_FunctionInjector_getCallerFnState_2949(callerFn)
		_ = callerState
		forbidTemps = ((callerState == StaticGet_FunctionInjector_REFERENCES_EVAL()) || !((callerState == StaticGet_FunctionInjector_NO_FUNCTIONS()) || (callerState == calleeFn)))
	}
	if calleeContainsVars && forbidTemps {
		return false
	}
	if forbidTemps {
		var args *FunctionArgumentInjector_ParamArgList = this_.functionArgumentInjector.AcquireFunctionCallParameterMap(calleeFn, callRef.GetFunctionInjector_ReferenceBase().callNode, this_.safeNameIdSupplier)
		_ = args
		var didRet_1 bool
		var retVal_2 bool
		func() {
			defer func() {
				this_.functionArgumentInjector.ReleaseParamArgList(args)
			}()
			var hasArgs bool = !(args.IsEmpty())
			_ = hasArgs
			if hasArgs {
				if !(namesToAlias.IsEmpty()) {
					retVal_2 = false
					didRet_1 = true
					return
				}
				var allTemps j2gort.Ext_com_google_common_collect_ImmutableSet = this_.functionArgumentInjector.GatherCallArgumentsNeedingTemps(this_.compiler, calleeFn, args, namesToAlias, this_.compiler.GetCodingConvention())
				_ = allTemps
				if !(allTemps.IsEmpty()) {
					retVal_2 = false
					didRet_1 = true
					return
				}
			}
		}()
		if didRet_1 {
			return retVal_2
		}
	}
	return true
}

// Java method: FunctionInjector.canInlineReferenceDirectly(com.google.javascript.jscomp.FunctionInjector.Reference ref, com.google.javascript.rhino.Node fnNode, com.google.common.collect.ImmutableSet<java.lang.String> namesToAlias)
func (this_ *FunctionInjector) priv_FunctionInjector_canInlineReferenceDirectly_2951(ref FunctionInjector_Reference, fnNode pkg_rhino.Node, namesToAlias j2gort.Ext_com_google_common_collect_ImmutableSet) FunctionInjector_CanInlineResult {
	_ = ref
	_ = fnNode
	_ = namesToAlias
	if !(FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode)) {
		return StaticGet_FunctionInjector_CanInlineResult_NO()
	}
	var callNode pkg_rhino.Node = ref.GetFunctionInjector_ReferenceBase().callNode
	_ = callNode
	var cArg pkg_rhino.Node = callNode.GetSecondChild()
	_ = cArg
	if !(callNode.GetFirstChild().IsName()) {
		if NodeUtil_IsFunctionObjectCall(callNode) {
			if !(fnNode.IsArrowFunction()) && ((cArg == 0) || !(cArg.IsThis())) {
				return StaticGet_FunctionInjector_CanInlineResult_NO()
			}
			if cArg != 0 {
				cArg = cArg.GetNext()
			}
		} else {
			if !(!(NodeUtil_IsFunctionObjectApply(callNode))) {
				j2gort.CheckState(false)
			}
		}
	}
	if !(namesToAlias.IsEmpty()) {
		return StaticGet_FunctionInjector_CanInlineResult_NO()
	}
	var args *FunctionArgumentInjector_ParamArgList = this_.functionArgumentInjector.AcquireFunctionCallParameterMap(fnNode, callNode, this_.throwawayNameSupplier)
	_ = args
	var didRet_3 bool
	var retVal_4 FunctionInjector_CanInlineResult
	func() {
		defer func() {
			this_.functionArgumentInjector.ReleaseParamArgList(args)
		}()
		var hasArgs bool = !(args.IsEmpty())
		_ = hasArgs
		if hasArgs {
			var allTemps j2gort.Ext_com_google_common_collect_ImmutableSet = this_.functionArgumentInjector.GatherCallArgumentsNeedingTemps(this_.compiler, fnNode, args, namesToAlias, this_.compiler.GetCodingConvention())
			_ = allTemps
			if !(allTemps.IsEmpty()) {
				retVal_4 = StaticGet_FunctionInjector_CanInlineResult_NO()
				didRet_3 = true
				return
			}
		}
	}()
	if didRet_3 {
		return retVal_4
	}
	return StaticGet_FunctionInjector_CanInlineResult_YES()
}

// Java method: FunctionInjector.inliningLowersCost(com.google.javascript.jscomp.JSChunk fnChunk, com.google.javascript.rhino.Node fnNode, java.util.Collection<? extends com.google.javascript.jscomp.FunctionInjector.Reference> refs, java.util.Set<java.lang.String> namesToAlias, boolean isRemovable, boolean referencesThis)
func (this_ *FunctionInjector) InliningLowersCost(fnChunk *JSChunk, fnNode pkg_rhino.Node, refs j2gort.Ext_java_util_Collection, namesToAlias j2gort.Ext_java_util_Set, isRemovable bool, referencesThis bool) bool {
	_ = fnChunk
	_ = fnNode
	_ = refs
	_ = namesToAlias
	_ = isRemovable
	_ = referencesThis
	var referenceCount int32 = refs.Size()
	_ = referenceCount
	if referenceCount == int32(0) {
		return true
	}
	var referencesUsingBlockInlining int32 = int32(0)
	_ = referencesUsingBlockInlining
	var checkChunks bool = (isRemovable && (fnChunk != nil))
	_ = checkChunks
	var chunkGraph *JSChunkGraph = this_.compiler.GetChunkGraph()
	_ = chunkGraph
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(refs))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			ref := j2gort.AsIface[FunctionInjector_Reference](rawElem_6)
			_ = ref
			if ref.GetFunctionInjector_ReferenceBase().mode == StaticGet_FunctionInjector_InliningMode_BLOCK() {
				referencesUsingBlockInlining = referencesUsingBlockInlining + 1
			}
			if checkChunks && (ref.GetFunctionInjector_ReferenceBase().chunk != nil) {
				if (ref.GetFunctionInjector_ReferenceBase().chunk != fnChunk) && !(chunkGraph.DependsOn(ref.GetFunctionInjector_ReferenceBase().chunk, fnChunk)) {
					isRemovable = false
					checkChunks = false
				}
			}
		}
	}
	var referencesUsingDirectInlining int32 = (referenceCount - referencesUsingBlockInlining)
	_ = referencesUsingDirectInlining
	return this_.priv_FunctionInjector_inliningLowersCostForCounts_2954(fnNode, referenceCount, referencesUsingDirectInlining, referencesUsingBlockInlining, namesToAlias, isRemovable, referencesThis)
}

// Java method: FunctionInjector.inliningLowersCostInline(com.google.javascript.jscomp.JSChunk fnChunk, com.google.javascript.rhino.Node fnNode, com.google.javascript.jscomp.FunctionInjector.@org.jspecify.annotations.Nullable Reference ref0, com.google.javascript.jscomp.FunctionInjector.@org.jspecify.annotations.Nullable Reference ref1, java.util.Set<java.lang.String> namesToAlias, boolean isRemovable, boolean referencesThis)
func (this_ *FunctionInjector) InliningLowersCostInline(fnChunk *JSChunk, fnNode pkg_rhino.Node, ref0 FunctionInjector_Reference, ref1 FunctionInjector_Reference, namesToAlias j2gort.Ext_java_util_Set, isRemovable bool, referencesThis bool) bool {
	_ = fnChunk
	_ = fnNode
	_ = ref0
	_ = ref1
	_ = namesToAlias
	_ = isRemovable
	_ = referencesThis
	if ref0 == nil || ref0.IsJavaNil_() {
		return true
	}
	base_ref0_ := ref0.GetFunctionInjector_ReferenceBase()
	var referenceCount int32
	_ = referenceCount
	if ref1 != nil && !ref1.IsJavaNil_() {
		referenceCount = int32(2)
	} else {
		referenceCount = int32(1)
	}
	var referencesUsingBlockInlining int32 = int32(0)
	_ = referencesUsingBlockInlining
	if base_ref0_.mode == StaticGet_FunctionInjector_InliningMode_BLOCK() {
		referencesUsingBlockInlining = referencesUsingBlockInlining + 1
	}
	if (ref1 != nil && !ref1.IsJavaNil_()) && (ref1.GetFunctionInjector_ReferenceBase().mode == StaticGet_FunctionInjector_InliningMode_BLOCK()) {
		referencesUsingBlockInlining = referencesUsingBlockInlining + 1
	}
	if isRemovable && (fnChunk != nil) {
		var chunkGraph *JSChunkGraph = this_.compiler.GetChunkGraph()
		_ = chunkGraph
		if ((base_ref0_.chunk != nil) && (base_ref0_.chunk != fnChunk)) && !(chunkGraph.DependsOn(base_ref0_.chunk, fnChunk)) {
			isRemovable = false
		} else if (((ref1 != nil && !ref1.IsJavaNil_()) && (ref1.GetFunctionInjector_ReferenceBase().chunk != nil)) && (ref1.GetFunctionInjector_ReferenceBase().chunk != fnChunk)) && !(chunkGraph.DependsOn(ref1.GetFunctionInjector_ReferenceBase().chunk, fnChunk)) {
			isRemovable = false
		}
	}
	var referencesUsingDirectInlining int32 = (referenceCount - referencesUsingBlockInlining)
	_ = referencesUsingDirectInlining
	return this_.priv_FunctionInjector_inliningLowersCostForCounts_2954(fnNode, referenceCount, referencesUsingDirectInlining, referencesUsingBlockInlining, namesToAlias, isRemovable, referencesThis)
}

// Java method: FunctionInjector.inliningLowersCostForCounts(com.google.javascript.rhino.Node fnNode, int referenceCount, int referencesUsingDirectInlining, int referencesUsingBlockInlining, java.util.Set<java.lang.String> namesToAlias, boolean isRemovable, boolean referencesThis)
func (this_ *FunctionInjector) priv_FunctionInjector_inliningLowersCostForCounts_2954(fnNode pkg_rhino.Node, referenceCount int32, referencesUsingDirectInlining int32, referencesUsingBlockInlining int32, namesToAlias j2gort.Ext_java_util_Set, isRemovable bool, referencesThis bool) bool {
	_ = fnNode
	_ = referenceCount
	_ = referencesUsingDirectInlining
	_ = referencesUsingBlockInlining
	_ = namesToAlias
	_ = isRemovable
	_ = referencesThis
	if ((referenceCount == int32(1)) && isRemovable) && (referencesUsingDirectInlining == int32(1)) {
		return true
	}
	var callCost int32 = FunctionInjector_EstimateCallCost(fnNode, referencesThis)
	_ = callCost
	var overallCallCost int32 = (callCost * referenceCount)
	_ = overallCallCost
	var costDeltaDirect int32 = FunctionInjector_InlineCostDelta(fnNode, namesToAlias, StaticGet_FunctionInjector_InliningMode_DIRECT())
	_ = costDeltaDirect
	var costDeltaBlock int32 = FunctionInjector_InlineCostDelta(fnNode, namesToAlias, StaticGet_FunctionInjector_InliningMode_BLOCK())
	_ = costDeltaBlock
	return FunctionInjector_DoesLowerCost(fnNode, overallCallCost, referencesUsingDirectInlining, costDeltaDirect, referencesUsingBlockInlining, costDeltaBlock, isRemovable)
}

// Java static method: FunctionInjector.doesLowerCost(com.google.javascript.rhino.Node fnNode, int callCost, int directInlines, int costDeltaDirect, int blockInlines, int costDeltaBlock, boolean removable)
func FunctionInjector_DoesLowerCost(fnNode pkg_rhino.Node, callCost int32, directInlines int32, costDeltaDirect int32, blockInlines int32, costDeltaBlock int32, removable bool) bool {
	clinit_FunctionInjector()
	_ = fnNode
	_ = callCost
	_ = directInlines
	_ = costDeltaDirect
	_ = blockInlines
	_ = costDeltaBlock
	_ = removable
	var fnInstanceCount int32 = ((directInlines + blockInlines) - (func() int32 {
		if removable {
			return int32(1)
		}
		return int32(0)
	}()))
	_ = fnInstanceCount
	if fnInstanceCount == int32(0) {
		return ((blockInlines <= int32(0)) || (costDeltaBlock <= int32(0)))
	}
	var costDelta int32 = ((directInlines * -(costDeltaDirect)) + (blockInlines * -(costDeltaBlock)))
	_ = costDelta
	var threshold int32 = ((callCost + costDelta) / fnInstanceCount)
	_ = threshold
	return (InlineCostEstimator_GetCost_2(fnNode, (threshold+int32(1))) <= threshold)
}

// Java static method: FunctionInjector.estimateCallCost(com.google.javascript.rhino.Node fnNode, boolean referencesThis)
func FunctionInjector_EstimateCallCost(fnNode pkg_rhino.Node, referencesThis bool) int32 {
	clinit_FunctionInjector()
	_ = fnNode
	_ = referencesThis
	var argsNode pkg_rhino.Node = NodeUtil_GetFunctionParameters(fnNode)
	_ = argsNode
	var numArgs int32 = argsNode.GetChildCount()
	_ = numArgs
	var callCost int32 = (StaticGet_FunctionInjector_NAME_COST_ESTIMATE() + int32(2))
	_ = callCost
	if numArgs > int32(0) {
		callCost = (callCost + ((numArgs * StaticGet_FunctionInjector_NAME_COST_ESTIMATE()) + ((numArgs - int32(1)) * int32(1))))
	}
	if referencesThis {
		callCost = (callCost + (int32(5) + int32(5)))
	}
	return callCost
}

// Java static method: FunctionInjector.inlineCostDelta(com.google.javascript.rhino.Node fnNode, java.util.Set<java.lang.String> namesToAlias, com.google.javascript.jscomp.FunctionInjector.InliningMode mode)
func FunctionInjector_InlineCostDelta(fnNode pkg_rhino.Node, namesToAlias j2gort.Ext_java_util_Set, mode FunctionInjector_InliningMode) int32 {
	clinit_FunctionInjector()
	_ = fnNode
	_ = namesToAlias
	_ = mode
	var paramCount int32 = NodeUtil_GetFunctionParameters(fnNode).GetChildCount()
	_ = paramCount
	var commaCount int32
	_ = commaCount
	if paramCount > int32(1) {
		commaCount = (paramCount - int32(1))
	} else {
		commaCount = int32(0)
	}
	var costDeltaFunctionOverhead int32 = ((int32(15) + commaCount) + (paramCount * StaticGet_InlineCostEstimator_ESTIMATED_IDENTIFIER_COST()))
	_ = costDeltaFunctionOverhead
	var block pkg_rhino.Node = fnNode.GetLastChild()
	_ = block
	if !(block.HasChildren()) {
		return -(costDeltaFunctionOverhead)
	}
	if mode == StaticGet_FunctionInjector_InliningMode_DIRECT() {
		return -(costDeltaFunctionOverhead + int32(7))
	} else {
		var aliasCount int32 = namesToAlias.Size()
		_ = aliasCount
		var inlineBlockOverhead int32 = int32(4)
		_ = inlineBlockOverhead
		var perReturnOverhead int32 = int32(2)
		_ = perReturnOverhead
		var perReturnResultOverhead int32 = int32(3)
		_ = perReturnResultOverhead
		var perAliasOverhead int32 = int32(3)
		_ = perAliasOverhead
		var returnCount int32 = NodeUtil_GetNodeTypeReferenceCount(block, pkg_rhino.StaticGet_Token_RETURN(), j2gort.ToIface_Ext_com_google_common_base_Predicate(StaticGet_FunctionInjector_MATCH_SHALLOW_STATEMENT()))
		_ = returnCount
		var resultCount int32
		_ = resultCount
		if returnCount > int32(0) {
			resultCount = (returnCount - int32(1))
		} else {
			resultCount = int32(0)
		}
		var baseOverhead int32
		_ = baseOverhead
		if returnCount > int32(0) {
			baseOverhead = inlineBlockOverhead
		} else {
			baseOverhead = int32(0)
		}
		var overhead int32 = (((baseOverhead + (returnCount * perReturnOverhead)) + (resultCount * perReturnResultOverhead)) + (aliasCount * perAliasOverhead))
		_ = overhead
		return (overhead - costDeltaFunctionOverhead)
	}
}

// Java method: FunctionInjector.setKnownConstantFunctions(com.google.common.collect.ImmutableSet<java.lang.String> knownConstantFunctions)
func (this_ *FunctionInjector) SetKnownConstantFunctions(knownConstantFunctions j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = knownConstantFunctions
	if !(this_.knownConstantFunctions.IsEmpty()) {
		j2gort.CheckState(false)
	}
	this_.knownConstantFunctions = knownConstantFunctions
}

// Java class: com.google.javascript.jscomp.FunctionInjector.Builder
type FunctionInjector_Builder struct {
	compiler                 AbstractCompiler                       `gson:"compiler"`
	safeNameIdSupplier       j2gort.Ext_java_util_function_Supplier `gson:"safeNameIdSupplier"`
	functionArgumentInjector *FunctionArgumentInjector              `gson:"functionArgumentInjector"`
	idHash_                  atomic.Int32
	assumeStrictThis         bool `gson:"assumeStrictThis"`
	allowDecomposition       bool `gson:"allowDecomposition"`
}

func (this_ *FunctionInjector_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionInjector_Builder) IsJavaNil_() bool { return this_ == nil }
func (*FunctionInjector_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionInjector$Builder", "Builder"
}
func (this_ *FunctionInjector_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionInjector_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionInjector.Builder", "com.google.javascript.jscomp.FunctionInjector$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionInjector_Builder](v_)
	})
}

func (this_ *FunctionInjector_Builder) initFields_FunctionInjector_Builder() {
	_ = this_

	this_.assumeStrictThis = true
	this_.allowDecomposition = true

}

// Java constructor: Builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewFunctionInjector_Builder(compiler AbstractCompiler) *FunctionInjector_Builder {
	this_ := &FunctionInjector_Builder{}
	InitFunctionInjector_Builder(this_, compiler)
	return this_
}

func InitFunctionInjector_Builder(this_ *FunctionInjector_Builder, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_FunctionInjector_Builder()
	this_.compiler = j2gort.CheckNotNullIface(compiler)
}

// Java method: Builder.safeNameIdSupplier(java.util.function.Supplier<java.lang.String> safeNameIdSupplier)
func (this_ *FunctionInjector_Builder) SafeNameIdSupplier(safeNameIdSupplier j2gort.Ext_java_util_function_Supplier) *FunctionInjector_Builder {
	_ = safeNameIdSupplier
	this_.safeNameIdSupplier = j2gort.CheckNotNull(safeNameIdSupplier)
	return this_
}

// Java method: Builder.allowDecomposition(boolean allowDecomposition)
func (this_ *FunctionInjector_Builder) AllowDecomposition(allowDecomposition bool) *FunctionInjector_Builder {
	_ = allowDecomposition
	this_.allowDecomposition = allowDecomposition
	return this_
}

// Java method: Builder.assumeStrictThis(boolean assumeStrictThis)
func (this_ *FunctionInjector_Builder) AssumeStrictThis(assumeStrictThis bool) *FunctionInjector_Builder {
	_ = assumeStrictThis
	this_.assumeStrictThis = assumeStrictThis
	return this_
}

// Java method: Builder.functionArgumentInjector(com.google.javascript.jscomp.FunctionArgumentInjector functionArgumentInjector)
func (this_ *FunctionInjector_Builder) FunctionArgumentInjector(functionArgumentInjector *FunctionArgumentInjector) *FunctionInjector_Builder {
	_ = functionArgumentInjector
	this_.functionArgumentInjector = j2gort.CheckNotNullPtr(functionArgumentInjector)
	return this_
}

// Java method: Builder.build()
func (this_ *FunctionInjector_Builder) Build() *FunctionInjector {
	if this_.safeNameIdSupplier == nil {
		this_.safeNameIdSupplier = j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](this_.compiler.GetUniqueNameIdSupplier())
	}
	if this_.functionArgumentInjector == nil {
		this_.functionArgumentInjector = NewFunctionArgumentInjector(j2gort.CheckNotNullPtr(this_.compiler.GetAstAnalyzer()))
	}
	return NewFunctionInjector(this_)
}

// Java enum: com.google.javascript.jscomp.FunctionInjector.InliningMode
type FunctionInjector_InliningMode uint8

const (
	enum_FunctionInjector_InliningMode_DIRECT FunctionInjector_InliningMode = 1
	enum_FunctionInjector_InliningMode_BLOCK  FunctionInjector_InliningMode = 2
)

var enumBaseTable_FunctionInjector_InliningMode = [...]j2gort.EnumBase{
	{Name_: "DIRECT", Ordinal_: 0},
	{Name_: "BLOCK", Ordinal_: 1},
}
var enumValues_FunctionInjector_InliningMode = []FunctionInjector_InliningMode{1, 2}

func StaticGet_FunctionInjector_InliningMode_DIRECT() FunctionInjector_InliningMode {
	return enum_FunctionInjector_InliningMode_DIRECT
}
func StaticGet_FunctionInjector_InliningMode_BLOCK() FunctionInjector_InliningMode {
	return enum_FunctionInjector_InliningMode_BLOCK
}

func FunctionInjector_InliningMode_Values() []FunctionInjector_InliningMode {
	return j2gort.CloneSlice(enumValues_FunctionInjector_InliningMode)
}

func FunctionInjector_InliningMode_ValueOf(name string) FunctionInjector_InliningMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_FunctionInjector_InliningMode {
		if enumBaseTable_FunctionInjector_InliningMode[i_].Name_ == name {
			return FunctionInjector_InliningMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.FunctionInjector.InliningMode." + name))
}

func (this_ FunctionInjector_InliningMode) IsJavaNil_() bool                              { return this_ == 0 }
func (FunctionInjector_InliningMode) HasIdentityEquals_()                                 {}
func (FunctionInjector_InliningMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (FunctionInjector_InliningMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ FunctionInjector_InliningMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (FunctionInjector_InliningMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionInjector$InliningMode", "InliningMode"
}
func (this_ FunctionInjector_InliningMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_FunctionInjector_InliningMode[this_-1]
}
func (this_ FunctionInjector_InliningMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ FunctionInjector_InliningMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ FunctionInjector_InliningMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FunctionInjector_InliningMode[this_-1].Name_
}
func (this_ FunctionInjector_InliningMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_FunctionInjector_InliningMode[this_-1].Name_
}
func (this_ FunctionInjector_InliningMode) Equals(other_ any) bool {
	o_, ok_ := other_.(FunctionInjector_InliningMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ FunctionInjector_InliningMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_FunctionInjector_InliningMode[this_-1].HashCode()
}
func (this_ FunctionInjector_InliningMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FunctionInjector_InliningMode[this_-1].CompareTo(other_)
}
func (this_ FunctionInjector_InliningMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_FunctionInjector_InliningMode[this_-1].Name_)
}
func init() {
	vals_ := []any{FunctionInjector_InliningMode(1), FunctionInjector_InliningMode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.FunctionInjector.InliningMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.FunctionInjector$InliningMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionInjector.InliningMode", "com.google.javascript.jscomp.FunctionInjector$InliningMode", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionInjector_InliningMode](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.FunctionInjector.Reference
type FunctionInjector_Reference interface {
	Implements_jscomp_FunctionInjector_Reference()
	IsJavaNil_() bool
	GetFunctionInjector_ReferenceBase() *FunctionInjector_ReferenceBase
	String() string
}

func ToIface_FunctionInjector_Reference[T interface {
	FunctionInjector_Reference
	comparable
}](v T) FunctionInjector_Reference {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_FunctionInjector_Reference[T FunctionInjector_Reference](v T) FunctionInjector_Reference {
	var vi FunctionInjector_Reference = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.FunctionInjector.Reference
type FunctionInjector_ReferenceBase struct {
	scope    *Scope         `gson:"scope"`
	chunk    *JSChunk       `gson:"chunk"`
	callNode pkg_rhino.Node `gson:"callNode"`
	idHash_  atomic.Int32
	mode     FunctionInjector_InliningMode `gson:"mode"`
}

func (this_ *FunctionInjector_ReferenceBase) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionInjector_ReferenceBase) IsJavaNil_() bool { return this_ == nil }
func (*FunctionInjector_ReferenceBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionInjector$Reference", "Reference"
}
func (this_ *FunctionInjector_ReferenceBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionInjector.Reference", "com.google.javascript.jscomp.FunctionInjector$Reference", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionInjector_Reference](v_)
	})
}
func (b *FunctionInjector_ReferenceBase) GetFunctionInjector_ReferenceBase() *FunctionInjector_ReferenceBase {
	return b
}
func (*FunctionInjector_ReferenceBase) Implements_jscomp_FunctionInjector_Reference() {}

func (this_ *FunctionInjector_ReferenceBase) initFields_FunctionInjector_Reference() {
	_ = this_
}

// Java constructor: Reference(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.Scope scope, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.FunctionInjector.InliningMode mode)
func NewFunctionInjector_Reference(callNode pkg_rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) FunctionInjector_Reference {
	this_ := &FunctionInjector_ReferenceBase{}
	InitFunctionInjector_Reference(this_, callNode, scope, chunk, mode)
	return this_
}

func InitFunctionInjector_Reference(this_ *FunctionInjector_ReferenceBase, callNode pkg_rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) {
	_ = callNode
	_ = scope
	_ = chunk
	_ = mode
	this_.initFields_FunctionInjector_Reference()
	this_.callNode = callNode
	this_.scope = scope
	this_.chunk = chunk
	this_.mode = mode
}

// Java method: Reference.toString()
func (this_ *FunctionInjector_ReferenceBase) String() string {
	return ("Reference @ " + j2gort.StrVal(this_.callNode))
}

// Java enum: com.google.javascript.jscomp.FunctionInjector.CallSiteType
type FunctionInjector_CallSiteType struct {
	j2gort.EnumBase
	fn_Prepare func(*FunctionInjector_CallSiteType, *FunctionInjector, FunctionInjector_Reference)
}

func (this_ *FunctionInjector_CallSiteType) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *FunctionInjector_CallSiteType) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*FunctionInjector_CallSiteType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*FunctionInjector_CallSiteType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *FunctionInjector_CallSiteType) initFields_FunctionInjector_CallSiteType() {
	_ = this_
}

func initEnumFunctionInjector_CallSiteType(this_ *FunctionInjector_CallSiteType, enumName_ string, enumOrd_ int32) *FunctionInjector_CallSiteType {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_FunctionInjector_CallSiteType()
	return this_
}

var enum_FunctionInjector_CallSiteType_UNSUPPORTED *FunctionInjector_CallSiteType

func StaticGet_FunctionInjector_CallSiteType_UNSUPPORTED() *FunctionInjector_CallSiteType {
	clinit_FunctionInjector_CallSiteType()
	return enum_FunctionInjector_CallSiteType_UNSUPPORTED
}

var enum_FunctionInjector_CallSiteType_SIMPLE_CALL *FunctionInjector_CallSiteType

func StaticGet_FunctionInjector_CallSiteType_SIMPLE_CALL() *FunctionInjector_CallSiteType {
	clinit_FunctionInjector_CallSiteType()
	return enum_FunctionInjector_CallSiteType_SIMPLE_CALL
}

var enum_FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT *FunctionInjector_CallSiteType

func StaticGet_FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT() *FunctionInjector_CallSiteType {
	clinit_FunctionInjector_CallSiteType()
	return enum_FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT
}

var enum_FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT *FunctionInjector_CallSiteType

func StaticGet_FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT() *FunctionInjector_CallSiteType {
	clinit_FunctionInjector_CallSiteType()
	return enum_FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT
}

var enum_FunctionInjector_CallSiteType_EXPRESSION *FunctionInjector_CallSiteType

func StaticGet_FunctionInjector_CallSiteType_EXPRESSION() *FunctionInjector_CallSiteType {
	clinit_FunctionInjector_CallSiteType()
	return enum_FunctionInjector_CallSiteType_EXPRESSION
}

var enum_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION *FunctionInjector_CallSiteType

func StaticGet_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION() *FunctionInjector_CallSiteType {
	clinit_FunctionInjector_CallSiteType()
	return enum_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION
}

var clinitDone_FunctionInjector_CallSiteType uint32
var clinitState_FunctionInjector_CallSiteType j2gort.ClinitState

func clinit_FunctionInjector_CallSiteType() {
	if atomic.LoadUint32(&clinitDone_FunctionInjector_CallSiteType) != 0 {
		return
	}
	clinitSlow_FunctionInjector_CallSiteType()
}

func clinitSlow_FunctionInjector_CallSiteType() {
	if !j2gort.EnterClinit(&clinitState_FunctionInjector_CallSiteType, &clinitDone_FunctionInjector_CallSiteType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FunctionInjector_CallSiteType, &clinitDone_FunctionInjector_CallSiteType)
	enum_FunctionInjector_CallSiteType_UNSUPPORTED = &FunctionInjector_CallSiteType{}
	enum_FunctionInjector_CallSiteType_UNSUPPORTED.fn_Prepare = func(this_ *FunctionInjector_CallSiteType, injector *FunctionInjector, ref FunctionInjector_Reference) {
		_ = injector
		_ = ref
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected: "+j2gort.StrVal(ref))))
	}
	initEnumFunctionInjector_CallSiteType(enum_FunctionInjector_CallSiteType_UNSUPPORTED, "UNSUPPORTED", int32(0))
	enum_FunctionInjector_CallSiteType_SIMPLE_CALL = &FunctionInjector_CallSiteType{}
	enum_FunctionInjector_CallSiteType_SIMPLE_CALL.fn_Prepare = func(this_ *FunctionInjector_CallSiteType, injector *FunctionInjector, ref FunctionInjector_Reference) {
		_ = injector
		_ = ref
	}
	initEnumFunctionInjector_CallSiteType(enum_FunctionInjector_CallSiteType_SIMPLE_CALL, "SIMPLE_CALL", int32(1))
	enum_FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT = &FunctionInjector_CallSiteType{}
	enum_FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT.fn_Prepare = func(this_ *FunctionInjector_CallSiteType, injector *FunctionInjector, ref FunctionInjector_Reference) {
		_ = injector
		_ = ref
	}
	initEnumFunctionInjector_CallSiteType(enum_FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT, "SIMPLE_ASSIGNMENT", int32(2))
	enum_FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT = &FunctionInjector_CallSiteType{}
	enum_FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT.fn_Prepare = func(this_ *FunctionInjector_CallSiteType, injector *FunctionInjector, ref FunctionInjector_Reference) {
		_ = injector
		_ = ref
	}
	initEnumFunctionInjector_CallSiteType(enum_FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT, "VAR_DECL_SIMPLE_ASSIGNMENT", int32(3))
	enum_FunctionInjector_CallSiteType_EXPRESSION = &FunctionInjector_CallSiteType{}
	enum_FunctionInjector_CallSiteType_EXPRESSION.fn_Prepare = func(this_ *FunctionInjector_CallSiteType, injector *FunctionInjector, ref FunctionInjector_Reference) {
		_ = injector
		_ = ref
		base_ref_ := ref.GetFunctionInjector_ReferenceBase()
		var callNode pkg_rhino.Node = base_ref_.callNode
		_ = callNode
		injector.priv_FunctionInjector_getDecomposer_2944(base_ref_.scope).MoveExpression(callNode)
		var callSiteType *FunctionInjector_CallSiteType = injector.priv_FunctionInjector_classifyCallSite_2943(ref)
		_ = callSiteType
		if !(this_ != callSiteType) {
			j2gort.CheckState(false)
		}
		callSiteType.Prepare(injector, ref)
	}
	initEnumFunctionInjector_CallSiteType(enum_FunctionInjector_CallSiteType_EXPRESSION, "EXPRESSION", int32(4))
	enum_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION = &FunctionInjector_CallSiteType{}
	enum_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION.fn_Prepare = func(this_ *FunctionInjector_CallSiteType, injector *FunctionInjector, ref FunctionInjector_Reference) {
		_ = injector
		_ = ref
		base_ref_ := ref.GetFunctionInjector_ReferenceBase()
		var callNode pkg_rhino.Node = base_ref_.callNode
		_ = callNode
		injector.priv_FunctionInjector_getDecomposer_2944(base_ref_.scope).MaybeExposeExpression(callNode)
		var callSiteType *FunctionInjector_CallSiteType = injector.priv_FunctionInjector_classifyCallSite_2943(ref)
		_ = callSiteType
		if !(this_ != callSiteType) {
			j2gort.CheckState(false)
		}
		callSiteType.Prepare(injector, ref)
	}
	initEnumFunctionInjector_CallSiteType(enum_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION, "DECOMPOSABLE_EXPRESSION", int32(5))
	enumValues_FunctionInjector_CallSiteType = []*FunctionInjector_CallSiteType{enum_FunctionInjector_CallSiteType_UNSUPPORTED, enum_FunctionInjector_CallSiteType_SIMPLE_CALL, enum_FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT, enum_FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT, enum_FunctionInjector_CallSiteType_EXPRESSION, enum_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.FunctionInjector.CallSiteType", []any{enum_FunctionInjector_CallSiteType_UNSUPPORTED, enum_FunctionInjector_CallSiteType_SIMPLE_CALL, enum_FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT, enum_FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT, enum_FunctionInjector_CallSiteType_EXPRESSION, enum_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.FunctionInjector$CallSiteType", []any{enum_FunctionInjector_CallSiteType_UNSUPPORTED, enum_FunctionInjector_CallSiteType_SIMPLE_CALL, enum_FunctionInjector_CallSiteType_SIMPLE_ASSIGNMENT, enum_FunctionInjector_CallSiteType_VAR_DECL_SIMPLE_ASSIGNMENT, enum_FunctionInjector_CallSiteType_EXPRESSION, enum_FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.FunctionInjector.CallSiteType", "com.google.javascript.jscomp.FunctionInjector$CallSiteType", func() {
		clinit_FunctionInjector_CallSiteType()
	})
}

var enumValues_FunctionInjector_CallSiteType []*FunctionInjector_CallSiteType

func FunctionInjector_CallSiteType_Values() []*FunctionInjector_CallSiteType {
	clinit_FunctionInjector_CallSiteType()
	return j2gort.CloneSlice(enumValues_FunctionInjector_CallSiteType)
}

func FunctionInjector_CallSiteType_ValueOf(name string) *FunctionInjector_CallSiteType {
	clinit_FunctionInjector_CallSiteType()
	for _, c := range enumValues_FunctionInjector_CallSiteType {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.FunctionInjector.CallSiteType." + name))
}

// Java method: CallSiteType.prepare(com.google.javascript.jscomp.FunctionInjector injector, com.google.javascript.jscomp.FunctionInjector.Reference ref)
func (this_ *FunctionInjector_CallSiteType) Prepare(injector *FunctionInjector, ref FunctionInjector_Reference) {
	if this_.fn_Prepare != nil {
		this_.fn_Prepare(this_, injector, ref)
		return
	}
	this_.baseEnum_Prepare(injector, ref)
}

func (this_ *FunctionInjector_CallSiteType) baseEnum_Prepare(injector *FunctionInjector, ref FunctionInjector_Reference) {
	_ = injector
	_ = ref
	panic("abstract or native method invoked")
}

func (this_ *FunctionInjector_CallSiteType) IsJavaNil_() bool { return this_ == nil }
func (*FunctionInjector_CallSiteType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionInjector$CallSiteType", "CallSiteType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionInjector.CallSiteType", "com.google.javascript.jscomp.FunctionInjector$CallSiteType", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionInjector_CallSiteType](v_)
	})
}

// Java enum: com.google.javascript.jscomp.FunctionInjector.CanInlineResult
type FunctionInjector_CanInlineResult uint8

const (
	enum_FunctionInjector_CanInlineResult_YES               FunctionInjector_CanInlineResult = 1
	enum_FunctionInjector_CanInlineResult_AFTER_PREPARATION FunctionInjector_CanInlineResult = 2
	enum_FunctionInjector_CanInlineResult_NO                FunctionInjector_CanInlineResult = 3
)

var enumBaseTable_FunctionInjector_CanInlineResult = [...]j2gort.EnumBase{
	{Name_: "YES", Ordinal_: 0},
	{Name_: "AFTER_PREPARATION", Ordinal_: 1},
	{Name_: "NO", Ordinal_: 2},
}
var enumValues_FunctionInjector_CanInlineResult = []FunctionInjector_CanInlineResult{1, 2, 3}

func StaticGet_FunctionInjector_CanInlineResult_YES() FunctionInjector_CanInlineResult {
	return enum_FunctionInjector_CanInlineResult_YES
}
func StaticGet_FunctionInjector_CanInlineResult_AFTER_PREPARATION() FunctionInjector_CanInlineResult {
	return enum_FunctionInjector_CanInlineResult_AFTER_PREPARATION
}
func StaticGet_FunctionInjector_CanInlineResult_NO() FunctionInjector_CanInlineResult {
	return enum_FunctionInjector_CanInlineResult_NO
}

func FunctionInjector_CanInlineResult_Values() []FunctionInjector_CanInlineResult {
	return j2gort.CloneSlice(enumValues_FunctionInjector_CanInlineResult)
}

func FunctionInjector_CanInlineResult_ValueOf(name string) FunctionInjector_CanInlineResult {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_FunctionInjector_CanInlineResult {
		if enumBaseTable_FunctionInjector_CanInlineResult[i_].Name_ == name {
			return FunctionInjector_CanInlineResult(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.FunctionInjector.CanInlineResult." + name))
}

func (this_ FunctionInjector_CanInlineResult) IsJavaNil_() bool                              { return this_ == 0 }
func (FunctionInjector_CanInlineResult) HasIdentityEquals_()                                 {}
func (FunctionInjector_CanInlineResult) Implements_j2gort_Ext_java_lang_Enum()               {}
func (FunctionInjector_CanInlineResult) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ FunctionInjector_CanInlineResult) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (FunctionInjector_CanInlineResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionInjector$CanInlineResult", "CanInlineResult"
}
func (this_ FunctionInjector_CanInlineResult) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_FunctionInjector_CanInlineResult[this_-1]
}
func (this_ FunctionInjector_CanInlineResult) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ FunctionInjector_CanInlineResult) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ FunctionInjector_CanInlineResult) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FunctionInjector_CanInlineResult[this_-1].Name_
}
func (this_ FunctionInjector_CanInlineResult) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_FunctionInjector_CanInlineResult[this_-1].Name_
}
func (this_ FunctionInjector_CanInlineResult) Equals(other_ any) bool {
	o_, ok_ := other_.(FunctionInjector_CanInlineResult)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ FunctionInjector_CanInlineResult) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_FunctionInjector_CanInlineResult[this_-1].HashCode()
}
func (this_ FunctionInjector_CanInlineResult) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FunctionInjector_CanInlineResult[this_-1].CompareTo(other_)
}
func (this_ FunctionInjector_CanInlineResult) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_FunctionInjector_CanInlineResult[this_-1].Name_)
}
func init() {
	vals_ := []any{FunctionInjector_CanInlineResult(1), FunctionInjector_CanInlineResult(2), FunctionInjector_CanInlineResult(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.FunctionInjector.CanInlineResult", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.FunctionInjector$CanInlineResult", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionInjector.CanInlineResult", "com.google.javascript.jscomp.FunctionInjector$CanInlineResult", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionInjector_CanInlineResult](v_)
	})
}

// Java class:
type FunctionInjector__Anon_45 struct {
	outer_  *FunctionInjector
	nextId  int32 `gson:"nextId"`
	idHash_ atomic.Int32
}

func (this_ *FunctionInjector__Anon_45) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionInjector__Anon_45) IsJavaNil_() bool { return this_ == nil }
func (*FunctionInjector__Anon_45) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionInjector$1", ""
}
func (this_ *FunctionInjector__Anon_45) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionInjector__Anon_45) HasIdentityEquals_()                                {}
func (this_ *FunctionInjector__Anon_45) SetOuter_(o *FunctionInjector)                { this_.outer_ = o }
func (this_ *FunctionInjector__Anon_45) GetOuter_() *FunctionInjector                 { return this_.outer_ }
func (*FunctionInjector__Anon_45) Implements_j2gort_Ext_java_util_function_Supplier() {}

func (this_ *FunctionInjector__Anon_45) initFields_FunctionInjector__Anon_45() {
	_ = this_

}

func NewFunctionInjector__Anon_45(outer_ *FunctionInjector) *FunctionInjector__Anon_45 {
	this_ := &FunctionInjector__Anon_45{}
	this_.outer_ = outer_
	this_.initFields_FunctionInjector__Anon_45()
	return this_
}

// Java method: .get()
func (this_ *FunctionInjector__Anon_45) Get() any {
	var iifeArg_1 int32
	{

		old_7 := this_.nextId
		nw_8 := old_7 + 1
		this_.nextId = nw_8
		iifeArg_1 = old_7
	}
	return strconv.FormatInt(int64(iifeArg_1), 10)
}

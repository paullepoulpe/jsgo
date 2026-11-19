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
// Port of com/google/javascript/jscomp/CheckMissingReturn.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_check_missing_return_1 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_check_missing_return_2 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_check_missing_return_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_1 uint32) bool {
		input := pkg_rhino.Node(rawNode_1)
		_ = input
		return ((input != 0) && input.IsReturn())
	})
	lambdaSingleton_check_missing_return_2 = j2gort.NewLambda_Ext_com_google_common_base_Predicate(func(mrp0_ any) bool {
		return CheckMissingReturn_GoesThroughTrueCondition(j2gort.AsIface[pkg_graph.DiGraph_DiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch]](mrp0_))
	})
}

// Java class: com.google.javascript.jscomp.CheckMissingReturn
type CheckMissingReturn struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler   AbstractCompiler `gson:"compiler"`
	convention CodingConvention `gson:"convention"`
}

func (this_ *CheckMissingReturn) BindSelf_(s any) {
	_ = s
	this_.NodeTraversal_AbstractCfgCallbackBase.BindSelf_(s)
}
func (this_ *CheckMissingReturn) IsJavaNil_() bool { return this_ == nil }
func (*CheckMissingReturn) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckMissingReturn", "CheckMissingReturn"
}
func (this_ *CheckMissingReturn) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractCfgCallbackBase.IdentityHash_()
}
func (*CheckMissingReturn) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckMissingReturn", "com.google.javascript.jscomp.CheckMissingReturn", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckMissingReturn](v_)
	})
}
func (*CheckMissingReturn) Implements_jscomp_NodeTraversal_AbstractCfgCallback() {}
func (*CheckMissingReturn) Implements_jscomp_NodeTraversal_ScopedCallback()      {}
func (*CheckMissingReturn) Implements_jscomp_NodeTraversal_Callback()            {}

var static_CheckMissingReturn_MISSING_RETURN_STATEMENT *DiagnosticType

func StaticGet_CheckMissingReturn_MISSING_RETURN_STATEMENT() *DiagnosticType {
	clinit_CheckMissingReturn()
	return static_CheckMissingReturn_MISSING_RETURN_STATEMENT
}
func StaticSet_CheckMissingReturn_MISSING_RETURN_STATEMENT(v *DiagnosticType) {
	clinit_CheckMissingReturn()
	static_CheckMissingReturn_MISSING_RETURN_STATEMENT = v
}

var clinitDone_CheckMissingReturn uint32
var clinitState_CheckMissingReturn j2gort.ClinitState

func clinit_CheckMissingReturn() {
	if atomic.LoadUint32(&clinitDone_CheckMissingReturn) != 0 {
		return
	}
	clinitSlow_CheckMissingReturn()
}

func clinitSlow_CheckMissingReturn() {
	if !j2gort.EnterClinit(&clinitState_CheckMissingReturn, &clinitDone_CheckMissingReturn) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckMissingReturn, &clinitDone_CheckMissingReturn)
	static_CheckMissingReturn_MISSING_RETURN_STATEMENT = DiagnosticType_Warning("JSC_MISSING_RETURN_STATEMENT", "Missing return statement. Function expected to return {0}.")
}

func (this_ *CheckMissingReturn) initFields_CheckMissingReturn() {
	_ = this_
}

// Java constructor: CheckMissingReturn(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckMissingReturn(compiler AbstractCompiler) *CheckMissingReturn {
	clinit_CheckMissingReturn()
	this_ := &CheckMissingReturn{}
	this_.BindSelf_(this_)
	InitCheckMissingReturn(this_, compiler)
	return this_
}

func InitCheckMissingReturn(this_ *CheckMissingReturn, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractCfgCallback(&this_.NodeTraversal_AbstractCfgCallbackBase)
	this_.initFields_CheckMissingReturn()
	this_.compiler = compiler
	this_.convention = compiler.GetCodingConvention()
}

// Java static method: CheckMissingReturn.goesThroughTrueCondition(com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge<com.google.javascript.rhino.Node,com.google.javascript.jscomp.ControlFlowGraph.Branch> input)
func CheckMissingReturn_GoesThroughTrueCondition(input pkg_graph.DiGraph_DiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch]) bool {
	clinit_CheckMissingReturn()
	_ = input
	var branch ControlFlowGraph_Branch = input.GetValue()
	_ = branch
	if branch == StaticGet_ControlFlowGraph_Branch_ON_EX() {
		return false
	} else if branch.IsConditional() {
		var condition pkg_rhino.Node = NodeUtil_GetConditionExpression(input.GetSource().GetValue())
		_ = condition
		if condition != 0 {
			var val pkg_base.Tri = NodeUtil_GetBooleanValue(condition)
			_ = val
			if val != pkg_base.StaticGet_Tri_UNKNOWN() {
				return (val.ToBoolean(true) == (branch == StaticGet_ControlFlowGraph_Branch_ON_TRUE()))
			}
		}
	}
	return true
}

// Java method: CheckMissingReturn.enterScopeWithCfg(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CheckMissingReturn) EnterScopeWithCfg(t *NodeTraversal) {
	_ = t
	var n pkg_rhino.Node = t.GetScopeRoot()
	_ = n
	var returnType pkg_rhino.JSType = this_.priv_CheckMissingReturn_getExplicitReturnTypeIfExpected_1102(n)
	_ = returnType
	if returnType == nil || returnType.IsJavaNil_() {
		return
	}
	if n.IsArrowFunction() {
		var functionBody pkg_rhino.Node = NodeUtil_GetFunctionBody(n)
		_ = functionBody
		if !(functionBody.IsBlock()) {
			return
		}
	}
	var cfg ControlFlowGraph[pkg_rhino.Node] = this_.GetControlFlowGraph(this_.compiler)
	_ = cfg
	if this_.priv_CheckMissingReturn_fastAllPathsReturnCheck_1101(cfg) {
		return
	}
	var test *pkg_graph.CheckPathsBetweenNodes[pkg_rhino.Node, ControlFlowGraph_Branch] = pkg_graph.NewCheckPathsBetweenNodes_5[pkg_rhino.Node, ControlFlowGraph_Branch](j2gort.AsIface[pkg_graph.DiGraph[pkg_rhino.Node, ControlFlowGraph_Branch]](cfg), cfg.GetEntry(), cfg.GetImplicitReturn(), lambdaSingleton_check_missing_return_1, lambdaSingleton_check_missing_return_2)
	_ = test
	if !(test.AllPathsSatisfyPredicate()) {
		this_.compiler.Report_1(JSError_Make_3(t.GetScopeRoot(), StaticGet_CheckMissingReturn_MISSING_RETURN_STATEMENT(), []string{returnType.String()}))
	}
}

// Java method: CheckMissingReturn.fastAllPathsReturnCheck(com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> cfg)
func (this_ *CheckMissingReturn) priv_CheckMissingReturn_fastAllPathsReturnCheck_1101(cfg ControlFlowGraph[pkg_rhino.Node]) bool {
	_ = cfg
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(cfg.GetImplicitReturn().GetInEdges()))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			s := j2gort.AsIface[pkg_graph.DiGraph_DiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch]](rawElem_3)
			_ = s
			var n pkg_rhino.Node = s.GetSource().GetValue()
			_ = n
			if !(n.IsReturn()) && !(this_.convention.IsFunctionCallThatAlwaysThrows(n)) {
				return false
			}
		}
	}
	return true
}

// Java method: CheckMissingReturn.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckMissingReturn) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	return true
}

// Java method: CheckMissingReturn.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckMissingReturn) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}

// Java method: CheckMissingReturn.getExplicitReturnTypeIfExpected(com.google.javascript.rhino.Node scopeRoot)
func (this_ *CheckMissingReturn) priv_CheckMissingReturn_getExplicitReturnTypeIfExpected_1102(scopeRoot pkg_rhino.Node) pkg_rhino.JSType {
	_ = scopeRoot
	if !(scopeRoot.IsFunction()) {
		return nil
	}
	var scopeType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(scopeRoot.GetJSType())
	_ = scopeType
	if scopeType == nil {
		return nil
	}
	if CheckMissingReturn_IsEmptyFunction(scopeRoot) {
		return nil
	}
	if scopeType.IsConstructor() {
		return nil
	}
	var returnType pkg_rhino.JSType = scopeType.GetReturnType()
	_ = returnType
	if returnType == nil || returnType.IsJavaNil_() {
		return nil
	}
	if scopeRoot.IsAsyncFunction() {
		returnType = Promises_GetTemplateTypeOfThenable(this_.compiler.GetTypeRegistry(), returnType)
	} else if scopeRoot.IsGeneratorFunction() {
		returnType = JsIterables_GetReturnElementType(returnType, this_.compiler.GetTypeRegistry())
	}
	if !(this_.priv_CheckMissingReturn_isVoidOrUnknown_1103(returnType)) {
		return returnType
	}
	return nil
}

// Java static method: CheckMissingReturn.isEmptyFunction(com.google.javascript.rhino.Node function)
func CheckMissingReturn_IsEmptyFunction(function pkg_rhino.Node) bool {
	clinit_CheckMissingReturn()
	_ = function
	return (function.HasXChildren(int32(3)) && !(function.GetSecondChild().GetNext().HasChildren()))
}

// Java method: CheckMissingReturn.isVoidOrUnknown(com.google.javascript.rhino.jstype.JSType returnType)
func (this_ *CheckMissingReturn) priv_CheckMissingReturn_isVoidOrUnknown_1103(returnType pkg_rhino.JSType) bool {
	_ = returnType
	var voidType pkg_rhino.JSType = this_.compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
	_ = voidType
	return voidType.IsSubtypeOf_1(returnType)
}
func (*CheckMissingReturn) AlwaysShouldTraverse_()      {}
func (*CheckMissingReturn) NoOpVisit_()                 {}
func (*CheckMissingReturn) CanReleaseSyntacticScopes_() {}

// Port of com/google/javascript/jscomp/CheckMissingReturn.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckMissingReturn checks functions for missing return statements. Return
// statements are only expected for functions with return type information.
// Functions with empty bodies are ignored.
type CheckMissingReturn struct {
	NodeTraversal_AbstractCfgCallbackBase

	compiler        AbstractCompiler
	convention      CodingConvention
	sharedCfgSource *CheckUnreachableCode
}

// Java: CheckMissingReturn.MISSING_RETURN_STATEMENT
var CheckMissingReturn_MISSING_RETURN_STATEMENT = DiagnosticType_Warning(
	"JSC_MISSING_RETURN_STATEMENT",
	"Missing return statement. Function expected to return {0}.")

// CheckMissingReturn_GoesThroughTrueCondition skips all exception edges and
// impossible edges.
//
// Java: CheckMissingReturn.goesThroughTrueCondition(DiGraphEdge<Node, ControlFlowGraph.Branch>)
func CheckMissingReturn_GoesThroughTrueCondition(
	input DiGraph_DiGraphEdge[rhino.Node, ControlFlowGraph_Branch]) bool {
	// First skill all exceptions.
	branch := input.GetValue()
	if branch == ControlFlowGraph_Branch_ON_EX {
		return false
	} else if branch.IsConditional() {
		condition := NodeUtil_GetConditionExpression(input.GetSource().GetValue())
		// TODO(user): We CAN make this bit smarter just looking at
		// constants. We DO have a full blown ReverseAbstractInterupter and
		// type system that can evaluate some impressions' boolean value but
		// for now we will keep this pass lightweight.
		if !condition.IsNil() {
			val := NodeUtil_GetBooleanValue(condition)
			if val != base.Tri_UNKNOWN {
				return val.ToBoolean(true) == (branch == ControlFlowGraph_Branch_ON_TRUE)
			}
		}
	}
	return true
}

func checkMissingReturn_isReturn(input rhino.Node) bool {
	return !input.IsNil() && input.IsReturn()
}

// Java: CheckMissingReturn(AbstractCompiler)
func NewCheckMissingReturn(compiler AbstractCompiler) *CheckMissingReturn {
	c := &CheckMissingReturn{}
	c.InitNodeTraversal_AbstractCfgCallbackBase(c)
	c.compiler = compiler
	c.convention = compiler.GetCodingConvention()
	return c
}

// Java: CheckMissingReturn.enterScopeWithCfg(NodeTraversal)
func (c *CheckMissingReturn) EnterScopeWithCfg(t *NodeTraversal) {
	n := t.GetScopeRoot()
	returnType := c.GetExplicitReturnTypeIfExpected(n)

	if returnType == nil {
		// No return value is expected, so nothing to check.
		return
	}

	if n.IsArrowFunction() {
		functionBody := NodeUtil_GetFunctionBody(n)
		if !functionBody.IsBlock() {
			// Body is an expression, which is the implicit return value.
			return
		}
	}

	if straightLine, endsWithReturn := isStraightLineFunction(n); straightLine && endsWithReturn {
		return
	}

	var cfg ControlFlowGraph[rhino.Node]
	if c.sharedCfgSource != nil {
		if shared := c.sharedCfgSource.PeekCurrentScopeCfg(); shared != nil &&
			shared.GetEntry() != nil &&
			shared.GetEntry().GetValue() == ControlFlowAnalysis_ComputeFallThrough(n) {
			cfg = shared
		}
	}
	if cfg == nil {
		cfg = c.GetControlFlowGraph(c.compiler)
	}
	if c.FastAllPathsReturnCheck(cfg) {
		return
	}

	test := CheckPathsBetweenNodes[rhino.Node, ControlFlowGraph_Branch]{
		graph:         cfg,
		start:         cfg.GetEntry(),
		end:           cfg.GetImplicitReturn(),
		nodePredicate: checkMissingReturn_isReturn,
		edgePredicate: CheckMissingReturn_GoesThroughTrueCondition,
		inclusive:     true,
	}

	if !test.AllPathsSatisfyPredicate() {
		c.compiler.Report(
			JSError_MakeWithNode(t.GetScopeRoot(), CheckMissingReturn_MISSING_RETURN_STATEMENT, returnType.String()))
	}
}

// FastAllPathsReturnCheck is a fast check to see if all execution paths
// contain a return statement. May spuriously report that a return statement
// is missing.
//
// Returns true if all paths return, converse not necessarily true.
//
// Java: CheckMissingReturn.fastAllPathsReturnCheck(ControlFlowGraph<Node>)
func (c *CheckMissingReturn) FastAllPathsReturnCheck(cfg ControlFlowGraph[rhino.Node]) bool {
	for _, s := range cfg.GetImplicitReturn().GetInEdges() {
		n := s.GetSource().GetValue()
		// NOTE(dimvar): it is possible to change ControlFlowAnalysis.java, so
		// that the calls that always throw are treated in the same way as THROW
		// in the CFG. Then, we would not need to use the coding convention here.
		if !n.IsReturn() && !c.convention.IsFunctionCallThatAlwaysThrows(n) {
			return false
		}
	}
	return true
}

// Java: CheckMissingReturn.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckMissingReturn) ShouldTraverse(
	nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// Java: CheckMissingReturn.visit(NodeTraversal, Node, Node)
func (c *CheckMissingReturn) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
}

func (*CheckMissingReturn) nodeTraversalNoVisit()          {}
func (*CheckMissingReturn) nodeTraversalNoShouldTraverse() {}

// GetExplicitReturnTypeIfExpected determines if the given scope should
// explicitly return. All functions with non-void or non-unknown return types
// must have explicit returns.
//
// Exception: Constructors which specifically specify a return type are used
// to allow invocation without requiring the "new" keyword. They have an
// implicit return type. See unit tests.
//
// If a return type is expected, returns it. Otherwise, returns null.
//
// Java: CheckMissingReturn.getExplicitReturnTypeIfExpected(Node)
func (c *CheckMissingReturn) GetExplicitReturnTypeIfExpected(scopeRoot rhino.Node) rhino.JSType {
	if !scopeRoot.IsFunction() {
		// Nothing to do in a global/module/block scope.
		return nil
	}
	scopeType := rhino.JSType_ToMaybeFunctionType(scopeRoot.GetJSType())

	if scopeType == nil {
		return nil
	}

	if CheckMissingReturn_IsEmptyFunction(scopeRoot) {
		return nil
	}

	if scopeType.IsConstructor() {
		return nil
	}

	returnType := scopeType.GetReturnType()

	if returnType == nil {
		return nil
	}

	if scopeRoot.IsAsyncFunction() {
		// Unwrap the declared return type (e.g. "!Promise<number>" becomes "number")
		returnType = Promises_GetTemplateTypeOfThenable(c.compiler.GetTypeRegistry(), returnType)
	} else if scopeRoot.IsGeneratorFunction() {
		// Unwrap the declared return type (e.g. "!Generator<string, number>" becomes "number")
		returnType = JsIterables_GetReturnElementType(returnType, c.compiler.GetTypeRegistry())
	}

	if !c.IsVoidOrUnknown(returnType) {
		return returnType
	}

	return nil
}

// CheckMissingReturn_IsEmptyFunction returns true if function represents a
// JavaScript function with an empty body.
//
// Java: CheckMissingReturn.isEmptyFunction(Node)
func CheckMissingReturn_IsEmptyFunction(function rhino.Node) bool {
	return function.HasXChildren(3) && !function.GetSecondChild().GetNext().HasChildren()
}

// IsVoidOrUnknown returns true if returnType is void, unknown, or a union
// containing void or unknown.
//
// Java: CheckMissingReturn.isVoidOrUnknown(JSType)
func (c *CheckMissingReturn) IsVoidOrUnknown(returnType rhino.JSType) bool {
	voidType := c.compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	return voidType.IsSubtypeOf(returnType)
}

// Port of com/google/javascript/jscomp/CheckUnreachableCode.java @ closure-compiler v20261015

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

// CheckUnreachableCode uses ControlFlowGraph and GraphReachability to inform
// user about unreachable code.
type CheckUnreachableCode struct {
	NodeTraversal_AbstractCfgCallbackBase

	compiler            AbstractCompiler
	reachability        GraphReachability[*rhino.Node, ControlFlowGraph_Branch]
	hasUnreachableStack []bool
	hasUnreachableBuf   [8]bool
}

// Java: CheckUnreachableCode.UNREACHABLE_CODE
var CheckUnreachableCode_UNREACHABLE_CODE = DiagnosticType_Warning(
	"JSC_UNREACHABLE_CODE", "unreachable code")

// Java: CheckUnreachableCode(AbstractCompiler)
func NewCheckUnreachableCode(compiler AbstractCompiler) *CheckUnreachableCode {
	c := &CheckUnreachableCode{}
	c.InitNodeTraversal_AbstractCfgCallbackBase(c)
	c.compiler = compiler
	c.hasUnreachableStack = c.hasUnreachableBuf[:0]
	return c
}

func (*CheckUnreachableCode) nodeTraversalNoVisit() {}

func hasUnreachableCfgNode(cfg ControlFlowGraph[*rhino.Node]) bool {
	if base := cfg.GetControlFlowGraphBase(); base != nil {
		for n, gNode := range base.nodes.All() {
			if n != nil && gNode.GetAnnotation() != GraphReachability_REACHABLE &&
				n.GetLineno() != -1 && !n.IsEmpty() && !n.IsBreak() {
				return true
			}
		}
		return false
	}
	for _, gNode := range cfg.GetNodes() {
		n := gNode.GetValue()
		if n != nil && gNode.GetAnnotation() != GraphReachability_REACHABLE &&
			n.GetLineno() != -1 && !n.IsEmpty() && !n.IsBreak() {
			return true
		}
	}
	return false
}

// Java: CheckUnreachableCode.enterScopeWithCfg(NodeTraversal)
func (c *CheckUnreachableCode) EnterScopeWithCfg(t *NodeTraversal) {
	if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) {
		cfg := c.GetControlFlowGraph(c.compiler)
		c.InitScope(cfg)
		c.hasUnreachableStack = append(c.hasUnreachableStack, hasUnreachableCfgNode(cfg))
	}
}

// Java: CheckUnreachableCode.exitScopeWithCfg(NodeTraversal)
func (c *CheckUnreachableCode) ExitScopeWithCfg(t *NodeTraversal) {
	if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) && len(c.hasUnreachableStack) > 0 {
		c.hasUnreachableStack = c.hasUnreachableStack[:len(c.hasUnreachableStack)-1]
	}
}

// Java: CheckUnreachableCode.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckUnreachableCode) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if len(c.hasUnreachableStack) > 0 && !c.hasUnreachableStack[len(c.hasUnreachableStack)-1] {
		return true
	}
	// Only report error when there are some line number informations.
	// There are synthetic nodes with no line number informations, nodes
	// introduce by other passes (although not likely since this pass should
	// be executed early) or some rhino bug.
	// Allow spurious semi-colons and spurious breaks.
	if n.GetLineno() == -1 || n.IsEmpty() || n.IsBreak() {
		return true
	}
	cfg := c.GetControlFlowGraph(c.compiler)
	gNode := cfg.GetNode(n)
	if gNode != nil && gNode.GetAnnotation() != GraphReachability_REACHABLE {
		c.compiler.Report(JSError_MakeWithNode(n, CheckUnreachableCode_UNREACHABLE_CODE))
		// From now on, we are going to assume the user fixed the error and not
		// give more warning related to code section reachable from this node.
		c.reachability.graph = cfg
		c.reachability.edgePredicate = nil
		c.reachability.Recompute(n)

		// Saves time by not traversing children.
		return false
	}
	return true
}

// Java: CheckUnreachableCode.initScope(ControlFlowGraph<Node>)
func (c *CheckUnreachableCode) InitScope(controlFlowGraph ControlFlowGraph[*rhino.Node]) {
	c.reachability.graph = controlFlowGraph
	c.reachability.edgePredicate = CheckUnreachableCode_IsReachable
	c.reachability.Compute(controlFlowGraph.GetEntry().GetValue())
}

// Java: CheckUnreachableCode.isReachable(EdgeTuple<Node, Branch>)
func CheckUnreachableCode_IsReachable(input *GraphReachability_EdgeTuple[*rhino.Node, ControlFlowGraph_Branch]) bool {
	branch := input.edge
	if !branch.IsConditional() {
		return true
	}
	predecessor := input.sourceNode
	condition := NodeUtil_GetConditionExpression(predecessor)

	// TODO(user): Handle more complicated expression like true == true,
	// etc....
	if condition != nil {
		val := NodeUtil_GetBooleanValue(condition)
		if val != base.Tri_UNKNOWN {
			return val.ToBoolean(true) == (branch == ControlFlowGraph_Branch_ON_TRUE)
		}
	}
	return true
}

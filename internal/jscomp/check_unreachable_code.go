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

	compiler             AbstractCompiler
	reachability         GraphReachability[rhino.Node, ControlFlowGraph_Branch]
	hasUnreachableStack  []bool
	hasUnreachableBuf    [8]bool
	scriptCheckCfg       ReusableCfg
	rootHasNoUnreachable bool
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

func hasUnreachableCfgNode(cfg ControlFlowGraph[rhino.Node]) bool {
	if base := cfg.GetControlFlowGraphBase(); base != nil {
		for n, gNode := range base.nodes.All() {
			if !n.IsNil() && gNode.GetAnnotation() != GraphReachability_REACHABLE &&
				n.GetLineno() != -1 && !n.IsEmpty() && !n.IsBreak() {
				return true
			}
		}
		return false
	}
	for _, gNode := range cfg.GetNodes() {
		n := gNode.GetValue()
		if !n.IsNil() && gNode.GetAnnotation() != GraphReachability_REACHABLE &&
			n.GetLineno() != -1 && !n.IsEmpty() && !n.IsBreak() {
			return true
		}
	}
	return false
}

// CurrentScopeHasNoUnreachable reports whether the current CFG scope is known
// to have no unreachable CFG nodes.
func (c *CheckUnreachableCode) CurrentScopeHasNoUnreachable() bool {
	return len(c.hasUnreachableStack) > 0 && !c.hasUnreachableStack[len(c.hasUnreachableStack)-1]
}

// isStraightLineFunction reports whether fn's parameter list and block body
// contain only straight-line statements that cannot create unreachable code,
// and whether the body ends with an explicit return statement.
func isStraightLineFunction(fn rhino.Node) (straightLine bool, endsWithReturn bool) {
	if fn.IsNil() || !fn.IsFunction() {
		return false, false
	}
	paramList := fn.GetSecondChild()
	if paramList.IsNil() {
		return false, false
	}
	for p := paramList.GetFirstChild(); !p.IsNil(); p = p.GetNext() {
		if !p.IsName() {
			return false, false
		}
	}
	body := paramList.GetNext()
	if body.IsNil() || !body.IsBlock() || body.IsSyntheticBlock() {
		return false, false
	}
	terminated := false
	for child := body.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		tok := child.GetToken()
		if tok == rhino.Token_FUNCTION {
			continue
		}
		if terminated {
			return false, false
		}
		switch tok {
		case rhino.Token_EXPR_RESULT, rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_EMPTY:
		case rhino.Token_RETURN:
			terminated = true
			endsWithReturn = true
		case rhino.Token_THROW:
			terminated = true
		default:
			return false, false
		}
	}
	return true, endsWithReturn
}

func isStraightLineScript(script rhino.Node) bool {
	container := script
	if script.HasOneChild() && script.GetFirstChild().IsModuleBody() {
		container = script.GetFirstChild()
	}
	for child := container.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		switch child.GetToken() {
		case rhino.Token_EXPR_RESULT, rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST,
			rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_EMPTY,
			rhino.Token_IMPORT, rhino.Token_EXPORT:
		default:
			return false
		}
	}
	return true
}

func hasTopLevelReturn(n rhino.Node) bool {
	if n.IsFunction() {
		return false
	}
	if n.IsReturn() {
		return true
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if hasTopLevelReturn(c) {
			return true
		}
	}
	return false
}

func (c *CheckUnreachableCode) hasNoUnreachableCodeFast(scopeRoot rhino.Node) bool {
	switch scopeRoot.GetToken() {
	case rhino.Token_FUNCTION:
		straightLine, _ := isStraightLineFunction(scopeRoot)
		return straightLine
	case rhino.Token_SCRIPT, rhino.Token_MODULE_BODY:
		return isStraightLineScript(scopeRoot)
	case rhino.Token_ROOT:
		if c.rootHasNoUnreachable {
			return true
		}
		for s := scopeRoot.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
			if !s.IsScript() {
				return false
			}
			if isStraightLineScript(s) {
				continue
			}
			if hasTopLevelReturn(s) {
				return false
			}
			cfg := c.scriptCheckCfg.Compute(
				c.compiler,
				s,
				/* shouldTraverseFunctions= */ false,
				/* edgeAnnotations= */ true)
			c.InitScope(cfg)
			unreachable := hasUnreachableCfgNode(cfg) || cfg.GetImplicitReturn().GetAnnotation() != GraphReachability_REACHABLE
			c.scriptCheckCfg.release()
			if unreachable {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// Java: CheckUnreachableCode.enterScopeWithCfg(NodeTraversal)
func (c *CheckUnreachableCode) EnterScopeWithCfg(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	if NodeUtil_IsValidCfgRoot(scopeRoot) {
		if c.hasNoUnreachableCodeFast(scopeRoot) {
			c.hasUnreachableStack = append(c.hasUnreachableStack, false)
			if t != nil && t.callback == c {
				t.noShouldTraverse = true
			}
			return
		}
		cfg := c.GetControlFlowGraph(c.compiler)
		c.InitScope(cfg)
		hasUnreachable := hasUnreachableCfgNode(cfg)
		c.hasUnreachableStack = append(c.hasUnreachableStack, hasUnreachable)
		if t != nil && t.callback == c {
			t.noShouldTraverse = !hasUnreachable
		}
	}
}

// Java: CheckUnreachableCode.exitScopeWithCfg(NodeTraversal)
func (c *CheckUnreachableCode) ExitScopeWithCfg(t *NodeTraversal) {
	if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) && len(c.hasUnreachableStack) > 0 {
		c.hasUnreachableStack = c.hasUnreachableStack[:len(c.hasUnreachableStack)-1]
		if t != nil && t.callback == c {
			t.noShouldTraverse = c.CurrentScopeHasNoUnreachable()
		}
	}
}

// Java: CheckUnreachableCode.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckUnreachableCode) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
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
func (c *CheckUnreachableCode) InitScope(controlFlowGraph ControlFlowGraph[rhino.Node]) {
	c.reachability.graph = controlFlowGraph
	c.reachability.edgePredicate = CheckUnreachableCode_IsReachable
	c.reachability.Compute(controlFlowGraph.GetEntry().GetValue())
}

// Java: CheckUnreachableCode.isReachable(EdgeTuple<Node, Branch>)
func CheckUnreachableCode_IsReachable(input *GraphReachability_EdgeTuple[rhino.Node, ControlFlowGraph_Branch]) bool {
	branch := input.edge
	if !branch.IsConditional() {
		return true
	}
	predecessor := input.sourceNode
	condition := NodeUtil_GetConditionExpression(predecessor)

	// TODO(user): Handle more complicated expression like true == true,
	// etc....
	if !condition.IsNil() {
		val := NodeUtil_GetBooleanValue(condition)
		if val != base.Tri_UNKNOWN {
			return val.ToBoolean(true) == (branch == ControlFlowGraph_Branch_ON_TRUE)
		}
	}
	return true
}

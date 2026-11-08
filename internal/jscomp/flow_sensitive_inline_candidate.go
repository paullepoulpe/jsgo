// Port of com/google/javascript/jscomp/FlowSensitiveInlineVariables.java @ closure-compiler v20261104

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// hasSideEffectsInSubtree reports whether the intra-CFG subtree rooted at n
// contains a side effect (or, when checkDeps is true, mutates a property or a
// variable that defDeps depends on).
//
// Java: FlowSensitiveInlineVariables.SideEffectPredicate.apply(Node)
func (f *FlowSensitiveInlineVariables) hasSideEffectsInSubtree(
	n rhino.Node, defDeps *MustBeReachingVariableDef_Definition, checkDeps bool,
) bool {
	return walkCfgNode(n, func(cur, _ rhino.Node) bool {
		if checkDeps &&
			cur.IsName() &&
			f.reachingDef.hasDependentNameID(defDeps, cur.GetStringID()) &&
			NodeUtil_IsLValue(cur) {
			// the name is being written to. this is a problem, unless it is part of a top-level assign
			// chain and the write will take place after all CFG node subexpressions are evaluated
			return !fsivIsTopLevelAssignTarget(cur)
		}

		astAnalyzer := f.getAstAnalyzer()
		if (cur.IsCall() || cur.IsOptChainCall() || cur.IsTaggedTemplateLit()) &&
			astAnalyzer.FunctionCallHasSideEffects(cur) {
			return true
		}
		if cur.IsNew() && astAnalyzer.ConstructorCallHasSideEffects(cur) {
			return true
		}
		if cur.IsDelProp() {
			return true
		}
		if (cur.IsGetProp() || cur.IsGetElem()) && NodeUtil_IsLValue(cur) {
			return !checkDeps || !fsivIsTopLevelAssignTarget(cur)
		}
		return false
	})
}

func (f *FlowSensitiveInlineVariables) hasNodeSideEffects(n rhino.Node) bool {
	return f.hasSideEffectsInSubtree(n, nil, false)
}

// fsivIsTopLevelAssignTarget reports whether the given node is the target of a
// (possibly chained) assignment statement.
//
// Java: FlowSensitiveInlineVariables.isTopLevelAssignTarget(Node)
func fsivIsTopLevelAssignTarget(n rhino.Node) bool {
	ancestor := n.GetParent()
	for ancestor.IsAssign() {
		ancestor = ancestor.GetParent()
	}
	return ancestor.IsExprResult()
}

// flowSensitiveInlineCandidate models the connection between a definition and
// a use of that definition.
type flowSensitiveInlineCandidate struct {
	pass *FlowSensitiveInlineVariables

	// Pool ID of the variable name.
	nameID uint32

	// Nodes related to the definition.
	def         rhino.Node
	defMetadata *MustBeReachingVariableDef_Definition

	// Nodes related to the use.
	use        rhino.Node
	useCfgNode rhino.Node

	// Number of uses of the variable within the current CFG node.
	numUsesWithinCfgNode int

	// Go-only: marks candidates retired or invalidated during the value-slice drain loop.
	removed bool
}

// Java: FlowSensitiveInlineVariables.Candidate.getDefCfgNode()
func (c *flowSensitiveInlineCandidate) getDefCfgNode() rhino.Node {
	return c.defMetadata.node
}

// Java: FlowSensitiveInlineVariables.Candidate.canInline(Scope)
func (c *flowSensitiveInlineCandidate) canInline(scope *Scope) bool {
	defCfgNode := c.getDefCfgNode()
	// Cannot inline a parameter, or a definition/use detached by an earlier candidate.
	if !defCfgNode.IsNil() && (defCfgNode.IsFunction() || defCfgNode.GetParent().IsNil()) {
		return false
	}
	if c.use.IsNil() || c.use.GetParent().IsNil() ||
		(!c.useCfgNode.IsNil() && c.useCfgNode.GetParent().IsNil()) {
		return false
	}

	// Make sure that the name is not within a loop
	if NodeUtil_IsWithinLoop(c.use) {
		return false
	}

	c.getDefinition(defCfgNode)

	// Definition was not found, or its initializer/RHS was already detached.
	if c.def.IsNil() || c.def.GetParent().IsNil() {
		return false
	}
	rhs, hasRhs := c.defValueOpt()
	if !hasRhs {
		return false
	}

	// Check that the assignment isn't used as a R-Value.
	if c.def.IsAssign() && ast.RoleOf(c.def) != ast.RoleExprStmtExpr {
		return false
	}

	if !c.isRhsNodeSafeToInline(rhs.Node(), scope) {
		return false
	}

	if c.pass.getAstAnalyzer().MayHaveSideEffects(rhs.Node()) {
		return false
	}

	// Make sure that there are no more than one use in the program and in
	// the CFG node. Even when semantically valid, inlining twice increases
	// code size.
	c.getNumUseInUseCfgNode(c.useCfgNode)
	if c.numUsesWithinCfgNode != 1 {
		return false
	}

	// A subexpression evaluated after the variable has a side effect.
	// Example, for x:
	// x = readProp(b), modifyProp(b); print(x);
	if c.checkPostExpressions(c.def, defCfgNode) {
		return false
	}

	// Similar check as the above but this time, all the sub-expressions
	// evaluated before the variable.
	// x = readProp(b); modifyProp(b), print(x);
	if c.checkPreExpressions(c.use, c.useCfgNode) {
		return false
	}

	if !c.pass.getReachingUses().HasExactlyOneUseByID(c.nameID, defCfgNode) {
		return false
	}

	// We can skip the side effect check along the paths of two nodes if
	// they are just next to each other.
	if NodeUtil_IsStatementBlock(defCfgNode.GetParent()) &&
		defCfgNode.GetNext() != c.useCfgNode {
		// Similar side effect check as above but this time the side effect is
		// elsewhere along the path.
		// x = readProp(b); while(modifyProp(b)) {}; print(x);
		cfg := c.pass.cfg
		c.pass.pathChecker.ResetLinkedAllEdges(
			&cfg.LinkedDirectedGraph,
			cfg.Nodes.Get(defCfgNode),
			cfg.Nodes.Get(c.useCfgNode),
			c.pass.sideEffectPredicateFunc(),
			false)
		if c.pass.pathChecker.SomePathsSatisfyPredicate() {
			return false
		}
	}

	return true
}

// inlineVariable performs the actual transformation.
//
// Java: FlowSensitiveInlineVariables.Candidate.inlineVariable()
func (c *flowSensitiveInlineCandidate) inlineVariable() {
	defParent := c.def.GetParent()
	useParent := c.use.GetParent()
	if assign, ok := ast.AsAssign(c.def); ok {
		rhs := assign.Value().Detach()
		jutil.CheckState(defParent.IsExprResult())
		for defParent.GetParent().IsLabel() {
			defParent = defParent.GetParent()
		}
		c.pass.reportChange(defParent)
		defParent.Detach()
		c.use.ReplaceWith(rhs.Node())
	} else if decl, ok := ast.AsNameDeclaration(defParent); ok {
		rhs := c.defValue()
		if decl.Token() == rhino.Token_CONST {
			// If it is a const var we don't want to remove the rhs of the variable
			rhs.ReplaceWith(ast.NewName("undefined").AsExpr())
			c.use.ReplaceWith(rhs.Node())
		} else {
			rhs.Detach()
			c.use.ReplaceWith(rhs.Node())
		}
	} else {
		panic(jutil.NewIllegalStateExceptionMsg("No other definitions can be inlined."))
	}
	c.pass.reportChange(useParent)
}

// isDefOfVar reports whether n defines the candidate's variable: the
// variable declared with a value, or an assignment to it.
func (c *flowSensitiveInlineCandidate) isDefOfVar(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		return n.GetStringID() == c.nameID && n.HasChildren()
	case rhino.Token_ASSIGN:
		assign := ast.MustAssign(n)
		lhs, ok := assign.Target().AsName()
		return ok && lhs.Node().GetStringID() == c.nameID && !assign.Value().Node().IsNil()
	}
	return false
}

// defValueOpt returns the value that the definition assigns, or false if the
// definition's initializer or RHS has already been detached.
func (c *flowSensitiveInlineCandidate) defValueOpt() (ast.Expr, bool) {
	if assign, ok := ast.AsAssign(c.def); ok {
		v := assign.Value()
		return v, !v.Node().IsNil()
	}
	if name, ok := ast.AsName(c.def); ok {
		return name.Init()
	}
	return ast.Expr{}, false
}

// defValue returns the value that the definition assigns: the value of an
// assignment, or the initial value of a declared name.
func (c *flowSensitiveInlineCandidate) defValue() ast.Expr {
	value, ok := c.defValueOpt()
	jutil.CheckState(ok)
	return value
}

// getDefinition sets the def node.
//
// n is a node that has a corresponding CFG node in the CFG.
//
// Java: FlowSensitiveInlineVariables.Candidate.getDefinition(Node)
func (c *flowSensitiveInlineCandidate) getDefinition(n rhino.Node) {
	jutil.CheckNotNull(c.pass.compiler)
	if n.IsNil() {
		throwNilCfgNodeTraversal(c.pass.compiler)
		return
	}
	walkCfgEvalOrder(n, func(cur, _ rhino.Node) {
		if c.isDefOfVar(cur) {
			c.def = cur
		}
	})
}

func fsivIsAssignChain(child rhino.Node, ancestor rhino.Node) bool {
	for n := child; n != ancestor; n = n.GetParent() {
		if !n.IsAssign() {
			return false
		}
	}
	return true
}

// getNumUseInUseCfgNode computes the number of uses of the variable varName
// and stores it in numUsesWithinCfgNode.
//
// Java: FlowSensitiveInlineVariables.Candidate.getNumUseInUseCfgNode(Node)
func (c *flowSensitiveInlineCandidate) getNumUseInUseCfgNode(cfgNode rhino.Node) {
	c.numUsesWithinCfgNode = 0
	jutil.CheckNotNull(c.pass.compiler)
	if cfgNode.IsNil() {
		throwNilCfgNodeTraversal(c.pass.compiler)
		return
	}
	walkCfgNode(cfgNode, func(n, parent rhino.Node) bool {
		if n.IsName() && n.GetStringID() == c.nameID {
			if !(ast.RoleOf(n) == ast.RoleAssignTarget && fsivIsAssignChain(parent, cfgNode)) {
				c.numUsesWithinCfgNode++
			}
		}
		return false
	})
}

func (c *flowSensitiveInlineCandidate) hasOutOfScopeName(n rhino.Node, usageScope *Scope) bool {
	if n.IsName() {
		id := n.GetStringID()
		if id > strpool.EmptyID && usageScope.GetVarByID(id) == nil {
			if c.pass == nil || !c.pass.hasEnclosingSlotByID(id) {
				return true
			}
		}
	}
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if c.hasOutOfScopeName(child, usageScope) {
			return true
		}
	}
	return false
}

// isRhsSafeToInline checks if the definition we're considering inlining has
// anything that makes inlining unsafe (that hasn't already been caught).
//
// Java: FlowSensitiveInlineVariables.Candidate.isRhsSafeToInline(Scope)
func (c *flowSensitiveInlineCandidate) isRhsSafeToInline(usageScope *Scope) bool {
	rhs, ok := c.defValueOpt()
	return ok && c.isRhsNodeSafeToInline(rhs.Node(), usageScope)
}

func (c *flowSensitiveInlineCandidate) isRhsNodeSafeToInline(rhs rhino.Node, usageScope *Scope) bool {
	if unsafe, _ := fsivInspectRhs(rhs); unsafe {
		return false
	}
	return !c.hasOutOfScopeName(rhs, usageScope)
}

func (c *flowSensitiveInlineCandidate) hasPeerSideEffects(expr rhino.Node) bool {
	return c.pass.hasSideEffectsInSubtree(expr, c.defMetadata, true)
}

func fsivEvalPeerAfterState(c *flowSensitiveInlineCandidate, expr rhino.Node) cfgCondEvalState {
	if c.hasPeerSideEffects(expr) {
		return cfgCondEvalRead
	}
	return cfgCondEvalContinue
}

func (c *flowSensitiveInlineCandidate) checkPostExpressions(n rhino.Node, expressionRoot rhino.Node) bool {
	for p := n; p != expressionRoot; {
		parent := p.GetParent()
		if cm, ok := ast.AsComputedMember(parent); ok && parent.GetParent().IsClassMembers() {
			if parent == expressionRoot {
				if cfd, isField := cm.AsComputedFieldDef(); isField && cfd.IsStatic() && cfd.Key().Is(p) {
					if val, hasVal := cfd.Value(); hasVal && c.hasPeerSideEffects(val.Node()) {
						return true
					}
				}
				return false
			}
			classMembers := parent.GetParent()
			inPass2 := !cm.Key().Is(p)
			if (!inPass2 || cm.IsStatic()) &&
				walkClassMembersAfterState(classMembers, parent, inPass2, c, fsivEvalPeerAfterState) != cfgCondEvalContinue {
				return true
			}
			p = classMembers
			continue
		}
		if parent.IsClassMembers() {
			if p.IsMemberFieldDef() && p.IsStaticMember() &&
				walkClassMembersAfterState(parent, p, true, c, fsivEvalPeerAfterState) != cfgCondEvalContinue {
				return true
			}
			p = parent
			continue
		}
		for cur := p.GetNext(); !cur.IsNil(); cur = cur.GetNext() {
			if c.hasPeerSideEffects(cur) {
				return true
			}
		}
		p = parent
	}
	return false
}

func (c *flowSensitiveInlineCandidate) checkPreExpressions(n rhino.Node, expressionRoot rhino.Node) bool {
	for p := n; p != expressionRoot; {
		parent := p.GetParent()
		if cm, ok := ast.AsComputedMember(parent); ok && parent.GetParent().IsClassMembers() {
			if parent == expressionRoot {
				if cfd, isField := cm.AsComputedFieldDef(); isField && cfd.IsStatic() && !cfd.Key().Is(p) {
					if key := cfd.Key().Node(); !key.IsNil() && c.hasPeerSideEffects(key) {
						return true
					}
				}
				return false
			}
			classMembers := parent.GetParent()
			inPass2 := !cm.Key().Is(p)
			if (!inPass2 || cm.IsStatic()) &&
				walkClassMembersBefore(classMembers, parent, inPass2, c, (*flowSensitiveInlineCandidate).hasPeerSideEffects) {
				return true
			}
			p = classMembers
			continue
		}
		if parent.IsClassMembers() {
			if p.IsMemberFieldDef() && p.IsStaticMember() &&
				walkClassMembersBefore(parent, p, true, c, (*flowSensitiveInlineCandidate).hasPeerSideEffects) {
				return true
			}
			p = parent
			continue
		}
		oldestSibling := parent.GetFirstChild()
		if oldestSibling.IsDestructuringPattern() {
			if p.IsDestructuringPattern() {
				if !p.GetNext().IsNil() && c.hasPeerSideEffects(p.GetNext()) {
					return true
				}
			}
			p = parent
			continue
		}
		for cur := oldestSibling; cur != p; cur = cur.GetNext() {
			if c.hasPeerSideEffects(cur) {
				return true
			}
		}
		p = parent
	}
	return false
}

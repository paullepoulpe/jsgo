// Port of com/google/javascript/jscomp/FlowSensitiveInlineVariables.java @ closure-compiler v20261021

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
	"iter"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// FlowSensitiveInlineVariables_SideEffectPredicate is a Predicate<Node>.
type FlowSensitiveInlineVariables_SideEffectPredicate struct {
	outer *FlowSensitiveInlineVariables

	// Check if there are side effects affecting the value of any of these
	// names (but not properties defined on that name)
	namesToCheck *jutil.LinkedHashSet[string]
}

// Java: FlowSensitiveInlineVariables.SideEffectPredicate()
func NewFlowSensitiveInlineVariables_SideEffectPredicate(outer *FlowSensitiveInlineVariables) *FlowSensitiveInlineVariables_SideEffectPredicate {
	return &FlowSensitiveInlineVariables_SideEffectPredicate{outer: outer, namesToCheck: nil}
}

// Java: FlowSensitiveInlineVariables.SideEffectPredicate(Set<String>)
func NewFlowSensitiveInlineVariables_SideEffectPredicateWithNames(outer *FlowSensitiveInlineVariables, names *jutil.LinkedHashSet[string]) *FlowSensitiveInlineVariables_SideEffectPredicate {
	return &FlowSensitiveInlineVariables_SideEffectPredicate{outer: outer, namesToCheck: names}
}

// Java: FlowSensitiveInlineVariables.SideEffectPredicate.apply(Node)
func (p *FlowSensitiveInlineVariables_SideEffectPredicate) Apply(n rhino.Node) bool {
	// When the node is null it means, we reached the implicit return
	// where the function returns (possibly without a return statement)
	if n.IsNil() {
		return false
	}

	if p.namesToCheck != nil &&
		n.IsName() &&
		p.namesToCheck.Contains(n.GetString()) &&
		NodeUtil_IsLValue(n) {
		// the name is being written to. this is a problem, unless it is part of a top-level assign
		// chain and the write will take place after all CFG node subexpressions are evaluated
		return !FlowSensitiveInlineVariables_IsTopLevelAssignTarget(n)
	}

	astAnalyzer := p.outer.getAstAnalyzer()
	if (n.IsCall() || n.IsOptChainCall() || n.IsTaggedTemplateLit()) &&
		astAnalyzer.FunctionCallHasSideEffects(n) {
		return true
	}

	if n.IsNew() && astAnalyzer.ConstructorCallHasSideEffects(n) {
		return true
	}

	if n.IsDelProp() {
		return true
	}

	if (n.IsGetProp() || n.IsGetElem()) && NodeUtil_IsLValue(n) {
		return p.namesToCheck == nil || !FlowSensitiveInlineVariables_IsTopLevelAssignTarget(n)
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !ControlFlowGraph_IsEnteringNewCfgNode(c) && p.Apply(c) {
			return true
		}
	}
	return false
}

// FlowSensitiveInlineVariables_IsTopLevelAssignTarget reports whether the
// given node is the target of a (possibly chained) assignment.
//
// Java: FlowSensitiveInlineVariables.isTopLevelAssignTarget(Node)
func FlowSensitiveInlineVariables_IsTopLevelAssignTarget(n rhino.Node) bool {
	ancestor := n.GetParent()
	for ancestor.IsAssign() {
		ancestor = ancestor.GetParent()
	}
	return ancestor.IsExprResult()
}

// FlowSensitiveInlineVariables_Candidate models the connection between a
// definition and a use of that definition.
type FlowSensitiveInlineVariables_Candidate struct {
	outer *FlowSensitiveInlineVariables

	// Name of the variable.
	varName string

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

// Java: FlowSensitiveInlineVariables.Candidate(String, Definition, Node, Node)
func NewFlowSensitiveInlineVariables_Candidate(
	outer *FlowSensitiveInlineVariables,
	varName string,
	defMetadata *MustBeReachingVariableDef_Definition,
	use rhino.Node,
	useCfgNode rhino.Node,
) *FlowSensitiveInlineVariables_Candidate {
	jutil.CheckArgument(use.IsName())
	return &FlowSensitiveInlineVariables_Candidate{
		outer:       outer,
		varName:     varName,
		defMetadata: defMetadata,
		use:         use,
		useCfgNode:  useCfgNode,
	}
}

// Java: FlowSensitiveInlineVariables.Candidate.getDefCfgNode()
func (c *FlowSensitiveInlineVariables_Candidate) GetDefCfgNode() rhino.Node {
	return c.defMetadata.node
}

// Java: FlowSensitiveInlineVariables.Candidate.canInline(Scope)
func (c *FlowSensitiveInlineVariables_Candidate) CanInline(scope *Scope) bool {
	// Cannot inline a parameter.
	if c.GetDefCfgNode().IsFunction() {
		return false
	}

	// Make sure that the name is not within a loop
	if NodeUtil_IsWithinLoop(c.use) {
		return false
	}

	c.GetDefinition(c.GetDefCfgNode())

	// Definition was not found.
	if c.def.IsNil() {
		return false
	}

	// Check that the assignment isn't used as a R-Value.
	if c.def.IsAssign() && ast.RoleOf(c.def) != ast.RoleExprStmtExpr {
		return false
	}

	if !c.IsRhsSafeToInline(scope) {
		return false
	}

	if c.outer.getAstAnalyzer().MayHaveSideEffects(c.defValue().Node()) {
		return false
	}

	// Make sure that there are no more than one use in the program and in
	// the CFG node. Even when semantically valid, inlining twice increases
	// code size.
	c.GetNumUseInUseCfgNode(c.useCfgNode)
	if c.numUsesWithinCfgNode != 1 {
		return false
	}

	namesToCheck := &c.outer.namesToCheckSet
	namesToCheck.Clear()
	c.outer.reachingDef.AddDependencyNames(c.defMetadata, namesToCheck)

	sideEffectPredicateWithNames := FlowSensitiveInlineVariables_SideEffectPredicate{
		outer:        c.outer,
		namesToCheck: namesToCheck,
	}

	// A subexpression evaluated after the variable has a side effect.
	// Example, for x:
	// x = readProp(b), modifyProp(b); print(x);
	if fsivCheckPostExpressions(c.def, c.GetDefCfgNode(), &sideEffectPredicateWithNames) {
		return false
	}

	// Similar check as the above but this time, all the sub-expressions
	// evaluated before the variable.
	// x = readProp(b); modifyProp(b), print(x);
	if fsivCheckPreExpressions(c.use, c.useCfgNode, &sideEffectPredicateWithNames) {
		return false
	}

	if !c.outer.getReachingUses().HasExactlyOneUse(c.varName, c.GetDefCfgNode()) {
		return false
	}

	// We can skip the side effect check along the paths of two nodes if
	// they are just next to each other.
	if NodeUtil_IsStatementBlock(c.GetDefCfgNode().GetParent()) &&
		c.GetDefCfgNode().GetNext() != c.useCfgNode {
		// Similar side effect check as above but this time the side effect is
		// elsewhere along the path.
		// x = readProp(b); while(modifyProp(b)) {}; print(x);
		cfg := c.outer.cfg
		c.outer.pathChecker.ResetLinkedAllEdges(
			&cfg.LinkedDirectedGraph,
			cfg.Nodes.Get(c.GetDefCfgNode()),
			cfg.Nodes.Get(c.useCfgNode),
			c.outer.sideEffectPredicateFn,
			false)
		if c.outer.pathChecker.SomePathsSatisfyPredicate() {
			return false
		}
	}

	return true
}

// Java: FlowSensitiveInlineVariables.Candidate.hasExactlyOne(Iterable<Node>)
func (c *FlowSensitiveInlineVariables_Candidate) HasExactlyOne(iterable iter.Seq[rhino.Node]) bool {
	count := 0
	for range iterable {
		count++
		if count > 1 {
			break
		}
	}
	return count == 1
}

// InlineVariable performs the actual transformation.
//
// Java: FlowSensitiveInlineVariables.Candidate.inlineVariable()
func (c *FlowSensitiveInlineVariables_Candidate) InlineVariable() {
	defParent := c.def.GetParent()
	useParent := c.use.GetParent()
	if assign, ok := ast.AsAssign(c.def); ok {
		rhs := assign.Value().Detach()
		jutil.CheckState(defParent.IsExprResult())
		for defParent.GetParent().IsLabel() {
			defParent = defParent.GetParent()
		}
		c.outer.reportChange(defParent)
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
	c.outer.reportChange(useParent)
}

func (c *FlowSensitiveInlineVariables_Candidate) findDefWalk(n, parent rhino.Node) {
	if !parent.IsNil() && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return
	}
	switch n.GetToken() {
	case rhino.Token_CLASS:
		// In NodeTraversal's order: the extends clause, then the name
		// (which defines nothing), then the body.
		class := ast.MustClass(n)
		if extends, ok := class.Extends(); ok {
			c.findDefWalk(extends.Node(), n)
		}
		c.findDefWalk(class.Body().Node(), n)
	case rhino.Token_CLASS_MEMBERS:
		fsivWalkClassMembers(ast.MustClassBody(n), c.findDefWalk)
	default:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			c.findDefWalk(child, n)
		}
	}
	if c.isDefOfVar(n) {
		c.def = n
	}
}

// isDefOfVar reports whether n defines the candidate's variable: the
// variable declared with a value, or an assignment to it.
func (c *FlowSensitiveInlineVariables_Candidate) isDefOfVar(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		return n.GetString() == c.varName && n.HasChildren()
	case rhino.Token_ASSIGN:
		lhs, ok := ast.MustAssign(n).Target().AsName()
		return ok && lhs.Ident() == c.varName
	}
	return false
}

// defValue returns the value that the definition assigns: the value of an
// assignment, or the initial value of a declared name.
func (c *FlowSensitiveInlineVariables_Candidate) defValue() ast.Expr {
	if assign, ok := ast.AsAssign(c.def); ok {
		return assign.Value()
	}
	value, ok := ast.MustName(c.def).Init()
	jutil.CheckState(ok)
	return value
}

// fsivWalkClassMembers walks the parts of a class body that belong to the
// CFG node of the class, in NodeTraversal's order: the computed keys
// first, then the values of computed members (unless they start CFG nodes
// of their own) and the other members.
func fsivWalkClassMembers(body ast.ClassBody, walk func(n, parent rhino.Node)) {
	for m := range body.Members() {
		switch m.Token() {
		case rhino.Token_COMPUTED_PROP:
			walk(m.MustComputedProp().Key().Node(), m.Node())
		case rhino.Token_COMPUTED_FIELD_DEF:
			walk(m.MustComputedFieldDef().Key().Node(), m.Node())
		}
	}
	for m := range body.Members() {
		switch m.Token() {
		case rhino.Token_COMPUTED_PROP:
			if value, ok := m.MustComputedProp().Value(); ok && !ControlFlowGraph_IsEnteringNewCfgNode(m.Node()) {
				walk(value.Node(), m.Node())
			}
		case rhino.Token_COMPUTED_FIELD_DEF:
			if value, ok := m.MustComputedFieldDef().Value(); ok && !ControlFlowGraph_IsEnteringNewCfgNode(m.Node()) {
				walk(value.Node(), m.Node())
			}
		default:
			walk(m.Node(), body.Node())
		}
	}
}

// GetDefinition sets the def node.
//
// n is a node that has a corresponding CFG node in the CFG.
//
// Java: FlowSensitiveInlineVariables.Candidate.getDefinition(Node)
func (c *FlowSensitiveInlineVariables_Candidate) GetDefinition(n rhino.Node) {
	jutil.CheckNotNull(c.outer.compiler)
	if n.IsNil() {
		gatherCb := flowSensitiveInlineVariables_definitionCallback{candidate: c}
		traverseCfgNode(c.outer.compiler, n, &gatherCb)
		return
	}
	c.findDefWalk(n, rhino.NilNode)
}

// flowSensitiveInlineVariables_definitionCallback is the anonymous
// AbstractCfgNodeTraversalCallback of Candidate.getDefinition.
type flowSensitiveInlineVariables_definitionCallback struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	candidate *FlowSensitiveInlineVariables_Candidate
}

// Java: FlowSensitiveInlineVariables.Candidate.getDefinition(...).new AbstractCfgNodeTraversalCallback() {...}.visit(NodeTraversal, Node, Node)
func (cb *flowSensitiveInlineVariables_definitionCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if cb.candidate.isDefOfVar(n) {
		cb.candidate.def = n
	}
}

func (c *FlowSensitiveInlineVariables_Candidate) countUseWalk(n, parent, cfgNode rhino.Node) {
	if !parent.IsNil() && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return
	}
	if n.GetToken() == rhino.Token_CLASS_MEMBERS {
		fsivWalkClassMembers(ast.MustClassBody(n), func(n, parent rhino.Node) { c.countUseWalk(n, parent, cfgNode) })
		return
	}
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		c.countUseWalk(child, n, cfgNode)
	}
	if n.IsName() && n.GetString() == c.varName {
		if ast.RoleOf(n) == ast.RoleAssignTarget && fsivIsAssignChain(parent, cfgNode) {
			return
		}
		c.numUsesWithinCfgNode++
	}
}

func fsivIsAssignChain(child rhino.Node, ancestor rhino.Node) bool {
	for n := child; n != ancestor; n = n.GetParent() {
		if !n.IsAssign() {
			return false
		}
	}
	return true
}

// GetNumUseInUseCfgNode computes the number of uses of the variable varName
// and stores it in numUsesWithinCfgNode.
//
// Java: FlowSensitiveInlineVariables.Candidate.getNumUseInUseCfgNode(Node)
func (c *FlowSensitiveInlineVariables_Candidate) GetNumUseInUseCfgNode(cfgNode rhino.Node) {
	c.numUsesWithinCfgNode = 0
	jutil.CheckNotNull(c.outer.compiler)
	if cfgNode.IsNil() {
		gatherCb := flowSensitiveInlineVariables_numUseCallback{candidate: c, cfgNode: cfgNode}
		traverseCfgNode(c.outer.compiler, cfgNode, &gatherCb)
		return
	}
	c.countUseWalk(cfgNode, rhino.NilNode, cfgNode)
}

// flowSensitiveInlineVariables_numUseCallback is the anonymous
// AbstractCfgNodeTraversalCallback of Candidate.getNumUseInUseCfgNode.
type flowSensitiveInlineVariables_numUseCallback struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	candidate *FlowSensitiveInlineVariables_Candidate
	cfgNode   rhino.Node
}

// Java: FlowSensitiveInlineVariables.Candidate.getNumUseInUseCfgNode(...).new AbstractCfgNodeTraversalCallback() {...}.visit(NodeTraversal, Node, Node)
func (cb *flowSensitiveInlineVariables_numUseCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsName() && n.GetString() == cb.candidate.varName {
		if ast.RoleOf(n) == ast.RoleAssignTarget && cb.IsAssignChain(parent, cb.cfgNode) {
			return
		}
		cb.candidate.numUsesWithinCfgNode++
	}
}

// Java: FlowSensitiveInlineVariables.Candidate.getNumUseInUseCfgNode(...).new AbstractCfgNodeTraversalCallback() {...}.isAssignChain(Node, Node)
func (cb *flowSensitiveInlineVariables_numUseCallback) IsAssignChain(child rhino.Node, ancestor rhino.Node) bool {
	return fsivIsAssignChain(child, ancestor)
}

func fsivHasUnsafeRhsNode(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_GETELEM,
		rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETPROP,
		rhino.Token_OPTCHAIN_GETELEM,
		rhino.Token_CLASS,
		rhino.Token_ARRAYLIT,
		rhino.Token_OBJECTLIT,
		rhino.Token_REGEXP,
		rhino.Token_NEW:
		return true
	case rhino.Token_FUNCTION:
		return false
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if fsivHasUnsafeRhsNode(c) {
			return true
		}
	}
	return false
}

func (c *FlowSensitiveInlineVariables_Candidate) hasOutOfScopeName(n rhino.Node, usageScope *Scope) bool {
	if n.IsName() {
		name := n.GetString()
		if name != "" && !usageScope.HasSlot(name) {
			if c.outer == nil || !c.outer.hasEnclosingSlot(name) {
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

// IsRhsSafeToInline checks if the definition we're considering inlining has
// anything that makes inlining unsafe (that hasn't already been caught).
//
// Java: FlowSensitiveInlineVariables.Candidate.isRhsSafeToInline(Scope)
func (c *FlowSensitiveInlineVariables_Candidate) IsRhsSafeToInline(usageScope *Scope) bool {
	rhs := c.defValue().Node()
	if fsivHasUnsafeRhsNode(rhs) {
		return false
	}
	if c.hasOutOfScopeName(rhs, usageScope) {
		return false
	}
	return true
}

func fsivCheckPostExpressions(
	n rhino.Node, expressionRoot rhino.Node, predicate *FlowSensitiveInlineVariables_SideEffectPredicate,
) bool {
	for p := n; p != expressionRoot; p = p.GetParent() {
		for cur := p.GetNext(); !cur.IsNil(); cur = cur.GetNext() {
			if predicate.Apply(cur) {
				return true
			}
		}
	}
	return false
}

func fsivCheckPreExpressions(
	n rhino.Node, expressionRoot rhino.Node, predicate *FlowSensitiveInlineVariables_SideEffectPredicate,
) bool {
	for p := n; p != expressionRoot; p = p.GetParent() {
		oldestSibling := p.GetParent().GetFirstChild()
		if oldestSibling.IsDestructuringPattern() {
			if p.IsDestructuringPattern() {
				if !p.GetNext().IsNil() && predicate.Apply(p.GetNext()) {
					return true
				}
			}
			continue
		}
		for cur := oldestSibling; cur != p; cur = cur.GetNext() {
			if predicate.Apply(cur) {
				return true
			}
		}
	}
	return false
}

// FlowSensitiveInlineVariables_CheckPostExpressions returns true, given an
// expression by its root and sub-expression n, if the predicate is true for
// some expression evaluated after n.
//
// Java: FlowSensitiveInlineVariables.checkPostExpressions(Node, Node, Predicate<Node>)
func FlowSensitiveInlineVariables_CheckPostExpressions(
	n rhino.Node, expressionRoot rhino.Node, predicate func(rhino.Node) bool,
) bool {
	for p := n; p != expressionRoot; p = p.GetParent() {
		for cur := p.GetNext(); !cur.IsNil(); cur = cur.GetNext() {
			if predicate(cur) {
				return true
			}
		}
	}
	return false
}

// FlowSensitiveInlineVariables_CheckPreExpressions returns true, given an
// expression by its root and sub-expression n, if the predicate is true for
// some expression evaluated before n.
//
// Java: FlowSensitiveInlineVariables.checkPreExpressions(Node, Node, Predicate<Node>)
func FlowSensitiveInlineVariables_CheckPreExpressions(
	n rhino.Node, expressionRoot rhino.Node, predicate func(rhino.Node) bool,
) bool {
	for p := n; p != expressionRoot; p = p.GetParent() {
		oldestSibling := p.GetParent().GetFirstChild()
		// Evaluate a destructuring assignment right-to-left.
		if oldestSibling.IsDestructuringPattern() {
			if p.IsDestructuringPattern() {
				if !p.GetNext().IsNil() && predicate(p.GetNext()) {
					return true
				}
			}
			continue
		}
		for cur := oldestSibling; cur != p; cur = cur.GetNext() {
			if predicate(cur) {
				return true
			}
		}
	}
	return false
}

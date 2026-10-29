// Port of com/google/javascript/jscomp/ControlFlowAnalysis.java @ closure-compiler v20261021

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
	"iter"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type cfaFinallyEntry struct {
	jump     rhino.Node
	head     [4]rhino.Node
	nHead    uint8
	overflow []rhino.Node
}

func (e *cfaFinallyEntry) addTarget(target rhino.Node) {
	for i := uint8(0); i < e.nHead; i++ {
		if e.head[i] == target {
			return
		}
	}
	if int(e.nHead) < len(e.head) {
		e.head[e.nHead] = target
		e.nHead++
		return
	}
	for _, existing := range e.overflow {
		if existing == target {
			return
		}
	}
	e.overflow = append(e.overflow, target)
}

type cfaFinallyTargets struct {
	entry *cfaFinallyEntry
}

func (t cfaFinallyTargets) Size() int {
	if t.entry == nil {
		return 0
	}
	return int(t.entry.nHead) + len(t.entry.overflow)
}

func (t cfaFinallyTargets) All() iter.Seq[rhino.Node] {
	return func(yield func(rhino.Node) bool) {
		if t.entry == nil {
			return
		}
		for i := uint8(0); i < t.entry.nHead; i++ {
			if !yield(t.entry.head[i]) {
				return
			}
		}
		for _, n := range t.entry.overflow {
			if !yield(n) {
				return
			}
		}
	}
}

type cfaFinallyMap struct {
	entries []cfaFinallyEntry
}

func (m *cfaFinallyMap) reset() {
	if len(m.entries) == 0 {
		return
	}
	if cap(m.entries) > cfaReuseMaxSize {
		m.entries = nil
		return
	}
	clear(m.entries)
	m.entries = m.entries[:0]
}

func (m *cfaFinallyMap) put(jump, target rhino.Node) {
	for i := range m.entries {
		if m.entries[i].jump == jump {
			m.entries[i].addTarget(target)
			return
		}
	}
	m.entries = append(m.entries, cfaFinallyEntry{
		jump:  jump,
		head:  [4]rhino.Node{0: target},
		nHead: 1,
	})
}

func (m *cfaFinallyMap) find(jump rhino.Node) *cfaFinallyEntry {
	for i := range m.entries {
		if m.entries[i].jump == jump {
			return &m.entries[i]
		}
	}
	return nil
}

func (m *cfaFinallyMap) Get(jump rhino.Node) cfaFinallyTargets {
	return cfaFinallyTargets{entry: m.find(jump)}
}

func (m *cfaFinallyMap) KeySet() iter.Seq[rhino.Node] {
	return func(yield func(rhino.Node) bool) {
		for i := range m.entries {
			if !yield(m.entries[i].jump) {
				return
			}
		}
	}
}

func (c *ControlFlowAnalysis) putFinally(jump rhino.Node, target rhino.Node) {
	c.finallyMap.put(jump, target)
}

// Java: ControlFlowAnalysis.shouldTraverse(NodeTraversal, Node, Node)
func (c *ControlFlowAnalysis) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if c.shouldTraverseIntoChildren(n, parent) {
		// Any AST node that will later have a corresponding CFG node must be in astPosition.
		// To avoid having astPosition grow too large, we exclude AST nodes that are not traversed
		// further, as they usually are not put in the CFG.
		c.astPosition[n] = c.astPositionCounter
		c.astPositionCounter++
		return true
	}
	return false
}

// exceptionHandlerPush is ArrayDeque.push(e).
func (c *ControlFlowAnalysis) exceptionHandlerPush(n rhino.Node) {
	c.exceptionHandler = append(c.exceptionHandler, n)
}

// exceptionHandlerPeek is ArrayDeque.peek(): the top, or null when empty.
func (c *ControlFlowAnalysis) exceptionHandlerPeek() rhino.Node {
	if len(c.exceptionHandler) == 0 {
		return rhino.NilNode
	}
	return c.exceptionHandler[len(c.exceptionHandler)-1]
}

// exceptionHandlerPop is ArrayDeque.pop(): removes the top, or throws
// NoSuchElementException when empty.
func (c *ControlFlowAnalysis) exceptionHandlerPop() rhino.Node {
	if len(c.exceptionHandler) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	top := c.exceptionHandler[len(c.exceptionHandler)-1]
	c.exceptionHandler[len(c.exceptionHandler)-1] = rhino.NilNode
	c.exceptionHandler = c.exceptionHandler[:len(c.exceptionHandler)-1]
	return top
}

// shouldTraverseIntoChildren returns whether the children of this node
// should be traversed as part of this control-flow analysis, i.e. whether
// the control-flow graph being built may require any edges into children of
// this node.
//
// Java: ControlFlowAnalysis.shouldTraverseIntoChildren(Node, Node)
func (c *ControlFlowAnalysis) shouldTraverseIntoChildren(n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		if c.shouldTraverseFunctions || n == c.cfg.GetEntry().GetValue() {
			c.exceptionHandlerPush(n)
			return true
		}
		return false
	case rhino.Token_TRY:
		c.exceptionHandlerPush(n)
		return true
	default:
	}

	/*
	 * We are going to stop the traversal depending on what the node's parent
	 * is.
	 *
	 * We are only interested in adding edges between nodes that change control
	 * flow. The most obvious ones are loops and IF-ELSE's. A statement
	 * transfers control to its next sibling.
	 *
	 * In case of an expression tree, there is no control flow within the tree
	 * even when there are short circuited operators and conditionals. When we
	 * are doing data flow analysis, we will simply synthesize lattices up the
	 * expression tree by finding the meet at each expression node.
	 *
	 * For example: within a Token.SWITCH, the expression in question does not
	 * change the control flow and need not to be considered.
	 */
	if !parent.IsNil() {
		switch parent.GetToken() {
		case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			// Only traverse the body of the for loop.
			if ast.MustLoop(parent).Body().Is(n) {
				return true
			}
			// The control-flow graph contains edges from a FOR node to all its children, even
			// though only the body is actually traversed. So put the other children in astPosition.
			c.astPosition[n] = c.astPositionCounter
			c.astPositionCounter++
			return false
		case rhino.Token_DO:
			// Only traverse the body of the do-while.
			return ast.MustDo(parent).Body().Is(n)
		// Skip conditions, and only traverse the body of the cases. (The
		// head of each statement below, its condition, object, discriminant,
		// case value, catch param or label name, has no node in the graph:
		// the statement's own node stands for it.)
		case rhino.Token_IF:
			return !ast.MustIf(parent).Cond().Is(n)
		case rhino.Token_WHILE:
			return ast.MustWhile(parent).Body().Is(n)
		case rhino.Token_WITH:
			return ast.MustWith(parent).Body().Is(n)
		case rhino.Token_SWITCH:
			return ast.MustSwitch(parent).Body().Is(n)
		case rhino.Token_CASE:
			return ast.MustCase(parent).Body().Is(n)
		case rhino.Token_CATCH:
			return ast.MustCatch(parent).Body().Is(n)
		case rhino.Token_LABEL:
			return ast.MustLabel(parent).Body().Is(n)
		case rhino.Token_FUNCTION:
			return ast.MustFunction(parent).Body().Is(n)
		case rhino.Token_CLASS:
			return c.shouldTraverseFunctions && ast.MustClass(parent).Body().Is(n)
		case rhino.Token_COMPUTED_PROP,
			rhino.Token_CONTINUE,
			rhino.Token_BREAK,
			rhino.Token_EXPR_RESULT,
			rhino.Token_VAR,
			rhino.Token_LET,
			rhino.Token_CONST,
			rhino.Token_EXPORT,
			rhino.Token_IMPORT,
			rhino.Token_RETURN,
			rhino.Token_THROW,
			rhino.Token_MEMBER_FUNCTION_DEF,
			rhino.Token_MEMBER_FIELD_DEF,
			rhino.Token_COMPUTED_FIELD_DEF:
			return false
		case rhino.Token_TRY:
			/* When we are done with the TRY block and there is no FINALLY block,
			 * or done with both the TRY and CATCH block, then no more exceptions
			 * can be handled at this TRY statement, so it can be taken out of the
			 * stack.
			 */
			try := ast.MustTry(parent)
			finally, hasFinally := try.Finally()
			if (!hasFinally && try.Catches().Is(n)) || (hasFinally && finally.Is(n)) {
				jutil.CheckState(c.exceptionHandlerPeek() == parent)
				c.exceptionHandlerPop()
			}
		default:
		}
		// Don't traverse further in an arrow function expression
		if parent.HasParent() && parent.GetParent().IsArrowFunction() && !parent.IsBlock() {
			return false
		}
	}
	return true
}

// cfaJumpLabelID returns the interned label ID of a break or continue.
func cfaJumpLabelID(jump ast.Jump) (uint32, bool) {
	name, ok := jump.Label()
	if !ok {
		return 0, false
	}
	return name.Node().GetStringID(), true
}

// cfaFinallyLeft returns the FINALLY block that runs when a jump leaves cur,
// coming from its child previous: ok is false if cur is not a TRY with a
// FINALLY, or if the jump leaves from the FINALLY itself.
func cfaFinallyLeft(cur, previous rhino.Node) (ast.Block, bool) {
	try, ok := ast.AsTry(cur)
	if !ok {
		return ast.Block{}, false
	}
	finally, ok := try.Finally()
	return finally, ok && !finally.Is(previous)
}

// Java: ControlFlowAnalysis.handleBreak(Node)
func (c *ControlFlowAnalysis) handleBreak(brk ast.Break) {
	node := brk.Node()
	// See if it is a break with label.
	labelID, hasLabel := cfaJumpLabelID(brk.AsJump())
	var cur rhino.Node
	var previous rhino.Node
	var lastJump rhino.Node
	parent := node.GetParent()
	/*
	 * Continuously look up the ancestor tree for the BREAK target or the target
	 * with the corresponding label and connect to it. If along the path we
	 * discover a FINALLY, we will connect the BREAK to that FINALLY. From then
	 * on, we will just record the control flow changes in the finallyMap. This
	 * is due to the fact that we need to connect any node that leaves its own
	 * FINALLY block to the outer FINALLY or the BREAK's target but those nodes
	 * are not known yet due to the way we traverse the nodes.
	 */
	for cur, lastJump = node, node; !cfaIsBreakTargetByID(cur, labelID, hasLabel); cur, parent = parent, parent.GetParent() {
		if finally, ok := cfaFinallyLeft(cur, previous); ok {
			if lastJump == node {
				c.createEdge(lastJump, ControlFlowGraph_Branch_UNCOND, finally.Node())
			} else {
				c.putFinally(lastJump, finally.Node())
			}
			lastJump = cur
		}
		if parent.IsNil() {
			if c.compiler.GetOptions().CanContinueAfterErrors() {
				// In IDE mode, we expect that the data flow graph may
				// not be well-formed.
				return
			} else {
				panic(jutil.NewIllegalStateExceptionMsg("Cannot find break target."))
			}
		}
		previous = cur
	}
	if lastJump == node {
		c.createEdge(lastJump, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFollowNodeWithCfa(cur, c))
	} else {
		c.putFinally(lastJump, ControlFlowAnalysis_ComputeFollowNodeWithCfa(cur, c))
	}
}

// Java: ControlFlowAnalysis.handleContinue(Node)
func (c *ControlFlowAnalysis) handleContinue(cont ast.Continue) {
	node := cont.Node()
	labelID, hasLabel := cfaJumpLabelID(cont.AsJump())
	var cur rhino.Node
	var previous rhino.Node
	var lastJump rhino.Node

	// Similar to handBreak's logic with a few minor variation.
	for cur, lastJump = node, node; !cfaIsContinueTargetByID(cur, labelID, hasLabel); cur = cur.GetParent() {
		if finally, ok := cfaFinallyLeft(cur, previous); ok {
			if lastJump == node {
				c.createEdge(lastJump, ControlFlowGraph_Branch_UNCOND, finally.Node())
			} else {
				c.putFinally(lastJump, finally.Node())
			}
			lastJump = cur
		}
		jutil.CheckStateMsg(cur.HasParent(), "Cannot find continue target.")
		previous = cur
	}
	iter := cur
	if forNode, ok := ast.AsFor(cur); ok {
		// the increment expression happens after the continue
		iter = cfaForUpdate(forNode)
	}

	if lastJump == node {
		c.createEdge(node, ControlFlowGraph_Branch_UNCOND, iter)
	} else {
		c.putFinally(lastJump, iter)
	}
}

// Java: ControlFlowAnalysis.handleReturn(Node)
func (c *ControlFlowAnalysis) handleReturn(ret ast.Return) {
	node := ret.Node()
	var lastJump rhino.Node
	for i := len(c.exceptionHandler) - 1; i >= 0; i-- {
		curHandler := c.exceptionHandler[i]
		if curHandler.IsFunction() {
			break
		}
		if finally, ok := ast.MustTry(curHandler).Finally(); ok {
			if lastJump.IsNil() {
				c.createEdge(node, ControlFlowGraph_Branch_UNCOND, finally.Node())
			} else {
				c.putFinally(lastJump, finally.Node())
			}
			lastJump = curHandler
		}
	}

	if value, ok := ret.Value(); ok {
		c.connectToPossibleExceptionHandler(node, value.Node())
	}

	if lastJump.IsNil() {
		c.createEdge(node, ControlFlowGraph_Branch_UNCOND, rhino.NilNode)
	} else {
		c.putFinally(lastJump, rhino.NilNode)
	}
}

// Java: ControlFlowAnalysis.computeFollowNode(Node, ControlFlowAnalysis)
func ControlFlowAnalysis_ComputeFollowNodeWithCfa(node rhino.Node, cfa *ControlFlowAnalysis) rhino.Node {
	return ControlFlowAnalysis_ComputeFollowNodeFrom(node, node, cfa)
}

// Java: ControlFlowAnalysis.computeFollowNode(Node)
func ControlFlowAnalysis_ComputeFollowNode(node rhino.Node) rhino.Node {
	return ControlFlowAnalysis_ComputeFollowNodeFrom(node, node, nil)
}

// ControlFlowAnalysis_ComputeFollowNodeFrom computes the follow() node of a
// given node and its parent. There is a side effect when calling this
// function. If this function computed an edge that exists a FINALLY, it'll
// attempt to connect the fromNode to the outer FINALLY according to the
// finallyMap.
//
// fromNode is the original source node since node is changed during
// recursion; node is the node that follow() should compute.
//
// Java: ControlFlowAnalysis.computeFollowNode(Node, Node, ControlFlowAnalysis)
func ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode rhino.Node, node rhino.Node, cfa *ControlFlowAnalysis) rhino.Node {
	/*
	 * This is the case where:
	 *
	 * 1. Parent is null implies that we are transferring control to the end of
	 * the script.
	 *
	 * 2. Parent is a function implies that we are transferring control back to
	 * the caller of the function.
	 *
	 * 3. If the node is a return statement, we should also transfer control
	 * back to the caller of the function.
	 *
	 * 4. If the node is root then we have reached the end of what we have been
	 * asked to traverse.
	 *
	 * In all cases we should transfer control to a "symbolic return" node.
	 * This will make life easier for DFAs.
	 */
	parent := node.GetParent()
	if parent.IsNil() || parent.IsFunction() || (cfa != nil && node == cfa.root) {
		return rhino.NilNode
	}

	// If we are just before a IF/WHILE/DO/FOR:
	switch parent.GetToken() {
	case rhino.Token_IF:
		// The follow() of any of the path from IF would be what follows IF.
		return ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode, parent, cfa)
	case rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
		// After the body of a CASE, the control goes to the body of the next
		// case, without having to go to the case condition.
		if next, ok := ast.MustSwitchCase(parent).Next(); ok {
			return next.Body().Node()
		}
		return ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode, parent, cfa)
	case rhino.Token_FOR:
		// After the body, the update runs (its EMPTY node if the loop has
		// none), then the FOR checks the condition again.
		return cfaForUpdate(ast.MustFor(parent))
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_WHILE, rhino.Token_DO:
		return parent
	case rhino.Token_TRY:
		try := ast.MustTry(parent)
		finally, hasFinally := try.Finally()
		switch {
		// If we are coming out of the TRY block or the CATCH block...
		case try.Body().Is(node), try.Catches().Is(node):
			if hasFinally { // and have FINALLY block.
				return finally.Node()
			}
			// and have no FINALLY.
			return ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode, parent, cfa)
		// If we are coming out of the FINALLY block...
		case hasFinally && finally.Is(node):
			if cfa != nil && len(cfa.finallyMap.entries) > 0 {
				if e := cfa.finallyMap.find(parent); e != nil {
					for i := uint8(0); i < e.nHead; i++ {
						cfa.createEdge(fromNode, ControlFlowGraph_Branch_ON_EX, e.head[i])
					}
					for _, finallyNode := range e.overflow {
						cfa.createEdge(fromNode, ControlFlowGraph_Branch_ON_EX, finallyNode)
					}
				}
			}
			return ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode, parent, cfa)
		}
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected TRY child " + node.String()))
	default:
	}

	// Now that we are done with the special cases follow should be its
	// immediate sibling, unless its sibling is a function
	nextSibling := node.GetNext()

	// Skip function declarations because control doesn't get pass into it.
	for !nextSibling.IsNil() && nextSibling.IsFunction() {
		nextSibling = nextSibling.GetNext()
	}

	if !nextSibling.IsNil() {
		return ControlFlowAnalysis_ComputeFallThrough(nextSibling)
	} else {
		// If there are no more siblings, control is transferred up the AST.
		return ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode, parent, cfa)
	}
}

// ControlFlowAnalysis_ComputeFallThrough computes the destination node of n
// when we want to fallthrough into the subtree of n. We don't always create
// a CFG edge into n itself because of DOs and FORs.
//
// Loops and labels start inside; every other node, a block included, starts
// at itself.
//
// Java: ControlFlowAnalysis.computeFallThrough(Node)
func ControlFlowAnalysis_ComputeFallThrough(n rhino.Node) rhino.Node {
	switch n.GetToken() {
	case rhino.Token_DO:
		// A do-while starts with its body. (Java recurses into it, but a
		// block starts at its own node.)
		return ast.MustDo(n).Body().Node()
	case rhino.Token_FOR:
		// A for loop starts with its init (its EMPTY node if it has none).
		// (Java recurses into it, but an init, a declaration or an
		// expression, starts at its own node.)
		return cfaForInit(ast.MustFor(n))
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// A for-in/of loop starts by evaluating its collection.
		return ast.MustEnhancedFor(n).Collection().Node()
	case rhino.Token_LABEL:
		return ControlFlowAnalysis_ComputeFallThrough(ast.MustLabel(n).Body().Node())
	default:
		return n
	}
}

func cfaHasEdgeInDirection(
	source *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
	branch ControlFlowGraph_Branch,
	dest *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
) bool {
	outEdges := source.OutEdgeList.Slice()
	inEdges := dest.InEdgeList.Slice()
	if len(outEdges) < len(inEdges) {
		for _, outEdge := range outEdges {
			if outEdge.DestNode == dest && outEdge.Value == branch {
				return true
			}
		}
	} else {
		for _, inEdge := range inEdges {
			if inEdge.SourceNode == source && inEdge.Value == branch {
				return true
			}
		}
	}
	return false
}

// createEdge connects the two nodes in the control flow graph.
//
// fromNode is the source; toNode is the destination.
//
// Java: ControlFlowAnalysis.createEdge(Node, ControlFlowGraph.Branch, Node)
func (c *ControlFlowAnalysis) createEdge(fromNode rhino.Node, branch ControlFlowGraph_Branch, toNode rhino.Node) {
	from := c.cfg.CreateLinkedDiGraphNode(fromNode)
	to := c.cfg.CreateLinkedDiGraphNode(toNode)
	if !cfaHasEdgeInDirection(from, branch, to) && !cfaHasEdgeInDirection(to, branch, from) {
		c.cfg.ConnectLinkedNodes(from, branch, to)
	}
}

// connectToPossibleExceptionHandler connects cfgNode to the proper CATCH
// block if target subtree might throw an exception. If there are FINALLY
// blocks reached before a CATCH, it will make the corresponding entry in
// finallyMap.
//
// Java: ControlFlowAnalysis.connectToPossibleExceptionHandler(Node, Node)
func (c *ControlFlowAnalysis) connectToPossibleExceptionHandler(cfgNode rhino.Node, target rhino.Node) {
	if len(c.exceptionHandler) == 0 || c.exceptionHandler[len(c.exceptionHandler)-1].IsFunction() {
		return
	}
	if ControlFlowAnalysis_MayThrowException(target) {
		c.connectToPossibleExceptionHandlerUnconditionally(cfgNode)
	}
}

// Java: ControlFlowAnalysis.connectToPossibleExceptionHandlerUnconditionally(Node)
func (c *ControlFlowAnalysis) connectToPossibleExceptionHandlerUnconditionally(cfgNode rhino.Node) {
	if len(c.exceptionHandler) == 0 {
		return
	}
	lastJump := cfgNode
	for i := len(c.exceptionHandler) - 1; i >= 0; i-- {
		handler := c.exceptionHandler[i]
		if handler.IsFunction() {
			return
		}
		try := ast.MustTry(handler)
		catchBlock := try.Catches()

		lastJumpInCatchBlock := false
		for ancestor := lastJump; !ancestor.IsNil(); ancestor = ancestor.GetParent() {
			if ancestor == handler {
				break
			} else if catchBlock.Is(ancestor) {
				lastJumpInCatchBlock = true
				break
			}
		}

		// No catch but a FINALLY, or lastJump is inside the catch block.
		if _, hasCatch := catchBlock.Catch(); !hasCatch || lastJumpInCatchBlock {
			// The TRY has a FINALLY: a TRY without catch must have one, and a
			// TRY without FINALLY leaves the exceptionHandler stack before its
			// catch block is traversed (shouldTraverseIntoChildren), so
			// lastJump is in the catch block only of a TRY with a FINALLY.
			// (Java takes the last child of the TRY; Go-only: the check
			// states the invariant.)
			finally, ok := try.Finally()
			jutil.CheckState(ok)
			if lastJump == cfgNode {
				c.createEdge(cfgNode, ControlFlowGraph_Branch_ON_EX, finally.Node())
			} else {
				c.putFinally(lastJump, finally.Node())
			}
		} else { // Has a catch.
			if lastJump == cfgNode {
				c.createEdge(cfgNode, ControlFlowGraph_Branch_ON_EX, catchBlock.Node())
				return
			} else {
				c.putFinally(lastJump, catchBlock.Node())
			}
		}
		lastJump = handler
	}
}

func cfaIsBreakTargetByID(target rhino.Node, labelID uint32, hasLabel bool) bool {
	return ControlFlowAnalysis_IsBreakStructure(target, hasLabel) && cfaMatchLabelByID(target.GetParent(), labelID, hasLabel)
}

func cfaIsContinueTargetByID(target rhino.Node, labelID uint32, hasLabel bool) bool {
	return NodeUtil_IsLoopStructure(target) && cfaMatchLabelByID(target.GetParent(), labelID, hasLabel)
}

func cfaMatchLabelByID(target rhino.Node, labelID uint32, hasLabel bool) bool {
	if !hasLabel {
		return true
	}
	for l, ok := ast.AsLabel(target); ok; l, ok = ast.AsLabel(l.Node().GetParent()) {
		if l.Name().Node().GetStringID() == labelID {
			return true
		}
	}
	return false
}

// ControlFlowAnalysis_MayThrowException determines if the subtree might
// throw an exception.
//
// Java: ControlFlowAnalysis.mayThrowException(Node)
func ControlFlowAnalysis_MayThrowException(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_CALL,
		rhino.Token_OPTCHAIN_CALL,
		rhino.Token_TAGGED_TEMPLATELIT,
		rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETPROP,
		rhino.Token_GETELEM,
		rhino.Token_OPTCHAIN_GETELEM,
		rhino.Token_THROW,
		rhino.Token_NEW,
		rhino.Token_ASSIGN,
		rhino.Token_INC,
		rhino.Token_DEC,
		rhino.Token_INSTANCEOF,
		rhino.Token_IN,
		rhino.Token_YIELD,
		rhino.Token_AWAIT:
		return true
	case rhino.Token_FUNCTION:
		return false
	default:
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !ControlFlowGraph_IsEnteringNewCfgNode(c) && ControlFlowAnalysis_MayThrowException(c) {
			return true
		}
	}
	return false
}

// ControlFlowAnalysis_IsBreakStructure determines whether the given node can
// be terminated with a BREAK node.
//
// Java: ControlFlowAnalysis.isBreakStructure(Node, boolean)
func ControlFlowAnalysis_IsBreakStructure(n rhino.Node, labeled bool) bool {
	switch n.GetToken() {
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_DO, rhino.Token_WHILE, rhino.Token_SWITCH:
		return true
	case rhino.Token_BLOCK, rhino.Token_IF, rhino.Token_TRY, rhino.Token_BREAK:
		return labeled
	default:
		return false
	}
}

// ControlFlowAnalysis_GetExceptionHandler gets the TRY block with a CATCH
// that would be run if n throws an exception.
//
// Returns the CATCH node or null if it there isn't a CATCH before the
// function terminates.
//
// Java: ControlFlowAnalysis.getExceptionHandler(Node)
func ControlFlowAnalysis_GetExceptionHandler(n rhino.Node) rhino.Node {
	for cur := n; !cur.IsScript() && !cur.IsFunction(); cur = cur.GetParent() {
		catchNode := ControlFlowAnalysis_GetCatchHandlerForBlock(cur)
		if !catchNode.IsNil() {
			return catchNode
		}
	}
	return rhino.NilNode
}

// ControlFlowAnalysis_GetCatchHandlerForBlock locates the catch BLOCK given
// the first block in a TRY.
//
// Returns the CATCH node or null there is no catch handler.
//
// Java: ControlFlowAnalysis.getCatchHandlerForBlock(Node)
func ControlFlowAnalysis_GetCatchHandlerForBlock(block rhino.Node) rhino.Node {
	// (Java looks for a catch in the blocks after the TRY block: the catch
	// holder, then the FINALLY block, which never starts with a CATCH.)
	if try, ok := ast.AsTry(block.GetParent()); ok && try.Body().Is(block) {
		if catch, ok := try.Catch(); ok {
			return catch.Node()
		}
	}
	return rhino.NilNode
}

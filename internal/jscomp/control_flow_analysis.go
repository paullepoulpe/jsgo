// Port of com/google/javascript/jscomp/ControlFlowAnalysis.java @ closure-compiler v20260415

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
	"cmp"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ControlFlowAnalysis computes a ControlFlowGraph for a given AST
//
// Example usage:
//
//	cfg := ControlFlowAnalysis_Builder_().
//		SetCompiler(compiler).
//		SetCfgRoot(functionRoot).
//		SetIncludeEdgeAnnotations(true).
//		ComputeCfg()
type ControlFlowAnalysis struct {
	// Based roughly on the first few pages of
	//
	// "Declarative Intraprocedural Flow Analysis of Java Source Code by
	// Nilsson-Nyman, Hedin, Magnusson & Ekman",
	//
	// This pass computes the control flow graph from the AST. However, a full
	// attribute grammar is not necessary. We will compute the flow edges with a
	// single post order traversal. The "follow()" of a given node will be
	// computed recursively in a demand driven fashion.
	//
	// At the moment, we are not performing any inter-procedural analysis.
	compiler AbstractCompiler

	cfg ControlFlowGraph[*rhino.Node]

	// Go-only: language-adapter — the Integer values are *int, so that an
	// absent key (Java null) stays distinguishable from position 0.
	astPosition *base.LinkedIdentityHashMap[*rhino.Node, *int]

	// We order CFG nodes by by looking at the AST positions.
	// CFG nodes that come first lexically should be visited first, because
	// they will often be executed first in the source program.
	priorityComparator func(a, b DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) int

	astPositionCounter int
	priorityCounter    int

	shouldTraverseFunctions bool
	edgeAnnotations         bool

	// We need to store where we started, in case we aren't doing a flow
	// analysis for the whole scope. This happens, for example, when running
	// type inference on only the externs.
	root *rhino.Node

	// This stack captures the structure of nested TRY blocks. The top of the
	// stack is the inner most TRY block. A FUNCTION node in this stack
	// implies that the handler is determined by the caller of the function at
	// runtime.
	//
	// Go-only: language-adapter — the ArrayDeque used as a stack is a slice
	// whose last element is the top (push = addFirst, peek = peekFirst,
	// pop = removeFirst; iteration runs from the top down).
	exceptionHandler []*rhino.Node

	// This map is used to handle the follow of FINALLY. For example:
	//
	//	while(x) {
	//	 try {
	//	   try {
	//	     break;
	//	   } catch (a) {
	//	   } finally {
	//	     foo();
	//	   }
	//	   fooFollow();
	//	 } catch (b) {
	//	 } finally {
	//	   bar();
	//	 }
	//	 barFollow();
	//	}
	//	END();
	//
	// In this case finallyMap will contain a map from:
	//
	//	first FINALLY -> bar()
	//	second FINALLY -> END()
	//
	// When we are connecting foo() and bar() to to their respective follow,
	// we must also look up this map and connect:
	//
	//	foo() -> bar()
	//	bar() -> END
	//
	// Go-only: language-adapter — Java HashMap order: a HashMultimap keyed by
	// Nodes, which hash by identity, so Java iterates the values of a key
	// (the ON_EX edges created when leaving a FINALLY block) in an order that
	// differs from run to run. The port iterates them in insertion order
	// (like IdentityHashMap → LinkedHashMap, PORTING.md §6).
	finallyMap *jutil.SetMultimap[*rhino.Node, *rhino.Node]
}

// newControlFlowAnalysis is the constructor. Should only be called from
// within the Builder.
//
// compiler is the compiler instance; shouldTraverseFunctions is whether
// functions should be traversed; edgeAnnotations is whether to allow edge
// annotations.
//
// Java: ControlFlowAnalysis(AbstractCompiler, boolean, boolean)
func newControlFlowAnalysis(compiler AbstractCompiler, shouldTraverseFunctions bool, edgeAnnotations bool) *ControlFlowAnalysis {
	c := &ControlFlowAnalysis{
		exceptionHandler: []*rhino.Node{},
		finallyMap:       jutil.NewSetMultimap[*rhino.Node, *rhino.Node](),
	}
	c.priorityComparator = func(a, b DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) int {
		// Comparator.comparingInt(keyExtractor)
		return cmp.Compare(c.astPositionOf(a), c.astPositionOf(b))
	}
	c.compiler = compiler
	c.shouldTraverseFunctions = shouldTraverseFunctions
	c.edgeAnnotations = edgeAnnotations
	return c
}

// astPositionOf is the key extractor of priorityComparator:
// checkNotNull(astPosition.get(digraphNode.getValue()), digraphNode.getValue()).
func (c *ControlFlowAnalysis) astPositionOf(digraphNode DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) int {
	p := c.astPosition.Get(digraphNode.GetValue())
	if p == nil {
		// checkNotNull(T, Object errorMessage)
		panic(jutil.NewNullPointerExceptionMsg(jutil.ToJavaString(digraphNode.GetValue())))
	}
	return *p
}

// ControlFlowAnalysis_Builder configures a ControlFlowAnalysis instance then
// computes the ControlFlowGraph
type ControlFlowAnalysis_Builder struct {
	compiler                AbstractCompiler
	cfgRoot                 *rhino.Node
	shouldTraverseFunctions bool
	edgeAnnotations         bool
}

// Java: ControlFlowAnalysis.Builder()
func newControlFlowAnalysis_Builder() *ControlFlowAnalysis_Builder {
	return &ControlFlowAnalysis_Builder{
		shouldTraverseFunctions: false,
		edgeAnnotations:         false,
	}
}

// Java: ControlFlowAnalysis.Builder.setCompiler(AbstractCompiler)
func (b *ControlFlowAnalysis_Builder) SetCompiler(compiler AbstractCompiler) *ControlFlowAnalysis_Builder {
	b.compiler = compiler
	return b
}

// Java: ControlFlowAnalysis.Builder.setCfgRoot(Node)
func (b *ControlFlowAnalysis_Builder) SetCfgRoot(cfgRoot *rhino.Node) *ControlFlowAnalysis_Builder {
	b.cfgRoot = cfgRoot
	return b
}

// Java: ControlFlowAnalysis.Builder.setTraverseFunctions(boolean)
func (b *ControlFlowAnalysis_Builder) SetTraverseFunctions(shouldTraverseFunctions bool) *ControlFlowAnalysis_Builder {
	b.shouldTraverseFunctions = shouldTraverseFunctions
	return b
}

// Java: ControlFlowAnalysis.Builder.setIncludeEdgeAnnotations(boolean)
func (b *ControlFlowAnalysis_Builder) SetIncludeEdgeAnnotations(includeEdgeAnnotations bool) *ControlFlowAnalysis_Builder {
	b.edgeAnnotations = includeEdgeAnnotations
	return b
}

// Java: ControlFlowAnalysis.Builder.computeCfg()
func (b *ControlFlowAnalysis_Builder) ComputeCfg() ControlFlowGraph[*rhino.Node] {
	jutil.CheckNotNullMsg(b.compiler, "Need to call setCompiler()")
	jutil.CheckNotNullMsg(b.cfgRoot, "Need to call setCfgRoot()")
	cfa := newControlFlowAnalysis(b.compiler, b.shouldTraverseFunctions, b.edgeAnnotations)
	cfa.computeCfg(b.cfgRoot)
	return cfa.cfg
}

// Java: ControlFlowAnalysis.builder()
func ControlFlowAnalysis_Builder_() *ControlFlowAnalysis_Builder {
	return newControlFlowAnalysis_Builder()
}

// Java: ControlFlowAnalysis.getCfg()
func (c *ControlFlowAnalysis) GetCfg() ControlFlowGraph[*rhino.Node] {
	return c.cfg
}

// Java: ControlFlowAnalysis.computeCfg(Node)
func (c *ControlFlowAnalysis) computeCfg(root *rhino.Node) {
	jutil.CheckArgumentMsg1(
		NodeUtil_IsValidCfgRoot(root), "Unexpected control flow graph root %s", root)
	c.root = root
	c.astPosition = base.NewLinkedIdentityHashMap[*rhino.Node, *int]()
	c.astPositionCounter = 0
	c.cfg = newControlFlowAnalysis_AstControlFlowGraph(ControlFlowAnalysis_ComputeFallThrough(root), c.edgeAnnotations)

	// Traverse the graph.
	NodeTraversal_Traverse(c.compiler, root, c)

	// Insert an implicit return last.
	c.astPositionCounter++
	c.astPosition.Put(nil, newControlFlowAnalysis_Position(c.astPositionCounter))

	// Now, generate the priority of nodes by doing a depth-first
	// search on the CFG.
	c.priorityCounter = 0
	entry := c.cfg.GetEntry()
	c.prioritizeFromEntryNode(entry)

	if c.shouldTraverseFunctions {
		// If we're traversing inner functions, we need to rank the
		// priority of them too.
		for _, candidate := range c.cfg.GetDiGraphNodes() {
			value := candidate.GetValue()
			if value != nil && value.IsFunction() {
				c.prioritizeFromEntryNode(candidate)
			}
		}
	}

	// Free our large table now that it is no longer needed
	c.astPosition = nil

	// At this point, all reachable nodes have been given a priority, but
	// unreachable nodes have not been given a priority. Put them last.
	// Presumably, it doesn't really matter what priority they get, since
	// this shouldn't happen in real code.
	for _, candidate := range c.cfg.GetDiGraphNodes() {
		if !candidate.HasPriority() {
			c.priorityCounter++
			candidate.SetPriority(c.priorityCounter)
		}
	}

	// Again, the implicit return node is always last.
	c.priorityCounter++
	c.cfg.GetImplicitReturn().SetPriority(c.priorityCounter)
}

// newControlFlowAnalysis_Position boxes an AST position (Integer.valueOf).
func newControlFlowAnalysis_Position(p int) *int {
	return &p
}

// prioritizeFromEntryNode: given an entry node, find all the nodes reachable
// from that node and prioritize them.
//
// Java: ControlFlowAnalysis.prioritizeFromEntryNode(DiGraphNode<Node, Branch>)
func (c *ControlFlowAnalysis) prioritizeFromEntryNode(entry DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) {
	worklist := jutil.NewPriorityQueueWithCapacity(10, c.priorityComparator)
	worklist.Add(entry)

	for !worklist.IsEmpty() {
		current := worklist.Remove()
		if current.HasPriority() {
			continue
		}

		c.priorityCounter++
		current.SetPriority(c.priorityCounter)

		successors := c.cfg.GetDirectedSuccNodes(current)
		worklist.AddAll(successors)
	}
}

// Java: ControlFlowAnalysis.shouldTraverse(NodeTraversal, Node, Node)
func (c *ControlFlowAnalysis) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if c.shouldTraverseIntoChildren(n, parent) {
		// Any AST node that will later have a corresponding CFG node must be in astPosition.
		// To avoid having astPosition grow too large, we exclude AST nodes that are not traversed
		// further, as they usually are not put in the CFG.
		c.astPosition.Put(n, newControlFlowAnalysis_Position(c.astPositionCounter))
		c.astPositionCounter++
		return true
	}
	return false
}

// exceptionHandlerPush is ArrayDeque.push(e).
func (c *ControlFlowAnalysis) exceptionHandlerPush(n *rhino.Node) {
	c.exceptionHandler = append(c.exceptionHandler, n)
}

// exceptionHandlerPeek is ArrayDeque.peek(): the top, or null when empty.
func (c *ControlFlowAnalysis) exceptionHandlerPeek() *rhino.Node {
	if len(c.exceptionHandler) == 0 {
		return nil
	}
	return c.exceptionHandler[len(c.exceptionHandler)-1]
}

// exceptionHandlerPop is ArrayDeque.pop(): removes the top, or throws
// NoSuchElementException when empty.
func (c *ControlFlowAnalysis) exceptionHandlerPop() *rhino.Node {
	if len(c.exceptionHandler) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	top := c.exceptionHandler[len(c.exceptionHandler)-1]
	c.exceptionHandler[len(c.exceptionHandler)-1] = nil
	c.exceptionHandler = c.exceptionHandler[:len(c.exceptionHandler)-1]
	return top
}

// shouldTraverseIntoChildren returns whether the children of this node
// should be traversed as part of this control-flow analysis, i.e. whether
// the control-flow graph being built may require any edges into children of
// this node.
//
// Java: ControlFlowAnalysis.shouldTraverseIntoChildren(Node, Node)
func (c *ControlFlowAnalysis) shouldTraverseIntoChildren(n *rhino.Node, parent *rhino.Node) bool {
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
	if parent != nil {
		switch parent.GetToken() {
		case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			// Only traverse the body of the for loop.
			shouldTraverseForChild := n == parent.GetLastChild()
			if !shouldTraverseForChild {
				// The control-flow graph contains edges from a FOR node to all its children, even
				// though only the body is actually traversed. So put the other children in astPosition.
				c.astPosition.Put(n, newControlFlowAnalysis_Position(c.astPositionCounter))
				c.astPositionCounter++
			}
			return shouldTraverseForChild
		case rhino.Token_DO:
			// Only traverse the body of the do-while.
			return n != parent.GetSecondChild()
		case rhino.Token_IF, rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_CATCH, rhino.Token_LABEL:
			// Skip conditions, and only traverse the body of the cases
			return n != parent.GetFirstChild()
		case rhino.Token_FUNCTION:
			return n == parent.GetLastChild()
		case rhino.Token_CLASS:
			return c.shouldTraverseFunctions && n == parent.GetLastChild()
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
			if (!NodeUtil_HasFinally(parent) && n == NodeUtil_GetCatchBlock(parent)) ||
				NodeUtil_IsTryFinallyNode(parent, n) {
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

// Java: ControlFlowAnalysis.visit(NodeTraversal, Node, Node)
func (c *ControlFlowAnalysis) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_IF:
		c.handleIf(n)
	case rhino.Token_WHILE:
		c.handleWhile(n)
	case rhino.Token_DO:
		c.handleDo(n)
	case rhino.Token_FOR:
		c.handleFor(n)
	case rhino.Token_FOR_OF, rhino.Token_FOR_IN, rhino.Token_FOR_AWAIT_OF:
		c.handleEnhancedFor(n)
	case rhino.Token_SWITCH:
		c.handleSwitch(n)
	case rhino.Token_CASE:
		c.handleCase(n)
	case rhino.Token_DEFAULT_CASE:
		c.handleDefault(n)
	case rhino.Token_BLOCK, rhino.Token_ROOT, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY:
		c.handleStmtList(n)
	case rhino.Token_FUNCTION:
		c.handleFunction(n)
	case rhino.Token_EXPR_RESULT:
		c.handleExpr(n)
	case rhino.Token_THROW:
		c.handleThrow(n)
	case rhino.Token_TRY:
		c.handleTry(n)
	case rhino.Token_CATCH:
		c.handleCatch(n)
	case rhino.Token_BREAK:
		c.handleBreak(n)
	case rhino.Token_CONTINUE:
		c.handleContinue(n)
	case rhino.Token_RETURN:
		c.handleReturn(n)
	case rhino.Token_WITH:
		c.handleWith(n)
	case rhino.Token_SWITCH_BODY,
		rhino.Token_LABEL,
		rhino.Token_CLASS_MEMBERS,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_COMPUTED_FIELD_DEF:
	default:
		c.handleStmt(n)
	}
}

// Java: ControlFlowAnalysis.handleIf(Node)
func (c *ControlFlowAnalysis) handleIf(node *rhino.Node) {
	thenBlock := node.GetSecondChild()
	elseBlock := thenBlock.GetNext()
	c.createEdge(node, ControlFlowGraph_Branch_ON_TRUE, ControlFlowAnalysis_ComputeFallThrough(thenBlock))

	if elseBlock == nil {
		c.createEdge(node, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c)) // not taken branch
	} else {
		c.createEdge(node, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFallThrough(elseBlock))
	}
	c.connectToPossibleExceptionHandler(node, NodeUtil_GetConditionExpression(node))
}

// Java: ControlFlowAnalysis.handleWhile(Node)
func (c *ControlFlowAnalysis) handleWhile(node *rhino.Node) {
	cond := node.GetFirstChild()
	// Control goes to the first statement if the condition evaluates to true.
	c.createEdge(node, ControlFlowGraph_Branch_ON_TRUE, ControlFlowAnalysis_ComputeFallThrough(cond.GetNext()))

	if !cond.IsTrue() {
		// Control goes to the follow() if the condition evaluates to false.
		c.createEdge(node, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c))
	}
	c.connectToPossibleExceptionHandler(node, NodeUtil_GetConditionExpression(node))
}

// Java: ControlFlowAnalysis.handleDo(Node)
func (c *ControlFlowAnalysis) handleDo(node *rhino.Node) {
	body := node.GetFirstChild()
	// The first edge can be the initial iteration as well as the iterations
	// after.
	c.createEdge(node, ControlFlowGraph_Branch_ON_TRUE, ControlFlowAnalysis_ComputeFallThrough(body))
	cond := body.GetNext()
	if !cond.IsTrue() {
		// The edge that leaves the do loop if the condition fails.
		c.createEdge(node, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c))
	}
	c.connectToPossibleExceptionHandler(node, NodeUtil_GetConditionExpression(node))
}

// Java: ControlFlowAnalysis.handleEnhancedFor(Node)
func (c *ControlFlowAnalysis) handleEnhancedFor(forNode *rhino.Node) {
	// We have:  for (index in collection) { body }
	// or:       for (item of collection) { body }
	// or:       for await (item of collection) { body }
	item := forNode.GetFirstChild()
	collection := item.GetNext()
	body := collection.GetNext()
	// The collection behaves like init.
	c.createEdge(collection, ControlFlowGraph_Branch_UNCOND, forNode)
	// The edge that transfer control to the beginning of the loop body.
	c.createEdge(forNode, ControlFlowGraph_Branch_ON_TRUE, ControlFlowAnalysis_ComputeFallThrough(body))
	// The edge to end of the loop.
	c.createEdge(forNode, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(forNode, c))
	c.connectToPossibleExceptionHandler(forNode, collection)
}

// Java: ControlFlowAnalysis.handleFor(Node)
func (c *ControlFlowAnalysis) handleFor(forNode *rhino.Node) {
	// We have for (init; cond; iter) { body }
	init := forNode.GetFirstChild()
	cond := init.GetNext()
	iter := cond.GetNext()
	body := iter.GetNext()
	// After initialization, we transfer to the FOR which is in charge of
	// checking the condition (for the first time).
	c.createEdge(init, ControlFlowGraph_Branch_UNCOND, forNode)
	// The edge that transfer control to the beginning of the loop body.
	c.createEdge(forNode, ControlFlowGraph_Branch_ON_TRUE, ControlFlowAnalysis_ComputeFallThrough(body))
	// The edge to end of the loop.
	if !cond.IsEmpty() && !cond.IsTrue() {
		c.createEdge(forNode, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(forNode, c))
	}
	// The end of the body will have a unconditional branch to our iter
	// (handled by calling computeFollowNode of the last instruction of the
	// body. Our iter will jump to the forNode again to another condition
	// check.
	c.createEdge(iter, ControlFlowGraph_Branch_UNCOND, forNode)
	c.connectToPossibleExceptionHandler(init, init)
	c.connectToPossibleExceptionHandler(forNode, cond)
	c.connectToPossibleExceptionHandler(iter, iter)
}

// Java: ControlFlowAnalysis.handleSwitch(Node)
func (c *ControlFlowAnalysis) handleSwitch(node *rhino.Node) {
	// Transfer to the first non-DEFAULT CASE. if there are none, transfer
	// to the DEFAULT or the EMPTY node.
	switchBody := node.GetLastChild()
	c.createEdge(node, ControlFlowGraph_Branch_UNCOND, switchBody)
	next := ControlFlowAnalysis_GetNextSiblingOfType(switchBody.GetFirstChild(), rhino.Token_CASE, rhino.Token_EMPTY)
	if next != nil { // Has at least one CASE or EMPTY
		c.createEdge(switchBody, ControlFlowGraph_Branch_UNCOND, next)
	} else { // Has no CASE but possibly a DEFAULT
		if switchBody.HasChildren() {
			c.createEdge(switchBody, ControlFlowGraph_Branch_UNCOND, switchBody.GetFirstChild())
		} else { // No CASE, no DEFAULT
			c.createEdge(switchBody, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c))
		}
	}
	c.connectToPossibleExceptionHandler(node, node.GetFirstChild())
}

// Java: ControlFlowAnalysis.handleCase(Node)
func (c *ControlFlowAnalysis) handleCase(caseNode *rhino.Node) {
	// Case is a bit tricky....First it goes into the body if condition is true.
	c.createEdge(caseNode, ControlFlowGraph_Branch_ON_TRUE, caseNode.GetSecondChild())

	// Look for the next CASE, skipping over DEFAULT.
	next := ControlFlowAnalysis_GetNextSiblingOfType(caseNode.GetNext(), rhino.Token_CASE)
	if next != nil { // Found a CASE
		jutil.CheckState(next.IsCase())
		c.createEdge(caseNode, ControlFlowGraph_Branch_ON_FALSE, next)
	} else { // No more CASE found, go back and search for a DEFAULT.
		switchBody := caseNode.GetParent()
		deflt := ControlFlowAnalysis_GetNextSiblingOfType(switchBody.GetFirstChild(), rhino.Token_DEFAULT_CASE)
		if deflt != nil { // Has a DEFAULT
			c.createEdge(caseNode, ControlFlowGraph_Branch_ON_FALSE, deflt)
		} else { // No DEFAULT found, go to the follow of the SWITCH.
			c.createEdge(caseNode, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(caseNode, c))
		}
	}
	c.connectToPossibleExceptionHandler(caseNode, caseNode.GetFirstChild())
}

// Java: ControlFlowAnalysis.handleDefault(Node)
func (c *ControlFlowAnalysis) handleDefault(node *rhino.Node) {
	// Directly goes to the body. It should not transfer to the next case.
	c.createEdge(node, ControlFlowGraph_Branch_UNCOND, node.GetFirstChild())
}

// Java: ControlFlowAnalysis.handleWith(Node)
func (c *ControlFlowAnalysis) handleWith(node *rhino.Node) {
	// Directly goes to the body. It should not transfer to the next case.
	c.createEdge(node, ControlFlowGraph_Branch_UNCOND, node.GetLastChild())
	c.connectToPossibleExceptionHandler(node, node.GetFirstChild())
}

// Java: ControlFlowAnalysis.handleStmtList(Node)
func (c *ControlFlowAnalysis) handleStmtList(node *rhino.Node) {
	parent := node.GetParent()
	// Special case, don't add a block of empty CATCH block to the graph.
	if node.IsBlock() &&
		parent.IsTry() &&
		NodeUtil_GetCatchBlock(parent) == node &&
		!NodeUtil_HasCatchHandler(node) {
		return
	}

	// A block transfer control to its first child if it is not empty.
	child := node.GetFirstChild()

	// Function declarations are skipped since control doesn't go into that
	// function (unless it is called)
	for child != nil && child.IsFunction() {
		child = child.GetNext()
	}

	if child != nil {
		c.createEdge(node, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFallThrough(child))
	} else {
		c.createEdge(node, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c))
	}

	// Synthetic blocks
	if parent != nil {
		switch parent.GetToken() {
		case rhino.Token_DEFAULT_CASE, rhino.Token_CASE, rhino.Token_TRY:
		case rhino.Token_ROOT:
			// TODO(b/71873602): why is this path necessary?
			if node.IsRoot() && node.GetNext() != nil {
				c.createEdge(node, ControlFlowGraph_Branch_UNCOND, node.GetNext())
			}
		default:
			if node.IsBlock() && node.IsSyntheticBlock() {
				c.createEdge(node, ControlFlowGraph_Branch_SYN_BLOCK, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c))
			}
		}
	}
}

// Java: ControlFlowAnalysis.handleFunction(Node)
func (c *ControlFlowAnalysis) handleFunction(node *rhino.Node) {
	// A block transfer control to its first child if it is not empty.
	jutil.CheckState(node.IsFunction())
	jutil.CheckState(node.HasXChildren(3))
	c.createEdge(node, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFallThrough(node.GetLastChild()))
	jutil.CheckState(c.exceptionHandlerPeek() == node)
	c.exceptionHandlerPop()
}

// Java: ControlFlowAnalysis.handleExpr(Node)
func (c *ControlFlowAnalysis) handleExpr(node *rhino.Node) {
	c.createEdge(node, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c))
	c.connectToPossibleExceptionHandler(node, node)
}

// Java: ControlFlowAnalysis.handleThrow(Node)
func (c *ControlFlowAnalysis) handleThrow(node *rhino.Node) {
	c.connectToPossibleExceptionHandler(node, node)
}

// Java: ControlFlowAnalysis.handleTry(Node)
func (c *ControlFlowAnalysis) handleTry(node *rhino.Node) {
	c.createEdge(node, ControlFlowGraph_Branch_UNCOND, node.GetFirstChild())
}

// Java: ControlFlowAnalysis.handleCatch(Node)
func (c *ControlFlowAnalysis) handleCatch(node *rhino.Node) {
	c.createEdge(node, ControlFlowGraph_Branch_UNCOND, node.GetLastChild())
}

// Java: ControlFlowAnalysis.handleBreak(Node)
func (c *ControlFlowAnalysis) handleBreak(node *rhino.Node) {
	var label *string
	// See if it is a break with label.
	if node.HasChildren() {
		s := node.GetFirstChild().GetString()
		label = &s
	}
	var cur *rhino.Node
	var previous *rhino.Node
	var lastJump *rhino.Node
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
	for cur, lastJump = node, node; !ControlFlowAnalysis_IsBreakTarget(cur, label); cur, parent = parent, parent.GetParent() {
		if cur.IsTry() && NodeUtil_HasFinally(cur) && cur.GetLastChild() != previous {
			if lastJump == node {
				c.createEdge(lastJump, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFallThrough(cur.GetLastChild()))
			} else {
				c.finallyMap.Put(lastJump, ControlFlowAnalysis_ComputeFallThrough(cur.GetLastChild()))
			}
			lastJump = cur
		}
		if parent == nil {
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
		c.finallyMap.Put(lastJump, ControlFlowAnalysis_ComputeFollowNodeWithCfa(cur, c))
	}
}

// Java: ControlFlowAnalysis.handleContinue(Node)
func (c *ControlFlowAnalysis) handleContinue(node *rhino.Node) {
	var label *string
	if node.HasChildren() {
		s := node.GetFirstChild().GetString()
		label = &s
	}
	var cur *rhino.Node
	var previous *rhino.Node
	var lastJump *rhino.Node

	// Similar to handBreak's logic with a few minor variation.
	for cur, lastJump = node, node; !ControlFlowAnalysis_IsContinueTarget(cur, label); cur = cur.GetParent() {
		if cur.IsTry() && NodeUtil_HasFinally(cur) && cur.GetLastChild() != previous {
			if lastJump == node {
				c.createEdge(lastJump, ControlFlowGraph_Branch_UNCOND, cur.GetLastChild())
			} else {
				c.finallyMap.Put(lastJump, ControlFlowAnalysis_ComputeFallThrough(cur.GetLastChild()))
			}
			lastJump = cur
		}
		jutil.CheckStateMsg(cur.HasParent(), "Cannot find continue target.")
		previous = cur
	}
	iter := cur
	if cur.IsVanillaFor() {
		// the increment expression happens after the continue
		iter = cur.GetChildAtIndex(2)
	}

	if lastJump == node {
		c.createEdge(node, ControlFlowGraph_Branch_UNCOND, iter)
	} else {
		c.finallyMap.Put(lastJump, iter)
	}
}

// Java: ControlFlowAnalysis.handleReturn(Node)
func (c *ControlFlowAnalysis) handleReturn(node *rhino.Node) {
	var lastJump *rhino.Node
	for i := len(c.exceptionHandler) - 1; i >= 0; i-- {
		curHandler := c.exceptionHandler[i]
		if curHandler.IsFunction() {
			break
		}
		if NodeUtil_HasFinally(curHandler) {
			if lastJump == nil {
				c.createEdge(node, ControlFlowGraph_Branch_UNCOND, curHandler.GetLastChild())
			} else {
				c.finallyMap.Put(lastJump, ControlFlowAnalysis_ComputeFallThrough(curHandler.GetLastChild()))
			}
			lastJump = curHandler
		}
	}

	if node.HasChildren() {
		c.connectToPossibleExceptionHandler(node, node.GetFirstChild())
	}

	if lastJump == nil {
		c.createEdge(node, ControlFlowGraph_Branch_UNCOND, nil)
	} else {
		c.finallyMap.Put(lastJump, nil)
	}
}

// Java: ControlFlowAnalysis.handleStmt(Node)
func (c *ControlFlowAnalysis) handleStmt(node *rhino.Node) {
	// Simply transfer to the next line.
	c.createEdge(node, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c))
	c.connectToPossibleExceptionHandler(node, node)
}

// Java: ControlFlowAnalysis.computeFollowNode(Node, ControlFlowAnalysis)
func ControlFlowAnalysis_ComputeFollowNodeWithCfa(node *rhino.Node, cfa *ControlFlowAnalysis) *rhino.Node {
	return ControlFlowAnalysis_ComputeFollowNodeFrom(node, node, cfa)
}

// Java: ControlFlowAnalysis.computeFollowNode(Node)
func ControlFlowAnalysis_ComputeFollowNode(node *rhino.Node) *rhino.Node {
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
func ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode *rhino.Node, node *rhino.Node, cfa *ControlFlowAnalysis) *rhino.Node {
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
	if parent == nil || parent.IsFunction() || (cfa != nil && node == cfa.root) {
		return nil
	}

	// If we are just before a IF/WHILE/DO/FOR:
	switch parent.GetToken() {
	case rhino.Token_IF:
		// The follow() of any of the path from IF would be what follows IF.
		return ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode, parent, cfa)
	case rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
		// After the body of a CASE, the control goes to the body of the next
		// case, without having to go to the case condition.
		if parent.GetNext() != nil {
			if parent.GetNext().IsCase() {
				return parent.GetNext().GetSecondChild()
			} else if parent.GetNext().IsDefaultCase() {
				return parent.GetNext().GetFirstChild()
			} else {
				panic(jutil.NewIllegalStateExceptionMsg("Not reachable"))
			}
		} else {
			return ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode, parent, cfa)
		}
	case rhino.Token_FOR:
		return parent.GetSecondChild().GetNext()
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_WHILE, rhino.Token_DO:
		return parent
	case rhino.Token_TRY:
		// If we are coming out of the TRY block...
		if parent.GetFirstChild() == node {
			if NodeUtil_HasFinally(parent) { // and have FINALLY block.
				return ControlFlowAnalysis_ComputeFallThrough(parent.GetLastChild())
			} else { // and have no FINALLY.
				return ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode, parent, cfa)
			}
			// CATCH block.
		} else if NodeUtil_GetCatchBlock(parent) == node {
			if NodeUtil_HasFinally(parent) { // and have FINALLY block.
				return ControlFlowAnalysis_ComputeFallThrough(node.GetNext())
			} else {
				return ControlFlowAnalysis_ComputeFollowNodeFrom(fromNode, parent, cfa)
			}
			// If we are coming out of the FINALLY block...
		} else if parent.GetLastChild() == node {
			if cfa != nil {
				for finallyNode := range cfa.finallyMap.Get(parent).All() {
					cfa.createEdge(fromNode, ControlFlowGraph_Branch_ON_EX, finallyNode)
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
	for nextSibling != nil && nextSibling.IsFunction() {
		nextSibling = nextSibling.GetNext()
	}

	if nextSibling != nil {
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
// Java: ControlFlowAnalysis.computeFallThrough(Node)
func ControlFlowAnalysis_ComputeFallThrough(n *rhino.Node) *rhino.Node {
	switch n.GetToken() {
	case rhino.Token_DO, rhino.Token_FOR:
		return ControlFlowAnalysis_ComputeFallThrough(n.GetFirstChild())
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		return n.GetSecondChild()
	case rhino.Token_LABEL:
		return ControlFlowAnalysis_ComputeFallThrough(n.GetLastChild())
	default:
		return n
	}
}

// createEdge connects the two nodes in the control flow graph.
//
// fromNode is the source; toNode is the destination.
//
// Java: ControlFlowAnalysis.createEdge(Node, ControlFlowGraph.Branch, Node)
func (c *ControlFlowAnalysis) createEdge(fromNode *rhino.Node, branch ControlFlowGraph_Branch, toNode *rhino.Node) {
	c.cfg.CreateNode(fromNode)
	c.cfg.CreateNode(toNode)
	c.cfg.ConnectIfNotFound(fromNode, branch, toNode)
}

// connectToPossibleExceptionHandler connects cfgNode to the proper CATCH
// block if target subtree might throw an exception. If there are FINALLY
// blocks reached before a CATCH, it will make the corresponding entry in
// finallyMap.
//
// Java: ControlFlowAnalysis.connectToPossibleExceptionHandler(Node, Node)
func (c *ControlFlowAnalysis) connectToPossibleExceptionHandler(cfgNode *rhino.Node, target *rhino.Node) {
	if ControlFlowAnalysis_MayThrowException(target) && len(c.exceptionHandler) != 0 {
		lastJump := cfgNode
		for i := len(c.exceptionHandler) - 1; i >= 0; i-- {
			handler := c.exceptionHandler[i]
			if handler.IsFunction() {
				return
			}
			jutil.CheckState(handler.IsTry())
			catchBlock := NodeUtil_GetCatchBlock(handler)

			lastJumpInCatchBlock := false
			for ancestor := range lastJump.GetAncestors().All() {
				if ancestor == handler {
					break
				} else if ancestor == catchBlock {
					lastJumpInCatchBlock = true
					break
				}
			}

			// No catch but a FINALLY, or lastJump is inside the catch block.
			if !NodeUtil_HasCatchHandler(catchBlock) || lastJumpInCatchBlock {
				if lastJump == cfgNode {
					c.createEdge(cfgNode, ControlFlowGraph_Branch_ON_EX, handler.GetLastChild())
				} else {
					c.finallyMap.Put(lastJump, handler.GetLastChild())
				}
			} else { // Has a catch.
				if lastJump == cfgNode {
					c.createEdge(cfgNode, ControlFlowGraph_Branch_ON_EX, catchBlock)
					return
				} else {
					c.finallyMap.Put(lastJump, catchBlock)
				}
			}
			lastJump = handler
		}
	}
}

// ControlFlowAnalysis_GetNextSiblingOfType gets the next sibling (including
// itself) of one of the given types.
//
// Java: ControlFlowAnalysis.getNextSiblingOfType(Node, Token...)
func ControlFlowAnalysis_GetNextSiblingOfType(first *rhino.Node, types ...rhino.Token) *rhino.Node {
	for c := first; c != nil; c = c.GetNext() {
		for _, type_ := range types {
			if c.GetToken() == type_ {
				return c
			}
		}
	}
	return nil
}

// ControlFlowAnalysis_IsBreakTarget checks if target is actually the break
// target of labeled continue. The label can be null if it is an unlabeled
// break.
//
// Java: ControlFlowAnalysis.isBreakTarget(Node, String)
func ControlFlowAnalysis_IsBreakTarget(target *rhino.Node, label *string) bool {
	return ControlFlowAnalysis_IsBreakStructure(target, label != nil) && ControlFlowAnalysis_MatchLabel(target.GetParent(), label)
}

// ControlFlowAnalysis_IsContinueTarget checks if target is actually the
// continue target of labeled continue. The label can be null if it is an
// unlabeled continue.
//
// Java: ControlFlowAnalysis.isContinueTarget(Node, String)
func ControlFlowAnalysis_IsContinueTarget(target *rhino.Node, label *string) bool {
	return NodeUtil_IsLoopStructure(target) && ControlFlowAnalysis_MatchLabel(target.GetParent(), label)
}

// ControlFlowAnalysis_MatchLabel checks if label is actually referencing the
// target control structure. If label is null, it always returns true.
//
// Java: ControlFlowAnalysis.matchLabel(Node, String)
func ControlFlowAnalysis_MatchLabel(target *rhino.Node, label *string) bool {
	if label == nil {
		return true
	}
	for target.IsLabel() {
		if target.GetFirstChild().GetString() == *label {
			return true
		}
		target = target.GetParent()
	}
	return false
}

// ControlFlowAnalysis_MayThrowException determines if the subtree might
// throw an exception.
//
// Java: ControlFlowAnalysis.mayThrowException(Node)
func ControlFlowAnalysis_MayThrowException(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_CALL,
		rhino.Token_TAGGED_TEMPLATELIT,
		rhino.Token_GETPROP,
		rhino.Token_GETELEM,
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
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
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
func ControlFlowAnalysis_IsBreakStructure(n *rhino.Node, labeled bool) bool {
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
func ControlFlowAnalysis_GetExceptionHandler(n *rhino.Node) *rhino.Node {
	for cur := n; !cur.IsScript() && !cur.IsFunction(); cur = cur.GetParent() {
		catchNode := ControlFlowAnalysis_GetCatchHandlerForBlock(cur)
		if catchNode != nil {
			return catchNode
		}
	}
	return nil
}

// ControlFlowAnalysis_GetCatchHandlerForBlock locates the catch BLOCK given
// the first block in a TRY.
//
// Returns the CATCH node or null there is no catch handler.
//
// Java: ControlFlowAnalysis.getCatchHandlerForBlock(Node)
func ControlFlowAnalysis_GetCatchHandlerForBlock(block *rhino.Node) *rhino.Node {
	if block.IsBlock() &&
		block.GetParent().IsTry() &&
		block.GetParent().GetFirstChild() == block {
		for s := block.GetNext(); s != nil; s = s.GetNext() {
			if NodeUtil_HasCatchHandler(s) {
				return s.GetFirstChild()
			}
		}
	}
	return nil
}

// ControlFlowAnalysis_AstControlFlowGraph is a ControlFlowGraph which
// provides a node comparator based on the pre-order traversal of the AST.
type ControlFlowAnalysis_AstControlFlowGraph struct {
	ControlFlowGraphBase[*rhino.Node]
}

// newControlFlowAnalysis_AstControlFlowGraph constructs this graph.
//
// entry is the entry node; edgeAnnotations is the annotation.
//
// Java: ControlFlowAnalysis.AstControlFlowGraph(Node, boolean)
func newControlFlowAnalysis_AstControlFlowGraph(entry *rhino.Node, edgeAnnotations bool) *ControlFlowAnalysis_AstControlFlowGraph {
	g := &ControlFlowAnalysis_AstControlFlowGraph{}
	g.InitControlFlowGraphBase(g, entry /* nodeAnnotations= */, true, edgeAnnotations)
	return g
}

// GetOptionalNodeComparator returns a node comparator based on the
// pre-order traversal of the AST.
//
// isForward: x 'before' y in the pre-order traversal implies x 'less than'
// y (if true) and x 'greater than' y (if false).
//
// Java: ControlFlowAnalysis.AstControlFlowGraph.getOptionalNodeComparator(boolean)
func (g *ControlFlowAnalysis_AstControlFlowGraph) GetOptionalNodeComparator(isForward bool) func(a, b DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) int {
	if isForward {
		return func(a, b DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) int {
			return cmp.Compare(a.GetPriority(), b.GetPriority())
		}
	}
	// Comparator.reversed()
	return func(a, b DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) int {
		return cmp.Compare(b.GetPriority(), a.GetPriority())
	}
}

var (
	_ NodeTraversal_Callback        = (*ControlFlowAnalysis)(nil)
	_ ControlFlowGraph[*rhino.Node] = (*ControlFlowAnalysis_AstControlFlowGraph)(nil)
)

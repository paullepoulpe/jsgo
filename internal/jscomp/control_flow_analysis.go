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
	"cmp"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass reads statement shapes through the typed view of
// package ast (ifs.Else(), try.Finally(), sw.Body().FirstCase()...) instead
// of child positions. The graph it builds stays a
// ControlFlowGraph[rhino.Node], and the helpers that other passes call keep
// their rhino.Node signatures.

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

	cfg *ControlFlowAnalysis_AstControlFlowGraph

	// Go-only: language-adapter — the LinkedIdentityHashMap is a Go map: it is
	// only looked up, never iterated.
	astPosition map[rhino.Node]int32

	astPositionCounter int32
	priorityCounter    int

	shouldTraverseFunctions bool
	edgeAnnotations         bool

	// We need to store where we started, in case we aren't doing a flow
	// analysis for the whole scope. This happens, for example, when running
	// type inference on only the externs.
	root rhino.Node

	// This stack captures the structure of nested TRY blocks. The top of the
	// stack is the inner most TRY block. A FUNCTION node in this stack
	// implies that the handler is determined by the caller of the function at
	// runtime.
	//
	// Go-only: language-adapter — the ArrayDeque used as a stack is a slice
	// whose last element is the top (push = addFirst, peek = peekFirst,
	// pop = removeFirst; iteration runs from the top down).
	exceptionHandler []rhino.Node

	// The PriorityQueue of prioritizeFromEntryNode, kept for reuse.
	worklist []cfaWorkItem

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
	// A LinkedHashMultimap: the values of a key iterate in insertion order.
	finallyMap *jutil.SetMultimap[rhino.Node, rhino.Node]
}

// cfaEmptyFinallyMap stands for an empty finallyMap until putFinally needs a
// real one (most CFGs have no FINALLY).
var cfaEmptyFinallyMap = jutil.NewSetMultimap[rhino.Node, rhino.Node]()

// cfaReuseMaxSize is the largest astPosition map or worklist that a
// ControlFlowAnalysis keeps for its next CFG.
const cfaReuseMaxSize = 1024

// ReusableCfg computes control flow graphs into storage that it reuses from
// one graph to the next. A graph that it returns is valid until the next
// Compute or release.
//
// Go-only: performance — passes compute a CFG for each function that they
// analyze (tens of thousands per compilation); reusing the analysis and
// graph storage of the previous function avoids allocating them again.
type ReusableCfg struct {
	cfa ControlFlowAnalysis
	cfg ControlFlowAnalysis_AstControlFlowGraph
}

// Compute is ControlFlowAnalysis.builder() with these settings followed by
// computeCfg(), into rc's storage.
func (rc *ReusableCfg) Compute(
	compiler AbstractCompiler,
	cfgRoot rhino.Node,
	shouldTraverseFunctions bool,
	edgeAnnotations bool,
) ControlFlowGraph[rhino.Node] {
	jutil.CheckNotNullMsg(compiler, "Need to call setCompiler()")
	rhino.CheckNotNullMsg(cfgRoot, "Need to call setCfgRoot()")
	rc.cfa.init(compiler, shouldTraverseFunctions, edgeAnnotations)
	rc.cfa.computeCfg(cfgRoot, &rc.cfg)
	return rc.cfa.cfg
}

// release clears the graph computed last, so that rc does not keep it (or
// more storage than is worth reusing) alive once its user is done with it.
func (rc *ReusableCfg) release() {
	rc.cfg.clearForReuse()
	rc.cfa.compiler = nil
	rc.cfa.cfg = nil
	rc.cfa.root = rhino.NilNode
	rc.cfa.finallyMap = cfaEmptyFinallyMap
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
	c := &ControlFlowAnalysis{}
	c.init(compiler, shouldTraverseFunctions, edgeAnnotations)
	return c
}

// init runs the constructor on c, keeping the storage of a previous use.
func (c *ControlFlowAnalysis) init(compiler AbstractCompiler, shouldTraverseFunctions bool, edgeAnnotations bool) {
	c.compiler = compiler
	c.shouldTraverseFunctions = shouldTraverseFunctions
	c.edgeAnnotations = edgeAnnotations
	clear(c.exceptionHandler)
	c.exceptionHandler = c.exceptionHandler[:0]
	c.finallyMap = cfaEmptyFinallyMap
}

func (c *ControlFlowAnalysis) putFinally(jump rhino.Node, target rhino.Node) {
	if c.finallyMap == cfaEmptyFinallyMap {
		c.finallyMap = jutil.NewSetMultimap[rhino.Node, rhino.Node]()
	}
	c.finallyMap.Put(jump, target)
}

// astPositionOfNode is checkNotNull(astPosition.get(v), v).
func (c *ControlFlowAnalysis) astPositionOfNode(v rhino.Node) int32 {
	p, ok := c.astPosition[v]
	if !ok {
		// checkNotNull(T, Object errorMessage)
		panic(jutil.NewNullPointerExceptionMsg(jutil.ToJavaString(v)))
	}
	return p
}

// ControlFlowAnalysis_Builder configures a ControlFlowAnalysis instance then
// computes the ControlFlowGraph
type ControlFlowAnalysis_Builder struct {
	compiler                AbstractCompiler
	cfgRoot                 rhino.Node
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
func (b *ControlFlowAnalysis_Builder) SetCfgRoot(cfgRoot rhino.Node) *ControlFlowAnalysis_Builder {
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
func (b *ControlFlowAnalysis_Builder) ComputeCfg() ControlFlowGraph[rhino.Node] {
	jutil.CheckNotNullMsg(b.compiler, "Need to call setCompiler()")
	rhino.CheckNotNullMsg(b.cfgRoot, "Need to call setCfgRoot()")
	cfa := newControlFlowAnalysis(b.compiler, b.shouldTraverseFunctions, b.edgeAnnotations)
	cfa.computeCfg(b.cfgRoot, nil)
	return cfa.cfg
}

// Java: ControlFlowAnalysis.builder()
func ControlFlowAnalysis_Builder_() *ControlFlowAnalysis_Builder {
	return newControlFlowAnalysis_Builder()
}

// Java: ControlFlowAnalysis.getCfg()
func (c *ControlFlowAnalysis) GetCfg() ControlFlowGraph[rhino.Node] {
	return c.cfg
}

// computeCfg builds the CFG of root into storage, a graph to reuse, or into a
// new graph when storage is nil.
//
// Java: ControlFlowAnalysis.computeCfg(Node)
func (c *ControlFlowAnalysis) computeCfg(root rhino.Node, storage *ControlFlowAnalysis_AstControlFlowGraph) {
	jutil.CheckArgumentMsg1(
		NodeUtil_IsValidCfgRoot(root), "Unexpected control flow graph root %s", root)
	c.root = root
	if c.astPosition == nil {
		c.astPosition = make(map[rhino.Node]int32)
	}
	c.astPositionCounter = 0
	entry := ControlFlowAnalysis_ComputeFallThrough(root)
	if storage != nil {
		storage.reinit(entry, c.edgeAnnotations)
		c.cfg = storage
	} else {
		c.cfg = newControlFlowAnalysis_AstControlFlowGraph(entry, c.edgeAnnotations)
	}

	// Traverse the graph.
	c.traverse(root)

	// Insert an implicit return last.
	c.astPositionCounter++
	c.astPosition[rhino.NilNode] = c.astPositionCounter

	// Now, generate the priority of nodes by doing a depth-first
	// search on the CFG.
	c.priorityCounter = 0
	c.prioritizeFromEntryNode(c.cfg.entry.(*LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]))

	cfgNodes := &c.cfg.nodes
	if c.shouldTraverseFunctions {
		// If we're traversing inner functions, we need to rank the
		// priority of them too.
		for candidate := range cfgNodes.Values() {
			value := candidate.value
			if !value.IsNil() && value.IsFunction() {
				c.prioritizeFromEntryNode(candidate)
			}
		}
	}

	// Free our large table now that it is no longer needed
	if len(c.astPosition) > cfaReuseMaxSize {
		c.astPosition = nil
	} else {
		clear(c.astPosition)
	}
	if cap(c.worklist) > cfaReuseMaxSize {
		c.worklist = nil
	}

	// At this point, all reachable nodes have been given a priority, but
	// unreachable nodes have not been given a priority. Put them last.
	// Presumably, it doesn't really matter what priority they get, since
	// this shouldn't happen in real code.
	for candidate := range cfgNodes.Values() {
		if candidate.priority < 0 {
			c.priorityCounter++
			candidate.priority = int32(c.priorityCounter)
		}
	}

	// Again, the implicit return node is always last.
	c.priorityCounter++
	c.cfg.GetImplicitReturn().SetPriority(c.priorityCounter)
}

func (c *ControlFlowAnalysis) traverse(root rhino.Node) {
	jutil.CheckNotNull(c.compiler)
	if root.IsNil() {
		b := NodeTraversal_Builder{compiler: c.compiler, callback: newControlFlowAnalysis(c.compiler, c.shouldTraverseFunctions, c.edgeAnnotations)}
		newNodeTraversal(&b).Traverse(root)
		return
	}
	var currentNode rhino.Node
	var currentScript rhino.Node
	jutil.Try(func() {
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		currentNode = root
		c.walkBranch(root, rhino.NilNode, &currentNode, &currentScript)
	}, func(unexpectedException jutil.Throwable) {
		b := NodeTraversal_Builder{compiler: c.compiler, callback: newControlFlowAnalysis(c.compiler, c.shouldTraverseFunctions, c.edgeAnnotations)}
		t := newNodeTraversal(&b)
		t.InitTraversal(root)
		if !currentScript.IsNil() {
			t.currentScript = currentScript
		}
		t.currentNode = currentNode
		t.ThrowUnexpectedException(unexpectedException)
	})
}

func (c *ControlFlowAnalysis) walkBranch(n, parent rhino.Node, currentNode, currentScript *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		*currentNode = n
		*currentScript = n
		if c.ShouldTraverse(nil, n, parent) {
			for child := n.GetFirstChild(); !child.IsNil(); {
				next := child.GetNext()
				c.walkBranch(child, n, currentNode, currentScript)
				child = next
			}
			*currentNode = n
			c.Visit(nil, n, parent)
		}
		return
	case rhino.Token_FUNCTION:
		*currentNode = n
		if c.ShouldTraverse(nil, n, parent) {
			fnName := n.GetFirstChild()
			args := fnName.GetNext()
			body := args.GetNext()
			c.walkBranch(fnName, n, currentNode, currentScript)
			c.walkBranch(args, n, currentNode, currentScript)
			c.walkBranch(body, n, currentNode, currentScript)
			*currentNode = n
			c.Visit(nil, n, parent)
		}
		return
	case rhino.Token_CLASS:
		*currentNode = n
		if c.ShouldTraverse(nil, n, parent) {
			className := n.GetFirstChild()
			extendsClause := className.GetNext()
			body := extendsClause.GetNext()
			c.walkBranch(extendsClause, n, currentNode, currentScript)
			c.walkBranch(className, n, currentNode, currentScript)
			c.walkBranch(body, n, currentNode, currentScript)
			*currentNode = n
			c.Visit(nil, n, parent)
		}
		return
	case rhino.Token_CLASS_MEMBERS:
		*currentNode = n
		if !c.ShouldTraverse(nil, n, parent) {
			return
		}
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			if child.IsComputedProp() || child.IsComputedFieldDef() {
				c.walkBranch(child.GetFirstChild(), child, currentNode, currentScript)
			}
			child = next
		}
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			switch child.GetToken() {
			case rhino.Token_COMPUTED_PROP:
				*currentNode = n
				if c.ShouldTraverse(nil, child, n) {
					c.walkBranch(child.GetLastChild(), child, currentNode, currentScript)
					*currentNode = n
					c.Visit(nil, child, n)
				}
			case rhino.Token_COMPUTED_FIELD_DEF:
				*currentNode = n
				if c.ShouldTraverse(nil, child, n) {
					if child.HasTwoChildren() {
						c.walkBranch(child.GetLastChild(), child, currentNode, currentScript)
					}
					*currentNode = n
					c.Visit(nil, child, n)
				}
			case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_BLOCK, rhino.Token_MEMBER_FUNCTION_DEF,
				rhino.Token_MEMBER_VARIABLE_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
				c.walkBranch(child, n, currentNode, currentScript)
			default:
				panic(jutil.NewIllegalStateExceptionMsg("Invalid class member: " + child.GetToken().String()))
			}
			child = next
		}
		*currentNode = n
		c.Visit(nil, n, parent)
		return
	default:
	}

	*currentNode = n
	if !c.ShouldTraverse(nil, n, parent) {
		return
	}
	for child := n.GetFirstChild(); !child.IsNil(); {
		next := child.GetNext()
		c.walkBranch(child, n, currentNode, currentScript)
		child = next
	}
	*currentNode = n
	c.Visit(nil, n, parent)
}

type cfaWorkItem struct {
	node *LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	pos  int32
}

func (c *ControlFlowAnalysis) cfaSiftUp(worklist []cfaWorkItem, i int, x cfaWorkItem) {
	for i > 0 {
		parent := (i - 1) / 2
		p := worklist[parent]
		if x.pos >= p.pos {
			break
		}
		worklist[i] = p
		i = parent
	}
	worklist[i] = x
}

func (c *ControlFlowAnalysis) cfaSiftDown(worklist []cfaWorkItem, i int, x cfaWorkItem) {
	n := len(worklist)
	if 2*i+1 >= n {
		worklist[i] = x
		return
	}
	xPos := x.pos
	for {
		pick := 2*i + 1
		if pick >= n {
			break
		}
		smaller := worklist[pick]
		smallerPos := smaller.pos
		if sibling := pick + 1; sibling < n {
			sib := worklist[sibling]
			if sibPos := sib.pos; smallerPos > sibPos {
				pick = sibling
				smaller = sib
				smallerPos = sibPos
			}
		}
		if xPos <= smallerPos {
			break
		}
		worklist[i] = smaller
		i = pick
	}
	worklist[i] = x
}

// prioritizeFromEntryNode: given an entry node, find all the nodes reachable
// from that node and prioritize them.
//
// Java: ControlFlowAnalysis.prioritizeFromEntryNode(DiGraphNode<Node, Branch>)
func (c *ControlFlowAnalysis) prioritizeFromEntryNode(entry *LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]) {
	if entry == nil {
		panic(jutil.NewNullPointerException())
	}
	worklist := append(c.worklist[:0], cfaWorkItem{node: entry, pos: -1})

	for len(worklist) > 0 {
		n := len(worklist)
		current := worklist[0].node
		last := worklist[n-1]
		worklist[n-1] = cfaWorkItem{}
		worklist = worklist[:n-1]
		if n > 1 {
			c.cfaSiftDown(worklist, 0, last)
		}
		if current.priority >= 0 {
			continue
		}

		c.priorityCounter++
		current.priority = int32(c.priorityCounter)

		for _, edge := range current.outEdgeList.Slice() {
			succ := edge.destNode
			if succ == nil {
				panic(jutil.NewNullPointerException())
			}
			if len(worklist) == 0 {
				worklist = append(worklist, cfaWorkItem{node: succ, pos: -1})
			} else {
				if worklist[0].pos < 0 {
					worklist[0].pos = c.astPositionOfNode(worklist[0].node.value)
				}
				succPos := c.astPositionOfNode(succ.value)
				worklist = append(worklist, cfaWorkItem{})
				c.cfaSiftUp(worklist, len(worklist)-1, cfaWorkItem{node: succ, pos: succPos})
			}
		}
	}
	c.worklist = worklist[:0]
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

// Java: ControlFlowAnalysis.visit(NodeTraversal, Node, Node)
func (c *ControlFlowAnalysis) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_IF:
		c.handleIf(ast.MustIf(n))
	case rhino.Token_WHILE:
		c.handleWhile(ast.MustWhile(n))
	case rhino.Token_DO:
		c.handleDo(ast.MustDo(n))
	case rhino.Token_FOR:
		c.handleFor(ast.MustFor(n))
	case rhino.Token_FOR_OF, rhino.Token_FOR_IN, rhino.Token_FOR_AWAIT_OF:
		c.handleEnhancedFor(ast.MustEnhancedFor(n))
	case rhino.Token_SWITCH:
		c.handleSwitch(ast.MustSwitch(n))
	case rhino.Token_CASE:
		c.handleCase(ast.MustCase(n))
	case rhino.Token_DEFAULT_CASE:
		c.handleDefault(ast.MustDefaultCase(n))
	case rhino.Token_BLOCK, rhino.Token_ROOT, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY:
		c.handleStmtList(n)
	case rhino.Token_FUNCTION:
		c.handleFunction(ast.MustFunction(n))
	case rhino.Token_EXPR_RESULT:
		c.handleExpr(ast.MustExprStmt(n))
	case rhino.Token_THROW:
		c.handleThrow(ast.MustThrow(n))
	case rhino.Token_TRY:
		c.handleTry(ast.MustTry(n))
	case rhino.Token_CATCH:
		c.handleCatch(ast.MustCatch(n))
	case rhino.Token_BREAK:
		c.handleBreak(ast.MustBreak(n))
	case rhino.Token_CONTINUE:
		c.handleContinue(ast.MustContinue(n))
	case rhino.Token_RETURN:
		c.handleReturn(ast.MustReturn(n))
	case rhino.Token_WITH:
		c.handleWith(ast.MustWith(n))
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

// The handlers below connect a statement to the blocks it runs (the
// branches of an if, the body of a loop...) directly: control falls through
// into a block at the block itself (see computeFallThrough), so Java's
// computeFallThrough(block) calls are left out.

// Java: ControlFlowAnalysis.handleIf(Node)
func (c *ControlFlowAnalysis) handleIf(ifs ast.If) {
	n := ifs.Node()
	c.createEdge(n, ControlFlowGraph_Branch_ON_TRUE, ifs.Then().Node())

	if els, ok := ifs.Else(); ok {
		c.createEdge(n, ControlFlowGraph_Branch_ON_FALSE, els.Node())
	} else {
		c.createEdge(n, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(n, c)) // not taken branch
	}
	c.connectToPossibleExceptionHandler(n, ifs.Cond().Node())
}

// Java: ControlFlowAnalysis.handleWhile(Node)
func (c *ControlFlowAnalysis) handleWhile(while ast.While) {
	n := while.Node()
	// Control goes to the first statement if the condition evaluates to true.
	c.createEdge(n, ControlFlowGraph_Branch_ON_TRUE, while.Body().Node())

	if while.Cond().Token() != rhino.Token_TRUE {
		// Control goes to the follow() if the condition evaluates to false.
		c.createEdge(n, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(n, c))
	}
	c.connectToPossibleExceptionHandler(n, while.Cond().Node())
}

// Java: ControlFlowAnalysis.handleDo(Node)
func (c *ControlFlowAnalysis) handleDo(do ast.Do) {
	n := do.Node()
	// The first edge can be the initial iteration as well as the iterations
	// after.
	c.createEdge(n, ControlFlowGraph_Branch_ON_TRUE, do.Body().Node())
	if do.Cond().Token() != rhino.Token_TRUE {
		// The edge that leaves the do loop if the condition fails.
		c.createEdge(n, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(n, c))
	}
	c.connectToPossibleExceptionHandler(n, do.Cond().Node())
}

// Java: ControlFlowAnalysis.handleEnhancedFor(Node)
func (c *ControlFlowAnalysis) handleEnhancedFor(forNode ast.EnhancedFor) {
	// We have:  for (index in collection) { body }
	// or:       for (item of collection) { body }
	// or:       for await (item of collection) { body }
	n := forNode.Node()
	collection := forNode.Collection().Node()
	// The collection behaves like init.
	c.createEdge(collection, ControlFlowGraph_Branch_UNCOND, n)
	// The edge that transfer control to the beginning of the loop body.
	c.createEdge(n, ControlFlowGraph_Branch_ON_TRUE, forNode.Body().Node())
	// The edge to end of the loop.
	c.createEdge(n, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(n, c))
	if forNode.Token() == rhino.Token_FOR_IN {
		c.connectToPossibleExceptionHandler(n, collection)
	} else {
		// A for-of or for-await-of loop calls the iterator, which may throw.
		c.connectToPossibleExceptionHandlerUnconditionally(n)
	}
}

// Java: ControlFlowAnalysis.handleFor(Node)
func (c *ControlFlowAnalysis) handleFor(forNode ast.For) {
	// We have for (init; cond; iter) { body }
	n := forNode.Node()
	init := cfaForInit(forNode)
	iter := cfaForUpdate(forNode)
	// After initialization, we transfer to the FOR which is in charge of
	// checking the condition (for the first time).
	c.createEdge(init, ControlFlowGraph_Branch_UNCOND, n)
	// The edge that transfer control to the beginning of the loop body.
	c.createEdge(n, ControlFlowGraph_Branch_ON_TRUE, forNode.Body().Node())
	// The edge to end of the loop.
	cond, hasCond := forNode.Cond()
	if hasCond && cond.Token() != rhino.Token_TRUE {
		c.createEdge(n, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(n, c))
	}
	// The end of the body will have a unconditional branch to our iter
	// (handled by calling computeFollowNode of the last instruction of the
	// body. Our iter will jump to the forNode again to another condition
	// check.
	c.createEdge(iter, ControlFlowGraph_Branch_UNCOND, n)
	c.connectToPossibleExceptionHandler(init, init)
	if hasCond {
		c.connectToPossibleExceptionHandler(n, cond.Node())
	}
	c.connectToPossibleExceptionHandler(iter, iter)
}

// cfaForInit returns the init of a for loop as a node of the graph: the
// EMPTY node in its place if the loop has none (For.Init reports it as
// missing, but the graph has a node for it).
func cfaForInit(forNode ast.For) rhino.Node { return forNode.Node().GetFirstChild() }

// cfaForUpdate returns the update of a for loop as a node of the graph: the
// EMPTY node in its place if the loop has none (see cfaForInit).
func cfaForUpdate(forNode ast.For) rhino.Node { return forNode.Node().GetChildAtIndex(2) }

// Java: ControlFlowAnalysis.handleSwitch(Node)
func (c *ControlFlowAnalysis) handleSwitch(sw ast.Switch) {
	// Transfer to the first non-DEFAULT CASE. if there are none, transfer
	// to the DEFAULT.
	//
	// (Java also looks for an EMPTY node, which a SWITCH_BODY never holds:
	// its children are cases.)
	n := sw.Node()
	switchBody := sw.Body()
	c.createEdge(n, ControlFlowGraph_Branch_UNCOND, switchBody.Node())
	if next, ok := cfaNextCase(switchBody.FirstCase()); ok { // Has at least one CASE
		c.createEdge(switchBody.Node(), ControlFlowGraph_Branch_UNCOND, next.Node())
	} else if deflt, ok := switchBody.FirstCase(); ok { // Has no CASE, so this is the DEFAULT
		c.createEdge(switchBody.Node(), ControlFlowGraph_Branch_UNCOND, deflt.Node())
	} else { // No CASE, no DEFAULT
		c.createEdge(switchBody.Node(), ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFollowNodeWithCfa(n, c))
	}
	c.connectToPossibleExceptionHandler(n, sw.Discriminant().Node())
}

// Java: ControlFlowAnalysis.handleCase(Node)
func (c *ControlFlowAnalysis) handleCase(caseNode ast.Case) {
	n := caseNode.Node()
	// Case is a bit tricky....First it goes into the body if condition is true.
	c.createEdge(n, ControlFlowGraph_Branch_ON_TRUE, caseNode.Body().Node())

	// Look for the next CASE, skipping over DEFAULT.
	if next, ok := cfaNextCase(caseNode.AsSwitchCase().Next()); ok { // Found a CASE
		c.createEdge(n, ControlFlowGraph_Branch_ON_FALSE, next.Node())
	} else if deflt, ok := cfaDefaultCase(ast.MustSwitchBody(n.GetParent())); ok {
		// No more CASE found, but the switch has a DEFAULT.
		c.createEdge(n, ControlFlowGraph_Branch_ON_FALSE, deflt.Node())
	} else { // No DEFAULT found, go to the follow of the SWITCH.
		c.createEdge(n, ControlFlowGraph_Branch_ON_FALSE, ControlFlowAnalysis_ComputeFollowNodeWithCfa(n, c))
	}
	c.connectToPossibleExceptionHandler(n, caseNode.Value().Node())
}

// cfaNextCase returns the first CASE from sc on (sc itself if it is one),
// skipping DEFAULT; ok is false if there is none. It takes what a
// FirstCase() or Next() returns.
//
// Java: getNextSiblingOfType(sc, Token.CASE)
func cfaNextCase(sc ast.SwitchCase, ok bool) (ast.Case, bool) {
	for ; ok; sc, ok = sc.Next() {
		if caseNode, isCase := sc.AsCase(); isCase {
			return caseNode, true
		}
	}
	return ast.Case{}, false
}

// cfaDefaultCase returns the DEFAULT of a switch; ok is false if it has
// none.
//
// Java: getNextSiblingOfType(switchBody.getFirstChild(), Token.DEFAULT_CASE)
func cfaDefaultCase(switchBody ast.SwitchBody) (ast.DefaultCase, bool) {
	for sc := range switchBody.Cases() {
		if deflt, ok := sc.AsDefaultCase(); ok {
			return deflt, true
		}
	}
	return ast.DefaultCase{}, false
}

// Java: ControlFlowAnalysis.handleDefault(Node)
func (c *ControlFlowAnalysis) handleDefault(deflt ast.DefaultCase) {
	// Directly goes to the body. It should not transfer to the next case.
	c.createEdge(deflt.Node(), ControlFlowGraph_Branch_UNCOND, deflt.Body().Node())
}

// Java: ControlFlowAnalysis.handleWith(Node)
func (c *ControlFlowAnalysis) handleWith(with ast.With) {
	// Directly goes to the body. It should not transfer to the next case.
	c.createEdge(with.Node(), ControlFlowGraph_Branch_UNCOND, with.Body().Node())
	c.connectToPossibleExceptionHandler(with.Node(), with.Object().Node())
}

// handleStmtList handles the nodes that hold a list of statements: a block,
// a script, a module body, a root (which holds scripts), and the block that
// holds the catch clause of a TRY. They have no common view, so node stays
// a rhino.Node.
//
// Java: ControlFlowAnalysis.handleStmtList(Node)
func (c *ControlFlowAnalysis) handleStmtList(node rhino.Node) {
	parent := node.GetParent()
	// Special case, don't add a block of empty CATCH block to the graph.
	if try, ok := ast.AsTry(parent); ok && try.Catches().Is(node) {
		if _, hasCatch := try.Catch(); !hasCatch {
			return
		}
	}

	// A block transfer control to its first child if it is not empty.
	child := node.GetFirstChild()

	// Function declarations are skipped since control doesn't go into that
	// function (unless it is called)
	for !child.IsNil() && child.IsFunction() {
		child = child.GetNext()
	}

	if !child.IsNil() {
		c.createEdge(node, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFallThrough(child))
	} else {
		c.createEdge(node, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c))
	}

	// Synthetic blocks
	if !parent.IsNil() {
		switch parent.GetToken() {
		case rhino.Token_DEFAULT_CASE, rhino.Token_CASE, rhino.Token_TRY:
		case rhino.Token_ROOT:
			// TODO(b/71873602): why is this path necessary?
			if node.IsRoot() && !node.GetNext().IsNil() {
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
func (c *ControlFlowAnalysis) handleFunction(fn ast.Function) {
	// A block transfer control to its first child if it is not empty.
	//
	// (The body is a block, the expression of an arrow function, or EMPTY in
	// a type summary: control falls through into each at the node itself.)
	c.createEdge(fn.Node(), ControlFlowGraph_Branch_UNCOND, fn.Body().Node())
	jutil.CheckState(c.exceptionHandlerPeek() == fn.Node())
	c.exceptionHandlerPop()
}

// Java: ControlFlowAnalysis.handleExpr(Node)
func (c *ControlFlowAnalysis) handleExpr(expr ast.ExprStmt) {
	n := expr.Node()
	c.createEdge(n, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFollowNodeWithCfa(n, c))
	c.connectToPossibleExceptionHandler(n, n)
}

// Java: ControlFlowAnalysis.handleThrow(Node)
func (c *ControlFlowAnalysis) handleThrow(throw ast.Throw) {
	c.connectToPossibleExceptionHandler(throw.Node(), throw.Node())
}

// Java: ControlFlowAnalysis.handleTry(Node)
func (c *ControlFlowAnalysis) handleTry(try ast.Try) {
	c.createEdge(try.Node(), ControlFlowGraph_Branch_UNCOND, try.Body().Node())
}

// Java: ControlFlowAnalysis.handleCatch(Node)
func (c *ControlFlowAnalysis) handleCatch(catch ast.Catch) {
	c.createEdge(catch.Node(), ControlFlowGraph_Branch_UNCOND, catch.Body().Node())
}

// cfaJumpLabel returns the label of a break or continue, or nil if it has
// none.
func cfaJumpLabel(jump ast.Jump) *string {
	name, ok := jump.Label()
	if !ok {
		return nil
	}
	s := name.Ident()
	return &s
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
	label := cfaJumpLabel(brk.AsJump())
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
	for cur, lastJump = node, node; !ControlFlowAnalysis_IsBreakTarget(cur, label); cur, parent = parent, parent.GetParent() {
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
	label := cfaJumpLabel(cont.AsJump())
	var cur rhino.Node
	var previous rhino.Node
	var lastJump rhino.Node

	// Similar to handBreak's logic with a few minor variation.
	for cur, lastJump = node, node; !ControlFlowAnalysis_IsContinueTarget(cur, label); cur = cur.GetParent() {
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

// Java: ControlFlowAnalysis.handleStmt(Node)
func (c *ControlFlowAnalysis) handleStmt(node rhino.Node) {
	// Simply transfer to the next line.
	c.createEdge(node, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c))
	c.connectToPossibleExceptionHandler(node, node)
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
			if cfa != nil && cfa.finallyMap != nil && cfa.finallyMap.ContainsKey(parent) {
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
	source *LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
	branch ControlFlowGraph_Branch,
	dest *LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
) bool {
	outEdges := source.outEdgeList.Slice()
	inEdges := dest.inEdgeList.Slice()
	if len(outEdges) < len(inEdges) {
		for _, outEdge := range outEdges {
			if outEdge.GetDestination() == DiGraph_DiGraphNode[rhino.Node, ControlFlowGraph_Branch](dest) && outEdge.GetValue() == branch {
				return true
			}
		}
	} else {
		for _, inEdge := range inEdges {
			if inEdge.GetSource() == DiGraph_DiGraphNode[rhino.Node, ControlFlowGraph_Branch](source) && inEdge.GetValue() == branch {
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
	base := &c.cfg.GetControlFlowGraphBase().LinkedDirectedGraph
	from := base.CreateLinkedDiGraphNode(fromNode)
	to := base.CreateLinkedDiGraphNode(toNode)
	if !cfaHasEdgeInDirection(from, branch, to) && !cfaHasEdgeInDirection(to, branch, from) {
		base.connectLinkedNodes(from, branch, to)
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

// ControlFlowAnalysis_IsBreakTarget checks if target is actually the break
// target of labeled continue. The label can be null if it is an unlabeled
// break.
//
// Java: ControlFlowAnalysis.isBreakTarget(Node, String)
func ControlFlowAnalysis_IsBreakTarget(target rhino.Node, label *string) bool {
	return ControlFlowAnalysis_IsBreakStructure(target, label != nil) && ControlFlowAnalysis_MatchLabel(target.GetParent(), label)
}

// ControlFlowAnalysis_IsContinueTarget checks if target is actually the
// continue target of labeled continue. The label can be null if it is an
// unlabeled continue.
//
// Java: ControlFlowAnalysis.isContinueTarget(Node, String)
func ControlFlowAnalysis_IsContinueTarget(target rhino.Node, label *string) bool {
	return NodeUtil_IsLoopStructure(target) && ControlFlowAnalysis_MatchLabel(target.GetParent(), label)
}

// ControlFlowAnalysis_MatchLabel checks if label is actually referencing the
// target control structure. If label is null, it always returns true.
//
// Java: ControlFlowAnalysis.matchLabel(Node, String)
func ControlFlowAnalysis_MatchLabel(target rhino.Node, label *string) bool {
	if label == nil {
		return true
	}
	// A statement can have several labels (a: b: while (x) ...).
	for l, ok := ast.AsLabel(target); ok; l, ok = ast.AsLabel(l.Node().GetParent()) {
		if l.Name().Ident() == *label {
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

// ControlFlowAnalysis_AstControlFlowGraph is a ControlFlowGraph which
// provides a node comparator based on the pre-order traversal of the AST.
type ControlFlowAnalysis_AstControlFlowGraph struct {
	ControlFlowGraphBase[rhino.Node]
}

// newControlFlowAnalysis_AstControlFlowGraph constructs this graph.
//
// entry is the entry node; edgeAnnotations is the annotation.
//
// Java: ControlFlowAnalysis.AstControlFlowGraph(Node, boolean)
func newControlFlowAnalysis_AstControlFlowGraph(entry rhino.Node, edgeAnnotations bool) *ControlFlowAnalysis_AstControlFlowGraph {
	g := &ControlFlowAnalysis_AstControlFlowGraph{}
	g.InitControlFlowGraphBase(g, entry /* nodeAnnotations= */, true, edgeAnnotations)
	return g
}

// reinit makes g a new AstControlFlowGraph(entry, edgeAnnotations), reusing
// its storage.
//
// Go-only: performance — see ReusableCfg.
func (g *ControlFlowAnalysis_AstControlFlowGraph) reinit(entry rhino.Node, edgeAnnotations bool) {
	g.clearForReuse()
	g.InitControlFlowGraphBase(g, entry /* nodeAnnotations= */, true, edgeAnnotations)
}

func cfaForwardNodeComparator(a, b DiGraph_DiGraphNode[rhino.Node, ControlFlowGraph_Branch]) int {
	return cmp.Compare(a.GetPriority(), b.GetPriority())
}

func cfaBackwardNodeComparator(a, b DiGraph_DiGraphNode[rhino.Node, ControlFlowGraph_Branch]) int {
	return cmp.Compare(b.GetPriority(), a.GetPriority())
}

// GetOptionalNodeComparator returns a node comparator based on the
// pre-order traversal of the AST.
//
// isForward: x 'before' y in the pre-order traversal implies x 'less than'
// y (if true) and x 'greater than' y (if false).
//
// Java: ControlFlowAnalysis.AstControlFlowGraph.getOptionalNodeComparator(boolean)
func (g *ControlFlowAnalysis_AstControlFlowGraph) GetOptionalNodeComparator(isForward bool) func(a, b DiGraph_DiGraphNode[rhino.Node, ControlFlowGraph_Branch]) int {
	if isForward {
		return cfaForwardNodeComparator
	}
	// Comparator.reversed()
	return cfaBackwardNodeComparator
}

var (
	_ NodeTraversal_Callback       = (*ControlFlowAnalysis)(nil)
	_ ControlFlowGraph[rhino.Node] = (*ControlFlowAnalysis_AstControlFlowGraph)(nil)
)

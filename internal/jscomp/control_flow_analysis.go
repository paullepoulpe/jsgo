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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/graph"
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

	cfg *ControlFlowGraph[rhino.Node]

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
	// Insertion-ordered unique targets per jump node (avoids SetMultimap heap allocations).
	finallyMap cfaFinallyMap
}

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
	cfg ControlFlowGraph[rhino.Node]
}

// Compute is ControlFlowAnalysis.builder() with these settings followed by
// computeCfg(), into rc's storage.
func (rc *ReusableCfg) Compute(
	compiler AbstractCompiler,
	cfgRoot rhino.Node,
	shouldTraverseFunctions bool,
	edgeAnnotations bool,
) *ControlFlowGraph[rhino.Node] {
	jutil.CheckNotNullMsg(compiler, "Need to call setCompiler()")
	rhino.CheckNotNullMsg(cfgRoot, "Need to call setCfgRoot()")
	rc.cfa.init(compiler, shouldTraverseFunctions, edgeAnnotations)
	rc.cfa.computeCfg(cfgRoot, &rc.cfg)
	return rc.cfa.cfg
}

// release clears the graph computed last, so that rc does not keep it (or
// more storage than is worth reusing) alive once its user is done with it.
func (rc *ReusableCfg) release() {
	rc.cfg.ClearForReuse()
	rc.cfa.compiler = nil
	rc.cfa.cfg = nil
	rc.cfa.root = rhino.NilNode
	rc.cfa.finallyMap.reset()
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
	c.finallyMap.reset()
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
func (b *ControlFlowAnalysis_Builder) ComputeCfg() *ControlFlowGraph[rhino.Node] {
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
func (c *ControlFlowAnalysis) GetCfg() *ControlFlowGraph[rhino.Node] {
	return c.cfg
}

// computeCfg builds the CFG of root into storage, a graph to reuse, or into a
// new graph when storage is nil.
//
// Java: ControlFlowAnalysis.computeCfg(Node)
func (c *ControlFlowAnalysis) computeCfg(root rhino.Node, storage *ControlFlowGraph[rhino.Node]) {
	jutil.CheckArgumentMsg1(
		NodeUtil_IsValidCfgRoot(root), "Unexpected control flow graph root %s", root)
	c.root = root
	if c.astPosition == nil {
		c.astPosition = make(map[rhino.Node]int32)
	}
	c.astPositionCounter = 0
	entry := ControlFlowAnalysis_ComputeFallThrough(root)
	if storage != nil {
		storage.reinitAst(entry, c.edgeAnnotations)
		c.cfg = storage
	} else {
		c.cfg = newAstControlFlowGraph(entry, c.edgeAnnotations)
	}

	// Traverse the graph.
	c.traverse(root)

	// Insert an implicit return last.
	c.astPositionCounter++
	c.astPosition[rhino.NilNode] = c.astPositionCounter

	// Now, generate the priority of nodes by doing a depth-first
	// search on the CFG.
	c.priorityCounter = 0
	c.prioritizeFromEntryNode(c.cfg.entry)

	cfgNodes := &c.cfg.Nodes
	if c.shouldTraverseFunctions {
		// If we're traversing inner functions, we need to rank the
		// priority of them too.
		for candidate := range cfgNodes.Values() {
			value := candidate.Value
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
		if candidate.Priority < 0 {
			c.priorityCounter++
			candidate.Priority = int32(c.priorityCounter)
		}
	}

	// Again, the implicit return node is always last.
	c.priorityCounter++
	c.cfg.implicitReturn.SetPriority(c.priorityCounter)
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
	defer func() {
		if r := recover(); r != nil {
			b := NodeTraversal_Builder{compiler: c.compiler, callback: newControlFlowAnalysis(c.compiler, c.shouldTraverseFunctions, c.edgeAnnotations)}
			t := newNodeTraversal(&b)
			t.InitTraversal(root)
			if !currentScript.IsNil() {
				t.currentScript = currentScript
			}
			t.currentNode = currentNode
			t.ThrowUnexpectedException(jutil.AsThrowable(r))
		}
	}()
	if Platform_IsThreadInterrupted() {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
	}
	currentNode = root
	c.walkBranch(root, rhino.NilNode, &currentNode, &currentScript)
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
	node *graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
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
func (c *ControlFlowAnalysis) prioritizeFromEntryNode(entry *graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]) {
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
		if current.Priority >= 0 {
			continue
		}

		c.priorityCounter++
		current.Priority = int32(c.priorityCounter)

		for _, edge := range current.OutEdgeList.Slice() {
			succ := edge.DestNode
			if succ == nil {
				panic(jutil.NewNullPointerException())
			}
			if len(worklist) == 0 {
				worklist = append(worklist, cfaWorkItem{node: succ, pos: -1})
			} else {
				if worklist[0].pos < 0 {
					worklist[0].pos = c.astPositionOfNode(worklist[0].node.Value)
				}
				succPos := c.astPositionOfNode(succ.Value)
				worklist = append(worklist, cfaWorkItem{})
				c.cfaSiftUp(worklist, len(worklist)-1, cfaWorkItem{node: succ, pos: succPos})
			}
		}
	}
	c.worklist = worklist[:0]
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

// Java: ControlFlowAnalysis.handleStmt(Node)
func (c *ControlFlowAnalysis) handleStmt(node rhino.Node) {
	// Simply transfer to the next line.
	c.createEdge(node, ControlFlowGraph_Branch_UNCOND, ControlFlowAnalysis_ComputeFollowNodeWithCfa(node, c))
	c.connectToPossibleExceptionHandler(node, node)
}

var (
	_ NodeTraversal_Callback = (*ControlFlowAnalysis)(nil)
)

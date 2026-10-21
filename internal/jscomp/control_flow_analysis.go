// Port of com/google/javascript/jscomp/ControlFlowAnalysis.java @ closure-compiler v20261015

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
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type cfaPosEntry struct {
	node *rhino.Node
	pos  int32
	gen  uint32
}

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

	cfg        ControlFlowGraph[*rhino.Node]
	cfgStorage *ControlFlowAnalysis_AstControlFlowGraph

	posTable []cfaPosEntry
	posGen   uint32
	posCount int32
	nilPos   int32

	astPositionCounter int32
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
	exceptionHandler    []*rhino.Node
	exceptionHandlerBuf [4]*rhino.Node

	worklist    []cfaWorkItem
	worklistBuf [16]cfaWorkItem

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

var cfaEmptyFinallyMap = jutil.NewSetMultimap[*rhino.Node, *rhino.Node]()

var cfaRootNodesPool = sync.Pool{
	New: func() any {
		return *jutil.NewLinkedHashMapWithCapacity[*rhino.Node, *LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch]](4096)
	},
}

// ReusableCfg pairs a reusable ControlFlowAnalysis and AstControlFlowGraph.
type ReusableCfg struct {
	cfa ControlFlowAnalysis
	cfg ControlFlowAnalysis_AstControlFlowGraph
}

var reusableCfgPool = sync.Pool{
	New: func() any {
		return &ReusableCfg{}
	},
}

var cachedReusableCfgs [4]atomic.Pointer[ReusableCfg]

func evictCachedReusableCfgs() {
	for i := range cachedReusableCfgs {
		cachedReusableCfgs[i].Store(nil)
	}
}

func acquireReusableCfg() *ReusableCfg {
	for {
		bestIdx := -1
		bestChunks := -1
		for i := range cachedReusableCfgs {
			if rc := cachedReusableCfgs[i].Load(); rc != nil {
				if n := len(rc.cfg.nodeChunks); n > bestChunks {
					bestChunks = n
					bestIdx = i
				}
			}
		}
		if bestIdx < 0 {
			break
		}
		if rc := cachedReusableCfgs[bestIdx].Swap(nil); rc != nil {
			return rc
		}
	}
	if v := reusableCfgPool.Get(); v != nil {
		return v.(*ReusableCfg)
	}
	return &ReusableCfg{}
}

func releaseReusableCfg(rc *ReusableCfg) {
	if rc == nil {
		return
	}
	rc.cfg.ReleaseRootNodes()
	rc.cfg.ClearSlabsForPool()
	rc.cfg.self = nil
	rc.cfg.entry = nil
	rc.cfg.implicitReturn = nil
	rc.cfa.compiler = nil
	rc.cfa.root = nil
	rc.cfa.cfg = nil
	rc.cfa.cfgStorage = nil
	rc.cfa.finallyMap = cfaEmptyFinallyMap
	clear(rc.cfa.exceptionHandler)
	rc.cfa.exceptionHandler = rc.cfa.exceptionHandler[:0]
	clear(rc.cfa.worklist)
	rc.cfa.worklist = rc.cfa.worklist[:0]
	if len(rc.cfa.posTable) > 8192 {
		rc.cfa.posTable = nil
		rc.cfa.posGen = 0
		rc.cfa.posCount = 0
	} else if len(rc.cfa.posTable) > 0 {
		clear(rc.cfa.posTable)
		rc.cfa.posGen = 0
		rc.cfa.posCount = 0
	}
	for i := range cachedReusableCfgs {
		if cachedReusableCfgs[i].CompareAndSwap(nil, rc) {
			return
		}
	}
	minIdx := -1
	minChunks := len(rc.cfg.nodeChunks)
	for i := range cachedReusableCfgs {
		if old := cachedReusableCfgs[i].Load(); old != nil {
			if n := len(old.cfg.nodeChunks); n < minChunks {
				minChunks = n
				minIdx = i
			}
		}
	}
	if minIdx >= 0 {
		if old := cachedReusableCfgs[minIdx].Swap(rc); old != nil {
			reusableCfgPool.Put(old)
			return
		}
	}
	reusableCfgPool.Put(rc)
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
		compiler:                compiler,
		shouldTraverseFunctions: shouldTraverseFunctions,
		edgeAnnotations:         edgeAnnotations,
		finallyMap:              cfaEmptyFinallyMap,
	}
	c.exceptionHandler = c.exceptionHandlerBuf[:0]
	return c
}

func (c *ControlFlowAnalysis) putFinally(jump *rhino.Node, target *rhino.Node) {
	if c.finallyMap == nil || c.finallyMap == cfaEmptyFinallyMap {
		c.finallyMap = jutil.NewSetMultimap[*rhino.Node, *rhino.Node]()
	}
	c.finallyMap.Put(jump, target)
}

func (c *ControlFlowAnalysis) resetPosTable() {
	c.posGen++
	if c.posGen == 0 {
		clear(c.posTable)
		c.posGen = 1
	}
	c.posCount = 0
	c.nilPos = -1
	if len(c.posTable) == 0 {
		c.posTable = make([]cfaPosEntry, 256)
	}
}

func (c *ControlFlowAnalysis) growPosTable() {
	newTable := make([]cfaPosEntry, len(c.posTable)*2)
	mask := uint32(len(newTable) - 1)
	gen := c.posGen
	for _, e := range c.posTable {
		if e.gen == gen {
			idx := (uint32(uintptr(unsafe.Pointer(e.node))>>4) * 0x9e3779b1) & mask
			for newTable[idx].gen == gen {
				idx = (idx + 1) & mask
			}
			newTable[idx] = e
		}
	}
	c.posTable = newTable
}

func (c *ControlFlowAnalysis) setAstPosition(n *rhino.Node, pos int32) {
	if n == nil {
		c.nilPos = pos
		return
	}
	if int(c.posCount+1)*4 >= len(c.posTable)*3 {
		c.growPosTable()
	}
	mask := uint32(len(c.posTable) - 1)
	idx := (uint32(uintptr(unsafe.Pointer(n))>>4) * 0x9e3779b1) & mask
	gen := c.posGen
	for {
		e := c.posTable[idx]
		if e.gen != gen {
			c.posTable[idx] = cfaPosEntry{node: n, pos: pos, gen: gen}
			c.posCount++
			return
		}
		if e.node == n {
			c.posTable[idx].pos = pos
			return
		}
		idx = (idx + 1) & mask
	}
}

func (c *ControlFlowAnalysis) getAstPosition(n *rhino.Node) int32 {
	if n == nil {
		return c.nilPos
	}
	mask := uint32(len(c.posTable) - 1)
	idx := (uint32(uintptr(unsafe.Pointer(n))>>4) * 0x9e3779b1) & mask
	gen := c.posGen
	for {
		e := c.posTable[idx]
		if e.gen != gen {
			return -1
		}
		if e.node == n {
			return e.pos
		}
		idx = (idx + 1) & mask
	}
}

func (c *ControlFlowAnalysis) astPositionOfNode(v *rhino.Node) int32 {
	p := c.getAstPosition(v)
	if p < 0 {
		// checkNotNull(T, Object errorMessage)
		panic(jutil.NewNullPointerExceptionMsg(jutil.ToJavaString(v)))
	}
	return p
}

// astPositionOf is the key extractor of priorityComparator:
// checkNotNull(astPosition.get(digraphNode.getValue()), digraphNode.getValue()).
func (c *ControlFlowAnalysis) astPositionOf(digraphNode DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) int32 {
	return c.astPositionOfNode(digraphNode.GetValue())
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
	var cfa ControlFlowAnalysis
	cfa.compiler = b.compiler
	cfa.shouldTraverseFunctions = b.shouldTraverseFunctions
	cfa.edgeAnnotations = b.edgeAnnotations
	cfa.finallyMap = cfaEmptyFinallyMap
	cfa.exceptionHandler = cfa.exceptionHandlerBuf[:0]
	cfa.computeCfg(b.cfgRoot)
	return cfa.cfg
}

// Java: ControlFlowAnalysis.builder()
func ControlFlowAnalysis_Builder_() *ControlFlowAnalysis_Builder {
	return newControlFlowAnalysis_Builder()
}

// ComputeReusableCfg computes a CFG reusing c and cfgStorage across calls.
func (c *ControlFlowAnalysis) ComputeReusableCfg(
	compiler AbstractCompiler,
	cfgRoot *rhino.Node,
	shouldTraverseFunctions bool,
	edgeAnnotations bool,
	cfgStorage *ControlFlowAnalysis_AstControlFlowGraph,
) ControlFlowGraph[*rhino.Node] {
	jutil.CheckNotNullMsg(compiler, "Need to call setCompiler()")
	jutil.CheckNotNullMsg(cfgRoot, "Need to call setCfgRoot()")
	c.compiler = compiler
	c.shouldTraverseFunctions = shouldTraverseFunctions
	c.edgeAnnotations = edgeAnnotations
	c.cfgStorage = cfgStorage
	c.finallyMap = cfaEmptyFinallyMap
	if cap(c.exceptionHandler) == 0 {
		c.exceptionHandler = c.exceptionHandlerBuf[:0]
	} else {
		clear(c.exceptionHandler)
		c.exceptionHandler = c.exceptionHandler[:0]
	}
	c.computeCfg(cfgRoot)
	return c.cfg
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
	isRoot := root.IsRoot()
	c.resetPosTable()
	c.astPositionCounter = 0
	entryNode := ControlFlowAnalysis_ComputeFallThrough(root)
	if c.cfgStorage != nil {
		c.cfgStorage.Reset(entryNode, c.edgeAnnotations)
		c.cfg = c.cfgStorage
	} else {
		c.cfg = newControlFlowAnalysis_AstControlFlowGraph(entryNode, c.edgeAnnotations)
	}

	// Traverse the graph.
	c.traverse(root)

	// Insert an implicit return last.
	c.astPositionCounter++
	c.setAstPosition(nil, c.astPositionCounter)

	// Now, generate the priority of nodes by doing a depth-first
	// search on the CFG.
	c.priorityCounter = 0
	entry := c.cfg.GetEntry().(*LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch])
	c.prioritizeFromEntryNode(entry)

	cfgNodes := &c.cfg.GetControlFlowGraphBase().nodes
	if c.shouldTraverseFunctions {
		// If we're traversing inner functions, we need to rank the
		// priority of them too.
		for candidate := range cfgNodes.Values() {
			value := candidate.value
			if value != nil && value.IsFunction() {
				c.prioritizeFromEntryNode(candidate)
			}
		}
	}

	// Free our large table now that it is no longer needed
	c.astPositionCounter = 0

	if c.cfgStorage != nil {
		if isRoot {
			c.cfgStorage.wasRootCfg = true
			c.cfgStorage.pendingRootNodes = *cfgNodes
		} else if cfgNodes.Size() <= 256 {
			c.cfgStorage.savedNodes = *cfgNodes
		}
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

func (c *ControlFlowAnalysis) traverse(root *rhino.Node) {
	jutil.CheckNotNull(c.compiler)
	if root == nil {
		b := NodeTraversal_Builder{compiler: c.compiler, callback: newControlFlowAnalysis(c.compiler, c.shouldTraverseFunctions, c.edgeAnnotations)}
		newNodeTraversal(&b).Traverse(root)
		return
	}
	var currentNode *rhino.Node
	var currentScript *rhino.Node
	jutil.Try(func() {
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		currentNode = root
		c.walkBranch(root, nil, &currentNode, &currentScript)
	}, func(unexpectedException jutil.Throwable) {
		b := NodeTraversal_Builder{compiler: c.compiler, callback: newControlFlowAnalysis(c.compiler, c.shouldTraverseFunctions, c.edgeAnnotations)}
		t := newNodeTraversal(&b)
		t.InitTraversal(root)
		if currentScript != nil {
			t.currentScript = currentScript
		}
		t.currentNode = currentNode
		t.ThrowUnexpectedException(unexpectedException)
	})
}

func (c *ControlFlowAnalysis) walkBranch(n, parent *rhino.Node, currentNode, currentScript **rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		*currentNode = n
		*currentScript = n
		if c.ShouldTraverse(nil, n, parent) {
			for child := n.GetFirstChild(); child != nil; {
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
			if parent != nil {
				_ = NodeUtil_IsFunctionDeclaration(n)
			}
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
			_ = NodeUtil_IsClassExpression(n)
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
		for child := n.GetFirstChild(); child != nil; {
			next := child.GetNext()
			if child.IsComputedProp() || child.IsComputedFieldDef() {
				c.walkBranch(child.GetFirstChild(), child, currentNode, currentScript)
			}
			child = next
		}
		for child := n.GetFirstChild(); child != nil; {
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
	for child := n.GetFirstChild(); child != nil; {
		next := child.GetNext()
		c.walkBranch(child, n, currentNode, currentScript)
		child = next
	}
	*currentNode = n
	c.Visit(nil, n, parent)
}

type cfaWorkItem struct {
	node *LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch]
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
func (c *ControlFlowAnalysis) prioritizeFromEntryNode(entry *LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) {
	if entry == nil {
		panic(jutil.NewNullPointerException())
	}
	var worklist []cfaWorkItem
	if cap(c.worklist) == 0 {
		worklist = c.worklistBuf[:1]
	} else {
		worklist = c.worklist[:1]
	}
	worklist[0] = cfaWorkItem{node: entry, pos: -1}

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
			succ := edge.(*LinkedDirectedGraph_LinkedDiGraphEdge[*rhino.Node, ControlFlowGraph_Branch]).destNode
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
func (c *ControlFlowAnalysis) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if c.shouldTraverseIntoChildren(n, parent) {
		// Any AST node that will later have a corresponding CFG node must be in astPosition.
		// To avoid having astPosition grow too large, we exclude AST nodes that are not traversed
		// further, as they usually are not put in the CFG.
		c.setAstPosition(n, c.astPositionCounter)
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
				c.setAstPosition(n, c.astPositionCounter)
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
	if forNode.IsForOf() || forNode.IsForAwaitOf() {
		c.connectToPossibleExceptionHandlerUnconditionally(forNode)
	} else {
		c.connectToPossibleExceptionHandler(forNode, collection)
	}
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
				c.putFinally(lastJump, ControlFlowAnalysis_ComputeFallThrough(cur.GetLastChild()))
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
		c.putFinally(lastJump, ControlFlowAnalysis_ComputeFollowNodeWithCfa(cur, c))
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
				c.putFinally(lastJump, ControlFlowAnalysis_ComputeFallThrough(cur.GetLastChild()))
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
		c.putFinally(lastJump, iter)
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
				c.putFinally(lastJump, ControlFlowAnalysis_ComputeFallThrough(curHandler.GetLastChild()))
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
		c.putFinally(lastJump, nil)
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

func cfaHasEdgeInDirection(
	source *LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch],
	branch ControlFlowGraph_Branch,
	dest *LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch],
) bool {
	outEdges := source.outEdgeList.Slice()
	inEdges := dest.inEdgeList.Slice()
	if len(outEdges) < len(inEdges) {
		for _, outEdge := range outEdges {
			if outEdge.GetDestination() == DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch](dest) && outEdge.GetValue() == branch {
				return true
			}
		}
	} else {
		for _, inEdge := range inEdges {
			if inEdge.GetSource() == DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch](source) && inEdge.GetValue() == branch {
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
func (c *ControlFlowAnalysis) createEdge(fromNode *rhino.Node, branch ControlFlowGraph_Branch, toNode *rhino.Node) {
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
func (c *ControlFlowAnalysis) connectToPossibleExceptionHandler(cfgNode *rhino.Node, target *rhino.Node) {
	if len(c.exceptionHandler) == 0 || c.exceptionHandler[len(c.exceptionHandler)-1].IsFunction() {
		return
	}
	if ControlFlowAnalysis_MayThrowException(target) {
		c.connectToPossibleExceptionHandlerUnconditionally(cfgNode)
	}
}

// Java: ControlFlowAnalysis.connectToPossibleExceptionHandlerUnconditionally(Node)
func (c *ControlFlowAnalysis) connectToPossibleExceptionHandlerUnconditionally(cfgNode *rhino.Node) {
	if len(c.exceptionHandler) == 0 {
		return
	}
	lastJump := cfgNode
	for i := len(c.exceptionHandler) - 1; i >= 0; i-- {
		handler := c.exceptionHandler[i]
		if handler.IsFunction() {
			return
		}
		jutil.CheckState(handler.IsTry())
		catchBlock := NodeUtil_GetCatchBlock(handler)

		lastJumpInCatchBlock := false
		for ancestor := lastJump; ancestor != nil; ancestor = ancestor.GetParent() {
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
				c.putFinally(lastJump, handler.GetLastChild())
			}
		} else { // Has a catch.
			if lastJump == cfgNode {
				c.createEdge(cfgNode, ControlFlowGraph_Branch_ON_EX, catchBlock)
				return
			} else {
				c.putFinally(lastJump, catchBlock)
			}
		}
		lastJump = handler
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
	pendingRootNodes jutil.LinkedHashMap[*rhino.Node, *LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch]]
	wasRootCfg       bool
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

// ReleaseRootNodes returns any pending root CFG node map to cfaRootNodesPool.
func (g *ControlFlowAnalysis_AstControlFlowGraph) ReleaseRootNodes() {
	if g.wasRootCfg {
		g.wasRootCfg = false
		g.nodes = jutil.LinkedHashMap[*rhino.Node, *LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch]]{}
		if g.pendingRootNodes.HasBacking() {
			g.pendingRootNodes.Clear()
			cfaRootNodesPool.Put(g.pendingRootNodes)
			g.pendingRootNodes = jutil.LinkedHashMap[*rhino.Node, *LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch]]{}
		}
	}
}

// Reset resets g for reuse on a new CFG root while retaining node/edge slabs.
func (g *ControlFlowAnalysis_AstControlFlowGraph) Reset(entry *rhino.Node, edgeAnnotations bool) {
	g.self = g
	g.ReleaseRootNodes()
	g.ResetLinkedDirectedGraph(g /* useNodeAnnotations= */, true, edgeAnnotations)
	if entry != nil && entry.IsRoot() {
		if pooled, ok := cfaRootNodesPool.Get().(jutil.LinkedHashMap[*rhino.Node, *LinkedDirectedGraph_LinkedDiGraphNode[*rhino.Node, ControlFlowGraph_Branch]]); ok && pooled.HasBacking() {
			pooled.Clear()
			g.nodes = pooled
		}
	}
	g.implicitReturn = g.CreateDiGraphNode(nil)
	g.entry = g.CreateDiGraphNode(entry)
}

func cfaForwardNodeComparator(a, b DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) int {
	return cmp.Compare(a.GetPriority(), b.GetPriority())
}

func cfaBackwardNodeComparator(a, b DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]) int {
	return cmp.Compare(b.GetPriority(), a.GetPriority())
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
		return cfaForwardNodeComparator
	}
	// Comparator.reversed()
	return cfaBackwardNodeComparator
}

var (
	_ NodeTraversal_Callback        = (*ControlFlowAnalysis)(nil)
	_ ControlFlowGraph[*rhino.Node] = (*ControlFlowAnalysis_AstControlFlowGraph)(nil)
)

// Port of com/google/javascript/jscomp/ControlFlowGraph.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ControlFlowGraph is the control flow graph.
//
// N is the instruction type of the control flow graph.
//
// ControlFlowGraph has a subclass (ControlFlowAnalysis.AstControlFlowGraph),
// so it is an interface (PORTING.md §4); ControlFlowGraphBase holds its
// fields and implementations, and is also what the constructor returns. The
// interface declares, besides ControlFlowGraph's own methods, the public
// methods inherited from LinkedDirectedGraph.
type ControlFlowGraph[N comparable] interface {
	DiGraph[N, ControlFlowGraph_Branch]
	GraphvizGraph

	// Java: LinkedDirectedGraph.newSubGraph()
	NewSubGraph() SubGraph[N, ControlFlowGraph_Branch]
	// Java: LinkedDirectedGraph.connect(DiGraphNode<N, E>, E, DiGraphNode<N, E>)
	ConnectNodes(src DiGraph_DiGraphNode[N, ControlFlowGraph_Branch], edgeValue ControlFlowGraph_Branch, dest DiGraph_DiGraphNode[N, ControlFlowGraph_Branch])
	// Java: LinkedDirectedGraph.connectIfNotConnectedInDirection(N, E, N)
	ConnectIfNotConnectedInDirection(srcValue N, edgeValue ControlFlowGraph_Branch, destValue N)
	// Java: LinkedDirectedGraph.getNodes() (Collection<LinkedDiGraphNode<N, E>>)
	GetLinkedDiGraphNodes() []*LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]
	// Java: LinkedDirectedGraph.getNode(N) (LinkedDiGraphNode<N, E>)
	GetLinkedDiGraphNode(nodeValue N) *LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]
	// Java: LinkedDirectedGraph.createNode(N) (LinkedDiGraphNode<N, E>)
	CreateLinkedDiGraphNode(nodeValue N) *LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]
	// Java: LinkedDirectedGraph.isConnectedInDirection(LinkedDiGraphNode<N, E>, Predicate<E>, LinkedDiGraphNode<N, E>)
	IsConnectedInDirectionOfNodes(
		source *LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], edgeFilter func(ControlFlowGraph_Branch) bool, dest *LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]) bool

	// Java: ControlFlowGraph.getImplicitReturn()
	GetImplicitReturn() DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]
	// Java: ControlFlowGraph.getEntry()
	GetEntry() DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]
	// Java: ControlFlowGraph.isImplicitReturn(DiGraphNode<N, Branch>)
	IsImplicitReturn(node DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]) bool
	// Java: ControlFlowGraph.getOptionalNodeComparator(boolean)
	GetOptionalNodeComparator(isForward bool) func(a, b DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]) int
	// Java: ControlFlowGraph.toString()
	String() string

	GetControlFlowGraphBase() *ControlFlowGraphBase[N]
}

// ControlFlowGraphBase holds ControlFlowGraph's fields and implementations.
type ControlFlowGraphBase[N comparable] struct {
	LinkedDirectedGraph[N, ControlFlowGraph_Branch]

	self ControlFlowGraph[N]

	// A special node marked by the node value key null to a singleton
	// "return" when control is transferred outside of the current control
	// flow graph.
	implicitReturn DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]

	entry DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]
}

// NewControlFlowGraph is the constructor.
//
// Go-only: language-adapter — the Java null node value (the implicit
// return's key) is the zero value of N (nil for the pointer types that N is
// instantiated with).
//
// Java: ControlFlowGraph(N, boolean, boolean)
func NewControlFlowGraph[N comparable](entry N, nodeAnnotations bool, edgeAnnotations bool) *ControlFlowGraphBase[N] {
	c := &ControlFlowGraphBase[N]{}
	c.InitControlFlowGraphBase(c, entry, nodeAnnotations, edgeAnnotations)
	return c
}

// InitControlFlowGraphBase runs the ControlFlowGraph(N, boolean, boolean)
// constructor on c (for subclasses, which embed ControlFlowGraphBase); self
// is the object being constructed.
//
// Java: ControlFlowGraph(N, boolean, boolean)
func (c *ControlFlowGraphBase[N]) InitControlFlowGraphBase(self ControlFlowGraph[N], entry N, nodeAnnotations bool, edgeAnnotations bool) {
	c.self = self
	c.InitLinkedDirectedGraph(self, nodeAnnotations, edgeAnnotations)
	var null N
	c.implicitReturn = self.CreateDiGraphNode(null)
	c.entry = self.CreateDiGraphNode(entry)
}

// GetControlFlowGraphBase gives access to the base through the interface.
func (c *ControlFlowGraphBase[N]) GetControlFlowGraphBase() *ControlFlowGraphBase[N] { return c }

// GetImplicitReturn gets the implicit return node.
//
// Returns the return node.
//
// Java: ControlFlowGraph.getImplicitReturn()
func (c *ControlFlowGraphBase[N]) GetImplicitReturn() DiGraph_DiGraphNode[N, ControlFlowGraph_Branch] {
	return c.implicitReturn
}

// GetEntry gets the entry point of the control flow graph. In general, this
// should be the beginning of the global script or beginning of a function.
//
// Returns the entry point.
//
// Java: ControlFlowGraph.getEntry()
func (c *ControlFlowGraphBase[N]) GetEntry() DiGraph_DiGraphNode[N, ControlFlowGraph_Branch] {
	return c.entry
}

// IsImplicitReturn checks whether node is the implicit return.
//
// node is the node; returns true if the node is the implicit return.
//
// Java: ControlFlowGraph.isImplicitReturn(DiGraphNode<N, Branch>)
func (c *ControlFlowGraphBase[N]) IsImplicitReturn(node DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]) bool {
	return node == c.implicitReturn
}

// GetOptionalNodeComparator gets a comparator for the nodes. The default
// implementation returns null. See
// ControlFlowGraph#getOptionalNodeComparator.
//
// isForward is whether the comparator sorts the nodes in the direction of
// the flow; returns a comparator or null (in particular, if not overridden).
//
// Java: ControlFlowGraph.getOptionalNodeComparator(boolean)
func (c *ControlFlowGraphBase[N]) GetOptionalNodeComparator(isForward bool) func(a, b DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]) int {
	return nil
}

// ControlFlowGraph_Branch is the edge object for the control flow graph.
type ControlFlowGraph_Branch int32

const (
	// Edge is taken if the condition is true.
	ControlFlowGraph_Branch_ON_TRUE ControlFlowGraph_Branch = iota
	// Edge is taken if the condition is false.
	ControlFlowGraph_Branch_ON_FALSE
	// Unconditional branch.
	ControlFlowGraph_Branch_UNCOND
	// Exception-handling code paths. Conflates two kind of control flow
	// passing:
	//   - An exception is thrown, and falls into a catch or finally block
	//   - During exception handling, a finally block finishes and control
	//     passes to the next finally block.
	// In theory, we need 2 different edge types. In practice, we can just
	// treat them as "the edges we can't really optimize".
	ControlFlowGraph_Branch_ON_EX
	// Possible folded-away template
	ControlFlowGraph_Branch_SYN_BLOCK
)

var controlFlowGraph_Branch_names = [...]string{"ON_TRUE", "ON_FALSE", "UNCOND", "ON_EX", "SYN_BLOCK"}

// ControlFlowGraph_Branch_values is Branch.values().
func ControlFlowGraph_Branch_values() []ControlFlowGraph_Branch {
	return []ControlFlowGraph_Branch{
		ControlFlowGraph_Branch_ON_TRUE,
		ControlFlowGraph_Branch_ON_FALSE,
		ControlFlowGraph_Branch_UNCOND,
		ControlFlowGraph_Branch_ON_EX,
		ControlFlowGraph_Branch_SYN_BLOCK,
	}
}

// String is Enum.name().
func (b ControlFlowGraph_Branch) String() string { return controlFlowGraph_Branch_names[b] }

// Ordinal is Enum.ordinal().
func (b ControlFlowGraph_Branch) Ordinal() int { return int(b) }

// CompareTo is Enum.compareTo(E): the ordinal difference.
func (b ControlFlowGraph_Branch) CompareTo(o ControlFlowGraph_Branch) int { return int(b) - int(o) }

// Java: ControlFlowGraph.Branch.isConditional()
func (b ControlFlowGraph_Branch) IsConditional() bool {
	return b == ControlFlowGraph_Branch_ON_TRUE || b == ControlFlowGraph_Branch_ON_FALSE
}

// ControlFlowGraph_AbstractCfgNodeTraversalCallback is an abstract callback
// to visit a control flow graph node without going into subtrees of the node
// that are also represented by other control flow graph nodes.
//
// For example, traversing an IF node as root will visit the two subtrees
// pointed by the ControlFlowGraph.Branch#ON_TRUE and
// ControlFlowGraph.Branch#ON_FALSE edges.
//
// Subclasses embed ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase,
// which provides the final shouldTraverse, and implement Visit.
type ControlFlowGraph_AbstractCfgNodeTraversalCallback interface {
	NodeTraversal_Callback

	GetControlFlowGraph_AbstractCfgNodeTraversalCallbackBase() *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
}

// ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase holds the base
// implementation of ControlFlowGraph_AbstractCfgNodeTraversalCallback.
type ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase struct{}

// GetControlFlowGraph_AbstractCfgNodeTraversalCallbackBase gives access to
// the base through the interface.
func (b *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) GetControlFlowGraph_AbstractCfgNodeTraversalCallbackBase() *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase {
	return b
}

// ShouldTraverse is final in Java.
//
// Java: ControlFlowGraph.AbstractCfgNodeTraversalCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if parent.IsNil() {
		return true
	}
	return !ControlFlowGraph_IsEnteringNewCfgNode(n)
}

// ControlFlowGraph_IsEnteringNewCfgNode returns true if n should be
// represented by a new CFG node in the control flow graph.
//
// Java: ControlFlowGraph.isEnteringNewCfgNode(Node)
func ControlFlowGraph_IsEnteringNewCfgNode(n rhino.Node) bool {
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_ROOT, rhino.Token_SCRIPT, rhino.Token_TRY, rhino.Token_SWITCH_BODY:
		return true
	case rhino.Token_FUNCTION:
		// A function node represents the start of a function where the name
		// bleeds into the local scope and parameters are assigned
		// to the formal argument names. The node includes the name of the
		// function and the PARAM_LIST since we assume the whole set up process
		// is atomic without change in control flow. The next change of
		// control is going into the function's body, represented by the second
		// child.
		return n != parent.GetSecondChild()
	case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_IF:
		// These control structures are represented by a node that holds the
		// condition. Each of them is a branch node based on its condition.
		return NodeUtil_GetConditionExpression(parent) != n
	case rhino.Token_FOR:
		// The FOR(;;) node differs from other control structures in that
		// it has an initialization and an increment statement. Those
		// two statements have corresponding CFG nodes to represent them.
		// The FOR node only represents the condition check for each iteration.
		// That way the following:
		// for(var x = 0; x < 10; x++) { } has a graph that is isomorphic to
		// var x = 0; while(x<10) {  x++; }
		return NodeUtil_GetConditionExpression(parent) != n
	case rhino.Token_FOR_IN:
		// TODO(user): Investigate how we should handle the case where
		// we have a very complex expression inside the FOR-IN header.
		return n != parent.GetFirstChild()
	case rhino.Token_CASE, rhino.Token_CATCH, rhino.Token_WITH:
		return n != parent.GetFirstChild() // exclude the condition
	default:
		return false
	}
}

// Java: ControlFlowGraph.toString()
func (c *ControlFlowGraphBase[N]) String() string {
	s := "CFG:\n"
	for _, e := range c.self.GetGraphvizEdges() {
		s += jutil.ToJavaString(e) + "\n"
	}
	return s
}

// traverseCfgNode walks a single CFG node's subtree with a scope-free
// ControlFlowGraph_AbstractCfgNodeTraversalCallback without allocating a
// NodeTraversal or walking parent pointers to the script root.
func traverseCfgNode[T ControlFlowGraph_AbstractCfgNodeTraversalCallback](compiler AbstractCompiler, root rhino.Node, cb T) {
	jutil.CheckNotNull(compiler)
	if root.IsNil() {
		b := NodeTraversal_Builder{compiler: compiler, callback: cb}
		newNodeTraversal(&b).Traverse(root)
		return
	}
	var currentNode rhino.Node
	jutil.Try(func() {
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		currentNode = root
		walkCfgNodeBranch(root, rhino.Node{}, cb, &currentNode)
	}, func(unexpectedException jutil.Throwable) {
		b := NodeTraversal_Builder{compiler: compiler, callback: cb}
		t := newNodeTraversal(&b)
		t.InitTraversal(root)
		t.currentNode = currentNode
		t.ThrowUnexpectedException(unexpectedException)
	})
}

func walkCfgNodeBranch[T ControlFlowGraph_AbstractCfgNodeTraversalCallback](n, parent rhino.Node, cb T, currentNode *rhino.Node) {
	*currentNode = n
	if !parent.IsNil() && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return
	}
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			walkCfgNodeBranch(child, n, cb, currentNode)
			child = next
		}
	case rhino.Token_FUNCTION:
		fnName := n.GetFirstChild()
		if !parent.IsNil() {
			_ = NodeUtil_IsFunctionDeclaration(n)
		}
		args := fnName.GetNext()
		body := args.GetNext()
		walkCfgNodeBranch(fnName, n, cb, currentNode)
		walkCfgNodeBranch(args, n, cb, currentNode)
		walkCfgNodeBranch(body, n, cb, currentNode)
	case rhino.Token_CLASS:
		className := n.GetFirstChild()
		extendsClause := className.GetNext()
		body := extendsClause.GetNext()
		_ = NodeUtil_IsClassExpression(n)
		walkCfgNodeBranch(extendsClause, n, cb, currentNode)
		walkCfgNodeBranch(className, n, cb, currentNode)
		walkCfgNodeBranch(body, n, cb, currentNode)
	case rhino.Token_CLASS_MEMBERS:
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			if child.IsComputedProp() || child.IsComputedFieldDef() {
				walkCfgNodeBranch(child.GetFirstChild(), child, cb, currentNode)
			}
			child = next
		}
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			switch child.GetToken() {
			case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
				*currentNode = n
				if !ControlFlowGraph_IsEnteringNewCfgNode(child) {
					if child.HasTwoChildren() {
						walkCfgNodeBranch(child.GetLastChild(), child, cb, currentNode)
					}
					*currentNode = n
					cb.Visit(nil, child, n)
				}
			case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_BLOCK, rhino.Token_MEMBER_FUNCTION_DEF,
				rhino.Token_MEMBER_VARIABLE_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
				walkCfgNodeBranch(child, n, cb, currentNode)
			default:
				panic(jutil.NewIllegalStateExceptionMsg("Invalid class member: " + child.GetToken().String()))
			}
			child = next
		}
	default:
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			walkCfgNodeBranch(child, n, cb, currentNode)
			child = next
		}
	}
	*currentNode = n
	cb.Visit(nil, n, parent)
}

var (
	_ ControlFlowGraph[rhino.Node] = (*ControlFlowGraphBase[rhino.Node])(nil)
)

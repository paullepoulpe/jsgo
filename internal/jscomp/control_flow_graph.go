// Port of com/google/javascript/jscomp/ControlFlowGraph.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ControlFlowGraph is the control flow graph.
//
// N is the instruction type of the control flow graph.
type ControlFlowGraph[N comparable] struct {
	graph.LinkedDirectedGraph[N, ControlFlowGraph_Branch]

	// A special node marked by the node value key null to a singleton
	// "return" when control is transferred outside of the current control
	// flow graph.
	implicitReturn *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]

	entry *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]

	hasNodePriorities bool
}

// NewControlFlowGraph is the constructor.
//
// Go-only: language-adapter — the Java null node value (the implicit
// return's key) is the zero value of N (nil for the pointer types that N is
// instantiated with).
//
// Java: ControlFlowGraph(N, boolean, boolean)
func NewControlFlowGraph[N comparable](entry N, nodeAnnotations bool, edgeAnnotations bool) *ControlFlowGraph[N] {
	c := &ControlFlowGraph[N]{}
	c.init(entry, nodeAnnotations, edgeAnnotations, false)
	return c
}

// newAstControlFlowGraph constructs a ControlFlowGraph which provides a node
// comparator based on the pre-order traversal of the AST.
//
// Java: ControlFlowAnalysis.AstControlFlowGraph(Node, boolean)
func newAstControlFlowGraph(entry rhino.Node, edgeAnnotations bool) *ControlFlowGraph[rhino.Node] {
	c := &ControlFlowGraph[rhino.Node]{}
	c.init(entry /* nodeAnnotations= */, true, edgeAnnotations, true)
	return c
}

// reinitAst resets c as a new AST ControlFlowGraph(entry, edgeAnnotations),
// reusing its underlying slab and node-map storage.
func (c *ControlFlowGraph[N]) reinitAst(entry N, edgeAnnotations bool) {
	c.ClearForReuse()
	c.init(entry /* nodeAnnotations= */, true, edgeAnnotations, true)
}

func (c *ControlFlowGraph[N]) init(entry N, nodeAnnotations bool, edgeAnnotations bool, hasNodePriorities bool) {
	c.Init(nodeAnnotations, edgeAnnotations)
	c.hasNodePriorities = hasNodePriorities
	var null N
	c.implicitReturn = c.CreateLinkedDiGraphNode(null)
	c.entry = c.CreateLinkedDiGraphNode(entry)
}

// GetImplicitReturn gets the implicit return node.
//
// Returns the return node.
//
// Java: ControlFlowGraph.getImplicitReturn()
func (c *ControlFlowGraph[N]) GetImplicitReturn() *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch] {
	return c.implicitReturn
}

// GetEntry gets the entry point of the control flow graph. In general, this
// should be the beginning of the global script or beginning of a function.
//
// Returns the entry point.
//
// Java: ControlFlowGraph.getEntry()
func (c *ControlFlowGraph[N]) GetEntry() *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch] {
	return c.entry
}

// IsImplicitReturn checks whether node is the implicit return.
//
// node is the node; returns true if the node is the implicit return.
//
// Java: ControlFlowGraph.isImplicitReturn(DiGraphNode<N, Branch>)
func (c *ControlFlowGraph[N]) IsImplicitReturn(node *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]) bool {
	return node == c.implicitReturn
}

func cfgForwardNodeComparator[N comparable](a, b *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]) int {
	return cmp.Compare(a.GetPriority(), b.GetPriority())
}

func cfgBackwardNodeComparator[N comparable](a, b *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]) int {
	return cmp.Compare(b.GetPriority(), a.GetPriority())
}

// GetOptionalNodeComparator gets a comparator for the nodes. The default
// implementation returns nil unless node priorities were computed by
// ControlFlowAnalysis.
//
// isForward is whether the comparator sorts the nodes in the direction of
// the flow; returns a comparator or nil.
//
// Java: ControlFlowGraph.getOptionalNodeComparator(boolean)
func (c *ControlFlowGraph[N]) GetOptionalNodeComparator(isForward bool) func(a, b *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]) int {
	if !c.hasNodePriorities {
		return nil
	}
	if isForward {
		return cfgForwardNodeComparator[N]
	}
	return cfgBackwardNodeComparator[N]
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
		// Go-only: note that parent.GetSecondChild() in the AST is actually
		// Params() (child index 1), not Body() (child index 2); only Params()
		// returns false here, matching Java's `n != parent.getSecondChild()`.
		return !ast.MustFunction(parent).Params().Is(n)
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
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// TODO(user): Investigate how we should handle the case where
		// we have a very complex expression inside the FOR-IN header.
		return !ast.MustEnhancedFor(parent).Target().Is(n)
	case rhino.Token_CASE:
		return !ast.MustCase(parent).Value().Is(n)
	case rhino.Token_CATCH:
		return parent.GetFirstChild() != n // exclude catch param (or its EMPTY placeholder)
	case rhino.Token_WITH:
		return !ast.MustWith(parent).Object().Is(n)
	default:
		return false
	}
}

// Java: ControlFlowGraph.toString()
func (c *ControlFlowGraph[N]) String() string {
	s := "CFG:\n"
	for _, e := range c.GetGraphvizEdges() {
		s += jutil.ToJavaString(e) + "\n"
	}
	return s
}

func throwNilCfgNodeTraversal(compiler AbstractCompiler) {
	jutil.CheckNotNull(compiler)
	if Platform_IsThreadInterrupted() {
		exc := jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException())
		compiler.ThrowInternalError(jutil.ToJavaString(jutil.GetMessageOrNull(exc)), exc)
		return
	}
	compiler.ThrowInternalError("null", jutil.NewNullPointerException())
}

var (
	_ graph.DiGraph[rhino.Node, ControlFlowGraph_Branch] = (*ControlFlowGraph[rhino.Node])(nil)
	_ graph.GraphvizGraph                                = (*ControlFlowGraph[rhino.Node])(nil)
)

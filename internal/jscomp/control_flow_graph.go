// Port of com/google/javascript/jscomp/ControlFlowGraph.java @ closure-compiler v20261104

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
	"strings"

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

	hasNodePriorities  bool
	hasInnerScopeRoots bool
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
	c.hasInnerScopeRoots = false
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
	case rhino.Token_WHILE:
		return !ast.MustWhile(parent).Cond().Is(n)
	case rhino.Token_DO:
		return !ast.MustDo(parent).Cond().Is(n)
	case rhino.Token_IF:
		return !ast.MustIf(parent).Cond().Is(n)
	case rhino.Token_FOR:
		// The FOR(;;) node differs from other control structures in that
		// it has an initialization and an increment statement. Those
		// two statements have corresponding CFG nodes to represent them.
		// The FOR node only represents the condition check for each iteration.
		// That way the following:
		// for(var x = 0; x < 10; x++) { } has a graph that is isomorphic to
		// var x = 0; while(x<10) {  x++; }
		return ast.RoleIn(n, parent) != ast.RoleForCond
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// TODO(user): Investigate how we should handle the case where
		// we have a very complex expression inside the FOR-IN header.
		return !ast.MustEnhancedFor(parent).Target().Is(n)
	case rhino.Token_CASE:
		return !ast.MustCase(parent).Value().Is(n)
	case rhino.Token_CATCH:
		return ast.RoleIn(n, parent) != ast.RoleCatchParam // exclude catch param (or its EMPTY placeholder)
	case rhino.Token_WITH:
		return !ast.MustWith(parent).Object().Is(n)
	default:
		return false
	}
}

// Java: ControlFlowGraph.toString()
func (c *ControlFlowGraph[N]) String() string {
	var sb strings.Builder
	sb.WriteString("CFG:\n")
	for _, e := range c.GetGraphvizEdges() {
		sb.WriteString(jutil.ToJavaString(e))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func throwNilCfgNodeTraversal(compiler AbstractCompiler) {
	jutil.CheckNotNull(compiler)
	if Platform_IsThreadInterrupted() {
		exc := jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException())
		compiler.ThrowInternalError(throwableMessageOrNullString(exc), exc)
		return
	}
	compiler.ThrowInternalError("null", jutil.NewNullPointerException())
}

// isCfgSubtreeBoundary reports whether n crosses out of the current CFG node's
// synchronous execution into a nested CFG node, nested function, instance field
// initializer, or class static block.
func isCfgSubtreeBoundary(n rhino.Node, isRoot bool, allowRootFunction bool) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		if !isRoot || !allowRootFunction {
			return true
		}
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_VARIABLE_DEF:
		return true
	case rhino.Token_MEMBER_FIELD_DEF:
		if !ast.MustMemberFieldDef(n).IsStatic() {
			return true
		}
	case rhino.Token_BLOCK:
		if !isRoot && n.GetParent().IsClassMembers() {
			return true
		}
	}
	return !isRoot && ControlFlowGraph_IsEnteringNewCfgNode(n)
}

func isCfgCondEvalStopToken(tok rhino.Token) bool {
	switch tok {
	case rhino.Token_SCRIPT,
		rhino.Token_ROOT,
		rhino.Token_FUNCTION,
		rhino.Token_BLOCK,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF,
		rhino.Token_MEMBER_VARIABLE_DEF,
		rhino.Token_INDEX_SIGNATURE,
		rhino.Token_CALL_SIGNATURE:
		return true
	default:
		return false
	}
}

func isCfgCondEvalStopNode(n rhino.Node) bool {
	if isCfgCondEvalStopToken(n.GetToken()) {
		return true
	}
	mf, ok := ast.AsMemberFieldDef(n)
	return ok && !mf.IsStatic()
}

// walkCfgNode walks the AST subtree of cfgRoot without crossing into nested CFG
// nodes, nested functions, or deferred class members, calling visit(n, parent)
// in pre-order and returning true immediately if visit returns true.
func walkCfgNode(cfgRoot rhino.Node, visit func(n, parent rhino.Node) bool) bool {
	if cfgRoot.IsNil() {
		return false
	}
	return walkCfgSubtree(cfgRoot, rhino.NilNode, true, visit)
}

func walkCfgSubtree(n, parent rhino.Node, isRoot bool, visit func(n, parent rhino.Node) bool) bool {
	if n.IsNil() || isCfgSubtreeBoundary(n, isRoot, false) {
		return false
	}
	if visit(n, parent) {
		return true
	}
	switch n.GetToken() {
	case rhino.Token_CLASS:
		cls := ast.MustClass(n)
		return walkCfgSubtree(cls.ExtendsNode(), n, false, visit) ||
			walkCfgSubtree(cls.NameNode(), n, false, visit) ||
			walkCfgSubtree(cls.Body().Node(), n, false, visit)
	case rhino.Token_CLASS_MEMBERS:
		body := ast.MustClassBody(n)
		for cm := range body.ComputedMembers() {
			if walkCfgSubtree(cm.Key().Node(), cm.Node(), false, visit) {
				return true
			}
		}
		for m, ok := body.FirstMember(); ok; m, ok = m.Next() {
			switch m.Token() {
			case rhino.Token_COMPUTED_FIELD_DEF:
				cf := m.MustComputedFieldDef()
				if cf.IsStatic() {
					if val, hasVal := cf.Value(); hasVal && walkCfgSubtree(val.Node(), cf.Node(), false, visit) {
						return true
					}
				}
			case rhino.Token_MEMBER_FIELD_DEF:
				if m.MustMemberFieldDef().IsStatic() && walkCfgSubtree(m.Node(), n, false, visit) {
					return true
				}
			}
		}
		return false
	case rhino.Token_COMPUTED_FIELD_DEF:
		cf := ast.MustComputedFieldDef(n)
		if walkCfgSubtree(cf.Key().Node(), n, false, visit) {
			return true
		}
		if cf.IsStatic() {
			if val, hasVal := cf.Value(); hasVal {
				return walkCfgSubtree(val.Node(), n, false, visit)
			}
		}
		return false
	default:
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if walkCfgSubtree(c, n, false, visit) {
				return true
			}
		}
		return false
	}
}

// walkCfgEvalOrder walks the AST subtree of cfgRoot without crossing into
// nested CFG nodes, nested functions, or deferred class members, calling
// visit(n, parent) in post-order according to JavaScript runtime evaluation order.
func walkCfgEvalOrder(cfgRoot rhino.Node, visit func(n, parent rhino.Node)) {
	if !cfgRoot.IsNil() {
		walkCfgEvalSubtree(cfgRoot, rhino.NilNode, true, visit)
	}
}

func walkCfgEvalSubtree(n, parent rhino.Node, isRoot bool, visit func(n, parent rhino.Node)) {
	if n.IsNil() || isCfgSubtreeBoundary(n, isRoot, true) {
		return
	}
	if isRoot {
		parent = n.GetParent()
	}
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		for c := n.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			walkCfgEvalSubtree(c, n, false, visit)
			c = next
		}
	case rhino.Token_DESTRUCTURING_LHS:
		dlhs := ast.MustDestructuringLhs(n)
		if init, hasInit := dlhs.Init(); hasInit {
			walkCfgEvalSubtree(init.Node(), n, false, visit)
		}
		walkCfgEvalSubtree(dlhs.Pattern().Node(), n, false, visit)
	case rhino.Token_ASSIGN:
		assign := ast.MustAssign(n)
		if pat, isPat := assign.Target().AsPattern(); isPat {
			walkCfgEvalSubtree(assign.Value().Node(), n, false, visit)
			walkCfgEvalSubtree(pat.Node(), n, false, visit)
		} else {
			walkCfgEvalSubtree(assign.Target().Node(), n, false, visit)
			walkCfgEvalSubtree(assign.Value().Node(), n, false, visit)
		}
	case rhino.Token_DEFAULT_VALUE:
		dv := ast.MustDefaultValue(n)
		walkCfgEvalSubtree(dv.Value().Node(), n, false, visit)
		walkCfgEvalSubtree(dv.Target().Node(), n, false, visit)
	case rhino.Token_CLASS:
		cls := ast.MustClass(n)
		walkCfgEvalSubtree(cls.ExtendsNode(), n, false, visit)
		walkCfgEvalSubtree(cls.NameNode(), n, false, visit)
		walkCfgEvalSubtree(cls.Body().Node(), n, false, visit)
	case rhino.Token_CLASS_MEMBERS:
		body := ast.MustClassBody(n)
		for cm := range body.ComputedMembers() {
			walkCfgEvalSubtree(cm.Key().Node(), cm.Node(), false, visit)
		}
		for m, ok := body.FirstMember(); ok; {
			next, nextOk := m.Next()
			switch m.Token() {
			case rhino.Token_COMPUTED_FIELD_DEF:
				cf := m.MustComputedFieldDef()
				if cf.IsStatic() {
					if val, hasVal := cf.Value(); hasVal {
						walkCfgEvalSubtree(val.Node(), cf.Node(), false, visit)
					}
				}
			case rhino.Token_MEMBER_FIELD_DEF:
				if m.MustMemberFieldDef().IsStatic() {
					walkCfgEvalSubtree(m.Node(), n, false, visit)
				}
			}
			m, ok = next, nextOk
		}
	case rhino.Token_COMPUTED_FIELD_DEF:
		cf := ast.MustComputedFieldDef(n)
		walkCfgEvalSubtree(cf.Key().Node(), n, false, visit)
		if cf.IsStatic() {
			if val, hasVal := cf.Value(); hasVal {
				walkCfgEvalSubtree(val.Node(), n, false, visit)
			}
		}
	default:
		for c := n.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			walkCfgEvalSubtree(c, n, false, visit)
			c = next
		}
	}
	visit(n, parent)
}

type cfgCondEvalState uint8

const (
	cfgCondEvalContinue cfgCondEvalState = iota
	cfgCondEvalRead
	cfgCondEvalKill
)

func joinCfgCondEvalBranches(v1, v2 cfgCondEvalState) cfgCondEvalState {
	if v1 == cfgCondEvalRead || v2 == cfgCondEvalRead {
		return cfgCondEvalRead
	}
	if v1 == cfgCondEvalKill && v2 == cfgCondEvalKill {
		return cfgCondEvalKill
	}
	return cfgCondEvalContinue
}

// walkCfgCondEvalForward walks the intra-CFG expression subtree rooted at n in
// forward JavaScript runtime evaluation order, propagating conditional=true
// into short-circuiting and optional-chaining branches and stopping early when
// a branch returns non-cfgCondEvalContinue.
func walkCfgCondEvalForward[T any](
	n rhino.Node,
	conditional bool,
	ctx T,
	visit func(ctx T, n rhino.Node, conditional bool) (cfgCondEvalState, bool),
) cfgCondEvalState {
	if n.IsNil() {
		return cfgCondEvalContinue
	}
	if state, handled := visit(ctx, n, conditional); handled {
		return state
	}
	if isCfgCondEvalStopNode(n) {
		return cfgCondEvalContinue
	}
	switch n.GetToken() {
	case rhino.Token_WHILE:
		return walkCfgCondEvalForward(ast.MustWhile(n).Cond().Node(), conditional, ctx, visit)
	case rhino.Token_DO:
		return walkCfgCondEvalForward(ast.MustDo(n).Cond().Node(), conditional, ctx, visit)
	case rhino.Token_IF:
		return walkCfgCondEvalForward(ast.MustIf(n).Cond().Node(), conditional, ctx, visit)
	case rhino.Token_FOR:
		if cond, ok := ast.MustFor(n).Cond(); ok {
			return walkCfgCondEvalForward(cond.Node(), conditional, ctx, visit)
		}
		return cfgCondEvalContinue
	case rhino.Token_OPTCHAIN_GETPROP:
		return walkCfgCondEvalForward(ast.MustOptChainGetProp(n).Object().Node(), conditional, ctx, visit)
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		bin := ast.MustBinary(n)
		if v1 := walkCfgCondEvalForward(bin.Left().Node(), conditional, ctx, visit); v1 != cfgCondEvalContinue {
			return v1
		}
		if walkCfgCondEvalForward(bin.Right().Node(), true, ctx, visit) == cfgCondEvalRead {
			return cfgCondEvalRead
		}
		return cfgCondEvalContinue
	case rhino.Token_OPTCHAIN_GETELEM:
		ge := ast.MustOptChainGetElem(n)
		if v1 := walkCfgCondEvalForward(ge.Object().Node(), conditional, ctx, visit); v1 != cfgCondEvalContinue {
			return v1
		}
		if walkCfgCondEvalForward(ge.Index().Node(), true, ctx, visit) == cfgCondEvalRead {
			return cfgCondEvalRead
		}
		return cfgCondEvalContinue
	case rhino.Token_OPTCHAIN_CALL:
		call := ast.MustOptChainCall(n)
		if v1 := walkCfgCondEvalForward(call.Callee().Node(), conditional, ctx, visit); v1 != cfgCondEvalContinue {
			return v1
		}
		for arg, ok := call.FirstArg(); ok; arg, ok = arg.Next() {
			if state := walkCfgCondEvalForward(arg.Node(), true, ctx, visit); state != cfgCondEvalContinue {
				if state == cfgCondEvalRead {
					return cfgCondEvalRead
				}
				return cfgCondEvalContinue
			}
		}
		return cfgCondEvalContinue
	case rhino.Token_HOOK:
		hook := ast.MustHook(n)
		if first := walkCfgCondEvalForward(hook.Cond().Node(), conditional, ctx, visit); first != cfgCondEvalContinue {
			return first
		}
		return joinCfgCondEvalBranches(
			walkCfgCondEvalForward(hook.Then().Node(), true, ctx, visit),
			walkCfgCondEvalForward(hook.Else().Node(), true, ctx, visit),
		)
	case rhino.Token_CLASS:
		cls := ast.MustClass(n)
		if state := walkCfgCondEvalForward(cls.NameNode(), conditional, ctx, visit); state != cfgCondEvalContinue {
			return state
		}
		if state := walkCfgCondEvalForward(cls.ExtendsNode(), conditional, ctx, visit); state != cfgCondEvalContinue {
			return state
		}
		return walkCfgCondEvalForward(cls.Body().Node(), conditional, ctx, visit)
	case rhino.Token_CLASS_MEMBERS:
		body := ast.MustClassBody(n)
		for cm := range body.ComputedMembers() {
			if state := walkCfgCondEvalForward(cm.Key().Node(), conditional, ctx, visit); state != cfgCondEvalContinue {
				return state
			}
		}
		for m, ok := body.FirstMember(); ok; m, ok = m.Next() {
			if fd, isField := m.AsFieldDef(); isField && fd.IsStatic() {
				if val, hasVal := fd.Value(); hasVal {
					if state := walkCfgCondEvalForward(val.Node(), conditional, ctx, visit); state != cfgCondEvalContinue {
						return state
					}
				}
			}
		}
		return cfgCondEvalContinue
	case rhino.Token_COMPUTED_FIELD_DEF:
		cf := ast.MustComputedFieldDef(n)
		if state := walkCfgCondEvalForward(cf.Key().Node(), conditional, ctx, visit); state != cfgCondEvalContinue {
			return state
		}
		if cf.IsStatic() {
			if val, hasVal := cf.Value(); hasVal {
				return walkCfgCondEvalForward(val.Node(), conditional, ctx, visit)
			}
		}
		return cfgCondEvalContinue
	default:
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if state := walkCfgCondEvalForward(c, conditional, ctx, visit); state != cfgCondEvalContinue {
				return state
			}
		}
		return cfgCondEvalContinue
	}
}

// walkCfgCondEvalBackward walks the intra-CFG expression subtree rooted at n in
// reverse JavaScript runtime evaluation order, propagating conditional=true
// into short-circuiting and optional-chaining branches.
func walkCfgCondEvalBackward[T any](
	n rhino.Node,
	conditional bool,
	ctx T,
	visit func(ctx T, n rhino.Node, conditional bool) bool,
) {
	if n.IsNil() || visit(ctx, n, conditional) {
		return
	}
	if isCfgCondEvalStopNode(n) {
		return
	}
	switch n.GetToken() {
	case rhino.Token_WHILE:
		walkCfgCondEvalBackward(ast.MustWhile(n).Cond().Node(), conditional, ctx, visit)
	case rhino.Token_DO:
		walkCfgCondEvalBackward(ast.MustDo(n).Cond().Node(), conditional, ctx, visit)
	case rhino.Token_IF:
		walkCfgCondEvalBackward(ast.MustIf(n).Cond().Node(), conditional, ctx, visit)
	case rhino.Token_FOR:
		if cond, ok := ast.MustFor(n).Cond(); ok {
			walkCfgCondEvalBackward(cond.Node(), conditional, ctx, visit)
		}
	case rhino.Token_OPTCHAIN_GETPROP:
		walkCfgCondEvalBackward(ast.MustOptChainGetProp(n).Object().Node(), conditional, ctx, visit)
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		bin := ast.MustBinary(n)
		walkCfgCondEvalBackward(bin.Right().Node(), true, ctx, visit)
		walkCfgCondEvalBackward(bin.Left().Node(), conditional, ctx, visit)
	case rhino.Token_OPTCHAIN_GETELEM:
		ge := ast.MustOptChainGetElem(n)
		walkCfgCondEvalBackward(ge.Index().Node(), true, ctx, visit)
		walkCfgCondEvalBackward(ge.Object().Node(), conditional, ctx, visit)
	case rhino.Token_OPTCHAIN_CALL:
		call := ast.MustOptChainCall(n)
		callee := call.Callee().Node()
		for c := n.GetLastChild(); c != callee; c = c.GetPrevious() {
			walkCfgCondEvalBackward(c, true, ctx, visit)
		}
		walkCfgCondEvalBackward(callee, conditional, ctx, visit)
	case rhino.Token_HOOK:
		hook := ast.MustHook(n)
		walkCfgCondEvalBackward(hook.Else().Node(), true, ctx, visit)
		walkCfgCondEvalBackward(hook.Then().Node(), true, ctx, visit)
		walkCfgCondEvalBackward(hook.Cond().Node(), conditional, ctx, visit)
	case rhino.Token_CLASS:
		cls := ast.MustClass(n)
		walkCfgCondEvalBackward(cls.Body().Node(), conditional, ctx, visit)
		walkCfgCondEvalBackward(cls.ExtendsNode(), conditional, ctx, visit)
		walkCfgCondEvalBackward(cls.NameNode(), conditional, ctx, visit)
	case rhino.Token_CLASS_MEMBERS:
		body := ast.MustClassBody(n)
		for m, ok := body.LastMember(); ok; m, ok = m.Previous() {
			if fd, isField := m.AsFieldDef(); isField && fd.IsStatic() {
				if val, hasVal := fd.Value(); hasVal {
					walkCfgCondEvalBackward(val.Node(), conditional, ctx, visit)
				}
			}
		}
		for m, ok := body.LastMember(); ok; m, ok = m.Previous() {
			if cm, isComp := m.AsComputedMember(); isComp {
				walkCfgCondEvalBackward(cm.Key().Node(), conditional, ctx, visit)
			}
		}
	case rhino.Token_COMPUTED_FIELD_DEF:
		cf := ast.MustComputedFieldDef(n)
		if cf.IsStatic() {
			if val, hasVal := cf.Value(); hasVal {
				walkCfgCondEvalBackward(val.Node(), conditional, ctx, visit)
			}
		}
		walkCfgCondEvalBackward(cf.Key().Node(), conditional, ctx, visit)
	default:
		for c := n.GetLastChild(); !c.IsNil(); c = c.GetPrevious() {
			walkCfgCondEvalBackward(c, conditional, ctx, visit)
		}
	}
}

// walkClassMembersBefore visits every intra-CFG expression inside classMembers
// that evaluates synchronously strictly before (member, inPass2), where
// inPass2 == false denotes member's Pass 1 computed key and inPass2 == true
// denotes member's Pass 2 static field initializer.
func walkClassMembersBefore[T any](
	classMembers rhino.Node,
	member rhino.Node,
	inPass2 bool,
	ctx T,
	visit func(ctx T, expr rhino.Node) bool,
) bool {
	body := ast.MustClassBody(classMembers)
	if !inPass2 {
		for m, ok := body.FirstMember(); ok && !m.Is(member); m, ok = m.Next() {
			if cm, isComp := m.AsComputedMember(); isComp {
				if key := cm.Key().Node(); !key.IsNil() && visit(ctx, key) {
					return true
				}
			}
		}
		return false
	}
	for cm := range body.ComputedMembers() {
		if key := cm.Key().Node(); !key.IsNil() && visit(ctx, key) {
			return true
		}
	}
	for m, ok := body.FirstMember(); ok && !m.Is(member); m, ok = m.Next() {
		if fd, isField := m.AsFieldDef(); isField && fd.IsStatic() {
			if val, hasVal := fd.Value(); hasVal && visit(ctx, val.Node()) {
				return true
			}
		}
	}
	return false
}

// walkClassMembersAfterState visits every intra-CFG expression inside
// classMembers that evaluates synchronously strictly after (member, inPass2),
// stopping early as soon as visit returns non-cfgCondEvalContinue.
func walkClassMembersAfterState[T any](
	classMembers rhino.Node,
	member rhino.Node,
	inPass2 bool,
	ctx T,
	visit func(ctx T, expr rhino.Node) cfgCondEvalState,
) cfgCondEvalState {
	cur := ast.MustClassMember(member)
	startStatic, hasStatic := cur.Next()
	if !inPass2 {
		for m, ok := startStatic, hasStatic; ok; m, ok = m.Next() {
			if cm, isComp := m.AsComputedMember(); isComp {
				if key := cm.Key().Node(); !key.IsNil() {
					if state := visit(ctx, key); state != cfgCondEvalContinue {
						return state
					}
				}
			}
		}
		startStatic, hasStatic = ast.MustClassBody(classMembers).FirstMember()
	}
	for m, ok := startStatic, hasStatic; ok; m, ok = m.Next() {
		if fd, isField := m.AsFieldDef(); isField && fd.IsStatic() {
			if val, hasVal := fd.Value(); hasVal {
				if state := visit(ctx, val.Node()); state != cfgCondEvalContinue {
					return state
				}
			}
		}
	}
	return cfgCondEvalContinue
}

var (
	_ graph.DiGraph[rhino.Node, ControlFlowGraph_Branch] = (*ControlFlowGraph[rhino.Node])(nil)
	_ graph.GraphvizGraph                                = (*ControlFlowGraph[rhino.Node])(nil)
)

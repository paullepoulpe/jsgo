// Port of com/google/javascript/jscomp/NodeIterators.java
// @ closure-compiler v20260415.
//
// Provides common iteration patterns over the AST. All iterators
// are forward, post-order traversals unless otherwise noted.
//
// FunctionlessLocalScope traverses the local scope skipping
// function nodes. LocalVarMotion builds on it to find nodes where
// a variable retains its original value (used by inlining passes).
//
// LocalVarMotion depends on AbstractCompiler.GetAstAnalyzer() and
// AstAnalyzer.MayHaveSideEffects / NodeTypeMayHaveSideEffects.
// The astAnalyzer interface is implemented by AstAnalyzer in ast_analyzer.go.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// astAnalyzer is a minimal subset of AstAnalyzer used by LocalVarMotion.
// Matches the method signatures called in NodeIterators.java.
type astAnalyzer interface {
	MayHaveSideEffects(n *rhino.Node) bool
	NodeTypeMayHaveSideEffects(n *rhino.Node) bool
}

// abstractCompiler is a minimal subset of AbstractCompiler used by
// LocalVarMotion. Matches getAstAnalyzer() called in NodeIterators.java.
type abstractCompiler interface {
	GetAstAnalyzer() astAnalyzer
}

// FunctionlessLocalScope traverses the local scope, skipping all
// function nodes. Mirrors Java's NodeIterators.FunctionlessLocalScope.
type FunctionlessLocalScope struct {
	// ancestors holds the chain from deepest to shallowest; the last
	// element is the "current" node (Java uses a Deque with addLast /
	// peekLast / removeLast — i.e. a stack whose top is the tail).
	ancestors []*rhino.Node
}

// NewFunctionlessLocalScope constructs a FunctionlessLocalScope starting
// at the given ancestors (deepest first). The start node itself is not
// exposed in the iteration.
//
// Mirrors the Java constructor:
//
//	FunctionlessLocalScope(Node ... ancestors)
func NewFunctionlessLocalScope(ancestors ...*rhino.Node) *FunctionlessLocalScope {
	if len(ancestors) == 0 {
		panic("FunctionlessLocalScope: at least one ancestor required")
	}
	s := &FunctionlessLocalScope{}
	for _, n := range ancestors {
		if n.Token() == rhino.FUNCTION {
			break
		}
		s.ancestors = append(s.ancestors, n)
	}
	return s
}

// HasNext reports whether there is a next node to iterate. Mirrors
// Java's Iterator<Node>.hasNext().
func (s *FunctionlessLocalScope) HasNext() bool {
	return !(len(s.ancestors) == 1 && s.ancestors[len(s.ancestors)-1].Next() == nil)
}

// Next advances the iterator and returns the next node in post-order.
// Mirrors Java's Iterator<Node>.next().
func (s *FunctionlessLocalScope) Next() *rhino.Node {
	current := s.ancestors[len(s.ancestors)-1]
	s.ancestors = s.ancestors[:len(s.ancestors)-1]

	if current.Next() == nil {
		current = s.ancestors[len(s.ancestors)-1]

		// If this is a function node, skip it.
		if current.Token() == rhino.FUNCTION {
			return s.Next()
		}
	} else {
		current = current.Next()
		s.ancestors = append(s.ancestors, current)

		// If this is a function node, skip it.
		if current.Token() == rhino.FUNCTION {
			return s.Next()
		}

		for current.FirstChild() != nil {
			current = current.FirstChild()
			s.ancestors = append(s.ancestors, current)

			// If this is a function node, skip it.
			if current.Token() == rhino.FUNCTION {
				return s.Next()
			}
		}
	}

	return current
}

// current returns the node most recently returned by Next(). Mirrors
// Java's protected Node current().
func (s *FunctionlessLocalScope) current() *rhino.Node {
	return s.ancestors[len(s.ancestors)-1]
}

// currentParent returns the parent of the node most recently returned
// by Next(), or nil when the current node is at the top of the tracked
// chain. Mirrors Java's protected @Nullable Node currentParent().
func (s *FunctionlessLocalScope) currentParent() *rhino.Node {
	if len(s.ancestors) >= 2 {
		return s.current().Parent()
	}
	return nil
}

// CurrentAncestors returns the ancestors of the current node, deepest
// first. Only exposed for testing purposes — mirrors Java's
// List<Node> currentAncestors().
func (s *FunctionlessLocalScope) CurrentAncestors() []*rhino.Node {
	result := make([]*rhino.Node, len(s.ancestors))
	for i, n := range s.ancestors {
		result[len(s.ancestors)-1-i] = n
	}
	return result
}

// LocalVarMotion iterates past the declaration of a variable, finding
// all nodes where the variable is guaranteed to retain its original
// value. Mirrors Java's NodeIterators.LocalVarMotion.
//
// Consider:
//
//	var X = 1;
//	var Y = 3; // X is still 1
//	if (Y) {
//	  // X is still 1
//	} else {
//	  X = 5;
//	}
//	// X may not be 1
//
// The iterator stops at the assignment X = 5.
type LocalVarMotion struct {
	compiler            abstractCompiler
	valueHasSideEffects bool
	iterator            *FunctionlessLocalScope
	varName             string
	lookAhead           *rhino.Node
}

// ForVar creates a LocalVarMotion for use with moving a value assigned
// at a variable declaration. Mirrors Java's LocalVarMotion.forVar.
func ForVar(compiler abstractCompiler, name, varNode, block *rhino.Node) *LocalVarMotion {
	if !isNameDeclarationToken(varNode.Token()) {
		panic("ForVar: var must be a name declaration (VAR/LET/CONST)")
	}
	if !isStatementToken(varNode.Token()) {
		panic("ForVar: var must be a statement")
	}
	// The FunctionlessLocalScope must start at "name" as this may be used
	// before the Normalize pass, and thus the VAR node may define multiple
	// names and the "name" node may have siblings. The actual assigned
	// value is skipped as it is a child of name.
	return newLocalVarMotion(compiler, name, NewFunctionlessLocalScope(name, varNode, block))
}

// ForAssign creates a LocalVarMotion for use with moving a value
// assigned as part of a simple assignment expression ("a = b;").
// Mirrors Java's LocalVarMotion.forAssign.
func ForAssign(compiler abstractCompiler, name, assign, expr, block *rhino.Node) *LocalVarMotion {
	if assign.Token() != rhino.ASSIGN {
		panic("ForAssign: assign must be an ASSIGN node")
	}
	if expr.Token() != rhino.EXPR_RESULT {
		panic("ForAssign: expr must be an EXPR_RESULT node")
	}
	// The FunctionlessLocalScope must start at "assign", to skip the value
	// assigned to "name" (which would be its sibling).
	return newLocalVarMotion(compiler, name, NewFunctionlessLocalScope(assign, expr, block))
}

// newLocalVarMotion is the internal constructor. Mirrors Java's private
// LocalVarMotion(AbstractCompiler, Node, FunctionlessLocalScope).
func newLocalVarMotion(compiler abstractCompiler, nameNode *rhino.Node, iterator *FunctionlessLocalScope) *LocalVarMotion {
	if nameNode.Token() != rhino.NAME {
		panic("newLocalVarMotion: nameNode must be a NAME node")
	}
	valueNode := getAssignedValue(nameNode)
	var valueHasSideEffects bool
	if valueNode != nil && compiler != nil {
		valueHasSideEffects = compiler.GetAstAnalyzer().MayHaveSideEffects(valueNode)
	}
	m := &LocalVarMotion{
		compiler:            compiler,
		varName:             nameNode.String_(),
		valueHasSideEffects: valueHasSideEffects,
		iterator:            iterator,
	}
	m.advanceLookAhead(true)
	return m
}

// HasNext reports whether there is a next node. Mirrors Java's hasNext().
func (m *LocalVarMotion) HasNext() bool {
	return m.lookAhead != nil
}

// Next returns the next node and advances the iterator.
// Mirrors Java's next().
func (m *LocalVarMotion) Next() *rhino.Node {
	next := m.lookAhead
	m.advanceLookAhead(false)
	return next
}

// advanceLookAhead computes the next valid look-ahead node, or sets
// lookAhead to nil if iteration must stop. Mirrors Java's
// advanceLookAhead(boolean atStart).
func (m *LocalVarMotion) advanceLookAhead(atStart bool) {
	if !atStart {
		if m.lookAhead == nil {
			return
		}

		// Don't advance past a reference to the variable we're inlining.
		curNode := m.iterator.current()
		if curNode.Token() == rhino.NAME && m.varName == curNode.String_() {
			m.lookAhead = nil
			return
		}
	}

	if !m.iterator.HasNext() {
		m.lookAhead = nil
		return
	}

	nextNode := m.iterator.Next()
	nextParent := m.iterator.currentParent()
	typ := nextNode.Token()

	if m.valueHasSideEffects {
		// Reject anything that might read state.
		readsState := false

		if // Any read of a different variable.
		(nextNode.Token() == rhino.NAME && m.varName != nextNode.String_()) ||
			// Any read of a property.
			nextNode.Token() == rhino.GETPROP || nextNode.Token() == rhino.GETELEM {

			// If this is a simple assign, we'll be ok.
			if nextParent == nil || !isNameDeclOrSimpleAssignLhs(nextNode, nextParent) {
				readsState = true
			}

		} else if nextNode.Token() == rhino.CALL || nextNode.Token() == rhino.NEW {
			// Calls and new-expressions may read state. Included for
			// completeness — see Java comment in LocalVarMotion.
			readsState = true
		}

		if readsState {
			m.lookAhead = nil
			return
		}
	}

	// Reject anything that might modify relevant state.
	var typeMaySideEffect bool
	if m.compiler != nil {
		typeMaySideEffect = m.compiler.GetAstAnalyzer().NodeTypeMayHaveSideEffects(nextNode)
	}
	if (typeMaySideEffect && typ != rhino.NAME) ||
		(typ == rhino.NAME && nextParent != nil && nextParent.Token() == rhino.CATCH) {
		m.lookAhead = nil
		return
	}

	m.lookAhead = nextNode
}

// isNameDeclarationToken reports whether t is VAR, LET, or CONST.
// Mirrors NodeUtil.isNameDeclaration used in LocalVarMotion.forVar.
func isNameDeclarationToken(t rhino.Token) bool {
	return t == rhino.VAR || t == rhino.LET || t == rhino.CONST
}

// isStatementToken reports whether t corresponds to a statement node.
// Narrow approximation of NodeUtil.isStatement used by forVar's
// precondition check. VAR / LET / CONST at block level are statements.
func isStatementToken(t rhino.Token) bool {
	switch t {
	case rhino.VAR, rhino.LET, rhino.CONST,
		rhino.EXPR_RESULT, rhino.RETURN, rhino.IF, rhino.BLOCK,
		rhino.WHILE, rhino.DO_WHILE, rhino.FOR, rhino.FOR_IN, rhino.FOR_OF,
		rhino.SWITCH, rhino.TRY, rhino.THROW, rhino.BREAK, rhino.CONTINUE,
		rhino.LABEL, rhino.FUNCTION, rhino.CLASS:
		return true
	}
	return false
}

// getAssignedValue returns the value node that is assigned to nameNode,
// or nil if there is none. Mirrors NodeUtil.getAssignedValue in Java:
// for a NAME child of VAR/LET/CONST the assigned value is the NAME's
// first child; for a NAME on the LHS of an ASSIGN it is the RHS.
func getAssignedValue(nameNode *rhino.Node) *rhino.Node {
	parent := nameNode.Parent()
	if parent == nil {
		return nil
	}
	switch parent.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		return nameNode.FirstChild()
	case rhino.ASSIGN:
		if parent.FirstChild() == nameNode {
			kids := parent.Children()
			if len(kids) >= 2 {
				return kids[1]
			}
		}
	}
	return nil
}

// isNameDeclOrSimpleAssignLhs reports whether n is the left-hand side
// of a name declaration or a simple assignment. Mirrors Java's
// NodeUtil.isNameDeclOrSimpleAssignLhs.
func isNameDeclOrSimpleAssignLhs(n, parent *rhino.Node) bool {
	if parent.Token() == rhino.ASSIGN && parent.FirstChild() == n {
		return true
	}
	if isNameDeclarationToken(parent.Token()) && parent.FirstChild() == n {
		return true
	}
	return false
}

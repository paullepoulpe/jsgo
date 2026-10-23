// Port of com/google/javascript/jscomp/NodeIterators.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NodeIterators is a package for common iteration patterns.
//
// All iterators are forward, post-order traversals unless otherwise noted.
//
// All static (Java: private NodeIterators()).
type NodeIterators struct{}

// NodeIterators_FunctionlessLocalScope traverses the local scope, skipping
// all function nodes.
type NodeIterators_FunctionlessLocalScope struct {
	// Go-only: language-adapter — the ArrayDeque, first element at index 0
	// and last element at the end of the slice.
	ancestors    []*rhino.Node
	ancestorsBuf [16]*rhino.Node
}

// NewNodeIterators_FunctionlessLocalScope: ancestors are the ancestors of
// the point where iteration will start, beginning with the deepest
// ancestor. The start node will not be exposed in the iteration.
//
// Java: NodeIterators.FunctionlessLocalScope(Node...)
func NewNodeIterators_FunctionlessLocalScope(ancestors ...*rhino.Node) *NodeIterators_FunctionlessLocalScope {
	jutil.CheckArgument(len(ancestors) > 0)
	s := &NodeIterators_FunctionlessLocalScope{}
	s.ancestors = s.ancestorsBuf[:0]

	for _, n := range ancestors {
		if n.IsFunction() {
			break
		}

		s.addFirst(n)
	}
	return s
}

func (s *NodeIterators_FunctionlessLocalScope) init3(a, b, c *rhino.Node) {
	s.ancestors = s.ancestorsBuf[:0]
	if a.IsFunction() {
		return
	}
	jutil.CheckNotNull(a)
	if b.IsFunction() {
		s.ancestors = append(s.ancestors, a)
		return
	}
	jutil.CheckNotNull(b)
	if c.IsFunction() {
		s.ancestors = append(s.ancestors, b, a)
		return
	}
	jutil.CheckNotNull(c)
	s.ancestors = append(s.ancestors, c, b, a)
}

// Go-only: language-adapter — ArrayDeque.addFirst(E) (null-hostile).
func (s *NodeIterators_FunctionlessLocalScope) addFirst(n *rhino.Node) {
	jutil.CheckNotNull(n)
	s.ancestors = append(s.ancestors, nil)
	copy(s.ancestors[1:], s.ancestors)
	s.ancestors[0] = n
}

// Go-only: language-adapter — ArrayDeque.addLast(E) (null-hostile).
func (s *NodeIterators_FunctionlessLocalScope) addLast(n *rhino.Node) {
	jutil.CheckNotNull(n)
	s.ancestors = append(s.ancestors, n)
}

// Go-only: language-adapter — ArrayDeque.peekLast(): null when empty.
func (s *NodeIterators_FunctionlessLocalScope) peekLast() *rhino.Node {
	if len(s.ancestors) == 0 {
		return nil
	}
	return s.ancestors[len(s.ancestors)-1]
}

// Go-only: language-adapter — ArrayDeque.removeLast(): throws NoSuchElementException
// when empty.
func (s *NodeIterators_FunctionlessLocalScope) removeLast() *rhino.Node {
	if len(s.ancestors) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	last := s.ancestors[len(s.ancestors)-1]
	s.ancestors[len(s.ancestors)-1] = nil
	s.ancestors = s.ancestors[:len(s.ancestors)-1]
	return last
}

// Java: NodeIterators.FunctionlessLocalScope.hasNext()
func (s *NodeIterators_FunctionlessLocalScope) HasNext() bool {
	// Check if the current node has any nodes after it.
	return !(len(s.ancestors) == 1 && s.peekLast().GetNext() == nil)
}

// Java: NodeIterators.FunctionlessLocalScope.next()
func (s *NodeIterators_FunctionlessLocalScope) Next() *rhino.Node {
	current := s.removeLast()
	if current.GetNext() == nil {
		current = s.peekLast()

		// If this is a function node, skip it.
		if current.IsFunction() {
			return s.Next()
		}
	} else {
		current = current.GetNext()
		s.addLast(current)

		// If this is a function node, skip it.
		if current.IsFunction() {
			return s.Next()
		}

		for current.HasChildren() {
			current = current.GetFirstChild()
			s.addLast(current)

			// If this is a function node, skip it.
			if current.IsFunction() {
				return s.Next()
			}
		}
	}

	return current
}

// Java: NodeIterators.FunctionlessLocalScope.remove()
func (s *NodeIterators_FunctionlessLocalScope) Remove() {
	panic(jutil.NewUnsupportedOperationExceptionMsg("Not implemented"))
}

// Current gets the node most recently returned by next().
//
// Java: NodeIterators.FunctionlessLocalScope.current()
func (s *NodeIterators_FunctionlessLocalScope) Current() *rhino.Node {
	return s.peekLast()
}

// CurrentParent gets the parent of the node most recently returned by
// next().
//
// Java: NodeIterators.FunctionlessLocalScope.currentParent()
func (s *NodeIterators_FunctionlessLocalScope) CurrentParent() *rhino.Node {
	if len(s.ancestors) >= 2 {
		return s.Current().GetParent()
	}
	return nil
}

// CurrentAncestors gets the ancestors of the current node, with the deepest
// node first. Only exposed for testing purposes.
//
// Java: NodeIterators.FunctionlessLocalScope.currentAncestors()
func (s *NodeIterators_FunctionlessLocalScope) CurrentAncestors() *jutil.ArrayList[*rhino.Node] {
	list := jutil.NewArrayListFrom(append([]*rhino.Node(nil), s.ancestors...))
	// Collections.reverse(list)
	elems := list.Slice()
	for i, j := 0, len(elems)-1; i < j; i, j = i+1, j-1 {
		elems[i], elems[j] = elems[j], elems[i]
	}
	return list
}

// NodeIterators_LocalVarMotion is an iterator to help with variable
// inlining. Given a variable declaration, find all the nodes in post-order
// where the variable is guaranteed to retain its original value.
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
// In the above example, the iterator will iterate past the declaration of Y
// and into the first block of the IF branch, and will stop at the
// assignment X = 5.
type NodeIterators_LocalVarMotion struct {
	compiler            AbstractCompiler
	valueHasSideEffects bool
	iterator            *NodeIterators_FunctionlessLocalScope
	iteratorStorage     NodeIterators_FunctionlessLocalScope
	varName             string
	lookAhead           *rhino.Node
}

// NodeIterators_LocalVarMotion_ForVar: the name is a bit of a misnomer;
// this works with let and const as well.
//
// Returns a LocalVarMotion for use with moving a value assigned at a
// variable declaration.
//
// Java: NodeIterators.LocalVarMotion.forVar(AbstractCompiler, Node, Node, Node)
func NodeIterators_LocalVarMotion_ForVar(compiler AbstractCompiler, name *rhino.Node, var_ *rhino.Node, block *rhino.Node) *NodeIterators_LocalVarMotion {
	m := &NodeIterators_LocalVarMotion{}
	m.initForVar(compiler, name, var_, block)
	return m
}

func (m *NodeIterators_LocalVarMotion) initForVar(compiler AbstractCompiler, name *rhino.Node, var_ *rhino.Node, block *rhino.Node) {
	jutil.CheckArgument(NodeUtil_IsNameDeclaration(var_))
	jutil.CheckArgument(NodeUtil_IsStatement(var_))
	// The FunctionlessLocalScope must start at "name" as this may be used
	// before the Normalize pass, and thus the VAR node may define multiple
	// names and the "name" node may have siblings. The actual assigned value
	// is skipped as it is a child of name.
	m.iteratorStorage.init3(name, var_, block)
	m.init(compiler, name, &m.iteratorStorage)
}

// NodeIterators_LocalVarMotion_ForAssign returns a LocalVarMotion for use
// with moving a value assigned as part of a simple assignment expression
// ("a = b;").
//
// Java: NodeIterators.LocalVarMotion.forAssign(AbstractCompiler, Node, Node, Node, Node)
func NodeIterators_LocalVarMotion_ForAssign(
	compiler AbstractCompiler, name *rhino.Node, assign *rhino.Node, expr *rhino.Node, block *rhino.Node) *NodeIterators_LocalVarMotion {
	m := &NodeIterators_LocalVarMotion{}
	m.initForAssign(compiler, name, assign, expr, block)
	return m
}

func (m *NodeIterators_LocalVarMotion) initForAssign(
	compiler AbstractCompiler, name *rhino.Node, assign *rhino.Node, expr *rhino.Node, block *rhino.Node) {
	jutil.CheckArgument(assign.IsAssign())
	jutil.CheckArgument(expr.IsExprResult())
	// The FunctionlessLocalScope must start at "assign", to skip the value
	// assigned to "name" (which would be its sibling).
	m.iteratorStorage.init3(assign, expr, block)
	m.init(compiler, name, &m.iteratorStorage)
}

// newNodeIterators_LocalVarMotion: iterator is the iterator to use while
// inspecting the node.
//
// Java: NodeIterators.LocalVarMotion(AbstractCompiler, Node, FunctionlessLocalScope)
func newNodeIterators_LocalVarMotion(
	compiler AbstractCompiler, nameNode *rhino.Node, iterator *NodeIterators_FunctionlessLocalScope) *NodeIterators_LocalVarMotion {
	m := &NodeIterators_LocalVarMotion{}
	m.init(compiler, nameNode, iterator)
	return m
}

func (m *NodeIterators_LocalVarMotion) init(
	compiler AbstractCompiler, nameNode *rhino.Node, iterator *NodeIterators_FunctionlessLocalScope) {
	jutil.CheckArgument(nameNode.IsName())
	valueNode := NodeUtil_GetAssignedValue(nameNode)
	m.compiler = jutil.CheckNotNull(compiler)
	m.varName = nameNode.GetString()
	m.valueHasSideEffects = valueNode != nil && compiler.GetAstAnalyzer().MayHaveSideEffects(valueNode)
	m.iterator = iterator
	m.lookAhead = nil
	m.AdvanceLookAhead(true)
}

// Java: NodeIterators.LocalVarMotion.hasNext()
func (m *NodeIterators_LocalVarMotion) HasNext() bool {
	return m.lookAhead != nil
}

// Java: NodeIterators.LocalVarMotion.next()
func (m *NodeIterators_LocalVarMotion) Next() *rhino.Node {
	next := m.lookAhead
	m.AdvanceLookAhead(false)
	return next
}

// Java: NodeIterators.LocalVarMotion.remove()
func (m *NodeIterators_LocalVarMotion) Remove() {
	panic(jutil.NewUnsupportedOperationExceptionMsg("Not implemented"))
}

// Java: NodeIterators.LocalVarMotion.advanceLookAhead(boolean)
func (m *NodeIterators_LocalVarMotion) AdvanceLookAhead(atStart bool) {
	if !atStart {
		if m.lookAhead == nil {
			return
		}

		// Don't advance past a reference to the variable that we're trying
		// to inline.
		curNode := m.iterator.Current()
		if curNode.IsName() && m.varName == curNode.GetString() {
			m.lookAhead = nil
			return
		}
	}

	if !m.iterator.HasNext() {
		m.lookAhead = nil
		return
	}

	nextNode := m.iterator.Next()
	nextParent := m.iterator.CurrentParent()
	type_ := nextNode.GetToken()

	if m.valueHasSideEffects {
		// Reject anything that might read state
		readsState := false

		if // Any read of a different variable.
		(nextNode.IsName() && m.varName != nextNode.GetString()) ||
			// Any read of a property.
			(nextNode.IsGetProp() ||
				nextNode.IsGetElem() ||
				nextNode.IsOptChainGetProp() ||
				nextNode.IsOptChainGetElem()) {

			// If this is a simple assign, we'll be ok.
			if nextParent == nil || !NodeUtil_IsNameDeclOrSimpleAssignLhs(nextNode, nextParent) {
				readsState = true
			}

		} else if nextNode.IsCall() || nextNode.IsNew() || nextNode.IsOptChainCall() {
			// This isn't really an important case. In most cases when we use
			// CALL or NEW, we're invoking it on a NAME or a GETPROP. And in the
			// few cases where we're not, it's because we have an anonymous
			// function that escapes the variable we're worried about. But we
			// include this for completeness.
			readsState = true
		}

		if readsState {
			m.lookAhead = nil
			return
		}
	}

	// Reject anything that might modify relevant state. We assume that
	// nobody relies on variables being undeclared, which will break
	// constructions like:
	//   var a = b;
	//   var b = 3;
	//   alert(a);
	if (m.compiler.GetAstAnalyzer().NodeTypeMayHaveSideEffects(nextNode) && type_ != rhino.Token_NAME) ||
		(type_ == rhino.Token_NAME && nextParent.IsCatch()) {
		m.lookAhead = nil
		return
	}

	m.lookAhead = nextNode
}

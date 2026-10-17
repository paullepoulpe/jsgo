// Port of com/google/javascript/jscomp/ClosureCodeRemoval.java @ closure-compiler v20260415

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

// ClosureCodeRemoval is a compiler pass that removes Closure-specific code
// patterns.
//
// Currently does the following:
//
//   - Instead of setting abstract methods to a function that throws an
//     informative error, this pass allows some binary size reduction by
//     removing these methods altogether for production builds.
//   - Remove calls to assertion functions (like goog.asserts.assert). If the
//     return value of the assertion function is used, then the first argument
//     (the asserted value) will be directly inlined. Otherwise, the entire
//     call will be removed. It is well-known that this is not provably safe,
//     much like the equivalent assert statement in Java.
type ClosureCodeRemoval struct {
	// Reference to the JS compiler
	compiler AbstractCompiler

	removeAbstractMethods bool
	removeAssertionCalls  bool

	// List of names referenced in successive generations of finding
	// referenced nodes.
	abstractMethodAssignmentNodes *jutil.ArrayList[*ClosureCodeRemoval_RemovableAssignment]

	// List of member function definition nodes annotated with @abstract.
	abstractMemberFunctionNodes *jutil.ArrayList[*rhino.Node]

	// List of assertion functions.
	assertionCalls *jutil.ArrayList[*rhino.Node]
}

var _ CompilerPass = (*ClosureCodeRemoval)(nil)

// ClosureCodeRemoval_ABSTRACT_METHOD_NAME is the name used to denote an
// abstract function.
const ClosureCodeRemoval_ABSTRACT_METHOD_NAME = "goog.abstractMethod"

// ClosureCodeRemoval_RemovableAssignment is a utility class to track a node
// and its parent (an inner class: outer is the enclosing
// ClosureCodeRemoval).
type ClosureCodeRemoval_RemovableAssignment struct {
	outer *ClosureCodeRemoval

	// The node
	node *rhino.Node

	// Its parent
	parent *rhino.Node

	// Full chain of ASSIGN ancestors
	assignAncestors *jutil.ArrayList[*rhino.Node]

	// The last ancestor
	lastAncestor *rhino.Node
}

// NewRemovableAssignment is the data structure for information about a
// removable assignment. nameNode is the LHS; assignNode the parent ASSIGN
// node; traversal gives access to further levels, assumed to start at 1.
//
// Java: ClosureCodeRemoval.RemovableAssignment(Node, Node, NodeTraversal)
func (c *ClosureCodeRemoval) NewRemovableAssignment(nameNode *rhino.Node, assignNode *rhino.Node, traversal *NodeTraversal) *ClosureCodeRemoval_RemovableAssignment {
	a := &ClosureCodeRemoval_RemovableAssignment{outer: c, assignAncestors: jutil.NewArrayList[*rhino.Node]()}
	a.node = nameNode
	a.parent = assignNode

	ancestor := assignNode
	for {
		ancestor = ancestor.GetParent()
		a.assignAncestors.Add(ancestor)
		if !(ancestor.IsAssign() &&
			ancestor.GetFirstChild().IsQualifiedName()) {
			break
		}
	}
	a.lastAncestor = ancestor.GetParent()
	return a
}

// Remove removes this node.
//
// Java: ClosureCodeRemoval.RemovableAssignment.remove()
func (a *ClosureCodeRemoval_RemovableAssignment) Remove() {
	compiler := a.outer.compiler
	rhs := a.node.GetNext()
	last := a.parent
	for ancestor := range a.assignAncestors.All() {
		if ancestor.IsExprResult() {
			ancestor.Detach()
			NodeUtil_MarkFunctionsDeleted(ancestor, compiler)
		} else {
			rhs.Detach()
			last.ReplaceWith(rhs)
		}
		last = ancestor
	}
	compiler.ReportChangeToEnclosingScope(a.lastAncestor)
}

// ClosureCodeRemoval_FindAbstractMethods identifies all assignments of the
// abstract method to a variable and all methods annotated with "@abstract"
// in their JSDoc (an inner class).
type ClosureCodeRemoval_FindAbstractMethods struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *ClosureCodeRemoval
}

// Java: ClosureCodeRemoval.FindAbstractMethods.visit(NodeTraversal, Node, Node)
func (f *ClosureCodeRemoval_FindAbstractMethods) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	c := f.outer
	if n.IsAssign() {
		nameNode := n.GetFirstChild()
		valueNode := n.GetLastChild()

		if nameNode.IsQualifiedName() &&
			valueNode.IsQualifiedName() &&
			valueNode.MatchesQualifiedName(ClosureCodeRemoval_ABSTRACT_METHOD_NAME) {
			// Foo.prototype.bar = goog.abstractMethod
			c.abstractMethodAssignmentNodes.Add(
				c.NewRemovableAssignment(n.GetFirstChild(), n, t))
		} else if n.GetJSDocInfo() != nil &&
			n.GetJSDocInfo().IsAbstract() &&
			NodeUtil_IsEmptyFunctionExpression(valueNode) &&
			!n.GetJSDocInfo().IsConstructor() {
			// @abstract
			c.abstractMethodAssignmentNodes.Add(
				c.NewRemovableAssignment(n.GetFirstChild(), n, t))
		}
	} else if n.IsMemberFunctionDef() && parent.IsClassMembers() {
		if n.GetJSDocInfo() != nil && n.GetJSDocInfo().IsAbstract() {
			c.abstractMemberFunctionNodes.Add(n)
		}
	}
}

// ClosureCodeRemoval_FindAssertionCalls identifies all assertion calls (an
// inner class).
type ClosureCodeRemoval_FindAssertionCalls struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer          *ClosureCodeRemoval
	assertionNames *CodingConvention_AssertionFunctionLookup
}

// Java: ClosureCodeRemoval.FindAssertionCalls()
func (c *ClosureCodeRemoval) NewFindAssertionCalls() *ClosureCodeRemoval_FindAssertionCalls {
	f := &ClosureCodeRemoval_FindAssertionCalls{outer: c}
	f.assertionNames =
		CodingConvention_AssertionFunctionLookup_Of(c.compiler.GetCodingConvention().GetAssertionFunctions())
	return f
}

// Java: ClosureCodeRemoval.FindAssertionCalls.visit(NodeTraversal, Node, Node)
func (f *ClosureCodeRemoval_FindAssertionCalls) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !n.IsCall() {
		return
	}

	callee := n.GetFirstChild()

	if f.assertionNames.LookupByCallee(callee) != nil || // type-based
		(callee.GetColor() != nil && callee.GetColor().IsClosureAssert()) { // color-based
		f.outer.assertionCalls.Add(n)
	}
}

// NewClosureCodeRemoval creates a Closure code remover. compiler is the
// AbstractCompiler; removeAbstractMethods removes declarations of abstract
// methods; removeAssertionCalls removes calls to goog.assert functions.
//
// Java: ClosureCodeRemoval(AbstractCompiler, boolean, boolean)
func NewClosureCodeRemoval(compiler AbstractCompiler, removeAbstractMethods bool, removeAssertionCalls bool) *ClosureCodeRemoval {
	return &ClosureCodeRemoval{
		compiler:                      compiler,
		removeAbstractMethods:         removeAbstractMethods,
		removeAssertionCalls:          removeAssertionCalls,
		abstractMethodAssignmentNodes: jutil.NewArrayList[*ClosureCodeRemoval_RemovableAssignment](),
		abstractMemberFunctionNodes:   jutil.NewArrayList[*rhino.Node](),
		assertionCalls:                jutil.NewArrayList[*rhino.Node](),
	}
}

// Java: ClosureCodeRemoval.process(Node, Node)
func (c *ClosureCodeRemoval) Process(externs *rhino.Node, root *rhino.Node) {
	var passes []NodeTraversal_Callback
	if c.removeAbstractMethods {
		passes = append(passes, &ClosureCodeRemoval_FindAbstractMethods{outer: c})
	}
	if c.removeAssertionCalls {
		passes = append(passes, c.NewFindAssertionCalls())
	}
	CombinedCompilerPass_Traverse(c.compiler, root, passes)

	for assignment := range c.abstractMethodAssignmentNodes.All() {
		assignment.Remove()
	}

	for memberFunction := range c.abstractMemberFunctionNodes.All() {
		c.compiler.ReportFunctionDeleted(memberFunction.GetFirstChild())
		parent := memberFunction.GetParent()
		memberFunction.Detach()
		c.compiler.ReportChangeToEnclosingScope(parent)
	}

	for call := range c.assertionCalls.All() {
		// If the assertion is an expression, just strip the whole thing.
		c.compiler.ReportChangeToEnclosingScope(call)
		parent := call.GetParent()
		if parent.IsExprResult() {
			parent.Detach()
			NodeUtil_MarkFunctionsDeleted(parent, c.compiler)
		} else {
			// Otherwise, replace the assertion with its first argument,
			// which is the return value of the assertion.
			firstArg := call.GetSecondChild()
			if firstArg == nil {
				call.ReplaceWith(NodeUtil_NewUndefinedNode(call))
			} else {
				replacement := firstArg.Detach()
				replacement.CopyTypeFrom(call)
				call.ReplaceWith(replacement)
			}
			NodeUtil_MarkFunctionsDeleted(call, c.compiler)
		}
	}
}

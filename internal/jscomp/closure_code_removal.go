// Port of com/google/javascript/jscomp/ClosureCodeRemoval.java @ closure-compiler v20261028

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	ccrGoogID           = uint32(strpool.Intern("goog"))
	ccrAbstractMethodID = uint32(strpool.Intern("abstractMethod"))

	closureCodeRemovalAbstractVisitTokens = makeVisitTokenSet(
		rhino.Token_ASSIGN,
		rhino.Token_MEMBER_FUNCTION_DEF,
	)
	closureCodeRemovalAssertVisitTokens = makeVisitTokenSet(rhino.Token_CALL)
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
	compiler              AbstractCompiler
	removeAbstractMethods bool
	removeAssertionCalls  bool

	// Removable abstract-method assignments discovered during traversal.
	abstractMethodAssignmentNodes []removableAssignment

	// Member function definition nodes annotated with @abstract.
	abstractMemberFunctionNodes []ast.MemberFunctionDef

	// Calls to assertion functions.
	assertionCalls []ast.Call
}

var _ CompilerPass = (*ClosureCodeRemoval)(nil)

// removableAssignment tracks an abstract-method assignment and its enclosing
// chain of ASSIGN ancestors so it can be removed cleanly.
type removableAssignment struct {
	outer           *ClosureCodeRemoval
	assign          ast.Assign
	assignAncestors []rhino.Node
	lastAncestor    rhino.Node
}

func (c *ClosureCodeRemoval) makeRemovableAssignment(assign ast.Assign) removableAssignment {
	a := removableAssignment{
		outer:  c,
		assign: assign,
	}

	ancestor := assign.Node()
	for {
		ancestor = ancestor.GetParent()
		a.assignAncestors = append(a.assignAncestors, ancestor)
		parentAssign, ok := ast.AsAssign(ancestor)
		if !ok || !parentAssign.Target().Node().IsQualifiedName() {
			break
		}
	}
	a.lastAncestor = ancestor.GetParent()
	return a
}

func (a *removableAssignment) remove() {
	compiler := a.outer.compiler
	rhs := a.assign.Value().Node()
	last := a.assign.Node()
	for _, ancestor := range a.assignAncestors {
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

// abstractMethodFinder identifies all assignments of goog.abstractMethod to a
// variable and all methods annotated with "@abstract" in their JSDoc.
type abstractMethodFinder struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *ClosureCodeRemoval
}

func (*abstractMethodFinder) noVisitLeaves() {}

var _ nodeTraversal_visitTokens = (*abstractMethodFinder)(nil)

func (*abstractMethodFinder) visitTokens() *ast.TokenSet {
	return &closureCodeRemovalAbstractVisitTokens
}

func (f *abstractMethodFinder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	c := f.outer
	if assign, ok := ast.AsAssign(n); ok {
		lhs := assign.Target().Node()
		rhs := assign.Value().Node()

		if getProp, ok := assign.Value().AsGetProp(); ok &&
			getProp.Node().GetStringID() == ccrAbstractMethodID &&
			getProp.Object().Node().IsName() &&
			getProp.Object().Node().GetStringID() == ccrGoogID &&
			lhs.IsQualifiedName() {
			// Foo.prototype.bar = goog.abstractMethod
			c.abstractMethodAssignmentNodes = append(c.abstractMethodAssignmentNodes, c.makeRemovableAssignment(assign))
		} else if jsdoc := n.GetJSDocInfo(); jsdoc != nil &&
			jsdoc.IsAbstract() &&
			!jsdoc.IsConstructor() &&
			NodeUtil_IsEmptyFunctionExpression(rhs) {
			// @abstract
			c.abstractMethodAssignmentNodes = append(c.abstractMethodAssignmentNodes, c.makeRemovableAssignment(assign))
		}
	} else if m, ok := ast.AsMemberFunctionDef(n); ok && parent.IsClassMembers() {
		if jsdoc := n.GetJSDocInfo(); jsdoc != nil && jsdoc.IsAbstract() {
			c.abstractMemberFunctionNodes = append(c.abstractMemberFunctionNodes, m)
		}
	}
}

// assertionCallFinder identifies all calls to Closure assertion functions.
type assertionCallFinder struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer          *ClosureCodeRemoval
	assertionNames *CodingConvention_AssertionFunctionLookup
}

func (*assertionCallFinder) noVisitLeaves() {}

var _ nodeTraversal_visitTokens = (*assertionCallFinder)(nil)

func (*assertionCallFinder) visitTokens() *ast.TokenSet {
	return &closureCodeRemovalAssertVisitTokens
}

func (c *ClosureCodeRemoval) newAssertionCallFinder() *assertionCallFinder {
	return &assertionCallFinder{
		outer:          c,
		assertionNames: CodingConvention_AssertionFunctionLookup_Of(c.compiler.GetCodingConvention().GetAssertionFunctions()),
	}
}

func (f *assertionCallFinder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	call, ok := ast.AsCall(n)
	if !ok {
		return
	}

	callee := call.Callee().Node()
	if f.assertionNames.LookupByCallee(callee) != nil { // type-based
		f.outer.assertionCalls = append(f.outer.assertionCalls, call)
	} else if color := callee.GetColor(); color != nil && color.IsClosureAssert() { // color-based
		f.outer.assertionCalls = append(f.outer.assertionCalls, call)
	}
}

// NewClosureCodeRemoval creates a Closure code remover. compiler is the
// AbstractCompiler; removeAbstractMethods removes declarations of abstract
// methods; removeAssertionCalls removes calls to goog.assert functions.
//
// Java: ClosureCodeRemoval(AbstractCompiler, boolean, boolean)
func NewClosureCodeRemoval(compiler AbstractCompiler, removeAbstractMethods bool, removeAssertionCalls bool) *ClosureCodeRemoval {
	return &ClosureCodeRemoval{
		compiler:              compiler,
		removeAbstractMethods: removeAbstractMethods,
		removeAssertionCalls:  removeAssertionCalls,
	}
}

// Java: ClosureCodeRemoval.process(Node, Node)
func (c *ClosureCodeRemoval) Process(externs rhino.Node, root rhino.Node) {
	var passes []NodeTraversal_Callback
	if c.removeAbstractMethods {
		passes = append(passes, &abstractMethodFinder{outer: c})
	}
	if c.removeAssertionCalls {
		passes = append(passes, c.newAssertionCallFinder())
	}
	CombinedCompilerPass_Traverse(c.compiler, root, passes)

	for i := range c.abstractMethodAssignmentNodes {
		c.abstractMethodAssignmentNodes[i].remove()
	}

	for _, memberFn := range c.abstractMemberFunctionNodes {
		c.compiler.ReportFunctionDeleted(memberFn.Function().Node())
		parent := memberFn.Node().GetParent()
		memberFn.Detach()
		c.compiler.ReportChangeToEnclosingScope(parent)
	}

	for _, call := range c.assertionCalls {
		callNode := call.Node()
		// If the assertion is an expression statement, strip the whole statement.
		c.compiler.ReportChangeToEnclosingScope(callNode)
		parent := callNode.GetParent()
		if parent.IsExprResult() {
			parent.Detach()
			NodeUtil_MarkFunctionsDeleted(parent, c.compiler)
		} else {
			// Otherwise, replace the assertion with its first argument,
			// which is the return value of the assertion.
			if firstArg, ok := call.FirstArg(); !ok {
				callNode.ReplaceWith(NodeUtil_NewUndefinedNode(callNode))
			} else {
				replacement := firstArg.Detach().Node()
				replacement.CopyTypeFrom(callNode)
				callNode.ReplaceWith(replacement)
			}
			NodeUtil_MarkFunctionsDeleted(callNode, c.compiler)
		}
	}
}

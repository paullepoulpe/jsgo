// Port of com/google/javascript/jscomp/CheckGlobalThis.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckGlobalThis checks for certain uses of the {@code this} keyword that are
// considered unsafe because they are likely to reference the global
// {@code this} object unintentionally.
//
// A use of {@code this} is considered unsafe if it's on the left side of an
// assignment or a property access, and not inside one of the following:
//
//  1. a class static initialization block
//  2. a prototype method
//  3. a function annotated with {@code @constructor}
//  4. a function annotated with {@code @this}.
//  5. a function where there's no logical place to put a {@code this}
//     annotation.
//
// Note that this check does not track assignments of {@code this} to
// variables or objects. The code
//
//	function evil() {
//	  var a = this;
//	  a.useful = undefined;
//	}
//
// will not get flagged, even though it is semantically equivalent to
//
//	function evil() {
//	  this.useful = undefined;
//	}
//
// which would get flagged.
type CheckGlobalThis struct {
	compiler AbstractCompiler

	// If assignLhsChild != null, then the node being traversed is a
	// descendant of the first child of an ASSIGN node. assignLhsChild's parent
	// is this ASSIGN node.
	assignLhsChild rhino.Node
}

// Java: CheckGlobalThis.GLOBAL_THIS
var CheckGlobalThis_GLOBAL_THIS = DiagnosticType_Warning(
	"JSC_USED_GLOBAL_THIS", "dangerous use of the global 'this' object")

// Java: CheckGlobalThis(AbstractCompiler)
func NewCheckGlobalThis(compiler AbstractCompiler) *CheckGlobalThis {
	c := &CheckGlobalThis{}
	c.assignLhsChild = rhino.NilNode
	c.compiler = compiler
	return c
}

// ShouldTraverse: since this pass reports errors only when a global
// {@code this} keyword is encountered, there is no reason to traverse non
// global contexts.
//
// Java: CheckGlobalThis.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckGlobalThis) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {

	if n.IsFunction() {
		// Arrow functions automatically get the "this" value from the
		// enclosing scope. e.g. the "this" in
		//   Foo.prototype.getBar = () => this.bar;
		// is the global "this", not an instance of Foo.
		if n.IsArrowFunction() {
			return true
		}

		// Don't traverse functions that are constructors or have the @this
		// or @override annotation.
		jsDoc := NodeUtil_GetBestJSDocInfo(n)
		if jsDoc != nil &&
			(jsDoc.IsConstructor() ||
				jsDoc.IsInterface() ||
				jsDoc.HasThisType() ||
				jsDoc.IsOverride()) {
			return false
		}

		if jsDoc != nil {
			functionType := jsDoc.GetType()
			if functionType != nil {
				functionNode := functionType.GetRoot()
				if !functionNode.IsNil() && functionNode.IsFunction() {
					// function(this: ThisType, ...)
					// `this:` is only allowed as the very first child of the
					// FUNCTION node.
					thisNode := functionNode.GetFirstChild()
					if !thisNode.IsNil() && thisNode.IsThis() {
						// Type of `this` is specified, so no need to check further.
						return false
					}
				}
			}
		}

		// Don't traverse functions unless they would normally
		// be able to have a @this annotation associated with them. e.g.,
		// var a = function() { }; // or
		// function a() {} // or
		// a.x = function() {}; // or
		// var a = {x: function() {}};
		pType := parent.GetToken()
		if !(pType == rhino.Token_BLOCK ||
			pType == rhino.Token_SCRIPT ||
			pType == rhino.Token_NAME ||
			pType == rhino.Token_ASSIGN ||
			// object literal keys
			pType == rhino.Token_STRING_KEY) {
			return false
		}

		// Don't traverse functions that are getting lent to a prototype.
		grandparent := parent.GetParent()
		if NodeUtil_MayBeObjectLitKey(parent) {
			maybeLends := grandparent.GetJSDocInfo()
			if maybeLends != nil &&
				maybeLends.HasLendsName() &&
				strings.HasSuffix(maybeLends.GetLendsName().GetRoot().GetString(), ".prototype") {
				return false
			}
		}
	}

	// Don't traverse class static blocks, 'this' is never the global 'this'
	if NodeUtil_IsClassStaticBlock(n) {
		return false
	}

	// `this` in class field definitions refers to the instance (for non-static fields) or the class
	// itself (for static fields), never the global 'this'.
	if (n.IsMemberFieldDef() || n.IsComputedFieldDef()) && n.GetParent().IsClassMembers() {
		return false
	}

	if !parent.IsNil() && parent.IsAssign() {
		lhs := parent.GetFirstChild()

		if n == lhs {
			// Always traverse the left side of the assignment. To handle
			// nested assignments properly (e.g., (a = this).property = c;),
			// assignLhsChild should not be overridden.
			if c.assignLhsChild.IsNil() {
				c.assignLhsChild = lhs
			}
		} else {
			// Only traverse the right side if it's not an assignment to a prototype
			// property or subproperty.
			if NodeUtil_IsNormalGet(lhs) {
				if lhs.IsGetProp() && lhs.GetString() == "prototype" {
					return false
				}
				llhs := lhs.GetFirstChild()
				if llhs.IsGetProp() && llhs.GetString() == "prototype" {
					return false
				}
			}
		}
	}

	return true
}

// Java: CheckGlobalThis.visit(NodeTraversal, Node, Node)
func (c *CheckGlobalThis) Visit(unused *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsThis() && c.ShouldReportThis(n) {
		c.compiler.Report(JSError_MakeWithNode(n, CheckGlobalThis_GLOBAL_THIS))
	}
	if n == c.assignLhsChild {
		c.assignLhsChild = rhino.NilNode
	}
}

// Java: CheckGlobalThis.shouldReportThis(Node)
func (c *CheckGlobalThis) ShouldReportThis(n rhino.Node) bool {
	parent := n.GetParent()
	if !c.assignLhsChild.IsNil() {
		// Always report a THIS on the left side of an assign.
		return true
	}

	// Also report a THIS with a property access.
	return !parent.IsNil() && NodeUtil_IsNormalGet(parent)
}

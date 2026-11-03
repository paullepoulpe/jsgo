// Port of com/google/javascript/jscomp/CheckGlobalThis.java @ closure-compiler v20261028

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var checkGlobalThisVisitTokens = makeVisitTokenSet(rhino.Token_THIS)

var _ nodeTraversal_visitTokens = (*CheckGlobalThis)(nil)

func (*CheckGlobalThis) visitTokens() *ast.TokenSet { return &checkGlobalThisVisitTokens }

// CheckGlobalThis checks for certain uses of the `this` keyword that are
// considered unsafe because they are likely to reference the global `this`
// object unintentionally.
//
// A use of `this` is considered unsafe if it's on the left side of an
// assignment or a property access, and not inside one of the following:
//
//  1. a class static initialization block
//  2. a prototype method
//  3. a function annotated with `@constructor`
//  4. a function annotated with `@this`
//  5. a function where there's no logical place to put a `@this` annotation.
type CheckGlobalThis struct {
	compiler AbstractCompiler
}

// Java: CheckGlobalThis.GLOBAL_THIS
var CheckGlobalThis_GLOBAL_THIS = DiagnosticType_Warning(
	"JSC_USED_GLOBAL_THIS", "dangerous use of the global 'this' object")

// Java: CheckGlobalThis(AbstractCompiler)
func NewCheckGlobalThis(compiler AbstractCompiler) *CheckGlobalThis {
	return &CheckGlobalThis{compiler: compiler}
}

// ShouldTraverse skips non-global `this` contexts.
//
// Java: CheckGlobalThis.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckGlobalThis) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if fn, ok := ast.AsFunction(n); ok {
		// Arrow functions automatically inherit `this` from the enclosing scope;
		// e.g. in `Foo.prototype.getBar = () => this.bar;`, `this` is global `this`.
		if fn.IsArrow() {
			return true
		}

		// Don't traverse functions that are constructors or have @this / @override.
		jsDoc := NodeUtil_GetBestJSDocInfo(n)
		if jsDoc != nil &&
			(jsDoc.IsConstructor() ||
				jsDoc.IsInterface() ||
				jsDoc.HasThisType() ||
				jsDoc.IsOverride()) {
			return false
		}

		if jsDoc != nil {
			if functionType := jsDoc.GetType(); functionType != nil {
				// Note: functionNode is a JSDoc type expression node, not a JS ast.Function.
				functionNode := functionType.GetRoot()
				if functionNode.IsFunction() && functionNode.GetFirstChild().IsThis() {
					return false
				}
			}
		}

		// Only traverse functions where a @this annotation could normally be placed:
		//   var a = function() {};
		//   function a() {}
		//   a.x = function() {};
		//   var a = {x: function() {}};
		switch parent.GetToken() {
		case rhino.Token_BLOCK, rhino.Token_SCRIPT, rhino.Token_NAME, rhino.Token_ASSIGN, rhino.Token_STRING_KEY:
		default:
			return false
		}

		// Don't traverse functions lent to a prototype via @lends.
		if NodeUtil_MayBeObjectLitKey(parent) {
			if maybeLends := parent.GetParent().GetJSDocInfo(); maybeLends != nil &&
				maybeLends.HasLendsName() &&
				strings.HasSuffix(maybeLends.GetLendsName().GetRoot().GetString(), ".prototype") {
				return false
			}
		}
	}

	// Don't traverse class static blocks; `this` is never global `this`.
	if NodeUtil_IsClassStaticBlock(n) {
		return false
	}

	// `this` in class field definitions refers to the instance or the class, never global `this`.
	if _, ok := ast.AsFieldDef(n); ok && parent.IsClassMembers() {
		return false
	}

	if assign, ok := ast.AsAssign(parent); ok && !assign.Target().Is(n) {
		// Only traverse the right side if it's not an assignment to a prototype
		// property or subproperty.
		if lhsGet, ok := ast.AsNormalGet(assign.Target().Node()); ok {
			if gp, ok := lhsGet.AsGetProp(); ok && gp.Property() == "prototype" {
				return false
			}
			if objGp, ok := lhsGet.Object().AsGetProp(); ok && objGp.Property() == "prototype" {
				return false
			}
		}
	}

	return true
}

// Java: CheckGlobalThis.visit(NodeTraversal, Node, Node)
func (c *CheckGlobalThis) Visit(unused *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsThis() && shouldReportGlobalThis(n) {
		c.compiler.Report(JSError_MakeWithNode(n, CheckGlobalThis_GLOBAL_THIS))
	}
}

func shouldReportGlobalThis(n rhino.Node) bool {
	parent := n.GetParent()
	if _, ok := ast.AsNormalGet(parent); ok {
		// Report a THIS with a property access.
		return true
	}
	// Always report a THIS on the left side of an assign.
	for cur, p := n, parent; !p.IsNil() && !p.IsScript(); cur, p = p, p.GetParent() {
		if ast.RoleIn(cur, p) == ast.RoleAssignTarget {
			return true
		}
	}
	return false
}

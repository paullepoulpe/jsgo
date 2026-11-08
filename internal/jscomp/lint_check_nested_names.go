// Port of com/google/javascript/jscomp/lint/CheckNestedNames.java @ closure-compiler v20261104

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var checkNestedNamesVisitTokens = makeVisitTokenSet(rhino.Token_GETPROP)

var _ nodeTraversal_visitTokens = (*CheckNestedNames)(nil)

func (*CheckNestedNames) noVisitLeaves()             {}
func (*CheckNestedNames) visitTokens() *ast.TokenSet { return &checkNestedNamesVisitTokens }

// CheckNestedNames checks if code has a module-level static property
// assignment (`X.Y`) inside a `goog.module`. If yes, reports a linter warning
// on the `X.Y` if:
//
//  1. X is any name {object, class, function, interface, enum, typedef}.
//  2. Y is {class, interface, enum, typedef}
//
// For example, reports a linter finding in the following:
//
//	class C {}
//	let obj = {...};
//
//	/**
//	 * Reports nested name E on class C.
//	 * @enum {number}
//	 * /
//	C.E={...}; // WARNING:
//
//	/**
//	 * Reports nested interface I on variable obj.
//	 * @interface
//	 * /
//	obj.I=class {}; // WARNING:
//
//	// reports nested name D on class C.
//	C.D=class {};
//
// Does not report on nested functions, or any nested property assigned to a
// value:
//
//	C.F = function() {}
//	obj.S = '';
//	obj.Some = new Something();
type CheckNestedNames struct {
	compiler AbstractCompiler
}

var CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE = DiagnosticType_Disabled(
	"JSC_NESTED_NAME_IN_GOOG_MODULE",
	"A nested {0} is created on the name `{1}`."+
		" Fix this linter finding by converting the module-level static property"+
		" assignment on `{1}` into a module-level flat name (i.e. change `{1}.prop = ...`"+
		" into `{1}_prop = ...`. You can (if required) export this flat name using named"+
		" exports (`exports.{1}_prop = {1}_prop`).")

// NewCheckNestedNames constructs a CheckNestedNames pass.
//
// Java: CheckNestedNames(AbstractCompiler)
func NewCheckNestedNames(compiler AbstractCompiler) *CheckNestedNames {
	return &CheckNestedNames{compiler: compiler}
}

func (*CheckNestedNames) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckNestedNames(c)
}

// Java: CheckNestedNames.process(Node, Node)
func (c *CheckNestedNames) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckNestedNames.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckNestedNames) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		// Only traverse inside goog.modules.
		return n.GetBooleanProp(rhino.Node_GOOG_MODULE)
	}
	// Only warn on static property assignments of module-level names.
	return nodeTraversal.InGlobalOrModuleScope()
}

// Java: CheckNestedNames.visit(NodeTraversal, Node, Node)
func (c *CheckNestedNames) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	jutil.CheckArgument(t.InGlobalOrModuleScope())
	gp, ok := ast.AsGetProp(n)
	if !ok || !n.IsQualifiedName() {
		return
	}
	role := ast.RoleIn(n, parent)
	if role != ast.RoleAssignTarget && role != ast.RoleExprStmtExpr {
		return
	}
	target := outermostTargetGetProp(gp)
	if target.Property() == "prototype" {
		return
	}
	owner, ok := target.Object().AsName()
	if !ok || owner.Ident() == "exports" {
		return // For example `this` or `super` or `exports.SomeEnum = {}`
	}
	if nestedKind, ok := nestedDeclarationKind(n); ok {
		// Found a nested name that created an inner class, interface, typedef, or enum.
		t.Report(target.Node(), CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE, nestedKind, owner.Ident())
	}
}

// outermostTargetGetProp returns `a.b` from `a.b.c.d...`.
func outermostTargetGetProp(gp ast.GetProp) ast.GetProp {
	for inner, ok := gp.Object().AsGetProp(); ok; inner, ok = gp.Object().AsGetProp() {
		gp = inner
	}
	return gp
}

func nestedDeclarationKind(lhs rhino.Node) (string, bool) {
	// Handle `/** @typedef {...} */ X.Y` or `/** @const */ X.Y = SomeTypeDefName`
	if lhs.GetTypedefTypeProp() != nil {
		return "typedef", true
	}

	if type_ := lhs.GetJSType(); type_ != nil {
		// Check whether it's an enum, interface, or class constructor.
		if type_.IsEnumType() {
			return "enum", true
		} else if type_.IsInterface() {
			return "interface", true
		} else if type_.IsConstructor() {
			return "class", true
		}
	}

	return "", false
}

// Port of com/google/javascript/jscomp/lint/CheckNestedNames.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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

// CheckNestedNames_DeclarationKind is the private enum DeclarationKind.
type CheckNestedNames_DeclarationKind int32

const (
	CheckNestedNames_DeclarationKind_CLASS CheckNestedNames_DeclarationKind = iota
	CheckNestedNames_DeclarationKind_ENUM
	CheckNestedNames_DeclarationKind_TYPEDEF
	CheckNestedNames_DeclarationKind_INTERFACE
)

var checkNestedNames_DeclarationKind_names = [...]string{"CLASS", "ENUM", "TYPEDEF", "INTERFACE"}

// CheckNestedNames_DeclarationKind_values is DeclarationKind.values().
func CheckNestedNames_DeclarationKind_values() []CheckNestedNames_DeclarationKind {
	return []CheckNestedNames_DeclarationKind{
		CheckNestedNames_DeclarationKind_CLASS,
		CheckNestedNames_DeclarationKind_ENUM,
		CheckNestedNames_DeclarationKind_TYPEDEF,
		CheckNestedNames_DeclarationKind_INTERFACE,
	}
}

// Name is Enum.name().
func (k CheckNestedNames_DeclarationKind) Name() string {
	return checkNestedNames_DeclarationKind_names[k]
}

// String is Enum.toString().
func (k CheckNestedNames_DeclarationKind) String() string { return k.Name() }

// Ordinal is Enum.ordinal().
func (k CheckNestedNames_DeclarationKind) Ordinal() int { return int(k) }

// Go-only: language-adapter — a Java DeclarationKind reference may be null
// (getNestedDeclarationKind returns null); such values are
// *CheckNestedNames_DeclarationKind here, and Ptr returns the canonical
// pointer to a constant.
var checkNestedNames_DeclarationKind_ptrs = [...]CheckNestedNames_DeclarationKind{
	CheckNestedNames_DeclarationKind_CLASS,
	CheckNestedNames_DeclarationKind_ENUM,
	CheckNestedNames_DeclarationKind_TYPEDEF,
	CheckNestedNames_DeclarationKind_INTERFACE,
}

// Ptr returns the canonical non-null *CheckNestedNames_DeclarationKind for k.
func (k CheckNestedNames_DeclarationKind) Ptr() *CheckNestedNames_DeclarationKind {
	return &checkNestedNames_DeclarationKind_ptrs[k]
}

var checkNestedNames_NESTED_KIND_TO_REPORT = func() *jutil.LinkedHashMap[CheckNestedNames_DeclarationKind, string] {
	m := jutil.NewLinkedHashMap[CheckNestedNames_DeclarationKind, string]()
	m.Put(CheckNestedNames_DeclarationKind_CLASS, "class")
	m.Put(CheckNestedNames_DeclarationKind_ENUM, "enum")
	m.Put(CheckNestedNames_DeclarationKind_INTERFACE, "interface")
	m.Put(CheckNestedNames_DeclarationKind_TYPEDEF, "typedef")
	return m
}()

// NewCheckNestedNames is new CheckNestedNames(compiler).
//
// Java: CheckNestedNames(AbstractCompiler)
func NewCheckNestedNames(compiler AbstractCompiler) *CheckNestedNames {
	return &CheckNestedNames{compiler: compiler}
}

// Java: CheckNestedNames.process(Node, Node)
func (c *CheckNestedNames) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckNestedNames.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckNestedNames) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		// only traverse inside goog.modules
		return n.GetBooleanProp(rhino.Node_GOOG_MODULE)
	}
	// only warn on static property assignments of module-level names
	return nodeTraversal.InGlobalOrModuleScope()
}

// Java: CheckNestedNames.visit(NodeTraversal, Node, Node)
func (c *CheckNestedNames) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	jutil.CheckArgument(t.InGlobalOrModuleScope())
	if !n.IsGetProp() || !n.IsQualifiedName() {
		return
	}
	targetGetProp := CheckNestedNames_GetTargetGetProp(n)
	if (parent.IsAssign() || parent.IsExprResult()) &&
		n.IsFirstChildOf(parent) &&
		!CheckNestedNames_IsDotPrototype(targetGetProp) {
		owner := targetGetProp.GetFirstChild()
		if !owner.IsName() || owner.GetString() == "exports" {
			return //  For example `this` or `super` or `exports.SomeEnum = {}`
		}
		ownerName := owner.GetString()
		declarationKind := c.GetNestedDeclarationKind(n, parent)
		if declarationKind != nil && checkNestedNames_NESTED_KIND_TO_REPORT.ContainsKey(*declarationKind) {
			// Found a nested name that created an inner class, interface, typedef or enum.
			nestedKind := checkNestedNames_NESTED_KIND_TO_REPORT.Get(*declarationKind)
			t.Report(targetGetProp, CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE, nestedKind, ownerName)
		}
	}
}

// CheckNestedNames_GetTargetGetProp gives `a.b` from `a.b.c.d...`.
//
// Java: CheckNestedNames.getTargetGetProp(Node)
func CheckNestedNames_GetTargetGetProp(n rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(n.IsGetProp(), "%s", n)
	if n.GetFirstChild().IsGetProp() {
		return CheckNestedNames_GetTargetGetProp(n.GetFirstChild())
	}
	return n
}

// CheckNestedNames_IsDotPrototype: true for `someExpression.prototype`.
//
// Java: CheckNestedNames.isDotPrototype(Node)
func CheckNestedNames_IsDotPrototype(getProp rhino.Node) bool {
	return getProp.IsGetProp() && getProp.GetString() == "prototype"
}

// Java: CheckNestedNames.getNestedDeclarationKind(Node, Node)
func (c *CheckNestedNames) GetNestedDeclarationKind(lhs rhino.Node, parent rhino.Node) *CheckNestedNames_DeclarationKind {
	jutil.CheckArgumentMsg1(lhs.IsFirstChildOf(parent), "%s", lhs)
	// Handle `/** @typedef {...} */` X.Y` or `/** @const */ X.Y = SomeTypeDefName`
	if lhs.GetTypedefTypeProp() != nil {
		return CheckNestedNames_DeclarationKind_TYPEDEF.Ptr()
	}

	type_ := lhs.GetJSType()
	if type_ != nil {
		// check whether it's an enum, interface or typedef first.
		if type_.IsEnumType() {
			return CheckNestedNames_DeclarationKind_ENUM.Ptr()
		} else if type_.IsInterface() {
			return CheckNestedNames_DeclarationKind_INTERFACE.Ptr()
		} else if type_.IsConstructor() {
			return CheckNestedNames_DeclarationKind_CLASS.Ptr()
		}
	}

	return nil
}

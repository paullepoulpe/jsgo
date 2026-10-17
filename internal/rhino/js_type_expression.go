// Port of com/google/javascript/rhino/JSTypeExpression.java @ closure-compiler v20261015

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSTypeExpression: when parsing a jsdoc, a type-annotation string is parsed
// to a type AST. Somewhat confusingly, we use the Node class both for type
// ASTs and for the source-code AST. JSTypeExpression wraps a type AST.
// During type checking, type ASTs are evaluated to JavaScript types.
type JSTypeExpression struct {
	// The root of the AST.
	root *Node

	// The source name where the type expression appears.
	sourceName string
}

// Java: JSTypeExpression.serialVersionUID
const JSTypeExpression_serialVersionUID int64 = 1

// Java: JSTypeExpression.IMPLICIT_TEMPLATE_BOUND_SOURCE
const JSTypeExpression_IMPLICIT_TEMPLATE_BOUND_SOURCE = "<IMPLICIT_TEMPLATE_BOUND>"

// JSTypeExpression_IMPLICIT_TEMPLATE_BOUND (with the static initializer
// that sets its root's source file).
//
// Java: JSTypeExpression.IMPLICIT_TEMPLATE_BOUND
var JSTypeExpression_IMPLICIT_TEMPLATE_BOUND = func() *JSTypeExpression {
	e := NewJSTypeExpression(NewNode(Token_QMARK), JSTypeExpression_IMPLICIT_TEMPLATE_BOUND_SOURCE)
	e.GetRoot().SetStaticSourceFile(
		NewSimpleSourceFile(
			JSTypeExpression_IMPLICIT_TEMPLATE_BOUND_SOURCE, StaticSourceFile_SourceKind_STRONG))
	return e
}()

// Java: JSTypeExpression(Node, String)
func NewJSTypeExpression(root *Node, sourceName string) *JSTypeExpression {
	return &JSTypeExpression{root: root, sourceName: sourceName}
}

// ReplaceNamesWithUnknownType replaces given names in this type expression
// with unknown.
//
// Java: JSTypeExpression.replaceNamesWithUnknownType(Set<String>)
func (e *JSTypeExpression) ReplaceNamesWithUnknownType(names *jutil.LinkedHashSet[string]) *JSTypeExpression {
	oldExprRoot := e.root.CloneTree()
	newExprRoot := JSTypeExpression_ReplaceNames(oldExprRoot, names)
	newTypeExpression := NewJSTypeExpression(newExprRoot, e.sourceName)
	return newTypeExpression
}

// JSTypeExpression_ReplaceNames recursively traverses over the type tree and
// replaces matched types with unknown type.
//
// n: root of the JsTypeExpression on which replacement is applied. names:
// the set of names to replace in this type expression. Returns the new root
// after replacing the names.
//
// Java: JSTypeExpression.replaceNames(Node, Set<String>)
func JSTypeExpression_ReplaceNames(n *Node, names *jutil.LinkedHashSet[string]) *Node {
	if n == nil {
		return nil
	}
	for child := n.GetFirstChild(); child != nil; {
		next := child.GetNext()
		JSTypeExpression_ReplaceNames(child, names)
		child = next
	}
	if n.IsStringLit() && names.Contains(JSTypeExpression_BaseName(n.GetString())) {
		qMark := NewNode(Token_QMARK)
		if n.HasChildren() && !n.GetFirstChild().IsBlock() {
			// preserve the children of the type expression unless it's a block -
			// a BLOCK child represents a templatized type Foo<string, number>, and
			// `?` cannot be templatized.
			qMark.AddChildrenToBack(n.RemoveChildren())
		}
		if n.HasParent() {
			n.ReplaceWith(qMark)
		}
		return qMark
	}
	return n
}

// Java: JSTypeExpression.baseName(String)
func JSTypeExpression_BaseName(name string) string {
	dotIndex := jutil.StringIndexOfChar(name, '.')
	if dotIndex == -1 {
		return name
	}
	return jutil.StringSubstring(name, 0, dotIndex)
}

// GetAllTypeNodes returns a list of all type nodes in this type expression.
//
// Java: JSTypeExpression.getAllTypeNodes()
func (e *JSTypeExpression) GetAllTypeNodes() []*Node {
	var builder []*Node
	JSTypeExpression_VisitAllTypeNodes(e.root, func(n *Node) { builder = append(builder, n) })
	return builder
}

// GetAllTypeNames returns a set of all string names in this type expression.
//
// Java: JSTypeExpression.getAllTypeNames()
func (e *JSTypeExpression) GetAllTypeNames() *jutil.LinkedHashSet[string] {
	builder := jutil.NewLinkedHashSet[string]()
	JSTypeExpression_VisitAllTypeNodes(e.root, func(n *Node) { builder.Add(n.GetString()) })
	return builder
}

// JSTypeExpression_VisitAllTypeNodes recursively traverses the type tree and
// visits all type nodes.
//
// Java: JSTypeExpression.visitAllTypeNodes(Node, Consumer<Node>)
func JSTypeExpression_VisitAllTypeNodes(n *Node, visitor func(*Node)) {
	if n == nil {
		return
	}
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		JSTypeExpression_VisitAllTypeNodes(child, visitor)
	}
	if n.IsStringLit() {
		visitor(n)
	}
}

// JSTypeExpression_MakeOptionalArg makes the given type expression into an
// optional type expression, if possible.
//
// Java: JSTypeExpression.makeOptionalArg(JSTypeExpression)
func JSTypeExpression_MakeOptionalArg(expr *JSTypeExpression) *JSTypeExpression {
	if expr.IsOptionalArg() || expr.IsVarArgs() {
		return expr
	} else {
		equals := NewNodeWithChild(Token_EQUALS, expr.root)
		equals.ClonePropsFrom(expr.root)
		return NewJSTypeExpression(equals, expr.sourceName)
	}
}

// JSTypeExpression_MakeVarArgs makes the given type expression into an
// optional type expression, if possible.
//
// Java: JSTypeExpression.makeVarArgs(JSTypeExpression)
func JSTypeExpression_MakeVarArgs(expr *JSTypeExpression) *JSTypeExpression {
	if expr.IsOptionalArg() || expr.IsVarArgs() {
		return expr
	} else {
		equals := NewNodeWithChild(Token_ITER_REST, expr.root)
		equals.ClonePropsFrom(expr.root)
		return NewJSTypeExpression(equals, expr.sourceName)
	}
}

// IsOptionalArg reports whether this expression denotes an optional @param.
//
// Java: JSTypeExpression.isOptionalArg()
func (e *JSTypeExpression) IsOptionalArg() bool {
	return e.root.GetToken() == Token_EQUALS
}

// IsVarArgs reports whether this expression denotes a rest args @param.
//
// Java: JSTypeExpression.isVarArgs()
func (e *JSTypeExpression) IsVarArgs() bool {
	return e.root.GetToken() == Token_ITER_REST
}

// Evaluate evaluates the type expression into a JSType object.
//
// Java: JSTypeExpression.evaluate(StaticTypedScope, JSTypeRegistry)
func (e *JSTypeExpression) Evaluate(scope StaticTypedScope, registry *JSTypeRegistry) JSType {
	type_ := registry.CreateTypeFromCommentNodeWithScope(e.root, e.sourceName, scope)
	e.root.SetJSType(type_)
	return type_
}

// IsEquivalentTo reports whether this object represents a type expression
// that is equivalent to the other one.
//
// Java: JSTypeExpression.isEquivalentTo(JSTypeExpression)
func (e *JSTypeExpression) IsEquivalentTo(other *JSTypeExpression) bool {
	return other != nil && e.root.IsEquivalentTo(other.root)
}

// GetRoot returns the source for this type expression. Note that it will not
// contain an expression if there's an @override tag.
//
// Java: JSTypeExpression.getRoot()
func (e *JSTypeExpression) GetRoot() *Node {
	return e.root
}

// Java: JSTypeExpression.getSourceName()
func (e *JSTypeExpression) GetSourceName() string {
	return e.sourceName
}

// Java: JSTypeExpression.toString()
func (e *JSTypeExpression) String() string {
	return "type: " + e.root.ToStringTree()
}

// Java: JSTypeExpression.copy()
func (e *JSTypeExpression) Copy() *JSTypeExpression {
	return NewJSTypeExpression(e.root.CloneTree(), e.sourceName)
}

// IsExplicitUnknownTemplateBound reports whether this expression is an
// explicit unknown template bound.
//
// Java: JSTypeExpression.isExplicitUnknownTemplateBound()
func (e *JSTypeExpression) IsExplicitUnknownTemplateBound() bool {
	return base.JSCompObjects_Identical(e.root.GetToken(), Token_QMARK) &&
		!e.root.HasChildren() &&
		e.sourceName != JSTypeExpression_IMPLICIT_TEMPLATE_BOUND_SOURCE
}

// GetRecordPropertyNames returns a set of keys of all record types (e.g.
// {{key : string }}) present in this JSTypeExpression.
//
// Java: JSTypeExpression.getRecordPropertyNames()
func (e *JSTypeExpression) GetRecordPropertyNames() *jutil.LinkedHashSet[string] {
	builder := jutil.NewLinkedHashSet[string]()
	JSTypeExpression_GetRecordPropertyNamesRecursive(e.root, builder)
	return builder
}

// Java: JSTypeExpression.getRecordPropertyNamesRecursive(Node, ImmutableSet.Builder<String>)
func JSTypeExpression_GetRecordPropertyNamesRecursive(n *Node, names *jutil.LinkedHashSet[string]) {
	if n == nil {
		return
	}

	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		JSTypeExpression_GetRecordPropertyNamesRecursive(child, names)
	}
	if n.IsStringKey() {
		names.Add(n.GetString())
	}
}

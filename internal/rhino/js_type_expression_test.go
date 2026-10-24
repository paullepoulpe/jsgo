// Port of com/google/javascript/rhino/JSTypeExpressionTest.java @ closure-compiler v20261021

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
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// assertContainsExactly is Truth's containsExactly (any order).
func assertContainsExactly(t *testing.T, got []string, want ...string) {
	t.Helper()
	g, w := slices.Clone(got), slices.Clone(want)
	slices.Sort(g)
	slices.Sort(w)
	if !slices.Equal(g, w) {
		t.Errorf("got %q, want exactly %q", got, want)
	}
}

func TestJSTypeExpression_GetAllTypeNames(t *testing.T) {
	assertContainsExactly(t, getTestExpression().GetAllTypeNames().ToSlice(),
		"foo.Bar", "string", "Object")
}

func TestJSTypeExpression_GetAllTypeNodes(t *testing.T) {
	var names []string
	for _, n := range getTestExpression().GetAllTypeNodes() {
		names = append(names, n.GetString())
	}
	assertContainsExactly(t, names, "foo.Bar", "string", "Object", "string")
}

func TestJSTypeExpression_IsExplicitUnknownTemplateBound(t *testing.T) {
	// Explicit unknown template bound should return true.
	if !NewJSTypeExpression(NewNode(Token_QMARK), "").IsExplicitUnknownTemplateBound() {
		t.Error("explicit ? bound: want true")
	}

	// The implicit bound should return false, including copies of it.
	implicitBoundCopy := JSTypeExpression_IMPLICIT_TEMPLATE_BOUND.Copy()
	if implicitBoundCopy == JSTypeExpression_IMPLICIT_TEMPLATE_BOUND {
		t.Error("copy is the same instance as IMPLICIT_TEMPLATE_BOUND")
	}
	if JSTypeExpression_IMPLICIT_TEMPLATE_BOUND.IsExplicitUnknownTemplateBound() {
		t.Error("IMPLICIT_TEMPLATE_BOUND: want false")
	}
	if implicitBoundCopy.IsExplicitUnknownTemplateBound() {
		t.Error("copy of IMPLICIT_TEMPLATE_BOUND: want false")
	}

	// Non-unknown bounds should return false, even if they start with a "?".
	if NewJSTypeExpression(NewNode(Token_STAR), "").IsExplicitUnknownTemplateBound() {
		t.Error("* bound: want false")
	}
	if NewJSTypeExpression(NewNodeWithChild(Token_QMARK, Node_NewString("number")), "").
		IsExplicitUnknownTemplateBound() {
		t.Error("?number bound: want false")
	}
}

func getTestExpression() *JSTypeExpression {
	a := Node_NewString("foo.Bar")
	b := Node_NewString("string")
	c := NewNode(Token_PIPE)
	c.AddChildToBack(a)
	c.AddChildToBack(b)
	d := Node_NewString("Object")
	e := Node_NewString("string")
	f := NewNode(Token_PIPE)
	f.AddChildToBack(c)
	f.AddChildToBack(d)
	f.AddChildToBack(e)
	return NewJSTypeExpression(f, "<JSTypeExpressionTest>")
}

func TestJSTypeExpression_MakeOptionalArgAndVarArgsDoesNotMutateRoot(t *testing.T) {
	root := NewNode(Token_NAME)
	expr := NewJSTypeExpression(root, "source.js")

	optArg := JSTypeExpression_MakeOptionalArg(expr)
	if got := optArg.GetRoot().GetToken(); got != Token_EQUALS {
		t.Errorf("optArg root token: got %v, want EQUALS", got)
	}
	if !optArg.GetRoot().GetParent().IsNil() {
		t.Error("optArg root has a parent")
	}

	varArg := JSTypeExpression_MakeVarArgs(expr)
	if got := varArg.GetRoot().GetToken(); got != Token_ITER_REST {
		t.Errorf("varArg root token: got %v, want ITER_REST", got)
	}
	if !varArg.GetRoot().GetParent().IsNil() {
		t.Error("varArg root has a parent")
	}

	// root remains detached
	if !expr.GetRoot().GetParent().IsNil() {
		t.Error("expr root has a parent")
	}

	// We can call it multiple times on the same expression without throwing
	optArg2 := JSTypeExpression_MakeOptionalArg(expr)
	if got := optArg2.GetRoot().GetToken(); got != Token_EQUALS {
		t.Errorf("optArg2 root token: got %v, want EQUALS", got)
	}
	if !optArg2.GetRoot().GetParent().IsNil() {
		t.Error("optArg2 root has a parent")
	}

	// replaceNamesWithUnknownType does not mutate the original root
	replaced := expr.ReplaceNamesWithUnknownType(jutil.NewLinkedHashSet("foo"))
	if !replaced.GetRoot().GetParent().IsNil() {
		t.Error("replaced root has a parent")
	}
	if !expr.GetRoot().GetParent().IsNil() {
		t.Error("expr root has a parent after replaceNamesWithUnknownType")
	}
}

func TestJSTypeExpression_ConstructWithAttachedNodeClonesRoot(t *testing.T) {
	parent := NewNode(Token_BLOCK)
	child := NewNode(Token_NAME)
	parent.AddChildToBack(child)
	if child.GetParent() != parent {
		t.Fatal("child's parent is not parent")
	}

	expr := NewJSTypeExpression(child, "source.js")
	if expr.GetRoot() == child {
		t.Error("expr root is the attached child")
	}
	if !expr.GetRoot().GetParent().IsNil() {
		t.Error("expr root has a parent")
	}
	if child.GetParent() != parent {
		t.Error("child's parent changed")
	}
}

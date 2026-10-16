// Port of com/google/javascript/rhino/JSTypeExpressionTest.java @ closure-compiler v20260415

package rhino

import (
	"slices"
	"testing"
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
	return NewJSTypeExpression(f, "")
}

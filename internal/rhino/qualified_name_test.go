// Port of com/google/javascript/rhino/QualifiedNameTest.java @ closure-compiler v20261015

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func qname(root *Node, props ...string) *Node {
	n := root
	for _, p := range props {
		n = IR_Getprop(n, p)
	}
	return n
}

func assertJoin(t *testing.T, n QualifiedName, want string) {
	t.Helper()
	if got := n.Join(); got != want {
		t.Errorf("join() = %q, want %q", got, want)
	}
}

func assertJoinOnSeparator(t *testing.T, n QualifiedName, separator uint16, want string) {
	t.Helper()
	if got := n.JoinWithSeparator(separator); got != want {
		t.Errorf("join(%q) = %q, want %q", rune(separator), got, want)
	}
}

func assertComponents(t *testing.T, n QualifiedName, want ...string) {
	t.Helper()
	if got := n.Components(); !slices.Equal(got, want) {
		t.Errorf("components() = %q, want %q", got, want)
	}
}

func assertComponent(t *testing.T, n QualifiedName, want string) {
	t.Helper()
	if got := n.GetComponent(); got != want {
		t.Errorf("getComponent() = %q, want %q", got, want)
	}
}

func assertMatches(t *testing.T, n QualifiedName, node *Node, want bool) {
	t.Helper()
	if got := n.Matches(node); got != want {
		t.Errorf("%s.matches(%s) = %v, want %v", n.Join(), node.ToStringTree(), got, want)
	}
}

func assertComponentCount(t *testing.T, n QualifiedName, want int) {
	t.Helper()
	if got := n.GetComponentCount(); got != want {
		t.Errorf("getComponentCount() = %d, want %d", got, want)
	}
}

func TestQualifiedName_Join_FromString(t *testing.T) {
	n := QualifiedName_Of("foo.bar.baz")
	assertJoin(t, n, "foo.bar.baz")
}

func TestQualifiedName_Join_OnSeparator_FromString(t *testing.T) {
	n := QualifiedName_Of("foo.bar.baz")
	assertJoinOnSeparator(t, n, '$', "foo$bar$baz")
}

func TestQualifiedName_Join_FromNode(t *testing.T) {
	n := qname(IR_Name("foo"), "bar", "baz").GetQualifiedNameObject()
	assertJoin(t, n, "foo.bar.baz")

	n = qname(IR_ThisNode(), "bar", "baz").GetQualifiedNameObject()
	assertJoin(t, n, "this.bar.baz")

	n = qname(IR_SuperNode(), "bar", "baz").GetQualifiedNameObject()
	assertJoin(t, n, "super.bar.baz")
}

func TestQualifiedName_Join_OnSeparator_FromNode(t *testing.T) {
	n := qname(IR_Name("foo"), "bar", "baz").GetQualifiedNameObject()
	assertJoinOnSeparator(t, n, '$', "foo$bar$baz")

	n = qname(IR_ThisNode(), "bar", "baz").GetQualifiedNameObject()
	assertJoinOnSeparator(t, n, '$', "this$bar$baz")

	n = qname(IR_SuperNode(), "bar", "baz").GetQualifiedNameObject()
	assertJoinOnSeparator(t, n, '$', "super$bar$baz")
}

func TestQualifiedName_Join_FromGetprop(t *testing.T) {
	n := QualifiedName_Of("foo").Getprop("bar").Getprop("baz")
	assertJoin(t, n, "foo.bar.baz")
}

func TestQualifiedName_Join_OnSeparator_FromGetprop(t *testing.T) {
	n := QualifiedName_Of("foo").Getprop("bar").Getprop("baz")
	assertJoinOnSeparator(t, n, '$', "foo$bar$baz")
}

func TestQualifiedName_Components_FromString(t *testing.T) {
	n := QualifiedName_Of("foo.bar.baz")
	assertComponents(t, n, "foo", "bar", "baz")
}

func TestQualifiedName_Components_FromNode(t *testing.T) {
	n := qname(IR_Name("foo"), "bar", "baz").GetQualifiedNameObject()
	assertComponents(t, n, "foo", "bar", "baz")
}

func TestQualifiedName_Components_FromGetprop(t *testing.T) {
	n := QualifiedName_Of("foo").Getprop("bar").Getprop("baz")
	assertComponents(t, n, "foo", "bar", "baz")
}

func TestQualifiedName_GetOwner_FromString(t *testing.T) {
	assertJoin(t, QualifiedName_Of("foo.bar.baz").GetOwner(), "foo.bar")
	if o := QualifiedName_Of("foo").GetOwner(); o != nil {
		t.Errorf("owner of foo = %v, want null", o)
	}
}

func TestQualifiedName_GetOwner_FromNode(t *testing.T) {
	n := qname(IR_Name("foo"), "bar", "baz").GetQualifiedNameObject()
	assertJoin(t, n.GetOwner(), "foo.bar")
	for _, root := range []*Node{IR_Name("foo"), IR_ThisNode(), IR_SuperNode()} {
		if o := root.GetQualifiedNameObject().GetOwner(); o != nil {
			t.Errorf("owner of %s = %v, want null", root, o)
		}
	}
}

func TestQualifiedName_GetOwner_FromGetprop(t *testing.T) {
	foo := QualifiedName_Of("foo")
	fooBar := foo.Getprop("bar")
	if fooBar.GetOwner() != foo {
		t.Error("fooBar.getOwner() != foo")
	}
	if fooBar.Getprop("baz").GetOwner() != fooBar {
		t.Error("fooBar.getprop(baz).getOwner() != fooBar")
	}
}

func TestQualifiedName_GetComponent_FromString(t *testing.T) {
	assertComponent(t, QualifiedName_Of("foo.bar.baz"), "baz")
	assertComponent(t, QualifiedName_Of("foo"), "foo")
}

func TestQualifiedName_GetComponent_FromNode(t *testing.T) {
	n := qname(IR_Name("foo"), "bar", "baz").GetQualifiedNameObject()
	assertComponent(t, n, "baz")
	assertComponent(t, IR_Name("foo").GetQualifiedNameObject(), "foo")
	assertComponent(t, IR_ThisNode().GetQualifiedNameObject(), "this")
	assertComponent(t, IR_SuperNode().GetQualifiedNameObject(), "super")
}

func TestQualifiedName_GetComponent_FromGetprop(t *testing.T) {
	foo := QualifiedName_Of("foo")
	fooBar := foo.Getprop("bar")
	assertComponent(t, fooBar, "bar")
	assertComponent(t, fooBar.Getprop("baz"), "baz")
}

func TestQualifiedName_Match_FromString(t *testing.T) {
	n := QualifiedName_Of("foo.bar.baz")
	assertMatches(t, n, qname(IR_Name("foo"), "bar", "baz"), true)
	assertMatches(t, n, qname(IR_ThisNode(), "bar", "baz"), false)
	assertMatches(t, n, qname(IR_SuperNode(), "bar", "baz"), false)
	assertMatches(t, n, qname(IR_Name("foo"), "baz", "bar"), false)

	n = QualifiedName_Of("this.qux")
	assertMatches(t, n, qname(IR_ThisNode(), "qux"), true)
	assertMatches(t, n, qname(IR_Name("x"), "qux"), false)

	n = QualifiedName_Of("super.qux")
	assertMatches(t, n, qname(IR_SuperNode(), "qux"), true)
	assertMatches(t, n, qname(IR_Name("x"), "qux"), false)
}

func TestQualifiedName_Match_FromNode(t *testing.T) {
	n := qname(IR_Name("foo"), "bar", "baz").GetQualifiedNameObject()
	assertMatches(t, n, qname(IR_Name("foo"), "bar", "baz"), true)
	assertMatches(t, n, qname(IR_ThisNode(), "bar", "baz"), false)
	assertMatches(t, n, qname(IR_SuperNode(), "bar", "baz"), false)
	assertMatches(t, n, qname(IR_Name("foo"), "baz", "bar"), false)

	n = qname(IR_ThisNode(), "qux").GetQualifiedNameObject()
	assertMatches(t, n, qname(IR_ThisNode(), "qux"), true)
	assertMatches(t, n, qname(IR_Name("x"), "qux"), false)

	n = qname(IR_SuperNode(), "qux").GetQualifiedNameObject()
	assertMatches(t, n, qname(IR_SuperNode(), "qux"), true)
	assertMatches(t, n, qname(IR_Name("x"), "qux"), false)
}

func TestQualifiedName_Match_FromGetprop(t *testing.T) {
	n := QualifiedName_Of("foo").Getprop("bar").Getprop("baz")
	assertMatches(t, n, qname(IR_Name("foo"), "bar", "baz"), true)
	assertMatches(t, n, qname(IR_ThisNode(), "bar", "baz"), false)
	assertMatches(t, n, qname(IR_SuperNode(), "bar", "baz"), false)
	assertMatches(t, n, qname(IR_Name("foo"), "baz", "bar"), false)

	n = QualifiedName_Of("this").Getprop("qux")
	assertMatches(t, n, qname(IR_ThisNode(), "qux"), true)
	assertMatches(t, n, qname(IR_Name("x"), "qux"), false)

	n = QualifiedName_Of("super").Getprop("qux")
	assertMatches(t, n, qname(IR_SuperNode(), "qux"), true)
	assertMatches(t, n, qname(IR_Name("x"), "qux"), false)
}

func TestQualifiedName_EdgeCases_FromString(t *testing.T) {
	assertComponents(t, QualifiedName_Of(""), "")
	assertComponents(t, QualifiedName_Of("."), "", "")
	assertComponents(t, QualifiedName_Of(".foo.bar"), "", "foo", "bar")
	assertComponents(t, QualifiedName_Of("foo.bar."), "foo", "bar", "")
	assertComponents(t, QualifiedName_Of("foo..bar"), "foo", "", "bar")
}

func TestQualifiedName_GetComponentCount_FromString(t *testing.T) {
	assertComponentCount(t, QualifiedName_Of("foo"), 1)
	assertComponentCount(t, QualifiedName_Of("foo.bar"), 2)
	assertComponentCount(t, QualifiedName_Of("foo.bar.baz"), 3)
}

func TestQualifiedName_GetComponentCount_FromNode(t *testing.T) {
	assertComponentCount(t, IR_Name("foo").GetQualifiedNameObject(), 1)
	assertComponentCount(t, qname(IR_Name("foo"), "bar").GetQualifiedNameObject(), 2)
	assertComponentCount(t, qname(IR_Name("foo"), "bar", "baz").GetQualifiedNameObject(), 3)
}

func TestQualifiedName_GetComponentCount_FromGetprop(t *testing.T) {
	assertComponentCount(t, QualifiedName_Of("foo").Getprop("bar"), 2)
}

// Go-only: QualifiedName and its subclasses do not override toString, so a
// message that embeds one (Es6RewriteModules' "unexpected qualified name:
// %s") prints java.lang.Object.toString(): the binary class name, '@' and
// a hex hash. The port printed Go's %v of the struct.
func TestQualifiedName_ObjectToString(t *testing.T) {
	check := func(q QualifiedName, class string) {
		t.Helper()
		s := jutil.ToJavaString(q)
		prefix := "com.google.javascript.rhino.QualifiedName$" + class + "@"
		hash, ok := strings.CutPrefix(s, prefix)
		if !ok || hash == "" || strings.Trim(hash, "0123456789abcdef") != "" {
			t.Errorf("toString() = %q, want %s<hex>", s, prefix)
		}
	}
	check(QualifiedName_Of("a.b"), "StringListQname")
	check(QualifiedName_Of("a").Getprop("b"), "GetpropQname")
	check(qname(IR_Name("a"), "b").GetQualifiedNameObject(), "NodeQname")
}

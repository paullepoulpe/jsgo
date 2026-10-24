// Port of com/google/javascript/rhino/RhinoStringPoolTest.java @ closure-compiler v20261021

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
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"strconv"
	"testing"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// note - RhinoStringPool intentionally uses static state, so these tests are
// not really order-independent.

// sameInstance is Truth's isSameInstanceAs for strings.
//
// Go-only: language-adapter — Java compares String references; the Go
// equivalent of an interned instance is a string sharing the canonical
// backing array, so the data pointers are compared.
func sameInstance(a, b string) bool {
	return len(a) == len(b) && (len(a) == 0 || unsafe.StringData(a) == unsafe.StringData(b))
}

// newString is `new String(s)`: a fresh copy with its own backing array.
func newString(s string) string {
	return string([]byte(s))
}

func TestRhinoStringPool_ReturnsDotEqualsString(t *testing.T) {
	newFoo := RhinoStringPool_AddOrGet("foo")
	if newFoo != "foo" {
		t.Fatalf("got %q", newFoo)
	}
}

func TestRhinoStringPool_MultipleCallsReturnSameInstance(t *testing.T) {
	var foos [100]string
	for i := 0; i < 100; i++ {
		newFoo := newString("foo") // use a copy to ensure different objects
		foos[i] = RhinoStringPool_AddOrGet(newFoo)
	}
	for i := 1; i < 100; i++ {
		if !sameInstance(foos[i], foos[0]) {
			t.Fatalf("foos[%d] is not the same instance as foos[0]", i)
		}
	}
}

func TestRhinoStringPool_LazyInternedStringList_OutOfBounds_ThrowsException(t *testing.T) {
	list := NewRhinoStringPool_LazyInternedStringList([]string{"foo"})
	assertThrows[jutil.IArrayIndexOutOfBoundsException](t, func() { list.Get(100) })
}

func TestRhinoStringPool_LazyInternedStringList_NegativeIndex_ThrowsException(t *testing.T) {
	list := NewRhinoStringPool_LazyInternedStringList([]string{})
	assertThrows[jutil.IArrayIndexOutOfBoundsException](t, func() { list.Get(-1) })
}

func TestRhinoStringPool_LazyInternedStringList_GetsEqualValuesAsInputList(t *testing.T) {
	list := NewRhinoStringPool_LazyInternedStringList([]string{"foo", "bar", "baz"})
	for i, want := range []string{"foo", "bar", "baz"} {
		if got := list.Get(i); got != want {
			t.Errorf("list.Get(%d) = %q, want %q", i, got, want)
		}
	}
}

func TestRhinoStringPool_LazyInternedStringList_GetsSameInstanceAsRhinoStringPool(t *testing.T) {
	// use copies to ensure we're passing unique String instances
	strs := []string{newString("foo"), newString("bar"), newString("baz")}
	list := NewRhinoStringPool_LazyInternedStringList(strs)
	for i, s := range []string{"foo", "bar", "baz"} {
		if !sameInstance(list.Get(i), RhinoStringPool_AddOrGet(s)) {
			t.Errorf("list.Get(%d) is not the pooled instance", i)
		}
	}
}

func TestRhinoStringPool_LazyInternedStringList_VeryLargeListWorks(t *testing.T) {
	var strs []string
	for i := 0; i < 10000; i++ {
		strs = append(strs, newString("foo"+strconv.Itoa(i)))
	}
	list := NewRhinoStringPool_LazyInternedStringList(strs)
	for i := 0; i < 10000; i++ {
		if !sameInstance(list.Get(i), RhinoStringPool_AddOrGet("foo"+strconv.Itoa(i))) {
			t.Fatalf("list.Get(%d) is not the pooled instance", i)
		}
	}
}

func TestRhinoStringPool_WriteOnlyBitset_InitializedToFalse(t *testing.T) {
	bitset := NewRhinoStringPool_WriteOnlyBitset(5)
	for i := 0; i < 5; i++ {
		if bitset.Get(i) {
			t.Errorf("bit %d set", i)
		}
	}
}

func TestRhinoStringPool_WriteOnlyBitset_SetToTrue_ThenGetReturnsTrue(t *testing.T) {
	bitset := NewRhinoStringPool_WriteOnlyBitset(2)
	bitset.Set(0)
	if !bitset.Get(0) {
		t.Error("bit 0 not set")
	}
	if bitset.Get(1) {
		t.Error("bit 1 set")
	}
}

func TestRhinoStringPool_WriteOnlyBitset_WordBoundary_SetToTrue_ThenGetReturnsTrue(t *testing.T) {
	bitset := NewRhinoStringPool_WriteOnlyBitset(33)
	bitset.Set(32)
	if !bitset.Get(32) {
		t.Error("bit 32 not set")
	}
}

func TestRhinoStringPool_WriteOnlyBitset_MultipleSetsToTrue_AllReturnTrue(t *testing.T) {
	bitset := NewRhinoStringPool_WriteOnlyBitset(5)
	bitset.Set(1)
	bitset.Set(2)
	bitset.Set(3)
	for i, want := range []bool{false, true, true, true, false} {
		if bitset.Get(i) != want {
			t.Errorf("bit %d = %v, want %v", i, bitset.Get(i), want)
		}
	}
}

// Port of com/google/javascript/rhino/HamtPMapTest.java @ closure-compiler v20260415

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

// Java's boxed Integer and String values (which may be null) are carried as
// `any` here, so that null is nil and equality is Objects.equals
// (jutil.ObjectsEquals), as in the Java test.

import (
	"iter"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// assertValuesExactly is Truth's containsExactly (any order, with
// multiplicity) on an iterable, using Objects.equals.
func assertValuesExactly[T any](t *testing.T, seq iter.Seq[T], want ...any) {
	t.Helper()
	var got []any
	for v := range seq {
		got = append(got, any(v))
	}
	remaining := slices.Clone(want)
	ok := len(got) == len(want)
	for _, g := range got {
		if !ok {
			break
		}
		i := slices.IndexFunc(remaining, func(w any) bool { return jutil.ObjectsEquals(g, w) })
		if i < 0 {
			ok = false
			break
		}
		remaining = slices.Delete(remaining, i, i+1)
	}
	if !ok {
		t.Errorf("got %v, want exactly %v", got, want)
	}
}

func assertGet[K any](t *testing.T, m PMap[K, any], key K, want any) {
	t.Helper()
	if got := m.Get(key); !jutil.ObjectsEquals(got, want) {
		t.Errorf("get(%v) = %v, want %v", key, got, want)
	}
}

func assertIsEmpty[K any](t *testing.T, m PMap[K, any], want bool) {
	t.Helper()
	if got := m.IsEmpty(); got != want {
		t.Errorf("isEmpty() = %v, want %v", got, want)
	}
}

func assertBool(t *testing.T, what string, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", what, got, want)
	}
}

func TestHamtPMap_Empty(t *testing.T) {
	var m PMap[string, any] = HamtPMap_Empty[string, any]()
	assertGet(t, m, "foo", nil)
	assertValuesExactly(t, m.Values())
	assertIsEmpty(t, m, true)
}

func TestHamtPMap_Plus(t *testing.T) {
	var empty PMap[string, any] = HamtPMap_Empty[string, any]()
	m := empty.Plus("foo", 42)
	assertGet(t, m, "foo", 42)
	assertValuesExactly(t, m.Values(), 42)
	assertIsEmpty(t, m, false)
}

func TestHamtPMap_Plus_AlreadyExistsReturnsSame(t *testing.T) {
	var empty PMap[string, any] = HamtPMap_Empty[string, any]()
	m := empty.Plus("foo", 42)
	if m.Plus("foo", 42) != m {
		t.Error("plus of an existing entry is not the same instance")
	}
}

func TestHamtPMap_Plus_UpdatesExistingKey(t *testing.T) {
	var empty PMap[string, any] = HamtPMap_Empty[string, any]()
	m := empty.Plus("foo", 42).Plus("foo", 23)
	assertGet(t, m, "foo", 23)
	assertValuesExactly(t, m.Values(), 23)
	assertIsEmpty(t, m, false)
}

func TestHamtPMap_Minus(t *testing.T) {
	m := HamtPMap_Empty[string, any]().Plus("foo", 42)
	empty := m.Minus("foo")
	assertGet(t, empty, "foo", nil)
	assertValuesExactly(t, empty.Values())
	assertIsEmpty(t, empty, true)
}

func TestHamtPMap_Minus_NonexistentKeyReturnsSame(t *testing.T) {
	var empty PMap[string, any] = HamtPMap_Empty[string, any]()
	m := empty.Plus("foo", 42)
	if m.Minus("bar") != m {
		t.Error("minus of a missing key is not the same instance")
	}
}

func TestHamtPMap_PlusAndMinus(t *testing.T) {
	var empty PMap[string, any] = HamtPMap_Empty[string, any]()
	m := empty.Plus("foo", 42).Plus("bar", 19).Plus("baz", 12)
	assertGet(t, m, "foo", 42)
	assertGet(t, m, "bar", 19)
	assertGet(t, m, "baz", 12)
	assertValuesExactly(t, m.Values(), 12, 19, 42)
	assertIsEmpty(t, m, false)

	m = m.Minus("bar")
	assertGet(t, m, "foo", 42)
	assertGet(t, m, "bar", nil)
	assertGet(t, m, "baz", 12)
	assertValuesExactly(t, m.Values(), 12, 42)
	assertIsEmpty(t, m, false)

	m = m.Minus("baz")
	assertGet(t, m, "foo", 42)
	assertGet(t, m, "baz", nil)
	assertValuesExactly(t, m.Values(), 42)
	assertIsEmpty(t, m, false)

	m = m.Minus("foo")
	assertGet(t, m, "foo", nil)
	assertValuesExactly(t, m.Values())
	assertIsEmpty(t, m, true)
}

func TestHamtPMap_Reconcile_CustomJoiner(t *testing.T) {
	var empty PMap[int, any] = HamtPMap_Empty[int, any]()
	left := empty.Plus(1, "1").Plus(3, "3").Plus(5, "5").Plus(7, "7").Plus(9, "9")
	right := empty.Plus(3, "C").Plus(4, "D").Plus(5, "5").Plus(6, "F").Plus(7, "G")

	expectations := newJoinExpectations(t).
		expect(1, "1", nil, "1").
		expect(3, "3", "C", "c").
		expect(4, nil, "D", "D").
		expect(6, nil, "F", "F").
		expect(7, "7", "G", "g").
		expect(9, "9", nil, "9")

	joined := left.Reconcile(right, expectations.merge)

	assertValuesExactly(t, joined.Values(), "1", "c", "D", "5", "F", "g", "9")
	expectations.verify()
}

func TestHamtPMap_Reconcile_EmptyMaps(t *testing.T) {
	var empty PMap[int, any] = HamtPMap_Empty[int, any]()
	m := empty.Plus(1, "1").Plus(3, "3").Plus(5, "5").Plus(7, "7").Plus(9, "9")

	joined := m.Reconcile(
		empty,
		func(k int, a, b any) any {
			if b != nil {
				t.Errorf("b = %v, want null", b)
			}
			return a
		})
	assertBool(t, "joined.equivalent(map)", joined.Equivalent(m, objectsEquals), true)

	joined = empty.Reconcile(
		m,
		func(k int, a, b any) any {
			if a != nil {
				t.Errorf("a = %v, want null", a)
			}
			return b
		})
	assertBool(t, "joined.equivalent(map)", joined.Equivalent(m, objectsEquals), true)
}

func TestHamtPMap_Reconcile_DifferentSizes(t *testing.T) {
	left := buildHamt(t, 1, 6, 2, 19, 4, 23, 5, 8, 42, 12, 18, 33)
	right := buildHamt(t).Plus(3, 1).Plus(4, 1).Plus(8, 1).Plus(19, 1).Plus(25, 1).Plus(42, 1)
	expected := left.Plus(3, 1).Plus(4, 5).Plus(8, 9).Plus(19, 20).Plus(25, 1).Plus(42, 43)

	joined := left.Reconcile(right, func(k int, a, b any) any {
		if a == nil {
			return b
		}
		if b == nil {
			return a
		}
		return a.(int) + b.(int)
	})
	assertBool(t, "joined.equivalent(expected)", joined.Equivalent(expected, objectsEquals), true)
}

func TestHamtPMap_Reconcile_RejectsNullResult(t *testing.T) {
	left := buildHamt(t, 1, 3, 5, 7)
	right := buildHamt(t, 2, 4, 6, 8)

	assertThrows[jutil.Throwable](t, func() {
		left.Reconcile(right, func(k int, a, b any) any { return nil })
	})
}

func TestHamtPMap_Equivalent_Equivalent(t *testing.T) {
	left := buildHamt(t, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024)
	right := buildHamt(t, 2, 512, 32, 8, 1024, 1, 64, 4, 16, 256, 128)

	assertBool(t, "left.equivalent(right)", left.Equivalent(right, objectsEquals), true)
	assertBool(t, "right.equivalent(left)", right.Equivalent(left, objectsEquals), true)
}

func TestHamtPMap_Equivalent_Equivalent_HashCodeCollisions(t *testing.T) {
	fixedHashCode := int32(42)
	left := buildHamtWithCustomHashCode(t, fixedHashCode, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024)
	right := buildHamtWithCustomHashCode(t, fixedHashCode, 2, 512, 32, 8, 1024, 1, 64, 4, 16, 256, 128)

	assertBool(t, "left.equivalent(right)", left.Equivalent(right, objectsEquals), true)
	assertBool(t, "right.equivalent(left)", right.Equivalent(left, objectsEquals), true)
}

func TestHamtPMap_Equivalent_OneDifferentValue(t *testing.T) {
	// Insertion order doesn't matter
	left1 := buildHamt(t, 1, 3, 5, 7)
	left2 := buildHamt(t, 7, 3, 5, 1)
	right1 := buildHamt(t).Plus(1, 2).Plus(3, 3).Plus(5, 5).Plus(7, 7)
	right2 := buildHamt(t, 1, 3, 5).Plus(7, 8)

	assertBool(t, "left1.equivalent(right1)", left1.Equivalent(right1, objectsEquals), false)
	assertBool(t, "right1.equivalent(left1)", right1.Equivalent(left1, objectsEquals), false)
	assertBool(t, "left2.equivalent(right1)", left2.Equivalent(right1, objectsEquals), false)
	assertBool(t, "right1.equivalent(left2)", right1.Equivalent(left2, objectsEquals), false)
	assertBool(t, "left1.equivalent(right2)", left1.Equivalent(right2, objectsEquals), false)
	assertBool(t, "right2.equivalent(left1)", right2.Equivalent(left1, objectsEquals), false)
	assertBool(t, "left2.equivalent(right2)", left2.Equivalent(right2, objectsEquals), false)
	assertBool(t, "right2.equivalent(left2)", right2.Equivalent(left2, objectsEquals), false)
}

func TestHamtPMap_Equivalent_OneDifferentKey(t *testing.T) {
	left1 := buildHamt(t).Plus(2, 1).Plus(3, 3).Plus(5, 5).Plus(7, 7)
	left2 := buildHamt(t, 1, 3, 5).Plus(7, 8)
	right := buildHamt(t, 1, 3, 5, 7)

	assertBool(t, "left1.equivalent(right)", left1.Equivalent(right, objectsEquals), false)
	assertBool(t, "right.equivalent(left1)", right.Equivalent(left1, objectsEquals), false)
	assertBool(t, "left2.equivalent(right)", left2.Equivalent(right, objectsEquals), false)
	assertBool(t, "right.equivalent(left2)", right.Equivalent(left2, objectsEquals), false)
}

func TestHamtPMap_Equivalent_OneDifferentKey_WithHashCodeCollisions(t *testing.T) {
	fixedHashCode := int32(42)
	left1 := buildHamtWithCustomHashCode(t, 42).
		Plus(newIntegerWithCustomHashCode(2, fixedHashCode), 1).
		Plus(newIntegerWithCustomHashCode(3, fixedHashCode), 3).
		Plus(newIntegerWithCustomHashCode(5, fixedHashCode), 5).
		Plus(newIntegerWithCustomHashCode(7, fixedHashCode), 7)
	left2 := buildHamtWithCustomHashCode(t, fixedHashCode, 1, 3, 5).
		Plus(newIntegerWithCustomHashCode(7, fixedHashCode), 8)
	right := buildHamtWithCustomHashCode(t, fixedHashCode, 1, 3, 5, 7)

	assertBool(t, "left1.equivalent(right)", left1.Equivalent(right, objectsEquals), false)
	assertBool(t, "right.equivalent(left1)", right.Equivalent(left1, objectsEquals), false)
	assertBool(t, "left2.equivalent(right)", left2.Equivalent(right, objectsEquals), false)
	assertBool(t, "right.equivalent(left2)", right.Equivalent(left2, objectsEquals), false)
}

func TestHamtPMap_Equivalent_OneDifferentKey_WithHashCodeCollisions_DoesNotPassNullToEquivalenceTest(t *testing.T) {
	// Regression test for b/323185572, where HamtPMap would sometimes pass a
	// `null` value to the 'BiPredicate<V, V> equivalence' parameter of
	// HamtPMap.equivalent. `null` values are an internal implementation
	// detail of HamtPMap.pivot and not part of the public API.
	fixedHashCode := int32(42)
	left1 := buildHamtWithCustomHashCode(t, 42).
		Plus(newIntegerWithCustomHashCode(2, fixedHashCode), 1).
		Plus(newIntegerWithCustomHashCode(3, fixedHashCode), 3).
		Plus(newIntegerWithCustomHashCode(5, fixedHashCode), 5).
		Plus(newIntegerWithCustomHashCode(7, fixedHashCode), 7)
	left2 := buildHamtWithCustomHashCode(t, fixedHashCode, 1, 3, 5).
		Plus(newIntegerWithCustomHashCode(7, fixedHashCode), 8)
	right := buildHamtWithCustomHashCode(t, fixedHashCode, 1, 3, 5, 7)

	assertBool(t, "left1.equivalent(right)", left1.Equivalent(right, equalsAndCheckNotNull), false)
	assertBool(t, "right.equivalent(left1)", right.Equivalent(left1, equalsAndCheckNotNull), false)
	assertBool(t, "left2.equivalent(right)", left2.Equivalent(right, equalsAndCheckNotNull), false)
	assertBool(t, "right.equivalent(left2)", right.Equivalent(left2, equalsAndCheckNotNull), false)
}

func TestHamtPMap_Equivalent_OneMissingKey(t *testing.T) {
	left := buildHamt(t, 1, 5, 4, 3, 9)
	right := buildHamt(t, 1, 5, 3, 9)

	assertBool(t, "left.equivalent(right)", left.Equivalent(right, objectsEquals), false)
	assertBool(t, "right.equivalent(left)", right.Equivalent(left, objectsEquals), false)
}

func TestHamtPMap_Equivalent_EqualsWithCustomEquivalence(t *testing.T) {
	var empty PMap[string, any] = HamtPMap_Empty[string, any]()
	left := empty.Plus("abc", "AAa").Plus("def", "Ddd").Plus("ghi", "ggG")
	right := empty.Plus("abc", "aAA").Plus("def", "ddD").Plus("ghi", "GGG")

	assertBool(t, "left.equivalent(right)", left.Equivalent(right, func(l, r any) bool {
		return strings.ToLower(l.(string)) == strings.ToLower(r.(string))
	}), true)
}

func TestHamtPMap_Equivalent_NotEquivalentShortCircuits(t *testing.T) {
	var empty PMap[int, any] = HamtPMap_Empty[int, any]()
	left := empty.Plus(1, "1").Plus(2, "2")
	right := empty.Plus(1, "3").Plus(2, "4")

	calls := 0
	assertBool(t, "left.equivalent(right)", left.Equivalent(right, func(l, r any) bool {
		calls++
		return false
	}), false)
	if calls != 1 {
		t.Errorf("calls = %d, want 1", calls)
	}
}

func TestHamtPMap_Keys_ContainsExactlyKeys(t *testing.T) {
	var empty PMap[int, any] = HamtPMap_Empty[int, any]()
	actual := empty.Plus(1, "a").Plus(2, "b").Plus(100, "x").Plus(-847, "h").Minus(2)

	assertValuesExactly(t, actual.Keys(), 1, 100, -847)
}

func TestHamtPMap_Values_ContainsExactlyValues(t *testing.T) {
	var empty PMap[int, any] = HamtPMap_Empty[int, any]()
	actual := empty.Plus(1, "a").Plus(2, "b").Plus(3, "kkdmw").Plus(4, ":::").Minus(1)

	assertValuesExactly(t, actual.Values(), "b", "kkdmw", ":::")
}

// joinExpectations is HamtPMapTest.JoinExpectations, a Reconciler (its
// merge method value is the PMap_Reconciler).
type joinExpectations struct {
	t            *testing.T
	expectations *jutil.LinkedHashMap[int, []any]
}

func newJoinExpectations(t *testing.T) *joinExpectations {
	return &joinExpectations{t: t, expectations: jutil.NewLinkedHashMap[int, []any]()}
}

func (j *joinExpectations) merge(key int, left any, right any) any {
	j.t.Helper()
	if !j.expectations.ContainsKey(key) {
		j.t.Fatalf("unexpected merge of key %d", key)
	}
	e := j.expectations.Get(key)
	if !jutil.ObjectsEquals(e[0], left) {
		j.t.Errorf("merge(%d): left = %v, want %v", key, left, e[0])
	}
	if !jutil.ObjectsEquals(e[1], right) {
		j.t.Errorf("merge(%d): right = %v, want %v", key, right, e[1])
	}
	j.expectations.Remove(key)
	return e[2]
}

func (j *joinExpectations) expect(key int, left, right, result any) *joinExpectations {
	if j.expectations.ContainsKey(key) {
		j.t.Fatalf("duplicate expectation for key %d", key)
	}
	j.expectations.Put(key, []any{left, right, result})
	return j
}

func (j *joinExpectations) verify() {
	if !j.expectations.IsEmpty() {
		j.t.Errorf("unmet expectations for keys %v", slices.Collect(j.expectations.Keys()))
	}
}

func TestHamtPMap_Integration(t *testing.T) {
	ref := map[int]bool{} // TreeSet<Integer>
	refValues := func() []any {
		var vs []any
		for j := 1; j < 127; j++ {
			if ref[j] {
				vs = append(vs, strconv.Itoa(j))
			}
		}
		return vs
	}
	var m PMap[int, any] = HamtPMap_Empty[int, any]()

	check := func() {
		t.Helper()
		m.(*HamtPMap[int, any]).AssertCorrectStructure()
		assertIsEmpty(t, m, len(ref) == 0)
		for j := 1; j < 127; j++ {
			var want any
			if ref[j] {
				want = strconv.Itoa(j)
			}
			assertGet(t, m, j, want)
			if ref[j] {
				if m.Plus(j, strconv.Itoa(j)) != m {
					t.Fatalf("plus(%d) of an existing entry is not the same instance", j)
				}
			} else {
				if m.Minus(j) != m {
					t.Fatalf("minus(%d) of a missing key is not the same instance", j)
				}
			}
		}
		assertValuesExactly(t, m.Values(), refValues()...)
	}

	// Note: multiplying by 43 produces every nonzero integer mod 127.
	for i := 17; !ref[i]; i = (i * 43) % 127 {
		ref[i] = true
		m = m.Plus(i, strconv.Itoa(i))
		check()
	}
	if len(ref) != 126 { // (quick check that we actually did add everything)
		t.Fatalf("ref has %d elements, want 126", len(ref))
	}

	// Note: multiplying by 39 produces every nonzero integer mod 127 in a
	// different order.
	for i := 12; ref[i]; i = (i * 39) % 127 {
		delete(ref, i)
		m = m.Minus(i)
		check()
	}
	if len(ref) != 0 {
		t.Fatalf("ref has %d elements, want 0", len(ref))
	}
	if m != PMap[int, any](HamtPMap_Empty[int, any]()) {
		t.Error("map is not the same instance as HamtPMap.empty()")
	}
}

func buildHamt(t *testing.T, values ...int) PMap[int, any] {
	t.Helper()
	var m PMap[int, any] = HamtPMap_Empty[int, any]()
	for _, value := range values {
		m = m.Plus(value, value)
		m.(*HamtPMap[int, any]).AssertCorrectStructure()
	}
	return m
}

func buildHamtWithCustomHashCode(t *testing.T, fixedHashCode int32, values ...int) PMap[*integerWithCustomHashCode, any] {
	t.Helper()
	var m PMap[*integerWithCustomHashCode, any] = HamtPMap_Empty[*integerWithCustomHashCode, any]()
	for _, value := range values {
		m = m.Plus(newIntegerWithCustomHashCode(value, fixedHashCode), value)
		m.(*HamtPMap[*integerWithCustomHashCode, any]).AssertCorrectStructure()
	}
	return m
}

// integerWithCustomHashCode is a helper class to test HamtPMap behavior
// given keys with hash code collisions.
type integerWithCustomHashCode struct {
	i        int
	hashCode int32
}

func newIntegerWithCustomHashCode(i int, hashCode int32) *integerWithCustomHashCode {
	return &integerWithCustomHashCode{i: i, hashCode: hashCode}
}

func (x *integerWithCustomHashCode) HashCode() int32 { return x.hashCode }

func (x *integerWithCustomHashCode) Equals(o any) bool {
	if that, ok := o.(*integerWithCustomHashCode); ok && that != nil {
		return x.i == that.i
	}
	return false
}

// objectsEquals is Objects::equals as a BiPredicate.
func objectsEquals(a, b any) bool { return jutil.ObjectsEquals(a, b) }

func equalsAndCheckNotNull(o1, o2 any) bool {
	jutil.CheckNotNull(o1)
	jutil.CheckNotNull(o2)
	return jutil.ObjectsEquals(o1, o2)
}

// Port of com/google/javascript/rhino/HamtPMapTest.java @ closure-compiler v20261028

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
	"github.com/paullepoulpe/jsgo/internal/strpool"
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

func assertGet[K comparable](t *testing.T, m *HamtPMap[K, any], key K, want any) {
	t.Helper()
	if got := m.Get(key); !jutil.ObjectsEquals(got, want) {
		t.Errorf("get(%v) = %v, want %v", key, got, want)
	}
}

func assertIsEmpty[K comparable](t *testing.T, m *HamtPMap[K, any], want bool) {
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
	m := HamtPMap_Empty[string, any]()
	assertGet(t, m, "foo", nil)
	assertValuesExactly(t, m.Values())
	assertIsEmpty(t, m, true)

	var nilMap *HamtPMap[string, any]
	if nilMap != HamtPMap_Empty[string, any]() {
		t.Error("HamtPMap_Empty() != nil")
	}
	assertGet(t, nilMap, "foo", nil)
	assertValuesExactly(t, nilMap.Values())
	assertIsEmpty(t, nilMap, true)
	if got := nilMap.String(); got != "{}" {
		t.Errorf("nilMap.String() = %q, want {}", got)
	}
}

func TestHamtPMap_Plus(t *testing.T) {
	empty := HamtPMap_Empty[string, any]()
	m := empty.Plus("foo", 42)
	assertGet(t, m, "foo", 42)
	assertValuesExactly(t, m.Values(), 42)
	assertIsEmpty(t, m, false)
}

func TestHamtPMap_Plus_AlreadyExistsReturnsSame(t *testing.T) {
	empty := HamtPMap_Empty[string, any]()
	m := empty.Plus("foo", 42)
	if m.Plus("foo", 42) != m {
		t.Error("plus of an existing entry is not the same instance")
	}
}

func TestHamtPMap_Plus_UpdatesExistingKey(t *testing.T) {
	empty := HamtPMap_Empty[string, any]()
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
	empty := HamtPMap_Empty[string, any]()
	m := empty.Plus("foo", 42)
	if m.Minus("bar") != m {
		t.Error("minus of a missing key is not the same instance")
	}
}

func TestHamtPMap_PlusAndMinus(t *testing.T) {
	empty := HamtPMap_Empty[string, any]()
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
	empty := HamtPMap_Empty[int, any]()
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
	empty := HamtPMap_Empty[int, any]()
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
	empty := HamtPMap_Empty[string, any]()
	left := empty.Plus("abc", "AAa").Plus("def", "Ddd").Plus("ghi", "ggG")
	right := empty.Plus("abc", "aAA").Plus("def", "ddD").Plus("ghi", "GGG")

	assertBool(t, "left.equivalent(right)", left.Equivalent(right, func(l, r any) bool {
		return strings.ToLower(l.(string)) == strings.ToLower(r.(string))
	}), true)
}

func TestHamtPMap_Equivalent_NotEquivalentShortCircuits(t *testing.T) {
	empty := HamtPMap_Empty[int, any]()
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
	empty := HamtPMap_Empty[int, any]()
	actual := empty.Plus(1, "a").Plus(2, "b").Plus(100, "x").Plus(-847, "h").Minus(2)

	assertValuesExactly(t, actual.Keys(), 1, 100, -847)
}

func TestHamtPMap_Values_ContainsExactlyValues(t *testing.T) {
	empty := HamtPMap_Empty[int, any]()
	actual := empty.Plus(1, "a").Plus(2, "b").Plus(3, "kkdmw").Plus(4, ":::").Minus(1)

	assertValuesExactly(t, actual.Values(), "b", "kkdmw", ":::")
}

func TestHamtPMap_All_OrderAndEarlyStop(t *testing.T) {
	var nilMap *HamtPMap[int, any]
	for range nilMap.All() {
		t.Fatal("nil map should yield no entries")
	}

	m := HamtPMap_Empty[int, any]()
	for i := 1; i <= 100; i++ {
		m = m.Plus(i*17, strconv.Itoa(i))
	}

	var wantKeys []int
	var wantVals []any
	for q := []*HamtPMap[int, any]{m}; len(q) > 0; {
		top := q[0]
		q = q[1:]
		wantKeys = append(wantKeys, top.key)
		wantVals = append(wantVals, top.value)
		for i := len(top.children) - 1; i >= 0; i-- {
			q = append(q, top.children[i])
		}
	}

	var gotKeys []int
	var gotVals []any
	for k, v := range m.All() {
		gotKeys = append(gotKeys, k)
		gotVals = append(gotVals, v)
	}
	if !slices.Equal(gotKeys, wantKeys) {
		t.Fatalf("All() keys = %v, want %v", gotKeys, wantKeys)
	}
	if !slices.Equal(gotVals, wantVals) {
		t.Fatalf("All() vals = %v, want %v", gotVals, wantVals)
	}
	if !slices.Equal(slices.Collect(m.Keys()), wantKeys) {
		t.Fatalf("Keys() order mismatch")
	}
	if !slices.Equal(slices.Collect(m.Values()), wantVals) {
		t.Fatalf("Values() order mismatch")
	}

	count := 0
	for range m.All() {
		count++
		if count == 5 {
			break
		}
	}
	if count != 5 {
		t.Fatalf("early stop count = %d, want 5", count)
	}
}

// joinExpectations is HamtPMapTest.JoinExpectations, a Reconciler (its
// merge method value is the HamtPMap_Reconciler).
type joinExpectations struct {
	t            *testing.T
	expectations *jutil.OrderedMap[int, []any]
}

func newJoinExpectations(t *testing.T) *joinExpectations {
	return &joinExpectations{t: t, expectations: jutil.NewOrderedMap[int, []any]()}
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
	m := HamtPMap_Empty[int, any]()

	check := func() {
		t.Helper()
		m.AssertCorrectStructure()
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
	if m != HamtPMap_Empty[int, any]() {
		t.Error("map is not the same instance as HamtPMap.empty()")
	}
}

func buildHamt(t *testing.T, values ...int) *HamtPMap[int, any] {
	t.Helper()
	m := HamtPMap_Empty[int, any]()
	for _, value := range values {
		m = m.Plus(value, value)
		m.AssertCorrectStructure()
	}
	return m
}

func buildHamtWithCustomHashCode(t *testing.T, fixedHashCode int32, values ...int) *HamtPMap[integerWithCustomHashCode, any] {
	t.Helper()
	m := HamtPMap_Empty[integerWithCustomHashCode, any]()
	for _, value := range values {
		m = m.Plus(newIntegerWithCustomHashCode(value, fixedHashCode), value)
		m.AssertCorrectStructure()
	}
	return m
}

// integerWithCustomHashCode is a helper class to test HamtPMap behavior
// given keys with hash code collisions.
type integerWithCustomHashCode struct {
	i        int
	hashCode int32
}

func newIntegerWithCustomHashCode(i int, hashCode int32) integerWithCustomHashCode {
	return integerWithCustomHashCode{i: i, hashCode: hashCode}
}

func (x integerWithCustomHashCode) HashCode() int32 { return x.hashCode }

// objectsEquals is Objects::equals as a BiPredicate.
func objectsEquals(a, b any) bool { return jutil.ObjectsEquals(a, b) }

func equalsAndCheckNotNull(o1, o2 any) bool {
	jutil.CheckNotNull(o1)
	jutil.CheckNotNull(o2)
	return jutil.ObjectsEquals(o1, o2)
}

type dummyHamtJSType struct {
	JSType
	tag int
}

func (*dummyHamtJSType) GetJSTypeBase() *JSTypeBase { return nil }
func (d *dummyHamtJSType) Equals(o any) bool {
	other, ok := o.(*dummyHamtJSType)
	return ok && other != nil && d.tag == other.tag
}

func TestHamtPMapArenaZeroAllocs(t *testing.T) {
	var strArena HamtPMapArena[string, JSType]
	var idArena HamtPMapArena[strpool.ID, JSType]
	keys := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", "theta"}
	var t1 JSType = &dummyHamtJSType{tag: 1}
	var t2 JSType = &dummyHamtJSType{tag: 2}
	joinerStr := func(_ string, v1, v2 JSType) JSType {
		if v1 != nil {
			return v1
		}
		return v2
	}
	joinerID := func(_ strpool.ID, v1, v2 JSType) JSType {
		if v1 != nil {
			return v1
		}
		return v2
	}

	allocs := testing.AllocsPerRun(20, func() {
		strArena.Reset()
		var m1, m2 *HamtPMap[string, JSType]
		for i, k := range keys {
			if i%2 == 0 {
				m1 = m1.PlusMapWithArena(&strArena, k, t1)
			} else {
				m2 = m2.PlusMapWithArena(&strArena, k, t2)
			}
		}
		merged := m1.ReconcileMapWithArena(&strArena, m2, joinerStr)
		if merged.Get("alpha") != t1 || merged.Get("beta") != t2 {
			t.Fatal("unexpected merged values")
		}
		_ = merged.MinusMapWithArena(&strArena, "alpha")

		idArena.Reset()
		var im1, im2 *HamtPMap[strpool.ID, JSType]
		for i := 1; i <= 8; i++ {
			id := strpool.ID(i)
			if i%2 == 0 {
				im1 = im1.PlusMapWithArena(&idArena, id, t1)
			} else {
				im2 = im2.PlusMapWithArena(&idArena, id, t2)
			}
		}
		imerged := im1.ReconcileMapWithArena(&idArena, im2, joinerID)
		if imerged.Get(2) != t1 || imerged.Get(1) != t2 {
			t.Fatal("unexpected id merged values")
		}
	})
	if allocs != 0 {
		t.Fatalf("got %v allocs/op, want 0", allocs)
	}
	strArena.TrimLargeChunks()
}

func TestHamtPMapArenaTrimLargeChunksRetainsWorkingSet(t *testing.T) {
	reg := NewJSTypeRegistry(nil)
	t1 := reg.GetNativeType(JSTypeNative_NUMBER_TYPE)

	var arena HamtPMapArena[strpool.ID, JSType]
	// Grow past jutil.ReuseSlab's 15-chunk cap (>3,000 nodes) to simulate a
	// large global pass.
	var big *HamtPMap[strpool.ID, JSType]
	for i := 1; i <= 3500; i++ {
		big = big.PlusMapWithArena(&arena, strpool.ID(i), t1)
	}
	if big.Get(3500) != t1 {
		t.Fatal("expected big map to contain 3500")
	}

	// TrimLargeChunks must drop excess chunks while retaining the working set
	// so subsequent function-level passes run with zero allocations.
	arena.TrimLargeChunks()

	allocs := testing.AllocsPerRun(50, func() {
		arena.Reset()
		var m *HamtPMap[strpool.ID, JSType]
		for i := 1; i <= 16; i++ {
			m = m.PlusMapWithArena(&arena, strpool.ID(i), t1)
		}
		if m.Get(16) != t1 {
			t.Fatal("expected working-set map to contain 16")
		}
	})
	if allocs != 0 {
		t.Fatalf("after TrimLargeChunks, got %v allocs/op, want 0", allocs)
	}
}

func TestHamtPMapNullAndEqualitySemantics(t *testing.T) {
	if hamtIsNull("") {
		t.Fatal("empty string must not be treated as null")
	}
	if hamtIsNull(0) {
		t.Fatal("int(0) must not be treated as null")
	}
	if !hamtIsNull((JSType)(nil)) {
		t.Fatal("nil JSType must be treated as null")
	}
	if !hamtIsNull((*TemplateType)(nil)) {
		t.Fatal("nil *TemplateType must be treated as null")
	}
	if !hamtIsNull(NilNode) {
		t.Fatal("NilNode must be treated as null")
	}
	reg := NewJSTypeRegistry(nil)
	num := reg.GetNativeType(JSTypeNative_NUMBER_TYPE)
	str := reg.GetNativeType(JSTypeNative_STRING_TYPE)
	ta := reg.CreateUnionType(num, str)
	tb := reg.CreateUnionType(str, num)
	if !hamtValueEquals(ta, tb) {
		t.Fatal("equal JSType values should match via jsTypeEquals")
	}
	if hamtValueEquals(ta, num) || hamtValueEquals(ta, nil) {
		t.Fatal("unequal or nil JSType values must not match")
	}
}

func TestHamtPMapZeroNonPointerValuesAndZeroHashableKeys(t *testing.T) {
	type zeroStructVal struct{ a, b int }
	if hamtIsNull(false) || hamtIsNull(int32(0)) || hamtIsNull(uint32(0)) || hamtIsNull(zeroStructVal{}) {
		t.Fatal("non-pointer zero values must not be classified as null")
	}

	boolMap := HamtPMap_Empty[string, bool]().Plus("f", false).Plus("t", true)
	if boolMap.Get("f") != false || boolMap.Get("t") != true {
		t.Fatal("bool HamtPMap Get mismatch")
	}
	i32Map := HamtPMap_Empty[string, int32]().Plus("z", 0)
	u32Map := HamtPMap_Empty[string, uint32]().Plus("z", 0)
	structMap := HamtPMap_Empty[string, zeroStructVal]().Plus("z", zeroStructVal{})
	if i32Map.IsEmpty() || u32Map.IsEmpty() || structMap.IsEmpty() {
		t.Fatal("maps storing non-pointer zero values must not be empty")
	}

	// Zero-valued struct key implementing HashCode() int32 must not panic with NPE.
	zeroKey := newIntegerWithCustomHashCode(0, 0)
	otherKey := newIntegerWithCustomHashCode(1, 0)
	missingKey := newIntegerWithCustomHashCode(2, 0)
	m1 := HamtPMap_Empty[integerWithCustomHashCode, int]().
		Plus(zeroKey, 0).
		Plus(otherKey, 7)
	if got := m1.Get(zeroKey); got != 0 {
		t.Fatalf("Get(zeroKey) = %d, want 0", got)
	}
	if got := m1.Get(otherKey); got != 7 {
		t.Fatalf("Get(otherKey) = %d, want 7", got)
	}

	// Hash-collision Equivalent when V has zero value 0 and key is absent on right side.
	m2 := HamtPMap_Empty[integerWithCustomHashCode, int]().
		Plus(otherKey, 7).
		Plus(missingKey, 0)
	if m1.Equivalent(m2, func(a, b int) bool { return a == b }) {
		t.Fatal("maps with distinct colliding keys and 0 values must not be equivalent")
	}

	// Hash-collision Reconcile when V has zero value 0 and key is absent on right side.
	seenZeroKey := false
	merged := m1.Reconcile(m2, func(k integerWithCustomHashCode, l, r int) int {
		if k == zeroKey {
			seenZeroKey = true
			return 99
		}
		return l + r
	})
	if !seenZeroKey || merged.Get(zeroKey) != 99 {
		t.Fatalf("expected joiner to run on colliding key with zero value; seen=%v, val=%d", seenZeroKey, merged.Get(zeroKey))
	}
}

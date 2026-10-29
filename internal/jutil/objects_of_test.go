package jutil

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

type objectsOfKey struct{ id int }

func (k objectsOfKey) Equals(o any) bool { x, ok := o.(objectsOfKey); return ok && x.id == k.id }
func (k objectsOfKey) HashCode() int32   { return int32(k.id * 31) }

type ptrHashableKey struct{ id int }

func (k *ptrHashableKey) Equals(o any) bool {
	if k == o {
		return true
	}
	x, ok := o.(*ptrHashableKey)
	return ok && x != nil && x.id == k.id
}
func (k *ptrHashableKey) HashCode() int32 { return int32(k.id*31 + 7) }

type plainPtrKey struct{ id int }

type testNodeID uint32

func (n testNodeID) IsNil() bool       { return n == 0 }
func (n testNodeID) HashCode() int32   { return int32(uint32(n) * 0x9E3779B1) }
func (n testNodeID) Equals(o any) bool { m, ok := o.(testNodeID); return ok && m == n }

type plainU32 uint32

func checkObjectsOf[T any](t *testing.T, vals []T) {
	t.Helper()
	for _, a := range vals {
		if got, want := isNilOf(a), IsNil(a); got != want {
			t.Errorf("isNilOf(%#v) = %v, IsNil %v", a, got, want)
		}
		if got, want := objectHashCodeOf(a), ObjectHashCode(a); got != want {
			t.Errorf("objectHashCodeOf(%#v) = %d, ObjectHashCode %d", a, got, want)
		}
		for _, b := range vals {
			if got, want := ObjectsEqualsOf(a, b), ObjectsEquals(a, b); got != want {
				t.Errorf("ObjectsEqualsOf(%#v, %#v) = %v, ObjectsEquals %v", a, b, got, want)
			}
		}
	}
}

func TestObjectsEqualsOfAndHashCodeOfMatchObjects(t *testing.T) {
	checkObjectsOf(t, []string{"", "a", "b", "ab", "é", "\U0001F600", string([]byte{0xed, 0xa0, 0x80})})
	checkObjectsOf(t, []int32{0, 1, -1, 1 << 30})
	checkObjectsOf(t, []int{0, 1, -1, 1 << 40})
	checkObjectsOf(t, []int64{0, 1, -1, math.MinInt64, math.MaxInt64})
	checkObjectsOf(t, []int16{0, 1, -1, math.MinInt16, math.MaxInt16})
	checkObjectsOf(t, []int8{0, 1, -1, math.MinInt8, math.MaxInt8})
	checkObjectsOf(t, []uint{0, 1, 1 << 40})
	checkObjectsOf(t, []uint64{0, 1, math.MaxUint64})
	checkObjectsOf(t, []uint32{0, 1, math.MaxUint32})
	checkObjectsOf(t, []uint16{0, 1, math.MaxUint16})
	checkObjectsOf(t, []uint8{0, 1, math.MaxUint8})
	checkObjectsOf(t, []bool{false, true})
	negZero64 := math.Copysign(0, -1)
	altNaN64 := math.Float64frombits(0x7ff8000000000001)
	checkObjectsOf(t, []float64{0.0, negZero64, 1.0, -1.0, math.Inf(1), math.Inf(-1), math.NaN(), altNaN64})
	negZero32 := float32(negZero64)
	altNaN32 := math.Float32frombits(0x7fc00001)
	checkObjectsOf(t, []float32{0.0, negZero32, 1.0, -1.0, float32(math.Inf(1)), float32(math.NaN()), altNaN32})
	checkObjectsOf(t, []*string{nil, StringPtr(""), StringPtr("x"), StringPtr("x"), StringPtr("y")})

	pk1a, pk1b, pk2 := &ptrHashableKey{1}, &ptrHashableKey{1}, &ptrHashableKey{2}
	checkObjectsOf(t, []*ptrHashableKey{nil, pk1a, pk1b, pk2})

	pp1, pp2 := &plainPtrKey{1}, &plainPtrKey{1}
	checkObjectsOf(t, []*plainPtrKey{nil, pp1, pp2})

	checkObjectsOf(t, []testNodeID{0, 1, 2, 1})
	n1 := testNodeID(1)
	checkObjectsOf(t, []*testNodeID{nil, &n1})
	checkObjectsOf(t, []plainU32{0, 1, 2})

	type plainStruct struct{ a, b int }
	checkObjectsOf(t, []plainStruct{{0, 0}, {1, 2}, {1, 2}})

	checkObjectsOf(t, []any{
		nil, "a", StringPtr("a"), StringPtr("a"), (*string)(nil),
		int32(1), 1, int64(1), uint32(0), uint32(1), uint64(1),
		0.0, negZero64, 1.0, math.NaN(), altNaN64,
		float32(0.0), negZero32, float32(math.NaN()),
		objectsOfKey{1}, objectsOfKey{2},
		(*int)(nil), (*ptrHashableKey)(nil), (*testNodeID)(nil),
		pk1a, pk1b, pk2, pp1, pp2,
		testNodeID(0), testNodeID(1), testNodeID(2),
		plainU32(0), plainU32(1),
		true, false,
	})
	checkObjectsOf(t, []objectsOfKey{{0}, {1}, {2}, {1}})
	type named string
	checkObjectsOf(t, []named{"a", "b"})
}

func objectsOfStringOps[T any](a, b T) bool {
	return ObjectsEqualsOf(a, b) && objectHashCodeOf(a) == objectHashCodeOf(b)
}

func TestObjectsEqualsOfDoesNotBoxStrings(t *testing.T) {
	a, b := "some key", string([]byte("some key"))
	if n := testing.AllocsPerRun(100, func() {
		if !objectsOfStringOps(a, b) {
			t.Fatal("not equal")
		}
	}); n != 0 {
		t.Errorf("%v allocations per call", n)
	}
}

func TestObjectsOfZeroAllocsScalarsAndPointers(t *testing.T) {
	k1 := &ptrHashableKey{42}
	k1Dup := &ptrHashableKey{42}
	k2 := &ptrHashableKey{99}
	var kNil *ptrHashableKey
	p1, p2 := &plainPtrKey{1}, &plainPtrKey{2}
	s1, s2 := StringPtr("hello"), StringPtr("hello")
	node := testNodeID(7)

	m := NewObjectLinkedHashMap[*ptrHashableKey, int]()
	for i := 0; i < 32; i++ {
		m.Put(&ptrHashableKey{i}, i)
	}
	m.Put(k1, 42)
	emptyMap := NewObjectLinkedHashMap[*ptrHashableKey, int]()

	if n := testing.AllocsPerRun(100, func() {
		if !objectsOfStringOps(123, 123) ||
			!objectsOfStringOps(int64(1)<<40, int64(1)<<40) ||
			!objectsOfStringOps(uint32(77), uint32(77)) ||
			!objectsOfStringOps(3.14, 3.14) ||
			!objectsOfStringOps(true, true) ||
			!objectsOfStringOps(s1, s2) ||
			!objectsOfStringOps(k1, k1) ||
			!objectsOfStringOps(k1, k1Dup) ||
			!objectsOfStringOps(kNil, kNil) ||
			ObjectsEqualsOf(k1, k2) ||
			ObjectsEqualsOf(k1, kNil) ||
			!objectsOfStringOps(p1, p1) ||
			ObjectsEqualsOf(p1, p2) ||
			!ObjectsEqualsOf(node, node) {
			t.Fatal("unexpected result")
		}
		if m.Get(k1) != 42 || m.Get(k1Dup) != 42 || m.ContainsKey(k2) || emptyMap.ContainsKey(k1) {
			t.Fatal("unexpected map lookup result")
		}
		_ = CheckNotNull("abc")
		_ = CheckNotNull(123)
		_ = CheckNotNull(k1)
		_ = CheckNotNull(node)
	}); n != 0 {
		t.Fatalf("expected 0 allocations, got %v", n)
	}
}

func TestTryParseHelpersMatchJava(t *testing.T) {
	for _, r := range readTSV(t, "parse_double.tsv") {
		in := StringFromUTF16(decodeHexUTF16(r[0]))
		got, ok := TryParseDouble(in)
		wantErr := strings.HasPrefix(r[1], "E:")
		if ok == wantErr {
			t.Fatalf("TryParseDouble(%q) ok=%v, Java=%s", in, ok, r[1])
		}
		if ok {
			if gotHex := strconv.FormatUint(math.Float64bits(got), 16); gotHex != r[1] {
				t.Fatalf("TryParseDouble(%q) = %s, Java %s", in, gotHex, r[1])
			}
		}
	}
	for _, r := range readTSV(t, "bigint_radix.tsv") {
		in := StringFromUTF16(decodeHexUTF16(r[0]))
		radix, _ := strconv.Atoi(r[1])
		got, ok := TryParseBigInt(in, radix)
		wantErr := strings.HasPrefix(r[2], "E:")
		if ok == wantErr {
			t.Fatalf("TryParseBigInt(%q, %d) ok=%v, Java=%s", in, radix, ok, r[2])
		}
		if ok && got.String() != r[2] {
			t.Fatalf("TryParseBigInt(%q, %d) = %s, Java %s", in, radix, got, r[2])
		}
	}
	if _, ok := TryParseLong("123", 1); ok {
		t.Fatal("radix 1 should fail")
	}
	if _, ok := TryParseLong("123", 37); ok {
		t.Fatal("radix 37 should fail")
	}
	if v, ok := TryParseLong("-9223372036854775808", 10); !ok || v != math.MinInt64 {
		t.Fatalf("MinInt64: got %d, %v", v, ok)
	}
	if _, ok := TryParseLong("9223372036854775808", 10); ok {
		t.Fatal("MaxInt64+1 should fail")
	}
}

func BenchmarkObjectLinkedHashMapPointerHashable(b *testing.B) {
	m := NewObjectLinkedHashMap[*ptrHashableKey, int]()
	keys := make([]*ptrHashableKey, 32)
	dups := make([]*ptrHashableKey, 32)
	for i := range keys {
		keys[i] = &ptrHashableKey{i}
		dups[i] = &ptrHashableKey{i}
		m.Put(keys[i], i)
	}
	miss := &ptrHashableKey{999}
	b.Run("samePtrHit", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if m.Get(keys[i&31]) != i&31 {
				b.Fatal("miss")
			}
		}
	})
	b.Run("equalPtrHit", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if m.Get(dups[i&31]) != i&31 {
				b.Fatal("miss")
			}
		}
	})
	b.Run("miss", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if m.ContainsKey(miss) {
				b.Fatal("unexpected hit")
			}
		}
	})
}

// IdentityHashCodeOfPointer gives IdentityHashCode's value without
// allocating.
func TestIdentityHashCodeOfPointer(t *testing.T) {
	type node struct{ a, b int }
	ns := make([]*node, 64)
	for i := range ns {
		ns[i] = &node{a: i}
		if got, want := IdentityHashCodeOfPointer(ns[i]), IdentityHashCode(ns[i]); got != want {
			t.Fatalf("IdentityHashCodeOfPointer = %d, IdentityHashCode = %d", got, want)
		}
	}
	var none *node
	if got, want := IdentityHashCodeOfPointer(none), IdentityHashCode(none); got != want || got != 0 {
		t.Fatalf("nil: IdentityHashCodeOfPointer = %d, IdentityHashCode = %d", got, want)
	}
	if n := testing.AllocsPerRun(100, func() { _ = IdentityHashCodeOfPointer(ns[3]) }); n != 0 {
		t.Fatalf("allocated %v times", n)
	}
}

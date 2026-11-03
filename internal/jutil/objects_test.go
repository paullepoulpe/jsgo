package jutil

import (
	"sync"
	"testing"
)

type negSentinelHandle int32

func (h negSentinelHandle) IsNil() bool { return h < 0 }

type nullableValStruct struct{ present bool }

func (s nullableValStruct) IsNil() bool { return !s.present }

func TestOrderedCollectionsEqualsAndHash(t *testing.T) {
	m1 := NewOrderedMap[int, *ptrHashableKey]()
	m2 := NewOrderedMap[int, *ptrHashableKey]()
	ls1 := NewOrderedSet[string]()
	ls2 := NewOrderedSet[string]()

	for i := 0; i < 20; i++ {
		m1.Put(i, &ptrHashableKey{i})
		m2.Put(19-i, &ptrHashableKey{19 - i})
	}
	ls1.AddAll("alpha", "beta", "gamma")
	ls2.AddAll("gamma", "alpha", "beta")

	if !m1.Equals(m2) || m1.HashCode() != m2.HashCode() {
		t.Fatal("OrderedMap Equals/HashCode mismatch")
	}
	if !ls1.Equals(ls2) || ls1.HashCode() != ls2.HashCode() {
		t.Fatal("OrderedSet Equals/HashCode mismatch")
	}
	strSlice := []string{"alpha", "beta", "gamma"}
	if SetHashCode(strSlice) != ls1.HashCode() {
		t.Fatal("SetHashCode mismatch")
	}

	if n := testing.AllocsPerRun(100, func() {
		if !ls1.Equals(ls2) {
			t.Fatal("unexpected mismatch")
		}
	}); n != 0 {
		t.Fatalf("expected 0 allocations for OrderedSet.Equals, got %v", n)
	}
}

func TestObjectsEqualsFastPathsAndEdgeCases(t *testing.T) {
	type plainStruct struct{ x int }
	p1, p2 := &plainStruct{1}, &plainStruct{1}
	if !ObjectsEquals(p1, p1) {
		t.Fatal("identical plain pointers must be equal")
	}
	if ObjectsEquals(p1, p2) {
		t.Fatal("distinct plain pointers without Equaler must use identity")
	}

	k1a, k1b, k2 := &ptrHashableKey{10}, &ptrHashableKey{10}, &ptrHashableKey{20}
	if !ObjectsEquals(k1a, k1a) || !ObjectsEquals(k1a, k1b) || ObjectsEquals(k1a, k2) {
		t.Fatal("ptrHashableKey ObjectsEquals mismatch")
	}

	var nilKey *ptrHashableKey
	var nilPlain *plainStruct
	if !ObjectsEquals(nilKey, nilPlain) || !ObjectsEquals(nilKey, nil) || !ObjectsEquals(nil, nilKey) {
		t.Fatal("typed-nil pointers must equal nil and each other")
	}
	if ObjectsEquals(k1a, nilKey) || ObjectsEquals(nilKey, k1a) || ObjectsEquals(p1, nilPlain) {
		t.Fatal("non-nil pointer must not equal typed-nil")
	}

	if !ObjectsEquals(testNodeID(0), nil) || !ObjectsEquals(testNodeID(0), nilKey) || !ObjectsEquals(nilKey, testNodeID(0)) {
		t.Fatal("zero testNodeID (IsNil=true) must equal nil")
	}
	if ObjectsEquals(uint32(0), nil) || ObjectsEquals(testNodeID(0), uint32(0)) {
		t.Fatal("plain uint32(0) is not nil")
	}
	if !ObjectsEquals(testNodeID(7), testNodeID(7)) || ObjectsEquals(testNodeID(7), testNodeID(8)) {
		t.Fatal("non-zero testNodeID equality mismatch")
	}

	if !ObjectsEquals(plainStruct{3}, plainStruct{3}) || ObjectsEquals(plainStruct{3}, plainStruct{4}) {
		t.Fatal("value struct equality mismatch")
	}

	s := []int{1, 2, 3}
	sAlias := s
	sCopy := []int{1, 2, 3}
	var nilSlice []int
	if !ObjectsEquals(nilSlice, nilKey) || ObjectsEquals(s, nilSlice) {
		t.Fatal("nil slice handling mismatch")
	}
	if !ObjectsEquals(s, sAlias) || ObjectsEquals(s, sCopy) || !ObjectsEquals([]int{}, []int{}) {
		t.Fatal("slice identity semantics mismatch")
	}

	if n := testing.AllocsPerRun(100, func() {
		_ = ObjectsEquals(k1a, k1b)
		_ = ObjectsEquals(k1a, k1a)
		_ = ObjectsEquals(p1, p2)
		_ = ObjectsEquals(nilKey, nilPlain)
	}); n != 0 {
		t.Fatalf("expected 0 allocations for pointer ObjectsEquals, got %v", n)
	}
}

func TestCustomIsNilTypesAndHelpfulNPEConcurrency(t *testing.T) {
	if !IsNil(negSentinelHandle(-1)) || IsNil(negSentinelHandle(0)) {
		t.Fatal("negSentinelHandle IsNil mismatch")
	}
	if !IsNil(nullableValStruct{present: false}) || IsNil(nullableValStruct{present: true}) {
		t.Fatal("nullableValStruct IsNil mismatch")
	}
	if !ObjectsEquals(negSentinelHandle(-1), nil) || !ObjectsEquals(nullableValStruct{}, negSentinelHandle(-1)) {
		t.Fatal("custom IsNil ObjectsEquals mismatch")
	}
	if ObjectsEquals(negSentinelHandle(0), nil) {
		t.Fatal("negSentinelHandle(0) must not equal nil")
	}
	_ = CheckNotNull(negSentinelHandle(0))
	_ = CheckNotNull(nullableValStruct{present: true})
	caught := false
	Try(func() { CheckNotNull(negSentinelHandle(-1)) }, func(INullPointerException) { caught = true })
	if !caught {
		t.Fatal("expected CheckNotNull(negSentinelHandle(-1)) to panic with NullPointerException")
	}

	orig := SetShowCodeDetailsInExceptionMessages(true)
	defer SetShowCodeDetailsInExceptionMessages(orig)
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				SetShowCodeDetailsInExceptionMessages((id+i)%2 == 0)
				_ = NewHelpfulNullPointerException("detail")
			}
		}(g)
	}
	wg.Wait()
}

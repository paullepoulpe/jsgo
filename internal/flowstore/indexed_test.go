// Copyright 2026 The jsgo Authors.
// Licensed under the Apache License, Version 2.0.

package flowstore

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type testValue struct{ version int }
type testMap = rhino.PMap[string, *testValue]

func choose(_ string, a, b *testValue) *testValue {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.version >= b.version {
		return a
	}
	return b
}
func equal(a, b *testValue) bool { return a.version%3 == b.version%3 }

func checkMap(t *testing.T, actual, reference testMap, names []string) {
	t.Helper()
	if actual.IsEmpty() != reference.IsEmpty() {
		t.Fatal("empty mismatch")
	}
	seen := map[string]bool{}
	for key := range actual.Keys() {
		if seen[key] {
			t.Fatalf("duplicate key %q", key)
		}
		seen[key] = true
	}
	count := 0
	for value := range actual.Values() {
		if value == nil {
			t.Fatal("nil value")
		}
		count++
	}
	if count != len(seen) {
		t.Fatalf("keys %d != values %d", len(seen), count)
	}
	for _, name := range names {
		got, want := actual.Get(name), reference.Get(name)
		if got != want {
			t.Fatalf("Get(%q): %v != %v", name, got, want)
		}
		if seen[name] != (want != nil) {
			t.Fatalf("key membership %q differs", name)
		}
	}
}

// The original Closure HAMT is the differential oracle. Persistent historical
// branches are repeatedly updated, deleted, joined and rechecked, including
// Java string hash collisions, astral Unicode, sparse IDs and the overflow path.
func TestStoreDifferentialPersistentBranches(t *testing.T) {
	for _, size := range []int{1, 4, 17, 64, MaxID, MaxID + 44} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			names := []string{"alpha", "beta", "gamma", "delta", "😀", "𝄞"}
			for len(names) < size {
				names = append(names, fmt.Sprintf("v%d.x", len(names)))
			}
			names = names[:size]
			owners := []*IDs{NewIDs()}
			actual := []testMap{New[*testValue](owners[0]), New[*testValue](owners[0]), New[*testValue](owners[0])}
			refs := []testMap{rhino.HamtPMap_Empty[string, *testValue](), rhino.HamtPMap_Empty[string, *testValue](), rhino.HamtPMap_Empty[string, *testValue]()}
			rng := rand.New(rand.NewSource(int64(size)))
			for step := 0; step < 1500; step++ {
				pick := rng.Intn(len(actual))
				key := names[rng.Intn(len(names))]
				got, want := actual[pick], refs[pick]
				op := rng.Intn(5)
				switch op {
				case 0, 1:
					value := &testValue{step + 1}
					got = got.Plus(key, value)
					want = want.Plus(key, value)
				case 2:
					got = got.Minus(key)
					want = want.Minus(key)
				case 3:
					other := rng.Intn(len(actual))
					var calls, referenceCalls []string
					got = got.Reconcile(actual[other], func(k string, a, b *testValue) *testValue {
						calls = append(calls, fmt.Sprintf("%s:%v:%v", k, a, b))
						return choose(k, a, b)
					})
					want = want.Reconcile(refs[other], func(k string, a, b *testValue) *testValue {
						referenceCalls = append(referenceCalls, fmt.Sprintf("%s:%v:%v", k, a, b))
						return choose(k, a, b)
					})
					if !reflect.DeepEqual(calls, referenceCalls) {
						t.Fatalf("join callback order differs at %d: %v != %v", step, calls, referenceCalls)
					}
				case 4:
					other := rng.Intn(len(actual))
					var calls, referenceCalls []string
					actualEqual := got.Equivalent(actual[other], func(a, b *testValue) bool {
						if a != b {
							calls = append(calls, fmt.Sprintf("%v:%v", a, b))
						}
						return equal(a, b)
					})
					referenceEqual := want.Equivalent(refs[other], func(a, b *testValue) bool {
						if a != b {
							referenceCalls = append(referenceCalls, fmt.Sprintf("%v:%v", a, b))
						}
						return equal(a, b)
					})
					if !reflect.DeepEqual(calls, referenceCalls) {
						t.Fatalf("equality callback order differs at %d: %v != %v", step, calls, referenceCalls)
					}
					if actualEqual != referenceEqual {
						t.Fatalf("equivalence differs at %d", step)
					}
				}
				if step%13 == 0 && len(actual) < 48 {
					actual = append(actual, got)
					refs = append(refs, want)
				} else {
					actual[pick] = got
					refs[pick] = want
				}
				checkMap(t, got, want, names)
				if step%67 == 0 {
					for i := range actual {
						checkMap(t, actual[i], refs[i], names)
					}
				}
			}
		})
	}
}

func TestStoreOverflowAndHistoricalVersion(t *testing.T) {
	ids := NewIDs()
	base := New[*testValue](ids)
	saved := &testValue{1}
	old := base.Plus("saved", saved)
	current := old
	for i := 0; i < MaxID+20; i++ {
		current = current.Plus(fmt.Sprint(i), &testValue{i + 2})
	}
	if current.(*Store[*testValue]).fallback == nil {
		t.Fatal("overflow did not activate HAMT")
	}
	if old.Get("saved") != saved || old.Get("1") != nil {
		t.Fatal("historical version changed")
	}
	joined := old.Reconcile(current, choose)
	for i := 0; i < MaxID+20; i++ {
		if joined.Get(fmt.Sprint(i)) != current.Get(fmt.Sprint(i)) {
			t.Fatal("overflow join lost entry")
		}
	}
	for key := range joined.Keys() {
		joined = joined.Minus(key)
	}
	if !joined.IsEmpty() {
		t.Fatal("deletion did not empty fallback")
	}
}

func TestStoreNoChangeAndSharedPages(t *testing.T) {
	m := New[*testValue](NewIDs())
	value := &testValue{1}
	one := m.Plus("x", value)
	if one.Plus("x", value) != one || one.Minus("missing") != one {
		t.Fatal("no-op update lost identity")
	}
	if one.Reconcile(one, func(string, *testValue, *testValue) *testValue { t.Fatal("same map joined"); return nil }) != one {
		t.Fatal("self join lost identity")
	}
	two := one.Plus("y", &testValue{2})
	calls := 0
	merged := one.Reconcile(two, func(_ string, a, b *testValue) *testValue {
		calls++
		if a != nil {
			t.Fatal("shared slot should skip join")
		}
		return b
	})
	if calls != 1 || merged != two {
		t.Fatal("join failed to preserve unchanged version")
	}
}

func TestStoreConcurrentSymbolInterning(t *testing.T) {
	base := New[*testValue](NewIDs())
	// The underlying original HAMT singleton cache is initialized serially.
	rhino.HamtPMap_Empty[string, *testValue]()
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			var m testMap = base
			for i := 0; i < 50; i++ {
				key := fmt.Sprintf("%d.%d", worker, i)
				value := &testValue{i + 1}
				m = m.Plus(key, value)
				if m.Get(key) != value {
					t.Errorf("lost %s", key)
				}
			}
			if !base.IsEmpty() {
				t.Error("shared empty version changed")
			}
		}(worker)
	}
	wg.Wait()
}

func BenchmarkStoreVersions(b *testing.B) {
	for _, size := range []int{4, 16, 64, 256} {
		for _, native := range []bool{false, true} {
			label := "hamt"
			if native {
				label = "indexed"
			}
			b.Run(fmt.Sprintf("%s/%d", label, size), func(b *testing.B) {
				var m testMap = rhino.HamtPMap_Empty[string, *testValue]()
				if native {
					m = New[*testValue](NewIDs())
				}
				names := make([]string, size)
				values := make([]*testValue, size)
				for i := range names {
					names[i] = fmt.Sprintf("v%d", i)
					values[i] = &testValue{i + 1}
					m = m.Plus(names[i], values[i])
				}
				alternative := make([]*testValue, size)
				for i := range alternative {
					alternative[i] = &testValue{i + size + 1}
				}
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					k := i % size
					next := m.Plus(names[k], alternative[k])
					if next.Get(names[k]) != alternative[k] {
						b.Fatal("incorrect update")
					}
					if !m.Equivalent(m, equal) {
						b.Fatal("self equality")
					}
					if next.Reconcile(m, choose).Get(names[k]) != alternative[k] {
						b.Fatal("incorrect join")
					}
				}
			})
		}
	}
}

func TestStoreDifferentOwnersAndReferenceMap(t *testing.T) {
	a, b := New[*testValue](NewIDs()), New[*testValue](NewIDs())
	x, y := &testValue{1}, &testValue{2}
	left := a.Plus("x", x)
	right := b.Plus("y", y)
	reference := rhino.HamtPMap_Empty[string, *testValue]().Plus("y", y)
	for _, other := range []testMap{right, reference} {
		merged := left.Reconcile(other, choose)
		if merged.Get("x") != x || merged.Get("y") != y {
			t.Fatal("mixed-owner join lost values")
		}
		if !merged.Equivalent(left.Reconcile(reference, choose), equal) {
			t.Fatal("mixed-owner equality differs")
		}
	}
}

func TestStoreJavaHashCollisionPreservesOriginalHistory(t *testing.T) {
	names := []string{"Aa", "BB", "AaAa", "BBBB"}
	ids := NewIDs()
	actual := []testMap{New[*testValue](ids)}
	refs := []testMap{rhino.HamtPMap_Empty[string, *testValue]()}
	rng := rand.New(rand.NewSource(4))
	for step := 0; step < 2000; step++ {
		pick := rng.Intn(len(actual))
		key := names[rng.Intn(len(names))]
		got, want := actual[pick], refs[pick]
		switch rng.Intn(4) {
		case 0, 1:
			value := &testValue{step + 1}
			got = got.Plus(key, value)
			want = want.Plus(key, value)
		case 2:
			got = got.Minus(key)
			want = want.Minus(key)
		case 3:
			other := rng.Intn(len(actual))
			got = got.Reconcile(actual[other], choose)
			want = want.Reconcile(refs[other], choose)
		}
		if step%17 == 0 && len(actual) < 64 {
			actual = append(actual, got)
			refs = append(refs, want)
		} else {
			actual[pick] = got
			refs[pick] = want
		}
		checkMap(t, got, want, names)
	}
}

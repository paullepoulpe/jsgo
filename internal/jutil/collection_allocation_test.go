package jutil

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBatchedCollectionMutationDuringIteration(t *testing.T) {
	m := NewOrderedMap[int, int]()
	for i := 0; i < 300; i++ {
		m.Put(i, i+1000)
	}
	var got []int
	for k, v := range m.All() {
		if v != k+1000 {
			t.Fatalf("value lost for %d", k)
		}
		got = append(got, k)
		if k == 10 {
			if old := m.Remove(11); old != 1011 {
				t.Fatalf("remove returned %d", old)
			}
			m.Put(300, 1300)
		}
		if k == 20 {
			m.Clear()
			m.Put(400, 1400)
		}
	}
	var want []int
	last := 20
	for i := 0; i <= last; i++ {
		if i != 11 {
			want = append(want, i)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("iteration: %v, want %v", got, want)
	}
	for k, v := range m.All() {
		if k != 400 || v != 1400 {
			t.Fatalf("clear/reinsert: %d=%d", k, v)
		}
	}
	for i := 0; i < 300; i++ {
		m.Put(i, i+1000)
	}
	for i := 0; i < 300; i++ {
		if m.Remove(i) != i+1000 {
			t.Fatalf("reused map lost %d", i)
		}
	}
}

func TestCollectionIterationZeroAllocs(t *testing.T) {
	lm := NewOrderedMap[int, int]()
	ls := NewOrderedSet[int]()
	for i := 0; i < 16; i++ {
		lm.Put(i, i)
		ls.Add(i)
	}
	var sink int
	if n := testing.AllocsPerRun(100, func() {
		for k, v := range lm.All() {
			sink += k + v
		}
		for k := range lm.Keys() {
			sink += k
		}
		for v := range lm.Values() {
			sink += v
		}
		for k := range ls.All() {
			sink += k
		}
	}); n != 0 {
		t.Fatalf("expected 0 allocations during ordered collection iteration, got %v (sink=%d)", n, sink)
	}
}

func BenchmarkOrderedMapBuild(b *testing.B) {
	for _, size := range []int{0, 1, 8, 16, 128, 1024} {
		b.Run(fmt.Sprintf("ordered/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m := NewOrderedMap[int, int]()
				for j := 0; j < size; j++ {
					m.Put(j, j)
				}
				if m.Size() != size {
					b.Fatal(m.Size())
				}
			}
		})
	}
}

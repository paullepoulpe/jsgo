package jutil

import (
	"fmt"
	"iter"
	"reflect"
	"testing"
)

type allocationTestMap interface {
	Put(int, int) int
	Remove(int) int
	Clear()
	All() iter.Seq2[int, int]
}

func TestBatchedCollectionMutationDuringIteration(t *testing.T) {
	for name, create := range map[string]func() allocationTestMap{
		"linked": func() allocationTestMap { return NewLinkedHashMap[int, int]() },
		"object": func() allocationTestMap { return NewObjectLinkedHashMap[int, int]() },
	} {
		t.Run(name, func(t *testing.T) {
			m := create()
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
			if name == "object" {
				// ObjectLinkedHashMap.Clear replaces the map's storage;
				// an iterator already on the old chain finishes that chain.
				last = 300
			}
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
		})
	}
}

func BenchmarkOrderedMapBuild(b *testing.B) {
	for _, size := range []int{0, 1, 8, 16, 128, 1024} {
		b.Run(fmt.Sprintf("linked/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m := NewLinkedHashMap[int, int]()
				for j := 0; j < size; j++ {
					m.Put(j, j)
				}
				if m.Size() != size {
					b.Fatal(m.Size())
				}
			}
		})
		b.Run(fmt.Sprintf("object/%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m := NewObjectLinkedHashMap[int, int]()
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

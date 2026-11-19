// Clean-room specification based on Google Guava Immutable Collections (Apache 2.0).

package guava

import (
	"iter"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
)

// ImmutableList represents an immutable list of elements.
type ImmutableList[T any] struct {
	elems []T
}

// ImmutableList_Of creates an ImmutableList containing elems.
func ImmutableList_Of[T any](elems ...T) *ImmutableList[T] {
	cp := make([]T, len(elems))
	copy(cp, elems)
	return &ImmutableList[T]{elems: cp}
}

// ImmutableList_CopyOf creates an ImmutableList copying from slice.
func ImmutableList_CopyOf[T any](elems []T) *ImmutableList[T] {
	return ImmutableList_Of(elems...)
}

func (l *ImmutableList[T]) Size() int {
	if l == nil {
		return 0
	}
	return len(l.elems)
}

func (l *ImmutableList[T]) IsEmpty() bool {
	return l.Size() == 0
}

func (l *ImmutableList[T]) Get(i int) T {
	if l == nil || i < 0 || i >= len(l.elems) {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(i, l.Size()))
	}
	return l.elems[i]
}

func (l *ImmutableList[T]) Contains(elem T) bool {
	if l == nil {
		return false
	}
	for i := range l.elems {
		if javalang.ObjectsEquals(l.elems[i], elem) {
			return true
		}
	}
	return false
}

func (l *ImmutableList[T]) ToSlice() []T {
	if l == nil || len(l.elems) == 0 {
		return nil
	}
	cp := make([]T, len(l.elems))
	copy(cp, l.elems)
	return cp
}

func (l *ImmutableList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if l == nil {
			return
		}
		for _, e := range l.elems {
			if !yield(e) {
				return
			}
		}
	}
}

func (l *ImmutableList[T]) Equals(other any) bool {
	o, ok := other.(*ImmutableList[T])
	if !ok {
		return false
	}
	if l == o {
		return true
	}
	if l.Size() != o.Size() {
		return false
	}
	for i := range l.elems {
		if !javalang.ObjectsEquals(l.elems[i], o.elems[i]) {
			return false
		}
	}
	return true
}

func (l *ImmutableList[T]) HashCode() int32 {
	if l == nil {
		return 0
	}
	var h int32 = 1
	for _, e := range l.elems {
		h = 31*h + javalang.ObjectHashCode(e)
	}
	return h
}

func (l *ImmutableList[T]) String() string {
	if l == nil {
		return "null"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, e := range l.elems {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(javalang.ToJavaString(e))
	}
	sb.WriteByte(']')
	return sb.String()
}

// ---------------------------------------------------------------------------
// ImmutableSet
// ---------------------------------------------------------------------------

// ImmutableSet represents an immutable set preserving insertion order.
type ImmutableSet[T comparable] struct {
	set *javautil.LinkedHashSet[T]
}

// ImmutableSet_Of creates an ImmutableSet containing elems.
func ImmutableSet_Of[T comparable](elems ...T) *ImmutableSet[T] {
	return &ImmutableSet[T]{set: javautil.NewLinkedHashSet(elems...)}
}

// ImmutableSet_CopyOf creates an ImmutableSet copying from slice elems.
func ImmutableSet_CopyOf[T comparable](elems []T) *ImmutableSet[T] {
	return ImmutableSet_Of(elems...)
}

func (s *ImmutableSet[T]) Size() int {
	if s == nil || s.set == nil {
		return 0
	}
	return s.set.Size()
}

func (s *ImmutableSet[T]) IsEmpty() bool {
	return s.Size() == 0
}

func (s *ImmutableSet[T]) Contains(elem T) bool {
	if s == nil || s.set == nil {
		return false
	}
	return s.set.Contains(elem)
}

func (s *ImmutableSet[T]) ToSlice() []T {
	if s == nil || s.set == nil {
		return nil
	}
	return s.set.ToSlice()
}

func (s *ImmutableSet[T]) All() iter.Seq[T] {
	if s == nil || s.set == nil {
		return func(yield func(T) bool) {}
	}
	return s.set.All()
}

func (s *ImmutableSet[T]) Equals(other any) bool {
	o, ok := other.(*ImmutableSet[T])
	if !ok {
		return false
	}
	if s == o {
		return true
	}
	return s.set.Equals(o.set)
}

func (s *ImmutableSet[T]) HashCode() int32 {
	if s == nil || s.set == nil {
		return 0
	}
	return s.set.HashCode()
}

func (s *ImmutableSet[T]) String() string {
	if s == nil || s.set == nil {
		return "null"
	}
	return s.set.String()
}

// ---------------------------------------------------------------------------
// ImmutableMap
// ---------------------------------------------------------------------------

// ImmutableMap represents an immutable map preserving insertion order.
type ImmutableMap[K comparable, V any] struct {
	m *javautil.LinkedHashMap[K, V]
}

func ImmutableMap_Of[K comparable, V any]() *ImmutableMap[K, V] {
	return &ImmutableMap[K, V]{m: javautil.NewLinkedHashMap[K, V]()}
}

func ImmutableMap_Of1[K comparable, V any](k1 K, v1 V) *ImmutableMap[K, V] {
	m := javautil.NewLinkedHashMap[K, V]()
	m.Put(k1, v1)
	return &ImmutableMap[K, V]{m: m}
}

func ImmutableMap_Of2[K comparable, V any](k1 K, v1 V, k2 K, v2 V) *ImmutableMap[K, V] {
	m := javautil.NewLinkedHashMap[K, V]()
	m.Put(k1, v1)
	m.Put(k2, v2)
	return &ImmutableMap[K, V]{m: m}
}

func ImmutableMap_CopyOf[K comparable, V any](src *javautil.LinkedHashMap[K, V]) *ImmutableMap[K, V] {
	m := javautil.NewLinkedHashMapWithCapacity[K, V](src.Size())
	m.PutAll(src)
	return &ImmutableMap[K, V]{m: m}
}

func (m *ImmutableMap[K, V]) Size() int {
	if m == nil || m.m == nil {
		return 0
	}
	return m.m.Size()
}

func (m *ImmutableMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

func (m *ImmutableMap[K, V]) Get(k K) (V, bool) {
	var zero V
	if m == nil || m.m == nil {
		return zero, false
	}
	return m.m.Get(k)
}

func (m *ImmutableMap[K, V]) GetVal(k K) V {
	v, _ := m.Get(k)
	return v
}

func (m *ImmutableMap[K, V]) ContainsKey(k K) bool {
	if m == nil || m.m == nil {
		return false
	}
	return m.m.ContainsKey(k)
}

func (m *ImmutableMap[K, V]) All() iter.Seq2[K, V] {
	if m == nil || m.m == nil {
		return func(yield func(K, V) bool) {}
	}
	return m.m.All()
}

func (m *ImmutableMap[K, V]) Keys() iter.Seq[K] {
	if m == nil || m.m == nil {
		return func(yield func(K) bool) {}
	}
	return m.m.Keys()
}

func (m *ImmutableMap[K, V]) Equals(other any) bool {
	o, ok := other.(*ImmutableMap[K, V])
	if !ok {
		return false
	}
	if m == o {
		return true
	}
	return m.m.Equals(o.m)
}

func (m *ImmutableMap[K, V]) HashCode() int32 {
	if m == nil || m.m == nil {
		return 0
	}
	return m.m.HashCode()
}

func (m *ImmutableMap[K, V]) String() string {
	if m == nil || m.m == nil {
		return "null"
	}
	return m.m.String()
}

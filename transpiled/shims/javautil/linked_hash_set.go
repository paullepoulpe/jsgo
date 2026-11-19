// Clean-room specification based on java.util.LinkedHashSet & java.util.HashSet Javadoc.

package javautil

import (
	"iter"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// LinkedHashSet is a deterministic, insertion-ordered set backed by LinkedHashMap.
// Ranging over the set yields elements in predictable insertion order.
type LinkedHashSet[T comparable] struct {
	m LinkedHashMap[T, struct{}]
}

// HashSet is an alias to LinkedHashSet. In this port, all hash sets preserve insertion
// order to guarantee deterministic iteration and byte-for-byte parity.
type HashSet[T comparable] = LinkedHashSet[T]

// NewLinkedHashSet creates a LinkedHashSet containing the optional initial elements.
func NewLinkedHashSet[T comparable](elems ...T) *LinkedHashSet[T] {
	s := &LinkedHashSet[T]{}
	for _, e := range elems {
		s.Add(e)
	}
	return s
}

// NewLinkedHashSetWithCapacity creates a LinkedHashSet with pre-allocated capacity.
func NewLinkedHashSetWithCapacity[T comparable](capHint int) *LinkedHashSet[T] {
	return &LinkedHashSet[T]{m: *NewLinkedHashMapWithCapacity[T, struct{}](capHint)}
}

// NewHashSet creates a new HashSet with deterministic insertion order.
func NewHashSet[T comparable](elems ...T) *HashSet[T] {
	return NewLinkedHashSet(elems...)
}

// NewHashSetWithCapacity creates a new HashSet with capacity hint.
func NewHashSetWithCapacity[T comparable](capHint int) *HashSet[T] {
	return NewLinkedHashSetWithCapacity[T](capHint)
}

// Add adds elem to the set, returning true if not already present.
func (s *LinkedHashSet[T]) Add(elem T) bool {
	return s.m.PutIfAbsent(elem, struct{}{})
}

// Remove removes elem from the set, returning true if present.
func (s *LinkedHashSet[T]) Remove(elem T) bool {
	_, ok := s.m.Remove(elem)
	return ok
}

// Contains reports whether elem is in the set.
func (s *LinkedHashSet[T]) Contains(elem T) bool {
	return s.m.ContainsKey(elem)
}

// Size returns the number of elements in the set.
func (s *LinkedHashSet[T]) Size() int {
	return s.m.Size()
}

// IsEmpty reports whether the set contains no elements.
func (s *LinkedHashSet[T]) IsEmpty() bool {
	return s.m.IsEmpty()
}

// Clear removes all elements from the set.
func (s *LinkedHashSet[T]) Clear() {
	s.m.Clear()
}

// ToSlice returns all elements as a slice in insertion order.
func (s *LinkedHashSet[T]) ToSlice() []T {
	if s == nil || s.Size() == 0 {
		return nil
	}
	res := make([]T, 0, s.Size())
	for elem := range s.All() {
		res = append(res, elem)
	}
	return res
}

// Iterator returns an Iterator over the elements in insertion order.
func (s *LinkedHashSet[T]) Iterator() Iterator[T] {
	slice := s.ToSlice()
	return &setIterator[T]{set: s, slice: slice, idx: 0, lastRet: -1}
}

// All returns an iter.Seq[T] yielding elements in insertion order.
func (s *LinkedHashSet[T]) All() iter.Seq[T] {
	return s.m.Keys()
}

// AddAll adds all elements to the set.
func (s *LinkedHashSet[T]) AddAll(elems ...T) bool {
	modified := false
	for _, e := range elems {
		if s.Add(e) {
			modified = true
		}
	}
	return modified
}

// AddAllSet adds all elements from other to this set.
func (s *LinkedHashSet[T]) AddAllSet(other *LinkedHashSet[T]) bool {
	if other == nil {
		return false
	}
	modified := false
	for elem := range other.All() {
		if s.Add(elem) {
			modified = true
		}
	}
	return modified
}

// RetainAll retains only elements that are contained in other.
func (s *LinkedHashSet[T]) RetainAll(other *LinkedHashSet[T]) bool {
	if other == nil {
		s.Clear()
		return true
	}
	modified := false
	for _, elem := range s.ToSlice() {
		if !other.Contains(elem) {
			s.Remove(elem)
			modified = true
		}
	}
	return modified
}

// RemoveAll removes all elements contained in other.
func (s *LinkedHashSet[T]) RemoveAll(other *LinkedHashSet[T]) bool {
	if other == nil {
		return false
	}
	modified := false
	for _, elem := range s.ToSlice() {
		if other.Contains(elem) {
			if s.Remove(elem) {
				modified = true
			}
		}
	}
	return modified
}

// Equals reports whether other is a set containing the exact same elements.
func (s *LinkedHashSet[T]) Equals(other any) bool {
	o, ok := other.(*LinkedHashSet[T])
	if !ok {
		return false
	}
	if s == o {
		return true
	}
	if s.Size() != o.Size() {
		return false
	}
	for elem := range s.All() {
		if !o.Contains(elem) {
			return false
		}
	}
	return true
}

// HashCode returns the Java set hash code: sum of elem.hashCode().
func (s *LinkedHashSet[T]) HashCode() int32 {
	if s == nil {
		return 0
	}
	var h int32
	for elem := range s.All() {
		h += javalang.ObjectHashCode(elem)
	}
	return h
}

// String formats the set as "[elem0, elem1, ...]".
func (s *LinkedHashSet[T]) String() string {
	if s == nil {
		return "null"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for elem := range s.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(javalang.ToJavaString(elem))
	}
	sb.WriteByte(']')
	return sb.String()
}

type setIterator[T comparable] struct {
	set     *LinkedHashSet[T]
	slice   []T
	idx     int
	lastRet int
}

func (it *setIterator[T]) HasNext() bool {
	return it.idx < len(it.slice)
}

func (it *setIterator[T]) Next() T {
	if !it.HasNext() {
		panic(javalang.NewNoSuchElementException("No more elements in set iterator"))
	}
	v := it.slice[it.idx]
	it.lastRet = it.idx
	it.idx++
	return v
}

func (it *setIterator[T]) Remove() {
	if it.lastRet < 0 {
		panic(javalang.NewIllegalStateException("Remove called before Next or twice"))
	}
	it.set.Remove(it.slice[it.lastRet])
	it.lastRet = -1
}

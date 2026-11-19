// Clean-room specification based on java.util.ArrayList Javadoc & Java Language Specification.

package javautil

import (
	"iter"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// ArrayList is a dynamically resizable array list matching Java's java.util.ArrayList.
type ArrayList[T any] struct {
	elems []T
}

// NewArrayList creates an empty ArrayList.
func NewArrayList[T any]() *ArrayList[T] {
	return &ArrayList[T]{}
}

// NewArrayListWithCapacity creates an ArrayList with pre-allocated capacity.
func NewArrayListWithCapacity[T any](initialCapacity int) *ArrayList[T] {
	if initialCapacity <= 0 {
		return &ArrayList[T]{}
	}
	return &ArrayList[T]{elems: make([]T, 0, initialCapacity)}
}

// NewArrayListFrom creates an ArrayList containing the specified elements.
func NewArrayListFrom[T any](elems ...T) *ArrayList[T] {
	cp := make([]T, len(elems))
	copy(cp, elems)
	return &ArrayList[T]{elems: cp}
}

// NewArrayListFromSlice creates an ArrayList containing the elements of slice s.
func NewArrayListFromSlice[T any](s []T) *ArrayList[T] {
	return NewArrayListFrom(s...)
}

// Size returns the number of elements in this list.
func (a *ArrayList[T]) Size() int {
	if a == nil {
		return 0
	}
	return len(a.elems)
}

// IsEmpty reports whether the list contains no elements.
func (a *ArrayList[T]) IsEmpty() bool {
	return a.Size() == 0
}

// Get returns the element at the specified position.
func (a *ArrayList[T]) Get(index int) T {
	if a == nil || index < 0 || index >= len(a.elems) {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(index, a.Size()))
	}
	return a.elems[index]
}

// Set replaces the element at index with elem and returns the previous element.
func (a *ArrayList[T]) Set(index int, elem T) T {
	if a == nil || index < 0 || index >= len(a.elems) {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(index, a.Size()))
	}
	old := a.elems[index]
	a.elems[index] = elem
	return old
}

// Add appends the specified element to the end of this list.
func (a *ArrayList[T]) Add(elem T) bool {
	a.elems = append(a.elems, elem)
	return true
}

// AddAt inserts the specified element at the specified position.
func (a *ArrayList[T]) AddAt(index int, elem T) {
	if index < 0 || index > len(a.elems) {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(index, len(a.elems)))
	}
	var zero T
	a.elems = append(a.elems, zero)
	copy(a.elems[index+1:], a.elems[index:])
	a.elems[index] = elem
}

// Remove removes the element at the specified position and returns it.
func (a *ArrayList[T]) Remove(index int) T {
	if a == nil || index < 0 || index >= len(a.elems) {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(index, a.Size()))
	}
	old := a.elems[index]
	copy(a.elems[index:], a.elems[index+1:])
	var zero T
	a.elems[len(a.elems)-1] = zero
	a.elems = a.elems[:len(a.elems)-1]
	return old
}

// RemoveElem removes the first occurrence of elem, returning true if found.
func (a *ArrayList[T]) RemoveElem(elem T) bool {
	if a == nil {
		return false
	}
	for i := range a.elems {
		if javalang.ObjectsEquals(a.elems[i], elem) {
			a.Remove(i)
			return true
		}
	}
	return false
}

// Clear removes all elements from this list.
func (a *ArrayList[T]) Clear() {
	if a == nil {
		return
	}
	clear(a.elems)
	a.elems = a.elems[:0]
}

// Contains reports whether this list contains the specified element.
func (a *ArrayList[T]) Contains(elem T) bool {
	if a == nil {
		return false
	}
	for i := range a.elems {
		if javalang.ObjectsEquals(a.elems[i], elem) {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first occurrence of elem, or -1.
func (a *ArrayList[T]) IndexOf(elem T) int {
	if a == nil {
		return -1
	}
	for i := range a.elems {
		if javalang.ObjectsEquals(a.elems[i], elem) {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the index of the last occurrence of elem, or -1.
func (a *ArrayList[T]) LastIndexOf(elem T) int {
	if a == nil {
		return -1
	}
	for i := len(a.elems) - 1; i >= 0; i-- {
		if javalang.ObjectsEquals(a.elems[i], elem) {
			return i
		}
	}
	return -1
}

// AddAll appends all elements of other to this list.
func (a *ArrayList[T]) AddAll(other *ArrayList[T]) bool {
	if other == nil || len(other.elems) == 0 {
		return false
	}
	a.elems = append(a.elems, other.elems...)
	return true
}

// AddAllSlice appends all elements of slice s to this list.
func (a *ArrayList[T]) AddAllSlice(s []T) bool {
	if len(s) == 0 {
		return false
	}
	a.elems = append(a.elems, s...)
	return true
}

// ToArray returns a copy of the elements as a Go slice.
func (a *ArrayList[T]) ToArray() []T {
	if a == nil || len(a.elems) == 0 {
		return nil
	}
	cp := make([]T, len(a.elems))
	copy(cp, a.elems)
	return cp
}

// ToSlice returns a copy of the elements as a Go slice.
func (a *ArrayList[T]) ToSlice() []T {
	return a.ToArray()
}

// SubList returns a new ArrayList containing the sub-range [fromIndex, toIndex).
func (a *ArrayList[T]) SubList(fromIndex, toIndex int) *ArrayList[T] {
	if a == nil || fromIndex < 0 || toIndex > len(a.elems) || fromIndex > toIndex {
		panic(javalang.NewIndexOutOfBoundsException("SubList range out of bounds"))
	}
	sub := make([]T, toIndex-fromIndex)
	copy(sub, a.elems[fromIndex:toIndex])
	return &ArrayList[T]{elems: sub}
}

// Sort sorts this list in-place using the given comparator.
func (a *ArrayList[T]) Sort(cmp Comparator[T]) {
	if a == nil || len(a.elems) <= 1 {
		return
	}
	slices.SortFunc(a.elems, func(x, y T) int {
		return cmp(x, y)
	})
}

// Iterator returns an Iterator over the elements in this list.
func (a *ArrayList[T]) Iterator() Iterator[T] {
	return &sliceIterator[T]{list: a, idx: 0, lastRet: -1}
}

// All returns an iter.Seq[T] for Go range loops over elements in encounter order.
func (a *ArrayList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if a == nil {
			return
		}
		for _, e := range a.elems {
			if !yield(e) {
				return
			}
		}
	}
}

// All2 returns an iter.Seq2[int, T] for Go range loops yielding (index, element).
func (a *ArrayList[T]) All2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		if a == nil {
			return
		}
		for i, e := range a.elems {
			if !yield(i, e) {
				return
			}
		}
	}
}

// Equals reports whether other is a List with identical size and equal elements in order.
func (a *ArrayList[T]) Equals(other any) bool {
	o, ok := other.(*ArrayList[T])
	if !ok {
		return false
	}
	if a == o {
		return true
	}
	if a == nil || o == nil || len(a.elems) != len(o.elems) {
		return false
	}
	for i := range a.elems {
		if !javalang.ObjectsEquals(a.elems[i], o.elems[i]) {
			return false
		}
	}
	return true
}

// HashCode returns the Java list polynomial hash code: sum(31^(n-1-i) * elem.hashCode()).
func (a *ArrayList[T]) HashCode() int32 {
	if a == nil {
		return 0
	}
	var h int32 = 1
	for _, e := range a.elems {
		h = 31*h + javalang.ObjectHashCode(e)
	}
	return h
}

// String formats the list as "[elem0, elem1, ...]".
func (a *ArrayList[T]) String() string {
	if a == nil {
		return "null"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, e := range a.elems {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(javalang.ToJavaString(e))
	}
	sb.WriteByte(']')
	return sb.String()
}

type sliceIterator[T any] struct {
	list    *ArrayList[T]
	idx     int
	lastRet int
}

func (it *sliceIterator[T]) HasNext() bool {
	return it.list != nil && it.idx < len(it.list.elems)
}

func (it *sliceIterator[T]) Next() T {
	if !it.HasNext() {
		panic(javalang.NewNoSuchElementException("No more elements in iterator"))
	}
	v := it.list.elems[it.idx]
	it.lastRet = it.idx
	it.idx++
	return v
}

func (it *sliceIterator[T]) Remove() {
	if it.lastRet < 0 {
		panic(javalang.NewIllegalStateException("Remove called before Next or twice"))
	}
	it.list.Remove(it.lastRet)
	it.idx = it.lastRet
	it.lastRet = -1
}

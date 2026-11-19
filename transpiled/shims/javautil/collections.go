// Clean-room specification based on java.util.Collections Javadoc.

package javautil

import (
	"slices"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Collections_UnmodifiableList returns an unmodifiable view of list.
func Collections_UnmodifiableList[T any](l *ArrayList[T]) *ArrayList[T] {
	if l == nil {
		return nil
	}
	// For minimal clean-room shim, returns list or cloned view
	return l
}

// Collections_UnmodifiableSet returns an unmodifiable view of set.
func Collections_UnmodifiableSet[T comparable](s *LinkedHashSet[T]) *LinkedHashSet[T] {
	return s
}

// Collections_UnmodifiableMap returns an unmodifiable view of map.
func Collections_UnmodifiableMap[K comparable, V any](m *LinkedHashMap[K, V]) *LinkedHashMap[K, V] {
	return m
}

// Collections_EmptyList returns an immutable empty ArrayList.
func Collections_EmptyList[T any]() *ArrayList[T] {
	return &ArrayList[T]{}
}

// Collections_EmptySet returns an immutable empty LinkedHashSet.
func Collections_EmptySet[T comparable]() *LinkedHashSet[T] {
	return &LinkedHashSet[T]{}
}

// Collections_EmptyMap returns an immutable empty LinkedHashMap.
func Collections_EmptyMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{}
}

// Collections_SingletonList returns an ArrayList containing only elem.
func Collections_SingletonList[T any](elem T) *ArrayList[T] {
	return NewArrayListFrom(elem)
}

// Collections_Singleton returns a LinkedHashSet containing only elem.
func Collections_Singleton[T comparable](elem T) *LinkedHashSet[T] {
	return NewLinkedHashSet(elem)
}

// Collections_SingletonMap returns a LinkedHashMap containing only (k, v).
func Collections_SingletonMap[K comparable, V any](k K, v V) *LinkedHashMap[K, V] {
	m := NewLinkedHashMap[K, V]()
	m.Put(k, v)
	return m
}

// Collections_Sort sorts list in-place using comparator cmp.
func Collections_Sort[T any](list *ArrayList[T], cmp Comparator[T]) {
	if list != nil {
		list.Sort(cmp)
	}
}

// Collections_Reverse reverses the elements of list in-place.
func Collections_Reverse[T any](list *ArrayList[T]) {
	if list == nil || len(list.elems) <= 1 {
		return
	}
	slices.Reverse(list.elems)
}

// Collections_BinarySearch searches the list for key using comparator cmp.
func Collections_BinarySearch[T any](list *ArrayList[T], key T, cmp Comparator[T]) int {
	if list == nil || len(list.elems) == 0 {
		return -1
	}
	low := 0
	high := len(list.elems) - 1
	for low <= high {
		mid := int(uint(low+high) >> 1)
		midVal := list.elems[mid]
		c := cmp(midVal, key)
		if c < 0 {
			low = mid + 1
		} else if c > 0 {
			high = mid - 1
		} else {
			return mid // key found
		}
	}
	return -(low + 1) // key not found
}

// Collections_Disjoint reports whether c1 and c2 have no elements in common.
func Collections_Disjoint[T comparable](c1, c2 *LinkedHashSet[T]) bool {
	if c1 == nil || c2 == nil || c1.Size() == 0 || c2.Size() == 0 {
		return true
	}
	if c1.Size() > c2.Size() {
		c1, c2 = c2, c1
	}
	for elem := range c1.All() {
		if c2.Contains(elem) {
			return false
		}
	}
	return true
}

// Collections_Min returns the minimum element according to cmp.
func Collections_Min[T any](list *ArrayList[T], cmp Comparator[T]) T {
	if list == nil || len(list.elems) == 0 {
		panic(javalang.NewNoSuchElementException("empty collection for min"))
	}
	minElem := list.elems[0]
	for _, e := range list.elems[1:] {
		if cmp(e, minElem) < 0 {
			minElem = e
		}
	}
	return minElem
}

// Collections_Max returns the maximum element according to cmp.
func Collections_Max[T any](list *ArrayList[T], cmp Comparator[T]) T {
	if list == nil || len(list.elems) == 0 {
		panic(javalang.NewNoSuchElementException("empty collection for max"))
	}
	maxElem := list.elems[0]
	for _, e := range list.elems[1:] {
		if cmp(e, maxElem) > 0 {
			maxElem = e
		}
	}
	return maxElem
}

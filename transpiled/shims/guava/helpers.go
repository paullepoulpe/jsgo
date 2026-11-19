// Clean-room specification based on Google Guava Lists, Maps, Sets (Apache 2.0).

package guava

import (
	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
)

// Lists_NewArrayList creates a new empty ArrayList.
func Lists_NewArrayList[T any]() *javautil.ArrayList[T] {
	return javautil.NewArrayList[T]()
}

// Lists_NewArrayListWithCapacity creates an ArrayList with pre-allocated capacity.
func Lists_NewArrayListWithCapacity[T any](cap int) *javautil.ArrayList[T] {
	return javautil.NewArrayListWithCapacity[T](cap)
}

// Lists_NewArrayListFrom creates an ArrayList containing elems.
func Lists_NewArrayListFrom[T any](elems ...T) *javautil.ArrayList[T] {
	return javautil.NewArrayListFrom(elems...)
}

// Lists_Reverse reverses list in-place.
func Lists_Reverse[T any](list *javautil.ArrayList[T]) {
	javautil.Collections_Reverse(list)
}

// Lists_Transform creates a new ArrayList by applying fn to each element.
func Lists_Transform[F, T any](list *javautil.ArrayList[F], fn func(F) T) *javautil.ArrayList[T] {
	if list == nil {
		return nil
	}
	out := javautil.NewArrayListWithCapacity[T](list.Size())
	for elem := range list.All() {
		out.Add(fn(elem))
	}
	return out
}

// Maps_NewHashMap creates a new deterministic insertion-ordered map.
func Maps_NewHashMap[K comparable, V any]() *javautil.LinkedHashMap[K, V] {
	return javautil.NewLinkedHashMap[K, V]()
}

// Maps_NewLinkedHashMap creates a new insertion-ordered map.
func Maps_NewLinkedHashMap[K comparable, V any]() *javautil.LinkedHashMap[K, V] {
	return javautil.NewLinkedHashMap[K, V]()
}

// Maps_NewHashMapWithExpectedSize creates a map sized for expectedSize.
func Maps_NewHashMapWithExpectedSize[K comparable, V any](expectedSize int) *javautil.LinkedHashMap[K, V] {
	return javautil.NewLinkedHashMapWithCapacity[K, V](expectedSize)
}

// Maps_UniqueIndex indexes values by the key returned by keyFn.
func Maps_UniqueIndex[K comparable, V any](values []V, keyFn func(V) K) *javautil.LinkedHashMap[K, V] {
	m := javautil.NewLinkedHashMapWithCapacity[K, V](len(values))
	for _, v := range values {
		k := keyFn(v)
		m.Put(k, v)
	}
	return m
}

// Sets_NewHashSet creates a new deterministic set with optional initial elements.
func Sets_NewHashSet[T comparable](elems ...T) *javautil.LinkedHashSet[T] {
	return javautil.NewLinkedHashSet[T](elems...)
}

// Sets_NewLinkedHashSet creates a new insertion-ordered set with optional initial elements.
func Sets_NewLinkedHashSet[T comparable](elems ...T) *javautil.LinkedHashSet[T] {
	return javautil.NewLinkedHashSet[T](elems...)
}

// Sets_NewHashSetWithExpectedSize creates a set sized for expectedSize.
func Sets_NewHashSetWithExpectedSize[T comparable](expectedSize int) *javautil.LinkedHashSet[T] {
	return javautil.NewLinkedHashSetWithCapacity[T](expectedSize)
}

// Sets_Difference returns an unmodifiable view or new set of elements in s1 not in s2.
func Sets_Difference[T comparable](s1, s2 *javautil.LinkedHashSet[T]) *javautil.LinkedHashSet[T] {
	diff := javautil.NewLinkedHashSet[T]()
	if s1 == nil {
		return diff
	}
	for elem := range s1.All() {
		if s2 == nil || !s2.Contains(elem) {
			diff.Add(elem)
		}
	}
	return diff
}

// Sets_Intersection returns an unmodifiable view or new set of elements in both s1 and s2.
func Sets_Intersection[T comparable](s1, s2 *javautil.LinkedHashSet[T]) *javautil.LinkedHashSet[T] {
	inter := javautil.NewLinkedHashSet[T]()
	if s1 == nil || s2 == nil {
		return inter
	}
	for elem := range s1.All() {
		if s2.Contains(elem) {
			inter.Add(elem)
		}
	}
	return inter
}

// Sets_Union returns a new set containing all elements in s1 or s2.
func Sets_Union[T comparable](s1, s2 *javautil.LinkedHashSet[T]) *javautil.LinkedHashSet[T] {
	union := javautil.NewLinkedHashSet[T]()
	if s1 != nil {
		for elem := range s1.All() {
			union.Add(elem)
		}
	}
	if s2 != nil {
		for elem := range s2.All() {
			union.Add(elem)
		}
	}
	return union
}

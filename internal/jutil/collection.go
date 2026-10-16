package jutil

import "iter"

// Go-only: language-adapter — the read side of java.util.Collection and
// java.util.Map that copy constructors (new HashSet<>(c), new HashMap<>(m))
// and bulk operations (addAll, removeAll, retainAll, putAll) take. Go cannot
// infer a type argument from a method set, so callers of the generic
// functions that take these name the type arguments explicitly.

// Collection is the part of java.util.Collection<T> that bulk operations
// read: the size, membership, and the elements in iteration order.
// *ArrayList, *LinkedHashSet, *ObjectLinkedHashSet, *JavaHashSet and
// JavaHashMap.KeySet() implement it.
type Collection[T any] interface {
	Size() int
	Contains(e T) bool
	All() iter.Seq[T]
}

// MapSource is the part of java.util.Map<K, V> that copy constructors and
// putAll read: the size and the entries in iteration order. *LinkedHashMap,
// *ObjectLinkedHashMap, *TreeMap and *JavaHashMap implement it.
type MapSource[K any, V any] interface {
	Size() int
	All() iter.Seq2[K, V]
}

var (
	_ Collection[string]     = (*ArrayList[string])(nil)
	_ Collection[string]     = (*LinkedHashSet[string])(nil)
	_ Collection[string]     = (*ObjectLinkedHashSet[string])(nil)
	_ MapSource[string, int] = (*LinkedHashMap[string, int])(nil)
	_ MapSource[string, int] = (*ObjectLinkedHashMap[string, int])(nil)
)

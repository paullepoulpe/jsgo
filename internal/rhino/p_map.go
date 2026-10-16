// Port of com/google/javascript/rhino/PMap.java @ closure-compiler v20260415

package rhino

import "iter"

// PMap is a minimal interface for null-hostile, persistent immutable maps.
//
// Go-only: language-adapter — Java null values are the values for which
// jutil.IsNil is true (nil pointers and interfaces), so V is normally a
// pointer or interface type; Iterable results are iter.Seq.
type PMap[K any, V any] interface {
	// IsEmpty returns whether this map is empty.
	IsEmpty() bool

	// Values returns an iterable for the values in this map.
	Values() iter.Seq[V]

	// Keys returns an iterable for the keys in this map.
	Keys() iter.Seq[K]

	// Get retrieves the given key from the map, or returns null if it is not
	// present.
	Get(key K) V

	// Plus returns a new map with the given key-value pair added. If the
	// value is already present, then this same map will be returned.
	Plus(key K, value V) PMap[K, V]

	// Minus returns a new map with the given key removed. If the key was not
	// present in the first place, then this same map will be returned.
	Minus(key K) PMap[K, V]

	// Reconcile performs a reconcile operation to merge `this` and `that`.
	//
	// `joiner` is called for each pair of entries, one from each map, which
	// share the same key and whose values are not Object#equals. This
	// includes entries that are absent from one of the maps, for which `null`
	// is passed as the absent value.
	//
	// The return of calling `joiner` will appear in the merged map at the key
	// of the original entry pair. The return may not be null. If the values
	// in a pair of entries are Object#equals, that value will be used directly
	// in the result without calling `joiner`.
	//
	// The first value passed to `joiner` comes from `this`, and the second
	// value comes from `that`. There are no guarantees on the source of
	// `key`. Note that `that` map must be the same implementation.
	Reconcile(that PMap[K, V], joiner PMap_Reconciler[K, V]) PMap[K, V]

	// Equivalent checks equality recursively based on the given equivalence.
	// Short-circuits as soon as a 'false' result is found, or if a key in one
	// map is missing from the other. The equivalence will only be called on
	// two non-null values. Note that `that` map must be the same
	// implementation. Note that the equivalence MUST be reflective (i.e.
	// equivalence.test(x, x) == true).
	Equivalent(that PMap[K, V], equivalence func(V, V) bool) bool
}

// PMap_Reconciler merges two values; see PMap.Reconcile.
//
// Go-only: language-adapter — a functional interface, so a Go func type.
//
// Java: PMap.Reconciler<K, V>.merge(K, V, V)
type PMap_Reconciler[K any, V any] func(key K, thisVal V, thatVal V) V

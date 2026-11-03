package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

const orderedMapIndexThreshold = 8

type (
	OrderedMapEntry[K comparable, V any] = jutil.OrderedMapEntry[K, V]
	OrderedMap[K comparable, V any]      = jutil.OrderedMap[K, V]
	OrderedStringEntry[V any]            = jutil.OrderedStringEntry[V]
	OrderedStringMap[V any]              = jutil.OrderedStringMap[V]
	OrderedSet[K comparable]             = jutil.OrderedSet[K]
)

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return jutil.NewOrderedMap[K, V]()
}

func NewOrderedStringMap[V any]() *OrderedStringMap[V] {
	return jutil.NewOrderedStringMap[V]()
}

func NewOrderedMapWithCapacity[K comparable, V any](capHint int) *OrderedMap[K, V] {
	return jutil.NewOrderedMapWithCapacity[K, V](capHint)
}

func NewOrderedStringMapWithCapacity[V any](capHint int) *OrderedStringMap[V] {
	return jutil.NewOrderedStringMapWithCapacity[V](capHint)
}

func OrderedMapEqual[K comparable, V comparable](a, b *OrderedMap[K, V]) bool {
	return jutil.OrderedMapEqual(a, b)
}

func NewOrderedSet[K comparable](elems ...K) *OrderedSet[K] {
	return jutil.NewOrderedSet(elems...)
}

func NewOrderedSetWithCapacity[K comparable](capHint int) *OrderedSet[K] {
	return jutil.NewOrderedSetWithCapacity[K](capHint)
}

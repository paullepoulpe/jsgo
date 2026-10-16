package jutil

import (
	"math"
	"strconv"
)

// Go-only: language-adapter — the hash-map factories of Guava's
// com.google.common.collect.Maps, as the Guava shaded into
// closure-compiler.jar implements them (checked with javap).

// Maps_Capacity is Maps.capacity(int): the HashMap initial capacity that
// holds expectedSize entries without resizing: expectedSize+1 below 3,
// (int) Math.ceil(expectedSize / 0.75) below 2^30, and Integer.MAX_VALUE
// beyond.
func Maps_Capacity(expectedSize int) int {
	if expectedSize < 3 {
		collectPreconditionsCheckNonnegative(expectedSize, "expectedSize")
		return expectedSize + 1
	}
	if expectedSize < 1<<30 {
		return int(D2I(math.Ceil(float64(expectedSize) / 0.75)))
	}
	return math.MaxInt32
}

// Maps_NewHashMap is Maps.newHashMap(): new HashMap<>().
func Maps_NewHashMap[K any, V any]() *JavaHashMap[K, V] { return NewJavaHashMap[K, V]() }

// Maps_NewHashMapFrom is Maps.newHashMap(Map): new HashMap<>(map).
func Maps_NewHashMapFrom[K any, V any](m MapSource[K, V]) *JavaHashMap[K, V] {
	return NewJavaHashMapFrom[K, V](m)
}

// Maps_NewHashMapWithExpectedSize is Maps.newHashMapWithExpectedSize(int):
// new HashMap<>(Maps.capacity(expectedSize)).
func Maps_NewHashMapWithExpectedSize[K any, V any](expectedSize int) *JavaHashMap[K, V] {
	return NewJavaHashMapWithCapacity[K, V](Maps_Capacity(expectedSize))
}

// collectPreconditionsCheckNonnegative is Guava's
// CollectPreconditions.checkNonnegative(int, String).
func collectPreconditionsCheckNonnegative(value int, name string) int {
	if value < 0 {
		panic(NewIllegalArgumentExceptionMsg(name + " cannot be negative but was: " + strconv.Itoa(value)))
	}
	return value
}

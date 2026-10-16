// Port of com/google/javascript/rhino/jstype/JSTypeIterations.java
// @ closure-compiler v20260415
//
// Stream-like helpers for querying slices / unions of JSType values.
// The Java version uses ImmutableList<? extends JSType>; Go uses
// []JSType slices throughout since we have no persistent-list lib.
// Predicate<? super JSType> → func(JSType) bool.
// Function<? super JSType, T> → func(JSType) JSType (T fixed to JSType).
//
// These functions minimise allocations because they are expected to
// be called very frequently — no iterators, no boxing, no closures
// that force heap escapes beyond what the caller already owns.

package jstype

// AnyTypeMatches reports whether any element of types satisfies
// predicate. Port of JSTypeIterations.anyTypeMatches(Predicate, ImmutableList).
func AnyTypeMatches(predicate func(JSType) bool, types []JSType) bool {
	for i := 0; i < len(types); i++ {
		if predicate(types[i]) {
			return true
		}
	}
	return false
}

// AnyTypeMatchesUnion reports whether any alternate of union satisfies
// predicate. Port of JSTypeIterations.anyTypeMatches(Predicate, UnionType).
func AnyTypeMatchesUnion(predicate func(JSType) bool, union *UnionType) bool {
	return AnyTypeMatches(predicate, union.Members())
}

// AllTypesMatch reports whether all elements of types satisfy predicate.
// Port of JSTypeIterations.allTypesMatch(Predicate, ImmutableList).
func AllTypesMatch(predicate func(JSType) bool, types []JSType) bool {
	for i := 0; i < len(types); i++ {
		if !predicate(types[i]) {
			return false
		}
	}
	return true
}

// AllTypesMatchUnion reports whether all alternates of union satisfy
// predicate. Port of JSTypeIterations.allTypesMatch(Predicate, UnionType).
func AllTypesMatchUnion(predicate func(JSType) bool, union *UnionType) bool {
	return AllTypesMatch(predicate, union.Members())
}

// MapTypes applies mapper to every element of types and returns the
// results as a new slice. Port of JSTypeIterations.mapTypes(Function, ImmutableList).
func MapTypes(mapper func(JSType) JSType, types []JSType) []JSType {
	out := make([]JSType, len(types))
	for i := 0; i < len(types); i++ {
		out[i] = mapper(types[i])
	}
	return out
}

// MapTypesUnion applies mapper to every alternate of union and returns
// their union. Port of JSTypeIterations.mapTypes(Function, UnionType).
func MapTypesUnion(mapper func(JSType) JSType, union *UnionType) JSType {
	mapped := MapTypes(mapper, union.Members())
	return NewUnion(mapped...)
}

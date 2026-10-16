// Port of com/google/javascript/rhino/jstype/JSTypeIterations.java @ closure-compiler v20260415

package rhino

// JSTypeIterations holds Stream-like functions for efficient querying of
// collections of JSTypes.
//
// These methods are designed to minimize allocations since they are
// expected to be called very often. That's why they don't:
//
//   - instantiate Iterators
//   - instantiate Streams
//   - (un)box primitives
//   - expect closure-generating lambdas
type JSTypeIterations struct{}

// JSTypeIterations_AnyTypeMatches returns whether any element of types
// matches predicate.
//
// Java: JSTypeIterations.anyTypeMatches(Predicate<? super JSType>, ImmutableList<? extends JSType>)
func JSTypeIterations_AnyTypeMatches[S JSType](predicate func(JSType) bool, types []S) bool {
	for i := 0; i < len(types); i++ {
		if predicate(types[i]) {
			return true
		}
	}
	return false
}

// JSTypeIterations_AnyTypeMatchesInUnion returns whether any element of
// union matches predicate.
//
// Java: JSTypeIterations.anyTypeMatches(Predicate<? super JSType>, UnionType)
func JSTypeIterations_AnyTypeMatchesInUnion(predicate func(JSType) bool, union *UnionType) bool {
	return JSTypeIterations_AnyTypeMatches(predicate, union.GetAlternates())
}

// JSTypeIterations_AllTypesMatch returns whether all elements of types
// match predicate.
//
// Java: JSTypeIterations.allTypesMatch(Predicate<? super JSType>, ImmutableList<? extends JSType>)
func JSTypeIterations_AllTypesMatch[S JSType](predicate func(JSType) bool, types []S) bool {
	for i := 0; i < len(types); i++ {
		if !predicate(types[i]) {
			return false
		}
	}
	return true
}

// JSTypeIterations_AllTypesMatchInUnion returns whether all elements of
// union match predicate.
//
// Java: JSTypeIterations.allTypesMatch(Predicate<? super JSType>, UnionType)
func JSTypeIterations_AllTypesMatchInUnion(predicate func(JSType) bool, union *UnionType) bool {
	return JSTypeIterations_AllTypesMatch(predicate, union.GetAlternates())
}

// JSTypeIterations_MapTypes returns the result of applying mapper to all
// elements of types.
//
// Java: JSTypeIterations.mapTypes(Function<? super JSType, T>, ImmutableList<? extends JSType>)
func JSTypeIterations_MapTypes[S JSType, T any](mapper func(S) T, types []S) []T {
	builder := make([]T, 0, len(types))
	for i := 0; i < len(types); i++ {
		builder = append(builder, mapper(types[i]))
	}
	return builder
}

// JSTypeIterations_MapTypesInUnion returns the union of applying mapper to
// all elements of union.
//
// Java: JSTypeIterations.mapTypes(Function<? super JSType, ? extends JSType>, UnionType)
func JSTypeIterations_MapTypesInUnion(mapper func(JSType) JSType, union *UnionType) JSType {
	return UnionType_Builder_(union.registry).
		AddAlternates(JSTypeIterations_MapTypes(mapper, union.GetAlternates())).
		Build()
}

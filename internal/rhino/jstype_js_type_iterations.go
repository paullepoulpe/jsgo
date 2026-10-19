// Port of com/google/javascript/rhino/jstype/JSTypeIterations.java @ closure-compiler v20261015

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

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
	b := UnionType_Builder_(union.registry)
	for _, alt := range union.GetAlternates() {
		b.AddAlternate(mapper(alt))
	}
	return b.Build()
}

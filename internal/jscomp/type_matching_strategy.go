// Port of com/google/javascript/jscomp/TypeMatchingStrategy.java @ closure-compiler v20260415

/*
 * Copyright 2015 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypeMatchingStrategy is the different strategies for matching the JSType
// of nodes.
type TypeMatchingStrategy int32

const (
	// TypeMatchingStrategy_LOOSE matches type or any subtype. Matches types
	// with different nullability/voidability. Allows loose matches.
	TypeMatchingStrategy_LOOSE TypeMatchingStrategy = iota

	// TypeMatchingStrategy_STRICT_NULLABILITY matches type or any subtype.
	// Does not match types with different nullability/voidability. Allows
	// loose matches.
	TypeMatchingStrategy_STRICT_NULLABILITY

	// TypeMatchingStrategy_SUBTYPES matches type or any subtype. Does not
	// match types with different nullability/voidability. Does not allow
	// loose matches.
	TypeMatchingStrategy_SUBTYPES

	// TypeMatchingStrategy_EXACT does not match subtypes. Does not match
	// types with different nullability/voidability. Does not allow loose
	// matches.
	TypeMatchingStrategy_EXACT
)

var typeMatchingStrategy_names = [...]string{"LOOSE", "STRICT_NULLABILITY", "SUBTYPES", "EXACT"}

// The constructor arguments of each constant, indexed by ordinal.
//
// Java: TypeMatchingStrategy(boolean, boolean, boolean)
var (
	typeMatchingStrategy_allowSubtypes     = [...]bool{true, true, true, false}
	typeMatchingStrategy_ignoreNullability = [...]bool{true, false, false, false}
	typeMatchingStrategy_allowLooseMatches = [...]bool{true, true, false, false}
)

// TypeMatchingStrategy_values is TypeMatchingStrategy.values().
func TypeMatchingStrategy_values() []TypeMatchingStrategy {
	return []TypeMatchingStrategy{
		TypeMatchingStrategy_LOOSE,
		TypeMatchingStrategy_STRICT_NULLABILITY,
		TypeMatchingStrategy_SUBTYPES,
		TypeMatchingStrategy_EXACT,
	}
}

// Name is Enum.name().
func (s TypeMatchingStrategy) Name() string { return typeMatchingStrategy_names[s] }

// String is Enum.toString().
func (s TypeMatchingStrategy) String() string { return s.Name() }

// Ordinal is Enum.ordinal().
func (s TypeMatchingStrategy) Ordinal() int { return int(s) }

// Java: TypeMatchingStrategy.match(JSType, JSType)
func (s TypeMatchingStrategy) Match(templateType rhino.JSType, type_ rhino.JSType) TypeMatchingStrategy_MatchResult {
	allowSubtypes := typeMatchingStrategy_allowSubtypes[s]
	ignoreNullability := typeMatchingStrategy_ignoreNullability[s]
	allowLooseMatches := typeMatchingStrategy_allowLooseMatches[s]

	if templateType.IsUnknownType() {
		// If the template type is '?', then any type is a match and this is not considered a loose
		// match.
		return TypeMatchingStrategy_MatchResult_MATCH
	}

	if type_ == nil || type_.IsUnknownType() || type_.IsAllType() {
		if allowLooseMatches {
			return TypeMatchingStrategy_MatchResult_LOOSE_MATCH
		}
		return TypeMatchingStrategy_MatchResult_NO_MATCH
	}

	if allowSubtypes {
		if ignoreNullability {
			type_ = type_.RestrictByNotNullOrUndefined()
		}
		if type_.IsSubtypeOf(templateType) {
			return TypeMatchingStrategy_MatchResult_MATCH
		}
	}

	nullableMismatch := templateType.IsNullable() != type_.IsNullable()
	voidableMismatch := templateType.IsVoidable() != type_.IsVoidable()
	if !ignoreNullability && (nullableMismatch || voidableMismatch) {
		return TypeMatchingStrategy_MatchResult_NO_MATCH
	}

	if type_.Equals(templateType) {
		return TypeMatchingStrategy_MatchResult_MATCH
	}
	return TypeMatchingStrategy_MatchResult_NO_MATCH
}

// TypeMatchingStrategy_MatchResult is the result of comparing two different
// JSType instances.
type TypeMatchingStrategy_MatchResult int32

const (
	TypeMatchingStrategy_MatchResult_MATCH TypeMatchingStrategy_MatchResult = iota
	TypeMatchingStrategy_MatchResult_NO_MATCH
	TypeMatchingStrategy_MatchResult_LOOSE_MATCH
)

var typeMatchingStrategy_MatchResult_names = [...]string{"MATCH", "NO_MATCH", "LOOSE_MATCH"}

// TypeMatchingStrategy_MatchResult_values is MatchResult.values().
func TypeMatchingStrategy_MatchResult_values() []TypeMatchingStrategy_MatchResult {
	return []TypeMatchingStrategy_MatchResult{
		TypeMatchingStrategy_MatchResult_MATCH,
		TypeMatchingStrategy_MatchResult_NO_MATCH,
		TypeMatchingStrategy_MatchResult_LOOSE_MATCH,
	}
}

// Name is Enum.name().
func (r TypeMatchingStrategy_MatchResult) Name() string {
	return typeMatchingStrategy_MatchResult_names[r]
}

// String is Enum.toString().
func (r TypeMatchingStrategy_MatchResult) String() string { return r.Name() }

// Ordinal is Enum.ordinal().
func (r TypeMatchingStrategy_MatchResult) Ordinal() int { return int(r) }

// Java: TypeMatchingStrategy.MatchResult.isMatch()
func (r TypeMatchingStrategy_MatchResult) IsMatch() bool {
	return r != TypeMatchingStrategy_MatchResult_NO_MATCH
}

// Java: TypeMatchingStrategy.MatchResult.isLooseMatch()
func (r TypeMatchingStrategy_MatchResult) IsLooseMatch() bool {
	return r == TypeMatchingStrategy_MatchResult_LOOSE_MATCH
}

// Port of com/google/javascript/rhino/jstype/EqualityChecker.java @ closure-compiler v20260415
/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1/GPL 2.0
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
 * Alternatively, the contents of this file may be used under the terms of
 * the GNU General Public License Version 2 or later (the "GPL"), in which
 * case the provisions of the GPL are applicable instead of those above. If
 * you wish to allow use of your version of this file only under the terms of
 * the GPL and not to allow others to use your version of this file under the
 * MPL, indicate your decision by deleting the provisions above and replacing
 * them with the notice and other provisions required by the GPL. If you do
 * not delete the provisions above, a recipient may use your version of this
 * file under either the MPL or the GPL.
 *
 * ***** END LICENSE BLOCK ***** */

// Port of com/google/javascript/rhino/jstype/EqualityChecker.java
// @ closure-compiler v20260415
//
// Represents the computation of a single equality relationship.
//
// Each instance can only be used once. It acts as a scope for the computation
// and may accumulate state.

package jstype

// EqMethod represents different ways for comparing equality among types.
// Port of EqualityChecker.EqMethod in Java.
type EqMethod int

const (
	// EqMethodIdentity indicates that the two types should behave exactly the
	// same under all type operations.
	// Thus, {string} != {?} and {Unresolved} != {?}
	EqMethodIdentity EqMethod = iota

	// EqMethodDataFlow indicates that the two types are almost exactly the
	// same, and that a data flow analysis algorithm comparing them should
	// consider them equal.
	// Thus, {string} != {?} and {Unresolved} ~= {?}
	EqMethodDataFlow
)

// eqMatchStatus mirrors JSType.MatchStatus in Java, used by EqualityChecker.
type eqMatchStatus int

const (
	eqProcessing eqMatchStatus = iota
	eqMatch
	eqNoMatch
)

func (s eqMatchStatus) eqValue() bool {
	return s == eqMatch
}

func eqMatchStatusOf(b bool) eqMatchStatus {
	if b {
		return eqMatch
	}
	return eqNoMatch
}

// eqCacheKey is a pair of JSType values used as a cache key.
// Port of EqualityChecker.CacheKey in Java.
// Note: unlike the subtype cache, equality is symmetric, so we normalize
// the key by identity hash (pointer value).
type eqCacheKey struct {
	left  JSType
	right JSType
}

// potentiallyCyclicEqRecursionDepth is an arbitrary depth at which to start
// checking for cyclic recursion. Port of POTENTIALLY_CYCLIC_RECURSION_DEPTH.
const potentiallyCyclicEqRecursionDepth = 20

// EqualityChecker represents the computation of a single equality relationship.
// Each instance can only be used once.
// Port of com/google/javascript/rhino/jstype/EqualityChecker.java.
type EqualityChecker struct {
	eqMethod *EqMethod

	eqCache        map[eqCacheKey]eqMatchStatus
	recursionDepth int
	hasRun         bool
}

// NewEqualityChecker constructs an EqualityChecker.
// Port of EqualityChecker() constructor in Java.
func NewEqualityChecker() *EqualityChecker {
	return &EqualityChecker{}
}

// SetEqMethod sets the equality method. Must be called before Check().
func (e *EqualityChecker) SetEqMethod(x EqMethod) *EqualityChecker {
	e.checkHasNotRun()
	if e.eqMethod != nil {
		panic("eqMethod already set")
	}
	e.eqMethod = &x
	return e
}

func (e *EqualityChecker) checkHasNotRun() {
	if e.hasRun {
		panic("EqualityChecker has already run")
	}
}

// Check runs the equality check for two types. May only be called once.
// Port of EqualityChecker.check() in Java.
func (e *EqualityChecker) Check(left, right JSType) bool {
	e.checkHasNotRun()
	e.hasRun = true
	return e.areEqualCaching(left, right)
}

// areEqualCaching is the caching entrypoint for equality logic.
// Port of EqualityChecker.areEqualCaching() in Java.
func (e *EqualityChecker) areEqualCaching(left, right JSType) bool {
	// Wait to instantiate/use the cache until we have some hint that there
	// may be recursion.
	if e.recursionDepth > potentiallyCyclicEqRecursionDepth {
		if e.eqCache == nil {
			e.eqCache = make(map[eqCacheKey]eqMatchStatus)
		}
	}

	// Once the cache exists, use it consistently.
	if e.eqCache == nil {
		e.recursionDepth++
		result := e.areEqualInternal(left, right)
		e.recursionDepth--
		return result
	}

	key := eqCacheKey{left: left, right: right}
	if cached, ok := e.eqCache[key]; ok {
		if cached == eqProcessing {
			e.eqCache[key] = eqMatch
			return true
		}
		return cached.eqValue()
	}

	e.eqCache[key] = eqProcessing
	result := e.areEqualInternal(left, right)
	e.eqCache[key] = eqMatchStatusOf(result)
	return result
}

// areEqualInternal is the core equality logic.
// Port of EqualityChecker.areEqualInternal() in Java.
func (e *EqualityChecker) areEqualInternal(left, right JSType) bool {
	if left == right {
		return true
	}
	if left == nil || right == nil {
		return false
	}

	leftUnknown := isUnknown(left)
	rightUnknown := isUnknown(right)
	if leftUnknown || rightUnknown {
		if e.eqMethod != nil && *e.eqMethod == EqMethodDataFlow {
			// If we're checking data flow, then two types are the same if
			// they're both unknown.
			return leftUnknown && rightUnknown
		}
	}

	// Union equality: both unions.
	leftUnion, leftIsUnion := left.(*UnionType)
	rightUnion, rightIsUnion := right.(*UnionType)
	if leftIsUnion && rightIsUnion {
		return e.areUnionEqual(leftUnion, rightUnion)
	} else if leftIsUnion {
		// A union with one alternate is equal to that alternate.
		if len(leftUnion.members) == 1 {
			return e.areEqualInternal(leftUnion.members[0], right)
		}
	} else if rightIsUnion {
		if len(rightUnion.members) == 1 {
			return e.areEqualInternal(left, rightUnion.members[0])
		}
	}

	// Function equality.
	leftFn, leftIsFn := left.(*FunctionType)
	rightFn, rightIsFn := right.(*FunctionType)
	if leftIsFn && rightIsFn {
		return e.areFunctionEqual(leftFn, rightFn)
	}

	// Template type map equality check (simplified: our TemplateTypeMap doesn't
	// need cross-type equality; TemplatizedType handles it below).

	// Record type equality.
	leftObj, leftIsObj := left.(*ObjectType)
	rightObj, rightIsObj := right.(*ObjectType)
	if leftIsObj && rightIsObj {
		if leftObj.IsStructural && rightObj.IsStructural {
			return e.areRecordEqual(leftObj, rightObj)
		}
		// Nominal types: compare by pointer identity (same registry instance).
		return leftObj == rightObj
	}

	// TemplatizedType equality.
	leftTpl, leftIsTpl := left.(*TemplatizedType)
	rightTpl, rightIsTpl := right.(*TemplatizedType)
	if leftIsTpl && rightIsTpl {
		return e.areTemplatizedEqual(leftTpl, rightTpl)
	}

	// TemplateType uses identity semantics.
	if _, ok := left.(*TemplateType); ok {
		return left == right
	}
	if _, ok := right.(*TemplateType); ok {
		return left == right
	}

	// For all other types (primitives, etc.) that reach here, use TypeName
	// equality as a proxy for value equality (all primitive types are
	// package-level singletons). This matches Java's reference-equality on
	// interned types.
	return left.TypeName() == right.TypeName()
}

// areUnionEqual checks whether two union types are equal.
// Two union types are equal if, after flattening nested union types, they
// have the same number of alternates and all alternates are equal.
// Port of EqualityChecker.areUnionEqual() in Java.
func (e *EqualityChecker) areUnionEqual(left, right *UnionType) bool {
	leftAlts := left.members
	rightAlts := right.members

	if e.eqMethod != nil && *e.eqMethod == EqMethodIdentity && len(leftAlts) != len(rightAlts) {
		return false
	}

	// For each right alternate, find a matching left alternate.
outer:
	for _, rightAlt := range rightAlts {
		for _, leftAlt := range leftAlts {
			if e.areEqualCaching(leftAlt, rightAlt) {
				continue outer
			}
		}
		return false
	}

	return true
}

// areFunctionEqual checks whether two function types are equal.
// Two function types are equal if their signatures match.
// Port of EqualityChecker.areFunctionEqual() in Java.
func (e *EqualityChecker) areFunctionEqual(left, right *FunctionType) bool {
	if left == right {
		return true
	}

	// Constructors and interfaces use identity semantics.
	if left.IsConstructor != right.IsConstructor {
		return false
	}
	if left.IsConstructor {
		// identity check already done above
		return false
	}

	// Ordinary functions: compare return type and params.
	if !e.areEqualCaching(left.Return, right.Return) {
		return false
	}

	if len(left.Params) != len(right.Params) {
		return false
	}
	for i := range left.Params {
		if !e.areEqualCaching(left.Params[i], right.Params[i]) {
			return false
		}
	}
	return true
}

// areRecordEqual checks equality for structural record types.
// Port of EqualityChecker.areRecordEqual() in Java.
//
// Only inline record literal types can be compared for equality based on
// structure. Don't check inherited properties; checking them is both
// incorrect and slow.
func (e *EqualityChecker) areRecordEqual(left, right *ObjectType) bool {
	if len(left.members) != len(right.members) {
		return false
	}
	for key, leftType := range left.members {
		rightType, ok := right.members[key]
		if !ok {
			return false
		}
		if !e.areEqualCaching(leftType, rightType) {
			return false
		}
	}
	return true
}

// areTemplatizedEqual checks equality for templatized types.
// Port of the templatized-type branch in EqualityChecker.areEqualInternal().
func (e *EqualityChecker) areTemplatizedEqual(left, right *TemplatizedType) bool {
	if left.Base != right.Base {
		return false
	}
	if len(left.Args) != len(right.Args) {
		return false
	}
	for i := range left.Args {
		if !e.areEqualCaching(left.Args[i], right.Args[i]) {
			return false
		}
	}
	return true
}

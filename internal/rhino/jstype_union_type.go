// Port of com/google/javascript/rhino/jstype/UnionType.java @ closure-compiler v20261015

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

import (
	"iter"
	"slices"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// UnionType is a type that may be any one of a set of types, and thus has
// the intersection of the properties of those types.
//
// The UnionType implements a common JavaScript idiom in which the code is
// specifically designed to work with multiple input types. Because
// JavaScript always knows the run-time type of an object value, this is
// safer than a C union.
//
// For instance, values of the union type (String,boolean) can be of type
// String or of type boolean. The commutativity of the statement is captured
// by making (String,boolean) and (boolean,String) equal.
//
// The implementation of this class prevents the creation of nested unions.
type UnionType struct {
	JSTypeBase

	// The types that are unioned.
	//
	// We use a List rather than a Set so that we can iterate by index. This is
	// substantially more efficient than using the for-each idiom due to the
	// overhead of allocating iterators. All loops over this collection should
	// be indexed.
	alternates       []JSType
	inlineAlternates [2]JSType

	// How many alternates were resolved the last time this union was
	// (re)built?
	//
	// Used to skip rebuilds that cannot change anything. -1 means "never
	// rebuilt since construction".
	resolvedAlternatesAtLastBuild int32

	// Were all the types used in the most recent build of this union resolved
	// before that build?
	//
	// If this is true, there is no need to rebuild this union ever again,
	// including during resolution. The result will not change.
	alternatesResolvedBeforeBuild bool
}

// Java: UnionType.TYPE_CLASS
const UnionType_TYPE_CLASS = JSTypeClass_UNION

// UnionType_MAX_UNION_SIZE: generally, if the best we can do is say "this
// object is one of thirty things", then we should just give up and admit
// that we have no clue.
//
// Java: UnionType.MAX_UNION_SIZE
const UnionType_MAX_UNION_SIZE = 30

// NewUnionType creates a union.
//
// This ctor is private because all instances are created using a Builder.
// The builder is also responsible for setting the alternates, which is why
// they aren't passed as a parameter.
//
// Java: UnionType(Builder)
func NewUnionType(builder *UnionType_Builder) *UnionType {
	t := &UnionType{}
	t.objFlags = jsKindUnion
	t.InitJSTypeBase(t, builder.registry)
	t.resolvedAlternatesAtLastBuild = -1
	t.FillFromBuilder(builder)

	t.GetRegistry().GetResolver().ResolveIfClosed(t, UnionType_TYPE_CLASS)
	return t
}

// Java: UnionType.getTypeClass()
func (t *UnionType) GetTypeClass() JSTypeClass {
	return UnionType_TYPE_CLASS
}

// UnionType_Builder_ creates a Builder for a new UnionType.
//
// Java: UnionType.builder(JSTypeRegistry)
func UnionType_Builder_(registry *JSTypeRegistry) *UnionType_Builder {
	return NewUnionType_Builder(registry)
}

// GetAlternates gets the alternate types of this union type. The returned
// set is immutable.
//
// Java: UnionType.getAlternates()
func (t *UnionType) GetAlternates() []JSType {
	if !t.IsResolved() && !t.alternatesResolvedBeforeBuild {
		// Only rebuild if something could actually have changed since the
		// last rebuild.
		//
		// `alternates` is only ever replaced by a rebuild, and an individual
		// alternate only ever goes from unresolved to resolved. So for as long
		// as the list is stable the resolved count can only rise, and a rise
		// is exactly the signal that a rebuild could now reach a different
		// answer.
		//
		// Recording the count *before* rebuilding is what makes this
		// re-entrant-safe: the rebuild runs subtype checks that can read this
		// same union, and `alternates` is not replaced until the rebuild
		// finishes, so the nested read sees no change and returns instead of
		// starting another rebuild.
		resolvedAlternates := int32(t.CountResolvedAlternates())

		if resolvedAlternates != t.resolvedAlternatesAtLastBuild {
			t.resolvedAlternatesAtLastBuild = resolvedAlternates

			b := t.GetRegistry().acquireUnionBuilder()
			b.rebuildTarget = t
			b.AddAlternates(t.alternates)
			// Double checked rebuilds, in case the union is involved in a cycle.
			if !t.IsResolved() && !t.alternatesResolvedBeforeBuild {
				b.Build()

				// The rebuild may have replaced `alternates` with a different
				// list, so the count taken above no longer describes this union.
				// Notably the count can go *down*: an alternate that resolved to
				// a union is flattened into its members, and those members may
				// themselves be unresolved. Re-baseline against the list we
				// actually have now, so the next read compares like with like.
				t.resolvedAlternatesAtLastBuild = int32(t.CountResolvedAlternates())
			}
			t.GetRegistry().releaseUnionBuilder(b)
		}
	}
	return t.alternates
}

// CountResolvedAlternates: how many of the current alternates are resolved?
//
// Java: UnionType.countResolvedAlternates()
func (t *UnionType) CountResolvedAlternates() int {
	count := 0
	for i := 0; i < len(t.alternates); i++ {
		if t.alternates[i].IsResolved() {
			count++
		}
	}
	return count
}

// Java: UnionType.fillFromBuilder(Builder)
func (t *UnionType) FillFromBuilder(builder *UnionType_Builder) {
	jutil.CheckState(!t.alternatesResolvedBeforeBuild)
	if builder.finalAlternates != nil {
		n := len(builder.finalAlternates)
		jutil.CheckState(n != 0)
		if n <= len(t.inlineAlternates) {
			t.inlineAlternates[0] = builder.finalAlternates[0]
			if n == 2 {
				t.inlineAlternates[1] = builder.finalAlternates[1]
			} else {
				t.inlineAlternates[1] = nil
			}
			t.alternates = t.inlineAlternates[:n]
		} else {
			t.alternates = builder.finalAlternates
		}
	} else {
		n := builder.inlineLen
		jutil.CheckState(n != 0)
		if n <= len(t.inlineAlternates) {
			t.inlineAlternates[0] = builder.inlineAlternates[0]
			if n == 2 {
				t.inlineAlternates[1] = builder.inlineAlternates[1]
			} else {
				t.inlineAlternates[1] = nil
			}
			t.alternates = t.inlineAlternates[:n]
		} else {
			t.alternates = slices.Clone(builder.inlineAlternates[:n])
		}
	}
	t.alternatesResolvedBeforeBuild = builder.alternatesResolvedBeforeBuild
}

// MatchesNumberContext is used to test whether a given type can appear in a
// numeric context, such as an operand of a multiply operator.
//
// Java: UnionType.matchesNumberContext()
func (t *UnionType) MatchesNumberContext() bool {
	// TODO(user): Reverse this logic to make it correct instead of generous.
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].MatchesNumberContext() {
			return true
		}
	}
	return false
}

// MatchesStringContext is used to test whether a given type can appear in a
// String context, such as an operand of a string concat (+) operator.
//
// All types have at least the potential for converting to String. When we
// add externally defined types, such as a browser OM, we may choose to add
// types that do not automatically convert to String.
//
// Java: UnionType.matchesStringContext()
func (t *UnionType) MatchesStringContext() bool {
	// TODO(user): Reverse this logic to make it correct instead of generous.
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].MatchesStringContext() {
			return true
		}
	}
	return false
}

// MatchesSymbolContext is used to test whether a given type can appear in a
// Symbol context.
//
// Java: UnionType.matchesSymbolContext()
func (t *UnionType) MatchesSymbolContext() bool {
	// TODO(user): Reverse this logic to make it correct instead of generous.
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].MatchesSymbolContext() {
			return true
		}
	}
	return false
}

// MatchesObjectContext is used to test whether a given type can appear in an
// Object context, such as the expression in a with statement.
//
// Most types we will encounter, except notably null, have at least the
// potential for converting to Object. Host defined objects can get peculiar.
//
// VOID type is included here because while it is not part of the JavaScript
// language, functions returning 'void' type can't be used as operands of
// any operator or statement.
//
// Java: UnionType.matchesObjectContext()
func (t *UnionType) MatchesObjectContext() bool {
	// TODO(user): Reverse this logic to make it correct instead of generous.
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].MatchesObjectContext() {
			return true
		}
	}
	return false
}

// Java: UnionType.findPropertyTypeWithoutConsideringTemplateTypes(Property.Key)
func (t *UnionType) FindPropertyTypeWithoutConsideringTemplateTypes(propertyName Property_Key) JSType {
	var propertyType JSType

	for _, alternate := range t.alternates {
		// Filter out the null/undefined type.
		if alternate.IsNullType() || alternate.IsVoidType() {
			continue
		}

		altPropertyType := alternate.FindPropertyTypeByKey(propertyName)
		if altPropertyType == nil {
			continue
		}

		if propertyType == nil {
			propertyType = altPropertyType
		} else {
			propertyType = propertyType.GetLeastSupertype(altPropertyType)
		}
	}

	return propertyType
}

// Java: UnionType.canBeCalled()
func (t *UnionType) CanBeCalled() bool {
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if !alts[i].CanBeCalled() {
			return false
		}
	}
	return true
}

// Java: UnionType.autobox()
func (t *UnionType) Autobox() JSType {
	b := t.GetRegistry().acquireUnionBuilder()
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		b.AddAlternate(alts[i].Autobox())
	}
	res := b.Build()
	t.GetRegistry().releaseUnionBuilder(b)
	return res
}

// Java: UnionType.restrictByNotNullOrUndefined()
func (t *UnionType) RestrictByNotNullOrUndefined() JSType {
	alts := t.GetAlternates()
	if t.alternatesResolvedBeforeBuild && len(alts) == 2 {
		r0 := alts[0].RestrictByNotNullOrUndefined()
		r1 := alts[1].RestrictByNotNullOrUndefined()
		if r0 == alts[0] && r1.IsNoType() {
			return r0
		}
		if r0.IsNoType() && r1 == alts[1] {
			return r1
		}
		b := t.GetRegistry().acquireUnionBuilder()
		b.AddAlternate(r0)
		b.AddAlternate(r1)
		res := b.Build()
		t.GetRegistry().releaseUnionBuilder(b)
		return res
	}
	b := t.GetRegistry().acquireUnionBuilder()
	for i := 0; i < len(alts); i++ {
		b.AddAlternate(alts[i].RestrictByNotNullOrUndefined())
	}
	res := b.Build()
	t.GetRegistry().releaseUnionBuilder(b)
	return res
}

// Java: UnionType.restrictByNotUndefined()
func (t *UnionType) RestrictByNotUndefined() JSType {
	alts := t.GetAlternates()
	if t.alternatesResolvedBeforeBuild && len(alts) == 2 {
		r0 := alts[0].RestrictByNotUndefined()
		r1 := alts[1].RestrictByNotUndefined()
		if r0 == alts[0] && r1.IsNoType() {
			return r0
		}
		if r0.IsNoType() && r1 == alts[1] {
			return r1
		}
		b := t.GetRegistry().acquireUnionBuilder()
		b.AddAlternate(r0)
		b.AddAlternate(r1)
		res := b.Build()
		t.GetRegistry().releaseUnionBuilder(b)
		return res
	}
	b := t.GetRegistry().acquireUnionBuilder()
	for i := 0; i < len(alts); i++ {
		b.AddAlternate(alts[i].RestrictByNotUndefined())
	}
	res := b.Build()
	t.GetRegistry().releaseUnionBuilder(b)
	return res
}

// Java: UnionType.restrictByNotNull()
func (t *UnionType) RestrictByNotNull() JSType {
	alts := t.GetAlternates()
	if t.alternatesResolvedBeforeBuild && len(alts) == 2 {
		r0 := alts[0].RestrictByNotNull()
		r1 := alts[1].RestrictByNotNull()
		if r0 == alts[0] && r1.IsNoType() {
			return r0
		}
		if r0.IsNoType() && r1 == alts[1] {
			return r1
		}
		b := t.GetRegistry().acquireUnionBuilder()
		b.AddAlternate(r0)
		b.AddAlternate(r1)
		res := b.Build()
		t.GetRegistry().releaseUnionBuilder(b)
		return res
	}
	b := t.GetRegistry().acquireUnionBuilder()
	for i := 0; i < len(alts); i++ {
		b.AddAlternate(alts[i].RestrictByNotNull())
	}
	res := b.Build()
	t.GetRegistry().releaseUnionBuilder(b)
	return res
}

// Java: UnionType.testForEquality(JSType)
func (t *UnionType) TestForEquality(that JSType) base.Tri {
	result := jsType_TRI_NULL
	for i := 0; i < len(t.alternates); i++ {
		alt := t.alternates[i]
		test := alt.TestForEquality(that)
		if result == jsType_TRI_NULL {
			result = test
		} else {
			jsTypeCheckTriNotNull(result)
			if result != test {
				return base.Tri_UNKNOWN
			}
		}
	}
	return result
}

// IsNullable determines whether objects of this type can have the null
// value, and therefore can appear in contexts where null is expected.
//
// Java: UnionType.isNullable()
func (t *UnionType) IsNullable() bool {
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].IsNullable() {
			return true
		}
	}
	return false
}

// IsVoidable tests whether this type is voidable.
//
// Java: UnionType.isVoidable()
func (t *UnionType) IsVoidable() bool {
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].IsVoidable() {
			return true
		}
	}
	return false
}

// IsExplicitlyVoidable tests whether this type explicitly allows undefined
// (as opposed to ? or *).
//
// Java: UnionType.isExplicitlyVoidable()
func (t *UnionType) IsExplicitlyVoidable() bool {
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].IsExplicitlyVoidable() {
			return true
		}
	}
	return false
}

// Java: UnionType.isUnknownType()
func (t *UnionType) IsUnknownType() bool {
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].IsUnknownType() {
			return true
		}
	}
	return false
}

// Java: UnionType.isStruct()
func (t *UnionType) IsStruct() bool {
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].IsStruct() {
			return true
		}
	}
	return false
}

// Java: UnionType.isDict()
func (t *UnionType) IsDict() bool {
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].IsDict() {
			return true
		}
	}
	return false
}

// Java: UnionType.getLeastSupertype(JSType)
func (t *UnionType) GetLeastSupertype(that JSType) JSType {
	if t == that {
		return t
	}
	if !that.IsUnknownType() && !that.IsUnionType() {
		for i := 0; i < len(t.alternates); i++ {
			alternate := t.alternates[i]
			if alternate == that || (!alternate.IsUnknownType() && that.IsSubtypeOf(alternate)) {
				return t
			}
		}
	}

	return JSType_GetLeastSupertype(t, that)
}

// Java: UnionType.getGreatestSubtype(UnionType, JSType)
func UnionType_GetGreatestSubtype(union *UnionType, that JSType) JSType {
	// This method is implemented as a static because we don't want
	// polymorphism. Ideally all the `greatestSubtype` code would be in one
	// place. Until then, using static calls minimizes confusion.

	registry := union.GetRegistry()
	builder := registry.acquireUnionBuilder()

	for i := 0; i < len(union.alternates); i++ {
		alternate := union.alternates[i]
		if alternate.IsSubtypeOf(that) {
			builder.AddAlternate(alternate)
		}
	}

	if that.IsUnionType() {
		thoseAlternates := that.ToMaybeUnionType().GetAlternates()
		for i := 0; i < len(thoseAlternates); i++ {
			otherAlternate := thoseAlternates[i]
			if otherAlternate.IsSubtypeOf(union) {
				builder.AddAlternate(otherAlternate)
			}
		}
	} else if that.IsSubtypeOf(union) {
		builder.AddAlternate(that)
	}

	result := builder.Build()
	registry.releaseUnionBuilder(builder)
	if !result.IsNoType() {
		return result
	} else if union.IsObject() && (that.IsObject() && !that.IsNoType()) {
		return registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE)
	} else {
		return registry.GetNativeType(JSTypeNative_NO_TYPE)
	}
}

// Java: UnionType.getPropertyKind(Property.Key, boolean)
func (t *UnionType) GetPropertyKindByKeyWithAutobox(pname Property_Key, autobox bool) JSType_HasPropertyKind {
	found := false
	always := true
	for i := 0; i < len(t.alternates); i++ {
		alternate := t.alternates[i]
		if alternate.IsNullType() || alternate.IsVoidType() {
			continue
		}
		switch alternate.GetPropertyKindByKeyWithAutobox(pname, autobox) {
		case JSType_HasPropertyKind_KNOWN_PRESENT:
			found = true
		case JSType_HasPropertyKind_ABSENT:
			always = false
		case JSType_HasPropertyKind_MAYBE_PRESENT:
			found = true
			always = false
		}
		if found && !always {
			break
		}
	}
	if found {
		if always {
			return JSType_HasPropertyKind_KNOWN_PRESENT
		}
		return JSType_HasPropertyKind_MAYBE_PRESENT
	}
	return JSType_HasPropertyKind_ABSENT
}

// Java: UnionType.recursionUnsafeHashCode()
func (t *UnionType) RecursionUnsafeHashCode() int32 {
	hashCode := int32(len(t.alternates))
	for i := 0; i < len(t.alternates); i++ {
		// To be determinisitic this aggregation must be order-independent.
		// Using a commutative operatator (multiplication) allows us to achieve
		// that without sorting. Multiplication also has some nice properties
		// about reducing collisions compared to addition or xor.
		hashCode *= t.alternates[i].HashCode()
	}
	return hashCode
}

// Java: JSType.isUnionType()
func (t *UnionType) IsUnionType() bool {
	return true
}

// Java: UnionType.toMaybeUnionType()
func (t *UnionType) ToMaybeUnionType() *UnionType {
	return t
}

// Java: UnionType.isObject()
func (t *UnionType) IsObject() bool {
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if !alts[i].IsObject() {
			return false
		}
	}
	return true
}

// Java: JSType.isObjectType()
func (t *UnionType) IsObjectType() bool {
	return t.IsObject()
}

// Contains: a UnionType contains a given type (alternate) iff the member
// vector contains it.
//
// Java: UnionType.contains(JSType)
func (t *UnionType) Contains(type_ JSType) bool {
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if type_.Equals(alts[i]) {
			return true
		}
	}
	return false
}

// GetRestrictedUnion returns a more restricted union type than this one, in
// which all subtypes of type have been removed.
//
// Examples:
//
//   - (number,string) restricted by number is string
//   - (null, EvalError, URIError) restricted by Error is null
//
// Java: UnionType.getRestrictedUnion(JSType)
func (t *UnionType) GetRestrictedUnion(type_ JSType) JSType {
	restricted := t.GetRegistry().acquireUnionBuilder()
	for i := 0; i < len(t.alternates); i++ {
		alt := t.alternates[i]
		// Keep all unknown/unresolved types.
		if alt.IsUnknownType() || !alt.IsSubtypeOf(type_) {
			restricted.AddAlternate(alt)
		}
	}
	res := restricted.Build()
	t.GetRegistry().releaseUnionBuilder(restricted)
	return res
}

// Java: UnionType.appendTo(TypeStringBuilder)
func (t *UnionType) AppendTo(sb *TypeStringBuilder) {
	sb.Append("(")

	names := []string{}
	for _, alt := range t.alternates {
		// Clone the config to preserve indentation.
		names = append(names, sb.CloneWithConfig().AppendType(alt).Build())
	}
	// Sort types by stringification to get deterministic behaviour.
	slices.SortStableFunc(names, jutil.StringCompareTo)
	elements := make([]any, len(names))
	for i, n := range names {
		elements[i] = n
	}
	sb.AppendAll(elements, "|")

	sb.Append(")")
}

// Java: UnionType.getRestrictedTypeGivenOutcome(Outcome)
func (t *UnionType) GetRestrictedTypeGivenOutcome(outcome Outcome) JSType {
	b := t.GetRegistry().acquireUnionBuilder()
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		b.AddAlternate(alts[i].GetRestrictedTypeGivenOutcome(outcome))
	}
	res := b.Build()
	t.GetRegistry().releaseUnionBuilder(b)
	return res
}

// Java: UnionType.getPossibleToBooleanOutcomes()
func (t *UnionType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	literals := BooleanLiteralSet_EMPTY
	for i := 0; i < len(t.alternates); i++ {
		element := t.alternates[i]
		literals = literals.Union(element.GetPossibleToBooleanOutcomes())
		if literals == BooleanLiteralSet_BOTH {
			break
		}
	}
	return literals
}

// Java: UnionType.getTypesUnderEquality(JSType)
func (t *UnionType) GetTypesUnderEquality(that JSType) *JSType_TypePair {
	resA, resB := t.GetTypesUnderEqualityPair(that)
	return NewJSType_TypePair(resA, resB)
}

func (t *UnionType) GetTypesUnderEqualityPair(that JSType) (JSType, JSType) {
	thisRestricted := t.GetRegistry().acquireUnionBuilder()
	thatRestricted := t.GetRegistry().acquireUnionBuilder()
	for i := 0; i < len(t.alternates); i++ {
		element := t.alternates[i]
		typeA, typeB := element.GetTypesUnderEqualityPair(that)
		if typeA != nil {
			thisRestricted.AddAlternate(typeA)
		}
		if typeB != nil {
			thatRestricted.AddAlternate(typeB)
		}
	}
	resA := thisRestricted.Build()
	resB := thatRestricted.Build()
	t.GetRegistry().releaseUnionBuilder(thatRestricted)
	t.GetRegistry().releaseUnionBuilder(thisRestricted)
	return resA, resB
}

// Java: UnionType.getTypesUnderInequality(JSType)
func (t *UnionType) GetTypesUnderInequality(that JSType) *JSType_TypePair {
	resA, resB := t.GetTypesUnderInequalityPair(that)
	return NewJSType_TypePair(resA, resB)
}

func (t *UnionType) GetTypesUnderInequalityPair(that JSType) (JSType, JSType) {
	thisRestricted := t.GetRegistry().acquireUnionBuilder()
	thatRestricted := t.GetRegistry().acquireUnionBuilder()
	for i := 0; i < len(t.alternates); i++ {
		element := t.alternates[i]
		typeA, typeB := element.GetTypesUnderInequalityPair(that)
		if typeA != nil {
			thisRestricted.AddAlternate(typeA)
		}
		if typeB != nil {
			thatRestricted.AddAlternate(typeB)
		}
	}
	resA := thisRestricted.Build()
	resB := thatRestricted.Build()
	t.GetRegistry().releaseUnionBuilder(thatRestricted)
	t.GetRegistry().releaseUnionBuilder(thisRestricted)
	return resA, resB
}

// Java: UnionType.getTypesUnderShallowInequality(JSType)
func (t *UnionType) GetTypesUnderShallowInequality(that JSType) *JSType_TypePair {
	resA, resB := t.GetTypesUnderShallowInequalityPair(that)
	return NewJSType_TypePair(resA, resB)
}

func (t *UnionType) GetTypesUnderShallowInequalityPair(that JSType) (JSType, JSType) {
	thisRestricted := t.GetRegistry().acquireUnionBuilder()
	thatRestricted := t.GetRegistry().acquireUnionBuilder()
	for i := 0; i < len(t.alternates); i++ {
		element := t.alternates[i]
		typeA, typeB := element.GetTypesUnderShallowInequalityPair(that)
		if typeA != nil {
			thisRestricted.AddAlternate(typeA)
		}
		if typeB != nil {
			thatRestricted.AddAlternate(typeB)
		}
	}
	resA := thisRestricted.Build()
	resB := thatRestricted.Build()
	t.GetRegistry().releaseUnionBuilder(thatRestricted)
	t.GetRegistry().releaseUnionBuilder(thisRestricted)
	return resA, resB
}

// Java: UnionType.visit(Visitor<T>)
func (t *UnionType) Visit(visitor Visitor) any {
	return visitor.CaseUnionType(t)
}

// Java: UnionType.visit(RelationshipVisitor<T>, JSType)
func (t *UnionType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseUnionType(t, that)
}

// Java: UnionType.resolveInternal(ErrorReporter)
func (t *UnionType) ResolveInternal(reporter ErrorReporter) JSType {
	if t.alternatesResolvedBeforeBuild {
		return t
	}

	for i := 0; i < len(t.alternates); i++ {
		t.alternates[i].Resolve(reporter)
	}
	b := t.GetRegistry().acquireUnionBuilder()
	b.rebuildTarget = t
	res := b.AddAlternates(t.alternates).Build()
	t.GetRegistry().releaseUnionBuilder(b)
	return res
}

// Java: UnionType.collapseUnion()
func (t *UnionType) CollapseUnion() JSType {
	var currentValue JSType
	var currentCommonSuper ObjectType
	for i := 0; i < len(t.alternates); i++ {
		a := t.alternates[i]
		if a.IsUnknownType() {
			return t.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}

		obj := a.ToObjectType()
		if obj == nil {
			if currentValue == nil && currentCommonSuper == nil {
				// If obj is not an object, then it must be a value.
				currentValue = a
			} else {
				// Multiple values and objects will always collapse to the ALL_TYPE.
				return t.GetNativeType(JSTypeNative_ALL_TYPE)
			}
		} else if currentValue != nil {
			// Values and objects will always collapse to the ALL_TYPE.
			return t.GetNativeType(JSTypeNative_ALL_TYPE)
		} else if currentCommonSuper == nil {
			currentCommonSuper = obj
		} else {
			currentCommonSuper = t.GetRegistry().FindCommonSuperObject(currentCommonSuper, obj)
		}
	}
	return jsTypeOrNil(currentCommonSuper)
}

// Java: UnionType.matchConstraint(JSType)
func (t *UnionType) MatchConstraint(constraint JSType) {
	for i := 0; i < len(t.alternates); i++ {
		alternate := t.alternates[i]
		alternate.MatchConstraint(constraint)
	}
}

// Java: UnionType.hasAnyTemplateTypesInternal()
func (t *UnionType) HasAnyTemplateTypesInternal() bool {
	alts := t.GetAlternates()
	for i := 0; i < len(alts); i++ {
		if alts[i].HasAnyTemplateTypes() {
			return true
		}
	}
	return false
}

// UnionType_Builder implements type unioning logic, since UnionTypes only
// actually need to perform unioning operations when being (re)built.
//
// Builders exist in two forms. One for assembing a new union and one for
// updating an existing union. Only the former is exposed.
//
// Most users of this class should prefer JSTypeRegistry#createUnionType
// instead.
type UnionType_Builder struct {
	rebuildTarget *UnionType
	registry      *JSTypeRegistry

	inlineAlternates [32]JSType
	inlineLen        int
	alternates       []JSType
	finalAlternates  []JSType

	// If a union has ? or *, we do not care about any other types, except for
	// undefined (for optional properties).
	containsVoidType              bool
	isAllType                     bool
	isNativeUnknownType           bool
	areAllUnknownsChecked         bool
	alternatesResolvedBeforeBuild bool
	hasBuilt                      bool

	// Every UnionType may have at most one structural function in it.
	//
	// NOTE(nicksantos): I've read some literature that says that
	// type-inferenced languages are fundamentally incompatible with union
	// types. I refuse to believe this. But they do make the type lattice much
	// more complicated.
	//
	// For this reason, when we deal with function types, we actually merge
	// some nodes on the lattice, and treat them as fundamentally equivalent.
	// For example, we treat
	//
	//	function(): string | function(): number
	//
	// as equivalent to
	//
	//	function(): (string|number)
	//
	// and normalize the first type into the second type.
	//
	// To perform this normalization, we've modified Builder to disallow
	// multiple structural functions in a union. We always delegate to
	// FunctionType::getLeastSupertype, which either merges the functions into
	// one structural function, or just bails out and uses the top function
	// type.
	functionTypePosition int
}

func (b *UnionType_Builder) size() int {
	if b.alternates != nil {
		return len(b.alternates)
	}
	return b.inlineLen
}

func (b *UnionType_Builder) get(i int) JSType {
	if b.alternates != nil {
		return b.alternates[i]
	}
	return b.inlineAlternates[i]
}

func (b *UnionType_Builder) set(i int, t JSType) {
	if b.alternates != nil {
		b.alternates[i] = t
	} else {
		b.inlineAlternates[i] = t
	}
}

func (b *UnionType_Builder) removeAt(i int) {
	if b.alternates == nil {
		last := b.inlineLen - 1
		copy(b.inlineAlternates[i:last], b.inlineAlternates[i+1:b.inlineLen])
		b.inlineAlternates[last] = nil
		b.inlineLen = last
		return
	}
	last := len(b.alternates) - 1
	copy(b.alternates[i:], b.alternates[i+1:])
	b.alternates[last] = nil
	b.alternates = b.alternates[:last]
}

func (b *UnionType_Builder) add(t JSType) {
	if b.alternates != nil {
		b.alternates = append(b.alternates, t)
		return
	}
	if b.inlineLen < len(b.inlineAlternates) {
		b.inlineAlternates[b.inlineLen] = t
		b.inlineLen++
		return
	}
	b.alternates = make([]JSType, len(b.inlineAlternates)+1, 64)
	copy(b.alternates, b.inlineAlternates[:])
	b.alternates[len(b.inlineAlternates)] = t
	clear(b.inlineAlternates[:])
	b.inlineLen = 0
}

func initUnionType_Builder(b *UnionType_Builder, registry *JSTypeRegistry) {
	b.rebuildTarget = nil
	b.registry = registry
	b.inlineLen = 0
	b.alternates = nil
	b.finalAlternates = nil
	b.containsVoidType = false
	b.isAllType = false
	b.isNativeUnknownType = false
	b.areAllUnknownsChecked = true
	b.alternatesResolvedBeforeBuild = true
	b.hasBuilt = false
	b.functionTypePosition = -1
}

// NewUnionType_Builder creates a builder for a new union.
//
// Java: UnionType.Builder(JSTypeRegistry)
func NewUnionType_Builder(registry *JSTypeRegistry) *UnionType_Builder {
	return &UnionType_Builder{
		registry:                      registry,
		areAllUnknownsChecked:         true,
		alternatesResolvedBeforeBuild: true,
		functionTypePosition:          -1,
	}
}

// NewUnionType_BuilderOfRebuildTarget creates a re-builder for an existing
// union.
//
// Java: UnionType.Builder(UnionType)
func NewUnionType_BuilderOfRebuildTarget(rebuildTarget *UnionType) *UnionType_Builder {
	b := NewUnionType_Builder(rebuildTarget.GetRegistry())
	b.rebuildTarget = rebuildTarget
	return b
}

// Java: UnionType.Builder.isSubtype(JSType, JSType)
func UnionType_Builder_IsSubtype(rightType JSType, leftType JSType) bool {
	return rightType.IsSubtypeWithoutStructuralTyping(leftType)
}

// Java: UnionType.Builder.addAlternates(Collection<? extends JSType>)
func (b *UnionType_Builder) AddAlternatesOfCollection(c iter.Seq[JSType]) *UnionType_Builder {
	for type_ := range c {
		b.AddAlternate(type_)
	}
	return b
}

// AddAlternates: a specific override that avoid creating an iterator. This
// version is currently used when adding a union as an alternate.
//
// Java: UnionType.Builder.addAlternates(List<? extends JSType>)
func (b *UnionType_Builder) AddAlternates(list []JSType) *UnionType_Builder {
	for i := 0; i < len(list); i++ {
		b.AddAlternate(list[i])
	}
	return b
}

// AddAlternate adds an alternate to the union type under construction.
//
// Returns this for easy chaining.
//
// Java: UnionType.Builder.addAlternate(JSType)
func (b *UnionType_Builder) AddAlternate(alternate JSType) *UnionType_Builder {
	b.CheckHasNotBuilt()

	altBase := jsTypeBase(alternate)
	altKind := altBase.objFlags & objKindMask
	if altKind == jsKindUnion {
		b.AddAlternates((*UnionType)(unsafe.Pointer(altBase)).GetAlternates())
		return b
	}
	if ((uint32(1)<<altKind)&jsKindProxyOrFallbackMask) != 0 && alternate.IsUnionType() {
		b.AddAlternates(alternate.ToMaybeUnionType().GetAlternates())
		return b
	}

	if b.size() > UnionType_MAX_UNION_SIZE {
		return b
	}

	// Defer removing duplicate elements until all alternates in the union are
	// resolved. JSType operations like equality and subtyping are not reliable
	// pre-resolution.
	if altBase.resolveResult == nil {
		b.alternatesResolvedBeforeBuild = false

		for i := 0; i < b.size(); i++ {
			if jsTypeBase(b.get(i)) == altBase {
				return b
			}
		}
		b.add(alternate)
		return b
	}

	// build() returns the bottom type by default, so we can just bail out
	// early here.
	if altKind == objKindNo || (((uint32(1)<<altKind)&jsKindProxyOrFallbackMask) != 0 && alternate.IsNoType()) {
		return b
	}

	if !b.isAllType && (altKind == jsKindAll || (((uint32(1)<<altKind)&jsKindProxyOrFallbackMask) != 0 && alternate.IsAllType())) {
		b.isAllType = true
	}
	if !b.containsVoidType && (altKind == jsKindVoid || (altKind == objKindFallback && alternate.IsVoidType())) {
		b.containsVoidType = true
	}

	if altKind == objKindUnknown {
		b.isNativeUnknownType = true
		b.areAllUnknownsChecked = b.areAllUnknownsChecked && (*UnknownType)(unsafe.Pointer(altBase)).isChecked
	}

	if b.isAllType || b.isNativeUnknownType {
		return b
	}

	// Function types are special, because they have their own bizarre
	// sub-lattice. See the comments on FunctionType#supAndInf helper and above
	// at functionTypePosition.
	altIsFunc := jsTypeIsFunction(alternate, altKind)
	if altIsFunc && b.functionTypePosition != -1 {
		// See the comments on functionTypePosition above.
		other := b.get(b.functionTypePosition).ToMaybeFunctionType()
		supremum := alternate.ToMaybeFunctionType().SupAndInfHelper(other, true)
		b.set(b.functionTypePosition, supremum)
		return b
	}

	altUnknown := jsTypeIsUnknown(alternate, altBase, altKind)
	altHasTemplate := !altUnknown && jsTypeHasAnyTemplateTypes(altBase)
	altTemplatized := !altUnknown && !altHasTemplate && jsTypeIsTemplatized(alternate, altKind)

	// Look through the alternates we've got so far, and check if any of them
	// are duplicates of one another.
	for index := 0; index < b.size(); index++ {
		removeCurrent := false
		current := b.get(index)
		curBase := jsTypeBase(current)
		if curBase.resolveResult == nil { // Defer de-duplicating unresolved alternates.
			continue
		}
		if curBase == altBase {
			return b
		}
		curKind := curBase.objFlags & objKindMask

		// Unknown types may just be names that haven't been resolved yet. So
		// keep these in the union, and just use equality checking for simple
		// de-duping.
		if altUnknown ||
			jsTypeIsUnknown(current, curBase, curKind) ||
			altHasTemplate ||
			jsTypeHasAnyTemplateTypes(curBase) {
			if equalityChecker_AreIdentityEqualWithBases(alternate, altBase, current, curBase) {
				// Alternate is unnecessary.
				return b
			}
		} else if curTemplatized := jsTypeIsTemplatized(current, curKind); altTemplatized || curTemplatized {
			// Because "Foo" and "Foo<?>" are roughly equivalent templatized types,
			// special care is needed when building the union. For example:
			//   Object is consider a subtype of Object<string>
			// but we want to leave "Object" not "Object<string>" when building the
			// subtype.
			//

			// Cases:
			// 1) alternate:Array<string> and current:Object ==> Object
			// 2) alternate:Array<string> and current:Array ==> Array
			// 3) alternate:Object<string> and
			//    current:Array ==> Array|Object<string>
			// 4) alternate:Object and current:Array<string> ==> Object
			// 5) alternate:Array and current:Array<string> ==> Array
			// 6) alternate:Array and
			//    current:Object<string> ==> Array|Object<string>
			// 7a) alternate:Object<string> and
			//     current:Object<number> ==> Object<?>
			// 7b) alternate:ReadonlyArray<string> and
			//     current:ReadonlyArray<number> ==> ReadonlyArray<string>|ReadonlyArray<number>
			// 8) alternate:Array<string> and
			//    current:Array<string> ==> Array<string>
			// 9) alternate:Array<string> and
			//    current:Object<string> ==> Object<string>|Array<string>

			if !curTemplatized {
				if UnionType_Builder_IsSubtype(alternate, current) {
					// case 1, 2
					return b
				}
				// case 3: leave current, add alternate
			} else if !altTemplatized {
				if UnionType_Builder_IsSubtype(current, alternate) {
					// case 4, 5
					removeCurrent = true
				}
				// case 6: leave current, add alternate
			} else {
				templatizedAlternate := alternate.ToMaybeTemplatizedType()
				templatizedCurrent := current.ToMaybeTemplatizedType()

				if templatizedCurrent.WrapsSameRawType(templatizedAlternate) {

					if equalityChecker_AreIdentityEqualWithBases(current, curBase, alternate, altBase) {
						// case 8
						return b
					} else if !templatizedCurrent.GetReferencedType().IsReadonlyArrayType() {
						// case 7a: replace with a merged alternate specialized on `?`.
						rawType := templatizedCurrent.GetReferencedObjTypeInternal()
						// Providing no type-parameter values specializes `rawType` on `?`
						// by default.
						alternate = b.registry.CreateTemplatizedTypeOfList(rawType, nil)
						altBase = jsTypeBase(alternate)
						altKind = altBase.objFlags & objKindMask
						altIsFunc = jsTypeIsFunction(alternate, altKind)
						altUnknown = jsTypeIsUnknown(alternate, altBase, altKind)
						altHasTemplate = !altUnknown && jsTypeHasAnyTemplateTypes(altBase)
						altTemplatized = !altUnknown && !altHasTemplate && jsTypeIsTemplatized(alternate, altKind)
						removeCurrent = true
					}
				}
				// case 9 & 7b: leave current, add alternate
			}
			// Otherwise leave both templatized types.
		} else {
			bothBits := (uint32(1) << altKind) | (uint32(1) << curKind)
			if (bothBits&jsKindSimpleDisjointValueMask) == 0 ||
				(bothBits&jsKindCanRelateToValueMask) != 0 {
				if UnionType_Builder_IsSubtype(alternate, current) {
					// Alternate is unnecessary.
					b.MayRegisterDroppedProperties(alternate, current)
					return b
				} else if UnionType_Builder_IsSubtype(current, alternate) {
					// Alternate makes current obsolete
					b.MayRegisterDroppedProperties(current, alternate)
					removeCurrent = true
				}
			}
		}

		if removeCurrent {
			b.removeAt(index)

			if index == b.functionTypePosition {
				b.functionTypePosition = -1
			} else if index < b.functionTypePosition {
				b.functionTypePosition--
			}

			index--
		}
	}

	if altIsFunc {
		// See the comments on functionTypePosition above.
		jutil.CheckState(b.functionTypePosition == -1)
		b.functionTypePosition = b.size()
	}

	b.add(alternate)
	return b
}

// Java: UnionType.Builder.mayRegisterDroppedProperties(JSType, JSType)
func (b *UnionType_Builder) MayRegisterDroppedProperties(subtype JSType, supertype JSType) {
	subKind := jsTypeBase(subtype).objFlags & objKindMask
	supKind := jsTypeBase(supertype).objFlags & objKindMask
	if ((uint32(1) << subKind) & jsKindMaybeRecordMask) == 0 ||
		((uint32(1) << supKind) & jsKindMaybeRecordMask) == 0 {
		return
	}
	if subRec := subtype.ToMaybeRecordType(); subRec != nil {
		if supRec := supertype.ToMaybeRecordType(); supRec != nil && subRec != supRec {
			b.registry.RegisterDroppedPropertiesInUnion(subRec, supRec)
		}
	}
}

// Build returns a type, not necessarily a UnionType, that represents the
// union of the inputs.
//
// The Builder cannot be used again once this method is called.
//
// Java: UnionType.Builder.build()
func (b *UnionType_Builder) Build() JSType {
	if b.rebuildTarget == nil {
		b.CheckHasNotBuilt()
		b.hasBuilt = true

		wildcard := b.GetNativeWildcardType()
		if wildcard != nil {
			if !b.containsVoidType {
				return wildcard
			}
			b.inlineAlternates[0] = wildcard
			b.inlineAlternates[1] = b.registry.GetNativeType(JSTypeNative_VOID_TYPE)
			b.inlineLen = 2
			b.alternates = nil
			b.finalAlternates = nil
			return NewUnionType(b)
		}
		n := b.size()
		if n == 0 {
			return b.registry.GetNativeType(JSTypeNative_NO_TYPE)
		}
		if n == 1 {
			return b.get(0)
		}
		if n > UnionType_MAX_UNION_SIZE {
			return b.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
		if b.alternates != nil {
			b.finalAlternates = b.alternates
		}
		return NewUnionType(b)
	}

	b.BuildInternal()
	b.rebuildTarget.FillFromBuilder(b)
	if len(b.rebuildTarget.alternates) == 1 {
		return b.rebuildTarget.alternates[0]
	}
	return b.rebuildTarget
}

// BuildInternal creates the final set of alternates for either a new union
// or a union being rebuilt.
//
// Java: UnionType.Builder.buildInternal()
func (b *UnionType_Builder) BuildInternal() *UnionType_Builder {
	b.CheckHasNotBuilt()
	b.hasBuilt = true

	wildcard := b.GetNativeWildcardType()
	if wildcard != nil {
		if b.containsVoidType {
			b.inlineAlternates[0] = wildcard
			b.inlineAlternates[1] = b.registry.GetNativeType(JSTypeNative_VOID_TYPE)
			b.inlineLen = 2
		} else {
			b.inlineAlternates[0] = wildcard
			b.inlineLen = 1
		}
		b.alternates = nil
		b.finalAlternates = nil
	} else if b.size() == 0 {
		// To simplify the typesystem, empty union types are forbidden. Using a
		// single `bottom` makes it essentially a proxy instead.
		b.inlineAlternates[0] = b.registry.GetNativeType(JSTypeNative_NO_TYPE)
		b.inlineLen = 1
		b.alternates = nil
		b.finalAlternates = nil
	} else if b.size() > UnionType_MAX_UNION_SIZE {
		b.inlineAlternates[0] = b.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		b.inlineLen = 1
		b.alternates = nil
		b.finalAlternates = nil
	} else if b.alternates != nil {
		b.finalAlternates = b.alternates
	}

	return b
}

// GetNativeWildcardType returns ALL_TYPE, UNKNOWN_TYPE,
// CHECKED_UNKNOWN_TYPE, or null as specified by the flags.
//
// Java: UnionType.Builder.getNativeWildcardType()
func (b *UnionType_Builder) GetNativeWildcardType() JSType {
	if b.isAllType {
		return b.registry.GetNativeType(JSTypeNative_ALL_TYPE)
	} else if b.isNativeUnknownType {
		if b.areAllUnknownsChecked {
			return b.registry.GetNativeType(JSTypeNative_CHECKED_UNKNOWN_TYPE)
		} else {
			return b.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
	}
	return nil
}

// Java: UnionType.Builder.checkHasNotBuilt()
func (b *UnionType_Builder) CheckHasNotBuilt() {
	if b.hasBuilt {
		panic(jutil.NewIllegalStateExceptionMsg("Cannot reuse a `UnionType.Builder` that has already built."))
	}
}

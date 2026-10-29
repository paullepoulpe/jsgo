// Port of com/google/javascript/rhino/jstype/UnionType.java @ closure-compiler v20261021

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
	"slices"

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
	inlineAlternates [3]JSType
	heapAlternates   *[]JSType
	inlineLen        uint8

	// How many alternates were resolved the last time this union was
	// (re)built?
	//
	// Used to skip rebuilds that cannot change anything. -1 means "never
	// rebuilt since construction".
	resolvedAlternatesAtLastBuild int8

	// Were all the types used in the most recent build of this union resolved
	// before that build?
	//
	// If this is true, there is no need to rebuild this union ever again,
	// including during resolution. The result will not change.
	alternatesResolvedBeforeBuild bool
}

func (t *UnionType) rawAlternates() []JSType {
	if t.heapAlternates != nil {
		return *t.heapAlternates
	}
	return t.inlineAlternates[:t.inlineLen]
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
		resolvedAlternates := int8(t.CountResolvedAlternates())

		if resolvedAlternates != t.resolvedAlternatesAtLastBuild {
			t.resolvedAlternatesAtLastBuild = resolvedAlternates

			b := t.GetRegistry().acquireUnionBuilder()
			b.rebuildTarget = t
			b.AddAlternates(t.rawAlternates())
			// Double checked rebuilds, in case the union is involved in a cycle.
			if !t.IsResolved() && !t.alternatesResolvedBeforeBuild {
				b.Build()

				// The rebuild may have replaced `alternates` with a different
				// list, so the count taken above no longer describes this union.
				// Notably the count can go *down*: an alternate that resolved to
				// a union is flattened into its members, and those members may
				// themselves be unresolved. Re-baseline against the list we
				// actually have now, so the next read compares like with like.
				t.resolvedAlternatesAtLastBuild = int8(t.CountResolvedAlternates())
			}
			t.GetRegistry().releaseUnionBuilder(b)
		}
	}
	return t.rawAlternates()
}

// CountResolvedAlternates: how many of the current alternates are resolved?
//
// Java: UnionType.countResolvedAlternates()
func (t *UnionType) CountResolvedAlternates() int {
	alts := t.rawAlternates()
	count := 0
	for i := 0; i < len(alts); i++ {
		if alts[i].IsResolved() {
			count++
		}
	}
	return count
}

// Java: UnionType.fillFromBuilder(Builder)
func (t *UnionType) FillFromBuilder(builder *UnionType_Builder) {
	jutil.CheckState(!t.alternatesResolvedBeforeBuild)
	src := builder.finalAlternates
	owned := src != nil
	if !owned {
		src = builder.inlineAlternates[:builder.inlineLen]
	}
	jutil.CheckState(len(src) != 0)
	t.inlineAlternates = [3]JSType{}
	if len(src) <= len(t.inlineAlternates) {
		copy(t.inlineAlternates[:], src)
		t.inlineLen = uint8(len(src))
		t.heapAlternates = nil
	} else {
		if !owned {
			src = slices.Clone(src)
		}
		t.inlineLen = 0
		t.heapAlternates = &src
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

	for _, alternate := range t.rawAlternates() {
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
	alts := t.rawAlternates()
	for i := 0; i < len(alts); i++ {
		alt := alts[i]
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
		alts := t.rawAlternates()
		for i := 0; i < len(alts); i++ {
			alternate := alts[i]
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

	unionAlts := union.rawAlternates()
	for i := 0; i < len(unionAlts); i++ {
		alternate := unionAlts[i]
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
	alts := t.rawAlternates()
	for i := 0; i < len(alts); i++ {
		alternate := alts[i]
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
	alts := t.rawAlternates()
	hashCode := int32(len(alts))
	for i := 0; i < len(alts); i++ {
		// To be determinisitic this aggregation must be order-independent.
		// Using a commutative operatator (multiplication) allows us to achieve
		// that without sorting. Multiplication also has some nice properties
		// about reducing collisions compared to addition or xor.
		hashCode *= alts[i].HashCode()
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
	alts := t.rawAlternates()
	for i := 0; i < len(alts); i++ {
		alt := alts[i]
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

	alts := t.rawAlternates()
	var namesBuf [8]string
	names := namesBuf[:0]
	if len(alts) > len(namesBuf) {
		names = make([]string, 0, len(alts))
	}
	var clone TypeStringBuilder
	for _, alt := range alts {
		// Clone the config to preserve indentation.
		clone.builder.Reset()
		clone.isForAnnotations = sb.isForAnnotations
		clone.indentDepth = sb.indentDepth
		clone.AppendType(alt)
		names = append(names, clone.Build())
	}
	// Sort types by stringification to get deterministic behaviour.
	slices.SortStableFunc(names, jutil.StringCompareTo)
	sb.AppendStrings(names, "|")

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
	alts := t.rawAlternates()
	for i := 0; i < len(alts); i++ {
		element := alts[i]
		literals = literals.Union(element.GetPossibleToBooleanOutcomes())
		if literals == BooleanLiteralSet_BOTH {
			break
		}
	}
	return literals
}

// Java: UnionType.getTypesUnderEquality(JSType)
func (t *UnionType) GetTypesUnderEqualityPair(that JSType) (JSType, JSType) {
	thisRestricted := t.GetRegistry().acquireUnionBuilder()
	thatRestricted := t.GetRegistry().acquireUnionBuilder()
	alts := t.rawAlternates()
	for i := 0; i < len(alts); i++ {
		element := alts[i]
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
func (t *UnionType) GetTypesUnderInequalityPair(that JSType) (JSType, JSType) {
	thisRestricted := t.GetRegistry().acquireUnionBuilder()
	thatRestricted := t.GetRegistry().acquireUnionBuilder()
	alts := t.rawAlternates()
	for i := 0; i < len(alts); i++ {
		element := alts[i]
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
func (t *UnionType) GetTypesUnderShallowInequalityPair(that JSType) (JSType, JSType) {
	thisRestricted := t.GetRegistry().acquireUnionBuilder()
	thatRestricted := t.GetRegistry().acquireUnionBuilder()
	alts := t.rawAlternates()
	for i := 0; i < len(alts); i++ {
		element := alts[i]
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

// Java: UnionType.resolveInternal(ErrorReporter)
func (t *UnionType) ResolveInternal(reporter ErrorReporter) JSType {
	if t.alternatesResolvedBeforeBuild {
		return t
	}

	alts := t.rawAlternates()
	for i := 0; i < len(alts); i++ {
		alts[i].Resolve(reporter)
	}
	b := t.GetRegistry().acquireUnionBuilder()
	b.rebuildTarget = t
	res := b.AddAlternates(t.rawAlternates()).Build()
	t.GetRegistry().releaseUnionBuilder(b)
	return res
}

// Java: UnionType.collapseUnion()
func (t *UnionType) CollapseUnion() JSType {
	var currentValue JSType
	var currentCommonSuper ObjectType
	alts := t.rawAlternates()
	for i := 0; i < len(alts); i++ {
		a := alts[i]
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
	alts := t.rawAlternates()
	for i := 0; i < len(alts); i++ {
		alternate := alts[i]
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

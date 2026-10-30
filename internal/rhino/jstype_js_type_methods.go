// Port of com/google/javascript/rhino/jstype/JSType.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IsStruct returns true iff this can be a struct. UnionType overrides the
// method, assume this is not a union here.
//
// Java: JSType.isStruct()
func (t *JSTypeBase) IsStruct() bool {
	if t.self().IsObject() {
		objType := t.self().ToObjectType()
		ctor := objType.GetConstructor()
		// This test is true for object literals
		if ctor == nil {
			info := objType.GetJSDocInfo()
			return info != nil && info.MakesStructs()
		} else {
			return ctor.MakesStructs()
		}
	}
	return false
}

// IsDict returns true iff this can be a dict. UnionType overrides the
// method, assume this is not a union here.
//
// Java: JSType.isDict()
func (t *JSTypeBase) IsDict() bool {
	if t.self().IsObject() {
		objType := t.self().ToObjectType()
		ctor := objType.GetConstructor()
		// This test is true for object literals
		if ctor == nil {
			info := objType.GetJSDocInfo()
			return info != nil && info.MakesDicts()
		} else {
			return ctor.MakesDicts()
		}
	}
	return false
}

// ContainsReferenceAncestor searches for a type `target` in the reference
// chain of ProxyObjectTypes.
//
// Java: JSType.containsReferenceAncestor(JSType)
func (t *JSTypeBase) ContainsReferenceAncestor(target JSType) bool {
	return containsUpperBoundSuperTypePresent(t.self(), target)
}

// Java: JSType.isLiteralObject()
func (t *JSTypeBase) IsLiteralObject() bool {
	return (t.protoFlags & protoFlagAnonymous) != 0
}

// ToMaybeUnionType downcasts this to a UnionType, or returns null if this is
// not a UnionType.
//
// Named in honor of Haskell's Maybe type constructor.
//
// Java: JSType.toMaybeUnionType()
func (t *JSTypeBase) ToMaybeUnionType() *UnionType {
	return nil
}

// IsGlobalThisType returns true if this is a global this type.
//
// Java: JSType.isGlobalThisType()
func (t *JSTypeBase) IsGlobalThisType() bool {
	return t.self() == t.GetRegistry().GetNativeType(JSTypeNative_GLOBAL_THIS)
}

// IsFunctionType returns true if toMaybeFunctionType returns a non-null
// FunctionType.
//
// Java: JSType.isFunctionType()
func (t *JSTypeBase) IsFunctionType() bool {
	return false
}

// ToMaybeFunctionType downcasts this to a FunctionType, or returns null if
// this is not a function.
//
// For the purposes of this function, we define a MaybeFunctionType as any
// type in the sub-lattice { x | LEAST_FUNCTION_TYPE <= x <=
// GREATEST_FUNCTION_TYPE } This definition excludes bottom types like NoType
// and NoObjectType.
//
// This definition is somewhat arbitrary and axiomatic, but this is the
// definition that makes the most sense for the most callers.
//
// Java: JSType.toMaybeFunctionType()
func (t *JSTypeBase) ToMaybeFunctionType() FunctionType {
	return nil
}

// AssertFunctionType returns this object cast to FunctionType or throws an
// exception if it isn't a FunctionType.
//
// Java: JSType.assertFunctionType()
func (t *JSTypeBase) AssertFunctionType() FunctionType {
	result := t.self().ToMaybeFunctionType()
	if result == nil {
		jutil.CheckNotNullMsg1(result, "not a FunctionType: %s", t.self())
	}
	return result
}

// AssertObjectType returns this object cast to ObjectType or throws an
// exception if it isn't an ObjectType.
//
// Java: JSType.assertObjectType()
func (t *JSTypeBase) AssertObjectType() ObjectType {
	result := t.self().ToMaybeObjectType()
	if result == nil {
		jutil.CheckNotNullMsg1(result, "Not an ObjectType: %s", t.self())
	}
	return result
}

// JSType_ToMaybeFunctionType is a null-safe version of
// toMaybeFunctionType().
//
// Java: JSType.toMaybeFunctionType(JSType)
func JSType_ToMaybeFunctionType(type_ JSType) FunctionType {
	if type_ == nil {
		return nil
	}
	return type_.ToMaybeFunctionType()
}

// Java: JSType.toMaybeKnownSymbolType()
func (t *JSTypeBase) ToMaybeKnownSymbolType() *KnownSymbolType {
	return nil
}

// Java: JSType.isEnumElementType()
func (t *JSTypeBase) IsEnumElementType() bool {
	return false
}

// Java: JSType.getEnumeratedTypeOfEnumElement()
func (t *JSTypeBase) GetEnumeratedTypeOfEnumElement() JSType {
	e := t.self().ToMaybeEnumElementType()
	if e == nil {
		return nil
	}
	return e.GetPrimitiveType()
}

// ToMaybeEnumElementType downcasts this to an EnumElementType, or returns
// null if this is not an EnumElementType.
//
// Java: JSType.toMaybeEnumElementType()
func (t *JSTypeBase) ToMaybeEnumElementType() *EnumElementType {
	return nil
}

// Java: JSType.isEnumType()
func (t *JSTypeBase) IsEnumType() bool {
	// TODO(sdh): Consider changing this to isEnumObjectType(), though this
	// would be inconsistent with the EnumType class and toMaybeEnumType(), so
	// we would need to consider changing them, too.
	return false
}

// ToMaybeEnumType downcasts this to an EnumType, or returns null if this is
// not an EnumType.
//
// Java: JSType.toMaybeEnumType()
func (t *JSTypeBase) ToMaybeEnumType() *EnumType {
	return nil
}

// Java: JSType.isNamedType()
func (t *JSTypeBase) IsNamedType() bool {
	return false
}

// Java: JSType.toMaybeNamedType()
func (t *JSTypeBase) ToMaybeNamedType() *NamedType {
	return nil
}

// Java: JSType.isRecordType()
func (t *JSTypeBase) IsRecordType() bool {
	return false
}

// Java: JSType.isStructuralInterface()
func (t *JSTypeBase) IsStructuralInterface() bool {
	return false
}

// Java: JSType.isStructuralType()
func (t *JSTypeBase) IsStructuralType() bool {
	return false
}

// ToMaybeRecordType downcasts this to a RecordType, or returns null if this
// is not a RecordType.
//
// Java: JSType.toMaybeRecordType()
func (t *JSTypeBase) ToMaybeRecordType() *RecordType {
	return nil
}

// Java: JSType.isTemplatizedType()
func (t *JSTypeBase) IsTemplatizedType() bool {
	return false
}

// ToMaybeTemplatizedType downcasts this to a TemplatizedType, or returns
// null if this is not a templatized type.
//
// Java: JSType.toMaybeTemplatizedType()
func (t *JSTypeBase) ToMaybeTemplatizedType() *TemplatizedType {
	return nil
}

// Java: JSType.isTemplateType()
func (t *JSTypeBase) IsTemplateType() bool {
	return false
}

// ToMaybeTemplateType downcasts this to a TemplateType, or returns null if
// this is not a template type.
//
// Java: JSType.toMaybeTemplateType()
func (t *JSTypeBase) ToMaybeTemplateType() *TemplateType {
	return nil
}

// Java: JSType.hasAnyTemplateTypes()
func (t *JSTypeBase) HasAnyTemplateTypes() bool {
	var st JSTypeStack
	return t.hasAnyTemplateTypesIn(&st)
}

// jsTypeHasAnyTemplateTypesIn is x.HasAnyTemplateTypes() inside a
// HasAnyTemplateTypes computation.
func jsTypeHasAnyTemplateTypesIn(x JSType, st *JSTypeStack) bool {
	return x.GetJSTypeBase().hasAnyTemplateTypesIn(st)
}

// Go-only: the types in progress are kept on the stack of the outermost
// call, not as a flag on each type (see HashCode).
func (t *JSTypeBase) hasAnyTemplateTypesIn(st *JSTypeStack) bool {
	// If the result has been cached use that.
	if r := t.templateCheckResult(); r != 0 {
		return r == 2
	}

	if st.contains(t) {
		// prevent infinite recursion, this is "not yet".
		return false
	}

	st.push(t)
	result := t.self().HasAnyTemplateTypesInternal(st)
	st.pop()

	// Cache the result if the type has been resolved.
	if t.self().IsResolved() {
		if result {
			t.setTemplateCheckResult(2)
		} else {
			t.setTemplateCheckResult(1)
		}
	}
	return result
}

// Java: JSType.hasAnyTemplateTypesInternal()
func (t *JSTypeBase) HasAnyTemplateTypesInternal(st *JSTypeStack) bool {
	return t.self().GetTemplateTypeMap().HasAnyTemplateTypesInternal(st)
}

// GetTemplateTypeMap returns the template type map associated with this
// type.
//
// Java: JSType.getTemplateTypeMap()
func (t *JSTypeBase) GetTemplateTypeMap() *TemplateTypeMap {
	return t.templateTypeMap
}

// GetTypeParameters returns, in order, the sequence of type parameters
// declared for this type.
//
// In general, this value corresponds to an element for every `@template`
// declaration on the type definition. It does not include template
// parameters from superclasses or superinterfaces.
//
// Java: JSType.getTypeParameters()
func (t *JSTypeBase) GetTypeParameters() []*TemplateType {
	m := t.self().GetTemplateTypeMap()
	return jsTypeSubList(m.GetTemplateKeys(), m.Size()-t.self().GetTemplateParamCount(), m.Size())
}

// GetTemplateParamCount returns the number of template parameters declared
// for this type.
//
// In general, this value corresponds to the number of `@template`
// declarations on the type definition. It does not include template
// parameters from superclasses or superinterfaces.
//
// Java: JSType.getTemplateParamCount()
func (t *JSTypeBase) GetTemplateParamCount() int {
	return 0
}

// MergeSupertypeTemplateTypes prepends the template type map associated
// with this type, merging in the keys and values of the specified map.
//
// Java: JSType.mergeSupertypeTemplateTypes(ObjectType)
func (t *JSTypeBase) MergeSupertypeTemplateTypes(other ObjectType) {
	t.self().MaybeLoosenTypecheckingDueToForwardReferencedSupertype(other)

	if other.IsRawTypeOfTemplatizedType() {
		// Before a type can be prepended it needs to be fully specialized. This
		// can happen, for example, when type arguments are not specified in an
		// `@extends` annotation.
		other = t.GetRegistry().CreateTemplatizedTypeOfList(other, []JSType{})
	}
	t.templateTypeMap = other.GetTemplateTypeMap().CopyWithExtension(t.self().GetTemplateTypeMap())
}

// IsObject tests whether this type is an Object, or any subtype thereof.
// Returns this <: Object.
//
// Java: JSType.isObject()
func (t *JSTypeBase) IsObject() bool {
	return false
}

// IsObjectType tests whether this type is an Object, or any subtype thereof.
//
// Java: JSType.isObjectType()
func (t *JSTypeBase) IsObjectType() bool {
	return false
}

// IsConstructor reports whether this type is a FunctionType that is a
// constructor or a named type that points to such a type.
//
// Java: JSType.isConstructor()
func (t *JSTypeBase) IsConstructor() bool {
	return false
}

// IsNominalType reports whether this type is a nominal type (a named
// instance object or a named enum).
//
// Java: JSType.isNominalType()
func (t *JSTypeBase) IsNominalType() bool {
	return false
}

// IsNominalConstructorOrInterface reports whether this type is the original
// constructor of a nominal type. Does not include structural constructors.
//
// Java: JSType.isNominalConstructorOrInterface()
func (t *JSTypeBase) IsNominalConstructorOrInterface() bool {
	if t.self().IsConstructor() || t.self().IsInterface() {
		fn := t.self().ToMaybeFunctionType()
		if fn == nil {
			return false
		}

		// Programmer-defined constructors will have a link back to the original
		// function in the source tree. Structural constructors will not.
		if !fn.GetSource().IsNil() {
			return true
		}

		// Native constructors are always nominal.
		return fn.IsNativeObjectType()
	}
	return false
}

// Java: JSType.isNativeObjectType()
func (t *JSTypeBase) IsNativeObjectType() bool {
	return false
}

// LoosenTypecheckingDueToForwardReferencedSupertype reports whether, during
// type definition, one of the supertypes of this type was a forward
// reference.
//
// This is a hack to work around the fact that inheritance chains and
// template types aren't wired up correctly when this happens, which causes
// various false positives in checks. The known bugs associated are
// b/145145406, b/144327372, and b/132980305.
//
// This method should only be used to suppress potential false positives
// caused by one of the above bugs, in the case where we think suppressing
// typechecking is preferable to emitting a false positive.
//
// Java: JSType.loosenTypecheckingDueToForwardReferencedSupertype()
func (t *JSTypeBase) LoosenTypecheckingDueToForwardReferencedSupertype() bool {
	return t.hasFlag(jsTypeFlagLoosenTypechecking)
}

// Java: JSType.maybeLoosenTypecheckingDueToForwardReferencedSupertype(JSType)
func (t *JSTypeBase) MaybeLoosenTypecheckingDueToForwardReferencedSupertype(supertype JSType) {
	// In the common case this method is never called with a resolved
	// NamedType, but in some edge cases it could be, and typechecking should
	// not be loosened.
	if (supertype.IsNamedType() && !supertype.IsResolved()) ||
		supertype.LoosenTypecheckingDueToForwardReferencedSupertype() {
		t.setFlag(jsTypeFlagLoosenTypechecking)
	}
}

// IsInstanceType reports whether this type is an Instance object of some
// constructor. Does not necessarily mean this is an InstanceObjectType.
//
// Java: JSType.isInstanceType()
func (t *JSTypeBase) IsInstanceType() bool {
	return false
}

// IsInterface reports whether this type is a FunctionType that is an
// interface or a named type that points to such a type.
//
// Java: JSType.isInterface()
func (t *JSTypeBase) IsInterface() bool {
	return false
}

// IsOrdinaryFunction reports whether this type is a FunctionType that is an
// ordinary function (i.e. not a constructor, nominal interface, or record
// interface), or a named type that points to such a type.
//
// Java: JSType.isOrdinaryFunction()
func (t *JSTypeBase) IsOrdinaryFunction() bool {
	return false
}

// Java: JSType.equals(Object)
func (t *JSTypeBase) Equals(other any) bool {
	jSType, ok := other.(JSType)
	if !ok {
		return false
	}
	otherBase := jsTypeBase(jSType)
	if otherBase == nil {
		return false
	}
	if t == otherBase {
		return true
	}
	return equalityChecker_AreIdentityEqualWithBases(t.self_, t, jSType, otherBase)
}

// DiffersFrom reports whether this type is meaningfully different from that
// type for the purposes of data flow analysis.
//
// This is a trickier check than pure equality, because it has to properly
// handle unknown types. See EqMethod for more info.
//
// Java: JSType.differsFrom(JSType)
func (t *JSTypeBase) DiffersFrom(that JSType) bool {
	if jsTypeBase(that) == t {
		return false
	}
	c := EqualityChecker{
		eqMethod: optEqMethodDataFlow,
		hasRun:   true,
	}
	return !c.AreEqualInternal(t.self(), that)
}

// JSType_AreSimilar reports whether the first type is meaningfully same as
// second type for the purposes of data flow analysis.
//
// Java: JSType.areSimilar(JSType, JSType)
func JSType_AreSimilar(first JSType, second JSType) bool {
	fb := jsTypeBase(first)
	sb := jsTypeBase(second)
	if fb == sb {
		return true
	}
	if fb == nil || sb == nil {
		return false
	}
	return !fb.DiffersFrom(second)
}

// HashCode calculates a hash of the object as per Object#hashCode().
//
// Java marks the type as "hash code in progress" to cut recursion (and is
// therefore unsafe for concurrent use). Go-only: the types in progress are
// kept on a stack owned by this call instead, which gives the same values
// and lets several goroutines hash the same types at once.
//
// Java: JSType.hashCode()
func (t *JSTypeBase) HashCode() int32 {
	var hs JSTypeStack
	return t.hashCodeIn(&hs)
}

// JSTypeStack holds the types in progress in one recursive computation:
// one HashCode or HasAnyTemplateTypes call (Java's hashCodeInProgress and
// inTemplatedCheckLoop flags on each type).
type JSTypeStack struct {
	n      int
	inline [8]*JSTypeBase
	more   []*JSTypeBase
}

func (s *JSTypeStack) contains(t *JSTypeBase) bool {
	for i := 0; i < s.n && i < len(s.inline); i++ {
		if s.inline[i] == t {
			return true
		}
	}
	for _, x := range s.more {
		if x == t {
			return true
		}
	}
	return false
}

func (s *JSTypeStack) push(t *JSTypeBase) {
	if s.n < len(s.inline) {
		s.inline[s.n] = t
	} else {
		s.more = append(s.more, t)
	}
	s.n++
}

func (s *JSTypeStack) pop() {
	s.n--
	if s.n >= len(s.inline) {
		s.more = s.more[:len(s.more)-1]
	}
}

func (t *JSTypeBase) hashCodeIn(hs *JSTypeStack) int32 {
	if hs.contains(t) {
		return -1 // Recursive base-case.
	}
	hs.push(t)
	hashCode := t.self().RecursionUnsafeHashCode(hs)
	hs.pop()
	return hashCode
}

// jsTypeHashCodeIn is x.HashCode() inside a hash code computation: 0 for
// nil, like Objects.hashCode.
func jsTypeHashCodeIn(x JSType, hs *JSTypeStack) int32 {
	if x == nil {
		return 0
	}
	return x.GetJSTypeBase().hashCodeIn(hs)
}

// MatchesNumberContext is used to test whether a given type can appear in a
// numeric context, such as an operand of a multiply operator.
//
// Java: JSType.matchesNumberContext()
func (t *JSTypeBase) MatchesNumberContext() bool {
	return false
}

// MatchesStringContext is used to test whether a given type can appear in a
// String context, such as an operand of a string concat (+) operator.
//
// All types have at least the potential for converting to String. When we
// add externally defined types, such as a browser OM, we may choose to add
// types that do not automatically convert to String.
//
// Java: JSType.matchesStringContext()
func (t *JSTypeBase) MatchesStringContext() bool {
	return false
}

// MatchesSymbolContext is used to test whether a given type can appear in a
// symbol context such as property access.
//
// Java: JSType.matchesSymbolContext()
func (t *JSTypeBase) MatchesSymbolContext() bool {
	return false
}

// MatchesObjectContext is used to test whether a given type can appear in
// an Object context, such as the expression in a with statement.
//
// Most types we will encounter, except notably null, have at least the
// potential for converting to Object. Host defined objects can get
// peculiar.
//
// Java: JSType.matchesObjectContext()
func (t *JSTypeBase) MatchesObjectContext() bool {
	return false
}

// FindPropertyTypeByKey coerces this type to an Object type, then gets the
// type of the property whose name is given.
//
// Unlike ObjectType#getPropertyType, returns null if the property is not
// found.
//
// Returns the property's type. null if the current type cannot have
// properties, or if the type is not found.
//
// Java: JSType.findPropertyType(Property.Key)
func (t *JSTypeBase) FindPropertyTypeByKey(propertyName Property_Key) JSType {
	propertyType := t.self().FindPropertyTypeWithoutConsideringTemplateTypes(propertyName)
	if propertyType == nil {
		return nil
	}

	// Do templatized type replacing logic here, and make this method final, to
	// prevent a subclass from forgetting to replace template types
	if t.self().GetTemplateTypeMap().IsEmpty() || !propertyType.HasAnyTemplateTypes() {
		return propertyType
	}

	typeMap := t.self().GetTemplateTypeMap()
	replacer := TemplateTypeReplacer_MakeForPartialReplacement(t.GetRegistry(), typeMap)
	return replacer.VisitJSType(propertyType)
}

// FindPropertyType coerces this type to an Object type, then gets the type
// of the property whose name is given.
//
// Java: JSType.findPropertyType(String)
func (t *JSTypeBase) FindPropertyType(propertyName string) JSType {
	return t.self().FindPropertyTypeByKey(propertyKeyLookup(propertyName))
}

// FindPropertyTypeWithoutConsideringTemplateTypes looks up a property on
// this type, but without properly replacing any templates in the result.
//
// Subclasses can override this if they need more complicated logic for
// property lookup than just autoboxing to an object.
//
// This is only for use by findPropertyType(JSType). Call that method
// instead if you need to lookup a property on a random JSType
//
// Java: JSType.findPropertyTypeWithoutConsideringTemplateTypes(Property.Key)
func (t *JSTypeBase) FindPropertyTypeWithoutConsideringTemplateTypes(propertyName Property_Key) JSType {
	autoboxObjType := ObjectType_Cast(t.self().AutoboxesTo())
	if autoboxObjType != nil {
		return autoboxObjType.FindPropertyTypeByKey(propertyName)
	}
	return nil
}

// CanBeCalled is used to test whether a given type can be used as the
// 'function' in a function call. Returns true if this type might be
// callable.
//
// Java: JSType.canBeCalled()
func (t *JSTypeBase) CanBeCalled() bool {
	return false
}

// JSType_WithSourceRef is a type with source information.
type JSType_WithSourceRef interface {
	// GetSource returns the node representing the declaration of this type.
	GetSource() Node

	// GetGoogModuleId returns the ID of the goog.module in which this type
	// was declared.
	GetGoogModuleId() *string
}

// String is a string representation of this type, suitable for printing in
// warnings.
//
// Java: JSType.toString()
func (t *JSTypeBase) String() string {
	var sb TypeStringBuilder
	sb.AppendType(t.self())
	return sb.Build()
}

// ToAnnotationString: don't call from this package; use appendAsNonNull
// instead.
//
// Java: JSType.toAnnotationString(Nullability)
func (t *JSTypeBase) ToAnnotationString(nullability JSType_Nullability) string {
	builder := TypeStringBuilder{isForAnnotations: true}
	if nullability == JSType_Nullability_EXPLICIT {
		builder.AppendNonNull(t.self())
	} else {
		builder.AppendType(t.self())
	}
	return builder.Build()
}

// Java: JSType.toMaybeObjectType()
func (t *JSTypeBase) ToMaybeObjectType() ObjectType {
	return t.self().ToObjectType()
}

// JSType_Nullability specifies how to express nullability of reference
// types in annotation strings and error messages. Note that this only
// applies to the outer-most type. Nullability of generic type arguments is
// always explicit.
type JSType_Nullability int32

const (
	// JSType_Nullability_EXPLICIT includes an explicit '!' for non-nullable
	// reference types. This is suitable for use in most type contexts
	// (particularly 'type', 'param', and 'return' annotations).
	JSType_Nullability_EXPLICIT JSType_Nullability = iota
	// JSType_Nullability_IMPLICIT omits the explicit '!' from the outermost
	// non-nullable reference type. This is suitable for use in cases where a
	// single reference type is expected (e.g. 'extends' and 'implements').
	JSType_Nullability_IMPLICIT
)

var jsType_Nullability_names = [...]string{"EXPLICIT", "IMPLICIT"}

// String is Enum.name().
func (n JSType_Nullability) String() string { return jsType_Nullability_names[n] }

// jsTypeSubList is ImmutableList.subList(from, to): it throws
// IndexOutOfBoundsException for an invalid range.
func jsTypeSubList[T any](list []T, from int, to int) []T {
	if from < 0 || to > len(list) || from > to {
		panic(jutil.NewIndexOutOfBoundsException(
			"fromIndex: " + jutil.ToJavaString(from) + ", toIndex: " + jutil.ToJavaString(to) +
				", size: " + jutil.ToJavaString(len(list))))
	}
	return list[from:to:to]
}

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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSType represents JavaScript value types.
//
// Types are split into two separate families: value types and object types.
//
// A special UnknownType exists to represent a wildcard type on which no
// information can be gathered. In particular, it can assign to everyone, is
// a subtype of everyone (and everyone is a subtype of it).
//
// If you remove the UnknownType, the set of types in the type system forms
// a lattice with the #isSubtype relation defining the partial order of
// types. All types are united at the top of the lattice by the AllType and
// at the bottom by the NoType.
//
// Java: abstract class JSType. The Go interface declares every non-private
// instance method; JSTypeBase holds the fields and the base implementations
// and dispatches virtual calls through self (PORTING.md §4).
//
// Overloads (PORTING.md §3): the String-keyed property overloads keep the
// plain name and the Property.Key ones take a `ByKey` suffix, e.g.
// getPropertyType(String) → GetPropertyType, getPropertyType(Property.Key)
// → GetPropertyTypeByKey; an extra boolean/mode parameter is a `With…`
// suffix (GetPropertyKindWithAutobox, IsSubtypeOfWithMode). Static methods
// are package functions JSType_Foo.
type JSType interface {
	// Java: JSType.getTypeClass()
	GetTypeClass() JSTypeClass
	// Java: JSType.getNativeType(JSTypeNative)
	GetNativeType(typeId JSTypeNative) JSType
	// Java: JSType.getJSDocInfo()
	GetJSDocInfo() *JSDocInfo
	// Java: JSType.getDisplayName()
	GetDisplayName() *string
	// Java: JSType.hasDisplayName()
	HasDisplayName() bool
	// Java: JSType.getPropertyKind(Property.Key)
	GetPropertyKindByKey(pname Property_Key) JSType_HasPropertyKind
	// Java: JSType.getPropertyKind(Property.Key, boolean)
	GetPropertyKindByKeyWithAutobox(pname Property_Key, autobox bool) JSType_HasPropertyKind
	// Java: JSType.getPropertyKind(String, boolean)
	GetPropertyKindWithAutobox(pname string, autobox bool) JSType_HasPropertyKind
	// Java: JSType.getPropertyKind(String)
	GetPropertyKind(pname string) JSType_HasPropertyKind
	// Java: JSType.hasProperty(String)
	HasProperty(pname string) bool
	// Java: JSType.hasProperty(Property.Key)
	HasPropertyByKey(pname Property_Key) bool
	IsNoType() bool
	IsNoResolvedType() bool
	IsNoObjectType() bool
	IsEmptyType() bool
	IsNumberObjectType() bool
	IsNumberValueType() bool
	IsBigIntObjectType() bool
	IsBigIntValueType() bool
	IsFunctionPrototypeType() bool
	IsStringObjectType() bool
	IsSymbolObjectType() bool
	IsTheObjectType() bool
	IsStringValueType() bool
	IsSymbolValueType() bool
	IsKnownSymbolValueType() bool
	IsString() bool
	IsNumber() bool
	IsSymbol() bool
	IsBigIntOrNumber() bool
	IsOnlyBigInt() bool
	IsArrayType() bool
	IsReadonlyArrayType() bool
	IsBooleanObjectType() bool
	IsBooleanValueType() bool
	IsRegexpType() bool
	IsDateType() bool
	IsNullType() bool
	IsVoidType() bool
	IsAllType() bool
	IsUnknownType() bool
	IsSomeUnknownType() bool
	IsCheckedUnknownType() bool
	IsUnionType() bool
	IsRawTypeOfTemplatizedType() bool
	IsStruct() bool
	IsDict() bool
	ContainsReferenceAncestor(target JSType) bool
	IsLiteralObject() bool
	ToMaybeUnionType() *UnionType
	IsGlobalThisType() bool
	IsFunctionType() bool
	ToMaybeFunctionType() FunctionType
	AssertFunctionType() FunctionType
	AssertObjectType() ObjectType
	ToMaybeKnownSymbolType() *KnownSymbolType
	IsEnumElementType() bool
	GetEnumeratedTypeOfEnumElement() JSType
	ToMaybeEnumElementType() *EnumElementType
	IsEnumType() bool
	ToMaybeEnumType() *EnumType
	IsNamedType() bool
	ToMaybeNamedType() *NamedType
	IsRecordType() bool
	IsStructuralInterface() bool
	IsStructuralType() bool
	ToMaybeRecordType() *RecordType
	IsTemplatizedType() bool
	ToMaybeTemplatizedType() *TemplatizedType
	IsTemplateType() bool
	ToMaybeTemplateType() *TemplateType
	HasAnyTemplateTypes() bool
	HasAnyTemplateTypesInternal() bool
	GetTemplateTypeMap() *TemplateTypeMap
	GetTypeParameters() []*TemplateType
	GetTemplateParamCount() int
	MergeSupertypeTemplateTypes(other ObjectType)
	IsObject() bool
	IsObjectType() bool
	IsConstructor() bool
	IsNominalType() bool
	IsNominalConstructorOrInterface() bool
	IsNativeObjectType() bool
	LoosenTypecheckingDueToForwardReferencedSupertype() bool
	MaybeLoosenTypecheckingDueToForwardReferencedSupertype(supertype JSType)
	IsInstanceType() bool
	IsInterface() bool
	IsOrdinaryFunction() bool
	// Java: JSType.equals(Object)
	Equals(other any) bool
	DiffersFrom(that JSType) bool
	// Java: JSType.hashCode()
	HashCode() int32
	RecursionUnsafeHashCode() int32
	MatchesNumberContext() bool
	MatchesStringContext() bool
	MatchesSymbolContext() bool
	MatchesObjectContext() bool
	// Java: JSType.findPropertyType(Property.Key)
	FindPropertyTypeByKey(propertyName Property_Key) JSType
	// Java: JSType.findPropertyType(String)
	FindPropertyType(propertyName string) JSType
	FindPropertyTypeWithoutConsideringTemplateTypes(propertyName Property_Key) JSType
	CanBeCalled() bool
	CanCastTo(that JSType) bool
	AutoboxesTo() JSType
	IsBoxableScalar() bool
	ToObjectType() ObjectType
	Autobox() JSType
	Dereference() ObjectType
	CanTestForEqualityWith(that JSType) bool
	TestForEquality(that JSType) base.Tri
	TestForEqualityHelper(aType JSType, bType JSType) base.Tri
	CanTestForShallowEqualityWith(that JSType) bool
	IsNullable() bool
	IsVoidable() bool
	IsExplicitlyVoidable() bool
	CollapseUnion() JSType
	GetLeastSupertype(that JSType) JSType
	GetGreatestSubtype(that JSType) JSType
	GetRestrictedTypeGivenOutcome(outcome Outcome) JSType
	GetPossibleToBooleanOutcomes() BooleanLiteralSet
	GetTypesUnderEquality(that JSType) *JSType_TypePair
	GetTypesUnderEqualityPair(that JSType) (JSType, JSType)
	GetTypesUnderInequality(that JSType) *JSType_TypePair
	GetTypesUnderInequalityPair(that JSType) (JSType, JSType)
	GetTypesUnderShallowEquality(that JSType) *JSType_TypePair
	GetTypesUnderShallowEqualityPair(that JSType) (JSType, JSType)
	GetTypesUnderShallowInequality(that JSType) *JSType_TypePair
	GetTypesUnderShallowInequalityPair(that JSType) (JSType, JSType)
	GetUnionMembers() []JSType
	RestrictByNotNullOrUndefined() JSType
	RestrictByNotUndefined() JSType
	RestrictByNotNull() JSType
	IsSubtypeWithoutStructuralTyping(supertype JSType) bool
	IsSubtype(supertype JSType) bool
	IsSubtypeWithMode(supertype JSType, mode JSType_SubtypingMode) bool
	IsSubtypeOf(supertype JSType) bool
	IsSubtypeOfWithMode(supertype JSType, mode JSType_SubtypingMode) bool
	Resolve(reporter ErrorReporter) JSType
	ResolveInternal(reporter ErrorReporter) JSType
	EagerlyResolveToSelf()
	IsResolved() bool
	IsSuccessfullyResolved() bool
	IsUnsuccessfullyResolved() bool
	SetValidator(validator func(JSType) bool) bool
	// Java: JSType.toString()
	String() string
	ToAnnotationString(nullability JSType_Nullability) string
	AppendTo(sb *TypeStringBuilder)
	MatchConstraint(constraint JSType)
	ToMaybeObjectType() ObjectType

	// GetJSTypeBase gives access to JSType's fields through the interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4). It also keeps non-JSType values (such as
	// *colors.Color, which Node stores in the same field) from satisfying
	// the interface.
	GetJSTypeBase() *JSTypeBase
}

// JSTypeBase holds the fields and base method implementations of JSType.
type JSTypeBase struct {
	// self_ is the concrete JSType that embeds this base, the receiver of
	// virtual calls (PORTING.md §4).
	self_ JSType
	// The result of Resolve; nil until the type is resolved.
	resolveResult *JSTypeBase
	// The type's template type map, never nil; its registry is the type's
	// registry (Java's JSType.registry field; see GetRegistry).
	templateTypeMap *TemplateTypeMap

	// nodeHandle is the type's handle in the node type table once a Node
	// referred to it (see node_tables.go); 0 before.
	nodeHandle uint32

	// flags holds the rarely used booleans (jsTypeFlag*) and the 2-bit
	// templateCheckResult (0 = unset, 1 = false, 2 = true).
	flags    uint8
	unknown  bool
	objFlags uint8
}

const (
	jsTypeFlagLoosenTypechecking uint8 = 1 << 0 // loosenTypecheckingDueToForwardReferencedSupertype
	jsTypeFlagHashCodeInProgress uint8 = 1 << 1
	jsTypeFlagInTemplatedCheck   uint8 = 1 << 2
	jsTypeFlagVisited            uint8 = 1 << 3
	jsTypeFlagResolvedToSelf     uint8 = 1 << 4
	jsTypeTemplateCheckShift           = 5
	jsTypeTemplateCheckMask      uint8 = 3 << jsTypeTemplateCheckShift
)

func (t *JSTypeBase) templateCheckResult() uint8 {
	return (t.flags & jsTypeTemplateCheckMask) >> jsTypeTemplateCheckShift
}

func (t *JSTypeBase) setTemplateCheckResult(v uint8) {
	t.flags = (t.flags &^ jsTypeTemplateCheckMask) | (v << jsTypeTemplateCheckShift)
}

const (
	objKindFallback uint8 = iota
	objKindFunction
	objKindPrototype
	objKindInstance
	objKindNamed
	objKindTemplatized
	objKindRecord
	objKindEnum
	objKindUnknown
	objKindNoObject
	objKindNo
	objKindNoResolved
	objKindEnumElement
	objKindTemplate
	objKindProxy
	jsKindUnion
	jsKindAll
	jsKindNull
	jsKindVoid
	jsKindScalarValue
	jsKindKnownSymbol
)

const (
	objKindMask   uint8 = 0x7f
	objHasDocFlag uint8 = 0x80
)

const (
	jsKindProxyOrFallbackMask uint32 = (1 << objKindFallback) |
		(1 << objKindNamed) |
		(1 << objKindTemplatized) |
		(1 << objKindTemplate) |
		(1 << objKindProxy)
	jsKindMaybeUnknownMask uint32 = jsKindProxyOrFallbackMask | (1 << jsKindUnion)
	jsKindComplexEqMask    uint32 = jsKindMaybeUnknownMask | (1 << objKindUnknown)
	jsKindValueMask        uint32 = (1 << jsKindNull) |
		(1 << jsKindVoid) |
		(1 << jsKindScalarValue) |
		(1 << jsKindKnownSymbol)
	jsKindNonObjectMask           uint32 = jsKindValueMask | (1 << jsKindAll)
	jsKindMaybeRecordMask         uint32 = (1 << objKindRecord) | jsKindProxyOrFallbackMask
	jsKindSimpleDisjointValueMask uint32 = (1 << jsKindNull) | (1 << jsKindVoid) | (1 << jsKindScalarValue)
	jsKindCanRelateToValueMask    uint32 = jsKindProxyOrFallbackMask |
		(1 << jsKindKnownSymbol) |
		(1 << objKindEnumElement) |
		(1 << objKindNoObject) |
		(1 << objKindNo) |
		(1 << jsKindAll) |
		(1 << jsKindUnion)
)

func jsTypeBase(x JSType) *JSTypeBase {
	if x == nil {
		return nil
	}
	return x.GetJSTypeBase()
}

func jsTypeEquals(a, b JSType) bool {
	ab := jsTypeBase(a)
	bb := jsTypeBase(b)
	if ab == nil || bb == nil {
		return false
	}
	if ab == bb {
		return true
	}
	return equalityChecker_AreIdentityEqualWithBases(a, ab, b, bb)
}

func jsTypeIsUnknown(t JSType, tb *JSTypeBase, kind uint8) bool {
	if ((uint32(1)<<kind)&jsKindMaybeUnknownMask) == 0 && !tb.unknown {
		return false
	}
	if kind == objKindUnknown {
		return true
	}
	return t.IsUnknownType()
}

func jsTypeHasAnyTemplateTypes(tb *JSTypeBase) bool {
	if r := tb.templateCheckResult(); r != 0 {
		return r == 2
	}
	return tb.HasAnyTemplateTypes()
}

func jsTypeIsTemplatized(t JSType, kind uint8) bool {
	if kind == objKindTemplatized {
		return true
	}
	if ((uint32(1) << kind) & jsKindProxyOrFallbackMask) != 0 {
		return t.IsTemplatizedType()
	}
	return false
}

func jsTypeIsFunction(t JSType, kind uint8) bool {
	if kind == objKindFunction {
		return true
	}
	if ((uint32(1) << kind) & jsKindProxyOrFallbackMask) != 0 {
		return t.IsFunctionType()
	}
	return false
}

// Java: JSType.CAN_CAST_TO_VISITOR
var JSType_CAN_CAST_TO_VISITOR = NewCanCastToVisitor()

// InitJSTypeBase is the Go form of the JSType(JSTypeRegistry) constructor;
// self is the concrete type.
//
// Java: JSType(JSTypeRegistry)
func (t *JSTypeBase) InitJSTypeBase(self JSType, registry *JSTypeRegistry) {
	t.InitJSTypeBaseWithTemplateTypeMap(self, registry, nil)
}

// InitJSTypeBaseWithTemplateTypeMap is the Go form of the
// JSType(JSTypeRegistry, TemplateTypeMap) constructor.
//
// Java: JSType(JSTypeRegistry, TemplateTypeMap)
func (t *JSTypeBase) InitJSTypeBaseWithTemplateTypeMap(
	self JSType, registry *JSTypeRegistry, templateTypeMap *TemplateTypeMap) {
	t.self_ = self

	if templateTypeMap == nil {
		t.templateTypeMap = registry.GetEmptyTemplateTypeMap()
	} else {
		jutil.CheckState(templateTypeMap.registry == registry)
		t.templateTypeMap = templateTypeMap
	}

	registry.GetResolver().AddUnresolved(self)
}

// GetJSTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (t *JSTypeBase) GetJSTypeBase() *JSTypeBase { return t }

// jsTypeOf returns the JSType whose base is b, or nil for nil.
func jsTypeOf(b *JSTypeBase) JSType {
	if b == nil {
		return nil
	}
	return b.self_
}

// self returns the concrete JSType embedding this base.
func (t *JSTypeBase) self() JSType {
	return t.self_
}

// Self returns the concrete JSType embedding this base.
func (t *JSTypeBase) Self() JSType { return t.self() }

// GetRegistry returns the type's registry (Java's JSType.registry field).
//
// Go-only: performance — read from the template type map, which belongs to
// the same registry, instead of a field of its own in every type.
func (t *JSTypeBase) GetRegistry() *JSTypeRegistry {
	if t.templateTypeMap == nil {
		return nil
	}
	return t.templateTypeMap.registry
}

// GetNativeType is a utility method for less verbose code.
//
// Java: JSType.getNativeType(JSTypeNative)
func (t *JSTypeBase) GetNativeType(typeId JSTypeNative) JSType {
	return t.GetRegistry().GetNativeType(typeId)
}

// GetJSDocInfo gets the docInfo for this type. By default, documentation
// cannot be attached to arbitrary types. This must be overridden for
// programmer-defined types.
//
// Java: JSType.getJSDocInfo()
func (t *JSTypeBase) GetJSDocInfo() *JSDocInfo {
	return nil
}

// GetDisplayName returns a user meaningful label for the JSType instance.
// For example, Functions and Enums will return their declaration name (if
// they have one). Some types will not have a meaningful display name. Calls
// to hasDisplayName() will return true IFF getDisplayName() will return null
// or a zero length string.
//
// Returns the display name of the type, or null if one is not available.
//
// Java: JSType.getDisplayName()
func (t *JSTypeBase) GetDisplayName() *string {
	return nil
}

// HasDisplayName returns true if the JSType has a user meaningful label.
//
// Java: JSType.hasDisplayName()
func (t *JSTypeBase) HasDisplayName() bool {
	displayName := t.self().GetDisplayName()
	return displayName != nil && *displayName != ""
}

// JSType_HasPropertyKind is a tristate value returned from
// canPropertyBeDefined.
type JSType_HasPropertyKind int32

const (
	JSType_HasPropertyKind_ABSENT        JSType_HasPropertyKind = iota // The property is not known to be part of this type
	JSType_HasPropertyKind_KNOWN_PRESENT                               // The properties is known to be defined on a type or its super types
	JSType_HasPropertyKind_MAYBE_PRESENT                               // The property is loosely associated with a type, typically one of its subtypes
)

var jsType_HasPropertyKind_names = [...]string{"ABSENT", "KNOWN_PRESENT", "MAYBE_PRESENT"}

// String is Enum.name().
func (k JSType_HasPropertyKind) String() string { return jsType_HasPropertyKind_names[k] }

// Java: JSType.HasPropertyKind.of(boolean)
func JSType_HasPropertyKind_Of(has bool) JSType_HasPropertyKind {
	if has {
		return JSType_HasPropertyKind_KNOWN_PRESENT
	}
	return JSType_HasPropertyKind_ABSENT
}

// GetPropertyKindByKey checks whether the property is present on the
// object.
//
// Java: JSType.getPropertyKind(Property.Key)
func (t *JSTypeBase) GetPropertyKindByKey(pname Property_Key) JSType_HasPropertyKind {
	return t.self().GetPropertyKindByKeyWithAutobox(pname, true)
}

// GetPropertyKindByKeyWithAutobox checks whether the property is present on
// the object. autobox: whether to check for the presents on an autoboxed
// type.
//
// Java: JSType.getPropertyKind(Property.Key, boolean)
func (t *JSTypeBase) GetPropertyKindByKeyWithAutobox(pname Property_Key, autobox bool) JSType_HasPropertyKind {
	return JSType_HasPropertyKind_ABSENT
}

// GetPropertyKindWithAutobox checks whether the property is present on the
// object.
//
// Java: JSType.getPropertyKind(String, boolean)
func (t *JSTypeBase) GetPropertyKindWithAutobox(pname string, autobox bool) JSType_HasPropertyKind {
	return t.self().GetPropertyKindByKeyWithAutobox(NewProperty_StringKey(pname), autobox)
}

// GetPropertyKind checks whether the property is present on the object.
//
// Java: JSType.getPropertyKind(String)
func (t *JSTypeBase) GetPropertyKind(pname string) JSType_HasPropertyKind {
	return t.self().GetPropertyKindByKey(NewProperty_StringKey(pname))
}

// HasProperty checks whether the property is present on the object.
//
// Java: JSType.hasProperty(String)
func (t *JSTypeBase) HasProperty(pname string) bool {
	return t.self().HasPropertyByKey(NewProperty_StringKey(pname))
}

// Java: JSType.hasProperty(Property.Key)
func (t *JSTypeBase) HasPropertyByKey(pname Property_Key) bool {
	return t.self().GetPropertyKindByKeyWithAutobox(pname, false) != JSType_HasPropertyKind_ABSENT
}

// Java: JSType.isNoType()
func (t *JSTypeBase) IsNoType() bool {
	return false
}

// Java: JSType.isNoResolvedType()
func (t *JSTypeBase) IsNoResolvedType() bool {
	return false
}

// Java: JSType.isNoObjectType()
func (t *JSTypeBase) IsNoObjectType() bool {
	return false
}

// Java: JSType.isEmptyType()
func (t *JSTypeBase) IsEmptyType() bool {
	return t.self().IsNoType() ||
		t.self().IsNoObjectType() ||
		t.self().IsNoResolvedType() ||
		JSType(t.self()) == JSType(t.GetRegistry().GetNativeFunctionType(JSTypeNative_LEAST_FUNCTION_TYPE))
}

// Java: JSType.isNumberObjectType()
func (t *JSTypeBase) IsNumberObjectType() bool {
	return false
}

// Java: JSType.isNumberValueType()
func (t *JSTypeBase) IsNumberValueType() bool {
	return false
}

// Java: JSType.isBigIntObjectType()
func (t *JSTypeBase) IsBigIntObjectType() bool {
	return false
}

// Java: JSType.isBigIntValueType()
func (t *JSTypeBase) IsBigIntValueType() bool {
	return false
}

// IsFunctionPrototypeType reports whether this is the prototype of a
// function.
//
// Java: JSType.isFunctionPrototypeType()
func (t *JSTypeBase) IsFunctionPrototypeType() bool {
	// TODO(sdh): consider renaming this to isPrototypeObject.
	return false
}

// Java: JSType.isStringObjectType()
func (t *JSTypeBase) IsStringObjectType() bool {
	return false
}

// Java: JSType.isSymbolObjectType()
func (t *JSTypeBase) IsSymbolObjectType() bool {
	return false
}

// Java: JSType.isTheObjectType()
func (t *JSTypeBase) IsTheObjectType() bool {
	return false
}

// Java: JSType.isStringValueType()
func (t *JSTypeBase) IsStringValueType() bool {
	return false
}

// Java: JSType.isSymbolValueType()
func (t *JSTypeBase) IsSymbolValueType() bool {
	return false
}

// Java: JSType.isKnownSymbolValueType()
func (t *JSTypeBase) IsKnownSymbolValueType() bool {
	return false
}

// IsString tests whether the type is a string (value or Object). Returns
// this <: (String, string).
//
// Java: JSType.isString()
func (t *JSTypeBase) IsString() bool {
	return t.self().IsSubtypeOf(t.GetNativeType(JSTypeNative_STRING_TYPE)) ||
		t.self().IsSubtypeOf(t.GetNativeType(JSTypeNative_STRING_OBJECT_TYPE))
}

// IsNumber tests whether the type is a number (value or Object). Returns
// this <: (Number, number).
//
// Java: JSType.isNumber()
func (t *JSTypeBase) IsNumber() bool {
	return t.self().IsSubtypeOf(t.GetNativeType(JSTypeNative_NUMBER_TYPE)) ||
		t.self().IsSubtypeOf(t.GetNativeType(JSTypeNative_NUMBER_OBJECT_TYPE))
}

// Java: JSType.isSymbol()
func (t *JSTypeBase) IsSymbol() bool {
	return t.self().IsSubtypeOf(t.GetNativeType(JSTypeNative_SYMBOL_TYPE)) ||
		t.self().IsSubtypeOf(t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE))
}

// Java: JSType.isBigIntOrNumber()
func (t *JSTypeBase) IsBigIntOrNumber() bool {
	// NOTE: This method will return false for unions that mix the primitive
	// and object types, e.g. (bigint|BigInt). This is technically wrong, but
	// is being done on purpose because 1) it maintains consistency with
	// previously existing methods such as isNumber() above, 2) The object
	// versions of numeric types are very rare in practice and shouldn't be
	// used in general, and 3) a value that could be either the primitive or
	// the object type seems likely to be the result of a coding error.
	return t.self().IsSubtypeOf(t.GetNativeType(JSTypeNative_BIGINT_NUMBER)) ||
		t.self().IsSubtypeOf(t.GetNativeType(JSTypeNative_BIGINT_NUMBER_OBJECT))
}

// IsOnlyBigInt checks whether the type is a bigint and *only* a bigint
// (value or Object).
//
// Java: JSType.isOnlyBigInt()
func (t *JSTypeBase) IsOnlyBigInt() bool {
	return t.self().IsBigIntValueType() || t.self().IsBigIntObjectType()
}

// Java: JSType.isArrayType()
func (t *JSTypeBase) IsArrayType() bool {
	return false
}

// Java: JSType.isReadonlyArrayType()
func (t *JSTypeBase) IsReadonlyArrayType() bool {
	return false
}

// Java: JSType.isBooleanObjectType()
func (t *JSTypeBase) IsBooleanObjectType() bool {
	return false
}

// Java: JSType.isBooleanValueType()
func (t *JSTypeBase) IsBooleanValueType() bool {
	return false
}

// Java: JSType.isRegexpType()
func (t *JSTypeBase) IsRegexpType() bool {
	return false
}

// Java: JSType.isDateType()
func (t *JSTypeBase) IsDateType() bool {
	return false
}

// Java: JSType.isNullType()
func (t *JSTypeBase) IsNullType() bool {
	return false
}

// Java: JSType.isVoidType()
func (t *JSTypeBase) IsVoidType() bool {
	return false
}

// Java: JSType.isAllType()
func (t *JSTypeBase) IsAllType() bool {
	return false
}

// Java: JSType.isUnknownType()
func (t *JSTypeBase) IsUnknownType() bool {
	return false
}

// Java: JSType.isSomeUnknownType()
func (t *JSTypeBase) IsSomeUnknownType() bool {
	// OTI's notion of isUnknownType already accounts for looseness (see
	// override in ObjectType).
	return t.self().IsUnknownType()
}

// Java: JSType.isCheckedUnknownType()
func (t *JSTypeBase) IsCheckedUnknownType() bool {
	return false
}

// Java: JSType.isUnionType()
func (t *JSTypeBase) IsUnionType() bool {
	return false
}

// Java: JSType.isRawTypeOfTemplatizedType()
func (t *JSTypeBase) IsRawTypeOfTemplatizedType() bool {
	return t.self().GetTemplateParamCount() > 0 && !t.self().IsTemplatizedType()
}

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
	visitor := NewContainsUpperBoundSuperTypeVisitor(target)
	return VisitJSType[ContainsUpperBoundSuperTypeVisitor_Result](t.self(), visitor) == ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
}

// Java: JSType.isLiteralObject()
func (t *JSTypeBase) IsLiteralObject() bool {
	if prototypeObjectType, ok := t.self().(PrototypeObjectType); ok {
		return prototypeObjectType.IsAnonymous()
	}
	return false
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
	// If the result has been cached use that.
	if r := t.templateCheckResult(); r != 0 {
		return r == 2
	}

	if t.flags&jsTypeFlagInTemplatedCheck != 0 {
		// prevent infinite recursion, this is "not yet".
		return false
	}

	t.flags |= jsTypeFlagInTemplatedCheck
	result := t.self().HasAnyTemplateTypesInternal()
	t.flags &^= jsTypeFlagInTemplatedCheck

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
func (t *JSTypeBase) HasAnyTemplateTypesInternal() bool {
	return t.self().GetTemplateTypeMap().HasAnyTemplateTypesInternal()
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
	return t.flags&jsTypeFlagLoosenTypechecking != 0
}

// Java: JSType.maybeLoosenTypecheckingDueToForwardReferencedSupertype(JSType)
func (t *JSTypeBase) MaybeLoosenTypecheckingDueToForwardReferencedSupertype(supertype JSType) {
	// In the common case this method is never called with a resolved
	// NamedType, but in some edge cases it could be, and typechecking should
	// not be loosened.
	if (supertype.IsNamedType() && !supertype.IsResolved()) ||
		supertype.LoosenTypecheckingDueToForwardReferencedSupertype() {
		t.flags |= jsTypeFlagLoosenTypechecking
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
// This method is unsafe for multi-threaded use. The implementation mutates
// instance state to prevent recursion and therefore expects sole access.
//
// Java: JSType.hashCode()
func (t *JSTypeBase) HashCode() int32 {
	if t.flags&jsTypeFlagHashCodeInProgress != 0 {
		return -1 // Recursive base-case.
	}

	t.flags |= jsTypeFlagHashCodeInProgress
	hashCode := t.self().RecursionUnsafeHashCode()
	t.flags &^= jsTypeFlagHashCodeInProgress
	return hashCode
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
	return t.self().FindPropertyTypeByKey(NewProperty_StringKey(propertyName))
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
	return NewTypeStringBuilder(false).AppendType(t.self()).Build()
}

// ToAnnotationString: don't call from this package; use appendAsNonNull
// instead.
//
// Java: JSType.toAnnotationString(Nullability)
func (t *JSTypeBase) ToAnnotationString(nullability JSType_Nullability) string {
	builder := NewTypeStringBuilder(true)
	if nullability == JSType_Nullability_EXPLICIT {
		return builder.AppendNonNull(t.self()).Build()
	}
	return builder.AppendType(t.self()).Build()
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

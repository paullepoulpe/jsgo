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
	GetTypesUnderEqualityPair(that JSType) (JSType, JSType)
	GetTypesUnderInequalityPair(that JSType) (JSType, JSType)
	GetTypesUnderShallowEqualityPair(that JSType) (JSType, JSType)
	GetTypesUnderShallowInequalityPair(that JSType) (JSType, JSType)
	GetUnionMembers() []JSType
	RestrictByNotNullOrUndefined() JSType
	RestrictByNotUndefined() JSType
	RestrictByNotNull() JSType
	IsSubtypeWithoutStructuralTyping(supertype JSType) bool
	IsSubtype(supertype JSType) bool
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
	flags      uint8
	unknown    bool
	objFlags   uint8
	protoFlags uint8
}

const (
	protoFlagNative         uint8 = 1 << 0
	protoFlagAnonymous      uint8 = 1 << 1
	protoFlagPrettyPrint    uint8 = 1 << 2
	protoFlagRecordDeclared uint8 = 1 << 3
	protoFlagRecordFrozen   uint8 = 1 << 4
)

const (
	jsTypeFlagLoosenTypechecking      uint8 = 1 << 0 // loosenTypecheckingDueToForwardReferencedSupertype
	jsTypeFlagHashCodeInProgress      uint8 = 1 << 1
	jsTypeFlagInTemplatedCheck        uint8 = 1 << 2
	jsTypeFlagVisited                 uint8 = 1 << 3
	jsTypeFlagResolvedToSelf          uint8 = 1 << 4
	jsTypeTemplateCheckShift                = 5
	jsTypeTemplateCheckMask           uint8 = 3 << jsTypeTemplateCheckShift
	jsTypeFlagWellFormedTemplateValid uint8 = 1 << 7
)

// IsWellFormedTemplateValid reports whether TypeValidator has already verified
// that this resolved type has well-formed bounded generic arguments.
func (t *JSTypeBase) IsWellFormedTemplateValid() bool {
	return t.flags&jsTypeFlagWellFormedTemplateValid != 0
}

// MarkWellFormedTemplateValidIfResolved records a positive well-formed
// templatized-type check result once the type is resolved and the resolver is
// closed (so referenced types and bounds are final).
func (t *JSTypeBase) MarkWellFormedTemplateValidIfResolved() {
	if t.resolveResult != nil {
		if reg := t.GetRegistry(); reg != nil && reg.resolver != nil && reg.resolver.state == JSTypeResolver_State_CLOSED {
			t.flags |= jsTypeFlagWellFormedTemplateValid
		}
	}
}

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
	return t.self().GetPropertyKindByKeyWithAutobox(propertyKeyLookup(pname), autobox)
}

// GetPropertyKind checks whether the property is present on the object.
//
// Java: JSType.getPropertyKind(String)
func (t *JSTypeBase) GetPropertyKind(pname string) JSType_HasPropertyKind {
	return t.self().GetPropertyKindByKey(propertyKeyLookup(pname))
}

// HasProperty checks whether the property is present on the object.
//
// Java: JSType.hasProperty(String)
func (t *JSTypeBase) HasProperty(pname string) bool {
	return t.self().HasPropertyByKey(propertyKeyLookup(pname))
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

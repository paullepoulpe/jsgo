// Port of com/google/javascript/rhino/jstype/JSType.java @ closure-compiler v20261015

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
	"unsafe"

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
// suffix (GetPropertyKindWithAutobox, IsSubtypeOfWithMode);
// visit(RelationshipVisitor, JSType) is VisitRelationship. Static methods
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
	// Java: JSType.visit(Visitor<T>)
	Visit(visitor Visitor) any
	// Java: JSType.visit(RelationshipVisitor<T>, JSType)
	VisitRelationship(visitor RelationshipVisitor, that JSType) any
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
	// selfTab is the itab word of the JSType interface of the concrete type
	// that embeds this base; with the base's own address (the base is the
	// first field of every type) it forms self(), the receiver of virtual
	// calls (PORTING.md §4).
	//
	// Go-only: performance — half the size of a JSType field, and not a
	// pointer for the GC (itabs are never freed).
	selfTab uintptr
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

const (
	_ = -uint(unsafe.Offsetof(FunctionTypeBase{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(PrototypeObjectTypeBase{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(InstanceObjectType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(NamedType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(TemplatizedType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(RecordType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(EnumType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(UnknownType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(NoObjectTypeBase{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(NoType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(NoResolvedType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(EnumElementType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(TemplateType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(ProxyObjectTypeBase{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(UnionType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(AllType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(NullType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(VoidType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(BooleanType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(NumberType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(StringType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(BigIntType{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(SymbolTypeBase{}.JSTypeBase))
	_ = -uint(unsafe.Offsetof(KnownSymbolType{}.JSTypeBase))
)

// Size guards for the most numerous parts of the type graph: each bound is a
// Go size class, so growing one of these structs past it costs a whole class.
const (
	_ = uintptr(32) - unsafe.Sizeof(JSTypeBase{})
	_ = uintptr(96) - unsafe.Sizeof(PrototypeObjectTypeBase{})
	_ = uintptr(176) - unsafe.Sizeof(FunctionTypeBase{})
	_ = uintptr(96) - unsafe.Sizeof(UnionType{})
	_ = uintptr(96) - unsafe.Sizeof(TemplatizedType{})
	_ = uintptr(48) - unsafe.Sizeof(Property{})
	_ = uintptr(16) - unsafe.Sizeof(FunctionType_Parameter{})
)

type jsTypeIface struct {
	tab  unsafe.Pointer
	data *JSTypeBase
}

type anyEface struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}

func jsTypeBase(x JSType) *JSTypeBase {
	return (*(*jsTypeIface)(unsafe.Pointer(&x))).data
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

// Java: JSType.serialVersionUID
const JSType_serialVersionUID int64 = 1

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
	iface := (*jsTypeIface)(unsafe.Pointer(&self))
	if iface.data != t {
		panic("JSTypeBase must be the first field of " + jutil.ToJavaString(self))
	}
	t.selfTab = uintptr(iface.tab)

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
	return b.self()
}

// self returns the concrete JSType embedding this base.
func (t *JSTypeBase) self() JSType {
	var x JSType
	*(*jsTypeIface)(unsafe.Pointer(&x)) = jsTypeIface{*(*unsafe.Pointer)(unsafe.Pointer(&t.selfTab)), t}
	return x
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

// JSType_HasPropertyKind_values is HasPropertyKind.values().
func JSType_HasPropertyKind_values() []JSType_HasPropertyKind {
	return []JSType_HasPropertyKind{
		JSType_HasPropertyKind_ABSENT, JSType_HasPropertyKind_KNOWN_PRESENT, JSType_HasPropertyKind_MAYBE_PRESENT,
	}
}

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
	return t.self().Visit(visitor) == ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
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
	if (*(*anyEface)(unsafe.Pointer(&other))).data == unsafe.Pointer(t) {
		return any(t.self()) == other
	}
	jSType, ok := other.(JSType)
	if !ok {
		return false
	}
	otherBase := jsTypeBase(jSType)
	if otherBase == nil {
		return false
	}
	return equalityChecker_AreIdentityEqualWithBases(t.self(), t, jSType, otherBase)
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
		eqMethod: &equalityChecker_eqMethods[EqualityChecker_EqMethod_DATA_FLOW],
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
	replacer := MakeTemplateTypeReplacerForPartialReplacement(t.GetRegistry(), typeMap)
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

// CanCastTo tests whether values of this type can be safely assigned to
// values of that type.
//
// The default implementation verifies that this is a subtype of that.
//
// Java: JSType.canCastTo(JSType)
func (t *JSTypeBase) CanCastTo(that JSType) bool {
	return t.self().VisitRelationship(JSType_CAN_CAST_TO_VISITOR, that).(bool)
}

// AutoboxesTo turns a scalar type to the corresponding object type. Returns
// the auto-boxed type or null if this type is not a scalar.
//
// Java: JSType.autoboxesTo()
func (t *JSTypeBase) AutoboxesTo() JSType {
	return nil
}

// Java: JSType.isBoxableScalar()
func (t *JSTypeBase) IsBoxableScalar() bool {
	return t.self().AutoboxesTo() != nil
}

// ToObjectType casts this to an ObjectType, or returns null if this is not
// an ObjectType. If this is a scalar type, it will *not* be converted to an
// object type. If you want to simulate JS autoboxing or dereferencing, you
// should use autoboxesTo() or dereference().
//
// Java: JSType.toObjectType()
func (t *JSTypeBase) ToObjectType() ObjectType {
	k := t.objFlags & objKindMask
	if k >= jsKindUnion {
		return nil
	}
	if k >= objKindFunction && k <= objKindEnum {
		return (*ObjectTypeBase)(unsafe.Pointer(t)).selfObj()
	}
	if objectType, ok := t.self().(ObjectType); ok {
		return objectType
	}
	return nil
}

// Autobox dereferences a type for property access.
//
// Filters null/undefined and autoboxes the resulting type. Never returns
// null.
//
// Java: JSType.autobox()
func (t *JSTypeBase) Autobox() JSType {
	restricted := t.self().RestrictByNotNullOrUndefined()
	autobox := restricted.AutoboxesTo()
	if autobox == nil {
		return restricted
	}
	return autobox
}

// Dereference dereferences a type for property access.
//
// Filters null/undefined, autoboxes the resulting type, and returns it iff
// it's an object. If not an object, returns null.
//
// Java: JSType.dereference()
func (t *JSTypeBase) Dereference() ObjectType {
	return t.self().Autobox().ToObjectType()
}

// CanTestForEqualityWith tests whether this and that are meaningfully
// comparable. By meaningfully, we mean compatible types that do not lead to
// step 22 of the definition of the Abstract Equality Comparison Algorithm
// (11.9.3, page 55–56) of the ECMA-262 specification.
//
// Java: JSType.canTestForEqualityWith(JSType)
func (t *JSTypeBase) CanTestForEqualityWith(that JSType) bool {
	result := t.self().TestForEquality(that)
	jsTypeCheckTriNotNull(result)
	return result == base.Tri_UNKNOWN
}

// jsType_TRI_NULL stands for a null Tri reference.
//
// Go-only: language-adapter — testForEquality returns a nullable Tri in
// Java: the base JSType implementation (testForEqualityHelper) returns null
// when it cannot decide, and the overrides test `result != null`. base.Tri is
// a value type, so null is this out-of-range value, which never escapes the
// package's own overrides.
const jsType_TRI_NULL base.Tri = -1

// jsTypeCheckTriNotNull throws NullPointerException where Java dereferences
// a null Tri.
func jsTypeCheckTriNotNull(t base.Tri) {
	if t == jsType_TRI_NULL {
		panic(jutil.NewNullPointerException())
	}
}

// TestForEquality compares this and that.
//
// Returns
//
//   - Tri#TRUE if the comparison of values of this type and that always
//     succeed (such as undefined compared to null)
//   - Tri#FALSE if the comparison of values of this type and that always
//     fails (such as undefined compared to number)
//   - Tri#UNKNOWN if the comparison can succeed or fail depending on the
//     concrete values
//
// Java: JSType.testForEquality(JSType)
func (t *JSTypeBase) TestForEquality(that JSType) base.Tri {
	return t.self().TestForEqualityHelper(t.self(), that)
}

// Java: JSType.testForEqualityHelper(JSType, JSType)
func (t *JSTypeBase) TestForEqualityHelper(aType JSType, bType JSType) base.Tri {
	if bType.IsAllType() || bType.IsUnknownType() ||
		bType.IsNoResolvedType() ||
		aType.IsAllType() || aType.IsUnknownType() ||
		aType.IsNoResolvedType() {
		return base.Tri_UNKNOWN
	}

	aIsEmpty := aType.IsEmptyType()
	bIsEmpty := bType.IsEmptyType()
	if aIsEmpty || bIsEmpty {
		if aIsEmpty && bIsEmpty {
			return base.Tri_TRUE
		} else {
			return base.Tri_UNKNOWN
		}
	}

	if bType.IsUnionType() {
		return bType.TestForEquality(aType)
	}

	if aType.IsFunctionType() || bType.IsFunctionType() {
		otherType := aType
		if aType.IsFunctionType() {
			otherType = bType
		}

		// TODO(johnlenz): tighten function type comparisons in general.
		if otherType.IsSymbol() {
			return base.Tri_FALSE
		}

		// In theory, functions are comparable to anything except null/undefined.
		// For example, on FF3:
		// function() {} == 'function () {\n}'
		// In practice, how a function serializes to a string is
		// implementation-dependent, so it does not really make sense to test
		// for equality with a string.
		greatestSubtype := otherType.GetGreatestSubtype(t.GetNativeType(JSTypeNative_OBJECT_TYPE))
		if greatestSubtype.IsNoType() || greatestSubtype.IsNoObjectType() {
			return base.Tri_FALSE
		} else {
			return base.Tri_UNKNOWN
		}
	}

	if bType.IsEnumElementType() {
		return bType.TestForEquality(aType)
	}

	// If this is a "Symbol" or that is "symbol" or "Symbol"
	if aType.IsSymbol() {
		if bType.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_TYPE)) ||
			bType.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE)) {
			return base.Tri_UNKNOWN
		}
		return base.Tri_FALSE
	}

	if bType.IsSymbol() {
		if aType.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_TYPE)) ||
			aType.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE)) {
			return base.Tri_UNKNOWN
		}
		return base.Tri_FALSE
	}

	return jsType_TRI_NULL
}

// CanTestForShallowEqualityWith tests whether this and that are
// meaningfully comparable using shallow comparison. By meaningfully, we mean
// compatible types that are not rejected by step 1 of the definition of the
// Strict Equality Comparison Algorithm (11.9.6, page 56–57) of the ECMA-262
// specification.
//
// Java: JSType.canTestForShallowEqualityWith(JSType)
func (t *JSTypeBase) CanTestForShallowEqualityWith(that JSType) bool {
	if t.self().IsNoResolvedType() || that.IsNoResolvedType() {
		return true
	}
	if t.self().IsEmptyType() || that.IsEmptyType() {
		return t.self().IsSubtypeOf(that) || that.IsSubtypeOf(t.self())
	}

	inf := t.self().GetGreatestSubtype(that)
	return !inf.IsEmptyType() ||
		// Our getGreatestSubtype relation on functions is pretty bad. Let's
		// just say it's always ok to compare two functions. Once the TODO in
		// FunctionType is fixed, we should be able to remove this.
		inf == t.GetRegistry().GetNativeType(JSTypeNative_LEAST_FUNCTION_TYPE)
}

// IsNullable tests whether this type is nullable.
//
// Java: JSType.isNullable()
func (t *JSTypeBase) IsNullable() bool {
	return false
}

// IsVoidable tests whether this type is voidable.
//
// Java: JSType.isVoidable()
func (t *JSTypeBase) IsVoidable() bool {
	return false
}

// IsExplicitlyVoidable tests whether this type explicitly allows undefined,
// as opposed to ? or *. This is required for a property to be optional.
//
// Java: JSType.isExplicitlyVoidable()
func (t *JSTypeBase) IsExplicitlyVoidable() bool {
	return false
}

// CollapseUnion gets the least supertype of this that's not a union.
//
// Java: JSType.collapseUnion()
func (t *JSTypeBase) CollapseUnion() JSType {
	return t.self()
}

// GetLeastSupertype gets the least supertype of this and that. The least
// supertype is the join (∨) or supremum of both types in the type lattice.
//
// Examples:
//
//   - number ∨ * = *
//   - number ∨ Object = (number, Object)
//   - Number ∨ Object = Object
//
// Java: JSType.getLeastSupertype(JSType)
func (t *JSTypeBase) GetLeastSupertype(that JSType) JSType {
	if t.self() == that {
		return t.self()
	}

	if that.IsUnionType() {
		// Union types have their own implementation of getLeastSupertype.
		return that.ToMaybeUnionType().GetLeastSupertype(t.self())
	}
	return JSType_GetLeastSupertype(t.self(), that)
}

// JSType_GetLeastSupertype is a generic implementation meant to be used as a
// helper for common getLeastSupertype implementations.
//
// Java: JSType.getLeastSupertype(JSType, JSType)
func JSType_GetLeastSupertype(thisType JSType, thatType JSType) JSType {
	thisTb := jsTypeBase(thisType)
	thatTb := jsTypeBase(thatType)
	if thisTb == thatTb {
		return thisType
	}
	registry := thisTb.GetRegistry()
	if registry != nil && registry.resolver.state == JSTypeResolver_State_CLOSED && thisTb.flags&jsTypeFlagResolvedToSelf != 0 && thatTb != nil {
		if thatTb.flags&jsTypeFlagResolvedToSelf != 0 {
			if cached := registry.lookupLeastSupertype(thisTb, thatTb); cached != nil {
				return cached
			}
			areEquivalent := equalityChecker_AreIdentityEqualWithBases(thisType, thisTb, thatType, thatTb)
			if areEquivalent {
				registry.storeLeastSupertype(thisTb, thatTb, thisType)
				return thisType
			}
			res := registry.CreateUnionType(thisType, thatType)
			resTb := jsTypeBase(res)
			if resTb == thisTb || resTb == thatTb {
				registry.storeLeastSupertype(thisTb, thatTb, res)
			}
			return res
		}
	}
	areEquivalent := jsTypeEquals(thisType, thatType)
	if areEquivalent {
		return thisType
	}
	return registry.CreateUnionType(thisType, thatType)
}

// GetGreatestSubtype gets the greatest subtype of this and that. The
// greatest subtype is the meet (∧) or infimum of both types in the type
// lattice.
//
// Examples
//
//   - Number ∧ Any = Any
//   - number ∧ Object = Any
//   - Number ∧ Object = Number
//
// Java: JSType.getGreatestSubtype(JSType)
func (t *JSTypeBase) GetGreatestSubtype(that JSType) JSType {
	return JSType_GetGreatestSubtype(t.self(), that)
}

// JSType_GetGreatestSubtype is a generic implementation meant to be used as
// a helper for common getGreatestSubtype implementations.
//
// Java: JSType.getGreatestSubtype(JSType, JSType)
func JSType_GetGreatestSubtype(thisType JSType, thatType JSType) JSType {
	if thisType.IsFunctionType() && thatType.IsFunctionType() {
		// The FunctionType sub-lattice is not well-defined. i.e., the
		// proposition
		// A < B => sup(A, B) == B
		// does not hold because of unknown parameters and return types.
		// See the comment in supAndInfHelper for more info on this.
		return thisType.ToMaybeFunctionType().SupAndInfHelper(thatType.ToMaybeFunctionType(), false)
	} else if thisType.Equals(thatType) {
		return thisType
	} else if thisType.IsUnknownType() {
		// Handle unknowns by returning the tightest type possible:
		//  given a known and unknown type, return the known type
		//  given a checked unknown type and some unknown type, return a checked unknown
		//  otherwise return an unchecked unknown
		if !thatType.IsUnknownType() {
			return thatType
		} else if thisType.IsCheckedUnknownType() || thatType.IsCheckedUnknownType() {
			return thisType.GetNativeType(JSTypeNative_CHECKED_UNKNOWN_TYPE)
		} else {
			return thisType.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
	} else if thatType.IsUnknownType() {
		// thisType must be known
		return thisType
	} else if thisType.IsUnionType() {
		return UnionType_GetGreatestSubtype(thisType.ToMaybeUnionType(), thatType)
	} else if thatType.IsUnionType() {
		return UnionType_GetGreatestSubtype(thatType.ToMaybeUnionType(), thisType)
	} else if thisType.IsTemplatizedType() {
		return thisType.ToMaybeTemplatizedType().GetGreatestSubtypeHelper(thatType)
	} else if thatType.IsTemplatizedType() {
		return thatType.ToMaybeTemplatizedType().GetGreatestSubtypeHelper(thisType)
	} else if thisType.IsSubtypeOf(thatType) {
		return thisType
	} else if thatType.IsSubtypeOf(thisType) {
		return thatType
	} else if thisType.IsRecordType() {
		return thisType.ToMaybeRecordType().GetGreatestSubtypeHelper(thatType)
	} else if thatType.IsRecordType() {
		return thatType.ToMaybeRecordType().GetGreatestSubtypeHelper(thisType)
	}

	if thisType.IsEnumElementType() {
		inf := EnumElementType_GetGreatestSubtype(thisType.ToMaybeEnumElementType(), thatType)
		if inf != nil {
			return inf
		}
	} else if thatType.IsEnumElementType() {
		inf := EnumElementType_GetGreatestSubtype(thatType.ToMaybeEnumElementType(), thisType)
		if inf != nil {
			return inf
		}
	}

	if thisType.IsObject() && thatType.IsObject() {
		return thisType.GetNativeType(JSTypeNative_NO_OBJECT_TYPE)
	}
	return thisType.GetNativeType(JSTypeNative_NO_TYPE)
}

// GetRestrictedTypeGivenOutcome computes the restricted type of this type
// knowing that the ToBoolean predicate has a specific value. For more
// information about the ToBoolean predicate, see
// #getPossibleToBooleanOutcomes.
//
// Returns the restricted type, or the Any Type if the underlying type could
// not have yielded this ToBoolean value.
//
// TODO(user): Move this method to the SemanticRAI and use the visit method
// of types to get the restricted type.
//
// Java: JSType.getRestrictedTypeGivenOutcome(Outcome)
func (t *JSTypeBase) GetRestrictedTypeGivenOutcome(outcome Outcome) JSType {
	if outcome.IsTruthy() && t.self() == t.GetNativeType(JSTypeNative_UNKNOWN_TYPE) {
		return t.GetNativeType(JSTypeNative_CHECKED_UNKNOWN_TYPE)
	}

	// Only `NULL_TYPE` and `VOID_TYPE` can really be nullish, so the result
	// for any other non-union type must be NO_TYPE. Note that the UnionType
	// class is responsible for handling the union case.
	if outcome.IsNullish().ToBoolean(false) {
		if t.self() == t.GetNativeType(JSTypeNative_VOID_TYPE) ||
			t.self() == t.GetNativeType(JSTypeNative_NULL_TYPE) {
			return t.self()
		} else {
			return t.GetNativeType(JSTypeNative_NO_TYPE)
		}
	}

	literals := t.self().GetPossibleToBooleanOutcomes()
	if literals.Contains(outcome.IsTruthy()) {
		return t.self()
	} else {
		return t.GetNativeType(JSTypeNative_NO_TYPE)
	}
}

// GetTypesUnderEquality computes the subset of this and that types if
// equality is observed. If a value v1 of type null is equal to a value v2 of
// type (undefined,number), we can infer that the type of v1 is null and the
// type of v2 is undefined.
//
// Returns a pair containing the restricted type of this as the first
// component and the restricted type of that as the second element. The
// returned pair is never null even though its components may be null
//
// Java: JSType.getTypesUnderEquality(JSType)
func (t *JSTypeBase) GetTypesUnderEquality(that JSType) *JSType_TypePair {
	a, b := t.GetTypesUnderEqualityPair(that)
	return NewJSType_TypePair(a, b)
}

func (t *JSTypeBase) GetTypesUnderEqualityPair(that JSType) (JSType, JSType) {
	// unions types
	if that.IsUnionType() {
		a, b := that.ToMaybeUnionType().GetTypesUnderEqualityPair(t.self())
		return b, a
	}

	// other types
	r := t.self().TestForEquality(that)
	jsTypeCheckTriNotNull(r)
	switch r {
	case base.Tri_FALSE:
		return nil, nil
	}
	return t.self(), that // TRUE, UNKNOWN
}

// GetTypesUnderInequality computes the subset of this and that types if
// inequality is observed. If a value v1 of type number is not equal to a
// value v2 of type (undefined,number), we can infer that the type of v1 is
// number and the type of v2 is number as well.
//
// Java: JSType.getTypesUnderInequality(JSType)
func (t *JSTypeBase) GetTypesUnderInequality(that JSType) *JSType_TypePair {
	a, b := t.GetTypesUnderInequalityPair(that)
	return NewJSType_TypePair(a, b)
}

func (t *JSTypeBase) GetTypesUnderInequalityPair(that JSType) (JSType, JSType) {
	// unions types
	if that.IsUnionType() {
		a, b := that.ToMaybeUnionType().GetTypesUnderInequalityPair(t.self())
		return b, a
	}

	// For the remaining types, the only way to restrict anything is if both
	// types are null or void, in which case we're left with nothing.
	if t.IsNullTypeOrVoidType() && jsTypeIsNullTypeOrVoidType(that) {
		noType := t.GetRegistry().GetNativeType(JSTypeNative_NO_TYPE)
		return noType, noType
	}
	return t.self(), that
}

// GetTypesUnderShallowEquality computes the subset of this and that types
// under shallow equality.
//
// Java: JSType.getTypesUnderShallowEquality(JSType)
func (t *JSTypeBase) GetTypesUnderShallowEquality(that JSType) *JSType_TypePair {
	a, b := t.GetTypesUnderShallowEqualityPair(that)
	return NewJSType_TypePair(a, b)
}

func (t *JSTypeBase) GetTypesUnderShallowEqualityPair(that JSType) (JSType, JSType) {
	commonType := t.self().GetGreatestSubtype(that)
	return commonType, commonType
}

// GetTypesUnderShallowInequality computes the subset of this and that types
// under shallow inequality.
//
// Java: JSType.getTypesUnderShallowInequality(JSType)
func (t *JSTypeBase) GetTypesUnderShallowInequality(that JSType) *JSType_TypePair {
	a, b := t.GetTypesUnderShallowInequalityPair(that)
	return NewJSType_TypePair(a, b)
}

func (t *JSTypeBase) GetTypesUnderShallowInequalityPair(that JSType) (JSType, JSType) {
	// union types
	if that.IsUnionType() {
		a, b := that.ToMaybeUnionType().GetTypesUnderShallowInequalityPair(t.self())
		return b, a
	}

	// Other types. There are only two types whose shallow inequality is
	// deterministically true -- null and undefined. We can just enumerate
	// them.
	if t.self() == that && t.IsNullTypeOrVoidType() {
		return nil, nil
	}
	// Since unions have already been removed, if this and that don't have
	// the same unit type then there is no further narrowing that can be done
	// on them.
	return t.self(), that
}

// Java: JSType.IsNullTypeOrVoidType()
func (t *JSTypeBase) IsNullTypeOrVoidType() bool {
	// TODO(sdh): Should this be exposed more publicly?
	return t.self().IsNullType() || t.self().IsVoidType()
}

// jsTypeIsNullTypeOrVoidType is that.IsNullTypeOrVoidType() (a private
// method called on another instance).
func jsTypeIsNullTypeOrVoidType(that JSType) bool {
	return that.GetJSTypeBase().IsNullTypeOrVoidType()
}

// Java: JSType.getUnionMembers()
func (t *JSTypeBase) GetUnionMembers() []JSType {
	if t.self().IsUnionType() {
		return t.self().ToMaybeUnionType().GetAlternates()
	}
	return nil
}

// RestrictByNotNullOrUndefined returns, if this is a union type, a union
// type that does not include the null or undefined type.
//
// Java: JSType.restrictByNotNullOrUndefined()
func (t *JSTypeBase) RestrictByNotNullOrUndefined() JSType {
	return t.self()
}

// RestrictByNotUndefined returns, if this is a union type, a union type
// that does not include the undefined type.
//
// Java: JSType.restrictByNotUndefined()
func (t *JSTypeBase) RestrictByNotUndefined() JSType {
	return t.self()
}

// RestrictByNotNull returns, if this is a union type, a union type that does
// not include the null type.
//
// Java: JSType.restrictByNotNull()
func (t *JSTypeBase) RestrictByNotNull() JSType {
	return t.self()
}

// IsSubtypeWithoutStructuralTyping: the logic of this method is similar to
// isSubtype, except that it does not perform structural interface matching
//
// This function is added for disambiguate properties, and is deprecated for
// the other use cases.
//
// Java: JSType.isSubtypeWithoutStructuralTyping(JSType)
func (t *JSTypeBase) IsSubtypeWithoutStructuralTyping(supertype JSType) bool {
	supBase := jsTypeBase(supertype)
	if supBase == nil {
		failJSTypeNull()
	}
	if t == supBase {
		return true
	}
	c := SubtypeChecker{
		initialSubtype:          t.self(),
		initialSupertype:        supertype,
		registry:                t.GetRegistry(),
		isUsingStructuralTyping: &subtypeChecker_booleans[0],
		subtypingMode:           &subtypeChecker_subtypingModes[JSType_SubtypingMode_NORMAL],
		hasRun:                  true,
	}
	return c.IsSubtypeDispatching(t.self(), supertype)
}

// JSType_SubtypingMode: in files translated from Java, we typecheck null
// and undefined loosely.
type JSType_SubtypingMode int32

const (
	JSType_SubtypingMode_NORMAL JSType_SubtypingMode = iota
	JSType_SubtypingMode_IGNORE_NULL_UNDEFINED
)

var jsType_SubtypingMode_names = [...]string{"NORMAL", "IGNORE_NULL_UNDEFINED"}

// JSType_SubtypingMode_values is SubtypingMode.values().
func JSType_SubtypingMode_values() []JSType_SubtypingMode {
	return []JSType_SubtypingMode{JSType_SubtypingMode_NORMAL, JSType_SubtypingMode_IGNORE_NULL_UNDEFINED}
}

// String is Enum.name().
func (m JSType_SubtypingMode) String() string { return jsType_SubtypingMode_names[m] }

// IsSubtype checks whether this is a subtype of that.
//
// Note this function also returns true if this type structurally matches
// the protocol define by that type (if that type is an interface function
// type)
//
// Subtyping rules:
//
//   - (unknown) — every type is a subtype of the Unknown type.
//   - (no) — the No type is a subtype of every type.
//   - (no-object) — the NoObject type is a subtype of every object type
//     (i.e. subtypes of the Object type).
//   - (ref) — a type is a subtype of itself.
//   - (union-l) — A union type is a subtype of a type U if all the union
//     type's constituents are a subtype of U. Formally (T1, …, Tn) <: U if
//     and only Tk <: U for all k ∈ 1..n.
//   - (union-r) — A type U is a subtype of a union type if it is a subtype
//     of one of the union type's constituents. Formally U <: (T1, …, Tn) if
//     and only if U <: Tk for some index k.
//   - (objects) — an Object O1 is a subtype of an object O2 if it has more
//     properties than O2 and all common properties are pairwise subtypes.
//
// Returns this <: that.
//
// Java: JSType.isSubtype(JSType)
func (t *JSTypeBase) IsSubtype(supertype JSType) bool {
	return t.IsSubtypeOf(supertype)
}

// IsSubtypeWithMode is deprecated: prefer isSubtype(JSType) instead.
//
// Java: JSType.isSubtype(JSType, SubtypingMode)
func (t *JSTypeBase) IsSubtypeWithMode(supertype JSType, mode JSType_SubtypingMode) bool {
	return t.IsSubtypeOfWithMode(supertype, mode)
}

// Java: JSType.isSubtypeOf(JSType)
func (t *JSTypeBase) IsSubtypeOf(supertype JSType) bool {
	supBase := jsTypeBase(supertype)
	if supBase == nil {
		failJSTypeNull()
	}
	if t == supBase {
		return true
	}
	c := SubtypeChecker{
		initialSubtype:          t.self(),
		initialSupertype:        supertype,
		registry:                t.GetRegistry(),
		isUsingStructuralTyping: &subtypeChecker_booleans[1],
		subtypingMode:           &subtypeChecker_subtypingModes[JSType_SubtypingMode_NORMAL],
		hasRun:                  true,
	}
	return c.IsSubtypeDispatching(t.self(), supertype)
}

// Java: JSType.isSubtypeOf(JSType, SubtypingMode)
func (t *JSTypeBase) IsSubtypeOfWithMode(supertype JSType, mode JSType_SubtypingMode) bool {
	supBase := jsTypeBase(supertype)
	if supBase == nil {
		failJSTypeNull()
	}
	if t == supBase {
		return true
	}
	var modePtr *JSType_SubtypingMode
	if mode >= 0 && int(mode) < len(subtypeChecker_subtypingModes) {
		modePtr = &subtypeChecker_subtypingModes[mode]
	} else {
		other := mode
		modePtr = &other
	}
	c := SubtypeChecker{
		initialSubtype:          t.self(),
		initialSupertype:        supertype,
		registry:                t.GetRegistry(),
		isUsingStructuralTyping: &subtypeChecker_booleans[1],
		subtypingMode:           modePtr,
		hasRun:                  true,
	}
	return c.IsSubtypeDispatching(t.self(), supertype)
}

// JSType_VisitForJSType is type.visit(visitor) for a Visitor<JSType>: the
// erased result converted back to JSType.
//
// Go-only: language-adapter — see Visitor.
func JSType_VisitForJSType(type_ JSType, visitor Visitor) JSType {
	if r, ok := visitor.(*TemplateTypeReplacer); ok {
		return r.visitJSType(type_)
	}
	r, _ := type_.Visit(visitor).(JSType)
	return r
}

// Resolve resolves this type in the given scope.
//
// The returned value must be equal to this, as defined by #equals. It may
// or may not be the same object. This method may modify the internal state
// of this, as long as it does so in a way that preserves Object equality.
//
// For efficiency, we should only resolve a type once per compilation job.
// For incremental compilations, one compilation job may need the artifacts
// from a previous generation, so we will eventually need a generational flag
// instead of a boolean one.
//
// Java: JSType.resolve(ErrorReporter)
func (t *JSTypeBase) Resolve(reporter ErrorReporter) JSType {
	t.GetRegistry().GetResolver().AssertLegalToResolveTypes()
	if t.resolveResult == nil {
		// Prevent infinite recursion in cyclically defined types.
		//
		// If resolve is called a twice on a type, before the first call
		// completes, there is a cycle somewhere. The only thing we can do is
		// return `this`. Doing so is safe because the resolve result *should*
		// be equal to `this`, even though the different calls will return
		// different objects.
		t.resolveResult = t
		res := t.self().ResolveInternal(reporter)
		if res != nil {
			t.resolveResult = res.GetJSTypeBase()
		} else {
			t.resolveResult = nil
		}
		if t.resolveResult == t {
			t.flags |= jsTypeFlagResolvedToSelf
		}
		jutil.CheckState(t.resolveResult != nil)
	}
	return t.resolveResult.self()
}

// EagerlyResolveToSelf bypasses the normal call to resolve(ErrorReporter).
// Useful if you are aware that no resolution work needs to be done.
//
// Java: JSType.eagerlyResolveToSelf()
func (t *JSTypeBase) EagerlyResolveToSelf() {
	jutil.CheckStateMsg1(!t.self().IsResolved(), "%s", t.self())
	t.resolveResult = t
	t.flags |= jsTypeFlagResolvedToSelf
	t.GetRegistry().GetResolver().ResolveIfClosed(t.self(), t.self().GetTypeClass())
}

// IsResolved returns whether the type has undergone resolution.
//
// A value of true does not indicate that resolution was successful, only
// that it was attempted and has finished.
//
// Java: JSType.isResolved()
func (t *JSTypeBase) IsResolved() bool {
	return t.resolveResult != nil
}

// IsSuccessfullyResolved returns whether the type has undergone resolution
// and resolved to a "useful" type.
//
// Java: JSType.isSuccessfullyResolved()
func (t *JSTypeBase) IsSuccessfullyResolved() bool {
	return t.self().IsResolved() && !t.self().IsNoResolvedType()
}

// IsUnsuccessfullyResolved returns whether the type has undergone
// resolution and resolved to a "useless" type.
//
// Java: JSType.isUnsuccessfullyResolved()
func (t *JSTypeBase) IsUnsuccessfullyResolved() bool {
	return t.self().IsResolved() && t.self().IsNoResolvedType()
}

// JSType_SafeResolve is a null-safe resolve.
//
// Java: JSType.safeResolve(JSType, ErrorReporter)
func JSType_SafeResolve(type_ JSType, reporter ErrorReporter) JSType {
	if type_ == nil {
		return nil
	}
	return type_.Resolve(reporter)
}

// SetValidator: certain types have constraints on them at resolution-time.
// For example, a type in an @extends annotation must be an object. Clients
// should inject a validator that emits a warning if the type does not
// validate, and return false.
//
// Java: JSType.setValidator(Predicate<JSType>)
func (t *JSTypeBase) SetValidator(validator func(JSType) bool) bool {
	return validator(t.self())
}

// JSType_TypePair represents a pair of types.
type JSType_TypePair struct {
	TypeA JSType
	TypeB JSType
}

// Java: JSType.TypePair(JSType, JSType)
func NewJSType_TypePair(typeA JSType, typeB JSType) *JSType_TypePair {
	return &JSType_TypePair{TypeA: typeA, TypeB: typeB}
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

// MatchConstraint modifies this type so that it matches the specified type.
//
// This is useful for reverse type-inference, where we want to infer that an
// object literal matches its constraint (much like how the java compiler
// does reverse-inference to figure out generics).
//
// Java: JSType.matchConstraint(JSType)
func (t *JSTypeBase) MatchConstraint(constraint JSType) {}

// Java: JSType.toMaybeObjectType()
func (t *JSTypeBase) ToMaybeObjectType() ObjectType {
	return t.self().ToObjectType()
}

// JSType_MatchStatus describes the status of checking that a function
// implicitly implements an interface.
//
// It also be used to describe the status of checking that a record type
// structurally matches another record type
//
// A function implicitly implements an interface if the function does not use
// @implements to declare that it implements the interface, but its class
// structure complies with the protocol defined by the interface
type JSType_MatchStatus int32

const (
	// JSType_MatchStatus_MATCH indicates that a function implicitly
	// implements an interface (i.e., the function structurally complies with
	// the protocol defined in interface) or a record type matches another
	// record type
	JSType_MatchStatus_MATCH JSType_MatchStatus = iota
	// JSType_MatchStatus_NOT_MATCH indicates that a function does not
	// implicitly implements an interface (i.e., the function does not
	// structurally comply with the protocol defined in interface) or a record
	// type does not match another record type
	JSType_MatchStatus_NOT_MATCH
	// JSType_MatchStatus_PROCESSING indicates that the interface and function
	// relationship is under processing
	JSType_MatchStatus_PROCESSING
)

var jsType_MatchStatus_names = [...]string{"MATCH", "NOT_MATCH", "PROCESSING"}

// jsType_MatchStatus_isSubtype is the isSubtype field of each constant.
var jsType_MatchStatus_isSubtype = [...]bool{true, false, true}

// JSType_MatchStatus_values is MatchStatus.values().
func JSType_MatchStatus_values() []JSType_MatchStatus {
	return []JSType_MatchStatus{JSType_MatchStatus_MATCH, JSType_MatchStatus_NOT_MATCH, JSType_MatchStatus_PROCESSING}
}

// String is Enum.name().
func (s JSType_MatchStatus) String() string { return jsType_MatchStatus_names[s] }

// Java: JSType.MatchStatus.subtypeValue()
func (s JSType_MatchStatus) SubtypeValue() bool {
	return jsType_MatchStatus_isSubtype[s]
}

// Java: JSType.MatchStatus.valueOf(boolean)
func JSType_MatchStatus_ValueOf(match bool) JSType_MatchStatus {
	if match {
		return JSType_MatchStatus_MATCH
	}
	return JSType_MatchStatus_NOT_MATCH
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

// JSType_Nullability_values is Nullability.values().
func JSType_Nullability_values() []JSType_Nullability {
	return []JSType_Nullability{JSType_Nullability_EXPLICIT, JSType_Nullability_IMPLICIT}
}

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

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
// Port of com/google/javascript/rhino/jstype/JSType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.rhino.jstype.JSType
type JSType interface {
	Implements_rhino_JSType()
	IsJavaNil_() bool
	GetJSTypeBase() *JSTypeBase
	FastTypeClass() JSTypeClass
	GetKindMask() int32
	GetTypeId() int32
	GetWellFormedTemplatizedKind() int32
	MarkWellFormedTemplatizedOk()
	MayHaveNonStringifiableKey() bool
	HasKindMaskProxy(int32) bool
	HasKindMask(int32) bool
	GetTypeClass() JSTypeClass
	GetNativeType(JSTypeNative) JSType
	GetJSDocInfo() *JSDocInfo
	GetDisplayName() string
	HasDisplayName() bool
	GetPropertyKind_1_Property_Key(Property_Key) JSType_HasPropertyKind
	GetPropertyKind_2_Property_Key_bool(Property_Key, bool) JSType_HasPropertyKind
	GetPropertyKind_2_string_bool(string, bool) JSType_HasPropertyKind
	GetPropertyKind_1_string(string) JSType_HasPropertyKind
	HasProperty_1_string(string) bool
	HasProperty_1_Property_Key(Property_Key) bool
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
	IsUnknownTypeSlow() bool
	IsSomeUnknownType() bool
	IsCheckedUnknownType() bool
	IsUnionType() bool
	IsRawTypeOfTemplatizedType() bool
	IsStruct() bool
	IsDict() bool
	ContainsReferenceAncestor(JSType) bool
	IsLiteralObject() bool
	ToMaybeUnionType() *UnionType
	IsGlobalThisType() bool
	IsFunctionType() bool
	ToMaybeFunctionType() *FunctionType
	AssertFunctionType() *FunctionType
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
	GetTypeParameters() j2gort.Ext_com_google_common_collect_ImmutableList
	GetTemplateParamCount() int32
	MergeSupertypeTemplateTypes(ObjectType)
	IsObject() bool
	IsObjectType() bool
	IsConstructor() bool
	IsNominalType() bool
	IsNominalConstructorOrInterface() bool
	IsNativeObjectType() bool
	LoosenTypecheckingDueToForwardReferencedSupertype() bool
	MaybeLoosenTypecheckingDueToForwardReferencedSupertype(JSType)
	IsInstanceType() bool
	IsInterface() bool
	IsOrdinaryFunction() bool
	Equals(any) bool
	EqualsJSType(JSType) bool
	DiffersFrom(JSType) bool
	HashCode() int32
	RecursionUnsafeHashCode() int32
	MatchesNumberContext() bool
	MatchesStringContext() bool
	MatchesSymbolContext() bool
	MatchesObjectContext() bool
	FindPropertyType_1_Property_Key(Property_Key) JSType
	FindPropertyType_1_string(string) JSType
	FindPropertyTypeWithoutConsideringTemplateTypes_1_string(string) JSType
	FindPropertyTypeWithoutConsideringTemplateTypes_1_Property_Key(Property_Key) JSType
	CanBeCalled() bool
	CanCastTo(JSType) bool
	AutoboxesTo() JSType
	IsBoxableScalar() bool
	ToObjectType() ObjectType
	Autobox() JSType
	Dereference() ObjectType
	CanTestForEqualityWith(JSType) bool
	TestForEquality(JSType) pkg_base.Tri
	TestForEqualityHelper(JSType, JSType) pkg_base.Tri
	CanTestForShallowEqualityWith(JSType) bool
	IsNullable() bool
	IsVoidable() bool
	IsExplicitlyVoidable() bool
	CollapseUnion() JSType
	GetLeastSupertype(JSType) JSType
	GetGreatestSubtype(JSType) JSType
	GetRestrictedTypeGivenOutcome(*Outcome) JSType
	GetPossibleToBooleanOutcomes() BooleanLiteralSet
	GetTypesUnderEquality(JSType) JSType_TypePair
	GetTypesUnderInequality(JSType) JSType_TypePair
	GetTypesUnderShallowEquality(JSType) JSType_TypePair
	GetTypesUnderShallowInequality(JSType) JSType_TypePair
	GetUnionMembers() j2gort.Ext_java_lang_Iterable
	RestrictByNotNullOrUndefined() JSType
	RestrictByNotUndefined() JSType
	RestrictByNotNull() JSType
	IsSubtypeWithoutStructuralTyping(JSType) bool
	IsSubtype_1(JSType) bool
	IsSubtype_2(JSType, JSType_SubtypingMode) bool
	IsSubtypeOf_1(JSType) bool
	IsSubtypeOf_2(JSType, JSType_SubtypingMode) bool
	Visit_1(Visitor) any
	Visit_2(RelationshipVisitor, JSType) any
	Resolve(ErrorReporter) JSType
	GetCachedResolveResult() JSType
	SetCachedResolveResult(JSType)
	ClearCachedResolveResult()
	MarkResolved()
	ResolveInternal(ErrorReporter) JSType
	EagerlyResolveToSelf()
	IsResolved() bool
	IsSuccessfullyResolved() bool
	IsUnsuccessfullyResolved() bool
	SetValidator(j2gort.Ext_com_google_common_base_Predicate) bool
	String() string
	ToAnnotationString(JSType_Nullability) string
	AppendTo(*TypeStringBuilder)
	MatchConstraint(JSType)
	ToMaybeObjectType() ObjectType
}

func ToIface_JSType[T interface {
	JSType
	comparable
}](v T) JSType {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JSType[T JSType](v T) JSType {
	var vi JSType = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.rhino.jstype.JSType
type JSTypeBase struct {
	self_                       JSType
	templateTypeMap             *TemplateTypeMap `gson:"templateTypeMap"`
	registry                    *JSTypeRegistry  `gson:"registry"`
	_tcHandle                   atomic.Uint32
	subSourceNode               Node  `gson:"subSourceNode"`
	typeId                      int32 `gson:"typeId"`
	typeKindMask                int32 `gson:"typeKindMask"`
	cachedHashCode              int32 `gson:"cachedHashCode"`
	jstypeFlags                 int8  `gson:"jstypeFlags"`
	protoFlags                  int8  `gson:"protoFlags"`
	fnFlags                     int8  `gson:"fnFlags"`
	protoTemplateParamCountByte int8  `gson:"protoTemplateParamCountByte"`
}

func (this_ *JSTypeBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(JSType)
}
func (this_ *JSTypeBase) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSType", "JSType"
}
func (this_ *JSTypeBase) TCHandlePtr_() *atomic.Uint32 { return &this_._tcHandle }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSType", "com.google.javascript.rhino.jstype.JSType", func(v_ any) bool {
		return j2gort.InstanceOf[JSType](v_)
	})
}
func (b *JSTypeBase) GetJSTypeBase() *JSTypeBase { return b }
func (*JSTypeBase) Implements_rhino_JSType()     {}

func StaticGet_JSType_SerialVersionUID() int64                     { return int64(1) }
func StaticGet_JSType_JSTYPE_FLAG_LOOSEN_FWD_REF() int8            { return int8(1) }
func StaticGet_JSType_JSTYPE_FLAG_HASH_IN_PROGRESS() int8          { return int8(2) }
func StaticGet_JSType_JSTYPE_FLAG_IN_TEMPLATED_VISIT() int8        { return int8(4) }
func StaticGet_JSType_JSTYPE_FLAG_RESOLVED() int8                  { return int8(8) }
func StaticGet_JSType_JSTYPE_FLAG_OBJ_VISITED() int8               { return int8(16) }
func StaticGet_JSType_JSTYPE_FLAG_OBJ_NOT_UNKNOWN() int8           { return int8(32) }
func StaticGet_JSType_JSTYPE_TEMPLATE_CHECK_SHIFT() int32          { return int32(6) }
func StaticGet_JSType_PROTO_FLAG_WELL_FORMED_TEMPLATIZED_OK() int8 { return int8(-128) }

var static_JSType_CAN_CAST_TO_VISITOR *CanCastToVisitor

func StaticGet_JSType_CAN_CAST_TO_VISITOR() *CanCastToVisitor {
	clinit_JSType()
	return static_JSType_CAN_CAST_TO_VISITOR
}
func StaticSet_JSType_CAN_CAST_TO_VISITOR(v *CanCastToVisitor) {
	clinit_JSType()
	static_JSType_CAN_CAST_TO_VISITOR = v
}
func StaticGet_JSType_KIND_PROXY() int32                      { return int32(1) }
func StaticGet_JSType_KIND_NO_TYPE() int32                    { return int32(2) }
func StaticGet_JSType_KIND_NO_OBJECT_TYPE() int32             { return int32(4) }
func StaticGet_JSType_KIND_NO_RESOLVED_TYPE() int32           { return int32(8) }
func StaticGet_JSType_KIND_ALL() int32                        { return int32(16) }
func StaticGet_JSType_KIND_CHECKED_UNKNOWN() int32            { return int32(32) }
func StaticGet_JSType_KIND_UNION() int32                      { return int32(64) }
func StaticGet_JSType_KIND_FUNCTION() int32                   { return int32(128) }
func StaticGet_JSType_KIND_ENUM_ELEMENT() int32               { return int32(256) }
func StaticGet_JSType_KIND_ENUM() int32                       { return int32(512) }
func StaticGet_JSType_KIND_NAMED() int32                      { return int32(1024) }
func StaticGet_JSType_KIND_RECORD() int32                     { return int32(2048) }
func StaticGet_JSType_KIND_TEMPLATIZED() int32                { return int32(4096) }
func StaticGet_JSType_KIND_TEMPLATE() int32                   { return int32(8192) }
func StaticGet_JSType_KIND_NUMBER_VALUE() int32               { return int32(16384) }
func StaticGet_JSType_KIND_STRING_VALUE() int32               { return int32(32768) }
func StaticGet_JSType_KIND_BOOLEAN_VALUE() int32              { return int32(65536) }
func StaticGet_JSType_KIND_NULL() int32                       { return int32(131072) }
func StaticGet_JSType_KIND_VOID() int32                       { return int32(262144) }
func StaticGet_JSType_KIND_BIGINT_VALUE() int32               { return int32(524288) }
func StaticGet_JSType_KIND_SYMBOL_VALUE() int32               { return int32(1048576) }
func StaticGet_JSType_KIND_KNOWN_SYMBOL_VALUE() int32         { return int32(2097152) }
func StaticGet_JSType_KIND_CONSTRUCTOR() int32                { return int32(4194304) }
func StaticGet_JSType_KIND_INTERFACE() int32                  { return int32(8388608) }
func StaticGet_JSType_KIND_ORDINARY_FUNCTION() int32          { return int32(16777216) }
func StaticGet_JSType_KIND_NATIVE_OBJECT() int32              { return int32(33554432) }
func StaticGet_JSType_KIND_HAS_NON_EMPTY_TEMPLATE_MAP() int32 { return int32(67108864) }
func StaticGet_JSType_TYPE_CLASS_SHIFT() int32                { return int32(27) }

var static_JSType_TYPE_CLASSES []JSTypeClass

func StaticGet_JSType_TYPE_CLASSES() []JSTypeClass {
	clinit_JSType()
	return static_JSType_TYPE_CLASSES
}
func StaticSet_JSType_TYPE_CLASSES(v []JSTypeClass) {
	clinit_JSType()
	static_JSType_TYPE_CLASSES = v
}
func StaticGet_JSType_MAY_HAVE_NON_EMPTY_TEMPLATE_MAP_MASK() int32 { return int32(67108865) }
func StaticGet_JSType_EMPTY_TYPE_MASK() int32                      { return int32(14) }
func StaticGet_JSType_WELL_FORMED_TEMPLATIZED_ACTIVE_MASK() int32  { return int32(13760) }
func StaticGet_JSType_NON_STRINGIFIABLE_KEY_CANDIDATE_MASK() int32 { return int32(16781377) }
func StaticGet_JSType_CAN_RESOLVE_TO_OTHER_MASK() int32            { return int32(5184) }

var clinitDone_JSType uint32
var clinitState_JSType j2gort.ClinitState

func clinit_JSType() {
	if atomic.LoadUint32(&clinitDone_JSType) != 0 {
		return
	}
	clinitSlow_JSType()
}

func clinitSlow_JSType() {
	if !j2gort.EnterClinit(&clinitState_JSType, &clinitDone_JSType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSType, &clinitDone_JSType)
	static_JSType_CAN_CAST_TO_VISITOR = NewCanCastToVisitor()
	static_JSType_TYPE_CLASSES = JSTypeClass_Values()
}

func (this_ *JSTypeBase) initFields_JSType() {
	_ = this_
	this_.jstypeFlags = int8(int32(0))
	this_.protoFlags = int8(int32(0))
	this_.fnFlags = int8(int32(0))
	this_.protoTemplateParamCountByte = int8(int32(0))
	this_.cachedHashCode = int32(0)
}

func InitJSType_1(this_ *JSTypeBase, registry *JSTypeRegistry) {
	_ = registry
	InitJSType_3(this_, registry, nil, int32(0))
}

func InitJSType_2(this_ *JSTypeBase, registry *JSTypeRegistry, templateTypeMap *TemplateTypeMap) {
	_ = registry
	_ = templateTypeMap
	InitJSType_3(this_, registry, templateTypeMap, int32(0))
}

func InitJSType_3(this_ *JSTypeBase, registry *JSTypeRegistry, templateTypeMap *TemplateTypeMap, typeKindMask int32) {
	_ = registry
	_ = templateTypeMap
	_ = typeKindMask
	this_.initFields_JSType()
	this_.registry = registry
	if registry != nil {
		this_.typeId = func() int32 {
			r_1 := registry
			old_2 := r_1.nextTypeId
			nw_3 := old_2 + 1
			r_1.nextTypeId = nw_3
			return nw_3
		}()
	} else {
		this_.typeId = int32(0)
	}
	this_.typeKindMask = typeKindMask
	this_.typeKindMask = (typeKindMask | (javalang.Shl32((this_.self_.GetTypeClass()).Ordinal(), int(int32(27)))))
	var tm *TemplateTypeMap
	_ = tm
	if templateTypeMap == nil {
		tm = registry.GetEmptyTemplateTypeMap()
	} else {
		tm = templateTypeMap
	}
	this_.templateTypeMap = tm
	if (tm != nil) && !(tm.IsEmpty()) {
		this_.typeKindMask = (this_.typeKindMask | int32(67108864))
	}
	registry.GetResolver().AddUnresolved(this_.self_)
}

// Java method: JSType.fastTypeClass()
func (this_ *JSTypeBase) FastTypeClass() JSTypeClass {
	return StaticGet_JSType_TYPE_CLASSES()[javalang.Ushr32(this_.typeKindMask, int(int32(27)))]
}

// Java method: JSType.getKindMask()
func (this_ *JSTypeBase) GetKindMask() int32 {
	return this_.typeKindMask
}

// Java method: JSType.getTypeId()
func (this_ *JSTypeBase) GetTypeId() int32 {
	return this_.typeId
}

// Java method: JSType.getWellFormedTemplatizedKind()
func (this_ *JSTypeBase) GetWellFormedTemplatizedKind() int32 {
	if (int32(this_.protoFlags) & int32(int8(-128))) != int32(0) {
		return int32(0)
	}
	return (this_.typeKindMask & int32(13760))
}

// Java method: JSType.markWellFormedTemplatizedOk()
func (this_ *JSTypeBase) MarkWellFormedTemplatizedOk() {
	if this_.registry.IsResolverClosed() {
		this_.protoFlags = int8((int32(this_.protoFlags) | int32(int8(-128))))
	}
}

// Java method: JSType.mayHaveNonStringifiableKey()
func (this_ *JSTypeBase) MayHaveNonStringifiableKey() bool {
	return ((this_.typeKindMask & int32(16781377)) != int32(0))
}

// Java method: JSType.hasKindMaskProxy(int mask)
func (this_ *JSTypeBase) HasKindMaskProxy(mask int32) bool {
	_ = mask
	return false
}

// Java method: JSType.hasKindMask(int mask)
func (this_ *JSTypeBase) HasKindMask(mask int32) bool {
	_ = mask
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (mask | int32(1))) == int32(0) {
		return false
	}
	return (((m & mask) != int32(0)) || this_.self_.HasKindMaskProxy(mask))
}

// Java method: JSType.getTypeClass()
func (this_ *JSTypeBase) GetTypeClass() JSTypeClass {
	panic("abstract or native method invoked")
}

// Java method: JSType.getNativeType(com.google.javascript.rhino.jstype.JSTypeNative typeId)
func (this_ *JSTypeBase) GetNativeType(typeId JSTypeNative) JSType {
	_ = typeId
	return this_.registry.GetNativeType(typeId)
}

// Java method: JSType.getJSDocInfo()
func (this_ *JSTypeBase) GetJSDocInfo() *JSDocInfo {
	return nil
}

// Java method: JSType.getDisplayName()
func (this_ *JSTypeBase) GetDisplayName() string {
	return j2gort.NullString
}

// Java method: JSType.hasDisplayName()
func (this_ *JSTypeBase) HasDisplayName() bool {
	var displayName string = this_.self_.GetDisplayName()
	_ = displayName
	return ((displayName != j2gort.NullString) && !(j2gort.CheckNotNullString(displayName) == ""))
}

// Java method: JSType.getPropertyKind(com.google.javascript.rhino.jstype.Property.Key pname)
func (this_ *JSTypeBase) GetPropertyKind_1_Property_Key(pname Property_Key) JSType_HasPropertyKind {
	_ = pname
	return this_.self_.GetPropertyKind_2_Property_Key_bool(pname, true)
}

// Java method: JSType.getPropertyKind(com.google.javascript.rhino.jstype.Property.Key pname, boolean autobox)
func (this_ *JSTypeBase) GetPropertyKind_2_Property_Key_bool(pname Property_Key, autobox bool) JSType_HasPropertyKind {
	_ = pname
	_ = autobox
	return StaticGet_JSType_HasPropertyKind_ABSENT()
}

// Java method: JSType.getPropertyKind(java.lang.String pname, boolean autobox)
func (this_ *JSTypeBase) GetPropertyKind_2_string_bool(pname string, autobox bool) JSType_HasPropertyKind {
	_ = pname
	_ = autobox
	j2gort.CheckNotNullString(pname)
	return StaticGet_JSType_HasPropertyKind_ABSENT()
}

// Java method: JSType.getPropertyKind(java.lang.String pname)
func (this_ *JSTypeBase) GetPropertyKind_1_string(pname string) JSType_HasPropertyKind {
	_ = pname
	return this_.self_.GetPropertyKind_2_string_bool(pname, true)
}

// Java method: JSType.hasProperty(java.lang.String pname)
func (this_ *JSTypeBase) HasProperty_1_string(pname string) bool {
	_ = pname
	return !(func() bool {
		re_4 := this_.self_.GetPropertyKind_2_string_bool(pname, false)
		ae_5 := StaticGet_JSType_HasPropertyKind_ABSENT()
		if re_4 == 0 {
			j2gort.ThrowNPE()
		}
		return re_4 == ae_5
	}())
}

// Java method: JSType.hasProperty(com.google.javascript.rhino.jstype.Property.Key pname)
func (this_ *JSTypeBase) HasProperty_1_Property_Key(pname Property_Key) bool {
	_ = pname
	return !(func() bool {
		re_6 := this_.self_.GetPropertyKind_2_Property_Key_bool(pname, false)
		ae_7 := StaticGet_JSType_HasPropertyKind_ABSENT()
		if re_6 == 0 {
			j2gort.ThrowNPE()
		}
		return re_6 == ae_7
	}())
}

// Java method: JSType.isNoType()
func (this_ *JSTypeBase) IsNoType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(2) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(2)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(2)))
}

// Java method: JSType.isNoResolvedType()
func (this_ *JSTypeBase) IsNoResolvedType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(8) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(8)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(8)))
}

// Java method: JSType.isNoObjectType()
func (this_ *JSTypeBase) IsNoObjectType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(4) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(4)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(4)))
}

// Java method: JSType.isEmptyType()
func (this_ *JSTypeBase) IsEmptyType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & ((int32(14) | int32(1)) | int32(128))) == int32(0) {
		return false
	}
	if (m & int32(14)) != int32(0) {
		return true
	}
	if (m & int32(1)) != int32(0) {
		return this_.self_.HasKindMaskProxy(int32(14))
	}
	return j2gort.IfaceEquals[JSType](this_.self_, this_.registry.GetNativeFunctionType(StaticGet_JSTypeNative_LEAST_FUNCTION_TYPE()))
}

// Java method: JSType.isNumberObjectType()
func (this_ *JSTypeBase) IsNumberObjectType() bool {
	return false
}

// Java method: JSType.isNumberValueType()
func (this_ *JSTypeBase) IsNumberValueType() bool {
	return ((this_.typeKindMask & int32(16384)) != int32(0))
}

// Java method: JSType.isBigIntObjectType()
func (this_ *JSTypeBase) IsBigIntObjectType() bool {
	return false
}

// Java method: JSType.isBigIntValueType()
func (this_ *JSTypeBase) IsBigIntValueType() bool {
	return ((this_.typeKindMask & int32(524288)) != int32(0))
}

// Java method: JSType.isFunctionPrototypeType()
func (this_ *JSTypeBase) IsFunctionPrototypeType() bool {
	return false
}

// Java method: JSType.isStringObjectType()
func (this_ *JSTypeBase) IsStringObjectType() bool {
	return false
}

// Java method: JSType.isSymbolObjectType()
func (this_ *JSTypeBase) IsSymbolObjectType() bool {
	return false
}

// Java method: JSType.isTheObjectType()
func (this_ *JSTypeBase) IsTheObjectType() bool {
	return false
}

// Java method: JSType.isStringValueType()
func (this_ *JSTypeBase) IsStringValueType() bool {
	return ((this_.typeKindMask & int32(32768)) != int32(0))
}

// Java method: JSType.isSymbolValueType()
func (this_ *JSTypeBase) IsSymbolValueType() bool {
	return ((this_.typeKindMask & int32(1048576)) != int32(0))
}

// Java method: JSType.isKnownSymbolValueType()
func (this_ *JSTypeBase) IsKnownSymbolValueType() bool {
	return ((this_.typeKindMask & int32(2097152)) != int32(0))
}

// Java method: JSType.isString()
func (this_ *JSTypeBase) IsString() bool {
	return (this_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_STRING_TYPE())) || this_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_STRING_OBJECT_TYPE())))
}

// Java method: JSType.isNumber()
func (this_ *JSTypeBase) IsNumber() bool {
	return (this_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_NUMBER_TYPE())) || this_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_NUMBER_OBJECT_TYPE())))
}

// Java method: JSType.isSymbol()
func (this_ *JSTypeBase) IsSymbol() bool {
	return (this_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_TYPE())) || this_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_OBJECT_TYPE())))
}

// Java method: JSType.isBigIntOrNumber()
func (this_ *JSTypeBase) IsBigIntOrNumber() bool {
	return (this_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_BIGINT_NUMBER())) || this_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_BIGINT_NUMBER_OBJECT())))
}

// Java method: JSType.isOnlyBigInt()
func (this_ *JSTypeBase) IsOnlyBigInt() bool {
	return (this_.IsBigIntValueType() || this_.self_.IsBigIntObjectType())
}

// Java method: JSType.isArrayType()
func (this_ *JSTypeBase) IsArrayType() bool {
	return false
}

// Java method: JSType.isReadonlyArrayType()
func (this_ *JSTypeBase) IsReadonlyArrayType() bool {
	return false
}

// Java method: JSType.isBooleanObjectType()
func (this_ *JSTypeBase) IsBooleanObjectType() bool {
	return false
}

// Java method: JSType.isBooleanValueType()
func (this_ *JSTypeBase) IsBooleanValueType() bool {
	return ((this_.typeKindMask & int32(65536)) != int32(0))
}

// Java method: JSType.isRegexpType()
func (this_ *JSTypeBase) IsRegexpType() bool {
	return false
}

// Java method: JSType.isDateType()
func (this_ *JSTypeBase) IsDateType() bool {
	return false
}

// Java method: JSType.isNullType()
func (this_ *JSTypeBase) IsNullType() bool {
	return ((this_.typeKindMask & int32(131072)) != int32(0))
}

// Java method: JSType.isVoidType()
func (this_ *JSTypeBase) IsVoidType() bool {
	return ((this_.typeKindMask & int32(262144)) != int32(0))
}

// Java method: JSType.isAllType()
func (this_ *JSTypeBase) IsAllType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(16) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(16)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(16)))
}

// Java method: JSType.isUnknownType()
func (this_ *JSTypeBase) IsUnknownType() bool {
	if (int32(this_.jstypeFlags) & int32(int8(32))) != int32(0) {
		return false
	}
	return this_.self_.IsUnknownTypeSlow()
}

// Java method: JSType.isUnknownTypeSlow()
func (this_ *JSTypeBase) IsUnknownTypeSlow() bool {
	return false
}

// Java method: JSType.isSomeUnknownType()
func (this_ *JSTypeBase) IsSomeUnknownType() bool {
	return this_.IsUnknownType()
}

// Java method: JSType.isCheckedUnknownType()
func (this_ *JSTypeBase) IsCheckedUnknownType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(32) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(32)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(32)))
}

// Java method: JSType.isUnionType()
func (this_ *JSTypeBase) IsUnionType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(64) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(64)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(64)))
}

// Java method: JSType.isRawTypeOfTemplatizedType()
func (this_ *JSTypeBase) IsRawTypeOfTemplatizedType() bool {
	return ((this_.self_.GetTemplateParamCount() > int32(0)) && !(this_.IsTemplatizedType()))
}

// Java method: JSType.isStruct()
func (this_ *JSTypeBase) IsStruct() bool {
	if this_.self_.IsObject() {
		var objType ObjectType = this_.ToObjectType()
		_ = objType
		var ctor *FunctionType = objType.GetConstructor()
		_ = ctor
		if ctor == nil {
			var info *JSDocInfo = objType.GetJSDocInfo()
			_ = info
			return ((info != nil) && info.MakesStructs())
		} else {
			return ctor.MakesStructs()
		}
	}
	return false
}

// Java method: JSType.isDict()
func (this_ *JSTypeBase) IsDict() bool {
	if this_.self_.IsObject() {
		var objType ObjectType = this_.ToObjectType()
		_ = objType
		var ctor *FunctionType = objType.GetConstructor()
		_ = ctor
		if ctor == nil {
			var info *JSDocInfo = objType.GetJSDocInfo()
			_ = info
			return ((info != nil) && info.MakesDicts())
		} else {
			return ctor.MakesDicts()
		}
	}
	return false
}

// Java method: JSType.containsReferenceAncestor(com.google.javascript.rhino.jstype.JSType target)
func (this_ *JSTypeBase) ContainsReferenceAncestor(target JSType) bool {
	_ = target
	var visitor *ContainsUpperBoundSuperTypeVisitor = NewContainsUpperBoundSuperTypeVisitor(target)
	_ = visitor
	return (j2gort.Cast[ContainsUpperBoundSuperTypeVisitor_Result](this_.self_.Visit_1(visitor)) == StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_PRESENT())
}

// Java method: JSType.isLiteralObject()
func (this_ *JSTypeBase) IsLiteralObject() bool {
	var prototypeObjectType PrototypeObjectType
	_ = prototypeObjectType
	var ok_8 bool
	prototypeObjectType, ok_8 = j2gort.TryCast[PrototypeObjectType](this_.self_)
	if ok_8 {
		return prototypeObjectType.GetPrototypeObjectTypeBase().IsAnonymous()
	}
	return false
}

// Java method: JSType.toMaybeUnionType()
func (this_ *JSTypeBase) ToMaybeUnionType() *UnionType {
	return nil
}

// Java method: JSType.isGlobalThisType()
func (this_ *JSTypeBase) IsGlobalThisType() bool {
	return j2gort.IfaceEquals[JSType](this_.self_, this_.registry.GetNativeType(StaticGet_JSTypeNative_GLOBAL_THIS()))
}

// Java method: JSType.isFunctionType()
func (this_ *JSTypeBase) IsFunctionType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(128) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(128)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(128)))
}

// Java method: JSType.toMaybeFunctionType()
func (this_ *JSTypeBase) ToMaybeFunctionType() *FunctionType {
	return nil
}

// Java method: JSType.assertFunctionType()
func (this_ *JSTypeBase) AssertFunctionType() *FunctionType {
	var result *FunctionType = j2gort.CheckNotNullPtr(this_.self_.ToMaybeFunctionType())
	_ = result
	return result
}

// Java method: JSType.assertObjectType()
func (this_ *JSTypeBase) AssertObjectType() ObjectType {
	var result ObjectType = j2gort.CheckNotNullIface(this_.ToMaybeObjectType())
	_ = result
	return result
}

// Java static method: JSType.toMaybeFunctionType(com.google.javascript.rhino.jstype.JSType type)
func JSType_ToMaybeFunctionType(type_ JSType) *FunctionType {
	clinit_JSType()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return nil
	} else {
		return type_.ToMaybeFunctionType()
	}
}

// Java method: JSType.toMaybeKnownSymbolType()
func (this_ *JSTypeBase) ToMaybeKnownSymbolType() *KnownSymbolType {
	return nil
}

// Java method: JSType.isEnumElementType()
func (this_ *JSTypeBase) IsEnumElementType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(256) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(256)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(256)))
}

// Java method: JSType.getEnumeratedTypeOfEnumElement()
func (this_ *JSTypeBase) GetEnumeratedTypeOfEnumElement() JSType {
	var e *EnumElementType = this_.self_.ToMaybeEnumElementType()
	_ = e
	if e == nil {
		return nil
	} else {
		return e.GetPrimitiveType()
	}
}

// Java method: JSType.toMaybeEnumElementType()
func (this_ *JSTypeBase) ToMaybeEnumElementType() *EnumElementType {
	return nil
}

// Java method: JSType.isEnumType()
func (this_ *JSTypeBase) IsEnumType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(512) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(512)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(512)))
}

// Java method: JSType.toMaybeEnumType()
func (this_ *JSTypeBase) ToMaybeEnumType() *EnumType {
	return nil
}

// Java method: JSType.isNamedType()
func (this_ *JSTypeBase) IsNamedType() bool {
	return ((this_.typeKindMask & int32(1024)) != int32(0))
}

// Java method: JSType.toMaybeNamedType()
func (this_ *JSTypeBase) ToMaybeNamedType() *NamedType {
	return nil
}

// Java method: JSType.isRecordType()
func (this_ *JSTypeBase) IsRecordType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(2048) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(2048)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(2048)))
}

// Java method: JSType.isStructuralInterface()
func (this_ *JSTypeBase) IsStructuralInterface() bool {
	return false
}

// Java method: JSType.isStructuralType()
func (this_ *JSTypeBase) IsStructuralType() bool {
	return false
}

// Java method: JSType.toMaybeRecordType()
func (this_ *JSTypeBase) ToMaybeRecordType() *RecordType {
	return nil
}

// Java method: JSType.isTemplatizedType()
func (this_ *JSTypeBase) IsTemplatizedType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(4096) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(4096)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(4096)))
}

// Java method: JSType.toMaybeTemplatizedType()
func (this_ *JSTypeBase) ToMaybeTemplatizedType() *TemplatizedType {
	return nil
}

// Java method: JSType.isTemplateType()
func (this_ *JSTypeBase) IsTemplateType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(8192) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(8192)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(8192)))
}

// Java method: JSType.toMaybeTemplateType()
func (this_ *JSTypeBase) ToMaybeTemplateType() *TemplateType {
	return nil
}

// Java method: JSType.hasAnyTemplateTypes()
func (this_ *JSTypeBase) HasAnyTemplateTypes() bool {
	var cached int32 = ((javalang.Ushr32(int32(this_.jstypeFlags), int(int32(6)))) & int32(3))
	_ = cached
	if cached != int32(0) {
		return (cached == int32(2))
	}
	if (int32(this_.jstypeFlags) & int32(int8(4))) != int32(0) {
		return false
	}
	this_.jstypeFlags = int8((int32(this_.jstypeFlags) | int32(int8(4))))
	var result bool = this_.self_.HasAnyTemplateTypesInternal()
	_ = result
	this_.jstypeFlags = int8((int32(this_.jstypeFlags) & ^(int32(int8(4)))))
	if this_.IsResolved() || this_.registry.IsResolverClosed() {
		var code int32
		_ = code
		if result {
			code = int32(2)
		} else {
			code = int32(1)
		}
		this_.jstypeFlags = int8(((int32(this_.jstypeFlags) & ^(javalang.Shl32(int32(3), int(int32(6))))) | (javalang.Shl32(code, int(int32(6))))))
	}
	return result
}

// Java method: JSType.hasAnyTemplateTypesInternal()
func (this_ *JSTypeBase) HasAnyTemplateTypesInternal() bool {
	return this_.self_.GetTemplateTypeMap().HasAnyTemplateTypesInternal()
}

// Java method: JSType.getTemplateTypeMap()
func (this_ *JSTypeBase) GetTemplateTypeMap() *TemplateTypeMap {
	return this_.templateTypeMap
}

// Java method: JSType.getTypeParameters()
func (this_ *JSTypeBase) GetTypeParameters() j2gort.Ext_com_google_common_collect_ImmutableList {
	var map_ *TemplateTypeMap = this_.self_.GetTemplateTypeMap()
	_ = map_
	var size int32 = map_.Size()
	_ = size
	return map_.GetTemplateKeysSlice((size - this_.self_.GetTemplateParamCount()), size)
}

// Java method: JSType.getTemplateParamCount()
func (this_ *JSTypeBase) GetTemplateParamCount() int32 {
	return int32(0)
}

// Java method: JSType.mergeSupertypeTemplateTypes(com.google.javascript.rhino.jstype.ObjectType other)
func (this_ *JSTypeBase) MergeSupertypeTemplateTypes(other ObjectType) {
	_ = other
	this_.MaybeLoosenTypecheckingDueToForwardReferencedSupertype(Upcast_JSType(other))
	if other.GetObjectTypeBase().IsRawTypeOfTemplatizedType() {
		other = ToIface_ObjectType(this_.registry.CreateTemplatizedType_2_ObjectType_Ext_com_google_common_collect_ImmutableList(other, j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()))
	}
	this_.templateTypeMap = other.GetTemplateTypeMap().CopyWithExtension_1(this_.self_.GetTemplateTypeMap())
	if !(this_.templateTypeMap.IsEmpty()) {
		this_.typeKindMask = (this_.typeKindMask | int32(67108864))
	}
}

// Java method: JSType.isObject()
func (this_ *JSTypeBase) IsObject() bool {
	return false
}

// Java method: JSType.isObjectType()
func (this_ *JSTypeBase) IsObjectType() bool {
	return this_.self_.IsObject()
}

// Java method: JSType.isConstructor()
func (this_ *JSTypeBase) IsConstructor() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(4194304) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(4194304)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(4194304)))
}

// Java method: JSType.isNominalType()
func (this_ *JSTypeBase) IsNominalType() bool {
	return false
}

// Java method: JSType.isNominalConstructorOrInterface()
func (this_ *JSTypeBase) IsNominalConstructorOrInterface() bool {
	if this_.HasKindMask((int32(4194304) | int32(8388608))) {
		var fn *FunctionType = this_.self_.ToMaybeFunctionType()
		_ = fn
		if fn == nil {
			return false
		}
		if fn.GetSource() != 0 {
			return true
		}
		return fn.IsNativeObjectType()
	}
	return false
}

// Java method: JSType.isNativeObjectType()
func (this_ *JSTypeBase) IsNativeObjectType() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(33554432) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(33554432)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(33554432)))
}

// Java method: JSType.loosenTypecheckingDueToForwardReferencedSupertype()
func (this_ *JSTypeBase) LoosenTypecheckingDueToForwardReferencedSupertype() bool {
	return ((int32(this_.jstypeFlags) & int32(int8(1))) != int32(0))
}

// Java method: JSType.maybeLoosenTypecheckingDueToForwardReferencedSupertype(com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *JSTypeBase) MaybeLoosenTypecheckingDueToForwardReferencedSupertype(supertype JSType) {
	_ = supertype
	base_supertype_ := supertype.GetJSTypeBase()
	if (base_supertype_.IsNamedType() && !(base_supertype_.IsResolved())) || supertype.LoosenTypecheckingDueToForwardReferencedSupertype() {
		this_.jstypeFlags = int8((int32(this_.jstypeFlags) | int32(int8(1))))
	}
}

// Java method: JSType.isInstanceType()
func (this_ *JSTypeBase) IsInstanceType() bool {
	return false
}

// Java method: JSType.isInterface()
func (this_ *JSTypeBase) IsInterface() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(8388608) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(8388608)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(8388608)))
}

// Java method: JSType.isOrdinaryFunction()
func (this_ *JSTypeBase) IsOrdinaryFunction() bool {
	var m int32 = this_.typeKindMask
	_ = m
	if (m & (int32(16777216) | int32(1))) == int32(0) {
		return false
	}
	return (((m & int32(16777216)) != int32(0)) || this_.self_.HasKindMaskProxy(int32(16777216)))
}

// Java method: JSType.equals(java.lang.@org.jspecify.annotations.Nullable Object other)
func (this_ *JSTypeBase) Equals(arg_9 any) bool {
	other := arg_9
	_ = other
	if j2gort.RefEquals(this_.self_, other) {
		return true
	}
	var jSType JSType
	_ = jSType
	var ok_10 bool
	jSType, ok_10 = j2gort.TryCast[JSType](other)
	return ok_10 && (EqualityChecker_Check(StaticGet_EqualityChecker_EqMethod_IDENTITY(), this_.self_, jSType))
}

// Java method: JSType.equalsJSType(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType that)
func (this_ *JSTypeBase) EqualsJSType(that JSType) bool {
	_ = that
	if this_.self_ == that || ((this_.self_ == nil || this_.self_.IsJavaNil_()) && (that == nil || that.IsJavaNil_())) {
		return true
	}
	if that == nil || that.IsJavaNil_() {
		return false
	}
	return EqualityChecker_Check(StaticGet_EqualityChecker_EqMethod_IDENTITY(), this_.self_, that)
}

// Java method: JSType.differsFrom(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) DiffersFrom(that JSType) bool {
	_ = that
	return !(EqualityChecker_Check(StaticGet_EqualityChecker_EqMethod_DATA_FLOW(), this_.self_, that))
}

// Java static method: JSType.areSimilar(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType first, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType second)
func JSType_AreSimilar(first JSType, second JSType) bool {
	clinit_JSType()
	_ = first
	_ = second
	if (first == nil || first.IsJavaNil_()) || (second == nil || second.IsJavaNil_()) {
		return first == second || ((first == nil || first.IsJavaNil_()) && (second == nil || second.IsJavaNil_()))
	}
	return !(first.GetJSTypeBase().DiffersFrom(second))
}

// Java method: JSType.hashCode()
func (this_ *JSTypeBase) HashCode() int32 {
	if (int32(this_.jstypeFlags) & int32(int8(2))) != int32(0) {
		return int32(-1)
	}
	var reg *JSTypeRegistry = this_.registry
	_ = reg
	var canCache bool = ((((reg != nil) && (reg.hashDepth == int32(0))) && ((int32(this_.jstypeFlags) & int32(int8(8))) != int32(0))) && reg.GetResolver().IsClosed())
	_ = canCache
	if canCache {
		var cached int32 = this_.cachedHashCode
		_ = cached
		if cached != int32(0) {
			return cached
		}
	}
	this_.jstypeFlags = int8((int32(this_.jstypeFlags) | int32(int8(2))))
	if reg != nil {
		r_11 := reg
		r_11.hashDepth = r_11.hashDepth + 1
	}
	var hashCode int32 = int32(0)
	_ = hashCode
	func() {
		defer func() {
			if reg != nil {
				r_12 := reg
				r_12.hashDepth = r_12.hashDepth - 1
			}
			this_.jstypeFlags = int8((int32(this_.jstypeFlags) & ^(int32(int8(2)))))
		}()
		hashCode = this_.self_.RecursionUnsafeHashCode()
	}()
	if canCache && (hashCode != int32(0)) {
		this_.cachedHashCode = hashCode
	}
	return hashCode
}

// Java method: JSType.recursionUnsafeHashCode()
func (this_ *JSTypeBase) RecursionUnsafeHashCode() int32 {
	panic("abstract or native method invoked")
}

// Java method: JSType.matchesNumberContext()
func (this_ *JSTypeBase) MatchesNumberContext() bool {
	return false
}

// Java method: JSType.matchesStringContext()
func (this_ *JSTypeBase) MatchesStringContext() bool {
	return false
}

// Java method: JSType.matchesSymbolContext()
func (this_ *JSTypeBase) MatchesSymbolContext() bool {
	return false
}

// Java method: JSType.matchesObjectContext()
func (this_ *JSTypeBase) MatchesObjectContext() bool {
	return false
}

// Java method: JSType.findPropertyType(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *JSTypeBase) FindPropertyType_1_Property_Key(propertyName Property_Key) JSType {
	_ = propertyName
	var propertyType JSType = this_.self_.FindPropertyTypeWithoutConsideringTemplateTypes_1_Property_Key(propertyName)
	_ = propertyType
	if propertyType == nil || propertyType.IsJavaNil_() {
		return nil
	}
	var typeMap *TemplateTypeMap = this_.self_.GetTemplateTypeMap()
	_ = typeMap
	if typeMap.IsEmpty() || !(propertyType.GetJSTypeBase().HasAnyTemplateTypes()) {
		return propertyType
	}
	return this_.registry.ReplaceTemplateTypesPartial(propertyType, typeMap)
}

// Java method: JSType.findPropertyType(java.lang.String propertyName)
func (this_ *JSTypeBase) FindPropertyType_1_string(propertyName string) JSType {
	_ = propertyName
	j2gort.CheckNotNullString(propertyName)
	var propertyType JSType = this_.self_.FindPropertyTypeWithoutConsideringTemplateTypes_1_string(propertyName)
	_ = propertyType
	if propertyType == nil || propertyType.IsJavaNil_() {
		return nil
	}
	var typeMap *TemplateTypeMap = this_.self_.GetTemplateTypeMap()
	_ = typeMap
	if typeMap.IsEmpty() || !(propertyType.GetJSTypeBase().HasAnyTemplateTypes()) {
		return propertyType
	}
	return this_.registry.ReplaceTemplateTypesPartial(propertyType, typeMap)
}

// Java method: JSType.findPropertyTypeWithoutConsideringTemplateTypes(java.lang.String propertyName)
func (this_ *JSTypeBase) FindPropertyTypeWithoutConsideringTemplateTypes_1_string(propertyName string) JSType {
	_ = propertyName
	var autoboxObjType ObjectType = ObjectType_Cast(this_.self_.AutoboxesTo())
	_ = autoboxObjType
	if autoboxObjType != nil && !autoboxObjType.IsJavaNil_() {
		return autoboxObjType.GetObjectTypeBase().FindPropertyType_1_string(propertyName)
	}
	return nil
}

// Java method: JSType.findPropertyTypeWithoutConsideringTemplateTypes(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *JSTypeBase) FindPropertyTypeWithoutConsideringTemplateTypes_1_Property_Key(propertyName Property_Key) JSType {
	_ = propertyName
	var autoboxObjType ObjectType = ObjectType_Cast(this_.self_.AutoboxesTo())
	_ = autoboxObjType
	if autoboxObjType != nil && !autoboxObjType.IsJavaNil_() {
		return autoboxObjType.GetObjectTypeBase().FindPropertyType_1_Property_Key(propertyName)
	}
	return nil
}

// Java method: JSType.canBeCalled()
func (this_ *JSTypeBase) CanBeCalled() bool {
	return false
}

// Java method: JSType.canCastTo(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) CanCastTo(that JSType) bool {
	_ = that
	return j2gort.UnboxBoolean(j2gort.CastPtr[j2gort.BoxBoolean](this_.self_.Visit_2(ToIface_RelationshipVisitor(StaticGet_JSType_CAN_CAST_TO_VISITOR()), that)))
}

// Java method: JSType.autoboxesTo()
func (this_ *JSTypeBase) AutoboxesTo() JSType {
	return nil
}

// Java method: JSType.isBoxableScalar()
func (this_ *JSTypeBase) IsBoxableScalar() bool {
	return !(j2gort.IfaceEquals[JSType](this_.self_.AutoboxesTo(), nil))
}

// Java method: JSType.toObjectType()
func (this_ *JSTypeBase) ToObjectType() ObjectType {
	var objectType ObjectType
	_ = objectType
	var ok_13 bool
	objectType, ok_13 = j2gort.TryCast[ObjectType](this_.self_)
	if ok_13 {
		return objectType
	} else {
		return nil
	}
}

// Java method: JSType.autobox()
func (this_ *JSTypeBase) Autobox() JSType {
	var restricted JSType = this_.self_.RestrictByNotNullOrUndefined()
	_ = restricted
	var autobox JSType = restricted.AutoboxesTo()
	_ = autobox
	if autobox == nil || autobox.IsJavaNil_() {
		return restricted
	} else {
		return autobox
	}
}

// Java method: JSType.dereference()
func (this_ *JSTypeBase) Dereference() ObjectType {
	return this_.self_.Autobox().GetJSTypeBase().ToObjectType()
}

// Java method: JSType.canTestForEqualityWith(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) CanTestForEqualityWith(that JSType) bool {
	_ = that
	{

		re_14 := this_.self_.TestForEquality(that)
		ae_15 := pkg_base.StaticGet_Tri_UNKNOWN()
		if re_14 == 0 {
			j2gort.ThrowNPE()
		}
		return re_14 == ae_15
	}

}

// Java method: JSType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	return this_.TestForEqualityHelper(this_.self_, that)
}

// Java method: JSType.testForEqualityHelper(com.google.javascript.rhino.jstype.JSType aType, com.google.javascript.rhino.jstype.JSType bType)
func (this_ *JSTypeBase) TestForEqualityHelper(aType JSType, bType JSType) pkg_base.Tri {
	_ = aType
	_ = bType
	base_bType_ := bType.GetJSTypeBase()
	if ((((base_bType_.IsAllType() || base_bType_.IsUnknownType()) || base_bType_.IsNoResolvedType()) || aType.GetJSTypeBase().IsAllType()) || aType.GetJSTypeBase().IsUnknownType()) || aType.GetJSTypeBase().IsNoResolvedType() {
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	base_aType_ := aType.GetJSTypeBase()
	var aIsEmpty bool = base_aType_.IsEmptyType()
	_ = aIsEmpty
	var bIsEmpty bool = base_bType_.IsEmptyType()
	_ = bIsEmpty
	if aIsEmpty || bIsEmpty {
		if aIsEmpty && bIsEmpty {
			return pkg_base.StaticGet_Tri_TRUE()
		} else {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
	}
	if base_bType_.IsUnionType() {
		return bType.TestForEquality(aType)
	}
	if base_aType_.IsFunctionType() || base_bType_.IsFunctionType() {
		var otherType JSType
		_ = otherType
		if base_aType_.IsFunctionType() {
			otherType = bType
		} else {
			otherType = aType
		}
		if otherType.GetJSTypeBase().IsSymbol() {
			return pkg_base.StaticGet_Tri_FALSE()
		}
		var greatestSubtype JSType = otherType.GetJSTypeBase().GetGreatestSubtype(this_.GetNativeType(StaticGet_JSTypeNative_OBJECT_TYPE()))
		_ = greatestSubtype
		if greatestSubtype.GetJSTypeBase().IsNoType() || greatestSubtype.GetJSTypeBase().IsNoObjectType() {
			return pkg_base.StaticGet_Tri_FALSE()
		} else {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
	}
	if base_bType_.IsEnumElementType() {
		return bType.TestForEquality(aType)
	}
	if base_aType_.IsSymbol() {
		if base_bType_.CanCastTo(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_TYPE())) || base_bType_.CanCastTo(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_OBJECT_TYPE())) {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		} else {
			return pkg_base.StaticGet_Tri_FALSE()
		}
	}
	if base_bType_.IsSymbol() {
		if base_aType_.CanCastTo(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_TYPE())) || base_aType_.CanCastTo(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_OBJECT_TYPE())) {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		} else {
			return pkg_base.StaticGet_Tri_FALSE()
		}
	}
	return pkg_base.Tri(0)
}

// Java method: JSType.canTestForShallowEqualityWith(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) CanTestForShallowEqualityWith(that JSType) bool {
	_ = that
	if this_.IsNoResolvedType() || that.GetJSTypeBase().IsNoResolvedType() {
		return true
	}
	base_that_ := that.GetJSTypeBase()
	if this_.IsEmptyType() || base_that_.IsEmptyType() {
		return (this_.IsSubtypeOf_1(that) || base_that_.IsSubtypeOf_1(this_.self_))
	}
	var inf JSType = this_.GetGreatestSubtype(that)
	_ = inf
	return (!(inf.GetJSTypeBase().IsEmptyType()) || j2gort.IfaceEquals[JSType](inf, this_.registry.GetNativeType(StaticGet_JSTypeNative_LEAST_FUNCTION_TYPE())))
}

// Java method: JSType.isNullable()
func (this_ *JSTypeBase) IsNullable() bool {
	return false
}

// Java method: JSType.isVoidable()
func (this_ *JSTypeBase) IsVoidable() bool {
	return false
}

// Java method: JSType.isExplicitlyVoidable()
func (this_ *JSTypeBase) IsExplicitlyVoidable() bool {
	return false
}

// Java method: JSType.collapseUnion()
func (this_ *JSTypeBase) CollapseUnion() JSType {
	return this_.self_
}

// Java method: JSType.getLeastSupertype(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) GetLeastSupertype(that JSType) JSType {
	_ = that
	if this_.self_ == that || ((this_.self_ == nil || this_.self_.IsJavaNil_()) && (that == nil || that.IsJavaNil_())) {
		return this_.self_
	}
	if that.GetJSTypeBase().IsUnionType() {
		return that.ToMaybeUnionType().GetLeastSupertype(this_.self_)
	}
	return JSType_GetLeastSupertype(this_.self_, that)
}

// Java static method: JSType.getLeastSupertype(com.google.javascript.rhino.jstype.JSType thisType, com.google.javascript.rhino.jstype.JSType thatType)
func JSType_GetLeastSupertype(thisType JSType, thatType JSType) JSType {
	clinit_JSType()
	_ = thisType
	_ = thatType
	base_thisType_ := thisType.GetJSTypeBase()
	var areEquivalent bool = base_thisType_.EqualsJSType(thatType)
	_ = areEquivalent
	if areEquivalent {
		return thisType
	} else {
		return base_thisType_.registry.CreateUnionType_2(thisType, thatType)
	}
}

// Java method: JSType.getGreatestSubtype(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) GetGreatestSubtype(that JSType) JSType {
	_ = that
	return JSType_GetGreatestSubtype(this_.self_, that)
}

// Java static method: JSType.getGreatestSubtype(com.google.javascript.rhino.jstype.JSType thisType, com.google.javascript.rhino.jstype.JSType thatType)
func JSType_GetGreatestSubtype(thisType JSType, thatType JSType) JSType {
	clinit_JSType()
	_ = thisType
	_ = thatType
	base_thisType_ := thisType.GetJSTypeBase()
	if base_thisType_.IsFunctionType() && thatType.GetJSTypeBase().IsFunctionType() {
		return ToIface_JSType(thisType.ToMaybeFunctionType().SupAndInfHelper(thatType.ToMaybeFunctionType(), false))
	} else if base_thisType_.EqualsJSType(thatType) {
		return thisType
	} else if base_thisType_.IsUnknownType() {
		if !(thatType.GetJSTypeBase().IsUnknownType()) {
			return thatType
		} else if base_thisType_.IsCheckedUnknownType() || thatType.GetJSTypeBase().IsCheckedUnknownType() {
			return base_thisType_.GetNativeType(StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE())
		} else {
			return base_thisType_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}
	} else if thatType.GetJSTypeBase().IsUnknownType() {
		return thisType
	} else if base_thisType_.IsUnionType() {
		return UnionType_GetGreatestSubtype(thisType.ToMaybeUnionType(), thatType)
	} else if thatType.GetJSTypeBase().IsUnionType() {
		return UnionType_GetGreatestSubtype(thatType.ToMaybeUnionType(), thisType)
	} else if base_thisType_.IsTemplatizedType() {
		return thisType.ToMaybeTemplatizedType().GetGreatestSubtypeHelper(thatType)
	} else if thatType.GetJSTypeBase().IsTemplatizedType() {
		return thatType.ToMaybeTemplatizedType().GetGreatestSubtypeHelper(thisType)
	} else if base_thisType_.IsSubtypeOf_1(thatType) {
		return thisType
	} else if thatType.GetJSTypeBase().IsSubtypeOf_1(thisType) {
		return thatType
	} else if base_thisType_.IsRecordType() {
		return thisType.ToMaybeRecordType().GetGreatestSubtypeHelper(thatType)
	} else if thatType.GetJSTypeBase().IsRecordType() {
		return thatType.ToMaybeRecordType().GetGreatestSubtypeHelper(thisType)
	}
	if base_thisType_.IsEnumElementType() {
		var inf JSType = EnumElementType_GetGreatestSubtype(thisType.ToMaybeEnumElementType(), thatType)
		_ = inf
		if inf != nil && !inf.IsJavaNil_() {
			return inf
		}
	} else if thatType.GetJSTypeBase().IsEnumElementType() {
		var inf_2 JSType = EnumElementType_GetGreatestSubtype(thatType.ToMaybeEnumElementType(), thisType)
		_ = inf_2
		if inf_2 != nil && !inf_2.IsJavaNil_() {
			return inf_2
		}
	}
	if thisType.IsObject() && thatType.IsObject() {
		return base_thisType_.GetNativeType(StaticGet_JSTypeNative_NO_OBJECT_TYPE())
	}
	return base_thisType_.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
}

// Java method: JSType.getRestrictedTypeGivenOutcome(com.google.javascript.rhino.Outcome outcome)
func (this_ *JSTypeBase) GetRestrictedTypeGivenOutcome(outcome *Outcome) JSType {
	_ = outcome
	if outcome.IsTruthy() && j2gort.IfaceEquals[JSType](this_.self_, this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())) {
		return this_.GetNativeType(StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE())
	}
	if outcome.IsNullish().ToBoolean(false) {
		if j2gort.IfaceEquals[JSType](this_.self_, this_.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE())) || j2gort.IfaceEquals[JSType](this_.self_, this_.GetNativeType(StaticGet_JSTypeNative_NULL_TYPE())) {
			return this_.self_
		} else {
			return this_.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
		}
	}
	var literals BooleanLiteralSet = this_.self_.GetPossibleToBooleanOutcomes()
	_ = literals
	if literals.Contains(outcome.IsTruthy()) {
		return this_.self_
	} else {
		return this_.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
	}
}

// Java method: JSType.getPossibleToBooleanOutcomes()
func (this_ *JSTypeBase) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	panic("abstract or native method invoked")
}

// Java method: JSType.getTypesUnderEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) GetTypesUnderEquality(that JSType) JSType_TypePair {
	_ = that
	if that.GetJSTypeBase().IsUnionType() {
		var p JSType_TypePair = that.ToMaybeUnionType().GetTypesUnderEquality(this_.self_)
		_ = p
		return NewJSType_TypePair(p.typeB, p.typeA)
	}
	switch (this_.self_.TestForEquality(that)).Ordinal() {
	case int32(0) /* Tri.FALSE */ :
		return NewJSType_TypePair(nil, nil)
	case int32(2) /* Tri.TRUE */, int32(1) /* Tri.UNKNOWN */ :
		return NewJSType_TypePair(this_.self_, that)
	}
	panic("unreachable switch expression")
}

// Java method: JSType.getTypesUnderInequality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) GetTypesUnderInequality(that JSType) JSType_TypePair {
	_ = that
	base_that_ := that.GetJSTypeBase()
	if base_that_.IsUnionType() {
		var p JSType_TypePair = that.ToMaybeUnionType().GetTypesUnderInequality(this_.self_)
		_ = p
		return NewJSType_TypePair(p.typeB, p.typeA)
	}
	if this_.priv_JSType_isNullTypeOrVoidType_9170() && base_that_.priv_JSType_isNullTypeOrVoidType_9170() {
		var noType JSType = this_.registry.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
		_ = noType
		return NewJSType_TypePair(noType, noType)
	}
	return NewJSType_TypePair(this_.self_, that)
}

// Java method: JSType.getTypesUnderShallowEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) GetTypesUnderShallowEquality(that JSType) JSType_TypePair {
	_ = that
	var commonType JSType = this_.GetGreatestSubtype(that)
	_ = commonType
	return NewJSType_TypePair(commonType, commonType)
}

// Java method: JSType.getTypesUnderShallowInequality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) GetTypesUnderShallowInequality(that JSType) JSType_TypePair {
	_ = that
	if that.GetJSTypeBase().IsUnionType() {
		var p JSType_TypePair = that.ToMaybeUnionType().GetTypesUnderShallowInequality(this_.self_)
		_ = p
		return NewJSType_TypePair(p.typeB, p.typeA)
	}
	if (this_.self_ == that || ((this_.self_ == nil || this_.self_.IsJavaNil_()) && (that == nil || that.IsJavaNil_()))) && this_.priv_JSType_isNullTypeOrVoidType_9170() {
		return NewJSType_TypePair(nil, nil)
	}
	return NewJSType_TypePair(this_.self_, that)
}

// Java method: JSType.isNullTypeOrVoidType()
func (this_ *JSTypeBase) priv_JSType_isNullTypeOrVoidType_9170() bool {
	return ((this_.typeKindMask & (int32(131072) | int32(262144))) != int32(0))
}

// Java method: JSType.getUnionMembers()
func (this_ *JSTypeBase) GetUnionMembers() j2gort.Ext_java_lang_Iterable {
	if this_.IsUnionType() {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.self_.ToMaybeUnionType().GetAlternates())
	} else {
		return nil
	}
}

// Java method: JSType.restrictByNotNullOrUndefined()
func (this_ *JSTypeBase) RestrictByNotNullOrUndefined() JSType {
	return this_.self_
}

// Java method: JSType.restrictByNotUndefined()
func (this_ *JSTypeBase) RestrictByNotUndefined() JSType {
	return this_.self_
}

// Java method: JSType.restrictByNotNull()
func (this_ *JSTypeBase) RestrictByNotNull() JSType {
	return this_.self_
}

// Java method: JSType.isSubtypeWithoutStructuralTyping(com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *JSTypeBase) IsSubtypeWithoutStructuralTyping(supertype JSType) bool {
	_ = supertype
	return SubtypeChecker_Check(this_.registry, this_.self_, supertype, false, StaticGet_JSType_SubtypingMode_NORMAL())
}

// Java method: JSType.isSubtype(com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *JSTypeBase) IsSubtype_1(supertype JSType) bool {
	_ = supertype
	return this_.IsSubtypeOf_1(supertype)
}

// Java method: JSType.isSubtype(com.google.javascript.rhino.jstype.JSType supertype, com.google.javascript.rhino.jstype.JSType.SubtypingMode mode)
func (this_ *JSTypeBase) IsSubtype_2(supertype JSType, mode JSType_SubtypingMode) bool {
	_ = supertype
	_ = mode
	return this_.IsSubtypeOf_2(supertype, mode)
}

// Java method: JSType.isSubtypeOf(com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *JSTypeBase) IsSubtypeOf_1(supertype JSType) bool {
	_ = supertype
	return SubtypeChecker_Check(this_.registry, this_.self_, supertype, true, StaticGet_JSType_SubtypingMode_NORMAL())
}

// Java method: JSType.isSubtypeOf(com.google.javascript.rhino.jstype.JSType supertype, com.google.javascript.rhino.jstype.JSType.SubtypingMode mode)
func (this_ *JSTypeBase) IsSubtypeOf_2(supertype JSType, mode JSType_SubtypingMode) bool {
	_ = supertype
	_ = mode
	return SubtypeChecker_Check(this_.registry, this_.self_, supertype, true, mode)
}

// Java method: JSType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *JSTypeBase) Visit_1(visitor Visitor) any {
	_ = visitor
	panic("abstract or native method invoked")
}

// Java method: JSType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *JSTypeBase) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	panic("abstract or native method invoked")
}

// Java method: JSType.resolve(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *JSTypeBase) Resolve(reporter ErrorReporter) JSType {
	_ = reporter
	this_.registry.GetResolver().AssertLegalToResolveTypes()
	if !(this_.IsResolved()) {
		this_.self_.SetCachedResolveResult(this_.self_)
		var resolvedType JSType = this_.self_.ResolveInternal(reporter)
		_ = resolvedType
		this_.self_.SetCachedResolveResult(resolvedType)
		if !(this_.IsResolved()) {
			j2gort.CheckState(false)
		}
		return resolvedType
	}
	if (this_.typeKindMask & int32(5184)) == int32(0) {
		return this_.self_
	}
	return this_.self_.GetCachedResolveResult()
}

// Java method: JSType.getCachedResolveResult()
func (this_ *JSTypeBase) GetCachedResolveResult() JSType {
	if this_.IsResolved() {
		return this_.self_
	} else {
		return nil
	}
}

// Java method: JSType.setCachedResolveResult(com.google.javascript.rhino.jstype.JSType result)
func (this_ *JSTypeBase) SetCachedResolveResult(result JSType) {
	_ = result
	j2gort.CheckNotNullIface(result)
	if !(result == this_.self_ || ((result == nil || result.IsJavaNil_()) && (this_.self_ == nil || this_.self_.IsJavaNil_()))) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected non-self resolve result on "+j2gort.StrVal(j2gort.GetObjectClass(this_.self_)))))
	}
	this_.jstypeFlags = int8((int32(this_.jstypeFlags) | int32(int8(8))))
}

// Java method: JSType.clearCachedResolveResult()
func (this_ *JSTypeBase) ClearCachedResolveResult() {
	this_.jstypeFlags = int8((int32(this_.jstypeFlags) & ^(int32(int8(8)))))
	this_.cachedHashCode = int32(0)
}

// Java method: JSType.markResolved()
func (this_ *JSTypeBase) MarkResolved() {
	this_.jstypeFlags = int8((int32(this_.jstypeFlags) | int32(int8(8))))
}

// Java method: JSType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *JSTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	panic("abstract or native method invoked")
}

// Java method: JSType.eagerlyResolveToSelf()
func (this_ *JSTypeBase) EagerlyResolveToSelf() {
	if !(!(this_.IsResolved())) {
		j2gort.CheckStateMsg(false, this_.self_)
	}
	this_.self_.SetCachedResolveResult(this_.self_)
	this_.registry.GetResolver().ResolveIfClosed(this_.self_, this_.FastTypeClass())
}

// Java method: JSType.isResolved()
func (this_ *JSTypeBase) IsResolved() bool {
	return ((int32(this_.jstypeFlags) & int32(int8(8))) != int32(0))
}

// Java method: JSType.isSuccessfullyResolved()
func (this_ *JSTypeBase) IsSuccessfullyResolved() bool {
	return (this_.IsResolved() && !(this_.IsNoResolvedType()))
}

// Java method: JSType.isUnsuccessfullyResolved()
func (this_ *JSTypeBase) IsUnsuccessfullyResolved() bool {
	return (this_.IsResolved() && this_.IsNoResolvedType())
}

// Java static method: JSType.safeResolve(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.ErrorReporter reporter)
func JSType_SafeResolve(type_ JSType, reporter ErrorReporter) JSType {
	clinit_JSType()
	_ = type_
	_ = reporter
	if type_ == nil || type_.IsJavaNil_() {
		return nil
	} else {
		return type_.GetJSTypeBase().Resolve(reporter)
	}
}

// Java method: JSType.setValidator(com.google.common.base.Predicate<com.google.javascript.rhino.jstype.JSType> validator)
func (this_ *JSTypeBase) SetValidator(validator j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = validator
	return validator.Apply(this_.self_)
}

// Java method: JSType.toString()
func (this_ *JSTypeBase) String() string {
	return NewTypeStringBuilder(false).Append_1_JSType(this_.self_).Build()
}

// Java method: JSType.toAnnotationString(com.google.javascript.rhino.jstype.JSType.Nullability nullability)
func (this_ *JSTypeBase) ToAnnotationString(nullability JSType_Nullability) string {
	_ = nullability
	var val_16 TypeStringBuilder
	InitTypeStringBuilder(&val_16, true)
	var builder *TypeStringBuilder = &val_16
	_ = builder
	return (func() *TypeStringBuilder {
		if nullability == StaticGet_JSType_Nullability_EXPLICIT() {
			return builder.AppendNonNull(this_.self_)
		}
		return builder.Append_1_JSType(this_.self_)
	}()).Build()
}

// Java method: JSType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *JSTypeBase) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	panic("abstract or native method invoked")
}

// Java method: JSType.matchConstraint(com.google.javascript.rhino.jstype.JSType constraint)
func (this_ *JSTypeBase) MatchConstraint(constraint JSType) {
	_ = constraint
}

// Java method: JSType.toMaybeObjectType()
func (this_ *JSTypeBase) ToMaybeObjectType() ObjectType {
	return this_.ToObjectType()
}

// Java enum: com.google.javascript.rhino.jstype.JSType.HasPropertyKind
type JSType_HasPropertyKind uint8

const (
	enum_JSType_HasPropertyKind_ABSENT        JSType_HasPropertyKind = 1
	enum_JSType_HasPropertyKind_KNOWN_PRESENT JSType_HasPropertyKind = 2
	enum_JSType_HasPropertyKind_MAYBE_PRESENT JSType_HasPropertyKind = 3
)

var enumBaseTable_JSType_HasPropertyKind = [...]j2gort.EnumBase{
	{Name_: "ABSENT", Ordinal_: 0},
	{Name_: "KNOWN_PRESENT", Ordinal_: 1},
	{Name_: "MAYBE_PRESENT", Ordinal_: 2},
}
var enumValues_JSType_HasPropertyKind = []JSType_HasPropertyKind{1, 2, 3}

func StaticGet_JSType_HasPropertyKind_ABSENT() JSType_HasPropertyKind {
	return enum_JSType_HasPropertyKind_ABSENT
}
func StaticGet_JSType_HasPropertyKind_KNOWN_PRESENT() JSType_HasPropertyKind {
	return enum_JSType_HasPropertyKind_KNOWN_PRESENT
}
func StaticGet_JSType_HasPropertyKind_MAYBE_PRESENT() JSType_HasPropertyKind {
	return enum_JSType_HasPropertyKind_MAYBE_PRESENT
}

func JSType_HasPropertyKind_Values() []JSType_HasPropertyKind {
	return j2gort.CloneSlice(enumValues_JSType_HasPropertyKind)
}

func JSType_HasPropertyKind_ValueOf(name string) JSType_HasPropertyKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSType_HasPropertyKind {
		if enumBaseTable_JSType_HasPropertyKind[i_].Name_ == name {
			return JSType_HasPropertyKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.JSType.HasPropertyKind." + name))
}

// Java static method: HasPropertyKind.of(boolean has)
func JSType_HasPropertyKind_Of(has bool) JSType_HasPropertyKind {
	_ = has
	if has {
		return StaticGet_JSType_HasPropertyKind_KNOWN_PRESENT()
	} else {
		return StaticGet_JSType_HasPropertyKind_ABSENT()
	}
}

func (this_ JSType_HasPropertyKind) IsJavaNil_() bool                              { return this_ == 0 }
func (JSType_HasPropertyKind) HasIdentityEquals_()                                 {}
func (JSType_HasPropertyKind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JSType_HasPropertyKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSType_HasPropertyKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JSType_HasPropertyKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSType$HasPropertyKind", "HasPropertyKind"
}
func (this_ JSType_HasPropertyKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSType_HasPropertyKind[this_-1]
}
func (this_ JSType_HasPropertyKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSType_HasPropertyKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSType_HasPropertyKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSType_HasPropertyKind[this_-1].Name_
}
func (this_ JSType_HasPropertyKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSType_HasPropertyKind[this_-1].Name_
}
func (this_ JSType_HasPropertyKind) Equals(other_ any) bool {
	o_, ok_ := other_.(JSType_HasPropertyKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSType_HasPropertyKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSType_HasPropertyKind[this_-1].HashCode()
}
func (this_ JSType_HasPropertyKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSType_HasPropertyKind[this_-1].CompareTo(other_)
}
func (this_ JSType_HasPropertyKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSType_HasPropertyKind[this_-1].Name_)
}
func init() {
	vals_ := []any{JSType_HasPropertyKind(1), JSType_HasPropertyKind(2), JSType_HasPropertyKind(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSType.HasPropertyKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSType$HasPropertyKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSType.HasPropertyKind", "com.google.javascript.rhino.jstype.JSType$HasPropertyKind", func(v_ any) bool {
		return j2gort.InstanceOf[JSType_HasPropertyKind](v_)
	})
}

// Java enum: com.google.javascript.rhino.jstype.JSType.SubtypingMode
type JSType_SubtypingMode uint8

const (
	enum_JSType_SubtypingMode_NORMAL                JSType_SubtypingMode = 1
	enum_JSType_SubtypingMode_IGNORE_NULL_UNDEFINED JSType_SubtypingMode = 2
)

var enumBaseTable_JSType_SubtypingMode = [...]j2gort.EnumBase{
	{Name_: "NORMAL", Ordinal_: 0},
	{Name_: "IGNORE_NULL_UNDEFINED", Ordinal_: 1},
}
var enumValues_JSType_SubtypingMode = []JSType_SubtypingMode{1, 2}

func StaticGet_JSType_SubtypingMode_NORMAL() JSType_SubtypingMode {
	return enum_JSType_SubtypingMode_NORMAL
}
func StaticGet_JSType_SubtypingMode_IGNORE_NULL_UNDEFINED() JSType_SubtypingMode {
	return enum_JSType_SubtypingMode_IGNORE_NULL_UNDEFINED
}

func JSType_SubtypingMode_Values() []JSType_SubtypingMode {
	return j2gort.CloneSlice(enumValues_JSType_SubtypingMode)
}

func JSType_SubtypingMode_ValueOf(name string) JSType_SubtypingMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSType_SubtypingMode {
		if enumBaseTable_JSType_SubtypingMode[i_].Name_ == name {
			return JSType_SubtypingMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.JSType.SubtypingMode." + name))
}

func (this_ JSType_SubtypingMode) IsJavaNil_() bool                              { return this_ == 0 }
func (JSType_SubtypingMode) HasIdentityEquals_()                                 {}
func (JSType_SubtypingMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JSType_SubtypingMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSType_SubtypingMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JSType_SubtypingMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSType$SubtypingMode", "SubtypingMode"
}
func (this_ JSType_SubtypingMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSType_SubtypingMode[this_-1]
}
func (this_ JSType_SubtypingMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSType_SubtypingMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSType_SubtypingMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSType_SubtypingMode[this_-1].Name_
}
func (this_ JSType_SubtypingMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSType_SubtypingMode[this_-1].Name_
}
func (this_ JSType_SubtypingMode) Equals(other_ any) bool {
	o_, ok_ := other_.(JSType_SubtypingMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSType_SubtypingMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSType_SubtypingMode[this_-1].HashCode()
}
func (this_ JSType_SubtypingMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSType_SubtypingMode[this_-1].CompareTo(other_)
}
func (this_ JSType_SubtypingMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSType_SubtypingMode[this_-1].Name_)
}
func init() {
	vals_ := []any{JSType_SubtypingMode(1), JSType_SubtypingMode(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSType.SubtypingMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSType$SubtypingMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSType.SubtypingMode", "com.google.javascript.rhino.jstype.JSType$SubtypingMode", func(v_ any) bool {
		return j2gort.InstanceOf[JSType_SubtypingMode](v_)
	})
}

// Java class: com.google.javascript.rhino.jstype.JSType.TypePair
type JSType_TypePair struct {
	typeA  JSType `gson:"typeA"`
	typeB  JSType `gson:"typeB"`
	valid_ bool   `gson:"-"`
}

func (this_ JSType_TypePair) BindSelf_(s any) {
	_ = s
}
func (this_ JSType_TypePair) IsJavaNil_() bool { return !this_.valid_ }
func (JSType_TypePair) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSType$TypePair", "TypePair"
}
func (JSType_TypePair) HasIdentityEquals_() {}
func (this_ JSType_TypePair) IdentityHash_() int32 {
	if !this_.valid_ {
		return 0
	}
	var h_ int32 = 1
	{
		hi_17 := this_.typeA
		h_ = 31 * h_
		if hi_17 != nil && !hi_17.IsJavaNil_() {
			h_ += hi_17.GetJSTypeBase().HashCode()
		}
	}
	{
		hi_18 := this_.typeB
		h_ = 31 * h_
		if hi_18 != nil && !hi_18.IsJavaNil_() {
			h_ += hi_18.GetJSTypeBase().HashCode()
		}
	}
	return h_
}
func (this_ JSType_TypePair) FastValHash_() (int32, bool, bool) {
	if this_.IsJavaNil_() {
		return 0, false, true
	}
	return this_.IdentityHash_(), true, true
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSType.TypePair", "com.google.javascript.rhino.jstype.JSType$TypePair", func(v_ any) bool {
		return j2gort.InstanceOf[JSType_TypePair](v_)
	})
}
func (this_ JSType_TypePair) FieldGet_typeA() JSType {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return this_.typeA
}
func (this_ JSType_TypePair) FieldGet_typeB() JSType {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return this_.typeB
}

func (this_ *JSType_TypePair) initFields_JSType_TypePair() {
	_ = this_
	this_.valid_ = true
}

// Java constructor: TypePair(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType typeA, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType typeB)
func NewJSType_TypePair(typeA JSType, typeB JSType) JSType_TypePair {
	this_ := JSType_TypePair{}
	InitJSType_TypePair(&this_, typeA, typeB)
	return this_
}

func InitJSType_TypePair(this_ *JSType_TypePair, typeA JSType, typeB JSType) {
	_ = typeA
	_ = typeB
	this_.initFields_JSType_TypePair()
	this_.typeA = typeA
	this_.typeB = typeB
}

// Java interface: com.google.javascript.rhino.jstype.JSType.WithSourceRef
type JSType_WithSourceRef interface {
	Implements_rhino_JSType_WithSourceRef()
	IsJavaNil_() bool
	GetSource() Node
	GetGoogModuleId() string
}

func ToIface_JSType_WithSourceRef[T interface {
	JSType_WithSourceRef
	comparable
}](v T) JSType_WithSourceRef {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JSType_WithSourceRef[T JSType_WithSourceRef](v T) JSType_WithSourceRef {
	var vi JSType_WithSourceRef = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSType.WithSourceRef", "com.google.javascript.rhino.jstype.JSType$WithSourceRef", func(v_ any) bool {
		return j2gort.InstanceOf[JSType_WithSourceRef](v_)
	})
}

// Java enum: com.google.javascript.rhino.jstype.JSType.MatchStatus
type JSType_MatchStatus uint8

const (
	enum_JSType_MatchStatus_MATCH      JSType_MatchStatus = 1
	enum_JSType_MatchStatus_NOT_MATCH  JSType_MatchStatus = 2
	enum_JSType_MatchStatus_PROCESSING JSType_MatchStatus = 3
)

var enumBaseTable_JSType_MatchStatus = [...]j2gort.EnumBase{
	{Name_: "MATCH", Ordinal_: 0},
	{Name_: "NOT_MATCH", Ordinal_: 1},
	{Name_: "PROCESSING", Ordinal_: 2},
}
var enumValues_JSType_MatchStatus = []JSType_MatchStatus{1, 2, 3}

var enumField_JSType_MatchStatus_isSubtype = [...]bool{
	true,
	false,
	true,
}

func (this_ JSType_MatchStatus) FieldGet_isSubtype() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_JSType_MatchStatus_isSubtype[this_-1]
}

func StaticGet_JSType_MatchStatus_MATCH() JSType_MatchStatus {
	return enum_JSType_MatchStatus_MATCH
}
func StaticGet_JSType_MatchStatus_NOT_MATCH() JSType_MatchStatus {
	return enum_JSType_MatchStatus_NOT_MATCH
}
func StaticGet_JSType_MatchStatus_PROCESSING() JSType_MatchStatus {
	return enum_JSType_MatchStatus_PROCESSING
}

func JSType_MatchStatus_Values() []JSType_MatchStatus {
	return j2gort.CloneSlice(enumValues_JSType_MatchStatus)
}

func JSType_MatchStatus_ValueOf(name string) JSType_MatchStatus {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSType_MatchStatus {
		if enumBaseTable_JSType_MatchStatus[i_].Name_ == name {
			return JSType_MatchStatus(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.JSType.MatchStatus." + name))
}

// Java method: MatchStatus.subtypeValue()
func (this_ JSType_MatchStatus) SubtypeValue() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_isSubtype()
}

// Java static method: MatchStatus.valueOf(boolean match)
func JSType_MatchStatus_ValueOf_1_bool(match bool) JSType_MatchStatus {
	_ = match
	if match {
		return StaticGet_JSType_MatchStatus_MATCH()
	} else {
		return StaticGet_JSType_MatchStatus_NOT_MATCH()
	}
}

func (this_ JSType_MatchStatus) IsJavaNil_() bool                              { return this_ == 0 }
func (JSType_MatchStatus) HasIdentityEquals_()                                 {}
func (JSType_MatchStatus) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JSType_MatchStatus) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSType_MatchStatus) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JSType_MatchStatus) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSType$MatchStatus", "MatchStatus"
}
func (this_ JSType_MatchStatus) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSType_MatchStatus[this_-1]
}
func (this_ JSType_MatchStatus) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSType_MatchStatus) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSType_MatchStatus) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSType_MatchStatus[this_-1].Name_
}
func (this_ JSType_MatchStatus) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSType_MatchStatus[this_-1].Name_
}
func (this_ JSType_MatchStatus) Equals(other_ any) bool {
	o_, ok_ := other_.(JSType_MatchStatus)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSType_MatchStatus) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSType_MatchStatus[this_-1].HashCode()
}
func (this_ JSType_MatchStatus) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSType_MatchStatus[this_-1].CompareTo(other_)
}
func (this_ JSType_MatchStatus) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSType_MatchStatus[this_-1].Name_)
}
func init() {
	vals_ := []any{JSType_MatchStatus(1), JSType_MatchStatus(2), JSType_MatchStatus(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSType.MatchStatus", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSType$MatchStatus", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSType.MatchStatus", "com.google.javascript.rhino.jstype.JSType$MatchStatus", func(v_ any) bool {
		return j2gort.InstanceOf[JSType_MatchStatus](v_)
	})
}

// Java enum: com.google.javascript.rhino.jstype.JSType.Nullability
type JSType_Nullability uint8

const (
	enum_JSType_Nullability_EXPLICIT JSType_Nullability = 1
	enum_JSType_Nullability_IMPLICIT JSType_Nullability = 2
)

var enumBaseTable_JSType_Nullability = [...]j2gort.EnumBase{
	{Name_: "EXPLICIT", Ordinal_: 0},
	{Name_: "IMPLICIT", Ordinal_: 1},
}
var enumValues_JSType_Nullability = []JSType_Nullability{1, 2}

func StaticGet_JSType_Nullability_EXPLICIT() JSType_Nullability {
	return enum_JSType_Nullability_EXPLICIT
}
func StaticGet_JSType_Nullability_IMPLICIT() JSType_Nullability {
	return enum_JSType_Nullability_IMPLICIT
}

func JSType_Nullability_Values() []JSType_Nullability {
	return j2gort.CloneSlice(enumValues_JSType_Nullability)
}

func JSType_Nullability_ValueOf(name string) JSType_Nullability {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSType_Nullability {
		if enumBaseTable_JSType_Nullability[i_].Name_ == name {
			return JSType_Nullability(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.JSType.Nullability." + name))
}

func (this_ JSType_Nullability) IsJavaNil_() bool                              { return this_ == 0 }
func (JSType_Nullability) HasIdentityEquals_()                                 {}
func (JSType_Nullability) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JSType_Nullability) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSType_Nullability) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JSType_Nullability) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.JSType$Nullability", "Nullability"
}
func (this_ JSType_Nullability) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSType_Nullability[this_-1]
}
func (this_ JSType_Nullability) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSType_Nullability) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSType_Nullability) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSType_Nullability[this_-1].Name_
}
func (this_ JSType_Nullability) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSType_Nullability[this_-1].Name_
}
func (this_ JSType_Nullability) Equals(other_ any) bool {
	o_, ok_ := other_.(JSType_Nullability)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSType_Nullability) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSType_Nullability[this_-1].HashCode()
}
func (this_ JSType_Nullability) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSType_Nullability[this_-1].CompareTo(other_)
}
func (this_ JSType_Nullability) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSType_Nullability[this_-1].Name_)
}
func init() {
	vals_ := []any{JSType_Nullability(1), JSType_Nullability(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSType.Nullability", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSType$Nullability", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSType.Nullability", "com.google.javascript.rhino.jstype.JSType$Nullability", func(v_ any) bool {
		return j2gort.InstanceOf[JSType_Nullability](v_)
	})
}

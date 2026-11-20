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
// Port of com/google/javascript/rhino/jstype/FunctionType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.FunctionType
type FunctionType struct {
	PrototypeObjectTypeBase
	parameters []FunctionType_Parameter `gson:"parameters"`
	returnType JSType                   `gson:"returnType"`
	typeOfThis JSType                   `gson:"typeOfThis"`
}

func (this_ *FunctionType) BindSelf_(s any) {
	_ = s
	this_.PrototypeObjectTypeBase.BindSelf_(s)
}
func (this_ *FunctionType) IsJavaNil_() bool { return this_ == nil }
func (*FunctionType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.FunctionType", "FunctionType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.FunctionType", "com.google.javascript.rhino.jstype.FunctionType", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionType](v_)
	})
}
func (*FunctionType) Implements_rhino_PrototypeObjectType()  {}
func (*FunctionType) Implements_rhino_JSType_WithSourceRef() {}
func (*FunctionType) Implements_rhino_ObjectType()           {}
func (*FunctionType) Implements_rhino_JSType()               {}

var static_FunctionType_TYPE_CLASS JSTypeClass

func StaticGet_FunctionType_TYPE_CLASS() JSTypeClass {
	clinit_FunctionType()
	return static_FunctionType_TYPE_CLASS
}
func StaticSet_FunctionType_TYPE_CLASS(v JSTypeClass) {
	clinit_FunctionType()
	static_FunctionType_TYPE_CLASS = v
}

var static_FunctionType_EMPTY_PARAMETERS []FunctionType_Parameter

func StaticGet_FunctionType_EMPTY_PARAMETERS() []FunctionType_Parameter {
	clinit_FunctionType()
	return static_FunctionType_EMPTY_PARAMETERS
}
func StaticSet_FunctionType_EMPTY_PARAMETERS(v []FunctionType_Parameter) {
	clinit_FunctionType()
	static_FunctionType_EMPTY_PARAMETERS = v
}

var static_FunctionType_EMPTY_PARAMETER_LIST j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_FunctionType_EMPTY_PARAMETER_LIST() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_FunctionType()
	return static_FunctionType_EMPTY_PARAMETER_LIST
}
func StaticSet_FunctionType_EMPTY_PARAMETER_LIST(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_FunctionType()
	static_FunctionType_EMPTY_PARAMETER_LIST = v
}

var static_FunctionType_EMPTY_TEMPLATE_KEYS j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_FunctionType_EMPTY_TEMPLATE_KEYS() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_FunctionType()
	return static_FunctionType_EMPTY_TEMPLATE_KEYS
}
func StaticSet_FunctionType_EMPTY_TEMPLATE_KEYS(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_FunctionType()
	static_FunctionType_EMPTY_TEMPLATE_KEYS = v
}

var static_FunctionType_EMPTY_TEMPLATE_VALUES j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_FunctionType_EMPTY_TEMPLATE_VALUES() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_FunctionType()
	return static_FunctionType_EMPTY_TEMPLATE_VALUES
}
func StaticSet_FunctionType_EMPTY_TEMPLATE_VALUES(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_FunctionType()
	static_FunctionType_EMPTY_TEMPLATE_VALUES = v
}

var static_FunctionType_EMPTY_CTOR_KEYS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_FunctionType_EMPTY_CTOR_KEYS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_FunctionType()
	return static_FunctionType_EMPTY_CTOR_KEYS
}
func StaticSet_FunctionType_EMPTY_CTOR_KEYS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_FunctionType()
	static_FunctionType_EMPTY_CTOR_KEYS = v
}
func StaticGet_FunctionType_NO_KIND_NORMAL() int8             { return int8(0) }
func StaticGet_FunctionType_NO_KIND_NO_OBJECT() int8          { return int8(1) }
func StaticGet_FunctionType_NO_KIND_NO_TYPE() int8            { return int8(2) }
func StaticGet_FunctionType_FN_KIND_ORDINARY() int32          { return int32(0) }
func StaticGet_FunctionType_FN_KIND_CONSTRUCTOR() int32       { return int32(1) }
func StaticGet_FunctionType_FN_KIND_INTERFACE() int32         { return int32(2) }
func StaticGet_FunctionType_FN_KIND_NONE() int32              { return int32(3) }
func StaticGet_FunctionType_FN_NO_KIND_SHIFT() int32          { return int32(2) }
func StaticGet_FunctionType_FN_NO_KIND_MASK() int32           { return int32(12) }
func StaticGet_FunctionType_FN_NO_KIND_NO_OBJECT_BITS() int32 { return int32(4) }
func StaticGet_FunctionType_FN_NO_KIND_NO_TYPE_BITS() int32   { return int32(8) }
func StaticGet_FunctionType_FN_FLAG_RETURN_INFERRED() int8    { return int8(16) }
func StaticGet_FunctionType_FN_FLAG_STRUCTURAL_IFACE() int8   { return int8(32) }
func StaticGet_FunctionType_FN_FLAG_ABSTRACT() int8           { return int8(64) }
func StaticGet_FunctionType_AMBIGUITY_UNKNOWN() int8          { return int8(0) }
func StaticGet_FunctionType_AMBIGUITY_AMBIGUOUS() int8        { return int8(1) }
func StaticGet_FunctionType_AMBIGUITY_UNAMBIGUOUS() int8      { return int8(2) }
func StaticGet_FunctionType_PROP_ACCESS_UNSET() int8          { return int8(0) }
func StaticGet_FunctionType_PROP_ACCESS_ANY() int8            { return int8(1) }
func StaticGet_FunctionType_PROP_ACCESS_ANY_EXPLICIT() int8   { return int8(2) }
func StaticGet_FunctionType_PROP_ACCESS_STRUCT() int8         { return int8(3) }
func StaticGet_FunctionType_PROP_ACCESS_DICT() int8           { return int8(4) }

var static_FunctionType_DELEGATE_SUFFIX string

func StaticGet_FunctionType_DELEGATE_SUFFIX() string {
	clinit_FunctionType()
	return static_FunctionType_DELEGATE_SUFFIX
}
func StaticSet_FunctionType_DELEGATE_SUFFIX(v string) {
	clinit_FunctionType()
	static_FunctionType_DELEGATE_SUFFIX = v
}

var clinitDone_FunctionType uint32
var clinitState_FunctionType j2gort.ClinitState

func clinit_FunctionType() {
	if atomic.LoadUint32(&clinitDone_FunctionType) != 0 {
		return
	}
	clinitSlow_FunctionType()
}

func clinitSlow_FunctionType() {
	if !j2gort.EnterClinit(&clinitState_FunctionType, &clinitDone_FunctionType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FunctionType, &clinitDone_FunctionType)
	static_FunctionType_TYPE_CLASS = StaticGet_JSTypeClass_FUNCTION()
	static_FunctionType_EMPTY_PARAMETERS = make([]FunctionType_Parameter, int32(0))
	static_FunctionType_EMPTY_PARAMETER_LIST = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	static_FunctionType_EMPTY_TEMPLATE_KEYS = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	static_FunctionType_EMPTY_TEMPLATE_VALUES = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	static_FunctionType_EMPTY_CTOR_KEYS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	static_FunctionType_DELEGATE_SUFFIX = ObjectType_CreateDelegateSuffix("Proxy")
}

func (this_ *FunctionType) initFields_FunctionType() {
	_ = this_
}

// Java constructor: FunctionType(com.google.javascript.rhino.jstype.FunctionType.Builder builder)
func NewFunctionType(builder *FunctionType_Builder) *FunctionType {
	clinit_FunctionType()
	this_ := &FunctionType{}
	this_.BindSelf_(this_)
	InitFunctionType(this_, builder)
	return this_
}
func allocPooled_FunctionType(slab_ *j2gort.PooledSlab[FunctionType], builder *FunctionType_Builder) *FunctionType {
	clinit_FunctionType()
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitFunctionType(this_, builder)
	return this_
}

func InitFunctionType(this_ *FunctionType, builder *FunctionType_Builder) {
	_ = builder
	initPrototypeObjectTypeFromBase_(&this_.PrototypeObjectTypeBase, &builder.PrototypeObjectType_BuilderBase, FunctionType_ComputeFunctionKindMask(builder))
	this_.initFields_FunctionType()
	this_.SetPrettyPrint(true)
	var source Node = builder.sourceNode
	_ = source
	if !(((source == 0) || source.IsFunction()) || source.IsClass()) {
		j2gort.CheckArgument(false)
	}
	this_.GetJSTypeBase().subSourceNode = source
	var kind FunctionType_Kind = builder.kind
	_ = kind
	var flags int8 = int8((FunctionType_EncodeKind(kind) | (javalang.Shl32((int32(builder.noTypeKind) & int32(3)), int(int32(2))))))
	_ = flags
	if builder.returnTypeIsInferred {
		flags = int8((int32(flags) | int32(int8(16))))
	}
	if builder.isAbstract {
		flags = int8((int32(flags) | int32(int8(64))))
	}
	this_.GetJSTypeBase().fnFlags = flags
	if builder.googModuleId != j2gort.NullString {
		{
			r_1 := this_.EnsureCold()
			r_1.googModuleId = builder.googModuleId
		}
	}
	if builder.isKnownAmbiguous {
		{
			r_2 := this_.EnsureCold()
			r_2.constructorAmbiguity = int8(1)
		}
	}
	if builder.parameterArray != nil {
		this_.parameters = builder.parameterArray
	} else if builder.parameters != nil {
		this_.parameters = FunctionType_CopyListToParameterArray(builder.parameters)
	} else {
		this_.parameters = this_.GetJSTypeBase().registry.GetDefaultVarArgsParameterArray()
	}
	if !(builder.returnsOwnInstanceType) {
		if builder.returnType == nil || builder.returnType.IsJavaNil_() {
			this_.returnType = this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
		} else {
			this_.returnType = builder.returnType
		}
	}
	if builder.typeOfThis != nil && !builder.typeOfThis.IsJavaNil_() {
		this_.typeOfThis = builder.typeOfThis
	} else {
		switch (kind).Ordinal() {
		case int32(1) /* Kind.CONSTRUCTOR */, int32(2) /* Kind.INTERFACE */ :
			var typeOfThisBuilder *InstanceObjectType_Builder = InstanceObjectType_BuilderForCtor(this_)
			_ = typeOfThisBuilder
			var ctorKeys j2gort.Ext_com_google_common_collect_ImmutableSet = builder.constructorOnlyKeys
			_ = ctorKeys
			if !(ctorKeys.IsEmpty()) {
				j2gort.CastPtr[InstanceObjectType_Builder](typeOfThisBuilder.SetTemplateTypeMap(this_.GetJSTypeBase().templateTypeMap.CopyWithoutKeys(j2gort.AsIface[j2gort.Ext_java_util_Set](ctorKeys)))).SetTemplateParamCount((this_.GetTemplateParamCount() - ctorKeys.Size()))
			}
			this_.typeOfThis = ToIface_JSType(j2gort.CastPtr[InstanceObjectType](typeOfThisBuilder.Build()))
		case int32(0) /* Kind.ORDINARY */ :
			this_.typeOfThis = Upcast_JSType(this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_UNKNOWN_TYPE()))
		case int32(3) /* Kind.NONE */ :
			this_.typeOfThis = ToIface_JSType(this_)
		}
	}
	if builder.returnsOwnInstanceType {
		if this_.typeOfThis == nil || this_.typeOfThis.IsJavaNil_() {
			this_.returnType = this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
		} else {
			this_.returnType = this_.typeOfThis
		}
	}
	if builder.primitiveId != 0 {
		{
			r_3 := this_.EnsureCold()
			r_3.closurePrimitive = builder.primitiveId
		}
	}
	var canonicalRepresentation *FunctionType = builder.canonicalRepresentation
	_ = canonicalRepresentation
	if !((canonicalRepresentation == nil) || (kind == StaticGet_FunctionType_Kind_CONSTRUCTOR())) {
		j2gort.CheckArgumentStr(false, "Only constructors should have canonical representations")
	}
	if canonicalRepresentation != nil {
		{
			r_4 := this_.EnsureCold()
			r_4.canonicalRepresentation = canonicalRepresentation
		}
	}
	if builder.setPrototypeBasedOn != nil && !builder.setPrototypeBasedOn.IsJavaNil_() {
		this_.SetPrototypeBasedOn_1(builder.setPrototypeBasedOn)
	}
	this_.priv_FunctionType_refreshWellFormedTemplatizedFlag_8953()
	if int32(builder.noTypeKind) != int32(int8(0)) {
		this_.EagerlyResolveToSelf()
	} else {
		this_.GetJSTypeBase().registry.GetResolver().ResolveIfClosed(ToIface_JSType(this_), StaticGet_FunctionType_TYPE_CLASS())
	}
}

// Java static method: FunctionType.createNoObjectType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func FunctionType_CreateNoObjectType(registry *JSTypeRegistry) *FunctionType {
	clinit_FunctionType()
	_ = registry
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, registry)
	(&fb_).WithKind(StaticGet_FunctionType_Kind_NONE())
	(&fb_).WithReturnsOwnInstanceType()
	(&fb_).ForNativeType()
	(&fb_).priv_Builder_withNoTypeKind_9051(int8(1))
	return j2gort.CastPtr[FunctionType]((&fb_).Build())
}

// Java static method: FunctionType.createNoType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func FunctionType_CreateNoType(registry *JSTypeRegistry) *FunctionType {
	clinit_FunctionType()
	_ = registry
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, registry)
	(&fb_).WithKind(StaticGet_FunctionType_Kind_NONE())
	(&fb_).WithReturnsOwnInstanceType()
	(&fb_).ForNativeType()
	(&fb_).priv_Builder_withNoTypeKind_9051(int8(2))
	return j2gort.CastPtr[FunctionType]((&fb_).Build())
}

// Java method: FunctionType.hasPrototypeSlot()
func (this_ *FunctionType) HasPrototypeSlot() bool {
	return ((this_.GetPrototypeObjectTypeBase().cold != nil) && (this_.GetPrototypeObjectTypeBase().cold.prototypeSlot != nil))
}

// Java method: FunctionType.getPrototypeSlot()
func (this_ *FunctionType) priv_FunctionType_getPrototypeSlot_8951() *Property {
	if this_.GetPrototypeObjectTypeBase().cold != nil {
		return this_.GetPrototypeObjectTypeBase().cold.prototypeSlot
	} else {
		return nil
	}
}

// Java method: FunctionType.getPropAccessCode()
func (this_ *FunctionType) priv_FunctionType_getPropAccessCode_8952() int8 {
	var c *PrototypeObjectType_ColdData = this_.GetPrototypeObjectTypeBase().cold
	_ = c
	if (c != nil) && (int32(c.propAccess) != int32(int8(0))) {
		return c.propAccess
	}
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(3)) == int32(1) {
		return int8(1)
	} else {
		return int8(0)
	}
}

// Java static method: FunctionType.encodeKind(com.google.javascript.rhino.jstype.FunctionType.Kind kind)
func FunctionType_EncodeKind(kind FunctionType_Kind) int32 {
	clinit_FunctionType()
	_ = kind
	switch (kind).Ordinal() {
	case int32(0) /* Kind.ORDINARY */ :
		return int32(0)
	case int32(1) /* Kind.CONSTRUCTOR */ :
		return int32(1)
	case int32(2) /* Kind.INTERFACE */ :
		return int32(2)
	case int32(3) /* Kind.NONE */ :
		return int32(3)
	}
	panic("unreachable switch expression")
}

// Java static method: FunctionType.copyListToParameterArray(java.util.List<com.google.javascript.rhino.jstype.FunctionType.Parameter> list)
func FunctionType_CopyListToParameterArray(list j2gort.Ext_java_util_List) []FunctionType_Parameter {
	clinit_FunctionType()
	_ = list
	var n int32 = list.Size()
	_ = n
	if n == int32(0) {
		return StaticGet_FunctionType_EMPTY_PARAMETERS()
	}
	var arr []FunctionType_Parameter = make([]FunctionType_Parameter, n)
	_ = arr
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			{
				a_5 := arr
				i_6 := i
				a_5[i_6] = j2gort.ListGetVal[FunctionType_Parameter](list, i)
			}
		}
	}
	return arr
}

// Java static method: FunctionType.computeFunctionKindMask(com.google.javascript.rhino.jstype.FunctionType.Builder builder)
func FunctionType_ComputeFunctionKindMask(builder *FunctionType_Builder) int32 {
	clinit_FunctionType()
	_ = builder
	var mask int32 = int32(0)
	_ = mask
	switch builder.noTypeKind {
	case int8(1):
		mask = int32(4)
	case int8(2):
		mask = int32(2)
	default:
		mask = int32(128)
	}
	switch (builder.kind).Ordinal() {
	case int32(1) /* Kind.CONSTRUCTOR */ :
		mask = (mask | int32(4194304))
	case int32(2) /* Kind.INTERFACE */ :
		mask = (mask | int32(8388608))
	case int32(0) /* Kind.ORDINARY */ :
		mask = (mask | int32(16777216))
	default:
	}
	return mask
}

// Java method: FunctionType.refreshWellFormedTemplatizedFlag()
func (this_ *FunctionType) priv_FunctionType_refreshWellFormedTemplatizedFlag_8953() {
	var ret JSType = this_.returnType
	_ = ret
	if (ret == nil || ret.IsJavaNil_()) || (ret.GetJSTypeBase().GetWellFormedTemplatizedKind() != int32(0)) {
		this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) & ^(int32(int8(-128)))))
		return
	}
	var params []FunctionType_Parameter = this_.parameters
	_ = params
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(params))
		_ = len_
		for ; i < len_; i = i + 1 {
			var pt JSType = params[i].GetJSType()
			_ = pt
			if (pt == nil || pt.IsJavaNil_()) || (pt.GetJSTypeBase().GetWellFormedTemplatizedKind() != int32(0)) {
				this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) & ^(int32(int8(-128)))))
				return
			}
		}
	}
	this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) | int32(int8(-128))))
}

// Java method: FunctionType.getTypeClass()
func (this_ *FunctionType) GetTypeClass() JSTypeClass {
	var m int32 = this_.GetJSTypeBase().typeKindMask
	_ = m
	if (m & (int32(4) | int32(2))) != int32(0) {
		if (m & int32(4)) != int32(0) {
			return StaticGet_JSTypeClass_NO_OBJECT()
		} else {
			return StaticGet_JSTypeClass_NO()
		}
	}
	return StaticGet_JSTypeClass_FUNCTION()
}

// Java method: FunctionType.getConstructor()
func (this_ *FunctionType) GetConstructor() *FunctionType {
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0) {
		return nil
	}
	return j2gort.CheckNotNullPtr(this_.GetJSTypeBase().registry.GetNativeFunctionType(StaticGet_JSTypeNative_FUNCTION_FUNCTION_TYPE()))
}

// Java method: FunctionType.isNullable()
func (this_ *FunctionType) IsNullable() bool {
	return ((int32(this_.GetJSTypeBase().fnFlags) & int32(12)) == int32(8))
}

// Java method: FunctionType.isVoidable()
func (this_ *FunctionType) IsVoidable() bool {
	return ((int32(this_.GetJSTypeBase().fnFlags) & int32(12)) == int32(8))
}

// Java method: FunctionType.getPossibleToBooleanOutcomes()
func (this_ *FunctionType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) == int32(8) {
		return StaticGet_BooleanLiteralSet_BOTH()
	} else {
		return this_.PrototypeObjectTypeBase.GetPossibleToBooleanOutcomes()
	}
}

// Java method: FunctionType.getImplicitPrototype()
func (this_ *FunctionType) GetImplicitPrototype() ObjectType {
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0) {
		return nil
	} else {
		return this_.PrototypeObjectTypeBase.GetImplicitPrototype()
	}
}

// Java method: FunctionType.getReferenceName()
func (this_ *FunctionType) GetReferenceName() string {
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0) {
		return j2gort.NullString
	} else {
		return this_.PrototypeObjectTypeBase.GetReferenceName()
	}
}

// Java method: FunctionType.matchesNumberContext()
func (this_ *FunctionType) MatchesNumberContext() bool {
	return (((int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0)) || this_.PrototypeObjectTypeBase.MatchesNumberContext())
}

// Java method: FunctionType.matchesStringContext()
func (this_ *FunctionType) MatchesStringContext() bool {
	return (((int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0)) || this_.PrototypeObjectTypeBase.MatchesStringContext())
}

// Java method: FunctionType.matchesSymbolContext()
func (this_ *FunctionType) MatchesSymbolContext() bool {
	return (((int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0)) || this_.PrototypeObjectTypeBase.MatchesSymbolContext())
}

// Java method: FunctionType.setPropertyJSDocInfo(java.lang.String propertyName, com.google.javascript.rhino.JSDocInfo info)
func (this_ *FunctionType) SetPropertyJSDocInfo_2_string_JSDocInfo(propertyName string, info *JSDocInfo) {
	_ = propertyName
	_ = info
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0) {
		j2gort.CheckNotNullString(propertyName)
		return
	}
	this_.PrototypeObjectTypeBase.SetPropertyJSDocInfo_2_string_JSDocInfo(propertyName, info)
}

// Java method: FunctionType.setPropertyJSDocInfo(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.JSDocInfo info)
func (this_ *FunctionType) SetPropertyJSDocInfo_2_Property_Key_JSDocInfo(propertyName Property_Key, info *JSDocInfo) {
	_ = propertyName
	_ = info
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0) {
		return
	}
	this_.PrototypeObjectTypeBase.SetPropertyJSDocInfo_2_Property_Key_JSDocInfo(propertyName, info)
}

// Java method: FunctionType.isInstanceType()
func (this_ *FunctionType) IsInstanceType() bool {
	return j2gort.IfaceEquals[JSType](this_, this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_FUNCTION_TYPE()))
}

// Java method: FunctionType.getKind()
func (this_ *FunctionType) GetKind() FunctionType_Kind {
	switch int32(this_.GetJSTypeBase().fnFlags) & int32(3) {
	case int32(0):
		return StaticGet_FunctionType_Kind_ORDINARY()
	case int32(1):
		return StaticGet_FunctionType_Kind_CONSTRUCTOR()
	case int32(2):
		return StaticGet_FunctionType_Kind_INTERFACE()
	default:
		return StaticGet_FunctionType_Kind_NONE()
	}
	panic("unreachable switch expression")
}

// Java method: FunctionType.makesStructs()
func (this_ *FunctionType) MakesStructs() bool {
	var access int8 = this_.priv_FunctionType_getPropAccessCode_8952()
	_ = access
	if !(this_.HasInstanceType()) || (int32(access) == int32(int8(0))) {
		return false
	}
	{
		var superc *FunctionType = nil
		_ = superc
		switch access {
		case int8(3):
			return true
		case int8(4), int8(2):
			return false
		case int8(1):
			superc = this_.GetSuperClassConstructor()
			if (superc != nil) && superc.MakesStructs() {
				this_.SetStruct()
				return true
			}
			return false
		default:
			return false
		}
	}
	panic("unreachable")
}

// Java method: FunctionType.makesDicts()
func (this_ *FunctionType) MakesDicts() bool {
	var access int8 = this_.priv_FunctionType_getPropAccessCode_8952()
	_ = access
	if !(this_.IsConstructor()) || (int32(access) == int32(int8(0))) {
		return false
	}
	{
		var superc *FunctionType = nil
		_ = superc
		switch access {
		case int8(4):
			return true
		case int8(3), int8(2):
			return false
		case int8(1):
			superc = this_.GetSuperClassConstructor()
			if (superc != nil) && superc.MakesDicts() {
				this_.SetDict()
				return true
			}
			return false
		default:
			return false
		}
	}
	panic("unreachable")
}

// Java method: FunctionType.setStruct()
func (this_ *FunctionType) SetStruct() {
	{
		r_7 := this_.EnsureCold()
		r_7.propAccess = int8(3)
	}
}

// Java method: FunctionType.setDict()
func (this_ *FunctionType) SetDict() {
	{
		r_8 := this_.EnsureCold()
		r_8.propAccess = int8(4)
	}
}

// Java method: FunctionType.setExplicitUnrestricted()
func (this_ *FunctionType) SetExplicitUnrestricted() {
	{
		r_9 := this_.EnsureCold()
		r_9.propAccess = int8(2)
	}
}

// Java method: FunctionType.toMaybeFunctionType()
func (this_ *FunctionType) ToMaybeFunctionType() *FunctionType {
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) == int32(0) {
		return this_
	} else {
		return nil
	}
}

// Java method: FunctionType.canBeCalled()
func (this_ *FunctionType) CanBeCalled() bool {
	return true
}

// Java method: FunctionType.hasImplementedInterfaces()
func (this_ *FunctionType) HasImplementedInterfaces() bool {
	if !(this_.GetOwnImplementedInterfaces().IsEmpty()) {
		return true
	}
	var superCtor *FunctionType
	_ = superCtor
	if this_.IsConstructor() {
		superCtor = this_.GetSuperClassConstructor()
	} else {
		superCtor = nil
	}
	if superCtor != nil {
		return superCtor.HasImplementedInterfaces()
	}
	return false
}

// Java method: FunctionType.getParameterArray()
func (this_ *FunctionType) GetParameterArray() []FunctionType_Parameter {
	return this_.parameters
}

// Java method: FunctionType.getParameterCount()
func (this_ *FunctionType) GetParameterCount() int32 {
	return int32(len(this_.parameters))
}

// Java method: FunctionType.getParameter(int index)
func (this_ *FunctionType) GetParameter(index int32) FunctionType_Parameter {
	_ = index
	return this_.parameters[index]
}

// Java method: FunctionType.getParameters()
func (this_ *FunctionType) GetParameters() j2gort.Ext_com_google_common_collect_ImmutableList {
	var params []FunctionType_Parameter = this_.parameters
	_ = params
	var n int32 = int32(len(params))
	_ = n
	if n == int32(0) {
		return StaticGet_FunctionType_EMPTY_PARAMETER_LIST()
	}
	if n == int32(1) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.BoxNullableVal(params[int32(0)]))
	}
	if n == int32(2) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.BoxNullableVal(params[int32(0)]), j2gort.BoxNullableVal(params[int32(1)]))
	}
	if n == int32(3) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_3(j2gort.BoxNullableVal(params[int32(0)]), j2gort.BoxNullableVal(params[int32(1)]), j2gort.BoxNullableVal(params[int32(2)]))
	}
	if n == int32(4) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_4(j2gort.BoxNullableVal(params[int32(0)]), j2gort.BoxNullableVal(params[int32(1)]), j2gort.BoxNullableVal(params[int32(2)]), j2gort.BoxNullableVal(params[int32(3)]))
	}
	var b *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, n)
	_ = b
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			b.AddVoid_(j2gort.BoxNullableVal(params[i]))
		}
	}
	return b.BuildImmList_()
}

// Java method: FunctionType.getMinArity()
func (this_ *FunctionType) GetMinArity() int32 {
	var params []FunctionType_Parameter = this_.parameters
	_ = params
	var min_ int32 = int32(0)
	_ = min_
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(params))
		_ = size
		for ; i < size; i = i + 1 {
			var parameter FunctionType_Parameter = params[i]
			_ = parameter
			if !(parameter.IsOptional()) && !(parameter.IsVariadic()) {
				min_ = (i + int32(1))
			}
		}
	}
	return min_
}

// Java method: FunctionType.getMaxArity()
func (this_ *FunctionType) GetMaxArity() int32 {
	var params []FunctionType_Parameter = this_.parameters
	_ = params
	var len_ int32 = int32(len(params))
	_ = len_
	if len_ == int32(0) {
		return int32(0)
	}
	var lastParam FunctionType_Parameter = params[(len_ - int32(1))]
	_ = lastParam
	if !(lastParam.IsVariadic()) {
		return len_
	}
	return int32(2147483647)
}

// Java method: FunctionType.getReturnType()
func (this_ *FunctionType) GetReturnType() JSType {
	return this_.returnType
}

// Java method: FunctionType.isReturnTypeInferred()
func (this_ *FunctionType) IsReturnTypeInferred() bool {
	return ((int32(this_.GetJSTypeBase().fnFlags) & int32(int8(16))) != int32(0))
}

// Java method: FunctionType.getSlot(com.google.javascript.rhino.jstype.Property.Key name)
func (this_ *FunctionType) GetSlot_1_Property_Key(name Property_Key) *Property {
	_ = name
	if name.Matches("prototype") {
		this_.GetPrototype()
		return this_.priv_FunctionType_getPrototypeSlot_8951()
	} else {
		return this_.PrototypeObjectTypeBase.GetSlot_1_Property_Key(name)
	}
}

// Java method: FunctionType.getOwnPropertyNames()
func (this_ *FunctionType) GetOwnPropertyNames() j2gort.Ext_java_util_Set {
	if this_.priv_FunctionType_getPrototypeSlot_8951() == nil {
		return this_.PrototypeObjectTypeBase.GetOwnPropertyNames()
	} else {
		var names j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
		_ = names
		j2gort.ImmSetBuilderAddStrVoid(names, "prototype")
		names.AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.PrototypeObjectTypeBase.GetOwnPropertyNames()))
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](names.Build())
	}
}

// Java method: FunctionType.getPrototypeProperty()
func (this_ *FunctionType) GetPrototypeProperty() ObjectType {
	return this_.GetPrototype()
}

// Java method: FunctionType.getPrototype()
func (this_ *FunctionType) GetPrototype() ObjectType {
	var slot *Property = this_.priv_FunctionType_getPrototypeSlot_8951()
	_ = slot
	if slot == nil {
		var refName string = this_.GetReferenceName()
		_ = refName
		if refName == j2gort.NullString {
			this_.priv_FunctionType_setPrototypeNoCheck_8980(this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_UNKNOWN_TYPE()), Node(0))
		} else {
			this_.SetPrototype(Upcast_ObjectType(j2gort.AsIface[PrototypeObjectType](j2gort.AsIface[PrototypeObjectType_Builder](PrototypeObjectType_S_Builder(this_.GetJSTypeBase().registry).SetName((j2gort.StrValString(this_.GetReferenceName()) + ".prototype"))).GetPrototypeObjectType_BuilderBase().SetImplicitPrototype(this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_OBJECT_TYPE())).GetPrototypeObjectType_BuilderBase().SetNative(this_.IsNativeObjectType()).Build())), Node(0))
		}
		slot = this_.priv_FunctionType_getPrototypeSlot_8951()
	}
	return j2gort.AsIface[ObjectType](slot.GetType())
}

// Java method: FunctionType.setPrototypeBasedOn(com.google.javascript.rhino.jstype.ObjectType baseType)
func (this_ *FunctionType) SetPrototypeBasedOn_1(baseType ObjectType) {
	_ = baseType
	this_.priv_FunctionType_setPrototypeBasedOn_8978(baseType, Node(0))
}

// Java method: FunctionType.setPrototypeBasedOn(com.google.javascript.rhino.jstype.ObjectType baseType, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node propertyNode)
func (this_ *FunctionType) priv_FunctionType_setPrototypeBasedOn_8978(baseType ObjectType, propertyNode Node) {
	_ = baseType
	_ = propertyNode
	if this_.IsNativeObjectType() {
		return
	}
	var src Node = this_.GetJSTypeBase().subSourceNode
	_ = src
	if (src != 0) && src.IsClass() {
		var superCtor *FunctionType = baseType.GetConstructor()
		_ = superCtor
		if superCtor != nil {
			this_.SetImplicitPrototype(ToIface_ObjectType(superCtor))
		}
		this_.MaybeLoosenTypecheckingDueToForwardReferencedSupertype(Upcast_JSType(baseType))
	}
	if (baseType.GetObjectTypeBase().HasReferenceName() || this_.IsNativeObjectType()) || baseType.GetObjectTypeBase().IsFunctionPrototypeType() {
		if ((this_.priv_FunctionType_getPrototypeSlot_8951() != nil) && this_.HasInstanceType()) && baseType.GetObjectTypeBase().EqualsJSType(Upcast_JSType(this_.GetInstanceType())) {
			return
		}
		baseType = Upcast_ObjectType(j2gort.AsIface[PrototypeObjectType](j2gort.AsIface[PrototypeObjectType_Builder](PrototypeObjectType_S_Builder(this_.GetJSTypeBase().registry).SetName((j2gort.StrValString(this_.GetReferenceName()) + ".prototype"))).GetPrototypeObjectType_BuilderBase().SetImplicitPrototype(baseType).Build()))
	}
	this_.SetPrototype(baseType, propertyNode)
}

// Java method: FunctionType.setPrototype(com.google.javascript.rhino.jstype.ObjectType prototype, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node propertyNode)
func (this_ *FunctionType) SetPrototype(prototype ObjectType, propertyNode Node) bool {
	_ = prototype
	_ = propertyNode
	if prototype == nil || prototype.IsJavaNil_() {
		return false
	}
	if this_.IsConstructor() && j2gort.IfaceEquals[ObjectType](prototype, this_.GetInstanceType()) {
		return false
	}
	return this_.priv_FunctionType_setPrototypeNoCheck_8980(prototype, propertyNode)
}

// Java method: FunctionType.setPrototypeNoCheck(com.google.javascript.rhino.jstype.ObjectType prototype, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node propertyNode)
func (this_ *FunctionType) priv_FunctionType_setPrototypeNoCheck_8980(prototype ObjectType, propertyNode Node) bool {
	_ = prototype
	_ = propertyNode
	var slot *Property = this_.priv_FunctionType_getPrototypeSlot_8951()
	_ = slot
	var oldPrototype ObjectType
	_ = oldPrototype
	if slot == nil {
		oldPrototype = nil
	} else {
		oldPrototype = j2gort.AsIface[ObjectType](slot.GetType())
	}
	var replacedPrototype bool = (oldPrototype != nil && !oldPrototype.IsJavaNil_())
	_ = replacedPrototype
	{
		r_10 := this_.EnsureCold()
		r_10.prototypeSlot = NewProperty_4_string_JSType_bool_Node("prototype", Upcast_JSType(prototype), true, func() Node {
			if propertyNode == 0 {
				return this_.GetJSTypeBase().subSourceNode
			}
			return propertyNode
		}())
	}
	prototype.SetOwnerFunction(this_)
	if oldPrototype != nil && !oldPrototype.IsJavaNil_() {
		oldPrototype.SetOwnerFunction(nil)
	}
	if replacedPrototype {
		this_.ClearCachedValues()
	}
	return true
}

// Java method: FunctionType.getAllImplementedInterfaces()
func (this_ *FunctionType) GetAllImplementedInterfaces() j2gort.Ext_java_lang_Iterable {
	var implemented []ObjectType = j2gort.ImmListToSliceNullable[ObjectType](this_.GetImplementedInterfaces())
	_ = implemented
	var n int32 = int32(len(implemented))
	_ = n
	if n == int32(0) {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST())
	}
	if n == int32(1) {
		var single ObjectType = implemented[int32(0)]
		_ = single
		var ctor *FunctionType = single.GetConstructor()
		_ = ctor
		if (ctor == nil) || !(ctor.IsInterface()) {
			return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST())
		}
		if ctor.GetExtendedInterfacesCount() == int32(0) {
			return j2gort.ImmListFromTypedSlice(implemented)
		}
	}
	var interfaces j2gort.Ext_java_util_Set = nil
	_ = interfaces
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			interfaces = this_.priv_FunctionType_addRelatedInterfaces_8982(implemented[i], interfaces)
		}
	}
	if interfaces != nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](interfaces)
	} else {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST())
	}
}

// Java method: FunctionType.addRelatedInterfaces(com.google.javascript.rhino.jstype.ObjectType instance, java.util.@org.jspecify.annotations.Nullable Set<com.google.javascript.rhino.jstype.ObjectType> set)
func (this_ *FunctionType) priv_FunctionType_addRelatedInterfaces_8982(instance ObjectType, set j2gort.Ext_java_util_Set) j2gort.Ext_java_util_Set {
	_ = instance
	_ = set
	var constructor *FunctionType = instance.GetConstructor()
	_ = constructor
	if constructor != nil {
		if !(constructor.IsInterface()) {
			return set
		}
		if set == nil {
			set = j2gort.NewExt_java_util_LinkedHashSet()
		}
		if !(set.Add_1(instance)) {
			return set
		}
		if constructor.GetExtendedInterfacesCount() > int32(0) {
			{
				cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](instance.GetCtorExtendedInterfaces())))
				for cur_11.HasNext() {
					rawElem_12 := cur_11.Next()
					interfaceType := j2gort.AsIface[ObjectType](rawElem_12)
					_ = interfaceType
					set = this_.priv_FunctionType_addRelatedInterfaces_8982(interfaceType, set)
				}
			}
		}
	}
	return set
}

// Java method: FunctionType.getAncestorInterfaces()
func (this_ *FunctionType) GetAncestorInterfaces() j2gort.Ext_java_util_Collection {
	var result j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = result
	if this_.IsConstructor() {
		result.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.GetImplementedInterfaces()))
	} else {
		result.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.GetExtendedInterfaces()))
	}
	return result
}

// Java method: FunctionType.getImplementedInterfaces()
func (this_ *FunctionType) GetImplementedInterfaces() j2gort.Ext_com_google_common_collect_ImmutableList {
	var ownImplemented []ObjectType = j2gort.ImmListToSliceNullable[ObjectType](this_.GetOwnImplementedInterfaces())
	_ = ownImplemented
	var superCtor *FunctionType
	_ = superCtor
	if this_.IsConstructor() {
		superCtor = this_.GetSuperClassConstructor()
	} else {
		superCtor = nil
	}
	for (superCtor != nil) && superCtor.GetOwnImplementedInterfaces().IsEmpty() {
		superCtor = superCtor.GetSuperClassConstructor()
	}
	if superCtor == nil {
		return j2gort.ImmListFromTypedSlice(ownImplemented)
	}
	var builder []ObjectType = []ObjectType{}
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(ownImplemented))
		_ = n
		for ; i < n; i = i + 1 {
			builder = append(builder, j2gort.CheckNotNullElem(ownImplemented[i]))
		}
	}
	for superCtor != nil {
		var superOwn []ObjectType = j2gort.ImmListToSliceNullable[ObjectType](superCtor.GetOwnImplementedInterfaces())
		_ = superOwn
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var n_2 int32 = int32(len(superOwn))
			_ = n_2
			for ; i_2 < n_2; i_2 = i_2 + 1 {
				builder = append(builder, j2gort.CheckNotNullElem(superOwn[i_2]))
			}
		}
		superCtor = superCtor.GetSuperClassConstructor()
	}
	return j2gort.ImmListFromTypedSlice(builder)
}

// Java method: FunctionType.getOwnImplementedInterfaces()
func (this_ *FunctionType) GetOwnImplementedInterfaces() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.GetPrototypeObjectTypeBase().cold != nil {
		return this_.GetPrototypeObjectTypeBase().cold.implementedInterfaces
	} else {
		return StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST()
	}
}

// Java method: FunctionType.setImplementedInterfaces(java.util.List<com.google.javascript.rhino.jstype.ObjectType> implementedInterfaces)
func (this_ *FunctionType) SetImplementedInterfaces(implementedInterfaces j2gort.Ext_java_util_List) {
	_ = implementedInterfaces
	if !(this_.IsConstructor()) {
		j2gort.CheckState(false)
	}
	var copy_ j2gort.Ext_com_google_common_collect_ImmutableList
	_ = copy_
	if implementedInterfaces.IsEmpty() {
		copy_ = StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST()
	} else {
		copy_ = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](implementedInterfaces))
	}
	if !(copy_.IsEmpty()) || (this_.GetPrototypeObjectTypeBase().cold != nil) {
		{
			r_13 := this_.EnsureCold()
			r_13.implementedInterfaces = copy_
		}
	}
	{
		cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNull(implementedInterfaces))
		for cur_14.HasNext() {
			rawElem_15 := cur_14.Next()
			type_ := j2gort.AsIface[ObjectType](rawElem_15)
			_ = type_
			this_.typeOfThis.GetJSTypeBase().MergeSupertypeTemplateTypes(type_)
		}
	}
}

// Java method: FunctionType.getExtendedInterfaces()
func (this_ *FunctionType) GetExtendedInterfaces() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.GetPrototypeObjectTypeBase().cold != nil {
		return this_.GetPrototypeObjectTypeBase().cold.extendedInterfaces
	} else {
		return StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST()
	}
}

// Java method: FunctionType.getExtendedInterfacesCount()
func (this_ *FunctionType) GetExtendedInterfacesCount() int32 {
	if this_.GetPrototypeObjectTypeBase().cold != nil {
		return this_.GetPrototypeObjectTypeBase().cold.extendedInterfaces.Size()
	} else {
		return int32(0)
	}
}

// Java method: FunctionType.setExtendedInterfaces(java.util.List<com.google.javascript.rhino.jstype.ObjectType> extendedInterfaces)
func (this_ *FunctionType) SetExtendedInterfaces(extendedInterfaces j2gort.Ext_java_util_List) {
	_ = extendedInterfaces
	if !(this_.IsInterface()) {
		j2gort.CheckState(false)
	}
	var copy_ j2gort.Ext_com_google_common_collect_ImmutableList
	_ = copy_
	if extendedInterfaces.IsEmpty() {
		copy_ = StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST()
	} else {
		copy_ = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](extendedInterfaces))
	}
	if !(copy_.IsEmpty()) || (this_.GetPrototypeObjectTypeBase().cold != nil) {
		{
			r_16 := this_.EnsureCold()
			r_16.extendedInterfaces = copy_
		}
	}
	{
		cur_17 := j2gort.NewIterCursor(j2gort.CheckNotNull(extendedInterfaces))
		for cur_17.HasNext() {
			rawElem_18 := cur_17.Next()
			extendedInterface := j2gort.AsIface[ObjectType](rawElem_18)
			_ = extendedInterface
			this_.typeOfThis.GetJSTypeBase().MergeSupertypeTemplateTypes(extendedInterface)
		}
	}
}

// Java method: FunctionType.getPropertyType(com.google.javascript.rhino.jstype.Property.Key name)
func (this_ *FunctionType) GetPropertyType_1_Property_Key(name Property_Key) JSType {
	_ = name
	if !(this_.HasOwnProperty_1_Property_Key(name)) {
		var isCall bool = name.Matches("call")
		_ = isCall
		var isBind bool = name.Matches("bind")
		_ = isBind
		if isCall || isBind {
			this_.DefineDeclaredProperty_3_Property_Key_JSType_Node(name, ToIface_JSType(this_.priv_FunctionType_getCallOrBindSignature_8992(isCall)), this_.GetJSTypeBase().subSourceNode)
		} else if name.Matches("apply") {
			var val_19 FunctionParamBuilder
			InitFunctionParamBuilder_2(&val_19, this_.GetJSTypeBase().registry, int32(2))
			var builder *FunctionParamBuilder = &val_19
			_ = builder
			builder.AddOptionalParams([]JSType{this_.GetJSTypeBase().registry.CreateNullableType(this_.GetTypeOfThis()), this_.GetJSTypeBase().registry.CreateNullableType(this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_OBJECT_TYPE()))})
			var fb_ FunctionType_Builder
			InitFunctionType_Builder(&fb_, this_.GetJSTypeBase().registry)
			(&fb_).WithParameterArray(builder.BuildArray())
			(&fb_).WithReturnType_1(this_.GetReturnType())
			(&fb_).WithTemplateKeysFrom(this_.GetTemplateTypeMap())
			this_.DefineDeclaredProperty_3_Property_Key_JSType_Node(name, ToIface_JSType(j2gort.CastPtr[FunctionType]((&fb_).Build())), this_.GetJSTypeBase().subSourceNode)
		}
	}
	return this_.PrototypeObjectTypeBase.GetPropertyType_1_Property_Key(name)
}

// Java method: FunctionType.getBindReturnType(int argsToBind)
func (this_ *FunctionType) GetBindReturnType(argsToBind int32) *FunctionType {
	_ = argsToBind
	var builder *FunctionType_Builder = FunctionType_S_Builder(this_.GetJSTypeBase().registry).WithReturnType_1(this_.GetReturnType()).WithTemplateKeysFrom(this_.GetTemplateTypeMap())
	_ = builder
	if argsToBind < int32(0) {
		return j2gort.CastPtr[FunctionType](builder.Build())
	}
	var origParams []FunctionType_Parameter = this_.parameters
	_ = origParams
	var origLen int32 = int32(len(origParams))
	_ = origLen
	var skip int32 = int32(0)
	_ = skip
	{
		var i int32 = int32(1)
		_ = i
		for ; (i < argsToBind) && (skip < origLen); i = i + 1 {
			if origParams[skip].IsVariadic() {
				break
			}
			skip = skip + 1
		}
	}
	if skip == int32(0) {
		builder.WithParameterArray(origParams)
	} else if skip >= origLen {
		builder.WithParameterArray(StaticGet_FunctionType_EMPTY_PARAMETERS())
	} else {
		var rem int32 = (origLen - skip)
		_ = rem
		var boundParams []FunctionType_Parameter = make([]FunctionType_Parameter, rem)
		_ = boundParams
		func() {
			acSrc_20 := origParams
			acSP_21 := skip
			acDst_22 := boundParams
			acDP_23 := int32(0)
			acLen_24 := rem
			j2gort.CheckNotNullSlice(acSrc_20)
			j2gort.CheckNotNullSlice(acDst_22)
			if acSP_21 < 0 || acDP_23 < 0 || acLen_24 < 0 || int64(acSP_21)+int64(acLen_24) > int64(len(acSrc_20)) || int64(acDP_23)+int64(acLen_24) > int64(len(acDst_22)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_22[int(acDP_23):int(acDP_23)+int(acLen_24)], acSrc_20[int(acSP_21):int(acSP_21)+int(acLen_24)])
		}()
		builder.WithParameterArray(boundParams)
	}
	return j2gort.CastPtr[FunctionType](builder.Build())
}

// Java method: FunctionType.getCallOrBindSignature(boolean isCall)
func (this_ *FunctionType) priv_FunctionType_getCallOrBindSignature_8992(isCall bool) *FunctionType {
	_ = isCall
	var isBind bool = !(isCall)
	_ = isBind
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, this_.GetJSTypeBase().registry)
	(&fb_).WithReturnType_1(func() JSType {
		if isCall {
			return this_.GetReturnType()
		}
		return ToIface_JSType(this_.GetBindReturnType(int32(-1)))
	}())
	(&fb_).WithTemplateKeysFrom(this_.GetTemplateTypeMap())
	var builder *FunctionType_Builder = &fb_

	_ = builder
	var origParams []FunctionType_Parameter = this_.parameters
	_ = origParams
	var origLen int32 = int32(len(origParams))
	_ = origLen
	var params []FunctionType_Parameter = make([]FunctionType_Parameter, (origLen + int32(1)))
	_ = params
	var thisJSType JSType = this_.GetJSTypeBase().registry.CreateOptionalNullableType(this_.GetTypeOfThis())
	_ = thisJSType
	var firstOptional bool = false
	_ = firstOptional
	if !(isBind) {
		var firstArg FunctionType_Parameter
		_ = firstArg
		if origLen > int32(0) {
			firstArg = origParams[int32(0)]
		} else {
			firstArg = (FunctionType_Parameter{})
		}
		if ((firstArg == (FunctionType_Parameter{})) || firstArg.IsOptional()) || firstArg.IsVariadic() {
			firstOptional = true
		}
	}
	{
		a_25 := params
		i_26 := int32(0)
		a_25[i_26] = FunctionType_Parameter_Create(thisJSType, firstOptional, false)
	}
	if isBind {
		{
			var i int32 = int32(0)
			_ = i
			for ; i < origLen; i = i + 1 {
				var current FunctionType_Parameter = origParams[i]
				_ = current
				{
					a_27 := params
					i_28 := (i + int32(1))
					a_27[i_28] = FunctionType_Parameter_Create(current.GetJSType(), true, current.IsVariadic())
				}
			}
		}
	} else {
		func() {
			acSrc_29 := origParams
			acSP_30 := int32(0)
			acDst_31 := params
			acDP_32 := int32(1)
			acLen_33 := origLen
			j2gort.CheckNotNullSlice(acSrc_29)
			j2gort.CheckNotNullSlice(acDst_31)
			if acSP_30 < 0 || acDP_32 < 0 || acLen_33 < 0 || int64(acSP_30)+int64(acLen_33) > int64(len(acSrc_29)) || int64(acDP_32)+int64(acLen_33) > int64(len(acDst_31)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_31[int(acDP_32):int(acDP_32)+int(acLen_33)], acSrc_29[int(acSP_30):int(acSP_30)+int(acLen_33)])
		}()
	}
	builder.WithParameterArray(params)
	return j2gort.CastPtr[FunctionType](builder.Build())
}

// Java method: FunctionType.defineProperty(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *FunctionType) DefineProperty_4_string_JSType_bool_Node(name string, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = name
	_ = type_
	_ = inferred
	_ = propertyNode
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0) {
		return true
	}
	if "prototype" == name {
		var objType ObjectType = type_.GetJSTypeBase().ToObjectType()
		_ = objType
		if objType != nil && !objType.IsJavaNil_() {
			var slot *Property = this_.priv_FunctionType_getPrototypeSlot_8951()
			_ = slot
			if (slot != nil) && objType.GetObjectTypeBase().EqualsJSType(slot.GetType()) {
				return true
			}
			this_.priv_FunctionType_setPrototypeBasedOn_8978(objType, propertyNode)
			return true
		} else {
			return false
		}
	}
	return this_.PrototypeObjectTypeBase.DefineProperty_4_string_JSType_bool_Node(name, type_, inferred, propertyNode)
}

// Java method: FunctionType.defineProperty(com.google.javascript.rhino.jstype.Property.Key name, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *FunctionType) DefineProperty_4_Property_Key_JSType_bool_Node(name Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = name
	_ = type_
	_ = inferred
	_ = propertyNode
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0) {
		return true
	}
	if name.Matches("prototype") {
		var objType ObjectType = type_.GetJSTypeBase().ToObjectType()
		_ = objType
		if objType != nil && !objType.IsJavaNil_() {
			var slot *Property = this_.priv_FunctionType_getPrototypeSlot_8951()
			_ = slot
			if (slot != nil) && objType.GetObjectTypeBase().EqualsJSType(slot.GetType()) {
				return true
			}
			this_.priv_FunctionType_setPrototypeBasedOn_8978(objType, propertyNode)
			return true
		} else {
			return false
		}
	}
	return this_.PrototypeObjectTypeBase.DefineProperty_4_Property_Key_JSType_bool_Node(name, type_, inferred, propertyNode)
}

// Java method: FunctionType.supAndInfHelper(com.google.javascript.rhino.jstype.FunctionType that, boolean leastSuper)
func (this_ *FunctionType) SupAndInfHelper(that *FunctionType, leastSuper bool) *FunctionType {
	_ = that
	_ = leastSuper
	j2gort.CheckNotNullPtr(that)
	if this_.EqualsJSType(ToIface_JSType(that)) {
		return this_
	}
	if ((this_.IsOrdinaryFunction() && that.IsOrdinaryFunction()) && !(this_.priv_FunctionType_hasUnknownParamsOrReturn_8994())) && !(that.priv_FunctionType_hasUnknownParamsOrReturn_8994()) {
		var isSubtypeOfThat bool = this_.IsSubtype_1(ToIface_JSType(that))
		_ = isSubtypeOfThat
		var isSubtypeOfThis bool = that.IsSubtype_1(this_)
		_ = isSubtypeOfThis
		if isSubtypeOfThat && !(isSubtypeOfThis) {
			if leastSuper {
				return that
			} else {
				return this_
			}
		} else if isSubtypeOfThis && !(isSubtypeOfThat) {
			if leastSuper {
				return this_
			} else {
				return that
			}
		}
		var merged *FunctionType = this_.priv_FunctionType_tryMergeFunctionPiecewise_8995(that, leastSuper)
		_ = merged
		if merged != nil {
			return merged
		}
	}
	var functionInstance JSType = this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_FUNCTION_TYPE())
	_ = functionInstance
	base_functionInstance_ := functionInstance.GetJSTypeBase()
	if base_functionInstance_.EqualsJSType(ToIface_JSType(that)) {
		if leastSuper {
			return that
		} else {
			return this_
		}
	} else if base_functionInstance_.EqualsJSType(this_) {
		if leastSuper {
			return this_
		} else {
			return that
		}
	}
	var greatestFn *FunctionType = this_.GetJSTypeBase().registry.GetNativeFunctionType(StaticGet_JSTypeNative_FUNCTION_TYPE())
	_ = greatestFn
	var leastFn *FunctionType = this_.GetJSTypeBase().registry.GetNativeFunctionType(StaticGet_JSTypeNative_LEAST_FUNCTION_TYPE())
	_ = leastFn
	if leastSuper {
		return greatestFn
	} else {
		return leastFn
	}
}

// Java method: FunctionType.hasUnknownParamsOrReturn()
func (this_ *FunctionType) priv_FunctionType_hasUnknownParamsOrReturn_8994() bool {
	var params []FunctionType_Parameter = this_.parameters
	_ = params
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(params))
		_ = size
		for ; i < size; i = i + 1 {
			if params[i].GetJSType().GetJSTypeBase().IsUnknownType() {
				return true
			}
		}
	}
	return ((this_.returnType == nil || this_.returnType.IsJavaNil_()) || this_.returnType.GetJSTypeBase().IsUnknownType())
}

// Java method: FunctionType.tryMergeFunctionPiecewise(com.google.javascript.rhino.jstype.FunctionType other, boolean leastSuper)
func (this_ *FunctionType) priv_FunctionType_tryMergeFunctionPiecewise_8995(other *FunctionType, leastSuper bool) *FunctionType {
	_ = other
	_ = leastSuper
	var newParamsNode []FunctionType_Parameter = nil
	_ = newParamsNode
	if EqualityChecker_CheckParameters(StaticGet_EqualityChecker_EqMethod_IDENTITY(), this_, other) {
		newParamsNode = this_.parameters
	} else {
		return nil
	}
	var newReturnType JSType
	_ = newReturnType
	if leastSuper {
		newReturnType = this_.returnType.GetLeastSupertype(other.returnType)
	} else {
		newReturnType = this_.returnType.GetJSTypeBase().GetGreatestSubtype(other.returnType)
	}
	var newTypeOfThis JSType = nil
	_ = newTypeOfThis
	if func() bool {
		if this_.typeOfThis == nil || this_.typeOfThis.IsJavaNil_() {
			return other.typeOfThis == nil || other.typeOfThis.IsJavaNil_()
		}
		return this_.typeOfThis.GetJSTypeBase().EqualsJSType(other.typeOfThis)
	}() {
		newTypeOfThis = this_.typeOfThis
	} else {
		var maybeNewTypeOfThis JSType
		_ = maybeNewTypeOfThis
		if leastSuper {
			maybeNewTypeOfThis = this_.typeOfThis.GetLeastSupertype(other.typeOfThis)
		} else {
			maybeNewTypeOfThis = this_.typeOfThis.GetJSTypeBase().GetGreatestSubtype(other.typeOfThis)
		}
		newTypeOfThis = maybeNewTypeOfThis
	}
	var newReturnTypeInferred bool = (this_.IsReturnTypeInferred() || other.IsReturnTypeInferred())
	_ = newReturnTypeInferred
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, this_.GetJSTypeBase().registry)
	(&fb_).WithParameterArray(newParamsNode)
	(&fb_).WithReturnType_2(newReturnType, newReturnTypeInferred)
	(&fb_).WithTypeOfThis(newTypeOfThis)
	return j2gort.CastPtr[FunctionType]((&fb_).Build())
}

// Java method: FunctionType.getSuperClassConstructor()
func (this_ *FunctionType) GetSuperClassConstructor() *FunctionType {
	if !(this_.IsConstructor() || this_.IsInterface()) {
		j2gort.CheckArgument(false)
	}
	var maybeSuperInstanceType ObjectType = this_.GetPrototype().GetImplicitPrototype()
	_ = maybeSuperInstanceType
	if maybeSuperInstanceType == nil || maybeSuperInstanceType.IsJavaNil_() {
		return nil
	}
	return maybeSuperInstanceType.GetConstructor()
}

// Java method: FunctionType.recursionUnsafeHashCode()
func (this_ *FunctionType) RecursionUnsafeHashCode() int32 {
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0) {
		return j2gort.IdentityHash(this_)
	}
	var kind FunctionType_Kind = this_.GetKind()
	_ = kind
	var hc int32 = j2gort.CheckNotNullEnum(kind).HashCode()
	_ = hc
	switch (kind).Ordinal() {
	case int32(1) /* Kind.CONSTRUCTOR */, int32(2) /* Kind.INTERFACE */ :
		return ((int32(31) * hc) + j2gort.IdentityHash(this_))
	case int32(0) /* Kind.ORDINARY */ :
		hc = ((int32(31) * hc) + this_.typeOfThis.GetJSTypeBase().HashCode())
		var callHash int32
	iifeSwitch_1:
		switch {
		default:
			hi_34 := this_.returnType
			if hi_34 == nil || hi_34.IsJavaNil_() {
				callHash = 0
				break iifeSwitch_1
			}
			callHash = hi_34.GetJSTypeBase().HashCode()
			break iifeSwitch_1
		}

		_ = callHash
		var params []FunctionType_Parameter = this_.parameters
		_ = params
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = int32(len(params))
			_ = size
			for ; i < size; i = i + 1 {
				callHash = ((callHash * int32(31)) + func() int32 {
					hi_35 := params[i].GetJSType()
					if hi_35 == nil || hi_35.IsJavaNil_() {
						return 0
					}
					return hi_35.GetJSTypeBase().HashCode()
				}())
			}
		}
		hc = ((int32(31) * hc) + callHash)
		hc = ((int32(31) * hc) + (this_.GetClosurePrimitive()).HashCode())
		return hc
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	panic("unreachable switch expression")
}

// Java method: FunctionType.hasEqualCallType(com.google.javascript.rhino.jstype.FunctionType that)
func (this_ *FunctionType) HasEqualCallType(that *FunctionType) bool {
	_ = that
	return EqualityChecker_CheckCallType(StaticGet_EqualityChecker_EqMethod_IDENTITY(), this_, that)
}

// Java method: FunctionType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *FunctionType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	var noKindBits int32 = (int32(this_.GetJSTypeBase().fnFlags) & int32(12))
	_ = noKindBits
	if noKindBits != int32(0) {
		if noKindBits == int32(8) {
			var iifeArg_1 string
		iifeSwitch_2:
			switch {
			default:

				if sb.IsForAnnotations() {
					iifeArg_1 = "?"
					break iifeSwitch_2
				}
				iifeArg_1 = "None"
				break iifeSwitch_2
			}
			sb.Append_1_string(iifeArg_1)
		} else {
			var iifeArg_2 string
		iifeSwitch_3:
			switch {
			default:

				if sb.IsForAnnotations() {
					iifeArg_2 = "?"
					break iifeSwitch_3
				}
				iifeArg_2 = "NoObject"
				break iifeSwitch_3
			}
			sb.Append_1_string(iifeArg_2)
		}
		return
	}
	if !(this_.IsPrettyPrint()) || j2gort.IfaceEquals[JSType](this_, this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_FUNCTION_TYPE())) {
		var iifeArg_3 string
	iifeSwitch_4:
		switch {
		default:

			if sb.IsForAnnotations() {
				iifeArg_3 = "!Function"
				break iifeSwitch_4
			}
			iifeArg_3 = "Function"
			break iifeSwitch_4
		}
		sb.Append_1_string(iifeArg_3)
		return
	}
	if this_.HasInstanceType() && (this_.GetSource() != 0) {
		sb.Append_1_string("(typeof ").Append_1_JSType(Upcast_JSType(this_.GetInstanceType())).Append_1_string(")")
		return
	}
	this_.SetPrettyPrint(false)
	sb.Append_1_string("function(")
	var params []FunctionType_Parameter = this_.parameters
	_ = params
	var paramNum int32 = int32(len(params))
	_ = paramNum
	var hasKnownTypeOfThis bool = !(j2gort.InstanceOf[*UnknownType](this_.typeOfThis))
	_ = hasKnownTypeOfThis
	if hasKnownTypeOfThis {
		if this_.IsConstructor() {
			sb.Append_1_string("new:")
		} else {
			sb.Append_1_string("this:")
		}
		sb.Append_1_JSType(this_.typeOfThis)
	}
	if paramNum > int32(0) {
		if hasKnownTypeOfThis {
			sb.Append_1_string(", ")
		}
		var p FunctionType_Parameter = params[int32(0)]
		_ = p
		this_.priv_FunctionType_appendArgString_8998(sb, p)
		{
			var i int32 = int32(1)
			_ = i
			for ; i < paramNum; i = i + 1 {
				p = params[i]
				sb.Append_1_string(", ")
				this_.priv_FunctionType_appendArgString_8998(sb, p)
			}
		}
	}
	sb.Append_1_string("): ")
	sb.AppendNonNull(this_.returnType)
	this_.SetPrettyPrint(true)
}

// Java method: FunctionType.appendArgString(com.google.javascript.rhino.jstype.TypeStringBuilder sb, com.google.javascript.rhino.jstype.FunctionType.Parameter p)
func (this_ *FunctionType) priv_FunctionType_appendArgString_8998(sb *TypeStringBuilder, p FunctionType_Parameter) {
	_ = sb
	_ = p
	if p.IsVariadic() {
		this_.priv_FunctionType_appendVarArgsString_8999(sb, p.GetJSType())
	} else if p.IsOptional() {
		this_.priv_FunctionType_appendOptionalArgString_9000(sb, p.GetJSType())
	} else {
		sb.AppendNonNull(p.GetJSType())
	}
}

// Java method: FunctionType.appendVarArgsString(com.google.javascript.rhino.jstype.TypeStringBuilder sb, com.google.javascript.rhino.jstype.JSType paramType)
func (this_ *FunctionType) priv_FunctionType_appendVarArgsString_8999(sb *TypeStringBuilder, paramType JSType) {
	_ = sb
	_ = paramType
	sb.Append_1_string("...").AppendNonNull(paramType)
}

// Java method: FunctionType.appendOptionalArgString(com.google.javascript.rhino.jstype.TypeStringBuilder sb, com.google.javascript.rhino.jstype.JSType paramType)
func (this_ *FunctionType) priv_FunctionType_appendOptionalArgString_9000(sb *TypeStringBuilder, paramType JSType) {
	_ = sb
	_ = paramType
	if paramType.GetJSTypeBase().IsUnionType() {
		paramType = paramType.ToMaybeUnionType().GetRestrictedUnion(this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE()))
	}
	sb.AppendNonNull(paramType).Append_1_string("=")
}

// Java method: FunctionType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *FunctionType) Visit_1(visitor Visitor) any {
	_ = visitor
	switch int32(this_.GetJSTypeBase().fnFlags) & int32(12) {
	case int32(8):
		return visitor.CaseNoType(this_)
	case int32(4):
		return visitor.CaseNoObjectType()
	default:
		return visitor.CaseFunctionType(this_)
	}
	panic("unreachable switch expression")
}

// Java method: FunctionType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *FunctionType) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	switch int32(this_.GetJSTypeBase().fnFlags) & int32(12) {
	case int32(8):
		return visitor.CaseNoType(that)
	case int32(4):
		return visitor.CaseNoObjectType(that)
	default:
		return visitor.CaseFunctionType(this_, that)
	}
	panic("unreachable switch expression")
}

// Java method: FunctionType.getInstanceType()
func (this_ *FunctionType) GetInstanceType() ObjectType {
	if !(this_.HasInstanceType()) {
		j2gort.CheckState(false)
	}
	return this_.typeOfThis.GetJSTypeBase().ToObjectType()
}

// Java method: FunctionType.hasInstanceType()
func (this_ *FunctionType) HasInstanceType() bool {
	return (this_.IsConstructor() || this_.IsInterface())
}

// Java method: FunctionType.getTypeOfThis()
func (this_ *FunctionType) GetTypeOfThis() JSType {
	if this_.typeOfThis.GetJSTypeBase().IsEmptyType() {
		return Upcast_JSType(this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_UNKNOWN_TYPE()))
	} else {
		return this_.typeOfThis
	}
}

// Java method: FunctionType.getSource()
func (this_ *FunctionType) GetSource() Node {
	return this_.GetJSTypeBase().subSourceNode
}

// Java method: FunctionType.getGoogModuleId()
func (this_ *FunctionType) GetGoogModuleId() string {
	if this_.GetPrototypeObjectTypeBase().cold != nil {
		return this_.GetPrototypeObjectTypeBase().cold.googModuleId
	} else {
		return j2gort.NullString
	}
}

// Java method: FunctionType.setSource(com.google.javascript.rhino.Node source)
func (this_ *FunctionType) SetSource(source Node) {
	_ = source
	var slot *Property = this_.priv_FunctionType_getPrototypeSlot_8951()
	_ = slot
	if slot != nil {
		if (source == 0) || (slot.GetNode() == 0) {
			{
				r_36 := this_.EnsureCold()
				r_36.prototypeSlot = NewProperty_4_string_JSType_bool_Node(slot.GetName(), slot.GetType(), slot.IsTypeInferred(), source)
			}
		}
	}
	this_.GetJSTypeBase().subSourceNode = source
}

// Java method: FunctionType.clearCachedValues()
func (this_ *FunctionType) ClearCachedValues() {
	this_.PrototypeObjectTypeBase.ClearCachedValues()
	if !(this_.IsNativeObjectType()) {
		if this_.HasInstanceType() {
			this_.GetInstanceType().ClearCachedValues()
		}
		var slot *Property = this_.priv_FunctionType_getPrototypeSlot_8951()
		_ = slot
		if slot != nil {
			(j2gort.AsIface[ObjectType](slot.GetType())).ClearCachedValues()
		}
	}
}

// Java method: FunctionType.hasCachedValues()
func (this_ *FunctionType) HasCachedValues() bool {
	return ((this_.priv_FunctionType_getPrototypeSlot_8951() != nil) || this_.PrototypeObjectTypeBase.HasCachedValues())
}

// Java method: FunctionType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *FunctionType) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	if (int32(this_.GetJSTypeBase().fnFlags) & int32(12)) != int32(0) {
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	this_.returnType = JSType_SafeResolve(this_.returnType, reporter)
	var params []FunctionType_Parameter = this_.parameters
	_ = params
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(params))
		_ = size
		for ; i < size; i = i + 1 {
			params[i].GetJSType().GetJSTypeBase().Resolve(reporter)
		}
	}
	var slot *Property = this_.priv_FunctionType_getPrototypeSlot_8951()
	_ = slot
	if slot != nil {
		slot.SetType(JSType_SafeResolve(slot.GetType(), reporter))
	}
	var maybeTypeOfThis JSType = JSType_SafeResolve(this_.typeOfThis, reporter)
	_ = maybeTypeOfThis
	if maybeTypeOfThis != nil && !maybeTypeOfThis.IsJavaNil_() {
		if maybeTypeOfThis.GetJSTypeBase().IsNullType() || maybeTypeOfThis.GetJSTypeBase().IsVoidType() {
			this_.typeOfThis = maybeTypeOfThis
		} else {
			maybeTypeOfThis = Upcast_JSType(ObjectType_Cast(maybeTypeOfThis.RestrictByNotNullOrUndefined()))
			if maybeTypeOfThis != nil && !maybeTypeOfThis.IsJavaNil_() {
				this_.typeOfThis = maybeTypeOfThis
			}
		}
	}
	var c *PrototypeObjectType_ColdData = this_.GetPrototypeObjectTypeBase().cold
	_ = c
	if c != nil {
		if !(c.implementedInterfaces.IsEmpty()) {
			var resolvedImplemented j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_FunctionType_resolveTypeListHelper_9007(c.implementedInterfaces, reporter)
			_ = resolvedImplemented
			if resolvedImplemented != nil {
				{
					r_37 := c
					r_37.implementedInterfaces = resolvedImplemented
				}
			}
		}
		if !(c.extendedInterfaces.IsEmpty()) {
			var resolvedExtended j2gort.Ext_com_google_common_collect_ImmutableList = this_.priv_FunctionType_resolveTypeListHelper_9007(c.extendedInterfaces, reporter)
			_ = resolvedExtended
			if resolvedExtended != nil {
				{
					r_38 := c
					r_38.extendedInterfaces = resolvedExtended
				}
			}
		}
	}
	var resolved JSType = this_.PrototypeObjectTypeBase.ResolveInternal(reporter)
	_ = resolved
	this_.priv_FunctionType_refreshWellFormedTemplatizedFlag_8953()
	return resolved
}

// Java method: FunctionType.resolveTypeListHelper(com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.ObjectType> list, com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *FunctionType) priv_FunctionType_resolveTypeListHelper_9007(list j2gort.Ext_com_google_common_collect_ImmutableList, reporter ErrorReporter) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = list
	_ = reporter
	var changed bool = false
	_ = changed
	var resolvedList []ObjectType = []ObjectType{}
	_ = resolvedList
	var src Node = this_.GetJSTypeBase().subSourceNode
	_ = src
	{
		cur_39 := j2gort.NewIterCursor(j2gort.CheckNotNull(list))
		for cur_39.HasNext() {
			rawElem_40 := cur_39.Next()
			type_ := j2gort.AsIface[ObjectType](rawElem_40)
			_ = type_
			var rt JSType = type_.GetObjectTypeBase().Resolve(reporter)
			_ = rt
			if !(rt.GetJSTypeBase().IsObjectType()) {
				reporter.Warning(("not an object type: " + j2gort.StrVal(rt) + " (at " + j2gort.StrValString(this_.String()) + ")"), src.GetSourceFileName(), src.GetLineno(), src.GetCharno())
				continue
			}
			var resolved ObjectType = rt.GetJSTypeBase().ToObjectType()
			_ = resolved
			resolvedList = append(resolvedList, j2gort.CheckNotNullElem(resolved))
			{

				ba_41 := changed
				bb_42 := !(resolved == type_ || ((resolved == nil || resolved.IsJavaNil_()) && (type_ == nil || type_.IsJavaNil_())))
				changed = ba_41 || bb_42
			}

		}
	}
	if changed {
		return j2gort.ImmListFromTypedSlice(resolvedList)
	} else {
		return nil
	}
}

// Java method: FunctionType.hasAnyTemplateTypesInternal()
func (this_ *FunctionType) HasAnyTemplateTypesInternal() bool {
	return ((((this_.GetTemplateParamCount() > int32(0)) || this_.typeOfThis.GetJSTypeBase().HasAnyTemplateTypes()) || this_.returnType.GetJSTypeBase().HasAnyTemplateTypes()) || this_.priv_FunctionType_hasTemplatedParameterType_9009())
}

// Java method: FunctionType.hasTemplatedParameterType()
func (this_ *FunctionType) priv_FunctionType_hasTemplatedParameterType_9009() bool {
	var params []FunctionType_Parameter = this_.parameters
	_ = params
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(params))
		_ = size
		for ; i < size; i = i + 1 {
			if params[i].GetJSType().GetJSTypeBase().HasAnyTemplateTypes() {
				return true
			}
		}
	}
	return false
}

// Java method: FunctionType.hasProperties()
func (this_ *FunctionType) HasProperties() bool {
	var map_ *PropertyMap = this_.GetOwnPropertyMapOrNull()
	_ = map_
	return ((map_ != nil) && (map_.GetOwnPropertiesCount() > int32(0)))
}

// Java method: FunctionType.setImplicitMatch(boolean flag)
func (this_ *FunctionType) SetImplicitMatch(flag bool) {
	_ = flag
	if !(this_.IsInterface()) {
		j2gort.CheckState(false)
	}
	if flag {
		this_.GetJSTypeBase().fnFlags = int8((int32(this_.GetJSTypeBase().fnFlags) | int32(int8(32))))
	} else {
		this_.GetJSTypeBase().fnFlags = int8((int32(this_.GetJSTypeBase().fnFlags) & ^(int32(int8(32)))))
	}
}

// Java method: FunctionType.isStructuralInterface()
func (this_ *FunctionType) IsStructuralInterface() bool {
	return (this_.IsInterface() && ((int32(this_.GetJSTypeBase().fnFlags) & int32(int8(32))) != int32(0)))
}

// Java method: FunctionType.isAbstract()
func (this_ *FunctionType) IsAbstract() bool {
	return ((int32(this_.GetJSTypeBase().fnFlags) & int32(int8(64))) != int32(0))
}

// Java method: FunctionType.getPropertyTypeMap()
func (this_ *FunctionType) GetPropertyTypeMap() j2gort.Ext_java_util_Map {
	var propTypeMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = propTypeMap
	FunctionType_UpdatePropertyTypeMap(this_, propTypeMap, j2gort.NewExt_java_util_LinkedHashSet())
	return propTypeMap
}

// Java static method: FunctionType.updatePropertyTypeMap(com.google.javascript.rhino.jstype.FunctionType type, java.util.Map<java.lang.String,com.google.javascript.rhino.jstype.JSType> propTypeMap, java.util.LinkedHashSet<com.google.javascript.rhino.jstype.FunctionType> cache)
func FunctionType_UpdatePropertyTypeMap(type_ *FunctionType, propTypeMap j2gort.Ext_java_util_Map, cache *j2gort.Ext_java_util_LinkedHashSet) {
	clinit_FunctionType()
	_ = type_
	_ = propTypeMap
	_ = cache
	if type_ == nil {
		return
	}
	var prototype ObjectType = type_.GetPrototype()
	_ = prototype
	if prototype != nil && !prototype.IsJavaNil_() {
		var propNames j2gort.Ext_java_util_Set = prototype.GetOwnPropertyNames()
		_ = propNames
		{
			cur_43 := j2gort.NewIterCursor(j2gort.CheckNotNull(propNames))
			for cur_43.HasNext() {
				name := cur_43.NextStr()
				_ = name
				if !(j2gort.MapContainsKeyStr(propTypeMap, name)) {
					var propType JSType = prototype.GetObjectTypeBase().GetPropertyType_1_string(name)
					_ = propType
					j2gort.MapPutStr(propTypeMap, name, propType)
				}
			}
		}
	}
	var iterable []ObjectType = j2gort.ImmListToSliceNullable[ObjectType](type_.GetExtendedInterfaces())
	_ = iterable
	if iterable != nil {
		{

			for _, interfaceType := range iterable {

				_ = interfaceType
				var superConstructor *FunctionType = interfaceType.GetConstructor()
				_ = superConstructor
				if (superConstructor == nil) || cache.Contains(j2gort.PtrToAny(superConstructor)) {
					continue
				}
				cache.Add_1(j2gort.PtrToAny(superConstructor))
				FunctionType_UpdatePropertyTypeMap(superConstructor, propTypeMap, cache)
				cache.Remove_1_any(j2gort.PtrToAny(superConstructor))
			}
		}
	}
}

// Java method: FunctionType.checkExtendsLoop()
func (this_ *FunctionType) CheckExtendsLoop_0() j2gort.Ext_java_util_List {
	return this_.priv_FunctionType_checkExtendsLoop_9016(j2gort.NewExt_java_util_LinkedHashSet(), j2gort.NewExt_java_util_ArrayList())
}

// Java method: FunctionType.checkExtendsLoop(java.util.Set<com.google.javascript.rhino.jstype.FunctionType> cache, java.util.List<com.google.javascript.rhino.jstype.FunctionType> path)
func (this_ *FunctionType) priv_FunctionType_checkExtendsLoop_9016(cache j2gort.Ext_java_util_Set, path j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	_ = cache
	_ = path
	var iterable []ObjectType = j2gort.ImmListToSliceNullable[ObjectType](this_.GetExtendedInterfaces())
	_ = iterable
	if iterable != nil {
		{

			for _, interfaceType := range iterable {

				_ = interfaceType
				var superConstructor *FunctionType = interfaceType.GetConstructor()
				_ = superConstructor
				if superConstructor == nil {
					continue
				}
				if cache.Contains(j2gort.PtrToAny(superConstructor)) {
					path.Add_1(j2gort.PtrToAny(superConstructor))
					for !(j2gort.CastPtr[FunctionType](path.Get(int32(0))) == superConstructor) {
						path.Remove_1_int32(int32(0))
					}
					return path
				}
				cache.Add_1(j2gort.PtrToAny(superConstructor))
				path.Add_1(j2gort.PtrToAny(superConstructor))
				var result j2gort.Ext_java_util_List = superConstructor.priv_FunctionType_checkExtendsLoop_9016(cache, path)
				_ = result
				if result != nil {
					return result
				}
				cache.Remove_1_any(j2gort.PtrToAny(superConstructor))
				path.Remove_1_int32((path.Size() - int32(1)))
			}
		}
	}
	return nil
}

// Java method: FunctionType.acceptsArguments(java.util.List<? extends com.google.javascript.rhino.jstype.JSType> argumentTypes)
func (this_ *FunctionType) AcceptsArguments(argumentTypes j2gort.Ext_java_util_List) bool {
	_ = argumentTypes
	var arguments j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](argumentTypes.Iterator())
	_ = arguments
	var params []FunctionType_Parameter = this_.parameters
	_ = params
	var paramLen int32 = int32(len(params))
	_ = paramLen
	var paramIdx int32 = int32(0)
	_ = paramIdx
	var parameter FunctionType_Parameter = (FunctionType_Parameter{})
	_ = parameter
	var argument JSType = nil
	_ = argument
	for arguments.HasNext() && ((paramIdx < paramLen) || ((parameter != (FunctionType_Parameter{})) && parameter.IsVariadic())) {
		if paramIdx < paramLen {
			parameter = params[func() int32 { old_48 := paramIdx; nw_49 := old_48 + 1; paramIdx = nw_49; return old_48 }()]
		}
		argument = j2gort.AsIface[JSType](arguments.Next())
		if !(argument.GetJSTypeBase().IsSubtypeOf_1(parameter.GetJSType())) {
			return false
		}
	}
	var numArgs int32 = argumentTypes.Size()
	_ = numArgs
	return ((this_.GetMinArity() <= numArgs) && (numArgs <= this_.GetMaxArity()))
}

// Java method: FunctionType.forgetParameterAndReturnTypes()
func (this_ *FunctionType) ForgetParameterAndReturnTypes() *FunctionType {
	var fb_ FunctionType_Builder
	InitFunctionType_Builder(&fb_, this_.GetJSTypeBase().registry)
	(&fb_).WithName(this_.GetReferenceName())
	(&fb_).WithSourceNode(this_.GetJSTypeBase().subSourceNode)
	(&fb_).WithTypeOfThis(Upcast_JSType(this_.GetInstanceType()))
	(&fb_).WithKind(this_.GetKind())
	(&fb_).priv_Builder_withCanonicalRepresentation_9052(this_)
	var result *FunctionType = j2gort.CastPtr[FunctionType]((&fb_).Build())
	_ = result
	result.SetPrototypeBasedOn_1(this_.GetInstanceType())
	return result
}

// Java method: FunctionType.getConstructorOnlyTemplateParameters()
func (this_ *FunctionType) GetConstructorOnlyTemplateParameters() j2gort.Ext_com_google_common_collect_ImmutableList {
	if !(this_.IsConstructor()) {
		j2gort.CheckStateMsg(false, this_)
	}
	if this_.GetJSTypeBase().subSourceNode == 0 {
		return StaticGet_FunctionType_EMPTY_TEMPLATE_KEYS()
	}
	var map_ *TemplateTypeMap = this_.GetTemplateTypeMap()
	_ = map_
	var ctorOnlyKeyCount int32 = (this_.GetTemplateParamCount() - this_.GetInstanceType().GetTemplateParamCount())
	_ = ctorOnlyKeyCount
	var size int32 = map_.Size()
	_ = size
	return map_.GetTemplateKeysSlice((size - ctorOnlyKeyCount), size)
}

// Java method: FunctionType.isAmbiguousConstructor()
func (this_ *FunctionType) IsAmbiguousConstructor() bool {
	var ambiguity int8
	_ = ambiguity
	if this_.GetPrototypeObjectTypeBase().cold != nil {
		ambiguity = this_.GetPrototypeObjectTypeBase().cold.constructorAmbiguity
	} else {
		ambiguity = int8(0)
	}
	if int32(ambiguity) == int32(int8(0)) {
		ambiguity = this_.priv_FunctionType_calculateConstructorAmbiguityCode_9021()
		{
			r_50 := this_.EnsureCold()
			r_50.constructorAmbiguity = ambiguity
		}
	}
	return (int32(ambiguity) == int32(int8(1)))
}

// Java method: FunctionType.calculateConstructorAmbiguityCode()
func (this_ *FunctionType) priv_FunctionType_calculateConstructorAmbiguityCode_9021() int8 {
	if this_.IsUnknownType() {
		return int8(1)
	} else if this_.IsNativeObjectType() {
		return int8(2)
	} else {
		var superConstructor *FunctionType = this_.GetSuperClassConstructor()
		_ = superConstructor
		if superConstructor == nil {
			return int8(1)
		} else if superConstructor.IsAmbiguousConstructor() {
			return int8(1)
		} else if this_.GetJSTypeBase().subSourceNode != 0 {
			return int8(2)
		} else if this_.priv_FunctionType_isDelegateProxy_9022() {
			return int8(2)
		} else {
			return int8(1)
		}
	}
}

// Java method: FunctionType.isDelegateProxy()
func (this_ *FunctionType) priv_FunctionType_isDelegateProxy_9022() bool {
	return (this_.HasReferenceName() && javalang.StringEndsWith(j2gort.CheckNotNullString(this_.GetReferenceName()), j2gort.CheckNotNullString(StaticGet_FunctionType_DELEGATE_SUFFIX())))
}

// Java method: FunctionType.getClosurePrimitive()
func (this_ *FunctionType) GetClosurePrimitive() ClosurePrimitive {
	if this_.GetPrototypeObjectTypeBase().cold != nil {
		return this_.GetPrototypeObjectTypeBase().cold.closurePrimitive
	} else {
		return ClosurePrimitive(0)
	}
}

// Java static method: FunctionType.builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func FunctionType_S_Builder(registry *JSTypeRegistry) *FunctionType_Builder {
	clinit_FunctionType()
	_ = registry
	return NewFunctionType_Builder(registry)
}

// Java method: FunctionType.toBuilder()
func (this_ *FunctionType) ToBuilder() *FunctionType_Builder {
	builder := &FunctionType_Builder{}
	builder.BindSelf_(builder)
	this_.InitBuilder_(builder)
	return builder

}
func (this_ *FunctionType) InitBuilder_(builder *FunctionType_Builder) {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	InitFunctionType_Builder(builder, this_.GetJSTypeBase().registry)
	builder.SetGoogModuleId(this_.GetGoogModuleId())
	builder.className = this_.GetReferenceName()
	builder.nativeType = this_.IsNativeObjectType()
	builder.templateTypeMap = this_.GetTemplateTypeMap()
	builder.priv_Builder_withCanonicalRepresentation_9052(this_.GetCanonicalRepresentation())
	builder.SetIsKnownAmbiguous(((this_.GetPrototypeObjectTypeBase().cold != nil) && (int32(this_.GetPrototypeObjectTypeBase().cold.constructorAmbiguity) == int32(int8(1)))))
	builder.WithClosurePrimitiveId(this_.GetClosurePrimitive())
	builder.WithIsAbstract(this_.IsAbstract())
	builder.WithKind(this_.GetKind())
	builder.WithParameterArray(this_.parameters)
	builder.WithReturnType_2(this_.GetReturnType(), this_.IsReturnTypeInferred())
	builder.WithSourceNode(this_.GetSource())
	builder.WithTypeOfThis(this_.GetTypeOfThis())
	builder.templateParamCount = this_.GetTemplateParamCount()

	if this_.IsConstructor() {
		builder.WithConstructorTemplateKeys(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.GetConstructorOnlyTemplateParameters()))
	}
}

// Java method: FunctionType.getCanonicalRepresentation()
func (this_ *FunctionType) GetCanonicalRepresentation() *FunctionType {
	if this_.GetPrototypeObjectTypeBase().cold != nil {
		return this_.GetPrototypeObjectTypeBase().cold.canonicalRepresentation
	} else {
		return nil
	}
}

// Java enum: com.google.javascript.rhino.jstype.FunctionType.Kind
type FunctionType_Kind uint8

const (
	enum_FunctionType_Kind_ORDINARY    FunctionType_Kind = 1
	enum_FunctionType_Kind_CONSTRUCTOR FunctionType_Kind = 2
	enum_FunctionType_Kind_INTERFACE   FunctionType_Kind = 3
	enum_FunctionType_Kind_NONE        FunctionType_Kind = 4
)

var enumBaseTable_FunctionType_Kind = [...]j2gort.EnumBase{
	{Name_: "ORDINARY", Ordinal_: 0},
	{Name_: "CONSTRUCTOR", Ordinal_: 1},
	{Name_: "INTERFACE", Ordinal_: 2},
	{Name_: "NONE", Ordinal_: 3},
}
var enumValues_FunctionType_Kind = []FunctionType_Kind{1, 2, 3, 4}

func StaticGet_FunctionType_Kind_ORDINARY() FunctionType_Kind {
	return enum_FunctionType_Kind_ORDINARY
}
func StaticGet_FunctionType_Kind_CONSTRUCTOR() FunctionType_Kind {
	return enum_FunctionType_Kind_CONSTRUCTOR
}
func StaticGet_FunctionType_Kind_INTERFACE() FunctionType_Kind {
	return enum_FunctionType_Kind_INTERFACE
}
func StaticGet_FunctionType_Kind_NONE() FunctionType_Kind {
	return enum_FunctionType_Kind_NONE
}

func FunctionType_Kind_Values() []FunctionType_Kind {
	return j2gort.CloneSlice(enumValues_FunctionType_Kind)
}

func FunctionType_Kind_ValueOf(name string) FunctionType_Kind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_FunctionType_Kind {
		if enumBaseTable_FunctionType_Kind[i_].Name_ == name {
			return FunctionType_Kind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.FunctionType.Kind." + name))
}

func (this_ FunctionType_Kind) IsJavaNil_() bool                              { return this_ == 0 }
func (FunctionType_Kind) HasIdentityEquals_()                                 {}
func (FunctionType_Kind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (FunctionType_Kind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ FunctionType_Kind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (FunctionType_Kind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.FunctionType$Kind", "Kind"
}
func (this_ FunctionType_Kind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_FunctionType_Kind[this_-1]
}
func (this_ FunctionType_Kind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ FunctionType_Kind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ FunctionType_Kind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FunctionType_Kind[this_-1].Name_
}
func (this_ FunctionType_Kind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_FunctionType_Kind[this_-1].Name_
}
func (this_ FunctionType_Kind) Equals(other_ any) bool {
	o_, ok_ := other_.(FunctionType_Kind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ FunctionType_Kind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_FunctionType_Kind[this_-1].HashCode()
}
func (this_ FunctionType_Kind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FunctionType_Kind[this_-1].CompareTo(other_)
}
func (this_ FunctionType_Kind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_FunctionType_Kind[this_-1].Name_)
}
func init() {
	vals_ := []any{FunctionType_Kind(1), FunctionType_Kind(2), FunctionType_Kind(3), FunctionType_Kind(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.FunctionType.Kind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.FunctionType$Kind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.FunctionType.Kind", "com.google.javascript.rhino.jstype.FunctionType$Kind", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionType_Kind](v_)
	})
}

// Java enum: com.google.javascript.rhino.jstype.FunctionType.ConstructorAmbiguity
type FunctionType_ConstructorAmbiguity uint8

const (
	enum_FunctionType_ConstructorAmbiguity_UNKNOWN                    FunctionType_ConstructorAmbiguity = 1
	enum_FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR   FunctionType_ConstructorAmbiguity = 2
	enum_FunctionType_ConstructorAmbiguity_IS_UNAMBIGUOUS_CONSTRUCTOR FunctionType_ConstructorAmbiguity = 3
)

var enumBaseTable_FunctionType_ConstructorAmbiguity = [...]j2gort.EnumBase{
	{Name_: "UNKNOWN", Ordinal_: 0},
	{Name_: "IS_AMBIGUOUS_CONSTRUCTOR", Ordinal_: 1},
	{Name_: "IS_UNAMBIGUOUS_CONSTRUCTOR", Ordinal_: 2},
}
var enumValues_FunctionType_ConstructorAmbiguity = []FunctionType_ConstructorAmbiguity{1, 2, 3}

func StaticGet_FunctionType_ConstructorAmbiguity_UNKNOWN() FunctionType_ConstructorAmbiguity {
	return enum_FunctionType_ConstructorAmbiguity_UNKNOWN
}
func StaticGet_FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR() FunctionType_ConstructorAmbiguity {
	return enum_FunctionType_ConstructorAmbiguity_IS_AMBIGUOUS_CONSTRUCTOR
}
func StaticGet_FunctionType_ConstructorAmbiguity_IS_UNAMBIGUOUS_CONSTRUCTOR() FunctionType_ConstructorAmbiguity {
	return enum_FunctionType_ConstructorAmbiguity_IS_UNAMBIGUOUS_CONSTRUCTOR
}

func FunctionType_ConstructorAmbiguity_Values() []FunctionType_ConstructorAmbiguity {
	return j2gort.CloneSlice(enumValues_FunctionType_ConstructorAmbiguity)
}

func FunctionType_ConstructorAmbiguity_ValueOf(name string) FunctionType_ConstructorAmbiguity {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_FunctionType_ConstructorAmbiguity {
		if enumBaseTable_FunctionType_ConstructorAmbiguity[i_].Name_ == name {
			return FunctionType_ConstructorAmbiguity(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.FunctionType.ConstructorAmbiguity." + name))
}

func (this_ FunctionType_ConstructorAmbiguity) IsJavaNil_() bool                              { return this_ == 0 }
func (FunctionType_ConstructorAmbiguity) HasIdentityEquals_()                                 {}
func (FunctionType_ConstructorAmbiguity) Implements_j2gort_Ext_java_lang_Enum()               {}
func (FunctionType_ConstructorAmbiguity) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ FunctionType_ConstructorAmbiguity) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (FunctionType_ConstructorAmbiguity) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.FunctionType$ConstructorAmbiguity", "ConstructorAmbiguity"
}
func (this_ FunctionType_ConstructorAmbiguity) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_FunctionType_ConstructorAmbiguity[this_-1]
}
func (this_ FunctionType_ConstructorAmbiguity) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ FunctionType_ConstructorAmbiguity) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ FunctionType_ConstructorAmbiguity) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FunctionType_ConstructorAmbiguity[this_-1].Name_
}
func (this_ FunctionType_ConstructorAmbiguity) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_FunctionType_ConstructorAmbiguity[this_-1].Name_
}
func (this_ FunctionType_ConstructorAmbiguity) Equals(other_ any) bool {
	o_, ok_ := other_.(FunctionType_ConstructorAmbiguity)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ FunctionType_ConstructorAmbiguity) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_FunctionType_ConstructorAmbiguity[this_-1].HashCode()
}
func (this_ FunctionType_ConstructorAmbiguity) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FunctionType_ConstructorAmbiguity[this_-1].CompareTo(other_)
}
func (this_ FunctionType_ConstructorAmbiguity) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_FunctionType_ConstructorAmbiguity[this_-1].Name_)
}
func init() {
	vals_ := []any{FunctionType_ConstructorAmbiguity(1), FunctionType_ConstructorAmbiguity(2), FunctionType_ConstructorAmbiguity(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.FunctionType.ConstructorAmbiguity", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.FunctionType$ConstructorAmbiguity", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.FunctionType.ConstructorAmbiguity", "com.google.javascript.rhino.jstype.FunctionType$ConstructorAmbiguity", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionType_ConstructorAmbiguity](v_)
	})
}

// Java class: com.google.javascript.rhino.jstype.FunctionType.Builder
type FunctionType_Builder struct {
	PrototypeObjectType_BuilderBase
	googModuleId            string                                            `gson:"googModuleId"`
	parameters              j2gort.Ext_java_util_List                         `gson:"parameters"`
	parameterArray          []FunctionType_Parameter                          `gson:"parameterArray"`
	returnType              JSType                                            `gson:"returnType"`
	typeOfThis              JSType                                            `gson:"typeOfThis"`
	setPrototypeBasedOn     ObjectType                                        `gson:"setPrototypeBasedOn"`
	constructorOnlyKeys     j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"constructorOnlyKeys"`
	canonicalRepresentation *FunctionType                                     `gson:"canonicalRepresentation"`
	sourceNode              Node                                              `gson:"sourceNode"`
	kind                    FunctionType_Kind                                 `gson:"kind"`
	isAbstract              bool                                              `gson:"isAbstract"`
	isKnownAmbiguous        bool                                              `gson:"isKnownAmbiguous"`
	returnTypeIsInferred    bool                                              `gson:"returnTypeIsInferred"`
	returnsOwnInstanceType  bool                                              `gson:"returnsOwnInstanceType"`
	noTypeKind              int8                                              `gson:"noTypeKind"`
	primitiveId             ClosurePrimitive                                  `gson:"primitiveId"`
}

func (this_ *FunctionType_Builder) BindSelf_(s any) {
	_ = s
	this_.PrototypeObjectType_BuilderBase.BindSelf_(s)
}
func (this_ *FunctionType_Builder) IsJavaNil_() bool { return this_ == nil }
func (*FunctionType_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.FunctionType$Builder", "Builder"
}
func (this_ *FunctionType_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PrototypeObjectType_BuilderBase.IdentityHash_()
}
func (*FunctionType_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.FunctionType.Builder", "com.google.javascript.rhino.jstype.FunctionType$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionType_Builder](v_)
	})
}
func (*FunctionType_Builder) Implements_rhino_PrototypeObjectType_Builder() {}

func (this_ *FunctionType_Builder) initFields_FunctionType_Builder() {
	_ = this_
	this_.sourceNode = Node(0)
	this_.googModuleId = j2gort.NullString
	this_.parameters = nil
	this_.parameterArray = nil
	this_.returnType = nil
	this_.typeOfThis = nil
	this_.setPrototypeBasedOn = nil
	this_.constructorOnlyKeys = StaticGet_FunctionType_EMPTY_CTOR_KEYS()
	this_.kind = StaticGet_FunctionType_Kind_ORDINARY()
	this_.isKnownAmbiguous = false
	this_.noTypeKind = int8(0)
	this_.primitiveId = ClosurePrimitive(0)
	this_.canonicalRepresentation = nil
}

// Java constructor: Builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewFunctionType_Builder(registry *JSTypeRegistry) *FunctionType_Builder {
	this_ := &FunctionType_Builder{}
	this_.BindSelf_(this_)
	InitFunctionType_Builder(this_, registry)
	return this_
}

func InitFunctionType_Builder(this_ *FunctionType_Builder, registry *JSTypeRegistry) {
	_ = registry
	InitPrototypeObjectType_Builder(&this_.PrototypeObjectType_BuilderBase, registry)
	this_.initFields_FunctionType_Builder()
	this_.implicitPrototype = j2gort.CheckNotNullIface(registry.GetNativeObjectType(StaticGet_JSTypeNative_FUNCTION_PROTOTYPE()))
}

// Java method: Builder.reset()
func (this_ *FunctionType_Builder) Reset() *FunctionType_Builder {
	this_.ResetBase()
	this_.SetImplicitPrototype(j2gort.CheckNotNullIface(this_.GetPrototypeObjectType_BuilderBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_FUNCTION_PROTOTYPE())))
	this_.sourceNode = Node(0)
	this_.googModuleId = j2gort.NullString
	this_.parameters = nil
	this_.parameterArray = nil
	this_.returnType = nil
	this_.typeOfThis = nil
	this_.setPrototypeBasedOn = nil
	this_.constructorOnlyKeys = StaticGet_FunctionType_EMPTY_CTOR_KEYS()
	this_.kind = StaticGet_FunctionType_Kind_ORDINARY()
	this_.isAbstract = false
	this_.isKnownAmbiguous = false
	this_.returnTypeIsInferred = false
	this_.returnsOwnInstanceType = false
	this_.noTypeKind = int8(0)
	this_.primitiveId = ClosurePrimitive(0)
	this_.canonicalRepresentation = nil
	return this_
}

// Java method: Builder.withName(java.lang.String name)
func (this_ *FunctionType_Builder) WithName(name string) *FunctionType_Builder {
	_ = name
	this_.className = name
	return this_
}

// Java method: Builder.withSourceNode(com.google.javascript.rhino.Node sourceNode)
func (this_ *FunctionType_Builder) WithSourceNode(sourceNode Node) *FunctionType_Builder {
	_ = sourceNode
	this_.sourceNode = sourceNode
	return this_
}

// Java method: Builder.setGoogModuleId(java.lang.String x)
func (this_ *FunctionType_Builder) SetGoogModuleId(x string) *FunctionType_Builder {
	_ = x
	this_.googModuleId = x
	return this_
}

// Java method: Builder.withParameters(java.util.List<com.google.javascript.rhino.jstype.FunctionType.Parameter> parameters)
func (this_ *FunctionType_Builder) WithParameters_1(parameters j2gort.Ext_java_util_List) *FunctionType_Builder {
	_ = parameters
	this_.parameters = parameters
	if (parameters != nil) && parameters.IsEmpty() {
		this_.parameterArray = StaticGet_FunctionType_EMPTY_PARAMETERS()
	} else {
		this_.parameterArray = nil
	}
	return this_
}

// Java method: Builder.withParameterArray(com.google.javascript.rhino.jstype.FunctionType.Parameter @org.jspecify.annotations.Nullable [] parameters)
func (this_ *FunctionType_Builder) WithParameterArray(parameters []FunctionType_Parameter) *FunctionType_Builder {
	_ = parameters
	this_.parameters = nil
	if parameters == nil {
		this_.parameterArray = nil
	} else if int32(len(parameters)) == int32(0) {
		this_.parameterArray = StaticGet_FunctionType_EMPTY_PARAMETERS()
	} else {
		this_.parameterArray = parameters
	}
	return this_
}

// Java method: Builder.withParameters()
func (this_ *FunctionType_Builder) WithParameters_0() *FunctionType_Builder {
	this_.parameters = nil
	this_.parameterArray = StaticGet_FunctionType_EMPTY_PARAMETERS()
	return this_
}

// Java method: Builder.withReturnType(com.google.javascript.rhino.jstype.JSType returnType)
func (this_ *FunctionType_Builder) WithReturnType_1(returnType JSType) *FunctionType_Builder {
	_ = returnType
	this_.returnType = returnType
	return this_
}

// Java method: Builder.withReturnType(com.google.javascript.rhino.jstype.JSType returnType, boolean inferred)
func (this_ *FunctionType_Builder) WithReturnType_2(returnType JSType, inferred bool) *FunctionType_Builder {
	_ = returnType
	_ = inferred
	this_.returnType = returnType
	this_.returnTypeIsInferred = inferred
	return this_
}

// Java method: Builder.withReturnsOwnInstanceType()
func (this_ *FunctionType_Builder) WithReturnsOwnInstanceType() *FunctionType_Builder {
	this_.returnsOwnInstanceType = true
	return this_
}

// Java method: Builder.withInferredReturnType(com.google.javascript.rhino.jstype.JSType returnType)
func (this_ *FunctionType_Builder) WithInferredReturnType(returnType JSType) *FunctionType_Builder {
	_ = returnType
	this_.returnType = returnType
	this_.returnTypeIsInferred = true
	return this_
}

// Java method: Builder.withTypeOfThis(com.google.javascript.rhino.jstype.JSType typeOfThis)
func (this_ *FunctionType_Builder) WithTypeOfThis(typeOfThis JSType) *FunctionType_Builder {
	_ = typeOfThis
	this_.typeOfThis = typeOfThis
	return this_
}

// Java method: Builder.withTemplateKeys(com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.TemplateType> templateKeys)
func (this_ *FunctionType_Builder) WithTemplateKeys_1_Ext_com_google_common_collect_ImmutableList(templateKeys j2gort.Ext_com_google_common_collect_ImmutableList) *FunctionType_Builder {
	_ = templateKeys
	if (templateKeys == nil) || templateKeys.IsEmpty() {
		this_.templateTypeMap = this_.GetPrototypeObjectType_BuilderBase().registry.GetEmptyTemplateTypeMap()
		this_.templateParamCount = int32(0)
		return this_
	}
	this_.templateTypeMap = this_.GetPrototypeObjectType_BuilderBase().registry.GetEmptyTemplateTypeMap().CopyWithExtension_2(templateKeys, StaticGet_FunctionType_EMPTY_TEMPLATE_VALUES())
	this_.templateParamCount = templateKeys.Size()
	return this_
}

// Java method: Builder.withTemplateKeysArray(com.google.javascript.rhino.jstype.TemplateType[] templateKeys)
func (this_ *FunctionType_Builder) WithTemplateKeysArray(templateKeys []*TemplateType) *FunctionType_Builder {
	_ = templateKeys
	if (templateKeys == nil) || (int32(len(templateKeys)) == int32(0)) {
		this_.templateTypeMap = this_.GetPrototypeObjectType_BuilderBase().registry.GetEmptyTemplateTypeMap()
		this_.templateParamCount = int32(0)
		return this_
	}
	this_.templateTypeMap = this_.GetPrototypeObjectType_BuilderBase().registry.GetEmptyTemplateTypeMap().CopyWithExtensionKeysOnly(templateKeys)
	this_.templateParamCount = int32(len(templateKeys))
	return this_
}

// Java method: Builder.withTemplateKeysFrom(com.google.javascript.rhino.jstype.TemplateTypeMap sourceMap)
func (this_ *FunctionType_Builder) WithTemplateKeysFrom(sourceMap *TemplateTypeMap) *FunctionType_Builder {
	_ = sourceMap
	var count int32 = sourceMap.Size()
	_ = count
	if count == int32(0) {
		this_.templateTypeMap = this_.GetPrototypeObjectType_BuilderBase().registry.GetEmptyTemplateTypeMap()
		this_.templateParamCount = int32(0)
		return this_
	}
	this_.templateTypeMap = sourceMap.AsKeysOnlySingleSubMap()
	this_.templateParamCount = count
	return this_
}

// Java method: Builder.withTemplateKeys(com.google.javascript.rhino.jstype.TemplateType[] templateKeys)
func (this_ *FunctionType_Builder) WithTemplateKeys_1_TemplateType(templateKeys []*TemplateType) *FunctionType_Builder {
	_ = templateKeys
	return this_.WithTemplateKeysArray(templateKeys)
}

// Java method: Builder.withConstructorTemplateKeys(java.lang.Iterable<com.google.javascript.rhino.jstype.TemplateType> constructorOnlyKeys)
func (this_ *FunctionType_Builder) WithConstructorTemplateKeys(constructorOnlyKeys j2gort.Ext_java_lang_Iterable) *FunctionType_Builder {
	_ = constructorOnlyKeys
	this_.constructorOnlyKeys = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_lang_Iterable(constructorOnlyKeys)
	return this_
}

// Java method: Builder.withKind(com.google.javascript.rhino.jstype.FunctionType.Kind kind)
func (this_ *FunctionType_Builder) WithKind(kind FunctionType_Kind) *FunctionType_Builder {
	_ = kind
	this_.kind = kind
	return this_
}

// Java method: Builder.forConstructor()
func (this_ *FunctionType_Builder) ForConstructor() *FunctionType_Builder {
	this_.kind = StaticGet_FunctionType_Kind_CONSTRUCTOR()
	return this_
}

// Java method: Builder.forInterface()
func (this_ *FunctionType_Builder) ForInterface() *FunctionType_Builder {
	this_.kind = StaticGet_FunctionType_Kind_INTERFACE()
	this_.parameters = nil
	this_.parameterArray = StaticGet_FunctionType_EMPTY_PARAMETERS()
	return this_
}

// Java method: Builder.forNativeType()
func (this_ *FunctionType_Builder) ForNativeType() *FunctionType_Builder {
	this_.nativeType = true
	return this_
}

// Java method: Builder.withIsAbstract(boolean isAbstract)
func (this_ *FunctionType_Builder) WithIsAbstract(isAbstract bool) *FunctionType_Builder {
	_ = isAbstract
	this_.isAbstract = isAbstract
	return this_
}

// Java method: Builder.setIsKnownAmbiguous(boolean x)
func (this_ *FunctionType_Builder) SetIsKnownAmbiguous(x bool) *FunctionType_Builder {
	_ = x
	this_.isKnownAmbiguous = x
	return this_
}

// Java method: Builder.withPrototypeBasedOn(com.google.javascript.rhino.jstype.ObjectType setPrototypeBasedOn)
func (this_ *FunctionType_Builder) WithPrototypeBasedOn(setPrototypeBasedOn ObjectType) *FunctionType_Builder {
	_ = setPrototypeBasedOn
	this_.setPrototypeBasedOn = setPrototypeBasedOn
	return this_
}

// Java method: Builder.withClosurePrimitiveId(com.google.javascript.rhino.ClosurePrimitive id)
func (this_ *FunctionType_Builder) WithClosurePrimitiveId(id ClosurePrimitive) *FunctionType_Builder {
	_ = id
	this_.primitiveId = id
	return this_
}

// Java method: Builder.withNoTypeKind(byte noTypeKind)
func (this_ *FunctionType_Builder) priv_Builder_withNoTypeKind_9051(noTypeKind int8) *FunctionType_Builder {
	_ = noTypeKind
	this_.noTypeKind = noTypeKind
	return this_
}

// Java method: Builder.withCanonicalRepresentation(com.google.javascript.rhino.jstype.FunctionType representation)
func (this_ *FunctionType_Builder) priv_Builder_withCanonicalRepresentation_9052(representation *FunctionType) *FunctionType_Builder {
	_ = representation
	this_.canonicalRepresentation = representation
	return this_
}

// Java method: Builder.build()
func (this_ *FunctionType_Builder) Build() any {
	switch (this_.kind).Ordinal() {
	case int32(1) /* Kind.CONSTRUCTOR */, int32(2) /* Kind.INTERFACE */ :
	case int32(3) /* Kind.NONE */ :
		if !(this_.returnsOwnInstanceType) {
			j2gort.CheckState(false)
		}
	case int32(0) /* Kind.ORDINARY */ :
		if !(!(this_.returnsOwnInstanceType)) {
			j2gort.CheckState(false)
		}
	}
	if this_.returnsOwnInstanceType {
		if this_.returnType != nil && !this_.returnType.IsJavaNil_() {
			j2gort.CheckState(false)
		}
		if this_.typeOfThis != nil && !this_.typeOfThis.IsJavaNil_() {
			j2gort.CheckState(false)
		}
	}
	switch (this_.kind).Ordinal() {
	case int32(1) /* Kind.CONSTRUCTOR */, int32(2) /* Kind.INTERFACE */ :
	case int32(3) /* Kind.NONE */, int32(0) /* Kind.ORDINARY */ :
		if !(this_.constructorOnlyKeys.IsEmpty()) {
			j2gort.CheckState(false)
		}
	}
	return j2gort.PtrToAny(allocPooled_FunctionType(&this_.registry.

		// Java method: Builder.buildAndResolve()
		slab_FunctionType, this_))
}

func (this_ *FunctionType_Builder) BuildAndResolve() *FunctionType {
	return j2gort.CheckNotNullPtr(j2gort.CastPtr[FunctionType](this_.Build()).ToMaybeFunctionType())
}

// Java class: com.google.javascript.rhino.jstype.FunctionType.Parameter
type FunctionType_Parameter struct {
	jsType   JSType `gson:"jsType"`
	optional bool   `gson:"optional"`
	variadic bool   `gson:"variadic"`
}

func (this_ FunctionType_Parameter) BindSelf_(s any) {
	_ = s
}
func (this_ FunctionType_Parameter) IsJavaNil_() bool { return this_ == (FunctionType_Parameter{}) }
func (FunctionType_Parameter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.FunctionType$Parameter", "Parameter"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.FunctionType.Parameter", "com.google.javascript.rhino.jstype.FunctionType$Parameter", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionType_Parameter](v_)
	})
}

func (this_ *FunctionType_Parameter) initFields_FunctionType_Parameter() {
	_ = this_
}

// Java constructor: Parameter(com.google.javascript.rhino.jstype.JSType jsType, boolean optional, boolean variadic)
func NewFunctionType_Parameter(jsType JSType, optional bool, variadic bool) FunctionType_Parameter {
	this_ := FunctionType_Parameter{}
	InitFunctionType_Parameter(&this_, jsType, optional, variadic)
	return this_
}

func InitFunctionType_Parameter(this_ *FunctionType_Parameter, jsType JSType, optional bool, variadic bool) {
	_ = jsType
	_ = optional
	_ = variadic
	this_.initFields_FunctionType_Parameter()
	if jsType == nil || jsType.IsJavaNil_() {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null JSType")))
	}
	this_.jsType = jsType
	this_.optional = optional
	this_.variadic = variadic
}

// Java static method: Parameter.create(com.google.javascript.rhino.jstype.JSType type, boolean isOptional, boolean isVariadic)
func FunctionType_Parameter_Create(type_ JSType, isOptional bool, isVariadic bool) FunctionType_Parameter {
	_ = type_
	_ = isOptional
	_ = isVariadic
	return NewFunctionType_Parameter(type_, isOptional, isVariadic)
}

// Java method: Parameter.getJSType()
func (this_ FunctionType_Parameter) GetJSType() JSType {
	if this_ == (FunctionType_Parameter{}) {
		j2gort.ThrowNPE()
	}
	return this_.jsType
}

// Java method: Parameter.isOptional()
func (this_ FunctionType_Parameter) IsOptional() bool {
	if this_ == (FunctionType_Parameter{}) {
		j2gort.ThrowNPE()
	}
	return this_.optional
}

// Java method: Parameter.isVariadic()
func (this_ FunctionType_Parameter) IsVariadic() bool {
	if this_ == (FunctionType_Parameter{}) {
		j2gort.ThrowNPE()
	}
	return this_.variadic
}

// Java method: Parameter.toString()
func (this_ FunctionType_Parameter) String() string {
	if this_ == (FunctionType_Parameter{}) {
		j2gort.ThrowNPE()
	}
	return ("Parameter{JSType=" + j2gort.StrVal(this_.jsType) + ", optional=" + j2gort.StrValBool(this_.optional) + ", variadic=" + j2gort.StrValBool(this_.variadic) + "}")
}

// Java method: Parameter.equals(java.lang.Object o)
func (this_ FunctionType_Parameter) Equals(arg_51 any) bool {
	if this_ == (FunctionType_Parameter{}) {
		j2gort.ThrowNPE()
	}
	o := arg_51
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	var that FunctionType_Parameter
	_ = that
	var ok_52 bool
	that, ok_52 = j2gort.TryCast[FunctionType_Parameter](o)
	if ok_52 {
		return ((this_.jsType.GetJSTypeBase().EqualsJSType(that.jsType) && (this_.optional == that.optional)) && (this_.variadic == that.variadic))
	}
	return false
}

// Java method: Parameter.hashCode()
func (this_ FunctionType_Parameter) HashCode() int32 {
	if this_ == (FunctionType_Parameter{}) {
		j2gort.ThrowNPE()
	}
	var h int32 = int32(1)
	_ = h
	h = (h * int32(1000003))
	h = (h ^ this_.jsType.GetJSTypeBase().HashCode())
	h = (h * int32(1000003))
	h = (h ^ func() int32 {
		if this_.optional {
			return int32(1231)
		}
		return int32(1237)
	}())
	h = (h * int32(1000003))
	h = (h ^ func() int32 {
		if this_.variadic {
			return int32(1231)
		}
		return int32(1237)
	}())
	return h
}

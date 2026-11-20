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
// Port of com/google/javascript/rhino/jstype/FunctionParamBuilder.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.FunctionParamBuilder
type FunctionParamBuilder struct {
	registry                 *JSTypeRegistry                `gson:"registry"`
	parameters               j2gort.Ext_java_util_ArrayList `gson:"parameters"`
	initialCapacity          int32                          `gson:"initialCapacity"`
	pureUnknownRequiredCount int32                          `gson:"pureUnknownRequiredCount"`
	idHash_                  atomic.Int32
}

func (this_ *FunctionParamBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionParamBuilder) IsJavaNil_() bool { return this_ == nil }
func (*FunctionParamBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.FunctionParamBuilder", "FunctionParamBuilder"
}
func (this_ *FunctionParamBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionParamBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.FunctionParamBuilder", "com.google.javascript.rhino.jstype.FunctionParamBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionParamBuilder](v_)
	})
}

var static_FunctionParamBuilder_EMPTY_PARAMETER_LIST j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_FunctionParamBuilder_EMPTY_PARAMETER_LIST() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_FunctionParamBuilder()
	return static_FunctionParamBuilder_EMPTY_PARAMETER_LIST
}
func StaticSet_FunctionParamBuilder_EMPTY_PARAMETER_LIST(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_FunctionParamBuilder()
	static_FunctionParamBuilder_EMPTY_PARAMETER_LIST = v
}

var clinitDone_FunctionParamBuilder uint32
var clinitState_FunctionParamBuilder j2gort.ClinitState

func clinit_FunctionParamBuilder() {
	if atomic.LoadUint32(&clinitDone_FunctionParamBuilder) != 0 {
		return
	}
	clinitSlow_FunctionParamBuilder()
}

func clinitSlow_FunctionParamBuilder() {
	if !j2gort.EnterClinit(&clinitState_FunctionParamBuilder, &clinitDone_FunctionParamBuilder) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FunctionParamBuilder, &clinitDone_FunctionParamBuilder)
	static_FunctionParamBuilder_EMPTY_PARAMETER_LIST = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
}

func (this_ *FunctionParamBuilder) initFields_FunctionParamBuilder() {
	_ = this_
}

// Java constructor: FunctionParamBuilder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewFunctionParamBuilder_1(registry *JSTypeRegistry) *FunctionParamBuilder {
	this_ := &FunctionParamBuilder{}
	InitFunctionParamBuilder_1(this_, registry)
	return this_
}

func InitFunctionParamBuilder_1(this_ *FunctionParamBuilder, registry *JSTypeRegistry) {
	_ = registry
	this_.initFields_FunctionParamBuilder()
	this_.registry = registry
	this_.parameters = nil
	this_.initialCapacity = int32(4)
	this_.pureUnknownRequiredCount = int32(0)
}

// Java constructor: FunctionParamBuilder(com.google.javascript.rhino.jstype.JSTypeRegistry registry, int initialParameterCapacity)
func NewFunctionParamBuilder_2(registry *JSTypeRegistry, initialParameterCapacity int32) *FunctionParamBuilder {
	this_ := &FunctionParamBuilder{}
	InitFunctionParamBuilder_2(this_, registry, initialParameterCapacity)
	return this_
}

func InitFunctionParamBuilder_2(this_ *FunctionParamBuilder, registry *JSTypeRegistry, initialParameterCapacity int32) {
	_ = registry
	_ = initialParameterCapacity
	this_.initFields_FunctionParamBuilder()
	this_.registry = registry
	this_.parameters = nil
	this_.initialCapacity = initialParameterCapacity
	this_.pureUnknownRequiredCount = int32(0)
}

// Java method: FunctionParamBuilder.reset()
func (this_ *FunctionParamBuilder) Reset() {
	this_.pureUnknownRequiredCount = int32(0)
	if this_.parameters != nil {
		this_.parameters.Clear()
	}
}

// Java method: FunctionParamBuilder.flushPureUnknownRequired()
func (this_ *FunctionParamBuilder) priv_FunctionParamBuilder_flushPureUnknownRequired_8936() {
	var count int32 = this_.pureUnknownRequiredCount
	_ = count
	if count > int32(0) {
		this_.pureUnknownRequiredCount = int32(0)
		var list j2gort.Ext_java_util_ArrayList = this_.priv_FunctionParamBuilder_ensureParameters_8937()
		_ = list
		var p FunctionType_Parameter = this_.registry.GetUnknownRequiredParameter()
		_ = p
		{
			var i int32 = int32(0)
			_ = i
			for ; i < count; i = i + 1 {
				list.Add_1(j2gort.BoxNullableVal(p))
			}
		}
	}
}

// Java method: FunctionParamBuilder.ensureParameters()
func (this_ *FunctionParamBuilder) priv_FunctionParamBuilder_ensureParameters_8937() j2gort.Ext_java_util_ArrayList {
	var list j2gort.Ext_java_util_ArrayList = this_.parameters
	_ = list
	if list == nil {
		var iifeArg_1 int32
	iifeSwitch_1:
		switch {
		default:

			if this_.initialCapacity > int32(0) {
				iifeArg_1 = this_.initialCapacity
				break iifeSwitch_1
			}
			iifeArg_1 = int32(4)
			break iifeSwitch_1
		}
		list = j2gort.NewExt_java_util_ArrayList_1_int32(iifeArg_1)
		this_.parameters = list
	}
	return list
}

// Java method: FunctionParamBuilder.addRequiredParam(com.google.javascript.rhino.jstype.JSType type)
func (this_ *FunctionParamBuilder) AddRequiredParam(type_ JSType) bool {
	_ = type_
	var list j2gort.Ext_java_util_ArrayList = this_.parameters
	_ = list
	if (list == nil) || list.IsEmpty() {
		if j2gort.IfaceEquals[JSType](type_, this_.registry.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())) {
			this_.pureUnknownRequiredCount = this_.pureUnknownRequiredCount + 1
			return true
		}
		this_.priv_FunctionParamBuilder_flushPureUnknownRequired_8936()
		this_.priv_FunctionParamBuilder_newParameter_8945(type_, false, false)
		return true
	}
	if this_.priv_FunctionParamBuilder_hasOptionalOrVarArgs_8946() {
		return false
	}
	this_.priv_FunctionParamBuilder_newParameter_8945(type_, false, false)
	return true
}

// Java method: FunctionParamBuilder.addRequiredParams(com.google.javascript.rhino.jstype.JSType[] types)
func (this_ *FunctionParamBuilder) AddRequiredParams(types []JSType) bool {
	_ = types
	if this_.priv_FunctionParamBuilder_hasOptionalOrVarArgs_8946() {
		return false
	}
	for _, el_1 := range j2gort.CheckNotNullSlice(types) {
		type_ := el_1
		_ = type_
		this_.AddRequiredParam(type_)
	}
	return true
}

// Java method: FunctionParamBuilder.addOptionalParam(com.google.javascript.rhino.jstype.JSType type)
func (this_ *FunctionParamBuilder) AddOptionalParam(type_ JSType) bool {
	_ = type_
	if this_.HasVarArgs() {
		return false
	}
	this_.priv_FunctionParamBuilder_flushPureUnknownRequired_8936()
	this_.priv_FunctionParamBuilder_newParameter_8945(this_.registry.CreateOptionalType(type_), true, false)
	return true
}

// Java method: FunctionParamBuilder.addOptionalParams(com.google.javascript.rhino.jstype.JSType[] types)
func (this_ *FunctionParamBuilder) AddOptionalParams(types []JSType) bool {
	_ = types
	if this_.HasVarArgs() {
		return false
	}
	this_.priv_FunctionParamBuilder_flushPureUnknownRequired_8936()
	for _, el_2 := range j2gort.CheckNotNullSlice(types) {
		type_ := el_2
		_ = type_
		this_.priv_FunctionParamBuilder_newParameter_8945(this_.registry.CreateOptionalType(type_), true, false)
	}
	return true
}

// Java method: FunctionParamBuilder.addVarArgs(com.google.javascript.rhino.jstype.JSType type)
func (this_ *FunctionParamBuilder) AddVarArgs(type_ JSType) bool {
	_ = type_
	if this_.HasVarArgs() {
		return false
	}
	this_.priv_FunctionParamBuilder_flushPureUnknownRequired_8936()
	this_.priv_FunctionParamBuilder_newParameter_8945(type_, false, true)
	return true
}

// Java method: FunctionParamBuilder.newParameterFrom(com.google.javascript.rhino.jstype.FunctionType.Parameter n)
func (this_ *FunctionParamBuilder) NewParameterFrom(n FunctionType_Parameter) {
	_ = n
	this_.priv_FunctionParamBuilder_flushPureUnknownRequired_8936()
	this_.priv_FunctionParamBuilder_ensureParameters_8937().Add_1(j2gort.BoxNullableVal(n))
}

// Java method: FunctionParamBuilder.newOptionalParameterFrom(com.google.javascript.rhino.jstype.FunctionType.Parameter p)
func (this_ *FunctionParamBuilder) NewOptionalParameterFrom(p FunctionType_Parameter) {
	_ = p
	this_.priv_FunctionParamBuilder_flushPureUnknownRequired_8936()
	var isOptional bool = (p.IsOptional() || !(p.IsVariadic()))
	_ = isOptional
	if isOptional != p.IsOptional() {
		this_.priv_FunctionParamBuilder_newParameter_8945(p.GetJSType(), isOptional, p.IsVariadic())
	} else {
		this_.priv_FunctionParamBuilder_ensureParameters_8937().Add_1(j2gort.BoxNullableVal(p))
	}
}

// Java method: FunctionParamBuilder.newParameter(com.google.javascript.rhino.jstype.JSType type, boolean isOptional, boolean isVariadic)
func (this_ *FunctionParamBuilder) priv_FunctionParamBuilder_newParameter_8945(type_ JSType, isOptional bool, isVariadic bool) {
	_ = type_
	_ = isOptional
	_ = isVariadic
	var parameter FunctionType_Parameter
	_ = parameter
	if (!(isOptional) && !(isVariadic)) && j2gort.IfaceEquals[JSType](type_, this_.registry.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())) {
		parameter = this_.registry.GetUnknownRequiredParameter()
	} else {
		parameter = FunctionType_Parameter_Create(type_, isOptional, isVariadic)
	}
	this_.priv_FunctionParamBuilder_ensureParameters_8937().Add_1(j2gort.BoxNullableVal(parameter))
}

// Java method: FunctionParamBuilder.hasOptionalOrVarArgs()
func (this_ *FunctionParamBuilder) priv_FunctionParamBuilder_hasOptionalOrVarArgs_8946() bool {
	if this_.pureUnknownRequiredCount > int32(0) {
		return false
	}
	var list j2gort.Ext_java_util_ArrayList = this_.parameters
	_ = list
	if (list == nil) || list.IsEmpty() {
		return false
	}
	var lastParam FunctionType_Parameter = j2gort.Cast[FunctionType_Parameter](list.Get((list.Size() - int32(1))))
	_ = lastParam
	return (lastParam.IsOptional() || lastParam.IsVariadic())
}

// Java method: FunctionParamBuilder.hasVarArgs()
func (this_ *FunctionParamBuilder) HasVarArgs() bool {
	if this_.pureUnknownRequiredCount > int32(0) {
		return false
	}
	var list j2gort.Ext_java_util_ArrayList = this_.parameters
	_ = list
	return (((list != nil) && !(list.IsEmpty())) && j2gort.Cast[FunctionType_Parameter](list.Get((list.Size() - int32(1)))).IsVariadic())
}

// Java method: FunctionParamBuilder.buildArray()
func (this_ *FunctionParamBuilder) BuildArray() []FunctionType_Parameter {
	if this_.pureUnknownRequiredCount > int32(0) {
		return this_.registry.GetUnknownRequiredParameterArray(this_.pureUnknownRequiredCount)
	}
	var list j2gort.Ext_java_util_ArrayList = this_.parameters
	_ = list
	if (list == nil) || list.IsEmpty() {
		return StaticGet_FunctionType_EMPTY_PARAMETERS()
	}
	var n int32 = list.Size()
	_ = n
	var arr []FunctionType_Parameter = make([]FunctionType_Parameter, n)
	_ = arr
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			{
				a_3 := arr
				i_4 := i
				a_3[i_4] = j2gort.Cast[FunctionType_Parameter](list.Get(i))
			}
		}
	}
	return arr
}

// Java method: FunctionParamBuilder.build()
func (this_ *FunctionParamBuilder) Build() j2gort.Ext_com_google_common_collect_ImmutableList {
	this_.priv_FunctionParamBuilder_flushPureUnknownRequired_8936()
	var list j2gort.Ext_java_util_ArrayList = this_.parameters
	_ = list
	if (list == nil) || list.IsEmpty() {
		return StaticGet_FunctionParamBuilder_EMPTY_PARAMETER_LIST()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](list))
	}
}

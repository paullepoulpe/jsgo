/*
 * Copyright 2006 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Port of com/google/javascript/jscomp/JsMessage.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.JsMessage
type JsMessage struct {
	getSourceName                       string                                            `gson:"getSourceName"`
	getKey                              string                                            `gson:"getKey"`
	getId                               string                                            `gson:"getId"`
	getParts                            []JsMessage_Part                                  `gson:"getParts"`
	getGenderedMessagesMap              j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"getGenderedMessagesMap"`
	getAlternateId                      string                                            `gson:"getAlternateId"`
	getDesc                             string                                            `gson:"getDesc"`
	getMeaning                          string                                            `gson:"getMeaning"`
	getPlaceholderNameToExampleMap      j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"getPlaceholderNameToExampleMap"`
	getPlaceholderNameToOriginalCodeMap j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"getPlaceholderNameToOriginalCodeMap"`
	jsPlaceholderNames                  j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"jsPlaceholderNames"`
	canonicalPlaceholderNames           j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"canonicalPlaceholderNames"`
	isAnonymous                         bool                                              `gson:"isAnonymous"`
	isExternal                          bool                                              `gson:"isExternal"`
}

func (this_ *JsMessage) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessage) IsJavaNil_() bool { return this_ == nil }
func (*JsMessage) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessage", "JsMessage"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessage", "com.google.javascript.jscomp.JsMessage", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessage](v_)
	})
}
func (*JsMessage) Implements_j2gort_Ext_java_lang_Record() {}

func StaticGet_JsMessage_PH_JS_PREFIX() string { return "{$" }
func StaticGet_JsMessage_PH_JS_SUFFIX() string { return "}" }

var static_JsMessage_JS_PLACEHOLDER_NAME_RE *j2gort.Ext_java_util_regex_Pattern

func StaticGet_JsMessage_JS_PLACEHOLDER_NAME_RE() *j2gort.Ext_java_util_regex_Pattern {
	clinit_JsMessage()
	return static_JsMessage_JS_PLACEHOLDER_NAME_RE
}
func StaticSet_JsMessage_JS_PLACEHOLDER_NAME_RE(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_JsMessage()
	static_JsMessage_JS_PLACEHOLDER_NAME_RE = v
}

var static_JsMessage_CANONICAL_PLACEHOLDER_NAME_RE *j2gort.Ext_java_util_regex_Pattern

func StaticGet_JsMessage_CANONICAL_PLACEHOLDER_NAME_RE() *j2gort.Ext_java_util_regex_Pattern {
	clinit_JsMessage()
	return static_JsMessage_CANONICAL_PLACEHOLDER_NAME_RE
}
func StaticSet_JsMessage_CANONICAL_PLACEHOLDER_NAME_RE(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_JsMessage()
	static_JsMessage_CANONICAL_PLACEHOLDER_NAME_RE = v
}

var clinitDone_JsMessage uint32
var clinitState_JsMessage j2gort.ClinitState

func clinit_JsMessage() {
	if atomic.LoadUint32(&clinitDone_JsMessage) != 0 {
		return
	}
	clinitSlow_JsMessage()
}

func clinitSlow_JsMessage() {
	if !j2gort.EnterClinit(&clinitState_JsMessage, &clinitDone_JsMessage) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JsMessage, &clinitDone_JsMessage)
	static_JsMessage_JS_PLACEHOLDER_NAME_RE = j2gort.Ext_java_util_regex_Pattern_Compile_1("[a-z][a-zA-Z\\d]*[_\\d]*")
	static_JsMessage_CANONICAL_PLACEHOLDER_NAME_RE = j2gort.Ext_java_util_regex_Pattern_Compile_1("[A-Z\\d_]*")
}

func (this_ *JsMessage) initFields_JsMessage() {
	_ = this_
}

// Java constructor: JsMessage(java.lang.@org.jspecify.annotations.Nullable String getSourceName, java.lang.String getKey, boolean isAnonymous, boolean isExternal, java.lang.String getId, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JsMessage.Part> getParts, com.google.common.collect.ImmutableMap<com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase,com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JsMessage.Part>> getGenderedMessagesMap, java.lang.@org.jspecify.annotations.Nullable String getAlternateId, java.lang.@org.jspecify.annotations.Nullable String getDesc, java.lang.@org.jspecify.annotations.Nullable String getMeaning, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> getPlaceholderNameToExampleMap, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> getPlaceholderNameToOriginalCodeMap, com.google.common.collect.ImmutableSet<java.lang.String> jsPlaceholderNames, com.google.common.collect.ImmutableSet<java.lang.String> canonicalPlaceholderNames)
func NewJsMessage(getSourceName string, getKey string, isAnonymous bool, isExternal bool, getId string, getParts j2gort.Ext_com_google_common_collect_ImmutableList, getGenderedMessagesMap j2gort.Ext_com_google_common_collect_ImmutableMap, getAlternateId string, getDesc string, getMeaning string, getPlaceholderNameToExampleMap j2gort.Ext_com_google_common_collect_ImmutableMap, getPlaceholderNameToOriginalCodeMap j2gort.Ext_com_google_common_collect_ImmutableMap, jsPlaceholderNames j2gort.Ext_com_google_common_collect_ImmutableSet, canonicalPlaceholderNames j2gort.Ext_com_google_common_collect_ImmutableSet) *JsMessage {
	clinit_JsMessage()
	this_ := &JsMessage{}
	InitJsMessage(this_, getSourceName, getKey, isAnonymous, isExternal, getId, getParts, getGenderedMessagesMap, getAlternateId, getDesc, getMeaning, getPlaceholderNameToExampleMap, getPlaceholderNameToOriginalCodeMap, jsPlaceholderNames, canonicalPlaceholderNames)
	return this_
}

func InitJsMessage(this_ *JsMessage, getSourceName string, getKey string, isAnonymous bool, isExternal bool, getId string, getParts j2gort.Ext_com_google_common_collect_ImmutableList, getGenderedMessagesMap j2gort.Ext_com_google_common_collect_ImmutableMap, getAlternateId string, getDesc string, getMeaning string, getPlaceholderNameToExampleMap j2gort.Ext_com_google_common_collect_ImmutableMap, getPlaceholderNameToOriginalCodeMap j2gort.Ext_com_google_common_collect_ImmutableMap, jsPlaceholderNames j2gort.Ext_com_google_common_collect_ImmutableSet, canonicalPlaceholderNames j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = getSourceName
	_ = getKey
	_ = isAnonymous
	_ = isExternal
	_ = getId
	_ = getParts
	_ = getGenderedMessagesMap
	_ = getAlternateId
	_ = getDesc
	_ = getMeaning
	_ = getPlaceholderNameToExampleMap
	_ = getPlaceholderNameToOriginalCodeMap
	_ = jsPlaceholderNames
	_ = canonicalPlaceholderNames
	this_.initFields_JsMessage()
	j2gort.CheckNotNullString(getKey)
	j2gort.CheckNotNullString(getId)
	j2gort.CheckNotNull(getParts)
	j2gort.CheckNotNull(getPlaceholderNameToExampleMap)
	j2gort.CheckNotNull(getPlaceholderNameToOriginalCodeMap)
	j2gort.CheckNotNull(jsPlaceholderNames)
	j2gort.CheckNotNull(canonicalPlaceholderNames)
	this_.getSourceName = getSourceName
	this_.getKey = getKey
	this_.isAnonymous = isAnonymous
	this_.isExternal = isExternal
	this_.getId = getId
	this_.getParts = j2gort.ImmListToSliceNullable[JsMessage_Part](getParts)
	this_.getGenderedMessagesMap = getGenderedMessagesMap
	this_.getAlternateId = getAlternateId
	this_.getDesc = getDesc
	this_.getMeaning = getMeaning
	this_.getPlaceholderNameToExampleMap = getPlaceholderNameToExampleMap
	this_.getPlaceholderNameToOriginalCodeMap = getPlaceholderNameToOriginalCodeMap
	this_.jsPlaceholderNames = jsPlaceholderNames
	this_.canonicalPlaceholderNames = canonicalPlaceholderNames
}

// Java method: JsMessage.getPlaceholderOriginalCode(com.google.javascript.jscomp.JsMessage.PlaceholderReference placeholderReference)
func (this_ *JsMessage) GetPlaceholderOriginalCode(placeholderReference *JsMessage_PlaceholderReference) string {
	_ = placeholderReference
	return j2gort.CastString(j2gort.MapGetOrDefaultStr(this_.GetPlaceholderNameToOriginalCodeMap(), placeholderReference.StoredPlaceholderName(), "-"))
}

// Java method: JsMessage.getPlaceholderExample(com.google.javascript.jscomp.JsMessage.PlaceholderReference placeholderReference)
func (this_ *JsMessage) GetPlaceholderExample(placeholderReference *JsMessage_PlaceholderReference) string {
	_ = placeholderReference
	return j2gort.CastString(j2gort.MapGetOrDefaultStr(this_.GetPlaceholderNameToExampleMap(), placeholderReference.StoredPlaceholderName(), "-"))
}

// Java method: JsMessage.getGenderedMessageParts(com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase genderCase)
func (this_ *JsMessage) GetGenderedMessageParts(genderCase JsMessage_GrammaticalGenderCase) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = genderCase
	if this_.GetGenderedMessagesMap().IsEmpty() {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Message does not contain grammatical gendered variants.")))
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](this_.GetGenderedMessagesMap().Get(j2gort.BoxNullableVal(genderCase)))
}

// Java method: JsMessage.getGenderedMessageVariants()
func (this_ *JsMessage) GetGenderedMessageVariants() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.GetGenderedMessagesMap().KeySet_0()).AsList()
}

// Java method: JsMessage.asJsMessageString()
func (this_ *JsMessage) AsJsMessageString_0() string {
	if !(this_.GetGenderedMessagesMap().IsEmpty()) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "asJsMessageString() is not supported for messages with gendered variants. Please provide a grammatical gender case as an argument.")))
	}
	var sb j2gort.LocalStringBuilder

	{

		for _, p := range j2gort.CheckNotNullSlice(this_.
			Slice_GetParts()) {

			_ = p
			if p.IsPlaceholder() {
				sb.AppendStr("{$").AppendStr(p.GetJsPlaceholderName()).AppendStr("}")
			} else {
				sb.AppendStr(p.GetString())
			}
		}
	}
	return sb.String()
}

// Java method: JsMessage.asJsMessageString(com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase genderCase)
func (this_ *JsMessage) AsJsMessageString_1(genderCase JsMessage_GrammaticalGenderCase) string {
	_ = genderCase
	if !(this_.GetGenderedMessagesMap().ContainsKey(j2gort.BoxNullableVal(genderCase))) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "No message for grammtical gender case: "+j2gort.StrVal(genderCase))))
	}
	var sb j2gort.LocalStringBuilder

	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](this_.GetGenderedMessagesMap().Get(j2gort.BoxNullableVal(genderCase)))))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			p := j2gort.AsIface[JsMessage_Part](rawElem_4)
			_ = p
			if p.IsPlaceholder() {
				sb.AppendStr("{$").AppendStr(p.GetJsPlaceholderName()).AppendStr("}")
			} else {
				sb.AppendStr(p.GetString())
			}
		}
	}
	return sb.String()
}

// Java method: JsMessage.asIcuMessageString()
func (this_ *JsMessage) AsIcuMessageString_0() string {
	if !(this_.GetGenderedMessagesMap().IsEmpty()) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "asIcuMessageString() is not supported for messages with gendered variants. Please provide a grammatical gender case as an argument.")))
	}
	var sb j2gort.LocalStringBuilder

	{

		for _, p := range j2gort.CheckNotNullSlice(this_.
			Slice_GetParts()) {

			_ = p
			if p.IsPlaceholder() {
				sb.AppendChar(uint16(123)).AppendStr(p.GetCanonicalPlaceholderName()).AppendChar(uint16(125))
			} else {
				sb.AppendStr(p.GetString())
			}
		}
	}
	return sb.String()
}

// Java method: JsMessage.asIcuMessageString(com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase genderCase)
func (this_ *JsMessage) AsIcuMessageString_1(genderCase JsMessage_GrammaticalGenderCase) string {
	_ = genderCase
	if !(this_.GetGenderedMessagesMap().ContainsKey(j2gort.BoxNullableVal(genderCase))) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "No message for grammtical gender case: "+j2gort.StrVal(genderCase))))
	}
	var sb j2gort.LocalStringBuilder

	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](this_.GetGenderedMessagesMap().Get(j2gort.BoxNullableVal(genderCase)))))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			p := j2gort.AsIface[JsMessage_Part](rawElem_8)
			_ = p
			if p.IsPlaceholder() {
				sb.AppendChar(uint16(123)).AppendStr(p.GetCanonicalPlaceholderName()).AppendChar(uint16(125))
			} else {
				sb.AppendStr(p.GetString())
			}
		}
	}
	return sb.String()
}

// Java method: JsMessage.isEmpty()
func (this_ *JsMessage) IsEmpty() bool {
	if this_.GetGenderedMessagesMap().IsEmpty() {
		{

			for _, part := range j2gort.CheckNotNullSlice(this_.
				Slice_GetParts()) {

				_ = part
				if part.IsPlaceholder() || (j2gort.JStrLen(part.GetString()) > int32(0)) {
					return false
				}
			}
		}
	} else {
		{
			cur_11 := j2gort.NewMapValCursor(this_.GetGenderedMessagesMap())
			for cur_11.HasNext() {
				rawElem_12 := cur_11.Next()
				parts := j2gort.ImmListToSliceNullable[JsMessage_Part](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](rawElem_12))
				_ = parts
				{

					for _, part_2 := range parts {

						_ = part_2
						if part_2.IsPlaceholder() || (j2gort.JStrLen(part_2.GetString()) > int32(0)) {
							return false
						}
					}
				}
			}
		}
	}
	return true
}

// Java static method: JsMessage.isLowerCamelCaseWithNumericSuffixes(java.lang.String input)
func JsMessage_IsLowerCamelCaseWithNumericSuffixes(input string) bool {
	clinit_JsMessage()
	_ = input
	return StaticGet_JsMessage_JS_PLACEHOLDER_NAME_RE().Matcher(j2gort.BoxString(input)).Matches()
}

// Java static method: JsMessage.toLowerCamelCaseWithNumericSuffixes(java.lang.String input)
func JsMessage_ToLowerCamelCaseWithNumericSuffixes(input string) string {
	clinit_JsMessage()
	_ = input
	var suffixStart int32 = j2gort.JStrLen(input)
	_ = suffixStart
	for suffixStart > int32(0) {
		var ch uint16 = uint16(0)
		_ = ch
		var numberStart int32 = suffixStart
		_ = numberStart
		for numberStart > int32(0) {
			ch = j2gort.JStrCharAt(input, (numberStart - int32(1)))
			if j2gort.BoxChar_IsDigit_1_uint16(ch) {
				numberStart = numberStart - 1
			} else {
				break
			}
		}
		if ((numberStart > int32(0)) && (numberStart < suffixStart)) && (int32(ch) == int32(uint16(95))) {
			suffixStart = (numberStart - int32(1))
		} else {
			break
		}
	}
	if suffixStart == j2gort.JStrLen(input) {
		return j2gort.Ext_com_google_common_base_CaseFormat_F_UPPER_UNDERSCORE.To(j2gort.Ext_com_google_common_base_CaseFormat_F_LOWER_CAMEL, input)
	} else {
		return (j2gort.StrValString(j2gort.Ext_com_google_common_base_CaseFormat_F_UPPER_UNDERSCORE.To(j2gort.Ext_com_google_common_base_CaseFormat_F_LOWER_CAMEL, j2gort.JStrSub(input, int32(0), suffixStart))) + j2gort.StrValString(j2gort.JStrSubFrom(input, suffixStart)))
	}
}

// Java static method: JsMessage.isCanonicalPlaceholderNameFormat(java.lang.String name)
func JsMessage_IsCanonicalPlaceholderNameFormat(name string) bool {
	clinit_JsMessage()
	_ = name
	return StaticGet_JsMessage_CANONICAL_PLACEHOLDER_NAME_RE().Matcher(j2gort.BoxString(name)).Matches()
}

func (this_ *JsMessage) GetSourceName() string { return this_.getSourceName }
func (this_ *JsMessage) GetKey() string        { return this_.getKey }
func (this_ *JsMessage) IsAnonymous() bool     { return this_.isAnonymous }
func (this_ *JsMessage) IsExternal() bool      { return this_.isExternal }
func (this_ *JsMessage) GetId() string         { return this_.getId }
func (this_ *JsMessage) GetParts() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(this_.getParts)
}
func (this_ *JsMessage) GetGenderedMessagesMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.getGenderedMessagesMap
}
func (this_ *JsMessage) GetAlternateId() string { return this_.getAlternateId }
func (this_ *JsMessage) GetDesc() string        { return this_.getDesc }
func (this_ *JsMessage) GetMeaning() string     { return this_.getMeaning }
func (this_ *JsMessage) GetPlaceholderNameToExampleMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.getPlaceholderNameToExampleMap
}
func (this_ *JsMessage) GetPlaceholderNameToOriginalCodeMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.getPlaceholderNameToOriginalCodeMap
}
func (this_ *JsMessage) JsPlaceholderNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.jsPlaceholderNames
}
func (this_ *JsMessage) CanonicalPlaceholderNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.canonicalPlaceholderNames
}
func (this_ *JsMessage) Equals(other_ any) bool {
	o_, ok_ := other_.(*JsMessage)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.getSourceName != o_.getSourceName {
		return false
	}
	if this_.getKey != o_.getKey {
		return false
	}
	if this_.isAnonymous != o_.isAnonymous {
		return false
	}
	if this_.isExternal != o_.isExternal {
		return false
	}
	if this_.getId != o_.getId {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.ImmListFromTypedSlice(this_.getParts), j2gort.ImmListFromTypedSlice(o_.getParts)) {
		return false
	}
	if !javalang.ObjectsEquals(this_.getGenderedMessagesMap, o_.getGenderedMessagesMap) {
		return false
	}
	if this_.getAlternateId != o_.getAlternateId {
		return false
	}
	if this_.getDesc != o_.getDesc {
		return false
	}
	if this_.getMeaning != o_.getMeaning {
		return false
	}
	if !javalang.ObjectsEquals(this_.getPlaceholderNameToExampleMap, o_.getPlaceholderNameToExampleMap) {
		return false
	}
	if !javalang.ObjectsEquals(this_.getPlaceholderNameToOriginalCodeMap, o_.getPlaceholderNameToOriginalCodeMap) {
		return false
	}
	if !javalang.ObjectsEquals(this_.jsPlaceholderNames, o_.jsPlaceholderNames) {
		return false
	}
	if !javalang.ObjectsEquals(this_.canonicalPlaceholderNames, o_.canonicalPlaceholderNames) {
		return false
	}
	return true
}
func (this_ *JsMessage) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_15 := this_.getSourceName
		h_ = 31 * h_
		if hs_15 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_15)
		}
	}
	{
		hs_16 := this_.getKey
		h_ = 31 * h_
		if hs_16 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_16)
		}
	}
	if this_.isAnonymous {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	if this_.isExternal {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	{
		hs_17 := this_.getId
		h_ = 31 * h_
		if hs_17 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_17)
		}
	}
	h_ = 31*h_ + j2gort.ObjectHash(j2gort.ImmListFromTypedSlice(this_.getParts))
	h_ = 31*h_ + j2gort.ObjectHash(this_.getGenderedMessagesMap)
	{
		hs_18 := this_.getAlternateId
		h_ = 31 * h_
		if hs_18 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_18)
		}
	}
	{
		hs_19 := this_.getDesc
		h_ = 31 * h_
		if hs_19 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_19)
		}
	}
	{
		hs_20 := this_.getMeaning
		h_ = 31 * h_
		if hs_20 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_20)
		}
	}
	h_ = 31*h_ + j2gort.ObjectHash(this_.getPlaceholderNameToExampleMap)
	h_ = 31*h_ + j2gort.ObjectHash(this_.getPlaceholderNameToOriginalCodeMap)
	h_ = 31*h_ + j2gort.ObjectHash(this_.jsPlaceholderNames)
	h_ = 31*h_ + j2gort.ObjectHash(this_.canonicalPlaceholderNames)
	return h_
}
func (this_ *JsMessage) String() string {
	_ = fmt.Sprint
	return "JsMessage[" + "getSourceName=" + j2gort.StrVal(j2gort.BoxString(this_.getSourceName)) + ", " + "getKey=" + j2gort.StrVal(j2gort.BoxString(this_.getKey)) + ", " + "isAnonymous=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.isAnonymous)) + ", " + "isExternal=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.isExternal)) + ", " + "getId=" + j2gort.StrVal(j2gort.BoxString(this_.getId)) + ", " + "getParts=" + j2gort.StrVal(j2gort.ImmListFromTypedSlice(this_.getParts)) + ", " + "getGenderedMessagesMap=" + j2gort.StrVal(this_.getGenderedMessagesMap) + ", " + "getAlternateId=" + j2gort.StrVal(j2gort.BoxString(this_.getAlternateId)) + ", " + "getDesc=" + j2gort.StrVal(j2gort.BoxString(this_.getDesc)) + ", " + "getMeaning=" + j2gort.StrVal(j2gort.BoxString(this_.getMeaning)) + ", " + "getPlaceholderNameToExampleMap=" + j2gort.StrVal(this_.getPlaceholderNameToExampleMap) + ", " + "getPlaceholderNameToOriginalCodeMap=" + j2gort.StrVal(this_.getPlaceholderNameToOriginalCodeMap) + ", " + "jsPlaceholderNames=" + j2gort.StrVal(this_.jsPlaceholderNames) + ", " + "canonicalPlaceholderNames=" + j2gort.StrVal(this_.canonicalPlaceholderNames) + "]"
}

// Java enum: com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase
type JsMessage_GrammaticalGenderCase uint8

const (
	enum_JsMessage_GrammaticalGenderCase_MASCULINE JsMessage_GrammaticalGenderCase = 1
	enum_JsMessage_GrammaticalGenderCase_FEMININE  JsMessage_GrammaticalGenderCase = 2
	enum_JsMessage_GrammaticalGenderCase_NEUTER    JsMessage_GrammaticalGenderCase = 3
	enum_JsMessage_GrammaticalGenderCase_OTHER     JsMessage_GrammaticalGenderCase = 4
)

var enumBaseTable_JsMessage_GrammaticalGenderCase = [...]j2gort.EnumBase{
	{Name_: "MASCULINE", Ordinal_: 0},
	{Name_: "FEMININE", Ordinal_: 1},
	{Name_: "NEUTER", Ordinal_: 2},
	{Name_: "OTHER", Ordinal_: 3},
}
var enumValues_JsMessage_GrammaticalGenderCase = []JsMessage_GrammaticalGenderCase{1, 2, 3, 4}

func StaticGet_JsMessage_GrammaticalGenderCase_MASCULINE() JsMessage_GrammaticalGenderCase {
	return enum_JsMessage_GrammaticalGenderCase_MASCULINE
}
func StaticGet_JsMessage_GrammaticalGenderCase_FEMININE() JsMessage_GrammaticalGenderCase {
	return enum_JsMessage_GrammaticalGenderCase_FEMININE
}
func StaticGet_JsMessage_GrammaticalGenderCase_NEUTER() JsMessage_GrammaticalGenderCase {
	return enum_JsMessage_GrammaticalGenderCase_NEUTER
}
func StaticGet_JsMessage_GrammaticalGenderCase_OTHER() JsMessage_GrammaticalGenderCase {
	return enum_JsMessage_GrammaticalGenderCase_OTHER
}

func JsMessage_GrammaticalGenderCase_Values() []JsMessage_GrammaticalGenderCase {
	return j2gort.CloneSlice(enumValues_JsMessage_GrammaticalGenderCase)
}

func JsMessage_GrammaticalGenderCase_ValueOf(name string) JsMessage_GrammaticalGenderCase {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JsMessage_GrammaticalGenderCase {
		if enumBaseTable_JsMessage_GrammaticalGenderCase[i_].Name_ == name {
			return JsMessage_GrammaticalGenderCase(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase." + name))
}

func (this_ JsMessage_GrammaticalGenderCase) IsJavaNil_() bool                              { return this_ == 0 }
func (JsMessage_GrammaticalGenderCase) HasIdentityEquals_()                                 {}
func (JsMessage_GrammaticalGenderCase) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JsMessage_GrammaticalGenderCase) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JsMessage_GrammaticalGenderCase) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JsMessage_GrammaticalGenderCase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessage$GrammaticalGenderCase", "GrammaticalGenderCase"
}
func (this_ JsMessage_GrammaticalGenderCase) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JsMessage_GrammaticalGenderCase[this_-1]
}
func (this_ JsMessage_GrammaticalGenderCase) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JsMessage_GrammaticalGenderCase) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JsMessage_GrammaticalGenderCase) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsMessage_GrammaticalGenderCase[this_-1].Name_
}
func (this_ JsMessage_GrammaticalGenderCase) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JsMessage_GrammaticalGenderCase[this_-1].Name_
}
func (this_ JsMessage_GrammaticalGenderCase) Equals(other_ any) bool {
	o_, ok_ := other_.(JsMessage_GrammaticalGenderCase)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JsMessage_GrammaticalGenderCase) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JsMessage_GrammaticalGenderCase[this_-1].HashCode()
}
func (this_ JsMessage_GrammaticalGenderCase) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsMessage_GrammaticalGenderCase[this_-1].CompareTo(other_)
}
func (this_ JsMessage_GrammaticalGenderCase) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JsMessage_GrammaticalGenderCase[this_-1].Name_)
}
func init() {
	vals_ := []any{JsMessage_GrammaticalGenderCase(1), JsMessage_GrammaticalGenderCase(2), JsMessage_GrammaticalGenderCase(3), JsMessage_GrammaticalGenderCase(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.JsMessage$GrammaticalGenderCase", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase", "com.google.javascript.jscomp.JsMessage$GrammaticalGenderCase", func(v_ any) bool {
		return j2gort.InstanceOf[JsMessage_GrammaticalGenderCase](v_)
	})
}

// Java class: com.google.javascript.jscomp.JsMessage.PlaceholderFormatException
type JsMessage_PlaceholderFormatException struct {
	javalang.JavaException
	idHash_ atomic.Int32
}

func (this_ *JsMessage_PlaceholderFormatException) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessage_PlaceholderFormatException) IsJavaNil_() bool { return this_ == nil }
func (*JsMessage_PlaceholderFormatException) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessage$PlaceholderFormatException", "PlaceholderFormatException"
}
func (this_ *JsMessage_PlaceholderFormatException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessage.PlaceholderFormatException", "com.google.javascript.jscomp.JsMessage$PlaceholderFormatException", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessage_PlaceholderFormatException](v_)
	})
}

func (this_ *JsMessage_PlaceholderFormatException) initFields_JsMessage_PlaceholderFormatException() {
	_ = this_
}

// Java constructor: PlaceholderFormatException(java.lang.String msg)
func NewJsMessage_PlaceholderFormatException(msg string) *JsMessage_PlaceholderFormatException {
	this_ := &JsMessage_PlaceholderFormatException{}
	InitJsMessage_PlaceholderFormatException(this_, msg)
	return this_
}

func InitJsMessage_PlaceholderFormatException(this_ *JsMessage_PlaceholderFormatException, msg string) {
	_ = msg
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.JsMessage.PlaceholderFormatException", j2gort.BoxString(msg))
	this_.initFields_JsMessage_PlaceholderFormatException()
}

// Java interface: com.google.javascript.jscomp.JsMessage.Part
type JsMessage_Part interface {
	Implements_jscomp_JsMessage_Part()
	IsJavaNil_() bool
	IsPlaceholder() bool
	GetJsPlaceholderName() string
	GetCanonicalPlaceholderName() string
	GetString() string
}

func ToIface_JsMessage_Part[T interface {
	JsMessage_Part
	comparable
}](v T) JsMessage_Part {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JsMessage_Part[T JsMessage_Part](v T) JsMessage_Part {
	var vi JsMessage_Part = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessage.Part", "com.google.javascript.jscomp.JsMessage$Part", func(v_ any) bool {
		return j2gort.InstanceOf[JsMessage_Part](v_)
	})
}

// Java record: com.google.javascript.jscomp.JsMessage.StringPart
type JsMessage_StringPart struct {
	string_ string `gson:"string"`
}

func (this_ *JsMessage_StringPart) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessage_StringPart) IsJavaNil_() bool { return this_ == nil }
func (*JsMessage_StringPart) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessage$StringPart", "StringPart"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessage.StringPart", "com.google.javascript.jscomp.JsMessage$StringPart", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessage_StringPart](v_)
	})
}
func (*JsMessage_StringPart) Implements_j2gort_Ext_java_lang_Record() {}
func (*JsMessage_StringPart) Implements_jscomp_JsMessage_Part()       {}

func (this_ *JsMessage_StringPart) initFields_JsMessage_StringPart() {
	_ = this_
}

// Java constructor: StringPart(java.lang.String string)
func NewJsMessage_StringPart(string_ string) *JsMessage_StringPart {
	this_ := &JsMessage_StringPart{}
	InitJsMessage_StringPart(this_, string_)
	return this_
}

func InitJsMessage_StringPart(this_ *JsMessage_StringPart, string_ string) {
	_ = string_
	this_.initFields_JsMessage_StringPart()
	j2gort.CheckNotNullString(string_)
	this_.string_ = string_
}

// Java method: StringPart.getString()
func (this_ *JsMessage_StringPart) GetString() string {
	return this_.M_String()
}

// Java static method: StringPart.create(java.lang.String str)
func JsMessage_StringPart_Create(str string) *JsMessage_StringPart {
	_ = str
	return NewJsMessage_StringPart(str)
}

// Java method: StringPart.isPlaceholder()
func (this_ *JsMessage_StringPart) IsPlaceholder() bool {
	return false
}

// Java method: StringPart.getJsPlaceholderName()
func (this_ *JsMessage_StringPart) GetJsPlaceholderName() string {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("not a placeholder: '%s'", []any{j2gort.BoxString(this_.M_String())})))))
}

// Java method: StringPart.getCanonicalPlaceholderName()
func (this_ *JsMessage_StringPart) GetCanonicalPlaceholderName() string {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("not a placeholder: '%s'", []any{j2gort.BoxString(this_.M_String())})))))
}

func (this_ *JsMessage_StringPart) M_String() string { return this_.string_ }
func (this_ *JsMessage_StringPart) Equals(other_ any) bool {
	o_, ok_ := other_.(*JsMessage_StringPart)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.string_ != o_.string_ {
		return false
	}
	return true
}
func (this_ *JsMessage_StringPart) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_15 := this_.string_
		h_ = 31 * h_
		if hs_15 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_15)
		}
	}
	return h_
}
func (this_ *JsMessage_StringPart) String() string {
	_ = fmt.Sprint
	return "StringPart[" + "string=" + j2gort.StrVal(j2gort.BoxString(this_.string_)) + "]"
}

// Java record: com.google.javascript.jscomp.JsMessage.PlaceholderReference
type JsMessage_PlaceholderReference struct {
	storedPlaceholderName string `gson:"storedPlaceholderName"`
	canonicalFormat       bool   `gson:"canonicalFormat"`
}

func (this_ *JsMessage_PlaceholderReference) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessage_PlaceholderReference) IsJavaNil_() bool { return this_ == nil }
func (*JsMessage_PlaceholderReference) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessage$PlaceholderReference", "PlaceholderReference"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessage.PlaceholderReference", "com.google.javascript.jscomp.JsMessage$PlaceholderReference", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessage_PlaceholderReference](v_)
	})
}
func (*JsMessage_PlaceholderReference) Implements_j2gort_Ext_java_lang_Record() {}
func (*JsMessage_PlaceholderReference) Implements_jscomp_JsMessage_Part()       {}

func (this_ *JsMessage_PlaceholderReference) initFields_JsMessage_PlaceholderReference() {
	_ = this_
}

// Java constructor: PlaceholderReference(java.lang.String storedPlaceholderName, boolean canonicalFormat)
func NewJsMessage_PlaceholderReference(storedPlaceholderName string, canonicalFormat bool) *JsMessage_PlaceholderReference {
	this_ := &JsMessage_PlaceholderReference{}
	InitJsMessage_PlaceholderReference(this_, storedPlaceholderName, canonicalFormat)
	return this_
}

func InitJsMessage_PlaceholderReference(this_ *JsMessage_PlaceholderReference, storedPlaceholderName string, canonicalFormat bool) {
	_ = storedPlaceholderName
	_ = canonicalFormat
	this_.initFields_JsMessage_PlaceholderReference()
	j2gort.CheckNotNullString(storedPlaceholderName)
	this_.storedPlaceholderName = storedPlaceholderName
	this_.canonicalFormat = canonicalFormat
}

// Java method: PlaceholderReference.getStoredPlaceholderName()
func (this_ *JsMessage_PlaceholderReference) GetStoredPlaceholderName() string {
	return this_.StoredPlaceholderName()
}

// Java method: PlaceholderReference.isCanonicalFormat()
func (this_ *JsMessage_PlaceholderReference) IsCanonicalFormat() bool {
	return this_.CanonicalFormat()
}

// Java static method: PlaceholderReference.createForJsName(java.lang.String name)
func JsMessage_PlaceholderReference_CreateForJsName(name string) *JsMessage_PlaceholderReference {
	_ = name
	if !(JsMessage_IsLowerCamelCaseWithNumericSuffixes(name)) {
		j2gort.CheckArgumentFmt(false, "invalid JS placeholder name format: '%s'", j2gort.BoxString(name))
	}
	return NewJsMessage_PlaceholderReference(name, false)
}

// Java static method: PlaceholderReference.createForCanonicalName(java.lang.String name)
func JsMessage_PlaceholderReference_CreateForCanonicalName(name string) *JsMessage_PlaceholderReference {
	_ = name
	if !(JsMessage_IsCanonicalPlaceholderNameFormat(name)) {
		j2gort.CheckArgumentFmt(false, "not a canonical placeholder name format: '%s'", j2gort.BoxString(name))
	}
	return NewJsMessage_PlaceholderReference(name, true)
}

// Java method: PlaceholderReference.isPlaceholder()
func (this_ *JsMessage_PlaceholderReference) IsPlaceholder() bool {
	return true
}

// Java method: PlaceholderReference.getJsPlaceholderName()
func (this_ *JsMessage_PlaceholderReference) GetJsPlaceholderName() string {
	var storedPlaceholderName string = this_.StoredPlaceholderName()
	_ = storedPlaceholderName
	if this_.CanonicalFormat() {
		return JsMessage_ToLowerCamelCaseWithNumericSuffixes(storedPlaceholderName)
	} else {
		return storedPlaceholderName
	}
}

// Java method: PlaceholderReference.getCanonicalPlaceholderName()
func (this_ *JsMessage_PlaceholderReference) GetCanonicalPlaceholderName() string {
	var storedPlaceholderName string = this_.StoredPlaceholderName()
	_ = storedPlaceholderName
	if this_.CanonicalFormat() {
		return storedPlaceholderName
	} else {
		return j2gort.Ext_com_google_common_base_CaseFormat_F_LOWER_CAMEL.To(j2gort.Ext_com_google_common_base_CaseFormat_F_UPPER_UNDERSCORE, storedPlaceholderName)
	}
}

// Java method: PlaceholderReference.getString()
func (this_ *JsMessage_PlaceholderReference) GetString() string {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("not a string part: '%s'", []any{j2gort.BoxString(this_.GetJsPlaceholderName())})))))
}

func (this_ *JsMessage_PlaceholderReference) StoredPlaceholderName() string {
	return this_.storedPlaceholderName
}
func (this_ *JsMessage_PlaceholderReference) CanonicalFormat() bool { return this_.canonicalFormat }
func (this_ *JsMessage_PlaceholderReference) Equals(other_ any) bool {
	o_, ok_ := other_.(*JsMessage_PlaceholderReference)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.storedPlaceholderName != o_.storedPlaceholderName {
		return false
	}
	if this_.canonicalFormat != o_.canonicalFormat {
		return false
	}
	return true
}
func (this_ *JsMessage_PlaceholderReference) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_15 := this_.storedPlaceholderName
		h_ = 31 * h_
		if hs_15 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_15)
		}
	}
	if this_.canonicalFormat {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	return h_
}
func (this_ *JsMessage_PlaceholderReference) String() string {
	_ = fmt.Sprint
	return "PlaceholderReference[" + "storedPlaceholderName=" + j2gort.StrVal(j2gort.BoxString(this_.storedPlaceholderName)) + ", " + "canonicalFormat=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.canonicalFormat)) + "]"
}

// Java class: com.google.javascript.jscomp.JsMessage.Builder
type JsMessage_Builder struct {
	key                              string                                                                        `gson:"key"`
	meaning                          string                                                                        `gson:"meaning"`
	desc                             string                                                                        `gson:"desc"`
	id                               string                                                                        `gson:"id"`
	alternateId                      string                                                                        `gson:"alternateId"`
	parts                            j2gort.Ext_java_util_List                                                     `gson:"parts"`
	genderedMessageMap               j2gort.OrderedMap[JsMessage_GrammaticalGenderCase, j2gort.Ext_java_util_List] `gson:"genderedMessageMap"`
	jsPlaceholderNames               j2gort.OrderedSet[string]                                                     `gson:"jsPlaceholderNames"`
	canonicalPlaceholderNames        j2gort.OrderedSet[string]                                                     `gson:"canonicalPlaceholderNames"`
	placeholderNameToExampleMap      j2gort.Ext_com_google_common_collect_ImmutableMap                             `gson:"placeholderNameToExampleMap"`
	placeholderNameToOriginalCodeMap j2gort.Ext_com_google_common_collect_ImmutableMap                             `gson:"placeholderNameToOriginalCodeMap"`
	sourceName                       string                                                                        `gson:"sourceName"`
	idHash_                          atomic.Int32
	isAnonymous                      bool `gson:"isAnonymous"`
	isExternal                       bool `gson:"isExternal"`
}

func (this_ *JsMessage_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessage_Builder) IsJavaNil_() bool { return this_ == nil }
func (*JsMessage_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessage$Builder", "Builder"
}
func (this_ *JsMessage_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessage_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessage.Builder", "com.google.javascript.jscomp.JsMessage$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessage_Builder](v_)
	})
}

func (this_ *JsMessage_Builder) initFields_JsMessage_Builder() {
	_ = this_
	this_.key = j2gort.NullString
	this_.meaning = j2gort.NullString
	this_.desc = j2gort.NullString

	this_.id = j2gort.NullString
	this_.alternateId = j2gort.NullString
	this_.parts = j2gort.NewExt_java_util_ArrayList()
	this_.genderedMessageMap = j2gort.NewOrderedMap[JsMessage_GrammaticalGenderCase, j2gort.Ext_java_util_List]()
	this_.jsPlaceholderNames = j2gort.NewOrderedSet[string]()
	this_.canonicalPlaceholderNames = j2gort.NewOrderedSet[string]()
	this_.placeholderNameToExampleMap = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	this_.placeholderNameToOriginalCodeMap = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	this_.sourceName = j2gort.NullString
}

// Java constructor: Builder()
func NewJsMessage_Builder() *JsMessage_Builder {
	this_ := &JsMessage_Builder{}
	InitJsMessage_Builder(this_)
	return this_
}

func InitJsMessage_Builder(this_ *JsMessage_Builder) {
	this_.initFields_JsMessage_Builder()
}

// Java method: Builder.getKey()
func (this_ *JsMessage_Builder) GetKey() string {
	return this_.key
}

// Java method: Builder.setKey(java.lang.String key)
func (this_ *JsMessage_Builder) SetKey(key string) *JsMessage_Builder {
	_ = key
	this_.key = key
	return this_
}

// Java method: Builder.setSourceName(java.lang.String sourceName)
func (this_ *JsMessage_Builder) SetSourceName(sourceName string) *JsMessage_Builder {
	_ = sourceName
	this_.sourceName = sourceName
	return this_
}

// Java method: Builder.appendPart(com.google.javascript.jscomp.JsMessage.Part part)
func (this_ *JsMessage_Builder) AppendPart(part JsMessage_Part) *JsMessage_Builder {
	_ = part
	j2gort.CheckNotNullIface(part)
	this_.parts.Add_1(part)
	if part.IsPlaceholder() {
		this_.jsPlaceholderNames.Add(part.GetJsPlaceholderName())
		this_.canonicalPlaceholderNames.Add(part.GetCanonicalPlaceholderName())
	}
	return this_
}

// Java method: Builder.appendParts(java.util.List<com.google.javascript.jscomp.JsMessage.Part> parts)
func (this_ *JsMessage_Builder) AppendParts(parts j2gort.Ext_java_util_List) *JsMessage_Builder {
	_ = parts
	{
		cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(parts))
		for cur_15.HasNext() {
			rawElem_16 := cur_15.Next()
			part := j2gort.AsIface[JsMessage_Part](rawElem_16)
			_ = part
			this_.AppendPart(part)
		}
	}
	return this_
}

// Java method: Builder.appendJsPlaceholderReference(java.lang.String name)
func (this_ *JsMessage_Builder) AppendJsPlaceholderReference_1(name string) *JsMessage_Builder {
	_ = name
	j2gort.CheckNotNullString(name)
	this_.parts.Add_1(j2gort.PtrToAny(JsMessage_PlaceholderReference_CreateForJsName(name)))
	this_.jsPlaceholderNames.Add(name)
	return this_
}

// Java method: Builder.appendJsPlaceholderReference(com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase genderCase, java.lang.String name)
func (this_ *JsMessage_Builder) AppendJsPlaceholderReference_2(genderCase JsMessage_GrammaticalGenderCase, name string) *JsMessage_Builder {
	_ = genderCase
	_ = name
	j2gort.CheckNotNullString(name)
	this_.genderedMessageMap.Get(genderCase).Add_1(j2gort.PtrToAny(JsMessage_PlaceholderReference_CreateForJsName(name)))
	this_.jsPlaceholderNames.Add(name)
	return this_
}

// Java method: Builder.appendCanonicalPlaceholderReference(java.lang.String name)
func (this_ *JsMessage_Builder) AppendCanonicalPlaceholderReference_1(name string) *JsMessage_Builder {
	_ = name
	j2gort.CheckNotNullString(name)
	var placeholder *JsMessage_PlaceholderReference = JsMessage_PlaceholderReference_CreateForCanonicalName(name)
	_ = placeholder
	this_.parts.Add_1(j2gort.PtrToAny(placeholder))
	this_.jsPlaceholderNames.Add(placeholder.GetJsPlaceholderName())
	return this_
}

// Java method: Builder.appendCanonicalPlaceholderReference(com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase genderCase, java.lang.String name)
func (this_ *JsMessage_Builder) AppendCanonicalPlaceholderReference_2(genderCase JsMessage_GrammaticalGenderCase, name string) *JsMessage_Builder {
	_ = genderCase
	_ = name
	j2gort.CheckNotNullString(name)
	var placeholder *JsMessage_PlaceholderReference = JsMessage_PlaceholderReference_CreateForCanonicalName(name)
	_ = placeholder
	this_.genderedMessageMap.Get(genderCase).Add_1(j2gort.PtrToAny(placeholder))
	this_.jsPlaceholderNames.Add(placeholder.GetJsPlaceholderName())
	return this_
}

// Java method: Builder.appendStringPart(java.lang.String part)
func (this_ *JsMessage_Builder) AppendStringPart_1(part string) *JsMessage_Builder {
	_ = part
	j2gort.CheckNotNullString(part)
	this_.parts.Add_1(j2gort.PtrToAny(JsMessage_StringPart_Create(part)))
	return this_
}

// Java method: Builder.appendStringPart(com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase genderCase, java.lang.String part)
func (this_ *JsMessage_Builder) AppendStringPart_2(genderCase JsMessage_GrammaticalGenderCase, part string) *JsMessage_Builder {
	_ = genderCase
	_ = part
	j2gort.CheckNotNullEnum(genderCase)
	j2gort.CheckNotNullString(part)
	this_.genderedMessageMap.Get(genderCase).Add_1(j2gort.PtrToAny(JsMessage_StringPart_Create(part)))
	return this_
}

// Java method: Builder.addGenderedMessageKey(com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase part)
func (this_ *JsMessage_Builder) AddGenderedMessageKey(part JsMessage_GrammaticalGenderCase) *JsMessage_Builder {
	_ = part
	j2gort.CheckNotNullEnum(part)
	this_.genderedMessageMap.Put(part, j2gort.NewExt_java_util_ArrayList())
	return this_
}

// Java method: Builder.setPlaceholderNameToExampleMap(java.util.Map<java.lang.String,java.lang.String> map)
func (this_ *JsMessage_Builder) SetPlaceholderNameToExampleMap(map_ j2gort.Ext_java_util_Map) *JsMessage_Builder {
	_ = map_
	this_.placeholderNameToExampleMap = j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(map_)
	return this_
}

// Java method: Builder.setPlaceholderNameToOriginalCodeMap(java.util.Map<java.lang.String,java.lang.String> map)
func (this_ *JsMessage_Builder) SetPlaceholderNameToOriginalCodeMap(map_ j2gort.Ext_java_util_Map) *JsMessage_Builder {
	_ = map_
	this_.placeholderNameToOriginalCodeMap = j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(map_)
	return this_
}

// Java method: Builder.setDesc(java.lang.@org.jspecify.annotations.Nullable String desc)
func (this_ *JsMessage_Builder) SetDesc(desc string) *JsMessage_Builder {
	_ = desc
	this_.desc = desc
	return this_
}

// Java method: Builder.setMeaning(java.lang.@org.jspecify.annotations.Nullable String meaning)
func (this_ *JsMessage_Builder) SetMeaning(meaning string) *JsMessage_Builder {
	_ = meaning
	this_.meaning = meaning
	return this_
}

// Java method: Builder.setAlternateId(java.lang.@org.jspecify.annotations.Nullable String alternateId)
func (this_ *JsMessage_Builder) SetAlternateId(alternateId string) *JsMessage_Builder {
	_ = alternateId
	this_.alternateId = alternateId
	return this_
}

// Java method: Builder.setIsAnonymous(boolean isAnonymous)
func (this_ *JsMessage_Builder) SetIsAnonymous(isAnonymous bool) *JsMessage_Builder {
	_ = isAnonymous
	this_.isAnonymous = isAnonymous
	return this_
}

// Java method: Builder.setId(java.lang.String id)
func (this_ *JsMessage_Builder) SetId(id string) *JsMessage_Builder {
	_ = id
	if !(this_.id == j2gort.NullString) {
		j2gort.CheckStateFmt(false, "id already set to '%s': cannot change it to '%s'", j2gort.BoxString(this_.id), j2gort.BoxString(id))
	}
	this_.id = id
	return this_
}

// Java method: Builder.setIsExternalMsg(boolean isExternalMsg)
func (this_ *JsMessage_Builder) SetIsExternalMsg(isExternalMsg bool) *JsMessage_Builder {
	_ = isExternalMsg
	this_.isExternal = isExternalMsg
	return this_
}

// Java method: Builder.hasParts()
func (this_ *JsMessage_Builder) HasParts() bool {
	return !(this_.parts.IsEmpty())
}

// Java method: Builder.hasGenderedVariants()
func (this_ *JsMessage_Builder) HasGenderedVariants() bool {
	return !(this_.genderedMessageMap.IsEmpty())
}

// Java method: Builder.getParts()
func (this_ *JsMessage_Builder) GetParts() j2gort.Ext_java_util_List {
	return this_.parts
}

// Java method: Builder.build()
func (this_ *JsMessage_Builder) Build() *JsMessage {
	j2gort.CheckNotNullString(this_.key)
	j2gort.CheckNotNullString(this_.id)
	if !(!(this_.isExternal) || !(this_.isAnonymous)) {
		j2gort.CheckStateStr(false, "a message cannot be both anonymous and external")
	}
	if (this_.alternateId != j2gort.NullString) && func() bool {
		sr_17 := this_.alternateId
		se_18 := this_.id
		return j2gort.CheckNotNullString(sr_17) == se_18
	}() {
		this_.alternateId = j2gort.NullString
	}
	var immutableGenderedPartsBuilder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = immutableGenderedPartsBuilder
	{
		cur_19 := this_.genderedMessageMap.IterEntries()
		for cur_19.NextEntry() {
			immutableGenderedPartsBuilder.Put_2(j2gort.BoxNullableVal(cur_19.Key()), j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(cur_19.Value()))
		}
	}
	var immutableGenderedMessageMap j2gort.Ext_com_google_common_collect_ImmutableMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](immutableGenderedPartsBuilder.BuildOrThrow())
	_ = immutableGenderedMessageMap
	return NewJsMessage(this_.sourceName, this_.key, this_.isAnonymous, this_.isExternal, this_.id, j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.parts)), immutableGenderedMessageMap, this_.alternateId, this_.desc, this_.meaning, this_.placeholderNameToExampleMap, this_.placeholderNameToOriginalCodeMap, this_.jsPlaceholderNames.ToImmutableSet(), this_.canonicalPlaceholderNames.ToImmutableSet())
}

// Java class: com.google.javascript.jscomp.JsMessage.Hash
type JsMessage_Hash struct {
	idHash_ atomic.Int32
}

func (this_ *JsMessage_Hash) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessage_Hash) IsJavaNil_() bool { return this_ == nil }
func (*JsMessage_Hash) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessage$Hash", "Hash"
}
func (this_ *JsMessage_Hash) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessage_Hash) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessage.Hash", "com.google.javascript.jscomp.JsMessage$Hash", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessage_Hash](v_)
	})
}

func StaticGet_JsMessage_Hash_SEED64() int64     { return int64(3141592653589793238) }
func StaticGet_JsMessage_Hash_CONSTANT64() int64 { return int64(-2266404186210603134) }
func (this_ *JsMessage_Hash) initFields_JsMessage_Hash() {
	_ = this_
}

// Java constructor: Hash()
func NewJsMessage_Hash() *JsMessage_Hash {
	this_ := &JsMessage_Hash{}
	InitJsMessage_Hash(this_)
	return this_
}

func InitJsMessage_Hash(this_ *JsMessage_Hash) {
	this_.initFields_JsMessage_Hash()
}

// Java static method: Hash.hash64(java.lang.@org.jspecify.annotations.Nullable String value)
func JsMessage_Hash_Hash64_1(value string) int64 {
	_ = value
	return JsMessage_Hash_Hash64_2_string_int64(value, int64(3141592653589793238))
}

// Java static method: Hash.hash64(java.lang.@org.jspecify.annotations.Nullable String value, long seed)
func JsMessage_Hash_Hash64_2_string_int64(value string, seed int64) int64 {
	_ = value
	_ = seed
	if value == j2gort.NullString {
		return JsMessage_Hash_Hash64_4(nil, int32(0), int32(0), seed)
	}
	return JsMessage_Hash_Hash64_2_slice_int8_int64(func() []int8 {
		_ = j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8
		return j2gort.JStrGetBytes(value)
	}(), seed)
}

// Java static method: Hash.hash64(byte[] value, long seed)
func JsMessage_Hash_Hash64_2_slice_int8_int64(value []int8, seed int64) int64 {
	_ = value
	_ = seed
	return JsMessage_Hash_Hash64_4(value, int32(0), func() int32 {
		if value == nil {
			return int32(0)
		}
		return int32(len(value))
	}(), seed)
}

// Java static method: Hash.hash64(byte @org.jspecify.annotations.Nullable [] value, int offset, int length, long seed)
func JsMessage_Hash_Hash64_4(value []int8, offset int32, length int32, seed int64) int64 {
	_ = value
	_ = offset
	_ = length
	_ = seed
	var a int64 = int64(-2266404186210603134)
	_ = a
	var b int64 = a
	_ = b
	var c int64 = seed
	_ = c
	var keylen int32 = int32(0)
	_ = keylen
	{
		keylen = length
		for ; keylen >= int32(24); func() { keylen = (keylen - int32(24)); offset = (offset + int32(24)) }() {
			a = (a + JsMessage_Hash_Word64At(value, offset))
			b = (b + JsMessage_Hash_Word64At(value, (offset+int32(8))))
			c = (c + JsMessage_Hash_Word64At(value, (offset+int32(16))))
			a = (a - b)
			a = (a - c)
			a = (a ^ javalang.Ushr64(c, int(int32(43))))
			b = (b - c)
			b = (b - a)
			b = (b ^ javalang.Shl64(a, int(int32(9))))
			c = (c - a)
			c = (c - b)
			c = (c ^ javalang.Ushr64(b, int(int32(8))))
			a = (a - b)
			a = (a - c)
			a = (a ^ javalang.Ushr64(c, int(int32(38))))
			b = (b - c)
			b = (b - a)
			b = (b ^ javalang.Shl64(a, int(int32(23))))
			c = (c - a)
			c = (c - b)
			c = (c ^ javalang.Ushr64(b, int(int32(5))))
			a = (a - b)
			a = (a - c)
			a = (a ^ javalang.Ushr64(c, int(int32(35))))
			b = (b - c)
			b = (b - a)
			b = (b ^ javalang.Shl64(a, int(int32(49))))
			c = (c - a)
			c = (c - b)
			c = (c ^ javalang.Ushr64(b, int(int32(11))))
			a = (a - b)
			a = (a - c)
			a = (a ^ javalang.Ushr64(c, int(int32(12))))
			b = (b - c)
			b = (b - a)
			b = (b ^ javalang.Shl64(a, int(int32(18))))
			c = (c - a)
			c = (c - b)
			c = (c ^ javalang.Ushr64(b, int(int32(22))))
		}
	}
	c = (c + int64(length))
	if keylen >= int32(16) {
		if keylen == int32(23) {
			c = (c + javalang.Shl64((int64(value[(offset+int32(22))])), int(int32(56))))
		}
		if keylen >= int32(22) {
			c = (c + javalang.Shl64((int64(value[(offset+int32(21))])&int64(255)), int(int32(48))))
		}
		if keylen >= int32(21) {
			c = (c + javalang.Shl64((int64(value[(offset+int32(20))])&int64(255)), int(int32(40))))
		}
		if keylen >= int32(20) {
			c = (c + javalang.Shl64((int64(value[(offset+int32(19))])&int64(255)), int(int32(32))))
		}
		if keylen >= int32(19) {
			c = (c + javalang.Shl64((int64(value[(offset+int32(18))])&int64(255)), int(int32(24))))
		}
		if keylen >= int32(18) {
			c = (c + javalang.Shl64((int64(value[(offset+int32(17))])&int64(255)), int(int32(16))))
		}
		if keylen >= int32(17) {
			c = (c + javalang.Shl64((int64(value[(offset+int32(16))])&int64(255)), int(int32(8))))
		}
		if keylen >= int32(16) {
			b = (b + JsMessage_Hash_Word64At(value, (offset+int32(8))))
			a = (a + JsMessage_Hash_Word64At(value, offset))
		}
	} else if keylen >= int32(8) {
		if keylen == int32(15) {
			b = (b + javalang.Shl64((int64(value[(offset+int32(14))])&int64(255)), int(int32(48))))
		}
		if keylen >= int32(14) {
			b = (b + javalang.Shl64((int64(value[(offset+int32(13))])&int64(255)), int(int32(40))))
		}
		if keylen >= int32(13) {
			b = (b + javalang.Shl64((int64(value[(offset+int32(12))])&int64(255)), int(int32(32))))
		}
		if keylen >= int32(12) {
			b = (b + javalang.Shl64((int64(value[(offset+int32(11))])&int64(255)), int(int32(24))))
		}
		if keylen >= int32(11) {
			b = (b + javalang.Shl64((int64(value[(offset+int32(10))])&int64(255)), int(int32(16))))
		}
		if keylen >= int32(10) {
			b = (b + javalang.Shl64((int64(value[(offset+int32(9))])&int64(255)), int(int32(8))))
		}
		if keylen >= int32(9) {
			b = (b + (int64(value[(offset+int32(8))]) & int64(255)))
		}
		if keylen >= int32(8) {
			a = (a + JsMessage_Hash_Word64At(value, offset))
		}
	} else {
		if keylen == int32(7) {
			a = (a + javalang.Shl64((int64(value[(offset+int32(6))])&int64(255)), int(int32(48))))
		}
		if keylen >= int32(6) {
			a = (a + javalang.Shl64((int64(value[(offset+int32(5))])&int64(255)), int(int32(40))))
		}
		if keylen >= int32(5) {
			a = (a + javalang.Shl64((int64(value[(offset+int32(4))])&int64(255)), int(int32(32))))
		}
		if keylen >= int32(4) {
			a = (a + javalang.Shl64((int64(value[(offset+int32(3))])&int64(255)), int(int32(24))))
		}
		if keylen >= int32(3) {
			a = (a + javalang.Shl64((int64(value[(offset+int32(2))])&int64(255)), int(int32(16))))
		}
		if keylen >= int32(2) {
			a = (a + javalang.Shl64((int64(value[(offset+int32(1))])&int64(255)), int(int32(8))))
		}
		if keylen >= int32(1) {
			a = (a + (int64(value[(offset+int32(0))]) & int64(255)))
		}
	}
	return JsMessage_Hash_Mix64(a, b, c)
}

// Java static method: Hash.word64At(byte[] bytes, int offset)
func JsMessage_Hash_Word64At(bytes []int8, offset int32) int64 {
	_ = bytes
	_ = offset
	return ((((((((int64(bytes[(offset+int32(0))]) & int64(255)) + (javalang.Shl64((int64(bytes[(offset+int32(1))]) & int64(255)), int(int32(8))))) + (javalang.Shl64((int64(bytes[(offset+int32(2))]) & int64(255)), int(int32(16))))) + (javalang.Shl64((int64(bytes[(offset+int32(3))]) & int64(255)), int(int32(24))))) + (javalang.Shl64((int64(bytes[(offset+int32(4))]) & int64(255)), int(int32(32))))) + (javalang.Shl64((int64(bytes[(offset+int32(5))]) & int64(255)), int(int32(40))))) + (javalang.Shl64((int64(bytes[(offset+int32(6))]) & int64(255)), int(int32(48))))) + (javalang.Shl64((int64(bytes[(offset+int32(7))]) & int64(255)), int(int32(56)))))
}

// Java static method: Hash.mix64(long a, long b, long c)
func JsMessage_Hash_Mix64(a int64, b int64, c int64) int64 {
	_ = a
	_ = b
	_ = c
	a = (a - b)
	a = (a - c)
	a = (a ^ javalang.Ushr64(c, int(int32(43))))
	b = (b - c)
	b = (b - a)
	b = (b ^ javalang.Shl64(a, int(int32(9))))
	c = (c - a)
	c = (c - b)
	c = (c ^ javalang.Ushr64(b, int(int32(8))))
	a = (a - b)
	a = (a - c)
	a = (a ^ javalang.Ushr64(c, int(int32(38))))
	b = (b - c)
	b = (b - a)
	b = (b ^ javalang.Shl64(a, int(int32(23))))
	c = (c - a)
	c = (c - b)
	c = (c ^ javalang.Ushr64(b, int(int32(5))))
	a = (a - b)
	a = (a - c)
	a = (a ^ javalang.Ushr64(c, int(int32(35))))
	b = (b - c)
	b = (b - a)
	b = (b ^ javalang.Shl64(a, int(int32(49))))
	c = (c - a)
	c = (c - b)
	c = (c ^ javalang.Ushr64(b, int(int32(11))))
	a = (a - b)
	a = (a - c)
	a = (a ^ javalang.Ushr64(c, int(int32(12))))
	b = (b - c)
	b = (b - a)
	b = (b ^ javalang.Shl64(a, int(int32(18))))
	c = (c - a)
	c = (c - b)
	c = (c ^ javalang.Ushr64(b, int(int32(22))))
	return c
}

// Java interface: com.google.javascript.jscomp.JsMessage.IdGenerator
type JsMessage_IdGenerator interface {
	Implements_jscomp_JsMessage_IdGenerator()
	IsJavaNil_() bool
	GenerateId(string, j2gort.Ext_java_util_List) string
}

func ToIface_JsMessage_IdGenerator[T interface {
	JsMessage_IdGenerator
	comparable
}](v T) JsMessage_IdGenerator {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JsMessage_IdGenerator[T JsMessage_IdGenerator](v T) JsMessage_IdGenerator {
	var vi JsMessage_IdGenerator = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessage.IdGenerator", "com.google.javascript.jscomp.JsMessage$IdGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[JsMessage_IdGenerator](v_)
	})
}
func (this_ *JsMessage) Slice_GetParts() []JsMessage_Part {
	return this_.getParts
}

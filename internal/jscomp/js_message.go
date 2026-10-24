// Port of com/google/javascript/jscomp/JsMessage.java @ closure-compiler v20261021

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

package jscomp

import (
	"regexp"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JsMessage is a representation of a translatable message in JavaScript
// source code.
//
// Instances are created using a JsMessage.Builder, like this:
//
//	JsMessage m = new JsMessage.Builder(key)
//	    .appendPart("Hi ")
//	    .appendPlaceholderReference("firstName")
//	    .appendPart("!")
//	    .setDesc("A welcome message")
//	    .build();
//
// The record components:
//   - getSourceName: the message's sourceName.
//   - getKey: the message's key, or name (e.g. "MSG_HELLO").
//   - getId: the message's id, or name (e.g. "92430284230902938293").
//   - getParts: a read-only list of the parts of this message. Each part is
//     either a String or a PlaceholderReference.
//   - getAlternateId: the message's alternate ID (e.g.
//     "92430284230902938293"), if available. This will be used if a
//     translation for `id` is not available.
//   - getDesc: the description associated with this message, intended to
//     help translators, or null if this message has no description.
//   - getMeaning: the meaning annotated to the message, intended to force
//     different translations.
//   - jsPlaceholderNames: a set of the registered placeholders in this
//     message in lowerCamelCase format. This is the format used in
//     `goog.getMsg()` declarations.
//   - canonicalPlaceholderNames: a set of the registered placeholders in
//     this message in UPPER_SNAKE_CASE format. This is the format stored
//     into XMB / XTB files and used in `declareIcuTemplate()
//     declarations.
//
// Go-only: language-adapter — a record: fields named after the components,
// accessor methods, and the record's equals/hashCode/toString.
type JsMessage struct {
	getSourceName                       *string
	getKey                              string
	isAnonymous                         bool
	isExternal                          bool
	getId                               string
	getParts                            []JsMessage_Part
	getGenderedMessagesMap              *jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part]
	getAlternateId                      *string
	getDesc                             *string
	getMeaning                          *string
	getPlaceholderNameToExampleMap      *jutil.LinkedHashMap[string, string]
	getPlaceholderNameToOriginalCodeMap *jutil.LinkedHashMap[string, string]
	jsPlaceholderNames                  *jutil.LinkedHashSet[string]
	canonicalPlaceholderNames           *jutil.LinkedHashSet[string]
}

// jsMessage_requireNonNull is Objects.requireNonNull(obj, message).
func jsMessage_requireNonNull(isNull bool, message string) {
	if isNull {
		panic(jutil.NewNullPointerExceptionMsg(message))
	}
}

// NewJsMessage is the record's canonical constructor.
//
// Java: JsMessage(String, String, boolean, boolean, String, ImmutableList<Part>, ImmutableMap<GrammaticalGenderCase, ImmutableList<Part>>, String, String, String, ImmutableMap<String, String>, ImmutableMap<String, String>, ImmutableSet<String>, ImmutableSet<String>)
func NewJsMessage(
	getSourceName *string,
	getKey string,
	isAnonymous bool,
	isExternal bool,
	getId string,
	getParts []JsMessage_Part,
	getGenderedMessagesMap *jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part],
	getAlternateId *string,
	getDesc *string,
	getMeaning *string,
	getPlaceholderNameToExampleMap *jutil.LinkedHashMap[string, string],
	getPlaceholderNameToOriginalCodeMap *jutil.LinkedHashMap[string, string],
	jsPlaceholderNames *jutil.LinkedHashSet[string],
	canonicalPlaceholderNames *jutil.LinkedHashSet[string]) *JsMessage {
	// requireNonNull(getKey, "getKey"), requireNonNull(getId, "getId"): a Go
	// string is never null. requireNonNull(getParts, "getParts"): a nil Go
	// slice is the empty list.
	jsMessage_requireNonNull(getPlaceholderNameToExampleMap == nil, "getPlaceholderNameToExampleMap")
	jsMessage_requireNonNull(getPlaceholderNameToOriginalCodeMap == nil, "getPlaceholderNameToOriginalCodeMap")
	jsMessage_requireNonNull(jsPlaceholderNames == nil, "jsPlaceholderNames")
	jsMessage_requireNonNull(canonicalPlaceholderNames == nil, "canonicalPlaceholderNames")
	return &JsMessage{
		getSourceName:                       getSourceName,
		getKey:                              getKey,
		isAnonymous:                         isAnonymous,
		isExternal:                          isExternal,
		getId:                               getId,
		getParts:                            getParts,
		getGenderedMessagesMap:              getGenderedMessagesMap,
		getAlternateId:                      getAlternateId,
		getDesc:                             getDesc,
		getMeaning:                          getMeaning,
		getPlaceholderNameToExampleMap:      getPlaceholderNameToExampleMap,
		getPlaceholderNameToOriginalCodeMap: getPlaceholderNameToOriginalCodeMap,
		jsPlaceholderNames:                  jsPlaceholderNames,
		canonicalPlaceholderNames:           canonicalPlaceholderNames,
	}
}

// Java: JsMessage.getSourceName()
func (m *JsMessage) GetSourceName() *string { return m.getSourceName }

// Java: JsMessage.getKey()
func (m *JsMessage) GetKey() string { return m.getKey }

// Java: JsMessage.isAnonymous()
func (m *JsMessage) IsAnonymous() bool { return m.isAnonymous }

// Java: JsMessage.isExternal()
func (m *JsMessage) IsExternal() bool { return m.isExternal }

// Java: JsMessage.getId()
func (m *JsMessage) GetId() string { return m.getId }

// Java: JsMessage.getParts()
func (m *JsMessage) GetParts() []JsMessage_Part { return m.getParts }

// Java: JsMessage.getGenderedMessagesMap()
func (m *JsMessage) GetGenderedMessagesMap() *jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part] {
	return m.getGenderedMessagesMap
}

// Java: JsMessage.getAlternateId()
func (m *JsMessage) GetAlternateId() *string { return m.getAlternateId }

// Java: JsMessage.getDesc()
func (m *JsMessage) GetDesc() *string { return m.getDesc }

// Java: JsMessage.getMeaning()
func (m *JsMessage) GetMeaning() *string { return m.getMeaning }

// Java: JsMessage.getPlaceholderNameToExampleMap()
func (m *JsMessage) GetPlaceholderNameToExampleMap() *jutil.LinkedHashMap[string, string] {
	return m.getPlaceholderNameToExampleMap
}

// Java: JsMessage.getPlaceholderNameToOriginalCodeMap()
func (m *JsMessage) GetPlaceholderNameToOriginalCodeMap() *jutil.LinkedHashMap[string, string] {
	return m.getPlaceholderNameToOriginalCodeMap
}

// Java: JsMessage.jsPlaceholderNames()
func (m *JsMessage) JsPlaceholderNames() *jutil.LinkedHashSet[string] { return m.jsPlaceholderNames }

// Java: JsMessage.canonicalPlaceholderNames()
func (m *JsMessage) CanonicalPlaceholderNames() *jutil.LinkedHashSet[string] {
	return m.canonicalPlaceholderNames
}

// Java: JsMessage.PH_JS_PREFIX
const JsMessage_PH_JS_PREFIX = "{$"

// Java: JsMessage.PH_JS_SUFFIX
const JsMessage_PH_JS_SUFFIX = "}"

// JsMessage_GrammaticalGenderCase is the enum for grammatical gender
// cases.
type JsMessage_GrammaticalGenderCase int32

const (
	JsMessage_GrammaticalGenderCase_MASCULINE JsMessage_GrammaticalGenderCase = iota
	JsMessage_GrammaticalGenderCase_FEMININE
	JsMessage_GrammaticalGenderCase_NEUTER
	JsMessage_GrammaticalGenderCase_OTHER
)

var jsMessage_GrammaticalGenderCase_names = [...]string{"MASCULINE", "FEMININE", "NEUTER", "OTHER"}

// JsMessage_GrammaticalGenderCase_values is GrammaticalGenderCase.values().
func JsMessage_GrammaticalGenderCase_values() []JsMessage_GrammaticalGenderCase {
	return []JsMessage_GrammaticalGenderCase{
		JsMessage_GrammaticalGenderCase_MASCULINE,
		JsMessage_GrammaticalGenderCase_FEMININE,
		JsMessage_GrammaticalGenderCase_NEUTER,
		JsMessage_GrammaticalGenderCase_OTHER,
	}
}

// JsMessage_GrammaticalGenderCase_ValueOf is GrammaticalGenderCase.valueOf(String).
func JsMessage_GrammaticalGenderCase_ValueOf(name string) JsMessage_GrammaticalGenderCase {
	for i, n := range jsMessage_GrammaticalGenderCase_names {
		if n == name {
			return JsMessage_GrammaticalGenderCase(i)
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(
		"No enum constant com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase." + name))
}

// Name is Enum.name().
func (g JsMessage_GrammaticalGenderCase) Name() string {
	return jsMessage_GrammaticalGenderCase_names[g]
}

// String is Enum.toString().
func (g JsMessage_GrammaticalGenderCase) String() string { return g.Name() }

// Ordinal is Enum.ordinal().
func (g JsMessage_GrammaticalGenderCase) Ordinal() int { return int(g) }

// JsMessage_PlaceholderFormatException is thrown when parsing a message
// string into parts fails because of a misformatted place holder.
type JsMessage_PlaceholderFormatException struct {
	jutil.JavaException
}

func (*JsMessage_PlaceholderFormatException) placeholderFormatException() {}

// IJsMessage_PlaceholderFormatException matches PlaceholderFormatException.
type IJsMessage_PlaceholderFormatException interface {
	jutil.IException
	placeholderFormatException()
}

// Java: JsMessage.PlaceholderFormatException(String)
func NewJsMessage_PlaceholderFormatException(msg string) *JsMessage_PlaceholderFormatException {
	e := &JsMessage_PlaceholderFormatException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.JsMessage$PlaceholderFormatException", jutil.StringPtr(msg), nil)
	return e
}

// Java: JsMessage.getPlaceholderOriginalCode(PlaceholderReference)
func (m *JsMessage) GetPlaceholderOriginalCode(placeholderReference *JsMessage_PlaceholderReference) string {
	return m.GetPlaceholderNameToOriginalCodeMap().
		GetOrDefault(placeholderReference.StoredPlaceholderName(), "-")
}

// Java: JsMessage.getPlaceholderExample(PlaceholderReference)
func (m *JsMessage) GetPlaceholderExample(placeholderReference *JsMessage_PlaceholderReference) string {
	return m.GetPlaceholderNameToExampleMap().
		GetOrDefault(placeholderReference.StoredPlaceholderName(), "-")
}

// GetGenderedMessageParts returns the parts for genderCase (nil, Java
// null, if there are none for it).
//
// Java: JsMessage.getGenderedMessageParts(GrammaticalGenderCase)
func (m *JsMessage) GetGenderedMessageParts(genderCase JsMessage_GrammaticalGenderCase) []JsMessage_Part {
	if m.GetGenderedMessagesMap().IsEmpty() {
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"Message does not contain grammatical gendered variants."))
	}
	return m.GetGenderedMessagesMap().Get(genderCase)
}

// GetGenderedMessageVariants gets the list of grammatical gender cases for
// the message.
//
// Java: JsMessage.getGenderedMessageVariants()
func (m *JsMessage) GetGenderedMessageVariants() []JsMessage_GrammaticalGenderCase {
	return m.GetGenderedMessagesMap().KeySlice()
}

// AsJsMessageString returns a single string representing the message.
//
// In the returned string all placeholders are joined with the literal
// string parts in the form "literal string part {$jsPlaceholderName} more
// literal string", which is how placeholders are represented in
// `goog.getMsg()` text strings.
//
// If the message has gendered variants, the returned string will be a
// multiline string with each line containing the grammatical gender case,
// followed by its message.
//
// Throws UnsupportedOperationException if the message has gendered
// variants.
//
// Java: JsMessage.asJsMessageString()
func (m *JsMessage) AsJsMessageString() string {
	if !m.GetGenderedMessagesMap().IsEmpty() {
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"asJsMessageString() is not supported for messages with gendered variants. Please provide" +
				" a grammatical gender case as an argument."))
	}
	sb := jutil.NewStringBuilder()
	for _, p := range m.GetParts() {
		if p.IsPlaceholder() {
			sb.Append(JsMessage_PH_JS_PREFIX).Append(p.GetJsPlaceholderName()).Append(JsMessage_PH_JS_SUFFIX)
		} else {
			sb.Append(p.GetString())
		}
	}

	return sb.String()
}

// AsJsMessageStringForGender is the overloaded method that returns a single
// string representing the message for a given gender case.
//
// In the returned string all placeholders are joined with the literal
// string parts in the form "literal string part {$jsPlaceholderName} more
// literal string", which is how placeholders are represented in
// `goog.getMsg()` text strings.
//
// Throws IllegalArgumentException if the message does not have a message
// for the given gender case.
//
// Java: JsMessage.asJsMessageString(GrammaticalGenderCase)
func (m *JsMessage) AsJsMessageStringForGender(genderCase JsMessage_GrammaticalGenderCase) string {
	if !m.GetGenderedMessagesMap().ContainsKey(genderCase) {
		panic(jutil.NewIllegalArgumentExceptionMsg("No message for grammtical gender case: " + genderCase.String()))
	}
	sb := jutil.NewStringBuilder()
	for _, p := range m.GetGenderedMessagesMap().Get(genderCase) {
		if p.IsPlaceholder() {
			sb.Append(JsMessage_PH_JS_PREFIX).Append(p.GetJsPlaceholderName()).Append(JsMessage_PH_JS_SUFFIX)
		} else {
			sb.Append(p.GetString())
		}
	}
	return sb.String()
}

// AsIcuMessageString returns a single string representing the message.
//
// In the returned string all placeholders are joined with the literal
// string parts in the form "literal string part
// {CANONICAL_PLACEHOLDER_NAME} more literal string", which is how
// placeholders are represented in ICU messages.
//
// Throws UnsupportedOperationException if the message has gendered
// variants.
//
// Java: JsMessage.asIcuMessageString()
func (m *JsMessage) AsIcuMessageString() string {
	if !m.GetGenderedMessagesMap().IsEmpty() {
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"asIcuMessageString() is not supported for messages with gendered variants. Please" +
				" provide a grammatical gender case as an argument."))
	}
	sb := jutil.NewStringBuilder()
	for _, p := range m.GetParts() {
		if p.IsPlaceholder() {
			sb.AppendChar('{').Append(p.GetCanonicalPlaceholderName()).AppendChar('}')
		} else {
			sb.Append(p.GetString())
		}
	}
	return sb.String()
}

// AsIcuMessageStringForGender is the overloaded method that returns a
// single string representing the message for a given gender case.
//
// In the returned string all placeholders are joined with the literal
// string parts in the form "literal string part
// {CANONICAL_PLACEHOLDER_NAME} more literal string", which is how
// placeholders are represented in ICU messages.
//
// Throws IllegalArgumentException if the message does not have a message
// for the given gender case.
//
// Java: JsMessage.asIcuMessageString(GrammaticalGenderCase)
func (m *JsMessage) AsIcuMessageStringForGender(genderCase JsMessage_GrammaticalGenderCase) string {
	if !m.GetGenderedMessagesMap().ContainsKey(genderCase) {
		panic(jutil.NewIllegalArgumentExceptionMsg("No message for grammtical gender case: " + genderCase.String()))
	}
	sb := jutil.NewStringBuilder()
	for _, p := range m.GetGenderedMessagesMap().Get(genderCase) {
		if p.IsPlaceholder() {
			sb.AppendChar('{').Append(p.GetCanonicalPlaceholderName()).AppendChar('}')
		} else {
			sb.Append(p.GetString())
		}
	}
	return sb.String()
}

// IsEmpty returns false iff the message is represented by empty string.
//
// Java: JsMessage.isEmpty()
func (m *JsMessage) IsEmpty() bool {
	if m.GetGenderedMessagesMap().IsEmpty() {
		for _, part := range m.GetParts() {
			if part.IsPlaceholder() || len(part.GetString()) > 0 {
				return false
			}
		}
	} else {
		for parts := range m.GetGenderedMessagesMap().Values() {
			for _, part := range parts {
				if part.IsPlaceholder() || len(part.GetString()) > 0 {
					return false
				}
			}
		}
	}

	return true
}

// Equals is the record's equals: every component is equal.
//
// Java: JsMessage.equals(Object)
func (m *JsMessage) Equals(o any) bool {
	that, ok := o.(*JsMessage)
	if !ok || that == nil {
		return false
	}
	return jutil.ObjectsEquals(m.getSourceName, that.getSourceName) &&
		m.getKey == that.getKey &&
		m.isAnonymous == that.isAnonymous &&
		m.isExternal == that.isExternal &&
		m.getId == that.getId &&
		jsMessage_partsEqual(m.getParts, that.getParts) &&
		jsMessage_genderedEqual(m.getGenderedMessagesMap, that.getGenderedMessagesMap) &&
		jutil.ObjectsEquals(m.getAlternateId, that.getAlternateId) &&
		jutil.ObjectsEquals(m.getDesc, that.getDesc) &&
		jutil.ObjectsEquals(m.getMeaning, that.getMeaning) &&
		jutil.ObjectsEquals(m.getPlaceholderNameToExampleMap, that.getPlaceholderNameToExampleMap) &&
		jutil.ObjectsEquals(m.getPlaceholderNameToOriginalCodeMap, that.getPlaceholderNameToOriginalCodeMap) &&
		jutil.ObjectsEquals(m.jsPlaceholderNames, that.jsPlaceholderNames) &&
		jutil.ObjectsEquals(m.canonicalPlaceholderNames, that.canonicalPlaceholderNames)
}

func jsMessage_partsEqual(a, b []JsMessage_Part) bool {
	return slices.EqualFunc(a, b, func(x, y JsMessage_Part) bool { return jutil.ObjectsEquals(x, y) })
}

func jsMessage_genderedEqual(a, b *jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part]) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if a.Size() != b.Size() {
		return false
	}
	for k, v := range a.All() {
		w, ok := b.Lookup(k)
		if !ok || !jsMessage_partsEqual(v, w) {
			return false
		}
	}
	return true
}

func jsMessage_partsHashCode(parts []JsMessage_Part) int32 {
	// List.hashCode()
	var h int32 = 1
	for _, p := range parts {
		h = 31*h + jutil.ObjectHashCode(p)
	}
	return h
}

func jsMessage_boolHashCode(b bool) int32 {
	if b {
		return 1231
	}
	return 1237
}

// HashCode is the record's hashCode: h = 31*h + hash(component), from 0.
//
// Java: JsMessage.hashCode()
func (m *JsMessage) HashCode() int32 {
	var h int32 = 0
	h = 31*h + jutil.ObjectHashCode(m.getSourceName)
	h = 31*h + jutil.StringHashCode(m.getKey)
	h = 31*h + jsMessage_boolHashCode(m.isAnonymous)
	h = 31*h + jsMessage_boolHashCode(m.isExternal)
	h = 31*h + jutil.StringHashCode(m.getId)
	h = 31*h + jsMessage_partsHashCode(m.getParts)
	var gendered int32 = 0
	if m.getGenderedMessagesMap != nil {
		// Map.hashCode(): the sum of key.hashCode() ^ value.hashCode().
		for k, v := range m.getGenderedMessagesMap.All() {
			gendered += jutil.ObjectHashCode(k) ^ jsMessage_partsHashCode(v)
		}
	}
	h = 31*h + gendered
	h = 31*h + jutil.ObjectHashCode(m.getAlternateId)
	h = 31*h + jutil.ObjectHashCode(m.getDesc)
	h = 31*h + jutil.ObjectHashCode(m.getMeaning)
	h = 31*h + jutil.ObjectHashCode(m.getPlaceholderNameToExampleMap)
	h = 31*h + jutil.ObjectHashCode(m.getPlaceholderNameToOriginalCodeMap)
	h = 31*h + jutil.ObjectHashCode(m.jsPlaceholderNames)
	h = 31*h + jutil.ObjectHashCode(m.canonicalPlaceholderNames)
	return h
}

// String is the record's toString.
//
// Java: JsMessage.toString()
func (m *JsMessage) String() string {
	var gendered string
	if m.getGenderedMessagesMap == nil {
		gendered = "null"
	} else {
		var sb strings.Builder
		sb.WriteByte('{')
		first := true
		for k, v := range m.getGenderedMessagesMap.All() {
			if !first {
				sb.WriteString(", ")
			}
			first = false
			sb.WriteString(k.String())
			sb.WriteByte('=')
			sb.WriteString(jutil.CollectionToString(v))
		}
		sb.WriteByte('}')
		gendered = sb.String()
	}
	return "JsMessage[" +
		"getSourceName=" + jutil.ToJavaString(m.getSourceName) +
		", getKey=" + m.getKey +
		", isAnonymous=" + jutil.ToJavaString(m.isAnonymous) +
		", isExternal=" + jutil.ToJavaString(m.isExternal) +
		", getId=" + m.getId +
		", getParts=" + jutil.CollectionToString(m.getParts) +
		", getGenderedMessagesMap=" + gendered +
		", getAlternateId=" + jutil.ToJavaString(m.getAlternateId) +
		", getDesc=" + jutil.ToJavaString(m.getDesc) +
		", getMeaning=" + jutil.ToJavaString(m.getMeaning) +
		", getPlaceholderNameToExampleMap=" + jutil.ToJavaString(m.getPlaceholderNameToExampleMap) +
		", getPlaceholderNameToOriginalCodeMap=" + jutil.ToJavaString(m.getPlaceholderNameToOriginalCodeMap) +
		", jsPlaceholderNames=" + jutil.ToJavaString(m.jsPlaceholderNames) +
		", canonicalPlaceholderNames=" + jutil.ToJavaString(m.canonicalPlaceholderNames) +
		"]"
}

// JsMessage_Part represents part of a message.
type JsMessage_Part interface {
	// IsPlaceholder: true for placeholders, false for literal string parts.
	//
	// Java: JsMessage.Part.isPlaceholder()
	IsPlaceholder() bool

	// GetJsPlaceholderName gets the name of the placeholder as it would
	// appear in JS code.
	//
	// In JS code placeholders are in lower camel case with zero or more
	// trailing numeric suffixes (e.g. myPlaceholderName_123_456).
	//
	// Throws UnsupportedOperationException if this is not a
	// PlaceholderReference.
	//
	// Java: JsMessage.Part.getJsPlaceholderName()
	GetJsPlaceholderName() string

	// GetCanonicalPlaceholderName gets the name of the placeholder as it
	// would appear in XMB or XTB files, and also the form it takes when
	// generating message IDs.
	//
	// This format is UPPER_SNAKE_CASE.
	//
	// Throws UnsupportedOperationException if this is not a
	// PlaceholderReference.
	//
	// Java: JsMessage.Part.getCanonicalPlaceholderName()
	GetCanonicalPlaceholderName() string

	// GetString gets the literal string for this message part.
	//
	// Throws UnsupportedOperationException if this is not a StringPart.
	//
	// Java: JsMessage.Part.getString()
	GetString() string
}

// JsMessage_StringPart represents a literal string part of a message.
//
// Go-only: language-adapter — a record; its component accessor string() is
// String_() because String() is toString().
type JsMessage_StringPart struct {
	string_ string
}

// Java: JsMessage.StringPart(String)
func NewJsMessage_StringPart(s string) *JsMessage_StringPart {
	// requireNonNull(string, "string"): a Go string is never null.
	return &JsMessage_StringPart{string_: s}
}

// String_ is the record component accessor string().
//
// Java: JsMessage.StringPart.string()
func (p *JsMessage_StringPart) String_() string { return p.string_ }

// Java: JsMessage.StringPart.getString()
func (p *JsMessage_StringPart) GetString() string {
	return p.String_()
}

// Java: JsMessage.StringPart.create(String)
func JsMessage_StringPart_Create(str string) *JsMessage_StringPart {
	return NewJsMessage_StringPart(str)
}

// Java: JsMessage.StringPart.isPlaceholder()
func (p *JsMessage_StringPart) IsPlaceholder() bool {
	return false
}

// Java: JsMessage.StringPart.getJsPlaceholderName()
func (p *JsMessage_StringPart) GetJsPlaceholderName() string {
	panic(jutil.NewUnsupportedOperationExceptionMsg(jutil.StringFormat("not a placeholder: '%s'", p.String_())))
}

// Java: JsMessage.StringPart.getCanonicalPlaceholderName()
func (p *JsMessage_StringPart) GetCanonicalPlaceholderName() string {
	panic(jutil.NewUnsupportedOperationExceptionMsg(jutil.StringFormat("not a placeholder: '%s'", p.String_())))
}

// Java: JsMessage.StringPart.equals(Object)
func (p *JsMessage_StringPart) Equals(o any) bool {
	that, ok := o.(*JsMessage_StringPart)
	return ok && that != nil && p.string_ == that.string_
}

// Java: JsMessage.StringPart.hashCode()
func (p *JsMessage_StringPart) HashCode() int32 {
	return jutil.StringHashCode(p.string_)
}

// Java: JsMessage.StringPart.toString()
func (p *JsMessage_StringPart) String() string {
	return "StringPart[string=" + p.string_ + "]"
}

// jsMessage_JS_PLACEHOLDER_NAME_RE: in JS code we expect placeholder names
// to be lowerCamelCase with an optional _123_456 suffix.
var jsMessage_JS_PLACEHOLDER_NAME_RE = regexp.MustCompile(`\A(?:[a-z][a-zA-Z0-9]*[_0-9]*)\z`)

// JsMessage_IsLowerCamelCaseWithNumericSuffixes returns whether a string is
// nonempty, begins with a lowercase letter, and contains only digits and
// underscores after the first underscore.
//
// Java: JsMessage.isLowerCamelCaseWithNumericSuffixes(String)
func JsMessage_IsLowerCamelCaseWithNumericSuffixes(input string) bool {
	return jsMessage_JS_PLACEHOLDER_NAME_RE.MatchString(input)
}

// JsMessage_ToLowerCamelCaseWithNumericSuffixes converts the given string
// from upper-underscore case to lower-camel case, preserving numeric
// suffixes. For example: "NAME" -> "name" "A4_LETTER" -> "a4Letter"
// "START_SPAN_1_23" -> "startSpan_1_23".
//
// Java: JsMessage.toLowerCamelCaseWithNumericSuffixes(String)
func JsMessage_ToLowerCamelCaseWithNumericSuffixes(input string) string {
	units := jutil.StringToUTF16(input)
	// Determine where the numeric suffixes begin
	suffixStart := len(units)
	for suffixStart > 0 {
		var ch uint16 = 0
		numberStart := suffixStart
		for numberStart > 0 {
			ch = units[numberStart-1]
			if jutil.CharacterIsDigit(ch) {
				numberStart--
			} else {
				break
			}
		}
		if (numberStart > 0) && (numberStart < suffixStart) && (ch == '_') {
			suffixStart = numberStart - 1
		} else {
			break
		}
	}

	if suffixStart == len(units) {
		return jutil.CaseFormat_UpperUnderscoreToLowerCamel(input)
	} else {
		return jutil.CaseFormat_UpperUnderscoreToLowerCamel(jutil.StringFromUTF16(units[:suffixStart])) +
			jutil.StringFromUTF16(units[suffixStart:])
	}
}

// jsMessage_CANONICAL_PLACEHOLDER_NAME_RE: in XMB and XTB files we expect
// the placeholder name to be in UPPER_SNAKE_CASE with optional _123_456
// suffix.
//
// This pattern will allow alpha sections after numeric sections (e.g.
// NOT_A_1234_VALID_NAME), but it is still probably good enough, since these
// name are usually automatically generated from the JS formatted
// lowerCamelCase_with_optional_12_34 names.
//
// The pattern also allows leading and trailing underscores. There are known
// cases of messages containing those.
var jsMessage_CANONICAL_PLACEHOLDER_NAME_RE = regexp.MustCompile(`\A(?:[A-Z0-9_]*)\z`)

// JsMessage_IsCanonicalPlaceholderNameFormat: is the name in the canonical
// format for placeholder names in XTB and XMB files.
//
// Java: JsMessage.isCanonicalPlaceholderNameFormat(String)
func JsMessage_IsCanonicalPlaceholderNameFormat(name string) bool {
	return jsMessage_CANONICAL_PLACEHOLDER_NAME_RE.MatchString(name)
}

// JsMessage_PlaceholderReference is a reference to a placeholder in a
// translatable message.
//
// Go-only: language-adapter — a record.
type JsMessage_PlaceholderReference struct {
	storedPlaceholderName string
	canonicalFormat       bool
}

// Java: JsMessage.PlaceholderReference(String, boolean)
func NewJsMessage_PlaceholderReference(storedPlaceholderName string, canonicalFormat bool) *JsMessage_PlaceholderReference {
	// requireNonNull(storedPlaceholderName, "storedPlaceholderName"): a Go
	// string is never null.
	return &JsMessage_PlaceholderReference{storedPlaceholderName: storedPlaceholderName, canonicalFormat: canonicalFormat}
}

// Java: JsMessage.PlaceholderReference.storedPlaceholderName()
func (p *JsMessage_PlaceholderReference) StoredPlaceholderName() string {
	return p.storedPlaceholderName
}

// Java: JsMessage.PlaceholderReference.canonicalFormat()
func (p *JsMessage_PlaceholderReference) CanonicalFormat() bool {
	return p.canonicalFormat
}

// Java: JsMessage.PlaceholderReference.getStoredPlaceholderName()
func (p *JsMessage_PlaceholderReference) GetStoredPlaceholderName() string {
	return p.StoredPlaceholderName()
}

// Java: JsMessage.PlaceholderReference.isCanonicalFormat()
func (p *JsMessage_PlaceholderReference) IsCanonicalFormat() bool {
	return p.CanonicalFormat()
}

// Java: JsMessage.PlaceholderReference.createForJsName(String)
func JsMessage_PlaceholderReference_CreateForJsName(name string) *JsMessage_PlaceholderReference {
	jutil.CheckArgumentMsg1(
		JsMessage_IsLowerCamelCaseWithNumericSuffixes(name),
		"invalid JS placeholder name format: '%s'",
		name)
	return NewJsMessage_PlaceholderReference(name /* canonicalFormat= */, false)
}

// Java: JsMessage.PlaceholderReference.createForCanonicalName(String)
func JsMessage_PlaceholderReference_CreateForCanonicalName(name string) *JsMessage_PlaceholderReference {
	jutil.CheckArgumentMsg1(
		JsMessage_IsCanonicalPlaceholderNameFormat(name),
		"not a canonical placeholder name format: '%s'",
		name)
	return NewJsMessage_PlaceholderReference(name /* canonicalFormat= */, true)
}

// Java: JsMessage.PlaceholderReference.isPlaceholder()
func (p *JsMessage_PlaceholderReference) IsPlaceholder() bool {
	return true
}

// Java: JsMessage.PlaceholderReference.getJsPlaceholderName()
func (p *JsMessage_PlaceholderReference) GetJsPlaceholderName() string {
	storedPlaceholderName := p.StoredPlaceholderName()
	if p.CanonicalFormat() {
		return JsMessage_ToLowerCamelCaseWithNumericSuffixes(storedPlaceholderName)
	}
	return storedPlaceholderName
}

// Java: JsMessage.PlaceholderReference.getCanonicalPlaceholderName()
func (p *JsMessage_PlaceholderReference) GetCanonicalPlaceholderName() string {
	storedPlaceholderName := p.StoredPlaceholderName()
	if p.CanonicalFormat() {
		return storedPlaceholderName
	}
	return jutil.CaseFormat_LowerCamelToUpperUnderscore(storedPlaceholderName)
}

// Java: JsMessage.PlaceholderReference.getString()
func (p *JsMessage_PlaceholderReference) GetString() string {
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		jutil.StringFormat("not a string part: '%s'", p.GetJsPlaceholderName())))
}

// Java: JsMessage.PlaceholderReference.equals(Object)
func (p *JsMessage_PlaceholderReference) Equals(o any) bool {
	that, ok := o.(*JsMessage_PlaceholderReference)
	return ok && that != nil && p.storedPlaceholderName == that.storedPlaceholderName &&
		p.canonicalFormat == that.canonicalFormat
}

// Java: JsMessage.PlaceholderReference.hashCode()
func (p *JsMessage_PlaceholderReference) HashCode() int32 {
	return 31*jutil.StringHashCode(p.storedPlaceholderName) + jsMessage_boolHashCode(p.canonicalFormat)
}

// Java: JsMessage.PlaceholderReference.toString()
func (p *JsMessage_PlaceholderReference) String() string {
	return "PlaceholderReference[storedPlaceholderName=" + p.storedPlaceholderName +
		", canonicalFormat=" + jutil.ToJavaString(p.canonicalFormat) + "]"
}

// JsMessage_Builder contains functionality for creating JS messages.
// Generates authoritative keys and fingerprints for a message that must
// stay constant over time.
//
// This implementation correctly processes unnamed messages and creates a
// key for them that looks like MSG_<fingerprint value>;.
type JsMessage_Builder struct {
	key *string

	meaning *string

	desc        *string
	isAnonymous bool
	isExternal  bool

	id          *string
	alternateId *string

	parts              *jutil.ArrayList[JsMessage_Part]
	genderedMessageMap *jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, *jutil.ArrayList[JsMessage_Part]]
	// Placeholder names in JS code format (lowerCamelCase), which is used
	// for `goog.getMsg()` messages
	jsPlaceholderNames *jutil.LinkedHashSet[string]
	// Placeholder names in canonical format (UPPER_SNAKE_CASE) which is used
	// in XMB / XTB files and `declareIcuTemplate()` messages.
	canonicalPlaceholderNames        *jutil.LinkedHashSet[string]
	placeholderNameToExampleMap      *jutil.LinkedHashMap[string, string]
	placeholderNameToOriginalCodeMap *jutil.LinkedHashMap[string, string]

	sourceName *string
}

// NewJsMessage_Builder is new JsMessage.Builder().
func NewJsMessage_Builder() *JsMessage_Builder {
	return &JsMessage_Builder{
		parts:                            jutil.NewArrayList[JsMessage_Part](),
		genderedMessageMap:               jutil.NewLinkedHashMap[JsMessage_GrammaticalGenderCase, *jutil.ArrayList[JsMessage_Part]](),
		jsPlaceholderNames:               jutil.NewLinkedHashSet[string](),
		canonicalPlaceholderNames:        jutil.NewLinkedHashSet[string](),
		placeholderNameToExampleMap:      jutil.NewLinkedHashMap[string, string](),
		placeholderNameToOriginalCodeMap: jutil.NewLinkedHashMap[string, string](),
	}
}

// GetKey gets the message's key (e.g. "MSG_HELLO").
//
// Java: JsMessage.Builder.getKey()
func (b *JsMessage_Builder) GetKey() *string {
	return b.key
}

// SetKey: key is a key that should uniquely identify this message;
// typically it is the message's name (e.g. "MSG_HELLO").
//
// Java: JsMessage.Builder.setKey(String)
func (b *JsMessage_Builder) SetKey(key string) *JsMessage_Builder {
	b.key = &key
	return b
}

// SetSourceName: sourceName is the message's sourceName.
//
// Java: JsMessage.Builder.setSourceName(String)
func (b *JsMessage_Builder) SetSourceName(sourceName *string) *JsMessage_Builder {
	b.sourceName = sourceName
	return b
}

// Java: JsMessage.Builder.appendPart(Part)
func (b *JsMessage_Builder) AppendPart(part JsMessage_Part) *JsMessage_Builder {
	jutil.CheckNotNull(part)
	b.parts.Add(part)
	if part.IsPlaceholder() {
		b.jsPlaceholderNames.Add(part.GetJsPlaceholderName())
		b.canonicalPlaceholderNames.Add(part.GetCanonicalPlaceholderName())
	}
	return b
}

// Java: JsMessage.Builder.appendParts(List<Part>)
func (b *JsMessage_Builder) AppendParts(parts []JsMessage_Part) *JsMessage_Builder {
	for _, part := range parts {
		b.AppendPart(part)
	}
	return b
}

// genderedParts is genderedMessageMap.get(genderCase), which must exist
// (Java throws NullPointerException on the null list).
func (b *JsMessage_Builder) genderedParts(genderCase JsMessage_GrammaticalGenderCase) *jutil.ArrayList[JsMessage_Part] {
	parts := b.genderedMessageMap.Get(genderCase)
	if parts == nil {
		panic(jutil.NewNullPointerException())
	}
	return parts
}

// AppendJsPlaceholderReference appends a placeholder reference to the
// message.
//
// name is the placeholder name in the format used in JS code
// (lowerCamelCaseWithOptional_12_34).
//
// Java: JsMessage.Builder.appendJsPlaceholderReference(String)
func (b *JsMessage_Builder) AppendJsPlaceholderReference(name string) *JsMessage_Builder {
	// checkNotNull(name, "Placeholder name must not be null"): a Go string
	// is never null.
	b.parts.Add(JsMessage_PlaceholderReference_CreateForJsName(name))
	b.jsPlaceholderNames.Add(name)
	return b
}

// AppendJsPlaceholderReferenceForGender is the overloaded method that
// appends a placeholder reference to the corresponding gendered message.
//
// name is the placeholder name in the format used in JS code
// (lowerCamelCaseWithOptional_12_34).
//
// Java: JsMessage.Builder.appendJsPlaceholderReference(GrammaticalGenderCase, String)
func (b *JsMessage_Builder) AppendJsPlaceholderReferenceForGender(genderCase JsMessage_GrammaticalGenderCase, name string) *JsMessage_Builder {
	// checkNotNull(name, "Placeholder name must not be null"): a Go string
	// is never null.
	b.genderedParts(genderCase).Add(JsMessage_PlaceholderReference_CreateForJsName(name))
	b.jsPlaceholderNames.Add(name)
	return b
}

// AppendCanonicalPlaceholderReference appends a placeholder reference to
// the message.
//
// name is the placeholder name in the format used in XML files
// (UPPER_SNAKE_CASE_12_34).
//
// Java: JsMessage.Builder.appendCanonicalPlaceholderReference(String)
func (b *JsMessage_Builder) AppendCanonicalPlaceholderReference(name string) *JsMessage_Builder {
	// checkNotNull(name, "Placeholder name must not be null"): a Go string
	// is never null.
	placeholder := JsMessage_PlaceholderReference_CreateForCanonicalName(name)
	b.parts.Add(placeholder)
	b.jsPlaceholderNames.Add(placeholder.GetJsPlaceholderName())
	return b
}

// AppendCanonicalPlaceholderReferenceForGender is the overloaded method that
// appends a placeholder reference to the corresponding gendered message.
//
// name is the placeholder name in the format used in XML files
// (UPPER_SNAKE_CASE_12_34).
//
// Java: JsMessage.Builder.appendCanonicalPlaceholderReference(GrammaticalGenderCase, String)
func (b *JsMessage_Builder) AppendCanonicalPlaceholderReferenceForGender(
	genderCase JsMessage_GrammaticalGenderCase, name string) *JsMessage_Builder {
	// checkNotNull(name, "Placeholder name must not be null"): a Go string
	// is never null.
	placeholder := JsMessage_PlaceholderReference_CreateForCanonicalName(name)
	b.genderedParts(genderCase).Add(placeholder)
	b.jsPlaceholderNames.Add(placeholder.GetJsPlaceholderName())
	return b
}

// AppendStringPart appends a translatable string literal to the message.
//
// Java: JsMessage.Builder.appendStringPart(String)
func (b *JsMessage_Builder) AppendStringPart(part string) *JsMessage_Builder {
	// checkNotNull(part, "String part of the message must not be null"): a
	// Go string is never null.
	b.parts.Add(JsMessage_StringPart_Create(part))
	return b
}

// AppendStringPartForGender is the overloaded method that appends a
// translatable string literal to the corresponding gendered message.
//
// Java: JsMessage.Builder.appendStringPart(GrammaticalGenderCase, String)
func (b *JsMessage_Builder) AppendStringPartForGender(genderCase JsMessage_GrammaticalGenderCase, part string) *JsMessage_Builder {
	// checkNotNull(genderCase, ...), checkNotNull(part, ...): Go enums and
	// strings are never null.
	b.genderedParts(genderCase).Add(JsMessage_StringPart_Create(part))
	return b
}

// AddGenderedMessageKey adds a gendered message key to
// genderedMessageParts map.
//
// Java: JsMessage.Builder.addGenderedMessageKey(GrammaticalGenderCase)
func (b *JsMessage_Builder) AddGenderedMessageKey(part JsMessage_GrammaticalGenderCase) *JsMessage_Builder {
	// checkNotNull(part, "Gendered message key must not be null"): a Go
	// enum is never null.
	b.genderedMessageMap.Put(part, jutil.NewArrayList[JsMessage_Part]())
	return b
}

// Java: JsMessage.Builder.setPlaceholderNameToExampleMap(Map<String, String>)
func (b *JsMessage_Builder) SetPlaceholderNameToExampleMap(m *jutil.LinkedHashMap[string, string]) *JsMessage_Builder {
	b.placeholderNameToExampleMap = jutil.NewLinkedHashMapFrom(m)
	return b
}

// Java: JsMessage.Builder.setPlaceholderNameToOriginalCodeMap(Map<String, String>)
func (b *JsMessage_Builder) SetPlaceholderNameToOriginalCodeMap(m *jutil.LinkedHashMap[string, string]) *JsMessage_Builder {
	b.placeholderNameToOriginalCodeMap = jutil.NewLinkedHashMapFrom(m)
	return b
}

// SetDesc sets the description of the message, which helps translators.
//
// Java: JsMessage.Builder.setDesc(String)
func (b *JsMessage_Builder) SetDesc(desc *string) *JsMessage_Builder {
	b.desc = desc
	return b
}

// SetMeaning sets the programmer-specified meaning of this message, which
// forces this message to translate differently.
//
// Java: JsMessage.Builder.setMeaning(String)
func (b *JsMessage_Builder) SetMeaning(meaning *string) *JsMessage_Builder {
	b.meaning = meaning
	return b
}

// SetAlternateId sets the alternate message ID, to be used if the primary
// ID is not yet translated.
//
// Java: JsMessage.Builder.setAlternateId(String)
func (b *JsMessage_Builder) SetAlternateId(alternateId *string) *JsMessage_Builder {
	b.alternateId = alternateId
	return b
}

// Java: JsMessage.Builder.setIsAnonymous(boolean)
func (b *JsMessage_Builder) SetIsAnonymous(isAnonymous bool) *JsMessage_Builder {
	b.isAnonymous = isAnonymous
	return b
}

// Java: JsMessage.Builder.setId(String)
func (b *JsMessage_Builder) SetId(id string) *JsMessage_Builder {
	if b.id != nil {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat(
			"id already set to '%s': cannot change it to '%s'", *b.id, id)))
	}
	b.id = &id
	return b
}

// Java: JsMessage.Builder.setIsExternalMsg(boolean)
func (b *JsMessage_Builder) SetIsExternalMsg(isExternalMsg bool) *JsMessage_Builder {
	b.isExternal = isExternalMsg
	return b
}

// HasParts gets whether at least one part has been appended.
//
// Java: JsMessage.Builder.hasParts()
func (b *JsMessage_Builder) HasParts() bool {
	return !b.parts.IsEmpty()
}

// HasGenderedVariants gets whether the message has gendered variants.
//
// Java: JsMessage.Builder.hasGenderedVariants()
func (b *JsMessage_Builder) HasGenderedVariants() bool {
	return !b.genderedMessageMap.IsEmpty()
}

// Java: JsMessage.Builder.getParts()
func (b *JsMessage_Builder) GetParts() *jutil.ArrayList[JsMessage_Part] {
	return b.parts
}

// Java: JsMessage.Builder.build()
func (b *JsMessage_Builder) Build() *JsMessage {
	if b.key == nil {
		panic(jutil.NewNullPointerExceptionMsg("key has not been set"))
	}
	if b.id == nil {
		panic(jutil.NewNullPointerExceptionMsg("id has not been set"))
	}
	jutil.CheckStateMsg(!b.isExternal || !b.isAnonymous, "a message cannot be both anonymous and external")

	// An alternate ID that points to itself is a no-op, so just omit it.
	if (b.alternateId != nil) && *b.alternateId == *b.id {
		b.alternateId = nil
	}

	immutableGenderedPartsBuilder := jutil.NewImmutableMapBuilder[JsMessage_GrammaticalGenderCase, []JsMessage_Part]()
	for key, value := range b.genderedMessageMap.All() {
		immutableGenderedPartsBuilder.Put(key, value.ToSlice())
	}
	immutableGenderedMessageMap := immutableGenderedPartsBuilder.BuildOrThrow()

	return NewJsMessage(
		b.sourceName,
		*b.key,
		b.isAnonymous,
		b.isExternal,
		*b.id,
		b.parts.ToSlice(),
		immutableGenderedMessageMap,
		b.alternateId,
		b.desc,
		b.meaning,
		b.placeholderNameToExampleMap,
		b.placeholderNameToOriginalCodeMap,
		jutil.NewLinkedHashSet(b.jsPlaceholderNames.ToSlice()...),
		jutil.NewLinkedHashSet(b.canonicalPlaceholderNames.ToSlice()...))
}

// JsMessage_Hash contains routines for hashing.
//
// The hash takes a byte array representing arbitrary data (a number,
// String, or Object) and turns it into a small, hopefully unique, number.
// There are additional convenience functions which hash int, long, and
// String types.
//
// Note: this hash has weaknesses in the two most-significant key bits and
// in the three least-significant seed bits. The weaknesses are small and
// practically speaking, will not affect the distribution of hash values.
// Still, it would be good practice not to choose seeds 0, 1, 2, 3, ..., n
// to yield n, independent hash functions. Use pseudo-random seeds instead.
//
// This code is based on the work of Craig Silverstein and Sanjay Ghemawat
// in, then forked from com.google.common.
//
// The original code for the hash function is courtesy Bob Jenkins
// (http://burtleburtle.net/bob/hash/evahash.html).
//
// TODO(anatol): Add stream hashing functionality.
type JsMessage_Hash struct{}

const (
	// Default hash seed (64 bit)
	jsMessage_Hash_SEED64 int64 = 0x2b992ddfa23249d6 // part of pi, arbitrary

	// Hash constant (64 bit)
	jsMessage_Hash_CONSTANT64 int64 = -0x1f73e299748a907e // 0xe08c1d668b756f82L; part of golden ratio, arbitrary
)

// JsMessage_Hash_Hash64 hashes a string to a 64 bit value. The digits of
// pi are used for the hash seed.
//
// Java: JsMessage.Hash.hash64(String)
func JsMessage_Hash_Hash64(value *string) int64 {
	return jsMessage_Hash_hash64WithSeed(value, jsMessage_Hash_SEED64)
}

// jsMessage_Hash_hash64WithSeed hashes a string to a 64 bit value using
// the supplied seed.
//
// Java: JsMessage.Hash.hash64(String, long)
func jsMessage_Hash_hash64WithSeed(value *string, seed int64) int64 {
	if value == nil {
		return jsMessage_Hash_hash64Range(nil, 0, 0, seed)
	}
	return jsMessage_Hash_hash64Bytes(jutil.StringGetBytesUTF8(*value), seed)
}

// jsMessage_Hash_hash64Bytes hashes byte array to a 64 bit value using the
// supplied seed.
//
// Java: JsMessage.Hash.hash64(byte[], long)
func jsMessage_Hash_hash64Bytes(value []byte, seed int64) int64 {
	length := 0
	if value != nil {
		length = len(value)
	}
	return jsMessage_Hash_hash64Range(value, 0, length, seed)
}

// jsMessage_Hash_hash64Range hashes byte array to a 64 bit value using the
// supplied seed.
//
// value is the bytes to hash; offset the starting position of value where
// bytes are used for the hash computation; length the number of bytes of
// value that are used for the hash computation; seed the seed.
//
// Java: JsMessage.Hash.hash64(byte[], int, int, long)
func jsMessage_Hash_hash64Range(value []byte, offset int, length int, seed int64) int64 {
	a := jsMessage_Hash_CONSTANT64
	b := a
	c := seed
	var keylen int

	ff := func(i int) int64 { return int64(value[i]) & 0xff }

	for keylen = length; keylen >= 24; keylen, offset = keylen-24, offset+24 {
		a += jsMessage_Hash_word64At(value, offset)
		b += jsMessage_Hash_word64At(value, offset+8)
		c += jsMessage_Hash_word64At(value, offset+16)

		// Mix
		a -= b
		a -= c
		a ^= jutil.Ushr64(c, 43)
		b -= c
		b -= a
		b ^= a << 9
		c -= a
		c -= b
		c ^= jutil.Ushr64(b, 8)
		a -= b
		a -= c
		a ^= jutil.Ushr64(c, 38)
		b -= c
		b -= a
		b ^= a << 23
		c -= a
		c -= b
		c ^= jutil.Ushr64(b, 5)
		a -= b
		a -= c
		a ^= jutil.Ushr64(c, 35)
		b -= c
		b -= a
		b ^= a << 49
		c -= a
		c -= b
		c ^= jutil.Ushr64(b, 11)
		a -= b
		a -= c
		a ^= jutil.Ushr64(c, 12)
		b -= c
		b -= a
		b ^= a << 18
		c -= a
		c -= b
		c ^= jutil.Ushr64(b, 22)
	}

	c += int64(length)
	if keylen >= 16 {
		if keylen == 23 {
			c += int64(int8(value[offset+22])) << 56
		}
		if keylen >= 22 {
			c += ff(offset+21) << 48
		}
		if keylen >= 21 {
			c += ff(offset+20) << 40
		}
		if keylen >= 20 {
			c += ff(offset+19) << 32
		}
		if keylen >= 19 {
			c += ff(offset+18) << 24
		}
		if keylen >= 18 {
			c += ff(offset+17) << 16
		}
		if keylen >= 17 {
			c += ff(offset+16) << 8
			// the first byte of c is reserved for the length
		}
		if keylen >= 16 {
			b += jsMessage_Hash_word64At(value, offset+8)
			a += jsMessage_Hash_word64At(value, offset)
		}
	} else if keylen >= 8 {
		if keylen == 15 {
			b += ff(offset+14) << 48
		}
		if keylen >= 14 {
			b += ff(offset+13) << 40
		}
		if keylen >= 13 {
			b += ff(offset+12) << 32
		}
		if keylen >= 12 {
			b += ff(offset+11) << 24
		}
		if keylen >= 11 {
			b += ff(offset+10) << 16
		}
		if keylen >= 10 {
			b += ff(offset+9) << 8
		}
		if keylen >= 9 {
			b += ff(offset + 8)
		}
		if keylen >= 8 {
			a += jsMessage_Hash_word64At(value, offset)
		}
	} else {
		if keylen == 7 {
			a += ff(offset+6) << 48
		}
		if keylen >= 6 {
			a += ff(offset+5) << 40
		}
		if keylen >= 5 {
			a += ff(offset+4) << 32
		}
		if keylen >= 4 {
			a += ff(offset+3) << 24
		}
		if keylen >= 3 {
			a += ff(offset+2) << 16
		}
		if keylen >= 2 {
			a += ff(offset+1) << 8
		}
		if keylen >= 1 {
			a += ff(offset + 0)
			// case 0: nothing left to add
		}
	}
	return jsMessage_Hash_mix64(a, b, c)
}

// Java: JsMessage.Hash.word64At(byte[], int)
func jsMessage_Hash_word64At(bytes []byte, offset int) int64 {
	return (int64(bytes[offset+0]) & 0xff) +
		((int64(bytes[offset+1]) & 0xff) << 8) +
		((int64(bytes[offset+2]) & 0xff) << 16) +
		((int64(bytes[offset+3]) & 0xff) << 24) +
		((int64(bytes[offset+4]) & 0xff) << 32) +
		((int64(bytes[offset+5]) & 0xff) << 40) +
		((int64(bytes[offset+6]) & 0xff) << 48) +
		((int64(bytes[offset+7]) & 0xff) << 56)
}

// jsMessage_Hash_mix64 mixes longs a, b, and c, and returns the final value
// of c.
//
// Java: JsMessage.Hash.mix64(long, long, long)
func jsMessage_Hash_mix64(a int64, b int64, c int64) int64 {
	a -= b
	a -= c
	a ^= jutil.Ushr64(c, 43)
	b -= c
	b -= a
	b ^= a << 9
	c -= a
	c -= b
	c ^= jutil.Ushr64(b, 8)
	a -= b
	a -= c
	a ^= jutil.Ushr64(c, 38)
	b -= c
	b -= a
	b ^= a << 23
	c -= a
	c -= b
	c ^= jutil.Ushr64(b, 5)
	a -= b
	a -= c
	a ^= jutil.Ushr64(c, 35)
	b -= c
	b -= a
	b ^= a << 49
	c -= a
	c -= b
	c ^= jutil.Ushr64(b, 11)
	a -= b
	a -= c
	a ^= jutil.Ushr64(c, 12)
	b -= c
	b -= a
	b ^= a << 18
	c -= a
	c -= b
	c ^= jutil.Ushr64(b, 22)
	return c
}

// JsMessage_IdGenerator is the ID generator.
type JsMessage_IdGenerator interface {
	// GenerateId generates the ID for the message. Messages with the same
	// messageParts and meaning will get the same id. Messages with the same
	// id will get the same translation.
	//
	// meaning is the programmer-specified meaning. If no @meaning
	// annotation appears, we will use the name of the variable it's
	// assigned to. If the variable is unnamed, then we will just use a
	// fingerprint of the message. messageParts are the parts of the
	// message, including the main message text.
	//
	// Java: JsMessage.IdGenerator.generateId(String, List<Part>)
	GenerateId(meaning *string, messageParts []JsMessage_Part) string
}

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
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	jsMessageEmptyGenderMap jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part]
	jsMessageEmptyStringMap jutil.LinkedHashMap[string, string]
	jsMessageEmptyStringSet jutil.LinkedHashSet[string]
	jsMessageEmptyNodeMap   jutil.LinkedHashMap[string, rhino.Node]
)

// JsMessageDefinition is the container struct that holds a JsMessage and
// information about the `goog.getMsg()` call it was built from.
type JsMessageDefinition struct {
	message               *JsMessage
	placeholderValueMap   *jutil.LinkedHashMap[string, rhino.Node]
	messageNode           rhino.Node
	templateTextNode      rhino.Node
	placeholderValuesNode rhino.Node
	escapeLessThan        bool
	unescapeHtmlEntities  bool
}

// GetMessage returns the JsMessage built from the `goog.getMsg()` call.
func (d JsMessageDefinition) GetMessage() *JsMessage { return d.message }

// GetMessageNode returns the RHS node of the message assignment statement.
func (d JsMessageDefinition) GetMessageNode() rhino.Node { return d.messageNode }

// GetTemplateTextNode returns the Node representing the message text template.
func (d JsMessageDefinition) GetTemplateTextNode() rhino.Node { return d.templateTextNode }

// GetPlaceholderValuesNode returns the object literal Node that maps
// placeholder names to expressions providing their values, or NilNode.
func (d JsMessageDefinition) GetPlaceholderValuesNode() rhino.Node { return d.placeholderValuesNode }

// GetPlaceholderValueMap returns a map from placeholder name to the Node
// assigned to it in the values map argument of `goog.getMsg()`.
func (d JsMessageDefinition) GetPlaceholderValueMap() *jutil.LinkedHashMap[string, rhino.Node] {
	if d.placeholderValueMap == nil {
		return &jsMessageEmptyNodeMap
	}
	return d.placeholderValueMap
}

// ShouldEscapeLessThan returns the value of the 'html' options key in the
// options bag argument.
func (d JsMessageDefinition) ShouldEscapeLessThan() bool { return d.escapeLessThan }

// ShouldUnescapeHtmlEntities returns the value of the 'unescapeHtmlEntities'
// options key in the options bag argument.
func (d JsMessageDefinition) ShouldUnescapeHtmlEntities() bool { return d.unescapeHtmlEntities }

// IcuTemplateDefinition contains a JsMessage representing an ICU template
// message and information about the
// `goog.i18n.messages.declareIcuTemplate()` call from which it was extracted.
type IcuTemplateDefinition struct {
	message          *JsMessage
	messageNode      rhino.Node
	templateTextNode rhino.Node
}

// GetMessage returns the JsMessage built from the `declareIcuTemplate()` call.
func (d IcuTemplateDefinition) GetMessage() *JsMessage { return d.message }

// GetMessageNode returns the RHS node of the message assignment statement.
func (d IcuTemplateDefinition) GetMessageNode() rhino.Node { return d.messageNode }

// GetTemplateTextNode returns the Node representing the message text template.
func (d IcuTemplateDefinition) GetTemplateTextNode() rhino.Node { return d.templateTextNode }

// JsMessage is a representation of a translatable message in JavaScript
// source code.
//
// Instances are created using a JsMessage_Builder. Optional maps and sets
// remain nil when empty to avoid eager collection allocations per message.
type JsMessage struct {
	sourceName                       string
	key                              string
	id                               string
	alternateId                      string
	desc                             string
	meaning                          string
	parts                            []JsMessage_Part
	genderedMessagesMap              *jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part]
	placeholderNameToExampleMap      *jutil.LinkedHashMap[string, string]
	placeholderNameToOriginalCodeMap *jutil.LinkedHashMap[string, string]
	jsPlaceholderNames               *jutil.LinkedHashSet[string]
	canonicalPlaceholderNames        *jutil.LinkedHashSet[string]
	hasSourceName                    bool
	hasAlternateId                   bool
	hasDesc                          bool
	hasMeaning                       bool
	isAnonymous                      bool
	isExternal                       bool
}

// Java: JsMessage.getSourceName()
func (m *JsMessage) GetSourceName() *string {
	if !m.hasSourceName {
		return nil
	}
	return &m.sourceName
}

// Java: JsMessage.getKey()
func (m *JsMessage) GetKey() string { return m.key }

// Java: JsMessage.isAnonymous()
func (m *JsMessage) IsAnonymous() bool { return m.isAnonymous }

// Java: JsMessage.isExternal()
func (m *JsMessage) IsExternal() bool { return m.isExternal }

// Java: JsMessage.getId()
func (m *JsMessage) GetId() string { return m.id }

// Java: JsMessage.getParts()
func (m *JsMessage) GetParts() []JsMessage_Part { return m.parts }

// Java: JsMessage.getGenderedMessagesMap()
func (m *JsMessage) GetGenderedMessagesMap() *jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part] {
	if m.genderedMessagesMap == nil {
		return &jsMessageEmptyGenderMap
	}
	return m.genderedMessagesMap
}

// HasGenderedVariants reports whether the message has gendered variants.
func (m *JsMessage) HasGenderedVariants() bool {
	return !m.genderedMessagesMap.IsEmpty()
}

// Java: JsMessage.getAlternateId()
func (m *JsMessage) GetAlternateId() *string {
	if !m.hasAlternateId {
		return nil
	}
	return &m.alternateId
}

// AlternateIdOpt returns the alternate ID and whether one is present.
func (m *JsMessage) AlternateIdOpt() (string, bool) { return m.alternateId, m.hasAlternateId }

// Java: JsMessage.getDesc()
func (m *JsMessage) GetDesc() *string {
	if !m.hasDesc {
		return nil
	}
	return &m.desc
}

// DescOpt returns the description and whether one is present.
func (m *JsMessage) DescOpt() (string, bool) { return m.desc, m.hasDesc }

// Java: JsMessage.getMeaning()
func (m *JsMessage) GetMeaning() *string {
	if !m.hasMeaning {
		return nil
	}
	return &m.meaning
}

// MeaningOpt returns the meaning and whether one is present.
func (m *JsMessage) MeaningOpt() (string, bool) { return m.meaning, m.hasMeaning }

// Java: JsMessage.getPlaceholderNameToExampleMap()
func (m *JsMessage) GetPlaceholderNameToExampleMap() *jutil.LinkedHashMap[string, string] {
	if m.placeholderNameToExampleMap == nil {
		return &jsMessageEmptyStringMap
	}
	return m.placeholderNameToExampleMap
}

// Java: JsMessage.getPlaceholderNameToOriginalCodeMap()
func (m *JsMessage) GetPlaceholderNameToOriginalCodeMap() *jutil.LinkedHashMap[string, string] {
	if m.placeholderNameToOriginalCodeMap == nil {
		return &jsMessageEmptyStringMap
	}
	return m.placeholderNameToOriginalCodeMap
}

// Java: JsMessage.jsPlaceholderNames()
func (m *JsMessage) JsPlaceholderNames() *jutil.LinkedHashSet[string] {
	if m.jsPlaceholderNames == nil {
		return &jsMessageEmptyStringSet
	}
	return m.jsPlaceholderNames
}

// Java: JsMessage.canonicalPlaceholderNames()
func (m *JsMessage) CanonicalPlaceholderNames() *jutil.LinkedHashSet[string] {
	if m.canonicalPlaceholderNames == nil {
		return &jsMessageEmptyStringSet
	}
	return m.canonicalPlaceholderNames
}

// Java: JsMessage.PH_JS_PREFIX
const JsMessage_PH_JS_PREFIX = "{$"

// Java: JsMessage.PH_JS_SUFFIX
const JsMessage_PH_JS_SUFFIX = "}"

// JsMessage_GrammaticalGenderCase is the enum for grammatical gender cases.
type JsMessage_GrammaticalGenderCase int32

const (
	JsMessage_GrammaticalGenderCase_MASCULINE JsMessage_GrammaticalGenderCase = iota
	JsMessage_GrammaticalGenderCase_FEMININE
	JsMessage_GrammaticalGenderCase_NEUTER
	JsMessage_GrammaticalGenderCase_OTHER
)

var jsMessage_GrammaticalGenderCase_names = [...]string{"MASCULINE", "FEMININE", "NEUTER", "OTHER"}

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

// Java: JsMessage.getPlaceholderOriginalCode(PlaceholderReference)
func (m *JsMessage) GetPlaceholderOriginalCode(placeholderReference JsMessage_Part) string {
	return m.placeholderNameToOriginalCodeMap.GetOrDefault(placeholderReference.StoredPlaceholderName(), "-")
}

// Java: JsMessage.getPlaceholderExample(PlaceholderReference)
func (m *JsMessage) GetPlaceholderExample(placeholderReference JsMessage_Part) string {
	return m.placeholderNameToExampleMap.GetOrDefault(placeholderReference.StoredPlaceholderName(), "-")
}

// GetGenderedMessageParts returns the parts for genderCase (nil if there are
// none for it).
//
// Java: JsMessage.getGenderedMessageParts(GrammaticalGenderCase)
func (m *JsMessage) GetGenderedMessageParts(genderCase JsMessage_GrammaticalGenderCase) []JsMessage_Part {
	if m.genderedMessagesMap.IsEmpty() {
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"Message does not contain grammatical gendered variants."))
	}
	return m.genderedMessagesMap.Get(genderCase)
}

// GetGenderedMessageVariants gets the list of grammatical gender cases for
// the message.
//
// Java: JsMessage.getGenderedMessageVariants()
func (m *JsMessage) GetGenderedMessageVariants() []JsMessage_GrammaticalGenderCase {
	return m.genderedMessagesMap.KeySlice()
}

// AsJsMessageString returns a single string representing the message in
// `goog.getMsg()` placeholder format ("{$jsPlaceholderName}").
//
// Java: JsMessage.asJsMessageString()
func (m *JsMessage) AsJsMessageString() string {
	if !m.genderedMessagesMap.IsEmpty() {
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"asJsMessageString() is not supported for messages with gendered variants. Please provide" +
				" a grammatical gender case as an argument."))
	}
	return jsMessage_formatPartsAsJs(m.parts)
}

// AsJsMessageStringForGender returns a single string representing the message
// for a given gender case in `goog.getMsg()` placeholder format.
//
// Java: JsMessage.asJsMessageString(GrammaticalGenderCase)
func (m *JsMessage) AsJsMessageStringForGender(genderCase JsMessage_GrammaticalGenderCase) string {
	parts, ok := m.genderedMessagesMap.Lookup(genderCase)
	if !ok {
		panic(jutil.NewIllegalArgumentExceptionMsg("No message for grammtical gender case: " + genderCase.String()))
	}
	return jsMessage_formatPartsAsJs(parts)
}

func jsMessage_formatPartsAsJs(parts []JsMessage_Part) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 && !parts[0].IsPlaceholder() {
		return parts[0].GetString()
	}
	var sb strings.Builder
	for _, p := range parts {
		if p.IsPlaceholder() {
			sb.WriteString(JsMessage_PH_JS_PREFIX)
			sb.WriteString(p.GetJsPlaceholderName())
			sb.WriteString(JsMessage_PH_JS_SUFFIX)
		} else {
			sb.WriteString(p.GetString())
		}
	}
	return sb.String()
}

// AsIcuMessageString returns a single string representing the message in ICU
// placeholder format ("{CANONICAL_PLACEHOLDER_NAME}").
//
// Java: JsMessage.asIcuMessageString()
func (m *JsMessage) AsIcuMessageString() string {
	if !m.genderedMessagesMap.IsEmpty() {
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"asIcuMessageString() is not supported for messages with gendered variants. Please" +
				" provide a grammatical gender case as an argument."))
	}
	return jsMessage_formatPartsAsIcu(m.parts)
}

// AsIcuMessageStringForGender returns a single string representing the message
// for a given gender case in ICU placeholder format.
//
// Java: JsMessage.asIcuMessageString(GrammaticalGenderCase)
func (m *JsMessage) AsIcuMessageStringForGender(genderCase JsMessage_GrammaticalGenderCase) string {
	parts, ok := m.genderedMessagesMap.Lookup(genderCase)
	if !ok {
		panic(jutil.NewIllegalArgumentExceptionMsg("No message for grammtical gender case: " + genderCase.String()))
	}
	return jsMessage_formatPartsAsIcu(parts)
}

func jsMessage_formatPartsAsIcu(parts []JsMessage_Part) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 && !parts[0].IsPlaceholder() {
		return parts[0].GetString()
	}
	var sb strings.Builder
	for _, p := range parts {
		if p.IsPlaceholder() {
			sb.WriteByte('{')
			sb.WriteString(p.GetCanonicalPlaceholderName())
			sb.WriteByte('}')
		} else {
			sb.WriteString(p.GetString())
		}
	}
	return sb.String()
}

// IsEmpty returns true iff the message is represented by empty string.
//
// Java: JsMessage.isEmpty()
func (m *JsMessage) IsEmpty() bool {
	if m.genderedMessagesMap.IsEmpty() {
		for _, part := range m.parts {
			if part.IsPlaceholder() || len(part.str) > 0 {
				return false
			}
		}
	} else {
		for parts := range m.genderedMessagesMap.Values() {
			for _, part := range parts {
				if part.IsPlaceholder() || len(part.str) > 0 {
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
	if m == that {
		return true
	}
	return m.hasSourceName == that.hasSourceName && m.sourceName == that.sourceName &&
		m.key == that.key &&
		m.isAnonymous == that.isAnonymous &&
		m.isExternal == that.isExternal &&
		m.id == that.id &&
		slices.Equal(m.parts, that.parts) &&
		jsMessage_genderedEqual(m.genderedMessagesMap, that.genderedMessagesMap) &&
		m.hasAlternateId == that.hasAlternateId && m.alternateId == that.alternateId &&
		m.hasDesc == that.hasDesc && m.desc == that.desc &&
		m.hasMeaning == that.hasMeaning && m.meaning == that.meaning &&
		jsMessage_stringMapEqual(m.placeholderNameToExampleMap, that.placeholderNameToExampleMap) &&
		jsMessage_stringMapEqual(m.placeholderNameToOriginalCodeMap, that.placeholderNameToOriginalCodeMap) &&
		jsMessage_stringSetEqual(m.jsPlaceholderNames, that.jsPlaceholderNames) &&
		jsMessage_stringSetEqual(m.canonicalPlaceholderNames, that.canonicalPlaceholderNames)
}

func jsMessage_genderedEqual(a, b *jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part]) bool {
	if a.Size() != b.Size() {
		return false
	}
	for k, v := range a.All() {
		w, ok := b.Lookup(k)
		if !ok || !slices.Equal(v, w) {
			return false
		}
	}
	return true
}

func jsMessage_stringMapEqual(a, b *jutil.LinkedHashMap[string, string]) bool {
	if a.Size() != b.Size() {
		return false
	}
	for k, v := range a.All() {
		w, ok := b.Lookup(k)
		if !ok || v != w {
			return false
		}
	}
	return true
}

func jsMessage_stringSetEqual(a, b *jutil.LinkedHashSet[string]) bool {
	if a.Size() != b.Size() {
		return false
	}
	for k := range a.All() {
		if !b.Contains(k) {
			return false
		}
	}
	return true
}

func jsMessage_partsHashCode(parts []JsMessage_Part) int32 {
	var h int32 = 1
	for _, p := range parts {
		h = 31*h + p.HashCode()
	}
	return h
}

func jsMessage_boolHashCode(b bool) int32 {
	if b {
		return 1231
	}
	return 1237
}

func jsMessage_optStringHashCode(s string, ok bool) int32 {
	if !ok {
		return 0
	}
	return jutil.StringHashCode(s)
}

func jsMessage_optStringFormat(s string, ok bool) string {
	if !ok {
		return "null"
	}
	return s
}

func jsMessage_stringMapHashCode(m *jutil.LinkedHashMap[string, string]) int32 {
	var h int32
	for k, v := range m.All() {
		h += jutil.StringHashCode(k) ^ jutil.StringHashCode(v)
	}
	return h
}

func jsMessage_stringSetHashCode(s *jutil.LinkedHashSet[string]) int32 {
	var h int32
	for k := range s.All() {
		h += jutil.StringHashCode(k)
	}
	return h
}

// HashCode is the record's hashCode: h = 31*h + hash(component), from 0.
//
// Java: JsMessage.hashCode()
func (m *JsMessage) HashCode() int32 {
	var h int32 = 0
	h = 31*h + jsMessage_optStringHashCode(m.sourceName, m.hasSourceName)
	h = 31*h + jutil.StringHashCode(m.key)
	h = 31*h + jsMessage_boolHashCode(m.isAnonymous)
	h = 31*h + jsMessage_boolHashCode(m.isExternal)
	h = 31*h + jutil.StringHashCode(m.id)
	h = 31*h + jsMessage_partsHashCode(m.parts)
	var gendered int32
	for k, v := range m.genderedMessagesMap.All() {
		gendered += int32(k) ^ jsMessage_partsHashCode(v)
	}
	h = 31*h + gendered
	h = 31*h + jsMessage_optStringHashCode(m.alternateId, m.hasAlternateId)
	h = 31*h + jsMessage_optStringHashCode(m.desc, m.hasDesc)
	h = 31*h + jsMessage_optStringHashCode(m.meaning, m.hasMeaning)
	h = 31*h + jsMessage_stringMapHashCode(m.placeholderNameToExampleMap)
	h = 31*h + jsMessage_stringMapHashCode(m.placeholderNameToOriginalCodeMap)
	h = 31*h + jsMessage_stringSetHashCode(m.jsPlaceholderNames)
	h = 31*h + jsMessage_stringSetHashCode(m.canonicalPlaceholderNames)
	return h
}

// String is the record's toString.
//
// Java: JsMessage.toString()
func (m *JsMessage) String() string {
	var gendered string
	if m.genderedMessagesMap.IsEmpty() {
		gendered = "{}"
	} else {
		var sb strings.Builder
		sb.WriteByte('{')
		first := true
		for k, v := range m.genderedMessagesMap.All() {
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
		"getSourceName=" + jsMessage_optStringFormat(m.sourceName, m.hasSourceName) +
		", getKey=" + m.key +
		", isAnonymous=" + strconv.FormatBool(m.isAnonymous) +
		", isExternal=" + strconv.FormatBool(m.isExternal) +
		", getId=" + m.id +
		", getParts=" + jutil.CollectionToString(m.parts) +
		", getGenderedMessagesMap=" + gendered +
		", getAlternateId=" + jsMessage_optStringFormat(m.alternateId, m.hasAlternateId) +
		", getDesc=" + jsMessage_optStringFormat(m.desc, m.hasDesc) +
		", getMeaning=" + jsMessage_optStringFormat(m.meaning, m.hasMeaning) +
		", getPlaceholderNameToExampleMap=" + m.GetPlaceholderNameToExampleMap().String() +
		", getPlaceholderNameToOriginalCodeMap=" + m.GetPlaceholderNameToOriginalCodeMap().String() +
		", jsPlaceholderNames=" + jutil.CollectionToString(m.jsPlaceholderNames.ToSlice()) +
		", canonicalPlaceholderNames=" + jutil.CollectionToString(m.canonicalPlaceholderNames.ToSlice()) +
		"]"
}

// JsMessage_Part represents a piece of a JsMessage: either a literal string
// part (!isPlaceholder) or a placeholder reference (isPlaceholder, stored in
// canonical UPPER_UNDERSCORE format when canonical is true, or JS
// lowerCamelCase format when canonical is false).
type JsMessage_Part struct {
	str           string
	canonical     bool
	isPlaceholder bool
}

// JsMessage_StringPart_Create constructs a literal string message part.
//
// Java: JsMessage.StringPart.create(String)
func JsMessage_StringPart_Create(str string) JsMessage_Part {
	return JsMessage_Part{str: str}
}

// JsMessage_PlaceholderReference_CreateForJsName constructs a placeholder
// reference from a JS placeholder name (lowerCamelCase with optional numeric
// suffixes).
//
// Java: JsMessage.PlaceholderReference.createForJsName(String)
func JsMessage_PlaceholderReference_CreateForJsName(name string) JsMessage_Part {
	jutil.CheckArgumentMsg1(
		JsMessage_IsLowerCamelCaseWithNumericSuffixes(name),
		"invalid JS placeholder name format: '%s'",
		name)
	return JsMessage_Part{str: name, isPlaceholder: true}
}

// JsMessage_PlaceholderReference_CreateForCanonicalName constructs a
// placeholder reference from a canonical placeholder name (UPPER_SNAKE_CASE).
//
// Java: JsMessage.PlaceholderReference.createForCanonicalName(String)
func JsMessage_PlaceholderReference_CreateForCanonicalName(name string) JsMessage_Part {
	jutil.CheckArgumentMsg1(
		JsMessage_IsCanonicalPlaceholderNameFormat(name),
		"not a canonical placeholder name format: '%s'",
		name)
	return JsMessage_Part{str: name, canonical: true, isPlaceholder: true}
}

// IsPlaceholder reports whether this part is a placeholder reference.
//
// Java: JsMessage.Part.isPlaceholder()
func (p JsMessage_Part) IsPlaceholder() bool { return p.isPlaceholder }

// GetString returns the literal string for a non-placeholder part.
//
// Java: JsMessage.Part.getString()
func (p JsMessage_Part) GetString() string {
	if p.isPlaceholder {
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			jutil.StringFormat("not a string part: '%s'", p.GetJsPlaceholderName())))
	}
	return p.str
}

// StoredPlaceholderName returns the raw stored placeholder name.
//
// Java: JsMessage.PlaceholderReference.storedPlaceholderName()
func (p JsMessage_Part) StoredPlaceholderName() string {
	if !p.isPlaceholder {
		panic(jutil.NewUnsupportedOperationExceptionMsg(jutil.StringFormat("not a placeholder: '%s'", p.str)))
	}
	return p.str
}

// CanonicalFormat reports whether a placeholder reference is stored in
// canonical format.
//
// Java: JsMessage.PlaceholderReference.canonicalFormat()
func (p JsMessage_Part) CanonicalFormat() bool { return p.canonical }

// GetJsPlaceholderName gets the name of the placeholder as it would appear in
// JS code (lowerCamelCase with optional _123_456 suffixes).
//
// Java: JsMessage.Part.getJsPlaceholderName()
func (p JsMessage_Part) GetJsPlaceholderName() string {
	if !p.isPlaceholder {
		panic(jutil.NewUnsupportedOperationExceptionMsg(jutil.StringFormat("not a placeholder: '%s'", p.str)))
	}
	if p.canonical {
		return JsMessage_ToLowerCamelCaseWithNumericSuffixes(p.str)
	}
	return p.str
}

// GetCanonicalPlaceholderName gets the name of the placeholder in XMB/XTB
// UPPER_SNAKE_CASE format.
//
// Java: JsMessage.Part.getCanonicalPlaceholderName()
func (p JsMessage_Part) GetCanonicalPlaceholderName() string {
	if !p.isPlaceholder {
		panic(jutil.NewUnsupportedOperationExceptionMsg(jutil.StringFormat("not a placeholder: '%s'", p.str)))
	}
	if p.canonical {
		return p.str
	}
	return jutil.CaseFormat_LowerCamelToUpperUnderscore(p.str)
}

// Equals compares two JsMessage_Part values.
func (p JsMessage_Part) Equals(o any) bool {
	that, ok := o.(JsMessage_Part)
	return ok && p == that
}

// HashCode computes the Java record hashCode for StringPart or PlaceholderReference.
func (p JsMessage_Part) HashCode() int32 {
	if !p.isPlaceholder {
		return jutil.StringHashCode(p.str)
	}
	return 31*jutil.StringHashCode(p.str) + jsMessage_boolHashCode(p.canonical)
}

// String renders the Java record toString for StringPart or PlaceholderReference.
func (p JsMessage_Part) String() string {
	if !p.isPlaceholder {
		return "StringPart[string=" + p.str + "]"
	}
	return "PlaceholderReference[storedPlaceholderName=" + p.str +
		", canonicalFormat=" + strconv.FormatBool(p.canonical) + "]"
}

// JsMessage_IsLowerCamelCaseWithNumericSuffixes returns whether a string is
// nonempty, begins with a lowercase letter, and contains only digits and
// underscores after the first underscore (`\A[a-z][a-zA-Z0-9]*[_0-9]*\z`).
//
// Java: JsMessage.isLowerCamelCaseWithNumericSuffixes(String)
func JsMessage_IsLowerCamelCaseWithNumericSuffixes(input string) bool {
	if len(input) == 0 || input[0] < 'a' || input[0] > 'z' {
		return false
	}
	i := 1
	for i < len(input) {
		c := input[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			i++
		} else {
			break
		}
	}
	for ; i < len(input); i++ {
		c := input[i]
		if c != '_' && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

// JsMessage_ToLowerCamelCaseWithNumericSuffixes converts the given string
// from upper-underscore case to lower-camel case, preserving numeric
// suffixes. For example: "NAME" -> "name", "A4_LETTER" -> "a4Letter",
// "START_SPAN_1_23" -> "startSpan_1_23".
//
// Java: JsMessage.toLowerCamelCaseWithNumericSuffixes(String)
func JsMessage_ToLowerCamelCaseWithNumericSuffixes(input string) string {
	// Determine where the numeric suffixes begin by scanning backwards.
	// Surrogate code units are never Character.isDigit, so BMP UTF-8 decoding
	// finds the exact same suffix boundary as UTF-16 code-unit scanning with
	// zero intermediate slice allocations.
	suffixStart := len(input)
	if jutil.IsASCII(input) {
		for suffixStart > 0 {
			numberStart := suffixStart
			for numberStart > 0 && input[numberStart-1] >= '0' && input[numberStart-1] <= '9' {
				numberStart--
			}
			if numberStart > 0 && numberStart < suffixStart && input[numberStart-1] == '_' {
				suffixStart = numberStart - 1
			} else {
				break
			}
		}
	} else {
		for suffixStart > 0 {
			numberStart := suffixStart
			var lastRune rune
			var lastSize int
			for numberStart > 0 {
				r, size := utf8.DecodeLastRuneInString(input[:numberStart])
				if r < 0x10000 && (r != utf8.RuneError || size != 1) && jutil.CharacterIsDigit(uint16(r)) {
					numberStart -= size
				} else {
					lastRune = r
					lastSize = size
					break
				}
			}
			if numberStart > 0 && numberStart < suffixStart && lastSize == 1 && lastRune == '_' {
				suffixStart = numberStart - 1
			} else {
				break
			}
		}
	}

	if suffixStart == len(input) {
		return jutil.CaseFormat_UpperUnderscoreToLowerCamel(input)
	}
	return jutil.CaseFormat_UpperUnderscoreToLowerCamel(input[:suffixStart]) + input[suffixStart:]
}

// JsMessage_IsCanonicalPlaceholderNameFormat reports whether name matches
// `\A[A-Z0-9_]*\z` (the canonical format for placeholder names in XTB/XMB).
//
// Java: JsMessage.isCanonicalPlaceholderNameFormat(String)
func JsMessage_IsCanonicalPlaceholderNameFormat(name string) bool {
	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			return false
		}
	}
	return true
}

// JsMessage_Builder contains functionality for creating JS messages.
type JsMessage_Builder struct {
	key         string
	meaning     string
	desc        string
	id          string
	alternateId string
	sourceName  string

	parts                            []JsMessage_Part
	genderedMessageMap               *jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part]
	jsPlaceholderNames               *jutil.LinkedHashSet[string]
	canonicalPlaceholderNames        *jutil.LinkedHashSet[string]
	placeholderNameToExampleMap      *jutil.LinkedHashMap[string, string]
	placeholderNameToOriginalCodeMap *jutil.LinkedHashMap[string, string]

	hasKey         bool
	hasMeaning     bool
	hasDesc        bool
	hasId          bool
	hasAlternateId bool
	hasSourceName  bool
	isAnonymous    bool
	isExternal     bool
}

// NewJsMessage_Builder is new JsMessage.Builder().
func NewJsMessage_Builder() *JsMessage_Builder {
	return &JsMessage_Builder{}
}

func (b *JsMessage_Builder) addJsPlaceholderName(name string) {
	if b.jsPlaceholderNames == nil {
		b.jsPlaceholderNames = jutil.NewLinkedHashSet[string]()
	}
	b.jsPlaceholderNames.Add(name)
}

func (b *JsMessage_Builder) addCanonicalPlaceholderName(name string) {
	if b.canonicalPlaceholderNames == nil {
		b.canonicalPlaceholderNames = jutil.NewLinkedHashSet[string]()
	}
	b.canonicalPlaceholderNames.Add(name)
}

// GetKey gets the message's key (e.g. "MSG_HELLO").
//
// Java: JsMessage.Builder.getKey()
func (b *JsMessage_Builder) GetKey() *string {
	if !b.hasKey {
		return nil
	}
	return &b.key
}

// SetKey sets a key that uniquely identifies this message (e.g. "MSG_HELLO").
//
// Java: JsMessage.Builder.setKey(String)
func (b *JsMessage_Builder) SetKey(key string) *JsMessage_Builder {
	b.key = key
	b.hasKey = true
	return b
}

// SetSourceName sets the message's sourceName.
//
// Java: JsMessage.Builder.setSourceName(String)
func (b *JsMessage_Builder) SetSourceName(sourceName *string) *JsMessage_Builder {
	if sourceName == nil {
		b.sourceName = ""
		b.hasSourceName = false
	} else {
		b.sourceName = *sourceName
		b.hasSourceName = true
	}
	return b
}

// SetSourceNameString sets a non-null sourceName without pointer boxing.
func (b *JsMessage_Builder) SetSourceNameString(sourceName string) *JsMessage_Builder {
	b.sourceName = sourceName
	b.hasSourceName = true
	return b
}

// Java: JsMessage.Builder.appendPart(Part)
func (b *JsMessage_Builder) AppendPart(part JsMessage_Part) *JsMessage_Builder {
	b.parts = append(b.parts, part)
	if part.IsPlaceholder() {
		b.addJsPlaceholderName(part.GetJsPlaceholderName())
		b.addCanonicalPlaceholderName(part.GetCanonicalPlaceholderName())
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

// appendGenderedPart appends to genderedMessageMap.get(genderCase), which must
// exist (Java throws NullPointerException on the null list).
func (b *JsMessage_Builder) appendGenderedPart(genderCase JsMessage_GrammaticalGenderCase, part JsMessage_Part) {
	parts, ok := b.genderedMessageMap.Lookup(genderCase)
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	b.genderedMessageMap.Put(genderCase, append(parts, part))
}

// AppendJsPlaceholderReference appends a placeholder reference to the message.
//
// Java: JsMessage.Builder.appendJsPlaceholderReference(String)
func (b *JsMessage_Builder) AppendJsPlaceholderReference(name string) *JsMessage_Builder {
	b.parts = append(b.parts, JsMessage_PlaceholderReference_CreateForJsName(name))
	b.addJsPlaceholderName(name)
	return b
}

// AppendJsPlaceholderReferenceForGender appends a placeholder reference to the
// corresponding gendered message.
//
// Java: JsMessage.Builder.appendJsPlaceholderReference(GrammaticalGenderCase, String)
func (b *JsMessage_Builder) AppendJsPlaceholderReferenceForGender(genderCase JsMessage_GrammaticalGenderCase, name string) *JsMessage_Builder {
	b.appendGenderedPart(genderCase, JsMessage_PlaceholderReference_CreateForJsName(name))
	b.addJsPlaceholderName(name)
	return b
}

// AppendCanonicalPlaceholderReference appends a canonical placeholder
// reference to the message.
//
// Java: JsMessage.Builder.appendCanonicalPlaceholderReference(String)
func (b *JsMessage_Builder) AppendCanonicalPlaceholderReference(name string) *JsMessage_Builder {
	placeholder := JsMessage_PlaceholderReference_CreateForCanonicalName(name)
	b.parts = append(b.parts, placeholder)
	b.addJsPlaceholderName(placeholder.GetJsPlaceholderName())
	return b
}

// AppendCanonicalPlaceholderReferenceForGender appends a canonical placeholder
// reference to the corresponding gendered message.
//
// Java: JsMessage.Builder.appendCanonicalPlaceholderReference(GrammaticalGenderCase, String)
func (b *JsMessage_Builder) AppendCanonicalPlaceholderReferenceForGender(
	genderCase JsMessage_GrammaticalGenderCase, name string) *JsMessage_Builder {
	placeholder := JsMessage_PlaceholderReference_CreateForCanonicalName(name)
	b.appendGenderedPart(genderCase, placeholder)
	b.addJsPlaceholderName(placeholder.GetJsPlaceholderName())
	return b
}

// AppendStringPart appends a translatable string literal to the message.
//
// Java: JsMessage.Builder.appendStringPart(String)
func (b *JsMessage_Builder) AppendStringPart(part string) *JsMessage_Builder {
	b.parts = append(b.parts, JsMessage_StringPart_Create(part))
	return b
}

// AppendStringPartForGender appends a translatable string literal to the
// corresponding gendered message.
//
// Java: JsMessage.Builder.appendStringPart(GrammaticalGenderCase, String)
func (b *JsMessage_Builder) AppendStringPartForGender(genderCase JsMessage_GrammaticalGenderCase, part string) *JsMessage_Builder {
	b.appendGenderedPart(genderCase, JsMessage_StringPart_Create(part))
	return b
}

// AddGenderedMessageKey adds a gendered message key to genderedMessageParts map.
//
// Java: JsMessage.Builder.addGenderedMessageKey(GrammaticalGenderCase)
func (b *JsMessage_Builder) AddGenderedMessageKey(part JsMessage_GrammaticalGenderCase) *JsMessage_Builder {
	if b.genderedMessageMap == nil {
		b.genderedMessageMap = jutil.NewLinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part]()
	}
	b.genderedMessageMap.Put(part, nil)
	return b
}

// Java: JsMessage.Builder.setPlaceholderNameToExampleMap(Map<String, String>)
func (b *JsMessage_Builder) SetPlaceholderNameToExampleMap(m *jutil.LinkedHashMap[string, string]) *JsMessage_Builder {
	if m.IsEmpty() {
		b.placeholderNameToExampleMap = nil
	} else {
		b.placeholderNameToExampleMap = jutil.NewLinkedHashMapFrom(m)
	}
	return b
}

// Java: JsMessage.Builder.setPlaceholderNameToOriginalCodeMap(Map<String, String>)
func (b *JsMessage_Builder) SetPlaceholderNameToOriginalCodeMap(m *jutil.LinkedHashMap[string, string]) *JsMessage_Builder {
	if m.IsEmpty() {
		b.placeholderNameToOriginalCodeMap = nil
	} else {
		b.placeholderNameToOriginalCodeMap = jutil.NewLinkedHashMapFrom(m)
	}
	return b
}

// SetDesc sets the description of the message, which helps translators.
//
// Java: JsMessage.Builder.setDesc(String)
func (b *JsMessage_Builder) SetDesc(desc *string) *JsMessage_Builder {
	if desc == nil {
		b.desc = ""
		b.hasDesc = false
	} else {
		b.desc = *desc
		b.hasDesc = true
	}
	return b
}

// SetDescString sets a non-null description without pointer boxing.
func (b *JsMessage_Builder) SetDescString(desc string) *JsMessage_Builder {
	b.desc = desc
	b.hasDesc = true
	return b
}

// SetMeaning sets the programmer-specified meaning of this message, which
// forces this message to translate differently.
//
// Java: JsMessage.Builder.setMeaning(String)
func (b *JsMessage_Builder) SetMeaning(meaning *string) *JsMessage_Builder {
	if meaning == nil {
		b.meaning = ""
		b.hasMeaning = false
	} else {
		b.meaning = *meaning
		b.hasMeaning = true
	}
	return b
}

// SetMeaningString sets a non-null meaning without pointer boxing.
func (b *JsMessage_Builder) SetMeaningString(meaning string) *JsMessage_Builder {
	b.meaning = meaning
	b.hasMeaning = true
	return b
}

// SetAlternateId sets the alternate message ID, to be used if the primary ID
// is not yet translated.
//
// Java: JsMessage.Builder.setAlternateId(String)
func (b *JsMessage_Builder) SetAlternateId(alternateId *string) *JsMessage_Builder {
	if alternateId == nil {
		b.alternateId = ""
		b.hasAlternateId = false
	} else {
		b.alternateId = *alternateId
		b.hasAlternateId = true
	}
	return b
}

// SetAlternateIdString sets a non-null alternate message ID without pointer boxing.
func (b *JsMessage_Builder) SetAlternateIdString(alternateId string) *JsMessage_Builder {
	b.alternateId = alternateId
	b.hasAlternateId = true
	return b
}

// Java: JsMessage.Builder.setIsAnonymous(boolean)
func (b *JsMessage_Builder) SetIsAnonymous(isAnonymous bool) *JsMessage_Builder {
	b.isAnonymous = isAnonymous
	return b
}

// Java: JsMessage.Builder.setId(String)
func (b *JsMessage_Builder) SetId(id string) *JsMessage_Builder {
	if b.hasId {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat(
			"id already set to '%s': cannot change it to '%s'", b.id, id)))
	}
	b.id = id
	b.hasId = true
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
	return len(b.parts) > 0
}

// HasGenderedVariants gets whether the message has gendered variants.
//
// Java: JsMessage.Builder.hasGenderedVariants()
func (b *JsMessage_Builder) HasGenderedVariants() bool {
	return !b.genderedMessageMap.IsEmpty()
}

// Java: JsMessage.Builder.getParts()
func (b *JsMessage_Builder) GetParts() []JsMessage_Part {
	return b.parts
}

// Java: JsMessage.Builder.build()
func (b *JsMessage_Builder) Build() *JsMessage {
	if !b.hasKey {
		panic(jutil.NewNullPointerExceptionMsg("key has not been set"))
	}
	if !b.hasId {
		panic(jutil.NewNullPointerExceptionMsg("id has not been set"))
	}
	jutil.CheckStateMsg(!b.isExternal || !b.isAnonymous, "a message cannot be both anonymous and external")

	// An alternate ID that points to itself is a no-op, so just omit it.
	hasAltId := b.hasAlternateId
	altId := b.alternateId
	if hasAltId && altId == b.id {
		hasAltId = false
		altId = ""
	}

	var immutableGenderedMessageMap *jutil.LinkedHashMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part]
	if !b.genderedMessageMap.IsEmpty() {
		immutableGenderedMessageMap = jutil.NewLinkedHashMapWithCapacity[JsMessage_GrammaticalGenderCase, []JsMessage_Part](b.genderedMessageMap.Size())
		for key, value := range b.genderedMessageMap.All() {
			immutableGenderedMessageMap.PutNew(key, slices.Clone(value))
		}
	}

	var exampleMap *jutil.LinkedHashMap[string, string]
	if !b.placeholderNameToExampleMap.IsEmpty() {
		exampleMap = jutil.NewLinkedHashMapFrom(b.placeholderNameToExampleMap)
	}

	var originalCodeMap *jutil.LinkedHashMap[string, string]
	if !b.placeholderNameToOriginalCodeMap.IsEmpty() {
		originalCodeMap = jutil.NewLinkedHashMapFrom(b.placeholderNameToOriginalCodeMap)
	}

	var jsPhNames *jutil.LinkedHashSet[string]
	if !b.jsPlaceholderNames.IsEmpty() {
		jsPhNames = jutil.NewLinkedHashSetWithCapacity[string](b.jsPlaceholderNames.Size())
		for name := range b.jsPlaceholderNames.All() {
			jsPhNames.AddNew(name)
		}
	}

	var canonicalPhNames *jutil.LinkedHashSet[string]
	if !b.canonicalPlaceholderNames.IsEmpty() {
		canonicalPhNames = jutil.NewLinkedHashSetWithCapacity[string](b.canonicalPlaceholderNames.Size())
		for name := range b.canonicalPlaceholderNames.All() {
			canonicalPhNames.AddNew(name)
		}
	}

	return &JsMessage{
		sourceName:                       b.sourceName,
		hasSourceName:                    b.hasSourceName,
		key:                              b.key,
		isAnonymous:                      b.isAnonymous,
		isExternal:                       b.isExternal,
		id:                               b.id,
		parts:                            slices.Clone(b.parts),
		genderedMessagesMap:              immutableGenderedMessageMap,
		alternateId:                      altId,
		hasAlternateId:                   hasAltId,
		desc:                             b.desc,
		hasDesc:                          b.hasDesc,
		meaning:                          b.meaning,
		hasMeaning:                       b.hasMeaning,
		placeholderNameToExampleMap:      exampleMap,
		placeholderNameToOriginalCodeMap: originalCodeMap,
		jsPlaceholderNames:               jsPhNames,
		canonicalPlaceholderNames:        canonicalPhNames,
	}
}

// JsMessage_Hash contains routines for hashing.
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
	if value == nil {
		return jsMessage_Hash_hash64Range(nil, 0, 0, jsMessage_Hash_SEED64)
	}
	return jsMessage_Hash_hash64String(*value)
}

func jsMessage_Hash_hash64String(value string) int64 {
	b := jutil.StringGetBytesUTF8(value)
	return jsMessage_Hash_hash64Range(b, 0, len(b), jsMessage_Hash_SEED64)
}

// jsMessage_Hash_hash64Range hashes byte array to a 64 bit value using the
// supplied seed.
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
	// Java: JsMessage.IdGenerator.generateId(String, List<Part>)
	GenerateId(meaning *string, messageParts []JsMessage_Part) string
}

// Port of com/google/javascript/jscomp/JsMessage.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/MessageBundle.java @ closure-compiler v20261021

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

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"iter"
	"maps"
	"slices"
	"strconv"
	"strings"

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
	placeholderValueMap   map[string]rhino.Node
	placeholderValueOrder []string
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
	if len(d.placeholderValueMap) == 0 {
		return &jsMessageEmptyNodeMap
	}
	out := jutil.NewLinkedHashMapWithCapacity[string, rhino.Node](len(d.placeholderValueOrder))
	for _, k := range d.placeholderValueOrder {
		out.PutNew(k, d.placeholderValueMap[k])
	}
	return out
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
	sourceName                         string
	key                                string
	id                                 string
	alternateId                        string
	desc                               string
	meaning                            string
	parts                              []JsMessage_Part
	genderedMessagesMap                map[JsMessage_GrammaticalGenderCase][]JsMessage_Part
	genderedMessagesOrder              []JsMessage_GrammaticalGenderCase
	placeholderNameToExampleMap        map[string]string
	placeholderNameToExampleOrder      []string
	placeholderNameToOriginalCodeMap   map[string]string
	placeholderNameToOriginalCodeOrder []string
	jsPlaceholderNames                 map[string]struct{}
	jsPlaceholderNamesOrder            []string
	canonicalPlaceholderNames          map[string]struct{}
	canonicalPlaceholderNamesOrder     []string
	hasSourceName                      bool
	hasAlternateId                     bool
	hasDesc                            bool
	hasMeaning                         bool
	isAnonymous                        bool
	isExternal                         bool
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
	if len(m.genderedMessagesMap) == 0 {
		return &jsMessageEmptyGenderMap
	}
	out := jutil.NewLinkedHashMapWithCapacity[JsMessage_GrammaticalGenderCase, []JsMessage_Part](len(m.genderedMessagesOrder))
	for _, k := range m.genderedMessagesOrder {
		out.PutNew(k, m.genderedMessagesMap[k])
	}
	return out
}

// HasGenderedVariants reports whether the message has gendered variants.
func (m *JsMessage) HasGenderedVariants() bool {
	return len(m.genderedMessagesMap) > 0
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
	if len(m.placeholderNameToExampleMap) == 0 {
		return &jsMessageEmptyStringMap
	}
	out := jutil.NewLinkedHashMapWithCapacity[string, string](len(m.placeholderNameToExampleOrder))
	for _, k := range m.placeholderNameToExampleOrder {
		out.PutNew(k, m.placeholderNameToExampleMap[k])
	}
	return out
}

// Java: JsMessage.getPlaceholderNameToOriginalCodeMap()
func (m *JsMessage) GetPlaceholderNameToOriginalCodeMap() *jutil.LinkedHashMap[string, string] {
	if len(m.placeholderNameToOriginalCodeMap) == 0 {
		return &jsMessageEmptyStringMap
	}
	out := jutil.NewLinkedHashMapWithCapacity[string, string](len(m.placeholderNameToOriginalCodeOrder))
	for _, k := range m.placeholderNameToOriginalCodeOrder {
		out.PutNew(k, m.placeholderNameToOriginalCodeMap[k])
	}
	return out
}

// Java: JsMessage.jsPlaceholderNames()
func (m *JsMessage) JsPlaceholderNames() *jutil.LinkedHashSet[string] {
	if len(m.jsPlaceholderNamesOrder) == 0 {
		return &jsMessageEmptyStringSet
	}
	out := jutil.NewLinkedHashSetWithCapacity[string](len(m.jsPlaceholderNamesOrder))
	for _, name := range m.jsPlaceholderNamesOrder {
		out.AddNew(name)
	}
	return out
}

// Java: JsMessage.canonicalPlaceholderNames()
func (m *JsMessage) CanonicalPlaceholderNames() *jutil.LinkedHashSet[string] {
	if len(m.canonicalPlaceholderNamesOrder) == 0 {
		return &jsMessageEmptyStringSet
	}
	out := jutil.NewLinkedHashSetWithCapacity[string](len(m.canonicalPlaceholderNamesOrder))
	for _, name := range m.canonicalPlaceholderNamesOrder {
		out.AddNew(name)
	}
	return out
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
	if v, ok := m.placeholderNameToOriginalCodeMap[placeholderReference.StoredPlaceholderName()]; ok {
		return v
	}
	return "-"
}

// Java: JsMessage.getPlaceholderExample(PlaceholderReference)
func (m *JsMessage) GetPlaceholderExample(placeholderReference JsMessage_Part) string {
	if v, ok := m.placeholderNameToExampleMap[placeholderReference.StoredPlaceholderName()]; ok {
		return v
	}
	return "-"
}

// GetGenderedMessageParts returns the parts for genderCase (nil if there are
// none for it).
//
// Java: JsMessage.getGenderedMessageParts(GrammaticalGenderCase)
func (m *JsMessage) GetGenderedMessageParts(genderCase JsMessage_GrammaticalGenderCase) []JsMessage_Part {
	if len(m.genderedMessagesMap) == 0 {
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"Message does not contain grammatical gendered variants."))
	}
	return m.genderedMessagesMap[genderCase]
}

// GetGenderedMessageVariants gets the list of grammatical gender cases for
// the message.
//
// Java: JsMessage.getGenderedMessageVariants()
func (m *JsMessage) GetGenderedMessageVariants() []JsMessage_GrammaticalGenderCase {
	if len(m.genderedMessagesOrder) == 0 {
		return []JsMessage_GrammaticalGenderCase{}
	}
	return slices.Clone(m.genderedMessagesOrder)
}

// AsJsMessageString returns a single string representing the message in
// `goog.getMsg()` placeholder format ("{$jsPlaceholderName}").
//
// Java: JsMessage.asJsMessageString()
func (m *JsMessage) AsJsMessageString() string {
	if len(m.genderedMessagesMap) > 0 {
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
	parts, ok := m.genderedMessagesMap[genderCase]
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
	if len(m.genderedMessagesMap) > 0 {
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
	parts, ok := m.genderedMessagesMap[genderCase]
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
	if len(m.genderedMessagesMap) == 0 {
		for _, part := range m.parts {
			if part.IsPlaceholder() || len(part.str) > 0 {
				return false
			}
		}
	} else {
		for _, k := range m.genderedMessagesOrder {
			for _, part := range m.genderedMessagesMap[k] {
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
		maps.EqualFunc(m.genderedMessagesMap, that.genderedMessagesMap, slices.Equal) &&
		m.hasAlternateId == that.hasAlternateId && m.alternateId == that.alternateId &&
		m.hasDesc == that.hasDesc && m.desc == that.desc &&
		m.hasMeaning == that.hasMeaning && m.meaning == that.meaning &&
		maps.Equal(m.placeholderNameToExampleMap, that.placeholderNameToExampleMap) &&
		maps.Equal(m.placeholderNameToOriginalCodeMap, that.placeholderNameToOriginalCodeMap) &&
		maps.Equal(m.jsPlaceholderNames, that.jsPlaceholderNames) &&
		maps.Equal(m.canonicalPlaceholderNames, that.canonicalPlaceholderNames)
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

func jsMessage_stringMapHashCode(m map[string]string, order []string) int32 {
	var h int32
	for _, k := range order {
		h += jutil.StringHashCode(k) ^ jutil.StringHashCode(m[k])
	}
	return h
}

func jsMessage_stringSliceHashCode(order []string) int32 {
	var h int32
	for _, k := range order {
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
	for _, k := range m.genderedMessagesOrder {
		gendered += int32(k) ^ jsMessage_partsHashCode(m.genderedMessagesMap[k])
	}
	h = 31*h + gendered
	h = 31*h + jsMessage_optStringHashCode(m.alternateId, m.hasAlternateId)
	h = 31*h + jsMessage_optStringHashCode(m.desc, m.hasDesc)
	h = 31*h + jsMessage_optStringHashCode(m.meaning, m.hasMeaning)
	h = 31*h + jsMessage_stringMapHashCode(m.placeholderNameToExampleMap, m.placeholderNameToExampleOrder)
	h = 31*h + jsMessage_stringMapHashCode(m.placeholderNameToOriginalCodeMap, m.placeholderNameToOriginalCodeOrder)
	h = 31*h + jsMessage_stringSliceHashCode(m.jsPlaceholderNamesOrder)
	h = 31*h + jsMessage_stringSliceHashCode(m.canonicalPlaceholderNamesOrder)
	return h
}

func jsMessage_formatPartsList(parts []JsMessage_Part) string {
	if len(parts) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, p := range parts {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(p.String())
	}
	sb.WriteByte(']')
	return sb.String()
}

func jsMessage_formatStringMap(m map[string]string, order []string) string {
	if len(order) == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, k := range order {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(m[k])
	}
	sb.WriteByte('}')
	return sb.String()
}

func jsMessage_formatStringSlice(order []string) string {
	if len(order) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, elem := range order {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(elem)
	}
	sb.WriteByte(']')
	return sb.String()
}

// String is the record's toString.
//
// Java: JsMessage.toString()
func (m *JsMessage) String() string {
	var gendered string
	if len(m.genderedMessagesOrder) == 0 {
		gendered = "{}"
	} else {
		var sb strings.Builder
		sb.WriteByte('{')
		for i, k := range m.genderedMessagesOrder {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(k.String())
			sb.WriteByte('=')
			sb.WriteString(jsMessage_formatPartsList(m.genderedMessagesMap[k]))
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
		", getParts=" + jsMessage_formatPartsList(m.parts) +
		", getGenderedMessagesMap=" + gendered +
		", getAlternateId=" + jsMessage_optStringFormat(m.alternateId, m.hasAlternateId) +
		", getDesc=" + jsMessage_optStringFormat(m.desc, m.hasDesc) +
		", getMeaning=" + jsMessage_optStringFormat(m.meaning, m.hasMeaning) +
		", getPlaceholderNameToExampleMap=" + jsMessage_formatStringMap(m.placeholderNameToExampleMap, m.placeholderNameToExampleOrder) +
		", getPlaceholderNameToOriginalCodeMap=" + jsMessage_formatStringMap(m.placeholderNameToOriginalCodeMap, m.placeholderNameToOriginalCodeOrder) +
		", jsPlaceholderNames=" + jsMessage_formatStringSlice(m.jsPlaceholderNamesOrder) +
		", canonicalPlaceholderNames=" + jsMessage_formatStringSlice(m.canonicalPlaceholderNamesOrder) +
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
		panic(jutil.NewUnsupportedOperationExceptionMsg("not a string part: '" + p.GetJsPlaceholderName() + "'"))
	}
	return p.str
}

// StoredPlaceholderName returns the raw stored placeholder name.
//
// Java: JsMessage.PlaceholderReference.storedPlaceholderName()
func (p JsMessage_Part) StoredPlaceholderName() string {
	if !p.isPlaceholder {
		panic(jutil.NewUnsupportedOperationExceptionMsg("not a placeholder: '" + p.str + "'"))
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
		panic(jutil.NewUnsupportedOperationExceptionMsg("not a placeholder: '" + p.str + "'"))
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
		panic(jutil.NewUnsupportedOperationExceptionMsg("not a placeholder: '" + p.str + "'"))
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

// JsMessage_IdGenerator is the ID generator.
type JsMessage_IdGenerator interface {
	// GenerateId generates the ID for the message. Messages with the same
	// messageParts and meaning will get the same id. Messages with the same
	// id will get the same translation.
	//
	// Java: JsMessage.IdGenerator.generateId(String, List<Part>)
	GenerateId(meaning *string, messageParts []JsMessage_Part) string
}

func generateMessageId(gen JsMessage_IdGenerator, meaning string, parts []JsMessage_Part) string {
	if g, ok := gen.(*GoogleJsMessageIdGenerator); ok {
		return g.generateId(meaning, parts)
	}
	return gen.GenerateId(&meaning, parts)
}

// MessageBundle is an interface for providing alternative values for
// user-visible messages in JavaScript code.
type MessageBundle interface {
	// IdGenerator gets the message ID generator to use to compute message
	// IDs for this type of bundle.
	//
	// Returns the idGenerator instance or null if we do not want to use any
	// custom id generation. In case if idGenerator is null caller should
	// decide how to create id by itself. In the most cases using the message
	// key is enough.
	//
	// Java: MessageBundle.idGenerator()
	IdGenerator() JsMessage_IdGenerator

	// GetMessage gets a message replacement.
	//
	// id is the id of the message being replaced; the key is message ID
	// generated by JsMessage.IdGenerator. Returns the message replacement,
	// which may be null.
	//
	// Java: MessageBundle.getMessage(String)
	GetMessage(id string) *JsMessage

	// GetAllMessages returns an iterable over the keys that this object has
	// replacements for.
	//
	// Java: MessageBundle.getAllMessages()
	GetAllMessages() iter.Seq[*JsMessage]
}

// EmptyMessageBundle is an implementation of MessageBundle that has no
// translations.
type EmptyMessageBundle struct{}

// NewEmptyMessageBundle is new EmptyMessageBundle().
func NewEmptyMessageBundle() *EmptyMessageBundle { return &EmptyMessageBundle{} }

// IdGenerator gets a dummy message ID generator.
//
// Java: EmptyMessageBundle.idGenerator()
func (b *EmptyMessageBundle) IdGenerator() JsMessage_IdGenerator {
	return nil
}

// GetMessage returns null, to indicate it has no message replacements.
//
// Java: EmptyMessageBundle.getMessage(String)
func (b *EmptyMessageBundle) GetMessage(id string) *JsMessage {
	return nil
}

// GetAllMessages returns an empty list of messages.
//
// Java: EmptyMessageBundle.getAllMessages()
func (b *EmptyMessageBundle) GetAllMessages() iter.Seq[*JsMessage] {
	return func(yield func(*JsMessage) bool) {}
}

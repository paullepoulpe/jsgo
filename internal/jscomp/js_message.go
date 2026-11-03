// Port of com/google/javascript/jscomp/JsMessage.java @ closure-compiler v20261028
// Port of com/google/javascript/jscomp/MessageBundle.java @ closure-compiler v20261028

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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsMessageDefinition is the container struct that holds a JsMessage and
// information about the `goog.getMsg()` call it was built from.
type JsMessageDefinition struct {
	message               *JsMessage
	placeholderValues     OrderedStringMap[rhino.Node]
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
func (d *JsMessageDefinition) GetPlaceholderValueMap() *OrderedStringMap[rhino.Node] {
	return &d.placeholderValues
}

// LookupPlaceholderNode looks up the value Node assigned to placeholder name.
func (d *JsMessageDefinition) LookupPlaceholderNode(name string) (rhino.Node, bool) {
	return d.placeholderValues.Lookup(name)
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
	sourceName                string
	key                       string
	id                        string
	alternateId               string
	desc                      string
	meaning                   string
	parts                     []JsMessage_Part
	genderedMessages          OrderedMap[JsMessage_GrammaticalGenderCase, []JsMessage_Part]
	placeholderExamples       OrderedStringMap[string]
	placeholderOriginalCode   OrderedStringMap[string]
	jsPlaceholderNames        OrderedSet[string]
	canonicalPlaceholderNames OrderedSet[string]
	hasSourceName             bool
	hasAlternateId            bool
	hasDesc                   bool
	hasMeaning                bool
	isAnonymous               bool
	isExternal                bool
}

// SourceNameOpt returns the source name and whether one is present.
//
// Java: JsMessage.getSourceName()
func (m *JsMessage) SourceNameOpt() (string, bool) { return m.sourceName, m.hasSourceName }

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

// HasGenderedVariants reports whether the message has gendered variants.
func (m *JsMessage) HasGenderedVariants() bool {
	return !m.genderedMessages.IsEmpty()
}

// GenderedMessagesAll yields gendered message variants in insertion order.
func (m *JsMessage) GenderedMessagesAll() iter.Seq2[JsMessage_GrammaticalGenderCase, []JsMessage_Part] {
	return m.genderedMessages.All()
}

// AlternateIdOpt returns the alternate ID and whether one is present.
//
// Java: JsMessage.getAlternateId()
func (m *JsMessage) AlternateIdOpt() (string, bool) { return m.alternateId, m.hasAlternateId }

// DescOpt returns the description and whether one is present.
//
// Java: JsMessage.getDesc()
func (m *JsMessage) DescOpt() (string, bool) { return m.desc, m.hasDesc }

// MeaningOpt returns the meaning and whether one is present.
//
// Java: JsMessage.getMeaning()
func (m *JsMessage) MeaningOpt() (string, bool) { return m.meaning, m.hasMeaning }

// Java: JsMessage.getPlaceholderNameToExampleMap()
func (m *JsMessage) GetPlaceholderNameToExampleMap() *OrderedStringMap[string] {
	return &m.placeholderExamples
}

// Java: JsMessage.getPlaceholderNameToOriginalCodeMap()
func (m *JsMessage) GetPlaceholderNameToOriginalCodeMap() *OrderedStringMap[string] {
	return &m.placeholderOriginalCode
}

// Java: JsMessage.jsPlaceholderNames()
func (m *JsMessage) JsPlaceholderNames() *OrderedSet[string] {
	return &m.jsPlaceholderNames
}

// HasJsPlaceholders reports whether the message references any JS placeholders.
func (m *JsMessage) HasJsPlaceholders() bool {
	return !m.jsPlaceholderNames.IsEmpty()
}

// JsPlaceholderSlice returns the ordered JS placeholder names.
func (m *JsMessage) JsPlaceholderSlice() []string {
	return m.jsPlaceholderNames.ToSlice()
}

func (m *JsMessage) sameJsPlaceholders(other *JsMessage) bool {
	return m.jsPlaceholderNames.Equal(&other.jsPlaceholderNames)
}

func (m *JsMessage) formatJsPlaceholders() string {
	return jsMessage_formatStringSet(&m.jsPlaceholderNames)
}

// Java: JsMessage.canonicalPlaceholderNames()
func (m *JsMessage) CanonicalPlaceholderNames() *OrderedSet[string] {
	return &m.canonicalPlaceholderNames
}

// HasCanonicalPlaceholders reports whether the message references any canonical placeholders.
func (m *JsMessage) HasCanonicalPlaceholders() bool {
	return !m.canonicalPlaceholderNames.IsEmpty()
}

// CanonicalPlaceholderAll yields canonical placeholder names in insertion order.
func (m *JsMessage) CanonicalPlaceholderAll() iter.Seq[string] {
	return m.canonicalPlaceholderNames.All()
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

// String is Enum.toString().
func (g JsMessage_GrammaticalGenderCase) String() string {
	return jsMessage_GrammaticalGenderCase_names[g]
}

// Java: JsMessage.getPlaceholderOriginalCode(PlaceholderReference)
func (m *JsMessage) GetPlaceholderOriginalCode(placeholderReference JsMessage_Part) string {
	if v, ok := m.placeholderOriginalCode.Lookup(placeholderReference.StoredPlaceholderName()); ok {
		return v
	}
	return "-"
}

// Java: JsMessage.getPlaceholderExample(PlaceholderReference)
func (m *JsMessage) GetPlaceholderExample(placeholderReference JsMessage_Part) string {
	if v, ok := m.placeholderExamples.Lookup(placeholderReference.StoredPlaceholderName()); ok {
		return v
	}
	return "-"
}

// GetGenderedMessageParts returns the parts for genderCase (nil if there are
// none for it).
//
// Java: JsMessage.getGenderedMessageParts(GrammaticalGenderCase)
func (m *JsMessage) GetGenderedMessageParts(genderCase JsMessage_GrammaticalGenderCase) []JsMessage_Part {
	if m.genderedMessages.IsEmpty() {
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"Message does not contain grammatical gendered variants."))
	}
	return m.genderedMessages.Get(genderCase)
}

// GetGenderedMessageVariants gets the list of grammatical gender cases for
// the message.
//
// Java: JsMessage.getGenderedMessageVariants()
func (m *JsMessage) GetGenderedMessageVariants() []JsMessage_GrammaticalGenderCase {
	if m.genderedMessages.IsEmpty() {
		return []JsMessage_GrammaticalGenderCase{}
	}
	return m.genderedMessages.KeySlice()
}

// AsJsMessageString returns a single string representing the message in
// `goog.getMsg()` placeholder format ("{$jsPlaceholderName}").
//
// Java: JsMessage.asJsMessageString()
func (m *JsMessage) AsJsMessageString() string {
	if !m.genderedMessages.IsEmpty() {
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
	parts, ok := m.genderedMessages.Lookup(genderCase)
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
	if !m.genderedMessages.IsEmpty() {
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
	parts, ok := m.genderedMessages.Lookup(genderCase)
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
	if m.genderedMessages.IsEmpty() {
		for _, part := range m.parts {
			if part.IsPlaceholder() || len(part.str) > 0 {
				return false
			}
		}
	} else {
		for _, e := range m.genderedMessages.Entries() {
			for _, part := range e.Value {
				if part.IsPlaceholder() || len(part.str) > 0 {
					return false
				}
			}
		}
	}
	return true
}

func jsMessage_optStringFormat(s string, ok bool) string {
	if !ok {
		return "null"
	}
	return s
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

func jsMessage_formatStringMap(m *OrderedStringMap[string]) string {
	entries := m.Entries()
	if len(entries) == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, e := range entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(e.Key)
		sb.WriteByte('=')
		sb.WriteString(e.Value)
	}
	sb.WriteByte('}')
	return sb.String()
}

func jsMessage_formatStringSet(s *OrderedSet[string]) string {
	entries := s.Entries()
	if len(entries) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, e := range entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(e.Key)
	}
	sb.WriteByte(']')
	return sb.String()
}

// String is the record's toString.
//
// Java: JsMessage.toString()
func (m *JsMessage) String() string {
	var gendered string
	entries := m.genderedMessages.Entries()
	if len(entries) == 0 {
		gendered = "{}"
	} else {
		var sb strings.Builder
		sb.WriteByte('{')
		for i, e := range entries {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(e.Key.String())
			sb.WriteByte('=')
			sb.WriteString(jsMessage_formatPartsList(e.Value))
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
		", getPlaceholderNameToExampleMap=" + jsMessage_formatStringMap(&m.placeholderExamples) +
		", getPlaceholderNameToOriginalCodeMap=" + jsMessage_formatStringMap(&m.placeholderOriginalCode) +
		", jsPlaceholderNames=" + jsMessage_formatStringSet(&m.jsPlaceholderNames) +
		", canonicalPlaceholderNames=" + jsMessage_formatStringSet(&m.canonicalPlaceholderNames) +
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
	return caseFormatLowerCamelToUpperUnderscore(p.str)
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

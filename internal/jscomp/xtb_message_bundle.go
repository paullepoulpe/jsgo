// Port of com/google/javascript/jscomp/XtbMessageBundle.java @ closure-compiler v20261021

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
	"io"
	"iter"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// XtbMessageBundle is a MessageBundle that parses messages from an XML
// Translation Bundle (XTB) file.
type XtbMessageBundle struct {
	messages    map[string]*JsMessage
	order       []string
	idGenerator JsMessage_IdGenerator
}

func (b *XtbMessageBundle) putMessage(key string, msg *JsMessage) {
	if b.messages == nil {
		b.messages = make(map[string]*JsMessage)
	}
	if _, exists := b.messages[key]; !exists {
		b.order = append(b.order, key)
	}
	b.messages[key] = msg
}

var xtbMessageBundle_NOOP_RESOLVER = &xtbMessageBundle_SecureEntityResolver{}

// NewXtbMessageBundle creates an instance and initializes it with the
// messages in an XTB file.
//
// xtb is the XTB file as a byte stream; projectId is the translation console
// project id (i.e. name)
//
// Java: XtbMessageBundle(InputStream, String)
func NewXtbMessageBundle(xtb io.Reader, projectId *string) *XtbMessageBundle {
	b := &XtbMessageBundle{
		idGenerator: NewGoogleJsMessageIdGenerator(projectId),
	}

	jutil.Try(func() {
		// Use a SAX parser for speed and less memory usage.
		parser := xtbMessageBundle_createSAXParser()
		reader := parser.GetXMLReader()
		contentHandler := &xtbMessageBundle_Handler{outer: b}
		reader.SetContentHandler(contentHandler)
		reader.Parse(jutil.NewInputSource(xtb))
	}, func(e jutil.IException) {
		switch e.(type) {
		case *jutil.ParserConfigurationException, jutil.IIOException, jutil.ISAXException:
			panic(jutil.NewRuntimeExceptionWithCause(e))
		}
		panic(e)
	})
	return b
}

// xtbMessageBundle_createSAXParser is inlined from guava-internal.
//
// Java: XtbMessageBundle.createSAXParser()
func xtbMessageBundle_createSAXParser() *jutil.SAXParser {
	factory := jutil.SAXParserFactory_NewInstance()
	factory.SetValidating(false)
	factory.SetXIncludeAware(false)
	factory.SetFeature(
		"http://xml.org/sax/features/external-general-entities", false)
	factory.SetFeature(
		"http://xml.org/sax/features/external-parameter-entities", false)
	factory.SetFeature(
		"http://apache.org/xml/features/nonvalidating/load-external-dtd",
		false)
	factory.SetFeature(jutil.XMLConstants_FEATURE_SECURE_PROCESSING, true)

	parser := factory.NewSAXParser()
	xmlReader := parser.GetXMLReader()
	xmlReader.SetEntityResolver(xtbMessageBundle_NOOP_RESOLVER)
	return parser
}

// Java: XtbMessageBundle.getMessage(String)
func (b *XtbMessageBundle) GetMessage(id string) *JsMessage {
	return b.messages[id]
}

// Java: XtbMessageBundle.idGenerator()
func (b *XtbMessageBundle) IdGenerator() JsMessage_IdGenerator {
	return b.idGenerator
}

// Java: XtbMessageBundle.getAllMessages()
func (b *XtbMessageBundle) GetAllMessages() iter.Seq[*JsMessage] {
	return func(yield func(*JsMessage) bool) {
		for _, k := range b.order {
			if !yield(b.messages[k]) {
				return
			}
		}
	}
}

const (
	xtbMessageBundle_Handler_BUNDLE_ELEM_NAME          = "translationbundle"
	xtbMessageBundle_Handler_LANG_ATT_NAME             = "lang"
	xtbMessageBundle_Handler_TRANSLATION_ELEM_NAME     = "translation"
	xtbMessageBundle_Handler_MESSAGE_ID_ATT_NAME       = "id"
	xtbMessageBundle_Handler_PLACEHOLDER_ELEM_NAME     = "ph"
	xtbMessageBundle_Handler_PLACEHOLDER_NAME_ATT_NAME = "name"
	xtbMessageBundle_Handler_BRANCH_ELEM_NAME          = "branch"
	xtbMessageBundle_Handler_BRANCH_NAME_ATT_NAME      = "variants"
	xtbMessageBundle_Handler_MALE_GENDER_CASE          = "MASCULINE"
	xtbMessageBundle_Handler_FEMALE_GENDER_CASE        = "FEMININE"
	xtbMessageBundle_Handler_NEUTER_GENDER_CASE        = "NEUTER"
	xtbMessageBundle_Handler_OTHER_GENDER_CASE         = "OTHER"
	xtbMessageBundle_Handler_GENDER_VARIANT_PREFIX     = "grammatical_gender_case:"
)

// parseXtbGenderVariant matches `grammatical_gender_case:\s*(\w+)` on ASCII
// whitespace and word characters without regex compilation or allocation.
func parseXtbGenderVariant(s string) (string, bool) {
	for pos := 0; pos < len(s); {
		idx := strings.Index(s[pos:], xtbMessageBundle_Handler_GENDER_VARIANT_PREFIX)
		if idx < 0 {
			return "", false
		}
		i := pos + idx + len(xtbMessageBundle_Handler_GENDER_VARIANT_PREFIX)
		for i < len(s) {
			switch s[i] {
			case ' ', '\t', '\n', '\v', '\f', '\r':
				i++
				continue
			}
			break
		}
		start := i
		for i < len(s) {
			c := s[i]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
				i++
			} else {
				break
			}
		}
		if i > start {
			return s[start:i], true
		}
		pos = start
	}
	return "", false
}

// xtbMessageBundle_Handler is a ContentHandler that creates a JsMessage for
// each message parsed from an XML Translation Bundle (XTB) file.
//
// Java: private class XtbMessageBundle.Handler implements ContentHandler
type xtbMessageBundle_Handler struct {
	outer         *XtbMessageBundle
	msgBuilder    *JsMessage_Builder
	genderCase    JsMessage_GrammaticalGenderCase
	hasGenderCase bool
	hasLang       bool
}

func (h *xtbMessageBundle_Handler) requireGenderCase() JsMessage_GrammaticalGenderCase {
	if !h.hasGenderCase {
		panic(jutil.NewNullPointerException())
	}
	return h.genderCase
}

// Java: XtbMessageBundle.Handler.setDocumentLocator(Locator)
func (h *xtbMessageBundle_Handler) SetDocumentLocator(locator jutil.Locator) {}

// Java: XtbMessageBundle.Handler.startDocument()
func (h *xtbMessageBundle_Handler) StartDocument() {}

// Java: XtbMessageBundle.Handler.endDocument()
func (h *xtbMessageBundle_Handler) EndDocument() {}

// Java: XtbMessageBundle.Handler.startPrefixMapping(String, String)
func (h *xtbMessageBundle_Handler) StartPrefixMapping(prefix string, uri string) {}

// Java: XtbMessageBundle.Handler.endPrefixMapping(String)
func (h *xtbMessageBundle_Handler) EndPrefixMapping(prefix string) {}

// Java: XtbMessageBundle.Handler.startElement(String, String, String, Attributes)
func (h *xtbMessageBundle_Handler) StartElement(uri string, localName string, qName string,
	atts jutil.Attributes) {
	switch qName {
	case xtbMessageBundle_Handler_BUNDLE_ELEM_NAME:
		jutil.CheckState(!h.hasLang)
		lang := atts.GetValue(xtbMessageBundle_Handler_LANG_ATT_NAME)
		jutil.CheckState(lang != nil && *lang != "")
		h.hasLang = true
	case xtbMessageBundle_Handler_TRANSLATION_ELEM_NAME:
		jutil.CheckState(h.msgBuilder == nil)
		id := atts.GetValue(xtbMessageBundle_Handler_MESSAGE_ID_ATT_NAME)
		jutil.CheckState(id != nil && *id != "")
		h.msgBuilder = NewJsMessage_Builder().SetKey(*id).SetId(*id)
	case xtbMessageBundle_Handler_PLACEHOLDER_ELEM_NAME:
		jutil.CheckState(h.msgBuilder != nil)
		phRef := *jutil.CheckNotNullMsg(
			atts.GetValue(xtbMessageBundle_Handler_PLACEHOLDER_NAME_ATT_NAME),
			"Placeholder name must not be null")
		if h.msgBuilder.HasGenderedVariants() {
			h.msgBuilder.AppendCanonicalPlaceholderReferenceForGender(h.requireGenderCase(), phRef)
		} else {
			h.msgBuilder.AppendCanonicalPlaceholderReference(phRef)
		}
	case xtbMessageBundle_Handler_BRANCH_ELEM_NAME:
		jutil.CheckState(h.msgBuilder != nil)
		gender := atts.GetValue(xtbMessageBundle_Handler_BRANCH_NAME_ATT_NAME)
		if variant, ok := parseXtbGenderVariant(*jutil.CheckNotNull(gender)); ok {
			var parsed JsMessage_GrammaticalGenderCase
			switch variant {
			case xtbMessageBundle_Handler_MALE_GENDER_CASE:
				parsed = JsMessage_GrammaticalGenderCase_MASCULINE
			case xtbMessageBundle_Handler_FEMALE_GENDER_CASE:
				parsed = JsMessage_GrammaticalGenderCase_FEMININE
			case xtbMessageBundle_Handler_NEUTER_GENDER_CASE:
				parsed = JsMessage_GrammaticalGenderCase_NEUTER
			case xtbMessageBundle_Handler_OTHER_GENDER_CASE:
				parsed = JsMessage_GrammaticalGenderCase_OTHER
			default:
				panic(jutil.NewIllegalStateExceptionMsg(
					"Gender case must be one of the following: MASCULINE, FEMININE, NEUTER, or OTHER."))
			}
			h.genderCase = parsed
			h.hasGenderCase = true
			h.msgBuilder.AddGenderedMessageKey(parsed)
		}
	}
}

// Java: XtbMessageBundle.Handler.endElement(String, String, String)
func (h *xtbMessageBundle_Handler) EndElement(uri string, localName string, qName string) {
	if xtbMessageBundle_Handler_TRANSLATION_ELEM_NAME == qName {
		jutil.CheckState(h.msgBuilder != nil)
		if !h.msgBuilder.HasParts() {
			h.msgBuilder.AppendStringPart("")
		}
		h.outer.putMessage(h.msgBuilder.key, h.msgBuilder.Build())
		h.msgBuilder = nil
	}
}

// Java: XtbMessageBundle.Handler.characters(char[], int, int)
func (h *xtbMessageBundle_Handler) Characters(ch []uint16, start int, length int) {
	if h.msgBuilder != nil {
		part := jutil.StringFromUTF16(ch[start : start+length])
		if h.msgBuilder.HasGenderedVariants() {
			h.msgBuilder.AppendStringPartForGender(h.requireGenderCase(), part)
		} else {
			h.msgBuilder.AppendStringPart(part)
		}
	}
}

// Java: XtbMessageBundle.Handler.ignorableWhitespace(char[], int, int)
func (h *xtbMessageBundle_Handler) IgnorableWhitespace(ch []uint16, start int, length int) {
	if h.msgBuilder != nil {
		part := jutil.StringFromUTF16(ch[start : start+length])
		if h.msgBuilder.HasGenderedVariants() {
			h.msgBuilder.AppendStringPartForGender(h.requireGenderCase(), part)
		} else {
			h.msgBuilder.AppendStringPart(part)
		}
	}
}

// Java: XtbMessageBundle.Handler.processingInstruction(String, String)
func (h *xtbMessageBundle_Handler) ProcessingInstruction(target string, data string) {}

// Java: XtbMessageBundle.Handler.skippedEntity(String)
func (h *xtbMessageBundle_Handler) SkippedEntity(name string) {}

// xtbMessageBundle_SecureEntityResolver returns an empty string in response to
// any attempt to resolve an external entity.
type xtbMessageBundle_SecureEntityResolver struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: XtbMessageBundle.SecureEntityResolver.resolveEntity(String, String)
func (r *xtbMessageBundle_SecureEntityResolver) ResolveEntity(publicId *string, systemId *string) *jutil.InputSource {
	return jutil.NewInputSourceOfCharacterStream(strings.NewReader(""))
}

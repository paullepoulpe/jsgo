// Port of com/google/javascript/jscomp/XtbMessageBundle.java @ closure-compiler v20260415

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
	messages    *jutil.LinkedHashMap[string, *JsMessage]
	idGenerator JsMessage_IdGenerator
}

var xtbMessageBundle_NOOP_RESOLVER = newXtbMessageBundle_SecureEntityResolver()

// NewXtbMessageBundle creates an instance and initializes it with the
// messages in an XTB file.
//
// xtb is the XTB file as a byte stream; projectId is the translation console
// project id (i.e. name)
//
// Java: XtbMessageBundle(InputStream, String)
func NewXtbMessageBundle(xtb io.Reader, projectId *string) *XtbMessageBundle {
	b := &XtbMessageBundle{}
	b.messages = jutil.NewLinkedHashMap[string, *JsMessage]()
	b.idGenerator = NewGoogleJsMessageIdGenerator(projectId)

	jutil.Try(func() {
		// Use a SAX parser for speed and less memory usage.
		parser := xtbMessageBundle_createSAXParser()
		reader := parser.GetXMLReader()
		contentHandler := newXtbMessageBundle_Handler(b)
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
	return b.messages.Get(id)
}

// Java: XtbMessageBundle.idGenerator()
func (b *XtbMessageBundle) IdGenerator() JsMessage_IdGenerator {
	return b.idGenerator
}

// Java: XtbMessageBundle.getAllMessages()
func (b *XtbMessageBundle) GetAllMessages() iter.Seq[*JsMessage] {
	return b.messages.Values()
}

const xtbMessageBundle_Handler_BUNDLE_ELEM_NAME = "translationbundle"
const xtbMessageBundle_Handler_LANG_ATT_NAME = "lang"

const xtbMessageBundle_Handler_TRANSLATION_ELEM_NAME = "translation"
const xtbMessageBundle_Handler_MESSAGE_ID_ATT_NAME = "id"

const xtbMessageBundle_Handler_PLACEHOLDER_ELEM_NAME = "ph"
const xtbMessageBundle_Handler_PLACEHOLDER_NAME_ATT_NAME = "name"

const xtbMessageBundle_Handler_BRANCH_ELEM_NAME = "branch"
const xtbMessageBundle_Handler_BRANCH_NAME_ATT_NAME = "variants"
const xtbMessageBundle_Handler_MALE_GENDER_CASE = "MASCULINE"
const xtbMessageBundle_Handler_FEMALE_GENDER_CASE = "FEMININE"
const xtbMessageBundle_Handler_NEUTER_GENDER_CASE = "NEUTER"
const xtbMessageBundle_Handler_OTHER_GENDER_CASE = "OTHER"

const xtbMessageBundle_Handler_GENDER_VARIANT_PATTERN = `grammatical_gender_case:\s*(\w+)`

var xtbMessageBundle_Handler_GET_STRING_PATTERN = jutil.Pattern_Compile(xtbMessageBundle_Handler_GENDER_VARIANT_PATTERN)

// xtbMessageBundle_Handler is a ContentHandler that creates a JsMessage for
// each message parsed from an XML Translation Bundle (XTB) file.
//
// Java: private class XtbMessageBundle.Handler implements ContentHandler
type xtbMessageBundle_Handler struct {
	outer *XtbMessageBundle

	genderCase *string

	lang       *string
	msgBuilder *JsMessage_Builder
}

func newXtbMessageBundle_Handler(outer *XtbMessageBundle) *xtbMessageBundle_Handler {
	return &xtbMessageBundle_Handler{outer: outer}
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
		jutil.CheckState(h.lang == nil)
		h.lang = atts.GetValue(xtbMessageBundle_Handler_LANG_ATT_NAME)
		jutil.CheckState(h.lang != nil && *h.lang != "")
	case xtbMessageBundle_Handler_TRANSLATION_ELEM_NAME:
		jutil.CheckState(h.msgBuilder == nil)
		id := atts.GetValue(xtbMessageBundle_Handler_MESSAGE_ID_ATT_NAME)
		jutil.CheckState(id != nil && *id != "")
		h.msgBuilder = NewJsMessage_Builder().SetKey(*id).SetId(*id)
	case xtbMessageBundle_Handler_PLACEHOLDER_ELEM_NAME:
		jutil.CheckState(h.msgBuilder != nil)
		phRef := atts.GetValue(xtbMessageBundle_Handler_PLACEHOLDER_NAME_ATT_NAME)
		if h.msgBuilder.HasGenderedVariants() {
			h.msgBuilder.AppendCanonicalPlaceholderReferenceForGender(
				JsMessage_GrammaticalGenderCase_ValueOf(*h.genderCase), *jutil.CheckNotNullMsg(phRef, "Placeholder name must not be null"))
		} else {
			h.msgBuilder.AppendCanonicalPlaceholderReference(*jutil.CheckNotNullMsg(phRef, "Placeholder name must not be null"))
		}
	case xtbMessageBundle_Handler_BRANCH_ELEM_NAME:
		jutil.CheckState(h.msgBuilder != nil)
		gender := atts.GetValue(xtbMessageBundle_Handler_BRANCH_NAME_ATT_NAME)
		// Gender case must be one of the following: MALE, FEMALE, NEUTER, or OTHER.

		matcher := xtbMessageBundle_Handler_GET_STRING_PATTERN.Matcher(*jutil.CheckNotNull(gender))
		if matcher.Find() {
			h.genderCase = matcher.Group(1)
			jutil.CheckStateMsg(
				*h.genderCase == xtbMessageBundle_Handler_MALE_GENDER_CASE ||
					*h.genderCase == xtbMessageBundle_Handler_FEMALE_GENDER_CASE ||
					*h.genderCase == xtbMessageBundle_Handler_NEUTER_GENDER_CASE ||
					*h.genderCase == xtbMessageBundle_Handler_OTHER_GENDER_CASE,
				"Gender case must be one of the following: MASCULINE, FEMININE, NEUTER, or OTHER.")
			h.msgBuilder.AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_ValueOf(*h.genderCase))
		}
	default:
	}
}

// Java: XtbMessageBundle.Handler.endElement(String, String, String)
func (h *xtbMessageBundle_Handler) EndElement(uri string, localName string, qName string) {
	if xtbMessageBundle_Handler_TRANSLATION_ELEM_NAME == qName {
		jutil.CheckState(h.msgBuilder != nil)
		if !h.msgBuilder.HasParts() {
			h.msgBuilder.AppendStringPart("")
		}
		key := h.msgBuilder.GetKey()
		h.outer.messages.Put(*key, h.msgBuilder.Build())
		h.msgBuilder = nil
	}
}

// Java: XtbMessageBundle.Handler.characters(char[], int, int)
func (h *xtbMessageBundle_Handler) Characters(ch []uint16, start int, length int) {
	if h.msgBuilder != nil {
		part := jutil.StringFromUTF16(ch[start : start+length])
		// Append a string literal to the message.
		if h.msgBuilder.HasGenderedVariants() {
			h.msgBuilder.AppendStringPartForGender(JsMessage_GrammaticalGenderCase_ValueOf(*h.genderCase), part)
		} else {
			h.msgBuilder.AppendStringPart(part)
		}
	}
}

// Java: XtbMessageBundle.Handler.ignorableWhitespace(char[], int, int)
func (h *xtbMessageBundle_Handler) IgnorableWhitespace(ch []uint16, start int, length int) {
	if h.msgBuilder != nil {
		// Preserve whitespace in messages.
		if h.msgBuilder.HasGenderedVariants() {
			h.msgBuilder.AppendStringPartForGender(
				JsMessage_GrammaticalGenderCase_ValueOf(*h.genderCase),
				jutil.StringFromUTF16(ch[start:start+length]))
		} else {
			h.msgBuilder.AppendStringPart(jutil.StringFromUTF16(ch[start : start+length]))
		}
	}
}

// Java: XtbMessageBundle.Handler.processingInstruction(String, String)
func (h *xtbMessageBundle_Handler) ProcessingInstruction(target string, data string) {}

// Java: XtbMessageBundle.Handler.skippedEntity(String)
func (h *xtbMessageBundle_Handler) SkippedEntity(name string) {}

// xtbMessageBundle_SecureEntityResolver is a secure EntityResolver that
// returns an empty string in response to any attempt to resolve an external
// entity. The class is used by our secure version of the internal saxon SAX
// parser.
//
// Java: private static final class XtbMessageBundle.SecureEntityResolver implements EntityResolver
type xtbMessageBundle_SecureEntityResolver struct {
	_ byte // Go-only: language-adapter — distinct identity
}

func newXtbMessageBundle_SecureEntityResolver() *xtbMessageBundle_SecureEntityResolver {
	return &xtbMessageBundle_SecureEntityResolver{}
}

// Java: XtbMessageBundle.SecureEntityResolver.resolveEntity(String, String)
func (r *xtbMessageBundle_SecureEntityResolver) ResolveEntity(publicId *string, systemId *string) *jutil.InputSource {
	return jutil.NewInputSourceOfCharacterStream(strings.NewReader(""))
}

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
// Port of com/google/javascript/jscomp/XtbMessageBundle.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.XtbMessageBundle
type XtbMessageBundle struct {
	messages    j2gort.OrderedMap[string, *JsMessage] `gson:"messages"`
	idGenerator JsMessage_IdGenerator                 `gson:"idGenerator"`
	idHash_     atomic.Int32
}

func (this_ *XtbMessageBundle) BindSelf_(s any) {
	_ = s
}
func (this_ *XtbMessageBundle) IsJavaNil_() bool { return this_ == nil }
func (*XtbMessageBundle) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.XtbMessageBundle", "XtbMessageBundle"
}
func (this_ *XtbMessageBundle) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*XtbMessageBundle) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.XtbMessageBundle", "com.google.javascript.jscomp.XtbMessageBundle", func(v_ any) bool {
		return j2gort.InstanceOf[*XtbMessageBundle](v_)
	})
}
func (*XtbMessageBundle) Implements_jscomp_MessageBundle() {}

var static_XtbMessageBundle_NOOP_RESOLVER *XtbMessageBundle_SecureEntityResolver

func StaticGet_XtbMessageBundle_NOOP_RESOLVER() *XtbMessageBundle_SecureEntityResolver {
	clinit_XtbMessageBundle()
	return static_XtbMessageBundle_NOOP_RESOLVER
}
func StaticSet_XtbMessageBundle_NOOP_RESOLVER(v *XtbMessageBundle_SecureEntityResolver) {
	clinit_XtbMessageBundle()
	static_XtbMessageBundle_NOOP_RESOLVER = v
}

var clinitDone_XtbMessageBundle uint32
var clinitState_XtbMessageBundle j2gort.ClinitState

func clinit_XtbMessageBundle() {
	if atomic.LoadUint32(&clinitDone_XtbMessageBundle) != 0 {
		return
	}
	clinitSlow_XtbMessageBundle()
}

func clinitSlow_XtbMessageBundle() {
	if !j2gort.EnterClinit(&clinitState_XtbMessageBundle, &clinitDone_XtbMessageBundle) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_XtbMessageBundle, &clinitDone_XtbMessageBundle)
	static_XtbMessageBundle_NOOP_RESOLVER = NewXtbMessageBundle_SecureEntityResolver()
}

func (this_ *XtbMessageBundle) initFields_XtbMessageBundle() {
	_ = this_
}

// Java constructor: XtbMessageBundle(java.io.InputStream xtb, java.lang.@org.jspecify.annotations.Nullable String projectId)
func NewXtbMessageBundle(xtb j2gort.Ext_java_io_InputStream, projectId string) *XtbMessageBundle {
	clinit_XtbMessageBundle()
	this_ := &XtbMessageBundle{}
	InitXtbMessageBundle(this_, xtb, projectId)
	return this_
}

func InitXtbMessageBundle(this_ *XtbMessageBundle, xtb j2gort.Ext_java_io_InputStream, projectId string) {
	_ = xtb
	_ = projectId
	this_.initFields_XtbMessageBundle()
	this_.messages = j2gort.NewOrderedMap[string, *JsMessage]()
	this_.idGenerator = NewGoogleJsMessageIdGenerator(projectId)
	func() {
		defer func() {
			if rec_1 := recover(); rec_1 != nil {
				jex_2 := javalang.RecoverException(rec_1)
				_ = jex_2
				if j2gort.InstanceOf[j2gort.Ext_javax_xml_parsers_ParserConfigurationException](jex_2) || j2gort.InstanceOf[javaio.IIOException](jex_2) || j2gort.InstanceOf[j2gort.Ext_org_xml_sax_SAXException](jex_2) {
					e := jex_2
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_1)
				}
			}
		}()
		var parser j2gort.Ext_javax_xml_parsers_SAXParser = XtbMessageBundle_CreateSAXParser()
		_ = parser
		var reader j2gort.Ext_org_xml_sax_XMLReader = parser.GetXMLReader()
		_ = reader
		clinit_XtbMessageBundle_Handler()
		var val_3 XtbMessageBundle_Handler
		val_3.outer_ = this_
		InitXtbMessageBundle_Handler(&val_3)
		var contentHandler *XtbMessageBundle_Handler = &val_3
		_ = contentHandler
		reader.SetContentHandler(contentHandler)
		reader.Parse_1_Ext_org_xml_sax_InputSource(j2gort.NewExt_org_xml_sax_InputSource_1_Ext_java_io_InputStream(xtb))
	}()
}

// Java static method: XtbMessageBundle.createSAXParser()
func XtbMessageBundle_CreateSAXParser() j2gort.Ext_javax_xml_parsers_SAXParser {
	clinit_XtbMessageBundle()
	var factory j2gort.Ext_javax_xml_parsers_SAXParserFactory = j2gort.Ext_javax_xml_parsers_SAXParserFactory_NewInstance_0()
	_ = factory
	factory.SetValidating(false)
	factory.SetXIncludeAware(false)
	factory.SetFeature("http://xml.org/sax/features/external-general-entities", false)
	factory.SetFeature("http://xml.org/sax/features/external-parameter-entities", false)
	factory.SetFeature("http://apache.org/xml/features/nonvalidating/load-external-dtd", false)
	factory.SetFeature("http://javax.xml.XMLConstants/feature/secure-processing", true)
	var parser j2gort.Ext_javax_xml_parsers_SAXParser = factory.NewSAXParser()
	_ = parser
	var xmlReader j2gort.Ext_org_xml_sax_XMLReader = parser.GetXMLReader()
	_ = xmlReader
	xmlReader.SetEntityResolver(j2gort.ToIface_Ext_org_xml_sax_EntityResolver(StaticGet_XtbMessageBundle_NOOP_RESOLVER()))
	return parser
}

// Java method: XtbMessageBundle.getMessage(java.lang.String id)
func (this_ *XtbMessageBundle) GetMessage(id string) *JsMessage {
	_ = id
	return this_.messages.Get(id)
}

// Java method: XtbMessageBundle.idGenerator()
func (this_ *XtbMessageBundle) IdGenerator() JsMessage_IdGenerator {
	return this_.idGenerator
}

// Java method: XtbMessageBundle.getAllMessages()
func (this_ *XtbMessageBundle) GetAllMessages() any {
	return j2gort.Ext_com_google_common_collect_Iterables_S_UnmodifiableIterable_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.messages.LiveValues()))
}

// Java class: com.google.javascript.jscomp.XtbMessageBundle.Handler
type XtbMessageBundle_Handler struct {
	outer_     *XtbMessageBundle
	genderCase string             `gson:"genderCase"`
	lang       string             `gson:"lang"`
	msgBuilder *JsMessage_Builder `gson:"msgBuilder"`
	idHash_    atomic.Int32
}

func (this_ *XtbMessageBundle_Handler) BindSelf_(s any) {
	_ = s
}
func (this_ *XtbMessageBundle_Handler) IsJavaNil_() bool { return this_ == nil }
func (*XtbMessageBundle_Handler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.XtbMessageBundle$Handler", "Handler"
}
func (this_ *XtbMessageBundle_Handler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*XtbMessageBundle_Handler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.XtbMessageBundle.Handler", "com.google.javascript.jscomp.XtbMessageBundle$Handler", func(v_ any) bool {
		return j2gort.InstanceOf[*XtbMessageBundle_Handler](v_)
	})
}
func (this_ *XtbMessageBundle_Handler) SetOuter_(o *XtbMessageBundle)               { this_.outer_ = o }
func (this_ *XtbMessageBundle_Handler) GetOuter_() *XtbMessageBundle                { return this_.outer_ }
func (*XtbMessageBundle_Handler) Implements_j2gort_Ext_org_xml_sax_ContentHandler() {}

func StaticGet_XtbMessageBundle_Handler_BUNDLE_ELEM_NAME() string          { return "translationbundle" }
func StaticGet_XtbMessageBundle_Handler_LANG_ATT_NAME() string             { return "lang" }
func StaticGet_XtbMessageBundle_Handler_TRANSLATION_ELEM_NAME() string     { return "translation" }
func StaticGet_XtbMessageBundle_Handler_MESSAGE_ID_ATT_NAME() string       { return "id" }
func StaticGet_XtbMessageBundle_Handler_PLACEHOLDER_ELEM_NAME() string     { return "ph" }
func StaticGet_XtbMessageBundle_Handler_PLACEHOLDER_NAME_ATT_NAME() string { return "name" }
func StaticGet_XtbMessageBundle_Handler_BRANCH_ELEM_NAME() string          { return "branch" }
func StaticGet_XtbMessageBundle_Handler_BRANCH_NAME_ATT_NAME() string      { return "variants" }
func StaticGet_XtbMessageBundle_Handler_MALE_GENDER_CASE() string          { return "MASCULINE" }
func StaticGet_XtbMessageBundle_Handler_FEMALE_GENDER_CASE() string        { return "FEMININE" }
func StaticGet_XtbMessageBundle_Handler_NEUTER_GENDER_CASE() string        { return "NEUTER" }
func StaticGet_XtbMessageBundle_Handler_OTHER_GENDER_CASE() string         { return "OTHER" }
func StaticGet_XtbMessageBundle_Handler_GENDER_VARIANT_PATTERN() string {
	return "grammatical_gender_case:\\s*(\\w+)"
}

var static_XtbMessageBundle_Handler_GET_STRING_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_XtbMessageBundle_Handler_GET_STRING_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_XtbMessageBundle_Handler()
	return static_XtbMessageBundle_Handler_GET_STRING_PATTERN
}
func StaticSet_XtbMessageBundle_Handler_GET_STRING_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_XtbMessageBundle_Handler()
	static_XtbMessageBundle_Handler_GET_STRING_PATTERN = v
}

var clinitDone_XtbMessageBundle_Handler uint32
var clinitState_XtbMessageBundle_Handler j2gort.ClinitState

func clinit_XtbMessageBundle_Handler() {
	if atomic.LoadUint32(&clinitDone_XtbMessageBundle_Handler) != 0 {
		return
	}
	clinitSlow_XtbMessageBundle_Handler()
}

func clinitSlow_XtbMessageBundle_Handler() {
	if !j2gort.EnterClinit(&clinitState_XtbMessageBundle_Handler, &clinitDone_XtbMessageBundle_Handler) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_XtbMessageBundle_Handler, &clinitDone_XtbMessageBundle_Handler)
	static_XtbMessageBundle_Handler_GET_STRING_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("grammatical_gender_case:\\s*(\\w+)")
}

func (this_ *XtbMessageBundle_Handler) initFields_XtbMessageBundle_Handler() {
	_ = this_
	this_.genderCase = j2gort.NullString
	this_.lang = j2gort.NullString
}

// Java constructor: Handler()
func NewXtbMessageBundle_Handler(outer_ *XtbMessageBundle) *XtbMessageBundle_Handler {
	clinit_XtbMessageBundle_Handler()
	this_ := &XtbMessageBundle_Handler{}
	this_.outer_ = outer_
	InitXtbMessageBundle_Handler(this_)
	return this_
}

func InitXtbMessageBundle_Handler(this_ *XtbMessageBundle_Handler) {
	this_.initFields_XtbMessageBundle_Handler()
}

// Java method: Handler.setDocumentLocator(org.xml.sax.Locator locator)
func (this_ *XtbMessageBundle_Handler) SetDocumentLocator(locator j2gort.Ext_org_xml_sax_Locator) {
	_ = locator
}

// Java method: Handler.startDocument()
func (this_ *XtbMessageBundle_Handler) StartDocument() {
}

// Java method: Handler.endDocument()
func (this_ *XtbMessageBundle_Handler) EndDocument() {
}

// Java method: Handler.startPrefixMapping(java.lang.String prefix, java.lang.String uri)
func (this_ *XtbMessageBundle_Handler) StartPrefixMapping(prefix string, uri string) {
	_ = prefix
	_ = uri
}

// Java method: Handler.endPrefixMapping(java.lang.String prefix)
func (this_ *XtbMessageBundle_Handler) EndPrefixMapping(prefix string) {
	_ = prefix
}

// Java method: Handler.startElement(java.lang.String uri, java.lang.String localName, java.lang.String qName, org.xml.sax.Attributes atts)
func (this_ *XtbMessageBundle_Handler) StartElement(uri string, localName string, qName string, atts j2gort.Ext_org_xml_sax_Attributes) {
	_ = uri
	_ = localName
	_ = qName
	_ = atts
	switch j2gort.CheckNotNullString(qName) {
	case "translationbundle":
		if !(this_.lang == j2gort.NullString) {
			j2gort.CheckState(false)
		}
		this_.lang = atts.GetValue_1_string("lang")
		if !((this_.lang != j2gort.NullString) && !(j2gort.CheckNotNullString(this_.lang) == "")) {
			j2gort.CheckState(false)
		}
	case "translation":
		if !(this_.msgBuilder == nil) {
			j2gort.CheckState(false)
		}
		var id string = atts.GetValue_1_string("id")
		_ = id
		if !((id != j2gort.NullString) && !(j2gort.CheckNotNullString(id) == "")) {
			j2gort.CheckState(false)
		}
		this_.msgBuilder = NewJsMessage_Builder().SetKey(id).SetId(id)
	case "ph":
		if !(this_.msgBuilder != nil) {
			j2gort.CheckState(false)
		}
		var phRef string = atts.GetValue_1_string("name")
		_ = phRef
		if this_.msgBuilder.HasGenderedVariants() {
			this_.msgBuilder.AppendCanonicalPlaceholderReference_2(JsMessage_GrammaticalGenderCase_ValueOf(this_.genderCase), phRef)
		} else {
			this_.msgBuilder.AppendCanonicalPlaceholderReference_1(phRef)
		}
	case "branch":
		if !(this_.msgBuilder != nil) {
			j2gort.CheckState(false)
		}
		var gender string = atts.GetValue_1_string("variants")
		_ = gender
		var matcher *j2gort.Ext_java_util_regex_Matcher = StaticGet_XtbMessageBundle_Handler_GET_STRING_PATTERN().Matcher(j2gort.BoxString(gender))
		_ = matcher
		if matcher.Find_0() {
			this_.genderCase = matcher.Group_1_int32(int32(1))
			{
				var iifeCond_1 bool
				{

					sr_4 := this_.genderCase
					se_5 := "MASCULINE"
					iifeCond_1 = j2gort.CheckNotNullString(sr_4) == se_5
				}
				if !(((iifeCond_1 || func() bool {
					sr_6 := this_.genderCase
					se_7 := "FEMININE"
					return j2gort.CheckNotNullString(sr_6) == se_7
				}()) || func() bool {
					sr_8 := this_.genderCase
					se_9 := "NEUTER"
					return j2gort.CheckNotNullString(sr_8) == se_9
				}()) || func() bool {
					sr_10 := this_.genderCase
					se_11 := "OTHER"
					return j2gort.CheckNotNullString(sr_10) == se_11
				}()) {
					j2gort.CheckStateStr(false, "Gender case must be one of the following: MASCULINE, FEMININE, NEUTER, or OTHER.")
				}
			}
			this_.msgBuilder.AddGenderedMessageKey(JsMessage_GrammaticalGenderCase_ValueOf(this_.genderCase))
		}
	default:
	}
}

// Java method: Handler.endElement(java.lang.String uri, java.lang.String localName, java.lang.String qName)
func (this_ *XtbMessageBundle_Handler) EndElement(uri string, localName string, qName string) {
	_ = uri
	_ = localName
	_ = qName
	if "translation" == qName {
		if !(this_.msgBuilder != nil) {
			j2gort.CheckState(false)
		}
		if !(this_.msgBuilder.HasParts()) {
			this_.msgBuilder.AppendStringPart_1("")
		}
		var key string = this_.msgBuilder.GetKey()
		_ = key
		this_.outer_.messages.Put(key, this_.msgBuilder.Build())
		this_.msgBuilder = nil
	}
}

// Java method: Handler.characters(char[] ch, int start, int length)
func (this_ *XtbMessageBundle_Handler) Characters(ch []uint16, start int32, length int32) {
	_ = ch
	_ = start
	_ = length
	if this_.msgBuilder != nil {
		var part string = j2gort.Ext_java_lang_String_ValueOf_3(ch, start, length)
		_ = part
		if this_.msgBuilder.HasGenderedVariants() {
			this_.msgBuilder.AppendStringPart_2(JsMessage_GrammaticalGenderCase_ValueOf(this_.genderCase), part)
		} else {
			this_.msgBuilder.AppendStringPart_1(part)
		}
	}
}

// Java method: Handler.ignorableWhitespace(char[] ch, int start, int length)
func (this_ *XtbMessageBundle_Handler) IgnorableWhitespace(ch []uint16, start int32, length int32) {
	_ = ch
	_ = start
	_ = length
	if this_.msgBuilder != nil {
		if this_.msgBuilder.HasGenderedVariants() {
			this_.msgBuilder.AppendStringPart_2(JsMessage_GrammaticalGenderCase_ValueOf(this_.genderCase), j2gort.Ext_java_lang_String_ValueOf_3(ch, start, length))
		} else {
			this_.msgBuilder.AppendStringPart_1(j2gort.Ext_java_lang_String_ValueOf_3(ch, start, length))
		}
	}
}

// Java method: Handler.processingInstruction(java.lang.String target, java.lang.String data)
func (this_ *XtbMessageBundle_Handler) ProcessingInstruction(target string, data string) {
	_ = target
	_ = data
}

// Java method: Handler.skippedEntity(java.lang.String name)
func (this_ *XtbMessageBundle_Handler) SkippedEntity(name string) {
	_ = name
}

// Java class: com.google.javascript.jscomp.XtbMessageBundle.SecureEntityResolver
type XtbMessageBundle_SecureEntityResolver struct {
	idHash_ atomic.Int32
}

func (this_ *XtbMessageBundle_SecureEntityResolver) BindSelf_(s any) {
	_ = s
}
func (this_ *XtbMessageBundle_SecureEntityResolver) IsJavaNil_() bool { return this_ == nil }
func (*XtbMessageBundle_SecureEntityResolver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.XtbMessageBundle$SecureEntityResolver", "SecureEntityResolver"
}
func (this_ *XtbMessageBundle_SecureEntityResolver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*XtbMessageBundle_SecureEntityResolver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.XtbMessageBundle.SecureEntityResolver", "com.google.javascript.jscomp.XtbMessageBundle$SecureEntityResolver", func(v_ any) bool {
		return j2gort.InstanceOf[*XtbMessageBundle_SecureEntityResolver](v_)
	})
}
func (*XtbMessageBundle_SecureEntityResolver) Implements_j2gort_Ext_org_xml_sax_EntityResolver() {}

func (this_ *XtbMessageBundle_SecureEntityResolver) initFields_XtbMessageBundle_SecureEntityResolver() {
	_ = this_
}

// Java constructor: SecureEntityResolver()
func NewXtbMessageBundle_SecureEntityResolver() *XtbMessageBundle_SecureEntityResolver {
	this_ := &XtbMessageBundle_SecureEntityResolver{}
	InitXtbMessageBundle_SecureEntityResolver(this_)
	return this_
}

func InitXtbMessageBundle_SecureEntityResolver(this_ *XtbMessageBundle_SecureEntityResolver) {
	this_.initFields_XtbMessageBundle_SecureEntityResolver()
}

// Java method: SecureEntityResolver.resolveEntity(java.lang.String publicId, java.lang.String systemId)
func (this_ *XtbMessageBundle_SecureEntityResolver) ResolveEntity(publicId string, systemId string) *j2gort.Ext_org_xml_sax_InputSource {
	_ = publicId
	_ = systemId
	return j2gort.NewExt_org_xml_sax_InputSource_1_Ext_java_io_Reader(j2gort.NewExt_java_io_StringReader(""))
}

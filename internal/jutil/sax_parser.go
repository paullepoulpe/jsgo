package jutil

// Go-only: language-adapter — the javax.xml.parsers / org.xml.sax API that
// XtbMessageBundle uses (SAXParserFactory, SAXParser, XMLReader,
// ContentHandler, Attributes, Locator, EntityResolver, InputSource and the
// SAX exceptions), declared from the public Javadoc, over a clean-room XML
// parser (sax_scanner.go, sax_input.go) written from the XML 1.0
// specification and checked black-box against the JDK's SAXParser
// (tools/oracle/java/JutilSaxVectors.java, testdata/sax_vectors.txt.gz).
//
// Only the configuration XtbMessageBundle uses is emulated: no validation,
// no XInclude, no namespaces, no external entities and no external DTD,
// secure processing on.

import (
	"io"
	"strconv"
)

// XMLConstants_FEATURE_SECURE_PROCESSING is
// javax.xml.XMLConstants.FEATURE_SECURE_PROCESSING.
const XMLConstants_FEATURE_SECURE_PROCESSING = "http://javax.xml.XMLConstants/feature/secure-processing"

const (
	saxFeatureExternalGeneralEntities   = "http://xml.org/sax/features/external-general-entities"
	saxFeatureExternalParameterEntities = "http://xml.org/sax/features/external-parameter-entities"
	saxFeatureLoadExternalDtd           = "http://apache.org/xml/features/nonvalidating/load-external-dtd"
)

// SAXException is org.xml.sax.SAXException.
type SAXException struct{ JavaException }

func (*SAXException) saxException() {}

// ISAXException matches SAXException and its subclasses.
type ISAXException interface {
	IException
	saxException()
}

// SAXParseException is org.xml.sax.SAXParseException. The parser's
// exceptions have no public or system id (the InputSource has none).
type SAXParseException struct {
	SAXException
	lineNumber   int
	columnNumber int
}

func (*SAXParseException) saxParseException() {}

// NewSAXParseException is new SAXParseException(String, Locator) without a
// locator.
func NewSAXParseException(message string) *SAXParseException {
	return newSAXParseExceptionAt(message, -1, -1)
}

func newSAXParseExceptionAt(message string, line, column int) *SAXParseException {
	e := &SAXParseException{lineNumber: line, columnNumber: column}
	InitThrowable(e, "org.xml.sax.SAXParseException", &message, nil)
	return e
}

// GetLineNumber is SAXParseException.getLineNumber(): -1 when unknown.
func (e *SAXParseException) GetLineNumber() int { return e.lineNumber }

// GetColumnNumber is SAXParseException.getColumnNumber(): -1 when unknown.
func (e *SAXParseException) GetColumnNumber() int { return e.columnNumber }

// Error is SAXParseException.toString(), which the Javadoc only describes
// as more detailed. Observed behaviour (black-box): the class name, then
// "; lineNumber: L" and "; columnNumber: C" for the known ones, then
// "; " and the message.
func (e *SAXParseException) Error() string {
	s := "org.xml.sax.SAXParseException"
	if e.lineNumber != -1 {
		s += "; lineNumber: " + strconv.Itoa(e.lineNumber)
	}
	if e.columnNumber != -1 {
		s += "; columnNumber: " + strconv.Itoa(e.columnNumber)
	}
	if e.hasMsg {
		s += "; " + e.message
	}
	return s
}

// ParserConfigurationException is
// javax.xml.parsers.ParserConfigurationException.
type ParserConfigurationException struct{ JavaException }

func (*ParserConfigurationException) parserConfigurationException() {}

// Locator is org.xml.sax.Locator.
type Locator interface {
	GetPublicId() *string
	GetSystemId() *string
	GetLineNumber() int
	GetColumnNumber() int
}

// Attributes is org.xml.sax.Attributes (the queries by index and by qName).
type Attributes interface {
	GetLength() int
	GetQName(index int) string
	GetValueAt(index int) string
	// GetValue is getValue(String qName): nil when there is no such
	// attribute.
	GetValue(qName string) *string
}

// ContentHandler is org.xml.sax.ContentHandler. Characters and
// IgnorableWhitespace receive UTF-16 code units, like Java's char[].
type ContentHandler interface {
	SetDocumentLocator(locator Locator)
	StartDocument()
	EndDocument()
	StartPrefixMapping(prefix string, uri string)
	EndPrefixMapping(prefix string)
	StartElement(uri string, localName string, qName string, atts Attributes)
	EndElement(uri string, localName string, qName string)
	Characters(ch []uint16, start int, length int)
	IgnorableWhitespace(ch []uint16, start int, length int)
	ProcessingInstruction(target string, data string)
	SkippedEntity(name string)
}

// EntityResolver is org.xml.sax.EntityResolver.
type EntityResolver interface {
	ResolveEntity(publicId *string, systemId *string) *InputSource
}

// InputSource is org.xml.sax.InputSource with a byte stream or a character
// stream.
type InputSource struct {
	byteStream      io.Reader
	characterStream io.Reader
}

// NewInputSource is new InputSource(InputStream).
func NewInputSource(byteStream io.Reader) *InputSource {
	return &InputSource{byteStream: byteStream}
}

// NewInputSourceOfCharacterStream is new InputSource(Reader); the Go reader
// yields UTF-8 text.
func NewInputSourceOfCharacterStream(characterStream io.Reader) *InputSource {
	return &InputSource{characterStream: characterStream}
}

// SAXParserFactory is javax.xml.parsers.SAXParserFactory.
type SAXParserFactory struct {
	validating    bool
	xIncludeAware bool
	features      map[string]bool
}

// SAXParserFactory_NewInstance is SAXParserFactory.newInstance().
func SAXParserFactory_NewInstance() *SAXParserFactory {
	return &SAXParserFactory{features: map[string]bool{
		saxFeatureExternalGeneralEntities:      true,
		saxFeatureExternalParameterEntities:    true,
		saxFeatureLoadExternalDtd:              true,
		XMLConstants_FEATURE_SECURE_PROCESSING: true,
	}}
}

// SetValidating is SAXParserFactory.setValidating(boolean).
func (f *SAXParserFactory) SetValidating(validating bool) { f.validating = validating }

// SetXIncludeAware is SAXParserFactory.setXIncludeAware(boolean).
func (f *SAXParserFactory) SetXIncludeAware(state bool) { f.xIncludeAware = state }

// SetFeature is SAXParserFactory.setFeature(String, boolean).
func (f *SAXParserFactory) SetFeature(name string, value bool) {
	if _, ok := f.features[name]; !ok {
		panic(NewUnsupportedOperationExceptionMsg("SAX feature not emulated: " + name))
	}
	f.features[name] = value
}

// NewSAXParser is SAXParserFactory.newSAXParser().
func (f *SAXParserFactory) NewSAXParser() *SAXParser {
	if f.validating || f.xIncludeAware ||
		f.features[saxFeatureExternalGeneralEntities] ||
		f.features[saxFeatureExternalParameterEntities] ||
		f.features[saxFeatureLoadExternalDtd] {
		panic(NewUnsupportedOperationExceptionMsg(
			"SAXParser emulated only without validation, XInclude, external entities and external DTD loading"))
	}
	return &SAXParser{reader: &XMLReader{}}
}

// SAXParser is javax.xml.parsers.SAXParser.
type SAXParser struct {
	reader *XMLReader
}

// GetXMLReader is SAXParser.getXMLReader().
func (p *SAXParser) GetXMLReader() *XMLReader { return p.reader }

// XMLReader is org.xml.sax.XMLReader.
type XMLReader struct {
	contentHandler ContentHandler
	entityResolver EntityResolver
}

// SetEntityResolver is XMLReader.setEntityResolver(EntityResolver).
// Observed behaviour (black-box): with the external DTD and external
// entities off, the JDK never calls the resolver, so neither does this
// parser.
func (r *XMLReader) SetEntityResolver(resolver EntityResolver) { r.entityResolver = resolver }

// SetContentHandler is XMLReader.setContentHandler(ContentHandler).
func (r *XMLReader) SetContentHandler(handler ContentHandler) { r.contentHandler = handler }

// Parse is XMLReader.parse(InputSource): it reports the document to the
// ContentHandler and throws SAXParseException at the first well-formedness
// error (see sax_scanner.go).
func (r *XMLReader) Parse(input *InputSource) {
	newSaxScanner(r.contentHandler, input).parse()
}

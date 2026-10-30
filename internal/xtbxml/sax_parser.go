package xtbxml

// Go-only: language-adapter — clean-room XML 1.0 streaming parser
// (sax_scanner.go, sax_input.go) checked black-box against the JDK's SAXParser
// (tools/oracle/java/JutilSaxVectors.java, testdata/sax_vectors.txt.gz).
//
// Only the secure configuration XtbMessageBundle uses is supported: no
// validation, no XInclude, no namespaces, no external entities and no
// external DTD.

import (
	"io"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SAXException is org.xml.sax.SAXException.
type SAXException struct{ jutil.JavaException }

func (*SAXException) saxException() {}

// ISAXException matches SAXException and its subclasses.
type ISAXException interface {
	jutil.IException
	saxException()
}

// SAXParseException is org.xml.sax.SAXParseException.
type SAXParseException struct {
	SAXException
	lineNumber   int
	columnNumber int
}

func (*SAXParseException) saxParseException() {}

func newSAXParseExceptionAt(message string, line, column int) *SAXParseException {
	e := &SAXParseException{lineNumber: line, columnNumber: column}
	jutil.InitThrowable(e, "org.xml.sax.SAXParseException", &message, nil)
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
	if e.HasMessage() {
		s += "; " + e.GetMessage()
	}
	return s
}

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

// DefaultHandler provides no-op implementations of ContentHandler.
type DefaultHandler struct{}

func (DefaultHandler) SetDocumentLocator(Locator)                      {}
func (DefaultHandler) StartDocument()                                  {}
func (DefaultHandler) EndDocument()                                    {}
func (DefaultHandler) StartPrefixMapping(string, string)               {}
func (DefaultHandler) EndPrefixMapping(string)                         {}
func (DefaultHandler) StartElement(string, string, string, Attributes) {}
func (DefaultHandler) EndElement(string, string, string)               {}
func (DefaultHandler) Characters([]uint16, int, int)                   {}
func (DefaultHandler) IgnorableWhitespace([]uint16, int, int)          {}
func (DefaultHandler) ProcessingInstruction(string, string)            {}
func (DefaultHandler) SkippedEntity(string)                            {}

type saxBailout struct {
	err jutil.Throwable
}

// Parse streams XML 1.0 events from r to handler and returns the exception
// the JDK's SAXParser would throw at the first well-formedness, encoding, or
// I/O error.
func Parse(r io.Reader, handler ContentHandler) (err jutil.Throwable) {
	defer func() {
		if rec := recover(); rec != nil {
			if b, ok := rec.(saxBailout); ok {
				err = b.err
				return
			}
			panic(rec)
		}
	}()
	newSaxScanner(handler, r).parse()
	return nil
}

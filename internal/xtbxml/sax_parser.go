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

func newSAXParseExceptionAt(message string, line, column int) *SAXParseException {
	e := &SAXParseException{lineNumber: line, columnNumber: column}
	jutil.InitThrowableMsg(e, "org.xml.sax.SAXParseException", message, nil)
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
	GetPublicId() (string, bool)
	GetSystemId() (string, bool)
	GetLineNumber() int
	GetColumnNumber() int
}

// saxLocator is the Locator given to the ContentHandler: the scanner's
// current position, and -1 once the document has ended.
type saxLocator struct{ s *saxScanner }

func (l *saxLocator) GetPublicId() (string, bool) { return "", false }
func (l *saxLocator) GetSystemId() (string, bool) { return "", false }
func (l *saxLocator) GetLineNumber() int {
	if l.s.ended {
		return -1
	}
	return l.s.in.line
}
func (l *saxLocator) GetColumnNumber() int {
	if l.s.ended {
		return -1
	}
	return l.s.in.column
}

// Attributes is org.xml.sax.Attributes (the queries by index and by qName).
type Attributes interface {
	GetLength() int
	GetQName(index int) string
	GetValueAt(index int) string
	// Value returns the value for qName and whether the attribute is present.
	Value(qName string) (string, bool)
}

// saxAttributes is the Attributes of a start tag, in document order.
type saxAttributes struct {
	names  []string
	values []string
}

func (a *saxAttributes) GetLength() int { return len(a.names) }

func (a *saxAttributes) GetQName(index int) string {
	if index < 0 || index >= len(a.names) {
		return ""
	}
	return a.names[index]
}

func (a *saxAttributes) GetValueAt(index int) string {
	if index < 0 || index >= len(a.values) {
		return ""
	}
	return a.values[index]
}

func (a *saxAttributes) Value(qName string) (string, bool) {
	for i, n := range a.names {
		if n == qName {
			return a.values[i], true
		}
	}
	return "", false
}

// ContentHandler is org.xml.sax.ContentHandler (namespace-unaware subset).
// Characters and IgnorableWhitespace receive UTF-16 code units, like Java's char[].
type ContentHandler interface {
	SetDocumentLocator(locator Locator)
	StartDocument()
	EndDocument()
	StartElement(qName string, atts Attributes)
	EndElement(qName string)
	Characters(ch []uint16, start int, length int)
	IgnorableWhitespace(ch []uint16, start int, length int)
	ProcessingInstruction(target string, data string)
	SkippedEntity(name string)
}

// DefaultHandler provides no-op implementations of ContentHandler.
type DefaultHandler struct{}

func (DefaultHandler) SetDocumentLocator(Locator)             {}
func (DefaultHandler) StartDocument()                         {}
func (DefaultHandler) EndDocument()                           {}
func (DefaultHandler) StartElement(string, Attributes)        {}
func (DefaultHandler) EndElement(string)                      {}
func (DefaultHandler) Characters([]uint16, int, int)          {}
func (DefaultHandler) IgnorableWhitespace([]uint16, int, int) {}
func (DefaultHandler) ProcessingInstruction(string, string)   {}
func (DefaultHandler) SkippedEntity(string)                   {}

// Parse streams XML 1.0 events from r to handler and returns the exception
// the JDK's SAXParser would throw at the first well-formedness, encoding, or
// I/O error.
func Parse(r io.Reader, handler ContentHandler) jutil.Throwable {
	return newSaxScanner(handler, r).parse()
}

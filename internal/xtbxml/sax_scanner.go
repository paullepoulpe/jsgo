package xtbxml

// Go-only: language-adapter — the XML scanner behind XMLReader.Parse
// (sax_parser.go): a non-validating, namespace-unaware XML 1.0 processor
// for documents without external entities, written from the XML 1.0
// specification:
//   - §2.8 prolog and document type declaration, §4.3.3 encoding
//     declaration, §2.9 standalone declaration;
//   - §3.1 start-tags, end-tags and empty-element tags, §3.3.3 attribute
//     value normalization (attributes are CDATA: no declarations are read);
//   - §2.4 character data, §2.7 CDATA sections, §2.5 comments, §2.6
//     processing instructions, §4.1 character and entity references and
//     §4.6 the predefined entities, §2.11 end-of-line handling;
//   - §4.1 WFC: Entity Declared — references to undeclared entities are
//     reported through ContentHandler.skippedEntity when the document has
//     an external subset the parser does not read and is not standalone.
//
// What XML 1.0 leaves open and the JDK makes observable was determined
// black-box with tools/oracle/java/JutilSaxVectors.java:
//   - the error messages and the Locator position of each error and event;
//   - where character data is cut into ContentHandler.characters() calls.
//     Observed behaviour (black-box): a run of text is scanned in steps
//     that stop at markup ('<', '&'), at ']', at a line end, at a
//     character outside the BMP and at the end of the input buffer
//     (sax_input.go). One characters() call reports at most two such steps
//     (with the ']' runs after each), a supplementary character ends a
//     call, and a reference is reported alone. A step that starts with
//     line ends takes the whole run of them, but stops after a line end
//     that leaves at most one char in the buffer.
//   - Names follow the character classes of XML 1.0 (Fourth Edition),
//     Appendix B, not the Fifth Edition's ranges (determined black-box by
//     trying every BMP character; see saxNameStartRanges).
//   - The limits of the JDK's secure processing: names of more than 1000
//     chars and elements with more than 10000 attributes are errors.
//
// Not emulated (UnsupportedOperationException): XML 1.1 documents,
// encodings other than UTF-8, UTF-16, ISO-8859-1 and US-ASCII, and the
// internal-subset declarations that would change the events: references to
// internal entities (general or parameter), attribute-list declarations
// with attribute definitions, and element declarations with element
// content (whitespace in such elements is ignorable).

import (
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	saxStopEOF    = -1
	saxStopBufEnd = -2
)

const saxMsgEOFInEntity = "XML document structures must start and end within the same entity."

type saxEntityDecl struct {
	external bool
	unparsed bool
}

type saxScanner struct {
	in      *saxInput
	handler ContentHandler
	locator saxLocator

	text []uint16 // character data not yet reported
	val  []uint16 // attribute value, PI data, literal being scanned

	open  [][]uint16 // names of the open elements
	attrs saxAttributes

	hasDecl    bool
	inSubset   bool
	declStage  int
	versionPad int
	// versionLoose: the pre-scan stayed in the first read and consumed
	// more than "<?xml version=" (see preScanDeclaration).
	versionLoose   bool
	sawDoctype     bool
	externalSubset bool
	standalone     bool
	entities       map[string]saxEntityDecl
	paramEntities  map[string]saxEntityDecl
	elementContent map[string]bool
	ended          bool
}

// saxLocator is the Locator given to the ContentHandler: the scanner's
// current position, and -1 once the document has ended.
type saxLocator struct{ s *saxScanner }

func (l *saxLocator) GetPublicId() *string { return nil }
func (l *saxLocator) GetSystemId() *string { return nil }
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

func (a *saxAttributes) GetValue(qName string) *string {
	for i, n := range a.names {
		if n == qName {
			return &a.values[i]
		}
	}
	return nil
}

type saxNoopHandler struct{}

func (saxNoopHandler) SetDocumentLocator(Locator)                      {}
func (saxNoopHandler) StartDocument()                                  {}
func (saxNoopHandler) EndDocument()                                    {}
func (saxNoopHandler) StartPrefixMapping(string, string)               {}
func (saxNoopHandler) EndPrefixMapping(string)                         {}
func (saxNoopHandler) StartElement(string, string, string, Attributes) {}
func (saxNoopHandler) EndElement(string, string, string)               {}
func (saxNoopHandler) Characters([]uint16, int, int)                   {}
func (saxNoopHandler) IgnorableWhitespace([]uint16, int, int)          {}
func (saxNoopHandler) ProcessingInstruction(string, string)            {}
func (saxNoopHandler) SkippedEntity(string)                            {}

func newSaxScanner(handler ContentHandler, source *InputSource) *saxScanner {
	if handler == nil {
		handler = saxNoopHandler{}
	}
	s := &saxScanner{in: newSaxInput(source), handler: handler}
	s.locator.s = s
	s.in.fail = s.fatal
	return s
}

// ---------------------------------------------------------------------------
// Errors.

func (s *saxScanner) fatal(msg string) {
	s.fatalAt(msg, s.in.line, s.in.column)
}

// fatalAt throws the SAXParseException of a fatal error. With no
// ErrorHandler (XMLReader has none here), the JDK also prints the error on
// System.err; observed behaviour (black-box): "[Fatal Error] :L:C: msg".
func (s *saxScanner) fatalAt(msg string, line, column int) {
	jutil.System_Err().Println("[Fatal Error] :" + strconv.Itoa(line) + ":" + strconv.Itoa(column) + ": " + msg)
	panic(newSAXParseExceptionAt(msg, line, column))
}

// fatalOrEOF reports msg, or the end-of-input error when the input has
// ended.
func (s *saxScanner) fatalOrEOF(msg string) {
	if s.peek() == saxStopEOF {
		s.fatal(s.eofMessage())
	}
	s.fatal(msg)
}

// eofMessage is the error for an input that ends inside a construct.
// Observed behaviour (black-box): inside the internal subset it is
// "Premature end of file.".
func (s *saxScanner) eofMessage() string {
	if s.inSubset {
		return "Premature end of file."
	}
	return saxMsgEOFInEntity
}

func saxUnsupported(what string) {
	panic(jutil.NewUnsupportedOperationExceptionMsg("SAX parser: " + what + " not emulated"))
}

func saxHex(c int) string {
	return "0x" + strconv.FormatInt(int64(c), 16)
}

// saxGrouped formats n with comma thousands separators, as the JDK's
// messages print numbers (in an English locale).
func saxGrouped(n int) string {
	t := strconv.Itoa(n)
	var b strings.Builder
	for i, c := range t {
		if i > 0 && (len(t)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// Character classes.

func saxIsSpace(c int) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// saxIsXMLChar is production [2] Char for a BMP code unit (surrogates are
// checked in pairs by the callers).
func saxIsXMLChar(c int) bool {
	return c == 0x9 || c == 0xA || c == 0xD || (c >= 0x20 && c <= 0xD7FF) || (c >= 0xE000 && c <= 0xFFFD)
}

// saxIsTextChar is a char that the character-data step scans through.
func saxIsTextChar(c uint16) bool {
	switch {
	case c >= 0x20 && c <= 0xD7FF:
		return c != '<' && c != '&' && c != ']'
	case c == 0x9:
		return true
	case c >= 0xE000:
		return c <= 0xFFFD
	}
	return false
}

func saxIsHighSurrogate(c int) bool { return c >= 0xD800 && c <= 0xDBFF }
func saxIsLowSurrogate(c int) bool  { return c >= 0xDC00 && c <= 0xDFFF }

var saxNameStartBits, saxNameCharBits = saxBuildNameBits(saxNameStartRanges[:]), saxBuildNameBits(saxNameCharRanges[:])

func saxBuildNameBits(ranges []uint16) []uint64 {
	bits := make([]uint64, 0x10000/64)
	for i := 0; i < len(ranges); i += 2 {
		for c := int(ranges[i]); c <= int(ranges[i+1]); c++ {
			bits[c>>6] |= 1 << (c & 63)
		}
	}
	return bits
}

func saxIsNameStart(c int) bool {
	return c >= 0 && c < 0x10000 && saxNameStartBits[c>>6]&(1<<(c&63)) != 0
}

func saxIsNameChar(c int) bool {
	return c >= 0 && c < 0x10000 && saxNameCharBits[c>>6]&(1<<(c&63)) != 0
}

// ---------------------------------------------------------------------------
// Scanning primitives. Each one reads more input exactly when the JDK's
// scanner does (see sax_input.go), because that decides where later
// character data is cut.

// peek returns the next char without consuming it, or saxStopEOF.
func (s *saxScanner) peek() int {
	in := s.in
	if in.pos == in.count && in.load(in.pos) == 0 {
		return saxStopEOF
	}
	return int(in.buf[in.pos])
}

// next consumes one char, normalizing a line end to '\n'. Observed
// behaviour (black-box): a line end that is the last char of the buffer is
// consumed, then kept in the buffer while more input is read.
func (s *saxScanner) next() int {
	in := s.in
	c := s.peek()
	switch c {
	case saxStopEOF:
	case '\n', '\r':
		in.newline(1)
		if in.pos == in.count {
			in.load(in.pos - 1)
		}
		if c == '\r' && in.pos < in.count && in.buf[in.pos] == '\n' {
			in.pos++
		}
		c = '\n'
	default:
		in.advance(1)
	}
	return c
}

// skipChar consumes c (not a line end) if it is next.
func (s *saxScanner) skipChar(c int) bool {
	if s.peek() != c {
		return false
	}
	s.in.advance(1)
	return true
}

// skipSpaces consumes white space (production [3] S); it reports whether
// there was any.
func (s *saxScanner) skipSpaces() bool {
	if !saxIsSpace(s.peek()) {
		return false
	}
	for saxIsSpace(s.peek()) {
		s.next()
	}
	return true
}

// skipUnits consumes str if it is next. It first makes sure the buffer
// holds len(str) chars, reading into the free end of the buffer.
func (s *saxScanner) skipUnits(str []uint16) bool {
	in := s.in
	for in.count-in.pos < len(str) {
		if in.loadMore() == 0 {
			return false
		}
	}
	for i, c := range str {
		if in.buf[in.pos+i] != c {
			return false
		}
	}
	in.advance(len(str))
	return true
}

func (s *saxScanner) skipASCII(str string) bool {
	in := s.in
	for in.count-in.pos < len(str) {
		if in.loadMore() == 0 {
			return false
		}
	}
	for i := 0; i < len(str); i++ {
		if in.buf[in.pos+i] != uint16(str[i]) {
			return false
		}
	}
	in.advance(len(str))
	return true
}

const saxMaxNameLength = 1000

// scanName scans a Name (production [5], Fourth Edition classes); it
// returns nil when no name starts here. A name that reaches the end of
// the buffer is kept in the buffer while more is read. Names longer than
// 1000 chars hit the JDK's secure-processing limit.
func (s *saxScanner) scanName() []uint16 {
	in := s.in
	if in.pos == in.count && in.load(in.pos) == 0 {
		return nil
	}
	if !saxIsNameStart(int(in.buf[in.pos])) {
		return nil
	}
	start := in.pos
	end := start + 1
	for {
		for end < in.count && saxIsNameChar(int(in.buf[end])) {
			end++
		}
		if end < in.count {
			break
		}
		n := end - start
		if n > saxMaxNameLength {
			// Observed behaviour (black-box): checked whenever the name
			// reaches the end of the buffer, before reading more.
			s.fatal("JAXP00010005: The length of entity \"[xml]\" is \"" + saxGrouped(n) +
				"\" that exceeds the \"1,000\" limit set by \"jaxp.properties\".")
		}
		read := in.load(start)
		start = in.pos
		end = start + n
		if read == 0 {
			break
		}
	}
	name := append([]uint16(nil), in.buf[start:end]...)
	in.advance(end - start)
	if len(name) > saxMaxNameLength {
		s.fatal("JAXP00010005: The length of entity \"[xml]\" is \"" + saxGrouped(len(name)) +
			"\" that exceeds the \"1,000\" limit set by \"jaxp.properties\".")
	}
	return name
}

// flushText reports the pending character data.
func (s *saxScanner) flushText() {
	if len(s.text) > 0 {
		s.handler.Characters(s.text, 0, len(s.text))
		s.text = s.text[:0]
	}
}

// ---------------------------------------------------------------------------
// Document.

func (s *saxScanner) parse() {
	s.in.detectEncoding()
	s.preScanDeclaration()
	// Observed behaviour (black-box): the Locator reads 1:1 during
	// setDocumentLocator and startDocument, whatever the pre-scan consumed.
	line, column := s.in.line, s.in.column
	s.in.line, s.in.column = 1, 1
	s.handler.SetDocumentLocator(&s.locator)
	s.handler.StartDocument()
	s.in.line, s.in.column = line, column
	if s.hasDecl {
		s.xmlDecl()
	}
	s.prolog()
	s.trailingMisc()
	s.ended = true
	s.handler.EndDocument()
}

// preScanDeclaration looks for an XML declaration and its version before
// the document is reported (the JDK reads the version first, to choose
// between XML 1.0 and 1.1). It consumes "<?xml", white space, "version",
// white space, '=' and white space, as far as they are there, then looks
// at the next five chars (the version literal and one more) without
// consuming them; the declaration is then scanned from where it stopped.
// Observed behaviour (black-box):
//   - the end of input during the pre-scan is a "Premature end of file."
//     error at line and column -1, before setDocumentLocator; a decoding
//     error is reported where the pre-scan met it, also before;
//   - after the pre-scan the position is line 1, column one more than
//     the number of chars consumed: line ends in the pre-scanned part
//     count as columns. But when the version literal is not "1.0" in a
//     matching pair of quotes, the column at the literal's quote is 15
//     (that of "<?xml version='"), and a literal of four chars or more
//     gets as many spaces after its fourth char as the pre-scanned part is
//     longer than "<?xml version=" (the spaces count as columns too);
//   - when the pre-scan (look-ahead included) had to read past the first
//     read of the input (32 bytes, so always in UTF-16), the column is
//     that of the shortest declaration in the same state ("<?xml ",
//     "<?xml version", "<?xml version="), and there are no spaces;
//   - when the pre-scan stayed in the first read and consumed more than
//     "<?xml version=", the pseudo attribute after the version needs no
//     white space before it;
//   - when "<?xml" is followed by a name char instead (a processing
//     instruction such as "<?xml-stylesheet"), the scan restarts at the
//     beginning but keeps the five columns it skipped, so positions on
//     line 1 are five columns too far.
func (s *saxScanner) preScanDeclaration() {
	in := s.in
	if !s.skipASCII("<?xml") {
		return
	}
	eof := func() { s.fatalAt("Premature end of file.", -1, -1) }
	c := s.peek()
	if c == saxStopEOF {
		eof()
	}
	if !saxIsSpace(c) {
		in.pos = 0
		return
	}
	s.hasDecl = true
	in.declMode = true
	firstEnd := in.base + in.count
	beyond := func() bool { return in.base+in.count > firstEnd }
	spaces := func() {
		s.skipSpaces()
		if s.peek() == saxStopEOF {
			eof()
		}
	}
	// finish sets the position after the pre-scan; shortest is the length
	// of the shortest declaration prefix in the same state.
	finish := func(shortest int) {
		in.line, in.column = 1, in.base+in.pos+1
		if beyond() {
			in.column = shortest + 1
		}
	}
	spaces()
	if !s.skipASCII("version") {
		finish(len("<?xml "))
		return
	}
	s.declStage = 1
	spaces()
	if s.peek() != '=' {
		finish(len("<?xml version"))
		return
	}
	in.advance(1)
	s.declStage = 2
	spaces()
	// The literal is read one char at a time: a decoding error is
	// reported after the chars read so far.
	for k := 1; k <= 5; k++ {
		for in.count-in.pos < k {
			line, column := in.line, in.column
			for _, u := range in.buf[in.pos : in.pos+k-1] {
				if u == '\n' || u == '\r' {
					in.line, in.column = in.line+1, 1
				} else {
					in.column++
				}
			}
			n := in.loadMore()
			in.line, in.column = line, column
			if n == 0 {
				eof()
			}
		}
	}
	b := in.buf[in.pos : in.pos+5]
	quoted := (b[0] == '"' || b[0] == '\'') && b[4] == b[0]
	if quoted && b[1] == '1' && b[2] == '.' && b[3] == '1' {
		saxUnsupported("XML 1.1 documents are")
	}
	s.versionLoose = !beyond() && in.base+in.pos > len("<?xml version=")
	if quoted && b[1] == '1' && b[2] == '.' && b[3] == '0' {
		finish(len("<?xml version="))
		return
	}
	if !beyond() {
		s.versionPad = in.base + in.pos - len("<?xml version=")
	}
	in.line, in.column = 1, len("<?xml version=")+1
}

// xmlDecl scans the rest of the XML declaration (production [23]) from
// where the pre-scan stopped (s.declStage: 0 before the first pseudo
// attribute, 1 before the '=' of version, 2 before its value).
func (s *saxScanner) xmlDecl() {
	seen := 0
	var encoding string
	haveEncoding, haveStandalone := false, false
	for {
		var which string
		sawSpace := true
		stage := 0
		if seen == 0 {
			stage = s.declStage
		}
		if stage == 0 {
			if seen > 0 {
				sawSpace = s.skipSpaces() || seen == 1 && s.versionLoose
			}
			c := s.peek()
			if c == '?' || c == saxStopEOF {
				if seen == 0 {
					s.fatal("The version is required in the XML declaration.")
				}
				break
			}
			switch {
			case s.skipASCII("version"):
				which = "version"
			case s.skipASCII("encoding"):
				which = "encoding"
			case s.skipASCII("standalone"):
				which = "standalone"
			default:
				s.fatal("A pseudo attribute name is expected.")
			}
			s.skipSpaces()
		} else {
			which = "version"
		}
		if stage < 2 {
			if !s.skipChar('=') {
				s.fatalOrEOF("The ' = ' character must follow \"" + which + "\" in the XML declaration.")
			}
			s.skipSpaces()
		}
		quote := s.peek()
		if quote != '"' && quote != '\'' {
			s.fatalOrEOF("The value following \"" + which + "\" in the XML declaration must be a quoted string.")
		}
		s.in.advance(1)
		if s.in.pos == s.in.count-1 {
			// Observed behaviour (black-box): when the value starts with
			// the last char of the buffer, more is read before it.
			s.in.load(s.in.pos)
		}
		var value []uint16
		padded := seen != 0
		for {
			c := s.peek()
			if c == saxStopEOF {
				s.fatal(s.eofMessage())
			}
			if saxIsHighSurrogate(c) && saxIsLowSurrogate(s.peekSecond()) {
				value = append(value, uint16(c), s.in.buf[s.in.pos+1])
				s.in.advance(2)
			} else if saxIsHighSurrogate(c) {
				s.loneHighSurrogate(c)
			} else {
				if !saxIsXMLChar(c) {
					s.fatal("An invalid XML character (Unicode: " + saxHex(c) + ") was found in the XML declaration.")
				}
				if c = s.next(); c == quote {
					break
				}
				value = append(value, uint16(c))
			}
			if !padded && len(value) >= 4 {
				for i := 0; i < s.versionPad; i++ {
					value = append(value, ' ')
				}
				s.in.column += s.versionPad
				padded = true
			}
		}
		v := jutil.StringFromUTF16(value)
		switch {
		case seen == 0 && which != "version":
			s.fatal("The version is required in the XML declaration.")
		case seen == 0:
			if v != "1.0" {
				s.fatal("XML version \"" + v + "\" is not supported, only XML 1.0 is supported.")
			}
		// Observed behaviour (black-box): the messages for a pseudo
		// attribute out of order.
		case haveStandalone:
			s.fatal("No more pseudo attributes are allowed.")
		case haveEncoding && which != "standalone":
			s.fatal("The standalone name in XML declaration may be misspelled.")
		case which == "version":
			s.fatal("The encoding declaration is required in the text declaration.")
		case !sawSpace:
			// Observed behaviour (black-box): the message names encoding
			// for standalone too.
			s.fatal("White space is required before the encoding pseudo attribute in the XML declaration.")
		case which == "encoding":
			haveEncoding = true
			encoding = v
		default:
			if v != "yes" && v != "no" {
				s.fatal("The standalone document declaration value must be \"yes\" or \"no\", not \"" + v + "\".")
			}
			haveStandalone = true
			s.standalone = v == "yes"
		}
		seen++
	}
	if !s.skipChar('?') {
		s.fatalOrEOF("The XML declaration must end with \"?>\".")
	}
	if !s.skipChar('>') {
		s.fatalOrEOF("The XML declaration must end with \"?>\".")
	}
	s.in.declMode = false
	if haveEncoding {
		if !saxIsEncName(encoding) {
			s.fatal("Invalid encoding name \"" + encoding + "\".")
		}
		enc, ok := saxEncodingByName(encoding)
		if !ok {
			if saxKnownEncoding(encoding) {
				saxUnsupported("encoding \"" + encoding + "\" is")
			}
			// Observed behaviour (black-box): the JDK throws
			// java.io.UnsupportedEncodingException for an encoding it does
			// not know.
			e := &saxUnsupportedEncodingException{}
			jutil.InitThrowable(e, "java.io.UnsupportedEncodingException", &encoding, nil)
			panic(e)
		}
		s.in.setEncoding(enc, encoding)
	}
}

// saxIsEncName is production [81] EncName.
func saxIsEncName(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		letter := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		if i == 0 {
			if !letter {
				return false
			}
		} else if !letter && !(c >= '0' && c <= '9') && c != '.' && c != '_' && c != '-' {
			return false
		}
	}
	return true
}

// prolog scans Misc* (doctypedecl Misc*)? and the root element.
func (s *saxScanner) prolog() {
	for {
		s.skipSpaces()
		c := s.peek()
		if c == saxStopEOF {
			s.fatal("Premature end of file.")
		}
		if c == '&' {
			s.in.advance(1)
			s.fatal("Reference is not allowed in prolog.")
		}
		if c != '<' {
			s.fatal("Content is not allowed in prolog.")
		}
		s.in.advance(1)
		c = s.peek()
		switch {
		case c == '?':
			s.in.advance(1)
			s.processingInstruction(true)
		case c == '!':
			s.in.advance(1)
			if s.peek() == '-' {
				s.in.advance(1)
				if !s.skipChar('-') {
					s.fatalOrEOF("Comment must start with \"<!--\".")
				}
				s.comment()
			} else if s.skipASCII("DOCTYPE") {
				if s.sawDoctype {
					s.fatal("Already seen doctype.")
				}
				s.doctype()
			} else {
				s.fatalOrEOF("The markup in the document preceding the root element must be well-formed.")
			}
		case saxIsNameStart(c):
			s.element()
			return
		default:
			s.fatalOrEOF("The markup in the document preceding the root element must be well-formed.")
		}
	}
}

// trailingMisc scans the Misc* after the root element.
func (s *saxScanner) trailingMisc() {
	for {
		s.skipSpaces()
		c := s.peek()
		if c == saxStopEOF {
			return
		}
		if c != '<' {
			s.fatal("Content is not allowed in trailing section.")
		}
		s.in.advance(1)
		c = s.peek()
		switch c {
		case '?':
			s.in.advance(1)
			s.processingInstruction(true)
		case '!':
			s.in.advance(1)
			if !s.skipASCII("--") {
				s.fatal("Comment must start with \"<!--\".")
			}
			s.comment()
		case '/':
			// Observed behaviour (black-box): reported after "</".
			s.in.advance(1)
			s.fatal("The markup in the document following the root element must be well-formed.")
		default:
			s.fatalOrEOF("The markup in the document following the root element must be well-formed.")
		}
	}
}

// ---------------------------------------------------------------------------
// Elements and content.

// element scans the root element: its start tag (the '<' is consumed and
// a name start is next) and, unless it is empty, its content.
func (s *saxScanner) element() {
	if s.startTag() {
		return
	}
	s.content()
}

// startTag scans a start tag or empty-element tag after '<' and reports
// it; it returns whether the element was empty. Otherwise the element is
// pushed on s.open.
func (s *saxScanner) startTag() bool {
	nameUnits := s.scanName()
	name := jutil.StringFromUTF16(nameUnits)
	s.attrs.names = s.attrs.names[:0]
	s.attrs.values = s.attrs.values[:0]
	empty := false
	for {
		sawSpace := s.skipSpaces()
		c := s.peek()
		if c == '>' {
			s.in.advance(1)
			break
		}
		if c == '/' {
			s.in.advance(1)
			if !s.skipChar('>') {
				s.fatalOrEOF("Element type \"" + name + "\" must be followed by either attribute specifications, \">\" or \"/>\".")
			}
			empty = true
			break
		}
		if c == saxStopEOF {
			s.fatal(s.eofMessage())
		}
		if !sawSpace || !saxIsNameStart(c) {
			s.fatal("Element type \"" + name + "\" must be followed by either attribute specifications, \">\" or \"/>\".")
		}
		attr := jutil.StringFromUTF16(s.scanName())
		s.skipSpaces()
		if !s.skipChar('=') {
			s.fatalOrEOF("Attribute name \"" + attr + "\" associated with an element type \"" + name + "\" must be followed by the ' = ' character.")
		}
		s.skipSpaces()
		quote := s.peek()
		if quote != '"' && quote != '\'' {
			s.fatalOrEOF("Open quote is expected for attribute \"" + attr + "\" associated with an  element type  \"" + name + "\".")
		}
		s.in.advance(1)
		value := s.attributeValue(quote, name, attr)
		for _, n := range s.attrs.names {
			if n == attr {
				s.fatal("Attribute \"" + attr + "\" was already specified for element \"" + name + "\".")
			}
		}
		if len(s.attrs.names) == 10000 {
			s.fatal("JAXP00010002:  Element \"" + name + "\" has more than \"10,000\" attributes, \"10,000\" is the limit imposed by the JDK.")
		}
		s.attrs.names = append(s.attrs.names, attr)
		s.attrs.values = append(s.attrs.values, value)
	}
	if !empty && s.elementContent[name] {
		saxUnsupported("elements declared with element content (their white space is ignorable) are")
	}
	s.handler.StartElement("", "", name, &s.attrs)
	if empty {
		s.handler.EndElement("", "", name)
		return true
	}
	s.open = append(s.open, nameUnits)
	return false
}

// endTag scans an end tag after "</" and reports it.
func (s *saxScanner) endTag() {
	nameUnits := s.open[len(s.open)-1]
	name := jutil.StringFromUTF16(nameUnits)
	if !s.skipUnits(nameUnits) {
		s.fatal("The element type \"" + name + "\" must be terminated by the matching end-tag \"</" + name + ">\".")
	}
	s.skipSpaces()
	if !s.skipChar('>') {
		s.fatalOrEOF("The end-tag for element type \"" + name + "\" must end with a '>' delimiter.")
	}
	s.handler.EndElement("", "", name)
	s.open = s.open[:len(s.open)-1]
}

// content scans the content of the root element up to its end tag.
func (s *saxScanner) content() {
	for len(s.open) > 0 {
		// Observed behaviour (black-box): markup and references are
		// recognized before a text step starts (so without its reading
		// rule for the last char of the buffer).
		c := s.peek()
		if c != '<' && c != '&' && c != saxStopEOF {
			c = s.textStep(true)
		}
		switch c {
		case '<':
			s.in.advance(1)
			switch s.peek() {
			case '/':
				s.in.advance(1)
				s.flushText()
				s.endTag()
				continue
			case saxStopEOF:
				s.fatal(s.eofMessage())
			}
			s.flushText()
			s.markup()
		case '&':
			s.in.advance(1)
			s.flushText()
			s.reference()
		case saxStopEOF:
			s.fatal(s.eofMessage())
		default:
			if c != ']' && c != '\n' && c != '\r' && c != saxStopBufEnd && !saxIsTextChar(uint16(c)) {
				s.otherChar(c)
				s.flushText()
				continue
			}
			if c == ']' {
				s.brackets()
			}
			c = s.textStep(false)
			if c == saxStopEOF {
				s.fatal(s.eofMessage())
			}
			if c == ']' {
				s.brackets()
			}
			s.flushText()
		}
	}
}

// textStep scans character data into s.text (see the file comment for
// where it stops). It returns the char it stopped at (not consumed),
// saxStopBufEnd at the end of the buffer, or saxStopEOF; after line ends
// near the end of the buffer, '\n' unless '<' or a char that is not
// character data follows.
func (s *saxScanner) textStep(first bool) int {
	in := s.in
	if in.pos == in.count {
		if in.load(in.pos) == 0 {
			return saxStopEOF
		}
	} else if in.pos == in.count-1 {
		in.load(in.pos)
	}
	if c := in.buf[in.pos]; c == '\n' || c == '\r' {
		var stopped bool
		s.text, stopped = s.lineEnds(s.text, '\n')
		if stopped {
			// Observed behaviour (black-box): what follows the line ends
			// matters only when it is '<' or, with a char left in the
			// buffer, not character data. In the first step of a
			// characters() call, the JDK reads more input to see the next
			// char when the line ends end the buffer.
			if in.pos < in.count {
				c := in.buf[in.pos]
				if c == '<' || (!saxIsTextChar(c) && c != '&' && c != ']' && c != '\n' && c != '\r') {
					return int(c)
				}
			} else if first && in.load(in.pos) > 0 && in.buf[in.pos] == '<' {
				return '<'
			}
			return '\n'
		}
	}
	start := in.pos
	end := start
	for end < in.count && saxIsTextChar(in.buf[end]) {
		end++
	}
	s.text = append(s.text, in.buf[start:end]...)
	in.advance(end - start)
	if end == in.count {
		// Observed behaviour (black-box): as after line ends, the first
		// step reads on to see whether '<' follows.
		if first && in.load(in.pos) > 0 && in.buf[in.pos] == '<' {
			return '<'
		}
		return saxStopBufEnd
	}
	return int(in.buf[end])
}

// lineEnds consumes the run of line ends at pos (§2.11: CR LF and CR are
// read as LF), appending repl to dst for each. It stops early, returning
// true, after a line end that leaves at most one char in the buffer.
// Observed behaviour (black-box): after the run the JDK's column is 1
// minus the number of lone CRs (not followed by LF) in it, and stays that
// far off for the rest of the line; except when the run stops with one
// char left in the buffer.
func (s *saxScanner) lineEnds(dst []uint16, repl uint16) ([]uint16, bool) {
	in := s.in
	lone := 0
	stopped := false
	for {
		c := in.buf[in.pos]
		if c == '\n' {
			in.newline(1)
		} else if c == '\r' {
			if in.pos+1 < in.count && in.buf[in.pos+1] == '\n' {
				in.newline(2)
			} else {
				// Observed behaviour (black-box): a CR that ends the input
				// does not count.
				if in.pos+1 < in.count {
					lone++
				}
				in.newline(1)
			}
		} else {
			break
		}
		dst = append(dst, repl)
		if in.pos >= in.count-1 {
			stopped = true
			break
		}
	}
	if !stopped || in.pos != in.count-1 {
		in.column -= lone
	}
	return dst, stopped
}

// brackets scans a run of ']' in character data; "]]>" is an error.
func (s *saxScanner) brackets() {
	s.in.advance(1)
	s.text = append(s.text, ']')
	n := 1
	for {
		c := s.peek()
		if c == ']' {
			s.in.advance(1)
			s.text = append(s.text, ']')
			n++
			continue
		}
		if c == '>' && n >= 2 {
			s.in.advance(1)
			s.fatal("The character sequence \"]]>\" must not appear in content unless used to mark the end of a CDATA section.")
		}
		if c == saxStopEOF {
			s.fatal(s.eofMessage())
		}
		return
	}
}

// otherChar handles the char that stopped textStep when it is not markup:
// a surrogate pair, or an invalid character.
func (s *saxScanner) otherChar(c int) {
	if saxIsHighSurrogate(c) {
		low := s.peekSecond()
		if saxIsLowSurrogate(low) {
			s.text = append(s.text, uint16(c), uint16(low))
			s.in.advance(2)
			return
		}
		s.loneHighSurrogate(c)
	}
	s.fatal("An invalid XML character (Unicode: " + saxHex(c) + ") was found in the element content of the document.")
}

// loneHighSurrogate reports a high surrogate that is not followed by a low
// one. Observed behaviour (black-box): in every construct, the JDK reports
// it after consuming it, as found in element content; at the end of the
// input it reports the end of input instead.
func (s *saxScanner) loneHighSurrogate(c int) {
	s.in.advance(1)
	if s.peek() == saxStopEOF {
		s.fatal(s.eofMessage())
	}
	s.fatal("An invalid XML character (Unicode: " + saxHex(c) + ") was found in the element content of the document.")
}

// peekSecond returns the char after the next one, reading more input
// (keeping the next char) when it is not in the buffer.
func (s *saxScanner) peekSecond() int {
	in := s.in
	if in.pos+1 >= in.count {
		in.load(in.pos)
		if in.pos+1 >= in.count {
			return saxStopEOF
		}
	}
	return int(in.buf[in.pos+1])
}

// markup scans the markup after '<' in content (not an end tag).
func (s *saxScanner) markup() {
	c := s.peek()
	switch {
	case c == '!':
		s.in.advance(1)
		switch s.peek() {
		case '-':
			s.in.advance(1)
			if !s.skipChar('-') {
				s.fatalOrEOF("Comment must start with \"<!--\".")
			}
			s.comment()
		case '[':
			if !s.skipASCII("[CDATA[") {
				s.fatalOrEOF("The content of elements must consist of well-formed character data or markup.")
			}
			s.cdata()
		case 'D':
			if s.skipASCII("DOCTYPE") {
				// Observed behaviour (black-box): the JDK fails with this
				// internal-state SAXException (not a SAXParseException).
				e := &SAXException{}
				msg := "Scanner State 24 not Recognized "
				jutil.InitThrowable(e, "org.xml.sax.SAXException", &msg, nil)
				panic(e)
			}
			s.fatalOrEOF("The content of elements must consist of well-formed character data or markup.")
		default:
			s.fatalOrEOF("The content of elements must consist of well-formed character data or markup.")
		}
	case c == '?':
		s.in.advance(1)
		s.processingInstruction(true)
	case saxIsNameStart(c):
		s.startTag()
	default:
		s.fatalOrEOF("The content of elements must consist of well-formed character data or markup.")
	}
}

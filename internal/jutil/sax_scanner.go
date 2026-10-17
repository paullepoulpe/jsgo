package jutil

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
	System_Err().Println("[Fatal Error] :" + strconv.Itoa(line) + ":" + strconv.Itoa(column) + ": " + msg)
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
	panic(NewUnsupportedOperationExceptionMsg("SAX parser: " + what + " not emulated"))
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
		v := StringFromUTF16(value)
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
			InitThrowable(e, "java.io.UnsupportedEncodingException", &encoding, nil)
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
	name := StringFromUTF16(nameUnits)
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
		attr := StringFromUTF16(s.scanName())
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
	name := StringFromUTF16(nameUnits)
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
				InitThrowable(e, "org.xml.sax.SAXException", &msg, nil)
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

// reference scans a reference in content after '&' and reports it.
func (s *saxScanner) reference() {
	if s.peek() == '#' {
		s.in.advance(1)
		s.text = s.charRef(s.text)
		s.flushText()
		return
	}
	nameUnits := s.scanName()
	if nameUnits == nil {
		s.fatalOrEOF("The entity name must immediately follow the '&' in the entity reference.")
	}
	name := StringFromUTF16(nameUnits)
	if !s.skipChar(';') {
		s.fatalOrEOF("The reference to entity \"" + name + "\" must end with the ';' delimiter.")
	}
	if c, ok := saxPredefined(name); ok {
		s.text = append(s.text, c)
		s.flushText()
		return
	}
	decl, declared := s.entities[name]
	switch {
	case !declared:
		if !s.externalSubset || s.standalone {
			s.fatal("The entity \"" + name + "\" was referenced, but not declared.")
		}
		s.handler.SkippedEntity(name)
	case decl.unparsed:
		s.fatal("The unparsed entity reference \"&" + name + ";\" is not permitted.")
	case decl.external:
		// External general entities are not read (the external-general-
		// entities feature is off): SAX reports them as skipped.
		s.handler.SkippedEntity(name)
	default:
		saxUnsupported("references to internal entities are")
	}
}

func saxPredefined(name string) (uint16, bool) {
	switch name {
	case "amp":
		return '&', true
	case "lt":
		return '<', true
	case "gt":
		return '>', true
	case "quot":
		return '"', true
	case "apos":
		return '\'', true
	}
	return 0, false
}

// charRef scans a character reference after "&#" (production [66]) and
// appends the character to out.
func (s *saxScanner) charRef(out []uint16) []uint16 {
	hex := false
	if s.peek() == 'x' {
		s.in.advance(1)
		hex = true
	}
	var digits []uint16
	value := 0
	for {
		c := s.peek()
		var d int
		switch {
		case c >= '0' && c <= '9':
			d = c - '0'
		case hex && c >= 'a' && c <= 'f':
			d = c - 'a' + 10
		case hex && c >= 'A' && c <= 'F':
			d = c - 'A' + 10
		default:
			d = -1
		}
		if d < 0 {
			break
		}
		digits = append(digits, uint16(c))
		if value <= 0x10FFFF {
			if hex {
				value = value*16 + d
			} else {
				value = value*10 + d
			}
		}
		s.in.advance(1)
	}
	if len(digits) == 0 {
		if hex {
			s.fatalOrEOF("A hexadecimal representation must immediately follow the \"&#x\" in a character reference.")
		}
		s.fatalOrEOF("A decimal representation must immediately follow the \"&#\" in a character reference.")
	}
	if !s.skipChar(';') {
		s.fatalOrEOF("The character reference must end with the ';' delimiter.")
	}
	valid := value == 0x9 || value == 0xA || value == 0xD || (value >= 0x20 && value <= 0xD7FF) ||
		(value >= 0xE000 && value <= 0xFFFD) || (value >= 0x10000 && value <= 0x10FFFF)
	if !valid {
		prefix := "&#"
		if hex {
			prefix = "&#x"
		}
		s.fatal("Character reference \"" + prefix + StringFromUTF16(digits) + "\" is an invalid XML character.")
	}
	if value >= 0x10000 {
		v := value - 0x10000
		return append(out, uint16(0xD800+v>>10), uint16(0xDC00+v&0x3FF))
	}
	return append(out, uint16(value))
}

// attributeValue scans an attribute value after its opening quote, with
// the normalization of XML 1.0 §3.3.3 for CDATA attributes: each white
// space char (a line end counting as one) becomes a space, references are
// replaced. It is scanned in the same steps as character data.
func (s *saxScanner) attributeValue(quote int, elem, attr string) string {
	in := s.in
	s.val = s.val[:0]
	for {
		if in.pos == in.count {
			if in.load(in.pos) == 0 {
				s.fatal(s.eofMessage())
			}
		} else if in.pos == in.count-1 {
			in.load(in.pos)
		}
		c := int(in.buf[in.pos])
		if c == '\n' || c == '\r' {
			var stopped bool
			s.val, stopped = s.lineEnds(s.val, ' ')
			if stopped {
				continue
			}
		}
		bracket := false
		for in.pos < in.count {
			c = int(in.buf[in.pos])
			if c == quote {
				in.advance(1)
				return StringFromUTF16(s.val)
			}
			if c == '\t' {
				s.val = append(s.val, ' ')
				in.advance(1)
				continue
			}
			if c == ']' {
				// Observed behaviour (black-box): a ']' ends the step, so
				// the reading rule applies again after it.
				s.val = append(s.val, ']')
				in.advance(1)
				bracket = true
				break
			}
			if !saxIsTextChar(uint16(c)) {
				break
			}
			s.val = append(s.val, uint16(c))
			in.advance(1)
		}
		if bracket || in.pos == in.count {
			continue
		}
		switch {
		case c == '\n' || c == '\r':
		case c == '<':
			s.fatal("The value of attribute \"" + attr + "\" associated with an element type \"" + elem + "\" must not contain the '<' character.")
		case c == '&':
			in.advance(1)
			s.attributeReference(elem, attr)
		case saxIsHighSurrogate(c) && saxIsLowSurrogate(s.peekSecond()):
			s.val = append(s.val, uint16(c), in.buf[in.pos+1])
			in.advance(2)
		case saxIsHighSurrogate(c):
			s.loneHighSurrogate(c)
		default:
			s.fatal("An invalid XML character (Unicode: " + saxHex(c) + ") was found in the value of attribute \"" + attr + "\" and element is \"" + elem + "\".")
		}
	}
}

// attributeReference scans a reference in an attribute value after '&'.
func (s *saxScanner) attributeReference(elem, attr string) {
	if s.peek() == '#' {
		s.in.advance(1)
		s.val = s.charRef(s.val)
		return
	}
	nameUnits := s.scanName()
	if nameUnits == nil {
		s.fatalOrEOF("The entity name must immediately follow the '&' in the entity reference.")
	}
	name := StringFromUTF16(nameUnits)
	if !s.skipChar(';') {
		s.fatalOrEOF("The reference to entity \"" + name + "\" must end with the ';' delimiter.")
	}
	if c, ok := saxPredefined(name); ok {
		s.val = append(s.val, c)
		return
	}
	decl, declared := s.entities[name]
	switch {
	case !declared:
		if !s.externalSubset || s.standalone {
			s.fatal("The entity \"" + name + "\" was referenced, but not declared.")
		}
		// Observed behaviour (black-box): with an unread external subset,
		// an undeclared entity contributes nothing to the value.
	case decl.unparsed:
		s.fatal("The unparsed entity reference \"&" + name + ";\" is not permitted.")
	case decl.external:
		s.fatal("The external entity reference \"&" + name + ";\" is not permitted in an attribute value.")
	default:
		saxUnsupported("references to internal entities are")
	}
}

// cdata scans a CDATA section after "<![CDATA[" and reports its text in
// one characters() call.
func (s *saxScanner) cdata() {
	data := s.scanUntil("]]>", "CDATA section", func(part []uint16) {
		s.handler.Characters(part, 0, len(part))
	})
	if len(data) > 0 {
		s.handler.Characters(data, 0, len(data))
	}
}

// comment scans a comment after "<!--".
func (s *saxScanner) comment() {
	s.scanUntil("--", "comment", nil)
	if !s.skipChar('>') {
		s.fatalOrEOF("The string \"--\" is not permitted within comments.")
	}
}

// processingInstruction scans a processing instruction after "<?";
// report says whether to report it (PIs in the internal subset are not).
func (s *saxScanner) processingInstruction(report bool) {
	targetUnits := s.scanName()
	if targetUnits == nil {
		s.fatalOrEOF("The processing instruction must begin with the name of the target.")
	}
	target := StringFromUTF16(targetUnits)
	if len(targetUnits) == 3 && strings.EqualFold(target, "xml") {
		s.fatal("The processing instruction target matching \"[xX][mM][lL]\" is not allowed.")
	}
	data := ""
	if s.peek() != '?' || !s.skipASCII("?>") {
		if !s.skipSpaces() {
			s.fatalOrEOF("White space is required between the processing instruction target and data.")
		}
		data = StringFromUTF16(s.scanUntil("?>", "processing instruction", nil))
	}
	if report {
		s.handler.ProcessingInstruction(target, data)
	}
}

// scanUntil scans chars up to delim (consumed) into s.val, normalizing
// line ends; what names the construct in the invalid-character message.
// It goes in steps like character data: a step starts with the line ends
// there and stops at a line end, at delim's first char or at the end of
// the buffer. Observed behaviour (black-box): it reads more input,
// keeping the unscanned chars, whenever fewer chars than the delimiter's
// length are left in the buffer at a step start, or when the buffer ends
// inside what could be the delimiter.
func (s *saxScanner) scanUntil(delim, what string, flush func([]uint16)) []uint16 {
	in := s.in
	s.val = s.val[:0]
	d0 := uint16(delim[0])
	for {
		if in.count-in.pos < len(delim) {
			for in.load(in.pos) > 0 && in.count-in.pos < len(delim) {
			}
			if in.count-in.pos < len(delim) {
				// Observed behaviour (black-box): at the end of input the
				// chars left are skipped as plain chars (a line end counts
				// as a column).
				s.val = append(s.val, in.buf[in.pos:in.count]...)
				in.advance(in.count - in.pos)
				s.fatal(s.eofMessage())
			}
		}
		if c := in.buf[in.pos]; c == '\n' || c == '\r' {
			var stopped bool
			s.val, stopped = s.lineEnds(s.val, '\n')
			if stopped {
				// Observed behaviour (black-box): when the line ends
				// leave one char in the buffer that is not a line end, a
				// CDATA section reports what it has so far.
				if flush != nil && in.pos == in.count-1 && len(s.val) > 0 {
					flush(s.val)
					s.val = s.val[:0]
				}
				continue
			}
		}
	chars:
		for in.pos < in.count {
			c := in.buf[in.pos]
			switch {
			case c == d0:
				match := true
				for i := 1; i < len(delim); i++ {
					if in.pos+i == in.count {
						// The buffer ends inside a possible delimiter.
						for in.load(in.pos) > 0 && in.pos+i == in.count {
						}
						if in.pos+i == in.count {
							s.val = append(s.val, in.buf[in.pos:in.count]...)
							in.advance(in.count - in.pos)
							s.fatal(s.eofMessage())
						}
					}
					if in.buf[in.pos+i] != uint16(delim[i]) {
						match = false
						break
					}
				}
				if match {
					in.advance(len(delim))
					return s.val
				}
				s.val = append(s.val, c)
				in.advance(1)
			case c == '\n' || c == '\r':
				break chars
			case saxIsXMLChar(int(c)):
				s.val = append(s.val, c)
				in.advance(1)
			case saxIsHighSurrogate(int(c)) && saxIsLowSurrogate(s.peekSecond()):
				// Observed behaviour (black-box): a supplementary
				// character ends the step.
				s.val = append(s.val, c, in.buf[in.pos+1])
				in.advance(2)
				break chars
			case saxIsHighSurrogate(int(c)):
				s.loneHighSurrogate(int(c))
			default:
				s.fatal("An invalid XML character (Unicode: " + saxHex(int(c)) + ") was found in the " + what + ".")
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Document type declaration (production [28]). The external subset is not
// read (the load-external-dtd feature is off), so the EntityResolver is
// never called; the internal subset is scanned for its entity
// declarations.

func (s *saxScanner) doctype() {
	s.sawDoctype = true
	if !s.skipSpaces() {
		s.fatalOrEOF("White space is required after \"<!DOCTYPE\" in the document type declaration.")
	}
	rootUnits := s.scanName()
	if rootUnits == nil {
		s.fatalOrEOF("The root element type must appear after \"<!DOCTYPE\" in the document type declaration.")
	}
	root := StringFromUTF16(rootUnits)
	if s.skipSpaces() {
		if s.skipASCII("SYSTEM") {
			if !s.skipSpaces() {
				s.fatalOrEOF("White space is required after keyword SYSTEM in DOCTYPE decl.")
			}
			s.systemLiteral()
			s.externalSubset = true
		} else if s.skipASCII("PUBLIC") {
			if !s.skipSpaces() {
				s.fatalOrEOF("White spaces are required after keyword PUBLIC in DOCTYPE decl.")
			}
			s.pubidLiteral()
			if !s.skipSpaces() {
				s.fatalOrEOF("White spaces are required between publicId and systemId.")
			}
			s.systemLiteral()
			s.externalSubset = true
		}
		s.skipSpaces()
	}
	if s.peek() == '[' {
		s.in.advance(1)
		s.inSubset = true
		s.internalSubset()
		s.inSubset = false
		s.skipSpaces()
		if s.peek() == saxStopEOF {
			// Observed behaviour (black-box): reported without a position.
			s.fatalAt("Premature end of file.", -1, -1)
		}
	}
	if !s.skipChar('>') {
		s.fatalOrEOF("The document type declaration for root element type \"" + root + "\" must end with '>'.")
	}
}

func (s *saxScanner) systemLiteral() {
	quote := s.peek()
	if quote != '"' && quote != '\'' {
		s.fatalOrEOF("The system identifier must begin with either a single or double quote character.")
	}
	s.in.advance(1)
	s.literal(quote)
}

func (s *saxScanner) pubidLiteral() {
	quote := s.peek()
	if quote != '"' && quote != '\'' {
		// Observed behaviour (black-box): reported after that char.
		if s.next() == saxStopEOF {
			s.fatal(s.eofMessage())
		}
		s.fatal("The public identifier must begin with either a single or double quote character.")
	}
	s.in.advance(1)
	s.pubidChars(quote)
}

// pubidChars scans a public identifier literal after its opening quote.
// Observed behaviour (black-box): it is read one char at a time (more input
// is read only when the buffer is used up), and the char after a line end
// counts two columns (after CR LF, the column is 2 at once).
func (s *saxScanner) pubidChars(quote int) {
	in := s.in
	bump := false
	for {
		c := s.peek()
		if c == saxStopEOF {
			s.fatal(s.eofMessage())
		}
		if bump {
			in.column++
			bump = false
		}
		switch {
		case c == quote:
			in.advance(1)
			return
		case c == '\n' || c == '\r':
			at := in.base + in.pos
			s.next()
			if in.base+in.pos-at == 2 {
				// Observed behaviour (black-box): after CR LF the column
				// is 2 at once.
				in.column = 2
			} else {
				bump = true
			}
		case !saxIsPubidChar(c):
			in.advance(1)
			s.fatal("An invalid XML character (Unicode: " + saxHex(c) + ") was found in the public identifier.")
		default:
			in.advance(1)
		}
	}
}

// literal scans a system identifier literal after its opening quote.
// Observed behaviour (black-box): it is scanned in steps like character
// data, with the same reading rule for the last char of the buffer and the
// same line-end runs, and each '&', '<', ']' or '%' ends a step; it takes
// the chars of production [2] Char in the BMP only.
func (s *saxScanner) literal(quote int) {
	in := s.in
	for {
		if in.pos == in.count {
			if in.load(in.pos) == 0 {
				s.fatal(s.eofMessage())
			}
		} else if in.pos == in.count-1 {
			in.load(in.pos)
		}
		if c := in.buf[in.pos]; (c == '\n' || c == '\r') && int(c) != quote {
			if _, stopped := s.lineEnds(nil, '\n'); stopped {
				if c := in.buf[in.pos]; in.pos == in.count-1 && (c == '&' || c == '<' || c == ']' || c == '%') {
					// Observed behaviour (black-box): such a char left
					// after the run is consumed before more is read.
					in.advance(1)
				}
				continue
			}
		}
		for in.pos < in.count {
			c := int(in.buf[in.pos])
			if c == quote {
				in.advance(1)
				return
			}
			if c == '\n' || c == '\r' {
				break
			}
			if !saxIsXMLChar(c) {
				s.fatal("An invalid XML character (Unicode: " + saxHex(c) + ") was found in the system identifier.")
			}
			in.advance(1)
			if c == '&' || c == '<' || c == ']' || c == '%' {
				break
			}
		}
	}
}

// saxIsPubidChar is production [13] PubidChar.
func saxIsPubidChar(c int) bool {
	switch {
	case c == ' ' || c == '\n' || c == '\r':
		return true
	case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		return true
	}
	return strings.ContainsRune("-'()+,./:=?;!*#@$_%", rune(c))
}

func (s *saxScanner) internalSubset() {
	const notWellFormed = "The markup declarations contained or pointed to by the document type declaration must be well-formed."
	for {
		s.skipSpaces()
		c := s.peek()
		switch c {
		case saxStopEOF:
			s.fatalAt("Premature end of file.", -1, -1)
		case ']':
			s.in.advance(1)
			return
		case '%':
			s.in.advance(1)
			name := s.peReference(false)
			if _, ok := s.paramEntities[name]; ok {
				saxUnsupported("parameter entity references are")
			}
			// Observed behaviour (black-box): a reference to an undeclared
			// parameter entity is ignored.
		case '<':
			s.in.advance(1)
			switch s.peek() {
			case '?':
				s.in.advance(1)
				s.processingInstruction(false)
			case '!':
				s.in.advance(1)
				switch {
				case s.peek() == '-':
					s.in.advance(1)
					if !s.skipChar('-') {
						s.fatalOrEOF(notWellFormed)
					}
					s.comment()
				case s.skipASCII("ELEMENT"):
					s.elementDecl()
				case s.skipASCII("ATTLIST"):
					s.attlistDecl()
				case s.skipASCII("ENTITY"):
					s.entityDecl()
				case s.skipASCII("NOTATION"):
					s.notationDecl()
				case s.peek() == '[':
					// Observed behaviour (black-box): a conditional section
					// is reported after its '['.
					s.in.advance(1)
					s.fatal(notWellFormed)
				default:
					s.fatalOrEOF(notWellFormed)
				}
			default:
				s.fatalOrEOF(notWellFormed)
			}
		default:
			s.fatal(notWellFormed)
		}
	}
}

// elementDecl scans an element type declaration (production [45]).
func (s *saxScanner) elementDecl() {
	if !s.skipSpaces() {
		s.fatalOrEOF("White space is required after \"<!ELEMENT\" in the element type declaration.")
	}
	nameUnits := s.scanName()
	if nameUnits == nil {
		s.fatalOrEOF("The element type is required in the element type declaration.")
	}
	name := StringFromUTF16(nameUnits)
	if !s.skipSpaces() {
		s.fatalOrEOF("White space is required after the element type \"" + name + "\" in the element type declaration.")
	}
	switch {
	case s.skipASCII("EMPTY"), s.skipASCII("ANY"):
	case s.peek() == '(':
		s.in.advance(1)
		s.skipSpaces()
		if !s.skipASCII("#PCDATA") {
			s.childrenGroup(name)
			if s.elementContent == nil {
				s.elementContent = map[string]bool{}
			}
			s.elementContent[name] = true
			break
		}
		s.skipSpaces()
		names := 0
		for s.peek() == '|' {
			s.in.advance(1)
			s.skipSpaces()
			if s.scanName() == nil {
				s.fatalOrEOF("An element type is required in the declaration of element type \"" + name + "\".")
			}
			names++
			s.skipSpaces()
		}
		if c := s.peek(); c != ')' && names == 0 {
			s.fatalOrEOF("A ')' is required in the declaration of element type \"" + name + "\".")
		}
		if names > 0 {
			if !s.skipASCII(")*") {
				s.fatalOrEOF("The mixed content model \"" + name + "\" must end with \")*\" when the types of child elements are constrained.")
			}
		} else if !s.skipASCII(")*") && !s.skipChar(')') {
			// Observed behaviour (black-box): ")*" is looked for first,
			// which reads more input before the ')' is consumed.
			s.fatalOrEOF("The mixed content model \"" + name + "\" must end with \")*\" when the types of child elements are constrained.")
		}
	default:
		s.fatalOrEOF("A '(' character or an element type is required in the declaration of element type \"" + name + "\".")
	}
	s.skipSpaces()
	if !s.skipChar('>') {
		s.fatalOrEOF("The declaration for element type \"" + name + "\" must end with '>'.")
	}
}

// childrenGroup scans the rest of a children content model group after
// its '(' and white space (productions [47] to [50]).
func (s *saxScanner) childrenGroup(name string) {
	s.contentParticle(name)
	s.skipSpaces()
	separator := 0
	for {
		c := s.peek()
		if c == ')' {
			s.in.advance(1)
			break
		}
		if (c != '|' && c != ',') || (separator != 0 && c != separator) {
			s.fatalOrEOF("A ')' is required in the declaration of element type \"" + name + "\".")
		}
		separator = c
		s.in.advance(1)
		s.skipSpaces()
		s.contentParticle(name)
		s.skipSpaces()
	}
	s.occurrence()
}

func (s *saxScanner) contentParticle(name string) {
	if s.peek() == '(' {
		s.in.advance(1)
		s.skipSpaces()
		s.childrenGroup(name)
		return
	}
	if s.scanName() == nil {
		s.fatalOrEOF("A '(' character or an element type is required in the declaration of element type \"" + name + "\".")
	}
	s.occurrence()
}

func (s *saxScanner) occurrence() {
	if c := s.peek(); c == '?' || c == '*' || c == '+' {
		s.in.advance(1)
	}
}

// attlistDecl scans an attribute-list declaration (production [52]); only
// one without attribute definitions is emulated.
func (s *saxScanner) attlistDecl() {
	if !s.skipSpaces() {
		s.fatalOrEOF("White space is required after \"<!ATTLIST\" in the attribute-list declaration.")
	}
	if s.scanName() == nil {
		s.fatalOrEOF("The element type is required in the attribute-list declaration.")
	}
	s.skipSpaces()
	if !s.skipChar('>') {
		saxUnsupported("attribute-list declarations are")
	}
}

// entityDecl scans an entity declaration (production [70]).
func (s *saxScanner) entityDecl() {
	if !s.skipSpaces() {
		s.fatalOrEOF("White space is required after \"<!ENTITY\" in the entity declaration.")
	}
	isParam := false
	if s.peek() == '%' {
		s.in.advance(1)
		if !s.skipSpaces() {
			s.fatalOrEOF("White space is required after \"<!ENTITY\" in the entity declaration.")
		}
		isParam = true
	}
	nameUnits := s.scanName()
	if nameUnits == nil {
		s.fatalOrEOF("The name of the entity is required in the entity declaration.")
	}
	name := StringFromUTF16(nameUnits)
	if !s.skipSpaces() {
		s.fatalOrEOF("White space is required between the entity name \"" + name + "\" and the definition in the entity declaration.")
	}
	decl := saxEntityDecl{}
	// Observed behaviour (black-box): the keywords are tried first (which
	// reads ahead), then the quote.
	switch {
	case s.skipASCII("SYSTEM"):
		if !s.skipSpaces() {
			s.fatalOrEOF("White space is required after keyword SYSTEM in DOCTYPE decl.")
		}
		s.systemLiteral()
		decl.external = true
	case s.skipASCII("PUBLIC"):
		if !s.skipSpaces() {
			s.fatalOrEOF("White spaces are required after keyword PUBLIC in DOCTYPE decl.")
		}
		s.pubidLiteral()
		if !s.skipSpaces() {
			s.fatalOrEOF("White spaces are required between publicId and systemId.")
		}
		s.systemLiteral()
		decl.external = true
	case s.peek() == '"' || s.peek() == '\'':
		quote := s.peek()
		s.in.advance(1)
		s.entityValue(quote)
	default:
		// Observed behaviour (black-box): the JDK reports the bare message
		// key, after the char that is not a quote.
		if s.next() == saxStopEOF {
			s.fatal(s.eofMessage())
		}
		s.fatal("OpenQuoteMissingInDecl")
	}
	s.skipSpaces()
	// Observed behaviour (black-box): NDATA is looked for (reading ahead)
	// whether or not white space precedes.
	if decl.external && !isParam && s.skipASCII("NDATA") {
		if !s.skipSpaces() {
			s.fatalOrEOF("White space is required between \"NDATA\" and the notation name in the declaration for the entity \"" + name + "\".")
		}
		if s.scanName() == nil {
			s.fatalOrEOF("The name of the notation is required after \"NDATA\" in the declaration for the entity \"" + name + "\".")
		}
		decl.unparsed = true
		s.skipSpaces()
	}
	if !s.skipChar('>') {
		s.fatalOrEOF("The declaration for the entity \"" + name + "\" must end with '>'.")
	}
	if isParam {
		if s.paramEntities == nil {
			s.paramEntities = map[string]saxEntityDecl{}
		}
		if _, ok := s.paramEntities[name]; !ok {
			s.paramEntities[name] = decl
		}
		return
	}
	if _, predefined := saxPredefined(name); predefined {
		return
	}
	if s.entities == nil {
		s.entities = map[string]saxEntityDecl{}
	}
	if _, ok := s.entities[name]; !ok {
		s.entities[name] = decl
	}
}

// entityValue scans an entity value literal after its opening quote
// (production [9]), checking its references. Observed behaviour
// (black-box): it is scanned in steps like character data (a step ends
// after a reference and at the end of the buffer, with the same reading
// rule for the last char of the buffer and the same line-end runs at the
// start of a step), except that a line end met inside a step (or left
// over by a run that stopped at the end of the buffer) is consumed on its
// own (like a line end left as the last char of the buffer: it is
// consumed before more input is read), ends the step, and makes the next
// char count two columns. A '<' or ']' also ends the step.
func (s *saxScanner) entityValue(quote int) {
	in := s.in
	continued := false
	for {
		if in.pos == in.count {
			if in.load(in.pos) == 0 {
				s.fatal("Premature end of file.")
			}
		} else if c := in.buf[in.pos]; in.pos == in.count-1 && !(continued && (c == '\n' || c == '\r')) {
			in.load(in.pos)
		}
		if c := in.buf[in.pos]; (c == '\n' || c == '\r') && !continued {
			if _, stopped := s.lineEnds(nil, '\n'); stopped {
				continued = true
				if c := in.buf[in.pos]; in.pos == in.count-1 && c != '\n' && c != '\r' {
					// Observed behaviour (black-box): when the run leaves
					// one other char, the column is one too far while
					// more input is read (a decoding error shows it).
					in.column++
					in.load(in.pos)
					in.column--
				}
				continue
			}
		}
		continued = false
	step:
		for in.pos < in.count {
			c := int(in.buf[in.pos])
			switch {
			case c == quote:
				in.advance(1)
				return
			case c == '\n' || c == '\r':
				s.next()
				in.column++
				break step
			case c == '<' || c == ']':
				in.advance(1)
				break step
			case c == '%':
				in.advance(1)
				s.peReference(true)
				break step
			case c == '&':
				in.advance(1)
				if s.peek() == '#' {
					in.advance(1)
					s.charRef(nil)
					break step
				}
				nameUnits := s.scanName()
				if nameUnits == nil {
					s.fatalOrEOF("The entity name must immediately follow the '&' in the entity reference.")
				}
				if !s.skipChar(';') {
					s.fatalOrEOF("The reference to entity \"" + StringFromUTF16(nameUnits) + "\" must end with the ';' delimiter.")
				}
				break step
			case saxIsHighSurrogate(c) && saxIsLowSurrogate(s.peekSecond()):
				in.advance(2)
			case saxIsHighSurrogate(c):
				s.loneHighSurrogate(c)
			case saxIsXMLChar(c):
				in.advance(1)
			default:
				// Observed behaviour (black-box): the JDK reports the bare
				// message key.
				s.fatal("InvalidCharInLiteral")
			}
		}
	}
}

// peReference scans a parameter entity reference after its '%' and
// returns its name. In an entity value (inValue) any reference is an
// error.
func (s *saxScanner) peReference(inValue bool) string {
	nameUnits := s.scanName()
	if nameUnits == nil {
		s.fatalOrEOF("The entity name must immediately follow the '%' in the parameter entity reference.")
	}
	name := StringFromUTF16(nameUnits)
	if !s.skipChar(';') {
		s.fatalOrEOF("The parameter entity reference \"%" + name + ";\" must end with the ';' delimiter.")
	}
	if inValue {
		s.fatal("The parameter entity reference \"%" + name + ";\" cannot occur within markup in the internal subset of the DTD.")
	}
	return name
}

// notationDecl scans a notation declaration (production [82]).
func (s *saxScanner) notationDecl() {
	if !s.skipSpaces() {
		s.fatalOrEOF("White space is required after \"<!NOTATION\" in the notation declaration.")
	}
	nameUnits := s.scanName()
	if nameUnits == nil {
		s.fatalOrEOF("The name of the notation is required in the notation declaration.")
	}
	name := StringFromUTF16(nameUnits)
	if !s.skipSpaces() {
		s.fatalOrEOF("White space is required after the notation name \"" + name + "\" in the notation declaration.")
	}
	switch {
	case s.skipASCII("SYSTEM"):
		if !s.skipSpaces() {
			s.fatalOrEOF("White space is required after keyword SYSTEM in DOCTYPE decl.")
		}
		s.systemLiteral()
	case s.skipASCII("PUBLIC"):
		if !s.skipSpaces() {
			s.fatalOrEOF("White spaces are required after keyword PUBLIC in DOCTYPE decl.")
		}
		s.pubidLiteral()
		if s.skipSpaces() {
			if q := s.peek(); q == '"' || q == '\'' {
				s.systemLiteral()
			}
		}
	default:
		s.fatalOrEOF("The declaration for the notation \"" + name + "\" must include a system or public identifier.")
	}
	s.skipSpaces()
	if !s.skipChar('>') {
		s.fatalOrEOF("The declaration for the notation \"" + name + "\" must end with '>'.")
	}
}

// saxNameStartRanges and saxNameCharRanges are the inclusive ranges of the
// BMP chars that may start a Name and that may occur in one. Observed
// behaviour (black-box): measured by parsing "<X/>" and "<aX/>" with the
// JDK for every BMP char X; they are the Letter, '_' and ':' and the
// NameChar classes of XML 1.0 (Fourth Edition), Appendix B. No
// supplementary character is accepted.
var saxNameStartRanges = [...]uint16{
	0x003A, 0x003A, 0x0041, 0x005A, 0x005F, 0x005F, 0x0061, 0x007A, 0x00C0, 0x00D6,
	0x00D8, 0x00F6, 0x00F8, 0x0131, 0x0134, 0x013E, 0x0141, 0x0148, 0x014A, 0x017E,
	0x0180, 0x01C3, 0x01CD, 0x01F0, 0x01F4, 0x01F5, 0x01FA, 0x0217, 0x0250, 0x02A8,
	0x02BB, 0x02C1, 0x0386, 0x0386, 0x0388, 0x038A, 0x038C, 0x038C, 0x038E, 0x03A1,
	0x03A3, 0x03CE, 0x03D0, 0x03D6, 0x03DA, 0x03DA, 0x03DC, 0x03DC, 0x03DE, 0x03DE,
	0x03E0, 0x03E0, 0x03E2, 0x03F3, 0x0401, 0x040C, 0x040E, 0x044F, 0x0451, 0x045C,
	0x045E, 0x0481, 0x0490, 0x04C4, 0x04C7, 0x04C8, 0x04CB, 0x04CC, 0x04D0, 0x04EB,
	0x04EE, 0x04F5, 0x04F8, 0x04F9, 0x0531, 0x0556, 0x0559, 0x0559, 0x0561, 0x0586,
	0x05D0, 0x05EA, 0x05F0, 0x05F2, 0x0621, 0x063A, 0x0641, 0x064A, 0x0671, 0x06B7,
	0x06BA, 0x06BE, 0x06C0, 0x06CE, 0x06D0, 0x06D3, 0x06D5, 0x06D5, 0x06E5, 0x06E6,
	0x0905, 0x0939, 0x093D, 0x093D, 0x0958, 0x0961, 0x0985, 0x098C, 0x098F, 0x0990,
	0x0993, 0x09A8, 0x09AA, 0x09B0, 0x09B2, 0x09B2, 0x09B6, 0x09B9, 0x09DC, 0x09DD,
	0x09DF, 0x09E1, 0x09F0, 0x09F1, 0x0A05, 0x0A0A, 0x0A0F, 0x0A10, 0x0A13, 0x0A28,
	0x0A2A, 0x0A30, 0x0A32, 0x0A33, 0x0A35, 0x0A36, 0x0A38, 0x0A39, 0x0A59, 0x0A5C,
	0x0A5E, 0x0A5E, 0x0A72, 0x0A74, 0x0A85, 0x0A8B, 0x0A8D, 0x0A8D, 0x0A8F, 0x0A91,
	0x0A93, 0x0AA8, 0x0AAA, 0x0AB0, 0x0AB2, 0x0AB3, 0x0AB5, 0x0AB9, 0x0ABD, 0x0ABD,
	0x0AE0, 0x0AE0, 0x0B05, 0x0B0C, 0x0B0F, 0x0B10, 0x0B13, 0x0B28, 0x0B2A, 0x0B30,
	0x0B32, 0x0B33, 0x0B36, 0x0B39, 0x0B3D, 0x0B3D, 0x0B5C, 0x0B5D, 0x0B5F, 0x0B61,
	0x0B85, 0x0B8A, 0x0B8E, 0x0B90, 0x0B92, 0x0B95, 0x0B99, 0x0B9A, 0x0B9C, 0x0B9C,
	0x0B9E, 0x0B9F, 0x0BA3, 0x0BA4, 0x0BA8, 0x0BAA, 0x0BAE, 0x0BB5, 0x0BB7, 0x0BB9,
	0x0C05, 0x0C0C, 0x0C0E, 0x0C10, 0x0C12, 0x0C28, 0x0C2A, 0x0C33, 0x0C35, 0x0C39,
	0x0C60, 0x0C61, 0x0C85, 0x0C8C, 0x0C8E, 0x0C90, 0x0C92, 0x0CA8, 0x0CAA, 0x0CB3,
	0x0CB5, 0x0CB9, 0x0CDE, 0x0CDE, 0x0CE0, 0x0CE1, 0x0D05, 0x0D0C, 0x0D0E, 0x0D10,
	0x0D12, 0x0D28, 0x0D2A, 0x0D39, 0x0D60, 0x0D61, 0x0E01, 0x0E2E, 0x0E30, 0x0E30,
	0x0E32, 0x0E33, 0x0E40, 0x0E45, 0x0E81, 0x0E82, 0x0E84, 0x0E84, 0x0E87, 0x0E88,
	0x0E8A, 0x0E8A, 0x0E8D, 0x0E8D, 0x0E94, 0x0E97, 0x0E99, 0x0E9F, 0x0EA1, 0x0EA3,
	0x0EA5, 0x0EA5, 0x0EA7, 0x0EA7, 0x0EAA, 0x0EAB, 0x0EAD, 0x0EAE, 0x0EB0, 0x0EB0,
	0x0EB2, 0x0EB3, 0x0EBD, 0x0EBD, 0x0EC0, 0x0EC4, 0x0F40, 0x0F47, 0x0F49, 0x0F69,
	0x10A0, 0x10C5, 0x10D0, 0x10F6, 0x1100, 0x1100, 0x1102, 0x1103, 0x1105, 0x1107,
	0x1109, 0x1109, 0x110B, 0x110C, 0x110E, 0x1112, 0x113C, 0x113C, 0x113E, 0x113E,
	0x1140, 0x1140, 0x114C, 0x114C, 0x114E, 0x114E, 0x1150, 0x1150, 0x1154, 0x1155,
	0x1159, 0x1159, 0x115F, 0x1161, 0x1163, 0x1163, 0x1165, 0x1165, 0x1167, 0x1167,
	0x1169, 0x1169, 0x116D, 0x116E, 0x1172, 0x1173, 0x1175, 0x1175, 0x119E, 0x119E,
	0x11A8, 0x11A8, 0x11AB, 0x11AB, 0x11AE, 0x11AF, 0x11B7, 0x11B8, 0x11BA, 0x11BA,
	0x11BC, 0x11C2, 0x11EB, 0x11EB, 0x11F0, 0x11F0, 0x11F9, 0x11F9, 0x1E00, 0x1E9B,
	0x1EA0, 0x1EF9, 0x1F00, 0x1F15, 0x1F18, 0x1F1D, 0x1F20, 0x1F45, 0x1F48, 0x1F4D,
	0x1F50, 0x1F57, 0x1F59, 0x1F59, 0x1F5B, 0x1F5B, 0x1F5D, 0x1F5D, 0x1F5F, 0x1F7D,
	0x1F80, 0x1FB4, 0x1FB6, 0x1FBC, 0x1FBE, 0x1FBE, 0x1FC2, 0x1FC4, 0x1FC6, 0x1FCC,
	0x1FD0, 0x1FD3, 0x1FD6, 0x1FDB, 0x1FE0, 0x1FEC, 0x1FF2, 0x1FF4, 0x1FF6, 0x1FFC,
	0x2126, 0x2126, 0x212A, 0x212B, 0x212E, 0x212E, 0x2180, 0x2182, 0x3007, 0x3007,
	0x3021, 0x3029, 0x3041, 0x3094, 0x30A1, 0x30FA, 0x3105, 0x312C, 0x4E00, 0x9FA5,
	0xAC00, 0xD7A3,
}

var saxNameCharRanges = [...]uint16{
	0x002D, 0x002E, 0x0030, 0x003A, 0x0041, 0x005A, 0x005F, 0x005F, 0x0061, 0x007A,
	0x00B7, 0x00B7, 0x00C0, 0x00D6, 0x00D8, 0x00F6, 0x00F8, 0x0131, 0x0134, 0x013E,
	0x0141, 0x0148, 0x014A, 0x017E, 0x0180, 0x01C3, 0x01CD, 0x01F0, 0x01F4, 0x01F5,
	0x01FA, 0x0217, 0x0250, 0x02A8, 0x02BB, 0x02C1, 0x02D0, 0x02D1, 0x0300, 0x0345,
	0x0360, 0x0361, 0x0386, 0x038A, 0x038C, 0x038C, 0x038E, 0x03A1, 0x03A3, 0x03CE,
	0x03D0, 0x03D6, 0x03DA, 0x03DA, 0x03DC, 0x03DC, 0x03DE, 0x03DE, 0x03E0, 0x03E0,
	0x03E2, 0x03F3, 0x0401, 0x040C, 0x040E, 0x044F, 0x0451, 0x045C, 0x045E, 0x0481,
	0x0483, 0x0486, 0x0490, 0x04C4, 0x04C7, 0x04C8, 0x04CB, 0x04CC, 0x04D0, 0x04EB,
	0x04EE, 0x04F5, 0x04F8, 0x04F9, 0x0531, 0x0556, 0x0559, 0x0559, 0x0561, 0x0586,
	0x0591, 0x05A1, 0x05A3, 0x05B9, 0x05BB, 0x05BD, 0x05BF, 0x05BF, 0x05C1, 0x05C2,
	0x05C4, 0x05C4, 0x05D0, 0x05EA, 0x05F0, 0x05F2, 0x0621, 0x063A, 0x0640, 0x0652,
	0x0660, 0x0669, 0x0670, 0x06B7, 0x06BA, 0x06BE, 0x06C0, 0x06CE, 0x06D0, 0x06D3,
	0x06D5, 0x06E8, 0x06EA, 0x06ED, 0x06F0, 0x06F9, 0x0901, 0x0903, 0x0905, 0x0939,
	0x093C, 0x094D, 0x0951, 0x0954, 0x0958, 0x0963, 0x0966, 0x096F, 0x0981, 0x0983,
	0x0985, 0x098C, 0x098F, 0x0990, 0x0993, 0x09A8, 0x09AA, 0x09B0, 0x09B2, 0x09B2,
	0x09B6, 0x09B9, 0x09BC, 0x09BC, 0x09BE, 0x09C4, 0x09C7, 0x09C8, 0x09CB, 0x09CD,
	0x09D7, 0x09D7, 0x09DC, 0x09DD, 0x09DF, 0x09E3, 0x09E6, 0x09F1, 0x0A02, 0x0A02,
	0x0A05, 0x0A0A, 0x0A0F, 0x0A10, 0x0A13, 0x0A28, 0x0A2A, 0x0A30, 0x0A32, 0x0A33,
	0x0A35, 0x0A36, 0x0A38, 0x0A39, 0x0A3C, 0x0A3C, 0x0A3E, 0x0A42, 0x0A47, 0x0A48,
	0x0A4B, 0x0A4D, 0x0A59, 0x0A5C, 0x0A5E, 0x0A5E, 0x0A66, 0x0A74, 0x0A81, 0x0A83,
	0x0A85, 0x0A8B, 0x0A8D, 0x0A8D, 0x0A8F, 0x0A91, 0x0A93, 0x0AA8, 0x0AAA, 0x0AB0,
	0x0AB2, 0x0AB3, 0x0AB5, 0x0AB9, 0x0ABC, 0x0AC5, 0x0AC7, 0x0AC9, 0x0ACB, 0x0ACD,
	0x0AE0, 0x0AE0, 0x0AE6, 0x0AEF, 0x0B01, 0x0B03, 0x0B05, 0x0B0C, 0x0B0F, 0x0B10,
	0x0B13, 0x0B28, 0x0B2A, 0x0B30, 0x0B32, 0x0B33, 0x0B36, 0x0B39, 0x0B3C, 0x0B43,
	0x0B47, 0x0B48, 0x0B4B, 0x0B4D, 0x0B56, 0x0B57, 0x0B5C, 0x0B5D, 0x0B5F, 0x0B61,
	0x0B66, 0x0B6F, 0x0B82, 0x0B83, 0x0B85, 0x0B8A, 0x0B8E, 0x0B90, 0x0B92, 0x0B95,
	0x0B99, 0x0B9A, 0x0B9C, 0x0B9C, 0x0B9E, 0x0B9F, 0x0BA3, 0x0BA4, 0x0BA8, 0x0BAA,
	0x0BAE, 0x0BB5, 0x0BB7, 0x0BB9, 0x0BBE, 0x0BC2, 0x0BC6, 0x0BC8, 0x0BCA, 0x0BCD,
	0x0BD7, 0x0BD7, 0x0BE7, 0x0BEF, 0x0C01, 0x0C03, 0x0C05, 0x0C0C, 0x0C0E, 0x0C10,
	0x0C12, 0x0C28, 0x0C2A, 0x0C33, 0x0C35, 0x0C39, 0x0C3E, 0x0C44, 0x0C46, 0x0C48,
	0x0C4A, 0x0C4D, 0x0C55, 0x0C56, 0x0C60, 0x0C61, 0x0C66, 0x0C6F, 0x0C82, 0x0C83,
	0x0C85, 0x0C8C, 0x0C8E, 0x0C90, 0x0C92, 0x0CA8, 0x0CAA, 0x0CB3, 0x0CB5, 0x0CB9,
	0x0CBE, 0x0CC4, 0x0CC6, 0x0CC8, 0x0CCA, 0x0CCD, 0x0CD5, 0x0CD6, 0x0CDE, 0x0CDE,
	0x0CE0, 0x0CE1, 0x0CE6, 0x0CEF, 0x0D02, 0x0D03, 0x0D05, 0x0D0C, 0x0D0E, 0x0D10,
	0x0D12, 0x0D28, 0x0D2A, 0x0D39, 0x0D3E, 0x0D43, 0x0D46, 0x0D48, 0x0D4A, 0x0D4D,
	0x0D57, 0x0D57, 0x0D60, 0x0D61, 0x0D66, 0x0D6F, 0x0E01, 0x0E2E, 0x0E30, 0x0E3A,
	0x0E40, 0x0E4E, 0x0E50, 0x0E59, 0x0E81, 0x0E82, 0x0E84, 0x0E84, 0x0E87, 0x0E88,
	0x0E8A, 0x0E8A, 0x0E8D, 0x0E8D, 0x0E94, 0x0E97, 0x0E99, 0x0E9F, 0x0EA1, 0x0EA3,
	0x0EA5, 0x0EA5, 0x0EA7, 0x0EA7, 0x0EAA, 0x0EAB, 0x0EAD, 0x0EAE, 0x0EB0, 0x0EB9,
	0x0EBB, 0x0EBD, 0x0EC0, 0x0EC4, 0x0EC6, 0x0EC6, 0x0EC8, 0x0ECD, 0x0ED0, 0x0ED9,
	0x0F18, 0x0F19, 0x0F20, 0x0F29, 0x0F35, 0x0F35, 0x0F37, 0x0F37, 0x0F39, 0x0F39,
	0x0F3E, 0x0F47, 0x0F49, 0x0F69, 0x0F71, 0x0F84, 0x0F86, 0x0F8B, 0x0F90, 0x0F95,
	0x0F97, 0x0F97, 0x0F99, 0x0FAD, 0x0FB1, 0x0FB7, 0x0FB9, 0x0FB9, 0x10A0, 0x10C5,
	0x10D0, 0x10F6, 0x1100, 0x1100, 0x1102, 0x1103, 0x1105, 0x1107, 0x1109, 0x1109,
	0x110B, 0x110C, 0x110E, 0x1112, 0x113C, 0x113C, 0x113E, 0x113E, 0x1140, 0x1140,
	0x114C, 0x114C, 0x114E, 0x114E, 0x1150, 0x1150, 0x1154, 0x1155, 0x1159, 0x1159,
	0x115F, 0x1161, 0x1163, 0x1163, 0x1165, 0x1165, 0x1167, 0x1167, 0x1169, 0x1169,
	0x116D, 0x116E, 0x1172, 0x1173, 0x1175, 0x1175, 0x119E, 0x119E, 0x11A8, 0x11A8,
	0x11AB, 0x11AB, 0x11AE, 0x11AF, 0x11B7, 0x11B8, 0x11BA, 0x11BA, 0x11BC, 0x11C2,
	0x11EB, 0x11EB, 0x11F0, 0x11F0, 0x11F9, 0x11F9, 0x1E00, 0x1E9B, 0x1EA0, 0x1EF9,
	0x1F00, 0x1F15, 0x1F18, 0x1F1D, 0x1F20, 0x1F45, 0x1F48, 0x1F4D, 0x1F50, 0x1F57,
	0x1F59, 0x1F59, 0x1F5B, 0x1F5B, 0x1F5D, 0x1F5D, 0x1F5F, 0x1F7D, 0x1F80, 0x1FB4,
	0x1FB6, 0x1FBC, 0x1FBE, 0x1FBE, 0x1FC2, 0x1FC4, 0x1FC6, 0x1FCC, 0x1FD0, 0x1FD3,
	0x1FD6, 0x1FDB, 0x1FE0, 0x1FEC, 0x1FF2, 0x1FF4, 0x1FF6, 0x1FFC, 0x20D0, 0x20DC,
	0x20E1, 0x20E1, 0x2126, 0x2126, 0x212A, 0x212B, 0x212E, 0x212E, 0x2180, 0x2182,
	0x3005, 0x3005, 0x3007, 0x3007, 0x3021, 0x302F, 0x3031, 0x3035, 0x3041, 0x3094,
	0x3099, 0x309A, 0x309D, 0x309E, 0x30A1, 0x30FA, 0x30FC, 0x30FE, 0x3105, 0x312C,
	0x4E00, 0x9FA5, 0xAC00, 0xD7A3,
}

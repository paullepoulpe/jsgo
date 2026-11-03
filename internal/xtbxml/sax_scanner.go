package xtbxml

// Go-only: language-adapter — the XML scanner behind XMLReader.Parse
// (sax_parser.go): a non-validating, namespace-unaware XML 1.0 processor
// for documents without external entities, written from the XML 1.0
// specification:
//   - §2.8 prolog and document type declaration, §4.3.3 encoding
//     declaration, §2.9 standalone declaration;
//   - §3.1 start-tags, end-tags and empty-element tags, §3.3 attribute-list
//     declarations (default values and non-CDATA normalization) and §3.2
//     element-content declarations (ignorable whitespace);
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
// encodings other than UTF-8, UTF-16, ISO-8859-1 and US-ASCII, and
// references to internal entities (general or parameter).

import (
	"io"
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

	open            [][]uint16 // names of the open elements
	openElemContent []bool     // whether each open element has element-only content
	attrs           saxAttributes

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
	attlists       map[string][]saxAttrDecl
	ended          bool
	err            jutil.Throwable
}

func newSaxScanner(handler ContentHandler, src io.Reader) *saxScanner {
	if handler == nil {
		handler = DefaultHandler{}
	}
	s := &saxScanner{handler: handler}
	s.in = newSaxInput(s, src)
	s.locator.s = s
	return s
}

// ---------------------------------------------------------------------------
// Errors.

func (s *saxScanner) setErr(err jutil.Throwable) {
	if s.err == nil {
		s.err = err
	}
}

func (s *saxScanner) fatal(msg string) {
	s.fatalAt(msg, s.in.line, s.in.column)
}

// fatalAt reports the SAXParseException of a fatal error. With no
// ErrorHandler, the JDK also prints the error on System.err; observed
// behaviour (black-box): "[Fatal Error] :L:C: msg".
func (s *saxScanner) fatalAt(msg string, line, column int) {
	if s.err != nil {
		return
	}
	jutil.System_Err().Println("[Fatal Error] :" + strconv.Itoa(line) + ":" + strconv.Itoa(column) + ": " + msg)
	s.err = newSAXParseExceptionAt(msg, line, column)
}

// fatalOrEOF reports msg, or the end-of-input error when the input has
// ended.
func (s *saxScanner) fatalOrEOF(msg string) {
	if s.peek() == saxStopEOF {
		s.fatal(s.eofMessage())
		return
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
	if s.err != nil {
		return saxStopEOF
	}
	in := s.in
	if (in.pos == in.count && in.load(in.pos) == 0) || s.err != nil {
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
	if s.err != nil {
		return false
	}
	in := s.in
	for in.count-in.pos < len(str) {
		if in.loadMore() == 0 || s.err != nil {
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
	if s.err != nil {
		return false
	}
	in := s.in
	for in.count-in.pos < len(str) {
		if in.loadMore() == 0 || s.err != nil {
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
	if s.err != nil {
		return nil
	}
	in := s.in
	if (in.pos == in.count && in.load(in.pos) == 0) || s.err != nil {
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
			return nil
		}
		read := in.load(start)
		if s.err != nil {
			return nil
		}
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
		return nil
	}
	return name
}

// flushText reports the pending character data, routing whitespace-only runs
// inside element-content elements to ContentHandler.IgnorableWhitespace.
func (s *saxScanner) flushText() {
	if s.err == nil && len(s.text) > 0 {
		if len(s.openElemContent) > 0 && s.openElemContent[len(s.openElemContent)-1] && saxAllSpaces(s.text) {
			s.handler.IgnorableWhitespace(s.text, 0, len(s.text))
		} else {
			s.handler.Characters(s.text, 0, len(s.text))
		}
		s.text = s.text[:0]
	}
}

func saxAllSpaces(units []uint16) bool {
	for _, u := range units {
		if !saxIsSpace(int(u)) {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Document.

func (s *saxScanner) parse() jutil.Throwable {
	s.in.detectEncoding()
	if s.err != nil {
		return s.err
	}
	s.preScanDeclaration()
	if s.err != nil {
		return s.err
	}
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
	if s.err == nil {
		s.prolog()
	}
	if s.err == nil {
		s.trailingMisc()
	}
	if s.err != nil {
		return s.err
	}
	s.ended = true
	s.handler.EndDocument()
	return nil
}

// prolog scans Misc* (doctypedecl Misc*)? and the root element.
func (s *saxScanner) prolog() {
	for s.err == nil {
		s.skipSpaces()
		c := s.peek()
		if c == saxStopEOF {
			s.fatal("Premature end of file.")
			return
		}
		if c == '&' {
			s.in.advance(1)
			s.fatal("Reference is not allowed in prolog.")
			return
		}
		if c != '<' {
			s.fatal("Content is not allowed in prolog.")
			return
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
					return
				}
				s.comment()
			} else if s.skipASCII("DOCTYPE") {
				if s.sawDoctype {
					s.fatal("Already seen doctype.")
					return
				}
				s.doctype()
			} else {
				s.fatalOrEOF("The markup in the document preceding the root element must be well-formed.")
				return
			}
		case saxIsNameStart(c):
			s.element()
			return
		default:
			s.fatalOrEOF("The markup in the document preceding the root element must be well-formed.")
			return
		}
	}
}

// trailingMisc scans the Misc* after the root element.
func (s *saxScanner) trailingMisc() {
	for s.err == nil {
		s.skipSpaces()
		c := s.peek()
		if c == saxStopEOF {
			return
		}
		if c != '<' {
			s.fatal("Content is not allowed in trailing section.")
			return
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
				return
			}
			s.comment()
		case '/':
			// Observed behaviour (black-box): reported after "</".
			s.in.advance(1)
			s.fatal("The markup in the document following the root element must be well-formed.")
			return
		default:
			s.fatalOrEOF("The markup in the document following the root element must be well-formed.")
			return
		}
	}
}

// ---------------------------------------------------------------------------
// Elements and content.

// element scans the root element: its start tag (the '<' is consumed and
// a name start is next) and, unless it is empty, its content.
func (s *saxScanner) element() {
	if s.startTag() || s.err != nil {
		return
	}
	s.content()
}

// startTag scans a start tag or empty-element tag after '<' and reports
// it; it returns whether the element was empty. Otherwise the element is
// pushed on s.open.
func (s *saxScanner) startTag() bool {
	nameUnits := s.scanName()
	if s.err != nil {
		return false
	}
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
				return false
			}
			empty = true
			break
		}
		if c == saxStopEOF {
			s.fatal(s.eofMessage())
			return false
		}
		if !sawSpace || !saxIsNameStart(c) {
			s.fatal("Element type \"" + name + "\" must be followed by either attribute specifications, \">\" or \"/>\".")
			return false
		}
		attrUnits := s.scanName()
		if s.err != nil {
			return false
		}
		attr := jutil.StringFromUTF16(attrUnits)
		s.skipSpaces()
		if !s.skipChar('=') {
			s.fatalOrEOF("Attribute name \"" + attr + "\" associated with an element type \"" + name + "\" must be followed by the ' = ' character.")
			return false
		}
		s.skipSpaces()
		quote := s.peek()
		if quote != '"' && quote != '\'' {
			s.fatalOrEOF("Open quote is expected for attribute \"" + attr + "\" associated with an  element type  \"" + name + "\".")
			return false
		}
		s.in.advance(1)
		value := s.attributeValue(quote, name, attr)
		if s.err != nil {
			return false
		}
		if s.isNonCDATAAttr(name, attr) {
			value = saxNormalizeNonCDATA(value)
		}
		for _, n := range s.attrs.names {
			if n == attr {
				s.fatal("Attribute \"" + attr + "\" was already specified for element \"" + name + "\".")
				return false
			}
		}
		if len(s.attrs.names) == 10000 {
			s.fatal("JAXP00010002:  Element \"" + name + "\" has more than \"10,000\" attributes, \"10,000\" is the limit imposed by the JDK.")
			return false
		}
		s.attrs.names = append(s.attrs.names, attr)
		s.attrs.values = append(s.attrs.values, value)
	}
	if decls := s.attlists[name]; len(decls) > 0 {
		for _, d := range decls {
			if _, ok := s.attrs.Value(d.name); d.hasDefault && !ok {
				if len(s.attrs.names) == 10000 {
					s.fatal("JAXP00010002:  Element \"" + name + "\" has more than \"10,000\" attributes, \"10,000\" is the limit imposed by the JDK.")
					return false
				}
				s.attrs.names = append(s.attrs.names, d.name)
				s.attrs.values = append(s.attrs.values, d.defaultVal)
			}
		}
	}
	s.handler.StartElement(name, &s.attrs)
	if empty {
		s.handler.EndElement(name)
		return true
	}
	s.open = append(s.open, nameUnits)
	s.openElemContent = append(s.openElemContent, s.elementContent[name])
	return false
}

// endTag scans an end tag after "</" and reports it.
func (s *saxScanner) endTag() {
	nameUnits := s.open[len(s.open)-1]
	name := jutil.StringFromUTF16(nameUnits)
	if !s.skipUnits(nameUnits) {
		s.fatal("The element type \"" + name + "\" must be terminated by the matching end-tag \"</" + name + ">\".")
		return
	}
	s.skipSpaces()
	if !s.skipChar('>') {
		s.fatalOrEOF("The end-tag for element type \"" + name + "\" must end with a '>' delimiter.")
		return
	}
	s.handler.EndElement(name)
	s.open = s.open[:len(s.open)-1]
	s.openElemContent = s.openElemContent[:len(s.openElemContent)-1]
}

// content scans the content of the root element up to its end tag.
func (s *saxScanner) content() {
	for len(s.open) > 0 && s.err == nil {
		// Observed behaviour (black-box): markup and references are
		// recognized before a text step starts (so without its reading
		// rule for the last char of the buffer).
		c := s.peek()
		if c != '<' && c != '&' && c != saxStopEOF {
			c = s.textStep(true)
		}
		if s.err != nil {
			return
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
				return
			}
			s.flushText()
			s.markup()
		case '&':
			s.in.advance(1)
			s.flushText()
			s.reference()
		case saxStopEOF:
			s.fatal(s.eofMessage())
			return
		default:
			if c != ']' && c != '\n' && c != '\r' && c != saxStopBufEnd && !saxIsTextChar(uint16(c)) {
				s.otherChar(c)
				s.flushText()
				continue
			}
			if c == ']' {
				if s.brackets(); s.err != nil {
					return
				}
			}
			c = s.textStep(false)
			if s.err != nil {
				return
			}
			if c == saxStopEOF {
				s.fatal(s.eofMessage())
				return
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
		if s.err != nil {
			return saxStopEOF
		}
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
			return
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
		return
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
		return
	}
	s.fatal("An invalid XML character (Unicode: " + saxHex(c) + ") was found in the element content of the document.")
}

// peekSecond returns the char after the next one, reading more input
// (keeping the next char) when it is not in the buffer.
func (s *saxScanner) peekSecond() int {
	in := s.in
	if in.pos+1 >= in.count {
		in.load(in.pos)
		if in.pos+1 >= in.count || s.err != nil {
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
				return
			}
			s.comment()
		case '[':
			if !s.skipASCII("[CDATA[") {
				s.fatalOrEOF("The content of elements must consist of well-formed character data or markup.")
				return
			}
			s.cdata()
		case 'D':
			if s.skipASCII("DOCTYPE") {
				// Observed behaviour (black-box): the JDK fails with this
				// internal-state SAXException (not a SAXParseException).
				e := &SAXException{}
				jutil.InitThrowableMsg(e, "org.xml.sax.SAXException", "Scanner State 24 not Recognized ", nil)
				s.setErr(e)
				return
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

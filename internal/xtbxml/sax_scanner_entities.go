package xtbxml

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

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
	name := jutil.StringFromUTF16(nameUnits)
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
		s.fatal("Character reference \"" + prefix + jutil.StringFromUTF16(digits) + "\" is an invalid XML character.")
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
				return jutil.StringFromUTF16(s.val)
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
	name := jutil.StringFromUTF16(nameUnits)
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
	target := jutil.StringFromUTF16(targetUnits)
	if len(targetUnits) == 3 && strings.EqualFold(target, "xml") {
		s.fatal("The processing instruction target matching \"[xX][mM][lL]\" is not allowed.")
	}
	data := ""
	if s.peek() != '?' || !s.skipASCII("?>") {
		if !s.skipSpaces() {
			s.fatalOrEOF("White space is required between the processing instruction target and data.")
		}
		data = jutil.StringFromUTF16(s.scanUntil("?>", "processing instruction", nil))
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
	root := jutil.StringFromUTF16(rootUnits)
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
	name := jutil.StringFromUTF16(nameUnits)
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
	name := jutil.StringFromUTF16(nameUnits)
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
					s.fatalOrEOF("The reference to entity \"" + jutil.StringFromUTF16(nameUnits) + "\" must end with the ';' delimiter.")
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
	name := jutil.StringFromUTF16(nameUnits)
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
	name := jutil.StringFromUTF16(nameUnits)
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

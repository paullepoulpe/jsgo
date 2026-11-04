package xtbxml

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type attrDecl struct {
	name       string
	nonCDATA   bool
	hasDefault bool
	defaultVal string
}

// attlistDecl scans an attribute-list declaration (XML 1.0 §3.3, production [52]),
// recording the first declaration of each (element, attribute) pair for default
// value injection and non-CDATA whitespace normalization.
func (s *scanner) attlistDecl() {
	if !s.skipSpaces() {
		s.fatalOrEOF("White space is required after \"<!ATTLIST\" in the attribute-list declaration.")
		return
	}
	elemUnits := s.scanName()
	if elemUnits == nil {
		s.fatalOrEOF("The element type is required in the attribute-list declaration.")
		return
	}
	elem := jutil.StringFromUTF16(elemUnits)
	for attrIdx := 0; s.err == nil; attrIdx++ {
		sawSpace := s.skipSpaces()
		if s.skipChar('>') {
			return
		}
		if attrIdx == 0 && !sawSpace {
			s.fatalOrEOF("White space is required before the attribute name in the attribute-list declaration for element \"" + elem + "\".")
			return
		}
		attrUnits := s.scanName()
		if attrUnits == nil {
			s.fatalOrEOF("The attribute name must be specified in the attribute-list declaration for element \"" + elem + "\".")
			return
		}
		attr := jutil.StringFromUTF16(attrUnits)
		if !s.skipSpaces() {
			s.fatalOrEOF("White space is required before the attribute type in the declaration of attribute \"" + attr + "\" for element \"" + elem + "\".")
			return
		}
		nonCDATA := s.scanAttType(elem, attr)
		if s.err != nil {
			return
		}
		if !s.skipSpaces() {
			s.fatalOrEOF("White space is required before the attribute default in the declaration of attribute \"" + attr + "\" for element \"" + elem + "\".")
			return
		}
		hasDefault, defaultVal := s.scanDefaultDecl(elem, attr, nonCDATA)
		if s.err != nil {
			return
		}
		s.recordAttrDecl(elem, attrDecl{
			name:       attr,
			nonCDATA:   nonCDATA,
			hasDefault: hasDefault,
			defaultVal: defaultVal,
		})
	}
}

func (s *scanner) scanAttType(elem, attr string) (nonCDATA bool) {
	if s.skipChar('(') {
		for {
			s.skipSpaces()
			if !s.scanNmtoken() {
				s.fatalOrEOF("The name token is required in the enumerated type list for the \"" + attr + "\" attribute declaration.")
				return false
			}
			s.skipSpaces()
			if s.skipChar('|') {
				continue
			}
			if !s.skipChar(')') {
				s.fatalOrEOF("The enumerated type list must end with ')' in the \"" + attr + "\" attribute declaration.")
				return false
			}
			return true
		}
	}
	typeUnits := s.scanName()
	if typeUnits == nil {
		s.fatalOrEOF("The attribute type is required in the declaration of attribute \"" + attr + "\" for element \"" + elem + "\".")
		return false
	}
	switch jutil.StringFromUTF16(typeUnits) {
	case "CDATA":
		return false
	case "ID", "IDREF", "IDREFS", "ENTITY", "ENTITIES", "NMTOKEN", "NMTOKENS":
		return true
	case "NOTATION":
		if !s.skipSpaces() {
			s.fatalOrEOF("White space must appear after \"NOTATION\" in the \"" + attr + "\" attribute declaration.")
			return false
		}
		if !s.skipChar('(') {
			s.fatalOrEOF("The '(' character must follow \"NOTATION\" in the \"" + attr + "\" attribute declaration.")
			return false
		}
		for {
			s.skipSpaces()
			if s.scanName() == nil {
				s.fatalOrEOF("The notation name is required in the notation type list for the \"" + attr + "\" attribute declaration.")
				return false
			}
			s.skipSpaces()
			if s.skipChar('|') {
				continue
			}
			if !s.skipChar(')') {
				s.fatalOrEOF("The notation type list must end with ')' in the \"" + attr + "\" attribute declaration.")
				return false
			}
			return true
		}
	default:
		s.fatal("The attribute type is required in the declaration of attribute \"" + attr + "\" for element \"" + elem + "\".")
		return false
	}
}

func (s *scanner) scanDefaultDecl(elem, attr string, nonCDATA bool) (bool, string) {
	switch {
	case s.skipASCII("#REQUIRED"), s.skipASCII("#IMPLIED"):
		return false, ""
	case s.skipASCII("#FIXED"):
		if !s.skipSpaces() {
			s.fatalOrEOF("White space must appear after \"FIXED\" in the \"" + attr + "\" attribute declaration.")
			return false, ""
		}
	}
	quote := s.peek()
	if quote != '"' && quote != '\'' {
		s.fatalOrEOF("Open quote is expected for attribute \"" + attr + "\" associated with an  element type  \"" + elem + "\".")
		return false, ""
	}
	s.in.advance(1)
	val := s.attributeValue(quote, elem, attr)
	if s.err != nil {
		return false, ""
	}
	if nonCDATA {
		val = normalizeNonCDATA(val)
	}
	return true, val
}

func (s *scanner) scanNmtoken() bool {
	in := s.in
	if (in.pos == in.count && in.load(in.pos) == 0) || s.err != nil {
		return false
	}
	if !isNameChar(int(in.buf[in.pos])) {
		return false
	}
	n := 0
	for {
		start := in.pos
		for in.pos < in.count && isNameChar(int(in.buf[in.pos])) {
			in.advance(1)
		}
		n += in.pos - start
		if n > maxNameLength {
			s.fatal("JAXP00010005: The length of entity \"[xml]\" is \"" + formatGroupedInt(n) +
				"\" that exceeds the \"1,000\" limit set by \"jaxp.properties\".")
			return false
		}
		if in.pos < in.count || in.load(in.pos) == 0 || s.err != nil {
			break
		}
	}
	return s.err == nil
}

func (s *scanner) recordAttrDecl(elem string, decl attrDecl) {
	if s.attlists == nil {
		s.attlists = make(map[string][]attrDecl)
	}
	for _, existing := range s.attlists[elem] {
		if existing.name == decl.name {
			return
		}
	}
	s.attlists[elem] = append(s.attlists[elem], decl)
}

func (s *scanner) isNonCDATAAttr(elem, attr string) bool {
	for _, d := range s.attlists[elem] {
		if d.name == attr {
			return d.nonCDATA
		}
	}
	return false
}

// normalizeNonCDATA applies the additional XML 1.0 §3.3.3 normalization step
// for attributes whose declared type is not CDATA: discard leading and trailing
// space (#x20) characters, and replace sequences of spaces with a single space.
func normalizeNonCDATA(v string) string {
	start := 0
	for start < len(v) && v[start] == ' ' {
		start++
	}
	end := len(v)
	for end > start && v[end-1] == ' ' {
		end--
	}
	hasRun := false
	for i := start; i+1 < end; i++ {
		if v[i] == ' ' && v[i+1] == ' ' {
			hasRun = true
			break
		}
	}
	if !hasRun {
		return v[start:end]
	}
	var b strings.Builder
	b.Grow(end - start - 1)
	prevSpace := false
	for i := start; i < end; i++ {
		c := v[i]
		if c == ' ' {
			if prevSpace {
				continue
			}
			prevSpace = true
		} else {
			prevSpace = false
		}
		b.WriteByte(c)
	}
	return b.String()
}

// attributeValue scans an attribute value after its opening quote, with
// the normalization of XML 1.0 §3.3.3 for CDATA attributes: each white
// space char (a line end counting as one) becomes a space, references are
// replaced. It is scanned in the same steps as character data.
func (s *scanner) attributeValue(quote int, elem, attr string) string {
	in := s.in
	s.val = s.val[:0]
	for s.err == nil {
		if in.pos == in.count {
			if in.load(in.pos) == 0 {
				s.fatal(s.eofMessage())
				return ""
			}
		} else if in.pos == in.count-1 {
			in.load(in.pos)
			if s.err != nil {
				return ""
			}
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
			if !isTextChar(uint16(c)) {
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
			return ""
		case c == '&':
			in.advance(1)
			s.attributeReference(elem, attr)
		case isHighSurrogate(c) && isLowSurrogate(s.peekSecond()):
			s.val = append(s.val, uint16(c), in.buf[in.pos+1])
			in.advance(2)
		case isHighSurrogate(c):
			s.loneHighSurrogate(c)
			return ""
		default:
			s.fatal("An invalid XML character (Unicode: " + formatHex(c) + ") was found in the value of attribute \"" + attr + "\" and element is \"" + elem + "\".")
			return ""
		}
	}
	return ""
}

// attributeReference scans a reference in an attribute value after '&'.
func (s *scanner) attributeReference(elem, attr string) {
	if s.peek() == '#' {
		s.in.advance(1)
		s.val = s.charRef(s.val)
		return
	}
	nameUnits := s.scanName()
	if nameUnits == nil {
		s.fatalOrEOF("The entity name must immediately follow the '&' in the entity reference.")
		return
	}
	name := jutil.StringFromUTF16(nameUnits)
	if !s.skipChar(';') {
		s.fatalOrEOF("The reference to entity \"" + name + "\" must end with the ';' delimiter.")
		return
	}
	if c, ok := predefinedEntity(name); ok {
		s.val = append(s.val, c)
		return
	}
	decl, declared := s.entities[name]
	switch {
	case !declared:
		if !s.externalSubset || s.standalone {
			s.fatal("The entity \"" + name + "\" was referenced, but not declared.")
			return
		}
		// Observed behaviour (black-box): with an unread external subset,
		// an undeclared entity contributes nothing to the value.
	case decl.unparsed:
		s.fatal("The unparsed entity reference \"&" + name + ";\" is not permitted.")
	case decl.external:
		s.fatal("The external entity reference \"&" + name + ";\" is not permitted in an attribute value.")
	default:
		panicUnsupported("references to internal entities are")
	}
}

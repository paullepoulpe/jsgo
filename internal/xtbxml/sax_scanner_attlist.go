package xtbxml

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type saxAttrDecl struct {
	name       string
	nonCDATA   bool
	hasDefault bool
	defaultVal string
}

// attlistDecl scans an attribute-list declaration (XML 1.0 §3.3, production [52]),
// recording the first declaration of each (element, attribute) pair for default
// value injection and non-CDATA whitespace normalization.
func (s *saxScanner) attlistDecl() {
	if !s.skipSpaces() {
		s.fatalOrEOF("White space is required after \"<!ATTLIST\" in the attribute-list declaration.")
	}
	elemUnits := s.scanName()
	if elemUnits == nil {
		s.fatalOrEOF("The element type is required in the attribute-list declaration.")
	}
	elem := jutil.StringFromUTF16(elemUnits)
	for attrIdx := 0; ; attrIdx++ {
		sawSpace := s.skipSpaces()
		if s.skipChar('>') {
			return
		}
		if attrIdx == 0 && !sawSpace {
			s.fatalOrEOF("White space is required before the attribute name in the attribute-list declaration for element \"" + elem + "\".")
		}
		attrUnits := s.scanName()
		if attrUnits == nil {
			s.fatalOrEOF("The attribute name must be specified in the attribute-list declaration for element \"" + elem + "\".")
		}
		attr := jutil.StringFromUTF16(attrUnits)
		if !s.skipSpaces() {
			s.fatalOrEOF("White space is required before the attribute type in the declaration of attribute \"" + attr + "\" for element \"" + elem + "\".")
		}
		nonCDATA := s.scanAttType(elem, attr)
		if !s.skipSpaces() {
			s.fatalOrEOF("White space is required before the attribute default in the declaration of attribute \"" + attr + "\" for element \"" + elem + "\".")
		}
		hasDefault, defaultVal := s.scanDefaultDecl(elem, attr, nonCDATA)
		s.recordAttrDecl(elem, saxAttrDecl{
			name:       attr,
			nonCDATA:   nonCDATA,
			hasDefault: hasDefault,
			defaultVal: defaultVal,
		})
	}
}

func (s *saxScanner) scanAttType(elem, attr string) (nonCDATA bool) {
	if s.skipChar('(') {
		for {
			s.skipSpaces()
			if !s.scanNmtoken() {
				s.fatalOrEOF("The name token is required in the enumerated type list for the \"" + attr + "\" attribute declaration.")
			}
			s.skipSpaces()
			if s.skipChar('|') {
				continue
			}
			if !s.skipChar(')') {
				s.fatalOrEOF("The enumerated type list must end with ')' in the \"" + attr + "\" attribute declaration.")
			}
			return true
		}
	}
	typeUnits := s.scanName()
	if typeUnits == nil {
		s.fatalOrEOF("The attribute type is required in the declaration of attribute \"" + attr + "\" for element \"" + elem + "\".")
	}
	switch jutil.StringFromUTF16(typeUnits) {
	case "CDATA":
		return false
	case "ID", "IDREF", "IDREFS", "ENTITY", "ENTITIES", "NMTOKEN", "NMTOKENS":
		return true
	case "NOTATION":
		if !s.skipSpaces() {
			s.fatalOrEOF("White space must appear after \"NOTATION\" in the \"" + attr + "\" attribute declaration.")
		}
		if !s.skipChar('(') {
			s.fatalOrEOF("The '(' character must follow \"NOTATION\" in the \"" + attr + "\" attribute declaration.")
		}
		for {
			s.skipSpaces()
			if s.scanName() == nil {
				s.fatalOrEOF("The notation name is required in the notation type list for the \"" + attr + "\" attribute declaration.")
			}
			s.skipSpaces()
			if s.skipChar('|') {
				continue
			}
			if !s.skipChar(')') {
				s.fatalOrEOF("The notation type list must end with ')' in the \"" + attr + "\" attribute declaration.")
			}
			return true
		}
	default:
		s.fatal("The attribute type is required in the declaration of attribute \"" + attr + "\" for element \"" + elem + "\".")
		return false
	}
}

func (s *saxScanner) scanDefaultDecl(elem, attr string, nonCDATA bool) (bool, string) {
	switch {
	case s.skipASCII("#REQUIRED"), s.skipASCII("#IMPLIED"):
		return false, ""
	case s.skipASCII("#FIXED"):
		if !s.skipSpaces() {
			s.fatalOrEOF("White space must appear after \"FIXED\" in the \"" + attr + "\" attribute declaration.")
		}
	}
	quote := s.peek()
	if quote != '"' && quote != '\'' {
		s.fatalOrEOF("Open quote is expected for attribute \"" + attr + "\" associated with an  element type  \"" + elem + "\".")
	}
	s.in.advance(1)
	val := s.attributeValue(quote, elem, attr)
	if nonCDATA {
		val = saxNormalizeNonCDATA(val)
	}
	return true, val
}

func (s *saxScanner) scanNmtoken() bool {
	in := s.in
	if in.pos == in.count && in.load(in.pos) == 0 {
		return false
	}
	if !saxIsNameChar(int(in.buf[in.pos])) {
		return false
	}
	n := 0
	for {
		start := in.pos
		for in.pos < in.count && saxIsNameChar(int(in.buf[in.pos])) {
			in.advance(1)
		}
		n += in.pos - start
		if n > saxMaxNameLength {
			s.fatal("JAXP00010005: The length of entity \"[xml]\" is \"" + saxGrouped(n) +
				"\" that exceeds the \"1,000\" limit set by \"jaxp.properties\".")
		}
		if in.pos < in.count || in.load(in.pos) == 0 {
			break
		}
	}
	return true
}

func (s *saxScanner) recordAttrDecl(elem string, decl saxAttrDecl) {
	if s.attlists == nil {
		s.attlists = make(map[string][]saxAttrDecl)
	}
	for _, existing := range s.attlists[elem] {
		if existing.name == decl.name {
			return
		}
	}
	s.attlists[elem] = append(s.attlists[elem], decl)
}

func (s *saxScanner) isNonCDATAAttr(elem, attr string) bool {
	for _, d := range s.attlists[elem] {
		if d.name == attr {
			return d.nonCDATA
		}
	}
	return false
}

// saxNormalizeNonCDATA applies the additional XML 1.0 §3.3.3 normalization step
// for attributes whose declared type is not CDATA: discard leading and trailing
// space (#x20) characters, and replace sequences of spaces with a single space.
func saxNormalizeNonCDATA(v string) string {
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

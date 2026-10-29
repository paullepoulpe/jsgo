package jutil

import (
	"strconv"
	"strings"
)

// Go-only: language-adapter — java.net.URI as the URI(String) constructor
// builds it: a URI reference parsed by the grammar of RFC 2396 (Appendix A)
// with the deviations listed in the Javadoc of java.net.URI (an empty
// authority before a path, query or fragment; empty relative paths; IPv4
// hosts checked numerically; one-label host names that start with a digit;
// IPv6 literals in brackets (RFC 2732/2373); and "other" characters —
// non-ASCII, neither ISO control nor space characters — wherever escaped
// octets are allowed).
//
// Positions are UTF-16 indices, as in URISyntaxException.getIndex().
//
// The authority is server-based (user info, host, port) when it parses as
// such, otherwise registry-based (see the Javadoc of
// URI.parseServerAuthority). Equality follows the Javadoc of URI.equals:
// schemes and hosts compare without regard to case, and the hexadecimal
// digits of escaped octets without regard to case.
//
// Observed behaviour (black-box, JDK 21). What the Javadoc leaves open was
// determined by running new URI(s) on many inputs:
//   - the scheme is the text before the first ':' when no '/', '?' or '#'
//     comes before it. An empty one is "Expected scheme name"; a bad
//     character in it is "Illegal character in scheme name". Nothing (or
//     only a fragment) after the ':' is "Expected scheme-specific part";
//   - "//" with nothing at all after it is "Expected authority";
//   - each component is checked from left to right; the first bad character
//     gives "Illegal character in <opaque part|authority|path|query|
//     fragment>", and a '%' without two hex digits "Malformed escape pair",
//     at its index;
//   - the authority is first checked against the registry-name characters.
//     If they all fit, a server-based parse is tried and a failure means a
//     registry-based authority. Otherwise, when the rest of the authority
//     holds only registry characters, '[', ']' and '%', the server-based
//     parse must succeed and its error is thrown; else the error is
//     "Illegal character in authority" at the first non-registry character;
//   - server-based parse: user info runs to the first '@' ("Illegal
//     character in user info"). A host in brackets runs to the first ']'
//     ("Expected closing bracket for IPv6 address" at the end of the
//     authority, or just after '[' for "[]"), with an optional scope id
//     after the first '%' (letters, digits, '_' and '.'; "Illegal character
//     in scope id", or "scope id expected" without an index when empty);
//     only ':' may follow it ("Expected port number"). Other hosts are a
//     dotted quad of decimal numbers up to 255 (any number of digits) when
//     followed by the end or ':', else a host name: labels of letters,
//     digits and inner hyphens separated by dots, an optional final dot,
//     and a last label that starts with a letter unless it is the only one.
//     A host name that stops before the end or ':' is "Illegal character in
//     hostname" at that point (or "Expected hostname" when empty and
//     followed by ':'); a bad last label is the same error at its start;
//   - the port is decimal digits ("Illegal character in port number"); an
//     empty port is undefined and one above 2^31-1 is "Malformed port
//     number" at its first digit;
//   - IPv6 literals (indices relative to the URI; s is the first character
//     after '['): groups of up to four hex digits ("IPv6 hexadecimal digit
//     sequence too long" at the group), separated by ':', with at most one
//     "::" and an optional final dotted quad. A ':' not followed by hex
//     digits is "Expected digits for an IPv6 address"; a group followed by
//     '.' starts the dotted quad. The quad must be the rest of the literal:
//     if its digits-and-dots run is empty or does not reach the end, the
//     error is "Expected IPv4 address" (after a plain ':') or "Expected hex
//     digits or IPv4 address" (after "::"), at the start of the quad;
//     otherwise a bad quad is "Malformed IPv4 address" where it goes wrong
//     (at the start of a number above 255). A literal that starts with
//     neither hex digits nor "::", or has text left over, is "Malformed
//     IPv6 address" at s; more than 16 bytes is "IPv6 address too long",
//     fewer without "::" "IPv6 address too short", and 16 bytes with "::"
//     "Malformed IPv6 address", all at s;
//   - hashCode: starting from 0, each defined component in turn is folded
//     in as h = 127*h + hash(c), where hash is String.hashCode after the hex
//     digits of escaped octets are upper-cased: the scheme (lower-cased,
//     folded as 31*h + c per char), the fragment, then the scheme-specific
//     part of an opaque URI, or the path and query of a hierarchical one
//     followed, for a server-based authority, by the user info, the host
//     (lower-cased, folded as 31*h + c per char) and 1949*port added, and
//     for a registry-based one by the authority.

// URI is java.net.URI.
type URI struct {
	str       string
	path      *string // nil for opaque URIs
	hasScheme bool
}

// ParseURI parses str as a java.net.URI without panicking.
func ParseURI(str string) (*URI, error) {
	p := &uriParser{input: str, s: StringToUTF16(str)}
	u := &URI{str: str}
	if err := p.parse(u); err != nil {
		return nil, err
	}
	return u, nil
}

// IsAbsolute is URI.isAbsolute(): whether the URI has a scheme.
func (u *URI) IsAbsolute() bool { return u.hasScheme }

// GetPath is URI.getPath(): the decoded path component, nil if the path is
// undefined.
func (u *URI) GetPath() *string {
	if u.path == nil {
		return nil
	}
	s := uriDecodeEscapedOctets(*u.path)
	return &s
}

// uriDecodeEscapedOctets replaces each maximal run of escaped octets
// ("%XX%YY...") in s by the UTF-8 decoding of those octets. A URI that
// parsed only contains well-formed escapes.
func uriDecodeEscapedOctets(s string) string {
	if !strings.Contains(s, "%") {
		return s
	}
	units := StringToUTF16(s)
	out := make([]uint16, 0, len(units))
	for i := 0; i < len(units); {
		if units[i] != '%' {
			out = append(out, units[i])
			i++
			continue
		}
		var octets []byte
		for i+2 < len(units) && units[i] == '%' {
			octets = append(octets, byte(uriGetPathHexValue(units[i+1])<<4|uriGetPathHexValue(units[i+2])))
			i += 3
		}
		out = append(out, StringToUTF16(StandardCharsets_UTF_8.Decode(octets))...)
	}
	return StringFromUTF16(out)
}

// uriGetPathHexValue is the value of the hexadecimal digit c.
func uriGetPathHexValue(c uint16) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	default:
		return int(c-'A') + 10
	}
}

// String is URI.toString(): the string the URI was parsed from.
func (u *URI) String() string { return u.str }

// ---------------------------------------------------------------------------
// Character classes (Javadoc of java.net.URI, "Categories of characters").

func uriIsAlpha(c uint16) bool { return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' }

func uriIsDigit(c uint16) bool { return c >= '0' && c <= '9' }

func uriIsAlphanum(c uint16) bool { return uriIsAlpha(c) || uriIsDigit(c) }

func uriIsHex(c uint16) bool {
	return uriIsDigit(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

func uriIn(c uint16, set string) bool {
	for i := 0; i < len(set); i++ {
		if uint16(set[i]) == c {
			return true
		}
	}
	return false
}

func uriIsUnreserved(c uint16) bool { return uriIsAlphanum(c) || uriIn(c, "_-!.~'()*") }

// uriIsOther is the "other" category: non-ASCII characters that are neither
// ISO control characters nor space characters (Zs, Zl, Zp).
func uriIsOther(c uint16) bool {
	if c <= 0x9F {
		return false
	}
	switch {
	case c == 0x00A0, c == 0x1680, c >= 0x2000 && c <= 0x200A, c == 0x2028, c == 0x2029,
		c == 0x202F, c == 0x205F, c == 0x3000:
		return false
	}
	return true
}

// Characters allowed (besides escaped octets) in each component.
func uriIsPathChar(c uint16) bool {
	return uriIsUnreserved(c) || uriIsOther(c) || uriIn(c, ";:@&=+$,/")
}

func uriIsUricChar(c uint16) bool {
	return uriIsUnreserved(c) || uriIsOther(c) || uriIn(c, ";/?:@&=+$,[]")
}

func uriIsUserInfoChar(c uint16) bool {
	return uriIsUnreserved(c) || uriIsOther(c) || uriIn(c, ";:&=+$,")
}

func uriIsRegNameChar(c uint16) bool {
	return uriIsUnreserved(c) || uriIsOther(c) || uriIn(c, "$,;:@&=+")
}

// ---------------------------------------------------------------------------
// Parser

type uriParser struct {
	input string
	s     []uint16
}

func (p *uriParser) fail(reason string, index int) *uriSyntaxError {
	return &uriSyntaxError{input: p.input, reason: reason, index: index}
}

func (p *uriParser) str(from, to int) *string {
	v := StringFromUTF16(p.s[from:to])
	return &v
}

// find returns the index of the first character in [from, to) that is in
// set, or to.
func (p *uriParser) find(from, to int, set string) int {
	for i := from; i < to; i++ {
		if uriIn(p.s[i], set) {
			return i
		}
	}
	return to
}

// check validates [from, to) against legal (plus escaped octets).
func (p *uriParser) check(from, to int, legal func(uint16) bool, component string) *uriSyntaxError {
	for i := from; i < to; i++ {
		c := p.s[i]
		if c == '%' {
			if err := p.checkEscape(i, to); err != nil {
				return err
			}
			i += 2
			continue
		}
		if !legal(c) {
			return p.fail("Illegal character in "+component, i)
		}
	}
	return nil
}

// checkEscape checks the escaped octet at i ('%' and two hex digits).
func (p *uriParser) checkEscape(i, to int) *uriSyntaxError {
	if i+2 >= to || !uriIsHex(p.s[i+1]) || !uriIsHex(p.s[i+2]) {
		return p.fail("Malformed escape pair", i)
	}
	return nil
}

func (p *uriParser) parse(u *URI) *uriSyntaxError {
	n := len(p.s)
	start := 0
	if c := p.find(0, n, ":/?#"); c < n && p.s[c] == ':' {
		if c == 0 {
			return p.fail("Expected scheme name", 0)
		}
		if !uriIsAlpha(p.s[0]) {
			return p.fail("Illegal character in scheme name", 0)
		}
		for i := 1; i < c; i++ {
			if x := p.s[i]; !uriIsAlphanum(x) && x != '+' && x != '-' && x != '.' {
				return p.fail("Illegal character in scheme name", i)
			}
		}
		u.hasScheme = true
		start = c + 1
		if start == n || p.s[start] == '#' {
			return p.fail("Expected scheme-specific part", start)
		}
		if p.s[start] != '/' {
			end := p.find(start, n, "#")
			if err := p.check(start, end, uriIsUricChar, "opaque part"); err != nil {
				return err
			}
			return p.parseFragment(end)
		}
	}
	i := start
	if i+1 < n && p.s[i] == '/' && p.s[i+1] == '/' {
		authStart := i + 2
		authEnd := p.find(authStart, n, "/?#")
		if authEnd == authStart {
			if authEnd == n {
				return p.fail("Expected authority", authEnd)
			}
		} else {
			if err := p.parseAuthority(authStart, authEnd); err != nil {
				return err
			}
		}
		i = authEnd
	}
	pathEnd := p.find(i, n, "?#")
	if err := p.check(i, pathEnd, uriIsPathChar, "path"); err != nil {
		return err
	}
	u.path = p.str(i, pathEnd)
	i = pathEnd
	if i < n && p.s[i] == '?' {
		queryEnd := p.find(i+1, n, "#")
		if err := p.check(i+1, queryEnd, uriIsUricChar, "query"); err != nil {
			return err
		}
		i = queryEnd
	}
	return p.parseFragment(i)
}

func (p *uriParser) parseFragment(i int) *uriSyntaxError {
	n := len(p.s)
	if i < n && p.s[i] == '#' {
		if err := p.check(i+1, n, uriIsUricChar, "fragment"); err != nil {
			return err
		}
	}
	return nil
}

// parseAuthority parses the non-empty authority [from, to).
func (p *uriParser) parseAuthority(from, to int) *uriSyntaxError {
	regEnd := from
	for regEnd < to {
		c := p.s[regEnd]
		if c == '%' {
			if err := p.checkEscape(regEnd, to); err != nil {
				return err
			}
			regEnd += 3
			continue
		}
		if !uriIsRegNameChar(c) {
			break
		}
		regEnd++
	}
	if regEnd == to {
		_ = p.parseServer(from, to)
		return nil
	}
	for i := regEnd; i < to; i++ {
		if c := p.s[i]; !uriIsRegNameChar(c) && c != '[' && c != ']' && c != '%' {
			return p.fail("Illegal character in authority", regEnd)
		}
	}
	return p.parseServer(from, to)
}

// parseServer parses [from, to) as [user-info@]host[:port].
func (p *uriParser) parseServer(from, to int) *uriSyntaxError {
	i := from
	if at := p.find(from, to, "@"); at < to {
		for k := from; k < at; k++ {
			c := p.s[k]
			if c == '%' {
				if err := p.checkEscape(k, at); err != nil {
					return err
				}
				k += 2
				continue
			}
			if !uriIsUserInfoChar(c) {
				return p.fail("Illegal character in user info", k)
			}
		}
		i = at + 1
	}
	var hostEnd int
	if i < to && p.s[i] == '[' {
		close := p.find(i+1, to, "]")
		if close == to || close == i+1 {
			return p.fail("Expected closing bracket for IPv6 address", close)
		}
		scope := p.find(i+1, close, "%")
		if err := p.parseIPv6(i+1, scope); err != nil {
			return err
		}
		if scope < close {
			if scope+1 == close {
				return p.fail("scope id expected", -1)
			}
			for k := scope + 1; k < close; k++ {
				if c := p.s[k]; !uriIsAlphanum(c) && c != '_' && c != '.' {
					return p.fail("Illegal character in scope id", k)
				}
			}
		}
		hostEnd = close + 1
		if hostEnd < to && p.s[hostEnd] != ':' {
			return p.fail("Expected port number", hostEnd)
		}
	} else {
		var err *uriSyntaxError
		hostEnd, err = p.parseHostName(i, to)
		if err != nil {
			return err
		}
	}
	if hostEnd < to {
		// p.s[hostEnd] == ':'
		portStart := hostEnd + 1
		port := int64(0)
		for k := portStart; k < to; k++ {
			c := p.s[k]
			if !uriIsDigit(c) {
				return p.fail("Illegal character in port number", k)
			}
			if port <= 1<<31 {
				port = port*10 + int64(c-'0')
			}
		}
		if portStart < to && port > 1<<31-1 {
			return p.fail("Malformed port number", portStart)
		}
	}
	return nil
}

// parseHostName parses a dotted-quad or a host name starting at from and
// returns its end, which is to or the index of a ':'.
func (p *uriParser) parseHostName(from, to int) (int, *uriSyntaxError) {
	if end, ok := p.dottedQuad(from, to); ok && (end == to || p.s[end] == ':') {
		return end, nil
	}
	i := from
	labels, lastLabel := 0, from
	for i < to && uriIsAlphanum(p.s[i]) {
		labelStart := i
		last := i
		for i < to && (uriIsAlphanum(p.s[i]) || p.s[i] == '-') {
			if p.s[i] != '-' {
				last = i
			}
			i++
		}
		i = last + 1 // a label cannot end with '-'
		labels++
		lastLabel = labelStart
		if i < to && p.s[i] == '.' {
			i++
			continue
		}
		break
	}
	if i < to && p.s[i] != ':' {
		return 0, p.fail("Illegal character in hostname", i)
	}
	if labels == 0 {
		return 0, p.fail("Expected hostname", i)
	}
	if labels > 1 && !uriIsAlpha(p.s[lastLabel]) {
		return 0, p.fail("Illegal character in hostname", lastLabel)
	}
	return i, nil
}

// dottedQuad reads four dot-separated decimal numbers up to 255 at from.
func (p *uriParser) dottedQuad(from, to int) (int, bool) {
	i := from
	for part := 0; part < 4; part++ {
		if part > 0 {
			if i >= to || p.s[i] != '.' {
				return i, false
			}
			i++
		}
		start, value := i, 0
		for i < to && uriIsDigit(p.s[i]) {
			if value <= 255 {
				value = value*10 + int(p.s[i]-'0')
			}
			i++
		}
		if i == start || value > 255 {
			return start, false
		}
	}
	return i, true
}

// parseIPv6 checks the IPv6 literal [from, to) (see the observed behaviour
// above).
func (p *uriParser) parseIPv6(from, to int) *uriSyntaxError {
	units := 0 // 16-bit pieces
	compressed := false
	hexRun := func(i int) int {
		j := i
		for j < to && uriIsHex(p.s[j]) {
			j++
		}
		return j
	}
	startsQuad := func(i int) bool {
		j := hexRun(i)
		return j < to && p.s[j] == '.'
	}
	// groups reads hex groups separated by single colons from i, which holds
	// hex digits that do not start a dotted quad. It stops before "::", and
	// before a ':' that introduces a dotted quad.
	groups := func(i int) (int, *uriSyntaxError) {
		for {
			j := hexRun(i)
			if j-i > 4 {
				return 0, p.fail("IPv6 hexadecimal digit sequence too long", i)
			}
			units++
			i = j
			if i >= to || p.s[i] != ':' || i+1 < to && p.s[i+1] == ':' {
				return i, nil
			}
			if hexRun(i+1) == i+1 {
				return 0, p.fail("Expected digits for an IPv6 address", i+1)
			}
			if startsQuad(i + 1) {
				return i, nil
			}
			i++
		}
	}
	quad := func(i int, expected string) (int, *uriSyntaxError) {
		j := i
		for j < to && (uriIsDigit(p.s[j]) || p.s[j] == '.') {
			j++
		}
		if j == i || j != to {
			return 0, p.fail("Expected "+expected, i)
		}
		k := i
		for part := 0; part < 4; part++ {
			if part > 0 {
				if k >= j || p.s[k] != '.' {
					return 0, p.fail("Malformed IPv4 address", k)
				}
				k++
			}
			start, value := k, 0
			for k < j && uriIsDigit(p.s[k]) {
				if value <= 255 {
					value = value*10 + int(p.s[k]-'0')
				}
				k++
			}
			if k == start || value > 255 {
				return 0, p.fail("Malformed IPv4 address", start)
			}
		}
		if k != j {
			return 0, p.fail("Malformed IPv4 address", k)
		}
		units += 2
		return k, nil
	}
	afterCompression := func(i int) (int, *uriSyntaxError) {
		if i == to {
			return i, nil
		}
		if hexRun(i) > i && !startsQuad(i) {
			var err *uriSyntaxError
			i, err = groups(i)
			if err != nil {
				return 0, err
			}
			if i < to && p.s[i] == ':' {
				return quad(i+1, "hex digits or IPv4 address")
			}
			return i, nil
		}
		return quad(i, "hex digits or IPv4 address")
	}
	i := from
	var err *uriSyntaxError
	switch {
	case hexRun(i) > i && !startsQuad(i):
		i, err = groups(i)
		if err != nil {
			return err
		}
		if i+1 < to && p.s[i] == ':' && p.s[i+1] == ':' {
			compressed = true
			i, err = afterCompression(i + 2)
			if err != nil {
				return err
			}
		} else if i < to && p.s[i] == ':' {
			i, err = quad(i+1, "IPv4 address")
			if err != nil {
				return err
			}
		}
	case i+1 < to && p.s[i] == ':' && p.s[i+1] == ':':
		compressed = true
		i, err = afterCompression(i + 2)
		if err != nil {
			return err
		}
	default:
		return p.fail("Malformed IPv6 address", from)
	}
	switch {
	case i < to:
		return p.fail("Malformed IPv6 address", from)
	case units > 8:
		return p.fail("IPv6 address too long", from)
	case !compressed && units < 8:
		return p.fail("IPv6 address too short", from)
	case compressed && units == 8:
		return p.fail("Malformed IPv6 address", from)
	}
	return nil
}

type uriSyntaxError struct {
	input  string
	reason string
	index  int
}

func (e *uriSyntaxError) Error() string {
	msg := e.reason
	if e.index > -1 {
		msg += " at index " + strconv.Itoa(e.index)
	}
	return "java.net.URISyntaxException: " + msg + ": " + e.input
}

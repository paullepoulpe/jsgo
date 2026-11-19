package javanio

import (
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// URISyntaxException models java.net.URISyntaxException.
type URISyntaxException struct {
	Input  string
	Reason string
	Index  int
}

func NewURISyntaxException(input, reason string) *URISyntaxException {
	return &URISyntaxException{Input: input, Reason: reason, Index: -1}
}

func NewURISyntaxExceptionIndex(input, reason string, index int) *URISyntaxException {
	return &URISyntaxException{Input: input, Reason: reason, Index: index}
}

func (e *URISyntaxException) GetInput() string  { return e.Input }
func (e *URISyntaxException) GetReason() string { return e.Reason }
func (e *URISyntaxException) GetIndex() int     { return e.Index }

func (e *URISyntaxException) GetMessage() string {
	msg := e.Reason
	if e.Index > -1 {
		msg += " at index " + strconv.Itoa(e.Index)
	}
	return msg + ": " + e.Input
}

func (e *URISyntaxException) Error() string {
	return "java.net.URISyntaxException: " + e.GetMessage()
}

// URI models java.net.URI (RFC 2396 + Java deviations).
type URI struct {
	str          string
	scheme       string
	hasScheme    bool
	authority    string
	hasAuthority bool
	path         string
	hasPath      bool
	query        string
	hasQuery     bool
	fragment     string
	hasFragment  bool
}

// NewURI parses str as java.net.URI(String), returning URISyntaxException on syntax error.
func NewURI(str string) (*URI, error) {
	p := &uriParser{input: str, s: javalang.StringToUTF16(str)}
	u := &URI{str: str}
	if err := p.parse(u); err != nil {
		return nil, err
	}
	return u, nil
}

// URI_Create implements URI.create(String), panicking with IllegalArgumentException on syntax error.
func URI_Create(str string) *URI {
	u, err := NewURI(str)
	if err != nil {
		panic(javalang.NewIllegalArgumentException(err.Error()))
	}
	return u
}

func (u *URI) IsAbsolute() bool { return u.hasScheme }
func (u *URI) IsOpaque() bool   { return !u.hasPath }

func (u *URI) GetScheme() (string, bool) {
	return u.scheme, u.hasScheme
}

func (u *URI) GetSchemeOrEmpty() string {
	return u.scheme
}

func (u *URI) GetRawAuthority() (string, bool) {
	return u.authority, u.hasAuthority
}

func (u *URI) GetAuthority() (string, bool) {
	if !u.hasAuthority {
		return "", false
	}
	return uriDecodeEscapedOctets(u.authority), true
}

func (u *URI) GetRawPath() (string, bool) {
	return u.path, u.hasPath
}

func (u *URI) GetPath() (string, bool) {
	if !u.hasPath {
		return "", false
	}
	return uriDecodeEscapedOctets(u.path), true
}

func (u *URI) GetPathOrEmpty() string {
	p, ok := u.GetPath()
	if !ok {
		return ""
	}
	return p
}

func (u *URI) GetRawQuery() (string, bool) {
	return u.query, u.hasQuery
}

func (u *URI) GetQuery() (string, bool) {
	if !u.hasQuery {
		return "", false
	}
	return uriDecodeEscapedOctets(u.query), true
}

func (u *URI) GetRawFragment() (string, bool) {
	return u.fragment, u.hasFragment
}

func (u *URI) GetFragment() (string, bool) {
	if !u.hasFragment {
		return "", false
	}
	return uriDecodeEscapedOctets(u.fragment), true
}

func (u *URI) Normalize() *URI {
	if !u.hasPath || u.path == "" {
		return u
	}
	np := normalizeURIPath(u.path)
	if np == u.path {
		return u
	}
	cp := *u
	cp.path = np
	cp.str = cp.rebuildString()
	return &cp
}

func (u *URI) Resolve(child *URI) *URI {
	if child == nil {
		panic(javalang.NewNullPointerExceptionMsg("uri is null"))
	}
	if child.hasScheme || u.IsOpaque() {
		return child
	}
	// Fragment-only relative reference
	if !child.hasAuthority && child.path == "" && !child.hasQuery && child.hasFragment {
		cp := *u
		cp.fragment = child.fragment
		cp.hasFragment = true
		cp.str = cp.rebuildString()
		return &cp
	}
	res := &URI{
		scheme:      u.scheme,
		hasScheme:   u.hasScheme,
		fragment:    child.fragment,
		hasFragment: child.hasFragment,
	}
	if child.hasAuthority {
		res.authority = child.authority
		res.hasAuthority = true
		res.path = child.path
		res.hasPath = child.hasPath
		res.query = child.query
		res.hasQuery = child.hasQuery
	} else {
		res.authority = u.authority
		res.hasAuthority = u.hasAuthority
		res.hasPath = true
		if strings.HasPrefix(child.path, "/") {
			res.path = child.path
			res.query = child.query
			res.hasQuery = child.hasQuery
		} else if child.path == "" {
			res.path = u.path
			if child.hasQuery {
				res.query = child.query
				res.hasQuery = true
			} else {
				res.query = u.query
				res.hasQuery = u.hasQuery
			}
		} else {
			idx := strings.LastIndexByte(u.path, '/')
			var base string
			if idx >= 0 {
				base = u.path[:idx+1] + child.path
			} else {
				base = child.path
			}
			res.path = normalizeURIPath(base)
			res.query = child.query
			res.hasQuery = child.hasQuery
		}
	}
	res.str = res.rebuildString()
	return res
}

func (u *URI) ResolveString(str string) *URI {
	return u.Resolve(URI_Create(str))
}

func (u *URI) rebuildString() string {
	var sb strings.Builder
	if u.hasScheme {
		sb.WriteString(u.scheme)
		sb.WriteByte(':')
	}
	if u.hasAuthority {
		sb.WriteString("//")
		sb.WriteString(u.authority)
	}
	if u.hasPath {
		sb.WriteString(u.path)
	}
	if u.hasQuery {
		sb.WriteByte('?')
		sb.WriteString(u.query)
	}
	if u.hasFragment {
		sb.WriteByte('#')
		sb.WriteString(u.fragment)
	}
	return sb.String()
}

func normalizeURIPath(p string) string {
	if p == "" {
		return ""
	}
	leadingSlash := strings.HasPrefix(p, "/")
	trailingSlash := len(p) > 1 && strings.HasSuffix(p, "/")
	parts := strings.Split(p, "/")
	var stack []string
	for i, seg := range parts {
		if seg == "" {
			if i == 0 && leadingSlash {
				continue
			}
			if i == len(parts)-1 && trailingSlash {
				continue
			}
			continue
		}
		if seg == "." {
			if i == len(parts)-1 {
				trailingSlash = true
			}
			continue
		}
		if seg == ".." {
			if len(stack) > 0 && stack[len(stack)-1] != ".." {
				stack = stack[:len(stack)-1]
				if i == len(parts)-1 {
					trailingSlash = true
				}
				continue
			}
			if leadingSlash {
				stack = append(stack, "..")
				continue
			}
		}
		stack = append(stack, seg)
	}
	var sb strings.Builder
	if leadingSlash {
		sb.WriteByte('/')
	}
	sb.WriteString(strings.Join(stack, "/"))
	if trailingSlash && sb.Len() > 0 && sb.String() != "/" {
		sb.WriteByte('/')
	}
	res := sb.String()
	if !leadingSlash && len(stack) > 0 && strings.Contains(stack[0], ":") {
		res = "./" + res
	}
	return res
}

func (u *URI) String() string { return u.str }

func (u *URI) Equals(other any) bool {
	o, ok := other.(*URI)
	if !ok || o == nil {
		return false
	}
	if u.hasScheme != o.hasScheme || !strings.EqualFold(u.scheme, o.scheme) {
		return false
	}
	if u.hasFragment != o.hasFragment || !uriEqualComponent(u.fragment, o.fragment) {
		return false
	}
	return uriEqualComponent(u.str, o.str) ||
		(u.hasAuthority == o.hasAuthority &&
			strings.EqualFold(u.authority, o.authority) &&
			u.hasPath == o.hasPath && uriEqualComponent(u.path, o.path) &&
			u.hasQuery == o.hasQuery && uriEqualComponent(u.query, o.query))
}

func uriEqualComponent(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] == '%' && b[i] == '%' && i+2 < len(a) {
			if !strings.EqualFold(a[i+1:i+3], b[i+1:i+3]) {
				return false
			}
			i += 2
			continue
		}
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (u *URI) HashCode() int32 {
	return javalang.StringHashCode(javalang.StringToLowerCaseASCII(u.str))
}

// --- RFC 2396 + Java URI Parser ---

func uriDecodeEscapedOctets(s string) string {
	if !strings.Contains(s, "%") {
		return s
	}
	units := javalang.StringToUTF16(s)
	out := make([]uint16, 0, len(units))
	for i := 0; i < len(units); {
		if units[i] != '%' {
			out = append(out, units[i])
			i++
			continue
		}
		var octets []byte
		for i+2 < len(units) && units[i] == '%' {
			octets = append(octets, byte(uriHexVal(units[i+1])<<4|uriHexVal(units[i+2])))
			i += 3
		}
		out = append(out, javalang.StringToUTF16(StandardCharsets_UTF_8.Decode(octets))...)
	}
	return javalang.StringFromUTF16(out)
}

func uriHexVal(c uint16) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	default:
		return int(c-'A') + 10
	}
}

func uriIsAlpha(c uint16) bool    { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func uriIsDigit(c uint16) bool    { return c >= '0' && c <= '9' }
func uriIsAlphanum(c uint16) bool { return uriIsAlpha(c) || uriIsDigit(c) }
func uriIsHex(c uint16) bool {
	return uriIsDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
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

type uriParser struct {
	input string
	s     []uint16
}

func (p *uriParser) fail(reason string, index int) *URISyntaxException {
	return &URISyntaxException{Input: p.input, Reason: reason, Index: index}
}

func (p *uriParser) str(from, to int) string {
	return javalang.StringFromUTF16(p.s[from:to])
}

func (p *uriParser) find(from, to int, set string) int {
	for i := from; i < to; i++ {
		if uriIn(p.s[i], set) {
			return i
		}
	}
	return to
}

func (p *uriParser) checkEscape(i, to int) *URISyntaxException {
	if i+2 >= to || !uriIsHex(p.s[i+1]) || !uriIsHex(p.s[i+2]) {
		return p.fail("Malformed escape pair", i)
	}
	return nil
}

func (p *uriParser) check(from, to int, legal func(uint16) bool, component string) *URISyntaxException {
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

func (p *uriParser) parse(u *URI) *URISyntaxException {
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
		u.scheme = p.str(0, c)
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
			return p.parseFragment(u, end)
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
			u.authority = p.str(authStart, authEnd)
			u.hasAuthority = true
		}
		i = authEnd
	}
	pathEnd := p.find(i, n, "?#")
	if err := p.check(i, pathEnd, uriIsPathChar, "path"); err != nil {
		return err
	}
	u.path, u.hasPath = p.str(i, pathEnd), true
	i = pathEnd
	if i < n && p.s[i] == '?' {
		queryEnd := p.find(i+1, n, "#")
		if err := p.check(i+1, queryEnd, uriIsUricChar, "query"); err != nil {
			return err
		}
		u.query, u.hasQuery = p.str(i+1, queryEnd), true
		i = queryEnd
	}
	return p.parseFragment(u, i)
}

func (p *uriParser) parseFragment(u *URI, i int) *URISyntaxException {
	n := len(p.s)
	if i < n && p.s[i] == '#' {
		if err := p.check(i+1, n, uriIsUricChar, "fragment"); err != nil {
			return err
		}
		u.fragment, u.hasFragment = p.str(i+1, n), true
	}
	return nil
}

func (p *uriParser) parseAuthority(from, to int) *URISyntaxException {
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
		return nil
	}
	for i := regEnd; i < to; i++ {
		if c := p.s[i]; !uriIsRegNameChar(c) && c != '[' && c != ']' && c != '%' {
			return p.fail("Illegal character in authority", regEnd)
		}
	}
	return p.parseServer(from, to)
}

func (p *uriParser) parseServer(from, to int) *URISyntaxException {
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
		closeIdx := p.find(i+1, to, "]")
		if closeIdx == to || closeIdx == i+1 {
			return p.fail("Expected closing bracket for IPv6 address", closeIdx)
		}
		// Validate hex/colon/dot/scope characters inside IPv6 brackets
		scope := p.find(i+1, closeIdx, "%")
		for k := i + 1; k < scope; k++ {
			if c := p.s[k]; !uriIsHex(c) && c != ':' && c != '.' {
				return p.fail("Malformed IPv6 address", k)
			}
		}
		if scope < closeIdx {
			if scope+1 == closeIdx {
				return p.fail("scope id expected", -1)
			}
			for k := scope + 1; k < closeIdx; k++ {
				if c := p.s[k]; !uriIsAlphanum(c) && c != '_' && c != '.' {
					return p.fail("Illegal character in scope id", k)
				}
			}
		}
		hostEnd = closeIdx + 1
		if hostEnd < to && p.s[hostEnd] != ':' {
			return p.fail("Expected port number", hostEnd)
		}
	} else {
		hostEnd = p.find(i, to, ":")
	}
	if hostEnd < to {
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
		if portStart < to && port > (1<<31-1) {
			return p.fail("Malformed port number", portStart)
		}
	}
	return nil
}

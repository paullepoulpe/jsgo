package jutil

// Go-only: language-adapter — java.util.regex.Pattern on top of Go's RE2.
//
// Java regexes are translated to RE2 syntax so that the constructs they have
// in common keep Java's meaning:
//   - `.` excludes Java's line terminators (\n, \r, \u0085, \u2028, \u2029),
//     not only \n as in RE2 (unless (?s) is active);
//   - \s, \S (and \s inside classes) are Java's [ \t\n\x0B\f\r];
//   - \v, \h (also inside classes), \V and \H are the vertical and
//     horizontal whitespace classes of the Pattern Javadoc;
//   - \Q...\E quoting is resolved here;
//   - Matcher.matches() anchors the whole input with \A(?:...)\z;
//   - `^` and `$` keep their default (non-MULTILINE) Java meaning, see
//     below.
//
// Constructs whose Java semantics RE2 cannot reproduce (back-references,
// lookaround, possessive quantifiers, atomic groups, class intersections,
// nested classes, \b \B \z \Z \G, a `$` that is not at the end of the
// pattern or of an alternative, embedded flags) make Pattern_Compile panic
// with UnsupportedOperationException rather than silently match
// differently.
//
// Anchors. The Javadoc of java.util.regex.Pattern ("Line terminators")
// says that without MULTILINE `^` matches only at the beginning of the
// whole input, and `$` at its end and also just before a line terminator
// that ends the input. The line terminators are \n, \r, \r\n, \u0085,
// \u2028 and \u2029. Observed (black-box, JDK 21): `$` does not match
// between the \r and the \n of a final \r\n; `^` does not match where a
// later find() resumes. Here:
//   - `^` becomes RE2's `^` (start of text) for a search from offset 0, and
//     a class that matches nothing for a search from further on (the
//     search runs on the input's suffix, where RE2's `^` would match at its
//     first char);
//   - `$` becomes `\z` in the anchored regex of matches() (a match that
//     must end at the end of the input), and, for searches, end of text or
//     the input's final line terminator followed by end of text. The
//     terminator is consumed, in a capturing group of its own, so the
//     search result is cut back to where that group starts: nothing can
//     follow such a `$` in its alternative, so every offset after the `$`
//     is the `$`'s own. Which terminators qualify depends on the input: an
//     input ending in \r\n gets a regex that accepts only \r\n there, so
//     that the lone \n after the \r cannot satisfy `$`.

import (
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

// Pattern is a compiled java.util.regex.Pattern.
type Pattern struct {
	pattern string
	groups  int            // the Java pattern's capturing groups
	full    *regexp.Regexp // anchored at both ends (matches)
	// searchers[resumed][crlf] are the unanchored regexes searches use:
	// resumed for a search starting past offset 0 (`^` cannot match),
	// crlf for an input that ends in "\r\n" (see the anchors above).
	// Variants that would be identical share one compiled regex.
	searchers [2][2]searcher
}

// searcher is one unanchored translation of the pattern.
type searcher struct {
	re *regexp.Regexp
	// marks are the RE2 group numbers, in increasing order, of the groups
	// that hold the line terminator a `$` consumed; empty without `$`.
	marks []int
}

// javaRegexMode selects how translateJavaRegex renders the anchors.
type javaRegexMode struct {
	anchored bool // for matches(): `$` is end of text
	resumed  bool // the search starts past offset 0: `^` never matches
	crlf     bool // the input ends in "\r\n"
}

// javaRegexInfo is what translateJavaRegex learns about the pattern.
type javaRegexInfo struct {
	groups int   // capturing groups of the Java pattern
	marks  []int // see searcher.marks
	caret  bool  // the pattern uses `^` (outside classes)
	dollar bool  // the pattern uses `$` (outside classes)
}

const javaDotClass = `[^\n\r\x{85}\x{2028}\x{2029}]`
const javaSpaceChars = `\t\n\x0B\f\r `

// javaVerticalSpaceChars and javaHorizontalSpaceChars are the classes \v and
// \h of the Pattern Javadoc: [\n\x0B\f\r\x85\u2028\u2029] and
// [ \t\xA0\u1680\u180e\u2000-\u200a\u202f\u205f\u3000].
const javaVerticalSpaceChars = `\n\x0B\f\r\x{85}\x{2028}\x{2029}`
const javaHorizontalSpaceChars = ` \t\x{A0}\x{1680}\x{180E}\x{2000}-\x{200A}\x{202F}\x{205F}\x{3000}`

// Pattern_Compile is Pattern.compile(String).
//
// Go-only: performance — translating a Java regex and compiling its RE2
// variants costs far more than Pattern.compile does on the JVM, and ported
// code compiles the same constant regexes again and again (for example the
// CompilerOptions field initializer that every CodePrinter.Builder runs).
// Successful compilations are cached by regex. A compiled Pattern is
// immutable, so each call still returns a distinct Pattern (Java identity)
// sharing the cached translation; a regex that fails to compile is not
// cached and fails the same way every time.
func Pattern_Compile(regex string) *Pattern {
	if cached, ok := patternCache.Load(regex); ok {
		p := *cached.(*Pattern)
		return &p
	}
	p := patternCompile(regex)
	if patternCacheSize.Add(1) <= patternCacheLimit {
		cached := *p
		patternCache.Store(regex, &cached)
	}
	return p
}

// patternCache maps a regex to its compiled Pattern (see Pattern_Compile).
// It stops growing at patternCacheLimit entries, so that regexes built from
// input cannot make it grow without bound.
var patternCache sync.Map // string -> *Pattern

var patternCacheSize atomic.Int64

const patternCacheLimit = 4096

// patternCompile compiles regex (Pattern_Compile without the cache).
func patternCompile(regex string) *Pattern {
	whole, info := translateJavaRegex(regex, javaRegexMode{anchored: true})
	fullRe := `\A(?:` + whole + `)\z`
	p := &Pattern{
		pattern: regex,
		groups:  info.groups,
		full:    compileTranslatedJavaRegex(regex, fullRe),
	}
	for resumed := range 2 {
		for crlf := range 2 {
			switch {
			case resumed == 1 && !info.caret:
				p.searchers[1][crlf] = p.searchers[0][crlf]
			case crlf == 1 && !info.dollar:
				p.searchers[resumed][1] = p.searchers[resumed][0]
			default:
				t, ti := translateJavaRegex(regex, javaRegexMode{resumed: resumed == 1, crlf: crlf == 1})
				p.searchers[resumed][crlf] = searcher{compileTranslatedJavaRegex(regex, t), ti.marks}
			}
		}
	}
	return p
}

// searcherFor picks the translation for a search of input from byte offset
// from.
func (p *Pattern) searcherFor(input string, from int) searcher {
	resumed, crlf := 0, 0
	if from > 0 {
		resumed = 1
	}
	if strings.HasSuffix(input, "\r\n") {
		crlf = 1
	}
	return p.searchers[resumed][crlf]
}

// findFrom is the leftmost match of the pattern in input that starts at
// byte offset from or later, as byte offsets in Go's submatch layout over
// the Java pattern's groups, or nil. RE2's leftmost-first search picks the
// match and submatches that Java's backtracking finds first.
func (p *Pattern) findFrom(input string, from int) []int {
	s := p.searcherFor(input, from)
	loc := s.re.FindStringSubmatchIndex(input[from:])
	if loc == nil {
		return nil
	}
	for i, off := range loc {
		if off >= 0 {
			loc[i] = off + from
		}
	}
	if len(s.marks) == 0 {
		return loc
	}
	// Drop the terminator groups; if one took part, the `$` sat where it
	// starts, and every later offset moves back there.
	cut := -1
	for _, g := range s.marks {
		if at := loc[2*g]; at >= 0 && (cut < 0 || at < cut) {
			cut = at
		}
	}
	out := make([]int, 0, 2*(p.groups+1))
	next := 0
	for g := 0; 2*g < len(loc); g++ {
		if next < len(s.marks) && s.marks[next] == g {
			next++
			continue
		}
		for _, off := range loc[2*g : 2*g+2] {
			if cut >= 0 && off > cut {
				off = cut
			}
			out = append(out, off)
		}
	}
	return out
}

// Pattern_Quote is Pattern.quote(String): a regex that matches s literally
// (see the Javadoc of java.util.regex.Pattern.quote). The text is wrapped in
// \Q...\E.
//
// Observed behaviour (black-box): the result is always "\Q" + s' + "\E",
// where every "\E" inside s is written as "\E\\E\Q" (close the quote,
// match a backslash and an 'E', reopen the quote). For example a\Eb becomes
// \Qa\E\\E\Qb\E, and the empty string becomes \Q\E.
func Pattern_Quote(s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 4)
	sb.WriteString(`\Q`)
	for {
		i := strings.Index(s, `\E`)
		if i < 0 {
			break
		}
		sb.WriteString(s[:i])
		sb.WriteString(`\E\\E\Q`)
		s = s[i+2:]
	}
	sb.WriteString(s)
	sb.WriteString(`\E`)
	return sb.String()
}

// Pattern is Pattern.pattern().
func (p *Pattern) Pattern() string { return p.pattern }

// String is Pattern.toString().
func (p *Pattern) String() string { return p.pattern }

// Matches is p.matcher(s).matches().
func (p *Pattern) Matches(s string) bool { return p.full.MatchString(s) }

// Find is p.matcher(s).find() from the start of s.
func (p *Pattern) Find(s string) bool { return p.searcherFor(s, 0).re.MatchString(s) }

func unsupportedJavaRegex(regex, what string) {
	panic(NewUnsupportedOperationExceptionMsg("Java regex construct not supported (" + what + "): " + regex))
}

// javaFinalTerminators is the class of the one-char line terminators that
// may end the input after a `$`; javaFinalCRLF replaces it for inputs that
// end in "\r\n".
const (
	javaFinalTerminators = `[\n\r\x{85}\x{2028}\x{2029}]`
	javaFinalCRLF        = `\r\n`
	// javaNoMatch is a class no char belongs to.
	javaNoMatch = `[^\x00-\x{10FFFF}]`
)

// translateJavaRegex rewrites a Java regex into an equivalent RE2 regex,
// rendering the anchors as mode asks (see the file comment).
func translateJavaRegex(regex string, mode javaRegexMode) (string, javaRegexInfo) {
	var sb strings.Builder
	var info javaRegexInfo
	multiline := false
	dotall := false
	rs := []rune(regex)
	n := len(rs)
	inClass := false
	classStart := 0
	// Group nesting, and the RE2 number of the last group opened (the
	// Java pattern's groups and the terminator groups of `$`).
	depth, lastGroup := 0, 0
	// After a `$`, only the end of its alternative may follow: tailDepth is
	// the nesting depth at which that still has to hold (-1: no `$`
	// pending), and endsInDollar[d] records that an earlier alternative of
	// the open group at depth d ended with `$`, so that the group's end is
	// a tail position as well.
	tailDepth := -1
	var endsInDollar []bool
	afterGroupClose := func(i int) {
		if i+1 < n && strings.ContainsRune("*+?{", rs[i+1]) {
			unsupportedJavaRegex(regex, "quantified group containing $")
		}
	}
	for i := 0; i < n; i++ {
		c := rs[i]
		if !multiline && tailDepth >= 0 && !inClass && c != ')' && c != '|' && c != '$' {
			unsupportedJavaRegex(regex, "$ not at the end of an alternative")
		}
		switch {
		case c == '\\':
			if i+1 >= n {
				panic(NewPatternSyntaxException("Unexpected internal error", regex, n))
			}
			i++
			e := rs[i]
			switch {
			case e == 'Q':
				// Quote until \E (or the end).
				end := -1
				for j := i + 1; j+1 < n; j++ {
					if rs[j] == '\\' && rs[j+1] == 'E' {
						end = j
						break
					}
				}
				var lit string
				if end == -1 {
					lit = string(rs[i+1:])
					i = n
				} else {
					lit = string(rs[i+1 : end])
					i = end + 1
				}
				if inClass {
					for _, r := range lit {
						sb.WriteString(regexp.QuoteMeta(string(r)))
					}
				} else {
					sb.WriteString(regexp.QuoteMeta(lit))
				}
			case e == 'E':
				// A stray \E is ignored by Java.
			case e == 's':
				if inClass {
					sb.WriteString(javaSpaceChars)
				} else {
					sb.WriteString(`[` + javaSpaceChars + `]`)
				}
			case e == 'S':
				if inClass {
					unsupportedJavaRegex(regex, `\S in class`)
				}
				sb.WriteString(`[^` + javaSpaceChars + `]`)
			case e == 'v' || e == 'h':
				chars := javaVerticalSpaceChars
				if e == 'h' {
					chars = javaHorizontalSpaceChars
				}
				if inClass {
					sb.WriteString(chars)
				} else {
					sb.WriteString(`[` + chars + `]`)
				}
			case e == 'V' || e == 'H':
				if inClass {
					unsupportedJavaRegex(regex, `\`+string(e)+` in class`)
				}
				chars := javaVerticalSpaceChars
				if e == 'H' {
					chars = javaHorizontalSpaceChars
				}
				sb.WriteString(`[^` + chars + `]`)
			case e == 'd' || e == 'D' || e == 'w' || e == 'W' || e == 't' || e == 'n' || e == 'r' || e == 'f':
				sb.WriteByte('\\')
				sb.WriteRune(e)
			case e == 'e':
				sb.WriteString(`\x1B`)
			case e == 'a':
				sb.WriteString(`\x07`)
			case e == '0' || (e >= '1' && e <= '9') || e == 'k':
				unsupportedJavaRegex(regex, "back-reference or octal escape")
			case e == 'b' || e == 'B' || e == 'z' || e == 'Z' || e == 'G' || e == 'A' || e == 'R' || e == 'X' ||
				e == 'p' || e == 'P' || e == 'x' || e == 'u' || e == 'c' || e == 'N':
				unsupportedJavaRegex(regex, `\`+string(e))
			case (e >= 'a' && e <= 'z') || (e >= 'A' && e <= 'Z'):
				panic(NewPatternSyntaxException("Illegal/unsupported escape sequence", regex, i))
			default:
				sb.WriteString(regexp.QuoteMeta(string(e)))
			}
		case inClass:
			switch c {
			case '[':
				unsupportedJavaRegex(regex, "nested character class")
			case '&':
				if i+1 < n && rs[i+1] == '&' {
					unsupportedJavaRegex(regex, "class intersection")
				}
				sb.WriteRune(c)
			case ']':
				if i == classStart || (i == classStart+1 && rs[classStart] == '^') {
					// A leading ']' is a literal in Java as in RE2.
					sb.WriteString(`\]`)
				} else {
					inClass = false
					sb.WriteRune(c)
				}
			default:
				sb.WriteRune(c)
			}
		case c == '[':
			inClass = true
			classStart = i + 1
			sb.WriteRune(c)
		case c == '.':
			if dotall {
				sb.WriteByte('.')
			} else {
				sb.WriteString(javaDotClass)
			}
		case c == '^':
			info.caret = true
			if mode.resumed && !multiline {
				sb.WriteString(javaNoMatch)
			} else {
				sb.WriteByte('^')
			}
		case c == '$':
			info.dollar = true
			if multiline {
				sb.WriteByte('$')
			} else {
				tailDepth = depth
				switch {
				case mode.anchored:
					sb.WriteString(`\z`)
				default:
					lastGroup++
					info.marks = append(info.marks, lastGroup)
					term := javaFinalTerminators
					if mode.crlf {
						term = javaFinalCRLF
					}
					sb.WriteString(`(?:\z|(` + term + `)\z)`)
				}
			}
		case c == '|':
			if tailDepth == depth {
				for len(endsInDollar) <= depth {
					endsInDollar = append(endsInDollar, false)
				}
				endsInDollar[depth] = true
				tailDepth = -1
			}
			sb.WriteRune(c)
		case c == ')':
			if tailDepth == depth || depth < len(endsInDollar) && endsInDollar[depth] {
				if depth < len(endsInDollar) {
					endsInDollar[depth] = false
				}
				tailDepth = depth - 1
				afterGroupClose(i)
			}
			depth--
			sb.WriteRune(c)
		case c == '(':
			depth++
			if i+1 < n && rs[i+1] == '?' {
				if i+2 < n && rs[i+2] == ':' {
					sb.WriteString("(?:")
					i += 2
					continue
				}
				j := i + 2
				neg := false
				mSet, mVal := false, false
				sSet, sVal := false, false
				validFlags := j < n
				for ; j < n && rs[j] != ')' && rs[j] != ':'; j++ {
					switch rs[j] {
					case '-':
						neg = true
					case 'i':
					case 'm':
						mSet, mVal = true, !neg
					case 's':
						sSet, sVal = true, !neg
					default:
						validFlags = false
					}
				}
				if validFlags && j > i+2 && j < n && (rs[j] == ')' || rs[j] == ':') {
					if mSet {
						multiline = mVal
					}
					if sSet {
						dotall = sVal
					}
					sb.WriteString(string(rs[i : j+1]))
					if rs[j] == ')' {
						depth--
					}
					i = j
					continue
				}
				unsupportedJavaRegex(regex, "special group")
			}
			lastGroup++
			info.groups++
			sb.WriteRune(c)
		case c == '*' || c == '+' || c == '?' || c == '}':
			sb.WriteRune(c)
			if i+1 < n && rs[i+1] == '+' {
				unsupportedJavaRegex(regex, "possessive quantifier")
			}
		default:
			sb.WriteRune(c)
		}
	}
	if inClass {
		panic(NewPatternSyntaxException("Unclosed character class", regex, n-1))
	}
	return sb.String(), info
}

// Matcher is java.util.regex.Matcher.
type Matcher struct {
	p            *Pattern
	input        string
	groups       []int
	resumeAt     int
	prevWasEmpty bool
}

// Matcher is Pattern.matcher(CharSequence).
func (p *Pattern) Matcher(input string) *Matcher {
	return &Matcher{p: p, input: input}
}

// Reset is Matcher.reset(): no current match, and the next find() looks
// from the start again.
func (m *Matcher) Reset() *Matcher {
	m.groups, m.resumeAt, m.prevWasEmpty = nil, 0, false
	return m
}

func (m *Matcher) record(loc []int) bool {
	m.groups = loc
	if loc == nil {
		m.prevWasEmpty = false
		return false
	}
	m.resumeAt = loc[1]
	m.prevWasEmpty = loc[0] == loc[1]
	return true
}

// Matches is Matcher.matches(): the entire input must match.
func (m *Matcher) Matches() bool {
	return m.record(m.p.full.FindStringSubmatchIndex(m.input))
}

// LookingAt is Matcher.lookingAt(): a match must start at the beginning of
// the input but need not reach its end.
func (m *Matcher) LookingAt() bool {
	loc := m.p.findFrom(m.input, 0)
	if loc != nil && loc[0] != 0 {
		loc = nil
	}
	return m.record(loc)
}

// Find is Matcher.find(): the next match of the pattern.
func (m *Matcher) Find() bool {
	from := m.resumeAt
	if m.prevWasEmpty {
		if from == len(m.input) {
			m.groups = nil
			return false
		}
		_, size := decodeWTF8(m.input[from:])
		if size == 4 {
			unsupportedJavaRegex(m.p.pattern, "find() resuming inside a surrogate pair after an empty match")
		}
		from += size
	}
	return m.record(m.p.findFrom(m.input, from))
}

// groupCount is Matcher.groupCount().
func (m *Matcher) groupCount() int { return m.p.groups }

func (m *Matcher) requireMatch() {
	if m.groups == nil {
		panic(NewIllegalStateExceptionMsg("No match found"))
	}
}

func (m *Matcher) requireGroupIndex(group int) {
	if group < 0 || group > m.groupCount() {
		panic(NewIndexOutOfBoundsException("No group " + itoa(group)))
	}
}

// Start is Matcher.start(): the UTF-16 index of the start of the match.
func (m *Matcher) Start() int {
	m.requireMatch()
	return StringLength(m.input[:m.groups[0]])
}

// End is Matcher.end(): the UTF-16 index after the end of the match.
func (m *Matcher) End() int {
	m.requireMatch()
	return StringLength(m.input[:m.groups[1]])
}

// startGroup is Matcher.start(int): -1 when the group did not match.
func (m *Matcher) startGroup(group int) int {
	m.requireMatch()
	m.requireGroupIndex(group)
	if m.groups[2*group] < 0 {
		return -1
	}
	return StringLength(m.input[:m.groups[2*group]])
}

// endGroup is Matcher.end(int): -1 when the group did not match.
func (m *Matcher) endGroup(group int) int {
	m.requireMatch()
	m.requireGroupIndex(group)
	if m.groups[2*group+1] < 0 {
		return -1
	}
	return StringLength(m.input[:m.groups[2*group+1]])
}

// Group is Matcher.group(int): nil when the group did not take part in the
// match.
func (m *Matcher) Group(group int) *string {
	m.requireMatch()
	m.requireGroupIndex(group)
	if m.groups[2*group] < 0 {
		return nil
	}
	s := m.input[m.groups[2*group]:m.groups[2*group+1]]
	return &s
}

func (m *Matcher) appendReplacement(sb *strings.Builder, appendPos int, replacement string) {
	m.requireMatch()
	sb.WriteString(m.input[appendPos:m.groups[0]])
	r := replacement
	for i := 0; i < len(r); {
		switch r[i] {
		case '\\':
			i++
			if i == len(r) {
				panic(NewIllegalArgumentExceptionMsg("character to be escaped is missing"))
			}
			sb.WriteByte(r[i])
			i++
		case '$':
			i++
			if i == len(r) {
				panic(NewIllegalArgumentExceptionMsg("Illegal group reference: group index is missing"))
			}
			var group int
			switch c := r[i]; {
			case c == '{':
				i++
				start := i
				for i < len(r) && isASCIIAlphanumeric(r[i]) {
					i++
				}
				name := r[start:i]
				switch {
				case name == "":
					panic(NewIllegalArgumentExceptionMsg("named capturing group has 0 length name"))
				case i == len(r) || r[i] != '}':
					panic(NewIllegalArgumentExceptionMsg("named capturing group is missing trailing '}'"))
				case name[0] >= '0' && name[0] <= '9':
					panic(NewIllegalArgumentExceptionMsg("capturing group name {" + name + "} starts with digit character"))
				}
				panic(NewIllegalArgumentExceptionMsg("No group with name {" + name + "}"))
			case c >= '0' && c <= '9':
				group = int(c - '0')
				i++
				count := m.groupCount()
				for i < len(r) && r[i] >= '0' && r[i] <= '9' {
					longer := group*10 + int(r[i]-'0')
					if longer > count {
						break
					}
					group = longer
					i++
				}
				m.requireGroupIndex(group)
			default:
				panic(NewIllegalArgumentExceptionMsg("Illegal group reference"))
			}
			if lo := m.groups[2*group]; lo >= 0 {
				sb.WriteString(m.input[lo:m.groups[2*group+1]])
			}
		default:
			j := i + 1
			for j < len(r) && r[j] != '\\' && r[j] != '$' {
				j++
			}
			sb.WriteString(r[i:j])
			i = j
		}
	}
}

func isASCIIAlphanumeric(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// ReplaceFirst is Matcher.replaceFirst(String).
func (m *Matcher) ReplaceFirst(replacement string) string {
	m.Reset()
	if !m.Find() {
		return m.input
	}
	var sb strings.Builder
	m.appendReplacement(&sb, 0, replacement)
	sb.WriteString(m.input[m.groups[1]:])
	return sb.String()
}

// ReplaceAll is Matcher.replaceAll(String).
func (m *Matcher) ReplaceAll(replacement string) string {
	m.Reset()
	if !m.Find() {
		return m.input
	}
	var sb strings.Builder
	appendPos := 0
	for {
		m.appendReplacement(&sb, appendPos, replacement)
		appendPos = m.groups[1]
		if !m.Find() {
			break
		}
	}
	sb.WriteString(m.input[appendPos:])
	return sb.String()
}

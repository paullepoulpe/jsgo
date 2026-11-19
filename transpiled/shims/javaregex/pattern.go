// Clean-room specification based on java.util.regex.Pattern, Matcher, and PatternSyntaxException Javadoc.

package javaregex

import (
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Pattern flag constants (java.util.regex.Pattern Javadoc).
const (
	Pattern_UNIX_LINES       = 0x01
	Pattern_CASE_INSENSITIVE = 0x02
	Pattern_COMMENTS         = 0x04
	Pattern_MULTILINE        = 0x08
	Pattern_LITERAL          = 0x10
	Pattern_DOTALL           = 0x20
	Pattern_UNICODE_CASE     = 0x40
)

// PatternSyntaxException implements java.util.regex.PatternSyntaxException.
type PatternSyntaxException struct {
	javalang.IllegalArgumentException
	desc    string
	pattern string
	index   int
}

// NewPatternSyntaxException creates a PatternSyntaxException with Java-compatible message layout.
func NewPatternSyntaxException(desc, regex string, index int) *PatternSyntaxException {
	e := &PatternSyntaxException{desc: desc, pattern: regex, index: index}
	var sb strings.Builder
	sb.WriteString(desc)
	if index >= 0 {
		sb.WriteString(" near index ")
		sb.WriteString(strconv.Itoa(index))
	}
	sb.WriteString("\n")
	sb.WriteString(regex)
	if index >= 0 && index < javalang.StringLength(regex) {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat(" ", index))
		sb.WriteString("^")
	}
	javalang.InitThrowable(e, "java.util.regex.PatternSyntaxException", sb.String(), nil)
	return e
}

func (e *PatternSyntaxException) GetDescription() string { return e.desc }
func (e *PatternSyntaxException) GetIndex() int          { return e.index }
func (e *PatternSyntaxException) GetPattern() string     { return e.pattern }

// Pattern represents a compiled java.util.regex.Pattern (and com.google.re2j.Pattern).
type Pattern struct {
	pattern                 string
	flags                   int
	groups                  int
	namedGroups             map[string]int
	full                    *regexp.Regexp
	searchers               [2][2]searcher
	fullPrefilter           literalPrefilter
	findPrefilter           literalPrefilter
	matchesCore             *regexp.Regexp
	matchesCoreMode         uint8
	matchesCoreCanMatchTerm bool
}

type searcher struct {
	re    *regexp.Regexp
	marks []int
}

type javaRegexMode struct {
	anchored bool
	resumed  bool
	crlf     bool
}

type javaRegexInfo struct {
	groups      int
	namedGroups map[string]int
	marks       []int
	caret       bool
	dollar      bool
}

const (
	javaDotClass             = `[^\n\r\x{85}\x{2028}\x{2029}]`
	javaSpaceChars           = `\t\n\x0B\f\r `
	javaVerticalSpaceChars   = `\n\x0B\f\r\x{85}\x{2028}\x{2029}`
	javaHorizontalSpaceChars = ` \t\x{A0}\x{1680}\x{180E}\x{2000}-\x{200A}\x{202F}\x{205F}\x{3000}`
	javaFinalTerminators     = `[\n\r\x{85}\x{2028}\x{2029}]`
	javaFinalCRLF            = `\r\n`
	javaNoMatch              = `[^\x00-\x{10FFFF}]`
	patternCacheLimit        = 4096
)

type cacheKey struct {
	regex string
	flags int
}

var (
	patternCache     sync.Map // cacheKey -> *Pattern
	patternCacheSize atomic.Int64
)

// Pattern_Compile implements Pattern.compile(String).
func Pattern_Compile(regex string) *Pattern {
	return Pattern_CompileFlags(regex, 0)
}

// Pattern_CompileFlags implements Pattern.compile(String, int).
func Pattern_CompileFlags(regex string, flags int) *Pattern {
	key := cacheKey{regex: regex, flags: flags}
	if cached, ok := patternCache.Load(key); ok {
		p := *cached.(*Pattern)
		return &p
	}
	p := patternCompile(regex, flags)
	if patternCacheSize.Load() < patternCacheLimit {
		if patternCacheSize.Add(1) <= patternCacheLimit {
			cached := *p
			if _, loaded := patternCache.LoadOrStore(key, &cached); loaded {
				patternCacheSize.Add(-1)
			}
		} else {
			patternCacheSize.Add(-1)
		}
	}
	return p
}

// Pattern_Matches implements Pattern.matches(String, CharSequence).
func Pattern_Matches(regex, input string) bool {
	return Pattern_Compile(regex).Matches(input)
}

// Pattern_Quote implements Pattern.quote(String).
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

// Matcher_QuoteReplacement implements Matcher.quoteReplacement(String).
func Matcher_QuoteReplacement(s string) string {
	if !strings.ContainsAny(s, `\$`) {
		return s
	}
	var sb strings.Builder
	sb.Grow(len(s) + 8)
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' || s[i] == '$' {
			sb.WriteByte('\\')
		}
		sb.WriteByte(s[i])
	}
	return sb.String()
}

func patternCompile(regex string, flags int) *Pattern {
	effective := regex
	if flags&Pattern_LITERAL != 0 {
		effective = Pattern_Quote(regex)
	}
	var flagPrefix strings.Builder
	if flags&Pattern_CASE_INSENSITIVE != 0 {
		flagPrefix.WriteByte('i')
	}
	if flags&Pattern_MULTILINE != 0 {
		flagPrefix.WriteByte('m')
	}
	if flags&Pattern_DOTALL != 0 {
		flagPrefix.WriteByte('s')
	}
	if flagPrefix.Len() > 0 {
		effective = "(?" + flagPrefix.String() + ")" + effective
	}

	whole, info := translateJavaRegex(effective, javaRegexMode{anchored: true}, regex)
	fullRe := `\A(?:` + whole + `)\z`
	p := &Pattern{
		pattern:     regex,
		flags:       flags,
		groups:      info.groups,
		namedGroups: info.namedGroups,
		full:        compileTranslatedJavaRegex(regex, fullRe),
	}
	var search00 string
	for resumed := range 2 {
		for crlf := range 2 {
			switch {
			case resumed == 1 && !info.caret:
				p.searchers[1][crlf] = p.searchers[0][crlf]
			case crlf == 1 && !info.dollar:
				p.searchers[resumed][1] = p.searchers[resumed][0]
			default:
				t, ti := translateJavaRegex(effective, javaRegexMode{resumed: resumed == 1, crlf: crlf == 1}, regex)
				if resumed == 0 && crlf == 0 {
					search00 = t
				}
				p.searchers[resumed][crlf] = searcher{compileTranslatedJavaRegex(regex, t), ti.marks}
			}
		}
	}
	populatePatternFastPaths(p, whole, search00)
	return p
}

func compileTranslatedJavaRegex(regex, translated string) *regexp.Regexp {
	re, err := regexp.Compile(translated)
	if err == nil {
		return re
	}
	if se, ok := err.(*syntax.Error); ok {
		switch se.Code {
		case syntax.ErrInvalidRepeatSize, syntax.ErrInvalidRepeatOp, syntax.ErrLarge,
			syntax.ErrNestingDepth, syntax.ErrInternalError:
			unsupportedJavaRegex(regex, "RE2: "+se.Code.String())
		}
		panic(NewPatternSyntaxException(se.Code.String(), regex, -1))
	}
	panic(err)
}

func unsupportedJavaRegex(regex, what string) {
	panic(javalang.NewUnsupportedOperationException("Java regex construct not supported (" + what + "): " + regex))
}

func (p *Pattern) Pattern() string { return p.pattern }
func (p *Pattern) Flags() int      { return p.flags }
func (p *Pattern) String() string  { return p.pattern }

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

func (p *Pattern) findFrom(input string, from int) []int {
	tail := input[from:]
	if !p.findPrefilter.match(tail) {
		return nil
	}
	s := p.searcherFor(input, from)
	loc := s.re.FindStringSubmatchIndex(tail)
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

// Matches implements p.matcher(s).matches().
func (p *Pattern) Matches(s string) bool {
	if !p.fullPrefilter.match(s) {
		return false
	}
	switch p.matchesCoreMode {
	case matchesCoreDotAll:
		if p.matchesCore == nil {
			return true
		}
		return p.matchesCore.MatchString(s)
	case matchesCoreJavaDot:
		if !containsJavaLineTerminator(s) {
			if p.matchesCore == nil {
				return true
			}
			return p.matchesCore.MatchString(s)
		}
		if !p.matchesCoreCanMatchTerm {
			return false
		}
	}
	return p.full.MatchString(s)
}

// Find reports whether any subsequence of s matches p.
func (p *Pattern) Find(s string) bool {
	if !p.findPrefilter.match(s) {
		return false
	}
	return p.searcherFor(s, 0).re.MatchString(s)
}

// Split implements Pattern.split(CharSequence) (limit = 0: trims trailing empty strings).
func (p *Pattern) Split(input string) []string {
	return p.SplitLimit(input, 0)
}

// SplitLimit implements Pattern.split(CharSequence, int) per Java Pattern.split Javadoc.
func (p *Pattern) SplitLimit(input string, limit int) []string {
	m := p.Matcher(input)
	var matchList []string
	index := 0
	matchLimited := limit > 0
	matchedAny := false

	for m.Find() {
		matchedAny = true
		if !matchLimited || len(matchList) < limit-1 {
			startByte, endByte := m.groups[0], m.groups[1]
			if index == 0 && index == startByte && startByte == endByte {
				// Java 8+: leading empty match does not produce a leading empty substring.
				continue
			}
			matchList = append(matchList, input[index:startByte])
			index = endByte
		} else if len(matchList) == limit-1 {
			matchList = append(matchList, input[index:])
			index = len(input)
			break
		}
	}

	if !matchedAny {
		return []string{input}
	}
	if !matchLimited || len(matchList) < limit {
		matchList = append(matchList, input[index:])
	}
	resultSize := len(matchList)
	if limit == 0 {
		for resultSize > 0 && matchList[resultSize-1] == "" {
			resultSize--
		}
	}
	return matchList[:resultSize]
}

func javaPosixPropertyChars(name string) (string, bool) {
	switch name {
	case "Lower":
		return `a-z`, true
	case "Upper":
		return `A-Z`, true
	case "ASCII":
		return `\x00-\x7F`, true
	case "Alpha":
		return `A-Za-z`, true
	case "Digit":
		return `0-9`, true
	case "Alnum":
		return `A-Za-z0-9`, true
	case "Punct":
		return `\x21-\x2F\x3A-\x40\x5B-\x60\x7B-\x7E`, true
	case "Graph":
		return `\x21-\x7E`, true
	case "Print":
		return `\x20-\x7E`, true
	case "Blank":
		return ` \t`, true
	case "Cntrl":
		return `\x00-\x1F\x7F`, true
	case "XDigit":
		return `0-9a-fA-F`, true
	case "Space":
		return ` \t\n\x0B\f\r`, true
	default:
		return "", false
	}
}

func translateJavaRegex(regex string, mode javaRegexMode, origRegex string) (string, javaRegexInfo) {
	var sb strings.Builder
	var info javaRegexInfo
	multiline := false
	dotall := false
	rs := []rune(regex)
	n := len(rs)
	inClass := false
	classStart := 0
	depth, lastGroup := 0, 0
	tailDepth := -1
	var endsInDollar []bool
	afterGroupClose := func(i int) {
		if i+1 < n && strings.ContainsRune("*+?{", rs[i+1]) {
			unsupportedJavaRegex(origRegex, "quantified group containing $")
		}
	}
	for i := 0; i < n; i++ {
		c := rs[i]
		if tailDepth >= 0 && !inClass && javaTokenCanMatchTerminator(rs, i, dotall) {
			unsupportedJavaRegex(origRegex, "$ followed by potential line terminator")
		}
		switch {
		case c == '\\':
			if i+1 >= n {
				panic(NewPatternSyntaxException("Unexpected internal error", origRegex, len([]rune(origRegex))))
			}
			i++
			e := rs[i]
			switch {
			case e == 'Q':
				i = writeJavaQuotedLiteral(&sb, rs, i, inClass)
			case e == 'E':
			case e == 's':
				if inClass {
					sb.WriteString(javaSpaceChars)
				} else {
					sb.WriteString(`[` + javaSpaceChars + `]`)
				}
			case e == 'S':
				if inClass {
					unsupportedJavaRegex(origRegex, `\S in class`)
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
					unsupportedJavaRegex(origRegex, `\`+string(e)+` in class`)
				}
				chars := javaVerticalSpaceChars
				if e == 'H' {
					chars = javaHorizontalSpaceChars
				}
				sb.WriteString(`[^` + chars + `]`)
			case e == 'd' || e == 'D' || e == 'w' || e == 'W' || e == 't' || e == 'n' || e == 'r' || e == 'f':
				sb.WriteByte('\\')
				sb.WriteRune(e)
			case (e == 'b' || e == 'B') && !inClass:
				sb.WriteByte('\\')
				sb.WriteRune(e)
			case e == 'e':
				sb.WriteString(`\x1B`)
			case e == 'a':
				sb.WriteString(`\x07`)
			case e == '0' || (e >= '1' && e <= '9') || e == 'k':
				unsupportedJavaRegex(origRegex, "back-reference or octal escape")
			case e == 'A' || e == 'z':
				if inClass {
					unsupportedJavaRegex(origRegex, `\`+string(e)+` in class`)
				}
				sb.WriteByte('\\')
				sb.WriteRune(e)
			case e == 'p' || e == 'P':
				if i+1 < n && rs[i+1] == '{' {
					closeIdx := -1
					for k := i + 2; k < n; k++ {
						if rs[k] == '}' {
							closeIdx = k
							break
						}
					}
					if closeIdx > i+2 {
						propName := string(rs[i+2 : closeIdx])
						if posixChars, ok := javaPosixPropertyChars(propName); ok {
							if inClass {
								if e == 'P' {
									unsupportedJavaRegex(origRegex, `\P{`+propName+`} in class`)
								}
								sb.WriteString(posixChars)
							} else if e == 'p' {
								sb.WriteString(`[` + posixChars + `]`)
							} else {
								sb.WriteString(`[^` + posixChars + `]`)
							}
						} else {
							sb.WriteByte('\\')
							sb.WriteRune(e)
							sb.WriteString(string(rs[i+1 : closeIdx+1]))
						}
						i = closeIdx
						break
					}
				} else if i+1 < n {
					sb.WriteByte('\\')
					sb.WriteRune(e)
					sb.WriteRune(rs[i+1])
					i++
					break
				}
				panic(NewPatternSyntaxException("Illegal character property", origRegex, i))
			case e == 'x':
				if i+1 < n && rs[i+1] == '{' {
					closeIdx := -1
					for k := i + 2; k < n; k++ {
						if rs[k] == '}' {
							closeIdx = k
							break
						}
					}
					if closeIdx > i+2 {
						sb.WriteString(`\x`)
						sb.WriteString(string(rs[i+1 : closeIdx+1]))
						i = closeIdx
						break
					}
				} else if i+2 < n {
					sb.WriteString(`\x`)
					sb.WriteRune(rs[i+1])
					sb.WriteRune(rs[i+2])
					i += 2
					break
				}
				panic(NewPatternSyntaxException("Illegal hexadecimal escape sequence", origRegex, i))
			case e == 'u':
				if i+4 < n {
					sb.WriteString(`\x{`)
					sb.WriteString(string(rs[i+1 : i+5]))
					sb.WriteByte('}')
					i += 4
					break
				}
				panic(NewPatternSyntaxException("Illegal Unicode escape sequence", origRegex, i))
			case e == 'Z' || e == 'G' || e == 'R' || e == 'X' || e == 'c' || e == 'N':
				unsupportedJavaRegex(origRegex, `\`+string(e))
			case (e >= 'a' && e <= 'z') || (e >= 'A' && e <= 'Z'):
				panic(NewPatternSyntaxException("Illegal/unsupported escape sequence", origRegex, i))
			default:
				sb.WriteString(regexp.QuoteMeta(string(e)))
			}
		case inClass:
			switch c {
			case '[':
				unsupportedJavaRegex(origRegex, "nested character class")
			case '&':
				if i+1 < n && rs[i+1] == '&' {
					unsupportedJavaRegex(origRegex, "class intersection")
				}
				sb.WriteRune(c)
			case ']':
				if isLeadingClassClose(rs, i, classStart) {
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
				if i+1 < n && strings.ContainsRune("*+?{", rs[i+1]) {
					unsupportedJavaRegex(origRegex, "quantified $")
				}
				tailDepth = depth
				if mode.anchored {
					sb.WriteString(`\z`)
				} else {
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
			if tailDepth == depth || (depth < len(endsInDollar) && endsInDollar[depth]) {
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
				if i+2 < n && rs[i+2] == '<' && (i+3 >= n || (rs[i+3] != '=' && rs[i+3] != '!')) {
					closeIdx := -1
					for k := i + 3; k < n; k++ {
						if rs[k] == '>' {
							closeIdx = k
							break
						}
					}
					if closeIdx <= i+3 {
						panic(NewPatternSyntaxException("Capturing group is missing trailing '>'", origRegex, i+2))
					}
					gname := string(rs[i+3 : closeIdx])
					lastGroup++
					info.groups++
					if info.namedGroups == nil {
						info.namedGroups = make(map[string]int)
					}
					if _, dup := info.namedGroups[gname]; dup {
						panic(NewPatternSyntaxException("Named capturing group <"+gname+"> is already defined", origRegex, i+3))
					}
					info.namedGroups[gname] = info.groups
					sb.WriteString("(?P<")
					sb.WriteString(gname)
					sb.WriteByte('>')
					i = closeIdx
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
				unsupportedJavaRegex(origRegex, "special group")
			}
			lastGroup++
			info.groups++
			sb.WriteRune(c)
		case c == '*' || c == '+' || c == '?' || c == '}':
			sb.WriteRune(c)
			if i+1 < n && rs[i+1] == '+' {
				unsupportedJavaRegex(origRegex, "possessive quantifier")
			}
		default:
			sb.WriteRune(c)
		}
	}
	if inClass {
		panic(NewPatternSyntaxException("Unclosed character class", origRegex, len([]rune(origRegex))-1))
	}
	return sb.String(), info
}

func writeJavaQuotedLiteral(sb *strings.Builder, rs []rune, i int, inClass bool) int {
	n := len(rs)
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
	return i
}

func isLeadingClassClose(rs []rune, i, classStart int) bool {
	return i == classStart || (i == classStart+1 && rs[classStart] == '^')
}

func isJavaLineTerminatorRune(r rune) bool {
	return r == '\n' || r == '\r' || r == '\u0085' || r == '\u2028' || r == '\u2029'
}

func javaTokenCanMatchTerminator(rs []rune, i int, dotall bool) bool {
	c := rs[i]
	switch {
	case isJavaLineTerminatorRune(c):
		return true
	case c == '.':
		return dotall
	case c == '[':
		return true
	case c == '\\':
		if i+1 >= len(rs) {
			return false
		}
		e := rs[i+1]
		switch e {
		case 'n', 'r', 's', 'S', 'v', 'D', 'W', 'H':
			return true
		case 'Q':
			for j := i + 2; j < len(rs); j++ {
				if j+1 < len(rs) && rs[j] == '\\' && rs[j+1] == 'E' {
					return false
				}
				if isJavaLineTerminatorRune(rs[j]) {
					return true
				}
			}
			return false
		default:
			return isJavaLineTerminatorRune(e)
		}
	}
	return false
}

// IsValidPattern reports whether pattern is syntactically valid as a Java Pattern.
func IsValidPattern(pattern string) bool {
	var sb strings.Builder
	rs := []rune(pattern)
	n := len(rs)
	inClass := false
	classStart := 0
	for i := 0; i < n; i++ {
		c := rs[i]
		switch {
		case c == '\\':
			if i+1 >= n {
				return false
			}
			i++
			e := rs[i]
			switch {
			case e == 'Q':
				i = writeJavaQuotedLiteral(&sb, rs, i, inClass)
			case e == 'E':
			case e == 's', e == 'S', e == 'v', e == 'h', e == 'V', e == 'H',
				e == 'e', e == 'a',
				e == '0', e >= '1' && e <= '9', e == 'k',
				e == 'b', e == 'B', e == 'z', e == 'Z', e == 'G', e == 'A', e == 'R', e == 'X',
				e == 'p', e == 'P', e == 'x', e == 'u', e == 'c', e == 'N':
				sb.WriteByte('a')
			case e == 'd', e == 'D', e == 'w', e == 'W', e == 't', e == 'n', e == 'r', e == 'f':
				sb.WriteByte('\\')
				sb.WriteRune(e)
			case (e >= 'a' && e <= 'z') || (e >= 'A' && e <= 'Z'):
				return false
			default:
				sb.WriteString(regexp.QuoteMeta(string(e)))
			}
		case inClass:
			if c == ']' {
				if isLeadingClassClose(rs, i, classStart) {
					sb.WriteString(`\]`)
				} else {
					inClass = false
					sb.WriteByte(']')
				}
			} else {
				sb.WriteRune(c)
			}
		case c == '[':
			inClass = true
			classStart = i + 1
			sb.WriteByte('[')
		default:
			sb.WriteRune(c)
		}
	}
	if inClass {
		return false
	}
	if _, err := syntax.Parse(sb.String(), syntax.Perl); err != nil {
		if se, ok := err.(*syntax.Error); ok {
			switch se.Code {
			case syntax.ErrInvalidRepeatSize, syntax.ErrInvalidRepeatOp, syntax.ErrLarge,
				syntax.ErrNestingDepth, syntax.ErrInternalError:
				return true
			}
		}
		return false
	}
	return true
}

// ---------------------------------------------------------------------------
// Matcher
// ---------------------------------------------------------------------------

// Matcher implements java.util.regex.Matcher.
type Matcher struct {
	p                *Pattern
	input            string
	groups           []int
	resumeAt         int
	appendPos        int
	prevWasEmpty     bool
	hasRegion        bool
	regionStartUTF16 int
	regionEndUTF16   int
	regionStartByte  int
	regionEndByte    int
	hitEnd           bool
}

// Matcher creates a Matcher for input against p.
func (p *Pattern) Matcher(input string) *Matcher {
	return &Matcher{p: p, input: input}
}

// Reset resets this matcher's match state and region.
func (m *Matcher) Reset() *Matcher {
	m.groups, m.resumeAt, m.appendPos, m.prevWasEmpty = nil, 0, 0, false
	m.hasRegion, m.regionStartUTF16, m.regionEndUTF16 = false, 0, 0
	m.regionStartByte, m.regionEndByte, m.hitEnd = 0, 0, false
	return m
}

// ResetInput resets this matcher with a new input string.
func (m *Matcher) ResetInput(input string) *Matcher {
	m.input = input
	return m.Reset()
}

// Region sets the limits of this matcher's region (in UTF-16 code units).
func (m *Matcher) Region(startUTF16, endUTF16 int) *Matcher {
	inLen := javalang.StringLength(m.input)
	if startUTF16 < 0 || endUTF16 < 0 || startUTF16 > inLen || endUTF16 > inLen || startUTF16 > endUTF16 {
		panic(javalang.NewIndexOutOfBoundsException("Illegal region"))
	}
	m.Reset()
	m.hasRegion = true
	m.regionStartUTF16 = startUTF16
	m.regionEndUTF16 = endUTF16
	m.regionStartByte = utf16IndexToByteOffset(m.input, startUTF16)
	m.regionEndByte = utf16IndexToByteOffset(m.input, endUTF16)
	m.resumeAt = m.regionStartByte
	return m
}

// RegionStart returns the start index (in UTF-16 code units) of this matcher's region.
func (m *Matcher) RegionStart() int {
	if m.hasRegion {
		return m.regionStartUTF16
	}
	return 0
}

// RegionEnd returns the end index (in UTF-16 code units) of this matcher's region.
func (m *Matcher) RegionEnd() int {
	if m.hasRegion {
		return m.regionEndUTF16
	}
	return javalang.StringLength(m.input)
}

// HitEnd reports whether the end of the region was hit in the last match operation.
func (m *Matcher) HitEnd() bool {
	return m.hitEnd
}

func (m *Matcher) effBounds() (int, int) {
	if m.hasRegion {
		return m.regionStartByte, m.regionEndByte
	}
	return 0, len(m.input)
}

func adjustLoc(loc []int, delta int) []int {
	if loc == nil || delta == 0 {
		return loc
	}
	for i, off := range loc {
		if off >= 0 {
			loc[i] = off + delta
		}
	}
	return loc
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

// Matches implements Matcher.matches().
func (m *Matcher) Matches() bool {
	rStart, rEnd := m.effBounds()
	sub := m.input[rStart:rEnd]
	if !m.p.fullPrefilter.match(sub) {
		m.hitEnd = true
		return m.record(nil)
	}
	loc := m.p.full.FindStringSubmatchIndex(sub)
	if loc == nil {
		m.hitEnd = true
		return m.record(nil)
	}
	loc = adjustLoc(loc, rStart)
	m.hitEnd = loc[1] == rEnd
	return m.record(loc)
}

// LookingAt implements Matcher.lookingAt().
func (m *Matcher) LookingAt() bool {
	rStart, rEnd := m.effBounds()
	sub := m.input[rStart:rEnd]
	loc := m.p.findFrom(sub, 0)
	if loc != nil && loc[0] != 0 {
		loc = nil
	}
	if loc == nil {
		m.hitEnd = len(sub) == 0
		return m.record(nil)
	}
	loc = adjustLoc(loc, rStart)
	m.hitEnd = loc[1] == rEnd
	return m.record(loc)
}

// Find implements Matcher.find().
func (m *Matcher) Find() bool {
	rStart, rEnd := m.effBounds()
	from := m.resumeAt
	if from < rStart {
		from = rStart
	}
	if m.prevWasEmpty {
		if from >= rEnd {
			m.groups = nil
			m.hitEnd = true
			return false
		}
		// Advance by one WTF-8 code point (or 1 UTF-16 unit if BMP/ASCII).
		b := m.input[from]
		switch {
		case b < 0x80:
			from++
		case b < 0xE0:
			from += 2
		case b < 0xF0:
			from += 3
		default:
			unsupportedJavaRegex(m.p.pattern, "find() resuming inside a surrogate pair after an empty match")
		}
	}
	if from > rEnd {
		m.groups = nil
		m.hitEnd = true
		return false
	}
	sub := m.input[rStart:rEnd]
	loc := m.p.findFrom(sub, from-rStart)
	if loc == nil {
		m.hitEnd = true
		return m.record(nil)
	}
	loc = adjustLoc(loc, rStart)
	m.hitEnd = loc[1] == rEnd
	return m.record(loc)
}

// FindFrom implements Matcher.find(int startUTF16).
func (m *Matcher) FindFrom(startUTF16 int) bool {
	inLen := javalang.StringLength(m.input)
	if startUTF16 < 0 || startUTF16 > inLen {
		panic(javalang.NewIndexOutOfBoundsException("Illegal start index"))
	}
	m.Reset()
	byteOff := utf16IndexToByteOffset(m.input, startUTF16)
	return m.record(m.p.findFrom(m.input, byteOff))
}

func utf16IndexToByteOffset(s string, targetUTF16 int) int {
	if targetUTF16 == 0 || javalang.IsASCII(s) {
		return targetUTF16
	}
	uIdx := 0
	for i := 0; i < len(s); {
		if uIdx == targetUTF16 {
			return i
		}
		b := s[i]
		switch {
		case b < 0x80:
			i++
			uIdx++
		case b < 0xE0:
			i += 2
			uIdx++
		case b < 0xF0:
			i += 3
			uIdx++
		default:
			i += 4
			uIdx += 2
		}
	}
	return len(s)
}

// GroupCount implements Matcher.groupCount().
func (m *Matcher) GroupCount() int { return m.p.groups }

func (m *Matcher) requireMatch() {
	if m.groups == nil {
		panic(javalang.NewIllegalStateException("No match found"))
	}
}

func (m *Matcher) requireGroupIndex(group int) {
	if group < 0 || group > m.GroupCount() {
		panic(javalang.NewIndexOutOfBoundsException("No group " + strconv.Itoa(group)))
	}
}

// Start implements Matcher.start().
func (m *Matcher) Start() int {
	m.requireMatch()
	return javalang.StringLength(m.input[:m.groups[0]])
}

// End implements Matcher.end().
func (m *Matcher) End() int {
	m.requireMatch()
	return javalang.StringLength(m.input[:m.groups[1]])
}

// StartGroup implements Matcher.start(int).
func (m *Matcher) StartGroup(group int) int {
	m.requireMatch()
	m.requireGroupIndex(group)
	if m.groups[2*group] < 0 {
		return -1
	}
	return javalang.StringLength(m.input[:m.groups[2*group]])
}

// EndGroup implements Matcher.end(int).
func (m *Matcher) EndGroup(group int) int {
	m.requireMatch()
	m.requireGroupIndex(group)
	if m.groups[2*group+1] < 0 {
		return -1
	}
	return javalang.StringLength(m.input[:m.groups[2*group+1]])
}

// Group0 implements Matcher.group() (entire match).
func (m *Matcher) Group0() string {
	s, _ := m.Group(0)
	return s
}

// Group implements Matcher.group(int), returning (subsequence, matched).
func (m *Matcher) Group(group int) (string, bool) {
	m.requireMatch()
	m.requireGroupIndex(group)
	if m.groups[2*group] < 0 {
		return "", false
	}
	return m.input[m.groups[2*group]:m.groups[2*group+1]], true
}

// GroupPtr implements Matcher.group(int) returning *string (nil when unmatched).
func (m *Matcher) GroupPtr(group int) *string {
	s, ok := m.Group(group)
	if !ok {
		return nil
	}
	return &s
}

// GroupNamed implements Matcher.group(String), returning (subsequence, matched).
func (m *Matcher) GroupNamed(name string) (string, bool) {
	m.requireMatch()
	idx, ok := m.p.namedGroups[name]
	if !ok {
		panic(javalang.NewIllegalArgumentException("No group with name <" + name + ">"))
	}
	return m.Group(idx)
}

func (m *Matcher) appendReplacementAt(sb *strings.Builder, appendPos int, replacement string) {
	m.requireMatch()
	sb.WriteString(m.input[appendPos:m.groups[0]])
	r := replacement
	for i := 0; i < len(r); {
		switch r[i] {
		case '\\':
			i++
			if i == len(r) {
				panic(javalang.NewIllegalArgumentException("character to be escaped is missing"))
			}
			sb.WriteByte(r[i])
			i++
		case '$':
			i++
			if i == len(r) {
				panic(javalang.NewIllegalArgumentException("Illegal group reference: group index is missing"))
			}
			c := r[i]
			if c == '{' {
				i++
				start := i
				for i < len(r) && isASCIIAlphanumeric(r[i]) {
					i++
				}
				name := r[start:i]
				switch {
				case name == "":
					panic(javalang.NewIllegalArgumentException("named capturing group has 0 length name"))
				case i == len(r) || r[i] != '}':
					panic(javalang.NewIllegalArgumentException("named capturing group is missing trailing '}'"))
				case name[0] >= '0' && name[0] <= '9':
					panic(javalang.NewIllegalArgumentException("capturing group name {" + name + "} starts with digit character"))
				default:
					idx, ok := m.p.namedGroups[name]
					if !ok {
						panic(javalang.NewIllegalArgumentException("No group with name {" + name + "}"))
					}
					i++ // consume '}'
					if lo := m.groups[2*idx]; lo >= 0 {
						sb.WriteString(m.input[lo:m.groups[2*idx+1]])
					}
					continue
				}
			}
			if c < '0' || c > '9' {
				panic(javalang.NewIllegalArgumentException("Illegal group reference"))
			}
			group := int(c - '0')
			i++
			count := m.GroupCount()
			for i < len(r) && r[i] >= '0' && r[i] <= '9' {
				longer := group*10 + int(r[i]-'0')
				if longer > count {
					break
				}
				group = longer
				i++
			}
			m.requireGroupIndex(group)
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
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// AppendReplacement implements Matcher.appendReplacement(StringBuilder, String).
func (m *Matcher) AppendReplacement(sb *strings.Builder, replacement string) *Matcher {
	m.appendReplacementAt(sb, m.appendPos, replacement)
	m.appendPos = m.groups[1]
	return m
}

// AppendTail implements Matcher.appendTail(StringBuilder).
func (m *Matcher) AppendTail(sb *strings.Builder) *strings.Builder {
	sb.WriteString(m.input[m.appendPos:])
	return sb
}

// ReplaceFirst implements Matcher.replaceFirst(String).
func (m *Matcher) ReplaceFirst(replacement string) string {
	m.Reset()
	if !m.Find() {
		return m.input
	}
	var sb strings.Builder
	m.appendReplacementAt(&sb, 0, replacement)
	sb.WriteString(m.input[m.groups[1]:])
	return sb.String()
}

// ReplaceAll implements Matcher.replaceAll(String).
func (m *Matcher) ReplaceAll(replacement string) string {
	m.Reset()
	if !m.Find() {
		return m.input
	}
	var sb strings.Builder
	appendPos := 0
	for {
		m.appendReplacementAt(&sb, appendPos, replacement)
		appendPos = m.groups[1]
		if !m.Find() {
			break
		}
	}
	sb.WriteString(m.input[appendPos:])
	return sb.String()
}

// ---------------------------------------------------------------------------
// String regex convenience helpers (String.matches / replaceAll / replaceFirst / split)
// ---------------------------------------------------------------------------

func StringMatches(s, regex string) bool {
	return Pattern_Compile(regex).Matches(s)
}

func StringReplaceAll(s, regex, replacement string) string {
	return Pattern_Compile(regex).Matcher(s).ReplaceAll(replacement)
}

func StringReplaceFirst(s, regex, replacement string) string {
	return Pattern_Compile(regex).Matcher(s).ReplaceFirst(replacement)
}

func StringSplit(s, regex string) []string {
	return Pattern_Compile(regex).SplitLimit(s, 0)
}

func StringSplitLimit(s, regex string, limit int) []string {
	return Pattern_Compile(regex).SplitLimit(s, limit)
}

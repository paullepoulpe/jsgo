package jutil

// Go-only: language-adapter — java.util.regex.Pattern on top of Go's RE2.
//
// Java regexes are translated to RE2 syntax so that the constructs they have
// in common keep Java's meaning:
//   - `.` excludes Java's line terminators (\n, \r, \u0085,  ,  ),
//     not only \n as in RE2;
//   - \s, \S (and \s inside classes) are Java's [ \t\n\x0B\f\r];
//   - \Q...\E quoting is resolved here;
//   - Matcher.matches() anchors the whole input.
//
// Constructs whose Java semantics RE2 cannot reproduce (back-references,
// lookaround, possessive quantifiers, atomic groups, class intersections,
// nested classes, \b \B \z \Z \G, `$`, embedded flags) make Pattern_Compile
// panic with UnsupportedOperationException rather than silently match
// differently.

import (
	"regexp"
	"strings"
)

// Pattern is a compiled java.util.regex.Pattern.
type Pattern struct {
	pattern string
	re      *regexp.Regexp // unanchored (find)
	full    *regexp.Regexp // anchored at both ends (matches)
}

const javaDotClass = `[^\n\r\x{85}\x{2028}\x{2029}]`
const javaSpaceChars = `\t\n\x0B\f\r `

// Pattern_Compile is Pattern.compile(String).
func Pattern_Compile(regex string) *Pattern {
	t := translateJavaRegex(regex)
	return &Pattern{
		pattern: regex,
		re:      regexp.MustCompile(t),
		full:    regexp.MustCompile(`^(?:` + t + `)$`),
	}
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
func (p *Pattern) Find(s string) bool { return p.re.MatchString(s) }

func unsupportedJavaRegex(regex, what string) {
	panic(NewUnsupportedOperationExceptionMsg("Java regex construct not supported (" + what + "): " + regex))
}

// translateJavaRegex rewrites a Java regex into an equivalent RE2 regex.
func translateJavaRegex(regex string) string {
	var sb strings.Builder
	rs := []rune(regex)
	n := len(rs)
	inClass := false
	classStart := 0
	for i := 0; i < n; i++ {
		c := rs[i]
		switch {
		case c == '\\':
			if i+1 >= n {
				panic(NewIllegalArgumentExceptionMsg("Unexpected internal error near index " + itoa(n) + "\n" + regex))
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
				e == 'h' || e == 'H' || e == 'v' || e == 'V' || e == 'p' || e == 'P' || e == 'x' || e == 'u' || e == 'c' || e == 'N':
				unsupportedJavaRegex(regex, `\`+string(e))
			case (e >= 'a' && e <= 'z') || (e >= 'A' && e <= 'Z'):
				panic(NewIllegalArgumentExceptionMsg("Illegal/unsupported escape sequence near index " + itoa(i) + "\n" + regex))
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
			sb.WriteString(javaDotClass)
		case c == '$':
			unsupportedJavaRegex(regex, "$")
		case c == '(':
			if i+1 < n && rs[i+1] == '?' {
				if i+2 < n && rs[i+2] == ':' {
					sb.WriteString("(?:")
					i += 2
					continue
				}
				unsupportedJavaRegex(regex, "special group")
			}
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
		panic(NewIllegalArgumentExceptionMsg("Unclosed character class near index " + itoa(n-1) + "\n" + regex))
	}
	return sb.String()
}

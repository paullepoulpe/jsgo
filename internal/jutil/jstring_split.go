package jutil

import (
	"regexp"
	"regexp/syntax"
	"sync"
	"unicode/utf8"
)

// Go-only: language-adapter — String.split(String regex) and
// Pattern.split(CharSequence, int) for a regex given as a Go regexp whose
// meaning is the same in RE2 and java.util.regex (no back-references, no
// lookaround), see the Javadoc of String.split(String, int).
//
// The rules the Javadoc gives:
//   - the input is cut before every match and the match is dropped; the
//     input itself is the only element when nothing matches;
//   - a zero-width match at the very beginning of the input never produces
//     an empty leading element;
//   - limit > 0 makes at most limit-1 cuts, the last element holding the
//     rest of the input; limit < 0 makes every cut; limit 0 makes every cut
//     and then drops the trailing empty strings (so the result may be empty).
//
// The matches are found like successive calls of Matcher.find(). That
// differs from Go's FindAll in one place, which this file reproduces: after
// a non-empty match, find() accepts an empty match at the very next position
// (Go skips it).
//
// Observed behaviour (black-box, JDK 21):
//   - "".split("") is [""], and a skipped zero-width match at the start does
//     not count towards the limit ("abc".split("", 2) is ["a", "bc"]);
//   - "baaac".split("a*") is ["b", "", "c"]: find() yields (0,0) (1,4) (4,4)
//     (5,5).
//
// Not reproduced: after an empty match, find() moves on by one UTF-16 char,
// so it can stop between the two halves of a surrogate pair (Java splits
// "\U0001F600x" with "" into three elements). A Go regexp can only stop at
// code point boundaries.

// StringSplitRegexp is s.split(regex) (limit 0) where re is regex.
func StringSplitRegexp(s string, re *regexp.Regexp) []string {
	return stringSplitRegexpLimit(s, re, 0)
}

// stringSplitRegexpLimit is s.split(regex, limit).
func stringSplitRegexpLimit(s string, re *regexp.Regexp, limit int) []string {
	var parts []string
	index := 0
	for _, m := range javaFindAll(s, re) {
		if m[1] == 0 {
			// A zero-width match at the beginning (only one can end at 0).
			continue
		}
		if limit > 0 && len(parts) == limit-1 {
			break
		}
		parts = append(parts, s[index:m[0]])
		index = m[1]
	}
	if index == 0 {
		return []string{s}
	}
	parts = append(parts, s[index:])
	if limit == 0 {
		n := len(parts)
		for n > 0 && parts[n-1] == "" {
			n--
		}
		parts = parts[:n]
	}
	return parts
}

// javaFindAll lists the matches successive Matcher.find() calls return, as
// byte offsets into s.
func javaFindAll(s string, re *regexp.Regexp) [][]int {
	all := re.FindAllStringIndex(s, -1)
	if !regexCanMatchEmpty(re) {
		// Without empty matches both iterations are the same.
		return all
	}
	// Go drops an empty match that starts where the previous match ended;
	// find() keeps it when the previous match was not empty. Each such match
	// can only sit at the end of a non-empty match that Go reports, and Go
	// then goes on searching one code point further, as find() does.
	var out [][]int
	for i, m := range all {
		out = append(out, m)
		if m[0] == m[1] {
			continue
		}
		e := m[1]
		if i+1 < len(all) && all[i+1][0] == e {
			continue // the next match starts here (it is not empty)
		}
		if matchesEmptyAt(s, re, e) {
			out = append(out, []int{e, e})
		}
	}
	// Go walks code points; the WTF-8 form of a lone surrogate is three
	// bytes that Go sees as three invalid ones. Java sees one char, so
	// positions inside it do not exist.
	kept := out[:0]
	for _, m := range out {
		if isCharBoundary(s, m[0]) && isCharBoundary(s, m[1]) {
			kept = append(kept, m)
		}
	}
	return kept
}

func isCharBoundary(s string, i int) bool {
	return i == 0 || i == len(s) || s[i]&0xC0 != 0x80
}

// matchesEmptyAt reports whether the preferred match of re that starts at
// byte offset e (e > 0) is empty. The regexp is run on s from the code point
// before e with that code point consumed first, so that assertions such as
// \b and ^ see the right context.
func matchesEmptyAt(s string, re *regexp.Regexp, e int) bool {
	start := e - 1
	for start > 0 && !isCharBoundary(s, start) {
		start--
	}
	// A valid code point is one `.` for Go; the three bytes of a WTF-8 lone
	// surrogate are three invalid bytes, one `.` each.
	width := 1
	if _, size := utf8.DecodeRuneInString(s[start:]); size != e-start {
		width = e - start
	}
	m := emptyProbe(re, width).FindStringIndex(s[start:])
	return m != nil && m[1] == e-start
}

type emptyProbeKey struct {
	re    *regexp.Regexp
	width int
}

var emptyProbes sync.Map // emptyProbeKey -> *regexp.Regexp

// emptyProbe compiles `^(?s:.{width})(?:re)`.
func emptyProbe(re *regexp.Regexp, width int) *regexp.Regexp {
	key := emptyProbeKey{re, width}
	if p, ok := emptyProbes.Load(key); ok {
		return p.(*regexp.Regexp)
	}
	p := regexp.MustCompile(`^(?s:.{` + itoa(width) + `})(?:` + re.String() + `)`)
	emptyProbes.Store(key, p)
	return p
}

// regexCanMatchEmpty reports whether re can match the empty string somewhere
// (assertions count as empty).
func regexCanMatchEmpty(re *regexp.Regexp) bool {
	t, err := syntax.Parse(re.String(), syntax.Perl)
	if err != nil {
		return true
	}
	return syntaxNullable(t.Simplify())
}

func syntaxNullable(t *syntax.Regexp) bool {
	switch t.Op {
	case syntax.OpEmptyMatch, syntax.OpBeginLine, syntax.OpEndLine, syntax.OpBeginText, syntax.OpEndText,
		syntax.OpWordBoundary, syntax.OpNoWordBoundary, syntax.OpStar, syntax.OpQuest:
		return true
	case syntax.OpNoMatch, syntax.OpLiteral, syntax.OpCharClass, syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		return t.Op == syntax.OpLiteral && len(t.Rune) == 0
	case syntax.OpCapture, syntax.OpPlus:
		return syntaxNullable(t.Sub[0])
	case syntax.OpRepeat:
		return t.Min == 0 || syntaxNullable(t.Sub[0])
	case syntax.OpConcat:
		for _, sub := range t.Sub {
			if !syntaxNullable(sub) {
				return false
			}
		}
		return true
	case syntax.OpAlternate:
		for _, sub := range t.Sub {
			if syntaxNullable(sub) {
				return true
			}
		}
		return false
	}
	return true
}

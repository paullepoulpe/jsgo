package jutil

// Go-only: language-adapter — java.util.regex.Matcher over a Pattern
// (pattern.go), for a whole input string (no regions, no hitEnd/requireEnd,
// no named groups). Indices (start/end) are UTF-16 offsets, as in Java.
//
// Searching resumes from the end of the previous match exactly as
// Matcher.find() does (an empty match advances the next search by one
// char; resuming inside a surrogate pair panics with
// UnsupportedOperationException). Because RE2 cannot start a search in the middle of a string with
// the preceding text as context, a pattern that uses `^` outside a
// character class panics with UnsupportedOperationException when find()
// has to resume past index 0.

import (
	"strings"
)

// Matcher is java.util.regex.Matcher.
type Matcher struct {
	p     *Pattern
	input string
	// groups holds the byte offsets of the last match (Go's submatch index
	// layout), or nil when there is no current match.
	groups []int
	// first and last delimit the previous match, as byte offsets into input
	// (first == -1 before any match); find() resumes at last.
	first int
	last  int
}

// Matcher is Pattern.matcher(CharSequence).
func (p *Pattern) Matcher(input string) *Matcher {
	return &Matcher{p: p, input: input, first: -1, last: 0}
}

// Pattern is Matcher.pattern().
func (m *Matcher) Pattern() *Pattern { return m.p }

// Reset is Matcher.reset().
func (m *Matcher) Reset() *Matcher {
	m.groups = nil
	m.first = -1
	m.last = 0
	return m
}

func (m *Matcher) setMatch(loc []int) bool {
	if loc == nil {
		m.groups = nil
		m.first = -1
		return false
	}
	m.groups = loc
	m.first = loc[0]
	m.last = loc[1]
	return true
}

// Matches is Matcher.matches(): the entire input must match.
func (m *Matcher) Matches() bool {
	return m.setMatch(m.p.full.FindStringSubmatchIndex(m.input))
}

// LookingAt is Matcher.lookingAt(): a match must start at the beginning of
// the input but need not reach its end. RE2's leftmost-first search returns
// a match starting at 0 whenever one exists, with the extent Java's
// backtracking would choose there.
func (m *Matcher) LookingAt() bool {
	loc := m.p.re.FindStringSubmatchIndex(m.input)
	if loc != nil && loc[0] != 0 {
		loc = nil
	}
	return m.setMatch(loc)
}

// Find is Matcher.find(): the next match, starting where the previous one
// ended (one char later after an empty match).
func (m *Matcher) Find() bool {
	next := m.last
	if next == m.first {
		// Empty match: advance by one char (one UTF-16 unit).
		if next >= len(m.input) {
			m.groups = nil
			m.first = -1
			return false
		}
		_, size := decodeWTF8(m.input[next:])
		if size == 4 {
			// Java resumes between the two surrogates of a supplementary
			// character, which a WTF-8 offset cannot express.
			unsupportedJavaRegex(m.p.pattern, "find() resuming inside a surrogate pair after an empty match")
		}
		next += size
	}
	if next > len(m.input) {
		m.groups = nil
		m.first = -1
		return false
	}
	return m.search(next)
}

func (m *Matcher) search(from int) bool {
	if from > 0 && patternHasCaret(m.p.pattern) {
		unsupportedJavaRegex(m.p.pattern, "^ when find() resumes past the start")
	}
	loc := m.p.re.FindStringSubmatchIndex(m.input[from:])
	if loc != nil {
		for i := range loc {
			if loc[i] >= 0 {
				loc[i] += from
			}
		}
	}
	return m.setMatch(loc)
}

// patternHasCaret reports whether a Java regex uses `^` outside a character
// class.
func patternHasCaret(regex string) bool {
	inClass := false
	for i := 0; i < len(regex); i++ {
		switch c := regex[i]; {
		case c == '\\':
			i++
		case inClass:
			if c == ']' {
				inClass = false
			}
		case c == '[':
			inClass = true
			if i+1 < len(regex) && regex[i+1] == '^' {
				i++
			}
			if i+1 < len(regex) && regex[i+1] == ']' {
				i++
			}
		case c == '^':
			return true
		}
	}
	return false
}

// GroupCount is Matcher.groupCount().
func (m *Matcher) GroupCount() int { return m.p.re.NumSubexp() }

func (m *Matcher) checkMatch() {
	if m.groups == nil {
		panic(NewIllegalStateExceptionMsg("No match found"))
	}
}

func (m *Matcher) checkGroup(group int) {
	if group < 0 || group > m.GroupCount() {
		panic(NewIndexOutOfBoundsException("No group " + itoa(group)))
	}
}

// Start is Matcher.start(): the UTF-16 index of the start of the match.
func (m *Matcher) Start() int {
	m.checkMatch()
	return StringLength(m.input[:m.groups[0]])
}

// End is Matcher.end(): the UTF-16 index after the end of the match.
func (m *Matcher) End() int {
	m.checkMatch()
	return StringLength(m.input[:m.groups[1]])
}

// StartGroup is Matcher.start(int): -1 when the group did not match.
func (m *Matcher) StartGroup(group int) int {
	m.checkMatch()
	m.checkGroup(group)
	if m.groups[2*group] < 0 {
		return -1
	}
	return StringLength(m.input[:m.groups[2*group]])
}

// EndGroup is Matcher.end(int): -1 when the group did not match.
func (m *Matcher) EndGroup(group int) int {
	m.checkMatch()
	m.checkGroup(group)
	if m.groups[2*group+1] < 0 {
		return -1
	}
	return StringLength(m.input[:m.groups[2*group+1]])
}

// Group is Matcher.group(int): nil when the group did not take part in the
// match.
func (m *Matcher) Group(group int) *string {
	m.checkMatch()
	m.checkGroup(group)
	if m.groups[2*group] < 0 {
		return nil
	}
	s := m.input[m.groups[2*group]:m.groups[2*group+1]]
	return &s
}

// GroupAll is Matcher.group() (the whole match).
func (m *Matcher) GroupAll() string {
	return *m.Group(0)
}

// appendReplacement is Matcher.appendReplacement(StringBuilder, String)
// with the append position passed in (a byte offset into the input; the
// callers move it to the end of each match). It appends the input from
// appendPos up to the current match, then the replacement with its group
// references expanded.
//
// The replacement syntax is the one the Javadoc of appendReplacement gives:
//   - a backslash quotes the character after it, which is appended as is
//     (`\$` is a dollar sign, `\\` a backslash);
//   - `$n` is capturing group n. Its first (ASCII) digit always belongs to
//     the reference; each following digit is taken only while the number
//     still names a group of the pattern, so with fewer than 12 groups
//     `$12` is group 1 followed by a literal "2";
//   - `${name}` is a named group (a letter, then letters and digits, all
//     ASCII);
//   - a group that did not take part in the match expands to nothing.
//
// Observed behaviour (black-box, JDK 21): the replacement is read left to
// right, and the first malformed reference throws:
//   - IllegalArgumentException "character to be escaped is missing" for a
//     trailing backslash;
//   - IllegalArgumentException "Illegal group reference: group index is
//     missing" for a trailing `$`, and "Illegal group reference" when `$` is
//     followed by anything but an ASCII digit or `{`;
//   - IndexOutOfBoundsException "No group <n>" when n exceeds groupCount();
//   - for `${`, IllegalArgumentException "named capturing group has 0
//     length name" when no letter or digit follows, then "named capturing
//     group is missing trailing '}'", then "capturing group name {<name>}
//     starts with digit character", and finally "No group with name
//     {<name>}" for a well-formed name that the pattern does not define.
func (m *Matcher) appendReplacement(sb *strings.Builder, appendPos int, replacement string) {
	m.checkMatch()
	sb.WriteString(m.input[appendPos:m.groups[0]])
	r := replacement
	// Every special character is ASCII, so scanning the WTF-8 bytes is
	// enough: the bytes of other characters are copied as they are.
	for i := 0; i < len(r); {
		switch r[i] {
		case '\\':
			i++
			if i == len(r) {
				panic(NewIllegalArgumentExceptionMsg("character to be escaped is missing"))
			}
			// The rest of a multi-byte character is copied by the
			// literal case below.
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
				// Pattern has no named groups (pattern.go rejects `(?<`),
				// so no name can be resolved.
				panic(NewIllegalArgumentExceptionMsg("No group with name {" + name + "}"))
			case c >= '0' && c <= '9':
				group = int(c - '0')
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
				m.checkGroup(group)
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

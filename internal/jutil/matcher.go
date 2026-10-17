package jutil

// Go-only: language-adapter — java.util.regex.Matcher over a Pattern
// (pattern.go), for a whole input string (no regions, no hitEnd/requireEnd,
// no named groups). Indices (start/end) are UTF-16 offsets, as in Java.
//
// Where find() looks (Javadoc of Matcher.find(), plus black-box runs of
// JDK 21 for what the Javadoc leaves open):
//   - the first find() after creation or reset() looks from the start;
//   - after a successful match (by find(), matches() or lookingAt()) it
//     looks from where that match ended, or one char further when the
//     match was empty, so that it cannot return the same empty match;
//   - a failed attempt (find(), matches() or lookingAt()) keeps that
//     position but cancels the one-char step: `^` matches empty at 0 in
//     "ab", the next find() fails from 1, and the one after that returns the
//     empty match at 0 again;
//   - except when the one-char step itself would leave the input: find()
//     then fails without changing anything, so it keeps failing until a
//     successful matches()/lookingAt() or a reset().
// Taking the one-char step from just before a supplementary character would
// land between its two surrogates, a position the WTF-8 input cannot hold;
// find() panics with UnsupportedOperationException there.

import (
	"strings"
)

// Matcher is java.util.regex.Matcher.
type Matcher struct {
	p     *Pattern
	input string
	// groups holds the byte offsets of the current match in Go's submatch
	// index layout (the Java pattern's own groups only), or nil when there
	// is no current match.
	groups []int
	// resumeAt is the byte offset where the next find() starts looking: the
	// end of the latest successful match, or 0.
	resumeAt int
	// prevWasEmpty is set when that match was empty (and no failed attempt
	// came since): find() then starts one char after resumeAt.
	prevWasEmpty bool
}

// Matcher is Pattern.matcher(CharSequence).
func (p *Pattern) Matcher(input string) *Matcher {
	return &Matcher{p: p, input: input}
}

// Pattern is Matcher.pattern().
func (m *Matcher) Pattern() *Pattern { return m.p }

// Reset is Matcher.reset(): no current match, and the next find() looks
// from the start again.
func (m *Matcher) Reset() *Matcher {
	m.groups, m.resumeAt, m.prevWasEmpty = nil, 0, false
	return m
}

// record installs the outcome of a match attempt (loc nil: it failed) and
// updates where find() resumes, as described in the file comment.
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
// the input but need not reach its end. RE2's leftmost-first search returns
// a match starting at 0 whenever one exists, with the extent Java's
// backtracking would choose there.
func (m *Matcher) LookingAt() bool {
	loc := m.p.findFrom(m.input, 0)
	if loc != nil && loc[0] != 0 {
		loc = nil
	}
	return m.record(loc)
}

// Find is Matcher.find(): the next match of the pattern, looked for from
// the position the file comment describes.
func (m *Matcher) Find() bool {
	from := m.resumeAt
	if m.prevWasEmpty {
		if from == len(m.input) {
			// Observed (black-box): find() fails here and leaves the state
			// as it is, so a further find() fails too. (JDK 21 then still
			// answers start() and end() with the empty match, but group()
			// with null; this port has no current match, as the Javadoc of
			// start() describes after a failed match.)
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

// GroupCount is Matcher.groupCount().
func (m *Matcher) GroupCount() int { return m.p.groups }

func (m *Matcher) requireMatch() {
	if m.groups == nil {
		panic(NewIllegalStateExceptionMsg("No match found"))
	}
}

func (m *Matcher) requireGroupIndex(group int) {
	if group < 0 || group > m.GroupCount() {
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

// StartGroup is Matcher.start(int): -1 when the group did not match.
func (m *Matcher) StartGroup(group int) int {
	m.requireMatch()
	m.requireGroupIndex(group)
	if m.groups[2*group] < 0 {
		return -1
	}
	return StringLength(m.input[:m.groups[2*group]])
}

// EndGroup is Matcher.end(int): -1 when the group did not match.
func (m *Matcher) EndGroup(group int) int {
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
	m.requireMatch()
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

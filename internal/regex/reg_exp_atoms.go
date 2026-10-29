// Port of com/google/javascript/jscomp/regex/RegExpTree.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package regex

import (
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func writeUTF16Char(sb *strings.Builder, ch uint16) {
	switch {
	case ch < 0x80:
		sb.WriteByte(byte(ch))
	case ch < 0x800:
		buf := [2]byte{
			byte(0xC0 | (ch >> 6)),
			byte(0x80 | (ch & 0x3F)),
		}
		sb.Write(buf[:])
	default:
		buf := [3]byte{
			byte(0xE0 | (ch >> 12)),
			byte(0x80 | ((ch >> 6) & 0x3F)),
			byte(0x80 | (ch & 0x3F)),
		}
		sb.Write(buf[:])
	}
}

func stringFromUTF16Char(ch uint16) string {
	switch {
	case ch < 0x80:
		return string(rune(ch))
	case ch < 0x800:
		buf := [2]byte{
			byte(0xC0 | (ch >> 6)),
			byte(0x80 | (ch & 0x3F)),
		}
		return string(buf[:])
	default:
		buf := [3]byte{
			byte(0xE0 | (ch >> 12)),
			byte(0x80 | ((ch >> 6) & 0x3F)),
			byte(0x80 | (ch & 0x3F)),
		}
		return string(buf[:])
	}
}

func stringFromCodePoint(cp int) string {
	if cp <= 0xFFFF {
		return stringFromUTF16Char(uint16(cp))
	}
	return string(rune(cp))
}

func decodeWTF8Rune(s string) (rune, int) {
	c := s[0]
	if c < 0x80 {
		return rune(c), 1
	}
	if c == 0xED && len(s) >= 3 && s[1] >= 0xA0 && s[1] <= 0xBF && s[2]&0xC0 == 0x80 {
		return rune(c&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F), 3
	}
	return utf8.DecodeRuneInString(s)
}

// regExpTree_singleCharUnit reports whether t is a Text of length 1 UTF-16
// code unit and returns that code unit without allocating.
func regExpTree_singleCharUnit(t RegExpTree) (uint16, bool) {
	text, ok := t.(*RegExpTree_Text)
	if !ok || len(text.text) == 0 {
		return 0, false
	}
	s := text.text
	if len(s) == 1 {
		return uint16(s[0]), true
	}
	if s[0] < 0x80 {
		return 0, false
	}
	r, size := decodeWTF8Rune(s)
	if size == len(s) && r < 0x10000 {
		return uint16(r), true
	}
	return 0, false
}

// regExpTree_isSingleCharText is `x instanceof Text text &&
// text.text.length() == 1`.
//
// Go-only: language-adapter — pattern-matching instanceof as a helper.
func regExpTree_isSingleCharText(t RegExpTree) bool {
	_, ok := regExpTree_singleCharUnit(t)
	return ok
}

func writeDecimalInt(sb *strings.Builder, n int) {
	var buf [12]byte
	sb.Write(strconv.AppendInt(buf[:0], int64(n), 10))
}

var emptyRegExpChildren = []RegExpTree{}

// RegExpTree_Empty represents an empty portion of a RegExp such as the middle
// of "||".
//
// Java: private static final class RegExpTree.Empty
type RegExpTree_Empty struct{}

// NewRegExpTree_Empty constructs an Empty.
//
// Java: Empty()
func NewRegExpTree_Empty() *RegExpTree_Empty {
	return &RegExpTree_Empty{}
}

// RegExpTree_Empty_INSTANCE is the shared Empty.
var RegExpTree_Empty_INSTANCE = NewRegExpTree_Empty()

func (e *RegExpTree_Empty) IsCaseSensitive() bool   { return false }
func (e *RegExpTree_Empty) ContainsAnchor() bool    { return false }
func (e *RegExpTree_Empty) HasCapturingGroup() bool { return false }
func (e *RegExpTree_Empty) NumCapturingGroups() int { return 0 }
func (e *RegExpTree_Empty) Children() []RegExpTree  { return emptyRegExpChildren }

// Simplify returns e.
//
// Java: Empty.simplify(String)
func (e *RegExpTree_Empty) Simplify(flags string) RegExpTree {
	return e
}

// AppendSourceCode emits nothing.
//
// Java: Empty.appendSourceCode(StringBuilder)
func (e *RegExpTree_Empty) AppendSourceCode(sb *strings.Builder) {
	// No output
}

// AppendDebugInfo emits nothing.
//
// Java: Empty.appendDebugInfo(StringBuilder)
func (e *RegExpTree_Empty) AppendDebugInfo(sb *strings.Builder) {
	// No output
}

// Equals reports whether o is an Empty.
//
// Java: Empty.equals(Object)
func (e *RegExpTree_Empty) Equals(o any) bool {
	_, ok := o.(*RegExpTree_Empty)
	return ok
}

// HashCode is a constant.
//
// Java: Empty.hashCode()
func (e *RegExpTree_Empty) HashCode() int32 {
	return 0x7ee06141
}

// Java: RegExpTree.toString()
func (e *RegExpTree_Empty) String() string { return regExpTreeString(e) }

// RegExpTree_MatchesWholeInput reports true if, but not necessarily always
// when, the given regular expression must match the whole input or none of
// it.
//
// Java: RegExpTree.matchesWholeInput(RegExpTree, String)
func RegExpTree_MatchesWholeInput(t RegExpTree, flags string) bool {
	if strings.IndexByte(flags, 'm') >= 0 {
		return false
	}

	c, ok := t.(*RegExpTree_Concatenation)
	if !ok {
		return false
	}

	if len(c.elements) == 0 {
		return false
	}
	first := c.elements[0]
	last := c.elements[len(c.elements)-1]
	firstAnchor, ok1 := first.(*RegExpTree_Anchor)
	lastAnchor, ok2 := last.(*RegExpTree_Anchor)
	if !(ok1 && ok2) {
		return false
	}
	return firstAnchor.type_ == '^' && lastAnchor.type_ == '$'
}

// RegExpTree_Anchor represents an anchor, namely ^ or $.
//
// Java: public static final class RegExpTree.Anchor
type RegExpTree_Anchor struct {
	type_ uint16
}

// NewRegExpTree_Anchor constructs an Anchor.
//
// Java: Anchor(char)
func NewRegExpTree_Anchor(type_ uint16) *RegExpTree_Anchor {
	return &RegExpTree_Anchor{type_: type_}
}

func (a *RegExpTree_Anchor) IsCaseSensitive() bool   { return false }
func (a *RegExpTree_Anchor) HasCapturingGroup() bool { return false }
func (a *RegExpTree_Anchor) NumCapturingGroups() int { return 0 }
func (a *RegExpTree_Anchor) Children() []RegExpTree  { return emptyRegExpChildren }

// Simplify returns a.
//
// Java: Anchor.simplify(String)
func (a *RegExpTree_Anchor) Simplify(flags string) RegExpTree {
	return a
}

// ContainsAnchor is true.
//
// Java: Anchor.containsAnchor()
func (a *RegExpTree_Anchor) ContainsAnchor() bool {
	return true
}

// AppendSourceCode emits ^ or $.
//
// Java: Anchor.appendSourceCode(StringBuilder)
func (a *RegExpTree_Anchor) AppendSourceCode(sb *strings.Builder) {
	writeUTF16Char(sb, a.type_)
}

// AppendDebugInfo emits the anchor character.
//
// Java: Anchor.appendDebugInfo(StringBuilder)
func (a *RegExpTree_Anchor) AppendDebugInfo(sb *strings.Builder) {
	writeUTF16Char(sb, a.type_)
}

// Equals compares anchor types.
//
// Java: Anchor.equals(Object)
func (a *RegExpTree_Anchor) Equals(o any) bool {
	anchor, ok := o.(*RegExpTree_Anchor)
	return ok && a.type_ == anchor.type_
}

// HashCode hashes the anchor type.
//
// Java: Anchor.hashCode()
func (a *RegExpTree_Anchor) HashCode() int32 {
	return int32(uint32(a.type_) ^ 0xe85317ff)
}

// Java: RegExpTree.toString()
func (a *RegExpTree_Anchor) String() string { return regExpTreeString(a) }

// RegExpTree_WordBoundary represents \b or \B.
//
// Java: public static final class RegExpTree.WordBoundary
type RegExpTree_WordBoundary struct {
	type_ uint16
}

// NewRegExpTree_WordBoundary constructs a WordBoundary.
//
// Java: WordBoundary(char)
func NewRegExpTree_WordBoundary(type_ uint16) *RegExpTree_WordBoundary {
	return &RegExpTree_WordBoundary{type_: type_}
}

func (w *RegExpTree_WordBoundary) IsCaseSensitive() bool   { return false }
func (w *RegExpTree_WordBoundary) ContainsAnchor() bool    { return false }
func (w *RegExpTree_WordBoundary) HasCapturingGroup() bool { return false }
func (w *RegExpTree_WordBoundary) NumCapturingGroups() int { return 0 }
func (w *RegExpTree_WordBoundary) Children() []RegExpTree  { return emptyRegExpChildren }

// Simplify returns w.
//
// Java: WordBoundary.simplify(String)
func (w *RegExpTree_WordBoundary) Simplify(flags string) RegExpTree {
	return w
}

// AppendSourceCode emits \b or \B.
//
// Java: WordBoundary.appendSourceCode(StringBuilder)
func (w *RegExpTree_WordBoundary) AppendSourceCode(sb *strings.Builder) {
	sb.WriteByte('\\')
	writeUTF16Char(sb, w.type_)
}

// AppendDebugInfo emits the boundary type.
//
// Java: WordBoundary.appendDebugInfo(StringBuilder)
func (w *RegExpTree_WordBoundary) AppendDebugInfo(sb *strings.Builder) {
	writeUTF16Char(sb, w.type_)
}

// Equals compares boundary types.
//
// Java: WordBoundary.equals(Object)
func (w *RegExpTree_WordBoundary) Equals(o any) bool {
	wordBoundary, ok := o.(*RegExpTree_WordBoundary)
	return ok && w.type_ == wordBoundary.type_
}

// HashCode hashes the boundary type.
//
// Java: WordBoundary.hashCode()
func (w *RegExpTree_WordBoundary) HashCode() int32 {
	return 0x5673aa29 ^ int32(w.type_)
}

// Java: RegExpTree.toString()
func (w *RegExpTree_WordBoundary) String() string { return regExpTreeString(w) }

// RegExpTree_BackReference represents a reference to a previous group such
// as \1 or \2.
//
// Java: public static final class RegExpTree.BackReference
type RegExpTree_BackReference struct {
	groupIndex int
}

// NewRegExpTree_BackReference constructs a BackReference.
//
// Java: BackReference(int)
func NewRegExpTree_BackReference(groupIndex int) *RegExpTree_BackReference {
	jutil.CheckArgument(groupIndex >= 0 && groupIndex <= 99)
	return &RegExpTree_BackReference{groupIndex: groupIndex}
}

func (b *RegExpTree_BackReference) IsCaseSensitive() bool   { return false }
func (b *RegExpTree_BackReference) ContainsAnchor() bool    { return false }
func (b *RegExpTree_BackReference) HasCapturingGroup() bool { return false }
func (b *RegExpTree_BackReference) NumCapturingGroups() int { return 0 }
func (b *RegExpTree_BackReference) Children() []RegExpTree  { return emptyRegExpChildren }

// Simplify returns b.
//
// Java: BackReference.simplify(String)
func (b *RegExpTree_BackReference) Simplify(flags string) RegExpTree {
	return b
}

// AppendSourceCode emits \N.
//
// Java: BackReference.appendSourceCode(StringBuilder)
func (b *RegExpTree_BackReference) AppendSourceCode(sb *strings.Builder) {
	sb.WriteByte('\\')
	writeDecimalInt(sb, b.groupIndex)
}

// AppendDebugInfo emits the group index.
//
// Java: BackReference.appendDebugInfo(StringBuilder)
func (b *RegExpTree_BackReference) AppendDebugInfo(sb *strings.Builder) {
	writeDecimalInt(sb, b.groupIndex)
}

// Equals compares group indices.
//
// Java: BackReference.equals(Object)
func (b *RegExpTree_BackReference) Equals(o any) bool {
	backReference, ok := o.(*RegExpTree_BackReference)
	return ok && b.groupIndex == backReference.groupIndex
}

// HashCode hashes the group index.
//
// Java: BackReference.hashCode()
func (b *RegExpTree_BackReference) HashCode() int32 {
	return int32(uint32(0xff072663) ^ uint32(int32(b.groupIndex)))
}

// Java: RegExpTree.toString()
func (b *RegExpTree_BackReference) String() string { return regExpTreeString(b) }

// RegExpTree_NamedBackReference represents a reference to a previous named
// group.
//
// Java: public static final class RegExpTree.NamedBackReference
type RegExpTree_NamedBackReference struct {
	groupName string
}

// NewRegExpTree_NamedBackReference constructs a NamedBackReference.
//
// Java: NamedBackReference(String)
func NewRegExpTree_NamedBackReference(groupName string) *RegExpTree_NamedBackReference {
	return &RegExpTree_NamedBackReference{groupName: groupName}
}

func (n *RegExpTree_NamedBackReference) IsCaseSensitive() bool   { return false }
func (n *RegExpTree_NamedBackReference) ContainsAnchor() bool    { return false }
func (n *RegExpTree_NamedBackReference) HasCapturingGroup() bool { return false }
func (n *RegExpTree_NamedBackReference) NumCapturingGroups() int { return 0 }
func (n *RegExpTree_NamedBackReference) Children() []RegExpTree  { return emptyRegExpChildren }

// Simplify returns n.
//
// Java: NamedBackReference.simplify(String)
func (n *RegExpTree_NamedBackReference) Simplify(flags string) RegExpTree {
	return n
}

// AppendSourceCode emits \k<name>.
//
// Java: NamedBackReference.appendSourceCode(StringBuilder)
func (n *RegExpTree_NamedBackReference) AppendSourceCode(sb *strings.Builder) {
	sb.WriteString("\\k<")
	sb.WriteString(n.groupName)
	sb.WriteByte('>')
}

// AppendDebugInfo emits the group name.
//
// Java: NamedBackReference.appendDebugInfo(StringBuilder)
func (n *RegExpTree_NamedBackReference) AppendDebugInfo(sb *strings.Builder) {
	sb.WriteString(n.groupName)
}

// Equals compares group names.
//
// Java: NamedBackReference.equals(Object)
func (n *RegExpTree_NamedBackReference) Equals(o any) bool {
	namedBackReference, ok := o.(*RegExpTree_NamedBackReference)
	return ok && n.groupName == namedBackReference.groupName
}

// HashCode hashes the group name.
//
// Java: NamedBackReference.hashCode()
func (n *RegExpTree_NamedBackReference) HashCode() int32 {
	return jutil.StringHashCode(n.groupName) // Java: Objects.hashCode(groupName)
}

// Java: RegExpTree.toString()
func (n *RegExpTree_NamedBackReference) String() string { return regExpTreeString(n) }

// RegExpTree_NamedCaptureGroup represents a named capture group.
//
// Java: public static final class RegExpTree.NamedCaptureGroup
type RegExpTree_NamedCaptureGroup struct {
	body RegExpTree
	name string
}

// NewRegExpTree_NamedCaptureGroup constructs a NamedCaptureGroup.
//
// Java: NamedCaptureGroup(RegExpTree, String)
func NewRegExpTree_NamedCaptureGroup(body RegExpTree, name string) *RegExpTree_NamedCaptureGroup {
	return &RegExpTree_NamedCaptureGroup{
		body: body,
		name: name,
	}
}

func (n *RegExpTree_NamedCaptureGroup) IsCaseSensitive() bool  { return n.body.IsCaseSensitive() }
func (n *RegExpTree_NamedCaptureGroup) ContainsAnchor() bool   { return n.body.ContainsAnchor() }
func (n *RegExpTree_NamedCaptureGroup) Children() []RegExpTree { return []RegExpTree{n.body} }

// Simplify simplifies the body.
//
// Java: NamedCaptureGroup.simplify(String)
func (n *RegExpTree_NamedCaptureGroup) Simplify(flags string) RegExpTree {
	return NewRegExpTree_NamedCaptureGroup(n.body.Simplify(flags), n.name)
}

// Java: RegExpTree.hasCapturingGroup()
func (n *RegExpTree_NamedCaptureGroup) HasCapturingGroup() bool {
	return true
}

// NumCapturingGroups is one plus the body's.
//
// Java: NamedCaptureGroup.numCapturingGroups()
func (n *RegExpTree_NamedCaptureGroup) NumCapturingGroups() int {
	return 1 + n.body.NumCapturingGroups()
}

// AppendSourceCode emits (?<name>...).
//
// Java: NamedCaptureGroup.appendSourceCode(StringBuilder)
func (n *RegExpTree_NamedCaptureGroup) AppendSourceCode(sb *strings.Builder) {
	sb.WriteString("(?<")
	sb.WriteString(n.name)
	sb.WriteByte('>')
	n.body.AppendSourceCode(sb)
	sb.WriteByte(')')
}

// AppendDebugInfo emits the name.
//
// Java: NamedCaptureGroup.appendDebugInfo(StringBuilder)
func (n *RegExpTree_NamedCaptureGroup) AppendDebugInfo(sb *strings.Builder) {
	sb.WriteString(" name=")
	sb.WriteString(n.name)
}

// Equals compares name and body.
//
// Java: NamedCaptureGroup.equals(Object)
func (n *RegExpTree_NamedCaptureGroup) Equals(o any) bool {
	namedCaptureGroup, ok := o.(*RegExpTree_NamedCaptureGroup)
	return ok &&
		n.name == namedCaptureGroup.name &&
		n.body.Equals(namedCaptureGroup.body)
}

// HashCode hashes name and body.
//
// Java: NamedCaptureGroup.hashCode()
func (n *RegExpTree_NamedCaptureGroup) HashCode() int32 {
	return jutil.StringHashCode(n.name) ^ n.body.HashCode() // Java: Objects.hashCode(name)
}

// Java: RegExpTree.toString()
func (n *RegExpTree_NamedCaptureGroup) String() string { return regExpTreeString(n) }

// RegExpTree_Text represents a run of non-special characters such as
// "foobar".
//
// Java: public static final class RegExpTree.Text
type RegExpTree_Text struct {
	text string
}

// NewRegExpTree_Text constructs a Text.
//
// Java: Text(String)
func NewRegExpTree_Text(text string) *RegExpTree_Text {
	return &RegExpTree_Text{text: text}
}

func (t *RegExpTree_Text) ContainsAnchor() bool    { return false }
func (t *RegExpTree_Text) HasCapturingGroup() bool { return false }
func (t *RegExpTree_Text) NumCapturingGroups() int { return 0 }
func (t *RegExpTree_Text) Children() []RegExpTree  { return emptyRegExpChildren }

// RegExpTree_Text_EscapeRegularCharOnto escapes ch onto sb. ch is the
// code-unit to escape; next is the next code-unit or -1 if indeterminable.
//
// Java: Text.escapeRegularCharOnto(char, int, StringBuilder)
func RegExpTree_Text_EscapeRegularCharOnto(ch uint16, next int, sb *strings.Builder) {
	switch ch {
	case '$', '^', '*', '(', ')', '+', '[', '|', '.', '/', '?':
		sb.WriteByte('\\')
		sb.WriteByte(byte(ch))
	case '{':
		// If possibly part of a repetition, then escape.
		// Concatenation is handled by the digitsMightBleed check.
		if '0' <= next && next <= '9' {
			sb.WriteByte('\\')
		}
		sb.WriteByte('{')
	default:
		RegExpTree_EscapeCharOnto(ch, sb)
	}
}

// Simplify drops empty text and case-canonicalizes under the 'i' flag.
//
// Java: Text.simplify(String)
func (t *RegExpTree_Text) Simplify(flags string) RegExpTree {
	if t.text == "" {
		return RegExpTree_Empty_INSTANCE
	}
	if strings.IndexByte(flags, 'i') >= 0 {
		canonicalized := CaseCanonicalize_CaseCanonicalize(t.text)
		if t.text != canonicalized {
			return NewRegExpTree_Text(canonicalized)
		}
	}
	return t
}

// IsCaseSensitive reports whether any character is case sensitive.
//
// Java: Text.isCaseSensitive()
func (t *RegExpTree_Text) IsCaseSensitive() bool {
	s := t.text
	for i := 0; i < len(s); {
		if s[i] < 0x80 {
			if CaseCanonicalize_CASE_SENSITIVE.Contains(int(s[i])) {
				return true
			}
			i++
			continue
		}
		r, size := decodeWTF8Rune(s[i:])
		i += size
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			if CaseCanonicalize_CASE_SENSITIVE.Contains(int(uint16(hi))) ||
				CaseCanonicalize_CASE_SENSITIVE.Contains(int(uint16(lo))) {
				return true
			}
		} else if CaseCanonicalize_CASE_SENSITIVE.Contains(int(uint16(r))) {
			return true
		}
	}
	return false
}

// AppendSourceCode emits the escaped text.
//
// Java: Text.appendSourceCode(StringBuilder)
func (t *RegExpTree_Text) AppendSourceCode(sb *strings.Builder) {
	s := t.text
	if jutil.IsASCII(s) {
		for i, n := 0, len(s); i < n; i++ {
			next := -1
			if i+1 < n {
				next = int(s[i+1])
			}
			RegExpTree_Text_EscapeRegularCharOnto(uint16(s[i]), next, sb)
		}
		return
	}
	nextUnit := func(pos int, pendingLo int) (uint16, int, int, bool) {
		if pendingLo >= 0 {
			return uint16(pendingLo), pos, -1, true
		}
		if pos >= len(s) {
			return 0, pos, -1, false
		}
		r, size := decodeWTF8Rune(s[pos:])
		pos += size
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			return uint16(hi), pos, int(uint16(lo)), true
		}
		return uint16(r), pos, -1, true
	}
	ch, pos, pendingLo, ok := nextUnit(0, -1)
	for ok {
		peekCh, nextPos, nextPendingLo, hasNext := nextUnit(pos, pendingLo)
		next := -1
		if hasNext {
			next = int(peekCh)
		}
		RegExpTree_Text_EscapeRegularCharOnto(ch, next, sb)
		ch, pos, pendingLo, ok = peekCh, nextPos, nextPendingLo, hasNext
	}
}

// AppendDebugInfo emits `text`.
//
// Java: Text.appendDebugInfo(StringBuilder)
func (t *RegExpTree_Text) AppendDebugInfo(sb *strings.Builder) {
	sb.WriteByte('`')
	sb.WriteString(t.text)
	sb.WriteByte('`')
}

// Equals compares the text.
//
// Java: Text.equals(Object)
func (t *RegExpTree_Text) Equals(o any) bool {
	text, ok := o.(*RegExpTree_Text)
	return ok && t.text == text.text
}

// HashCode hashes the text.
//
// Java: Text.hashCode()
func (t *RegExpTree_Text) HashCode() int32 {
	return jutil.StringHashCode(t.text) ^ 0x617e310
}

// Java: RegExpTree.toString()
func (t *RegExpTree_Text) String() string { return regExpTreeString(t) }

// RegExpTree_UnicodePropertyEscape represents a Unicode Property Escape such
// as in /\p{Script=Greek}/u.
//
// Java: public static final class RegExpTree.UnicodePropertyEscape
type RegExpTree_UnicodePropertyEscape struct {
	propertyName  string // empty for the shorthand form \p{Value}
	propertyValue string
	negated       bool
}

func (u *RegExpTree_UnicodePropertyEscape) IsCaseSensitive() bool   { return false }
func (u *RegExpTree_UnicodePropertyEscape) ContainsAnchor() bool    { return false }
func (u *RegExpTree_UnicodePropertyEscape) HasCapturingGroup() bool { return false }
func (u *RegExpTree_UnicodePropertyEscape) NumCapturingGroups() int { return 0 }
func (u *RegExpTree_UnicodePropertyEscape) Children() []RegExpTree  { return emptyRegExpChildren }

func newUnicodePropertyEscape(propertyName string, hasName bool, propertyValue string, negated bool) (*RegExpTree_UnicodePropertyEscape, *RegExpParseError) {
	if hasName && propertyName == "" {
		return nil, illegalArgError("if '=' is present in a unicode property escape, the name cannot be empty")
	}
	if propertyValue == "" {
		return nil, illegalArgError("unicode property escape value cannot be empty")
	}
	return &RegExpTree_UnicodePropertyEscape{
		propertyName:  propertyName,
		propertyValue: propertyValue,
		negated:       negated,
	}, nil
}

// NewRegExpTree_UnicodePropertyEscape constructs a UnicodePropertyEscape;
// propertyName is "" for the shorthand form \p{Value}.
//
// Java: UnicodePropertyEscape(@Nullable String, String, boolean)
func NewRegExpTree_UnicodePropertyEscape(propertyName string, propertyValue string, negated bool) *RegExpTree_UnicodePropertyEscape {
	u, err := newUnicodePropertyEscape(propertyName, false, propertyValue, negated)
	if err != nil {
		panic(err.toThrowable())
	}
	return u
}

// Simplify returns u.
//
// Java: UnicodePropertyEscape.simplify(String)
func (u *RegExpTree_UnicodePropertyEscape) Simplify(flags string) RegExpTree {
	return u
}

// AppendSourceCode emits \p{...} or \P{...}.
//
// Java: UnicodePropertyEscape.appendSourceCode(StringBuilder)
func (u *RegExpTree_UnicodePropertyEscape) AppendSourceCode(sb *strings.Builder) {
	if u.negated {
		sb.WriteString("\\P{")
	} else {
		sb.WriteString("\\p{")
	}
	if u.propertyName != "" {
		sb.WriteString(u.propertyName)
		sb.WriteByte('=')
	}
	sb.WriteString(u.propertyValue)
	sb.WriteByte('}')
}

// AppendDebugInfo emits nothing.
//
// Java: UnicodePropertyEscape.appendDebugInfo(StringBuilder)
func (u *RegExpTree_UnicodePropertyEscape) AppendDebugInfo(sb *strings.Builder) {
	// Nothing besides properties and possible negation.
}

// Equals compares negation, name and value.
//
// Java: UnicodePropertyEscape.equals(Object)
func (u *RegExpTree_UnicodePropertyEscape) Equals(o any) bool {
	unicodePropertyEscape, ok := o.(*RegExpTree_UnicodePropertyEscape)
	return ok &&
		u.negated == unicodePropertyEscape.negated &&
		u.propertyName == unicodePropertyEscape.propertyName &&
		u.propertyValue == unicodePropertyEscape.propertyValue
}

// HashCode hashes negation, name and value.
//
// Java: UnicodePropertyEscape.hashCode()
func (u *RegExpTree_UnicodePropertyEscape) HashCode() int32 {
	// Java: Objects.hash(negated, propertyName, propertyValue), i.e.
	// Arrays.hashCode(new Object[] {negated, propertyName, propertyValue}).
	negatedHash := int32(1237) // Boolean.hashCode(false)
	if u.negated {
		negatedHash = 1231
	}
	nameHash := int32(0)
	if u.propertyName != "" {
		nameHash = jutil.StringHashCode(u.propertyName)
	}
	result := int32(1)
	result = 31*result + negatedHash
	result = 31*result + nameHash
	result = 31*result + jutil.StringHashCode(u.propertyValue)
	return result
}

// Java: RegExpTree.toString()
func (u *RegExpTree_UnicodePropertyEscape) String() string { return regExpTreeString(u) }

// RegExpTree_EscapeCharOnto escapes a code unit for regular expression
// source.
//
// Java: RegExpTree.escapeCharOnto(char, StringBuilder)
func RegExpTree_EscapeCharOnto(ch uint16, sb *strings.Builder) {
	const hexDigits = "0123456789abcdef"
	switch ch {
	case '\u0000':
		sb.WriteString("\\0")
	case '\f':
		sb.WriteString("\\f")
	case '\t':
		sb.WriteString("\\t")
	case '\n':
		sb.WriteString("\\n")
	case '\r':
		sb.WriteString("\\r")
	case '\\':
		sb.WriteString("\\\\")
	default:
		if ch < 0x20 || ch >= 0x7f {
			if ch >= 0x100 {
				buf := [6]byte{
					'\\', 'u',
					hexDigits[(ch>>12)&0xf],
					hexDigits[(ch>>8)&0xf],
					hexDigits[(ch>>4)&0xf],
					hexDigits[ch&0xf],
				}
				sb.Write(buf[:])
			} else {
				buf := [4]byte{
					'\\', 'x',
					hexDigits[(ch>>4)&0xf],
					hexDigits[ch&0xf],
				}
				sb.Write(buf[:])
			}
		} else {
			sb.WriteByte(byte(ch))
		}
	}
}

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
	"math"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// javaIntegerMaxValue is Integer.MAX_VALUE.
const javaIntegerMaxValue = math.MaxInt32

// RegExpTree is an AST for JavaScript regular expressions.
//
// Java: abstract class RegExpTree.
type RegExpTree interface {
	// Simplify returns a simpler regular expression that is semantically the
	// same assuming the given flags (e.g. "igm").
	Simplify(flags string) RegExpTree
	// IsCaseSensitive reports whether the presence or absence of an "i" flag
	// would change the meaning of this regular expression.
	IsCaseSensitive() bool
	// ContainsAnchor reports whether the regular expression contains an
	// anchor: ^ or $.
	ContainsAnchor() bool
	// HasCapturingGroup reports whether the regular expression contains
	// capturing groups.
	HasCapturingGroup() bool
	// NumCapturingGroups is the number of capturing groups.
	NumCapturingGroups() int
	// Children returns the children of this node.
	Children() []RegExpTree
	// AppendSourceCode appends this regular expression source to the given
	// buffer.
	AppendSourceCode(sb *jutil.StringBuilder)
	AppendDebugInfo(sb *jutil.StringBuilder)
	String() string
	Equals(o any) bool
	HashCode() int32
}

// regExpTreeString renders the regular expression as a literal, e.g. "/foo/".
//
// Java: RegExpTree.toString()
func regExpTreeString(t RegExpTree) string {
	sb := jutil.NewStringBuilder()
	sb.AppendChar('/')
	t.AppendSourceCode(sb)
	// Don't emit a regular expression that looks like a line comment start.
	if sb.Length() == 1 {
		sb.Append("(?:)")
	}
	sb.AppendChar('/')
	return sb.String()
}

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

// RegExpTree_RegExpTreeAtomBase holds RegExpTreeAtom's base implementations.
// It represents a node that never has children such as an anchor or charset.
//
// Java: abstract static class RegExpTree.RegExpTreeAtom
type RegExpTree_RegExpTreeAtomBase struct{}

// IsCaseSensitive is false for atoms unless overridden.
//
// Java: RegExpTreeAtom.isCaseSensitive()
func (a *RegExpTree_RegExpTreeAtomBase) IsCaseSensitive() bool {
	return false
}

// ContainsAnchor is false for atoms unless overridden.
//
// Java: RegExpTreeAtom.containsAnchor()
func (a *RegExpTree_RegExpTreeAtomBase) ContainsAnchor() bool {
	return false
}

// HasCapturingGroup is false for atoms.
//
// Java: RegExpTree.hasCapturingGroup()
func (a *RegExpTree_RegExpTreeAtomBase) HasCapturingGroup() bool {
	return false
}

// NumCapturingGroups is 0 for atoms.
//
// Java: RegExpTreeAtom.numCapturingGroups()
func (a *RegExpTree_RegExpTreeAtomBase) NumCapturingGroups() int {
	return 0
}

// Children is empty for atoms.
//
// Java: RegExpTreeAtom.children()
func (a *RegExpTree_RegExpTreeAtomBase) Children() []RegExpTree {
	return []RegExpTree{}
}

// RegExpTree_Empty represents an empty portion of a RegExp such as the middle
// of "||".
//
// Java: private static final class RegExpTree.Empty
type RegExpTree_Empty struct {
	RegExpTree_RegExpTreeAtomBase
}

// NewRegExpTree_Empty constructs an Empty.
//
// Java: Empty()
func NewRegExpTree_Empty() *RegExpTree_Empty {
	return &RegExpTree_Empty{}
}

// RegExpTree_Empty_INSTANCE is the shared Empty.
var RegExpTree_Empty_INSTANCE = NewRegExpTree_Empty()

// Simplify returns e.
//
// Java: Empty.simplify(String)
func (e *RegExpTree_Empty) Simplify(flags string) RegExpTree {
	return e
}

// AppendSourceCode emits nothing.
//
// Java: Empty.appendSourceCode(StringBuilder)
func (e *RegExpTree_Empty) AppendSourceCode(sb *jutil.StringBuilder) {
	// No output
}

// AppendDebugInfo emits nothing.
//
// Java: Empty.appendDebugInfo(StringBuilder)
func (e *RegExpTree_Empty) AppendDebugInfo(sb *jutil.StringBuilder) {
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

// RegExpTree_Anchor represents an anchor, namely ^ or $.
//
// Java: public static final class RegExpTree.Anchor
type RegExpTree_Anchor struct {
	RegExpTree_RegExpTreeAtomBase
	type_ uint16
}

// NewRegExpTree_Anchor constructs an Anchor.
//
// Java: Anchor(char)
func NewRegExpTree_Anchor(type_ uint16) *RegExpTree_Anchor {
	return &RegExpTree_Anchor{type_: type_}
}

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
func (a *RegExpTree_Anchor) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.AppendChar(a.type_)
}

// AppendDebugInfo emits the anchor character.
//
// Java: Anchor.appendDebugInfo(StringBuilder)
func (a *RegExpTree_Anchor) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.AppendChar(a.type_)
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
	RegExpTree_RegExpTreeAtomBase
	type_ uint16
}

// NewRegExpTree_WordBoundary constructs a WordBoundary.
//
// Java: WordBoundary(char)
func NewRegExpTree_WordBoundary(type_ uint16) *RegExpTree_WordBoundary {
	return &RegExpTree_WordBoundary{type_: type_}
}

// Simplify returns w.
//
// Java: WordBoundary.simplify(String)
func (w *RegExpTree_WordBoundary) Simplify(flags string) RegExpTree {
	return w
}

// AppendSourceCode emits \b or \B.
//
// Java: WordBoundary.appendSourceCode(StringBuilder)
func (w *RegExpTree_WordBoundary) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.AppendChar('\\').AppendChar(w.type_)
}

// AppendDebugInfo emits the boundary type.
//
// Java: WordBoundary.appendDebugInfo(StringBuilder)
func (w *RegExpTree_WordBoundary) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.AppendChar(w.type_)
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
	RegExpTree_RegExpTreeAtomBase
	groupIndex int
}

// NewRegExpTree_BackReference constructs a BackReference.
//
// Java: BackReference(int)
func NewRegExpTree_BackReference(groupIndex int) *RegExpTree_BackReference {
	jutil.CheckArgument(groupIndex >= 0 && groupIndex <= 99)
	return &RegExpTree_BackReference{groupIndex: groupIndex}
}

// Simplify returns b.
//
// Java: BackReference.simplify(String)
func (b *RegExpTree_BackReference) Simplify(flags string) RegExpTree {
	return b
}

// AppendSourceCode emits \N.
//
// Java: BackReference.appendSourceCode(StringBuilder)
func (b *RegExpTree_BackReference) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.AppendChar('\\').AppendInt(int64(b.groupIndex))
}

// AppendDebugInfo emits the group index.
//
// Java: BackReference.appendDebugInfo(StringBuilder)
func (b *RegExpTree_BackReference) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.AppendInt(int64(b.groupIndex))
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
	RegExpTree_RegExpTreeAtomBase
	groupName string
}

// NewRegExpTree_NamedBackReference constructs a NamedBackReference.
//
// Java: NamedBackReference(String)
func NewRegExpTree_NamedBackReference(groupName string) *RegExpTree_NamedBackReference {
	return &RegExpTree_NamedBackReference{groupName: groupName}
}

// Simplify returns n.
//
// Java: NamedBackReference.simplify(String)
func (n *RegExpTree_NamedBackReference) Simplify(flags string) RegExpTree {
	return n
}

// AppendSourceCode emits \k<name>.
//
// Java: NamedBackReference.appendSourceCode(StringBuilder)
func (n *RegExpTree_NamedBackReference) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.Append("\\k<").Append(n.groupName).AppendChar('>')
}

// AppendDebugInfo emits the group name.
//
// Java: NamedBackReference.appendDebugInfo(StringBuilder)
func (n *RegExpTree_NamedBackReference) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.Append(n.groupName)
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

// RegExpTree_Text represents a run of non-special characters such as
// "foobar".
//
// Java: public static final class RegExpTree.Text
type RegExpTree_Text struct {
	RegExpTree_RegExpTreeAtomBase
	text string
}

// NewRegExpTree_Text constructs a Text.
//
// Java: Text(String)
func NewRegExpTree_Text(text string) *RegExpTree_Text {
	return &RegExpTree_Text{text: text}
}

// RegExpTree_Text_EscapeRegularCharOnto escapes ch onto sb. ch is the
// code-unit to escape; next is the next code-unit or -1 if indeterminable.
//
// Java: Text.escapeRegularCharOnto(char, int, StringBuilder)
func RegExpTree_Text_EscapeRegularCharOnto(ch uint16, next int, sb *jutil.StringBuilder) {
	switch ch {
	case '$', '^', '*', '(', ')', '+', '[', '|', '.', '/', '?':
		sb.AppendChar('\\').AppendChar(ch)
	case '{':
		// If possibly part of a repetition, then escape.
		// Concatenation is handled by the digitsMightBleed check.
		if '0' <= next && next <= '9' {
			sb.AppendChar('\\')
		}
		sb.AppendChar(ch)
	default:
		RegExpTree_EscapeCharOnto(ch, sb)
	}
}

// Simplify drops empty text and case-canonicalizes under the 'i' flag.
//
// Java: Text.simplify(String)
func (t *RegExpTree_Text) Simplify(flags string) RegExpTree {
	n := jutil.StringLength(t.text)
	if n == 0 {
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
	chars := jutil.StringToUTF16(t.text) // Java: text.charAt(i)
	for i, n := 0, len(chars); i < n; i++ {
		if CaseCanonicalize_CASE_SENSITIVE.Contains(int(chars[i])) {
			return true
		}
	}
	return false
}

// AppendSourceCode emits the escaped text.
//
// Java: Text.appendSourceCode(StringBuilder)
func (t *RegExpTree_Text) AppendSourceCode(sb *jutil.StringBuilder) {
	chars := jutil.StringToUTF16(t.text) // Java: text.charAt(i)
	for i, n := 0, len(chars); i < n; i++ {
		next := -1
		if i+1 < n {
			next = int(chars[i+1])
		}
		RegExpTree_Text_EscapeRegularCharOnto(chars[i], next, sb)
	}
}

// AppendDebugInfo emits `text`.
//
// Java: Text.appendDebugInfo(StringBuilder)
func (t *RegExpTree_Text) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.AppendChar('`').Append(t.text).AppendChar('`')
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

// RegExpTree_Repetition represents a repeating item such as ...+, ...*, or
// ...{0,1}.
//
// Java: public static final class RegExpTree.Repetition
type RegExpTree_Repetition struct {
	body   RegExpTree
	min    int
	max    int
	greedy bool
}

// NewRegExpTree_Repetition constructs a Repetition.
//
// Java: Repetition(RegExpTree, int, int, boolean)
func NewRegExpTree_Repetition(body RegExpTree, min_ int, max_ int, greedy bool) *RegExpTree_Repetition {
	return &RegExpTree_Repetition{
		body:   body,
		min:    min_,
		max:    max_,
		greedy: greedy,
	}
}

// Simplify folds nested and trivial repetitions.
//
// Java: Repetition.simplify(String)
func (r *RegExpTree_Repetition) Simplify(flags string) RegExpTree {
	body := r.body.Simplify(flags)
	if r.max == 0 && !body.HasCapturingGroup() {
		return RegExpTree_Empty_INSTANCE
	}
	if _, ok := body.(*RegExpTree_Empty); ok || RegExpTree_NEVER_MATCHES.Equals(body) {
		return body
	}
	min_ := r.min
	max_ := r.max
	if rbody, ok := body.(*RegExpTree_Repetition); ok {
		if rbody.greedy == r.greedy {
			lmin := int64(min_) * int64(rbody.min)
			lmax := int64(max_) * int64(rbody.max)
			if lmin < javaIntegerMaxValue {
				body = rbody.body
				min_ = int(lmin)
				if lmax >= javaIntegerMaxValue {
					max_ = javaIntegerMaxValue
				} else {
					max_ = int(lmax)
				}
			}
		}
	}
	if min_ == 1 && max_ == 1 {
		return body
	}
	greedy := r.greedy || min_ == max_
	if body.Equals(r.body) && min_ == r.min && max_ == r.max && greedy == r.greedy {
		return r
	}
	return NewRegExpTree_Repetition(body, min_, max_, greedy).Simplify(flags)
}

// IsCaseSensitive delegates to the body.
//
// Java: Repetition.isCaseSensitive()
func (r *RegExpTree_Repetition) IsCaseSensitive() bool {
	return r.body.IsCaseSensitive()
}

// ContainsAnchor delegates to the body.
//
// Java: Repetition.containsAnchor()
func (r *RegExpTree_Repetition) ContainsAnchor() bool {
	return r.body.ContainsAnchor()
}

// Java: RegExpTree.hasCapturingGroup()
func (r *RegExpTree_Repetition) HasCapturingGroup() bool {
	return r.body.HasCapturingGroup()
}

// NumCapturingGroups delegates to the body.
//
// Java: Repetition.numCapturingGroups()
func (r *RegExpTree_Repetition) NumCapturingGroups() int {
	return r.body.NumCapturingGroups()
}

// Children is [body].
//
// Java: Repetition.children()
func (r *RegExpTree_Repetition) Children() []RegExpTree {
	return []RegExpTree{r.body}
}

// AppendBodySourceCode emits the body, parenthesized when needed.
//
// Java: Repetition.appendBodySourceCode(StringBuilder)
func (r *RegExpTree_Repetition) AppendBodySourceCode(sb *jutil.StringBuilder) {
	_, isAlternation := r.body.(*RegExpTree_Alternation)
	_, isConcatenation := r.body.(*RegExpTree_Concatenation)
	_, isRepetition := r.body.(*RegExpTree_Repetition)
	text, isText := r.body.(*RegExpTree_Text)
	if isAlternation ||
		isConcatenation ||
		isRepetition ||
		(isText && jutil.StringLength(text.text) > 1) {
		sb.Append("(?:")
		r.body.AppendSourceCode(sb)
		sb.AppendChar(')')
	} else {
		r.body.AppendSourceCode(sb)
	}
}

// RegExpTree_Repetition_SuffixLen is the length of the suffix that
// AppendSourceCode renders for min and max.
//
// Java: Repetition.suffixLen(int, int)
func RegExpTree_Repetition_SuffixLen(min_ int, max_ int) int {
	// This mirrors the branches that renders a suffix in appendSourceCode below.
	if max_ == javaIntegerMaxValue {
		switch min_ {
		case 0, // *
			1: // +
			return 1
		default: // {3,}
			return 3 + RegExpTree_Repetition_NumDecimalDigits(min_)
		}
	}
	if min_ == 0 && max_ == 1 {
		return 1 // ?
	}
	if min_ == max_ {
		if min_ == 1 {
			return 0 // No suffix needed for {1}.
		}
		return 2 + RegExpTree_Repetition_NumDecimalDigits(min_) // {4}
	}
	return 3 + RegExpTree_Repetition_NumDecimalDigits(min_) + RegExpTree_Repetition_NumDecimalDigits(max_) // {2,7}
}

// RegExpTree_Repetition_NumDecimalDigits is the number of decimal digits of
// a non-negative n.
//
// Java: Repetition.numDecimalDigits(int)
func RegExpTree_Repetition_NumDecimalDigits(n int) int {
	if n < 0 {
		// Negative values should not be passed in.
		panic(jutil.NewAssertionErrorNoMessage())
		// If changing this code to support negative values,
		// Integer.MIN_VALUE is a corner-case..
	}
	nDigits := 1
	for n >= 10 {
		nDigits++
		n /= 10
	}
	return nDigits
}

// AppendSourceCode emits the body and the shortest suffix.
//
// Java: Repetition.appendSourceCode(StringBuilder)
func (r *RegExpTree_Repetition) AppendSourceCode(sb *jutil.StringBuilder) {
	bodyStart := sb.Length()
	r.AppendBodySourceCode(sb)
	bodyEnd := sb.Length()
	bodyLen := bodyEnd - bodyStart
	min_ := r.min
	max_ := r.max
	if (min_ >= 2 && max_ == javaIntegerMaxValue) || max_-min_ <= 1 {
		var expanded int
		// If min == max then we want to try expanding to the limit and
		// attach the empty suffix, which is equivalent to min = max = 1,
		// i.e. /a/ vs /a{1}/.
		if min_ == max_ ||
			// Give aa+ preference over aaa*.
			max_ == javaIntegerMaxValue {
			expanded = min_ - 1
		} else {
			expanded = min_
		}
		expandedMin := min_ - expanded
		expandedMax := max_ - expanded
		if max_ == javaIntegerMaxValue {
			expandedMax = max_
		}
		suffixLen := RegExpTree_Repetition_SuffixLen(min_, max_)
		expandedSuffixLen := RegExpTree_Repetition_SuffixLen(expandedMin, expandedMax)
		// Go-only: language-adapter — Java int arithmetic (wraps at 32 bits).
		if int32(bodyLen)*int32(expanded)+int32(expandedSuffixLen) < int32(suffixLen) && !r.body.HasCapturingGroup() {
			// a{2} -> aa
			// a{2,} -> aa+
			// a{2,3} -> aaa?
			for expanded--; expanded >= 0; expanded-- {
				sb.AppendCharSequenceRange(sb, bodyStart, bodyEnd)
			}
			min_ = expandedMin
			max_ = expandedMax
		}
	}

	if max_ == javaIntegerMaxValue {
		switch min_ {
		case 0:
			sb.AppendChar('*')
		case 1:
			sb.AppendChar('+')
		default:
			sb.AppendChar('{').AppendInt(int64(min_)).Append(",}")
		}
	} else if min_ == 0 && max_ == 1 {
		sb.AppendChar('?')
	} else if min_ == max_ {
		if min_ != 1 {
			sb.AppendChar('{').AppendInt(int64(min_)).AppendChar('}')
		}
	} else {
		sb.AppendChar('{').AppendInt(int64(min_)).AppendChar(',').AppendInt(int64(max_)).AppendChar('}')
	}
	if !r.greedy {
		sb.AppendChar('?')
	}
}

// AppendDebugInfo emits the bounds and greediness.
//
// Java: Repetition.appendDebugInfo(StringBuilder)
func (r *RegExpTree_Repetition) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.Append(" min=").AppendInt(int64(r.min)).Append(", max=").AppendInt(int64(r.max))
	if !r.greedy {
		sb.Append("  not_greedy")
	}
}

// Equals compares body, bounds and greediness.
//
// Java: Repetition.equals(Object)
func (r *RegExpTree_Repetition) Equals(o any) bool {
	that, ok := o.(*RegExpTree_Repetition)
	if !ok {
		return false
	}
	return r.body.Equals(that.body) &&
		r.min == that.min &&
		r.max == that.max &&
		r.greedy == that.greedy
}

// HashCode hashes body, bounds and greediness.
//
// Java: Repetition.hashCode()
func (r *RegExpTree_Repetition) HashCode() int32 {
	greedy := int32(0)
	if r.greedy {
		greedy = 1
	}
	return int32(r.min) + 31*(int32(r.max)+31*(greedy+31*r.body.HashCode()))
}

// Java: RegExpTree.toString()
func (r *RegExpTree_Repetition) String() string { return regExpTreeString(r) }

// RegExpTree_Alternation represents the possibilities ["foo", "bar" ] for a
// RegExp /foo|bar/.
//
// Java: public static final class RegExpTree.Alternation
type RegExpTree_Alternation struct {
	alternatives []RegExpTree
}

// NewRegExpTree_Alternation constructs an Alternation over a copy of
// alternatives.
//
// Java: Alternation(List<? extends RegExpTree>)
func NewRegExpTree_Alternation(alternatives []RegExpTree) *RegExpTree_Alternation {
	return &RegExpTree_Alternation{
		alternatives: slices.Clone(alternatives), // Java: ImmutableList.copyOf
	}
}

// Simplify flattens, deduplicates and collapses character alternatives.
//
// Java: Alternation.simplify(String)
func (a *RegExpTree_Alternation) Simplify(flags string) RegExpTree {
	alternatives := []RegExpTree{}
	for _, alternative := range a.alternatives {
		alternative = alternative.Simplify(flags)
		if alternation, ok := alternative.(*RegExpTree_Alternation); ok {
			alternatives = append(alternatives, alternation.alternatives...)
		} else {
			alternatives = append(alternatives, alternative)
		}
	}
	// Remove duplicates
	var last RegExpTree
	for i := 0; i < len(alternatives); { // Java: Iterator with it.remove()
		alternative := alternatives[i]
		if alternative.Equals(RegExpTree_NEVER_MATCHES) {
			i++
			continue
		}
		if alternative.Equals(last) && !alternative.HasCapturingGroup() {
			alternatives = slices.Delete(alternatives, i, i+1)
		} else {
			last = alternative
			i++
		}
	}
	// Collapse character alternatives into character sets.
	for i, n := 0, len(alternatives); i < n; i++ {
		alternative := alternatives[i]
		_, isCharset := alternative.(*RegExpTree_Charset)
		if regExpTree_isSingleCharText(alternative) || isCharset {
			end := i
			nCharsets := 0
			for end < n {
				follower := alternatives[end]
				if _, ok := follower.(*RegExpTree_Charset); ok {
					nCharsets++
				} else if !regExpTree_isSingleCharText(follower) {
					break
				}
				end++
			}
			if end-i >= 3 || (nCharsets != 0 && end-i >= 2) {
				members := make([]int, end-i-nCharsets)
				memberIdx := 0
				chars := CharRanges_EMPTY
				ieExplicits := CharRanges_EMPTY
				charAlternatives := alternatives[i:end] // Java: alternatives.subList(i, end)
				for _, charAlternative := range charAlternatives {
					if text, ok := charAlternative.(*RegExpTree_Text); ok {
						ch := jutil.StringCharAt(text.text, 0)
						members[memberIdx] = int(ch)
						memberIdx++
						if RegExpTree_IE_SPEC_ERRORS.Contains(int(ch)) {
							ieExplicits = ieExplicits.Union(CharRanges_Inclusive(int(ch), int(ch)))
						}
					} else if cs, ok := charAlternative.(*RegExpTree_Charset); ok {
						chars = chars.Union(cs.ranges)
						ieExplicits = ieExplicits.Union(cs.ieExplicits)
					}
				}
				chars = chars.Union(CharRanges_WithMembers(members...))
				// Java: charAlternatives.clear(); charAlternatives.add(...)
				alternatives = slices.Replace(alternatives, i, end,
					NewRegExpTree_Charset(chars, ieExplicits).Simplify(flags))
				n = len(alternatives)
			}
		}
	}
	switch len(alternatives) {
	case 0:
		return RegExpTree_Empty_INSTANCE
	case 1:
		return alternatives[0]
	case 2:
		if _, ok := alternatives[1].(*RegExpTree_Empty); ok { // (?:a|) -> a?
			return NewRegExpTree_Repetition(alternatives[0], 0, 1, true)
		} else if _, ok := alternatives[0].(*RegExpTree_Empty); ok {
			return NewRegExpTree_Repetition(alternatives[1], 0, 1, false)
		}
	}
	// TODO: maybe pull out common prefix or suffix
	if regExpTree_listEquals(alternatives, a.alternatives) {
		return a
	}
	return NewRegExpTree_Alternation(alternatives)
}

// regExpTree_isSingleCharText is `x instanceof Text text &&
// text.text.length() == 1`.
//
// Go-only: language-adapter — pattern-matching instanceof as a helper.
func regExpTree_isSingleCharText(t RegExpTree) bool {
	text, ok := t.(*RegExpTree_Text)
	return ok && jutil.StringLength(text.text) == 1
}

// IsCaseSensitive reports whether any alternative is.
//
// Java: Alternation.isCaseSensitive()
func (a *RegExpTree_Alternation) IsCaseSensitive() bool {
	for _, alternative := range a.alternatives {
		if alternative.IsCaseSensitive() {
			return true
		}
	}
	return false
}

// ContainsAnchor reports whether any alternative does.
//
// Java: Alternation.containsAnchor()
func (a *RegExpTree_Alternation) ContainsAnchor() bool {
	for _, alternative := range a.alternatives {
		if alternative.ContainsAnchor() {
			return true
		}
	}
	return false
}

// Java: RegExpTree.hasCapturingGroup()
func (a *RegExpTree_Alternation) HasCapturingGroup() bool {
	for _, alternative := range a.alternatives {
		if alternative.HasCapturingGroup() {
			return true
		}
	}
	return false
}

// NumCapturingGroups sums over the alternatives.
//
// Java: Alternation.numCapturingGroups()
func (a *RegExpTree_Alternation) NumCapturingGroups() int {
	n := 0
	for _, alternative := range a.alternatives {
		n += alternative.NumCapturingGroups()
	}
	return n
}

// Children returns the alternatives.
//
// Java: Alternation.children()
func (a *RegExpTree_Alternation) Children() []RegExpTree {
	return a.alternatives
}

// AppendSourceCode emits the alternatives separated by '|'.
//
// Java: Alternation.appendSourceCode(StringBuilder)
func (a *RegExpTree_Alternation) AppendSourceCode(sb *jutil.StringBuilder) {
	for i, n := 0, len(a.alternatives); i < n; i++ {
		if i != 0 {
			sb.AppendChar('|')
		}
		a.alternatives[i].AppendSourceCode(sb)
	}
}

// AppendDebugInfo emits nothing besides children.
//
// Java: Alternation.appendDebugInfo(StringBuilder)
func (a *RegExpTree_Alternation) AppendDebugInfo(sb *jutil.StringBuilder) {
	// Nothing besides children.
}

// Equals compares the alternatives.
//
// Java: Alternation.equals(Object)
func (a *RegExpTree_Alternation) Equals(o any) bool {
	if alternation, ok := o.(*RegExpTree_Alternation); ok {
		return a == alternation || regExpTree_listEquals(a.alternatives, alternation.alternatives)
	}
	return false
}

// HashCode hashes the alternatives.
//
// Java: Alternation.hashCode()
func (a *RegExpTree_Alternation) HashCode() int32 {
	return 0x51b57cd1 ^ regExpTree_listHashCode(a.alternatives)
}

// Java: RegExpTree.toString()
func (a *RegExpTree_Alternation) String() string { return regExpTreeString(a) }

// regExpTree_listEquals is List.equals on lists of RegExpTree.
//
// Go-only: language-adapter — java.util.List.equals.
func regExpTree_listEquals(a, b []RegExpTree) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equals(b[i]) {
			return false
		}
	}
	return true
}

// regExpTree_listHashCode is List.hashCode on lists of RegExpTree.
//
// Go-only: language-adapter — java.util.List.hashCode.
func regExpTree_listHashCode(l []RegExpTree) int32 {
	hashCode := int32(1)
	for _, e := range l {
		hashCode = 31*hashCode + e.HashCode()
	}
	return hashCode
}

// RegExpTree_NEVER_MATCHES is (?!), which matches nothing.
var RegExpTree_NEVER_MATCHES RegExpTree = NewRegExpTree_LookaheadAssertion(RegExpTree_Empty_INSTANCE, false)

// RegExpTree_LookaheadAssertion represents a lookahead assertion such as
// (?=...) or (?!...).
//
// Java: private static final class RegExpTree.LookaheadAssertion
type RegExpTree_LookaheadAssertion struct {
	body     RegExpTree
	positive bool
}

// NewRegExpTree_LookaheadAssertion constructs a LookaheadAssertion.
//
// Java: LookaheadAssertion(RegExpTree, boolean)
func NewRegExpTree_LookaheadAssertion(body RegExpTree, positive bool) *RegExpTree_LookaheadAssertion {
	return &RegExpTree_LookaheadAssertion{
		body:     body,
		positive: positive,
	}
}

// Simplify simplifies the body; (?=) is always true.
//
// Java: LookaheadAssertion.simplify(String)
func (l *RegExpTree_LookaheadAssertion) Simplify(flags string) RegExpTree {
	simpleBody := l.body.Simplify(flags)
	if _, ok := simpleBody.(*RegExpTree_Empty); ok {
		if l.positive { // Always true
			return simpleBody
		}
	}
	return NewRegExpTree_LookaheadAssertion(simpleBody, l.positive)
}

// IsCaseSensitive delegates to the body.
//
// Java: LookaheadAssertion.isCaseSensitive()
func (l *RegExpTree_LookaheadAssertion) IsCaseSensitive() bool {
	return l.body.IsCaseSensitive()
}

// ContainsAnchor delegates to the body.
//
// Java: LookaheadAssertion.containsAnchor()
func (l *RegExpTree_LookaheadAssertion) ContainsAnchor() bool {
	return l.body.ContainsAnchor()
}

// Java: RegExpTree.hasCapturingGroup()
func (l *RegExpTree_LookaheadAssertion) HasCapturingGroup() bool {
	return l.body.HasCapturingGroup()
}

// NumCapturingGroups delegates to the body.
//
// Java: LookaheadAssertion.numCapturingGroups()
func (l *RegExpTree_LookaheadAssertion) NumCapturingGroups() int {
	return l.body.NumCapturingGroups()
}

// Children is [body].
//
// Java: LookaheadAssertion.children()
func (l *RegExpTree_LookaheadAssertion) Children() []RegExpTree {
	return []RegExpTree{l.body}
}

// AppendSourceCode emits (?=...) or (?!...).
//
// Java: LookaheadAssertion.appendSourceCode(StringBuilder)
func (l *RegExpTree_LookaheadAssertion) AppendSourceCode(sb *jutil.StringBuilder) {
	if l.positive {
		sb.Append("(?=")
	} else {
		sb.Append("(?!")
	}
	l.body.AppendSourceCode(sb)
	sb.AppendChar(')')
}

// AppendDebugInfo emits positive or negative.
//
// Java: LookaheadAssertion.appendDebugInfo(StringBuilder)
func (l *RegExpTree_LookaheadAssertion) AppendDebugInfo(sb *jutil.StringBuilder) {
	if l.positive {
		sb.Append("positive")
	} else {
		sb.Append("negative")
	}
}

// Equals compares polarity and body.
//
// Java: LookaheadAssertion.equals(Object)
func (l *RegExpTree_LookaheadAssertion) Equals(o any) bool {
	that, ok := o.(*RegExpTree_LookaheadAssertion)
	if !ok {
		return false
	}
	return l.positive == that.positive && l.body.Equals(that.body)
}

// HashCode hashes the body.
//
// Java: LookaheadAssertion.hashCode()
func (l *RegExpTree_LookaheadAssertion) HashCode() int32 {
	return 0x723aba9 ^ l.body.HashCode()
}

// Java: RegExpTree.toString()
func (l *RegExpTree_LookaheadAssertion) String() string { return regExpTreeString(l) }

// RegExpTree_LookbehindAssertion represents a lookbehind assertion such as
// (?<=...) or (?<!...).
//
// Java: public static final class RegExpTree.LookbehindAssertion
type RegExpTree_LookbehindAssertion struct {
	body     RegExpTree
	positive bool
}

// NewRegExpTree_LookbehindAssertion constructs a LookbehindAssertion.
//
// Java: LookbehindAssertion(RegExpTree, boolean)
func NewRegExpTree_LookbehindAssertion(body RegExpTree, positive bool) *RegExpTree_LookbehindAssertion {
	return &RegExpTree_LookbehindAssertion{
		body:     body,
		positive: positive,
	}
}

// Simplify simplifies the body; (?<=) is always true.
//
// Java: LookbehindAssertion.simplify(String)
func (l *RegExpTree_LookbehindAssertion) Simplify(flags string) RegExpTree {
	simpleBody := l.body.Simplify(flags)
	if _, ok := simpleBody.(*RegExpTree_Empty); ok {
		if l.positive { // Always true
			return simpleBody
		}
	}
	return NewRegExpTree_LookbehindAssertion(simpleBody, l.positive)
}

// IsCaseSensitive delegates to the body.
//
// Java: LookbehindAssertion.isCaseSensitive()
func (l *RegExpTree_LookbehindAssertion) IsCaseSensitive() bool {
	return l.body.IsCaseSensitive()
}

// ContainsAnchor delegates to the body.
//
// Java: LookbehindAssertion.containsAnchor()
func (l *RegExpTree_LookbehindAssertion) ContainsAnchor() bool {
	return l.body.ContainsAnchor()
}

// Java: RegExpTree.hasCapturingGroup()
func (l *RegExpTree_LookbehindAssertion) HasCapturingGroup() bool {
	return l.body.HasCapturingGroup()
}

// NumCapturingGroups delegates to the body.
//
// Java: LookbehindAssertion.numCapturingGroups()
func (l *RegExpTree_LookbehindAssertion) NumCapturingGroups() int {
	return l.body.NumCapturingGroups()
}

// Children is [body].
//
// Java: LookbehindAssertion.children()
func (l *RegExpTree_LookbehindAssertion) Children() []RegExpTree {
	return []RegExpTree{l.body}
}

// AppendSourceCode emits (?<=...) or (?<!...).
//
// Java: LookbehindAssertion.appendSourceCode(StringBuilder)
func (l *RegExpTree_LookbehindAssertion) AppendSourceCode(sb *jutil.StringBuilder) {
	if l.positive {
		sb.Append("(?<=")
	} else {
		sb.Append("(?<!")
	}
	l.body.AppendSourceCode(sb)
	sb.AppendChar(')')
}

// AppendDebugInfo emits positive or negative.
//
// Java: LookbehindAssertion.appendDebugInfo(StringBuilder)
func (l *RegExpTree_LookbehindAssertion) AppendDebugInfo(sb *jutil.StringBuilder) {
	if l.positive {
		sb.Append("positive")
	} else {
		sb.Append("negative")
	}
}

// Equals compares polarity and body.
//
// Java: LookbehindAssertion.equals(Object)
func (l *RegExpTree_LookbehindAssertion) Equals(o any) bool {
	that, ok := o.(*RegExpTree_LookbehindAssertion)
	if !ok {
		return false
	}
	return l.positive == that.positive && l.body.Equals(that.body)
}

// HashCode hashes the body.
//
// Java: LookbehindAssertion.hashCode()
func (l *RegExpTree_LookbehindAssertion) HashCode() int32 {
	return 0x723aba9 ^ l.body.HashCode()
}

// Java: RegExpTree.toString()
func (l *RegExpTree_LookbehindAssertion) String() string { return regExpTreeString(l) }

// RegExpTree_CapturingGroup represents a capturing group such as (asdf).
//
// Java: public static final class RegExpTree.CapturingGroup
type RegExpTree_CapturingGroup struct {
	body RegExpTree
}

// NewRegExpTree_CapturingGroup constructs a CapturingGroup.
//
// Java: CapturingGroup(RegExpTree)
func NewRegExpTree_CapturingGroup(body RegExpTree) *RegExpTree_CapturingGroup {
	return &RegExpTree_CapturingGroup{body: body}
}

// Simplify simplifies the body.
//
// Java: CapturingGroup.simplify(String)
func (c *RegExpTree_CapturingGroup) Simplify(flags string) RegExpTree {
	return NewRegExpTree_CapturingGroup(c.body.Simplify(flags))
}

// IsCaseSensitive delegates to the body.
//
// Java: CapturingGroup.isCaseSensitive()
func (c *RegExpTree_CapturingGroup) IsCaseSensitive() bool {
	return c.body.IsCaseSensitive()
}

// ContainsAnchor delegates to the body.
//
// Java: CapturingGroup.containsAnchor()
func (c *RegExpTree_CapturingGroup) ContainsAnchor() bool {
	return c.body.ContainsAnchor()
}

// Java: RegExpTree.hasCapturingGroup()
func (c *RegExpTree_CapturingGroup) HasCapturingGroup() bool {
	return true
}

// NumCapturingGroups is one plus the body's.
//
// Java: CapturingGroup.numCapturingGroups()
func (c *RegExpTree_CapturingGroup) NumCapturingGroups() int {
	return 1 + c.body.NumCapturingGroups()
}

// Children is [body].
//
// Java: CapturingGroup.children()
func (c *RegExpTree_CapturingGroup) Children() []RegExpTree {
	return []RegExpTree{c.body}
}

// AppendSourceCode emits (...).
//
// Java: CapturingGroup.appendSourceCode(StringBuilder)
func (c *RegExpTree_CapturingGroup) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.AppendChar('(')
	c.body.AppendSourceCode(sb)
	sb.AppendChar(')')
}

// AppendDebugInfo emits nothing besides children.
//
// Java: CapturingGroup.appendDebugInfo(StringBuilder)
func (c *RegExpTree_CapturingGroup) AppendDebugInfo(sb *jutil.StringBuilder) {
	// Nothing besides children.
}

// Equals compares the body.
//
// Java: CapturingGroup.equals(Object)
func (c *RegExpTree_CapturingGroup) Equals(o any) bool {
	capturingGroup, ok := o.(*RegExpTree_CapturingGroup)
	return ok && c.body.Equals(capturingGroup.body)
}

// HashCode hashes the body.
//
// Java: CapturingGroup.hashCode()
func (c *RegExpTree_CapturingGroup) HashCode() int32 {
	return 0x55781738 ^ c.body.HashCode()
}

// Java: RegExpTree.toString()
func (c *RegExpTree_CapturingGroup) String() string { return regExpTreeString(c) }

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

// Simplify simplifies the body.
//
// Java: NamedCaptureGroup.simplify(String)
func (n *RegExpTree_NamedCaptureGroup) Simplify(flags string) RegExpTree {
	return NewRegExpTree_NamedCaptureGroup(n.body.Simplify(flags), n.name)
}

// IsCaseSensitive delegates to the body.
//
// Java: NamedCaptureGroup.isCaseSensitive()
func (n *RegExpTree_NamedCaptureGroup) IsCaseSensitive() bool {
	return n.body.IsCaseSensitive()
}

// ContainsAnchor delegates to the body.
//
// Java: NamedCaptureGroup.containsAnchor()
func (n *RegExpTree_NamedCaptureGroup) ContainsAnchor() bool {
	return n.body.ContainsAnchor()
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

// Children is [body].
//
// Java: NamedCaptureGroup.children()
func (n *RegExpTree_NamedCaptureGroup) Children() []RegExpTree {
	return []RegExpTree{n.body}
}

// AppendSourceCode emits (?<name>...).
//
// Java: NamedCaptureGroup.appendSourceCode(StringBuilder)
func (n *RegExpTree_NamedCaptureGroup) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.Append("(?<")
	sb.Append(n.name)
	sb.AppendChar('>')
	n.body.AppendSourceCode(sb)
	sb.AppendChar(')')
}

// AppendDebugInfo emits the name.
//
// Java: NamedCaptureGroup.appendDebugInfo(StringBuilder)
func (n *RegExpTree_NamedCaptureGroup) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.Append(" name=").Append(n.name)
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

// RegExpTree_Concatenation represents a series of nodes chained one after
// another such as (?:...)[a-z]*(...).
//
// Java: public static final class RegExpTree.Concatenation
type RegExpTree_Concatenation struct {
	elements []RegExpTree
}

// NewRegExpTree_ConcatenationOfPair constructs the concatenation of a and b.
//
// Java: Concatenation(RegExpTree, RegExpTree)
func NewRegExpTree_ConcatenationOfPair(a RegExpTree, b RegExpTree) *RegExpTree_Concatenation {
	return &RegExpTree_Concatenation{elements: []RegExpTree{a, b}}
}

// NewRegExpTree_Concatenation constructs a concatenation over a copy of
// elements.
//
// Java: Concatenation(List<? extends RegExpTree>)
func NewRegExpTree_Concatenation(elements []RegExpTree) *RegExpTree_Concatenation {
	return &RegExpTree_Concatenation{
		elements: slices.Clone(elements), // Java: ImmutableList.copyOf
	}
}

// Simplify flattens nested concatenations and folds adjacent texts and
// repetitions.
//
// Java: Concatenation.simplify(String)
func (c *RegExpTree_Concatenation) Simplify(flags string) RegExpTree {
	s := NewRegExpTree_Concatenation_Simplifier(flags)
	for _, element := range c.elements {
		s.Simplify(element.Simplify(flags))
	}

	switch len(s.simplified) {
	case 0:
		return RegExpTree_Empty_INSTANCE
	case 1:
		return s.simplified[0]
	default:
		return NewRegExpTree_Concatenation(s.simplified)
	}
}

// RegExpTree_Concatenation_Simplifier accumulates simplified elements.
//
// Java: local class Simplifier in Concatenation.simplify(String); the
// captured flags are a constructor parameter.
type RegExpTree_Concatenation_Simplifier struct {
	flags      string
	simplified []RegExpTree
}

// NewRegExpTree_Concatenation_Simplifier constructs a Simplifier.
//
// Java: new Simplifier() (inside Concatenation.simplify)
func NewRegExpTree_Concatenation_Simplifier(flags string) *RegExpTree_Concatenation_Simplifier {
	s := &RegExpTree_Concatenation_Simplifier{}
	s.flags = flags
	s.simplified = []RegExpTree{}
	return s
}

// Simplify appends t, flattening and folding pairwise.
//
// Java: Simplifier.simplify(RegExpTree)
func (s *RegExpTree_Concatenation_Simplifier) Simplify(t RegExpTree) {
	if concatenation, ok := t.(*RegExpTree_Concatenation); ok {
		for _, child := range concatenation.elements {
			s.Simplify(child)
		}
	} else if _, ok := t.(*RegExpTree_Empty); ok {
		// Do nothing
	} else {
		lastIndex := len(s.simplified) - 1
		if lastIndex >= 0 {
			pairwise := s.SimplifyPairwise(s.simplified[lastIndex], t)
			if pairwise != nil {
				s.simplified[lastIndex] = pairwise
				return
			}
		}
		s.simplified = append(s.simplified, t)
	}
}

// SimplifyPairwise folds two adjacent elements, or returns nil.
//
// Java: Simplifier.simplifyPairwise(RegExpTree, RegExpTree)
func (s *RegExpTree_Concatenation_Simplifier) SimplifyPairwise(before RegExpTree, after RegExpTree) RegExpTree {
	beforeText, ok1 := before.(*RegExpTree_Text)
	afterText, ok2 := after.(*RegExpTree_Text)
	if ok1 && ok2 {
		return NewRegExpTree_Text(beforeText.text + afterText.text).Simplify(s.flags)
	}
	// Fold adjacent repetitions.
	beforeMin := 1
	beforeMax := 1
	beforeBody := before
	beforeGreedy := false
	if r, ok := before.(*RegExpTree_Repetition); ok {
		beforeMin = r.min
		beforeMax = r.max
		beforeBody = r.body
		beforeGreedy = r.greedy
	}
	afterMin := 1
	afterMax := 1
	afterBody := after
	afterGreedy := false
	if r, ok := after.(*RegExpTree_Repetition); ok {
		afterMin = r.min
		afterMax = r.max
		afterBody = r.body
		afterGreedy = r.greedy
	}
	if beforeBody.Equals(afterBody) && !beforeBody.HasCapturingGroup() {
		lmin := int64(beforeMin) + int64(afterMin)
		lmax := int64(beforeMax) + int64(afterMax)
		if lmin < javaIntegerMaxValue {
			min_ := int(lmin)
			max_ := javaIntegerMaxValue
			if lmax < javaIntegerMaxValue {
				max_ = int(lmax)
			}
			return NewRegExpTree_Repetition(
				beforeBody, min_, max_, beforeGreedy || afterGreedy || min_ == max_)
		}
	}
	return nil
}

// IsCaseSensitive reports whether any element is.
//
// Java: Concatenation.isCaseSensitive()
func (c *RegExpTree_Concatenation) IsCaseSensitive() bool {
	for _, element := range c.elements {
		if element.IsCaseSensitive() {
			return true
		}
	}
	return false
}

// ContainsAnchor reports whether any element does.
//
// Java: Concatenation.containsAnchor()
func (c *RegExpTree_Concatenation) ContainsAnchor() bool {
	for _, element := range c.elements {
		if element.ContainsAnchor() {
			return true
		}
	}
	return false
}

// Java: RegExpTree.hasCapturingGroup()
func (c *RegExpTree_Concatenation) HasCapturingGroup() bool {
	for _, element := range c.elements {
		if element.HasCapturingGroup() {
			return true
		}
	}
	return false
}

// NumCapturingGroups sums over the elements.
//
// Java: Concatenation.numCapturingGroups()
func (c *RegExpTree_Concatenation) NumCapturingGroups() int {
	n := 0
	for _, element := range c.elements {
		n += element.NumCapturingGroups()
	}
	return n
}

// Children returns the elements.
//
// Java: Concatenation.children()
func (c *RegExpTree_Concatenation) Children() []RegExpTree {
	return c.elements
}

// AppendSourceCode emits the elements, guarding against digits that would
// bleed into a preceding back reference or brace.
//
// Java: Concatenation.appendSourceCode(StringBuilder)
func (c *RegExpTree_Concatenation) AppendSourceCode(sb *jutil.StringBuilder) {
	// True if the last content written might consume
	// decimal digits written subsequently.
	digitsMightBleed := false
	for _, element := range c.elements {
		parenthesize := false
		_, isAlternation := element.(*RegExpTree_Alternation)
		_, isConcatenation := element.(*RegExpTree_Concatenation)
		if isAlternation || isConcatenation {
			parenthesize = true
		}
		if parenthesize {
			sb.Append("(?:")
			element.AppendSourceCode(sb)
			sb.AppendChar(')')
		} else {
			start := sb.Length()
			element.AppendSourceCode(sb)
			if digitsMightBleed && sb.Length() > start {
				firstChar := sb.CharAt(start)
				if '0' <= firstChar && firstChar <= '9' {
					// Bleeding happened.
					// If the last character would be ambiguous
					// with a repetition, escape it.
					if sb.CharAt(start-1) == '{' {
						// Concatenation from optimization of
						// /{(?:0,}/ -> /\{0,}/
						sb.InsertChar(start-1, '\\')
					} else {
						// Or parenthesize otherwise.
						// Concatenation from optimization of
						// /(.)\1(?:0)/ -> /(.)\1(?:0)/.
						sb.Insert(start, "(?:").AppendChar(')')
					}
				}
			}
		}
		backReference, isBackReference := element.(*RegExpTree_BackReference)
		text, isText := element.(*RegExpTree_Text)
		digitsMightBleed =
			// \1(?:0) bleeds if there are 10 or more
			// capturing groups preceding.
			(isBackReference && backReference.groupIndex < 10) ||
				// foo{(?:10}) bleeds.
				(isText && strings.HasSuffix(text.text, "{"))
	}
}

// AppendDebugInfo emits nothing besides children.
//
// Java: Concatenation.appendDebugInfo(StringBuilder)
func (c *RegExpTree_Concatenation) AppendDebugInfo(sb *jutil.StringBuilder) {
	// Nothing besides children.
}

// Equals compares the elements.
//
// Java: Concatenation.equals(Object)
func (c *RegExpTree_Concatenation) Equals(o any) bool {
	concatenation, ok := o.(*RegExpTree_Concatenation)
	return ok && regExpTree_listEquals(c.elements, concatenation.elements)
}

// HashCode hashes the elements.
//
// Java: Concatenation.hashCode()
func (c *RegExpTree_Concatenation) HashCode() int32 {
	return 0x20997e3e ^ regExpTree_listHashCode(c.elements)
}

// Java: RegExpTree.toString()
func (c *RegExpTree_Concatenation) String() string { return regExpTreeString(c) }

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

// Package regex is the port of Closure Compiler's
// com.google.javascript.jscomp.regex package: the JavaScript regular
// expression AST (RegExpTree and its node classes), the sparse code-unit set
// CharRanges and the ES5 Canonicalize tables (CaseCanonicalize).
//
// Patterns are parsed as UTF-16 code units, like Java chars, so offsets in
// exception messages and the code units that appendSourceCode escapes are
// exactly Java's. reg_exp_tree_oracle_test.go diffs this package against the
// Java classes (tools/oracle/java/RegExpDump.java).
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
	AppendSourceCode(sb *strings.Builder)
	AppendDebugInfo(sb *strings.Builder)
	String() string
	Equals(o any) bool
	HashCode() int32
}

// regExpTreeString renders the regular expression as a literal, e.g. "/foo/".
//
// Java: RegExpTree.toString()
func regExpTreeString(t RegExpTree) string {
	var sb strings.Builder
	sb.WriteByte('/')
	t.AppendSourceCode(&sb)
	// Don't emit a regular expression that looks like a line comment start.
	if sb.Len() == 1 {
		sb.WriteString("(?:)")
	}
	sb.WriteByte('/')
	return sb.String()
}

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

func (r *RegExpTree_Repetition) IsCaseSensitive() bool   { return r.body.IsCaseSensitive() }
func (r *RegExpTree_Repetition) ContainsAnchor() bool    { return r.body.ContainsAnchor() }
func (r *RegExpTree_Repetition) HasCapturingGroup() bool { return r.body.HasCapturingGroup() }
func (r *RegExpTree_Repetition) NumCapturingGroups() int { return r.body.NumCapturingGroups() }
func (r *RegExpTree_Repetition) Children() []RegExpTree  { return []RegExpTree{r.body} }

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

// AppendBodySourceCode emits the body, parenthesized when needed.
//
// Java: Repetition.appendBodySourceCode(StringBuilder)
func (r *RegExpTree_Repetition) AppendBodySourceCode(sb *strings.Builder) {
	_, isAlternation := r.body.(*RegExpTree_Alternation)
	_, isConcatenation := r.body.(*RegExpTree_Concatenation)
	_, isRepetition := r.body.(*RegExpTree_Repetition)
	_, isText := r.body.(*RegExpTree_Text)
	if isAlternation ||
		isConcatenation ||
		isRepetition ||
		(isText && !regExpTree_isSingleCharText(r.body) && r.body.(*RegExpTree_Text).text != "") {
		sb.WriteString("(?:")
		r.body.AppendSourceCode(sb)
		sb.WriteByte(')')
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
func (r *RegExpTree_Repetition) AppendSourceCode(sb *strings.Builder) {
	bodyStart := sb.Len()
	r.AppendBodySourceCode(sb)
	bodyLen := jutil.StringLength(sb.String()[bodyStart:])
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
			body := sb.String()[bodyStart:]
			for expanded--; expanded >= 0; expanded-- {
				sb.WriteString(body)
			}
			min_ = expandedMin
			max_ = expandedMax
		}
	}

	if max_ == javaIntegerMaxValue {
		switch min_ {
		case 0:
			sb.WriteByte('*')
		case 1:
			sb.WriteByte('+')
		default:
			sb.WriteByte('{')
			writeDecimalInt(sb, min_)
			sb.WriteString(",}")
		}
	} else if min_ == 0 && max_ == 1 {
		sb.WriteByte('?')
	} else if min_ == max_ {
		if min_ != 1 {
			sb.WriteByte('{')
			writeDecimalInt(sb, min_)
			sb.WriteByte('}')
		}
	} else {
		sb.WriteByte('{')
		writeDecimalInt(sb, min_)
		sb.WriteByte(',')
		writeDecimalInt(sb, max_)
		sb.WriteByte('}')
	}
	if !r.greedy {
		sb.WriteByte('?')
	}
}

// AppendDebugInfo emits the bounds and greediness.
//
// Java: Repetition.appendDebugInfo(StringBuilder)
func (r *RegExpTree_Repetition) AppendDebugInfo(sb *strings.Builder) {
	sb.WriteString(" min=")
	writeDecimalInt(sb, r.min)
	sb.WriteString(", max=")
	writeDecimalInt(sb, r.max)
	if !r.greedy {
		sb.WriteString("  not_greedy")
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
	alternatives := make([]RegExpTree, 0, len(a.alternatives))
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
					if ch, ok := regExpTree_singleCharUnit(charAlternative); ok {
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
	return &RegExpTree_Alternation{alternatives: alternatives}
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
func (a *RegExpTree_Alternation) AppendSourceCode(sb *strings.Builder) {
	for i, n := 0, len(a.alternatives); i < n; i++ {
		if i != 0 {
			sb.WriteByte('|')
		}
		a.alternatives[i].AppendSourceCode(sb)
	}
}

// AppendDebugInfo emits nothing besides children.
//
// Java: Alternation.appendDebugInfo(StringBuilder)
func (a *RegExpTree_Alternation) AppendDebugInfo(sb *strings.Builder) {
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

func (l *RegExpTree_LookaheadAssertion) IsCaseSensitive() bool   { return l.body.IsCaseSensitive() }
func (l *RegExpTree_LookaheadAssertion) ContainsAnchor() bool    { return l.body.ContainsAnchor() }
func (l *RegExpTree_LookaheadAssertion) HasCapturingGroup() bool { return l.body.HasCapturingGroup() }
func (l *RegExpTree_LookaheadAssertion) NumCapturingGroups() int { return l.body.NumCapturingGroups() }
func (l *RegExpTree_LookaheadAssertion) Children() []RegExpTree  { return []RegExpTree{l.body} }

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

// AppendSourceCode emits (?=...) or (?!...).
//
// Java: LookaheadAssertion.appendSourceCode(StringBuilder)
func (l *RegExpTree_LookaheadAssertion) AppendSourceCode(sb *strings.Builder) {
	if l.positive {
		sb.WriteString("(?=")
	} else {
		sb.WriteString("(?!")
	}
	l.body.AppendSourceCode(sb)
	sb.WriteByte(')')
}

// AppendDebugInfo emits positive or negative.
//
// Java: LookaheadAssertion.appendDebugInfo(StringBuilder)
func (l *RegExpTree_LookaheadAssertion) AppendDebugInfo(sb *strings.Builder) {
	if l.positive {
		sb.WriteString("positive")
	} else {
		sb.WriteString("negative")
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

func (l *RegExpTree_LookbehindAssertion) IsCaseSensitive() bool   { return l.body.IsCaseSensitive() }
func (l *RegExpTree_LookbehindAssertion) ContainsAnchor() bool    { return l.body.ContainsAnchor() }
func (l *RegExpTree_LookbehindAssertion) HasCapturingGroup() bool { return l.body.HasCapturingGroup() }
func (l *RegExpTree_LookbehindAssertion) NumCapturingGroups() int { return l.body.NumCapturingGroups() }
func (l *RegExpTree_LookbehindAssertion) Children() []RegExpTree  { return []RegExpTree{l.body} }

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

// AppendSourceCode emits (?<=...) or (?<!...).
//
// Java: LookbehindAssertion.appendSourceCode(StringBuilder)
func (l *RegExpTree_LookbehindAssertion) AppendSourceCode(sb *strings.Builder) {
	if l.positive {
		sb.WriteString("(?<=")
	} else {
		sb.WriteString("(?<!")
	}
	l.body.AppendSourceCode(sb)
	sb.WriteByte(')')
}

// AppendDebugInfo emits positive or negative.
//
// Java: LookbehindAssertion.appendDebugInfo(StringBuilder)
func (l *RegExpTree_LookbehindAssertion) AppendDebugInfo(sb *strings.Builder) {
	if l.positive {
		sb.WriteString("positive")
	} else {
		sb.WriteString("negative")
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

func (c *RegExpTree_CapturingGroup) IsCaseSensitive() bool  { return c.body.IsCaseSensitive() }
func (c *RegExpTree_CapturingGroup) ContainsAnchor() bool   { return c.body.ContainsAnchor() }
func (c *RegExpTree_CapturingGroup) Children() []RegExpTree { return []RegExpTree{c.body} }

// Simplify simplifies the body.
//
// Java: CapturingGroup.simplify(String)
func (c *RegExpTree_CapturingGroup) Simplify(flags string) RegExpTree {
	return NewRegExpTree_CapturingGroup(c.body.Simplify(flags))
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

// AppendSourceCode emits (...).
//
// Java: CapturingGroup.appendSourceCode(StringBuilder)
func (c *RegExpTree_CapturingGroup) AppendSourceCode(sb *strings.Builder) {
	sb.WriteByte('(')
	c.body.AppendSourceCode(sb)
	sb.WriteByte(')')
}

// AppendDebugInfo emits nothing besides children.
//
// Java: CapturingGroup.appendDebugInfo(StringBuilder)
func (c *RegExpTree_CapturingGroup) AppendDebugInfo(sb *strings.Builder) {
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
	s := RegExpTree_Concatenation_Simplifier{
		flags:      flags,
		simplified: make([]RegExpTree, 0, len(c.elements)),
	}
	for _, element := range c.elements {
		s.Simplify(element.Simplify(flags))
	}

	switch len(s.simplified) {
	case 0:
		return RegExpTree_Empty_INSTANCE
	case 1:
		return s.simplified[0]
	default:
		return &RegExpTree_Concatenation{elements: s.simplified}
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
	return &RegExpTree_Concatenation_Simplifier{
		flags:      flags,
		simplified: []RegExpTree{},
	}
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
func (c *RegExpTree_Concatenation) AppendSourceCode(sb *strings.Builder) {
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
			sb.WriteString("(?:")
			element.AppendSourceCode(sb)
			sb.WriteByte(')')
		} else {
			start := sb.Len()
			element.AppendSourceCode(sb)
			if digitsMightBleed && sb.Len() > start {
				s := sb.String()
				firstChar := s[start]
				if '0' <= firstChar && firstChar <= '9' {
					// Bleeding happened.
					// If the last character would be ambiguous
					// with a repetition, escape it.
					saved := strings.Clone(s)
					sb.Reset()
					if saved[start-1] == '{' {
						// Concatenation from optimization of
						// /{(?:0,}/ -> /\{0,}/
						sb.WriteString(saved[:start-1])
						sb.WriteString("\\{")
						sb.WriteString(saved[start:])
					} else {
						// Or parenthesize otherwise.
						// Concatenation from optimization of
						// /(.)\1(?:0)/ -> /(.)\1(?:0)/.
						sb.WriteString(saved[:start])
						sb.WriteString("(?:")
						sb.WriteString(saved[start:])
						sb.WriteByte(')')
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
func (c *RegExpTree_Concatenation) AppendDebugInfo(sb *strings.Builder) {
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

// Clean-room specification based on Google Guava base, collect, primitives, escape, and net utilities (Apache 2.0).

package guava

import (
	"fmt"
	"iter"
	"strings"
	"unicode/utf16"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javanio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaregex"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
)

// ---------------------------------------------------------------------------
// Immutable Collection Builders & Iterables / Chars
// ---------------------------------------------------------------------------

type ImmutableList_Builder[T any] struct {
	elems []T
}

func ImmutableList_NewBuilder[T any]() *ImmutableList_Builder[T] {
	return &ImmutableList_Builder[T]{}
}

func (b *ImmutableList_Builder[T]) Add(v T) *ImmutableList_Builder[T] {
	CheckNotNull(v)
	b.elems = append(b.elems, v)
	return b
}

func (b *ImmutableList_Builder[T]) AddAll(vs ...T) *ImmutableList_Builder[T] {
	for _, v := range vs {
		CheckNotNull(v)
		b.elems = append(b.elems, v)
	}
	return b
}

func (b *ImmutableList_Builder[T]) AddAllSeq(seq iter.Seq[T]) *ImmutableList_Builder[T] {
	for v := range seq {
		CheckNotNull(v)
		b.elems = append(b.elems, v)
	}
	return b
}

func (b *ImmutableList_Builder[T]) Build() *ImmutableList[T] {
	return ImmutableList_CopyOf(b.elems)
}

type ImmutableSet_Builder[T comparable] struct {
	set *javautil.LinkedHashSet[T]
}

func ImmutableSet_NewBuilder[T comparable]() *ImmutableSet_Builder[T] {
	return &ImmutableSet_Builder[T]{set: javautil.NewLinkedHashSet[T]()}
}

func (b *ImmutableSet_Builder[T]) Add(v T) *ImmutableSet_Builder[T] {
	CheckNotNull(v)
	b.set.Add(v)
	return b
}

func (b *ImmutableSet_Builder[T]) AddAll(vs ...T) *ImmutableSet_Builder[T] {
	for _, v := range vs {
		CheckNotNull(v)
		b.set.Add(v)
	}
	return b
}

func (b *ImmutableSet_Builder[T]) AddAllSeq(seq iter.Seq[T]) *ImmutableSet_Builder[T] {
	for v := range seq {
		CheckNotNull(v)
		b.set.Add(v)
	}
	return b
}

func (b *ImmutableSet_Builder[T]) Build() *ImmutableSet[T] {
	return ImmutableSet_CopyOf(b.set.ToSlice())
}

type ImmutableMap_Builder[K comparable, V any] struct {
	entries []struct {
		k K
		v V
	}
}

func ImmutableMap_NewBuilder[K comparable, V any]() *ImmutableMap_Builder[K, V] {
	return &ImmutableMap_Builder[K, V]{}
}

func (b *ImmutableMap_Builder[K, V]) Put(k K, v V) *ImmutableMap_Builder[K, V] {
	CheckNotNull(k)
	CheckNotNull(v)
	b.entries = append(b.entries, struct {
		k K
		v V
	}{k: k, v: v})
	return b
}

func (b *ImmutableMap_Builder[K, V]) PutAll(m *javautil.LinkedHashMap[K, V]) *ImmutableMap_Builder[K, V] {
	if m == nil {
		panic(javalang.NewNullPointerException())
	}
	for k, v := range m.All() {
		b.Put(k, v)
	}
	return b
}

func (b *ImmutableMap_Builder[K, V]) BuildOrThrow() *ImmutableMap[K, V] {
	m := javautil.NewLinkedHashMapWithCapacity[K, V](len(b.entries))
	for _, e := range b.entries {
		if m.ContainsKey(e.k) {
			panic(javalang.NewIllegalArgumentException(fmt.Sprintf("Multiple entries with same key: %v", e.k)))
		}
		m.Put(e.k, e.v)
	}
	return &ImmutableMap[K, V]{m: m}
}

func (b *ImmutableMap_Builder[K, V]) BuildKeepingLast() *ImmutableMap[K, V] {
	m := javautil.NewLinkedHashMapWithCapacity[K, V](len(b.entries))
	for _, e := range b.entries {
		m.Put(e.k, e.v)
	}
	return &ImmutableMap[K, V]{m: m}
}

func (b *ImmutableMap_Builder[K, V]) Build() *ImmutableMap[K, V] {
	return b.BuildOrThrow()
}

// Iterables_GetOnlyElement returns the single element contained in elems, or panics with
// NoSuchElementException (if empty) or IllegalArgumentException (if size > 1).
func Iterables_GetOnlyElement[T any](elems []T) T {
	if len(elems) == 0 {
		panic(javalang.NewNoSuchElementException(""))
	}
	if len(elems) > 1 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("expected one element but was: <%v, ...>", elems[0])))
	}
	return elems[0]
}

// Iterables_GetOnlyElementSeq returns the single element yielded by seq.
func Iterables_GetOnlyElementSeq[T any](seq iter.Seq[T]) T {
	var first T
	count := 0
	for v := range seq {
		if count == 0 {
			first = v
		} else {
			panic(javalang.NewIllegalArgumentException("expected one element"))
		}
		count++
	}
	if count == 0 {
		panic(javalang.NewNoSuchElementException(""))
	}
	return first
}

// Iterables_Concat concatenates slices into a single slice.
func Iterables_Concat[T any](slices ...[]T) []T {
	total := 0
	for _, s := range slices {
		total += len(s)
	}
	out := make([]T, 0, total)
	for _, s := range slices {
		out = append(out, s...)
	}
	return out
}

// Chars_AsList wraps a []uint16 (Java char[]) into a slice copy / ArrayList.
func Chars_AsList(backing []uint16) *javautil.ArrayList[uint16] {
	CheckNotNull(backing)
	return javautil.NewArrayListFrom(backing...)
}

// Chars_ToArray converts a collection slice of uint16 into a fresh []uint16.
func Chars_ToArray(elems []uint16) []uint16 {
	out := make([]uint16, len(elems))
	copy(out, elems)
	return out
}

// Chars_Contains reports whether target is present in array.
func Chars_Contains(array []uint16, target uint16) bool {
	for _, c := range array {
		if c == target {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// com.google.common.base.Strings & Ascii
// ---------------------------------------------------------------------------

// Strings_IsNullOrEmpty reports whether s is empty (or null represented via *string helper).
func Strings_IsNullOrEmpty(s string) bool {
	return s == ""
}

func Strings_IsNullOrEmptyPtr(s *string) bool {
	return s == nil || *s == ""
}

func Strings_NullToEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Strings_EmptyToNull(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func Strings_LenientFormat(template string, args ...any) string {
	return LenientFormat(template, args...)
}

func Strings_Repeat(s string, count int) string {
	if count < 0 {
		panic(javalang.NewIllegalArgumentException(fmt.Sprintf("invalid count: %d", count)))
	}
	if count <= 1 {
		if count == 0 {
			return ""
		}
		return s
	}
	return strings.Repeat(s, count)
}

func Strings_PadStart(s string, minLength int, padChar uint16) string {
	n := javalang.StringLength(s)
	if n >= minLength {
		return s
	}
	padStr := javalang.StringFromUTF16([]uint16{padChar})
	return strings.Repeat(padStr, minLength-n) + s
}

func Strings_PadEnd(s string, minLength int, padChar uint16) string {
	n := javalang.StringLength(s)
	if n >= minLength {
		return s
	}
	padStr := javalang.StringFromUTF16([]uint16{padChar})
	return s + strings.Repeat(padStr, minLength-n)
}

func Strings_CommonPrefix(a, b string) string {
	na := javalang.StringLength(a)
	nb := javalang.StringLength(b)
	limit := min(na, nb)
	p := 0
	for p < limit && javalang.StringCharAt(a, p) == javalang.StringCharAt(b, p) {
		p++
	}
	// Do not split a valid surrogate pair if only the high surrogate matched.
	if p > 0 {
		prev := javalang.StringCharAt(a, p-1)
		if prev >= 0xD800 && prev <= 0xDBFF {
			if (p < na && javalang.StringCharAt(a, p) >= 0xDC00 && javalang.StringCharAt(a, p) <= 0xDFFF) ||
				(p < nb && javalang.StringCharAt(b, p) >= 0xDC00 && javalang.StringCharAt(b, p) <= 0xDFFF) {
				p--
			}
		}
	}
	return javalang.StringSubstring(a, 0, p)
}

func Strings_CommonSuffix(a, b string) string {
	na := javalang.StringLength(a)
	nb := javalang.StringLength(b)
	limit := min(na, nb)
	s := 0
	for s < limit && javalang.StringCharAt(a, na-s-1) == javalang.StringCharAt(b, nb-s-1) {
		s++
	}
	if s > 0 {
		first := javalang.StringCharAt(a, na-s)
		if first >= 0xDC00 && first <= 0xDFFF {
			if (na-s-1 >= 0 && javalang.StringCharAt(a, na-s-1) >= 0xD800 && javalang.StringCharAt(a, na-s-1) <= 0xDBFF) ||
				(nb-s-1 >= 0 && javalang.StringCharAt(b, nb-s-1) >= 0xD800 && javalang.StringCharAt(b, nb-s-1) <= 0xDBFF) {
				s--
			}
		}
	}
	return javalang.StringSubstring(a, na-s, na)
}

func Ascii_ToLowerCase(s string) string {
	return javalang.StringToLowerCaseASCII(s)
}

func Ascii_ToUpperCase(s string) string {
	return javalang.StringToUpperCaseASCII(s)
}

func Ascii_ToLowerCaseChar(c uint16) uint16 {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

func Ascii_ToUpperCaseChar(c uint16) uint16 {
	if c >= 'a' && c <= 'z' {
		return c - ('a' - 'A')
	}
	return c
}

func Ascii_IsLowerCase(c uint16) bool { return c >= 'a' && c <= 'z' }
func Ascii_IsUpperCase(c uint16) bool { return c >= 'A' && c <= 'Z' }

func Ascii_EqualsIgnoreCase(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		c1 := s1[i]
		c2 := s2[i]
		if c1 == c2 {
			continue
		}
		if c1 >= 'A' && c1 <= 'Z' {
			c1 += 'a' - 'A'
		}
		if c2 >= 'A' && c2 <= 'Z' {
			c2 += 'a' - 'A'
		}
		if c1 != c2 || c1 < 'a' || c1 > 'z' {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// com.google.common.base.Joiner
// ---------------------------------------------------------------------------

type Joiner struct {
	sep        string
	skipNulls  bool
	hasNullStr bool
	nullText   string
}

func Joiner_On(separator string) *Joiner {
	return &Joiner{sep: separator}
}

func Joiner_OnChar(separator uint16) *Joiner {
	return &Joiner{sep: javalang.StringFromUTF16([]uint16{separator})}
}

func (j *Joiner) SkipNulls() *Joiner {
	if j.hasNullStr {
		panic(javalang.NewUnsupportedOperationException("already specified useForNull"))
	}
	cp := *j
	cp.skipNulls = true
	return &cp
}

func (j *Joiner) UseForNull(nullText string) *Joiner {
	if j.skipNulls || j.hasNullStr {
		panic(javalang.NewUnsupportedOperationException("already specified null behavior"))
	}
	cp := *j
	cp.hasNullStr = true
	cp.nullText = nullText
	return &cp
}

func (j *Joiner) JoinStrings(parts []string) string {
	return strings.Join(parts, j.sep)
}

func (j *Joiner) Join(parts ...any) string {
	return j.JoinAny(parts)
}

func (j *Joiner) JoinAny(parts []any) string {
	var sb strings.Builder
	first := true
	for _, p := range parts {
		isNil := javalang.IsNil(p)
		if isNil {
			if j.skipNulls {
				continue
			}
			if !j.hasNullStr {
				panic(javalang.NewNullPointerException())
			}
		}
		if !first {
			sb.WriteString(j.sep)
		}
		first = false
		if isNil {
			sb.WriteString(j.nullText)
		} else if ch, ok := p.(uint16); ok {
			sb.WriteString(javalang.StringFromUTF16([]uint16{ch}))
		} else {
			sb.WriteString(javalang.ToJavaString(p))
		}
	}
	return sb.String()
}

// ---------------------------------------------------------------------------
// com.google.common.base.CharMatcher
// ---------------------------------------------------------------------------

type CharMatcher struct {
	pred func(uint16) bool
}

func CharMatcher_Any() *CharMatcher {
	return &CharMatcher{pred: func(uint16) bool { return true }}
}

func CharMatcher_None() *CharMatcher {
	return &CharMatcher{pred: func(uint16) bool { return false }}
}

func CharMatcher_Ascii() *CharMatcher {
	return &CharMatcher{pred: func(c uint16) bool { return c <= 0x7F }}
}

// CharMatcher_Whitespace implements Guava's Unicode table of whitespace characters.
func CharMatcher_Whitespace() *CharMatcher {
	return &CharMatcher{pred: isGuavaWhitespace}
}

func CharMatcher_JavaLetterOrDigit() *CharMatcher {
	return &CharMatcher{pred: javalang.CharacterIsLetterOrDigit}
}

func isGuavaWhitespace(c uint16) bool {
	switch c {
	case '\t', '\n', 0x000B, '\f', '\r', ' ', 0x0085, 0x1680,
		0x2028, 0x2029, 0x205F, 0x3000:
		return true
	}
	return (c >= 0x2000 && c <= 0x200A && c != 0x2007)
}

func CharMatcher_Is(match uint16) *CharMatcher {
	return &CharMatcher{pred: func(c uint16) bool { return c == match }}
}

func CharMatcher_IsNot(match uint16) *CharMatcher {
	return &CharMatcher{pred: func(c uint16) bool { return c != match }}
}

func CharMatcher_InRange(startInclusive, endInclusive uint16) *CharMatcher {
	CheckArgument(endInclusive >= startInclusive)
	return &CharMatcher{pred: func(c uint16) bool { return c >= startInclusive && c <= endInclusive }}
}

func CharMatcher_AnyOf(sequence string) *CharMatcher {
	units := javalang.StringToUTF16(sequence)
	switch len(units) {
	case 0:
		return CharMatcher_None()
	case 1:
		return CharMatcher_Is(units[0])
	default:
		set := make(map[uint16]struct{}, len(units))
		for _, u := range units {
			set[u] = struct{}{}
		}
		return &CharMatcher{pred: func(c uint16) bool {
			_, ok := set[c]
			return ok
		}}
	}
}

func (m *CharMatcher) Or(other *CharMatcher) *CharMatcher {
	CheckNotNull(other)
	return &CharMatcher{pred: func(c uint16) bool { return m.pred(c) || other.pred(c) }}
}

func (m *CharMatcher) And(other *CharMatcher) *CharMatcher {
	CheckNotNull(other)
	return &CharMatcher{pred: func(c uint16) bool { return m.pred(c) && other.pred(c) }}
}

func (m *CharMatcher) Negate() *CharMatcher {
	return &CharMatcher{pred: func(c uint16) bool { return !m.pred(c) }}
}

func (m *CharMatcher) Matches(c uint16) bool { return m.pred(c) }

func (m *CharMatcher) MatchesAllOf(s string) bool {
	units := javalang.StringToUTF16(s)
	for i := len(units) - 1; i >= 0; i-- {
		if !m.pred(units[i]) {
			return false
		}
	}
	return true
}

func (m *CharMatcher) MatchesNoneOf(s string) bool {
	return m.IndexIn(s) == -1
}

func (m *CharMatcher) MatchesAnyOf(s string) bool {
	return !m.MatchesNoneOf(s)
}

func (m *CharMatcher) IndexIn(s string) int {
	return m.IndexInFrom(s, 0)
}

func (m *CharMatcher) IndexInFrom(s string, start int) int {
	units := javalang.StringToUTF16(s)
	CheckPositionIndex(start, len(units))
	for i := start; i < len(units); i++ {
		if m.pred(units[i]) {
			return i
		}
	}
	return -1
}

func (m *CharMatcher) LastIndexIn(s string) int {
	units := javalang.StringToUTF16(s)
	for i := len(units) - 1; i >= 0; i-- {
		if m.pred(units[i]) {
			return i
		}
	}
	return -1
}

func (m *CharMatcher) CountIn(s string) int {
	units := javalang.StringToUTF16(s)
	n := 0
	for _, u := range units {
		if m.pred(u) {
			n++
		}
	}
	return n
}

func (m *CharMatcher) TrimFrom(s string) string {
	units := javalang.StringToUTF16(s)
	first := 0
	for first < len(units) && m.pred(units[first]) {
		first++
	}
	last := len(units) - 1
	for last > first && m.pred(units[last]) {
		last--
	}
	if first == 0 && last == len(units)-1 {
		return s
	}
	if first > last {
		return ""
	}
	return javalang.StringFromUTF16(units[first : last+1])
}

func (m *CharMatcher) TrimLeadingFrom(s string) string {
	units := javalang.StringToUTF16(s)
	first := 0
	for first < len(units) && m.pred(units[first]) {
		first++
	}
	if first == 0 {
		return s
	}
	return javalang.StringFromUTF16(units[first:])
}

func (m *CharMatcher) TrimTrailingFrom(s string) string {
	units := javalang.StringToUTF16(s)
	last := len(units) - 1
	for last >= 0 && m.pred(units[last]) {
		last--
	}
	if last == len(units)-1 {
		return s
	}
	return javalang.StringFromUTF16(units[:last+1])
}

func (m *CharMatcher) RemoveFrom(s string) string {
	units := javalang.StringToUTF16(s)
	var out []uint16
	for i, u := range units {
		if m.pred(u) {
			if out == nil {
				out = make([]uint16, 0, len(units)-1)
				out = append(out, units[:i]...)
			}
			continue
		}
		if out != nil {
			out = append(out, u)
		}
	}
	if out == nil {
		return s
	}
	return javalang.StringFromUTF16(out)
}

func (m *CharMatcher) RetainFrom(s string) string {
	return m.Negate().RemoveFrom(s)
}

func (m *CharMatcher) ReplaceFrom(s string, replacement uint16) string {
	units := javalang.StringToUTF16(s)
	changed := false
	for i, u := range units {
		if m.pred(u) {
			units[i] = replacement
			changed = true
		}
	}
	if !changed {
		return s
	}
	return javalang.StringFromUTF16(units)
}

func (m *CharMatcher) CollapseFrom(s string, replacement uint16) string {
	units := javalang.StringToUTF16(s)
	var out []uint16
	inRun := false
	for _, u := range units {
		if m.pred(u) {
			if !inRun {
				out = append(out, replacement)
				inRun = true
			}
		} else {
			out = append(out, u)
			inRun = false
		}
	}
	return javalang.StringFromUTF16(out)
}

// ---------------------------------------------------------------------------
// com.google.common.base.Splitter
// ---------------------------------------------------------------------------

type splitterStrategy func(s string) iter.Seq[struct{ start, end int }]

type Splitter struct {
	strategy    splitterStrategy
	omitEmpty   bool
	trimmer     *CharMatcher
	limit       int // 0 means unlimited (Integer.MAX_VALUE)
}

func Splitter_OnChar(separator uint16) *Splitter {
	return Splitter_OnMatcher(CharMatcher_Is(separator))
}

func Splitter_OnMatcher(matcher *CharMatcher) *Splitter {
	CheckNotNull(matcher)
	return &Splitter{
		trimmer: CharMatcher_None(),
		strategy: func(s string) iter.Seq[struct{ start, end int }] {
			return func(yield func(struct{ start, end int }) bool) {
				units := javalang.StringToUTF16(s)
				offset := 0
				for offset <= len(units) {
					sepIdx := -1
					for i := offset; i < len(units); i++ {
						if matcher.Matches(units[i]) {
							sepIdx = i
							break
						}
					}
					if sepIdx == -1 {
						yield(struct{ start, end int }{offset, len(units)})
						return
					}
					if !yield(struct{ start, end int }{offset, sepIdx}) {
						return
					}
					offset = sepIdx + 1
				}
			}
		},
	}
}

func Splitter_On(separator string) *Splitter {
	sepLen := javalang.StringLength(separator)
	CheckArgumentMsg(sepLen != 0, "The separator may not be the empty string.")
	if sepLen == 1 {
		return Splitter_OnChar(javalang.StringCharAt(separator, 0))
	}
	return &Splitter{
		trimmer: CharMatcher_None(),
		strategy: func(s string) iter.Seq[struct{ start, end int }] {
			return func(yield func(struct{ start, end int }) bool) {
				n := javalang.StringLength(s)
				offset := 0
				for offset <= n {
					idx := javalang.StringIndexOfFrom(s, separator, offset)
					if idx == -1 {
						yield(struct{ start, end int }{offset, n})
						return
					}
					if !yield(struct{ start, end int }{offset, idx}) {
						return
					}
					offset = idx + sepLen
				}
			}
		},
	}
}

func Splitter_OnPattern(separatorPattern string) *Splitter {
	pat := javaregex.Pattern_Compile(separatorPattern)
	CheckArgumentMsg(!pat.Matcher("").Matches(), "The pattern may not match the empty string: %s", separatorPattern)
	return &Splitter{
		trimmer: CharMatcher_None(),
		strategy: func(s string) iter.Seq[struct{ start, end int }] {
			return func(yield func(struct{ start, end int }) bool) {
				n := javalang.StringLength(s)
				m := pat.Matcher(s)
				offset := 0
				for offset <= n {
					if !m.FindFrom(offset) {
						yield(struct{ start, end int }{offset, n})
						return
					}
					ms := m.Start()
					me := m.End()
					if !yield(struct{ start, end int }{offset, ms}) {
						return
					}
					if ms == me {
						offset = me + 1
					} else {
						offset = me
					}
				}
			}
		},
	}
}

func (sp *Splitter) OmitEmptyStrings() *Splitter {
	cp := *sp
	cp.omitEmpty = true
	return &cp
}

func (sp *Splitter) TrimResults() *Splitter {
	return sp.TrimResultsWith(CharMatcher_Whitespace())
}

func (sp *Splitter) TrimResultsWith(trimmer *CharMatcher) *Splitter {
	CheckNotNull(trimmer)
	cp := *sp
	cp.trimmer = trimmer
	return &cp
}

func (sp *Splitter) Limit(maxItems int) *Splitter {
	CheckArgumentMsg1(maxItems > 0, "must be greater than zero: %s", maxItems)
	cp := *sp
	cp.limit = maxItems
	return &cp
}

func (sp *Splitter) Split(sequence string) iter.Seq[string] {
	return func(yield func(string) bool) {
		units := javalang.StringToUTF16(sequence)
		remaining := sp.limit
		next, stop := iter.Pull(sp.strategy(sequence))
		defer stop()

		for {
			seg, ok := next()
			if !ok {
				return
			}
			start, end := seg.start, seg.end
			if sp.limit > 0 && remaining == 1 {
				end = len(units)
				stop()
			}
			for start < end && sp.trimmer.Matches(units[start]) {
				start++
			}
			for end > start && sp.trimmer.Matches(units[end-1]) {
				end--
			}
			if sp.omitEmpty && start == end {
				continue
			}
			if sp.limit > 0 {
				remaining--
				if remaining == 0 {
					// Re-trim trailing after extending to len(units) above.
					if !yield(javalang.StringFromUTF16(units[start:end])) {
						return
					}
					return
				}
			}
			if !yield(javalang.StringFromUTF16(units[start:end])) {
				return
			}
		}
	}
}

func (sp *Splitter) SplitToList(sequence string) []string {
	var out []string
	for part := range sp.Split(sequence) {
		out = append(out, part)
	}
	if out == nil {
		return []string{}
	}
	return out
}

// ---------------------------------------------------------------------------
// com.google.common.base.MoreObjects
// ---------------------------------------------------------------------------

func MoreObjects_FirstNonNull[T any](first, second T) T {
	if !javalang.IsNil(first) {
		return first
	}
	if !javalang.IsNil(second) {
		return second
	}
	panic(javalang.NewNullPointerExceptionMsg("Both parameters are null"))
}

type ToStringHelper struct {
	className string
	omitNulls bool
	fields    []struct {
		hasName bool
		name    string
		value   any
	}
}

func MoreObjects_ToStringHelper(className string) *ToStringHelper {
	return &ToStringHelper{className: className}
}

func (h *ToStringHelper) OmitNullValues() *ToStringHelper {
	h.omitNulls = true
	return h
}

func (h *ToStringHelper) Add(name string, value any) *ToStringHelper {
	CheckNotNull(name)
	h.fields = append(h.fields, struct {
		hasName bool
		name    string
		value   any
	}{hasName: true, name: name, value: value})
	return h
}

func (h *ToStringHelper) AddValue(value any) *ToStringHelper {
	h.fields = append(h.fields, struct {
		hasName bool
		name    string
		value   any
	}{hasName: false, value: value})
	return h
}

func (h *ToStringHelper) String() string {
	var sb strings.Builder
	sb.WriteString(h.className)
	sb.WriteByte('{')
	sep := ""
	for _, f := range h.fields {
		if h.omitNulls && javalang.IsNil(f.value) {
			continue
		}
		sb.WriteString(sep)
		sep = ", "
		if f.hasName {
			sb.WriteString(f.name)
			sb.WriteByte('=')
		}
		sb.WriteString(javalang.ToJavaString(f.value))
	}
	sb.WriteByte('}')
	return sb.String()
}

// ---------------------------------------------------------------------------
// com.google.common.escape.Escaper & com.google.common.net.PercentEscaper
// ---------------------------------------------------------------------------

type Escaper interface {
	Escape(s string) string
}

// PercentEscaper implements Guava's RFC 3986 / custom safeChars UTF-8 percent-escaper.
type PercentEscaper struct {
	safeAscii    [128]bool
	plusForSpace bool
}

func NewPercentEscaper(safeChars string, plusForSpace bool) *PercentEscaper {
	units := javalang.StringToUTF16(safeChars)
	for _, u := range units {
		if (u >= '0' && u <= '9') || (u >= 'a' && u <= 'z') || (u >= 'A' && u <= 'Z') {
			panic(javalang.NewIllegalArgumentException("Alphanumeric characters are always 'safe' and should not be explicitly specified"))
		}
	}
	for _, u := range units {
		if u == ' ' && plusForSpace {
			panic(javalang.NewIllegalArgumentException("plusForSpace cannot be specified when space is a 'safe' character"))
		}
		if u > 0x7F {
			panic(javalang.NewIllegalArgumentException("PercentEscaper only supports ASCII 'safe' characters"))
		}
	}
	pe := &PercentEscaper{plusForSpace: plusForSpace}
	for c := byte('0'); c <= '9'; c++ {
		pe.safeAscii[c] = true
	}
	for c := byte('a'); c <= 'z'; c++ {
		pe.safeAscii[c] = true
	}
	for c := byte('A'); c <= 'Z'; c++ {
		pe.safeAscii[c] = true
	}
	for _, u := range units {
		pe.safeAscii[u] = true
	}
	return pe
}

const upperHex = "0123456789ABCDEF"

func (pe *PercentEscaper) Escape(s string) string {
	units := javalang.StringToUTF16(s)
	// Fast path: all safe ASCII.
	allSafe := true
	for _, u := range units {
		if u >= 128 || !pe.safeAscii[u] {
			allSafe = false
			break
		}
	}
	if allSafe {
		return s
	}

	var sb strings.Builder
	sb.Grow(len(s) + 16)
	for i := 0; i < len(units); i++ {
		u := units[i]
		if u < 128 && pe.safeAscii[u] {
			sb.WriteByte(byte(u))
			continue
		}
		if u == ' ' && pe.plusForSpace {
			sb.WriteByte('+')
			continue
		}
		var cp rune
		if u >= 0xD800 && u <= 0xDBFF {
			if i+1 >= len(units) || units[i+1] < 0xDC00 || units[i+1] > 0xDFFF {
				panic(javalang.NewIllegalArgumentException("Unpaired high surrogate"))
			}
			cp = utf16.DecodeRune(rune(u), rune(units[i+1]))
			i++
		} else if u >= 0xDC00 && u <= 0xDFFF {
			panic(javalang.NewIllegalArgumentException("Unexpected low surrogate"))
		} else {
			cp = rune(u)
		}
		utf8Bytes := javanio.StandardCharsets_UTF_8.Encode(string(cp))
		for _, b := range utf8Bytes {
			sb.WriteByte('%')
			sb.WriteByte(upperHex[b>>4])
			sb.WriteByte(upperHex[b&0x0F])
		}
	}
	return sb.String()
}

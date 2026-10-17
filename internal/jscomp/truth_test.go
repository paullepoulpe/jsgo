package jscomp

// Go-only: test-helper — Google Truth (com.google.common.truth) and the
// org.junit.Assert helpers used by the ported jscomp tests.
//
// Java's static imports map to package functions:
//
//	assertThat(x).isEqualTo(y)                    →  AssertThat(x).IsEqualTo(y)
//	assertWithMessage("m %s", a).that(x).isTrue() →  AssertWithMessage("m %s", a).That(x).IsTrue()
//	assertThrows(Foo.class, () -> ...)            →  AssertThrows[*Foo](func() { ... })
//	assert_().fail() / fail("m")                  →  Fail("m")
//
// Truth picks a Subject by the static type of the argument (BooleanSubject,
// StringSubject, IterableSubject, MapSubject, IntegerSubject...). Go has no
// overloading, so a single Subject carries all of their methods and
// dispatches on the dynamic type of the actual value (Go-only:
// language-adapter). Equality is Objects.equals (Equals methods, element-wise
// for lists), as in rhinotesting.Equal.
//
// A failed assertion panics with *jutil.AssertionError, like Java; RunJUnit4
// reports it.

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// AssertThat is Truth.assertThat(actual).
func AssertThat(actual any) *Subject {
	return &Subject{actual: actual}
}

// AssertWithMessage is Truth.assertWithMessage(format, args...): the message
// prefixes the failure of the next assertion. The format uses Guava's
// lenient %s placeholders.
func AssertWithMessage(format string, args ...any) *StandardSubjectBuilder {
	return &StandardSubjectBuilder{message: jutil.LenientFormat(format, args...)}
}

// Fail is Assert.fail(message) / Truth's assert_().fail().
func Fail(message string) {
	panic(jutil.NewAssertionError(message))
}

// AssertThrows is Assert.assertThrows(E.class, body): body must throw (panic
// with) an exception assignable to E, which is returned.
func AssertThrows[E any](body func()) E {
	var caught E
	ok := false
	jutil.Try(body, func(e E) {
		caught = e
		ok = true
	})
	if !ok {
		Fail(fmt.Sprintf("expected %s to be thrown, but nothing was thrown", reflect.TypeOf((*E)(nil)).Elem()))
	}
	return caught
}

// StandardSubjectBuilder is Truth's StandardSubjectBuilder.
type StandardSubjectBuilder struct {
	message string
}

// That is StandardSubjectBuilder.that(actual).
func (b *StandardSubjectBuilder) That(actual any) *Subject {
	return &Subject{actual: actual, message: b.message}
}

// Fail is StandardSubjectBuilder.fail().
func (b *StandardSubjectBuilder) Fail() {
	panic(jutil.NewAssertionError(b.message))
}

// Subject is Truth's Subject together with the methods of its subclasses
// (BooleanSubject, ComparableSubject, StringSubject, IterableSubject,
// MapSubject).
type Subject struct {
	actual  any
	message string
}

// Actual is the value under test.
func (s *Subject) Actual() any { return s.actual }

// failf fails the assertion; facts are rendered like Truth's "key: value"
// lines, after the message given to AssertWithMessage.
func (s *Subject) failf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if s.message != "" {
		msg = s.message + "\n" + msg
	}
	panic(jutil.NewAssertionError(msg))
}

// show renders a value like Truth (String.valueOf, collections as [a, b]).
func show(v any) string { return rhinotesting.Show(v) }

// truthEqual is Objects.equals with List/Set/Map equality.
func truthEqual(a, b any) bool {
	a, b = truthNormalize(a), truthNormalize(b)
	if ma, ok := truthAsMap(a); ok {
		if mb, ok := truthAsMap(b); ok {
			return truthMapsEqual(ma, mb)
		}
	}
	if ia, ok := truthAsInt64(a); ok {
		if ib, ok := truthAsInt64(b); ok {
			return ia == ib
		}
	}
	if truthIsSet(a) && truthIsSet(b) {
		// Set.equals: the same elements in any order.
		return jutil.ObjectsEquals(a, b)
	}
	return rhinotesting.Equal(a, b)
}

// truthIsSet reports whether v is one of the jutil java.util.Set types
// (whose Equals is Set.equals).
func truthIsSet(v any) bool {
	if jutil.IsNil(v) {
		return false
	}
	if _, ok := v.(interface{ Equals(any) bool }); !ok {
		return false
	}
	name := reflect.TypeOf(v).String()
	return strings.HasPrefix(name, "*jutil.") && strings.Contains(name, "Set[")
}

// truthNormalize maps Go stand-ins to the Java values they represent: a
// *string to the String (or null), a pointer to another basic value (a
// nullable boxed primitive or enum, e.g. Symbol.getVisibility()'s
// *JSDocInfo_Visibility) to that value (or null), an iter.Seq to its
// elements.
func truthNormalize(v any) any {
	if p, ok := v.(*string); ok {
		if p == nil {
			return nil
		}
		return *p
	}
	if rv := reflect.ValueOf(v); rv.IsValid() && rv.Kind() == reflect.Pointer {
		switch rv.Type().Elem().Kind() {
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			if rv.IsNil() {
				return nil
			}
			return rv.Elem().Interface()
		}
	}
	if l, ok := truthSeqElements(v); ok {
		return l
	}
	return v
}

// truthSeqElements collects, in order, the elements of a non-nil
// iter.Seq[T] (func(yield func(T) bool)), the Go form of a Java Iterable
// (PORTING.md §3).
func truthSeqElements(v any) ([]any, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() != reflect.Func || rv.IsNil() {
		return nil, false
	}
	ft := rv.Type()
	if ft.NumIn() != 1 || ft.NumOut() != 0 || ft.IsVariadic() {
		return nil, false
	}
	yt := ft.In(0)
	if yt.Kind() != reflect.Func || yt.NumIn() != 1 || yt.NumOut() != 1 || yt.Out(0).Kind() != reflect.Bool {
		return nil, false
	}
	out := []any{}
	yield := reflect.MakeFunc(yt, func(args []reflect.Value) []reflect.Value {
		out = append(out, args[0].Interface())
		return []reflect.Value{reflect.ValueOf(true)}
	})
	rv.Call([]reflect.Value{yield})
	return out, true
}

// truthAsInt64 unboxes a Java integral value (Integer, Long, Short, Byte,
// Character), so that e.g. an int and an int32 compare like Java's
// Integer.equals on the same value.
func truthAsInt64(v any) (int64, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return 0, false
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv.Type().PkgPath() != "" {
			return 0, false // enums and other named types compare as themselves
		}
		return rv.Int(), true
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		if rv.Type().PkgPath() != "" {
			return 0, false
		}
		return int64(rv.Uint()), true
	}
	return 0, false
}

// truthEntry is a Map.Entry.
type truthEntry struct {
	key, value any
}

// truthAsMap returns the entries of a Java map stand-in (a jutil map with
// KeySlice and Get, or a Go map), in iteration order.
func truthAsMap(v any) ([]truthEntry, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil, false
	}
	if rv.Kind() == reflect.Map {
		var out []truthEntry
		it := rv.MapRange()
		for it.Next() {
			out = append(out, truthEntry{it.Key().Interface(), it.Value().Interface()})
		}
		return out, true
	}
	keys := rv.MethodByName("KeySlice")
	get := rv.MethodByName("Get")
	if !keys.IsValid() || !get.IsValid() || keys.Type().NumIn() != 0 || get.Type().NumIn() != 1 {
		return nil, false
	}
	if rv.Kind() == reflect.Pointer && rv.IsNil() {
		return nil, false
	}
	ks := keys.Call(nil)[0]
	out := make([]truthEntry, ks.Len())
	for i := range out {
		k := ks.Index(i)
		out[i] = truthEntry{k.Interface(), get.Call([]reflect.Value{k})[0].Interface()}
	}
	return out, true
}

func truthMapLookup(m []truthEntry, key any) (any, bool) {
	for _, e := range m {
		if truthEqual(e.key, key) {
			return e.value, true
		}
	}
	return nil, false
}

func truthMapsEqual(a, b []truthEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for _, e := range a {
		v, ok := truthMapLookup(b, e.key)
		if !ok || !truthEqual(e.value, v) {
			return false
		}
	}
	return true
}

func showMap(m []truthEntry) string {
	var sb strings.Builder
	sb.WriteByte('{')
	for i, e := range m {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(show(e.key))
		sb.WriteByte('=')
		sb.WriteString(show(e.value))
	}
	sb.WriteByte('}')
	return sb.String()
}

// showActual renders the actual value, maps as {k=v}.
func (s *Subject) showActual() string {
	if m, ok := truthAsMap(truthNormalize(s.actual)); ok {
		return showMap(m)
	}
	return show(s.actual)
}

// IsEqualTo is Subject.isEqualTo(expected).
func (s *Subject) IsEqualTo(expected any) {
	if !truthEqual(s.actual, expected) {
		s.failf("expected: %s\nbut was : %s", show(expected), s.showActual())
	}
}

// IsNotEqualTo is Subject.isNotEqualTo(unexpected).
func (s *Subject) IsNotEqualTo(unexpected any) {
	if truthEqual(s.actual, unexpected) {
		s.failf("expected not to be: %s", show(unexpected))
	}
}

// IsNull is Subject.isNull().
func (s *Subject) IsNull() {
	if !jutil.IsNil(truthNormalize(s.actual)) {
		s.failf("expected: null\nbut was : %s", s.showActual())
	}
}

// IsNotNull is Subject.isNotNull().
func (s *Subject) IsNotNull() {
	if jutil.IsNil(truthNormalize(s.actual)) {
		s.failf("expected not to be null")
	}
}

// IsSameInstanceAs is Subject.isSameInstanceAs(expected).
func (s *Subject) IsSameInstanceAs(expected any) {
	if !truthSameInstance(s.actual, expected) {
		s.failf("expected specific instance: %s\nbut was           : %s", show(expected), s.showActual())
	}
}

// IsNotSameInstanceAs is Subject.isNotSameInstanceAs(unexpected).
func (s *Subject) IsNotSameInstanceAs(unexpected any) {
	if truthSameInstance(s.actual, unexpected) {
		s.failf("expected not to be specific instance: %s", show(unexpected))
	}
}

func truthSameInstance(a, b any) bool {
	if jutil.IsNil(a) || jutil.IsNil(b) {
		return jutil.IsNil(a) && jutil.IsNil(b)
	}
	ra, rb := reflect.ValueOf(a), reflect.ValueOf(b)
	if ra.Type() != rb.Type() {
		return false
	}
	switch ra.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Func, reflect.Chan:
		return ra.Pointer() == rb.Pointer()
	case reflect.Slice:
		return ra.Pointer() == rb.Pointer() && ra.Len() == rb.Len()
	}
	if ra.Comparable() {
		return a == b
	}
	return false
}

// IsAnyOf is Subject.isAnyOf(first, second, rest...).
func (s *Subject) IsAnyOf(values ...any) {
	for _, v := range values {
		if truthEqual(s.actual, v) {
			return
		}
	}
	s.failf("expected any of: %s\nbut was        : %s", show(values), s.showActual())
}

// IsNoneOf is Subject.isNoneOf(first, second, rest...).
func (s *Subject) IsNoneOf(values ...any) {
	for _, v := range values {
		if truthEqual(s.actual, v) {
			s.failf("expected not to be any of: %s\nbut was                  : %s", show(values), s.showActual())
		}
	}
}

// IsIn is Subject.isIn(iterable).
func (s *Subject) IsIn(iterable any) {
	s.IsAnyOf(s.listOf(iterable)...)
}

// IsNotIn is Subject.isNotIn(iterable).
func (s *Subject) IsNotIn(iterable any) {
	s.IsNoneOf(s.listOf(iterable)...)
}

// IsInstanceOf is Subject.isInstanceOf(T.class).
func IsInstanceOf[T any](s *Subject) {
	if _, ok := s.actual.(T); !ok || jutil.IsNil(s.actual) {
		s.failf("expected instance of: %s\nbut was             : %s (%T)", reflect.TypeOf((*T)(nil)).Elem(), s.showActual(), s.actual)
	}
}

// IsNotInstanceOf is Subject.isNotInstanceOf(T.class).
func IsNotInstanceOf[T any](s *Subject) {
	if _, ok := s.actual.(T); ok && !jutil.IsNil(s.actual) {
		s.failf("expected not to be an instance of: %s\nbut was: %s (%T)", reflect.TypeOf((*T)(nil)).Elem(), s.showActual(), s.actual)
	}
}

// --- BooleanSubject

// IsTrue is BooleanSubject.isTrue().
func (s *Subject) IsTrue() {
	if b, ok := s.actual.(bool); !ok || !b {
		s.failf("expected to be true, but was %s", s.showActual())
	}
}

// IsFalse is BooleanSubject.isFalse().
func (s *Subject) IsFalse() {
	if b, ok := s.actual.(bool); !ok || b {
		s.failf("expected to be false, but was %s", s.showActual())
	}
}

// --- ComparableSubject (numbers)

func truthAsFloat(v any) (float64, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return 0, false
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint()), true
	case reflect.Float32, reflect.Float64:
		return rv.Float(), true
	}
	return 0, false
}

func (s *Subject) compare(other any) int {
	a, ok1 := truthAsFloat(s.actual)
	b, ok2 := truthAsFloat(other)
	if !ok1 || !ok2 {
		if c, ok := s.actual.(interface{ CompareTo(any) int }); ok {
			return c.CompareTo(other)
		}
		s.failf("cannot compare %s with %s", s.showActual(), show(other))
	}
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// IsGreaterThan is ComparableSubject.isGreaterThan(other).
func (s *Subject) IsGreaterThan(other any) {
	if s.compare(other) <= 0 {
		s.failf("expected to be greater than: %s\nbut was                    : %s", show(other), s.showActual())
	}
}

// IsLessThan is ComparableSubject.isLessThan(other).
func (s *Subject) IsLessThan(other any) {
	if s.compare(other) >= 0 {
		s.failf("expected to be less than: %s\nbut was                 : %s", show(other), s.showActual())
	}
}

// IsAtLeast is ComparableSubject.isAtLeast(other).
func (s *Subject) IsAtLeast(other any) {
	if s.compare(other) < 0 {
		s.failf("expected to be at least: %s\nbut was                : %s", show(other), s.showActual())
	}
}

// IsAtMost is ComparableSubject.isAtMost(other).
func (s *Subject) IsAtMost(other any) {
	if s.compare(other) > 0 {
		s.failf("expected to be at most: %s\nbut was               : %s", show(other), s.showActual())
	}
}

// --- StringSubject

func (s *Subject) str() (string, bool) {
	str, ok := truthNormalize(s.actual).(string)
	return str, ok
}

func (s *Subject) mustStr() string {
	str, ok := s.str()
	if !ok {
		s.failf("expected a string, but was %s", s.showActual())
	}
	return str
}

// StartsWith is StringSubject.startsWith(prefix).
func (s *Subject) StartsWith(prefix string) {
	if !strings.HasPrefix(s.mustStr(), prefix) {
		s.failf("expected to start with: %q\nbut was               : %q", prefix, s.mustStr())
	}
}

// EndsWith is StringSubject.endsWith(suffix).
func (s *Subject) EndsWith(suffix string) {
	if !strings.HasSuffix(s.mustStr(), suffix) {
		s.failf("expected to end with: %q\nbut was             : %q", suffix, s.mustStr())
	}
}

// Matches is StringSubject.matches(regex): the whole string matches.
func (s *Subject) Matches(regex string) {
	if !jutil.Pattern_Compile(regex).Matches(s.mustStr()) {
		s.failf("expected to match: %s\nbut was          : %q", regex, s.mustStr())
	}
}

// DoesNotMatch is StringSubject.doesNotMatch(regex).
func (s *Subject) DoesNotMatch(regex string) {
	if jutil.Pattern_Compile(regex).Matches(s.mustStr()) {
		s.failf("expected not to match: %s\nbut was              : %q", regex, s.mustStr())
	}
}

// ContainsMatch is StringSubject.containsMatch(regex).
func (s *Subject) ContainsMatch(regex string) {
	if !jutil.Pattern_Compile(regex).Find(s.mustStr()) {
		s.failf("expected to contain a match for: %s\nbut was                        : %q", regex, s.mustStr())
	}
}

// DoesNotContainMatch is StringSubject.doesNotContainMatch(regex).
func (s *Subject) DoesNotContainMatch(regex string) {
	if jutil.Pattern_Compile(regex).Find(s.mustStr()) {
		s.failf("expected not to contain a match for: %s\nbut was                            : %q", regex, s.mustStr())
	}
}

// HasLength is StringSubject.hasLength(n), in UTF-16 units.
func (s *Subject) HasLength(n int) {
	if l := jutil.StringLength(s.mustStr()); l != n {
		s.failf("expected to have length: %d\nbut was                : %q", n, s.mustStr())
	}
}

// --- IterableSubject / MapSubject / StringSubject shared methods

func (s *Subject) listOf(v any) []any {
	v = truthNormalize(v)
	if truthIsSlice(v) {
		// A nil slice is the Go empty list.
		l, _ := rhinotesting.AsList(v)
		return l
	}
	if jutil.IsNil(v) {
		return nil
	}
	if l, ok := rhinotesting.AsList(v); ok {
		return l
	}
	s.failf("expected a collection, but was %s", show(v))
	return nil
}

func (s *Subject) list() []any {
	if v := truthNormalize(s.actual); jutil.IsNil(v) && !truthIsSlice(v) {
		s.failf("expected a collection, but was null")
	}
	return s.listOf(s.actual)
}

// truthIsSlice reports whether v is a Go slice (nil included): a Java list
// whose nil value the ports use for the empty list.
func truthIsSlice(v any) bool {
	return v != nil && reflect.TypeOf(v).Kind() == reflect.Slice
}

// size is the size of a string (UTF-16 units), collection or map.
func (s *Subject) size() int {
	if str, ok := s.str(); ok {
		return jutil.StringLength(str)
	}
	if m, ok := truthAsMap(truthNormalize(s.actual)); ok {
		return len(m)
	}
	return len(s.list())
}

// HasSize is IterableSubject.hasSize(n) / MapSubject.hasSize(n).
func (s *Subject) HasSize(n int) {
	if s.size() != n {
		s.failf("expected to have size: %d\nbut was              : %s", n, s.showActual())
	}
}

// IsEmpty is IterableSubject/MapSubject/StringSubject.isEmpty().
func (s *Subject) IsEmpty() {
	if s.size() != 0 {
		s.failf("expected to be empty, but was %s", s.showActual())
	}
}

// IsNotEmpty is IterableSubject/MapSubject/StringSubject.isNotEmpty().
func (s *Subject) IsNotEmpty() {
	if s.size() == 0 {
		s.failf("expected not to be empty")
	}
}

// Contains is IterableSubject.contains(element), or StringSubject.contains
// for a string actual.
func (s *Subject) Contains(element any) {
	if str, ok := s.str(); ok {
		sub, _ := truthNormalize(element).(string)
		if !strings.Contains(str, sub) {
			s.failf("expected to contain: %q\nbut was            : %q", sub, str)
		}
		return
	}
	for _, e := range s.list() {
		if truthEqual(e, element) {
			return
		}
	}
	s.failf("expected to contain: %s\nbut was            : %s", show(element), s.showActual())
}

// DoesNotContain is IterableSubject.doesNotContain(element) or
// StringSubject.doesNotContain.
func (s *Subject) DoesNotContain(element any) {
	if str, ok := s.str(); ok {
		sub, _ := truthNormalize(element).(string)
		if strings.Contains(str, sub) {
			s.failf("expected not to contain: %q\nbut was                : %q", sub, str)
		}
		return
	}
	for _, e := range s.list() {
		if truthEqual(e, element) {
			s.failf("expected not to contain: %s\nbut was                : %s", show(element), s.showActual())
		}
	}
}

// ContainsNoDuplicates is IterableSubject.containsNoDuplicates().
func (s *Subject) ContainsNoDuplicates() {
	l := s.list()
	for i := range l {
		for j := i + 1; j < len(l); j++ {
			if truthEqual(l[i], l[j]) {
				s.failf("expected not to contain duplicates, but contained %s\nfull contents: %s", show(l[i]), s.showActual())
			}
		}
	}
}

// ContainsAnyOf is IterableSubject.containsAnyOf(first, second, rest...).
func (s *Subject) ContainsAnyOf(elements ...any) {
	s.ContainsAnyIn(elements)
}

// ContainsAnyIn is IterableSubject.containsAnyIn(expected).
func (s *Subject) ContainsAnyIn(expected any) {
	for _, e := range s.listOf(expected) {
		for _, a := range s.list() {
			if truthEqual(a, e) {
				return
			}
		}
	}
	s.failf("expected to contain any of: %s\nbut was                   : %s", show(expected), s.showActual())
}

// ContainsNoneOf is IterableSubject.containsNoneOf(first, second, rest...).
func (s *Subject) ContainsNoneOf(elements ...any) {
	s.ContainsNoneIn(elements)
}

// ContainsNoneIn is IterableSubject.containsNoneIn(excluded).
func (s *Subject) ContainsNoneIn(excluded any) {
	for _, e := range s.listOf(excluded) {
		for _, a := range s.list() {
			if truthEqual(a, e) {
				s.failf("expected not to contain any of: %s\nbut contained                 : %s\nfull contents: %s", show(excluded), show(e), s.showActual())
			}
		}
	}
}

// Ordered is Truth's Ordered: the result of containsExactly/containsAtLeast.
type Ordered struct {
	s        *Subject
	expected []any
	matches  func(actual, expected any) bool
	exact    bool
}

// InOrder is Ordered.inOrder().
func (o *Ordered) InOrder() {
	if o == nil {
		return // the containment assertion already failed
	}
	actual := o.s.list()
	j := 0
	for _, a := range actual {
		if j < len(o.expected) && o.matches(a, o.expected[j]) {
			j++
		} else if o.exact {
			break
		}
	}
	if j != len(o.expected) {
		o.s.failf("contents match, but order was wrong\nexpected: %s\nbut was : %s", show(o.expected), o.s.showActual())
	}
}

// ContainsExactly is IterableSubject.containsExactly(elements...).
func (s *Subject) ContainsExactly(elements ...any) *Ordered {
	return s.containsExactlyUsing(elements, truthEqual, "")
}

// ContainsExactlyElementsIn is IterableSubject.containsExactlyElementsIn(expected).
func (s *Subject) ContainsExactlyElementsIn(expected any) *Ordered {
	return s.containsExactlyUsing(s.listOf(expected), truthEqual, "")
}

func (s *Subject) containsExactlyUsing(expected []any, matches func(a, e any) bool, how string) *Ordered {
	actual := s.list()
	used := make([]bool, len(actual))
	var missing, unexpected []any
	for _, e := range expected {
		found := false
		for i, a := range actual {
			if !used[i] && matches(a, e) {
				used[i] = true
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, e)
		}
	}
	for i, u := range used {
		if !u {
			unexpected = append(unexpected, actual[i])
		}
	}
	if len(missing) > 0 || len(unexpected) > 0 {
		var sb strings.Builder
		if len(missing) > 0 {
			fmt.Fprintf(&sb, "missing (%d)   : %s\n", len(missing), show(missing))
		}
		if len(unexpected) > 0 {
			fmt.Fprintf(&sb, "unexpected (%d): %s\n", len(unexpected), show(unexpected))
		}
		fmt.Fprintf(&sb, "---\nexpected%s: %s\nbut was : %s", how, show(expected), s.showActual())
		s.failf("%s", sb.String())
	}
	return &Ordered{s: s, expected: expected, matches: matches, exact: true}
}

// ContainsAtLeast is IterableSubject.containsAtLeast(elements...).
func (s *Subject) ContainsAtLeast(elements ...any) *Ordered {
	return s.containsAtLeastUsing(elements, truthEqual)
}

// ContainsAtLeastElementsIn is IterableSubject.containsAtLeastElementsIn(expected).
func (s *Subject) ContainsAtLeastElementsIn(expected any) *Ordered {
	return s.containsAtLeastUsing(s.listOf(expected), truthEqual)
}

func (s *Subject) containsAtLeastUsing(expected []any, matches func(a, e any) bool) *Ordered {
	actual := s.list()
	used := make([]bool, len(actual))
	var missing []any
	for _, e := range expected {
		found := false
		for i, a := range actual {
			if !used[i] && matches(a, e) {
				used[i] = true
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, e)
		}
	}
	if len(missing) > 0 {
		s.failf("missing (%d): %s\n---\nexpected to contain at least: %s\nbut was                     : %s", len(missing), show(missing), show(expected), s.showActual())
	}
	return &Ordered{s: s, expected: expected, matches: matches}
}

// --- MapSubject

func (s *Subject) entries() []truthEntry {
	m, ok := truthAsMap(truthNormalize(s.actual))
	if !ok {
		s.failf("expected a map, but was %s", s.showActual())
	}
	return m
}

// ContainsKey is MapSubject.containsKey(key).
func (s *Subject) ContainsKey(key any) {
	if _, ok := truthMapLookup(s.entries(), key); !ok {
		s.failf("expected to contain key: %s\nbut was                : %s", show(key), s.showActual())
	}
}

// DoesNotContainKey is MapSubject.doesNotContainKey(key).
func (s *Subject) DoesNotContainKey(key any) {
	if _, ok := truthMapLookup(s.entries(), key); ok {
		s.failf("expected not to contain key: %s\nbut was                    : %s", show(key), s.showActual())
	}
}

// ContainsEntry is MapSubject.containsEntry(key, value).
func (s *Subject) ContainsEntry(key, value any) {
	if v, ok := truthMapLookup(s.entries(), key); !ok || !truthEqual(v, value) {
		s.failf("expected to contain entry: %s=%s\nbut was                  : %s", show(key), show(value), s.showActual())
	}
}

// DoesNotContainEntry is MapSubject.doesNotContainEntry(key, value).
func (s *Subject) DoesNotContainEntry(key, value any) {
	if v, ok := truthMapLookup(s.entries(), key); ok && truthEqual(v, value) {
		s.failf("expected not to contain entry: %s=%s\nbut was                      : %s", show(key), show(value), s.showActual())
	}
}

// ContainsExactlyEntriesIn is MapSubject.containsExactlyEntriesIn(expected).
func (s *Subject) ContainsExactlyEntriesIn(expected any) {
	exp, ok := truthAsMap(truthNormalize(expected))
	if !ok {
		s.failf("expected a map argument, but was %s", show(expected))
	}
	if !truthMapsEqual(s.entries(), exp) {
		s.failf("expected: %s\nbut was : %s", showMap(exp), s.showActual())
	}
}

// ContainsAtLeastEntriesIn is MapSubject.containsAtLeastEntriesIn(expected).
func (s *Subject) ContainsAtLeastEntriesIn(expected any) {
	exp, ok := truthAsMap(truthNormalize(expected))
	if !ok {
		s.failf("expected a map argument, but was %s", show(expected))
	}
	actual := s.entries()
	for _, e := range exp {
		if v, ok := truthMapLookup(actual, e.key); !ok || !truthEqual(v, e.value) {
			s.failf("missing entry: %s=%s\n---\nexpected to contain at least: %s\nbut was                     : %s", show(e.key), show(e.value), showMap(exp), s.showActual())
		}
	}
}

// --- Correspondence

// Correspondence is Truth's Correspondence<A, E>: how an actual element
// matches an expected one, with a description and an optional diff
// formatter.
type Correspondence struct {
	compare     func(actual, expected any) bool
	description string
	formatDiff  func(actual, expected any) string
}

// Correspondence_From is Correspondence.from(predicate, description).
func Correspondence_From[A, E any](predicate func(actual A, expected E) bool, description string) *Correspondence {
	return &Correspondence{
		compare: func(a, e any) bool {
			av, ok1 := truthCast[A](a)
			ev, ok2 := truthCast[E](e)
			return ok1 && ok2 && predicate(av, ev)
		},
		description: description,
	}
}

// Correspondence_Transforming is Correspondence.transforming(fn, description):
// an actual element matches if fn(actual) equals the expected element.
func Correspondence_Transforming[A, E any](fn func(actual A) E, description string) *Correspondence {
	return Correspondence_From(func(a A, e E) bool { return truthEqual(fn(a), e) }, description)
}

// FormattingDiffsUsing is Correspondence.formattingDiffsUsing(formatter).
func FormattingDiffsUsing[A, E any](c *Correspondence, formatter func(actual A, expected E) string) *Correspondence {
	return &Correspondence{
		compare:     c.compare,
		description: c.description,
		formatDiff: func(a, e any) string {
			av, _ := truthCast[A](a)
			ev, _ := truthCast[E](e)
			return formatter(av, ev)
		},
	}
}

// truthCast is a Java reference cast that lets null through.
func truthCast[T any](v any) (T, bool) {
	var zero T
	if v == nil {
		return zero, true
	}
	t, ok := v.(T)
	return t, ok
}

// String is Correspondence.toString().
func (c *Correspondence) String() string { return c.description }

// UsingCorrespondence is IterableSubject.UsingCorrespondence.
type UsingCorrespondence struct {
	s *Subject
	c *Correspondence
}

// ComparingElementsUsing is IterableSubject.comparingElementsUsing(correspondence).
func (s *Subject) ComparingElementsUsing(c *Correspondence) *UsingCorrespondence {
	return &UsingCorrespondence{s: s, c: c}
}

// ContainsExactlyElementsIn is UsingCorrespondence.containsExactlyElementsIn(expected).
func (u *UsingCorrespondence) ContainsExactlyElementsIn(expected any) *Ordered {
	return u.containsExactly(u.s.listOf(expected))
}

// ContainsExactly is UsingCorrespondence.containsExactly(expected...).
func (u *UsingCorrespondence) ContainsExactly(expected ...any) *Ordered {
	return u.containsExactly(expected)
}

func (u *UsingCorrespondence) containsExactly(expected []any) *Ordered {
	defer u.addDiffs(expected)
	return u.s.containsExactlyUsing(expected, u.c.compare, " (elements that "+u.c.description+")")
}

// addDiffs adds the formatted diffs of a failed assertion between an actual
// element and the expected element of the same index, like Truth's "diff"
// facts.
func (u *UsingCorrespondence) addDiffs(expected []any) {
	r := recover()
	if r == nil {
		return
	}
	ae, ok := r.(*jutil.AssertionError)
	if !ok || u.c.formatDiff == nil {
		panic(r)
	}
	msg := ae.GetMessage()
	actual := u.s.listOf(u.s.actual)
	for i := 0; i < len(actual) && i < len(expected); i++ {
		if !u.c.compare(actual[i], expected[i]) {
			msg += fmt.Sprintf("\ndiff (#%d): %s", i, u.c.formatDiff(actual[i], expected[i]))
		}
	}
	panic(jutil.NewAssertionError(msg))
}

// ContainsAtLeastElementsIn is UsingCorrespondence.containsAtLeastElementsIn(expected).
func (u *UsingCorrespondence) ContainsAtLeastElementsIn(expected any) *Ordered {
	return u.s.containsAtLeastUsing(u.s.listOf(expected), u.c.compare)
}

// ContainsAtLeast is UsingCorrespondence.containsAtLeast(expected...).
func (u *UsingCorrespondence) ContainsAtLeast(expected ...any) *Ordered {
	return u.s.containsAtLeastUsing(expected, u.c.compare)
}

// Contains is UsingCorrespondence.contains(expected).
func (u *UsingCorrespondence) Contains(expected any) {
	for _, a := range u.s.list() {
		if u.c.compare(a, expected) {
			return
		}
	}
	u.s.failf("expected to contain at least one element that %s %s\nbut was: %s", u.c.description, show(expected), u.s.showActual())
}

// DoesNotContain is UsingCorrespondence.doesNotContain(excluded).
func (u *UsingCorrespondence) DoesNotContain(excluded any) {
	for _, a := range u.s.list() {
		if u.c.compare(a, excluded) {
			u.s.failf("expected not to contain any element that %s %s\nbut contained %s", u.c.description, show(excluded), show(a))
		}
	}
}

// Fact is com.google.common.truth.Fact: a key and an optional value of a
// failure message.
type Fact struct {
	key   string
	value *string
}

// Fact_Fact is Fact.fact(key, value): the value is rendered with
// String.valueOf.
func Fact_Fact(key string, value any) Fact {
	v := show(value)
	return Fact{key: key, value: &v}
}

// Fact_SimpleFact is Fact.simpleFact(key): a fact with no value.
func Fact_SimpleFact(key string) Fact {
	return Fact{key: key}
}

// fact_makeMessage is Fact.makeMessage: the facts one per line, values
// aligned after the longest key, or each multi-line value indented on the
// lines after its key.
func fact_makeMessage(messages []string, facts []Fact) string {
	longestKeyLength := 0
	seenNewlineInValue := false
	for _, f := range facts {
		if f.value != nil {
			longestKeyLength = max(longestKeyLength, len(f.key))
			seenNewlineInValue = seenNewlineInValue || strings.Contains(*f.value, "\n")
		}
	}
	var b strings.Builder
	for _, m := range messages {
		b.WriteString(m)
		b.WriteString("\n")
	}
	for _, f := range facts {
		if f.value == nil {
			b.WriteString(f.key)
		} else if seenNewlineInValue {
			b.WriteString(f.key)
			b.WriteString(":\n    ")
			b.WriteString(strings.ReplaceAll(*f.value, "\n", "\n    "))
		} else {
			b.WriteString(f.key)
			b.WriteString(strings.Repeat(" ", longestKeyLength-len(f.key)))
			b.WriteString(": ")
			b.WriteString(*f.value)
		}
		b.WriteString("\n")
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// FailWithoutActual is Subject.failWithoutActual(Fact, Fact...): fails
// with the given facts (after the AssertWithMessage message, if any).
func (s *Subject) FailWithoutActual(facts ...Fact) {
	var messages []string
	if s.message != "" {
		messages = append(messages, s.message)
	}
	panic(jutil.NewAssertionError(fact_makeMessage(messages, facts)))
}

// TolerantDoubleComparison is Truth's DoubleSubject.TolerantDoubleComparison.
type TolerantDoubleComparison struct {
	s         *Subject
	tolerance float64
}

// IsWithin is DoubleSubject.isWithin(tolerance): call Of(expected) on the
// result.
func (s *Subject) IsWithin(tolerance float64) *TolerantDoubleComparison {
	return &TolerantDoubleComparison{s: s, tolerance: tolerance}
}

// truthCheckTolerance is DoubleSubject.checkTolerance.
func truthCheckTolerance(tolerance float64) {
	jutil.CheckArgumentMsg(!math.IsNaN(tolerance), "tolerance cannot be NaN")
	jutil.CheckArgumentMsg1(tolerance >= 0.0, "tolerance (%s) cannot be negative", tolerance)
	jutil.CheckArgumentMsg(tolerance != math.MaxFloat64, "tolerance cannot be Double.MAX_VALUE")
	jutil.CheckArgumentMsg(!math.IsInf(tolerance, 1), "tolerance cannot be POSITIVE_INFINITY")
}

// Of is TolerantDoubleComparison.of(expected): the actual double is
// finite and within the tolerance of the finite expected value.
func (c *TolerantDoubleComparison) Of(expected float64) {
	actual, ok := truthAsFloat(c.s.actual)
	if !ok {
		c.s.failf("expected a double\nbut was: %s", c.s.showActual())
	}
	truthCheckTolerance(c.tolerance)
	// MathUtil.equalWithinTolerance
	if !(math.Abs(actual-expected) <= math.Abs(c.tolerance)) {
		c.s.FailWithoutActual(
			Fact_Fact("expected", jutil.DoubleToString(expected)),
			Fact_Fact("but was", jutil.DoubleToString(actual)),
			Fact_Fact("outside tolerance", jutil.DoubleToString(c.tolerance)))
	}
}

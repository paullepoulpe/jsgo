// Package rhinotesting is the Go port of com.google.javascript.rhino.testing
// (closure-compiler v20260415): the base test case and Truth subjects used by
// the rhino and jscomp unit tests.
package rhinotesting

// Go-only: language-adapter — a minimal Google Truth (com.google.common.truth)
// for the ported tests. Java's Subject throws AssertionError on failure,
// which ends the test; here a failure calls t.Fatalf. Equality is
// java.util.Objects.equals (jutil.ObjectsEquals); lists compare element by
// element like List.equals, and a *string compares like the String it
// stands for.

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Truth carries the test that assertions report to.
type Truth struct {
	T testing.TB
}

// NewTruth returns assertions reporting to t.
func NewTruth(t testing.TB) Truth {
	return Truth{T: t}
}

// AssertThat is Truth.assertThat(actual).
func (tr Truth) AssertThat(actual any) *Subject {
	return &Subject{t: tr.T, actual: actual}
}

// AssertWithMessage is Truth.assertWithMessage(format, args...).
func (tr Truth) AssertWithMessage(format string, args ...any) *StandardSubjectBuilder {
	return &StandardSubjectBuilder{t: tr.T, message: jutil.LenientFormat(format, args...)}
}

// StandardSubjectBuilder is Truth's StandardSubjectBuilder: a message
// prefix for the next assertion.
type StandardSubjectBuilder struct {
	t       testing.TB
	message string
}

// That is StandardSubjectBuilder.that(actual).
func (b *StandardSubjectBuilder) That(actual any) *Subject {
	return &Subject{t: b.t, actual: actual, message: b.message}
}

// Subject is Truth's Subject (with the IterableSubject methods).
type Subject struct {
	t       testing.TB
	actual  any
	message string
}

// Actual is the value under test.
func (s *Subject) Actual() any { return s.actual }

func (s *Subject) fail(format string, args ...any) {
	s.t.Helper()
	msg := fmt.Sprintf(format, args...)
	if s.message != "" {
		msg = s.message + "\n" + msg
	}
	s.t.Fatal(msg)
}

// normalize maps Go stand-ins to the Java values they represent.
func normalize(v any) any {
	if p, ok := v.(*string); ok {
		if p == nil {
			return nil
		}
		return *p
	}
	return v
}

// Show renders v like Truth does (String.valueOf, lists as [a, b]).
func Show(v any) string {
	v = normalize(v)
	if jutil.IsNil(v) {
		return "null"
	}
	if elems, ok := AsList(v); ok {
		var sb strings.Builder
		sb.WriteByte('[')
		for i, e := range elems {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(Show(e))
		}
		sb.WriteByte(']')
		return sb.String()
	}
	return jutil.ToJavaString(v)
}

// AsList returns the elements of a Go slice or of a collection with a
// ToSlice method, in iteration order.
func AsList(v any) ([]any, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil, false
	}
	if rv.Kind() == reflect.Slice {
		out := make([]any, rv.Len())
		for i := range out {
			out[i] = rv.Index(i).Interface()
		}
		return out, true
	}
	if m := rv.MethodByName("ToSlice"); m.IsValid() && m.Type().NumIn() == 0 && m.Type().NumOut() == 1 {
		return AsList(m.Call(nil)[0].Interface())
	}
	return nil, false
}

// Equal is Objects.equals with List.equals for lists.
func Equal(a, b any) bool {
	a, b = normalize(a), normalize(b)
	if jutil.IsNil(a) || jutil.IsNil(b) {
		return jutil.IsNil(a) && jutil.IsNil(b)
	}
	la, aList := AsList(a)
	lb, bList := AsList(b)
	if aList && bList {
		if len(la) != len(lb) {
			return false
		}
		for i := range la {
			if !Equal(la[i], lb[i]) {
				return false
			}
		}
		return true
	}
	return jutil.ObjectsEquals(a, b)
}

// IsTrue is BooleanSubject.isTrue().
func (s *Subject) IsTrue() {
	s.t.Helper()
	if b, ok := s.actual.(bool); !ok || !b {
		s.fail("expected to be true, but was %s", Show(s.actual))
	}
}

// IsFalse is BooleanSubject.isFalse().
func (s *Subject) IsFalse() {
	s.t.Helper()
	if b, ok := s.actual.(bool); !ok || b {
		s.fail("expected to be false, but was %s", Show(s.actual))
	}
}

// IsEqualTo is Subject.isEqualTo(expected).
func (s *Subject) IsEqualTo(expected any) {
	s.t.Helper()
	if !Equal(s.actual, expected) {
		s.fail("expected: %s\nbut was : %s", Show(expected), Show(s.actual))
	}
}

// IsNotEqualTo is Subject.isNotEqualTo(unexpected).
func (s *Subject) IsNotEqualTo(unexpected any) {
	s.t.Helper()
	if Equal(s.actual, unexpected) {
		s.fail("expected not to be: %s", Show(unexpected))
	}
}

// IsNull is Subject.isNull().
func (s *Subject) IsNull() {
	s.t.Helper()
	if !jutil.IsNil(normalize(s.actual)) {
		s.fail("expected: null\nbut was : %s", Show(s.actual))
	}
}

// IsNotNull is Subject.isNotNull().
func (s *Subject) IsNotNull() {
	s.t.Helper()
	if jutil.IsNil(normalize(s.actual)) {
		s.fail("expected not to be null")
	}
}

func sameInstance(a, b any) bool {
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
		// The same Java collection instance: the same backing array.
		return ra.Pointer() == rb.Pointer() && ra.Len() == rb.Len()
	}
	if ra.Comparable() {
		return a == b
	}
	return false
}

// IsSameInstanceAs is Subject.isSameInstanceAs(expected).
func (s *Subject) IsSameInstanceAs(expected any) {
	s.t.Helper()
	if !sameInstance(s.actual, expected) {
		s.fail("expected specific instance: %s\nbut was           : %s", Show(expected), Show(s.actual))
	}
}

// IsNotSameInstanceAs is Subject.isNotSameInstanceAs(unexpected).
func (s *Subject) IsNotSameInstanceAs(unexpected any) {
	s.t.Helper()
	if sameInstance(s.actual, unexpected) {
		s.fail("expected not to be specific instance: %s", Show(unexpected))
	}
}

// IsAtLeast is ComparableSubject.isAtLeast(other) for ints.
func (s *Subject) IsAtLeast(other int) {
	s.t.Helper()
	if v, ok := s.actual.(int); !ok || v < other {
		s.fail("expected to be at least: %d\nbut was                : %s", other, Show(s.actual))
	}
}

func (s *Subject) list() []any {
	s.t.Helper()
	l, ok := AsList(s.actual)
	if !ok {
		if jutil.IsNil(s.actual) {
			return nil
		}
		s.fail("expected a collection, but was %s", Show(s.actual))
	}
	return l
}

// HasSize is IterableSubject.hasSize(n).
func (s *Subject) HasSize(n int) {
	s.t.Helper()
	if l := s.list(); len(l) != n {
		s.fail("expected to have size: %d\nbut was: %s", n, Show(s.actual))
	}
}

// IsEmpty is IterableSubject.isEmpty().
func (s *Subject) IsEmpty() {
	s.t.Helper()
	if str, ok := normalize(s.actual).(string); ok {
		if str != "" {
			s.fail("expected to be empty, but was %q", str)
		}
		return
	}
	if l := s.list(); len(l) != 0 {
		s.fail("expected to be empty, but was %s", Show(s.actual))
	}
}

// IsNotEmpty is IterableSubject.isNotEmpty().
func (s *Subject) IsNotEmpty() {
	s.t.Helper()
	if l := s.list(); len(l) == 0 {
		s.fail("expected not to be empty")
	}
}

// Contains is IterableSubject.contains(element).
func (s *Subject) Contains(element any) {
	s.t.Helper()
	for _, e := range s.list() {
		if Equal(e, element) {
			return
		}
	}
	s.fail("expected to contain: %s\nbut was: %s", Show(element), Show(s.actual))
}

// DoesNotContain is IterableSubject.doesNotContain(element).
func (s *Subject) DoesNotContain(element any) {
	s.t.Helper()
	for _, e := range s.list() {
		if Equal(e, element) {
			s.fail("expected not to contain: %s\nbut was: %s", Show(element), Show(s.actual))
		}
	}
}

// Ordered is Truth's Ordered: the result of a containsExactly assertion.
type Ordered struct {
	s        *Subject
	expected []any
}

// ContainsExactly is IterableSubject.containsExactly(elements...).
func (s *Subject) ContainsExactly(elements ...any) *Ordered {
	s.t.Helper()
	return s.containsExactly(elements)
}

// ContainsExactlyElementsIn is IterableSubject.containsExactlyElementsIn(expected).
func (s *Subject) ContainsExactlyElementsIn(expected any) *Ordered {
	s.t.Helper()
	l, ok := AsList(expected)
	if !ok && !jutil.IsNil(expected) {
		s.fail("expected a collection argument, but was %s", Show(expected))
	}
	return s.containsExactly(l)
}

func (s *Subject) containsExactly(expected []any) *Ordered {
	s.t.Helper()
	actual := s.list()
	used := make([]bool, len(actual))
	for _, e := range expected {
		found := false
		for i, a := range actual {
			if !used[i] && Equal(a, e) {
				used[i] = true
				found = true
				break
			}
		}
		if !found {
			s.fail("missing: %s\nexpected exactly: %s\nbut was: %s", Show(e), Show(expected), Show(s.actual))
		}
	}
	for i, u := range used {
		if !u {
			s.fail("unexpected: %s\nexpected exactly: %s\nbut was: %s", Show(actual[i]), Show(expected), Show(s.actual))
		}
	}
	return &Ordered{s: s, expected: expected}
}

// InOrder is Ordered.inOrder().
func (o *Ordered) InOrder() {
	o.s.t.Helper()
	actual := o.s.list()
	for i := range actual {
		if !Equal(actual[i], o.expected[i]) {
			o.s.fail("contents match, but order was wrong\nexpected: %s\nbut was : %s", Show(o.expected), Show(o.s.actual))
		}
	}
}

// IsInstanceOf is Subject.isInstanceOf(T.class).
func IsInstanceOf[T any](s *Subject) {
	s.t.Helper()
	if _, ok := s.actual.(T); !ok || jutil.IsNil(s.actual) {
		s.fail("expected instance of: %s\nbut was: %s (%T)", reflect.TypeOf((*T)(nil)).Elem(), Show(s.actual), s.actual)
	}
}

// IsNotInstanceOf is Subject.isNotInstanceOf(T.class).
func IsNotInstanceOf[T any](s *Subject) {
	s.t.Helper()
	if _, ok := s.actual.(T); ok && !jutil.IsNil(s.actual) {
		s.fail("expected not to be an instance of: %s\nbut was: %s (%T)", reflect.TypeOf((*T)(nil)).Elem(), Show(s.actual), s.actual)
	}
}

// AssertThrows is org.junit.Assert.assertThrows(E.class, body): body must
// throw an exception assignable to E, which is returned.
func AssertThrows[E any](t testing.TB, body func()) E {
	t.Helper()
	var caught E
	ok := false
	jutil.Try(body, func(e E) {
		caught = e
		ok = true
	})
	if !ok {
		t.Fatalf("expected %s to be thrown", reflect.TypeOf((*E)(nil)).Elem())
	}
	return caught
}

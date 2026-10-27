package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/testjunit"
)

type (
	JUnitTestCase            = testjunit.JUnitTestCase
	StandardSubjectBuilder   = testjunit.StandardSubjectBuilder
	Subject                  = testjunit.Subject
	Ordered                  = testjunit.Ordered
	Correspondence           = testjunit.Correspondence
	UsingCorrespondence      = testjunit.UsingCorrespondence
	Fact                     = testjunit.Fact
	TolerantDoubleComparison = testjunit.TolerantDoubleComparison
)

func RunJUnit4(t *testing.T, newTest func() JUnitTestCase) {
	testjunit.RunJUnit4(t, newTest)
}

func AssertThat(actual any) *Subject {
	return testjunit.AssertThat(actual)
}

func AssertWithMessage(format string, args ...any) *StandardSubjectBuilder {
	return testjunit.AssertWithMessage(format, args...)
}

func Fail(message string) {
	testjunit.Fail(message)
}

func AssertThrows[E any](body func()) E {
	return testjunit.AssertThrows[E](body)
}

func IsInstanceOf[T any](s *Subject) {
	testjunit.IsInstanceOf[T](s)
}

func IsNotInstanceOf[T any](s *Subject) {
	testjunit.IsNotInstanceOf[T](s)
}

func Correspondence_From[A, E any](predicate func(actual A, expected E) bool, description string) *Correspondence {
	return testjunit.Correspondence_From(predicate, description)
}

func Correspondence_Transforming[A, E any](fn func(actual A) E, description string) *Correspondence {
	return testjunit.Correspondence_Transforming(fn, description)
}

func Correspondence_TransformingBoth[A, E, T any](
	actualTransform func(actual A) T, expectedTransform func(expected E) T, description string) *Correspondence {
	return testjunit.Correspondence_TransformingBoth(actualTransform, expectedTransform, description)
}

func FormattingDiffsUsing[A, E any](c *Correspondence, formatter func(actual A, expected E) string) *Correspondence {
	return testjunit.FormattingDiffsUsing(c, formatter)
}

func Fact_Fact(key string, value any) Fact {
	return testjunit.Fact_Fact(key, value)
}

func Fact_SimpleFact(key string) Fact {
	return testjunit.Fact_SimpleFact(key)
}

func show(v any) string {
	return testjunit.Show(v)
}

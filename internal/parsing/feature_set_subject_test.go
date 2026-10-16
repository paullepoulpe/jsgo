// Port of com/google/javascript/jscomp/parsing/parser/testing/FeatureSetSubject.java @ closure-compiler v20260415

package parsing

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

// featureSetSubject is a Truth Subject for FeatureSet.
type featureSetSubject struct {
	t      *testing.T
	actual *parser.FeatureSet
}

// Java: FeatureSetSubject.assertFS(FeatureSet)
func assertFS(t *testing.T, fs *parser.FeatureSet) *featureSetSubject {
	return &featureSetSubject{t: t, actual: fs}
}

// Java: FeatureSetSubject.contains(FeatureSet)
func (s *featureSetSubject) Contains(other *parser.FeatureSet) {
	s.t.Helper()
	if !s.actual.Contains(other) {
		s.t.Fatal(jutil.LenientFormat("Expected a FeatureSet containing: %s\nBut got: %s", other, s.actual))
	}
}

// Java: FeatureSetSubject.containsNoneOf(FeatureSet)
func (s *featureSetSubject) ContainsNoneOf(other *parser.FeatureSet) {
	s.t.Helper()
	if !other.Without(s.actual).Equals(other) {
		s.t.Fatal(jutil.LenientFormat("Expected a FeatureSet containing none of: %s\nBut got: %s", other, s.actual))
	}
}

// Java: FeatureSetSubject.has(Feature)
func (s *featureSetSubject) Has(feature parser.FeatureSet_Feature) {
	s.t.Helper()
	if !s.actual.Has(feature) {
		s.t.Fatal(jutil.LenientFormat("Expected a FeatureSet that has: %s\nBut got: %s", feature, s.actual))
	}
}

// Java: FeatureSetSubject.doesNotHave(Feature)
func (s *featureSetSubject) DoesNotHave(feature parser.FeatureSet_Feature) {
	s.t.Helper()
	if s.actual.Has(feature) {
		s.t.Fatal(jutil.LenientFormat("Expected a FeatureSet that doesn't have: %s\nBut got: %s", feature, s.actual))
	}
}

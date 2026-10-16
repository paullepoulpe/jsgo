// Port of com/google/javascript/jscomp/parsing/parser/FeatureSetTest.java @ closure-compiler v20260415

package parser

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// featureSetSubject stands in for parser.testing.FeatureSetSubject (Truth),
// which is not part of the pinned test sources: has/doesNotHave check
// FeatureSet.has, contains checks FeatureSet.contains(FeatureSet) and equals
// checks FeatureSet.equals.
type featureSetSubject struct {
	t      *testing.T
	actual *FeatureSet
}

// assertFS is FeatureSetSubject.assertFS(FeatureSet).
func assertFS(t *testing.T, fs *FeatureSet) *featureSetSubject {
	t.Helper()
	return &featureSetSubject{t: t, actual: fs}
}

func (s *featureSetSubject) has(feature FeatureSet_Feature) {
	s.t.Helper()
	if !s.actual.Has(feature) {
		s.t.Errorf("Expected a FeatureSet that has: %s\nBut got: %s", feature, s.actual)
	}
}

func (s *featureSetSubject) doesNotHave(feature FeatureSet_Feature) {
	s.t.Helper()
	if s.actual.Has(feature) {
		s.t.Errorf("Expected a FeatureSet that does not have: %s\nBut got: %s", feature, s.actual)
	}
}

func (s *featureSetSubject) contains(other *FeatureSet) {
	s.t.Helper()
	if !s.actual.Contains(other) {
		s.t.Errorf("Expected a FeatureSet containing: %s\nBut got: %s", other, s.actual)
	}
}

func (s *featureSetSubject) equals(other *FeatureSet) {
	s.t.Helper()
	if !s.actual.Equals(other) {
		s.t.Errorf("Expected a FeatureSet equal to: %s\nBut got: %s", other, s.actual)
	}
}

func assertVersion(t *testing.T, fs *FeatureSet, want string) {
	t.Helper()
	if got := fs.Version(); got != want {
		t.Errorf("version() = %q, want %q", got, want)
	}
}

func TestFeatureSetContains(t *testing.T) {
	assertFS(t, FeatureSet_ALL).has(FeatureSet_Feature_TYPE_ANNOTATION)
	assertFS(t, FeatureSet_ALL).has(FeatureSet_Feature_MODULES)
}

func TestFeatureSetWithoutModules(t *testing.T) {
	assertFS(t, FeatureSet_ALL.WithoutFeatures(FeatureSet_Feature_MODULES)).has(FeatureSet_Feature_TYPE_ANNOTATION)
	assertFS(t, FeatureSet_ALL.WithoutFeatures(FeatureSet_Feature_MODULES)).doesNotHave(FeatureSet_Feature_MODULES)
}

func TestFeatureSetWithoutTypes(t *testing.T) {
	assertFS(t, FeatureSet_ALL.WithoutTypes()).doesNotHave(FeatureSet_Feature_TYPE_ANNOTATION)
	assertFS(t, FeatureSet_ALL.WithoutTypes()).has(FeatureSet_Feature_MODULES)
}

func TestFeatureSetEsOrdering(t *testing.T) {
	assertFS(t, FeatureSet_ALL).contains(FeatureSet_ES_UNSUPPORTED)
	assertFS(t, FeatureSet_ES_UNSUPPORTED).contains(FeatureSet_ES_UNSTABLE)
	assertFS(t, FeatureSet_ES_UNSTABLE).contains(FeatureSet_ES_NEXT)
	assertFS(t, FeatureSet_ES_NEXT).contains(FeatureSet_ES2022)
	assertFS(t, FeatureSet_ES2022).contains(FeatureSet_ES2021)
	assertFS(t, FeatureSet_ES2021).contains(FeatureSet_ES2020)
	assertFS(t, FeatureSet_ES2020).contains(FeatureSet_ES2019)
	assertFS(t, FeatureSet_ES2019).contains(FeatureSet_ES2018)
	assertFS(t, FeatureSet_ES2018).contains(FeatureSet_ES2017)
	assertFS(t, FeatureSet_ES2017).contains(FeatureSet_ES2016)
	assertFS(t, FeatureSet_ES2016).contains(FeatureSet_ES2015)
	assertFS(t, FeatureSet_ES2015).contains(FeatureSet_ES5)
	assertFS(t, FeatureSet_ES5).contains(FeatureSet_ES3)
	assertFS(t, FeatureSet_ES3).contains(FeatureSet_BARE_MINIMUM)
}

func TestFeatureSetEsModuleOrdering(t *testing.T) {
	assertFS(t, FeatureSet_ES2022_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)).equals(FeatureSet_ES2022)
	assertFS(t, FeatureSet_ES2021_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)).equals(FeatureSet_ES2021)
	assertFS(t, FeatureSet_ES2020_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)).equals(FeatureSet_ES2020)
	assertFS(t, FeatureSet_ES2019_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)).equals(FeatureSet_ES2019)
	assertFS(t, FeatureSet_ES2018_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)).equals(FeatureSet_ES2018)
	assertFS(t, FeatureSet_ES2017_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)).equals(FeatureSet_ES2017)
	assertFS(t, FeatureSet_ES2016_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)).equals(FeatureSet_ES2016)
	assertFS(t, FeatureSet_ES2015_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)).equals(FeatureSet_ES2015)
}

func TestFeatureSetVersionForDebugging(t *testing.T) {
	// ES_NEXT, ES_UNSUPPORTED are tested separately - see below
	assertVersion(t, FeatureSet_ES3, "es3")
	assertVersion(t, FeatureSet_ES5, "es5")
	assertVersion(t, FeatureSet_ES2015, "es6")
	assertVersion(t, FeatureSet_ES2015_MODULES, "es6")
	assertVersion(t, FeatureSet_ES2016, "es7")
	assertVersion(t, FeatureSet_ES2016_MODULES, "es7")
	assertVersion(t, FeatureSet_ES2017, "es8")
	assertVersion(t, FeatureSet_ES2017_MODULES, "es8")
	assertVersion(t, FeatureSet_ES2018, "es9")
	assertVersion(t, FeatureSet_ES2018_MODULES, "es9")
	assertVersion(t, FeatureSet_ES2019, "es_2019")
	assertVersion(t, FeatureSet_ES2019_MODULES, "es_2019")
	assertVersion(t, FeatureSet_ES2020, "es_2020")
	assertVersion(t, FeatureSet_ES2020_MODULES, "es_2020")
	assertVersion(t, FeatureSet_ES2021, "es_2021")
	assertVersion(t, FeatureSet_ES2021_MODULES, "es_2021")
	assertVersion(t, FeatureSet_ES2022, "es_2022")
	assertVersion(t, FeatureSet_ES2022_MODULES, "es_2022")
	assertVersion(t, FeatureSet_ALL, "all")
}

func TestFeatureSetEsNextAndNewer(t *testing.T) {
	// ES_NEXT and ES_UNSTABLE contain *_RUNTIME features to make these feature sets NOT moving
	// targets because polyfills can use "es_next" or "es_unstable" as their fromLang and we don't
	// want to incorrectly prune the polyfills in the case that no features remain in the set.
	assertVersion(t, FeatureSet_ES_NEXT, "es_next")
	assertVersion(t, FeatureSet_ES_UNSTABLE, "es_unstable")

	// ES_UNSUPPORTED is a moving target that may not have any unique features in which case it will
	// resolve to "es_unstable". This will change as new features are added and removed.
	assertVersion(t, FeatureSet_ES_UNSUPPORTED, "es_unsupported")
}

func TestFeatureSetValueOf(t *testing.T) {
	assertFS(t, FeatureSet_ValueOf("es3")).equals(FeatureSet_ES3)
	assertFS(t, FeatureSet_ValueOf("es5")).equals(FeatureSet_ES5)
	assertFS(t, FeatureSet_ValueOf("es6")).equals(FeatureSet_ES2015)
	assertFS(t, FeatureSet_ValueOf("es7")).equals(FeatureSet_ES2016)
	assertFS(t, FeatureSet_ValueOf("es8")).equals(FeatureSet_ES2017)
	assertFS(t, FeatureSet_ValueOf("es_2018")).equals(FeatureSet_ES2018)
	assertFS(t, FeatureSet_ValueOf("es9")).equals(FeatureSet_ES2018)
	assertFS(t, FeatureSet_ValueOf("es_2019")).equals(FeatureSet_ES2019)
	assertFS(t, FeatureSet_ValueOf("es_2020")).equals(FeatureSet_ES2020)
	assertFS(t, FeatureSet_ValueOf("es_2021")).equals(FeatureSet_ES2021)
	assertFS(t, FeatureSet_ValueOf("es_2022")).equals(FeatureSet_ES2022)
	assertFS(t, FeatureSet_ValueOf("es_next")).equals(FeatureSet_ES_NEXT)
	assertFS(t, FeatureSet_ValueOf("es_unstable")).equals(FeatureSet_ES_UNSTABLE)
	assertFS(t, FeatureSet_ValueOf("es_unsupported")).equals(FeatureSet_ES_UNSUPPORTED)
	assertFS(t, FeatureSet_ValueOf("all")).equals(FeatureSet_ALL)
	threw := false
	jutil.Try(func() { FeatureSet_ValueOf("bad") }, func(e jutil.IIllegalArgumentException) { threw = true })
	if !threw {
		t.Errorf("FeatureSet_ValueOf(\"bad\") did not throw IllegalArgumentException")
	}
}

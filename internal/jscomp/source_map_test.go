// Port of test/com/google/javascript/jscomp/SourceMapTest.java @ closure-compiler v20260415
//
// SourceMapTest extends SourceMapTestCase, whose compile() runs the whole
// Compiler and CodePrinter; every test is skipped until they are ported. The
// location mappings each test configures are built here, so that the
// configuration half of each test already runs.

package jscomp

import "testing"

// sourceMapTest_lambdaMapping is the lambda (location) -> "mapped/" + location.
type sourceMapTest_lambdaMapping struct{ _ byte }

func (sourceMapTest_lambdaMapping) Map(location string) *string {
	s := "mapped/" + location
	return &s
}

func sourceMapTest_skip(t *testing.T, mappings []SourceMap_LocationMapping) {
	t.Helper()
	options := NewCompilerOptions()
	if mappings != nil {
		options.SetSourceMapLocationMappings(mappings)
	}
	t.Skip("needs Compiler and SourceMapTestCase.compile (CodePrinter)")
}

func TestSourceMap_PrefixReplacement1(t *testing.T) {
	// mapping can be used to remove a prefix
	sourceMapTest_skip(t, []SourceMap_LocationMapping{NewSourceMap_PrefixLocationMapping("pre/", "")})
}

func TestSourceMap_PrefixReplacement2(t *testing.T) {
	// mapping can be used to replace a prefix
	sourceMapTest_skip(t, []SourceMap_LocationMapping{NewSourceMap_PrefixLocationMapping("pre/file", "src")})
}

func TestSourceMap_PrefixReplacement3(t *testing.T) {
	// multiple mappings can be applied
	sourceMapTest_skip(t, []SourceMap_LocationMapping{
		NewSourceMap_PrefixLocationMapping("file1", "x"),
		NewSourceMap_PrefixLocationMapping("file2", "y"),
	})
}

func TestSourceMap_PrefixReplacement4(t *testing.T) {
	// first match wins
	sourceMapTest_skip(t, []SourceMap_LocationMapping{
		NewSourceMap_PrefixLocationMapping("file1", "x"),
		NewSourceMap_PrefixLocationMapping("file", "y"),
	})
}

func TestSourceMap_LambdaReplacement(t *testing.T) {
	sourceMapTest_skip(t, []SourceMap_LocationMapping{sourceMapTest_lambdaMapping{}})
}

// This is taken from SourceMapJsLangTest.
func TestSourceMap_RepeatedCompilation(t *testing.T) {
	sourceMapTest_skip(t, nil)
}

func TestSourceMap_IntermediateFilesOmitUnmappedCode(t *testing.T) {
	sourceMapTest_skip(t, nil)
}

func TestSourceMap_EmptySourceMapsAreIgnored(t *testing.T) {
	sourceMapTest_skip(t, nil)
}

// Go-only: test-helper — PrefixLocationMapping's own behaviour (map, toString,
// equals, hashCode), which the skipped tests rely on.
func TestSourceMap_PrefixLocationMapping(t *testing.T) {
	m := NewSourceMap_PrefixLocationMapping("pre/file", "src")
	if got := m.Map("pre/file1"); got == nil || *got != "src1" {
		t.Errorf("map(pre/file1) = %v", got)
	}
	if got := m.Map("other"); got != nil {
		t.Errorf("map(other) = %q, want nil", *got)
	}
	assertEqual(t, "toString", m.String(), "(pre/file|src)")
	if !m.Equals(NewSourceMap_PrefixLocationMapping("pre/file", "src")) ||
		m.Equals(NewSourceMap_PrefixLocationMapping("pre/file", "x")) {
		t.Errorf("equals")
	}
}

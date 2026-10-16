// Port of com/google/javascript/rhino/testing/TestErrorReporter.java @ closure-compiler v20260415

package rhinotesting

import (
	"slices"
	"strconv"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TestErrorReporter is an error reporter for testing that verifies that
// messages reported to the reporter are expected.
//
// Use expectAllErrors/expectAllWarnings before the code that reports them,
// and VerifyHasEncounteredAllWarningsAndErrors to check that all were seen.
type TestErrorReporter struct {
	expectedErrors   []string
	expectedWarnings []string
	// Seen diagnostics are stored in a LinkedHashMap to deduplicate
	// identical diagnostics. The default
	// com.google.javascript.jscomp.SortingErrorManager behaves in this way.
	seenErrors   *jutil.LinkedHashMap[string, string]
	seenWarnings *jutil.LinkedHashMap[string, string]
}

// NewTestErrorReporter returns a reporter that expects no diagnostics.
//
// Java: TestErrorReporter()
func NewTestErrorReporter() *TestErrorReporter {
	return &TestErrorReporter{
		expectedErrors:   []string{},
		expectedWarnings: []string{},
		seenErrors:       jutil.NewLinkedHashMap[string, string](),
		seenWarnings:     jutil.NewLinkedHashMap[string, string](),
	}
}

// Java: TestErrorReporter.error(String, String, int, int)
func (r *TestErrorReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.seenErrors.Put(r.fmtDiagnosticKey(message, sourceName, line, lineOffset), message)
}

// Java: TestErrorReporter.warning(String, String, int, int)
func (r *TestErrorReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.seenWarnings.Put(r.fmtDiagnosticKey(message, sourceName, line, lineOffset), message)
}

// Java: TestErrorReporter.expectAllErrors(String...)
func (r *TestErrorReporter) ExpectAllErrors(errors ...string) *TestErrorReporter {
	r.expectedErrors = append(r.expectedErrors, errors...)
	return r
}

// Java: TestErrorReporter.expectAllWarnings(String...)
func (r *TestErrorReporter) ExpectAllWarnings(warnings ...string) *TestErrorReporter {
	r.expectedWarnings = append(r.expectedWarnings, warnings...)
	return r
}

// VerifyHasEncounteredAllWarningsAndErrors checks that exactly the expected
// diagnostics were reported, in order. Go-only: language-adapter — takes
// the test that Truth reports to.
//
// Java: TestErrorReporter.verifyHasEncounteredAllWarningsAndErrors()
func (r *TestErrorReporter) VerifyHasEncounteredAllWarningsAndErrors(t testing.TB) {
	t.Helper()
	tr := Truth{T: t}
	tr.AssertThat(slices.Collect(r.seenWarnings.Values())).ContainsExactlyElementsIn(r.expectedWarnings).InOrder()
	tr.AssertThat(slices.Collect(r.seenErrors.Values())).ContainsExactlyElementsIn(r.expectedErrors).InOrder()
}

// Java: TestErrorReporter.fmtDiagnosticKey(String, String, int, int)
func (r *TestErrorReporter) fmtDiagnosticKey(message string, sourceName string, line int, lineOffset int) string {
	return message + ":" + sourceName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(lineOffset)
}

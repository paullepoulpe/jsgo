// Port of test/com/google/javascript/jscomp/RhinoErrorReporterTest.java @ closure-compiler v20260415

package jscomp

import "testing"

func TestRhinoErrorReporter_MissingTypeWarnings(t *testing.T) {
	t.Skip("needs Compiler")
}

func TestRhinoErrorReporter_MissingCurlyBraceWarning(t *testing.T) {
	t.Skip("needs Compiler")
}

func TestRhinoErrorReporter_LanguageFeatureInHigherLanguageInError(t *testing.T) {
	t.Skip("needs Compiler")
}

// jsaRecordingErrorHandler records the reports it receives.
type jsaRecordingErrorHandler struct {
	levels []*CheckLevel
	errors []*JSError
}

func (h *jsaRecordingErrorHandler) Report(level *CheckLevel, error *JSError) {
	h.levels = append(h.levels, level)
	h.errors = append(h.errors, error)
}

// Go-only: extra coverage — the reporter maps Rhino messages to diagnostic
// types and keeps the default level only for unmapped messages.
func TestRhinoErrorReporter_ForOldRhino(t *testing.T) {
	h := &jsaRecordingErrorHandler{}
	r := RhinoErrorReporter_ForOldRhino(h)
	r.Warning("Missing type declaration.", "a.js", 1, 4)
	r.Error("something odd", "b.js", 2, 3)
	r.Warning("Duplicate parameter name \"x\"", "c.js", 3, 0)
	jsaEq(t, len(h.errors), 3)
	jsaTrue(t, h.levels[0] == nil && h.levels[1] == nil)
	jsaTrue(t, h.errors[0].Type() == RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING)
	jsaEq(t, h.errors[0].GetDefaultLevel(), CheckLevel_OFF)
	jsaEq(t, *h.errors[0].SourceName(), "a.js")
	jsaEq(t, h.errors[0].Charno(), 4)
	jsaTrue(t, h.errors[1].Type() == RhinoErrorReporter_PARSE_ERROR)
	jsaEq(t, h.errors[1].Description(), "Parse error. something odd")
	jsaEq(t, h.errors[1].GetDefaultLevel(), CheckLevel_ERROR)
	jsaTrue(t, h.errors[2].Type() == RhinoErrorReporter_DUPLICATE_PARAM)
	jsaEq(t, h.errors[2].GetDefaultLevel(), CheckLevel_ERROR)
	jsaTrue(t, RhinoErrorReporter_MapError("no match") == nil)
	jsaTrue(t, RhinoErrorReporter_MapError("Too many template parameters\nx") == RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS)
	jsaTrue(t, RhinoErrorReporter_MapError("a Unknown type b\nc") == RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}

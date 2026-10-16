// Port of com/google/javascript/jscomp/DiagnosticGroupPathSuppressingWarningsGuard.java @ closure-compiler v20260415

package jscomp

import "strings"

// DiagnosticGroupPathSuppressingWarningsGuard is a warnings guard that
// suppresses warnings for a particular diagnostic group for a file that
// contains the specified substring.
type DiagnosticGroupPathSuppressingWarningsGuard struct {
	WarningsGuardBase
	group *DiagnosticGroup
	part  string
}

var _ WarningsGuard = (*DiagnosticGroupPathSuppressingWarningsGuard)(nil)

// Java: DiagnosticGroupPathSuppressingWarningsGuard(DiagnosticGroup, String)
func NewDiagnosticGroupPathSuppressingWarningsGuard(group *DiagnosticGroup, part string) *DiagnosticGroupPathSuppressingWarningsGuard {
	return &DiagnosticGroupPathSuppressingWarningsGuard{group: group, part: part}
}

// Level does not touch warnings in other paths. It may return nil.
//
// Java: DiagnosticGroupPathSuppressingWarningsGuard.level(JSError)
func (g *DiagnosticGroupPathSuppressingWarningsGuard) Level(error *JSError) *CheckLevel {
	if error.SourceName() == nil || !strings.Contains(*error.SourceName(), g.part) {
		return nil
	}
	if g.group.Matches(error) {
		return CheckLevel_OFF.Ptr()
	}
	return nil
}

// Java: DiagnosticGroupPathSuppressingWarningsGuard.toString()
func (g *DiagnosticGroupPathSuppressingWarningsGuard) String() string {
	return g.group.String() + "(" + g.part + ")"
}

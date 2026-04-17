// Port of com/google/javascript/jscomp/JSError.java and
// com/google/javascript/jscomp/DiagnosticType.java @ closure-compiler v20260415
//
// JSError in Java carries source position, diagnostic type (code +
// default-level), and format-substituted arguments. For Tier-0 we track
// only the minimum needed for tuple-level conformance (path, line, col,
// code) plus a human-readable message shaped like Java's output so
// dogfooding the CLI feels right.

package diag

import "fmt"

type Level int

const (
	Warning Level = iota
	Error
)

func (l Level) String() string {
	if l == Error {
		return "ERROR"
	}
	return "WARNING"
}

type Diagnostic struct {
	Path    string
	Line    int // 1-indexed, matches Java Node.getLineno()
	Col     int // 0-indexed, matches Java Node.getCharno()
	Code    string
	Level   Level
	Message string
}

// FormatJavaLike reproduces the shape of Closure Compiler's default
// output ("<path>:<line>:<col>: <LEVEL> - [<CODE>] <msg>"). Strings are
// not part of the conformance check (see SPEC.md), but matching the shape
// keeps the CLI comfortable to use.
func (d Diagnostic) FormatJavaLike() string {
	return fmt.Sprintf("%s:%d:%d: %s - [%s] %s", d.Path, d.Line, d.Col, d.Level, d.Code, d.Message)
}

// Well-known diagnostic codes. Only Tier-0 entries live here; add more as
// corpus grows. Codes must match the JSC_* identifiers used by Java
// Closure Compiler (see DiagnosticType registrations in TypeCheck.java).
const (
	JSC_TYPE_MISMATCH = "JSC_TYPE_MISMATCH"
)

// Port of com/google/javascript/jscomp/JSError.java +
// com/google/javascript/jscomp/DiagnosticType.java
// @ closure-compiler v20260415 (subset).
//
// Java's JSError carries position, a DiagnosticType (code +
// default-level + format template), substituted args, and links to
// the source file. DiagnosticType is a full registry with
// per-category warning levels. This file tracks only the minimum
// the tuple-level conformance check consumes: (path, line, col,
// code) plus a Java-shaped message so CLI output reads the same.
// A proper DiagnosticType registry lands when we want to honour
// `--jscomp_warning=<group>` style flags.

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

// Well-known diagnostic codes. Add new ones as the corpus grows. Codes
// must match the JSC_* identifiers used by Java Closure Compiler (see
// DiagnosticType registrations in TypeCheck.java / VarCheck.java /
// VariableReferenceCheck.java).
const (
	JSC_TYPE_MISMATCH               = "JSC_TYPE_MISMATCH"
	JSC_INEXISTENT_PROPERTY         = "JSC_INEXISTENT_PROPERTY"
	JSC_UNDEFINED_VARIABLE          = "JSC_UNDEFINED_VARIABLE"
	JSC_VAR_MULTIPLY_DECLARED_ERROR = "JSC_VAR_MULTIPLY_DECLARED_ERROR"
	JSC_MISSING_RETURN_STATEMENT    = "JSC_MISSING_RETURN_STATEMENT"
	JSC_USELESS_CODE                = "JSC_USELESS_CODE"
	JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH = "JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH"
	JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH  = "JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH"
	JSC_ILLEGAL_PROPERTY_ACCESS             = "JSC_ILLEGAL_PROPERTY_ACCESS"
	JSC_NOT_A_CONSTRUCTOR                   = "JSC_NOT_A_CONSTRUCTOR"
)

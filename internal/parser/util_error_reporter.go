// Port of com/google/javascript/jscomp/parsing/parser/util/ErrorReporter.java @ closure-compiler v20260415

package parser

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ErrorReporter is a conduit for reporting errors and warnings to the user.
//
// Java: abstract class ErrorReporter. Subclasses embed ErrorReporterBase,
// call InitErrorReporterBase(self) and implement ReportError and
// ReportWarning (the abstract two-argument overloads).
type ErrorReporter interface {
	GetErrorReporterBase() *ErrorReporterBase
	ReportErrorFormat(location *SourcePosition, format string, arguments ...any)
	ReportError(location *SourcePosition, message string)
	ReportWarningFormat(location *SourcePosition, format string, arguments ...any)
	ReportWarning(location *SourcePosition, message string)
	HadError() bool
}

// ErrorReporterBase holds ErrorReporter's state and final methods.
type ErrorReporterBase struct {
	self ErrorReporter

	hadError bool
}

// InitErrorReporterBase initializes the base; self is the concrete reporter
// whose ReportError/ReportWarning the final methods dispatch to.
//
// Java: ErrorReporter()
func (e *ErrorReporterBase) InitErrorReporterBase(self ErrorReporter) {
	e.self = self
}

// GetErrorReporterBase returns the base struct.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (e *ErrorReporterBase) GetErrorReporterBase() *ErrorReporterBase {
	return e
}

// ReportErrorFormat formats the message with String.format semantics and
// reports it as an error.
//
// Java: ErrorReporter.reportError(SourcePosition, String, Object...)
func (e *ErrorReporterBase) ReportErrorFormat(location *SourcePosition, format string, arguments ...any) {
	e.hadError = true
	message := jutil.StringFormat(format, arguments...)
	e.self.ReportError(location, message)
}

// ReportWarningFormat formats the message with String.format semantics and
// reports it as a warning.
//
// Java: ErrorReporter.reportWarning(SourcePosition, String, Object...)
func (e *ErrorReporterBase) ReportWarningFormat(location *SourcePosition, format string, arguments ...any) {
	message := jutil.StringFormat(format, arguments...)
	e.self.ReportWarning(location, message)
}

// HadError reports whether an error was reported.
//
// Java: ErrorReporter.hadError()
func (e *ErrorReporterBase) HadError() bool {
	return e.hadError
}

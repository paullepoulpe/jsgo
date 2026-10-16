// Port of com/google/javascript/rhino/ErrorReporter.java @ closure-compiler v20260415

package rhino

// ErrorReporter defines a protocol for reporting errors during JavaScript
// translation or execution. Mirrors com.google.javascript.rhino.ErrorReporter.
type ErrorReporter interface {
	// Warning reports a warning. The implementing type may choose to ignore it.
	// Param: message a String describing the warning
	// Param: sourceName a String describing the JavaScript source where the warning occurred
	// Param: line the line number associated with the warning
	// Param: lineOffset the offset into lineSource where problem was detected
	Warning(message, sourceName string, line, lineOffset int)

	// Error reports an error. The implementing type is free to panic if it desires.
	// Param: message a String describing the error
	// Param: sourceName a String describing the JavaScript source where the error occurred
	// Param: line the line number associated with the error
	// Param: lineOffset the offset into lineSource where problem was detected
	Error(message, sourceName string, line, lineOffset int)
}

// NullErrorReporter is an ErrorReporter that silently discards all warnings
// and errors. Mirrors ErrorReporter.NULL_INSTANCE.
var NullErrorReporter ErrorReporter = &nullErrorReporter{}

type nullErrorReporter struct{}

func (n *nullErrorReporter) Warning(_, _ string, _, _ int) {}
func (n *nullErrorReporter) Error(_, _ string, _, _ int)   {}

// AlwaysThrowsErrorReporter is an ErrorReporter that panics on any warning or
// error. Mirrors ErrorReporter.ALWAYS_THROWS_INSTANCE.
var AlwaysThrowsErrorReporter ErrorReporter = &alwaysThrowsErrorReporter{}

type alwaysThrowsErrorReporter struct{}

func (a *alwaysThrowsErrorReporter) Warning(message, _ string, _, _ int) {
	panic(message)
}

func (a *alwaysThrowsErrorReporter) Error(message, _ string, _, _ int) {
	panic(message)
}

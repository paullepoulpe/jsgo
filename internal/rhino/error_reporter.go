// Port of com/google/javascript/rhino/ErrorReporter.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ErrorReporter defines a protocol for the reporting of errors during
// JavaScript translation or execution.
type ErrorReporter interface {
	// Warning reports a warning.
	//
	// The implementing class may choose to ignore the warning if it desires.
	//
	// message: a String describing the warning. sourceName: a String
	// describing the JavaScript source where the warning occurred; typically a
	// filename or URL. line: the line number associated with the warning.
	// lineOffset: the offset into lineSource where problem was detected.
	Warning(message string, sourceName string, line int, lineOffset int)

	// Error reports an error.
	//
	// The implementing class is free to throw an exception if it desires.
	//
	// If execution has not yet begun, the JavaScript engine is free to find
	// additional errors rather than terminating the translation. It will not
	// execute a script that had errors, however.
	//
	// message: a String describing the error. sourceName: a String
	// describing the JavaScript source where the error occurred; typically a
	// filename or URL. line: the line number associated with the error.
	// lineOffset: the offset into lineSource where problem was detected.
	Error(message string, sourceName string, line int, lineOffset int)
}

type errorReporter_nullInstance struct{}

func (errorReporter_nullInstance) Warning(message string, sourceName string, line int, lineOffset int) {
}

func (errorReporter_nullInstance) Error(message string, sourceName string, line int, lineOffset int) {
}

// Java: ErrorReporter.NULL_INSTANCE
var ErrorReporter_NULL_INSTANCE ErrorReporter = errorReporter_nullInstance{}

type errorReporter_alwaysThrowsInstance struct{}

func (errorReporter_alwaysThrowsInstance) Warning(message string, sourceName string, line int, lineOffset int) {
	panic(jutil.NewIllegalStateExceptionMsg(message))
}

func (errorReporter_alwaysThrowsInstance) Error(message string, sourceName string, line int, lineOffset int) {
	panic(jutil.NewIllegalStateExceptionMsg(message))
}

// Java: ErrorReporter.ALWAYS_THROWS_INSTANCE
var ErrorReporter_ALWAYS_THROWS_INSTANCE ErrorReporter = errorReporter_alwaysThrowsInstance{}

// Port of com/google/javascript/jscomp/parsing/parser/util/ErrorReporter.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parser

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ErrorReporter is a conduit for reporting errors and warnings to the user.
//
// Java: abstract class ErrorReporter. Subclasses embed ErrorReporterBase,
// call InitErrorReporterBase(self) and implement ReportError and
// ReportWarning (the abstract two-argument overloads).
type ErrorReporter interface {
	ReportErrorFormat(location LinePosition, format string, arguments ...any)
	ReportError(location LinePosition, message string)
	ReportWarningFormat(location LinePosition, format string, arguments ...any)
	ReportWarning(location LinePosition, message string)
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

// ReportErrorFormat formats the message with String.format semantics and
// reports it as an error.
//
// Java: ErrorReporter.reportError(SourcePosition, String, Object...)
func (e *ErrorReporterBase) ReportErrorFormat(location LinePosition, format string, arguments ...any) {
	e.hadError = true
	message := jutil.StringFormat(format, arguments...)
	e.self.ReportError(location, message)
}

// ReportWarningFormat formats the message with String.format semantics and
// reports it as a warning.
//
// Java: ErrorReporter.reportWarning(SourcePosition, String, Object...)
func (e *ErrorReporterBase) ReportWarningFormat(location LinePosition, format string, arguments ...any) {
	message := jutil.StringFormat(format, arguments...)
	e.self.ReportWarning(location, message)
}

// HadError reports whether an error was reported.
//
// Java: ErrorReporter.hadError()
func (e *ErrorReporterBase) HadError() bool {
	return e.hadError
}

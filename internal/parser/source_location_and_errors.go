// Port of com/google/javascript/jscomp/parsing/parser/util/SourcePosition.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/util/SourceRange.java @ closure-compiler v20261021
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

// Package parser is the port of Closure Compiler's
// com.google.javascript.jscomp.parsing.parser package together with its
// .trees and .util subpackages (one Go package because the three Java
// packages import each other; see PORTING.md §2).
//
// Source text is held as UTF-16 code units ([]uint16, like Java chars), so
// every offset, line and column is in UTF-16 units exactly as in Java.
package parser

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourcePosition is a position in a source string: its offset, in UTF-16
// code units. Its line and column come from the file's line table
// (SourceFile.LinePosition).
//
// Go-only: performance — Java's SourcePosition also holds the source file,
// the line and the column. Every token and parse tree has two positions,
// and most are only ever compared or turned into a node's line and column
// once, so they keep the offset alone (4 bytes instead of 12); the code
// that needs lines (comments, errors) uses LinePosition.
type SourcePosition struct {
	Offset int32
}

// GetOffset returns the offset.
//
// Java: SourcePosition.getOffset()
func (s SourcePosition) GetOffset() int {
	return int(s.Offset)
}

// NewSourcePosition returns the position at offset.
//
// Java: SourcePosition(SourceFile, int, int, int)
func NewSourcePosition(offset int) SourcePosition {
	return SourcePosition{Offset: int32(offset)}
}

// String renders the position as "@<offset>".
//
// Java: SourcePosition.toString() — Go-only: language-adapter — the line and
// column need the file (see LinePosition.String).
func (s SourcePosition) String() string {
	return "@" + strconv.Itoa(int(s.Offset))
}

// LinePosition is a position with its line and column (0-based), as Java's
// SourcePosition has them: what comments and errors keep.
//
// Go-only: see SourcePosition.
type LinePosition struct {
	Offset int32
	Line   int32
	Column int32
}

// Position returns the position without its line and column.
func (p LinePosition) Position() SourcePosition {
	return SourcePosition{Offset: p.Offset}
}

// GetOffset returns the offset.
func (p LinePosition) GetOffset() int {
	return int(p.Offset)
}

// String renders the position as "(<line>, <column>)", 1-based.
//
// Java: SourcePosition.toString() — Go-only: language-adapter — without the
// short source name.
func (p LinePosition) String() string {
	return "(" + strconv.Itoa(int(p.Line)+1) + ", " + strconv.Itoa(int(p.Column)+1) + ")"
}

// SourceRange is a range of positions in a source string.
// Start is inclusive.
// End is exclusive.
//
// Go-only: performance — a pointer-free value of two offsets (see
// SourcePosition).
type SourceRange struct {
	Start SourcePosition
	End   SourcePosition
}

// NewSourceRange constructs a SourceRange.
//
// Java: SourceRange(SourcePosition, SourcePosition)
func NewSourceRange(start, end SourcePosition) SourceRange {
	return SourceRange{Start: start, End: end}
}

// String renders the range as "<start - end>".
//
// Java: SourceRange.toString()
func (s SourceRange) String() string {
	return "<" + s.Start.String() + " - " + s.End.String() + ">"
}

// LineRange is a SourceRange whose positions have their lines and columns
// (see LinePosition): the location of a comment.
type LineRange struct {
	Start LinePosition
	End   LinePosition
}

// Range returns the range without the lines and columns.
func (r LineRange) Range() SourceRange {
	return SourceRange{Start: r.Start.Position(), End: r.End.Position()}
}

// String renders the range as "<start - end>".
func (r LineRange) String() string {
	return "<" + r.Start.String() + " - " + r.End.String() + ">"
}

// ErrorReporter is a conduit for reporting errors and warnings to the user.
//
// Java: abstract class ErrorReporter. Concrete Go implementations implement
// ReportError, ReportWarning, and HadError directly without self-dispatch.
type ErrorReporter interface {
	ReportError(location LinePosition, message string)
	ReportWarning(location LinePosition, message string)
	HadError() bool
}

// ErrorReporterBase provides minimal HadError tracking without a self
// ErrorReporter reference for external ErrorReporter implementations.
type ErrorReporterBase struct {
	hadError bool
}

// InitErrorReporterBase is a no-op retained for compatibility with external
// benchmark harnesses.
func (e *ErrorReporterBase) InitErrorReporterBase(_ ErrorReporter) {}

// HadError reports whether an error was reported.
//
// Java: ErrorReporter.hadError()
func (e *ErrorReporterBase) HadError() bool {
	return e.hadError
}

func (e *ErrorReporterBase) markError() {
	e.hadError = true
}

func reportErrorFormat(r ErrorReporter, location LinePosition, format string, arguments ...any) {
	if b, ok := r.(interface{ markError() }); ok {
		b.markError()
	}
	r.ReportError(location, jutil.StringFormat(format, arguments...))
}

func reportWarningFormat(r ErrorReporter, location LinePosition, format string, arguments ...any) {
	r.ReportWarning(location, jutil.StringFormat(format, arguments...))
}

// Port of com/google/javascript/jscomp/LightweightMessageFormatter.java @ closure-compiler v20260415

/*
 * Copyright 2007 The Closure Compiler Authors.
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

package jscomp

import (
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// LightweightMessageFormatter is a lightweight message formatter. The format
// of messages this formatter produces is very compact and to the point.
type LightweightMessageFormatter struct {
	AbstractMessageFormatterBase
	defaultFormat   SourceExcerptProvider_SourceExcerpt
	includeLocation bool
	includeLevel    bool
}

var _ AbstractMessageFormatter = (*LightweightMessageFormatter)(nil)

var lightweightMessageFormatter_excerptFormatter SourceExcerptProvider_ExcerptFormatter = NewLightweightMessageFormatter_LineNumberingFormatter()

const lightweightMessageFormatter_MAX_MULTILINE_ERROR_LENGTH = 4 // must be even.

// newLightweightMessageFormatterWithoutSource is the constructor for when
// the client doesn't care about source information.
//
// Java: LightweightMessageFormatter()
func newLightweightMessageFormatterWithoutSource() *LightweightMessageFormatter {
	f := &LightweightMessageFormatter{includeLocation: true, includeLevel: true}
	f.InitAbstractMessageFormatterBase(nil)
	f.defaultFormat = SourceExcerptProvider_SourceExcerpt_LINE
	return f
}

// Java: LightweightMessageFormatter(SourceExcerptProvider)
func NewLightweightMessageFormatter(source SourceExcerptProvider) *LightweightMessageFormatter {
	return NewLightweightMessageFormatterWithFormat(source, SourceExcerptProvider_SourceExcerpt_LINE)
}

// Java: LightweightMessageFormatter(SourceExcerptProvider, SourceExcerpt)
func NewLightweightMessageFormatterWithFormat(source SourceExcerptProvider, defaultFormat SourceExcerptProvider_SourceExcerpt) *LightweightMessageFormatter {
	f := &LightweightMessageFormatter{includeLocation: true, includeLevel: true}
	f.InitAbstractMessageFormatterBase(source)
	jutil.CheckNotNull(source)
	f.defaultFormat = defaultFormat
	return f
}

// Java: LightweightMessageFormatter.withoutSource()
func LightweightMessageFormatter_WithoutSource() *LightweightMessageFormatter {
	return newLightweightMessageFormatterWithoutSource()
}

// Java: LightweightMessageFormatter.setIncludeLocation(boolean)
func (f *LightweightMessageFormatter) SetIncludeLocation(includeLocation bool) *LightweightMessageFormatter {
	f.includeLocation = includeLocation
	return f
}

// Java: LightweightMessageFormatter.setIncludeLevel(boolean)
func (f *LightweightMessageFormatter) SetIncludeLevel(includeLevel bool) *LightweightMessageFormatter {
	f.includeLevel = includeLevel
	return f
}

// Java: LightweightMessageFormatter.formatError(JSError)
func (f *LightweightMessageFormatter) FormatError(error *JSError) string {
	return f.Format(error, false)
}

// Java: LightweightMessageFormatter.formatWarning(JSError)
func (f *LightweightMessageFormatter) FormatWarning(warning *JSError) string {
	return f.Format(warning, true)
}

// Java: LightweightMessageFormatter.format(JSError, boolean)
func (f *LightweightMessageFormatter) Format(error *JSError, warning bool) string {
	source := f.GetSource()
	sourceName := error.SourceName()
	lineNumber := error.GetLineNumber()
	charno := error.Charno()

	// Format the non-reverse-mapped position.
	var b strings.Builder
	var boldLine strings.Builder

	var mapping *sourcemap.Mapping_OriginalMapping
	if source != nil {
		mapping = source.GetSourceMapping(error.SourceName(), error.GetLineNumber(), error.Charno())
	}

	// Check if we can reverse-map the source.
	if f.includeLocation {
		if mapping != nil {
			lightweightMessageFormatter_appendPosition(&b, sourceName, lineNumber, charno)

			originalFile := mapping.GetOriginalFile()
			sourceName = &originalFile
			lineNumber = mapping.GetLineNumber()
			charno = mapping.GetColumnPosition()

			b.WriteString("\nOriginally at:\n")
		}

		lightweightMessageFormatter_appendPosition(&boldLine, sourceName, lineNumber, charno)
	}

	if f.includeLevel {
		level := CheckLevel_ERROR
		if warning {
			level = CheckLevel_WARNING
		}
		boldLine.WriteString(f.GetLevelName(level))
		boldLine.WriteString(" - [")
		boldLine.WriteString(error.Type().key)
		boldLine.WriteString("] ")
	}

	boldLine.WriteString(error.Description())

	b.WriteString(f.MaybeEmbolden(boldLine.String()))
	b.WriteByte('\n')

	// For reverse-mapped sources, fall back to a single line excerpt because the excerpt length
	// cannot be reliably mapped.
	format := f.defaultFormat
	if mapping != nil {
		format = SourceExcerptProvider_SourceExcerpt_LINE
	}
	sourceExcerptWithPosition :=
		f.GetExcerptWithPositionAndFormat(
			error, sourceName, lineNumber, charno, format)

	b.WriteString(sourceExcerptWithPosition)

	return b.String()
}

// Java: LightweightMessageFormatter.getExcerptWithPosition(JSError)
func (f *LightweightMessageFormatter) GetExcerptWithPosition(error *JSError) string {
	return f.GetExcerptWithPositionAndLength(
		error,
		error.SourceName(),
		error.GetLineNumber(),
		error.Charno(),
		error.Length(),
		f.defaultFormat)
}

// Java: LightweightMessageFormatter.getExcerptWithPosition(JSError, String, int, int, SourceExcerpt)
func (f *LightweightMessageFormatter) GetExcerptWithPositionAndFormat(
	error *JSError, sourceName *string, lineNumber int, charno int, format SourceExcerptProvider_SourceExcerpt) string {
	nodeLength := error.Length()
	return f.GetExcerptWithPositionAndLength(error, sourceName, lineNumber, charno, nodeLength, format)
}

// Java: LightweightMessageFormatter.getExcerptWithPosition(JSError, String, int, int, int, SourceExcerpt)
func (f *LightweightMessageFormatter) GetExcerptWithPositionAndLength(
	error *JSError,
	sourceName *string,
	lineNumber int,
	charno int,
	length int,
	format SourceExcerptProvider_SourceExcerpt) string {
	var b strings.Builder

	source := f.GetSource()
	cookedLength := -1
	if charno >= 0 && length >= 0 {
		cookedLength = charno + length
	}

	var sourceExcerpt *string
	if source != nil {
		sourceExcerpt = format.GetWithLength(source, sourceName, lineNumber, cookedLength, lightweightMessageFormatter_excerptFormatter)
	}

	if sourceExcerpt == nil {
		return ""
	}

	success := f.PadSourceExcerpt(error, *sourceExcerpt, charno, length, format, &b)
	if success {
		return b.String()
	}
	message :=
		"Source excerpt could not be formatted. This may indicate a bug in the compiler."
	return message
}

// Java: LightweightMessageFormatter.padSourceExcerpt(JSError, String, int, int, SourceExcerpt, StringBuilder)
func (f *LightweightMessageFormatter) PadSourceExcerpt(
	error *JSError,
	sourceExcerpt string,
	charno int,
	length int,
	format SourceExcerptProvider_SourceExcerpt,
	b *strings.Builder) bool {
	if format == SourceExcerptProvider_SourceExcerpt_FULL {
		if charno >= 0 {
			success := f.PadMultipleLines(charno, sourceExcerpt, b, error.Node())
			if !success {
				return false
			}
		} else {
			b.WriteString(sourceExcerpt)
			b.WriteByte('\n')
		}
	} else {
		b.WriteString(sourceExcerpt)
		b.WriteByte('\n')

		// charno == sourceExcerpt.length() means something is missing
		// at the end of the line
		if format == SourceExcerptProvider_SourceExcerpt_LINE && 0 <= charno && charno <= jutil.StringLength(sourceExcerpt) {
			success := f.PadLine(charno, sourceExcerpt, b, length, error.Node())
			if !success {
				return false
			}
		}
	}
	return true
}

// Java: LightweightMessageFormatter.appendPosition(StringBuilder, String, int, int)
func lightweightMessageFormatter_appendPosition(
	b *strings.Builder, sourceName *string, lineNumber int, charno int) {
	if sourceName != nil {
		b.WriteString(*sourceName)
		if lineNumber > 0 {
			b.WriteByte(':')
			b.WriteString(strconv.Itoa(lineNumber))
			if charno >= 0 {
				b.WriteByte(':')
				b.WriteString(strconv.Itoa(charno))
			}
		}
		b.WriteString(": ")
	}
}

// PadLine appends the given excerpt, attempting to add "^^^^" highlighting
// to the parts of the excerpt covered by the given error node.
//
// charno is the (positive) charno representing the index into the first
// line; sourceExcerpt is the original source, possibly multiple lines
// separated by '\n'; errLength is the length of the error node, or -1 if
// there is no error node; errorNode is the error node, or nil if there is
// no error node. It returns true if the source excerpt was successfully
// added to the builder, false if some error condition was encountered (and
// the builder's contents should be discarded).
//
// Java: LightweightMessageFormatter.padLine(int, String, StringBuilder, int, Node)
func (f *LightweightMessageFormatter) PadLine(
	charno int, sourceExcerpt string, b *strings.Builder, errLength int, errorNode *rhino.Node) bool {
	// This can occur when somehow code is requesting to "pad" a line that is shorter than the
	// starting charno where we will insert '^' characters.
	// This typically indicates a bug somewhere else in the compiler, and by using checkState we
	// will automatically report a user-"friendly" message that asks them to file a bug.
	units := jutil.StringToUTF16(sourceExcerpt)
	if charno > len(units) {
		return false
	}
	// Append leading whitespace
	for i := 0; i < charno; i++ {
		c := units[i]
		if rhino.TokenUtil_IsWhitespace(int(c)) {
			// Whitespace characters are never surrogates.
			b.WriteRune(rune(c))
		} else {
			b.WriteByte(' ')
		}
	}
	if errorNode == nil {
		b.WriteString("^")
	} else {
		length := max(1, min(errLength, len(units)-charno))
		for i := 0; i < length; i++ {
			b.WriteString("^")
		}
	}
	b.WriteString("\n")
	return true
}

// PadMultipleLines appends lines from the given excerpt, attempting to add
// "^^^^" highlighting to the parts of the excerpt covered by the given error
// node.
//
// startCharno is the (positive) charno representing the index into the
// first line; sourceExcerpt is the original source, possibly multiple lines
// separated by '\n'; b is the StringBuilder to append to; errorNode is the
// error node, or nil if there is no error node. This is used to determine
// what parts of the source excerpt to highlight. It returns true if the
// source excerpt was successfully added to the builder, false if some error
// condition was encountered (and the builder's contents should be
// discarded).
//
// Java: LightweightMessageFormatter.padMultipleLines(int, String, StringBuilder, Node)
func (f *LightweightMessageFormatter) PadMultipleLines(
	startCharno int, sourceExcerpt string, b *strings.Builder, errorNode *rhino.Node) bool {
	if errorNode == nil {
		b.WriteString(sourceExcerpt)
		b.WriteString("\n")
		charWithLineNumberOffset := startCharno + jutil.StringIndexOfChar(sourceExcerpt, '|') + 2
		return f.PadLine(charWithLineNumberOffset, sourceExcerpt, b, -1, nil)
	}
	lines := strings.Split(sourceExcerpt, "\n")
	requiresTruncation := len(lines) > lightweightMessageFormatter_MAX_MULTILINE_ERROR_LENGTH
	truncationStart := lightweightMessageFormatter_MAX_MULTILINE_ERROR_LENGTH / 2
	truncationEnd := len(lines) - lightweightMessageFormatter_MAX_MULTILINE_ERROR_LENGTH/2

	remainingLength := errorNode.GetLength()
	charno := startCharno
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if requiresTruncation && i == truncationStart {
			b.WriteString("...\n")
		}
		shouldPrintLine := !requiresTruncation || i < truncationStart || i >= truncationEnd
		// The LineNumberingFormatter below will append "  5| " to the start of each line, so subtract
		// that from the source line length in order to match the actual error length.
		charWithLineNumberOffset := charno + jutil.StringIndexOfChar(line, '|') + 2

		if shouldPrintLine {
			b.WriteString(line)
			b.WriteString("\n")
			success := f.PadLine(charWithLineNumberOffset, line, b, remainingLength, errorNode)
			if !success {
				return false
			}
		}

		// add 1 to represent the newline; subtract the offset of the "  5| ".
		remainingLength -= (jutil.StringLength(line) + 1 - charWithLineNumberOffset)
		charno = 0
	}
	return true
}

// LightweightMessageFormatter_LineNumberingFormatter formats a region by
// appending line numbers in front, e.g.
//
//	 9| if (foo) {
//	10|   alert('bar');
//	11| }
//
// and return line excerpt without any modification.
type LightweightMessageFormatter_LineNumberingFormatter struct{}

var _ SourceExcerptProvider_ExcerptFormatter = (*LightweightMessageFormatter_LineNumberingFormatter)(nil)

// Java: LightweightMessageFormatter.LineNumberingFormatter()
func NewLightweightMessageFormatter_LineNumberingFormatter() *LightweightMessageFormatter_LineNumberingFormatter {
	return &LightweightMessageFormatter_LineNumberingFormatter{}
}

// Java: LightweightMessageFormatter.LineNumberingFormatter.formatLine(String, int)
func (*LightweightMessageFormatter_LineNumberingFormatter) FormatLine(line *string, lineNumber int) *string {
	return line
}

// FormatRegion may return nil.
//
// Go-only: language-adapter — the line breaks are ASCII, so cutting the
// code at their byte offsets yields the same substrings as Java's UTF-16
// indices.
//
// Java: LightweightMessageFormatter.LineNumberingFormatter.formatRegion(Region)
func (*LightweightMessageFormatter_LineNumberingFormatter) FormatRegion(region Region) *string {
	if region == nil {
		return nil
	}
	code := region.GetSourceExcerpt()
	if code == "" {
		return nil
	}

	// max length of the number display
	numberLength := len(strconv.Itoa(region.GetEndingLineNumber()))

	// formatting
	var builder strings.Builder
	builder.Grow(len(code) * 2)
	start := 0
	end := lightweightMessageFormatter_indexOfNewline(code, start)
	lineNumber := region.GetBeginningLineNumber()
	for start >= 0 {
		// line extraction
		var line string
		if end < 0 {
			line = code[start:]
			if line == "" {
				s := builder.String()
				s = s[:len(s)-1]
				return &s
			}
		} else {
			line = code[start:end]
		}
		builder.WriteString("  ")

		// nice spaces for the line number
		spaces := numberLength - len(strconv.Itoa(lineNumber))
		if spaces < 0 {
			// StringBuilder.repeat(CharSequence, int)
			panic(jutil.NewIllegalArgumentExceptionMsg("count is negative: " + strconv.Itoa(spaces)))
		}
		builder.WriteString(strings.Repeat(" ", spaces))
		builder.WriteString(strconv.Itoa(lineNumber))
		builder.WriteString("| ")

		// end & update
		if end < 0 {
			builder.WriteString(line)
			start = -1
		} else {
			builder.WriteString(line)
			builder.WriteByte('\n')
			start = end + 1
			end = lightweightMessageFormatter_indexOfNewline(code, start)
			lineNumber++
		}
	}
	s := builder.String()
	return &s
}

// lightweightMessageFormatter_indexOfNewline is code.indexOf('\n', from)
// with byte offsets.
func lightweightMessageFormatter_indexOfNewline(code string, from int) int {
	if from >= len(code) {
		return -1
	}
	i := strings.IndexByte(code[from:], '\n')
	if i < 0 {
		return -1
	}
	return from + i
}

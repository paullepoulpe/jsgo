// Port of com/google/javascript/jscomp/LightweightMessageFormatter.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/AbstractMessageFormatter.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/ErrorFormat.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/VerboseMessageFormatter.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/SourceExcerptProvider.java @ closure-compiler v20261021

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

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// MessageFormatter formats warnings and errors. This interface may be used
// by implementations of ErrorManager to request message formatting
// capabilities.
type MessageFormatter interface {
	// FormatError formats an error.
	FormatError(error *JSError) string

	// FormatWarning formats a warning.
	FormatWarning(warning *JSError) string
}

// Java: AbstractMessageFormatter.termSupportsColor(String)
func AbstractMessageFormatter_TermSupportsColor(term string) bool {
	switch term {
	case "xterm", "xterm-color", "xterm-256color", "screen-bce":
		return true
	default:
		return false
	}
}

// Java: AbstractMessageFormatter.Color
type abstractMessageFormatter_Color int32

const (
	abstractMessageFormatter_Color_ERROR abstractMessageFormatter_Color = iota
	abstractMessageFormatter_Color_WARNING
	abstractMessageFormatter_Color_NO_COLOR
	abstractMessageFormatter_Color_BOLD
	abstractMessageFormatter_Color_UNBOLD
)

var abstractMessageFormatter_Color_controlCharacter = [...]string{
	"\u001b[31m",
	"\u001b[35m",
	"\u001b[39m",
	"\u001b[1m",
	"\u001b[0m",
}

func (c abstractMessageFormatter_Color) controlCharacter() string {
	return abstractMessageFormatter_Color_controlCharacter[c]
}

func messageFormatterGetLevelName(level CheckLevel, colorize bool) string {
	switch level {
	case CheckLevel_ERROR:
		return messageFormatterMaybeColorize("ERROR", abstractMessageFormatter_Color_ERROR, colorize)
	case CheckLevel_WARNING:
		return messageFormatterMaybeColorize("WARNING", abstractMessageFormatter_Color_WARNING, colorize)
	default:
		return level.String()
	}
}

func messageFormatterMaybeEmbolden(text string, colorize bool) string {
	if !colorize {
		return text
	}
	return abstractMessageFormatter_Color_BOLD.controlCharacter() +
		text + abstractMessageFormatter_Color_UNBOLD.controlCharacter()
}

func messageFormatterMaybeColorize(text string, color abstractMessageFormatter_Color, colorize bool) string {
	if !colorize {
		return text
	}
	return color.controlCharacter() +
		text + abstractMessageFormatter_Color_NO_COLOR.controlCharacter()
}

// ErrorFormat lists the error formats available.
type ErrorFormat int32

const (
	ErrorFormat_SINGLELINE ErrorFormat = iota
	ErrorFormat_FULL
	ErrorFormat_MULTILINE
	ErrorFormat_SOURCELESS
)

var errorFormat_names = [...]string{"SINGLELINE", "FULL", "MULTILINE", "SOURCELESS"}

// Java: ErrorFormat.values()
func ErrorFormat_values() []ErrorFormat {
	return []ErrorFormat{ErrorFormat_SINGLELINE, ErrorFormat_FULL, ErrorFormat_MULTILINE, ErrorFormat_SOURCELESS}
}

// String is Enum.name().
func (f ErrorFormat) String() string { return errorFormat_names[f] }

// Ordinal is Enum.ordinal().
func (f ErrorFormat) Ordinal() int { return int(f) }

// ToFormatter converts to a concrete formatter.
//
// Java: ErrorFormat.toFormatter(SourceExcerptProvider, boolean)
func (f ErrorFormat) ToFormatter(source SourceExcerptProvider, colorize bool) MessageFormatter {
	var formatter *LightweightMessageFormatter
	switch f {
	case ErrorFormat_SINGLELINE:
		formatter = NewLightweightMessageFormatter(source)
	case ErrorFormat_FULL:
		formatter = NewLightweightMessageFormatterWithFormat(source, SourceExcerptProvider_SourceExcerpt_FULL)
	case ErrorFormat_MULTILINE:
		formatter = NewLightweightMessageFormatterWithFormat(source, SourceExcerptProvider_SourceExcerpt_REGION)
	default: // SOURCELESS
		formatter = LightweightMessageFormatter_WithoutSource()
	}
	formatter.SetColorize(colorize)
	return formatter
}

// VerboseMessageFormatter is a verbose message formatter. This formatter
// generates very loud and long messages with multi-line source excerpts.
type VerboseMessageFormatter struct {
	source   SourceExcerptProvider
	colorize bool
}

var _ MessageFormatter = (*VerboseMessageFormatter)(nil)

// Java: VerboseMessageFormatter(SourceExcerptProvider)
func NewVerboseMessageFormatter(source SourceExcerptProvider) *VerboseMessageFormatter {
	return &VerboseMessageFormatter{source: source}
}

// Java: AbstractMessageFormatter.setColorize(boolean)
func (f *VerboseMessageFormatter) SetColorize(colorize bool) {
	f.colorize = colorize
}

// Java: AbstractMessageFormatter.getSource()
func (f *VerboseMessageFormatter) GetSource() SourceExcerptProvider {
	return f.source
}

// Java: VerboseMessageFormatter.formatError(JSError)
func (f *VerboseMessageFormatter) FormatError(error *JSError) string {
	return messageFormatterGetLevelName(CheckLevel_ERROR, f.colorize) + ": " + f.Format(error)
}

// Java: VerboseMessageFormatter.formatWarning(JSError)
func (f *VerboseMessageFormatter) FormatWarning(warning *JSError) string {
	return messageFormatterGetLevelName(CheckLevel_WARNING, f.colorize) + ": " + f.Format(warning)
}

// Java: VerboseMessageFormatter.format(JSError)
func (f *VerboseMessageFormatter) Format(message *JSError) string {
	description := message.Description()
	sourceName := message.SourceName()
	lineNumber := message.GetLineNumber()
	sourceRegion := f.source.GetSourceRegion(sourceName, lineNumber)

	source := "(unknown source)"
	if sourceName != nil && *sourceName != "" {
		source = *sourceName
	}
	line := "(unknown line)"
	if lineNumber >= 0 {
		line = strconv.Itoa(lineNumber)
	}
	excerpt := "."
	if sourceRegion != nil {
		excerpt = ":\n\n" + sourceRegion.GetSourceExcerpt()
	}
	return description + " at " + source + " line " + line + " " + excerpt
}

// LightweightMessageFormatter is a lightweight message formatter. The format
// of messages this formatter produces is very compact and to the point.
type LightweightMessageFormatter struct {
	source          SourceExcerptProvider
	colorize        bool
	defaultFormat   SourceExcerptProvider_SourceExcerpt
	includeLocation bool
	includeLevel    bool
}

var _ MessageFormatter = (*LightweightMessageFormatter)(nil)

var lightweightMessageFormatter_excerptFormatter SourceExcerptProvider_ExcerptFormatter = NewLightweightMessageFormatter_LineNumberingFormatter()

const lightweightMessageFormatter_MAX_MULTILINE_ERROR_LENGTH = 4 // must be even.

// newLightweightMessageFormatterWithoutSource is the constructor for when
// the client doesn't care about source information.
//
// Java: LightweightMessageFormatter()
func newLightweightMessageFormatterWithoutSource() *LightweightMessageFormatter {
	return &LightweightMessageFormatter{
		defaultFormat:   SourceExcerptProvider_SourceExcerpt_LINE,
		includeLocation: true,
		includeLevel:    true,
	}
}

// Java: LightweightMessageFormatter(SourceExcerptProvider)
func NewLightweightMessageFormatter(source SourceExcerptProvider) *LightweightMessageFormatter {
	return NewLightweightMessageFormatterWithFormat(source, SourceExcerptProvider_SourceExcerpt_LINE)
}

// Java: LightweightMessageFormatter(SourceExcerptProvider, SourceExcerpt)
func NewLightweightMessageFormatterWithFormat(source SourceExcerptProvider, defaultFormat SourceExcerptProvider_SourceExcerpt) *LightweightMessageFormatter {
	jutil.CheckNotNull(source)
	return &LightweightMessageFormatter{
		source:          source,
		defaultFormat:   defaultFormat,
		includeLocation: true,
		includeLevel:    true,
	}
}

// Java: AbstractMessageFormatter.setColorize(boolean)
func (f *LightweightMessageFormatter) SetColorize(colorize bool) {
	f.colorize = colorize
}

// Java: AbstractMessageFormatter.getSource()
func (f *LightweightMessageFormatter) GetSource() SourceExcerptProvider {
	return f.source
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

func (f *LightweightMessageFormatter) appendHeaderLine(
	b *strings.Builder,
	error *JSError,
	warning bool,
	sourceName string,
	hasSource bool,
	lineNumber int,
	charno int,
) {
	if f.includeLocation && hasSource {
		lightweightMessageFormatter_appendPosition(b, sourceName, lineNumber, charno)
	}
	if f.includeLevel {
		level := CheckLevel_ERROR
		if warning {
			level = CheckLevel_WARNING
		}
		b.WriteString(messageFormatterGetLevelName(level, f.colorize))
		b.WriteString(" - [")
		b.WriteString(error.Type().key)
		b.WriteString("] ")
	}
	b.WriteString(error.Description())
}

// Java: LightweightMessageFormatter.format(JSError, boolean)
func (f *LightweightMessageFormatter) Format(error *JSError, warning bool) string {
	source := f.source
	sourceNamePtr := error.SourceName()
	var sourceName string
	hasSource := sourceNamePtr != nil
	if hasSource {
		sourceName = *sourceNamePtr
	}
	lineNumber := error.GetLineNumber()
	charno := error.Charno()

	// Format the non-reverse-mapped position.
	var b strings.Builder
	b.Grow(256)

	var mapping *sourcemap.Mapping_OriginalMapping
	if source != nil {
		mapping = source.GetSourceMapping(sourceNamePtr, lineNumber, charno)
	}

	// Check if we can reverse-map the source.
	if f.includeLocation && mapping != nil {
		if hasSource {
			lightweightMessageFormatter_appendPosition(&b, sourceName, lineNumber, charno)
		}

		originalFile := mapping.GetOriginalFile()
		sourceName = originalFile
		hasSource = true
		sourceNamePtr = &originalFile
		lineNumber = mapping.GetLineNumber()
		charno = mapping.GetColumnPosition()

		b.WriteString("\nOriginally at:\n")
	}

	if !f.colorize {
		f.appendHeaderLine(&b, error, warning, sourceName, hasSource, lineNumber, charno)
	} else {
		var boldLine strings.Builder
		f.appendHeaderLine(&boldLine, error, warning, sourceName, hasSource, lineNumber, charno)
		b.WriteString(messageFormatterMaybeEmbolden(boldLine.String(), f.colorize))
	}
	b.WriteByte('\n')

	// For reverse-mapped sources, fall back to a single line excerpt because the excerpt length
	// cannot be reliably mapped.
	format := f.defaultFormat
	if mapping != nil {
		format = SourceExcerptProvider_SourceExcerpt_LINE
	}
	f.appendExcerptWithPositionAndFormat(
		&b, error, sourceNamePtr, lineNumber, charno, format)

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
	var b strings.Builder
	f.appendExcerptWithPositionAndFormat(&b, error, sourceName, lineNumber, charno, format)
	return b.String()
}

func (f *LightweightMessageFormatter) appendExcerptWithPositionAndFormat(
	b *strings.Builder,
	error *JSError,
	sourceName *string,
	lineNumber int,
	charno int,
	format SourceExcerptProvider_SourceExcerpt) {
	f.appendExcerptWithPositionAndLength(b, error, sourceName, lineNumber, charno, error.Length(), format)
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
	f.appendExcerptWithPositionAndLength(&b, error, sourceName, lineNumber, charno, length, format)
	return b.String()
}

func (f *LightweightMessageFormatter) appendExcerptWithPositionAndLength(
	b *strings.Builder,
	error *JSError,
	sourceName *string,
	lineNumber int,
	charno int,
	length int,
	format SourceExcerptProvider_SourceExcerpt) {
	source := f.source
	if source == nil {
		return
	}
	cookedLength := -1
	if charno >= 0 && length >= 0 {
		cookedLength = charno + length
	}

	sourceExcerpt := format.GetWithLength(source, sourceName, lineNumber, cookedLength, lightweightMessageFormatter_excerptFormatter)
	if sourceExcerpt == nil {
		return
	}

	success := f.PadSourceExcerpt(error, *sourceExcerpt, charno, length, format, b)
	if !success {
		b.WriteString("Source excerpt could not be formatted. This may indicate a bug in the compiler.")
	}
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
	b *strings.Builder, sourceName string, lineNumber int, charno int) {
	b.WriteString(sourceName)
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
	charno int, sourceExcerpt string, b *strings.Builder, errLength int, errorNode rhino.Node) bool {
	// Fast path for pure-ASCII source lines (avoids allocating a []uint16 slice
	// on the heap for every warning).
	isASCII := true
	for i := 0; i < len(sourceExcerpt); i++ {
		if sourceExcerpt[i] >= 0x80 {
			isASCII = false
			break
		}
	}
	if isASCII {
		if charno > len(sourceExcerpt) {
			return false
		}
		b.Grow(max(charno, 0) + 16)
		for i := 0; i < charno; i++ {
			c := sourceExcerpt[i]
			if rhino.TokenUtil_IsWhitespace(int(c)) {
				b.WriteByte(c)
			} else {
				b.WriteByte(' ')
			}
		}
		if errorNode.IsNil() {
			b.WriteByte('^')
		} else {
			length := max(1, min(errLength, len(sourceExcerpt)-charno))
			for i := 0; i < length; i++ {
				b.WriteByte('^')
			}
		}
		b.WriteByte('\n')
		return true
	}

	// This can occur when somehow code is requesting to "pad" a line that is shorter than the
	// starting charno where we will insert '^' characters.
	// This typically indicates a bug somewhere else in the compiler, and by using checkState we
	// will automatically report a user-"friendly" message that asks them to file a bug.
	unitsLen := jutil.StringLength(sourceExcerpt)
	if charno > unitsLen {
		return false
	}
	b.Grow(max(charno, 0) + 16)
	// Append leading whitespace
	utf16Pos := 0
	for _, r := range sourceExcerpt {
		if utf16Pos >= charno {
			break
		}
		units := 1
		if r >= 0x10000 {
			units = 2
		}
		if units == 1 && rhino.TokenUtil_IsWhitespace(int(r)) {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
			if units == 2 && utf16Pos+1 < charno {
				b.WriteByte(' ')
			}
		}
		utf16Pos += units
	}
	for utf16Pos < charno {
		b.WriteByte(' ')
		utf16Pos++
	}
	if errorNode.IsNil() {
		b.WriteByte('^')
	} else {
		length := max(1, min(errLength, unitsLen-charno))
		for i := 0; i < length; i++ {
			b.WriteByte('^')
		}
	}
	b.WriteByte('\n')
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
	startCharno int, sourceExcerpt string, b *strings.Builder, errorNode rhino.Node) bool {
	if errorNode.IsNil() {
		charWithLineNumberOffset := startCharno + jutil.StringIndexOfChar(sourceExcerpt, '|') + 2
		if charWithLineNumberOffset > jutil.StringLength(sourceExcerpt) {
			return false
		}
		b.WriteString(sourceExcerpt)
		b.WriteByte('\n')
		return f.PadLine(charWithLineNumberOffset, sourceExcerpt, b, -1, rhino.NilNode)
	}
	numLines := strings.Count(sourceExcerpt, "\n") + 1
	requiresTruncation := numLines > lightweightMessageFormatter_MAX_MULTILINE_ERROR_LENGTH
	truncationStart := lightweightMessageFormatter_MAX_MULTILINE_ERROR_LENGTH / 2
	truncationEnd := numLines - lightweightMessageFormatter_MAX_MULTILINE_ERROR_LENGTH/2

	remainingLength := errorNode.GetLength()
	charno := startCharno
	rem := sourceExcerpt
	for i := 0; i < numLines; i++ {
		var line string
		if nl := strings.IndexByte(rem, '\n'); nl >= 0 {
			line = rem[:nl]
			rem = rem[nl+1:]
		} else {
			line = rem
		}
		if requiresTruncation && i == truncationStart {
			b.WriteString("...\n")
		}
		shouldPrintLine := !requiresTruncation || i < truncationStart || i >= truncationEnd
		// The LineNumberingFormatter below will append "  5| " to the start of each line, so subtract
		// that from the source line length in order to match the actual error length.
		charWithLineNumberOffset := charno + jutil.StringIndexOfChar(line, '|') + 2
		lineLen := jutil.StringLength(line)

		if shouldPrintLine {
			if charWithLineNumberOffset > lineLen {
				return false
			}
			b.WriteString(line)
			b.WriteByte('\n')
			success := f.PadLine(charWithLineNumberOffset, line, b, remainingLength, errorNode)
			if !success {
				return false
			}
		}

		// add 1 to represent the newline; subtract the offset of the "  5| ".
		remainingLength -= (lineLen + 1 - charWithLineNumberOffset)
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
	numLines := strings.Count(code, "\n") + 1
	builder.Grow(len(code) + (numberLength+4)*numLines)
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
		lineNumberStr := strconv.Itoa(lineNumber)
		spaces := numberLength - len(lineNumberStr)
		if spaces < 0 {
			// StringBuilder.repeat(CharSequence, int)
			panic(jutil.NewIllegalArgumentExceptionMsg("count is negative: " + strconv.Itoa(spaces)))
		}
		for range spaces {
			builder.WriteByte(' ')
		}
		builder.WriteString(lineNumberStr)
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

// Region is a source code region.
type Region interface {
	// GetSourceExcerpt gets the source region.
	GetSourceExcerpt() string

	// GetBeginningLineNumber gets the beginning line number.
	GetBeginningLineNumber() int

	// GetEndingLineNumber gets the ending line number.
	GetEndingLineNumber() int
}

// SimpleRegion is a simple region.
type SimpleRegion struct {
	beginningLineNumber int
	endingLineNumber    int
	source              string
}

// Java: SimpleRegion(int, int, String)
func NewSimpleRegion(beginningLineNumber int, endingLineNumber int, source string) *SimpleRegion {
	return &SimpleRegion{beginningLineNumber: beginningLineNumber, endingLineNumber: endingLineNumber, source: source}
}

// Java: SimpleRegion.getBeginningLineNumber()
func (r *SimpleRegion) GetBeginningLineNumber() int { return r.beginningLineNumber }

// Java: SimpleRegion.getEndingLineNumber()
func (r *SimpleRegion) GetEndingLineNumber() int { return r.endingLineNumber }

// Java: SimpleRegion.getSourceExcerpt()
func (r *SimpleRegion) GetSourceExcerpt() string { return r.source }

// SourceExcerptProvider is responsible for building source code excerpt of
// specific locations, such as a specific line or a region around a given
// line number.
//
// Go-only: language-adapter — Java Strings that may be null (source names,
// lines, formatted excerpts) are *string; nil is null.
type SourceExcerptProvider interface {
	// GetSourceLine gets the line indicated by the line number (1 being the
	// first line of the file). This call will return only the specific line,
	// or nil if it does not exist.
	GetSourceLine(sourceName *string, lineNumber int) *string

	// GetSourceLines gets the specific lines returned by the beginning and
	// excerpt length. Must not start or end with a carriage return.
	// Implementations may decide to truncate lines but will always include
	// the first line and never be longer than the specified length.
	//
	// length is the desired length of the excerpt (in chars). If -1, returns
	// just a single line. Otherwise, returns as many lines as needed
	// (including the remainder of the last line). The result is nil if it
	// does not exist.
	GetSourceLines(sourceName *string, lineNumber int, length int) Region

	// GetSourceRegion gets a region around the indicated line number. The
	// exact definition of a region is implementation specific, but it must
	// contain the line indicated by the line number. A region must not start
	// or end by a carriage return. The result is nil if it does not exist.
	GetSourceRegion(sourceName *string, lineNumber int) Region

	// GetSourceMapping does a lookup for the original file position based on
	// this source map, if the source file is associated with an
	// inputSourceMap.
	GetSourceMapping(sourceName *string, lineNumber int, columnNumber int) *sourcemap.Mapping_OriginalMapping
}

// SourceExcerptProvider_SourceExcerpt is the source excerpt variety.
type SourceExcerptProvider_SourceExcerpt int32

const (
	// Line excerpt.
	SourceExcerptProvider_SourceExcerpt_LINE SourceExcerptProvider_SourceExcerpt = iota
	// Multiple lines excerpt.
	SourceExcerptProvider_SourceExcerpt_FULL
	// Region excerpt.
	SourceExcerptProvider_SourceExcerpt_REGION
)

var sourceExcerptProvider_SourceExcerpt_names = [...]string{"LINE", "FULL", "REGION"}

// String is Enum.name().
func (e SourceExcerptProvider_SourceExcerpt) String() string {
	return sourceExcerptProvider_SourceExcerpt_names[e]
}

// Get gets a source excerpt string based on the type of the source excerpt.
//
// Java: SourceExcerptProvider.SourceExcerpt.get(SourceExcerptProvider, String, int, ExcerptFormatter)
func (e SourceExcerptProvider_SourceExcerpt) Get(
	source SourceExcerptProvider, sourceName *string, lineNumber int, formatter SourceExcerptProvider_ExcerptFormatter) *string {
	return e.GetWithLength(source, sourceName, lineNumber, -1, formatter)
}

// GetWithLength gets a source excerpt string based on the type of the
// source excerpt.
//
// Java: SourceExcerptProvider.SourceExcerpt.get(SourceExcerptProvider, String, int, int, ExcerptFormatter)
func (e SourceExcerptProvider_SourceExcerpt) GetWithLength(
	source SourceExcerptProvider, sourceName *string, lineNumber int, length int,
	formatter SourceExcerptProvider_ExcerptFormatter) *string {
	switch e {
	case SourceExcerptProvider_SourceExcerpt_LINE:
		return formatter.FormatLine(source.GetSourceLine(sourceName, lineNumber), lineNumber)
	case SourceExcerptProvider_SourceExcerpt_FULL:
		return formatter.FormatRegion(source.GetSourceLines(sourceName, lineNumber, length))
	default: // REGION
		return formatter.FormatRegion(source.GetSourceRegion(sourceName, lineNumber))
	}
}

// SourceExcerptProvider_ExcerptFormatter is responsible of formatting
// source excerpts.
type SourceExcerptProvider_ExcerptFormatter interface {
	// FormatLine formats a line excerpt.
	FormatLine(line *string, lineNumber int) *string

	// FormatRegion formats a region excerpt.
	FormatRegion(region Region) *string
}

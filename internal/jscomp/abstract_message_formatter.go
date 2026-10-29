// Port of com/google/javascript/jscomp/AbstractMessageFormatter.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/ErrorFormat.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/VerboseMessageFormatter.java @ closure-compiler v20261021

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

import "strconv"

// MessageFormatter formats warnings and errors. This interface may be used
// by implementations of ErrorManager to request message formatting
// capabilities.
type MessageFormatter interface {
	// FormatError formats an error.
	FormatError(error *JSError) string

	// FormatWarning formats a warning.
	FormatWarning(warning *JSError) string
}

// AbstractMessageFormatter is an abstract message formatter providing
// default behavior for implementations of MessageFormatter needing a
// SourceExcerptProvider.
type AbstractMessageFormatter interface {
	MessageFormatter
	SetColorize(colorize bool)
	GetAbstractMessageFormatterBase() *AbstractMessageFormatterBase
}

// AbstractMessageFormatterBase holds the fields and base implementations of
// AbstractMessageFormatter; subclasses embed it.
type AbstractMessageFormatterBase struct {
	source   SourceExcerptProvider
	colorize bool
}

// InitAbstractMessageFormatterBase is the Go form of super(source).
//
// Java: AbstractMessageFormatter(SourceExcerptProvider)
func (f *AbstractMessageFormatterBase) InitAbstractMessageFormatterBase(source SourceExcerptProvider) {
	f.source = source
}

// GetAbstractMessageFormatterBase gives access to the base fields.
func (f *AbstractMessageFormatterBase) GetAbstractMessageFormatterBase() *AbstractMessageFormatterBase {
	return f
}

// Java: AbstractMessageFormatter.setColorize(boolean)
func (f *AbstractMessageFormatterBase) SetColorize(colorize bool) {
	f.colorize = colorize
}

// GetSource gets the source excerpt provider.
//
// Java: AbstractMessageFormatter.getSource()
func (f *AbstractMessageFormatterBase) GetSource() SourceExcerptProvider {
	return f.source
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

// Java: AbstractMessageFormatter.Color.getControlCharacter()
func (c abstractMessageFormatter_Color) GetControlCharacter() string {
	return abstractMessageFormatter_Color_controlCharacter[c]
}

// Java: AbstractMessageFormatter.getLevelName(CheckLevel)
func (f *AbstractMessageFormatterBase) GetLevelName(level CheckLevel) string {
	switch level {
	case CheckLevel_ERROR:
		return f.MaybeColorize("ERROR", abstractMessageFormatter_Color_ERROR)
	case CheckLevel_WARNING:
		return f.MaybeColorize("WARNING", abstractMessageFormatter_Color_WARNING)
	default:
		return level.String()
	}
}

// Java: AbstractMessageFormatter.maybeEmbolden(String)
func (f *AbstractMessageFormatterBase) MaybeEmbolden(text string) string {
	if !f.colorize {
		return text
	}
	return abstractMessageFormatter_Color_BOLD.GetControlCharacter() +
		text + abstractMessageFormatter_Color_UNBOLD.GetControlCharacter()
}

// Java: AbstractMessageFormatter.maybeColorize(String, Color)
func (f *AbstractMessageFormatterBase) MaybeColorize(text string, color abstractMessageFormatter_Color) string {
	if !f.colorize {
		return text
	}
	return color.GetControlCharacter() +
		text + abstractMessageFormatter_Color_NO_COLOR.GetControlCharacter()
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
	AbstractMessageFormatterBase
}

var _ AbstractMessageFormatter = (*VerboseMessageFormatter)(nil)

// Java: VerboseMessageFormatter(SourceExcerptProvider)
func NewVerboseMessageFormatter(source SourceExcerptProvider) *VerboseMessageFormatter {
	f := &VerboseMessageFormatter{}
	f.InitAbstractMessageFormatterBase(source)
	return f
}

// Java: VerboseMessageFormatter.formatError(JSError)
func (f *VerboseMessageFormatter) FormatError(error *JSError) string {
	return f.GetLevelName(CheckLevel_ERROR) + ": " + f.Format(error)
}

// Java: VerboseMessageFormatter.formatWarning(JSError)
func (f *VerboseMessageFormatter) FormatWarning(warning *JSError) string {
	return f.GetLevelName(CheckLevel_WARNING) + ": " + f.Format(warning)
}

// Java: VerboseMessageFormatter.format(JSError)
func (f *VerboseMessageFormatter) Format(message *JSError) string {
	description := message.Description()
	sourceName := message.SourceName()
	lineNumber := message.GetLineNumber()
	sourceRegion := f.GetSource().GetSourceRegion(sourceName, lineNumber)

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

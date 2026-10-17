// Port of com/google/javascript/jscomp/AbstractMessageFormatter.java @ closure-compiler v20261015

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

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

var abstractMessageFormatter_SUPPORTED_COLOR_TERMINALS = jutil.NewLinkedHashSet("xterm", "xterm-color", "xterm-256color", "screen-bce")

// Java: AbstractMessageFormatter.termSupportsColor(String)
func AbstractMessageFormatter_TermSupportsColor(term string) bool {
	return abstractMessageFormatter_SUPPORTED_COLOR_TERMINALS.Contains(term)
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

// Port of com/google/javascript/jscomp/ErrorFormat.java @ closure-compiler v20261015

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
	switch f {
	case ErrorFormat_SINGLELINE:
		formatter := NewLightweightMessageFormatter(
			source)
		formatter.SetColorize(colorize)
		return formatter
	case ErrorFormat_FULL:
		formatter :=
			NewLightweightMessageFormatterWithFormat(source, SourceExcerptProvider_SourceExcerpt_FULL)
		formatter.SetColorize(colorize)
		return formatter
	case ErrorFormat_MULTILINE:
		formatter :=
			NewLightweightMessageFormatterWithFormat(source, SourceExcerptProvider_SourceExcerpt_REGION)
		formatter.SetColorize(colorize)
		return formatter
	default: // SOURCELESS
		formatter :=
			LightweightMessageFormatter_WithoutSource()
		formatter.SetColorize(colorize)
		return formatter
	}
}

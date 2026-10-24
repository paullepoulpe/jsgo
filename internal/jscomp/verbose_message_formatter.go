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

package jscomp

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

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
	var lineSource *string
	if sourceRegion != nil {
		excerpt := sourceRegion.GetSourceExcerpt()
		lineSource = &excerpt
	}
	source := "(unknown source)"
	if sourceName != nil && *sourceName != "" {
		source = *sourceName
	}
	line := strconv.Itoa(lineNumber)
	if lineNumber < 0 {
		line = "(unknown line)"
	}
	excerpt := "."
	if lineSource != nil {
		excerpt = ":\n\n" + *lineSource
	}
	return jutil.StringFormat(
		"%s at %s line %s %s",
		description,
		source,
		line,
		excerpt)
}

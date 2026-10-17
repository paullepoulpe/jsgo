// Port of com/google/javascript/jscomp/JsonErrorReportGenerator.java @ closure-compiler v20260415

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsonErrorReportGenerator is an error report generator that prints error
// and warning data to the print stream as an array of JSON objects.
type JsonErrorReportGenerator struct {
	stream                *jutil.PrintStream
	sourceExcerptProvider SourceExcerptProvider
}

var _ SortingErrorManager_ErrorReportGenerator = (*JsonErrorReportGenerator)(nil)

var jsonErrorReportGenerator_excerptFormatter = NewLightweightMessageFormatter_LineNumberingFormatter()

// NewJsonErrorReportGenerator creates an error manager. stream is the
// stream on which the errors and warnings should be printed (this class
// does not close the stream); sourceExcerptProvider is used to retrieve the
// source context which generated the error.
//
// Java: JsonErrorReportGenerator(PrintStream, SourceExcerptProvider)
func NewJsonErrorReportGenerator(stream *jutil.PrintStream, sourceExcerptProvider SourceExcerptProvider) *JsonErrorReportGenerator {
	return &JsonErrorReportGenerator{stream: stream, sourceExcerptProvider: sourceExcerptProvider}
}

// Java: JsonErrorReportGenerator.generateReport(SortingErrorManager)
func (g *JsonErrorReportGenerator) GenerateReport(manager SortingErrorManager) {
	// new JsonWriter(new OutputStreamWriter(new ByteArrayOutputStream(), UTF_8)): the
	// characters are collected here and encoded when appended to the stream.
	var bufferedStream strings.Builder
	jutil.Try(func() {
		jsonWriter := jutil.NewJsonWriter(&bufferedStream)
		// try-with-resources: close on exit; an exception from close() while
		// another one propagates is suppressed.
		defer func() {
			if r := recover(); r != nil {
				func() {
					defer func() { _ = recover() }()
					jsonWriter.Close()
				}()
				panic(r)
			}
			jsonWriter.Close()
		}()
		jsonWriter.BeginArray()
		for _, message := range manager.GetSortedDiagnostics() {
			sourceName := message.Error.SourceName()
			lineNumber := message.Error.GetLineNumber()
			charno := message.Error.Charno()

			jsonWriter.BeginObject()
			level := "warning"
			if message.Level != nil && *message.Level == CheckLevel_ERROR {
				level = "error"
			}
			jsonWriter.Name("level").ValueString(&level)
			description := message.Error.Description()
			jsonWriter.Name("description").ValueString(&description)
			key := message.Error.Type().key
			jsonWriter.Name("key").ValueString(&key)
			if message.Error.Requirement() != nil {
				jsonWriter.Name("requirement").BeginObject()
				ruleId := message.Error.Requirement().GetRuleId()
				jsonWriter.Name("ruleId").ValueString(&ruleId)
				jsonWriter.Name("configFiles").BeginArray()
				for _, configFile := range message.Error.Requirement().GetConfigFileList() {
					jsonWriter.ValueString(&configFile)
				}
				jsonWriter.EndArray()
				jsonWriter.EndObject()
			}
			jsonWriter.Name("source").ValueString(sourceName)
			jsonWriter.Name("line").ValueLong(int64(lineNumber))
			jsonWriter.Name("column").ValueLong(int64(charno))
			node := message.Error.Node()
			regionLength := message.Error.Length()
			if node != nil && regionLength > 0 {
				jsonWriter.Name("length").ValueLong(int64(message.Error.Length()))
			}

			// extract source excerpt
			sourceExcerpt :=
				SourceExcerptProvider_SourceExcerpt_LINE.Get(g.sourceExcerptProvider, sourceName, lineNumber, jsonErrorReportGenerator_excerptFormatter)
			if sourceExcerpt != nil {
				var b strings.Builder
				b.WriteString(*sourceExcerpt)
				b.WriteString("\n")

				// padding equal to the excerpt and arrow at the end
				// charno == sourceExcerpt.length() means something is missing
				// at the end of the line
				units := jutil.StringToUTF16(*sourceExcerpt)
				if 0 <= charno && charno <= len(units) {
					for i := 0; i < charno; i++ {
						c := units[i]
						if rhino.TokenUtil_IsWhitespace(int(c)) {
							// Whitespace characters are never surrogates.
							b.WriteRune(rune(c))
						} else {
							b.WriteByte(' ')
						}
					}
					if node == nil {
						b.WriteString("^")
					} else {
						length := max(1, min(regionLength, len(units)-charno))
						for i := 0; i < length; i++ {
							b.WriteString("^")
						}
					}
				}
				context := b.String()
				jsonWriter.Name("context").ValueString(&context)
			}

			mapping :=
				g.sourceExcerptProvider.GetSourceMapping(
					sourceName, message.Error.GetLineNumber(), message.Error.Charno())
			if mapping != nil {
				jsonWriter.Name("originalLocation").BeginObject()
				originalFile := mapping.GetOriginalFile()
				jsonWriter.Name("source").ValueString(&originalFile)
				jsonWriter.Name("line").ValueLong(int64(mapping.GetLineNumber()))
				jsonWriter.Name("column").ValueLong(int64(mapping.GetColumnPosition()))
				jsonWriter.EndObject()
			}
			jsonWriter.EndObject()
		}

		var summaryBuilder strings.Builder
		if manager.GetTypedPercent() > 0.0 {
			summaryBuilder.WriteString(
				jutil.StringFormat(
					"%d error(s), %d warning(s), %s%% typed",
					manager.GetErrorCount(), manager.GetWarningCount(), jutil.Formatter_FormatDecimalFloat(manager.GetTypedPercent(), 1)))
		} else {
			summaryBuilder.WriteString(
				jutil.StringFormat(
					"%d error(s), %d warning(s)", manager.GetErrorCount(), manager.GetWarningCount()))
		}

		jsonWriter.BeginObject()
		info := "info"
		jsonWriter.Name("level").ValueString(&info)
		summary := summaryBuilder.String()
		jsonWriter.Name("description").ValueString(&summary)
		jsonWriter.EndObject()

		jsonWriter.EndArray()
		jsonWriter.Flush()
	}, func(e jutil.IIOException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
	g.stream.Append(bufferedStream.String())
}

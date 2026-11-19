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
// Port of com/google/javascript/jscomp/JsonErrorReportGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.JsonErrorReportGenerator
type JsonErrorReportGenerator struct {
	stream                *j2gort.Ext_java_io_PrintStream `gson:"stream"`
	sourceExcerptProvider SourceExcerptProvider           `gson:"sourceExcerptProvider"`
	idHash_               atomic.Int32
}

func (this_ *JsonErrorReportGenerator) BindSelf_(s any) {
	_ = s
}
func (this_ *JsonErrorReportGenerator) IsJavaNil_() bool { return this_ == nil }
func (*JsonErrorReportGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsonErrorReportGenerator", "JsonErrorReportGenerator"
}
func (this_ *JsonErrorReportGenerator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsonErrorReportGenerator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsonErrorReportGenerator", "com.google.javascript.jscomp.JsonErrorReportGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[*JsonErrorReportGenerator](v_)
	})
}
func (*JsonErrorReportGenerator) Implements_jscomp_SortingErrorManager_ErrorReportGenerator() {}

var static_JsonErrorReportGenerator_excerptFormatter *LightweightMessageFormatter_LineNumberingFormatter

func StaticGet_JsonErrorReportGenerator_ExcerptFormatter() *LightweightMessageFormatter_LineNumberingFormatter {
	clinit_JsonErrorReportGenerator()
	return static_JsonErrorReportGenerator_excerptFormatter
}
func StaticSet_JsonErrorReportGenerator_ExcerptFormatter(v *LightweightMessageFormatter_LineNumberingFormatter) {
	clinit_JsonErrorReportGenerator()
	static_JsonErrorReportGenerator_excerptFormatter = v
}

var clinitDone_JsonErrorReportGenerator uint32
var clinitState_JsonErrorReportGenerator j2gort.ClinitState

func clinit_JsonErrorReportGenerator() {
	if atomic.LoadUint32(&clinitDone_JsonErrorReportGenerator) != 0 {
		return
	}
	clinitSlow_JsonErrorReportGenerator()
}

func clinitSlow_JsonErrorReportGenerator() {
	if !j2gort.EnterClinit(&clinitState_JsonErrorReportGenerator, &clinitDone_JsonErrorReportGenerator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JsonErrorReportGenerator, &clinitDone_JsonErrorReportGenerator)
	static_JsonErrorReportGenerator_excerptFormatter = NewLightweightMessageFormatter_LineNumberingFormatter()
}

func (this_ *JsonErrorReportGenerator) initFields_JsonErrorReportGenerator() {
	_ = this_
}

// Java constructor: JsonErrorReportGenerator(java.io.PrintStream stream, com.google.javascript.jscomp.SourceExcerptProvider sourceExcerptProvider)
func NewJsonErrorReportGenerator(stream *j2gort.Ext_java_io_PrintStream, sourceExcerptProvider SourceExcerptProvider) *JsonErrorReportGenerator {
	clinit_JsonErrorReportGenerator()
	this_ := &JsonErrorReportGenerator{}
	InitJsonErrorReportGenerator(this_, stream, sourceExcerptProvider)
	return this_
}

func InitJsonErrorReportGenerator(this_ *JsonErrorReportGenerator, stream *j2gort.Ext_java_io_PrintStream, sourceExcerptProvider SourceExcerptProvider) {
	_ = stream
	_ = sourceExcerptProvider
	this_.initFields_JsonErrorReportGenerator()
	this_.stream = stream
	this_.sourceExcerptProvider = sourceExcerptProvider
}

// Java method: JsonErrorReportGenerator.generateReport(com.google.javascript.jscomp.SortingErrorManager manager)
func (this_ *JsonErrorReportGenerator) GenerateReport(manager SortingErrorManager) {
	_ = manager
	var bufferedStream *j2gort.Ext_java_io_ByteArrayOutputStream = j2gort.NewExt_java_io_ByteArrayOutputStream()
	_ = bufferedStream
	func() {
		defer func() {
			if rec_1 := recover(); rec_1 != nil {
				jex_2 := javalang.RecoverException(rec_1)
				_ = jex_2
				if j2gort.InstanceOf[javaio.IIOException](jex_2) {
					e := j2gort.AsIface[javaio.IIOException](jex_2)
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_1)
				}
			}
		}()
		var jsonWriter *j2gort.Ext_com_google_gson_stream_JsonWriter = j2gort.NewExt_com_google_gson_stream_JsonWriter(j2gort.NewExt_java_io_OutputStreamWriter_2_Ext_java_io_OutputStream_Ext_java_nio_charset_Charset(bufferedStream, j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8))
		_ = jsonWriter
		defer j2gort.CloseResource(jsonWriter)
		jsonWriter.BeginArray()
		{
			cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(manager.GetSortedDiagnostics()))
			for cur_3.HasNext() {
				rawElem_4 := cur_3.Next()
				message := j2gort.CastPtr[SortingErrorManager_ErrorWithLevel](rawElem_4)
				_ = message
				var sourceName string = message.error_.SourceName()
				_ = sourceName
				var lineNumber int32 = message.error_.GetLineNumber()
				_ = lineNumber
				var charno int32 = message.error_.Charno()
				_ = charno
				jsonWriter.BeginObject()
				jsonWriter.Name("level").Value_1_string(func() string {
					if message.level == StaticGet_CheckLevel_ERROR() {
						return "error"
					}
					return "warning"
				}())
				jsonWriter.Name("description").Value_1_string(message.error_.Description())
				jsonWriter.Name("key").Value_1_string(message.error_.Type().key)
				if message.error_.Requirement() != nil {
					jsonWriter.Name("requirement").BeginObject()
					jsonWriter.Name("ruleId").Value_1_string(message.error_.Requirement().GetRuleId())
					jsonWriter.Name("configFiles").BeginArray()
					{
						cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_protobuf_ProtocolStringList](message.error_.Requirement().GetConfigFileList())))
						for cur_5.HasNext() {
							configFile := cur_5.NextStr()
							_ = configFile
							jsonWriter.Value_1_string(configFile)
						}
					}
					jsonWriter.EndArray()
					jsonWriter.EndObject()
				}
				jsonWriter.Name("source").Value_1_string(sourceName)
				jsonWriter.Name("line").Value_1_int64(int64(lineNumber))
				jsonWriter.Name("column").Value_1_int64(int64(charno))
				var node pkg_rhino.Node = message.error_.Node()
				_ = node
				var regionLength int32 = message.error_.Length()
				_ = regionLength
				if (node != 0) && (regionLength > int32(0)) {
					jsonWriter.Name("length").Value_1_int64(int64(message.error_.Length()))
				}
				var sourceExcerpt string = StaticGet_SourceExcerptProvider_SourceExcerpt_LINE().Get_4(this_.sourceExcerptProvider, sourceName, lineNumber, ToIface_SourceExcerptProvider_ExcerptFormatter(StaticGet_JsonErrorReportGenerator_ExcerptFormatter()))
				_ = sourceExcerpt
				if sourceExcerpt != j2gort.NullString {
					var b j2gort.LocalStringBuilder
					b.InitStr(sourceExcerpt)

					b.AppendStr("\n")
					if (int32(0) <= charno) && (charno <= j2gort.JStrLen(sourceExcerpt)) {
						{
							var i int32 = int32(0)
							_ = i
							for ; i < charno; i = i + 1 {
								var c uint16 = j2gort.JStrCharAt(sourceExcerpt, i)
								_ = c
								if pkg_rhino.TokenUtil_IsWhitespace(int32(c)) {
									b.AppendChar(c)
								} else {
									b.AppendChar(uint16(32))
								}
							}
						}
						if node == 0 {
							b.AppendStr("^")
						} else {
							var length int32 = max(int32(1), min(regionLength, (j2gort.JStrLen(sourceExcerpt)-charno)))
							_ = length
							{
								var i_2 int32 = int32(0)
								_ = i_2
								for ; i_2 < length; i_2 = i_2 + 1 {
									b.AppendStr("^")
								}
							}
						}
					}
					jsonWriter.Name("context").Value_1_string(b.String())
				}
				var mapping *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping = this_.sourceExcerptProvider.GetSourceMapping(sourceName, message.error_.GetLineNumber(), message.error_.Charno())
				_ = mapping
				if mapping != nil {
					jsonWriter.Name("originalLocation").BeginObject()
					jsonWriter.Name("source").Value_1_string(mapping.GetOriginalFile())
					jsonWriter.Name("line").Value_1_int64(int64(mapping.GetLineNumber()))
					jsonWriter.Name("column").Value_1_int64(int64(mapping.GetColumnPosition()))
					jsonWriter.EndObject()
				}
				jsonWriter.EndObject()
			}
		}
		var summaryBuilder j2gort.LocalStringBuilder

		if manager.GetTypedPercent() > float64(0.0) {
			summaryBuilder.AppendStr(j2gort.Ext_java_lang_String_Format_2("%d error(s), %d warning(s), %.1f%% typed", []any{j2gort.ValueOfInt(manager.GetErrorCount()), j2gort.ValueOfInt(manager.GetWarningCount()), j2gort.ValueOfDouble(manager.GetTypedPercent())}))
		} else {
			summaryBuilder.AppendStr(j2gort.Ext_java_lang_String_Format_2("%d error(s), %d warning(s)", []any{j2gort.ValueOfInt(manager.GetErrorCount()), j2gort.ValueOfInt(manager.GetWarningCount())}))
		}
		jsonWriter.BeginObject()
		jsonWriter.Name("level").Value_1_string("info")
		jsonWriter.Name("description").Value_1_string(summaryBuilder.String())
		jsonWriter.EndObject()
		jsonWriter.EndArray()
		jsonWriter.Flush()
	}()
	this_.stream.Append_1_any(j2gort.BoxString(bufferedStream.ToString_1_Ext_java_nio_charset_Charset(j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)))
}

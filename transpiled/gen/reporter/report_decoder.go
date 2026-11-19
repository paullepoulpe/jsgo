/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/instrumentation/reporter/ReportDecoder.java (generated mechanically by j2go).

package reporter

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_sourcemap "github.com/paullepoulpe/jsgo/transpiled/gen/sourcemap"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_report_decoder_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_report_decoder_2 func(any) string
var lambdaSingleton_report_decoder_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_report_decoder_4 j2gort.Ext_java_util_function_Function
var lambdaSingleton_report_decoder_5 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_report_decoder_6 func(any) string

func init() {
	lambdaSingleton_report_decoder_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.PtrToAny(j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Type_ValueOf_1_string(j2gort.CastString(mrp0_)))
	})
	lambdaSingleton_report_decoder_2 = func(mrp0_ any) string {
		return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile](mrp0_).GetFileName()
	}
	lambdaSingleton_report_decoder_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_14 any) any {
		entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_14)
		_ = entry
		return j2gort.PtrToAny(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile](j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile_NewBuilder_0().SetFileName(j2gort.CastString(entry.GetKey())).AddAllInstrumentationPointsStats(j2gort.AsIface[j2gort.Ext_java_util_List](entry.GetValue())).Build()))
	})
	lambdaSingleton_report_decoder_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_15 any) any {
		pointStats := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats](lp_15)
		_ = pointStats
		return j2gort.BoxString(pointStats.GetPoint().GetFileName())
	})
	lambdaSingleton_report_decoder_5 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(lp_16 any) int32 {
		pointStats := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats](lp_16)
		_ = pointStats
		return pointStats.GetPoint().GetLineNumber()
	})
	lambdaSingleton_report_decoder_6 = func(lp_17 any) string {
		pointStats := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats](lp_17)
		_ = pointStats
		return pointStats.GetPoint().GetFileName()
	}
}

// Java class: com.google.javascript.jscomp.instrumentation.reporter.ReportDecoder
type ReportDecoder struct {
	idHash_ atomic.Int32
}

func (this_ *ReportDecoder) BindSelf_(s any) {
	_ = s
}
func (this_ *ReportDecoder) IsJavaNil_() bool { return this_ == nil }
func (*ReportDecoder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.reporter.ReportDecoder", "ReportDecoder"
}
func (this_ *ReportDecoder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReportDecoder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.reporter.ReportDecoder", "com.google.javascript.jscomp.instrumentation.reporter.ReportDecoder", func(v_ any) bool {
		return j2gort.InstanceOf[*ReportDecoder](v_)
	})
}

func (this_ *ReportDecoder) initFields_ReportDecoder() {
	_ = this_
}

// Java constructor: ReportDecoder()
func NewReportDecoder() *ReportDecoder {
	this_ := &ReportDecoder{}
	InitReportDecoder(this_)
	return this_
}

func InitReportDecoder(this_ *ReportDecoder) {
	this_.initFields_ReportDecoder()
}

// Java static method: ReportDecoder.parseMappingFromFile(java.lang.String mappingFilePath)
func ReportDecoder_ParseMappingFromFile(mappingFilePath string) j2gort.Ext_java_util_Map {
	_ = mappingFilePath
	var fileContent string = j2gort.Ext_com_google_common_io_CharStreams_ToString(j2gort.AsIface[j2gort.Ext_java_lang_Readable](j2gort.Ext_java_nio_file_Files_NewBufferedReader_2(j2gort.Ext_java_nio_file_Path_Of_2(mappingFilePath, []string{}), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)))
	_ = fileContent
	return ReportDecoder_ParseMapping(fileContent)
}

// Java static method: ReportDecoder.parseMapping(java.lang.String mappingFileContent)
func ReportDecoder_ParseMapping(mappingFileContent string) j2gort.Ext_java_util_Map {
	_ = mappingFileContent
	var mapping j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = mapping
	var linesOfFile j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(10)).OmitEmptyStrings(), mappingFileContent)
	_ = linesOfFile
	if !(linesOfFile.Size() >= int32(3)) {
		j2gort.CheckStateFmt(false, "Malformed report %s", linesOfFile)
	}
	var fileNamesLine string = javalang.StringTrim(j2gort.CheckNotNullString(j2gort.ListGetStr(linesOfFile, int32(0))))
	_ = fileNamesLine
	var functionNamesLine string = javalang.StringTrim(j2gort.CheckNotNullString(j2gort.ListGetStr(linesOfFile, int32(1))))
	_ = functionNamesLine
	var typesLine string = javalang.StringTrim(j2gort.CheckNotNullString(j2gort.ListGetStr(linesOfFile, int32(2))))
	_ = typesLine
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(fileNamesLine), "FileNames:")) {
		j2gort.CheckState(false)
	}
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(functionNamesLine), "FunctionNames:")) {
		j2gort.CheckState(false)
	}
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(typesLine), "Types:")) {
		j2gort.CheckState(false)
	}
	var fileNamesAsJsonArray string = j2gort.JStrSubFrom(fileNamesLine, (j2gort.JStrIndexOf(fileNamesLine, ":") + int32(1)))
	_ = fileNamesAsJsonArray
	var functionNamesAsJsonArray string = j2gort.JStrSubFrom(functionNamesLine, (j2gort.JStrIndexOf(functionNamesLine, ":") + int32(1)))
	_ = functionNamesAsJsonArray
	var typesAsJsonArray string = j2gort.JStrSubFrom(typesLine, (j2gort.JStrIndexOf(typesLine, ":") + int32(1)))
	_ = typesAsJsonArray
	var stringListType any = NewReportDecoder__Anon_122().GetType()
	_ = stringListType
	var gson *j2gort.Ext_com_google_gson_Gson = j2gort.NewExt_com_google_gson_GsonBuilder().DisableHtmlEscaping().Create()
	_ = gson
	var localFileNames j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](gson.FromJson_2_string_any(fileNamesAsJsonArray, stringListType))
	_ = localFileNames
	var functionNames j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](gson.FromJson_2_string_any(functionNamesAsJsonArray, stringListType))
	_ = functionNames
	var typesAsStringList j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](gson.FromJson_2_string_any(typesAsJsonArray, stringListType))
	_ = typesAsStringList
	var types j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](typesAsStringList.Stream().Map(lambdaSingleton_report_decoder_1).Collect_1(j2gort.Ext_java_util_stream_Collectors_ToList()))
	_ = types
	{
		var i int32 = int32(3)
		_ = i
		for ; i < linesOfFile.Size(); i = i + 1 {
			var lineItem string = j2gort.ListGetStr(linesOfFile, i)
			_ = lineItem
			var id string = j2gort.JStrSub(lineItem, int32(0), j2gort.JStrIndexOfCp(lineItem, int32(uint16(58))))
			_ = id
			var encodedDetails string = j2gort.JStrSubFrom(lineItem, (j2gort.JStrIndexOfCp(lineItem, int32(uint16(58))) + int32(1)))
			_ = encodedDetails
			var encodedDetailsAsCharIt *ReportDecoder_StringCharIterator = NewReportDecoder_StringCharIterator(encodedDetails)
			_ = encodedDetailsAsCharIt
			var temp *j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_NewBuilder_0().SetFileName(j2gort.ListGetStr(localFileNames, pkg_sourcemap.Base64VLQ_Decode(encodedDetailsAsCharIt))).SetFunctionName(j2gort.ListGetStr(functionNames, pkg_sourcemap.Base64VLQ_Decode(encodedDetailsAsCharIt))).SetType(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Type](types.Get(pkg_sourcemap.Base64VLQ_Decode(encodedDetailsAsCharIt)))).SetLineNumber(pkg_sourcemap.Base64VLQ_Decode(encodedDetailsAsCharIt)).SetColumnNumber(pkg_sourcemap.Base64VLQ_Decode(encodedDetailsAsCharIt)).Build())
			_ = temp
			mapping.PutIfAbsent(j2gort.BoxString(id), j2gort.PtrToAny(temp))
		}
	}
	return mapping
}

// Java static method: ReportDecoder.decodeReport(java.util.Map<java.lang.String,com.google.javascript.jscomp.instrumentation.reporter.proto.InstrumentationPoint> mapping, java.util.Map<java.lang.String,java.lang.Long> frequencies)
func ReportDecoder_DecodeReport(mapping j2gort.Ext_java_util_Map, frequencies j2gort.Ext_java_util_Map) *j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile {
	_ = mapping
	_ = frequencies
	var instrumentationPoints j2gort.Ext_java_util_stream_Stream = j2gort.AsIface[j2gort.Ext_java_util_Set](mapping.EntrySet()).Stream().Map(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
		entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_1)
		_ = entry
		return j2gort.PtrToAny(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats](j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_NewBuilder_0().SetPoint_1_Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](entry.GetValue())).SetTimesExecuted(j2gort.UnboxLong(j2gort.CastPtr[j2gort.BoxLong](j2gort.MapGetOrDefaultStr(frequencies, j2gort.CastString(entry.GetKey()), j2gort.ValueOfLong(int64(0)))))).Build()))
	}))
	_ = instrumentationPoints
	return ReportDecoder_CreateReportProfile(instrumentationPoints)
}

// Java static method: ReportDecoder.createProfileOfStaticallyUsedCode(java.util.Map<java.lang.String,com.google.javascript.jscomp.instrumentation.reporter.proto.InstrumentationPoint> mapping, java.lang.String fileContent, java.lang.String instrumentationArrayName)
func ReportDecoder_CreateProfileOfStaticallyUsedCode(mapping j2gort.Ext_java_util_Map, fileContent string, instrumentationArrayName string) *j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile {
	_ = mapping
	_ = fileContent
	_ = instrumentationArrayName
	var matcher *j2gort.Ext_java_util_regex_Matcher = j2gort.Ext_java_util_regex_Pattern_Compile_1((j2gort.StrValString(instrumentationArrayName) + "\\.push\\(['\"](.*?)['\"]\\)")).Matcher(j2gort.BoxString(fileContent))
	_ = matcher
	var encodedPoints j2gort.PlainSet[string] = j2gort.NewPlainSetVal[string]()
	_ = encodedPoints
	for matcher.Find_0() {
		encodedPoints.AddVoid(matcher.Group_1_int32(int32(1)))
	}
	var instrumentationPoints j2gort.Ext_java_util_stream_Stream = j2gort.AsIface[j2gort.Ext_java_util_Set](mapping.EntrySet()).Stream().Map(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_2 any) any {
		entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_2)
		_ = entry
		return j2gort.PtrToAny(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats](j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_NewBuilder_0().SetPoint_1_Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](entry.GetValue())).SetTimesExecuted(int64(int32(0))).SetPointPresence(func() *j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence {
			if encodedPoints.Contains(j2gort.CastString(entry.GetKey())) {
				return j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence_F_PRESENT
			}
			return j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence_F_STATICALLY_REMOVED
		}()).Build()))
	}))
	_ = instrumentationPoints
	return ReportDecoder_CreateReportProfile(instrumentationPoints)
}

// Java static method: ReportDecoder.mergeProfiles(java.util.List<com.google.javascript.jscomp.instrumentation.reporter.proto.ReportProfile> profiles)
func ReportDecoder_MergeProfiles(profiles j2gort.Ext_java_util_List) *j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile {
	_ = profiles
	var frequencies j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_HashMap()
	_ = frequencies
	var presences j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_HashMap()
	_ = presences
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(profiles))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			profile := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile](rawElem_4)
			_ = profile
			{
				cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(profile.GetFileProfileList()))
				for cur_5.HasNext() {
					rawElem_6 := cur_5.Next()
					fileProfile := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile](rawElem_6)
					_ = fileProfile
					{
						cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(fileProfile.GetInstrumentationPointsStatsList()))
						for cur_7.HasNext() {
							rawElem_8 := cur_7.Next()
							pointStats := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats](rawElem_8)
							_ = pointStats
							frequencies.Compute(j2gort.PtrToAny(pointStats.GetPoint()), j2gort.NewLambda_Ext_java_util_function_BiFunction(func(lp_9 any, lp_10 any) any {
								p := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](lp_9)
								_ = p
								existing := j2gort.CastPtr[j2gort.BoxLong](lp_10)
								_ = existing
								return j2gort.ValueOfLong(((func() int64 {
									if existing == nil {
										return int64(int32(0))
									}
									return j2gort.UnboxLong(existing)
								}()) + pointStats.GetTimesExecuted()))
							}))
							presences.Compute(j2gort.PtrToAny(pointStats.GetPoint()), j2gort.NewLambda_Ext_java_util_function_BiFunction(func(lp_11 any, lp_12 any) any {
								p := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](lp_11)
								_ = p
								existing := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence](lp_12)
								_ = existing
								var iifeArg_1 *j2gort.
									Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence
							iifeSwitch_1:
								switch {
								default:

									if ((pointStats.GetPointPresence() == j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence_F_PRESENT) || (existing == j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence_F_PRESENCE_UNKNOWN)) || (existing == nil) {
										iifeArg_1 = pointStats.GetPointPresence()
										break iifeSwitch_1
									}
									iifeArg_1 = existing
									break iifeSwitch_1
								}
								return j2gort.PtrToAny(iifeArg_1)
							}))
						}
					}
				}
			}
		}
	}
	var finalPoints j2gort.Ext_java_util_stream_Stream = j2gort.AsIface[j2gort.Ext_java_util_Set](frequencies.EntrySet()).Stream().Map(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_13 any) any {
		entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_13)
		_ = entry
		return j2gort.PtrToAny(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats](j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_NewBuilder_0().SetPoint_1_Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](entry.GetKey())).SetTimesExecuted(j2gort.UnboxLong(j2gort.CastPtr[j2gort.BoxLong](entry.GetValue()))).SetPointPresence(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence](presences.Get(j2gort.PtrToAny(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](entry.GetKey()))))).Build()))
	}))
	_ = finalPoints
	return ReportDecoder_CreateReportProfile(finalPoints)
}

// Java static method: ReportDecoder.createReportProfile(java.util.stream.Stream<com.google.javascript.jscomp.instrumentation.reporter.proto.InstrumentationPointStats> instrumentationPoints)
func ReportDecoder_CreateReportProfile(instrumentationPoints j2gort.Ext_java_util_stream_Stream) *j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile {
	_ = instrumentationPoints
	var fileProfiles j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.AsIface[j2gort.Ext_java_util_Map](instrumentationPoints.Sorted_1(j2gort.ComparatorComparingStr(lambdaSingleton_report_decoder_6).ThenComparingInt(lambdaSingleton_report_decoder_5)).Collect_1(j2gort.Ext_java_util_stream_Collectors_GroupingBy_1(lambdaSingleton_report_decoder_4))).EntrySet()).Stream().Map(lambdaSingleton_report_decoder_3).Sorted_1(j2gort.ComparatorComparingStr(lambdaSingleton_report_decoder_2)).Collect_1(j2gort.Ext_java_util_stream_Collectors_ToList()))
	_ = fileProfiles
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile](j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile_NewBuilder_0().AddAllFileProfile(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](fileProfiles)).Build())
}

// Java class: com.google.javascript.jscomp.instrumentation.reporter.ReportDecoder.StringCharIterator
type ReportDecoder_StringCharIterator struct {
	content string `gson:"content"`
	length  int32  `gson:"length"`
	current int32  `gson:"current"`
	idHash_ atomic.Int32
}

func (this_ *ReportDecoder_StringCharIterator) BindSelf_(s any) {
	_ = s
}
func (this_ *ReportDecoder_StringCharIterator) IsJavaNil_() bool { return this_ == nil }
func (*ReportDecoder_StringCharIterator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.reporter.ReportDecoder$StringCharIterator", "StringCharIterator"
}
func (this_ *ReportDecoder_StringCharIterator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReportDecoder_StringCharIterator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.reporter.ReportDecoder.StringCharIterator", "com.google.javascript.jscomp.instrumentation.reporter.ReportDecoder$StringCharIterator", func(v_ any) bool {
		return j2gort.InstanceOf[*ReportDecoder_StringCharIterator](v_)
	})
}
func (*ReportDecoder_StringCharIterator) Implements_sourcemap_Base64VLQ_CharIterator() {}

func (this_ *ReportDecoder_StringCharIterator) initFields_ReportDecoder_StringCharIterator() {
	_ = this_
	this_.content = j2gort.NullString

}

// Java constructor: StringCharIterator(java.lang.String content)
func NewReportDecoder_StringCharIterator(content string) *ReportDecoder_StringCharIterator {
	this_ := &ReportDecoder_StringCharIterator{}
	InitReportDecoder_StringCharIterator(this_, content)
	return this_
}

func InitReportDecoder_StringCharIterator(this_ *ReportDecoder_StringCharIterator, content string) {
	_ = content
	this_.initFields_ReportDecoder_StringCharIterator()
	this_.content = content
	this_.length = j2gort.JStrLen(content)
}

// Java method: StringCharIterator.next()
func (this_ *ReportDecoder_StringCharIterator) Next() uint16 {
	return j2gort.JStrCharAt(this_.content, func() int32 { old_18 := this_.current; nw_19 := old_18 + 1; this_.current = nw_19; return old_18 }())
}

// Java method: StringCharIterator.hasNext()
func (this_ *ReportDecoder_StringCharIterator) HasNext() bool {
	return (this_.current < this_.length)
}

// Java class:
type ReportDecoder__Anon_122 struct {
	j2gort.Ext_com_google_gson_reflect_TypeTokenBase
}

func (this_ *ReportDecoder__Anon_122) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_gson_reflect_TypeTokenBase.BindExtSelf_(s)
	this_.Ext_com_google_gson_reflect_TypeTokenBase.Core_ = "java.util.List<java.lang.String>"
}
func (this_ *ReportDecoder__Anon_122) IsJavaNil_() bool { return this_ == nil }
func (*ReportDecoder__Anon_122) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.reporter.ReportDecoder$1", ""
}
func (*ReportDecoder__Anon_122) Implements_j2gort_Ext_com_google_gson_reflect_TypeToken() {}

func (this_ *ReportDecoder__Anon_122) initFields_ReportDecoder__Anon_122() {
	_ = this_
}

func NewReportDecoder__Anon_122() *ReportDecoder__Anon_122 {
	this_ := &ReportDecoder__Anon_122{}
	this_.BindSelf_(this_)
	j2gort.InitExt_com_google_gson_reflect_TypeToken(&this_.Ext_com_google_gson_reflect_TypeTokenBase)
	this_.initFields_ReportDecoder__Anon_122()
	return this_
}

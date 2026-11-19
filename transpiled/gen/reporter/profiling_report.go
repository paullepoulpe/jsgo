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
// Port of com/google/javascript/jscomp/instrumentation/reporter/ProfilingReport.java (generated mechanically by j2go).

package reporter

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

var lambdaSingleton_profiling_report_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_profiling_report_2 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_profiling_report_3 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_profiling_report_4 j2gort.Ext_java_util_function_ToLongFunction
var lambdaSingleton_profiling_report_5 j2gort.Ext_java_util_function_Function
var lambdaSingleton_profiling_report_6 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_profiling_report_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_2 any) any {
		point := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](lp_2)
		_ = point
		return NewProfilingReport_InstrumentationPointMetrics()
	})
	lambdaSingleton_profiling_report_2 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(lp_11 any) int32 {
		point := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_11)
		_ = point
		return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](point.GetKey()).GetType().GetNumber()
	})
	lambdaSingleton_profiling_report_3 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(lp_12 any) int32 {
		point := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_12)
		_ = point
		return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](point.GetKey()).GetLineNumber()
	})
	lambdaSingleton_profiling_report_4 = j2gort.NewLambda_Ext_java_util_function_ToLongFunction(func(lp_13 any) int64 {
		point := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_13)
		_ = point
		return -(j2gort.CastPtr[ProfilingReport_InstrumentationPointMetrics](point.GetValue()).totalTimesExecuted)
	})
	lambdaSingleton_profiling_report_5 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_14 any) any {
		point := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_14)
		_ = point
		return j2gort.BoxString(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](point.GetKey()).GetFileName())
	})
	lambdaSingleton_profiling_report_6 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_17 any) any {
		row := j2gort.AsIface[j2gort.Ext_java_util_List](lp_17)
		_ = row
		return j2gort.BoxString(j2gort.Ext_java_lang_String_Join_2_any_Ext_java_lang_Iterable("\t", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](row)))
	})
}

// Java class: com.google.javascript.jscomp.instrumentation.reporter.ProfilingReport
type ProfilingReport struct {
	idHash_ atomic.Int32
}

func (this_ *ProfilingReport) BindSelf_(s any) {
	_ = s
}
func (this_ *ProfilingReport) IsJavaNil_() bool { return this_ == nil }
func (*ProfilingReport) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.reporter.ProfilingReport", "ProfilingReport"
}
func (this_ *ProfilingReport) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProfilingReport) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.reporter.ProfilingReport", "com.google.javascript.jscomp.instrumentation.reporter.ProfilingReport", func(v_ any) bool {
		return j2gort.InstanceOf[*ProfilingReport](v_)
	})
}

func (this_ *ProfilingReport) initFields_ProfilingReport() {
	_ = this_
}

// Java constructor: ProfilingReport()
func NewProfilingReport() *ProfilingReport {
	this_ := &ProfilingReport{}
	InitProfilingReport(this_)
	return this_
}

func InitProfilingReport(this_ *ProfilingReport) {
	this_.initFields_ProfilingReport()
}

// Java static method: ProfilingReport.createTabSeparatedProfilingReport(java.util.Map<java.lang.String,com.google.javascript.jscomp.instrumentation.reporter.proto.InstrumentationPoint> instrumentationMapping, com.google.common.collect.ImmutableList<java.util.Map<java.lang.String,java.lang.Long>> allInstrumentationReports)
func ProfilingReport_CreateTabSeparatedProfilingReport(instrumentationMapping j2gort.Ext_java_util_Map, allInstrumentationReports j2gort.Ext_com_google_common_collect_ImmutableList) string {
	_ = instrumentationMapping
	_ = allInstrumentationReports
	var profiles j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](allInstrumentationReports.Stream().Map(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
		report := j2gort.AsIface[j2gort.Ext_java_util_Map](lp_1)
		_ = report
		return j2gort.PtrToAny(ReportDecoder_DecodeReport(instrumentationMapping, report))
	})).Collect_1(j2gort.Ext_java_util_stream_Collectors_ToList()))
	_ = profiles
	var metrics j2gort.Ext_java_util_Map = j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.AsIface[j2gort.Ext_java_util_Collection](instrumentationMapping.Values()).Stream().Collect_1(j2gort.Ext_java_util_stream_Collectors_ToMap_2(j2gort.Ext_java_util_function_Function_Identity(), lambdaSingleton_profiling_report_1)))
	_ = metrics
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
							var metric *ProfilingReport_InstrumentationPointMetrics = j2gort.CastPtr[ProfilingReport_InstrumentationPointMetrics](metrics.Get(j2gort.PtrToAny(pointStats.GetPoint())))
							_ = metric
							r_9 := metric
							r_9.totalTimesExecuted = (r_9.totalTimesExecuted + pointStats.GetTimesExecuted())
							r_10 := metric
							r_10.numberOfReports = r_10.numberOfReports + 1
						}
					}
				}
			}
		}
	}
	var result []j2gort.Ext_java_util_Map_Entry = j2gort.CopyCollectionToSlice[j2gort.Ext_java_util_Map_Entry](j2gort.AsIface[j2gort.Ext_java_util_Set](metrics.EntrySet()))
	_ = result
	j2gort.SliceSortComparator(result, j2gort.Ext_java_util_Comparator_Comparing_1(lambdaSingleton_profiling_report_5).ThenComparingLong(lambdaSingleton_profiling_report_4).ThenComparingInt(lambdaSingleton_profiling_report_3).ThenComparingInt(lambdaSingleton_profiling_report_2))
	var output j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = output
	output.Add_1(j2gort.Ext_com_google_common_collect_ImmutableList_Of_6("File", "Function name", "Line", "Type", "Total executed", "Reports occurred"))
	for _, entry := range result {
		_ = entry
		output.Add_1(j2gort.Ext_com_google_common_collect_ImmutableList_Of_6(j2gort.BoxString(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](entry.GetKey()).GetFileName()), j2gort.BoxString(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](entry.GetKey()).GetFunctionName()), strconv.FormatInt(int64(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](entry.GetKey()).GetLineNumber()), 10), j2gort.BoxString(j2gort.EnumName(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint](entry.GetKey()).GetType())), strconv.FormatInt(j2gort.CastPtr[ProfilingReport_InstrumentationPointMetrics](entry.GetValue()).totalTimesExecuted, 10), strconv.FormatInt(int64(j2gort.CastPtr[ProfilingReport_InstrumentationPointMetrics](entry.GetValue()).numberOfReports), 10)))
	}

	return j2gort.CastString(output.Stream().Map(lambdaSingleton_profiling_report_6).Collect_1(j2gort.Ext_java_util_stream_Collectors_Joining_1("\n")))
}

// Java class: com.google.javascript.jscomp.instrumentation.reporter.ProfilingReport.InstrumentationPointMetrics
type ProfilingReport_InstrumentationPointMetrics struct {
	totalTimesExecuted int64 `gson:"totalTimesExecuted"`
	numberOfReports    int32 `gson:"numberOfReports"`
	idHash_            atomic.Int32
}

func (this_ *ProfilingReport_InstrumentationPointMetrics) BindSelf_(s any) {
	_ = s
}
func (this_ *ProfilingReport_InstrumentationPointMetrics) IsJavaNil_() bool { return this_ == nil }
func (*ProfilingReport_InstrumentationPointMetrics) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.reporter.ProfilingReport$InstrumentationPointMetrics", "InstrumentationPointMetrics"
}
func (this_ *ProfilingReport_InstrumentationPointMetrics) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProfilingReport_InstrumentationPointMetrics) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.reporter.ProfilingReport.InstrumentationPointMetrics", "com.google.javascript.jscomp.instrumentation.reporter.ProfilingReport$InstrumentationPointMetrics", func(v_ any) bool {
		return j2gort.InstanceOf[*ProfilingReport_InstrumentationPointMetrics](v_)
	})
}

func (this_ *ProfilingReport_InstrumentationPointMetrics) initFields_ProfilingReport_InstrumentationPointMetrics() {
	_ = this_

}

// Java constructor: InstrumentationPointMetrics()
func NewProfilingReport_InstrumentationPointMetrics() *ProfilingReport_InstrumentationPointMetrics {
	this_ := &ProfilingReport_InstrumentationPointMetrics{}
	InitProfilingReport_InstrumentationPointMetrics(this_)
	return this_
}

func InitProfilingReport_InstrumentationPointMetrics(this_ *ProfilingReport_InstrumentationPointMetrics) {
	this_.initFields_ProfilingReport_InstrumentationPointMetrics()
}

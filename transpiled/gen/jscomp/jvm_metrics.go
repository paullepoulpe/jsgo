/*
 * Copyright 2012 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/JvmMetrics.java (generated mechanically by j2go).

package jscomp

import (
	"strings"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.JvmMetrics
type JvmMetrics struct {
	idHash_ atomic.Int32
}

func (this_ *JvmMetrics) BindSelf_(s any) {
	_ = s
}
func (this_ *JvmMetrics) IsJavaNil_() bool { return this_ == nil }
func (*JvmMetrics) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JvmMetrics", "JvmMetrics"
}
func (this_ *JvmMetrics) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JvmMetrics) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JvmMetrics", "com.google.javascript.jscomp.JvmMetrics", func(v_ any) bool {
		return j2gort.InstanceOf[*JvmMetrics](v_)
	})
}

func StaticGet_JvmMetrics_TABULAR_COLON_POS() int32 { return int32(40) }
func StaticGet_JvmMetrics_ONE_KILO_BYTE() int64     { return int64(1024) }
func StaticGet_JvmMetrics_ONE_MEGA_BYTE() int64     { return int64(1048576) }
func StaticGet_JvmMetrics_ONE_GIGA_BYTE() int64     { return int64(1073741824) }
func (this_ *JvmMetrics) initFields_JvmMetrics() {
	_ = this_
}

// Java constructor: JvmMetrics()
func NewJvmMetrics() *JvmMetrics {
	this_ := &JvmMetrics{}
	InitJvmMetrics(this_)
	return this_
}

func InitJvmMetrics(this_ *JvmMetrics) {
	this_.initFields_JvmMetrics()
}

// Java static method: JvmMetrics.maybeWriteJvmMetrics(java.io.PrintStream out, java.lang.String options)
func JvmMetrics_MaybeWriteJvmMetrics(out *j2gort.Ext_java_io_PrintStream, options string) {
	_ = out
	_ = options
	if options == j2gort.NullString {
		return
	}
	var verboseMode bool = false
	_ = verboseMode
	var prettyMode bool = false
	_ = prettyMode
	var st *j2gort.Ext_java_util_StringTokenizer = j2gort.NewExt_java_util_StringTokenizer_2(options, ":")
	_ = st
	if st.HasMoreTokens() {
		var mode string = st.NextToken_0()
		_ = mode
		if strings.EqualFold(j2gort.CheckNotNullString(mode), "verbose") {
			verboseMode = true
		}
	}
	if st.HasMoreTokens() {
		var mode_2 string = st.NextToken_0()
		_ = mode_2
		if strings.EqualFold(j2gort.CheckNotNullString(mode_2), "pretty") {
			prettyMode = true
		}
	}
	if st.HasMoreTokens() {
		for st.HasMoreTokens() {
			var types string = st.NextToken_0()
			_ = types
			var typeSt *j2gort.Ext_java_util_StringTokenizer = j2gort.NewExt_java_util_StringTokenizer_2(types, ",")
			_ = typeSt
			for typeSt.HasMoreElements() {
				var type_ string = typeSt.NextToken_0()
				_ = type_
				JvmMetrics_WriteMetrics(out, type_, verboseMode, prettyMode)
			}
		}
	} else {
		JvmMetrics_WriteMetrics(out, "all", verboseMode, prettyMode)
	}
}

// Java static method: JvmMetrics.writeMetrics(java.io.PrintStream out, java.lang.String type, boolean verbose, boolean pretty)
func JvmMetrics_WriteMetrics(out *j2gort.Ext_java_io_PrintStream, type_ string, verbose bool, pretty bool) {
	_ = out
	_ = type_
	_ = verbose
	_ = pretty
	if (j2gort.CheckNotNullString(type_) == "gc") || strings.EqualFold(j2gort.CheckNotNullString(type_), "all") {
		JvmMetrics_WriteGarbageCollectionStats(out, verbose, pretty)
	}
	if (j2gort.CheckNotNullString(type_) == "mem") || strings.EqualFold(j2gort.CheckNotNullString(type_), "all") {
		JvmMetrics_WriteMemoryMetrics(out, verbose, pretty)
	}
	if (j2gort.CheckNotNullString(type_) == "jit") || strings.EqualFold(j2gort.CheckNotNullString(type_), "all") {
		JvmMetrics_WriteJitMetrics(out, verbose, pretty)
	}
}

// Java static method: JvmMetrics.writeJitMetrics(java.io.PrintStream out, boolean verbose, boolean pretty)
func JvmMetrics_WriteJitMetrics(out *j2gort.Ext_java_io_PrintStream, verbose bool, pretty bool) {
	_ = out
	_ = verbose
	_ = pretty
	var cBean j2gort.Ext_java_lang_management_CompilationMXBean = j2gort.Ext_java_lang_management_ManagementFactory_GetCompilationMXBean()
	_ = cBean
	var name string = j2gort.NullString
	_ = name
	if verbose {
		name = cBean.GetName()
	} else {
		name = "total"
	}
	if pretty {
		out.Println_1_string("\nJIT Stats")
		out.Printf_2("\t%s jit time: %d ms%n", []any{j2gort.BoxString(name), j2gort.ValueOfLong(cBean.GetTotalCompilationTime())})
	} else {
		out.Println_1_string(JvmMetrics_NormalizeTabularColonPos(j2gort.Ext_java_lang_String_Format_2("%s-jit-time-ms : %d", []any{j2gort.BoxString(JvmMetrics_NormalizeName(name)), j2gort.ValueOfLong(cBean.GetTotalCompilationTime())})))
	}
}

// Java static method: JvmMetrics.writeOverallMemoryUsage(java.io.PrintStream out, java.lang.management.MemoryUsage usage, java.lang.String prefix, boolean pretty)
func JvmMetrics_WriteOverallMemoryUsage(out *j2gort.Ext_java_io_PrintStream, usage *j2gort.Ext_java_lang_management_MemoryUsage, prefix string, pretty bool) {
	_ = out
	_ = usage
	_ = prefix
	_ = pretty
	if pretty {
		out.Format_2("\t%s\n", []any{j2gort.BoxString(prefix)})
		out.Format_2("\t\tavailable         : %s\n", []any{j2gort.BoxString(JvmMetrics_FormatBytes(usage.GetMax()))})
		out.Format_2("\t\tcurrent           : %s\n", []any{j2gort.BoxString(JvmMetrics_FormatBytes(usage.GetUsed()))})
	} else {
		prefix = JvmMetrics_NormalizeName(prefix)
		out.Println_1_string(JvmMetrics_NormalizeTabularColonPos(j2gort.Ext_java_lang_String_Format_2("%s-available-bytes : %d", []any{j2gort.BoxString(prefix), j2gort.ValueOfLong(usage.GetMax())})))
		out.Println_1_string(JvmMetrics_NormalizeTabularColonPos(j2gort.Ext_java_lang_String_Format_2("%s-current-bytes : %d", []any{j2gort.BoxString(prefix), j2gort.ValueOfLong(usage.GetUsed())})))
	}
}

// Java static method: JvmMetrics.writePoolMemoryUsage(java.io.PrintStream out, java.lang.management.MemoryUsage usage, java.lang.management.MemoryUsage peakUsage, java.lang.@org.jspecify.annotations.Nullable String prefix, boolean pretty)
func JvmMetrics_WritePoolMemoryUsage(out *j2gort.Ext_java_io_PrintStream, usage *j2gort.Ext_java_lang_management_MemoryUsage, peakUsage *j2gort.Ext_java_lang_management_MemoryUsage, prefix string, pretty bool) {
	_ = out
	_ = usage
	_ = peakUsage
	_ = prefix
	_ = pretty
	if pretty {
		out.Format_2("\t\tavailable         : %s\n", []any{j2gort.BoxString(JvmMetrics_FormatBytes(usage.GetMax()))})
		out.Format_2("\t\tpeak              : %s\n", []any{j2gort.BoxString(JvmMetrics_FormatBytes(peakUsage.GetUsed()))})
		out.Format_2("\t\tcurrent           : %s\n", []any{j2gort.BoxString(JvmMetrics_FormatBytes(usage.GetUsed()))})
	} else {
		out.Println_1_string(JvmMetrics_NormalizeTabularColonPos(j2gort.Ext_java_lang_String_Format_2("%s-available-bytes : %d", []any{j2gort.BoxString(prefix), j2gort.ValueOfLong(usage.GetMax())})))
		out.Println_1_string(JvmMetrics_NormalizeTabularColonPos(j2gort.Ext_java_lang_String_Format_2("%s-peak-bytes : %d", []any{j2gort.BoxString(prefix), j2gort.ValueOfLong(peakUsage.GetUsed())})))
		out.Println_1_string(JvmMetrics_NormalizeTabularColonPos(j2gort.Ext_java_lang_String_Format_2("%s-current-bytes : %d", []any{j2gort.BoxString(prefix), j2gort.ValueOfLong(usage.GetUsed())})))
	}
}

// Java static method: JvmMetrics.writeMemoryMetrics(java.io.PrintStream out, boolean verbose, boolean pretty)
func JvmMetrics_WriteMemoryMetrics(out *j2gort.Ext_java_io_PrintStream, verbose bool, pretty bool) {
	_ = out
	_ = verbose
	_ = pretty
	if pretty {
		out.Println_1_string("\nMemory usage")
	}
	if verbose {
		var overallMemBean j2gort.Ext_java_lang_management_MemoryMXBean = j2gort.Ext_java_lang_management_ManagementFactory_GetMemoryMXBean()
		_ = overallMemBean
		var usage *j2gort.Ext_java_lang_management_MemoryUsage = overallMemBean.GetHeapMemoryUsage()
		_ = usage
		JvmMetrics_WriteOverallMemoryUsage(out, usage, "Heap", pretty)
		usage = overallMemBean.GetNonHeapMemoryUsage()
		JvmMetrics_WriteOverallMemoryUsage(out, usage, "Non-heap", pretty)
	}
	if verbose {
		var mpBeans j2gort.Ext_java_util_List = j2gort.Ext_java_lang_management_ManagementFactory_GetMemoryPoolMXBeans()
		_ = mpBeans
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(mpBeans))
			for cur_1.HasNext() {
				rawElem_2 := cur_1.Next()
				mpBean := j2gort.AsIface[j2gort.Ext_java_lang_management_MemoryPoolMXBean](rawElem_2)
				_ = mpBean
				var currentUsage *j2gort.Ext_java_lang_management_MemoryUsage = mpBean.GetUsage()
				_ = currentUsage
				var peakUsage *j2gort.Ext_java_lang_management_MemoryUsage = mpBean.GetPeakUsage()
				_ = peakUsage
				if pretty {
					out.Println_1_string(("\tPool " + j2gort.StrValString(mpBean.GetName())))
					JvmMetrics_WritePoolMemoryUsage(out, currentUsage, peakUsage, j2gort.NullString, true)
				} else {
					JvmMetrics_WritePoolMemoryUsage(out, currentUsage, peakUsage, ("mem-pool-" + j2gort.StrValString(JvmMetrics_NormalizeName(mpBean.GetName()))), false)
				}
			}
		}
	} else {
		var available int64 = int64(int32(0))
		_ = available
		var current int64 = int64(int32(0))
		_ = current
		var peak int64 = int64(int32(0))
		_ = peak
		var mpBeans_2 j2gort.Ext_java_util_List = j2gort.Ext_java_lang_management_ManagementFactory_GetMemoryPoolMXBeans()
		_ = mpBeans_2
		{
			cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(mpBeans_2))
			for cur_3.HasNext() {
				rawElem_4 := cur_3.Next()
				mpBean_2 := j2gort.AsIface[j2gort.Ext_java_lang_management_MemoryPoolMXBean](rawElem_4)
				_ = mpBean_2
				var currentUsage_2 *j2gort.Ext_java_lang_management_MemoryUsage = mpBean_2.GetUsage()
				_ = currentUsage_2
				available = (available + currentUsage_2.GetMax())
				current = (current + currentUsage_2.GetUsed())
				var peakUsage_2 *j2gort.Ext_java_lang_management_MemoryUsage = mpBean_2.GetPeakUsage()
				_ = peakUsage_2
				peak = (peak + peakUsage_2.GetUsed())
			}
		}
		var summaryUsage *j2gort.Ext_java_lang_management_MemoryUsage = j2gort.NewExt_java_lang_management_MemoryUsage(int64(int32(0)), current, current, available)
		_ = summaryUsage
		var summaryPeakUsage *j2gort.Ext_java_lang_management_MemoryUsage = j2gort.NewExt_java_lang_management_MemoryUsage(int64(int32(0)), peak, peak, peak)
		_ = summaryPeakUsage
		if pretty {
			out.Format_2("\tAggregate of %d memory pools\n", []any{j2gort.ValueOfInt(mpBeans_2.Size())})
			JvmMetrics_WritePoolMemoryUsage(out, summaryUsage, summaryPeakUsage, j2gort.NullString, true)
		} else {
			JvmMetrics_WritePoolMemoryUsage(out, summaryUsage, summaryPeakUsage, "mem", false)
		}
	}
}

// Java static method: JvmMetrics.writeGarbageCollectionStats(java.io.PrintStream out, boolean verbose, boolean pretty)
func JvmMetrics_WriteGarbageCollectionStats(out *j2gort.Ext_java_io_PrintStream, verbose bool, pretty bool) {
	_ = out
	_ = verbose
	_ = pretty
	var gcBeans j2gort.Ext_java_util_List = j2gort.Ext_java_lang_management_ManagementFactory_GetGarbageCollectorMXBeans()
	_ = gcBeans
	if verbose {
		if pretty {
			out.Println_1_string("\nGarbage collection stats")
			{
				cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(gcBeans))
				for cur_5.HasNext() {
					rawElem_6 := cur_5.Next()
					gcBean := j2gort.AsIface[j2gort.Ext_java_lang_management_GarbageCollectorMXBean](rawElem_6)
					_ = gcBean
					out.Println_1_string(("\tCollector " + j2gort.StrValString(gcBean.GetName())))
					out.Format_2("\t\tcollection count   : %d\n", []any{j2gort.ValueOfLong(gcBean.GetCollectionCount())})
					out.Format_2("\t\tcollection time    : %d ms\n", []any{j2gort.ValueOfLong(gcBean.GetCollectionTime())})
				}
			}
		} else {
			{
				cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(gcBeans))
				for cur_7.HasNext() {
					rawElem_8 := cur_7.Next()
					gcBean_2 := j2gort.AsIface[j2gort.Ext_java_lang_management_GarbageCollectorMXBean](rawElem_8)
					_ = gcBean_2
					var name string = JvmMetrics_NormalizeName(gcBean_2.GetName())
					_ = name
					out.Println_1_string(JvmMetrics_NormalizeTabularColonPos(j2gort.Ext_java_lang_String_Format_2(("gc-" + j2gort.StrValString(name) + "-collection-count : %d"), []any{j2gort.ValueOfLong(gcBean_2.GetCollectionCount())})))
					out.Println_1_string(JvmMetrics_NormalizeTabularColonPos(j2gort.Ext_java_lang_String_Format_2(("gc-" + j2gort.StrValString(name) + "-collection-time-ms : %d"), []any{j2gort.ValueOfLong(gcBean_2.GetCollectionTime())})))
				}
			}
		}
	} else {
		var collectionCount int64 = int64(int32(0))
		_ = collectionCount
		var collectionTime int64 = int64(int32(0))
		_ = collectionTime
		var collectorCount int32 = gcBeans.Size()
		_ = collectorCount
		{
			cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(gcBeans))
			for cur_9.HasNext() {
				rawElem_10 := cur_9.Next()
				gcBean_3 := j2gort.AsIface[j2gort.Ext_java_lang_management_GarbageCollectorMXBean](rawElem_10)
				_ = gcBean_3
				collectionCount = (collectionCount + gcBean_3.GetCollectionCount())
				collectionTime = (collectionTime + gcBean_3.GetCollectionTime())
			}
		}
		if pretty {
			out.Println_1_string("\nGarbage collection stats")
			out.Format_2("\tAggregate of %d collectors\n", []any{j2gort.ValueOfInt(collectorCount)})
			out.Format_2("\t\tcollection count   : %d\n", []any{j2gort.ValueOfLong(collectionCount)})
			out.Format_2("\t\tcollection time    : %d ms\n", []any{j2gort.ValueOfLong(collectionTime)})
		} else {
			var name_2 string = JvmMetrics_NormalizeName("aggregate")
			_ = name_2
			out.Println_1_string(JvmMetrics_NormalizeTabularColonPos(j2gort.Ext_java_lang_String_Format_2("gc-%s-collection-count : %d", []any{j2gort.BoxString(name_2), j2gort.ValueOfLong(collectionCount)})))
			out.Println_1_string(JvmMetrics_NormalizeTabularColonPos(j2gort.Ext_java_lang_String_Format_2("gc-%s-collection-time-ms : %d", []any{j2gort.BoxString(name_2), j2gort.ValueOfLong(collectionTime)})))
		}
	}
}

// Java static method: JvmMetrics.normalizeName(java.lang.String name)
func JvmMetrics_NormalizeName(name string) string {
	_ = name
	return javalang.StringToLowerCaseASCII(j2gort.JStrReplaceChar(name, uint16(32), uint16(95)))
}

// Java static method: JvmMetrics.normalizeTabularColonPos(java.lang.String string)
func JvmMetrics_NormalizeTabularColonPos(string_ string) string {
	_ = string_
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_string(string_)
	_ = sb
	var index int32 = sb.IndexOf_1(":")
	_ = index
	for ; index < int32(40); index = index + 1 {
		sb.Insert_2_int32_uint16(index, uint16(32))
	}
	return sb.String()
}

// Java static method: JvmMetrics.formatBytes(long numBytes)
func JvmMetrics_FormatBytes(numBytes int64) string {
	_ = numBytes
	if numBytes < int64(1024) {
		return j2gort.Ext_java_lang_String_Format_2("%d B", []any{j2gort.ValueOfLong(numBytes)})
	} else if numBytes < int64(1048576) {
		return j2gort.Ext_java_lang_String_Format_2("%d KB", []any{j2gort.ValueOfLong((numBytes / int64(1024)))})
	} else if numBytes < int64(1073741824) {
		return j2gort.Ext_java_lang_String_Format_2("%d MB", []any{j2gort.ValueOfLong((numBytes / int64(1048576)))})
	} else {
		return j2gort.Ext_java_lang_String_Format_2("%d GB", []any{j2gort.ValueOfLong((numBytes / int64(1073741824)))})
	}
}

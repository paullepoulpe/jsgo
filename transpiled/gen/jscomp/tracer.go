/*
 * Copyright 2002 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Tracer.java (generated mechanically by j2go).

package jscomp

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.Tracer
type Tracer struct {
	extraTracingValues []int64                      `gson:"extraTracingValues"`
	type_              string                       `gson:"type"`
	comment            string                       `gson:"comment"`
	startThread        *j2gort.Ext_java_lang_Thread `gson:"startThread"`
	startTimeMs        int64                        `gson:"startTimeMs"`
	stopTimeMs         int64                        `gson:"stopTimeMs"`
	idHash_            atomic.Int32
}

func (this_ *Tracer) BindSelf_(s any) {
	_ = s
}
func (this_ *Tracer) IsJavaNil_() bool { return this_ == nil }
func (*Tracer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Tracer", "Tracer"
}
func (this_ *Tracer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Tracer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Tracer", "com.google.javascript.jscomp.Tracer", func(v_ any) bool {
		return j2gort.InstanceOf[*Tracer](v_)
	})
}

var staticSyncMu_Tracer sync.Mutex
var static_Tracer_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_Tracer_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_Tracer()
	return static_Tracer_logger
}
func StaticSet_Tracer_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_Tracer()
	static_Tracer_logger = v
}

var static_Tracer_defaultPrettyPrint bool

func StaticGet_Tracer_DefaultPrettyPrint() bool {
	clinit_Tracer()
	return static_Tracer_defaultPrettyPrint
}
func StaticSet_Tracer_DefaultPrettyPrint(v bool) {
	clinit_Tracer()
	static_Tracer_defaultPrettyPrint = v
}

var static_Tracer_extraTracingStatistics j2gort.Ext_java_util_List

func StaticGet_Tracer_ExtraTracingStatistics() j2gort.Ext_java_util_List {
	clinit_Tracer()
	return static_Tracer_extraTracingStatistics
}
func StaticSet_Tracer_ExtraTracingStatistics(v j2gort.Ext_java_util_List) {
	clinit_Tracer()
	static_Tracer_extraTracingStatistics = v
}
func StaticGet_Tracer_MAX_TRACE_SIZE() int32 { return int32(1000) }

var static_Tracer_clock Tracer_InternalClock

func StaticGet_Tracer_Clock() Tracer_InternalClock {
	clinit_Tracer()
	return static_Tracer_clock
}
func StaticSet_Tracer_Clock(v Tracer_InternalClock) {
	clinit_Tracer()
	static_Tracer_clock = v
}

var static_Tracer_typeToCountMap *Tracer_AtomicTracerStatMap

func StaticGet_Tracer_TypeToCountMap() *Tracer_AtomicTracerStatMap {
	clinit_Tracer()
	return static_Tracer_typeToCountMap
}
func StaticSet_Tracer_TypeToCountMap(v *Tracer_AtomicTracerStatMap) {
	clinit_Tracer()
	static_Tracer_typeToCountMap = v
}

var static_Tracer_typeToSilentMap *Tracer_AtomicTracerStatMap

func StaticGet_Tracer_TypeToSilentMap() *Tracer_AtomicTracerStatMap {
	clinit_Tracer()
	return static_Tracer_typeToSilentMap
}
func StaticSet_Tracer_TypeToSilentMap(v *Tracer_AtomicTracerStatMap) {
	clinit_Tracer()
	static_Tracer_typeToSilentMap = v
}

var static_Tracer_typeToTimeMap *Tracer_AtomicTracerStatMap

func StaticGet_Tracer_TypeToTimeMap() *Tracer_AtomicTracerStatMap {
	clinit_Tracer()
	return static_Tracer_typeToTimeMap
}
func StaticSet_Tracer_TypeToTimeMap(v *Tracer_AtomicTracerStatMap) {
	clinit_Tracer()
	static_Tracer_typeToTimeMap = v
}

var static_Tracer_ZERO_STAT *Tracer_Stat

func StaticGet_Tracer_ZERO_STAT() *Tracer_Stat {
	clinit_Tracer()
	return static_Tracer_ZERO_STAT
}
func StaticSet_Tracer_ZERO_STAT(v *Tracer_Stat) {
	clinit_Tracer()
	static_Tracer_ZERO_STAT = v
}

var static_Tracer_traces *j2gort.Ext_java_lang_ThreadLocal

func StaticGet_Tracer_Traces() *j2gort.Ext_java_lang_ThreadLocal {
	clinit_Tracer()
	return static_Tracer_traces
}
func StaticSet_Tracer_Traces(v *j2gort.Ext_java_lang_ThreadLocal) {
	clinit_Tracer()
	static_Tracer_traces = v
}

var clinitDone_Tracer uint32
var clinitState_Tracer j2gort.ClinitState

func clinit_Tracer() {
	if atomic.LoadUint32(&clinitDone_Tracer) != 0 {
		return
	}
	clinitSlow_Tracer()
}

func clinitSlow_Tracer() {
	if !j2gort.EnterClinit(&clinitState_Tracer, &clinitDone_Tracer) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Tracer, &clinitDone_Tracer)
	static_Tracer_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.Tracer", "Tracer")).GetName())
	static_Tracer_extraTracingStatistics = j2gort.NewExt_java_util_concurrent_CopyOnWriteArrayList()
	static_Tracer_clock = NewTracer__Anon_98()
	static_Tracer_ZERO_STAT = NewTracer_Stat()
	static_Tracer_traces = j2gort.NewExt_java_lang_ThreadLocal()
}

func (this_ *Tracer) initFields_Tracer() {
	_ = this_
	this_.type_ = j2gort.NullString
	this_.comment = j2gort.NullString
}

// Java constructor: Tracer(java.lang.@org.jspecify.annotations.Nullable String type, java.lang.@org.jspecify.annotations.Nullable String comment)
func NewTracer_2(type_ string, comment string) *Tracer {
	clinit_Tracer()
	this_ := &Tracer{}
	InitTracer_2(this_, type_, comment)
	return this_
}

func InitTracer_2(this_ *Tracer, type_ string, comment string) {
	_ = type_
	_ = comment
	this_.initFields_Tracer()
	this_.type_ = type_
	this_.comment = j2gort.Ext_com_google_common_base_Strings_NullToEmpty(comment)
	this_.startTimeMs = StaticGet_Tracer_Clock().CurrentTimeMillis()
	this_.startThread = j2gort.Ext_java_lang_Thread_CurrentThread()
	if !(StaticGet_Tracer_ExtraTracingStatistics().IsEmpty()) {
		var size int32 = StaticGet_Tracer_ExtraTracingStatistics().Size()
		_ = size
		this_.extraTracingValues = make([]int64, size)
		var i int32 = int32(0)
		_ = i
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_Tracer_ExtraTracingStatistics()))
			for cur_1.HasNext() {
				rawElem_2 := cur_1.Next()
				tracingStatistic := j2gort.AsIface[Tracer_TracingStatistic](rawElem_2)
				_ = tracingStatistic
				{
					a_3 := this_.extraTracingValues
					i_4 := i
					a_3[i_4] = tracingStatistic.Start(this_.startThread)
				}
				i = i + 1
			}
		}
	}
	var trace *Tracer_ThreadTrace = Tracer_GetThreadTrace()
	_ = trace
	if !(trace.IsInitialized()) {
		return
	}
	if int32(len(trace.events)) >= int32(1000) {
		StaticGet_Tracer_Logger().Log_3_Ext_java_util_logging_Level_string_Throwable(j2gort.Ext_java_util_logging_Level_F_WARNING, "Giant thread trace. Too many Tracers created. Clearing to avoid memory leak.", j2gort.ConstructShimThrowable("java.lang.Throwable", j2gort.BoxString(trace.String())))
		trace.TruncateEvents()
	}
	if trace.outstandingEvents.Size() >= int32(1000) {
		StaticGet_Tracer_Logger().Log_3_Ext_java_util_logging_Level_string_Throwable(j2gort.Ext_java_util_logging_Level_F_WARNING, "Too many outstanding Tracers. Tracer.stop() is missing or Tracer.stop() is not wrapped in a try/finally block. Clearing to avoid memory leak.", j2gort.ConstructShimThrowable("java.lang.Throwable", j2gort.BoxString(trace.String())))
		trace.TruncateOutstandingEvents()
	}
	trace.StartEvent(this_)
}

// Java constructor: Tracer(java.lang.String comment)
func NewTracer_1(comment string) *Tracer {
	clinit_Tracer()
	this_ := &Tracer{}
	InitTracer_1(this_, comment)
	return this_
}

func InitTracer_1(this_ *Tracer, comment string) {
	_ = comment
	InitTracer_2(this_, j2gort.NullString, comment)
}

// Java static method: Tracer.longToPaddedString(long v, int digitsColumnWidth)
func Tracer_LongToPaddedString(v int64, digitsColumnWidth int32) string {
	clinit_Tracer()
	_ = v
	_ = digitsColumnWidth
	var digitWidth int32 = Tracer_NumDigits(v)
	_ = digitWidth
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	Tracer_AppendSpaces(sb, (digitsColumnWidth - digitWidth))
	sb.Append_1_int64(v)
	return sb.String()
}

// Java static method: Tracer.numDigits(long n)
func Tracer_NumDigits(n int64) int32 {
	clinit_Tracer()
	_ = n
	var i int32 = int32(0)
	_ = i
	for doFirst_5 := true; doFirst_5 || (n > int64(int32(0))); doFirst_5 = false {
		i = i + 1
		n = (n / int64(int32(10)))
	}
	return i
}

// Java static method: Tracer.appendSpaces(java.lang.StringBuilder sb, int numSpaces)
func Tracer_AppendSpaces(sb *j2gort.Ext_java_lang_StringBuilder, numSpaces int32) {
	clinit_Tracer()
	_ = sb
	_ = numSpaces
	if numSpaces > int32(16) {
		StaticGet_Tracer_Logger().Warning_1_string("Tracer.appendSpaces called with large numSpaces")
		numSpaces = int32(16)
	}
	for numSpaces >= int32(5) {
		sb.Append_1_string("     ")
		numSpaces = (numSpaces - int32(5))
	}
	switch numSpaces {
	case int32(1):
		sb.Append_1_string(" ")
	case int32(2):
		sb.Append_1_string("  ")
	case int32(3):
		sb.Append_1_string("   ")
	case int32(4):
		sb.Append_1_string("    ")
	}
}

// Java static method: Tracer.addTracingStatistic(com.google.javascript.jscomp.Tracer.TracingStatistic tracingStatistic)
func Tracer_AddTracingStatistic(tracingStatistic Tracer_TracingStatistic) int32 {
	clinit_Tracer()
	_ = tracingStatistic
	if tracingStatistic.Enable() {
		StaticGet_Tracer_ExtraTracingStatistics().Add_1(tracingStatistic)
		return StaticGet_Tracer_ExtraTracingStatistics().LastIndexOf_1(tracingStatistic)
	} else {
		return int32(-1)
	}
}

// Java static method: Tracer.clearTracingStatisticsTestingOnly()
func Tracer_ClearTracingStatisticsTestingOnly() {
	clinit_Tracer()
	StaticGet_Tracer_ExtraTracingStatistics().Clear()
}

// Java method: Tracer.stop(int silenceThreshold)
func (this_ *Tracer) Stop_1(silenceThreshold int32) int64 {
	_ = silenceThreshold
	if !(j2gort.Ext_java_lang_Thread_CurrentThread() == this_.startThread) {
		j2gort.CheckState(false)
	}
	var trace *Tracer_ThreadTrace = Tracer_GetThreadTrace()
	_ = trace
	if !(trace.IsInitialized()) {
		return int64(int32(0))
	}
	this_.stopTimeMs = StaticGet_Tracer_Clock().CurrentTimeMillis()
	if this_.extraTracingValues != nil {
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(this_.extraTracingValues)); i = i + 1 {
				var value int64 = j2gort.AsIface[Tracer_TracingStatistic](StaticGet_Tracer_ExtraTracingStatistics().Get(i)).Stop(this_.startThread)
				_ = value
				{
					a_6 := this_.extraTracingValues
					i_7 := i
					a_6[i_7] = (value - this_.extraTracingValues[i])
				}
			}
		}
	}
	if !(trace.IsInitialized()) {
		return int64(int32(0))
	}
	trace.EndEvent(this_, silenceThreshold)
	return (this_.stopTimeMs - this_.startTimeMs)
}

// Java method: Tracer.stop()
func (this_ *Tracer) Stop_0() int64 {
	return this_.Stop_1(int32(-1))
}

// Java method: Tracer.toString()
func (this_ *Tracer) String() string {
	if this_.type_ == j2gort.NullString {
		return this_.comment
	} else {
		return ("[" + j2gort.StrValString(this_.type_) + "] " + j2gort.StrValString(this_.comment))
	}
}

// Java static method: Tracer.setDefaultSilenceThreshold(int threshold)
func Tracer_SetDefaultSilenceThreshold(threshold int32) {
	clinit_Tracer()
	_ = threshold
	{
		r_8 := Tracer_GetThreadTrace()
		r_8.defaultSilenceThreshold = threshold
	}
}

// Java static method: Tracer.initCurrentThreadTrace()
func Tracer_InitCurrentThreadTrace_0() {
	clinit_Tracer()
	var events *Tracer_ThreadTrace = Tracer_GetThreadTrace()
	_ = events
	if !(events.IsEmpty()) {
		StaticGet_Tracer_Logger().Log_3_Ext_java_util_logging_Level_string_Throwable(j2gort.Ext_java_util_logging_Level_F_WARNING, ("Non-empty timer log:\n" + j2gort.StrVal(events)), j2gort.ConstructShimThrowable("java.lang.Throwable"))
		Tracer_ClearThreadTrace()
		events = Tracer_GetThreadTrace()
	}
	events.Init()
}

// Java static method: Tracer.initCurrentThreadTrace(int defaultSilenceThreshold)
func Tracer_InitCurrentThreadTrace_1(defaultSilenceThreshold int32) {
	clinit_Tracer()
	_ = defaultSilenceThreshold
	Tracer_InitCurrentThreadTrace_0()
	Tracer_SetDefaultSilenceThreshold(defaultSilenceThreshold)
}

// Java static method: Tracer.getCurrentThreadTraceReport()
func Tracer_GetCurrentThreadTraceReport() string {
	clinit_Tracer()
	return Tracer_GetThreadTrace().String()
}

// Java static method: Tracer.logCurrentThreadTrace()
func Tracer_LogCurrentThreadTrace() {
	clinit_Tracer()
	var trace *Tracer_ThreadTrace = Tracer_GetThreadTrace()
	_ = trace
	if !(trace.IsInitialized()) {
		StaticGet_Tracer_Logger().Log_3_Ext_java_util_logging_Level_string_Throwable(j2gort.Ext_java_util_logging_Level_F_INFO, "Tracer log requested for this thread but was not initialized using Tracer.initCurrentThreadTrace().", j2gort.ConstructShimThrowable("java.lang.Throwable"))
		return
	}
	if !(trace.IsEmpty()) {
		StaticGet_Tracer_Logger().Log_3_Ext_java_util_logging_Level_string_any(j2gort.Ext_java_util_logging_Level_F_INFO, "timers:\n{0}", j2gort.BoxString(Tracer_GetCurrentThreadTraceReport()))
	}
}

// Java static method: Tracer.clearCurrentThreadTrace()
func Tracer_ClearCurrentThreadTrace() {
	clinit_Tracer()
	Tracer_ClearThreadTrace()
}

// Java static method: Tracer.logAndClearCurrentThreadTrace()
func Tracer_LogAndClearCurrentThreadTrace() {
	clinit_Tracer()
	Tracer_LogCurrentThreadTrace()
	Tracer_ClearThreadTrace()
}

// Java static method: Tracer.setPrettyPrint(boolean enabled)
func Tracer_SetPrettyPrint(enabled bool) {
	clinit_Tracer()
	_ = enabled
	StaticSet_Tracer_DefaultPrettyPrint(enabled)
}

// Java static method: Tracer.enableTypeMaps()
func Tracer_EnableTypeMaps() {
	clinit_Tracer()
	staticSyncMu_Tracer.Lock()
	defer staticSyncMu_Tracer.Unlock()
	if StaticGet_Tracer_TypeToCountMap() == nil {
		StaticSet_Tracer_TypeToCountMap(NewTracer_AtomicTracerStatMap())
		StaticSet_Tracer_TypeToSilentMap(NewTracer_AtomicTracerStatMap())
		StaticSet_Tracer_TypeToTimeMap(NewTracer_AtomicTracerStatMap())
	}
}

// Java static method: Tracer.getTypeToCountMap()
func Tracer_GetTypeToCountMap() j2gort.Ext_java_util_Map {
	clinit_Tracer()
	if StaticGet_Tracer_TypeToCountMap() != nil {
		return StaticGet_Tracer_TypeToCountMap().GetMap()
	} else {
		return nil
	}
}

// Java static method: Tracer.getTypeToSilentMap()
func Tracer_GetTypeToSilentMap() j2gort.Ext_java_util_Map {
	clinit_Tracer()
	if StaticGet_Tracer_TypeToSilentMap() != nil {
		return StaticGet_Tracer_TypeToSilentMap().GetMap()
	} else {
		return nil
	}
}

// Java static method: Tracer.getTypeToTimeMap()
func Tracer_GetTypeToTimeMap() j2gort.Ext_java_util_Map {
	clinit_Tracer()
	if StaticGet_Tracer_TypeToTimeMap() != nil {
		return StaticGet_Tracer_TypeToTimeMap().GetMap()
	} else {
		return nil
	}
}

// Java static method: Tracer.getStatsForType(java.lang.String type)
func Tracer_GetStatsForType(type_ string) *Tracer_Stat {
	clinit_Tracer()
	_ = type_
	var stat *Tracer_Stat = Tracer_GetThreadTrace().stats.Get(type_)
	_ = stat
	if stat != nil {
		return stat
	} else {
		return StaticGet_Tracer_ZERO_STAT()
	}
}

// Java static method: Tracer.formatTime(long time)
func Tracer_FormatTime(time int64) string {
	clinit_Tracer()
	_ = time
	var sec int32 = int32(((time / int64(int32(1000))) % int64(int32(60))))
	_ = sec
	var ms int32 = int32((time % int64(int32(1000))))
	_ = ms
	return j2gort.Ext_java_lang_String_Format_2("%02d.%03d", []any{j2gort.ValueOfInt(sec), j2gort.ValueOfInt(ms)})
}

// Java static method: Tracer.getThreadTrace()
func Tracer_GetThreadTrace() *Tracer_ThreadTrace {
	clinit_Tracer()
	var t *Tracer_ThreadTrace = j2gort.CastPtr[Tracer_ThreadTrace](StaticGet_Tracer_Traces().Get())
	_ = t
	if t == nil {
		t = NewTracer_ThreadTrace()
		{
			r_9 := t
			r_9.prettyPrint = StaticGet_Tracer_DefaultPrettyPrint()
		}
		StaticGet_Tracer_Traces().Set(j2gort.PtrToAny(t))
	}
	return t
}

// Java static method: Tracer.clearThreadTrace()
func Tracer_ClearThreadTrace() {
	clinit_Tracer()
	StaticGet_Tracer_Traces().Remove()
}

// Java interface: com.google.javascript.jscomp.Tracer.InternalClock
type Tracer_InternalClock interface {
	Implements_jscomp_Tracer_InternalClock()
	IsJavaNil_() bool
	CurrentTimeMillis() int64
}

func ToIface_Tracer_InternalClock[T interface {
	Tracer_InternalClock
	comparable
}](v T) Tracer_InternalClock {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Tracer_InternalClock[T Tracer_InternalClock](v T) Tracer_InternalClock {
	var vi Tracer_InternalClock = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Tracer.InternalClock", "com.google.javascript.jscomp.Tracer$InternalClock", func(v_ any) bool {
		return j2gort.InstanceOf[Tracer_InternalClock](v_)
	})
}

// Java class: com.google.javascript.jscomp.Tracer.Stat
type Tracer_Stat struct {
	extraInfo []int32 `gson:"extraInfo"`
	count     int32   `gson:"count"`
	silent    int32   `gson:"silent"`
	clockTime int32   `gson:"clockTime"`
	idHash_   atomic.Int32
}

func (this_ *Tracer_Stat) BindSelf_(s any) {
	_ = s
}
func (this_ *Tracer_Stat) IsJavaNil_() bool { return this_ == nil }
func (*Tracer_Stat) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Tracer$Stat", "Stat"
}
func (this_ *Tracer_Stat) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Tracer_Stat) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Tracer.Stat", "com.google.javascript.jscomp.Tracer$Stat", func(v_ any) bool {
		return j2gort.InstanceOf[*Tracer_Stat](v_)
	})
}

func (this_ *Tracer_Stat) initFields_Tracer_Stat() {
	_ = this_
}

// Java constructor: Stat()
func NewTracer_Stat() *Tracer_Stat {
	this_ := &Tracer_Stat{}
	InitTracer_Stat(this_)
	return this_
}

func InitTracer_Stat(this_ *Tracer_Stat) {
	this_.initFields_Tracer_Stat()
}

// Java method: Stat.getCount()
func (this_ *Tracer_Stat) GetCount() int32 {
	return this_.count
}

// Java method: Stat.getSilentCount()
func (this_ *Tracer_Stat) GetSilentCount() int32 {
	return this_.silent
}

// Java method: Stat.getTotalTime()
func (this_ *Tracer_Stat) GetTotalTime() int32 {
	return this_.clockTime
}

// Java method: Stat.getExtraInfo(int index)
func (this_ *Tracer_Stat) GetExtraInfo(index int32) int32 {
	_ = index
	if index >= int32(len(this_.extraInfo)) {
		return int32(0)
	} else {
		return this_.extraInfo[index]
	}
}

// Java class: com.google.javascript.jscomp.Tracer.Event
type Tracer_Event struct {
	tracer  *Tracer `gson:"tracer"`
	idHash_ atomic.Int32
	isStart bool `gson:"isStart"`
}

func (this_ *Tracer_Event) BindSelf_(s any) {
	_ = s
}
func (this_ *Tracer_Event) IsJavaNil_() bool { return this_ == nil }
func (*Tracer_Event) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Tracer$Event", "Event"
}
func (this_ *Tracer_Event) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Tracer_Event) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Tracer.Event", "com.google.javascript.jscomp.Tracer$Event", func(v_ any) bool {
		return j2gort.InstanceOf[*Tracer_Event](v_)
	})
}

func (this_ *Tracer_Event) initFields_Tracer_Event() {
	_ = this_
}

// Java constructor: Event(boolean start, com.google.javascript.jscomp.Tracer t)
func NewTracer_Event(start bool, t *Tracer) *Tracer_Event {
	this_ := &Tracer_Event{}
	InitTracer_Event(this_, start, t)
	return this_
}

func InitTracer_Event(this_ *Tracer_Event, start bool, t *Tracer) {
	_ = start
	_ = t
	this_.initFields_Tracer_Event()
	this_.isStart = start
	this_.tracer = t
}

// Java method: Event.eventTime()
func (this_ *Tracer_Event) EventTime() int64 {
	if this_.isStart {
		return this_.tracer.startTimeMs
	} else {
		return this_.tracer.stopTimeMs
	}
}

// Java method: Event.toString(long prevEventTime, java.lang.String indent, int digitsColWidth)
func (this_ *Tracer_Event) ToString(prevEventTime int64, indent string, digitsColWidth int32) string {
	_ = prevEventTime
	_ = indent
	_ = digitsColWidth
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_int32(int32(120))
	_ = sb
	if prevEventTime == int64(int32(-1)) {
		Tracer_AppendSpaces(sb, digitsColWidth)
	} else {
		sb.Append_1_string(Tracer_LongToPaddedString((this_.EventTime() - prevEventTime), digitsColWidth))
	}
	sb.Append_1_uint16(uint16(32))
	sb.Append_1_string(Tracer_FormatTime(this_.EventTime()))
	if this_.isStart {
		sb.Append_1_string(" Start ")
		Tracer_AppendSpaces(sb, digitsColWidth)
		sb.Append_1_string("   ")
	} else {
		sb.Append_1_string(" Done ")
		var delta int64 = (this_.tracer.stopTimeMs - this_.tracer.startTimeMs)
		_ = delta
		sb.Append_1_string(Tracer_LongToPaddedString(delta, digitsColWidth))
		sb.Append_1_string(" ms ")
		if this_.tracer.extraTracingValues != nil {
			{
				var i int32 = int32(0)
				_ = i
				for ; i < int32(len(this_.tracer.extraTracingValues)); i = i + 1 {
					delta = this_.tracer.extraTracingValues[i]
					sb.Append_1_string(j2gort.Ext_java_lang_String_Format_2("%4d", []any{j2gort.ValueOfLong(delta)}))
					sb.Append_1_string(j2gort.AsIface[Tracer_TracingStatistic](StaticGet_Tracer_ExtraTracingStatistics().Get(i)).GetUnits())
					sb.Append_1_string(";  ")
				}
			}
		}
	}
	sb.Append_1_string(indent)
	sb.Append_1_any_9887(j2gort.PtrToAny(this_.tracer))
	return sb.String()
}

// Java class: com.google.javascript.jscomp.Tracer.ThreadTrace
type Tracer_ThreadTrace struct {
	events                       []*Tracer_Event                         `gson:"events"`
	outstandingEvents            j2gort.OrderedSet[*Tracer]              `gson:"outstandingEvents"`
	stats                        j2gort.OrderedMap[string, *Tracer_Stat] `gson:"stats"`
	defaultSilenceThreshold      int32                                   `gson:"defaultSilenceThreshold"`
	idHash_                      atomic.Int32
	isOutstandingEventsTruncated bool `gson:"isOutstandingEventsTruncated"`
	isEventsTruncated            bool `gson:"isEventsTruncated"`
	isInitialized                bool `gson:"isInitialized"`
	prettyPrint                  bool `gson:"prettyPrint"`
}

func (this_ *Tracer_ThreadTrace) BindSelf_(s any) {
	_ = s
}
func (this_ *Tracer_ThreadTrace) IsJavaNil_() bool { return this_ == nil }
func (*Tracer_ThreadTrace) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Tracer$ThreadTrace", "ThreadTrace"
}
func (this_ *Tracer_ThreadTrace) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Tracer_ThreadTrace) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Tracer.ThreadTrace", "com.google.javascript.jscomp.Tracer$ThreadTrace", func(v_ any) bool {
		return j2gort.InstanceOf[*Tracer_ThreadTrace](v_)
	})
}

func (this_ *Tracer_ThreadTrace) initFields_Tracer_ThreadTrace() {
	_ = this_
	this_.events = []*Tracer_Event{}
	this_.outstandingEvents = j2gort.NewOrderedSet[*Tracer]()
	this_.stats = j2gort.NewOrderedMap[string, *Tracer_Stat]()

}

// Java constructor: ThreadTrace()
func NewTracer_ThreadTrace() *Tracer_ThreadTrace {
	this_ := &Tracer_ThreadTrace{}
	InitTracer_ThreadTrace(this_)
	return this_
}

func InitTracer_ThreadTrace(this_ *Tracer_ThreadTrace) {
	this_.initFields_Tracer_ThreadTrace()
}

// Java method: ThreadTrace.init()
func (this_ *Tracer_ThreadTrace) Init() {
	this_.isInitialized = true
}

// Java method: ThreadTrace.isInitialized()
func (this_ *Tracer_ThreadTrace) IsInitialized() bool {
	return this_.isInitialized
}

// Java method: ThreadTrace.startEvent(com.google.javascript.jscomp.Tracer t)
func (this_ *Tracer_ThreadTrace) StartEvent(t *Tracer) {
	_ = t
	this_.events = append(this_.events, NewTracer_Event(true, t))
	var notAlreadyOutstanding bool = this_.outstandingEvents.Add(t)
	_ = notAlreadyOutstanding
	if !(notAlreadyOutstanding) {
		j2gort.CheckState(false)
	}
}

// Java method: ThreadTrace.endEvent(com.google.javascript.jscomp.Tracer t, int silenceThreshold)
func (this_ *Tracer_ThreadTrace) EndEvent(t *Tracer, silenceThreshold int32) {
	_ = t
	_ = silenceThreshold
	var wasOutstanding bool = this_.outstandingEvents.Remove(t)
	_ = wasOutstanding
	if !(wasOutstanding) {
		if this_.isOutstandingEventsTruncated {
			StaticGet_Tracer_Logger().Log_3_Ext_java_util_logging_Level_string_Throwable(j2gort.Ext_java_util_logging_Level_F_WARNING, "event not found, probably because the event stack overflowed and was truncated", j2gort.ConstructShimThrowable("java.lang.Throwable"))
		} else {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
		}
	}
	var elapsed int64 = (t.stopTimeMs - t.startTimeMs)
	_ = elapsed
	if silenceThreshold == int32(-1) {
		silenceThreshold = this_.defaultSilenceThreshold
	}
	if elapsed < int64(silenceThreshold) {
		var removed bool = false
		_ = removed
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(this_.events)); i = i + 1 {
				var e *Tracer_Event = this_.events[i]
				_ = e
				if e.tracer == t {
					if !(e.isStart) {
						j2gort.CheckState(false)
					}
					this_.events = j2gort.SliceRemoveAt(this_.events, i)
					removed = true
					break
				}
			}
		}
		if !(removed || this_.isEventsTruncated) {
			j2gort.CheckState(false)
		}
	} else {
		this_.events = append(this_.events, NewTracer_Event(false, t))
	}
	if t.type_ != j2gort.NullString {
		var stat *Tracer_Stat = this_.stats.Get(t.type_)
		_ = stat
		if stat == nil {
			stat = NewTracer_Stat()
			if !(StaticGet_Tracer_ExtraTracingStatistics().IsEmpty()) {
				{
					r_10 := stat
					r_10.extraInfo = make([]int32, StaticGet_Tracer_ExtraTracingStatistics().Size())
				}
			}
			this_.stats.Put(t.type_, stat)
		}
		r_11 := stat
		r_11.count = r_11.count + 1
		if StaticGet_Tracer_TypeToCountMap() != nil {
			StaticGet_Tracer_TypeToCountMap().IncrementBy(t.type_, int64(int32(1)))
		}
		r_12 := stat
		r_12.clockTime = int32((int64(r_12.clockTime) + elapsed))
		if StaticGet_Tracer_TypeToTimeMap() != nil {
			StaticGet_Tracer_TypeToTimeMap().IncrementBy(t.type_, elapsed)
		}
		if (stat.extraInfo != nil) && (t.extraTracingValues != nil) {
			var overlapLength int32 = min(int32(len(stat.extraInfo)), int32(len(t.extraTracingValues)))
			_ = overlapLength
			{
				var i_2 int32 = int32(0)
				_ = i_2
				for ; i_2 < overlapLength; i_2 = i_2 + 1 {
					a_13 := stat.extraInfo
					i_14 := i_2
					a_13[i_14] = int32((int64(a_13[i_14]) + t.extraTracingValues[i_2]))
					var map_ *Tracer_AtomicTracerStatMap = j2gort.AsIface[Tracer_TracingStatistic](StaticGet_Tracer_ExtraTracingStatistics().Get(i_2)).GetTracingStat()
					_ = map_
					if map_ != nil {
						map_.IncrementBy(t.type_, t.extraTracingValues[i_2])
					}
				}
			}
		}
		if elapsed < int64(silenceThreshold) {
			r_15 := stat
			r_15.silent = r_15.silent + 1
			if StaticGet_Tracer_TypeToSilentMap() != nil {
				StaticGet_Tracer_TypeToSilentMap().IncrementBy(t.type_, int64(int32(1)))
			}
		}
	}
}

// Java method: ThreadTrace.isEmpty()
func (this_ *Tracer_ThreadTrace) IsEmpty() bool {
	return (len(this_.events) == 0 && this_.outstandingEvents.IsEmpty())
}

// Java method: ThreadTrace.truncateOutstandingEvents()
func (this_ *Tracer_ThreadTrace) TruncateOutstandingEvents() {
	this_.isOutstandingEventsTruncated = true
	this_.outstandingEvents.Clear()
}

// Java method: ThreadTrace.truncateEvents()
func (this_ *Tracer_ThreadTrace) TruncateEvents() {
	this_.isEventsTruncated = true
	clear(this_.events)
	this_.events = this_.events[:0]
}

// Java method: ThreadTrace.toString()
func (this_ *Tracer_ThreadTrace) String() string {
	var numDigits int32 = this_.priv_ThreadTrace_getMaxDigits_5377()
	_ = numDigits
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	var etime int64 = int64(int32(-1))
	_ = etime
	var indent j2gort.Ext_java_util_Deque
	_ = indent
	if this_.prettyPrint {
		indent = j2gort.NewExt_java_util_ArrayDeque()
	} else {
		indent = nil
	}
	for _, e := range this_.events {
		_ = e
		if (this_.prettyPrint && !(e.isStart)) && !(indent.IsEmpty()) {
			indent.Pop()
		}
		sb.Append_1_string(" ")
		if this_.prettyPrint {
			sb.Append_1_string(e.ToString(etime, j2gort.Ext_com_google_common_base_Joiner_On_1_string("").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](indent)), numDigits))
		} else {
			sb.Append_1_string(e.ToString(etime, "", int32(4)))
		}
		etime = e.EventTime()
		sb.Append_1_uint16(uint16(10))
		if this_.prettyPrint && e.isStart {
			indent.Push("|  ")
		}
	}

	if !(this_.outstandingEvents.IsEmpty()) {
		var now int64 = StaticGet_Tracer_Clock().CurrentTimeMillis()
		_ = now
		sb.Append_1_string(" Unstopped timers:\n")
		{
			cur_18 := this_.outstandingEvents.Iter()
			for cur_18.HasNext() {

				t := cur_18.
					Next()
				_ = t
				j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string("  ")).Append_1_any_9887(j2gort.PtrToAny(t))).Append_1_string(" (")).Append_1_int64((now - t.startTimeMs))).Append_1_string(" ms, started at ")).Append_1_string(Tracer_FormatTime(t.startTimeMs))).Append_1_string(")\n")
			}
		}
	}
	{
		cur_20 := this_.stats.IterEntries()
		for cur_20.NextEntry() {
			var stat *Tracer_Stat = cur_20.Value()
			_ = stat
			if stat.count > int32(1) {
				j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string(" TOTAL ")).Append_1_string(cur_20.Key())).Append_1_string(" ")).Append_1_int32(stat.count)).Append_1_string(" (")).Append_1_int32(stat.clockTime)).Append_1_string(" ms")
				if stat.extraInfo != nil {
					{
						var i int32 = int32(0)
						_ = i
						for ; i < int32(len(stat.extraInfo)); i = i + 1 {
							sb.Append_1_string("; ")
							j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_int32(stat.extraInfo[i])).Append_1_uint16(uint16(32))).Append_1_string(j2gort.AsIface[Tracer_TracingStatistic](StaticGet_Tracer_ExtraTracingStatistics().Get(i)).GetUnits())
						}
					}
				}
				sb.Append_1_string(")\n")
			}
		}
	}
	return sb.String()
}

// Java method: ThreadTrace.getMaxDigits()
func (this_ *Tracer_ThreadTrace) priv_ThreadTrace_getMaxDigits_5377() int32 {
	var etime int64 = int64(int32(-1))
	_ = etime
	var maxTime int64 = int64(int32(0))
	_ = maxTime
	for _, e := range this_.events {
		_ = e
		if etime != int64(int32(-1)) {
			var time int64 = (e.EventTime() - etime)
			_ = time
			maxTime = max(maxTime, time)
		}
		if !(e.isStart) {
			var time_2 int64 = (e.tracer.stopTimeMs - e.tracer.startTimeMs)
			_ = time_2
			maxTime = max(maxTime, time_2)
		}
		etime = e.EventTime()
	}

	return max(int32(3), Tracer_NumDigits(maxTime))
}

// Java interface: com.google.javascript.jscomp.Tracer.TracingStatistic
type Tracer_TracingStatistic interface {
	Implements_jscomp_Tracer_TracingStatistic()
	IsJavaNil_() bool
	Start(*j2gort.Ext_java_lang_Thread) int64
	Stop(*j2gort.Ext_java_lang_Thread) int64
	Enable() bool
	GetTracingStat() *Tracer_AtomicTracerStatMap
	GetUnits() string
}

func ToIface_Tracer_TracingStatistic[T interface {
	Tracer_TracingStatistic
	comparable
}](v T) Tracer_TracingStatistic {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Tracer_TracingStatistic[T Tracer_TracingStatistic](v T) Tracer_TracingStatistic {
	var vi Tracer_TracingStatistic = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Tracer.TracingStatistic", "com.google.javascript.jscomp.Tracer$TracingStatistic", func(v_ any) bool {
		return j2gort.InstanceOf[Tracer_TracingStatistic](v_)
	})
}

// Java class: com.google.javascript.jscomp.Tracer.AtomicTracerStatMap
type Tracer_AtomicTracerStatMap struct {
	map_    j2gort.Ext_java_util_concurrent_ConcurrentMap `gson:"map"`
	idHash_ atomic.Int32
}

func (this_ *Tracer_AtomicTracerStatMap) BindSelf_(s any) {
	_ = s
}
func (this_ *Tracer_AtomicTracerStatMap) IsJavaNil_() bool { return this_ == nil }
func (*Tracer_AtomicTracerStatMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Tracer$AtomicTracerStatMap", "AtomicTracerStatMap"
}
func (this_ *Tracer_AtomicTracerStatMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Tracer_AtomicTracerStatMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Tracer.AtomicTracerStatMap", "com.google.javascript.jscomp.Tracer$AtomicTracerStatMap", func(v_ any) bool {
		return j2gort.InstanceOf[*Tracer_AtomicTracerStatMap](v_)
	})
}

func (this_ *Tracer_AtomicTracerStatMap) initFields_Tracer_AtomicTracerStatMap() {
	_ = this_
	this_.map_ = j2gort.NewExt_java_util_concurrent_ConcurrentHashMap()
}

// Java constructor: AtomicTracerStatMap()
func NewTracer_AtomicTracerStatMap() *Tracer_AtomicTracerStatMap {
	this_ := &Tracer_AtomicTracerStatMap{}
	InitTracer_AtomicTracerStatMap(this_)
	return this_
}

func InitTracer_AtomicTracerStatMap(this_ *Tracer_AtomicTracerStatMap) {
	this_.initFields_Tracer_AtomicTracerStatMap()
}

// Java method: AtomicTracerStatMap.incrementBy(java.lang.String key, long delta)
func (this_ *Tracer_AtomicTracerStatMap) IncrementBy(key string, delta int64) {
	_ = key
	_ = delta
	var oldValue *j2gort.BoxLong = j2gort.CastPtr[j2gort.BoxLong](j2gort.MapGetStr(this_.map_, key))
	_ = oldValue
	if oldValue == nil {
		oldValue = j2gort.CastPtr[j2gort.BoxLong](this_.map_.PutIfAbsent(j2gort.BoxString(key), j2gort.ValueOfLong(delta)))
		if oldValue == nil {
			return
		} else {
		}
	}
	for {
		if this_.map_.Replace_3(j2gort.BoxString(key), j2gort.PtrToAny(oldValue), j2gort.ValueOfLong((j2gort.UnboxLong(oldValue) + delta))) {
			break
		}
		oldValue = j2gort.CastPtr[j2gort.BoxLong](j2gort.MapGetStr(this_.map_, key))
	}
}

// Java method: AtomicTracerStatMap.getMap()
func (this_ *Tracer_AtomicTracerStatMap) GetMap() j2gort.Ext_java_util_Map {
	return j2gort.AsIface[j2gort.Ext_java_util_Map](this_.map_)
}

// Java class:
type Tracer__Anon_98 struct {
	idHash_ atomic.Int32
}

func (this_ *Tracer__Anon_98) BindSelf_(s any) {
	_ = s
}
func (this_ *Tracer__Anon_98) IsJavaNil_() bool { return this_ == nil }
func (*Tracer__Anon_98) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Tracer$1", ""
}
func (this_ *Tracer__Anon_98) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Tracer__Anon_98) HasIdentityEquals_()                     {}
func (*Tracer__Anon_98) Implements_jscomp_Tracer_InternalClock() {}

func (this_ *Tracer__Anon_98) initFields_Tracer__Anon_98() {
	_ = this_
}

func NewTracer__Anon_98() *Tracer__Anon_98 {
	this_ := &Tracer__Anon_98{}
	this_.initFields_Tracer__Anon_98()
	return this_
}

// Java method: .currentTimeMillis()
func (this_ *Tracer__Anon_98) CurrentTimeMillis() int64 {
	return javalang.System_CurrentTimeMillis()
}

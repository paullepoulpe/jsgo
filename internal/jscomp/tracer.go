// Port of com/google/javascript/jscomp/Tracer.java @ closure-compiler v20261021

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

package jscomp

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Tracer provides a simple way to trace the handling of a request.
//
// A thread must call Tracer_InitCurrentThreadTrace() to enable Tracer logging,
// otherwise Tracer does nothing.
type Tracer struct {
	// The type for grouping traces, may be null
	type_ *string

	// A comment string for the report
	comment string

	// Start time of the trace
	startTimeMs int64

	// Stop time of the trace, non-final
	stopTimeMs int64

	// Record our starter thread in order to trap Traces that are started in
	// one thread and stopped in another
	startThread *jutil.Thread
}

// Tracer_logger is package-private for access from unit tests.
var Tracer_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.Tracer")

// Tracer_MAX_TRACE_SIZE limits the number of events in a Trace in order to
// catch memory leaks (a thread that keeps logging events and never clears them).
const Tracer_MAX_TRACE_SIZE = 1000

// Tracer_InternalClock is for unit testing.
type Tracer_InternalClock interface {
	// Java: Tracer.InternalClock.currentTimeMillis()
	CurrentTimeMillis() int64
}

type tracer_systemClock struct{}

// CurrentTimeMillis is System.currentTimeMillis().
func (tracer_systemClock) CurrentTimeMillis() int64 { return time.Now().UnixMilli() }

// Tracer_clock is the default clock that calls through to the system clock.
var Tracer_clock Tracer_InternalClock = tracer_systemClock{}

// tracer_Throwable is `new Throwable(message)`, used to log a stack trace marker.
type tracer_Throwable struct{ jutil.ThrowableBase }

func newTracer_Throwable(message *string) *tracer_Throwable {
	t := &tracer_Throwable{}
	jutil.InitThrowable(t, "java.lang.Throwable", message, nil)
	return t
}

// NewTracerWithType creates and starts a tracer. Both type and comment may be null.
//
// Java: Tracer(String, String)
func NewTracerWithType(type_ *string, comment *string) *Tracer {
	t := &Tracer{}
	t.type_ = type_
	if comment != nil {
		t.comment = *comment
	}
	t.startTimeMs = Tracer_clock.CurrentTimeMillis()
	t.startThread = jutil.Thread_CurrentThread()

	trace := Tracer_GetThreadTrace()

	// Do nothing if the current thread trace wasn't initialized.
	if !trace.IsInitialized() {
		return t
	}

	// Check if we are creating too many Tracers.
	if len(trace.events) >= Tracer_MAX_TRACE_SIZE {
		Tracer_logger.LogWithThrown(jutil.Level_WARNING,
			"Giant thread trace. Too many Tracers created. "+
				"Clearing to avoid memory leak.",
			newTracer_Throwable(jutil.StringPtr(trace.String())))
		trace.TruncateEvents()
	}

	// Check if we forgot to close the Tracers.
	if trace.outstandingEvents.Size() >= Tracer_MAX_TRACE_SIZE {
		Tracer_logger.LogWithThrown(jutil.Level_WARNING,
			"Too many outstanding Tracers. Tracer.stop() is missing "+
				"or Tracer.stop() is not wrapped in a "+
				"try/finally block. "+
				"Clearing to avoid memory leak.",
			newTracer_Throwable(jutil.StringPtr(trace.String())))
		trace.TruncateOutstandingEvents()
	}

	trace.StartEvent(t)
	return t
}

// NewTracer creates a tracer that isn't summed as part of a total.
//
// Java: Tracer(String)
func NewTracer(comment string) *Tracer {
	return NewTracerWithType(nil, &comment)
}

// tracer_longToPaddedString converts 'v' to a string and pads it with up
// to 16 spaces for improved alignment.
//
// Java: Tracer.longToPaddedString(long, int)
func tracer_longToPaddedString(v int64, digitsColumnWidth int) string {
	digitWidth := tracer_numDigits(v)
	var sb strings.Builder
	Tracer_AppendSpaces(&sb, digitsColumnWidth-digitWidth)
	sb.WriteString(strconv.FormatInt(v, 10))
	return sb.String()
}

// tracer_numDigits gets the number of digits in an integer when printed in
// base 10. Assumes a positive integer.
//
// Java: Tracer.numDigits(long)
func tracer_numDigits(n int64) int {
	i := 0
	for {
		i++
		n = n / 10
		if !(n > 0) {
			break
		}
	}
	return i
}

// Tracer_AppendSpaces gets a string of spaces of the length specified.
//
// Java: Tracer.appendSpaces(StringBuilder, int)
func Tracer_AppendSpaces(sb *strings.Builder, numSpaces int) {
	if numSpaces > 16 {
		Tracer_logger.Warning("Tracer.appendSpaces called with large numSpaces")
		numSpaces = 16
	}
	if numSpaces > 0 {
		const spaces16 = "                "
		sb.WriteString(spaces16[:numSpaces])
	}
}

// StopWithThreshold stops the trace. This may only be done once and must
// be done from the same thread that started it.
//
// Traces for time less than silenceThreshold ms will be left out of the
// trace report. A value of -1 indicates that the current ThreadTrace
// silence_threshold should be used.
//
// Returns the time that this trace actually ran.
//
// Java: Tracer.stop(int)
func (t *Tracer) StopWithThreshold(silenceThreshold int) int64 {
	jutil.CheckState(jutil.Thread_CurrentThread() == t.startThread)

	trace := Tracer_GetThreadTrace()
	// Do nothing if the thread trace was not initialized.
	if !trace.IsInitialized() {
		return 0
	}

	t.stopTimeMs = Tracer_clock.CurrentTimeMillis()

	trace.EndEvent(t, silenceThreshold)
	return t.stopTimeMs - t.startTimeMs
}

// Stop stops the trace using the default silenceThreshold.
//
// Returns the time that this trace actually ran.
//
// Java: Tracer.stop()
func (t *Tracer) Stop() int64 {
	return t.StopWithThreshold(-1)
}

// Java: Tracer.toString()
func (t *Tracer) String() string {
	if t.type_ == nil {
		return t.comment
	}
	return "[" + *t.type_ + "] " + t.comment
}

// Java: Tracer.setDefaultSilenceThreshold(int)
func Tracer_SetDefaultSilenceThreshold(threshold int) {
	Tracer_GetThreadTrace().defaultSilenceThreshold = threshold
}

// Tracer_InitCurrentThreadTrace initializes the trace associated with the
// current thread by clearing out any existing trace.
//
// Java: Tracer.initCurrentThreadTrace()
func Tracer_InitCurrentThreadTrace() {
	events := Tracer_GetThreadTrace()
	if !events.IsEmpty() {
		Tracer_logger.LogWithThrown(jutil.Level_WARNING,
			"Non-empty timer log:\n"+events.String(),
			newTracer_Throwable(nil))
		Tracer_ClearThreadTrace()
		events = Tracer_GetThreadTrace()
	}

	// Mark the thread trace as initialized.
	events.Init()
}

// Java: Tracer.initCurrentThreadTrace(int)
func Tracer_InitCurrentThreadTraceWithThreshold(defaultSilenceThreshold int) {
	Tracer_InitCurrentThreadTrace()
	Tracer_SetDefaultSilenceThreshold(defaultSilenceThreshold)
}

// Tracer_GetCurrentThreadTraceReport returns a timer report.
//
// Java: Tracer.getCurrentThreadTraceReport()
func Tracer_GetCurrentThreadTraceReport() string {
	return Tracer_GetThreadTrace().String()
}

// Tracer_LogCurrentThreadTrace logs a timer report.
//
// Java: Tracer.logCurrentThreadTrace()
func Tracer_LogCurrentThreadTrace() {
	trace := Tracer_GetThreadTrace()
	if !trace.IsInitialized() {
		Tracer_logger.LogWithThrown(
			jutil.Level_INFO,
			"Tracer log requested for this thread but was not "+
				"initialized using Tracer.initCurrentThreadTrace().",
			newTracer_Throwable(nil))
		return
	}

	if !trace.IsEmpty() {
		Tracer_logger.LogWithParam(jutil.Level_INFO, "timers:\n{0}", Tracer_GetCurrentThreadTraceReport())
	}
}

// Tracer_ClearCurrentThreadTrace throws away any Trace associated with the
// current thread.
//
// Java: Tracer.clearCurrentThreadTrace()
func Tracer_ClearCurrentThreadTrace() {
	Tracer_ClearThreadTrace()
}

// Tracer_LogAndClearCurrentThreadTrace is logCurrentThreadTrace() then
// clearCurrentThreadTrace().
//
// Java: Tracer.logAndClearCurrentThreadTrace()
func Tracer_LogAndClearCurrentThreadTrace() {
	Tracer_LogCurrentThreadTrace()
	Tracer_ClearThreadTrace()
}

// tracer_Stat holds the statistics for a given tracer type.
type tracer_Stat struct {
	count     int
	silent    int
	clockTime int
}

// tracer_formatTime returns the sec.ms part of time (if time =
// "20:06:11.566", "11.566").
//
// Java: Tracer.formatTime(long)
func tracer_formatTime(t int64) string {
	sec := int32((t / 1000) % 60)
	ms := int32(t % 1000)
	return jutil.StringFormat("%02d.%03d", sec, ms)
}

// tracer_Event is created every time a Tracer is created or stopped.
type tracer_Event struct {
	isStart bool // else is_stop
	tracer  *Tracer
}

// Java: Tracer.Event(boolean, Tracer)
func newTracer_Event(start bool, t *Tracer) *tracer_Event {
	return &tracer_Event{isStart: start, tracer: t}
}

// Java: Tracer.Event.eventTime()
func (e *tracer_Event) EventTime() int64 {
	if e.isStart {
		return e.tracer.startTimeMs
	}
	return e.tracer.stopTimeMs
}

// ToString converts the event to a formatted string.
//
// Java: Tracer.Event.toString(long, String, int)
func (e *tracer_Event) ToString(prevEventTime int64, indent string, digitsColWidth int) string {
	var sb strings.Builder

	if prevEventTime == -1 {
		Tracer_AppendSpaces(&sb, digitsColWidth)
	} else {
		sb.WriteString(tracer_longToPaddedString(
			e.EventTime()-prevEventTime, digitsColWidth))
	}

	sb.WriteByte(' ')
	sb.WriteString(tracer_formatTime(e.EventTime()))
	if e.isStart {
		sb.WriteString(" Start ")
		Tracer_AppendSpaces(&sb, digitsColWidth)
		sb.WriteString("   ")
	} else {
		sb.WriteString(" Done ")
		delta := e.tracer.stopTimeMs - e.tracer.startTimeMs
		sb.WriteString(tracer_longToPaddedString(delta, digitsColWidth))
		sb.WriteString(" ms ")
	}
	sb.WriteString(indent)
	sb.WriteString(e.tracer.String())
	return sb.String()
}

// Tracer_ThreadTrace stores a thread's Trace.
type Tracer_ThreadTrace struct {
	// Events taking less than this number of milliseconds are not reported.
	defaultSilenceThreshold int

	// The Events corresponding to each startEvent/stopEvent
	events []*tracer_Event

	// Tracers that have not had their .stop() called
	outstandingEvents *jutil.LinkedHashSet[*Tracer]

	// Map from type to Stat object
	stats *jutil.LinkedHashMap[string, *tracer_Stat]

	// True if outstandingEvents has been cleared because we exceeded the max trace limit.
	isOutstandingEventsTruncated bool

	// True if events has been cleared because we exceeded the max trace limit.
	isEventsTruncated bool

	// Set to true if Tracer#initCurrentThreadTrace() was called by the current thread.
	isInitialized bool
}

// NewTracer_ThreadTrace is new ThreadTrace().
func NewTracer_ThreadTrace() *Tracer_ThreadTrace {
	return &Tracer_ThreadTrace{
		outstandingEvents: jutil.NewLinkedHashSet[*Tracer](),
		stats:             jutil.NewLinkedHashMap[string, *tracer_Stat](),
	}
}

// Java: Tracer.ThreadTrace.init()
func (tt *Tracer_ThreadTrace) Init() {
	tt.isInitialized = true
}

// Java: Tracer.ThreadTrace.isInitialized()
func (tt *Tracer_ThreadTrace) IsInitialized() bool {
	return tt.isInitialized
}

// Java: Tracer.ThreadTrace.startEvent(Tracer)
func (tt *Tracer_ThreadTrace) StartEvent(t *Tracer) {
	tt.events = append(tt.events, newTracer_Event(true, t))
	notAlreadyOutstanding := tt.outstandingEvents.Add(t)
	jutil.CheckState(notAlreadyOutstanding)
}

// Java: Tracer.ThreadTrace.endEvent(Tracer, int)
func (tt *Tracer_ThreadTrace) EndEvent(t *Tracer, silenceThreshold int) {
	wasOutstanding := tt.outstandingEvents.Remove(t)
	if !wasOutstanding {
		if tt.isOutstandingEventsTruncated {
			Tracer_logger.LogWithThrown(jutil.Level_WARNING,
				"event not found, probably because the event stack "+
					"overflowed and was truncated",
				newTracer_Throwable(nil))
		} else {
			panic(jutil.NewIllegalStateException())
		}
	}

	elapsed := t.stopTimeMs - t.startTimeMs

	if silenceThreshold == -1 {
		silenceThreshold = tt.defaultSilenceThreshold
	}

	if elapsed < int64(silenceThreshold) {
		removed := false
		for i, e := range tt.events {
			if e.tracer == t {
				jutil.CheckState(e.isStart)
				tt.events = slices.Delete(tt.events, i, i+1)
				removed = true
				break
			}
		}
		jutil.CheckState(removed || tt.isEventsTruncated)
	} else {
		tt.events = append(tt.events, newTracer_Event(false, t))
	}

	if t.type_ != nil {
		type_ := *t.type_
		stat := tt.stats.Get(type_)
		if stat == nil {
			stat = &tracer_Stat{}
			tt.stats.Put(type_, stat)
		}

		stat.count++
		stat.clockTime = int(int32(int64(stat.clockTime) + elapsed))
		if elapsed < int64(silenceThreshold) {
			stat.silent++
		}
	}
}

// Java: Tracer.ThreadTrace.isEmpty()
func (tt *Tracer_ThreadTrace) IsEmpty() bool {
	return len(tt.events) == 0 && tt.outstandingEvents.IsEmpty()
}

// Java: Tracer.ThreadTrace.truncateOutstandingEvents()
func (tt *Tracer_ThreadTrace) TruncateOutstandingEvents() {
	tt.isOutstandingEventsTruncated = true
	tt.outstandingEvents.Clear()
}

// Java: Tracer.ThreadTrace.truncateEvents()
func (tt *Tracer_ThreadTrace) TruncateEvents() {
	tt.isEventsTruncated = true
	tt.events = nil
}

// Java: Tracer.ThreadTrace.toString()
func (tt *Tracer_ThreadTrace) String() string {
	var sb strings.Builder
	var etime int64 = -1
	for _, e := range tt.events {
		sb.WriteByte(' ')
		sb.WriteString(e.ToString(etime, "", 4))
		etime = e.EventTime()
		sb.WriteByte('\n')
	}

	if tt.outstandingEvents.Size() != 0 {
		now := Tracer_clock.CurrentTimeMillis()

		sb.WriteString(" Unstopped timers:\n")
		for t := range tt.outstandingEvents.All() {
			sb.WriteString("  ")
			sb.WriteString(t.String())
			sb.WriteString(" (")
			sb.WriteString(strconv.FormatInt(now-t.startTimeMs, 10))
			sb.WriteString(" ms, started at ")
			sb.WriteString(tracer_formatTime(t.startTimeMs))
			sb.WriteString(")\n")
		}
	}

	for key, stat := range tt.stats.All() {
		if stat.count > 1 {
			sb.WriteString(" TOTAL ")
			sb.WriteString(key)
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(stat.count))
			sb.WriteString(" (")
			sb.WriteString(strconv.Itoa(stat.clockTime))
			sb.WriteString(" ms)\n")
		}
	}
	return sb.String()
}

// Holds the ThreadTrace for each thread.
var tracer_traces *Tracer_ThreadTrace

// Tracer_GetThreadTrace gets the ThreadTrace for the current thread,
// creating one if necessary.
//
// Java: Tracer.getThreadTrace()
func Tracer_GetThreadTrace() *Tracer_ThreadTrace {
	t := tracer_traces
	if t == nil {
		t = NewTracer_ThreadTrace()
		tracer_traces = t
	}
	return t
}

// Tracer_ClearThreadTrace removes any ThreadTrace associated with the
// current thread.
//
// Java: Tracer.clearThreadTrace()
func Tracer_ClearThreadTrace() {
	tracer_traces = nil
}

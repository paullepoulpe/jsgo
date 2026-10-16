// Port of com/google/javascript/jscomp/Tracer.java @ closure-compiler v20260415

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
	"strings"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Tracer provides a simple way to trace the handling of a request.
//
// By timing likely slow points in the code you can quickly pinpoint why a
// request is slow.
//
// Example usage:
//
//	Tracer.initCurrentThreadTrace(); // must be called in each Thread
//	Tracer wholeRequest = new Tracer(null, "Request " + params);
//	try {
//	  ...
//	  t = new Tracer("Database", "getName()");
//	  try {
//	    name = database.getName();
//	  } finally {
//	    t.stop();
//	  }
//	  ...
//	  t = new Tracer(null, "call sendmail");
//	  try {
//	    sendMessage();
//	  } finally {
//	    t.stop();
//	  }
//	  ...
//	  t = new Tracer("Database", "updateinfo()");
//	  try {
//	    database.updateinfo("new info");
//	  } finally {
//	    t.stop();
//	  }
//	  ...
//	} finally {
//	  if (wholeRequest.stop() > 1000) {
//	    // more than a second, better log
//	    Tracer.logAndClearCurrentThreadTrace();
//	  } else {
//	    Tracer.clearCurrentThreadTrace();
//	  }
//	}
//
// Now slow requests will produce a report like this:
//
//	     10.452 Start        Request cmd=dostuff
//	   3 10.455 Start        [Database] getName()
//	  34 10.489 Done   34 ms [Database] getName()
//	   3 10.491 Start        call sendmail
//	1042 11.533 Done 1042 ms call sendmail
//	   0 11.533 Start        [Database] updateinfo()
//	   3 11.536 Done    3 ms [Database] updateinfo()
//	  64 11.600 Done 1148 ms Request cmd=dostuff
//	 TOTAL Database 2 (37 ms)
//
// If you enabled pretty-printing by calling Tracer#setPrettyPrint, it will
// print more easily readable reports that use indentation to visualize the
// tracer hierarchy and dynamically adjusts the padding to handle large
// durations. Like:
//
//	     10.452 Start        Request cmd=dostuff
//	   3 10.455 Start        | [Database] getName()
//	  34 10.489 Done   34 ms | [Database] getName()
//	   3 10.491 Start        | call sendmail
//	1042 11.533 Done 1042 ms | call sendmail
//	   0 11.533 Start        | [Database] updateinfo()
//	   3 11.536 Done    3 ms | [Database] updateinfo()
//	  64 11.600 Done 1148 ms Request cmd=dostuff
//	 TOTAL Database 2 (37 ms)
//
// Pretty-printing is an application global setting and should only be
// called in the main setup of an application, not in library code.
//
// Now you can easily see that sendmail is causing your problems, not the
// two database calls.
//
// You can easily add additional tracing statistics to your Trace output by
// adding additional tracing statistics. Simply add to your initialization
// code:
//
//	Tracer.addTracingStatistic(myTracingStatistic)
//
// where myTracingStatistic implements the TracingStatistic interface. The
// class com.google.monitoring.tracing.TracingStatistics contains several
// useful statistics such as CPU time, wait time, and memory usage. If you
// add your own tracing statistics, the output is not quite as pretty, but
// includes additional useful information.
//
// If a Trace is given a type (the first argument to the constructor) and
// multiple Traces are done on that type then a "TOTAL line will be produced
// showing the total number of traces and the sum of the time ("TOTAL
// Database 2 (37 ms)" in our example). These traces should be mutually
// exclusive or else the sum won't make sense (the time will be double
// counted if the second starts before the first ends).
//
// It is also possible to have a "silent" Tracer which does not appear in
// the trace because it was faster than the silence threshold. This
// threshold can be set for the for the current ThreadTrace with
// setDefaultSilenceThreshold(threshold), or on a per-Tracer basis with
// t.stop(threshold). Silent tracers are still counted in the type totals,
// so these events are not completely lost.
//
// WARNING: This code makes a big assumption that everything for a given
// trace is done within a single thread. It uses threads to identify
// requests. It is fine to have multiple requests traced in multiple
// simultaneous threads but it is not ok to have any given request traced
// in multiple threads. (the results will be scattered across reports).
//
// Java objects do not support destructors (as in C++) so Tracer is not
// robust when exceptions are thrown. Each Tracer object should be wrapped
// in a try/finally block so that if an exception is thrown, the
// Tracer.stop() method is guaranteed to be called.
//
// A thread must call Tracer#initCurrentThreadTrace() to enable the Tracer
// logging, otherwise Tracer does nothing. The requirement to call
// initCurrentThreadTrace avoids the situation where Tracer is called
// without the explicit knowledge of the application authors because they
// happen to use a class in another package that uses Tracer. If
// Tracer#logCurrentThreadTrace is called without calling
// Tracer#initCurrentThreadTrace(), then a Third Eye WARNING message is
// logged, which should help track down the problem.
//
// Go-only: language-adapter — the port is sequential (PORTING.md §5): the
// per-thread trace (a ThreadLocal in Java) is a single package variable,
// and the start thread is jutil.Thread_CurrentThread().
type Tracer struct {
	// Values returned by extraTracingStatistics
	extraTracingValues []int64

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

// Whether pretty printing is enabled. This is intended to be set once at
// application startup.
var tracer_defaultPrettyPrint bool

// This list is guaranteed to only increase in length. It contains a list
// of additional statistics that the user wants to keep track of.
var tracer_extraTracingStatistics []Tracer_TracingStatistic

// Tracer_MAX_TRACE_SIZE: we limit the number of events in a Trace in order
// to catch memory leaks (a thread that keeps logging events and never
// clears them). This number is arbitrary and can be increased if necessary
// (though if there are more than 1000 events then the Tracer is probably
// being misused).
const Tracer_MAX_TRACE_SIZE = 1000

// Tracer_InternalClock is for unit testing. Can't use
// com.google.common.time.Clock because this code is in base and has
// minimal dependencies.
type Tracer_InternalClock interface {
	// Java: Tracer.InternalClock.currentTimeMillis()
	CurrentTimeMillis() int64
}

type tracer_systemClock struct{}

// CurrentTimeMillis is System.currentTimeMillis().
func (tracer_systemClock) CurrentTimeMillis() int64 { return time.Now().UnixMilli() }

// Tracer_clock is the default clock that calls through to the system
// clock. Can be overridden in unit tests.
var Tracer_clock Tracer_InternalClock = tracer_systemClock{}

// tracer_Throwable is `new Throwable(message)`, used to log a stack trace
// marker.
type tracer_Throwable struct{ jutil.ThrowableBase }

func newTracer_Throwable(message *string) *tracer_Throwable {
	t := &tracer_Throwable{}
	jutil.InitThrowable(t, "java.lang.Throwable", message, nil)
	return t
}

// NewTracerWithType creates and starts a tracer. Both type and comment may
// be null. See class comment for usage.
//
// type_ is the type for totaling; comment is a comment about this tracer.
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
	if len(tracer_extraTracingStatistics) != 0 {
		size := len(tracer_extraTracingStatistics)
		t.extraTracingValues = make([]int64, size)
		i := 0
		for _, tracingStatistic := range tracer_extraTracingStatistics {
			t.extraTracingValues[i] = tracingStatistic.Start(t.startThread)
			i++
		}
	}

	trace := Tracer_GetThreadTrace()

	// Do nothing if the current thread trace wasn't initialized.
	if !trace.IsInitialized() {
		return t
	}

	// Check if we are creating too many Tracers.
	if trace.events.Size() >= Tracer_MAX_TRACE_SIZE {
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
	sb := jutil.NewStringBuilder()
	Tracer_AppendSpaces(sb, digitsColumnWidth-digitWidth)
	sb.AppendInt(v)
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
func Tracer_AppendSpaces(sb *jutil.StringBuilder, numSpaces int) {
	if numSpaces > 16 {
		Tracer_logger.Warning("Tracer.appendSpaces called with large numSpaces")
		// Avoid long loop in case some bug in the caller
		numSpaces = 16
	}
	for numSpaces >= 5 {
		sb.Append("     ")
		numSpaces -= 5
	}

	// We know it's less than 5 now
	switch numSpaces {
	case 1:
		sb.Append(" ")
	case 2:
		sb.Append("  ")
	case 3:
		sb.Append("   ")
	case 4:
		sb.Append("    ")
	}
}

// Tracer_AddTracingStatistic adds a new tracing statistic to a trace.
//
// Returns the index of this statistic (for use with stat.extraInfo()), or
// -1 if the statistic is not enabled.
//
// Java: Tracer.addTracingStatistic(TracingStatistic)
func Tracer_AddTracingStatistic(tracingStatistic Tracer_TracingStatistic) int {
	// Check to see if we can enable the tracing statistic before actually
	// adding it.
	if tracingStatistic.Enable() {
		// No synchronization needed, since this is a copy-on-write array.
		tracer_extraTracingStatistics = append(tracer_extraTracingStatistics, tracingStatistic)
		// 99.9% of the time, this will be O(1) and return
		// extraTracingStatistics.length - 1
		for i := len(tracer_extraTracingStatistics) - 1; i >= 0; i-- {
			if jutil.ObjectsEquals(tracer_extraTracingStatistics[i], tracingStatistic) {
				return i
			}
		}
		return -1
	} else {
		return -1
	}
}

// Tracer_ClearTracingStatisticsTestingOnly is for testing purposes only.
// These removes all current tracers. Severe errors can occur if there are
// any active tracers going on when this is called.
//
// The test suite uses this to remove any tracers that it has added.
//
// Java: Tracer.clearTracingStatisticsTestingOnly()
func Tracer_ClearTracingStatisticsTestingOnly() {
	tracer_extraTracingStatistics = nil
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
	if t.extraTracingValues != nil {
		// We use extraTracingValues.length rather than
		// extraTracingStatistics.size() because a new statistic may have
		// been added
		for i := 0; i < len(t.extraTracingValues); i++ {
			value := tracer_extraTracingStatistics[i].Stop(t.startThread)
			t.extraTracingValues[i] = value - t.extraTracingValues[i]
		}
	}

	// Do nothing if the thread trace was not initialized.
	if !trace.IsInitialized() {
		return 0
	}

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
	} else {
		return "[" + *t.type_ + "] " + t.comment
	}
}

// Java: Tracer.setDefaultSilenceThreshold(int)
func Tracer_SetDefaultSilenceThreshold(threshold int) {
	Tracer_GetThreadTrace().defaultSilenceThreshold = threshold
}

// Tracer_InitCurrentThreadTrace initializes the trace associated with the
// current thread by clearing out any existing trace. There shouldn't be a
// trace so if one is found we log it as an error.
//
// Java: Tracer.initCurrentThreadTrace()
func Tracer_InitCurrentThreadTrace() {
	events := Tracer_GetThreadTrace()
	if !events.IsEmpty() {
		Tracer_logger.LogWithThrown(jutil.Level_WARNING,
			"Non-empty timer log:\n"+events.String(),
			newTracer_Throwable(nil))
		Tracer_ClearThreadTrace()

		// Grab a new thread trace if we find a previous non-empty
		// ThreadTrace.
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

// Tracer_GetCurrentThreadTraceReport returns a timer report similar to the
// one described in the class comment.
//
// Java: Tracer.getCurrentThreadTraceReport()
func Tracer_GetCurrentThreadTraceReport() string {
	return Tracer_GetThreadTrace().String()
}

// Tracer_LogCurrentThreadTrace logs a timer report similar to the one
// described in the class comment.
//
// Java: Tracer.logCurrentThreadTrace()
func Tracer_LogCurrentThreadTrace() {
	trace := Tracer_GetThreadTrace()

	// New threads must call Tracer.initCurrentThreadTrace() before Tracer
	// statistics are gathered. This is a recent change (Jun 2007) that
	// prevents spurious Third Eye messages when an application uses a class
	// in a different package that happens to call Tracer without knowledge
	// of the application authors.
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

// Tracer_SetPrettyPrint sets whether pretty printing is enabled. See
// class-level comment. This only affects tracers created after this is
// called.
//
// Java: Tracer.setPrettyPrint(boolean)
func Tracer_SetPrettyPrint(enabled bool) {
	tracer_defaultPrettyPrint = enabled
}

// Tracer_Stat holds the statistics for a given tracer type.
type Tracer_Stat struct {
	count     int
	silent    int
	clockTime int
	extraInfo []int
}

// GetCount is the total count of tracers of a type, including silent.
//
// Java: Tracer.Stat.getCount()
func (s *Tracer_Stat) GetCount() int { return s.count }

// GetSilentCount is the total count of silent tracers of a type.
//
// Java: Tracer.Stat.getSilentCount()
func (s *Tracer_Stat) GetSilentCount() int { return s.silent }

// GetTotalTime is the total time spent in tracers of a type, in ms.
//
// Java: Tracer.Stat.getTotalTime()
func (s *Tracer_Stat) GetTotalTime() int { return s.clockTime }

// GetExtraInfo is the total time spent doing additional things that we are
// clocking.
//
// Java: Tracer.Stat.getExtraInfo(int)
func (s *Tracer_Stat) GetExtraInfo(index int) int {
	if s.extraInfo == nil {
		panic(jutil.NewNullPointerException())
	}
	if index >= len(s.extraInfo) {
		return 0
	}
	return s.extraInfo[index]
}

// This map tracks counts of tracers for each type over all time.
var tracer_typeToCountMap *Tracer_AtomicTracerStatMap

// This map tracks counts of silent tracers for each type over all time.
var tracer_typeToSilentMap *Tracer_AtomicTracerStatMap

// This map tracks time (ms) for each type over all time.
var tracer_typeToTimeMap *Tracer_AtomicTracerStatMap

// Tracer_EnableTypeMaps MUST be called before getTypeToCountMap (and
// friends) will return a valid map. This is because computing this
// information imposes a synchronization penalty on every Tracer that is
// stopped.
//
// Java: Tracer.enableTypeMaps()
func Tracer_EnableTypeMaps() {
	if tracer_typeToCountMap == nil {
		tracer_typeToCountMap = NewTracer_AtomicTracerStatMap()
		tracer_typeToSilentMap = NewTracer_AtomicTracerStatMap()
		tracer_typeToTimeMap = NewTracer_AtomicTracerStatMap()
	}
}

// Tracer_GetTypeToCountMap is used for exporting this data via varz.
// Accesses to this map must be synchronized on the map. If enableTypeMaps
// has not been called, this will return null.
//
// Java: Tracer.getTypeToCountMap()
func Tracer_GetTypeToCountMap() *jutil.LinkedHashMap[string, int64] {
	if tracer_typeToCountMap != nil {
		return tracer_typeToCountMap.GetMap()
	}
	return nil
}

// Tracer_GetTypeToSilentMap is used for exporting this data via varz.
// Accesses to this map must be synchronized on the map. If enableTypeMaps
// has not been called, this will return null.
//
// Java: Tracer.getTypeToSilentMap()
func Tracer_GetTypeToSilentMap() *jutil.LinkedHashMap[string, int64] {
	if tracer_typeToSilentMap != nil {
		return tracer_typeToSilentMap.GetMap()
	}
	return nil
}

// Tracer_GetTypeToTimeMap is used for exporting this data via varz.
// Accesses to this map must be synchronized on the map. If enableTypeMaps
// has not been called, this will return null.
//
// Java: Tracer.getTypeToTimeMap()
func Tracer_GetTypeToTimeMap() *jutil.LinkedHashMap[string, int64] {
	if tracer_typeToTimeMap != nil {
		return tracer_typeToTimeMap.GetMap()
	}
	return nil
}

// Tracer_GetStatsForType gets the Stat for a tracer type; never returns
// null.
//
// Java: Tracer.getStatsForType(String)
func Tracer_GetStatsForType(type_ string) *Tracer_Stat {
	stat := Tracer_GetThreadTrace().stats.Get(type_)
	if stat != nil {
		return stat
	}
	return tracer_ZERO_STAT
}

var tracer_ZERO_STAT = &Tracer_Stat{}

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
// prevEventTime is the time of the previous event which appears at the
// left most part of the trace line; indent is the indentation to put
// before the tracer to show the hierarchy; digitsColWidth is how many
// characters the digits should use.
//
// Java: Tracer.Event.toString(long, String, int)
func (e *tracer_Event) ToString(prevEventTime int64, indent string, digitsColWidth int) string {
	sb := jutil.NewStringBuilder()

	if prevEventTime == -1 {
		Tracer_AppendSpaces(sb, digitsColWidth)
	} else {
		sb.Append(tracer_longToPaddedString(
			e.EventTime()-prevEventTime, digitsColWidth))
	}

	sb.AppendChar(' ')
	sb.Append(tracer_formatTime(e.EventTime()))
	if e.isStart {
		sb.Append(" Start ")
		Tracer_AppendSpaces(sb, digitsColWidth)
		sb.Append("   ")
	} else {
		sb.Append(" Done ")
		delta := e.tracer.stopTimeMs - e.tracer.startTimeMs
		sb.Append(tracer_longToPaddedString(delta, digitsColWidth))
		sb.Append(" ms ")
		if e.tracer.extraTracingValues != nil {
			for i := 0; i < len(e.tracer.extraTracingValues); i++ {
				delta = e.tracer.extraTracingValues[i]
				sb.Append(jutil.StringFormat("%4d", delta))
				sb.Append(tracer_extraTracingStatistics[i].GetUnits())
				sb.Append(";  ")
			}
		}
	}
	sb.Append(indent)
	sb.Append(e.tracer.String())
	return sb.String()
}

// Tracer_ThreadTrace stores a thread's Trace.
type Tracer_ThreadTrace struct {
	// Events taking less than this number of milliseconds are not reported.
	defaultSilenceThreshold int // non-final

	// The Events corresponding to each startEvent/stopEvent
	events *jutil.ArrayList[*tracer_Event]

	// Tracers that have not had their .stop() called
	outstandingEvents *jutil.LinkedHashSet[*Tracer]

	// Map from type to Stat object
	stats *jutil.LinkedHashMap[string, *Tracer_Stat]

	// True if outstandingEvents has been cleared because we exceeded the
	// max trace limit.
	isOutstandingEventsTruncated bool

	// True if events has been cleared because we exceeded the max trace
	// limit.
	isEventsTruncated bool

	// Set to true if Tracer#initCurrentThreadTrace() was called by the
	// current thread.
	isInitialized bool

	// Whether pretty printing is enabled for the trace.
	prettyPrint bool
}

// NewTracer_ThreadTrace is new ThreadTrace().
func NewTracer_ThreadTrace() *Tracer_ThreadTrace {
	return &Tracer_ThreadTrace{
		events:            jutil.NewArrayList[*tracer_Event](),
		outstandingEvents: jutil.NewLinkedHashSet[*Tracer](),
		stats:             jutil.NewLinkedHashMap[string, *Tracer_Stat](),
	}
}

// Init initializes the trace.
//
// Java: Tracer.ThreadTrace.init()
func (tt *Tracer_ThreadTrace) Init() {
	tt.isInitialized = true
}

// IsInitialized: is initialized?
//
// Java: Tracer.ThreadTrace.isInitialized()
func (tt *Tracer_ThreadTrace) IsInitialized() bool {
	return tt.isInitialized
}

// StartEvent is called by the constructor Tracer#Tracer(String, String) to
// create a start event.
//
// Java: Tracer.ThreadTrace.startEvent(Tracer)
func (tt *Tracer_ThreadTrace) StartEvent(t *Tracer) {
	tt.events.Add(newTracer_Event(true, t))
	notAlreadyOutstanding := tt.outstandingEvents.Add(t)
	jutil.CheckState(notAlreadyOutstanding)
}

// EndEvent is called by Tracer#stop() to create a stop event.
//
// Java: Tracer.ThreadTrace.endEvent(Tracer, int)
func (tt *Tracer_ThreadTrace) EndEvent(t *Tracer, silenceThreshold int) {
	wasOutstanding := tt.outstandingEvents.Remove(t)
	if !wasOutstanding {
		if tt.isOutstandingEventsTruncated {
			// The events stack overflowed and was truncated, so just log a
			// warning. Otherwise, we get an exception which is extremely
			// confusing.
			Tracer_logger.LogWithThrown(jutil.Level_WARNING,
				"event not found, probably because the event stack "+
					"overflowed and was truncated",
				newTracer_Throwable(nil))
		} else {
			// throw an exception if the event was not found and the events
			// stack is pristine
			panic(jutil.NewIllegalStateException())
		}
	}

	elapsed := t.stopTimeMs - t.startTimeMs

	if silenceThreshold == -1 { // use default
		silenceThreshold = tt.defaultSilenceThreshold
	}

	if elapsed < int64(silenceThreshold) {
		// If this one is silent then we need to remove the start Event
		removed := false
		for i := 0; i < tt.events.Size(); i++ {
			e := tt.events.Get(i)
			if e.tracer == t {
				jutil.CheckState(e.isStart)
				tt.events.RemoveAt(i)
				removed = true
				break
			}
		}

		// Only assert if we didn't find the original and the events weren't
		// truncated.
		jutil.CheckState(removed || tt.isEventsTruncated)
	} else {
		tt.events.Add(newTracer_Event(false, t))
	}

	if t.type_ != nil {
		type_ := *t.type_
		stat := tt.stats.Get(type_)
		if stat == nil {
			stat = &Tracer_Stat{}
			if len(tracer_extraTracingStatistics) != 0 {
				stat.extraInfo = make([]int, len(tracer_extraTracingStatistics))
			}
			tt.stats.Put(type_, stat)
		}

		stat.count++
		if tracer_typeToCountMap != nil {
			tracer_typeToCountMap.IncrementBy(type_, 1)
		}

		// Java: stat.clockTime += elapsed (a compound assignment narrowing
		// long to int).
		stat.clockTime = int(int32(int64(stat.clockTime) + elapsed))
		if tracer_typeToTimeMap != nil {
			tracer_typeToTimeMap.IncrementBy(type_, elapsed)
		}

		if stat.extraInfo != nil && t.extraTracingValues != nil {
			overlapLength := min(len(stat.extraInfo), len(t.extraTracingValues))
			for i := 0; i < overlapLength; i++ {
				stat.extraInfo[i] = int(int32(int64(stat.extraInfo[i]) + t.extraTracingValues[i]))
				m := tracer_extraTracingStatistics[i].GetTracingStat()
				if m != nil {
					m.IncrementBy(type_, t.extraTracingValues[i])
				}
			}
		}

		if elapsed < int64(silenceThreshold) {
			stat.silent++
			if tracer_typeToSilentMap != nil {
				tracer_typeToSilentMap.IncrementBy(type_, 1)
			}
		}
	}
}

// Java: Tracer.ThreadTrace.isEmpty()
func (tt *Tracer_ThreadTrace) IsEmpty() bool {
	return tt.events.IsEmpty() && tt.outstandingEvents.IsEmpty()
}

// Java: Tracer.ThreadTrace.truncateOutstandingEvents()
func (tt *Tracer_ThreadTrace) TruncateOutstandingEvents() {
	tt.isOutstandingEventsTruncated = true
	tt.outstandingEvents.Clear()
}

// Java: Tracer.ThreadTrace.truncateEvents()
func (tt *Tracer_ThreadTrace) TruncateEvents() {
	tt.isEventsTruncated = true
	tt.events.Clear()
}

// String produces the lovely Trace seen in the class comments.
//
// Java: Tracer.ThreadTrace.toString()
func (tt *Tracer_ThreadTrace) String() string {
	numDigits := tt.getMaxDigits()
	sb := jutil.NewStringBuilder()
	var etime int64 = -1
	// Java: an ArrayDeque used as a stack (push/pop at the front).
	var indent []string
	for _, e := range tt.events.Slice() {
		if tt.prettyPrint && !e.isStart && len(indent) != 0 {
			indent = indent[1:]
		}
		sb.Append(" ")
		if tt.prettyPrint {
			sb.Append(e.ToString(etime, strings.Join(indent, ""), numDigits))
		} else {
			sb.Append(e.ToString(etime, "", 4))
		}
		etime = e.EventTime()
		sb.AppendChar('\n')
		if tt.prettyPrint && e.isStart {
			indent = append([]string{"|  "}, indent...)
		}
	}

	if tt.outstandingEvents.Size() != 0 {
		now := Tracer_clock.CurrentTimeMillis()

		sb.Append(" Unstopped timers:\n")
		for t := range tt.outstandingEvents.All() {
			sb.Append("  ").
				Append(t.String()).
				Append(" (").
				AppendInt(now - t.startTimeMs).
				Append(" ms, started at ").
				Append(tracer_formatTime(t.startTimeMs)).
				Append(")\n")
		}
	}

	for key, stat := range tt.stats.All() {
		if stat.count > 1 {
			sb.Append(" TOTAL ").
				Append(key).
				Append(" ").
				AppendInt(int64(stat.count)).
				Append(" (").
				AppendInt(int64(stat.clockTime)).
				Append(" ms")
			if stat.extraInfo != nil {
				for i := 0; i < len(stat.extraInfo); i++ {
					sb.Append("; ")
					sb.AppendInt(int64(stat.extraInfo[i])).
						AppendChar(' ').
						Append(tracer_extraTracingStatistics[i].GetUnits())
				}
			}
			sb.Append(")\n")
		}
	}
	return sb.String()
}

// getMaxDigits gets the maximum number of digits that can appear in the
// tracer output in the gaps between tracers or the duration of a tracer.
// This is used by the pretty printing case so that all of the tracers are
// aligned.
//
// Java: Tracer.ThreadTrace.getMaxDigits()
func (tt *Tracer_ThreadTrace) getMaxDigits() int {
	var etime int64 = -1
	var maxTime int64 = 0
	for _, e := range tt.events.Slice() {
		if etime != -1 {
			t := e.EventTime() - etime
			maxTime = max(maxTime, t)
		}
		if !e.isStart {
			t := e.tracer.stopTimeMs - e.tracer.startTimeMs
			maxTime = max(maxTime, t)
		}
		etime = e.EventTime()
	}
	// Minimum is 3 to preserve an indent even when max is small.
	return max(3, tracer_numDigits(maxTime))
}

// Holds the ThreadTrace for each thread.
//
// Go-only: language-adapter — a ThreadLocal; the port has one thread.
var tracer_traces *Tracer_ThreadTrace

// Tracer_GetThreadTrace gets the ThreadTrace for the current thread,
// creating one if necessary.
//
// Java: Tracer.getThreadTrace()
func Tracer_GetThreadTrace() *Tracer_ThreadTrace {
	t := tracer_traces
	if t == nil {
		t = NewTracer_ThreadTrace()
		t.prettyPrint = tracer_defaultPrettyPrint
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

// Tracer_TracingStatistic allows the program to add additional optional
// statistics to the trace output.
//
// The class com.google.monitoring.tracing.TracingStatistics contains
// several useful tracing statistics
type Tracer_TracingStatistic interface {
	// Start is called at the start of the trace. It should return a numeric
	// result indicating the amount of the specific resource in use before
	// the call started.
	//
	// Java: Tracer.TracingStatistic.start(Thread)
	Start(thread *jutil.Thread) int64

	// Stop is called at the end of the trace. It should return a numeric
	// result indicating the amount of the specific resource in use after
	// the call ends. The actual reported result will be the result end() -
	// start().
	//
	// Java: Tracer.TracingStatistic.stop(Thread)
	Stop(thread *jutil.Thread) int64

	// Enable is called when this tracing statistic is first enabled. A
	// return value of True indicates that this statistic can successfully
	// run in the current JVM.
	//
	// Java: Tracer.TracingStatistic.enable()
	Enable() bool

	// GetTracingStat returns this tracing statistic's trace map.
	//
	// Java: Tracer.TracingStatistic.getTracingStat()
	GetTracingStat() *Tracer_AtomicTracerStatMap

	// GetUnits is a string that should be appended to the numeric output
	// indicating what this is.
	//
	// Java: Tracer.TracingStatistic.getUnits()
	GetUnits() string
}

// Tracer_AtomicTracerStatMap encapsulates a map for keeping track of
// tracing statistics. It allows the caller to atomically increment named
// fields.
//
// Go-only: language-adapter — Java ConcurrentHashMap order: the map
// iterates in insertion order here; Java's ConcurrentHashMap iterates in
// hash order. Nothing in the compiler iterates it.
type Tracer_AtomicTracerStatMap struct {
	m *jutil.LinkedHashMap[string, int64]
}

// NewTracer_AtomicTracerStatMap is new AtomicTracerStatMap().
func NewTracer_AtomicTracerStatMap() *Tracer_AtomicTracerStatMap {
	return &Tracer_AtomicTracerStatMap{m: jutil.NewLinkedHashMap[string, int64]()}
}

// IncrementBy atomically increments the specified field by the specified
// amount.
//
// Java: Tracer.AtomicTracerStatMap.incrementBy(String, long)
func (a *Tracer_AtomicTracerStatMap) IncrementBy(key string, delta int64) {
	// Go-only: language-adapter — sequential, so the compare-and-set loop
	// of the Java code is a plain read-modify-write.
	oldValue, ok := a.m.Lookup(key)
	if !ok {
		a.m.Put(key, delta)
		return
	}
	a.m.Put(key, oldValue+delta)
}

// GetMap returns a map of key:value pairs.
//
// Java: Tracer.AtomicTracerStatMap.getMap()
func (a *Tracer_AtomicTracerStatMap) GetMap() *jutil.LinkedHashMap[string, int64] {
	return a.m
}

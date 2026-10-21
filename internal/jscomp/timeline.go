// Port of com/google/javascript/jscomp/Timeline.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Timeline is an ordered set that moves values to the front when added
// (even if already contained) and which allows for the marking of named
// points in time in the series from which point an accumulation of all
// following values can be (efficiently) requested.
//
// TODO(stalcup): save memory by unlinking events older than the oldest
// marked time.
type Timeline[T comparable] struct {
	eventsByTime  map[string]*timeline_Event[T]
	eventsByValue map[T]*timeline_Event[T]
	headEvent     *timeline_Event[T]
	freeEvent     *timeline_Event[T]
	eventSlab     []timeline_Event[T]
}

// timeline_Event is Timeline.Event<T>.
type timeline_Event[T comparable] struct {
	nextEvent     *timeline_Event[T]
	previousEvent *timeline_Event[T]
	value         T
	isTime        bool
}

func (tl *Timeline[T]) allocEvent(value T, isTime bool) *timeline_Event[T] {
	if e := tl.freeEvent; e != nil {
		tl.freeEvent = e.nextEvent
		*e = timeline_Event[T]{value: value, isTime: isTime}
		return e
	}
	if len(tl.eventSlab) == 0 {
		tl.eventSlab = make([]timeline_Event[T], 256)
	}
	e := &tl.eventSlab[0]
	tl.eventSlab = tl.eventSlab[1:]
	e.value = value
	e.isTime = isTime
	return e
}

// NewTimeline is new Timeline<>().
func NewTimeline[T comparable]() *Timeline[T] {
	tl := &Timeline[T]{
		eventsByTime:  make(map[string]*timeline_Event[T]),
		eventsByValue: make(map[T]*timeline_Event[T]),
	}
	var zero T
	tl.headEvent = tl.allocEvent(zero, true)
	return tl
}

// Java: Timeline.add(T)
func (tl *Timeline[T]) Add(value T) {
	jutil.CheckNotNull(value)
	event := tl.eventsByValue[value]
	if tl.headEvent == event {
		return
	}
	if event != nil {
		event.previousEvent.nextEvent = event.nextEvent
		event.nextEvent.previousEvent = event.previousEvent
		event.nextEvent = nil
	} else {
		event = tl.allocEvent(value, false)
		tl.eventsByValue[value] = event
	}
	event.previousEvent = tl.headEvent
	tl.headEvent.nextEvent = event
	tl.headEvent = event
}

// Java: Timeline.remove(T)
func (tl *Timeline[T]) Remove(value T) {
	event := tl.eventsByValue[value]
	if event != nil {
		delete(tl.eventsByValue, value)
		if event.nextEvent != nil {
			event.nextEvent.previousEvent = event.previousEvent
		} else {
			tl.headEvent = event.previousEvent
		}
		event.previousEvent.nextEvent = event.nextEvent

		var zero T
		event.value = zero
		event.previousEvent = nil
		event.nextEvent = tl.freeEvent
		tl.freeEvent = event
	}
}

// Java: Timeline.mark(String)
func (tl *Timeline[T]) Mark(timeName string) {
	event := tl.eventsByTime[timeName]
	if tl.headEvent == event {
		return
	}
	if event != nil {
		event.previousEvent.nextEvent = event.nextEvent
		event.nextEvent.previousEvent = event.previousEvent
		event.nextEvent = nil
	} else {
		var zero T
		event = tl.allocEvent(zero, true)
		tl.eventsByTime[timeName] = event
	}
	event.previousEvent = tl.headEvent
	tl.headEvent.nextEvent = event
	tl.headEvent = event
}

// GetSince returns the values added since the named time, or nil (Java
// null) if the time was never marked.
//
// Java: Timeline.getSince(String)
func (tl *Timeline[T]) GetSince(timeName string) []T {
	firstEvent := tl.eventsByTime[timeName]
	if firstEvent == nil {
		return nil
	}
	values := []T{}
	for event := firstEvent; event != nil; event = event.nextEvent {
		if !event.isTime {
			values = append(values, event.value)
		}
	}
	return values
}

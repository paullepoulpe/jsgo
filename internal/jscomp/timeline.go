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
	eventsByTime  *jutil.LinkedHashMap[timeline_Time, *timeline_Event]
	eventsByValue *jutil.LinkedHashMap[T, *timeline_Event]

	// In practice headEvent is Event<Time> or Event<T>
	headEvent *timeline_Event
}

// timeline_Event is Timeline.Event<T>.
//
// Go-only: language-adapter — Event<?> holds either a T or a Time, so the
// value is stored as `any`.
type timeline_Event struct {
	nextEvent     *timeline_Event
	previousEvent *timeline_Event
	value         any
}

// Java: Timeline.Event(T)
func newTimeline_Event(value any) *timeline_Event {
	jutil.CheckNotNull(value)
	return &timeline_Event{value: value}
}

// timeline_Time is Timeline.Time; its equals/hashCode compare the name, so
// it is a Go value type keyed by that name.
type timeline_Time struct {
	name string
}

// Java: Timeline.Time(String)
func newTimeline_Time(name string) timeline_Time {
	return timeline_Time{name: name}
}

// NewTimeline is new Timeline<>().
func NewTimeline[T comparable]() *Timeline[T] {
	return &Timeline[T]{
		eventsByTime:  jutil.NewLinkedHashMap[timeline_Time, *timeline_Event](),
		eventsByValue: jutil.NewLinkedHashMap[T, *timeline_Event](),
		headEvent:     newTimeline_Event(newTimeline_Time("-beginning-")),
	}
}

// Java: Timeline.add(T)
func (tl *Timeline[T]) Add(value T) {
	timeline_addEvent(tl, value, tl.eventsByValue)
}

// Java: Timeline.remove(T)
func (tl *Timeline[T]) Remove(value T) {
	event := tl.eventsByValue.Remove(value)

	// If the event already exists somewhere else in the history then...
	if event != nil {
		// make the rest of the list not reference it...
		if event.nextEvent != nil {
			event.nextEvent.previousEvent = event.previousEvent
		} else {
			// if it was the head element then back the head reference up one
			// node.
			tl.headEvent = event.previousEvent
		}
		event.previousEvent.nextEvent = event.nextEvent

		// make it not reference the rest of the list.
		event.nextEvent = nil
		event.previousEvent = nil
	}
}

// Java: Timeline.mark(String)
func (tl *Timeline[T]) Mark(timeName string) {
	timeline_addEvent(tl, newTimeline_Time(timeName), tl.eventsByTime)
}

// GetSince returns the values added since the named time, or nil (Java
// null) if the time was never marked.
//
// Java: Timeline.getSince(String)
func (tl *Timeline[T]) GetSince(timeName string) []T {
	values := []T{}
	firstEvent := tl.eventsByTime.Get(newTimeline_Time(timeName))
	if firstEvent == nil {
		return nil
	}

	for event := firstEvent; event != nil; event = event.nextEvent {
		// If the event contains a user provided value, accumulate it.
		if _, isTime := event.value.(timeline_Time); !isTime {
			values = append(values, event.value.(T))
		}
	}
	return values
}

// In practice V is either T or Time.
//
// Java: Timeline.addEvent(V, Map<V, Event<V>>)
func timeline_addEvent[T comparable, V comparable](tl *Timeline[T], value V, eventsByKey *jutil.LinkedHashMap[V, *timeline_Event]) {
	event := eventsByKey.Get(value)

	// If the event already exists and is already at the front, do nothing.
	if tl.headEvent == event {
		return
	}

	// If the event already exists somewhere else in the history then...
	if event != nil {
		// cut it out of the linked list...
		event.previousEvent.nextEvent = event.nextEvent
		event.nextEvent.previousEvent = event.previousEvent
		event.nextEvent = nil
	} else {
		// Otherwise create and track an event for the given value.
		event = newTimeline_Event(value)
		eventsByKey.Put(value, event)
	}

	// Regardless, stick the event at the front.
	event.previousEvent = tl.headEvent
	tl.headEvent.nextEvent = event
	tl.headEvent = event
}

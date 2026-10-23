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

const (
	timelineChunkShift = 8
	timelineChunkSize  = 1 << timelineChunkShift
	timelineChunkMask  = timelineChunkSize - 1
	timelineTimeBit    = uint32(1) << 31
	timelinePrevMask   = ^timelineTimeBit
)

// Timeline is an ordered set that moves values to the front when added
// (even if already contained) and which allows for the marking of named
// points in time in the series from which point an accumulation of all
// following values can be (efficiently) requested.
//
// TODO(stalcup): save memory by unlinking events older than the oldest
// marked time.
type Timeline[T comparable] struct {
	eventsByTime  map[string]uint32
	eventsByValue map[T]uint32
	chunks        []*[timelineChunkSize]timeline_Event[T]
	nextIdx       uint32
	headEvent     uint32
	freeEvent     uint32
}

// timeline_Event is Timeline.Event<T>.
type timeline_Event[T comparable] struct {
	value       T
	next        uint32
	prevAndTime uint32
}

func (e *timeline_Event[T]) prev() uint32   { return e.prevAndTime & timelinePrevMask }
func (e *timeline_Event[T]) isTime() bool   { return e.prevAndTime&timelineTimeBit != 0 }
func (e *timeline_Event[T]) setPrev(p uint32) {
	e.prevAndTime = (e.prevAndTime & timelineTimeBit) | p
}

func (tl *Timeline[T]) event(id uint32) *timeline_Event[T] {
	idx := id - 1
	return &tl.chunks[idx>>timelineChunkShift][idx&timelineChunkMask]
}

func (tl *Timeline[T]) allocEvent(value T, isTime bool) uint32 {
	var pt uint32
	if isTime {
		pt = timelineTimeBit
	}
	if id := tl.freeEvent; id != 0 {
		e := tl.event(id)
		tl.freeEvent = e.next
		*e = timeline_Event[T]{value: value, prevAndTime: pt}
		return id
	}
	idx := tl.nextIdx
	tl.nextIdx = idx + 1
	if idx&timelineChunkMask == 0 {
		tl.chunks = append(tl.chunks, new([timelineChunkSize]timeline_Event[T]))
	}
	e := &tl.chunks[idx>>timelineChunkShift][idx&timelineChunkMask]
	e.value = value
	e.prevAndTime = pt
	return idx + 1
}

// NewTimeline is new Timeline<>().
func NewTimeline[T comparable]() *Timeline[T] {
	tl := &Timeline[T]{
		eventsByTime:  make(map[string]uint32),
		eventsByValue: make(map[T]uint32),
	}
	var zero T
	tl.headEvent = tl.allocEvent(zero, true)
	return tl
}

// Java: Timeline.add(T)
func (tl *Timeline[T]) Add(value T) {
	jutil.CheckNotNull(value)
	id := tl.eventsByValue[value]
	if tl.headEvent == id {
		return
	}
	var e *timeline_Event[T]
	if id != 0 {
		e = tl.event(id)
		prevId := e.prev()
		nextId := e.next
		tl.event(prevId).next = nextId
		tl.event(nextId).setPrev(prevId)
		e.next = 0
	} else {
		id = tl.allocEvent(value, false)
		tl.eventsByValue[value] = id
		e = tl.event(id)
	}
	e.setPrev(tl.headEvent)
	tl.event(tl.headEvent).next = id
	tl.headEvent = id
}

// Java: Timeline.remove(T)
func (tl *Timeline[T]) Remove(value T) {
	id := tl.eventsByValue[value]
	if id != 0 {
		delete(tl.eventsByValue, value)
		e := tl.event(id)
		prevId := e.prev()
		nextId := e.next
		if nextId != 0 {
			tl.event(nextId).setPrev(prevId)
		} else {
			tl.headEvent = prevId
		}
		tl.event(prevId).next = nextId

		var zero T
		e.value = zero
		e.prevAndTime = 0
		e.next = tl.freeEvent
		tl.freeEvent = id
	}
}

// Java: Timeline.mark(String)
func (tl *Timeline[T]) Mark(timeName string) {
	id := tl.eventsByTime[timeName]
	if tl.headEvent == id {
		return
	}
	var e *timeline_Event[T]
	if id != 0 {
		e = tl.event(id)
		prevId := e.prev()
		nextId := e.next
		tl.event(prevId).next = nextId
		tl.event(nextId).setPrev(prevId)
		e.next = 0
	} else {
		var zero T
		id = tl.allocEvent(zero, true)
		tl.eventsByTime[timeName] = id
		e = tl.event(id)
	}
	e.setPrev(tl.headEvent)
	tl.event(tl.headEvent).next = id
	tl.headEvent = id
}

// Unmark removes a named time point from the timeline if present.
func (tl *Timeline[T]) Unmark(timeName string) {
	id := tl.eventsByTime[timeName]
	if id == 0 {
		return
	}
	delete(tl.eventsByTime, timeName)
	e := tl.event(id)
	prevId := e.prev()
	nextId := e.next
	if nextId != 0 {
		tl.event(nextId).setPrev(prevId)
	} else {
		tl.headEvent = prevId
	}
	tl.event(prevId).next = nextId

	var zero T
	e.value = zero
	e.prevAndTime = 0
	e.next = tl.freeEvent
	tl.freeEvent = id
}

// Values calls yield for every value in the timeline, oldest first.
//
// Go-only: lets ConvertTypesToColors reach the detached scope roots the
// timeline still references.
func (tl *Timeline[T]) Values(yield func(T) bool) {
	for id := uint32(1); id != 0; {
		e := tl.event(id)
		if !e.isTime() && !yield(e.value) {
			return
		}
		id = e.next
	}
}

// GetSince returns the values added since the named time, or nil (Java
// null) if the time was never marked.
//
// Java: Timeline.getSince(String)
func (tl *Timeline[T]) GetSince(timeName string) []T {
	firstId := tl.eventsByTime[timeName]
	if firstId == 0 {
		return nil
	}
	values := []T{}
	for id := firstId; id != 0; {
		e := tl.event(id)
		if !e.isTime() {
			values = append(values, e.value)
		}
		id = e.next
	}
	return values
}

// Port of com/google/javascript/jscomp/TimelineTest.java @ closure-compiler v20260415

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
	"testing"
)

// timelineTestEvent is TimelineTest.Event, whose equals/hashCode compare
// the hashcode field.
//
// Go-only: language-adapter — Timeline keys values with Go ==, so a value
// type gives Event's equals() semantics.
type timelineTestEvent struct {
	hashcode int
}

var (
	timelineTest_ATE_GREEN_EGGS = timelineTestEvent{1}
	timelineTest_ATE_HAM        = timelineTestEvent{2}
	timelineTest_ATE_HASHBROWNS = timelineTestEvent{3}
)

func TestTimeline_RemoveSimple(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	timeline.Mark("Monday")
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Remove(timelineTest_ATE_GREEN_EGGS)

	assertSliceEqual(t, "since Monday", timeline.GetSince("Monday"), []timelineTestEvent{})
}

func TestTimeline_RemoveMaintainsOrder(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	timeline.Mark("Monday")
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Add(timelineTest_ATE_HAM)
	timeline.Add(timelineTest_ATE_HASHBROWNS)
	timeline.Remove(timelineTest_ATE_HAM)

	assertSameElements(t, "since Monday", timeline.GetSince("Monday"),
		[]timelineTestEvent{timelineTest_ATE_GREEN_EGGS, timelineTest_ATE_HASHBROWNS})
}

func TestTimeline_UnknownTimesReturnNull(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	if got := timeline.GetSince("Monday"); got != nil {
		t.Errorf("got %v, want null", got)
	}
}

func TestTimeline_DoesntReturnValuesBeforeTime(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Mark("Monday")

	got := timeline.GetSince("Monday")
	if got == nil || len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}

func TestTimeline_ReturnsValuesAfterTime(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	timeline.Mark("Monday")
	timeline.Add(timelineTest_ATE_GREEN_EGGS)

	assertSameElements(t, "since Monday", timeline.GetSince("Monday"),
		[]timelineTestEvent{timelineTest_ATE_GREEN_EGGS})
}

func TestTimeline_MaintainsOrder1(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	timeline.Mark("Monday")
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Add(timelineTest_ATE_HAM)

	assertSliceEqual(t, "since Monday", timeline.GetSince("Monday"),
		[]timelineTestEvent{timelineTest_ATE_GREEN_EGGS, timelineTest_ATE_HAM})
}

func TestTimeline_MaintainsOrder2(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	timeline.Mark("Monday")
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Add(timelineTest_ATE_HAM)

	assertSliceEqual(t, "since Monday", timeline.GetSince("Monday"),
		[]timelineTestEvent{timelineTest_ATE_GREEN_EGGS, timelineTest_ATE_HAM})
}

func TestTimeline_MaintainsOrder3(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	timeline.Mark("Monday")
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Add(timelineTest_ATE_HAM)
	timeline.Add(timelineTest_ATE_GREEN_EGGS)

	assertSliceEqual(t, "since Monday", timeline.GetSince("Monday"),
		[]timelineTestEvent{timelineTest_ATE_HAM, timelineTest_ATE_GREEN_EGGS})
}

func TestTimeline_MaintainsOrder4(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	timeline.Mark("Monday")
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Add(timelineTest_ATE_HAM)
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Add(timelineTest_ATE_HAM)

	assertSliceEqual(t, "since Monday", timeline.GetSince("Monday"),
		[]timelineTestEvent{timelineTest_ATE_GREEN_EGGS, timelineTest_ATE_HAM})
}

func TestTimeline_UsesEqualityAndHashcode(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	// Are object and hashcode equal.
	ateBerries := timelineTestEvent{10}
	ateGrapes := timelineTestEvent{10}

	timeline.Mark("Monday")
	timeline.Add(ateBerries)
	timeline.Add(ateGrapes)

	assertSliceEqual(t, "since Monday", timeline.GetSince("Monday"), []timelineTestEvent{ateBerries})
	assertSliceEqual(t, "since Monday", timeline.GetSince("Monday"), []timelineTestEvent{ateGrapes})
}

func TestTimeline_UpdatesExistingTimes(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	timeline.Mark("Monday")
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Mark("Monday")

	assertSliceEqual(t, "since Monday", timeline.GetSince("Monday"), []timelineTestEvent{})
}

func TestTimeline_ManagesDifferentTimes(t *testing.T) {
	timeline := NewTimeline[timelineTestEvent]()
	timeline.Mark("Monday")
	timeline.Add(timelineTest_ATE_GREEN_EGGS)
	timeline.Mark("Thursday")
	timeline.Add(timelineTest_ATE_HAM)

	assertSliceEqual(t, "since Monday", timeline.GetSince("Monday"),
		[]timelineTestEvent{timelineTest_ATE_GREEN_EGGS, timelineTest_ATE_HAM})
	assertSliceEqual(t, "since Thursday", timeline.GetSince("Thursday"),
		[]timelineTestEvent{timelineTest_ATE_HAM})
}

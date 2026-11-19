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
// Port of com/google/javascript/jscomp/Timeline.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.Timeline
type Timeline struct {
	eventsByTime  j2gort.OrderedMap[string, *Timeline_Event] `gson:"eventsByTime"`
	eventsByValue *j2gort.Ext_java_util_LinkedHashMap        `gson:"eventsByValue"`
	headEvent     *Timeline_Event                            `gson:"headEvent"`
	idHash_       atomic.Int32
	eventsByNode_ j2gort.OrderedMap[pkg_rhino.Node, *Timeline_Event]
	eventBuf_     []Timeline_Event
	freeEvent_    *Timeline_Event
}

func (this_ *Timeline) BindSelf_(s any) {
	_ = s
}
func (this_ *Timeline) IsJavaNil_() bool { return this_ == nil }
func (*Timeline) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Timeline", "Timeline"
}
func (this_ *Timeline) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Timeline) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Timeline", "com.google.javascript.jscomp.Timeline", func(v_ any) bool {
		return j2gort.InstanceOf[*Timeline](v_)
	})
}

func (this_ *Timeline) initFields_Timeline() {
	this_.eventsByTime = j2gort.NewOrderedMap[string, *Timeline_Event]()
	this_.eventsByNode_ = j2gort.NewOrderedMap[pkg_rhino.Node, *Timeline_Event]()

	ev0_ := this_.allocEvent_()
	ev0_.kind_ = 1
	this_.headEvent = ev0_

}

// Java constructor: Timeline()
func NewTimeline() *Timeline {
	this_ := &Timeline{}
	InitTimeline(this_)
	return this_
}

func InitTimeline(this_ *Timeline) {
	this_.initFields_Timeline()
}

// Java method: Timeline.add(T value)
func (this_ *Timeline) Add(value any) {
	if n_, ok_ := value.(pkg_rhino.Node); ok_ {
		this_.AddNode_(n_)
		return
	}
	if this_ == nil || value == nil {
		j2gort.ThrowNPE()
	}
	if this_.eventsByValue == nil {
		this_.eventsByValue = j2gort.NewExt_java_util_LinkedHashMap()
	}
	this_.priv_Timeline_addEvent_5356(value, this_.eventsByValue)

}

// Java method: Timeline.remove(T value)
func (this_ *Timeline) Remove(value any) {
	if n_, ok_ := value.(pkg_rhino.Node); ok_ {
		this_.RemoveNode_(n_)
		return
	}
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	if value == nil || this_.eventsByValue == nil {
		return
	}
	event, _ := this_.eventsByValue.Remove_1(value).(*Timeline_Event)
	if event != nil {
		this_.unlinkAndFreeEvent_(event)
	}

}

// Java method: Timeline.mark(java.lang.String timeName)
func (this_ *Timeline) Mark(timeName string) {
	if this_ == nil || timeName == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	event := this_.eventsByTime.Get(timeName)
	if this_.headEvent == event {
		return
	}
	if event != nil {
		event.previousEvent.nextEvent = event.nextEvent
		event.nextEvent.previousEvent = event.previousEvent
		event.nextEvent = nil
	} else {
		event = this_.allocEvent_()
		event.kind_ = 1
		this_.eventsByTime.Put(timeName, event)
	}
	event.previousEvent = this_.headEvent
	this_.headEvent.nextEvent = event
	this_.headEvent = event

}

// Java method: Timeline.getSince(java.lang.String timeName)
func (this_ *Timeline) GetSince(timeName string) j2gort.Ext_java_util_List {
	if this_ == nil || timeName == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	firstEvent := this_.eventsByTime.Get(timeName)
	if firstEvent == nil {
		return nil
	}
	if this_.eventsByValue == nil {
		values := []pkg_rhino.Node{}
		for event := firstEvent; event != nil; event = event.nextEvent {
			if event.kind_ == 2 {
				values = append(values, event.node_)
			}
		}
		return j2gort.WrapTypedArrayList(values)
	}
	values := j2gort.NewExt_java_util_ArrayList()
	for event := firstEvent; event != nil; event = event.nextEvent {
		if event.kind_ == 2 {
			values.Add_1(j2gort.BoxNode(event.node_))
		} else if event.kind_ == 0 {
			values.Add_1(event.value)
		}
	}
	return values

}

// Java method: Timeline.addEvent(V value, java.util.Map<V,com.google.javascript.jscomp.Timeline.Event<V>> eventsByKey)
func (this_ *Timeline) priv_Timeline_addEvent_5356(value any, eventsByKey j2gort.Ext_java_util_Map) {
	event, _ := eventsByKey.Get(value).(*Timeline_Event)
	if this_.headEvent == event {
		return
	}
	if event != nil {
		event.previousEvent.nextEvent = event.nextEvent
		event.nextEvent.previousEvent = event.previousEvent
		event.nextEvent = nil
	} else {
		event = this_.allocEvent_()
		event.value = value
		event.kind_ = 0
		eventsByKey.Put_2(value, event)
	}
	event.previousEvent = this_.headEvent
	this_.headEvent.nextEvent = event
	this_.headEvent = event

}

// Java class: com.google.javascript.jscomp.Timeline.Event
type Timeline_Event struct {
	nextEvent     *Timeline_Event `gson:"nextEvent"`
	previousEvent *Timeline_Event `gson:"previousEvent"`
	value         any             `gson:"value"`
	idHash_       atomic.Int32
	node_         pkg_rhino.Node
	kind_         uint8
}

func (this_ *Timeline_Event) BindSelf_(s any) {
	_ = s
}
func (this_ *Timeline_Event) IsJavaNil_() bool { return this_ == nil }
func (*Timeline_Event) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Timeline$Event", "Event"
}
func (this_ *Timeline_Event) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Timeline_Event) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Timeline.Event", "com.google.javascript.jscomp.Timeline$Event", func(v_ any) bool {
		return j2gort.InstanceOf[*Timeline_Event](v_)
	})
}

func (this_ *Timeline_Event) initFields_Timeline_Event() {
	_ = this_
}

// Java constructor: Event(T value)
func NewTimeline_Event(value any) *Timeline_Event {
	this_ := &Timeline_Event{}
	InitTimeline_Event(this_, value)
	return this_
}

func InitTimeline_Event(this_ *Timeline_Event, value any) {
	_ = value
	this_.initFields_Timeline_Event()
	j2gort.CheckNotNull(value)
	this_.value = value
}

// Java class: com.google.javascript.jscomp.Timeline.Time
type Timeline_Time struct {
	name string `gson:"name"`
}

func (this_ *Timeline_Time) BindSelf_(s any) {
	_ = s
}
func (this_ *Timeline_Time) IsJavaNil_() bool { return this_ == nil }
func (*Timeline_Time) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Timeline$Time", "Time"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Timeline.Time", "com.google.javascript.jscomp.Timeline$Time", func(v_ any) bool {
		return j2gort.InstanceOf[*Timeline_Time](v_)
	})
}

func (this_ *Timeline_Time) initFields_Timeline_Time() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: Time(java.lang.String name)
func NewTimeline_Time(name string) *Timeline_Time {
	this_ := &Timeline_Time{}
	InitTimeline_Time(this_, name)
	return this_
}

func InitTimeline_Time(this_ *Timeline_Time, name string) {
	_ = name
	this_.initFields_Timeline_Time()
	j2gort.CheckNotNullString(name)
	this_.name = name
}

// Java method: Time.equals(java.lang.Object anObject)
func (this_ *Timeline_Time) Equals(arg_10 any) bool {
	anObject := arg_10
	_ = anObject
	var time *Timeline_Time
	_ = time
	var ok_11 bool
	time, ok_11 = j2gort.TryCast[*Timeline_Time](anObject)
	if ok_11 {
		{
			sr_12 := this_.name
			se_13 := time.name
			return j2gort.CheckNotNullString(sr_12) == se_13
		}
	}
	return false
}

// Java method: Time.hashCode()
func (this_ *Timeline_Time) HashCode() int32 {
	return javalang.StringHashCode(j2gort.CheckNotNullString(this_.name))
}
func (this_ *Timeline) allocEvent_() *Timeline_Event {
	if ev := this_.freeEvent_; ev != nil {
		this_.freeEvent_ = ev.nextEvent
		ev.nextEvent = nil
		return ev
	}
	if len(this_.eventBuf_) == cap(this_.eventBuf_) {
		this_.eventBuf_ = make([]Timeline_Event, 0, 64)
	}
	i := len(this_.eventBuf_)
	this_.eventBuf_ = this_.eventBuf_[:i+1]
	return &this_.eventBuf_[i]
}

func (this_ *Timeline) unlinkAndFreeEvent_(event *Timeline_Event) {
	if event.nextEvent != nil {
		event.nextEvent.previousEvent = event.previousEvent
	} else {
		this_.headEvent = event.previousEvent
	}
	event.previousEvent.nextEvent = event.nextEvent
	event.previousEvent = nil
	event.value = nil
	event.node_ = 0
	event.kind_ = 0
	event.nextEvent = this_.freeEvent_
	this_.freeEvent_ = event
}

func (this_ *Timeline) AddNode_(n pkg_rhino.Node) {
	if this_ == nil || n == 0 {
		j2gort.ThrowNPE()
	}
	event := this_.eventsByNode_.Get(n)
	if this_.headEvent == event {
		return
	}
	if event != nil {
		event.previousEvent.nextEvent = event.nextEvent
		event.nextEvent.previousEvent = event.previousEvent
		event.nextEvent = nil
	} else {
		event = this_.allocEvent_()
		event.node_ = n
		event.kind_ = 2
		this_.eventsByNode_.Put(n, event)
	}
	event.previousEvent = this_.headEvent
	this_.headEvent.nextEvent = event
	this_.headEvent = event
}

func (this_ *Timeline) RemoveNode_(n pkg_rhino.Node) {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	if n == 0 {
		return
	}
	event := this_.eventsByNode_.Remove(n)
	if event != nil {
		this_.unlinkAndFreeEvent_(event)
	}
}

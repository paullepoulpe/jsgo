/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/base/LinkedIdentityHashMap.java (generated mechanically by j2go).

package base

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.base.LinkedIdentityHashMap
type LinkedIdentityHashMap[K comparable, V any] struct {
	innerMap j2gort.IdentityOrderedMap[K, V] `gson:"innerMap"`

	idHash_ atomic.Int32
}

func (this_ *LinkedIdentityHashMap[K, V]) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedIdentityHashMap[K, V]) IsJavaNil_() bool { return this_ == nil }
func (*LinkedIdentityHashMap[K, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.base.LinkedIdentityHashMap", "LinkedIdentityHashMap"
}
func (this_ *LinkedIdentityHashMap[K, V]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedIdentityHashMap[K, V]) HasIdentityEquals_()                        {}
func (*LinkedIdentityHashMap[K, V]) IsGenericInst_base_LinkedIdentityHashMap_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.LinkedIdentityHashMap", "com.google.javascript.jscomp.base.LinkedIdentityHashMap", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_base_LinkedIdentityHashMap_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

func (this_ *LinkedIdentityHashMap[K, V]) initFields_LinkedIdentityHashMap() {
	_ = this_
	this_.innerMap = j2gort.NewIdentityOrderedMap[K, V]()
}

// Java constructor: LinkedIdentityHashMap()

func NewLinkedIdentityHashMap[K comparable, V any]() *LinkedIdentityHashMap[K, V] {
	this_ := &LinkedIdentityHashMap[K, V]{}
	InitLinkedIdentityHashMap[K, V](this_)
	return this_
}

func InitLinkedIdentityHashMap[K comparable, V any](this_ *LinkedIdentityHashMap[K, V]) {
	this_.initFields_LinkedIdentityHashMap()
}

// Java method: LinkedIdentityHashMap.get(K key)
func (this_ *LinkedIdentityHashMap[K, V]) Get(key K) V {
	if v_, ok_ := this_.innerMap.GetOk(key); ok_ {
		return v_
	}
	if _, ok_ := any((*V)(nil)).(*string); ok_ {
		return any(j2gort.NullString).(V)
	}
	var z_ V
	return z_
}

// Java method: LinkedIdentityHashMap.put(K key, V value)
func (this_ *LinkedIdentityHashMap[K, V]) Put(key K, value V) V {
	if old_, ok_ := this_.innerMap.PutOk(key, value); ok_ {
		return old_
	}
	if _, ok_ := any((*V)(nil)).(*string); ok_ {
		return any(j2gort.NullString).(V)
	}
	var z_ V
	return z_
}

// Java method: LinkedIdentityHashMap.computeIfAbsent(K key, java.util.function.Function<? super K,? extends V> fn)
func (this_ *LinkedIdentityHashMap[K, V]) ComputeIfAbsent(key K, fn j2gort.Ext_java_util_function_Function) V {
	_ = key
	_ = fn
	var oldValue V = this_.Get(key)
	_ = oldValue
	if !j2gort.IsNullGeneric(oldValue) {
		return oldValue
	}
	var newValue V
iifeSwitch_1:
	switch {
	default:
		tv_3 := fn.Apply(j2gort.ToAny(key))
		if tv_3 == nil {
			{

				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					newValue = any(j2gort.NullString).(V)
					break iifeSwitch_1
				}
				var z_ V
				newValue = z_
				break iifeSwitch_1
			}

		}
		newValue = j2gort.Cast[V](tv_3)
		break iifeSwitch_1
	}

	_ = newValue
	this_.Put(key, newValue)
	return newValue
}

// Java method: LinkedIdentityHashMap.merge(K key, V value, java.util.function.BiFunction<? super V,? super V,? extends V> fn)
func (this_ *LinkedIdentityHashMap[K, V]) Merge(key K, value V, fn j2gort.Ext_java_util_function_BiFunction) V {
	_ = key
	_ = value
	_ = fn
	var oldValue V = this_.Get(key)
	_ = oldValue
	if j2gort.IsNullGeneric(oldValue) {
		this_.Put(key, value)
		return value
	}
	var mergeValue V
iifeSwitch_2:
	switch {
	default:
		tv_4 := fn.Apply(j2gort.ToAny(oldValue), j2gort.ToAny(value))
		if tv_4 == nil {
			{

				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					mergeValue = any(j2gort.NullString).(V)
					break iifeSwitch_2
				}
				var z_ V
				mergeValue = z_
				break iifeSwitch_2
			}

		}
		mergeValue = j2gort.Cast[V](tv_4)
		break iifeSwitch_2
	}

	_ = mergeValue
	this_.Put(key, mergeValue)
	return mergeValue
}

// Java method: LinkedIdentityHashMap.forEach(java.util.function.BiConsumer<? super K,? super V> fn)
func (this_ *LinkedIdentityHashMap[K, V]) ForEach(fn j2gort.Ext_java_util_function_BiConsumer) {
	cur_ := this_.innerMap.IterEntries()
	for cur_.NextEntry() {
		fn.Accept(j2gort.ToAny(cur_.Key()), j2gort.ToAny(cur_.Value()))
	}
}

// Java class: com.google.javascript.jscomp.base.LinkedIdentityHashMap.Entry
type LinkedIdentityHashMap_Entry[K comparable, V any] struct {
	key     K                                  `gson:"key"`
	value   V                                  `gson:"value"`
	next    *LinkedIdentityHashMap_Entry[K, V] `gson:"next"`
	idHash_ atomic.Int32
}

func (this_ *LinkedIdentityHashMap_Entry[K, V]) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedIdentityHashMap_Entry[K, V]) IsJavaNil_() bool { return this_ == nil }
func (*LinkedIdentityHashMap_Entry[K, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.base.LinkedIdentityHashMap$Entry", "Entry"
}
func (this_ *LinkedIdentityHashMap_Entry[K, V]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedIdentityHashMap_Entry[K, V]) HasIdentityEquals_()                              {}
func (*LinkedIdentityHashMap_Entry[K, V]) IsGenericInst_base_LinkedIdentityHashMap_Entry_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.LinkedIdentityHashMap.Entry", "com.google.javascript.jscomp.base.LinkedIdentityHashMap$Entry", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_base_LinkedIdentityHashMap_Entry_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

func (this_ *LinkedIdentityHashMap_Entry[K, V]) initFields_LinkedIdentityHashMap_Entry() {
	_ = this_
}

// Java constructor: Entry(K key, V value)
func NewLinkedIdentityHashMap_Entry[K comparable, V any](key K, value V) *LinkedIdentityHashMap_Entry[K, V] {
	this_ := &LinkedIdentityHashMap_Entry[K, V]{}
	InitLinkedIdentityHashMap_Entry[K, V](this_, key, value)
	return this_
}

func InitLinkedIdentityHashMap_Entry[K comparable, V any](this_ *LinkedIdentityHashMap_Entry[K, V], key K, value V) {
	_ = key
	_ = value
	this_.initFields_LinkedIdentityHashMap_Entry()
	this_.key = key
	this_.value = value
}

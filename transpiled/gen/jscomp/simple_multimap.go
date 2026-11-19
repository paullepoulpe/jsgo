/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/SimpleMultimap.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.SimpleMultimap
type SimpleMultimap[K comparable, V any] struct {
	map_      j2gort.Ext_java_util_Map `gson:"map"`
	totalSize int32                    `gson:"totalSize"`
}

func (this_ *SimpleMultimap[K, V]) BindSelf_(s any) {
	_ = s
}
func (this_ *SimpleMultimap[K, V]) IsJavaNil_() bool { return this_ == nil }
func (*SimpleMultimap[K, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SimpleMultimap", "SimpleMultimap"
}
func (*SimpleMultimap[K, V]) IsGenericInst_jscomp_SimpleMultimap_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SimpleMultimap", "com.google.javascript.jscomp.SimpleMultimap", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_jscomp_SimpleMultimap_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

func StaticGet_SimpleMultimap_SerialVersionUID() int64 { return int64(1) }
func (this_ *SimpleMultimap[K, V]) initFields_SimpleMultimap() {
	_ = this_
}

// Java constructor: SimpleMultimap()
func NewSimpleMultimap[K comparable, V any]() *SimpleMultimap[K, V] {
	this_ := &SimpleMultimap[K, V]{}
	InitSimpleMultimap[K, V](this_)
	return this_
}

func InitSimpleMultimap[K comparable, V any](this_ *SimpleMultimap[K, V]) {
	this_.initFields_SimpleMultimap()
	this_.map_ = j2gort.NewExt_java_util_LinkedHashMap()
	this_.totalSize = int32(0)
}

// Java constructor: SimpleMultimap(com.google.javascript.jscomp.SimpleMultimap<? extends K,? extends V> other)
func NewSimpleMultimap_1[K comparable, V any](other *SimpleMultimap[K, V]) *SimpleMultimap[K, V] {
	this_ := &SimpleMultimap[K, V]{}
	InitSimpleMultimap_1[K, V](this_, other)
	return this_
}

func InitSimpleMultimap_1[K comparable, V any](this_ *SimpleMultimap[K, V], other *SimpleMultimap[K, V]) {
	_ = other
	InitSimpleMultimap[K, V](this_)
	this_.PutAll_1(other)
}

// Java static method: SimpleMultimap.create()
func SimpleMultimap_Create_0[K comparable, V any]() *SimpleMultimap[K, V] {
	return NewSimpleMultimap[K, V]()
}

// Java static method: SimpleMultimap.create(com.google.javascript.jscomp.SimpleMultimap<? extends K,? extends V> other)
func SimpleMultimap_Create_1[K comparable, V any](other *SimpleMultimap[K, V]) *SimpleMultimap[K, V] {
	_ = other
	return NewSimpleMultimap_1[K, V](other)
}

// Java method: SimpleMultimap.size()
func (this_ *SimpleMultimap[K, V]) Size() int32 {
	return this_.totalSize
}

// Java method: SimpleMultimap.isEmpty()
func (this_ *SimpleMultimap[K, V]) IsEmpty() bool {
	return (this_.totalSize == int32(0))
}

// Java method: SimpleMultimap.containsKey(java.lang.Object key)
func (this_ *SimpleMultimap[K, V]) ContainsKey_1_any(key any) bool {
	_ = key
	return this_.map_.ContainsKey(key)
}

// Java method: SimpleMultimap.containsKey(java.lang.String key)
func (this_ *SimpleMultimap[K, V]) ContainsKey_1_string(key string) bool {
	_ = key
	return j2gort.MapContainsKeyStr(this_.map_, key)
}

// Java method: SimpleMultimap.containsValue(java.lang.Object value)
func (this_ *SimpleMultimap[K, V]) ContainsValue(value any) bool {
	_ = value
	{
		cur_1 := j2gort.NewMapValCursor(this_.map_)
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			collection := j2gort.AsIface[j2gort.Ext_java_util_List](rawElem_2)
			_ = collection
			if collection.Contains(value) {
				return true
			}
		}
	}
	return false
}

// Java method: SimpleMultimap.containsEntry(java.lang.Object key, java.lang.Object value)
func (this_ *SimpleMultimap[K, V]) ContainsEntry(key any, value any) bool {
	_ = key
	_ = value
	var collection j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](this_.map_.Get(key))
	_ = collection
	return ((collection != nil) && collection.Contains(value))
}

// Java method: SimpleMultimap.put(K key, V value)
func (this_ *SimpleMultimap[K, V]) Put(key K, value V) bool {
	_ = key
	_ = value
	var list j2gort.Ext_java_util_List = func() j2gort.Ext_java_util_List {
		m_3 := this_.map_
		k_4 := j2gort.ToAny(key)
		if v_5, ok_6 := j2gort.MapTryGetComputeAny(m_3, k_4); ok_6 {
			return j2gort.AsIface[j2gort.Ext_java_util_List](v_5)
		}
		return j2gort.AsIface[j2gort.Ext_java_util_List](m_3.ComputeIfAbsent(k_4, j2gort.NewLambda_Ext_java_util_function_Function(func(lp_7 any) any {
			k := func() K {
				tv_8 := lp_7
				if tv_8 == nil {
					return func() K {
						if _, ok_ := any((*K)(nil)).(*string); ok_ {
							return any(j2gort.NullString).(K)
						}
						var z_ K
						return z_
					}()
				}
				return j2gort.Cast[K](tv_8)
			}()
			_ = k
			return j2gort.NewExt_java_util_ArrayList()
		})))
	}()
	_ = list
	var added bool = list.Add_1(j2gort.ToAny(value))
	_ = added
	if added {
		this_.totalSize = this_.totalSize + 1
	}
	return added
}

// Java method: SimpleMultimap.remove(java.lang.Object key, java.lang.Object value)
func (this_ *SimpleMultimap[K, V]) Remove(key any, value any) bool {
	_ = key
	_ = value
	var list j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](this_.map_.Get(key))
	_ = list
	if list == nil {
		return false
	}
	var removed bool = list.Remove_1_any(value)
	_ = removed
	if removed {
		this_.totalSize = this_.totalSize - 1
		if list.IsEmpty() {
			this_.map_.Remove_1(key)
		}
	}
	return removed
}

// Java method: SimpleMultimap.putAll(K key, java.lang.Iterable<? extends V> values)
func (this_ *SimpleMultimap[K, V]) PutAll_2(key K, values j2gort.Ext_java_lang_Iterable) bool {
	_ = key
	_ = values
	j2gort.CheckNotNull(values)
	var modified bool = false
	_ = modified
	{
		cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(values))
		for cur_9.HasNext() {
			rawElem_10 := cur_9.Next()
			value := func() V {
				tv_11 := rawElem_10
				if tv_11 == nil {
					return func() V {
						if _, ok_ := any((*V)(nil)).(*string); ok_ {
							return any(j2gort.NullString).(V)
						}
						var z_ V
						return z_
					}()
				}
				return j2gort.Cast[V](tv_11)
			}()
			_ = value
			modified = func() bool { ba_12 := modified; bb_13 := this_.Put(key, value); return ba_12 || bb_13 }()
		}
	}
	return modified
}

// Java method: SimpleMultimap.putAll(com.google.javascript.jscomp.SimpleMultimap<? extends K,? extends V> other)
func (this_ *SimpleMultimap[K, V]) PutAll_1(other *SimpleMultimap[K, V]) bool {
	_ = other
	j2gort.CheckNotNullPtr(other)
	var modified bool = false
	_ = modified
	{
		cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNull(other.Entries()))
		for cur_14.HasNext() {
			rawElem_15 := cur_14.Next()
			entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](rawElem_15)
			_ = entry
			modified = func() bool {
				ba_18 := modified
				bb_19 := this_.Put(func() K {
					tv_16 := entry.GetKey()
					if tv_16 == nil {
						return func() K {
							if _, ok_ := any((*K)(nil)).(*string); ok_ {
								return any(j2gort.NullString).(K)
							}
							var z_ K
							return z_
						}()
					}
					return j2gort.Cast[K](tv_16)
				}(), func() V {
					tv_17 := entry.GetValue()
					if tv_17 == nil {
						return func() V {
							if _, ok_ := any((*V)(nil)).(*string); ok_ {
								return any(j2gort.NullString).(V)
							}
							var z_ V
							return z_
						}()
					}
					return j2gort.Cast[V](tv_17)
				}())
				return ba_18 || bb_19
			}()
		}
	}
	return modified
}

// Java method: SimpleMultimap.get(java.lang.Object key)
func (this_ *SimpleMultimap[K, V]) Get_1_any(key any) j2gort.Ext_java_util_List {
	_ = key
	var list j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](this_.map_.Get(key))
	_ = list
	if list == nil {
		return j2gort.Ext_java_util_Collections_EmptyList()
	}
	return list
}

// Java method: SimpleMultimap.get(java.lang.String key)
func (this_ *SimpleMultimap[K, V]) Get_1_string(key string) j2gort.Ext_java_util_List {
	_ = key
	var list j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapGetStr(this_.map_, key))
	_ = list
	if list == nil {
		return j2gort.Ext_java_util_Collections_EmptyList()
	}
	return list
}

// Java method: SimpleMultimap.removeAll(java.lang.Object key)
func (this_ *SimpleMultimap[K, V]) RemoveAll_1_any(key any) j2gort.Ext_java_util_List {
	_ = key
	var list j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](this_.map_.Remove_1(key))
	_ = list
	if list == nil {
		return j2gort.Ext_java_util_Collections_EmptyList()
	}
	this_.totalSize = (this_.totalSize - list.Size())
	return list
}

// Java method: SimpleMultimap.removeAll(java.lang.String key)
func (this_ *SimpleMultimap[K, V]) RemoveAll_1_string(key string) j2gort.Ext_java_util_List {
	_ = key
	var list j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapRemoveStr(this_.map_, key))
	_ = list
	if list == nil {
		return j2gort.Ext_java_util_Collections_EmptyList()
	}
	this_.totalSize = (this_.totalSize - list.Size())
	return list
}

// Java method: SimpleMultimap.clear()
func (this_ *SimpleMultimap[K, V]) Clear() {
	this_.map_.Clear()
	this_.totalSize = int32(0)
}

func (this_ *SimpleMultimap[K, V]) ContainsKeyStr_(key string) bool {
	return j2gort.MapContainsKeyStr(this_.map_, key)
}

func (this_ *SimpleMultimap[K, V]) ContainsKeyNode_(key uint32) bool {
	return j2gort.MapContainsKeyNode(this_.map_, key)
}

func (this_ *SimpleMultimap[K, V]) ContainsEntryStr_(key string, value any) bool {
	list := j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapGetStr(this_.map_, key))
	return list != nil && list.Contains(j2gort.ToAny(value))
}

func (this_ *SimpleMultimap[K, V]) ContainsEntryNode_(key uint32, value any) bool {
	list := j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapGetNode(this_.map_, key))
	return list != nil && list.Contains(j2gort.ToAny(value))
}

func (this_ *SimpleMultimap[K, V]) PutStr_(key string, value any) bool {
	value = j2gort.ToAny(value)
	list := j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapGetStr(this_.map_, key))
	if list == nil {
		al := j2gort.NewExt_java_util_ArrayList()
		j2gort.MapPutStr(this_.map_, key, al)
		list = al
	}
	if list.Add_1(value) {
		this_.totalSize++
		return true
	}
	return false
}

func (this_ *SimpleMultimap[K, V]) PutNode_(key uint32, value any) bool {
	value = j2gort.ToAny(value)
	list := j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapGetNode(this_.map_, key))
	if list == nil {
		al := j2gort.NewExt_java_util_ArrayList()
		j2gort.MapPutNode(this_.map_, key, al)
		list = al
	}
	if list.Add_1(value) {
		this_.totalSize++
		return true
	}
	return false
}

func (this_ *SimpleMultimap[K, V]) GetStr_(key string) j2gort.Ext_java_util_List {
	list := j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapGetStr(this_.map_, key))
	if list == nil {
		return j2gort.Ext_java_util_Collections_EmptyList()
	}
	return list
}

func (this_ *SimpleMultimap[K, V]) GetNode_(key uint32) j2gort.Ext_java_util_List {
	list := j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapGetNode(this_.map_, key))
	if list == nil {
		return j2gort.Ext_java_util_Collections_EmptyList()
	}
	return list
}

func (this_ *SimpleMultimap[K, V]) RemoveStr_(key string, value any) bool {
	list := j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapGetStr(this_.map_, key))
	if list == nil {
		return false
	}
	if list.Remove_1_any(j2gort.ToAny(value)) {
		this_.totalSize--
		if list.IsEmpty() {
			j2gort.MapRemoveStr(this_.map_, key)
		}
		return true
	}
	return false
}

func (this_ *SimpleMultimap[K, V]) RemoveNode_(key uint32, value any) bool {
	list := j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapGetNode(this_.map_, key))
	if list == nil {
		return false
	}
	if list.Remove_1_any(j2gort.ToAny(value)) {
		this_.totalSize--
		if list.IsEmpty() {
			j2gort.MapRemoveNode(this_.map_, key)
		}
		return true
	}
	return false
}

func (this_ *SimpleMultimap[K, V]) RemoveAllStr_(key string) j2gort.Ext_java_util_List {
	list := j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapRemoveStr(this_.map_, key))
	if list == nil {
		return j2gort.Ext_java_util_Collections_EmptyList()
	}
	this_.totalSize -= list.Size()
	return list
}

func (this_ *SimpleMultimap[K, V]) RemoveAllNode_(key uint32) j2gort.Ext_java_util_List {
	list := j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MapRemoveNode(this_.map_, key))
	if list == nil {
		return j2gort.Ext_java_util_Collections_EmptyList()
	}
	this_.totalSize -= list.Size()
	return list
}

func (this_ *SimpleMultimap[K, V]) PutAllStr_(key string, values j2gort.Ext_java_lang_Iterable) bool {
	j2gort.CheckNotNull(values)
	modified := false
	cur := j2gort.NewIterCursor(values)
	for cur.HasNext() {
		if this_.PutStr_(key, cur.Next()) {
			modified = true
		}
	}
	return modified
}

func (this_ *SimpleMultimap[K, V]) PutAllNode_(key uint32, values j2gort.Ext_java_lang_Iterable) bool {
	j2gort.CheckNotNull(values)
	modified := false
	cur := j2gort.NewIterCursor(values)
	for cur.HasNext() {
		if this_.PutNode_(key, cur.Next()) {
			modified = true
		}
	}
	return modified
}

// Java method: SimpleMultimap.keySet()
func (this_ *SimpleMultimap[K, V]) KeySet() j2gort.Ext_java_util_Set {
	return j2gort.AsIface[j2gort.Ext_java_util_Set](this_.map_.KeySet_0())
}

// Java method: SimpleMultimap.values()
func (this_ *SimpleMultimap[K, V]) Values() j2gort.Ext_java_util_Collection {
	return NewSimpleMultimap__Anon_80[K, V](this_)
}

// Java method: SimpleMultimap.entries()
func (this_ *SimpleMultimap[K, V]) Entries() j2gort.Ext_java_util_Collection {
	var entries []j2gort.Ext_java_util_AbstractMap_SimpleImmutableEntry = j2gort.MakeNonNilSlice[j2gort.Ext_java_util_AbstractMap_SimpleImmutableEntry](this_.totalSize)
	_ = entries
	{
		cur_20 := j2gort.NewMapEntryCursor(this_.map_)
		for cur_20.NextEntry() {
			var key K = func() K {
				tv_21 := cur_20.GetKey()
				if tv_21 == nil {
					return func() K {
						if _, ok_ := any((*K)(nil)).(*string); ok_ {
							return any(j2gort.NullString).(K)
						}
						var z_ K
						return z_
					}()
				}
				return j2gort.Cast[K](tv_21)
			}()
			_ = key
			{
				cur_22 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_List](cur_20.GetValue())))
				for cur_22.HasNext() {
					rawElem_23 := cur_22.Next()
					value := func() V {
						tv_24 := rawElem_23
						if tv_24 == nil {
							return func() V {
								if _, ok_ := any((*V)(nil)).(*string); ok_ {
									return any(j2gort.NullString).(V)
								}
								var z_ V
								return z_
							}()
						}
						return j2gort.Cast[V](tv_24)
					}()
					_ = value
					entries = append(entries, j2gort.NewExt_java_util_AbstractMap_SimpleImmutableEntry_2(j2gort.ToAny(key), j2gort.ToAny(value)))
				}
			}
		}
	}
	return j2gort.WrapTypedArrayList(entries)
}

// Java method: SimpleMultimap.asMap()
func (this_ *SimpleMultimap[K, V]) AsMap() j2gort.Ext_java_util_Map {
	return j2gort.Ext_java_util_Collections_UnmodifiableMap(this_.map_)
}

// Java method: SimpleMultimap.equals(java.lang.Object o)
func (this_ *SimpleMultimap[K, V]) Equals(arg_25 any) bool {
	o := arg_25
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	{
		var iifeCond_1 bool
		{

			v_, ok_ := o.(interface {
				IsGenericInst_jscomp_SimpleMultimap_()
				IsJavaNil_() bool
			})
			iifeCond_1 = ok_ && !v_.IsJavaNil_()
		}
		if !(iifeCond_1) {
			return false
		}
	}
	var other *SimpleMultimap[K, V] = j2gort.CastPtr[SimpleMultimap[K, V]](j2gort.CastPtr[SimpleMultimap[any, any]](o))
	_ = other
	return this_.map_.Equals(other.map_)
}

// Java method: SimpleMultimap.hashCode()
func (this_ *SimpleMultimap[K, V]) HashCode() int32 {
	return this_.map_.HashCode()
}

// Java method: SimpleMultimap.toString()
func (this_ *SimpleMultimap[K, V]) String() string {
	return javalang.ToJavaString(j2gort.CheckNotNull(this_.map_))
}

// Java class:
type SimpleMultimap__Anon_80[K comparable, V any] struct {
	j2gort.Ext_java_util_AbstractCollectionBase
	outer_  *SimpleMultimap[K, V]
	idHash_ atomic.Int32
}

func (this_ *SimpleMultimap__Anon_80[K, V]) BindSelf_(s any) {
	_ = s
	this_.Ext_java_util_AbstractCollectionBase.BindExtSelf_(s)
}
func (this_ *SimpleMultimap__Anon_80[K, V]) IsJavaNil_() bool { return this_ == nil }
func (*SimpleMultimap__Anon_80[K, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SimpleMultimap$1", ""
}
func (this_ *SimpleMultimap__Anon_80[K, V]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SimpleMultimap__Anon_80[K, V]) IsGenericInst_jscomp_SimpleMultimap__Anon_80_()      {}
func (this_ *SimpleMultimap__Anon_80[K, V]) SetOuter_(o *SimpleMultimap[K, V])             { this_.outer_ = o }
func (this_ *SimpleMultimap__Anon_80[K, V]) GetOuter_() *SimpleMultimap[K, V]              { return this_.outer_ }
func (*SimpleMultimap__Anon_80[K, V]) Implements_j2gort_Ext_java_util_AbstractCollection() {}
func (*SimpleMultimap__Anon_80[K, V]) Implements_j2gort_Ext_java_util_Collection()         {}
func (*SimpleMultimap__Anon_80[K, V]) Implements_j2gort_Ext_java_lang_Iterable()           {}

func (this_ *SimpleMultimap__Anon_80[K, V]) initFields_SimpleMultimap__Anon_80() {
	_ = this_
}

func NewSimpleMultimap__Anon_80[K comparable, V any](outer_ *SimpleMultimap[K, V]) *SimpleMultimap__Anon_80[K, V] {
	this_ := &SimpleMultimap__Anon_80[K, V]{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	j2gort.InitExt_java_util_AbstractCollection(&this_.Ext_java_util_AbstractCollectionBase)
	this_.initFields_SimpleMultimap__Anon_80()
	return this_
}

// Java method: .iterator()
func (this_ *SimpleMultimap__Anon_80[K, V]) Iterator() any {
	return NewSimpleMultimap__Anon_81[K, V](this_)
}

// Java method: .size()
func (this_ *SimpleMultimap__Anon_80[K, V]) Size() int32 {
	return this_.outer_.totalSize
}

// Java class:
type SimpleMultimap__Anon_81[K comparable, V any] struct {
	outer_       *SimpleMultimap__Anon_80[K, V]
	mapIterator  j2gort.Ext_java_util_Iterator `gson:"mapIterator"`
	listIterator j2gort.Ext_java_util_Iterator `gson:"listIterator"`
	idHash_      atomic.Int32
}

func (this_ *SimpleMultimap__Anon_81[K, V]) BindSelf_(s any) {
	_ = s
}
func (this_ *SimpleMultimap__Anon_81[K, V]) IsJavaNil_() bool { return this_ == nil }
func (*SimpleMultimap__Anon_81[K, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SimpleMultimap$1$1", ""
}
func (this_ *SimpleMultimap__Anon_81[K, V]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SimpleMultimap__Anon_81[K, V]) HasIdentityEquals_()                            {}
func (*SimpleMultimap__Anon_81[K, V]) IsGenericInst_jscomp_SimpleMultimap__Anon_81_() {}
func (this_ *SimpleMultimap__Anon_81[K, V]) SetOuter_(o *SimpleMultimap__Anon_80[K, V]) {
	this_.outer_ = o
}
func (this_ *SimpleMultimap__Anon_81[K, V]) GetOuter_() *SimpleMultimap__Anon_80[K, V] {
	return this_.outer_
}
func (*SimpleMultimap__Anon_81[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}

func (this_ *SimpleMultimap__Anon_81[K, V]) initFields_SimpleMultimap__Anon_81() {
	_ = this_
	this_.mapIterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.outer_.outer_.map_.Values()).Iterator())
	this_.listIterator = j2gort.Ext_java_util_Collections_EmptyIterator()
}

func NewSimpleMultimap__Anon_81[K comparable, V any](outer_ *SimpleMultimap__Anon_80[K, V]) *SimpleMultimap__Anon_81[K, V] {
	this_ := &SimpleMultimap__Anon_81[K, V]{}
	this_.outer_ = outer_
	this_.initFields_SimpleMultimap__Anon_81()
	return this_
}

// Java method: .hasNext()
func (this_ *SimpleMultimap__Anon_81[K, V]) HasNext() bool {
	for !(this_.listIterator.HasNext()) && this_.mapIterator.HasNext() {
		this_.listIterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.AsIface[j2gort.Ext_java_util_List](this_.mapIterator.Next()).Iterator())
	}
	return this_.listIterator.HasNext()
}

// Java method: .next()
func (this_ *SimpleMultimap__Anon_81[K, V]) Next() any {
	return j2gort.ToAny(func() V {
		tv_26 := this_.listIterator.Next()
		if tv_26 == nil {
			return func() V {
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					return any(j2gort.NullString).(V)
				}
				var z_ V
				return z_
			}()
		}
		return j2gort.Cast[V](tv_26)
	}())
}

func (this_ *SimpleMultimap__Anon_81[K, V]) Remove() {
	panic(javalang.NewUnsupportedOperationException("remove"))
}

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/HamtPMap.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.HamtPMap
type HamtPMap[K comparable, V comparable] struct {
	key      K                 `gson:"key"`
	value    V                 `gson:"value"`
	children []*HamtPMap[K, V] `gson:"children"`
	hash     int32             `gson:"hash"`
	mask     int32             `gson:"mask"`
	idHash_  atomic.Int32
	present  bool `gson:"present"`
}

func (this_ *HamtPMap[K, V]) BindSelf_(s any) {
	_ = s
}
func (this_ *HamtPMap[K, V]) IsJavaNil_() bool { return this_ == nil }
func (*HamtPMap[K, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.HamtPMap", "HamtPMap"
}
func (this_ *HamtPMap[K, V]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*HamtPMap[K, V]) HasIdentityEquals_()            {}
func (*HamtPMap[K, V]) IsGenericInst_rhino_HamtPMap_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.HamtPMap", "com.google.javascript.rhino.HamtPMap", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_rhino_HamtPMap_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*HamtPMap[K, V]) Implements_rhino_PMap() {}

func StaticGet_HamtPMap_BITS() int32       { return int32(4) }
func StaticGet_HamtPMap_BITS_SHIFT() int32 { return int32(28) }
func (this_ *HamtPMap[K, V]) initFields_HamtPMap() {
	_ = this_
}

// Java constructor: HamtPMap(@org.jspecify.annotations.Nullable K key, int hash, @org.jspecify.annotations.Nullable V value, boolean present, int mask, com.google.javascript.rhino.HamtPMap<K,V> @org.jspecify.annotations.Nullable [] children)
func NewHamtPMap[K comparable, V comparable](key K, hash int32, value V, present bool, mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	this_ := &HamtPMap[K, V]{}
	InitHamtPMap[K, V](this_, key, hash, value, present, mask, children)
	return this_
}

func InitHamtPMap[K comparable, V comparable](this_ *HamtPMap[K, V], key K, hash int32, value V, present bool, mask int32, children []*HamtPMap[K, V]) {
	_ = key
	_ = hash
	_ = value
	_ = present
	_ = mask
	_ = children
	this_.initFields_HamtPMap()
	this_.key = key
	this_.hash = hash
	this_.value = value
	this_.present = present
	this_.mask = mask
	this_.children = children
}

// Java static method: HamtPMap.empty()
func HamtPMap_Empty[K comparable, V comparable]() *HamtPMap[K, V] {
	return NewHamtPMap[K, V](func() K {
		if _, ok_ := any((*K)(nil)).(*string); ok_ {
			return any(j2gort.NullString).(K)
		}
		var z_ K
		return z_
	}(), int32(0), func() V {
		if _, ok_ := any((*V)(nil)).(*string); ok_ {
			return any(j2gort.NullString).(V)
		}
		var z_ V
		return z_
	}(), false, int32(0), nil)
}

// Java method: HamtPMap.toString()
func (this_ *HamtPMap[K, V]) String() string {
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.NewExt_java_lang_StringBuilder().Append_1_string("{"))
	_ = sb
	if !(this_.IsEmpty()) {
		this_.priv_HamtPMap_appendTo_7984(sb)
	}
	return j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string("}")).String()
}

// Java method: HamtPMap.appendTo(java.lang.StringBuilder sb)
func (this_ *HamtPMap[K, V]) priv_HamtPMap_appendTo_7984(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	if sb.Length() > int32(1) {
		sb.Append_1_string(", ")
	}
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_any_9859(j2gort.ToAny(this_.key))).Append_1_string(": ")).Append_1_any_9859(j2gort.ToAny(this_.value))
	if this_.children != nil {
		for _, el_1 := range j2gort.CheckNotNullSlice(this_.children) {
			child := el_1
			_ = child
			child.priv_HamtPMap_appendTo_7984(sb)
		}
	}
}

// Java method: HamtPMap.isEmpty()
func (this_ *HamtPMap[K, V]) IsEmpty() bool {
	return (!(this_.present) && (this_.mask == int32(0)))
}

// Java method: HamtPMap.values()
func (this_ *HamtPMap[K, V]) Values() j2gort.Ext_java_lang_Iterable {
	if this_.IsEmpty() {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_EmptyList())
	}
	return j2gort.NewLambda_Ext_java_lang_Iterable(func() any {
		return NewHamtPMap_ValueIter[K, V](this_)
	})
}

// Java method: HamtPMap.keys()
func (this_ *HamtPMap[K, V]) Keys() j2gort.Ext_java_lang_Iterable {
	if this_.IsEmpty() {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_EmptyList())
	}
	return j2gort.NewLambda_Ext_java_lang_Iterable(func() any {
		return NewHamtPMap_KeyIter[K, V](this_)
	})
}

// Java method: HamtPMap.get(K key)
func (this_ *HamtPMap[K, V]) Get_1(key K) V {
	_ = key
	if !(this_.IsEmpty()) {
		return this_.priv_HamtPMap_get_7995(key, HamtPMap_Hash[K](key))
	} else {
		return func() V {
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				return any(j2gort.NullString).(V)
			}
			var z_ V
			return z_
		}()
	}
}

// Java method: HamtPMap.containsHash(int keyHash)
func (this_ *HamtPMap[K, V]) ContainsHash(keyHash int32) bool {
	_ = keyHash
	if this_.IsEmpty() {
		return false
	}
	var h int32 = j2gort.BoxInt_Reverse(keyHash)
	_ = h
	var cur *HamtPMap[K, V] = this_
	_ = cur
	for cur != nil {
		if cur.present && (cur.hash == h) {
			return true
		}
		if HamtPMap_CompareUnsigned(h, cur.hash) < int32(0) {
			return false
		}
		var bucketMask int32 = javalang.Shl32(int32(1), int(HamtPMap_Bucket(h)))
		_ = bucketMask
		if (cur.mask & bucketMask) == int32(0) {
			return false
		}
		cur = cur.children[cur.priv_HamtPMap_index_8005(bucketMask)]
		h = HamtPMap_Shift(h)
	}
	return false
}

// Java static method: HamtPMap.singleton(K key, V value)
func HamtPMap_Singleton[K comparable, V comparable](key K, value V) *HamtPMap[K, V] {
	_ = key
	_ = value
	func() V {
		tv_2 := value
		if _, ok_ := any((*V)(nil)).(*string); ok_ {
			if tv_2 == any(j2gort.NullString).(V) {
				j2gort.ThrowNPE()
			}
			return tv_2
		}
		var z_ V
		if tv_2 == z_ {
			j2gort.ThrowNPE()
		}
		return tv_2
	}()
	return NewHamtPMap[K, V](key, HamtPMap_Hash[K](key), value, true, int32(0), nil)
}

// Java static method: HamtPMap.ofTwo(K k1, V v1, K k2, V v2)
func HamtPMap_OfTwo[K comparable, V comparable](k1 K, v1 V, k2 K, v2 V) *HamtPMap[K, V] {
	_ = k1
	_ = v1
	_ = k2
	_ = v2
	func() V {
		tv_3 := v1
		if _, ok_ := any((*V)(nil)).(*string); ok_ {
			if tv_3 == any(j2gort.NullString).(V) {
				j2gort.ThrowNPE()
			}
			return tv_3
		}
		var z_ V
		if tv_3 == z_ {
			j2gort.ThrowNPE()
		}
		return tv_3
	}()
	func() V {
		tv_4 := v2
		if _, ok_ := any((*V)(nil)).(*string); ok_ {
			if tv_4 == any(j2gort.NullString).(V) {
				j2gort.ThrowNPE()
			}
			return tv_4
		}
		var z_ V
		if tv_4 == z_ {
			j2gort.ThrowNPE()
		}
		return tv_4
	}()
	var h1 int32 = HamtPMap_Hash[K](k1)
	_ = h1
	var h2 int32 = HamtPMap_Hash[K](k2)
	_ = h2
	if (h1 == h2) && func() bool {
		rk_5 := k2
		ak_6 := k1
		if _, ok_ := any((*K)(nil)).(*string); ok_ {
			_ = j2gort.CheckNotNullString(any(rk_5).(string))
			return rk_5 == ak_6
		}
		var z_ K
		if rk_5 == z_ {
			j2gort.ThrowNPE()
		}
		if rk_5 == ak_6 {
			return true
		}
		return javalang.ObjectsEquals(j2gort.ToAny(rk_5), j2gort.ToAny(ak_6))
	}() {
		var sameVal bool = func() bool {
			rk_7 := v2
			ak_8 := v1
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				_ = j2gort.CheckNotNullString(any(rk_7).(string))
				return rk_7 == ak_8
			}
			var z_ V
			if rk_7 == z_ {
				j2gort.ThrowNPE()
			}
			if rk_7 == ak_8 {
				return true
			}
			return javalang.ObjectsEquals(j2gort.ToAny(rk_7), j2gort.ToAny(ak_8))
		}()
		_ = sameVal
		return NewHamtPMap[K, V](func() K {
			if sameVal {
				return k1
			}
			return k2
		}(), h1, func() V {
			if sameVal {
				return v1
			}
			return v2
		}(), true, int32(0), nil)
	}
	if HamtPMap_CompareUnsigned(h2, h1) < int32(0) {
		var child *HamtPMap[K, V] = NewHamtPMap[K, V](k1, HamtPMap_Shift(h1), v1, true, int32(0), nil)
		_ = child
		var arr []*HamtPMap[K, V] = []*HamtPMap[K, V]{child}
		_ = arr
		return NewHamtPMap[K, V](k2, h2, v2, true, javalang.Shl32(int32(1), int(HamtPMap_Bucket(h1))), arr)
	}
	var child_2 *HamtPMap[K, V] = NewHamtPMap[K, V](k2, HamtPMap_Shift(h2), v2, true, int32(0), nil)
	_ = child_2
	var arr_2 []*HamtPMap[K, V] = []*HamtPMap[K, V]{child_2}
	_ = arr_2
	return NewHamtPMap[K, V](k1, h1, v1, true, javalang.Shl32(int32(1), int(HamtPMap_Bucket(h2))), arr_2)
}

// Java method: HamtPMap.isSingleEntry()
func (this_ *HamtPMap[K, V]) IsSingleEntry() bool {
	return (this_.present && (this_.mask == int32(0)))
}

// Java method: HamtPMap.rootKey()
func (this_ *HamtPMap[K, V]) RootKey() K {
	if this_.present {
		return this_.key
	} else {
		return func() K {
			if _, ok_ := any((*K)(nil)).(*string); ok_ {
				return any(j2gort.NullString).(K)
			}
			var z_ K
			return z_
		}()
	}
}

// Java method: HamtPMap.rootValue()
func (this_ *HamtPMap[K, V]) RootValue() V {
	if this_.present {
		return this_.value
	} else {
		return func() V {
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				return any(j2gort.NullString).(V)
			}
			var z_ V
			return z_
		}()
	}
}

// Java method: HamtPMap.childCount()
func (this_ *HamtPMap[K, V]) ChildCount() int32 {
	if this_.children == nil {
		return int32(0)
	} else {
		return int32(len(this_.children))
	}
}

// Java method: HamtPMap.childAt(int index)
func (this_ *HamtPMap[K, V]) ChildAt(index int32) *HamtPMap[K, V] {
	_ = index
	return this_.children[index]
}

// Java method: HamtPMap.get(K key, int hash)
func (this_ *HamtPMap[K, V]) priv_HamtPMap_get_7995(key K, hash int32) V {
	_ = key
	_ = hash
	if (hash == this_.hash) && func() bool {
		rk_9 := key
		ak_10 := this_.key
		if _, ok_ := any((*K)(nil)).(*string); ok_ {
			_ = j2gort.CheckNotNullString(any(rk_9).(string))
			return rk_9 == ak_10
		}
		var z_ K
		if rk_9 == z_ {
			j2gort.ThrowNPE()
		}
		if rk_9 == ak_10 {
			return true
		}
		return javalang.ObjectsEquals(j2gort.ToAny(rk_9), j2gort.ToAny(ak_10))
	}() {
		return this_.value
	}
	var bucket int32 = HamtPMap_Bucket(hash)
	_ = bucket
	var bucketMask int32 = javalang.Shl32(int32(1), int(bucket))
	_ = bucketMask
	if (this_.mask & bucketMask) != int32(0) {
		return this_.children[this_.priv_HamtPMap_index_8005(bucketMask)].priv_HamtPMap_get_7995(key, HamtPMap_Shift(hash))
	} else {
		return func() V {
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				return any(j2gort.NullString).(V)
			}
			var z_ V
			return z_
		}()
	}
}

// Java method: HamtPMap.plus(K key, V value)
func (this_ *HamtPMap[K, V]) Plus_2(key K, value V) PMap[K, V] {
	_ = key
	_ = value
	func() V {
		tv_11 := value
		if _, ok_ := any((*V)(nil)).(*string); ok_ {
			if tv_11 == any(j2gort.NullString).(V) {
				j2gort.ThrowNPE()
			}
			return tv_11
		}
		var z_ V
		if tv_11 == z_ {
			j2gort.ThrowNPE()
		}
		return tv_11
	}()
	if !(this_.IsEmpty()) {
		return ToIface_PMap[K, V](this_.priv_HamtPMap_plus_7997(key, HamtPMap_Hash[K](key), value))
	} else {
		return ToIface_PMap[K, V](NewHamtPMap[K, V](key, HamtPMap_Hash[K](key), value, true, int32(0), nil))
	}
}

// Java method: HamtPMap.plus(K key, int hash, V value)
func (this_ *HamtPMap[K, V]) priv_HamtPMap_plus_7997(key K, hash int32, value V) *HamtPMap[K, V] {
	_ = key
	_ = hash
	_ = value
	if (hash == this_.hash) && func() bool {
		rk_12 := key
		ak_13 := this_.key
		if _, ok_ := any((*K)(nil)).(*string); ok_ {
			_ = j2gort.CheckNotNullString(any(rk_12).(string))
			return rk_12 == ak_13
		}
		var z_ K
		if rk_12 == z_ {
			j2gort.ThrowNPE()
		}
		if rk_12 == ak_13 {
			return true
		}
		return javalang.ObjectsEquals(j2gort.ToAny(rk_12), j2gort.ToAny(ak_13))
	}() {
		if func() bool {
			rk_14 := value
			ak_15 := this_.value
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				_ = j2gort.CheckNotNullString(any(rk_14).(string))
				return rk_14 == ak_15
			}
			var z_ V
			if rk_14 == z_ {
				j2gort.ThrowNPE()
			}
			if rk_14 == ak_15 {
				return true
			}
			return javalang.ObjectsEquals(j2gort.ToAny(rk_14), j2gort.ToAny(ak_15))
		}() {
			return this_
		} else {
			return NewHamtPMap[K, V](key, hash, value, true, this_.mask, this_.children)
		}
	}
	if HamtPMap_CompareUnsigned(hash, this_.hash) < int32(0) {
		return this_.priv_HamtPMap_replaceRoot_7998(key, hash, value)
	}
	var bucket int32 = HamtPMap_Bucket(hash)
	_ = bucket
	hash = HamtPMap_Shift(hash)
	var bucketMask int32 = javalang.Shl32(int32(1), int(bucket))
	_ = bucketMask
	var index int32 = this_.priv_HamtPMap_index_8005(bucketMask)
	_ = index
	if (this_.mask & bucketMask) != int32(0) {
		var child *HamtPMap[K, V] = this_.children[index]
		_ = child
		var newChild *HamtPMap[K, V] = child.priv_HamtPMap_plus_7997(key, hash, value)
		_ = newChild
		if child == newChild {
			return this_
		} else {
			return this_.priv_HamtPMap_withChildren_8010(this_.mask, HamtPMap_ReplaceChild[K, V](this_.children, index, newChild))
		}
	} else {
		var newChild_2 *HamtPMap[K, V] = NewHamtPMap[K, V](key, hash, value, true, int32(0), nil)
		_ = newChild_2
		return this_.priv_HamtPMap_withChildren_8010((this_.mask | bucketMask), HamtPMap_InsertChild[K, V](this_.children, index, newChild_2))
	}
}

// Java method: HamtPMap.replaceRoot(K key, int hash, V value)
func (this_ *HamtPMap[K, V]) priv_HamtPMap_replaceRoot_7998(key K, hash int32, value V) *HamtPMap[K, V] {
	_ = key
	_ = hash
	_ = value
	var bucket int32 = HamtPMap_Bucket(this_.hash)
	_ = bucket
	var leafHash int32 = HamtPMap_Shift(this_.hash)
	_ = leafHash
	var bucketMask int32 = javalang.Shl32(int32(1), int(bucket))
	_ = bucketMask
	var index int32 = this_.priv_HamtPMap_index_8005(bucketMask)
	_ = index
	var newChildren []*HamtPMap[K, V] = nil
	_ = newChildren
	if (this_.mask & bucketMask) != int32(0) {
		newChildren = HamtPMap_ReplaceChild[K, V](this_.children, index, this_.children[index].priv_HamtPMap_plus_7997(this_.key, leafHash, this_.value))
	} else {
		var newChild *HamtPMap[K, V] = NewHamtPMap[K, V](this_.key, leafHash, this_.value, true, int32(0), nil)
		_ = newChild
		newChildren = HamtPMap_InsertChild[K, V](this_.children, index, newChild)
	}
	return NewHamtPMap[K, V](key, hash, value, true, (this_.mask | bucketMask), newChildren)
}

// Java method: HamtPMap.minus(K key)
func (this_ *HamtPMap[K, V]) Minus(key K) PMap[K, V] {
	_ = key
	if this_.IsEmpty() {
		return ToIface_PMap[K, V](this_)
	}
	var result *HamtPMap[K, V] = this_.priv_HamtPMap_minusOrNull_8001(key, HamtPMap_Hash[K](key), nil)
	_ = result
	if result != nil {
		return ToIface_PMap[K, V](result)
	} else {
		return ToIface_PMap[K, V](HamtPMap_Empty[K, V]())
	}
}

// Java method: HamtPMap.minusOrNull(K key)
func (this_ *HamtPMap[K, V]) MinusOrNull_1(key K) *HamtPMap[K, V] {
	_ = key
	if !(this_.IsEmpty()) {
		return this_.priv_HamtPMap_minusOrNull_8001(key, HamtPMap_Hash[K](key), nil)
	} else {
		return nil
	}
}

// Java method: HamtPMap.minusOrNull(K key, int hash, com.google.javascript.rhino.HamtPMap.@org.jspecify.annotations.Nullable Holder<V> removed)
func (this_ *HamtPMap[K, V]) priv_HamtPMap_minusOrNull_8001(key K, hash int32, removed *HamtPMap_Holder[V]) *HamtPMap[K, V] {
	_ = key
	_ = hash
	_ = removed
	if (hash == this_.hash) && func() bool {
		rk_16 := key
		ak_17 := this_.key
		if _, ok_ := any((*K)(nil)).(*string); ok_ {
			_ = j2gort.CheckNotNullString(any(rk_16).(string))
			return rk_16 == ak_17
		}
		var z_ K
		if rk_16 == z_ {
			j2gort.ThrowNPE()
		}
		if rk_16 == ak_17 {
			return true
		}
		return javalang.ObjectsEquals(j2gort.ToAny(rk_16), j2gort.ToAny(ak_17))
	}() {
		var result *HamtPMap[K, V] = HamtPMap_DeleteRoot[K, V](this_.mask, this_.children)
		_ = result
		if removed != nil {
			{
				r_18 := removed
				r_18.found = true
			}
			{
				r_19 := removed
				r_19.value = this_.value
			}
		}
		return result
	}
	var bucket int32 = HamtPMap_Bucket(hash)
	_ = bucket
	var bucketMask int32 = javalang.Shl32(int32(1), int(bucket))
	_ = bucketMask
	if (this_.mask & bucketMask) == int32(0) {
		return this_
	}
	hash = HamtPMap_Shift(hash)
	var index int32 = this_.priv_HamtPMap_index_8005(bucketMask)
	_ = index
	var child *HamtPMap[K, V] = this_.children[index]
	_ = child
	var newChild *HamtPMap[K, V] = child.priv_HamtPMap_minusOrNull_8001(key, hash, removed)
	_ = newChild
	if newChild == child {
		return this_
	} else if newChild == nil {
		return this_.priv_HamtPMap_withChildren_8010((this_.mask & ^(bucketMask)), HamtPMap_DeleteChild[K, V](this_.children, index))
	} else {
		return this_.priv_HamtPMap_withChildren_8010(this_.mask, HamtPMap_ReplaceChild[K, V](this_.children, index, newChild))
	}
}

// Java method: HamtPMap.reconcile(com.google.javascript.rhino.PMap<K,V> that, com.google.javascript.rhino.PMap.Reconciler<K,V> joiner)
func (this_ *HamtPMap[K, V]) Reconcile(that PMap[K, V], joiner PMap_Reconciler[K, V]) PMap[K, V] {
	_ = that
	_ = joiner
	if this_.IsEmpty() && that.IsEmpty() {
		return ToIface_PMap[K, V](this_)
	}
	var result *HamtPMap[K, V] = HamtPMap_Reconcile[K, V](func() *HamtPMap[K, V] {
		if !(this_.IsEmpty()) {
			return this_
		}
		return nil
	}(), func() *HamtPMap[K, V] {
		if !(that.IsEmpty()) {
			return j2gort.CastPtr[HamtPMap[K, V]](that)
		}
		return nil
	}(), joiner)
	_ = result
	if result != nil {
		return ToIface_PMap[K, V](result)
	} else {
		return ToIface_PMap[K, V](HamtPMap_Empty[K, V]())
	}
}

// Java method: HamtPMap.reconcileOrNull(com.google.javascript.rhino.HamtPMap<K,V> that, com.google.javascript.rhino.PMap.Reconciler<K,V> joiner)
func (this_ *HamtPMap[K, V]) ReconcileOrNull(that *HamtPMap[K, V], joiner PMap_Reconciler[K, V]) *HamtPMap[K, V] {
	_ = that
	_ = joiner
	return HamtPMap_Reconcile[K, V](func() *HamtPMap[K, V] {
		if !(this_.IsEmpty()) {
			return this_
		}
		return nil
	}(), func() *HamtPMap[K, V] {
		if !(that.IsEmpty()) {
			return that
		}
		return nil
	}(), joiner)
}

// Java static method: HamtPMap.reconcile(com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<K,V> t1, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<K,V> t2, com.google.javascript.rhino.PMap.Reconciler<K,V> joiner)
func HamtPMap_Reconcile[K comparable, V comparable](t1 *HamtPMap[K, V], t2 *HamtPMap[K, V], joiner PMap_Reconciler[K, V]) *HamtPMap[K, V] {
	_ = t1
	_ = t2
	_ = joiner
	if t1 == t2 {
		return t1
	} else if t1 == nil {
		var newValue V = func() V {
			tv_20 := joiner.Merge(t2.key, func() V {
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					return any(j2gort.NullString).(V)
				}
				var z_ V
				return z_
			}(), t2.value)
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				if tv_20 == any(j2gort.NullString).(V) {
					j2gort.ThrowNPE()
				}
				return tv_20
			}
			var z_ V
			if tv_20 == z_ {
				j2gort.ThrowNPE()
			}
			return tv_20
		}()
		_ = newValue
		var childMask int32 = t2.mask
		_ = childMask
		if t2.children == nil {
			if t2.present && func() bool {
				rk_21 := t2.value
				ak_22 := newValue
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					_ = j2gort.CheckNotNullString(any(rk_21).(string))
					return rk_21 == ak_22
				}
				var z_ V
				if rk_21 == z_ {
					j2gort.ThrowNPE()
				}
				if rk_21 == ak_22 {
					return true
				}
				return javalang.ObjectsEquals(j2gort.ToAny(rk_21), j2gort.ToAny(ak_22))
			}() {
				return t2
			} else {
				return NewHamtPMap[K, V](t2.key, t2.hash, newValue, true, int32(0), nil)
			}
		}
		var copied []*HamtPMap[K, V] = nil
		_ = copied
		var rem int32 = childMask
		_ = rem
		var outIdx int32 = int32(0)
		_ = outIdx
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = int32(len(t2.children))
			_ = len_
			for ; i < len_; i = i + 1 {
				var bit int32 = j2gort.IntLowestOneBit(rem)
				_ = bit
				rem = (rem & ^(bit))
				var origChild *HamtPMap[K, V] = t2.children[i]
				_ = origChild
				var c *HamtPMap[K, V] = HamtPMap_Reconcile[K, V](nil, origChild, joiner)
				_ = c
				if copied == nil {
					if c == origChild {
						outIdx = outIdx + 1
						continue
					}
					var alloc []*HamtPMap[K, V] = make([]*HamtPMap[K, V], len_)
					_ = alloc
					{
						var j int32 = int32(0)
						_ = j
						for ; j < i; j = j + 1 {
							{
								a_23 := alloc
								i_24 := j
								a_23[i_24] = t2.children[j]
							}
						}
					}
					copied = alloc
				}
				if c != nil {
					{
						a_25 := copied
						i_26 := func() int32 { old_27 := outIdx; nw_28 := old_27 + 1; outIdx = nw_28; return old_27 }()
						a_25[i_26] = c
					}
				} else {
					childMask = (childMask & ^(bit))
				}
			}
		}
		if copied == nil {
			if t2.present && func() bool {
				rk_29 := t2.value
				ak_30 := newValue
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					_ = j2gort.CheckNotNullString(any(rk_29).(string))
					return rk_29 == ak_30
				}
				var z_ V
				if rk_29 == z_ {
					j2gort.ThrowNPE()
				}
				if rk_29 == ak_30 {
					return true
				}
				return javalang.ObjectsEquals(j2gort.ToAny(rk_29), j2gort.ToAny(ak_30))
			}() {
				return t2
			} else {
				return NewHamtPMap[K, V](t2.key, t2.hash, newValue, true, childMask, t2.children)
			}
		}
		var newChildren []*HamtPMap[K, V] = HamtPMap_TrimChildren[K, V](copied, outIdx)
		_ = newChildren
		return NewHamtPMap[K, V](t2.key, t2.hash, newValue, true, childMask, newChildren)
	} else if t2 == nil {
		var newValue_2 V = func() V {
			tv_31 := joiner.Merge(t1.key, t1.value, func() V {
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					return any(j2gort.NullString).(V)
				}
				var z_ V
				return z_
			}())
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				if tv_31 == any(j2gort.NullString).(V) {
					j2gort.ThrowNPE()
				}
				return tv_31
			}
			var z_ V
			if tv_31 == z_ {
				j2gort.ThrowNPE()
			}
			return tv_31
		}()
		_ = newValue_2
		var childMask_2 int32 = t1.mask
		_ = childMask_2
		if t1.children == nil {
			if t1.present && func() bool {
				rk_32 := t1.value
				ak_33 := newValue_2
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					_ = j2gort.CheckNotNullString(any(rk_32).(string))
					return rk_32 == ak_33
				}
				var z_ V
				if rk_32 == z_ {
					j2gort.ThrowNPE()
				}
				if rk_32 == ak_33 {
					return true
				}
				return javalang.ObjectsEquals(j2gort.ToAny(rk_32), j2gort.ToAny(ak_33))
			}() {
				return t1
			} else {
				return NewHamtPMap[K, V](t1.key, t1.hash, newValue_2, true, int32(0), nil)
			}
		}
		var copied_2 []*HamtPMap[K, V] = nil
		_ = copied_2
		var rem_2 int32 = childMask_2
		_ = rem_2
		var outIdx_2 int32 = int32(0)
		_ = outIdx_2
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var len__2 int32 = int32(len(t1.children))
			_ = len__2
			for ; i_2 < len__2; i_2 = i_2 + 1 {
				var bit_2 int32 = j2gort.IntLowestOneBit(rem_2)
				_ = bit_2
				rem_2 = (rem_2 & ^(bit_2))
				var origChild_2 *HamtPMap[K, V] = t1.children[i_2]
				_ = origChild_2
				var c_2 *HamtPMap[K, V] = HamtPMap_Reconcile[K, V](origChild_2, nil, joiner)
				_ = c_2
				if copied_2 == nil {
					if c_2 == origChild_2 {
						outIdx_2 = outIdx_2 + 1
						continue
					}
					var alloc_2 []*HamtPMap[K, V] = make([]*HamtPMap[K, V], len__2)
					_ = alloc_2
					{
						var j_2 int32 = int32(0)
						_ = j_2
						for ; j_2 < i_2; j_2 = j_2 + 1 {
							{
								a_34 := alloc_2
								i_35 := j_2
								a_34[i_35] = t1.children[j_2]
							}
						}
					}
					copied_2 = alloc_2
				}
				if c_2 != nil {
					{
						a_36 := copied_2
						i_37 := func() int32 { old_38 := outIdx_2; nw_39 := old_38 + 1; outIdx_2 = nw_39; return old_38 }()
						a_36[i_37] = c_2
					}
				} else {
					childMask_2 = (childMask_2 & ^(bit_2))
				}
			}
		}
		if copied_2 == nil {
			if t1.present && func() bool {
				rk_40 := t1.value
				ak_41 := newValue_2
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					_ = j2gort.CheckNotNullString(any(rk_40).(string))
					return rk_40 == ak_41
				}
				var z_ V
				if rk_40 == z_ {
					j2gort.ThrowNPE()
				}
				if rk_40 == ak_41 {
					return true
				}
				return javalang.ObjectsEquals(j2gort.ToAny(rk_40), j2gort.ToAny(ak_41))
			}() {
				return t1
			} else {
				return NewHamtPMap[K, V](t1.key, t1.hash, newValue_2, true, childMask_2, t1.children)
			}
		}
		var newChildren_2 []*HamtPMap[K, V] = HamtPMap_TrimChildren[K, V](copied_2, outIdx_2)
		_ = newChildren_2
		return NewHamtPMap[K, V](t1.key, t1.hash, newValue_2, true, childMask_2, newChildren_2)
	}
	var sameChildrenAs1 bool = true
	_ = sameChildrenAs1
	var sameChildrenAs2 bool = true
	_ = sameChildrenAs2
	var hashCmp int32 = HamtPMap_CompareUnsigned(t1.hash, t2.hash)
	_ = hashCmp
	if ((((hashCmp != int32(0)) && (t1.mask == int32(0))) && (t2.mask == int32(0))) && t1.present) && t2.present {
		if hashCmp < int32(0) {
			var v1 V = func() V {
				tv_42 := joiner.Merge(t1.key, t1.value, func() V {
					if _, ok_ := any((*V)(nil)).(*string); ok_ {
						return any(j2gort.NullString).(V)
					}
					var z_ V
					return z_
				}())
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					if tv_42 == any(j2gort.NullString).(V) {
						j2gort.ThrowNPE()
					}
					return tv_42
				}
				var z_ V
				if tv_42 == z_ {
					j2gort.ThrowNPE()
				}
				return tv_42
			}()
			_ = v1
			var v2 V = func() V {
				tv_43 := joiner.Merge(t2.key, func() V {
					if _, ok_ := any((*V)(nil)).(*string); ok_ {
						return any(j2gort.NullString).(V)
					}
					var z_ V
					return z_
				}(), t2.value)
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					if tv_43 == any(j2gort.NullString).(V) {
						j2gort.ThrowNPE()
					}
					return tv_43
				}
				var z_ V
				if tv_43 == z_ {
					j2gort.ThrowNPE()
				}
				return tv_43
			}()
			_ = v2
			var child *HamtPMap[K, V] = NewHamtPMap[K, V](t2.key, HamtPMap_Shift(t2.hash), v2, true, int32(0), nil)
			_ = child
			var arr []*HamtPMap[K, V] = []*HamtPMap[K, V]{child}
			_ = arr
			return NewHamtPMap[K, V](t1.key, t1.hash, v1, true, javalang.Shl32(int32(1), int(HamtPMap_Bucket(t2.hash))), arr)
		} else {
			var v2_2 V = func() V {
				tv_44 := joiner.Merge(t2.key, func() V {
					if _, ok_ := any((*V)(nil)).(*string); ok_ {
						return any(j2gort.NullString).(V)
					}
					var z_ V
					return z_
				}(), t2.value)
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					if tv_44 == any(j2gort.NullString).(V) {
						j2gort.ThrowNPE()
					}
					return tv_44
				}
				var z_ V
				if tv_44 == z_ {
					j2gort.ThrowNPE()
				}
				return tv_44
			}()
			_ = v2_2
			var v1_2 V = func() V {
				tv_45 := joiner.Merge(t1.key, t1.value, func() V {
					if _, ok_ := any((*V)(nil)).(*string); ok_ {
						return any(j2gort.NullString).(V)
					}
					var z_ V
					return z_
				}())
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					if tv_45 == any(j2gort.NullString).(V) {
						j2gort.ThrowNPE()
					}
					return tv_45
				}
				var z_ V
				if tv_45 == z_ {
					j2gort.ThrowNPE()
				}
				return tv_45
			}()
			_ = v1_2
			var child_2 *HamtPMap[K, V] = NewHamtPMap[K, V](t1.key, HamtPMap_Shift(t1.hash), v1_2, true, int32(0), nil)
			_ = child_2
			var arr_2 []*HamtPMap[K, V] = []*HamtPMap[K, V]{child_2}
			_ = arr_2
			return NewHamtPMap[K, V](t2.key, t2.hash, v2_2, true, javalang.Shl32(int32(1), int(HamtPMap_Bucket(t1.hash))), arr_2)
		}
	}
	var key K = t1.key
	_ = key
	var hash int32 = t1.hash
	_ = hash
	if hashCmp < int32(0) {
		t2 = t2.priv_HamtPMap_vacateRoot_8009()
		sameChildrenAs2 = false
	} else if hashCmp > int32(0) {
		t1 = t1.priv_HamtPMap_vacateRoot_8009()
		sameChildrenAs1 = false
		key = t2.key
		hash = t2.hash
	} else if !(func() bool {
		rk_46 := t1.key
		ak_47 := t2.key
		if _, ok_ := any((*K)(nil)).(*string); ok_ {
			_ = j2gort.CheckNotNullString(any(rk_46).(string))
			return rk_46 == ak_47
		}
		var z_ K
		if rk_46 == z_ {
			j2gort.ThrowNPE()
		}
		if rk_46 == ak_47 {
			return true
		}
		return javalang.ObjectsEquals(j2gort.ToAny(rk_46), j2gort.ToAny(ak_47))
	}()) {
		t2 = t2.priv_HamtPMap_pivot_8007(t1.key, t1.hash)
		sameChildrenAs2 = false
	}
	var newValue_3 V = func() V {
		if _, ok_ := any((*V)(nil)).(*string); ok_ {
			return any(j2gort.NullString).(V)
		}
		var z_ V
		return z_
	}()
	_ = newValue_3
	if !(t1.present) {
		newValue_3 = func() V {
			tv_48 := joiner.Merge(t2.key, func() V {
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					return any(j2gort.NullString).(V)
				}
				var z_ V
				return z_
			}(), t2.value)
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				if tv_48 == any(j2gort.NullString).(V) {
					j2gort.ThrowNPE()
				}
				return tv_48
			}
			var z_ V
			if tv_48 == z_ {
				j2gort.ThrowNPE()
			}
			return tv_48
		}()
	} else if !(t2.present) {
		newValue_3 = func() V {
			tv_49 := joiner.Merge(t1.key, t1.value, func() V {
				if _, ok_ := any((*V)(nil)).(*string); ok_ {
					return any(j2gort.NullString).(V)
				}
				var z_ V
				return z_
			}())
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				if tv_49 == any(j2gort.NullString).(V) {
					j2gort.ThrowNPE()
				}
				return tv_49
			}
			var z_ V
			if tv_49 == z_ {
				j2gort.ThrowNPE()
			}
			return tv_49
		}()
	} else if func() bool {
		rk_50 := t1.value
		ak_51 := t2.value
		if _, ok_ := any((*V)(nil)).(*string); ok_ {
			_ = j2gort.CheckNotNullString(any(rk_50).(string))
			return rk_50 == ak_51
		}
		var z_ V
		if rk_50 == z_ {
			j2gort.ThrowNPE()
		}
		if rk_50 == ak_51 {
			return true
		}
		return javalang.ObjectsEquals(j2gort.ToAny(rk_50), j2gort.ToAny(ak_51))
	}() {
		newValue_3 = t1.value
	} else {
		newValue_3 = func() V {
			tv_52 := joiner.Merge(t1.key, t1.value, t2.value)
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				if tv_52 == any(j2gort.NullString).(V) {
					j2gort.ThrowNPE()
				}
				return tv_52
			}
			var z_ V
			if tv_52 == z_ {
				j2gort.ThrowNPE()
			}
			return tv_52
		}()
	}
	var origMask int32 = (t1.mask | t2.mask)
	_ = origMask
	var newMask int32 = origMask
	_ = newMask
	sameChildrenAs1 = func() bool { ba_53 := sameChildrenAs1; bb_54 := (newMask == t1.mask); return ba_53 && bb_54 }()
	sameChildrenAs2 = func() bool { ba_55 := sameChildrenAs2; bb_56 := (newMask == t2.mask); return ba_55 && bb_56 }()
	if origMask == int32(0) {
		if (sameChildrenAs1 && t1.present) && func() bool {
			rk_57 := t1.value
			ak_58 := newValue_3
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				_ = j2gort.CheckNotNullString(any(rk_57).(string))
				return rk_57 == ak_58
			}
			var z_ V
			if rk_57 == z_ {
				j2gort.ThrowNPE()
			}
			if rk_57 == ak_58 {
				return true
			}
			return javalang.ObjectsEquals(j2gort.ToAny(rk_57), j2gort.ToAny(ak_58))
		}() {
			return t1
		} else if (sameChildrenAs2 && t2.present) && func() bool {
			rk_59 := t2.value
			ak_60 := newValue_3
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				_ = j2gort.CheckNotNullString(any(rk_59).(string))
				return rk_59 == ak_60
			}
			var z_ V
			if rk_59 == z_ {
				j2gort.ThrowNPE()
			}
			if rk_59 == ak_60 {
				return true
			}
			return javalang.ObjectsEquals(j2gort.ToAny(rk_59), j2gort.ToAny(ak_60))
		}() {
			return t2
		}
		return NewHamtPMap[K, V](key, hash, newValue_3, true, int32(0), nil)
	}
	var totalSlots int32 = j2gort.IntBitCount(origMask)
	_ = totalSlots
	var baseChildren []*HamtPMap[K, V]
	_ = baseChildren
	if origMask == t1.mask {
		baseChildren = t1.children
	} else if origMask == t2.mask {
		baseChildren = t2.children
	} else {
		baseChildren = nil
	}
	var newChildren_3 []*HamtPMap[K, V] = nil
	_ = newChildren_3
	if baseChildren == nil {
		var alloc_3 []*HamtPMap[K, V] = make([]*HamtPMap[K, V], totalSlots)
		_ = alloc_3
		newChildren_3 = alloc_3
	}
	var mask int32 = origMask
	_ = mask
	var slotIdx int32 = int32(0)
	_ = slotIdx
	var index int32 = int32(0)
	_ = index
	for mask != int32(0) {
		var childBit int32 = j2gort.IntLowestOneBit(mask)
		_ = childBit
		mask = (mask & ^(childBit))
		var child1 *HamtPMap[K, V] = t1.priv_HamtPMap_getChild_8006(childBit)
		_ = child1
		var child2 *HamtPMap[K, V] = t2.priv_HamtPMap_getChild_8006(childBit)
		_ = child2
		var rec *HamtPMap[K, V] = HamtPMap_Reconcile[K, V](child1, child2, joiner)
		_ = rec
		sameChildrenAs1 = func() bool { ba_61 := sameChildrenAs1; bb_62 := (rec == child1); return ba_61 && bb_62 }()
		sameChildrenAs2 = func() bool { ba_63 := sameChildrenAs2; bb_64 := (rec == child2); return ba_63 && bb_64 }()
		if newChildren_3 == nil {
			if rec == baseChildren[slotIdx] {
				slotIdx = slotIdx + 1
				index = index + 1
				continue
			}
			var alloc_4 []*HamtPMap[K, V] = make([]*HamtPMap[K, V], totalSlots)
			_ = alloc_4
			{
				var j_3 int32 = int32(0)
				_ = j_3
				for ; j_3 < slotIdx; j_3 = j_3 + 1 {
					{
						a_65 := alloc_4
						i_66 := j_3
						a_65[i_66] = baseChildren[j_3]
					}
				}
			}
			newChildren_3 = alloc_4
		}
		slotIdx = slotIdx + 1
		if rec != nil {
			{
				a_67 := newChildren_3
				i_68 := func() int32 { old_69 := index; nw_70 := old_69 + 1; index = nw_70; return old_69 }()
				a_67[i_68] = rec
			}
		} else {
			newMask = (newMask & ^(childBit))
		}
	}
	if (sameChildrenAs1 && t1.present) && func() bool {
		rk_71 := t1.value
		ak_72 := newValue_3
		if _, ok_ := any((*V)(nil)).(*string); ok_ {
			_ = j2gort.CheckNotNullString(any(rk_71).(string))
			return rk_71 == ak_72
		}
		var z_ V
		if rk_71 == z_ {
			j2gort.ThrowNPE()
		}
		if rk_71 == ak_72 {
			return true
		}
		return javalang.ObjectsEquals(j2gort.ToAny(rk_71), j2gort.ToAny(ak_72))
	}() {
		return t1
	} else if (sameChildrenAs2 && t2.present) && func() bool {
		rk_73 := t2.value
		ak_74 := newValue_3
		if _, ok_ := any((*V)(nil)).(*string); ok_ {
			_ = j2gort.CheckNotNullString(any(rk_73).(string))
			return rk_73 == ak_74
		}
		var z_ V
		if rk_73 == z_ {
			j2gort.ThrowNPE()
		}
		if rk_73 == ak_74 {
			return true
		}
		return javalang.ObjectsEquals(j2gort.ToAny(rk_73), j2gort.ToAny(ak_74))
	}() {
		return t2
	}
	var finalChildren []*HamtPMap[K, V]
	_ = finalChildren
	if newChildren_3 == nil {
		finalChildren = baseChildren
	} else {
		finalChildren = HamtPMap_TrimChildren[K, V](newChildren_3, index)
	}
	return NewHamtPMap[K, V](key, hash, newValue_3, true, newMask, finalChildren)
}

// Java static method: HamtPMap.trimChildren(com.google.javascript.rhino.HamtPMap<K,V> @org.jspecify.annotations.Nullable [] children, int count)
func HamtPMap_TrimChildren[K comparable, V comparable](children []*HamtPMap[K, V], count int32) []*HamtPMap[K, V] {
	_ = children
	_ = count
	if (count == int32(0)) || (children == nil) {
		return nil
	}
	if count == int32(len(children)) {
		return children
	}
	var trimmed []*HamtPMap[K, V] = make([]*HamtPMap[K, V], count)
	_ = trimmed
	{
		var i int32 = int32(0)
		_ = i
		for ; i < count; i = i + 1 {
			{
				a_75 := trimmed
				i_76 := i
				a_75[i_76] = children[i]
			}
		}
	}
	return trimmed
}

// Java method: HamtPMap.equivalent(com.google.javascript.rhino.PMap<K,V> that, java.util.function.BiPredicate<V,V> equivalence)
func (this_ *HamtPMap[K, V]) Equivalent(that PMap[K, V], equivalence j2gort.Ext_java_util_function_BiPredicate) bool {
	_ = that
	_ = equivalence
	return HamtPMap_Equivalent[K, V](func() *HamtPMap[K, V] {
		if !(this_.IsEmpty()) {
			return this_
		}
		return nil
	}(), func() *HamtPMap[K, V] {
		if !(that.IsEmpty()) {
			return j2gort.CastPtr[HamtPMap[K, V]](that)
		}
		return nil
	}(), equivalence)
}

// Java static method: HamtPMap.equivalent(com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<K,V> t1, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<K,V> t2, java.util.function.BiPredicate<V,V> equivalence)
func HamtPMap_Equivalent[K comparable, V comparable](t1 *HamtPMap[K, V], t2 *HamtPMap[K, V], equivalence j2gort.Ext_java_util_function_BiPredicate) bool {
	_ = t1
	_ = t2
	_ = equivalence
	if t1 == t2 {
		return true
	} else if (t1 == nil) || (t2 == nil) {
		return false
	}
	if t1.hash != t2.hash {
		return false
	} else if !(func() bool {
		rk_77 := t1.key
		ak_78 := t2.key
		if _, ok_ := any((*K)(nil)).(*string); ok_ {
			_ = j2gort.CheckNotNullString(any(rk_77).(string))
			return rk_77 == ak_78
		}
		var z_ K
		if rk_77 == z_ {
			j2gort.ThrowNPE()
		}
		if rk_77 == ak_78 {
			return true
		}
		return javalang.ObjectsEquals(j2gort.ToAny(rk_77), j2gort.ToAny(ak_78))
	}()) {
		t2 = t2.priv_HamtPMap_pivot_8007(t1.key, t1.hash)
		if !(t2.present) {
			return false
		}
	}
	if !(equivalence.Test(j2gort.ToAny(t1.value), j2gort.ToAny(t2.value))) {
		return false
	}
	var mask int32 = (t1.mask | t2.mask)
	_ = mask
	for mask != int32(0) {
		var childBit int32 = j2gort.IntLowestOneBit(mask)
		_ = childBit
		mask = (mask & ^(childBit))
		if !(HamtPMap_Equivalent[K, V](t1.priv_HamtPMap_getChild_8006(childBit), t2.priv_HamtPMap_getChild_8006(childBit), equivalence)) {
			return false
		}
	}
	return true
}

// Java method: HamtPMap.index(int bit)
func (this_ *HamtPMap[K, V]) priv_HamtPMap_index_8005(bit int32) int32 {
	_ = bit
	return j2gort.IntBitCount((this_.mask & (bit - int32(1))))
}

// Java method: HamtPMap.getChild(int bit)
func (this_ *HamtPMap[K, V]) priv_HamtPMap_getChild_8006(bit int32) *HamtPMap[K, V] {
	_ = bit
	if (this_.mask & bit) != int32(0) {
		return this_.children[this_.priv_HamtPMap_index_8005(bit)]
	} else {
		return nil
	}
}

// Java static method: HamtPMap.hash(K key)
func HamtPMap_Hash[K comparable](key K) int32 {
	_ = key
	return j2gort.BoxInt_Reverse(func() int32 {
		hk_79 := key
		if _, ok_ := any((*K)(nil)).(*string); ok_ {
			return javalang.StringHashCode(j2gort.CheckNotNullString(any(hk_79).(string)))
		}
		return javalang.ObjectHashCode(j2gort.CheckNotNull(j2gort.ToAny(hk_79)))
	}())
}

// Java static method: HamtPMap.bucket(int hash)
func HamtPMap_Bucket(hash int32) int32 {
	_ = hash
	return javalang.Ushr32(hash, int(int32(28)))
}

// Java static method: HamtPMap.shift(int hash)
func HamtPMap_Shift(hash int32) int32 {
	_ = hash
	return javalang.Shl32(hash, int(int32(4)))
}

// Java static method: HamtPMap.unshift(int hash, int bucket)
func HamtPMap_Unshift(hash int32, bucket int32) int32 {
	_ = hash
	_ = bucket
	return ((javalang.Ushr32(hash, int(int32(4)))) | (javalang.Shl32(bucket, int(int32(28)))))
}

// Java static method: HamtPMap.compareUnsigned(int left, int right)
func HamtPMap_CompareUnsigned(left int32, right int32) int32 {
	_ = left
	_ = right
	var diff int32 = ((javalang.Ushr32(left, int(int32(2)))) - (javalang.Ushr32(right, int(int32(2)))))
	_ = diff
	if diff != int32(0) {
		return diff
	} else {
		return ((left & int32(3)) - (right & int32(3)))
	}
}

// Java method: HamtPMap.pivot(K key, int hash)
func (this_ *HamtPMap[K, V]) priv_HamtPMap_pivot_8007(key K, hash int32) *HamtPMap[K, V] {
	_ = key
	_ = hash
	return this_.priv_HamtPMap_pivot_8008(key, hash, nil, NewHamtPMap_Holder[V]())
}

// Java method: HamtPMap.pivot(K key, int hash, com.google.javascript.rhino.@org.jspecify.annotations.Nullable HamtPMap<K,V> parent, com.google.javascript.rhino.HamtPMap.Holder<V> result)
func (this_ *HamtPMap[K, V]) priv_HamtPMap_pivot_8008(key K, hash int32, parent *HamtPMap[K, V], result *HamtPMap_Holder[V]) *HamtPMap[K, V] {
	_ = key
	_ = hash
	_ = parent
	_ = result
	var newMask int32 = this_.mask
	_ = newMask
	var newChildren []*HamtPMap[K, V] = this_.children
	_ = newChildren
	if (hash == this_.hash) && func() bool {
		rk_80 := key
		ak_81 := this_.key
		if _, ok_ := any((*K)(nil)).(*string); ok_ {
			_ = j2gort.CheckNotNullString(any(rk_80).(string))
			return rk_80 == ak_81
		}
		var z_ K
		if rk_80 == z_ {
			j2gort.ThrowNPE()
		}
		if rk_80 == ak_81 {
			return true
		}
		return javalang.ObjectsEquals(j2gort.ToAny(rk_80), j2gort.ToAny(ak_81))
	}() {
		{
			r_82 := result
			r_82.found = true
		}
		{
			r_83 := result
			r_83.value = this_.value
		}
	} else {
		var searchBucket int32 = HamtPMap_Bucket(hash)
		_ = searchBucket
		var replacementBucket int32 = HamtPMap_Bucket(this_.hash)
		_ = replacementBucket
		if searchBucket == replacementBucket {
			var bucketMask int32 = javalang.Shl32(int32(1), int(searchBucket))
			_ = bucketMask
			if (this_.mask & bucketMask) != int32(0) {
				var index int32 = this_.priv_HamtPMap_index_8005(bucketMask)
				_ = index
				var child *HamtPMap[K, V] = newChildren[index]
				_ = child
				var newChild *HamtPMap[K, V] = child.priv_HamtPMap_pivot_8008(key, HamtPMap_Shift(hash), this_, result)
				_ = newChild
				newChildren = HamtPMap_ReplaceChild[K, V](newChildren, index, newChild)
			}
		} else {
			var searchMask int32 = javalang.Shl32(int32(1), int(searchBucket))
			_ = searchMask
			if (this_.mask & searchMask) != int32(0) {
				var index_2 int32 = this_.priv_HamtPMap_index_8005(searchMask)
				_ = index_2
				var child_2 *HamtPMap[K, V] = newChildren[index_2]
				_ = child_2
				var newChild_2 *HamtPMap[K, V] = child_2.priv_HamtPMap_minusOrNull_8001(key, HamtPMap_Shift(hash), result)
				_ = newChild_2
				if newChild_2 != nil {
					newChildren = HamtPMap_ReplaceChild[K, V](newChildren, index_2, newChild_2)
				} else {
					newChildren = HamtPMap_DeleteChild[K, V](newChildren, index_2)
					newMask = (newMask & ^(searchMask))
				}
			}
			var replacementMask int32 = javalang.Shl32(int32(1), int(replacementBucket))
			_ = replacementMask
			var index_3 int32 = j2gort.IntBitCount((newMask & (replacementMask - int32(1))))
			_ = index_3
			if (this_.mask & replacementMask) != int32(0) {
				var child_3 *HamtPMap[K, V] = newChildren[index_3]
				_ = child_3
				var newChild_3 *HamtPMap[K, V] = child_3.priv_HamtPMap_plus_7997(this_.key, HamtPMap_Shift(this_.hash), this_.value)
				_ = newChild_3
				newChildren = HamtPMap_ReplaceChild[K, V](newChildren, index_3, newChild_3)
			} else {
				newChildren = HamtPMap_InsertChild[K, V](newChildren, index_3, NewHamtPMap[K, V](this_.key, HamtPMap_Shift(this_.hash), this_.value, true, int32(0), nil))
				newMask = (newMask | replacementMask)
			}
		}
	}
	if parent != nil {
		return NewHamtPMap[K, V](parent.key, HamtPMap_Shift(parent.hash), parent.value, true, newMask, newChildren)
	} else {
		return NewHamtPMap[K, V](key, hash, result.value, result.found, newMask, newChildren)
	}
}

// Java method: HamtPMap.vacateRoot()
func (this_ *HamtPMap[K, V]) priv_HamtPMap_vacateRoot_8009() *HamtPMap[K, V] {
	var bucket int32 = HamtPMap_Bucket(this_.hash)
	_ = bucket
	var bucketMask int32 = javalang.Shl32(int32(1), int(bucket))
	_ = bucketMask
	var index int32 = this_.priv_HamtPMap_index_8005(bucketMask)
	_ = index
	if (this_.mask & bucketMask) != int32(0) {
		var newChild *HamtPMap[K, V] = this_.children[index].priv_HamtPMap_plus_7997(this_.key, HamtPMap_Shift(this_.hash), this_.value)
		_ = newChild
		return NewHamtPMap[K, V](func() K {
			if _, ok_ := any((*K)(nil)).(*string); ok_ {
				return any(j2gort.NullString).(K)
			}
			var z_ K
			return z_
		}(), int32(0), func() V {
			if _, ok_ := any((*V)(nil)).(*string); ok_ {
				return any(j2gort.NullString).(V)
			}
			var z_ V
			return z_
		}(), false, this_.mask, HamtPMap_ReplaceChild[K, V](this_.children, index, newChild))
	}
	var newChild_2 *HamtPMap[K, V] = NewHamtPMap[K, V](this_.key, HamtPMap_Shift(this_.hash), this_.value, true, int32(0), nil)
	_ = newChild_2
	return NewHamtPMap[K, V](func() K {
		if _, ok_ := any((*K)(nil)).(*string); ok_ {
			return any(j2gort.NullString).(K)
		}
		var z_ K
		return z_
	}(), int32(0), func() V {
		if _, ok_ := any((*V)(nil)).(*string); ok_ {
			return any(j2gort.NullString).(V)
		}
		var z_ V
		return z_
	}(), false, (this_.mask | bucketMask), HamtPMap_InsertChild[K, V](this_.children, index, newChild_2))
}

// Java method: HamtPMap.withChildren(int mask, com.google.javascript.rhino.HamtPMap<K,V> @org.jspecify.annotations.Nullable [] children)
func (this_ *HamtPMap[K, V]) priv_HamtPMap_withChildren_8010(mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	_ = mask
	_ = children
	if (mask == this_.mask) && j2gort.SliceRefEquals(children, this_.children) {
		return this_
	} else {
		return NewHamtPMap[K, V](this_.key, this_.hash, this_.value, this_.present, mask, children)
	}
}

// Java static method: HamtPMap.deleteRoot(int mask, com.google.javascript.rhino.HamtPMap<K,V> @org.jspecify.annotations.Nullable [] children)
func HamtPMap_DeleteRoot[K comparable, V comparable](mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	_ = mask
	_ = children
	if mask == int32(0) {
		return nil
	}
	var child *HamtPMap[K, V] = children[int32(0)]
	_ = child
	var hashBits int32 = j2gort.IntTrailingZeros(mask)
	_ = hashBits
	var newHash int32 = HamtPMap_Unshift(child.hash, hashBits)
	_ = newHash
	var newChild *HamtPMap[K, V] = HamtPMap_DeleteRoot[K, V](child.mask, child.children)
	_ = newChild
	if newChild == nil {
		var newMask int32 = (mask & ^(j2gort.IntLowestOneBit(mask)))
		_ = newMask
		return NewHamtPMap[K, V](child.key, newHash, child.value, true, newMask, HamtPMap_DeleteChild[K, V](children, int32(0)))
	} else {
		return NewHamtPMap[K, V](child.key, newHash, child.value, true, mask, HamtPMap_ReplaceChild[K, V](children, int32(0), newChild))
	}
}

// Java static method: HamtPMap.insertChild(com.google.javascript.rhino.HamtPMap<K,V> @org.jspecify.annotations.Nullable [] children, int index, com.google.javascript.rhino.HamtPMap<K,V> child)
func HamtPMap_InsertChild[K comparable, V comparable](children []*HamtPMap[K, V], index int32, child *HamtPMap[K, V]) []*HamtPMap[K, V] {
	_ = children
	_ = index
	_ = child
	var oldLen int32
	_ = oldLen
	if children == nil {
		oldLen = int32(0)
	} else {
		oldLen = int32(len(children))
	}
	var newChildren []*HamtPMap[K, V] = make([]*HamtPMap[K, V], (oldLen + int32(1)))
	_ = newChildren
	{
		var i int32 = int32(0)
		_ = i
		for ; i < index; i = i + 1 {
			{
				a_84 := newChildren
				i_85 := i
				a_84[i_85] = children[i]
			}
		}
	}
	{
		a_86 := newChildren
		i_87 := index
		a_86[i_87] = child
	}
	{
		var i_2 int32 = index
		_ = i_2
		for ; i_2 < oldLen; i_2 = i_2 + 1 {
			{
				a_88 := newChildren
				i_89 := (i_2 + int32(1))
				a_88[i_89] = children[i_2]
			}
		}
	}
	return newChildren
}

// Java static method: HamtPMap.replaceChild(com.google.javascript.rhino.HamtPMap<K,V>[] children, int index, com.google.javascript.rhino.HamtPMap<K,V> child)
func HamtPMap_ReplaceChild[K comparable, V comparable](children []*HamtPMap[K, V], index int32, child *HamtPMap[K, V]) []*HamtPMap[K, V] {
	_ = children
	_ = index
	_ = child
	var newChildren []*HamtPMap[K, V] = make([]*HamtPMap[K, V], int32(len(children)))
	_ = newChildren
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(children)); i = i + 1 {
			{
				a_90 := newChildren
				i_91 := i
				a_90[i_91] = children[i]
			}
		}
	}
	{
		a_92 := newChildren
		i_93 := index
		a_92[i_93] = child
	}
	return newChildren
}

// Java static method: HamtPMap.deleteChild(com.google.javascript.rhino.HamtPMap<K,V>[] children, int index)
func HamtPMap_DeleteChild[K comparable, V comparable](children []*HamtPMap[K, V], index int32) []*HamtPMap[K, V] {
	_ = children
	_ = index
	if int32(len(children)) == int32(1) {
		return nil
	}
	var newChildren []*HamtPMap[K, V] = make([]*HamtPMap[K, V], (int32(len(children)) - int32(1)))
	_ = newChildren
	{
		var i int32 = int32(0)
		_ = i
		for ; i < index; i = i + 1 {
			{
				a_94 := newChildren
				i_95 := i
				a_94[i_95] = children[i]
			}
		}
	}
	{
		var i_2 int32 = (index + int32(1))
		_ = i_2
		for ; i_2 < int32(len(children)); i_2 = i_2 + 1 {
			{
				a_96 := newChildren
				i_97 := (i_2 - int32(1))
				a_96[i_97] = children[i_2]
			}
		}
	}
	return newChildren
}

// Java method: HamtPMap.assertCorrectStructure()
func (this_ *HamtPMap[K, V]) AssertCorrectStructure() *HamtPMap[K, V] {
	if this_.IsEmpty() {
		return this_
	}
	var hash int32 = HamtPMap_Hash[K](this_.key)
	_ = hash
	if this_.children != nil {
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(this_.children)); i = i + 1 {
				var childHash int32 = HamtPMap_Hash[K](this_.children[i].key)
				_ = childHash
				if HamtPMap_CompareUnsigned(childHash, hash) < int32(0) {
					panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.BoxString(("Invalid map has decreasing hash " + j2gort.StrVal(j2gort.ToAny(this_.children[i].key)) + "(" + javalang.IntegerToHexString(childHash) + ") beneath " + j2gort.StrVal(j2gort.ToAny(this_.key)) + "(" + javalang.IntegerToHexString(hash) + ": " + j2gort.StrVal(this_))))))
				}
				this_.children[i].AssertCorrectStructure()
			}
		}
	}
	return this_
}

// Java class: com.google.javascript.rhino.HamtPMap.Holder
type HamtPMap_Holder[V comparable] struct {
	value   V `gson:"value"`
	idHash_ atomic.Int32
	found   bool `gson:"found"`
}

func (this_ *HamtPMap_Holder[V]) BindSelf_(s any) {
	_ = s
}
func (this_ *HamtPMap_Holder[V]) IsJavaNil_() bool { return this_ == nil }
func (*HamtPMap_Holder[V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.HamtPMap$Holder", "Holder"
}
func (this_ *HamtPMap_Holder[V]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*HamtPMap_Holder[V]) HasIdentityEquals_()                   {}
func (*HamtPMap_Holder[V]) IsGenericInst_rhino_HamtPMap_Holder_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.HamtPMap.Holder", "com.google.javascript.rhino.HamtPMap$Holder", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_rhino_HamtPMap_Holder_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

func (this_ *HamtPMap_Holder[V]) initFields_HamtPMap_Holder() {
	_ = this_
}

// Java constructor: Holder()
func NewHamtPMap_Holder[V comparable]() *HamtPMap_Holder[V] {
	this_ := &HamtPMap_Holder[V]{}
	InitHamtPMap_Holder[V](this_)
	return this_
}

func InitHamtPMap_Holder[V comparable](this_ *HamtPMap_Holder[V]) {
	this_.initFields_HamtPMap_Holder()
}

// Java class: com.google.javascript.rhino.HamtPMap.KeyIter
type HamtPMap_KeyIter[K comparable, V comparable] struct {
	queue   j2gort.Deque[*HamtPMap[K, V]] `gson:"queue"`
	idHash_ atomic.Int32
}

func (this_ *HamtPMap_KeyIter[K, V]) BindSelf_(s any) {
	_ = s
}
func (this_ *HamtPMap_KeyIter[K, V]) IsJavaNil_() bool { return this_ == nil }
func (*HamtPMap_KeyIter[K, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.HamtPMap$KeyIter", "KeyIter"
}
func (this_ *HamtPMap_KeyIter[K, V]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*HamtPMap_KeyIter[K, V]) HasIdentityEquals_()                    {}
func (*HamtPMap_KeyIter[K, V]) IsGenericInst_rhino_HamtPMap_KeyIter_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.HamtPMap.KeyIter", "com.google.javascript.rhino.HamtPMap$KeyIter", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_rhino_HamtPMap_KeyIter_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*HamtPMap_KeyIter[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}

func (this_ *HamtPMap_KeyIter[K, V]) initFields_HamtPMap_KeyIter() {
	_ = this_
	this_.queue = j2gort.Deque[*HamtPMap[K, V]]{}
}

// Java constructor: KeyIter(com.google.javascript.rhino.HamtPMap<K,V> map)
func NewHamtPMap_KeyIter[K comparable, V comparable](map_ *HamtPMap[K, V]) *HamtPMap_KeyIter[K, V] {
	this_ := &HamtPMap_KeyIter[K, V]{}
	InitHamtPMap_KeyIter[K, V](this_, map_)
	return this_
}

func InitHamtPMap_KeyIter[K comparable, V comparable](this_ *HamtPMap_KeyIter[K, V], map_ *HamtPMap[K, V]) {
	_ = map_
	this_.initFields_HamtPMap_KeyIter()
	if !(map_.IsEmpty()) {
		this_.queue.AddLast(map_)
	}
}

// Java method: KeyIter.hasNext()
func (this_ *HamtPMap_KeyIter[K, V]) HasNext() bool {
	return !(this_.queue.IsEmpty())
}

// Java method: KeyIter.next()
func (this_ *HamtPMap_KeyIter[K, V]) Next() any {
	var top *HamtPMap[K, V] = this_.queue.RemoveFirst()
	_ = top
	if top.children != nil {
		{
			var i int32 = (int32(len(top.children)) - int32(1))
			_ = i
			for ; i >= int32(0); i = i - 1 {
				this_.queue.AddLast(top.children[i])
			}
		}
	}
	return j2gort.ToAny(top.key)
}

// Java method: KeyIter.remove()
func (this_ *HamtPMap_KeyIter[K, V]) Remove() {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java class: com.google.javascript.rhino.HamtPMap.ValueIter
type HamtPMap_ValueIter[K comparable, V comparable] struct {
	queue   j2gort.Deque[*HamtPMap[K, V]] `gson:"queue"`
	idHash_ atomic.Int32
}

func (this_ *HamtPMap_ValueIter[K, V]) BindSelf_(s any) {
	_ = s
}
func (this_ *HamtPMap_ValueIter[K, V]) IsJavaNil_() bool { return this_ == nil }
func (*HamtPMap_ValueIter[K, V]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.HamtPMap$ValueIter", "ValueIter"
}
func (this_ *HamtPMap_ValueIter[K, V]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*HamtPMap_ValueIter[K, V]) HasIdentityEquals_()                      {}
func (*HamtPMap_ValueIter[K, V]) IsGenericInst_rhino_HamtPMap_ValueIter_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.HamtPMap.ValueIter", "com.google.javascript.rhino.HamtPMap$ValueIter", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_rhino_HamtPMap_ValueIter_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*HamtPMap_ValueIter[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}

func (this_ *HamtPMap_ValueIter[K, V]) initFields_HamtPMap_ValueIter() {
	_ = this_
	this_.queue = j2gort.Deque[*HamtPMap[K, V]]{}
}

// Java constructor: ValueIter(com.google.javascript.rhino.HamtPMap<K,V> map)
func NewHamtPMap_ValueIter[K comparable, V comparable](map_ *HamtPMap[K, V]) *HamtPMap_ValueIter[K, V] {
	this_ := &HamtPMap_ValueIter[K, V]{}
	InitHamtPMap_ValueIter[K, V](this_, map_)
	return this_
}

func InitHamtPMap_ValueIter[K comparable, V comparable](this_ *HamtPMap_ValueIter[K, V], map_ *HamtPMap[K, V]) {
	_ = map_
	this_.initFields_HamtPMap_ValueIter()
	if !(map_.IsEmpty()) {
		this_.queue.AddLast(map_)
	}
}

// Java method: ValueIter.hasNext()
func (this_ *HamtPMap_ValueIter[K, V]) HasNext() bool {
	return !(this_.queue.IsEmpty())
}

// Java method: ValueIter.next()
func (this_ *HamtPMap_ValueIter[K, V]) Next() any {
	var top *HamtPMap[K, V] = this_.queue.RemoveFirst()
	_ = top
	if top.children != nil {
		{
			var i int32 = (int32(len(top.children)) - int32(1))
			_ = i
			for ; i >= int32(0); i = i - 1 {
				this_.queue.AddLast(top.children[i])
			}
		}
	}
	return j2gort.ToAny(top.value)
}

// Java method: ValueIter.remove()
func (this_ *HamtPMap_ValueIter[K, V]) Remove() {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

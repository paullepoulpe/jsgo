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
// Port of com/google/javascript/rhino/RhinoStringPool.java (generated mechanically by j2go).

package rhino

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

// Java class: com.google.javascript.rhino.RhinoStringPool
type RhinoStringPool struct {
	idHash_ atomic.Int32
}

func (this_ *RhinoStringPool) BindSelf_(s any) {
	_ = s
}
func (this_ *RhinoStringPool) IsJavaNil_() bool { return this_ == nil }
func (*RhinoStringPool) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.RhinoStringPool", "RhinoStringPool" }
func (this_ *RhinoStringPool) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RhinoStringPool) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.RhinoStringPool", "com.google.javascript.rhino.RhinoStringPool", func(v_ any) bool {
		return j2gort.InstanceOf[*RhinoStringPool](v_)
	})
}

func (this_ *RhinoStringPool) initFields_RhinoStringPool() {
	_ = this_
}

// Java constructor: RhinoStringPool()
func NewRhinoStringPool() *RhinoStringPool {
	this_ := &RhinoStringPool{}
	InitRhinoStringPool(this_)
	return this_
}

func InitRhinoStringPool(this_ *RhinoStringPool) {
	this_.initFields_RhinoStringPool()
}

// Java static method: RhinoStringPool.uncheckedEquals(java.lang.@org.jspecify.annotations.Nullable String a, java.lang.@org.jspecify.annotations.Nullable String b)
func RhinoStringPool_UncheckedEquals(a string, b string) bool {
	_ = a
	_ = b
	return ((a == b) || (((a != j2gort.NullString) && (j2gort.CheckNotNullString(a) == b))))
}

// Java static method: RhinoStringPool.addOrGet(java.lang.String s)
func RhinoStringPool_AddOrGet(s string) string {
	_ = s
	return j2gort.CheckNotNullString(s)
}

// Java static method: RhinoStringPool.fastHash32(java.lang.String s)
func RhinoStringPool_FastHash32(s string) int32 {
	_ = s
	return j2gort.FastStrHash32(s)
}

// Java class: com.google.javascript.rhino.RhinoStringPool.LazyInternedStringList
type RhinoStringPool_LazyInternedStringList struct {
	pool []string `gson:"pool"`
	ids_ []uint32 `gson:"-"`
	epoch_ uint64 `gson:"-"`
	size int32 `gson:"size"`
	idHash_ atomic.Int32
}

func (this_ *RhinoStringPool_LazyInternedStringList) BindSelf_(s any) {
	_ = s
}
func (this_ *RhinoStringPool_LazyInternedStringList) IsJavaNil_() bool { return this_ == nil }
func (*RhinoStringPool_LazyInternedStringList) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.RhinoStringPool$LazyInternedStringList", "LazyInternedStringList" }
func (this_ *RhinoStringPool_LazyInternedStringList) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RhinoStringPool_LazyInternedStringList) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.RhinoStringPool.LazyInternedStringList", "com.google.javascript.rhino.RhinoStringPool$LazyInternedStringList", func(v_ any) bool {
		return j2gort.InstanceOf[*RhinoStringPool_LazyInternedStringList](v_)
	})
}
var lazyStringPoolEpochMu_ sync.Mutex

func (this_ *RhinoStringPool_LazyInternedStringList) InternID_(offset int32) nodestore.StringID {
	return internIDFast_(this_, offset, (*RhinoStringPool_LazyInternedStringList).internIDSlow_)
}

func internIDFast_(this_ *RhinoStringPool_LazyInternedStringList, offset int32, slow func(*RhinoStringPool_LazyInternedStringList, int32) nodestore.StringID) nodestore.StringID {
	if this_ != nil {
		if ids := this_.ids_; uint(offset) < uint(len(ids)) && atomic.LoadUint64(&this_.epoch_) == nodestore.Epoch() {
			if id := atomic.LoadUint32(&ids[offset]); id != 0 {
				return nodestore.StringID(id)
			}
		}
	}
	return slow(this_, offset)
}

//go:noinline
func (this_ *RhinoStringPool_LazyInternedStringList) internIDSlow_(offset int32) nodestore.StringID {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	if offset < 0 || offset >= this_.size {
		_ = this_.Get(offset)
	}
	ep := nodestore.Epoch()
	if atomic.LoadUint64(&this_.epoch_) != ep {
		lazyStringPoolEpochMu_.Lock()
		if ep = nodestore.Epoch(); atomic.LoadUint64(&this_.epoch_) != ep {
			for i := range this_.ids_ {
				atomic.StoreUint32(&this_.ids_[i], 0)
			}
			atomic.StoreUint64(&this_.epoch_, ep)
		}
		lazyStringPoolEpochMu_.Unlock()
	}
	s := this_.pool[offset]
	if s == j2gort.NullString {
		return 0
	}
	id := nodestore.Intern(s)
	if int(offset) < len(this_.ids_) {
		atomic.StoreUint32(&this_.ids_[offset], uint32(id))
	}
	return id
}

func (this_ *RhinoStringPool_LazyInternedStringList) initFields_RhinoStringPool_LazyInternedStringList() {
	_ = this_
}

// Java constructor: LazyInternedStringList(java.util.List<java.lang.String> src)
func NewRhinoStringPool_LazyInternedStringList_1(src j2gort.Ext_java_util_List) *RhinoStringPool_LazyInternedStringList {
	this_ := &RhinoStringPool_LazyInternedStringList{}
	InitRhinoStringPool_LazyInternedStringList_1(this_, src)
	return this_
}

func InitRhinoStringPool_LazyInternedStringList_1(this_ *RhinoStringPool_LazyInternedStringList, src j2gort.Ext_java_util_List) {
	_ = src
	this_.initFields_RhinoStringPool_LazyInternedStringList()
	var n int32 = src.Size()
	_ = n
	var copy_ []string = j2gort.NewStringArray(n)
	_ = copy_
	{
		var i int32 = int32(0)
		_ = i
		for ; (i < n); i = i + 1 {
			{ a_1 := copy_; i_2 := i; a_1[i_2] = j2gort.CheckNotNullString(j2gort.ListGetStr(src, i)) }
		}
	}
	this_.pool = copy_
	this_.size = n
	if this_.size > 0 {
		this_.ids_ = make([]uint32, int(this_.size))
		atomic.StoreUint64(&this_.epoch_, nodestore.Epoch())
	}
}

// Java constructor: LazyInternedStringList(java.lang.String[] trustedArray, int size)
func NewRhinoStringPool_LazyInternedStringList_2(trustedArray []string, size int32) *RhinoStringPool_LazyInternedStringList {
	this_ := &RhinoStringPool_LazyInternedStringList{}
	InitRhinoStringPool_LazyInternedStringList_2(this_, trustedArray, size)
	return this_
}

func InitRhinoStringPool_LazyInternedStringList_2(this_ *RhinoStringPool_LazyInternedStringList, trustedArray []string, size int32) {
	_ = trustedArray
	_ = size
	this_.initFields_RhinoStringPool_LazyInternedStringList()
	this_.pool = trustedArray
	this_.size = size
	if this_.size > 0 {
		this_.ids_ = make([]uint32, int(this_.size))
		atomic.StoreUint64(&this_.epoch_, nodestore.Epoch())
	}
}

// Java method: LazyInternedStringList.size()
func (this_ *RhinoStringPool_LazyInternedStringList) Size() int32 {
	return this_.size
}

// Java method: LazyInternedStringList.get(int offset)
func (this_ *RhinoStringPool_LazyInternedStringList) Get(offset int32) string {
	_ = offset
	if (((offset < int32(0)) || (offset >= this_.size))) {
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.BoxString(("Index: " + j2gort.StrValInt32(offset) + ", Size: " + j2gort.StrValInt32(this_.size))))))
	}
	return this_.pool[offset]
}


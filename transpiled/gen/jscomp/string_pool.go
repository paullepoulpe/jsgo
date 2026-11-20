/*
 * Copyright 2021 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/serialization/StringPool.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

// Java class: com.google.javascript.jscomp.serialization.StringPool
type StringPool struct {
	pool         *pkg_rhino.RhinoStringPool_LazyInternedStringList `gson:"pool"`
	maxLength    int32                                             `gson:"maxLength"`
	idHash_      atomic.Int32
	internedIDs_ []atomic.Uint32
	internEpoch_ atomic.Uint64
}

func (this_ *StringPool) BindSelf_(s any) {
	_ = s
}
func (this_ *StringPool) IsJavaNil_() bool { return this_ == nil }
func (*StringPool) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.StringPool", "StringPool"
}
func (this_ *StringPool) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*StringPool) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.StringPool", "com.google.javascript.jscomp.serialization.StringPool", func(v_ any) bool {
		return j2gort.InstanceOf[*StringPool](v_)
	})
}

var static_StringPool_EMPTY *StringPool

func StaticGet_StringPool_EMPTY() *StringPool {
	clinit_StringPool()
	return static_StringPool_EMPTY
}
func StaticSet_StringPool_EMPTY(v *StringPool) {
	clinit_StringPool()
	static_StringPool_EMPTY = v
}

var clinitDone_StringPool uint32
var clinitState_StringPool j2gort.ClinitState

func clinit_StringPool() {
	if atomic.LoadUint32(&clinitDone_StringPool) != 0 {
		return
	}
	clinitSlow_StringPool()
}

func clinitSlow_StringPool() {
	if !j2gort.EnterClinit(&clinitState_StringPool, &clinitDone_StringPool) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_StringPool, &clinitDone_StringPool)
	static_StringPool_EMPTY = NewStringPool(int32(0), []string{""}, int32(1))
}

func (this_ *StringPool) initFields_StringPool() {
	_ = this_
}

// Java constructor: StringPool(int maxLength, java.lang.String[] pool, int size)
func NewStringPool(maxLength int32, pool []string, size int32) *StringPool {
	clinit_StringPool()
	this_ := &StringPool{}
	InitStringPool(this_, maxLength, pool, size)
	return this_
}

func InitStringPool(this_ *StringPool, maxLength int32, pool []string, size int32) {
	_ = maxLength
	_ = pool
	_ = size
	this_.initFields_StringPool()
	this_.maxLength = maxLength
	this_.pool = pkg_rhino.NewRhinoStringPool_LazyInternedStringList_2(pool, size)
	if !((size > int32(0)) && (j2gort.CheckNotNullString(pool[int32(0)]) == "")) {
		j2gort.CheckState(false)
	}
	if this_.pool != nil {
		this_.internedIDs_ = make([]atomic.Uint32, this_.pool.Size())
		this_.internEpoch_.Store(nodestore.Epoch())
	}

}

// Java static method: StringPool.fromProto(com.google.javascript.jscomp.serialization.StringPoolProto proto)
func StringPool_FromProto(proto *j2gort.Ext_com_google_javascript_jscomp_serialization_StringPoolProto) *StringPool {
	clinit_StringPool()
	_ = proto
	var decoder *Wtf8_Decoder = Wtf8_S_Decoder(proto.GetMaxLength())
	_ = decoder
	var stringsList j2gort.Ext_java_util_List = proto.GetStringsList()
	_ = stringsList
	var size int32 = stringsList.Size()
	_ = size
	var pool []string = j2gort.NewStringArray((size + int32(1)))
	_ = pool
	{
		a_1 := pool
		i_2 := int32(0)
		a_1[i_2] = ""
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < size; i = i + 1 {
			var s j2gort.Ext_com_google_protobuf_ByteString = j2gort.ListGetByteString(stringsList, i)
			_ = s
			{
				a_3 := pool
				i_4 := (i + int32(1))
				a_3[i_4] = j2gort.CheckNotNullString(decoder.Decode(s))
			}
		}
	}
	return NewStringPool(proto.GetMaxLength(), pool, (size + int32(1)))
}

// Java static method: StringPool.empty()
func StringPool_Empty() *StringPool {
	clinit_StringPool()
	return StaticGet_StringPool_EMPTY()
}

// Java method: StringPool.get(int offset)
func (this_ *StringPool) Get(offset int32) string {
	_ = offset
	return this_.pool.Get(offset)
}

// Java method: StringPool.getInternedStrings()
func (this_ *StringPool) GetInternedStrings() *pkg_rhino.RhinoStringPool_LazyInternedStringList {
	return this_.pool
}

// Java method: StringPool.toProto()
func (this_ *StringPool) ToProto() *j2gort.Ext_com_google_javascript_jscomp_serialization_StringPoolProto {
	var builder *j2gort.Ext_com_google_javascript_jscomp_serialization_StringPoolProto_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_StringPoolProto_NewBuilder_0().SetMaxLength(this_.maxLength)
	_ = builder
	{
		var i int32 = int32(1)
		_ = i
		var n int32 = this_.pool.Size()
		_ = n
		for ; i < n; i = i + 1 {
			builder.AddStrings(Wtf8_EncodeToWtf8(this_.pool.Get(i)))
		}
	}
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_StringPoolProto](builder.Build())
}

// Java static method: StringPool.builder()
func StringPool_S_Builder() *StringPool_Builder {
	clinit_StringPool()
	return NewStringPool_Builder(int32(1024), int32(2048))
}

// Java static method: StringPool.smallBuilder()
func StringPool_SmallBuilder() *StringPool_Builder {
	clinit_StringPool()
	return NewStringPool_Builder(int32(8), int32(16))
}

// Java class: com.google.javascript.jscomp.serialization.StringPool.Builder
type StringPool_Builder struct {
	pool            []string `gson:"pool"`
	table           []int64  `gson:"table"`
	idToPool        []int32  `gson:"idToPool"`
	initialPoolCap  int32    `gson:"initialPoolCap"`
	initialTableCap int32    `gson:"initialTableCap"`
	maxLength       int32    `gson:"maxLength"`
	size            int32    `gson:"size"`
	idHash_         atomic.Int32
	built           bool `gson:"built"`
}

func (this_ *StringPool_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *StringPool_Builder) IsJavaNil_() bool { return this_ == nil }
func (*StringPool_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.StringPool$Builder", "Builder"
}
func (this_ *StringPool_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*StringPool_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.StringPool.Builder", "com.google.javascript.jscomp.serialization.StringPool$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*StringPool_Builder](v_)
	})
}

func (this_ *StringPool_Builder) initFields_StringPool_Builder() {
	_ = this_

	this_.size = int32(1)

}

// Java constructor: Builder(int initialPoolCap, int initialTableCap)
func NewStringPool_Builder(initialPoolCap int32, initialTableCap int32) *StringPool_Builder {
	this_ := &StringPool_Builder{}
	InitStringPool_Builder(this_, initialPoolCap, initialTableCap)
	return this_
}

func InitStringPool_Builder(this_ *StringPool_Builder, initialPoolCap int32, initialTableCap int32) {
	_ = initialPoolCap
	_ = initialTableCap
	this_.initFields_StringPool_Builder()
	this_.initialPoolCap = initialPoolCap
	this_.initialTableCap = initialTableCap
}

// Java method: Builder.putNodeString(com.google.javascript.rhino.Node n)
func (this_ *StringPool_Builder) PutNodeString(n pkg_rhino.Node) int32 {
	_ = n
	var id int32 = pkg_rhino.RhinoStringPool_GetNodeStringId(n)
	_ = id
	if id == int32(0) {
		return this_.Put(n.GetString())
	}
	if id == int32(1) {
		return int32(0)
	}
	return this_.PutByInternedId(id)
}

// Java method: Builder.putByInternedId(int id)
func (this_ *StringPool_Builder) PutByInternedId(id int32) int32 {
	_ = id
	if id <= int32(1) {
		return int32(0)
	}
	var map_ []int32 = this_.idToPool
	_ = map_
	if map_ == nil {
		var cap_ int32 = int32(4096)
		_ = cap_
		for cap_ <= id {
			cap_ = javalang.Shl32(cap_, int(int32(1)))
		}
		map_ = make([]int32, cap_)
		this_.idToPool = map_
	} else if id < int32(len(map_)) {
		var cached int32 = map_[id]
		_ = cached
		if cached != int32(0) {
			return (cached - int32(1))
		}
	} else {
		var newCap int32 = javalang.Shl32(int32(len(map_)), int(int32(1)))
		_ = newCap
		for newCap <= id {
			newCap = javalang.Shl32(newCap, int(int32(1)))
		}
		{

			arrSrc_5 := map_
			arrLen_6 := newCap
			j2gort.CheckNotNullSlice(arrSrc_5)
			if arrLen_6 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_7 := make([]int32, arrLen_6)
			copy(arrOut_7, arrSrc_5)
			map_ = arrOut_7
		}

		this_.idToPool = map_
	}
	var ptr int32 = this_.Put(j2gort.CheckNotNullString(pkg_rhino.RhinoStringPool_StringById(id)))
	_ = ptr
	{
		a_8 := map_
		i_9 := id
		a_8[i_9] = (ptr + int32(1))
	}
	return ptr
}

// Java method: Builder.put(java.lang.String string)
func (this_ *StringPool_Builder) Put(string_ string) int32 {
	_ = string_
	j2gort.CheckNotNullString(string_)
	if j2gort.CheckNotNullString(string_) == "" {
		return int32(0)
	}
	var hash int32 = j2gort.FastStrHash32(string_)
	_ = hash
	var poolArr []string = this_.pool
	_ = poolArr
	var tbl []int64 = this_.table
	_ = tbl
	var mask int32 = int32(0)
	_ = mask
	var slot int32 = int32(0)
	_ = slot
	if poolArr == nil {
		poolArr = j2gort.NewStringArray(this_.initialPoolCap)
		{
			a_10 := poolArr
			i_11 := int32(0)
			a_10[i_11] = ""
		}
		this_.pool = poolArr
		tbl = make([]int64, this_.initialTableCap)
		this_.table = tbl
		mask = (int32(len(tbl)) - int32(1))
		slot = (hash & mask)
	} else {
		mask = (int32(len(tbl)) - int32(1))
		slot = (hash & mask)
		var entry int64 = int64(0)
		_ = entry
		for (func() int64 { var asgn_12 int64 = tbl[slot]; entry = asgn_12; return asgn_12 }()) != int64(0) {
			if (int32((javalang.Ushr64(entry, int(int32(32)))))) == hash {
				var idx int32 = ((int32(entry)) - int32(1))
				_ = idx
				if j2gort.CheckNotNullString(poolArr[idx]) == string_ {
					return idx
				}
			}
			slot = ((slot + int32(1)) & mask)
		}
		if this_.built {
			{

				arrSrc_13 := poolArr
				arrLen_14 := max(int32(len(poolArr)), (this_.size + int32(16)))
				j2gort.CheckNotNullSlice(arrSrc_13)
				if arrLen_14 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_15 := j2gort.NewStringArray(arrLen_14)
				copy(arrOut_15, arrSrc_13)
				poolArr = arrOut_15
			}

			this_.pool = poolArr
			this_.built = false
		} else if this_.size == int32(len(poolArr)) {
			{

				arrSrc_16 := poolArr
				arrLen_17 := (int32(len(poolArr)) * int32(2))
				j2gort.CheckNotNullSlice(arrSrc_16)
				if arrLen_17 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_18 := j2gort.NewStringArray(arrLen_17)
				copy(arrOut_18, arrSrc_16)
				poolArr = arrOut_18
			}

			this_.pool = poolArr
		}
	}
	var len_ int32 = j2gort.JStrLen(string_)
	_ = len_
	if len_ > this_.maxLength {
		this_.maxLength = len_
	}
	var id int32 = this_.size
	_ = id
	{
		a_19 := poolArr
		i_20 := id
		a_19[i_20] = string_
	}
	var newSize int32 = (id + int32(1))
	_ = newSize
	this_.size = newSize
	if (javalang.Shl32(newSize, int(int32(2)))) > (int32(len(tbl)) * int32(3)) {
		tbl = this_.priv_Builder_rehash_7828(javalang.Shl32(int32(len(tbl)), int(int32(1))))
		mask = (int32(len(tbl)) - int32(1))
		slot = (hash & mask)
		for tbl[slot] != int64(0) {
			slot = ((slot + int32(1)) & mask)
		}
	}
	{
		a_21 := tbl
		i_22 := slot
		a_21[i_22] = ((javalang.Shl64((int64(hash)), int(int32(32)))) | (int64(newSize) & int64(4294967295)))
	}
	return id
}

// Java method: Builder.rehash(int newCap)
func (this_ *StringPool_Builder) priv_Builder_rehash_7828(newCap int32) []int64 {
	_ = newCap
	var oldTbl []int64 = this_.table
	_ = oldTbl
	var newTbl []int64 = make([]int64, newCap)
	_ = newTbl
	var newMask int32 = (newCap - int32(1))
	_ = newMask
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(oldTbl))
		_ = len_
		for ; i < len_; i = i + 1 {
			var entry int64 = oldTbl[i]
			_ = entry
			if entry != int64(0) {
				var slot int32 = ((int32((javalang.Ushr64(entry, int(int32(32)))))) & newMask)
				_ = slot
				for newTbl[slot] != int64(0) {
					slot = ((slot + int32(1)) & newMask)
				}
				{
					a_23 := newTbl
					i_24 := slot
					a_23[i_24] = entry
				}
			}
		}
	}
	this_.table = newTbl
	return newTbl
}

// Java method: Builder.putAnd(java.lang.String string)
func (this_ *StringPool_Builder) PutAnd(string_ string) *StringPool_Builder {
	_ = string_
	this_.Put(string_)
	return this_
}

// Java method: Builder.build()
func (this_ *StringPool_Builder) Build() *StringPool {
	if this_.pool == nil {
		return StaticGet_StringPool_EMPTY()
	}
	this_.built = true
	return NewStringPool(this_.maxLength, this_.pool, this_.size)
}
func (this_ *StringPool) internAt_(offset int32, alloc_ *nodestore.NodeAllocator) (nodestore.StringID, bool) {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	if uint32(offset) < uint32(len(this_.internedIDs_)) && this_.internEpoch_.Load() == nodestore.Epoch() {
		if id_ := this_.internedIDs_[offset].Load(); id_ != 0 {
			return nodestore.StringID(id_), true
		}
		s_ := this_.pool.Get(offset)
		if s_ == j2gort.NullString {
			return 0, false
		}
		id_ := nodestore.InternOwned(s_)
		this_.internedIDs_[offset].Store(uint32(id_))
		return id_, true
	}
	s_ := this_.pool.Get(offset)
	if s_ == j2gort.NullString {
		return 0, false
	}
	return nodestore.InternOwned(s_), true
}

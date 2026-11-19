/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/colors/ColorId.java (generated mechanically by j2go).

package colors

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_color_id_1 ColorId_ByteAt
var lambdaSingleton_color_id_2 ColorId_ByteAt
var lambdaSingleton_color_id_3 ColorId_ByteAt

func init() {
	lambdaSingleton_color_id_1 = NewLambda_ColorId_ByteAt(func(lp_1 any, i int32) int8 {
		s := j2gort.CastString(lp_1)
		_ = s
		_ = i
		var c int32 = j2gort.JStrCodePointAt(s, i)
		_ = c
		if !((int32(0) <= c) && (c < int32(128))) {
			j2gort.CheckStateMsg(false, j2gort.ValueOfInt(c))
		}
		return int8(c)
	})
	lambdaSingleton_color_id_2 = NewLambda_ColorId_ByteAt(func(mrp0_ any, mrp1_ int32) int8 {
		return j2gort.Cast[j2gort.Ext_com_google_protobuf_ByteString](mrp0_).ByteAt(mrp1_)
	})
	lambdaSingleton_color_id_3 = NewLambda_ColorId_ByteAt(func(lp_2 any, i int32) int8 {
		s := j2gort.Cast[[]int8](lp_2)
		_ = s
		_ = i
		return s[i]
	})
}

// Java class: com.google.javascript.jscomp.colors.ColorId
type ColorId struct {
	rightAligned int64 `gson:"rightAligned"`
	valid_       bool  `gson:"-"`
}

func (this_ ColorId) BindSelf_(s any) {
	_ = s
}
func (this_ ColorId) IsJavaNil_() bool { return !this_.valid_ }
func (ColorId) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.colors.ColorId", "ColorId"
}
func (this_ ColorId) FastValHash_() (int32, bool, bool) {
	if !this_.valid_ {
		return 0, false, false
	}
	return int32(this_.rightAligned), true, false
}
func (this_ ColorId) RawColorId_() (int64, bool) { return this_.rightAligned, this_.valid_ }
func init() {
	j2gort.RegisterColorIdBoxer(func(raw_ int64) any { return ColorId{rightAligned: raw_, valid_: true} })
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.colors.ColorId", "com.google.javascript.jscomp.colors.ColorId", func(v_ any) bool {
		return j2gort.InstanceOf[ColorId](v_)
	})
}
func (ColorId) Implements_diagnostic_LogsGson_Able() {}

var static_ColorId_FARM_64 j2gort.Ext_com_google_common_hash_HashFunction

func StaticGet_ColorId_FARM_64() j2gort.Ext_com_google_common_hash_HashFunction {
	clinit_ColorId()
	return static_ColorId_FARM_64
}
func StaticSet_ColorId_FARM_64(v j2gort.Ext_com_google_common_hash_HashFunction) {
	clinit_ColorId()
	static_ColorId_FARM_64 = v
}

var static_ColorId_ASCII_AT ColorId_ByteAt

func StaticGet_ColorId_ASCII_AT() ColorId_ByteAt {
	clinit_ColorId()
	return static_ColorId_ASCII_AT
}
func StaticSet_ColorId_ASCII_AT(v ColorId_ByteAt) {
	clinit_ColorId()
	static_ColorId_ASCII_AT = v
}

var clinitDone_ColorId uint32
var clinitState_ColorId j2gort.ClinitState

func clinit_ColorId() {
	if atomic.LoadUint32(&clinitDone_ColorId) != 0 {
		return
	}
	clinitSlow_ColorId()
}

func clinitSlow_ColorId() {
	if !j2gort.EnterClinit(&clinitState_ColorId, &clinitDone_ColorId) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ColorId, &clinitDone_ColorId)
	static_ColorId_FARM_64 = j2gort.Ext_com_google_common_hash_Hashing_FarmHashFingerprint64()
	static_ColorId_ASCII_AT = lambdaSingleton_color_id_1
}

func (this_ *ColorId) initFields_ColorId() {
	_ = this_
	this_.valid_ = true
}

// Java constructor: ColorId(long rightAligned)
func NewColorId(rightAligned int64) ColorId {
	clinit_ColorId()
	this_ := ColorId{}
	InitColorId(&this_, rightAligned)
	return this_
}

func InitColorId(this_ *ColorId, rightAligned int64) {
	_ = rightAligned
	this_.initFields_ColorId()
	this_.rightAligned = rightAligned
}

// Java static method: ColorId.fromUnsigned(long x)
func ColorId_FromUnsigned_1_int64(x int64) ColorId {
	clinit_ColorId()
	_ = x
	return NewColorId(x)
}

// Java static method: ColorId.fromUnsigned(int x)
func ColorId_FromUnsigned_1_int32(x int32) ColorId {
	clinit_ColorId()
	_ = x
	return NewColorId((int64(x) & int64(4294967295)))
}

// Java static method: ColorId.fromUnsigned(byte x)
func ColorId_FromUnsigned_1_int8(x int8) ColorId {
	clinit_ColorId()
	_ = x
	return NewColorId((int64(x) & int64(255)))
}

// Java static method: ColorId.fromBytes(com.google.protobuf.ByteString bytes)
func ColorId_FromBytes_1_Ext_com_google_protobuf_ByteString(bytes j2gort.Ext_com_google_protobuf_ByteString) ColorId {
	clinit_ColorId()
	_ = bytes
	return ColorId_FromBytesAt(j2gort.BoxByteString(bytes), bytes.Size(), lambdaSingleton_color_id_2)
}

// Java static method: ColorId.fromBytes(byte[] bytes)
func ColorId_FromBytes_1_slice_int8(bytes []int8) ColorId {
	clinit_ColorId()
	_ = bytes
	return ColorId_FromBytesAt(j2gort.SliceToAny(bytes), int32(len(bytes)), lambdaSingleton_color_id_3)
}

// Java static method: ColorId.fromAscii(java.lang.String str)
func ColorId_FromAscii(str string) ColorId {
	clinit_ColorId()
	_ = str
	return ColorId_FromBytesAt(j2gort.BoxString(str), j2gort.JStrLen(str), StaticGet_ColorId_ASCII_AT())
}

// Java static method: ColorId.fromBytesAt(T source, int length, com.google.javascript.jscomp.colors.ColorId.ByteAt<T> byteAt)
func ColorId_FromBytesAt(source any, length int32, byteAt ColorId_ByteAt) ColorId {
	clinit_ColorId()
	_ = source
	_ = length
	_ = byteAt
	if !(length <= int32(8)) {
		j2gort.CheckStateMsg(false, j2gort.ValueOfInt(length))
	}
	var rightAligned int64 = int64(int32(0))
	_ = rightAligned
	{
		var i int32 = int32(0)
		_ = i
		for ; i < length; i = i + 1 {
			rightAligned = javalang.Shl64(rightAligned, int(int32(8)))
			rightAligned = (rightAligned | ((int64(byteAt.Get(source, i))) & int64(int32(255))))
		}
	}
	return NewColorId(rightAligned)
}

// Java static method: ColorId.union(java.util.Set<com.google.javascript.jscomp.colors.ColorId> ids)
func ColorId_Union(ids j2gort.Ext_java_util_Set) ColorId {
	clinit_ColorId()
	_ = ids
	if ids.Size() <= int32(1) {
		return j2gort.Cast[ColorId](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](ids)))
	}
	var index int32 = int32(0)
	_ = index
	var sorted []int64 = make([]int64, ids.Size())
	_ = sorted
	{
		cur_3 := j2gort.NewTypedIterCursor[ColorId](j2gort.CheckNotNull(ids))
		for cur_3.HasNext() {
			id := cur_3.Next()
			_ = id
			{
				a_4 := sorted
				i_5 := func() int32 { old_6 := index; nw_7 := old_6 + 1; index = nw_7; return old_6 }()
				a_4[i_5] = id.rightAligned
			}
		}
	}
	j2gort.Ext_java_util_Arrays_Sort_1_slice_int64(sorted)
	var hasher j2gort.Ext_com_google_common_hash_Hasher = StaticGet_ColorId_FARM_64().NewHasher_0()
	_ = hasher
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(sorted)); i = i + 1 {
			hasher.PutLong(sorted[i])
		}
	}
	return NewColorId(hasher.Hash().AsLong())
}

// Java static method: ColorId.unionFromColors(com.google.javascript.jscomp.colors.Color[] colors)
func ColorId_UnionFromColors(colors []*Color) ColorId {
	clinit_ColorId()
	_ = colors
	var len_ int32 = int32(len(colors))
	_ = len_
	if len_ == int32(1) {
		return colors[int32(0)].GetId()
	}
	if len_ == int32(2) {
		return ColorId_Union2(colors[int32(0)].GetId(), colors[int32(1)].GetId())
	}
	if len_ == int32(3) {
		return ColorId_Union3(colors[int32(0)].GetId(), colors[int32(1)].GetId(), colors[int32(2)].GetId())
	}
	if len_ == int32(4) {
		return ColorId_Union4(colors[int32(0)].GetId(), colors[int32(1)].GetId(), colors[int32(2)].GetId(), colors[int32(3)].GetId())
	}
	var sorted []int64 = make([]int64, len_)
	_ = sorted
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			{
				a_8 := sorted
				i_9 := i
				a_8[i_9] = colors[i].GetId().rightAligned
			}
		}
	}
	j2gort.Ext_java_util_Arrays_Sort_1_slice_int64(sorted)
	var hasher j2gort.Ext_com_google_common_hash_Hasher = StaticGet_ColorId_FARM_64().NewHasher_0()
	_ = hasher
	var prev int64 = sorted[int32(0)]
	_ = prev
	hasher.PutLong(prev)
	{
		var i_2 int32 = int32(1)
		_ = i_2
		for ; i_2 < len_; i_2 = i_2 + 1 {
			var cur int64 = sorted[i_2]
			_ = cur
			if cur != prev {
				hasher.PutLong(cur)
				prev = cur
			}
		}
	}
	return NewColorId(hasher.Hash().AsLong())
}

// Java static method: ColorId.union2(com.google.javascript.jscomp.colors.ColorId a, com.google.javascript.jscomp.colors.ColorId b)
func ColorId_Union2(a ColorId, b ColorId) ColorId {
	clinit_ColorId()
	_ = a
	_ = b
	var x int64 = a.rightAligned
	_ = x
	var y int64 = b.rightAligned
	_ = y
	if x == y {
		return a
	}
	if x > y {
		var t int64 = x
		_ = t
		x = y
		y = t
	}
	var hasher j2gort.Ext_com_google_common_hash_Hasher = StaticGet_ColorId_FARM_64().NewHasher_0()
	_ = hasher
	hasher.PutLong(x)
	hasher.PutLong(y)
	return NewColorId(hasher.Hash().AsLong())
}

// Java static method: ColorId.union3(com.google.javascript.jscomp.colors.ColorId a, com.google.javascript.jscomp.colors.ColorId b, com.google.javascript.jscomp.colors.ColorId c)
func ColorId_Union3(a ColorId, b ColorId, c ColorId) ColorId {
	clinit_ColorId()
	_ = a
	_ = b
	_ = c
	var x int64 = a.rightAligned
	_ = x
	var y int64 = b.rightAligned
	_ = y
	var z int64 = c.rightAligned
	_ = z
	if x > y {
		var t int64 = x
		_ = t
		x = y
		y = t
	}
	if y > z {
		var t_2 int64 = y
		_ = t_2
		y = z
		z = t_2
	}
	if x > y {
		var t_3 int64 = x
		_ = t_3
		x = y
		y = t_3
	}
	var hasher j2gort.Ext_com_google_common_hash_Hasher = StaticGet_ColorId_FARM_64().NewHasher_0()
	_ = hasher
	hasher.PutLong(x)
	if y != x {
		hasher.PutLong(y)
	}
	if z != y {
		hasher.PutLong(z)
	}
	return NewColorId(hasher.Hash().AsLong())
}

// Java static method: ColorId.union4(com.google.javascript.jscomp.colors.ColorId a, com.google.javascript.jscomp.colors.ColorId b, com.google.javascript.jscomp.colors.ColorId c, com.google.javascript.jscomp.colors.ColorId d)
func ColorId_Union4(a ColorId, b ColorId, c ColorId, d ColorId) ColorId {
	clinit_ColorId()
	_ = a
	_ = b
	_ = c
	_ = d
	var v0 int64 = a.rightAligned
	_ = v0
	var v1 int64 = b.rightAligned
	_ = v1
	var v2 int64 = c.rightAligned
	_ = v2
	var v3 int64 = d.rightAligned
	_ = v3
	if v0 > v1 {
		var t int64 = v0
		_ = t
		v0 = v1
		v1 = t
	}
	if v2 > v3 {
		var t_2 int64 = v2
		_ = t_2
		v2 = v3
		v3 = t_2
	}
	if v0 > v2 {
		var t_3 int64 = v0
		_ = t_3
		v0 = v2
		v2 = t_3
	}
	if v1 > v3 {
		var t_4 int64 = v1
		_ = t_4
		v1 = v3
		v3 = t_4
	}
	if v1 > v2 {
		var t_5 int64 = v1
		_ = t_5
		v1 = v2
		v2 = t_5
	}
	var hasher j2gort.Ext_com_google_common_hash_Hasher = StaticGet_ColorId_FARM_64().NewHasher_0()
	_ = hasher
	hasher.PutLong(v0)
	if v1 != v0 {
		hasher.PutLong(v1)
	}
	if v2 != v1 {
		hasher.PutLong(v2)
	}
	if v3 != v2 {
		hasher.PutLong(v3)
	}
	return NewColorId(hasher.Hash().AsLong())
}

// Java method: ColorId.equals(java.lang.Object x)
func (this_ ColorId) Equals(arg_10 any) bool {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	x := arg_10
	_ = x
	var that ColorId
	_ = that
	var ok_11 bool
	that, ok_11 = j2gort.TryCast[ColorId](x)
	if !ok_11 {
		return false
	}
	return (this_.rightAligned == that.rightAligned)
}

// Java method: ColorId.hashCode()
func (this_ ColorId) HashCode() int32 {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return int32(this_.rightAligned)
}

// Java method: ColorId.asByteString()
func (this_ ColorId) AsByteString() j2gort.Ext_com_google_protobuf_ByteString {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	var copy_ int64 = this_.rightAligned
	_ = copy_
	var out []int8 = make([]int8, int32(8))
	_ = out
	{
		var i int32 = int32(7)
		_ = i
		for ; i >= int32(0); i = i - 1 {
			{
				a_12 := out
				i_13 := i
				a_12[i_13] = int8(copy_)
			}
			copy_ = javalang.Ushr64(copy_, int(int32(8)))
		}
	}
	return j2gort.Ext_com_google_protobuf_ByteString_CopyFrom_1_slice_int8(out)
}

// Java method: ColorId.toString()
func (this_ ColorId) String() string {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return javalang.LongToHexString(this_.rightAligned)
}

// Java method: ColorId.toLogsGson()
func (this_ ColorId) ToLogsGson() any {
	if !this_.valid_ {
		j2gort.ThrowNPE()
	}
	return j2gort.BoxString(this_.String())
}

// Java interface: com.google.javascript.jscomp.colors.ColorId.ByteAt
type ColorId_ByteAt interface {
	Implements_colors_ColorId_ByteAt()
	IsJavaNil_() bool
	Get(any, int32) int8
}

func ToIface_ColorId_ByteAt[T interface {
	ColorId_ByteAt
	comparable
}](v T) ColorId_ByteAt {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ColorId_ByteAt[T ColorId_ByteAt](v T) ColorId_ByteAt {
	var vi ColorId_ByteAt = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.colors.ColorId.ByteAt", "com.google.javascript.jscomp.colors.ColorId$ByteAt", func(v_ any) bool {
		return j2gort.InstanceOf[ColorId_ByteAt](v_)
	})
}

type lambda_ColorId_ByteAt struct {
	fn func(any, int32) int8
}

func (l *lambda_ColorId_ByteAt) IsJavaNil_() bool { return l == nil }
func NewLambda_ColorId_ByteAt(fn func(any, int32) int8) ColorId_ByteAt {
	if fn == nil {
		return nil
	}
	return &lambda_ColorId_ByteAt{fn: fn}
}

func (*lambda_ColorId_ByteAt) Implements_colors_ColorId_ByteAt() {}
func (l *lambda_ColorId_ByteAt) Get(p0_ any, p1_ int32) int8 {
	return l.fn(p0_, p1_)
}

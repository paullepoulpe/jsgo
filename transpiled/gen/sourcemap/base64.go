/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/debugging/sourcemap/Base64.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.debugging.sourcemap.Base64
type Base64 struct {
	idHash_ atomic.Int32
}

func (this_ *Base64) BindSelf_(s any) {
	_ = s
}
func (this_ *Base64) IsJavaNil_() bool { return this_ == nil }
func (*Base64) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.Base64", "Base64"
}
func (this_ *Base64) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Base64) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.Base64", "com.google.debugging.sourcemap.Base64", func(v_ any) bool {
		return j2gort.InstanceOf[*Base64](v_)
	})
}

func StaticGet_Base64_BASE64_MAP() string {
	return "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
}

var static_Base64_BASE64_DECODE_MAP []int32

func StaticGet_Base64_BASE64_DECODE_MAP() []int32 {
	clinit_Base64()
	return static_Base64_BASE64_DECODE_MAP
}
func StaticSet_Base64_BASE64_DECODE_MAP(v []int32) {
	clinit_Base64()
	static_Base64_BASE64_DECODE_MAP = v
}

var clinitDone_Base64 uint32
var clinitState_Base64 j2gort.ClinitState

func clinit_Base64() {
	if atomic.LoadUint32(&clinitDone_Base64) != 0 {
		return
	}
	clinitSlow_Base64()
}

func clinitSlow_Base64() {
	if !j2gort.EnterClinit(&clinitState_Base64, &clinitDone_Base64) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Base64, &clinitDone_Base64)
	static_Base64_BASE64_DECODE_MAP = make([]int32, int32(256))
	j2gort.ArraysFill(StaticGet_Base64_BASE64_DECODE_MAP(), int32(-1))
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"); i = i + 1 {
			{
				a_1 := StaticGet_Base64_BASE64_DECODE_MAP()
				i_2 := int32(j2gort.JStrCharAt("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/", i))
				a_1[i_2] = i
			}
		}
	}
}

func (this_ *Base64) initFields_Base64() {
	_ = this_
}

// Java constructor: Base64()
func NewBase64() *Base64 {
	clinit_Base64()
	this_ := &Base64{}
	InitBase64(this_)
	return this_
}

func InitBase64(this_ *Base64) {
	this_.initFields_Base64()
}

// Java static method: Base64.toBase64(int value)
func Base64_ToBase64(value int32) uint16 {
	clinit_Base64()
	_ = value
	if !((value <= int32(63)) && (value >= int32(0))) {
		panic(javalang.NewAssertionError("value out of range:" + j2gort.StrValInt32(value)))
	}
	return j2gort.JStrCharAt("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/", value)
}

// Java static method: Base64.fromBase64(char c)
func Base64_FromBase64(c uint16) int32 {
	clinit_Base64()
	_ = c
	var result int32 = StaticGet_Base64_BASE64_DECODE_MAP()[int32(c)]
	_ = result
	if !(result != int32(-1)) {
		panic(javalang.NewAssertionError("invalid char"))
	}
	return result
}

// Java static method: Base64.base64EncodeInt(int value)
func Base64_Base64EncodeInt(value int32) string {
	clinit_Base64()
	_ = value
	var c []uint16 = make([]uint16, int32(6))
	_ = c
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(5); i = i + 1 {
			{
				a_3 := c
				i_4 := i
				a_3[i_4] = Base64_ToBase64(((javalang.Shr32(value, int((int32(26) - (i * int32(6)))))) & int32(63)))
			}
		}
	}
	{
		a_5 := c
		i_6 := int32(5)
		a_5[i_6] = Base64_ToBase64(((javalang.Shl32(value, int(int32(4)))) & int32(63)))
	}
	return javalang.StringFromUTF16(c)
}

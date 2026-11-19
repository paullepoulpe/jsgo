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
// Port of com/google/debugging/sourcemap/Base64VLQ.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.debugging.sourcemap.Base64VLQ
type Base64VLQ struct {
	idHash_ atomic.Int32
}

func (this_ *Base64VLQ) BindSelf_(s any) {
	_ = s
}
func (this_ *Base64VLQ) IsJavaNil_() bool { return this_ == nil }
func (*Base64VLQ) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.Base64VLQ", "Base64VLQ"
}
func (this_ *Base64VLQ) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Base64VLQ) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.Base64VLQ", "com.google.debugging.sourcemap.Base64VLQ", func(v_ any) bool {
		return j2gort.InstanceOf[*Base64VLQ](v_)
	})
}

func StaticGet_Base64VLQ_VLQ_BASE_SHIFT() int32       { return int32(5) }
func StaticGet_Base64VLQ_VLQ_BASE() int32             { return int32(32) }
func StaticGet_Base64VLQ_VLQ_BASE_MASK() int32        { return int32(31) }
func StaticGet_Base64VLQ_VLQ_CONTINUATION_BIT() int32 { return int32(32) }
func (this_ *Base64VLQ) initFields_Base64VLQ() {
	_ = this_
}

// Java constructor: Base64VLQ()
func NewBase64VLQ() *Base64VLQ {
	this_ := &Base64VLQ{}
	InitBase64VLQ(this_)
	return this_
}

func InitBase64VLQ(this_ *Base64VLQ) {
	this_.initFields_Base64VLQ()
}

// Java static method: Base64VLQ.toVLQSigned(int value)
func Base64VLQ_ToVLQSigned(value int32) int32 {
	_ = value
	if value < int32(0) {
		return ((javalang.Shl32((-(value)), int(int32(1)))) + int32(1))
	} else {
		return ((javalang.Shl32(value, int(int32(1)))) + int32(0))
	}
}

// Java static method: Base64VLQ.fromVLQSigned(int value)
func Base64VLQ_FromVLQSigned(value int32) int32 {
	_ = value
	var negate bool = ((value & int32(1)) == int32(1))
	_ = negate
	value = javalang.Ushr32(value, int(int32(1)))
	if !(negate) {
		return value
	}
	return (-(value) | int32(-2147483648))
}

// Java static method: Base64VLQ.encode(java.lang.Appendable out, int value)
func Base64VLQ_Encode(out j2gort.Ext_java_lang_Appendable, value int32) {
	_ = out
	_ = value
	value = Base64VLQ_ToVLQSigned(value)
	for doFirst_1 := true; doFirst_1 || (value > int32(0)); doFirst_1 = false {
		var digit int32 = (value & int32(31))
		_ = digit
		value = javalang.Ushr32(value, int(int32(5)))
		if value > int32(0) {
			digit = (digit | int32(32))
		}
		out.Append_1_uint16(Base64_ToBase64(digit))
	}
}

// Java static method: Base64VLQ.decode(com.google.debugging.sourcemap.Base64VLQ.CharIterator in)
func Base64VLQ_Decode(in Base64VLQ_CharIterator) int32 {
	_ = in
	var result int32 = int32(0)
	_ = result
	var continuation bool = false
	_ = continuation
	var shift int32 = int32(0)
	_ = shift
	for doFirst_2 := true; doFirst_2 || (continuation); doFirst_2 = false {
		var c uint16 = in.Next()
		_ = c
		var digit int32 = Base64_FromBase64(c)
		_ = digit
		continuation = ((digit & int32(32)) != int32(0))
		digit = (digit & int32(31))
		result = (result + (javalang.Shl32(digit, int(shift))))
		shift = (shift + int32(5))
	}
	return Base64VLQ_FromVLQSigned(result)
}

// Java interface: com.google.debugging.sourcemap.Base64VLQ.CharIterator
type Base64VLQ_CharIterator interface {
	Implements_sourcemap_Base64VLQ_CharIterator()
	IsJavaNil_() bool
	HasNext() bool
	Next() uint16
}

func ToIface_Base64VLQ_CharIterator[T interface {
	Base64VLQ_CharIterator
	comparable
}](v T) Base64VLQ_CharIterator {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Base64VLQ_CharIterator[T Base64VLQ_CharIterator](v T) Base64VLQ_CharIterator {
	var vi Base64VLQ_CharIterator = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.Base64VLQ.CharIterator", "com.google.debugging.sourcemap.Base64VLQ$CharIterator", func(v_ any) bool {
		return j2gort.InstanceOf[Base64VLQ_CharIterator](v_)
	})
}

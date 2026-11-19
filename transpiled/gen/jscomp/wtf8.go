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
// Port of com/google/javascript/jscomp/serialization/Wtf8.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.serialization.Wtf8
type Wtf8 struct {
	idHash_ atomic.Int32
}

func (this_ *Wtf8) BindSelf_(s any) {
	_ = s
}
func (this_ *Wtf8) IsJavaNil_() bool { return this_ == nil }
func (*Wtf8) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.serialization.Wtf8", "Wtf8" }
func (this_ *Wtf8) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Wtf8) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.Wtf8", "com.google.javascript.jscomp.serialization.Wtf8", func(v_ any) bool {
		return j2gort.InstanceOf[*Wtf8](v_)
	})
}

func StaticGet_Wtf8_CONTINUATION_MASK() int8 { return int8(63) }
func (this_ *Wtf8) initFields_Wtf8() {
	_ = this_
}

// Java constructor: Wtf8()
func NewWtf8() *Wtf8 {
	this_ := &Wtf8{}
	InitWtf8(this_)
	return this_
}

func InitWtf8(this_ *Wtf8) {
	this_.initFields_Wtf8()
}

// Java static method: Wtf8.isAscii(java.lang.String s, int length)
func Wtf8_IsAscii(s string, length int32) bool {
	_ = s
	_ = length
	{
		var i int32 = int32(0)
		_ = i
		for ; (i < length); i = i + 1 {
			if ((int32(j2gort.JStrCharAt(s, i)) >= int32(128))) {
				return false
			}
		}
	}
	return true
}

// Java static method: Wtf8.encodeToWtf8(java.lang.String s)
func Wtf8_EncodeToWtf8(s string) j2gort.Ext_com_google_protobuf_ByteString {
	_ = s
	var length int32 = j2gort.JStrLen(s)
	_ = length
	if ((length == int32(0))) {
		return j2gort.Ext_com_google_protobuf_ByteString_F_EMPTY
	}
	if (Wtf8_IsAscii(s, length)) {
		return j2gort.Ext_com_google_protobuf_ByteString_CopyFromUtf8(s)
	}
	var output *j2gort.Ext_com_google_protobuf_ByteString_Output = j2gort.Ext_com_google_protobuf_ByteString_NewOutput_1(length)
	_ = output
	{
		var i int32 = int32(0)
		_ = i
		for ; (i < length); i = i + 1 {
			var codepoint int32 = j2gort.JStrCodePointAt(s, i)
			_ = codepoint
			if ((codepoint < int32(128))) {
				output.Write_1_int32(codepoint)
			} else if ((codepoint < int32(2048))) {
				output.Write_1_int32((int32(192) | ((int32(31) & (javalang.Ushr32(codepoint, int(int32(6))))))))
				output.Write_1_int32((int32(128) | ((int32(int8(63)) & codepoint))))
			} else if ((codepoint < int32(65536))) {
				output.Write_1_int32((int32(224) | ((int32(15) & (javalang.Ushr32(codepoint, int(int32(12))))))))
				output.Write_1_int32((int32(128) | ((int32(int8(63)) & (javalang.Ushr32(codepoint, int(int32(6))))))))
				output.Write_1_int32((int32(128) | ((int32(int8(63)) & codepoint))))
			} else {
				i = i + 1
				output.Write_1_int32((int32(240) | ((int32(7) & (javalang.Ushr32(codepoint, int(int32(18))))))))
				output.Write_1_int32((int32(128) | ((int32(int8(63)) & (javalang.Ushr32(codepoint, int(int32(12))))))))
				output.Write_1_int32((int32(128) | ((int32(int8(63)) & (javalang.Ushr32(codepoint, int(int32(6))))))))
				output.Write_1_int32((int32(128) | ((int32(int8(63)) & codepoint))))
			}
		}
	}
	return output.ToByteString()
}

// Java static method: Wtf8.decoder(int maxLength)
func Wtf8_S_Decoder(maxLength int32) *Wtf8_Decoder {
	_ = maxLength
	return NewWtf8_Decoder(maxLength)
}

// Java class: com.google.javascript.jscomp.serialization.Wtf8.Decoder
type Wtf8_Decoder struct {
	codepointBuffer []int32 `gson:"codepointBuffer"`
	idHash_ atomic.Int32
}

func (this_ *Wtf8_Decoder) BindSelf_(s any) {
	_ = s
}
func (this_ *Wtf8_Decoder) IsJavaNil_() bool { return this_ == nil }
func (*Wtf8_Decoder) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.serialization.Wtf8$Decoder", "Decoder" }
func (this_ *Wtf8_Decoder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Wtf8_Decoder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.Wtf8.Decoder", "com.google.javascript.jscomp.serialization.Wtf8$Decoder", func(v_ any) bool {
		return j2gort.InstanceOf[*Wtf8_Decoder](v_)
	})
}

func (this_ *Wtf8_Decoder) initFields_Wtf8_Decoder() {
	_ = this_
}

// Java constructor: Decoder(int maxLength)
func NewWtf8_Decoder(maxLength int32) *Wtf8_Decoder {
	this_ := &Wtf8_Decoder{}
	InitWtf8_Decoder(this_, maxLength)
	return this_
}

func InitWtf8_Decoder(this_ *Wtf8_Decoder, maxLength int32) {
	_ = maxLength
	this_.initFields_Wtf8_Decoder()
	this_.codepointBuffer = make([]int32, maxLength)
}

// Java static method: Decoder.isAsciiBytes(com.google.protobuf.ByteString encoded, int length)
func Wtf8_Decoder_IsAsciiBytes(encoded j2gort.Ext_com_google_protobuf_ByteString, length int32) bool {
	_ = encoded
	_ = length
	{
		var i int32 = int32(0)
		_ = i
		for ; (i < length); i = i + 1 {
			if ((int32(encoded.ByteAt(i)) < int32(0))) {
				return false
			}
		}
	}
	return true
}

// Java method: Decoder.decode(com.google.protobuf.ByteString encoded)
func (this_ *Wtf8_Decoder) Decode(encoded j2gort.Ext_com_google_protobuf_ByteString) string {
	_ = encoded
	var encodedByteCount int32 = encoded.Size()
	_ = encodedByteCount
	if ((encodedByteCount == int32(0))) {
		return ""
	}
	if (Wtf8_Decoder_IsAsciiBytes(encoded, encodedByteCount)) {
		return encoded.ToStringUtf8()
	}
	var codepointCount int32 = int32(0)
	_ = codepointCount
	{
		var i int32 = int32(0)
		_ = i
		for ; (i < encodedByteCount); i = i + 1 {
			var b int8 = encoded.ByteAt(i)
			_ = b
			var codepoint int32 = int32(0)
			_ = codepoint
			if ((((int32(b) & int32(128))) == int32(0))) {
				codepoint = int32(b)
			} else if ((((int32(b) & int32(224))) == int32(192))) {
				var firstByte int32 = (int32(31) & int32(b))
				_ = firstByte
				var secondByte int32 = (int32(encoded.ByteAt(func() int32 { old_1 := i; nw_2 := old_1 + 1; i = nw_2; return nw_2 }())) & int32(int8(63)))
				_ = secondByte
				codepoint = ((javalang.Shl32(firstByte, int(int32(6)))) | secondByte)
			} else if ((((int32(b) & int32(240))) == int32(224))) {
				var firstByte_2 int32 = (int32(15) & int32(b))
				_ = firstByte_2
				var secondByte_2 int32 = (int32(encoded.ByteAt(func() int32 { old_3 := i; nw_4 := old_3 + 1; i = nw_4; return nw_4 }())) & int32(int8(63)))
				_ = secondByte_2
				var thirdByte int32 = (int32(encoded.ByteAt(func() int32 { old_5 := i; nw_6 := old_5 + 1; i = nw_6; return nw_6 }())) & int32(int8(63)))
				_ = thirdByte
				codepoint = (((javalang.Shl32(firstByte_2, int(int32(12)))) | (javalang.Shl32(secondByte_2, int(int32(6))))) | thirdByte)
			} else if ((((int32(b) & int32(248))) == int32(240))) {
				var firstByte_3 int32 = (int32(7) & int32(b))
				_ = firstByte_3
				var secondByte_3 int32 = (int32(encoded.ByteAt(func() int32 { old_7 := i; nw_8 := old_7 + 1; i = nw_8; return nw_8 }())) & int32(int8(63)))
				_ = secondByte_3
				var thirdByte_2 int32 = (int32(encoded.ByteAt(func() int32 { old_9 := i; nw_10 := old_9 + 1; i = nw_10; return nw_10 }())) & int32(int8(63)))
				_ = thirdByte_2
				var fourthByte int32 = (int32(encoded.ByteAt(func() int32 { old_11 := i; nw_12 := old_11 + 1; i = nw_12; return nw_12 }())) & int32(int8(63)))
				_ = fourthByte
				codepoint = ((((javalang.Shl32(firstByte_3, int(int32(18)))) | (javalang.Shl32(secondByte_3, int(int32(12))))) | (javalang.Shl32(thirdByte_2, int(int32(6))))) | fourthByte)
			} else {
				panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
			}
			{ a_13 := this_.codepointBuffer; i_14 := func() int32 { old_15 := codepointCount; nw_16 := old_15 + 1; codepointCount = nw_16; return old_15 }(); a_13[i_14] = codepoint }
		}
	}
	return javalang.StringFromCodePoints(this_.codepointBuffer, int32(0), codepointCount)
}


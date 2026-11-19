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
// Port of com/google/debugging/sourcemap/Util.java (generated mechanically by j2go).

package sourcemap

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.debugging.sourcemap.Util
type Util struct {
	idHash_ atomic.Int32
}

func (this_ *Util) BindSelf_(s any) {
	_ = s
}
func (this_ *Util) IsJavaNil_() bool           { return this_ == nil }
func (*Util) JavaClassInfo_() (string, string) { return "com.google.debugging.sourcemap.Util", "Util" }
func (this_ *Util) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Util) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.Util", "com.google.debugging.sourcemap.Util", func(v_ any) bool {
		return j2gort.InstanceOf[*Util](v_)
	})
}

var static_Util_HEX_CHARS []uint16

func StaticGet_Util_HEX_CHARS() []uint16 {
	clinit_Util()
	return static_Util_HEX_CHARS
}
func StaticSet_Util_HEX_CHARS(v []uint16) {
	clinit_Util()
	static_Util_HEX_CHARS = v
}

var clinitDone_Util uint32
var clinitState_Util j2gort.ClinitState

func clinit_Util() {
	if atomic.LoadUint32(&clinitDone_Util) != 0 {
		return
	}
	clinitSlow_Util()
}

func clinitSlow_Util() {
	if !j2gort.EnterClinit(&clinitState_Util, &clinitDone_Util) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Util, &clinitDone_Util)
	static_Util_HEX_CHARS = []uint16{uint16(48), uint16(49), uint16(50), uint16(51), uint16(52), uint16(53), uint16(54), uint16(55), uint16(56), uint16(57), uint16(97), uint16(98), uint16(99), uint16(100), uint16(101), uint16(102)}
}

func (this_ *Util) initFields_Util() {
	_ = this_
}

// Java constructor: Util()
func NewUtil() *Util {
	clinit_Util()
	this_ := &Util{}
	InitUtil(this_)
	return this_
}

func InitUtil(this_ *Util) {
	this_.initFields_Util()
}

// Java static method: Util.escapeString(java.lang.String s)
func Util_EscapeString(s string) string {
	clinit_Util()
	_ = s
	var quote uint16 = uint16(34)
	_ = quote
	var doublequoteEscape string = "\\\""
	_ = doublequoteEscape
	var singlequoteEscape string = "'"
	_ = singlequoteEscape
	var backslashEscape string = "\\\\"
	_ = backslashEscape
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_int32((j2gort.JStrLen(s) + int32(2)))
	_ = sb
	sb.Append_1_uint16(quote)
	var region *Util__Local_1_UnescapedRegion = NewUtil__Local_1_UnescapedRegion(sb, s)
	_ = region
	var length int32 = j2gort.JStrLen(s)
	_ = length
	{
		var i int32 = int32(0)
		_ = i
		for ; i < length; i = i + 1 {
			var c uint16 = j2gort.JStrCharAt(s, i)
			_ = c
			switch c {
			case uint16(10):
				region.AppendForEscapedChar("\\n")
			case uint16(13):
				region.AppendForEscapedChar("\\r")
			case uint16(9):
				region.AppendForEscapedChar("\\t")
			case uint16(92):
				region.AppendForEscapedChar(backslashEscape)
			case uint16(34):
				region.AppendForEscapedChar(doublequoteEscape)
			case uint16(39):
				region.AppendForEscapedChar(singlequoteEscape)
			case uint16(62):
				if (i >= int32(2)) && (((int32(j2gort.JStrCharAt(s, (i-int32(1)))) == int32(uint16(45))) && (int32(j2gort.JStrCharAt(s, (i-int32(2)))) == int32(uint16(45)))) || ((int32(j2gort.JStrCharAt(s, (i-int32(1)))) == int32(uint16(93))) && (int32(j2gort.JStrCharAt(s, (i-int32(2)))) == int32(uint16(93))))) {
					region.AppendForEscapedChar("\\u003e")
				} else {
					region.IncrementForNormalChar()
				}
			case uint16(60):
				var l_END_SCRIPT string = "/script"
				_ = l_END_SCRIPT
				var l_START_COMMENT string = "!--"
				_ = l_START_COMMENT
				if j2gort.JStrRegionMatches(s, true, (i + int32(1)), l_END_SCRIPT, int32(0), j2gort.JStrLen(l_END_SCRIPT)) {
					region.AppendForEscapedChar("\\u003c")
				} else if j2gort.JStrRegionMatches(s, false, (i + int32(1)), l_START_COMMENT, int32(0), j2gort.JStrLen(l_START_COMMENT)) {
					region.AppendForEscapedChar("\\u003c")
				} else {
					region.IncrementForNormalChar()
				}
			default:
				if (int32(c) > int32(31)) && (int32(c) <= int32(127)) {
					region.IncrementForNormalChar()
				} else {
					region.AppendUnescaped()
					region.IncrementForEscapedChar()
					Util_AppendHexJavaScriptRepresentation_2_Ext_java_lang_StringBuilder_uint16(sb, c)
				}
			}
		}
	}
	region.AppendUnescaped()
	sb.Append_1_uint16(quote)
	return sb.String()
}

// Java static method: Util.appendHexJavaScriptRepresentation(java.lang.StringBuilder sb, char c)
func Util_AppendHexJavaScriptRepresentation_2_Ext_java_lang_StringBuilder_uint16(sb *j2gort.Ext_java_lang_StringBuilder, c uint16) {
	clinit_Util()
	_ = sb
	_ = c
	func() {
		defer func() {
			if rec_1 := recover(); rec_1 != nil {
				jex_2 := javalang.RecoverException(rec_1)
				_ = jex_2
				if j2gort.InstanceOf[javaio.IIOException](jex_2) {
					ex := j2gort.AsIface[javaio.IIOException](jex_2)
					_ = ex
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", ex)))
				} else {
					panic(rec_1)
				}
			}
		}()
		Util_AppendHexJavaScriptRepresentation_2_int32_Ext_java_lang_Appendable(int32(c), j2gort.AsIface[j2gort.Ext_java_lang_Appendable](sb))
	}()
}

// Java static method: Util.appendHexJavaScriptRepresentation(int codePoint, java.lang.Appendable out)
func Util_AppendHexJavaScriptRepresentation_2_int32_Ext_java_lang_Appendable(codePoint int32, out j2gort.Ext_java_lang_Appendable) {
	clinit_Util()
	_ = codePoint
	_ = out
	if j2gort.BoxChar_IsSupplementaryCodePoint(codePoint) {
		var surrogates []uint16 = j2gort.BoxChar_ToChars_1(codePoint)
		_ = surrogates
		Util_AppendHexJavaScriptRepresentation_2_int32_Ext_java_lang_Appendable(int32(surrogates[int32(0)]), out)
		Util_AppendHexJavaScriptRepresentation_2_int32_Ext_java_lang_Appendable(int32(surrogates[int32(1)]), out)
		return
	}
	j2gort.AsIface[j2gort.Ext_java_lang_Appendable](j2gort.AsIface[j2gort.Ext_java_lang_Appendable](j2gort.AsIface[j2gort.Ext_java_lang_Appendable](j2gort.AsIface[j2gort.Ext_java_lang_Appendable](out.Append_1_any("\\u")).Append_1_uint16(StaticGet_Util_HEX_CHARS()[((javalang.Ushr32(codePoint, int(int32(12)))) & int32(15))])).Append_1_uint16(StaticGet_Util_HEX_CHARS()[((javalang.Ushr32(codePoint, int(int32(8)))) & int32(15))])).Append_1_uint16(StaticGet_Util_HEX_CHARS()[((javalang.Ushr32(codePoint, int(int32(4)))) & int32(15))])).Append_1_uint16(StaticGet_Util_HEX_CHARS()[(codePoint & int32(15))])
}

// Java class: UnescapedRegion
type Util__Local_1_UnescapedRegion struct {
	cap_0_sb             *j2gort.Ext_java_lang_StringBuilder
	cap_1_s              string
	unescapedRegionStart int32 `gson:"unescapedRegionStart"`
	unescapedRegionEnd   int32 `gson:"unescapedRegionEnd"`
	idHash_              atomic.Int32
}

func (this_ *Util__Local_1_UnescapedRegion) BindSelf_(s any) {
	_ = s
}
func (this_ *Util__Local_1_UnescapedRegion) IsJavaNil_() bool { return this_ == nil }
func (*Util__Local_1_UnescapedRegion) JavaClassInfo_() (string, string) {
	return "com.google.debugging.sourcemap.Util$1UnescapedRegion", "UnescapedRegion"
}
func (this_ *Util__Local_1_UnescapedRegion) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Util__Local_1_UnescapedRegion) HasIdentityEquals_() {}

func (this_ *Util__Local_1_UnescapedRegion) initFields_Util__Local_1_UnescapedRegion() {
	_ = this_
	this_.unescapedRegionStart = int32(0)
	this_.unescapedRegionEnd = int32(0)
}

// Java constructor: UnescapedRegion()
func NewUtil__Local_1_UnescapedRegion(cap_p_sb *j2gort.Ext_java_lang_StringBuilder, cap_p_s string) *Util__Local_1_UnescapedRegion {
	this_ := &Util__Local_1_UnescapedRegion{}
	this_.cap_0_sb = cap_p_sb
	this_.cap_1_s = cap_p_s
	InitUtil__Local_1_UnescapedRegion(this_)
	return this_
}

func InitUtil__Local_1_UnescapedRegion(this_ *Util__Local_1_UnescapedRegion) {
	this_.initFields_Util__Local_1_UnescapedRegion()
}

// Java method: UnescapedRegion.appendUnescaped()
func (this_ *Util__Local_1_UnescapedRegion) AppendUnescaped() {
	if this_.unescapedRegionStart != this_.unescapedRegionEnd {
		this_.cap_0_sb.Append_1_string(j2gort.JStrSub(this_.cap_1_s, this_.unescapedRegionStart, this_.unescapedRegionEnd))
	}
	this_.unescapedRegionStart = this_.unescapedRegionEnd
}

// Java method: UnescapedRegion.incrementForNormalChar()
func (this_ *Util__Local_1_UnescapedRegion) IncrementForNormalChar() {
	this_.unescapedRegionEnd = this_.unescapedRegionEnd + 1
}

// Java method: UnescapedRegion.incrementForEscapedChar()
func (this_ *Util__Local_1_UnescapedRegion) IncrementForEscapedChar() {
	if this_.unescapedRegionStart != this_.unescapedRegionEnd {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
	}
	this_.unescapedRegionStart = this_.unescapedRegionStart + 1
	this_.unescapedRegionEnd = this_.unescapedRegionEnd + 1
}

// Java method: UnescapedRegion.appendForEscapedChar(java.lang.String escaped)
func (this_ *Util__Local_1_UnescapedRegion) AppendForEscapedChar(escaped string) {
	_ = escaped
	this_.AppendUnescaped()
	this_.IncrementForEscapedChar()
	this_.cap_0_sb.Append_1_string(escaped)
}

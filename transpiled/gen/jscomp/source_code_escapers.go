/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/SourceCodeEscapers.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.SourceCodeEscapers
type SourceCodeEscapers struct {
	idHash_ atomic.Int32
}

func (this_ *SourceCodeEscapers) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceCodeEscapers) IsJavaNil_() bool { return this_ == nil }
func (*SourceCodeEscapers) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.SourceCodeEscapers", "SourceCodeEscapers"
}
func (this_ *SourceCodeEscapers) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SourceCodeEscapers) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.SourceCodeEscapers", "com.google.javascript.jscomp.deps.SourceCodeEscapers", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceCodeEscapers](v_)
	})
}

func StaticGet_SourceCodeEscapers_PRINTABLE_ASCII_MIN() uint16 { return uint16(32) }
func StaticGet_SourceCodeEscapers_PRINTABLE_ASCII_MAX() uint16 { return uint16(126) }

var static_SourceCodeEscapers_HEX_DIGITS []uint16

func StaticGet_SourceCodeEscapers_HEX_DIGITS() []uint16 {
	clinit_SourceCodeEscapers()
	return static_SourceCodeEscapers_HEX_DIGITS
}
func StaticSet_SourceCodeEscapers_HEX_DIGITS(v []uint16) {
	clinit_SourceCodeEscapers()
	static_SourceCodeEscapers_HEX_DIGITS = v
}

var static_SourceCodeEscapers_JAVASCRIPT_ESCAPER *SourceCodeEscapers_JavaScriptEscaper

func StaticGet_SourceCodeEscapers_JAVASCRIPT_ESCAPER() *SourceCodeEscapers_JavaScriptEscaper {
	clinit_SourceCodeEscapers()
	return static_SourceCodeEscapers_JAVASCRIPT_ESCAPER
}
func StaticSet_SourceCodeEscapers_JAVASCRIPT_ESCAPER(v *SourceCodeEscapers_JavaScriptEscaper) {
	clinit_SourceCodeEscapers()
	static_SourceCodeEscapers_JAVASCRIPT_ESCAPER = v
}

var clinitDone_SourceCodeEscapers uint32
var clinitState_SourceCodeEscapers j2gort.ClinitState

func clinit_SourceCodeEscapers() {
	if atomic.LoadUint32(&clinitDone_SourceCodeEscapers) != 0 {
		return
	}
	clinitSlow_SourceCodeEscapers()
}

func clinitSlow_SourceCodeEscapers() {
	if !j2gort.EnterClinit(&clinitState_SourceCodeEscapers, &clinitDone_SourceCodeEscapers) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SourceCodeEscapers, &clinitDone_SourceCodeEscapers)
	static_SourceCodeEscapers_HEX_DIGITS = javalang.StringToUTF16("0123456789abcdef")
	static_SourceCodeEscapers_JAVASCRIPT_ESCAPER = NewSourceCodeEscapers_JavaScriptEscaper()
}

func (this_ *SourceCodeEscapers) initFields_SourceCodeEscapers() {
	_ = this_
}

// Java constructor: SourceCodeEscapers()
func NewSourceCodeEscapers() *SourceCodeEscapers {
	clinit_SourceCodeEscapers()
	this_ := &SourceCodeEscapers{}
	InitSourceCodeEscapers(this_)
	return this_
}

func InitSourceCodeEscapers(this_ *SourceCodeEscapers) {
	this_.initFields_SourceCodeEscapers()
}

// Java static method: SourceCodeEscapers.javascriptEscaper()
func SourceCodeEscapers_JavascriptEscaper() j2gort.Ext_com_google_common_escape_Escaper {
	clinit_SourceCodeEscapers()
	return j2gort.ToIface_Ext_com_google_common_escape_Escaper(StaticGet_SourceCodeEscapers_JAVASCRIPT_ESCAPER())
}

// Java static method: SourceCodeEscapers.appendWithJavascriptEscaper(java.lang.CharSequence c, java.lang.Appendable to)
func SourceCodeEscapers_AppendWithJavascriptEscaper(c any, to j2gort.Ext_java_lang_Appendable) {
	clinit_SourceCodeEscapers()
	_ = c
	_ = to
	StaticGet_SourceCodeEscapers_JAVASCRIPT_ESCAPER().AppendTo(c, to)
}

// Java class: com.google.javascript.jscomp.deps.SourceCodeEscapers.JavaScriptEscaper
type SourceCodeEscapers_JavaScriptEscaper struct {
	j2gort.Ext_com_google_common_escape_EscaperBase
	idHash_ atomic.Int32
}

func (this_ *SourceCodeEscapers_JavaScriptEscaper) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_common_escape_EscaperBase.BindExtSelf_(s)
}
func (this_ *SourceCodeEscapers_JavaScriptEscaper) IsJavaNil_() bool { return this_ == nil }
func (*SourceCodeEscapers_JavaScriptEscaper) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.SourceCodeEscapers$JavaScriptEscaper", "JavaScriptEscaper"
}
func (this_ *SourceCodeEscapers_JavaScriptEscaper) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.SourceCodeEscapers.JavaScriptEscaper", "com.google.javascript.jscomp.deps.SourceCodeEscapers$JavaScriptEscaper", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceCodeEscapers_JavaScriptEscaper](v_)
	})
}
func (*SourceCodeEscapers_JavaScriptEscaper) Implements_j2gort_Ext_com_google_common_escape_Escaper() {
}

var static_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS [][]uint16

func StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS() [][]uint16 {
	clinit_SourceCodeEscapers_JavaScriptEscaper()
	return static_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS
}
func StaticSet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS(v [][]uint16) {
	clinit_SourceCodeEscapers_JavaScriptEscaper()
	static_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS = v
}

var clinitDone_SourceCodeEscapers_JavaScriptEscaper uint32
var clinitState_SourceCodeEscapers_JavaScriptEscaper j2gort.ClinitState

func clinit_SourceCodeEscapers_JavaScriptEscaper() {
	if atomic.LoadUint32(&clinitDone_SourceCodeEscapers_JavaScriptEscaper) != 0 {
		return
	}
	clinitSlow_SourceCodeEscapers_JavaScriptEscaper()
}

func clinitSlow_SourceCodeEscapers_JavaScriptEscaper() {
	if !j2gort.EnterClinit(&clinitState_SourceCodeEscapers_JavaScriptEscaper, &clinitDone_SourceCodeEscapers_JavaScriptEscaper) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SourceCodeEscapers_JavaScriptEscaper, &clinitDone_SourceCodeEscapers_JavaScriptEscaper)
	StaticSet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS(make([][]uint16, int32(256)))
	{
		var i uint16 = uint16(int32(0))
		_ = i
		for ; int32(i) < int32(uint16(32)); i = i + 1 {
			var c uint16 = i
			_ = c
			var r []uint16 = make([]uint16, int32(4))
			_ = r
			{
				a_1 := r
				i_2 := int32(3)
				a_1[i_2] = StaticGet_SourceCodeEscapers_HEX_DIGITS()[(int32(c) & int32(15))]
			}
			c = uint16((javalang.Ushr32(int32(c), int(int32(4)))))
			{
				a_3 := r
				i_4 := int32(2)
				a_3[i_4] = StaticGet_SourceCodeEscapers_HEX_DIGITS()[(int32(c) & int32(15))]
			}
			{
				a_5 := r
				i_6 := int32(1)
				a_5[i_6] = uint16(120)
			}
			{
				a_7 := r
				i_8 := int32(0)
				a_7[i_8] = uint16(92)
			}
			{
				a_9 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
				i_10 := int32(i)
				a_9[i_10] = r
			}
		}
	}
	{
		var i_2 uint16 = uint16((int32(uint16(126)) + int32(1)))
		_ = i_2
		for ; int32(i_2) < int32(256); i_2 = i_2 + 1 {
			var c_2 uint16 = i_2
			_ = c_2
			var r_2 []uint16 = make([]uint16, int32(4))
			_ = r_2
			{
				a_11 := r_2
				i_12 := int32(3)
				a_11[i_12] = StaticGet_SourceCodeEscapers_HEX_DIGITS()[(int32(c_2) & int32(15))]
			}
			c_2 = uint16((javalang.Ushr32(int32(c_2), int(int32(4)))))
			{
				a_13 := r_2
				i_14 := int32(2)
				a_13[i_14] = StaticGet_SourceCodeEscapers_HEX_DIGITS()[(int32(c_2) & int32(15))]
			}
			{
				a_15 := r_2
				i_16 := int32(1)
				a_15[i_16] = uint16(120)
			}
			{
				a_17 := r_2
				i_18 := int32(0)
				a_17[i_18] = uint16(92)
			}
			{
				a_19 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
				i_20 := int32(i_2)
				a_19[i_20] = r_2
			}
		}
	}
	{
		a_21 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_22 := int32(uint16(39))
		a_21[i_22] = javalang.StringToUTF16("\\x27")
	}
	{
		a_23 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_24 := int32(uint16(34))
		a_23[i_24] = javalang.StringToUTF16("\\x22")
	}
	{
		a_25 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_26 := int32(uint16(60))
		a_25[i_26] = javalang.StringToUTF16("\\x3c")
	}
	{
		a_27 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_28 := int32(uint16(61))
		a_27[i_28] = javalang.StringToUTF16("\\x3d")
	}
	{
		a_29 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_30 := int32(uint16(62))
		a_29[i_30] = javalang.StringToUTF16("\\x3e")
	}
	{
		a_31 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_32 := int32(uint16(38))
		a_31[i_32] = javalang.StringToUTF16("\\x26")
	}
	{
		a_33 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_34 := int32(uint16(8))
		a_33[i_34] = javalang.StringToUTF16("\\b")
	}
	{
		a_35 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_36 := int32(uint16(9))
		a_35[i_36] = javalang.StringToUTF16("\\t")
	}
	{
		a_37 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_38 := int32(uint16(10))
		a_37[i_38] = javalang.StringToUTF16("\\n")
	}
	{
		a_39 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_40 := int32(uint16(12))
		a_39[i_40] = javalang.StringToUTF16("\\f")
	}
	{
		a_41 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_42 := int32(uint16(13))
		a_41[i_42] = javalang.StringToUTF16("\\r")
	}
	{
		a_43 := StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()
		i_44 := int32(uint16(92))
		a_43[i_44] = javalang.StringToUTF16("\\\\")
	}
}

func (this_ *SourceCodeEscapers_JavaScriptEscaper) initFields_SourceCodeEscapers_JavaScriptEscaper() {
	_ = this_
}

// Java constructor: JavaScriptEscaper()
func NewSourceCodeEscapers_JavaScriptEscaper() *SourceCodeEscapers_JavaScriptEscaper {
	clinit_SourceCodeEscapers_JavaScriptEscaper()
	this_ := &SourceCodeEscapers_JavaScriptEscaper{}
	this_.BindSelf_(this_)
	InitSourceCodeEscapers_JavaScriptEscaper(this_)
	return this_
}

func InitSourceCodeEscapers_JavaScriptEscaper(this_ *SourceCodeEscapers_JavaScriptEscaper) {
	j2gort.InitExt_com_google_common_escape_Escaper(&this_.Ext_com_google_common_escape_EscaperBase)
	this_.initFields_SourceCodeEscapers_JavaScriptEscaper()
}

// Java method: JavaScriptEscaper.appendTo(java.lang.CharSequence cs, java.lang.Appendable to)
func (this_ *SourceCodeEscapers_JavaScriptEscaper) AppendTo(cs any, to j2gort.Ext_java_lang_Appendable) {
	_ = cs
	_ = to
	var last int32 = int32(0)
	_ = last
	var length int32 = j2gort.CharSeqLength(cs)
	_ = length
	{
		var i int32 = int32(0)
		_ = i
		for ; i < length; i = i + 1 {
			var c uint16 = j2gort.CharSeqCharAt(cs, i)
			_ = c
			var replacement []uint16 = nil
			_ = replacement
			if int32(c) < int32(256) {
				replacement = StaticGet_SourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS()[int32(c)]
				if replacement == nil {
					continue
				}
			} else {
				replacement = SourceCodeEscapers_JavaScriptEscaper_AsUnicodeHexEscape(c)
			}
			if last < i {
				to.Append_3_any_int32_int32(cs, last, i)
			}
			for _, el_45 := range j2gort.CheckNotNullSlice(replacement) {
				r := el_45
				_ = r
				to.Append_1_uint16(r)
			}
			last = (i + int32(1))
		}
	}
	if last < length {
		to.Append_3_any_int32_int32(cs, last, length)
	}
}

// Java static method: JavaScriptEscaper.asUnicodeHexEscape(char c)
func SourceCodeEscapers_JavaScriptEscaper_AsUnicodeHexEscape(c uint16) []uint16 {
	clinit_SourceCodeEscapers_JavaScriptEscaper()
	_ = c
	var r []uint16 = make([]uint16, int32(6))
	_ = r
	{
		a_46 := r
		i_47 := int32(0)
		a_46[i_47] = uint16(92)
	}
	{
		a_48 := r
		i_49 := int32(1)
		a_48[i_49] = uint16(117)
	}
	{
		a_50 := r
		i_51 := int32(5)
		a_50[i_51] = StaticGet_SourceCodeEscapers_HEX_DIGITS()[(int32(c) & int32(15))]
	}
	c = uint16((javalang.Ushr32(int32(c), int(int32(4)))))
	{
		a_52 := r
		i_53 := int32(4)
		a_52[i_53] = StaticGet_SourceCodeEscapers_HEX_DIGITS()[(int32(c) & int32(15))]
	}
	c = uint16((javalang.Ushr32(int32(c), int(int32(4)))))
	{
		a_54 := r
		i_55 := int32(3)
		a_54[i_55] = StaticGet_SourceCodeEscapers_HEX_DIGITS()[(int32(c) & int32(15))]
	}
	c = uint16((javalang.Ushr32(int32(c), int(int32(4)))))
	{
		a_56 := r
		i_57 := int32(2)
		a_56[i_57] = StaticGet_SourceCodeEscapers_HEX_DIGITS()[(int32(c) & int32(15))]
	}
	return r
}

// Java method: JavaScriptEscaper.escape(java.lang.String string)
func (this_ *SourceCodeEscapers_JavaScriptEscaper) Escape_1_string(string_ string) string {
	_ = string_
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	func() {
		defer func() {
			if rec_58 := recover(); rec_58 != nil {
				jex_59 := javalang.RecoverException(rec_58)
				_ = jex_59
				if j2gort.InstanceOf[javaio.IIOException](jex_59) {
					e := j2gort.AsIface[javaio.IIOException](jex_59)
					_ = e
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "This should never throw - StringBuilder.append doesn't actually throw IOException", e)))
				} else {
					panic(rec_58)
				}
			}
		}()
		this_.AppendTo(j2gort.BoxString(string_), sb)
	}()
	return sb.String()
}

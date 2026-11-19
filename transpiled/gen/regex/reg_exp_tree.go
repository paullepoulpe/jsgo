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
// Port of com/google/javascript/jscomp/regex/RegExpTree.java (generated mechanically by j2go).

package regex

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.regex.RegExpTree
type RegExpTree interface {
	Implements_regex_RegExpTree()
	IsJavaNil_() bool
	GetRegExpTreeBase() *RegExpTreeBase
	Simplify(string) RegExpTree
	IsCaseSensitive() bool
	ContainsAnchor() bool
	HasCapturingGroup() bool
	NumCapturingGroups() int32
	Children() any
	AppendSourceCode(*j2gort.Ext_java_lang_StringBuilder)
	AppendDebugInfo(*j2gort.Ext_java_lang_StringBuilder)
	String() string
	Equals(any) bool
	HashCode() int32
}

func ToIface_RegExpTree[T interface {
	RegExpTree
	comparable
}](v T) RegExpTree {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_RegExpTree[T RegExpTree](v T) RegExpTree {
	var vi RegExpTree = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree
type RegExpTreeBase struct {
	self_ RegExpTree
}

func (this_ *RegExpTreeBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(RegExpTree)
}
func (this_ *RegExpTreeBase) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTreeBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree", "RegExpTree"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree", "com.google.javascript.jscomp.regex.RegExpTree", func(v_ any) bool {
		return j2gort.InstanceOf[RegExpTree](v_)
	})
}
func (b *RegExpTreeBase) GetRegExpTreeBase() *RegExpTreeBase { return b }
func (*RegExpTreeBase) Implements_regex_RegExpTree()         {}

var static_RegExpTree_NEVER_MATCHES RegExpTree

func StaticGet_RegExpTree_NEVER_MATCHES() RegExpTree {
	clinit_RegExpTree()
	return static_RegExpTree_NEVER_MATCHES
}
func StaticSet_RegExpTree_NEVER_MATCHES(v RegExpTree) {
	clinit_RegExpTree()
	static_RegExpTree_NEVER_MATCHES = v
}

var static_RegExpTree_DIGITS *CharRanges

func StaticGet_RegExpTree_DIGITS() *CharRanges {
	clinit_RegExpTree()
	return static_RegExpTree_DIGITS
}
func StaticSet_RegExpTree_DIGITS(v *CharRanges) {
	clinit_RegExpTree()
	static_RegExpTree_DIGITS = v
}

var static_RegExpTree_UCASE_LETTERS *CharRanges

func StaticGet_RegExpTree_UCASE_LETTERS() *CharRanges {
	clinit_RegExpTree()
	return static_RegExpTree_UCASE_LETTERS
}
func StaticSet_RegExpTree_UCASE_LETTERS(v *CharRanges) {
	clinit_RegExpTree()
	static_RegExpTree_UCASE_LETTERS = v
}

var static_RegExpTree_LCASE_LETTERS *CharRanges

func StaticGet_RegExpTree_LCASE_LETTERS() *CharRanges {
	clinit_RegExpTree()
	return static_RegExpTree_LCASE_LETTERS
}
func StaticSet_RegExpTree_LCASE_LETTERS(v *CharRanges) {
	clinit_RegExpTree()
	static_RegExpTree_LCASE_LETTERS = v
}

var static_RegExpTree_LETTERS *CharRanges

func StaticGet_RegExpTree_LETTERS() *CharRanges {
	clinit_RegExpTree()
	return static_RegExpTree_LETTERS
}
func StaticSet_RegExpTree_LETTERS(v *CharRanges) {
	clinit_RegExpTree()
	static_RegExpTree_LETTERS = v
}

var static_RegExpTree_WORD_CHARS *CharRanges

func StaticGet_RegExpTree_WORD_CHARS() *CharRanges {
	clinit_RegExpTree()
	return static_RegExpTree_WORD_CHARS
}
func StaticSet_RegExpTree_WORD_CHARS(v *CharRanges) {
	clinit_RegExpTree()
	static_RegExpTree_WORD_CHARS = v
}

var static_RegExpTree_INVERSE_WORD_CHARS *CharRanges

func StaticGet_RegExpTree_INVERSE_WORD_CHARS() *CharRanges {
	clinit_RegExpTree()
	return static_RegExpTree_INVERSE_WORD_CHARS
}
func StaticSet_RegExpTree_INVERSE_WORD_CHARS(v *CharRanges) {
	clinit_RegExpTree()
	static_RegExpTree_INVERSE_WORD_CHARS = v
}

var static_RegExpTree_SPACE_CHARS *CharRanges

func StaticGet_RegExpTree_SPACE_CHARS() *CharRanges {
	clinit_RegExpTree()
	return static_RegExpTree_SPACE_CHARS
}
func StaticSet_RegExpTree_SPACE_CHARS(v *CharRanges) {
	clinit_RegExpTree()
	static_RegExpTree_SPACE_CHARS = v
}

var static_RegExpTree_IE_SPACE_CHARS *CharRanges

func StaticGet_RegExpTree_IE_SPACE_CHARS() *CharRanges {
	clinit_RegExpTree()
	return static_RegExpTree_IE_SPACE_CHARS
}
func StaticSet_RegExpTree_IE_SPACE_CHARS(v *CharRanges) {
	clinit_RegExpTree()
	static_RegExpTree_IE_SPACE_CHARS = v
}

var static_RegExpTree_IE_SPEC_ERRORS *CharRanges

func StaticGet_RegExpTree_IE_SPEC_ERRORS() *CharRanges {
	clinit_RegExpTree()
	return static_RegExpTree_IE_SPEC_ERRORS
}
func StaticSet_RegExpTree_IE_SPEC_ERRORS(v *CharRanges) {
	clinit_RegExpTree()
	static_RegExpTree_IE_SPEC_ERRORS = v
}

var static_RegExpTree_NAMED_CHAR_GROUPS j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_RegExpTree_NAMED_CHAR_GROUPS() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_RegExpTree()
	return static_RegExpTree_NAMED_CHAR_GROUPS
}
func StaticSet_RegExpTree_NAMED_CHAR_GROUPS(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_RegExpTree()
	static_RegExpTree_NAMED_CHAR_GROUPS = v
}

var static_RegExpTree_DOT_CHARSET *RegExpTree_Charset

func StaticGet_RegExpTree_DOT_CHARSET() *RegExpTree_Charset {
	clinit_RegExpTree()
	return static_RegExpTree_DOT_CHARSET
}
func StaticSet_RegExpTree_DOT_CHARSET(v *RegExpTree_Charset) {
	clinit_RegExpTree()
	static_RegExpTree_DOT_CHARSET = v
}

var clinitDone_RegExpTree uint32
var clinitState_RegExpTree j2gort.ClinitState

func clinit_RegExpTree() {
	if atomic.LoadUint32(&clinitDone_RegExpTree) != 0 {
		return
	}
	clinitSlow_RegExpTree()
}

func clinitSlow_RegExpTree() {
	if !j2gort.EnterClinit(&clinitState_RegExpTree, &clinitDone_RegExpTree) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RegExpTree, &clinitDone_RegExpTree)
	static_RegExpTree_NEVER_MATCHES = NewRegExpTree_LookaheadAssertion(ToIface_RegExpTree(StaticGet_RegExpTree_Empty_INSTANCE()), false)
	static_RegExpTree_DIGITS = CharRanges_Inclusive(int32(uint16(48)), int32(uint16(57)))
	static_RegExpTree_UCASE_LETTERS = CharRanges_Inclusive(int32(uint16(65)), int32(uint16(90)))
	static_RegExpTree_LCASE_LETTERS = CharRanges_Inclusive(int32(uint16(97)), int32(uint16(122)))
	static_RegExpTree_LETTERS = StaticGet_RegExpTree_UCASE_LETTERS().Union(StaticGet_RegExpTree_LCASE_LETTERS())
	static_RegExpTree_WORD_CHARS = StaticGet_RegExpTree_DIGITS().Union(StaticGet_RegExpTree_LETTERS()).Union(CharRanges_WithMembers([]int32{int32(uint16(95))}))
	static_RegExpTree_INVERSE_WORD_CHARS = StaticGet_CharRanges_ALL_CODE_UNITS().Difference(StaticGet_RegExpTree_WORD_CHARS())
	static_RegExpTree_SPACE_CHARS = CharRanges_WithMembers([]int32{int32(uint16(9)), int32(uint16(10)), int32(uint16(11)), int32(uint16(12)), int32(uint16(13)), int32(uint16(32)), int32(uint16(160)), int32(uint16(5760)), int32(uint16(6158)), int32(uint16(8192)), int32(uint16(8193)), int32(uint16(8194)), int32(uint16(8195)), int32(uint16(8196)), int32(uint16(8197)), int32(uint16(8198)), int32(uint16(8199)), int32(uint16(8200)), int32(uint16(8201)), int32(uint16(8202)), int32(uint16(8232)), int32(uint16(8233)), int32(uint16(8239)), int32(uint16(8287)), int32(uint16(12288)), int32(uint16(65279))})
	static_RegExpTree_IE_SPACE_CHARS = CharRanges_WithMembers([]int32{int32(uint16(9)), int32(uint16(10)), int32(uint16(11)), int32(uint16(12)), int32(uint16(13)), int32(uint16(32))})
	static_RegExpTree_IE_SPEC_ERRORS = StaticGet_RegExpTree_SPACE_CHARS().Difference(StaticGet_RegExpTree_IE_SPACE_CHARS())
	static_RegExpTree_NAMED_CHAR_GROUPS = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder().Put_2(j2gort.ValueOfChar(uint16(100)), j2gort.PtrToAny(StaticGet_RegExpTree_DIGITS()))).Put_2(j2gort.ValueOfChar(uint16(68)), j2gort.PtrToAny(StaticGet_CharRanges_ALL_CODE_UNITS().Difference(StaticGet_RegExpTree_DIGITS())))).Put_2(j2gort.ValueOfChar(uint16(115)), j2gort.PtrToAny(StaticGet_RegExpTree_SPACE_CHARS()))).Put_2(j2gort.ValueOfChar(uint16(83)), j2gort.PtrToAny(StaticGet_CharRanges_ALL_CODE_UNITS().Difference(StaticGet_RegExpTree_SPACE_CHARS())))).Put_2(j2gort.ValueOfChar(uint16(119)), j2gort.PtrToAny(StaticGet_RegExpTree_WORD_CHARS()))).Put_2(j2gort.ValueOfChar(uint16(87)), j2gort.PtrToAny(StaticGet_RegExpTree_INVERSE_WORD_CHARS()))).BuildOrThrow())
	static_RegExpTree_DOT_CHARSET = NewRegExpTree_Charset(StaticGet_CharRanges_ALL_CODE_UNITS().Difference(CharRanges_WithMembers([]int32{int32(uint16(10)), int32(uint16(13)), int32(uint16(8232)), int32(uint16(8233))})), StaticGet_CharRanges_EMPTY())
}

func (this_ *RegExpTreeBase) initFields_RegExpTree() {
	_ = this_
}

func InitRegExpTree(this_ *RegExpTreeBase) {
	this_.initFields_RegExpTree()
}

// Java method: RegExpTree.simplify(java.lang.String flags)
func (this_ *RegExpTreeBase) Simplify(flags string) RegExpTree {
	_ = flags
	panic("abstract or native method invoked")
}

// Java method: RegExpTree.isCaseSensitive()
func (this_ *RegExpTreeBase) IsCaseSensitive() bool {
	panic("abstract or native method invoked")
}

// Java method: RegExpTree.containsAnchor()
func (this_ *RegExpTreeBase) ContainsAnchor() bool {
	panic("abstract or native method invoked")
}

// Java method: RegExpTree.hasCapturingGroup()
func (this_ *RegExpTreeBase) HasCapturingGroup() bool {
	return (this_.self_.NumCapturingGroups() != int32(0))
}

// Java method: RegExpTree.numCapturingGroups()
func (this_ *RegExpTreeBase) NumCapturingGroups() int32 {
	panic("abstract or native method invoked")
}

// Java method: RegExpTree.children()
func (this_ *RegExpTreeBase) Children() any {
	panic("abstract or native method invoked")
}

// Java method: RegExpTree.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTreeBase) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	panic("abstract or native method invoked")
}

// Java method: RegExpTree.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTreeBase) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	panic("abstract or native method invoked")
}

// Java method: RegExpTree.toString()
func (this_ *RegExpTreeBase) String() string {
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	sb.Append_1_uint16(uint16(47))
	this_.self_.AppendSourceCode(sb)
	if sb.Length() == int32(1) {
		sb.Append_1_string("(?:)")
	}
	sb.Append_1_uint16(uint16(47))
	return sb.String()
}

// Java method: RegExpTree.equals(java.lang.Object o)
func (this_ *RegExpTreeBase) Equals(arg_1 any) bool {
	o := arg_1
	_ = o
	panic("abstract or native method invoked")
}

// Java method: RegExpTree.hashCode()
func (this_ *RegExpTreeBase) HashCode() int32 {
	panic("abstract or native method invoked")
}

// Java static method: RegExpTree.parseRegExp(java.lang.String pattern, java.lang.String flags)
func RegExpTree_ParseRegExp(pattern string, flags string) RegExpTree {
	clinit_RegExpTree()
	_ = pattern
	_ = flags
	return NewRegExpTree__Local_124_Parser(pattern, flags).ParseTopLevel()
}

// Java static method: RegExpTree.isIdentifierStart(char ch)
func RegExpTree_IsIdentifierStart(ch uint16) bool {
	clinit_RegExpTree()
	_ = ch
	if int32(ch) <= int32(127) {
		return ((((int32(ch) >= int32(uint16(65))) && (int32(ch) <= int32(uint16(90)))) || ((int32(ch) >= int32(uint16(97))) && (int32(ch) <= int32(uint16(122))))) || ((int32(ch) == int32(uint16(95))) || (int32(ch) == int32(uint16(36)))))
	}
	return j2gort.BoxChar_IsLetter_1_uint16(ch)
}

// Java static method: RegExpTree.isIdentifierPart(char ch)
func RegExpTree_IsIdentifierPart(ch uint16) bool {
	clinit_RegExpTree()
	_ = ch
	if int32(ch) <= int32(127) {
		return (((((int32(ch) >= int32(uint16(65))) && (int32(ch) <= int32(uint16(90)))) || ((int32(ch) >= int32(uint16(97))) && (int32(ch) <= int32(uint16(122))))) || ((int32(ch) >= int32(uint16(48))) && (int32(ch) <= int32(uint16(57))))) || ((int32(ch) == int32(uint16(95))) || (int32(ch) == int32(uint16(36)))))
	}
	return (RegExpTree_IsIdentifierStart(ch) || j2gort.BoxChar_IsDigit_1_uint16(ch))
}

// Java static method: RegExpTree.matchesWholeInput(com.google.javascript.jscomp.regex.RegExpTree t, java.lang.String flags)
func RegExpTree_MatchesWholeInput(t RegExpTree, flags string) bool {
	clinit_RegExpTree()
	_ = t
	_ = flags
	if j2gort.JStrIndexOfCp(flags, int32(uint16(109))) >= int32(0) {
		return false
	}
	var c *RegExpTree_Concatenation
	_ = c
	var ok_2 bool
	c, ok_2 = j2gort.TryCast[*RegExpTree_Concatenation](t)
	if !ok_2 {
		return false
	}
	if len(c.elements) == 0 {
		return false
	}
	var first RegExpTree = c.elements[int32(0)]
	_ = first
	var last RegExpTree = j2gort.SliceGetLast(c.elements)
	_ = last
	var firstAnchor *RegExpTree_Anchor
	_ = firstAnchor
	var lastAnchor *RegExpTree_Anchor
	_ = lastAnchor
	var ok_3 bool
	firstAnchor, ok_3 = j2gort.TryCast[*RegExpTree_Anchor](first)
	if ok_3 {
		lastAnchor, ok_3 = j2gort.TryCast[*RegExpTree_Anchor](last)
	}
	if !ok_3 {
		return false
	}
	return ((int32(firstAnchor.type_) == int32(uint16(94))) && (int32(lastAnchor.type_) == int32(uint16(36))))
}

// Java static method: RegExpTree.escapeCharOnto(char ch, java.lang.StringBuilder sb)
func RegExpTree_EscapeCharOnto(ch uint16, sb *j2gort.Ext_java_lang_StringBuilder) {
	clinit_RegExpTree()
	_ = ch
	_ = sb
	switch ch {
	case uint16(0):
		sb.Append_1_string("\\0")
	case uint16(12):
		sb.Append_1_string("\\f")
	case uint16(9):
		sb.Append_1_string("\\t")
	case uint16(10):
		sb.Append_1_string("\\n")
	case uint16(13):
		sb.Append_1_string("\\r")
	case uint16(92):
		sb.Append_1_string("\\\\")
	default:
		if (int32(ch) < int32(32)) || (int32(ch) >= int32(127)) {
			if int32(ch) >= int32(256) {
				sb.Append_1_string("\\u")
				sb.Append_1_uint16(j2gort.JStrCharAt("0123456789abcdef", ((javalang.Shr32(int32(ch), int(int32(12)))) & int32(15))))
				sb.Append_1_uint16(j2gort.JStrCharAt("0123456789abcdef", ((javalang.Shr32(int32(ch), int(int32(8)))) & int32(15))))
				sb.Append_1_uint16(j2gort.JStrCharAt("0123456789abcdef", ((javalang.Shr32(int32(ch), int(int32(4)))) & int32(15))))
				sb.Append_1_uint16(j2gort.JStrCharAt("0123456789abcdef", (int32((ch)) & int32(15))))
			} else {
				sb.Append_1_string("\\x")
				sb.Append_1_uint16(j2gort.JStrCharAt("0123456789abcdef", ((javalang.Shr32(int32(ch), int(int32(4)))) & int32(15))))
				sb.Append_1_uint16(j2gort.JStrCharAt("0123456789abcdef", (int32((ch)) & int32(15))))
			}
		} else {
			sb.Append_1_uint16(ch)
		}
	}
}

// Java enum: com.google.javascript.jscomp.regex.RegExpTree.ParentheticalType
type RegExpTree_ParentheticalType uint8

const (
	enum_RegExpTree_ParentheticalType_CAPTURING           RegExpTree_ParentheticalType = 1
	enum_RegExpTree_ParentheticalType_NONCAPTURING        RegExpTree_ParentheticalType = 2
	enum_RegExpTree_ParentheticalType_POSITIVE_LOOKAHEAD  RegExpTree_ParentheticalType = 3
	enum_RegExpTree_ParentheticalType_NEGATIVE_LOOKAHEAD  RegExpTree_ParentheticalType = 4
	enum_RegExpTree_ParentheticalType_POSITIVE_LOOKBEHIND RegExpTree_ParentheticalType = 5
	enum_RegExpTree_ParentheticalType_NEGATIVE_LOOKBEHIND RegExpTree_ParentheticalType = 6
	enum_RegExpTree_ParentheticalType_NAMED_GROUPS        RegExpTree_ParentheticalType = 7
)

var enumBaseTable_RegExpTree_ParentheticalType = [...]j2gort.EnumBase{
	{Name_: "CAPTURING", Ordinal_: 0},
	{Name_: "NONCAPTURING", Ordinal_: 1},
	{Name_: "POSITIVE_LOOKAHEAD", Ordinal_: 2},
	{Name_: "NEGATIVE_LOOKAHEAD", Ordinal_: 3},
	{Name_: "POSITIVE_LOOKBEHIND", Ordinal_: 4},
	{Name_: "NEGATIVE_LOOKBEHIND", Ordinal_: 5},
	{Name_: "NAMED_GROUPS", Ordinal_: 6},
}
var enumValues_RegExpTree_ParentheticalType = []RegExpTree_ParentheticalType{1, 2, 3, 4, 5, 6, 7}

func StaticGet_RegExpTree_ParentheticalType_CAPTURING() RegExpTree_ParentheticalType {
	return enum_RegExpTree_ParentheticalType_CAPTURING
}
func StaticGet_RegExpTree_ParentheticalType_NONCAPTURING() RegExpTree_ParentheticalType {
	return enum_RegExpTree_ParentheticalType_NONCAPTURING
}
func StaticGet_RegExpTree_ParentheticalType_POSITIVE_LOOKAHEAD() RegExpTree_ParentheticalType {
	return enum_RegExpTree_ParentheticalType_POSITIVE_LOOKAHEAD
}
func StaticGet_RegExpTree_ParentheticalType_NEGATIVE_LOOKAHEAD() RegExpTree_ParentheticalType {
	return enum_RegExpTree_ParentheticalType_NEGATIVE_LOOKAHEAD
}
func StaticGet_RegExpTree_ParentheticalType_POSITIVE_LOOKBEHIND() RegExpTree_ParentheticalType {
	return enum_RegExpTree_ParentheticalType_POSITIVE_LOOKBEHIND
}
func StaticGet_RegExpTree_ParentheticalType_NEGATIVE_LOOKBEHIND() RegExpTree_ParentheticalType {
	return enum_RegExpTree_ParentheticalType_NEGATIVE_LOOKBEHIND
}
func StaticGet_RegExpTree_ParentheticalType_NAMED_GROUPS() RegExpTree_ParentheticalType {
	return enum_RegExpTree_ParentheticalType_NAMED_GROUPS
}

func RegExpTree_ParentheticalType_Values() []RegExpTree_ParentheticalType {
	return j2gort.CloneSlice(enumValues_RegExpTree_ParentheticalType)
}

func RegExpTree_ParentheticalType_ValueOf(name string) RegExpTree_ParentheticalType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_RegExpTree_ParentheticalType {
		if enumBaseTable_RegExpTree_ParentheticalType[i_].Name_ == name {
			return RegExpTree_ParentheticalType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.regex.RegExpTree.ParentheticalType." + name))
}

func (this_ RegExpTree_ParentheticalType) IsJavaNil_() bool                              { return this_ == 0 }
func (RegExpTree_ParentheticalType) HasIdentityEquals_()                                 {}
func (RegExpTree_ParentheticalType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (RegExpTree_ParentheticalType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ RegExpTree_ParentheticalType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (RegExpTree_ParentheticalType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$ParentheticalType", "ParentheticalType"
}
func (this_ RegExpTree_ParentheticalType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_RegExpTree_ParentheticalType[this_-1]
}
func (this_ RegExpTree_ParentheticalType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ RegExpTree_ParentheticalType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ RegExpTree_ParentheticalType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_RegExpTree_ParentheticalType[this_-1].Name_
}
func (this_ RegExpTree_ParentheticalType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_RegExpTree_ParentheticalType[this_-1].Name_
}
func (this_ RegExpTree_ParentheticalType) Equals(other_ any) bool {
	o_, ok_ := other_.(RegExpTree_ParentheticalType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ RegExpTree_ParentheticalType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_RegExpTree_ParentheticalType[this_-1].HashCode()
}
func (this_ RegExpTree_ParentheticalType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_RegExpTree_ParentheticalType[this_-1].CompareTo(other_)
}
func (this_ RegExpTree_ParentheticalType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_RegExpTree_ParentheticalType[this_-1].Name_)
}
func init() {
	vals_ := []any{RegExpTree_ParentheticalType(1), RegExpTree_ParentheticalType(2), RegExpTree_ParentheticalType(3), RegExpTree_ParentheticalType(4), RegExpTree_ParentheticalType(5), RegExpTree_ParentheticalType(6), RegExpTree_ParentheticalType(7)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.regex.RegExpTree.ParentheticalType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.regex.RegExpTree$ParentheticalType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.ParentheticalType", "com.google.javascript.jscomp.regex.RegExpTree$ParentheticalType", func(v_ any) bool {
		return j2gort.InstanceOf[RegExpTree_ParentheticalType](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.regex.RegExpTree.RegExpTreeAtom
type RegExpTree_RegExpTreeAtom interface {
	Implements_regex_RegExpTree_RegExpTreeAtom()
	IsJavaNil_() bool
	GetRegExpTree_RegExpTreeAtomBase() *RegExpTree_RegExpTreeAtomBase
	RegExpTree
}

func ToIface_RegExpTree_RegExpTreeAtom[T interface {
	RegExpTree_RegExpTreeAtom
	comparable
}](v T) RegExpTree_RegExpTreeAtom {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_RegExpTree_RegExpTreeAtom[T RegExpTree_RegExpTreeAtom](v T) RegExpTree_RegExpTreeAtom {
	var vi RegExpTree_RegExpTreeAtom = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.RegExpTreeAtom
type RegExpTree_RegExpTreeAtomBase struct {
	RegExpTreeBase
}

func (this_ *RegExpTree_RegExpTreeAtomBase) BindSelf_(s any) {
	_ = s
	this_.RegExpTreeBase.BindSelf_(s)
}
func (this_ *RegExpTree_RegExpTreeAtomBase) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_RegExpTreeAtomBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$RegExpTreeAtom", "RegExpTreeAtom"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.RegExpTreeAtom", "com.google.javascript.jscomp.regex.RegExpTree$RegExpTreeAtom", func(v_ any) bool {
		return j2gort.InstanceOf[RegExpTree_RegExpTreeAtom](v_)
	})
}
func (b *RegExpTree_RegExpTreeAtomBase) GetRegExpTree_RegExpTreeAtomBase() *RegExpTree_RegExpTreeAtomBase {
	return b
}
func (*RegExpTree_RegExpTreeAtomBase) Implements_regex_RegExpTree_RegExpTreeAtom() {}
func (*RegExpTree_RegExpTreeAtomBase) Implements_regex_RegExpTree()                {}

func (this_ *RegExpTree_RegExpTreeAtomBase) initFields_RegExpTree_RegExpTreeAtom() {
	_ = this_
}

func InitRegExpTree_RegExpTreeAtom(this_ *RegExpTree_RegExpTreeAtomBase) {
	InitRegExpTree(&this_.RegExpTreeBase)
	this_.initFields_RegExpTree_RegExpTreeAtom()
}

// Java method: RegExpTreeAtom.isCaseSensitive()
func (this_ *RegExpTree_RegExpTreeAtomBase) IsCaseSensitive() bool {
	return false
}

// Java method: RegExpTreeAtom.containsAnchor()
func (this_ *RegExpTree_RegExpTreeAtomBase) ContainsAnchor() bool {
	return false
}

// Java method: RegExpTreeAtom.numCapturingGroups()
func (this_ *RegExpTree_RegExpTreeAtomBase) NumCapturingGroups() int32 {
	return int32(0)
}

// Java method: RegExpTreeAtom.children()
func (this_ *RegExpTree_RegExpTreeAtomBase) Children() any {
	return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.Empty
type RegExpTree_Empty struct {
	RegExpTree_RegExpTreeAtomBase
}

func (this_ *RegExpTree_Empty) BindSelf_(s any) {
	_ = s
	this_.RegExpTree_RegExpTreeAtomBase.BindSelf_(s)
}
func (this_ *RegExpTree_Empty) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_Empty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$Empty", "Empty"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.Empty", "com.google.javascript.jscomp.regex.RegExpTree$Empty", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_Empty](v_)
	})
}
func (*RegExpTree_Empty) Implements_regex_RegExpTree_RegExpTreeAtom() {}
func (*RegExpTree_Empty) Implements_regex_RegExpTree()                {}

var static_RegExpTree_Empty_INSTANCE *RegExpTree_Empty

func StaticGet_RegExpTree_Empty_INSTANCE() *RegExpTree_Empty {
	clinit_RegExpTree_Empty()
	return static_RegExpTree_Empty_INSTANCE
}
func StaticSet_RegExpTree_Empty_INSTANCE(v *RegExpTree_Empty) {
	clinit_RegExpTree_Empty()
	static_RegExpTree_Empty_INSTANCE = v
}

var clinitDone_RegExpTree_Empty uint32
var clinitState_RegExpTree_Empty j2gort.ClinitState

func clinit_RegExpTree_Empty() {
	if atomic.LoadUint32(&clinitDone_RegExpTree_Empty) != 0 {
		return
	}
	clinitSlow_RegExpTree_Empty()
}

func clinitSlow_RegExpTree_Empty() {
	if !j2gort.EnterClinit(&clinitState_RegExpTree_Empty, &clinitDone_RegExpTree_Empty) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RegExpTree_Empty, &clinitDone_RegExpTree_Empty)
	static_RegExpTree_Empty_INSTANCE = NewRegExpTree_Empty()
}

func (this_ *RegExpTree_Empty) initFields_RegExpTree_Empty() {
	_ = this_
}

// Java constructor: Empty()
func NewRegExpTree_Empty() *RegExpTree_Empty {
	clinit_RegExpTree_Empty()
	this_ := &RegExpTree_Empty{}
	this_.BindSelf_(this_)
	InitRegExpTree_Empty(this_)
	return this_
}

func InitRegExpTree_Empty(this_ *RegExpTree_Empty) {
	InitRegExpTree_RegExpTreeAtom(&this_.RegExpTree_RegExpTreeAtomBase)
	this_.initFields_RegExpTree_Empty()
}

// Java method: Empty.simplify(java.lang.String flags)
func (this_ *RegExpTree_Empty) Simplify(flags string) RegExpTree {
	_ = flags
	return this_
}

// Java method: Empty.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Empty) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
}

// Java method: Empty.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Empty) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
}

// Java method: Empty.equals(java.lang.Object o)
func (this_ *RegExpTree_Empty) Equals(arg_4 any) bool {
	o := arg_4
	_ = o
	return j2gort.InstanceOf[*RegExpTree_Empty](o)
}

// Java method: Empty.hashCode()
func (this_ *RegExpTree_Empty) HashCode() int32 {
	return int32(2128634177)
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.Anchor
type RegExpTree_Anchor struct {
	RegExpTree_RegExpTreeAtomBase
	type_ uint16 `gson:"type"`
}

func (this_ *RegExpTree_Anchor) BindSelf_(s any) {
	_ = s
	this_.RegExpTree_RegExpTreeAtomBase.BindSelf_(s)
}
func (this_ *RegExpTree_Anchor) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_Anchor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$Anchor", "Anchor"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.Anchor", "com.google.javascript.jscomp.regex.RegExpTree$Anchor", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_Anchor](v_)
	})
}
func (*RegExpTree_Anchor) Implements_regex_RegExpTree_RegExpTreeAtom() {}
func (*RegExpTree_Anchor) Implements_regex_RegExpTree()                {}

func (this_ *RegExpTree_Anchor) initFields_RegExpTree_Anchor() {
	_ = this_
}

// Java constructor: Anchor(char type)
func NewRegExpTree_Anchor(type_ uint16) *RegExpTree_Anchor {
	this_ := &RegExpTree_Anchor{}
	this_.BindSelf_(this_)
	InitRegExpTree_Anchor(this_, type_)
	return this_
}

func InitRegExpTree_Anchor(this_ *RegExpTree_Anchor, type_ uint16) {
	_ = type_
	InitRegExpTree_RegExpTreeAtom(&this_.RegExpTree_RegExpTreeAtomBase)
	this_.initFields_RegExpTree_Anchor()
	this_.type_ = type_
}

// Java method: Anchor.simplify(java.lang.String flags)
func (this_ *RegExpTree_Anchor) Simplify(flags string) RegExpTree {
	_ = flags
	return this_
}

// Java method: Anchor.containsAnchor()
func (this_ *RegExpTree_Anchor) ContainsAnchor() bool {
	return true
}

// Java method: Anchor.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Anchor) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_uint16(this_.type_)
}

// Java method: Anchor.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Anchor) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_uint16(this_.type_)
}

// Java method: Anchor.equals(java.lang.Object o)
func (this_ *RegExpTree_Anchor) Equals(arg_5 any) bool {
	o := arg_5
	_ = o
	var anchor *RegExpTree_Anchor
	_ = anchor
	var ok_6 bool
	anchor, ok_6 = j2gort.TryCast[*RegExpTree_Anchor](o)
	return ok_6 && (int32(this_.type_) == int32(anchor.type_))
}

// Java method: Anchor.hashCode()
func (this_ *RegExpTree_Anchor) HashCode() int32 {
	return (int32(this_.type_) ^ int32(-397207553))
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.WordBoundary
type RegExpTree_WordBoundary struct {
	RegExpTree_RegExpTreeAtomBase
	type_ uint16 `gson:"type"`
}

func (this_ *RegExpTree_WordBoundary) BindSelf_(s any) {
	_ = s
	this_.RegExpTree_RegExpTreeAtomBase.BindSelf_(s)
}
func (this_ *RegExpTree_WordBoundary) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_WordBoundary) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$WordBoundary", "WordBoundary"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.WordBoundary", "com.google.javascript.jscomp.regex.RegExpTree$WordBoundary", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_WordBoundary](v_)
	})
}
func (*RegExpTree_WordBoundary) Implements_regex_RegExpTree_RegExpTreeAtom() {}
func (*RegExpTree_WordBoundary) Implements_regex_RegExpTree()                {}

func (this_ *RegExpTree_WordBoundary) initFields_RegExpTree_WordBoundary() {
	_ = this_
}

// Java constructor: WordBoundary(char type)
func NewRegExpTree_WordBoundary(type_ uint16) *RegExpTree_WordBoundary {
	this_ := &RegExpTree_WordBoundary{}
	this_.BindSelf_(this_)
	InitRegExpTree_WordBoundary(this_, type_)
	return this_
}

func InitRegExpTree_WordBoundary(this_ *RegExpTree_WordBoundary, type_ uint16) {
	_ = type_
	InitRegExpTree_RegExpTreeAtom(&this_.RegExpTree_RegExpTreeAtomBase)
	this_.initFields_RegExpTree_WordBoundary()
	this_.type_ = type_
}

// Java method: WordBoundary.simplify(java.lang.String flags)
func (this_ *RegExpTree_WordBoundary) Simplify(flags string) RegExpTree {
	_ = flags
	return this_
}

// Java method: WordBoundary.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_WordBoundary) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(92))).Append_1_uint16(this_.type_)
}

// Java method: WordBoundary.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_WordBoundary) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_uint16(this_.type_)
}

// Java method: WordBoundary.equals(java.lang.Object o)
func (this_ *RegExpTree_WordBoundary) Equals(arg_7 any) bool {
	o := arg_7
	_ = o
	var wordBoundary *RegExpTree_WordBoundary
	_ = wordBoundary
	var ok_8 bool
	wordBoundary, ok_8 = j2gort.TryCast[*RegExpTree_WordBoundary](o)
	return ok_8 && (int32(this_.type_) == int32(wordBoundary.type_))
}

// Java method: WordBoundary.hashCode()
func (this_ *RegExpTree_WordBoundary) HashCode() int32 {
	return (int32(1450420777) ^ int32(this_.type_))
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.BackReference
type RegExpTree_BackReference struct {
	RegExpTree_RegExpTreeAtomBase
	groupIndex int32 `gson:"groupIndex"`
}

func (this_ *RegExpTree_BackReference) BindSelf_(s any) {
	_ = s
	this_.RegExpTree_RegExpTreeAtomBase.BindSelf_(s)
}
func (this_ *RegExpTree_BackReference) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_BackReference) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$BackReference", "BackReference"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.BackReference", "com.google.javascript.jscomp.regex.RegExpTree$BackReference", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_BackReference](v_)
	})
}
func (*RegExpTree_BackReference) Implements_regex_RegExpTree_RegExpTreeAtom() {}
func (*RegExpTree_BackReference) Implements_regex_RegExpTree()                {}

func (this_ *RegExpTree_BackReference) initFields_RegExpTree_BackReference() {
	_ = this_
}

// Java constructor: BackReference(int groupIndex)
func NewRegExpTree_BackReference(groupIndex int32) *RegExpTree_BackReference {
	this_ := &RegExpTree_BackReference{}
	this_.BindSelf_(this_)
	InitRegExpTree_BackReference(this_, groupIndex)
	return this_
}

func InitRegExpTree_BackReference(this_ *RegExpTree_BackReference, groupIndex int32) {
	_ = groupIndex
	InitRegExpTree_RegExpTreeAtom(&this_.RegExpTree_RegExpTreeAtomBase)
	this_.initFields_RegExpTree_BackReference()
	if !((groupIndex >= int32(0)) && (groupIndex <= int32(99))) {
		j2gort.CheckArgument(false)
	}
	this_.groupIndex = groupIndex
}

// Java method: BackReference.simplify(java.lang.String flags)
func (this_ *RegExpTree_BackReference) Simplify(flags string) RegExpTree {
	_ = flags
	return this_
}

// Java method: BackReference.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_BackReference) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(92))).Append_1_int32(this_.groupIndex)
}

// Java method: BackReference.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_BackReference) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_int32(this_.groupIndex)
}

// Java method: BackReference.equals(java.lang.Object o)
func (this_ *RegExpTree_BackReference) Equals(arg_9 any) bool {
	o := arg_9
	_ = o
	var backReference *RegExpTree_BackReference
	_ = backReference
	var ok_10 bool
	backReference, ok_10 = j2gort.TryCast[*RegExpTree_BackReference](o)
	return ok_10 && (this_.groupIndex == backReference.groupIndex)
}

// Java method: BackReference.hashCode()
func (this_ *RegExpTree_BackReference) HashCode() int32 {
	return (int32(-16308637) ^ this_.groupIndex)
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.NamedBackReference
type RegExpTree_NamedBackReference struct {
	RegExpTree_RegExpTreeAtomBase
	groupName string `gson:"groupName"`
}

func (this_ *RegExpTree_NamedBackReference) BindSelf_(s any) {
	_ = s
	this_.RegExpTree_RegExpTreeAtomBase.BindSelf_(s)
}
func (this_ *RegExpTree_NamedBackReference) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_NamedBackReference) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$NamedBackReference", "NamedBackReference"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.NamedBackReference", "com.google.javascript.jscomp.regex.RegExpTree$NamedBackReference", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_NamedBackReference](v_)
	})
}
func (*RegExpTree_NamedBackReference) Implements_regex_RegExpTree_RegExpTreeAtom() {}
func (*RegExpTree_NamedBackReference) Implements_regex_RegExpTree()                {}

func (this_ *RegExpTree_NamedBackReference) initFields_RegExpTree_NamedBackReference() {
	_ = this_
	this_.groupName = j2gort.NullString
}

// Java constructor: NamedBackReference(java.lang.String groupName)
func NewRegExpTree_NamedBackReference(groupName string) *RegExpTree_NamedBackReference {
	this_ := &RegExpTree_NamedBackReference{}
	this_.BindSelf_(this_)
	InitRegExpTree_NamedBackReference(this_, groupName)
	return this_
}

func InitRegExpTree_NamedBackReference(this_ *RegExpTree_NamedBackReference, groupName string) {
	_ = groupName
	InitRegExpTree_RegExpTreeAtom(&this_.RegExpTree_RegExpTreeAtomBase)
	this_.initFields_RegExpTree_NamedBackReference()
	this_.groupName = groupName
}

// Java method: NamedBackReference.simplify(java.lang.String flags)
func (this_ *RegExpTree_NamedBackReference) Simplify(flags string) RegExpTree {
	_ = flags
	return this_
}

// Java method: NamedBackReference.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_NamedBackReference) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string("\\k<")).Append_1_string(this_.groupName)).Append_1_uint16(uint16(62))
}

// Java method: NamedBackReference.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_NamedBackReference) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_string(this_.groupName)
}

// Java method: NamedBackReference.equals(java.lang.Object o)
func (this_ *RegExpTree_NamedBackReference) Equals(arg_11 any) bool {
	o := arg_11
	_ = o
	var namedBackReference *RegExpTree_NamedBackReference
	_ = namedBackReference
	var ok_12 bool
	namedBackReference, ok_12 = j2gort.TryCast[*RegExpTree_NamedBackReference](o)
	return ok_12 && (func() bool {
		sr_13 := this_.groupName
		se_14 := namedBackReference.groupName
		return j2gort.CheckNotNullString(sr_13) == se_14
	}())
}

// Java method: NamedBackReference.hashCode()
func (this_ *RegExpTree_NamedBackReference) HashCode() int32 {
	return func() int32 {
		hs_15 := this_.groupName
		if hs_15 == j2gort.NullString {
			return 0
		}
		return javalang.StringHashCode(hs_15)
	}()
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.Text
type RegExpTree_Text struct {
	RegExpTree_RegExpTreeAtomBase
	text string `gson:"text"`
}

func (this_ *RegExpTree_Text) BindSelf_(s any) {
	_ = s
	this_.RegExpTree_RegExpTreeAtomBase.BindSelf_(s)
}
func (this_ *RegExpTree_Text) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_Text) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$Text", "Text"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.Text", "com.google.javascript.jscomp.regex.RegExpTree$Text", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_Text](v_)
	})
}
func (*RegExpTree_Text) Implements_regex_RegExpTree_RegExpTreeAtom() {}
func (*RegExpTree_Text) Implements_regex_RegExpTree()                {}

func (this_ *RegExpTree_Text) initFields_RegExpTree_Text() {
	_ = this_
	this_.text = j2gort.NullString
}

// Java constructor: Text(java.lang.String text)
func NewRegExpTree_Text(text string) *RegExpTree_Text {
	this_ := &RegExpTree_Text{}
	this_.BindSelf_(this_)
	InitRegExpTree_Text(this_, text)
	return this_
}

func InitRegExpTree_Text(this_ *RegExpTree_Text, text string) {
	_ = text
	InitRegExpTree_RegExpTreeAtom(&this_.RegExpTree_RegExpTreeAtomBase)
	this_.initFields_RegExpTree_Text()
	this_.text = text
}

// Java static method: Text.escapeRegularCharOnto(char ch, int next, java.lang.StringBuilder sb)
func RegExpTree_Text_EscapeRegularCharOnto(ch uint16, next int32, sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = ch
	_ = next
	_ = sb
	switch ch {
	case uint16(36), uint16(94), uint16(42), uint16(40), uint16(41), uint16(43), uint16(91), uint16(124), uint16(46), uint16(47), uint16(63):
		j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(92))).Append_1_uint16(ch)
	case uint16(123):
		if (int32(uint16(48)) <= next) && (next <= int32(uint16(57))) {
			sb.Append_1_uint16(uint16(92))
		}
		sb.Append_1_uint16(ch)
	default:
		RegExpTree_EscapeCharOnto(ch, sb)
	}
}

// Java method: Text.simplify(java.lang.String flags)
func (this_ *RegExpTree_Text) Simplify(flags string) RegExpTree {
	_ = flags
	var n int32 = j2gort.JStrLen(this_.text)
	_ = n
	if n == int32(0) {
		return ToIface_RegExpTree(StaticGet_RegExpTree_Empty_INSTANCE())
	}
	if j2gort.JStrIndexOfCp(flags, int32(uint16(105))) >= int32(0) {
		var canonicalized string = CaseCanonicalize_CaseCanonicalize_1_string(this_.text)
		_ = canonicalized
		if !(j2gort.CheckNotNullString(this_.text) == canonicalized) {
			return NewRegExpTree_Text(canonicalized)
		}
	}
	return this_
}

// Java method: Text.isCaseSensitive()
func (this_ *RegExpTree_Text) IsCaseSensitive() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = j2gort.JStrLen(this_.text)
		_ = n
		for ; i < n; i = i + 1 {
			if StaticGet_CaseCanonicalize_CASE_SENSITIVE().Contains(int32(j2gort.JStrCharAt(this_.text, i))) {
				return true
			}
		}
	}
	return false
}

// Java method: Text.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Text) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = j2gort.JStrLen(this_.text)
		_ = n
		for ; i < n; i = i + 1 {
			RegExpTree_Text_EscapeRegularCharOnto(j2gort.JStrCharAt(this_.text, i), func() int32 {
				if (i + int32(1)) < n {
					return int32(j2gort.JStrCharAt(this_.text, (i + int32(1))))
				}
				return int32(-1)
			}(), sb)
		}
	}
}

// Java method: Text.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Text) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(96))).Append_1_string(this_.text)).Append_1_uint16(uint16(96))
}

// Java method: Text.equals(java.lang.Object o)
func (this_ *RegExpTree_Text) Equals(arg_16 any) bool {
	o := arg_16
	_ = o
	var text *RegExpTree_Text
	_ = text
	var ok_17 bool
	text, ok_17 = j2gort.TryCast[*RegExpTree_Text](o)
	return ok_17 && (func() bool { sr_18 := this_.text; se_19 := text.text; return j2gort.CheckNotNullString(sr_18) == se_19 }())
}

// Java method: Text.hashCode()
func (this_ *RegExpTree_Text) HashCode() int32 {
	return (javalang.StringHashCode(j2gort.CheckNotNullString(this_.text)) ^ int32(102228752))
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.Repetition
type RegExpTree_Repetition struct {
	RegExpTreeBase
	body   RegExpTree `gson:"body"`
	min_   int32      `gson:"min"`
	max_   int32      `gson:"max"`
	greedy bool       `gson:"greedy"`
}

func (this_ *RegExpTree_Repetition) BindSelf_(s any) {
	_ = s
	this_.RegExpTreeBase.BindSelf_(s)
}
func (this_ *RegExpTree_Repetition) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_Repetition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$Repetition", "Repetition"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.Repetition", "com.google.javascript.jscomp.regex.RegExpTree$Repetition", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_Repetition](v_)
	})
}
func (*RegExpTree_Repetition) Implements_regex_RegExpTree() {}

func (this_ *RegExpTree_Repetition) initFields_RegExpTree_Repetition() {
	_ = this_
}

// Java constructor: Repetition(com.google.javascript.jscomp.regex.RegExpTree body, int min, int max, boolean greedy)
func NewRegExpTree_Repetition(body RegExpTree, min_ int32, max_ int32, greedy bool) *RegExpTree_Repetition {
	this_ := &RegExpTree_Repetition{}
	this_.BindSelf_(this_)
	InitRegExpTree_Repetition(this_, body, min_, max_, greedy)
	return this_
}

func InitRegExpTree_Repetition(this_ *RegExpTree_Repetition, body RegExpTree, min_ int32, max_ int32, greedy bool) {
	_ = body
	_ = min_
	_ = max_
	_ = greedy
	InitRegExpTree(&this_.RegExpTreeBase)
	this_.initFields_RegExpTree_Repetition()
	this_.body = body
	this_.min_ = min_
	this_.max_ = max_
	this_.greedy = greedy
}

// Java method: Repetition.simplify(java.lang.String flags)
func (this_ *RegExpTree_Repetition) Simplify(flags string) RegExpTree {
	_ = flags
	var body RegExpTree = this_.body.Simplify(flags)
	_ = body
	if (this_.max_ == int32(0)) && !(body.GetRegExpTreeBase().HasCapturingGroup()) {
		return ToIface_RegExpTree(StaticGet_RegExpTree_Empty_INSTANCE())
	}
	if j2gort.InstanceOf[*RegExpTree_Empty](body) || StaticGet_RegExpTree_NEVER_MATCHES().Equals(body) {
		return body
	}
	var min_ int32 = this_.min_
	_ = min_
	var max_ int32 = this_.max_
	_ = max_
	var rbody *RegExpTree_Repetition
	_ = rbody
	var ok_20 bool
	rbody, ok_20 = j2gort.TryCast[*RegExpTree_Repetition](body)
	if ok_20 {
		if rbody.greedy == this_.greedy {
			var lmin int64 = ((int64(min_)) * int64(rbody.min_))
			_ = lmin
			var lmax int64 = ((int64(max_)) * int64(rbody.max_))
			_ = lmax
			if lmin < int64(int32(2147483647)) {
				body = rbody.body
				min_ = int32(lmin)
				if lmax >= int64(int32(2147483647)) {
					max_ = int32(2147483647)
				} else {
					max_ = int32(lmax)
				}
			}
		}
	}
	if (min_ == int32(1)) && (max_ == int32(1)) {
		return body
	}
	var greedy bool = (this_.greedy || (min_ == max_))
	_ = greedy
	if ((body.Equals(this_.body) && (min_ == this_.min_)) && (max_ == this_.max_)) && (greedy == this_.greedy) {
		return this_
	} else {
		return NewRegExpTree_Repetition(body, min_, max_, greedy).Simplify(flags)
	}
}

// Java method: Repetition.isCaseSensitive()
func (this_ *RegExpTree_Repetition) IsCaseSensitive() bool {
	return this_.body.IsCaseSensitive()
}

// Java method: Repetition.containsAnchor()
func (this_ *RegExpTree_Repetition) ContainsAnchor() bool {
	return this_.body.ContainsAnchor()
}

// Java method: Repetition.numCapturingGroups()
func (this_ *RegExpTree_Repetition) NumCapturingGroups() int32 {
	return this_.body.NumCapturingGroups()
}

// Java method: Repetition.children()
func (this_ *RegExpTree_Repetition) Children() any {
	return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(this_.body)
}

// Java method: Repetition.appendBodySourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Repetition) priv_Repetition_appendBodySourceCode_7673(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	var text *RegExpTree_Text
	_ = text
	if ((j2gort.InstanceOf[*RegExpTree_Alternation](this_.body) || j2gort.InstanceOf[*RegExpTree_Concatenation](this_.body)) || j2gort.InstanceOf[*RegExpTree_Repetition](this_.body)) || (j2gort.InstanceOfBind[*RegExpTree_Text](this_.body, &text) && (j2gort.JStrLen(text.text) > int32(1))) {
		sb.Append_1_string("(?:")
		this_.body.AppendSourceCode(sb)
		sb.Append_1_uint16(uint16(41))
	} else {
		this_.body.AppendSourceCode(sb)
	}
}

// Java static method: Repetition.suffixLen(int min, int max)
func RegExpTree_Repetition_SuffixLen(min_ int32, max_ int32) int32 {
	_ = min_
	_ = max_
	if max_ == int32(2147483647) {
		switch min_ {
		case int32(0), int32(1):
			return int32(1)
		default:
			return (int32(3) + RegExpTree_Repetition_NumDecimalDigits(min_))
		}
		panic("unreachable switch expression")
	}
	if (min_ == int32(0)) && (max_ == int32(1)) {
		return int32(1)
	}
	if min_ == max_ {
		if min_ == int32(1) {
			return int32(0)
		}
		return (int32(2) + RegExpTree_Repetition_NumDecimalDigits(min_))
	}
	return ((int32(3) + RegExpTree_Repetition_NumDecimalDigits(min_)) + RegExpTree_Repetition_NumDecimalDigits(max_))
}

// Java static method: Repetition.numDecimalDigits(int n)
func RegExpTree_Repetition_NumDecimalDigits(n int32) int32 {
	_ = n
	if n < int32(0) {
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	var nDigits int32 = int32(1)
	_ = nDigits
	for n >= int32(10) {
		nDigits = nDigits + 1
		n = (n / int32(10))
	}
	return nDigits
}

// Java method: Repetition.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Repetition) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	var bodyStart int32 = sb.Length()
	_ = bodyStart
	this_.priv_Repetition_appendBodySourceCode_7673(sb)
	var bodyEnd int32 = sb.Length()
	_ = bodyEnd
	var bodyLen int32 = (bodyEnd - bodyStart)
	_ = bodyLen
	var min_ int32 = this_.min_
	_ = min_
	var max_ int32 = this_.max_
	_ = max_
	if ((min_ >= int32(2)) && (max_ == int32(2147483647))) || ((max_ - min_) <= int32(1)) {
		var expanded int32
		_ = expanded
		if (min_ == max_) || (max_ == int32(2147483647)) {
			expanded = (min_ - int32(1))
		} else {
			expanded = min_
		}
		var expandedMin int32 = (min_ - expanded)
		_ = expandedMin
		var expandedMax int32
		_ = expandedMax
		if max_ == int32(2147483647) {
			expandedMax = max_
		} else {
			expandedMax = (max_ - expanded)
		}
		var suffixLen int32 = RegExpTree_Repetition_SuffixLen(min_, max_)
		_ = suffixLen
		var expandedSuffixLen int32 = RegExpTree_Repetition_SuffixLen(expandedMin, expandedMax)
		_ = expandedSuffixLen
		if (((bodyLen * expanded) + expandedSuffixLen) < suffixLen) && !(this_.body.GetRegExpTreeBase().HasCapturingGroup()) {
			for func() int32 { old_21 := expanded; nw_22 := old_21 - 1; expanded = nw_22; return nw_22 }() >= int32(0) {
				sb.Append_3_any_int32_int32(j2gort.PtrToAny(sb), bodyStart, bodyEnd)
			}
			min_ = expandedMin
			max_ = expandedMax
		}
	}
	if max_ == int32(2147483647) {
		switch min_ {
		case int32(0):
			sb.Append_1_uint16(uint16(42))
		case int32(1):
			sb.Append_1_uint16(uint16(43))
		default:
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(123))).Append_1_int32(min_)).Append_1_string(",}")
		}
	} else if (min_ == int32(0)) && (max_ == int32(1)) {
		sb.Append_1_uint16(uint16(63))
	} else if min_ == max_ {
		if min_ != int32(1) {
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(123))).Append_1_int32(min_)).Append_1_uint16(uint16(125))
		}
	} else {
		j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(123))).Append_1_int32(min_)).Append_1_uint16(uint16(44))).Append_1_int32(max_)).Append_1_uint16(uint16(125))
	}
	if !(this_.greedy) {
		sb.Append_1_uint16(uint16(63))
	}
}

// Java method: Repetition.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Repetition) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string(" min=")).Append_1_int32(this_.min_)).Append_1_string(", max=")).Append_1_int32(this_.max_)
	if !(this_.greedy) {
		sb.Append_1_string("  not_greedy")
	}
}

// Java method: Repetition.equals(java.lang.Object o)
func (this_ *RegExpTree_Repetition) Equals(arg_23 any) bool {
	o := arg_23
	_ = o
	var that *RegExpTree_Repetition
	_ = that
	var ok_24 bool
	that, ok_24 = j2gort.TryCast[*RegExpTree_Repetition](o)
	if !ok_24 {
		return false
	}
	return (((this_.body.Equals(that.body) && (this_.min_ == that.min_)) && (this_.max_ == that.max_)) && (this_.greedy == that.greedy))
}

// Java method: Repetition.hashCode()
func (this_ *RegExpTree_Repetition) HashCode() int32 {
	return (this_.min_ + (int32(31) * (this_.max_ + (int32(31) * ((func() int32 {
		if this_.greedy {
			return int32(1)
		}
		return int32(0)
	}()) + (int32(31) * this_.body.HashCode()))))))
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.Alternation
type RegExpTree_Alternation struct {
	RegExpTreeBase
	alternatives []RegExpTree `gson:"alternatives"`
}

func (this_ *RegExpTree_Alternation) BindSelf_(s any) {
	_ = s
	this_.RegExpTreeBase.BindSelf_(s)
}
func (this_ *RegExpTree_Alternation) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_Alternation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$Alternation", "Alternation"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.Alternation", "com.google.javascript.jscomp.regex.RegExpTree$Alternation", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_Alternation](v_)
	})
}
func (*RegExpTree_Alternation) Implements_regex_RegExpTree() {}

func (this_ *RegExpTree_Alternation) initFields_RegExpTree_Alternation() {
	_ = this_
}

// Java constructor: Alternation(java.util.List<? extends com.google.javascript.jscomp.regex.RegExpTree> alternatives)
func NewRegExpTree_Alternation(alternatives j2gort.Ext_java_util_List) *RegExpTree_Alternation {
	this_ := &RegExpTree_Alternation{}
	this_.BindSelf_(this_)
	InitRegExpTree_Alternation(this_, alternatives)
	return this_
}

func InitRegExpTree_Alternation(this_ *RegExpTree_Alternation, alternatives j2gort.Ext_java_util_List) {
	_ = alternatives
	InitRegExpTree(&this_.RegExpTreeBase)
	this_.initFields_RegExpTree_Alternation()
	this_.alternatives = j2gort.ImmListCopyToSlice[RegExpTree](j2gort.AsIface[j2gort.Ext_java_util_Collection](alternatives))
}

// Java method: Alternation.simplify(java.lang.String flags)
func (this_ *RegExpTree_Alternation) Simplify(flags string) RegExpTree {
	_ = flags
	var alternatives j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = alternatives
	{

		for _, alternative := range this_.alternatives {

			_ = alternative
			alternative = alternative.Simplify(flags)
			var alternation *RegExpTree_Alternation
			_ = alternation
			var ok_27 bool
			alternation, ok_27 = j2gort.TryCast[*RegExpTree_Alternation](alternative)
			if ok_27 {
				alternatives.AddAll_1(j2gort.ImmListFromTypedSlice(alternation.alternatives))
			} else {
				alternatives.Add_1(alternative)
			}
		}
	}
	var last RegExpTree = nil
	_ = last
	{
		var it j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](alternatives.Iterator())
		_ = it
		for it.HasNext() {
			var alternative_2 RegExpTree = j2gort.AsIface[RegExpTree](it.Next())
			_ = alternative_2
			if alternative_2.Equals(StaticGet_RegExpTree_NEVER_MATCHES()) {
				continue
			}
			if alternative_2.Equals(last) && !(alternative_2.GetRegExpTreeBase().HasCapturingGroup()) {
				it.Remove()
			} else {
				last = alternative_2
			}
		}
	}
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = alternatives.Size()
		_ = n
		for ; i < n; i = i + 1 {
			var alternative_3 RegExpTree = j2gort.AsIface[RegExpTree](alternatives.Get(i))
			_ = alternative_3
			var text *RegExpTree_Text
			_ = text
			if (j2gort.InstanceOfBind[*RegExpTree_Text](alternative_3, &text) && (j2gort.JStrLen(text.text) == int32(1))) || j2gort.InstanceOf[*RegExpTree_Charset](alternative_3) {
				var end int32 = i
				_ = end
				var nCharsets int32 = int32(0)
				_ = nCharsets
				for end < n {
					var follower RegExpTree = j2gort.AsIface[RegExpTree](alternatives.Get(end))
					_ = follower
					var text_2 *RegExpTree_Text
					_ = text_2
					if j2gort.InstanceOf[*RegExpTree_Charset](follower) {
						nCharsets = nCharsets + 1
					} else {
						var ok_28 bool
						text_2, ok_28 = j2gort.TryCast[*RegExpTree_Text](follower)
						if ok_28 {
							ok_28 = (j2gort.JStrLen(text_2.text) == int32(1))
						}
						if !ok_28 {
							break
						}
					}
					end = end + 1
				}
				if ((end - i) >= int32(3)) || ((nCharsets != int32(0)) && ((end - i) >= int32(2))) {
					var members []int32 = make([]int32, ((end - i) - nCharsets))
					_ = members
					var memberIdx int32 = int32(0)
					_ = memberIdx
					var chars *CharRanges = StaticGet_CharRanges_EMPTY()
					_ = chars
					var ieExplicits *CharRanges = StaticGet_CharRanges_EMPTY()
					_ = ieExplicits
					var charAlternatives j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](alternatives.SubList(i, end))
					_ = charAlternatives
					{
						cur_29 := j2gort.NewIterCursor(j2gort.CheckNotNull(charAlternatives))
						for cur_29.HasNext() {
							rawElem_30 := cur_29.Next()
							charAlternative := j2gort.AsIface[RegExpTree](rawElem_30)
							_ = charAlternative
							var text_3 *RegExpTree_Text
							_ = text_3
							var cs *RegExpTree_Charset
							_ = cs
							var ok_31 bool
							text_3, ok_31 = j2gort.TryCast[*RegExpTree_Text](charAlternative)
							if ok_31 {
								var ch uint16 = j2gort.JStrCharAt(text_3.text, int32(0))
								_ = ch
								{
									a_32 := members
									i_33 := func() int32 { old_34 := memberIdx; nw_35 := old_34 + 1; memberIdx = nw_35; return old_34 }()
									a_32[i_33] = int32(ch)
								}
								if StaticGet_RegExpTree_IE_SPEC_ERRORS().Contains(int32(ch)) {
									ieExplicits = ieExplicits.Union(CharRanges_Inclusive(int32(ch), int32(ch)))
								}
							} else {
								var ok_36 bool
								cs, ok_36 = j2gort.TryCast[*RegExpTree_Charset](charAlternative)
								if ok_36 {
									chars = chars.Union(cs.ranges)
									ieExplicits = ieExplicits.Union(cs.ieExplicits)
								}
							}
						}
					}
					chars = chars.Union(CharRanges_WithMembers(members))
					charAlternatives.Clear()
					charAlternatives.Add_1(NewRegExpTree_Charset(chars, ieExplicits).Simplify(flags))
					n = alternatives.Size()
				}
			}
		}
	}
	switch alternatives.Size() {
	case int32(0):
		return ToIface_RegExpTree(StaticGet_RegExpTree_Empty_INSTANCE())
	case int32(1):
		return j2gort.AsIface[RegExpTree](alternatives.Get(int32(0)))
	case int32(2):
		if j2gort.InstanceOf[*RegExpTree_Empty](j2gort.AsIface[RegExpTree](alternatives.Get(int32(1)))) {
			return NewRegExpTree_Repetition(j2gort.AsIface[RegExpTree](alternatives.Get(int32(0))), int32(0), int32(1), true)
		} else if j2gort.InstanceOf[*RegExpTree_Empty](j2gort.AsIface[RegExpTree](alternatives.Get(int32(0)))) {
			return NewRegExpTree_Repetition(j2gort.AsIface[RegExpTree](alternatives.Get(int32(1))), int32(0), int32(1), false)
		}
	}
	if alternatives.Equals(j2gort.ImmListFromTypedSlice(this_.alternatives)) {
		return this_
	} else {
		return NewRegExpTree_Alternation(alternatives)
	}
}

// Java method: Alternation.isCaseSensitive()
func (this_ *RegExpTree_Alternation) IsCaseSensitive() bool {
	{

		for _, alternative := range this_.alternatives {

			_ = alternative
			if alternative.IsCaseSensitive() {
				return true
			}
		}
	}
	return false
}

// Java method: Alternation.containsAnchor()
func (this_ *RegExpTree_Alternation) ContainsAnchor() bool {
	{

		for _, alternative := range this_.alternatives {

			_ = alternative
			if alternative.ContainsAnchor() {
				return true
			}
		}
	}
	return false
}

// Java method: Alternation.numCapturingGroups()
func (this_ *RegExpTree_Alternation) NumCapturingGroups() int32 {
	var n int32 = int32(0)
	_ = n
	{

		for _, alternative := range this_.alternatives {

			_ = alternative
			n = (n + alternative.NumCapturingGroups())
		}
	}
	return n
}

// Java method: Alternation.children()
func (this_ *RegExpTree_Alternation) Children() any {
	return j2gort.ImmListFromTypedSlice(this_.alternatives)
}

// Java method: Alternation.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Alternation) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(this_.alternatives))
		_ = n
		for ; i < n; i = i + 1 {
			if i != int32(0) {
				sb.Append_1_uint16(uint16(124))
			}
			this_.alternatives[i].AppendSourceCode(sb)
		}
	}
}

// Java method: Alternation.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Alternation) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
}

// Java method: Alternation.equals(java.lang.Object o)
func (this_ *RegExpTree_Alternation) Equals(arg_43 any) bool {
	o := arg_43
	_ = o
	var alternation *RegExpTree_Alternation
	_ = alternation
	return (j2gort.RefEquals(this_, o) || ((j2gort.InstanceOfBind[*RegExpTree_Alternation](o, &alternation)) && j2gort.ImmSliceEquals(this_.alternatives, j2gort.ImmListFromTypedSlice(alternation.alternatives))))
}

// Java method: Alternation.hashCode()
func (this_ *RegExpTree_Alternation) HashCode() int32 {
	return (int32(1370848465) ^ j2gort.ImmSliceHashCode(this_.alternatives))
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.LookaheadAssertion
type RegExpTree_LookaheadAssertion struct {
	RegExpTreeBase
	body     RegExpTree `gson:"body"`
	positive bool       `gson:"positive"`
}

func (this_ *RegExpTree_LookaheadAssertion) BindSelf_(s any) {
	_ = s
	this_.RegExpTreeBase.BindSelf_(s)
}
func (this_ *RegExpTree_LookaheadAssertion) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_LookaheadAssertion) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$LookaheadAssertion", "LookaheadAssertion"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.LookaheadAssertion", "com.google.javascript.jscomp.regex.RegExpTree$LookaheadAssertion", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_LookaheadAssertion](v_)
	})
}
func (*RegExpTree_LookaheadAssertion) Implements_regex_RegExpTree() {}

func (this_ *RegExpTree_LookaheadAssertion) initFields_RegExpTree_LookaheadAssertion() {
	_ = this_
}

// Java constructor: LookaheadAssertion(com.google.javascript.jscomp.regex.RegExpTree body, boolean positive)
func NewRegExpTree_LookaheadAssertion(body RegExpTree, positive bool) *RegExpTree_LookaheadAssertion {
	this_ := &RegExpTree_LookaheadAssertion{}
	this_.BindSelf_(this_)
	InitRegExpTree_LookaheadAssertion(this_, body, positive)
	return this_
}

func InitRegExpTree_LookaheadAssertion(this_ *RegExpTree_LookaheadAssertion, body RegExpTree, positive bool) {
	_ = body
	_ = positive
	InitRegExpTree(&this_.RegExpTreeBase)
	this_.initFields_RegExpTree_LookaheadAssertion()
	this_.body = body
	this_.positive = positive
}

// Java method: LookaheadAssertion.simplify(java.lang.String flags)
func (this_ *RegExpTree_LookaheadAssertion) Simplify(flags string) RegExpTree {
	_ = flags
	var simpleBody RegExpTree = this_.body.Simplify(flags)
	_ = simpleBody
	if j2gort.InstanceOf[*RegExpTree_Empty](simpleBody) {
		if this_.positive {
			return simpleBody
		}
	}
	return NewRegExpTree_LookaheadAssertion(simpleBody, this_.positive)
}

// Java method: LookaheadAssertion.isCaseSensitive()
func (this_ *RegExpTree_LookaheadAssertion) IsCaseSensitive() bool {
	return this_.body.IsCaseSensitive()
}

// Java method: LookaheadAssertion.containsAnchor()
func (this_ *RegExpTree_LookaheadAssertion) ContainsAnchor() bool {
	return this_.body.ContainsAnchor()
}

// Java method: LookaheadAssertion.numCapturingGroups()
func (this_ *RegExpTree_LookaheadAssertion) NumCapturingGroups() int32 {
	return this_.body.NumCapturingGroups()
}

// Java method: LookaheadAssertion.children()
func (this_ *RegExpTree_LookaheadAssertion) Children() any {
	return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(this_.body)
}

// Java method: LookaheadAssertion.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_LookaheadAssertion) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_string(func() string {
		if this_.positive {
			return "(?="
		}
		return "(?!"
	}())
	this_.body.AppendSourceCode(sb)
	sb.Append_1_uint16(uint16(41))
}

// Java method: LookaheadAssertion.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_LookaheadAssertion) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_string(func() string {
		if this_.positive {
			return "positive"
		}
		return "negative"
	}())
}

// Java method: LookaheadAssertion.equals(java.lang.Object o)
func (this_ *RegExpTree_LookaheadAssertion) Equals(arg_44 any) bool {
	o := arg_44
	_ = o
	var that *RegExpTree_LookaheadAssertion
	_ = that
	var ok_45 bool
	that, ok_45 = j2gort.TryCast[*RegExpTree_LookaheadAssertion](o)
	if !ok_45 {
		return false
	}
	return ((this_.positive == that.positive) && this_.body.Equals(that.body))
}

// Java method: LookaheadAssertion.hashCode()
func (this_ *RegExpTree_LookaheadAssertion) HashCode() int32 {
	return (int32(119778217) ^ this_.body.HashCode())
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.LookbehindAssertion
type RegExpTree_LookbehindAssertion struct {
	RegExpTreeBase
	body     RegExpTree `gson:"body"`
	positive bool       `gson:"positive"`
}

func (this_ *RegExpTree_LookbehindAssertion) BindSelf_(s any) {
	_ = s
	this_.RegExpTreeBase.BindSelf_(s)
}
func (this_ *RegExpTree_LookbehindAssertion) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_LookbehindAssertion) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$LookbehindAssertion", "LookbehindAssertion"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.LookbehindAssertion", "com.google.javascript.jscomp.regex.RegExpTree$LookbehindAssertion", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_LookbehindAssertion](v_)
	})
}
func (*RegExpTree_LookbehindAssertion) Implements_regex_RegExpTree() {}

func (this_ *RegExpTree_LookbehindAssertion) initFields_RegExpTree_LookbehindAssertion() {
	_ = this_
}

// Java constructor: LookbehindAssertion(com.google.javascript.jscomp.regex.RegExpTree body, boolean positive)
func NewRegExpTree_LookbehindAssertion(body RegExpTree, positive bool) *RegExpTree_LookbehindAssertion {
	this_ := &RegExpTree_LookbehindAssertion{}
	this_.BindSelf_(this_)
	InitRegExpTree_LookbehindAssertion(this_, body, positive)
	return this_
}

func InitRegExpTree_LookbehindAssertion(this_ *RegExpTree_LookbehindAssertion, body RegExpTree, positive bool) {
	_ = body
	_ = positive
	InitRegExpTree(&this_.RegExpTreeBase)
	this_.initFields_RegExpTree_LookbehindAssertion()
	this_.body = body
	this_.positive = positive
}

// Java method: LookbehindAssertion.simplify(java.lang.String flags)
func (this_ *RegExpTree_LookbehindAssertion) Simplify(flags string) RegExpTree {
	_ = flags
	var simpleBody RegExpTree = this_.body.Simplify(flags)
	_ = simpleBody
	if j2gort.InstanceOf[*RegExpTree_Empty](simpleBody) {
		if this_.positive {
			return simpleBody
		}
	}
	return NewRegExpTree_LookbehindAssertion(simpleBody, this_.positive)
}

// Java method: LookbehindAssertion.isCaseSensitive()
func (this_ *RegExpTree_LookbehindAssertion) IsCaseSensitive() bool {
	return this_.body.IsCaseSensitive()
}

// Java method: LookbehindAssertion.containsAnchor()
func (this_ *RegExpTree_LookbehindAssertion) ContainsAnchor() bool {
	return this_.body.ContainsAnchor()
}

// Java method: LookbehindAssertion.numCapturingGroups()
func (this_ *RegExpTree_LookbehindAssertion) NumCapturingGroups() int32 {
	return this_.body.NumCapturingGroups()
}

// Java method: LookbehindAssertion.children()
func (this_ *RegExpTree_LookbehindAssertion) Children() any {
	return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(this_.body)
}

// Java method: LookbehindAssertion.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_LookbehindAssertion) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_string(func() string {
		if this_.positive {
			return "(?<="
		}
		return "(?<!"
	}())
	this_.body.AppendSourceCode(sb)
	sb.Append_1_uint16(uint16(41))
}

// Java method: LookbehindAssertion.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_LookbehindAssertion) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_string(func() string {
		if this_.positive {
			return "positive"
		}
		return "negative"
	}())
}

// Java method: LookbehindAssertion.equals(java.lang.Object o)
func (this_ *RegExpTree_LookbehindAssertion) Equals(arg_46 any) bool {
	o := arg_46
	_ = o
	var that *RegExpTree_LookbehindAssertion
	_ = that
	var ok_47 bool
	that, ok_47 = j2gort.TryCast[*RegExpTree_LookbehindAssertion](o)
	if !ok_47 {
		return false
	}
	return ((this_.positive == that.positive) && this_.body.Equals(that.body))
}

// Java method: LookbehindAssertion.hashCode()
func (this_ *RegExpTree_LookbehindAssertion) HashCode() int32 {
	return (int32(119778217) ^ this_.body.HashCode())
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.CapturingGroup
type RegExpTree_CapturingGroup struct {
	RegExpTreeBase
	body RegExpTree `gson:"body"`
}

func (this_ *RegExpTree_CapturingGroup) BindSelf_(s any) {
	_ = s
	this_.RegExpTreeBase.BindSelf_(s)
}
func (this_ *RegExpTree_CapturingGroup) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_CapturingGroup) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$CapturingGroup", "CapturingGroup"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.CapturingGroup", "com.google.javascript.jscomp.regex.RegExpTree$CapturingGroup", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_CapturingGroup](v_)
	})
}
func (*RegExpTree_CapturingGroup) Implements_regex_RegExpTree() {}

func (this_ *RegExpTree_CapturingGroup) initFields_RegExpTree_CapturingGroup() {
	_ = this_
}

// Java constructor: CapturingGroup(com.google.javascript.jscomp.regex.RegExpTree body)
func NewRegExpTree_CapturingGroup(body RegExpTree) *RegExpTree_CapturingGroup {
	this_ := &RegExpTree_CapturingGroup{}
	this_.BindSelf_(this_)
	InitRegExpTree_CapturingGroup(this_, body)
	return this_
}

func InitRegExpTree_CapturingGroup(this_ *RegExpTree_CapturingGroup, body RegExpTree) {
	_ = body
	InitRegExpTree(&this_.RegExpTreeBase)
	this_.initFields_RegExpTree_CapturingGroup()
	this_.body = body
}

// Java method: CapturingGroup.simplify(java.lang.String flags)
func (this_ *RegExpTree_CapturingGroup) Simplify(flags string) RegExpTree {
	_ = flags
	return NewRegExpTree_CapturingGroup(this_.body.Simplify(flags))
}

// Java method: CapturingGroup.isCaseSensitive()
func (this_ *RegExpTree_CapturingGroup) IsCaseSensitive() bool {
	return this_.body.IsCaseSensitive()
}

// Java method: CapturingGroup.containsAnchor()
func (this_ *RegExpTree_CapturingGroup) ContainsAnchor() bool {
	return this_.body.ContainsAnchor()
}

// Java method: CapturingGroup.numCapturingGroups()
func (this_ *RegExpTree_CapturingGroup) NumCapturingGroups() int32 {
	return (int32(1) + this_.body.NumCapturingGroups())
}

// Java method: CapturingGroup.children()
func (this_ *RegExpTree_CapturingGroup) Children() any {
	return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(this_.body)
}

// Java method: CapturingGroup.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_CapturingGroup) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_uint16(uint16(40))
	this_.body.AppendSourceCode(sb)
	sb.Append_1_uint16(uint16(41))
}

// Java method: CapturingGroup.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_CapturingGroup) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
}

// Java method: CapturingGroup.equals(java.lang.Object o)
func (this_ *RegExpTree_CapturingGroup) Equals(arg_48 any) bool {
	o := arg_48
	_ = o
	var capturingGroup *RegExpTree_CapturingGroup
	_ = capturingGroup
	var ok_49 bool
	capturingGroup, ok_49 = j2gort.TryCast[*RegExpTree_CapturingGroup](o)
	return ok_49 && (this_.body.Equals(capturingGroup.body))
}

// Java method: CapturingGroup.hashCode()
func (this_ *RegExpTree_CapturingGroup) HashCode() int32 {
	return (int32(1433933624) ^ this_.body.HashCode())
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.NamedCaptureGroup
type RegExpTree_NamedCaptureGroup struct {
	RegExpTreeBase
	body RegExpTree `gson:"body"`
	name string     `gson:"name"`
}

func (this_ *RegExpTree_NamedCaptureGroup) BindSelf_(s any) {
	_ = s
	this_.RegExpTreeBase.BindSelf_(s)
}
func (this_ *RegExpTree_NamedCaptureGroup) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_NamedCaptureGroup) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$NamedCaptureGroup", "NamedCaptureGroup"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.NamedCaptureGroup", "com.google.javascript.jscomp.regex.RegExpTree$NamedCaptureGroup", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_NamedCaptureGroup](v_)
	})
}
func (*RegExpTree_NamedCaptureGroup) Implements_regex_RegExpTree() {}

func (this_ *RegExpTree_NamedCaptureGroup) initFields_RegExpTree_NamedCaptureGroup() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: NamedCaptureGroup(com.google.javascript.jscomp.regex.RegExpTree body, java.lang.String name)
func NewRegExpTree_NamedCaptureGroup(body RegExpTree, name string) *RegExpTree_NamedCaptureGroup {
	this_ := &RegExpTree_NamedCaptureGroup{}
	this_.BindSelf_(this_)
	InitRegExpTree_NamedCaptureGroup(this_, body, name)
	return this_
}

func InitRegExpTree_NamedCaptureGroup(this_ *RegExpTree_NamedCaptureGroup, body RegExpTree, name string) {
	_ = body
	_ = name
	InitRegExpTree(&this_.RegExpTreeBase)
	this_.initFields_RegExpTree_NamedCaptureGroup()
	this_.body = body
	this_.name = name
}

// Java method: NamedCaptureGroup.simplify(java.lang.String flags)
func (this_ *RegExpTree_NamedCaptureGroup) Simplify(flags string) RegExpTree {
	_ = flags
	return NewRegExpTree_NamedCaptureGroup(this_.body.Simplify(flags), this_.name)
}

// Java method: NamedCaptureGroup.isCaseSensitive()
func (this_ *RegExpTree_NamedCaptureGroup) IsCaseSensitive() bool {
	return this_.body.IsCaseSensitive()
}

// Java method: NamedCaptureGroup.containsAnchor()
func (this_ *RegExpTree_NamedCaptureGroup) ContainsAnchor() bool {
	return this_.body.ContainsAnchor()
}

// Java method: NamedCaptureGroup.numCapturingGroups()
func (this_ *RegExpTree_NamedCaptureGroup) NumCapturingGroups() int32 {
	return (int32(1) + this_.body.NumCapturingGroups())
}

// Java method: NamedCaptureGroup.children()
func (this_ *RegExpTree_NamedCaptureGroup) Children() any {
	return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(this_.body)
}

// Java method: NamedCaptureGroup.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_NamedCaptureGroup) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_string("(?<")
	sb.Append_1_string(this_.name)
	sb.Append_1_uint16(uint16(62))
	this_.body.AppendSourceCode(sb)
	sb.Append_1_uint16(uint16(41))
}

// Java method: NamedCaptureGroup.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_NamedCaptureGroup) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_string(" name=")).Append_1_string(this_.name)
}

// Java method: NamedCaptureGroup.equals(java.lang.Object o)
func (this_ *RegExpTree_NamedCaptureGroup) Equals(arg_50 any) bool {
	o := arg_50
	_ = o
	var namedCaptureGroup *RegExpTree_NamedCaptureGroup
	_ = namedCaptureGroup
	var ok_51 bool
	namedCaptureGroup, ok_51 = j2gort.TryCast[*RegExpTree_NamedCaptureGroup](o)
	return ok_51 && (func() bool {
		sr_52 := this_.name
		se_53 := namedCaptureGroup.name
		return j2gort.CheckNotNullString(sr_52) == se_53
	}()) && (this_.body.Equals(namedCaptureGroup.body))
}

// Java method: NamedCaptureGroup.hashCode()
func (this_ *RegExpTree_NamedCaptureGroup) HashCode() int32 {
	return (func() int32 {
		hs_54 := this_.name
		if hs_54 == j2gort.NullString {
			return 0
		}
		return javalang.StringHashCode(hs_54)
	}() ^ this_.body.HashCode())
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.UnicodePropertyEscape
type RegExpTree_UnicodePropertyEscape struct {
	RegExpTree_RegExpTreeAtomBase
	propertyName  string `gson:"propertyName"`
	propertyValue string `gson:"propertyValue"`
	negated       bool   `gson:"negated"`
}

func (this_ *RegExpTree_UnicodePropertyEscape) BindSelf_(s any) {
	_ = s
	this_.RegExpTree_RegExpTreeAtomBase.BindSelf_(s)
}
func (this_ *RegExpTree_UnicodePropertyEscape) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_UnicodePropertyEscape) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$UnicodePropertyEscape", "UnicodePropertyEscape"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.UnicodePropertyEscape", "com.google.javascript.jscomp.regex.RegExpTree$UnicodePropertyEscape", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_UnicodePropertyEscape](v_)
	})
}
func (*RegExpTree_UnicodePropertyEscape) Implements_regex_RegExpTree_RegExpTreeAtom() {}
func (*RegExpTree_UnicodePropertyEscape) Implements_regex_RegExpTree()                {}

func (this_ *RegExpTree_UnicodePropertyEscape) initFields_RegExpTree_UnicodePropertyEscape() {
	_ = this_
	this_.propertyName = j2gort.NullString
	this_.propertyValue = j2gort.NullString
}

// Java constructor: UnicodePropertyEscape(java.lang.@org.jspecify.annotations.Nullable String propertyName, java.lang.String propertyValue, boolean negated)
func NewRegExpTree_UnicodePropertyEscape(propertyName string, propertyValue string, negated bool) *RegExpTree_UnicodePropertyEscape {
	this_ := &RegExpTree_UnicodePropertyEscape{}
	this_.BindSelf_(this_)
	InitRegExpTree_UnicodePropertyEscape(this_, propertyName, propertyValue, negated)
	return this_
}

func InitRegExpTree_UnicodePropertyEscape(this_ *RegExpTree_UnicodePropertyEscape, propertyName string, propertyValue string, negated bool) {
	_ = propertyName
	_ = propertyValue
	_ = negated
	InitRegExpTree_RegExpTreeAtom(&this_.RegExpTree_RegExpTreeAtomBase)
	this_.initFields_RegExpTree_UnicodePropertyEscape()
	if !(propertyValue != j2gort.NullString) {
		j2gort.CheckState(false)
	}
	if !((propertyName == j2gort.NullString) || !(j2gort.CheckNotNullString(propertyName) == "")) {
		j2gort.CheckArgumentStr(false, "if '=' is present in a unicode property escape, the name cannot be empty")
	}
	if !(!(j2gort.CheckNotNullString(propertyValue) == "")) {
		j2gort.CheckArgumentStr(false, "unicode property escape value cannot be empty")
	}
	this_.propertyName = propertyName
	this_.propertyValue = propertyValue
	this_.negated = negated
}

// Java method: UnicodePropertyEscape.simplify(java.lang.String flags)
func (this_ *RegExpTree_UnicodePropertyEscape) Simplify(flags string) RegExpTree {
	_ = flags
	return this_
}

// Java method: UnicodePropertyEscape.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_UnicodePropertyEscape) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_string(func() string {
		if this_.negated {
			return "\\P{"
		}
		return "\\p{"
	}())
	if this_.propertyName != j2gort.NullString {
		sb.Append_1_string(this_.propertyName)
		sb.Append_1_uint16(uint16(61))
	}
	sb.Append_1_string(this_.propertyValue)
	sb.Append_1_string("}")
}

// Java method: UnicodePropertyEscape.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_UnicodePropertyEscape) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
}

// Java method: UnicodePropertyEscape.equals(java.lang.Object o)
func (this_ *RegExpTree_UnicodePropertyEscape) Equals(arg_55 any) bool {
	o := arg_55
	_ = o
	var unicodePropertyEscape *RegExpTree_UnicodePropertyEscape
	_ = unicodePropertyEscape
	var ok_56 bool
	unicodePropertyEscape, ok_56 = j2gort.TryCast[*RegExpTree_UnicodePropertyEscape](o)
	return ok_56 && (this_.negated == unicodePropertyEscape.negated) && (this_.propertyName == unicodePropertyEscape.propertyName) && (this_.propertyValue == unicodePropertyEscape.propertyValue)
}

// Java method: UnicodePropertyEscape.hashCode()
func (this_ *RegExpTree_UnicodePropertyEscape) HashCode() int32 {
	return func() int32 {
		var oh_57 int32 = 1
		if this_.negated {
			oh_57 = 31*oh_57 + 1231
		} else {
			oh_57 = 31*oh_57 + 1237
		}
		{
			hs_58 := this_.propertyName
			oh_57 = 31 * oh_57
			if hs_58 != j2gort.NullString {
				oh_57 += javalang.StringHashCode(hs_58)
			}
		}
		{
			hs_59 := this_.propertyValue
			oh_57 = 31 * oh_57
			if hs_59 != j2gort.NullString {
				oh_57 += javalang.StringHashCode(hs_59)
			}
		}
		return oh_57
	}()
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.Charset
type RegExpTree_Charset struct {
	RegExpTree_RegExpTreeAtomBase
	ranges      *CharRanges `gson:"ranges"`
	ieExplicits *CharRanges `gson:"ieExplicits"`
}

func (this_ *RegExpTree_Charset) BindSelf_(s any) {
	_ = s
	this_.RegExpTree_RegExpTreeAtomBase.BindSelf_(s)
}
func (this_ *RegExpTree_Charset) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_Charset) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$Charset", "Charset"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.Charset", "com.google.javascript.jscomp.regex.RegExpTree$Charset", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_Charset](v_)
	})
}
func (*RegExpTree_Charset) Implements_regex_RegExpTree_RegExpTreeAtom() {}
func (*RegExpTree_Charset) Implements_regex_RegExpTree()                {}

func (this_ *RegExpTree_Charset) initFields_RegExpTree_Charset() {
	_ = this_
}

// Java constructor: Charset(com.google.javascript.jscomp.regex.CharRanges ranges, com.google.javascript.jscomp.regex.CharRanges ieExplicits)
func NewRegExpTree_Charset(ranges *CharRanges, ieExplicits *CharRanges) *RegExpTree_Charset {
	this_ := &RegExpTree_Charset{}
	this_.BindSelf_(this_)
	InitRegExpTree_Charset(this_, ranges, ieExplicits)
	return this_
}

func InitRegExpTree_Charset(this_ *RegExpTree_Charset, ranges *CharRanges, ieExplicits *CharRanges) {
	_ = ranges
	_ = ieExplicits
	InitRegExpTree_RegExpTreeAtom(&this_.RegExpTree_RegExpTreeAtomBase)
	this_.initFields_RegExpTree_Charset()
	this_.ranges = ranges
	this_.ieExplicits = ieExplicits
}

// Java static method: Charset.complexityWordFolded(com.google.javascript.jscomp.regex.CharRanges ranges)
func RegExpTree_Charset_ComplexityWordFolded(ranges *CharRanges) int32 {
	_ = ranges
	return min(RegExpTree_Charset_ComplexityWordFoldedHelper(ranges), (int32(1) + RegExpTree_Charset_ComplexityWordFoldedHelper(StaticGet_CharRanges_ALL_CODE_UNITS().Difference(ranges))))
}

// Java static method: Charset.complexityWordFoldedHelper(com.google.javascript.jscomp.regex.CharRanges ranges)
func RegExpTree_Charset_ComplexityWordFoldedHelper(ranges *CharRanges) int32 {
	_ = ranges
	var complexity int32 = RegExpTree_DecomposedCharset_Complexity(ranges)
	_ = complexity
	if ranges.ContainsAll(StaticGet_RegExpTree_WORD_CHARS()) {
		complexity = min(complexity, (int32(1) + RegExpTree_DecomposedCharset_Complexity(ranges.Difference(StaticGet_RegExpTree_WORD_CHARS()))))
	}
	if ranges.ContainsAll(StaticGet_RegExpTree_INVERSE_WORD_CHARS()) {
		complexity = min(complexity, (int32(1) + RegExpTree_DecomposedCharset_Complexity(ranges.Difference(StaticGet_RegExpTree_INVERSE_WORD_CHARS()))))
	}
	return complexity
}

// Java method: Charset.simplify(java.lang.String flags)
func (this_ *RegExpTree_Charset) Simplify(flags string) RegExpTree {
	_ = flags
	if this_.ranges.IsEmpty() {
		return StaticGet_RegExpTree_NEVER_MATCHES()
	}
	var best *CharRanges = this_.ranges
	_ = best
	if j2gort.JStrIndexOfCp(flags, int32(uint16(105))) >= int32(0) {
		var options j2gort.CustomOrderedSet[*CharRanges] = j2gort.NewCustomOrderedSet[*CharRanges]()
		_ = options
		options.Add(CaseCanonicalize_ExpandToAllMatched(this_.ranges))
		options.Add(CaseCanonicalize_ReduceToMinimum(this_.ranges))
		var lcaseLetters *CharRanges = this_.ranges.Intersection(StaticGet_RegExpTree_LCASE_LETTERS())
		_ = lcaseLetters
		var ucaseLetters *CharRanges = this_.ranges.Intersection(StaticGet_RegExpTree_UCASE_LETTERS())
		_ = ucaseLetters
		var lcaseLettersToUpper *CharRanges = lcaseLetters.Shift(int32(-32))
		_ = lcaseLettersToUpper
		var ucaseLettersToLower *CharRanges = ucaseLetters.Shift(int32(32))
		_ = ucaseLettersToLower
		options.Add(this_.ranges.Union(ucaseLettersToLower))
		options.Add(this_.ranges.Union(lcaseLettersToUpper))
		options.Add(this_.ranges.Union(lcaseLettersToUpper).Union(ucaseLettersToLower))
		options.Add(this_.ranges.Union(ucaseLettersToLower).Difference(ucaseLetters))
		options.Add(this_.ranges.Union(lcaseLettersToUpper).Difference(lcaseLetters))
		var bestComplexity int32 = RegExpTree_Charset_ComplexityWordFolded(this_.ranges)
		_ = bestComplexity
		{
			cur_60 := options.Iter()
			for cur_60.HasNext() {

				option := cur_60.
					Next()
				_ = option
				var complexity int32 = RegExpTree_Charset_ComplexityWordFolded(option)
				_ = complexity
				if complexity < bestComplexity {
					bestComplexity = complexity
					best = option
				}
			}
		}
	}
	if (best.GetNumRanges() == int32(1)) && ((best.End(int32(0)) - best.Start(int32(0))) == int32(1)) {
		return NewRegExpTree_Text(javalang.StringFromChar(uint16(best.Start(int32(0)))))
	}
	if !(best.Equals(j2gort.PtrToAny(this_.ranges))) {
		return NewRegExpTree_Charset(best, this_.ieExplicits)
	}
	return this_
}

// Java method: Charset.isCaseSensitive()
func (this_ *RegExpTree_Charset) IsCaseSensitive() bool {
	var withoutNamedGroups *CharRanges = this_.Decompose_0().ranges
	_ = withoutNamedGroups
	return !(withoutNamedGroups.Equals(j2gort.PtrToAny(CaseCanonicalize_ExpandToAllMatched(withoutNamedGroups))))
}

// Java method: Charset.decompose(com.google.javascript.jscomp.regex.CharRanges ranges, boolean inverted)
func (this_ *RegExpTree_Charset) priv_Charset_decompose_7674(ranges *CharRanges, inverted bool) *RegExpTree_DecomposedCharset {
	_ = ranges
	_ = inverted
	var namedGroups j2gort.LocalStringBuilder

	var rangesInterIeExplicits *CharRanges = ranges.Intersection(this_.ieExplicits)
	_ = rangesInterIeExplicits
	for {
		var groupName uint16 = uint16(int32(0))
		_ = groupName
		var simplest *CharRanges = nil
		_ = simplest
		var minComplexity int32 = RegExpTree_DecomposedCharset_Complexity(ranges)
		_ = minComplexity
		{
			cur_62 := j2gort.NewMapEntryCursor(StaticGet_RegExpTree_NAMED_CHAR_GROUPS())
			for cur_62.NextEntry() {
				var group *CharRanges = j2gort.CastPtr[CharRanges](cur_62.GetValue())
				_ = group
				if ranges.ContainsAll(group) {
					var withoutGroup *CharRanges = ranges.Difference(group).Union(rangesInterIeExplicits)
					_ = withoutGroup
					var complexity int32 = RegExpTree_DecomposedCharset_Complexity(withoutGroup)
					_ = complexity
					if complexity < minComplexity {
						simplest = withoutGroup
						groupName = j2gort.CastPtr[j2gort.BoxChar](cur_62.GetKey()).CharValue()
						minComplexity = complexity
					}
				}
			}
		}
		if simplest != nil {
			namedGroups.AppendChar(uint16(92)).AppendChar(groupName)
			ranges = simplest
		} else {
			break
		}
	}
	return NewRegExpTree_DecomposedCharset(inverted, ranges, namedGroups.String())
}

// Java method: Charset.decompose()
func (this_ *RegExpTree_Charset) Decompose_0() *RegExpTree_DecomposedCharset {
	var negRanges *CharRanges = StaticGet_CharRanges_ALL_CODE_UNITS().Difference(this_.ranges)
	_ = negRanges
	if !(this_.ieExplicits.IsEmpty()) {
		if negRanges.Intersection(this_.ieExplicits).IsEmpty() {
			return this_.priv_Charset_decompose_7674(this_.ranges, false)
		} else if this_.ranges.Intersection(this_.ieExplicits).IsEmpty() {
			return this_.priv_Charset_decompose_7674(negRanges, true)
		}
	}
	var positive *RegExpTree_DecomposedCharset = this_.priv_Charset_decompose_7674(this_.ranges, false)
	_ = positive
	var negative *RegExpTree_DecomposedCharset = this_.priv_Charset_decompose_7674(negRanges, true)
	_ = negative
	if positive.Complexity() <= negative.Complexity() {
		return positive
	} else {
		return negative
	}
}

// Java method: Charset.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Charset) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	if StaticGet_RegExpTree_DOT_CHARSET().ranges.Equals(j2gort.PtrToAny(this_.ranges)) {
		sb.Append_1_uint16(uint16(46))
		return
	}
	this_.Decompose_0().AppendSourceCode(sb)
}

// Java method: Charset.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Charset) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	sb.Append_1_any_9850(j2gort.PtrToAny(this_.ranges))
}

// Java method: Charset.equals(java.lang.Object o)
func (this_ *RegExpTree_Charset) Equals(arg_63 any) bool {
	o := arg_63
	_ = o
	var charset *RegExpTree_Charset
	_ = charset
	var ok_64 bool
	charset, ok_64 = j2gort.TryCast[*RegExpTree_Charset](o)
	return ok_64 && (this_.ranges.Equals(j2gort.PtrToAny(charset.ranges)))
}

// Java method: Charset.hashCode()
func (this_ *RegExpTree_Charset) HashCode() int32 {
	return (this_.ranges.HashCode() ^ int32(-555867578))
}

// Java record: com.google.javascript.jscomp.regex.RegExpTree.DecomposedCharset
type RegExpTree_DecomposedCharset struct {
	ranges      *CharRanges `gson:"ranges"`
	namedGroups string      `gson:"namedGroups"`
	inverted    bool        `gson:"inverted"`
}

func (this_ *RegExpTree_DecomposedCharset) BindSelf_(s any) {
	_ = s
}
func (this_ *RegExpTree_DecomposedCharset) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_DecomposedCharset) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$DecomposedCharset", "DecomposedCharset"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.DecomposedCharset", "com.google.javascript.jscomp.regex.RegExpTree$DecomposedCharset", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_DecomposedCharset](v_)
	})
}
func (*RegExpTree_DecomposedCharset) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *RegExpTree_DecomposedCharset) initFields_RegExpTree_DecomposedCharset() {
	_ = this_
}

// Java constructor: DecomposedCharset(boolean inverted, com.google.javascript.jscomp.regex.CharRanges ranges, java.lang.String namedGroups)
func NewRegExpTree_DecomposedCharset(inverted bool, ranges *CharRanges, namedGroups string) *RegExpTree_DecomposedCharset {
	this_ := &RegExpTree_DecomposedCharset{}
	InitRegExpTree_DecomposedCharset(this_, inverted, ranges, namedGroups)
	return this_
}

func InitRegExpTree_DecomposedCharset(this_ *RegExpTree_DecomposedCharset, inverted bool, ranges *CharRanges, namedGroups string) {
	_ = inverted
	_ = ranges
	_ = namedGroups
	this_.initFields_RegExpTree_DecomposedCharset()
	this_.inverted = inverted
	this_.ranges = ranges
	this_.namedGroups = namedGroups
}

// Java method: DecomposedCharset.complexity()
func (this_ *RegExpTree_DecomposedCharset) Complexity() int32 {
	return (((func() int32 {
		if this_.inverted {
			return int32(1)
		}
		return int32(0)
	}()) + j2gort.JStrLen(this_.namedGroups)) + RegExpTree_DecomposedCharset_Complexity(this_.ranges))
}

// Java method: DecomposedCharset.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_DecomposedCharset) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	if this_.ranges.IsEmpty() {
		if !(this_.inverted) && (j2gort.JStrLen(this_.namedGroups) == int32(2)) {
			sb.Append_1_string(this_.namedGroups)
			return
		} else if this_.ranges.IsEmpty() && (j2gort.CheckNotNullString(this_.namedGroups) == "") {
			sb.Append_1_string(func() string {
				if this_.inverted {
					return "[\\S\\s]"
				}
				return "(?!)"
			}())
			return
		}
	}
	sb.Append_1_uint16(uint16(91))
	if this_.inverted {
		sb.Append_1_uint16(uint16(94))
	}
	sb.Append_1_string(this_.namedGroups)
	var rangesStartCharset bool = (!(this_.inverted) && (j2gort.CheckNotNullString(this_.namedGroups) == ""))
	_ = rangesStartCharset
	var emitDashAtEnd bool = false
	_ = emitDashAtEnd
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.ranges.GetNumRanges()
		_ = n
		for ; i < n; i = i + 1 {
			var start uint16 = uint16(this_.ranges.Start(i))
			_ = start
			var end uint16 = uint16((this_.ranges.End(i) - int32(1)))
			_ = end
			switch int32(end) - int32(start) {
			case int32(0):
				if int32(start) == int32(uint16(45)) {
					emitDashAtEnd = true
				} else {
					RegExpTree_DecomposedCharset_EscapeRangeCharOnto(start, rangesStartCharset, (i == int32(0)), ((i + int32(1)) == n), sb)
				}
			case int32(1):
				RegExpTree_DecomposedCharset_EscapeRangeCharOnto(start, rangesStartCharset, (i == int32(0)), false, sb)
				RegExpTree_DecomposedCharset_EscapeRangeCharOnto(end, rangesStartCharset, false, ((i + int32(1)) == n), sb)
			default:
				RegExpTree_DecomposedCharset_EscapeRangeCharOnto(start, rangesStartCharset, (i == int32(0)), false, sb)
				sb.Append_1_uint16(uint16(45))
				RegExpTree_DecomposedCharset_EscapeRangeCharOnto(end, rangesStartCharset, false, true, sb)
			}
		}
	}
	if emitDashAtEnd {
		sb.Append_1_uint16(uint16(45))
	}
	sb.Append_1_uint16(uint16(93))
}

// Java static method: DecomposedCharset.escapeRangeCharOnto(char ch, boolean startIsFlush, boolean atStart, boolean atEnd, java.lang.StringBuilder sb)
func RegExpTree_DecomposedCharset_EscapeRangeCharOnto(ch uint16, startIsFlush bool, atStart bool, atEnd bool, sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = ch
	_ = startIsFlush
	_ = atStart
	_ = atEnd
	_ = sb
	switch ch {
	case uint16(8):
		sb.Append_1_string("\\b")
	case uint16(94):
		sb.Append_1_string(func() string {
			if atStart && startIsFlush {
				return "\\^"
			}
			return "^"
		}())
	case uint16(45):
		sb.Append_1_string(func() string {
			if atStart || atEnd {
				return "-"
			}
			return "\\-"
		}())
	case uint16(92), uint16(93):
		j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_uint16(uint16(92))).Append_1_uint16(ch)
	default:
		RegExpTree_EscapeCharOnto(ch, sb)
	}
}

// Java static method: DecomposedCharset.complexity(com.google.javascript.jscomp.regex.CharRanges ranges)
func RegExpTree_DecomposedCharset_Complexity(ranges *CharRanges) int32 {
	_ = ranges
	var complexity int32 = int32(0)
	_ = complexity
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = ranges.GetNumRanges()
		_ = n
		for ; i < n; i = i + 1 {
			var start int32 = ranges.Start(i)
			_ = start
			var end int32 = (ranges.End(i) - int32(1))
			_ = end
			if (start < int32(32)) || (start >= int32(127)) {
				complexity = (complexity + func() int32 {
					if start >= int32(256) {
						return int32(6)
					}
					return int32(4)
				}())
			} else {
				complexity = complexity + 1
			}
			switch end - start {
			case int32(0):
				continue
			case int32(1):
			default:
				complexity = (complexity + int32(1))
			}
			if (end < int32(32)) || (end >= int32(127)) {
				complexity = (complexity + func() int32 {
					if end >= int32(256) {
						return int32(6)
					}
					return int32(4)
				}())
			} else {
				complexity = complexity + 1
			}
		}
	}
	return complexity
}

func (this_ *RegExpTree_DecomposedCharset) Inverted() bool      { return this_.inverted }
func (this_ *RegExpTree_DecomposedCharset) Ranges() *CharRanges { return this_.ranges }
func (this_ *RegExpTree_DecomposedCharset) NamedGroups() string { return this_.namedGroups }
func (this_ *RegExpTree_DecomposedCharset) Equals(other_ any) bool {
	o_, ok_ := other_.(*RegExpTree_DecomposedCharset)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.inverted != o_.inverted {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.ranges), j2gort.PtrToAny(o_.ranges)) {
		return false
	}
	if this_.namedGroups != o_.namedGroups {
		return false
	}
	return true
}
func (this_ *RegExpTree_DecomposedCharset) HashCode() int32 {
	var h_ int32 = 0
	if this_.inverted {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	{
		hp_65 := this_.ranges
		h_ = 31 * h_
		if hp_65 != nil {
			h_ += hp_65.HashCode()
		}
	}
	{
		hs_66 := this_.namedGroups
		h_ = 31 * h_
		if hs_66 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_66)
		}
	}
	return h_
}
func (this_ *RegExpTree_DecomposedCharset) String() string {
	_ = fmt.Sprint
	return "DecomposedCharset[" + "inverted=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.inverted)) + ", " + "ranges=" + j2gort.StrVal(j2gort.PtrToAny(this_.ranges)) + ", " + "namedGroups=" + j2gort.StrVal(j2gort.BoxString(this_.namedGroups)) + "]"
}

// Java class: com.google.javascript.jscomp.regex.RegExpTree.Concatenation
type RegExpTree_Concatenation struct {
	RegExpTreeBase
	elements []RegExpTree `gson:"elements"`
}

func (this_ *RegExpTree_Concatenation) BindSelf_(s any) {
	_ = s
	this_.RegExpTreeBase.BindSelf_(s)
}
func (this_ *RegExpTree_Concatenation) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree_Concatenation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$Concatenation", "Concatenation"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.RegExpTree.Concatenation", "com.google.javascript.jscomp.regex.RegExpTree$Concatenation", func(v_ any) bool {
		return j2gort.InstanceOf[*RegExpTree_Concatenation](v_)
	})
}
func (*RegExpTree_Concatenation) Implements_regex_RegExpTree() {}

func (this_ *RegExpTree_Concatenation) initFields_RegExpTree_Concatenation() {
	_ = this_
}

// Java constructor: Concatenation(com.google.javascript.jscomp.regex.RegExpTree a, com.google.javascript.jscomp.regex.RegExpTree b)
func NewRegExpTree_Concatenation_2(a RegExpTree, b RegExpTree) *RegExpTree_Concatenation {
	this_ := &RegExpTree_Concatenation{}
	this_.BindSelf_(this_)
	InitRegExpTree_Concatenation_2(this_, a, b)
	return this_
}

func InitRegExpTree_Concatenation_2(this_ *RegExpTree_Concatenation, a RegExpTree, b RegExpTree) {
	_ = a
	_ = b
	InitRegExpTree(&this_.RegExpTreeBase)
	this_.initFields_RegExpTree_Concatenation()
	this_.elements = []RegExpTree{j2gort.CheckNotNullElem(a), j2gort.CheckNotNullElem(b)}
}

// Java constructor: Concatenation(java.util.List<? extends com.google.javascript.jscomp.regex.RegExpTree> elements)
func NewRegExpTree_Concatenation_1(elements j2gort.Ext_java_util_List) *RegExpTree_Concatenation {
	this_ := &RegExpTree_Concatenation{}
	this_.BindSelf_(this_)
	InitRegExpTree_Concatenation_1(this_, elements)
	return this_
}

func InitRegExpTree_Concatenation_1(this_ *RegExpTree_Concatenation, elements j2gort.Ext_java_util_List) {
	_ = elements
	InitRegExpTree(&this_.RegExpTreeBase)
	this_.initFields_RegExpTree_Concatenation()
	this_.elements = j2gort.ImmListCopyToSlice[RegExpTree](j2gort.AsIface[j2gort.Ext_java_util_Collection](elements))
}

// Java method: Concatenation.simplify(java.lang.String flags)
func (this_ *RegExpTree_Concatenation) Simplify(flags string) RegExpTree {
	_ = flags
	var s *RegExpTree__Local_125_Simplifier = NewRegExpTree__Local_125_Simplifier(this_, flags)
	_ = s
	{

		for _, element := range this_.elements {

			_ = element
			s.Simplify(element.Simplify(flags))
		}
	}
	switch s.simplified.Size() {
	case int32(0):
		return ToIface_RegExpTree(StaticGet_RegExpTree_Empty_INSTANCE())
	case int32(1):
		return j2gort.AsIface[RegExpTree](s.simplified.Get(int32(0)))
	default:
		return NewRegExpTree_Concatenation_1(s.simplified)
	}
	panic("unreachable switch expression")
}

// Java method: Concatenation.isCaseSensitive()
func (this_ *RegExpTree_Concatenation) IsCaseSensitive() bool {
	{

		for _, element := range this_.elements {

			_ = element
			if element.IsCaseSensitive() {
				return true
			}
		}
	}
	return false
}

// Java method: Concatenation.containsAnchor()
func (this_ *RegExpTree_Concatenation) ContainsAnchor() bool {
	{

		for _, element := range this_.elements {

			_ = element
			if element.ContainsAnchor() {
				return true
			}
		}
	}
	return false
}

// Java method: Concatenation.numCapturingGroups()
func (this_ *RegExpTree_Concatenation) NumCapturingGroups() int32 {
	var n int32 = int32(0)
	_ = n
	{

		for _, element := range this_.elements {

			_ = element
			n = (n + element.NumCapturingGroups())
		}
	}
	return n
}

// Java method: Concatenation.children()
func (this_ *RegExpTree_Concatenation) Children() any {
	return j2gort.ImmListFromTypedSlice(this_.elements)
}

// Java method: Concatenation.appendSourceCode(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Concatenation) AppendSourceCode(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
	var digitsMightBleed bool = false
	_ = digitsMightBleed
	{

		for _, element := range this_.elements {

			_ = element
			var parenthesize bool = false
			_ = parenthesize
			if j2gort.InstanceOf[*RegExpTree_Alternation](element) || j2gort.InstanceOf[*RegExpTree_Concatenation](element) {
				parenthesize = true
			}
			if parenthesize {
				sb.Append_1_string("(?:")
				element.AppendSourceCode(sb)
				sb.Append_1_uint16(uint16(41))
			} else {
				var start int32 = sb.Length()
				_ = start
				element.AppendSourceCode(sb)
				if digitsMightBleed && (sb.Length() > start) {
					var firstChar uint16 = sb.CharAt(start)
					_ = firstChar
					if (int32(uint16(48)) <= int32(firstChar)) && (int32(firstChar) <= int32(uint16(57))) {
						if int32(sb.CharAt((start - int32(1)))) == int32(uint16(123)) {
							sb.Insert_2_int32_uint16((start - int32(1)), uint16(92))
						} else {
							j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Insert_2_int32_string(start, "(?:")).Append_1_uint16(uint16(41))
						}
					}
				}
			}
			var backReference *RegExpTree_BackReference
			_ = backReference
			var text *RegExpTree_Text
			_ = text
			digitsMightBleed = ((j2gort.InstanceOfBind[*RegExpTree_BackReference](element, &backReference) && (backReference.groupIndex < int32(10))) || (j2gort.InstanceOfBind[*RegExpTree_Text](element, &text) && javalang.StringEndsWith(j2gort.CheckNotNullString(text.text), "{")))
		}
	}
}

// Java method: Concatenation.appendDebugInfo(java.lang.StringBuilder sb)
func (this_ *RegExpTree_Concatenation) AppendDebugInfo(sb *j2gort.Ext_java_lang_StringBuilder) {
	_ = sb
}

// Java method: Concatenation.equals(java.lang.Object o)
func (this_ *RegExpTree_Concatenation) Equals(arg_75 any) bool {
	o := arg_75
	_ = o
	var concatenation *RegExpTree_Concatenation
	_ = concatenation
	var ok_76 bool
	concatenation, ok_76 = j2gort.TryCast[*RegExpTree_Concatenation](o)
	return ok_76 && (j2gort.ImmSliceEquals(this_.elements, j2gort.ImmListFromTypedSlice(concatenation.elements)))
}

// Java method: Concatenation.hashCode()
func (this_ *RegExpTree_Concatenation) HashCode() int32 {
	return (int32(546930238) ^ j2gort.ImmSliceHashCode(this_.elements))
}

// Java class: Parser
type RegExpTree__Local_124_Parser struct {
	cap_0_pattern                     string
	cap_1_flags                       string
	capturingGroupNames               j2gort.OrderedSet[string] `gson:"capturingGroupNames"`
	pos                               int32                     `gson:"pos"`
	numCapturingGroups                int32                     `gson:"numCapturingGroups"`
	limit                             int32                     `gson:"limit"`
	idHash_                           atomic.Int32
	lookForNamedCaptureBackreferences bool `gson:"lookForNamedCaptureBackreferences"`
}

func (this_ *RegExpTree__Local_124_Parser) BindSelf_(s any) {
	_ = s
}
func (this_ *RegExpTree__Local_124_Parser) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree__Local_124_Parser) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$1Parser", "Parser"
}
func (this_ *RegExpTree__Local_124_Parser) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RegExpTree__Local_124_Parser) HasIdentityEquals_() {}

func (this_ *RegExpTree__Local_124_Parser) initFields_RegExpTree__Local_124_Parser() {
	_ = this_
	this_.capturingGroupNames = j2gort.NewOrderedSet[string]()
	this_.limit = j2gort.JStrLen(this_.cap_0_pattern)
}

// Java constructor: Parser()
func NewRegExpTree__Local_124_Parser(cap_p_pattern string, cap_p_flags string) *RegExpTree__Local_124_Parser {
	this_ := &RegExpTree__Local_124_Parser{}
	this_.cap_0_pattern = cap_p_pattern
	this_.cap_1_flags = cap_p_flags
	InitRegExpTree__Local_124_Parser(this_)
	return this_
}

func InitRegExpTree__Local_124_Parser(this_ *RegExpTree__Local_124_Parser) {
	this_.initFields_RegExpTree__Local_124_Parser()
}

// Java method: Parser.parseTopLevel()
func (this_ *RegExpTree__Local_124_Parser) ParseTopLevel() RegExpTree {
	this_.pos = int32(0)
	this_.numCapturingGroups = int32(0)
	this_.lookForNamedCaptureBackreferences = false
	var out RegExpTree = this_.Parse()
	_ = out
	if !(this_.capturingGroupNames.IsEmpty()) {
		this_.pos = int32(0)
		this_.numCapturingGroups = int32(0)
		this_.lookForNamedCaptureBackreferences = true
		out = this_.Parse()
	}
	if this_.pos < this_.limit {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(j2gort.JStrSubFrom(this_.cap_0_pattern, this_.pos)))))
	}
	return out
}

// Java method: Parser.parse()
func (this_ *RegExpTree__Local_124_Parser) Parse() RegExpTree {
	var alternatives *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = nil
	_ = alternatives
	var preceder RegExpTree = nil
	_ = preceder
topLoop:
	for this_.pos < this_.limit {
		var ch uint16 = j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)
		_ = ch
		var atom RegExpTree = nil
		_ = atom
		switch ch {
		case uint16(91):
			atom = this_.priv_Parser_parseCharset_7666()
		case uint16(40):
			atom = this_.priv_Parser_parseParenthetical_7664()
		case uint16(41):
			break topLoop
		case uint16(92):
			atom = this_.priv_Parser_parseEscape_7668()
		case uint16(94), uint16(36):
			atom = NewRegExpTree_Anchor(ch)
			this_.pos = this_.pos + 1
		case uint16(46):
			atom = ToIface_RegExpTree(StaticGet_RegExpTree_DOT_CHARSET())
			this_.pos = this_.pos + 1
		case uint16(124):
			atom = ToIface_RegExpTree(StaticGet_RegExpTree_Empty_INSTANCE())
		default:
			var start int32 = this_.pos
			_ = start
			var end int32 = (this_.pos + int32(1))
			_ = end
		charsLoop:
			for end < this_.limit {
				switch j2gort.JStrCharAt(this_.cap_0_pattern, end) {
				case uint16(91), uint16(40), uint16(41), uint16(92), uint16(94), uint16(36), uint16(124), uint16(46), uint16(42), uint16(43), uint16(63), uint16(123):
					break charsLoop
				default:
					if ((end + int32(1)) >= this_.limit) || !(this_.priv_Parser_isRepetitionStart_7671(j2gort.JStrCharAt(this_.cap_0_pattern, (end + int32(1))))) {
						end = end + 1
					} else {
						break charsLoop
					}
				}
			}
			atom = NewRegExpTree_Text(j2gort.JStrSub(this_.cap_0_pattern, start, end))
			this_.pos = end
		}
		if (this_.pos < this_.limit) && this_.priv_Parser_isRepetitionStart_7671(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) {
			atom = this_.priv_Parser_parseRepetition_7672(atom)
		}
		if preceder == nil || preceder.IsJavaNil_() {
			preceder = atom
		} else {
			preceder = NewRegExpTree_Concatenation_2(preceder, atom)
		}
		if (this_.pos < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(124))) {
			if alternatives == nil {
				alternatives = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
			}
			alternatives.AddVoid_(preceder)
			preceder = nil
			this_.pos = this_.pos + 1
		}
	}
	if preceder == nil || preceder.IsJavaNil_() {
		preceder = ToIface_RegExpTree(StaticGet_RegExpTree_Empty_INSTANCE())
	}
	if alternatives != nil {
		alternatives.AddVoid_(preceder)
		return NewRegExpTree_Alternation(j2gort.AsIface[j2gort.Ext_java_util_List](alternatives.BuildImmList_()))
	} else {
		return preceder
	}
}

// Java method: Parser.parseParenthetical()
func (this_ *RegExpTree__Local_124_Parser) priv_Parser_parseParenthetical_7664() RegExpTree {
	if !(int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(40))) {
		j2gort.CheckState(false)
	}
	var start int32 = this_.pos
	_ = start
	this_.pos = this_.pos + 1
	var type_ RegExpTree_ParentheticalType = RegExpTree_ParentheticalType(0)
	_ = type_
	var captureName string = j2gort.NullString
	_ = captureName
	if (this_.pos < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(63))) {
		if (this_.pos + int32(1)) < this_.limit {
			var ch uint16 = j2gort.JStrCharAt(this_.cap_0_pattern, (this_.pos + int32(1)))
			_ = ch
			switch ch {
			case uint16(58):
				this_.pos = (this_.pos + int32(2))
				type_ = StaticGet_RegExpTree_ParentheticalType_NONCAPTURING()
			case uint16(61):
				this_.pos = (this_.pos + int32(2))
				type_ = StaticGet_RegExpTree_ParentheticalType_POSITIVE_LOOKAHEAD()
			case uint16(33):
				this_.pos = (this_.pos + int32(2))
				type_ = StaticGet_RegExpTree_ParentheticalType_NEGATIVE_LOOKAHEAD()
			case uint16(60):
				if ((this_.pos + int32(2)) < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, (this_.pos+int32(2)))) == int32(uint16(61))) {
					this_.pos = (this_.pos + int32(3))
					type_ = StaticGet_RegExpTree_ParentheticalType_POSITIVE_LOOKBEHIND()
				} else if ((this_.pos + int32(2)) < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, (this_.pos+int32(2)))) == int32(uint16(33))) {
					this_.pos = (this_.pos + int32(3))
					type_ = StaticGet_RegExpTree_ParentheticalType_NEGATIVE_LOOKBEHIND()
				} else {
					this_.pos = (this_.pos + int32(2))
					captureName = this_.priv_Parser_scanNamedGroupName_7665()
					this_.capturingGroupNames.Add(captureName)
					type_ = StaticGet_RegExpTree_ParentheticalType_NAMED_GROUPS()
				}
			default:
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Malformed parenthetical: "+j2gort.StrValString(j2gort.JStrSubFrom(this_.cap_0_pattern, start)))))
			}
		} else {
			panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Malformed parenthetical: "+j2gort.StrValString(j2gort.JStrSubFrom(this_.cap_0_pattern, start)))))
		}
	} else {
		type_ = StaticGet_RegExpTree_ParentheticalType_CAPTURING()
	}
	var body RegExpTree = this_.Parse()
	_ = body
	if (this_.pos < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(41))) {
		this_.pos = this_.pos + 1
	} else {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unclosed parenthetical group: "+j2gort.StrValString(j2gort.JStrSubFrom(this_.cap_0_pattern, start)))))
	}
	switch (type_).Ordinal() {
	case int32(0) /* ParentheticalType.CAPTURING */ :
		this_.numCapturingGroups = this_.numCapturingGroups + 1
		return NewRegExpTree_CapturingGroup(body)
	case int32(1) /* ParentheticalType.NONCAPTURING */ :
		return body
	case int32(2) /* ParentheticalType.POSITIVE_LOOKAHEAD */ :
		return NewRegExpTree_LookaheadAssertion(body, true)
	case int32(3) /* ParentheticalType.NEGATIVE_LOOKAHEAD */ :
		return NewRegExpTree_LookaheadAssertion(body, false)
	case int32(4) /* ParentheticalType.POSITIVE_LOOKBEHIND */ :
		return NewRegExpTree_LookbehindAssertion(body, true)
	case int32(5) /* ParentheticalType.NEGATIVE_LOOKBEHIND */ :
		return NewRegExpTree_LookbehindAssertion(body, false)
	case int32(6) /* ParentheticalType.NAMED_GROUPS */ :
		if captureName != j2gort.NullString {
			this_.numCapturingGroups = this_.numCapturingGroups + 1
			return NewRegExpTree_NamedCaptureGroup(body, captureName)
		} else {
			panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Malformed named capture group: "+j2gort.StrValString(j2gort.JStrSubFrom(this_.cap_0_pattern, start)))))
		}
	}
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unrecognized ParentheticalType "+j2gort.StrVal(type_))))
}

// Java method: Parser.scanNamedGroupName()
func (this_ *RegExpTree__Local_124_Parser) priv_Parser_scanNamedGroupName_7665() string {
	var start int32 = this_.pos
	_ = start
	var end int32 = this_.pos
	_ = end
	if !(RegExpTree_IsIdentifierStart(j2gort.JStrCharAt(this_.cap_0_pattern, start))) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Invalid capture group name: <"+j2gort.StrValString(j2gort.JStrSubFrom(this_.cap_0_pattern, start)))))
	}
	end = end + 1
	for end < this_.limit {
		if int32(j2gort.JStrCharAt(this_.cap_0_pattern, end)) == int32(uint16(62)) {
			this_.pos = (end + int32(1))
			return j2gort.JStrSub(this_.cap_0_pattern, start, end)
		} else if RegExpTree_IsIdentifierPart(j2gort.JStrCharAt(this_.cap_0_pattern, end)) {
			end = end + 1
		} else {
			panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Invalid capture group name: <"+j2gort.StrValString(j2gort.JStrSubFrom(this_.cap_0_pattern, start)))))
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Malformed named capture group: <"+j2gort.StrValString(j2gort.JStrSubFrom(this_.cap_0_pattern, start)))))
}

// Java method: Parser.parseCharset()
func (this_ *RegExpTree__Local_124_Parser) priv_Parser_parseCharset_7666() RegExpTree {
	if !(int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(91))) {
		j2gort.CheckState(false)
	}
	this_.pos = this_.pos + 1
	var isCaseInsensitive bool = (j2gort.JStrIndexOfCp(this_.cap_1_flags, int32(uint16(105))) >= int32(0))
	_ = isCaseInsensitive
	var inverse bool = ((this_.pos < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(94))))
	_ = inverse
	if inverse {
		this_.pos = this_.pos + 1
	}
	var ranges *CharRanges = StaticGet_CharRanges_EMPTY()
	_ = ranges
	var ieExplicits *CharRanges = StaticGet_CharRanges_EMPTY()
	_ = ieExplicits
	for (this_.pos < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) != int32(uint16(93))) {
		var ch uint16 = j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)
		_ = ch
		var start int32 = int32(0)
		_ = start
		if int32(ch) == int32(uint16(92)) {
			this_.pos = this_.pos + 1
			var possibleGroupName uint16 = j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)
			_ = possibleGroupName
			var group *CharRanges = j2gort.CastPtr[CharRanges](StaticGet_RegExpTree_NAMED_CHAR_GROUPS().Get(j2gort.ValueOfChar(possibleGroupName)))
			_ = group
			if group != nil {
				this_.pos = this_.pos + 1
				ranges = ranges.Union(group)
				continue
			}
			start = this_.priv_Parser_parseEscapeChar_7667()
		} else {
			start = int32(ch)
			this_.pos = this_.pos + 1
		}
		var end int32 = start
		_ = end
		if (((this_.pos + int32(1)) < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(45)))) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, (this_.pos+int32(1)))) != int32(uint16(93))) {
			this_.pos = this_.pos + 1
			ch = j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)
			if int32(ch) == int32(uint16(92)) {
				this_.pos = this_.pos + 1
				end = this_.priv_Parser_parseEscapeChar_7667()
			} else {
				end = int32(ch)
				this_.pos = this_.pos + 1
			}
		}
		var range_ *CharRanges = CharRanges_Inclusive(start, end)
		_ = range_
		ranges = ranges.Union(range_)
		if StaticGet_RegExpTree_IE_SPEC_ERRORS().Contains(start) && StaticGet_RegExpTree_IE_SPEC_ERRORS().Contains(end) {
			ieExplicits = ieExplicits.Union(range_.Intersection(StaticGet_RegExpTree_IE_SPEC_ERRORS()))
		}
		if isCaseInsensitive {
			ranges = CaseCanonicalize_ExpandToAllMatched(ranges)
		}
	}
	this_.pos = this_.pos + 1
	if inverse {
		ranges = StaticGet_CharRanges_ALL_CODE_UNITS().Difference(ranges)
	}
	return NewRegExpTree_Charset(ranges, ieExplicits)
}

// Java method: Parser.parseEscapeChar()
func (this_ *RegExpTree__Local_124_Parser) priv_Parser_parseEscapeChar_7667() int32 {
	var ch uint16 = j2gort.JStrCharAt(this_.cap_0_pattern, func() int32 { old_77 := this_.pos; nw_78 := old_77 + 1; this_.pos = nw_78; return old_77 }())
	_ = ch
	switch ch {
	case uint16(98):
		return int32(uint16(8))
	case uint16(102):
		return int32(uint16(12))
	case uint16(110):
		return int32(uint16(10))
	case uint16(114):
		return int32(uint16(13))
	case uint16(116):
		return int32(uint16(9))
	case uint16(117):
		return func() int32 {
			if ((j2gort.JStrIndexOf(this_.cap_1_flags, "u") >= 0) && (this_.pos < this_.limit)) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(123))) {
				return this_.priv_Parser_parseBracedUnicodeEscape_7670()
			}
			return this_.priv_Parser_parseHex_7669(int32(4))
		}()
	case uint16(118):
		return int32(uint16(11))
	case uint16(120):
		return this_.priv_Parser_parseHex_7669(int32(2))
	default:
		if (int32(uint16(48)) <= int32(ch)) && (int32(ch) <= int32(uint16(55))) {
			var codeUnit uint16 = uint16((int32(ch) - int32(uint16(48))))
			_ = codeUnit
			var octLimit int32 = min(this_.limit, ((this_.pos + (func() int32 {
				if int32(ch) <= int32(uint16(51)) {
					return int32(2)
				}
				return int32(1)
			}())) + (func() int32 {
				if int32(ch) == int32(uint16(48)) {
					return int32(1)
				}
				return int32(0)
			}())))
			_ = octLimit
			for this_.pos < octLimit {
				ch = j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)
				if (int32(uint16(48)) <= int32(ch)) && (int32(ch) <= int32(uint16(55))) {
					codeUnit = uint16(((javalang.Shl32(int32(codeUnit), int(int32(3)))) + (int32(ch) - int32(uint16(48)))))
					this_.pos = this_.pos + 1
				} else {
					break
				}
			}
			return int32(codeUnit)
		}
		return int32(ch)
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseEscape()
func (this_ *RegExpTree__Local_124_Parser) priv_Parser_parseEscape_7668() RegExpTree {
	if !(int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(92))) {
		j2gort.CheckState(false)
	}
	var start int32 = this_.pos
	_ = start
	this_.pos = this_.pos + 1
	var ch uint16 = j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)
	_ = ch
	if (int32(ch) == int32(uint16(98))) || (int32(ch) == int32(uint16(66))) {
		this_.pos = this_.pos + 1
		return NewRegExpTree_WordBoundary(ch)
	} else if ((int32(ch) == int32(uint16(112))) || (int32(ch) == int32(uint16(80)))) && (j2gort.JStrIndexOf(this_.cap_1_flags, "u") >= 0) {
		var negated bool = (int32(ch) == int32(uint16(80)))
		_ = negated
		this_.pos = this_.pos + 1
		if (this_.pos < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(123))) {
			var lhs j2gort.LocalStringBuilder

			for (func() int32 { old_79 := this_.pos; nw_80 := old_79 + 1; this_.pos = nw_80; return nw_80 }() < this_.limit) && ((((int32((func() uint16 {
				var asgn_81 uint16 = j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)
				ch = asgn_81
				return asgn_81
			}())) == int32(uint16(95))) || ((int32(uint16(97)) <= int32(ch)) && (int32(ch) <= int32(uint16(122))))) || ((int32(uint16(65)) <= int32(ch)) && (int32(ch) <= int32(uint16(90))))) || ((int32(uint16(48)) <= int32(ch)) && (int32(ch) <= int32(uint16(57))))) {
				lhs.AppendChar(ch)
			}
			if (this_.pos < this_.limit) && (int32(ch) == int32(uint16(125))) {
				this_.pos = this_.pos + 1
				return NewRegExpTree_UnicodePropertyEscape(j2gort.NullString, lhs.String(), negated)
			} else if (this_.pos < this_.limit) && (int32(ch) == int32(uint16(61))) {
				var rhs j2gort.LocalStringBuilder

				for (func() int32 { old_82 := this_.pos; nw_83 := old_82 + 1; this_.pos = nw_83; return nw_83 }() < this_.limit) && ((((int32((func() uint16 {
					var asgn_84 uint16 = j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)
					ch = asgn_84
					return asgn_84
				}())) == int32(uint16(95))) || ((int32(uint16(97)) <= int32(ch)) && (int32(ch) <= int32(uint16(122))))) || ((int32(uint16(65)) <= int32(ch)) && (int32(ch) <= int32(uint16(90))))) || ((int32(uint16(48)) <= int32(ch)) && (int32(ch) <= int32(uint16(57))))) {
					rhs.AppendChar(ch)
				}
				if (this_.pos < this_.limit) && (int32(ch) == int32(uint16(125))) {
					this_.pos = this_.pos + 1
					return NewRegExpTree_UnicodePropertyEscape(lhs.String(), rhs.String(), negated)
				} else {
					panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Malformed Unicode Property Escape: expected '}' after "+j2gort.StrValString(j2gort.JStrSub(this_.cap_0_pattern, start, this_.pos)))))
				}
			} else {
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Malformed Unicode Property Escape: expected '=' or '}' after "+j2gort.StrValString(j2gort.JStrSub(this_.cap_0_pattern, start, this_.pos)))))
			}
		} else {
			panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Malformed Unicode Property Escape: expected '{' after "+j2gort.StrValString(j2gort.JStrSub(this_.cap_0_pattern, start, this_.pos)))))
		}
	} else if (int32(uint16(49)) <= int32(ch)) && (int32(ch) <= int32(uint16(57))) {
		this_.pos = this_.pos + 1
		var possibleGroupIndex int32 = (int32(ch) - int32(uint16(48)))
		_ = possibleGroupIndex
		if this_.numCapturingGroups >= possibleGroupIndex {
			if this_.pos < this_.limit {
				var next uint16 = j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)
				_ = next
				if (int32(uint16(48)) <= int32(next)) && (int32(next) <= int32(uint16(57))) {
					var twoDigitGroupIndex int32 = ((possibleGroupIndex * int32(10)) + (int32(next) - int32(uint16(48))))
					_ = twoDigitGroupIndex
					if this_.numCapturingGroups >= twoDigitGroupIndex {
						this_.pos = this_.pos + 1
						possibleGroupIndex = twoDigitGroupIndex
					}
				}
			}
			return NewRegExpTree_BackReference(possibleGroupIndex)
		} else {
			return NewRegExpTree_Text(javalang.StringFromChar(func() uint16 {
				if possibleGroupIndex <= int32(7) {
					return uint16(possibleGroupIndex)
				}
				return ch
			}()))
		}
	} else if (((this_.lookForNamedCaptureBackreferences && (int32(ch) == int32(uint16(107)))) && ((this_.pos + int32(1)) < this_.limit)) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, (this_.pos+int32(1)))) == int32(uint16(60)))) && !(this_.capturingGroupNames.IsEmpty()) {
		this_.pos = (this_.pos + int32(2))
		var potentialName string = this_.priv_Parser_scanNamedGroupName_7665()
		_ = potentialName
		if !(this_.capturingGroupNames.Contains(potentialName)) {
			panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Invalid named capture referenced: "+j2gort.StrValString(j2gort.JStrSubFrom(this_.cap_0_pattern, start)))))
		}
		return NewRegExpTree_NamedBackReference(potentialName)
	} else {
		var charGroup *CharRanges = j2gort.CastPtr[CharRanges](StaticGet_RegExpTree_NAMED_CHAR_GROUPS().Get(j2gort.ValueOfChar(ch)))
		_ = charGroup
		if charGroup != nil {
			this_.pos = this_.pos + 1
			return NewRegExpTree_Charset(charGroup, StaticGet_CharRanges_EMPTY())
		}
		return NewRegExpTree_Text(javalang.StringFromUTF16(j2gort.BoxChar_ToChars_1(this_.priv_Parser_parseEscapeChar_7667())))
	}
}

// Java method: Parser.parseHex(int n)
func (this_ *RegExpTree__Local_124_Parser) priv_Parser_parseHex_7669(n int32) int32 {
	_ = n
	if (this_.pos + n) > this_.limit {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Abbreviated hex escape "+j2gort.StrValString(j2gort.JStrSubFrom(this_.cap_0_pattern, this_.pos)))))
	}
	if n > int32(7) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Cannot parse hexadecimal encoding wider than 28 bits: "+j2gort.StrValString(j2gort.JStrSub(this_.cap_0_pattern, this_.pos, (this_.pos+n))))))
	}
	var result int32 = int32(0)
	_ = result
	for func() int32 { old_85 := n; nw_86 := old_85 - 1; n = nw_86; return nw_86 }() >= int32(0) {
		var ch uint16 = j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)
		_ = ch
		var digit int32 = int32(0)
		_ = digit
		if (int32(uint16(48)) <= int32(ch)) && (int32(ch) <= int32(uint16(57))) {
			digit = (int32(ch) - int32(uint16(48)))
		} else if (int32(uint16(97)) <= int32(ch)) && (int32(ch) <= int32(uint16(102))) {
			digit = (int32(ch) + (int32(10) - int32(uint16(97))))
		} else if (int32(uint16(65)) <= int32(ch)) && (int32(ch) <= int32(uint16(70))) {
			digit = (int32(ch) + (int32(10) - int32(uint16(65))))
		} else {
			panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(j2gort.JStrSubFrom(this_.cap_0_pattern, this_.pos)))))
		}
		this_.pos = this_.pos + 1
		result = ((javalang.Shl32(result, int(int32(4)))) | digit)
	}
	return result
}

// Java method: Parser.parseBracedUnicodeEscape()
func (this_ *RegExpTree__Local_124_Parser) priv_Parser_parseBracedUnicodeEscape_7670() int32 {
	var openBrace int32 = this_.pos
	_ = openBrace
	if !(int32(j2gort.JStrCharAt(this_.cap_0_pattern, func() int32 { old_87 := this_.pos; nw_88 := old_87 + 1; this_.pos = nw_88; return old_87 }())) == int32(uint16(123))) {
		j2gort.CheckState(false)
	}
	var closeBrace int32 = this_.pos
	_ = closeBrace
	for (closeBrace < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, closeBrace)) != int32(uint16(125))) {
		closeBrace = closeBrace + 1
	}
	if closeBrace == this_.limit {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Malformed unicode escape: expected '}' after "+j2gort.StrValString(j2gort.JStrSubFrom(this_.cap_0_pattern, openBrace)))))
	} else if closeBrace == this_.pos {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Empty unicode escape")))
	}
	var result int32 = this_.priv_Parser_parseHex_7669((closeBrace - this_.pos))
	_ = result
	if result > int32(1114111) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unicode must be at most 0x10FFFF: "+j2gort.StrValString(j2gort.JStrSub(this_.cap_0_pattern, (openBrace+int32(1)), this_.pos)))))
	}
	this_.pos = this_.pos + 1
	return result
}

// Java method: Parser.isRepetitionStart(char ch)
func (this_ *RegExpTree__Local_124_Parser) priv_Parser_isRepetitionStart_7671(ch uint16) bool {
	_ = ch
	switch ch {
	case uint16(63), uint16(42), uint16(43), uint16(123):
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Parser.parseRepetition(com.google.javascript.jscomp.regex.RegExpTree body)
func (this_ *RegExpTree__Local_124_Parser) priv_Parser_parseRepetition_7672(body RegExpTree) RegExpTree {
	_ = body
	if this_.pos == this_.limit {
		return body
	}
	var min_ int32 = int32(0)
	_ = min_
	var max_ int32 = int32(0)
	_ = max_
	switch j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos) {
	case uint16(43):
		this_.pos = this_.pos + 1
		min_ = int32(1)
		max_ = int32(2147483647)
	case uint16(42):
		this_.pos = this_.pos + 1
		min_ = int32(0)
		max_ = int32(2147483647)
	case uint16(63):
		this_.pos = this_.pos + 1
		min_ = int32(0)
		max_ = int32(1)
	case uint16(123):
		this_.pos = this_.pos + 1
		var start int32 = this_.pos
		_ = start
		var end int32 = j2gort.JStrIndexOfCpFrom(this_.cap_0_pattern, int32(uint16(125)), start)
		_ = end
		if end < int32(0) {
			this_.pos = (start - int32(1))
			return body
		}
		var counts string = j2gort.JStrSub(this_.cap_0_pattern, start, end)
		_ = counts
		this_.pos = (end + int32(1))
		var comma int32 = j2gort.JStrIndexOfCp(counts, int32(uint16(44)))
		_ = comma
		func() {
			defer func() {
				if rec_89 := recover(); rec_89 != nil {
					jex_90 := javalang.RecoverException(rec_89)
					_ = jex_90
					if j2gort.InstanceOf[javalang.INumberFormatException](jex_90) {
						ex := j2gort.AsIface[javalang.INumberFormatException](jex_90)
						_ = ex
						min_ = func() int32 { var asgn_91 int32 = int32(-1); max_ = asgn_91; return asgn_91 }()
					} else {
						panic(rec_89)
					}
				}
			}()
			min_ = j2gort.ParseInt(func() string {
				if comma >= int32(0) {
					return j2gort.JStrSub(counts, int32(0), comma)
				}
				return counts
			}(), 10)
			if comma >= int32(0) {
				max_ = func() int32 {
					if (comma + int32(1)) != j2gort.JStrLen(counts) {
						return j2gort.ParseInt(j2gort.JStrSubFrom(counts, (comma+int32(1))), 10)
					}
					return int32(2147483647)
				}()
			} else {
				max_ = min_
			}
		}()
		if (min_ < int32(0)) || (min_ > max_) {
			this_.pos = (start - int32(1))
			return body
		}
	default:
		return body
	}
	var greedy bool = true
	_ = greedy
	if (this_.pos < this_.limit) && (int32(j2gort.JStrCharAt(this_.cap_0_pattern, this_.pos)) == int32(uint16(63))) {
		greedy = false
		this_.pos = this_.pos + 1
	}
	return NewRegExpTree_Repetition(body, min_, max_, greedy)
}

// Java class: Simplifier
type RegExpTree__Local_125_Simplifier struct {
	outer_      *RegExpTree_Concatenation
	cap_0_flags string
	simplified  j2gort.Ext_java_util_List `gson:"simplified"`
	idHash_     atomic.Int32
}

func (this_ *RegExpTree__Local_125_Simplifier) BindSelf_(s any) {
	_ = s
}
func (this_ *RegExpTree__Local_125_Simplifier) IsJavaNil_() bool { return this_ == nil }
func (*RegExpTree__Local_125_Simplifier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.RegExpTree$Concatenation$1Simplifier", "Simplifier"
}
func (this_ *RegExpTree__Local_125_Simplifier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RegExpTree__Local_125_Simplifier) HasIdentityEquals_() {}
func (this_ *RegExpTree__Local_125_Simplifier) SetOuter_(o *RegExpTree_Concatenation) {
	this_.outer_ = o
}
func (this_ *RegExpTree__Local_125_Simplifier) GetOuter_() *RegExpTree_Concatenation {
	return this_.outer_
}

func (this_ *RegExpTree__Local_125_Simplifier) initFields_RegExpTree__Local_125_Simplifier() {
	_ = this_
	this_.simplified = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: Simplifier()
func NewRegExpTree__Local_125_Simplifier(outer_ *RegExpTree_Concatenation, cap_p_flags string) *RegExpTree__Local_125_Simplifier {
	this_ := &RegExpTree__Local_125_Simplifier{}
	this_.outer_ = outer_
	this_.cap_0_flags = cap_p_flags
	InitRegExpTree__Local_125_Simplifier(this_)
	return this_
}

func InitRegExpTree__Local_125_Simplifier(this_ *RegExpTree__Local_125_Simplifier) {
	this_.initFields_RegExpTree__Local_125_Simplifier()
}

// Java method: Simplifier.simplify(com.google.javascript.jscomp.regex.RegExpTree t)
func (this_ *RegExpTree__Local_125_Simplifier) Simplify(t RegExpTree) {
	_ = t
	var concatenation *RegExpTree_Concatenation
	_ = concatenation
	var ok_92 bool
	concatenation, ok_92 = j2gort.TryCast[*RegExpTree_Concatenation](t)
	if ok_92 {
		{

			for _, child := range concatenation.elements {

				_ = child
				this_.Simplify(child)
			}
		}
	} else if j2gort.InstanceOf[*RegExpTree_Empty](t) {
	} else {
		var lastIndex int32 = (this_.simplified.Size() - int32(1))
		_ = lastIndex
		if lastIndex >= int32(0) {
			var pairwise RegExpTree = this_.SimplifyPairwise(j2gort.AsIface[RegExpTree](this_.simplified.Get(lastIndex)), t)
			_ = pairwise
			if pairwise != nil && !pairwise.IsJavaNil_() {
				this_.simplified.Set(lastIndex, pairwise)
				return
			}
		}
		this_.simplified.Add_1(t)
	}
}

// Java method: Simplifier.simplifyPairwise(com.google.javascript.jscomp.regex.RegExpTree before, com.google.javascript.jscomp.regex.RegExpTree after)
func (this_ *RegExpTree__Local_125_Simplifier) SimplifyPairwise(before RegExpTree, after RegExpTree) RegExpTree {
	_ = before
	_ = after
	var beforeText *RegExpTree_Text
	_ = beforeText
	var afterText *RegExpTree_Text
	_ = afterText
	var ok_95 bool
	beforeText, ok_95 = j2gort.TryCast[*RegExpTree_Text](before)
	if ok_95 {
		afterText, ok_95 = j2gort.TryCast[*RegExpTree_Text](after)
	}
	if ok_95 {
		return NewRegExpTree_Text((j2gort.StrValString(beforeText.text) + j2gort.StrValString(afterText.text))).Simplify(this_.cap_0_flags)
	}
	var beforeMin int32 = int32(1)
	_ = beforeMin
	var beforeMax int32 = int32(1)
	_ = beforeMax
	var beforeBody RegExpTree = before
	_ = beforeBody
	var beforeGreedy bool = false
	_ = beforeGreedy
	var r *RegExpTree_Repetition
	_ = r
	var ok_96 bool
	r, ok_96 = j2gort.TryCast[*RegExpTree_Repetition](before)
	if ok_96 {
		beforeMin = r.min_
		beforeMax = r.max_
		beforeBody = r.body
		beforeGreedy = r.greedy
	}
	var afterMin int32 = int32(1)
	_ = afterMin
	var afterMax int32 = int32(1)
	_ = afterMax
	var afterBody RegExpTree = after
	_ = afterBody
	var afterGreedy bool = false
	_ = afterGreedy
	var r_2 *RegExpTree_Repetition
	_ = r_2
	var ok_97 bool
	r_2, ok_97 = j2gort.TryCast[*RegExpTree_Repetition](after)
	if ok_97 {
		afterMin = r_2.min_
		afterMax = r_2.max_
		afterBody = r_2.body
		afterGreedy = r_2.greedy
	}
	if beforeBody.Equals(afterBody) && !(beforeBody.GetRegExpTreeBase().HasCapturingGroup()) {
		var lmin int64 = ((int64(beforeMin)) + int64(afterMin))
		_ = lmin
		var lmax int64 = ((int64(beforeMax)) + int64(afterMax))
		_ = lmax
		if lmin < int64(int32(2147483647)) {
			var min_ int32 = int32(lmin)
			_ = min_
			var max_ int32
			_ = max_
			if lmax >= int64(int32(2147483647)) {
				max_ = int32(2147483647)
			} else {
				max_ = int32(lmax)
			}
			return NewRegExpTree_Repetition(beforeBody, min_, max_, ((beforeGreedy || afterGreedy) || (min_ == max_)))
		}
	}
	return nil
}
func (this_ *RegExpTree_Alternation) Slice_Children() []RegExpTree {
	return this_.alternatives
}
func (this_ *RegExpTree_Concatenation) Slice_Children() []RegExpTree {
	return this_.elements
}

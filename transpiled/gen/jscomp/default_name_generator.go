/*
 * Copyright 2005 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/DefaultNameGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.DefaultNameGenerator
type DefaultNameGenerator struct {
	priorityLookupMap j2gort.Ext_java_util_Map             `gson:"priorityLookupMap"`
	reservedNames     j2gort.Ext_java_util_Set             `gson:"reservedNames"`
	prefix            string                               `gson:"prefix"`
	firstChars        []*DefaultNameGenerator_CharPriority `gson:"firstChars"`
	nonFirstChars     []*DefaultNameGenerator_CharPriority `gson:"nonFirstChars"`
	nameCount         int32                                `gson:"nameCount"`
	idHash_           atomic.Int32
}

func (this_ *DefaultNameGenerator) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultNameGenerator) IsJavaNil_() bool { return this_ == nil }
func (*DefaultNameGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultNameGenerator", "DefaultNameGenerator"
}
func (this_ *DefaultNameGenerator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultNameGenerator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DefaultNameGenerator", "com.google.javascript.jscomp.DefaultNameGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[*DefaultNameGenerator](v_)
	})
}
func (*DefaultNameGenerator) Implements_jscomp_NameGenerator() {}

var static_DefaultNameGenerator_FIRST_CHAR []uint16

func StaticGet_DefaultNameGenerator_FIRST_CHAR() []uint16 {
	clinit_DefaultNameGenerator()
	return static_DefaultNameGenerator_FIRST_CHAR
}
func StaticSet_DefaultNameGenerator_FIRST_CHAR(v []uint16) {
	clinit_DefaultNameGenerator()
	static_DefaultNameGenerator_FIRST_CHAR = v
}

var static_DefaultNameGenerator_NONFIRST_CHAR []uint16

func StaticGet_DefaultNameGenerator_NONFIRST_CHAR() []uint16 {
	clinit_DefaultNameGenerator()
	return static_DefaultNameGenerator_NONFIRST_CHAR
}
func StaticSet_DefaultNameGenerator_NONFIRST_CHAR(v []uint16) {
	clinit_DefaultNameGenerator()
	static_DefaultNameGenerator_NONFIRST_CHAR = v
}

var static_DefaultNameGenerator_BAD_NAMES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_DefaultNameGenerator_BAD_NAMES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_DefaultNameGenerator()
	return static_DefaultNameGenerator_BAD_NAMES
}
func StaticSet_DefaultNameGenerator_BAD_NAMES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_DefaultNameGenerator()
	static_DefaultNameGenerator_BAD_NAMES = v
}

var clinitDone_DefaultNameGenerator uint32
var clinitState_DefaultNameGenerator j2gort.ClinitState

func clinit_DefaultNameGenerator() {
	if atomic.LoadUint32(&clinitDone_DefaultNameGenerator) != 0 {
		return
	}
	clinitSlow_DefaultNameGenerator()
}

func clinitSlow_DefaultNameGenerator() {
	if !j2gort.EnterClinit(&clinitState_DefaultNameGenerator, &clinitDone_DefaultNameGenerator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DefaultNameGenerator, &clinitDone_DefaultNameGenerator)
	static_DefaultNameGenerator_FIRST_CHAR = javalang.StringToUTF16("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ$")
	static_DefaultNameGenerator_NONFIRST_CHAR = javalang.StringToUTF16("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789$")
	static_DefaultNameGenerator_BAD_NAMES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("let", "yield", "await", "eval")
}

func (this_ *DefaultNameGenerator) initFields_DefaultNameGenerator() {
	_ = this_
	this_.prefix = j2gort.NullString
}

// Java constructor: DefaultNameGenerator()
func NewDefaultNameGenerator() *DefaultNameGenerator {
	clinit_DefaultNameGenerator()
	this_ := &DefaultNameGenerator{}
	InitDefaultNameGenerator(this_)
	return this_
}

func InitDefaultNameGenerator(this_ *DefaultNameGenerator) {
	this_.initFields_DefaultNameGenerator()
	this_.priv_DefaultNameGenerator_buildPriorityLookupMap_2355()
	var reservedNames j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Sets_NewLinkedHashSetWithExpectedSize(int32(0)))
	_ = reservedNames
	this_.Reset_3(reservedNames, "", j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()))
}

// Java constructor: DefaultNameGenerator(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedCharacters)
func NewDefaultNameGenerator_3(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedCharacters j2gort.Ext_java_util_Set) *DefaultNameGenerator {
	clinit_DefaultNameGenerator()
	this_ := &DefaultNameGenerator{}
	InitDefaultNameGenerator_3(this_, reservedNames, prefix, reservedCharacters)
	return this_
}

func InitDefaultNameGenerator_3(this_ *DefaultNameGenerator, reservedNames j2gort.Ext_java_util_Set, prefix string, reservedCharacters j2gort.Ext_java_util_Set) {
	_ = reservedNames
	_ = prefix
	_ = reservedCharacters
	InitDefaultNameGenerator_4_Ext_java_util_Set_string_Ext_java_util_Set_Ext_java_util_Set(this_, reservedNames, prefix, reservedCharacters, reservedCharacters)
}

// Java constructor: DefaultNameGenerator(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedFirstCharacters, java.util.Set<java.lang.Character> reservedNonFirstCharacters)
func NewDefaultNameGenerator_4_Ext_java_util_Set_string_Ext_java_util_Set_Ext_java_util_Set(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set) *DefaultNameGenerator {
	clinit_DefaultNameGenerator()
	this_ := &DefaultNameGenerator{}
	InitDefaultNameGenerator_4_Ext_java_util_Set_string_Ext_java_util_Set_Ext_java_util_Set(this_, reservedNames, prefix, reservedFirstCharacters, reservedNonFirstCharacters)
	return this_
}

func InitDefaultNameGenerator_4_Ext_java_util_Set_string_Ext_java_util_Set_Ext_java_util_Set(this_ *DefaultNameGenerator, reservedNames j2gort.Ext_java_util_Set, prefix string, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set) {
	_ = reservedNames
	_ = prefix
	_ = reservedFirstCharacters
	_ = reservedNonFirstCharacters
	this_.initFields_DefaultNameGenerator()
	this_.priv_DefaultNameGenerator_buildPriorityLookupMap_2355()
	this_.Reset_4(reservedNames, prefix, reservedFirstCharacters, reservedNonFirstCharacters)
}

// Java constructor: DefaultNameGenerator(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedCharacters, java.util.Map<java.lang.Character,com.google.javascript.jscomp.DefaultNameGenerator.CharPriority> priorityLookupMap)
func NewDefaultNameGenerator_4_Ext_java_util_Set_string_Ext_java_util_Set_Ext_java_util_Map(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedCharacters j2gort.Ext_java_util_Set, priorityLookupMap j2gort.Ext_java_util_Map) *DefaultNameGenerator {
	clinit_DefaultNameGenerator()
	this_ := &DefaultNameGenerator{}
	InitDefaultNameGenerator_4_Ext_java_util_Set_string_Ext_java_util_Set_Ext_java_util_Map(this_, reservedNames, prefix, reservedCharacters, priorityLookupMap)
	return this_
}

func InitDefaultNameGenerator_4_Ext_java_util_Set_string_Ext_java_util_Set_Ext_java_util_Map(this_ *DefaultNameGenerator, reservedNames j2gort.Ext_java_util_Set, prefix string, reservedCharacters j2gort.Ext_java_util_Set, priorityLookupMap j2gort.Ext_java_util_Map) {
	_ = reservedNames
	_ = prefix
	_ = reservedCharacters
	_ = priorityLookupMap
	this_.initFields_DefaultNameGenerator()
	this_.priorityLookupMap = j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_Maps_NewLinkedHashMapWithExpectedSize(int32(len(StaticGet_DefaultNameGenerator_NONFIRST_CHAR()))))
	{
		cur_1 := j2gort.NewMapEntryCursor(priorityLookupMap)
		for cur_1.NextEntry() {
			this_.priorityLookupMap.Put_2(j2gort.PtrToAny(j2gort.CastPtr[j2gort.BoxChar](cur_1.GetKey())), j2gort.PtrToAny(j2gort.CastPtr[DefaultNameGenerator_CharPriority](cur_1.GetValue()).Clone()))
		}
	}
	this_.Reset_4(reservedNames, prefix, reservedCharacters, reservedCharacters)
}

// Java method: DefaultNameGenerator.buildPriorityLookupMap()
func (this_ *DefaultNameGenerator) priv_DefaultNameGenerator_buildPriorityLookupMap_2355() {
	this_.priorityLookupMap = j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_Maps_NewLinkedHashMapWithExpectedSize(int32(len(StaticGet_DefaultNameGenerator_NONFIRST_CHAR()))))
	var order int32 = int32(0)
	_ = order
	for _, el_2 := range j2gort.CheckNotNullSlice(StaticGet_DefaultNameGenerator_NONFIRST_CHAR()) {
		c := el_2
		_ = c
		this_.priorityLookupMap.Put_2(j2gort.ValueOfChar(c), NewDefaultNameGenerator_CharPriority(c, order))
		order = order + 1
	}
}

// Java static method: DefaultNameGenerator.isBadName(java.lang.String name)
func DefaultNameGenerator_IsBadName(name string) bool {
	clinit_DefaultNameGenerator()
	_ = name
	return j2gort.SetContainsStr(StaticGet_DefaultNameGenerator_BAD_NAMES(), name)
}

// Java method: DefaultNameGenerator.reset(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedCharacters)
func (this_ *DefaultNameGenerator) Reset_3(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedCharacters j2gort.Ext_java_util_Set) {
	_ = reservedNames
	_ = prefix
	_ = reservedCharacters
	this_.Reset_4(reservedNames, prefix, reservedCharacters, reservedCharacters)
}

// Java method: DefaultNameGenerator.reset(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedFirstCharacters, java.util.Set<java.lang.Character> reservedNonFirstCharacters)
func (this_ *DefaultNameGenerator) Reset_4(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set) {
	_ = reservedNames
	_ = prefix
	_ = reservedFirstCharacters
	_ = reservedNonFirstCharacters
	this_.reservedNames = reservedNames
	this_.prefix = prefix
	this_.nameCount = int32(0)
	this_.firstChars = this_.ReserveCharacters(StaticGet_DefaultNameGenerator_FIRST_CHAR(), reservedFirstCharacters)
	this_.nonFirstChars = this_.ReserveCharacters(StaticGet_DefaultNameGenerator_NONFIRST_CHAR(), reservedNonFirstCharacters)
	func() {
		arrMut_3 := this_.firstChars
		var arrBox_4 []any
		if arrMut_3 != nil {
			arrBox_4 = make([]any, len(arrMut_3))
			for i_, el_5 := range arrMut_3 {
				arrBox_4[i_] = j2gort.PtrToAny(el_5)
			}
		}
		j2gort.Ext_java_util_Arrays_Sort_1_slice_any(arrBox_4)
		for i_, el_5 := range arrBox_4 {
			arrMut_3[i_] = j2gort.CastPtr[DefaultNameGenerator_CharPriority](el_5)
		}
	}()
	func() {
		arrMut_6 := this_.nonFirstChars
		var arrBox_7 []any
		if arrMut_6 != nil {
			arrBox_7 = make([]any, len(arrMut_6))
			for i_, el_8 := range arrMut_6 {
				arrBox_7[i_] = j2gort.PtrToAny(el_8)
			}
		}
		j2gort.Ext_java_util_Arrays_Sort_1_slice_any(arrBox_7)
		for i_, el_8 := range arrBox_7 {
			arrMut_6[i_] = j2gort.CastPtr[DefaultNameGenerator_CharPriority](el_8)
		}
	}()
	this_.priv_DefaultNameGenerator_checkPrefix_2361(prefix)
}

// Java method: DefaultNameGenerator.clone(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedCharacters)
func (this_ *DefaultNameGenerator) Clone(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedCharacters j2gort.Ext_java_util_Set) NameGenerator {
	_ = reservedNames
	_ = prefix
	_ = reservedCharacters
	return NewDefaultNameGenerator_4_Ext_java_util_Set_string_Ext_java_util_Set_Ext_java_util_Map(reservedNames, prefix, reservedCharacters, this_.priorityLookupMap)
}

// Java method: DefaultNameGenerator.favors(java.lang.CharSequence sequence)
func (this_ *DefaultNameGenerator) Favors(sequence any) {
	_ = sequence
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.CharSeqLength(sequence); i = i + 1 {
			var c *DefaultNameGenerator_CharPriority = j2gort.CastPtr[DefaultNameGenerator_CharPriority](this_.priorityLookupMap.Get(j2gort.ValueOfChar(j2gort.CharSeqCharAt(sequence, i))))
			_ = c
			if c != nil {
				r_9 := c
				r_9.occurrence = r_9.occurrence + 1
			}
		}
	}
}

// Java method: DefaultNameGenerator.reserveCharacters(char[] chars, java.util.Set<java.lang.Character> reservedCharacterSet)
func (this_ *DefaultNameGenerator) ReserveCharacters(chars []uint16, reservedCharacterSet j2gort.Ext_java_util_Set) []*DefaultNameGenerator_CharPriority {
	_ = chars
	_ = reservedCharacterSet
	if reservedCharacterSet.IsEmpty() {
		var result []*DefaultNameGenerator_CharPriority = make([]*DefaultNameGenerator_CharPriority, int32(len(chars)))
		_ = result
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(chars)); i = i + 1 {
				{
					a_10 := result
					i_11 := i
					a_10[i_11] = j2gort.CastPtr[DefaultNameGenerator_CharPriority](this_.priorityLookupMap.Get(j2gort.ValueOfChar(chars[i])))
				}
			}
		}
		return result
	}
	var charSet j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_primitives_Chars_AsList(chars)))
	_ = charSet
	charSet.RemoveAll(j2gort.AsIface[j2gort.Ext_java_util_Collection](reservedCharacterSet))
	var result_2 []*DefaultNameGenerator_CharPriority = make([]*DefaultNameGenerator_CharPriority, charSet.Size())
	_ = result_2
	var index int32 = int32(0)
	_ = index
	{
		cur_12 := j2gort.NewIterCursor(charSet)
		for cur_12.HasNext() {
			rawElem_13 := cur_12.Next()
			c := j2gort.UnboxChar(j2gort.Cast[*j2gort.BoxChar](rawElem_13))
			_ = c
			{
				a_14 := result_2
				i_15 := func() int32 { old_16 := index; nw_17 := old_16 + 1; index = nw_17; return old_16 }()
				a_14[i_15] = j2gort.CastPtr[DefaultNameGenerator_CharPriority](this_.priorityLookupMap.Get(j2gort.ValueOfChar(c)))
			}
		}
	}
	return result_2
}

// Java method: DefaultNameGenerator.checkPrefix(java.lang.String prefix)
func (this_ *DefaultNameGenerator) priv_DefaultNameGenerator_checkPrefix_2361(prefix string) {
	_ = prefix
	if j2gort.JStrLen(prefix) > int32(0) {
		if !(DefaultNameGenerator_Contains(this_.firstChars, j2gort.JStrCharAt(prefix, int32(0)))) {
			var chars []uint16 = make([]uint16, int32(len(this_.firstChars)))
			_ = chars
			{
				var i int32 = int32(0)
				_ = i
				for ; i < int32(len(chars)); i = i + 1 {
					{
						a_18 := chars
						i_19 := i
						a_18[i_19] = this_.firstChars[i].name
					}
				}
			}
			panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "prefix must start with one of: "+j2gort.StrValString(j2gort.Ext_java_util_Arrays_ToString_1_slice_uint16(chars)))))
		}
		{
			var pos int32 = int32(1)
			_ = pos
			for ; pos < j2gort.JStrLen(prefix); pos = pos + 1 {
				var chars_2 []uint16 = make([]uint16, int32(len(this_.nonFirstChars)))
				_ = chars_2
				{
					var i_2 int32 = int32(0)
					_ = i_2
					for ; i_2 < int32(len(chars_2)); i_2 = i_2 + 1 {
						{
							a_20 := chars_2
							i_21 := i_2
							a_20[i_21] = this_.nonFirstChars[i_2].name
						}
					}
				}
				if !(DefaultNameGenerator_Contains(this_.nonFirstChars, j2gort.JStrCharAt(prefix, pos))) {
					panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "prefix has invalid characters, must be one of: "+j2gort.StrValString(j2gort.Ext_java_util_Arrays_ToString_1_slice_uint16(chars_2)))))
				}
			}
		}
	}
}

// Java static method: DefaultNameGenerator.contains(com.google.javascript.jscomp.DefaultNameGenerator.CharPriority[] arr, char c)
func DefaultNameGenerator_Contains(arr []*DefaultNameGenerator_CharPriority, c uint16) bool {
	clinit_DefaultNameGenerator()
	_ = arr
	_ = c
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(arr)); i = i + 1 {
			if int32(arr[i].name) == int32(c) {
				return true
			}
		}
	}
	return false
}

// Java method: DefaultNameGenerator.generateNextName()
func (this_ *DefaultNameGenerator) GenerateNextName() string {
	var name string = j2gort.NullString
	_ = name
	for doFirst_22 := true; doFirst_22 || ((pkg_rhino.TokenStream_IsKeyword(name) || j2gort.SetContainsStr(this_.reservedNames, name)) || DefaultNameGenerator_IsBadName(name)); doFirst_22 = false {
		var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = sb
		var i int32 = this_.nameCount
		_ = i
		if !(j2gort.CheckNotNullString(this_.prefix) == "") {
			sb.Append_1_string(this_.prefix)
		}
		if sb.Length() == int32(0) {
			var pos int32 = (i % int32(len(this_.firstChars)))
			_ = pos
			sb.Append_1_uint16(this_.firstChars[pos].name)
			i = (i / int32(len(this_.firstChars)))
		}
		for i > int32(0) {
			i = i - 1
			var pos_2 int32 = (i % int32(len(this_.nonFirstChars)))
			_ = pos_2
			sb.Append_1_uint16(this_.nonFirstChars[pos_2].name)
			i = (i / int32(len(this_.nonFirstChars)))
		}
		this_.nameCount = this_.nameCount + 1
		name = sb.String()
	}
	return name
}

// Java class: com.google.javascript.jscomp.DefaultNameGenerator.CharPriority
type DefaultNameGenerator_CharPriority struct {
	occurrence int32 `gson:"occurrence"`
	order      int32 `gson:"order"`
	idHash_    atomic.Int32
	name       uint16 `gson:"name"`
}

func (this_ *DefaultNameGenerator_CharPriority) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultNameGenerator_CharPriority) IsJavaNil_() bool { return this_ == nil }
func (*DefaultNameGenerator_CharPriority) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DefaultNameGenerator$CharPriority", "CharPriority"
}
func (this_ *DefaultNameGenerator_CharPriority) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultNameGenerator_CharPriority) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DefaultNameGenerator.CharPriority", "com.google.javascript.jscomp.DefaultNameGenerator$CharPriority", func(v_ any) bool {
		return j2gort.InstanceOf[*DefaultNameGenerator_CharPriority](v_)
	})
}

func (this_ *DefaultNameGenerator_CharPriority) initFields_DefaultNameGenerator_CharPriority() {
	_ = this_
}

// Java constructor: CharPriority(char name, int order)
func NewDefaultNameGenerator_CharPriority(name uint16, order int32) *DefaultNameGenerator_CharPriority {
	this_ := &DefaultNameGenerator_CharPriority{}
	InitDefaultNameGenerator_CharPriority(this_, name, order)
	return this_
}

func InitDefaultNameGenerator_CharPriority(this_ *DefaultNameGenerator_CharPriority, name uint16, order int32) {
	_ = name
	_ = order
	this_.initFields_DefaultNameGenerator_CharPriority()
	this_.name = name
	this_.order = order
	this_.occurrence = int32(0)
}

// Java method: CharPriority.clone()
func (this_ *DefaultNameGenerator_CharPriority) Clone() *DefaultNameGenerator_CharPriority {
	var val_23 DefaultNameGenerator_CharPriority
	InitDefaultNameGenerator_CharPriority(&val_23, this_.name, this_.order)
	var result *DefaultNameGenerator_CharPriority = &val_23
	_ = result
	{
		r_24 := result
		r_24.occurrence = this_.occurrence
	}
	return result
}

// Java method: CharPriority.compareTo(com.google.javascript.jscomp.DefaultNameGenerator.CharPriority other)
func (this_ *DefaultNameGenerator_CharPriority) CompareTo(arg_25 any) int32 {
	other := j2gort.CastPtr[DefaultNameGenerator_CharPriority](arg_25)
	_ = other
	var result int32 = (other.occurrence - this_.occurrence)
	_ = result
	if result != int32(0) {
		return result
	}
	result = (this_.order - other.order)
	return result
}

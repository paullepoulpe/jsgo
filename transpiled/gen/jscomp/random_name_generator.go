/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RandomNameGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.RandomNameGenerator
type RandomNameGenerator struct {
	firstChars       j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"firstChars"`
	nonFirstChars    j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"nonFirstChars"`
	random           *j2gort.Ext_java_util_Random                      `gson:"random"`
	reservedNames    j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"reservedNames"`
	prefix           string                                            `gson:"prefix"`
	shuffledFirst    string                                            `gson:"shuffledFirst"`
	shuffledNonFirst []string                                          `gson:"shuffledNonFirst"`
	nameCount        int32                                             `gson:"nameCount"`
	idHash_          atomic.Int32
}

func (this_ *RandomNameGenerator) BindSelf_(s any) {
	_ = s
}
func (this_ *RandomNameGenerator) IsJavaNil_() bool { return this_ == nil }
func (*RandomNameGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RandomNameGenerator", "RandomNameGenerator"
}
func (this_ *RandomNameGenerator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RandomNameGenerator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RandomNameGenerator", "com.google.javascript.jscomp.RandomNameGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[*RandomNameGenerator](v_)
	})
}
func (*RandomNameGenerator) Implements_jscomp_NameGenerator() {}

var static_RandomNameGenerator_FIRST_CHAR j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_RandomNameGenerator_FIRST_CHAR() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_RandomNameGenerator()
	return static_RandomNameGenerator_FIRST_CHAR
}
func StaticSet_RandomNameGenerator_FIRST_CHAR(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_RandomNameGenerator()
	static_RandomNameGenerator_FIRST_CHAR = v
}

var static_RandomNameGenerator_NONFIRST_CHAR j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_RandomNameGenerator_NONFIRST_CHAR() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_RandomNameGenerator()
	return static_RandomNameGenerator_NONFIRST_CHAR
}
func StaticSet_RandomNameGenerator_NONFIRST_CHAR(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_RandomNameGenerator()
	static_RandomNameGenerator_NONFIRST_CHAR = v
}
func StaticGet_RandomNameGenerator_NUM_SHUFFLES() int32 { return int32(16) }

var clinitDone_RandomNameGenerator uint32
var clinitState_RandomNameGenerator j2gort.ClinitState

func clinit_RandomNameGenerator() {
	if atomic.LoadUint32(&clinitDone_RandomNameGenerator) != 0 {
		return
	}
	clinitSlow_RandomNameGenerator()
}

func clinitSlow_RandomNameGenerator() {
	if !j2gort.EnterClinit(&clinitState_RandomNameGenerator, &clinitDone_RandomNameGenerator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RandomNameGenerator, &clinitDone_RandomNameGenerator)
	static_RandomNameGenerator_FIRST_CHAR = RandomNameGenerator_AsSet(StaticGet_DefaultNameGenerator_FIRST_CHAR())
	static_RandomNameGenerator_NONFIRST_CHAR = RandomNameGenerator_AsSet(StaticGet_DefaultNameGenerator_NONFIRST_CHAR())
}

func (this_ *RandomNameGenerator) initFields_RandomNameGenerator() {
	_ = this_
	this_.prefix = j2gort.NullString
	this_.shuffledFirst = j2gort.NullString
}

// Java constructor: RandomNameGenerator(java.util.Random random)
func NewRandomNameGenerator_1(random *j2gort.Ext_java_util_Random) *RandomNameGenerator {
	clinit_RandomNameGenerator()
	this_ := &RandomNameGenerator{}
	InitRandomNameGenerator_1(this_, random)
	return this_
}

func InitRandomNameGenerator_1(this_ *RandomNameGenerator, random *j2gort.Ext_java_util_Random) {
	_ = random
	this_.initFields_RandomNameGenerator()
	this_.random = random
	this_.Reset_3(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()), "", j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()))
}

// Java constructor: RandomNameGenerator(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedCharacters, java.util.Random random)
func NewRandomNameGenerator_4(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedCharacters j2gort.Ext_java_util_Set, random *j2gort.Ext_java_util_Random) *RandomNameGenerator {
	clinit_RandomNameGenerator()
	this_ := &RandomNameGenerator{}
	InitRandomNameGenerator_4(this_, reservedNames, prefix, reservedCharacters, random)
	return this_
}

func InitRandomNameGenerator_4(this_ *RandomNameGenerator, reservedNames j2gort.Ext_java_util_Set, prefix string, reservedCharacters j2gort.Ext_java_util_Set, random *j2gort.Ext_java_util_Random) {
	_ = reservedNames
	_ = prefix
	_ = reservedCharacters
	_ = random
	this_.initFields_RandomNameGenerator()
	this_.random = random
	this_.Reset_4(reservedNames, prefix, reservedCharacters, reservedCharacters)
}

// Java constructor: RandomNameGenerator(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedFirstCharacters, java.util.Set<java.lang.Character> reservedNonFirstCharacters, java.util.Random random)
func NewRandomNameGenerator_5(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set, random *j2gort.Ext_java_util_Random) *RandomNameGenerator {
	clinit_RandomNameGenerator()
	this_ := &RandomNameGenerator{}
	InitRandomNameGenerator_5(this_, reservedNames, prefix, reservedFirstCharacters, reservedNonFirstCharacters, random)
	return this_
}

func InitRandomNameGenerator_5(this_ *RandomNameGenerator, reservedNames j2gort.Ext_java_util_Set, prefix string, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set, random *j2gort.Ext_java_util_Random) {
	_ = reservedNames
	_ = prefix
	_ = reservedFirstCharacters
	_ = reservedNonFirstCharacters
	_ = random
	this_.initFields_RandomNameGenerator()
	this_.random = random
	this_.Reset_4(reservedNames, prefix, reservedFirstCharacters, reservedNonFirstCharacters)
}

// Java method: RandomNameGenerator.reset(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedCharacters)
func (this_ *RandomNameGenerator) Reset_3(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedCharacters j2gort.Ext_java_util_Set) {
	_ = reservedNames
	_ = prefix
	_ = reservedCharacters
	this_.Reset_4(reservedNames, prefix, reservedCharacters, reservedCharacters)
}

// Java method: RandomNameGenerator.reset(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedFirstCharacters, java.util.Set<java.lang.Character> reservedNonFirstCharacters)
func (this_ *RandomNameGenerator) Reset_4(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set) {
	_ = reservedNames
	_ = prefix
	_ = reservedFirstCharacters
	_ = reservedNonFirstCharacters
	this_.reservedNames = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](reservedNames))
	this_.prefix = prefix
	this_.nameCount = int32(0)
	this_.firstChars = j2gort.Ext_com_google_common_collect_Sets_Difference(j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_RandomNameGenerator_FIRST_CHAR()), reservedFirstCharacters).ImmutableCopy()
	this_.nonFirstChars = j2gort.Ext_com_google_common_collect_Sets_Difference(j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_RandomNameGenerator_NONFIRST_CHAR()), reservedNonFirstCharacters).ImmutableCopy()
	this_.priv_RandomNameGenerator_checkPrefix_4610(prefix)
	this_.priv_RandomNameGenerator_shuffleAlphabets_4611()
}

// Java method: RandomNameGenerator.clone(java.util.Set<java.lang.String> reservedNames, java.lang.String prefix, java.util.Set<java.lang.Character> reservedCharacters)
func (this_ *RandomNameGenerator) Clone(reservedNames j2gort.Ext_java_util_Set, prefix string, reservedCharacters j2gort.Ext_java_util_Set) NameGenerator {
	_ = reservedNames
	_ = prefix
	_ = reservedCharacters
	return NewRandomNameGenerator_4(reservedNames, prefix, reservedCharacters, this_.random)
}

// Java static method: RandomNameGenerator.asSet(char @org.jspecify.annotations.Nullable [] chars)
func RandomNameGenerator_AsSet(chars []uint16) j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_RandomNameGenerator()
	_ = chars
	if chars == nil {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_primitives_Chars_AsList(chars)))
	}
}

// Java method: RandomNameGenerator.checkPrefix(java.lang.String prefix)
func (this_ *RandomNameGenerator) priv_RandomNameGenerator_checkPrefix_4610(prefix string) {
	_ = prefix
	if j2gort.JStrLen(prefix) > int32(0) {
		if !(this_.firstChars.Contains(j2gort.ValueOfChar(j2gort.JStrCharAt(prefix, int32(0))))) {
			panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "prefix must start with one of: "+j2gort.StrValString(j2gort.Ext_com_google_common_base_Joiner_On_1_string(", ").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.firstChars))))))
		}
		{
			var pos int32 = int32(1)
			_ = pos
			for ; pos < j2gort.JStrLen(prefix); pos = pos + 1 {
				if !(this_.nonFirstChars.Contains(j2gort.ValueOfChar(j2gort.JStrCharAt(prefix, pos)))) {
					panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "prefix has invalid characters, must be one of: "+j2gort.StrValString(j2gort.Ext_com_google_common_base_Joiner_On_1_string(", ").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.nonFirstChars))))))
				}
			}
		}
	}
}

// Java static method: RandomNameGenerator.shuffleAndCopyAlphabet(java.util.Set<java.lang.Character> input, java.util.Random random)
func RandomNameGenerator_ShuffleAndCopyAlphabet(input j2gort.Ext_java_util_Set, random *j2gort.Ext_java_util_Random) string {
	clinit_RandomNameGenerator()
	_ = input
	_ = random
	var shuffled j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](input))
	_ = shuffled
	j2gort.Ext_java_util_Collections_Shuffle_2_Ext_java_util_List_Ext_java_util_Random(shuffled, random)
	return javalang.StringFromUTF16(j2gort.Ext_com_google_common_primitives_Chars_ToArray(shuffled))
}

// Java method: RandomNameGenerator.shuffleAlphabets()
func (this_ *RandomNameGenerator) priv_RandomNameGenerator_shuffleAlphabets_4611() {
	this_.shuffledFirst = RandomNameGenerator_ShuffleAndCopyAlphabet(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.firstChars), this_.random)
	var builder []string = []string{}
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(16); i = i + 1 {
			builder = append(builder, j2gort.CheckNotNullString(RandomNameGenerator_ShuffleAndCopyAlphabet(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.nonFirstChars), this_.random)))
		}
	}
	this_.shuffledNonFirst = builder
}

// Java method: RandomNameGenerator.getNameLength(int position, int nameIdx)
func (this_ *RandomNameGenerator) priv_RandomNameGenerator_getNameLength_4612(position int32, nameIdx int32) int32 {
	_ = position
	_ = nameIdx
	var length int32 = int32(0)
	_ = length
	nameIdx = nameIdx + 1
	for doFirst_1 := true; doFirst_1 || (nameIdx > int32(0)); doFirst_1 = false {
		nameIdx = nameIdx - 1
		var alphabetSize int32
		_ = alphabetSize
		if position == int32(0) {
			alphabetSize = this_.firstChars.Size()
		} else {
			alphabetSize = this_.nonFirstChars.Size()
		}
		nameIdx = (nameIdx / alphabetSize)
		position = position + 1
		length = length + 1
	}
	return length
}

// Java method: RandomNameGenerator.generateSuffix(int position, int nameIdx)
func (this_ *RandomNameGenerator) priv_RandomNameGenerator_generateSuffix_4613(position int32, nameIdx int32) string {
	_ = position
	_ = nameIdx
	var name *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = name
	var length int32 = this_.priv_RandomNameGenerator_getNameLength_4612(position, nameIdx)
	_ = length
	nameIdx = nameIdx + 1
	for doFirst_2 := true; doFirst_2 || (nameIdx > int32(0)); doFirst_2 = false {
		nameIdx = nameIdx - 1
		var alphabet string = j2gort.NullString
		_ = alphabet
		if position == int32(0) {
			alphabet = this_.shuffledFirst
		} else {
			var hasher j2gort.Ext_com_google_common_hash_Hasher = j2gort.Ext_com_google_common_hash_Hashing_Murmur3_128_0().NewHasher_0()
			_ = hasher
			hasher.PutInt(length)
			hasher.PutUnencodedChars(name)
			var alphabetIdx int32 = ((hasher.Hash().AsInt() & int32(2147483647)) % int32(16))
			_ = alphabetIdx
			alphabet = this_.shuffledNonFirst[alphabetIdx]
		}
		var alphabetSize int32 = j2gort.JStrLen(alphabet)
		_ = alphabetSize
		var character uint16 = j2gort.JStrCharAt(alphabet, (nameIdx % alphabetSize))
		_ = character
		name.Append_1_uint16(character)
		nameIdx = (nameIdx / alphabetSize)
		position = position + 1
	}
	return name.String()
}

// Java method: RandomNameGenerator.generateNextName()
func (this_ *RandomNameGenerator) GenerateNextName() string {
	for {
		var name string = (j2gort.StrValString(this_.prefix) + j2gort.StrValString(this_.priv_RandomNameGenerator_generateSuffix_4613(j2gort.JStrLen(this_.prefix), func() int32 { old_3 := this_.nameCount; nw_4 := old_3 + 1; this_.nameCount = nw_4; return old_3 }())))
		_ = name
		if (pkg_rhino.TokenStream_IsKeyword(name) || j2gort.SetContainsStr(this_.reservedNames, name)) || DefaultNameGenerator_IsBadName(name) {
			continue
		}
		return name
	}
	panic("unreachable")
}

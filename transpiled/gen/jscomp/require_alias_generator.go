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
// Port of com/google/javascript/refactoring/RequireAliasGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.refactoring.RequireAliasGenerator
type RequireAliasGenerator struct {
	parts      []string `gson:"parts"`
	idHash_    atomic.Int32
	capitalize bool `gson:"capitalize"`
}

func (this_ *RequireAliasGenerator) BindSelf_(s any) {
	_ = s
}
func (this_ *RequireAliasGenerator) IsJavaNil_() bool { return this_ == nil }
func (*RequireAliasGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.RequireAliasGenerator", "RequireAliasGenerator"
}
func (this_ *RequireAliasGenerator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RequireAliasGenerator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.RequireAliasGenerator", "com.google.javascript.refactoring.RequireAliasGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[*RequireAliasGenerator](v_)
	})
}
func (*RequireAliasGenerator) Implements_j2gort_Ext_java_lang_Iterable() {}

var static_RequireAliasGenerator_SKIPLISTED_ALIASES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_RequireAliasGenerator_SKIPLISTED_ALIASES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_RequireAliasGenerator()
	return static_RequireAliasGenerator_SKIPLISTED_ALIASES
}
func StaticSet_RequireAliasGenerator_SKIPLISTED_ALIASES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_RequireAliasGenerator()
	static_RequireAliasGenerator_SKIPLISTED_ALIASES = v
}

var static_RequireAliasGenerator_DOT_SPLITTER *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_RequireAliasGenerator_DOT_SPLITTER() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_RequireAliasGenerator()
	return static_RequireAliasGenerator_DOT_SPLITTER
}
func StaticSet_RequireAliasGenerator_DOT_SPLITTER(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_RequireAliasGenerator()
	static_RequireAliasGenerator_DOT_SPLITTER = v
}

var clinitDone_RequireAliasGenerator uint32
var clinitState_RequireAliasGenerator j2gort.ClinitState

func clinit_RequireAliasGenerator() {
	if atomic.LoadUint32(&clinitDone_RequireAliasGenerator) != 0 {
		return
	}
	clinitSlow_RequireAliasGenerator()
}

func clinitSlow_RequireAliasGenerator() {
	if !j2gort.EnterClinit(&clinitState_RequireAliasGenerator, &clinitDone_RequireAliasGenerator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RequireAliasGenerator, &clinitDone_RequireAliasGenerator)
	static_RequireAliasGenerator_SKIPLISTED_ALIASES = RequireAliasGenerator_CreateSkiplistedAliases()
	static_RequireAliasGenerator_DOT_SPLITTER = j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(46))
}

func (this_ *RequireAliasGenerator) initFields_RequireAliasGenerator() {
	_ = this_
}

// Java constructor: RequireAliasGenerator(java.lang.String namespace)
func NewRequireAliasGenerator(namespace string) *RequireAliasGenerator {
	clinit_RequireAliasGenerator()
	this_ := &RequireAliasGenerator{}
	InitRequireAliasGenerator(this_, namespace)
	return this_
}

func InitRequireAliasGenerator(this_ *RequireAliasGenerator, namespace string) {
	_ = namespace
	this_.initFields_RequireAliasGenerator()
	this_.parts = j2gort.ImmListToSliceNullable[string](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(j2gort.SplitterSplitStr(StaticGet_RequireAliasGenerator_DOT_SPLITTER(), namespace)).Reverse())
	this_.capitalize = j2gort.Ext_com_google_common_base_Ascii_IsUpperCase(j2gort.JStrCharAt(this_.parts[int32(0)], int32(0)))
}

// Java static method: RequireAliasGenerator.createSkiplistedAliases()
func RequireAliasGenerator_CreateSkiplistedAliases() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_RequireAliasGenerator()
	var seeds j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("", "array", "event", "map", "math", "object", []any{"promise", "set", "string", "thenable"})
	_ = seeds
	var builder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = builder
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(seeds))
		for cur_1.HasNext() {
			seed := cur_1.NextStr()
			_ = seed
			j2gort.ImmSetBuilderAddStrVoid(builder, RequireAliasGenerator_CaptializeFirstChar(seed))
			j2gort.ImmSetBuilderAddStrVoid(builder, j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(seed))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](builder.Build())
}

// Java static method: RequireAliasGenerator.over(java.lang.String namespace)
func RequireAliasGenerator_Over(namespace string) j2gort.Ext_java_lang_Iterable {
	clinit_RequireAliasGenerator()
	_ = namespace
	return NewRequireAliasGenerator(namespace)
}

// Java method: RequireAliasGenerator.iterator()
func (this_ *RequireAliasGenerator) Iterator() any {
	return NewRequireAliasGenerator_NameIterator(this_)
}

// Java static method: RequireAliasGenerator.captializeFirstChar(java.lang.String w)
func RequireAliasGenerator_CaptializeFirstChar(w string) string {
	clinit_RequireAliasGenerator()
	_ = w
	if (j2gort.CheckNotNullString(w) == "") || j2gort.Ext_com_google_common_base_Ascii_IsUpperCase(j2gort.JStrCharAt(w, int32(0))) {
		return w
	}
	return (j2gort.StrValChar(j2gort.Ext_com_google_common_base_Ascii_ToUpperCase_1_uint16(j2gort.JStrCharAt(w, int32(0)))) + j2gort.StrValString(j2gort.JStrSubFrom(w, int32(1))))
}

func (this_ *RequireAliasGenerator) ForEach(p0_ j2gort.Ext_java_util_function_Consumer) {
	j2gort.ExtDefault_Iterable_ForEach(this_, p0_)
}

// Java class: com.google.javascript.refactoring.RequireAliasGenerator.NameIterator
type RequireAliasGenerator_NameIterator struct {
	j2gort.Ext_com_google_common_collect_AbstractIteratorBase
	outer_  *RequireAliasGenerator
	concat  string `gson:"concat"`
	index   int32  `gson:"index"`
	suffix  int32  `gson:"suffix"`
	idHash_ atomic.Int32
}

func (this_ *RequireAliasGenerator_NameIterator) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_common_collect_AbstractIteratorBase.BindExtSelf_(s)
}
func (this_ *RequireAliasGenerator_NameIterator) IsJavaNil_() bool { return this_ == nil }
func (*RequireAliasGenerator_NameIterator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.RequireAliasGenerator$NameIterator", "NameIterator"
}
func (this_ *RequireAliasGenerator_NameIterator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.RequireAliasGenerator.NameIterator", "com.google.javascript.refactoring.RequireAliasGenerator$NameIterator", func(v_ any) bool {
		return j2gort.InstanceOf[*RequireAliasGenerator_NameIterator](v_)
	})
}
func (this_ *RequireAliasGenerator_NameIterator) SetOuter_(o *RequireAliasGenerator) {
	this_.outer_ = o
}
func (this_ *RequireAliasGenerator_NameIterator) GetOuter_() *RequireAliasGenerator {
	return this_.outer_
}
func (*RequireAliasGenerator_NameIterator) Implements_j2gort_Ext_com_google_common_collect_AbstractIterator() {
}
func (*RequireAliasGenerator_NameIterator) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {
}
func (*RequireAliasGenerator_NameIterator) Implements_j2gort_Ext_java_util_Iterator() {}

func (this_ *RequireAliasGenerator_NameIterator) initFields_RequireAliasGenerator_NameIterator() {
	_ = this_

}

// Java constructor: NameIterator()
func NewRequireAliasGenerator_NameIterator(outer_ *RequireAliasGenerator) *RequireAliasGenerator_NameIterator {
	this_ := &RequireAliasGenerator_NameIterator{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRequireAliasGenerator_NameIterator(this_)
	return this_
}

func InitRequireAliasGenerator_NameIterator(this_ *RequireAliasGenerator_NameIterator) {
	j2gort.InitExt_com_google_common_collect_AbstractIterator(&this_.Ext_com_google_common_collect_AbstractIteratorBase)
	this_.initFields_RequireAliasGenerator_NameIterator()
}

// Java method: NameIterator.computeNext()
func (this_ *RequireAliasGenerator_NameIterator) ComputeNext() any {
	var next string = j2gort.NullString
	_ = next
	for doFirst_2 := true; doFirst_2 || (j2gort.SetContainsStr(StaticGet_RequireAliasGenerator_SKIPLISTED_ALIASES(), next)); doFirst_2 = false {
		next = this_.priv_NameIterator_computeNextWithoutFiltering_7925()
	}
	return j2gort.BoxString(next)
}

// Java method: NameIterator.computeNextWithoutFiltering()
func (this_ *RequireAliasGenerator_NameIterator) priv_NameIterator_computeNextWithoutFiltering_7925() string {
	if this_.index >= int32(len(this_.outer_.parts)) {
		return (j2gort.StrValString(this_.concat) + j2gort.StrValInt32(func() int32 { old_3 := this_.suffix; nw_4 := old_3 + 1; this_.suffix = nw_4; return old_3 }()))
	}
	var part string = this_.outer_.parts[func() int32 { old_5 := this_.index; nw_6 := old_5 + 1; this_.index = nw_6; return old_5 }()]
	_ = part
	if this_.outer_.capitalize {
		part = RequireAliasGenerator_CaptializeFirstChar(part)
	}
	this_.concat = (j2gort.StrValString(part) + j2gort.StrValString(RequireAliasGenerator_CaptializeFirstChar(this_.concat)))
	return this_.concat
}

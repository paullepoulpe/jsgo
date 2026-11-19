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
// Port of com/google/javascript/refactoring/ApplySuggestedFixes.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_apply_suggested_fixes_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_apply_suggested_fixes_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_5 any) any {
		f := j2gort.CastPtr[SuggestedFix](lp_5)
		_ = f
		return j2gort.ValueOfInt(j2gort.ImmSliceSize(j2gort.CheckNotNullPtr(f).Slice_GetAlternatives()))
	})
}

// Java class: com.google.javascript.refactoring.ApplySuggestedFixes
type ApplySuggestedFixes struct {
	idHash_ atomic.Int32
}

func (this_ *ApplySuggestedFixes) BindSelf_(s any) {
	_ = s
}
func (this_ *ApplySuggestedFixes) IsJavaNil_() bool { return this_ == nil }
func (*ApplySuggestedFixes) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.ApplySuggestedFixes", "ApplySuggestedFixes"
}
func (this_ *ApplySuggestedFixes) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ApplySuggestedFixes) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.ApplySuggestedFixes", "com.google.javascript.refactoring.ApplySuggestedFixes", func(v_ any) bool {
		return j2gort.InstanceOf[*ApplySuggestedFixes](v_)
	})
}

var static_ApplySuggestedFixes_DOUBLE_LINE_JOINER *j2gort.Ext_com_google_common_base_Joiner

func StaticGet_ApplySuggestedFixes_DOUBLE_LINE_JOINER() *j2gort.Ext_com_google_common_base_Joiner {
	clinit_ApplySuggestedFixes()
	return static_ApplySuggestedFixes_DOUBLE_LINE_JOINER
}
func StaticSet_ApplySuggestedFixes_DOUBLE_LINE_JOINER(v *j2gort.Ext_com_google_common_base_Joiner) {
	clinit_ApplySuggestedFixes()
	static_ApplySuggestedFixes_DOUBLE_LINE_JOINER = v
}

var clinitDone_ApplySuggestedFixes uint32
var clinitState_ApplySuggestedFixes j2gort.ClinitState

func clinit_ApplySuggestedFixes() {
	if atomic.LoadUint32(&clinitDone_ApplySuggestedFixes) != 0 {
		return
	}
	clinitSlow_ApplySuggestedFixes()
}

func clinitSlow_ApplySuggestedFixes() {
	if !j2gort.EnterClinit(&clinitState_ApplySuggestedFixes, &clinitDone_ApplySuggestedFixes) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ApplySuggestedFixes, &clinitDone_ApplySuggestedFixes)
	static_ApplySuggestedFixes_DOUBLE_LINE_JOINER = j2gort.Ext_com_google_common_base_Joiner_On_1_string("\n\n")
}

func (this_ *ApplySuggestedFixes) initFields_ApplySuggestedFixes() {
	_ = this_
}

// Java constructor: ApplySuggestedFixes()
func NewApplySuggestedFixes() *ApplySuggestedFixes {
	clinit_ApplySuggestedFixes()
	this_ := &ApplySuggestedFixes{}
	InitApplySuggestedFixes(this_)
	return this_
}

func InitApplySuggestedFixes(this_ *ApplySuggestedFixes) {
	this_.initFields_ApplySuggestedFixes()
}

// Java static method: ApplySuggestedFixes.applySuggestedFixesToFiles(java.lang.Iterable<com.google.javascript.refactoring.SuggestedFix> fixes)
func ApplySuggestedFixes_ApplySuggestedFixesToFiles(fixes j2gort.Ext_java_lang_Iterable) {
	clinit_ApplySuggestedFixes()
	_ = fixes
	var filenames j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = filenames
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(fixes))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			fix := j2gort.CastPtr[SuggestedFix](rawElem_2)
			_ = fix
			j2gort.LocalOrderedSetAddAllStr(&filenames, j2gort.AsIface[j2gort.Ext_java_util_Set](fix.GetReplacements().KeySet()))
		}
	}
	var filenameToCodeMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_HashMap()
	_ = filenameToCodeMap
	{
		cur_3 := filenames.Iter()
		for cur_3.HasNext() {
			filename := cur_3.
				Next()
			_ = filename
			j2gort.MapPutStr(filenameToCodeMap, filename, j2gort.BoxString(j2gort.Ext_com_google_common_io_Files_AsCharSource(j2gort.NewExt_java_io_File_1_string(filename), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8).Read()))
		}
	}
	var newCode j2gort.Ext_com_google_common_collect_ImmutableMap = ApplySuggestedFixes_ApplySuggestedFixesToCode(fixes, filenameToCodeMap)
	_ = newCode
	{
		cur_4 := j2gort.NewMapEntryCursor(newCode)
		for cur_4.NextEntry() {
			j2gort.Ext_com_google_common_io_Files_AsCharSink(j2gort.NewExt_java_io_File_1_string(cur_4.KeyStr()), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8, []*j2gort.Ext_com_google_common_io_FileWriteMode{}).Write(j2gort.BoxString(j2gort.CastString(cur_4.GetValue())))
		}
	}
}

// Java static method: ApplySuggestedFixes.applyAllSuggestedFixChoicesToCode(java.lang.Iterable<com.google.javascript.refactoring.SuggestedFix> fixChoices, java.util.Map<java.lang.String,java.lang.String> fileNameToCodeMap)
func ApplySuggestedFixes_ApplyAllSuggestedFixChoicesToCode(fixChoices j2gort.Ext_java_lang_Iterable, fileNameToCodeMap j2gort.Ext_java_util_Map) j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_ApplySuggestedFixes()
	_ = fixChoices
	_ = fileNameToCodeMap
	if j2gort.Ext_com_google_common_collect_Iterables_IsEmpty(fixChoices) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0())
	}
	var alternativeCount int32 = j2gort.ImmSliceSize(j2gort.CheckNotNullPtr(j2gort.CastPtr[SuggestedFix](j2gort.Ext_com_google_common_collect_Iterables_GetFirst(fixChoices, nil))).Slice_GetAlternatives())
	_ = alternativeCount
	if !(j2gort.Ext_com_google_common_collect_Streams_Stream_1_Ext_java_lang_Iterable(fixChoices).Map(lambdaSingleton_apply_suggested_fixes_1).AllMatch(j2gort.Ext_java_util_function_Predicate_IsEqual(j2gort.ValueOfInt(alternativeCount)))) {
		j2gort.CheckArgumentStr(false, "All SuggestedFixAlternatives must offer an equal number of choices for this utility to make sense")
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](j2gort.Ext_java_util_stream_IntStream_Range(int32(0), alternativeCount).MapToObj(j2gort.NewLambda_Ext_java_util_function_IntFunction(func(i int32) any {
		_ = i
		return ApplySuggestedFixes_ApplySuggestedFixChoicesToCode(fixChoices, i, fileNameToCodeMap)
	})).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))
}

// Java static method: ApplySuggestedFixes.applySuggestedFixChoicesToCode(java.lang.Iterable<com.google.javascript.refactoring.SuggestedFix> fixChoices, int choiceIndex, java.util.Map<java.lang.String,java.lang.String> fileNameToCodeMap)
func ApplySuggestedFixes_ApplySuggestedFixChoicesToCode(fixChoices j2gort.Ext_java_lang_Iterable, choiceIndex int32, fileNameToCodeMap j2gort.Ext_java_util_Map) j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_ApplySuggestedFixes()
	_ = fixChoices
	_ = choiceIndex
	_ = fileNameToCodeMap
	var chosenFixes j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](j2gort.Ext_com_google_common_collect_Streams_Stream_1_Ext_java_lang_Iterable(fixChoices).Map(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_6 any) any {
		choices := j2gort.CastPtr[SuggestedFix](lp_6)
		_ = choices
		return j2gort.PtrToAny(j2gort.CheckNotNullSlice(j2gort.CheckNotNullPtr(choices).Slice_GetAlternatives())[choiceIndex])
	})).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))
	_ = chosenFixes
	return ApplySuggestedFixes_ApplySuggestedFixesToCode(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](chosenFixes), fileNameToCodeMap)
}

// Java static method: ApplySuggestedFixes.applySuggestedFixesToCode(java.lang.Iterable<com.google.javascript.refactoring.SuggestedFix> fixes, java.util.Map<java.lang.String,java.lang.String> filenameToCodeMap)
func ApplySuggestedFixes_ApplySuggestedFixesToCode(fixes j2gort.Ext_java_lang_Iterable, filenameToCodeMap j2gort.Ext_java_util_Map) j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_ApplySuggestedFixes()
	_ = fixes
	_ = filenameToCodeMap
	var val_7 ApplySuggestedFixes_ReplacementMap
	InitApplySuggestedFixes_ReplacementMap(&val_7)
	var map_ *ApplySuggestedFixes_ReplacementMap = &val_7
	_ = map_
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(fixes))
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			fix := j2gort.CastPtr[SuggestedFix](rawElem_9)
			_ = fix
			map_.PutIfNoOverlap(fix)
		}
	}
	var newCodeMap j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = newCodeMap
	{
		cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(map_.EntrySet()))
		for cur_10.HasNext() {
			rawElem_11 := cur_10.Next()
			entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](rawElem_11)
			_ = entry
			var filename string = j2gort.CastString(entry.GetKey())
			_ = filename
			if !(j2gort.MapContainsKeyStr(filenameToCodeMap, filename)) {
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "filenameToCodeMap missing code for file: "+j2gort.StrValString(filename))))
			}
			var replacements j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](entry.GetValue())
			_ = replacements
			var newCode string = ApplySuggestedFixes_ApplyCodeReplacements(j2gort.AsIface[j2gort.Ext_java_util_Collection](replacements), j2gort.CastString(j2gort.MapGetStr(filenameToCodeMap, filename)))
			_ = newCode
			newCodeMap.Put_2(j2gort.BoxString(filename), j2gort.BoxString(newCode))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](newCodeMap.BuildOrThrow())
}

// Java static method: ApplySuggestedFixes.applyCodeReplacements(java.util.Collection<com.google.javascript.refactoring.CodeReplacement> replacements, java.lang.String code)
func ApplySuggestedFixes_ApplyCodeReplacements(replacements j2gort.Ext_java_util_Collection, code string) string {
	clinit_ApplySuggestedFixes()
	_ = replacements
	_ = code
	var sortedReplacements j2gort.Ext_com_google_common_collect_ImmutableSortedSet = j2gort.Ext_com_google_common_collect_ImmutableSortedSet_CopyOf_1_Ext_java_util_Collection(replacements)
	_ = sortedReplacements
	j2gort.CheckArgumentFmt(!(ApplySuggestedFixes_ContainsOverlaps(sortedReplacements)), "Found overlap between code replacements:\n%s", j2gort.BoxString(StaticGet_ApplySuggestedFixes_DOUBLE_LINE_JOINER().Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](sortedReplacements))))
	var sb j2gort.LocalStringBuilder

	var lastIndex int32 = int32(0)
	_ = lastIndex
	{
		cur_12 := j2gort.NewIterCursor(j2gort.CheckNotNull(sortedReplacements))
		for cur_12.HasNext() {
			rawElem_13 := cur_12.Next()
			replacement := j2gort.AsIface[CodeReplacement](rawElem_13)
			_ = replacement
			sb.AppendSub(j2gort.BoxString(code), lastIndex, replacement.GetStartPosition())
			sb.AppendStr(replacement.GetNewContent())
			lastIndex = replacement.GetEndPosition()
		}
	}
	if lastIndex <= j2gort.JStrLen(code) {
		sb.AppendSub(j2gort.BoxString(code), lastIndex, j2gort.JStrLen(code))
	}
	return sb.String()
}

// Java static method: ApplySuggestedFixes.containsOverlaps(com.google.common.collect.ImmutableSortedSet<com.google.javascript.refactoring.CodeReplacement> replacements)
func ApplySuggestedFixes_ContainsOverlaps(replacements j2gort.Ext_com_google_common_collect_ImmutableSortedSet) bool {
	clinit_ApplySuggestedFixes()
	_ = replacements
	var start int32 = int32(-1)
	_ = start
	{
		cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNull(replacements))
		for cur_14.HasNext() {
			rawElem_15 := cur_14.Next()
			replacement := j2gort.AsIface[CodeReplacement](rawElem_15)
			_ = replacement
			if replacement.GetStartPosition() < start {
				return true
			}
			start = replacement.GetEndPosition()
		}
	}
	return false
}

// Java class: com.google.javascript.refactoring.ApplySuggestedFixes.ReplacementMap
type ApplySuggestedFixes_ReplacementMap struct {
	map_    j2gort.Ext_com_google_common_collect_SetMultimap `gson:"map"`
	idHash_ atomic.Int32
}

func (this_ *ApplySuggestedFixes_ReplacementMap) BindSelf_(s any) {
	_ = s
}
func (this_ *ApplySuggestedFixes_ReplacementMap) IsJavaNil_() bool { return this_ == nil }
func (*ApplySuggestedFixes_ReplacementMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.ApplySuggestedFixes$ReplacementMap", "ReplacementMap"
}
func (this_ *ApplySuggestedFixes_ReplacementMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ApplySuggestedFixes_ReplacementMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.ApplySuggestedFixes.ReplacementMap", "com.google.javascript.refactoring.ApplySuggestedFixes$ReplacementMap", func(v_ any) bool {
		return j2gort.InstanceOf[*ApplySuggestedFixes_ReplacementMap](v_)
	})
}

func (this_ *ApplySuggestedFixes_ReplacementMap) initFields_ApplySuggestedFixes_ReplacementMap() {
	_ = this_
	this_.map_ = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_HashMultimap_Create_0())
}

// Java constructor: ReplacementMap()
func NewApplySuggestedFixes_ReplacementMap() *ApplySuggestedFixes_ReplacementMap {
	this_ := &ApplySuggestedFixes_ReplacementMap{}
	InitApplySuggestedFixes_ReplacementMap(this_)
	return this_
}

func InitApplySuggestedFixes_ReplacementMap(this_ *ApplySuggestedFixes_ReplacementMap) {
	this_.initFields_ApplySuggestedFixes_ReplacementMap()
}

// Java method: ReplacementMap.putIfNoOverlap(com.google.javascript.refactoring.SuggestedFix fix)
func (this_ *ApplySuggestedFixes_ReplacementMap) PutIfNoOverlap(fix *SuggestedFix) {
	_ = fix
	if this_.priv_ReplacementMap_canPut_7855(fix) {
		this_.map_.PutAll_1(j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](fix.GetReplacements()))
	}
}

// Java method: ReplacementMap.canPut(com.google.javascript.refactoring.SuggestedFix fix)
func (this_ *ApplySuggestedFixes_ReplacementMap) priv_ReplacementMap_canPut_7855(fix *SuggestedFix) bool {
	_ = fix
	{
		cur_16 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](fix.GetReplacements().KeySet())))
		for cur_16.HasNext() {
			filename := cur_16.NextStr()
			_ = filename
			var replacements j2gort.Ext_com_google_common_collect_ImmutableSortedSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSortedSet](j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSortedSet_Builder](j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSortedSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSortedSet_NaturalOrder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.MultimapGetStr(this_.map_, filename)))).AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.MultimapGetStr(fix.GetReplacements(), filename)))).Build())
			_ = replacements
			if ApplySuggestedFixes_ContainsOverlaps(replacements) {
				return false
			}
		}
	}
	return true
}

// Java method: ReplacementMap.entrySet()
func (this_ *ApplySuggestedFixes_ReplacementMap) EntrySet() j2gort.Ext_java_util_Set {
	return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Multimaps_S_AsMap_1_Ext_com_google_common_collect_SetMultimap(this_.map_).EntrySet())
}

/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RescopeGlobalSymbolsInfo.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.RescopeGlobalSymbolsInfo
type RescopeGlobalSymbolsInfo struct {
	globalSymbolNamespace                                  string                                            `gson:"globalSymbolNamespace"`
	externNames                                            j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"externNames"`
	crossChunkNames                                        j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"crossChunkNames"`
	crossChunkNamesWithWriteFromOtherChunk                 j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"crossChunkNamesWithWriteFromOtherChunk"`
	crossChunkNamesUsedInNonDependentChunk                 j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"crossChunkNamesUsedInNonDependentChunk"`
	globalNamesWithReadInDefiningChunk                     j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"globalNamesWithReadInDefiningChunk"`
	globalNamesWithInnerScopeWriteInDefiningChunk          j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"globalNamesWithInnerScopeWriteInDefiningChunk"`
	globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk"`
	maybeReferencesThis                                    j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"maybeReferencesThis"`
	wrappedCrossChunkNames                                 j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"wrappedCrossChunkNames"`
}

func (this_ *RescopeGlobalSymbolsInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *RescopeGlobalSymbolsInfo) IsJavaNil_() bool { return this_ == nil }
func (*RescopeGlobalSymbolsInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RescopeGlobalSymbolsInfo", "RescopeGlobalSymbolsInfo"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RescopeGlobalSymbolsInfo", "com.google.javascript.jscomp.RescopeGlobalSymbolsInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*RescopeGlobalSymbolsInfo](v_)
	})
}
func (*RescopeGlobalSymbolsInfo) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *RescopeGlobalSymbolsInfo) initFields_RescopeGlobalSymbolsInfo() {
	_ = this_
}

// Java constructor: RescopeGlobalSymbolsInfo(java.lang.String globalSymbolNamespace, com.google.common.collect.ImmutableSet<java.lang.String> externNames, com.google.common.collect.ImmutableSet<java.lang.String> crossChunkNames, com.google.common.collect.ImmutableSet<java.lang.String> crossChunkNamesWithWriteFromOtherChunk, com.google.common.collect.ImmutableSet<java.lang.String> crossChunkNamesUsedInNonDependentChunk, com.google.common.collect.ImmutableSet<java.lang.String> globalNamesWithReadInDefiningChunk, com.google.common.collect.ImmutableSet<java.lang.String> globalNamesWithInnerScopeWriteInDefiningChunk, com.google.common.collect.ImmutableSet<java.lang.String> globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk, com.google.common.collect.ImmutableSet<java.lang.String> maybeReferencesThis, com.google.common.collect.ImmutableSet<java.lang.String> wrappedCrossChunkNames)
func NewRescopeGlobalSymbolsInfo_10(globalSymbolNamespace string, externNames j2gort.Ext_com_google_common_collect_ImmutableSet, crossChunkNames j2gort.Ext_com_google_common_collect_ImmutableSet, crossChunkNamesWithWriteFromOtherChunk j2gort.Ext_com_google_common_collect_ImmutableSet, crossChunkNamesUsedInNonDependentChunk j2gort.Ext_com_google_common_collect_ImmutableSet, globalNamesWithReadInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableSet, globalNamesWithInnerScopeWriteInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableSet, globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableSet, maybeReferencesThis j2gort.Ext_com_google_common_collect_ImmutableSet, wrappedCrossChunkNames j2gort.Ext_com_google_common_collect_ImmutableSet) *RescopeGlobalSymbolsInfo {
	this_ := &RescopeGlobalSymbolsInfo{}
	InitRescopeGlobalSymbolsInfo_10(this_, globalSymbolNamespace, externNames, crossChunkNames, crossChunkNamesWithWriteFromOtherChunk, crossChunkNamesUsedInNonDependentChunk, globalNamesWithReadInDefiningChunk, globalNamesWithInnerScopeWriteInDefiningChunk, globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk, maybeReferencesThis, wrappedCrossChunkNames)
	return this_
}

func InitRescopeGlobalSymbolsInfo_10(this_ *RescopeGlobalSymbolsInfo, globalSymbolNamespace string, externNames j2gort.Ext_com_google_common_collect_ImmutableSet, crossChunkNames j2gort.Ext_com_google_common_collect_ImmutableSet, crossChunkNamesWithWriteFromOtherChunk j2gort.Ext_com_google_common_collect_ImmutableSet, crossChunkNamesUsedInNonDependentChunk j2gort.Ext_com_google_common_collect_ImmutableSet, globalNamesWithReadInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableSet, globalNamesWithInnerScopeWriteInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableSet, globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk j2gort.Ext_com_google_common_collect_ImmutableSet, maybeReferencesThis j2gort.Ext_com_google_common_collect_ImmutableSet, wrappedCrossChunkNames j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = globalSymbolNamespace
	_ = externNames
	_ = crossChunkNames
	_ = crossChunkNamesWithWriteFromOtherChunk
	_ = crossChunkNamesUsedInNonDependentChunk
	_ = globalNamesWithReadInDefiningChunk
	_ = globalNamesWithInnerScopeWriteInDefiningChunk
	_ = globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk
	_ = maybeReferencesThis
	_ = wrappedCrossChunkNames
	this_.initFields_RescopeGlobalSymbolsInfo()
	this_.globalSymbolNamespace = globalSymbolNamespace
	this_.externNames = externNames
	this_.crossChunkNames = crossChunkNames
	this_.crossChunkNamesWithWriteFromOtherChunk = crossChunkNamesWithWriteFromOtherChunk
	this_.crossChunkNamesUsedInNonDependentChunk = crossChunkNamesUsedInNonDependentChunk
	this_.globalNamesWithReadInDefiningChunk = globalNamesWithReadInDefiningChunk
	this_.globalNamesWithInnerScopeWriteInDefiningChunk = globalNamesWithInnerScopeWriteInDefiningChunk
	this_.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk = globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk
	this_.maybeReferencesThis = maybeReferencesThis
	this_.wrappedCrossChunkNames = wrappedCrossChunkNames
}

// Java constructor: RescopeGlobalSymbolsInfo(java.lang.String globalSymbolNamespace, java.util.Set<java.lang.String> externNames, java.util.Set<java.lang.String> crossChunkNames, java.util.Set<java.lang.String> crossChunkNamesWithWriteFromOtherChunk, java.util.Set<java.lang.String> crossChunkNamesUsedInNonDependentChunk, java.util.Set<java.lang.String> globalNamesWithReadInDefiningChunk, java.util.Set<java.lang.String> globalNamesWithInnerScopeWriteInDefiningChunk, java.util.Set<java.lang.String> globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk, java.util.Set<java.lang.String> maybeReferencesThis)
func NewRescopeGlobalSymbolsInfo_9(globalSymbolNamespace string, externNames j2gort.Ext_java_util_Set, crossChunkNames j2gort.Ext_java_util_Set, crossChunkNamesWithWriteFromOtherChunk j2gort.Ext_java_util_Set, crossChunkNamesUsedInNonDependentChunk j2gort.Ext_java_util_Set, globalNamesWithReadInDefiningChunk j2gort.Ext_java_util_Set, globalNamesWithInnerScopeWriteInDefiningChunk j2gort.Ext_java_util_Set, globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk j2gort.Ext_java_util_Set, maybeReferencesThis j2gort.Ext_java_util_Set) *RescopeGlobalSymbolsInfo {
	this_ := &RescopeGlobalSymbolsInfo{}
	InitRescopeGlobalSymbolsInfo_9(this_, globalSymbolNamespace, externNames, crossChunkNames, crossChunkNamesWithWriteFromOtherChunk, crossChunkNamesUsedInNonDependentChunk, globalNamesWithReadInDefiningChunk, globalNamesWithInnerScopeWriteInDefiningChunk, globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk, maybeReferencesThis)
	return this_
}

func InitRescopeGlobalSymbolsInfo_9(this_ *RescopeGlobalSymbolsInfo, globalSymbolNamespace string, externNames j2gort.Ext_java_util_Set, crossChunkNames j2gort.Ext_java_util_Set, crossChunkNamesWithWriteFromOtherChunk j2gort.Ext_java_util_Set, crossChunkNamesUsedInNonDependentChunk j2gort.Ext_java_util_Set, globalNamesWithReadInDefiningChunk j2gort.Ext_java_util_Set, globalNamesWithInnerScopeWriteInDefiningChunk j2gort.Ext_java_util_Set, globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk j2gort.Ext_java_util_Set, maybeReferencesThis j2gort.Ext_java_util_Set) {
	_ = globalSymbolNamespace
	_ = externNames
	_ = crossChunkNames
	_ = crossChunkNamesWithWriteFromOtherChunk
	_ = crossChunkNamesUsedInNonDependentChunk
	_ = globalNamesWithReadInDefiningChunk
	_ = globalNamesWithInnerScopeWriteInDefiningChunk
	_ = globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk
	_ = maybeReferencesThis
	InitRescopeGlobalSymbolsInfo_10(this_, globalSymbolNamespace, j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](externNames)), j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](crossChunkNames)), j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](crossChunkNamesWithWriteFromOtherChunk)), j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](crossChunkNamesUsedInNonDependentChunk)), j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](globalNamesWithReadInDefiningChunk)), j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](globalNamesWithInnerScopeWriteInDefiningChunk)), j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk)), j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](maybeReferencesThis)), RescopeGlobalSymbolsInfo_CollectWrappedCrossChunkNames(externNames, crossChunkNames, crossChunkNamesWithWriteFromOtherChunk, crossChunkNamesUsedInNonDependentChunk, globalNamesWithInnerScopeWriteInDefiningChunk))
}

// Java static method: RescopeGlobalSymbolsInfo.collectWrappedCrossChunkNames(java.util.Set<java.lang.String> externNames, java.util.Set<java.lang.String> crossChunkNames, java.util.Set<java.lang.String> crossChunkNamesWithWriteFromOtherChunk, java.util.Set<java.lang.String> crossChunkNamesUsedInNonDependentChunk, java.util.Set<java.lang.String> globalNamesWithInnerScopeWriteInDefiningChunk)
func RescopeGlobalSymbolsInfo_CollectWrappedCrossChunkNames(externNames j2gort.Ext_java_util_Set, crossChunkNames j2gort.Ext_java_util_Set, crossChunkNamesWithWriteFromOtherChunk j2gort.Ext_java_util_Set, crossChunkNamesUsedInNonDependentChunk j2gort.Ext_java_util_Set, globalNamesWithInnerScopeWriteInDefiningChunk j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = externNames
	_ = crossChunkNames
	_ = crossChunkNamesWithWriteFromOtherChunk
	_ = crossChunkNamesUsedInNonDependentChunk
	_ = globalNamesWithInnerScopeWriteInDefiningChunk
	var candidates j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = candidates
	candidates.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](crossChunkNamesWithWriteFromOtherChunk))
	candidates.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](crossChunkNamesUsedInNonDependentChunk))
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(globalNamesWithInnerScopeWriteInDefiningChunk))
		for cur_1.HasNext() {
			name := cur_1.NextStr()
			_ = name
			if j2gort.SetContainsStr(crossChunkNames, name) {
				j2gort.SetAddStr(candidates, name)
			}
		}
	}
	candidates.RemoveAll(j2gort.AsIface[j2gort.Ext_java_util_Collection](externNames))
	return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(candidates)
}

func (this_ *RescopeGlobalSymbolsInfo) GlobalSymbolNamespace() string {
	return this_.globalSymbolNamespace
}
func (this_ *RescopeGlobalSymbolsInfo) ExternNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.externNames
}
func (this_ *RescopeGlobalSymbolsInfo) CrossChunkNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.crossChunkNames
}
func (this_ *RescopeGlobalSymbolsInfo) CrossChunkNamesWithWriteFromOtherChunk() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.crossChunkNamesWithWriteFromOtherChunk
}
func (this_ *RescopeGlobalSymbolsInfo) CrossChunkNamesUsedInNonDependentChunk() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.crossChunkNamesUsedInNonDependentChunk
}
func (this_ *RescopeGlobalSymbolsInfo) GlobalNamesWithReadInDefiningChunk() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.globalNamesWithReadInDefiningChunk
}
func (this_ *RescopeGlobalSymbolsInfo) GlobalNamesWithInnerScopeWriteInDefiningChunk() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.globalNamesWithInnerScopeWriteInDefiningChunk
}
func (this_ *RescopeGlobalSymbolsInfo) GlobalNamesWithUnsafeWriteForLocalAliasInDefiningChunk() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk
}
func (this_ *RescopeGlobalSymbolsInfo) MaybeReferencesThis() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.maybeReferencesThis
}
func (this_ *RescopeGlobalSymbolsInfo) WrappedCrossChunkNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.wrappedCrossChunkNames
}
func (this_ *RescopeGlobalSymbolsInfo) Equals(other_ any) bool {
	o_, ok_ := other_.(*RescopeGlobalSymbolsInfo)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.globalSymbolNamespace != o_.globalSymbolNamespace {
		return false
	}
	if !javalang.ObjectsEquals(this_.externNames, o_.externNames) {
		return false
	}
	if !javalang.ObjectsEquals(this_.crossChunkNames, o_.crossChunkNames) {
		return false
	}
	if !javalang.ObjectsEquals(this_.crossChunkNamesWithWriteFromOtherChunk, o_.crossChunkNamesWithWriteFromOtherChunk) {
		return false
	}
	if !javalang.ObjectsEquals(this_.crossChunkNamesUsedInNonDependentChunk, o_.crossChunkNamesUsedInNonDependentChunk) {
		return false
	}
	if !javalang.ObjectsEquals(this_.globalNamesWithReadInDefiningChunk, o_.globalNamesWithReadInDefiningChunk) {
		return false
	}
	if !javalang.ObjectsEquals(this_.globalNamesWithInnerScopeWriteInDefiningChunk, o_.globalNamesWithInnerScopeWriteInDefiningChunk) {
		return false
	}
	if !javalang.ObjectsEquals(this_.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk, o_.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk) {
		return false
	}
	if !javalang.ObjectsEquals(this_.maybeReferencesThis, o_.maybeReferencesThis) {
		return false
	}
	if !javalang.ObjectsEquals(this_.wrappedCrossChunkNames, o_.wrappedCrossChunkNames) {
		return false
	}
	return true
}
func (this_ *RescopeGlobalSymbolsInfo) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_2 := this_.globalSymbolNamespace
		h_ = 31 * h_
		if hs_2 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_2)
		}
	}
	h_ = 31*h_ + j2gort.ObjectHash(this_.externNames)
	h_ = 31*h_ + j2gort.ObjectHash(this_.crossChunkNames)
	h_ = 31*h_ + j2gort.ObjectHash(this_.crossChunkNamesWithWriteFromOtherChunk)
	h_ = 31*h_ + j2gort.ObjectHash(this_.crossChunkNamesUsedInNonDependentChunk)
	h_ = 31*h_ + j2gort.ObjectHash(this_.globalNamesWithReadInDefiningChunk)
	h_ = 31*h_ + j2gort.ObjectHash(this_.globalNamesWithInnerScopeWriteInDefiningChunk)
	h_ = 31*h_ + j2gort.ObjectHash(this_.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk)
	h_ = 31*h_ + j2gort.ObjectHash(this_.maybeReferencesThis)
	h_ = 31*h_ + j2gort.ObjectHash(this_.wrappedCrossChunkNames)
	return h_
}
func (this_ *RescopeGlobalSymbolsInfo) String() string {
	_ = fmt.Sprint
	return "RescopeGlobalSymbolsInfo[" + "globalSymbolNamespace=" + j2gort.StrVal(j2gort.BoxString(this_.globalSymbolNamespace)) + ", " + "externNames=" + j2gort.StrVal(this_.externNames) + ", " + "crossChunkNames=" + j2gort.StrVal(this_.crossChunkNames) + ", " + "crossChunkNamesWithWriteFromOtherChunk=" + j2gort.StrVal(this_.crossChunkNamesWithWriteFromOtherChunk) + ", " + "crossChunkNamesUsedInNonDependentChunk=" + j2gort.StrVal(this_.crossChunkNamesUsedInNonDependentChunk) + ", " + "globalNamesWithReadInDefiningChunk=" + j2gort.StrVal(this_.globalNamesWithReadInDefiningChunk) + ", " + "globalNamesWithInnerScopeWriteInDefiningChunk=" + j2gort.StrVal(this_.globalNamesWithInnerScopeWriteInDefiningChunk) + ", " + "globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk=" + j2gort.StrVal(this_.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk) + ", " + "maybeReferencesThis=" + j2gort.StrVal(this_.maybeReferencesThis) + ", " + "wrappedCrossChunkNames=" + j2gort.StrVal(this_.wrappedCrossChunkNames) + "]"
}

/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/disambiguate/UseSiteRenamer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.disambiguate.UseSiteRenamer
type UseSiteRenamer struct {
	mutationCb    j2gort.Ext_java_util_function_Consumer                             `gson:"mutationCb"`
	renamingIndex *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_Builder `gson:"renamingIndex"`
	idHash_       atomic.Int32
}

func (this_ *UseSiteRenamer) BindSelf_(s any) {
	_ = s
}
func (this_ *UseSiteRenamer) IsJavaNil_() bool { return this_ == nil }
func (*UseSiteRenamer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.UseSiteRenamer", "UseSiteRenamer"
}
func (this_ *UseSiteRenamer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*UseSiteRenamer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.UseSiteRenamer", "com.google.javascript.jscomp.disambiguate.UseSiteRenamer", func(v_ any) bool {
		return j2gort.InstanceOf[*UseSiteRenamer](v_)
	})
}

func StaticGet_UseSiteRenamer_INVALIDATED_NAME_VALUE() string { return "<INVALIDATED>" }
func (this_ *UseSiteRenamer) initFields_UseSiteRenamer() {
	_ = this_
	this_.renamingIndex = j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_S_Builder()
}

// Java constructor: UseSiteRenamer(java.util.function.Consumer<com.google.javascript.rhino.Node> mutationCb)
func NewUseSiteRenamer(mutationCb j2gort.Ext_java_util_function_Consumer) *UseSiteRenamer {
	this_ := &UseSiteRenamer{}
	InitUseSiteRenamer(this_, mutationCb)
	return this_
}

func InitUseSiteRenamer(this_ *UseSiteRenamer, mutationCb j2gort.Ext_java_util_function_Consumer) {
	_ = mutationCb
	this_.initFields_UseSiteRenamer()
	this_.mutationCb = mutationCb
}

// Java method: UseSiteRenamer.renameUses(com.google.javascript.jscomp.disambiguate.PropertyClustering prop)
func (this_ *UseSiteRenamer) RenameUses(prop *PropertyClustering) UseSiteRenamer_RenameUsesResult {
	_ = prop
	if prop.IsInvalidated() {
		j2gort.MultimapBuilderPutStr(this_.renamingIndex, prop.GetName(), "<INVALIDATED>")
		return StaticGet_UseSiteRenamer_RenameUsesResult_INVALIDATED()
	}
	var clusterNames j2gort.Ext_com_google_common_collect_ImmutableMap = UseSiteRenamer_CreateAllClusterNames(prop)
	_ = clusterNames
	if clusterNames.Size() <= int32(1) {
		j2gort.MultimapBuilderPutStr(this_.renamingIndex, prop.GetName(), j2gort.BoxString(prop.GetName()))
		return StaticGet_UseSiteRenamer_RenameUsesResult_ONLY_ONE_CLUSTER()
	}
	this_.renamingIndex.PutAll_2_any_Ext_java_lang_Iterable(j2gort.BoxString(prop.GetName()), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](clusterNames.Values()))
	prop.ForEachUseSite(j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(lp_1 any, lp_2 any) {
		site := j2gort.CastNode[pkg_rhino.Node](lp_1)
		_ = site
		owner := j2gort.CastPtr[ColorGraphNode](lp_2)
		_ = owner
		var flatRep *ColorGraphNode = prop.GetClusters().Find(owner)
		_ = flatRep
		var newName string = j2gort.CastString(clusterNames.Get(j2gort.PtrToAny(flatRep)))
		_ = newName
		if !(newName == site.GetString()) {
			site.SetString(newName)
			this_.mutationCb.Accept(j2gort.BoxNode(site))
		}
	}))
	return StaticGet_UseSiteRenamer_RenameUsesResult_DISAMBIGUATED()
}

// Java method: UseSiteRenamer.getRenamingIndex()
func (this_ *UseSiteRenamer) GetRenamingIndex() *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap {
	return j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSetMultimap](this_.renamingIndex.Build())
}

// Java static method: UseSiteRenamer.createAllClusterNames(com.google.javascript.jscomp.disambiguate.PropertyClustering prop)
func UseSiteRenamer_CreateAllClusterNames(prop *PropertyClustering) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = prop
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](prop.GetClusters().AllRepresentatives().Stream().Collect_1(j2gort.Ext_com_google_common_collect_ImmutableMap_ToImmutableMap_2(j2gort.Ext_java_util_function_Function_Identity(), j2gort.NewLambda_Ext_java_util_function_Function(func(lp_3 any) any {
		r := j2gort.CastPtr[ColorGraphNode](lp_3)
		_ = r
		return j2gort.BoxString(UseSiteRenamer_CreateClusterName(prop, r))
	}))))
}

// Java static method: UseSiteRenamer.createClusterName(com.google.javascript.jscomp.disambiguate.PropertyClustering prop, com.google.javascript.jscomp.disambiguate.ColorGraphNode rep)
func UseSiteRenamer_CreateClusterName(prop *PropertyClustering, rep *ColorGraphNode) string {
	_ = prop
	_ = rep
	if prop.GetOriginalNameClusterRep() == rep {
		return prop.GetName()
	}
	return ("JSC$" + j2gort.StrValInt32(rep.GetIndex()) + "_" + j2gort.StrValString(prop.GetName()))
}

// Java enum: com.google.javascript.jscomp.disambiguate.UseSiteRenamer.RenameUsesResult
type UseSiteRenamer_RenameUsesResult uint8

const (
	enum_UseSiteRenamer_RenameUsesResult_INVALIDATED      UseSiteRenamer_RenameUsesResult = 1
	enum_UseSiteRenamer_RenameUsesResult_ONLY_ONE_CLUSTER UseSiteRenamer_RenameUsesResult = 2
	enum_UseSiteRenamer_RenameUsesResult_DISAMBIGUATED    UseSiteRenamer_RenameUsesResult = 3
)

var enumBaseTable_UseSiteRenamer_RenameUsesResult = [...]j2gort.EnumBase{
	{Name_: "INVALIDATED", Ordinal_: 0},
	{Name_: "ONLY_ONE_CLUSTER", Ordinal_: 1},
	{Name_: "DISAMBIGUATED", Ordinal_: 2},
}
var enumValues_UseSiteRenamer_RenameUsesResult = []UseSiteRenamer_RenameUsesResult{1, 2, 3}

func StaticGet_UseSiteRenamer_RenameUsesResult_INVALIDATED() UseSiteRenamer_RenameUsesResult {
	return enum_UseSiteRenamer_RenameUsesResult_INVALIDATED
}
func StaticGet_UseSiteRenamer_RenameUsesResult_ONLY_ONE_CLUSTER() UseSiteRenamer_RenameUsesResult {
	return enum_UseSiteRenamer_RenameUsesResult_ONLY_ONE_CLUSTER
}
func StaticGet_UseSiteRenamer_RenameUsesResult_DISAMBIGUATED() UseSiteRenamer_RenameUsesResult {
	return enum_UseSiteRenamer_RenameUsesResult_DISAMBIGUATED
}

func UseSiteRenamer_RenameUsesResult_Values() []UseSiteRenamer_RenameUsesResult {
	return j2gort.CloneSlice(enumValues_UseSiteRenamer_RenameUsesResult)
}

func UseSiteRenamer_RenameUsesResult_ValueOf(name string) UseSiteRenamer_RenameUsesResult {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_UseSiteRenamer_RenameUsesResult {
		if enumBaseTable_UseSiteRenamer_RenameUsesResult[i_].Name_ == name {
			return UseSiteRenamer_RenameUsesResult(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.disambiguate.UseSiteRenamer.RenameUsesResult." + name))
}

func (this_ UseSiteRenamer_RenameUsesResult) IsJavaNil_() bool                              { return this_ == 0 }
func (UseSiteRenamer_RenameUsesResult) HasIdentityEquals_()                                 {}
func (UseSiteRenamer_RenameUsesResult) Implements_j2gort_Ext_java_lang_Enum()               {}
func (UseSiteRenamer_RenameUsesResult) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ UseSiteRenamer_RenameUsesResult) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (UseSiteRenamer_RenameUsesResult) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.UseSiteRenamer$RenameUsesResult", "RenameUsesResult"
}
func (this_ UseSiteRenamer_RenameUsesResult) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_UseSiteRenamer_RenameUsesResult[this_-1]
}
func (this_ UseSiteRenamer_RenameUsesResult) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ UseSiteRenamer_RenameUsesResult) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ UseSiteRenamer_RenameUsesResult) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_UseSiteRenamer_RenameUsesResult[this_-1].Name_
}
func (this_ UseSiteRenamer_RenameUsesResult) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_UseSiteRenamer_RenameUsesResult[this_-1].Name_
}
func (this_ UseSiteRenamer_RenameUsesResult) Equals(other_ any) bool {
	o_, ok_ := other_.(UseSiteRenamer_RenameUsesResult)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ UseSiteRenamer_RenameUsesResult) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_UseSiteRenamer_RenameUsesResult[this_-1].HashCode()
}
func (this_ UseSiteRenamer_RenameUsesResult) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_UseSiteRenamer_RenameUsesResult[this_-1].CompareTo(other_)
}
func (this_ UseSiteRenamer_RenameUsesResult) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_UseSiteRenamer_RenameUsesResult[this_-1].Name_)
}
func init() {
	vals_ := []any{UseSiteRenamer_RenameUsesResult(1), UseSiteRenamer_RenameUsesResult(2), UseSiteRenamer_RenameUsesResult(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.disambiguate.UseSiteRenamer.RenameUsesResult", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.disambiguate.UseSiteRenamer$RenameUsesResult", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.UseSiteRenamer.RenameUsesResult", "com.google.javascript.jscomp.disambiguate.UseSiteRenamer$RenameUsesResult", func(v_ any) bool {
		return j2gort.InstanceOf[UseSiteRenamer_RenameUsesResult](v_)
	})
}

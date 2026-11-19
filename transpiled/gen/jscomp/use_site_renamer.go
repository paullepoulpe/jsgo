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

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.disambiguate.UseSiteRenamer
type UseSiteRenamer struct {
	mutationCb    UseSiteRenamer_MutationCallback     `gson:"mutationCb"`
	recordedNames j2gort.OrderedMap[string, []string] `gson:"recordedNames"`
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

var static_UseSiteRenamer_INVALIDATED_MARKER []string

func StaticGet_UseSiteRenamer_INVALIDATED_MARKER() []string {
	clinit_UseSiteRenamer()
	return static_UseSiteRenamer_INVALIDATED_MARKER
}
func StaticSet_UseSiteRenamer_INVALIDATED_MARKER(v []string) {
	clinit_UseSiteRenamer()
	static_UseSiteRenamer_INVALIDATED_MARKER = v
}

var static_UseSiteRenamer_SINGLE_CLUSTER_MARKER []string

func StaticGet_UseSiteRenamer_SINGLE_CLUSTER_MARKER() []string {
	clinit_UseSiteRenamer()
	return static_UseSiteRenamer_SINGLE_CLUSTER_MARKER
}
func StaticSet_UseSiteRenamer_SINGLE_CLUSTER_MARKER(v []string) {
	clinit_UseSiteRenamer()
	static_UseSiteRenamer_SINGLE_CLUSTER_MARKER = v
}

var clinitDone_UseSiteRenamer uint32
var clinitState_UseSiteRenamer j2gort.ClinitState

func clinit_UseSiteRenamer() {
	if atomic.LoadUint32(&clinitDone_UseSiteRenamer) != 0 {
		return
	}
	clinitSlow_UseSiteRenamer()
}

func clinitSlow_UseSiteRenamer() {
	if !j2gort.EnterClinit(&clinitState_UseSiteRenamer, &clinitDone_UseSiteRenamer) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_UseSiteRenamer, &clinitDone_UseSiteRenamer)
	static_UseSiteRenamer_INVALIDATED_MARKER = []string{"<INVALIDATED>"}
	static_UseSiteRenamer_SINGLE_CLUSTER_MARKER = j2gort.NewStringArray(int32(0))
}

func (this_ *UseSiteRenamer) initFields_UseSiteRenamer() {
	_ = this_
	this_.recordedNames = j2gort.NewOrderedMap[string, []string]()
}

// Java constructor: UseSiteRenamer(com.google.javascript.jscomp.disambiguate.UseSiteRenamer.MutationCallback mutationCb)
func NewUseSiteRenamer(mutationCb UseSiteRenamer_MutationCallback) *UseSiteRenamer {
	clinit_UseSiteRenamer()
	this_ := &UseSiteRenamer{}
	InitUseSiteRenamer(this_, mutationCb)
	return this_
}

func InitUseSiteRenamer(this_ *UseSiteRenamer, mutationCb UseSiteRenamer_MutationCallback) {
	_ = mutationCb
	this_.initFields_UseSiteRenamer()
	this_.mutationCb = mutationCb
}

// Java method: UseSiteRenamer.renameUses(com.google.javascript.jscomp.disambiguate.PropertyClustering prop)
func (this_ *UseSiteRenamer) RenameUses(prop *PropertyClustering) UseSiteRenamer_RenameUsesResult {
	_ = prop
	var propName string = prop.GetName()
	_ = propName
	if prop.IsInvalidated() {
		this_.recordedNames.Put(propName, StaticGet_UseSiteRenamer_INVALIDATED_MARKER())
		return StaticGet_UseSiteRenamer_RenameUsesResult_INVALIDATED()
	}
	var clusters *pkg_graph.StandardUnionFind[*ColorGraphNode] = prop.GetClusters()
	_ = clusters
	if !(clusters.HasMultipleRepresentatives()) {
		this_.recordedNames.Put(propName, StaticGet_UseSiteRenamer_SINGLE_CLUSTER_MARKER())
		return StaticGet_UseSiteRenamer_RenameUsesResult_ONLY_ONE_CLUSTER()
	}
	var reps j2gort.Ext_java_util_ArrayList = clusters.GetRepresentativeList()
	_ = reps
	var repCount int32 = reps.Size()
	_ = repCount
	var clusterNames *pkg_base.LinkedIdentityHashMap[*ColorGraphNode, string] = pkg_base.NewLinkedIdentityHashMap[*ColorGraphNode, string]()
	_ = clusterNames
	var namesArr []string = j2gort.NewStringArray(repCount)
	_ = namesArr
	var origRep *ColorGraphNode = prop.GetOriginalNameClusterRep()
	_ = origRep
	{
		var i int32 = int32(0)
		_ = i
		for ; i < repCount; i = i + 1 {
			var rep *ColorGraphNode = j2gort.CastPtr[ColorGraphNode](reps.Get(i))
			_ = rep
			var cName string
			_ = cName
			if rep == origRep {
				cName = propName
			} else {
				cName = ("JSC$" + j2gort.StrValInt32(rep.GetIndex()) + "_" + j2gort.StrValString(propName))
			}
			clusterNames.Put(rep, cName)
			{
				a_1 := namesArr
				i_2 := i
				a_1[i_2] = cName
			}
		}
	}
	this_.recordedNames.Put(propName, namesArr)
	{
		cur_3 := prop.GetUseSites().IterEntries()
		for cur_3.NextEntry() {
			var site pkg_rhino.Node = cur_3.Key()
			_ = site
			var flatRep *ColorGraphNode = clusters.Find(cur_3.Value())
			_ = flatRep
			var newName string = clusterNames.Get(flatRep)
			_ = newName
			if !(newName == site.GetString()) {
				site.SetString(newName)
				this_.mutationCb.Accept(site)
			}
		}
	}
	return StaticGet_UseSiteRenamer_RenameUsesResult_DISAMBIGUATED()
}

// Java method: UseSiteRenamer.getRenamingIndex()
func (this_ *UseSiteRenamer) GetRenamingIndex() *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap {
	var builder *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_Builder = j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_S_Builder()
	_ = builder
	{
		cur_4 := this_.recordedNames.IterEntries()
		for cur_4.NextEntry() {
			var propName string = cur_4.Key()
			_ = propName
			var vals []string = cur_4.Value()
			_ = vals
			if j2gort.SliceRefEquals(vals, StaticGet_UseSiteRenamer_SINGLE_CLUSTER_MARKER()) {
				j2gort.MultimapBuilderPutStr(builder, propName, j2gort.BoxString(propName))
			} else if vals != nil {
				{
					var i int32 = int32(0)
					_ = i
					for ; i < int32(len(vals)); i = i + 1 {
						j2gort.MultimapBuilderPutStr(builder, propName, j2gort.BoxString(vals[i]))
					}
				}
			}
		}
	}
	return j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSetMultimap](builder.Build())
}

// Java interface: com.google.javascript.jscomp.disambiguate.UseSiteRenamer.MutationCallback
type UseSiteRenamer_MutationCallback interface {
	Implements_jscomp_UseSiteRenamer_MutationCallback()
	IsJavaNil_() bool
	Accept(pkg_rhino.Node)
}

func ToIface_UseSiteRenamer_MutationCallback[T interface {
	UseSiteRenamer_MutationCallback
	comparable
}](v T) UseSiteRenamer_MutationCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.UseSiteRenamer.MutationCallback", "com.google.javascript.jscomp.disambiguate.UseSiteRenamer$MutationCallback", func(v_ any) bool {
		return j2gort.InstanceOf[UseSiteRenamer_MutationCallback](v_)
	})
}

type lambda_UseSiteRenamer_MutationCallback struct {
	fn func(pkg_rhino.Node)
}

func (l *lambda_UseSiteRenamer_MutationCallback) IsJavaNil_() bool { return l == nil }
func NewLambda_UseSiteRenamer_MutationCallback(fn func(pkg_rhino.Node)) UseSiteRenamer_MutationCallback {
	if fn == nil {
		return nil
	}
	return &lambda_UseSiteRenamer_MutationCallback{fn: fn}
}

func (*lambda_UseSiteRenamer_MutationCallback) Implements_jscomp_UseSiteRenamer_MutationCallback() {}
func (l *lambda_UseSiteRenamer_MutationCallback) Accept(p0_ pkg_rhino.Node) {
	l.fn(p0_)
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

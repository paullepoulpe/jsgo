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
// Port of com/google/javascript/jscomp/DependencyOptions.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.DependencyOptions
type DependencyOptions struct {
	entryPoints []*ModuleIdentifier              `gson:"entryPoints"`
	mode        DependencyOptions_DependencyMode `gson:"mode"`
}

func (this_ *DependencyOptions) BindSelf_(s any) {
	_ = s
}
func (this_ *DependencyOptions) IsJavaNil_() bool { return this_ == nil }
func (*DependencyOptions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DependencyOptions", "DependencyOptions"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DependencyOptions", "com.google.javascript.jscomp.DependencyOptions", func(v_ any) bool {
		return j2gort.InstanceOf[*DependencyOptions](v_)
	})
}
func (*DependencyOptions) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *DependencyOptions) initFields_DependencyOptions() {
	_ = this_
}

// Java constructor: DependencyOptions(com.google.javascript.jscomp.DependencyOptions.DependencyMode mode, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.ModuleIdentifier> entryPoints)
func NewDependencyOptions(mode DependencyOptions_DependencyMode, entryPoints j2gort.Ext_com_google_common_collect_ImmutableList) *DependencyOptions {
	this_ := &DependencyOptions{}
	InitDependencyOptions(this_, mode, entryPoints)
	return this_
}

func InitDependencyOptions(this_ *DependencyOptions, mode DependencyOptions_DependencyMode, entryPoints j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = mode
	_ = entryPoints
	this_.initFields_DependencyOptions()
	j2gort.CheckNotNullEnum(mode)
	j2gort.CheckNotNull(entryPoints)
	this_.mode = mode
	this_.entryPoints = j2gort.ImmListToSliceNullable[

	// Java method: DependencyOptions.getMode()
	*ModuleIdentifier](entryPoints)
}

func (this_ *DependencyOptions) GetMode() DependencyOptions_DependencyMode {
	return this_.Mode()
}

// Java method: DependencyOptions.getEntryPoints()
func (this_ *DependencyOptions) GetEntryPoints() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.EntryPoints()
}

// Java method: DependencyOptions.needsManagement()
func (this_ *DependencyOptions) NeedsManagement() bool {
	return (this_.Mode() != StaticGet_DependencyOptions_DependencyMode_NONE())
}

// Java method: DependencyOptions.shouldSort()
func (this_ *DependencyOptions) ShouldSort() bool {
	return (this_.Mode() != StaticGet_DependencyOptions_DependencyMode_NONE())
}

// Java method: DependencyOptions.shouldPrune()
func (this_ *DependencyOptions) ShouldPrune() bool {
	return (((this_.Mode() == StaticGet_DependencyOptions_DependencyMode_PRUNE_LEGACY()) || (this_.Mode() == StaticGet_DependencyOptions_DependencyMode_PRUNE())) || (this_.Mode() == StaticGet_DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS()))
}

// Java method: DependencyOptions.shouldDropMoochers()
func (this_ *DependencyOptions) ShouldDropMoochers() bool {
	return ((this_.Mode() == StaticGet_DependencyOptions_DependencyMode_PRUNE()) || (this_.Mode() == StaticGet_DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS()))
}

// Java static method: DependencyOptions.none()
func DependencyOptions_None() *DependencyOptions {
	return NewDependencyOptions(StaticGet_DependencyOptions_DependencyMode_NONE(), j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
}

// Java static method: DependencyOptions.sortOnly()
func DependencyOptions_SortOnly() *DependencyOptions {
	return NewDependencyOptions(StaticGet_DependencyOptions_DependencyMode_SORT_ONLY(), j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
}

// Java static method: DependencyOptions.pruneLegacyForEntryPoints(java.lang.Iterable<com.google.javascript.jscomp.ModuleIdentifier> entryPoints)
func DependencyOptions_PruneLegacyForEntryPoints(entryPoints j2gort.Ext_java_lang_Iterable) *DependencyOptions {
	_ = entryPoints
	return NewDependencyOptions(StaticGet_DependencyOptions_DependencyMode_PRUNE_LEGACY(), j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(entryPoints))
}

// Java static method: DependencyOptions.pruneAllowNoEntryPoints(java.lang.Iterable<com.google.javascript.jscomp.ModuleIdentifier> entryPoints)
func DependencyOptions_PruneAllowNoEntryPoints(entryPoints j2gort.Ext_java_lang_Iterable) *DependencyOptions {
	_ = entryPoints
	return NewDependencyOptions(StaticGet_DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS(), j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(entryPoints))
}

// Java static method: DependencyOptions.pruneForEntryPoints(java.lang.Iterable<com.google.javascript.jscomp.ModuleIdentifier> entryPoints)
func DependencyOptions_PruneForEntryPoints(entryPoints j2gort.Ext_java_lang_Iterable) *DependencyOptions {
	_ = entryPoints
	if !(!(j2gort.Ext_com_google_common_collect_Iterables_IsEmpty(entryPoints))) {
		j2gort.CheckStateStr(false, "DependencyMode.PRUNE requires at least one entry point")
	}
	return NewDependencyOptions(StaticGet_DependencyOptions_DependencyMode_PRUNE(), j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(entryPoints))
}

// Java static method: DependencyOptions.fromFlags(com.google.javascript.jscomp.DependencyOptions.@org.jspecify.annotations.Nullable DependencyMode dependencyModeFlag, java.util.List<java.lang.String> entryPointFlag, java.util.List<java.lang.String> closureEntryPointFlag, java.lang.String commonJsEntryModuleFlag, boolean manageClosureDependenciesFlag, boolean onlyClosureDependenciesFlag)
func DependencyOptions_FromFlags(dependencyModeFlag DependencyOptions_DependencyMode, entryPointFlag j2gort.Ext_java_util_List, closureEntryPointFlag j2gort.Ext_java_util_List, commonJsEntryModuleFlag string, manageClosureDependenciesFlag bool, onlyClosureDependenciesFlag bool) *DependencyOptions {
	_ = dependencyModeFlag
	_ = entryPointFlag
	_ = closureEntryPointFlag
	_ = commonJsEntryModuleFlag
	_ = manageClosureDependenciesFlag
	_ = onlyClosureDependenciesFlag
	var hasEntryPoint bool = (((commonJsEntryModuleFlag != j2gort.NullString) || !(entryPointFlag.IsEmpty())) || !(closureEntryPointFlag.IsEmpty()))
	_ = hasEntryPoint
	if !(hasEntryPoint) && onlyClosureDependenciesFlag {
		panic(NewFlagUsageException("--only_closure_dependencies requires --entry_point."))
	}
	if !(hasEntryPoint) && (dependencyModeFlag == StaticGet_DependencyOptions_DependencyMode_PRUNE()) {
		panic(NewFlagUsageException("--dependency_mode=PRUNE requires --entry_point."))
	}
	if hasEntryPoint && ((dependencyModeFlag == StaticGet_DependencyOptions_DependencyMode_NONE()) || (dependencyModeFlag == StaticGet_DependencyOptions_DependencyMode_SORT_ONLY())) {
		panic(NewFlagUsageException(("--dependency_mode=" + j2gort.StrVal(dependencyModeFlag) + " cannot be used with --entry_point, --closure_entry_point or --common_js_entry_module.")))
	}
	if !(entryPointFlag.IsEmpty()) && !(closureEntryPointFlag.IsEmpty()) {
		panic(NewFlagUsageException("--closure_entry_point cannot be used with --entry_point."))
	}
	if (commonJsEntryModuleFlag != j2gort.NullString) && (!(entryPointFlag.IsEmpty()) || !(closureEntryPointFlag.IsEmpty())) {
		panic(NewFlagUsageException("--common_js_entry_module cannot be used with either --entry_point or --closure_entry_point."))
	}
	if manageClosureDependenciesFlag && onlyClosureDependenciesFlag {
		panic(NewFlagUsageException("--only_closure_dependencies cannot be used with --manage_closure_dependencies."))
	}
	if manageClosureDependenciesFlag && (dependencyModeFlag != 0) {
		panic(NewFlagUsageException("--manage_closure_dependencies cannot be used with --dependency_mode."))
	}
	if onlyClosureDependenciesFlag && (dependencyModeFlag != 0) {
		panic(NewFlagUsageException("--only_closure_dependencies cannot be used with --dependency_mode."))
	}
	var dependencyMode DependencyOptions_DependencyMode = DependencyOptions_DependencyMode(0)
	_ = dependencyMode
	if (dependencyModeFlag == StaticGet_DependencyOptions_DependencyMode_PRUNE()) || onlyClosureDependenciesFlag {
		dependencyMode = StaticGet_DependencyOptions_DependencyMode_PRUNE()
	} else if dependencyModeFlag == StaticGet_DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS() {
		dependencyMode = StaticGet_DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS()
	} else if ((dependencyModeFlag == StaticGet_DependencyOptions_DependencyMode_PRUNE_LEGACY()) || manageClosureDependenciesFlag) || hasEntryPoint {
		dependencyMode = StaticGet_DependencyOptions_DependencyMode_PRUNE_LEGACY()
	} else if dependencyModeFlag != 0 {
		dependencyMode = dependencyModeFlag
	} else {
		return nil
	}
	var entryPointsBuilder []*ModuleIdentifier = []*ModuleIdentifier{}
	_ = entryPointsBuilder
	if commonJsEntryModuleFlag != j2gort.NullString {
		entryPointsBuilder = append(entryPointsBuilder, j2gort.CheckNotNullElem(ModuleIdentifier_ForFile(commonJsEntryModuleFlag)))
	}
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(entryPointFlag))
		for cur_1.HasNext() {
			entryPoint := cur_1.NextStr()
			_ = entryPoint
			entryPointsBuilder = append(entryPointsBuilder, j2gort.CheckNotNullElem(ModuleIdentifier_ForFlagValue(entryPoint)))
		}
	}
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(closureEntryPointFlag))
		for cur_2.HasNext() {
			closureEntryPoint := cur_2.NextStr()
			_ = closureEntryPoint
			entryPointsBuilder = append(entryPointsBuilder, j2gort.CheckNotNullElem(ModuleIdentifier_ForClosure(closureEntryPoint)))
		}
	}
	switch (dependencyMode).Ordinal() {
	case int32(0) /* DependencyMode.NONE */ :
		return DependencyOptions_None()
	case int32(1) /* DependencyMode.SORT_ONLY */ :
		return DependencyOptions_SortOnly()
	case int32(2) /* DependencyMode.PRUNE_LEGACY */ :
		return DependencyOptions_PruneLegacyForEntryPoints(j2gort.ImmListFromTypedSlice(entryPointsBuilder))
	case int32(4) /* DependencyMode.PRUNE_ALLOW_NO_ENTRY_POINTS */ :
		return DependencyOptions_PruneAllowNoEntryPoints(j2gort.ImmListFromTypedSlice(entryPointsBuilder))
	case int32(3) /* DependencyMode.PRUNE */ :
		return DependencyOptions_PruneForEntryPoints(j2gort.ImmListFromTypedSlice(entryPointsBuilder))
	}
	panic("unreachable switch expression")
}

func (this_ *DependencyOptions) Mode() DependencyOptions_DependencyMode { return this_.mode }
func (this_ *DependencyOptions) EntryPoints() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(this_.entryPoints)
}
func (this_ *DependencyOptions) Equals(other_ any) bool {
	o_, ok_ := other_.(*DependencyOptions)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.mode != o_.mode {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.ImmListFromTypedSlice(this_.entryPoints), j2gort.ImmListFromTypedSlice(o_.entryPoints)) {
		return false
	}
	return true
}
func (this_ *DependencyOptions) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + (this_.mode).HashCode()
	h_ = 31*h_ + j2gort.ObjectHash(j2gort.ImmListFromTypedSlice(this_.entryPoints))
	return h_
}
func (this_ *DependencyOptions) String() string {
	_ = fmt.Sprint
	return "DependencyOptions[" + "mode=" + j2gort.StrVal(j2gort.BoxNullableVal(this_.mode)) + ", " + "entryPoints=" + j2gort.StrVal(j2gort.ImmListFromTypedSlice(this_.entryPoints)) + "]"
}

// Java enum: com.google.javascript.jscomp.DependencyOptions.DependencyMode
type DependencyOptions_DependencyMode uint8

const (
	enum_DependencyOptions_DependencyMode_NONE                        DependencyOptions_DependencyMode = 1
	enum_DependencyOptions_DependencyMode_SORT_ONLY                   DependencyOptions_DependencyMode = 2
	enum_DependencyOptions_DependencyMode_PRUNE_LEGACY                DependencyOptions_DependencyMode = 3
	enum_DependencyOptions_DependencyMode_PRUNE                       DependencyOptions_DependencyMode = 4
	enum_DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS DependencyOptions_DependencyMode = 5
)

var enumBaseTable_DependencyOptions_DependencyMode = [...]j2gort.EnumBase{
	{Name_: "NONE", Ordinal_: 0},
	{Name_: "SORT_ONLY", Ordinal_: 1},
	{Name_: "PRUNE_LEGACY", Ordinal_: 2},
	{Name_: "PRUNE", Ordinal_: 3},
	{Name_: "PRUNE_ALLOW_NO_ENTRY_POINTS", Ordinal_: 4},
}
var enumValues_DependencyOptions_DependencyMode = []DependencyOptions_DependencyMode{1, 2, 3, 4, 5}

func StaticGet_DependencyOptions_DependencyMode_NONE() DependencyOptions_DependencyMode {
	return enum_DependencyOptions_DependencyMode_NONE
}
func StaticGet_DependencyOptions_DependencyMode_SORT_ONLY() DependencyOptions_DependencyMode {
	return enum_DependencyOptions_DependencyMode_SORT_ONLY
}
func StaticGet_DependencyOptions_DependencyMode_PRUNE_LEGACY() DependencyOptions_DependencyMode {
	return enum_DependencyOptions_DependencyMode_PRUNE_LEGACY
}
func StaticGet_DependencyOptions_DependencyMode_PRUNE() DependencyOptions_DependencyMode {
	return enum_DependencyOptions_DependencyMode_PRUNE
}
func StaticGet_DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS() DependencyOptions_DependencyMode {
	return enum_DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS
}

func DependencyOptions_DependencyMode_Values() []DependencyOptions_DependencyMode {
	return j2gort.CloneSlice(enumValues_DependencyOptions_DependencyMode)
}

func DependencyOptions_DependencyMode_ValueOf(name string) DependencyOptions_DependencyMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_DependencyOptions_DependencyMode {
		if enumBaseTable_DependencyOptions_DependencyMode[i_].Name_ == name {
			return DependencyOptions_DependencyMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.DependencyOptions.DependencyMode." + name))
}

func (this_ DependencyOptions_DependencyMode) IsJavaNil_() bool                              { return this_ == 0 }
func (DependencyOptions_DependencyMode) HasIdentityEquals_()                                 {}
func (DependencyOptions_DependencyMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (DependencyOptions_DependencyMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ DependencyOptions_DependencyMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (DependencyOptions_DependencyMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DependencyOptions$DependencyMode", "DependencyMode"
}
func (this_ DependencyOptions_DependencyMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_DependencyOptions_DependencyMode[this_-1]
}
func (this_ DependencyOptions_DependencyMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ DependencyOptions_DependencyMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ DependencyOptions_DependencyMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_DependencyOptions_DependencyMode[this_-1].Name_
}
func (this_ DependencyOptions_DependencyMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_DependencyOptions_DependencyMode[this_-1].Name_
}
func (this_ DependencyOptions_DependencyMode) Equals(other_ any) bool {
	o_, ok_ := other_.(DependencyOptions_DependencyMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ DependencyOptions_DependencyMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_DependencyOptions_DependencyMode[this_-1].HashCode()
}
func (this_ DependencyOptions_DependencyMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_DependencyOptions_DependencyMode[this_-1].CompareTo(other_)
}
func (this_ DependencyOptions_DependencyMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_DependencyOptions_DependencyMode[this_-1].Name_)
}
func init() {
	vals_ := []any{DependencyOptions_DependencyMode(1), DependencyOptions_DependencyMode(2), DependencyOptions_DependencyMode(3), DependencyOptions_DependencyMode(4), DependencyOptions_DependencyMode(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.DependencyOptions.DependencyMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.DependencyOptions$DependencyMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DependencyOptions.DependencyMode", "com.google.javascript.jscomp.DependencyOptions$DependencyMode", func(v_ any) bool {
		return j2gort.InstanceOf[DependencyOptions_DependencyMode](v_)
	})
}
func (this_ *DependencyOptions) Slice_EntryPoints() []*ModuleIdentifier {
	return this_.entryPoints
}

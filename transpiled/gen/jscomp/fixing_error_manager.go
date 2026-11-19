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
// Port of com/google/javascript/refactoring/FixingErrorManager.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.refactoring.FixingErrorManager
type FixingErrorManager struct {
	BasicErrorManagerBase
	fixer           *ErrorToFixMapper                                 `gson:"fixer"`
	sureFixes       j2gort.CustomOrderedMap[*JSError, *SuggestedFix]  `gson:"sureFixes"`
	multiFixes      j2gort.Ext_com_google_common_collect_ListMultimap `gson:"multiFixes"`
	unfixableErrors j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"unfixableErrors"`
}

func (this_ *FixingErrorManager) BindSelf_(s any) {
	_ = s
	this_.BasicErrorManagerBase.BindSelf_(s)
}
func (this_ *FixingErrorManager) IsJavaNil_() bool { return this_ == nil }
func (*FixingErrorManager) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.FixingErrorManager", "FixingErrorManager"
}
func (this_ *FixingErrorManager) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.BasicErrorManagerBase.IdentityHash_()
}
func (*FixingErrorManager) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.FixingErrorManager", "com.google.javascript.refactoring.FixingErrorManager", func(v_ any) bool {
		return j2gort.InstanceOf[*FixingErrorManager](v_)
	})
}
func (*FixingErrorManager) Implements_jscomp_BasicErrorManager()   {}
func (*FixingErrorManager) Implements_jscomp_SortingErrorManager() {}
func (*FixingErrorManager) Implements_jscomp_ErrorManager()        {}
func (*FixingErrorManager) Implements_jscomp_ErrorHandler()        {}

func (this_ *FixingErrorManager) initFields_FixingErrorManager() {
	_ = this_
	this_.sureFixes = j2gort.NewCustomOrderedMap[*JSError, *SuggestedFix]()
	this_.multiFixes = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ListMultimap](j2gort.Ext_com_google_common_collect_ArrayListMultimap_Create_0())
}

// Java constructor: FixingErrorManager()
func NewFixingErrorManager() *FixingErrorManager {
	this_ := &FixingErrorManager{}
	this_.BindSelf_(this_)
	InitFixingErrorManager(this_)
	return this_
}

func InitFixingErrorManager(this_ *FixingErrorManager) {
	InitFixingErrorManager_1(this_, j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
}

// Java constructor: FixingErrorManager(com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.DiagnosticType> unfixableErrors)
func NewFixingErrorManager_1(unfixableErrors j2gort.Ext_com_google_common_collect_ImmutableSet) *FixingErrorManager {
	this_ := &FixingErrorManager{}
	this_.BindSelf_(this_)
	InitFixingErrorManager_1(this_, unfixableErrors)
	return this_
}

func InitFixingErrorManager_1(this_ *FixingErrorManager, unfixableErrors j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = unfixableErrors
	InitBasicErrorManager(&this_.BasicErrorManagerBase)
	this_.initFields_FixingErrorManager()
	this_.unfixableErrors = unfixableErrors
}

// Java method: FixingErrorManager.setCompiler(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *FixingErrorManager) SetCompiler(compiler AbstractCompiler) {
	_ = compiler
	this_.fixer = NewErrorToFixMapper(compiler)
}

// Java method: FixingErrorManager.report(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *FixingErrorManager) Report_2(level *CheckLevel, error_ *JSError) {
	_ = level
	_ = error_
	this_.BasicErrorManagerBase.Report_2(level, error_)
	if !(this_.unfixableErrors.Contains(j2gort.PtrToAny(error_.Type()))) {
		var fixes []*SuggestedFix = j2gort.ImmListToSliceNullable[*SuggestedFix](this_.fixer.GetFixesForJsError(error_))
		_ = fixes
		if int32(len(fixes)) == int32(1) {
			this_.sureFixes.Put(error_, fixes[int32(0)])
		} else {
			this_.multiFixes.PutAll_2(j2gort.PtrToAny(error_), j2gort.ImmListFromTypedSlice(fixes))
		}
	}
}

// Java method: FixingErrorManager.containsFixableShorthandModuleWarning()
func (this_ *FixingErrorManager) priv_FixingErrorManager_containsFixableShorthandModuleWarning_7897() bool {
	{
		cur_1 := this_.sureFixes.IterKeys()
		for cur_1.HasNext() {

			error_ := cur_1.
				Next()
			_ = error_
			if error_.Type().Equals(j2gort.PtrToAny(StaticGet_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME())) {
				return true
			}
		}
	}
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.multiFixes.KeySet())))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			error__2 := j2gort.CastPtr[JSError](rawElem_4)
			_ = error__2
			if error__2.Type().Equals(j2gort.PtrToAny(StaticGet_ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME())) {
				return true
			}
		}
	}
	return false
}

// Java method: FixingErrorManager.getFixesForJsError(com.google.javascript.jscomp.JSError error)
func (this_ *FixingErrorManager) GetFixesForJsError(error_ *JSError) j2gort.Ext_java_util_List {
	_ = error_
	if this_.sureFixes.ContainsKey(error_) {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(this_.sureFixes.Get(error_)))
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_List](this_.multiFixes.Get(j2gort.PtrToAny(error_)))
	}
}

// Java method: FixingErrorManager.getFixes(com.google.javascript.refactoring.FixingErrorManager.FixTypes fixTypes)
func (this_ *FixingErrorManager) priv_FixingErrorManager_getFixes_7899(fixTypes FixingErrorManager_FixTypes) j2gort.Ext_java_util_Collection {
	_ = fixTypes
	var containsFixableShorthandModuleWarning bool = this_.priv_FixingErrorManager_containsFixableShorthandModuleWarning_7897()
	_ = containsFixableShorthandModuleWarning
	var fixes []*SuggestedFix = []*SuggestedFix{}
	_ = fixes
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetErrors()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			error_ := j2gort.CastPtr[JSError](rawElem_6)
			_ = error_
			if containsFixableShorthandModuleWarning && error_.Type().Equals(j2gort.PtrToAny(StaticGet_CheckExtraRequires_EXTRA_REQUIRE_WARNING())) {
			} else {
				if (fixTypes == StaticGet_FixingErrorManager_FixTypes_ONE_FIX()) && this_.sureFixes.ContainsKey(error_) {
					fixes = append(fixes, this_.sureFixes.Get(error_))
				} else if fixTypes == StaticGet_FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES() {
					fixes = j2gort.AppendCollection[*SuggestedFix](fixes, j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.GetFixesForJsError(error_)))
				}
			}
		}
	}
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetWarnings()))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			warning := j2gort.CastPtr[JSError](rawElem_8)
			_ = warning
			if warning.Type().Equals(j2gort.PtrToAny(StaticGet_CheckExtraRequires_EXTRA_REQUIRE_WARNING())) && containsFixableShorthandModuleWarning {
			} else {
				if (fixTypes == StaticGet_FixingErrorManager_FixTypes_ONE_FIX()) && this_.sureFixes.ContainsKey(warning) {
					fixes = append(fixes, this_.sureFixes.Get(warning))
				} else if fixTypes == StaticGet_FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES() {
					fixes = j2gort.AppendCollection[*SuggestedFix](fixes, j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.GetFixesForJsError(warning)))
				}
			}
		}
	}
	return j2gort.WrapTypedArrayList(fixes)
}

// Java method: FixingErrorManager.getSureFixes()
func (this_ *FixingErrorManager) GetSureFixes() j2gort.Ext_java_util_Collection {
	return this_.priv_FixingErrorManager_getFixes_7899(StaticGet_FixingErrorManager_FixTypes_ONE_FIX())
}

// Java method: FixingErrorManager.getAllFixes()
func (this_ *FixingErrorManager) GetAllFixes() j2gort.Ext_java_util_Collection {
	return this_.priv_FixingErrorManager_getFixes_7899(StaticGet_FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES())
}

// Java method: FixingErrorManager.printSummary()
func (this_ *FixingErrorManager) PrintSummary() {
}

// Java method: FixingErrorManager.println(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *FixingErrorManager) Println(level *CheckLevel, error_ *JSError) {
	_ = level
	_ = error_
}

// Java enum: com.google.javascript.refactoring.FixingErrorManager.FixTypes
type FixingErrorManager_FixTypes uint8

const (
	enum_FixingErrorManager_FixTypes_ONE_FIX                 FixingErrorManager_FixTypes = 1
	enum_FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES FixingErrorManager_FixTypes = 2
)

var enumBaseTable_FixingErrorManager_FixTypes = [...]j2gort.EnumBase{
	{Name_: "ONE_FIX", Ordinal_: 0},
	{Name_: "ONE_FIX_AND_MULTI_FIXES", Ordinal_: 1},
}
var enumValues_FixingErrorManager_FixTypes = []FixingErrorManager_FixTypes{1, 2}

func StaticGet_FixingErrorManager_FixTypes_ONE_FIX() FixingErrorManager_FixTypes {
	return enum_FixingErrorManager_FixTypes_ONE_FIX
}
func StaticGet_FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES() FixingErrorManager_FixTypes {
	return enum_FixingErrorManager_FixTypes_ONE_FIX_AND_MULTI_FIXES
}

func FixingErrorManager_FixTypes_Values() []FixingErrorManager_FixTypes {
	return j2gort.CloneSlice(enumValues_FixingErrorManager_FixTypes)
}

func FixingErrorManager_FixTypes_ValueOf(name string) FixingErrorManager_FixTypes {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_FixingErrorManager_FixTypes {
		if enumBaseTable_FixingErrorManager_FixTypes[i_].Name_ == name {
			return FixingErrorManager_FixTypes(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.refactoring.FixingErrorManager.FixTypes." + name))
}

func (this_ FixingErrorManager_FixTypes) IsJavaNil_() bool                              { return this_ == 0 }
func (FixingErrorManager_FixTypes) HasIdentityEquals_()                                 {}
func (FixingErrorManager_FixTypes) Implements_j2gort_Ext_java_lang_Enum()               {}
func (FixingErrorManager_FixTypes) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ FixingErrorManager_FixTypes) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (FixingErrorManager_FixTypes) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.FixingErrorManager$FixTypes", "FixTypes"
}
func (this_ FixingErrorManager_FixTypes) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_FixingErrorManager_FixTypes[this_-1]
}
func (this_ FixingErrorManager_FixTypes) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ FixingErrorManager_FixTypes) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ FixingErrorManager_FixTypes) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FixingErrorManager_FixTypes[this_-1].Name_
}
func (this_ FixingErrorManager_FixTypes) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_FixingErrorManager_FixTypes[this_-1].Name_
}
func (this_ FixingErrorManager_FixTypes) Equals(other_ any) bool {
	o_, ok_ := other_.(FixingErrorManager_FixTypes)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ FixingErrorManager_FixTypes) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_FixingErrorManager_FixTypes[this_-1].HashCode()
}
func (this_ FixingErrorManager_FixTypes) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FixingErrorManager_FixTypes[this_-1].CompareTo(other_)
}
func (this_ FixingErrorManager_FixTypes) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_FixingErrorManager_FixTypes[this_-1].Name_)
}
func init() {
	vals_ := []any{FixingErrorManager_FixTypes(1), FixingErrorManager_FixTypes(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.refactoring.FixingErrorManager.FixTypes", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.refactoring.FixingErrorManager$FixTypes", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.FixingErrorManager.FixTypes", "com.google.javascript.refactoring.FixingErrorManager$FixTypes", func(v_ any) bool {
		return j2gort.InstanceOf[FixingErrorManager_FixTypes](v_)
	})
}

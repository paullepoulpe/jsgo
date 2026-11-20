/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/SortingErrorManager.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.SortingErrorManager
type SortingErrorManager interface {
	Implements_jscomp_SortingErrorManager()
	IsJavaNil_() bool
	GetSortingErrorManagerBase() *SortingErrorManagerBase
	ErrorManager
	GetSortedDiagnostics() j2gort.Ext_java_lang_Iterable
}

func ToIface_SortingErrorManager[T interface {
	SortingErrorManager
	comparable
}](v T) SortingErrorManager {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SortingErrorManager[T SortingErrorManager](v T) SortingErrorManager {
	var vi SortingErrorManager = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.SortingErrorManager
type SortingErrorManagerBase struct {
	self_                 SortingErrorManager
	messages              *j2gort.Ext_java_util_TreeSet                     `gson:"messages"`
	errorReportGenerators j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"errorReportGenerators"`
	typedPercent          float64                                           `gson:"typedPercent"`
	originalErrorCount    int32                                             `gson:"originalErrorCount"`
	promotedErrorCount    int32                                             `gson:"promotedErrorCount"`
	warningCount          int32                                             `gson:"warningCount"`
	idHash_               atomic.Int32
}

func (this_ *SortingErrorManagerBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(SortingErrorManager)
}
func (this_ *SortingErrorManagerBase) IsJavaNil_() bool { return this_ == nil }
func (*SortingErrorManagerBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SortingErrorManager", "SortingErrorManager"
}
func (this_ *SortingErrorManagerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SortingErrorManager", "com.google.javascript.jscomp.SortingErrorManager", func(v_ any) bool {
		return j2gort.InstanceOf[SortingErrorManager](v_)
	})
}
func (b *SortingErrorManagerBase) GetSortingErrorManagerBase() *SortingErrorManagerBase { return b }
func (*SortingErrorManagerBase) Implements_jscomp_SortingErrorManager()                 {}
func (*SortingErrorManagerBase) Implements_jscomp_ErrorManager()                        {}
func (*SortingErrorManagerBase) Implements_jscomp_ErrorHandler()                        {}

func (this_ *SortingErrorManagerBase) initFields_SortingErrorManager() {
	_ = this_
	this_.messages = j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Comparator(NewSortingErrorManager_LeveledJSErrorComparator())
	this_.originalErrorCount = int32(0)
	this_.promotedErrorCount = int32(0)
	this_.warningCount = int32(0)
	this_.typedPercent = float64(0.0)
}

// Java constructor: SortingErrorManager(java.util.Set<com.google.javascript.jscomp.SortingErrorManager.ErrorReportGenerator> errorReportGenerators)
func NewSortingErrorManager(errorReportGenerators j2gort.Ext_java_util_Set) SortingErrorManager {
	this_ := &SortingErrorManagerBase{}
	this_.BindSelf_(this_)
	InitSortingErrorManager(this_, errorReportGenerators)
	return this_
}

func InitSortingErrorManager(this_ *SortingErrorManagerBase, errorReportGenerators j2gort.Ext_java_util_Set) {
	_ = errorReportGenerators
	this_.initFields_SortingErrorManager()
	this_.errorReportGenerators = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](errorReportGenerators))
}

// Java method: SortingErrorManager.report(com.google.javascript.jscomp.CheckLevel level, com.google.javascript.jscomp.JSError error)
func (this_ *SortingErrorManagerBase) Report_2(level *CheckLevel, error_ *JSError) {
	_ = level
	_ = error_
	var val_1 SortingErrorManager_ErrorWithLevel
	InitSortingErrorManager_ErrorWithLevel(&val_1, error_, level)
	var e *SortingErrorManager_ErrorWithLevel = &val_1
	_ = e
	if this_.messages.Add_1(e) {
		if level == StaticGet_CheckLevel_ERROR() {
			if error_.Type().level == StaticGet_CheckLevel_ERROR() {
				this_.originalErrorCount = this_.originalErrorCount + 1
			} else {
				this_.promotedErrorCount = this_.promotedErrorCount + 1
			}
		} else if level == StaticGet_CheckLevel_WARNING() {
			this_.warningCount = this_.warningCount + 1
		}
	}
}

// Java method: SortingErrorManager.hasHaltingErrors()
func (this_ *SortingErrorManagerBase) HasHaltingErrors() bool {
	return (this_.originalErrorCount != int32(0))
}

// Java method: SortingErrorManager.getErrorCount()
func (this_ *SortingErrorManagerBase) GetErrorCount() int32 {
	return (this_.originalErrorCount + this_.promotedErrorCount)
}

// Java method: SortingErrorManager.getWarningCount()
func (this_ *SortingErrorManagerBase) GetWarningCount() int32 {
	return this_.warningCount
}

// Java method: SortingErrorManager.getErrors()
func (this_ *SortingErrorManagerBase) GetErrors() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.priv_SortingErrorManager_filterToListOf_5128(StaticGet_CheckLevel_ERROR())
}

// Java method: SortingErrorManager.getWarnings()
func (this_ *SortingErrorManagerBase) GetWarnings() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.priv_SortingErrorManager_filterToListOf_5128(StaticGet_CheckLevel_WARNING())
}

// Java method: SortingErrorManager.getSortedDiagnostics()
func (this_ *SortingErrorManagerBase) GetSortedDiagnostics() j2gort.Ext_java_lang_Iterable {
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.messages)))
}

// Java method: SortingErrorManager.setTypedPercent(double typedPercent)
func (this_ *SortingErrorManagerBase) SetTypedPercent(typedPercent float64) {
	_ = typedPercent
	this_.typedPercent = typedPercent
}

// Java method: SortingErrorManager.getTypedPercent()
func (this_ *SortingErrorManagerBase) GetTypedPercent() float64 {
	return this_.typedPercent
}

// Java method: SortingErrorManager.filterToListOf(com.google.javascript.jscomp.CheckLevel level)
func (this_ *SortingErrorManagerBase) priv_SortingErrorManager_filterToListOf_5128(level *CheckLevel) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = level
	var list []*JSError = []*JSError{}
	_ = list
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(this_.messages))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			e := j2gort.CastPtr[SortingErrorManager_ErrorWithLevel](rawElem_3)
			_ = e
			if e.level == level {
				list = append(list, e.error_)
			}
		}
	}
	return j2gort.ImmListFromSlice(list)
}

// Java method: SortingErrorManager.generateReport()
func (this_ *SortingErrorManagerBase) GenerateReport() {
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.errorReportGenerators))
		for cur_4.HasNext() {
			rawElem_5 := cur_4.Next()
			generator := j2gort.AsIface[SortingErrorManager_ErrorReportGenerator](rawElem_5)
			_ = generator
			generator.GenerateReport(this_.self_)
		}
	}
}

func (this_ *SortingErrorManagerBase) ShouldReportConformanceViolation(p0_ *j2gort.Ext_com_google_javascript_jscomp_Requirement, p1_ j2gort.Ext_com_google_common_base_Optional, p2_ *JSError, p3_ *j2gort.Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior, p4_ bool) bool {
	return IfaceDefault_ErrorManager_ShouldReportConformanceViolation(this_.self_, p0_, p1_, p2_, p3_, p4_)
}

// Java interface: com.google.javascript.jscomp.SortingErrorManager.ErrorReportGenerator
type SortingErrorManager_ErrorReportGenerator interface {
	Implements_jscomp_SortingErrorManager_ErrorReportGenerator()
	IsJavaNil_() bool
	GenerateReport(SortingErrorManager)
}

func ToIface_SortingErrorManager_ErrorReportGenerator[T interface {
	SortingErrorManager_ErrorReportGenerator
	comparable
}](v T) SortingErrorManager_ErrorReportGenerator {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SortingErrorManager_ErrorReportGenerator[T SortingErrorManager_ErrorReportGenerator](v T) SortingErrorManager_ErrorReportGenerator {
	var vi SortingErrorManager_ErrorReportGenerator = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SortingErrorManager.ErrorReportGenerator", "com.google.javascript.jscomp.SortingErrorManager$ErrorReportGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[SortingErrorManager_ErrorReportGenerator](v_)
	})
}

// Java class: com.google.javascript.jscomp.SortingErrorManager.LeveledJSErrorComparator
type SortingErrorManager_LeveledJSErrorComparator struct {
	idHash_ atomic.Int32
}

func (this_ *SortingErrorManager_LeveledJSErrorComparator) BindSelf_(s any) {
	_ = s
}
func (this_ *SortingErrorManager_LeveledJSErrorComparator) IsJavaNil_() bool { return this_ == nil }
func (*SortingErrorManager_LeveledJSErrorComparator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SortingErrorManager$LeveledJSErrorComparator", "LeveledJSErrorComparator"
}
func (this_ *SortingErrorManager_LeveledJSErrorComparator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SortingErrorManager_LeveledJSErrorComparator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SortingErrorManager.LeveledJSErrorComparator", "com.google.javascript.jscomp.SortingErrorManager$LeveledJSErrorComparator", func(v_ any) bool {
		return j2gort.InstanceOf[*SortingErrorManager_LeveledJSErrorComparator](v_)
	})
}
func (*SortingErrorManager_LeveledJSErrorComparator) Implements_j2gort_Ext_java_util_Comparator() {}

func StaticGet_SortingErrorManager_LeveledJSErrorComparator_P1_LT_P2() int32 { return int32(-1) }
func StaticGet_SortingErrorManager_LeveledJSErrorComparator_P1_GT_P2() int32 { return int32(1) }
func (this_ *SortingErrorManager_LeveledJSErrorComparator) initFields_SortingErrorManager_LeveledJSErrorComparator() {
	_ = this_
}

// Java constructor: LeveledJSErrorComparator()
func NewSortingErrorManager_LeveledJSErrorComparator() *SortingErrorManager_LeveledJSErrorComparator {
	this_ := &SortingErrorManager_LeveledJSErrorComparator{}
	InitSortingErrorManager_LeveledJSErrorComparator(this_)
	return this_
}

func InitSortingErrorManager_LeveledJSErrorComparator(this_ *SortingErrorManager_LeveledJSErrorComparator) {
	this_.initFields_SortingErrorManager_LeveledJSErrorComparator()
}

// Java method: LeveledJSErrorComparator.compare(com.google.javascript.jscomp.SortingErrorManager.ErrorWithLevel p1, com.google.javascript.jscomp.SortingErrorManager.ErrorWithLevel p2)
func (this_ *SortingErrorManager_LeveledJSErrorComparator) Compare(arg_6 any, arg_7 any) int32 {
	p1 := j2gort.CastPtr[SortingErrorManager_ErrorWithLevel](arg_6)
	_ = p1
	p2 := j2gort.CastPtr[SortingErrorManager_ErrorWithLevel](arg_7)
	_ = p2
	if p2 == nil {
		if p1 == nil {
			return int32(0)
		} else {
			return int32(1)
		}
	}
	if p1.level != p2.level {
		return j2gort.EnumCompareTo(p2.level, p1.level)
	}
	var source1 string = p1.error_.SourceName()
	_ = source1
	var source2 string = p2.error_.SourceName()
	_ = source2
	if (source1 != j2gort.NullString) && (source2 != j2gort.NullString) {
		var sourceCompare int32 = j2gort.JStrCompareTo(source1, source2)
		_ = sourceCompare
		if sourceCompare != int32(0) {
			return sourceCompare
		}
	} else if (source1 == j2gort.NullString) && (source2 != j2gort.NullString) {
		return int32(-1)
	} else if (source1 != j2gort.NullString) && (source2 == j2gort.NullString) {
		return int32(1)
	}
	var lineno1 int32 = p1.error_.GetLineNumber()
	_ = lineno1
	var lineno2 int32 = p2.error_.GetLineNumber()
	_ = lineno2
	if lineno1 != lineno2 {
		return (lineno1 - lineno2)
	} else if (lineno1 < int32(0)) && (int32(0) <= lineno2) {
		return int32(-1)
	} else if (int32(0) <= lineno1) && (lineno2 < int32(0)) {
		return int32(1)
	}
	var charno1 int32 = p1.error_.Charno()
	_ = charno1
	var charno2 int32 = p2.error_.Charno()
	_ = charno2
	if charno1 != charno2 {
		return (charno1 - charno2)
	} else if (charno1 < int32(0)) && (int32(0) <= charno2) {
		return int32(-1)
	} else if (int32(0) <= charno1) && (charno2 < int32(0)) {
		return int32(1)
	}
	return j2gort.JStrCompareTo(p1.error_.Description(), p2.error_.Description())
}

func (this_ *SortingErrorManager_LeveledJSErrorComparator) Reversed() j2gort.Ext_java_util_Comparator {
	return j2gort.ExtDefault_Comparator_Reversed(this_)
}

func (this_ *SortingErrorManager_LeveledJSErrorComparator) ThenComparing_1_Ext_java_util_Comparator(p0_ j2gort.Ext_java_util_Comparator) j2gort.Ext_java_util_Comparator {
	return j2gort.ExtDefault_Comparator_ThenComparing(this_, p0_)
}

func (this_ *SortingErrorManager_LeveledJSErrorComparator) ThenComparing_1_Ext_java_util_function_Function(p0_ j2gort.Ext_java_util_function_Function) j2gort.Ext_java_util_Comparator {
	return j2gort.ExtDefault_Comparator_ThenComparingKey(this_, p0_)
}

func (this_ *SortingErrorManager_LeveledJSErrorComparator) ThenComparingInt(p0_ j2gort.Ext_java_util_function_ToIntFunction) j2gort.Ext_java_util_Comparator {
	return j2gort.ExtDefault_Comparator_ThenComparingInt(this_, p0_)
}

func (this_ *SortingErrorManager_LeveledJSErrorComparator) ThenComparingLong(p0_ j2gort.Ext_java_util_function_ToLongFunction) j2gort.Ext_java_util_Comparator {
	return j2gort.ExtDefault_Comparator_ThenComparingLong(this_, p0_)
}

// Java class: com.google.javascript.jscomp.SortingErrorManager.ErrorWithLevel
type SortingErrorManager_ErrorWithLevel struct {
	error_ *JSError    `gson:"error"`
	level  *CheckLevel `gson:"level"`
}

func (this_ *SortingErrorManager_ErrorWithLevel) BindSelf_(s any) {
	_ = s
}
func (this_ *SortingErrorManager_ErrorWithLevel) IsJavaNil_() bool { return this_ == nil }
func (*SortingErrorManager_ErrorWithLevel) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SortingErrorManager$ErrorWithLevel", "ErrorWithLevel"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SortingErrorManager.ErrorWithLevel", "com.google.javascript.jscomp.SortingErrorManager$ErrorWithLevel", func(v_ any) bool {
		return j2gort.InstanceOf[*SortingErrorManager_ErrorWithLevel](v_)
	})
}

func (this_ *SortingErrorManager_ErrorWithLevel) initFields_SortingErrorManager_ErrorWithLevel() {
	_ = this_
}

// Java constructor: ErrorWithLevel(com.google.javascript.jscomp.JSError error, com.google.javascript.jscomp.CheckLevel level)
func NewSortingErrorManager_ErrorWithLevel(error_ *JSError, level *CheckLevel) *SortingErrorManager_ErrorWithLevel {
	this_ := &SortingErrorManager_ErrorWithLevel{}
	InitSortingErrorManager_ErrorWithLevel(this_, error_, level)
	return this_
}

func InitSortingErrorManager_ErrorWithLevel(this_ *SortingErrorManager_ErrorWithLevel, error_ *JSError, level *CheckLevel) {
	_ = error_
	_ = level
	this_.initFields_SortingErrorManager_ErrorWithLevel()
	this_.error_ = error_
	this_.level = level
}

// Java method: ErrorWithLevel.hashCode()
func (this_ *SortingErrorManager_ErrorWithLevel) HashCode() int32 {
	{
		var oh_8 int32 = 1
		{
			hp_9 := this_.level
			oh_8 = 31 * oh_8
			if hp_9 != nil {
				oh_8 += hp_9.HashCode()
			}
		}
		{
			hs_10 := this_.error_.Description()
			oh_8 = 31 * oh_8
			if hs_10 != j2gort.NullString {
				oh_8 += javalang.StringHashCode(hs_10)
			}
		}
		{
			hs_11 := this_.error_.SourceName()
			oh_8 = 31 * oh_8
			if hs_11 != j2gort.NullString {
				oh_8 += javalang.StringHashCode(hs_11)
			}
		}
		oh_8 = 31*oh_8 + (this_.error_.GetLineNumber())
		oh_8 = 31*oh_8 + (this_.error_.Charno())
		return oh_8
	}
}

// Java method: ErrorWithLevel.equals(java.lang.Object obj)
func (this_ *SortingErrorManager_ErrorWithLevel) Equals(arg_12 any) bool {
	obj := arg_12
	_ = obj
	var e *SortingErrorManager_ErrorWithLevel
	_ = e
	var ok_13 bool
	e, ok_13 = j2gort.TryCast[*SortingErrorManager_ErrorWithLevel](obj)
	if !ok_13 {
		return false
	}
	return (((((this_.level == e.level) && (this_.error_.Description() == e.error_.Description())) && (this_.error_.SourceName() == e.error_.SourceName())) && (this_.error_.GetLineNumber() == e.error_.GetLineNumber())) && (this_.error_.Charno() == e.error_.Charno()))
}

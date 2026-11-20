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
// Port of com/google/javascript/jscomp/DiagnosticToSuppressionMapper.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_diagnostic_to_suppression_mapper_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_diagnostic_to_suppression_mapper_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_diagnostic_to_suppression_mapper_3 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_diagnostic_to_suppression_mapper_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_5 any) any {
		entrySet := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_5)
		_ = entrySet
		return j2gort.BoxString(j2gort.CastPtr[DiagnosticType](entrySet.GetKey()).key)
	})
	lambdaSingleton_diagnostic_to_suppression_mapper_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](mrp0_).GetValue() })
	lambdaSingleton_diagnostic_to_suppression_mapper_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_7 any) any {
		e := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_7)
		_ = e
		return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("|%s|%s|", []any{j2gort.BoxString(j2gort.CastString(e.GetKey())), j2gort.BoxString(j2gort.CastString(e.GetValue()))}))
	})
}

// Java class: com.google.javascript.jscomp.DiagnosticToSuppressionMapper
type DiagnosticToSuppressionMapper struct {
	validSuppressions j2gort.Ext_java_util_Set                          `gson:"validSuppressions"`
	diagnosticGroups  j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"diagnosticGroups"`
	idHash_           atomic.Int32
}

func (this_ *DiagnosticToSuppressionMapper) BindSelf_(s any) {
	_ = s
}
func (this_ *DiagnosticToSuppressionMapper) IsJavaNil_() bool { return this_ == nil }
func (*DiagnosticToSuppressionMapper) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DiagnosticToSuppressionMapper", "DiagnosticToSuppressionMapper"
}
func (this_ *DiagnosticToSuppressionMapper) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DiagnosticToSuppressionMapper) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DiagnosticToSuppressionMapper", "com.google.javascript.jscomp.DiagnosticToSuppressionMapper", func(v_ any) bool {
		return j2gort.InstanceOf[*DiagnosticToSuppressionMapper](v_)
	})
}

func (this_ *DiagnosticToSuppressionMapper) initFields_DiagnosticToSuppressionMapper() {
	_ = this_
}

// Java constructor: DiagnosticToSuppressionMapper(java.util.Set<java.lang.String> validSuppressions, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.DiagnosticGroup> diagnosticGroups)
func NewDiagnosticToSuppressionMapper(validSuppressions j2gort.Ext_java_util_Set, diagnosticGroups j2gort.Ext_com_google_common_collect_ImmutableMap) *DiagnosticToSuppressionMapper {
	this_ := &DiagnosticToSuppressionMapper{}
	InitDiagnosticToSuppressionMapper(this_, validSuppressions, diagnosticGroups)
	return this_
}

func InitDiagnosticToSuppressionMapper(this_ *DiagnosticToSuppressionMapper, validSuppressions j2gort.Ext_java_util_Set, diagnosticGroups j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = validSuppressions
	_ = diagnosticGroups
	this_.initFields_DiagnosticToSuppressionMapper()
	this_.validSuppressions = validSuppressions
	this_.diagnosticGroups = diagnosticGroups
}

// Java method: DiagnosticToSuppressionMapper.printSuppressions(com.google.javascript.jscomp.DiagnosticToSuppressionMapper.OutputFormat output)
func (this_ *DiagnosticToSuppressionMapper) PrintSuppressions(output DiagnosticToSuppressionMapper_OutputFormat) {
	_ = output
	switch (output).Ordinal() {
	case int32(1) /* OutputFormat.JSON */ :
		DiagnosticToSuppressionMapper_PrintAsJson(this_.CreateSuppressionMap())
	case int32(0) /* OutputFormat.MD */ :
		DiagnosticToSuppressionMapper_PrintAsMarkdown(this_.CreateSuppressionMap())
	}
}

// Java method: DiagnosticToSuppressionMapper.createSuppressionMap()
func (this_ *DiagnosticToSuppressionMapper) CreateSuppressionMap() *j2gort.Ext_com_google_common_collect_ImmutableSortedMap {
	var validSuppressions j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_HashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.validSuppressions))
	_ = validSuppressions
	j2gort.SetRemoveStr(validSuppressions, "accessControls")
	j2gort.SetRemoveStr(validSuppressions, "missingSourcesWarnings")
	var diagnosticToSuppression *j2gort.Ext_java_util_LinkedHashMap = j2gort.NewExt_java_util_LinkedHashMap()
	_ = diagnosticToSuppression
	{
		cur_1 := j2gort.NewMapEntryCursor(this_.diagnosticGroups)
		for cur_1.NextEntry() {
			var suppression string = cur_1.KeyStr()
			_ = suppression
			if !(j2gort.SetContainsStr(validSuppressions, suppression)) {
				continue
			}
			var group *DiagnosticGroup = j2gort.CastPtr[DiagnosticGroup](cur_1.GetValue())
			_ = group
			{
				cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(group.GetTypes()))
				for cur_2.HasNext() {
					rawElem_3 := cur_2.Next()
					diagnosticType := j2gort.CastPtr[DiagnosticType](rawElem_3)
					_ = diagnosticType
					diagnosticToSuppression.Merge(j2gort.PtrToAny(diagnosticType), j2gort.BoxString(suppression), func() j2gort.Ext_java_util_function_BiFunction {
						bound_4 := this_
						return j2gort.NewLambda_Ext_java_util_function_BiFunction(func(mrp0_ any, mrp1_ any) any {
							return j2gort.BoxString(bound_4.priv_DiagnosticToSuppressionMapper_preferMoreTargetedGroup_2456(j2gort.CastString(mrp0_), j2gort.CastString(mrp1_)))
						})
					}())
				}
			}
		}
	}
	return j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSortedMap](j2gort.AsIface[j2gort.Ext_java_util_Set](diagnosticToSuppression.EntrySet()).Stream().Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSortedMap_ToImmutableSortedMap_3(j2gort.Ext_java_util_Comparator_NaturalOrder(), lambdaSingleton_diagnostic_to_suppression_mapper_1, lambdaSingleton_diagnostic_to_suppression_mapper_2)))
}

// Java method: DiagnosticToSuppressionMapper.preferMoreTargetedGroup(java.lang.String next, java.lang.String current)
func (this_ *DiagnosticToSuppressionMapper) priv_DiagnosticToSuppressionMapper_preferMoreTargetedGroup_2456(next string, current string) string {
	_ = next
	_ = current
	var nextGroup *DiagnosticGroup = j2gort.CastPtr[DiagnosticGroup](j2gort.MapGetStr(this_.diagnosticGroups, next))
	_ = nextGroup
	var currentGroup *DiagnosticGroup = j2gort.CastPtr[DiagnosticGroup](j2gort.MapGetStr(this_.diagnosticGroups, current))
	_ = currentGroup
	if nextGroup.GetTypes().Size() < currentGroup.GetTypes().Size() {
		return next
	} else {
		return current
	}
}

// Java static method: DiagnosticToSuppressionMapper.printAsJson(com.google.common.collect.ImmutableSortedMap<java.lang.String,java.lang.String> diagnosticToSuppression)
func DiagnosticToSuppressionMapper_PrintAsJson(diagnosticToSuppression *j2gort.Ext_com_google_common_collect_ImmutableSortedMap) {
	_ = diagnosticToSuppression
	j2gort.Ext_java_lang_System_F_Out.Println_1_string(j2gort.NewExt_com_google_gson_Gson().ToJson_1_any(j2gort.PtrToAny(diagnosticToSuppression)))
}

// Java static method: DiagnosticToSuppressionMapper.printAsMarkdown(com.google.common.collect.ImmutableSortedMap<java.lang.String,java.lang.String> diagnosticToSuppression)
func DiagnosticToSuppressionMapper_PrintAsMarkdown(diagnosticToSuppression *j2gort.Ext_com_google_common_collect_ImmutableSortedMap) {
	_ = diagnosticToSuppression
	j2gort.Ext_java_lang_System_F_Out.Println_1_string("| Error | Suppression tag |")
	j2gort.Ext_java_lang_System_F_Out.Println_1_string("|---|---|")
	j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](diagnosticToSuppression.EntrySet()).Stream().Map(lambdaSingleton_diagnostic_to_suppression_mapper_3).ForEachOrdered(func() j2gort.Ext_java_util_function_Consumer {
		bound_6 := j2gort.Ext_java_lang_System_F_Out
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) { bound_6.Println_1_string(j2gort.CastString(mrp0_)) })
	}())
}

// Java enum: com.google.javascript.jscomp.DiagnosticToSuppressionMapper.OutputFormat
type DiagnosticToSuppressionMapper_OutputFormat uint8

const (
	enum_DiagnosticToSuppressionMapper_OutputFormat_MD   DiagnosticToSuppressionMapper_OutputFormat = 1
	enum_DiagnosticToSuppressionMapper_OutputFormat_JSON DiagnosticToSuppressionMapper_OutputFormat = 2
)

var enumBaseTable_DiagnosticToSuppressionMapper_OutputFormat = [...]j2gort.EnumBase{
	{Name_: "MD", Ordinal_: 0},
	{Name_: "JSON", Ordinal_: 1},
}
var enumValues_DiagnosticToSuppressionMapper_OutputFormat = []DiagnosticToSuppressionMapper_OutputFormat{1, 2}

func StaticGet_DiagnosticToSuppressionMapper_OutputFormat_MD() DiagnosticToSuppressionMapper_OutputFormat {
	return enum_DiagnosticToSuppressionMapper_OutputFormat_MD
}
func StaticGet_DiagnosticToSuppressionMapper_OutputFormat_JSON() DiagnosticToSuppressionMapper_OutputFormat {
	return enum_DiagnosticToSuppressionMapper_OutputFormat_JSON
}

func DiagnosticToSuppressionMapper_OutputFormat_Values() []DiagnosticToSuppressionMapper_OutputFormat {
	return j2gort.CloneSlice(enumValues_DiagnosticToSuppressionMapper_OutputFormat)
}

func DiagnosticToSuppressionMapper_OutputFormat_ValueOf(name string) DiagnosticToSuppressionMapper_OutputFormat {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_DiagnosticToSuppressionMapper_OutputFormat {
		if enumBaseTable_DiagnosticToSuppressionMapper_OutputFormat[i_].Name_ == name {
			return DiagnosticToSuppressionMapper_OutputFormat(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.DiagnosticToSuppressionMapper.OutputFormat." + name))
}

func (this_ DiagnosticToSuppressionMapper_OutputFormat) IsJavaNil_() bool                { return this_ == 0 }
func (DiagnosticToSuppressionMapper_OutputFormat) HasIdentityEquals_()                   {}
func (DiagnosticToSuppressionMapper_OutputFormat) Implements_j2gort_Ext_java_lang_Enum() {}
func (DiagnosticToSuppressionMapper_OutputFormat) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ DiagnosticToSuppressionMapper_OutputFormat) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (DiagnosticToSuppressionMapper_OutputFormat) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DiagnosticToSuppressionMapper$OutputFormat", "OutputFormat"
}
func (this_ DiagnosticToSuppressionMapper_OutputFormat) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_DiagnosticToSuppressionMapper_OutputFormat[this_-1]
}
func (this_ DiagnosticToSuppressionMapper_OutputFormat) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ DiagnosticToSuppressionMapper_OutputFormat) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ DiagnosticToSuppressionMapper_OutputFormat) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_DiagnosticToSuppressionMapper_OutputFormat[this_-1].Name_
}
func (this_ DiagnosticToSuppressionMapper_OutputFormat) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_DiagnosticToSuppressionMapper_OutputFormat[this_-1].Name_
}
func (this_ DiagnosticToSuppressionMapper_OutputFormat) Equals(other_ any) bool {
	o_, ok_ := other_.(DiagnosticToSuppressionMapper_OutputFormat)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ DiagnosticToSuppressionMapper_OutputFormat) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_DiagnosticToSuppressionMapper_OutputFormat[this_-1].HashCode()
}
func (this_ DiagnosticToSuppressionMapper_OutputFormat) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_DiagnosticToSuppressionMapper_OutputFormat[this_-1].CompareTo(other_)
}
func (this_ DiagnosticToSuppressionMapper_OutputFormat) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_DiagnosticToSuppressionMapper_OutputFormat[this_-1].Name_)
}
func init() {
	vals_ := []any{DiagnosticToSuppressionMapper_OutputFormat(1), DiagnosticToSuppressionMapper_OutputFormat(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.DiagnosticToSuppressionMapper.OutputFormat", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.DiagnosticToSuppressionMapper$OutputFormat", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DiagnosticToSuppressionMapper.OutputFormat", "com.google.javascript.jscomp.DiagnosticToSuppressionMapper$OutputFormat", func(v_ any) bool {
		return j2gort.InstanceOf[DiagnosticToSuppressionMapper_OutputFormat](v_)
	})
}

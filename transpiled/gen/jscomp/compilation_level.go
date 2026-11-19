/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CompilationLevel.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.CompilationLevel
type CompilationLevel uint8

const (
	enum_CompilationLevel_BUNDLE CompilationLevel = 1
	enum_CompilationLevel_WHITESPACE_ONLY CompilationLevel = 2
	enum_CompilationLevel_SIMPLE_OPTIMIZATIONS CompilationLevel = 3
	enum_CompilationLevel_ADVANCED_OPTIMIZATIONS CompilationLevel = 4
	enum_CompilationLevel_TRANSPILE_ONLY CompilationLevel = 5
)

var enumBaseTable_CompilationLevel = [...]j2gort.EnumBase{
	{Name_: "BUNDLE", Ordinal_: 0},
	{Name_: "WHITESPACE_ONLY", Ordinal_: 1},
	{Name_: "SIMPLE_OPTIMIZATIONS", Ordinal_: 2},
	{Name_: "ADVANCED_OPTIMIZATIONS", Ordinal_: 3},
	{Name_: "TRANSPILE_ONLY", Ordinal_: 4},
}
var enumValues_CompilationLevel = []CompilationLevel{1, 2, 3, 4, 5}

func StaticGet_CompilationLevel_BUNDLE() CompilationLevel {
	return enum_CompilationLevel_BUNDLE
}
func StaticGet_CompilationLevel_WHITESPACE_ONLY() CompilationLevel {
	return enum_CompilationLevel_WHITESPACE_ONLY
}
func StaticGet_CompilationLevel_SIMPLE_OPTIMIZATIONS() CompilationLevel {
	return enum_CompilationLevel_SIMPLE_OPTIMIZATIONS
}
func StaticGet_CompilationLevel_ADVANCED_OPTIMIZATIONS() CompilationLevel {
	return enum_CompilationLevel_ADVANCED_OPTIMIZATIONS
}
func StaticGet_CompilationLevel_TRANSPILE_ONLY() CompilationLevel {
	return enum_CompilationLevel_TRANSPILE_ONLY
}

func CompilationLevel_Values() []CompilationLevel {
	return j2gort.CloneSlice(enumValues_CompilationLevel)
}

func CompilationLevel_ValueOf(name string) CompilationLevel {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CompilationLevel {
		if enumBaseTable_CompilationLevel[i_].Name_ == name {
			return CompilationLevel(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CompilationLevel." + name))
}

// Java static method: CompilationLevel.fromString(java.lang.String value)
func CompilationLevel_FromString(value string) CompilationLevel {
	_ = value
	if ((value == j2gort.NullString)) {
		return CompilationLevel(0)
	}
	switch j2gort.CheckNotNullString(value) {
	case "BUNDLE":
		return StaticGet_CompilationLevel_BUNDLE()
	case "WHITESPACE_ONLY", "WHITESPACE":
		return StaticGet_CompilationLevel_WHITESPACE_ONLY()
	case "TRANSPILE_ONLY":
		return StaticGet_CompilationLevel_TRANSPILE_ONLY()
	case "SIMPLE_OPTIMIZATIONS", "SIMPLE":
		return StaticGet_CompilationLevel_SIMPLE_OPTIMIZATIONS()
	case "ADVANCED_OPTIMIZATIONS", "ADVANCED":
		return StaticGet_CompilationLevel_ADVANCED_OPTIMIZATIONS()
	default:
		return CompilationLevel(0)
	}
	panic("unreachable switch expression")
}

// Java method: CompilationLevel.setOptionsForCompilationLevel(com.google.javascript.jscomp.CompilerOptions options)
func (this_ CompilationLevel) SetOptionsForCompilationLevel(options *CompilerOptions) {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = options
	switch (this_).Ordinal() {
	case int32(0) /* CompilationLevel.BUNDLE */:
	case int32(1) /* CompilationLevel.WHITESPACE_ONLY */:
		CompilationLevel_ApplyBasicCompilationOptions(options)
	case int32(2) /* CompilationLevel.SIMPLE_OPTIMIZATIONS */:
		CompilationLevel_ApplySafeCompilationOptions(options)
	case int32(4) /* CompilationLevel.TRANSPILE_ONLY */:
		CompilationLevel_ApplyTranspileOnlyOptions(options)
	case int32(3) /* CompilationLevel.ADVANCED_OPTIMIZATIONS */:
		CompilationLevel_ApplyFullCompilationOptions(options)
	}
}

// Java method: CompilationLevel.setDebugOptionsForCompilationLevel(com.google.javascript.jscomp.CompilerOptions options)
func (this_ CompilationLevel) SetDebugOptionsForCompilationLevel(options *CompilerOptions) {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = options
	options.SetGeneratePseudoNames(true)
	options.SetRemoveClosureAsserts(false)
	options.SetRemoveJ2clAsserts(true)
}

// Java static method: CompilationLevel.applyBasicCompilationOptions(com.google.javascript.jscomp.CompilerOptions options)
func CompilationLevel_ApplyBasicCompilationOptions(options *CompilerOptions) {
	_ = options
	options.SkipAllCompilerPasses()
}

// Java static method: CompilationLevel.applyTranspileOnlyOptions(com.google.javascript.jscomp.CompilerOptions options)
func CompilationLevel_ApplyTranspileOnlyOptions(options *CompilerOptions) {
	_ = options
	options.SetReplaceIdGenerators(false)
	options.SetClosurePass(true)
	options.SetRenamingPolicy(StaticGet_VariableRenamingPolicy_OFF(), StaticGet_PropertyRenamingPolicy_OFF())
	options.SetInlineVariables_1_CompilerOptions_Reach(StaticGet_CompilerOptions_Reach_NONE())
	options.SetInlineFunctions(StaticGet_CompilerOptions_Reach_NONE())
	options.SetFoldConstants(false)
	options.SetCoalesceVariableNames(false)
	options.SetDeadAssignmentElimination(false)
	options.SetCollapseVariableDeclarations(false)
	options.SetConvertToDottedProperties(false)
	options.SetLabelRenaming(false)
	options.SetRemoveUnusedVariables(StaticGet_CompilerOptions_Reach_NONE())
	options.SetCollapseObjectLiterals(false)
	options.SetProtectHiddenSideEffects(false)
}

// Java static method: CompilationLevel.applySafeCompilationOptions(com.google.javascript.jscomp.CompilerOptions options)
func CompilationLevel_ApplySafeCompilationOptions(options *CompilerOptions) {
	_ = options
	options.SetDependencyOptions(DependencyOptions_SortOnly())
	options.SetReplaceIdGenerators(false)
	options.SetClosurePass(true)
	options.SetRenamingPolicy(StaticGet_VariableRenamingPolicy_LOCAL(), StaticGet_PropertyRenamingPolicy_OFF())
	options.SetInlineVariables_1_CompilerOptions_Reach(StaticGet_CompilerOptions_Reach_LOCAL_ONLY())
	options.SetInlineFunctions(StaticGet_CompilerOptions_Reach_LOCAL_ONLY())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_GLOBAL_THIS(), StaticGet_CheckLevel_OFF())
	options.SetFoldConstants(true)
	options.SetCoalesceVariableNames(true)
	options.SetDeadAssignmentElimination(true)
	options.SetDeadPropertyAssignmentElimination(false)
	options.SetCollapseVariableDeclarations(true)
	options.SetConvertToDottedProperties(true)
	options.SetLabelRenaming(true)
	options.SetRemoveUnusedVariables(StaticGet_CompilerOptions_Reach_LOCAL_ONLY())
	options.SetCollapseObjectLiterals(true)
	options.SetProtectHiddenSideEffects(true)
}

// Java static method: CompilationLevel.applyFullCompilationOptions(com.google.javascript.jscomp.CompilerOptions options)
func CompilationLevel_ApplyFullCompilationOptions(options *CompilerOptions) {
	_ = options
	options.SetDependencyOptions(DependencyOptions_SortOnly())
	options.SetCheckSymbols(true)
	options.SetCheckTypes(true)
	options.SetClosurePass(true)
	options.SetFoldConstants(true)
	options.SetCoalesceVariableNames(true)
	options.SetDeadAssignmentElimination(true)
	options.SetExtractPrototypeMemberDeclarations_1_bool(true)
	options.SetCollapseVariableDeclarations(true)
	options.SetConvertToDottedProperties(true)
	options.SetLabelRenaming(true)
	options.SetCollapseObjectLiterals(true)
	options.SetProtectHiddenSideEffects(true)
	options.SetRemoveClosureAsserts(true)
	options.SetRemoveAbstractMethods(true)
	options.SetReserveRawExports(true)
	options.SetRenamingPolicy(StaticGet_VariableRenamingPolicy_ALL(), StaticGet_PropertyRenamingPolicy_ALL_UNQUOTED())
	options.SetRemoveUnusedPrototypeProperties(true)
	options.SetRemoveUnusedClassProperties(true)
	options.SetCollapseAnonymousFunctions(true)
	options.SetCollapsePropertiesLevel(StaticGet_CompilerOptions_PropertyCollapseLevel_ALL())
	options.SetWarningLevel(StaticGet_DiagnosticGroups_GLOBAL_THIS(), StaticGet_CheckLevel_WARNING())
	options.SetRewriteFunctionExpressions(false)
	options.SetSmartNameRemoval(true)
	options.SetInlineConstantVars(true)
	options.SetInlineFunctions(StaticGet_CompilerOptions_Reach_ALL())
	options.SetInlineVariables_1_CompilerOptions_Reach(StaticGet_CompilerOptions_Reach_ALL())
	options.SetComputeFunctionSideEffects(true)
	options.SetAssumeStrictThis(true)
	options.SetRemoveUnusedVariables(StaticGet_CompilerOptions_Reach_ALL())
	options.SetCrossChunkCodeMotion(true)
	options.SetCrossChunkMethodMotion(true)
	options.SetDevirtualizeMethods(true)
	options.SetOptimizeCalls(true)
	options.SetOptimizeESClassConstructors(true)
}

// Java method: CompilationLevel.setTypeBasedOptimizationOptions(com.google.javascript.jscomp.CompilerOptions options)
func (this_ CompilationLevel) SetTypeBasedOptimizationOptions(options *CompilerOptions) {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = options
	switch (this_).Ordinal() {
	case int32(3) /* CompilationLevel.ADVANCED_OPTIMIZATIONS */:
		options.SetDisambiguateProperties(true)
		options.SetAmbiguateProperties(true)
		options.SetInlineProperties(true)
		options.SetUseTypesForLocalOptimization(true)
	case int32(2) /* CompilationLevel.SIMPLE_OPTIMIZATIONS */, int32(1) /* CompilationLevel.WHITESPACE_ONLY */, int32(0) /* CompilationLevel.BUNDLE */, int32(4) /* CompilationLevel.TRANSPILE_ONLY */:
	}
}

// Java method: CompilationLevel.setWrappedOutputOptimizations(com.google.javascript.jscomp.CompilerOptions options)
func (this_ CompilationLevel) SetWrappedOutputOptimizations(options *CompilerOptions) {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = options
	options.SetReserveRawExports(false)
	switch (this_).Ordinal() {
	case int32(2) /* CompilationLevel.SIMPLE_OPTIMIZATIONS */:
		options.SetVariableRenaming(StaticGet_VariableRenamingPolicy_ALL())
		options.SetCollapsePropertiesLevel(StaticGet_CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT())
		options.SetCollapseAnonymousFunctions(true)
		options.SetInlineConstantVars(true)
		options.SetInlineFunctions(StaticGet_CompilerOptions_Reach_ALL())
		options.SetInlineVariables_1_CompilerOptions_Reach(StaticGet_CompilerOptions_Reach_ALL())
		options.SetRemoveUnusedVariables(StaticGet_CompilerOptions_Reach_ALL())
	case int32(3) /* CompilationLevel.ADVANCED_OPTIMIZATIONS */, int32(1) /* CompilationLevel.WHITESPACE_ONLY */, int32(0) /* CompilationLevel.BUNDLE */, int32(4) /* CompilationLevel.TRANSPILE_ONLY */:
	}
}

func (this_ CompilationLevel) IsJavaNil_() bool { return this_ == 0 }
func (CompilationLevel) HasIdentityEquals_() {}
func (CompilationLevel) Implements_j2gort_Ext_java_lang_Enum() {}
func (CompilationLevel) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CompilationLevel) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (CompilationLevel) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.CompilationLevel", "CompilationLevel" }
func (this_ CompilationLevel) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CompilationLevel[this_-1]
}
func (this_ CompilationLevel) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CompilationLevel) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CompilationLevel) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilationLevel[this_-1].Name_
}
func (this_ CompilationLevel) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CompilationLevel[this_-1].Name_
}
func (this_ CompilationLevel) Equals(other_ any) bool {
	o_, ok_ := other_.(CompilationLevel)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CompilationLevel) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CompilationLevel[this_-1].HashCode()
}
func (this_ CompilationLevel) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CompilationLevel[this_-1].CompareTo(other_)
}
func (this_ CompilationLevel) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CompilationLevel[this_-1].Name_)
}
func init() {
	vals_ := []any{CompilationLevel(1), CompilationLevel(2), CompilationLevel(3), CompilationLevel(4), CompilationLevel(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CompilationLevel", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilationLevel", "com.google.javascript.jscomp.CompilationLevel", func(v_ any) bool {
		return j2gort.InstanceOf[CompilationLevel](v_)
	})
}


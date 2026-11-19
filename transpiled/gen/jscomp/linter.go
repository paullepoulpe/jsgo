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
// Port of com/google/javascript/jscomp/Linter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
)

// Java class: com.google.javascript.jscomp.Linter
type Linter struct {
	options *CompilerOptions `gson:"options"`
	idHash_ atomic.Int32
}

func (this_ *Linter) BindSelf_(s any) {
	_ = s
}
func (this_ *Linter) IsJavaNil_() bool { return this_ == nil }
func (*Linter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Linter", "Linter"
}
func (this_ *Linter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Linter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Linter", "com.google.javascript.jscomp.Linter", func(v_ any) bool {
		return j2gort.InstanceOf[*Linter](v_)
	})
}

func StaticGet_Linter_MAX_FIXES() int32 { return int32(5) }
func (this_ *Linter) initFields_Linter() {
	_ = this_
}

// Java constructor: Linter(com.google.javascript.jscomp.CompilerOptions options)
func NewLinter(options *CompilerOptions) *Linter {
	this_ := &Linter{}
	InitLinter(this_, options)
	return this_
}

func InitLinter(this_ *Linter, options *CompilerOptions) {
	_ = options
	this_.initFields_Linter()
	this_.options = options
}

// Java static method: Linter.builder()
func Linter_S_Builder() *Linter_Builder {
	return NewLinter_Builder()
}

// Java method: Linter.lint(java.lang.String filename)
func (this_ *Linter) Lint_1(filename string) {
	_ = filename
	this_.Lint_2(j2gort.Ext_java_nio_file_Path_Of_2(filename, []string{}), NewCompiler_1_Ext_java_io_PrintStream(j2gort.Ext_java_lang_System_F_Out))
}

// Java method: Linter.lint(java.nio.file.Path path, com.google.javascript.jscomp.Compiler compiler)
func (this_ *Linter) Lint_2(path j2gort.Ext_java_nio_file_Path, compiler *Compiler) {
	_ = path
	_ = compiler
	var file *SourceFile = SourceFile_FromFile_1(path.String())
	_ = file
	compiler.SetPassConfig(NewLintPassConfig(this_.options))
	compiler.DisableThreads()
	var externs *SourceFile = SourceFile_FromCode_2("<Linter externs>", "")
	_ = externs
	compiler.Compile_3_Ext_java_util_List_Ext_java_util_List_CompilerOptions(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(externs)), j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(file)), this_.options)
}

// Java method: Linter.fixRepeatedly(java.lang.String filename)
func (this_ *Linter) FixRepeatedly_1(filename string) {
	_ = filename
	this_.FixRepeatedly_2(filename, j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
}

// Java method: Linter.fixRepeatedly(java.lang.String filename, com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.DiagnosticType> unfixableErrors)
func (this_ *Linter) FixRepeatedly_2(filename string, unfixableErrors j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = filename
	_ = unfixableErrors
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(5); i = i + 1 {
			if !(this_.priv_Linter_fix_3787(filename, unfixableErrors)) {
				break
			}
		}
	}
}

// Java method: Linter.fix(java.lang.String filename, com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.DiagnosticType> unfixableErrors)
func (this_ *Linter) priv_Linter_fix_3787(filename string, unfixableErrors j2gort.Ext_com_google_common_collect_ImmutableSet) bool {
	_ = filename
	_ = unfixableErrors
	var compiler *Compiler = NewCompiler_1_Ext_java_io_PrintStream(j2gort.Ext_java_lang_System_F_Out)
	_ = compiler
	var errorManager *FixingErrorManager = NewFixingErrorManager_1(unfixableErrors)
	_ = errorManager
	compiler.SetErrorManager(errorManager)
	errorManager.SetCompiler(compiler)
	this_.Lint_2(j2gort.Ext_java_nio_file_Path_Of_2(filename, []string{}), compiler)
	var fixes j2gort.Ext_java_util_Collection = errorManager.GetSureFixes()
	_ = fixes
	if !(fixes.IsEmpty()) {
		ApplySuggestedFixes_ApplySuggestedFixesToFiles(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](fixes))
		return true
	}
	return false
}

// Java class: com.google.javascript.jscomp.Linter.Builder
type Linter_Builder struct {
	options *CompilerOptions `gson:"options"`
	idHash_ atomic.Int32
}

func (this_ *Linter_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *Linter_Builder) IsJavaNil_() bool { return this_ == nil }
func (*Linter_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Linter$Builder", "Builder"
}
func (this_ *Linter_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Linter_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Linter.Builder", "com.google.javascript.jscomp.Linter$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*Linter_Builder](v_)
	})
}

func (this_ *Linter_Builder) initFields_Linter_Builder() {
	_ = this_
}

// Java constructor: Builder()
func NewLinter_Builder() *Linter_Builder {
	this_ := &Linter_Builder{}
	InitLinter_Builder(this_)
	return this_
}

func InitLinter_Builder(this_ *Linter_Builder) {
	this_.initFields_Linter_Builder()
	this_.options = NewCompilerOptions()
	this_.options.SetLanguage(StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_NEXT())
	this_.options.SetParseJsDocDocumentation(pkg_parsing.StaticGet_Config_JsDocParsing_INCLUDE_ALL_COMMENTS())
	this_.options.SetPreserveDetailedSourceInfo(true)
	this_.options.SetPrettyPrint(true)
	this_.options.SetPreserveTypeAnnotations(true)
	this_.options.SetPreferSingleQuotes(true)
	this_.options.SetEmitUseStrict(false)
	this_.options.SetParseJsDocDocumentation(pkg_parsing.StaticGet_Config_JsDocParsing_INCLUDE_ALL_COMMENTS())
	this_.options.SetCodingConvention(NewGoogleCodingConvention())
	this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_CHECK_TYPES(), StaticGet_CheckLevel_WARNING())
	this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_JSDOC_MISSING_TYPE(), StaticGet_CheckLevel_ERROR())
	this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_MISPLACED_MSG_ANNOTATION(), StaticGet_CheckLevel_WARNING())
	this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_LINT_CHECKS(), StaticGet_CheckLevel_WARNING())
	this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_UNUSED_LOCAL_VARIABLE(), StaticGet_CheckLevel_WARNING())
	this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_EXTRA_REQUIRE(), StaticGet_CheckLevel_ERROR())
	this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_MISPLACED_SUPPRESS(), StaticGet_CheckLevel_WARNING())
	this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_TYPE_IMPORT_CODE_REFERENCES(), StaticGet_CheckLevel_ERROR())
	this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_MODULE_LOAD(), StaticGet_CheckLevel_OFF())
	this_.options.SetWarningLevel(StaticGet_DiagnosticGroups_STRICT_MODULE_CHECKS(), StaticGet_CheckLevel_WARNING())
	this_.options.SetSummaryDetailLevel(int32(0))
}

// Java method: Builder.withModuleResolutionMode(com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode moduleResolutionMode)
func (this_ *Linter_Builder) WithModuleResolutionMode(moduleResolutionMode ModuleLoader_ResolutionMode) *Linter_Builder {
	_ = moduleResolutionMode
	this_.options.SetModuleResolutionMode(moduleResolutionMode)
	return this_
}

// Java method: Builder.withBrowserResolverPrefixReplacements(com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> replacements)
func (this_ *Linter_Builder) WithBrowserResolverPrefixReplacements(replacements j2gort.Ext_com_google_common_collect_ImmutableMap) *Linter_Builder {
	_ = replacements
	this_.options.SetBrowserResolverPrefixReplacements(replacements)
	return this_
}

// Java method: Builder.build()
func (this_ *Linter_Builder) Build() *Linter {
	return NewLinter(this_.options)
}

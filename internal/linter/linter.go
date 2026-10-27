// Port of com/google/javascript/jscomp/Linter.java @ closure-compiler v20261021

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

package linter

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/refactoring"
)

// Linter is a tool for running just the lint checks which can be run on a
// single file at a time.
//
// Run from LinterMain.
//
// Java: final class Linter
type Linter struct {
	options *jscomp.CompilerOptions
}

// Linter_Builder is a builder for a Linter that allows some customization.
//
// Note: A builder is not designed to generate multiple Linters, just a
// single one. The underlying CompilerOptions it builds is never copied
// defensively, so futher edits to this builder will affect previously built
// Linters.
//
// Java: static final class Linter.Builder
type Linter_Builder struct {
	options *jscomp.CompilerOptions
}

// Java: Linter.Builder()
func newLinter_Builder() *Linter_Builder {
	b := &Linter_Builder{}
	b.options = jscomp.NewCompilerOptions()
	b.options.SetLanguage(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	b.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	b.options.SetPreserveDetailedSourceInfo(true)

	// These are necessary to make sure that suggested fixes are printed
	// correctly.
	b.options.SetPrettyPrint(true)
	b.options.SetPreserveTypeAnnotations(true)
	b.options.SetPreferSingleQuotes(true)
	b.options.SetEmitUseStrict(false)

	b.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	b.options.SetCodingConvention(jscomp.NewGoogleCodingConvention())

	// Even though we're not running the typechecker, enable the checkTypes
	// DiagnosticGroup, since it contains some warnings we do want to report,
	// such as JSDoc parse warnings.
	b.options.SetWarningLevel(jscomp.DiagnosticGroups_CHECK_TYPES, jscomp.CheckLevel_WARNING)

	b.options.SetWarningLevel(jscomp.DiagnosticGroups_JSDOC_MISSING_TYPE, jscomp.CheckLevel_ERROR)
	b.options.SetWarningLevel(jscomp.DiagnosticGroups_MISPLACED_MSG_ANNOTATION, jscomp.CheckLevel_WARNING)
	b.options.SetWarningLevel(jscomp.DiagnosticGroups_LINT_CHECKS, jscomp.CheckLevel_WARNING)
	b.options.SetWarningLevel(jscomp.DiagnosticGroups_UNUSED_LOCAL_VARIABLE, jscomp.CheckLevel_WARNING)
	b.options.SetWarningLevel(jscomp.DiagnosticGroups_EXTRA_REQUIRE, jscomp.CheckLevel_ERROR)
	b.options.SetWarningLevel(jscomp.DiagnosticGroups_MISPLACED_SUPPRESS, jscomp.CheckLevel_WARNING)
	b.options.SetWarningLevel(jscomp.DiagnosticGroups_TYPE_IMPORT_CODE_REFERENCES, jscomp.CheckLevel_ERROR)
	b.options.SetWarningLevel(jscomp.DiagnosticGroups_MODULE_LOAD, jscomp.CheckLevel_OFF)
	b.options.SetWarningLevel(jscomp.DiagnosticGroups_STRICT_MODULE_CHECKS, jscomp.CheckLevel_WARNING)
	b.options.SetSummaryDetailLevel(0)
	return b
}

// Java: Linter.Builder.withModuleResolutionMode(ResolutionMode)
func (b *Linter_Builder) WithModuleResolutionMode(moduleResolutionMode jscomp.ModuleLoader_ResolutionMode) *Linter_Builder {
	b.options.SetModuleResolutionMode(moduleResolutionMode)
	return b
}

// Java: Linter.Builder.withBrowserResolverPrefixReplacements(ImmutableMap<String, String>)
func (b *Linter_Builder) WithBrowserResolverPrefixReplacements(
	replacements *jutil.LinkedHashMap[string, string]) *Linter_Builder {
	b.options.SetBrowserResolverPrefixReplacements(replacements)
	return b
}

// Java: Linter.Builder.build()
func (b *Linter_Builder) Build() *Linter {
	return newLinter(b.options)
}

// Options returns the underlying CompilerOptions configured on this builder.
func (b *Linter_Builder) Options() *jscomp.CompilerOptions {
	return b.options
}

// Linter_MAX_FIXES: don't try to apply fixes anymore, after trying this
// many times. This is to avoid the unlikely event of an infinite loop of
// fixes.
//
// Java: Linter.MAX_FIXES
const Linter_MAX_FIXES = 5

// Java: Linter(CompilerOptions)
func newLinter(options *jscomp.CompilerOptions) *Linter {
	l := &Linter{}
	l.options = options
	return l
}

// Linter_Builder_ is Linter.builder() (the trailing underscore keeps it
// apart from the type Linter_Builder).
//
// Java: Linter.builder()
func Linter_Builder_() *Linter_Builder {
	return newLinter_Builder()
}

// Java: Linter.lint(String)
func (l *Linter) Lint(filename string) {
	l.LintWithCompiler(jutil.Path_Of(filename), jscomp.NewCompilerWithOutStream(jutil.System_Out()))
}

// Java: Linter.lint(Path, Compiler)
func (l *Linter) LintWithCompiler(path *jutil.Path, compiler *jscomp.Compiler) {
	file := jscomp.SourceFile_FromFile(path.String())
	compiler.SetPassConfig(NewLintPassConfig(l.options))
	compiler.DisableThreads()
	externs := jscomp.SourceFile_FromCode("<Linter externs>", "")
	compiler.Compile([]*jscomp.SourceFile{externs}, []*jscomp.SourceFile{file}, l.options)
}

// FixRepeatedly keeps applying fixes to the given file until no more fixes
// can be found, or until fixes have been applied MAX_FIXES times. It throws
// IOException.
//
// Java: Linter.fixRepeatedly(String)
func (l *Linter) FixRepeatedly(filename string) {
	l.FixRepeatedlyWithUnfixableErrors(filename, jutil.NewObjectLinkedHashSet[*jscomp.DiagnosticType]())
}

// FixRepeatedlyWithUnfixableErrors keeps applying fixes to the given file
// until no more fixes can be found, or until fixes have been applied
// MAX_FIXES times. It throws IOException.
//
// Java: Linter.fixRepeatedly(String, ImmutableSet<DiagnosticType>)
func (l *Linter) FixRepeatedlyWithUnfixableErrors(filename string, unfixableErrors *jutil.ObjectLinkedHashSet[*jscomp.DiagnosticType]) {
	for i := 0; i < Linter_MAX_FIXES; i++ {
		if !l.fix(filename, unfixableErrors) {
			break
		}
	}
}

// fix returns whether any fixes were applied. It throws IOException.
//
// Java: Linter.fix(String, ImmutableSet<DiagnosticType>)
func (l *Linter) fix(filename string, unfixableErrors *jutil.ObjectLinkedHashSet[*jscomp.DiagnosticType]) bool {
	compiler := jscomp.NewCompilerWithOutStream(jutil.System_Out())
	errorManager := refactoring.NewFixingErrorManagerWithUnfixableErrors(unfixableErrors)
	compiler.SetErrorManager(errorManager)
	errorManager.SetCompiler(compiler)

	l.LintWithCompiler(jutil.Path_Of(filename), compiler)

	fixes := errorManager.GetSureFixes()
	if len(fixes) != 0 {
		refactoring.ApplySuggestedFixes_ApplySuggestedFixesToFiles(slices.Values(fixes))
		return true
	}
	return false
}

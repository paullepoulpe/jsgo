// Port of com/google/javascript/jscomp/Linter.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/LinterMain.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/LintPassConfig.java @ closure-compiler v20261021

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

// Package linter implements single-file Closure lint checking and iterative
// auto-fixing (Linter, LinterMain, and LintPassConfig).
package linter

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/args4j"
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
	return &Linter{options: options}
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
	l.FixRepeatedlyWithUnfixableErrors(filename, nil)
}

// FixRepeatedlyWithUnfixableErrors keeps applying fixes to the given file
// until no more fixes can be found, or until fixes have been applied
// MAX_FIXES times. It throws IOException.
//
// Java: Linter.fixRepeatedly(String, ImmutableSet<DiagnosticType>)
func (l *Linter) FixRepeatedlyWithUnfixableErrors(filename string, unfixableErrors []*jscomp.DiagnosticType) {
	for i := 0; i < Linter_MAX_FIXES; i++ {
		if !l.fix(filename, unfixableErrors) {
			break
		}
	}
}

// fix returns whether any fixes were applied. It throws IOException.
//
// Java: Linter.fix(String, ImmutableSet<DiagnosticType>)
func (l *Linter) fix(filename string, unfixableErrors []*jscomp.DiagnosticType) bool {
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

// LinterMain is a minimal binary that just runs the "lint" checks which can
// be run on a single file at a time. This means some checks in the
// lintChecks DiagnosticGroup are skipped, since they depend on type
// information.
//
// Java: final class LinterMain
type LinterMain struct {
	fix            bool
	resolutionMode jscomp.ModuleLoader_ResolutionMode
	files          []string
}

// linterMain_ResolutionModeClass is the Class of the enum
// ModuleLoader.ResolutionMode, the type of the --resolution_mode option.
//
// Go-only: language-adapter — reflection (see package args4j).
var linterMain_ResolutionModeClass = args4j.EnumClassOf(
	"com.google.javascript.jscomp.deps.ModuleLoader$ResolutionMode",
	jscomp.ModuleLoader_ResolutionMode_values())

// Java: LinterMain()
func newLinterMain() *LinterMain {
	return &LinterMain{
		fix:            false,
		resolutionMode: jscomp.ModuleLoader_ResolutionMode_BROWSER,
		files:          []string{},
	}
}

// Args4jFields lists the fields that carry args4j's @Option and @Argument
// annotations, in declaration order.
//
// Go-only: language-adapter — annotations (see package args4j).
func (m *LinterMain) Args4jFields() []args4j.AnnotatedField {
	const c = "com.google.javascript.jscomp.LinterMain"
	return []args4j.AnnotatedField{
		{
			Option: &args4j.Option{Name: "--fix", Usage: "Fix lint warnings automatically"},
			Field:  args4j.BooleanField(c, "fix", &m.fix),
		},
		{
			Option: &args4j.Option{
				Name:  "--resolution_mode",
				Usage: "Expected path resolution mode for resolving module specifiers."},
			Field: args4j.EnumField(c, "resolutionMode", linterMain_ResolutionModeClass, &m.resolutionMode),
		},
		{
			Argument: &args4j.Argument{},
			Field:    args4j.ListField(c, "files", args4j.Class_String, &m.files),
		},
	}
}

// LinterMain_Main throws IOException and CmdLineException.
//
// Java: LinterMain.main(String[])
func LinterMain_Main(args []string) {
	newLinterMain().run(args)
}

// Java: LinterMain.run(String[])
func (m *LinterMain) run(args []string) {
	parser := args4j.NewCmdLineParser(m)
	parser.ParseArgument(args...)

	linter := Linter_Builder_().WithModuleResolutionMode(m.resolutionMode).Build()

	for _, filename := range m.files {
		if m.fix {
			linter.FixRepeatedly(filename)
		} else {
			linter.Lint(filename)
		}
	}
}

// LintPassConfig is a PassConfig for the standalone linter, which runs on a
// single file at a time. This runs a similar set of checks to what you would
// get when running the compiler with the lintChecks DiagnosticGroup enabled,
// but some of the lint checks depend on type information, which is not
// available when looking at a single file, so those are omitted here.
//
// Java: class LintPassConfig extends PassConfig.PassConfigDelegate
type LintPassConfig struct {
	jscomp.PassConfigDelegate

	gatherModuleMetadataPass *jscomp.PassFactory
	earlyLintChecks          *jscomp.PassFactory
	variableReferenceCheck   *jscomp.PassFactory
	lateLintChecks           *jscomp.PassFactory
}

// Java: LintPassConfig(CompilerOptions)
func NewLintPassConfig(options *jscomp.CompilerOptions) *LintPassConfig {
	c := &LintPassConfig{
		PassConfigDelegate: jscomp.NewPassConfigDelegate(jscomp.NewDefaultPassConfig(options)),
	}

	c.gatherModuleMetadataPass =
		jscomp.PassFactory_Builder_().
			SetName(jscomp.PassNames_GATHER_MODULE_METADATA).
			SetInternalFactory(
				func(compiler jscomp.AbstractCompiler) jscomp.CompilerPass {
					return jscomp.NewGatherModuleMetadata(
						compiler,
						compiler.GetOptions().GetProcessCommonJSModules(),
						compiler.GetOptions().GetModuleResolutionMode())
				}).
			Build()

	c.earlyLintChecks =
		jscomp.PassFactory_Builder_().
			SetName("earlyLintChecks").
			SetInternalFactory(
				func(compiler jscomp.AbstractCompiler) jscomp.CompilerPass {
					return jscomp.NewCombinedCompilerPassOfList(
						compiler,
						[]jscomp.NodeTraversal_Callback{
							jscomp.NewCheckConstPrivateProperties(compiler),
							jscomp.NewCheckConstantCaseNames(compiler),
							jscomp.NewCheckDuplicateCase(compiler),
							jscomp.NewCheckEmptyStatements(compiler),
							jscomp.NewCheckEnums(compiler),
							jscomp.NewCheckExtraRequires(compiler, c.Options().GetUnusedImportsToRemove()),
							jscomp.NewCheckGoogModuleTypeScriptName(compiler),
							jscomp.NewCheckJSDocStyle(compiler),
							jscomp.NewCheckJSDoc(compiler),
							jscomp.NewCheckMissingSemicolon(compiler),
							jscomp.NewCheckSuper(compiler),
							jscomp.NewCheckPrimitiveAsObject(compiler),
							jscomp.NewClosureCheckModule(compiler, compiler.GetModuleMetadataMap()),
							jscomp.NewCheckNullabilityModifiers(compiler),
							jscomp.NewCheckProvidesSorted(jscomp.CheckProvidesSorted_Mode_COLLECT_AND_REPORT),
							jscomp.NewCheckRequiresSorted(jscomp.CheckRequiresSorted_Mode_COLLECT_AND_REPORT),
							jscomp.NewCheckSideEffects(
								compiler /* report= */, true /* protectSideEffectFreeCode= */, false),
							jscomp.NewCheckTypeImportCodeReferences(compiler),
							jscomp.NewCheckUnusedLabels(compiler),
							jscomp.NewCheckUnusedPrivateProperties(compiler),
							jscomp.NewCheckUselessBlocks(compiler),
							jscomp.NewCheckVar(compiler)})
				}).
			Build()

	c.variableReferenceCheck =
		jscomp.PassFactory_Builder_().
			SetName("variableReferenceCheck").
			SetRunInFixedPointLoop(true).
			SetInternalFactory(func(compiler jscomp.AbstractCompiler) jscomp.CompilerPass {
				return jscomp.NewVariableReferenceCheck(compiler)
			}).
			Build()

	c.lateLintChecks =
		jscomp.PassFactory_Builder_().
			SetName("lateLintChecks").
			SetInternalFactory(
				func(compiler jscomp.AbstractCompiler) jscomp.CompilerPass {
					return jscomp.NewCombinedCompilerPassOfList(
						compiler,
						[]jscomp.NodeTraversal_Callback{
							jscomp.NewCheckInterfaces(compiler), jscomp.NewCheckPrototypeProperties(compiler)})
				}).
			Build()
	return c
}

// Java: LintPassConfig.getChecks()
func (c *LintPassConfig) GetChecks() *jscomp.PassListBuilder {
	passes := jscomp.NewPassListBuilder(c.Options())
	passes.MaybeAdd(c.gatherModuleMetadataPass)
	passes.MaybeAdd(c.earlyLintChecks)
	passes.MaybeAdd(c.variableReferenceCheck)
	passes.MaybeAdd(c.lateLintChecks)
	return passes
}

// Java: LintPassConfig.getOptimizations()
func (c *LintPassConfig) GetOptimizations() *jscomp.PassListBuilder {
	return jscomp.NewPassListBuilder(c.Options())
}

// Java: LintPassConfig.getFinalizations()
func (c *LintPassConfig) GetFinalizations() *jscomp.PassListBuilder {
	return jscomp.NewPassListBuilder(c.Options())
}

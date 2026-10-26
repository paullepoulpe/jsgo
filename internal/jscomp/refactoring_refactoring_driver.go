// Port of com/google/javascript/refactoring/RefactoringDriver.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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

package jscomp

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RefactoringDriver is the primary driver of a refactoring. This class
// collects the inputs, runs the refactoring over the compiled input, and
// then collects the suggested fixes based on the refactoring.
//
// Java: final class refactoring.RefactoringDriver
type RefactoringDriver struct {
	compiler *Compiler
	rootNode rhino.Node
}

// Java: RefactoringDriver(List<SourceFile>, List<SourceFile>, CompilerOptions)
func newRefactoringDriver(
	inputs []*SourceFile,
	externs []*SourceFile,
	compilerOptions *CompilerOptions) *RefactoringDriver {
	d := &RefactoringDriver{}
	d.compiler = refactoringDriver_createCompiler(inputs, externs, compilerOptions)
	d.rootNode = d.compiler.GetRoot()
	return d
}

// DriveWithIncludeFilePattern runs a refactoring and returns any suggested
// fixes as a result. includeFilePattern may be nil.
//
// Java: RefactoringDriver.drive(Scanner, Pattern)
func (d *RefactoringDriver) DriveWithIncludeFilePattern(scanner Scanner, includeFilePattern *jutil.Pattern) *jutil.ArrayList[*SuggestedFix] {
	callback := NewJsFlumeCallback(scanner, includeFilePattern)
	NodeTraversal_Traverse(d.compiler, d.rootNode, callback)
	fixes := callback.GetFixes()
	fixes.AddAll(scanner.ProcessAllMatches(callback.GetMatches())...)
	return fixes
}

// Drive runs a refactoring and returns any suggested fixes as a result.
//
// Java: RefactoringDriver.drive(Scanner)
func (d *RefactoringDriver) Drive(scanner Scanner) *jutil.ArrayList[*SuggestedFix] {
	return d.DriveWithIncludeFilePattern(scanner, nil)
}

// Java: RefactoringDriver.getCompiler()
func (d *RefactoringDriver) GetCompiler() *Compiler {
	return d.compiler
}

// Java: RefactoringDriver.createCompiler(List<SourceFile>, List<SourceFile>, CompilerOptions)
func refactoringDriver_createCompiler(
	inputs []*SourceFile, externs []*SourceFile, compilerOptions *CompilerOptions) *Compiler {
	compiler := NewCompilerWithErrorManager(NewBlackHoleErrorManager())
	compiler.DisableThreads()
	compiler.Compile(externs, inputs, compilerOptions)
	return compiler
}

// TODO(tbreisacher): Make this package-private by refactoring tests so they
// don't need to call it directly.

// Java: RefactoringDriver.getCompilerOptions()
func RefactoringDriver_GetCompilerOptions() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	options.SetSummaryDetailLevel(0)

	options.SetDependencyOptions(DependencyOptions_SortOnly())

	options.SetChecksOnly(true)
	options.SetContinueAfterErrors(true)
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	options.SetCheckSuspiciousCode(true)
	options.SetCheckSymbols(true)
	options.SetCheckTypes(true)
	// TODO(bangert): Remove this -- we want to rewrite code before closure
	// syntax is removed. Unfortunately, setClosurePass is required, or code
	// doesn't type check.
	options.SetClosurePass(true)
	options.SetGenerateExports(true)
	options.SetPreserveClosurePrimitives(true)

	options.SetWarningLevel(DiagnosticGroups_EXTRA_REQUIRE, CheckLevel_WARNING)

	return options
}

// Java: static class RefactoringDriver.Builder
type RefactoringDriver_Builder struct {
	inputs          []*SourceFile
	externs         []*SourceFile
	compilerOptions *CompilerOptions
}

// Java: RefactoringDriver.Builder.TO_SOURCE_FILE_FN
func refactoringDriver_Builder_TO_SOURCE_FILE_FN(file string) *SourceFile {
	return SourceFile_Builder_().WithPath(file).Build()
}

// Java: RefactoringDriver.Builder()
func NewRefactoringDriver_Builder() *RefactoringDriver_Builder {
	return &RefactoringDriver_Builder{
		compilerOptions: RefactoringDriver_GetCompilerOptions(),
	}
}

// Java: RefactoringDriver.Builder.addExternsFromFile(String)
func (b *RefactoringDriver_Builder) AddExternsFromFile(filename string) *RefactoringDriver_Builder {
	b.externs = append(b.externs, SourceFile_FromFile(filename))
	return b
}

// Java: RefactoringDriver.Builder.addExternsFromFile(Iterable<String>)
func (b *RefactoringDriver_Builder) AddExternsFromFileOfList(externs []string) *RefactoringDriver_Builder {
	for _, f := range externs {
		b.externs = append(b.externs, refactoringDriver_Builder_TO_SOURCE_FILE_FN(f))
	}
	return b
}

// Java: RefactoringDriver.Builder.addExternsFromCode(String)
func (b *RefactoringDriver_Builder) AddExternsFromCode(code string) *RefactoringDriver_Builder {
	b.externs = append(b.externs, SourceFile_FromCode("externs", code))
	return b
}

// Java: RefactoringDriver.Builder.addExterns(Iterable<SourceFile>)
func (b *RefactoringDriver_Builder) AddExterns(externs []*SourceFile) *RefactoringDriver_Builder {
	b.externs = append(b.externs, externs...)
	return b
}

// Java: RefactoringDriver.Builder.addInputsFromFile(String)
func (b *RefactoringDriver_Builder) AddInputsFromFile(filename string) *RefactoringDriver_Builder {
	b.inputs = append(b.inputs, SourceFile_FromFile(filename))
	return b
}

// Java: RefactoringDriver.Builder.addInputsFromFile(Iterable<String>)
func (b *RefactoringDriver_Builder) AddInputsFromFileOfList(inputs []string) *RefactoringDriver_Builder {
	for _, f := range inputs {
		b.inputs = append(b.inputs, refactoringDriver_Builder_TO_SOURCE_FILE_FN(f))
	}
	return b
}

// Java: RefactoringDriver.Builder.addInputsFromCode(String)
func (b *RefactoringDriver_Builder) AddInputsFromCode(code string) *RefactoringDriver_Builder {
	return b.AddInputsFromCodeWithFilename(code, "input")
}

// Java: RefactoringDriver.Builder.addInputsFromCode(String, String)
func (b *RefactoringDriver_Builder) AddInputsFromCodeWithFilename(code string, filename string) *RefactoringDriver_Builder {
	b.inputs = append(b.inputs, SourceFile_FromCode(filename, code))
	return b
}

// Java: RefactoringDriver.Builder.addInputs(Iterable<SourceFile>)
func (b *RefactoringDriver_Builder) AddInputs(inputs []*SourceFile) *RefactoringDriver_Builder {
	b.inputs = append(b.inputs, inputs...)
	return b
}

// Java: RefactoringDriver.Builder.withCompilerOptions(CompilerOptions)
func (b *RefactoringDriver_Builder) WithCompilerOptions(compilerOptions *CompilerOptions) *RefactoringDriver_Builder {
	b.compilerOptions = jutil.CheckNotNull(compilerOptions)
	return b
}

// Java: RefactoringDriver.Builder.build()
func (b *RefactoringDriver_Builder) Build() *RefactoringDriver {
	return newRefactoringDriver(slices.Clone(b.inputs), slices.Clone(b.externs), b.compilerOptions)
}

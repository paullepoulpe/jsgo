// Port of com/google/javascript/jscomp/instrumentation/GoldenFileComparer.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/instrumentation/GoldenFileComparer.java)

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"os"
	"path/filepath"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// GoldenFileComparer tests the compiler by comparing output to data files
// aka golden files. Always start with the options object from this class.
//
// Java: public final class GoldenFileComparer
type GoldenFileComparer struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// GoldenFileComparer_DATA_DIR is Java's "test/com/google/javascript/jscomp/testdata/".
//
// Go-only: language-adapter — the data files are under the package's
// testdata directory (go test runs in the package directory).
var GoldenFileComparer_DATA_DIR = "testdata/"

// Java: GoldenFileComparer.coverageExterns()
func GoldenFileComparer_CoverageExterns() []*SourceFile {
	externs := SourceFile_FromCode("externs", "function Symbol() {}; var window; var self;")
	return []*SourceFile{externs}
}

// Java: GoldenFileComparer.compile(List<SourceFile>, List<SourceFile>, CompilerOptions)
func GoldenFileComparer_Compile(externs []*SourceFile, inputs []*SourceFile, options *CompilerOptions) string {
	compiler := NewCompiler()

	var allExterns []*SourceFile
	allExterns = append(allExterns, externs...)
	allExterns = append(allExterns, GoldenFileComparer_CoverageExterns()...)

	compiler.Compile(allExterns, inputs, options)
	return compiler.ToSource()
}

// Java: GoldenFileComparer.readFile(String)
func GoldenFileComparer_ReadFile(filePath string) string {
	data, err := os.ReadFile(filepath.FromSlash(filePath))
	if err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	return string(data)
}

// Java: GoldenFileComparer.readSource(String)
func GoldenFileComparer_ReadSource(fileName string) *SourceFile {
	source := GoldenFileComparer_ReadFile(GoldenFileComparer_ToFullPath(fileName))
	return SourceFile_FromCode(fileName, source)
}

// Java: GoldenFileComparer.compare(String, String)
func GoldenFileComparer_Compare(compiledSource string, goldenSource string) bool {
	compiledLines := jutil.StringSplitLiteral(compiledSource, "\n")
	goldenLines := jutil.StringSplitLiteral(goldenSource, "\n")

	AssertWithMessage("Num lines of compiled code do not match that of expected").
		That(len(compiledLines)).
		IsEqualTo(len(goldenLines))

	// Loop through each line to make it convenient to pin-point the faulty one
	for i := 0; i < len(compiledLines); i++ {
		AssertWithMessage(
			"Instrumented code does not match expected at line no: %s\nExpected:\n%s\nFound:\n%s\n",
			i+1, goldenLines[i], compiledLines[i]).
			That(compiledLines[i]).
			IsEqualTo(goldenLines[i])
	}

	return true
}

// Java: GoldenFileComparer.toFullPath(String)
func GoldenFileComparer_ToFullPath(fileName string) string {
	return GoldenFileComparer_DATA_DIR + fileName
}

// Java: GoldenFileComparer.compileAndCompare(String, CompilerOptions, List<SourceFile>, List<SourceFile>)
func GoldenFileComparer_CompileAndCompareWithLists(
	goldenFileName string,
	options *CompilerOptions,
	sourceFiles []*SourceFile,
	externsFiles []*SourceFile) {
	compiledSource := GoldenFileComparer_Compile(externsFiles, sourceFiles, options)
	referenceSource := GoldenFileComparer_ReadFile(GoldenFileComparer_ToFullPath(goldenFileName))
	GoldenFileComparer_Compare(compiledSource, referenceSource)
}

// GoldenFileComparer_CompileAndCompareSubsetOfActualToExpected compiles
// sourceFileName with given CompilerOptions and checks that all lines of code
// from goldenFileName are present and in the same order as in compiled
// sourceFile. When comparing a compiled version of the sourceFile which
// contains polyfill code, we only wish to compare the original instrumented
// source code to the golden.
//
// Java: GoldenFileComparer.compileAndCompareSubsetOfActualToExpected(String, CompilerOptions, String)
func GoldenFileComparer_CompileAndCompareSubsetOfActualToExpected(
	goldenFileName string, options *CompilerOptions, sourceFileName string) {
	sourceFiles := []*SourceFile{GoldenFileComparer_ReadSource(sourceFileName)}
	externsFiles := []*SourceFile{}
	compiledSource := GoldenFileComparer_Compile(externsFiles, sourceFiles, options)
	goldenSource := GoldenFileComparer_ReadFile(GoldenFileComparer_ToFullPath(goldenFileName))

	// Splitter.on('\n').split keeps every piece, empty ones included.
	compiledLines := strings.Split(compiledSource, "\n")
	goldenLines := strings.Split(goldenSource, "\n")

	identicalLinesBetweenFiles := 0
	for _, line := range compiledLines {
		if goldenLines[identicalLinesBetweenFiles] == line {
			identicalLinesBetweenFiles++
		}
	}

	AssertWithMessage(
		"Line %s from Golden was not found in compiled version of source file",
		identicalLinesBetweenFiles+1).
		That(identicalLinesBetweenFiles).
		IsEqualTo(len(goldenLines))
}

// GoldenFileComparer_CompileAndCompare compiles one input file and throws if
// the result does not match golden. Pass options from this class, mutated
// with desired options.
//
// Java: GoldenFileComparer.compileAndCompare(String, CompilerOptions, String)
func GoldenFileComparer_CompileAndCompare(
	goldenFileName string, options *CompilerOptions, sourceFileName string) {
	sourceFiles := []*SourceFile{GoldenFileComparer_ReadSource(sourceFileName)}
	externsFiles := []*SourceFile{}
	GoldenFileComparer_CompileAndCompareWithLists(goldenFileName, options, sourceFiles, externsFiles)
}

// GoldenFileComparer_CompileAndCompareTwoFiles compiles two input files and
// throws if the result does not match golden. Pass options from this class,
// mutated with desired options.
//
// Java: GoldenFileComparer.compileAndCompare(String, CompilerOptions, String, String, String)
func GoldenFileComparer_CompileAndCompareTwoFiles(
	goldenFileName string,
	options *CompilerOptions,
	sourceFileName1 string,
	sourceFileName2 string,
	externsFileName string) {
	// Prepare sources
	sourceFiles := []*SourceFile{GoldenFileComparer_ReadSource(sourceFileName1), GoldenFileComparer_ReadSource(sourceFileName2)}

	externsFiles := []*SourceFile{SourceFile_FromFile(GoldenFileComparer_ToFullPath(externsFileName))}

	GoldenFileComparer_CompileAndCompareWithLists(goldenFileName, options, sourceFiles, externsFiles)
}

// GoldenFileComparer_Options: always use these options, they set
// --pretty_print option for easy verification.
//
// Java: GoldenFileComparer.options()
func GoldenFileComparer_Options() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT3)
	// Instrumentation done
	options.SetPrettyPrint(true)
	return options
}

// Java: GoldenFileComparer()
func NewGoldenFileComparer() *GoldenFileComparer {
	return &GoldenFileComparer{}
}

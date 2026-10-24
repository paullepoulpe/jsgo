// Port of com/google/javascript/jscomp/CodePrinterTestBase.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/CodePrinterTestBase.java)

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

package jscomp

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CodePrinterTestBase is the base class for tests that exercise
// CodePrinter.
//
// Subclasses embed it; it is a JUnitTestCase (see junit4_test.go) whose
// @Before is SetUp.
type CodePrinterTestBase struct {
	t *testing.T

	// If this is set, ignore parse warnings and only fail the test
	// for parse errors.
	allowWarnings            bool
	trustedStrings           bool
	preserveTypeAnnotations  bool
	preserveNonJSDocComments bool
	languageMode             CompilerOptions_LanguageMode
	lastCompiler             *Compiler      // nil is null
	outputCharset            *jutil.Charset // nil is null
	diagnosticsToIgnore      []*DiagnosticGroup
}

// InitCodePrinterTestBase runs the field initializers of
// CodePrinterTestBase.
func (b *CodePrinterTestBase) InitCodePrinterTestBase() {
	b.allowWarnings = false
	b.trustedStrings = true
	b.preserveTypeAnnotations = false
	b.preserveNonJSDocComments = false
	b.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT5
	b.lastCompiler = nil
	b.outputCharset = nil
	b.diagnosticsToIgnore = []*DiagnosticGroup{}
}

// JUnitSetT gives the instance the Go test running its test method.
func (b *CodePrinterTestBase) JUnitSetT(t *testing.T) { b.t = t }

// JUnitBefore runs the @Before method.
func (b *CodePrinterTestBase) JUnitBefore() { b.SetUp() }

// JUnitAfter: CodePrinterTestBase has no @After method.
func (b *CodePrinterTestBase) JUnitAfter() {}

// Java: CodePrinterTestBase.setUp()
func (b *CodePrinterTestBase) SetUp() {
	b.allowWarnings = false
	b.preserveTypeAnnotations = false
	b.preserveNonJSDocComments = false
	b.trustedStrings = true
	b.lastCompiler = nil
	b.languageMode = CompilerOptions_LanguageMode_UNSUPPORTED
	b.outputCharset = nil
	b.diagnosticsToIgnore = []*DiagnosticGroup{}
}

// Java: CodePrinterTestBase.parse(String)
func (b *CodePrinterTestBase) Parse(js string) rhino.Node {
	return b.ParseTypeChecked(js /* typeChecked= */, false)
}

// Java: CodePrinterTestBase.parse(String, boolean)
func (b *CodePrinterTestBase) ParseTypeChecked(js string, typeChecked bool) rhino.Node {
	compiler := NewCompiler()
	b.lastCompiler = compiler
	options := NewCompilerOptions()
	options.SetTrustedStrings(b.trustedStrings)
	options.SetPreserveTypeAnnotations(b.preserveTypeAnnotations)
	options.SetOutputCharset(b.outputCharset)
	options.SetLanguageIn(b.languageMode)
	if b.preserveNonJSDocComments {
		options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
		options.SetPreserveNonJSDocComments(true)
	}
	options.SetContinueAfterErrors(true)
	compiler.SetPreferRegexParser(false)
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs", CompilerTestCase_MINIMAL_EXTERNS)},
		[]*SourceFile{SourceFile_FromCode("testcode", js)},
		options)
	externsAndJs := compiler.ParseInputs()
	b.CheckUnexpectedErrorsOrWarnings(compiler, 0)
	root := externsAndJs.GetLastChild()
	externs := externsAndJs.GetFirstChild()

	if typeChecked {
		passConfig := NewDefaultPassConfig(nil)
		inferTypes := passConfig.inferTypes.Create(compiler)
		inferTypes.Process(externs, root)
	}

	b.CheckUnexpectedErrorsOrWarnings(compiler, 0)
	return root.GetFirstChild()
}

// Java: CodePrinterTestBase.checkUnexpectedErrorsOrWarnings(Compiler, int)
func (b *CodePrinterTestBase) CheckUnexpectedErrorsOrWarnings(compiler *Compiler, expected int) {
	actual := 0
	var msg strings.Builder
	for _, err := range compiler.GetErrors() {
		if b.ShouldIgnore(err) {
			continue
		}
		actual++
		msg.WriteString("Error:" + err.String() + "\n")
	}
	if !b.allowWarnings {
		for _, err := range compiler.GetWarnings() {
			if b.ShouldIgnore(err) {
				continue
			}
			actual++
			msg.WriteString("Warning:" + err.String() + "\n")
		}
	}
	if actual != expected {
		AssertWithMessage("Unexpected warnings or errors.\n %s", msg.String()).
			That(actual).
			IsEqualTo(expected)
	}
}

// Java: CodePrinterTestBase.parsePrint(String, CompilerOptions)
func (b *CodePrinterTestBase) ParsePrint(js string, options *CompilerOptions) string {
	return NewCodePrinter_Builder(b.Parse(js)).SetCompilerOptions(options).Build()
}

// Java: CodePrinterTestBase.shouldIgnore(JSError)
func (b *CodePrinterTestBase) ShouldIgnore(error_ *JSError) bool {
	for _, diagnosticGroup := range b.diagnosticsToIgnore {
		if diagnosticGroup.Matches(error_) {
			return true
		}
	}
	return false
}

// CodePrinterTestBase_CompilerOptionBuilder is the abstract static class
// CompilerOptionBuilder.
type CodePrinterTestBase_CompilerOptionBuilder interface {
	// Java: CodePrinterTestBase.CompilerOptionBuilder.setOptions(CompilerOptions)
	SetOptions(options *CompilerOptions)
}

// codePrinterTestBase_CompilerOptionBuilderFunc is an anonymous
// CompilerOptionBuilder subclass.
type codePrinterTestBase_CompilerOptionBuilderFunc func(options *CompilerOptions)

func (f codePrinterTestBase_CompilerOptionBuilderFunc) SetOptions(options *CompilerOptions) {
	f(options)
}

// Java: CodePrinterTestBase.newCompilerOptions(CompilerOptionBuilder)
func (b *CodePrinterTestBase) NewCompilerOptions(builder CodePrinterTestBase_CompilerOptionBuilder) *CompilerOptions {
	options := NewCompilerOptions()
	options.SetOutputCharset(b.outputCharset)
	options.SetTrustedStrings(b.trustedStrings)
	options.SetPreserveTypeAnnotations(b.preserveTypeAnnotations)
	options.SetPreserveNonJSDocComments(b.preserveNonJSDocComments)
	options.SetLanguageOut(b.languageMode)
	builder.SetOptions(options)
	return options
}

// Java: CodePrinterTestBase.printNode(Node)
func (b *CodePrinterTestBase) PrintNode(n rhino.Node) string {
	options := NewCompilerOptions()
	options.SetLineLengthThreshold(CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD)
	options.SetLanguageOut(b.languageMode)
	options.SetOutputCharset(b.outputCharset)
	return NewCodePrinter_Builder(n).SetCompilerOptions(options).Build()
}

// Java: CodePrinterTestBase.assertPrintNode(String, Node)
func (b *CodePrinterTestBase) AssertPrintNode(expectedJs string, ast rhino.Node) {
	AssertThat(b.PrintNode(ast)).IsEqualTo(expectedJs)
}

// Java: CodePrinterTestBase.prettyPrintNode(Node)
func (b *CodePrinterTestBase) PrettyPrintNode(n rhino.Node) string {
	options := NewCompilerOptions()
	options.SetLineLengthThreshold(CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD)
	options.SetLanguageOut(b.languageMode)
	options.SetOutputCharset(b.outputCharset)
	options.SetPrettyPrint(true)
	return NewCodePrinter_Builder(n).SetCompilerOptions(options).Build()
}

// Java: CodePrinterTestBase.assertPrettyPrintNode(String, Node)
func (b *CodePrinterTestBase) AssertPrettyPrintNode(expectedJs string, ast rhino.Node) {
	AssertThat(b.PrettyPrintNode(ast)).IsEqualTo(expectedJs)
}

// Java: CodePrinterTestBase.assertPrint(String, String)
func (b *CodePrinterTestBase) AssertPrint(js string, expected string) {
	b.Parse(expected) // validate the expected string is valid JS
	AssertThat(
		jutil.StringTrim(
			b.ParsePrint(
				js,
				b.NewCompilerOptions(
					codePrinterTestBase_CompilerOptionBuilderFunc(func(options *CompilerOptions) {
						options.SetPrettyPrint(false)
						options.SetLineLengthThreshold(
							CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD)
					}))))).
		IsEqualTo(jutil.StringTrim(expected))
}

// Java: CodePrinterTestBase.assertPrintSame(String)
func (b *CodePrinterTestBase) AssertPrintSame(js string) {
	b.AssertPrint(js, js)
}

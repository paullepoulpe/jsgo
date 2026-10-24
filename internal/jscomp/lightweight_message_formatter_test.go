// Port of test/com/google/javascript/jscomp/LightweightMessageFormatterTest.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

var (
	jsaLmfFOO_TYPE             = DiagnosticType_Error("TEST_FOO", "error description here")
	jsaLmfORIGINAL_SOURCE_FILE = "original/source.html"
	jsaLmfORIGINAL_SOURCE      = sourcemap.Mapping_OriginalMapping_NewBuilder().
					SetOriginalFile(jsaLmfORIGINAL_SOURCE_FILE).
					SetLineNumber(3).
					SetColumnPosition(15).
					Build()
)

func jsaLmfFormat(region Region) *string {
	return NewLightweightMessageFormatter_LineNumberingFormatter().FormatRegion(region)
}

func jsaLmfRegion(startLine, endLine int, source string) Region {
	return NewSimpleRegion(startLine, endLine, source)
}

func jsaEqPtr(t *testing.T, got *string, want string) {
	t.Helper()
	if got == nil {
		t.Errorf("got nil, want %q", want)
		return
	}
	jsaEq(t, *got, want)
}

func TestLightweightMessageFormatter_Null(t *testing.T) {
	jsaTrue(t, jsaLmfFormat(nil) == nil)
}

func TestLightweightMessageFormatter_OneLineRegion(t *testing.T) {
	jsaEqPtr(t, jsaLmfFormat(jsaLmfRegion(5, 5, "hello world")), "  5| hello world")
}

func TestLightweightMessageFormatter_TwoLineRegion(t *testing.T) {
	jsaEqPtr(t, jsaLmfFormat(jsaLmfRegion(5, 6, "hello world\nfoo bar")),
		"  5| hello world\n"+
			"  6| foo bar")
}

func TestLightweightMessageFormatter_ThreeLineRegionAcrossNumberRange(t *testing.T) {
	region := jsaLmfFormat(jsaLmfRegion(9, 11, "hello world\nfoo bar\nanother one"))
	jsaEqPtr(t, region,
		"   9| hello world\n"+
			"  10| foo bar\n"+
			"  11| another one")
}

func TestLightweightMessageFormatter_ThreeLineRegionEmptyLine(t *testing.T) {
	region := jsaLmfFormat(jsaLmfRegion(7, 9, "hello world\n\nanother one"))
	jsaEqPtr(t, region,
		"  7| hello world\n"+
			"  8| \n"+
			"  9| another one")
}

func TestLightweightMessageFormatter_OnlyOneEmptyLine(t *testing.T) {
	jsaTrue(t, jsaLmfFormat(jsaLmfRegion(7, 7, "")) == nil)
}

func TestLightweightMessageFormatter_TwoEmptyLines(t *testing.T) {
	jsaEqPtr(t, jsaLmfFormat(jsaLmfRegion(7, 8, "\n")), "  7| ")
}

func TestLightweightMessageFormatter_ThreeLineRemoveLastEmptyLine(t *testing.T) {
	region := jsaLmfFormat(jsaLmfRegion(7, 9, "hello world\nfoobar\n"))
	jsaEqPtr(t, region,
		"  7| hello world\n"+
			"  8| foobar")
}

// Go-only: extra coverage — a region with more lines than its ending line
// number has digits for throws, as StringBuilder.repeat does.
func TestLightweightMessageFormatter_NegativeRepeatCount(t *testing.T) {
	e := jsaThrows[jutil.IIllegalArgumentException](t, func() { jsaLmfFormat(jsaLmfRegion(9, 9, "a\nb")) })
	if e != nil {
		jsaEq(t, e.Error(), "java.lang.IllegalArgumentException: count is negative: -1")
	}
}

func TestLightweightMessageFormatter_FormatErrorSpaces(t *testing.T) {
	n := rhino.Node_NewString("foobar").SetLinenoCharno(5, 8)
	n.SetLength(len("foobar"))
	n.SetSourceFileForTesting("javascript/complex.js")
	error := JSError_MakeWithNode(n, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatter("    if (foobar) {")
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:5:8: ERROR - [TEST_FOO] error description here\n"+
			"    if (foobar) {\n"+
			"        ^^^^^^\n")
}

func TestLightweightMessageFormatter_FormatErrorTabs(t *testing.T) {
	n := rhino.Node_NewString("foobar").SetLinenoCharno(5, 6)
	n.SetLength(len("foobar"))
	n.SetSourceFileForTesting("javascript/complex.js")
	error := JSError_MakeWithNode(n, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatter("\t\tif (foobar) {")
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:5:6: ERROR - [TEST_FOO] error description here\n"+
			"\t\tif (foobar) {\n"+
			"\t\t    ^^^^^^\n")
}

func TestLightweightMessageFormatter_FormatErrorSpaceEndOfLine1(t *testing.T) {
	error := JSError_MakeWithSourceLocation(jsaStr("javascript/complex.js"), 1, 10, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatter("assert (1;")
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:1:10: ERROR - [TEST_FOO] error description here\n"+
			"assert (1;\n"+
			"          ^\n")
}

func TestLightweightMessageFormatter_FormatErrorSpaceEndOfLine2(t *testing.T) {
	error := JSError_MakeWithSourceLocation(jsaStr("javascript/complex.js"), 6, 7, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatter("if (foo")
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:6:7: ERROR - [TEST_FOO] error description here\n"+
			"if (foo\n"+
			"       ^\n")
}

func TestLightweightMessageFormatter_FormatErrorOriginalSource(t *testing.T) {
	n := rhino.Node_NewString("foobar").SetLinenoCharno(5, 8)
	n.SetLength(len("foobar"))
	n.SetSourceFileForTesting("javascript/complex.js")
	error := JSError_MakeWithNode(n, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatterWithOriginal("    if (foobar) {", "<div ng-show='(foo'>")
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:5:8: \n"+
			"Originally at:\n"+
			"original/source.html:3:15: ERROR - [TEST_FOO] error description here\n"+
			"<div ng-show='(foo'>\n"+
			"               ^^^^^\n")
}

func TestLightweightMessageFormatter_Multiline_oneLineErrorMessage(t *testing.T) {
	n := rhino.Node_NewString("foobar").SetLinenoCharno(5, 8)
	n.SetLength(len("foobar"))
	n.SetSourceFileForTesting("javascript/complex.js")
	error := JSError_MakeWithNode(n, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatterFull("    if (foobar) {", SourceExcerptProvider_SourceExcerpt_FULL, 5)
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:5:8: ERROR - [TEST_FOO] error description here\n"+
			"  5|     if (foobar) {\n"+
			"             ^^^^^^\n")
}

func TestLightweightMessageFormatter_Multiline_twoLineErrorMessage(t *testing.T) {
	foobar := rhino.IR_String("foobar")
	baz := rhino.IR_String("baz")
	orNode := rhino.IR_Or(foobar, baz)
	orNode.SetLinenoCharno(5, 4)
	orNode.SetLength(len("foobar\n      || baz"))
	orNode.SetSourceFileForTesting("javascript/complex.js")

	error := JSError_MakeWithNode(orNode, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatterFull("if (foobar\n      || baz) {", SourceExcerptProvider_SourceExcerpt_FULL, 6)
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:5:4: ERROR - [TEST_FOO] error description here\n"+
			"  5| if (foobar\n"+
			"         ^^^^^^\n"+
			"  6|       || baz) {\n"+
			"     ^^^^^^^^^^^^\n")
}

func TestLightweightMessageFormatter_Multiline_charNoOutOfBoundsCrashes(t *testing.T) {
	n := rhino.Node_NewString("foobar").SetLinenoCharno(5, 800)
	n.SetLength(len("foobar"))
	n.SetSourceFileForTesting("javascript/complex.js")
	error := JSError_MakeWithNode(n, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatterFull("    if (foobar) {", SourceExcerptProvider_SourceExcerpt_FULL, 5)
	jsaTrue(t, strings.Contains(formatter.FormatError(error),
		"javascript/complex.js:5:800: ERROR - [TEST_FOO] error description here\n"+
			"Source excerpt could not be formatted. This may indicate a bug in the compiler."))
}

func TestLightweightMessageFormatter_Multiline_charNoOutOfBoundsCrashes2(t *testing.T) {
	n := rhino.Node_NewString("foobar").SetLinenoCharno(0, 2)
	n.SetLength(0)
	n.SetSourceFileForTesting("javascript/complex.js")
	error := JSError_MakeWithNode(n, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatterFull(
		// The first line is intentionally too short.
		"b\n"+
			"/** @export {(number|undefined)} */\n",
		SourceExcerptProvider_SourceExcerpt_FULL,
		5)

	jsaTrue(t, strings.Contains(formatter.FormatError(error),
		"javascript/complex.js: ERROR - [TEST_FOO] error description here\n"+
			"Source excerpt could not be formatted. This may indicate a bug in the compiler."))
}

func TestLightweightMessageFormatter_Multiline_longTruncatedErrorMessage(t *testing.T) {
	qname := rhino.IR_GetpropWithMoreProps(rhino.IR_Name("a"), "b", "c", "d", "e")
	qname.SetLinenoCharno(8, 0)
	qname.SetLength(len("a\n .b\n .c\n .d\n .e"))
	qname.SetSourceFileForTesting("javascript/complex.js")

	error := JSError_MakeWithNode(qname, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatterFull("a\n .b\n .c\n .d\n .e + 1;", SourceExcerptProvider_SourceExcerpt_FULL, 12)
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:8:0: ERROR - [TEST_FOO] error description here\n"+
			"   8| a\n"+
			"      ^\n"+
			"   9|  .b\n"+
			"      ^^^\n"+
			"...\n"+
			"  11|  .d\n"+
			"      ^^^\n"+
			"  12|  .e + 1;\n"+
			"      ^^^\n")
}

func TestLightweightMessageFormatter_Multiline_nodeLengthOutOfBounds(t *testing.T) {
	foobar := rhino.Node_NewString("foobar").SetLinenoCharno(5, 2)
	foobar.SetLength(len("foobar"))
	baz := rhino.Node_NewString("baz").SetLinenoCharno(6, 6)
	baz.SetLength(len("baz"))
	orNode := rhino.IR_Or(foobar, baz)
	orNode.SetLinenoCharno(5, 4)
	orNode.SetLength(1000) // intentionally too long
	orNode.SetSourceFileForTesting("javascript/complex.js")

	error := JSError_MakeWithNode(orNode, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatterFull("if (foobar", SourceExcerptProvider_SourceExcerpt_FULL, 5)
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:5:4: ERROR - [TEST_FOO] error description here\n"+
			"  5| if (foobar\n"+
			"         ^^^^^^\n")
}

func TestLightweightMessageFormatter_Multiline_errorWithoutAssociatedNode(t *testing.T) {
	error := JSError_MakeWithSourceLocation(jsaStr("javascript/complex.js"), 5, 4, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatterFull("if (foobar", SourceExcerptProvider_SourceExcerpt_FULL, 5)
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:5:4: ERROR - [TEST_FOO] error description here\n"+
			"  5| if (foobar\n"+
			"         ^\n")
}

func TestLightweightMessageFormatter_Multiline_errorWithoutAssociatedNodeAndNegativeChar(t *testing.T) {
	error := JSError_MakeWithSourceLocation(jsaStr("javascript/complex.js"), 5, -1, jsaLmfFOO_TYPE)
	formatter := jsaLmfFormatterFull("if (foobar", SourceExcerptProvider_SourceExcerpt_FULL, 5)
	jsaEq(t, formatter.FormatError(error),
		"javascript/complex.js:5: ERROR - [TEST_FOO] error description here\n"+
			"  5| if (foobar\n")
}

func jsaLmfFormatterFull(source string, excerpt SourceExcerptProvider_SourceExcerpt, endLine int) *LightweightMessageFormatter {
	return NewLightweightMessageFormatterWithFormat(jsaLmfSource(source, nil, endLine), excerpt)
}

func jsaLmfFormatter(source string) *LightweightMessageFormatter {
	return NewLightweightMessageFormatter(jsaLmfSource(source, nil, -1))
}

func jsaLmfFormatterWithOriginal(source string, originalSource string) *LightweightMessageFormatter {
	return NewLightweightMessageFormatter(jsaLmfSource(source, &originalSource, -1))
}

// jsaLmfSourceExcerptProvider is the anonymous SourceExcerptProvider of the
// Java test.
type jsaLmfSourceExcerptProvider struct {
	source         string
	originalSource *string
	endLineNumber  int
}

func jsaLmfSource(source string, originalSource *string, endLineNumber int) SourceExcerptProvider {
	return &jsaLmfSourceExcerptProvider{source: source, originalSource: originalSource, endLineNumber: endLineNumber}
}

func (p *jsaLmfSourceExcerptProvider) GetSourceLine(sourceName *string, lineNumber int) *string {
	if *sourceName == jsaLmfORIGINAL_SOURCE_FILE {
		return p.originalSource
	}
	return &p.source
}

func (p *jsaLmfSourceExcerptProvider) GetSourceLines(sourceName *string, lineNumber int, length int) Region {
	jutil.CheckStateMsg(p.endLineNumber != -1, "Must provide the end of the source lines")
	if *sourceName == jsaLmfORIGINAL_SOURCE_FILE {
		return NewSimpleRegion(lineNumber, p.endLineNumber, *p.originalSource)
	}
	return NewSimpleRegion(lineNumber, p.endLineNumber, p.source)
}

func (p *jsaLmfSourceExcerptProvider) GetSourceRegion(sourceName *string, lineNumber int) Region {
	panic(jutil.NewUnsupportedOperationException())
}

func (p *jsaLmfSourceExcerptProvider) GetSourceMapping(sourceName *string, lineNumber int, columnNumber int) *sourcemap.Mapping_OriginalMapping {
	if p.originalSource != nil {
		return jsaLmfORIGINAL_SOURCE
	}
	return nil
}

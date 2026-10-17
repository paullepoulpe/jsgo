// Port of test/com/google/javascript/jscomp/deps/JsFileLineParserTest.java @ closure-compiler v20260415

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestJsFileLineParser_SingleLine1(t *testing.T) {
	jsFileLineParserTest_assertStrip(t, "2", "// 1\n2\n")
}

func TestJsFileLineParser_SingleLine2(t *testing.T) {
	jsFileLineParserTest_assertStrip(t, "2 ", "// 1\n2 // 3 // 4\n")
}

func TestJsFileLineParser_MultiLine1(t *testing.T) {
	jsFileLineParserTest_assertStrip(t, "1", "/* hi */\n1\n")
}

func TestJsFileLineParser_MultiLine2(t *testing.T) {
	jsFileLineParserTest_assertStrip(t, "123", "1/* hi */2\n3\n")
}

func TestJsFileLineParser_MultiLine3(t *testing.T) {
	jsFileLineParserTest_assertStrip(t, "14", "1/* hi 2\n3*/4\n")
}

func TestJsFileLineParser_MultiLine4(t *testing.T) {
	jsFileLineParserTest_assertStrip(t, "15", "1/* hi x\ny\nz*/5\n")
}

func TestJsFileLineParser_MultiLine5(t *testing.T) {
	jsFileLineParserTest_assertStrip(t, "1234", "1/* hi */2/**/3/*\n/** bye */4\n")
}

func TestJsFileLineParser_MultiLine6(t *testing.T) {
	jsFileLineParserTest_assertStrip(t, "12", "1/*** hi *** 3 **/2\n")
}

func TestJsFileLineParser_MixedLine1(t *testing.T) {
	jsFileLineParserTest_assertStrip(t, "14", "1// /** 2 **/ 3\n4\n")
}

func TestJsFileLineParser_MixedLine2(t *testing.T) {
	jsFileLineParserTest_assertStrip(t, "1 34", "1/** // 2 **/ 3\n4\n")
}

func TestJsFileLineParser_BlockComment(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t, "/** one line */", "/** one line */")
}

func TestJsFileLineParser_BlockComment_singleQuotes(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t, "", "var x = '/** one line */';")
}

func TestJsFileLineParser_BlockComment_singleQuote(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t, "/** one line */", "'/** one line */")
}

func TestJsFileLineParser_BlockComment_singleQuotes_escaped(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t, "", "var x = '\\'/** one line */';")
}

func TestJsFileLineParser_BlockComment_singleQuotes_escapedSlash(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t, "", "var x = '\\'/** one line */\\\\';")
}

func TestJsFileLineParser_BlockComment_doubleQuotes(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t, "", "var x = \"/** one line */\";")
}

func TestJsFileLineParser_BlockComment_doubleQuote(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t, "/** one line */", "\"/** one line */;")
}

func TestJsFileLineParser_BlockComment_doubleQuotes_escaped(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t, "", "var x = \"\\\"/** one line */\";")
}

func TestJsFileLineParser_BlockComment_doubleQuotes_escapedSlash(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t, "", "var x = \"\\\"/** one line */\\\\\";")
}

func TestJsFileLineParser_InlineBlockComment(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t, "/** one line */", "var x; /** one line */")
	jsFileLineParserTest_assertBlocks(t, "/** one line */", "/** one line */ var y;")
	jsFileLineParserTest_assertBlocks(t, "/** one line */", "var x; /** one line */ var y;")
}

func TestJsFileLineParser_MultipleBlockComments(t *testing.T) {
	jsFileLineParserTest_assertBlocks(t,
		"/** first *//** * second */",
		"/** first */\n/**\n * second\n */\n")
}

func jsFileLineParserTest_assertStrip(t *testing.T, expected string, input string) {
	t.Helper()
	errorManager := NewPrintStreamErrorManager(jutil.System_Err())
	parser := newJsFileLineParserTest_TestParser(errorManager)
	parser.DoParse("file", jutil.NewStringReader(input))
	if got := parser.String(); got != expected {
		t.Errorf("strip %q: got %q, want %q", input, got, expected)
	}
}

func jsFileLineParserTest_assertBlocks(t *testing.T, expected string, input string) {
	t.Helper()
	errorManager := NewPrintStreamErrorManager(jutil.System_Err())
	parser := newJsFileLineParserTest_TestParser(errorManager)
	parser.DoParse("file", jutil.NewStringReader(input))
	if got := parser.comments.String(); got != expected {
		t.Errorf("blocks %q: got %q, want %q", input, got, expected)
	}
}

type jsFileLineParserTest_TestParser struct {
	JsFileLineParserBase
	sb       *jutil.StringBuilder
	comments *jutil.StringBuilder
}

func newJsFileLineParserTest_TestParser(errorManager ErrorManager) *jsFileLineParserTest_TestParser {
	p := &jsFileLineParserTest_TestParser{sb: jutil.NewStringBuilder(), comments: jutil.NewStringBuilder()}
	p.InitJsFileLineParserBase(p, errorManager)
	return p
}

func (p *jsFileLineParserTest_TestParser) ParseLine(line string) bool {
	p.sb.Append(line)
	return true
}

func (p *jsFileLineParserTest_TestParser) ParseJsDocCommentLine(line string) bool {
	p.comments.Append(line)
	return true
}

func (p *jsFileLineParserTest_TestParser) String() string {
	return p.sb.String()
}

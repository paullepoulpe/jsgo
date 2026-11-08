// Port of com/google/javascript/jscomp/parsing/JsDocInfoParserTest.java @ closure-compiler v20261104
// (test/com/google/javascript/jscomp/parsing/JsDocInfoParserTest.java): testJsDocInfoPosition.

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

// Go-only: test-helper — JsDocInfoParserTest's one test that reads
// positions through com.google.javascript.jscomp.SourceFile. jscomp imports
// parsing, so the test is in the external test package parsing_test; the
// rest of the class is js_doc_info_parser_test.go.
package parsing_test

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

func parseFullForPositionTest(t *testing.T, code string) rhino.Node {
	t.Helper()
	defaultCfg := parsing.ParserRunner_CreateConfigWithModes(
		parsing.Config_LanguageMode_ECMASCRIPT3,
		parsing.Config_JsDocParsing_TYPES_ONLY,
		parsing.Config_RunMode_STOP_AFTER_ERROR,
		nil,
		true,
		parsing.Config_StrictMode_SLOPPY)
	extraAnnotations := jutil.NewOrderedSet(defaultCfg.AnnotationNames()...)
	extraSuppressions := jutil.NewOrderedSet(parsing.ParserRunner_GetSuppressionNames().ToSlice()...)
	extraSuppressions.Add("x")
	extraSuppressions.Add("y")
	extraSuppressions.Add("z")

	testErrorReporter := rhinotesting.NewTestErrorReporter()
	config := parsing.NewConfigBuilder().
		SetExtraAnnotationNames(extraAnnotations.All()).
		SetJsDocParsingMode(parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE).
		SetRunMode(parsing.Config_RunMode_KEEP_GOING).
		SetSuppressionNames(extraSuppressions.All()).
		SetLanguageMode(parsing.Config_LanguageMode_ECMASCRIPT3).
		SetParseInlineSourceMaps(true).
		SetStrictMode(parsing.Config_StrictMode_SLOPPY).
		Build()

	result := parsing.ParserRunner_Parse(
		rhino.NewSimpleSourceFile("source", rhino.StaticSourceFile_SourceKind_STRONG),
		code,
		config,
		testErrorReporter)
	testErrorReporter.VerifyHasEncounteredAllWarningsAndErrors(t)
	return result.Ast
}

// Java: JsDocInfoParserTest.testJsDocInfoPosition()
func TestJsDocInfoPosition(t *testing.T) {
	truth := rhinotesting.NewTruth(t)
	sourceFile :=
		jscomp.SourceFile_FromCode(
			"comment-position-test.js",
			"   \n  /**\n   * A comment\n   */\n  function f(x) {}\n")
	script := parseFullForPositionTest(t, sourceFile.GetCode())
	jutil.CheckState(script.IsScript())
	fn := script.GetFirstChild()
	jutil.CheckState(fn.IsFunction())
	jsdoc := fn.GetJSDocInfo()

	truth.AssertThat(jsdoc.GetOriginalCommentPosition()).IsEqualTo(6)
	truth.AssertThat(sourceFile.GetLineOfOffset(jsdoc.GetOriginalCommentPosition())).IsEqualTo(2)
	truth.AssertThat(sourceFile.GetColumnOfOffset(jsdoc.GetOriginalCommentPosition())).IsEqualTo(2)
}

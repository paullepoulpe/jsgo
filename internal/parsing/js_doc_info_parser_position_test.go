// Port of com/google/javascript/jscomp/parsing/JsDocInfoParserTest.java @ closure-compiler v20261021
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
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// Java: JsDocInfoParserTest.testJsDocInfoPosition()
func TestJsDocInfoPosition(t *testing.T) {
	truth := rhinotesting.NewTruth(t)
	sourceFile :=
		jscomp.SourceFile_FromCode(
			"comment-position-test.js",
			"   \n  /**\n   * A comment\n   */\n  function f(x) {}\n")
	script := parsing.JsDocInfoParserTest_ParseFull(t, sourceFile.GetCode())
	jutil.CheckState(script.IsScript())
	fn := script.GetFirstChild()
	jutil.CheckState(fn.IsFunction())
	jsdoc := fn.GetJSDocInfo()

	truth.AssertThat(jsdoc.GetOriginalCommentPosition()).IsEqualTo(6)
	truth.AssertThat(sourceFile.GetLineOfOffset(jsdoc.GetOriginalCommentPosition())).IsEqualTo(2)
	truth.AssertThat(sourceFile.GetColumnOfOffset(jsdoc.GetOriginalCommentPosition())).IsEqualTo(2)
}

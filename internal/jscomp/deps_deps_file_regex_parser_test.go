// Port of com/google/javascript/jscomp/deps/DepsFileRegexParserTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/deps/DepsFileRegexParserTest.java)

/*
 * Copyright 2008 The Closure Compiler Authors.
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

// Tests for DepsFileRegexParser.
//
// Go-only: test-helper — the JUnit class is a plain Go test per @Test
// method; setUp is depsFileRegexParserTest_setUp.

const depsFileRegexParserTest_SRC_PATH = "/path/1.js"

// Java: DepsFileRegexParserTest.setUp()
func depsFileRegexParserTest_setUp() (*DepsFileRegexParser, ErrorManager) {
	errorManager := NewPrintStreamErrorManager(jutil.System_Err())
	parser := NewDepsFileRegexParser(errorManager)
	parser.SetShortcutMode(true)
	return parser, errorManager
}

// depsFileRegexParserTest_assertDepsEqual is
// assertThat(result).isEqualTo(expected) for lists of DependencyInfo.
//
// Go-only: test-helper — List.equals over DependencyInfo.equals.
func depsFileRegexParserTest_assertDepsEqual(t *testing.T, result []DependencyInfo, expected []DependencyInfo) {
	t.Helper()
	equal := len(result) == len(expected)
	for i := 0; equal && i < len(result); i++ {
		equal = result[i].(*SimpleDependencyInfo).Equals(expected[i])
	}
	if !equal {
		t.Fatalf("expected: %v\nbut was : %v", expected, result)
	}
}

func depsFileRegexParserTest_assertCounts(t *testing.T, errorManager ErrorManager, errors int, warnings int) {
	t.Helper()
	if got := errorManager.GetErrorCount(); got != errors {
		t.Errorf("error count: expected %d, got %d", errors, got)
	}
	if got := errorManager.GetWarningCount(); got != warnings {
		t.Errorf("warning count: expected %d, got %d", warnings, got)
	}
}

// TestDepsFileRegexParser_GoodParse tests:
//
//   - Parsing of comments,
//   - Parsing of different styles of quotes,
//   - Parsing of empty arrays,
//   - Parsing of non-empty arrays,
//   - Correct recording of what was parsed.
//
// Java: DepsFileRegexParserTest.testGoodParse()
func TestDepsFileRegexParser_GoodParse(t *testing.T) {
	parser, errorManager := depsFileRegexParserTest_setUp()
	contents :=
		"/*goog.addDependency('no1', [], []);*//*\n" +
			"goog.addDependency('no2', [ ], [ ]);\n" +
			"*/goog.addDependency('yes1', [], []);\n" +
			"/* blah */goog.addDependency(\"yes2\", [], [])/* blah*/\n" +
			"goog.addDependency('yes3', ['a','b'], ['c']); // goog.addDependency('no3', [], []);\n" +
			"// goog.addDependency('no4', [], []);\n" +
			"goog.addDependency(\"yes4\", [], [ \"a\",'b' , 'c' ]); //no new line at EOF\n"

	result := parser.ParseFileWithContents(depsFileRegexParserTest_SRC_PATH, contents)
	expected := []DependencyInfo{
		SimpleDependencyInfo_Builder_("yes1", depsFileRegexParserTest_SRC_PATH).Build(),
		SimpleDependencyInfo_Builder_("yes2", depsFileRegexParserTest_SRC_PATH).Build(),
		SimpleDependencyInfo_Builder_("yes3", depsFileRegexParserTest_SRC_PATH).
			SetProvides([]string{"a", "b"}).
			SetRequiresArray(DependencyInfo_Require_ParsedFromDeps("c")).
			Build(),
		SimpleDependencyInfo_Builder_("yes4", depsFileRegexParserTest_SRC_PATH).
			SetRequiresArray(
				DependencyInfo_Require_ParsedFromDeps("a"),
				DependencyInfo_Require_ParsedFromDeps("b"),
				DependencyInfo_Require_ParsedFromDeps("c")).
			Build(),
	}

	depsFileRegexParserTest_assertDepsEqual(t, result, expected)
	depsFileRegexParserTest_assertCounts(t, errorManager, 0, 0)
}

// Java: DepsFileRegexParserTest.testTooFewArgs()
func TestDepsFileRegexParser_TooFewArgs(t *testing.T) {
	parser, errorManager := depsFileRegexParserTest_setUp()
	parser.ParseFileWithContents(depsFileRegexParserTest_SRC_PATH, "goog.addDependency('a', []);")
	depsFileRegexParserTest_assertCounts(t, errorManager, 1, 0)
}

// Java: DepsFileRegexParserTest.testTooManyArgs1()
func TestDepsFileRegexParser_TooManyArgs1(t *testing.T) {
	parser, errorManager := depsFileRegexParserTest_setUp()
	parser.ParseFileWithContents(depsFileRegexParserTest_SRC_PATH, "goog.addDependency('a', [], [], []);")
	depsFileRegexParserTest_assertCounts(t, errorManager, 1, 0)
}

// Java: DepsFileRegexParserTest.testTooManyArgs2()
func TestDepsFileRegexParser_TooManyArgs2(t *testing.T) {
	parser, errorManager := depsFileRegexParserTest_setUp()
	parser.ParseFileWithContents(depsFileRegexParserTest_SRC_PATH, "goog.addDependency('a', [], [], false, []);")
	depsFileRegexParserTest_assertCounts(t, errorManager, 1, 0)
}

// Java: DepsFileRegexParserTest.testTooManyArgs3()
func TestDepsFileRegexParser_TooManyArgs3(t *testing.T) {
	parser, errorManager := depsFileRegexParserTest_setUp()
	parser.ParseFileWithContents(depsFileRegexParserTest_SRC_PATH, "goog.addDependency('a', [], [], {}, []);")
	depsFileRegexParserTest_assertCounts(t, errorManager, 1, 0)
}

// Java: DepsFileRegexParserTest.testBadLoadFlagsSyntax()
func TestDepsFileRegexParser_BadLoadFlagsSyntax(t *testing.T) {
	parser, errorManager := depsFileRegexParserTest_setUp()
	parser.ParseFileWithContents(depsFileRegexParserTest_SRC_PATH, "goog.addDependency('a', [], [], {module: 'goog'});")
	depsFileRegexParserTest_assertCounts(t, errorManager, 1, 0)
}

// Java: DepsFileRegexParserTest.testGoogModule()
func TestDepsFileRegexParser_GoogModule(t *testing.T) {
	parser, _ := depsFileRegexParserTest_setUp()
	result :=
		parser.ParseFileWithContents(
			depsFileRegexParserTest_SRC_PATH,
			"goog.addDependency('yes1', [], [], true);\n"+
				"goog.addDependency('yes2', [], [], false);\n")
	expected := []DependencyInfo{
		SimpleDependencyInfo_Builder_("yes1", depsFileRegexParserTest_SRC_PATH).SetGoogModule(true).Build(),
		SimpleDependencyInfo_Builder_("yes2", depsFileRegexParserTest_SRC_PATH).Build(),
	}
	depsFileRegexParserTest_assertDepsEqual(t, result, expected)
}

func depsFileRegexParserTest_mapOf(k string, v string) *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	m.Put(k, v)
	return m
}

// Java: DepsFileRegexParserTest.testEs6Module()
func TestDepsFileRegexParser_Es6Module(t *testing.T) {
	parser, _ := depsFileRegexParserTest_setUp()
	result :=
		parser.ParseFileWithContents(
			depsFileRegexParserTest_SRC_PATH,
			"goog.addDependency('path/from/closure.js', [], ['nexttoclosure.js'], {'module':'es6'});\n"+
				"goog.addDependency('nexttoclosure.js', [], [], {'module':'es6'});\n")
	expected := []DependencyInfo{
		SimpleDependencyInfo_Builder_("path/from/closure.js", depsFileRegexParserTest_SRC_PATH).
			SetLoadFlags(depsFileRegexParserTest_mapOf("module", "es6")).
			SetProvides([]string{"path/from/closure.js"}).
			SetRequiresArray(DependencyInfo_Require_ParsedFromDeps("nexttoclosure.js")).
			Build(),
		SimpleDependencyInfo_Builder_("nexttoclosure.js", depsFileRegexParserTest_SRC_PATH).
			SetLoadFlags(depsFileRegexParserTest_mapOf("module", "es6")).
			SetProvides([]string{"nexttoclosure.js"}).
			Build(),
	}
	depsFileRegexParserTest_assertDepsEqual(t, result, expected)
}

// Java: DepsFileRegexParserTest.testLoadFlags()
func TestDepsFileRegexParser_LoadFlags(t *testing.T) {
	parser, _ := depsFileRegexParserTest_setUp()
	result :=
		parser.ParseFileWithContents(
			depsFileRegexParserTest_SRC_PATH,
			"goog.addDependency('yes1', [], [], {'module': 'goog'});\n"+
				"goog.addDependency('yes2', [], [], {\"lang\": \"es6\"});\n"+
				"goog.addDependency('yes3', [], [], {});\n")
	expected := []DependencyInfo{
		SimpleDependencyInfo_Builder_("yes1", depsFileRegexParserTest_SRC_PATH).
			SetLoadFlags(depsFileRegexParserTest_mapOf("module", "goog")).
			Build(),
		SimpleDependencyInfo_Builder_("yes2", depsFileRegexParserTest_SRC_PATH).
			SetLoadFlags(depsFileRegexParserTest_mapOf("lang", "es6")).
			Build(),
		SimpleDependencyInfo_Builder_("yes3", depsFileRegexParserTest_SRC_PATH).Build(),
	}
	depsFileRegexParserTest_assertDepsEqual(t, result, expected)
}

// Java: DepsFileRegexParserTest.testShortcutMode()
func TestDepsFileRegexParser_ShortcutMode(t *testing.T) {
	parser, _ := depsFileRegexParserTest_setUp()
	result :=
		parser.ParseFileWithContents(
			depsFileRegexParserTest_SRC_PATH,
			"goog.addDependency('yes1', [], []); \n"+
				"foo();\n"+
				"goog.addDependency('no1', [], []);\n")
	expected := []DependencyInfo{SimpleDependencyInfo_Builder_("yes1", depsFileRegexParserTest_SRC_PATH).Build()}
	depsFileRegexParserTest_assertDepsEqual(t, result, expected)
}

// Java: DepsFileRegexParserTest.testNoShortcutMode()
func TestDepsFileRegexParser_NoShortcutMode(t *testing.T) {
	parser, _ := depsFileRegexParserTest_setUp()
	parser.SetShortcutMode(false)
	result :=
		parser.ParseFileWithContents(
			depsFileRegexParserTest_SRC_PATH,
			"goog.addDependency('yes1', [], []); \n"+
				"foo();\n"+
				"goog.addDependency('yes2', [], []);\n")
	expected := []DependencyInfo{
		SimpleDependencyInfo_Builder_("yes1", depsFileRegexParserTest_SRC_PATH).Build(),
		SimpleDependencyInfo_Builder_("yes2", depsFileRegexParserTest_SRC_PATH).Build(),
	}
	depsFileRegexParserTest_assertDepsEqual(t, result, expected)
}

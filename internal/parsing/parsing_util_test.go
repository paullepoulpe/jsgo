// Port of com/google/javascript/jscomp/parsing/ParsingUtilTest.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
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

package parsing

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// parsingUtilTest_NAME_EQUALITY is the Correspondence NAME_EQUALITY: a node
// corresponds to the string of a NAME node (anything else to null).
func parsingUtilTest_NAME_EQUALITY(n rhino.Node) *string {
	if n.IsName() {
		s := n.GetString()
		return &s
	}
	return nil
}

// Java: ParsingUtilTest.testNamePatternDeclaresName()
func TestNamePatternDeclaresName(t *testing.T) {
	parsingUtilTest_assertPatternDeclaresNames(t, "x", setOf("x"))

	// assertThrows(AssertionError.class, ...)
	if err := parsingUtilTest_patternDeclaresNames(t, "x", setOf("y")); err == nil {
		t.Errorf("expected an AssertionError")
	}
}

// Java: ParsingUtilTest.testEmptyPatternDeclaresNoNames()
func TestEmptyPatternDeclaresNoNames(t *testing.T) {
	parsingUtilTest_assertPatternDeclaresNames(t, "[]", setOf())
	parsingUtilTest_assertPatternDeclaresNames(t, "{}", setOf())
	parsingUtilTest_assertPatternDeclaresNames(t, "[,]", setOf())
}

// Java: ParsingUtilTest.testNamesDeclaredInArrayPattern()
func TestNamesDeclaredInArrayPattern(t *testing.T) {
	parsingUtilTest_assertPatternDeclaresNames(t, "[x]", setOf("x"))
	parsingUtilTest_assertPatternDeclaresNames(t, "[x, y]", setOf("x", "y"))
	parsingUtilTest_assertPatternDeclaresNames(t, "[x, y, ...z]", setOf("x", "y", "z"))
	parsingUtilTest_assertPatternDeclaresNames(t, "[...z]", setOf("z"))
	parsingUtilTest_assertPatternDeclaresNames(t, "[[x, y]]", setOf("x", "y"))
	parsingUtilTest_assertPatternDeclaresNames(t, "[...[x, y]]", setOf("x", "y"))
	parsingUtilTest_assertPatternDeclaresNames(t, "[[x], [y]]", setOf("x", "y"))
}

// Java: ParsingUtilTest.testNamesDeclaredInObjectPattern()
func TestNamesDeclaredInObjectPattern(t *testing.T) {
	parsingUtilTest_assertPatternDeclaresNames(t, "{x}", setOf("x"))
	parsingUtilTest_assertPatternDeclaresNames(t, "{x, y}", setOf("x", "y"))
	parsingUtilTest_assertPatternDeclaresNames(t, "{x, y, ...z}", setOf("x", "y", "z"))
	parsingUtilTest_assertPatternDeclaresNames(t, "{...z}", setOf("z"))
	parsingUtilTest_assertPatternDeclaresNames(t, "{x: y}", setOf("y"))
	parsingUtilTest_assertPatternDeclaresNames(t, "{[x]: y}", setOf("y"))
	parsingUtilTest_assertPatternDeclaresNames(t, "{[x()]: y}", setOf("y"))
}

// Java: ParsingUtilTest.testNamesDeclaredInDefaultValue()
func TestNamesDeclaredInDefaultValue(t *testing.T) {
	parsingUtilTest_assertPatternDeclaresNames(t, "x = 1", setOf("x"))
	parsingUtilTest_assertPatternDeclaresNames(t, "[x, y] = [1, 2]", setOf("x", "y"))
	parsingUtilTest_assertPatternDeclaresNames(t, "{x, y} = {x: 1, y: 2}", setOf("x", "y"))
}

// Java: ParsingUtilTest.testNamesDeclaredInRest()
func TestNamesDeclaredInRest(t *testing.T) {
	parsingUtilTest_assertPatternDeclaresNames(t, "...x", setOf("x"))
	parsingUtilTest_assertPatternDeclaresNames(t, "...[x, y, z]", setOf("x", "y", "z"))
}

// Java: ParsingUtilTest.testNamesDeclaredInParamList()
func TestNamesDeclaredInParamList(t *testing.T) {
	parsingUtilTest_assertPatternDeclaresNames(t, "x, y, z", setOf("x", "y", "z"))
	parsingUtilTest_assertPatternDeclaresNames(t,
		"[x1, x2], {y1, y2}, z = 0", setOf("x1", "x2", "y1", "y2", "z"))
}

// Java: ParsingUtilTest.assertPatternDeclaresNames(String, Set<String>)
func parsingUtilTest_assertPatternDeclaresNames(t *testing.T, pattern string, expected *jutil.LinkedHashSet[string]) {
	t.Helper()
	if err := parsingUtilTest_patternDeclaresNames(t, pattern, expected); err != nil {
		t.Fatal(err)
	}
}

// parsingUtilTest_patternDeclaresNames is the body of
// assertPatternDeclaresNames, returning the AssertionError of
// assertThat(seen).comparingElementsUsing(NAME_EQUALITY).containsExactlyElementsIn(expected)
// instead of throwing it (so that assertThrows can be expressed).
func parsingUtilTest_patternDeclaresNames(t *testing.T, pattern string, expected *jutil.LinkedHashSet[string]) error {
	t.Helper()
	seen := jutil.NewLinkedHashSet[rhino.Node]()
	ParsingUtil_GetParamOrPatternNames(parsingUtilTest_parsePattern(t, pattern), func(n rhino.Node) { seen.Add(n) })

	// containsExactlyElementsIn: every expected element corresponds to a
	// distinct actual element, and no actual element is left over.
	actual := seen.ToSlice()
	used := make([]bool, len(actual))
	missing := []string{}
	for e := range expected.All() {
		found := false
		for i, a := range actual {
			if name := parsingUtilTest_NAME_EQUALITY(a); !used[i] && name != nil && *name == e {
				used[i], found = true, true
				break
			}
		}
		if !found {
			missing = append(missing, e)
		}
	}
	unexpected := []string{}
	for i, a := range actual {
		if !used[i] {
			unexpected = append(unexpected, a.String())
		}
	}
	if len(missing) > 0 || len(unexpected) > 0 {
		return fmt.Errorf("not true that %v contains exactly elements that is NAME with string %v: missing %v, unexpected %v",
			actual, expected, missing, unexpected)
	}
	return nil
}

// parsingUtilTest_parsePattern inserts pattern into a functions parameter
// list and returns the Node that pattern parses to.
//
// Java: ParsingUtilTest.parsePattern(String)
func parsingUtilTest_parsePattern(t *testing.T, pattern string) rhino.Node {
	t.Helper()
	script := parsingUtilTest_parse(t, "function foo("+pattern+"){}")
	assertThat(t, script.IsScript()).IsTrue()
	function := script.GetFirstChild()
	assertThat(t, function.IsFunction()).IsTrue()
	paramList := function.GetSecondChild()
	assertThat(t, paramList.IsParamList()).IsTrue()
	return paramList
}

// Java: ParsingUtilTest.parse(String, String...)
func parsingUtilTest_parse(t *testing.T, source string, warnings ...string) rhino.Node {
	t.Helper()
	testErrorReporter := newTestErrorReporter().expectAllWarnings(warnings...)
	config := ParserRunner_CreateConfigWithModes(
		Config_LanguageMode_ES_NEXT,
		Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE,
		Config_RunMode_KEEP_GOING,
		nil,
		true,
		Config_StrictMode_STRICT)
	script := ParserRunner_Parse(
		rhino.NewSimpleSourceFile("input", rhino.StaticSourceFile_SourceKind_STRONG), source, config, testErrorReporter).
		Ast

	// verifying that all warnings were seen
	testErrorReporter.verifyHasEncounteredAllWarningsAndErrors(t)

	return script
}

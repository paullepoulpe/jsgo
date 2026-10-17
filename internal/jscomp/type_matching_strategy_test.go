// Port of com/google/javascript/jscomp/TypeMatchingStrategyTest.java @ closure-compiler v20260415

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// Java: TypeMatchingStrategyTest.EXTERNS
const typeMatchingStrategyTest_EXTERNS = `/** @constructor */
var OtherType = function() {};
/** @constructor */
var SuperType = function() {};
/** @constructor @extends {SuperType} */
var SubType = function() {};
`

func TestTypeMatchingStrategy_Match_default(t *testing.T) {
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "!SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "?SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "SuperType|undefined", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "SuperType|void", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "!SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "?SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "SubType|undefined", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "SubType|void", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "number", false, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "*", true, true)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "?", true, true)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "?", "string", true, false)
}

func TestTypeMatchingStrategy_Match_respectNullability(t *testing.T) {
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "!SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "?SuperType", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "SuperType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "SuperType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "!SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "?SubType", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "SubType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "SubType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "number", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "*", true, true)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "?", true, true)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "?", "string", true, false)
}

func TestTypeMatchingStrategy_Match_subtypes(t *testing.T) {
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "!SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "?SuperType", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "!OtherType", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "SuperType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "SuperType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "!SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "?SubType", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "SubType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "SubType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "number", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "*", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "?", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "!SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "?SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "?OtherType", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "SuperType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "SuperType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "SuperType|OtherType", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "!SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "?SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "SubType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "SubType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "number", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "*", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "?", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?", "string", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?", "?", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?", "*", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "*", "?", false, false)
}

func TestTypeMatchingStrategy_Match_exact(t *testing.T) {
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "!SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "?SuperType", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "SuperType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "SuperType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "!SubType", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "?SubType", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "SubType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "SubType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "number", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "*", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "?", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "?", "string", true, false)
}

func assertMatch(
	t *testing.T,
	typeMatchingStrategy TypeMatchingStrategy,
	templateType string,
	type_ string,
	isMatch bool,
	isLooseMatch bool) {
	t.Helper()
	truth := rhinotesting.NewTruth(t)

	// It's important that the test uses the same compiler to compile the template type and the
	// type to be matched. Otherwise, equal types won't be considered equal.
	compiler := NewCompiler()
	compiler.DisableThreads()
	options := NewCompilerOptions()
	options.SetCheckTypes(true)
	options.SetChecksOnly(true)

	compiler.Compile(
		[]*SourceFile{SourceFile_FromCode("externs", typeMatchingStrategyTest_EXTERNS)},
		[]*SourceFile{
			SourceFile_FromCode(
				"test",
				jutil.StringFormat(
					"/** @type {%s} */ var x; /** @type {%s} */ var y;", templateType, type_)),
		},
		options)
	script := compiler.GetRoot().GetLastChild().GetFirstChild()
	xNode := script.GetFirstChild()
	yNode := script.GetLastChild()
	templateJsType := xNode.GetFirstChild().GetJSType()
	jsType := yNode.GetFirstChild().GetJSType()

	matchResult := typeMatchingStrategy.Match(templateJsType, jsType)
	if isMatch {
		truth.AssertWithMessage("%s", "'"+templateJsType.String()+"' should match '"+jsType.String()+"'").
			That(matchResult.IsMatch()).IsEqualTo(isMatch)
	} else {
		truth.AssertWithMessage("%s", "'"+templateType+"' should not match '"+type_+"'").
			That(matchResult.IsMatch()).IsEqualTo(isMatch)
	}
	if isLooseMatch {
		truth.AssertWithMessage("%s", "'"+templateType+"' should loosely match '"+type_+"'").
			That(matchResult.IsLooseMatch()).IsEqualTo(isLooseMatch)
	} else {
		truth.AssertWithMessage("%s", "'"+templateType+"' should not loosely match '"+type_+"'").
			That(matchResult.IsLooseMatch()).IsEqualTo(isLooseMatch)
	}
}

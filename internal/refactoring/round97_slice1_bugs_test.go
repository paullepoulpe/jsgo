/*
 * Copyright 2026 The Closure Compiler Authors.
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

package refactoring

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func newExtraRequireFixTest(t *testing.T) *ErrorToFixMapperTest {
	t.Helper()
	c := NewErrorToFixMapperTest()
	c.JUnitSetT(t)
	c.JUnitBefore()
	c.options.SetCheckTypes(false)
	return c
}

func TestRound97ExtraRequireObjectRest(t *testing.T) {
	t.Run("UnusedObjectRestAfterUsedMember", func(t *testing.T) {
		c := newExtraRequireFixTest(t)

		c.AssertChanges(
			`goog.module('x');

const { Coordinate, ...rest } = goog.require('goog.math');

alert(Coordinate);
`,
			`goog.module('x');

const { Coordinate } = goog.require('goog.math');

alert(Coordinate);
`)
	})

	t.Run("SoleUnusedObjectRestRemovesDeclarationInTwoPasses", func(t *testing.T) {
		c := newExtraRequireFixTest(t)

		pass0 := `goog.module('x');

const {...rest} = goog.require('goog.math');
`
		pass1 := `goog.module('x');

const {} = goog.require('goog.math');
`
		pass2 := `goog.module('x');


`
		c.AssertChanges(pass0, pass1)

		c2 := newExtraRequireFixTest(t)
		c2.AssertChanges(pass1, pass2)
	})

	t.Run("UnusedObjectRestInRequireType", func(t *testing.T) {
		c := newExtraRequireFixTest(t)

		c.AssertChanges(
			`goog.module('x');

const { Coordinate, ...rest } = goog.requireType('goog.math');

/** @type {!Coordinate} */ let pt;
`,
			`goog.module('x');

const { Coordinate } = goog.requireType('goog.math');

/** @type {!Coordinate} */ let pt;
`)
	})
}

func TestRound97SuggestedFixDeleteObjectRest(t *testing.T) {
	t.Run("TrailingObjectRest", func(t *testing.T) {
		input := "const { Coordinate, ...rest } = obj;"
		compiler := suggestedFixTest_getCompiler(input)
		root := suggestedFixTest_compileToScriptRoot(compiler)
		objPat := root.GetFirstFirstChild().GetFirstChild()
		restProp := objPat.GetLastChild()

		fix := NewSuggestedFix_Builder().Delete(restProp).Build()
		wantReplacement := CodeReplacement_Create(jutil.StringIndexOf(input, ", ...rest"), jutil.StringLength(", ...rest"), "")
		SuggestedFixes_AssertReplacement(fix, wantReplacement)

		got := ApplyToCode([]*SuggestedFix{fix}, map[string]string{"test": input})["test"]
		if want := "const { Coordinate } = obj;"; got != want {
			t.Fatalf("ApplyToCode = %q, want %q", got, want)
		}
	})

	t.Run("SoleObjectRest", func(t *testing.T) {
		input := "const {...rest} = obj;"
		compiler := suggestedFixTest_getCompiler(input)
		root := suggestedFixTest_compileToScriptRoot(compiler)
		objPat := root.GetFirstFirstChild().GetFirstChild()
		restProp := objPat.GetFirstChild()

		fix := NewSuggestedFix_Builder().Delete(restProp).Build()
		wantReplacement := CodeReplacement_Create(jutil.StringIndexOf(input, "...rest"), jutil.StringLength("...rest"), "")
		SuggestedFixes_AssertReplacement(fix, wantReplacement)

		got := ApplyToCode([]*SuggestedFix{fix}, map[string]string{"test": input})["test"]
		if want := "const {} = obj;"; got != want {
			t.Fatalf("ApplyToCode = %q, want %q", got, want)
		}
	})
}

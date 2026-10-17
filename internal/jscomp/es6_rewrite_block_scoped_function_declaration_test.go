// Port of com/google/javascript/jscomp/Es6RewriteBlockScopedFunctionDeclarationTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/Es6RewriteBlockScopedFunctionDeclarationTest.java)

/*
 * Copyright 2014 The Closure Compiler Authors.
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
)

type Es6RewriteBlockScopedFunctionDeclarationTest struct {
	CompilerTestCaseBase
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest()
func NewEs6RewriteBlockScopedFunctionDeclarationTest() *Es6RewriteBlockScopedFunctionDeclarationTest {
	c := &Es6RewriteBlockScopedFunctionDeclarationTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestEs6RewriteBlockScopedFunctionDeclarationTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteBlockScopedFunctionDeclarationTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest.customSetUp()
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) CustomSetUp() {
	c.EnableNormalize()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest.getOptions()
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	return options
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest.getProcessor(Compiler)
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs6RewriteBlockScopedFunctionDeclaration(compiler)
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest.testRewritesBlockScopedFunctionDeclaration()
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) TestRewritesBlockScopedFunctionDeclaration() {
	c.Test("{ function f(){} }", "{ let f = function(){}; }")
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest.testHoistsFunctionToStartOfBlock()
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) TestHoistsFunctionToStartOfBlock() {
	c.Test(
		` // preserve newlines
function use(x) {}
{ use(f()); function f(){} }
`,
		` // preserve newlines
function use(x) {}
{ let f = function(){}; use(f()); }
`)
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest.testBlockScopedGeneratorFunction()
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) TestBlockScopedGeneratorFunction() {
	c.Test("{ function* f() {yield 1;} }", "{ let f = function*() { yield 1; }; }")
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest.testBlockNestedInsideFunction()
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) TestBlockNestedInsideFunction() {
	c.Test(
		`function f() {
  var x = 1;
  if (a) {
    x1();
    function x1() { return x1; }
  }
  return x;
}
`,
		`function f() {
  var x = 1;
  if (a) {
    let x1 = function() { return x1; };
    x1();
  }
  return x;
}
`)
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest.testFunctionInLoop()
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) TestFunctionInLoop() {
	c.Test(
		`for (var x of y) {
  y();
  function f() {
    let z;
  }
}
`,
		`var x;
for (x of y) {
  let f = function() {
    let z;
  };
  y();
}
`)
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest.testDoesNotRewriteTopLevelDeclarations()
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) TestDoesNotRewriteTopLevelDeclarations() {
	c.TestSame("function f(){}")
}

// Java: Es6RewriteBlockScopedFunctionDeclarationTest.testDoesNotRewriteFunctionScopedDeclarations()
func (c *Es6RewriteBlockScopedFunctionDeclarationTest) TestDoesNotRewriteFunctionScopedDeclarations() {
	c.TestSame("function g() {function f(){}}")
}

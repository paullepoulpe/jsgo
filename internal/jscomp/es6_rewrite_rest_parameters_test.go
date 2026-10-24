// Port of com/google/javascript/jscomp/Es6RewriteRestParametersTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/Es6RewriteRestParametersTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteRestParametersTest has the test cases for the transpilation of
// rest parameters.
type Es6RewriteRestParametersTest struct {
	CompilerTestCaseBase
}

// Java: Es6RewriteRestParametersTest.EXTERNS_BASE
var Es6RewriteRestParametersTest_EXTERNS_BASE = NewTestExternsBuilder().
	AddFunction().
	AddJSCompLibraries().
	AddExtra("$jscomp.getRestArguments = function(argument) {};").
	Build()

// Java: Es6RewriteRestParametersTest()
func NewEs6RewriteRestParametersTest() *Es6RewriteRestParametersTest {
	c := &Es6RewriteRestParametersTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, Es6RewriteRestParametersTest_EXTERNS_BASE)
	return c
}

func TestEs6RewriteRestParametersTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteRestParametersTest() })
}

// Java: Es6RewriteRestParametersTest.getProcessor(Compiler)
func (c *Es6RewriteRestParametersTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		NewInjectTranspilationRuntimeLibraries(compiler).Process(externs, root)
		NewEs6RewriteRestParameters(compiler).Process(externs, root)
	})
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *Es6RewriteRestParametersTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: Es6RewriteRestParametersTest.customSetUp()
func (c *Es6RewriteRestParametersTest) CustomSetUp() {
	c.EnableNormalize()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2016)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: Es6RewriteRestParametersTest.testUnusedRestParameterAtPositionZero()
func (c *Es6RewriteRestParametersTest) TestUnusedRestParameterAtPositionZero() {
	c.Test("function f(...zero) {}", "function f() {}")
}

// Java: Es6RewriteRestParametersTest.testUnusedRestParameterAtPositionOne()
func (c *Es6RewriteRestParametersTest) TestUnusedRestParameterAtPositionOne() {
	c.Test("function f(zero, ...one) {}", "function f(zero) {}")
}

// Java: Es6RewriteRestParametersTest.testUnusedRestParameterAtPositionTwo()
func (c *Es6RewriteRestParametersTest) TestUnusedRestParameterAtPositionTwo() {
	c.Test("function f(zero, one, ...two) {}", "function f(zero, one) {}")
}

// Java: Es6RewriteRestParametersTest.testUsedRestParameterAtPositionZero()
func (c *Es6RewriteRestParametersTest) TestUsedRestParameterAtPositionZero() {
	c.Test(
		"function f(...zero) { return zero; }",
		`function f() {
  let zero = $jscomp.getRestArguments.apply(0, arguments)
  return zero;
}
`)
}

// Java: Es6RewriteRestParametersTest.testUsedRestParameterAtPositionTwo()
func (c *Es6RewriteRestParametersTest) TestUsedRestParameterAtPositionTwo() {
	c.Test(
		"function f(zero, one, ...two) { return two; }",
		`function f(zero, one) {
  let two = $jscomp.getRestArguments.apply(2, arguments);
  return two;
}
`)
}

// Java: Es6RewriteRestParametersTest.testUsedRestParameterAtPositionTwo_maintainsNormalization()
func (c *Es6RewriteRestParametersTest) TestUsedRestParameterAtPositionTwo_maintainsNormalization() {
	c.Test(
		"function f(zero, one, ...two) { function inner() {} return two; }",
		`function f(zero, one) {
  function inner() {} // stays hoisted
  let two = $jscomp.getRestArguments.apply(2, arguments);
  return two;
}
`)
}

// Java: Es6RewriteRestParametersTest.testUsedRestParameterAtPositionTwo_maintainsNormalization_withoutReturn()
func (c *Es6RewriteRestParametersTest) TestUsedRestParameterAtPositionTwo_maintainsNormalization_withoutReturn() {
	c.Test(
		"function f(zero, one, ...two) { function inner() {} two; }",
		`function f(zero, one) {
  function inner() {} // stays hoisted
  let two = $jscomp.getRestArguments.apply(2, arguments);
  two;
}
`)
}

// Java: Es6RewriteRestParametersTest.testUnusedRestParameterAtPositionTwo_noGoodInsertionPoint()
func (c *Es6RewriteRestParametersTest) TestUnusedRestParameterAtPositionTwo_noGoodInsertionPoint() {
	c.Test(
		"function f(zero, one, ...two) { function inner() {} }",
		`function f(zero, one) {
  function inner() {} // stays hoisted
  let two = $jscomp.getRestArguments.apply(2, arguments); // declaration inserted
}
`)
}

// Java: Es6RewriteRestParametersTest.testUnusedRestParameterAtPositionZeroWithTypingOnFunction()
func (c *Es6RewriteRestParametersTest) TestUnusedRestParameterAtPositionZeroWithTypingOnFunction() {
	c.Test("/** @param {...number} zero */ function f(...zero) {}", "function f() {}")
}

// Java: Es6RewriteRestParametersTest.testUnusedRestParameterAtPositionZeroWithInlineTyping()
func (c *Es6RewriteRestParametersTest) TestUnusedRestParameterAtPositionZeroWithInlineTyping() {
	c.Test("function f(/** ...number */ ...zero) {}", "function f() {}")
}

// Java: Es6RewriteRestParametersTest.testUsedRestParameterAtPositionTwoWithTypingOnFunction()
func (c *Es6RewriteRestParametersTest) TestUsedRestParameterAtPositionTwoWithTypingOnFunction() {
	c.Test(
		"/** @param {...number} two */ function f(zero, one, ...two) { return two; }",
		`function f(zero, one) {
 let two = $jscomp.getRestArguments.apply(2, arguments);
 return two;
}
`)
}

// Java: Es6RewriteRestParametersTest.testUsedRestParameterAtPositionTwoWithTypingOnFunctionVariable()
func (c *Es6RewriteRestParametersTest) TestUsedRestParameterAtPositionTwoWithTypingOnFunctionVariable() {
	c.Test(
		"/** @param {...number} two */ var f = function(zero, one, ...two) { return two; }",
		`var f = function(zero, one) {
  let two = $jscomp.getRestArguments.apply(2, arguments);
  return two;
}
`)
}

// Java: Es6RewriteRestParametersTest.testUsedRestParameterAtPositionTwoWithTypingOnFunctionProperty()
func (c *Es6RewriteRestParametersTest) TestUsedRestParameterAtPositionTwoWithTypingOnFunctionProperty() {
	c.Test(
		"/** @param {...number} two */ ns.f = function(zero, one, ...two) { return two; }",
		`ns.f = function(zero, one) {
  let two = $jscomp.getRestArguments.apply(2, arguments);
  return two;
}
`)
}

// Java: Es6RewriteRestParametersTest.testUnusedRestParameterAtPositionTwoWithUsedParameterAtPositionOne()
func (c *Es6RewriteRestParametersTest) TestUnusedRestParameterAtPositionTwoWithUsedParameterAtPositionOne() {
	c.Test(
		"function f(zero, one, ...two) {one = (one === undefined) ? 1 : one;}",
		`function f(zero, one) {
  let two = $jscomp.getRestArguments.apply(2, arguments);
  one = (one === undefined) ? 1 : one;
}
`)
}

// Java: Es6RewriteRestParametersTest.testRestTranspiled_spreadExpressionsPreserved()
func (c *Es6RewriteRestParametersTest) TestRestTranspiled_spreadExpressionsPreserved() {
	c.Test(
		"function f(a, ...rest) { return f(...rest); }",
		`function f(a) {
  let rest = $jscomp.getRestArguments.apply(1, arguments);
  return f(...rest);
}
`)
}

// Java: Es6RewriteRestParametersTest.testArrayDestructuringRestDoesNotCrash()
func (c *Es6RewriteRestParametersTest) TestArrayDestructuringRestDoesNotCrash() {
	c.Test(
		"const keyFn = ([that, ...args]) => serializer(uidF, args);",
		`const keyFn = ([that, ...args]) => {
  return serializer(uidF, args);
};
`)
}

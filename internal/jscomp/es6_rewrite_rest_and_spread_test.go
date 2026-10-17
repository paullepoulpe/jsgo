// Port of com/google/javascript/jscomp/Es6RewriteRestAndSpreadTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/Es6RewriteRestAndSpreadTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteRestAndSpreadTest has the test cases for the transpilation of
// rest parameters and spread expressions.
type Es6RewriteRestAndSpreadTest struct {
	CompilerTestCaseBase
}

// Java: Es6RewriteRestAndSpreadTest.EXTERNS_BASE
var Es6RewriteRestAndSpreadTest_EXTERNS_BASE = NewTestExternsBuilder().
	AddFunction().
	AddJSCompLibraries().
	AddExtra("$jscomp.getRestArguments = function(argument) {};").
	Build()

// Java: Es6RewriteRestAndSpreadTest()
func NewEs6RewriteRestAndSpreadTest() *Es6RewriteRestAndSpreadTest {
	c := &Es6RewriteRestAndSpreadTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, Es6RewriteRestAndSpreadTest_EXTERNS_BASE)
	return c
}

func TestEs6RewriteRestAndSpreadTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteRestAndSpreadTest() })
}

// Java: Es6RewriteRestAndSpreadTest.getProcessor(Compiler)
func (c *Es6RewriteRestAndSpreadTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		NewInjectTranspilationRuntimeLibraries(compiler).Process(externs, root)
		NewEs6RewriteRestAndSpread(compiler).Process(externs, root)
	})
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *Es6RewriteRestAndSpreadTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: Es6RewriteRestAndSpreadTest.customSetUp()
func (c *Es6RewriteRestAndSpreadTest) CustomSetUp() {
	c.EnableNormalize()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2016)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadArrayLiteralIntoArrayLiteral()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadArrayLiteralIntoArrayLiteral() {
	c.Test("[...[1, 2], 3, ...[4], 5, 6, ...[], 7, 8]", "[1, 2, 3, 4, 5, 6, 7, 8]")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoArrayLiteral()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoArrayLiteral() {
	c.Test(
		"var arr = [1, 2, ...mid, 4, 5];",
		"var arr = [1, 2].concat((0, $jscomp.arrayFromIterable)(mid), [4, 5]);")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadFunctionReturnIntoArrayLiteral()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadFunctionReturnIntoArrayLiteral() {
	c.Test(
		"var arr = [1, 2, ...mid(), 4, 5];",
		"var arr = [1, 2].concat((0, $jscomp.arrayFromIterable)(mid()), [4, 5]);")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadFunctionArgumentsIntoArrayLiteral()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadFunctionArgumentsIntoArrayLiteral() {
	c.Test(
		"function f() { return [...arguments, 2]; };",
		`function f() {
  return [].concat((0, $jscomp.arrayFromIterable)(arguments), [2]);
};
`)
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableAndFunctionReturnIntoArrayLiteral()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableAndFunctionReturnIntoArrayLiteral() {
	c.Test(
		"var arr = [1, 2, ...mid, ...mid2(), 4, 5];",
		`var arr = [1,2].concat(
    (0, $jscomp.arrayFromIterable)(mid), (0, $jscomp.arrayFromIterable)(mid2()), [4, 5]);
`)
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadFunctionReturnIntoEntireArrayLiteral()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadFunctionReturnIntoEntireArrayLiteral() {
	c.Test("var arr = [...mid()];", "var arr = [].concat((0, $jscomp.arrayFromIterable)(mid()));")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadFunctionArgumentsIntoEntireArrayLiteral()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadFunctionArgumentsIntoEntireArrayLiteral() {
	c.Test(
		"function f() { return [...arguments]; };",
		`function f() {
  return [].concat((0, $jscomp.arrayFromIterable)(arguments));
};
`)
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadArrayLiteralIntoArrayLiteralWithinParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadArrayLiteralIntoArrayLiteralWithinParameterList() {
	c.Test("f(1, [2, ...[3], 4], 5);", "f(1, [2, 3, 4], 5);")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoArrayLiteralWithinParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoArrayLiteralWithinParameterList() {
	c.Test(
		"f(1, [2, ...mid, 4], 5);",
		"f(1, [2].concat((0, $jscomp.arrayFromIterable)(mid), [4]), 5);")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadFunctionReturnIntoArrayLiteralWithinParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadFunctionReturnIntoArrayLiteralWithinParameterList() {
	c.Test(
		"f(1, [2, ...mid(), 4], 5);",
		"f(1, [2].concat((0, $jscomp.arrayFromIterable)(mid()), [4]), 5);")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadArrayLiteralIntoEntireParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadArrayLiteralIntoEntireParameterList() {
	c.Test("f(...[0, 1, 2]);", "f.apply(null, [0, 1, 2]);")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadArrayLiteralIntoParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadArrayLiteralIntoParameterList() {
	c.Test("f(...[0, 1, 2], 3);", "f.apply(null, [0, 1, 2, 3]);")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoEntireParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoEntireParameterList() {
	c.Test("f(...arr);", "f.apply(null, (0, $jscomp.arrayFromIterable)(arr));")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoParameterList() {
	c.Test(
		"f(0, ...arr, 2);", "f.apply(null, [0].concat((0, $jscomp.arrayFromIterable)(arr), [2]));")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadFunctionReturnIntoEntireParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadFunctionReturnIntoEntireParameterList() {
	c.Test("f(...g());", "f.apply(null, (0, $jscomp.arrayFromIterable)(g()));")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadFunctionReturnIntoParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadFunctionReturnIntoParameterList() {
	c.Test(
		"f(0, ...g(), 2);", "f.apply(null, [0].concat((0, $jscomp.arrayFromIterable)(g()), [2]));")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoIifeParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoIifeParameterList() {
	c.Test(
		"(function() {})(...arr);",
		"(function() {}).apply(null, (0, $jscomp.arrayFromIterable)(arr))")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoAnonymousFunctionParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoAnonymousFunctionParameterList() {
	c.Test("getF()(...args);", "getF().apply(null, (0, $jscomp.arrayFromIterable)(args));")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoMethodParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoMethodParameterList() {
	c.Test(
		c.Externs(
			Es6RewriteRestAndSpreadTest_EXTERNS_BASE+
				`/**
 * @constructor
 * Skipping at-struct here to allow for string access.
 */
function TestClass() { }

/** @param {...string} args */
TestClass.prototype.testMethod = function(args) { }

/** @return {!TestClass} */
function testClassFactory() { }
`),
		c.Srcs(
			`var obj = new TestClass();
obj.testMethod(...arr);
obj['testMethod'](...arr);
`),
		c.Expected(
			`var obj = new TestClass();
obj.testMethod.apply(obj, (0, $jscomp.arrayFromIterable)(arr));
obj["testMethod"].apply(obj, (0, $jscomp.arrayFromIterable)(arr));
`))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoMethodParameterListInCast()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoMethodParameterListInCast() {
	c.Test(
		c.Externs(
			Es6RewriteRestAndSpreadTest_EXTERNS_BASE+
				`/**
 * @constructor
 * Skipping at-struct here to allow for string access.
 */
function TestClass() { }

/** @param {...string} args */
TestClass.prototype.testMethod = function(args) { }

/** @return {!TestClass} */
function testClassFactory() { }
`),
		c.Srcs(
			`var obj = new TestClass();
(/** @type {?} */ (obj.testMethod))(...arr);
(/** @type {?} */ (/** @type {?} */ (obj.testMethod)))(...arr);
(/** @type {?} */ (obj['testMethod']))(...arr);
`),
		c.Expected(
			`var obj = new TestClass();
obj.testMethod.apply(obj, (0, $jscomp.arrayFromIterable)(arr));
obj.testMethod.apply(obj, (0, $jscomp.arrayFromIterable)(arr));
obj['testMethod'].apply(obj, (0, $jscomp.arrayFromIterable)(arr));
`))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoDeepMethodParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoDeepMethodParameterList() {
	c.Test(
		c.Externs(
			Es6RewriteRestAndSpreadTest_EXTERNS_BASE+
				`/** @param {...number} args */ function numberVarargFn(args) { }

/** @type {!Iterable<number>} */ var numberIterable;
`),
		c.Srcs(
			`var x = {y: {z: {m: numberVarargFn}}};
x.y.z.m(...numberIterable);
`),
		c.Expected(
			`var x = {y: {z: {m: numberVarargFn}}};
x.y.z.m.apply(x.y.z, (0, $jscomp.arrayFromIterable)(numberIterable));
`))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoMethodParameterList_freeCall()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoMethodParameterList_freeCall() {
	c.Test(
		c.Externs(
			Es6RewriteRestAndSpreadTest_EXTERNS_BASE+
				"/**\n * @constructor\n * Skipping at-struct here to allow for string access.\n */\nfunction TestClass() { }\n\n// Adding @this {?} to allow calling testMethod without passing a TestClass as\n// `this`\n/** @param {...string} args @this {?} */\nTestClass.prototype.testMethod = function(args) { }\n\n/** @return {!TestClass} */\nfunction testClassFactory() { }\n"),
		c.Srcs(
			`var obj = new TestClass();
// The (0, obj.testMethod) tells the compiler that this is a 'free call'.
(0, obj.testMethod)(...arr);
(0, obj['testMethod'])(...arr);
`),
		c.Expected(
			`var obj = new TestClass();
obj.testMethod.apply(null, (0, $jscomp.arrayFromIterable)(arr));
obj["testMethod"].apply(null, (0, $jscomp.arrayFromIterable)(arr));
`))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadMultipleArrayLiteralsIntoParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadMultipleArrayLiteralsIntoParameterList() {
	c.Test(
		c.Externs(
			CompilerTestCase_MINIMAL_EXTERNS+"/** @param {...number} args */ function numberVarargFn(args) { }"),
		c.Srcs("numberVarargFn(...[1, 2], 3, ...[4, 5], 6, ...[7, 8])"),
		c.Expected("numberVarargFn.apply(null, [1, 2, 3, 4, 5, 6, 7, 8])"))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadMultipleVariablesIntoParameterList()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadMultipleVariablesIntoParameterList() {
	c.Test(
		c.Externs(
			Es6RewriteRestAndSpreadTest_EXTERNS_BASE+
				`/** @param {...number} args */ function numberVarargFn(args) { }
/** @type {!Iterable<number>} */ var numberIterable;
`),
		c.Srcs("numberVarargFn(0, ...numberIterable, 2, ...numberIterable, 4);"),
		c.Expected(
			`numberVarargFn.apply(
    null,
    [0].concat(
        (0, $jscomp.arrayFromIterable)(numberIterable),
        [2],
        (0, $jscomp.arrayFromIterable)(numberIterable),
        [4]));
`))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoMethodParameterListOnAnonymousRecieverWithSideEffects()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoMethodParameterListOnAnonymousRecieverWithSideEffects() {
	c.Test(
		c.Externs(
			Es6RewriteRestAndSpreadTest_EXTERNS_BASE+
				`/**
 * @constructor
 * @struct
 */
function TestClass() { }

/** @param {...string} args */
TestClass.prototype.testMethod = function(args) { }

/** @return {!TestClass} */
function testClassFactory() { }

/** @type {!Iterable<string>} */ var stringIterable;
`),
		c.Srcs("testClassFactory().testMethod(...stringIterable);"),
		c.Expected(
			`var $jscomp$spread$args0;
($jscomp$spread$args0 = testClassFactory()).testMethod.apply(
    $jscomp$spread$args0, (0, $jscomp.arrayFromIterable)(stringIterable));
`))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoMethodParameterListOnReceiverWithSideEffects_freeCall()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoMethodParameterListOnReceiverWithSideEffects_freeCall() {
	c.Test(
		c.Externs(
			Es6RewriteRestAndSpreadTest_EXTERNS_BASE+
				"/**\n * @constructor\n// Skip @struct to allow for bracket access\n */\nfunction TestClass() { }\n\n// Add @this {?} to allow calling testMethod without passing a TestClass as\n// `this`\n/** @param {...string} args @this {null} */\nTestClass.prototype.testMethod = function(args) { }\n\n/** @return {!TestClass} */\nfunction testClassFactory() { }\n\n/** @type {!Iterable<string>} */ var stringIterable;\n"),
		c.Srcs(
			`(0, testClassFactory().testMethod)(...stringIterable);
(0, testClassFactory()['testMethod'])(...stringIterable);
`),
		c.Expected(
			`testClassFactory().testMethod.apply(
    null, (0, $jscomp.arrayFromIterable)(stringIterable));
testClassFactory()["testMethod"].apply(
    null, (0, $jscomp.arrayFromIterable)(stringIterable));
`))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoMethodParameterListOnConditionalRecieverWithSideEffects()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoMethodParameterListOnConditionalRecieverWithSideEffects() {
	c.Test(
		c.Externs(
			Es6RewriteRestAndSpreadTest_EXTERNS_BASE+
				`/**
 * @constructor
 * @struct
 */
function TestClass() { }

/** @param {...string} args */
TestClass.prototype.testMethod = function(args) { }

/** @return {!TestClass} */
function testClassFactory() { }

/** @type {!Iterable<string>} */ var stringIterable;
`),
		c.Srcs("var x = b ? testClassFactory().testMethod(...stringIterable) : null;"),
		c.Expected(
			`var $jscomp$spread$args0;
var x = b ? ($jscomp$spread$args0 = testClassFactory()).testMethod.apply(
    $jscomp$spread$args0, (0, $jscomp.arrayFromIterable)(stringIterable))
        : null;
`))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoMethodParameterListOnConditionalRecieverWithSideEffectsInCast()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoMethodParameterListOnConditionalRecieverWithSideEffectsInCast() {
	c.Test(
		c.Externs(
			Es6RewriteRestAndSpreadTest_EXTERNS_BASE+
				`/**
 * @constructor
 * @struct
 */
function TestClass() { }

/** @param {...string} args */
TestClass.prototype.testMethod = function(args) { }

/** @return {!TestClass} */
function testClassFactory() { }

/** @type {!Iterable<string>} */ var stringIterable;
`),
		c.Srcs(
			`var x = b ?
    /** @type {?} */ (testClassFactory().testMethod)(...stringIterable) :
    null;
`),
		c.Expected(
			`var $jscomp$spread$args0;
var x = b ? ($jscomp$spread$args0 = testClassFactory()).testMethod.apply(
    $jscomp$spread$args0, (0, $jscomp.arrayFromIterable)(stringIterable))
        : null;
`))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoMethodParameterListOnRecieversWithSideEffectsMultipleTimesInOneScope()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoMethodParameterListOnRecieversWithSideEffectsMultipleTimesInOneScope() {
	c.Test(
		c.Externs(
			Es6RewriteRestAndSpreadTest_EXTERNS_BASE+
				`/**
 * @constructor
 * @struct
 */
function TestClass() { }

/** @param {...string} args */
TestClass.prototype.testMethod = function(args) { }

/** @return {!TestClass} */
function testClassFactory() { }

/** @type {!Iterable<string>} */ var stringIterable;
`),
		c.Srcs(
			`testClassFactory().testMethod(...stringIterable);
testClassFactory().testMethod(...stringIterable);
`),
		c.Expected(
			`var $jscomp$spread$args0;
($jscomp$spread$args0 = testClassFactory()).testMethod.apply(
    $jscomp$spread$args0, (0, $jscomp.arrayFromIterable)(stringIterable));
var $jscomp$spread$args1;
($jscomp$spread$args1 = testClassFactory()).testMethod.apply(
    $jscomp$spread$args1, (0, $jscomp.arrayFromIterable)(stringIterable));
`))
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadIntoSuperThrows()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadIntoSuperThrows() {
	// All ES6 classes must be transpiled away before this pass runs,
	// because we cannot spread into super() calls without creating invalid syntax.
	AssertThrows[jutil.IRuntimeException](
		func() {
			c.TestSame(
				`class A {
  constructor(...args) {
      this.p = args;
  }
}

class B extends A {
   constructor(...args) {
     super(0, ...args, 2);
   }
}
`)
		})
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoParameterListWithinArrayLiteral()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoParameterListWithinArrayLiteral() {
	c.Test(
		"[1, f(2, ...mid, 4), 5];",
		"[1, f.apply(null, [2].concat((0, $jscomp.arrayFromIterable)(mid), [4])), 5];")
}

// Java: Es6RewriteRestAndSpreadTest.testSpreadVariableIntoNew()
func (c *Es6RewriteRestAndSpreadTest) TestSpreadVariableIntoNew() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		c.Srcs("new F(...args);"),
		c.Expected(
			"new (Function.prototype.bind.apply(F, [null].concat((0, $jscomp.arrayFromIterable)(args))));"))
}

// Java: Es6RewriteRestAndSpreadTest.testUnusedRestParameterAtPositionZero()
func (c *Es6RewriteRestAndSpreadTest) TestUnusedRestParameterAtPositionZero() {
	c.Test("function f(...zero) {}", "function f() {}")
}

// Java: Es6RewriteRestAndSpreadTest.testUnusedRestParameterAtPositionOne()
func (c *Es6RewriteRestAndSpreadTest) TestUnusedRestParameterAtPositionOne() {
	c.Test("function f(zero, ...one) {}", "function f(zero) {}")
}

// Java: Es6RewriteRestAndSpreadTest.testUnusedRestParameterAtPositionTwo()
func (c *Es6RewriteRestAndSpreadTest) TestUnusedRestParameterAtPositionTwo() {
	c.Test("function f(zero, one, ...two) {}", "function f(zero, one) {}")
}

// Java: Es6RewriteRestAndSpreadTest.testUsedRestParameterAtPositionZero()
func (c *Es6RewriteRestAndSpreadTest) TestUsedRestParameterAtPositionZero() {
	c.Test(
		"function f(...zero) { return zero; }",
		`function f() {
  let zero = $jscomp.getRestArguments.apply(0, arguments)
  return zero;
}
`)
}

// Java: Es6RewriteRestAndSpreadTest.testUsedRestParameterAtPositionTwo()
func (c *Es6RewriteRestAndSpreadTest) TestUsedRestParameterAtPositionTwo() {
	c.Test(
		"function f(zero, one, ...two) { return two; }",
		`function f(zero, one) {
  let two = $jscomp.getRestArguments.apply(2, arguments);
  return two;
}
`)
}

// Java: Es6RewriteRestAndSpreadTest.testUsedRestParameterAtPositionTwo_maintainsNormalization()
func (c *Es6RewriteRestAndSpreadTest) TestUsedRestParameterAtPositionTwo_maintainsNormalization() {
	c.Test(
		"function f(zero, one, ...two) { function inner() {} return two; }",
		`function f(zero, one) {
  function inner() {} // stays hoisted
  let two = $jscomp.getRestArguments.apply(2, arguments);
  return two;
}
`)
}

// Java: Es6RewriteRestAndSpreadTest.testUsedRestParameterAtPositionTwo_maintainsNormalization_withoutReturn()
func (c *Es6RewriteRestAndSpreadTest) TestUsedRestParameterAtPositionTwo_maintainsNormalization_withoutReturn() {
	c.Test(
		"function f(zero, one, ...two) { function inner() {} two; }",
		`function f(zero, one) {
  function inner() {} // stays hoisted
  let two = $jscomp.getRestArguments.apply(2, arguments);
  two;
}
`)
}

// Java: Es6RewriteRestAndSpreadTest.testUnusedRestParameterAtPositionTwo_noGoodInsertionPoint()
func (c *Es6RewriteRestAndSpreadTest) TestUnusedRestParameterAtPositionTwo_noGoodInsertionPoint() {
	c.Test(
		"function f(zero, one, ...two) { function inner() {} }",
		`function f(zero, one) {
  function inner() {} // stays hoisted
  let two = $jscomp.getRestArguments.apply(2, arguments); // declaration inserted
}
`)
}

// Java: Es6RewriteRestAndSpreadTest.testUnusedRestParameterAtPositionZeroWithTypingOnFunction()
func (c *Es6RewriteRestAndSpreadTest) TestUnusedRestParameterAtPositionZeroWithTypingOnFunction() {
	c.Test("/** @param {...number} zero */ function f(...zero) {}", "function f() {}")
}

// Java: Es6RewriteRestAndSpreadTest.testUnusedRestParameterAtPositionZeroWithInlineTyping()
func (c *Es6RewriteRestAndSpreadTest) TestUnusedRestParameterAtPositionZeroWithInlineTyping() {
	c.Test("function f(/** ...number */ ...zero) {}", "function f() {}")
}

// Java: Es6RewriteRestAndSpreadTest.testUsedRestParameterAtPositionTwoWithTypingOnFunction()
func (c *Es6RewriteRestAndSpreadTest) TestUsedRestParameterAtPositionTwoWithTypingOnFunction() {
	c.Test(
		"/** @param {...number} two */ function f(zero, one, ...two) { return two; }",
		`function f(zero, one) {
 let two = $jscomp.getRestArguments.apply(2, arguments);
 return two;
}
`)
}

// Java: Es6RewriteRestAndSpreadTest.testUsedRestParameterAtPositionTwoWithTypingOnFunctionVariable()
func (c *Es6RewriteRestAndSpreadTest) TestUsedRestParameterAtPositionTwoWithTypingOnFunctionVariable() {
	c.Test(
		"/** @param {...number} two */ var f = function(zero, one, ...two) { return two; }",
		`var f = function(zero, one) {
  let two = $jscomp.getRestArguments.apply(2, arguments);
  return two;
}
`)
}

// Java: Es6RewriteRestAndSpreadTest.testUsedRestParameterAtPositionTwoWithTypingOnFunctionProperty()
func (c *Es6RewriteRestAndSpreadTest) TestUsedRestParameterAtPositionTwoWithTypingOnFunctionProperty() {
	c.Test(
		"/** @param {...number} two */ ns.f = function(zero, one, ...two) { return two; }",
		`ns.f = function(zero, one) {
  let two = $jscomp.getRestArguments.apply(2, arguments);
  return two;
}
`)
}

// Java: Es6RewriteRestAndSpreadTest.testUnusedRestParameterAtPositionTwoWithUsedParameterAtPositionOne()
func (c *Es6RewriteRestAndSpreadTest) TestUnusedRestParameterAtPositionTwoWithUsedParameterAtPositionOne() {
	c.Test(
		"function f(zero, one, ...two) {one = (one === undefined) ? 1 : one;}",
		`function f(zero, one) {
  let two = $jscomp.getRestArguments.apply(2, arguments);
  one = (one === undefined) ? 1 : one;
}
`)
}

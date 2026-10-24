// Port of com/google/javascript/jscomp/Es6RewriteSpreadExpressionsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/Es6RewriteSpreadExpressionsTest.java)

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

// Es6RewriteSpreadExpressionsTest has the test cases for the transpilation of
// spread expressions.
type Es6RewriteSpreadExpressionsTest struct {
	CompilerTestCaseBase
}

// Java: Es6RewriteSpreadExpressionsTest.EXTERNS_BASE
var Es6RewriteSpreadExpressionsTest_EXTERNS_BASE = NewTestExternsBuilder().
	AddFunction().
	AddJSCompLibraries().
	Build()

// Java: Es6RewriteSpreadExpressionsTest()
func NewEs6RewriteSpreadExpressionsTest() *Es6RewriteSpreadExpressionsTest {
	c := &Es6RewriteSpreadExpressionsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, Es6RewriteSpreadExpressionsTest_EXTERNS_BASE)
	return c
}

func TestEs6RewriteSpreadExpressionsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteSpreadExpressionsTest() })
}

// Java: Es6RewriteSpreadExpressionsTest.getProcessor(Compiler)
func (c *Es6RewriteSpreadExpressionsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		NewInjectTranspilationRuntimeLibraries(compiler).Process(externs, root)
		NewEs6RewriteSpreadExpressions(compiler).Process(externs, root)
	})
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *Es6RewriteSpreadExpressionsTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: Es6RewriteSpreadExpressionsTest.customSetUp()
func (c *Es6RewriteSpreadExpressionsTest) CustomSetUp() {
	c.EnableNormalize()
	c.SetGenericNameReplacements(transpileTests_ImmutableMapOf("SPREAD_ARGS", "$jscomp$spread$args"))
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2016)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadArrayLiteralIntoArrayLiteral()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadArrayLiteralIntoArrayLiteral() {
	c.Test("[...[1, 2], 3, ...[4], 5, 6, ...[], 7, 8]", "[1, 2, 3, 4, 5, 6, 7, 8]")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoArrayLiteral()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoArrayLiteral() {
	c.Test(
		"var arr = [1, 2, ...mid, 4, 5];",
		"var arr = [1, 2].concat((0, $jscomp.arrayFromIterable)(mid), [4, 5]);")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadFunctionReturnIntoArrayLiteral()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadFunctionReturnIntoArrayLiteral() {
	c.Test(
		"var arr = [1, 2, ...mid(), 4, 5];",
		"var arr = [1, 2].concat((0, $jscomp.arrayFromIterable)(mid()), [4, 5]);")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadFunctionArgumentsIntoArrayLiteral()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadFunctionArgumentsIntoArrayLiteral() {
	c.Test(
		"function f() { return [...arguments, 2]; };",
		`function f() {
  return [].concat((0, $jscomp.arrayFromIterable)(arguments), [2]);
};
`)
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableAndFunctionReturnIntoArrayLiteral()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableAndFunctionReturnIntoArrayLiteral() {
	c.Test(
		"var arr = [1, 2, ...mid, ...mid2(), 4, 5];",
		`var arr = [1,2].concat(
    (0, $jscomp.arrayFromIterable)(mid), (0, $jscomp.arrayFromIterable)(mid2()), [4, 5]);
`)
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadFunctionReturnIntoEntireArrayLiteral()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadFunctionReturnIntoEntireArrayLiteral() {
	c.Test("var arr = [...mid()];", "var arr = [].concat((0, $jscomp.arrayFromIterable)(mid()));")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadFunctionArgumentsIntoEntireArrayLiteral()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadFunctionArgumentsIntoEntireArrayLiteral() {
	c.Test(
		"function f() { return [...arguments]; };",
		`function f() {
  return [].concat((0, $jscomp.arrayFromIterable)(arguments));
};
`)
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadArrayLiteralIntoArrayLiteralWithinParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadArrayLiteralIntoArrayLiteralWithinParameterList() {
	c.Test("f(1, [2, ...[3], 4], 5);", "f(1, [2, 3, 4], 5);")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoArrayLiteralWithinParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoArrayLiteralWithinParameterList() {
	c.Test(
		"f(1, [2, ...mid, 4], 5);",
		"f(1, [2].concat((0, $jscomp.arrayFromIterable)(mid), [4]), 5);")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadFunctionReturnIntoArrayLiteralWithinParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadFunctionReturnIntoArrayLiteralWithinParameterList() {
	c.Test(
		"f(1, [2, ...mid(), 4], 5);",
		"f(1, [2].concat((0, $jscomp.arrayFromIterable)(mid()), [4]), 5);")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadArrayLiteralIntoEntireParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadArrayLiteralIntoEntireParameterList() {
	c.Test("f(...[0, 1, 2]);", "f.apply(null, [0, 1, 2]);")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadArrayLiteralIntoParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadArrayLiteralIntoParameterList() {
	c.Test("f(...[0, 1, 2], 3);", "f.apply(null, [0, 1, 2, 3]);")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoEntireParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoEntireParameterList() {
	c.Test("f(...arr);", "f.apply(null, (0, $jscomp.arrayFromIterable)(arr));")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoParameterList() {
	c.Test(
		"f(0, ...arr, 2);", "f.apply(null, [0].concat((0, $jscomp.arrayFromIterable)(arr), [2]));")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadFunctionReturnIntoEntireParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadFunctionReturnIntoEntireParameterList() {
	c.Test("f(...g());", "f.apply(null, (0, $jscomp.arrayFromIterable)(g()));")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadFunctionReturnIntoParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadFunctionReturnIntoParameterList() {
	c.Test(
		"f(0, ...g(), 2);", "f.apply(null, [0].concat((0, $jscomp.arrayFromIterable)(g()), [2]));")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoIifeParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoIifeParameterList() {
	c.Test(
		"(function() {})(...arr);",
		"(function() {}).apply(null, (0, $jscomp.arrayFromIterable)(arr))")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoAnonymousFunctionParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoAnonymousFunctionParameterList() {
	c.Test("getF()(...args);", "getF().apply(null, (0, $jscomp.arrayFromIterable)(args));")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoMethodParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoMethodParameterList() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
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

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoMethodParameterListInCast()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoMethodParameterListInCast() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
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

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoDeepMethodParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoDeepMethodParameterList() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
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

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoMethodParameterList_freeCall()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoMethodParameterList_freeCall() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
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

// Java: Es6RewriteSpreadExpressionsTest.testSpreadMultipleArrayLiteralsIntoParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadMultipleArrayLiteralsIntoParameterList() {
	c.Test(
		c.Externs(
			CompilerTestCase_MINIMAL_EXTERNS+"/** @param {...number} args */ function numberVarargFn(args) { }"),
		c.Srcs("numberVarargFn(...[1, 2], 3, ...[4, 5], 6, ...[7, 8])"),
		c.Expected("numberVarargFn.apply(null, [1, 2, 3, 4, 5, 6, 7, 8])"))
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadMultipleVariablesIntoParameterList()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadMultipleVariablesIntoParameterList() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
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

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoMethodParameterListOnAnonymousRecieverWithSideEffects()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoMethodParameterListOnAnonymousRecieverWithSideEffects() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
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
			`var SPREAD_ARGS$0;
(SPREAD_ARGS$0 = testClassFactory()).testMethod.apply(
    SPREAD_ARGS$0, (0, $jscomp.arrayFromIterable)(stringIterable));
`))
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoBracketAccessMethodParameterListOnAnonymousReceiverWithSideEffects()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoBracketAccessMethodParameterListOnAnonymousReceiverWithSideEffects() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
				`/**
 * @constructor
 */
function TestClass() { }

/** @param {...string} args */
TestClass.prototype.testMethod = function(args) { }

/** @return {!TestClass} */
function testClassFactory() { }

/** @type {!Iterable<string>} */ var stringIterable;
`),
		c.Srcs("testClassFactory()['testMethod'](...stringIterable);"),
		c.Expected(
			`var SPREAD_ARGS$0;
(SPREAD_ARGS$0 = testClassFactory())["testMethod"].apply(
    SPREAD_ARGS$0, (0, $jscomp.arrayFromIterable)(stringIterable));
`))
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoDynamicBracketAccessMethodParameterListOnReceiverWithSideEffects()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoDynamicBracketAccessMethodParameterListOnReceiverWithSideEffects() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
				`/** @type {!Array<!Object>} */ var queue;
/** @type {string} */ var cmd;
/** @type {!Iterable<string>} */ var stringIterable;
`),
		c.Srcs("queue.shift()[cmd](...stringIterable);"),
		c.Expected(
			`var SPREAD_ARGS$0;
(SPREAD_ARGS$0 = queue.shift())[cmd].apply(
    SPREAD_ARGS$0, (0, $jscomp.arrayFromIterable)(stringIterable));
`))
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoMethodParameterListOnReceiverWithSideEffects_freeCall()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoMethodParameterListOnReceiverWithSideEffects_freeCall() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
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

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoMethodParameterListOnConditionalRecieverWithSideEffects()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoMethodParameterListOnConditionalRecieverWithSideEffects() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
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
			`var SPREAD_ARGS$0;
var x = b ? (SPREAD_ARGS$0 = testClassFactory()).testMethod.apply(
    SPREAD_ARGS$0, (0, $jscomp.arrayFromIterable)(stringIterable))
        : null;
`))
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoMethodParameterListOnConditionalRecieverWithSideEffectsInCast()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoMethodParameterListOnConditionalRecieverWithSideEffectsInCast() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
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
			`var SPREAD_ARGS$0;
var x = b ? (SPREAD_ARGS$0 = testClassFactory()).testMethod.apply(
    SPREAD_ARGS$0, (0, $jscomp.arrayFromIterable)(stringIterable))
        : null;
`))
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoMethodParameterListOnRecieversWithSideEffectsMultipleTimesInOneScope()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoMethodParameterListOnRecieversWithSideEffectsMultipleTimesInOneScope() {
	c.Test(
		c.Externs(
			Es6RewriteSpreadExpressionsTest_EXTERNS_BASE+
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
			`var SPREAD_ARGS$0;
(SPREAD_ARGS$0 = testClassFactory()).testMethod.apply(
    SPREAD_ARGS$0, (0, $jscomp.arrayFromIterable)(stringIterable));
var SPREAD_ARGS$1;
(SPREAD_ARGS$1 = testClassFactory()).testMethod.apply(
    SPREAD_ARGS$1, (0, $jscomp.arrayFromIterable)(stringIterable));
`))
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadIntoSuperThrows()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadIntoSuperThrows() {
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

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoParameterListWithinArrayLiteral()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoParameterListWithinArrayLiteral() {
	c.Test(
		"[1, f(2, ...mid, 4), 5];",
		"[1, f.apply(null, [2].concat((0, $jscomp.arrayFromIterable)(mid), [4])), 5];")
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadVariableIntoNew()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadVariableIntoNew() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		c.Srcs("new F(...args);"),
		c.Expected(
			"new (Function.prototype.bind.apply(F, [null].concat((0, $jscomp.arrayFromIterable)(args))));"))
}

// Java: Es6RewriteSpreadExpressionsTest.testSpreadTranspiled_restParametersPreserved()
func (c *Es6RewriteSpreadExpressionsTest) TestSpreadTranspiled_restParametersPreserved() {
	c.Test(
		"function f(a, ...rest) { return f(...rest); }",
		"function f(a, ...rest) { return f.apply(null, (0, $jscomp.arrayFromIterable)(rest)); }")
}

// Java: Es6RewriteSpreadExpressionsTest.testArrayDestructuringRestPreserved()
func (c *Es6RewriteSpreadExpressionsTest) TestArrayDestructuringRestPreserved() {
	c.Test(
		"const keyFn = ([that, ...args]) => serializer(uidF, args);",
		`const keyFn = ([that, ...args]) => {
  return serializer(uidF, args);
};
`)
}

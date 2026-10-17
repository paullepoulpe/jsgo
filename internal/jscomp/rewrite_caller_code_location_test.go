// Port of com/google/javascript/jscomp/RewriteCallerCodeLocationTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/RewriteCallerCodeLocationTest.java)

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

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteCallerCodeLocationTest tests RewriteCallerCodeLocation.
type RewriteCallerCodeLocationTest struct {
	CompilerTestCaseBase
}

// Java: RewriteCallerCodeLocationTest()
func NewRewriteCallerCodeLocationTest() *RewriteCallerCodeLocationTest {
	c := &RewriteCallerCodeLocationTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteCallerCodeLocationTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteCallerCodeLocationTest() })
}

// Java: RewriteCallerCodeLocationTest.getProcessor(Compiler)
func (c *RewriteCallerCodeLocationTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		NewRewriteCallerCodeLocation(compiler).Process(externs, root)
	})
}

// Java: RewriteCallerCodeLocationTest.testSimpleCallee()
func (c *RewriteCallerCodeLocationTest) TestSimpleCallee() {
	c.Test(
		`function signal(here = goog.callerLocation()) {}
signal();
`,
		`function signal(here = goog.callerLocation()) {}
signal(goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:0'))
`)
}

// Java: RewriteCallerCodeLocationTest.testSimpleCallee2()
func (c *RewriteCallerCodeLocationTest) TestSimpleCallee2() {
	c.Test(
		`function signal(val, here = goog.callerLocation()) {}
const mySignal = (0, signal)(0);
`,
		`function signal(val, here = goog.callerLocation()) {}
const mySignal = signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:17'));
`)
}

// Java: RewriteCallerCodeLocationTest.testAlisedCalleeDoesNotRewrite()
func (c *RewriteCallerCodeLocationTest) TestAlisedCalleeDoesNotRewrite() {
	c.Test(
		`function signal(here = goog.callerLocation()) {}
const mySignal = signal;
mySignal();
`,
		`function signal(here = goog.callerLocation()) {}
const mySignal = signal;
mySignal();
`)
}

// Java: RewriteCallerCodeLocationTest.testRenamedExportDoesNotRewrite()
func (c *RewriteCallerCodeLocationTest) TestRenamedExportDoesNotRewrite() {
	c.TestSame(
		`function localSignal(val, here = goog.callerLocation()) {
  return val;
}
exports.signal = localSignal;
signal(0);
`)
}

// Java: RewriteCallerCodeLocationTest.testUndefined()
func (c *RewriteCallerCodeLocationTest) TestUndefined() {
	c.EnableTypeCheck()
	c.EnableParseTypeInfo()
	c.TestError(
		`function signal(val, here = goog.callerLocation()) {}
const mySignal = (0, signal)(0, undefined);
`,
		RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR)
}

// Java: RewriteCallerCodeLocationTest.testAssigned()
func (c *RewriteCallerCodeLocationTest) TestAssigned() {
	c.Test(
		`function signal(val, here = goog.callerLocation()) {}
const mySignal = signal(0);
`,
		`function signal(val, here = goog.callerLocation()) {}
const mySignal = signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:17'));
`)
}

// Java: RewriteCallerCodeLocationTest.testAliasedFunctionDoesNotRewrite()
func (c *RewriteCallerCodeLocationTest) TestAliasedFunctionDoesNotRewrite() {
	// aliased functions are not rewritten
	c.Test(
		"function signal(val, here = goog.callerLocation()) {}\nconst mySignal = signal(0); // we will rewrite this `signal(0)` call\n// mySignal will not be rewritten because we will not rewrite aliased functions\nconst myComputed = computed(() => mySignal[0]() % 2 === 0);\n",
		`function signal(val, here = goog.callerLocation()) {}
const mySignal = signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:17'));
const myComputed = computed(() => mySignal[0]() % 2 === 0);
`)
}

// Java: RewriteCallerCodeLocationTest.testArrayDestructuring()
func (c *RewriteCallerCodeLocationTest) TestArrayDestructuring() {
	c.Test(
		`function signal(val, here = goog.callerLocation()) {}
const [foo, setFoo] = signal(0);
`,
		`function signal(val, here = goog.callerLocation()) {}
const [foo, setFoo] = signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:22'));
`)
}

// Java: RewriteCallerCodeLocationTest.testObjectDestructuring()
func (c *RewriteCallerCodeLocationTest) TestObjectDestructuring() {
	c.Test(
		`function signal(val, here = goog.callerLocation()) {}
const {foo, setFoo} = signal(0);
`,
		`function signal(val, here = goog.callerLocation()) {}
const {foo, setFoo} = signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:22'));
`)
}

// Java: RewriteCallerCodeLocationTest.testPropertyAssignment()
func (c *RewriteCallerCodeLocationTest) TestPropertyAssignment() {
	c.Test(
		`function signal(val, here = goog.callerLocation()) {}
const obj = {prop: signal(0)};
`,
		`function signal(val, here = goog.callerLocation()) {}
const obj = {prop: signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:19'))};
`)
}

// Java: RewriteCallerCodeLocationTest.testTernaryAssignment()
func (c *RewriteCallerCodeLocationTest) TestTernaryAssignment() {
	c.Test(
		`function signal(val, here = goog.callerLocation()) {}
const maybeSignal = Math.random() ? signal(0) : null;
`,
		`function signal(val, here = goog.callerLocation()) {}
const maybeSignal = Math.random() ? signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:36')) : null;
`)
}

// Java: RewriteCallerCodeLocationTest.testIntermediateFunction()
func (c *RewriteCallerCodeLocationTest) TestIntermediateFunction() {
	c.Test(
		`function signal(val, here = goog.callerLocation()) {}
const intermediateFunction = (() => signal(0))();
`,
		`function signal(val, here = goog.callerLocation()) {}
const intermediateFunction = (() => signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:36')))();
`)
}

// Java: RewriteCallerCodeLocationTest.testArray()
func (c *RewriteCallerCodeLocationTest) TestArray() {
	c.Test(
		`function signal(val, here = goog.callerLocation()) {}
const signalArray = [signal(0), signal(1)];
`,
		`function signal(val, here = goog.callerLocation()) {}
const signalArray = [signal(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:21')), signal(1, goog.callerLocationIdInternalDoNotCallOrElse('testcode:2:32'))];
`)
}

// Java: RewriteCallerCodeLocationTest.testNoCallerLocationDefaultParameter()
func (c *RewriteCallerCodeLocationTest) TestNoCallerLocationDefaultParameter() {
	// no rewriting if `= goog.callerLocation()` is not set as default parameter
	c.TestSame(
		`function signal() {}
signal();
`)
}

// Java: RewriteCallerCodeLocationTest.testGoogCallerLocationDefinitionInBaseJsDoesNotError()
func (c *RewriteCallerCodeLocationTest) TestGoogCallerLocationDefinitionInBaseJsDoesNotError() {
	baseJsSource :=
		`var goog = goog || {};
goog.callerLocation = function() {};
`
	baseJs := SourceFile_FromCode("javascript/closure/base.js", baseJsSource)
	// this code is only allowed in base.js
	c.Test(c.Srcs(baseJs), c.Expected(baseJs))
}

// Java: RewriteCallerCodeLocationTest.testNoRewritingIfCallerLocationIsSpecified()
func (c *RewriteCallerCodeLocationTest) TestNoRewritingIfCallerLocationIsSpecified() {
	c.EnableTypeCheck()
	c.EnableParseTypeInfo()

	baseJsSource :=
		`var goog = goog || {};
goog.callerLocation = function() {};
`
	baseJs := SourceFile_FromCode("javascript/closure/base.js", baseJsSource)

	// no rewriting if CodeLocation is provided as an argument
	codeLocationProvided :=
		`function signal(here = goog.callerLocation()) {}
const mySignal = signal('path/to/file.ts:25');
`
	testJs := SourceFile_FromCode("codeLocationProvided.js", codeLocationProvided)
	c.TestSame(c.Srcs(baseJs, testJs))

	codeLocationProvided2 :=
		`function signal(val, here = goog.callerLocation()) {}
const mySignal = signal(0, xid('path/to/file.ts:25'));
`
	testJs2 := SourceFile_FromCode("codeLocationProvided2.js", codeLocationProvided2)
	c.TestSame(c.Srcs(baseJs, testJs2))

	codeLocationProvided3 :=
		`function signal(val, here = goog.callerLocation()) {}
const mySignal = signal(
  () => foo() % 2 === 0,
  xid('path/to/file.ts:25'),
);
`
	testJs3 := SourceFile_FromCode("codeLocationProvided3.js", codeLocationProvided3)
	c.TestSame(c.Srcs(baseJs, testJs3))

	// no rewriting if CodeLocation is provided as an argument, regardless of the type of the arg
	codeLocationProvided4 :=
		`function signal(here = goog.callerLocation()) {}
const mySignal = signal(foo());
`
	testJs4 := SourceFile_FromCode("codeLocationProvided4.js", codeLocationProvided4)
	c.TestSame(c.Srcs(baseJs, testJs4))
}

// Java: RewriteCallerCodeLocationTest.testNoRewriteOnClassMethods()
func (c *RewriteCallerCodeLocationTest) TestNoRewriteOnClassMethods() {
	c.TestError(
		`class Foo {
  bar(here = goog.callerLocation()) {}
}
const baz = new Foo();
baz.bar()
`,
		RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR)
}

// Java: RewriteCallerCodeLocationTest.testTransitiveUsage()
func (c *RewriteCallerCodeLocationTest) TestTransitiveUsage() {
	// `foo` has a goog.callerLocation() default parameter.
	// When `bar` calls `foo` internally, we will not rewrite the call sites of `bar` unless
	// `bar` itself has a goog.callerLocation() default.
	c.Test(
		"function foo(val, here = goog.callerLocation()) {}\nfunction bar(val) {\n  foo(val);\n}\n// `bar(0)` will not get rewritten because `bar` does not have a goog.callerLocation()\n// default parameter\nbar(0)\n",
		`function foo(val, here = goog.callerLocation()) {}
function bar(val) {
  foo(val, goog.callerLocationIdInternalDoNotCallOrElse('testcode:3:2'));
}
bar(0)
`)

	c.Test(
		"function foo(val, here = goog.callerLocation()) {}\nfunction bar(val, here = goog.callerLocation()) {\n  foo(val);\n}\n// `bar(0)` will get rewritten because `bar` has a goog.callerLocation() default\n// parameter\nbar(0)\n",
		`function foo(val, here = goog.callerLocation()) {}
function bar(val, here = goog.callerLocation()) {
  foo(val, goog.callerLocationIdInternalDoNotCallOrElse('testcode:3:2'));
}
bar(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:7:0'));
`)
}

// Java: RewriteCallerCodeLocationTest.testErrorOnCallerLocationNotUsedAsDefaultParam()
func (c *RewriteCallerCodeLocationTest) TestErrorOnCallerLocationNotUsedAsDefaultParam() {
	// no rewriting if `goog.callerLocation()` is not used in a default parameter
	c.TestError(
		`function foo() {
  const x = goog.callerLocation();
}
`,
		RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR)
}

// Java: RewriteCallerCodeLocationTest.testDestructuringOptionsBag()
func (c *RewriteCallerCodeLocationTest) TestDestructuringOptionsBag() {
	// Error if `goog.callerLocation()` is used in an object literal
	c.TestError(
		`function foo({val1, val2, here = goog.callerLocation()}) {}
foo({val1:0, val2:0})
`,
		RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR)
}

// Java: RewriteCallerCodeLocationTest.testMultipleDefaultParameters()
func (c *RewriteCallerCodeLocationTest) TestMultipleDefaultParameters() {
	c.TestError(
		`function foo(val1 = 0, val2 = 5, here = goog.callerLocation()) {}
foo()
foo(1)
foo(1, 2)
`,
		RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR)

	c.TestError(
		`function foo(val1 = 0, val2 = 5, here = goog.callerLocation(), val3 = 10, val4 = 15) {}
foo()
foo(1)
foo(1, 2)
foo(1, 2, goog.callerLocationIdInternalDoNotCallOrElse('customString'))
foo(1, 2, goog.callerLocationIdInternalDoNotCallOrElse('customString'), 3)
foo(1, 2, goog.callerLocationIdInternalDoNotCallOrElse('customString'), 3, 4)
`,
		RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR)
}

// Java: RewriteCallerCodeLocationTest.testTaggedTemplateLiteralsDoesNotRewrite()
func (c *RewriteCallerCodeLocationTest) TestTaggedTemplateLiteralsDoesNotRewrite() {
	c.TestSame(
		"function signal(strings, here = goog.callerLocation()) {}\nconst tpl = signal`not passing a here argument explicitly`;\nconst tpl = signal`passing a here argument ${'here!'}`;\n")
}

// Java: RewriteCallerCodeLocationTest.testAnonymousFunctionCallErrors()
func (c *RewriteCallerCodeLocationTest) TestAnonymousFunctionCallErrors() {
	c.TestError(
		`const add = function(a, b, here = goog.callerLocation()) {
  return a + b;
}
add(2, 3);
`,
		RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR)
}

// Java: RewriteCallerCodeLocationTest.testModuleRewriting()
func (c *RewriteCallerCodeLocationTest) TestModuleRewriting() {
	c.Test(
		"function module$contents$main_signal(val, here = goog.callerLocation()) {\n  return val;\n}\nmodule$exports$main.signal = module$contents$main_signal;\n// `module$exports$fileName.signal` is the call-site function name after\n// ClosureRewriteModule runs\n(0,module$exports$main.signal)(0);\n",
		`function module$contents$main_signal(val, here = goog.callerLocation()) {
  return val;
}
module$exports$main.signal = module$contents$main_signal;
(0,module$exports$main.signal)(0, goog.callerLocationIdInternalDoNotCallOrElse('testcode:7:0'));
`)
}

// Java: RewriteCallerCodeLocationTest.testScopeRewriting()
func (c *RewriteCallerCodeLocationTest) TestScopeRewriting() {
	c.Test(
		`function signal(here = goog.callerLocation()) {}
class Button {
  method(signal) {
    signal(); // this signal() call shouldn't be rewritten
  }
}
signal();
`,
		`function signal(here = goog.callerLocation()) {}
class Button {
  method(signal) {
    signal(); // not rewritten
  }
}
signal(goog.callerLocationIdInternalDoNotCallOrElse('testcode:7:0'));
`)

	c.Test(
		`function signal() {}
function outter() {
  function signal(here = goog.callerLocation()) {}
  signal(); // this signal() call should be rewritten
}
signal();
`,
		`function signal() {}
function outter() {
  function signal(here = goog.callerLocation()) {}
  signal(goog.callerLocationIdInternalDoNotCallOrElse('testcode:4:2')); // rewritten
}
signal();
`)
}

// Java: RewriteCallerCodeLocationTest.testFunctionWithCallerLocationParameterInExterns()
func (c *RewriteCallerCodeLocationTest) TestFunctionWithCallerLocationParameterInExterns() {
	c.Test(
		c.Externs(
			`/**
 * @fileoverview
 * @typeSummary
 * @externs
 */
var goog = {};
function signal(x = goog.callerLocation()) {};
`),
		c.Srcs("signal();"),
		c.Expected("signal(goog.callerLocationIdInternalDoNotCallOrElse('testcode:1:0'));"))
}

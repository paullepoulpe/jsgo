// Port of com/google/javascript/jscomp/TypeCheckNoTranspileTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/TypeCheckNoTranspileTest.java)

/*
 * Copyright 2017 The Closure Compiler Authors.
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

// TypeCheckNoTranspileTest tests the type checker.
type TypeCheckNoTranspileTest struct {
	typeCheckTests_noFixture
}

// Java: TypeCheckNoTranspileTest()
func NewTypeCheckNoTranspileTest() *TypeCheckNoTranspileTest {
	c := &TypeCheckNoTranspileTest{}
	return c
}

func TestTypeCheckNoTranspileTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckNoTranspileTest() })
}

// Java: TypeCheckNoTranspileTest.testCorrectSubtyping_ofRecursiveTemplateType()
func (c *TypeCheckNoTranspileTest) TestCorrectSubtyping_ofRecursiveTemplateType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @template T */\nclass Base { }\n\n/** @extends {Base<!Child>} */\nclass Child extends Base { }\n\n// Confirm that `Child` is seen as a subtype of `Base<Child>`.\nconst /** !Base<!Child> */ x = new Child();\n").
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrowInferredReturn()
func (c *TypeCheckNoTranspileTest) TestArrowInferredReturn() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`let fn = () => {
  return 1;
};
var /** null */ x = fn();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrowBlocklessInferredReturn()
func (c *TypeCheckNoTranspileTest) TestArrowBlocklessInferredReturn() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`let fn = () => 1;
var /** null */ x = fn();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrowRightScopeForBody()
func (c *TypeCheckNoTranspileTest) TestArrowRightScopeForBody() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string} */ let a = 's';
/**
  @param {number} a
  @return {null}
*/
let fn = (a) => {
  return a;
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrowRightBodyScopeForBlocklessBody()
func (c *TypeCheckNoTranspileTest) TestArrowRightBodyScopeForBlocklessBody() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string} */ let a = 's';
/**
  @param {number} a
  @return {null}
*/
let fn = (a) => a
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrowCorrectThis()
func (c *TypeCheckNoTranspileTest) TestArrowCorrectThis() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @this {String} */ function fn() {
  /**
    @return {null}
  */
  let fn = () => {
    return this;
  }
}
`).
		AddDiagnostic(`inconsistent return type
found   : String
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrowBlocklessCorrectThis()
func (c *TypeCheckNoTranspileTest) TestArrowBlocklessCorrectThis() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @this {String} */ function fn() {
  /**
    @return {null}
  */
  let fn = () => this;
}
`).
		AddDiagnostic(`inconsistent return type
found   : String
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrowCorrectArguments()
func (c *TypeCheckNoTranspileTest) TestArrowCorrectArguments() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @this {String} */ function fn() {
  {
    /** @type {number} */ let arguments = 1;
    /** @return {null} */
    let fn = () => {
      return arguments;
    }
  }
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrowBlocklessCorrectArguments()
func (c *TypeCheckNoTranspileTest) TestArrowBlocklessCorrectArguments() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @this {String} */ function fn() {
  {
    /** @type {number} */ let arguments = 1;
    /** @return {null} */
    let fn = () => arguments;
  }
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrowCorrectInheritsArguments()
func (c *TypeCheckNoTranspileTest) TestArrowCorrectInheritsArguments() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("/** @type {!Arguments} */ var arguments;").
		AddSource(`/** @this {String} */ function fn() {
  {
    /** @return {null} */
    let fn = () => {
      return arguments;
    }
  }
}
`).
		AddDiagnostic(`inconsistent return type
found   : Arguments
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrowBlocklessCorrectInheritsArguments()
func (c *TypeCheckNoTranspileTest) TestArrowBlocklessCorrectInheritsArguments() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("/** @type {!Arguments} */ var arguments;").
		AddSource(`/** @this {String} */ function fn() {
  {
    /** @return {null} */
    let fn = () => arguments;
  }
}
`).
		AddDiagnostic(`inconsistent return type
found   : Arguments
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncArrow_withValidBlocklessReturn_isAllowed()
func (c *TypeCheckNoTranspileTest) TestAsyncArrow_withValidBlocklessReturn_isAllowed() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesPromiseProvider(/** function():!Promise<number> */ getPromise) {}
takesPromiseProvider(async () => 1);
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncArrow_withInvalidBlocklessReturn_isError()
func (c *TypeCheckNoTranspileTest) TestAsyncArrow_withInvalidBlocklessReturn_isError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesPromiseProvider(/** function():!Promise<string> */ getPromise) {}
takesPromiseProvider(async () => 1);
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: (IThenable<string>|string)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncArrow_withInferredReturnType_ofValidUnionType_isAllowed()
func (c *TypeCheckNoTranspileTest) TestAsyncArrow_withInferredReturnType_ofValidUnionType_isAllowed() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function():(number|!Promise<string>)} getPromise */
function takesPromiseProvider(getPromise) {}

takesPromiseProvider(async () => '');
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncArrow_withInferredReturnType_ofInvalidUnionType_isError()
func (c *TypeCheckNoTranspileTest) TestAsyncArrow_withInferredReturnType_ofInvalidUnionType_isError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function():(number|!Promise<string>)} getPromise */
function takesPromiseProvider(getPromise) {}

takesPromiseProvider(async () => true);
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: (IThenable<string>|string)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefOfPropertyInBlock()
func (c *TypeCheckNoTranspileTest) TestTypedefOfPropertyInBlock() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("/** @interface */ function Foo() {}").
		AddSource(`/** @constructor */
function Bar(/** !Foo */ foo) {
  /** @type {!Foo} */
  this.foo = foo;
  {
    /** @typedef {boolean} */
    this.foo.bar;
    (() => this.foo.bar)();
  }
}
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefOfPropertyInFunctionScope()
func (c *TypeCheckNoTranspileTest) TestTypedefOfPropertyInFunctionScope() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("/** @interface */ function Foo() {}").
		AddSource(`/** @constructor */
function Bar(/** !Foo */ foo) {
  /** @type {!Foo} */
  this.foo = foo;
  /** @typedef {boolean} */
  this.foo.bar;
  {
    (() => this.foo.bar)();
  }
}
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayLitSpread()
func (c *TypeCheckNoTranspileTest) TestArrayLitSpread() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<string>} */
const strings = [];
/** @type {!Array<number>} */
const numbers = [...strings]; // This should generate an error
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayLitSpreadNonIterable()
func (c *TypeCheckNoTranspileTest) TestArrayLitSpreadNonIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<number>} */
const numbers = [...1];
`).
		AddDiagnostic(`Spread operator only applies to Iterable types
found   : number
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypecheckExpressionInArrayLitSpread()
func (c *TypeCheckNoTranspileTest) TestTypecheckExpressionInArrayLitSpread() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<string>} */
const strings = [];
/** @type {!Array<number>} */
let numbers = [];
const a = [...(numbers = strings)];
`).
		AddDiagnostic(`assignment
found   : Array<string>
required: Array<number>
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testInferTypesFromExpressionInArrayLitSpread()
func (c *TypeCheckNoTranspileTest) TestInferTypesFromExpressionInArrayLitSpread() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<string>} */
const strings = [];
let inferred = 1;
const a = [...(inferred = strings)];
/** @type {null} */
const n = inferred;
`).
		AddDiagnostic(`initializing variable
found   : Array<string>
required: null
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testSpreadAndFollowingParametersNotTypeChecked()
func (c *TypeCheckNoTranspileTest) TestSpreadAndFollowingParametersNotTypeChecked() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
 * @param {number} num
 * @param {string} str
 * @param {boolean} bool
 */
function use(num, str, bool) {}

`).
		AddSource(`/** @type {!Array<null>} */ const nulls = [];
use(1, ...nulls, null, null);
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testSpreadArgumentTypeCheckedForVarArgs()
func (c *TypeCheckNoTranspileTest) TestSpreadArgumentTypeCheckedForVarArgs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
 * @param {number} num
 * @param {...string} var_args
 */
function use(num, var_args) {}

`).
		AddSource(`/** @type {!Array<null>} */ const nulls = [];
use(1, ...nulls);
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testSpreadArgumentBackInference()
func (c *TypeCheckNoTranspileTest) TestSpreadArgumentBackInference() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
 * @param {number} num
 * @param {...{prop: number}} var_args
 */
function use(num, var_args) {}

`).
		AddSource("use(1, ...[{}]);").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testTooManyNonSpreadParameters()
func (c *TypeCheckNoTranspileTest) TestTooManyNonSpreadParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
 * @param {number} num
 * @param {string} str
 * @param {boolean} bool
 */
function use(num, str, bool) {}

`).
		AddSource(`/** @type {!Array<*>} */ const unusables = [];
use(1, 'hi', ...unusables, null, null);
`).
		AddDiagnostic("Function use: called with at least 4 argument(s). Function requires at least 3 argument(s) and no more than 3 argument(s).").
		Run()
}

// Java: TypeCheckNoTranspileTest.testArgumentSpreadDoesNotBlockTypeCheckOfInitialParameters()
func (c *TypeCheckNoTranspileTest) TestArgumentSpreadDoesNotBlockTypeCheckOfInitialParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
 * @param {number} num
 * @param {string} str
 * @param {boolean} bool
 */
function use(num, str, bool) {}

`).
		AddSource("use('should be number', ...[]);").
		AddDiagnostic(`actual parameter 1 of use does not match formal parameter
found   : string
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testArgumentSpreadNonIterable()
func (c *TypeCheckNoTranspileTest) TestArgumentSpreadNonIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("function use(x) {}").
		AddSource("use(...1);").
		AddDiagnostic(`Spread operator only applies to Iterable types
found   : number
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArgumentSpreadNonIterable_optChainCall()
func (c *TypeCheckNoTranspileTest) TestArgumentSpreadNonIterable_optChainCall() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("function use(x) {}").
		AddSource("use?.(...1);").
		AddDiagnostic(`Spread operator only applies to Iterable types
found   : number
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypecheckExpressionInArgumentSpread()
func (c *TypeCheckNoTranspileTest) TestTypecheckExpressionInArgumentSpread() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("function use(x) {}").
		AddSource(`/** @type {!Array<string>} */
const strings = [];
/** @type {!Array<number>} */
let numbers = [];
use(...(numbers = strings));
`).
		AddDiagnostic(`assignment
found   : Array<string>
required: Array<number>
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testInferTypesFromExpressionInArgumentSpread()
func (c *TypeCheckNoTranspileTest) TestInferTypesFromExpressionInArgumentSpread() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("function use(x) {}").
		AddSource(`/** @type {!Array<string>} */
const strings = [];
let inferred = 1;
use(...(inferred = strings));
/** @type {null} */
const n = inferred;
`).
		AddDiagnostic(`initializing variable
found   : Array<string>
required: null
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testOnlyRestParameterWithoutJSDocCalledWithNoArgs()
func (c *TypeCheckNoTranspileTest) TestOnlyRestParameterWithoutJSDocCalledWithNoArgs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function use(...numbers) {}
use(); // no args provided in call - should be OK
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testBadRestJSDoc()
func (c *TypeCheckNoTranspileTest) TestBadRestJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} numbers */
function f(...numbers) {
  var /** null */ n = numbers;
}
`).
		AddDiagnostic(`initializing variable
found   : Array<number>
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOnlyRestParameterWithoutJSDocCalledWithArgs()
func (c *TypeCheckNoTranspileTest) TestOnlyRestParameterWithoutJSDocCalledWithArgs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function use(...numbers) {}
use(1, 'hi', {});
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOnlyRestParameterWithJSDocCalledWithNoArgs()
func (c *TypeCheckNoTranspileTest) TestOnlyRestParameterWithJSDocCalledWithNoArgs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {...number} numbers
 */
function use(...numbers) {}
use(); // no args provided in call - should be OK
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOnlyRestParameterWithJSDocCalledWithGoodArgs()
func (c *TypeCheckNoTranspileTest) TestOnlyRestParameterWithJSDocCalledWithGoodArgs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {...number} numbers
 */
function use(...numbers) {}
use(1, 2, 3);
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOnlyRestParameterWithJSDocCalledWithBadArg()
func (c *TypeCheckNoTranspileTest) TestOnlyRestParameterWithJSDocCalledWithBadArg() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {...number} numbers
 */
function use(...numbers) {}
use(1, 'hi', 3);
`).
		AddDiagnostic(`actual parameter 2 of use does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testNormalAndRestParameterWithJSDocCalledWithOneArg()
func (c *TypeCheckNoTranspileTest) TestNormalAndRestParameterWithJSDocCalledWithOneArg() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} str
 * @param {...number} numbers
 */
function use(str, ...numbers) {}
use('hi'); // no rest args provided in call - should be OK
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testNormalAndRestParameterWithJSDocCalledWithGoodArgs()
func (c *TypeCheckNoTranspileTest) TestNormalAndRestParameterWithJSDocCalledWithGoodArgs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} str
 * @param {...number} numbers
 */
function use(str, ...numbers) {}
use('hi', 2, 3);
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOnlyRestParameterWithJSDocCalledWithBadNormalArg()
func (c *TypeCheckNoTranspileTest) TestOnlyRestParameterWithJSDocCalledWithBadNormalArg() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} str
 * @param {...number} numbers
 */
function use(str, ...numbers) {}
use(1, 2, 3);
`).
		AddDiagnostic(`actual parameter 1 of use does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOnlyRestParameterWithJSDocCalledWithBadRestArg()
func (c *TypeCheckNoTranspileTest) TestOnlyRestParameterWithJSDocCalledWithBadRestArg() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} str
 * @param {...number} numbers
 */
function use(str, ...numbers) {}
use('hi', 'there', 3);
`).
		AddDiagnostic(`actual parameter 2 of use does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testRestParameterInCallbackIsInferred()
func (c *TypeCheckNoTranspileTest) TestRestParameterInCallbackIsInferred() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(...number)} callback */
function f(callback) {}

f((...strings) => {
  const /** null */ n = strings; // verify that this causes a type mismatch
});
`).
		AddDiagnostic(`initializing variable
found   : Array<number>
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testExponent1()
func (c *TypeCheckNoTranspileTest) TestExponent1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function fn(someUnknown) {
  var x = someUnknown ** 2; // infer the result
  var /** null */ y = x;
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testExponent2()
func (c *TypeCheckNoTranspileTest) TestExponent2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function fn(someUnknown) {
  var x = someUnknown;
  x **= 2; // infer the result
  var /** null */ y = x;
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testExponent3()
func (c *TypeCheckNoTranspileTest) TestExponent3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function fn(someUnknown) {
  var y = true ** 3;
}
`).
		AddDiagnostic(`left operand
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testExponent4()
func (c *TypeCheckNoTranspileTest) TestExponent4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function fn(someUnknown) {
  var y = 1; y **= true;
}
`).
		AddDiagnostic(`right operand
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDuplicateCatchVarName()
func (c *TypeCheckNoTranspileTest) TestDuplicateCatchVarName() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`try { throw 1; } catch (/** @type {number} */ err) {}
try { throw 'error'; } catch (/** @type {string} */ err) {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefFieldInLoopLocal()
func (c *TypeCheckNoTranspileTest) TestTypedefFieldInLoopLocal() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{num: number, maybeNum: ?number}} */
let XType;

/** @param {!Array<!XType>} xlist */
function f(xlist) {
  for (let i = 0; i < xlist.length; i++) {
    /** @type {!XType} */
    const x = xlist[i];
    if (x.maybeNum === null) {
      continue;
    }
    x.num = x.maybeNum;
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefFieldInLoopGlobal()
func (c *TypeCheckNoTranspileTest) TestTypedefFieldInLoopGlobal() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`/** @typedef {{num: number, maybeNum: ?number}} */
let XType;

/** @type {!Array<!XType>} */
const xlist = [{maybeNum: null, num: 0}, {maybeNum: 1, num: 1}];

for (let i = 0; i < xlist.length; i++) {
  /** @type {!XType} */
  const x = xlist[i];
  if (x.maybeNum === null) {
    continue;
  }
// TODO(b/78364240): Compiler should realize that x.maybeNum must be a number here
  x.num = x.maybeNum;
}
`).
		AddDiagnostic(`assignment to property num of x
found   : (null|number)
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefAliasValueTypeIsUndefined()
func (c *TypeCheckNoTranspileTest) TestTypedefAliasValueTypeIsUndefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {number} */
var MyNumber;
var ns = {};
ns.MyNumber = MyNumber;
/** @type {string} */ (ns.MyNumber);
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefAliasOfLocalTypedef()
func (c *TypeCheckNoTranspileTest) TestTypedefAliasOfLocalTypedef() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @typedef {number} */
  var MyNumber;
  /** @const */
  var Alias = MyNumber;
  /** @type {Alias} */
  var x = 'x';
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefDestructuredAlias()
func (c *TypeCheckNoTranspileTest) TestTypedefDestructuredAlias() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  const ns = {};
  /** @typedef {number} */
  ns.MyNumber;
  const {MyNumber: Alias} = ns;
  /** @type {Alias} */
  var x = 'x';
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefDestructuredAlias_deeplyNested()
func (c *TypeCheckNoTranspileTest) TestTypedefDestructuredAlias_deeplyNested() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  const outer = {};
  /** @const */
  outer.inner = {};
  /** @typedef {number} */
  outer.inner.MyNumber;
  const alias = {};
  ({inner: /** @const */ alias.ns} = outer);
  /** @type {alias.ns.MyNumber} */
  var x = 'x';
}
`).
		AddDiagnostic(`Bad type annotation. Unknown type alias.ns.MyNumber
It's possible that 'alias.ns.MyNumber' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefLocalQualifiedName()
func (c *TypeCheckNoTranspileTest) TestTypedefLocalQualifiedName() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @const */
  var ns = {};
  /** @typedef {number} */
  ns.MyNumber;
  /** @type {ns.MyNumber} */
  var x = 'x';
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefLocalQualifiedNameAlias()
func (c *TypeCheckNoTranspileTest) TestTypedefLocalQualifiedNameAlias() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @typedef {number} */
  var MyNumber;
  /** @const */
  var ns = {};
  /** @const */
  ns.MyNumber = MyNumber;
  /** @type {ns.MyNumber} */
  var x = 'x';
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefLocalAliasOfGlobalTypedef()
func (c *TypeCheckNoTranspileTest) TestTypedefLocalAliasOfGlobalTypedef() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {number} */
var MyNumber;
function f() {
  /** @const */ var Alias = MyNumber;
  var /** Alias */ x = 'x';
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefOnAliasedNamespace()
func (c *TypeCheckNoTranspileTest) TestTypedefOnAliasedNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const ns = {};
/** @const */ ns.bar = 'x';
/** @typedef {number} */
ns.MyNumber;
const alias = ns;
/** @const */ alias.foo = 42
/** @type {alias.MyNumber} */ const x = 'str';
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefOnLocalAliasedNamespace()
func (c *TypeCheckNoTranspileTest) TestTypedefOnLocalAliasedNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  const ns = {};
  /** @typedef {number} */
  ns.MyNumber;
  const alias = ns;
  /** @const */ alias.foo = 42
  /** @type {alias.MyNumber} */ const x = 'str';
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefOnClassSideInheritedSubtypeInaccessible()
func (c *TypeCheckNoTranspileTest) TestTypedefOnClassSideInheritedSubtypeInaccessible() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {}
/** @typedef {number} */
Base.MyNumber;
class Sub extends Base {}
/** @type {Sub.MyNumber} */ let x;
`).
		AddDiagnostic(`Bad type annotation. Unknown type Sub.MyNumber
It's possible that 'Sub.MyNumber' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testBlockScopedVarInLoop1()
func (c *TypeCheckNoTranspileTest) TestBlockScopedVarInLoop1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`/** @constructor */ function Node() {};
function g(/** Node */ n){
  n.foo = {bar: 3};
}
function f(/** !Array<!Node> */ arr){
  for (var i = 0; i < arr.length; i++) {
    const tile = arr[i];
    const bar = tile.foo.bar;
// this assignment shouldn't cause 'tile.foo' to be inferred as undefined above.
    tile.foo = undefined
  }
}
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckNoTranspileTest.testBlockScopedVarInLoop2()
func (c *TypeCheckNoTranspileTest) TestBlockScopedVarInLoop2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`while (true) {
  let num;
  let /** undefined */ y = num;
// null assignment shouldn't make us think num could be null on the previous line.
  num = null;
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testBlockScopedVarInLoop3()
func (c *TypeCheckNoTranspileTest) TestBlockScopedVarInLoop3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number */ n) {}

function f(/** {num: ?number} */ obj) {
  for (const _ in {}) {
    const alias = obj;
    if (alias.num === null) {
     continue;
    }
    takesNumber(alias.num);
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGlobalEnumWithLet()
func (c *TypeCheckNoTranspileTest) TestGlobalEnumWithLet() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ // type defaults to {number}
let E = {A: 1, B: 2};

/**
 * @param {E} x
 * @return {number}
 */
function f(x) {return x}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGlobalEnumWithConst()
func (c *TypeCheckNoTranspileTest) TestGlobalEnumWithConst() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ // type defaults to {number}
const E = {A: 1, B: 2};

/**
 * @param {E} x
 * @return {number}
 */
function f(x) {return x}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testLocalEnumWithLet()
func (c *TypeCheckNoTranspileTest) TestLocalEnumWithLet() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`{
  /** @enum */ // type defaults to {number}
  let E = {A: 1, B: 2};

  /**
   * @param {E} x
   * @return {number}
   */
  function f(x) {return x}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testLocalEnumWithConst()
func (c *TypeCheckNoTranspileTest) TestLocalEnumWithConst() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`{
  /** @enum */ // type defaults to {number}
  const E = {A: 1, B: 2};

  /**
   * @param {E} x
   * @return {number}
   */
  function f(x) {return x}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGlobalTypedefWithLet()
func (c *TypeCheckNoTranspileTest) TestGlobalTypedefWithLet() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {number} */
let Bar;
/** @param {Bar} x */
function f(x) {}
f('3');
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testLocalTypedefWithLet()
func (c *TypeCheckNoTranspileTest) TestLocalTypedefWithLet() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`{
  /** @typedef {number} */
  let Bar;
  /** @param {Bar} x */
  function f(x) {}
  f('3');
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testConstWrongType()
func (c *TypeCheckNoTranspileTest) TestConstWrongType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {number} */ const x = 'hi';").
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testLetWrongType()
func (c *TypeCheckNoTranspileTest) TestLetWrongType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {number} */ let x = 'hi';").
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testLetInitializedToUndefined1()
func (c *TypeCheckNoTranspileTest) TestLetInitializedToUndefined1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("let foo; let /** number */ bar = foo;").
		AddDiagnostic(`initializing variable
found   : undefined
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testLetInitializedToUndefined2()
func (c *TypeCheckNoTranspileTest) TestLetInitializedToUndefined2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("let /** number */ foo; let /** number */ bar = foo;").
		Run()
}

// Java: TypeCheckNoTranspileTest.testLetInitializedToUndefined3()
func (c *TypeCheckNoTranspileTest) TestLetInitializedToUndefined3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`let foo;
/** @return {number} */
function getFoo() {
  return foo;
}
function setFoo(/** number */ num) {
  foo = num;
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf1()
func (c *TypeCheckNoTranspileTest) TestForOf1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Iterable} */ var it; for (var elem of it) {}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf2()
func (c *TypeCheckNoTranspileTest) TestForOf2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesString(/** string */ s) {}
/** @type {!Iterable<string>} */ var it;
for (var elem of it) { takesString(elem); }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf3()
func (c *TypeCheckNoTranspileTest) TestForOf3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesString(/** string */ s) {}
/** @type {!Iterable<number>} */ var it;
for (var elem of it) { takesString(elem); }
`).
		AddDiagnostic(`actual parameter 1 of takesString does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf4()
func (c *TypeCheckNoTranspileTest) TestForOf4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Iterable} */ var it; var obj = {}; for (obj.elem of it) {}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf5()
func (c *TypeCheckNoTranspileTest) TestForOf5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesString(/** string */ s) {}

function f(/** !Iterable<number> */ it) {
  var obj = {};
  for (obj.elem of it) {
    takesString(obj.elem);
  }
}
`).
		AddDiagnostic(`actual parameter 1 of takesString does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_wrongLoopVarType1()
func (c *TypeCheckNoTranspileTest) TestForOf_wrongLoopVarType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<number>} */
var numArray = [1, 2];
/** @type {string} */
var elem = '';
for (elem of numArray) {
}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_wrongLoopVarType2()
func (c *TypeCheckNoTranspileTest) TestForOf_wrongLoopVarType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<number>} */
var numArray = [1, 2];
for (let /** string */ elem of numArray) {
}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_wrongLoopVarType3()
func (c *TypeCheckNoTranspileTest) TestForOf_wrongLoopVarType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("for (var /** number */ x of 3) {}").
		AddDiagnostic(`Can only iterate over a (non-null) Iterable type
found   : number
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_wrongLoopVarType4a()
func (c *TypeCheckNoTranspileTest) TestForOf_wrongLoopVarType4a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<!Object>} */
var arr = [1, 2];
for (let /** ?Object */ elem of arr) {
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_wrongLoopVarType4b()
func (c *TypeCheckNoTranspileTest) TestForOf_wrongLoopVarType4b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<?Object>} */
var arr = [1, 2];
for (let /** !Object */ elem of arr) {
}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : (Object|null)
required: Object
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_wrongLoopVarType5()
func (c *TypeCheckNoTranspileTest) TestForOf_wrongLoopVarType5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<number>} */
var arr = [1, 2];
let n = null;
for (n of arr) {}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_wrongLoopVarType6a()
func (c *TypeCheckNoTranspileTest) TestForOf_wrongLoopVarType6a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`/** @type {!Array<number>} */
var arr = [1, 2, 3];
let /** string */ n = 'foo'; // n in global scope
for (let /** number */ n of arr) { // n in for of scope
  let /** null */ n = null; // n in inner block scope
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_wrongLoopVarType6b()
func (c *TypeCheckNoTranspileTest) TestForOf_wrongLoopVarType6b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`/** @type {!Array<string>} */
var arr = ['foo', 'bar'];
let /** string */ n = 'foo'; // n in global scope
for (let /** number */ n of arr) { // n in for of scope
  let /** null */ n = null; // n in inner block scope
}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_wrongLoopVarType7()
func (c *TypeCheckNoTranspileTest) TestForOf_wrongLoopVarType7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Iterable<string>} */ var it;
var /** !Object<string, number> */ obj = {};
for (obj['x'] of it) {}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_wrongLoopVarType8()
func (c *TypeCheckNoTranspileTest) TestForOf_wrongLoopVarType8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Iterable<string>} */ var it;
const /** @type {{x: number}} */ obj = {x: 5};
for (obj.x of it) {}
`).
		AddDiagnostic(`assignment to property x of obj
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_illegalPropertyCreation()
func (c *TypeCheckNoTranspileTest) TestForOf_illegalPropertyCreation() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Iterable<string>} */ var it;
const /** @struct */ obj = {};
for (obj.x of it) {}
`).
		AddDiagnostic("Cannot add a property to a struct instance after it is constructed. (If you already declared the property, make sure to give it a type.)").
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_badInterfaceMemberCreation()
func (c *TypeCheckNoTranspileTest) TestForOf_badInterfaceMemberCreation() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
function Foo() {}
for (Foo.prototype.bar of []) {}
`).
		AddDiagnostic("interface members can only be empty property declarations, empty functions, or goog.abstractMethod").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_badEnumCreation()
func (c *TypeCheckNoTranspileTest) TestForOf_badEnumCreation() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("for (var /** @enum */ myEnum of []) {}").
		AddDiagnostic("enum initializer must be an object literal or an enum").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_array1()
func (c *TypeCheckNoTranspileTest) TestForOf_array1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("for (var elem of [1, 2]) {}").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_array2()
func (c *TypeCheckNoTranspileTest) TestForOf_array2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<number>} */ var arr = [1, 2];
function takesString(/** string */ s) {}
for (var elem of arr) { takesString(elem); }
`).
		AddDiagnostic(`actual parameter 1 of takesString does not match formal parameter
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_array3()
func (c *TypeCheckNoTranspileTest) TestForOf_array3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<number>} */ var arr = [1, 2];
function takesNumber(/** number */ n) {}
for (var elem of arr) { takesNumber(elem); }
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_string1()
func (c *TypeCheckNoTranspileTest) TestForOf_string1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesString(/** string */ s) {}
for (var ch of 'a string') { takesString(ch); }
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_string2()
func (c *TypeCheckNoTranspileTest) TestForOf_string2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number */ n) {}
for (var ch of 'a string') { takesNumber(ch); }
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : string
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_StringObject1()
func (c *TypeCheckNoTranspileTest) TestForOf_StringObject1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesString(/** string */ s) {}
for (var ch of new String('boxed')) { takesString(elem); }
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_StringObject2()
func (c *TypeCheckNoTranspileTest) TestForOf_StringObject2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number */ n) {}
for (var ch of new String('boxed')) { takesNumber(elem); }
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_forbidsAsyncIterable()
func (c *TypeCheckNoTranspileTest) TestForOf_forbidsAsyncIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!AsyncIterable<string>} asyncIterable */
function f(asyncIterable) {
  for (var elem of asyncIterable) {}
}
`).
		IncludeDefaultExterns().
		AddDiagnostic(`Can only iterate over a (non-null) Iterable type
found   : AsyncIterable<string,?,?>
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_iterableTypeIsNotFirstTemplateType()
func (c *TypeCheckNoTranspileTest) TestForOf_iterableTypeIsNotFirstTemplateType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number */ n) {}

/**
 * @constructor
 * @implements {Iterable<T>}
 * @template S, T
 */
function MyIterable() {}

// Note that 'mi' is an Iterable<string>, not an Iterable<number>.
/** @type {!MyIterable<number, string>} */
var mi;

for (var t of mi) { takesNumber(t); }
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_unionType1()
func (c *TypeCheckNoTranspileTest) TestForOf_unionType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number */ n) {}
/** @param {(!Array<string>|undefined)} arr */
function f(arr) {
  for (let x of (arr || [])) {
    takesNumber(x);
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_unionType2()
func (c *TypeCheckNoTranspileTest) TestForOf_unionType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {(number|undefined)} n */
function f(n) {
  for (let x of (n || [])) {}
}
`).
		AddDiagnostic(`Can only iterate over a (non-null) Iterable type
found   : (Array<?>|number)
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_unionType_stringAndArray()
func (c *TypeCheckNoTranspileTest) TestForOf_unionType_stringAndArray() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`function takesNull(/** null */ n) {}

/** @param {string|!Array<number>} param */
function f(param) {
  for (let x of param) {
    takesNull(x);
  }
}
`).
		AddDiagnostic(`actual parameter 1 of takesNull does not match formal parameter
found   : (number|string)
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_unionType_readonlyArrayAndArray()
func (c *TypeCheckNoTranspileTest) TestForOf_unionType_readonlyArrayAndArray() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`function takesNull(/** null */ n) {}

/** @param {!ReadonlyArray<number>|!Array<string>} param */
function f(param) {
  for (let x of param) {
    takesNull(x);
  }
}
`).
		AddDiagnostic(`actual parameter 1 of takesNull does not match formal parameter
found   : (number|string)
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_nullable()
func (c *TypeCheckNoTranspileTest) TestForOf_nullable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {?Iterable} */ var it; for (var elem of it) {}").
		AddDiagnostic(`Can only iterate over a (non-null) Iterable type
found   : (Iterable|null)
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_null()
func (c *TypeCheckNoTranspileTest) TestForOf_null() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {null} */ var it = null; for (var elem of it) {}").
		AddDiagnostic(`Can only iterate over a (non-null) Iterable type
found   : null
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_maybeUndefined()
func (c *TypeCheckNoTranspileTest) TestForOf_maybeUndefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Iterable|undefined} */ var it; for (var elem of it) {}").
		AddDiagnostic(`Can only iterate over a (non-null) Iterable type
found   : (Iterable|undefined)
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_let()
func (c *TypeCheckNoTranspileTest) TestForOf_let() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Iterable} */ let it; for (let elem of it) {}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testForOf_const()
func (c *TypeCheckNoTranspileTest) TestForOf_const() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Iterable} */ const it = []; for (const elem of it) {}").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testImplicitCastInForOf()
func (c *TypeCheckNoTranspileTest) TestImplicitCastInForOf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */ function Element() {};
/**
 * @type {string}
 * @implicitCast
 */
Element.prototype.innerHTML;
`).
		AddSource(`/** @param {?Element} element
 * @param {!Array<string|number>} texts
 */
function f(element, texts) {
  for (element.innerHTML of texts) {};
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator1()
func (c *TypeCheckNoTranspileTest) TestGenerator1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Generator<?>} */ function* gen() {}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator2()
func (c *TypeCheckNoTranspileTest) TestGenerator2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Generator<number>} */ function* gen() { yield 1; }").
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator3()
func (c *TypeCheckNoTranspileTest) TestGenerator3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Generator<string>} */ function* gen() {  yield 1; }").
		AddDiagnostic(`Yielded type does not match declared return type.
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator4()
func (c *TypeCheckNoTranspileTest) TestGenerator4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Generator} */ // treat Generator as Generator<?>
function* gen() {
  yield 1;
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator5()
func (c *TypeCheckNoTranspileTest) TestGenerator5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Generator<{a: number, b: string}>} */
function *gen() {
  yield {a: 3, b: '4'};
}
var g = gen();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator6()
func (c *TypeCheckNoTranspileTest) TestGenerator6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`/** @return {!Generator<string>} */
function* gen() {
}
var g = gen();
var /** number */ n = g.next().value;
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_nextWithParameter()
func (c *TypeCheckNoTranspileTest) TestGenerator_nextWithParameter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`/** @return {!Generator<number>} */
function* gen() {
  var x = yield 1;
  yield x + 2;
}
var g = gen();
var /** number */ n = g.next().value; // 1
var /** number */ k = g.next('').value;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_yieldUndefined1()
func (c *TypeCheckNoTranspileTest) TestGenerator_yieldUndefined1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Generator<undefined>} */
function* gen() {
  yield undefined;
  yield; // yield undefined
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_yieldUndefined2()
func (c *TypeCheckNoTranspileTest) TestGenerator_yieldUndefined2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`/** @return {!Generator<number>} */
function* gen() {
  yield; // yield undefined
}
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : undefined
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_returnsIterable1()
func (c *TypeCheckNoTranspileTest) TestGenerator_returnsIterable1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Iterable<?>} */ function *gen() {}").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_returnsIterable2()
func (c *TypeCheckNoTranspileTest) TestGenerator_returnsIterable2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Iterable<string>} */ function* gen() {  yield 1; }").
		AddDiagnostic(`Yielded type does not match declared return type.
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_returnsIteratorLike1()
func (c *TypeCheckNoTranspileTest) TestGenerator_returnsIteratorLike1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!IteratorLike<?>} */ function *gen() {}").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_returnsIteratorLike2()
func (c *TypeCheckNoTranspileTest) TestGenerator_returnsIteratorLike2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!IteratorLike<string>} */ function* gen() {  yield 1; }").
		AddDiagnostic(`Yielded type does not match declared return type.
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_returnsStructuralIteratorLike()
func (c *TypeCheckNoTranspileTest) TestGenerator_returnsStructuralIteratorLike() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
function StructuralIterator() {}
/** @return {!IIterableResult<string>} */
StructuralIterator.prototype.next = function() {};

/** @return {!StructuralIterator} */
function* gen() { yield 1; }
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_returnsObjectLiteralIteratorLike()
func (c *TypeCheckNoTranspileTest) TestGenerator_returnsObjectLiteralIteratorLike() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{next: function(): !IIterableResult<string>}} */
var StructuralIterator;

/** @return {!StructuralIterator} */
function* gen() { yield 1; }
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_returnsIteratorIterable()
func (c *TypeCheckNoTranspileTest) TestGenerator_returnsIteratorIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!IteratorIterable<?>} */ function *gen() {}").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_cantReturnArray()
func (c *TypeCheckNoTranspileTest) TestGenerator_cantReturnArray() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Array<?>} */ function *gen() {}").
		AddDiagnostic(`A generator function must return a (supertype of) Generator
found   : Array<?>
required: Generator
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_notAConstructor()
func (c *TypeCheckNoTranspileTest) TestGenerator_notAConstructor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Generator<number>} */
function* gen() {
  yield 1;
}
var g = new gen;
`).
		AddDiagnostic("cannot instantiate non-constructor, found type: function(): Generator<number,?,?>").
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_noDeclaredReturnType1()
func (c *TypeCheckNoTranspileTest) TestGenerator_noDeclaredReturnType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function *gen() {} var /** !Generator<?> */ g = gen();").
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_noDeclaredReturnType2()
func (c *TypeCheckNoTranspileTest) TestGenerator_noDeclaredReturnType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function *gen() {} var /** !Generator<number> */ g = gen();").
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_noDeclaredReturnType3()
func (c *TypeCheckNoTranspileTest) TestGenerator_noDeclaredReturnType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function *gen() {
  yield 1;
  yield 2;
}
var /** string */ g = gen().next().value;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_return1()
func (c *TypeCheckNoTranspileTest) TestGenerator_return1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Generator<number, number>} */ function *gen() { return 1; }").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_return_noExplicitReturnStatement()
func (c *TypeCheckNoTranspileTest) TestGenerator_return_noExplicitReturnStatement() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Generator<number, string>} */ function *gen() {}").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_return_mismatch()
func (c *TypeCheckNoTranspileTest) TestGenerator_return_mismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Generator<number, string>} */ function *gen() { return 1; }").
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_returnIteratorLike_mismatch()
func (c *TypeCheckNoTranspileTest) TestGenerator_returnIteratorLike_mismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!IteratorLike<number, string>} */ function *gen() { return 1; }").
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_returnIterable_mismatch()
func (c *TypeCheckNoTranspileTest) TestGenerator_returnIterable_mismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Iterable<number, string>} */ function *gen() { return 1; }").
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_return_noSpecifiedReturnType()
func (c *TypeCheckNoTranspileTest) TestGenerator_return_noSpecifiedReturnType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Generator<string>} */ function *gen() {  return 1; }").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_return3()
func (c *TypeCheckNoTranspileTest) TestGenerator_return3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Generator<number>} */ function *gen() {  return; }").
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_yieldAll1()
func (c *TypeCheckNoTranspileTest) TestGenerator_yieldAll1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Generator<number>} */
function *gen() {
  yield* [1, 2, 3];
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_yieldAll2()
func (c *TypeCheckNoTranspileTest) TestGenerator_yieldAll2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Generator<number>} */ function *gen() { yield* 1; }").
		AddDiagnostic(`Expression yield* expects an iterable
found   : number
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_yieldAll3()
func (c *TypeCheckNoTranspileTest) TestGenerator_yieldAll3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Generator<number>} */
function *gen1() {
  yield 1;
}

/** @return {!Generator<number>} */
function *gen2() {
  yield* gen1();
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_yieldAll4()
func (c *TypeCheckNoTranspileTest) TestGenerator_yieldAll4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!Generator<string>} gen1
 * @return {!Generator<number>}
 */
function *gen2(gen1) {
  yield* gen1;
}
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : string
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_yieldAll4_returnMismatch()
func (c *TypeCheckNoTranspileTest) TestGenerator_yieldAll4_returnMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!Generator<number, string>} gen1
 * @return {!Generator<number>}
 */
function *gen2(gen1) {
  /** @type {number} */
  const str = yield* gen1;
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_yieldAll_string()
func (c *TypeCheckNoTranspileTest) TestGenerator_yieldAll_string() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Generator<string>} */
function *gen() {
  yield* 'some string';
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testGenerator_yieldAll_null()
func (c *TypeCheckNoTranspileTest) TestGenerator_yieldAll_null() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Generator<string>} */
function *gen() {
  yield* null;
}
`).
		AddDiagnostic(`Expression yield* expects an iterable
found   : null
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMemberFunctionDef1()
func (c *TypeCheckNoTranspileTest) TestMemberFunctionDef1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var obj = { // line break
  method (/** number */ n) {}
};
obj.method(1);
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMemberFunctionDef2()
func (c *TypeCheckNoTranspileTest) TestMemberFunctionDef2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var obj = { // line break
  method (/** string */ n) {}
};
obj.method(1);
`).
		AddDiagnostic(`actual parameter 1 of obj.method does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMemberFunctionDef3()
func (c *TypeCheckNoTranspileTest) TestMemberFunctionDef3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var obj = { method() {} }; new obj.method();").
		AddDiagnostic("cannot instantiate non-constructor, found type: function(): undefined").
		Run()
}

// Java: TypeCheckNoTranspileTest.testMemberFunctionDef_lends()
func (c *TypeCheckNoTranspileTest) TestMemberFunctionDef_lends() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function PolymerElement() {}
/** @param {...*} var_args */
PolymerElement.prototype.factoryImpl = function(var_args) {}
var Polymer = function(a) {};
`).
		AddSource(`/** @constructor @extends {PolymerElement} */
var X = function() {};
X = Polymer(/** @lends {X.prototype} */ {
 // Test that we can override PolymerElement.prototype.factoryImpl with a one-arg fn
  factoryImpl(e) {
    alert('Thank you for clicking');
  },
});
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMemberFunction_enum()
func (c *TypeCheckNoTranspileTest) TestMemberFunction_enum() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @enum */ var obj = {a() {}};").
		AddDiagnostic(`assignment to property a of enum{obj}
found   : function(): undefined
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp1()
func (c *TypeCheckNoTranspileTest) TestComputedProp1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var i = 1; var obj = { ['var' + i]: i, };").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp2a()
func (c *TypeCheckNoTranspileTest) TestComputedProp2a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var n;
var obj = {[n = 'foo']: i};
var /** number */ m = n;
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp2b()
func (c *TypeCheckNoTranspileTest) TestComputedProp2b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var /** number */ n = 1; // preserve new line
var obj = {
  [n = 'foo']: i
};
`).
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp2c()
func (c *TypeCheckNoTranspileTest) TestComputedProp2c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var n;
var obj = {[foo]: n = 'bar'};
var /** number */ m = n;
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp3()
func (c *TypeCheckNoTranspileTest) TestComputedProp3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var i = 1;
var obj = { ['var' + i]: i };
var x = obj.var1
`).
		AddDiagnostic("Property var1 never defined on obj").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp3b()
func (c *TypeCheckNoTranspileTest) TestComputedProp3b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var obj = { ['static']: 1 };
var /** number */ x = obj.static
`).
		AddDiagnostic("Property static never defined on obj").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp4()
func (c *TypeCheckNoTranspileTest) TestComputedProp4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesString(/** string */ str) {}

var obj = {
  /** @param {number} x */
  ['static']: (x) => {
    takesString(x);
  }
};
`).
		AddDiagnostic(`actual parameter 1 of takesString does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp_symbol()
func (c *TypeCheckNoTranspileTest) TestComputedProp_symbol() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var sym1 = Symbol('a'); var obj = {[sym1]: 1};").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp_number()
func (c *TypeCheckNoTranspileTest) TestComputedProp_number() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var obj = {[0]: 1};").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp_badKeyType()
func (c *TypeCheckNoTranspileTest) TestComputedProp_badKeyType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var foo = {}; var bar = {[foo]: 3};").
		AddDiagnostic(`property access
found   : {}
required: (string|symbol)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp_restrictedIndexType()
func (c *TypeCheckNoTranspileTest) TestComputedProp_restrictedIndexType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var /** !Object<string, *> */ obj = {[1]: 1};").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp_incorrectValueType1()
func (c *TypeCheckNoTranspileTest) TestComputedProp_incorrectValueType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var /** !Object<string, number> */ obj = {['x']: 'not numeric'};").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp_incorrectValueType2()
func (c *TypeCheckNoTranspileTest) TestComputedProp_incorrectValueType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = { /** @type {string} */ [1]: 12 };").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp_struct1()
func (c *TypeCheckNoTranspileTest) TestComputedProp_struct1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @struct */ var obj = {[1 + 2]: 3};").
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp_struct2()
func (c *TypeCheckNoTranspileTest) TestComputedProp_struct2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @struct */ var obj = {[Symbol.iterator]: function() {}};").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp_dict()
func (c *TypeCheckNoTranspileTest) TestComputedProp_dict() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ var obj = {[1]: 2};").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedProp_enum()
func (c *TypeCheckNoTranspileTest) TestComputedProp_enum() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @enum */ var obj = {[1]: 2};").
		AddDiagnostic("enum key must be a string or numeric literal").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedPropAllowedOnDictClass()
func (c *TypeCheckNoTranspileTest) TestComputedPropAllowedOnDictClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @dict */
class C {
  ['f']() {}
  [123]() {}
  [123n]() {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testNormalPropNotAllowedOnDictClass()
func (c *TypeCheckNoTranspileTest) TestNormalPropNotAllowedOnDictClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @dict */
class C {
  foo() {}
}
`).
		AddDiagnostic("Illegal key, the class is a dict").
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedPropNotAllowedOnStructClass()
func (c *TypeCheckNoTranspileTest) TestComputedPropNotAllowedOnStructClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C { // @struct is the default
  foo() {}
  ['f']() {}
}
`).
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckNoTranspileTest.testQuotedGetterPropNotAllowedOnStructClass()
func (c *TypeCheckNoTranspileTest) TestQuotedGetterPropNotAllowedOnStructClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C { // @struct is the default
  foo() {}
  get 'f'() {}
}
`).
		AddDiagnostic("Illegal key, the class is a struct").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTemplateLiteral1()
func (c *TypeCheckNoTranspileTest) TestTemplateLiteral1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var a, b\nvar /** string */ s = `template ${a} string ${b}`;\n").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTemplateLiteral2()
func (c *TypeCheckNoTranspileTest) TestTemplateLiteral2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var n; var s = `${n = 'str'}`; var /** number */ m = n;").
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTemplateLiteral3()
func (c *TypeCheckNoTranspileTest) TestTemplateLiteral3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var /** number */ n = 1; var s = `template ${n = 'str'} string`;").
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTemplateLiteral_substitutionsHaveAnyType()
func (c *TypeCheckNoTranspileTest) TestTemplateLiteral_substitutionsHaveAnyType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(/** * */ anyTypeParam) {\n  var /** string */ s = `template ${anyTypeParam} string`;\n}\n").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTemplateLiteral_isStringType()
func (c *TypeCheckNoTranspileTest) TestTemplateLiteral_isStringType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var /** number */ n = `${1}`;").
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_tagParameters1()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_tagParameters1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
 * @param {!ITemplateArray} template
 * @param {...*} var_args Substitution values.
 * @return {string}
 */
String.raw = function(template, var_args) {};
`).
		AddSource("String.raw`one ${1} two`").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_tagParameters2()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_tagParameters2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(/** !Array<string> */ strings){}\ntag`template string`;\n").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_tagParameters3()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_tagParameters3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(/** ?Array<string> */ strings){}\ntag`template string`;\n").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_tagParameters4()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_tagParameters4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(/** Object */ strings){}\ntag `template string`;\n").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_tagParameters5()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_tagParameters5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(/** ? */ strings){}\ntag `template string`;\n").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_invalidTagParameters1()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_invalidTagParameters1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(/** {a: number} */ strings){}\ntag `template string`;\n").
		AddDiagnostic(`Invalid type for the first parameter of tag function
found   : {a: number}
required: ITemplateArray
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_invalidTagParameters2()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_invalidTagParameters2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(/** !Array<number> */ strings) {}\ntag`template string`;\n").
		AddDiagnostic(`Invalid type for the first parameter of tag function
found   : Array<number>
required: ITemplateArray
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_invalidTagParameters3()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_invalidTagParameters3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(){} tag``;").
		AddDiagnostic("Function tag: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_tagNotAFunction()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_tagNotAFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const tag = 42; tag `template string`;").
		AddDiagnostic("number expressions are not callable").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_nullableTagFunction()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_nullableTagFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(/** ?function(!ITemplateArray) */ tag) {\n  tag `template string`;\n}\n").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_unknownTagFunction()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_unknownTagFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(/** ? */ tag) {\n  tag `template string`;\n}\n").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_tooFewArguments()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_tooFewArguments() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(strings, x, y) {} tag`${1}`;").
		AddDiagnostic("Function tag: called with 2 argument(s). Function requires at least 3 argument(s) and no more than 3 argument(s).").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_tooManyArguments()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_tooManyArguments() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(strings, x) {} tag`${0} ${1}`;").
		AddDiagnostic("Function tag: called with 3 argument(s). Function requires at least 2 argument(s) and no more than 2 argument(s).").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_argumentTypeMismatch()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_argumentTypeMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(strings, /** string */ s) {}\ntag`${123}`;\n").
		AddDiagnostic(`actual parameter 2 of tag does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_argumentWithCast()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_argumentWithCast() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(strings, /** string */ s) {}\ntag`${ /** @type {?} */ (123) }`;\n").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_optionalArguments()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_optionalArguments() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {number=} y */ function tag(strings, y){}\ntag``;\n").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_varArgs()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_varArgs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function tag(strings, /** ...number */ var_args){}\ntag`${1} ${'str'}`;\n").
		AddDiagnostic(`actual parameter 3 of tag does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_returnType1()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_returnType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function takesString(/** string */ s) {}\n\n/** @return {number} */\nfunction returnsNumber(strings){\n  return 1;\n}\ntakesString(returnsNumber`str`);\n").
		AddDiagnostic(`actual parameter 1 of takesString does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_returnType2()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_returnType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function takesString(/** string */ s) {}\n/**\n * @param {!ITemplateArray} strings\n * @param {T} subExpr\n * @param {*} var_args\n * @return {T}\n * @template T\n */\nfunction getFirstTemplateLitSub(strings, subExpr, var_args) { return subExpr; }\n\ntakesString(getFirstTemplateLitSub`${1}`);\n").
		AddDiagnostic(`actual parameter 1 of takesString does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTaggedTemplateLiteral_backInference()
func (c *TypeCheckNoTranspileTest) TestTaggedTemplateLiteral_backInference() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**\n* @param {T} x\n* @param {function(this:T, ...?)} z\n* @template T\n*/\nfunction f(x, z) {}\n// infers that \"this\" is ITemplateArray inside the function literal\nf`${ function() { /** @type {string} */ var x = this } }`;\n").
		AddDiagnostic(`initializing variable
found   : ITemplateArray
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testITemplateArray1()
func (c *TypeCheckNoTranspileTest) TestITemplateArray1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number */ n) {}
function f(/** !ITemplateArray */ arr) {
  for (let str of arr) {
    takesNumber(str);
  }
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : string
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationWithReturn()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationWithReturn() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var /** ?Foo */ cached = null;
class Foo {
  constructor() {
    if (cached) return cached;
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassErrorsReportedOnClassAndNotConstructor()
func (c *TypeCheckNoTranspileTest) TestClassErrorsReportedOnClassAndNotConstructor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @implements {number} */
class Foo {
  constructor() {
// Make sure there's an explicit constructor.
  }
}
`).
		AddDiagnostic("can only implement interfaces").
		Run()
}

// Java: TypeCheckNoTranspileTest.testInvalidInvocationOfClassConstructor()
func (c *TypeCheckNoTranspileTest) TestInvalidInvocationOfClassConstructor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor() {
  }
}
let /** ? */ x = Foo()
`).
		AddDiagnostic("Constructor (typeof Foo) should be called with the \"new\" keyword").
		Run()
}

// Java: TypeCheckNoTranspileTest.testInvalidInvocationOfClassConstructorWithReturnDeclaration()
func (c *TypeCheckNoTranspileTest) TestInvalidInvocationOfClassConstructorWithReturnDeclaration() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @return {!Array} */
  constructor() {
  }
}
let /** ? */ x = Foo()
`).
		AddDiagnostic("Constructor (typeof Foo) should be called with the \"new\" keyword").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationWithTemplate()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationWithTemplate() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class C {
  /** @param {T} a */
  constructor(a) {
  }
}
/** @type {null} */
const x = new C(0);
`).
		AddDiagnostic(`initializing variable
found   : C<number>
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclaration()
func (c *TypeCheckNoTranspileTest) TestClassDeclaration() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
var /** !Foo */ foo = new Foo();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationMismatch()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
class Bar {}
var /** !Foo */ foo = new Bar();
`).
		AddDiagnostic(`initializing variable
found   : Bar
required: Foo
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassGenerics()
func (c *TypeCheckNoTranspileTest) TestClassGenerics() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Foo {}
var /** !Foo<number> */ x = new Foo();
var /** !Foo<string> */ y = x;
`).
		AddDiagnostic(`initializing variable
found   : Foo<number>
required: Foo<string>
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassTooManyTypeParameters()
func (c *TypeCheckNoTranspileTest) TestClassTooManyTypeParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
var /** !Foo<number> */ x = new Foo();
`).
		AddDiagnosticType(RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassWithTemplatizedConstructorTooManyTypeParameters()
func (c *TypeCheckNoTranspileTest) TestClassWithTemplatizedConstructorTooManyTypeParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @template T */ constructor() {}
}
var /** !Foo<number> */ x = new Foo();
`).
		AddDiagnosticType(RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassWithTemplatizedClassAndConstructorTooManyTypeParameters()
func (c *TypeCheckNoTranspileTest) TestClassWithTemplatizedClassAndConstructorTooManyTypeParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Foo {
  /** @template U */ constructor() {}
}
var /** !Foo<number, number> */ x = new Foo();
`).
		AddDiagnosticType(RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationWithExtends()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationWithExtends() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
class Bar extends Foo {}
var /** !Foo */ foo = new Bar();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationWithExtendsMismatch()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationWithExtendsMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
class Bar extends Foo {}
var /** !Bar */ foo = new Foo();
`).
		AddDiagnostic(`initializing variable
found   : Foo
required: Bar
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationWithTransitiveExtends()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationWithTransitiveExtends() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
class Bar extends Foo {}
class Baz extends Bar {}
var /** !Foo */ foo = new Baz();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationWithAnonymousExtends()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationWithAnonymousExtends() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
class Bar extends class extends Foo {} {}
var /** !Foo */ foo = new Bar();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationInlineConstructorParameters()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationInlineConstructorParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor(/** number */ arg) {}
}
new Foo(42);
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationConstructorParametersMismatch()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationConstructorParametersMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor(/** number */ arg) {}
}
new Foo('xyz');
`).
		AddDiagnostic(`actual parameter 1 of Foo does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationTraditionalConstructorParameters()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationTraditionalConstructorParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @param {number} arg */
  constructor(arg) {}
}
new Foo(42);
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationTraditionalConstructorParametersMismatch()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationTraditionalConstructorParametersMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @param {number} arg */
  constructor(arg) {}
}
new Foo('xyz');
`).
		AddDiagnostic(`actual parameter 1 of Foo does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationInheritedConstructorParameters()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationInheritedConstructorParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor(/** number */ arg) {}
}
class Bar extends Foo {}
new Bar(42);
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationInheritedConstructorParametersMismatch()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationInheritedConstructorParametersMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor(/** number */ arg) {}
}
class Bar extends Foo {}
new Bar('xyz');
`).
		AddDiagnostic(`actual parameter 1 of Bar does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationWithSemicolonsBetweenMembers()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationWithSemicolonsBetweenMembers() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor() {};
  foo() {};
  bar() {};
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassPassedAsParameter()
func (c *TypeCheckNoTranspileTest) TestClassPassedAsParameter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
function foo(/** function(new: Foo) */ arg) {}
foo(class extends Foo {});
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassPassedAsParameterClassMismatch()
func (c *TypeCheckNoTranspileTest) TestClassPassedAsParameterClassMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
function foo(/** function(new: Foo) */ arg) {}
foo(class {});
`).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : (typeof <anonymous@testcode0:3>)
required: function(new:Foo): ?
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassPassedAsParameterConstructorParamsMismatch()
func (c *TypeCheckNoTranspileTest) TestClassPassedAsParameterConstructorParamsMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor(/** string */ arg) {}
}
function foo(/** function(new: Foo, number) */ arg) {}
foo(Foo);
`).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : (typeof Foo)
required: function(new:Foo, number): ?
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExpression()
func (c *TypeCheckNoTranspileTest) TestClassExpression() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var Foo = class Bar {}
var /** !Foo */ foo = new Foo();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExpressionDoesNotDefineTypeNameInOuterScope()
func (c *TypeCheckNoTranspileTest) TestClassExpressionDoesNotDefineTypeNameInOuterScope() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var Foo = class Bar {}
var /** !Bar */ foo = new Foo();
`).
		AddDiagnostic("Bad type annotation. Unknown type Bar").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExpressionDoesNotDefineConstructorReferenceInOuterScope()
func (c *TypeCheckNoTranspileTest) TestClassExpressionDoesNotDefineConstructorReferenceInOuterScope() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var Foo = class Bar {}
var /** null */ foo = new Bar();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExpressionAsStaticClassProeprty()
func (c *TypeCheckNoTranspileTest) TestClassExpressionAsStaticClassProeprty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
Foo.Bar = class extends Foo {}
var /** !Foo */ foo = new Foo.Bar();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSyntaxClassExtendsInterface()
func (c *TypeCheckNoTranspileTest) TestClassSyntaxClassExtendsInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Bar {}
class Foo extends Bar {}
`).
		AddDiagnostic("Foo cannot extend this type; constructors can only extend constructors").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSyntaxClassExtendsNonClass()
func (c *TypeCheckNoTranspileTest) TestClassSyntaxClassExtendsNonClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class Foo extends 42 {}").
		AddDiagnostic("Foo cannot extend this type; constructors can only extend constructors").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSyntaxInterfaceExtendsClass()
func (c *TypeCheckNoTranspileTest) TestClassSyntaxInterfaceExtendsClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Bar {}
/** @interface */
class Foo extends Bar {}
`).
		AddDiagnostic("Foo cannot extend this type; interfaces can only extend interfaces").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSyntaxInterfaceExtendsInterface()
func (c *TypeCheckNoTranspileTest) TestClassSyntaxInterfaceExtendsInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Bar {}
/** @interface */
class Foo extends Bar {}
var /** !Foo */ foo;
var /** !Bar */ bar = foo;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSyntaxInterfaceExtendsInterfaceMismatch()
func (c *TypeCheckNoTranspileTest) TestClassSyntaxInterfaceExtendsInterfaceMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Bar {}
/** @interface */
class Foo extends Bar {}
var /** !Bar */ bar;
var /** !Foo */ foo = bar;
`).
		AddDiagnostic(`initializing variable
found   : Bar
required: Foo
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSyntaxRecord()
func (c *TypeCheckNoTranspileTest) TestClassSyntaxRecord() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
class Rec {
  constructor() { /** @type {string} */ this.bar; }
  foo(/** number */ arg) {}
}
var /** !Rec */ rec = {bar: 'x', foo() {}};
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSyntaxRecordWithMethodMismatch()
func (c *TypeCheckNoTranspileTest) TestClassSyntaxRecordWithMethodMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
class Rec {
  foo(/** number */ arg) {}
}
var /** !Rec */ rec = {foo(/** string */ arg) {}};
`).
		AddDiagnostic(`initializing variable
found   : {foo: function(string): undefined}
required: Rec
missing : []
mismatch: [foo]
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSyntaxRecordWithPropertyMismatch()
func (c *TypeCheckNoTranspileTest) TestClassSyntaxRecordWithPropertyMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
class Rec {
  constructor() {
    /** @type {number} */
    this.foo;
  }
}
var /** !Rec */ rec = {foo: 'string'};
`).
		AddDiagnostic(`initializing variable
found   : {foo: string}
required: Rec
missing : []
mismatch: [foo]
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassJSDocExtendsInconsistentWithExtendsClause()
func (c *TypeCheckNoTranspileTest) TestClassJSDocExtendsInconsistentWithExtendsClause() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Bar {}
class Baz {}
/** @extends {Bar} */
class Foo extends Baz {}
`).
		AddDiagnostic(`mismatch in declaration of superclass type
found   : Baz
required: Bar
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassJSDocExtendsWithMissingExtendsClause()
func (c *TypeCheckNoTranspileTest) TestClassJSDocExtendsWithMissingExtendsClause() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Bar {}
/** @extends {Bar} */
class Foo {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testInterfaceHasBothExtendsClauseAndExtendsJSDoc()
func (c *TypeCheckNoTranspileTest) TestInterfaceHasBothExtendsClauseAndExtendsJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Bar {}
/** @interface */
class Baz {}
/**
 * @interface
 * @extends {Bar}
 */
class Foo extends Baz {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsGetElem()
func (c *TypeCheckNoTranspileTest) TestClassExtendsGetElem() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
/** @const {!Object<string, function(new:Foo)>} */
var obj = {};
class Bar extends obj['abc'] {}
var /** !Foo */ foo = new Bar();
`).
		AddDiagnostic("The right-hand side of an extends clause must be a qualified name, or else @extends must be specified in JSDoc").
		AddDiagnostic(`initializing variable
found   : Bar
required: Foo
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsFunctionCall()
func (c *TypeCheckNoTranspileTest) TestClassExtendsFunctionCall() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
/** @return {function(new:Foo)} */
function mixin() {}
class Bar extends mixin() {}
var /** !Foo */ foo = new Bar();
`).
		AddDiagnostic("The right-hand side of an extends clause must be a qualified name, or else @extends must be specified in JSDoc").
		AddDiagnostic(`initializing variable
found   : Bar
required: Foo
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInterfaceExtendsFunctionCall()
func (c *TypeCheckNoTranspileTest) TestClassInterfaceExtendsFunctionCall() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo {}
/** @return {function(new:Foo)} */
function mixin() {}
/** @interface */
class Bar extends mixin() {}
`).
		AddDiagnostic("The right-hand side of an extends clause must be a qualified name, or else @extends must be specified in JSDoc").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsFunctionCallWithJSDoc()
func (c *TypeCheckNoTranspileTest) TestClassExtendsFunctionCallWithJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor() { /** @type {number} */ this.foo; }
}
/** @return {function(new:Foo)} */
function mixin() {}
/** @extends {Foo} */
class Bar extends mixin() {}
var /** null */ x = new Bar().foo;
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsFunctionCallWithIncompatibleJSDoc()
func (c *TypeCheckNoTranspileTest) TestClassExtendsFunctionCallWithIncompatibleJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
class Baz {}
/** @return {function(new:Foo)} */
function mixin() {}
/** @extends {Baz} */
class Bar extends mixin() {}
`).
		AddDiagnostic(`mismatch in declaration of superclass type
found   : Foo
required: Baz
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassImplementsInterface()
func (c *TypeCheckNoTranspileTest) TestClassImplementsInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo { foo() {} }
/** @implements {Foo} */
class Bar {
  /** @override */
  foo() {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassImplementsInterfaceViaParent()
func (c *TypeCheckNoTranspileTest) TestClassImplementsInterfaceViaParent() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class IFoo { /** @return {*} */ foo() {} }
class Foo { /** @return {number} */ foo() {} }
/** @implements {IFoo} */
class Zoo extends Foo {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsAbstractClassesThatImplementsInterface()
func (c *TypeCheckNoTranspileTest) TestClassExtendsAbstractClassesThatImplementsInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class IFoo { foo() {} }
/** @abstract @implements {IFoo} */
class Foo { /** @override */ foo() {} }
/** @abstract @implements {IFoo} */
class Bar extends Foo {}
class Zoo extends Bar {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassMissingInterfaceMethod()
func (c *TypeCheckNoTranspileTest) TestClassMissingInterfaceMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo { foo() {} }
/** @implements {Foo} */
class Bar {}
`).
		AddDiagnostic("property foo on interface Foo is not implemented by type Bar").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassCannotImplementInterfaceWithAPrototypeAssignment()
func (c *TypeCheckNoTranspileTest) TestClassCannotImplementInterfaceWithAPrototypeAssignment() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
function MyInterface() {}
/** @type {string} */
MyInterface.prototype.foo;
/** @constructor @implements {MyInterface} */
function MyClass() {}
MyClass.prototype = MyInterface.prototype;
`).
		AddDiagnostic("property foo on interface MyInterface is not implemented by type MyClass").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassAbstractClassNeedNotExplicitlyOverrideUnimplementedInterfaceMethods()
func (c *TypeCheckNoTranspileTest) TestClassAbstractClassNeedNotExplicitlyOverrideUnimplementedInterfaceMethods() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo { foo() {} }
/** @abstract @implements {Foo} */
class Bar {
// Also make sure that we can call the interface method that is not re-declared within
// the abstract class itself.
    bar() { this.foo(); }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassMissingOverrideAnnotationForInterfaceMethod()
func (c *TypeCheckNoTranspileTest) TestClassMissingOverrideAnnotationForInterfaceMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo { foo() {} }
/** @implements {Foo} */
class Bar {
  foo() {}
}
`).
		AddDiagnostic("property foo already defined on interface Foo; use @override to override it").
		Run()
}

// Java: TypeCheckNoTranspileTest.testAbstractClassMissingOverrideAnnotationForInterfaceMethod()
func (c *TypeCheckNoTranspileTest) TestAbstractClassMissingOverrideAnnotationForInterfaceMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo { foo() {} }
/** @abstract @implements {Foo} */
class Bar {
  foo() {}
}
`).
		AddDiagnostic("property foo already defined on interface Foo; use @override to override it").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassMissingOverrideAnnotationForInterfaceInstanceProperty()
func (c *TypeCheckNoTranspileTest) TestClassMissingOverrideAnnotationForInterfaceInstanceProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @record */ // `@interface` would also trigger this.\nclass Foo {\n  constructor() {\n    /** @type {number} */\n    this.bar;\n  }\n}\n\n/** @implements {Foo} */\nclass MyFoo { }\n// No `@override`.\n// For some reason we only check this when assigning to prototype properties, not to\n// instance properties.\n/** @type {number} */\nMyFoo.prototype.bar = 0;\n").
		AddDiagnostic("property bar already defined on interface Foo; use @override to override it").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassIncompatibleInterfaceMethodImplementation()
func (c *TypeCheckNoTranspileTest) TestClassIncompatibleInterfaceMethodImplementation() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo {
  /** @return {number} */ foo() {}
}
/** @implements {Foo} */
class Bar {
  /** @override @return {number|string} */
  foo() {}
}
`).
		AddDiagnostic(`mismatch of the foo property on type Bar and the type of the property it overrides from interface Foo
original: function(this:Foo): number
override: function(this:Bar): (number|string)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassIncompatibleInterfaceMethodImplementationInheritedOverAbstractClass()
func (c *TypeCheckNoTranspileTest) TestClassIncompatibleInterfaceMethodImplementationInheritedOverAbstractClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo {
  /** @return {number} */ foo() {}
}
/** @abstract @implements {Foo} */
class Bar {}
class Zoo extends Bar {
  /** @override @return {number|string} */
  foo() {}
}
`).
		AddDiagnostic(`mismatch of the foo property on type Zoo and the type of the property it overrides from interface Foo
original: function(this:Foo): number
override: function(this:Zoo): (number|string)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAbstractClassIncompatibleInterfaceMethodImplementation()
func (c *TypeCheckNoTranspileTest) TestAbstractClassIncompatibleInterfaceMethodImplementation() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo {
  /** @return {number} */ foo() {}
}
/** @abstract @implements {Foo} */
class Bar {
  /** @override @return {number|string} */
  foo() {}
}
`).
		AddDiagnostic(`mismatch of the foo property on type Bar and the type of the property it overrides from interface Foo
original: function(this:Foo): number
override: function(this:Bar): (number|string)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassMissingTransitiveInterfaceMethod()
func (c *TypeCheckNoTranspileTest) TestClassMissingTransitiveInterfaceMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo { foo() {} }
/** @interface @extends {Foo} */
class Bar {}
/** @implements {Bar} */
class Baz {}
`).
		AddDiagnostic("property foo on interface Foo is not implemented by type Baz").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassMissingInterfaceInstanceProperty()
func (c *TypeCheckNoTranspileTest) TestClassMissingInterfaceInstanceProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @record */ // `@interface` would also trigger this.\nclass Foo {\n  constructor() {\n    /** @type {number} */\n    this.bar;\n  }\n}\n\n/** @implements {Foo} */\nclass MyFoo { }\n").
		AddDiagnostic("property bar on interface Foo is not implemented by type MyFoo").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassMissingTransitiveInterfaceInstanceProperty()
func (c *TypeCheckNoTranspileTest) TestClassMissingTransitiveInterfaceInstanceProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @record */ // `@interface` would also trigger this.\nclass Foo {\n  constructor() {\n    /** @type {number} */\n    this.bar;\n  }\n}\n\n/** @record */\nclass SubFoo extends Foo { }\n\n/** @implements {SubFoo} */\nclass MyFoo { }\n").
		AddDiagnostic("property bar on interface Foo is not implemented by type MyFoo").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInvalidOverrideOfInterfaceInstanceProperty()
func (c *TypeCheckNoTranspileTest) TestClassInvalidOverrideOfInterfaceInstanceProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @record */ // `@interface` would also trigger this.\nclass Foo {\n  constructor() {\n    /** @type {number} */\n    this.bar;\n  }\n}\n\n/** @implements {Foo} */\nclass MyFoo {\n  constructor() {\n    /** @type {string} */\n    this.bar;\n  }\n}\n").
		AddDiagnostic(`mismatch of the bar property on type MyFoo and the type of the property it overrides from interface Foo
original: number
override: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassPrototypeOverrideOfInterfaceInstanceProperty()
func (c *TypeCheckNoTranspileTest) TestClassPrototypeOverrideOfInterfaceInstanceProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @record */ // `@interface` would also trigger this.\nclass Foo {\n  constructor() {\n    /** @type {number} */\n    this.bar;\n  }\n}\n\n/** @implements {Foo} */\nclass MyFoo { }\n// It's legal to fulfill the interface using either instance or prototype properties.\n/** @override */\nMyFoo.prototype.bar;\n").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInheritedInterfaceMethod()
func (c *TypeCheckNoTranspileTest) TestClassInheritedInterfaceMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo { foo() {} bar() {} }
/** @abstract */
class Bar { foo() {} }
/** @implements {Foo} */
class Baz extends Bar { /** @override */ bar() {} }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassMixinAllowsNonOverriddenInterfaceMethods()
func (c *TypeCheckNoTranspileTest) TestClassMixinAllowsNonOverriddenInterfaceMethods() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo {
  /** @return {number} */ foo() {}
}
class Bar {}
// TODO(sdh): Intersection types would allow annotating this correctly.
/** @return {function(new:Bar)} */
function mixin() {}
/** @extends {Bar} @implements {Foo} */
class Baz extends mixin() {}
`).
		AddDiagnostic("property foo on interface Foo is not implemented by type Baz").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassDeclarationWithExtendsOnlyInJSDoc()
func (c *TypeCheckNoTranspileTest) TestClassDeclarationWithExtendsOnlyInJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
/** @extends {Foo} */
class Bar {}
var /** !Foo */ foo = new Bar();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassConstructorTypeParametersNotIncludedOnClass()
func (c *TypeCheckNoTranspileTest) TestClassConstructorTypeParametersNotIncludedOnClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Foo {
  /** @template U */
  constructor() {}
}
var /** !Foo<string, string> */ x = new Foo();
`).
		AddDiagnosticType(RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassConstructorTypeParametersNotVisibleFromOtherMethodSignatures()
func (c *TypeCheckNoTranspileTest) TestClassConstructorTypeParametersNotVisibleFromOtherMethodSignatures() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @template T */
  constructor() {}
  /** @param {T} p */
  foo(p) {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassConstructorTypeParametersNotVisibleFromOtherMethodBodies()
func (c *TypeCheckNoTranspileTest) TestClassConstructorTypeParametersNotVisibleFromOtherMethodBodies() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @template T */
  constructor() {}
  foo(p) {
    var /** T */ x;
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassTtlNotAllowedOnClass()
func (c *TypeCheckNoTranspileTest) TestClassTtlNotAllowedOnClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @template T := 'number' =: */ class Foo {}").
		AddDiagnostic("Template type transformation T not allowed on classes or interfaces").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassTtlAllowedOnConstructor()
func (c *TypeCheckNoTranspileTest) TestClassTtlAllowedOnConstructor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /**
   * @param {T} arg
   * @template T := 'number' =:
   */
  constructor(arg) {
    var /** T */ x = arg;
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassTtlAllowedOnMethod()
func (c *TypeCheckNoTranspileTest) TestClassTtlAllowedOnMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @template T := 'number' =: */
  foo(/** T */ arg) {
    var /** T */ x = arg;
  }
}
new Foo().foo('x')
`).
		AddDiagnostic(`actual parameter 1 of Foo.prototype.foo does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassConstructorTypeParametersChecked()
func (c *TypeCheckNoTranspileTest) TestClassConstructorTypeParametersChecked() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Foo {
  /** @template U */
  constructor(/** U */ arg1, /** function(U): T */ arg2) {}
}
/** @param {string} arg
    @return {number} */
function f(arg) {}
var /** !Foo<number> */ foo = new Foo('x', f);
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassConstructorTypeParametersWithClassTypeMismatch()
func (c *TypeCheckNoTranspileTest) TestClassConstructorTypeParametersWithClassTypeMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Foo {
  /** @template U */
  constructor(/** U */ arg1, /** function(U): T */ arg2) {}
}
/** @param {string} arg
    @return {number} */
function f(arg) {}
var /** !Foo<string> */ foo = new Foo('x', f);
`).
		AddDiagnostic(`initializing variable
found   : Foo<number>
required: Foo<string>
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassConstructorTypeParametersWithParameterTypeMismatch()
func (c *TypeCheckNoTranspileTest) TestClassConstructorTypeParametersWithParameterTypeMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Foo {
  /** @template U */
  constructor(/** U */ arg1, /** function(U): T */ arg2) {}
}
/** @param {string} arg
    @return {number} */
function f(arg) {}
var foo = new Foo(42, f);
`).
		AddDiagnostic(`actual parameter 2 of Foo does not match formal parameter
found   : function(string): number
required: function((number|string)): number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSideInheritanceFillsInParameterTypesWhenCheckingBody()
func (c *TypeCheckNoTranspileTest) TestClassSideInheritanceFillsInParameterTypesWhenCheckingBody() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  static foo(/** string */ arg) {}
}
class Bar extends Foo {
  /** @override */
  static foo(arg) {
    var /** null */ x = arg;
  }
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassMethodParameters()
func (c *TypeCheckNoTranspileTest) TestClassMethodParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @param {number} arg */
  m(arg) {}
}
new C().m('x');
`).
		AddDiagnostic(`actual parameter 1 of C.prototype.m does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInheritedMethodParameters()
func (c *TypeCheckNoTranspileTest) TestClassInheritedMethodParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var B = class {
  /** @param {boolean} arg */
  method(arg) {}
};
var C = class extends B {};
new C().method(1);
`).
		AddDiagnostic(`actual parameter 1 of B.prototype.method does not match formal parameter
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassMethodReturns()
func (c *TypeCheckNoTranspileTest) TestClassMethodReturns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var D = class {
  /** @return {number} */
  m() {}
}
var /** null */ x = new D().m();
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInheritedMethodReturns()
func (c *TypeCheckNoTranspileTest) TestClassInheritedMethodReturns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Q {
  /** @return {string} */
  method() {}
};
var P = class extends Q {};
var /** null */ x = new P().method();
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethodParameters()
func (c *TypeCheckNoTranspileTest) TestStaticMethodParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @param {number} arg */
  static m(arg) {}
}
C.m('x');
`).
		AddDiagnostic(`actual parameter 1 of C.m does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInheritedStaticMethodParameters()
func (c *TypeCheckNoTranspileTest) TestClassInheritedStaticMethodParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var B = class {
  /** @param {boolean} arg */
  static method(arg) {}
};
var C = class extends B {};
C.method(1);
`).
		AddDiagnostic(`actual parameter 1 of C.method does not match formal parameter
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethodReturns()
func (c *TypeCheckNoTranspileTest) TestStaticMethodReturns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var D = class {
  /** @return {number} */
  static m() {}
};
var /** null */ x = D.m();
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInheritedStaticMethodReturns()
func (c *TypeCheckNoTranspileTest) TestClassInheritedStaticMethodReturns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Q {
  /** @return {string} */
  static method() {}
}
class P extends Q {}
var /** null */ x = P.method();
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethodCalledOnInstance()
func (c *TypeCheckNoTranspileTest) TestStaticMethodCalledOnInstance() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  static m() {}
}
new C().m();
`).
		AddDiagnostic("Property m never defined on C").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInstanceMethodCalledOnClass()
func (c *TypeCheckNoTranspileTest) TestClassInstanceMethodCalledOnClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  m() {}
}
C.m();
`).
		AddDiagnostic("Property m never defined on C").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInstanceMethodOverriddenWithMissingOverrideAnnotation()
func (c *TypeCheckNoTranspileTest) TestClassInstanceMethodOverriddenWithMissingOverrideAnnotation() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  /** @param {string|number} arg */
  method(arg) {}
}
class Sub extends Base {
  method(arg) {}
}
`).
		AddDiagnostic("property method already defined on superclass Base; use @override to override it").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInstanceMethodOverriddenWithWidenedType()
func (c *TypeCheckNoTranspileTest) TestClassInstanceMethodOverriddenWithWidenedType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  /** @param {string} arg */
  method(arg) {}
}
class Sub extends Base {
  /** @override @param {string|number} arg */
  method(arg) {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInstanceMethodOverriddenWithIncompatibleType()
func (c *TypeCheckNoTranspileTest) TestClassInstanceMethodOverriddenWithIncompatibleType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  /** @param {string|number} arg */
  method(arg) {}
}
class Sub extends Base {
  /** @override @param {string} arg */
  method(arg) {}
}
`).
		AddDiagnostic(`mismatch of the method property type and the type of the property it overrides from superclass Base
original: function(this:Base, (number|string)): undefined
override: function(this:Sub, string): undefined
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassInstanceMethodOverriddenWithIncompatibleType2()
func (c *TypeCheckNoTranspileTest) TestClassInstanceMethodOverriddenWithIncompatibleType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  /** @return {string} */
  method() {}
}
class Sub extends Base {
  /** @override @return {string|number} */
  method() {}
}
`).
		AddDiagnostic(`mismatch of the method property type and the type of the property it overrides from superclass Base
original: function(this:Base): string
override: function(this:Sub): (number|string)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_overriddenInBody_withSubtype_atOverride_isOk()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_overriddenInBody_withSubtype_atOverride_isOk() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  /** @param {string} arg */
  static method(arg) {}
}

class Sub extends Base {
  /**
   * @override
// Method is a subtype due to parameter contravariance.
   * @param {string|number} arg
   */
  static method(arg) {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_overriddenInBody_notAtOverride_isBad()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_overriddenInBody_notAtOverride_isBad() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  /** @param {string} arg */
  static method(arg) {}
}

class Sub extends Base {
// Method is a subtype due to parameter contravariance.
  /** @param {string|number} arg */
  static method(arg) {}
}
`).
		AddDiagnostic("property method already defined on supertype (typeof Base); use @override to override it").
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_thatIsNotAnOverride_atOverride_isBad()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_thatIsNotAnOverride_atOverride_isBad() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  /**
   * @override
   * @param {string} arg
   */
  static method(arg) {}
}
`).
		AddDiagnostic("property method not defined on any supertype of (typeof Base)").
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_overriddenInBody_withSupertype_isBad()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_overriddenInBody_withSupertype_isBad() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  /** @param {string|number} arg */
  static method(arg) {}
}
class Sub extends Base {
  /**
   * @override
// Method is a supertype due to parameter contravariance.
   * @param {string} arg
   */
  static method(arg) {}
}
`).
		AddDiagnostic(`mismatch of the method property type and the type of the property it overrides from supertype (typeof Base)
original: function(this:(typeof Base), (number|string)): undefined
override: function(this:(typeof Sub), string): undefined
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_overriddenInBody_withSupertype_fromInline_isBad()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_overriddenInBody_withSupertype_fromInline_isBad() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  static method(/** string|number */ arg) {}
}
class Sub extends Base {
  /** @override */
  static method(/** string */ arg) {}
}
`).
		AddDiagnostic(`mismatch of the method property type and the type of the property it overrides from supertype (typeof Base)
original: function(this:(typeof Base), (number|string)): undefined
override: function(this:(typeof Sub), string): undefined
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_overriddenOutsideBody_withSubtype_atOverride_isOk()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_overriddenOutsideBody_withSubtype_atOverride_isOk() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  static method(/** string */ arg) {}
}

class Sub extends Base { }

/**
 * @override
// Method is a subtype due to parameter contravariance.
 * @param {string|number} arg
 */
Sub.method = function(arg) {};
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_overriddenOutsideBody_withSupertype_isBad()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_overriddenOutsideBody_withSupertype_isBad() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {
  static method(/** string|number */ arg) {}
}

class Sub extends Base { }

/**
 * @override
 * @param {string} arg
 */
Sub.method = function(arg) {};
`).
		AddDiagnostic(`mismatch of the method property type and the type of the property it overrides from supertype (typeof Base)
original: function(this:(typeof Base), (number|string)): undefined
override: function(string): undefined
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_onInterface_overriddenInBody_withSubtype_atOverride_isOk()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_onInterface_overriddenInBody_withSubtype_atOverride_isOk() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Base {
  /** @param {string} arg */
  static method(arg) {}
}

/** @interface */
class Sub extends Base {
  /**
   * @override
// Method is a subtype due to parameter contravariance.
   * @param {string|number} arg
   */
  static method(arg) {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_onNamespacedType_overriddenOutsideBody_withSubtype_atOverride_isOk()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_onNamespacedType_overriddenOutsideBody_withSubtype_atOverride_isOk() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const ns = {};

ns.Base = class {
  /** @param {string} arg */
  static method(arg) {}
};

ns.Sub = class extends ns.Base {};

/**
 * @override
// Method is a subtype due to parameter contravariance.
 * @param {string|number} arg
 */
// We specifically want to check that q-name lookups are checked.
ns.Sub.method = function(arg) {};
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_onNamespacedType_overriddenOutsideBody_notAtOverride_isBad()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_onNamespacedType_overriddenOutsideBody_notAtOverride_isBad() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const ns = {};

ns.Base = class {
  /** @param {string} arg */
  static method(arg) {}
};

ns.Sub = class extends ns.Base {};

// Method is a subtype due to parameter contravariance.
/** @param {string|number} arg */
// We specifically want to check that q-name lookups are checked.
ns.Sub.method = function(arg) {};
`).
		AddDiagnostic("property method already defined on supertype (typeof ns.Base); use @override to override it").
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticMethod_onNamespacedType_overridden_withNonSubtype_isBad()
func (c *TypeCheckNoTranspileTest) TestStaticMethod_onNamespacedType_overridden_withNonSubtype_isBad() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const ns = {};

ns.Base = class {
  /** @param {string} arg */
  static method(arg) {}
};

ns.Sub = class extends ns.Base {};

/**
 * @override
// Method is a subtype due to parameter contravariance.
 * @param {number} arg
 */
// We specifically want to check that q-name lookups are checked.
ns.Sub.method = function(arg) {};
`).
		AddDiagnostic(`mismatch of the method property type and the type of the property it overrides from supertype (typeof ns.Base)
original: function(this:(typeof ns.Base), string): undefined
override: function(number): undefined
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsForwardReference_staticMethodThatIsAnOverride_atOverride_isOk()
func (c *TypeCheckNoTranspileTest) TestClassExtendsForwardReference_staticMethodThatIsAnOverride_atOverride_isOk() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {function(new: Parent): ?} */
function mixin() {}
/** @extends {Parent} */
class Middle extends mixin() {
  /** @override */
  static method() {}
}

class Child extends Middle {
  /** @override */
  static method() {}
}

class Parent {
  method() {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassTreatedAsStruct()
func (c *TypeCheckNoTranspileTest) TestClassTreatedAsStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
var foo = new Foo();
foo.x = 42;
`).
		AddDiagnostic("Cannot add a property to a struct instance after it is constructed. (If you already declared the property, make sure to give it a type.)").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassTreatedAsStructSymbolAccess()
func (c *TypeCheckNoTranspileTest) TestClassTreatedAsStructSymbolAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
var foo = new Foo();
foo[Symbol.iterator] = 42;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassAnnotatedWithUnrestricted()
func (c *TypeCheckNoTranspileTest) TestClassAnnotatedWithUnrestricted() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @unrestricted */ class Foo {}
var foo = new Foo();
foo.x = 42;
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassAnnotatedWithDictDotAccess()
func (c *TypeCheckNoTranspileTest) TestClassAnnotatedWithDictDotAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @dict */ class Foo {}
var foo = new Foo();
foo.x = 42;
`).
		AddDiagnostic("Cannot do '.' access on a dict").
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassAnnotatedWithDictComputedAccess()
func (c *TypeCheckNoTranspileTest) TestClassAnnotatedWithDictComputedAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @dict */ class Foo {}
var foo = new Foo();
foo['x'] = 42;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperInConstructor()
func (c *TypeCheckNoTranspileTest) TestClassSuperInConstructor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor(/** number */ arg) {}
}
class Bar extends Foo {
  constructor(/** string */ arg) { super(1); }
}
var /** !Foo */ foo = new Bar('x');
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperConstructorParameterMismatch()
func (c *TypeCheckNoTranspileTest) TestClassSuperConstructorParameterMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor(/** number */ arg) {}
}
class Bar extends Foo {
  constructor() {
    super('x');
  }
}
`).
		AddDiagnostic(`actual parameter 1 of super does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperConstructorParameterCountMismatch()
func (c *TypeCheckNoTranspileTest) TestClassSuperConstructorParameterCountMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
class Bar extends Foo {
  constructor() {
    super(1);
  }
}
`).
		AddDiagnostic("Function super: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperMethodNotPresent()
func (c *TypeCheckNoTranspileTest) TestClassSuperMethodNotPresent() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
class Bar extends Foo {
  foo() {
    super.foo();
  }
}
`).
		AddDiagnostic("Property foo never defined on Foo").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperMethodParameterMismatch()
func (c *TypeCheckNoTranspileTest) TestClassSuperMethodParameterMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @param {string} arg */
  foo(arg) {}
}
class Bar extends Foo {
  /** @override */
  foo() {
    super.foo(42);
  }
}
`).
		AddDiagnostic(`actual parameter 1 of Foo.prototype.foo does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperMethodCalledFromArrow()
func (c *TypeCheckNoTranspileTest) TestClassSuperMethodCalledFromArrow() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @param {string} arg */
  foo(arg) {}
}
class Bar extends Foo {
  /** @override */
  foo() {
    () => super.foo(42);
  }
}
`).
		AddDiagnostic(`actual parameter 1 of Foo.prototype.foo does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperMethodReturnType()
func (c *TypeCheckNoTranspileTest) TestClassSuperMethodReturnType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @return {string} */
  foo() {}
}
class Bar extends Foo {
  /** @override */
  foo() {
    var /** null */ x = super.foo();
  }
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperMethodFromDifferentMethod()
func (c *TypeCheckNoTranspileTest) TestClassSuperMethodFromDifferentMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @return {string} */
  foo() {}
}
class Bar extends Foo {
  bar() {
    var /** null */ x = super.foo();
  }
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperMethodNotWidenedWhenOverrideWidens()
func (c *TypeCheckNoTranspileTest) TestClassSuperMethodNotWidenedWhenOverrideWidens() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @param {string} arg */
  foo(arg) {}
}
class Bar extends Foo {
  /** @override @param {string|number} arg */
  foo(arg) {}
  bar() {
    super.foo(42);
  }
}
`).
		AddDiagnostic(`actual parameter 1 of Foo.prototype.foo does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperMethodCallableInParameters()
func (c *TypeCheckNoTranspileTest) TestClassSuperMethodCallableInParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class Foo {\n  /** @return {string|number} */\n  foo() {\n    return 0;\n  }\n}\nclass Bar extends Foo {\n  /**\n   * @param {number=} param\n   * @return {string}\n   * @override\n   */\n//  super.foo() returns string|number, so `param` is typed as `string|number`\n  foo(param = super.foo()) {\n    return 'param: ' + param;\n  }\n}\n").
		AddDiagnostic(`default value has wrong type
found   : (number|string)
required: (number|undefined)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAbstractSuperMethodCall_warning()
func (c *TypeCheckNoTranspileTest) TestAbstractSuperMethodCall_warning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Foo {
  /** @abstract */
  foo() {}
}
class Bar extends Foo {
  /** @override */
  foo() {
    super.foo();
  }
}
`).
		AddDiagnostic("Abstract super method Foo.prototype.foo cannot be dereferenced").
		Run()
}

// Java: TypeCheckNoTranspileTest.testAbstractInheritedSuperMethodCall_warning()
func (c *TypeCheckNoTranspileTest) TestAbstractInheritedSuperMethodCall_warning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Foo {
  /** @abstract */
  foo() {}
}
/** @abstract */
class Bar extends Foo {}
class Baz extends Bar {
  /** @override */
  foo() {
    super.foo();
  }
}
`).
		AddDiagnostic("Abstract super method Foo.prototype.foo cannot be dereferenced").
		Run()
}

// Java: TypeCheckNoTranspileTest.testAbstractInheritedSuperMethodCallInAbstractClass_warning()
func (c *TypeCheckNoTranspileTest) TestAbstractInheritedSuperMethodCallInAbstractClass_warning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Base {
  /** @abstract */
  foo() {}
}
/** @abstract */
class Sub extends Base {
  bar() {
    super.foo();
  }
}
`).
		AddDiagnostic("Abstract super method Base.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testConcreteSuperMethodCall_noWarning()
func (c *TypeCheckNoTranspileTest) TestConcreteSuperMethodCall_noWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Foo {
  foo() {}
}
class Bar extends Foo {
  /** @override */
  foo() {
    super.foo();
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticSuperParameterMismatch()
func (c *TypeCheckNoTranspileTest) TestClassStaticSuperParameterMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  static foo(/** number */ arg) {}
}
class Bar extends Foo {
  /** @override */
  static foo() {
    super.foo('x');
  }
}
`).
		AddDiagnostic(`actual parameter 1 of super.foo does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticSuperCalledFromArrow()
func (c *TypeCheckNoTranspileTest) TestClassStaticSuperCalledFromArrow() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  static foo(/** number */ arg) {}
}
class Bar extends Foo {
  /** @override */
  static foo() {
    () => super.foo('x');
  }
}
`).
		AddDiagnostic(`actual parameter 1 of super.foo does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticSuperParameterCountMismatch()
func (c *TypeCheckNoTranspileTest) TestClassStaticSuperParameterCountMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  static foo() {}
}
class Bar extends Foo {
  /** @override */
  static foo() {
    super.foo(1);
  }
}
`).
		AddDiagnostic("Function super.foo: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticSuperNotPresent()
func (c *TypeCheckNoTranspileTest) TestClassStaticSuperNotPresent() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
class Bar extends Foo {
  static foo() {
    super.foo;
  }
}
`).
		AddDiagnostic("Property foo never defined on super").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticSuperCallsDifferentMethod()
func (c *TypeCheckNoTranspileTest) TestClassStaticSuperCallsDifferentMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @param {string} arg */
  static foo(arg) {}
}
class Bar extends Foo {
  /** @override */
  static foo(/** string|number */ arg) {}
  static bar() { super.foo(42); }
}
`).
		AddDiagnostic(`actual parameter 1 of super.foo does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassFields()
func (c *TypeCheckNoTranspileTest) TestClassFields() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class A {
  x = 2;
}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class B {
  x;
}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  x
}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class D {
  /** @type {string|undefined} */
  x;
}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class E {
  /** @type {string} @suppress {checkTypes} */
  x = 2;
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassFieldsThis()
func (c *TypeCheckNoTranspileTest) TestClassFieldsThis() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class F {
  /** @type {number} */
  x = 2;
  /** @type {boolean} */
  y = this.x
}
`).
		AddDiagnostic(`assignment to property y of F
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassFieldsSuper()
func (c *TypeCheckNoTranspileTest) TestClassFieldsSuper() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class G {
  /** @return {number} */
  getX() { return 2; }
}
class H extends G {
  /** @return {?} */
  /** @override*/ getX() {}
  /** @type {string} */
  y = super.getX();
}
`).
		AddDiagnostic(`assignment to property y of H
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedFields()
func (c *TypeCheckNoTranspileTest) TestComputedFields() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var /** number */ x = 1;
function takesNumber(/** number */ x) {}
/** @unrestricted */
class Foo {
  /** @type {boolean} */
  x = true;
  [this.x] = takesNumber(this.x);
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticClassFields()
func (c *TypeCheckNoTranspileTest) TestStaticClassFields() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class A {
  static x = 2;
}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class B {
  static x;
}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  static x
}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class D {
  /** @type {string|undefined} */
  static x;
}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class E {
  /** @type {string} @suppress {checkTypes} */
  static x = 2;
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticClassFieldsThis()
func (c *TypeCheckNoTranspileTest) TestStaticClassFieldsThis() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class F {
  /** @type {number} */
  static x = 2;
  /** @type {boolean} */
  static y = this.x
}
`).
		AddDiagnostic(`assignment to property y of F
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticComputedFields()
func (c *TypeCheckNoTranspileTest) TestStaticComputedFields() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var /** number */ x = 1;
function takesNumber(/** number */ x) {}
/** @unrestricted */
class Foo {
  /** @type {boolean} */
  static x = true;
  static [this.x] = takesNumber(this.x);
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStaticClassFieldsSuper()
func (c *TypeCheckNoTranspileTest) TestStaticClassFieldsSuper() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class G {
  /** @type {number} */
  static x = 2;
}
class H extends G {
  /** @type {string} */
  static y = super.x;
}
`).
		AddDiagnostic(`assignment to property y of H
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockVariablesWrongTypes()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockVariablesWrongTypes() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  static {
    /** @type {number} */
    let str = 'str';
    /** @type {boolean|string} */
    const num = 5;
    /** @type {?string} */
    var bool = true;
  }
};
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		AddDiagnostic(`initializing variable
found   : number
required: (boolean|string)
`).
		AddDiagnostic(`initializing variable
found   : boolean
required: (null|string)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockPropertyWithThis()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockPropertyWithThis() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number} */
var x = 4;
class Foo {
  static {
    this.x;
  }
};
`).
		AddDiagnostic("Property x never defined on this").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockWithWrongTypeThisRHS()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockWithWrongTypeThisRHS() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  static {
    /** @type {number} */
    this.num = 1;
    /** @type {string} */
    var str = this.num;
  }
};
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockWithWrongTypeThisLHS()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockWithWrongTypeThisLHS() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  static {
    /** @type {string} */
    this.str = 2;
  }
};
`).
		AddDiagnostic(`assignment to property str of this
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockWithSuper()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockWithSuper() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @type {string} */
  static str;
}
class Bar extends Foo {
  static {
    super.str = 'str';
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockWithWrongTypeSuper()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockWithWrongTypeSuper() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @type {string} */
  static str;
}
class Bar extends Foo {
  static {
    super.str = 5;
  }
};
`).
		AddDiagnostic(`assignment to property str of super
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockInheritanceWithClassName()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockInheritanceWithClassName() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  static foo(/** number */ arg) {}
}
class Bar extends Foo {
  static {
    Foo.foo('str');
  }
};
`).
		AddDiagnostic(`actual parameter 1 of Foo.foo does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockWithWrongTypeSuperParameter()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockWithWrongTypeSuperParameter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  static foo(/** number */ arg) {}
}
class Bar extends Foo {
  static {
    super.foo('str');
  }
};
`).
		AddDiagnostic(`actual parameter 1 of super.foo does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockTypeNarrowing()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockTypeNarrowing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  static {
    /** @param {?string} x */
    function foo(x) {
      if (x != null) {
        /** @type {string} */
        const noNull = x;
      }
    }
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockTypeNarrowing2()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockTypeNarrowing2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("/** @type {?string} */ var strOrNull;").
		AddSource(`class C {
  static {
    if (strOrNull != null) {
      /** @type {string} */
      const noNull = strOrNull;
    }
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticBlockOverrideSupertypeOnAnonymousClass()
func (c *TypeCheckNoTranspileTest) TestClassStaticBlockOverrideSupertypeOnAnonymousClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function use(ctor) {}

class Foo {
  static {
    /** @type {string} */
    this.str;
  }
}
use(class extends Foo {
  static { this.str = 3; }
});
`).
		AddDiagnostic(`assignment to property str of this
found   : number
required: string
`).
		AddDiagnostic("property str already defined on supertype (typeof Foo); use @override to override it").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassTypeOfThisInConstructor()
func (c *TypeCheckNoTranspileTest) TestClassTypeOfThisInConstructor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  constructor() {
    var /** null */ foo = this;
  }
}
`).
		AddDiagnostic(`initializing variable
found   : Foo
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassTypeOfThisInMethod()
func (c *TypeCheckNoTranspileTest) TestClassTypeOfThisInMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  foo() {
    var /** null */ foo = this;
  }
}
`).
		AddDiagnostic(`initializing variable
found   : Foo
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassTypeOfThisInStaticMethod()
func (c *TypeCheckNoTranspileTest) TestClassTypeOfThisInStaticMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  static foo() {
    var /** null */ foo = this;
  }
}
`).
		AddDiagnostic(`initializing variable
found   : (typeof Foo)
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassGetter()
func (c *TypeCheckNoTranspileTest) TestClassGetter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  get x() {}
}
var /** null */ y = new C().x;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassGetterMismatch()
func (c *TypeCheckNoTranspileTest) TestClassGetterMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @return {number} */
  get x() {}
}
var /** null */ y = new C().x;
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticGetter()
func (c *TypeCheckNoTranspileTest) TestClassStaticGetter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  static get x() {}
}
var /** null */ y = C.x;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticGetterMismatch()
func (c *TypeCheckNoTranspileTest) TestClassStaticGetterMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @return {number} */
  static get x() {}
}
var /** null */ y = C.x;
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSetter()
func (c *TypeCheckNoTranspileTest) TestClassSetter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  set x(arg) {}
}
new C().x = null;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSetterMismatch()
func (c *TypeCheckNoTranspileTest) TestClassSetterMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @param {number} arg */
  set x(arg) {}
}
new C().x = null;
`).
		AddDiagnostic(`assignment to property x of C
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSetterWithMissingParameter()
func (c *TypeCheckNoTranspileTest) TestClassSetterWithMissingParameter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class C { set a(b) {} }").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticSetter()
func (c *TypeCheckNoTranspileTest) TestClassStaticSetter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  static set x(arg) {}
}
C.x = null;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassStaticSetterMismatch()
func (c *TypeCheckNoTranspileTest) TestClassStaticSetterMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @param {number} arg */
  static set x(arg) {}
}
C.x = null;
`).
		AddDiagnostic(`assignment to property x of C
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassGetterAndSetter()
func (c *TypeCheckNoTranspileTest) TestClassGetterAndSetter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @return {number} */
  get x() {}
  /** @param {number} arg */
  set x(arg) {}
}
var /** number */ y = new C().x;
new C().x = 42;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassGetterAndSetterNoJsDoc()
func (c *TypeCheckNoTranspileTest) TestClassGetterAndSetterNoJsDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  get x() {}
  set x(arg) {}
}
var /** number */ y = new C().x;
new C().x = 42;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassGetterAndSetterDifferentTypes()
func (c *TypeCheckNoTranspileTest) TestClassGetterAndSetterDifferentTypes() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @return {number} */
  get x() {}
  /** @param {string} arg */
  set x(arg) {}
}
var /** null */ y = new C().x;
new C().x = null;
`).
		AddDiagnostic(`The types of the getter and setter for property 'x' do not match.
getter type is: number
setter type is: string
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		AddDiagnostic(`assignment to property x of C
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassGetterAndSetterWithSameStructuralTypeIsAllowed()
func (c *TypeCheckNoTranspileTest) TestClassGetterAndSetterWithSameStructuralTypeIsAllowed() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @return {{x: number}} */
  get x() { return {x: 0}; }
  /** @param {{x: number}} arg */
  set x(arg) {}
}
const c = new C();
c.x = {x: 3};
const /** {x: number} */ something = c.x;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassGetterAndSetterWithSameStructuralRecordAndNominalTypeIsAllowed()
func (c *TypeCheckNoTranspileTest) TestClassGetterAndSetterWithSameStructuralRecordAndNominalTypeIsAllowed() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
function xRecord() {}
/** @type {number} */
xRecord.prototype.x;

class C {
  /** @return {!xRecord} */
  get x() { return {x: 0}; }
  /** @param {{x: number}} arg */
  set x(arg) {}
}
const c = new C();
c.x = {x: 3};
const /** {x: number} */ something = c.x;
`).
		AddDiagnostic(`The types of the getter and setter for property 'x' do not match.
getter type is: xRecord
setter type is: {x: number}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassGetterWithShadowingDeclarationOnInstanceType()
func (c *TypeCheckNoTranspileTest) TestClassGetterWithShadowingDeclarationOnInstanceType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  constructor() {
    /** @type {string|undefined} */
    this.x;
  }
  /** @return {number} */
  get x() { return 0; }
}
new C().x = null;
`).
		AddDiagnostic(`assignment to property x of C
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassGetterWithDuplicateDeclarationLater()
func (c *TypeCheckNoTranspileTest) TestClassGetterWithDuplicateDeclarationLater() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @return {number} */
  get x() {}
}
/** @type {string} */
C.prototype.x;
new C().x = null;
`).
		AddDiagnostic(`assignment to property x of C
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassNewTargetInArrowFunction()
func (c *TypeCheckNoTranspileTest) TestClassNewTargetInArrowFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const f = () => { const /** null */ x = new.target; };").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassNewTargetInMethod()
func (c *TypeCheckNoTranspileTest) TestClassNewTargetInMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class Foo { foo() { const /** null */ x = new.target; } }").
		AddDiagnostic(`initializing variable
found   : undefined
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassNewTargetInVanillaFunction()
func (c *TypeCheckNoTranspileTest) TestClassNewTargetInVanillaFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f() { const /** null */ x = new.target; }").
		AddDiagnostic(`initializing variable
found   : (Function|undefined)
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassNewTargetInVanillaFunctionNestedArrow()
func (c *TypeCheckNoTranspileTest) TestClassNewTargetInVanillaFunctionNestedArrow() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f() { const f = () => { const /** null */ x = new.target; }; }").
		AddDiagnostic(`initializing variable
found   : (Function|undefined)
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassNewTargetInConstructor()
func (c *TypeCheckNoTranspileTest) TestClassNewTargetInConstructor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class Foo { constructor() { const /** null */ x = new.target; } };").
		AddDiagnostic(`initializing variable
found   : Function
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassNewTargetInConstructorNestedArrow()
func (c *TypeCheckNoTranspileTest) TestClassNewTargetInConstructorNestedArrow() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class Foo { constructor() { const f = () => { const /** null */ x = new.target; }; } };").
		AddDiagnostic(`initializing variable
found   : Function
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassEs5ClassCannotExtendEs6Class()
func (c *TypeCheckNoTranspileTest) TestClassEs5ClassCannotExtendEs6Class() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Base {}
/** @constructor @extends {Base} */
function Sub() {}
`).
		AddDiagnostic("ES5 class Sub cannot extend ES6 class Base").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassEs5ClassCanImplementEs6Interface()
func (c *TypeCheckNoTranspileTest) TestClassEs5ClassCanImplementEs6Interface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Inter {}
/** @constructor @implements {Inter} */
function Sub() {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsForwardReferencedClass()
func (c *TypeCheckNoTranspileTest) TestClassExtendsForwardReferencedClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
(function() {
  ns.Base = class {};
})();
class Sub extends ns.Base {}
var /** !ns.Base */ x = new Sub();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsItself()
func (c *TypeCheckNoTranspileTest) TestClassExtendsItself() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class Foo extends Foo {}").
		AddDiagnostic("Cycle detected in inheritance chain of type Foo").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsCycle()
func (c *TypeCheckNoTranspileTest) TestClassExtendsCycle() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo extends Bar {}
class Bar extends Foo {}
`).
		AddDiagnostic("Cycle detected in inheritance chain of type Bar").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsCycleOnlyInJsdoc()
func (c *TypeCheckNoTranspileTest) TestClassExtendsCycleOnlyInJsdoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Bar {}
/** @extends {Foo} */
class Foo extends Bar {}
`).
		AddDiagnostic("Cycle detected in inheritance chain of type Foo").
		AddDiagnostic("Could not resolve type in @extends tag of Foo").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsCycleOnlyInAst()
func (c *TypeCheckNoTranspileTest) TestClassExtendsCycleOnlyInAst() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Bar {}
/** @extends {Bar} */
class Foo extends Foo {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMixinFunction()
func (c *TypeCheckNoTranspileTest) TestMixinFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(new: ?, ...?)} ctor */
function mixin(ctor) {
// ctor isn't properly declared as a type,
// but we shouldn't generate an error,
// because it is a real value, not an annotation,
// and we need this coding pattern to work.
  class Foo extends ctor {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassImplementsForwardReferencedInterface()
func (c *TypeCheckNoTranspileTest) TestClassImplementsForwardReferencedInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
(function() {
  /** @interface */
  ns.Base = class {};
})();
/** @implements {ns.Base} */
class Sub {}
var /** !ns.Base */ x = new Sub();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassSuperCallResult()
func (c *TypeCheckNoTranspileTest) TestClassSuperCallResult() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Bar {}
class Foo extends Bar {
  constructor() {
    var /** null */ x = super();
  }
}
`).
		AddDiagnostic(`initializing variable
found   : Bar
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassComputedSymbolPropAllowed()
func (c *TypeCheckNoTranspileTest) TestClassComputedSymbolPropAllowed() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddIterable().Build()).
		AddSource("class Foo { [Symbol.iterator]() {} }").
		Run()
}

// Java: TypeCheckNoTranspileTest.testClassExtendsNonNativeObject()
func (c *TypeCheckNoTranspileTest) TestClassExtendsNonNativeObject() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Object {}
class Foo extends Object {
  /** @param {string} msg */
  constructor(msg) {
    super();
    this.msg = msg;
  }
}
`).
		AddDiagnostic(`attempted re-definition of type Object
found   : (typeof Object)
expected: (typeof Object)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunctionWithoutJSDoc()
func (c *TypeCheckNoTranspileTest) TestAsyncFunctionWithoutJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("async function f() { return 3; }").
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunctionInferredToReturnPromise_noExplicitReturns()
func (c *TypeCheckNoTranspileTest) TestAsyncFunctionInferredToReturnPromise_noExplicitReturns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("async function f() {} var /** null */ n = f();").
		AddDiagnostic(`initializing variable
found   : Promise<undefined>
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunctionInferredToReturnPromise_withExplicitReturns()
func (c *TypeCheckNoTranspileTest) TestAsyncFunctionInferredToReturnPromise_withExplicitReturns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("async function f() { return 3; } var /** null */ n = f();").
		AddDiagnostic(`initializing variable
found   : Promise<?>
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunction_cannotDeclareReturnToBe_Number()
func (c *TypeCheckNoTranspileTest) TestAsyncFunction_cannotDeclareReturnToBe_Number() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {number} */ async function f() {}").
		AddDiagnostic(`The return type of an async function must be a supertype of Promise
found: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunction_cannotDeclareReturnToBe_Array()
func (c *TypeCheckNoTranspileTest) TestAsyncFunction_cannotDeclareReturnToBe_Array() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Array} */ async function f() {}").
		AddDiagnostic(`The return type of an async function must be a supertype of Promise
found: Array
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunction_canDeclareReturnToBe_Object_andAccepts_undefined()
func (c *TypeCheckNoTranspileTest) TestAsyncFunction_canDeclareReturnToBe_Object_andAccepts_undefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Object} */ async function f() { return undefined; }").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunction_canDeclareReturnToBe_allType_andAccepts_undefined()
func (c *TypeCheckNoTranspileTest) TestAsyncFunction_canDeclareReturnToBe_allType_andAccepts_undefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {*} */ async function f() { return undefined; }").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncReturnsPromise1()
func (c *TypeCheckNoTranspileTest) TestAsyncReturnsPromise1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Promise<number>} */
async function getANumber() {
  return 1;
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncReturnsPromise2()
func (c *TypeCheckNoTranspileTest) TestAsyncReturnsPromise2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Promise<string>} */
async function getAString() {
  return 1;
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: (IThenable<string>|string)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunction_canDeclareReturnToBe_nullablePromise()
func (c *TypeCheckNoTranspileTest) TestAsyncFunction_canDeclareReturnToBe_nullablePromise() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {?Promise<string>} */
async function getAString() {
  return '';
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunction_canDeclareReturnToBe_unionOfPromiseAndNumber()
func (c *TypeCheckNoTranspileTest) TestAsyncFunction_canDeclareReturnToBe_unionOfPromiseAndNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {(number|!Promise<number>)} */
async function getAString() {
  return 1;
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testSuperclassMixinDoesntCollideWithAnotherScope()
func (c *TypeCheckNoTranspileTest) TestSuperclassMixinDoesntCollideWithAnotherScope() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class ParentOne {}
class ParentTwo {}
function fn(x) {
  let klass;
  if (x) {
    /** @constructor @extends {ParentOne} */
    let templatizedBase = SomeVar;
    klass = class TE extends templatizedBase {};
  } else {
    /** @constructor @extends {ParentTwo} */
    let templatizedBase = OtherVar;
    klass = class TY extends templatizedBase {};
  }
}
`).
		AddDiagnostic(`mismatch in declaration of superclass type
found   : templatizedBase
required: templatizedBase
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunction_cannotDeclareReturnToBe_aSubtypeOfPromise()
func (c *TypeCheckNoTranspileTest) TestAsyncFunction_cannotDeclareReturnToBe_aSubtypeOfPromise() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @extends {Promise<string>} */
class MyPromise extends Promise { }

/** @return {!MyPromise} */
async function getAString() {
  return '';
}
`).
		AddDiagnostic(`The return type of an async function must be a supertype of Promise
found: MyPromise
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunction_cannotDeclareReturnToBe_aSiblingOfPromise()
func (c *TypeCheckNoTranspileTest) TestAsyncFunction_cannotDeclareReturnToBe_aSiblingOfPromise() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @interface
 * @extends {IThenable<string>}
 */
class MyThenable { }

/** @return {!MyThenable} */
async function getAString() {
  return '';
}
`).
		AddDiagnostic(`The return type of an async function must be a supertype of Promise
found: MyThenable
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunction_canDeclareReturnToBe_IThenable1()
func (c *TypeCheckNoTranspileTest) TestAsyncFunction_canDeclareReturnToBe_IThenable1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!IThenable<string>} */
async function getAString() {
  return 1;
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: (IThenable<string>|string)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncFunction_checksReturnExpressionType_againstCorrectUpperBound()
func (c *TypeCheckNoTranspileTest) TestAsyncFunction_checksReturnExpressionType_againstCorrectUpperBound() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`/** @return {string|!IThenable<boolean|undefined>|!Promise<null>} */
async function getAString() {
  return {};
}
`).
		AddDiagnostic(`inconsistent return type
found   : {}
required: (IThenable<(boolean|null|undefined)>|boolean|null|undefined)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncReturnStatementIsResolved()
func (c *TypeCheckNoTranspileTest) TestAsyncReturnStatementIsResolved() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!IThenable<string>} */
async function getAString(/** !IThenable<number> */ iThenable) {
  return iThenable;
}
`).
		AddDiagnostic(`inconsistent return type
found   : IThenable<number>
required: (IThenable<string>|string)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitPromiseOfNumber1()
func (c *TypeCheckNoTranspileTest) TestAwaitPromiseOfNumber1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f(/** !Promise<number> */ p) {
  takesNumber(await p);
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitPromiseOfNumber2()
func (c *TypeCheckNoTranspileTest) TestAwaitPromiseOfNumber2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f(/** !Promise<string> */ p) {
  takesNumber(await p);
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitPromiseOfPromise()
func (c *TypeCheckNoTranspileTest) TestAwaitPromiseOfPromise() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f(/** !Promise<!Promise<number>> */ p) {
  takesNumber(await p);
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitPromiseOfUnknown()
func (c *TypeCheckNoTranspileTest) TestAwaitPromiseOfUnknown() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f(/** !Promise<?> */ p) {
  takesNumber(await p);
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitIThenable()
func (c *TypeCheckNoTranspileTest) TestAwaitIThenable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f(/** !IThenable<string> */ p) {
  takesNumber(await p);
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitNumber()
func (c *TypeCheckNoTranspileTest) TestAwaitNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f(/** string */ str) {
  takesNumber(await str);
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitDoesTypeInferenceWithin()
func (c *TypeCheckNoTranspileTest) TestAwaitDoesTypeInferenceWithin() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f() {
  var x = 1;
  await (x = 'some string'); // test we recognize that "x" is now a string.
  takesNumber(x);
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitUnionType1()
func (c *TypeCheckNoTranspileTest) TestAwaitUnionType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f(/** (number|!Promise<number>) */ param) {
  takesNumber(await param);
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitUnionType2()
func (c *TypeCheckNoTranspileTest) TestAwaitUnionType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f(/** (string|!Promise<number>) */ param) {
  takesNumber(await param);
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : (number|string)
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitUnionType3()
func (c *TypeCheckNoTranspileTest) TestAwaitUnionType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f(/** (number|!Promise<string>) */ param) {
  takesNumber(await param);
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : (number|string)
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitUnionOfPromiseAndIThenable()
func (c *TypeCheckNoTranspileTest) TestAwaitUnionOfPromiseAndIThenable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number*/ num) {}

async function f(/** (!IThenable<number>|!Promise<string>) */ param) {
  takesNumber(await param);
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : (number|string)
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitNullableIThenable()
func (c *TypeCheckNoTranspileTest) TestAwaitNullableIThenable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number */ n) {}

async function main(/** ?IThenable<number> */ iThenable) {
  takesNumber(await iThenable);
}
`).
		AddDiagnostic(`actual parameter 1 of takesNumber does not match formal parameter
found   : (null|number)
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAwaitThenable()
func (c *TypeCheckNoTranspileTest) TestAwaitThenable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesNumber(/** number */ n) {}

async function f(/** {then: string} */ thenable) {
  takesNumber(await thenable);
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultParameterWithCorrectType()
func (c *TypeCheckNoTranspileTest) TestDefaultParameterWithCorrectType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(/** number= */ n = 3) {}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultParameterWithWrongType()
func (c *TypeCheckNoTranspileTest) TestDefaultParameterWithWrongType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(/** number= */ n = 'foo') {}").
		AddDiagnostic(`default value has wrong type
found   : string
required: (number|undefined)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultParameterIsUndefined()
func (c *TypeCheckNoTranspileTest) TestDefaultParameterIsUndefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(/** number= */ n = undefined) {}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultParameter_IsVariableTypedAsUndefined()
func (c *TypeCheckNoTranspileTest) TestDefaultParameter_IsVariableTypedAsUndefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const alsoUndefined = undefined;

/** @param {string=} str */
function f(str = alsoUndefined) {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultParameter_IsKnownNotUndefinedInClosure()
func (c *TypeCheckNoTranspileTest) TestDefaultParameter_IsKnownNotUndefinedInClosure() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function takesString(/** string */ str) {}

/** @param {string=} str */
function f(str = '') {
  return () => takesString(str);
}
`).
		AddDiagnostic(`actual parameter 1 of takesString does not match formal parameter
found   : (string|undefined)
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultParameterInDestructuringIsUndefined()
func (c *TypeCheckNoTranspileTest) TestDefaultParameterInDestructuringIsUndefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{prop: (string|undefined)}} obj */
function f({prop = undefined}) {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultParameterInDestructuringIsVariableTypedUndefined()
func (c *TypeCheckNoTranspileTest) TestDefaultParameterInDestructuringIsVariableTypedUndefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const alsoUndefined = undefined;

/** @param {{prop: (string|undefined)}} obj */
function f({prop = alsoUndefined}) {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeCheckingOccursWithinDefaultParameter()
func (c *TypeCheckNoTranspileTest) TestTypeCheckingOccursWithinDefaultParameter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("let /** number */ age = 0; function f(x = age = 'foo') {}").
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultParameterWithTypeInferredFromCallback()
func (c *TypeCheckNoTranspileTest) TestDefaultParameterWithTypeInferredFromCallback() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** function(number=) */ callback) {}

f((x = 3) => {
  var /** number */ y = x;
})
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultParameterInIifeWithInferredType()
func (c *TypeCheckNoTranspileTest) TestDefaultParameterInIifeWithInferredType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var /** string|undefined */ stringOrUndefined;
(function f(x = 3) {
  var /** string */ str = x;
})(stringOrUndefined);
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultParameterWithNoJSDocTreatedAsOptional()
func (c *TypeCheckNoTranspileTest) TestDefaultParameterWithNoJSDocTreatedAsOptional() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(a = 3) {} f();").
		Run()
}

// Java: TypeCheckNoTranspileTest.testOverridingMethodHasDefaultParameterInPlaceOfRequiredParam()
func (c *TypeCheckNoTranspileTest) TestOverridingMethodHasDefaultParameterInPlaceOfRequiredParam() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Parent {
  /** @param {number} num */
  f(num) {}
}
class Child extends Parent {
  /** @override */
  f(num = undefined) {}
}
(new Child()).f();
(new Child()).f(undefined);
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOverridingMethodAddsOptionalParameterWithDefaultValue()
func (c *TypeCheckNoTranspileTest) TestOverridingMethodAddsOptionalParameterWithDefaultValue() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Parent {
  /** @param {number} num */
  f(num) {}
}
class Child extends Parent {
  /** @override */
  f(num, otherParam = undefined) {}
}
(new Child()).f(3);
(new Child()).f(3, 'str');
(new Child()).f(3, undefined);
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOptionalDestructuringParameterWithoutDefaultValue()
func (c *TypeCheckNoTranspileTest) TestOptionalDestructuringParameterWithoutDefaultValue() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {{x: number}=} opts */ function f({x}) {} f();").
		Run()
}

// Java: TypeCheckNoTranspileTest.testBasicArrayPatternDeclaration()
func (c *TypeCheckNoTranspileTest) TestBasicArrayPatternDeclaration() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterable<number> */ numbers) {
  const [/** string */ str] = numbers;
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testNestedDestructuringPatternDeclaration()
func (c *TypeCheckNoTranspileTest) TestNestedDestructuringPatternDeclaration() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterable<{x: number}> */ xNumberObjs) {
  const [{/** string */ x}] = xNumberObjs;
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayPatternDeclarationWithElision()
func (c *TypeCheckNoTranspileTest) TestArrayPatternDeclarationWithElision() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterable<number> */ numbers) {
  const [, /** number */ x, , /** string */ y] = numbers;
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testBasicArrayPatternAssign()
func (c *TypeCheckNoTranspileTest) TestBasicArrayPatternAssign() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterable<number> */ numbers) {
  var /** string */ str;
  [str] = numbers;
}
`).
		AddDiagnostic(`assignment
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testNestedDestructuringPatternAssign()
func (c *TypeCheckNoTranspileTest) TestNestedDestructuringPatternAssign() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterable<{x: number}> */ xNumberObjs) {
  var /** string */ x;
  [{x}] = xNumberObjs;
}
`).
		AddDiagnostic(`assignment
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testValidArrayPatternInForOfInitializer()
func (c *TypeCheckNoTranspileTest) TestValidArrayPatternInForOfInitializer() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterable<!Iterable<number>> */ numberLists) {
  for (const [/** number */ x, /** number */ y] of numberLists) {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayPatternInForOfInitializerWithTypeMismatch()
func (c *TypeCheckNoTranspileTest) TestArrayPatternInForOfInitializerWithTypeMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterable<!Iterable<number>> */ numberLists) {
  for (const [/** number */ x, /** string */ y] of numberLists) {}
}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testValidArrayPatternInForInInitializer()
func (c *TypeCheckNoTranspileTest) TestValidArrayPatternInForInInitializer() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`function f(/** !Object<string, number> */ obj) {
  for (const [/** string */ a, /** string */ b] in obj) {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayPatternInForInInitializerWithTypeMismatch()
func (c *TypeCheckNoTranspileTest) TestArrayPatternInForInInitializerWithTypeMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`function f(/** !Object<string, number> */ obj) {
  for (const [/** number */ a, /** number */ b] in obj) {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testBadDefaultValueInCatchCausesWarning()
func (c *TypeCheckNoTranspileTest) TestBadDefaultValueInCatchCausesWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("try { throw {x: undefined}; } catch ({/** string */ x = 3 + 4}) {}").
		AddDiagnostic(`default value has wrong type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testCannotAliasEnumThroughDestructuring()
func (c *TypeCheckNoTranspileTest) TestCannotAliasEnumThroughDestructuring() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/** @enum {number} */ const THINGS = {THING1: 1, THING2: 2};
// TODO(lharker): warn for putting @enum here
/** @enum */ const [OTHERTHINGS] = [THINGS];
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayPatternAssign_badInterfacePropertyCreation()
func (c *TypeCheckNoTranspileTest) TestArrayPatternAssign_badInterfacePropertyCreation() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource("/** @interface */ function Foo() {}; [Foo.prototype.bar] = [];").
		AddDiagnostic("interface members can only be empty property declarations, empty functions, or goog.abstractMethod").
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayPatternAssign_badPropertyAssignment()
func (c *TypeCheckNoTranspileTest) TestArrayPatternAssign_badPropertyAssignment() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Iterable<number>} numbers */
function f(numbers) {
  const /** {a: string} */ obj = {a: 'foo'};
  [obj.a] = numbers;
}
`).
		AddDiagnostic(`assignment to property a of obj
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testBadComputedPropertyKeyTypeInObjectPattern()
func (c *TypeCheckNoTranspileTest) TestBadComputedPropertyKeyTypeInObjectPattern() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const {[{}]: x} = {};").
		AddDiagnostic(`property access
found   : {}
required: (string|symbol)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testComputedPropertyAccessOnStructInObjectPattern()
func (c *TypeCheckNoTranspileTest) TestComputedPropertyAccessOnStructInObjectPattern() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @struct */ const myStruct = {a: 1}; const {['a']: a} = myStruct;").
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckNoTranspileTest.testQuotedPropertyAccessOnStructInObjectPattern()
func (c *TypeCheckNoTranspileTest) TestQuotedPropertyAccessOnStructInObjectPattern() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @struct */ const myStruct = {a: 1}; const {'a': a} = myStruct;").
		AddDiagnostic("Cannot do quoted access on a struct").
		Run()
}

// Java: TypeCheckNoTranspileTest.testNonQuotedPropertyAccessOnDictInObjectPattern()
func (c *TypeCheckNoTranspileTest) TestNonQuotedPropertyAccessOnDictInObjectPattern() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict*/ const myDict = {'a': 1}; const {a} = myDict;").
		AddDiagnostic("Cannot do unquoted access on a dict").
		Run()
}

// Java: TypeCheckNoTranspileTest.testRestrictedIndexTypeInComputedPropertyKeyInObjectPattern()
func (c *TypeCheckNoTranspileTest) TestRestrictedIndexTypeInComputedPropertyKeyInObjectPattern() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const /** !Object<number, number> */ obj = {3: 3, 4: 4};
const {['string']: x} = obj;
`).
		AddDiagnostic(`restricted index type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectDestructuringNullInDeclarationCausesWarning()
func (c *TypeCheckNoTranspileTest) TestObjectDestructuringNullInDeclarationCausesWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const {} = null;").
		AddDiagnostic(`cannot destructure 'null' or 'undefined'
found   : null
required: Object
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectDestructuringUndefinedInDeclarationCausesWarning()
func (c *TypeCheckNoTranspileTest) TestObjectDestructuringUndefinedInDeclarationCausesWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const {} = undefined;").
		AddDiagnostic(`cannot destructure 'null' or 'undefined'
found   : undefined
required: Object
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectDestructuringBoxablePrimitiveInDeclarationIsAllowed()
func (c *TypeCheckNoTranspileTest) TestObjectDestructuringBoxablePrimitiveInDeclarationIsAllowed() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const {} = 0;").
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectDestructuringNullInParametersCausesWarning()
func (c *TypeCheckNoTranspileTest) TestObjectDestructuringNullInParametersCausesWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {null} obj */ function f({}) {}").
		AddDiagnostic(`cannot destructure 'null' or 'undefined'
found   : null
required: Object
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectDestructuringNullInNestedPatternCausesWarning()
func (c *TypeCheckNoTranspileTest) TestObjectDestructuringNullInNestedPatternCausesWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const {a: {}} = {a: null};").
		AddDiagnostic(`cannot destructure 'null' or 'undefined'
found   : null
required: Object
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectDestructuringNullableDoesntCauseWarning()
func (c *TypeCheckNoTranspileTest) TestObjectDestructuringNullableDoesntCauseWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** ?{x: number} */ nullableObj) {
const {x} = nullableObj;
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayDestructuringNonIterableCausesWarning()
func (c *TypeCheckNoTranspileTest) TestArrayDestructuringNonIterableCausesWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const [] = 3;").
		AddDiagnostic(`array pattern destructuring requires an Iterable
found   : number
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayDestructuringStringIsAllowed()
func (c *TypeCheckNoTranspileTest) TestArrayDestructuringStringIsAllowed() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource("const [] = 'foobar';").
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayDestructuringNonIterableInParametersCausesWarning()
func (c *TypeCheckNoTranspileTest) TestArrayDestructuringNonIterableInParametersCausesWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {number} arr */ function f([]) {}").
		AddDiagnostic(`array pattern destructuring requires an Iterable
found   : number
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayDestructuringNonIterableInForOfLoopCausesWarning()
func (c *TypeCheckNoTranspileTest) TestArrayDestructuringNonIterableInForOfLoopCausesWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource("const /** !Iterable<number> */ iter = [0]; for (const [] of iter) {}").
		AddDiagnostic(`array pattern destructuring requires an Iterable
found   : number
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testBadDefaultValueInArrayPatternCausesWarning()
func (c *TypeCheckNoTranspileTest) TestBadDefaultValueInArrayPatternCausesWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource("const [/** string */ foo = 0] = [];").
		AddDiagnostic(`default value has wrong type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultValueForNestedArrayPatternMustBeIterable()
func (c *TypeCheckNoTranspileTest) TestDefaultValueForNestedArrayPatternMustBeIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource("const [[] = 0] = [];").
		AddDiagnostic(`array pattern destructuring requires an Iterable
found   : number
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayPatternParameterCanBeOptional()
func (c *TypeCheckNoTranspileTest) TestArrayPatternParameterCanBeOptional() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource("/** @param {!Array<string>=} arr */ function f([x, y] = []) {}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultValueForArrayPatternParameterMustBeIterable()
func (c *TypeCheckNoTranspileTest) TestDefaultValueForArrayPatternParameterMustBeIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource("function f([] = 0) {}").
		AddDiagnostic(`array pattern destructuring requires an Iterable
found   : number
required: Iterable
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultValueForNestedObjectPatternMustNotBeNull()
func (c *TypeCheckNoTranspileTest) TestDefaultValueForNestedObjectPatternMustNotBeNull() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource("const [{} = null] = [];").
		AddDiagnostic(`cannot destructure a 'null' or 'undefined' default value
found   : null
required: Object
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDefaultValueForObjectPatternParameterMustNotBeNull()
func (c *TypeCheckNoTranspileTest) TestDefaultValueForObjectPatternParameterMustNotBeNull() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f({} = null) {}").
		AddDiagnostic(`cannot destructure a 'null' or 'undefined' default value
found   : null
required: Object
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayDestructuringParameterWithElision()
func (c *TypeCheckNoTranspileTest) TestArrayDestructuringParameterWithElision() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource("/** @param {!Array<number>} numbers */ function f([, x, , y]) {}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectPatternDeclarationWithMissingPropertyWarning()
func (c *TypeCheckNoTranspileTest) TestObjectPatternDeclarationWithMissingPropertyWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** {a: number} */ obj) {
  const {a, b} = obj;
}
`).
		AddDiagnostic("Property b never defined on obj").
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectPatternDoesntCheckMissingPropertyForQuotedStringKey()
func (c *TypeCheckNoTranspileTest) TestObjectPatternDoesntCheckMissingPropertyForQuotedStringKey() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** {a: number} */ obj) {
  const {a, 'b': b} = obj;
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectPatternAssignWithMissingPropertyWarning()
func (c *TypeCheckNoTranspileTest) TestObjectPatternAssignWithMissingPropertyWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** {a: number} */ obj) {
  let a, b;
  ({a, b} = obj);
}
`).
		AddDiagnostic("Property b never defined on obj").
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectPatternDeclarationWithMissingPropertyWarningInForOf()
func (c *TypeCheckNoTranspileTest) TestObjectPatternDeclarationWithMissingPropertyWarningInForOf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterable<{a: number}> */ aNumberObj) {
  for (const {a, b} of aNumberObj) {}
}
`).
		AddDiagnostic("Property b never defined on {a: number}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfWithDestructuring()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfWithDestructuring() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`async function f(/** !Iterable<Promise<{a: number}>> */ o) {
  for await (const {a, b} of o) {}
}
`).
		AddDiagnostic("Property b never defined on {a: number}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectPatternWithMissingPropertyWarningInParameters()
func (c *TypeCheckNoTranspileTest) TestObjectPatternWithMissingPropertyWarningInParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{a: number}} obj */
function f(/** {a: number} */ {b}) {}
`).
		AddDiagnostic("Property b never defined on {a: number}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testArrayPatternAssignWithIllegalPropCreationInStruct()
func (c *TypeCheckNoTranspileTest) TestArrayPatternAssignWithIllegalPropCreationInStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  f(/** !Iterable<number> */ ) {
    [this.x] = arr;
  }
}
`).
		AddDiagnostic("Cannot add a property to a struct instance after it is constructed. (If you already declared the property, make sure to give it a type.)").
		AddDiagnostic("Property x never defined on C").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefAliasedThroughDestructuringFromLegacyNamespacePassesTypechecking()
func (c *TypeCheckNoTranspileTest) TestTypedefAliasedThroughDestructuringFromLegacyNamespacePassesTypechecking() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {number} */
let TypeOriginal;
class clientOpClass {}

const clientOp = clientOpClass;
/** @const */
clientOp.Type = TypeOriginal;
// The above pattern mimics some goog.module.declareLegacyNamespace() code

const {Type} = clientOp;
class C {
  /** @param {!Type} type */
  m(type) {
    type = 'cause a type error';
  }
}
`).
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testEnumAliasedThroughDestructuringPassesTypechecking()
func (c *TypeCheckNoTranspileTest) TestEnumAliasedThroughDestructuringPassesTypechecking() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const ns = {};
/** @enum {number} */
ns.myEnum = {FOO: 1, BAR: 2};

const {myEnum} = ns;
const /** myEnum */ n = myEnum.FOO;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testEnumAliasedThroughDestructuringReportsCorrectMissingPropWarning()
func (c *TypeCheckNoTranspileTest) TestEnumAliasedThroughDestructuringReportsCorrectMissingPropWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const ns = {};
/** @enum {number} */
ns.myEnum = {FOO: 1, BAR: 2};

const {myEnum} = ns;
const missing = myEnum.MISSING;
`).
		AddDiagnostic("element MISSING does not exist on this enum").
		Run()
}

// Java: TypeCheckNoTranspileTest.testAnnotatedObjectLiteralInDefaultParameterInitializer()
func (c *TypeCheckNoTranspileTest) TestAnnotatedObjectLiteralInDefaultParameterInitializer() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{g: function(number): undefined}=} x */
function f(x = {/** @param {string} x */ g(x) {}}) {}
`).
		AddDiagnostic(`default value has wrong type
found   : {g: function(string): undefined}
required: (undefined|{g: function(number): undefined})
missing : []
mismatch: [g]
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDictClass1()
func (c *TypeCheckNoTranspileTest) TestDictClass1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ var C = class { constructor() {} 'x'(){} };").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeCheckingOverriddenGetterFromSuperclass()
func (c *TypeCheckNoTranspileTest) TestTypeCheckingOverriddenGetterFromSuperclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Bar {
  /**
   * @abstract
   * @return {number}
   */
  get num() { return 1; }
}
/** @extends {Bar} */
class Baz extends Bar {
  /** @override */
  get num() { return 3; }
}
var /** string */ x = (new Baz).num;
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeCheckingOverriddenGetterFromSuperclassWithBadReturnType()
func (c *TypeCheckNoTranspileTest) TestTypeCheckingOverriddenGetterFromSuperclassWithBadReturnType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Bar {
  /**
   * @abstract
   * @return {number}
   */
  get num() { return 1; }
}
/** @extends {Bar} */
class Baz extends Bar {
  /** @override */
  get num() { return 'foo'; }
}
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGetterOverridesPrototypePropertyFromInterface()
func (c *TypeCheckNoTranspileTest) TestGetterOverridesPrototypePropertyFromInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Bar {}
/** @type {number} */
Bar.prototype.num;

/** @implements {Bar} */
class Baz {
  /** @override */
  get num() { return 3; }
}
var /** string */ x = (new Baz).num;
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGetterOverridesInstancePropertyFromInterface()
func (c *TypeCheckNoTranspileTest) TestGetterOverridesInstancePropertyFromInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Bar {
  constructor() {
    /** @type {number} */
    this.num;
  }
}
/** @implements {Bar} */
class Baz {
  /** @override */
  get num() { return 3; }
}
var /** string */ x = (new Baz).num;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOverriddenSetterFromSuperclass()
func (c *TypeCheckNoTranspileTest) TestOverriddenSetterFromSuperclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Bar {
  /**
   * @abstract
   * @param {number} x
   */
  set num(x) {}
}
/** @extends {Bar} */
class Baz extends Bar {
  /** @override */
  set num(x) {}
}
(new Baz).num = 'foo';
`).
		AddDiagnostic(`assignment to property num of Baz
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGetterOverridesMethod()
func (c *TypeCheckNoTranspileTest) TestGetterOverridesMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Bar {
  /** @return {number} */
  num() { return 1; }
}
/** @extends {Bar} */
class Baz extends Bar {
  /** @override */
  get num() { return 1; }
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: function(this:Bar): number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMisplacedOverrideOnGetter()
func (c *TypeCheckNoTranspileTest) TestMisplacedOverrideOnGetter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Bar {}
/** @extends {Bar} */
class Baz extends Bar {
  /** @override */
  get num() { return 3; }
}
var /** string */ x = (new Baz).num;
`).
		AddDiagnostic("property num not defined on any superclass of Baz").
		Run()
}

// Java: TypeCheckNoTranspileTest.testOverridingNonMethodWithMethodDoesntBlockTypeCheckingInsideMethod()
func (c *TypeCheckNoTranspileTest) TestOverridingNonMethodWithMethodDoesntBlockTypeCheckingInsideMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
/** @type {number} */
Foo.prototype.bar = 3;

class Bar extends Foo {
  /** @override */
  bar(/** number */ n) {
    var /** string */ str = n;
  }
}
`).
		AddDiagnostic(`mismatch of the bar property type and the type of the property it overrides from superclass Foo
original: number
override: function(this:Bar, number): undefined
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGetterWithTemplateTypeReturnIsTypeChecked()
func (c *TypeCheckNoTranspileTest) TestGetterWithTemplateTypeReturnIsTypeChecked() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @template T */
class C {
  /** @return {T} */
  get t() {}
}
/** @implements {C<string>} */
class CString {
  /** @override */
  get t() { return 3; } // inconsistent return type
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testStubMethodDeclarationDoesntBlockTypecheckingOfGetter()
func (c *TypeCheckNoTranspileTest) TestStubMethodDeclarationDoesntBlockTypecheckingOfGetter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo {}
/** @return {number} */
Foo.prototype.num;
/** @implements {Foo} */
class Bar {
  /** @override */
  get num() { return 1; }
}
var /** string */ x = (new Bar).num;
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: function(this:Foo): number
`).
		AddDiagnostic(`initializing variable
found   : function(this:Foo): number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOverrideSupertypeOnAnonymousClass()
func (c *TypeCheckNoTranspileTest) TestOverrideSupertypeOnAnonymousClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function use(ctor) {}

class Foo {
  constructor() {
    /** @type {string} */
    this.str;
  }
}
use(class extends Foo {
  f() { this.str = 3; }
});
`).
		AddDiagnostic(`assignment to property str of <anonymous@testcode0:9>
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOverrideSupertypeOnClassExpression()
func (c *TypeCheckNoTranspileTest) TestOverrideSupertypeOnClassExpression() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function use(ctor) {}

class Foo {
  constructor() {
    /** @type {string} */
    this.str;
  }
}
use(class Bar extends Foo {
  f() { this.str = 3; }
});
`).
		AddDiagnostic(`assignment to property str of <anonymous@testcode0:9>
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testOverrideInferredOnClassExpression()
func (c *TypeCheckNoTranspileTest) TestOverrideInferredOnClassExpression() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function use(ctor) {}

class Foo {
  f(/** number */ num) {}
}
use(class Bar extends Foo {
  /** @override */
  f(num) {
    var /** string */ str = num;
  }
});
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMixinWithUnknownTemplatedSupertypeDoesntCauseWarning()
func (c *TypeCheckNoTranspileTest) TestMixinWithUnknownTemplatedSupertypeDoesntCauseWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template T
 * @param {function(new:T)} superClass
 */
function mixin(superClass) {
  class Changed extends superClass {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMixinImplementingInterfaceAndUnknownTemplatedSuperclass()
func (c *TypeCheckNoTranspileTest) TestMixinImplementingInterfaceAndUnknownTemplatedSuperclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template T
 * @param {function(new:T)} superClass
 */
function mixin(superClass) {
  /** @implements {ChangedInterface} */
  class Changed extends superClass {
    /**
     * @override
     * @return {number}
     */
    method() {
      return 3;
    }
  }
}

/** @interface */
class ChangedInterface {
  /** @return {number} */
  method() {}
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGlobalAliasOfEnumIsNonNullable()
func (c *TypeCheckNoTranspileTest) TestGlobalAliasOfEnumIsNonNullable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
/** @enum {number} */
Foo.E = {A: 1};
const E = Foo.E;
/** @type {E} */ let e = undefined;
`).
		AddDiagnostic(`initializing variable
found   : undefined
required: Foo.E<number>
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeNameAliasOnAliasedNamespace()
func (c *TypeCheckNoTranspileTest) TestTypeNameAliasOnAliasedNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
/** @enum {number} */
Foo.E = {A: 1};
const F = Foo;
const E = F.E;
/** @type {E} */ let e = undefined;
`).
		AddDiagnostic(`initializing variable
found   : undefined
required: Foo.E<number>
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeNamePropertyOnAliasedNamespace()
func (c *TypeCheckNoTranspileTest) TestTypeNamePropertyOnAliasedNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
/** @enum {number} */
Foo.E = {A: 1};
const F = Foo;
/** @type {F.E} */ let e = undefined;
`).
		AddDiagnostic(`initializing variable
found   : undefined
required: (Foo.E<number>|null)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefNameAliasOnAliasedNamespace()
func (c *TypeCheckNoTranspileTest) TestTypedefNameAliasOnAliasedNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
/** @typedef {number|string} */
Foo.E;
const F = Foo;
const E = F.E;
/** @type {E} */ let e = undefined;
`).
		AddDiagnostic(`initializing variable
found   : undefined
required: (number|string)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeNameAliasOnAliasedClassSideNamespace()
func (c *TypeCheckNoTranspileTest) TestTypeNameAliasOnAliasedClassSideNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
/** @enum {number} */ Foo.E = {A: 1};
class Bar extends Foo {};
const B = Bar;
const E = B.E;
/** @type {E} */ let e = undefined;
`).
		AddDiagnostic(`initializing variable
found   : undefined
required: Foo.E<number>
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForwardDeclaredGlobalAliasOfEnumIsNonNullable_constDeclaration()
func (c *TypeCheckNoTranspileTest) TestForwardDeclaredGlobalAliasOfEnumIsNonNullable_constDeclaration() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {string} */
const Colors = {RED: 'red', YELLOW: 'yellow'};
const /** ColorsAlias */ c = null
const ColorsAlias = Colors;
`).
		AddDiagnostic(`initializing variable
found   : null
required: Colors<string>
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForwardDeclaredGlobalAliasOfEnumIsNonNullable_constJSDoc()
func (c *TypeCheckNoTranspileTest) TestForwardDeclaredGlobalAliasOfEnumIsNonNullable_constJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {string} */
const Colors = {RED: 'red', YELLOW: 'yellow'};
const /** ns.ColorsAlias */ c = null
const ns = {};
/** @const */ ns.ColorsAlias = Colors;
`).
		AddDiagnostic(`initializing variable
found   : null
required: Colors<string>
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testLocalEnumDoesNotInfluenceGlobalDefaultNullablity()
func (c *TypeCheckNoTranspileTest) TestLocalEnumDoesNotInfluenceGlobalDefaultNullablity() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {};
function f() {
  /** @enum {number} */ const Foo = {A: 1};
}
/** @type {Foo} */ let x = null;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGlobalEnumDoesNotInfluenceLocalDefaultNullablity()
func (c *TypeCheckNoTranspileTest) TestGlobalEnumDoesNotInfluenceLocalDefaultNullablity() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {number} */ const Foo = {A: 1};
function f() {
  class Foo {};
  /** @type {Foo} */ let x = null;
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testLocalEnumAliasDoesNotInfluenceGlobalDefaultNullablity()
func (c *TypeCheckNoTranspileTest) TestLocalEnumAliasDoesNotInfluenceGlobalDefaultNullablity() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {};
/** @enum {number} */ const Bar = {A: 1};
function f() {
  const Foo = Bar;
}
/** @type {Foo} */ let x = null;
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypedefInExtern()
func (c *TypeCheckNoTranspileTest) TestTypedefInExtern() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("/** @typedef {boolean} */ var ConstrainBoolean;").
		AddSource("var /** ConstrainBoolean */ x = 42;").
		AddDiagnostic(`initializing variable
found   : number
required: boolean
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckNoTranspileTest.testDeeplyNestedAliases()
func (c *TypeCheckNoTranspileTest) TestDeeplyNestedAliases() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const ns = {};
/** @typedef {number} */
ns.MyNumber;
const alias = {};
/** @const */
alias.child = ns;
const outer = {};
/** @const */
outer.inner = alias;
const /** outer.inner.child.MyNumber */ x = '';
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorNoReturnOrYield()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorNoReturnOrYield() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<?>} */
async function* asyncGen() {}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorDeclaredReturnMustBeSupertypeOfAsyncGenerator()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorDeclaredReturnMustBeSupertypeOfAsyncGenerator() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncIterator<?>} */
async function* asyncGen() {}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncIterable<?>} */
async function* asyncGen() {}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncIteratorIterable<?>} */
async function* asyncGen() {}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!Object} */
async function* asyncGen() {}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {*} */
async function* asyncGen() {}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {?} */
async function* asyncGen() {}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {number} */
async function* asyncGen() {}
`).
		AddDiagnostic(`An async generator function must return a (supertype of) AsyncGenerator
found   : number
required: AsyncGenerator
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithYield()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithYield() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield 0; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithYieldStarOtherAsyncGenerator()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithYieldStarOtherAsyncGenerator() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen0() { yield 0; }

/** @return {!AsyncGenerator<number>} */
async function* asyncGen1() { yield* asyncGen0(); }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldStarNonIterable()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldStarNonIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield* 0; }
`).
		AddDiagnostic(`Expression yield* expects an iterable or async iterable
found   : number
required: (AsyncIterator|Iterator)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldStarBoxableIterable()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldStarBoxableIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().AddString().
		Build()).
		AddSource(`let /** string */ boxable;
/** @return {!AsyncGenerator<string>} */
async function* asyncGen() { yield* 'boxable'; }
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().AddArray().
		Build()).
		AddSource(`let /** !Array<number> */ boxable;
/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield* boxable; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithYieldStarSyncGenerator()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithYieldStarSyncGenerator() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!Generator<number>} */
function* gen() { yield 0; }

/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield* gen(); }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithYieldStarSyncAndAsyncGeneratorUnion()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithYieldStarSyncAndAsyncGeneratorUnion() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !Generator<string>|!AsyncGenerator<number> */ gen;

/** @return {!AsyncGenerator<string|number>} */
async function* asyncGen() { yield* gen; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithYieldStarSyncAndAsyncGeneratorAndNonGeneratorUnion()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithYieldStarSyncAndAsyncGeneratorAndNonGeneratorUnion() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !Generator<string>|!AsyncGenerator<number>|number */ gen;

/** @return {!AsyncGenerator<string|number>} */
async function* asyncGen() { yield* gen; }
`).
		AddDiagnostic(`Expression yield* expects an iterable or async iterable
found   : (AsyncGenerator<number,?,?>|Generator<string,?,?>|number)
required: (AsyncIterator|Iterator)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithYieldStarSyncAndAsyncGeneratorUnionMismatch()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithYieldStarSyncAndAsyncGeneratorUnionMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !Generator<string>|!AsyncGenerator<number> */ gen;

/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield* gen; }
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : (number|string)
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithUnspecifiedReturn()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithUnspecifiedReturn() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { return 'str'; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithMismatchReturn()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithMismatchReturn() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<string, number>} */
async function* asyncGen() { return 'str'; }
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithMismatchReturn_returnAsyncIterator()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithMismatchReturn_returnAsyncIterator() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncIterator<string, number>} */
async function* asyncGen() { return 'str'; }
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithMismatchReturn_returnAsyncIterable()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithMismatchReturn_returnAsyncIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncIterable<string, number>} */
async function* asyncGen() { return 'str'; }
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorWithMismatchYield()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorWithMismatchYield() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield 'str'; }
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : string
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldAwaitNonThenable()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldAwaitNonThenable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield await 0; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldAwaitNonThenableMismatch()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldAwaitNonThenableMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield await 'str'; }
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : string
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldAwaitThenable()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldAwaitThenable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !IThenable<number> */ thenable;
/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield await thenable; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldAwaitThenableMismatch()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldAwaitThenableMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !IThenable<string> */ thenable;
/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield await thenable; }
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : string
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldAwaitPromise()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldAwaitPromise() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield await Promise.resolve(0); }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldAwaitPromiseMismatch()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldAwaitPromiseMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield await Promise.resolve('str'); }
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : string
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldPromise()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldPromise() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield Promise.resolve(0); }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldPromiseMismatch()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldPromiseMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield Promise.resolve('str'); }
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : Promise<string>
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldIThenable()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldIThenable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !IThenable<number> */ thenable;
/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield thenable; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldThenableMismatch()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldThenableMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !IThenable<string> */ thenable;
/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield thenable; }
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : IThenable<string>
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldIThenableUnionNonIThenable()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldIThenableUnionNonIThenable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** (!IThenable<number>|string) */ thenableOrString;
/** @return {!AsyncGenerator<number|string>} */
async function* asyncGen() { yield thenable; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorYieldIThenableUnionNonIThenableMismatch()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorYieldIThenableUnionNonIThenableMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** (!IThenable<number>|string) */ thenableOrString;
/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { yield thenableOrString; }
`).
		AddDiagnostic(`Yielded type does not match declared return type.
found   : (IThenable<number>|string)
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorReturnNothing()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorReturnNothing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { return; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorReturnSameType()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorReturnSameType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number>} */
async function* asyncGen() { return 0; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorReturnMismatchType()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorReturnMismatchType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`/** @return {!AsyncGenerator<number, number>} */
async function* asyncGen() { return 'str'; }
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorReturnVoidPromise()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorReturnVoidPromise() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !Promise<void> */ voidPromise;
/** @return {!AsyncGenerator<number, number>} */
async function* asyncGen() { return voidPromise; }
`).
		AddDiagnostic(`inconsistent return type
found   : Promise<undefined>
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorReturnUndefinedPromise()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorReturnUndefinedPromise() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !Promise<undefined> */ undefPromise;
/** @return {!AsyncGenerator<number, number>} */
async function* asyncGen() { return undefPromise; }
`).
		AddDiagnostic(`inconsistent return type
found   : Promise<undefined>
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorReturnPromise()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorReturnPromise() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !Promise<number> */ promise;
/** @return {!AsyncGenerator<number, number>} */
async function* asyncGen() { return promise; }
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorReturnPromiseMismatch()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorReturnPromiseMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !Promise<string> */ promise;
/** @return {!AsyncGenerator<number, number>} */
async function* asyncGen() { return promise; }
`).
		AddDiagnostic(`inconsistent return type
found   : Promise<string>
required: (IThenable<number>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorFunctionInferredToBeAsyncGenerator()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorFunctionInferredToBeAsyncGenerator() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`async function* asyncGen() { return 0; }
let /** !AsyncGenerator<number> */ g = asyncGen();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorFunctionInferredToBeAsyncGeneratorTemplateMismatch()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorFunctionInferredToBeAsyncGeneratorTemplateMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`async function* asyncGen() { return 0; }
let /** !AsyncGenerator<string> */ g = asyncGen();
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testAsyncGeneratorFunctionInferredToBeAsyncGeneratorMismatch()
func (c *TypeCheckNoTranspileTest) TestAsyncGeneratorFunctionInferredToBeAsyncGeneratorMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`async function* asyncGen() { return 0; }
let /** null */ g = asyncGen();
`).
		AddDiagnostic(`initializing variable
found   : AsyncGenerator<?,?,?>
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testObjectSpread_typedAsObject()
func (c *TypeCheckNoTranspileTest) TestObjectSpread_typedAsObject() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`let obj = {a: 1, b: 'str'};
let /** !{a: string, b: string, c: boolean} */ copy = {c: true, ...obj, a: 'hello'};
`).
		AddDiagnostic(`initializing variable
found   : Object
required: {
  a: string,
  b: string,
  c: boolean
}
missing : [a,b,c]
mismatch: []
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfNonIterable()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfNonIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`async function foo() {
  for await (const n of 0) {
  }
}
`).
		AddDiagnostic(`Can only async iterate over a (non-null) Iterable or AsyncIterable type
found   : number
required: (AsyncIterator|Iterator)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfAsyncIterator()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfAsyncIterator() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !AsyncIterable<number> */ gen;
async function foo() {
  for await (const /** number */ n of gen) {
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfNullableAsyncIterator()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfNullableAsyncIterator() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** ?AsyncIterable<number> */ gen;
async function foo() {
  for await (const /** number */ n of gen) {
  }
}
`).
		AddDiagnostic(`Can only async iterate over a (non-null) Iterable or AsyncIterable type
found   : (AsyncIterable<number,?,?>|null)
required: (AsyncIterator|Iterator)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfAsyncIteratorMismatch()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfAsyncIteratorMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().AddString().
		Build()).
		AddSource(`let /** !AsyncIterable<string> */ gen;
async function foo() {
  for await (const /** number */ n of gen) {
  }
}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfSynchronousIterable()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfSynchronousIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !Iterable<number> */ gen;
async function foo() {
  for await (const /** number */ n of gen) {
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfSynchronousIterableMismatch()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfSynchronousIterableMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().AddString().
		Build()).
		AddSource(`let /** !Iterable<string> */ gen;
async function foo() {
  for await (const /** number */ n of gen) {
  }
}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfBoxedSynchronousIterable()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfBoxedSynchronousIterable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().AddString().
		Build()).
		AddSource(`let /** string */ gen;
async function foo() {
  for await (const /** string */ n of gen) {
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfBoxedSynchronousIterableMismatch()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfBoxedSynchronousIterableMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().AddString().
		Build()).
		AddSource(`let /** string */ gen;
async function foo() {
  for await (const /** number */ n of gen) {
  }
}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfAsyncAndSynchronousIterableUnion()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfAsyncAndSynchronousIterableUnion() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** !Iterable<number>|!AsyncIterable<string> */ gen;
async function foo() {
  for await (const /** number|string */ n of gen) {
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfAsyncAndSynchronousIterableUnionMismatch()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfAsyncAndSynchronousIterableUnionMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().AddString().
		Build()).
		AddSource(`let /** !Iterable<number>|!AsyncIterable<string> */ gen;
async function foo() {
  for await (const /** boolean */ n of gen) {
  }
}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : (number|string)
required: boolean
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfAsyncAndBoxIterableUnion()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfAsyncAndBoxIterableUnion() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().AddString().
		Build()).
		AddSource(`let /** !Iterable<number>|string */ gen;
async function foo() {
  for await (const /** number|string */ n of gen) {
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfAsyncAndBoxIterableUnionMismatch()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfAsyncAndBoxIterableUnionMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().AddString().
		Build()).
		AddSource(`let /** !Iterable<number>|string */ gen;
async function foo() {
  for await (const /** boolean */ n of gen) {
  }
}
`).
		AddDiagnostic(`declared type of for-of loop variable does not match inferred type
found   : (number|string)
required: boolean
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOfUnknown()
func (c *TypeCheckNoTranspileTest) TestForAwaitOfUnknown() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddAsyncIterable().Build()).
		AddSource(`let /** ? */ gen;
async function foo() {
  for await (const /** null */ n of gen) {
  }
}
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOf_nonAsyncIterable_loopVarInferred()
func (c *TypeCheckNoTranspileTest) TestForAwaitOf_nonAsyncIterable_loopVarInferred() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`async function f(/** !Iterable<!Promise<string>> */ o) {
  for await (const s of o) {
    const /** number */ n = s;
  }
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOf_asyncIterable_loopVarInferred()
func (c *TypeCheckNoTranspileTest) TestForAwaitOf_asyncIterable_loopVarInferred() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`async function f(/** !AsyncIterable<string> */ o) {
  for await (const s of o) {
    const /** number */ n = s;
  }
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testForAwaitOf_unionOfIterableAndAsyncIterable_loopVarInferred()
func (c *TypeCheckNoTranspileTest) TestForAwaitOf_unionOfIterableAndAsyncIterable_loopVarInferred() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`async function f(/** !AsyncIterable<string>|!Iterable<number> */ o) {
  for await (const s of o) {
    const /** null */ n = s;
  }
}
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: null
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testNoCatchBinding()
func (c *TypeCheckNoTranspileTest) TestNoCatchBinding() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("try {} catch {}").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("try {} catch {} finally {}").
		Run()
}

// Java: TypeCheckNoTranspileTest.testMethodWithAtConstructorDoesNotDeclareType_staticClassMethod()
func (c *TypeCheckNoTranspileTest) TestMethodWithAtConstructorDoesNotDeclareType_staticClassMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @constructor */
  static Bar() { }
}

var /** !Foo.Bar */ x;
`).
		AddDiagnostic(`Bad type annotation. Unknown type Foo.Bar
It's possible that 'Foo.Bar' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMethodWithAtConstructorDoesNotDeclareType_namespaceMemberMethod()
func (c *TypeCheckNoTranspileTest) TestMethodWithAtConstructorDoesNotDeclareType_namespaceMemberMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const ns = {
  /** @constructor */
  Bar() { }
};

var /** !ns.Bar */ x;
`).
		AddDiagnostic(`Bad type annotation. Unknown type ns.Bar
It's possible that 'ns.Bar' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMethodWithAtInterfaceDoesNotDeclareType()
func (c *TypeCheckNoTranspileTest) TestMethodWithAtInterfaceDoesNotDeclareType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @interface */
  static Bar() { }
}

var /** !Foo.Bar */ x;
`).
		AddDiagnostic(`Bad type annotation. Unknown type Foo.Bar
It's possible that 'Foo.Bar' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testMethodWithAtRecordDoesNotDeclareType()
func (c *TypeCheckNoTranspileTest) TestMethodWithAtRecordDoesNotDeclareType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  /** @record */
  static Bar() { }
}

var /** !Foo.Bar */ x;
`).
		AddDiagnostic(`Bad type annotation. Unknown type Foo.Bar
It's possible that 'Foo.Bar' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeCheckingInsideGoogModule()
func (c *TypeCheckNoTranspileTest) TestTypeCheckingInsideGoogModule() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(CompilerTypeTestCase_DEFAULT_EXTERNS + CompilerTypeTestCase_CLOSURE_DEFS).
		AddSource(`goog.module('mod.A');
const /** number */ n = 'a string';
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testGoogModuleGet_hasTypeInferredInNestedExpression()
func (c *TypeCheckNoTranspileTest) TestGoogModuleGet_hasTypeInferredInNestedExpression() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(CompilerTypeTestCase_DEFAULT_EXTERNS + CompilerTypeTestCase_CLOSURE_DEFS).
		AddSource(`function takesString(/** string */ s) {}
goog.loadModule(function(exports) {
  goog.module('a');
  exports.NUM = 0;
  return exports;
});

(function() {
  takesString(goog.module.get('a').NUM);
})();
`).
		AddDiagnostic(`actual parameter 1 of takesString does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testJsdocCanReferToGoogModuleType_withoutNamedType()
func (c *TypeCheckNoTranspileTest) TestJsdocCanReferToGoogModuleType_withoutNamedType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(CompilerTypeTestCase_DEFAULT_EXTERNS).
		AddSource(`goog.loadModule(function(exports) {
  goog.module('a');
  exports.Foo = class {};
  return exports;
});
/** @type {!a.Foo<number>} */
let x;
`).
		AddDiagnosticType(RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS).
		Run()
}

// Java: TypeCheckNoTranspileTest.testJsdocCanReferToFunctionDeclarationType_withoutNamedType()
func (c *TypeCheckNoTranspileTest) TestJsdocCanReferToFunctionDeclarationType_withoutNamedType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(CompilerTypeTestCase_DEFAULT_EXTERNS + CompilerTypeTestCase_CLOSURE_DEFS).
		AddSource(`goog.provide('a.Foo');
/** @constructor */
a.Foo = function() {};

// file2
goog.loadModule(function(exports) {
  goog.module('b.Bar');

  const Foo = goog.require('a.Foo');
  /** @constructor @extends {Foo} */
  function Bar() {}
  exports = Bar;
  return exports;
});

// file3
/** @type {!b.Bar<number>} */
let x;
`).
		AddDiagnosticType(RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeCheckingEsModule_exportSpecs()
func (c *TypeCheckNoTranspileTest) TestTypeCheckingEsModule_exportSpecs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 0; export {x};").
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeCheckingEsExportedNameDecl()
func (c *TypeCheckNoTranspileTest) TestTypeCheckingEsExportedNameDecl() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("export const /** number */ x = 'not a number';").
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeCheckingInsideEsExportDefault()
func (c *TypeCheckNoTranspileTest) TestTypeCheckingInsideEsExportDefault() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("let /** number */ x; export default (x = 'not a number');").
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeCheckingEsModule_importSpecs()
func (c *TypeCheckNoTranspileTest) TestTypeCheckingEsModule_importSpecs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("import {x} from './input0';").
		Suppress(DiagnosticGroups_MODULE_LOAD).
		Run()
}

// Java: TypeCheckNoTranspileTest.testTypeCheckingEsModule_importStar()
func (c *TypeCheckNoTranspileTest) TestTypeCheckingEsModule_importStar() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("import * as mod from './input0';").
		Suppress(DiagnosticGroups_MODULE_LOAD).
		Run()
}

// Java: TypeCheckNoTranspileTest.testExplicitUnrestrictedOverridesSuperImplicitStruct()
func (c *TypeCheckNoTranspileTest) TestExplicitUnrestrictedOverridesSuperImplicitStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class A {}
/** @unrestricted */
class B extends A {
  foo() { this.x; this.x = 0; this[0]; this[0] = 0; }
}
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckNoTranspileTest.testImplicitStructOverridesSuperExplicitUnrestricted()
func (c *TypeCheckNoTranspileTest) TestImplicitStructOverridesSuperExplicitUnrestricted() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @unrestricted */
class A {}
class B extends A {
  foo() { this.x; this.x = 0; this[0]; this[0] = 0;}
}
`).
		AddDiagnostic("Property x never defined on B").
		AddDiagnostic("Cannot add a property to a struct instance after it is constructed. (If you already declared the property, make sure to give it a type.)").
		AddDiagnostic("Cannot do '[]' access on a struct").
		AddDiagnostic("Cannot do '[]' access on a struct").
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckNoTranspileTest.testExplicitUnrestrictedOverridesSuperExplicitStruct()
func (c *TypeCheckNoTranspileTest) TestExplicitUnrestrictedOverridesSuperExplicitStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @struct */
function A() {}
/** @unrestricted */
class B extends A {
  foo() { this.x; this.x = 0; this[0]; this[0] = 0;}
}
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckNoTranspileTest.testImplicitUnrestrictedDoesNotOverridesSuperExplicitStruct()
func (c *TypeCheckNoTranspileTest) TestImplicitUnrestrictedDoesNotOverridesSuperExplicitStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @struct */
function A() {}
/** @constructor @extends {A} */
function B() {}
B.prototype.foo = function() { this.x; this.x = 0; this[0]; this[0] = 0;};
`).
		AddDiagnostic("Property x never defined on B").
		AddDiagnostic("Cannot add a property to a struct instance after it is constructed. (If you already declared the property, make sure to give it a type.)").
		AddDiagnostic("Cannot do '[]' access on a struct").
		AddDiagnostic("Cannot do '[]' access on a struct").
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckNoTranspileTest.testUnion_forwardEnumRefAndNumber()
func (c *TypeCheckNoTranspileTest) TestUnion_forwardEnumRefAndNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {Type} */
const Enum = {A: 'a'};
/** @typedef {string} */ let Type;
const /** !Enum|number */ n = null;
`).
		AddDiagnostic(`initializing variable
found   : null
required: (Enum<string>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testUnion_numberAndForwardEnumRef()
func (c *TypeCheckNoTranspileTest) TestUnion_numberAndForwardEnumRef() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {Type} */
const Enum = {A: 'a'};
/** @typedef {string} */ let Type;
const /** number|!Enum */ n = null;
`).
		AddDiagnostic(`initializing variable
found   : null
required: (Enum<string>|number)
`).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDynamicImport()
func (c *TypeCheckNoTranspileTest) TestDynamicImport() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string} */
let foo = import('./foo.js');
`).
		AddDiagnostic(`initializing variable
found   : Promise<?>
required: string
`).
		Suppress(DiagnosticGroups_MODULE_LOAD).
		Run()
}

// Java: TypeCheckNoTranspileTest.testDynamicImportSpecifier()
func (c *TypeCheckNoTranspileTest) TestDynamicImportSpecifier() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const bar = null;
import(bar);
`).
		AddDiagnostic(`dynamic import specifier
found   : null
required: string
`).
		Run()
}

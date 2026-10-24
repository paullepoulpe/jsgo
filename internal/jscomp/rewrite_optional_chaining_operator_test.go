// Port of com/google/javascript/jscomp/RewriteOptionalChainingOperatorTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RewriteOptionalChainingOperatorTest.java)

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"fmt"
	"strconv"
	"testing"
)

// Test cases for transpilation pass that replaces the optional chaining
// operator (`?.`).
//
// Go-only: language-adapter — the JUnit Enclosed class RewriteOptionalChainingOperatorTest
// becomes one Go test per nested class; the @Parameterized BaseTestClass runs its test
// methods once per case, in a subtest named like the Java @Parameters name "{0} #{index}".

// Declares a class and variables to make construction of optional chains for
// test cases convenient.
//
// Java: RewriteOptionalChainingOperatorTest.TEST_BASE_EXTERNS
const RewriteOptionalChainingOperatorTest_TEST_BASE_EXTERNS = `class TestObject {
  constructor() {
    /** @const {!TestObject} */
    this.obj = this;
    /** @const {!Array<!TestObject>} */
    this.ary = [this];
    /** @const {function(number): !TestObject} */
    this.fun = (num) => this.ary[num];
    /** @const */
    this.num = 0;
  }
  /** @return {!TestObject} */
  getObj() { return this; }
  /** @return {!Array<!TestObject>} */
  getArr() { return this.ary; }
  /** @return {function(number): !TestObject} */
  getFun() { return this.fun; }
  /** @return {number} */
  getNum() { return 0; }
}

const obj = new TestObject();
const ary = obj.ary;
const fun = obj.fun;

/** @return {!TestObject} */
function getObj() {
  return obj;
}

/** @return {!Array<!TestObject>} */
function getAry() {
  return ary;
}

/** @return {function(number): !TestObject} */
function getFun() {
  return fun;
}
`

// RewriteOptionalChainingOperatorTest_BaseTestClass is
// RewriteOptionalChainingOperatorTest.BaseTestClass (@Parameterized).
type RewriteOptionalChainingOperatorTest_BaseTestClass struct {
	CompilerTestCaseBase

	// Java: @Parameter(0)
	JsSrc string

	// Java: @Parameter(1)
	JsOutput string
}

// Java: BaseTestClass()
func NewRewriteOptionalChainingOperatorTest_BaseTestClass() *RewriteOptionalChainingOperatorTest_BaseTestClass {
	c := &RewriteOptionalChainingOperatorTest_BaseTestClass{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteOptionalChainingOperatorTest_BaseTestClass(t *testing.T) {
	for index, kase := range RewriteOptionalChainingOperatorTest_BaseTestClass_Cases() {
		t.Run(fmt.Sprintf("%s #%d", kase[0], index), func(t *testing.T) {
			RunJUnit4(t, func() JUnitTestCase {
				c := NewRewriteOptionalChainingOperatorTest_BaseTestClass()
				c.JsSrc = kase[0]
				c.JsOutput = kase[1]
				return c
			})
		})
	}
}

// Java: BaseTestClass.cases()
func RewriteOptionalChainingOperatorTest_BaseTestClass_Cases() [][2]string {
	return [][2]string{

		{
			// Do rewriting within a function.
			// This will fail if the AST change is reported for the script's scope instead of
			// the function's scope.
			`function foo() {
  return obj?.num;
}
`,
			`function foo() {
  let tmp0;
  return (tmp0 = obj) == null ? void 0 : tmp0.num;
}
`,
		},
		{
			"eval?.('foo()');",
			"let tmp0;\n(tmp0 = eval) == null\n    ? void 0\n// The spec says that `eval?.()` must behave like an indirect\n// eval, so it is important that `eval?.()` not be transpiled to\n// anything that ends up containing `eval()`.\n// We must be sure to call it using the temporary variable.\n    : tmp0('foo()');\n",
		},
		{
			"obj?.ary[getNum()].obj.obj?.obj.ary",
			`let tmp0;
let tmp1;
(tmp0 = obj) == null
    ? void 0
    : (tmp1 = tmp0.ary[getNum()].obj.obj) == null
        ? void 0
        : tmp1.obj.ary
`,
		},
		{
			"(obj?.ary[getNum()]).obj.ary",
			`let tmp0;
((tmp0 = obj) == null
    ? void 0
    : tmp0.ary[getNum()]).obj.ary
`,
		},
		{
			"obj?.obj.obj?.fun(obj.getNum())",
			`let tmp0;
let tmp1;
(tmp0 = obj) == null
    ? void 0
    : (tmp1 = tmp0.obj.obj) == null
        ? void 0
        : tmp1.fun(obj.getNum())
`,
		},
		{
			"obj.ary?.[num].fun(obj.getNum?.())",
			`let tmp0;
let tmp1;
let tmp2;
(tmp2 = obj.ary) == null
    ? void 0
    : tmp2[num].fun(
        (tmp1 = (tmp0 = obj).getNum) == null
            ? void 0
            : tmp1.call(tmp0))
`,
		},
		{
			"obj?.obj",
			`let tmp0;
(tmp0 = obj) == null ? void 0 : tmp0.obj
`,
		},
		{
			"ary?.[num]",
			`let tmp0
(tmp0 = ary) == null ? void 0 : tmp0[num]
`,
		},
		{
			"obj.getObj?.()",
			`let tmp0;
let tmp1;
(tmp1 = (tmp0 = obj).getObj) == null
    ? void 0
    : tmp1.call(tmp0)
`,
		},
		{
			"getObj().getObj?.()",
			`let tmp0;
let tmp1;
(tmp1 = (tmp0 = getObj()).getObj) == null
    ? void 0
    : tmp1.call(tmp0)
`,
		},
		{
			"(getObj()?.getObj)()",
			`let tmp0;
let tmp1;
((tmp0 = getObj()) == null
    ? void 0
// Ideally we wouldn't generate a temporary to hold
// a temporary we already generated, but our logic
// is simpler if we don't worry about it.
// We will rely on optimizations to clean this up
// after transpilation.
    : (tmp1 = tmp0).getObj).call(tmp1)
`,
		},
		{
			"getAry()?.[num]",
			`let tmp0;
(tmp0 = getAry()) == null
    ? void 0
    : tmp0[num]
`,
		},
		{
			"getFun()?.(num)",
			`let tmp0;
(tmp0 = getFun()) == null
    ? void 0
    : tmp0(num)
`,
		},
		{
			"fun?.(obj?.getNum())",
			`let tmp0;
let tmp1;
(tmp1 = fun) == null
    ? void 0
    : tmp1(
          (tmp0 = obj) == null
              ? void 0
              : tmp0.getNum())
`,
		},
		{
			"obj?.fun(obj?.getNum())",
			`let tmp0;
let tmp1;
(tmp1 = obj) == null
    ? void 0
    : tmp1.fun(
        (tmp0 = obj) == null
            ? void 0
            : tmp0.getNum())
`,
		},
		{
			"while(obj = ary?.[obj?.getNum()]) {}",
			`let tmp0;
let tmp1;
for(;
    obj =
        (tmp1 = ary) == null
            ? void 0
            : tmp1[
                (tmp0 = obj) == null
                    ? void 0
                    : tmp0.getNum()
            ];) {
}
`,
		},
		{
			"let a = fun?.(num).obj.ary[obj?.getNum()]",
			`let tmp0;
let tmp1;
let a =
    (tmp1 = fun) == null
        ? void 0
        : tmp1(num).obj.ary[
            (tmp0 = obj) == null
                ? void 0
                : tmp0.getNum()
            ]
`,
		},
		{
			"() => {return foo(a?.b)}",
			"() => {let tmp0; return foo((tmp0 = a) == null ? void 0 : tmp0.b);}",
		},
		{
			"() => foo(a?.b)", //
			"() => { let tmp0; return foo((tmp0 = a) == null ? void 0 : tmp0.b);}",
		},
		{
			"(p = a?.b) => p", //
			"let tmp0; (p = (tmp0 = a) == null ? void 0 : tmp0.b) => { return p;}",
		},
		{
			"(p = a?.b?.c) => p", //
			`let tmp0;
let tmp1;
(p = (tmp0 = a) == null ? void 0 : (tmp1 = tmp0.b) == null ? void 0 : tmp1.c) =>
{ return p;}
`,
		},
		{
			`const a = { b: [3] };
label: for (const val of a?.b) {
  if (val != 3) {
    continue label;
  }
}
`,
			`const a = {b:[3]};
let tmp0;
label: for (const val of (tmp0 = a) == null ? void 0 : tmp0.b) {
  if (val != 3) {
  continue label;
  }
}
`,
		},
		{
			`{
  const x = 1;
  label: for (const a of b?.c) {}
}
`,
			`{
  const x = 1;
  let tmp0;
  label: for (const a of (tmp0 = b) == null ? void 0 : tmp0.c) {}
}
`,
		},
		{
			`const myObj = {
  getFun() {
    return { x: 3 };
  }
};
(0, myObj.getFun)?.().x;
`,
			`const myObj = {
  getFun() {
    return { x: 3 };
  }
};
let tmp0;
(tmp0 = myObj.getFun) == null ? void 0 : tmp0().x;
`,
		},
		{
			`const myObj = {
  getObj() {
    return this;
  }
};
(0, myObj?.getObj)();
`,
			`const myObj = {
  getObj() {
    return this;
  }
};
let tmp0;
((tmp0 = myObj) == null ? void 0 : tmp0.getObj)();
`,
		},
	}
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *RewriteOptionalChainingOperatorTest_BaseTestClass) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: BaseTestClass.customSetUp()
func (c *RewriteOptionalChainingOperatorTest_BaseTestClass) CustomSetUp() {
	c.EnableNormalize()
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: BaseTestClass.getProcessor(Compiler)
func (c *RewriteOptionalChainingOperatorTest_BaseTestClass) GetProcessor(compiler *Compiler) CompilerPass {
	// Just name temporary variables "tmp0", "tmp1", etc. to make the tests clearer.
	counter := 0
	testVarNameCreator := OptionalChainRewriter_TmpVarNameCreatorFunc(func() string {
		name := "tmp" + strconv.Itoa(counter)
		counter++
		return name
	})
	return NewRewriteOptionalChainingOperatorWithTmpVarNameCreator(compiler, testVarNameCreator)
}

// Java: BaseTestClass.doTest()
func (c *RewriteOptionalChainingOperatorTest_BaseTestClass) TestDoTest() {
	c.Test(c.Externs(RewriteOptionalChainingOperatorTest_TEST_BASE_EXTERNS), c.Srcs(c.JsSrc), c.Expected(c.JsOutput))
}

// RewriteOptionalChainingOperatorTest_DeleteOptChainTests is
// RewriteOptionalChainingOperatorTest.DeleteOptChainTests.
type RewriteOptionalChainingOperatorTest_DeleteOptChainTests struct {
	CompilerTestCaseBase
}

// Java: DeleteOptChainTests()
func NewRewriteOptionalChainingOperatorTest_DeleteOptChainTests() *RewriteOptionalChainingOperatorTest_DeleteOptChainTests {
	c := &RewriteOptionalChainingOperatorTest_DeleteOptChainTests{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteOptionalChainingOperatorTest_DeleteOptChainTests(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteOptionalChainingOperatorTest_DeleteOptChainTests() })
}

// Java: DeleteOptChainTests.setUp()
func (c *RewriteOptionalChainingOperatorTest_DeleteOptChainTests) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: DeleteOptChainTests.getProcessor(Compiler)
func (c *RewriteOptionalChainingOperatorTest_DeleteOptChainTests) GetProcessor(compiler *Compiler) CompilerPass {
	// Just name temporary variables "tmp0", "tmp1", etc. to make the tests clearer.
	counter := 0
	testVarNameCreator := OptionalChainRewriter_TmpVarNameCreatorFunc(func() string {
		name := "tmp" + strconv.Itoa(counter)
		counter++
		return name
	})
	return NewRewriteOptionalChainingOperatorWithTmpVarNameCreator(compiler, testVarNameCreator)
}

// Java: DeleteOptChainTests.testDeleteOptChainGetProp()
func (c *RewriteOptionalChainingOperatorTest_DeleteOptChainTests) TestDeleteOptChainGetProp() {
	c.Test(
		c.Externs(RewriteOptionalChainingOperatorTest_TEST_BASE_EXTERNS),
		c.Srcs("delete obj?.num;"),
		c.Expected(
			`let tmp0;
(tmp0 = obj) == null ? true : delete tmp0.num;
`))
}

// Java: DeleteOptChainTests.testDeleteOptChainGetProp2()
func (c *RewriteOptionalChainingOperatorTest_DeleteOptChainTests) TestDeleteOptChainGetProp2() {
	c.Test(
		c.Externs(RewriteOptionalChainingOperatorTest_TEST_BASE_EXTERNS),
		c.Srcs("delete this?.obj?.num;"),
		c.Expected(
			`let tmp0;
let tmp1
(tmp0 = this) == null ? true : (tmp1 = tmp0.obj) == null ? true : delete tmp1.num;
`))
}

// Java: DeleteOptChainTests.testDeleteOptChainGetProp3()
func (c *RewriteOptionalChainingOperatorTest_DeleteOptChainTests) TestDeleteOptChainGetProp3() {
	c.Test(
		c.Externs(RewriteOptionalChainingOperatorTest_TEST_BASE_EXTERNS),
		c.Srcs(
			// get the num-th !TestObject inside `ary`, and delete its `num` prop
			"delete getFun()?.(num).num"),
		c.Expected(
			`let tmp0;
(tmp0 = getFun()) == null
    ? true
    : delete tmp0(num).num
`))
}

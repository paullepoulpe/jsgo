// Port of com/google/javascript/jscomp/Es7RewriteExponentialOperatorTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/Es7RewriteExponentialOperatorTest.java)

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

import "testing"

// Es7RewriteExponentialOperatorTest has the test cases for the transpilation
// pass that replaces the exponential operator (`**`).
type Es7RewriteExponentialOperatorTest struct {
	CompilerTestCaseBase
}

// Java: Es7RewriteExponentialOperatorTest()
func NewEs7RewriteExponentialOperatorTest() *Es7RewriteExponentialOperatorTest {
	c := &Es7RewriteExponentialOperatorTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, NewTestExternsBuilder().AddMath().Build())
	return c
}

func TestEs7RewriteExponentialOperatorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs7RewriteExponentialOperatorTest() })
}

// Java: Es7RewriteExponentialOperatorTest.setUp()
func (c *Es7RewriteExponentialOperatorTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableNormalize()
	c.EnableMultistageCompilation()
}

// Java: Es7RewriteExponentialOperatorTest.getProcessor(Compiler)
func (c *Es7RewriteExponentialOperatorTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs7RewriteExponentialOperator(compiler)
}

// Java: Es7RewriteExponentialOperatorTest.testExponentiationOperator()
func (c *Es7RewriteExponentialOperatorTest) TestExponentiationOperator() {
	c.Test(c.Srcs("2 ** 2"), c.Expected("Math.pow(2, 2)"))
}

// Java: Es7RewriteExponentialOperatorTest.testExponentiationAssignmentOperatorSimple()
func (c *Es7RewriteExponentialOperatorTest) TestExponentiationAssignmentOperatorSimple() {
	c.Test(c.Srcs("x **= 2;"), c.Expected("x = Math.pow(x, 2)"))
}

// Java: Es7RewriteExponentialOperatorTest.testExponentiationAssignmentOperatorPropertyReference()
func (c *Es7RewriteExponentialOperatorTest) TestExponentiationAssignmentOperatorPropertyReference() {
	c.Test(
		c.Srcs("a.x **= 2;"),
		c.Expected(
			`let $jscomp$exp$assign$tmpm1146332801$0;
($jscomp$exp$assign$tmpm1146332801$0 = a).x
 = Math.pow($jscomp$exp$assign$tmpm1146332801$0.x, 2);
`))
}

// Java: Es7RewriteExponentialOperatorTest.testExponentiationAssignmentOperatorPropertyReferenceElement()
func (c *Es7RewriteExponentialOperatorTest) TestExponentiationAssignmentOperatorPropertyReferenceElement() {
	c.Test(
		c.Srcs("a[x] **= 2;"),
		c.Expected(
			`let $jscomp$exp$assign$tmpm1146332801$0;
let $jscomp$exp$assign$tmpindexm1146332801$0;
($jscomp$exp$assign$tmpm1146332801$0 = a)
[$jscomp$exp$assign$tmpindexm1146332801$0 = x]
 = Math.pow($jscomp$exp$assign$tmpm1146332801$0
[$jscomp$exp$assign$tmpindexm1146332801$0], 2);
`))
}

// @see <a href="https://github.com/google/closure-compiler/issues/2821">Issue 2821</a>
// Java: Es7RewriteExponentialOperatorTest.testExponentialOperatorInIfCondition()
func (c *Es7RewriteExponentialOperatorTest) TestExponentialOperatorInIfCondition() {
	c.Test(c.Srcs("if (2 ** 3 > 0) { }"), c.Expected("if (Math.pow(2, 3) > 0) { }"))
}

// @see <a href="https://github.com/google/closure-compiler/issues/3684">Issue 3684</a>
// Java: Es7RewriteExponentialOperatorTest.testExponentialOperatorWithBigIntCall()
func (c *Es7RewriteExponentialOperatorTest) TestExponentialOperatorWithBigIntCall() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	// Avoid transpiling to Math.pow() when BigInt is involved.
	c.TestError(c.Srcs("let a = BigInt(2) ** BigInt(3);"), Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT)
}

// @see <a href="https://github.com/google/closure-compiler/issues/3684">Issue 3684</a>
// Java: Es7RewriteExponentialOperatorTest.testAssignExponentialOperatorWithBigIntCall()
func (c *Es7RewriteExponentialOperatorTest) TestAssignExponentialOperatorWithBigIntCall() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	// Avoid transpiling to Math.pow() when BigInt is involved.
	c.TestError(c.Srcs("let a = 3n; a **= BigInt(2);"), Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT)
}

// @see <a href="https://github.com/google/closure-compiler/issues/3684">Issue 3684</a>
// Java: Es7RewriteExponentialOperatorTest.testExponentialOperatorWithBigIntCall_noWarningWithoutTypeChecking()
func (c *Es7RewriteExponentialOperatorTest) TestExponentialOperatorWithBigIntCall_noWarningWithoutTypeChecking() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.DisableTypeCheck()
	c.DisableTypeInfoValidation()
	// Ideally, we would like to avoid transpiling to Math.pow() when BigInt is involved. However,
	// when type information is missing, the transpile succeeds and generates this runtime-erroneous
	// code for <ES2016.
	c.Test(
		c.Srcs("let a = BigInt(2) ** BigInt(3);"),
		c.Expected("let a = Math.pow(BigInt(2), BigInt(3));"))
}

// @see <a href="https://github.com/google/closure-compiler/issues/3684">Issue 3684</a>
// Java: Es7RewriteExponentialOperatorTest.testExponentialOperatorWithBigIntLiteral()
func (c *Es7RewriteExponentialOperatorTest) TestExponentialOperatorWithBigIntLiteral() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	// Avoid transpiling to Math.pow() when BigInt is involved.
	c.TestError(c.Srcs("let a = 2n ** 3n;"), Es7RewriteExponentialOperator_TRANSPILE_EXPONENT_USING_BIGINT)
}

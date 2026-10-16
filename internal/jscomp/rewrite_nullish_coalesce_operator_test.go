// Port of com/google/javascript/jscomp/RewriteNullishCoalesceOperatorTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/RewriteNullishCoalesceOperatorTest.java)

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

import "testing"

// RewriteNullishCoalesceOperatorTest has the test cases for the
// transpilation pass that replaces the nullish coalesce operator (`??`).
type RewriteNullishCoalesceOperatorTest struct {
	CompilerTestCaseBase
}

// Java: RewriteNullishCoalesceOperatorTest()
func NewRewriteNullishCoalesceOperatorTest() *RewriteNullishCoalesceOperatorTest {
	c := &RewriteNullishCoalesceOperatorTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteNullishCoalesceOperatorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteNullishCoalesceOperatorTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *RewriteNullishCoalesceOperatorTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: RewriteNullishCoalesceOperatorTest.customSetUp()
func (c *RewriteNullishCoalesceOperatorTest) CustomSetUp() {
	c.EnableNormalize()
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: RewriteNullishCoalesceOperatorTest.getProcessor(Compiler)
func (c *RewriteNullishCoalesceOperatorTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRewriteNullishCoalesceOperator(compiler)
}

// Java: RewriteNullishCoalesceOperatorTest.simple()
func (c *RewriteNullishCoalesceOperatorTest) TestSimple() {
	c.Test(
		c.Srcs("a ?? b"),
		c.Expected(
			`let $jscomp$nullish$tmp0;
($jscomp$nullish$tmp0 = a) != null ? $jscomp$nullish$tmp0 : b
`))
}

// Java: RewriteNullishCoalesceOperatorTest.withOtherOperators()
func (c *RewriteNullishCoalesceOperatorTest) TestWithOtherOperators() {
	c.Test(
		c.Srcs("(x + y) ?? (a && b)"),
		c.Expected(
			`let $jscomp$nullish$tmp0;
($jscomp$nullish$tmp0 = x + y) != null ? $jscomp$nullish$tmp0 : (a && b)
`))
}

// Java: RewriteNullishCoalesceOperatorTest.insideArrowFunctionBody()
func (c *RewriteNullishCoalesceOperatorTest) TestInsideArrowFunctionBody() {
	c.Test(
		c.Srcs("() => (x + y) ?? (a && b)"),
		c.Expected(
			`() => {
  let $jscomp$nullish$tmp0;
  return ($jscomp$nullish$tmp0 = x + y) != null ? $jscomp$nullish$tmp0 : (a && b)
}
`))
}

// Java: RewriteNullishCoalesceOperatorTest.chain()
func (c *RewriteNullishCoalesceOperatorTest) TestChain() {
	c.Test(
		c.Srcs("a ?? b ?? c"),
		c.Expected(
			`let $jscomp$nullish$tmp0; let $jscomp$nullish$tmp1;
($jscomp$nullish$tmp1 = ($jscomp$nullish$tmp0 = a) != null ? $jscomp$nullish$tmp0 : b) != null
? $jscomp$nullish$tmp1 : c;
`))
}

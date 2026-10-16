// Port of com/google/javascript/jscomp/J2clUndefinedChecksRewriterPassTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/J2clUndefinedChecksRewriterPassTest.java)

/*
 * Copyright 2025 The Closure Compiler Authors.
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

// J2clUndefinedChecksRewriterPassTest tests J2clUndefinedChecksRewriterPass.
type J2clUndefinedChecksRewriterPassTest struct {
	CompilerTestCaseBase
}

// Java: J2clUndefinedChecksRewriterPassTest()
func NewJ2clUndefinedChecksRewriterPassTest() *J2clUndefinedChecksRewriterPassTest {
	c := &J2clUndefinedChecksRewriterPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestJ2clUndefinedChecksRewriterPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewJ2clUndefinedChecksRewriterPassTest() })
}

// Java: J2clUndefinedChecksRewriterPassTest.getProcessor(Compiler)
func (c *J2clUndefinedChecksRewriterPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPeepholeOptimizationsPass(
		compiler, c.GetName(), NewJ2clUndefinedChecksRewriterPass())
}

// Java: J2clUndefinedChecksRewriterPassTest.createCompiler()
func (c *J2clUndefinedChecksRewriterPassTest) CreateCompiler() *Compiler {
	compiler := c.CompilerTestCaseBase.CreateCompiler()
	J2clSourceFileChecker_MarkToRunJ2clPasses(compiler)
	return compiler
}

// Java: J2clUndefinedChecksRewriterPassTest.optimizeIsUndefined()
func (c *J2clUndefinedChecksRewriterPassTest) TestOptimizeIsUndefined() {
	c.Test("nativebootstrap$Util$impl$isUndefined(undefined)", "true")
	c.Test("nativebootstrap$Util$impl$isUndefined(void 0)", "true")
	c.Test("nativebootstrap$Util$impl$isUndefined(null)", "null, false")
	c.Test("nativebootstrap$Util$impl$isUndefined(1)", "1, false")
	c.Test("nativebootstrap$Util$impl$isUndefined('foo')", "'foo', false")
	c.Test("nativebootstrap$Util$impl$isUndefined([])", "[], false")
	c.Test("nativebootstrap$Util$impl$isUndefined({})", "({}, false)")
	c.Test("nativebootstrap$Util$impl$isUndefined(new Error())", "new Error(), false")

	c.TestSame(
		`var x = undefined;
isUndefined(x);
`)
	c.TestSame("nativebootstrap$Util$impl$isUndefined(this)")
}

// Java: J2clUndefinedChecksRewriterPassTest.optimizeCoerceToNull()
func (c *J2clUndefinedChecksRewriterPassTest) TestOptimizeCoerceToNull() {
	c.Test("nativebootstrap$Util$impl$coerceToNull(undefined)", "null")
	c.Test("nativebootstrap$Util$impl$coerceToNull(null)", "null")
	c.Test("nativebootstrap$Util$impl$coerceToNull(1)", "1")
	c.Test("nativebootstrap$Util$impl$coerceToNull('foo')", "'foo'")
	c.Test("nativebootstrap$Util$impl$coerceToNull([])", "[]")
	c.Test("nativebootstrap$Util$impl$coerceToNull({})", "({})")
	c.Test("nativebootstrap$Util$impl$coerceToNull(new Error())", "(new Error())")

	c.TestSame(
		`var x = 1;
nativebootstrap$Util$impl$coerceToNull(x);
`)
	c.TestSame("nativebootstrap$Util$impl$coerceToNull(this)")
}

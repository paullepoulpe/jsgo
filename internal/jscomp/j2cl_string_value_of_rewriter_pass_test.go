// Port of com/google/javascript/jscomp/J2clStringValueOfRewriterPassTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/J2clStringValueOfRewriterPassTest.java)

/*
 * Copyright 2018 The Closure Compiler Authors.
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

// J2clStringValueOfRewriterPassTest tests J2clStringValueOfRewriterPass.
type J2clStringValueOfRewriterPassTest struct {
	CompilerTestCaseBase
}

// Java: J2clStringValueOfRewriterPassTest()
func NewJ2clStringValueOfRewriterPassTest() *J2clStringValueOfRewriterPassTest {
	c := &J2clStringValueOfRewriterPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestJ2clStringValueOfRewriterPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewJ2clStringValueOfRewriterPassTest() })
}

// Java: J2clStringValueOfRewriterPassTest.getProcessor(Compiler)
func (c *J2clStringValueOfRewriterPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPeepholeOptimizationsPass(compiler, c.GetName(), NewJ2clStringValueOfRewriterPass())
}

// Java: J2clStringValueOfRewriterPassTest.createCompiler()
func (c *J2clStringValueOfRewriterPassTest) CreateCompiler() *Compiler {
	compiler := c.CompilerTestCaseBase.CreateCompiler()
	J2clSourceFileChecker_MarkToRunJ2clPasses(compiler)
	return compiler
}

// Java: J2clStringValueOfRewriterPassTest.testRemoveStringValueOf()
func (c *J2clStringValueOfRewriterPassTest) TestRemoveStringValueOf() {
	c.Test(
		"module$exports$java$lang$String$impl.m_valueOf__java_lang_Object__java_lang_String('')",
		"String('')")
	c.Test(
		"module$exports$java$lang$String$impl.m_valueOf__java_lang_Object__java_lang_String(null)",
		"String(null)")
	c.Test(
		`module$exports$java$lang$String$impl
    .m_valueOf__java_lang_Object__java_lang_String(undefined)
`,
		"'null'")
	c.Test(
		`module$exports$java$lang$String$impl
    .m_valueOf__java_lang_Object__java_lang_String(void 0)
`,
		"'null'")
	c.Test(
		`module$exports$java$lang$String$impl
   .m_valueOf__java_lang_Object__java_lang_String('foo' +
       module$exports$java$lang$String$impl
           .m_valueOf__java_lang_Object__java_lang_String(bar))
`,
		`String('foo' + module$exports$java$lang$String$impl
   .m_valueOf__java_lang_Object__java_lang_String(bar))
`)
	c.Test(
		`module$exports$java$lang$String$impl
   .m_valueOf__java_lang_Object__java_lang_String(foo + 'bar' + baz)
`,
		"String(foo + 'bar' + baz)")
	c.Test(
		`foo + module$exports$java$lang$String$impl
   .m_valueOf__java_lang_Object__java_lang_String('bar' + baz)
`,
		"foo + String('bar' + baz)")
	c.Test(
		"module$exports$java$lang$String$impl.m_valueOf__java_lang_Object__java_lang_String(1)",
		"String(1)")
	c.Test(
		`module$exports$java$lang$String$impl
    .m_valueOf__java_lang_Object__java_lang_String(null + foo)
`,
		"String(null + foo)")
	c.TestSame(
		`module$exports$java$lang$String$impl
   .m_valueOf__java_lang_Object__java_lang_String(window)
`)
	c.TestSame(
		`module$exports$java$lang$String$impl
    .m_valueOf__java_lang_Object__java_lang_String([])
`)
}

// Port of com/google/javascript/jscomp/J2clSourceFileCheckerTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/J2clSourceFileCheckerTest.java)

/*
 * Copyright 2016 The Closure Compiler Authors.
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

var j2clSourceFileCheckerTest_J2CL_SOURCE_FILE = SourceFile_FromCode(
	"jar:file://foo.js.zip!foo/bar.impl.java.js",
	"goog.provide('bar'); /** @constructor */ function bar(){};")

// J2clSourceFileCheckerTest is the tests for J2clSourceFileChecker.
type J2clSourceFileCheckerTest struct {
	CompilerTestCaseBase
	compiler AbstractCompiler
}

// Java: J2clSourceFileCheckerTest()
func NewJ2clSourceFileCheckerTest() *J2clSourceFileCheckerTest {
	c := &J2clSourceFileCheckerTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestJ2clSourceFileCheckerTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewJ2clSourceFileCheckerTest() })
}

// Java: J2clSourceFileCheckerTest.getProcessor(Compiler)
func (c *J2clSourceFileCheckerTest) GetProcessor(compiler *Compiler) CompilerPass {
	c.compiler = compiler
	return NewJ2clSourceFileChecker(compiler)
}

// Java: J2clSourceFileCheckerTest.testShouldRunJ2clPassesWithoutJ2clSources()
func (c *J2clSourceFileCheckerTest) TestShouldRunJ2clPassesWithoutJ2clSources() {
	c.TestSame("")
	AssertThat(J2clSourceFileChecker_ShouldRunJ2clPasses(c.compiler)).IsFalse()
}

// Java: J2clSourceFileCheckerTest.testShouldRunJ2clPassesWithJ2clSources()
func (c *J2clSourceFileCheckerTest) TestShouldRunJ2clPassesWithJ2clSources() {
	c.TestSame(c.Srcs(j2clSourceFileCheckerTest_J2CL_SOURCE_FILE))
	AssertThat(J2clSourceFileChecker_ShouldRunJ2clPasses(c.compiler)).IsTrue()
}

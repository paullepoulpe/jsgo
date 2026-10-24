// Port of com/google/javascript/jscomp/WhitespaceWrapGoogModulesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/WhitespaceWrapGoogModulesTest.java)

/*
 * Copyright 2015 The Closure Compiler Authors.
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

type WhitespaceWrapGoogModulesTest struct {
	CompilerTestCaseBase
}

func NewWhitespaceWrapGoogModulesTest() *WhitespaceWrapGoogModulesTest {
	c := &WhitespaceWrapGoogModulesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestWhitespaceWrapGoogModulesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewWhitespaceWrapGoogModulesTest() })
}

// Java: WhitespaceWrapGoogModulesTest.setUp()
func (c *WhitespaceWrapGoogModulesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()

	c.DisableCompareAsTree()
	// otherwise "use strict" in the expected output moves,
	// from where it should be (deliberately to match ClosureBundler),
	// to the top of the AST and breaks the comparison.
}

// Java: WhitespaceWrapGoogModulesTest.getProcessor(Compiler)
func (c *WhitespaceWrapGoogModulesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewWhitespaceWrapGoogModules(compiler)
}

// Java: WhitespaceWrapGoogModulesTest.testGoogModuleRewrite()
func (c *WhitespaceWrapGoogModulesTest) TestGoogModuleRewrite() {
	c.Test(
		`goog.module('test');
var f = 5;
exports = f;
`,
		"goog.loadModule(function(exports){\"use strict\";goog.module(\"test\");var f=5;exports=f;return exports})")
}

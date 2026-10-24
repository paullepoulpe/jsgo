// Port of com/google/javascript/jscomp/Es6NormalizeShorthandPropertiesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/Es6NormalizeShorthandPropertiesTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6NormalizeShorthandPropertiesTest holds test cases for the
// transpilation pass that eliminates shorthand property declarations.
//
// A shorthand property is where another identifier (variable) is used to
// define both the name *and* value of the property.
type Es6NormalizeShorthandPropertiesTest struct {
	CompilerTestCaseBase
}

// Java: Es6NormalizeShorthandPropertiesTest()
func NewEs6NormalizeShorthandPropertiesTest() *Es6NormalizeShorthandPropertiesTest {
	c := &Es6NormalizeShorthandPropertiesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestEs6NormalizeShorthandPropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6NormalizeShorthandPropertiesTest() })
}

// Java: Es6NormalizeShorthandPropertiesTest.setUp()
func (c *Es6NormalizeShorthandPropertiesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2016)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)

	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
	c.EnableMultistageCompilation()
}

// Java: Es6NormalizeShorthandPropertiesTest.getProcessor(Compiler)
func (c *Es6NormalizeShorthandPropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return PeepholeTranspilationsPass_Create(
		compiler, []AbstractPeepholeTranspilation{NewEs6NormalizeShorthandProperties(compiler)})
}

// Java: Es6NormalizeShorthandPropertiesTest.testNormalizationInSource()
func (c *Es6NormalizeShorthandPropertiesTest) TestNormalizationInSource() {
	c.Test(
		c.Externs(""),
		// Note that there's no *structural* AST difference here for the validator to notice.
		c.Srcs("const x = 5; const y = {x};"),
		c.Expected("const x = 5; const y = {x: x};"),
		c.Postcondition(
			func(compiler *Compiler) {
				NodeUtil_VisitPreOrder(
					compiler.GetRoot(),
					NodeUtil_VisitorFunc(Es6NormalizeShorthandPropertiesTest_AssertNotShorthandProperty))
			}))
}

// Java: Es6NormalizeShorthandPropertiesTest.assertNotShorthandProperty(Node)
func Es6NormalizeShorthandPropertiesTest_AssertNotShorthandProperty(node rhino.Node) {
	AssertWithMessage("Detected shorthand property node <%s>.", node).
		That(node.IsShorthandProperty()).
		IsFalse()
}

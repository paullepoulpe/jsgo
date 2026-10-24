// Port of com/google/javascript/jscomp/RewriteCatchWithNoBindingTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RewriteCatchWithNoBindingTest.java)

/*
 * Copyright 2019 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

type RewriteCatchWithNoBindingTest struct {
	CompilerTestCaseBase
}

// Java: RewriteCatchWithNoBindingTest()
func NewRewriteCatchWithNoBindingTest() *RewriteCatchWithNoBindingTest {
	c := &RewriteCatchWithNoBindingTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteCatchWithNoBindingTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteCatchWithNoBindingTest() })
}

// Java: RewriteCatchWithNoBindingTest.SPECIAL_VARIABLE_MAP
var RewriteCatchWithNoBindingTest_SPECIAL_VARIABLE_MAP = func() *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	m.Put("UNUSED_CATCH", "$jscomp$unused$catch$")
	return m
}()

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *RewriteCatchWithNoBindingTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: RewriteCatchWithNoBindingTest.customSetUp()
func (c *RewriteCatchWithNoBindingTest) CustomSetUp() {
	c.EnableNormalize()
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
	c.SetGenericNameReplacements(RewriteCatchWithNoBindingTest_SPECIAL_VARIABLE_MAP)
}

// Java: RewriteCatchWithNoBindingTest.getProcessor(Compiler)
func (c *RewriteCatchWithNoBindingTest) GetProcessor(compiler *Compiler) CompilerPass {
	return PeepholeTranspilationsPass_Create(
		compiler, []AbstractPeepholeTranspilation{NewRewriteCatchWithNoBinding(compiler)})
}

// Java: RewriteCatchWithNoBindingTest.getOptions()
func (c *RewriteCatchWithNoBindingTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2018)
	return options
}

// Java: RewriteCatchWithNoBindingTest.transpileCatchWithoutBinding()
func (c *RewriteCatchWithNoBindingTest) TestTranspileCatchWithoutBinding() {
	c.Test(
		`try {
  stuff();
} catch {
  onError();
}
`,
		`try {
  stuff();
} catch (UNUSED_CATCH$0) {
  onError();
}
`)
	AssertThat(c.GetLastCompiler().GetAllowableFeatures().ContainsFeature(parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING)).
		IsFalse()
}

// Java: RewriteCatchWithNoBindingTest.transpileCatchWithNoBindingNested()
func (c *RewriteCatchWithNoBindingTest) TestTranspileCatchWithNoBindingNested() {
	c.Test(
		`try {
  stuff();
} catch {
  try {
    onError();
  } catch {
    shruggie();
  }
}
`,
		`try {
  stuff();
} catch (UNUSED_CATCH$1) {
  try {
    onError();
  } catch (UNUSED_CATCH$0) {
    shruggie();
  }
}
`)
}

// Java: RewriteCatchWithNoBindingTest.typeOfAddedBindingIsUnknown()
func (c *RewriteCatchWithNoBindingTest) TestTypeOfAddedBindingIsUnknown() {
	c.Test(
		`try {
  stuff();
} catch {
  onError();
}
`,
		`try {
  stuff();
} catch (UNUSED_CATCH$0) {
  onError();
}
`)

	binding := c.GetLastCompiler().
		GetRoot().
		GetSecondChild(). // ROOT.
		GetFirstChild().  // SCRIPT.
		GetFirstChild().  // TRY.
		GetSecondChild(). // BLOCK.
		GetFirstChild().  // CATCH.
		GetFirstChild()   // NAME

	AssertNode(binding).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
}

// Java: RewriteCatchWithNoBindingTest.noTranspileCatchWithBinding()
func (c *RewriteCatchWithNoBindingTest) TestNoTranspileCatchWithBinding() {
	c.TestSame(
		`try {
  stuff();
} catch (err) {
  onError(err);
}
`)
}

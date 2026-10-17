// Port of com/google/javascript/jscomp/Es6TemplateLiteralsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/Es6TemplateLiteralsTest.java)

/*
 * Copyright 2026 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type Es6TemplateLiteralsTest struct {
	CompilerTestCaseBase
}

// Java: Es6TemplateLiteralsTest.REPLACEMENT_PREFIXES
var Es6TemplateLiteralsTest_REPLACEMENT_PREFIXES = func() *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	m.Put("TAGGED_TEMPLATE_TMP_VAR", "$jscomp$templatelit$")
	return m
}()

// Java: Es6TemplateLiteralsTest.RUNTIME_STUBS
const Es6TemplateLiteralsTest_RUNTIME_STUBS = `/** @const */
var $jscomp = {};
$jscomp.createTemplateTagFirstArg = function(arrayStrings) {};
$jscomp.createTemplateTagFirstArgWithRaw = function(anotherArray, rawArrayStrings) {};
`

// Java: Es6TemplateLiteralsTest()
func NewEs6TemplateLiteralsTest() *Es6TemplateLiteralsTest {
	c := &Es6TemplateLiteralsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_MINIMAL_EXTERNS+Es6TemplateLiteralsTest_RUNTIME_STUBS)
	return c
}

func TestEs6TemplateLiteralsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6TemplateLiteralsTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *Es6TemplateLiteralsTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: Es6TemplateLiteralsTest.customSetUp()
func (c *Es6TemplateLiteralsTest) CustomSetUp() {
	c.EnableNormalize()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
	c.SetGenericNameReplacements(Es6TemplateLiteralsTest_REPLACEMENT_PREFIXES)
}

// Java: Es6TemplateLiteralsTest.getProcessor(Compiler)
func (c *Es6TemplateLiteralsTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"injectTranspilationRuntimeLibraries", func(abstractCompiler AbstractCompiler) CompilerPass {
				return NewInjectTranspilationRuntimeLibraries(abstractCompiler)
			}))
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("lateEs6ToEs3Converter", func(abstractCompiler AbstractCompiler) CompilerPass {
		return NewLateEs6ToEs3Converter(abstractCompiler)
	}))
	return optimizer
}

// Java: Es6TemplateLiteralsTest.testTaggedTemplateLiteral_nameTag()
func (c *Es6TemplateLiteralsTest) TestTaggedTemplateLiteral_nameTag() {
	c.Test(
		"tag`hello`",
		`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['hello']);
tag(TAGGED_TEMPLATE_TMP_VAR$0);
`)
}

// Java: Es6TemplateLiteralsTest.testTaggedTemplateLiteral_getPropTag()
func (c *Es6TemplateLiteralsTest) TestTaggedTemplateLiteral_getPropTag() {
	c.Test(
		c.Externs(Es6TemplateLiteralsTest_RUNTIME_STUBS, "var obj = {}; obj.tag;"),
		c.Srcs("obj.tag`hello`"),
		c.Expected(
			`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['hello']);
obj.tag(TAGGED_TEMPLATE_TMP_VAR$0);
`))
}

// Java: Es6TemplateLiteralsTest.testTaggedTemplateLiteral_getElemTag()
func (c *Es6TemplateLiteralsTest) TestTaggedTemplateLiteral_getElemTag() {
	c.Test(
		c.Externs(Es6TemplateLiteralsTest_RUNTIME_STUBS, "var obj = {}; var tag; obj[tag];"),
		c.Srcs("obj[tag]`hello`"),
		c.Expected(
			`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['hello']);
obj[tag](TAGGED_TEMPLATE_TMP_VAR$0);
`))
}

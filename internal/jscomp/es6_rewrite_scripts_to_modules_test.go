// Port of com/google/javascript/jscomp/Es6RewriteScriptsToModulesTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/Es6RewriteScriptsToModulesTest.java)

/*
 * Copyright 2017 The Closure Compiler Authors.
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

type Es6RewriteScriptsToModulesTest struct {
	CompilerTestCaseBase
}

func NewEs6RewriteScriptsToModulesTest() *Es6RewriteScriptsToModulesTest {
	c := &Es6RewriteScriptsToModulesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestEs6RewriteScriptsToModulesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteScriptsToModulesTest() })
}

// Java: Es6RewriteScriptsToModulesTest.getProcessor(Compiler)
func (c *Es6RewriteScriptsToModulesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs6RewriteScriptsToModules(compiler)
}

// Java: Es6RewriteScriptsToModulesTest.getOptions()
func (c *Es6RewriteScriptsToModulesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	return options
}

// Java: Es6RewriteScriptsToModulesTest.testImportedScript()
func (c *Es6RewriteScriptsToModulesTest) TestImportedScript() {
	compiler := c.CreateCompiler()
	options := c.GetOptions()
	compiler.Init(
		[]*SourceFile{},
		[]*SourceFile{
			SourceFile_FromCode("/script.js", ""),
			SourceFile_FromCode("/module.js", "import '/script.js';")},
		options)
	root := compiler.ParseInputs()
	AssertThat(compiler.GetErrors()).IsEmpty()
	externsRoot := root.GetFirstChild()
	mainRoot := externsRoot.GetNext()
	c.GetProcessor(compiler).Process(externsRoot, mainRoot)
	AssertThat(mainRoot.GetFirstFirstChild().IsModuleBody()).IsTrue()
}

// Java: Es6RewriteScriptsToModulesTest.testNonImportedScript()
func (c *Es6RewriteScriptsToModulesTest) TestNonImportedScript() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode("/script.js", ""),
			SourceFile_FromCode("/module.js", "export default 'foo';")))
}

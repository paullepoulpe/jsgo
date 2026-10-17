// Port of com/google/javascript/jscomp/J2clUtilGetDefineRewriterPassTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/J2clUtilGetDefineRewriterPassTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clUtilGetDefineRewriterPassTest is the tests for J2clUtilGetDefineRewriterPass.
type J2clUtilGetDefineRewriterPassTest struct {
	CompilerTestCaseBase
}

// Java: J2clUtilGetDefineRewriterPassTest()
func NewJ2clUtilGetDefineRewriterPassTest() *J2clUtilGetDefineRewriterPassTest {
	c := &J2clUtilGetDefineRewriterPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestJ2clUtilGetDefineRewriterPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewJ2clUtilGetDefineRewriterPassTest() })
}

// Java: J2clUtilGetDefineRewriterPassTest.getProcessor(Compiler)
func (c *J2clUtilGetDefineRewriterPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		// We require the ProcessDefines pass to run first, so run it in optimize mode so that we
		// replicate the state we expect the code to be in when it reaches this pass.
		NewProcessDefines_Builder(compiler).SetMode(ProcessDefines_Mode_OPTIMIZE).Build().Process(externs, root)
		NewJ2clUtilGetDefineRewriterPass(compiler).Process(externs, root)
	})
}

// Java: J2clUtilGetDefineRewriterPassTest.createCompiler()
func (c *J2clUtilGetDefineRewriterPassTest) CreateCompiler() *Compiler {
	compiler := c.CompilerTestCaseBase.CreateCompiler()
	J2clSourceFileChecker_MarkToRunJ2clPasses(compiler)
	return compiler
}

// Java: J2clUtilGetDefineRewriterPassTest.testUtilGetDefine()
func (c *J2clUtilGetDefineRewriterPassTest) TestUtilGetDefine() {
	c.Test(
		`var a = {};
a.b = {}
/** @define {boolean} */ a.b.c = goog.define('a.b.c', true);
jre.addSystemPropertyFromGoogDefine('a.b.c', a.b.c);
nativebootstrap.Util.$getDefine('a.b.c', 'def');
`,
		`var a = {};
a.b = {}
/** @define {boolean} */ a.b.c = true;
var jscomp$defines$a$b$c = a.b.c;
jre.addSystemPropertyFromGoogDefine('a.b.c', a.b.c);
('def', String(jscomp$defines$a$b$c));
`)
	c.Test(
		`var a = {};
a.b = {}
/** @define {boolean} */ a.b.c = goog.define('a.b.c', true);
jre.addSystemPropertyFromGoogDefine('a.b.c', a.b.c);
nativebootstrap.Util.$getDefine('a.b.c');
`,
		`var a = {};
a.b = {}
/** @define {boolean} */ a.b.c = true;
var jscomp$defines$a$b$c = a.b.c;
jre.addSystemPropertyFromGoogDefine('a.b.c', a.b.c);
(null, String(jscomp$defines$a$b$c));
`)
	c.Test(
		`/** @define {boolean} */ var x = goog.define('x', 1);
jre.addSystemPropertyFromGoogDefine('x', x);
/** @define {boolean} */ var y = goog.define('y', x);
jre.addSystemPropertyFromGoogDefine('y', y);
nativebootstrap.Util.$getDefine('x');
nativebootstrap.Util.$getDefine('y');
`,
		`/** @define {boolean} */ var x = 1;
var jscomp$defines$x = x;
jre.addSystemPropertyFromGoogDefine('x', x);
/** @define {boolean} */ var y = x;
var jscomp$defines$y = y;
jre.addSystemPropertyFromGoogDefine('y', y);
(null, String(jscomp$defines$x));
(null, String(jscomp$defines$y));
`)
	c.Test(
		`jre.addSystemPropertyFromGoogDefine('COMPILED', COMPILED);
nativebootstrap.Util.$getDefine('COMPILED');
`,
		`jre.addSystemPropertyFromGoogDefine('COMPILED', COMPILED);
(null, String(COMPILED));
`)
}

// Java: J2clUtilGetDefineRewriterPassTest.testUtilGetDefine_notDefined()
func (c *J2clUtilGetDefineRewriterPassTest) TestUtilGetDefine_notDefined() {
	c.Test("nativebootstrap.Util.$getDefine('not.defined');", "null;")
	c.Test("nativebootstrap.Util.$getDefine('not.defined', 'def');", "'def';")
}

// Java: J2clUtilGetDefineRewriterPassTest.testUtilGetDefine_missingAddSystemProperty()
func (c *J2clUtilGetDefineRewriterPassTest) TestUtilGetDefine_missingAddSystemProperty() {
	c.TestError(
		`/** @define {boolean} */ var a = goog.define('a', true);
nativebootstrap.Util.$getDefine('a');
`,
		J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY)
}

// Java: J2clUtilGetDefineRewriterPassTest.testUtilGetDefine_notAStringLiteral()
func (c *J2clUtilGetDefineRewriterPassTest) TestUtilGetDefine_notAStringLiteral() {
	c.TestError(
		`var defineName = 'a';
/** @define {boolean} */ var a = goog.define('a', true);
nativebootstrap.Util.$getDefine(defineName);
`,
		J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME)
}

// Java: J2clUtilGetDefineRewriterPassTest.testAddSystemProperty_notDefined()
func (c *J2clUtilGetDefineRewriterPassTest) TestAddSystemProperty_notDefined() {
	c.TestError(
		`var notADefine = true;
jre.addSystemPropertyFromGoogDefine('not_a_define', notADefine);
`,
		J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE)
}

// Java: J2clUtilGetDefineRewriterPassTest.testAddSystemProperty_notAStringLiteral()
func (c *J2clUtilGetDefineRewriterPassTest) TestAddSystemProperty_notAStringLiteral() {
	c.TestError(
		`var defineName = 'someDefine';
/** @define {boolean} */ var someDefine = goog.define('someDefine', true);
jre.addSystemPropertyFromGoogDefine(defineName, someDefine);
`,
		J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME)
}

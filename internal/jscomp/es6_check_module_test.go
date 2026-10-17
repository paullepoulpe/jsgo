// Port of com/google/javascript/jscomp/Es6CheckModuleTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/Es6CheckModuleTest.java)

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

// Es6CheckModuleTest holds unit tests for Es6CheckModule
type Es6CheckModuleTest struct {
	CompilerTestCaseBase
}

func NewEs6CheckModuleTest() *Es6CheckModuleTest {
	c := &Es6CheckModuleTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestEs6CheckModuleTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6CheckModuleTest() })
}

// Java: Es6CheckModuleTest.getProcessor(Compiler)
func (c *Es6CheckModuleTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs6CheckModule(compiler)
}

// GetOptions specifies EcmaScript 2015 (ES6) for tests
//
// Java: Es6CheckModuleTest.getOptions()
func (c *Es6CheckModuleTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// Java: Es6CheckModuleTest.testEs6ThisWithExportModule()
func (c *Es6CheckModuleTest) TestEs6ThisWithExportModule() {
	c.TestWarning("export {};\nfoo.call(this, 1, 2, 3);", Es6CheckModule_ES6_MODULE_REFERENCES_THIS)
}

// Java: Es6CheckModuleTest.testEs6ThisWithImportModule()
func (c *Es6CheckModuleTest) TestEs6ThisWithImportModule() {
	c.TestWarning(
		`import ln from './other/x'
if (x) {
  alert(this);
}
`,
		Es6CheckModule_ES6_MODULE_REFERENCES_THIS)
}

// Java: Es6CheckModuleTest.testEs6ThisWithConstructor()
func (c *Es6CheckModuleTest) TestEs6ThisWithConstructor() {
	c.TestSame(
		`class Foo {
  constructor() {
    this.x = 5;
  }
}

exports = Foo;
`)
}

// Java: Es6CheckModuleTest.testThisWithStaticMethod()
func (c *Es6CheckModuleTest) TestThisWithStaticMethod() {
	c.TestSame("class Foo { static h() {var x = this.y;} }; exports = Foo;")
	c.TestSame("class Foo {static h() {this.x = 2; }}; exports = Foo;")
	c.TestSame("class Foo {static h() {this[this.x] = 3;}}; exports = Foo;")
	c.TestSame(
		`class Foo {
  static h() {
    function g() {
      return this.f() + 1;
    }
    var y = g() + 1;
  }
  static f() {return 1;}
}
exports = Foo;
`)
	c.TestSame(
		`class Foo {
  static h() {
    button.addEventListener('click', function () {
      this.click();
    });
  }
  static click() {}
};
exports = Foo;
`)
}

// Java: Es6CheckModuleTest.testThisWithStaticBlock()
func (c *Es6CheckModuleTest) TestThisWithStaticBlock() {
	c.TestSame("class Foo { static {var x = this.y;} }; exports = Foo;")
	c.TestSame("class Foo {static {this.x = 2; }}; exports = Foo;")
	c.TestSame("class Foo {static {this[this.x] = 3;}}; exports = Foo;")
	c.TestSame(
		`class Foo {
  static {
    function g() {
      return this.f() + 1;
    }
    var y = g() + 1;
  }
  static f() {return 1;}
}
exports = Foo;
`)
	c.TestSame(
		`class Foo {
  static {
    button.addEventListener('click', function () {
      this.click();
    });
  }
  static click() {}
};
exports = Foo;
`)
}

// just here to make sure import.meta doesn't break anything
//
// Java: Es6CheckModuleTest.testImportMeta()
func (c *Es6CheckModuleTest) TestImportMeta() {
	c.TestSame(
		`class Foo {
  constructor() {
    this.url = import.meta.url
  }
}

exports = Foo;
`)
}

// Java: Es6CheckModuleTest.testCannotRenameImport()
func (c *Es6CheckModuleTest) TestCannotRenameImport() {
	c.TestError("import { p } from '/other'; p = 2;", Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED)

	c.TestError(
		"import { p } from '/other'; ({p} = {});", Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED)
	c.TestError(
		"import { p } from '/other'; ({z:p} = {});", Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED)
	c.TestSame("import { p } from '/other'; ({p:z} = {});")

	c.TestError("import { p } from '/other'; [p] = [];", Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED)

	c.TestSame("import { p } from '/other'; p.x = 2;")
	c.TestSame("import { p } from '/other'; p['x'] = 2;")

	c.TestError(
		"import Default from '/other'; Default = 2;", Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED)

	c.TestSame("import Default from '/other'; Default.x = 2;")
	c.TestSame("import Default from '/other'; Default['x'] = 2;")

	c.TestError(
		"import * as Module from '/other'; Module = 2;",
		Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED)

	c.TestError(
		"import * as Module from '/other'; Module.x = 2;",
		Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED)

	c.TestError(
		"import * as Module from '/other'; Module['x'] = 2;",
		Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED)

	c.TestSame("import * as Module from '/other'; Module.x.y = 2;")

	c.TestSame("import * as Module from '/other'; Module['x'].y = 2;")

	// Handled by VariableReferenceCheck.
	c.TestSame("import { p } from '/other'; let p = 0;")
	c.TestSame("import { p } from '/other'; var {p} = {};")
	c.TestSame("import { p } from '/other'; var [p] = [];")
	c.TestSame("import { p } from '/other'; function p() {};")
	c.TestSame("import { p } from '/other'; class p {};")
}

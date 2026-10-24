// Port of com/google/javascript/jscomp/lint/CheckExtraRequiresTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckExtraRequiresTest.java)

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

// Tests for the "extra requires" check in {@link CheckMissingAndExtraRequires}.
type CheckExtraRequiresTest struct {
	CompilerTestCaseBase
}

func NewCheckExtraRequiresTest() *CheckExtraRequiresTest {
	c := &CheckExtraRequiresTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckExtraRequiresTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckExtraRequiresTest() })
}

// Java: CheckExtraRequiresTest.getOptions()
func (c *CheckExtraRequiresTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_EXTRA_REQUIRE, CheckLevel_ERROR)
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// Java: CheckExtraRequiresTest.getProcessor(Compiler)
func (c *CheckExtraRequiresTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckExtraRequires(compiler, nil)
}

// Java: CheckExtraRequiresTest.testExtraRequire()
func (c *CheckExtraRequiresTest) TestExtraRequire() {
	c.TestError("goog.require('foo.Bar');", CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testExtraImport()
func (c *CheckExtraRequiresTest) TestExtraImport() {
	c.TestError("import z from '/x.y';", CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testFailForwardDeclareInModule()
func (c *CheckExtraRequiresTest) TestFailForwardDeclareInModule() {
	c.TestError(
		`goog.module('example');

var Event = goog.forwardDeclare('goog.events.Event');
var Unused = goog.forwardDeclare('goog.events.Unused');

/**
 * @param {!Event} event
 */
function listener(event) {
  alert(event);
}

exports = listener;
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testShadowedUnusedImport()
func (c *CheckExtraRequiresTest) TestShadowedUnusedImport() {
	// It would be nice to catch this, but currently the pass is name based and thus misses
	// the fact that the import is unused.
	c.TestSame(
		`goog.module('example');

var Shadowed = goog.forwardDeclare('foo.Shadowed');

function f(Shadowed) {
  alert(Shadowed);
}
`)
}

// Java: CheckExtraRequiresTest.testFailForwardDeclare()
func (c *CheckExtraRequiresTest) TestFailForwardDeclare() {
	c.TestError(
		`goog.forwardDeclare('goog.events.Event');
goog.forwardDeclare('goog.events.Unused');

/**
 * @param {!goog.events.Event} event
 */
function listener(event) {
  alert(event);
}
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testNoWarning_require()
func (c *CheckExtraRequiresTest) TestNoWarning_require() {
	c.TestSame("goog.require('foo.Bar'); var x = new foo.Bar();")
	c.TestSame("goog.require('foo.Bar'); let x = new foo.Bar();")
	c.TestSame("goog.require('foo.Bar'); const x = new foo.Bar();")
	c.TestSame("goog.require('foo.Bar'); /** @type {foo.Bar} */ var x;")
	c.TestSame("goog.require('foo.Bar'); /** @type {Array<foo.Bar>} */ var x;")
	c.TestSame("goog.require('foo.Bar'); var x = new foo.Bar.Baz();")
	c.TestSame("goog.require('foo.bar'); var x = foo.bar();")
	c.TestSame("goog.require('foo.bar'); var x = /** @type {foo.bar} */ (null);")
	c.TestSame("goog.require('foo.bar'); function f(/** foo.bar */ x) {}")
	c.TestSame("goog.require('foo.bar'); alert(foo.bar.baz);")
	c.TestSame("/** @suppress {extraRequire} */ goog.require('foo.bar');")
	c.TestSame("goog.require('foo.bar'); goog.scope(function() { var bar = foo.bar; alert(bar); });")
	c.TestSame("goog.require('foo'); foo();")
	c.TestSame("goog.require('foo'); new foo();")
	c.TestSame("/** @suppress {extraRequire} */ var bar = goog.require('foo.bar');")
}

// Java: CheckExtraRequiresTest.testNoWarning_requireType()
func (c *CheckExtraRequiresTest) TestNoWarning_requireType() {
	c.TestSame("goog.requireType('foo.Bar'); /** @type {foo.Bar} */ var x;")
	c.TestSame("goog.requireType('foo.Bar'); /** @type {Array<foo.Bar>} */ var x;")
	c.TestSame("goog.requireType('foo.bar'); function f(/** foo.bar */ x) {}")
	c.TestSame("/** @suppress {extraRequire} */ goog.requireType('foo.bar');")
	c.TestSame("/** @suppress {extraRequire} */ var bar = goog.requireType('foo.bar');")
}

// Java: CheckExtraRequiresTest.testNoWarning_require_externsJsDoc()
func (c *CheckExtraRequiresTest) TestNoWarning_require_externsJsDoc() {
	c.TestSame(
		c.Externs([]*SourceFile{SourceFile_FromCode("externs", "/** @const */ var ns;")}),
		c.Srcs("goog.require('ns.Foo'); /** @type {ns.Foo} */ var f;"))
}

// Java: CheckExtraRequiresTest.testNoWarning_requireType_externsJsDoc()
func (c *CheckExtraRequiresTest) TestNoWarning_requireType_externsJsDoc() {
	c.TestSame(
		c.Externs([]*SourceFile{SourceFile_FromCode("externs", "/** @const */ var ns;")}),
		c.Srcs("goog.requireType('ns.Foo'); /** @type {ns.Foo} */ var f;"))
}

// Java: CheckExtraRequiresTest.testNoWarning_require_externsNew()
func (c *CheckExtraRequiresTest) TestNoWarning_require_externsNew() {
	c.TestSame(
		c.Externs([]*SourceFile{SourceFile_FromCode("externs", "/** @const */ var ns;")}),
		c.Srcs("goog.require('ns.Foo'); new ns.Foo();"))
}

// Java: CheckExtraRequiresTest.testNoWarning_requireType_externsNew()
func (c *CheckExtraRequiresTest) TestNoWarning_requireType_externsNew() {
	c.TestSame(
		c.Externs([]*SourceFile{SourceFile_FromCode("externs", "/** @const */ var ns;")}),
		c.Srcs("goog.requireType('ns.Foo'); new ns.Foo();"))
}

// Java: CheckExtraRequiresTest.testNoWarning_esImport_objlitShorthand()
func (c *CheckExtraRequiresTest) TestNoWarning_esImport_objlitShorthand() {
	c.TestSame(
		`import '/example.module';

import X from '/example.X';
alert({X});
`)
}

// Java: CheckExtraRequiresTest.testNoWarning_require_InnerClassInExtends()
func (c *CheckExtraRequiresTest) TestNoWarning_require_InnerClassInExtends() {
	c.TestSame(
		`var goog = {};
goog.require('goog.foo.Bar');

/** @constructor @extends {goog.foo.Bar.Inner} */
function SubClass() {}
`)
}

// Java: CheckExtraRequiresTest.testNoWarning_requireType_InnerClassInExtends()
func (c *CheckExtraRequiresTest) TestNoWarning_requireType_InnerClassInExtends() {
	c.TestSame(
		`var goog = {};
goog.requireType('goog.foo.Bar');

/** @constructor @extends {goog.foo.Bar.Inner} */
function SubClass() {}
`)
}

// Java: CheckExtraRequiresTest.testWarning_require()
func (c *CheckExtraRequiresTest) TestWarning_require() {
	c.TestError("goog.require('foo.bar');", CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError(
		`goog.require('Bar');
function func( {a} ){}
func( {a: 1} );
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError(
		`goog.require('Bar');
function func( a = 1 ){}
func(42);
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testWarning_requireType()
func (c *CheckExtraRequiresTest) TestWarning_requireType() {
	c.TestError("goog.requireType('foo.bar');", CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError(
		`goog.requireType('Bar');
/** @type {string} */ var x
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testNoWarningMultipleFiles()
func (c *CheckExtraRequiresTest) TestNoWarningMultipleFiles() {
	c.TestSame(
		c.Srcs(
			"goog.require('Foo'); var foo = new Foo();",
			"goog.require('Bar'); var bar = new Bar();"))
}

// Java: CheckExtraRequiresTest.testPassModule()
func (c *CheckExtraRequiresTest) TestPassModule() {
	c.TestSame(
		`import {Foo} from '/bar';
new Foo();
`)
	c.TestSame(
		`import Bar from '/bar';
new Bar();
`)
	c.TestSame(
		`import {CoolFeature as Foo} from '/bar';
new Foo();
`)
	c.TestSame(
		`import Bar, {CoolFeature as Foo, OtherThing as Baz} from '/bar';
new Foo(); new Bar(); new Baz();
`)
}

// Java: CheckExtraRequiresTest.testFailModule()
func (c *CheckExtraRequiresTest) TestFailModule() {
	c.TestError("import {Foo} from '/bar';", CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError("import {Foo as Foo} from '/bar';", CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError("import {Foo as Bar} from '/bar';", CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError(
		`import {Foo} from '/bar';
goog.require('example.ExtraRequire');
new Foo;
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError(
		`import {Foo} from '/bar';
goog.requireType('example.ExtraRequire');
new Foo;
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testPassForwardDeclareInModule()
func (c *CheckExtraRequiresTest) TestPassForwardDeclareInModule() {
	c.TestSame(
		`goog.module('example');

var Event = goog.forwardDeclare('goog.events.Event');

/**
 * @param {!Event} event
 */
function listener(event) {
  alert(event);
}

exports = listener;
`)
}

// Java: CheckExtraRequiresTest.testPassForwardDeclare()
func (c *CheckExtraRequiresTest) TestPassForwardDeclare() {
	c.TestSame(
		`goog.forwardDeclare('goog.events.Event');

/**
 * @param {!goog.events.Event} event
 */
function listener(event) {
  alert(event);
}
`)
}

// Java: CheckExtraRequiresTest.testGoogModuleGet()
func (c *CheckExtraRequiresTest) TestGoogModuleGet() {
	c.TestSame(
		`goog.provide('x.y');
goog.require('foo.bar');

goog.scope(function() {
var bar = goog.module.get('foo.bar');
x.y = function() {};
});
`)
}

// Java: CheckExtraRequiresTest.testGoogModuleWithAliasedRequire()
func (c *CheckExtraRequiresTest) TestGoogModuleWithAliasedRequire() {
	c.TestNoWarning(
		`goog.module('example');

const asserts = goog.require('goog.asserts');

exports = function() {
  asserts.assert(true);
};
`)
	c.TestError(
		`goog.module('example');

const asserts = goog.require('goog.asserts');

exports = function() {
  goog.asserts.assert(true);
};
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testGoogModuleWithDestructuringRequire()
func (c *CheckExtraRequiresTest) TestGoogModuleWithDestructuringRequire() {
	c.TestNoWarning(
		`goog.module('example');

const {assert} = goog.require('goog.asserts');

exports = function() {
  assert(true);
};
`)
	c.TestError(
		`goog.module('example');

const {assert} = goog.require('goog.asserts');

exports = function() {
  goog.asserts.assert(true);
};
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testGoogModuleWithDestructuringShortnameRequire()
func (c *CheckExtraRequiresTest) TestGoogModuleWithDestructuringShortnameRequire() {
	c.TestNoWarning(
		`goog.module('example');

const {assert: googAssert} = goog.require('goog.asserts');

exports = function() {
  googAssert(true);
};
`)
	c.TestError(
		`goog.module('example');

var {assert: googAssert} = goog.require('goog.asserts');

exports = function() {
  assert(true);
};
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError(
		`goog.module('example');

const {assert: googAssert} = goog.require('goog.asserts');

exports = function() {
  goog.asserts.assert(true);
};
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testGoogModuleWithPartiallyUnusedDestructuringRequire()
func (c *CheckExtraRequiresTest) TestGoogModuleWithPartiallyUnusedDestructuringRequire() {
	c.TestError(
		`goog.module('example');

const {assert, fail} = goog.require('goog.asserts');

exports = function() {
  assert(true);
};
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testGoogModuleWithEmptyDestructuringRequire()
func (c *CheckExtraRequiresTest) TestGoogModuleWithEmptyDestructuringRequire() {
	c.TestError(
		`goog.module('example');

var {} = goog.require('goog.asserts');
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testGoogModuleWithAliasedRequireType()
func (c *CheckExtraRequiresTest) TestGoogModuleWithAliasedRequireType() {
	c.TestNoWarning(
		`goog.module('example');

const color = goog.requireType('goog.color');

exports = /** @param {color.Rgb} x */ function(x) { alert(x); };
`)
	c.TestError(
		`goog.module('example');

const color = goog.requireType('goog.color');

exports = /** @param {goog.color.Rgb} x */ function(x) { alert(x); };
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testGoogModuleWithDestructuringRequireType()
func (c *CheckExtraRequiresTest) TestGoogModuleWithDestructuringRequireType() {
	c.TestError(
		`goog.module('example');

const {Rgb} = goog.requireType('goog.color');

exports = /** @param {goog.color.Rgb} x */ function(x) { alert(x); };
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testGoogModuleWithDestructuringShortnameRequireType()
func (c *CheckExtraRequiresTest) TestGoogModuleWithDestructuringShortnameRequireType() {
	c.TestNoWarning(
		`goog.module('example');

const {Rgb: googColorRgb} = goog.requireType('goog.color');

exports = /** @param {googColorRgb} x */ function(x) { alert(x); };
`)
	c.TestError(
		`goog.module('example');

const {Rgb: googColorRgb} = goog.requireType('goog.color');

exports = /** @param {Rgb} x */ function(x) { alert(x); };
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError(
		`goog.module('example');

const {Rgb: googColorRgb} = goog.requireType('goog.color');

exports = /** @param {goog.color.Rgb} x */ function(x) { alert(x); };
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testGoogModuleWithPartiallyUnusedDestructuringRequireType()
func (c *CheckExtraRequiresTest) TestGoogModuleWithPartiallyUnusedDestructuringRequireType() {
	c.TestError(
		`goog.module('example');

const {Rgb, Hsv} = goog.require('goog.color');

exports = /** @param {Rgb} x */ function(x) { alert(x); };
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testGoogModuleWithEmptyDestructuringRequireType()
func (c *CheckExtraRequiresTest) TestGoogModuleWithEmptyDestructuringRequireType() {
	c.TestError(
		`goog.module('example');

var {} = goog.requireType('goog.color');
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

// Java: CheckExtraRequiresTest.testES6ModuleWithDestructuringRequire()
func (c *CheckExtraRequiresTest) TestES6ModuleWithDestructuringRequire() {
	c.TestError(
		`import '/example';

import {assert, fail} from '/goog.asserts';

export default function() {
  assert(true);
};
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError(
		`import '/example';

import {assert as assert, fail as fail} from '/goog.asserts';

export default function() {
  assert(true);
};
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
	c.TestError(
		`import '/example';

import {assert as a, fail as f} from '/goog.asserts';

export default function() {
  a(true);
};
`,
		CheckExtraRequires_EXTRA_REQUIRE_WARNING)
}

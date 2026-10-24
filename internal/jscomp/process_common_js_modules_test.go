// Port of com/google/javascript/jscomp/ProcessCommonJSModulesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ProcessCommonJSModulesTest.java)

/*
 * Copyright 2011 The Closure Compiler Authors.
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

// Unit tests for ProcessCommonJSModules
type ProcessCommonJSModulesTest struct {
	CompilerTestCaseBase
	moduleRoots    []string
	resolutionMode ModuleLoader_ResolutionMode
}

func NewProcessCommonJSModulesTest() *ProcessCommonJSModulesTest {
	c := &ProcessCommonJSModulesTest{}
	c.moduleRoots = nil
	c.resolutionMode = ModuleLoader_ResolutionMode_NODE
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestProcessCommonJSModulesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewProcessCommonJSModulesTest() })
}

// Java: ProcessCommonJSModulesTest.getOptions()
func (c *ProcessCommonJSModulesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetProcessCommonJSModules(true)
	options.SetModuleResolutionMode(c.resolutionMode)

	if c.moduleRoots != nil {
		options.SetModuleRoots(c.moduleRoots)
	}

	return options
}

// Java: ProcessCommonJSModulesTest.getProcessor(Compiler)
func (c *ProcessCommonJSModulesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewProcessCommonJSModules(compiler)
}

// Java: ProcessCommonJSModulesTest.testModules(String, String, String)
func (c *ProcessCommonJSModulesTest) TestModules(filename string, input string, expected string) {
	ModulesTestUtils_TestModules(c, filename, input, expected)
}

// Java: ProcessCommonJSModulesTest.testWithoutExports()
func (c *ProcessCommonJSModulesTest) TestWithoutExports() {
	c.TestModules(
		"test.js",
		"var name = require('./other'); name.call(null)",
		`var name = module$other.default;
module$other.default.call(null);
`)
	c.Test(
		c.Srcs(
			SourceFile_FromCode(Compiler_JoinPathParts("mod", "name.js"), "module.exports = {};"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("test", "sub.js"),
				`var name = require('../mod/name');
(function() { let foo = name; foo(); })();
`)),
		c.Expected(
			SourceFile_FromCode(
				Compiler_JoinPathParts("mod", "name.js"),
				"/** @const */ var module$mod$name = {/** @const */ default: {}};"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("test", "sub.js"),
				`var name = module$mod$name.default;
(function() { let foo = module$mod$name.default; foo(); })();
`)))
}

// Java: ProcessCommonJSModulesTest.testExports()
func (c *ProcessCommonJSModulesTest) TestExports() {
	c.TestModules(
		"test.js",
		`var name = require('./other');
exports.foo = 1;
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
var name$$module$test = module$other.default;
module$test.default.foo = 1;
`)

	c.TestModules(
		"test.js",
		`var name = require('./other');
module.exports = function() {};
`,
		`/** @const */ var module$test = {};
var name$$module$test = module$other.default;
/** @const */ module$test.default = function () {};
`)
}

// Java: ProcessCommonJSModulesTest.testExportsInExpression()
func (c *ProcessCommonJSModulesTest) TestExportsInExpression() {
	c.TestModules(
		"test.js",
		`var name = require('./other');
var e;
e = module.exports = function() {};
`,
		`/** @const */ var module$test = {};
var name$$module$test = module$other.default;
var e$$module$test;
e$$module$test = /** @const */ module$test.default = function () {};
`)

	c.TestModules(
		"test.js",
		`var name = require('./other');
var e = module.exports = function() {};
`,
		`/** @const */ var module$test = {};
var name$$module$test = module$other.default;
var e$$module$test = /** @const */ module$test.default = function () {};
`)

	c.TestModules(
		"test.js",
		`var name = require('./other');
(module.exports = function() {})();
`,
		`/** @const */ var module$test = {};
var name$$module$test = module$other.default;
(/** @const */ module$test.default = function () {})();
`)
}

// Java: ProcessCommonJSModulesTest.testPropertyExports()
func (c *ProcessCommonJSModulesTest) TestPropertyExports() {
	c.TestModules(
		"test.js",
		"exports.one = 1; module.exports.obj = {}; module.exports.obj.two = 2;",
		`/** @const */ var module$test = {default: {}};
module$test.default.one = 1;
module$test.default.obj = {};
module$test.default.obj.two = 2;
`)
}

// Java: ProcessCommonJSModulesTest.testModuleExportsWrittenWithExportsRefs()
func (c *ProcessCommonJSModulesTest) TestModuleExportsWrittenWithExportsRefs() {
	c.TestModules(
		"test.js",
		`exports.one = 1;
module.exports = {};
`,
		"/** @const */ var module$test = { default: {}}; module$test.default.one = 1;")
}

// Java: ProcessCommonJSModulesTest.testVarRenaming()
func (c *ProcessCommonJSModulesTest) TestVarRenaming() {
	c.TestModules(
		"test.js",
		`module.exports = {};
var a = 1, b = 2;
(function() { var a; b = 4})();
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
var a$$module$test = 1;
var b$$module$test = 2;
(function() { var a; b$$module$test = 4})();
`)
}

// Java: ProcessCommonJSModulesTest.testDash()
func (c *ProcessCommonJSModulesTest) TestDash() {
	c.TestModules(
		"test-test.js",
		`var name = require('./other');
exports.foo = 1;
`,
		`/** @const */ var module$test_test = {/** @const */ default: {}};
var name$$module$test_test=module$other.default
module$test_test.default.foo = 1;
`)
}

// Java: ProcessCommonJSModulesTest.testIndex()
func (c *ProcessCommonJSModulesTest) TestIndex() {
	c.TestModules(
		"foo/index.js",
		`var name = require('../other');
exports.bar = 1;
`,
		`/** @const */ var module$foo$index = {/** @const */ default: {}};
var name$$module$foo$index = module$other.default;
module$foo$index.default.bar = 1;
`)
}

// Java: ProcessCommonJSModulesTest.testVarJsdocGoesOnAssignment()
func (c *ProcessCommonJSModulesTest) TestVarJsdocGoesOnAssignment() {
	c.TestModules(
		"testcode.js",
		`/**
 * @const
 * @enum {number}
 */
var MyEnum = { ONE: 1, TWO: 2 };
module.exports = {MyEnum: MyEnum};
`,
		`/** @const */ var module$testcode = {/** @const */ default: {}};
/**
 * @const
 * @enum {number}
 */
(module$testcode.default.MyEnum = {ONE:1, TWO:2});
`)
}

// Java: ProcessCommonJSModulesTest.testModuleName()
func (c *ProcessCommonJSModulesTest) TestModuleName() {
	c.TestModules(
		"foo/bar.js",
		`var name = require('../other');
module.exports = name;
`,
		`/** @const */ var module$foo$bar = {};
var name$$module$foo$bar = module$other.default;
/** @const */ module$foo$bar.default = module$other.default;
`)

	c.Test(
		c.Srcs(
			SourceFile_FromCode(Compiler_JoinPathParts("foo", "name.js"), ""),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`var name = require('./name');
module.exports = name;
`)),
		c.Expected(
			SourceFile_FromCode(Compiler_JoinPathParts("foo", "name.js"), ""),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @const */ var module$foo$bar = {};
var name$$module$foo$bar = module$foo$name.default;
/** @const */ module$foo$bar.default = module$foo$name.default;
`)))
}

// Java: ProcessCommonJSModulesTest.testModuleExportsScope()
func (c *ProcessCommonJSModulesTest) TestModuleExportsScope() {
	c.TestModules(
		"test.js",
		`var foo = function (module) {
  module.exports = {};
};
module.exports = foo;
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = function (module) {
  module.exports={};
};
`)

	c.TestModules(
		"test.js",
		`var foo = function () {
  var module = {};
  module.exports = {};
};
module.exports = foo;
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = function() {
  var module={};
  module.exports={}
};
`)

	c.TestModules(
		"test.js",
		`var foo = function () {
  if (true) var module = {};
  module.exports = {};
};
module.exports = foo;
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = function() {
  if (true) var module={};
  module.exports={}
};
`)
}

// Java: ProcessCommonJSModulesTest.testUMDPatternConversion()
func (c *ProcessCommonJSModulesTest) TestUMDPatternConversion() {
	c.TestModules(
		"test.js",
		`var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof window.define === 'function' && window.define.amd) {
  window.define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`var foobar = {foo: 'bar'};
if (typeof module === 'object' && module['exports']) {
  module['exports'] = foobar;
} else if (typeof define === 'function' && define['amd']) {
  define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`var foobar = {foo: 'bar'};
if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else {
  this.foobar = foobar;
}
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`var foobar = {foo: 'bar'};
if (typeof window.define === 'function' && window.define.amd) {
  window.define([], function() {return foobar;});
} else if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else {
  this.foobar = foobar;
}
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`var foobar = {foo: 'bar'};
if (typeof define === 'function' && define['amd']) {
  define([], function() {return foobar;});
} else if (typeof module === 'object' && module['exports']) {
  module['exports'] = foobar;
} else {
  this.foobar = foobar;
}
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
}
if (typeof define === 'function' && define.amd) {
  define([], function () {return foobar;});
}
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
}
if (typeof window.define === 'function' && window.define.amd) {
  window.define([], function () {return foobar;});
}
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`(function(global, factory) {
  true ? module.exports = factory(
             typeof angular === 'undefined' ? require('./other') :
                                              angular) :
         typeof define === 'function' && define.amd ?
         define('angular-cache', ['angular'], factory) :
         (global.angularCacheModuleName = factory(global.angular));
}(this, function(angular) {
  'use strict';
  console.log(angular);
  return angular;
}));
`,
		`/** @const */ var module$test = {};
var global$$module$test = this;
var factory$$module$test = function(angular) {
  console.log(angular);
  return angular;
};
/** @const */
module$test.default = factory$$module$test(typeof angular === "undefined" ? module$other.default : angular);
`)

	c.TestModules(
		"test.js",
		`(function(global, factory) {
  true ? module.exports = factory(
             typeof angular === 'undefined' ? require('./other') :
                                              angular) :
         typeof window.define === 'function' && window.define.amd ?
         window.define('angular-cache', ['angular'], factory) :
         (global.angularCacheModuleName = factory(global.angular));
}(this, function(angular) {
  'use strict';
  console.log(angular);
  return angular;
}));
`,
		`/** @const */ var module$test = {};
var global$$module$test = this;
var factory$$module$test = function(angular) {
  console.log(angular);
  return angular;
};
/** @const */
module$test.default = factory$$module$test(typeof angular === "undefined" ? module$other.default : angular);
`)
}

// Java: ProcessCommonJSModulesTest.testEs6ObjectShorthand()
func (c *ProcessCommonJSModulesTest) TestEs6ObjectShorthand() {
	c.TestModules(
		"test.js",
		`function foo() {}
module.exports = {
  prop: 'value',
  foo
};
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
module$test.default.foo = function () {};
module$test.default.prop = 'value';
`)

	c.TestModules(
		"test.js",
		`module.exports = {
  prop: 'value',
  foo() {
    console.log('bar');
  }
};
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
module$test.default.prop = 'value';
module$test.default.foo = function() {
  console.log('bar');
};
`)

	c.TestModules(
		"test.js",
		`var a = require('./other');
module.exports = {a: a};
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
var a$$module$test = module$other.default;
module$test.default.a = module$other.default;
`)

	c.TestModules(
		"test.js",
		`var a = require('./other');
module.exports = {a};
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
var a$$module$test = module$other.default;
module$test.default.a = module$other.default;
`)

	c.TestModules(
		"test.js",
		`var a = 4;
module.exports = {a};
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
module$test.default.a = 4;
`)
}

// Java: ProcessCommonJSModulesTest.testKeywordsInExports()
func (c *ProcessCommonJSModulesTest) TestKeywordsInExports() {
	c.TestModules(
		"testcode.js",
		`var a = 4;
module.exports = { else: a };
`,
		`/** @const */ var module$testcode = {/** @const */ default: {}};
module$testcode.default.else = 4;
`)
}

// Java: ProcessCommonJSModulesTest.testRequireResultUnused()
func (c *ProcessCommonJSModulesTest) TestRequireResultUnused() {
	c.TestModules("test.js", "require('./other');", "")
}

// Java: ProcessCommonJSModulesTest.testRequireEnsure()
func (c *ProcessCommonJSModulesTest) TestRequireEnsure() {
	c.TestModules(
		"test.js",
		`require.ensure(['./other'], function(require) {
  var other = require('./other');
  var bar = other;
});
`,
		`(function() {
  var other = module$other.default;
  var bar = module$other.default;
})()
`)
}

// Java: ProcessCommonJSModulesTest.testFunctionRewriting()
func (c *ProcessCommonJSModulesTest) TestFunctionRewriting() {
	c.TestModules(
		"test.js",
		`function foo() {}
foo.prototype = new Date();
module.exports = foo;
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = function() {};
module$test.default.prototype = new Date();
`)

	c.TestModules(
		"test.js",
		`function foo() {}
foo.prototype = new Date();
module.exports = {foo: foo};
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
module$test.default.foo = function () {};
module$test.default.foo.prototype = new Date();
`)
}

// Java: ProcessCommonJSModulesTest.testFunctionHoisting()
func (c *ProcessCommonJSModulesTest) TestFunctionHoisting() {
	c.TestModules(
		"test.js",
		`module.exports = foo;
function foo() {}
foo.prototype = new Date();
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = function() {};
module$test.default.prototype = new Date();
`)

	c.TestModules(
		"test.js",
		`function foo() {}
Object.assign(foo, { bar: foobar });
function foobar() {}
module.exports = foo;
module.exports.bar = foobar;
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = function () {};
module$test.default.bar = function() {};
Object.assign(module$test.default, { bar: module$test.default.bar });
`)
}

// Java: ProcessCommonJSModulesTest.testClassRewriting()
func (c *ProcessCommonJSModulesTest) TestClassRewriting() {
	c.TestModules(
		"test.js",
		`class foo extends Array {}
module.exports = foo;
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = class extends Array {};
`)

	c.TestModules(
		"test.js",
		`class foo {}
module.exports = foo;
`,
		"/** @const */ var module$test = {}; /** @const */ module$test.default = class {}")

	c.TestModules(
		"test.js",
		`class foo {}
module.exports.foo = foo;
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
module$test.default.foo = class {};
`)

	c.TestModules(
		"test.js",
		"module.exports = class { bar() { return 'bar'; }};",
		`/** @const */ var module$test = {};
/** @const */ module$test.default = class {
  bar() { return 'bar'; }
};
`)
}

// Java: ProcessCommonJSModulesTest.testMultipleAssignments()
func (c *ProcessCommonJSModulesTest) TestMultipleAssignments() {
	chunk := NewJSChunk("out")
	chunk.AddSourceFile(SourceFile_FromCode("other.js", "goog.provide('module$other');"))
	chunk.AddSourceFile(SourceFile_FromCode("yet_another.js", "goog.provide('module$yet_another');"))
	chunk.AddSourceFile(
		SourceFile_FromCode(
			"test",
			`/** @constructor */ function Hello() {}
module.exports = Hello;
/** @constructor */ function Bar() {}
Bar.prototype.foobar = function() { alert('foobar'); };
exports = Bar;
`))
	chunks := []*JSChunk{chunk}
	c.Test(
		c.Srcs(chunks),
		c.Expected([]string(nil)),
		newCompilerTestCase_Diagnostic(
			ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT.level,
			ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT,
			nil))
}

// Java: ProcessCommonJSModulesTest.testDestructuringImports()
func (c *ProcessCommonJSModulesTest) TestDestructuringImports() {
	c.TestModules(
		"test.js",
		`const {foo, bar} = require('./other');
var baz = foo + bar;
`,
		`const {foo, bar} = module$other.default;
var baz = module$other.default.foo + module$other.default.bar;
`)
}

// Java: ProcessCommonJSModulesTest.testDestructuringImports2()
func (c *ProcessCommonJSModulesTest) TestDestructuringImports2() {
	c.TestModules(
		"test.js",
		`const {foo, bar: {baz}} = require('./other');
module.exports = true;
`,
		`/** @const */ var module$test = {};
const {foo: foo$$module$test, bar: {baz: baz$$module$test}} = module$other.default;
/** @const */ module$test.default = true;
`)
}

// Java: ProcessCommonJSModulesTest.testAnnotationsCopied()
func (c *ProcessCommonJSModulesTest) TestAnnotationsCopied() {
	c.TestModules(
		"test.js",
		`/** @interface */ var a;
/** @type {string} */ a.prototype.foo;
module.exports.a = a;
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
/** @interface */ module$test.default.a;
/** @type {string} */ module$test.default.a.prototype.foo;
`)
}

// Java: ProcessCommonJSModulesTest.testUMDRemoveIIFE()
func (c *ProcessCommonJSModulesTest) TestUMDRemoveIIFE() {
	c.TestModules(
		"test.js",
		`(function(){
var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}})()
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`!function(){
var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}}()
`,
		"/** @const */ var module$test = {}; /** @const */ module$test.default = {foo: 'bar'};")

	c.TestModules(
		"test.js",
		`!function(){
var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}}()
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`;;;(function(){
var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}})()
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`(function(){
var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}}.call(this))
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {foo: 'bar'};
`)

	c.TestModules(
		"test.js",
		`;;;(function(global){
var foobar = {foo: 'bar'};
global.foobar = foobar;
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  global.foobar = foobar;
}})(this)
`,
		`/** @const */ var module$test = { default: {}};
module$test.default = {foo: 'bar'};
module$test.default.foobar = module$test.default;
`)

	c.TestModules(
		"test.js",
		`(function(global){
var foobar = {foo: 'bar'};
global.foobar = foobar;
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  global.foobar = foobar;
}}.call(this, this))
`,
		`/** @const */ var module$test = { default: {}};
module$test.default = {foo: 'bar'};
module$test.default.foobar = module$test.default;
`)

	// We can't remove IIFEs explict calls that don't use "this"
	c.TestModules(
		"test.js",
		`(function(){
var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}}.call(window))
`,
		`/** @const */ var module$test = {};
(function(){
  var foobar = {foo: 'bar'};
  /** @const */ module$test.default=foobar;
}).call(window);
`)

	// Can't remove IIFEs when there are sibling statements
	c.TestModules(
		"test.js",
		`(function(){
var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}})();
alert('foo');
`,
		`/** @const */ var module$test = {};
(function(){
  var foobar = {foo: 'bar'};
  /** @const */ module$test.default = foobar;
})();
alert('foo');
`)

	// Can't remove IIFEs when there are sibling statements
	c.TestModules(
		"test.js",
		`alert('foo');
(function(){
var foobar = {foo: 'bar'};
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  this.foobar = foobar;
}})();
`,
		`/** @const */ var module$test = {};
alert('foo');
(function(){
  var foobar={foo:"bar"};
  /** @const */ module$test.default=foobar;
})();
`)

	// Annotations for local names should be preserved
	c.TestModules(
		"test.js",
		`(function(global){
/** @param {...*} var_args */
function log(var_args) {}
var foobar = {foo: 'bar', log: function() { log.apply(null, arguments); } };
global.foobar = foobar;
if (typeof module === 'object' && module.exports) {
  module.exports = foobar;
} else if (typeof define === 'function' && define.amd) {
  define([], function() {return foobar;});
} else {
  global.foobar = foobar;
}}.call(this, this))
`,
		`/** @const */ var module$test = { default: {}};
/** @param {...*} var_args */
function log$$module$test(var_args){}
module$test.default = {
  foo: 'bar',
  log: function() { log$$module$test.apply(null,arguments); }
};
module$test.default.foobar = module$test.default;
`)
}

// Java: ProcessCommonJSModulesTest.testParamShadow()
func (c *ProcessCommonJSModulesTest) TestParamShadow() {
	c.TestModules(
		"test.js",
		`/** @constructor */ function Foo() {}
/** @constructor */ function Bar(Foo) { this.foo = new Foo(); }
Foo.prototype.test = new Bar(Foo);
module.exports = Foo;
`,
		`/** @const */ var module$test = {};
/** @const @constructor */ module$test.default = function () {};
/** @constructor */ function Bar$$module$test(Foo) { this.foo = new Foo(); }
module$test.default.prototype.test = new Bar$$module$test(module$test.default);
`)
}

// Java: ProcessCommonJSModulesTest.testIssue2308()
func (c *ProcessCommonJSModulesTest) TestIssue2308() {
	c.TestModules(
		"test.js",
		"exports.y = null; var x; x = exports.y;",
		`/** @const */ var module$test = {/** @const */ default: {}};
module$test.default.y = null;
var x$$module$test;
x$$module$test = module$test.default.y
`)
}

// Java: ProcessCommonJSModulesTest.testAbsoluteImportsWithModuleRoots()
func (c *ProcessCommonJSModulesTest) TestAbsoluteImportsWithModuleRoots() {
	c.moduleRoots = []string{"/base"}
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("base", "mod", "name.js"), "module.exports = {}"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("base", "test", "sub.js"),
				`var name = require('/mod/name');
(function() { let foo = name; foo(); })();
`)),
		c.Expected(
			SourceFile_FromCode(
				Compiler_JoinPathParts("base", "mod", "name.js"),
				"/** @const */ var module$mod$name = {/** @const */ default: {}};"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("base", "test", "sub.js"),
				`var name = module$mod$name.default;
(function() { let foo = module$mod$name.default; foo(); })();
`)))
}

// Java: ProcessCommonJSModulesTest.testIssue2510()
func (c *ProcessCommonJSModulesTest) TestIssue2510() {
	c.TestModules(
		"test.js",
		"module.exports = {a: 1, get b() { return 2; }};",
		`/** @const */ var module$test = {};
/** @const */ module$test.default = {
  get b() { return 2; }
};
module$test.default.a = 1;
`)
}

// Java: ProcessCommonJSModulesTest.testIssue2450()
func (c *ProcessCommonJSModulesTest) TestIssue2450() {
	c.TestModules(
		"test.js",
		`var BCRYPT_BLOCKS = 8,
    BCRYPT_HASHSIZE = 32;

module.exports = {
  BLOCKS: BCRYPT_BLOCKS,
  HASHSIZE: BCRYPT_HASHSIZE,
};
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
module$test.default.BLOCKS = 8;
module$test.default.HASHSIZE = 32;
`)
}

// Java: ProcessCommonJSModulesTest.testWebpackAmdPattern()
func (c *ProcessCommonJSModulesTest) TestWebpackAmdPattern() {
	c.TestModules(
		"test.js",
		`var __WEBPACK_AMD_DEFINE_ARRAY__, __WEBPACK_AMD_DEFINE_RESULT__;
!(__WEBPACK_AMD_DEFINE_ARRAY__ =
      [__webpack_require__(1), __webpack_require__(2)],
  __WEBPACK_AMD_DEFINE_RESULT__ =
      function(b, c) {
        console.log(b, c.exportA, c.exportB);
      }.apply(exports, __WEBPACK_AMD_DEFINE_ARRAY__),
  __WEBPACK_AMD_DEFINE_RESULT__ !== undefined &&
      (module.exports = __WEBPACK_AMD_DEFINE_RESULT__));
`,
		`/** @const */ var module$test = { default: {}};
var __WEBPACK_AMD_DEFINE_ARRAY__$$module$test;
!(__WEBPACK_AMD_DEFINE_ARRAY__$$module$test =
    [__webpack_require__(1), __webpack_require__(2)],
    module$test.default = function(b,c){console.log(b,c.exportA,c.exportB)}
        .apply(module$test.default,__WEBPACK_AMD_DEFINE_ARRAY__$$module$test),
    module$test.default!==undefined && module$test.default)
`)

	c.TestModules(
		"test.js",
		`/** @suppress {duplicate} */var __WEBPACK_AMD_DEFINE_RESULT__;(function() {
  var dialogPolyfill = {prop: 'DIALOG_POLYFILL'};

  if ('function' === 'function' && 'amd' in __webpack_require__(124)) {
    // AMD support
    !(__WEBPACK_AMD_DEFINE_RESULT__ = (function() { return dialogPolyfill; }).call(exports, __webpack_require__, exports, module),
        __WEBPACK_AMD_DEFINE_RESULT__ !== undefined && (module.exports = __WEBPACK_AMD_DEFINE_RESULT__));
  } else if (typeof module === 'object' && typeof module['exports'] === 'object') {
    // CommonJS support
    module['exports'] = dialogPolyfill;
  } else {
    // all others
    window['dialogPolyfill'] = dialogPolyfill;
  }
})();
`,
		`/** @const */ var module$test = {default: {}};
/** @suppress {duplicate} */
var __WEBPACK_AMD_DEFINE_RESULT__$$module$test;
(function () {
  var dialogPolyfill = {prop: "DIALOG_POLYFILL"};
  !(__WEBPACK_AMD_DEFINE_RESULT__$$module$test = function () {
    return dialogPolyfill
  }.call(module$test.default, __webpack_require__, module$test.default, {}),
  __WEBPACK_AMD_DEFINE_RESULT__$$module$test !== undefined && (module$test.default = __WEBPACK_AMD_DEFINE_RESULT__$$module$test))
})()
`)

	webpackModulesById :=
		esModulesTest_mapOf(
			"1", "other.js",
			"yet_another.js", "yet_another.js",
			"3", "test.js")

	c.SetWebpackModulesById(webpackModulesById)
	c.resolutionMode = ModuleLoader_ResolutionMode_WEBPACK

	c.TestModules(
		"test.js",
		`/** @suppress {duplicate} */var __WEBPACK_AMD_DEFINE_ARRAY__, __WEBPACK_AMD_DEFINE_RESULT__;(function (root, factory) {
  if (true) {
    !(__WEBPACK_AMD_DEFINE_ARRAY__ = [__webpack_require__(1),__webpack_require__('yet_another.js')], __WEBPACK_AMD_DEFINE_RESULT__ = function (a0,b1) {
      return (factory(a0,b1));
    }.apply(exports, __WEBPACK_AMD_DEFINE_ARRAY__),
      __WEBPACK_AMD_DEFINE_RESULT__ !== undefined && (module.exports = __WEBPACK_AMD_DEFINE_RESULT__));
  } else if (typeof module === 'object' && module.exports) {
    module.exports = factory(require('angular'),require('tinymce'));
  } else {
    root['banno.wysiwyg'] = factory(root['angular'],root['tinymce']);
  }
}(this, function (angular, tinymce) {
  console.log(angular, tinymce);
}))
`,
		`/** @const */ var module$test = {default: {}};
/** @suppress {duplicate} */
var __WEBPACK_AMD_DEFINE_ARRAY__$$module$test;
/** @suppress {duplicate} */
module$test.default;
var root$$module$test = this;
var factory$$module$test = function (angular, tinymce) {
  console.log(angular, tinymce)
};
!(__WEBPACK_AMD_DEFINE_ARRAY__$$module$test = [module$other.default,
  module$yet_another.default], module$test.default = function (a0, b1) {
  return factory$$module$test(a0, b1)
}.apply(module$test.default, __WEBPACK_AMD_DEFINE_ARRAY__$$module$test),
  module$test.default !== undefined && module$test.default)
`)
}

// Java: ProcessCommonJSModulesTest.testIssue2593()
func (c *ProcessCommonJSModulesTest) TestIssue2593() {
	c.TestModules(
		"test.js",
		`var first = 1,
    second = 2,
    third = 3,
    fourth = 4,
    fifth = 5;

module.exports = {};
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
var first$$module$test=1;
var second$$module$test=2;
var third$$module$test=3;
var fourth$$module$test=4;
var fifth$$module$test=5;
`)
}

// Java: ProcessCommonJSModulesTest.testTernaryUMDWrapper()
func (c *ProcessCommonJSModulesTest) TestTernaryUMDWrapper() {
	c.TestModules(
		"test.js",
		`var foobar = {foo: 'bar'};
typeof module === 'object' && module.exports ?
   module.exports = foobar :
   typeof define === 'function' && define.amd ?
     define([], function() {return foobar;}) :
     this.foobar = foobar;
`,
		"/** @const */ var module$test = {}; /** @const */ module$test.default = {foo: 'bar'};")
}

// Java: ProcessCommonJSModulesTest.testLeafletUMDWrapper()
func (c *ProcessCommonJSModulesTest) TestLeafletUMDWrapper() {
	c.TestModules(
		"test.js",
		`(function (global, factory) {
  typeof exports === 'object' && typeof module !== 'undefined' ?
    factory(exports) :
    typeof define === 'function' && define.amd ?
      define(['exports'], factory) :
      (factory((global.L = {})));
}(this, (function (exports) {
  'use strict';
  var webkit = userAgentContains('webkit');
  function userAgentContains(str) {
    return navigator.userAgent.toLowerCase().indexOf(str) >= 0;
  }
  exports.webkit = webkit
})));
`,
		`/** @const */ var module$test={/** @const */ default: {}};
var global$$module$test = this;
var factory$$module$test = function(exports) {
  var webkit = userAgentContains("webkit");
  function userAgentContains(str) {
    return navigator.userAgent.toLowerCase().indexOf(str) >= 0;
  }
  exports.webkit = webkit;
};
factory$$module$test(module$test.default);
`)
}

// Java: ProcessCommonJSModulesTest.testBowserUMDWrapper()
func (c *ProcessCommonJSModulesTest) TestBowserUMDWrapper() {
	c.TestModules(
		"test.js",
		`!function (root, name, definition) {
  if (typeof module != 'undefined' && module.exports)
    module.exports = definition()
  else if (typeof define == 'function' && define.amd)
    define(name, definition)
  else root[name] = definition()
}(this, 'foobar', function () {
  return {foo: 'bar'};
});
`,
		`/** @const */ var module$test={};
var root$$module$test = this;
var name$$module$test = "foobar";
var definition$$module$test = function() {
  return {foo: 'bar'};
};
/** @const */
module$test.default = definition$$module$test();
`)
}

// Java: ProcessCommonJSModulesTest.testDontSplitVarsInFor()
func (c *ProcessCommonJSModulesTest) TestDontSplitVarsInFor() {
	c.TestModules("test.js", "for (var a, b, c; ;) {}", "for (var a, b, c; ;) {}")
}

// Java: ProcessCommonJSModulesTest.testIssue2918()
func (c *ProcessCommonJSModulesTest) TestIssue2918() {
	c.TestModules(
		"test.js",
		`for (var a, b; a < 4; a++) {};
module.exports = {}
`,
		`/** @const */ var module$test = {/** @const */ default:{}};
for(var a$$module$test,b$$module$test;a$$module$test<4;a$$module$test++) {};
`)
}

// Java: ProcessCommonJSModulesTest.testExportsDirectAssignment()
func (c *ProcessCommonJSModulesTest) TestExportsDirectAssignment() {
	c.TestModules(
		"test.js",
		"exports = module.exports = {};",
		"/** @const */ var module$test = {/** @const */ default: {}};")
}

// Java: ProcessCommonJSModulesTest.testExportsPropertyHoisting()
func (c *ProcessCommonJSModulesTest) TestExportsPropertyHoisting() {
	c.TestModules(
		"test.js",
		`exports.Buffer = Buffer;
Buffer.TYPED_ARRAY_SUPPORT = {};
function Buffer() {}
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
module$test.default.Buffer = function() {};
module$test.default.Buffer.TYPED_ARRAY_SUPPORT = {};
`)
}

// Java: ProcessCommonJSModulesTest.testExportNameInParamList()
func (c *ProcessCommonJSModulesTest) TestExportNameInParamList() {
	c.TestModules(
		"test.js",
		`var tinymce = { foo: 'bar' };
function register(cb) { cb(tinymce); }
register(function(tinymce) { module.exports = tinymce; });
`,
		`/** @const */ var module$test = {};
var tinymce$$module$test = { foo: 'bar' };
function register$$module$test(cb) { cb(tinymce$$module$test); }
register$$module$test(function(tinymce) {
  /** @const */ module$test.default = tinymce;
});
`)
}

// Java: ProcessCommonJSModulesTest.testIssue2616()
func (c *ProcessCommonJSModulesTest) TestIssue2616() {
	c.TestModules(
		"test.js",
		`var foo = function foo() {
  return 1;
};
module.exports = {
  foo: foo,
};
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
module$test.default.foo = function foo() {
  return 1;
};
`)
}

// Java: ProcessCommonJSModulesTest.testFingerprintUmd()
func (c *ProcessCommonJSModulesTest) TestFingerprintUmd() {
	c.TestModules(
		"test.js",
		`(function (name, context, definition) {
  'use strict';
  if (typeof define === 'function' && define.amd) {
    define(definition);
  } else if (typeof module !== 'undefined' && module.exports) {
    module.exports = definition();
  } else if (context.exports) {
    context.exports = definition();
  } else {
    context[name] = definition();
  }
})('Fingerprint2', this, function() {
  var Fingerprint2 = function() {
    if (!(this instanceof Fingerprint2)) { return new Fingerprint2(); }
  };
  return Fingerprint2;
})
`,
		`/** @const */ var module$test = {};
var name$$module$test = "Fingerprint2";
var context$$module$test = this;
var definition$$module$test = function() {
  var Fingerprint2 = function() {
    if (!(this instanceof Fingerprint2)) {
      return new Fingerprint2();
    }
  };
  return Fingerprint2;
};
/** @const */
module$test.default = definition$$module$test();
`)
}

// Java: ProcessCommonJSModulesTest.testTypeofModuleReference()
func (c *ProcessCommonJSModulesTest) TestTypeofModuleReference() {
	c.TestModules(
		"test.js",
		`module.exports = 'foo';
console.log(typeof module);
console.log(typeof exports);
`,
		`/** @const */ var module$test={ default: {}};
module$test.default = 'foo';
console.log('object');
console.log('object');
`)
}

// Java: ProcessCommonJSModulesTest.testUpdateGenericTypeReferences()
func (c *ProcessCommonJSModulesTest) TestUpdateGenericTypeReferences() {
	c.TestModules(
		"test.js",
		`const Foo = require('./other');
/** @type {!Array<!Foo>} */ const bar = [];
module.exports = bar;
`,
		`/** @const */ var module$test={};
const Foo$$module$test = module$other.default;
/** @const  @type {!Array<!module$other.default>} */ module$test.default = [];
`)
}

// Java: ProcessCommonJSModulesTest.testMissingRequire()
func (c *ProcessCommonJSModulesTest) TestMissingRequire() {
	ModulesTestUtils_TestModulesError(c, "require('missing');", ModuleLoader_LOAD_WARNING)

	c.TestModules(
		"test.js",
		`/**
 * @fileoverview
 * @suppress {moduleLoad}
 */
var foo = require('missing');
`,
		`/**
 * @fileoverview
 * @suppress {moduleLoad}
 */
var foo = module$missing.default;
`)
}

// Java: ProcessCommonJSModulesTest.testExportsUsageInIf()
func (c *ProcessCommonJSModulesTest) TestExportsUsageInIf() {
	c.TestModules(
		"test.js",
		`exports.merge = function(source) {
  return Object.keys(source).reduce(function (acc, key) {
    if (Object.prototype.hasOwnProperty.call(acc, key)) {
      acc[key] = exports.merge(acc[key], value, options);
    } else {
      acc[key] = value;
    }
    return acc;
  }, {});
};
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
module$test.default.merge = function(source) {
  return Object.keys(source).reduce(function(acc,key) {
    if (Object.prototype.hasOwnProperty.call(acc,key)) {
      acc[key] = module$test.default.merge(acc[key],value,options);
    } else {
      acc[key] = value;
    }
    return acc;
  }, {});
}
`)
}

// Java: ProcessCommonJSModulesTest.testModuleId()
func (c *ProcessCommonJSModulesTest) TestModuleId() {
	c.TestModules(
		"test.js",
		"module.exports = module.id;",
		`/** @const */ var module$test = {};
/** @const */ module$test.default = 'test.js';
`)
}

// Java: ProcessCommonJSModulesTest.testModuleIdAlias()
func (c *ProcessCommonJSModulesTest) TestModuleIdAlias() {
	c.TestModules(
		"test.js",
		`module.exports = 'foo';
function foobar(module) { return module.id; }
`,
		`/** @const */ var module$test = {};
/** @const */ module$test.default = 'foo';
function foobar$$module$test(module) { return module.id; }
`)
}

// Java: ProcessCommonJSModulesTest.testWebpackRequire()
func (c *ProcessCommonJSModulesTest) TestWebpackRequire() {
	webpackModulesById :=
		esModulesTest_mapOf(
			"1", "other.js",
			"2", "yet_another.js",
			"3", "test.js")
	c.SetWebpackModulesById(webpackModulesById)
	c.resolutionMode = ModuleLoader_ResolutionMode_WEBPACK

	c.TestModules(
		"test.js",
		`var name = __webpack_require__(1);
exports.foo = 1;
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
var name$$module$test = module$other.default;
module$test.default.foo = 1;
`)
}

// Java: ProcessCommonJSModulesTest.testWebpackRequireString()
func (c *ProcessCommonJSModulesTest) TestWebpackRequireString() {
	webpackModulesById :=
		esModulesTest_mapOf(
			"1", "other.js",
			"yet_another.js", "yet_another.js",
			"3", "test.js")
	c.SetWebpackModulesById(webpackModulesById)
	c.resolutionMode = ModuleLoader_ResolutionMode_WEBPACK

	c.TestModules(
		"test.js",
		`var name = __webpack_require__('yet_another.js');
exports.foo = 1;
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
var name$$module$test = module$yet_another.default;
module$test.default.foo = 1;
`)
}

// Java: ProcessCommonJSModulesTest.testWebpackAMDModuleShim()
func (c *ProcessCommonJSModulesTest) TestWebpackAMDModuleShim() {
	webpackModulesById :=
		esModulesTest_mapOf(
			"1", "test.js",
			"2", "/webpack/buildin/module.js")
	c.SetWebpackModulesById(webpackModulesById)
	c.resolutionMode = ModuleLoader_ResolutionMode_WEBPACK

	// Shared with ProcessCommonJSModulesTest.
	inputs :=
		[]*SourceFile{
			SourceFile_FromCode(
				"test.js",
				`(function(module) {
  console.log(module.id);
})(__webpack_require__(2)(module))
`),
			SourceFile_FromCode(
				"/webpack/buildin/module.js",
				"module.exports = function(module) { return module; };")}
	expecteds :=
		[]*SourceFile{
			SourceFile_FromCode("test.js", "(function(){console.log('test.js')})()"),
			SourceFile_FromCode(
				"/webpack/buildin/module.js",
				`/** @const */ var module$webpack$buildin$module = {};
/** @const */ module$webpack$buildin$module.default =
    function(module) { return module; };
`)}
	c.Test(c.Srcs(inputs), c.Expected(expecteds))
}

// Java: ProcessCommonJSModulesTest.testComplexExportAssignment()
func (c *ProcessCommonJSModulesTest) TestComplexExportAssignment() {
	c.TestModules(
		"test.js",
		"const width = 800; const vwidth = exports.vwidth = width;",
		`/** @const */ var module$test = { /** @const */ default: {}};
module$test.default.vwidth = 800;
const vwidth$$module$test = module$test.default.vwidth;
`)
}

// Java: ProcessCommonJSModulesTest.testUMDRequiresIfTest()
func (c *ProcessCommonJSModulesTest) TestUMDRequiresIfTest() {
	c.TestModules(
		"test.js",
		"var foobar = {foo: 'bar'}; if (foobar) { module.exports = foobar; }",
		`/** @const */ var module$test = {};
var foobar$$module$test={foo:"bar"};
if(foobar$$module$test) {
  /** @const */ module$test.default = foobar$$module$test;
}
`)
}

// Java: ProcessCommonJSModulesTest.testObjectSpreadExport()
func (c *ProcessCommonJSModulesTest) TestObjectSpreadExport() {
	c.TestModules(
		"test.js",
		"var g = {}; module.exports = { ...g };",
		`/** @const */ var module$test = {};
var g$$module$test = {};
/** @const */ module$test.default = {
  ...g$$module$test
};
`)
}

// Java: ProcessCommonJSModulesTest.testBabelTranspiledESModules()
func (c *ProcessCommonJSModulesTest) TestBabelTranspiledESModules() {
	c.TestModules(
		"test.js",
		`'use strict';

Object.defineProperty(exports, '__esModule', {
  value: true
});
exports.default = toInteger;
function toInteger(dirtyNumber) { }
module.exports = exports['default'];
`,
		`/** @const */ var module$test = {default: {}};
module$test.default.default = function(dirtyNumber) { };
Object.defineProperty(module$test.default, '__esModule',{value:true})
module$test.default = module$test.default.default;
`)
}

// Java: ProcessCommonJSModulesTest.testLodashModulesCheck()
func (c *ProcessCommonJSModulesTest) TestLodashModulesCheck() {
	c.TestModules(
		"test.js",
		"/* Detect free variable `exports`. */\nconst freeExports = typeof exports == 'object' && exports !== null\n    && !exports.nodeType && exports\n/* Detect free variable `module`. */\nconst freeModule = freeExports && typeof module == 'object' && module !== null\n    && !module.nodeType && module\nconsole.log(freeExports, freeModule);\nmodule.exports = true;\n",
		`/** @const */ var module$test = {default: {}};
const freeExports$$module$test = 'object' == 'object' && module$test.default !== null
    && !module$test.default.nodeType && module$test.default;
const freeModule$$module$test = freeExports$$module$test && 'object' == 'object'
    && {} !== null && !{}.nodeType && {};
console.log(freeExports$$module$test, freeModule$$module$test);
module$test.default = true;
`)
}

// Java: ProcessCommonJSModulesTest.testIssue3051()
func (c *ProcessCommonJSModulesTest) TestIssue3051() {
	c.TestModules(
		"test.js",
		`class Base {}
exports.Base = Base;

class Impl extends exports.Base {
    getString() {
        return "test";
    }
}
exports.Impl = Impl;

const w = new exports.Impl("a")
console.log(w.getString());
`,
		`/** @const */ var module$test = {
    /** @const */ default: {}
};
module$test.default.Base = class {};
module$test.default.Impl = class extends module$test.default.Base {
    getString() {
        return "test"
    }
};
const w$$module$test = new module$test.default.Impl("a");
console.log(w$$module$test.getString());
`)
}

// Java: ProcessCommonJSModulesTest.testDestructuredImportExported()
func (c *ProcessCommonJSModulesTest) TestDestructuredImportExported() {
	c.TestModules(
		"test.js",
		"const {Foo} = require('./other.js'); Foo; exports.Foo = Foo;",
		`/** @const */ var module$test = {
    /** @const */ default: {}
};
const {Foo: Foo$$module$test} = module$other.default;
module$other.default.Foo;
module$test.default.Foo = module$other.default.Foo;
`)
}

// Java: ProcessCommonJSModulesTest.testDestructuredExports()
func (c *ProcessCommonJSModulesTest) TestDestructuredExports() {
	c.TestModules(
		"test.js",
		"const {b} = {b: 1}; module.exports = {b: b};",
		`/** @const */ var module$test = {
    /** @const */ default: {}
};
const {b: b$$module$test} = {b: 1};
module$test.default.b = b$$module$test;
`)
}

// Java: ProcessCommonJSModulesTest.testWebpackRequireNamespace()
func (c *ProcessCommonJSModulesTest) TestWebpackRequireNamespace() {
	webpackModulesById :=
		esModulesTest_mapOf(
			"1", "other.js",
			"yet_another.js", "yet_another.js",
			"3", "test.js")
	c.SetWebpackModulesById(webpackModulesById)
	c.resolutionMode = ModuleLoader_ResolutionMode_WEBPACK

	c.TestModules(
		"test.js",
		`var name = __webpack_require__.t('yet_another.js');
exports.foo = 1;
`,
		`/** @const */ var module$test = {/** @const */ default: {}};
var name$$module$test = module$yet_another;
module$test.default.foo = 1;
`)
}

// Java: ProcessCommonJSModulesTest.testGoogModuleUnaffected()
func (c *ProcessCommonJSModulesTest) TestGoogModuleUnaffected() {
	c.TestModules(
		"test.js", "goog.module('foo'); exports.y = 123;", "goog.module('foo'); exports.y = 123;")
}

// Java: ProcessCommonJSModulesTest.testGoogProvideUnaffected()
func (c *ProcessCommonJSModulesTest) TestGoogProvideUnaffected() {
	c.TestModules("test.js", "goog.provide('foo'); foo = 123;", "goog.provide('foo'); foo = 123;")
}

// Java: ProcessCommonJSModulesTest.testTopModuleCallNotRewritten()
func (c *ProcessCommonJSModulesTest) TestTopModuleCallNotRewritten() {
	// This test the case when some JS doesn't use common js but uses top-level module calls.
	// For example in Jasmine test framwork.
	c.TestModules("test.js", "module('foo.bar');", "module('foo.bar');")
}

// Java: ProcessCommonJSModulesTest.testEsModuleExportsNotRewritten()
func (c *ProcessCommonJSModulesTest) TestEsModuleExportsNotRewritten() {
	// There was a bug where logic in the ProcessCommonJsModules incorrectly split this double
	// declaration.
	code := "export var foo = 1, bar = 2;"
	c.TestModules("test.js", code, code)
}

// Java: ProcessCommonJSModulesTest.testRequireDot()
func (c *ProcessCommonJSModulesTest) TestRequireDot() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(Compiler_JoinPathParts("mod", "index.js"), "module.exports = {};"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("mod", "sub.js"),
				"var name = require('.');\n"+
					"(function() { let foo = name; foo(); })();\n")),
		c.Expected(
			SourceFile_FromCode(
				Compiler_JoinPathParts("mod", "index.js"),
				"/** @const */ var module$mod$index = {/** @const */ default: {}};"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("mod", "sub.js"),
				"var name = module$mod$index.default;\n"+
					"(function() { let foo = module$mod$index.default; foo(); })();\n")))
}

// Java: ProcessCommonJSModulesTest.testRequireDotDot()
func (c *ProcessCommonJSModulesTest) TestRequireDotDot() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(Compiler_JoinPathParts("mod", "index.js"), "module.exports = {};"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("mod", "sub", "sub.js"),
				"var name = require('..');\n"+
					"(function() { let foo = name; foo(); })();\n")),
		c.Expected(
			SourceFile_FromCode(
				Compiler_JoinPathParts("mod", "index.js"),
				"/** @const */ var module$mod$index = {/** @const */ default: {}};"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("mod", "sub", "sub.js"),
				"var name = module$mod$index.default;\n"+
					"(function() { let foo = module$mod$index.default; foo(); })();\n")))
}

// Java: ProcessCommonJSModulesTest.testFindModuleDependencies()
func (c *ProcessCommonJSModulesTest) TestFindModuleDependencies() {
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetProcessCommonJSModules(true)
	options.SetModuleResolutionMode(ModuleLoader_ResolutionMode_NODE)
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "")},
		[]*SourceFile{SourceFile_FromCode("mod.js", "var x = 1;")},
		options)
	root := compiler.ParseInputs()
	script := root.GetSecondChild().GetFirstChild()
	pass :=
		NewFindModuleDependencies(compiler, true, true, jutil.NewLinkedHashMap[string, string]())
	pass.Process(script)
	input := compiler.GetInput(script.GetInputId())
	AssertThat(input).IsNotNull()
	AssertThat(input.GetKnownProvides()).IsNotNull()
}

// Port of com/google/javascript/jscomp/Es6RewriteModulesWithGoogInteropTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/Es6RewriteModulesWithGoogInteropTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Unit tests for Es6RewriteModules that test interop with closure files.
type Es6RewriteModulesWithGoogInteropTest struct {
	CompilerTestCaseBase
	runTypeChecker bool
}

// Java: Es6RewriteModulesWithGoogInteropTest.CLOSURE_BASE
var es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE = SourceFile_FromCode("closure_base.js", NewTestExternsBuilder().AddClosureExterns().Build())

// Java: Es6RewriteModulesWithGoogInteropTest.CLOSURE_PROVIDE
var es6RewriteModulesWithGoogInteropTest_CLOSURE_PROVIDE = SourceFile_FromCode(
	"closure_provide.js",
	`goog.provide('closure.provide'); closure.provide = class {};
closure.provide.Type = class {};
closure.provide.a = class {};
closure.provide.a.d = class {};
closure.provide.b = class {};
`)

// Java: Es6RewriteModulesWithGoogInteropTest.CLOSURE_MODULE
var es6RewriteModulesWithGoogInteropTest_CLOSURE_MODULE = SourceFile_FromCode(
	"closure_module.js",
	`goog.module('closure.module');
exports = class {};
exports.x = class {};
exports.y = class {};
`)

// Java: Es6RewriteModulesWithGoogInteropTest.CLOSURE_MODULE_DESTRUCTURE
var es6RewriteModulesWithGoogInteropTest_CLOSURE_MODULE_DESTRUCTURE = SourceFile_FromCode(
	"closure_destructured_module.js",
	`goog.module('closure.destructured.module');
exports.Bar = 0;
exports.y = 0;
exports.b = class {};
exports.a = class {};
exports.a.d = class {};
`)

// Java: Es6RewriteModulesWithGoogInteropTest.CLOSURE_LEGACY_MODULE
var es6RewriteModulesWithGoogInteropTest_CLOSURE_LEGACY_MODULE = SourceFile_FromCode(
	"closure_legacy_destructured_module.js",
	`goog.module('closure.legacy.module');
goog.module.declareLegacyNamespace();
exports = class {};
exports.x = class {};
exports.y = class {};
`)

// Java: Es6RewriteModulesWithGoogInteropTest.CLOSURE_LEGACY_DESTRUCTURED_MODULE
var es6RewriteModulesWithGoogInteropTest_CLOSURE_LEGACY_DESTRUCTURED_MODULE = SourceFile_FromCode(
	"closure_legacy_module.js",
	`goog.module('closure.legacy.destructured.module');
goog.module.declareLegacyNamespace();
exports.Bar = 0;
exports.y = 0;
exports.b = class {};
exports.a = class {};
 exports.a.d = class {};
`)

func NewEs6RewriteModulesWithGoogInteropTest() *Es6RewriteModulesWithGoogInteropTest {
	c := &Es6RewriteModulesWithGoogInteropTest{}
	c.runTypeChecker = true
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestEs6RewriteModulesWithGoogInteropTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteModulesWithGoogInteropTest() })
}

// Java: Es6RewriteModulesWithGoogInteropTest.setUp()
func (c *Es6RewriteModulesWithGoogInteropTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// ECMASCRIPT5 to trigger module processing after parsing.
	c.EnableTypeInfoValidation()
	c.EnableCreateModuleMap()
}

// Java: Es6RewriteModulesWithGoogInteropTest.getOptions()
func (c *Es6RewriteModulesWithGoogInteropTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	// ECMASCRIPT5 to Trigger module processing after parsing.
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_ERROR)

	return options
}

// Java: Es6RewriteModulesWithGoogInteropTest.getProcessor(Compiler)
func (c *Es6RewriteModulesWithGoogInteropTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		var globalTypedScope *TypedScope
		if c.runTypeChecker {
			rai :=
				NewSemanticReverseAbstractInterpreter(compiler.GetTypeRegistry())
			globalTypedScope =
				jutil.CheckNotNull(
					NewTypeCheck(compiler, rai, compiler.GetTypeRegistry()).
						ProcessForTesting(externs, root))
			compiler.SetTypeCheckingHasRun(true)
		} else {
			globalTypedScope = nil
		}

		NewEs6RewriteModules(
			compiler,
			compiler.GetModuleMetadataMap(),
			compiler.GetModuleMap(),
			/* preprocessorSymbolTable= */ nil,
			globalTypedScope).
			Process(externs, root)
	})
}

// Java: Es6RewriteModulesWithGoogInteropTest.testModules(String, String)
func (c *Es6RewriteModulesWithGoogInteropTest) TestModules(input string, expected string) {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"closure_base.js", NewTestExternsBuilder().AddClosureExterns().Build()),
			es6RewriteModulesWithGoogInteropTest_CLOSURE_PROVIDE,
			es6RewriteModulesWithGoogInteropTest_CLOSURE_MODULE,
			es6RewriteModulesWithGoogInteropTest_CLOSURE_MODULE_DESTRUCTURE,
			es6RewriteModulesWithGoogInteropTest_CLOSURE_LEGACY_MODULE,
			es6RewriteModulesWithGoogInteropTest_CLOSURE_LEGACY_DESTRUCTURED_MODULE,
			SourceFile_FromCode("testcode", input)),
		c.Expected(
			SourceFile_FromCode(
				"closure_base.js", NewTestExternsBuilder().AddClosureExterns().Build()),
			es6RewriteModulesWithGoogInteropTest_CLOSURE_PROVIDE,
			es6RewriteModulesWithGoogInteropTest_CLOSURE_MODULE,
			es6RewriteModulesWithGoogInteropTest_CLOSURE_MODULE_DESTRUCTURE,
			es6RewriteModulesWithGoogInteropTest_CLOSURE_LEGACY_MODULE,
			es6RewriteModulesWithGoogInteropTest_CLOSURE_LEGACY_DESTRUCTURED_MODULE,
			SourceFile_FromCode("testcode", expected)))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testModulesError(String, DiagnosticType)
func (c *Es6RewriteModulesWithGoogInteropTest) TestModulesError(input string, error_ *DiagnosticType) {
	c.TestError(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_PROVIDE,
			es6RewriteModulesWithGoogInteropTest_CLOSURE_MODULE,
			es6RewriteModulesWithGoogInteropTest_CLOSURE_LEGACY_MODULE,
			SourceFile_FromCode("testcode", input)),
		c.Error(error_))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testClosureFilesUnchanged()
func (c *Es6RewriteModulesWithGoogInteropTest) TestClosureFilesUnchanged() {
	c.TestSame(c.Srcs(es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE, es6RewriteModulesWithGoogInteropTest_CLOSURE_PROVIDE, es6RewriteModulesWithGoogInteropTest_CLOSURE_MODULE, es6RewriteModulesWithGoogInteropTest_CLOSURE_LEGACY_MODULE))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireInNonEs6ModuleUnchanged()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireInNonEs6ModuleUnchanged() {
	c.TestSame(
		c.Srcs(NewTestExternsBuilder().AddClosureExterns().Build(), "goog.require('foo.bar');"))
	c.TestSame(
		c.Srcs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			"var bar = goog.require('foo.bar');"))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeInNonEs6ModuleUnchanged()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeInNonEs6ModuleUnchanged() {
	c.TestSame(
		c.Srcs(NewTestExternsBuilder().AddClosureExterns().Build(), "goog.requireType('foo.bar');"))
	c.TestSame(
		c.Srcs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			"var bar = goog.requireType('foo.bar');"))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testForwardDeclareInNonEs6ModuleUnchanged()
func (c *Es6RewriteModulesWithGoogInteropTest) TestForwardDeclareInNonEs6ModuleUnchanged() {
	c.TestSame(
		c.Srcs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			"goog.forwardDeclare('foo.bar');"))
	c.TestSame(
		c.Srcs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			"var bar = goog.forwardDeclare('foo.bar');"))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireInEs6ModuleDoesNotExistIsError()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireInEs6ModuleDoesNotExistIsError() {
	c.TestError(
		c.Srcs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			"export var x; goog.require('foo.bar');"),
		c.Error(ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeInEs6ModuleDoesNotExistIsError()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeInEs6ModuleDoesNotExistIsError() {
	c.TestError(
		c.Srcs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			"export var x; goog.requireType('foo.bar');"),
		c.Error(ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogModuleGetNonStringIsError()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogModuleGetNonStringIsError() {
	c.TestError(
		`const y = goog.module.get();
export {};
`,
		ClosurePrimitiveErrors_INVALID_GET_NAMESPACE)
	c.TestError(
		`const y = goog.module.get(0);
export {};
`,
		ClosurePrimitiveErrors_INVALID_GET_NAMESPACE)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireForProvide()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireForProvide() {
	c.TestModules(
		`const y = goog.require('closure.provide');
use(y, y.Type)
export {};
`,
		`use(closure.provide, closure.provide.Type);
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForProvide()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForProvide() {
	c.TestModules(
		`const y = goog.requireType('closure.provide');
/**
 * @param {y} a
 * @param {y.Type} b
*/
function f(a, b) {}
export {};
`,
		`/**
 * @param {closure.provide} a
 * @param {closure.provide.Type} b
 */
function f$$module$testcode(a, b) {}
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireForProvideWithDestructure()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireForProvideWithDestructure() {
	c.TestModules(
		`const {a, b:c} = goog.require('closure.provide');
use(a, a.d, c)
export {};
`,
		`use(closure.provide.a, closure.provide.a.d, closure.provide.b);
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForProvideWithDestructure()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForProvideWithDestructure() {
	c.TestModules(
		`const {a, b:c} = goog.requireType('closure.provide');
/**
 * @param {a} x
 * @param {a.d} y
 * @param {c} z
*/
function f(x, y, z) {}
export {};
`,
		`/**
 * @param {closure.provide.a} x
 * @param {closure.provide.a.d} y
 * @param {closure.provide.b} z
 */
function f$$module$testcode(x, y, z) {}
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireForGoogModule()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireForGoogModule() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.TestModules(
		`const y = goog.require('closure.module');
use(y, y.x)
export {};
`,
		`use(module$exports$closure$module, module$exports$closure$module.x);
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForGoogModule()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForGoogModule() {
	c.IgnoreWarnings(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
	c.TestModules(
		`const y = goog.requireType('closure.module');
/**
 * @param {y} a
 * @param {y.x} b
*/
function f(a, b) {}
export {};
`,
		`/**
 * @param {module$exports$closure$module} a
 * @param {module$exports$closure$module.x} b
 */
function f$$module$testcode(a, b) {}
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireForGoogModuleWithDestructure()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireForGoogModuleWithDestructure() {
	c.TestModules(
		`const {a, b:c} = goog.require('closure.destructured.module');
use(a, a.d, c)
export {};
`,
		`use(module$exports$closure$destructured$module.a, module$exports$closure$destructured$module.a.d,
  module$exports$closure$destructured$module.b);
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForGoogModuleWithDestructure()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForGoogModuleWithDestructure() {
	c.TestModules(
		`const {a, b:c} = goog.requireType('closure.destructured.module');
/**
 * @param {a} x
 * @param {a.d} y
 * @param {c} z
*/
function f(x, y, z) {}
export {};
`,
		`/**
 * @param {module$exports$closure$destructured$module.a} x
 * @param {module$exports$closure$destructured$module.a.d} y
 * @param {module$exports$closure$destructured$module.b} z
 */
function f$$module$testcode(x, y, z) {}
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireForLegacyGoogModule()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireForLegacyGoogModule() {
	c.TestModules(
		`const y = goog.require('closure.legacy.module');
use(y, y.x)
export {};
`,
		`use(closure.legacy.module, closure.legacy.module.x);
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForLegacyGoogModule()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForLegacyGoogModule() {
	c.TestModules(
		`const y = goog.requireType('closure.legacy.module');
/**
 * @param {y} a
 * @param {y.x} b
*/
function f(a, b) {}
export {};
`,
		`/**
 * @param {closure.legacy.module} a
 * @param {closure.legacy.module.x} b
 */
function f$$module$testcode(a, b) {}
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireForLegacyGoogModuleWithDestructure()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireForLegacyGoogModuleWithDestructure() {
	c.TestModules(
		`const {a, b:c} = goog.require('closure.legacy.destructured.module');
use(a, a.d, c)
export {};
`,
		`use(closure.legacy.destructured.module.a, closure.legacy.destructured.module.a.d, closure.legacy.destructured.module.b);
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForLegacyGoogModuleWithDestructure()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForLegacyGoogModuleWithDestructure() {
	c.TestModules(
		`const {a, b:c} = goog.requireType('closure.legacy.destructured.module');
/**
 * @param {a} x
 * @param {a.d} y
 * @param {c} z
*/
function f(x, y, z) {}
export {};
`,
		`/**
 * @param {closure.legacy.destructured.module.a} x
 * @param {closure.legacy.destructured.module.a.d} y
 * @param {closure.legacy.destructured.module.b} z
 */
function f$$module$testcode(x, y, z) {}
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogModuleGetForGoogModule()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogModuleGetForGoogModule() {
	c.TestModules(
		"function foo() { const y = goog.module.get('closure.module'); } export {};",
		`function foo$$module$testcode() { const y = module$exports$closure$module; }
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogModuleGetForLegacyGoogModule()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogModuleGetForLegacyGoogModule() {
	c.TestModules(
		"function foo(){ const y = goog.module.get('closure.legacy.module'); }\nexport {};",
		`function foo$$module$testcode(){ const y = closure.legacy.module; }
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogModuleGetForEs6Module()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogModuleGetForEs6Module() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export {}; goog.declareModuleId('es6');"),
			SourceFile_FromCode(
				"closure.js",
				"goog.module('my.module'); function f() { const y = goog.module.get('es6'); }")),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "/** @const */ var module$es6 = {};"),
			SourceFile_FromCode(
				"closure.js", "goog.module('my.module'); function f() { const y = module$es6; }")))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogModuleGetForEs6ModuleDeclaresNamespace()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogModuleGetForEs6ModuleDeclaresNamespace() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export let y; goog.declareModuleId('es6');"),
			SourceFile_FromCode(
				"closure.js",
				"goog.module('my.module'); function f() { return goog.module.get('es6').y; }")),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`let y$$module$es6;
/** @const */ var module$es6 = {};
/** @const */ module$es6.y = y$$module$es6;
`),
			SourceFile_FromCode(
				"closure.js", "goog.module('my.module'); function f() { return module$es6.y; }")))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testDeclareModuleId()
func (c *Es6RewriteModulesWithGoogInteropTest) TestDeclareModuleId() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export var x; goog.declareModuleId('my.es6');"),
			SourceFile_FromCode("goog.js", "const es6 = goog.require('my.es6'); use(es6, es6.x);")),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`var x$$module$es6;/** @const */ var module$es6={};
/** @const */ module$es6.x=x$$module$es6;
`),
			SourceFile_FromCode("goog.js", "const es6 = module$es6; use(es6, es6.x);")))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireForDeclareModuleId()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireForDeclareModuleId() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export var x; goog.declareModuleId('my.es6');"),
			SourceFile_FromCode("goog.js", "const es6 = goog.require('my.es6'); use(es6, es6.x);")),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`var x$$module$es6;/** @const */ var module$es6={};
/** @const */ module$es6.x=x$$module$es6;
`),
			SourceFile_FromCode("goog.js", "const es6 = module$es6; use(es6, es6.x);")))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireForDeclareModuleId_inGoogLoadModule()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireForDeclareModuleId_inGoogLoadModule() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export var x; goog.declareModuleId('my.es6');"),
			SourceFile_FromCode(
				"goog.js",
				`goog.loadModule(function(exports) {
goog.module('closure.module');
const es6 = goog.require('my.es6');
use(es6, es6.x);
return exports;
});
`)),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`var x$$module$es6;/** @const */ var module$es6={};
/** @const */ module$es6.x=x$$module$es6;
`),
			SourceFile_FromCode(
				"goog.js",
				`goog.loadModule(function(exports) {
goog.module('closure.module');
const es6 = module$es6;
use(es6, es6.x);
return exports;
});
`)))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForDeclareModuleId()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForDeclareModuleId() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export class x {}; goog.declareModuleId('my.es6');"),
			SourceFile_FromCode(
				"goog.js",
				`goog.module('bar')
const es6 = goog.requireType('my.es6');
/** @param {es6.x} a */
function f(a) {}
`)),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`class x$$module$es6 {}
/** @const */ var module$es6={};
/** @const */ module$es6.x = x$$module$es6;
`),
			SourceFile_FromCode(
				"goog.js",
				`goog.module('bar')
/**
 * @param {module$es6.x} a
 */
function f(a) {}
`)))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForDeclareModuleId_inLoadModule()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForDeclareModuleId_inLoadModule() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export class x {}; goog.declareModuleId('my.es6');"),
			SourceFile_FromCode(
				"goog.js",
				`goog.loadModule(function(exports) {
goog.module('bar')
const es6 = goog.requireType('my.es6');
/** @param {es6.x} a */
function f(a) {}
return exports;
});
`)),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`class x$$module$es6 {}
/** @const */ var module$es6={};
/** @const */ module$es6.x = x$$module$es6;
`),
			SourceFile_FromCode(
				"goog.js",
				`goog.loadModule(function(exports) {
goog.module('bar')
/**
 * @param {module$es6.x} a
 */
function f(a) {}
return exports;
});
`)))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForDeclareModuleId_lateLoaded()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForDeclareModuleId_lateLoaded() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"goog.js",
				`goog.module('bar')
const es6 = goog.requireType('my.es6');
/** @param {es6.x} a */
function f(a) {}
`),
			SourceFile_FromCode("es6.js", "export class x {}; goog.declareModuleId('my.es6');")),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"goog.js",
				`goog.module('bar')
/**
 * @param {module$es6.x} a
 */
function f(a) {}
`),
			SourceFile_FromCode(
				"es6.js",
				`class x$$module$es6 {}
/** @const */ var module$es6={};
/** @const */ module$es6.x = x$$module$es6;
`)))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireForDeclareModuleIdWithDestructure()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireForDeclareModuleIdWithDestructure() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export var x, z; goog.declareModuleId('my.es6');"),
			SourceFile_FromCode("goog.js", "const {x, z: y} = goog.require('my.es6'); use(x, y);")),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`var x$$module$es6, z$$module$es6;
/** @const */ var module$es6={};
/** @const */ module$es6.x=x$$module$es6;
/** @const */ module$es6.z=z$$module$es6;
`),
			SourceFile_FromCode(
				"goog.js",
				`const x = module$es6.x;
const y = module$es6.z;
use(x, y);
`)))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForDeclareModuleIdWithDestructure()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForDeclareModuleIdWithDestructure() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js", "export class x {} export class z {} goog.declareModuleId('my.es6');"),
			SourceFile_FromCode(
				"goog.js",
				`const {x, z: y} = goog.requireType('my.es6');
/**
 * @param {x} a
 * @param {y} b
 */
function f(a, b) {}
`)),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`class x$$module$es6 {}
class z$$module$es6 {}
/** @const */ var module$es6={};
/** @const */ module$es6.x=x$$module$es6;
/** @const */ module$es6.z=z$$module$es6;
`),
			SourceFile_FromCode(
				"goog.js",
				`/**
 * @param {module$es6.x} a
 * @param {module$es6.z} b
 */
function f(a, b) {}
`)))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireLhsNonConstIsError()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireLhsNonConstIsError() {
	c.TestModulesError(
		`var bar = goog.require('closure.provide');
export var x;
`,
		Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST)

	c.TestModulesError(
		`export var x;
var bar = goog.require('closure.provide');
`,
		Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST)

	c.TestModulesError(
		`export {};
var {foo, bar} = goog.require('closure.provide');
use(foo, bar);
`,
		Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST)

	c.TestModulesError(
		`export {};
let {foo, bar} = goog.require('closure.provide');
use(foo, bar);
`,
		Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeLhsNonConstIsError()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeLhsNonConstIsError() {
	c.TestModulesError(
		`var bar = goog.requireType('closure.provide');
export var x;
`,
		Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST)

	c.TestModulesError(
		`export var x;
var bar = goog.requireType('closure.provide');
`,
		Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST)

	c.TestModulesError(
		`export {};
var {foo, bar} = goog.requireType('closure.provide');
/**
 * @param {foo} x
 * @param {bar} y
 */
function f(x, y) {}
`,
		Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST)

	c.TestModulesError(
		`export {};
let {foo, bar} = goog.requireType('closure.provide');
/**
 * @param {foo} x
 * @param {bar} y
 */
function f(x, y) {}
`,
		Es6RewriteModules_LHS_OF_GOOG_REQUIRE_MUST_BE_CONST)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForNonEs6LhsNonConst1()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForNonEs6LhsNonConst1() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode("es6.js", "export var x; goog.declareModuleId('es6');"),
			SourceFile_FromCode(
				"closure.js",
				`goog.module('my.module');
var es6 = goog.requireType('es6');
`)),
		c.Error(Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeForNonEs6LhsNonConst2()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeForNonEs6LhsNonConst2() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode("es6.js", "export var x; goog.declareModuleId('es6');"),
			SourceFile_FromCode(
				"closure.js",
				`goog.module('my.module');
let { x } = goog.requireType('es6');
`)),
		c.Error(Es6RewriteModules_REQUIRE_TYPE_FOR_ES6_SHOULD_BE_CONST))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testNamespaceImports()
func (c *Es6RewriteModulesWithGoogInteropTest) TestNamespaceImports() {
	c.DisableTypeInfoValidation()
	c.runTypeChecker = false // typechecker does not support these namespace imports.
	c.TestModules(
		`import Foo from 'goog:closure.provide';
use(Foo);
`,
		"use(closure.provide); /** @const */ var module$testcode = {};")

	c.TestModules(
		`import {x, y} from 'goog:closure.provide';
use(x);
use(y);
`,
		"use(closure.provide.x); use(closure.provide.y);/** @const */ var module$testcode = {};")

	c.TestModules(
		`import Foo from 'goog:closure.provide';
/** @type {Foo} */ var foo = new Foo();
`,
		`/** @type {closure.provide} */
var foo$$module$testcode = new closure.provide();
/** @const */ var module$testcode = {};
`)

	c.TestModules(
		`import Foo from 'goog:closure.module';
/** @type {Foo.Bar} */ var foo = new Foo.Bar();
`,
		`/** @type {module$exports$closure$module.Bar} */
var foo$$module$testcode = new module$exports$closure$module.Bar();
/** @const */ var module$testcode = {};
`)

	c.TestModules(
		`import {Bar} from 'goog:closure.destructured.module';
/** @type {Bar} */ var foo = new Bar();
`,
		`/** @type {module$exports$closure$destructured$module.Bar} */
var foo$$module$testcode = new module$exports$closure$destructured$module.Bar();
/** @const */ var module$testcode = {};
`)

	c.TestModules(
		`import Foo from 'goog:closure.legacy.module';
/** @type {Foo.Bar} */ var foo = new Foo.Bar();
`,
		`/** @type {closure.legacy.module.Bar} */
var foo$$module$testcode = new closure.legacy.module.Bar();
/** @const */ var module$testcode = {};
`)

	c.TestModules(
		`import {Bar} from 'goog:closure.legacy.destructured.module';
/** @type {Bar} */ var foo = new Bar();
`,
		`/** @type {closure.legacy.destructured.module.Bar} */
var foo$$module$testcode = new closure.legacy.destructured.module.Bar();
/** @const */ var module$testcode = {};
`)

	c.TestModulesError(
		"import * as Foo from 'goog:closure.provide';", EsModuleProcessor_NAMESPACE_IMPORT_CANNOT_USE_STAR)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireMustBeModuleScope()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireMustBeModuleScope() {
	c.TestModulesError(
		"{ goog.require('closure.provide'); } export {}", ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeMustBeModuleScope()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeMustBeModuleScope() {
	c.TestModulesError(
		"{ goog.requireType('closure.provide'); } export {}", ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogModuleGetCannotBeModuleHoistScope()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogModuleGetCannotBeModuleHoistScope() {
	c.TestModulesError("goog.module.get('closure.module'); export {}", ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET)
	c.TestModulesError(
		"{ goog.module.get('closure.module'); } export {}", ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testExportSpecGoogRequire()
func (c *Es6RewriteModulesWithGoogInteropTest) TestExportSpecGoogRequire() {
	c.TestModules(
		`const {a} = goog.require('closure.provide');
export {a};
`,
		`/** @const */ var module$testcode={};
/** @const */ module$testcode.a = closure.provide.a;
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testRequireAndStoreGlobalUnqualifiedProvide()
func (c *Es6RewriteModulesWithGoogInteropTest) TestRequireAndStoreGlobalUnqualifiedProvide() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("provide.js", "goog.provide('foo')"),
			SourceFile_FromCode(
				"testcode",
				`const foo = goog.require('foo');
use(foo);
export {};
`)),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("provide.js", "goog.provide('foo')"),
			SourceFile_FromCode(
				"testcode",
				`use(foo);
/** @const */ var module$testcode={};
`)))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireAnnotationIsRenamed()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireAnnotationIsRenamed() {
	c.TestModules(
		`const {Type} = goog.require('closure.provide');
export let /** !Type */ value;
`,
		`let /** !closure.provide.Type */ value$$module$testcode;
/** @const */ var module$testcode={};
/** @const */ module$testcode.value = value$$module$testcode;
`)

	c.TestModules(
		`const Type = goog.require('closure.provide');
export let /** !Type */ value;
`,
		`let /** !closure.provide */ value$$module$testcode;
/** @const */ var module$testcode={};
/** @const */ module$testcode.value = value$$module$testcode;
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeAnnotationIsRenamed()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeAnnotationIsRenamed() {
	c.TestModules(
		`const {Type} = goog.requireType('closure.provide');
export let /** !Type */ value;
`,
		`let /** !closure.provide.Type */ value$$module$testcode;
/** @const */ var module$testcode={};
/** @const */ module$testcode.value = value$$module$testcode;
`)

	c.TestModules(
		`const Type = goog.requireType('closure.provide');
export let /** !Type */ value;
`,
		`let /** !closure.provide */ value$$module$testcode;
/** @const */ var module$testcode={};
/** @const */ module$testcode.value = value$$module$testcode;
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireCorrectAnnotationIsRenamed()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireCorrectAnnotationIsRenamed() {
	c.TestModules(
		`const {Type} = goog.require('closure.provide');
export let /** !Type */ value;
function foo() {
  class Type {}
  let /** !Type */ value;
}
`,
		`let /** !closure.provide.Type */ value$$module$testcode;
function foo$$module$testcode() {
  class Type {}
  let /** !Type */ value;
}
/** @const */ var module$testcode={};
/** @const */ module$testcode.value = value$$module$testcode;
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeCorrectAnnotationIsRenamed()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeCorrectAnnotationIsRenamed() {
	c.TestModules(
		`const {Type} = goog.requireType('closure.provide');
export let /** !Type */ value;
function foo() {
  class Type {}
  let /** !Type */ value;
}
`,
		`let /** !closure.provide.Type */ value$$module$testcode;
function foo$$module$testcode() {
  class Type {}
  let /** !Type */ value;
}
/** @const */ var module$testcode={};
/** @const */ module$testcode.value = value$$module$testcode;
`)
}

// Java: Es6RewriteModulesWithGoogInteropTest.testForwardDeclareEs6Module()
func (c *Es6RewriteModulesWithGoogInteropTest) TestForwardDeclareEs6Module() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js", "export {}; goog.declareModuleId('es6'); export class Type {}"),
			SourceFile_FromCode(
				"closure.js",
				`goog.module('my.module');
const alias = goog.forwardDeclare('es6');
let /** !alias.Type */ x;
alias;
`)),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`class Type$$module$es6 {}
/** @const */ var module$es6 = {};
/** @const */ module$es6.Type = Type$$module$es6;
`),
			SourceFile_FromCode(
				"closure.js",
				`goog.module('my.module');
let /** !module$es6.Type */ x;
alias;
`)))

	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js", "export {}; goog.declareModuleId('es6'); export class Type {}"),
			SourceFile_FromCode(
				"closure.js",
				`goog.module('my.module');
goog.forwardDeclare('es6');
let /** !es6.Type */ x;
es6;
`)),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))

	c.TestError(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export {}; goog.declareModuleId('es6');"),
			SourceFile_FromCode(
				"closure.js",
				`goog.module('my.module');
let alias = goog.forwardDeclare('es6');
let /** !alias.Type */ x;
alias = goog.modle.get('es6');
`)),
		c.Error(Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testWarnAboutRequireEs6FromEs6()
func (c *Es6RewriteModulesWithGoogInteropTest) TestWarnAboutRequireEs6FromEs6() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("first.js", "goog.declareModuleId('first'); export {};"),
			SourceFile_FromCode(
				"second.js", "const first = goog.require('first'); export {first};")),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("first.js", "/** @const */ var module$first = {};"),
			SourceFile_FromCode(
				"second.js",
				`const first$$module$second = module$first;
/** @const */ var module$second = {};
/** @const */ module$second.first = first$$module$second;
`)),
		c.Warning(Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE))

	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("first.js", "goog.declareModuleId('no.alias'); export {};"),
			SourceFile_FromCode("second.js", "goog.require('no.alias'); export {};")),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("first.js", "/** @const */ var module$first = {};"),
			SourceFile_FromCode("second.js", "/** @const */ var module$second = {};")),
		c.Warning(Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogRequireTypeEs6ModuleInEs6Module()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogRequireTypeEs6ModuleInEs6Module() {
	// TODO(tjgq): Make these warn once there's a requireType equivalent for ES6 modules.

	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "goog.declareModuleId('es6'); export class Type {};"),
			SourceFile_FromCode(
				"requiretype.js",
				`export {}
const alias = goog.requireType('es6');
let /** !alias.Type */ x;
`)),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`class Type$$module$es6 {}
/** @const */ var module$es6 = {};
/** @const */ module$es6.Type = Type$$module$es6;
`),
			SourceFile_FromCode(
				"requiretype.js",
				`let /** !module$es6.Type */ x$$module$requiretype;
/** @const */ var module$requiretype = {};
`)))

	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "goog.declareModuleId('es6'); export class Type {}"),
			SourceFile_FromCode(
				"requiretype.js",
				`export {};
goog.requireType('es6');
let /** !es6.Type */ x;
`)),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testGoogModuleGetEs6ModuleInEs6Module()
func (c *Es6RewriteModulesWithGoogInteropTest) TestGoogModuleGetEs6ModuleInEs6Module() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("first.js", "goog.declareModuleId('first'); export let x;"),
			SourceFile_FromCode(
				"second.js", "export function foo() { return goog.module.get('first').x; }")),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"first.js",
				`let x$$module$first;
/** @const */ var module$first = {};
/** @const */ module$first.x = x$$module$first;
`),
			SourceFile_FromCode(
				"second.js",
				`function foo$$module$second() {
  return module$first.x;
}
/** @const */ var module$second = {};
/** @const */ module$second.foo = foo$$module$second;
`)))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testForwardDeclareEs6ModuleInEs6Module()
func (c *Es6RewriteModulesWithGoogInteropTest) TestForwardDeclareEs6ModuleInEs6Module() {
	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export class Type {} goog.declareModuleId('es6');"),
			SourceFile_FromCode(
				"forwarddeclare.js",
				`export {}
const alias = goog.forwardDeclare('es6');
let /** !alias.Type */ x;
alias;
`)),
		c.Expected(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode(
				"es6.js",
				`class Type$$module$es6 {}
/** @const */ var module$es6 = {};
/** @const */ module$es6.Type = Type$$module$es6;
`),
			SourceFile_FromCode(
				"forwarddeclare.js",
				`let /** !module$es6.Type */ x$$module$forwarddeclare;
alias;
/** @const */ var module$forwarddeclare = {};
`)))

	c.Test(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export class Type {} goog.declareModuleId('es6');"),
			SourceFile_FromCode(
				"forwarddeclare.js",
				`export {};
goog.forwardDeclare('es6');
let /** !es6.Type */ x;
es6;
`)),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))

	c.TestError(
		c.Srcs(
			es6RewriteModulesWithGoogInteropTest_CLOSURE_BASE,
			SourceFile_FromCode("es6.js", "export class Type {} goog.declareModuleId('es6');"),
			SourceFile_FromCode(
				"closure.js",
				`export {};
let alias = goog.forwardDeclare('es6');
let /** !alias.Type */ x;
alias = goog.modle.get('es6');
`)),
		c.Error(Es6RewriteModules_FORWARD_DECLARE_FOR_ES6_SHOULD_BE_CONST))
}

// Java: Es6RewriteModulesWithGoogInteropTest.testMissingRequireAssumesGoogProvide()
func (c *Es6RewriteModulesWithGoogInteropTest) TestMissingRequireAssumesGoogProvide() {
	c.IgnoreWarnings(ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE, TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)

	c.Test(
		c.Srcs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			`const missing = goog.require('is.missing');
use(missing, missing.x);
export {};
`),
		c.Expected(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			`use(is.missing, is.missing.x);
/** @const */ var module$testcode1 = {};
`))

	c.Test(
		c.Srcs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			`const {prop} = goog.require('is.missing');
use(prop, prop.x);
export {};
`),
		c.Expected(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			`use(is.missing.prop, is.missing.prop.x);
/** @const */ var module$testcode1 = {};
`))

	c.Test(
		c.Srcs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			`import missing, {y} from 'goog:is.missing';
use(missing, missing.x, y);
export {};
`),
		c.Expected(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			`use(is.missing, is.missing.x, is.missing.y);
/** @const */ var module$testcode1 = {};
`))
}

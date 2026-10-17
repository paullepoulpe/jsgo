// Port of com/google/javascript/jscomp/ClosureRewriteModuleTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/ClosureRewriteModuleTest.java)

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

// Unit tests for ClosureRewriteModule
type ClosureRewriteModuleTest struct {
	CompilerTestCaseBase

	preserveClosurePrimitives bool
	allowMissingSources       bool
	runTypeCheck              bool
}

// Java: ClosureRewriteModuleTest()
func NewClosureRewriteModuleTest() *ClosureRewriteModuleTest {
	c := &ClosureRewriteModuleTest{runTypeCheck: true}
	c.InitCompilerTestCaseBaseWithExterns(c, NewTestExternsBuilder().AddClosureExterns().AddPromise().AddConsole().Build())
	return c
}

func TestClosureRewriteModuleTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewClosureRewriteModuleTest() })
}

// Java: ClosureRewriteModuleTest.getProcessor(Compiler)
func (c *ClosureRewriteModuleTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, main *rhino.Node) {
		var globalTypedScope *TypedScope
		if c.runTypeCheck {
			var rai ReverseAbstractInterpreter = NewSemanticReverseAbstractInterpreter(compiler.GetTypeRegistry())
			globalTypedScope =
				jutil.CheckNotNull(
					NewTypeCheck(compiler, rai, compiler.GetTypeRegistry()).
						ProcessForTesting(externs, main))
			compiler.SetTypeCheckingHasRun(true)
		}

		NewClosureRewriteModule(compiler, nil, globalTypedScope).Process(externs, main)
	})
}

// Java: ClosureRewriteModuleTest.setUp()
func (c *ClosureRewriteModuleTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.runTypeCheck = true
	c.preserveClosurePrimitives = false
	c.allowMissingSources = false
	c.EnableCreateModuleMap()
	c.EnableTypeInfoValidation()
}

// Java: ClosureRewriteModuleTest.getOptions()
func (c *ClosureRewriteModuleTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	options.SetPreserveClosurePrimitives(c.preserveClosurePrimitives)
	options.SetWarningLevel(DiagnosticGroups_MISSING_PROVIDE, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_WEAK_MODULE_GET, CheckLevel_ERROR)
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	if c.allowMissingSources {
		options.SetWarningLevel(DiagnosticGroups_MISSING_SOURCES_WARNINGS, CheckLevel_OFF)
	}
	options.SetPrettyPrint(true)
	return options
}

// Java: ClosureRewriteModuleTest.testEmpty()
func (c *ClosureRewriteModuleTest) TestEmpty() {
	c.TestSame("")
}

// Java: ClosureRewriteModuleTest.testProvide()
func (c *ClosureRewriteModuleTest) TestProvide() {
	c.TestSame("goog.provide('a');")
}

// Java: ClosureRewriteModuleTest.testModule()
func (c *ClosureRewriteModuleTest) TestModule() {
	c.Test("goog.module('a');", "/** @const */ var module$exports$a = {};")
}

// Java: ClosureRewriteModuleTest.testRequireModule()
func (c *ClosureRewriteModuleTest) TestRequireModule() {
	c.Test(
		c.Srcs(
			"goog.module('ns.b');",
			"goog.module('ns.c');",
			`goog.module('ns.a');
var b = goog.require('ns.b');
var c = goog.requireType('ns.c');
`),
		c.Expected(
			"/** @const */ var module$exports$ns$b = {};",
			"/** @const */ var module$exports$ns$c = {};",
			"/** @const */ var module$exports$ns$a = {};"))
}

// Java: ClosureRewriteModuleTest.testProvidedNamespaceGetsRequiredAndUsed()
func (c *ClosureRewriteModuleTest) TestProvidedNamespaceGetsRequiredAndUsed() {
	c.Test(
		c.Srcs(
			`goog.provide('goog.string.Const');
/** @constructor */ goog.string.Const = function() {};
goog.string.Const.from = function() {};
`,
			`goog.module('myModule');
var const1 = goog.require('goog.string.Const');
console.log(const1.from());
`),
		c.Expected(
			`goog.provide('goog.string.Const');
/** @constructor */ goog.string.Const = function() {};
goog.string.Const.from = function() {};
`,
			`/** @const */ var module$exports$myModule = {};
goog.require('goog.string.Const');
console.log(goog.string.Const.from());
`))
}

// Java: ClosureRewriteModuleTest.testProvidedNamespaceGetsRequiredAndUsed_destructuringImport()
func (c *ClosureRewriteModuleTest) TestProvidedNamespaceGetsRequiredAndUsed_destructuringImport() {
	c.Test(
		c.Srcs(
			`goog.provide('goog.string.Const');
/** @constructor */ goog.string.Const = function() {};
goog.string.Const.from = function() {};
`,
			`goog.module('myModule');
var {from} = goog.require('goog.string.Const');
console.log(from());
`),
		c.Expected(
			`goog.provide('goog.string.Const');
/** @constructor */ goog.string.Const = function() {};
goog.string.Const.from = function() {};
`,
			`/** @const */ var module$exports$myModule = {};
console.log((0, goog.string.Const.from)());
`))
}

// Java: ClosureRewriteModuleTest.testOutOfOrderRequireType()
func (c *ClosureRewriteModuleTest) TestOutOfOrderRequireType() {
	c.Test(
		c.Srcs(
			`goog.module('ns.a');
var b = goog.requireType('ns.b');
/** @type {b.Foo} */
let a;
`,
			`goog.module('ns.b');
exports.Foo = class {}
`),
		c.Expected(
			`/** @const */ var module$exports$ns$a = {};
/** @type {module$exports$ns$b.Foo} */ let module$contents$ns$a_a;
`,
			`/** @const */ var module$exports$ns$b = {};
/** @const */ module$exports$ns$b.Foo = class {};
`))
}

// Java: ClosureRewriteModuleTest.testRequireModuleMultivar()
func (c *ClosureRewriteModuleTest) TestRequireModuleMultivar() {
	c.Test(
		c.Srcs(
			"goog.module('ns.b');",
			"goog.module('ns.c');",
			`goog.module('ns.a');
var b = goog.require('ns.b'), c = goog.requireType('ns.c');
`),
		c.Expected(
			"/** @const */ var module$exports$ns$b = {};",
			"/** @const */ var module$exports$ns$c = {};",
			"/** @const */ var module$exports$ns$a = {};"))
}

// Java: ClosureRewriteModuleTest.testIjsModule()
func (c *ClosureRewriteModuleTest) TestIjsModule() {
	c.AllowExternsChanges()
	c.Test(
		c.Srcs(
			// .i.js files
			`/** @typeSummary */
goog.module('external1');
/** @constructor */
exports = function() {};
`,
			`/** @typeSummary */
goog.module('external2');
/** @constructor */
exports = function() {};
`,
			// source file
			`goog.module('ns.a');
var b = goog.require('external1');
var c = goog.requireType('external2');
/** @type {b} */ new b;
/** @typedef {c} */ var d;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$a = {};
/** @type {module$exports$external1} */ new module$exports$external1
/** @typedef {module$exports$external2} */ var module$contents$ns$a_d
`))
}

// Java: ClosureRewriteModuleTest.testDestructuringInsideModule()
func (c *ClosureRewriteModuleTest) TestDestructuringInsideModule() {
	// Array destrucuturing
	c.Test(
		`goog.module('a');
var [x, y, z] = foo();
`,
		`/** @const */ var module$exports$a = {};
var [module$contents$a_x, module$contents$a_y, module$contents$a_z] = foo();
`)

	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)

	// Object destructuring with explicit names
	c.Test(
		`goog.module('a');
var {p1: x, p2: y} = foo();
`,
		`/** @const */ var module$exports$a = {};
var {p1: module$contents$a_x, p2: module$contents$a_y} = foo();
`)

	// Object destructuring with short names
	c.Test(
		`goog.module('a');
var {x, y} = foo();
`,
		`/** @const */ var module$exports$a = {};
var {x: module$contents$a_x, y: module$contents$a_y} = foo();
`)
}

// Java: ClosureRewriteModuleTest.testShortObjectLiteralsInsideModule()
func (c *ClosureRewriteModuleTest) TestShortObjectLiteralsInsideModule() {
	c.Test(
		`goog.module('a');
var x = foo();
var o = {x};
`,
		`/** @const */ var module$exports$a = {};
var module$contents$a_x = foo();
var module$contents$a_o = {x: module$contents$a_x};
`)

	c.Test(
		`goog.module('a');
var x = foo();
exports = {x};
`,
		`/** @const */ var module$exports$a = {};
module$exports$a.x = foo();
`)
}

// Java: ClosureRewriteModuleTest.testNamedExportInliningJSDoc_atConstJSDoc()
func (c *ClosureRewriteModuleTest) TestNamedExportInliningJSDoc_atConstJSDoc() {
	c.Test(
		c.Srcs(
			`goog.module('ns.b');
/** @const */ var Foo = {a: 0};
exports.Foo = Foo;
`,
			`goog.module('ns.a');
var {Foo} = goog.require('ns.b');
var f = Foo;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
/** @const */ module$exports$ns$b.Foo = {a: 0};
`,
			`/** @const */ var module$exports$ns$a = {}
var module$contents$ns$a_f = module$exports$ns$b.Foo;
`))
}

// Java: ClosureRewriteModuleTest.testNamedExportInliningJSDoc_semanticConst()
func (c *ClosureRewriteModuleTest) TestNamedExportInliningJSDoc_semanticConst() {
	c.Test(
		c.Srcs(
			`goog.module('ns.b');
const Foo = {a: 0};
exports.Foo = Foo;
`,
			`goog.module('ns.a');
var {Foo} = goog.require('ns.b');
var f = Foo;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
/** @const */ module$exports$ns$b.Foo = {a: 0};
`,
			`/** @const */ var module$exports$ns$a = {}
var module$contents$ns$a_f = module$exports$ns$b.Foo;
`))
}

// Java: ClosureRewriteModuleTest.testNamedExportInliningJSDoc_semanticConstWithExistingJSDoc()
func (c *ClosureRewriteModuleTest) TestNamedExportInliningJSDoc_semanticConstWithExistingJSDoc() {
	c.Test(
		c.Srcs(
			`goog.module('ns.b');
/** @interface */ const Foo = class {};
exports.Foo = Foo;
`,
			`goog.module('ns.a');
var {Foo} = goog.require('ns.b');
var f = Foo;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
/** @const @interface */ module$exports$ns$b.Foo = class {};
`,
			`/** @const */ var module$exports$ns$a = {}
var module$contents$ns$a_f = module$exports$ns$b.Foo;
`))
}

// Java: ClosureRewriteModuleTest.testNamedExportInliningJSDoc_semanticConstWithInlineJSDoc()
func (c *ClosureRewriteModuleTest) TestNamedExportInliningJSDoc_semanticConstWithInlineJSDoc() {
	c.Test(
		c.Srcs(
			`goog.module('ns.b');
const /** number */ foo = 0;
exports.foo = foo;
`,
			`goog.module('ns.a');
var {foo} = goog.require('ns.b');
var f = foo;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
/** @const {number} */ module$exports$ns$b.foo = 0;
`,
			`/** @const */ var module$exports$ns$a = {}
var module$contents$ns$a_f = module$exports$ns$b.foo;
`))
}

// Java: ClosureRewriteModuleTest.testNamedExportInliningJSDoc_nonConstWithInlineJSDoc()
func (c *ClosureRewriteModuleTest) TestNamedExportInliningJSDoc_nonConstWithInlineJSDoc() {
	c.Test(
		c.Srcs(
			`goog.module('ns.b');
let /** number */ foo = 0;
exports.foo = foo;
`,
			`goog.module('ns.a');
var {foo} = goog.require('ns.b');
var f = foo;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
/** @type {number} */ module$exports$ns$b.foo = 0;
`,
			`/** @const */ var module$exports$ns$a = {}
var module$contents$ns$a_f = module$exports$ns$b.foo;
`))
}

// Java: ClosureRewriteModuleTest.testNamedExportInliningJSDoc_mutatedLocal()
func (c *ClosureRewriteModuleTest) TestNamedExportInliningJSDoc_mutatedLocal() {
	// TODO(lharker): Stop inlining 'exports.Foo -> Foo' in this case.
	c.Test(
		c.Srcs(
			`goog.module('ns.b');
let Foo = {a: 0};
exports.Foo = Foo;
Foo = {};
`,
			`goog.module('ns.a');
var {Foo} = goog.require('ns.b');
var f = Foo;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
module$exports$ns$b.Foo = {a: 0};
module$exports$ns$b.Foo = {};
`,
			`/** @const */ var module$exports$ns$a = {}
var module$contents$ns$a_f = module$exports$ns$b.Foo;
`))
}

// Java: ClosureRewriteModuleTest.testDeepNamespace()
func (c *ClosureRewriteModuleTest) TestDeepNamespace() {
	c.Test(
		c.Srcs(
			"goog.provide('ns.a.b.Foo'); /** @constructor */ ns.a.b.Foo = function() {};",
			`goog.module('ns.other');

var Foo = goog.require('ns.a.b.Foo');

/** @type {Foo} */
var f = new Foo;
`),
		c.Expected(
			"goog.provide('ns.a.b.Foo'); /** @constructor */ ns.a.b.Foo = function() {};",
			`/** @const */ var module$exports$ns$other = {}
goog.require('ns.a.b.Foo');
/** @type {ns.a.b.Foo} */
var module$contents$ns$other_f = new ns.a.b.Foo;
`))
}

// Java: ClosureRewriteModuleTest.testDestructuringImports()
func (c *ClosureRewriteModuleTest) TestDestructuringImports() {
	c.Test(
		c.Srcs(
			"goog.module('ns.b'); /** @constructor */ exports.Foo = function() {};",
			`goog.module('ns.a');

var {Foo} = goog.require('ns.b');

/** @type {Foo} */
var f = new Foo;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
/** @constructor @const */ module$exports$ns$b.Foo = function() {};
`,
			`/** @const */ var module$exports$ns$a = {}
/** @type {module$exports$ns$b.Foo} */
var module$contents$ns$a_f = new module$exports$ns$b.Foo;
`))

	c.Test(
		c.Srcs(
			"goog.module('ns.b'); /** @typedef {number} */ exports.Foo;",
			`goog.module('ns.a');

var {Foo} = goog.require('ns.b');

/** @type {Foo} */
var f = 4;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
/** @const @typedef {number} */ module$exports$ns$b.Foo;
`,
			`/** @const */ var module$exports$ns$a = {}
/** @type {module$exports$ns$b.Foo} */
var module$contents$ns$a_f = 4;
`))

	c.Test(
		c.Srcs(
			"goog.provide('ns.b'); /** @constructor */ ns.b.Foo = function() {};",
			`goog.module('ns.a');

var {Foo} = goog.require('ns.b');

/** @type {Foo} */
var f = new Foo;
`),
		c.Expected(
			"goog.provide('ns.b'); /** @constructor */ ns.b.Foo = function() {};",
			`/** @const */ var module$exports$ns$a = {}
/** @type {ns.b.Foo} */
var module$contents$ns$a_f = new ns.b.Foo;
`))

	c.Test(
		c.Srcs(
			`goog.module('ns.b');
goog.module.declareLegacyNamespace();

/** @constructor */ exports.Foo = function() {};
`,
			`goog.module('ns.a');

var {Foo} = goog.require('ns.b');

/** @type {Foo} */
var f = new Foo;
`),
		c.Expected(
			`goog.provide('ns.b');
/** @constructor @const */ ns.b.Foo = function() {};
`,
			`/** @const */ var module$exports$ns$a = {}
/** @type {ns.b.Foo} */
var module$contents$ns$a_f = new ns.b.Foo;
`))

	c.Test(
		c.Srcs(
			"goog.module('ns.b'); /** @constructor */ exports.Foo = function() {};",
			`goog.module('ns.a');

var {Foo: Bar} = goog.require('ns.b');

/** @type {Bar} */
var f = new Bar;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
/** @constructor @const */ module$exports$ns$b.Foo = function() {};
`,
			`/** @const */ var module$exports$ns$a = {}
/** @type {module$exports$ns$b.Foo} */
var module$contents$ns$a_f = new module$exports$ns$b.Foo;
`))

	c.Test(
		c.Srcs(
			"goog.module('modA'); class Foo {} exports.Foo = Foo;",
			`goog.module('modB');

var {Foo:importedFoo} = goog.require('modA');

/** @type {importedFoo} */
var f = new importedFoo;
`),
		c.Expected(
			`/** @const */ var module$exports$modA = {};
module$exports$modA.Foo = class {}
`,
			`/** @const */ var module$exports$modB = {}
/** @type {module$exports$modA.Foo} */
var module$contents$modB_f = new module$exports$modA.Foo;
`))

	c.Test(
		c.Srcs(
			"goog.module('modA'); class Foo {} exports.Foo = Foo;",
			`goog.module('modB');

var {Foo} = goog.require('modA');

/** @type {Foo} */
var f = new Foo;
`),
		c.Expected(
			`/** @const */ var module$exports$modA = {};
module$exports$modA.Foo = class {}
`,
			`/** @const */ var module$exports$modB = {}
/** @type {module$exports$modA.Foo} */
var module$contents$modB_f = new module$exports$modA.Foo;
`))

	c.Test(
		c.Srcs(
			"goog.module('modA'); class Foo {} exports = {Foo};",
			`goog.module('modB');

var {Foo} = goog.require('modA');

/** @type {Foo} */
var f = new Foo;
`),
		c.Expected(
			`/** @const */ var module$exports$modA = {};
module$exports$modA.Foo = class {};
`,
			`/** @const */ var module$exports$modB = {}
/** @type {module$exports$modA.Foo} */
var module$contents$modB_f = new module$exports$modA.Foo;
`))

	c.Test(
		c.Srcs(
			"goog.module('modA'); class Bar {} exports = class Foo {}; exports.Bar = Bar;",
			`goog.module('modB');

var Foo = goog.require('modA');

/** @type {Foo} */
var f = new Foo;
`),
		c.Expected(
			`class module$contents$modA_Bar {}
/** @const */ var module$exports$modA = class Foo {};
/** @const */ module$exports$modA.Bar = module$contents$modA_Bar;
`,
			`/** @const */ var module$exports$modB = {}
/** @type {module$exports$modA} */
var module$contents$modB_f = new module$exports$modA;
`))

	c.Test(
		c.Srcs(
			`goog.module('modA');
goog.module.declareLegacyNamespace();

class Foo {}
exports = {Foo};
`,
			`goog.module('modB');

var {Foo} = goog.require('modA');

/** @type {Foo} */
var f = new Foo;
`),
		c.Expected(
			`goog.provide('modA');
class module$contents$modA_Foo {}
/** @const */ modA.Foo = module$contents$modA_Foo;
`,
			`/** @const */ var module$exports$modB = {}
/** @type {modA.Foo} */
var module$contents$modB_f = new modA.Foo;
`))
}

// Java: ClosureRewriteModuleTest.testUninlinableExports()
func (c *ClosureRewriteModuleTest) TestUninlinableExports() {
	c.Test(
		c.Srcs(
			"goog.module('ns.b'); /** @constructor */ exports.Foo = function() {};",
			`goog.module('ns.a');
var {Foo} = goog.require('ns.b');
/** @type {Foo} */ var f = new Foo;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
/** @const @constructor */ module$exports$ns$b.Foo = function() {};
`,
			`/** @const */ var module$exports$ns$a = {}
/** @type {module$exports$ns$b.Foo} */
var module$contents$ns$a_f = new module$exports$ns$b.Foo;
`))

	c.Test(
		c.Srcs(
			"goog.module('ns.b'); /** @constructor */ exports.Foo = function() {};",
			`goog.module('ns.a');
var {Foo} = goog.requireType('ns.b');
/** @type {Foo} */ var f;
`),
		c.Expected(
			`/** @const */ var module$exports$ns$b = {};
/** @const @constructor */ module$exports$ns$b.Foo = function() {};
`,
			`/** @const */ var module$exports$ns$a = {}
/** @type {module$exports$ns$b.Foo} */
var module$contents$ns$a_f;
`))
}

// Java: ClosureRewriteModuleTest.testObjectLiteralDefaultExport()
func (c *ClosureRewriteModuleTest) TestObjectLiteralDefaultExport() {
	c.TestError(
		c.Srcs(
			`goog.module('modA');

class Foo {}
// This is not a named exports object because of the value literal
exports = {Foo, Bar: [1,2,3]};
`,
			`goog.module('modB');

var {Foo} = goog.require('modA');
`),
		ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS)

	c.TestError(
		c.Srcs(
			`goog.module('modA');

class Foo {}
// This is not a named exports object because of the value literal
exports = {Foo, Bar: [1,2,3]};
`,
			`goog.module('modB');

var {Foo} = goog.requireType('modA');
`),
		ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS)
}

// Java: ClosureRewriteModuleTest.testDontUseTheMangledModuleNameInCode()
func (c *ClosureRewriteModuleTest) TestDontUseTheMangledModuleNameInCode() {
	//
	c.TestError(
		c.Srcs(
			"/** @const */ var module$exports$modA = class {};", //
			"goog.module('modA');"),
		ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT)
}

// Java: ClosureRewriteModuleTest.testUninlinableNamedExports()
func (c *ClosureRewriteModuleTest) TestUninlinableNamedExports() {
	c.Test(
		c.Srcs(
			"goog.module('modA'); \n exports = class {};",
			`goog.module('modB');

var Foo = goog.require('modA');

exports.Foo = Foo;
`),
		c.Expected(
			"/** @const */ var module$exports$modA = class {};",
			`/** @const */ var module$exports$modB = {};
/** @const */ module$exports$modB.Foo = module$exports$modA;
`))

	c.Test(
		c.Srcs(
			"goog.module('modA'); \n exports = class {};",
			`goog.module('modB');

var Foo = goog.require('modA');

exports = {Foo};
`),
		c.Expected(
			"/** @const */ var module$exports$modA = class {};",
			`/** @const */ var module$exports$modB = {};
/** @const */ module$exports$modB.Foo = module$exports$modA;
`))

	c.Test(
		c.Srcs(
			"goog.module('modA'); \n exports = class {};",
			`goog.module('modB');

var Foo = goog.require('modA');

exports = {ExportedFoo: Foo};
`),
		c.Expected(
			"/** @const */ var module$exports$modA = class {};",
			`/** @const */ var module$exports$modB = {};
/** @const */ module$exports$modB.ExportedFoo = module$exports$modA;
`))

	c.Test(
		c.Srcs(
			"goog.module('modA'); \n exports = class {};",
			`goog.module('modB');

var Foo = goog.require('modA');
class Bar {}

exports = {Foo, Bar};
`),
		c.Expected(
			"/** @const */ var module$exports$modA = class {};",
			`/** @const */ var module$exports$modB = {};
module$exports$modB.Bar = class {};
/** @const */ module$exports$modB.Foo = module$exports$modA;
`))
}

// Java: ClosureRewriteModuleTest.testIllegalDestructuringImports()
func (c *ClosureRewriteModuleTest) TestIllegalDestructuringImports() {
	c.TestError(
		c.Srcs(
			`goog.module('p.A');
/** @constructor */ var A = function() {}
A.method = function() {}
exports = A
`,
			`goog.module('p.C');
var {method} = goog.require('p.A');
`),
		ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS)

	c.TestError(
		c.Srcs(
			`goog.module('p.A');
/** @constructor */ var A = function() {}
A.method = function() {}
exports = A
`,
			`goog.module('p.C');
var {method} = goog.requireType('p.A');
`),
		ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS)

	c.TestError(
		c.Srcs(
			`goog.module('p.A');
/** @constructor */ exports = class { static method() {} }
`,
			`goog.module('p.C');
var {method} = goog.require('p.A');
`),
		ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS)

	c.TestError(
		c.Srcs(
			`goog.module('p.A');
/** @constructor */ exports = class { static method() {} }
`,
			`goog.module('p.C');
var {method} = goog.requireType('p.A');
`),
		ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS)

	c.TestError(
		c.Srcs(
			`goog.module('p.A');

/** @constructor */ exports.Foo = class {};
/** @constructor */ exports.Bar = class {};
`,
			`goog.module('p.C');

var {Baz} = goog.require('p.A');
`),
		ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS)

	c.TestError(
		c.Srcs(
			`goog.module('p.A');

/** @constructor */ exports.Foo = class {};
/** @constructor */ exports.Bar = class {};
`,
			`goog.module('p.C');

var {Baz} = goog.requireType('p.A');
`),
		ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS)

	// TODO(blickly): We should warn for the next two as well, but it's harder to detect.

	c.Test(
		c.Srcs(
			`goog.provide('p.A');
/** @constructor */ p.A = function() {}
p.A.method = function() {}
`,
			`goog.module('p.C');
var {method} = goog.require('p.A');
`),
		c.Expected([]string(nil)))

	c.Test(
		c.Srcs(
			`goog.provide('p.A');
/** @constructor */ p.A = function() {}
p.A.method = function() {}
`,
			`goog.module('p.C');
var {method} = goog.requireType('p.A');
`),
		c.Expected([]string(nil)))
}

// Java: ClosureRewriteModuleTest.testDeclareLegacyNamespace()
func (c *ClosureRewriteModuleTest) TestDeclareLegacyNamespace() {
	c.Test("goog.module('ns.a'); goog.module.declareLegacyNamespace();", "goog.provide('ns.a');")
}

// Java: ClosureRewriteModuleTest.testRemovesPreventModuleExportSealing()
func (c *ClosureRewriteModuleTest) TestRemovesPreventModuleExportSealing() {
	c.Test(
		"goog.module('ns.a'); goog.module.preventModuleExportSealing();",
		"/** @const */ var module$exports$ns$a = {};")
}

// Java: ClosureRewriteModuleTest.testSideEffectOnlyModuleImport()
func (c *ClosureRewriteModuleTest) TestSideEffectOnlyModuleImport() {
	c.Test(
		c.Srcs(
			`goog.module('ns.b');
alert('hello world');
`,
			`goog.module('ns.c');
alert('hello world');
`,
			`goog.module('ns.a');
goog.require('ns.b');
goog.requireType('ns.c');
`),
		c.Expected(
			"/** @const */ var module$exports$ns$b = {}; alert('hello world');",
			"/** @const */ var module$exports$ns$c = {}; alert('hello world');",
			"/** @const */ var module$exports$ns$a = {};"))
}

// Java: ClosureRewriteModuleTest.testTypeOnlyModuleImport()
func (c *ClosureRewriteModuleTest) TestTypeOnlyModuleImport() {
	c.Test(
		c.Srcs(
			`goog.module('ns.B');
/** @constructor */ exports = function() {};
`,
			`goog.module('ns.C');
/** @constructor */ exports = function() {};
`,
			`goog.module('ns.a');
goog.require('ns.B');
goog.requireType('ns.C');
/** @type {ns.B} */ var b;
/** @type {ns.C} */ var c;
`),
		c.Expected(
			"/** @constructor @const */ var module$exports$ns$B = function() {};",
			"/** @constructor @const */ var module$exports$ns$C = function() {};",
			`/** @const */ var module$exports$ns$a = {};
/** @type {module$exports$ns$B} */ var module$contents$ns$a_b;
/** @type {module$exports$ns$C} */ var module$contents$ns$a_c;
`))
}

// Java: ClosureRewriteModuleTest.testSideEffectOnlyImportOfGoogProvide()
func (c *ClosureRewriteModuleTest) TestSideEffectOnlyImportOfGoogProvide() {
	c.Test(
		c.Srcs(
			`goog.provide('ns.b');
alert('hello world');
`,
			`goog.provide('ns.c');
alert('hello world');
`,
			`goog.module('ns.a');
goog.require('ns.b');
goog.requireType('ns.c');
`),
		c.Expected(
			"goog.provide('ns.b'); alert('hello world');",
			"goog.provide('ns.c'); alert('hello world');",
			`/** @const */ var module$exports$ns$a = {};
goog.require('ns.b');
goog.requireType('ns.c');
`))
}

// Java: ClosureRewriteModuleTest.testSideEffectOnlyImportOfLegacyGoogModule()
func (c *ClosureRewriteModuleTest) TestSideEffectOnlyImportOfLegacyGoogModule() {
	c.Test(
		c.Srcs(
			`goog.module('ns.b');
goog.module.declareLegacyNamespace();
alert('hello world');
`,
			`goog.module('ns.c');
goog.module.declareLegacyNamespace();
alert('hello world');
`,
			`goog.module('ns.a');
goog.require('ns.b');
goog.requireType('ns.c');
`),
		c.Expected(
			"goog.provide('ns.b'); alert('hello world');",
			"goog.provide('ns.c'); alert('hello world');",
			`/** @const */ var module$exports$ns$a = {};
goog.require('ns.b');
goog.requireType('ns.c');
`))
}

// Java: ClosureRewriteModuleTest.testTypeOnlyModuleImportFromLegacyFile()
func (c *ClosureRewriteModuleTest) TestTypeOnlyModuleImportFromLegacyFile() {
	c.Test(
		c.Srcs(
			`goog.module('ns.B');
/** @constructor */ exports = function() {};
`,
			`goog.module('ns.C');
/** @constructor */ exports = function() {};
`,
			`goog.provide('ns.a');
goog.require('ns.B');
goog.requireType('ns.C');
/** @type {ns.B} */ var b;
/** @type {ns.C} */ var c;
`),
		c.Expected(
			"/** @constructor @const */ var module$exports$ns$B = function() {};",
			"/** @constructor @const */ var module$exports$ns$C = function() {};",
			`goog.provide('ns.a');
/** @type {module$exports$ns$B} */ var b;
/** @type {module$exports$ns$C} */ var c;
`))
}

// Java: ClosureRewriteModuleTest.testTypeOnlyModuleImportFromLegacyFile_ofLegacyModule()
func (c *ClosureRewriteModuleTest) TestTypeOnlyModuleImportFromLegacyFile_ofLegacyModule() {
	c.Test(
		c.Srcs(
			`goog.module('ns.C');
goog.module.declareLegacyNamespace();
/** @constructor */ exports = function() {};
`,
			`goog.provide('ns.a');
goog.requireType('ns.C');
/** @type {ns.C} */ var c;
`),
		c.Expected(
			`goog.provide('ns.C');
/** @constructor @const */ ns.C = function() {};
`,
			`goog.provide('ns.a');
goog.requireType('ns.C');
/** @type {ns.C} */ var c;
`))
}

// Java: ClosureRewriteModuleTest.testBundle1()
func (c *ClosureRewriteModuleTest) TestBundle1() {
	c.Test(
		c.Srcs(
			"goog.module('ns.b');",
			`goog.loadModule(function(exports) {
  goog.module('ns.a');
  var b = goog.require('ns.b');
  exports.b = b;
  return exports;
});
`),
		c.Expected(
			"/** @const */ var module$exports$ns$b = {};",
			`/** @const */ var module$exports$ns$a = {};
/** @const */ module$exports$ns$a.b = module$exports$ns$b;
`))
}

// Java: ClosureRewriteModuleTest.testBundle2()
func (c *ClosureRewriteModuleTest) TestBundle2() {
	c.Test(
		`goog.loadModule(function(exports) {
  goog.module('ns.b');
  return exports;
});
goog.loadModule(function(exports) {
  goog.module('ns.a');
  var b = goog.require('ns.b');
  exports.b = b;
  return exports;
});
goog.loadModule(function(exports) {
  goog.module('ns.c');
  var b = goog.require('ns.b');
  exports.b = b;
  return exports;
});
`,
		`/** @const */ var module$exports$ns$b = {};
/** @const */ var module$exports$ns$a = {};
/** @const */ module$exports$ns$a.b = module$exports$ns$b;
/** @const */ var module$exports$ns$c = {};
/** @const */ module$exports$ns$c.b = module$exports$ns$b;
`)
}

// Java: ClosureRewriteModuleTest.testBundle3()
func (c *ClosureRewriteModuleTest) TestBundle3() {
	c.Test(
		`goog.loadModule(function(exports) {
  goog.module('ns.b');
  return exports;
});
goog.loadModule(function(exports) {
  'use strict';
  goog.module('ns.a');
  goog.module.declareLegacyNamespace();
  var b = goog.require('ns.b');
  return exports;
});
`,
		`/** @const */ var module$exports$ns$b = {};
goog.provide('ns.a');
`)
}

// Java: ClosureRewriteModuleTest.testBundle4()
func (c *ClosureRewriteModuleTest) TestBundle4() {
	c.Test(
		`goog.loadModule(function(exports) {
  goog.module('goog.asserts');
  return exports;
});
goog.loadModule(function(exports) {
  'use strict';
  goog.module('ns.a');
  var b = goog.require('goog.asserts');
  return exports;
});
`,
		`/** @const */ var module$exports$goog$asserts = {};
/** @const */ var module$exports$ns$a = {};
`)
}

// Java: ClosureRewriteModuleTest.testBundle5()
func (c *ClosureRewriteModuleTest) TestBundle5() {
	c.Test(
		`goog.loadModule(function(exports) {
  goog.module('goog.asserts');
  return exports;
});
goog.loadModule(function(exports) {
  'use strict';
  goog.module('xid');
  goog.module.declareLegacyNamespace();
  var asserts = goog.require('goog.asserts');
  exports = function(id) {
    return xid.internal_(id);
  };
  var xid = exports;
  exports.internal_ = function(id) {};
  return exports;
});
`,
		`/** @const */ var module$exports$goog$asserts = {};
goog.provide('xid');
/** @const */ xid = function(id) {
  return module$contents$xid_xid.internal_(id);
};
var module$contents$xid_xid = xid;
/** @const */
module$contents$xid_xid.internal_ = function(id) {};
`)
}

// Java: ClosureRewriteModuleTest.testBundle6()
func (c *ClosureRewriteModuleTest) TestBundle6() {
	c.Test(
		`goog.loadModule(function(exports) {
  goog.module('goog.asserts');
  return exports;
});
goog.loadModule(function(exports) {
  'use strict';
  goog.module('xid');
  goog.module.declareLegacyNamespace();
  var asserts = goog.require('goog.asserts');
  var xid = function(id) {
    return xid.internal_(id);
  };
  xid.internal_ = function(id) {};
  exports = xid;
  return exports;
});
`,
		`/** @const */ var module$exports$goog$asserts = {};
goog.provide('xid');
var module$contents$xid_xid = function(id) {
  return module$contents$xid_xid.internal_(id);
};
module$contents$xid_xid.internal_ = function(id) {};
/** @const */ xid = module$contents$xid_xid
`)
}

// Java: ClosureRewriteModuleTest.testBundleWithDestructuringImport()
func (c *ClosureRewriteModuleTest) TestBundleWithDestructuringImport() {
	c.Test(
		`goog.loadModule(function(exports) { 'use strict';
  goog.module('mod_B');

  /** @interface */ function B(){}

  exports.B = B;
  return exports;
});
goog.loadModule(function(exports) { 'use strict';
  goog.module('mod_A');

  var {B} = goog.require('mod_B');

  /** @constructor @implements {B} */
  function A() {}
  return exports;
});
`,
		`/** @const */ var module$exports$mod_B = {};
/** @interface */ function module$contents$mod_B_B(){}
/** @const */ module$exports$mod_B.B = module$contents$mod_B_B;

/** @const */ var module$exports$mod_A = {};
/** @constructor @implements {module$exports$mod_B.B} */
function module$contents$mod_A_A(){}
`)
}

// Java: ClosureRewriteModuleTest.testGoogLoadModule_missingReturn()
func (c *ClosureRewriteModuleTest) TestGoogLoadModule_missingReturn() {
	c.TestError(
		`goog.loadModule(function(exports) {
  goog.module('ns.b');
});
`,
		ClosureRewriteModule_LOAD_MODULE_FN_MISSING_RETURN)
}

// Java: ClosureRewriteModuleTest.testGoogLoadModuleString()
func (c *ClosureRewriteModuleTest) TestGoogLoadModuleString() {
	AssertThrows[jutil.IIllegalArgumentException](
		func() { c.TestSame("goog.loadModule(\"goog.module('a.b.c'); exports = class {};\");") })
}

// Java: ClosureRewriteModuleTest.testTopLevelNames1()
func (c *ClosureRewriteModuleTest) TestTopLevelNames1() {
	// Vars defined inside functions are not module top level.
	c.Test(
		`goog.module('a');
var a, b, c;
function Outer() {
  var a, b, c;
  function Inner() {
    var a, b, c;
  }
}
`,
		`/** @const */ var module$exports$a = {};
var module$contents$a_a, module$contents$a_b, module$contents$a_c;
function module$contents$a_Outer() {
  var a, b, c;
  function Inner() {
    var a, b, c
  }
}
`)
}

// Java: ClosureRewriteModuleTest.testTopLevelNames2()
func (c *ClosureRewriteModuleTest) TestTopLevelNames2() {
	// Vars in blocks are module top level because they are hoisted to the first execution context.
	c.Test(
		`goog.module('a.c');
if (true) {
  var a, b, c;
}
`,
		`/** @const */ var module$exports$a$c = {};
if (true) {
  var module$contents$a$c_a, module$contents$a$c_b, module$contents$a$c_c;
}
`)
}

// Java: ClosureRewriteModuleTest.testTopLevelNames3()
func (c *ClosureRewriteModuleTest) TestTopLevelNames3() {
	// Functions in blocks are not module top level because they are block scoped.
	c.Test(
		`goog.module('a.c');
if (true) {
  function a() {}
  function b() {}
  function c() {}
}
`,
		`/** @const */ var module$exports$a$c = {};
if (true) {
  function a() {}
  function b() {}
  function c() {}
}
`)
}

// Java: ClosureRewriteModuleTest.testThis()
func (c *ClosureRewriteModuleTest) TestThis() {
	// global "this" is retained.
	c.Test(
		`goog.module('a');
this;
`,
		`/** @const */ var module$exports$a = {};
this;
`)
}

// Java: ClosureRewriteModuleTest.testGoogImport_await()
func (c *ClosureRewriteModuleTest) TestGoogImport_await() {
	c.Test(
		c.Srcs(
			"goog.module('a.b.c');",
			`async function test() {
 var d = await goog.requireDynamic('a.b.c');
}
`),
		c.Expected(
			"/** @const */ var module$exports$a$b$c = {}",
			`async function test() {
 await goog.importHandler_('sG5M4c');
 var d = module$exports$a$b$c;
}
`))
}

// Java: ClosureRewriteModuleTest.testGoogImport_await_destructuring()
func (c *ClosureRewriteModuleTest) TestGoogImport_await_destructuring() {
	c.Test(
		c.Srcs(
			`goog.module('a.b.c');
exports.Foo=class{}
`,
			`async function test() {
 var {Foo} = await goog.requireDynamic('a.b.c');
}
`),
		c.Expected(
			`/** @const */ var module$exports$a$b$c = {};
/** @const */ module$exports$a$b$c.Foo = class {}
`,
			`async function test() {
 await goog.importHandler_('sG5M4c');
 var {Foo} = module$exports$a$b$c;
}
`))
}

// Java: ClosureRewriteModuleTest.testGoogRequireDynamic_then_destructuringPattern()
func (c *ClosureRewriteModuleTest) TestGoogRequireDynamic_then_destructuringPattern() {
	c.Test(
		c.Srcs(
			`goog.module('a.b.c');
exports.Foo=class{}
`,
			`async function test() {
  goog.requireDynamic('a.b.c').then(({Foo}) => {console.log(Foo);});
}
`),
		c.Expected(
			`/** @const */ var module$exports$a$b$c = {};
/** @const */ module$exports$a$b$c.Foo = class {}
`,
			`async function test() {
  goog.importHandler_('sG5M4c').then(() => {
     const {Foo} = module$exports$a$b$c;
     console.log(Foo);
  });
}
`))
}

// Java: ClosureRewriteModuleTest.testGoogRequireDynamic_then_destructuringPattern_withExpressionBody()
func (c *ClosureRewriteModuleTest) TestGoogRequireDynamic_then_destructuringPattern_withExpressionBody() {
	c.Test(
		c.Srcs(
			`goog.module('a.b.c');
exports.Foo=class{}
`,
			`async function test() {
  goog.requireDynamic('a.b.c').then(({Foo}) => Foo);
}
`),
		c.Expected(
			`/** @const */ var module$exports$a$b$c = {};
/** @const */ module$exports$a$b$c.Foo = class {}
`,
			`async function test() {
  goog.importHandler_('sG5M4c').then(() => {
     const {Foo} = module$exports$a$b$c;
     return Foo;
  });
}
`))
	c.Test(
		c.Srcs(
			`goog.module('a.b.c');
exports.Foo=class{}
`,
			`async function test() {
  goog.requireDynamic('a.b.c').then(({Foo}) => console.log(Foo));
}
`),
		c.Expected(
			`/** @const */ var module$exports$a$b$c = {};
/** @const */ module$exports$a$b$c.Foo = class {}
`,
			`async function test() {
  goog.importHandler_('sG5M4c').then(() => {
     const {Foo} = module$exports$a$b$c;
     return console.log(Foo);
  });
}
`))
}

// Java: ClosureRewriteModuleTest.testGoogRequireDynamic_then_name()
func (c *ClosureRewriteModuleTest) TestGoogRequireDynamic_then_name() {
	c.Test(
		c.Srcs(
			`goog.module('a.b.c');
exports.Foo=class{}
`,
			`async function test() {
  goog.requireDynamic('a.b.c').then((foo) => {console.log(foo.Foo);});
}
`),
		c.Expected(
			`/** @const */ var module$exports$a$b$c = {};
/** @const */ module$exports$a$b$c.Foo = class {}
`,
			`async function test() {
  goog.importHandler_('sG5M4c').then(() => {
     const foo = module$exports$a$b$c;
     console.log(foo.Foo);
  });
}
`))
}

// Java: ClosureRewriteModuleTest.testGoogRequireDynamic_then_name_laterLoadedModule()
func (c *ClosureRewriteModuleTest) TestGoogRequireDynamic_then_name_laterLoadedModule() {
	c.Test(
		c.Srcs(
			`async function test() {
  goog.requireDynamic('a.b.c').then((foo) => {console.log(foo.Foo);});
}
`,
			`goog.module('a.b.c');
exports.Foo=class{}
`),
		c.Expected(
			`async function test() {
  goog.importHandler_('sG5M4c').then(() => {
     const foo = module$exports$a$b$c;
     console.log(foo.Foo);
  });
}
`,
			`/** @const */ var module$exports$a$b$c = {};
/** @const */ module$exports$a$b$c.Foo = class {}
`))
}

// Java: ClosureRewriteModuleTest.testGoogRequireDynamic_then_missingSources()
func (c *ClosureRewriteModuleTest) TestGoogRequireDynamic_then_missingSources() {
	c.allowMissingSources = true
	c.Test(
		c.Srcs(
			`async function test() {
  goog.requireDynamic('a.b.c').then((foo) => {console.log(foo.Foo);});
}
`),
		c.Expected(
			`async function test() {
  null.then((foo) => {console.log(foo.Foo);});
}
`))
}

// Java: ClosureRewriteModuleTest.testRequireDynamic_illegal_await_lhs()
func (c *ClosureRewriteModuleTest) TestRequireDynamic_illegal_await_lhs() {
	c.TestError(
		c.Srcs(
			`goog.module('a.b.c');
exports.Foo=class{}
`,
			`async function test() {
await goog.requireDynamic('a.b.c');
}
`),
		ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT)
}

// Java: ClosureRewriteModuleTest.testGoogModuleGet1()
func (c *ClosureRewriteModuleTest) TestGoogModuleGet1() {
	c.Test(
		c.Srcs(
			"goog.module('a');",
			`function f() {
  var x = goog.module.get('a');
}
`),
		c.Expected(
			"/** @const */ var module$exports$a = {};",
			`function f() {
  var x = module$exports$a;
}
`))
}

// Java: ClosureRewriteModuleTest.testGoogModuleGet2()
func (c *ClosureRewriteModuleTest) TestGoogModuleGet2() {
	c.Test(
		c.Srcs(
			"goog.module('a.b.c');",
			`function f() {
  var x = goog.module.get('a.b.c');
}
`),
		c.Expected(
			"/** @const */ var module$exports$a$b$c = {};",
			`function f() {
  var x = module$exports$a$b$c;
}
`))
}

// Java: ClosureRewriteModuleTest.testGoogModuleGet3()
func (c *ClosureRewriteModuleTest) TestGoogModuleGet3() {
	c.Test(
		c.Srcs(
			`goog.module('a.b.c');
exports = class {};
`,
			`goog.module('x.y.z');

function f() {
  return new (goog.module.get('a.b.c'))();
}
`),
		c.Expected(
			"/** @const */ var module$exports$a$b$c = class {};",
			`/** @const */ var module$exports$x$y$z = {};
function module$contents$x$y$z_f() {
  return new module$exports$a$b$c();
}
`))
}

// Java: ClosureRewriteModuleTest.testGoogModuleGet_missingInExpression()
func (c *ClosureRewriteModuleTest) TestGoogModuleGet_missingInExpression() {
	c.AllowExternsChanges()
	c.Test(
		c.Srcs(
			`goog.module('x.y.z');

function f() {
  return goog.module.get('a.b.c');
}
f();
`),
		c.Expected(
			`/** @const */
var module$exports$x$y$z={};

function module$contents$x$y$z_f() {
  return null;
}

module$contents$x$y$z_f();
`))
}

// Java: ClosureRewriteModuleTest.testGoogModuleGet_missingAsRhs()
func (c *ClosureRewriteModuleTest) TestGoogModuleGet_missingAsRhs() {
	c.AllowExternsChanges()
	c.Test(
		c.Srcs(
			`goog.module('x.y.z');

function f() {
  const c = goog.module.get('a.b.c');
  return c;
}
f();
`),
		c.Expected(
			`/** @const */
var module$exports$x$y$z={};

function module$contents$x$y$z_f() {
  return c;
}

module$contents$x$y$z_f();
`))
}

// Java: ClosureRewriteModuleTest.testGoogRequire_missing_createSyntheticExterns()
func (c *ClosureRewriteModuleTest) TestGoogRequire_missing_createSyntheticExterns() {
	// typechecking reports errors due to us not including the Closure externs
	c.IgnoreWarnings(DiagnosticGroups_CHECK_TYPES)

	c.TestExternChanges(
		c.Srcs(
			`goog.module('mod');
const c = goog.require('a.b.c');
c;
`),
		c.Expected("var c;"))

	//
	c.TestExternChanges(
		c.Srcs(
			`goog.module('mod');
const {D: LocalD, E} = goog.require('a.b.c');
LocalD;
E
`),
		c.Expected("var LocalD;var E;"))
}

// Java: ClosureRewriteModuleTest.testGoogModuleGet_missing_createSyntheticExterns()
func (c *ClosureRewriteModuleTest) TestGoogModuleGet_missing_createSyntheticExterns() {
	// typechecking reports errors due to us not including the Closure externs
	c.IgnoreWarnings(DiagnosticGroups_CHECK_TYPES)

	c.TestExternChanges(
		c.Srcs(
			`(function() {
  const c = goog.module.get('a.b.c');
  const f = goog.module.get('d.e.f');
})
`),
		c.Expected("var c;var f;"))

	// Declare the 'lhs' even when goog.module.get is nested within the right-hand side.
	c.TestExternChanges(
		c.Srcs(
			`(function() {
  const result = process(goog.module.get('a.b.c').d) * 2;
})
`),
		c.Expected("var result;"))

	// Don't declare any names multiple times.
	c.TestExternChanges(
		c.Srcs(
			`(function() {
  const c = goog.module.get('a.b.c');
  c;
  {
    const c = goog.module.get('other.a.b.c');
  }
})
`),
		c.Expected("var c;"))
}

// Java: ClosureRewriteModuleTest.testGoogModuleGet_missing_noSyntheticExternForExistingGlobals()
func (c *ClosureRewriteModuleTest) TestGoogModuleGet_missing_noSyntheticExternForExistingGlobals() {
	c.TestNoWarning(
		c.Srcs(
			`var c;
(function() {
  const c = goog.module.get('a.b.c');
})
`))

	c.TestNoWarning(
		c.Srcs(
			`goog.provide('c');
(function() {
  const c = goog.module.get('a.b.c');
})
`))

	c.TestNoWarning(
		c.Srcs(
			`goog.provide('c.d');
(function() {
  const c = goog.module.get('a.b.c');
})
`))

	c.TestNoWarning(
		c.Srcs(
			`goog.module('c.d');
goog.module.declareLegacyNamespace();
`,
			`(function() {
  const c = goog.module.get('a.b.c');
})
`))
}

// Java: ClosureRewriteModuleTest.testGoogModuleGet_missing_preserve()
func (c *ClosureRewriteModuleTest) TestGoogModuleGet_missing_preserve() {
	c.preserveClosurePrimitives = true
	// Need to disable tree comparison because compiler adds MODULE_BODY token when parsing
	// expected output but it is not present in actual tree.
	c.DisableCompareAsTree()

	c.Test(
		c.Srcs(
			`goog.module('x.y.z');

function f() {
  return goog.module.get('a.b.c');
}
f();
`),
		c.Expected(
			`goog.module("x.y.z");
/** @const */
var module$exports$x$y$z = {};
function module$contents$x$y$z_f() {
  return goog.module.get("a.b.c");
}
module$contents$x$y$z_f();
`))
}

// Java: ClosureRewriteModuleTest.testAliasedGoogModuleGet1()
func (c *ClosureRewriteModuleTest) TestAliasedGoogModuleGet1() {
	c.Test(
		c.Srcs(
			"goog.module('b'); exports = class {};",
			`goog.module('a');
var x = goog.forwardDeclare('b');
function f() {
  x = goog.module.get('b');
  new x;
}
`),
		c.Expected(
			"/** @const */ var module$exports$b = class {};",
			`/** @const */ var module$exports$a = {};
function module$contents$a_f() {
  new module$exports$b;
}
`))
}

// Java: ClosureRewriteModuleTest.testAliasedGoogModuleGet2()
func (c *ClosureRewriteModuleTest) TestAliasedGoogModuleGet2() {
	c.Test(
		c.Srcs(
			"goog.module('x.y.z'); exports = class {};",
			`goog.module('a');
var x = goog.forwardDeclare('x.y.z');
function f() {
  x = goog.module.get('x.y.z');
  new x;
}
`),
		c.Expected(
			"/** @const */ var module$exports$x$y$z = class {};",
			`/** @const */ var module$exports$a = {};
function module$contents$a_f() {
  new module$exports$x$y$z;
}
`))
}

// Java: ClosureRewriteModuleTest.testAliasedGoogModuleGet3()
func (c *ClosureRewriteModuleTest) TestAliasedGoogModuleGet3() {
	c.Test(
		c.Srcs(
			`goog.module('a.b.c');
/** @constructor */ function C() {}
exports = C
`,
			`/** @type {a.b.c} */ var c;
function f() {
  var C = goog.module.get('a.b.c');
  c = new C;
}
`),
		c.Expected(
			`/** @constructor */ function module$contents$a$b$c_C() {}
/** @const */ var module$exports$a$b$c = module$contents$a$b$c_C;
`,
			`/** @type {module$exports$a$b$c} */ var c;
function f() {
  var C = module$exports$a$b$c;
  c = new C;
}
`))
}

// Java: ClosureRewriteModuleTest.testAliasedGoogModuleGet4()
func (c *ClosureRewriteModuleTest) TestAliasedGoogModuleGet4() {
	c.Test(
		c.Srcs(
			`goog.module('x.y.z');
/** @constructor */ function Z() {}
exports = Z
`,
			`goog.module('a');
/** @type {z} */ var c;
var z = goog.forwardDeclare('x.y.z');
function f() {
  z = goog.module.get('x.y.z');
  c = new z;
}
`),
		c.Expected(
			`/** @constructor */ function module$contents$x$y$z_Z() {}
/** @const */ var module$exports$x$y$z = module$contents$x$y$z_Z;
`,
			`/** @const */ var module$exports$a = {};
/** @type {module$contents$a_z} */ var module$contents$a_c;
function module$contents$a_f() {
  module$contents$a_c = new module$exports$x$y$z;
}
`))
}

// Java: ClosureRewriteModuleTest.testAliasedGoogModuleGet5()
func (c *ClosureRewriteModuleTest) TestAliasedGoogModuleGet5() {
	c.Test(
		c.Srcs(
			"goog.provide('b'); b = class {};",
			`goog.module('a');
var x = goog.forwardDeclare('b');
function f() {
  x = goog.module.get('b');
  new x;
}
`),
		c.Expected(
			"goog.provide('b'); b = class {};",
			`/** @const */ var module$exports$a = {};
goog.forwardDeclare('b');
function module$contents$a_f() {
  new b;
}
`))
}

// Java: ClosureRewriteModuleTest.testAliasedGoogModuleGet6()
func (c *ClosureRewriteModuleTest) TestAliasedGoogModuleGet6() {
	c.Test(
		c.Srcs(
			"goog.provide('x.y.z'); x.y.z = class {};",
			`goog.module('a');
var z = goog.forwardDeclare('x.y.z');
function f() {
  z = goog.module.get('x.y.z');
  new z;
}
`),
		c.Expected(
			"goog.provide('x.y.z'); x.y.z = class {};",
			`/** @const */ var module$exports$a = {};
goog.forwardDeclare('x.y.z');
function module$contents$a_f() {
  new x.y.z;
}
`))
}

// Java: ClosureRewriteModuleTest.testAliasedGoogModuleGet7()
func (c *ClosureRewriteModuleTest) TestAliasedGoogModuleGet7() {
	c.Test(
		c.Srcs(
			"goog.module('a.b.c.D');",
			`goog.require('a.b.c.D');
goog.scope(function() {
var D = goog.module.get('a.b.c.D');
});
`),
		c.Expected(
			"/** @const */ var module$exports$a$b$c$D = {};",
			`goog.scope(function() {
var D = module$exports$a$b$c$D;
});
`))
}

// Java: ClosureRewriteModuleTest.testAliasedGoogModuleGet8()
func (c *ClosureRewriteModuleTest) TestAliasedGoogModuleGet8() {
	c.Test(
		c.Srcs(
			"goog.module('a.b.c.D'); exports = class {};",
			`goog.require('a.b.c.D');
goog.scope(function() {
var D = goog.module.get('a.b.c.D');
var d = new D;
});
`),
		c.Expected(
			"/** @const */ var module$exports$a$b$c$D = class {};",
			`goog.scope(function() {
var D = module$exports$a$b$c$D;
var d = new D;
});
`))
}

// Java: ClosureRewriteModuleTest.testAliasedGoogModuleGetInGoogProvide()
func (c *ClosureRewriteModuleTest) TestAliasedGoogModuleGetInGoogProvide() {
	c.Test(
		c.Srcs(
			"goog.module('a.b.c.D');",
			`goog.provide('x.y.z');
goog.require('a.b.c.D');
x.y.z = goog.module.get('a.b.c.D');
`),
		c.Expected(
			"/** @const */ var module$exports$a$b$c$D = {};",
			`goog.provide('x.y.z');
x.y.z = module$exports$a$b$c$D;
`))
}

// Java: ClosureRewriteModuleTest.testWeakGoogModuleGet_error()
func (c *ClosureRewriteModuleTest) TestWeakGoogModuleGet_error() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCodeWithKind("b.js", "goog.module('b'); exports = class {};", rhino.StaticSourceFile_SourceKind_WEAK),
			SourceFile_FromCode(
				"a.js",
				`goog.module('a');
function f() {
  const b = goog.module.get('b');
}
`)),
		ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE)
}

// Java: ClosureRewriteModuleTest.testWeakGoogModuleGet_suppressWeakModuleGet()
func (c *ClosureRewriteModuleTest) TestWeakGoogModuleGet_suppressWeakModuleGet() {
	c.Test(
		c.Srcs(
			SourceFile_FromCodeWithKind("b.js", "goog.module('b'); exports = class {};", rhino.StaticSourceFile_SourceKind_WEAK),
			SourceFile_FromCode(
				"a.js",
				`goog.module('a');
function f() {
  /** @suppress {weakModuleGet} */
  const b = goog.module.get('b');
}
`)),
		c.Expected(
			`/** @const */ var module$exports$a = {};
function module$contents$a_f() {
  /** @suppress {weakModuleGet} */
  const b = null;
}
`,
			"/** @const */ var module$exports$b = class {};"))
}

// Java: ClosureRewriteModuleTest.testWeakGoogModuleGet_reverseFileOrder()
func (c *ClosureRewriteModuleTest) TestWeakGoogModuleGet_reverseFileOrder() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				"a.js",
				`goog.module('a');
function f() {
  const b = goog.module.get('b');
}
`),
			SourceFile_FromCodeWithKind("b.js", "goog.module('b'); exports = class {};", rhino.StaticSourceFile_SourceKind_WEAK)),
		ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE)
}

// Java: ClosureRewriteModuleTest.testWeakGoogProvideGet_error()
func (c *ClosureRewriteModuleTest) TestWeakGoogProvideGet_error() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCodeWithKind("b.js", "goog.provide('b'); b = class {};", rhino.StaticSourceFile_SourceKind_WEAK),
			SourceFile_FromCode(
				"a.js",
				`goog.module('a');
function f() {
  const b = goog.module.get('b');
}
`)),
		ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE)
}

// Java: ClosureRewriteModuleTest.testWeakAliasedGoogModuleGet_error()
func (c *ClosureRewriteModuleTest) TestWeakAliasedGoogModuleGet_error() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCodeWithKind("b.js", "goog.module('b'); exports = class {};", rhino.StaticSourceFile_SourceKind_WEAK),
			SourceFile_FromCode(
				"a.js",
				`goog.module('a');
var x = goog.forwardDeclare('b');
function f() {
  x = goog.module.get('b');
  new x;
}
`)),
		ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE)
}

// Java: ClosureRewriteModuleTest.testWeakAliasedGoogModuleGet_suppress()
func (c *ClosureRewriteModuleTest) TestWeakAliasedGoogModuleGet_suppress() {
	c.Test(
		c.Srcs(
			SourceFile_FromCodeWithKind("b.js", "goog.module('b'); exports = class {};", rhino.StaticSourceFile_SourceKind_WEAK),
			SourceFile_FromCode(
				"a.js",
				`goog.module('a');
var x = goog.forwardDeclare('b');
function f() {
  /** @suppress {weakModuleGet} */
  x = goog.module.get('b');
  new x;
}
`)),
		c.Expected(
			`/** @const */ var module$exports$a = {};
var module$contents$a_x = null;
function module$contents$a_f() {
  /** @suppress {weakModuleGet} */
  module$contents$a_x = null;
  new module$contents$a_x();
}
`,
			"/** @const */ var module$exports$b = class {};"))
}

// Java: ClosureRewriteModuleTest.testWeakGoogModuleGet_inExpression()
func (c *ClosureRewriteModuleTest) TestWeakGoogModuleGet_inExpression() {
	c.Test(
		c.Srcs(
			SourceFile_FromCodeWithKind("b.js", "goog.module('b'); exports = class {};", rhino.StaticSourceFile_SourceKind_WEAK),
			SourceFile_FromCode(
				"a.js",
				`goog.module('a');
/** @suppress {weakModuleGet} */
function f() {
  if (goog.module.get('b')) {
    return goog.module.get('b');
  }
  return null;
}
`)),
		c.Expected(
			`/** @const */ var module$exports$a = {};
/** @suppress {weakModuleGet} */
function module$contents$a_f() {
  if (null) {
    return null;
  }
  return null;
}
`,
			"/** @const */ var module$exports$b = class {};"))
}

// Java: ClosureRewriteModuleTest.testWeakCallingFile_doesNotReportError()
func (c *ClosureRewriteModuleTest) TestWeakCallingFile_doesNotReportError() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode("main.js", "goog.requireType('a');"),
			SourceFile_FromCodeWithKind("b.js", "goog.module('b'); exports = class {};", rhino.StaticSourceFile_SourceKind_WEAK),
			SourceFile_FromCodeWithKind(
				"a.js",
				`goog.provide('a');
goog.require('b');
function f() {
  const b = goog.module.get('b');
}
`, rhino.StaticSourceFile_SourceKind_WEAK)),
		c.Expected(
			"goog.requireType('a');",
			"/** @const */ var module$exports$b = class {};",
			`goog.provide('a');
function f() {
  const b = module$exports$b;
}
`))
}

// Java: ClosureRewriteModuleTest.testWeakGoogModuleGet_withGoogRequire_doesNotReportError()
func (c *ClosureRewriteModuleTest) TestWeakGoogModuleGet_withGoogRequire_doesNotReportError() {
	c.Test(
		c.Srcs(
			SourceFile_FromCodeWithKind("b.js", "goog.module('b'); exports = class {};", rhino.StaticSourceFile_SourceKind_WEAK),
			SourceFile_FromCode(
				"a.js",
				`goog.provide('a');
goog.require('b');
function f() {
  const b = goog.module.get('b');
}
`)),
		c.Expected(
			`goog.provide('a');
function f() {
  const b = module$exports$b;
}
`,
			"/** @const */ var module$exports$b = class {};"))
}

// Java: ClosureRewriteModuleTest.testInvalidGoogForwardDeclareParameter()
func (c *ClosureRewriteModuleTest) TestInvalidGoogForwardDeclareParameter() {
	// Wrong parameter count.
	c.TestError(
		`goog.module('a');
var x = goog.forwardDeclare();
`,
		ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE)

	// Wrong parameter count.
	c.TestError(
		`goog.module('a');
var x = goog.forwardDeclare('a', 'b');
`,
		ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE)

	// Wrong parameter type.
	c.TestError(
		`goog.module('a');
var x = goog.forwardDeclare({});
`,
		ClosurePrimitiveErrors_INVALID_FORWARD_DECLARE_NAMESPACE)
}

// Java: ClosureRewriteModuleTest.testInvalidGoogModuleGet1()
func (c *ClosureRewriteModuleTest) TestInvalidGoogModuleGet1() {
	c.TestError(
		`function f() {
  goog.module.get(a);
}
`,
		ClosurePrimitiveErrors_INVALID_GET_NAMESPACE)
}

// Java: ClosureRewriteModuleTest.testInvalidGoogModuleGet2()
func (c *ClosureRewriteModuleTest) TestInvalidGoogModuleGet2() {
	// This is checked earlier, in CheckClosureImports, not in the rewriting pass.
	c.TestNoWarning(c.Srcs("goog.module('a');", "goog.module.get('a');"))
}

// Java: ClosureRewriteModuleTest.testExtractableExport1()
func (c *ClosureRewriteModuleTest) TestExtractableExport1() {
	c.Test(
		`goog.module('xid');
var xid = function() {};
exports = xid;
`,
		"var module$exports$xid = function() {};")
}

// Java: ClosureRewriteModuleTest.testExtractableExport2()
func (c *ClosureRewriteModuleTest) TestExtractableExport2() {
	c.Test(
		`goog.module('xid');
function xid() {}
exports = xid;
`,
		`function module$contents$xid_xid() {}
/** @const */ var module$exports$xid = module$contents$xid_xid;
`)
}

// Java: ClosureRewriteModuleTest.testExtractableExport3()
func (c *ClosureRewriteModuleTest) TestExtractableExport3() {
	c.Test(
		`goog.module('Foo');
class Foo {}
exports = Foo;
`,
		"class module$exports$Foo {}")
}

// Java: ClosureRewriteModuleTest.testExtractableExport4()
func (c *ClosureRewriteModuleTest) TestExtractableExport4() {
	c.Test(
		`goog.module('Foo');
const Foo = class {}
exports = Foo;
`,
		"const module$exports$Foo = class {};")
}

// Java: ClosureRewriteModuleTest.testExport0()
func (c *ClosureRewriteModuleTest) TestExport0() {
	c.Test("goog.module('ns.a');", "/** @const */ var module$exports$ns$a = {};")
}

// Java: ClosureRewriteModuleTest.testExport1()
func (c *ClosureRewriteModuleTest) TestExport1() {
	c.Test(
		`goog.module('ns.a');
exports = {};
`,
		"/** @const */ var module$exports$ns$a = {};")
}

// Java: ClosureRewriteModuleTest.testExport2()
func (c *ClosureRewriteModuleTest) TestExport2() {
	c.Test(
		`goog.module('ns.a');
exports.x = 1;
`,
		`/** @const */ var module$exports$ns$a = {};
/** @const */ module$exports$ns$a.x = 1
`)
}

// Java: ClosureRewriteModuleTest.testExport4()
func (c *ClosureRewriteModuleTest) TestExport4() {
	c.Test(
		`goog.module('ns.a');
exports = { something: 1 };
`,
		"/** @const */ var module$exports$ns$a = { /** @const */ something : 1 };")
}

// Java: ClosureRewriteModuleTest.testExport5()
func (c *ClosureRewriteModuleTest) TestExport5() {
	c.Test(
		`goog.module('ns.a');
/** @typedef {string} */ var x;
exports.x = x;
`,
		`/** @const */ var module$exports$ns$a = {};
/** @typedef {string} */ module$exports$ns$a.x;
`)
}

// Java: ClosureRewriteModuleTest.testExport6()
func (c *ClosureRewriteModuleTest) TestExport6() {
	c.Test(
		`goog.module('ns.a');
/** @typedef {string} */ var x;
exports = { something: x };
`,
		`/** @const */ var module$exports$ns$a = {};
/** @typedef {string} */ module$exports$ns$a.something;
`)
}

// Java: ClosureRewriteModuleTest.testExport6_1()
func (c *ClosureRewriteModuleTest) TestExport6_1() {
	c.Test(
		`goog.module('ns.a');
/** @typedef {string} */ var x;
exports.something = x;
`,
		`/** @const */ var module$exports$ns$a = {};
/** @typedef {string} */ module$exports$ns$a.something;
`)
}

// Java: ClosureRewriteModuleTest.testExport7()
func (c *ClosureRewriteModuleTest) TestExport7() {
	c.Test(
		`goog.module('ns.a');
/** @constructor */
exports = function() {};
`,
		"/** @constructor @const */ var module$exports$ns$a = function() {};")
}

// Java: ClosureRewriteModuleTest.testExport9()
func (c *ClosureRewriteModuleTest) TestExport9() {
	// Doesn't legacy-to-binary bridge export a typedef.
	c.TestSame(
		`goog.provide('goog.ui.ControlContent');
/** @typedef {string} */ goog.ui.ControlContent;
`)
}

// Java: ClosureRewriteModuleTest.testExport10()
func (c *ClosureRewriteModuleTest) TestExport10() {
	// Doesn't rewrite exports in legacy scripts.
	c.TestSame(
		`(function() {
  /** @constructor */ function S(string) {}
  exports.S = S;
})();
`)
}

// Java: ClosureRewriteModuleTest.testExport11()
func (c *ClosureRewriteModuleTest) TestExport11() {
	// Does rewrite export typedefs and defensively creates the exports root object first.
	c.Test(
		`goog.module('a.B');
/** @typedef {string} */ exports.C;
`,
		`/** @const */ var module$exports$a$B = {};
/** @const @typedef {string} */ module$exports$a$B.C;
`)
}

// Java: ClosureRewriteModuleTest.testExport13()
func (c *ClosureRewriteModuleTest) TestExport13() {
	// Creates the exports root object before export object reads.
	c.Test(
		`goog.module('a.B');
var field = exports;
`,
		`/** @const */ var module$exports$a$B = {};
var module$contents$a$B_field = module$exports$a$B;
`)
}

// Java: ClosureRewriteModuleTest.testExport_initializedWithVar()
func (c *ClosureRewriteModuleTest) TestExport_initializedWithVar() {
	c.Test(
		// TODO(lharker): should `var exports = ...` be banned?
		`goog.module('ns.a');
/** @suppress {checkTypes} */
// this statement causes a "type mismatch" error
var exports = {};
`,
		`/** @const */ var module$exports$ns$a = {};
/** @suppress {checkTypes} */
var module$contents$ns$a_exports = {};
`)
}

// Java: ClosureRewriteModuleTest.testExport_dontMangleLocalVariableNamedExports()
func (c *ClosureRewriteModuleTest) TestExport_dontMangleLocalVariableNamedExports() {
	c.Test(
		`goog.module('ns.a');

function f(exports, a) {
// test the various syntactic froms of doing goog.module exports
  exports.prop = 0;
  exports = {a};
  exports = function() {};
  if (true) {
    const exports = {};
  }
  return exports;
}
`,
		`/** @const */
var module$exports$ns$a = {};

function module$contents$ns$a_f(exports, a) {
  exports.prop = 0;
  exports = {a};
  exports = function() {};
  if (true) {
    const exports = {};
  }
  return exports;
}
`)
}

// Java: ClosureRewriteModuleTest.testExportEnhancedObjectLiteral()
func (c *ClosureRewriteModuleTest) TestExportEnhancedObjectLiteral() {
	c.Test(
		`goog.module('ns.a');
class Something {}
exports = { Something };
`,
		`/** @const */ var module$exports$ns$a = {};
module$exports$ns$a.Something = class {};
`)

	c.TestError(
		`goog.module('ns.a');
exports = { [something]: 3 };
`,
		ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY)
}

// Java: ClosureRewriteModuleTest.testImport()
func (c *ClosureRewriteModuleTest) TestImport() {
	// A goog.module() that imports, jsdocs, and uses both another goog.module() and a legacy
	// script.
	c.Test(
		c.Srcs(
			`goog.module('p.A');
/** @constructor */ function A() {}
exports = A;
`,
			`goog.provide('p.B');
/** @constructor */ p.B = function() {}
`,
			`goog.module('p.C');
var A = goog.require('p.A');
var B = goog.require('p.B');
function main() {
  /** @type {A} */ var a = new A;
  /** @type {B} */ var b = new B;
}
`),
		c.Expected(
			`/** @constructor */ function module$contents$p$A_A() {}
/** @const */ var module$exports$p$A = module$contents$p$A_A;
`,
			`goog.provide('p.B');
/** @constructor */ p.B = function() {}
`,
			`/** @const */ var module$exports$p$C = {};
goog.require('p.B');
function module$contents$p$C_main() {
  /** @type {module$exports$p$A} */ var a = new module$exports$p$A;
  /** @type {p.B} */ var b = new p.B;
}
`))
}

// Java: ClosureRewriteModuleTest.testSetTestOnly()
func (c *ClosureRewriteModuleTest) TestSetTestOnly() {
	c.Test(
		`goog.module('ns.a');
goog.setTestOnly();
`,
		`/** @const */ var module$exports$ns$a = {};
goog.setTestOnly();
`)
}

// Java: ClosureRewriteModuleTest.testRewriteJsDocImportedType()
func (c *ClosureRewriteModuleTest) TestRewriteJsDocImportedType() {
	c.Test(
		c.Srcs(
			`goog.module('p.A');
/** @constructor */
function A() {}
exports = A;
`,
			`goog.module('p.B');
/** @constructor */
function B() {}
exports = B;
`,
			`goog.module('p.C');
var A = goog.require('p.A');
var B = goog.requireType('p.B');
function main() {
  /** @type {A} */
  var a = new A;
  /** @type {B} */
  var b;
}
`),
		c.Expected(
			`/** @constructor */ function module$contents$p$A_A() {}
/** @const */ var module$exports$p$A = module$contents$p$A_A;
`,
			`/** @constructor */ function module$contents$p$B_B() {}
/** @const */ var module$exports$p$B = module$contents$p$B_B;
`,
			`/** @const */ var module$exports$p$C = {};
function module$contents$p$C_main() {
  /** @type {module$exports$p$A} */
  var a = new module$exports$p$A;
  /** @type {module$exports$p$B} */
  var b;
}
`))
}

// Java: ClosureRewriteModuleTest.testRewriteJsDocOwnType()
func (c *ClosureRewriteModuleTest) TestRewriteJsDocOwnType() {
	c.Test(
		c.Srcs(
			`goog.module('p.b');
/** @constructor */
function B() {}
function main() {
  /** @type {B} */
  var b = new B;
}
`),
		c.Expected(
			`/** @const */ var module$exports$p$b = {};
/** @constructor */
function module$contents$p$b_B() {}
function module$contents$p$b_main() {
  /** @type {module$contents$p$b_B} */
  var b = new module$contents$p$b_B;
}
`))
}

// Java: ClosureRewriteModuleTest.testRewriteJsDocFullyQualified()
func (c *ClosureRewriteModuleTest) TestRewriteJsDocFullyQualified() {
	c.Test(
		c.Srcs(
			`goog.module('p.A');
/** @constructor */
function A() {}
exports = A;
`,
			`goog.module('p.B');
/** @constructor */
function B() {}
exports = B;
`,
			`goog.module('p.C');
var A = goog.require('p.A');
var B = goog.requireType('p.B');
function main() {
  /** @type {p.A} */
  var a = new A;
  /** @type {p.B} */
  var b;
}
`),
		c.Expected(
			`/** @constructor */
function module$contents$p$A_A() {}
/** @const */ var module$exports$p$A = module$contents$p$A_A;
`,
			`/** @constructor */
function module$contents$p$B_B() {}
/** @const */ var module$exports$p$B = module$contents$p$B_B;
`,
			`/** @const */ var module$exports$p$C = {};
function module$contents$p$C_main() {
  /** @type {module$exports$p$A} */
  var a = new module$exports$p$A;
  /** @type {module$exports$p$B} */
  var b;
}
`))
}

// Java: ClosureRewriteModuleTest.testRewriteJsDocCircularReferenceWithRequireType()
func (c *ClosureRewriteModuleTest) TestRewriteJsDocCircularReferenceWithRequireType() {
	c.Test(
		c.Srcs(
			`goog.module('p.A');
goog.requireType('p.B');
/** @constructor */
function A() {}
A.prototype.setB = function(/** p.B */ x) {}
exports = A;
`,
			`goog.module('p.B');
var A = goog.require('p.A');
/** @constructor @extends {A} */
function B() {}
B.prototype = new A;
exports = B;
`),
		c.Expected(
			`/** @constructor */
function module$contents$p$A_A() {}
module$contents$p$A_A.prototype.setB = function(/** module$exports$p$B */ x) {}
/** @const */ var module$exports$p$A = module$contents$p$A_A;
`,
			`/** @constructor @extends {module$exports$p$A} */
function module$contents$p$B_B() {}
module$contents$p$B_B.prototype = new module$exports$p$A;
/** @const */ var module$exports$p$B = module$contents$p$B_B;
`))
}

// Java: ClosureRewriteModuleTest.testRewriteJsDocCircularReferenceWithoutRequireType()
func (c *ClosureRewriteModuleTest) TestRewriteJsDocCircularReferenceWithoutRequireType() {
	c.Test(
		c.Srcs(
			`goog.module('p.A');
/** @constructor */
function A() {}
A.prototype.setB = function(/** p.B */ x) {}
exports = A;
`,
			`goog.module('p.B');
var A = goog.require('p.A');
/** @constructor @extends {A} */
function B() {}
B.prototype = new A;
exports = B;
`),
		c.Expected(
			`/** @constructor */
function module$contents$p$A_A() {}
module$contents$p$A_A.prototype.setB = function(/** module$exports$p$B */ x) {}
/** @const */ var module$exports$p$A = module$contents$p$A_A;
`,
			`/** @constructor @extends {module$exports$p$A} */
function module$contents$p$B_B() {}
module$contents$p$B_B.prototype = new module$exports$p$A;
/** @const */ var module$exports$p$B = module$contents$p$B_B;
`))
}

// Java: ClosureRewriteModuleTest.testRewriteJsDocOwnTypeExported()
func (c *ClosureRewriteModuleTest) TestRewriteJsDocOwnTypeExported() {
	c.Test(
		`goog.module('p.A');

/** @constructor */
function A() {}

/** @type {!A} */
var x = new A;

exports = A;
`,
		`/** @constructor */
function module$contents$p$A_A() {}
/** @type {!module$contents$p$A_A} */
var module$contents$p$A_x = new module$contents$p$A_A;
/** @const */ var module$exports$p$A = module$contents$p$A_A;
`)
}

// Java: ClosureRewriteModuleTest.testDuplicateModuleDoesntCrash()
func (c *ClosureRewriteModuleTest) TestDuplicateModuleDoesntCrash() {
	c.IgnoreWarnings(ClosurePrimitiveErrors_DUPLICATE_MODULE, ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT)
	c.Test(
		c.Srcs("goog.module('ns.a');", "goog.module('ns.a');"),
		c.Expected(
			"/** @const */ var module$exports$ns$a = {};",
			"/** @const */ var module$exports$ns$a = {};"))
}

// Java: ClosureRewriteModuleTest.testDuplicateRequireDoesntCrash()
func (c *ClosureRewriteModuleTest) TestDuplicateRequireDoesntCrash() {
	c.IgnoreWarnings(DiagnosticGroups_DUPLICATE_VARS)
	c.Test(
		c.Srcs(
			`goog.module('ns.b');
exports = 1;
`,
			`goog.module('ns.a');
const b = goog.require('ns.b');
const b = goog.require('ns.b');
console.log(b);
`),
		c.Expected(
			`/** @const */
var module$exports$ns$b = 1;
`,
			`/** @const */
var module$exports$ns$a = {};
console.log(module$exports$ns$b);
`))
}

// Java: ClosureRewriteModuleTest.testDuplicateBundledModuleDoesntCrash()
func (c *ClosureRewriteModuleTest) TestDuplicateBundledModuleDoesntCrash() {
	c.IgnoreWarnings(ClosurePrimitiveErrors_DUPLICATE_MODULE, ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT)
	bundledModule := "goog.loadModule(function(exports) { goog.module('ns.a'); return exports; });"
	c.Test(
		c.Srcs(bundledModule, bundledModule),
		c.Expected(
			"/** @const */ var module$exports$ns$a = {};",
			"/** @const */ var module$exports$ns$a = {};"))
}

// Java: ClosureRewriteModuleTest.testDuplicateNamespaceDoesntCrash()
func (c *ClosureRewriteModuleTest) TestDuplicateNamespaceDoesntCrash() {
	// The compiler emits a warning elsewhere for this code
	c.TestError(
		c.Srcs("goog.module('ns.a');", "goog.provide('ns.a');"), ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE)
}

// Java: ClosureRewriteModuleTest.testMisplacedGoogModuleDoesntCrash()
func (c *ClosureRewriteModuleTest) TestMisplacedGoogModuleDoesntCrash() {
	// The compiler emits a warning elsewhere for this code
	c.Test(
		c.Srcs("alert(goog.module('ns.a'));"), //
		c.Expected("alert(void 0);"))
}

// Java: ClosureRewriteModuleTest.testImportInliningDoesntShadow()
func (c *ClosureRewriteModuleTest) TestImportInliningDoesntShadow() {
	c.TestNoWarning(
		`goog.provide('a.b.c');
/** @const */ var a = a || {};
a.b.c = class {};
goog.loadModule(function(exports) { 'use strict';
  goog.module('a.b.d');
  goog.module.declareLegacyNamespace();
  var c = goog.require('a.b.c');
  exports.c = new c;
  return exports;
});
`)

	c.TestNoWarning(
		`goog.provide('a.b.c');
/** @const */ var a = a || {};
a.b.c = class {};
goog.loadModule(function(exports) { 'use strict';
  goog.module('a.b.d');
  goog.module.declareLegacyNamespace();
  var c = goog.requireType('a.b.c');
  exports.c = new c;
  return exports;
});
`)
}

// Java: ClosureRewriteModuleTest.testImportInliningShadowsVar()
func (c *ClosureRewriteModuleTest) TestImportInliningShadowsVar() {

	c.TestError(
		c.Srcs(
			`goog.provide('a.b.c');
a.b.c = 5;
`,
			`goog.module('a.b.d');
var c = goog.require('a.b.c');
function foo() {
  var a = 10;
  var b = c;
}
`),
		ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR)

	c.TestError(
		c.Srcs(
			`goog.provide('a.b.c');
a.b.c = 5;
`,
			`goog.module('a.b.d');
var c = goog.requireType('a.b.c');
function foo() {
  var a = 10;
  var b = c;
}
`),
		ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR)
}

// Java: ClosureRewriteModuleTest.testImportInliningShadowsDestructuredImport()
func (c *ClosureRewriteModuleTest) TestImportInliningShadowsDestructuredImport() {
	c.TestError(
		c.Srcs(
			`goog.provide('a.b.c');
a.b.c.d = 5;
`,
			`goog.module('a.b.d');
const {d} = goog.require('a.b.c');
function foo() {
  var a = 10;
  var b = d;
}
`),
		ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR)

	c.TestError(
		c.Srcs(
			`goog.provide('a.b.c');
a.b.c.d = 5;
`,
			`goog.module('a.b.d');
const {d} = goog.requireType('a.b.c');
function foo() {
  var a = 10;
  var b = d;
}
`),
		ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR)
}

// Java: ClosureRewriteModuleTest.testExportsShadowingAllowed()
func (c *ClosureRewriteModuleTest) TestExportsShadowingAllowed() {
	c.TestNoWarning(
		`goog.loadModule(function(exports) {
   goog.module('a.b.c');

   class Foo {}
   /** @const {*} */ Foo.prototype.x;
   exports.Foo = Foo;

   return exports;
});
`)
}

// Java: ClosureRewriteModuleTest.testExportRewritingShadows()
func (c *ClosureRewriteModuleTest) TestExportRewritingShadows() {
	c.Test(
		`goog.module('a.b.c');
function test() {}
function f(test) { return test; }
exports = test;
`,
		`function module$contents$a$b$c_test() {}
function module$contents$a$b$c_f(test) { return test; }
/** @const */ var module$exports$a$b$c = module$contents$a$b$c_test;
`)

	c.Test(
		`goog.module('a.b.c');
function test() {}
function f(test) { return test; }
exports.test = test;
`,
		`/** @const */ var module$exports$a$b$c = {};
function module$contents$a$b$c_test() {}
function module$contents$a$b$c_f(test) { return test; }
/** @const */ module$exports$a$b$c.test = module$contents$a$b$c_test;
`)
}

// Java: ClosureRewriteModuleTest.testEarlyRequireModule()
func (c *ClosureRewriteModuleTest) TestEarlyRequireModule() {
	c.Test(
		c.Srcs(
			`goog.module('ns.a');
goog.require('ns.b')
`,
			"goog.module('ns.b');"))
}

// Java: ClosureRewriteModuleTest.testEarlyRequireLegacyScript()
func (c *ClosureRewriteModuleTest) TestEarlyRequireLegacyScript() {
	c.Test(
		c.Srcs(
			`goog.module('ns.a');
goog.require('ns.b')
`,
			"goog.provide('ns.b');"))
}

// Java: ClosureRewriteModuleTest.testEarlyRequireTypeModule()
func (c *ClosureRewriteModuleTest) TestEarlyRequireTypeModule() {
	c.Test(
		c.Srcs(
			`goog.module('ns.a');
goog.requireType('ns.b')
`,
			"goog.module('ns.b');"))
}

// Java: ClosureRewriteModuleTest.testEarlyRequireTypeLegacyScript()
func (c *ClosureRewriteModuleTest) TestEarlyRequireTypeLegacyScript() {
	c.Test(
		c.Srcs(
			`goog.module('ns.a');
goog.requireType('ns.b')
`,
			"goog.provide('ns.b');"))
}

// Java: ClosureRewriteModuleTest.testValidEarlyGoogModuleGet()
func (c *ClosureRewriteModuleTest) TestValidEarlyGoogModuleGet() {
	// Legacy Script to Module goog.module.get.
	c.Test(
		c.Srcs(
			`goog.provide('ns.a');
function foo() {
  var b = goog.module.get('ns.b');
}
`,
			"goog.module('ns.b');"),
		c.Expected(
			"goog.provide('ns.a'); function foo() { var b = module$exports$ns$b; }",
			"/** @const */ var module$exports$ns$b = {};"))
}

// Java: ClosureRewriteModuleTest.testInnerScriptOuterModule()
func (c *ClosureRewriteModuleTest) TestInnerScriptOuterModule() {
	// Rewrites fully qualified JsDoc references to types but without writing a prefix as
	// module$exports when there's a longer prefix that references a script.

	// Test when the prefix is chunk unit longer than the module name
	c.Test(
		c.Srcs(
			`goog.module('A');
/** @constructor */
function A() {}
exports = A;
`,
			`goog.provide('A.B');
/** @constructor */
A.B = function () {}
function main() {
  /** @type {A.B} */
  var l = new A.B();
}
`),
		c.Expected(
			`/** @constructor */ function module$contents$A_A() {}
/** @const */ var module$exports$A = module$contents$A_A;
`,
			`goog.provide('A.B');
/** @constructor */
A.B = function() {};
function main() {
// Note A.B was NOT written to module$exports$A.B
  /** @type {A.B} */
  var l = new A.B();
}
`))

	// Test when the prefix is much longer than the module name
	c.Test(
		c.Srcs(
			`goog.module('A');
/** @constructor */
function A() {}
exports = A;
`,
			`goog.provide('A.b.c.D');
/** @constructor */
A.b.c.L = function () {}
function main() {
  /** @type {A.b.c.L} */
  var l = new A.b.c.L();
}
`),
		c.Expected(
			`/** @constructor */ function module$contents$A_A() {}
/** @const */ var module$exports$A = module$contents$A_A;
`,
			`goog.provide('A.b.c.D');
/** @constructor */
A.b.c.L = function() {}; // Note L not D
function main() {
// Note A.b.c.L was NOT written to module$exports$A.b.c.L.
  /** @type {A.b.c.L} */
  var l = new A.b.c.L();
}
`))
}

// Java: ClosureRewriteModuleTest.testModuleLevelVars()
func (c *ClosureRewriteModuleTest) TestModuleLevelVars() {
	c.Test(
		`goog.module('b.c.c');
/** @const */
var F = 0;
`,
		`/** @const */ var module$exports$b$c$c = {};
/** @const */ var module$contents$b$c$c_F = 0;
`)
}

// Java: ClosureRewriteModuleTest.testPublicExport()
func (c *ClosureRewriteModuleTest) TestPublicExport() {
	c.Test(
		`goog.module('a.b.c');
goog.module.declareLegacyNamespace();
/** @public */ exports = 5;
`,
		`goog.provide('a.b.c');
/** @const @public */ a.b.c = 5;
`)
}

// Java: ClosureRewriteModuleTest.testReferenceToNonLegacyGoogModuleName()
func (c *ClosureRewriteModuleTest) TestReferenceToNonLegacyGoogModuleName() {
	c.Test(
		c.Srcs("goog.module('a.b.c');", "use(a.b.c);"),
		c.Expected("/** @const */ var module$exports$a$b$c={}", "use(a.b.c);"),
		c.Warning(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY),
		c.Warning(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY))
}

// Java: ClosureRewriteModuleTest.testGoogModuleValidReferences1()
func (c *ClosureRewriteModuleTest) TestGoogModuleValidReferences1() {
	c.Test(
		c.Srcs(
			"goog.module('a.b.c');",
			`goog.module('x.y.z');
var c = goog.require('a.b.c');
use(c);
`),
		c.Expected(
			"/** @const */ var module$exports$a$b$c={};",
			`/** @const */ var module$exports$x$y$z={};
use(module$exports$a$b$c);
`))
}

// Java: ClosureRewriteModuleTest.testGoogModuleValidReferences2()
func (c *ClosureRewriteModuleTest) TestGoogModuleValidReferences2() {
	c.Test(
		c.Srcs(
			"goog.module('a.b.c');",
			`goog.require('a.b.c');
goog.scope(function() {
  var c = goog.module.get('a.b.c');
  use(c);
});
`),
		c.Expected(
			"/** @const */ var module$exports$a$b$c={};",
			"goog.scope(function() { var c = module$exports$a$b$c; use(c); });"))
}

// Java: ClosureRewriteModuleTest.testLegacyGoogModuleValidReferences()
func (c *ClosureRewriteModuleTest) TestLegacyGoogModuleValidReferences() {
	c.Test(
		c.Srcs(
			"goog.module('a.b.c'); goog.module.declareLegacyNamespace();",
			"goog.require('a.b.c'); use(a.b.c);"),
		c.Expected("goog.provide('a.b.c');", "goog.require('a.b.c'); use(a.b.c);"))

	c.Test(
		c.Srcs(
			"goog.module('a.b.c'); goog.module.declareLegacyNamespace();",
			"goog.module('x.y.z'); var c = goog.require('a.b.c'); use(c);"),
		c.Expected(
			"goog.provide('a.b.c');",
			"/** @const */ var module$exports$x$y$z={}; goog.require('a.b.c'); use(a.b.c);"))

	c.Test(
		c.Srcs(
			`goog.module('a.b.Foo');
goog.module.declareLegacyNamespace();

/** @constructor */ exports = function() {};
`,
			"/** @param {a.b.Foo} x */ function f(x) {}"),
		c.Expected(
			"goog.provide('a.b.Foo'); /** @constructor @const */ a.b.Foo = function() {};",
			"/** @param {a.b.Foo} x */ function f(x) {}"))

	c.Test(
		c.Srcs(
			`goog.module('a.b.c');
goog.module.declareLegacyNamespace();

exports = function() {};
`,
			"function f() { return goog.module.get('a.b.c'); }"),
		c.Expected(
			"goog.provide('a.b.c'); /** @const */ a.b.c = function() {};",
			"function f() { return a.b.c; }"))

	c.Test(
		c.Srcs(
			`goog.module('a.b.Foo');
goog.module.declareLegacyNamespace();

/** @constructor */ function Foo() {}

exports = Foo;
`,
			"/** @param {a.b.Foo} x */ function f(x) {}"),
		c.Expected(
			`goog.provide('a.b.Foo');
/** @constructor */ function module$contents$a$b$Foo_Foo() {}
/** @const */ a.b.Foo = module$contents$a$b$Foo_Foo;
`,
			"/** @param {a.b.Foo} x */ function f(x) {}"))

	c.Test(
		c.Srcs(
			`goog.module('a.b');
goog.module.declareLegacyNamespace();

/** @constructor */ function Foo() {};

exports.Foo = Foo;
`,
			"/** @param {a.b.Foo} x */ function f(x) {}"),
		c.Expected(
			`goog.provide('a.b');
/** @constructor */ function module$contents$a$b_Foo() {};
/** @const */ a.b.Foo = module$contents$a$b_Foo;
`,
			"/** @param {a.b.Foo} x */ function f(x) {}"))
}

// Java: ClosureRewriteModuleTest.testUselessUseStrict()
func (c *ClosureRewriteModuleTest) TestUselessUseStrict() {
	c.TestWarning(
		"'use strict'; goog.module('b.c.c');", ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE)
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleAliases1()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleAliases1() {
	c.Test(
		c.Srcs(
			`goog.module('base');

/** @constructor */ var Base = function() {}
exports = Base;
`,
			`goog.module('leaf');

var Base = goog.require('base');
exports = /** @constructor @extends {Base} */ function Foo() {}
`),
		c.Expected(
			"/** @constructor */ var module$exports$base = function() {};",
			`/** @const */ var module$exports$leaf =
/** @constructor @extends {module$exports$base} */ function Foo() {}
`))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleAliases2()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleAliases2() {
	c.Test(
		c.Srcs(
			`goog.module('ns.base');

/** @constructor */ var Base = function() {}
exports = Base;
`,
			`goog.module('leaf');

var Base = goog.require('ns.base');
exports = /** @constructor @extends {Base} */ function Foo() {}
`),
		c.Expected(
			"/** @constructor */ var module$exports$ns$base = function() {};",
			`/** @const */ var module$exports$leaf =
/** @constructor @extends {module$exports$ns$base} */ function Foo() {}
`))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleAliases3()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleAliases3() {
	c.Test(
		c.Srcs(
			`goog.module('ns.base');

/** @constructor */ var Base = function() {};
/** @constructor */ Base.Foo = function() {};
exports = Base;
`,
			`goog.module('leaf');

var Base = goog.require('ns.base');
exports = /** @constructor @extends {Base.Foo} */ function Foo() {}
`),
		c.Expected(
			`/** @constructor */ var module$exports$ns$base = function() {};
/** @constructor */ module$exports$ns$base.Foo = function() {};
`,
			`/** @const */ var module$exports$leaf =
/** @constructor @extends {module$exports$ns$base.Foo} */ function Foo() {}
`))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleAliases4()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleAliases4() {
	c.Test(
		c.Srcs(
			`goog.module('ns.base');

/** @constructor */ var Base = function() {}
exports = Base;
`,
			`goog.module('leaf');

var Base = goog.require('ns.base');
exports = new Base;
`),
		c.Expected(
			"/** @constructor */ var module$exports$ns$base = function() {};",
			"/** @const */ var module$exports$leaf = new module$exports$ns$base;"))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleAliases5()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleAliases5() {
	c.Test(
		c.Srcs(
			`goog.module('ns.base');

/** @constructor */ var Base = function() {}
exports = Base;
`,
			`goog.module('mid');

var Base = goog.require('ns.base');
exports = Base;
`,
			`goog.module('leaf')
var Base = goog.require('mid');
new Base;
`),
		c.Expected(
			"/** @constructor */ var module$exports$ns$base = function() {};",
			"/** @const */ var module$exports$mid = module$exports$ns$base;",
			"/** @const */ var module$exports$leaf = {}; new module$exports$mid;"))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleAliases6()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleAliases6() {
	c.Test(
		c.Srcs(
			`goog.module('base');

/** @constructor */ exports.Foo = function() {};
`,
			`goog.module('FooWrapper');

const {Foo} = goog.require('base');
exports = Foo;
`),
		c.Expected(
			`/** @const */ var module$exports$base = {};
/** @constructor @const */ module$exports$base.Foo = function() {};
`,
			"/** @const */ var module$exports$FooWrapper = module$exports$base.Foo;"))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleAliases7()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleAliases7() {
	c.Test(
		c.Srcs(
			`goog.module('base');

/** @constructor */ exports.Foo = function() {};
`,
			`goog.module('FooWrapper');

const {Foo: FooFromBaseModule} = goog.require('base');
exports = FooFromBaseModule;
`),
		c.Expected(
			`/** @const */ var module$exports$base = {};
/** @constructor @const */ module$exports$base.Foo = function() {};
`,
			"/** @const */ var module$exports$FooWrapper = module$exports$base.Foo;"))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleTypes_localShadowsModuleScopeVar()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleTypes_localShadowsModuleScopeVar() {
	c.Test(
		`goog.module('client');
class Foo {}
{
  class Foo {}
  let /** !Foo */ x;
}
`,
		`/** @const */ var module$exports$client = {};
class module$contents$client_Foo {}
{
  class Foo {}
  let /** !Foo */ x;
}
`)
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleTypes_localShadowsProvide()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleTypes_localShadowsProvide() {
	c.Test(
		c.Srcs(
			"goog.provide('Foo.Child');",
			"goog.module('Bar'); exports.Child = class {};",
			`goog.module('client');
const Foo = goog.require('Bar');
let /** !Foo.Child */ myStr;
`), // Foo.Child refers to the local, not the provide.
		c.Expected(
			"goog.provide('Foo.Child');",
			`/** @const */ var module$exports$Bar = {};
/** @const */ module$exports$Bar.Child = class {};
`,
			`/** @const */ var module$exports$client = {};
let /** !module$exports$Bar.Child */ module$contents$client_myStr;
`))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleTypes_localShadowsModule()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleTypes_localShadowsModule() {
	c.Test(
		c.Srcs(
			"goog.module('Foo.Child');",
			"goog.module('Bar'); exports.Child = class {};",
			`goog.module('client');
const Foo = goog.require('Bar');
let /** !Foo.Child */ myStr;
`), // Foo.Child refers to the local, not the module.
		c.Expected(
			"/** @const */ var module$exports$Foo$Child = {}",
			`/** @const */ var module$exports$Bar = {};
/** @const */ module$exports$Bar.Child = class {};
`,
			`/** @const */ var module$exports$client = {};
let /** !module$exports$Bar.Child */ module$contents$client_myStr;
`))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleAliasesWithPreservedPrimitives()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleAliasesWithPreservedPrimitives() {
	c.preserveClosurePrimitives = true
	// Need to disable tree comparison because compiler adds MODULE_BODY token when parsing
	// expected output but it is not present in actual tree.
	c.DisableCompareAsTree()
	c.Test(
		c.Srcs(
			`goog.module('foo');
/** @constructor */ exports = function() {};
`,
			`goog.module('bar');
exports.doBar = function() {};
`,
			`goog.module('baz');
exports.doBaz = function() {};
`,
			`goog.module('leaf1');
var Foo = goog.require('foo');
var {doBar} = goog.require('bar');
var {doBaz: doooBaz} = goog.require('baz');
`,
			`goog.module('leaf2');
var Foo = goog.requireType('foo');
var {doBar} = goog.requireType('bar');
var {doBaz: doooBaz} = goog.requireType('baz');
`),
		c.Expected(
			`goog.module("foo");
/** @const @constructor */
var module$exports$foo = function() {
};
`,
			`goog.module("bar");
/** @const */
var module$exports$bar = {};
/** @const */
module$exports$bar.doBar = function() {
};
`,
			`goog.module("baz");
/** @const */
var module$exports$baz = {};
/** @const */
module$exports$baz.doBaz = function() {
};
`,
			`goog.module("leaf1");
/** @const */
var module$exports$leaf1 = {};
var module$contents$leaf1_Foo = goog.require("foo");
var {doBar:module$contents$leaf1_doBar} = goog.require("bar");
var {doBaz:module$contents$leaf1_doooBaz} = goog.require("baz");
`,
			`goog.module("leaf2");
/** @const */
var module$exports$leaf2 = {};
var module$contents$leaf2_Foo = goog.requireType("foo");
var {doBar:module$contents$leaf2_doBar} = goog.requireType("bar");
var {doBaz:module$contents$leaf2_doooBaz} = goog.requireType("baz");
`))
}

// Java: ClosureRewriteModuleTest.testGoogModuleExportsProvidedName()
func (c *ClosureRewriteModuleTest) TestGoogModuleExportsProvidedName() {
	c.Test(
		c.Srcs(
			`goog.provide('Foo');

/** @constructor */ var Foo = function() {};
`,
			`goog.module('FooWrapper');

goog.require('Foo');

exports = Foo;
`),
		c.Expected(
			`goog.provide('Foo');

/** @constructor */ var Foo = function() {};
`,
			"goog.require('Foo'); /** @const */ var module$exports$FooWrapper = Foo;"))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleAliasesWithPrototypeGets1()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleAliasesWithPrototypeGets1() {
	c.Test(
		c.Srcs(
			`goog.module('mod_B');

/** @interface */ function B(){}
B.prototype.f = function(){};

exports = B;
`,
			`goog.module('mod_A');

var B = goog.require('mod_B');

/** @type {B} */
var b;
`),
		c.Expected(
			`/**@interface */ function module$contents$mod_B_B() {}
module$contents$mod_B_B.prototype.f = function() {};
/** @const */ var module$exports$mod_B = module$contents$mod_B_B;
`,
			`/** @const */ var module$exports$mod_A = {};
/**@type {module$exports$mod_B} */ var module$contents$mod_A_b;
`))
}

// Java: ClosureRewriteModuleTest.testRewriteGoogModuleAliasesWithPrototypeGets2()
func (c *ClosureRewriteModuleTest) TestRewriteGoogModuleAliasesWithPrototypeGets2() {
	c.Test(
		c.Srcs(
			`goog.module('mod_B');

/** @interface */ function B(){}

exports = B;
`,
			`goog.module('mod_A');

var B = goog.require('mod_B');
B.prototype;

/** @type {B} */
var b;
`),
		c.Expected(
			`/**@interface */ function module$contents$mod_B_B() {}
/** @const */ var module$exports$mod_B = module$contents$mod_B_B;
`,
			`/** @const */ var module$exports$mod_A = {}
module$exports$mod_B.prototype;
/**@type {module$exports$mod_B} */ var module$contents$mod_A_b;
`))
}

// Java: ClosureRewriteModuleTest.testLegacyModuleIsUninlined()
func (c *ClosureRewriteModuleTest) TestLegacyModuleIsUninlined() {
	c.Test(
		`goog.module('mod.ns');
goog.module.declareLegacyNamespace();

class Foo {}

exports.Foo = Foo;
`,
		`goog.provide('mod.ns');
class module$contents$mod$ns_Foo {}
/** @const */ mod.ns.Foo = module$contents$mod$ns_Foo;
`)
}

// Java: ClosureRewriteModuleTest.testLegacyModuleExportStillExported()
func (c *ClosureRewriteModuleTest) TestLegacyModuleExportStillExported() {
	c.Test(
		`goog.module('modA');
goog.module.declareLegacyNamespace();

class Foo {}
exports = { /** @export */ Foo};
`,
		`goog.provide('modA');
class module$contents$modA_Foo {}
/** @const @export */ modA.Foo = module$contents$modA_Foo;
`)
}

// Java: ClosureRewriteModuleTest.testMultiplyExportedSymbolDoesntCrash()
func (c *ClosureRewriteModuleTest) TestMultiplyExportedSymbolDoesntCrash() {
	c.Test(
		`goog.module('mod');

class Foo {}

exports.name1 = Foo;
exports.name2 = Foo;
`,
		`/** @const */ var module$exports$mod = {};
module$exports$mod.name1 = class {};
/** @const */ module$exports$mod.name2 = module$exports$mod.name1;
`)
}

// Java: ClosureRewriteModuleTest.testIjsFileInExterns()
func (c *ClosureRewriteModuleTest) TestIjsFileInExterns() {
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	c.AllowExternsChanges()
	c.TestNoWarning(
		c.Externs(
			`/** @typeSummary */
goog.module('mod_B');

/** @interface */ function B(){}

exports = B;
`),
		c.Srcs(
			`goog.module('mod_A');

var B = goog.require('mod_B');

/** @constructor @implements {B} */
function A() {}
`))

	c.TestNoWarning(
		c.Externs(
			`/** @typeSummary */
goog.loadModule(function(exports) { 'use strict';
  goog.module('mod_B');

  /** @interface */ function B(){}

  exports = B;
  return exports;
});
`),
		c.Srcs(
			`goog.module('mod_A');

var B = goog.require('mod_B');

/** @constructor @implements {B} */
function A() {}
`))

	c.TestNoWarning(
		c.Externs(
			`/** @typeSummary */
goog.loadModule(function(exports) { 'use strict';
  goog.module('mod_B');

  /** @interface */ function B(){}

  exports.B = B;
  return exports;
});
`),
		c.Srcs(
			`goog.module('mod_A');

var {B} = goog.require('mod_B');

/** @constructor @implements {B} */
function A() {}
`))
}

// Java: ClosureRewriteModuleTest.addFreeCallToNamedExports()
func (c *ClosureRewriteModuleTest) Test_addFreeCallToNamedExports() {
	c.DisableCompareAsTree() // necessary to compare 'free call' marking
	c.Test(
		c.Srcs(
			`goog.module('mod');
exports.fn = function() {};
`,
			`goog.module('client');
const mod = goog.require('mod');
mod.fn();
`),
		c.Expected(
			`/** @const */
var module$exports$mod = {};
/** @const */
module$exports$mod.fn = function() {
};
/** @const */
var module$exports$client = {};
(0,module$exports$mod.fn)();
`))
}

// Java: ClosureRewriteModuleTest.addFreeCallToNamedExportsLiteral()
func (c *ClosureRewriteModuleTest) Test_addFreeCallToNamedExportsLiteral() {
	c.DisableCompareAsTree() // necessary to compare 'free call' marking
	c.Test(
		c.Srcs(
			`goog.module('mod');
const fn = function() {};
exports = {fn}
`,
			`goog.module('client');
const mod = goog.require('mod');
mod.fn();
`),
		c.Expected(
			`/** @const */
var module$exports$mod = {};
/** @const */
module$exports$mod.fn = function() {
};
/** @const */
var module$exports$client = {};
(0,module$exports$mod.fn)();
`))
}

// Java: ClosureRewriteModuleTest.dontAddFreeCallToPropertiesOnDefaultExports()
func (c *ClosureRewriteModuleTest) Test_dontAddFreeCallToPropertiesOnDefaultExports() {
	c.Test(
		c.Srcs(
			`goog.module('mod.Bar');
/** @const */
exports = class Bar {
  static fn() { use(this); }
}
`,
			`goog.module('client');
const Bar = goog.require('mod.Bar');
Bar.fn();
`),
		c.Expected(
			`/** @const */
var module$exports$mod$Bar = class Bar {
  static fn() { use(this); }
}
`,
			`/** @const */
var module$exports$client = {};
// preserve the 'this' value of 'class Bar {' - don't make this a free call.
module$exports$mod$Bar.fn();
`))
}

// Java: ClosureRewriteModuleTest.dontAddFreeCallIfModuleNotCallee()
func (c *ClosureRewriteModuleTest) Test_dontAddFreeCallIfModuleNotCallee() {
	c.Test(
		c.Srcs(
			`goog.module('mod.Bar');
exports.NAME = 'BAR';
`,
			`goog.module('client');
const Bar = goog.require('mod.Bar');
alert(Bar.NAME);
`),
		c.Expected(
			`/** @const */
var module$exports$mod$Bar = {};
/** @const */
module$exports$mod$Bar.NAME = 'BAR'
`,
			`/** @const */
var module$exports$client = {};
// don't mistake the Bar.NAME argument for the callee
alert(module$exports$mod$Bar.NAME);
`))
}

// Java: ClosureRewriteModuleTest.testTypeAndSourceInfoOfGoogRequireFromModule()
func (c *ClosureRewriteModuleTest) TestTypeAndSourceInfoOfGoogRequireFromModule() {
	c.Test(
		c.Srcs(
			`goog.module('mod.one');
exports.Bar = class {};
`,
			`goog.module('mod.two');
const {Bar} = goog.require('mod.one');
new Bar();
`),
		c.Expected(
			`/** @const */
var module$exports$mod$one = {};
/** @const */
module$exports$mod$one.Bar = class {};
`,
			`/** @const */
var module$exports$mod$two = {};
new module$exports$mod$one.Bar();
`))

	// Verify the type of $module$exports$mod$one.Bar is correct
	secondScript := c.GetLastCompiler().GetJsRoot().GetSecondChild()
	moduleExportsDotBar := NodeUtil_FindPreorder(secondScript, func(node *rhino.Node) bool { return node.MatchesQualifiedName("module$exports$mod$one.Bar") }, func(*rhino.Node) bool { return true })

	AssertNode(moduleExportsDotBar).HasJSTypeThat().GetReferenceNameIsEqualTo("mod.one.Bar")
	// The source info for the rewritten name must match the source info for `Bar` in `new Bar();`
	AssertNode(moduleExportsDotBar).
		HasSourceFileName(secondScript.GetSourceFileName()).
		HasLineno(3).
		HasCharno(4).
		HasLength(3)
	AssertNode(moduleExportsDotBar.GetOnlyChild()).HasEqualSourceInfoTo(moduleExportsDotBar)
}

// Java: ClosureRewriteModuleTest.testTypeOfGoogRequireFromLegacyModule()
func (c *ClosureRewriteModuleTest) TestTypeOfGoogRequireFromLegacyModule() {
	c.Test(
		c.Srcs(
			`goog.module('mod.one');
goog.module.declareLegacyNamespace();
exports.Bar = class {};
`,
			`goog.module('mod.two');
const {Bar} = goog.require('mod.one');
new Bar();
`),
		c.Expected(
			`goog.provide('mod.one');
/** @const */
mod.one.Bar = class {};
`,
			`/** @const */
var module$exports$mod$two = {};
new mod.one.Bar();
`))

	// Verify the type of mod.one.Bar is correct
	secondScript := c.GetLastCompiler().GetJsRoot().GetSecondChild()
	moduleExportsDotBar := NodeUtil_FindPreorder(secondScript, func(node *rhino.Node) bool { return node.MatchesQualifiedName("mod.one.Bar") }, func(*rhino.Node) bool { return true })

	AssertNode(moduleExportsDotBar).HasJSTypeThat().GetReferenceNameIsEqualTo("mod.one.Bar")
}

// Java: ClosureRewriteModuleTest.rewriteLendsAnnotation()
func (c *ClosureRewriteModuleTest) Test_rewriteLendsAnnotation() {
	c.Test(
		`goog.module('mod');
class Foo {
  constructor() {
    this.x = /** @lends {Foo.prototype} */ {};
  }
}
`,
		`/** @const */ var module$exports$mod = {};
class module$contents$mod_Foo {
  constructor() {
    this.x = /** @lends {module$contents$mod_Foo.prototype} */  {};
  }
}
`)
}

// This pass only handles goog.modules. ES6 modules are left alone.
//
// Java: ClosureRewriteModuleTest.testEs6Module()
func (c *ClosureRewriteModuleTest) TestEs6Module() {
	c.TestSame("export var x;")
	c.TestSame("import {x} from '/y';")
}

// Java: ClosureRewriteModuleTest.testDuplicateExportOfName()
func (c *ClosureRewriteModuleTest) TestDuplicateExportOfName() {
	c.IgnoreWarnings(DiagnosticGroups_DUPLICATE_VARS)
	// ClosureCheckModule reports an error on this code, but that error is suppressed in some cases
	// so ClosureRewriteModule can't crash.
	c.Test(
		`goog.module('Foo');
class Foo {}
exports = Foo;
exports = Foo;
`,
		`class module$exports$Foo {
}
module$contents$Foo_Foo;
/** @const */ var module$exports$Foo = module$contents$Foo_Foo;
`)
}

// Java: ClosureRewriteModuleTest.testDuplicateExportOfExpressionUsesLatter()
func (c *ClosureRewriteModuleTest) TestDuplicateExportOfExpressionUsesLatter() {
	c.IgnoreWarnings(DiagnosticGroups_DUPLICATE_VARS)
	// ClosureCheckModule reports an error on this code, but that error is suppressed in some cases
	// so ClosureRewriteModule can't crash.
	c.Test(
		`goog.module('Foo');
exports = 0;
exports = 1;
`,
		`0;
/** @const */ var module$exports$Foo = 1;
`)
}

// Java: ClosureRewriteModuleTest.testTypeReferenceToDefaultExport()
func (c *ClosureRewriteModuleTest) TestTypeReferenceToDefaultExport() {
	c.Test(
		`goog.module('Foo');
exports = class Bar {};
/** @type {!exports} */
const e = new exports();
`,
		`/** @const */
var module$exports$Foo = class Bar {};
/** @type {!module$exports$Foo} */
const module$contents$Foo_e = new module$exports$Foo();
`)
}

// Java: ClosureRewriteModuleTest.testTypeReferenceToNamedExport()
func (c *ClosureRewriteModuleTest) TestTypeReferenceToNamedExport() {
	c.Test(
		`goog.module('Foo');
/** @enum {number} */
exports.Foo = {A: 1};
/** @type {!exports.Foo} */
const z = exports.Foo.A;
`,
		`/** @const */
var module$exports$Foo = {};
/** @const @enum {number} */
module$exports$Foo.Foo = {A: 1};
/** @type {!module$exports$Foo.Foo} */
const module$contents$Foo_z = module$exports$Foo.Foo.A;
`)
}

// Java: ClosureRewriteModuleTest.testLegacyNamespace_referenceToExports()
func (c *ClosureRewriteModuleTest) TestLegacyNamespace_referenceToExports() {
	c.Test(
		`goog.module('a.b.c');
goog.module.declareLegacyNamespace();
exports = class Foo {};
function useExports() {
  return exports;
}
`,
		`goog.provide('a.b.c');
/** @const */
a.b.c = class Foo {};
function module$contents$a$b$c_useExports() {
  return a.b.c;
}
`)
}

// Java: ClosureRewriteModuleTest.testLegacyNamespace_referenceToExportProperty()
func (c *ClosureRewriteModuleTest) TestLegacyNamespace_referenceToExportProperty() {
	c.Test(
		`goog.module('a.b.c');
goog.module.declareLegacyNamespace();
exports.Foo = class Foo {};
function useExports() {
  return exports.Foo;
}
`,
		`goog.provide('a.b.c');
/** @const */
a.b.c.Foo = class Foo {};
function module$contents$a$b$c_useExports() {
  return a.b.c.Foo;
}
`)
}

// Java: ClosureRewriteModuleTest.testUntypedLegacyNamespace_referenceToExports()
func (c *ClosureRewriteModuleTest) TestUntypedLegacyNamespace_referenceToExports() {
	c.runTypeCheck = false
	c.DisableTypeInfoValidation()
	c.Test(
		`goog.module('a.b.c');
goog.module.declareLegacyNamespace();
exports = class Foo {};
function useExports() {
  return exports;
}
`,
		`goog.provide('a.b.c');
/** @const */
a.b.c = class Foo {};
function module$contents$a$b$c_useExports() {
  return a.b.c;
}
`)
}

// Java: ClosureRewriteModuleTest.testUntypedLegacyNamespace_referenceToExportProperty()
func (c *ClosureRewriteModuleTest) TestUntypedLegacyNamespace_referenceToExportProperty() {
	c.runTypeCheck = false
	c.DisableTypeInfoValidation()
	c.Test(
		`goog.module('a.b.c');
goog.module.declareLegacyNamespace();
exports.Foo = class Foo {};
function useExports() {
  return exports.Foo;
}
`,
		`goog.provide('a.b.c');
/** @const */
a.b.c.Foo = class Foo {};
function module$contents$a$b$c_useExports() {
  return a.b.c.Foo;
}
`)
}

// Java: ClosureRewriteModuleTest.testUntypedGoogRequireDynamicInAwait_legacyNamespace()
func (c *ClosureRewriteModuleTest) TestUntypedGoogRequireDynamicInAwait_legacyNamespace() {
	c.runTypeCheck = false
	c.DisableTypeInfoValidation()
	c.Test(
		c.Srcs(
			`goog.module('a.b.c');
goog.module.declareLegacyNamespace();
exports = class Foo {};
`,
			`async function f() {
  const {Foo} = await goog.requireDynamic('a.b.c');
}
`),
		c.Expected(
			`goog.provide('a.b.c');
/** @const */
a.b.c = class Foo {};
`,
			`async function f() {
  await goog.importHandler_('sG5M4c');
  const {Foo} = a.b.c;
}
`))
}

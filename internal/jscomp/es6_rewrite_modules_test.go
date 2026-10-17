// Port of com/google/javascript/jscomp/Es6RewriteModulesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/Es6RewriteModulesTest.java)

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

// Unit tests for Es6RewriteModules
//
// See also Es6RewriteModulesBeforeTypeCheckingTest. This file tests only
// behavior when run after type checking.
//
// TODO(b/144593112): Remove the other test class when the pass is permanently
// moved after type checking.
type Es6RewriteModulesTest struct {
	CompilerTestCaseBase
	moduleRoots     []string
	chunkOutputType CompilerOptions_ChunkOutputType
}

// Java: Es6RewriteModulesTest.other
var es6RewriteModulesTest_other = SourceFile_FromCode(
	"other.js",
	`export default 0;
export let name, x, a, b, c;
export {x as class};
export class Parent {}
`)

// Java: Es6RewriteModulesTest.otherExpected
var es6RewriteModulesTest_otherExpected = SourceFile_FromCode(
	"other.js",
	`var $jscompDefaultExport$$module$other = 0;
let name$$module$other, x$$module$other, a$$module$other, b$$module$other,
  c$$module$other;
class Parent$$module$other {}
/** @const */ var module$other = {};
/** @const */ module$other.Parent = Parent$$module$other;
/** @const */ module$other.a = a$$module$other;
/** @const */ module$other.b = b$$module$other;
/** @const */ module$other.c = c$$module$other;
/** @const */ module$other.class = x$$module$other;
/** @const */ module$other.default = $jscompDefaultExport$$module$other;
/** @const */ module$other.name = name$$module$other;
/** @const */ module$other.x = x$$module$other;
`)

func NewEs6RewriteModulesTest() *Es6RewriteModulesTest {
	c := &Es6RewriteModulesTest{}
	c.moduleRoots = nil
	c.chunkOutputType = CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestEs6RewriteModulesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteModulesTest() })
}

// Java: Es6RewriteModulesTest.setUp()
func (c *Es6RewriteModulesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// ECMASCRIPT5 to trigger module processing after parsing.
	c.EnableCreateModuleMap()
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
}

// Java: Es6RewriteModulesTest.getOptions()
func (c *Es6RewriteModulesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	// ECMASCRIPT5 to Trigger module processing after parsing.
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_ERROR)
	options.SetPrettyPrint(true)

	if c.moduleRoots != nil {
		options.SetModuleRoots(c.moduleRoots)
	}

	return options
}

// Java: Es6RewriteModulesTest.getProcessor(Compiler)
func (c *Es6RewriteModulesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		rai :=
			NewSemanticReverseAbstractInterpreter(compiler.GetTypeRegistry())
		globalTypedScope :=
			jutil.CheckNotNull(
				NewTypeCheck(compiler, rai, compiler.GetTypeRegistry()).
					ProcessForTesting(externs, root))
		compiler.SetTypeCheckingHasRun(true)
		NewEs6RewriteModulesWithChunkOutputType(
			compiler,
			compiler.GetModuleMetadataMap(),
			compiler.GetModuleMap(),
			/* preprocessorSymbolTable= */ nil,
			globalTypedScope,
			c.chunkOutputType).
			Process(externs, root)

		if c.chunkOutputType == CompilerOptions_ChunkOutputType_ES_MODULES {
			NewConvertChunksToESModules(compiler).Process(externs, root)
		}
	})
}

// Java: Es6RewriteModulesTest.testModules(String, String)
func (c *Es6RewriteModulesTest) TestModules(input string, expected string) {
	c.Test(
		c.Srcs(es6RewriteModulesTest_other, SourceFile_FromCode("testcode", input)),
		c.Expected(es6RewriteModulesTest_otherExpected, SourceFile_FromCode("testcode", expected)))
}

// Java: Es6RewriteModulesTest.testImport()
func (c *Es6RewriteModulesTest) TestImport() {
	c.TestModules(
		`import name from './other.js';
use(name);
`,
		"use($jscompDefaultExport$$module$other); /** @const */ var module$testcode = {};")

	c.TestModules("import {a as name} from './other.js';", "/** @const */ var module$testcode = {};")

	c.TestModules(
		`import x, {a as foo, b as bar} from './other.js';
use(x);
`,
		"use($jscompDefaultExport$$module$other); /** @const */ var module$testcode = {};")

	c.TestModules(
		`import {default as name} from './other.js';
use(name);
`,
		"use($jscompDefaultExport$$module$other); /** @const */ var module$testcode = {};")

	c.TestModules(
		`import {class as name} from './other.js';
use(name);
`,
		"use(x$$module$other); /** @const */ var module$testcode = {};")
}

// Java: Es6RewriteModulesTest.testImport_missing()
func (c *Es6RewriteModulesTest) TestImport_missing() {
	ModulesTestUtils_TestModulesError(
		c, "import name from './does_not_exist';\n use(name);", ModuleLoader_LOAD_WARNING)

	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)

	// These are different as a side effect of the way that the fake bindings are made. The first
	// makes a binding for a fake variable in the fake module. The second creates a fake binding
	// for the fake module. When "dne.name" is resolved, the module does not have key "name", so
	// it chooses to rewrite to "module$does_not_exist.name", thinking that this could've been a
	// reference to an export that doesn't exist.
	c.TestModules(
		`import {name} from './does_not_exist';
use(name);
`,
		`use(name$$module$does_not_exist);
/** @const */ var module$testcode = {};
`)

	c.TestModules(
		`import * as dne from './does_not_exist';
use(dne.name);
`,
		`use(module$does_not_exist.name);
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesTest.testImportStar()
func (c *Es6RewriteModulesTest) TestImportStar() {
	c.TestModules(
		`import * as name from './other.js';
use(name.a);
`,
		"use(a$$module$other); /** @const */ var module$testcode = {};")
}

// Java: Es6RewriteModulesTest.testTypeNodeRewriting()
func (c *Es6RewriteModulesTest) TestTypeNodeRewriting() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"other.js",
				`export default 0;
export let name = 'George';
export let a = class {};
`),
			SourceFile_FromCode(
				"testcode",
				`import * as name from './other.js';
/** @type {name.a} */ var x;
`)),
		c.Expected(
			SourceFile_FromCode(
				"other.js",
				`var $jscompDefaultExport$$module$other = 0;
let name$$module$other = 'George';
let a$$module$other = class {};
/** @const */ var module$other = {};
/** @const */ module$other.a = a$$module$other;
/** @const */ module$other.default = $jscompDefaultExport$$module$other;
/** @const */ module$other.name = name$$module$other;
`),
			SourceFile_FromCode(
				"testcode",
				`/** @type {a$$module$other} */ var x$$module$testcode;
/** @const */ var module$testcode = {};
`)))
}

// Java: Es6RewriteModulesTest.testExport()
func (c *Es6RewriteModulesTest) TestExport() {
	c.TestModules(
		"export var a = 1, b = 2;",
		`var a$$module$testcode = 1, b$$module$testcode = 2;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.a = a$$module$testcode;
/** @const */ module$testcode.b = b$$module$testcode;
`)

	c.TestModules(
		"export var a;\nexport var b;",
		`var a$$module$testcode; var b$$module$testcode;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.a = a$$module$testcode;
/** @const */ module$testcode.b = b$$module$testcode;
`)

	c.TestModules(
		"export function f() {};",
		`function f$$module$testcode() {}
/** @const */ var module$testcode = {};
/** @const */ module$testcode.f = f$$module$testcode;
`)

	c.TestModules(
		"export function f() {};\nfunction g() { f(); }",
		`function f$$module$testcode() {}
function g$$module$testcode() { f$$module$testcode(); }
/** @const */ var module$testcode = {};
/** @const */ module$testcode.f = f$$module$testcode;
`)

	c.TestModules(
		`export function MyClass() {};
MyClass.prototype.foo = function() {};
`,
		`function MyClass$$module$testcode() {}
MyClass$$module$testcode.prototype.foo = function() {};
/** @const */ var module$testcode = {};
/** @const */ module$testcode.MyClass = MyClass$$module$testcode;
`)

	c.TestModules(
		"var f = 1;\nvar b = 2;\nexport {f as foo, b as bar};",
		`var f$$module$testcode = 1;
var b$$module$testcode = 2;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.bar = b$$module$testcode;
/** @const */ module$testcode.foo = f$$module$testcode;
`)

	c.TestModules(
		"var f = 1;\nexport {f as default};",
		`var f$$module$testcode = 1;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.default = f$$module$testcode;
`)

	c.TestModules(
		"var f = 1;\nexport {f as class};",
		`var f$$module$testcode = 1;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.class = f$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testModulesInExterns()
func (c *Es6RewriteModulesTest) TestModulesInExterns() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				"externsMod.js",
				`/** @fileoverview @externs */
export let /** !number */ externalName;
`)),
		c.Error(TranspilationUtil_CANNOT_CONVERT_YET))
}

// Java: Es6RewriteModulesTest.testModulesInTypeSummary()
func (c *Es6RewriteModulesTest) TestModulesInTypeSummary() {
	c.AllowExternsChanges()
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"mod1.js",
				`/** @fileoverview @typeSummary */
export let /** !number */ externalName;
`),
			SourceFile_FromCode(
				"mod2.js",
				`import {externalName as localName} from './mod1.js'
alert(localName);
`)),
		c.Expected(
			SourceFile_FromCode(
				"mod2.js",
				`alert(externalName$$module$mod1);
/** @const */ var module$mod2 = {};
`)))
}

// Java: Es6RewriteModulesTest.testMutableExport()
func (c *Es6RewriteModulesTest) TestMutableExport() {
	c.TestModules(
		`export var a = 1, b = 2;
function f() {
  a++;
  b++
}
`,
		`var a$$module$testcode = 1, b$$module$testcode = 2;
function f$$module$testcode() {
  a$$module$testcode++;
  b$$module$testcode++
}
/** @const */ var module$testcode = {
  get a() { return a$$module$testcode; },
  get b() { return b$$module$testcode; },
};
`)

	c.TestModules(
		`var a = 1, b = 2; export {a as A, b as B};
const f = () => {
  a++;
  b++
};
`,
		`var a$$module$testcode = 1, b$$module$testcode = 2;
const f$$module$testcode = () => {
  a$$module$testcode++;
  b$$module$testcode++
};
/** @const */ var module$testcode = {
  get A() { return a$$module$testcode; },
  get B() { return b$$module$testcode; },
};
`)

	c.TestModules(
		`export function f() {};
function g() {
  f = function() {};
}
`,
		`function f$$module$testcode() {}
function g$$module$testcode() {
  f$$module$testcode = function() {};
}
/** @const */ var module$testcode = {
  get f() { return f$$module$testcode; },
};
`)

	c.TestModules(
		`export default function f() {};
function g() {
  f = function() {};
}
`,
		`function f$$module$testcode() {}
function g$$module$testcode() {
  f$$module$testcode = function() {};
}
/** @const */ var module$testcode = {
  get default() { return f$$module$testcode; },
};
`)

	c.TestModules(
		`export class C {};
function g() {
  C = class {};
}
`,
		`class C$$module$testcode {}
function g$$module$testcode() {
  C$$module$testcode = class {};
}
/** @const */ var module$testcode = {
  get C() { return C$$module$testcode; },
};
`)

	c.TestModules(
		`export default class C {};
function g() {
  C = class {};
}
`,
		`class C$$module$testcode {}
function g$$module$testcode() {
  C$$module$testcode = class {};
}
/** @const */ var module$testcode = {
  get default() { return C$$module$testcode; },
};
`)

	c.TestModules(
		`export var IN, OF;
function f() {
  for (IN in {});
  for (OF of []);
}
`,
		`var IN$$module$testcode, OF$$module$testcode;
function f$$module$testcode() {
  for (IN$$module$testcode in {});
  for (OF$$module$testcode of []);
}
/** @const */ var module$testcode = {
  get IN() { return IN$$module$testcode; },
  get OF() { return OF$$module$testcode; },
};
`)

	c.TestModules(
		`export var ARRAY, OBJ, UNCHANGED;
function f() {
// TODO(lharker): this typechecking suppression seems like it shouldn't be necessary,
// but mutable exports are not really well-supported in the compiler so not a high
// priority to fix.
  /** @suppress {checkTypes} */
  ({OBJ} = {OBJ: {}});
  [ARRAY] = [];
  var x = {UNCHANGED: 0};
}
`,
		`var ARRAY$$module$testcode, OBJ$$module$testcode, UNCHANGED$$module$testcode;
function f$$module$testcode() {
  /** @suppress {checkTypes} */
  ({OBJ:OBJ$$module$testcode} = {OBJ: {}});
  [ARRAY$$module$testcode] = [];
  var x = {UNCHANGED: 0};
}
/** @const */ var module$testcode = {
  get ARRAY() { return ARRAY$$module$testcode; },
  get OBJ() { return OBJ$$module$testcode; },
};
/** @const */ module$testcode.UNCHANGED = UNCHANGED$$module$testcode
`)
}

// Java: Es6RewriteModulesTest.testConstClassExportIsConstant()
func (c *Es6RewriteModulesTest) TestConstClassExportIsConstant() {
	c.TestModules(
		"export const Class = class {}",
		`const Class$$module$testcode = class {}
/** @const */ var module$testcode = {};
/** @const */ module$testcode.Class = Class$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testTopLevelMutationIsNotMutable()
func (c *Es6RewriteModulesTest) TestTopLevelMutationIsNotMutable() {
	c.TestModules(
		`export var a = 1, b = 2;
a++;
b++
`,
		`var a$$module$testcode = 1, b$$module$testcode = 2;
a$$module$testcode++;
b$$module$testcode++
/** @const */ var module$testcode = {};
/** @const */ module$testcode.a = a$$module$testcode;
/** @const */ module$testcode.b = b$$module$testcode;
`)

	c.TestModules(
		`var a = 1, b = 2; export {a as A, b as B};
if (change) {
  a++;
  b++
}
`,
		`var a$$module$testcode = 1, b$$module$testcode = 2;
if (change) {
  a$$module$testcode++;
  b$$module$testcode++
}
/** @const */ var module$testcode = {};
/** @const */ module$testcode.A = a$$module$testcode;
/** @const */ module$testcode.B = b$$module$testcode;
`)

	c.TestModules(
		`export function f() {};
if (change) {
  f = function() {};
}
`,
		`function f$$module$testcode() {}
if (change) {
  f$$module$testcode = function() {};
}
/** @const */ var module$testcode = {};
/** @const */ module$testcode.f = f$$module$testcode;
`)

	c.TestModules(
		`export default function f() {};
try { f = function() {}; } catch (e) { f = function() {}; }
`,
		`function f$$module$testcode() {}
try { f$$module$testcode = function() {}; }
catch (e) { f$$module$testcode = function() {}; }
/** @const */ var module$testcode = {};
/** @const */ module$testcode.default = f$$module$testcode;
`)

	c.TestModules(
		`export class C {};
if (change) {
  C = class {};
}
`,
		`class C$$module$testcode {}
if (change) {
  C$$module$testcode = class {};
}
/** @const */ var module$testcode = {};
/** @const */ module$testcode.C = C$$module$testcode;
`)

	c.TestModules(
		`export default class C {};
{
  C = class {};
}
`,
		`class C$$module$testcode {}
{
  C$$module$testcode = class {};
}
/** @const */ var module$testcode = {};
/** @const */ module$testcode.default = C$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testExportWithJsDoc()
func (c *Es6RewriteModulesTest) TestExportWithJsDoc() {
	c.TestModules(
		"/** @constructor */\nexport function F() { return '';}",
		`/** @constructor */
function F$$module$testcode() { return ''; }
/** @const */ var module$testcode = {};
/** @const */ module$testcode.F = F$$module$testcode;
`)

	c.TestModules(
		"/** @return {string} */\nexport function f() { return '';}",
		`/** @return {string} */
function f$$module$testcode() { return ''; }
/** @const */ var module$testcode = {};
/** @const */ module$testcode.f = f$$module$testcode;
`)

	c.TestModules(
		"/** @return {string} */\nexport var f = function() { return '';}",
		`/** @return {string} */
var f$$module$testcode = function() { return ''; }
/** @const */ var module$testcode = {};
/** @const */ module$testcode.f = f$$module$testcode;
`)

	c.TestModules(
		"/** @type {number} */\nexport var x = 3",
		`/** @type {number} */
var x$$module$testcode = 3;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.x = x$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testImportAndExport()
func (c *Es6RewriteModulesTest) TestImportAndExport() {
	c.TestModules(
		`import {name as n} from './other.js';
use(n);
export {n as name};
`,
		`use(name$$module$other);
/** @const */ var module$testcode = {};
/** @const */ module$testcode.name = name$$module$other;
`)
}

// Java: Es6RewriteModulesTest.testExportFrom()
func (c *Es6RewriteModulesTest) TestExportFrom() {
	c.Test(
		c.Srcs(
			es6RewriteModulesTest_other,
			SourceFile_FromCode(
				"testcode",
				`export {name} from './other.js';
export {default} from './other.js';
export {class} from './other.js';
`)),
		c.Expected(
			es6RewriteModulesTest_otherExpected,
			SourceFile_FromCode(
				"testcode",
				`/** @const */ var module$testcode = {};
/** @const */ module$testcode.class = x$$module$other;
/** @const */ module$testcode.default = $jscompDefaultExport$$module$other;
/** @const */ module$testcode.name = name$$module$other;
`)))

	c.Test(
		c.Srcs(es6RewriteModulesTest_other, SourceFile_FromCode("testcode", "export {a, b as c, x} from './other.js';")),
		c.Expected(
			es6RewriteModulesTest_otherExpected,
			SourceFile_FromCode(
				"testcode",
				`/** @const */ var module$testcode = {}
/** @const */ module$testcode.a = a$$module$other;
/** @const */ module$testcode.c = b$$module$other;
/** @const */ module$testcode.x = x$$module$other;
`)))

	c.Test(
		c.Srcs(es6RewriteModulesTest_other, SourceFile_FromCode("testcode", "export {a as b, b as a} from './other.js';")),
		c.Expected(
			es6RewriteModulesTest_otherExpected,
			SourceFile_FromCode(
				"testcode",
				`/** @const */ var module$testcode = {}
/** @const */ module$testcode.a = b$$module$other;
/** @const */ module$testcode.b = a$$module$other;
`)))

	c.Test(
		c.Srcs(
			es6RewriteModulesTest_other,
			SourceFile_FromCode(
				"testcode",
				`export {default as a} from './other.js';
export {a as a2, default as b} from './other.js';
export {class as switch} from './other.js';
`)),
		c.Expected(
			es6RewriteModulesTest_otherExpected,
			SourceFile_FromCode(
				"testcode",
				`/** @const */ var module$testcode = {}
/** @const */ module$testcode.a = $jscompDefaultExport$$module$other;
/** @const */ module$testcode.a2 = a$$module$other;
/** @const */ module$testcode.b = $jscompDefaultExport$$module$other;
/** @const */ module$testcode.switch = x$$module$other;
`)))
}

// Java: Es6RewriteModulesTest.testExportDefault()
func (c *Es6RewriteModulesTest) TestExportDefault() {
	c.TestModules(
		"export default 'someString';",
		`var $jscompDefaultExport$$module$testcode = 'someString';
/** @const */ var module$testcode = {};
/** @const */ module$testcode.default = $jscompDefaultExport$$module$testcode;
`)

	c.TestModules(
		"var x = 5;\nexport default x;",
		`var x$$module$testcode = 5;
var $jscompDefaultExport$$module$testcode = x$$module$testcode;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.default = $jscompDefaultExport$$module$testcode;
`)

	c.TestModules(
		"export default function f(){};\n var x = f();",
		`function f$$module$testcode() {}
var x$$module$testcode = f$$module$testcode();
/** @const */ var module$testcode = {};
/** @const */ module$testcode.default = f$$module$testcode;
`)

	c.TestModules(
		"export default class Foo {};\n var x = new Foo;",
		`class Foo$$module$testcode {}
var x$$module$testcode = new Foo$$module$testcode;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.default = Foo$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testExportDefault_anonymous()
func (c *Es6RewriteModulesTest) TestExportDefault_anonymous() {
	c.TestModules(
		"export default class {};",
		`var $jscompDefaultExport$$module$testcode = class {};
/** @const */ var module$testcode = {};
/** @const */ module$testcode.default = $jscompDefaultExport$$module$testcode;
`)

	c.TestModules(
		"export default function() {}",
		`var $jscompDefaultExport$$module$testcode = function() {}
/** @const */ var module$testcode = {};
/** @const */ module$testcode.default = $jscompDefaultExport$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testExportDestructureDeclaration()
func (c *Es6RewriteModulesTest) TestExportDestructureDeclaration() {
	c.TestModules(
		"export let {a, c:b} = obj;",
		`let {a:a$$module$testcode, c:b$$module$testcode} = obj;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.a = a$$module$testcode;
/** @const */ module$testcode.b = b$$module$testcode;
`)

	c.TestModules(
		"export let [a, b] = obj;",
		`let [a$$module$testcode, b$$module$testcode] = obj;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.a = a$$module$testcode;
/** @const */ module$testcode.b = b$$module$testcode;
`)

	c.TestModules(
		"export let {a, b:[c,d]} = obj;",
		`let {a:a$$module$testcode, b:[c$$module$testcode, d$$module$testcode]} = obj;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.a = a$$module$testcode;
/** @const */ module$testcode.c = c$$module$testcode;
/** @const */ module$testcode.d = d$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testExtendImportedClass()
func (c *Es6RewriteModulesTest) TestExtendImportedClass() {
	c.TestModules(
		`import {Parent} from './other.js';
class Child extends Parent {
  /** @param {Parent} parent */
  useParent(parent) {}
}
`,
		`class Child$$module$testcode extends Parent$$module$other {
  /** @param {Parent$$module$other} parent */
  useParent(parent) {}
}
/** @const */ var module$testcode = {};
`)

	c.TestModules(
		`import {Parent} from './other.js';
export class Child extends Parent {
  /** @param {Parent} parent */
  useParent(parent) {}
}
`,
		`class Child$$module$testcode extends Parent$$module$other {
  /** @param {Parent$$module$other} parent */
  useParent(parent) {}
}
/** @const */ var module$testcode = {};
/** @const */ module$testcode.Child = Child$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testFixTypeNode()
func (c *Es6RewriteModulesTest) TestFixTypeNode() {
	c.TestModules(
		`export class Child {
  /** @param {Child} child */
  useChild(child) {}
}
`,
		`class Child$$module$testcode {
  /** @param {Child$$module$testcode} child */
  useChild(child) {}
}
/** @const */ var module$testcode = {};
/** @const */ module$testcode.Child = Child$$module$testcode;
`)

	c.TestModules(
		`export class Child {
  /** @param {Child.Foo.Bar.Baz} baz */
  useBaz(baz) {}
}

Child.Foo = class {};

Child.Foo.Bar = class {};

Child.Foo.Bar.Baz = class {};
`,
		`class Child$$module$testcode {
  /** @param {Child$$module$testcode.Foo.Bar.Baz} baz */
  useBaz(baz) {}
}
Child$$module$testcode.Foo=class{};
Child$$module$testcode.Foo.Bar=class{};
Child$$module$testcode.Foo.Bar.Baz=class{};
/** @const */ var module$testcode = {};
/** @const */ module$testcode.Child = Child$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testRenameTypedef()
func (c *Es6RewriteModulesTest) TestRenameTypedef() {
	c.TestModules(
		`import './other.js';
/** @typedef {string|!Object} */
export var UnionType;
`,
		`/** @typedef {string|!Object} */
var UnionType$$module$testcode;
/** @const */ var module$testcode = {};
/** @typedef {UnionType$$module$testcode} */
module$testcode.UnionType;
`)
}

// Java: Es6RewriteModulesTest.testNoInnerChange()
func (c *Es6RewriteModulesTest) TestNoInnerChange() {
	c.TestModules(
		`var Foo = (function () {
    /**  @param bar */
    function Foo(bar) {}
    return Foo;
}());
export { Foo };
`,
		`var Foo$$module$testcode = function() {
    /**  @param bar */
    function Foo(bar) {}
    return Foo;
}();
/** @const */ var module$testcode = {};
/** @const */ module$testcode.Foo = Foo$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testRenameImportedReference()
func (c *Es6RewriteModulesTest) TestRenameImportedReference() {
	c.TestModules(
		`import {a} from './other.js';
import {b as bar} from './other.js';
a();
function g() {
  a();
  bar++;
  function h() {
    var a = 3;
    { let a = 4; }
  }
}
`,
		`a$$module$other();
function g$$module$testcode() {
  a$$module$other();
  b$$module$other++;
  function h() {
    var a = 3;
    { let a = 4; }
  }
}
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesTest.testObjectDestructuringAndObjLitShorthand()
func (c *Es6RewriteModulesTest) TestObjectDestructuringAndObjLitShorthand() {
	c.TestModules(
		`import {c} from './other.js';
const foo = 1;
const {a, b} = c({foo});
use(a, b);
`,
		`const foo$$module$testcode = 1;
const {
  a: a$$module$testcode,
  b: b$$module$testcode,
} = c$$module$other({foo: foo$$module$testcode});
use(a$$module$testcode, b$$module$testcode);
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesTest.testObjectDestructuringAndObjLitShorthandWithDefaultValue()
func (c *Es6RewriteModulesTest) TestObjectDestructuringAndObjLitShorthandWithDefaultValue() {
	c.TestModules(
		`import {c} from './other.js';
const foo = 1;
const {a = 'A', b = 'B'} = c({foo});
use(a, b);
`,
		`const foo$$module$testcode = 1;
const {
  a: a$$module$testcode = 'A',
  b: b$$module$testcode = 'B',
} = c$$module$other({foo: foo$$module$testcode});
use(a$$module$testcode, b$$module$testcode);
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesTest.testImportWithoutReferences()
func (c *Es6RewriteModulesTest) TestImportWithoutReferences() {
	c.TestModules(
		"import './other.js';", //
		"/** @const */ var module$testcode = {};")
}

// Java: Es6RewriteModulesTest.testUselessUseStrict()
func (c *Es6RewriteModulesTest) TestUselessUseStrict() {
	ModulesTestUtils_TestModulesError(
		c,
		`'use strict';
export default undefined;
`,
		ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE)
}

// Java: Es6RewriteModulesTest.testUseStrict_noWarning()
func (c *Es6RewriteModulesTest) TestUseStrict_noWarning() {
	c.TestSame(
		`'use strict';
var x;
`)
}

// Java: Es6RewriteModulesTest.testAbsoluteImportsWithModuleRoots()
func (c *Es6RewriteModulesTest) TestAbsoluteImportsWithModuleRoots() {
	c.moduleRoots = []string{"/base"}
	c.Test(
		c.Srcs(
			SourceFile_FromCode(Compiler_JoinPathParts("base", "mod", "name.js"), ""),
			SourceFile_FromCode(
				Compiler_JoinPathParts("base", "test", "sub.js"),
				"import * as foo from '/mod/name.js';")),
		c.Expected(
			SourceFile_FromCode(Compiler_JoinPathParts("base", "mod", "name.js"), ""),
			SourceFile_FromCode(
				Compiler_JoinPathParts("base", "test", "sub.js"),
				"/** @const */ var module$test$sub = {};")))
}

// Java: Es6RewriteModulesTest.testUseImportInEs6ObjectLiteralShorthand()
func (c *Es6RewriteModulesTest) TestUseImportInEs6ObjectLiteralShorthand() {
	c.TestModules(
		`import {b} from './other.js';
var bar = {a: 1, b};
`,
		`var bar$$module$testcode={a: 1, b: b$$module$other};
/** @const */ var module$testcode = {};
`)

	c.TestModules(
		`import {a as foo} from './other.js';
var bar = {a: 1, foo};
`,
		`var bar$$module$testcode={a: 1, foo: a$$module$other};
/** @const */ var module$testcode = {};
`)

	c.TestModules(
		`import f from './other.js';
var bar = {a: 1, f};
`,
		`var bar$$module$testcode={a: 1, f: $jscompDefaultExport$$module$other};
/** @const */ var module$testcode = {};
`)

	c.TestModules(
		`import * as f from './other.js';
// TODO(lharker): this typechecking suppression seems like it shouldn't be necessary
/** @suppress {checkTypes} */
var bar = {a: 1, f};
`,
		`/** @suppress {checkTypes} */
var bar$$module$testcode={a: 1, f: module$other};
/** @const */ var module$testcode = {};
`)
}

// Java: Es6RewriteModulesTest.testImportAliasInTypeNode()
func (c *Es6RewriteModulesTest) TestImportAliasInTypeNode() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode("a.js", "export class A {}"),
			SourceFile_FromCode(
				"b.js",
				`import {A as B} from './a.js';
const /** !B */ b = new B();
`)),
		c.Expected(
			SourceFile_FromCode(
				"a.js",
				`class A$$module$a {}
/** @const */ var module$a = {};
/** @const */ module$a.A = A$$module$a;
`),
			SourceFile_FromCode(
				"b.js",
				`const /** !A$$module$a*/ b$$module$b = new A$$module$a();
/** @const */ var module$b = {};
`)))
}

// Java: Es6RewriteModulesTest.testExportStar()
func (c *Es6RewriteModulesTest) TestExportStar() {
	c.TestModules(
		"export * from './other.js';",
		`/** @const */ var module$testcode = {};
/** @const */ module$testcode.Parent = Parent$$module$other;
/** @const */ module$testcode.a = a$$module$other;
/** @const */ module$testcode.b = b$$module$other;
/** @const */ module$testcode.c = c$$module$other;
/** @const */ module$testcode.class = x$$module$other;
// no default
/** @const */ module$testcode.name = name$$module$other;
/** @const */ module$testcode.x = x$$module$other;
`)
}

// Java: Es6RewriteModulesTest.testExportStarWithLocalExport()
func (c *Es6RewriteModulesTest) TestExportStarWithLocalExport() {
	c.TestModules(
		`export * from './other.js';
export let baz, zed;
`,
		`let baz$$module$testcode, zed$$module$testcode;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.Parent = Parent$$module$other;
/** @const */ module$testcode.a = a$$module$other;
/** @const */ module$testcode.b = b$$module$other;
/** @const */ module$testcode.baz = baz$$module$testcode;
/** @const */ module$testcode.c = c$$module$other;
/** @const */ module$testcode.class = x$$module$other;
/** @const */ module$testcode.name = name$$module$other;
/** @const */ module$testcode.x = x$$module$other;
/** @const */ module$testcode.zed = zed$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testExportStarWithLocalExportOverride()
func (c *Es6RewriteModulesTest) TestExportStarWithLocalExportOverride() {
	c.TestModules(
		`export * from './other.js';
export let a, zed;
`,
		`let a$$module$testcode, zed$$module$testcode;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.Parent = Parent$$module$other;
/** @const */ module$testcode.a = a$$module$testcode;
/** @const */ module$testcode.b = b$$module$other;
/** @const */ module$testcode.c = c$$module$other;
/** @const */ module$testcode.class = x$$module$other;
/** @const */ module$testcode.name = name$$module$other;
/** @const */ module$testcode.x = x$$module$other;
/** @const */ module$testcode.zed = zed$$module$testcode;
`)
}

// Java: Es6RewriteModulesTest.testTransitiveImport()
func (c *Es6RewriteModulesTest) TestTransitiveImport() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode("a.js", "export class A {}"),
			SourceFile_FromCode("b.js", "export {A} from './a.js';"),
			SourceFile_FromCode(
				"c.js",
				`import {A} from './b.js';
let /** !A */ a = new A();
`)),
		c.Expected(
			SourceFile_FromCode(
				"a.js",
				`class A$$module$a {}
/** @const */ var module$a = {};
/** @const */ module$a.A = A$$module$a;
`),
			SourceFile_FromCode(
				"b.js",
				`/** @const */ var module$b = {};
/** @const */ module$b.A = A$$module$a;
`),
			SourceFile_FromCode(
				"c.js",
				`let /** !A$$module$a*/ a$$module$c = new A$$module$a();
/** @const */ var module$c = {};
`)))
	c.Test(
		c.Srcs(
			SourceFile_FromCode("a.js", "export class A {}"),
			SourceFile_FromCode("b.js", "export {A as B} from './a.js';"),
			SourceFile_FromCode(
				"c.js",
				`import {B as C} from './b.js';
let /** !C */ a = new C();
`)),
		c.Expected(
			SourceFile_FromCode(
				"a.js",
				`class A$$module$a {}
/** @const */ var module$a = {};
/** @const */ module$a.A = A$$module$a;
`),
			SourceFile_FromCode(
				"b.js",
				`/** @const */ var module$b = {};
/** @const */ module$b.B = A$$module$a;
`),
			SourceFile_FromCode(
				"c.js",
				`let /** !A$$module$a*/ a$$module$c = new A$$module$a();
/** @const */ var module$c = {};
`)))
}

// Java: Es6RewriteModulesTest.testMutableTransitiveImport()
func (c *Es6RewriteModulesTest) TestMutableTransitiveImport() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode("a.js", "export let A = class {}; () => (A = class {});"),
			SourceFile_FromCode("b.js", "export {A} from './a.js';"),
			SourceFile_FromCode(
				"c.js",
				`import {A} from './b.js';
let /** !A */ a = new A();
`)),
		c.Expected(
			SourceFile_FromCode(
				"a.js",
				`let A$$module$a = class {};
()=>A$$module$a = class {};
/** @const */ var module$a = {
  get A() { return A$$module$a; },
};
`),
			SourceFile_FromCode(
				"b.js",
				`/** @const */ var module$b = {
  get A() { return A$$module$a; },
};
`),
			SourceFile_FromCode(
				"c.js",
				`let /** !A$$module$a*/ a$$module$c = new A$$module$a();
/** @const */ var module$c = {};
`)))
	c.Test(
		c.Srcs(
			SourceFile_FromCode("a.js", "export let A = class {}; () => (A = class {});"),
			SourceFile_FromCode("b.js", "export {A as B} from './a.js';"),
			SourceFile_FromCode(
				"c.js",
				`import {B as C} from './b.js';
let /** !C */ a = new C();
`)),
		c.Expected(
			SourceFile_FromCode(
				"a.js",
				`let A$$module$a = class {};
()=>A$$module$a = class {};
/** @const */ var module$a = {
  get A() { return A$$module$a; },
};
`),
			SourceFile_FromCode(
				"b.js",
				`/** @const */ var module$b = {
  get B() { return A$$module$a; },
};
`),
			SourceFile_FromCode(
				"c.js",
				`let /** !A$$module$a*/ a$$module$c = new A$$module$a();
/** @const */ var module$c = {};
`)))
}

// Java: Es6RewriteModulesTest.testRewriteGetPropsWhileModuleReference()
func (c *Es6RewriteModulesTest) TestRewriteGetPropsWhileModuleReference() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode("a.js", "export class A {}"),
			SourceFile_FromCode(
				"b.js",
				`import * as a from './a.js';
export {a};
`),
			SourceFile_FromCode(
				"c.js",
				`import * as b from './b.js';
let /** !b.a.A */ a = new b.a.A();
`)),
		c.Expected(
			SourceFile_FromCode(
				"a.js",
				`class A$$module$a {}
/** @const */ var module$a = {};
/** @const */ module$a.A = A$$module$a;
`),
			SourceFile_FromCode(
				"b.js",
				`/** @const */ var module$b = {};
/** @const */ module$b.a = module$a;
`),
			SourceFile_FromCode(
				"c.js",
				`let /** !A$$module$a*/ a$$module$c = new A$$module$a();
/** @const */ var module$c = {};
`)))
}

// Java: Es6RewriteModulesTest.testRewritePropsWhenNotModuleReference()
func (c *Es6RewriteModulesTest) TestRewritePropsWhenNotModuleReference() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode("other.js", "export const name = {}, a = { Type: class {} };"),
			SourceFile_FromCode(
				"testcode",
				`import * as name from './other.js';
let /** !name.a.Type */ t = new name.a.Type();
`)),
		c.Expected(
			SourceFile_FromCode(
				"other.js",
				`const name$$module$other = {}, a$$module$other = { Type: class {} };
/** @const */ var module$other = {};
/** @const */ module$other.a = a$$module$other;
/** @const */ module$other.name = name$$module$other;
`),
			SourceFile_FromCode(
				"testcode",
				`let /** !a$$module$other.Type */ t$$module$testcode =
   new a$$module$other.Type();
/** @const */ var module$testcode = {};
`)))
}

// Java: Es6RewriteModulesTest.testExportsNotImplicitlyLocallyDeclared()
func (c *Es6RewriteModulesTest) TestExportsNotImplicitlyLocallyDeclared() {
	c.Test(
		c.Externs("var exports;"),
		c.Srcs("typeof exports; export {};"),
		// Regression test; compiler used to rewrite `exports` to `exports$$module$testcode`.
		c.Expected("typeof exports; /** @const */ var module$testcode = {};"))
}

// Java: Es6RewriteModulesTest.testImportMeta()
func (c *Es6RewriteModulesTest) TestImportMeta() {
	c.TestError("import.meta", TranspilationUtil_CANNOT_CONVERT)
}

// Java: Es6RewriteModulesTest.testImportMetaESModuleOutput()
func (c *Es6RewriteModulesTest) TestImportMetaESModuleOutput() {
	c.chunkOutputType = CompilerOptions_ChunkOutputType_ES_MODULES
	c.Test(
		`const url = import.meta.url;
export {url};
`,
		`const url$$module$testcode = import.meta.url;
/** @const */ var module$testcode = {};
/** @const */ module$testcode.url = url$$module$testcode;
export {};
`)
}

// Java: Es6RewriteModulesTest.testGoogColonImport_fromLegacyNamespace()
func (c *Es6RewriteModulesTest) TestGoogColonImport_fromLegacyNamespace() {
	closureBase := testExternsBuilder_GetClosureExternsAsSource()
	googAsserts :=
		`goog.module('goog.asserts');
goog.module.declareLegacyNamespace();
function assert(condition) { if (!condition) { throw new Error('Assertion failed')}}
exports = {assert};
`
	c.Test(
		c.Srcs(
			closureBase,
			googAsserts,
			`import {assert} from 'goog:goog.asserts';
assert(1 < 2);
`),
		c.Expected(
			closureBase,
			googAsserts,
			`(0, goog.asserts.assert)(1 < 2);
/** @const */
var module$testcode2 = {};
`))
}

// Java: Es6RewriteModulesTest.testImportDot()
func (c *Es6RewriteModulesTest) TestImportDot() {
	c.SetModuleResolutionMode(ModuleLoader_ResolutionMode_NODE)
	c.Test(
		c.Srcs(
			SourceFile_FromCode("mod/index.js", "export const foo = 1;"),
			SourceFile_FromCode("mod/sub.js", "import {foo} from '.'; use(foo);")),
		c.Expected(
			SourceFile_FromCode(
				"mod/index.js",
				`const foo$$module$mod$index = 1;
/** @const */ var module$mod$index = {};
/** @const */ module$mod$index.foo = foo$$module$mod$index;
`),
			SourceFile_FromCode(
				"mod/sub.js",
				`use(foo$$module$mod$index);
/** @const */ var module$mod$sub = {};
`)))
}

// Java: Es6RewriteModulesTest.testImportDotDot()
func (c *Es6RewriteModulesTest) TestImportDotDot() {
	c.SetModuleResolutionMode(ModuleLoader_ResolutionMode_NODE)
	c.Test(
		c.Srcs(
			SourceFile_FromCode("mod/index.js", "export const foo = 1;"),
			SourceFile_FromCode("mod/sub/sub.js", "import {foo} from '..'; use(foo);")),
		c.Expected(
			SourceFile_FromCode(
				"mod/index.js",
				`const foo$$module$mod$index = 1;
/** @const */ var module$mod$index = {};
/** @const */ module$mod$index.foo = foo$$module$mod$index;
`),
			SourceFile_FromCode(
				"mod/sub/sub.js",
				`use(foo$$module$mod$index);
/** @const */ var module$mod$sub$sub = {};
`)))
}

// Java: Es6RewriteModulesTest.testEs6RenameTypeReferences()
func (c *Es6RewriteModulesTest) TestEs6RenameTypeReferences() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode("mod/other.js", "export class Foo {}"),
			SourceFile_FromCode("mod/app.js", "import {Foo} from './other.js'; /** @type {Foo} */ var x;")),
		c.Expected(
			SourceFile_FromCode(
				"mod/other.js",
				`class Foo$$module$mod$other {}
/** @const */ var module$mod$other = {};
/** @const */ module$mod$other.Foo = Foo$$module$mod$other;
`),
			SourceFile_FromCode(
				"mod/app.js",
				`/** @type {Foo$$module$mod$other} */
var x$$module$mod$app;
/** @const */ var module$mod$app = {};
`)))
}

// Java: Es6RewriteModulesTest.testModuleRenaming()
func (c *Es6RewriteModulesTest) TestModuleRenaming() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode("mod/other.js", "export class Foo {}"),
			SourceFile_FromCode("mod/app.js", "import {Foo} from './other.js'; use(Foo);")),
		c.Expected(
			SourceFile_FromCode(
				"mod/other.js",
				`class Foo$$module$mod$other {}
/** @const */ var module$mod$other = {};
/** @const */ module$mod$other.Foo = Foo$$module$mod$other;
`),
			SourceFile_FromCode(
				"mod/app.js",
				`use(Foo$$module$mod$other);
/** @const */ var module$mod$app = {};
`)))
}

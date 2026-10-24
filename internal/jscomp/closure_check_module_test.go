// Port of com/google/javascript/jscomp/ClosureCheckModuleTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ClosureCheckModuleTest.java)

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

type ClosureCheckModuleTest struct {
	CompilerTestCaseBase
}

func NewClosureCheckModuleTest() *ClosureCheckModuleTest {
	c := &ClosureCheckModuleTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestClosureCheckModuleTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewClosureCheckModuleTest() })
}

// Java: ClosureCheckModuleTest.getProcessor(Compiler)
func (c *ClosureCheckModuleTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		gatherModuleMetadata :=
			NewGatherModuleMetadata(
				compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER)
		gatherModuleMetadata.Process(externs, root)
		NewClosureCheckModule(compiler, compiler.GetModuleMetadataMap()).Process(externs, root)
	})
}

// Java: ClosureCheckModuleTest.getOptions()
func (c *ClosureCheckModuleTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_ERROR)
	options.SetWarningLevel(DiagnosticGroups_STRICT_MODULE_CHECKS, CheckLevel_ERROR)
	return options
}

// Java: ClosureCheckModuleTest.testGoogModuleThisInSubclass()
func (c *ClosureCheckModuleTest) TestGoogModuleThisInSubclass() {
	c.TestSame(
		`goog.module('xyz');
class Foo {
  x = 5;
}
class Bar extends Foo {
   z = this.x;
}
`)

	c.TestSame(
		`goog.module('xyz');
class Foo {
  [x] = 5;
}
class Bar extends Foo {
   z = this[x];
}
`)

	c.TestSame(
		`goog.module('xyz');
class Foo {
  [x] = 5;
}
class Bar extends Foo {
   [z] = this[x];
}
`)
}

// Java: ClosureCheckModuleTest.testGoogModuleThisOnFields()
func (c *ClosureCheckModuleTest) TestGoogModuleThisOnFields() {
	// Ensure `this.property` is allowed as an assignment on public fields.
	c.TestSame(
		`goog.module('xyz');
class Foo {
  x = 5;
  y = this.x
}
`)

	c.TestSame(
		`goog.module('xyz');
class Foo {
  [x] = 5;
  [y] = this[x]
}
`)

	c.TestError(
		`goog.module('xyz');
class Foo {
  [x] = 5;
  [this.x] = 6
}
`,
		ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS)
}

// Java: ClosureCheckModuleTest.testStaticGoogModuleThisOnStaticFields()
func (c *ClosureCheckModuleTest) TestStaticGoogModuleThisOnStaticFields() {
	// Allow this to be allowed when a field is already declared
	c.TestSame(
		`goog.module('math')
class Foo {
  static x = 5;
  static y = this.x
}
`)

	c.TestSame(
		`goog.module('math')
class Foo {
  static [x] = 5;
  static [y] = this[x]
}
`)
}

// Java: ClosureCheckModuleTest.testGoogModuleReferencesThis()
func (c *ClosureCheckModuleTest) TestGoogModuleReferencesThis() {
	c.TestError("goog.module('xyz');\nfoo.call(this, 1, 2, 3);", ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS)

	c.TestError(
		`goog.module('xyz');

var x = goog.require('other.x');

if (x) {
  alert(this);
}
`,
		ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS)

	c.TestSame(
		`goog.module('xyz');

class Foo {
  constructor() {
    this.x = 5;
  }
}

exports = Foo;
`)
}

// Java: ClosureCheckModuleTest.testGoogModuleUsesThrow()
func (c *ClosureCheckModuleTest) TestGoogModuleUsesThrow() {
	c.TestError("goog.module('xyz');\nthrow 4;", ClosureCheckModule_GOOG_MODULE_USES_THROW)

	c.TestError(
		`goog.module('xyz');

var x = goog.require('other.x');

if (x) {
  throw 5;
}
`,
		ClosureCheckModule_GOOG_MODULE_USES_THROW)
}

// Java: ClosureCheckModuleTest.testGoogModuleGetAtTopLevel()
func (c *ClosureCheckModuleTest) TestGoogModuleGetAtTopLevel() {
	c.TestError("goog.module('xyz');\ngoog.module.get('abc');", ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET)

	c.TestError(
		`goog.module('xyz');

var x = goog.require('other.x');

if (x) {
  var y = goog.module.get('abc');
}
`,
		ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET)

	c.TestSame(
		`goog.module('xyz');

var x = goog.require('other.x');

function f() {
  var y = goog.module.get('abc');
}
`)
}

// Java: ClosureCheckModuleTest.testMultipleGoogModules()
func (c *ClosureCheckModuleTest) TestMultipleGoogModules() {
	c.TestError(
		`goog.module('xyz');
goog.module('abc');

var x = goog.require('other.x');
`,
		ClosureCheckModule_MULTIPLE_MODULES_IN_FILE)
}

// Java: ClosureCheckModuleTest.testMisplacedGoogModuleCall()
func (c *ClosureCheckModuleTest) TestMisplacedGoogModuleCall() {
	c.TestError(
		`var x;
goog.module('xyz');
`,
		ClosureCheckModule_GOOG_MODULE_MISPLACED)

	c.TestError(
		`;
goog.module('xyz');
`,
		ClosureCheckModule_GOOG_MODULE_MISPLACED)

	c.TestError(
		`function fn() {
  goog.module('xyz');
 }
`,
		ClosureCheckModule_GOOG_MODULE_MISPLACED)

	c.TestError("const mod =  goog.module('xyz');", ClosureCheckModule_GOOG_MODULE_MISPLACED)

	c.TestError(
		`goog.loadModule(function(exports) {
  var x;
  goog.module('xyz');
  return exports;
})
`,
		ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE)
}

// Java: ClosureCheckModuleTest.testBundledGoogModules()
func (c *ClosureCheckModuleTest) TestBundledGoogModules() {
	c.TestError(
		`goog.loadModule(function(exports){
  'use strict';
  goog.module('xyz');
  foo.call(this, 1, 2, 3);
  return exports;
});
`,
		ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS)

	c.TestError(
		`goog.loadModule(function(exports){
  'use strict';
  goog.module('foo.example.ClassName');
  /** @constructor @export */ function ClassName() {}
  exports = ClassName;
  return exports;
});
`,
		ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE)

	c.TestSame(
		`goog.loadModule(function(exports){
  'use strict';
  goog.module('Xyz');
  exports = class {}
  return exports;
});
goog.loadModule(function(exports){
  goog.module('abc');
  var Foo = goog.require('Xyz');
  var x = new Foo;
  return exports;
});
`)

	c.TestError(
		`goog.loadModule(function(exports){
  'use strict';
  goog.module('xyz');
  goog.module('abc');
  var x = goog.require('other.x');
  return exports;
});
`,
		ClosureCheckModule_MULTIPLE_MODULES_IN_FILE)
}

// Java: ClosureCheckModuleTest.testBundledGoogModulesAndProvides()
func (c *ClosureCheckModuleTest) TestBundledGoogModulesAndProvides() {
	c.Test(
		c.Srcs(
			`goog.provide('first.provide');
first.provide = 0;

goog.provide('second.provide');
second.provide = 0;

goog.loadModule(function(exports){
  'use strict';
  goog.module('Xyz');
  const first = goog.require('first.provide');
  exports = class {}
  return exports;
});

goog.loadModule(function(exports){
  goog.module('abc');
  const Foo = goog.require('Xyz');
  const second = goog.require('second.provide');
  var x = new Foo;
  return exports;
});
`),
		c.Error(ClosureCheckModule_USE_OF_GOOG_PROVIDE),
		c.Error(ClosureCheckModule_USE_OF_GOOG_PROVIDE))
}

// Java: ClosureCheckModuleTest.testBundledGoogModulesAndProvidesWithGoog()
func (c *ClosureCheckModuleTest) TestBundledGoogModulesAndProvidesWithGoog() {
	c.Test(
		c.Srcs(
			`/** @provideGoog */
var goog = {};

goog.provide('first.provide');
first.provide = 0;

goog.provide('second.provide');
second.provide = 0;

goog.loadModule(function(exports){
  'use strict';
  goog.module('Xyz');
  const first = goog.require('first.provide');
  exports = class {}
  return exports;
});

goog.loadModule(function(exports){
  goog.module('abc');
  const Foo = goog.require('Xyz');
  const second = goog.require('second.provide');
  var x = new Foo;
  return exports;
});
`),
		c.Error(ClosureCheckModule_USE_OF_GOOG_PROVIDE),
		c.Error(ClosureCheckModule_USE_OF_GOOG_PROVIDE))
}

// Java: ClosureCheckModuleTest.testGoogModuleReferencesGlobalName()
func (c *ClosureCheckModuleTest) TestGoogModuleReferencesGlobalName() {
	c.TestError("goog.module('x.y.z');\nx.y.z = function() {};", ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME)
}

// Java: ClosureCheckModuleTest.testIllegalAtExport()
func (c *ClosureCheckModuleTest) TestIllegalAtExport() {
	c.TestError(
		`goog.module('foo.example.ClassName');

/** @constructor @export */ function ClassName() {}

exports = ClassName;
`,
		ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE)

	c.TestError(
		`goog.module('foo.example.ClassName');

/** @export */ class ClassName {}

exports = ClassName;
`,
		ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE)

	c.TestError(
		`goog.module('foo.example.ClassName');

/** @constructor */ function ClassName() {}

/** @export */
exports = ClassName;
`,
		ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE)

	c.TestError(
		`goog.module('foo.example.ns');

/** @constructor */ function ClassName() {}

/** @export */
exports.ClassName = ClassName;
`,
		ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE)

	c.TestError(
		`goog.module('foo.example.ClassName');
goog.module.declareLegacyNamespace();

/** @constructor @export */ function ClassName() {}

exports = ClassName;
`,
		ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE)

	// TODO(b/123020550): warn for this, as the actual name exported will be
	//   module$contents$foo_example_ClassName.Builder
	c.TestSame(
		`goog.module('foo.example.ClassName');

class ClassName {}

/** @export */
ClassName.Builder = function() {};

exports = ClassName;
`)
}

// Java: ClosureCheckModuleTest.testLegalAtExport()
func (c *ClosureCheckModuleTest) TestLegalAtExport() {
	c.TestSame(
		`goog.module('foo.example.ClassName');

class ClassName {
  constructor() {
    /** @export */
    this.prop;
    /** @export */
    this.anotherProp = false;
  }
}

exports = ClassName;
`)

	c.TestSame(
		`goog.module('foo.example.ClassName');

var ClassName = class {
  constructor() {
    /** @export */
    this.prop;
    /** @export */
    this.anotherProp = false;
  }
}

exports = ClassName;
`)

	c.TestSame(
		`goog.module('foo.example.ClassName');

/** @constructor */
function ClassName() {
  /** @export */
  this.prop;
  /** @export */
  this.anotherProp = false;
}

exports = ClassName;
`)

	c.TestSame(
		`goog.module('foo.example.ClassName');

/** @constructor */
var ClassName = function() {
  /** @export */
  this.prop;
  /** @export */
  this.anotherProp = false;
};

exports = ClassName;
`)

	c.TestSame(
		`goog.module('foo.example.ClassName');
goog.module.declareLegacyNamespace();

/** @constructor */ function ClassName() {}

/** @export */
exports = ClassName;
`)

	c.TestSame(
		`goog.module('foo.example.ns');
goog.module.declareLegacyNamespace();

/** @constructor */ function ClassName() {}

/** @export */
exports.ClassName = ClassName;
`)

	c.TestSame(
		`goog.module('foo.example.ClassName');

/** @constructor */ var exports = function() {}

/** @export */
exports.prototype.fly = function() {};
`)
}

// Java: ClosureCheckModuleTest.testIllegalDeclareLegacyNamespace()
func (c *ClosureCheckModuleTest) TestIllegalDeclareLegacyNamespace() {
	c.Test(
		c.Srcs(
			`goog.provide('a.provided.namespace');
goog.module.declareLegacyNamespace();
`),
		c.Error(ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE),
		c.Error(ClosureCheckModule_USE_OF_GOOG_PROVIDE))

	c.TestError(
		`goog.module('xyz');
const ns = goog.module.declareLegacyNamespace();
`,
		ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL)

	c.TestError(
		`goog.module('xyz');
if (cond) {
  goog.module.declareLegacyNamespace();
}
`,
		ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE)

	c.TestError(
		`goog.module('xyz');
var foo = 0;
goog.module.declareLegacyNamespace();
`,
		ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE)

	c.TestError(
		`goog.loadModule(function(exports) {
  goog.module('xyz');
  var x;
  goog.module.declareLegacyNamespace();
  return exports
});
`,
		ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE)

	c.TestError(
		`goog.module('my.mod');
goog.module.declareLegacyNamespace('some comment');
`,
		ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT)
}

// Java: ClosureCheckModuleTest.testIllegalGoogRequires()
func (c *ClosureCheckModuleTest) TestIllegalGoogRequires() {
	c.TestError(
		`goog.module('xyz');

var foo = goog.require('other.x').foo;
`,
		ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL)

	c.TestError(
		`goog.module('xyz');

var x = goog.require('other.x').foo.toString();
`,
		ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL)

	c.TestError(
		`goog.module('xyz');

var moduleNames = [goog.require('other.x').name];
`,
		ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL)

	c.TestError(
		`goog.module('xyz');

exports = [goog.require('other.x').name];
`,
		ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL)

	c.TestError(
		`goog.module('xyz');

var a = goog.require('foo.a'), b = goog.require('foo.b');
`,
		ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION)

	c.TestError(
		`goog.module('xyz');

var [foo, bar] = goog.require('other.x');
`,
		ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE)

	c.TestError(
		`goog.module('xyz');

var {foo, bar = 'str'} = goog.require('other.x');
`,
		ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE)

	c.TestError(
		`goog.module('xyz');

var {foo, bar: {name}} = goog.require('other.x');
`,
		ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE)

	c.TestError(
		`goog.module('xyz');

var {[foo]: bar} = goog.require('other.x');
`,
		ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE)

	c.TestError(
		`goog.module('xyz');

var {...rest} = goog.require('other.x');
`,
		ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE)

	c.TestError(
		`goog.module('xyz');

var [foo, bar] = goog.requireType('other.x');
`,
		ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE)

	c.TestError(
		`goog.module('xyz');

var {foo, bar = 'str'} = goog.requireType('other.x');
`,
		ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE)

	c.TestError(
		`goog.module('xyz');

var {foo, bar: {name}} = goog.requireType('other.x');
`,
		ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE)

	c.TestError(
		`goog.module('xyz');

var foo = goog.require('abc.foo');
var foo = goog.require('def.foo');
`,
		ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE)

	c.TestError(
		`goog.module('xyz');

var {foo, bar} = goog.require('abc');
var foo = goog.require('def.foo');
`,
		ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE)
}

// Java: ClosureCheckModuleTest.testIllegalShortImportReferencedByLongName()
func (c *ClosureCheckModuleTest) TestIllegalShortImportReferencedByLongName() {
	c.TestError(
		`goog.module('x.y.z');

var A = goog.require('foo.A');

exports = function() { return new foo.A; };
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)
}

// Java: ClosureCheckModuleTest.testCorrectShortImportSuggestedInJsDoc()
func (c *ClosureCheckModuleTest) TestCorrectShortImportSuggestedInJsDoc() {
	c.TestError(
		c.Srcs(
			`goog.module('x.y.z');

var A = goog.require('foo.A');
var B = goog.require('foo.A.B');

/** @type {foo.A.B} */ var b;
`),
		c.Error(ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME).
			WithMessageContaining("Please use the short name 'B' instead."))
}

// Java: ClosureCheckModuleTest.testCorrectShortImportSuggestedInJsCode()
func (c *ClosureCheckModuleTest) TestCorrectShortImportSuggestedInJsCode() {
	c.TestError(
		c.Srcs(
			`goog.module('x.y.z');

var A = goog.require('foo.A');
var B = goog.require('foo.A.B');

ref(foo.A.B);
`),
		c.Error(ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME).
			WithMessageContaining("Please use the short name 'B' instead."))
}

// Java: ClosureCheckModuleTest.testIllegalShortImportReferencedByLongNameInJsDoc()
func (c *ClosureCheckModuleTest) TestIllegalShortImportReferencedByLongNameInJsDoc() {
	c.TestError(
		`goog.module('x.y.z');

var A = goog.require('foo.A');

/** @constructor @implements {foo.A} */ function B() {}
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)

	c.TestError(
		`goog.module('x.y.z');

var A = goog.requireType('foo.A');

/** @constructor @implements {foo.A} */ function B() {}
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)

	c.TestError(
		`goog.module('x.y.z');

var A = goog.require('foo.A');

/** @type {foo.A} */ var a;
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)

	c.TestError(
		`goog.module('x.y.z');

var A = goog.requireType('foo.A');

/** @type {foo.A} */ var a;
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)

	c.TestSame(
		`goog.module('x.y.z');

var A = goog.require('foo.A');

/** @type {A} */ var a;
`)

	c.TestSame(
		`goog.module('x.y.z');

var A = goog.requireType('foo.A');

/** @type {A} */ var a;
`)

	c.TestSame(
		`goog.module('x.y.z');

var Foo = goog.require('Foo');

/** @type {Foo} */ var a;
`)

	c.TestSame(
		`goog.module('x.y.z');

var Foo = goog.requireType('Foo');

/** @type {Foo} */ var a;
`)

	c.TestError(
		`goog.module('x.y.z');

var ns = goog.require('some.namespace');

/** @type {some.namespace.Foo} */ var foo;
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)

	c.TestError(
		`goog.module('x.y.z');

var ns = goog.requireType('some.namespace');

/** @type {some.namespace.Foo} */ var foo;
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)

	c.TestError(
		`goog.module('x.y.z');

var ns = goog.require('some.namespace');

/** @type {Array<some.namespace.Foo>} */ var foos;
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)

	c.TestError(
		`goog.module('x.y.z');

var ns = goog.requireType('some.namespace');

/** @type {Array<some.namespace.Foo>} */ var foos;
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)
}

// Java: ClosureCheckModuleTest.testIllegalShortImportDestructuring()
func (c *ClosureCheckModuleTest) TestIllegalShortImportDestructuring() {
	c.TestError(
		`goog.module('x.y.z');

var {doThing} = goog.require('foo.utils');

exports = function() { return foo.utils.doThing(''); };
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)

	c.TestError(
		`goog.module('x.y.z');

var {doThing: fooDoThing} = goog.require('foo.utils');

exports = function() { return foo.utils.doThing(''); };
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)
}

// Java: ClosureCheckModuleTest.testIllegalRequireNoAlias()
func (c *ClosureCheckModuleTest) TestIllegalRequireNoAlias() {
	c.TestError(
		`goog.module('x.y.z');

goog.require('foo.utils');

exports = function() { return foo.utils.doThing(''); };
`,
		ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME)
}

// Java: ClosureCheckModuleTest.testIllegalRequireTypeNoAlias()
func (c *ClosureCheckModuleTest) TestIllegalRequireTypeNoAlias() {
	c.TestError(
		`goog.module('x.y.z');

goog.requireType('foo.Bar');

/** @type {foo.Bar} */ var foo;
`,
		ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME)
}

// Java: ClosureCheckModuleTest.testLegalForwardDeclareNoAlias()
func (c *ClosureCheckModuleTest) TestLegalForwardDeclareNoAlias() {
	c.TestNoWarning(
		`goog.module('x.y.z');

goog.forwardDeclare('foo.utils');

exports = function() { return foo.utils.doThing(''); };
`)
}

// Java: ClosureCheckModuleTest.testSingleNameImportNoAlias1()
func (c *ClosureCheckModuleTest) TestSingleNameImportNoAlias1() {
	c.TestError(
		`goog.module('x.y.z');

goog.require('foo');

exports = function() { return foo.doThing(''); };
`,
		ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME)
}

// Java: ClosureCheckModuleTest.testSingleNameImportWithRenamingAlias()
func (c *ClosureCheckModuleTest) TestSingleNameImportWithRenamingAlias() {
	c.TestError(
		`goog.module('x.y.z');

var bar = goog.require('foo');

exports = function() { return foo.doThing(''); };
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)
}

// Java: ClosureCheckModuleTest.testSingleNameImportWithRenamingAlias_forwardDeclare()
func (c *ClosureCheckModuleTest) TestSingleNameImportWithRenamingAlias_forwardDeclare() {
	c.TestError(
		`goog.module('x.y.z');

var bar = goog.forwardDeclare('foo');

exports = function() { return foo.doThing(''); };
`,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME)
}

// Java: ClosureCheckModuleTest.testSingleNameImportCrossAlias()
func (c *ClosureCheckModuleTest) TestSingleNameImportCrossAlias() {
	c.TestSame(
		`goog.module('x.y.z');

var bar = goog.require('foo');
var foo = goog.require('bar');

exports = function() { return foo.doThing(''); };
`)
}

// Java: ClosureCheckModuleTest.testLegalSingleNameImport()
func (c *ClosureCheckModuleTest) TestLegalSingleNameImport() {
	c.TestSame(
		`goog.module('x.y.z');

var foo = goog.require('foo');

exports = function() { return foo.doThing(''); };
`)
}

// Java: ClosureCheckModuleTest.testSingleNameImportShadowed()
func (c *ClosureCheckModuleTest) TestSingleNameImportShadowed() {
	c.TestSame(
		`goog.module('x.y.z');

// for side-effects only
goog.require('foo');

exports = function(foo) { return foo.doThing(''); };
`)
}

// Java: ClosureCheckModuleTest.testIllegalLetShortRequire()
func (c *ClosureCheckModuleTest) TestIllegalLetShortRequire() {
	c.TestSame(
		`goog.module('xyz');

let a = goog.forwardDeclare('foo.a');
`)

	c.TestError(
		`goog.module('xyz');

let a = goog.require('foo.a');
`,
		ClosureCheckModule_LET_GOOG_REQUIRE)
}

// Java: ClosureCheckModuleTest.testIllegalDestructuringForwardDeclare()
func (c *ClosureCheckModuleTest) TestIllegalDestructuringForwardDeclare() {
	c.TestError(
		"goog.module('m'); var {x} = goog.forwardDeclare('a.b.c');",
		ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE)
}

// Java: ClosureCheckModuleTest.testLegalGoogRequires()
func (c *ClosureCheckModuleTest) TestLegalGoogRequires() {
	c.TestSame(
		`goog.module('xyz');

var {assert} = goog.require('goog.asserts');
`)

	c.TestSame(
		`goog.module('xyz');

const {assert} = goog.require('goog.asserts');
`)

	c.TestSame(
		`goog.module('xyz');

const {assert, fail} = goog.require('goog.asserts');
`)
}

// Java: ClosureCheckModuleTest.testShorthandNameConvention()
func (c *ClosureCheckModuleTest) TestShorthandNameConvention() {
	c.TestSame(
		`goog.module('xyz');

const googAsserts = goog.require('goog.asserts');
`)

	c.TestError(
		`goog.module('xyz');

const GoogAsserts = goog.require('goog.asserts');
`,
		ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION,
		"The capitalization of short name GoogAsserts is incorrect; it should be googAsserts.")

	c.TestSame(
		`goog.module('xyz');

const Event = goog.require('goog.events.Event');
`)

	c.TestError(
		`goog.module('xyz');

const event = goog.require('goog.events.Event');
`,
		ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION,
		"The capitalization of short name event is incorrect; it should be Event.")

	c.TestSame(
		`goog.module('xyz');

const Event = goog.require('google3.closure.events.event');
`)
}

// Java: ClosureCheckModuleTest.testNonModuleLevelExports()
func (c *ClosureCheckModuleTest) TestNonModuleLevelExports() {
	c.TestError(
		`goog.module('xyz');

if (window.exportMe) { exports = 5; }
`,
		ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE)

	c.TestError(
		`goog.module('xyz');

window.exportMe && (exports = 5);
`,
		ClosureCheckModule_EXPORT_NOT_A_STATEMENT)

	c.TestError(
		`goog.module('xyz');

exports.f = () => { exports.me = 5; }
`,
		ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE)

	c.TestSame(
		`goog.module('xyz');

exports = class {};
exports.staticMethod = () => { exports.me = 5; }
`)
}

// Java: ClosureCheckModuleTest.testRepeatedExports()
func (c *ClosureCheckModuleTest) TestRepeatedExports() {
	c.TestError(
		`goog.module('xyz');

exports = 5;
exports = 'str';
`,
		ClosureCheckModule_EXPORT_REPEATED_ERROR)

	c.TestError(
		`goog.module('xyz');

exports.y = 5;
exports.y = 'str';
`,
		ClosureCheckModule_EXPORT_REPEATED_ERROR)

	c.TestSame(
		`goog.module('xyz');

exports = class {};
exports.y = 'str';
exports.y = decorate(exports.y);
`)
}

// Java: ClosureCheckModuleTest.testDontCrashOnTrailingDot()
func (c *ClosureCheckModuleTest) TestDontCrashOnTrailingDot() {
	c.TestSame(
		`goog.module('foo');

var a = goog.require('abc.');
`)
}

// Java: ClosureCheckModuleTest.testAllowGoogProvideDeclaration()
func (c *ClosureCheckModuleTest) TestAllowGoogProvideDeclaration() {
	c.TestSame("const goog = {}; goog.provide = function(ns) {};")
}

// Java: ClosureCheckModuleTest.testWarnOnGoogProvideCall()
func (c *ClosureCheckModuleTest) TestWarnOnGoogProvideCall() {
	c.TestError("goog.provide('foo.bar');", ClosureCheckModule_USE_OF_GOOG_PROVIDE)
}

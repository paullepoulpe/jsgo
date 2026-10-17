// Port of com/google/javascript/jscomp/CheckMissingRequiresTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/CheckMissingRequiresTest.java)

/*
 * Copyright 2020 The Closure Compiler Authors.
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

// CheckMissingRequiresTest: Unit tests for CheckMissingRequires.
type CheckMissingRequiresTest struct {
	CompilerTestCaseBase
}

// Java: CheckMissingRequiresTest()
func NewCheckMissingRequiresTest() *CheckMissingRequiresTest {
	c := &CheckMissingRequiresTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckMissingRequiresTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckMissingRequiresTest() })
}

// Java: CheckMissingRequiresTest.getProcessor(Compiler)
func (c *CheckMissingRequiresTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		gatherModuleMetadata :=
			NewGatherModuleMetadata(
				compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER)
		gatherModuleMetadata.Process(externs, root)
		NewCheckMissingRequires(compiler, compiler.GetModuleMetadataMap()).Process(externs, root)
	})
}

// Java: CheckMissingRequiresTest.getOptions()
func (c *CheckMissingRequiresTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MISSING_REQUIRE, CheckLevel_WARNING)
	return options
}

// Java: CheckMissingRequiresTest.testNoWarning_existingRequire_withAlias()
func (c *CheckMissingRequiresTest) TestNoWarning_existingRequire_withAlias() {
	c.CheckNoWarning(
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
const Bar = goog.require('foo.Bar');
let x = new Bar();
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_existingRequire_withDestructure()
func (c *CheckMissingRequiresTest) TestNoWarning_existingRequire_withDestructure() {
	c.CheckNoWarning(
		`goog.module('foo.bar');
/** @constructor */
exports.Baz = function() {}
`,
		`goog.module('test');
const {Bar} = goog.require('foo.Bar');
let x = new Bar();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_inProvide()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_inProvide() {
	c.CheckRequireInProvidesFileWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
		`goog.provide('test');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_inScript()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_inScript() {
	c.CheckRequireInProvidesFileWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
		`goog.provide('test');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testMissingRequire_inEsModuleExport()
func (c *CheckMissingRequiresTest) TestMissingRequire_inEsModuleExport() {
	c.CheckRequireWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
		`let x = new foo.Bar();
export default 42;
`)
}

// Java: CheckMissingRequiresTest.testMissingRequire_inEsModuleImport()
func (c *CheckMissingRequiresTest) TestMissingRequire_inEsModuleImport() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.CheckRequireWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
		`import * as quux from './quux.js';
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequire_externs()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequire_externs() {
	c.CheckNoWarning(
		`/** @externs */
var foo = {}
/** @constructor */
foo.Bar = function() {};
`,
		`goog.provide('test');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequire_unknown()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequire_unknown() {
	c.CheckNoWarning(
		`goog.provide('test');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequire_sameProvide()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequire_sameProvide() {
	c.CheckNoWarning(
		`goog.provide('foo.Bar');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequire_sameNestedProvide()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequire_sameNestedProvide() {
	c.CheckNoWarning(
		"goog.provide('foo');",
		`goog.provide('foo.Bar');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_strongRef_withRequireType()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_strongRef_withRequireType() {
	c.CheckRequireInProvidesFileWarning(
		"a.b.C",
		`goog.provide('a.b.C');

/** @constructor */
a.b.C = function() { };
`,
		`goog.requireType('a.b.C');

new a.b.C;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequire_withParentRequire_fromSameFile()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequire_withParentRequire_fromSameFile() {
	c.CheckNoWarning(
		`goog.provide('a.b');
goog.provide('a.b.C');
`,
		`goog.require('a.b');

new a.b.C;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_withSiblingRequire_fromSameFile()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_withSiblingRequire_fromSameFile() {
	c.CheckRequireInProvidesFileWarning(
		"a.b.D",
		`goog.provide('a.b');
goog.provide('a.b.C');
goog.provide('a.b.D');
`,
		`goog.require('a.b.C');

new a.b.D;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_withParentRequire_fromDifferentFile()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_withParentRequire_fromDifferentFile() {
	c.CheckRequireInProvidesFileWarning(
		"a.b.C",
		"goog.provide('a.b');",
		`goog.require('a.b');
goog.provide('a.b.C');
`,
		`goog.require('a.b');

new a.b.C;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequire_sameModule()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequire_sameModule() {
	c.CheckNoWarning(
		`goog.module('foo.Bar');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequire_sameModule_nestedProvide()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequire_sameModule_nestedProvide() {
	c.CheckNoWarning(
		"goog.provide('foo.bar');",
		`goog.module('foo.bar.Baz');
/** @constructor */
exports.Baz = function() {}
let x = new foo.bar.Baz();
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequire_sameModuleWithLegacyNamespace_nestedProvide()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequire_sameModuleWithLegacyNamespace_nestedProvide() {
	c.CheckNoWarning(
		"goog.provide('foo.bar');",
		`goog.module('foo.bar.Baz');
goog.module.declareLegacyNamespace();
/** @constructor */
exports.Baz = function() {}
let x = new foo.bar.Baz();
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_existingRequireType_withAlias()
func (c *CheckMissingRequiresTest) TestNoWarning_existingRequireType_withAlias() {
	c.CheckNoWarning(
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
const Bar = goog.requireType('foo.Bar');
/** @type {!Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_existingRequireType_withDestructure()
func (c *CheckMissingRequiresTest) TestNoWarning_existingRequireType_withDestructure() {
	c.CheckNoWarning(
		`goog.module('foo.bar');
/** @constructor */
exports.Baz = function() {}
`,
		`goog.module('test');
const {Bar} = goog.requireType('foo.Bar');
/** @type {!Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_inProvide()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_inProvide() {
	c.CheckRequireTypeInProvidesFileWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
		`goog.provide('test');
/** @type {!foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_inScript()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_inScript() {
	c.CheckRequireTypeInProvidesFileWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
		`/** @type {!foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testMissingRequireType_inEsModuleExport()
func (c *CheckMissingRequiresTest) TestMissingRequireType_inEsModuleExport() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
		`/** @type {!foo.Bar} */
let x;
export default 42;
`)
}

// Java: CheckMissingRequiresTest.testMissingRequireType_inEsModuleImport()
func (c *CheckMissingRequiresTest) TestMissingRequireType_inEsModuleImport() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
		`import * as quux from './quux.js';
/** @type {!foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_forProvide()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_forProvide() {
	c.CheckRequireWarning(
		"foo.Bar",
		`goog.provide('foo.Bar');
/** @constructor */
foo.Bar = function() {}
`,
		`goog.module('test');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_forProvide_usingProperty()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_forProvide_usingProperty() {
	c.CheckRequireWarning(
		"foo.bar",
		`goog.provide('foo.bar');
foo.bar = {};
/** @constructor */
foo.bar.Baz = function() {}
`,
		`goog.module('test');
let x = new foo.bar.Baz();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_forNestedProvide()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_forNestedProvide() {
	c.CheckRequireWarning(
		"foo.Bar",
		"goog.provide('foo');",
		`goog.provide('foo.Bar');
/** @constructor */
foo.Bar = function() {}
`,
		`goog.module('test');
goog.require('foo');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_forNestedProvide_usingProperty()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_forNestedProvide_usingProperty() {
	c.CheckRequireWarning(
		"foo.bar",
		"goog.provide('foo');",
		`goog.provide('foo.bar');
foo.bar = {};
/** @constructor */
foo.bar.Baz = function() {}
`,
		`goog.module('test');
goog.require('foo');
let x = new foo.bar.Baz();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_forModule()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_forModule() {
	c.CheckRequireWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_forModule_usingProperty()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_forModule_usingProperty() {
	c.CheckRequireWarning(
		"foo.bar",
		`goog.module('foo.bar');
goog.module.declareLegacyNamespace();
/** @constructor */
function Baz() {}
exports = { Baz }
`,
		`goog.module('test');
let x = new foo.bar.Baz();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_forModuleWithLegacyNamespace()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_forModuleWithLegacyNamespace() {
	c.CheckRequireWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_forLateModule()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_forLateModule() {
	c.CheckRequireWarning(
		"foo.Bar",
		`goog.module('test');
let x = new foo.Bar();
`,
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`)
}

// Java: CheckMissingRequiresTest.testWarning_existingRequire_standalone()
func (c *CheckMissingRequiresTest) TestWarning_existingRequire_standalone() {
	c.CheckRequireWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
goog.require('foo.Bar');
let x = new foo.Bar();
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_forProvide()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_forProvide() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.provide('foo.Bar');
/** @constructor */
foo.Bar = function() {}
`,
		`goog.module('test');
/** @type {!foo.Bar} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_forProvide_usingProperty()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_forProvide_usingProperty() {
	c.CheckRequireTypeWarning(
		"foo.bar",
		`goog.provide('foo.bar');
foo.bar = {};
/** @constructor */
foo.bar.Baz = function() {}
`,
		`goog.module('test');
/** @type {!foo.bar.Baz} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_forNestedProvide()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_forNestedProvide() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		"goog.provide('foo');",
		`goog.provide('foo.Bar');
/** @constructor */
foo.Bar = function() {}
`,
		`goog.module('test');
goog.requireType('foo');
/** @type {!foo.Bar} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_forNestedProvide_usingProperty()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_forNestedProvide_usingProperty() {
	c.CheckRequireTypeWarning(
		"foo.bar",
		"goog.provide('foo');",
		`goog.provide('foo.bar');
foo.bar = {};
/** @constructor */
foo.bar.Baz = function() {}
`,
		`goog.module('test');
goog.require('foo');
/** @type {!foo.bar.Baz} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_forModule()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_forModule() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @type {!foo.Bar} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_forModule_usingProperty()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_forModule_usingProperty() {
	c.CheckRequireTypeWarning(
		"foo.bar",
		`goog.module('foo.bar');
goog.module.declareLegacyNamespace();
/** @constructor */
function Baz() {}
exports = { Baz }
`,
		`goog.module('test');
/** @type {!foo.bar.Baz} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_forModuleWithLegacyNamespace()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_forModuleWithLegacyNamespace() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @type {!foo.Bar} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_forLateModule()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_forLateModule() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.provide('foo.Bar');
/** @type {!foo.Bar} */ let x;
`,
		`goog.module('test');
/** @type {!foo.Bar} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_existingRequireType_standalone()
func (c *CheckMissingRequiresTest) TestWarning_existingRequireType_standalone() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
goog.requireType('foo.Bar');
/** @type {!foo.Bar} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequireType_externs()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequireType_externs() {
	c.CheckNoWarning(
		`/** @externs */
var foo = {}
/** @constructor */
foo.Bar = function() {};
`,
		`goog.provide('test');
/** @type {!foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequireType_unknown()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequireType_unknown() {
	c.CheckNoWarning(
		`goog.provide('test');
/** @type {!foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequireType_sameProvide()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequireType_sameProvide() {
	c.CheckNoWarning(
		`goog.provide('foo.Bar');
/** @type {!foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequireType_sameNestedProvide()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequireType_sameNestedProvide() {
	c.CheckNoWarning(
		"goog.provide('foo');",
		`goog.provide('foo.Bar');
/** @type {!foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequireType_sameModule()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequireType_sameModule() {
	c.CheckNoWarning(
		`goog.module('foo.Bar');
/** @type {!foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequireType_sameModule_nestedProvide()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequireType_sameModule_nestedProvide() {
	c.CheckNoWarning(
		"goog.provide('foo.bar');",
		`goog.module('foo.bar.Baz');
/** @constructor */
exports.Baz = function() {}
/** @type {!foo.bar.Baz} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_nestedProvide()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_nestedProvide() {
	c.CheckIncorrectNamespaceAliasRequireWarning(
		"foo.bar.Baz",
		"goog.provide('foo.bar');",
		"goog.provide('foo.bar.Baz');",
		`goog.module('another');
const {Baz} = goog.require('foo.bar');
function ref(a) {};
ref(Baz);
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_nestedProvide()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_nestedProvide() {
	c.CheckIncorrectNamespaceAliasRequireTypeWarning(
		"foo.bar.Baz",
		"goog.provide('foo.bar');",
		"goog.provide('foo.bar.Baz');foo.bar.Baz = class {};",
		`goog.module('another');
const {Baz} = goog.requireType('foo.bar');
let /** !Baz */ x;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequireType_nestedProvide()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequireType_nestedProvide() {
	c.CheckNoWarning(
		"goog.provide('foo.bar');",
		"goog.provide('foo.bar.Baz');foo.bar.Baz = class {};",
		`goog.module('another');
const {Baz} = goog.requireType('foo.bar.Baz');
let /** !Baz */ x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_nestedModule()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_nestedModule() {
	c.CheckIncorrectNamespaceAliasRequireWarning(
		"foo.bar.Baz",
		"goog.module('foo.bar'); goog.module.declareLegacyNamespace();",
		"goog.module('foo.bar.Baz'); goog.module.declareLegacyNamespace();",
		`goog.module('another');
const {Baz} = goog.require('foo.bar');
function ref(a) {};
ref(Baz);
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_nestedModule()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_nestedModule() {
	c.CheckIncorrectNamespaceAliasRequireTypeWarning(
		"foo.bar.Baz",
		"goog.module('foo.bar'); goog.module.declareLegacyNamespace();",
		"goog.module('foo.bar.Baz'); goog.module.declareLegacyNamespace(); exports = class {};",
		`goog.module('another');
const {Baz} = goog.requireType('foo.bar');
let /** Baz */ x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_destructure_nonLegacyModule_fromLegacyModule()
func (c *CheckMissingRequiresTest) TestWarning_destructure_nonLegacyModule_fromLegacyModule() {
	c.CheckIncorrectNamespaceAliasRequireWarning(
		"foo.bar.Baz",
		`goog.module('foo.bar');
goog.module.declareLegacyNamespace();
`,
		`goog.module('foo.bar.Baz');
/** @constructor */ exports = function() {};
`,
		`goog.module('test');
const {Baz} = goog.require('foo.bar');
function ref(a) {}
ref(Baz);
`)
}

// Java: CheckMissingRequiresTest.testWarning_destructure_nestedProvide_whenChildAlsoImported()
func (c *CheckMissingRequiresTest) TestWarning_destructure_nestedProvide_whenChildAlsoImported() {
	c.CheckIncorrectNamespaceAliasRequireWarning(
		"foo.bar.Baz",
		"goog.provide('foo.bar');",
		"goog.provide('foo.bar.Baz');",
		`goog.module('test');
const BarBaz = goog.require('foo.bar.Baz');
const {Baz} = goog.require('foo.bar');
function ref(a) {}
ref(Baz);
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_nestedProvideIndirectRef()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_nestedProvideIndirectRef() {
	c.CheckIndirectNamespaceRefRequireWarning(
		"foo.bar.Baz",
		"goog.provide('foo.bar');",
		"goog.provide('foo.bar.Baz');",
		`goog.module('another');
const bar = goog.require('foo.bar');
function ref(a) {};
ref(bar.Baz);
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_nestedProvideIndirectRef()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_nestedProvideIndirectRef() {
	c.CheckIndirectNamespaceRefRequireTypeWarning(
		"foo.bar.Baz",
		"goog.provide('foo.bar');",
		"goog.provide('foo.bar.Baz');foo.bar.Baz = class {};",
		`goog.module('another');
const bar = goog.requireType('foo.bar');
let /** bar.Baz */ x;
`)
}

// Java: CheckMissingRequiresTest.testIndirectNamespaceRef_errorLocation()
func (c *CheckMissingRequiresTest) TestIndirectNamespaceRef_errorLocation() {
	c.Test(
		c.Srcs(
			"goog.provide('goog.dom');",
			"goog.provide('goog.dom.TagName'); goog.dom.TagName = {B: 'b'};",
			`goog.module('test');
const dom = goog.require('goog.dom');
const x = dom.TagName.B;
`),
		c.Warning(CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE).
			WithMessage(
				"'goog.dom.TagName' should have its own goog.require.\n"+
					"Please add a separate goog.require and use that alias instead.").
			WithLocation(3, 10, 11))
}

// Java: CheckMissingRequiresTest.testIndirectNamespaceRef_jsdoc_errorLocation()
func (c *CheckMissingRequiresTest) TestIndirectNamespaceRef_jsdoc_errorLocation() {
	c.Test(
		c.Srcs(
			"goog.provide('goog.dom');",
			`goog.provide('goog.dom.TagName');
/** @const */
goog.dom.TagName = {};
goog.dom.TagName.C = class {};
`,
			`goog.module('test');
const dom = goog.require('goog.dom');
/** @type {dom.TagName.C} */
const x = null;
`),
		c.Warning(CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE).
			WithMessage(
				"'goog.dom.TagName' should have its own goog.requireType.\n"+
					"Please add a separate goog.requireType and use that alias instead.").
			WithLocation(3, 11, 11))
}

// Java: CheckMissingRequiresTest.testIndirectNamespaceRef_jsdoc_errorLocation_withFunction()
func (c *CheckMissingRequiresTest) TestIndirectNamespaceRef_jsdoc_errorLocation_withFunction() {
	c.Test(
		c.Srcs(
			"goog.provide('goog.dom');",
			"goog.provide('goog.dom.TagName'); /** @enum {string} */ goog.dom.TagName = {B: 'b'};",
			`goog.module('test');
const dom = goog.require('goog.dom');
/** @type {function(dom.TagName)} */
const x = null;
`),
		c.Warning(CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE).
			WithMessage(
				"'goog.dom.TagName' should have its own goog.requireType.\n"+
					"Please add a separate goog.requireType and use that alias instead.").
			WithLocation(3, 20, 11))
}

// Java: CheckMissingRequiresTest.testWarning_missingRequire_nestedModuleIndirectRef()
func (c *CheckMissingRequiresTest) TestWarning_missingRequire_nestedModuleIndirectRef() {
	c.CheckIndirectNamespaceRefRequireWarning(
		"foo.bar.Baz",
		"goog.module('foo.bar'); goog.module.declareLegacyNamespace();",
		"goog.module('foo.bar.Baz'); goog.module.declareLegacyNamespace();",
		`goog.module('another');
const bar = goog.require('foo.bar');
function ref(a) {};
ref(bar.Baz);
`)
}

// Java: CheckMissingRequiresTest.testWarning_missingRequireType_nestedModuleIndirectRef()
func (c *CheckMissingRequiresTest) TestWarning_missingRequireType_nestedModuleIndirectRef() {
	c.CheckIndirectNamespaceRefRequireTypeWarning(
		"foo.bar.Baz",
		"goog.module('foo.bar'); goog.module.declareLegacyNamespace();",
		"goog.module('foo.bar.Baz'); goog.module.declareLegacyNamespace();exports = class {};",
		`goog.module('another');
const bar = goog.requireType('foo.bar');
let /** bar.Baz */ x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_wrongAlias_nestedModuleIndirectRef()
func (c *CheckMissingRequiresTest) TestWarning_wrongAlias_nestedModuleIndirectRef() {
	c.CheckIndirectNamespaceRefRequireTypeWarning(
		"foo.bar.Baz",
		"goog.module('foo.bar'); goog.module.declareLegacyNamespace();",
		"goog.module('foo.bar.Baz'); goog.module.declareLegacyNamespace();exports = class {};",
		`goog.module('another');
const bar = goog.requireType('foo.bar');
const Baz = goog.requireType('foo.bar.Baz');
let /** bar.Baz */ x;
`) // error should be "Baz"
}

// Java: CheckMissingRequiresTest.testWarning_wrongAlias_nestedModuleIndirectRef2()
func (c *CheckMissingRequiresTest) TestWarning_wrongAlias_nestedModuleIndirectRef2() {
	c.CheckIndirectNamespaceRefRequireWarning(
		"foo.bar.Baz",
		"goog.module('foo.bar'); goog.module.declareLegacyNamespace();",
		"goog.module('foo.bar.Baz'); goog.module.declareLegacyNamespace();exports = class {};",
		`goog.module('another');
const bar = goog.require('foo.bar');
const Baz = goog.requireType('foo.bar.Baz');
let x = bar.Baz;
`) // error should be "Baz"
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequire_nestedDirectRef()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequire_nestedDirectRef() {
	c.CheckNoWarning(
		"goog.module('foo.bar'); goog.module.declareLegacyNamespace();",
		"goog.module('foo.bar.Baz'); goog.module.declareLegacyNamespace();exports = class {};",
		`goog.module('another');
const Baz = goog.require('foo.bar.Baz');
let x = new Baz.C();
`)
}

// Java: CheckMissingRequiresTest.testNoWarningForRequireType_nestedDirectRef()
func (c *CheckMissingRequiresTest) TestNoWarningForRequireType_nestedDirectRef() {
	c.CheckNoWarning(
		"goog.module('foo.bar'); goog.module.declareLegacyNamespace();",
		"goog.module('foo.bar.Baz'); goog.module.declareLegacyNamespace();exports = class {};",
		`goog.module('another');
const Baz = goog.requireType('foo.bar.Baz');
let /** !Baz.C */ x = null;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_overlapping_module_id_and_legacy_namespace()
func (c *CheckMissingRequiresTest) TestNoWarning_overlapping_module_id_and_legacy_namespace() {
	c.CheckNoWarning(
		"goog.module('jspb'); exports = {Message: class {}};",
		`goog.module('jspb.Message');
goog.module.declareLegacyNamespace();
const {Message} = goog.require('jspb');
exports = Message;
`,
		`goog.module('another');
const {Message} = goog.requireType('jspb');
let /** !Message */ x = null;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_overlapping_legacy_namespaces()
func (c *CheckMissingRequiresTest) TestNoWarning_overlapping_legacy_namespaces() {
	c.CheckNoWarning(
		"goog.provide('wiz');",
		"goog.provide('wiz.controller');",
		`goog.module('wiz.controller.idomcompatiblecontroller');
goog.module.declareLegacyNamespace();
class IdomCompatibleController {}
exports = {IdomCompatibleController};
`,
		`goog.module('another');
const SomeRandomThing = goog.require('wiz.controller.idomcompatiblecontroller');
let /** !SomeRandomThing.IdomCompatibleController */ x = null;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_missingRequireType_sameModuleWithLegacyNamespace_nestedProvide()
func (c *CheckMissingRequiresTest) TestNoWarning_missingRequireType_sameModuleWithLegacyNamespace_nestedProvide() {
	c.CheckNoWarning(
		"goog.provide('foo.bar');",
		`goog.module('foo.bar.Baz');
goog.module.declareLegacyNamespace();
/** @constructor */
exports.Baz = function() {}
/** @type {!foo.bar.Baz} */ let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocParam()
func (c *CheckMissingRequiresTest) TestWarning_jsDocParam() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @param {!foo.Bar} x */
function f(x) {}
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocThis()
func (c *CheckMissingRequiresTest) TestWarning_jsDocThis() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @this {foo.Bar} */
function f() {}
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocReturn()
func (c *CheckMissingRequiresTest) TestWarning_jsDocReturn() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @return {!foo.Bar} */
function f() {}
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocEnum()
func (c *CheckMissingRequiresTest) TestWarning_jsDocEnum() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @enum {foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocTypedef()
func (c *CheckMissingRequiresTest) TestWarning_jsDocTypedef() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @typedef {foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocExtendsClass()
func (c *CheckMissingRequiresTest) TestWarning_jsDocExtendsClass() {
	c.CheckRequireWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/**
 * @constructor
 * @extends {foo.Bar}
 */
function Bar() {}
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocExtendsInterface()
func (c *CheckMissingRequiresTest) TestWarning_jsDocExtendsInterface() {
	c.CheckRequireWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @interface */
exports = function() {}
`,
		`goog.module('test');
/**
 * @interface
 * @extends {foo.Bar}
 */
function Bar() {}
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocImplements()
func (c *CheckMissingRequiresTest) TestWarning_jsDocImplements() {
	c.CheckRequireWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @interface */
exports = function() {}
`,
		`goog.module('test');
/**
 * @constructor
 * @implements {foo.Bar}
 */
function Bar() {}
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocUnion()
func (c *CheckMissingRequiresTest) TestWarning_jsDocUnion() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @type {string|foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocTemplate()
func (c *CheckMissingRequiresTest) TestWarning_jsDocTemplate() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @template T */
class Quux {}
/** @type {!Quux<!foo.Bar>} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocRecord()
func (c *CheckMissingRequiresTest) TestWarning_jsDocRecord() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @type {{foo: !foo.Bar}} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocFunctionParam()
func (c *CheckMissingRequiresTest) TestWarning_jsDocFunctionParam() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @type {function(!foo.Bar)} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocFunctionReturn()
func (c *CheckMissingRequiresTest) TestWarning_jsDocFunctionReturn() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @type {function():!foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocFunctionNew()
func (c *CheckMissingRequiresTest) TestWarning_jsDocFunctionNew() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @type {function(new:foo.Bar)} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocFunctionThis()
func (c *CheckMissingRequiresTest) TestWarning_jsDocFunctionThis() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @type {function(this:foo.Bar)} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testWarning_jsDocTypeof()
func (c *CheckMissingRequiresTest) TestWarning_jsDocTypeof() {
	c.CheckRequireTypeWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
/** @constructor */
exports = function() {}
`,
		`goog.module('test');
/** @type {typeof foo.Bar} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_jsDocClassTemplateParam()
func (c *CheckMissingRequiresTest) TestNoWarning_jsDocClassTemplateParam() {
	c.CheckNoWarning(
		`goog.module('MyType');
goog.module.declareLegacyNamespace();
`,
		`goog.module('test');
/** @template MyType */
class Foo {
  /** @param {!MyType} x */
  constructor(x) {}
}
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_jsDocFunctionTemplateParam()
func (c *CheckMissingRequiresTest) TestNoWarning_jsDocFunctionTemplateParam() {
	c.CheckNoWarning(
		`goog.module('MyType');
goog.module.declareLegacyNamespace();
`,
		`goog.module('test');
/**
 * @param {!MyType} x
 * @return {!MyType}
 * @template MyType
 */
function foo(x) {
  return x;
}
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_jsDocTtlParam()
func (c *CheckMissingRequiresTest) TestNoWarning_jsDocTtlParam() {
	c.CheckNoWarning(
		`goog.module('MyType');
goog.module.declareLegacyNamespace();
`,
		`goog.module('YourType');
goog.module.declareLegacyNamespace();
`,
		`goog.module('test');
/**
 * @param {!MyType} x
 * @return {!YourType}
 * @template MyType
 * @template YourType := MyType =:
 */
function foo(x) {
  return x;
}
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_shadow_moduleScope()
func (c *CheckMissingRequiresTest) TestNoWarning_shadow_moduleScope() {
	c.CheckNoWarning(
		`goog.module('Bar');
goog.module.declareLegacyNamespace();
`,
		`goog.module('test');
class Bar {}
/** @type {!Bar} */
let x = new Bar();
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_shadow_localScope()
func (c *CheckMissingRequiresTest) TestNoWarning_shadow_localScope() {
	c.CheckNoWarning(
		`goog.module('Bar');
goog.module.declareLegacyNamespace();
`,
		`goog.module('test');
function foo() {
  class Bar {}
  /** @type {!Bar} */
  let x = new Bar();
}
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_shadow_parentScope()
func (c *CheckMissingRequiresTest) TestNoWarning_shadow_parentScope() {
	c.CheckNoWarning(
		`goog.module('Bar');
goog.module.declareLegacyNamespace();
`,
		`goog.module('test');
class Bar {}
function foo() {
  /** @type {!Bar} */
  let x = new Bar();
}
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_providedGoogModule()
func (c *CheckMissingRequiresTest) TestNoWarning_providedGoogModule() {
	c.CheckNoWarning(
		`goog.provide('goog.module');
goog.module = goog.module || {};
`,
		`goog.module('test');
goog.module.declareLegacyNamespace();
`)
}

// Java: CheckMissingRequiresTest.testWarning_providedGoogModule_missingRequireForNestedProvide()
func (c *CheckMissingRequiresTest) TestWarning_providedGoogModule_missingRequireForNestedProvide() {
	c.CheckRequireWarning(
		"goog.module.ModuleManager",
		`goog.provide('goog.module');
goog.module = goog.module || {};
`,
		`goog.provide('goog.module.ModuleManager');
/** @constructor */
goog.module.ModuleManager = function() {};
`,
		`goog.module('test');
let x = new goog.module.ModuleManager();
`)
}

// Java: CheckMissingRequiresTest.testWarning_providedGoogModule_missingRequireTypeForNestedProvide()
func (c *CheckMissingRequiresTest) TestWarning_providedGoogModule_missingRequireTypeForNestedProvide() {
	c.CheckRequireTypeWarning(
		"goog.module.ModuleManager",
		`goog.provide('goog.module');
goog.module = goog.module || {};
`,
		`goog.provide('goog.module.ModuleManager');
/** @constructor */
goog.module.ModuleManager = function() {};
`,
		`goog.module('test');
/** @type {!goog.module.ModuleManager} */
let x;
`)
}

// Java: CheckMissingRequiresTest.testNoCrashOnGetprops()
func (c *CheckMissingRequiresTest) TestNoCrashOnGetprops() {
	c.CheckNoWarning(
		`goog.module('test');
class Class {
  constructor() {
    foo().bar;
    foo['bar'].baz;
    this.foo;
    super.foo;
  }
}
`)
}

// Java: CheckMissingRequiresTest.testReferenceNonLegacyGoogModule_inScript_required_warns()
func (c *CheckMissingRequiresTest) TestReferenceNonLegacyGoogModule_inScript_required_warns() {
	c.Test(
		c.Srcs(
			"goog.module('test.a'); exports.x = 1;",
			"goog.provide('test.b'); goog.require('test.a'); console.log(test.a.x);"),
		c.Error(CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE))
}

// Java: CheckMissingRequiresTest.testReferenceNonLegacyGoogModule_inScript_typePosition_required_noWarning()
func (c *CheckMissingRequiresTest) TestReferenceNonLegacyGoogModule_inScript_typePosition_required_noWarning() {
	c.CheckNoWarning(
		`goog.module('test.a');
/** @interface */
exports.C = class {};
`,
		`goog.provide('test.b');
goog.require('test.a');
/** @implements {test.a.C} */
test.b.C = class {
  /** @param {!test.a.C} c */
  foo(c) {}
};
`)
}

// Java: CheckMissingRequiresTest.testReferenceNonLegacyGoogModule_inScript_typePosition_and_code_required_warning()
func (c *CheckMissingRequiresTest) TestReferenceNonLegacyGoogModule_inScript_typePosition_and_code_required_warning() {
	c.Test(
		c.Srcs(
			`goog.module('test.a');
/** @interface */
exports.C = class {};
`,
			`goog.provide('test.b');
goog.require('test.a');
/** @implements {test.a.C} */
test.b.C = class {
  /** @param {!test.a.C} c */
  foo(c) { test.a.C; } // reports error
};
`),
		c.Error(CheckMissingRequires_NON_LEGACY_GOOG_MODULE_REFERENCE))
}

// TestTesNonLegacyGoogModule_inScript_required_noWarning is the @Test method
// tesNonLegacyGoogModule_inScript_required_noWarning(); RunJUnit4 runs the
// methods whose name starts with Test (Go-only: language-adapter).
//
// Java: CheckMissingRequiresTest.tesNonLegacyGoogModule_inScript_required_noWarning()
func (c *CheckMissingRequiresTest) TestTesNonLegacyGoogModule_inScript_required_noWarning() {
	c.CheckNoWarning(
		"goog.module('test.a'); goog.module.declareLegacyNamespace(); exports.x = 1;",
		"goog.provide('test.b'); goog.require('test.a'); console.log(test.a.x);")
}

// Java: CheckMissingRequiresTest.testWarning_googScope_googModuleGet_noRequire()
func (c *CheckMissingRequiresTest) TestWarning_googScope_googModuleGet_noRequire() {
	c.CheckRequireForGoogModuleGetWarning(
		"foo.bar.Baz",
		`goog.module('foo.bar.Baz');
/** @constructor */
exports = function() {}
`,
		`goog.provide('test');

goog.scope(function() {
const Baz = goog.module.get('foo.bar.Baz');
exports.fn = function() { new Baz(); }
});
`)
}

// Java: CheckMissingRequiresTest.testWarning_script_googModuleGet_inProvideInitialization_noRequire()
func (c *CheckMissingRequiresTest) TestWarning_script_googModuleGet_inProvideInitialization_noRequire() {
	c.CheckRequireForGoogModuleGetWarning(
		"foo.bar.Baz",
		`goog.module('foo.bar.Baz');
/** @constructor */
exports.Boo = function() {}
`,
		`goog.provide('test');

test.boo = new (goog.module.get('foo.bar.Baz').Boo)();
`)
}

// Java: CheckMissingRequiresTest.testWarning_script_googModuleGet_inScriptIfBlock_noRequire()
func (c *CheckMissingRequiresTest) TestWarning_script_googModuleGet_inScriptIfBlock_noRequire() {
	c.CheckRequireForGoogModuleGetWarning(
		"foo.bar.Baz",
		`goog.module('foo.bar.Baz');
/** @constructor */
exports.Boo = function() {}
`,
		`goog.provide('test');

if (true) { console.log(new (goog.module.get('foo.bar.Baz').Boo)()); }
`)
}

// Java: CheckMissingRequiresTest.testWarning_script_googModuleGet_inScript_inTopLevelIIFE_noRequire()
func (c *CheckMissingRequiresTest) TestWarning_script_googModuleGet_inScript_inTopLevelIIFE_noRequire() {
	c.CheckRequireForGoogModuleGetWarning(
		"foo.bar.Baz",
		`goog.module('foo.bar.Baz');
/** @constructor */
exports.Boo = function() {}
`,
		`goog.provide('test');

(() => console.log(new (goog.module.get('foo.bar.Baz').Boo)()))();
`)
}

// Java: CheckMissingRequiresTest.testWarning_googScope_googModuleGet_withPropAccess_noRequire()
func (c *CheckMissingRequiresTest) TestWarning_googScope_googModuleGet_withPropAccess_noRequire() {
	c.CheckRequireForGoogModuleGetWarning(
		"foo.bar.Baz",
		`goog.module('foo.bar.Baz');
/** @constructor */
exports.Boo = function() {}
`,
		`goog.provide('test');

goog.scope(function() {
const Woo = goog.module.get('foo.bar.Baz').Boo.Goo.Woo;
exports.fn = function() { new Woo(); }
});
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_googScope_googModuleGet_hasRequire()
func (c *CheckMissingRequiresTest) TestNoWarning_googScope_googModuleGet_hasRequire() {
	c.CheckNoWarning(
		`goog.module('foo.bar.Baz');
/** @constructor */
exports = function() {}
`,
		`goog.provide('test');
goog.require('foo.bar.Baz');

goog.scope(function() {
const Baz = goog.module.get('foo.bar.Baz');
function fn() { new Baz(); }
});
`)
}

// Java: CheckMissingRequiresTest.testWarning_googScope_googModuleGet_IIFE_noRequire()
func (c *CheckMissingRequiresTest) TestWarning_googScope_googModuleGet_IIFE_noRequire() {
	c.CheckRequireForGoogModuleGetWarning(
		"foo.bar.Baz",
		`goog.module('foo.bar.Baz');
/** @constructor */
exports.Boo = function() {}
`,
		`goog.provide('test');

goog.scope(function() {
  (() => console.log(new (goog.module.get('foo.bar.Baz').Boo)()))();
});
`)
}

// Java: CheckMissingRequiresTest.testWarning_script_googModuleGet_inFunctionInScript_noRequire_noWarning()
func (c *CheckMissingRequiresTest) TestWarning_script_googModuleGet_inFunctionInScript_noRequire_noWarning() {
	c.CheckNoWarning(
		`goog.module('foo.bar.Baz');
/** @constructor */
exports.Boo = function() {}
`,
		`goog.provide('test');

// We allow this goog.module.get despite the lack of goog.require: it's possible that
// the foo.bar.Baz module will have been loaded by the time test.fn is called. It's up
// to the caller to ensure that it's loaded.
test.fn = function() {
  console.log(new (goog.module.get('foo.bar.Baz').Boo)());
}
`)
}

// Java: CheckMissingRequiresTest.testWarning_script_googModuleGet_inFunctionInScript_inNonTopLevelIIFE_noRequire_noWarning()
func (c *CheckMissingRequiresTest) TestWarning_script_googModuleGet_inFunctionInScript_inNonTopLevelIIFE_noRequire_noWarning() {
	c.CheckNoWarning(
		`goog.module('foo.bar.Baz');
/** @constructor */
exports.Boo = function() {}
`,
		`goog.provide('test');

// We allow this goog.module.get despite the lack of goog.require: it's possible that
// the foo.bar.Baz module will have been loaded by the time test.fn is called. It's up
// to the caller to ensure that it's loaded.
test.fn = function() {
  (() => console.log(new (goog.module.get('foo.bar.Baz').Boo)()))();
}
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_nestedInGoogScope_googModuleGet_noRequire()
func (c *CheckMissingRequiresTest) TestNoWarning_nestedInGoogScope_googModuleGet_noRequire() {
	c.CheckNoWarning(
		`goog.module('foo.bar.Baz');
/** @constructor */
exports = function() {}
`,
		`goog.provide('test');

goog.scope(function() {
test.fn = function() {
// We allow this goog.module.get despite the lack of goog.require: it's possible that
// the foo.bar.Baz module will have been loaded by the time test.fn is called. It's up
// to the caller to ensure that it's loaded.
  const Baz = goog.module.get('foo.bar.Baz');
  new Baz();
  new goog.module.get('foo.bar.Baz');
}
});
`)
}

// Java: CheckMissingRequiresTest.testWarning_googScopeAliasOfModule()
func (c *CheckMissingRequiresTest) TestWarning_googScopeAliasOfModule() {
	c.CheckRequireInProvidesFileWarning(
		"foo.Bar",
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
		`goog.provide('test');

goog.scope(function() {
const Bar = foo.Bar;
})(); // end goog.scope
`)
}

// Java: CheckMissingRequiresTest.testNoWarning_googScopeAliasOfModule()
func (c *CheckMissingRequiresTest) TestNoWarning_googScopeAliasOfModule() {
	c.CheckNoWarning(
		`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
exports = class {};
exports.Baz = class { static fn() {} };
`,
		`goog.provide('test');
goog.require('foo.Bar');

goog.scope(function() {
const Bar = foo.Bar;
const Baz = foo.Bar.Baz;
test.fn = function() {
  new Bar(); Bar.fn();
  new Baz(); Baz.fn();
  new Bar.Baz(); Bar.Baz.fn();
}
})(); // end goog.scope
`)
}

// Java: CheckMissingRequiresTest.testWarning_googScopeAlias_unprovidedParentOfModule()
func (c *CheckMissingRequiresTest) TestWarning_googScopeAlias_unprovidedParentOfModule() {
	c.Test(
		c.Srcs(
			`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
			`goog.provide('test');

goog.scope(function() {
const fooLocal = foo;
function fn() { return fooLocal.Bar; }
})(); // end goog.scope
`),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE))
}

// Java: CheckMissingRequiresTest.testWarning_googScopeAlias_unprovidedParentOfModule_typeRef()
func (c *CheckMissingRequiresTest) TestWarning_googScopeAlias_unprovidedParentOfModule_typeRef() {
	c.Test(
		c.Srcs(
			`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
			`goog.provide('test');

goog.scope(function() {
const fooLocal = foo;
function fn() { /** @type {!fooLocal.Bar} */ var x; }
})(); // end goog.scope
`),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE))
}

// Java: CheckMissingRequiresTest.testWarning_googScopeAlias_providedParentOfModule()
func (c *CheckMissingRequiresTest) TestWarning_googScopeAlias_providedParentOfModule() {
	c.Test(
		c.Srcs(
			"goog.provide('foo');",
			`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
			`goog.provide('test');

goog.scope(function() {
const fooLocal = foo;
function fn() { return fooLocal.Bar; }
})(); // end goog.scope
`),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE).
			WithMessageContaining("'foo'"),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE).
			WithMessageContaining("'foo.Bar'"))
}

// Java: CheckMissingRequiresTest.testWarning_googScopeAlias_providedParentOfModule_typeRef()
func (c *CheckMissingRequiresTest) TestWarning_googScopeAlias_providedParentOfModule_typeRef() {
	c.Test(
		c.Srcs(
			"goog.provide('foo');",
			`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
			`goog.provide('test');

goog.scope(function() {
const fooLocal = foo;
function fn() { /** @type {!fooLocal.Bar} */ var x; }
})(); // end goog.scope
`),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE).WithMessageContaining("foo"),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE).
			WithMessageContaining("foo.Bar"))
}

// Java: CheckMissingRequiresTest.testWarning_googScopeAlias_providedParentOfModule_requireParent()
func (c *CheckMissingRequiresTest) TestWarning_googScopeAlias_providedParentOfModule_requireParent() {
	c.Test(
		c.Srcs(
			"goog.provide('foo');",
			`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
			`goog.provide('test');
goog.require('foo');

goog.scope(function() {
const fooLocal = foo;
function fn() { return fooLocal.Bar; }
})(); // end goog.scope
`),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE).
			WithMessageContaining("'foo.Bar'"))
}

// Java: CheckMissingRequiresTest.testWarning_googScopeAlias_providedParentOfModule_googModuleGetParent()
func (c *CheckMissingRequiresTest) TestWarning_googScopeAlias_providedParentOfModule_googModuleGetParent() {
	c.Test(
		c.Srcs(
			"goog.provide('foo');",
			`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
			`goog.provide('test');
goog.require('foo');

goog.scope(function() {
const fooLocal = goog.module.get('foo');
function fn() { return fooLocal.Bar; }
})(); // end goog.scope
`),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE).
			WithMessageContaining("'foo.Bar'"))
}

// Java: CheckMissingRequiresTest.testWarning_googScopeAlias_providedParentOfModule_googModuleGetParent_propRef()
func (c *CheckMissingRequiresTest) TestWarning_googScopeAlias_providedParentOfModule_googModuleGetParent_propRef() {
	c.Test(
		c.Srcs(
			"goog.provide('foo');",
			`goog.module('foo.bar');
goog.module.declareLegacyNamespace();
`,
			`goog.module('foo.bar.Raz.Baz');
goog.module.declareLegacyNamespace();
exports = class {};
`,
			`goog.provide('test');
goog.require('foo.bar');

goog.scope(function() {
const fooBarRaz = goog.module.get('foo.bar').Raz;
function fn() { return fooBarRaz.Baz; }
})(); // end goog.scope
`),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE).
			WithMessageContaining("'foo.bar.Raz.Baz'"))
}

// Java: CheckMissingRequiresTest.testWarning_googScopeAlias_providedParentOfModule_typeRef_requireParent()
func (c *CheckMissingRequiresTest) TestWarning_googScopeAlias_providedParentOfModule_typeRef_requireParent() {
	c.Test(
		c.Srcs(
			`goog.provide('foo');
`,
			`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
			`goog.provide('test');
goog.require('foo');

goog.scope(function() {
const fooLocal = foo;
function fn() { /** @type {!fooLocal.Bar} */ var x; }
})(); // end goog.scope
`),
		// this warning could reasonably be INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE, but
		// the extra complexity to distinguish the two cases doesn't seem worthwhile.
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE).
			WithMessageContaining("foo.Bar"))
}

// Java: CheckMissingRequiresTest.testWarning_googScopeAliasOfParentOfModule_childStillRequired()
func (c *CheckMissingRequiresTest) TestWarning_googScopeAliasOfParentOfModule_childStillRequired() {
	c.Test(
		c.Srcs(
			`goog.module('foo.Bar');
goog.module.declareLegacyNamespace();
/** @constructor */
exports = function() {};
`,
			"goog.provide('test');\n// Note that there is a require for foo.Bar, which is correct,  but the warning is\n// still triggered because of the alias `const fooLocal = foo;`.\ngoog.require('foo.Bar');\n\ngoog.scope(function() {\nconst fooLocal = foo;\nfunction fn() { return fooLocal.Bar; }\n})(); // end goog.scope\n"),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE))
}

// Java: CheckMissingRequiresTest.checkNoWarning(String...)
func (c *CheckMissingRequiresTest) CheckNoWarning(js ...string) {
	c.TestSame(c.Srcs(js))
}

// Java: CheckMissingRequiresTest.checkRequireWarning(String, String...)
func (c *CheckMissingRequiresTest) CheckRequireWarning(namespace string, js ...string) {
	c.TestSame(
		c.Srcs(js),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE).WithMessageContaining("'"+namespace+"'"))
}

// Java: CheckMissingRequiresTest.checkRequireTypeWarning(String, String...)
func (c *CheckMissingRequiresTest) CheckRequireTypeWarning(namespace string, js ...string) {
	c.TestSame(
		c.Srcs(js),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_TYPE).
			WithMessageContaining("'"+namespace+"'"))
}

// Java: CheckMissingRequiresTest.checkRequireInProvidesFileWarning(String, String...)
func (c *CheckMissingRequiresTest) CheckRequireInProvidesFileWarning(namespace string, js ...string) {
	c.TestSame(
		c.Srcs(js),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE).
			WithMessageContaining("'"+namespace+"'"))
}

// Java: CheckMissingRequiresTest.checkRequireTypeInProvidesFileWarning(String, String...)
func (c *CheckMissingRequiresTest) CheckRequireTypeInProvidesFileWarning(namespace string, js ...string) {
	c.TestSame(
		c.Srcs(js),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE).
			WithMessageContaining("'"+namespace+"'"))
}

// Java: CheckMissingRequiresTest.checkIndirectNamespaceRefRequireWarning(String, String...)
func (c *CheckMissingRequiresTest) CheckIndirectNamespaceRefRequireWarning(namespace string, js ...string) {
	c.TestSame(
		c.Srcs(js),
		c.Warning(CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE).
			WithMessageContaining("'"+namespace+"'"))
}

// Java: CheckMissingRequiresTest.checkIndirectNamespaceRefRequireTypeWarning(String, String...)
func (c *CheckMissingRequiresTest) CheckIndirectNamespaceRefRequireTypeWarning(namespace string, js ...string) {
	c.TestSame(
		c.Srcs(js),
		c.Warning(CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE).
			WithMessageContaining("'"+namespace+"'"))
}

// Java: CheckMissingRequiresTest.checkIncorrectNamespaceAliasRequireWarning(String, String...)
func (c *CheckMissingRequiresTest) CheckIncorrectNamespaceAliasRequireWarning(namespace string, js ...string) {
	c.TestSame(
		c.Srcs(js),
		c.Warning(CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE).
			WithMessageContaining("'"+namespace+"'"))
}

// Java: CheckMissingRequiresTest.checkIncorrectNamespaceAliasRequireTypeWarning(String, String...)
func (c *CheckMissingRequiresTest) CheckIncorrectNamespaceAliasRequireTypeWarning(namespace string, js ...string) {
	c.TestSame(
		c.Srcs(js),
		c.Warning(CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE).
			WithMessageContaining("'"+namespace+"'"))
}

// Java: CheckMissingRequiresTest.checkRequireForGoogModuleGetWarning(String, String...)
func (c *CheckMissingRequiresTest) CheckRequireForGoogModuleGetWarning(namespace string, js ...string) {
	c.TestSame(
		c.Srcs(js),
		c.Warning(CheckMissingRequires_MISSING_REQUIRE_FOR_GOOG_MODULE_GET).
			WithMessageContaining("'"+namespace+"'"))
}

// Java: CheckMissingRequiresTest.testIndirectNamespaceRef_jsdoc_errorLocation_longNamespace()
func (c *CheckMissingRequiresTest) TestIndirectNamespaceRef_jsdoc_errorLocation_longNamespace() {
	c.Test(
		c.Srcs(
			"goog.provide('goog.dom');",
			"goog.provide('goog.dom.tags');",
			"goog.provide('goog.dom.tags.misc'); /** @const */ goog.dom.tags.misc = {FOO: 'foo'};",
			`goog.module('test');
const dom = goog.require('goog.dom');
/** @type {dom.tags.misc.FOO} */
const x = null;
`),
		c.Warning(CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE).
			WithMessage(
				"'goog.dom.tags.misc' should have its own goog.requireType.\n"+
								"Please add a separate goog.requireType and use that alias instead.").
			WithLocation(3, 11, 13)) // dom.tags.misc
}

// Port of com/google/javascript/jscomp/PolymerBehaviorExtractorTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/PolymerBehaviorExtractorTest.java)

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PolymerBehaviorExtractorTest is the unit tests for PolymerBehaviorExtractor.
type PolymerBehaviorExtractorTest struct {
	CompilerTypeTestCaseBase

	extractor         *PolymerBehaviorExtractor
	behaviorArray     rhino.Node
	moduleMetadataMap *ModuleMetadataMap
}

// Java: PolymerBehaviorExtractorTest()
func NewPolymerBehaviorExtractorTest() *PolymerBehaviorExtractorTest {
	c := &PolymerBehaviorExtractorTest{}
	c.InitCompilerTypeTestCaseBase(c)
	return c
}

func TestPolymerBehaviorExtractorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPolymerBehaviorExtractorTest() })
}

// Java: PolymerBehaviorExtractorTest.setUp()
func (c *PolymerBehaviorExtractorTest) SetUp() {
	c.CompilerTypeTestCaseBase.SetUp()
	c.behaviorArray = rhino.NilNode
}

// Java: PolymerBehaviorExtractorTest.testArrayBehavior()
func (c *PolymerBehaviorExtractorTest) TestArrayBehavior() {
	c.parseAndInitializeExtractor(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
  /** @param {string} funAmount */
  doSomethingFun: function(funAmount) { alert('Something ' + funAmount + ' fun!'); },
  /** @override */
  created: function() {}
};
/** @polymerBehavior */
var RadBehavior = {
  properties: {
    howRad: Number
  },
  /** @param {number} radAmount */
  doSomethingRad: function(radAmount) { alert('Something ' + radAmount + ' rad!'); },
  /** @override */
  ready: function() {}
};
/** @polymerBehavior */
var SuperCoolBehaviors = [FunBehavior, RadBehavior];
/** @polymerBehavior */
var BoringBehavior = {
  properties: {
    boringString: String
  },
  /** @param {boolean} boredYet */
  doSomething: function(boredYet) { alert(boredYet + ' ' + this.boringString); },
};
var A = Polymer({
  is: 'x-element',
  behaviors: [ SuperCoolBehaviors, BoringBehavior ],
});
`)

	defs := c.extractor.ExtractBehaviors(c.behaviorArray, nil)
	AssertThat(defs).HasSize(3)

	// TODO(jlklein): Actually verify the properties of the BehaviorDefinitions.
}

// Java: PolymerBehaviorExtractorTest.testArrayBehavior_fromGoogProvide()
func (c *PolymerBehaviorExtractorTest) TestArrayBehavior_fromGoogProvide() {
	c.parseAndInitializeExtractor(
		`goog.provide('my.project.RadBehavior');
goog.provide('my.project.FunBehavior');
/** @polymerBehavior */
my.project.FunBehavior = {
  properties: {
    isFun: Boolean
  },
  /** @param {string} funAmount */
  doSomethingFun: function(funAmount) { alert('Something ' + funAmount + ' fun!'); },
  /** @override */
  created: function() {}
};
/** @polymerBehavior */
my.project.RadBehavior = {
  properties: {
    howRad: Number
  },
  /** @param {number} radAmount */
  doSomethingRad: function(radAmount) { alert('Something ' + radAmount + ' rad!'); },
  /** @override */
  ready: function() {}
};
/** @polymerBehavior */
var SuperCoolBehaviors = [my.project.FunBehavior, my.project.RadBehavior];
/** @polymerBehavior */
var BoringBehavior = {
  properties: {
    boringString: String
  },
  /** @param {boolean} boredYet */
  doSomething: function(boredYet) { alert(boredYet + ' ' + this.boringString); },
};
var A = Polymer({
  is: 'x-element',
  behaviors: [ SuperCoolBehaviors, BoringBehavior ],
});
`)

	defs := c.extractor.ExtractBehaviors(c.behaviorArray, nil)
	AssertThat(defs).HasSize(3)
}

// Java: PolymerBehaviorExtractorTest.testBehaviorDeclaredAndUsedWithinModule()
func (c *PolymerBehaviorExtractorTest) TestBehaviorDeclaredAndUsedWithinModule() {
	c.parseAndInitializeExtractor(
		`goog.module('modules.rad');
` +
			c.createBehaviorDefinition("RadBehavior") +
			c.createBehaviorUsage("RadBehavior"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(
		c.moduleMetadataMap.GetModulesByGoogNamespace().Get("modules.rad"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_propertyOnGlobalSetInModule()
func (c *PolymerBehaviorExtractorTest) TestBehavior_propertyOnGlobalSetInModule() {
	c.parseAndInitializeExtractor(
		`function Polymer() {}
`,
		c.createBehaviorDefinition("RadBehavior")+
			`export {RadBehavior};
`,
		// Verify that we can resolve 'RadBehavior' when we see the 'Polymer.RadBehavior' assignment
		// despite the assignment being 'global', not module-level.
		`import {RadBehavior} from './testcode1';
Polymer.RadBehavior = RadBehavior;
`,
		c.createBehaviorUsage("Polymer.RadBehavior"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByPath().Get("testcode4"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_esExport_namedImport()
func (c *PolymerBehaviorExtractorTest) TestBehavior_esExport_namedImport() {
	c.parseAndInitializeExtractor(
		c.createBehaviorDefinition("RadBehaviorLocal")+
			`export {RadBehaviorLocal as RadBehavior};
`,
		`import {RadBehavior} from './testcode0'
`+
			c.createBehaviorUsage("RadBehavior"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByPath().Get("testcode1"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_esExport_invalidNamedImport()
func (c *PolymerBehaviorExtractorTest) TestBehavior_esExport_invalidNamedImport() {
	c.parseAndInitializeExtractor(
		`export const NotABehavior = {};
`,
		`import {NotABehavior} from './testcode0'
`+
			c.createBehaviorUsage("NotABehavior"))

	c.extractor.ExtractBehaviors(
		c.behaviorArray, c.moduleMetadataMap.GetModulesByPath().Get("testcode1"))
	AssertThat(c.compiler.GetErrors()).HasSize(1)
}

// Java: PolymerBehaviorExtractorTest.testBehavior_esExport_namedImport_ofArray()
func (c *PolymerBehaviorExtractorTest) TestBehavior_esExport_namedImport_ofArray() {
	c.parseAndInitializeExtractor(
		c.createBehaviorDefinition("RadBehaviorLocal")+
			`/** @polymerBehavior */
export const RadBehavior = [ RadBehaviorLocal ];
`,
		`import {RadBehavior} from './testcode0'
`+
			c.createBehaviorUsage("RadBehavior"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByPath().Get("testcode1"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_esExport_exportDeclaration()
func (c *PolymerBehaviorExtractorTest) TestBehavior_esExport_exportDeclaration() {
	c.parseAndInitializeExtractor(
		`/** @polymerBehavior */
export const RadBehavior = {
  /** @param {number} radAmount */
  doSomethingRad: function(radAmount) { alert('Something ' + radAmount + ' rad!'); },
};
`,
		`import {RadBehavior} from './testcode0'
`+
			c.createBehaviorUsage("RadBehavior"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByPath().Get("testcode1"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_esExport_defaultExport()
func (c *PolymerBehaviorExtractorTest) TestBehavior_esExport_defaultExport() {
	c.T().Skip("@Ignore") // @Ignore
	c.parseAndInitializeExtractor(
		c.createBehaviorDefinition("RadBehavior")+
			`export default RadBehavior
`,
		`import RadBehaviorImported from './testcod
`+
			c.createBehaviorUsage("RadBehaviorImported"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByPath().Get("testcode1"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_esExport_importStar()
func (c *PolymerBehaviorExtractorTest) TestBehavior_esExport_importStar() {
	c.parseAndInitializeExtractor(
		c.createBehaviorDefinition("RadBehaviorLocal")+
			`export {RadBehaviorLocal as RadBehavior};
`,
		`import * as mod from './testcode0'
`+
			c.createBehaviorUsage("mod.RadBehavior"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByPath().Get("testcode1"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_esExport_invalidImportStar()
func (c *PolymerBehaviorExtractorTest) TestBehavior_esExport_invalidImportStar() {
	c.parseAndInitializeExtractor(
		c.createBehaviorDefinition("RadBehaviorLocal")+
			`export {RadBehaviorLocal as RadBehavior};
`,
		`import * as mod from './testcode0'
`+
			// Try using 'mod' instead of 'mod.RadBehavior'.
			c.createBehaviorUsage("mod"))

	c.extractor.ExtractBehaviors(
		c.behaviorArray, c.moduleMetadataMap.GetModulesByPath().Get("testcode1"))
	AssertThat(c.compiler.GetErrors()).HasSize(1)
}

// Java: PolymerBehaviorExtractorTest.testBehavior_esExport_PropOnImportStar()
func (c *PolymerBehaviorExtractorTest) TestBehavior_esExport_PropOnImportStar() {
	c.parseAndInitializeExtractor(
		c.createBehaviorDefinition("RadBehaviorLocal")+
			`const ns = {};
ns.RadBehavior = RadBehaviorLocal;
export {ns};
`,
		`import * as mod from './testcode0'
`+
			c.createBehaviorUsage("mod.ns.RadBehavior"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByPath().Get("testcode1"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_defaultGoogProvideImport()
func (c *PolymerBehaviorExtractorTest) TestBehavior_defaultGoogProvideImport() {
	c.parseAndInitializeExtractor(
		`goog.provide('rad.RadBehavior');
`+
			c.createBehaviorDefinition("rad.RadBehavior"),
		`goog.module('a');
const RadBehaviorLocal = goog.require('rad.RadBehavior');
`+
			c.createBehaviorUsage("RadBehaviorLocal"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByGoogNamespace().Get("a"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_destructuringGoogProvideImport()
func (c *PolymerBehaviorExtractorTest) TestBehavior_destructuringGoogProvideImport() {
	c.parseAndInitializeExtractor(
		`goog.provide('rad');
/** @const */
rad.behaviors = {};
`+
			c.createBehaviorDefinition("rad.behaviors.RadBehavior"),
		`goog.module('a');
const {behaviors: b} = goog.require('rad');
`+
			c.createBehaviorUsage("b.RadBehavior"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByGoogNamespace().Get("a"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_legacyGoogModuleNamedExportOfLocalName()
func (c *PolymerBehaviorExtractorTest) TestBehavior_legacyGoogModuleNamedExportOfLocalName() {
	c.parseAndInitializeExtractor(
		`goog.module('modules.rad');
goog.module.declareLegacyNamespace();
`+
			c.createBehaviorDefinition("RadBehavior")+
			`exports = {RadBehavior};
`,
		`goog.require('modules.rad');
`+
			c.createBehaviorUsage("modules.rad.RadBehavior"))

	c.assertSingleBehaviorExtractionSucceeds()
}

// Java: PolymerBehaviorExtractorTest.testBehavior_legacyGoogModuleNamedExport_inLoadModuleCall()
func (c *PolymerBehaviorExtractorTest) TestBehavior_legacyGoogModuleNamedExport_inLoadModuleCall() {
	c.parseAndInitializeExtractor(
		`goog.loadModule(function(exports) {
goog.module('modules.rad');
goog.module.declareLegacyNamespace();
`+
			c.createBehaviorDefinition("exports.RadBehavior")+
			`  return exports;
});
`,
		`goog.require('modules.rad');
`+
			c.createBehaviorUsage("modules.rad.RadBehavior"))

	c.assertSingleBehaviorExtractionSucceeds()
}

// Java: PolymerBehaviorExtractorTest.testBehavior_defaultLegacyGoogModuleExport_exportsLocalVariable()
func (c *PolymerBehaviorExtractorTest) TestBehavior_defaultLegacyGoogModuleExport_exportsLocalVariable() {
	c.parseAndInitializeExtractor(
		`goog.module('modules.rad.RadBehavior');
goog.module.declareLegacyNamespace();
`+
			c.createBehaviorDefinition("RadBehavior")+
			`exports = RadBehavior;
`,
		`goog.require('rad.RadBehavior');
`+
			c.createBehaviorUsage("modules.rad.RadBehavior"))

	c.assertSingleBehaviorExtractionSucceeds()
}

// Java: PolymerBehaviorExtractorTest.testBehavior_defaultLegacyGoogModuleExport_assignedDirectly()
func (c *PolymerBehaviorExtractorTest) TestBehavior_defaultLegacyGoogModuleExport_assignedDirectly() {
	c.parseAndInitializeExtractor(
		`goog.module('modules.rad.RadBehavior');
goog.module.declareLegacyNamespace();
/** @polymerBehavior */
exports = {
  properties: {
    howRad: Number
  },
  /** @param {number} radAmount */
  doSomethingRad:
    function(radAmount) { alert('Something ' + radAmount + 'rad!'); },
  /** @override */
  ready: function() {}
};
`,
		`goog.require('modules.rad.RadBehavior');
`+
			c.createBehaviorUsage("modules.rad.RadBehavior"))

	c.assertSingleBehaviorExtractionSucceeds()
}

// Java: PolymerBehaviorExtractorTest.testBehavior_nonLegacyGoogModuleDoesNotResolve()
func (c *PolymerBehaviorExtractorTest) TestBehavior_nonLegacyGoogModuleDoesNotResolve() {
	c.parseAndInitializeExtractor(
		`goog.module('modules.rad');
`+
			c.createBehaviorDefinition("exports.RadBehavior"),
		`goog.require('modules.rad');
`+
			c.createBehaviorUsage("modules.rad.RadBehavior"))

	c.extractor.ExtractBehaviors(c.behaviorArray, nil)
	AssertThat(c.compiler.GetErrors()).HasSize(1)
}

// Java: PolymerBehaviorExtractorTest.testBehavior_googRequireDefaultGoogModuleExport()
func (c *PolymerBehaviorExtractorTest) TestBehavior_googRequireDefaultGoogModuleExport() {
	c.parseAndInitializeExtractor(
		`goog.module('rad.RadBehavior');
`+
			c.createBehaviorDefinition("RadBehaviorLocal")+
			`exports = RadBehaviorLocal;
`,
		`goog.module('a');
const RadBehaviorImported = goog.require('rad.RadBehavior');
`+
			c.createBehaviorUsage("RadBehaviorImported"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByGoogNamespace().Get("a"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_googRequireDefaultGoogModuleExportInLoadModule()
func (c *PolymerBehaviorExtractorTest) TestBehavior_googRequireDefaultGoogModuleExportInLoadModule() {
	c.parseAndInitializeExtractor(
		`goog.loadModule(function(exports) {
  goog.module('rad.RadBehavior');
`+
			c.createBehaviorDefinition("RadBehaviorLocal")+
			`  exports = RadBehaviorLocal;
  return exports;
});
`,
		`goog.module('a');
const RadBehaviorImported = goog.require('rad.RadBehavior');
`+
			c.createBehaviorUsage("RadBehaviorImported"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByGoogNamespace().Get("a"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_defaultGoogRequire_ofNamedExport()
func (c *PolymerBehaviorExtractorTest) TestBehavior_defaultGoogRequire_ofNamedExport() {
	c.parseAndInitializeExtractor(
		`goog.module('rad');
`+
			c.createBehaviorDefinition("RadBehaviorLocal")+
			`exports = {RadBehavior: RadBehaviorLocal};
`,
		`goog.module('a');
const rad = goog.require('rad');
`+
			c.createBehaviorUsage("rad.RadBehavior"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByGoogNamespace().Get("a"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_destructuringGoogRequire()
func (c *PolymerBehaviorExtractorTest) TestBehavior_destructuringGoogRequire() {
	c.parseAndInitializeExtractor(
		`goog.module('rad');
`+
			c.createBehaviorDefinition("RadBehaviorLocal")+
			`exports = {RadBehavior: RadBehaviorLocal};
`,
		`goog.module('a');
const {RadBehavior: RadBehaviorImported} = goog.require('rad');
`+
			c.createBehaviorUsage("RadBehaviorImported"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByGoogNamespace().Get("a"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_propertyOnDefaultExportOfGoogModule()
func (c *PolymerBehaviorExtractorTest) TestBehavior_propertyOnDefaultExportOfGoogModule() {
	c.T().Skip("@Ignore") // @Ignore
	c.parseAndInitializeExtractor(
		`goog.module('rad.SomeClass');
class SomeClass {}
`+
			c.createBehaviorDefinition("SomeClass.RadBehavior")+
			`exports = SomeClass;
`,
		`goog.module('a');
const SomeClass = goog.require('rad.SomeClass');
`+
			c.createBehaviorUsage("SomeClass.RadBehavior"))

	c.assertSingleBehaviorExtractionSucceedsWithMetadata(c.moduleMetadataMap.GetModulesByGoogNamespace().Get("a"))
}

// Java: PolymerBehaviorExtractorTest.testBehavior_googModuleGetOfModule_DirectUsage()
func (c *PolymerBehaviorExtractorTest) TestBehavior_googModuleGetOfModule_DirectUsage() {
	c.parseAndInitializeExtractor(
		`goog.module('rad');
`+
			c.createBehaviorDefinition("RadBehaviorLocal")+
			`exports = RadBehaviorLocal;
`,
		`(function() {
`+
			c.createBehaviorUsage("goog.module.get('rad')")+
			`})();
`)

	defs := c.extractor.ExtractBehaviors(c.behaviorArray, nil)
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	AssertThat(defs).HasSize(1)
}

// Java: PolymerBehaviorExtractorTest.testBehavior_googModuleGetOfProvide_DirectUsage()
func (c *PolymerBehaviorExtractorTest) TestBehavior_googModuleGetOfProvide_DirectUsage() {
	c.parseAndInitializeExtractor(
		`goog.provide('rad');
`+
			c.createBehaviorDefinition("rad"),
		`(function() {
`+
			c.createBehaviorUsage("goog.module.get('rad')")+
			`})();
`)

	defs := c.extractor.ExtractBehaviors(c.behaviorArray, nil)
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	AssertThat(defs).HasSize(1)
}

// Java: PolymerBehaviorExtractorTest.testBehavior_googModuleGetAssignedToVariable()
func (c *PolymerBehaviorExtractorTest) TestBehavior_googModuleGetAssignedToVariable() {
	c.T().Skip("@Ignore") // @Ignore
	// Note: this test fails because, currently, we only support resolving variables that are either
	// global or module-scoped. In theory the PolymerPass could be extended to understand other
	// locals but it hasn't been necessary thus far.
	c.parseAndInitializeExtractor(
		`goog.module('rad');
`+
			c.createBehaviorDefinition("RadBehaviorLocal")+
			`exports = RadBehaviorLocal;
`,
		`(function() {
const rad = goog.module.get('rad');
`+
			c.createBehaviorUsage("rad")+
			`})();
`)

	defs := c.extractor.ExtractBehaviors(c.behaviorArray, nil)
	AssertThat(c.compiler.GetErrors()).HasSize(1)
	AssertThat(defs).IsEmpty()
}

// Java: PolymerBehaviorExtractorTest.testBehavior_invalidGoogModuleGet()
func (c *PolymerBehaviorExtractorTest) TestBehavior_invalidGoogModuleGet() {
	// There is no module named 'rud'.
	c.parseAndInitializeExtractor(
		`(function() {
` +
			c.createBehaviorUsage("goog.module.get('rud')") +
			`})();
`)

	defs := c.extractor.ExtractBehaviors(c.behaviorArray, nil)
	AssertThat(c.compiler.GetErrors()).HasSize(1)
	AssertThat(defs).IsEmpty()
}

// Java: PolymerBehaviorExtractorTest.createBehaviorDefinition(String)
func (c *PolymerBehaviorExtractorTest) createBehaviorDefinition(name string) string {
	// Either `qualified.name` or `var simpleName`
	lhs := "var " + name
	if strings.IndexByte(name, '.') != -1 {
		lhs = name
	}
	return strings.ReplaceAll(`/** @polymerBehavior */
LHS = {
  properties: {
    howRad: Number
  },
  /** @param {number} radAmount */
  doSomethingRad:
    function(radAmount) { alert('Something ' + radAmount + 'rad!'); },
  /** @override */
  ready: function() {}
};
`, "LHS", lhs)
}

// Java: PolymerBehaviorExtractorTest.createBehaviorUsage(String)
func (c *PolymerBehaviorExtractorTest) createBehaviorUsage(behaviorName string) string {
	return strings.ReplaceAll(`var A = Polymer({
  is: 'x-element',
  behaviors: [ BEHAVIOR_NAME ],
});
`, "BEHAVIOR_NAME", behaviorName)
}

// assertSingleBehaviorExtractionSucceeds tests that the behaviorArray
// resolves to exactly one behavior without error
//
// Java: PolymerBehaviorExtractorTest.assertSingleBehaviorExtractionSucceeds()
func (c *PolymerBehaviorExtractorTest) assertSingleBehaviorExtractionSucceeds() {
	c.assertSingleBehaviorExtractionSucceedsWithMetadata( /* moduleMetadata */ nil)
}

// assertSingleBehaviorExtractionSucceedsWithMetadata tests that the
// behaviorArray resolves to exactly one behavior, using the provided
// ModuleMetadata for behavior extraction.
//
// Java: PolymerBehaviorExtractorTest.assertSingleBehaviorExtractionSucceeds(ModuleMetadata)
func (c *PolymerBehaviorExtractorTest) assertSingleBehaviorExtractionSucceedsWithMetadata(metadata *ModuleMetadataMap_ModuleMetadata) {
	defs := c.extractor.ExtractBehaviors(c.behaviorArray, metadata)
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	AssertThat(defs).HasSize(1)
}

// Java: PolymerBehaviorExtractorTest.testInlineLiteralBehavior()
func (c *PolymerBehaviorExtractorTest) TestInlineLiteralBehavior() {
	c.parseAndInitializeExtractor(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
  /** @param {string} funAmount */
  doSomethingFun: function(funAmount) { alert('Something ' + funAmount + ' fun!'); },
  /** @override */
  created: function() {}
};
/** @polymerBehavior */
var SuperCoolBehaviors = [FunBehavior, {
  properties: {
    howRad: Number
  },
  /** @param {number} radAmount */
  doSomethingRad: function(radAmount) { alert('Something ' + radAmount + ' rad!'); },
  /** @override */
  ready: function() {}
}];
var A = Polymer({
  is: 'x-element',
  behaviors: [ SuperCoolBehaviors ],
});
`)

	defs := c.extractor.ExtractBehaviors(c.behaviorArray, nil)
	AssertThat(defs).HasSize(2)

	// TODO(jlklein): Actually verify the properties of the BehaviorDefinitions.
}

// Java: PolymerBehaviorExtractorTest.testIsPropInBehavior()
func (c *PolymerBehaviorExtractorTest) TestIsPropInBehavior() {
	c.parseAndInitializeExtractor(
		`/** @polymerBehavior */
var FunBehavior = {
  is: 'fun-behavior',

  properties: {
    isFun: Boolean
  },
};
var A = Polymer({
  is: 'x-element',
  behaviors: [ FunBehavior ],
});
`)
	c.extractor.ExtractBehaviors(c.behaviorArray, nil)

	AssertThat(c.compiler.GetErrors()).HasSize(1)
	AssertError(c.compiler.GetErrors()[0]).HasType(PolymerPassErrors_POLYMER_INVALID_BEHAVIOR)
}

// Java: PolymerBehaviorExtractorTest.parseAndInitializeExtractor(String...)
func (c *PolymerBehaviorExtractorTest) parseAndInitializeExtractor(code ...string) {
	root := c.compiler.ParseTestCodeOfList(append([]string(nil), code...))

	NewGatherModuleMetadata(c.compiler, false, ModuleLoader_ResolutionMode_BROWSER).Process(rhino.IR_Root(), root)
	c.moduleMetadataMap = c.compiler.GetModuleMetadataMap()
	NewModuleMapCreator(c.compiler, c.moduleMetadataMap).Process(rhino.IR_Root(), root)
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	AssertThat(c.compiler.GetWarnings()).IsEmpty()
	globalNamespace := NewGlobalNamespace(c.compiler, root)
	c.extractor =
		NewPolymerBehaviorExtractor(
			c.compiler, globalNamespace, c.compiler.GetModuleMetadataMap(), c.compiler.GetModuleMap())

	NodeUtil_VisitPostOrder(
		root,
		NodeUtil_VisitorFunc(func(node rhino.Node) {
			if c.isBehaviorArrayDeclaration(node) {
				c.behaviorArray = node
			}
		}))

	AssertThat(c.behaviorArray).IsNotNull()
}

// Java: PolymerBehaviorExtractorTest.isBehaviorArrayDeclaration(Node)
func (c *PolymerBehaviorExtractorTest) isBehaviorArrayDeclaration(node rhino.Node) bool {
	return node.IsArrayLit() &&
		node.GetParent().IsStringKey() &&
		node.GetParent().GetString() == "behaviors"
}

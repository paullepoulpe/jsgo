// Port of com/google/javascript/jscomp/PolymerPassTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/PolymerPassTest.java)

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var polymerPassTest_EXTERNS_PREFIX = `/** @constructor */
var Element = function() {};
/** @constructor @extends {Element} */
var HTMLElement = function() {};
/** @constructor @extends {HTMLElement} */
var HTMLInputElement = function() {};
/** @constructor @extends {HTMLElement} */
var PolymerElement = function() {};
`

var polymerPassTest_EXTERNS_SUFFIX = `/**
 * @typedef {{
 *   type: !Function,
 *   value: (* | undefined),
 *   readOnly: (boolean | undefined),
 *   computed: (string | undefined),
 *   reflectToAttribute: (boolean | undefined),
 *   notify: (boolean | undefined),
 *   observer: (string | function(this:?, ?, ?) | undefined)
 * }}
 */
let PolymerElementPropertiesMeta;

/**
 * @typedef {Object<string, !Function|!PolymerElementPropertiesMeta>}
 */
let PolymerElementProperties;
/** @type {!Object} */
PolymerElement.prototype.$;
PolymerElement.prototype.created = function() {};
PolymerElement.prototype.ready = function() {};
PolymerElement.prototype.attached = function() {};
PolymerElement.prototype.domReady = function() {};
PolymerElement.prototype.detached = function() {};
/**
 * Call the callback after a timeout. Calling job again with the same name
 * resets the timer but will not result in additional calls to callback.
 *
 * @param {string} name
 * @param {Function} callback
 * @param {number} timeoutMillis The minimum delay in milliseconds before
 *     calling the callback.
 */
PolymerElement.prototype.job = function(name, callback, timeoutMillis) {};
/**
 * @param a {!Object}
 * @return {!Function}
 */
var Polymer = function(a) {};
// In the actual code Polymer.Element comes from a mixin. This definition is simplified.
/**
 * @polymer
 */
Polymer.Element = class {};

/** @type {!Object<string, !Element>} */
Polymer.Element.prototype.$;
var alert = function(msg) {};
`

var polymerPassTest_INPUT_EXTERNS = polymerPassTest_EXTERNS_PREFIX +
	`/** @constructor @extends {HTMLInputElement} */
var PolymerInputElement = function() {};
/** @type {!Object} */
PolymerInputElement.prototype.$;
PolymerInputElement.prototype.created = function() {};
PolymerInputElement.prototype.ready = function() {};
PolymerInputElement.prototype.attached = function() {};
PolymerInputElement.prototype.domReady = function() {};
PolymerInputElement.prototype.detached = function() {};
/**
 * Call the callback after a timeout. Calling job again with the same name
 * resets the timer but will not result in additional calls to callback.
 *
 * @param {string} name
 * @param {Function} callback
 * @param {number} timeoutMillis The minimum delay in milliseconds before
 *     calling the callback.
 */
PolymerInputElement.prototype.job = function(name, callback, timeoutMillis) {};
` +
	polymerPassTest_EXTERNS_SUFFIX

var polymerPassTest_REFLECT_OBJECT_DEF = `/** @const */ var $jscomp = $jscomp || {};
/** @const */ $jscomp.scope = {};
/**
 * @param {?Object} type
 * @param {T} object
 * @return {T}
 * @template T
 */
$jscomp.reflectObject = function (type, object) {};
`

var polymerPassTest_EXTERNS = polymerPassTest_EXTERNS_PREFIX + polymerPassTest_EXTERNS_SUFFIX + polymerPassTest_REFLECT_OBJECT_DEF

// PolymerPassTest is the unit tests for PolymerPass
type PolymerPassTest struct {
	CompilerTestCaseBase
}

// Java: PolymerPassTest()
func NewPolymerPassTest() *PolymerPassTest {
	c := &PolymerPassTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, polymerPassTest_EXTERNS)
	return c
}

func TestPolymerPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPolymerPassTest() })
}

// Java: PolymerPassTest.getProcessor(Compiler)
func (c *PolymerPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPolymerPass(compiler)
}

// Java: PolymerPassTest.getOptions()
func (c *PolymerPassTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetPreventLibraryInjection(true)
	return options
}

// Java: PolymerPassTest.setUp()
func (c *PolymerPassTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.AllowExternsChanges()
	c.EnableRunTypeCheckAfterProcessing()
	c.EnableParseTypeInfo()
	c.EnableCreateModuleMap()

	// Note: we don't have a version of testExternChanges that calls getActualExpected() here and
	// require test cases to hardcode the
	// full "Interface$m123..456$0" id instead. This is because "Interface$m123..456$0" gets put
	// in the synthetic externs, but named based on an input source file path, so getActualExpected
	// doesn't know what file path to use.
	// TODO(lharker): it would be a minor readability improvement to add a different helper to
	// support this case.

	// Helper to change the generic name string "Interface$UID$0" in the expected code with
	// actual string "Interface$m123..456$0" that will get produced
	c.SetGenericNameReplacements(polymerPassTest_genericNameReplacements())
}

// Java: PolymerPassTest.testPolymerRewriterGeneratesDeclarationOutsideLoadModule()
func (c *PolymerPassTest) TestPolymerRewriterGeneratesDeclarationOutsideLoadModule() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.loadModule(function(exports) {
  goog.module('ytu.app.ui.shared.YtuIcon');
  YtuIcon = Polymer({is: 'ytu-icon' });
  exports = YtuIcon;
  return exports;
})
`),
		c.Expected(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.loadModule(function(exports) {
goog.module('ytu.app.ui.shared.YtuIcon');
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerYtuIconInterface$UID$0}
 */
  var YtuIcon = function(){}
  YtuIcon = Polymer(/** @lends {YtuIcon.prototype} */ {is:"ytu-icon"});
  exports = YtuIcon;
  return exports;
})
`))
}

// Java: PolymerPassTest.testPolymerRewriterGeneratesDeclarationOutsideLoadModule2()
func (c *PolymerPassTest) TestPolymerRewriterGeneratesDeclarationOutsideLoadModule2() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.loadModule(function(exports) {
  goog.module('ytu.app.ui.shared.YtuIcon');
  Polymer({is: 'ytu-icon' });
  return exports;
})
`),
		c.Expected(
			`/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerYtuIconElementInterface$m1176578414$0}
 */
  var YtuIconElement = function(){}
`+
				testExternsBuilder_GetClosureExternsAsSource(),
			`goog.loadModule(function(exports) {
goog.module('ytu.app.ui.shared.YtuIcon');
  Polymer(/** @lends {YtuIconElement.prototype} */ {is:"ytu-icon"});
  return exports;
})
`))
}

// Java: PolymerPassTest.testAssignToGoogModuleExports()
func (c *PolymerPassTest) TestAssignToGoogModuleExports() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('modOne');
exports = Polymer({
  is: 'x-element',
});
`),
		c.Expected(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('modOne');
/** @constructor @extends {PolymerElement}
 * @implements {PolymerexportsForPolymer$jscomp0Interface$UID$0}
 */
var exportsForPolymer$jscomp0 = function() {};
exportsForPolymer$jscomp0 = Polymer(
  /** @lends {exportsForPolymer$jscomp0.prototype} */ {
  is: 'x-element',
});
exports = exportsForPolymer$jscomp0;
`))
}

// Java: PolymerPassTest.testAssignToGoogLoadModuleExports()
func (c *PolymerPassTest) TestAssignToGoogLoadModuleExports() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.loadModule(function(exports) {
goog.module('modOne');
exports = Polymer({
  is: 'x-element',
});
});
`),
		c.Expected(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.loadModule(function(exports) {
goog.module('modOne');
/** @constructor @extends {PolymerElement}
 * @implements {PolymerexportsForPolymer$jscomp0Interface$UID$0}
 */
var exportsForPolymer$jscomp0 = function() {};
exportsForPolymer$jscomp0 = Polymer(
  /** @lends {exportsForPolymer$jscomp0.prototype} */ {
  is: 'x-element',
});
exports = exportsForPolymer$jscomp0;
});
`))
}

// Java: PolymerPassTest.testSameNamedPolymerClassesInTwoModules()
func (c *PolymerPassTest) TestSameNamedPolymerClassesInTwoModules() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('modOne');
const Button = Polymer({
  is: 'x-element',
});
`,
			`goog.module('modTwo');
const Button = Polymer({
  is: 'y-element',
});
`),
		c.Expected(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('modOne');
/** @constructor @extends {PolymerElement} @implements {PolymerButtonInterface$UID$0} */
var Button = function() {};
Button = Polymer(/** @lends {Button.prototype} */ {
  is: 'x-element',
});
`,
			`goog.module('modTwo');
/** @constructor @extends {PolymerElement} @implements {PolymerButtonInterface$UID$0} */
var Button = function() {};
Button = Polymer(/** @lends {Button.prototype} */ {
  is: 'y-element',
});
`))
}

// Java: PolymerPassTest.testPolymerRewriterGeneratesDeclaration_OutsideGoogModule()
func (c *PolymerPassTest) TestPolymerRewriterGeneratesDeclaration_OutsideGoogModule() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
var X = Polymer({
  is: 'x-element',
});
`),
		c.Expected(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0} */
var X = function() {};
X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',
});
`))
}

// Java: PolymerPassTest.testPolymerRewriterGeneratesDeclaration_OutsideGoogModule2()
func (c *PolymerPassTest) TestPolymerRewriterGeneratesDeclaration_OutsideGoogModule2() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
Polymer({
  is: 'x-element',
});
`),
		c.Expected(
			`/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerXElementElementInterface$m1176578414$0}
 */
var XElementElement = function() {};
`+
				testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
Polymer(/** @lends {XElementElement.prototype} */ {
  is: 'x-element',
});
`))
}

// Java: PolymerPassTest.testGeneratesCodeAfterGoogRequire_WithLegacyNamespace()
func (c *PolymerPassTest) TestGeneratesCodeAfterGoogRequire_WithLegacyNamespace() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource()+
				`goog.provide('a.UnpluggedCancelOfferRenderer');
goog.provide('a.b');
/** @constructor */
a.UnpluggedCancelOfferRenderer = function() {
};
`,
			`goog.module('a.mod');
exports.Property = 'a property';
`,
			`goog.module('mod');
goog.module.declareLegacyNamespace()
const CancelOfferRenderer = goog.require('a.UnpluggedCancelOfferRenderer');
const {Property} = goog.require('a.mod');
Polymer(
 {is:'ytu-cancel-offer', properties:{
/** @type {CancelOfferRenderer} */
data:Object}});
`),
		c.Expected(
			strings.ReplaceAll(`/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerYtuCancelOfferElementInterface$m1176578413$0}
 */
var YtuCancelOfferElement = function() {};
CLOSURE_EXTERNS
goog.provide('a.UnpluggedCancelOfferRenderer');
goog.provide('a.b');
/** @constructor */ a.UnpluggedCancelOfferRenderer = function() {};
`, "CLOSURE_EXTERNS", testExternsBuilder_GetClosureExternsAsSource()),
			`goog.module('a.mod');
exports.Property = 'a property';
`,
			`goog.module('mod');
goog.module.declareLegacyNamespace()
const CancelOfferRenderer = goog.require('a.UnpluggedCancelOfferRenderer');
const {Property} = goog.require('a.mod');
/** @type {CancelOfferRenderer} */ YtuCancelOfferElement.prototype.data;
Polymer(/** @lends {YtuCancelOfferElement.prototype} */ {
  is:'ytu-cancel-offer',
  properties: $jscomp.reflectObject(YtuCancelOfferElement, {data:Object})
});
`))
}

// Java: PolymerPassTest.testGeneratesCodeAfterGoogRequire_WithSetTestOnly()
func (c *PolymerPassTest) TestGeneratesCodeAfterGoogRequire_WithSetTestOnly() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource()+
				`goog.provide('a.UnpluggedCancelOfferRenderer');
/** @constructor */
a.UnpluggedCancelOfferRenderer = function() {
};
`,
			`goog.module('mod');
goog.setTestOnly()
const CancelOfferRenderer = goog.require('a.UnpluggedCancelOfferRenderer');
Polymer(
 {is:'ytu-cancel-offer', properties:{
   /** @type {CancelOfferRenderer} */
   data:Object}});
`),
		c.Expected(
			strings.ReplaceAll(`/** @constructor @extends {PolymerElement} @implements
 {PolymerYtuCancelOfferElementInterface$m1176578414$0} */
var YtuCancelOfferElement = function() {};
CLOSURE_EXTERNS
goog.provide('a.UnpluggedCancelOfferRenderer');
/** @constructor */ a.UnpluggedCancelOfferRenderer = function() {};
`, "CLOSURE_EXTERNS", testExternsBuilder_GetClosureExternsAsSource()),
			`goog.module('mod');
goog.setTestOnly()
const CancelOfferRenderer = goog.require('a.UnpluggedCancelOfferRenderer');
/** @type {CancelOfferRenderer} */
YtuCancelOfferElement.prototype.data;
Polymer(
/** @lends {YtuCancelOfferElement.prototype} */
 {
    is:'ytu-cancel-offer',
    properties: $jscomp.reflectObject(YtuCancelOfferElement, {data:Object})
});
`))
}

// Java: PolymerPassTest.testPolymerRewriterGeneratesDeclaration_OutsideES6Module()
func (c *PolymerPassTest) TestPolymerRewriterGeneratesDeclaration_OutsideES6Module() {
	c.Test(
		c.Srcs(
			"", // empty script for getNodeForCodeInsertion
			`var X = Polymer({
  is: 'x-element',
});
export {X};
`),
		c.Expected(
			"",
			`/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0} */
var X = function() {};
X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',})
export {X};
`))
}

// Java: PolymerPassTest.testPolymerRewriterGeneratesDeclaration_OutsideES6Module2()
func (c *PolymerPassTest) TestPolymerRewriterGeneratesDeclaration_OutsideES6Module2() {
	c.Test(
		c.Srcs(
			"", // empty script for getNodeForCodeInsertion
			`var PI = 3.14
Polymer({
  is: 'x-element',
});
export {PI};
`),
		c.Expected(
			`/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerXElementElementInterface$m1176578414$0}
 */
var XElementElement = function() {};
`,
			`var PI = 3.14
Polymer(/** @lends {XElementElement.prototype} */ {
  is: 'x-element',})
export {PI};
`))
}

// Java: PolymerPassTest.testPolymerRewriterGeneratesDeclaration_OutsideModule_IIFE()
func (c *PolymerPassTest) TestPolymerRewriterGeneratesDeclaration_OutsideModule_IIFE() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
(function() {
  var X = Polymer({
    is: 'x-element',
  });
})();
`),
		c.Expected(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
(function() {
/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0} */
var X = function() {};
    X = Polymer(/** @lends {X.prototype} */ {
    is: 'x-element',
  });
})();
`))
}

// Java: PolymerPassTest.testPolymerRewriterGeneratesDeclaration_OutsideModule_IIFE2()
func (c *PolymerPassTest) TestPolymerRewriterGeneratesDeclaration_OutsideModule_IIFE2() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
(function() {
  Polymer({
    is: 'x-element',
  });
})();
`),
		c.Expected(
			`/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerXElementElementInterface$m1176578414$0}
 */
var XElementElement = function() {};
`+
				testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
(function() {
    Polymer(/** @lends {XElementElement.prototype} */ {
    is: 'x-element',
  });
})();
`))
}

// Java: PolymerPassTest.testPolymerRewriterGeneratesDeclaration_OutsideModule_WithRequires()
func (c *PolymerPassTest) TestPolymerRewriterGeneratesDeclaration_OutsideModule_WithRequires() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
const Component = goog.require('a');
goog.forwardDeclare('something.else');
const someLocal = (function() { return 0; })();
var X = Polymer({
  is: 'x-element',
});
`),
		c.Expected(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
const Component = goog.require('a');
goog.forwardDeclare('something.else');
const someLocal = (function() { return 0; })();
/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0} */
var X = function() {};
X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',
});
`))
}

// Java: PolymerPassTest.testPolymerRewriterGeneratesDeclaration_OutsideModule_WithRequires2()
func (c *PolymerPassTest) TestPolymerRewriterGeneratesDeclaration_OutsideModule_WithRequires2() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
const Component = goog.require('a');
goog.forwardDeclare('something.else');
const someLocal = (function() { return 0; })();
Polymer({
  is: 'x-element',
});
`),
		c.Expected(
			`/** @constructor @extends {PolymerElement} @implements {PolymerXElementElementInterface$m1176578414$0} */
var XElementElement = function() {};
`+
				testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('mod');
const Component = goog.require('a');
goog.forwardDeclare('something.else');
const someLocal = (function() { return 0; })();
Polymer(/** @lends {XElementElement.prototype} */ {
  is: 'x-element',
});
`))
}

// Java: PolymerPassTest.testVarTarget()
func (c *PolymerPassTest) TestVarTarget() {
	c.Test(
		`var X = Polymer({
  is: 'x-element',
});
`,
		`/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0} */
var X = function() {};
X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',
});
`)

	c.Test(
		`var X = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() { return {}; }
};
`,
		`/** @implements {PolymerXInterface$UID$0} */
var X = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() { return $jscomp.reflectObject(X, {}); }
};
`)
}

// Java: PolymerPassTest.testVarTargetMissingExterns()
func (c *PolymerPassTest) TestVarTargetMissingExterns() {
	c.AllowSourcelessWarnings() // the missing Polymer externs warning has no source, since it's not
	// about any particular file
	c.TestError(
		c.Externs(""),
		c.Srcs(
			`var X = Polymer({
  is: 'x-element',
});
`),
		c.Error(PolymerPassErrors_POLYMER_MISSING_EXTERNS))
}

// Java: PolymerPassTest.testLetTarget()
func (c *PolymerPassTest) TestLetTarget() {
	c.Test(
		`let X = Polymer({
  is: 'x-element',
});
`,
		`/**
 * @constructor
 * @implements {PolymerXInterface$UID$0}
 * @extends {PolymerElement}
 */
var X = function() {};
X = Polymer(/** @lends {X.prototype} */ {is:'x-element'});
`)

	c.Test(
		`let X = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() { return { }; }
};
`,
		`/** @implements {PolymerXInterface$UID$0} */
let X = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() { return $jscomp.reflectObject(X, {}); }
};
`)
}

// Java: PolymerPassTest.testConstTarget()
func (c *PolymerPassTest) TestConstTarget() {
	c.Test(
		`const X = Polymer({
  is: 'x-element',
});
`,
		`/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerXInterface$UID$0}
 */
var X=function(){};X=Polymer(/** @lends {X.prototype} */ {is:"x-element"})
`)

	c.Test(
		`const X = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() { return {}; }
};
`,
		`/** @implements {PolymerXInterface$UID$0} */
const X = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() { return $jscomp.reflectObject(X, {}); }
};
`)
}

// Java: PolymerPassTest.testDefaultTypeNameTarget()
func (c *PolymerPassTest) TestDefaultTypeNameTarget() {
	c.Test(
		`Polymer({
  is: 'x',
});
`,
		`/**
 * @implements {PolymerXElementInterface$UID$0}
 * @constructor @extends {PolymerElement}
 */
var XElement = function() {};
Polymer(/** @lends {XElement.prototype} */ {
  is: 'x',
});
`)
}

// Java: PolymerPassTest.testPathAssignmentTarget()
func (c *PolymerPassTest) TestPathAssignmentTarget() {
	c.Test(
		`/** @const */ var x = {};
x.Z = Polymer({
  is: 'x-element',
});
`,
		`/** @const */ var x = {};
/** @constructor @extends {PolymerElement} @implements {Polymerx_ZInterface$UID$0} */
x.Z = function() {};
x.Z = Polymer(/** @lends {x.Z.prototype} */ {
  is: 'x-element',
});
`)

	c.Test(
		`const x = {};
x.Z = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() { return {}; }
};
`,
		`const x = {};
/** @implements {Polymerx_ZInterface$UID$0} */
x.Z = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() { return $jscomp.reflectObject(x.Z, {}); }
};
`)

	c.Test(
		`var x = {};
x.Z = Polymer({
  is: 'x-element',
});
`,
		`var x = {};
/** @constructor @extends {PolymerElement} @implements {Polymerx_ZInterface$UID$0} */
x.Z = function() {};
x.Z = Polymer(/** @lends {x.Z.prototype} */ {
  is: 'x-element',
});
`)
}

// Java: PolymerPassTest.testComputedPropName()
func (c *PolymerPassTest) TestComputedPropName() {
	// TypeCheck cannot grab a name from a complicated computedPropName
	c.Test(
		"var X = Polymer({is:'x-element', [name + (() => 42)]: function() {return 42;}});",
		`/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0} */
var X = function() {}

X = Polymer(/** @lends {X.prototype} */{
  is: 'x-element',
  /** @this {X} */
  [name + (()=>42)]: function() {return 42;},
});
`)
}

// Java: PolymerPassTest.testIIFEExtractionInGlobalNamespace()
func (c *PolymerPassTest) TestIIFEExtractionInGlobalNamespace() {
	c.Test(
		`/** @const */ var x = {};
(function() {
  x.Z = Polymer({
    is: 'x-element',
    sayHi: function() { alert('hi'); },
  });
})()
`,
		`/** @const */ var x = {};
(function() {
  /** @constructor @extends {PolymerElement} @implements {Polymerx_ZInterface$UID$0}*/
  x.Z = function() {};
  x.Z = Polymer(/** @lends {x.Z.prototype} */ {
    is: 'x-element',
    /** @this {x.Z} */
    sayHi: function() { alert('hi'); },
  });
  /** @export */
  x.Z.prototype.sayHi;
})()
`)

	c.Test(
		`const x = {};
(function() {
  x.Z = class extends Polymer.Element {
    static get is() { return 'x-element'; }
    static get properties() { return {}; }
  };
})();
`,
		`const x = {};
(function() {
  /** @implements {Polymerx_ZInterface$UID$0} */
  x.Z = class extends Polymer.Element {
    /** @return {string} */
    static get is() { return 'x-element'; }
    /** @return {PolymerElementProperties} */
    static get properties() { return $jscomp.reflectObject(x.Z, {}); }
  };
})();
`)
}

// Java: PolymerPassTest.testIIFEExtractionNoAssignmentTarget()
func (c *PolymerPassTest) TestIIFEExtractionNoAssignmentTarget() {
	c.Test(
		`(function() {
  Polymer({
    is: 'x',
    sayHi: function() { alert('hi'); },
  });
})()
`,
		`/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerXElementInterface$UID$0}
 */
var XElement = function() {};
(function() {
  Polymer(/** @lends {XElement.prototype} */ {
    is: 'x',
    /** @this {XElement} */
    sayHi: function() { alert('hi'); },
  });
  /** @export */
  XElement.prototype.sayHi;
})()
`)

	c.Test(
		`(function() {
  const X = class extends Polymer.Element {
    static get is() { return 'x-element'; }
    static get properties() { return {}; }
  };
})();
`,
		`(function() {
  /** @implements {PolymerXInterface$UID$0} */
  const X = class extends Polymer.Element {
    /** @return {string} */
    static get is() { return 'x-element'; }
    /** @return {PolymerElementProperties} */
    static get properties() { return $jscomp.reflectObject(X, {}); }
  };
})();
`)
}

// Java: PolymerPassTest.testIIFEExtractionVarTarget()
func (c *PolymerPassTest) TestIIFEExtractionVarTarget() {
	c.Test(
		`(function() {
  Polymer({
    is: 'foo-thing',
    sayHi: function() { alert('hi'); },
  });
})()
`,
		`/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerFooThingElementInterface$UID$0}
 */
var FooThingElement = function() {};
(function() {
  Polymer(/** @lends {FooThingElement.prototype} */ {
    is: 'foo-thing',
    /** @this {FooThingElement} */
    sayHi: function() { alert('hi'); },
  });
  /** @export */
  FooThingElement.prototype.sayHi;
})()
`)
}

// Java: PolymerPassTest.testConstructorExtraction()
func (c *PolymerPassTest) TestConstructorExtraction() {
	c.Test(
		`var X = Polymer({
  is: 'x-element',
  /**
   * @param {string} name
   */
  factoryImpl: function(name) { alert('hi, ' + name); },
});
`,
		`/**
 * @param {string} name
 * @constructor @extends {PolymerElement}
 * @implements {PolymerXInterface$UID$0}
 */
var X = function(name) { alert('hi, ' + name); };
X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',
  factoryImpl: function(name) { alert('hi, ' + name); },
});
/** @export */
X.prototype.factoryImpl;
`)
}

// Java: PolymerPassTest.testShorthandConstructorExtraction()
func (c *PolymerPassTest) TestShorthandConstructorExtraction() {
	c.Test(
		`var X = Polymer({
  is: 'x-element',
  /**
   * @param {string} name
   */
  factoryImpl(name) { alert('hi, ' + name); },
});
`,
		`/**
 * @param {string} name
 * @constructor @extends {PolymerElement}
 * @implements {PolymerXInterface$UID$0}
 */
var X = function(name) { alert('hi, ' + name); };

X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',
  factoryImpl(name) { alert('hi, ' + name); },
});
/** @export */
X.prototype.factoryImpl;
`)
}

// Java: PolymerPassTest.testOtherKeysIgnored()
func (c *PolymerPassTest) TestOtherKeysIgnored() {
	c.Test(
		`var X = Polymer({
  is: 'x-element',
  listeners: {
    'click': 'handleClick',
  },

  handleClick: function(e) {
    alert('Thank you for clicking');
  },
});
`,
		`/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0}*/
var X = function() {};
X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',
  listeners: {
    'click': 'handleClick',
  },

  /** @this {X} */
  handleClick: function(e) {
    alert('Thank you for clicking');
  },
});
/** @export */
X.prototype.handleClick;
`)

	c.Test(
		`class X extends Polymer.Element {
  static get is() { return 'x-element'; }
  handleClick(e) {
    alert('Thank you for clicking');
  }
}
`,
		`/** @implements {PolymerXInterface$UID$0} */
class X extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  handleClick(e) {
    alert('Thank you for clicking');
  }
}
/** @export */
X.prototype.handleClick;
`)
}

// Java: PolymerPassTest.testListenersAndHostAttributeKeysQuoted()
func (c *PolymerPassTest) TestListenersAndHostAttributeKeysQuoted() {
	c.Test(
		`var X = Polymer({
  is: 'x-element',
  listeners: {
    click: 'handleClick',
    'foo-bar': 'handleClick',
  },
  hostAttributes: {
    role: 'button',
    'foo-bar': 'done',
    blah: 1,
  },

  handleClick: function(e) {
    alert('Thank you for clicking');
  },
});
`,
		`/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0}*/
var X = function() {};
X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',
  listeners: {
    'click': 'handleClick',
    'foo-bar': 'handleClick',
  },
  hostAttributes: {
    'role': 'button',
    'foo-bar': 'done',
    'blah': 1,
  },

  /** @this {X} */
  handleClick: function(e) {
    alert('Thank you for clicking');
  },
});
/** @export */
X.prototype.handleClick;
`)
}

// Java: PolymerPassTest.testPolymerBehaviorListenersAndHostAttributeKeysQuoted()
func (c *PolymerPassTest) TestPolymerBehaviorListenersAndHostAttributeKeysQuoted() {
	// Polymer behaviors with string key property on listeners and hostAttributes need to be marked
	// as quoted.
	// This is the bugfix for b/388337323.
	c.Test(
		c.Srcs(
			`/** @polymerBehavior */
exports.CheckedOnTapBehavior = {
  listeners: {tap: 'onTap'},
  hostAttributes: {foo: 1},
  onTap() {
    this.checked = true;
  },
};
`),
		c.Expected(
			`/** @polymerBehavior @nocollapse */
exports.CheckedOnTapBehavior = {
  listeners: {'tap': 'onTap'},
  hostAttributes: {'foo': 1},
  /** @suppress {checkTypes|globalThis|visibility} */
  onTap() {
    this.checked = true;
  },
};
`))
}

// Java: PolymerPassTest.testNativeElementExtension()
func (c *PolymerPassTest) TestNativeElementExtension() {
	js := `Polymer({
  is: 'x-input',
  extends: 'input',
});
`

	c.Test(
		js,
		`/**
 * @constructor @extends {PolymerInputElement}
 * @implements {PolymerXInputElementInterface$UID$0}
 */
var XInputElement = function() {};
Polymer(/** @lends {XInputElement.prototype} */ {
  is: 'x-input',
  extends: 'input',
});
`)

	c.TestExternChanges(
		c.Externs(polymerPassTest_EXTERNS),
		c.Srcs(js),
		c.Expected(
			`/** @interface */
var PolymerXInputElementInterface$m1146332801$0 = function() {};
`,
			polymerPassTest_INPUT_EXTERNS+polymerPassTest_REFLECT_OBJECT_DEF))
}

// Java: PolymerPassTest.testExtendNonExistentElement()
func (c *PolymerPassTest) TestExtendNonExistentElement() {
	js := `Polymer({
  is: 'x-input',
  extends: 'nonexist',
});
`

	c.TestError(js, PolymerPassErrors_POLYMER_INVALID_EXTENDS)
}

// Java: PolymerPassTest.testNativeElementExtensionExternsNotDuplicated()
func (c *PolymerPassTest) TestNativeElementExtensionExternsNotDuplicated() {
	js := `Polymer({
  is: 'x-input',
  extends: 'input',
});
Polymer({
  is: 'y-input',
  extends: 'input',
});
`
	newExterns := c.Expected(
		`/** @interface */
var PolymerXInputElementInterface$m1146332801$0 = function() {};
/** @interface */
var PolymerYInputElementInterface$m1146332801$1 = function() {};
`,
		polymerPassTest_INPUT_EXTERNS+polymerPassTest_REFLECT_OBJECT_DEF)

	c.TestExternChanges(c.Externs(polymerPassTest_EXTERNS), c.Srcs(js), newExterns)
}

// Java: PolymerPassTest.testPropertiesAddedToPrototype()
func (c *PolymerPassTest) TestPropertiesAddedToPrototype() {
	c.Test(
		`/** @constructor */
var User = function() {};
/** @const */ var a = {};
a.B = Polymer({
  is: 'x-element',
  properties: {
    /** @type {!User} @private */
    user_: Object,
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
    thingToDo: Function,
  },
});
`,
		`/** @constructor */
var User = function() {};
/** @const */ var a = {};
/** @constructor @extends {PolymerElement} @implements {Polymera_BInterface$UID$0}*/
a.B = function() {};
/** @type {!User} @private */
a.B.prototype.user_;
/** @type {!Array} */
a.B.prototype.pets;
/** @type {string} */
a.B.prototype.name;
/** @type {!Function} */
a.B.prototype.thingToDo;
a.B = Polymer(/** @lends {a.B.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(a.B, {
    user_: Object,
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
    thingToDo: Function,
  }),
});
`)

	c.Test(
		`/** @constructor */
var User = function() {};
var a = {};
a.B = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function,
    };
  }
};
`,
		`/** @constructor */
var User = function() {};
var a = {};
/** @implements {Polymera_BInterface$UID$0} */
a.B = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(a.B, {
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function,
    });
  }
};
/** @type {!User} @private */
a.B.prototype.user_;
/** @type {!Array} */
a.B.prototype.pets;
/** @type {string} */
a.B.prototype.name;
/** @type {!Function} */
a.B.prototype.thingToDo;
`)
}

// Java: PolymerPassTest.testPropertiesDefaultValueFunctions()
func (c *PolymerPassTest) TestPropertiesDefaultValueFunctions() {
	c.Test(
		`/** @constructor */
var User = function() {};
/** @const */ var a = {};
a.B = Polymer({
  is: 'x-element',
  properties: {
    /** @type {!User} @private */
    user_: {
      type: Object,
      value: function() { return new User(); },
    },
    pets: {
      type: Array,
      notify: true,
      value: function() { return [this.name]; },
    },
    name: String,
  },
});
`,
		`/** @constructor */
var User = function() {};
/** @const */ var a = {};
/** @constructor @extends {PolymerElement} @implements {Polymera_BInterface$UID$0}*/
a.B = function() {};
/** @type {!User} @private */
a.B.prototype.user_;
/** @type {!Array} */
a.B.prototype.pets;
/** @type {string} */
a.B.prototype.name;
a.B = Polymer(/** @lends {a.B.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(a.B, {
    user_: {
      type: Object,
      /** @this {a.B} @return {!User} */
      value: function() { return new User(); },
    },
    pets: {
      type: Array,
      notify: true,
      /** @this {a.B} @return {!Array} */
      value: function() { return [this.name]; },
    },
    name: String,
  }),
});
`)

	c.Test(
		`/** @constructor */
var User = function() {};
/** @const */ var a = {};
a.B = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: {
        type: Object,
        value: function() { return new User(); },
      },
      pets: {
        type: Array,
        notify: true,
        value: function() { return [this.name]; },
      },
      name: String,
    };
  }
};
`,
		`/** @constructor */
var User = function() {};
/** @const */ var a = {};
/** @implements {Polymera_BInterface$UID$0} */
a.B = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(a.B, {
      user_: {
        type: Object,
        /** @this {a.B} @return {!User} */
        value: function() { return new User(); },
      },
      pets: {
        type: Array,
        notify: true,
        /** @this {a.B} @return {!Array} */
        value: function() { return [this.name]; },
      },
      name: String,
    });
  }
};
/** @type {!User} @private */
a.B.prototype.user_;
/** @type {!Array} */
a.B.prototype.pets;
/** @type {string} */
a.B.prototype.name;
`)
}

// Java: PolymerPassTest.testPropertiesDefaultValueShortHandFunction()
func (c *PolymerPassTest) TestPropertiesDefaultValueShortHandFunction() {
	c.Test(
		`/** @constructor */
var User = function() {};
var ES6Test = Polymer({
  is: 'x-element',
  properties: {
    user: {
      type: Object,
      value() { return new User();},
    },
  },
});
`,
		`/** @constructor */
var User = function() {};

/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerES6TestInterface$UID$0}
 */
var ES6Test = function() {};
/** @type {!Object} */
ES6Test.prototype.user;

ES6Test = Polymer(/** @lends {ES6Test.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(ES6Test, {
    user: {
      type: Object,
      /** @this {ES6Test} @return {!Object} */
      value() { return new User(); },
    },
  }),
});
`)

	c.Test(
		`/** @constructor */
var User = function() {};
var a = {};
a.B = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: {
        type: Object,
        value() { return new User(); },
      },
      pets: {
        type: Array,
        notify: true,
        value() { return [this.name]; },
      },
      name: String,
    };
  }
};
`,
		`/** @constructor */
var User = function() {};
var a = {};
/** @implements {Polymera_BInterface$UID$0} */
a.B = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(a.B, {
      user_: {
        type: Object,
        /** @this {a.B} @return {!User} */
        value() { return new User(); },
      },
      pets: {
        type: Array,
        notify: true,
        /** @this {a.B} @return {!Array} */
        value() { return [this.name]; },
      },
      name: String,
    });
  }
};
/** @type {!User} @private */
a.B.prototype.user_;
/** @type {!Array} */
a.B.prototype.pets;
/** @type {string} */
a.B.prototype.name;
`)
}

// Java: PolymerPassTest.testReadOnlyPropertySetters()
func (c *PolymerPassTest) TestReadOnlyPropertySetters() {
	js := `/** @const */ var a = {};
a.B = Polymer({
  is: 'x-element',
  properties: {
    /** @type {!Array<string>} */
    pets: {
      type: Array,
      readOnly: true,
    },
    /** @private */
    name_: String,
  },
});
`

	c.Test(
		js,
		`/** @const */ var a = {};
/** @constructor @extends {PolymerElement} @implements {Polymera_BInterface$UID$0} */
a.B = function() {};
/** @type {!Array<string>} */
a.B.prototype.pets;
/** @private {string} */
a.B.prototype.name_;
/*** @param {!Array<string>} pets @override */
a.B.prototype._setPets = function(pets) {};
a.B = Polymer(/** @lends {a.B.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(a.B, {
    pets: {
      type: Array,
      readOnly: true,
    },
    name_: String,
  }),
});
`)

	c.TestExternChanges(
		c.Externs(polymerPassTest_EXTERNS),
		c.Srcs(js),
		c.Expected(
			`/** @interface */
var Polymera_BInterface$m1146332801$0 = function() {};
/** @type {?} */
Polymera_BInterface$m1146332801$0.prototype.pets;
/** @private {?} */
Polymera_BInterface$m1146332801$0.prototype.name_;
/** @param {?} pets **/
Polymera_BInterface$m1146332801$0.prototype._setPets;
`,
			polymerPassTest_EXTERNS))

	jsClass := `class A extends Polymer.Element {
  static get is() { return 'a-element'; }
  static get properties() {
    return {
      pets: {
        type: Array,
        readOnly: true
      }
    };
  }
}
`

	c.TestExternChanges(
		c.Externs(polymerPassTest_EXTERNS),
		c.Srcs(jsClass),
		c.Expected(
			`/** @interface */
var PolymerAInterface$m1146332801$0 = function() {};
/** @type {?} */
PolymerAInterface$m1146332801$0.prototype.pets;
/** @param {?} pets **/
PolymerAInterface$m1146332801$0.prototype._setPets;
`,
			polymerPassTest_EXTERNS))

	c.Test(
		jsClass,
		`/** @implements {PolymerAInterface$UID$0} */
class A extends Polymer.Element {
  /** @return {string} */ static get is() { return 'a-element'; }
  /** @return {PolymerElementProperties} */ static get properties() {
    return $jscomp.reflectObject(A, {
      pets: {
        type: Array,
        readOnly: true
      }
    });
  }
}
/** @type {!Array} */
A.prototype.pets;
/*** @param {!Array} pets @override */
A.prototype._setPets = function(pets) {};
`)
}

// Java: PolymerPassTest.testPolymerPropertyNamesGeneratedInExterns()
func (c *PolymerPassTest) TestPolymerPropertyNamesGeneratedInExterns() {
	js := `(function() {
  Polymer({
    is: 'foo',
    properties: {
     /** @type {{randomProperty: string}} */
     value: Object
  }
  });
})();

const obj = {randomProperty: 0, otherProperty: 1};
`
	c.Test(
		js,
		`/**
* @constructor
* @extends {PolymerElement}
* @implements {PolymerFooElementInterface$UID$0}
*/
var FooElement=function(){};
(function(){
   /** @type {{randomProperty:string}} */
   FooElement.prototype.value;
   Polymer(/** @lends {FooElement.prototype} */ {
     is:"foo",
     properties: $jscomp.reflectObject(FooElement, {value:Object})
   })})();
const obj={randomProperty:0,otherProperty:1}
`)

	c.TestExternChanges(
		c.Externs(polymerPassTest_INPUT_EXTERNS),
		c.Srcs(js),
		c.Expected(
			`/** @interface */ var PolymerFooElementInterface$m1146332801$0=function(){};
/** @type {{randomProperty:?}} */ var PolymerDummyVar0;
/** @type {?} */ PolymerFooElementInterface$m1146332801$0.prototype.value
`,
			polymerPassTest_INPUT_EXTERNS))
}

// Java: PolymerPassTest.testReflectToAttributeProperties()
func (c *PolymerPassTest) TestReflectToAttributeProperties() {
	js := `/** @const */ var a = {};
a.B = Polymer({
  is: 'x-element',
  properties: {
    /** @type {!Array<string>} */
    pets: {
      type: Array,
      readOnly: false,
    },
    name: {
      type: String,
      reflectToAttribute: true
    }
  },
});
`

	c.Test(
		js,
		`/** @const */ var a = {};
/** @constructor @extends {PolymerElement} @implements {Polymera_BInterface$m1146332801$0} */
a.B = function() {};
/** @type {!Array<string>} */
a.B.prototype.pets;
/** @type {string} */
a.B.prototype.name;
a.B = Polymer(/** @lends {a.B.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(a.B, {
    pets: {
      type: Array,
      readOnly: false,
    },
    name: {
      type: String,
      reflectToAttribute: true
    }
  }),
});
`)

	c.TestExternChanges(
		c.Externs(polymerPassTest_EXTERNS),
		c.Srcs(js),
		c.Expected(
			`/** @interface */
var Polymera_BInterface$m1146332801$0 = function() {};
/** @type {?} */
Polymera_BInterface$m1146332801$0.prototype.pets;
/** @type {?} */
Polymera_BInterface$m1146332801$0.prototype.name;
`,
			polymerPassTest_EXTERNS))
}

// Java: PolymerPassTest.testPolymerClassObserversTyped()
func (c *PolymerPassTest) TestPolymerClassObserversTyped() {
	c.Test(
		`class FooElement extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get observers() {
    return [];
  }
}
`,
		`/** @implements {PolymerFooElementInterface$UID$0} */
class FooElement extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {!Array<string>} */
  static get observers() {
    return [];
  }
}
`)
}

// Java: PolymerPassTest.testShorthandFunctionDefinition()
func (c *PolymerPassTest) TestShorthandFunctionDefinition() {
	c.Test(
		`var ES6Test = Polymer({
  is: 'x-element',
  sayHi() {
    alert('hi');
  },
});
`,
		`/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerES6TestInterface$UID$0}
 */
var ES6Test = function() {};

ES6Test = Polymer(/** @lends {ES6Test.prototype} */ {
  is: 'x-element',
  /** @this {ES6Test} */
  sayHi() {
    alert('hi');
  },
});
/** @export */
ES6Test.prototype.sayHi;
`)
}

// Java: PolymerPassTest.testArrowFunctionDefinition()
func (c *PolymerPassTest) TestArrowFunctionDefinition() {
	c.Test(
		`var ES6Test = Polymer({
  is: 'x-element',
  sayHi: () => 42,
});
`,
		`/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerES6TestInterface$UID$0}
 */
var ES6Test = function() {};

ES6Test = Polymer(/** @lends {ES6Test.prototype} */ {
  is: 'x-element',
  /** @this {ES6Test} */
  sayHi: () => 42,
});
/** @export */
ES6Test.prototype.sayHi;
`)
}

// Java: PolymerPassTest.testArrowFunctionDefinitionInIifeBehavior()
func (c *PolymerPassTest) TestArrowFunctionDefinitionInIifeBehavior() {
	c.Test(
		c.Srcs(
			`var es6 = {};
(function() {
  class Foo {}
  /** @polymerBehavior */
  es6.Behavior = {
    sayHi: () => new Foo(),
  };
})();
es6.Test = Polymer({
  is: 'x-element',
  behaviors: [es6.Behavior]
});
`),
		c.Expected(
			`var es6 = {};

(function() {
  class Foo {}
  /** @polymerBehavior @nocollapse */
  es6.Behavior = {
    /** @suppress {checkTypes,globalThis,visibility} */
    sayHi: () => new Foo(),
  };
})();

/**
 * @constructor @extends {PolymerElement}
 * @implements {Polymeres6_TestInterface$UID$0}
 */
es6.Test = function() {};
/**
 * @return {?}
 */
es6.Test.prototype.sayHi = () => void 0;

es6.Test = Polymer(/** @lends {es6.Test.prototype} */ {
  is: 'x-element',
  behaviors: [es6.Behavior]
});
/** @export */
es6.Test.prototype.sayHi;
`))
}

// Java: PolymerPassTest.testShorthandLifecycleCallbacks()
func (c *PolymerPassTest) TestShorthandLifecycleCallbacks() {
	c.Test(
		`var ES6Test = Polymer({
  is: 'x-element',
  /** @override */
  created() {
    alert('Shorthand created');
  },
});
`,
		`/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerES6TestInterface$UID$0}
 */
var ES6Test = function() {};

ES6Test = Polymer(/** @lends {ES6Test.prototype} */ {
  is: 'x-element',
  /** @override @this {ES6Test} */
  created() {
    alert('Shorthand created');
  },
});
/** @export */
ES6Test.prototype.created;
`)
}

// Java: PolymerPassTest.testShorthandFunctionDefinitionWithReturn()
func (c *PolymerPassTest) TestShorthandFunctionDefinitionWithReturn() {
	c.Test(
		`var ESTest = Polymer({
  is: 'x-element',
  sayHi() {
    return [1, 2];
  },
});
`,
		`/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerESTestInterface$UID$0}
 */
var ESTest = function() {};

ESTest = Polymer(/** @lends {ESTest.prototype} */ {
  is: 'x-element',
  /** @this {ESTest} */
  sayHi() {
    return [1, 2];
  },
});
/** @export */
ESTest.prototype.sayHi;
`)
}

// Java: PolymerPassTest.testThisTypeAddedToFunctions()
func (c *PolymerPassTest) TestThisTypeAddedToFunctions() {
	c.Test(
		`var X = Polymer({
  is: 'x-element',
  sayHi: function() {
    alert('hi');
  },
  /** @override */
  created: function() {
    this.sayHi();
    this.sayHelloTo_('Tester');
  },
  /**
   * @param {string} name
   * @private
   */
  sayHelloTo_: function(name) {
    alert('Hello, ' + name);
  },
});
`,
		`/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0} */
var X = function() {};
X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',
  /** @this {X} */
  sayHi: function() {
    alert('hi');
  },
  /** @override @this {X} */
  created: function() {
    this.sayHi();
    this.sayHelloTo_('Tester');
  },
  /**
   * @param {string} name
   * @private
   * @this {X}
   */
  sayHelloTo_: function(name) {
    alert('Hello, ' + name);
  },
});
/** @export @private */
X.prototype.sayHelloTo_;
/** @export */
X.prototype.created;
/** @export */
X.prototype.sayHi;
`)

	c.Test(
		`class Foo extends Polymer.Element {
  static get is() { return 'x-element'; }
  sayHi() {
    alert('hi');
  }
  connectedCallback() {
    this.sayHi();
    this.sayHelloTo_('Tester');
  }
  /**
   * @param {string} name
   * @private
   */
  sayHelloTo_(name) {
    alert('Hello, ' + name);
  }
}
`,
		`/** @implements {PolymerFooInterface$UID$0} */
class Foo extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  sayHi() {
    alert('hi');
  }
  connectedCallback() {
    this.sayHi();
    this.sayHelloTo_('Tester');
  }
  /**
   * @param {string} name
   * @private
   */
  sayHelloTo_(name) {
    alert('Hello, ' + name);
  }
}
/** @export @private */
Foo.prototype.sayHelloTo_;
/** @export */
Foo.prototype.connectedCallback;
/** @export */
Foo.prototype.sayHi;
`)
}

// Java: PolymerPassTest.testDollarSignPropsConvertedToBrackets_polymer1Style()
func (c *PolymerPassTest) TestDollarSignPropsConvertedToBrackets_polymer1Style() {
	c.Test(
		`/** @constructor */
var SomeType = function() {};
SomeType.prototype.toggle = function() {};
SomeType.prototype.switch = function() {};
SomeType.prototype.touch = function() {};
var X = Polymer({
  is: 'x-element',
  properties: {
    /** @type {!HTMLElement} */
    propName: {
      type: Object,
      value: function() { return this.$.id; },
    },
  },
  sayHi: function() {
    this.$.checkbox.toggle();
  },
  /** @override */
  created: function() {
    this.sayHi();
    this.$.radioButton.switch();
  },
  /**
   * @param {string} name
   * @private
   */
  sayHelloTo_: function(name) {
    this.$.otherThing.touch();
  },
});
`,
		`/** @constructor */
var SomeType = function() {};
SomeType.prototype.toggle = function() {};
SomeType.prototype.switch = function() {};
SomeType.prototype.touch = function() {};
/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0} */
var X = function() {};
/** @type {!HTMLElement} */
X.prototype.propName;
X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(X, {
    propName: {
      type: Object,
      /** @this {X} @return {!HTMLElement} */
      value: function() { return this.$['id']; },
    },
  }),
  /** @this {X} */
  sayHi: function() {
    this.$['checkbox'].toggle();
  },
  /** @override @this {X} */
  created: function() {
    this.sayHi();
    this.$['radioButton'].switch();
  },
  /**
   * @param {string} name
   * @private
   * @this {X}
   */
  sayHelloTo_: function(name) {
    this.$['otherThing'].touch();
  },
});
/** @export @private */
X.prototype.sayHelloTo_;
/** @export */
X.prototype.created;
/** @export */
X.prototype.sayHi;
`)
}

// Java: PolymerPassTest.testDollarSignPropsConvertedToBrackets_polymer2Style()
func (c *PolymerPassTest) TestDollarSignPropsConvertedToBrackets_polymer2Style() {
	c.IgnoreWarnings(DiagnosticGroups_MISSING_PROPERTIES)
	c.Test(
		`class Foo extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!HTMLElement} */
      propName: {
        type: Object,
        value: function() {
          return /** @type {!HTMLElement} */ (this.$['id']);
        },
      }
    };
  }
  sayHi() {
    this.$.checkbox.toggle();
  }
  connectedCallback() {
    this.sayHi();
    this.$.radioButton.switch();
  }
  /**
   * @param {string} name
   * @private
   */
  sayHelloTo_(name) {
    this.$.otherThing.touch();
  }
}
`,
		`/** @implements {PolymerFooInterface$UID$0} */
class Foo extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(Foo, {
      propName: {
        type: Object,
        /** @this {Foo} @return {!HTMLElement} */
        value: function() {
          return /** @type {!HTMLElement} */ (this.$['id']);
        },
      }
    });
  }
  sayHi() {
    this.$['checkbox'].toggle();
  }
  connectedCallback() {
    this.sayHi();
    this.$['radioButton'].switch();
  }
  /**
   * @param {string} name
   * @private
   */
  sayHelloTo_(name) {
    this.$['otherThing'].touch();
  }
}
/** @type {!HTMLElement} */
Foo.prototype.propName;
/** @export @private */
Foo.prototype.sayHelloTo_;
/** @export */
Foo.prototype.connectedCallback;
/** @export */
Foo.prototype.sayHi;
`)
}

// Java: PolymerPassTest.testDollarSignPropsInShorthandFunctionConvertedToBrackets()
func (c *PolymerPassTest) TestDollarSignPropsInShorthandFunctionConvertedToBrackets() {
	c.Test(
		`/** @constructor */
var SomeType = function() {};
SomeType.prototype.toggle = function() {};
var ES6Test = Polymer({
  is: 'x-element',
  sayHi() {
    this.$.checkbox.toggle();
  },
});
`,
		`/** @constructor */
var SomeType = function() {};
SomeType.prototype.toggle = function() {};

/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerES6TestInterface$UID$0}
 */
var ES6Test = function() {};

ES6Test = Polymer(/** @lends {ES6Test.prototype} */ {
  is: 'x-element',
  /** @this {ES6Test} */
  sayHi() {
    this.$['checkbox'].toggle();
  },
});
/** @export */
ES6Test.prototype.sayHi;
`)
}

// Java: PolymerPassTest.testDollarSignMethodCallInShorthandFunctionNotConvertedToBrackets()
func (c *PolymerPassTest) TestDollarSignMethodCallInShorthandFunctionNotConvertedToBrackets() {
	c.Test(
		`class ReceiverHelper {
  constructor() {}
  foo() {}
};
class Receiver {
  constructor() {
    this.$ = new ReceiverHelper();
  }
  toggle(el) {
    el.toggle();
  }
};
var ES6Test = Polymer({
  is: 'x-element',
  sayHi() {
    const receiver = new Receiver();
    receiver.$.foo();
    toggle(this.$.checkbox);
    receiver.toggle(this.$.checkbox);
  },
  toggle(el) {
    el.toggle();
  },
});
`,
		`class ReceiverHelper {
  constructor() {}
  foo() {}
};
class Receiver {
  constructor() {
    this.$ = new ReceiverHelper();
  }
  toggle(el) {
    el.toggle();
  }
};
/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerES6TestInterface$UID$0}
 */
var ES6Test = function() {};

ES6Test = Polymer(/** @lends {ES6Test.prototype} */ {
  is: 'x-element',
  /** @this {ES6Test} */
  sayHi() {
    const receiver = new Receiver();
    receiver.$.foo();
    toggle(this.$['checkbox']);
    receiver.toggle(this.$['checkbox']);
  },
  /** @this {ES6Test} */
  toggle(el) {
    el.toggle();
  },
});
/** @export */
ES6Test.prototype.toggle;
/** @export */
ES6Test.prototype.sayHi;
`)
}

// Java: PolymerPassTest.testBehaviorForMultipleElements()
func (c *PolymerPassTest) TestBehaviorForMultipleElements() {
	c.Test(
		`/** @constructor */
var FunObject = function() {};
/** @polymerBehavior */
var FunBehavior = {
  properties: {
    /** @type {!FunObject} */
    funThing: {
      type: Object,
    }
  },
};
var A = Polymer({
  is: 'x-element',
  properties: {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  },
  behaviors: [ FunBehavior ],
});
var B = Polymer({
  is: 'y-element',
  properties: {
    foo: String,
  },
  behaviors: [ FunBehavior ],
});
`,
		`/** @constructor */
var FunObject = function() {};
/** @polymerBehavior @nocollapse */
var FunBehavior = {
  properties: {
    funThing: {
      type: Object,
    }
  },
};
/** @constructor @extends {PolymerElement} @implements {PolymerAInterface$UID$0}*/
var A = function() {};
/** @type {!FunObject} */
A.prototype.funThing;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  }),
  behaviors: [ FunBehavior ],
});
/** @constructor @extends {PolymerElement} @implements {PolymerBInterface$m1146332801$1*/
var B = function() {};
/** @type {!FunObject} */
B.prototype.funThing;
/** @type {string} */
B.prototype.foo;
B = Polymer(/** @lends {B.prototype} */ {
  is: 'y-element',
  properties: $jscomp.reflectObject(B, {
    foo: String,
  }),
  behaviors: [ FunBehavior ],
});
`)
}

// Java: PolymerPassTest.testBehaviorInModule_NoBehaviorProps()
func (c *PolymerPassTest) TestBehaviorInModule_NoBehaviorProps() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('Data');

class Item {

}

exports.Item = Item;

/**

 * A Polymer behavior providing common data access and formatting methods.

 * @polymerBehavior

 */

exports.SummaryDataBehavior = {

  /**

   * @param {?Item} item

   * @export

   */

  getValue(item) {

    return this.getItemValue_(item);

  },

};
`,
			`goog.module('Client');
const Data = goog.require('Data');
var A = Polymer({
  is: 'x-element',
  behaviors: [ Data.SummaryDataBehavior ],
});
`))
}

// Java: PolymerPassTest.testBehaviorInModule_WithBehaviorProps_passesTypechecking()
func (c *PolymerPassTest) TestBehaviorInModule_WithBehaviorProps_passesTypechecking() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('Data');
const foo = goog.require('foo');
class Item {}
exports.Item = Item;
/** @polymerBehavior */
exports.SummaryDataBehavior = {
  properties: {
    /** @type {!Item} */
    ABC: {type: Object},
    /** @type {!foo.Item<string>} */
    DEF: {type: Object}
  },
  /** @return {?Item} */
  get x() {
    return null;
  }
};
`,
			`goog.module("client");
const Data = goog.require('Data');

var A = Polymer({
  is: 'x',
  behaviors: [Data.SummaryDataBehavior],
});
`))
}

// Java: PolymerPassTest.testSimpleBehavior()
func (c *PolymerPassTest) TestSimpleBehavior() {
	c.Test(
		c.Srcs(
			`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    /** @type {boolean} */
    isFun: {
      type: Boolean,
      value: true,
    }
  },
  listeners: {
    click: 'doSomethingFun',
  },
  /** @type {string} */
  foo: 'hooray',

  /** @return {number} */
  get someNumber() {
    return 5*7+2;
  },
  /** @param {string} funAmount */
  doSomethingFun: function(funAmount) { alert('Something ' + funAmount + ' fun!'); },
  /** @override */
  created: function() {}
};
var A = Polymer({
  is: 'x-element',
  properties: {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  },
  behaviors: [ FunBehavior ],
});
`),
		c.Expected(
			`/**
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  properties: {
    isFun: {
      type: Boolean,
      value: true,
    }
  },
  listeners: {
    'click': 'doSomethingFun',
  },
  /** @type {string} */
  foo: 'hooray',

  /**
   * @return {number}
   * @suppress {checkTypes,globalThis,visibility}
   */
  get someNumber() {
    return 5 * 7 + 2;
  },
  /**
   * @param {string} funAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingFun: function(funAmount) {
    alert('Something ' + funAmount + ' fun!');
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  created: function() {}
};
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/**
 * @param {string} funAmount
 */
A.prototype.doSomethingFun = function(funAmount) {
  alert('Something ' + funAmount + ' fun!');
};
/** @type {string} */
A.prototype.foo;
/** @type {number} */
A.prototype.someNumber;
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  }),
  behaviors: [FunBehavior],
});
/** @export */
A.prototype.doSomethingFun;
`),
		CompilerTestCase_PostconditionFunc(
			func(compiler *Compiler) {
				// The original doSomethingFun definition in FunBehavior is on line 21, so make sure
				// that line number is preserved when it's copied into the Polymer() call.
				root := compiler.GetRoot()
				visitor := newPolymerPassTest_DoSomethingFunFinder()
				NodeUtil_VisitPreOrder(root, visitor)
				AssertThat(visitor.found).IsTrue()
			}))
}

// Java: private static class PolymerPassTest.DoSomethingFunFinder implements Visitor
type PolymerPassTest_DoSomethingFunFinder struct {
	found bool
}

// Java: PolymerPassTest.DoSomethingFunFinder()
func newPolymerPassTest_DoSomethingFunFinder() *PolymerPassTest_DoSomethingFunFinder {
	return &PolymerPassTest_DoSomethingFunFinder{found: false}
}

// Java: PolymerPassTest.DoSomethingFunFinder.visit(Node)
func (f *PolymerPassTest_DoSomethingFunFinder) Visit(n *rhino.Node) {
	if n.MatchesQualifiedName("A.prototype.doSomethingFun") {
		AssertNode(n).HasLineno(21)
		f.found = true
	}
}

// polymerPassTest_genericNameReplacements is
// ImmutableMap.of("Interface$UID", "Interface$").
func polymerPassTest_genericNameReplacements() *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	m.Put("Interface$UID", "Interface$")
	return m
}

// Java: PolymerPassTest.testIndirectBehaviorAssignment()
func (c *PolymerPassTest) TestIndirectBehaviorAssignment() {
	c.Test(
		c.Srcs(
			`/** @polymerBehavior */
var MyBehavior = {
  properties: {
    behaviorProperty: Boolean
  }
};
var BehaviorAlias1 = MyBehavior;
var BehaviorAlias2 = BehaviorAlias1;
var MyElement = Polymer({
  is: 'my-element',
  behaviors: [ BehaviorAlias2 ]
});
`),
		c.Expected(
			`/** @polymerBehavior @nocollapse */
var MyBehavior = {
  properties: {
    behaviorProperty: Boolean
  }
};
var BehaviorAlias1 = MyBehavior;
var BehaviorAlias2 = BehaviorAlias1;
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerMyElementInterface$UID$0}
 */
var MyElement = function(){};
/** @type {boolean} */
MyElement.prototype.behaviorProperty;
MyElement = Polymer(/** @lends {MyElement.prototype} */ {
  is: 'my-element',
  behaviors: [ BehaviorAlias2 ]
});
`))
}

// Java: PolymerPassTest.testBehaviorWithProtectedMethod()
func (c *PolymerPassTest) TestBehaviorWithProtectedMethod() {
	c.EnableCheckAccessControls()
	c.Test(
		c.Srcs(
			`/** @polymerBehavior */
var FunBehavior = {
  /** @protected */
  doSomethingFun: function() {},
};
`,
			`var A = Polymer({
  is: 'x-element',
  callBehaviorMethod: function() {
    this.doSomethingFun();
  },
  behaviors: [ FunBehavior ],
});
`),
		c.Expected(
			`/**
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  /**
   * @protected
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingFun: function() {},
};
`,
			`/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};

/**
 * @public
 */
A.prototype.doSomethingFun = function(){};

A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  /** @this {A} */
  callBehaviorMethod: function(){ this.doSomethingFun(); },
  behaviors: [FunBehavior],
})
/** @export */
A.prototype.callBehaviorMethod;
/** @export @protected */
A.prototype.doSomethingFun;
`))
}

// Java: PolymerPassTest.testBehaviorWithPrivateMethod()
func (c *PolymerPassTest) TestBehaviorWithPrivateMethod() {
	c.EnableCheckAccessControls()
	c.TestWarning(
		c.Srcs(
			`/** @polymerBehavior */
var FunBehavior = {
  /** @private */
  doSomethingFun: function() {},
};
`,
			`var A = Polymer({
  is: 'x-element',
  callBehaviorMethod: function() {
    this.doSomethingFun();
  },
  behaviors: [ FunBehavior ],
});
`),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: PolymerPassTest.testBehaviorFunctionOverriddenByElement()
func (c *PolymerPassTest) TestBehaviorFunctionOverriddenByElement() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  /** @param {string} funAmount */
  doSomethingFun: function(funAmount) { alert('Something ' + funAmount + ' fun!'); },
};
var A = Polymer({
  is: 'x-element',
  properties: {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  },
  /** @param {string} funAmount */
  doSomethingFun: function(funAmount) {
    alert('Element doing something' + funAmount + ' fun!');
  },
  behaviors: [ FunBehavior ],
});
`,
		`/**
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  /**
   * @param {string} funAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingFun: function(funAmount) {
    alert('Something ' + funAmount + ' fun!');
  },
};
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  }),
  /**
   * @this {A}
   * @param {string} funAmount
   */
  doSomethingFun: function(funAmount) {
    alert('Element doing something' + funAmount + ' fun!');
  },
  behaviors: [FunBehavior],
});
/** @export */
A.prototype.doSomethingFun;
`)
}

// Java: PolymerPassTest.testBehaviorShorthandFunctionOverriddenByElement()
func (c *PolymerPassTest) TestBehaviorShorthandFunctionOverriddenByElement() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  /** @param {string} funAmount */
  doSomethingFun(funAmount) { alert('Something ' + funAmount + ' fun!'); },
};

var A = Polymer({
  is: 'x-element',
  properties: {
    name: String,
  },
  /** @param {string} funAmount */
  doSomethingFun(funAmount) {
    alert('Element doing something' + funAmount + ' fun!');
  },
  behaviors: [ FunBehavior ],
});
`,
		`/**
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  /**
   * @param {string} funAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingFun(funAmount) {
    alert("Something " + funAmount + " fun!");
  },
};

/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};

/** @type {string} */
A.prototype.name;

A = Polymer(/** @lends {A.prototype} */ {
  is: "x-element",
  properties: $jscomp.reflectObject(A, {
    name: String,
  }),
  /**
   * @this {A}
   * @param {string} funAmount
   */
  doSomethingFun(funAmount) {
    alert("Element doing something" + funAmount + " fun!");
  },
  behaviors: [FunBehavior],
});
/** @export */
A.prototype.doSomethingFun;
`)
}

// Java: PolymerPassTest.testBehaviorDefaultValueSuppression()
func (c *PolymerPassTest) TestBehaviorDefaultValueSuppression() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    /** @type {boolean} */
    isFun: {
      type: Boolean,
      value: true,
    },
    funObject: {
      type: Object,
      value: function() { return {fun: this.isFun }; },
    },
    funArray: {
      type: Array,
      value: function() { return [this.isFun]; },
    },
  },
};
var A = Polymer({
  is: 'x-element',
  properties: {
    name: String,
  },
  behaviors: [ FunBehavior ],
});
`,
		`/** @polymerBehavior @nocollapse */
var FunBehavior = {
  properties: {
    isFun: {
      type: Boolean,
      value: true,
    },
    funObject: {
      type: Object,
      /** @suppress {checkTypes|globalThis|visibility} */
      value: function() { return {fun: this.isFun }; },
    },
    funArray: {
      type: Array,
      /** @suppress {checkTypes|globalThis|visibility} */
      value: function() { return [this.isFun]; },
    },
  },
};
/** @constructor @extends {PolymerElement} @implements {PolymerAInterface$UID$0}*/
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
/** @type {!Object} */
A.prototype.funObject;
/** @type {!Array} */
A.prototype.funArray;
/** @type {string} */
A.prototype.name;
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    name: String,
  }),
  behaviors: [ FunBehavior ],
});
`)
}

// Java: PolymerPassTest.testArrayBehavior()
func (c *PolymerPassTest) TestArrayBehavior() {
	c.Test(
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
  properties: {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  },
  behaviors: [ SuperCoolBehaviors, BoringBehavior ],
});
`,
		`/**
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
  /**
   * @param {string} funAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingFun: function(funAmount) {
    alert('Something ' + funAmount + ' fun!');
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  created: function() {}
};
/**
 * @nocollapse
 * @polymerBehavior
 */
var RadBehavior = {
  properties: {
    howRad: Number
  },
  /**
   * @param {number} radAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingRad: function(radAmount) {
    alert('Something ' + radAmount + ' rad!');
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  ready: function() {}
};
/**
 * @nocollapse
 * @polymerBehavior
 */
var SuperCoolBehaviors = [FunBehavior, RadBehavior];
/**
 * @nocollapse
 * @polymerBehavior
 */
var BoringBehavior = {
  properties: {
    boringString: String
  },
  /**
   * @param {boolean} boredYet
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomething: function(boredYet) {
    alert(boredYet + ' ' + this.boringString);
  },
};
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
/** @type {number} */
A.prototype.howRad;
/** @type {string} */
A.prototype.boringString;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/**
 * @param {string} funAmount
 */
A.prototype.doSomethingFun = function(funAmount) {
  alert('Something ' + funAmount + ' fun!');
};
/**
 * @param {number} radAmount
 */
A.prototype.doSomethingRad = function(radAmount) {
  alert('Something ' + radAmount + ' rad!');
};
/**
 * @param {boolean} boredYet
 */
A.prototype.doSomething = function(boredYet) {
  alert(boredYet + ' ' + this.boringString);
};
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  }),
  behaviors: [SuperCoolBehaviors, BoringBehavior],
});
/** @export */
A.prototype.doSomething;
/** @export */
A.prototype.doSomethingRad;
/** @export */
A.prototype.doSomethingFun;
`)
}

// Java: PolymerPassTest.testInlineLiteralBehavior()
func (c *PolymerPassTest) TestInlineLiteralBehavior() {
	c.Test(
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
  properties: {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  },
  behaviors: [ SuperCoolBehaviors ],
});
`,
		`/**
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
  /**
   * @param {string} funAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingFun: function(funAmount) {
    alert('Something ' + funAmount + ' fun!');
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  created: function() {}
};
/**
 * @nocollapse
 * @polymerBehavior
 */
var SuperCoolBehaviors = [
  FunBehavior, {
    properties: {
      howRad: Number
    },
    /**
     * @param {number} radAmount
     * @suppress {checkTypes,globalThis,visibility}
     */
    doSomethingRad: function(radAmount) {
      alert('Something ' + radAmount + ' rad!');
    },
    /**
     * @override
     * @suppress {checkTypes,globalThis,visibility}
     */
    ready: function() {}
  }
];
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
/** @type {number} */
A.prototype.howRad;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/**
 * @param {string} funAmount
 */
A.prototype.doSomethingFun = function(funAmount) {
  alert('Something ' + funAmount + ' fun!');
};
/**
 * @param {number} radAmount
 */
A.prototype.doSomethingRad = function(radAmount) {
  alert('Something ' + radAmount + ' rad!');
};
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  }),
  behaviors: [SuperCoolBehaviors],
});
/** @export */
A.prototype.doSomethingRad;
/** @export */
A.prototype.doSomethingFun;
`)
}

// Java: PolymerPassTest.testGlobalBehaviorPropsPrototypesGetCreated()
func (c *PolymerPassTest) TestGlobalBehaviorPropsPrototypesGetCreated() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
}
var A = Polymer({
  is: 'x-element',
  behaviors: [ FunBehavior]
});
`,
		`/** @polymerBehavior @nocollapse */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
}
/** @constructor @extends {PolymerElement} @implements {PolymerAInterface$UID$0}*/
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  behaviors: [ FunBehavior ],
});
`)
}

// Java: PolymerPassTest.testGlobalBehaviorPropsPrototypesGetCreated_transitively()
func (c *PolymerPassTest) TestGlobalBehaviorPropsPrototypesGetCreated_transitively() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
}
/** @polymerBehavior */
var SuperCoolBehaviors = [FunBehavior];
var A = Polymer({
  is: 'x-element',
  behaviors: [ SuperCoolBehaviors]
});
`,
		`/** @polymerBehavior @nocollapse */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
}
/** @polymerBehavior @nocollapse */
var SuperCoolBehaviors = [FunBehavior];
/** @constructor @extends {PolymerElement} @implements {PolymerAInterface$UID$0}*/
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  behaviors: [ SuperCoolBehaviors ],
});
`)
}

// Java: PolymerPassTest.testBehaviorMethodsInPropertiesAndObservers()
func (c *PolymerPassTest) TestBehaviorMethodsInPropertiesAndObservers() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    computedProp: {
      type: String,
      computed: 'quotedProp(prop1, prop2)',
    },
    observedProp: {
      type: String,
      observer: 'observeProp',
    },
  },
  observers: [
    'observeList(item1, item2)',
    'observeAnotherList(item3)',
  ],
  'quotedProp': function(prop1, prop2) {},
  observeProp() {},
  observeList: function(item1, item2) {},
  observeAnotherList(item3) {},
  unreferencedMethod: function() {},
};
var A = Polymer({
  is: 'x-element',
  behaviors: [ FunBehavior ],
});
`,
		`/**
 * @polymerBehavior
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  properties: {
    computedProp: {
      type: String,
      computed: 'quotedProp(prop1, prop2)',
    },
    observedProp: {
      type: String,
      observer: 'observeProp',
    },
  },
  observers: [
    'observeList(item1, item2)',
    'observeAnotherList(item3)',
  ],
  /**
   * @export
   * @nocollapse
   * @suppress {checkTypes,globalThis,visibility}
   */
  'quotedProp': function(prop1, prop2) {},
  /**
   * @export
   * @nocollapse
   * @suppress {checkTypes,globalThis,visibility}
   */
  observeProp() {},
  /**
   * @export
   * @nocollapse
   * @suppress {checkTypes,globalThis,visibility}
   */
  observeList: function(item1, item2) {},
  /**
   * @export
   * @nocollapse
   * @suppress {checkTypes,globalThis,visibility}
   */
  observeAnotherList(item3) {},
  /**
   * @suppress {checkTypes,globalThis,visibility}
   */
  unreferencedMethod: function() {},
};
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};
/** @type {string} */
A.prototype.computedProp;
/** @type {string} */
A.prototype.observedProp;
A.prototype.quotedProp = function(prop1, prop2) {};
A.prototype.observeProp = function() {};
A.prototype.observeList = function(item1, item2) {};
A.prototype.observeAnotherList = function(item3) {};
A.prototype.unreferencedMethod = function() {};
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  behaviors: [FunBehavior],
});
/** @export */
A.prototype.unreferencedMethod;
/** @export */
A.prototype.observeAnotherList;
/** @export */
A.prototype.observeList;
/** @export */
A.prototype.observeProp;
/** @export */
A.prototype.quotedProp;
`)
}

// Java: PolymerPassTest.testBehaviorFunctionOverriding()
func (c *PolymerPassTest) TestBehaviorFunctionOverriding() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
  /** @param {boolean} boredYet */
  doSomething: function(boredYet) { alert(boredYet + ' ' + this.isFun); },
  /** @override */
  created: function() {}
};
/** @polymerBehavior */
var RadBehavior = {
  properties: {
    howRad: Number
  },
  /** @param {boolean} boredYet */
  doSomething: function(boredYet) { alert(boredYet + ' ' + this.howRad); },
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
  properties: {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  },
  behaviors: [ SuperCoolBehaviors, BoringBehavior ],
});
`,
		`/**
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
  /**
   * @param {boolean} boredYet
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomething: function(boredYet) {
    alert(boredYet + ' ' + this.isFun);
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  created: function() {}
};
/**
 * @nocollapse
 * @polymerBehavior
 */
var RadBehavior = {
  properties: {
    howRad: Number
  },
  /**
   * @param {boolean} boredYet
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomething: function(boredYet) {
    alert(boredYet + ' ' + this.howRad);
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  ready: function() {}
};
/**
 * @nocollapse
 * @polymerBehavior
 */
var SuperCoolBehaviors = [FunBehavior, RadBehavior];
/**
 * @nocollapse
 * @polymerBehavior
 */
var BoringBehavior = {
  properties: {
    boringString: String
  },
  /**
   * @param {boolean} boredYet
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomething: function(boredYet) {
    alert(boredYet + ' ' + this.boringString);
  },
};
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
/** @type {number} */
A.prototype.howRad;
/** @type {string} */
A.prototype.boringString;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/**
 * @param {boolean} boredYet
 */
A.prototype.doSomething = function(boredYet) {
  alert(boredYet + ' ' + this.boringString);
};
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  }),
  behaviors: [SuperCoolBehaviors, BoringBehavior],
});
/** @export */
A.prototype.doSomething;
`)
}

// Java: PolymerPassTest.testBehaviorShorthandFunctionOverriding()
func (c *PolymerPassTest) TestBehaviorShorthandFunctionOverriding() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
  /** @param {boolean} boredYet */
  doSomething(boredYet) { alert(boredYet + ' ' + this.isFun); },
};

/** @polymerBehavior */
var RadBehavior = {
  properties: {
    howRad: Number
  },
  /** @param {boolean} boredYet */
  doSomething(boredYet) { alert(boredYet + ' ' + this.howRad); },
};

/** @polymerBehavior */
var SuperCoolBehaviors = [FunBehavior, RadBehavior];

/** @polymerBehavior */
var BoringBehavior = {
  properties: {
    boringString: String
  },
  /** @param {boolean} boredYet */
  doSomething(boredYet) { alert(boredYet + ' ' + this.boringString); },
};

var A = Polymer({
  is: 'x-element',
  properties: {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  },
  behaviors: [ SuperCoolBehaviors, BoringBehavior ],
});
`,
		`/**
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
  /**
   * @param {boolean} boredYet
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomething(boredYet) {
    alert(boredYet + ' ' + this.isFun);
  },
};

/**
 * @nocollapse
 * @polymerBehavior
 */
var RadBehavior = {
  properties: {
    howRad: Number
  },
  /**
   * @param {boolean} boredYet
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomething(boredYet) {
    alert(boredYet + ' ' + this.howRad);
  },
};

/**
 * @nocollapse
 * @polymerBehavior
 */
var SuperCoolBehaviors = [FunBehavior, RadBehavior];
/**
 * @nocollapse
 * @polymerBehavior
 */
var BoringBehavior = {
  properties: {
    boringString: String
  },
  /**
   * @param {boolean} boredYet
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomething(boredYet) {
    alert(boredYet + ' ' + this.boringString);
  },
};

/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};

/** @type {boolean} */
A.prototype.isFun;
/** @type {number} */
A.prototype.howRad;
/** @type {string} */
A.prototype.boringString;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/**
 * @param {boolean} boredYet
 */
A.prototype.doSomething = function(boredYet) {
  alert(boredYet + ' ' + this.boringString);
};

A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  }),
  behaviors: [SuperCoolBehaviors, BoringBehavior],
});
/** @export */
A.prototype.doSomething;
`)
}

// Java: PolymerPassTest.testBehaviorReadOnlyProp()
func (c *PolymerPassTest) TestBehaviorReadOnlyProp() {
	js := `/** @polymerBehavior */
var FunBehavior = {
  properties: {
    isFun: {
      type: Boolean,
      readOnly: true,
    },
  },
  /** @param {string} funAmount */
  doSomethingFun: function(funAmount) { alert('Something ' + funAmount + ' fun!'); },
  /** @override */
  created: function() {}
};
var A = Polymer({
  is: 'x-element',
  properties: {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  },
  behaviors: [ FunBehavior ],
});
`

	c.Test(
		js,
		`/**
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  properties: {
    isFun: {
      type: Boolean,
      readOnly: true,
    },
  },
  /**
   * @param {string} funAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingFun: function(funAmount) {
    alert('Something ' + funAmount + ' fun!');
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  created: function() {}
};
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/**
 * @param {string} funAmount
 */
A.prototype.doSomethingFun = function(funAmount) {
  alert('Something ' + funAmount + ' fun!');
};
/**
 * @param {boolean} isFun
 * @override
 */
A.prototype._setIsFun = function(isFun) {};
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  }),
  behaviors: [FunBehavior],
});
/** @export */
A.prototype.doSomethingFun;
`)

	c.TestExternChanges(
		c.Externs(polymerPassTest_EXTERNS),
		c.Srcs(js),
		c.Expected(
			`/** @interface */
var PolymerAInterface$m1146332801$0 = function() {};
/** @type {?} */
PolymerAInterface$m1146332801$0.prototype.isFun;
/** @type {?} */
PolymerAInterface$m1146332801$0.prototype.pets;
/** @type {?} */
PolymerAInterface$m1146332801$0.prototype.name;
/** @param {?} isFun **/
PolymerAInterface$m1146332801$0.prototype._setIsFun;
`,
			polymerPassTest_EXTERNS))
}

// Java: PolymerPassTest.testBehaviorInIIFE()
func (c *PolymerPassTest) TestBehaviorInIIFE() {
	c.Test(
		`(function() {
  /** @polymerBehavior */
  Polymer.FunBehavior = {
    properties: {
      /** @type {boolean} */
      isFun: {
        type: Boolean,
        value: true,
      }
    },
    /** @type {string} */
    foo: 'hooray',
    /** @return {number} */
    get someNumber() {
      return 5*7+2;
    },
    /** @param {string} funAmount */
    doSomethingFun: function(funAmount) {
      alert('Something ' + funAmount + ' fun!');
    },
    /** @override */
    created: function() {}
  };
})();
var A = Polymer({
  is: 'x-element',
  properties: {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  },
  behaviors: [ Polymer.FunBehavior ],
});
`,
		`(function() {
/**
 * @nocollapse
 * @polymerBehavior
 */
Polymer.FunBehavior = {
  properties: {
    isFun: {
      type: Boolean,
      value: true,
    }
  },
  /** @type {string} */
  foo: 'hooray',
  /**
   * @return {number}
   * @suppress {checkTypes,globalThis,visibility}
   */
  get someNumber() {
    return 5 * 7 + 2;
  },
  /**
   * @param {string} funAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingFun: function(funAmount) {
    alert('Something ' + funAmount + ' fun!');
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  created: function() {}
};
})();
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/**
 * @param {string} funAmount
 * @return {?}
 */
A.prototype.doSomethingFun = function(funAmount) {};
/** @type {string} */
A.prototype.foo;
/** @type {number} */
A.prototype.someNumber;
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  }),
  behaviors: [Polymer.FunBehavior],
});
/** @export */
A.prototype.doSomethingFun;
`)
}

// Java: PolymerPassTest.testBehaviorInIIFE_NoTypeAnnotation()
func (c *PolymerPassTest) TestBehaviorInIIFE_NoTypeAnnotation() {
	c.Test(
		`(function() {
  /** @polymerBehavior */
  Polymer.FunBehavior = {
    properties: { // no annotation /** @type {number} */
      isFun: {
        type: Boolean,
        value: true,
      }
    },
  };
})();
var A = Polymer({
  is: 'x-element',
  behaviors: [ Polymer.FunBehavior ],
});
`,
		`(function() {
  /** @polymerBehavior @nocollapse */
  Polymer.FunBehavior = {
    properties: {
      isFun: {
        type: Boolean,
        value: true,
      }
    },
  };
})();
/** @constructor @extends {PolymerElement} @implements {PolymerAInterface$UID$0}*/
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  behaviors: [ Polymer.FunBehavior ],
});
`)
}

// Java: PolymerPassTest.testBehaviorInModule()
func (c *PolymerPassTest) TestBehaviorInModule() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('behaviors.CoolBehavior');
goog.module.declareLegacyNamespace();
const MODULE_LOCAL = 0;
/** @polymerBehavior */
Polymer.FunBehavior = {
  /** @param {string} funAmount */
  doSomethingFun: function(funAmount) {
    alert(MODULE_LOCAL);
    alert('Something ' + funAmount + ' fun!');
  },
  /** @override */
  created: function() {}
};
/** @polymerBehavior */
Polymer.VeryFunBehavior = [
  Polymer.FunBehavior,
  {
    doSomethingVeryFun: function(a, [b], c, veryFun = MODULE_LOCAL) {
      alert(MODULE_LOCAL);
      alert('Something very ' + veryFunAmount + ' fun!');
    },
    /** @override */
    created: function() {}
  }
]
/** @polymerBehavior */
exports = {
  /** @param {string} coolAmount */
  doSomethingCool: function(coolAmount) {
    alert(MODULE_LOCAL);
    alert('Something ' + funAmount + ' cool!');
  },
  /** @override */
  created: function() {}
};
`,
			`goog.require('behaviors.CoolBehavior');
var A = Polymer({
  is: 'x-element',
  properties: {
    name: String,
  },
  behaviors: [ Polymer.VeryFunBehavior, behaviors.CoolBehavior ],
});
`),
		c.Expected(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module("behaviors.CoolBehavior");
goog.module.declareLegacyNamespace();
const MODULE_LOCAL = 0;
/**
 * @nocollapse
 * @polymerBehavior
 */
Polymer.FunBehavior = {
  /**
   * @param {string} funAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingFun:function(funAmount) {
    alert(MODULE_LOCAL);
    alert("Something " + funAmount + " fun!");
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  created: function() {}
};
/**
 * @nocollapse
 * @polymerBehavior
 */
Polymer.VeryFunBehavior = [
  Polymer.FunBehavior,
  {
    /**
     * @suppress {checkTypes,globalThis,visibility}
     */
    doSomethingVeryFun: function(a, [b], c, veryFun = MODULE_LOCAL) {
      alert(MODULE_LOCAL);
      alert("Something very " + veryFunAmount + " fun!");
    },
    /**
     * @override
     * @suppress {checkTypes,globalThis,visibility}
     */
    created: function() {}
  }
];
/**
 * @nocollapse
 * @polymerBehavior
 */
exports = {
  /**
   * @param {string} coolAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingCool: function(coolAmount) {
    alert(MODULE_LOCAL);
    alert("Something " + funAmount + " cool!");
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  created: function() {}
};
`,
			`goog.require('behaviors.CoolBehavior');
/** @constructor @extends {PolymerElement} @implements {PolymerAInterface$UID$0}*/
var A = function() {};
/** @type {string} */ A.prototype.name;
/**
 * @param {string} funAmount
 * @return {?}
 */
A.prototype.doSomethingFun = function(funAmount) {};
/**
 * @return {?}
 */
A.prototype.doSomethingVeryFun =
    function(a, param$polymer$1, c, veryFun = void 0) {};
/**
 * @param {string} coolAmount
 * @return {?}
 */
A.prototype.doSomethingCool = function(coolAmount) {};
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    name: String,
  }),
  behaviors: [ Polymer.VeryFunBehavior, behaviors.CoolBehavior ],
});
/** @export */
A.prototype.doSomethingCool;
/** @export */
A.prototype.doSomethingVeryFun;
/** @export */
A.prototype.doSomethingFun;
`))
}

// Java: PolymerPassTest.testBehaviorInModule_destructuringRestParams()
func (c *PolymerPassTest) TestBehaviorInModule_destructuringRestParams() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module('behaviors.CoolBehavior');
goog.module.declareLegacyNamespace();
const MODULE_LOCAL = 0;
/** @polymerBehavior */
exports = {
  doSomethingCool: function([coolAmount = MODULE_LOCAL] = [], ...[more]) {
    alert(MODULE_LOCAL);
    alert('Something ' + funAmount + ' cool!');
  },
  /** @override */
  created: function() {}
};
`,
			`goog.require('behaviors.CoolBehavior');
var A = Polymer({
  is: 'x-element',
  properties: {
    name: String,
  },
  behaviors: [ behaviors.CoolBehavior ],
});
`),
		c.Expected(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.module("behaviors.CoolBehavior");
goog.module.declareLegacyNamespace();
const MODULE_LOCAL = 0;
/**
 * @nocollapse
 * @polymerBehavior
 */
exports = {
  /**
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingCool: function([coolAmount = MODULE_LOCAL] = [], ...[more]) {
    alert(MODULE_LOCAL);
    alert("Something " + funAmount + " cool!");
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  created: function() {}
};
`,
			`goog.require('behaviors.CoolBehavior');
/** @constructor @extends {PolymerElement} @implements {PolymerAInterface$UID$0}*/
var A = function() {};
/** @type {string} */ A.prototype.name;
/**
 * @return {?}
 */
A.prototype.doSomethingCool =
    function(param$polymer$0 = void 0, ...param$polymer$1) {};
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    name: String,
  }),
  behaviors: [ behaviors.CoolBehavior ],
});
/** @export */
A.prototype.doSomethingCool;
`))
}

// Java: PolymerPassTest.testDuplicatedBehaviorsAreCopiedOnce()
func (c *PolymerPassTest) TestDuplicatedBehaviorsAreCopiedOnce() {
	c.Test(
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
var A = Polymer({
  is: 'x-element',
  properties: {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  },
  behaviors: [ SuperCoolBehaviors, FunBehavior ],
});
`,
		`/**
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  properties: {
    isFun: Boolean
  },
  /**
   * @param {string} funAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingFun: function(funAmount) {
    alert('Something ' + funAmount + ' fun!');
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  created: function() {}
};
/**
 * @nocollapse
 * @polymerBehavior
 */
var RadBehavior = {
  properties: {
    howRad: Number
  },
  /**
   * @param {number} radAmount
   * @suppress {checkTypes,globalThis,visibility}
   */
  doSomethingRad: function(radAmount) {
    alert('Something ' + radAmount + ' rad!');
  },
  /**
   * @override
   * @suppress {checkTypes,globalThis,visibility}
   */
  ready: function() {}
};
/**
 * @nocollapse
 * @polymerBehavior
 */
var SuperCoolBehaviors = [FunBehavior, RadBehavior];
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerAInterface$UID$0}
 */
var A = function() {};
/** @type {boolean} */
A.prototype.isFun;
/** @type {number} */
A.prototype.howRad;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/**
 * @param {number} radAmount
 */
A.prototype.doSomethingRad = function(radAmount) {
  alert('Something ' + radAmount + ' rad!');
};
/**
 * @param {string} funAmount
 */
A.prototype.doSomethingFun = function(funAmount) {
  alert('Something ' + funAmount + ' fun!');
};
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
  }),
  behaviors: [SuperCoolBehaviors, FunBehavior],
});
/** @export */
A.prototype.doSomethingRad;
/** @export */
A.prototype.doSomethingFun;
`)
}

// Java: PolymerPassTest.testInvalid1()
func (c *PolymerPassTest) TestInvalid1() {
	c.TestWarning("var x = Polymer('blah');", PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID)
	c.Test(
		c.Srcs("var x = Polymer('foo-bar', {});"),
		c.Warning(PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID),
		c.Warning(TypeCheck_WRONG_ARGUMENT_COUNT))
	c.TestError("var x = Polymer({},'blah');", PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS)
	c.TestError("var x = Polymer({});", PolymerPassErrors_POLYMER_MISSING_IS)

	c.TestError(
		`var x = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() { return '' }
};
`,
		PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID)

	c.TestError(
		`var x = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  get properties() { return {} }
};
`,
		PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC)
}

// Java: PolymerPassTest.testInvalidProperties()
func (c *PolymerPassTest) TestInvalidProperties() {
	c.TestError(
		`Polymer({
  is: 'x-element',
  properties: {
    isHappy: true,
  },
});
`,
		PolymerPassErrors_POLYMER_INVALID_PROPERTY)

	c.TestError(
		`Polymer({
  is: 'x-element',
  properties: {
    isHappy: {
      value: true,
    },
  },
});
`,
		PolymerPassErrors_POLYMER_INVALID_PROPERTY)

	c.TestError(
		`var foo = {};
foo.bar = {};
Polymer({
  is: 'x-element',
  properties: {
    isHappy: {
      type: foo.Bar,
      value: true,
    },
  },
});
`,
		PolymerPassErrors_POLYMER_INVALID_PROPERTY)

	c.TestError(
		`var x = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      isHappy: true,
    };
  }
};
`,
		PolymerPassErrors_POLYMER_INVALID_PROPERTY)

	c.TestError(
		`var x = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      isHappy: {
        value: true,
      }
    };
  }
};
`,
		PolymerPassErrors_POLYMER_INVALID_PROPERTY)

	c.TestError(
		`var x = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      isHappy: {
        type: foo.Bar,
        value: true,
      }
    };
  }
};
`,
		PolymerPassErrors_POLYMER_INVALID_PROPERTY)
}

// Java: PolymerPassTest.testInvalidBehavior()
func (c *PolymerPassTest) TestInvalidBehavior() {
	c.TestError(
		`(function() {
  var isNotGloabl = {};
  Polymer({
    is: 'x-element',
    behaviors: [
      isNotGlobal
    ],
  });
})();
`,
		PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR)

	c.TestError(
		`var foo = {};
(function() {
  Polymer({
    is: 'x-element',
    behaviors: [
      foo.IsNotDefined
    ],
  });
})();
`,
		PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR)

	c.TestError(
		`var foo = {};
foo.Bar;
(function() {
  Polymer({
    is: 'x-element',
    behaviors: [
      foo.Bar
    ],
  });
})();
`,
		PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR)

	c.TestError(
		`Polymer({
  is: 'x-element',
  behaviors: [
    DoesNotExist
  ],
});
`,
		PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR)
}

// Java: PolymerPassTest.testUnannotatedBehavior()
func (c *PolymerPassTest) TestUnannotatedBehavior() {
	c.TestError(
		`var FunBehavior = {
  /** @override */
  created: function() {}
};
var A = Polymer({
  is: 'x-element',
  behaviors: [ FunBehavior ],
});
`,
		PolymerPassErrors_POLYMER_UNANNOTATED_BEHAVIOR)
}

// Java: PolymerPassTest.testInvalidTypeAssignment()
func (c *PolymerPassTest) TestInvalidTypeAssignment() {
	c.Test(
		`var X = Polymer({
  is: 'x-element',
  properties: {
    isHappy: Boolean,
  },
  /** @override */
  created: function() {
    this.isHappy = 7;
  },
});
`,
		`/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0} */
var X = function() {};
/** @type {boolean} */
X.prototype.isHappy;
X = Polymer(/** @lends {X.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(X, {
    isHappy: Boolean,
  }),
  /** @override @this {X} */
  created: function() {
    this.isHappy = 7;
  },
});
/** @export */
X.prototype.created;
`,
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))
}

// Java: PolymerPassTest.testFeaturesInFunctionBody()
func (c *PolymerPassTest) TestFeaturesInFunctionBody() {
	c.Test(
		"var X = Polymer({\n  is: 'x-element',\n  funcWithES6() {\n    var tag = 'tagged';\n    alert(`${tag}Template`);\n    var taggedTemp = `${tag}Template`;\n    var arrFunc = () => 42;\n    var num = arrFunc();\n    var obj = {one: 1, two: 2, three: 3};\n    var {one, two, three} = obj;\n    var arr = [1, 2, 3];\n    var [eins, zwei, drei] = arr;\n  },\n});\n",
		"/** @constructor @extends {PolymerElement} @implements {PolymerXInterface$UID$0} */\nvar X = function() {};\nX = Polymer(/** @lends {X.prototype} */ {\n  is: 'x-element',\n  /** @this {X} */\n  funcWithES6() {\n    var tag = 'tagged';\n    alert(`${tag}Template`);\n    var taggedTemp = `${tag}Template`;\n    var arrFunc = () => 42;\n    var num = arrFunc();\n    var obj = {one: 1, two: 2, three: 3};\n    var {one, two, three} = obj;\n    var arr = [1, 2, 3];\n    var [eins, zwei, drei] = arr;\n  },\n});\n/** @export */\nX.prototype.funcWithES6;\n")
}

// Java: PolymerPassTest.testPolymerElementAnnotation1()
func (c *PolymerPassTest) TestPolymerElementAnnotation1() {
	c.Test(
		`class Bar {}
/** @constructor */
var User = function() {};
/** @polymer */
class Foo extends Bar {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function
    };
  }
}
`,
		`class Bar {}
/** @constructor */
var User = function() {};
/** @polymer @implements {PolymerFooInterface$UID$0} */
class Foo extends Bar {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(Foo, {
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function
    });
  }
}
/** @type {!User} @private */
Foo.prototype.user_;
/** @type {!Array} */
Foo.prototype.pets;
/** @type {string} */
Foo.prototype.name;
/** @type {!Function} */
Foo.prototype.thingToDo;
`)
}

// Java: PolymerPassTest.testPolymerElementAnnotation2()
func (c *PolymerPassTest) TestPolymerElementAnnotation2() {
	c.Test(
		`class Foo {}
/** @constructor */
var User = function() {};
var a = {};
/** @polymer */
a.B = class extends Foo {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: {
      type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function
    };
  }
};
`,
		`class Foo {}
/** @constructor */
var User = function() {};
var a = {};
/** @polymer @implements {Polymera_BInterface$UID$0} */
a.B = class extends Foo {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(a.B, {
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function,
    });
  }
};
/** @type {!User} @private */
a.B.prototype.user_;
/** @type {!Array} */
a.B.prototype.pets;
/** @type {string} */
a.B.prototype.name;
/** @type {!Function} */
a.B.prototype.thingToDo;
`)
}

// Java: PolymerPassTest.testPolymerElementAnnotation3()
func (c *PolymerPassTest) TestPolymerElementAnnotation3() {
	c.Test(
		`class Foo {}
/** @interface */
function User() {};
/** @type {boolean} */ User.prototype.id;
var a = {};
/**
 * @polymer
 * @implements {User}
 */
a.B = class extends Foo {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      id: Boolean,
      other: {
        type: String,
        reflectToAttribute: true
      }
    };
  }
};
`,
		`class Foo {}
/** @interface */
function User() {};
/** @type {boolean} */ User.prototype.id;
var a = {};
/**
 * @polymer
 * @implements {User}
 * @implements {Polymera_BInterface$UID$0}
 */
a.B = class extends Foo {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(a.B, {
      id: Boolean,
      other: {
        type: String,
        reflectToAttribute: true
      }
    });
  }
};
/** @type {boolean} */
a.B.prototype.id;
/** @type {string} */
a.B.prototype.other;
`)
}

// Java: PolymerPassTest.testObjectReflectionAddedToConfigProperties1()
func (c *PolymerPassTest) TestObjectReflectionAddedToConfigProperties1() {
	c.Test(
		`/** @constructor */
var User = function() {};
/** @const */ var a = {};
a.B = Polymer({
  is: 'x-element',
  properties: {
    /** @type {!User} @private */
    user_: Object,
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
    thingToDo: Function,
  },
});
`,
		`/** @constructor */
var User = function() {};
/** @const */ var a = {};
/** @constructor @extends {PolymerElement} @implements {Polymera_BInterface$UID$0}*/
a.B = function() {};
/** @type {!User} @private */
a.B.prototype.user_;
/** @type {!Array} */
a.B.prototype.pets;
/** @type {string} */
a.B.prototype.name;
/** @type {!Function} */
a.B.prototype.thingToDo;
a.B = Polymer(/** @lends {a.B.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(a.B, {
    user_: Object,
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
    thingToDo: Function,
  }),
});
`)

	c.Test(
		`/** @constructor */
var User = function() {};
var a = {};
a.B = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function
    };
  }
};
`,
		`/** @constructor */
var User = function() {};
var a = {};
/** @implements {Polymera_BInterface$UID$0} */
a.B = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return  $jscomp.reflectObject(a.B, {
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function
    });
  }
};
/** @type {!User} @private */
a.B.prototype.user_;
/** @type {!Array} */
a.B.prototype.pets;
/** @type {string} */
a.B.prototype.name;
/** @type {!Function} */
a.B.prototype.thingToDo;
`)
}

// Java: PolymerPassTest.testObjectReflectionAddedToConfigProperties2()
func (c *PolymerPassTest) TestObjectReflectionAddedToConfigProperties2() {
	c.Test(
		`/** @constructor */
var User = function() {};
var A = Polymer({
  is: 'x-element',
  properties: {
    /** @type {!User} @private */
    user_: Object,
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
    thingToDo: Function,
  },
});
`,
		`/** @constructor */
var User = function() {};
/** @constructor @extends {PolymerElement} @implements {PolymerAInterface$UID$0}*/
var A = function() {};
/** @type {!User} @private */
A.prototype.user_;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/** @type {!Function} */
A.prototype.thingToDo;
A = Polymer(/** @lends {A.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(A, {
    user_: Object,
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
    thingToDo: Function,
  }),
});
`)

	c.Test(
		`/** @constructor */
var User = function() {};
const A = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function,
    };
  }
};
`,
		`/** @constructor */
var User = function() {};
/** @implements {PolymerAInterface$UID$0} */
const A = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(A, {
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function,
    });
  }
};
/** @type {!User} @private */
A.prototype.user_;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/** @type {!Function} */
A.prototype.thingToDo;
`)
}

// Java: PolymerPassTest.testObjectReflectionAddedToConfigProperties3()
func (c *PolymerPassTest) TestObjectReflectionAddedToConfigProperties3() {
	c.Test(
		`/** @constructor */
var User = function() {};
Polymer({
  is: 'x-element',
  properties: {
    /** @type {!User} @private */
    user_: Object,
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
    thingToDo: Function,
  },
});
`,
		`/** @constructor */
var User = function() {};
/**
 * @constructor @extends {PolymerElement}
 * @implements {PolymerXElementElementInterface$UID$0}
 */
var XElementElement = function() {};
/** @type {!User} @private */
XElementElement.prototype.user_;
/** @type {!Array} */
XElementElement.prototype.pets;
/** @type {string} */
XElementElement.prototype.name;
/** @type {!Function} */
XElementElement.prototype.thingToDo;
Polymer(/** @lends {XElementElement.prototype} */ {
  is: 'x-element',
  properties: $jscomp.reflectObject(XElementElement, {
    user_: Object,
    pets: {
      type: Array,
      notify: true,
    },
    name: String,
    thingToDo: Function,
  }),
});
`)

	c.Test(
		`/** @constructor */
var User = function() {};
class XElement extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function,
    };
  }
}
`,
		`/** @constructor */
var User = function() {};
/** @implements {PolymerXElementInterface$UID$0} */
class XElement extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(XElement, {
      user_: Object,
      pets: {
        type: Array,
        notify: true,
      },
      name: String,
      thingToDo: Function,
    });
  }
}
/** @type {!User} @private */
XElement.prototype.user_;
/** @type {!Array} */
XElement.prototype.pets;
/** @type {string} */
XElement.prototype.name;
/** @type {!Function} */
XElement.prototype.thingToDo;
`)
}

// Java: PolymerPassTest.testExportsMethodsFromClassBasedElement()
func (c *PolymerPassTest) TestExportsMethodsFromClassBasedElement() {
	c.Test(
		`class TestElement extends PolymerElement {
  /** @public */ method1() {}
  /** @private */ method2() {}
}
`,
		`/** @implements {PolymerTestElementInterface$UID$0} */
class TestElement extends PolymerElement {
  /** @public */ method1() {}
  /** @private */ method2() {}
}
/** @private @export */ TestElement.prototype.method2;
/** @public @export */ TestElement.prototype.method1;
`)
}

// Java: PolymerPassTest.testExportMethodsFromLegacyElement()
func (c *PolymerPassTest) TestExportMethodsFromLegacyElement() {
	c.Test(
		`Polymer({
  is: 'test-element',
  /** @public */ method1() {},
  /** @private */ method2() {},
});
`,
		`/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerTestElementElementInterface$UID$0}
 */
var TestElementElement = function() {};
Polymer(/** @lends {TestElementElement.prototype} */ {
  is: "test-element",
  /** @public @this {TestElementElement} */ method1() {},
  /** @private @this {TestElementElement} */ method2() {},
});
/** @private @export */ TestElementElement.prototype.method2;
/** @public @export */ TestElementElement.prototype.method1;
`)
}

// Java: PolymerPassTest.testExportsUniqueMethodsFromLegacyElementAndBehaviors()
func (c *PolymerPassTest) TestExportsUniqueMethodsFromLegacyElementAndBehaviors() {
	c.Test(
		`/** @polymerBehavior */
const Behavior1 = {
  /** @public */ onAll: function() {},
  /**
   * @public
// Note we include this @return annotation to test that we aren't including @return,
// @param and other redundant JSDoc in our generated @export statements, since that
// would cause a re-declaration error.
   * @return {void}
   */
   onBehavior1: function() {},
};
/** @polymerBehavior */
const Behavior2 = {
  /** @private */ onAll: function() {},
  /** @private */ onBehavior2: function() {},
};
Polymer({
  is: 'test-element',
  behaviors: [Behavior1, Behavior2],
  /** @private */ onAll: function() {},
  /** @private */ onElement: function() {},
});
`,
		`/**
 * @nocollapse
 * @polymerBehavior
 */
const Behavior1 = {
  /**
   * @public
   * @suppress {checkTypes,globalThis,visibility}
   */
  onAll: function() {},
  /**
   * @public
   * @return {void}
   * @suppress {checkTypes,globalThis,visibility}
   */
  onBehavior1: function() {},
};
/**
 * @nocollapse
 * @polymerBehavior
 */
const Behavior2 = {
  /**
   * @private
   * @suppress {checkTypes,globalThis,visibility}
   */
  onAll: function() {},
  /**
   * @private
   * @suppress {checkTypes,globalThis,visibility}
   */
  onBehavior2: function() {},
};
/**
 * @constructor
 * @extends {PolymerElement}
 * @implements {PolymerTestElementElementInterface$UID$0}
 */
var TestElementElement = function() {};
/**
 * @public
 * @return {void}
 */
TestElementElement.prototype.onBehavior1 = function() {};
/** @private */
TestElementElement.prototype.onBehavior2 = function() {};
Polymer(/** @lends {TestElementElement.prototype} */ {
  is: "test-element",
  behaviors: [Behavior1, Behavior2],
  /**
   * @private
   * @this {TestElementElement}
   */
  onAll: function() {},
  /**
   * @private
   * @this {TestElementElement}
   */
  onElement: function() {},
});
/** @private @export */
TestElementElement.prototype.onElement;
/** @private @export */
TestElementElement.prototype.onBehavior2;
/** @public @export */
TestElementElement.prototype.onBehavior1;
/** @private @export */
TestElementElement.prototype.onAll;
`)
}

// Java: PolymerPassTest.testSimpleObserverStringsConvertedToReferences1()
func (c *PolymerPassTest) TestSimpleObserverStringsConvertedToReferences1() {
	c.Test(
		`/** @constructor */
var User = function() {};
class XElement extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: {
        type: Array,
        observer: '_petsChanged'
      },
      name: {
        type: String,
        observer: '_nameChanged'
      },
      thingToDo: Function,
    };
  }
  _petsChanged(newValue, oldValue) {}
  _nameChanged(newValue, oldValue) {}
}
`,
		`/** @constructor */
var User = function() {};
/** @implements {PolymerXElementInterface$UID$0} */
class XElement extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(XElement, {
      user_: Object,
      pets: {
        type: Array,
        observer: XElement.prototype._petsChanged
      },
      name: {
        type: String,
        observer: XElement.prototype._nameChanged
      },
      thingToDo: Function,
    });
  }
  _petsChanged(newValue, oldValue) {}
  _nameChanged(newValue, oldValue) {}
}
/** @type {!User} @private */
XElement.prototype.user_;
/** @type {!Array} */
XElement.prototype.pets;
/** @type {string} */
XElement.prototype.name;
/** @type {!Function} */
XElement.prototype.thingToDo;
/** @export */
XElement.prototype._nameChanged;
/** @export */
XElement.prototype._petsChanged
`)
}

// Java: PolymerPassTest.testSimpleObserverStringsConvertedToReferences2()
func (c *PolymerPassTest) TestSimpleObserverStringsConvertedToReferences2() {
	c.Test(
		`/** @constructor */
var User = function() {};
var a = {};
a.B = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: {
        type: Array,
        observer: '_petsChanged'
      },
      name: {
        type: String,
        observer: '_nameChanged'
      },
      thingToDo: Function
    };
  }
  _petsChanged(newValue, oldValue) {}
  _nameChanged(newValue, oldValue) {}
};
`,
		`/** @constructor */
var User = function() {};
var a = {};
/** @implements {Polymera_BInterface$UID$0} */
a.B = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return  $jscomp.reflectObject(a.B, {
      user_: Object,
      pets: {
        type: Array,
        observer: a.B.prototype._petsChanged
      },
      name: {
        type: String,
        observer: a.B.prototype._nameChanged
      },
      thingToDo: Function
    });
  }
  _petsChanged(newValue, oldValue) {}
  _nameChanged(newValue, oldValue) {}
};
/** @type {!User} @private */
a.B.prototype.user_;
/** @type {!Array} */
a.B.prototype.pets;
/** @type {string} */
a.B.prototype.name;
/** @type {!Function} */
a.B.prototype.thingToDo;
/** @export */
a.B.prototype._nameChanged;
/** @export */
a.B.prototype._petsChanged
`)
}

// Java: PolymerPassTest.testSimpleObserverStringsConvertedToReferences3()
func (c *PolymerPassTest) TestSimpleObserverStringsConvertedToReferences3() {
	c.Test(
		`/** @constructor */
var User = function() {};
const A = class extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: {
        type: Array,
        observer: '_petsChanged'
      },
      name: {
        type: String,
        observer: '_nameChanged'
      },
      thingToDo: Function,
    };
  }
  _petsChanged(newValue, oldValue) {}
  _nameChanged(newValue, oldValue) {}
};
`,
		`/** @constructor */
var User = function() {};
/** @implements {PolymerAInterface$UID$0} */
const A = class extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(A, {
      user_: Object,
      pets: {
        type: Array,
        observer: A.prototype._petsChanged
      },
      name: {
        type: String,
        observer: A.prototype._nameChanged
      },
      thingToDo: Function,
    });
  }
  _petsChanged(newValue, oldValue) {}
  _nameChanged(newValue, oldValue) {}
};
/** @type {!User} @private */
A.prototype.user_;
/** @type {!Array} */
A.prototype.pets;
/** @type {string} */
A.prototype.name;
/** @type {!Function} */
A.prototype.thingToDo;
/** @export */
A.prototype._nameChanged;
/** @export */
A.prototype._petsChanged
`)
}

// Java: PolymerPassTest.testReflectionForComputedPropertyStrings1()
func (c *PolymerPassTest) TestReflectionForComputedPropertyStrings1() {
	c.Test(
		`/** @constructor */
var User = function() {};
class XElement extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: {
        type: Array,
        computed: '_computePets()'
      },
      name: {
        type: String,
        computed: '_computeName(user_, thingToDo)'
      },
      thingToDo: Function,
    };
  }
  _computePets() {}
  _computeName(user, thingToDo) {}
}
`,
		`/** @constructor */
var User = function() {};
/** @implements {PolymerXElementInterface$UID$0} */
class XElement extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(XElement, {
      user_: Object,
      pets: {
        type: Array,
        computed: '_computePets()'
      },
      name: {
        type: String,
        computed: '_computeName(user_, thingToDo)'
      },
      thingToDo: Function,
    });
  }
  _computePets() {}
  _computeName(user, thingToDo) {}
}
/** @type {!User} @private */
XElement.prototype.user_;
/** @type {!Array} */
XElement.prototype.pets;
/** @type {string} */
XElement.prototype.name;
/** @type {!Function} */
XElement.prototype.thingToDo;
/** @export */
XElement.prototype._computeName;
/** @export */
XElement.prototype._computePets;
`)
}

// Java: PolymerPassTest.testReflectionForComputedPropertyStrings2()
func (c *PolymerPassTest) TestReflectionForComputedPropertyStrings2() {
	c.Test(
		`/** @constructor */
var User = function() {};
/** @type {string} */ User.prototype.id;
class XElement extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: Array,
      name: {
        type: String,
        computed: '_computeName("user", 12.0, user_.id, user_.*)'
      },
      thingToDo: Function,
    };
  }
  _computePets() {}
  _computeName(user, thingToDo) {}
}
`,
		`/** @constructor */
var User = function() {};
/** @type {string} */ User.prototype.id;
/** @implements {PolymerXElementInterface$UID$0} */
class XElement extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(XElement, {
      user_: Object,
      pets: Array,
      name: {
        type: String,
        computed: '_computeName("user", 12.0, user_.id, user_.*)'
      },
      thingToDo: Function,
    });
  }
  _computePets() {}
  _computeName(user, thingToDo) {}
}
/** @type {!User} @private */
XElement.prototype.user_;
/** @type {!Array} */
XElement.prototype.pets;
/** @type {string} */
XElement.prototype.name;
/** @type {!Function} */
XElement.prototype.thingToDo;
/** @export */
XElement.prototype._computeName;
/** @export */
XElement.prototype._computePets;
`)
}

// Java: PolymerPassTest.testParseErrorForComputedPropertyStrings1()
func (c *PolymerPassTest) TestParseErrorForComputedPropertyStrings1() {
	c.CompilerTestCaseBase.TestError(
		`/** @constructor */
var User = function() {};
/** @type {string} */ User.prototype.id;
class XElement extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: Array,
      name: {
        type: String,
        computed: '_computeName("user", 12.0, user_.id'
      },
      thingToDo: Function,
    };
  }
  _computePets() {}
  _computeName(user, thingToDo) {}
}
`,
		PolymerPassErrors_POLYMER_UNPARSABLE_STRING)
}

// Java: PolymerPassTest.testParseErrorForComputedPropertyStrings2()
func (c *PolymerPassTest) TestParseErrorForComputedPropertyStrings2() {
	c.CompilerTestCaseBase.TestError(
		`/** @constructor */
var User = function() {};
/** @type {string} */ User.prototype.id;
class XElement extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: Array,
      name: {
        type: String,
        computed: '_computeName("user)'
      },
      thingToDo: Function,
    };
  }
  _computePets() {}
  _computeName(user, thingToDo) {}
}
`,
		PolymerPassErrors_POLYMER_UNPARSABLE_STRING)
}

// Java: PolymerPassTest.testReflectionForComplexObservers()
func (c *PolymerPassTest) TestReflectionForComplexObservers() {
	c.Test(
		`/** @constructor */
var User = function() {};
class XElement extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      /** @type {!User} @private */
      user_: Object,
      pets: Array,
      name: String,
      thingToDo: Function,
    };
  }
  static get observers() {
    return [
      '_userChanged(user_)',
      '_userOrThingToDoChanged(user_, thingToDo)'
    ];
  }
  _userChanged(user_) {}
  _userOrThingToDoChanged(user, thingToDo) {}
}
`,
		`/** @constructor */
var User = function() {};
/** @implements {PolymerXElementInterface$UID$0} */
class XElement extends Polymer.Element {
  /** @return {string} */
  static get is() { return 'x-element'; }
  /** @return {PolymerElementProperties} */
  static get properties() {
    return $jscomp.reflectObject(XElement, {
      user_: Object,
      pets: Array,
      name: String,
      thingToDo: Function,
    });
  }
  /** @return {!Array<string>} */
  static get observers() {
    return [
      '_userChanged(user_)',
      '_userOrThingToDoChanged(user_, thingToDo)'
    ];
  }
  _userChanged(user_) {}
  _userOrThingToDoChanged(user, thingToDo) {}
}
/** @type {!User} @private */
XElement.prototype.user_;
/** @type {!Array} */
XElement.prototype.pets;
/** @type {string} */
XElement.prototype.name;
/** @type {!Function} */
XElement.prototype.thingToDo;
/** @export */
XElement.prototype._userOrThingToDoChanged;
/** @export */
XElement.prototype._userChanged;
`)
}

// Java: PolymerPassTest.testPolymerInEsModuleExport()
func (c *PolymerPassTest) TestPolymerInEsModuleExport() {
	c.Test(
		`export let PaperMenuButton = Polymer({
  is: 'paper-menu-button',
  properties: {
    opened: {type: Boolean, value: false}
  }
});
`,
		`/** @constructor @extends {PolymerElement}
 * @implements {PolymerPaperMenuButtonInterface$UID$0}
 */
var PaperMenuButton = function() {};
/** @type {boolean} */
PaperMenuButton.prototype.opened;
PaperMenuButton = Polymer(/** @lends {PaperMenuButton.prototype} */ {
  is: 'paper-menu-button',
  properties: $jscomp.reflectObject(PaperMenuButton, {
    opened: {type: Boolean, value: false}
  })
});
export {PaperMenuButton};
`)
}

// Java: PolymerPassTest.testCanAccessPropertyOnPolymerElement_fromEsModuleTypeSummary()
func (c *PolymerPassTest) TestCanAccessPropertyOnPolymerElement_fromEsModuleTypeSummary() {
	c.Test(
		c.Externs(
			polymerPassTest_EXTERNS,
			`/** @typeSummary */
Polymer({
  is: 'foo',
  properties: {
    opened: {type: Boolean, value: false}
  },
  toggle: function() {}
})
export {}
`),
		c.Srcs(
			`function fn(/** !FooElement */ elem) {
  elem.toggle();
  return elem.opened;
}
`))
}

// Java: PolymerPassTest.testCanAccessPropertyOnPolymerElement_fromBundledGoogModuleTypeSummary()
func (c *PolymerPassTest) TestCanAccessPropertyOnPolymerElement_fromBundledGoogModuleTypeSummary() {
	c.Test(
		c.Externs(
			polymerPassTest_EXTERNS,
			`/** @typeSummary */
goog.loadModule(function(exports) {
goog.module('a.b.c');
Polymer({
  is: 'foo',
  properties: {
    opened: {type: Boolean, value: false}
  },
  toggle: function() {}
})
return exports;
});
`),
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`function fn(/** !FooElement */ elem) {
  elem.toggle();
  return elem.opened;
}
`))
}

// Java: PolymerPassTest.testCanAccessPropertyOnPolymerElement_fromBundledGoogModuleTypeSummaryBehavior()
func (c *PolymerPassTest) TestCanAccessPropertyOnPolymerElement_fromBundledGoogModuleTypeSummaryBehavior() {
	c.Test(
		c.Externs(
			polymerPassTest_EXTERNS,
			`/** @typeSummary */
goog.loadModule(function(exports) {
goog.module('SomeBehavior');
/** @polymerBehavior */
exports = {
  properties: {
    opened: {type: Boolean, value: false}
  },
  toggle: function() {}
};
return exports;
});
`,
			`/** @typeSummary */
goog.loadModule(function(exports) {
goog.module('a.b.c');
const SomeBehavior = goog.require('SomeBehavior');
Polymer({
  is: 'foo',
  behaviors: [SomeBehavior],
  properties: {
    opened: {type: Boolean, value: false}
  },
  toggle: function() {}
})
return exports;
});
`),
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`function fn(/** !FooElement */ elem) {
  elem.toggle();
  return elem.opened;
}
`))
}

// Java: PolymerPassTest.testPolymerImplicitGlobalNamingConflict()
func (c *PolymerPassTest) TestPolymerImplicitGlobalNamingConflict() {
	c.TestError(
		`var FooElement;
Polymer({is: 'foo'})
`,
		PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT)

	c.Test(
		c.Srcs(
			"var FooElement;",
			`goog.module('m');
Polymer({is: 'foo'})
`),
		c.Error(PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT))

	c.Test(
		c.Srcs(
			// First file in compilation needs to be a script, not a module, as the compiler injects
			// global declarations into the first file.
			"",
			`goog.module('m');
var FooElement;
Polymer({is: 'foo'})
`),
		c.Error(PolymerClassRewriter_IMPLICIT_GLOBAL_CONFLICT))

	c.TestNoWarning("var FooElement = Polymer({is: 'foo'});")
}

// Java: PolymerPassTest.testShadowPolymerElement()
func (c *PolymerPassTest) TestShadowPolymerElement() {
	c.TestError(
		`goog.module('m');
const Foo = Polymer({is: 'foo'});
var PolymerElement;
`,
		PolymerClassRewriter_POLYMER_ELEMENT_CONFLICT)
}

// Java: PolymerPassTest.testPolyerObjectPropertiesAccessedBeforeDefinition()
func (c *PolymerPassTest) TestPolyerObjectPropertiesAccessedBeforeDefinition() {
	c.Test(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`/** @polymerBehavior */
var PtGaUXBehavior = {
 _abc: function() {
    if (true) {}
 }
}
function func() {
 if (true) {
   Polymer({
     is: 'pt-test-component',
     behaviors: [PtGaUXBehavior],
     properties: {disabled: Boolean}
   });
 }
}
`),
		c.Expected(
			testExternsBuilder_GetClosureExternsAsSource(),
			"/**\n * @constructor\n * @extends {PolymerElement}\n * @implements {PolymerPtTestComponentElementInterface$m1176578414$0}\n */\n// This `var PtTestComponentElement = ...` is placed before any\n// `PtTestComponentElement.prototype.<...>` calls. We don't want to access\n// properties on `PtTestComponentElement` before it is defined.\nvar PtTestComponentElement = function() {\n};\n/**\n * @polymerBehavior\n * @nocollapse\n * @polymerBehavior\n */\nvar PtGaUXBehavior = {/**\n * @suppress {checkTypes,globalThis,visibility}\n */\n_abc:function() {\n  if (true) {\n  }\n}};\nfunction func() {\n  if (true) {\n    Polymer(/** @lends {PtTestComponentElement.prototype} */\n    {is:\"pt-test-component\", behaviors:[PtGaUXBehavior], properties:$jscomp.reflectObject(PtTestComponentElement, {disabled:Boolean})});\n    /** @export */\n    PtTestComponentElement.prototype._abc;\n  }\n}\n/** @type {boolean} */\nPtTestComponentElement.prototype.disabled;\nPtTestComponentElement.prototype._abc = function() {\n  if (true) {\n  }\n};\n"))
}

// Java: PolymerPassTest.testPolymerPropsAllExported_updatesNonPrivateProps()
func (c *PolymerPassTest) TestPolymerPropsAllExported_updatesNonPrivateProps() {
	c.Test(
		`class C extends Polymer.Element {
  static get is() { return 'x-element'; }

  static get properties() {
    return {
      explicitProp: {
        type: Number,
      },
    };
  }

  constructor() {
    super();
    this.explicitProp = 1;
    this.implicitProp = 2;
    // Note: This will also be detected as a method.
    this.implicitMethod = () => {};
    this.privateProp = 3;
  }
}
if (false) {
  /** @public @type {number} */
  C.prototype.explicitProp;
  /** @public @type {number} */
  C.prototype.implicitProp;
  /** @public @type {function()} */
  C.prototype.implicitMethod;
  /** @private @type {number} */
  C.prototype.privateProp;
}
`,
		`/**
 * @implements {PolymerCInterface$UID$0}
 */
class C extends Polymer.Element {
  /**
   * @return {string}
   */
  static get is() {
    return "x-element";
  }
  /**
   * @return {PolymerElementProperties}
   */
  static get properties() {
    return $jscomp.reflectObject(C, {explicitProp:{type:Number,},});
  }
  constructor() {
    super();
    this.explicitProp = 1;
    this.implicitProp = 2;
    this.implicitMethod = () => {};
    this.privateProp = 3;
  }
}
/** @type {number} */
C.prototype.explicitProp;
/** @export */
C.prototype.implicitMethod;
/** @export */
C.prototype.constructor;
if (false) {
  /** @export @type {number} */
  C.prototype.explicitProp;
  /** @export @type {number} */
  C.prototype.implicitProp;
  /** @export @type {function()} */
  C.prototype.implicitMethod;
  /** @private @type {number} */
  C.prototype.privateProp;
}
`)
}

// Java: PolymerPassTest.testPolymerPropsAllExported_ignoresOtherClass()
func (c *PolymerPassTest) TestPolymerPropsAllExported_ignoresOtherClass() {
	c.Test(
		`class C extends Polymer.Element {
  static get is() { return 'x-element'; }

  constructor() {
    super();
    this.explicitProp = 1;
  }
}
if (false) {
  /** @public @type {number} */
  C.prototype.explicitProp;
}

// Not a Polymer element
class D {
  constructor() {
    this.explicitProp = 1;
  }
}
if (false) {
  /** @public @type {number} */
  D.prototype.explicitProp;
}
`,
		`/**
 * @implements {PolymerCInterface$UID$0}
 */
class C extends Polymer.Element {
  /**
   * @return {string}
   */
  static get is() {
    return "x-element";
  }
  constructor() {
    super();
    this.explicitProp = 1;
  }
}
/** @export */
C.prototype.constructor;
if (false) {
  /** @export @type {number} */
  C.prototype.explicitProp;
}
class D {
  constructor() {
    this.explicitProp = 1;
  }
}
if (false) {
  /** @public @type {number} */
  D.prototype.explicitProp;
}
`)
}

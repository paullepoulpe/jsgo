// Port of com/google/javascript/jscomp/PolymerPassSuppressBehaviorsAndProtectKeysTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/PolymerPassSuppressBehaviorsAndProtectKeysTest.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var polymerPassSuppressBehaviorsAndProtectKeysTest_EXTERNS = CompilerTestCase_MINIMAL_EXTERNS +
	`/** @constructor */
var HTMLElement = function() {};
/** @constructor @extends {HTMLElement} */
var HTMLInputElement = function() {};
/** @constructor @extends {HTMLElement} */
var PolymerElement = function() {};
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
 * @return {!function()}
 */
var Polymer = function(a) {};
var alert = function(msg) {};
`

// PolymerPassSuppressBehaviorsAndProtectKeysTest is the unit tests for
// PolymerPassSuppressBehaviorsAndProtectKeys
type PolymerPassSuppressBehaviorsAndProtectKeysTest struct {
	CompilerTestCaseBase
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeysTest()
func NewPolymerPassSuppressBehaviorsAndProtectKeysTest() *PolymerPassSuppressBehaviorsAndProtectKeysTest {
	c := &PolymerPassSuppressBehaviorsAndProtectKeysTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, polymerPassSuppressBehaviorsAndProtectKeysTest_EXTERNS)
	return c
}

func TestPolymerPassSuppressBehaviorsAndProtectKeysTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPolymerPassSuppressBehaviorsAndProtectKeysTest() })
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeysTest.getProcessor(Compiler)
func (c *PolymerPassSuppressBehaviorsAndProtectKeysTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		NodeTraversal_Traverse(compiler, root, polymerPassSuppressBehaviorsAndProtectKeys{compiler: compiler})
	})
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeysTest.setUp()
func (c *PolymerPassSuppressBehaviorsAndProtectKeysTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.AllowExternsChanges()
	c.EnableRunTypeCheckAfterProcessing()
	c.EnableParseTypeInfo()
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeysTest.testPropertyTypeRemoval()
func (c *PolymerPassSuppressBehaviorsAndProtectKeysTest) TestPropertyTypeRemoval() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    /** @type {boolean} */
    isFun: {
      type: Boolean,
      value: true,
    },
    /** @type {!Element} */
    funObject: {
      type: Object,
    },
    /** @type {!Array<String>} */
    funArray: {
      type: Array,
    },
  },
};
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
    },
    funArray: {
      type: Array,
    },
  },
};
`)
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeysTest.testDefaultValueSuppression()
func (c *PolymerPassSuppressBehaviorsAndProtectKeysTest) TestDefaultValueSuppression() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
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
`)
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeysTest.testConstBehaviours()
func (c *PolymerPassSuppressBehaviorsAndProtectKeysTest) TestConstBehaviours() {
	c.DisableTypeCheck()

	c.Test(
		`/** @polymerBehavior */
const FunBehavior = {
};
`,
		`/** @polymerBehavior @nocollapse */
const FunBehavior = {
};
`)
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeysTest.testLetBehaviours()
func (c *PolymerPassSuppressBehaviorsAndProtectKeysTest) TestLetBehaviours() {
	c.DisableTypeCheck()

	c.Test(
		`/** @polymerBehavior */
let FunBehavior = {
};
`,
		`/** @polymerBehavior @nocollapse */
let FunBehavior = {
};
`)
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeysTest.testPolymerBehaviorWithoutRhsDoesntCrash()
func (c *PolymerPassSuppressBehaviorsAndProtectKeysTest) TestPolymerBehaviorWithoutRhsDoesntCrash() {
	c.DisableTypeCheck()

	c.TestError(
		"/** @polymerBehavior */ let FunBehavior;", PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR)
}

// Java: PolymerPassSuppressBehaviorsAndProtectKeysTest.testMethodsInPropertiesAndObserversAreProtected()
func (c *PolymerPassSuppressBehaviorsAndProtectKeysTest) TestMethodsInPropertiesAndObserversAreProtected() {
	c.Test(
		`/** @polymerBehavior */
var FunBehavior = {
  properties: {
    computedProp: {
      computed: 'quotedProp(prop1, prop2)',
    },
    observedProp: {
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
`,
		`/**
 * @polymerBehavior
 * @nocollapse
 * @polymerBehavior
 */
var FunBehavior = {
  properties: {
    computedProp: {
      computed: 'quotedProp(prop1, prop2)',
    },
    observedProp: {
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
`)
}

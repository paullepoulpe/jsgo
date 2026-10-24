// Port of com/google/javascript/jscomp/PolymerClassDefinitionTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/PolymerClassDefinitionTest.java)

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

// PolymerClassDefinitionTest is the unit test for PolymerClassDefinition.
type PolymerClassDefinitionTest struct {
	CompilerTypeTestCaseBase

	polymerCall rhino.Node
}

// Java: PolymerClassDefinitionTest()
func NewPolymerClassDefinitionTest() *PolymerClassDefinitionTest {
	c := &PolymerClassDefinitionTest{}
	c.InitCompilerTypeTestCaseBase(c)
	return c
}

func TestPolymerClassDefinitionTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPolymerClassDefinitionTest() })
}

// Java: PolymerClassDefinitionTest.setUp()
func (c *PolymerClassDefinitionTest) SetUp() {
	c.CompilerTypeTestCaseBase.SetUp()
	c.polymerCall = rhino.NilNode
}

// TODO(jlklein): Add more complex test cases and verify behaviors and descriptors.

// Java: PolymerClassDefinitionTest.testSimpleBehavior()
func (c *PolymerClassDefinitionTest) TestSimpleBehavior() {
	def := c.parseAndExtractClassDefFromCall(
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

  /** @param {string} funAmount */
  doSomethingFun: function(funAmount) {
    alert('Something ' + funAmount + ' fun!');
  },

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
`)

	AssertThat(def).IsNotNull()
	AssertThat(def.defType).IsEqualTo(PolymerClassDefinition_DefinitionType_ObjectLiteral)
	AssertNode(def.target).HasType(rhino.Token_NAME)
	AssertThat(def.target.GetString()).IsEqualTo("A")
	AssertThat(def.nativeBaseElement).IsNull()
	AssertThat(def.behaviors).HasSize(1)
	AssertThat(def.props).HasSize(2)
	AssertThat(def.props[0].name.GetString()).IsEqualTo("pets")
	AssertThat(def.props[1].name.GetString()).IsEqualTo("name")
	AssertThat(def.behaviorProps).HasSize(1) // 'isFun'
	AssertThat(def.behaviors[0].props[0].name.GetString()).IsEqualTo("isFun")
}

// Java: PolymerClassDefinitionTest.testBasicClass()
func (c *PolymerClassDefinitionTest) TestBasicClass() {
	def := c.parseAndExtractClassDefFromClass(
		`class A extends Polymer.Element {
  static get is() { return 'x-element'; }
  static get properties() {
    return {
      pets: {
        type: Array,
        notify: true,
      },
      name: String
    };
  }
}
`)

	AssertThat(def).IsNotNull()
	AssertThat(def.defType).IsEqualTo(PolymerClassDefinition_DefinitionType_ES6Class)
	AssertNode(def.target).HasType(rhino.Token_NAME)
	AssertThat(def.target.GetString()).IsEqualTo("A")
	AssertThat(def.descriptor).IsNotNull()
	AssertNode(def.descriptor).HasType(rhino.Token_OBJECTLIT)
	AssertThat(def.nativeBaseElement).IsNull()
	AssertThat(def.behaviors).IsNull()
	AssertThat(def.props).HasSize(2)
}

// Java: PolymerClassDefinitionTest.testClassMethods()
func (c *PolymerClassDefinitionTest) TestClassMethods() {
	def := c.parseAndExtractClassDefFromClass(
		`class A extends Polymer.Element {
  static get is() { return 'x-element'; }
  method() {}
  /** @type {function()} */
  fieldWithFunction = function() {};
  fieldWithArrowFunction = () => {};
  fieldWithNumber = 1;
  fieldWithoutInitializer;
  constructor() {
    super();
    this.propertyWithFunction = function() {};
    this.propertyWithArrowFunction = () => {};
    this.propertyWithNumber = 2;
  }
}
`)

	AssertThat(def).IsNotNull()
	var methodNames []string
	for _, m := range def.methods {
		methodNames = append(methodNames, m.name.GetString())
	}
	AssertThat(methodNames).
		ContainsExactly(
			"method",
			"fieldWithFunction",
			"fieldWithArrowFunction",
			"constructor",
			"propertyWithFunction",
			"propertyWithArrowFunction")
}

// Java: PolymerClassDefinitionTest.testDynamicDescriptor()
func (c *PolymerClassDefinitionTest) TestDynamicDescriptor() {
	def := c.parseAndExtractClassDefFromCall(
		`var A = Polymer({
  is: x,
});
`)

	AssertThat(def.target.GetString()).IsEqualTo("A")
}

// Java: PolymerClassDefinitionTest.testDynamicDescriptor1()
func (c *PolymerClassDefinitionTest) TestDynamicDescriptor1() {
	def := c.parseAndExtractClassDefFromCall(
		`Polymer({
  is: x,
});
`)

	AssertThat(def.target.GetString()).IsEqualTo("XElement")
}

// Java: PolymerClassDefinitionTest.testDynamicDescriptor2()
func (c *PolymerClassDefinitionTest) TestDynamicDescriptor2() {
	def := c.parseAndExtractClassDefFromCall(
		`Polymer({
  is: foo.bar,
});
`)

	AssertThat(def.target.GetString()).IsEqualTo("Foo$barElement")
}

// Java: PolymerClassDefinitionTest.testDynamicDescriptor3()
func (c *PolymerClassDefinitionTest) TestDynamicDescriptor3() {
	def := c.parseAndExtractClassDefFromCall(
		`Polymer({
  is: this.bar,
});
`)

	AssertThat(def.target.GetString()).IsEqualTo("This$barElement")
}

// Java: PolymerClassDefinitionTest.parseAndExtractClassDefFromCall(String)
func (c *PolymerClassDefinitionTest) parseAndExtractClassDefFromCall(code string) *PolymerClassDefinition {
	rootNode := c.compiler.ParseTestCode(code)
	globalNamespace := NewGlobalNamespace(c.compiler, rootNode)

	NodeUtil_VisitPostOrder(
		rootNode,
		NodeUtil_VisitorFunc(func(node rhino.Node) {
			if PolymerPassStaticUtils_IsPolymerCall(node) {
				c.polymerCall = node
			}
		}))

	AssertThat(c.polymerCall).IsNotNull()
	return PolymerClassDefinition_ExtractFromCallNode(
		c.polymerCall,
		c.compiler,
		/* moduleMetadata= */ nil,
		NewPolymerBehaviorExtractor(
			c.compiler,
			globalNamespace,
			ModuleMetadataMap_EmptyForTesting(),
			ModuleMap_EmptyForTesting()))
}

// Java: PolymerClassDefinitionTest.parseAndExtractClassDefFromClass(String)
func (c *PolymerClassDefinitionTest) parseAndExtractClassDefFromClass(code string) *PolymerClassDefinition {
	rootNode := c.compiler.ParseTestCode(code)

	NodeUtil_VisitPostOrder(
		rootNode,
		NodeUtil_VisitorFunc(func(node rhino.Node) {
			if PolymerPassStaticUtils_IsPolymerClass(node) {
				c.polymerCall = node
			}
		}))

	AssertThat(c.polymerCall).IsNotNull()
	return PolymerClassDefinition_ExtractFromClassNode(c.polymerCall, c.compiler)
}

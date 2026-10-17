// Port of com/google/javascript/jscomp/PolymerPassFindExternsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/PolymerPassFindExternsTest.java)

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

var polymerPassFindExternsTest_EXTERNS = CompilerTestCase_MINIMAL_EXTERNS +
	"/** @constructor */\n" +
	"var HTMLElement = function() {};\n" +
	"/** @constructor @extends {HTMLElement} */\n" +
	"var HTMLInputElement = function() {};\n" +
	"/** @constructor @extends {HTMLElement} */\n" +
	"var PolymerElement = function() {};\n" +
	"/** @type {!Object} */\n" +
	"PolymerElement.prototype.$;\n" +
	"PolymerElement.prototype.created = function() {};\n" +
	"PolymerElement.prototype.ready = function() {};\n" +
	"PolymerElement.prototype.attached = function() {};\n" +
	"PolymerElement.prototype.domReady = function() {};\n" +
	"PolymerElement.prototype.detached = function() {};\n" +
	"/**\n" +
	" * Call the callback after a timeout. Calling job again with the same name\n" +
	" * resets the timer but will not result in additional calls to callback.\n" +
	" *\n" +
	" * @param {string} name\n" +
	" * @param {Function} callback\n" +
	" * @param {number} timeoutMillis The minimum delay in milliseconds before\n" +
	" *     calling the callback.\n" +
	" */\n" +
	"PolymerElement.prototype.job = function(name, callback, timeoutMillis) {};\n" +
	"/**\n" +
	" * @param a {!Object}\n" +
	" * @return {!function()}\n" +
	" */\n" +
	"var Polymer = function(a) {};\n" +
	"var alert = function(msg) {};\n"

type PolymerPassFindExternsTest struct {
	CompilerTestCaseBase

	findExternsCallback *PolymerPassFindExterns
}

// Java: PolymerPassFindExternsTest()
func NewPolymerPassFindExternsTest() *PolymerPassFindExternsTest {
	c := &PolymerPassFindExternsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, polymerPassFindExternsTest_EXTERNS)
	return c
}

func TestPolymerPassFindExternsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPolymerPassFindExternsTest() })
}

// polymerPassFindExternsTest_pass is the anonymous CompilerPass of
// getProcessor.
type polymerPassFindExternsTest_pass struct {
	compiler            *Compiler
	findExternsCallback *PolymerPassFindExterns
}

func (p *polymerPassFindExternsTest_pass) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(p.compiler, externs, p.findExternsCallback)
}

// Java: PolymerPassFindExternsTest.getProcessor(Compiler)
func (c *PolymerPassFindExternsTest) GetProcessor(compiler *Compiler) CompilerPass {
	c.findExternsCallback = NewPolymerPassFindExterns()
	return &polymerPassFindExternsTest_pass{compiler: compiler, findExternsCallback: c.findExternsCallback}
}

// Java: PolymerPassFindExternsTest.setUp()
func (c *PolymerPassFindExternsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.AllowExternsChanges()
	c.EnableRunTypeCheckAfterProcessing()
	c.EnableParseTypeInfo()
}

// Java: PolymerPassFindExternsTest.testFindsPolymerElementRoot()
func (c *PolymerPassFindExternsTest) TestFindsPolymerElementRoot() {
	c.TestSame("")
	polymerElementNode := c.findExternsCallback.GetPolymerElementExterns()

	AssertThat(polymerElementNode).IsNotNull()
	AssertThat(polymerElementNode.IsVar()).IsTrue()
	AssertThat(polymerElementNode.GetFirstChild().MatchesQualifiedName("PolymerElement")).IsTrue()
}

// Java: PolymerPassFindExternsTest.testFindsPolymerElementRoot_inTypeSummary()
func (c *PolymerPassFindExternsTest) TestFindsPolymerElementRoot_inTypeSummary() {
	c.TestSame(c.Externs("/** @typeSummary */\n"+polymerPassFindExternsTest_EXTERNS), c.Srcs(""))
	polymerElementNode := c.findExternsCallback.GetPolymerElementExterns()

	AssertThat(polymerElementNode).IsNotNull()
	AssertThat(polymerElementNode.IsVar()).IsTrue()
	AssertThat(polymerElementNode.GetFirstChild().MatchesQualifiedName("PolymerElement")).IsTrue()
	AssertThat(polymerElementNode.GetParent().IsScript()).IsTrue()
	AssertThat(NodeUtil_IsFromTypeSummary(polymerElementNode.GetParent())).IsTrue()
}

// Java: PolymerPassFindExternsTest.testFindsPolymerElementRoot_inTypeSummary_ignoresModule()
func (c *PolymerPassFindExternsTest) TestFindsPolymerElementRoot_inTypeSummary_ignoresModule() {
	c.TestSame(
		c.Externs(
			"goog.module = function(id) { };",
			"/** @typeSummary */\n"+
				"goog.module('some.module');\n"+
				polymerPassFindExternsTest_EXTERNS),
		c.Srcs(""))
	polymerElementNode := c.findExternsCallback.GetPolymerElementExterns()

	AssertThat(polymerElementNode).IsNull()
}

// Java: PolymerPassFindExternsTest.testFindsPolymerElementProps()
func (c *PolymerPassFindExternsTest) TestFindsPolymerElementProps() {
	c.TestSame("")
	expectedProps :=
		[]string{"$", "created", "ready", "attached", "domReady", "detached", "job"}
	polymerElementProps := c.findExternsCallback.GetPolymerElementProps()

	AssertThat(polymerElementProps).IsNotNull()
	AssertThat(polymerElementProps).HasSize(len(expectedProps))
	for i := 0; i < len(polymerElementProps); i++ {
		AssertThat(c.getPropertyName(polymerElementProps[i])).IsEqualTo(expectedProps[i])
	}
}

// Java: PolymerPassFindExternsTest.getPropertyName(Node)
func (c *PolymerPassFindExternsTest) getPropertyName(node *rhino.Node) string {
	node = node.GetFirstChild()

	rvalue := node.GetSecondChild()
	if rvalue != nil && rvalue.IsFunction() && !node.IsGetProp() {
		node = node.GetFirstChild()
	}

	return node.GetString()
}

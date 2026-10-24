// Port of com/google/javascript/jscomp/GatherRawExportsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/GatherRawExportsTest.java)

/*
 * Copyright 2009 The Closure Compiler Authors.
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
)

// GatherRawExportsTest holds tests for GatherRawExports.
type GatherRawExportsTest struct {
	CompilerTestCaseBase

	last *GatherRawExports
}

const gatherRawExportsTest_EXTERNS = "var window;var self;"

// Java: GatherRawExportsTest()
func NewGatherRawExportsTest() *GatherRawExportsTest {
	c := &GatherRawExportsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, gatherRawExportsTest_EXTERNS)
	return c
}

func TestGatherRawExportsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewGatherRawExportsTest() })
}

// Java: GatherRawExportsTest.setUp()
func (c *GatherRawExportsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
}

// Java: GatherRawExportsTest.getProcessor(Compiler)
func (c *GatherRawExportsTest) GetProcessor(compiler *Compiler) CompilerPass {
	c.last = NewGatherRawExports(compiler)
	return c.last
}

// Java: GatherRawExportsTest.assertExported(String, String...)
func (c *GatherRawExportsTest) AssertExported(js string, names ...string) {
	setNames := jutil.NewLinkedHashSet(names...)
	c.TestSame(js)
	AssertThat(c.last.GetExportedVariableNames()).IsEqualTo(setNames)
}

// Java: GatherRawExportsTest.testTopLevelVar()
func (c *GatherRawExportsTest) TestTopLevelVar() {
	// Global vars are not what we are looking for.
	c.AssertExported("var a")
}

// Java: GatherRawExportsTest.testExportsFoundQuoted()
func (c *GatherRawExportsTest) TestExportsFoundQuoted() {
	c.AssertExported("window['a']", "a")
}

// Java: GatherRawExportsTest.testExportsFoundUnquoted()
func (c *GatherRawExportsTest) TestExportsFoundUnquoted() {
	c.AssertExported("window.a", "a")
}

// Java: GatherRawExportsTest.testExportWithOptChainFoundQuoted()
func (c *GatherRawExportsTest) TestExportWithOptChainFoundQuoted() {
	c.AssertExported("window?.['a']", "a")
}

// Java: GatherRawExportsTest.testExportWithOptChainFoundUnquoted()
func (c *GatherRawExportsTest) TestExportWithOptChainFoundUnquoted() {
	c.AssertExported("window?.a", "a")
}

// Java: GatherRawExportsTest.testExportsFound4()
func (c *GatherRawExportsTest) TestExportsFound4() {
	c.AssertExported("this['a']", "a")
}

// Java: GatherRawExportsTest.testExportsFound5()
func (c *GatherRawExportsTest) TestExportsFound5() {
	c.AssertExported("this.a", "a")
}

// Java: GatherRawExportsTest.testExportsFound6()
func (c *GatherRawExportsTest) TestExportsFound6() {
	c.AssertExported("function f() { this['a'] }")
}

// Java: GatherRawExportsTest.testExportsFound7()
func (c *GatherRawExportsTest) TestExportsFound7() {
	c.AssertExported("function f() { this.a }")
}

// Java: GatherRawExportsTest.testExportsFound8()
func (c *GatherRawExportsTest) TestExportsFound8() {
	c.AssertExported("window['foo']", "foo")
}

// Java: GatherRawExportsTest.testExportsFound9()
func (c *GatherRawExportsTest) TestExportsFound9() {
	c.AssertExported("window['a'] = 1;", "a")
}

// Java: GatherRawExportsTest.testExportsFound10()
func (c *GatherRawExportsTest) TestExportsFound10() {
	c.AssertExported("window['a']['b']['c'] = 1;", "a")
}

// Java: GatherRawExportsTest.testExportsFound11()
func (c *GatherRawExportsTest) TestExportsFound11() {
	c.AssertExported("if (window['a'] = 1) alert(x);", "a")
}

// Java: GatherRawExportsTest.testExportsFound12()
func (c *GatherRawExportsTest) TestExportsFound12() {
	c.AssertExported("function foo() { window['a'] = 1; }", "a")
}

// Java: GatherRawExportsTest.testExportsFound13()
func (c *GatherRawExportsTest) TestExportsFound13() {
	c.AssertExported("function foo() {var window; window['a'] = 1; }")
}

// Java: GatherRawExportsTest.testExportsFound14()
func (c *GatherRawExportsTest) TestExportsFound14() {
	c.AssertExported("var a={window:{}}; a.window['b']")
}

// Java: GatherRawExportsTest.testExportsFound15()
func (c *GatherRawExportsTest) TestExportsFound15() {
	c.AssertExported("window.window['b']", "window")
}

// Java: GatherRawExportsTest.testExportsFound16()
func (c *GatherRawExportsTest) TestExportsFound16() {
	// It would be nice to handle this case, hopefully inlining will take care
	// of it for us.
	c.AssertExported("var a = window; a['b']")
}

// Java: GatherRawExportsTest.testExportsFound17()
func (c *GatherRawExportsTest) TestExportsFound17() {
	// Gather "this" reference in a global if block.
	c.AssertExported("if (true) { this.a }", "a")
	// Does not gather "this" reference in a local if block.
	c.AssertExported("function f() { if (true) { this.a } }")
}

// Java: GatherRawExportsTest.testExportOnTopFound1()
func (c *GatherRawExportsTest) TestExportOnTopFound1() {
	c.AssertExported("top['a']", "a")
}

// Java: GatherRawExportsTest.testExportOntopFound2()
func (c *GatherRawExportsTest) TestExportOntopFound2() {
	c.AssertExported("top.a", "a")
}

// Java: GatherRawExportsTest.testExportOnGlobalThis1()
func (c *GatherRawExportsTest) TestExportOnGlobalThis1() {
	c.AssertExported("globalThis['a']", "a")
}

// Java: GatherRawExportsTest.testExportOnGlobalThis2()
func (c *GatherRawExportsTest) TestExportOnGlobalThis2() {
	c.AssertExported("globalThis.a", "a")
}

// Java: GatherRawExportsTest.testExportOnPolyfillGlobal1()
func (c *GatherRawExportsTest) TestExportOnPolyfillGlobal1() {
	c.AssertExported("$jscomp.global['a']", "a")
}

// Java: GatherRawExportsTest.testExportOnPolyfillGlobal2()
func (c *GatherRawExportsTest) TestExportOnPolyfillGlobal2() {
	c.AssertExported("$jscomp.global.a", "a")
}

// Java: GatherRawExportsTest.testExportOnPolyfillGlobal1Collapsed()
func (c *GatherRawExportsTest) TestExportOnPolyfillGlobal1Collapsed() {
	c.AssertExported("$jscomp$global['a']", "a")
}

// Java: GatherRawExportsTest.testExportOnPolyfillGlobal2Collapsed()
func (c *GatherRawExportsTest) TestExportOnPolyfillGlobal2Collapsed() {
	c.AssertExported("$jscomp$global.a", "a")
}

// Java: GatherRawExportsTest.testExportOnSelf1()
func (c *GatherRawExportsTest) TestExportOnSelf1() {
	c.AssertExported("self['a']", "a")
}

// Java: GatherRawExportsTest.testExportOnSelf2()
func (c *GatherRawExportsTest) TestExportOnSelf2() {
	c.AssertExported("self.a", "a")
}

// Java: GatherRawExportsTest.testNoExportOnLocalSelf1()
func (c *GatherRawExportsTest) TestNoExportOnLocalSelf1() {
	c.AssertExported("function fn(self) { self['a']; }")
}

// Java: GatherRawExportsTest.testNoExportOnLocalSelf2()
func (c *GatherRawExportsTest) TestNoExportOnLocalSelf2() {
	c.AssertExported("function fn(self) { self.a; }")
}

// Java: GatherRawExportsTest.testExportOnGoogGlobalFound1()
func (c *GatherRawExportsTest) TestExportOnGoogGlobalFound1() {
	c.AssertExported("goog.global['a']", "a")
}

// Java: GatherRawExportsTest.testExportOnGoogGlobalFound2()
func (c *GatherRawExportsTest) TestExportOnGoogGlobalFound2() {
	c.AssertExported("goog.global.a", "a")
}

// Java: GatherRawExportsTest.testExportOnGoogGlobalFound3()
func (c *GatherRawExportsTest) TestExportOnGoogGlobalFound3() {
	c.AssertExported("goog$global['a']", "a")
}

// Java: GatherRawExportsTest.testExportOnGoogGlobalFound4()
func (c *GatherRawExportsTest) TestExportOnGoogGlobalFound4() {
	c.AssertExported("goog$global.a", "a")
}

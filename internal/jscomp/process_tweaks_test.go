// Port of com/google/javascript/jscomp/ProcessTweaksTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ProcessTweaksTest.java)

/*
 * Copyright 2010 The Closure Compiler Authors.
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

// ProcessTweaksTest: @author agrieve@google.com (Andrew Grieve)
type ProcessTweaksTest struct {
	CompilerTestCaseBase

	defaultValueOverrides *jutil.LinkedHashMap[string, *rhino.Node]
	stripTweaks           bool
}

// Java: ProcessTweaksTest()
func NewProcessTweaksTest() *ProcessTweaksTest {
	c := &ProcessTweaksTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, "function alert(arg) {}")
	return c
}

func TestProcessTweaksTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewProcessTweaksTest() })
}

// Java: ProcessTweaksTest.setUp()
func (c *ProcessTweaksTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	c.defaultValueOverrides = jutil.NewLinkedHashMap[string, *rhino.Node]()
	c.stripTweaks = false
}

// processTweaksTest_pass is the anonymous CompilerPass of getProcessor.
type processTweaksTest_pass struct {
	outer    *ProcessTweaksTest
	compiler *Compiler
}

// Java: ProcessTweaksTest.getProcessor(Compiler).new CompilerPass() {...}.process(Node, Node)
func (p *processTweaksTest_pass) Process(externs *rhino.Node, root *rhino.Node) {
	processTweak := NewProcessTweaks(p.compiler, p.outer.stripTweaks)
	processTweak.Process(externs, root)

	if p.outer.stripTweaks {
		emptySet := jutil.NewLinkedHashSet[string]()
		stripCode :=
			NewStripCode(
				p.compiler, emptySet, emptySet, emptySet /* enableTweakStripping= */, true)
		stripCode.Process(externs, root)
	}
}

// Java: ProcessTweaksTest.getProcessor(Compiler)
func (c *ProcessTweaksTest) GetProcessor(compiler *Compiler) CompilerPass {
	return &processTweaksTest_pass{outer: c, compiler: compiler}
}

// Java: ProcessTweaksTest.testBasicTweak1()
func (c *ProcessTweaksTest) TestBasicTweak1() {
	c.TestSame(
		`goog.tweak.registerBoolean('Foo', 'Description');
goog.tweak.getBoolean('Foo')
`)
}

// Java: ProcessTweaksTest.testBasicTweak2()
func (c *ProcessTweaksTest) TestBasicTweak2() {
	c.TestSame(
		`goog.tweak.registerString('Foo', 'Description');
goog.tweak.getString('Foo')
`)
}

// Java: ProcessTweaksTest.testBasicTweak3()
func (c *ProcessTweaksTest) TestBasicTweak3() {
	c.TestSame(
		`goog.tweak.registerNumber('Foo', 'Description');
goog.tweak.getNumber('Foo')
`)
}

// Java: ProcessTweaksTest.testBasicTweak4()
func (c *ProcessTweaksTest) TestBasicTweak4() {
	c.TestSame("goog.tweak.registerButton('Foo', 'Description', function() {})")
}

// Java: ProcessTweaksTest.testBasicTweak5()
func (c *ProcessTweaksTest) TestBasicTweak5() {
	c.TestSame(
		`goog.tweak.registerBoolean('A.b_7', 'Description', true,
{ requiresRestart:false })
`)
}

// Java: ProcessTweaksTest.testBasicTweak6()
func (c *ProcessTweaksTest) TestBasicTweak6() {
	c.TestSame(
		`var opts = { requiresRestart:false };
goog.tweak.registerBoolean('Foo', 'Description', true, opts)
`)
}

// Java: ProcessTweaksTest.testNonLiteralId1()
func (c *ProcessTweaksTest) TestNonLiteralId1() {
	c.TestError(
		"goog.tweak.registerBoolean(3, 'Description')", ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR)
}

// Java: ProcessTweaksTest.testNonLiteralId2()
func (c *ProcessTweaksTest) TestNonLiteralId2() {
	c.TestError("goog.tweak.getBoolean('a' + 'b')", ProcessTweaks_NON_LITERAL_TWEAK_ID_ERROR)
}

// Java: ProcessTweaksTest.testInvalidId()
func (c *ProcessTweaksTest) TestInvalidId() {
	c.TestError("goog.tweak.registerBoolean('Some ID', 'a')", ProcessTweaks_INVALID_TWEAK_ID_ERROR)
}

// Java: ProcessTweaksTest.testInvalidDefaultValue1()
func (c *ProcessTweaksTest) TestInvalidDefaultValue1() {
	c.TestSame(
		"var val = true; goog.tweak.registerBoolean('Foo', 'desc', val)",
		ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING)
}

// Java: ProcessTweaksTest.testDuplicateTweak()
func (c *ProcessTweaksTest) TestDuplicateTweak() {
	c.TestError(
		`goog.tweak.registerBoolean('TweakA', 'desc');
goog.tweak.registerBoolean('TweakA', 'desc')
`,
		ProcessTweaks_TWEAK_MULTIPLY_REGISTERED_ERROR)
}

// Java: ProcessTweaksTest.testRegisterInNonGlobalScope()
func (c *ProcessTweaksTest) TestRegisterInNonGlobalScope() {
	c.TestError(
		"function foo() {goog.tweak.registerBoolean('TweakA', 'desc');};",
		ProcessTweaks_NON_GLOBAL_TWEAK_INIT_ERROR)
}

// Java: ProcessTweaksTest.testRegisterInIf()
func (c *ProcessTweaksTest) TestRegisterInIf() {
	c.TestSame("if (true) {goog.tweak.registerBoolean('TweakA', 'desc');};")
}

// Java: ProcessTweaksTest.testWrongGetter1()
func (c *ProcessTweaksTest) TestWrongGetter1() {
	c.TestSame(
		`goog.tweak.registerBoolean('TweakA', 'desc');
goog.tweak.getString('TweakA')
`,
		ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING)
}

// Java: ProcessTweaksTest.testWrongGetter2()
func (c *ProcessTweaksTest) TestWrongGetter2() {
	c.TestSame(
		`goog.tweak.registerString('TweakA', 'desc');
goog.tweak.getNumber('TweakA')
`,
		ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING)
}

// Java: ProcessTweaksTest.testWrongGetter3()
func (c *ProcessTweaksTest) TestWrongGetter3() {
	c.TestSame(
		`goog.tweak.registerNumber('TweakA', 'desc');
goog.tweak.getBoolean('TweakA')
`,
		ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING)
}

// Java: ProcessTweaksTest.testWithNoTweaks()
func (c *ProcessTweaksTest) TestWithNoTweaks() {
	c.TestSame("var DEF=true;var x={};x.foo={}")
}

// Java: ProcessTweaksTest.testStrippingWithImplicitDefaultValues()
func (c *ProcessTweaksTest) TestStrippingWithImplicitDefaultValues() {
	c.stripTweaks = true
	c.Test(
		`goog.tweak.registerNumber('TweakA', 'desc');
goog.tweak.registerBoolean('TweakB', 'desc');
goog.tweak.registerString('TweakC', 'desc');
alert(goog.tweak.getNumber('TweakA'));
alert(goog.tweak.getBoolean('TweakB'));
alert(goog.tweak.getString('TweakC'));
`,
		"void 0; void 0; void 0; alert(0); alert(false); alert('')")
}

// Java: ProcessTweaksTest.testStrippingWithExplicitDefaultValues()
func (c *ProcessTweaksTest) TestStrippingWithExplicitDefaultValues() {
	c.stripTweaks = true
	c.Test(
		`goog.tweak.registerNumber('TweakA', 'desc', 5);
goog.tweak.registerBoolean('TweakB', 'desc', true);
goog.tweak.registerString('TweakC', 'desc', '!');
alert(goog.tweak.getNumber('TweakA'));
alert(goog.tweak.getBoolean('TweakB'));
alert(goog.tweak.getString('TweakC'));
`,
		"void 0; void 0; void 0; alert(5); alert(true); alert('!')")
}

// Java: ProcessTweaksTest.testStrippingWithInCodeOverrides()
func (c *ProcessTweaksTest) TestStrippingWithInCodeOverrides() {
	c.stripTweaks = true
	c.Test(
		`goog.tweak.registerNumber('TweakA', 'desc');
goog.tweak.registerBoolean('TweakB', 'desc');
goog.tweak.registerString('TweakC', 'desc', 'foo');
alert(goog.tweak.getNumber('TweakA'));
alert(goog.tweak.getBoolean('TweakB'));
alert(goog.tweak.getString('TweakC'));
`,
		`void 0; void 0; void 0;
alert(0); alert(false); alert('foo');
`)
}

// Java: ProcessTweaksTest.testStrippingOfManuallyRegistered1()
func (c *ProcessTweaksTest) TestStrippingOfManuallyRegistered1() {
	c.stripTweaks = true
	c.Test(
		`var reg = goog.tweak.getRegistry();
if (reg) {
  reg.register(new goog.tweak.BooleanSetting('foo', 'desc'));
  reg.getEntry('foo').setDefaultValue(1);
}
`,
		"if (null);")
}

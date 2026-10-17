// Port of com/google/javascript/jscomp/ConstCheckTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ConstCheckTest.java)

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ConstCheckTest: Tests ConstCheck.
type ConstCheckTest struct {
	CompilerTestCaseBase
}

// Java: ConstCheckTest()
func NewConstCheckTest() *ConstCheckTest {
	c := &ConstCheckTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestConstCheckTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewConstCheckTest() })
}

// Java: ConstCheckTest.setUp()
func (c *ConstCheckTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableCreateModuleMap()
}

// Java: ConstCheckTest.getProcessor(Compiler)
func (c *ConstCheckTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		NewInferConsts(compiler).Process(externs, root)
		NewConstCheck(compiler).Process(externs, root)
	})
}

// TestWarning_ is the private overload testWarning(String); the trailing _
// keeps it apart from CompilerTestCase's testWarning(...) overloads (Go-only:
// language-adapter).
//
// Java: ConstCheckTest.testWarning(String)
func (c *ConstCheckTest) TestWarning_(js string) {
	c.TestWarning(js, ConstCheck_CONST_REASSIGNED_VALUE_ERROR)
}

// Java: ConstCheckTest.testConstantDefinition1()
func (c *ConstCheckTest) TestConstantDefinition1() {
	c.TestSame("/** @const */ var XYZ = 1;")
}

// Java: ConstCheckTest.testConstantDefinition2()
func (c *ConstCheckTest) TestConstantDefinition2() {
	c.TestSame("/** @const */ var a$b$XYZ = 1;")
}

// Java: ConstCheckTest.testConstantDefinition3()
func (c *ConstCheckTest) TestConstantDefinition3() {
	c.TestNoWarning("const xyz=1;")
}

// Java: ConstCheckTest.testConstantDefinition4()
func (c *ConstCheckTest) TestConstantDefinition4() {
	fmt.Println("HELLO")
	c.TestNoWarning("const a$b$xyz = 1;")
}

// Java: ConstCheckTest.testConstantInitializedInAnonymousNamespace1()
func (c *ConstCheckTest) TestConstantInitializedInAnonymousNamespace1() {
	c.TestSame("/** @const */ var XYZ; (function(){ XYZ = 1; })();")
}

// Java: ConstCheckTest.testConstantInitializedInAnonymousNamespace2()
func (c *ConstCheckTest) TestConstantInitializedInAnonymousNamespace2() {
	c.TestSame("/** @const */ var a$b$XYZ; (function(){ a$b$XYZ = 1; })();")
}

// Java: ConstCheckTest.testObjectModified()
func (c *ConstCheckTest) TestObjectModified() {
	c.TestSame("/** @const */ var IE = true, XYZ = {a:1,b:1}; if (IE) XYZ['c'] = 1;")
}

// Java: ConstCheckTest.testObjectPropertyInitializedLate()
func (c *ConstCheckTest) TestObjectPropertyInitializedLate() {
	c.TestSame("/** @const */ var XYZ = {}; for (var i = 0; i < 10; i++) { XYZ[i] = i; }")
}

// Java: ConstCheckTest.testObjectRedefined1()
func (c *ConstCheckTest) TestObjectRedefined1() {
	c.TestError_("/** @const */ var XYZ = {}; XYZ = 2;")
}

// Java: ConstCheckTest.testObjectRedefined2()
func (c *ConstCheckTest) TestObjectRedefined2() {
	c.TestError_("const xyz = {}; xyz = 2;")
}

// Java: ConstCheckTest.testConstantRedefined1()
func (c *ConstCheckTest) TestConstantRedefined1() {
	c.TestError_("/** @const */ var XYZ = 1; XYZ = 2;")
}

// Java: ConstCheckTest.testConstantRedefined2()
func (c *ConstCheckTest) TestConstantRedefined2() {
	c.TestError_("/** @const */ var a$b$XYZ = 1; a$b$XYZ = 2;")
}

// test will be caught be earlier pass, but demonstrates that it returns error upon const
// reassigning
// Java: ConstCheckTest.testConstantRedefined3()
func (c *ConstCheckTest) TestConstantRedefined3() {
	c.TestWarning_("const xyz = 1; xyz = 2;")
}

// Java: ConstCheckTest.testConstantRedefined4()
func (c *ConstCheckTest) TestConstantRedefined4() {
	c.TestError_("/** @const */ let XYZ = 1; XYZ = 2;")
}

// Java: ConstCheckTest.testConstantRedefinedInLocalScope1()
func (c *ConstCheckTest) TestConstantRedefinedInLocalScope1() {
	c.TestError_("/** @const */ var XYZ = 1; (function(){ XYZ = 2; })();")
}

// Java: ConstCheckTest.testConstantRedefinedInLocalScope2()
func (c *ConstCheckTest) TestConstantRedefinedInLocalScope2() {
	c.TestError_("/** @const */ var a$b$XYZ = 1; (function(){ a$b$XYZ = 2; })();")
}

// Java: ConstCheckTest.testConstantRedefinedInLocalScopeOutOfOrder()
func (c *ConstCheckTest) TestConstantRedefinedInLocalScopeOutOfOrder() {
	c.TestError_("function f() { XYZ = 2; } /** @const */ var XYZ = 1;")
}

// Java: ConstCheckTest.testConstantPostIncremented1()
func (c *ConstCheckTest) TestConstantPostIncremented1() {
	c.TestError_("/** @const */ var XYZ = 1; XYZ++;")
}

// Java: ConstCheckTest.testConstantPostIncremented2()
func (c *ConstCheckTest) TestConstantPostIncremented2() {
	c.TestError_("/** @const */ var a$b$XYZ = 1; a$b$XYZ++;")
}

// Java: ConstCheckTest.testConstantPreIncremented1()
func (c *ConstCheckTest) TestConstantPreIncremented1() {
	c.TestError_("/** @const */ var XYZ = 1; ++XYZ;")
}

// Java: ConstCheckTest.testConstantPreIncremented2()
func (c *ConstCheckTest) TestConstantPreIncremented2() {
	c.TestError_("/** @const */ var a$b$XYZ = 1; ++a$b$XYZ;")
}

// Java: ConstCheckTest.testConstantPostDecremented1()
func (c *ConstCheckTest) TestConstantPostDecremented1() {
	c.TestError_("/** @const */ var XYZ = 1; XYZ--;")
}

// Java: ConstCheckTest.testConstantPostDecremented2()
func (c *ConstCheckTest) TestConstantPostDecremented2() {
	c.TestError_("/** @const */ var a$b$XYZ = 1; a$b$XYZ--;")
}

// Java: ConstCheckTest.testConstantPreDecremented1()
func (c *ConstCheckTest) TestConstantPreDecremented1() {
	c.TestError_("/** @const */ var XYZ = 1; --XYZ;")
}

// Java: ConstCheckTest.testConstantPreDecremented2()
func (c *ConstCheckTest) TestConstantPreDecremented2() {
	c.TestError_("/** @const */ var a$b$XYZ = 1; --a$b$XYZ;")
}

// Java: ConstCheckTest.testConstantPreDecremented3()
func (c *ConstCheckTest) TestConstantPreDecremented3() {
	c.TestWarning_("const xyz = 1; --xyz;")
}

// Java: ConstCheckTest.testAbbreviatedArithmeticAssignment1()
func (c *ConstCheckTest) TestAbbreviatedArithmeticAssignment1() {
	c.TestError_("/** @const */ var XYZ = 1; XYZ += 2;")
}

// Java: ConstCheckTest.testAbbreviatedArithmeticAssignment2()
func (c *ConstCheckTest) TestAbbreviatedArithmeticAssignment2() {
	c.TestError_("/** @const */ var a$b$XYZ = 1; a$b$XYZ %= 2;")
}

// Java: ConstCheckTest.testAbbreviatedArithmeticAssignment3()
func (c *ConstCheckTest) TestAbbreviatedArithmeticAssignment3() {
	c.TestError_("/** @const */ var a$b$XYZ = 1; a$b$XYZ **= 2;")
}

// Java: ConstCheckTest.testAbbreviatedBitAssignment1()
func (c *ConstCheckTest) TestAbbreviatedBitAssignment1() {
	c.TestError_("/** @const */ var XYZ = 1; XYZ |= 2;")
}

// Java: ConstCheckTest.testAbbreviatedBitAssignment2()
func (c *ConstCheckTest) TestAbbreviatedBitAssignment2() {
	c.TestError_("/** @const */ var a$b$XYZ = 1; a$b$XYZ &= 2;")
}

// Java: ConstCheckTest.testAbbreviatedShiftAssignment1()
func (c *ConstCheckTest) TestAbbreviatedShiftAssignment1() {
	c.TestError_("/** @const */ var XYZ = 1; XYZ >>= 2;")
}

// Java: ConstCheckTest.testAbbreviatedShiftAssignment2()
func (c *ConstCheckTest) TestAbbreviatedShiftAssignment2() {
	c.TestError_("/** @const */ var a$b$XYZ = 1; a$b$XYZ <<= 2;")
}

// Java: ConstCheckTest.testConstAnnotation1()
func (c *ConstCheckTest) TestConstAnnotation1() {
	c.TestWarning_("/** @const */ var XYZ = 1; XYZ = 2;")
}

// Java: ConstCheckTest.testConstAnnotation2()
func (c *ConstCheckTest) TestConstAnnotation2() {
	c.TestWarning_("/** @const */ let x = 1; x = 2;")
}

// Java: ConstCheckTest.testConstAnnotation3()
func (c *ConstCheckTest) TestConstAnnotation3() {
	c.TestWarning_("/** @const */ const xyz = 1; xyz = 2;")
}

// Java: ConstCheckTest.testConstSuppressionInFileJsDoc()
func (c *ConstCheckTest) TestConstSuppressionInFileJsDoc() {
	c.TestSame(
		`/**
 * @fileoverview
 * @suppress {const}
 */
/** @const */ var xyz = 1; xyz = 3;
`)
}

// Java: ConstCheckTest.testConstSuppressionOnAssignment()
func (c *ConstCheckTest) TestConstSuppressionOnAssignment() {
	c.TestSame("/** @const */ var xyz = 1; /** @suppress {const} */ xyz = 3;")
}

// Java: ConstCheckTest.testConstSuppressionOnAddAssign()
func (c *ConstCheckTest) TestConstSuppressionOnAddAssign() {
	c.TestSame("/** @const */ var xyz = 1; /** @suppress {const} */ xyz += 1;")
}

// Java: ConstCheckTest.testConstSuppressionOnVar()
func (c *ConstCheckTest) TestConstSuppressionOnVar() {
	c.TestSame("/** @const */ var xyz = 1;\n/** @suppress {const} */ var xyz = 3;")
}

// If there are two 'var' statements for the same variable, one in externs and
// one in the JS, there is no normalization, and the suppression remains on the
// statement in the JS.
// Java: ConstCheckTest.testConstSuppressionOnVarFromExterns()
func (c *ConstCheckTest) TestConstSuppressionOnVarFromExterns() {
	externs := "/** @const */ var xyz;"
	js := "/** @suppress {const} */ var xyz = 3;"
	c.TestSame(c.Externs(externs), c.Srcs(js))
}

// Java: ConstCheckTest.testConstSuppressionOnInc()
func (c *ConstCheckTest) TestConstSuppressionOnInc() {
	c.TestSame("/** @const */ var xyz = 1; /** @suppress {const} */ xyz++;")
}

// Java: ConstCheckTest.testConstNameInExterns()
func (c *ConstCheckTest) TestConstNameInExterns() {
	externs := "/** @const */ var FOO;"
	js := "FOO = 1;"
	c.TestWarning(c.Externs(externs), c.Srcs(js), c.Warning(ConstCheck_CONST_REASSIGNED_VALUE_ERROR))
}

// Java: ConstCheckTest.testGoogModuleExportsShadowingGlobal()
func (c *ConstCheckTest) TestGoogModuleExportsShadowingGlobal() {
	c.TestSame(c.Srcs("/** @const */ var exports = {};", "goog.module('m'); exports = class {};"))
}

// Java: ConstCheckTest.testGoogModuleExportsReassigned()
func (c *ConstCheckTest) TestGoogModuleExportsReassigned() {
	c.TestWarning_("goog.module('m'); exports = class {}; exports = class {};")
}

// Java: ConstCheckTest.testGoogProvide_rootNamespaceExplicitlyDeclared()
func (c *ConstCheckTest) TestGoogProvide_rootNamespaceExplicitlyDeclared() {
	c.TestSame("goog.provide('a'); var a = {};")
	c.TestSame("goog.provide('a'); var a = {}; a = 1;")
	c.TestWarning(
		"goog.provide('a'); /** @const */ var a = {}; a = 1;",
		ConstCheck_CONST_REASSIGNED_VALUE_ERROR)
}

// Java: ConstCheckTest.testGoogProvide_rootNamespaceImplicitlyDeclared()
func (c *ConstCheckTest) TestGoogProvide_rootNamespaceImplicitlyDeclared() {
	c.TestSame("goog.provide('a'); a = {};")
	c.TestSame("goog.provide('a'); a = {}; a = 1;")
	c.TestSame("goog.provide('a'); a.B = class {};")
	c.TestSame("goog.provide('a.b.c'); a = class {};")
	c.TestSame(c.Srcs("goog.provide('a.b');", "goog.provide('a.c');"))
	c.TestSame("goog.provide('a'); a = 0; a++;")
}

// TestError_ is the private overload testError(String); see TestWarning_.
//
// Java: ConstCheckTest.testError(String)
func (c *ConstCheckTest) TestError_(js string) {
	c.TestWarning(js, ConstCheck_CONST_REASSIGNED_VALUE_ERROR)
}

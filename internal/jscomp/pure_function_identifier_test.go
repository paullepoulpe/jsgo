// Port of com/google/javascript/jscomp/PureFunctionIdentifierTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/PureFunctionIdentifierTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

type PureFunctionIdentifierTest struct {
	CompilerTestCaseBase
	noSideEffectCalls                []*rhino.Node
	regExpHaveSideEffects            bool
	enableArtificialPurityDebugError bool
}

// Java: PureFunctionIdentifierTest()
func NewPureFunctionIdentifierTest() *PureFunctionIdentifierTest {
	c := &PureFunctionIdentifierTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, PureFunctionIdentifierTest_TEST_EXTERNS)
	c.regExpHaveSideEffects = true
	c.enableArtificialPurityDebugError = false
	return c
}

func TestPureFunctionIdentifierTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPureFunctionIdentifierTest() })
}

var PureFunctionIdentifierTest_TEST_EXTERNS = CompilerTypeTestCase_DEFAULT_EXTERNS +
	"var window; window.setTimeout;\n" +
	"/**@nosideeffects*/ function externSENone(){}\n" +
	"/**@modifies{this}*/ function externSEThis(){}\n" +
	"/**@constructor\n" +
	" * @modifies{this}*/\n" +
	"function externObjSEThis(){}\n" +
	"/**\n" +
	" * @param {string} s id.\n" +
	" * @return {string}\n" +
	" * @modifies{this}\n" +
	" */\n" +
	"externObjSEThis.prototype.externObjSEThisMethod = function(s) {};\n" +
	"/**\n" +
	" * @param {string} s id.\n" +
	" * @return {string}\n" +
	" * @modifies{arguments}\n" +
	" */\n" +
	"externObjSEThis.prototype.externObjSEThisMethod2 = function(s) {};\n" +
	"/**@nosideeffects*/function Error(){}\n" +
	"function externSef1(){}\n" +
	"/**@nosideeffects*/function externNsef1(){}\n" +
	"var externSef2 = function(){};\n" +
	"/**@nosideeffects*/var externNsef2 = function(){};\n" +
	"var externNsef3 = /**@nosideeffects*/function(){};\n" +
	"var externObj;\n" +
	"externObj.sef1 = function(){};\n" +
	"/**@nosideeffects*/externObj.nsef1 = function(){};\n" +
	"externObj.nsef2 = /**@nosideeffects*/function(){};\n" +
	"externObj.partialFn = function(){};\n" +
	"externObj.partialSharedFn;\n" +
	"var externObj2;\n" +
	"externObj2.partialSharedFn = /**@nosideeffects*/function(){};\n" +
	"/**@constructor*/function externSefConstructor(){}\n" +
	"externSefConstructor.prototype.sefFnOfSefObj = function(){};\n" +
	"externSefConstructor.prototype.nsefFnOfSefObj =\n" +
	"  /**@nosideeffects*/function(){};\n" +
	"externSefConstructor.prototype.externShared = function(){};\n" +
	"/**@constructor@nosideeffects*/function externNsefConstructor(){}\n" +
	"externNsefConstructor.prototype.sefFnOfNsefObj = function(){};\n" +
	"externNsefConstructor.prototype.nsefFnOfNsefObj =\n" +
	"  /**@nosideeffects*/function(){};\n" +
	"externNsefConstructor.prototype.externShared =\n" +
	"  /**@nosideeffects*/function(){};\n" +
	"/**@constructor @nosideeffects*/function externNsefConstructor2(){}\n" +
	"externNsefConstructor2.prototype.externShared =\n" +
	"  /**@nosideeffects*/function(){};\n" +
	"externNsefConstructor.prototype.sharedPartialSef;\n" +
	"/**@nosideeffects*/externNsefConstructor.prototype.sharedPartialNsef;\n" +
	"// An externs definition with a stub before.\n" +
	"/**@constructor*/function externObj3(){}\n" +
	"externObj3.prototype.propWithStubBefore;\n" +
	"/**\n" +
	" * @param {string} s id.\n" +
	" * @return {string}\n" +
	" * @nosideeffects\n" +
	" */\n" +
	"externObj3.prototype.propWithStubBefore = function(s) {};\n" +
	"// useless JsDoc\n" +
	"/**\n" +
	" * @see {foo}\n" +
	" */\n" +
	"externObj3.prototype.propWithStubBeforeWithJSDoc;\n" +
	"/**\n" +
	" * @param {string} s id.\n" +
	" * @return {string}\n" +
	" * @nosideeffects\n" +
	" */\n" +
	"externObj3.prototype.propWithStubBeforeWithJSDoc = function(s) {};\n" +
	"// An externs definition with a stub after.\n" +
	"/**@constructor*/function externObj4(){}\n" +
	"/**\n" +
	" * @param {string} s id.\n" +
	" * @return {string}\n" +
	" * @nosideeffects\n" +
	" */\n" +
	"externObj4.prototype.propWithStubAfter = function(s) {};\n" +
	"externObj4.prototype.propWithStubAfter;\n" +
	"/**\n" +
	" * @param {string} s id.\n" +
	" * @return {string}\n" +
	" * @nosideeffects\n" +
	" */\n" +
	"externObj4.prototype.propWithStubAfterWithJSDoc = function(s) {};\n" +
	"// useless JsDoc\n" +
	"/**\n" +
	" * @see {foo}\n" +
	" */\n" +
	"externObj4.prototype.propWithStubAfterWithJSDoc;\n" +
	"goog.reflect = {};\n" +
	"goog.reflect.cache = function(a, b, c, opt_d) {};\n" +
	"/** @nosideeffects */\n" +
	"externObj.prototype.duplicateExternFunc = function() {};\n" +
	"externObj2.prototype.duplicateExternFunc = function() {};\n" +
	"externObj.prototype['weirdDefinition'] = function() {};\n"

// Java: PureFunctionIdentifierTest.setUp()
func (c *PureFunctionIdentifierTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()

	// Allow testing of features that are not yet fully supported.
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.DisableCompareJsDoc()
	c.EnableGatherExternProperties()
	c.EnableTypeCheck()
}

// Java: PureFunctionIdentifierTest.tearDown()
func (c *PureFunctionIdentifierTest) TearDown() {
	c.CompilerTestCaseBase.TearDown()
	c.regExpHaveSideEffects = true
	c.enableArtificialPurityDebugError = false
}

// Java: PureFunctionIdentifierTest.getOptions()
func (c *PureFunctionIdentifierTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	if c.enableArtificialPurityDebugError {
		options.SetWarningLevel(
			DiagnosticGroup_ForType(PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION),
			CheckLevel_ERROR)
	}
	return options
}

// PureFunctionIdentifierTest_NoSideEffectCallEnumerator runs
// PureFunctionIdentifier, then gathers a list of calls that are marked as
// having no side effects.
type PureFunctionIdentifierTest_NoSideEffectCallEnumerator struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer    *PureFunctionIdentifierTest
	compiler *Compiler
}

// Java: PureFunctionIdentifierTest.NoSideEffectCallEnumerator(Compiler)
func NewPureFunctionIdentifierTest_NoSideEffectCallEnumerator(outer *PureFunctionIdentifierTest, compiler *Compiler) *PureFunctionIdentifierTest_NoSideEffectCallEnumerator {
	e := &PureFunctionIdentifierTest_NoSideEffectCallEnumerator{outer: outer}
	e.compiler = compiler
	return e
}

// Java: PureFunctionIdentifierTest.NoSideEffectCallEnumerator.process(Node, Node)
func (e *PureFunctionIdentifierTest_NoSideEffectCallEnumerator) Process(externs *rhino.Node, root *rhino.Node) {
	e.outer.noSideEffectCalls = []*rhino.Node{}
	// TODO(nickreid): Move these into 'getOptions' and 'getCompiler' overrides.
	e.compiler.SetHasRegExpGlobalReferences(e.outer.regExpHaveSideEffects)
	e.compiler.GetOptions().SetUseTypesForLocalOptimization(true)

	NewPureFunctionIdentifier_Driver(e.compiler).Process(externs, root)
	NodeTraversal_Traverse(e.compiler, externs, e)
	NodeTraversal_Traverse(e.compiler, root, e)
}

// Java: PureFunctionIdentifierTest.NoSideEffectCallEnumerator.visit(NodeTraversal, Node, Node)
func (e *PureFunctionIdentifierTest_NoSideEffectCallEnumerator) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsNew() {
		if !e.compiler.GetAstAnalyzer().ConstructorCallHasSideEffects(n) {
			e.outer.noSideEffectCalls = append(e.outer.noSideEffectCalls, n.GetFirstChild())
		}
	} else if NodeUtil_IsInvocation(n) {
		if !e.compiler.GetAstAnalyzer().FunctionCallHasSideEffects(n) {
			e.outer.noSideEffectCalls = append(e.outer.noSideEffectCalls, n.GetFirstChild())
		}
	}
}

// Java: PureFunctionIdentifierTest.testOptionalChainGetProp()
func (c *PureFunctionIdentifierTest) TestOptionalChainGetProp() {
	c.AssertNoPureCalls("externObj?.sef1()")
	c.AssertPureCallsMarked("externObj?.nsef1()", []string{"externObj?.nsef1"})
}

// Java: PureFunctionIdentifierTest.testOptionalChainGetElem()
func (c *PureFunctionIdentifierTest) TestOptionalChainGetElem() {
	c.AssertNoPureCalls("externObj?.['sef1']()")
	// The use of `[]` hides the exact method being called from the compiler,
	// so it will assume there are side-effects.
	c.AssertNoPureCalls("externObj?.['nsef1']()")
}

// Java: PureFunctionIdentifierTest.testOptionalChainCall()
func (c *PureFunctionIdentifierTest) TestOptionalChainCall() {
	c.AssertNoPureCalls("externObj.sef1?.()")
	c.AssertNoPureCalls("externObj?.sef1?.()")
	c.AssertPureCallsMarked("externObj.nsef1?.()", []string{"externObj.nsef1"})
	c.AssertPureCallsMarked("externObj?.nsef1?.()", []string{"externObj?.nsef1"})
}

// Java: PureFunctionIdentifierTest.testOptionalGetterAccess_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestOptionalGetterAccess_hasSideEffects() {
	c.AssertNoPureCalls(
		"class Foo { get getter() { } }\n" +
			"\n" +
			"function foo(a) { a?.getter; }\n" +
			"foo(x);\n")
}

// Java: PureFunctionIdentifierTest.testOptionalNormalPropAccess_hasNoSideEffects()
func (c *PureFunctionIdentifierTest) TestOptionalNormalPropAccess_hasNoSideEffects() {
	c.AssertPureCallsMarked(
		"function foo(a) { a?.prop; }\n"+
			"foo(x);\n",
		[]string{"foo"})
}

// Java: PureFunctionIdentifierTest.testOptionalElementAccess_hasNoSideEffects()
func (c *PureFunctionIdentifierTest) TestOptionalElementAccess_hasNoSideEffects() {
	c.AssertPureCallsMarked(
		"function foo(a) { a?.['prop']; }\n"+
			"foo(x);\n",
		[]string{"foo"})
}

// Java: PureFunctionIdentifierTest.testIssue303()
func (c *PureFunctionIdentifierTest) TestIssue303() {
	source :=
		"/** @constructor */ function F() {\n" +
			"  var self = this;\n" +
			"  window.setTimeout(function() {\n" +
			"    window.location = self.location;\n" +
			"  }, 0);\n" +
			"}\n" +
			"F.prototype.setLocation = function(x) {\n" +
			"  this.location = x;\n" +
			"};\n" +
			"(new F()).setLocation('http://www.google.com/');\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testIssue303b()
func (c *PureFunctionIdentifierTest) TestIssue303b() {
	source :=
		"/** @constructor */ function F() {\n" +
			"  var self = this;\n" +
			"  window.setTimeout(function() {\n" +
			"    window.location = self.location;\n" +
			"  }, 0);\n" +
			"}\n" +
			"F.prototype.setLocation = function(x) {\n" +
			"  this.location = x;\n" +
			"};\n" +
			"function x() {\n" +
			"  (new F()).setLocation('http://www.google.com/');\n" +
			"} window['x'] = x;\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns_new1()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns_new1() {
	c.AssertPureCallsMarked("externSENone()", []string{"externSENone"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns_new2()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns_new2() {
	c.AssertNoPureCalls("externSEThis()")
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns_new3()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns_new3() {
	c.AssertPureCallsMarked("new externObjSEThis()", []string{"externObjSEThis"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns_new4()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns_new4() {
	// The entire expression containing "externObjSEThisMethod" is considered
	// side-effect free in this context.

	c.AssertPureCallsMarked(
		"new externObjSEThis().externObjSEThisMethod('')",
		[]string{"externObjSEThis", "new externObjSEThis().externObjSEThisMethod"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns_new5()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns_new5() {
	c.AssertPureCallsMarked(
		"function f() { new externObjSEThis() }; f();", []string{"externObjSEThis", "f"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns_new6()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns_new6() {
	// While "externObjSEThisMethod" has modifies "this"
	// it does not have global side-effects with "this" is
	// a known local value.
	// TODO(johnlenz): "f" is side-effect free but we need
	// to propagate that "externObjSEThisMethod" is modifying
	// a local object.
	source :=
		"function f() {\n" +
			"  new externObjSEThis().externObjSEThisMethod('')\n" +
			"};\n" +
			"f();\n"
	c.AssertPureCallsMarked(
		source, []string{"externObjSEThis", "new externObjSEThis().externObjSEThisMethod"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns_new7()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns_new7() {
	// While "externObjSEThisMethod" has modifies "this"
	// it does not have global side-effects with "this" is
	// a known local value.
	source :=
		"function f() {\n" +
			"  var x = new externObjSEThis();\n" +
			"  x.externObjSEThisMethod('')\n" +
			"};\n" +
			"f();\n"
	c.AssertPureCallsMarked(source, []string{"externObjSEThis"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns_new8()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns_new8() {
	// "externObjSEThisMethod" modifies "this", the 'this'
	// is not a known local value, so it must be assumed it is to
	// have global side-effects.
	source :=
		"function f(x) {\n" +
			"  x.externObjSEThisMethod('')\n" +
			"};\n" +
			"f(new externObjSEThis());\n"
	c.AssertPureCallsMarked(source, []string{"externObjSEThis"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns_new9()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns_new9() {
	// "externObjSEThisMethod" modifies "this", the 'this'
	// is not a known local value, so it must be assumed it is to
	// have global side-effects.  All possible values of "x" are considered
	// as no intraprocedural data flow is done.
	source :=
		"function f(x) {\n" +
			"  x = new externObjSEThis();\n" +
			"  x.externObjSEThisMethod('')\n" +
			"};\n" +
			"f(g);\n"
	c.AssertPureCallsMarked(source, []string{"externObjSEThis"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns_new10()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns_new10() {
	source :=
		"function f() {\n" +
			"  new externObjSEThis().externObjSEThisMethod2('')\n" +
			"};\n" +
			"f();\n"
	c.AssertPureCallsMarked(
		source,
		[]string{"externObjSEThis", "new externObjSEThis().externObjSEThisMethod2", "f"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns1()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns1() {
	c.AssertNoPureCalls("externSef1()")
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns2()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns2() {
	c.AssertNoPureCalls("externSef2()")
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns3()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns3() {
	c.AssertPureCallsMarked("externNsef1()", []string{"externNsef1"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns4()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns4() {
	c.AssertPureCallsMarked("externNsef2()", []string{"externNsef2"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExterns5()
func (c *PureFunctionIdentifierTest) TestAnnotationInExterns5() {
	c.AssertPureCallsMarked("externNsef3()", []string{"externNsef3"})
}

// Java: PureFunctionIdentifierTest.testNamespaceAnnotationInExterns1()
func (c *PureFunctionIdentifierTest) TestNamespaceAnnotationInExterns1() {
	c.AssertNoPureCalls("externObj.sef1()")
}

// Java: PureFunctionIdentifierTest.testNamespaceAnnotationInExterns2()
func (c *PureFunctionIdentifierTest) TestNamespaceAnnotationInExterns2() {
	c.AssertPureCallsMarked("externObj.nsef1()", []string{"externObj.nsef1"})
}

// Java: PureFunctionIdentifierTest.testNamespaceAnnotationInExterns3()
func (c *PureFunctionIdentifierTest) TestNamespaceAnnotationInExterns3() {
	c.AssertPureCallsMarked("externObj.nsef2()", []string{"externObj.nsef2"})
}

// Java: PureFunctionIdentifierTest.testNamespaceAnnotationInExterns4()
func (c *PureFunctionIdentifierTest) TestNamespaceAnnotationInExterns4() {
	c.AssertNoPureCalls("externObj.partialFn()")
}

// Java: PureFunctionIdentifierTest.testNamespaceAnnotationInExterns5()
func (c *PureFunctionIdentifierTest) TestNamespaceAnnotationInExterns5() {
	// Test that adding a second definition for a partially defined
	// function doesn't make us think that the function has no side
	// effects.
	templateSrc := "var o = {}; o.<fnName> = function(){}; o.<fnName>()"

	// Ensure that functions with name != "partialFn" get marked.
	c.AssertPureCallsMarked(
		jutil.StringReplace(templateSrc, "<fnName>", "notPartialFn"), []string{"o.notPartialFn"})

	c.AssertNoPureCalls(jutil.StringReplace(templateSrc, "<fnName>", "partialFn"))
}

// Java: PureFunctionIdentifierTest.testNamespaceAnnotationInExterns6()
func (c *PureFunctionIdentifierTest) TestNamespaceAnnotationInExterns6() {
	c.AssertNoPureCalls("externObj.partialSharedFn()")
}

// Java: PureFunctionIdentifierTest.testConstructorAnnotationInExterns1()
func (c *PureFunctionIdentifierTest) TestConstructorAnnotationInExterns1() {
	c.AssertNoPureCalls("new externSefConstructor()")
}

// Java: PureFunctionIdentifierTest.testConstructorAnnotationInExterns2()
func (c *PureFunctionIdentifierTest) TestConstructorAnnotationInExterns2() {
	source :=
		"var a = new externSefConstructor();\n" +
			"a.sefFnOfSefObj()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testConstructorAnnotationInExterns3()
func (c *PureFunctionIdentifierTest) TestConstructorAnnotationInExterns3() {
	source :=
		"var a = new externSefConstructor();\n" +
			"a.nsefFnOfSefObj()\n"
	c.AssertPureCallsMarked(source, []string{"a.nsefFnOfSefObj"})
}

// Java: PureFunctionIdentifierTest.testConstructorAnnotationInExterns4()
func (c *PureFunctionIdentifierTest) TestConstructorAnnotationInExterns4() {
	source :=
		"var a = new externSefConstructor();\n" +
			"a.externShared()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testConstructorAnnotationInExterns5()
func (c *PureFunctionIdentifierTest) TestConstructorAnnotationInExterns5() {
	c.AssertPureCallsMarked("new externNsefConstructor()", []string{"externNsefConstructor"})
}

// Java: PureFunctionIdentifierTest.testConstructorAnnotationInExterns6()
func (c *PureFunctionIdentifierTest) TestConstructorAnnotationInExterns6() {
	source :=
		"var a = new externNsefConstructor();\n" +
			"a.sefFnOfNsefObj()\n"
	c.AssertPureCallsMarked(source, []string{"externNsefConstructor"})
}

// Java: PureFunctionIdentifierTest.testConstructorAnnotationInExterns7()
func (c *PureFunctionIdentifierTest) TestConstructorAnnotationInExterns7() {
	source :=
		"var a = new externNsefConstructor();\n" +
			"a.nsefFnOfNsefObj()\n"
	c.AssertPureCallsMarked(source, []string{"externNsefConstructor", "a.nsefFnOfNsefObj"})
}

// Java: PureFunctionIdentifierTest.testConstructorAnnotationInExterns8()
func (c *PureFunctionIdentifierTest) TestConstructorAnnotationInExterns8() {
	source :=
		"var a = new externNsefConstructor();\n" +
			"a.externShared()\n"
	c.AssertPureCallsMarked(source, []string{"externNsefConstructor"})
}

// Java: PureFunctionIdentifierTest.testSharedFunctionName1()
func (c *PureFunctionIdentifierTest) TestSharedFunctionName1() {
	source :=
		"if (true) {\n" +
			"  a = new externNsefConstructor()\n" +
			"} else {\n" +
			"  a = new externSefConstructor()\n" +
			"}\n" +
			"a.externShared()\n"
	c.AssertPureCallsMarked(source, []string{"externNsefConstructor"})
}

// Java: PureFunctionIdentifierTest.testSharedFunctionName2()
func (c *PureFunctionIdentifierTest) TestSharedFunctionName2() {
	// Implementation for both externNsefConstructor and externNsefConstructor2
	// have no side effects.
	broken := true
	if broken {
		c.AssertPureCallsMarked(
			"var a;\n"+
				"if (true) {\n"+
				"  a = new externNsefConstructor()\n"+
				"} else {\n"+
				"  a = new externNsefConstructor2()\n"+
				"}\n"+
				"a.externShared()\n",
			[]string{"externNsefConstructor", "externNsefConstructor2"})
	} else {
		c.AssertPureCallsMarked(
			"var a;\n"+
				"if (true) {\n"+
				"  a = new externNsefConstructor()\n"+
				"} else {\n"+
				"  a = new externNsefConstructor2()\n"+
				"}\n"+
				"a.externShared()\n",
			[]string{"externNsefConstructor", "externNsefConstructor2", "a.externShared"})
	}
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsJsDocOnVariable()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsJsDocOnVariable() {
	c.AssertPureCallsMarked(
		"/** @nosideeffects */\n"+
			"function b() { throw 0; }\n"+
			"b();\n",
		[]string{"b"})
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsJsDocOnMethod()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsJsDocOnMethod() {
	c.AssertPureCallsMarked(
		"class Foo {\n"+
			"  /** @nosideeffects */ b() { throw 0; }\n"+
			"}\n"+
			"x.b();\n",
		[]string{"x.b"})
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsJsDoc_notPropagatedIntoNestedFunctions()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsJsDoc_notPropagatedIntoNestedFunctions() {
	c.AssertPureCallsMarked(
		"const obj = {};\n"+
			"/** @nosideeffects */\n"+
			"function pureOuter() {\n"+
			"  obj.impureInner = function() {\n"+
			"    throw 1;\n"+
			"  }\n"+
			"  /** @nosideeffects */\n"+
			"  obj.pureInner = function() {\n"+
			"    throw 1;\n"+
			"  }\n"+
			"  throw 0;\n"+
			"}\n"+
			"pureOuter();\n"+
			"obj.impureInner();\n"+
			"obj.pureInner();\n",
		[]string{"pureOuter", "obj.pureInner"})
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsJsDoc_astWithinFunctionBodyStillMarkedWithSideEffects()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsJsDoc_astWithinFunctionBodyStillMarkedWithSideEffects() {
	c.AssertPureCallsMarked(
		"/** @nosideeffects */\n"+
			"function pure() {\n"+
			"  // this is marked as a side-effectful call, but because of the earlier `@nosideeffects`\n"+
			"  // annotation, the call to `pure` is still marked as side-effect free.\n"+
			"  externSef1();\n"+
			"  externNsef1();\n"+
			"  throw 0;\n"+
			"}\n"+
			"pure();\n",
		[]string{"pure", "externNsef1"})
}

// Java: PureFunctionIdentifierTest.testAmbiguousExternAndSourceWithNoSideEffectsJsDoc()
func (c *PureFunctionIdentifierTest) TestAmbiguousExternAndSourceWithNoSideEffectsJsDoc() {
	c.Test(
		c.Externs("const a = {}; a.m = function() {};"),
		c.Srcs("class C { /** @nosideeffects */ m() {}} foo.m();"))

	AssertThat(c.noSideEffectCalls).IsEmpty()
}

// Java: PureFunctionIdentifierTest.testMultipleMethodDefinitionsOnlyOneWithNoSideEffectsJsDoc_otherIsPure()
func (c *PureFunctionIdentifierTest) TestMultipleMethodDefinitionsOnlyOneWithNoSideEffectsJsDoc_otherIsPure() {
	c.AssertPureCallsMarked(
		"class Foo {\n"+
			"  /** @nosideeffects */ b() { throw 0; }\n"+
			"}\n"+
			"class Bar {\n"+
			"  b() {}\n"+
			"}\n"+
			"x.b();\n",
		[]string{"x.b"})
}

// Java: PureFunctionIdentifierTest.testMultipleMethodDefinitionsOnlyOneWithNoSideEffectsJsDoc_otherHasSideEffects()
func (c *PureFunctionIdentifierTest) TestMultipleMethodDefinitionsOnlyOneWithNoSideEffectsJsDoc_otherHasSideEffects() {
	c.AssertNoPureCalls(
		"class Foo {\n" +
			"  /** @nosideeffects */ b() { throw 0; }\n" +
			"}\n" +
			"class Bar {\n" +
			"  b() { throw 0; }\n" +
			"}\n" +
			"x.b();\n")
}

// Java: PureFunctionIdentifierTest.testMultipleFunctionDefinitionsOnlyOneWithNoSideEffectsJsDoc()
func (c *PureFunctionIdentifierTest) TestMultipleFunctionDefinitionsOnlyOneWithNoSideEffectsJsDoc() {
	source :=
		"let b = /** @nosideeffects */(function(){ throw 0; });\n" +
			"if (cond) {\n" +
			"  b = () => {}\n" +
			"}\n" +
			"b()\n"
	c.Test(c.Externs(""), c.Srcs(source))
	AssertThat(c.noSideEffectCalls).HasSize(1)
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsJsDoc_errorIfCannotBeApplied_ambiguousDefModifyingArguments()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsJsDoc_errorIfCannotBeApplied_ambiguousDefModifyingArguments() {
	c.enableArtificialPurityDebugError = true
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  /** @nosideeffects */ b() { throw 0; } // line 2\n"+
				"}\n"+
				"class Bar {\n"+
				"  b(obj) { obj.prop = 34; } // line 5\n"+
				"}\n"+
				"x.b();\n"),
		c.Error(PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION).
			WithMessageContaining(
				"Side-effectful definitions:\n"+
					"  testcode:5:2"))
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsJsDoc_errorIfCannotBeApplied_ambiguousDefModifyingThis()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsJsDoc_errorIfCannotBeApplied_ambiguousDefModifyingThis() {
	c.enableArtificialPurityDebugError = true
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  /** @nosideeffects */ b() { throw 0; } // line 2\n"+
				"}\n"+
				"class Bar {\n"+
				"  b() { this.prop = 34; } // line 5\n"+
				"}\n"+
				"x.b();\n"),
		c.Error(PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION).
			WithMessageContaining(
				"Side-effectful definitions:\n"+
					"  testcode:5:2"))
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsJsDoc_errorIfCannotBeApplied_ambiguousDefWithThrow()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsJsDoc_errorIfCannotBeApplied_ambiguousDefWithThrow() {
	c.enableArtificialPurityDebugError = true
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  /** @nosideeffects */ b() { throw 0; } // line 2\n"+
				"}\n"+
				"class Bar {\n"+
				"  b() { throw 0; } // line 5\n"+
				"}\n"+
				"x.b();\n"),
		c.Error(PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION).
			WithMessageContaining(
				"Side-effectful definitions:\n"+
					"  testcode:5:2"))
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsJsDoc_errorIfCannotBeApplied_ambiguousExternFunction()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsJsDoc_errorIfCannotBeApplied_ambiguousExternFunction() {
	c.enableArtificialPurityDebugError = true
	c.Test(
		c.Externs(
			"class Bar {\n"+
				"  b() {}\n"+
				"}\n"),
		c.Srcs(
			"class Foo {\n"+
				"  /** @nosideeffects */ b() { throw 0; } // line 2\n"+
				"}\n"+
				"x.b();\n"),
		c.Error(PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION).
			WithMessageContaining(
				"Side-effectful definitions:\n"+
					"  externs:2:2"))
}

// Java: PureFunctionIdentifierTest.testMultipleFunctionDefinitions_onlySomeWithNoSideEffectsJsDoc_errorForDebugging()
func (c *PureFunctionIdentifierTest) TestMultipleFunctionDefinitions_onlySomeWithNoSideEffectsJsDoc_errorForDebugging() {
	c.enableArtificialPurityDebugError = true
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  /** @nosideeffects */ b() { throw 0; } // line 2\n"+
				"}\n"+
				"class Bar {\n"+
				"  b() { throw 0; } // line 5\n"+
				"}\n"+
				"class Baz {\n"+
				"  b() {} // line 8\n"+
				"}\n"+
				"class Qux {\n"+
				"  b() {\n"+
				"    // multiple side-effectful lines doesn't cause duplicate error message lines.\n"+
				"    if (cond) {\n"+
				"      throw 1;\n"+
				"    } else {\n"+
				"      throw 2;\n"+
				"    }\n"+
				"  } // line 11\n"+
				"}\n"+
				"x.b();\n"),
		c.Error(PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION).
			WithMessageContaining(
				"Side-effectful definitions:\n"+
					"  testcode:5:2\n"+
					"  testcode:11:2"))
}

// Java: PureFunctionIdentifierTest.testConstructorPropertyWithSideEffects_skipDebuggingError()
func (c *PureFunctionIdentifierTest) TestConstructorPropertyWithSideEffects_skipDebuggingError() {
	c.enableArtificialPurityDebugError = true
	c.TestNoWarning(
		c.Externs("Object.prototype.constructor = function() {};"),
		c.Srcs(
			"class Foo {\n"+
				"  /** @nosideeffects */\n"+
				"  constructor() { throw 0; }\n"+
				"}\n"+
				"new Foo();\n"))
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsOnConstructor()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsOnConstructor() {
	c.AssertPureCallsMarked(
		"class Foo {\n"+
			"  /** @nosideeffects */ constructor() { throw 0; }\n"+
			"}\n"+
			"new Foo();\n",
		[]string{"Foo"})
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsOnConstructor_ambiguousWithOtherSideEffects()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsOnConstructor_ambiguousWithOtherSideEffects() {
	c.AssertPureCallsMarked(
		"class Foo {\n"+
			"  /** @nosideeffects */ constructor() { throw 0; }\n"+
			"}\n"+
			"class Bar {\n"+
			"  constructor() { throw 0; }\n"+
			"}\n"+
			"new Foo();\n"+
			"new Bar();\n",
		[]string{"Foo"})
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsNoPropagationFromOtherCalls()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsNoPropagationFromOtherCalls() {
	source :=
		"var hasSideEffects = function() { throw 0; }\n" +
			"var foo = /** @nosideeffects */ function(){ hasSideEffects(); return 0; };\n" +
			"var bar = function(){ hasSideEffects(); return 0; };\n" +
			"foo();\n" +
			"bar();\n"
	// bar is not marked as pure because it inherits side effects from hasSideEffects
	c.AssertPureCallsMarked(source, []string{"foo"})
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsNoPropagationToOtherCalls()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsNoPropagationToOtherCalls() {
	source :=
		"var liesAboutSideEffects = /** @nosideeffects */ function() { throw 0; }\n" +
			"var foo = function(){ liesAboutSideEffects(); return 0; };\n" +
			"var bar = function(){ liesAboutSideEffects(); return 0; };\n" +
			"foo();\n" +
			"bar();\n"
	c.AssertPureCallsMarked(
		source, []string{"liesAboutSideEffects", "liesAboutSideEffects", "foo", "bar"})
}

// Java: PureFunctionIdentifierTest.testAnnotationInExternStubs1()
func (c *PureFunctionIdentifierTest) TestAnnotationInExternStubs1() {
	// In the case where a property is defined first as a stub and then with a FUNCTION:
	// we have to make a conservative assumption about the behaviour of the extern, since the stub
	// carried no information.
	c.AssertNoPureCalls("o.propWithStubBefore('a');")
}

// Java: PureFunctionIdentifierTest.testAnnotationInExternStubs1b()
func (c *PureFunctionIdentifierTest) TestAnnotationInExternStubs1b() {
	// In the case where a property is defined first as a stub and then with a FUNCTION:
	// we have to make a conservative assumption about the behaviour of the extern, since the stub
	// carried no information.
	c.AssertNoPureCalls("o.propWithStubBeforeWithJSDoc('a');")
}

// Java: PureFunctionIdentifierTest.testAnnotationInExternStubs2()
func (c *PureFunctionIdentifierTest) TestAnnotationInExternStubs2() {
	// In the case where a property is defined first with a FUNCTION and then as a stub:
	// we have to make a conservative assumption about the behaviour of the extern, since the stub
	// carried no information.
	c.AssertNoPureCalls("o.propWithStubAfter('a');")
}

// Java: PureFunctionIdentifierTest.testAnnotationInExternStubs3()
func (c *PureFunctionIdentifierTest) TestAnnotationInExternStubs3() {
	// In the case where a property is defined first with a FUNCTION and then as a stub:
	// we have to make a conservative assumption about the behaviour of the extern, since the stub
	// carried no information.
	c.AssertNoPureCalls("propWithAnnotatedStubAfter('a');")
}

// Java: PureFunctionIdentifierTest.testAnnotationInExternStubs4()
func (c *PureFunctionIdentifierTest) TestAnnotationInExternStubs4() {
	// An externs definition with a stub that differs from the declaration.
	// Verify our assumption is valid about this.
	externs :=
		"/**@constructor*/function externObj5(){}\n" +
			"externObj5.prototype.propWithAnnotatedStubAfter = function(s) {};\n" +
			"/**\n" +
			" * @param {string} s id.\n" +
			" * @return {string}\n" +
			" * @nosideeffects\n" +
			" */\n" +
			"externObj5.prototype.propWithAnnotatedStubAfter;\n"

	c.EnableTypeCheck()
	c.TestSame(c.Externs(externs), c.Srcs("o.prototype.propWithAnnotatedStubAfter"))
	AssertThat(c.noSideEffectCalls).IsEmpty()
}

// Java: PureFunctionIdentifierTest.testAnnotationInExternStubs5()
func (c *PureFunctionIdentifierTest) TestAnnotationInExternStubs5() {
	// An externs definition with a stub that differs from the declaration.
	// Verify our assumption is valid about this.
	externs :=
		"/**@constructor*/function externObj5(){}\n" +
			"/**\n" +
			" * @param {string} s id.\n" +
			" * @return {string}\n" +
			" * @nosideeffects\n" +
			" */\n" +
			"externObj5.prototype.propWithAnnotatedStubAfter = function(s) {};\n" +
			"/**\n" +
			" * @param {string} s id.\n" +
			" * @return {string}\n" +
			" */\n" +
			"externObj5.prototype.propWithAnnotatedStubAfter;\n"

	c.EnableTypeCheck()
	c.TestSame(c.Externs(externs), c.Srcs("o.prototype.propWithAnnotatedStubAfter"))
	AssertThat(c.noSideEffectCalls).IsEmpty()
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsSimple()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsSimple() {
	prefix := "function f(){"
	suffix := "} f()"
	expected := []string{"f"}

	c.AssertPureCallsMarked(prefix+""+suffix, expected)
	c.AssertPureCallsMarked(prefix+"return 1"+suffix, expected)
	c.AssertPureCallsMarked(prefix+"return 1 + 2"+suffix, expected)

	// local var
	c.AssertPureCallsMarked(prefix+"var a = 1; return a"+suffix, expected)

	// mutate local var
	c.AssertPureCallsMarked(prefix+"var a = 1; a = 2; return a"+suffix, expected)
	c.AssertPureCallsMarked(prefix+"var a = 1; a = 2; return a + 1"+suffix, expected)

	// read from obj literal
	c.AssertPureCallsMarked(prefix+"var a = {foo : 1}; return a.foo"+suffix, expected)
	c.AssertPureCallsMarked(prefix+"var a = {foo : 1}; return a.foo + 1"+suffix, expected)

	// read from extern
	c.AssertPureCallsMarked(prefix+"return externObj"+suffix, expected)
	c.AssertPureCallsMarked(
		"function g(x) { x.foo = 3; }" /* to suppress missing property */ +
			prefix+
			"return externObj.foo"+
			suffix,
		expected)
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsSimple2()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsSimple2() {
	c.regExpHaveSideEffects = false
	source :=
		"function f() {\n" +
			"  return ''.replace(/xyz/g, '');\n" +
			"}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"''.replace", "f"})
}

// Java: PureFunctionIdentifierTest.testNoSideEffectsSimple3()
func (c *PureFunctionIdentifierTest) TestNoSideEffectsSimple3() {
	c.regExpHaveSideEffects = false
	source :=
		"function f(/** string */ x) {\n" +
			"  return x.replace(/xyz/g, '');\n" +
			"}\n" +
			"f('')\n"
	c.AssertPureCallsMarked(source, []string{"x.replace", "f"})
}

// Java: PureFunctionIdentifierTest.testExternCalls()
func (c *PureFunctionIdentifierTest) TestExternCalls() {
	c.TestExternCallsForTypeInferenceMode( /* typeChecked= */ true)
}

// Java: PureFunctionIdentifierTest.testExternCallsNoTypeChecking()
func (c *PureFunctionIdentifierTest) TestExternCallsNoTypeChecking() {
	c.TestExternCallsForTypeInferenceMode( /* typeChecked= */ false)
}

// Java: PureFunctionIdentifierTest.testExternCallsForTypeInferenceMode(boolean)
func (c *PureFunctionIdentifierTest) TestExternCallsForTypeInferenceMode(typeChecked bool) {
	if typeChecked {
		c.EnableTypeCheck()
	} else {
		c.DisableTypeCheck()
	}
	prefix := "function f(){"
	suffix := "} f()"

	c.AssertPureCallsMarked(prefix+"externNsef1()"+suffix, []string{"externNsef1", "f"})
	c.AssertPureCallsMarked(
		prefix+"externObj.nsef1()"+suffix, []string{"externObj.nsef1", "f"})

	c.AssertNoPureCalls(prefix + "externSef1()" + suffix)
	c.AssertNoPureCalls(prefix + "externObj.sef1()" + suffix)
}

// Java: PureFunctionIdentifierTest.testApply()
func (c *PureFunctionIdentifierTest) TestApply() {
	source :=
		"function f() {return 42}\n" +
			"f.apply(null)\n"
	c.AssertPureCallsMarked(source, []string{"f.apply"})
}

// Java: PureFunctionIdentifierTest.testCall()
func (c *PureFunctionIdentifierTest) TestCall() {
	source :=
		"function f() {return 42}\n" +
			"f.call(null)\n"
	c.AssertPureCallsMarked(source, []string{"f.call"})
}

// Java: PureFunctionIdentifierTest.testApplyToUnknownDefinition()
func (c *PureFunctionIdentifierTest) TestApplyToUnknownDefinition() {
	source :=
		"var dict = {'func': function() {}};\n" +
			"function f() { var s = dict['func'];}\n" +
			"f.apply(null)\n"
	c.AssertPureCallsMarked(source, []string{"f.apply"})

	// Not marked because the definition cannot be found so unknown side effects.
	source =
		"var dict = {'func': function() {}};\n" +
			"function f() { var s = dict['func'].apply(null); }\n" +
			"f.apply(null)\n"
	c.AssertNoPureCalls(source)

	// Not marked because the definition cannot be found so unknown side effects.
	source =
		"var pure = function() {};\n" +
			"var dict = {'func': function() {}};\n" +
			"function f() { var s = (dict['func'] || pure)();}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)

	// Not marked because the definition cannot be found so unknown side effects.
	source =
		"var pure = function() {};\n" +
			"var dict = {'func': function() {}};\n" +
			"function f() { var s = (condition ? dict['func'] : pure)();}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testOptChainCall()
func (c *PureFunctionIdentifierTest) TestOptChainCall() {
	source :=
		"function f() {return g?.()}\n" +
			"function g() {return 42}\n" +
			"f?.()\n"
	c.AssertPureCallsMarked(source, []string{"g", "f"})
}

// Java: PureFunctionIdentifierTest.testInference1()
func (c *PureFunctionIdentifierTest) TestInference1() {
	source :=
		"function f() {return g()}\n" +
			"function g() {return 42}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"g", "f"})
}

// Java: PureFunctionIdentifierTest.testInference2()
func (c *PureFunctionIdentifierTest) TestInference2() {
	source :=
		"var a = 1;\n" +
			"function f() {g()}\n" +
			"function g() {a=2}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testInference3()
func (c *PureFunctionIdentifierTest) TestInference3() {
	source :=
		"var f = function() {return g()};\n" +
			"var g = function() {return 42};\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"g", "f"})
}

// Java: PureFunctionIdentifierTest.testInference4()
func (c *PureFunctionIdentifierTest) TestInference4() {
	source :=
		"var a = 1; var f = function() {g()};\n" +
			"var g = function() {a=2};\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testInference5()
func (c *PureFunctionIdentifierTest) TestInference5() {
	source :=
		"goog.f = function() {return goog.g()};\n" +
			"goog.g = function() {return 42};\n" +
			"goog.f()\n"
	c.AssertPureCallsMarked(source, []string{"goog.g", "goog.f"})
}

// Java: PureFunctionIdentifierTest.testInference6()
func (c *PureFunctionIdentifierTest) TestInference6() {
	source :=
		"var a = 1;\n" +
			"goog.f = function() {goog.g()};\n" +
			"goog.g = function() {a=2};\n" +
			"goog.f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects1()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects1() {
	// Returning a function that contains a modification of a local
	// is not a global side-effect.
	source :=
		"function f() {\n" +
			"  var x = {foo : 0}; return function() {x.foo++};\n" +
			"}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects2()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects2() {
	// Calling a function that contains a modification of a local
	// is a global side-effect (the value has escaped).
	source :=
		"function f() {\n" +
			"  var x = {foo : 0}; (function() {x.foo++})();\n" +
			"}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects3()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects3() {
	// A local that might be assigned a global value and whose properties
	// are modified must be considered a global side-effect.
	source :=
		"var g = {foo:1};\n" +
			"function f() {var x = g; x.foo++};\n" +
			"f();\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects4()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects4() {
	// An array is a local object, assigning a local array is not a global side-effect.
	c.AssertPureCallsMarked(
		"function f() {var x = []; x[0] = 1;} // preserve newline\n"+
			"f()\n",
		[]string{"f"})

	c.AssertPureCallsMarked(
		"function f() {const x = []; x[0] = 1;} // preserve newline\n"+
			"f()\n",
		[]string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects5()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects5() {
	// Assigning a local alias of a global is a global
	// side-effect.
	source :=
		"var g = [];\n" +
			"function f() {var x = g; x[0] = 1;};\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects6()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects6() {
	// Returning a local object that has been modified
	// is not a global side-effect.
	c.AssertPureCallsMarked(
		"function f() { // preserve newline\n"+
			"  var x = {}; x.foo = 1; return x;\n"+
			"}\n"+
			"f()\n",
		[]string{"f"})

	c.AssertPureCallsMarked(
		"function f() { // preserve newline\n"+
			"  const x = {}; x.foo = 1; return x;\n"+
			"}\n"+
			"f()\n",
		[]string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects7()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects7() {
	// Returning a local object that has been modified
	// is not a global side-effect.
	c.AssertPureCallsMarked(
		"/** @constructor A */ function A() {};\n"+
			"function f() {\n"+
			"  var a = []; a[1] = 1; return a;\n"+
			"}\n"+
			"f()\n",
		[]string{"f"})

	c.AssertPureCallsMarked(
		"/** @constructor A */ function A() {};\n"+
			"function f() {\n"+
			"  const a = []; a[1] = 1; return a;\n"+
			"}\n"+
			"f()\n",
		[]string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects8()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects8() {
	// Returning a local object that has been modified
	// is not a global side-effect.
	// TODO(tdeegan): Not yet. Propagate local object information.
	source :=
		"/** @constructor A */ function A() {};\n" +
			"function f() {\n" +
			"  var a = new A; a.foo = 1; return a;\n" +
			"}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"A"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects9()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects9() {
	// Returning a local object that has been modified
	// is not a global side-effect.
	// TODO(johnlenz): Not yet. Propagate local object information.
	source :=
		"/** @constructor A */ function A() {this.x = 1};\n" +
			"function f() {\n" +
			"  var a = new A; a.foo = 1; return a;\n" +
			"}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"A"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects10()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects10() {
	// Returning a local object that has been modified
	// is not a global side-effect.
	source :=
		"/** @constructor A */ function A() {};\n" +
			"A.prototype.g = function() {this.x = 1};\n" +
			"function f() {\n" +
			"  var a = new A; a.g(); return a;\n" +
			"}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"A"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects11()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects11() {
	// TODO(tdeegan): updateA is side effect free.
	// Calling a function of a local object that taints this.
	source :=
		"/** @constructor */\n" +
			"function A() {}\n" +
			"A.prototype.update = function() { this.x = 1; };\n" +
			"\n" +
			"/** @constructor */\n" +
			"function B() {\n" +
			"  this.a_ = new A();\n" +
			"}\n" +
			"B.prototype.updateA = function() {\n" +
			"  var b = this.a_;\n" +
			"  b.update();\n" +
			"};\n" +
			"\n" +
			"var x = new B();\n" +
			"x.updateA();\n"
	c.AssertPureCallsMarked(source, []string{"A", "B"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects12()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects12() {
	// An array is a local object, assigning a local array is not a global
	// side-effect. This tests the behavior if the access is in a block scope.
	c.AssertPureCallsMarked(
		"function f() {var x = []; { x[0] = 1; } } // preserve newline\n"+
			"f()\n",
		[]string{"f"})

	c.AssertPureCallsMarked(
		"function f() {const x = []; { x[0] = 1; } } // preserve newline\n"+
			"f()\n",
		[]string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects13()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects13() {
	source :=
		"function f() {var [x, y] = [3, 4]; }\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects14()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects14() {
	source :=
		"function f() {var x; if (true) { [x] = [5]; } }\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects15()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects15() {
	source :=
		"function f() {var {length} = 'a string'; }\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects16()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects16() {
	source :=
		"function f(someArray) {var [a, , b] = someArray; }\n" +
			"f([])\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects17()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects17() {

	source :=
		"function f(someObj) {var { very: { nested: { lhs: pattern }} } = someObj; }\n" +
			"f({very: {nested: {lhs: 0}}})\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects18()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects18() {
	source :=
		"/** @constructor */\n" +
			"function SomeCtor() { [this.x, this.y] = getCoordinates(); }\n" +
			"new SomeCtor()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects19()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects19() {
	source :=
		"/** @constructor */\n" +
			"function SomeCtor() { [this.x, this.y] = [0, 1]; }\n" +
			"new SomeCtor()\n"
	c.AssertPureCallsMarked(source, []string{"SomeCtor"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects20()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects20() {
	source :=
		"/** @constructor */\n" +
			"function SomeCtor() { this.x += 1; }\n" +
			"new SomeCtor()\n"
	c.AssertPureCallsMarked(source, []string{"SomeCtor"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects21()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects21() {
	source :=
		"function f(values) { const x = {}; [x.y, x.z] = values; }\n" +
			"f([])\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects22()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects22() {
	source :=
		"var x = {};\n" +
			"function f(values) { [x.y, x.z] = values; }\n" +
			"f([])\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects23()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects23() {
	source :=
		"function f(values) { const x = {}; [x.y, x.z = defaultNoSideEffects] = values; }\n" +
			"f([])\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testLocalizedSideEffects24()
func (c *PureFunctionIdentifierTest) TestLocalizedSideEffects24() {
	source :=
		"function f(values) { var x = {}; [x.y, x.z = defaultWithSideEffects()] = values; }\n" +
			"f([])\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testUnaryOperators1()
func (c *PureFunctionIdentifierTest) TestUnaryOperators1() {
	source :=
		"function f() {var x = 1; x++}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testUnaryOperators2()
func (c *PureFunctionIdentifierTest) TestUnaryOperators2() {
	source :=
		"var x = 1;\n" +
			"function f() {x++}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testUnaryOperators3()
func (c *PureFunctionIdentifierTest) TestUnaryOperators3() {
	c.AssertPureCallsMarked(
		"function f() {var x = {foo : 0}; x.foo++} // preserve newline\n"+
			"f()\n",
		[]string{"f"})

	c.AssertPureCallsMarked(
		"function f() {const x = {foo : 0}; x.foo++} // preserve newline\n"+
			"f()\n",
		[]string{"f"})
}

// Java: PureFunctionIdentifierTest.testUnaryOperators4()
func (c *PureFunctionIdentifierTest) TestUnaryOperators4() {
	source :=
		"var x = {foo : 0};\n" +
			"function f() {x.foo++}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testUnaryOperators5()
func (c *PureFunctionIdentifierTest) TestUnaryOperators5() {
	c.AssertPureCallsMarked(
		"function f(x) {x.foo++} // preserve newline\n"+
			"f({foo : 0})\n",
		[]string{"f"})

	c.AssertNoPureCalls(
		"function f({x}) {x.foo++} // preserve newline\n" +
			"f({x: {foo : 0}})\n")
}

// Java: PureFunctionIdentifierTest.testDeleteOperator1()
func (c *PureFunctionIdentifierTest) TestDeleteOperator1() {
	source :=
		"var x = {};\n" +
			"function f() {delete x}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testDeleteOperator1_getProp()
func (c *PureFunctionIdentifierTest) TestDeleteOperator1_getProp() {
	source :=
		"var x = {y:2};\n" +
			"function f() {delete x.y}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testDeleteOperator1_optChainGetProp()
func (c *PureFunctionIdentifierTest) TestDeleteOperator1_optChainGetProp() {
	source :=
		"var x = {y:2};\n" +
			"function f() {delete x?.y}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testDeleteOperator2()
func (c *PureFunctionIdentifierTest) TestDeleteOperator2() {
	source :=
		"function f() {var x = {}; delete x}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testDeleteOperator2_getProp()
func (c *PureFunctionIdentifierTest) TestDeleteOperator2_getProp() {
	source :=
		"function f() {var x = {y:2}; delete x.y}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testDeleteOperator2_optChainGetProp()
func (c *PureFunctionIdentifierTest) TestDeleteOperator2_optChainGetProp() {
	source :=
		"function f() {var x = {y:2}; delete x?.y}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Start pure alias

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_pureAlias_pureAlias()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_pureAlias_pureAlias() {
	c.AssertCallableExpressionPure(true, "externNsef1 || externNsef2")
	c.AssertCallableExpressionPure(true, "externNsef1 && externNsef2")
	c.AssertCallableExpressionPure(true, "externNsef1, externNsef2")
	c.AssertCallableExpressionPure(true, "x ? externNsef1 : externNsef2")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_pureAlias_pureLiteral()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_pureAlias_pureLiteral() {
	c.AssertCallableExpressionPure(true, "externNsef1 || function() { }")
	c.AssertCallableExpressionPure(true, "externNsef1 && function() { }")
	c.AssertCallableExpressionPure(true, "externNsef1, function() { }")
	c.AssertCallableExpressionPure(true, "x ? externNsef1 : function() { }")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_pureAlias_impureAlias()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_pureAlias_impureAlias() {
	c.AssertCallableExpressionPure(false, "externNsef1 || externSef2")
	c.AssertCallableExpressionPure(false, "externNsef1 && externSef2")
	c.AssertCallableExpressionPure(false, "externNsef1, externSef2")
	c.AssertCallableExpressionPure(false, "x ? externNsef1 : externSef2")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_pureAlias_impureLiteral()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_pureAlias_impureLiteral() {
	c.AssertCallableExpressionPure(false, "externNsef1 || function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "externNsef1 && function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "externNsef1, function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "x ? externNsef1 : function() { throw 0; }")
}

// Java: PureFunctionIdentifierTest.nullishCoalesce_pureAlias()
func (c *PureFunctionIdentifierTest) TestNullishCoalesce_pureAlias() {
	c.AssertCallableExpressionPure(true, "externNsef1 ?? externNsef2")
	c.AssertCallableExpressionPure(true, "externNsef1 ?? function() { }")
	c.AssertCallableExpressionPure(false, "externNsef1 ?? externSef2")
	c.AssertCallableExpressionPure(false, "externNsef1 ?? function() { throw 0; }")
}

// End pure alias, start pure literal

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_pureLiteral_pureAlias()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_pureLiteral_pureAlias() {
	c.AssertCallableExpressionPure(true, "function() { } || externNsef2")
	c.AssertCallableExpressionPure(true, "function() { } && externNsef2")
	c.AssertCallableExpressionPure(true, "function() { }, externNsef2")
	c.AssertCallableExpressionPure(true, "x ? function() { } : externNsef2")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_pureLiteral_pureLiteral()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_pureLiteral_pureLiteral() {
	c.AssertCallableExpressionPure(true, "function() { } || function() { }")
	c.AssertCallableExpressionPure(true, "function() { } && function() { }")
	c.AssertCallableExpressionPure(true, "function() { }, function() { }")
	c.AssertCallableExpressionPure(true, "x ? function() { } : function() { }")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_pureLiteral_impureAlias()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_pureLiteral_impureAlias() {
	c.AssertCallableExpressionPure(false, "function() { } || externSef2")
	c.AssertCallableExpressionPure(false, "function() { } && externSef2")
	c.AssertCallableExpressionPure(false, "function() { }, externSef2")
	c.AssertCallableExpressionPure(false, "x ? function() { } : externSef2")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_pureLiteral_impureLiteral()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_pureLiteral_impureLiteral() {
	c.AssertCallableExpressionPure(false, "function() { } || function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "function() { } && function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "function() { }, function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "x ? function() { } : function() { throw 0; }")
}

// Java: PureFunctionIdentifierTest.nullishCoalesce_pureLiteral()
func (c *PureFunctionIdentifierTest) TestNullishCoalesce_pureLiteral() {
	c.AssertCallableExpressionPure(true, "function() { } ?? externNsef2")
	c.AssertCallableExpressionPure(true, "function() { } ?? function() { }")
	c.AssertCallableExpressionPure(false, "function() { } ?? externSef2")
	c.AssertCallableExpressionPure(false, "function() { } ?? function() { throw 0; }")
}

// End pure literal, start impure alias

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_impureAlias_pureAlias()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_impureAlias_pureAlias() {
	c.AssertCallableExpressionPure(false, "externSef1 || externNsef2")
	c.AssertCallableExpressionPure(false, "externSef1 && externNsef2")
	c.AssertCallableExpressionPure(true, "externSef1, externNsef2")
	c.AssertCallableExpressionPure(false, "x ? externSef1 : externNsef2")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_impureAlias_pureLiteral()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_impureAlias_pureLiteral() {
	c.AssertCallableExpressionPure(false, "externSef1 || function() { }")
	c.AssertCallableExpressionPure(false, "externSef1 && function() { }")
	c.AssertCallableExpressionPure(true, "externSef1, function() { }")
	c.AssertCallableExpressionPure(false, "x ? externSef1 : function() { }")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_impureAlias_impureAlias()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_impureAlias_impureAlias() {
	c.AssertCallableExpressionPure(false, "externSef1 || externSef2")
	c.AssertCallableExpressionPure(false, "externSef1 && externSef2")
	c.AssertCallableExpressionPure(false, "externSef1, externSef2")
	c.AssertCallableExpressionPure(false, "x ? externSef1 : externSef2")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_impureAlias_impureLiteral()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_impureAlias_impureLiteral() {
	c.AssertCallableExpressionPure(false, "externSef1 || function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "externSef1 && function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "externSef1, function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "x ? externSef1 : function() { throw 0; }")
}

// Java: PureFunctionIdentifierTest.nullishCoalesce_impureAlias()
func (c *PureFunctionIdentifierTest) TestNullishCoalesce_impureAlias() {
	c.AssertCallableExpressionPure(false, "externSef1 ?? externNsef2")
	c.AssertCallableExpressionPure(false, "externSef1 ?? function() { }")
	c.AssertCallableExpressionPure(false, "externSef1 ?? externSef2")
	c.AssertCallableExpressionPure(false, "externSef1 ?? function() { throw 0; }")
}

// End impure alias, start impure literal

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_impureLiteral_pureAlias()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_impureLiteral_pureAlias() {
	c.AssertCallableExpressionPure(false, "function() { throw 0; } || externNsef2")
	c.AssertCallableExpressionPure(false, "function() { throw 0; } && externNsef2")
	c.AssertCallableExpressionPure(true, "function() { throw 0; }, externNsef2")
	c.AssertCallableExpressionPure(false, "x ? function() { throw 0; } : externNsef2")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_impureLiteral_pureLiteral()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_impureLiteral_pureLiteral() {
	c.AssertCallableExpressionPure(false, "function() { throw 0; } || function() { }")
	c.AssertCallableExpressionPure(false, "function() { throw 0; } && function() { }")
	c.AssertCallableExpressionPure(true, "function() { throw 0; }, function() { }")
	c.AssertCallableExpressionPure(false, "x ? function() { throw 0; } : function() { }")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_impureLiteral_impureAlias()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_impureLiteral_impureAlias() {
	c.AssertCallableExpressionPure(false, "function() { throw 0; } || externSef2")
	c.AssertCallableExpressionPure(false, "function() { throw 0; } && externSef2")
	c.AssertCallableExpressionPure(false, "function() { throw 0; }, externSef2")
	c.AssertCallableExpressionPure(false, "x ? function() { throw 0; } : externSef2")
}

// Java: PureFunctionIdentifierTest.testOperatorsTraversal_impureLiteral_impureLiteral()
func (c *PureFunctionIdentifierTest) TestOperatorsTraversal_impureLiteral_impureLiteral() {
	c.AssertCallableExpressionPure(false, "function() { throw 0; } || function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "function() { throw 0; } && function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "function() { throw 0; }, function() { throw 0; }")
	c.AssertCallableExpressionPure(false, "x ? function() { throw 0; } : function() { throw 0; }")
}

// Java: PureFunctionIdentifierTest.nullishCoalesce_impureLiteral()
func (c *PureFunctionIdentifierTest) TestNullishCoalesce_impureLiteral() {
	c.AssertCallableExpressionPure(false, "function() { throw 0; } ?? externNsef2")
	c.AssertCallableExpressionPure(false, "function() { throw 0; } ?? function() { }")
	c.AssertCallableExpressionPure(false, "function() { throw 0; } ?? externSef2")
	c.AssertCallableExpressionPure(false, "function() { throw 0; } ?? function() { throw 0; }")
}

// End impure literal

// Java: PureFunctionIdentifierTest.testKnownDynamicFunctions_areSkiplistedForAliasing()
func (c *PureFunctionIdentifierTest) TestKnownDynamicFunctions_areSkiplistedForAliasing() {
	for _, name := range []string{"call", "apply", "constructor"} {
		c.AssertNoPureCalls(
			// Create a known pure definition so we're sure the name isn't just undefined.
			jutil.StringReplace("foo.NAME = function() { };\n"+
				"const f = foo.NAME;\n"+
				"f();\n", "NAME", name))
	}
}

// Java: PureFunctionIdentifierTest.testThrow1()
func (c *PureFunctionIdentifierTest) TestThrow1() {
	source :=
		"function f(){throw Error()};\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"Error"})
}

// Java: PureFunctionIdentifierTest.testThrow2()
func (c *PureFunctionIdentifierTest) TestThrow2() {
	source :=
		"/**@constructor*/function A(){throw Error()};\n" +
			"function f(){return new A()}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"Error"})
}

// Java: PureFunctionIdentifierTest.testThrow_insideTry_withCatch_hasNoSideEffects()
func (c *PureFunctionIdentifierTest) TestThrow_insideTry_withCatch_hasNoSideEffects() {
	source :=
		"function f() {\n" +
			"  try {\n" +
			"    throw Error();\n" +
			"  } catch (e) {\n" +
			"  } finally {\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"f();\n"
	c.AssertPureCallsMarked(source, []string{"f", "Error"})
}

// Java: PureFunctionIdentifierTest.testThrow_insideTry_withoutCatch_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestThrow_insideTry_withoutCatch_hasSideEffects() {
	source :=
		"function f() {\n" +
			"  try {\n" +
			"    throw Error();\n" +
			"  } finally {\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"f();\n"
	c.AssertPureCallsMarked(source, []string{"Error"})
}

// Java: PureFunctionIdentifierTest.testThrow_insideCatch_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestThrow_insideCatch_hasSideEffects() {
	source :=
		"function f() {\n" +
			"  try {\n" +
			"  } catch (e) {\n" +
			"    throw Error();\n" +
			"  } finally {\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"f();\n"
	c.AssertPureCallsMarked(source, []string{"Error"})
}

// Java: PureFunctionIdentifierTest.testThrow_insideFinally_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestThrow_insideFinally_hasSideEffects() {
	source :=
		"function f() {\n" +
			"  try {\n" +
			"  } catch (e) {\n" +
			"  } finally {\n" +
			"    throw Error();\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"f();\n"
	c.AssertPureCallsMarked(source, []string{"Error"})
}

// Java: PureFunctionIdentifierTest.testThrow_insideTry_afterNestedTry_hasNoSideEffects()
func (c *PureFunctionIdentifierTest) TestThrow_insideTry_afterNestedTry_hasNoSideEffects() {
	/* Ensure we track the stack of trys correctly, rather than just toggling a boolean. */
	source :=
		"function f() {\n" +
			"  try {\n" +
			"    try {\n" +
			"    } finally {\n" +
			"    }\n" +
			"\n" +
			"    throw Error();\n" +
			"  } catch (e) {\n" +
			"  } finally {\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"f();\n"
	c.AssertPureCallsMarked(source, []string{"f", "Error"})
}

// Java: PureFunctionIdentifierTest.testThrow_insideFunction_insideTry_hasNoSideEffects()
func (c *PureFunctionIdentifierTest) TestThrow_insideFunction_insideTry_hasNoSideEffects() {
	/* Ensure we track the stack of trys correctly, rather than just toggling a boolean. */
	source :=
		"function f() {\n" +
			"  try {\n" +
			"    function g() {\n" +
			"      throw Error();\n" +
			"    }\n" +
			"  } catch (e) {\n" +
			"  } finally {\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"f();\n"
	c.AssertPureCallsMarked(source, []string{"f", "Error"})
}

// Java: PureFunctionIdentifierTest.testThrow_fromCallee_insideTry_withCatch_hasNoSideEffects()
func (c *PureFunctionIdentifierTest) TestThrow_fromCallee_insideTry_withCatch_hasNoSideEffects() {
	source :=
		"function h() {\n" +
			"  throw Error();\n" +
			"}\n" +
			"\n" +
			"function f() {\n" +
			"  try {\n" +
			"    h()\n" +
			"  } catch (e) {\n" +
			"  } finally {\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"h();\n" +
			"f();\n"
	c.AssertPureCallsMarked(source, []string{"f", "Error"})
}

// Java: PureFunctionIdentifierTest.testAssignmentOverride()
func (c *PureFunctionIdentifierTest) TestAssignmentOverride() {
	source :=
		"/**@constructor*/function A(){}\n" +
			"A.prototype.foo = function(){};\n" +
			"var a = new A;\n" +
			"a.foo();\n"
	c.AssertPureCallsMarked(source, []string{"A", "a.foo"})

	// Ideally inline aliases takes care of this.
	sourceOverride :=
		"/**@constructor*/ function A(){}\n" +
			"A.prototype.foo = function(){};\n" +
			"var x = 1\n" +
			"function f(){x = 10}\n" +
			"var a = new A;\n" +
			"a.foo = f;\n" +
			"a.foo();\n"
	c.AssertPureCallsMarked(sourceOverride, []string{"A"})
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitions()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitions() {
	source :=
		CompilerTypeTestCase_CLOSURE_DEFS +
			"var globalVar = 1;\n" +
			"A.f = function() {globalVar = 2;};\n" +
			"A.f = function() {};\n" +
			"function sideEffectCaller() { A.f() };\n" +
			"sideEffectCaller();\n"
	// Can't tell which f is being called so it assumes both.
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsCall()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsCall() {
	source :=
		CompilerTypeTestCase_CLOSURE_DEFS +
			"var globalVar = 1;\n" +
			"A.f = function() {globalVar = 2;};\n" +
			"A.f = function() {};\n" +
			"function sideEffectCaller() { A.f.call(null); };\n" +
			"sideEffectCaller();\n"
	// Can't tell which f is being called so it assumes both.
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsAllPropagationTypes()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsAllPropagationTypes() {
	source :=
		CompilerTypeTestCase_CLOSURE_DEFS +
			"/**@constructor*/A.f = function() { this.x = 5; };\n" +
			"/**@constructor*/B.f = function() {};\n" +
			"function sideEffectCaller() { new C.f() };\n" +
			"sideEffectCaller();\n"
	// "f" is only known to mutate `this`, but `this` is a local object when invoked with `new`.
	c.AssertPureCallsMarked(source, []string{"C.f", "sideEffectCaller"})
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsCallWithThis()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsCallWithThis() {
	source :=
		CompilerTypeTestCase_CLOSURE_DEFS +
			"A.modifiesThis = function() { this.x = 5; };\n" +
			"B.modifiesThis = function() {};\n" +
			"\n" +
			"/** @constructor */ function Constructor() { B.modifiesThis.call(this); };\n" +
			"\n" +
			"new Constructor();\n" +
			"A.modifiesThis();\n"

	// Can't tell which modifiesThis is being called so it assumes both.
	c.AssertPureCallsMarked(source, []string{"Constructor"})
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsBothCallThis()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsBothCallThis() {
	source :=
		"B.f = function() {\n" +
			"  this.x = 1;\n" +
			"}\n" +
			"/** @constructor */ function C() {\n" +
			"  this.f.apply(this);\n" +
			"}\n" +
			"C.prototype.f = function() {\n" +
			"  this.x = 2;\n" +
			"}\n" +
			"new C();\n"
	c.AssertPureCallsMarked(source, []string{"C"})
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsBothCallThisWithOptionalChain()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsBothCallThisWithOptionalChain() {
	source :=
		"B.f = function() {\n" +
			"  this.x = 1;\n" +
			"}\n" +
			"/** @constructor */ function C() {\n" +
			"  this.f?.apply(this);\n" +
			"}\n" +
			"C.prototype.f = function() {\n" +
			"  this.x = 2;\n" +
			"}\n" +
			"new C();\n"
	c.AssertPureCallsMarked(source, []string{"C"})
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsAllCallThis()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsAllCallThis() {
	source :=
		"A.f = function() { this.y = 1 };\n" +
			"C.f = function() { };\n" +
			"var g = function() {D.f()};\n" +
			"/** @constructor */ var h = function() {E.f.apply(this)};\n" +
			"// TODO(nickreid): Detect that `{}` being passed as `this` is local. With that\n" +
			"// understanding, we could determine that `i` has no side-effects.\n" +
			"var i = function() { F.f.apply({}); };\n" +
			"g();\n" +
			"new h();\n" +
			"i();\n"
	c.AssertPureCallsMarked(source, []string{"F.f.apply", "h"})
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsMutatesGlobalArgument()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsMutatesGlobalArgument() {
	source :=
		"// Mutates argument\n" +
			"A.a = function(argument) {\n" +
			"  argument.x = 2;\n" +
			"};\n" +
			"// No side effects\n" +
			"B.a = function() {};\n" +
			"var b = function(x) {C.a(x)};\n" +
			"b({});\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsMutatesLocalArgument()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsMutatesLocalArgument() {
	c.AssertPureCallsMarked(
		"// Mutates argument\n"+
			"A.a = function(argument) {\n"+
			"  argument.x = 2;\n"+
			"};\n"+
			"// No side effects\n"+
			"B.a = function() {};\n"+
			"var b = function() {\n"+
			"  C.a({});\n"+
			"};\n"+
			"b();\n",
		[]string{"C.a", "b"})
}

// Java: PureFunctionIdentifierTest.testAmbiguousExternDefinitions()
func (c *PureFunctionIdentifierTest) TestAmbiguousExternDefinitions() {
	c.AssertNoPureCalls("x.duplicateExternFunc()")

	// nsef1 is defined as no side effect in the externs.
	source :=
		"var global = 1;\n" +
			"// Overwrite the @nosideeffects with this side effect\n" +
			"A.nsef1 = function () {global = 2;};\n" +
			"externObj.nsef1();\n"
	c.AssertNoPureCalls(source)
}

/**
 * Test bug where the FunctionInformation for "A.x" and "a" were separate causing .x() calls to
 * appear pure because the global side effect was only registed for the function linked to "a".
 */
// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsDoubleDefinition()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsDoubleDefinition() {
	source :=
		"var global = 1;\n" +
			"A.x = function a() { global++; }\n" +
			"B.x = function() {}\n" +
			"B.x();\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsDoubleDefinition2()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsDoubleDefinition2() {
	source :=
		"var global = 1;\n" +
			"A.x = function a() { global++; }\n" +
			"var a = function() {} // This is the only `a` is in scope below.\n" +
			"B.x();\n" +
			"a();\n"
	c.AssertPureCallsMarked(source, []string{"a"})
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsDoubleDefinition3()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsDoubleDefinition3() {
	source :=
		"var global = 1;\n" +
			"A.x = function a() {}\n" +
			"var a = function() { global++; } // This is the only `a` is in scope below.\n" +
			"B.x();\n" +
			"a();\n"
	c.AssertPureCallsMarked(source, []string{"B.x"})
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsDoubleDefinition4()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsDoubleDefinition4() {
	source :=
		"var global = 1;\n" +
			"\n" +
			"A.x = function a() {}\n" +
			"B.x = function() { global++; }\n" +
			"\n" +
			"B.x();\n" +
			"a();\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsDoubleDefinition5()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsDoubleDefinition5() {
	source :=
		"var global = 1;\n" +
			"\n" +
			"A.x = cond ? function a() { global++ } : function b() {}\n" +
			"B.x = function() { global++; }\n" +
			"\n" +
			"B.x();\n" +
			"a(); // `a` isn't in scope here.\n" +
			"b();\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testAmbiguousDefinitionsDoubleDefinition6()
func (c *PureFunctionIdentifierTest) TestAmbiguousDefinitionsDoubleDefinition6() {
	source :=
		"var SetCustomData1 = function SetCustomData2(element, dataName, dataValue) {\n" +
			"    var x = element['_customData'];\n" +
			"    x[dataName] = dataValue;\n" +
			"}\n" +
			"SetCustomData1(window, \"foo\", \"bar\");\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testInnerFunction_isNeverPure()
func (c *PureFunctionIdentifierTest) TestInnerFunction_isNeverPure() {
	// TODO(b/129503101): `pure` should be marked pure in this case.
	c.AssertNoPureCalls(
		"function f() {\n" +
			"  function pure() { };\n" +
			"  pure();\n" +
			"}\n")
}

// Java: PureFunctionIdentifierTest.testNamedFunctionExpression_isNeverPure()
func (c *PureFunctionIdentifierTest) TestNamedFunctionExpression_isNeverPure() {
	// TODO(b/129503101): `pure` should be marked pure in this case.
	c.AssertNoPureCalls(
		"(function pure() {\n" +
			"  pure();\n" +
			"})\n")
}

// Java: PureFunctionIdentifierTest.testNameAliasPurity()
func (c *PureFunctionIdentifierTest) TestNameAliasPurity() {
	c.AssertAliasPurity(true, "const alias = pure;", false)
	c.AssertAliasPurity(true, "let alias = pure;", false)
	c.AssertAliasPurity(false, "let alias; alias = pure;", false)
	c.AssertAliasPurity(true, "var alias = pure;", false)
	c.AssertAliasPurity(false, "var alias; alias = pure;", false)
}

// Java: PureFunctionIdentifierTest.testPropertyAliasPurity()
func (c *PureFunctionIdentifierTest) TestPropertyAliasPurity() {
	c.AssertAliasPurity(true, "({alias: pure});", true)
	c.AssertAliasPurity(true, "obj.alias = pure;", true)
}

// Java: PureFunctionIdentifierTest.assertAliasPurity(boolean, String, boolean)
func (c *PureFunctionIdentifierTest) AssertAliasPurity(isPure bool, aliasDefinition string, isProperty bool) {
	var aliasExpression string
	if isProperty {
		aliasExpression = "obj.alias"
	} else {
		aliasExpression = "alias"
	}

	expected := []string{}
	expected = append(expected, "pure")
	if isPure {
		expected = append(expected, aliasExpression)
	}

	c.AssertPureCallsMarked(
		jutil.StringReplace(jutil.StringReplace("const /** ? */ obj = {};\n"+
			"function pure() { }\n"+
			"ALIAS_DEFINITION\n"+
			"\n"+
			"pure();\n"+
			"ALIAS_EXPRESSION();\n", "ALIAS_DEFINITION", aliasDefinition),
			"ALIAS_EXPRESSION", aliasExpression),
		expected)
}

// Java: PureFunctionIdentifierTest.testCallBeforeDefinition()
func (c *PureFunctionIdentifierTest) TestCallBeforeDefinition() {
	c.AssertPureCallsMarked("f(); function f(){}", []string{"f"})
	c.AssertPureCallsMarked("var a = {}; a.f(); a.f = function (){}", []string{"a.f"})
}

// Java: PureFunctionIdentifierTest.testSideEffectsArePropagated_toPrecedingCaller_fromFollowingCallee()
func (c *PureFunctionIdentifierTest) TestSideEffectsArePropagated_toPrecedingCaller_fromFollowingCallee() {
	c.AssertPureCallsMarked(
		"var following = function() { }; // This implementation is pure...\n"+
			"var preceding = function() { following(); }; // so this one should be initially...\n"+
			"following = function() { throw 'something'; }; //  but then it becomes impure.\n"+
			"\n"+
			"preceding();\n",
		[]string{})
}

// Java: PureFunctionIdentifierTest.testConstructorThatModifiesThis1()
func (c *PureFunctionIdentifierTest) TestConstructorThatModifiesThis1() {
	source :=
		"/**@constructor*/function A(){this.foo = 1}\n" +
			"function f() {return new A}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"A", "f"})
}

// Java: PureFunctionIdentifierTest.testConstructorThatModifiesThis2()
func (c *PureFunctionIdentifierTest) TestConstructorThatModifiesThis2() {
	source :=
		"/**@constructor*/function A(){this.foo()}\n" +
			"A.prototype.foo = function(){this.data=24};\n" +
			"function f() {return new A}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"A", "f"})
}

// Java: PureFunctionIdentifierTest.testOptChainCallInConstructorThatModifiesThis()
func (c *PureFunctionIdentifierTest) TestOptChainCallInConstructorThatModifiesThis() {
	source :=
		"/**@constructor*/function A(){this?.foo()}\n" +
			"A.prototype.foo = function(){this.data=24};\n" +
			"function f() {return new A}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"A", "f"})
}

// Java: PureFunctionIdentifierTest.testConstructorThatModifiesThis3()
func (c *PureFunctionIdentifierTest) TestConstructorThatModifiesThis3() {
	// test chained
	source :=
		"/**@constructor*/function A(){this.foo()}\n" +
			"A.prototype.foo = function(){this.bar()};\n" +
			"A.prototype.bar = function(){this.data=24};\n" +
			"function f() {return new A}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"A", "f"})
}

// Java: PureFunctionIdentifierTest.testConstructorThatModifiesThis4()
func (c *PureFunctionIdentifierTest) TestConstructorThatModifiesThis4() {
	source :=
		"/**@constructor*/function A(){foo.call(this)}\n" +
			"function foo(){this.data=24};\n" +
			"function f() {return new A}\n" +
			"f()\n"
	c.AssertPureCallsMarked(source, []string{"A", "f"})
}

// Java: PureFunctionIdentifierTest.testConstructorThatModifiesGlobal1()
func (c *PureFunctionIdentifierTest) TestConstructorThatModifiesGlobal1() {
	source :=
		"var b = 0;\n" +
			"/**@constructor*/function A(){b=1};\n" +
			"function f() {return new A}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testConstructorThatModifiesGlobal2()
func (c *PureFunctionIdentifierTest) TestConstructorThatModifiesGlobal2() {
	source :=
		"/**@constructor*/function A(){this.foo()}\n" +
			"A.prototype.foo = function(){b=1};\n" +
			"function f() {return new A}\n" +
			"f()\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testCallFunctionThatModifiesThis()
func (c *PureFunctionIdentifierTest) TestCallFunctionThatModifiesThis() {
	source :=
		"/**@constructor*/function A(){}\n" +
			"A.prototype.foo = function(){this.data=24};\n" +
			"function f(){var a = new A; return a}\n" +
			"function g(){var a = new A; a.foo(); return a}\n" +
			"f(); g()\n"
	c.AssertPureCallsMarked(source, []string{"A", "A", "f"})
}

// Java: PureFunctionIdentifierTest.testMutatesArguments1()
func (c *PureFunctionIdentifierTest) TestMutatesArguments1() {
	c.AssertPureCallsMarked(
		"function f(x) { x.y = 1; } // preserve newline\n"+
			"f({});\n",
		[]string{"f"})

	c.AssertPureCallsMarked(
		"function f(x = {}) { x.y = 1; } // preserve newline\n"+
			"f({});\n",
		[]string{"f"})
}

// Java: PureFunctionIdentifierTest.testMutatesArguments_notTrackedThroughDestructuring()
func (c *PureFunctionIdentifierTest) TestMutatesArguments_notTrackedThroughDestructuring() {
	c.AssertNoPureCalls(
		"const obj = {}; // preserve newline\n" +
			"function f([x]) { x.y = 1; }\n" +
			"f({obj});\n")

	// This is technically pure, but right now we don't track that the array being destructured was
	// a literal containing literals.
	c.AssertNoPureCalls(
		"function f([x]) { x.y = 1; } // preserve newline\n" +
			"f([{}]);\n")
}

// Java: PureFunctionIdentifierTest.testMutatesArguments2()
func (c *PureFunctionIdentifierTest) TestMutatesArguments2() {
	source :=
		"function f(x) { x.y = 1; }\n" +
			"f(window);\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testMutatesArguments3()
func (c *PureFunctionIdentifierTest) TestMutatesArguments3() {
	// We could do better here with better side-effect propagation.
	source :=
		"function f(x) { x.y = 1; }\n" +
			"function g(x) { f(x); }\n" +
			"g({});\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testMutatesArguments4()
func (c *PureFunctionIdentifierTest) TestMutatesArguments4() {
	c.AssertPureCallsMarked(
		"function f(x) { x.y = 1; } // preserve newline\n"+
			"function g(x) { f({}); x.y = 1; }\n"+
			"g({});\n",
		[]string{"f", "g"})
}

// Java: PureFunctionIdentifierTest.testMutatesArguments5()
func (c *PureFunctionIdentifierTest) TestMutatesArguments5() {
	source :=
		"function f(x) {\n" +
			"  function g() {\n" +
			"    x.prop = 5;\n" +
			"  }\n" +
			"  g();\n" +
			"}\n" +
			"f(window);\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testMutatesArgumentsArray1()
func (c *PureFunctionIdentifierTest) TestMutatesArgumentsArray1() {
	source :=
		"function f(x) { arguments[0] = 1; }\n" +
			"f({});\n"
	c.AssertPureCallsMarked(source, []string{"f"})
}

// Java: PureFunctionIdentifierTest.testMutatesArgumentsArray2()
func (c *PureFunctionIdentifierTest) TestMutatesArgumentsArray2() {
	// We could be smarter here.
	source :=
		"function f(x) { arguments[0].y = 1; }\n" +
			"f({});\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testMutatesArgumentsArray3()
func (c *PureFunctionIdentifierTest) TestMutatesArgumentsArray3() {
	source :=
		"function f(x) { arguments[0].y = 1; }\n" +
			"f(x);\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testCallGenerator1()
func (c *PureFunctionIdentifierTest) TestCallGenerator1() {
	source :=
		"var x = 0;\n" +
			"function* f() {\n" +
			"  x = 2\n" +
			"  while (true) {\n" +
			"    yield x;\n" +
			"  }\n" +
			"}\n" +
			"var g = f();\n"
	c.AssertNoPureCalls(source)
	lastRoot := c.GetLastCompiler().GetRoot().GetLastChild()
	call := c.FindQualifiedNameNode("f", lastRoot).GetParent()
	AssertThat(call.IsNoSideEffectsCall()).IsFalse()
	AssertThat(call.GetSideEffectFlags()).IsEqualTo(rhino.Node_SideEffectFlags_ALL_SIDE_EFFECTS)
}

// Java: PureFunctionIdentifierTest.testCallGenerator2()
func (c *PureFunctionIdentifierTest) TestCallGenerator2() {
	source :=
		"function* f() {\n" +
			"  while (true) {\n" +
			"    yield 1;\n" +
			"  }\n" +
			"}\n" +
			"var g = f();\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testCallFunctionFOrG()
func (c *PureFunctionIdentifierTest) TestCallFunctionFOrG() {
	source :=
		"function f(){}\n" +
			"function g(){}\n" +
			"function h(){ (f || g)() }\n" +
			"h()\n"
	c.AssertPureCallsMarked(source, []string{"f || g", "h"})
}

// Java: PureFunctionIdentifierTest.testCallFunctionFOrGViaHook()
func (c *PureFunctionIdentifierTest) TestCallFunctionFOrGViaHook() {
	source :=
		"function f(){}\n" +
			"function g(){}\n" +
			"function h(){ (false ? f : g)() }\n" +
			"h()\n"
	c.AssertPureCallsMarked(source, []string{"false ? f : g", "h"})
}

// Java: PureFunctionIdentifierTest.testCallFunctionForGorH()
func (c *PureFunctionIdentifierTest) TestCallFunctionForGorH() {
	source :=
		"function f(){}\n" +
			"function g(){}\n" +
			"function h(){}\n" +
			"function i(){ (false ? f : (g || h))() }\n" +
			"i()\n"
	c.AssertPureCallsMarked(source, []string{"false ? f : (g || h)", "i"})
}

// Java: PureFunctionIdentifierTest.testCallFunctionForGWithSideEffects()
func (c *PureFunctionIdentifierTest) TestCallFunctionForGWithSideEffects() {
	source :=
		"var x = 0;\n" +
			"function f(){x = 10}\n" +
			"function g(){}\n" +
			"function h(){ (f || g)() }\n" +
			"function i(){ (g || f)() }\n" +
			"function j(){ (f || f)() }\n" +
			"function k(){ (g || g)() }\n" +
			"h(); i(); j(); k()\n"
	c.AssertPureCallsMarked(source, []string{"g || g", "k"})
}

// Java: PureFunctionIdentifierTest.testCallFunctionFOrGViaHookWithSideEffects()
func (c *PureFunctionIdentifierTest) TestCallFunctionFOrGViaHookWithSideEffects() {
	source :=
		"var x = 0;\n" +
			"function f(){x = 10}\n" +
			"function g(){}\n" +
			"function h(){ (false ? f : g)() }\n" +
			"function i(){ (false ? g : f)() }\n" +
			"function j(){ (false ? f : f)() }\n" +
			"function k(){ (false ? g : g)() }\n" +
			"h(); i(); j(); k()\n"

	c.AssertPureCallsMarked(source, []string{"false ? g : g", "k"})
}

// Java: PureFunctionIdentifierTest.testCallRegExpWithSideEffects()
func (c *PureFunctionIdentifierTest) TestCallRegExpWithSideEffects() {
	source :=
		"var x = 0;\n" +
			"function k(){(/a/).exec('')}\n" +
			"k()\n"

	c.regExpHaveSideEffects = true
	c.AssertNoPureCalls(source)
	c.regExpHaveSideEffects = false
	c.AssertPureCallsMarked(source, []string{"/a/.exec", "k"})
}

// Java: PureFunctionIdentifierTest.testAnonymousFunction1()
func (c *PureFunctionIdentifierTest) TestAnonymousFunction1() {
	c.AssertPureCallsMarked("(function (){})();", []string{"function(){}"})
}

// Java: PureFunctionIdentifierTest.testAnonymousFunction2()
func (c *PureFunctionIdentifierTest) TestAnonymousFunction2() {
	source := "(Error || function (){})();"

	c.AssertPureCallsMarked(source, []string{"(Error || function(){})"})
}

// Java: PureFunctionIdentifierTest.testAnonymousFunction3()
func (c *PureFunctionIdentifierTest) TestAnonymousFunction3() {
	source := "var a = (Error || function (){})();"

	c.AssertPureCallsMarked(source, []string{"(Error || function(){})"})
}

// Indirect complex function definitions aren't yet supported.
// Java: PureFunctionIdentifierTest.testAnonymousFunction4()
func (c *PureFunctionIdentifierTest) TestAnonymousFunction4() {
	source :=
		"var a = (Error || function (){});\n" +
			"a();\n"

	c.AssertPureCallsMarked(source, []string{"a"})
}

// Java: PureFunctionIdentifierTest.testClassMethod1()
func (c *PureFunctionIdentifierTest) TestClassMethod1() {
	source := "class C { m() { alert(1); } }; (new C).m();"
	c.AssertPureCallsMarked(source, []string{"C"})
}

// Java: PureFunctionIdentifierTest.testClassMethod2()
func (c *PureFunctionIdentifierTest) TestClassMethod2() {
	source := "class C { m() { } }; (new C).m();"
	c.AssertPureCallsMarked(source, []string{"C", "new C().m"})
}

// Java: PureFunctionIdentifierTest.testClassMethod3()
func (c *PureFunctionIdentifierTest) TestClassMethod3() {
	source := "class C { m1() { } m2() { this.m1(); }}; (new C).m2();"
	c.AssertPureCallsMarked(source, []string{"C", "this.m1", "new C().m2"})
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_implicitCtor_noSuperclass_isPure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_implicitCtor_noSuperclass_isPure() {
	c.AssertPureCallsMarked(
		"class C { }\n"+
			"\n"+
			"new C()\n",
		[]string{"C"})
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_literal_implicitCtor_noSuperclass_isPure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_literal_implicitCtor_noSuperclass_isPure() {
	c.AssertPureCallsMarked("new class C {}", []string{"class C {}"})
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_pureCtor_noSuperclass_isPure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_pureCtor_noSuperclass_isPure() {
	c.AssertPureCallsMarked(
		"class C {\n"+
			"  constructor() { }\n"+
			"}\n"+
			"\n"+
			"new C()\n",
		[]string{"C"})
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_literal_pureCtor_noSuperclass_isPure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_literal_pureCtor_noSuperclass_isPure() {
	c.AssertPureCallsMarked(
		"new class C { constructor() {} }", []string{"class C { constructor() {} }"})
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_impureCtor_noSuperclass_isImpure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_impureCtor_noSuperclass_isImpure() {
	c.AssertNoPureCalls(
		"class C {\n" +
			"  constructor() { throw 0; }\n" +
			"}\n" +
			"\n" +
			"new C()\n")
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_implictCtor_pureSuperclassCtor_isPure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_implictCtor_pureSuperclassCtor_isPure() {
	c.AssertPureCallsMarked(
		"class A { } // Implict ctor is pure.\n"+
			"class C extends A { }\n"+
			"\n"+
			"new C()\n",
		[]string{"C"})
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_literal_implictCtor_pureSuperclassCtor_isPure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_literal_implictCtor_pureSuperclassCtor_isPure() {
	c.AssertPureCallsMarked(
		"class A { } // Implicit ctor is pure.\n"+
			"new class C extends A { }\n",
		[]string{"class C extends A { }"})
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_implictCtor_impureSuperclassCtor_isImpure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_implictCtor_impureSuperclassCtor_isImpure() {
	c.AssertNoPureCalls(
		"class A {\n" +
			"  constructor() { throw 0; }\n" +
			"}\n" +
			"class C extends A { }\n" +
			"\n" +
			"new C()\n")
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_literal_implictCtor_impureSuperclassCtor_isImpure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_literal_implictCtor_impureSuperclassCtor_isImpure() {
	c.AssertNoPureCalls(
		"class A {\n" +
			"  constructor() { throw 0; }\n" +
			"}\n" +
			"new class C extends A { }\n")
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_explitCtor_pureSuperclassCtor_isPure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_explitCtor_pureSuperclassCtor_isPure() {
	c.AssertPureCallsMarked(
		"class A { }\n"+
			"class C extends A {\n"+
			"  constructor() {\n"+
			"// This should be marked as \"mutates this\" despite being a call to a pure function.\n"+
			"// `super()` binds `this`, and more importantly is syntactically required. Marking it\n"+
			"// this way prevents removal.\n"+
			"    super();\n"+
			"  }\n"+
			"}\n"+
			"\n"+
			"new C()\n",
		[]string{"C"})
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_explitCtor_impureSuperclassCtor_isImpure()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_explitCtor_impureSuperclassCtor_isImpure() {
	c.AssertNoPureCalls(
		"class A {\n" +
			"  constructor() { throw 0; }\n" +
			"}\n" +
			"class C extends A {\n" +
			"  constructor() { super(); }\n" +
			"}\n" +
			"\n" +
			"new C()\n")
}

// Java: PureFunctionIdentifierTest.testClassInstantiation_mutatesThis_traversesSuper()
func (c *PureFunctionIdentifierTest) TestClassInstantiation_mutatesThis_traversesSuper() {
	c.AssertPureCallsMarked(
		"class A {\n"+
			"  constructor() { this.foo = 9; }\n"+
			"}\n"+
			"class C extends A {\n"+
			"  constructor() { super(); }\n"+
			"}\n"+
			"\n"+
			"new C()\n",
		[]string{"C"})
}

// Java: PureFunctionIdentifierTest.testMutatesThis_traversesThisAsReceiver()
func (c *PureFunctionIdentifierTest) TestMutatesThis_traversesThisAsReceiver() {
	c.AssertPureCallsMarked(
		"class C {\n"+
			"  constructor() {\n"+
			"    /** @type {number} */\n"+
			"    this.foo;\n"+
			"\n"+
			"    this.mutateThis();\n"+
			"  }\n"+
			"\n"+
			"  mutateThis() { this.foo = 9; }\n"+
			"}\n"+
			"\n"+
			"// Instantiation swallows \"mutates this\" side-effects.\n"+
			"var x = new C();\n"+
			"// But in general it is still a side-effect.\n"+
			"x.mutateThis()\n",
		[]string{"C"})
}

// Java: PureFunctionIdentifierTest.testMutatesThis_traversesSuperAsReceiver()
func (c *PureFunctionIdentifierTest) TestMutatesThis_traversesSuperAsReceiver() {
	c.AssertPureCallsMarked(
		"class C {\n"+
			"  constructor() {\n"+
			"    /** @type {number} */\n"+
			"    this.foo;\n"+
			"\n"+
			"// Ignore the fact that there's no superclass.\n"+
			"    super.mutateThis();\n"+
			"  }\n"+
			"\n"+
			"  mutateThis() { this.foo = 9; }\n"+
			"}\n"+
			"\n"+
			"// Instantiation swallows \"mutates this\" side-effects.\n"+
			"var x = new C();\n"+
			"// But in general it is still a side-effect.\n"+
			"x.mutateThis()\n",
		[]string{"C"})
}

// Java: PureFunctionIdentifierTest.testMutatesThis_expandsToMutatesGlobalScope_traversingArbitraryReceiver()
func (c *PureFunctionIdentifierTest) TestMutatesThis_expandsToMutatesGlobalScope_traversingArbitraryReceiver() {
	c.AssertNoPureCalls(
		"class C {\n" +
			"  constructor() {\n" +
			"    /** @type {number} */\n" +
			"    this.foo;\n" +
			"\n" +
			"    somethingElse.mutateThis();\n" +
			"  }\n" +
			"\n" +
			"  mutateThis() {\n" +
			"    this.foo = 9;\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"// Instantiation swallows \"mutates this\" side-effects.\n" +
			"var x = new C();\n" +
			"// But in general it is still a side-effect.\n" +
			"x.mutateThis()\n")
}

// Java: PureFunctionIdentifierTest.testGlobalScopeTaintedByWayOfThisPropertyAndForOfLoop()
func (c *PureFunctionIdentifierTest) TestGlobalScopeTaintedByWayOfThisPropertyAndForOfLoop() {
	source :=
		"class C {\n" +
			"  constructor(elements) {\n" +
			"    this.elements = elements;\n" +
			"    this.m1();\n" +
			"  }\n" +
			"  m1() {\n" +
			"    for (const element of this.elements) {\n" +
			"      element.someProp = 1;\n" +
			"    }\n" +
			"  }\n" +
			"}\n" +
			"new C([]).m1()\n"
	c.AssertPureCallsMarked(source, []string{})
}

// Java: PureFunctionIdentifierTest.testGlobalScopeTaintedByWayOfThisPropertyAndForOfLoopWithDestructuring()
func (c *PureFunctionIdentifierTest) TestGlobalScopeTaintedByWayOfThisPropertyAndForOfLoopWithDestructuring() {
	source :=
		"class C {\n" +
			"  constructor(elements) {\n" +
			"    this.elements = elements;\n" +
			"    this.m1();\n" +
			"  }\n" +
			"  m1() {\n" +
			"    this.elements[0].someProp = 1;\n" +
			"    for (const {prop} of this.elements) {\n" +
			"      prop.someProp = 1;\n" +
			"    }\n" +
			"  }\n" +
			"}\n" +
			"new C([{prop: {}}]).m1()\n"
	c.AssertPureCallsMarked(source, []string{})
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfFunctionScopedLet()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfFunctionScopedLet() {
	source :=
		"function m1(elements) {\n" +
			"  let e;\n" +
			"  for (let i = 0; i < elements.length; ++i) {\n" +
			"    e = elements[i];\n" +
			"    e.someProp = 1;\n" +
			"  }\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertPureCallsMarked(source, []string{})
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfFunctionScopedLetAssignedInForOf()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfFunctionScopedLetAssignedInForOf() {
	source :=
		"function m1(elements) {\n" +
			"  let e;\n" +
			"  for (e of elements) {\n" +
			"    e.someProp = 1;\n" +
			"  }\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertPureCallsMarked(source, []string{})
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfFunctionScopedLetAssignedInForOfWithDestructuring()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfFunctionScopedLetAssignedInForOfWithDestructuring() {
	source :=
		"function m1(elements) {\n" +
			"  let e;\n" +
			"  for ({e} of elements) {\n" +
			"    e.someProp = 1;\n" +
			"  }\n" +
			"}\n" +
			"m1([{e: {}}]);\n"
	c.AssertPureCallsMarked(source, []string{})
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfForOfAssignmentToQualifiedName()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfForOfAssignmentToQualifiedName() {
	source :=
		"function m1(obj, elements) {\n" +
			"  for (obj.e of elements) {\n" +
			"    if (obj.e != null) break;\n" +
			"  }\n" +
			"}\n" +
			"var globalObj = {};\n" +
			"m1(globalObj, []);\n"
	c.AssertPureCallsMarked(source, []string{})
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfForInAssignmentToQualifiedName()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfForInAssignmentToQualifiedName() {
	source :=
		"/**\n" +
			" *@param {!Object} obj\n" +
			" *@param {!Object} propObj\n" +
			" *@returns {number}\n" +
			" */\n" +
			"function m1(/** ? */ obj, /** Object */ propObj) {\n" +
			"  var propCharCount = 0;\n" +
			"  obj.e = '';\n" +
			"  for (obj.e in propObj) {\n" +
			"    propCharCount += obj.e.length;\n" +
			"  }\n" +
			"  return propCharCount;\n" +
			"}\n" +
			"var globalObj = {};\n" +
			"m1(globalObj, {x: 1});\n"
	c.AssertPureCallsMarked(source, []string{})
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfBlockScopedLet()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfBlockScopedLet() {
	source :=
		"function m1(elements) {\n" +
			"  for (let i = 0; i < elements.length; ++i) {\n" +
			"    let e;\n" +
			"    e = elements[i];\n" +
			"    e.someProp = 1;\n" +
			"  }\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertPureCallsMarked(source, []string{})
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfBlockScopedConst()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfBlockScopedConst() {
	source :=
		"function m1(elements) {\n" +
			"  for (let i = 0; i < elements.length; ++i) {\n" +
			"    const e = elements[i];\n" +
			"    e.someProp = 1;\n" +
			"  }\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertPureCallsMarked(source, []string{})
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfForOfScopedConst()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfForOfScopedConst() {
	source :=
		"function m1(elements) {\n" +
			"  for (const e of elements) {\n" +
			"    e.someProp = 1;\n" +
			"  }\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfForOfScopedConstWithDestructuring()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfForOfScopedConstWithDestructuring() {
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	source :=
		"function m1(elements) {\n" +
			"  for (const {e} of elements) {\n" +
			"    e.someProp = 1;\n" +
			"  }\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfNameDeclaration()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfNameDeclaration() {
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	source :=
		"function m1(obj) {\n" +
			"  let p = obj.p;\n" +
			"  p.someProp = 1;\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertPureCallsMarked(source, []string{})
}

// Java: PureFunctionIdentifierTest.testUnusedDestructuringNameDeclarationIsPure()
func (c *PureFunctionIdentifierTest) TestUnusedDestructuringNameDeclarationIsPure() {
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	source :=
		"function m1(obj) {\n" +
			"  let {p} = obj;\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertPureCallsMarked(source, []string{"m1"})
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfDestructuringNameDeclaration()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfDestructuringNameDeclaration() {
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	source :=
		"function m1(obj) {\n" +
			"  let {p} = obj;\n" +
			"  p.someProp = 1;\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfNestedDestructuringNameDeclaration()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfNestedDestructuringNameDeclaration() {
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	source :=
		"function m1(obj) {\n" +
			"  let {q: {p}} = obj;\n" +
			"  p.someProp = 1;\n" +
			"}\n" +
			"m1({});\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfDestructuringAssignment()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfDestructuringAssignment() {
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	source :=
		"function m1(obj) {\n" +
			"  let p;\n" +
			"  ({p} = obj);\n" +
			"  p.someProp = 1;\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testArgumentTaintedByWayOfArrayDestructuringAssignment()
func (c *PureFunctionIdentifierTest) TestArgumentTaintedByWayOfArrayDestructuringAssignment() {
	source :=
		"function m1(obj) {\n" +
			"  let p;\n" +
			"  ([p] = obj);\n" +
			"  p.someProp = 1;\n" +
			"}\n" +
			"m1([]);\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testDestructuringAssignMutatingGlobalState()
func (c *PureFunctionIdentifierTest) TestDestructuringAssignMutatingGlobalState() {
	source :=
		"var a = 0;\n" +
			"function m1() {\n" +
			"  ([a] = []);\n" +
			"}\n" +
			"m1();\n"
	c.AssertNoPureCalls(source)

	source =
		"var a = 0;\n" +
			"function m1() {\n" +
			"  ({a} = {a: 0});\n" +
			"}\n" +
			"m1();\n"
	c.AssertNoPureCalls(source)
}

// Java: PureFunctionIdentifierTest.testDefaultValueInitializers_areConsidered_whenAnalyzingDestructuring()
func (c *PureFunctionIdentifierTest) TestDefaultValueInitializers_areConsidered_whenAnalyzingDestructuring() {
	c.AssertNoPureCalls(
		"function impure() { throw 'something'; }\n" +
			"\n" +
			"function foo() { const {a = impure()} = {a: undefined}; }\n" +
			"foo();\n")

	c.AssertPureCallsMarked(
		"function foo() { const {a = 0} = {a: undefined}; }\n"+
			"foo();\n",
		[]string{"foo"})

	c.AssertNoPureCalls(
		"function impure() { throw 'something'; }\n" +
			"\n" +
			"function foo() { const [{a = impure()}] = [{a: undefined}]; }\n" +
			"foo();\n")
}

// Java: PureFunctionIdentifierTest.testDefaultValueInitializers_areConsidered_whenAnalyzingFunctions()
func (c *PureFunctionIdentifierTest) TestDefaultValueInitializers_areConsidered_whenAnalyzingFunctions() {
	c.AssertNoPureCalls(
		"function impure() { throw 'something'; }\n" +
			"\n" +
			"function foo(a = impure()) { }\n" +
			"foo();\n")

	c.AssertPureCallsMarked(
		"function pure() { }\n"+
			"\n"+
			"function foo(a = pure()) { }\n"+
			"foo();\n",
		[]string{"pure", "foo"})

	c.AssertPureCallsMarked(
		"function foo(a = 0) { }\n"+
			"foo();\n",
		[]string{"foo"})
}

// Java: PureFunctionIdentifierTest.testFunctionProperties1()
func (c *PureFunctionIdentifierTest) TestFunctionProperties1() {
	source :=
		"/** @constructor */\n" +
			"function F() { this.bar; }\n" +
			"function g() {\n" +
			"  this.bar = function() { alert(3); };\n" +
			"}\n" +
			"var x = new F();\n" +
			"g.call(x);\n" +
			"x.bar();\n"
	c.AssertPureCallsMarkedWithPostcondition(
		source,
		[]string{"F"},
		CompilerTestCase_PostconditionFunc(func(compiler *Compiler) {
			lastRoot := compiler.GetRoot()
			call := c.FindQualifiedNameNode("g.call", lastRoot).GetParent()
			AssertThat(call.GetSideEffectFlags()).
				IsEqualTo(
					rhino.NewNode_SideEffectFlags().ClearAllFlags().SetMutatesArguments().ValueOf())
		}))
}

// Java: PureFunctionIdentifierTest.testAwait_makesFunctionImpure()
func (c *PureFunctionIdentifierTest) TestAwait_makesFunctionImpure() {
	c.AssertNoPureCalls(
		"async function foo() { await 0; }\n" +
			"foo();\n")
}

// Java: PureFunctionIdentifierTest.testTaggedTemplatelit_propagatesCalleeSideEffects()
func (c *PureFunctionIdentifierTest) TestTaggedTemplatelit_propagatesCalleeSideEffects() {
	c.AssertPureCallsMarkedWithPostcondition(
		"function tag(a) { throw 'something'; }\n"+
			"\n"+
			"function foo() { tag`Hello World!`; }\n"+
			"foo();\n",
		[]string{},
		CompilerTestCase_PostconditionFunc(func(compiler *Compiler) {
			lastRoot := compiler.GetRoot()
			tagDef := c.FindQualifiedNameNode("tag", lastRoot).GetParent()
			fooDef := c.FindQualifiedNameNode("foo", lastRoot).GetParent()
			AssertThat(fooDef.GetSideEffectFlags()).IsEqualTo(tagDef.GetSideEffectFlags())
		}))
}

// Java: PureFunctionIdentifierTest.testIterableIteration_forOf_makesFunctionImpure()
func (c *PureFunctionIdentifierTest) TestIterableIteration_forOf_makesFunctionImpure() {
	// TODO(b/127862986): We could consider iteration over a known Array value to be pure.
	c.AssertNoPureCalls(
		"function foo(a) { for (const t of a) { } }\n" +
			"foo(x);\n")
}

// Java: PureFunctionIdentifierTest.testIterableIteration_forAwaitOf_makesFunctionImpure()
func (c *PureFunctionIdentifierTest) TestIterableIteration_forAwaitOf_makesFunctionImpure() {
	// TODO(b/127862986): We could consider iteration over a known Array value to be pure.
	c.AssertNoPureCalls(
		"async function foo(a) { for await (const t of a) { } }\n" +
			"foo(x);\n")
}

// Java: PureFunctionIdentifierTest.testIterableIteration_arraySpread_makesFunctionImpure()
func (c *PureFunctionIdentifierTest) TestIterableIteration_arraySpread_makesFunctionImpure() {
	// TODO(b/127862986): We could consider iteration over a known Array value to be pure.
	c.AssertNoPureCalls(
		"function foo(a) { [...a]; }\n" +
			"foo(x);\n")
}

// Java: PureFunctionIdentifierTest.testIterableIteration_callSpread_makesFunctionImpure()
func (c *PureFunctionIdentifierTest) TestIterableIteration_callSpread_makesFunctionImpure() {
	// TODO(b/127862986): We could consider iteration over a known Array value to be pure.
	c.AssertPureCallsMarked(
		"function pure(...rest) { }\n"+
			"\n"+
			"function foo(a) { pure(...a); }\n"+
			"foo(x);\n",
		[]string{"pure"})
}

// Java: PureFunctionIdentifierTest.testIterableIteration_newSpread_makesFunctionImpure()
func (c *PureFunctionIdentifierTest) TestIterableIteration_newSpread_makesFunctionImpure() {
	// TODO(b/127862986): We could consider iteration over a known Array value to be pure.
	c.AssertPureCallsMarked(
		"/** @constructor */\n"+
			"function pure(...rest) { }\n"+
			"\n"+
			"function foo(a) { new pure(...a); }\n"+
			"foo(x);\n",
		[]string{"pure"})
}

// Java: PureFunctionIdentifierTest.testIterableIteration_paramListRest_leavesFunctionPure()
func (c *PureFunctionIdentifierTest) TestIterableIteration_paramListRest_leavesFunctionPure() {
	c.AssertPureCallsMarked(
		"function foo(...a) { }\n"+
			"foo();\n",
		[]string{"foo"})
}

// Java: PureFunctionIdentifierTest.testIterableIteration_arrayPatternRest_makesFunctionImpure()
func (c *PureFunctionIdentifierTest) TestIterableIteration_arrayPatternRest_makesFunctionImpure() {
	// TODO(b/127862986): We could consider iteration over a known Array value to be pure.
	c.AssertNoPureCalls(
		"function foo(a) { const [...x] = a; }\n" +
			"foo(v);\n")
}

// Java: PureFunctionIdentifierTest.testIterableIteration_yieldStar_makesFunctionImpure()
func (c *PureFunctionIdentifierTest) TestIterableIteration_yieldStar_makesFunctionImpure() {
	// TODO(b/127862986): We could consider iteration over a known Array value to be pure.
	c.AssertNoPureCalls(
		"function* foo(a) { yield* a; }\n" +
			"foo(x);\n")
}

// Java: PureFunctionIdentifierTest.testObjectSpread_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestObjectSpread_hasSideEffects() {
	// Object-spread may trigger a getter.
	c.AssertNoPureCalls(
		"function foo(a) { ({...a}); }\n" +
			"foo(x);\n")
}

// Java: PureFunctionIdentifierTest.testObjectRest_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestObjectRest_hasSideEffects() {
	// Object-rest may trigger a getter.
	c.AssertNoPureCalls(
		"function foo(a, b) { ({...b} = a); }\n" +
			"foo(x, y);\n")
	c.AssertNoPureCalls(
		"function foo({...b}) { }\n" +
			"foo(x);\n")
}

// Java: PureFunctionIdentifierTest.testGetterAccess_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestGetterAccess_hasSideEffects() {
	c.AssertNoPureCalls(
		"class Foo { get getter() { } }\n" +
			"\n" +
			"function foo(a) { a.getter; }\n" +
			"foo(x);\n")
	c.AssertNoPureCalls(
		"class Foo { get getter() { } }\n" +
			"\n" +
			"function foo(a) { const {getter} = a; }\n" +
			"foo(x);\n")
}

// Java: PureFunctionIdentifierTest.testSetterAccess_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestSetterAccess_hasSideEffects() {
	c.AssertNoPureCalls(
		"class Foo { set setter(x) { } }\n" +
			"\n" +
			"function foo(a) { a.setter = 0; }\n" +
			"foo(x);\n")
}

// Java: PureFunctionIdentifierTest.testForAwaitOf_makesFunctionImpure()
func (c *PureFunctionIdentifierTest) TestForAwaitOf_makesFunctionImpure() {
	c.AssertNoPureCalls(
		"async function foo() { for await (const t of []) { } }\n" +
			"foo();\n")
}

// Java: PureFunctionIdentifierTest.testExistenceOfAGetter_makesPropertyNameImpure()
func (c *PureFunctionIdentifierTest) TestExistenceOfAGetter_makesPropertyNameImpure() {
	c.DeclareAccessor("foo", AccessorSummary_PropertyAccessKind_GETTER_ONLY)

	// Imagine the getter returned a function.
	c.AssertNoPureCalls(
		"class Foo {\n" +
			"  foo() {}\n" +
			"}\n" +
			"x.foo();\n")
}

// Java: PureFunctionIdentifierTest.testExistenceOfASetter_makesPropertyNameImpure()
func (c *PureFunctionIdentifierTest) TestExistenceOfASetter_makesPropertyNameImpure() {
	c.DeclareAccessor("foo", AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	// This doesn't actually seem like a risk but it's hard to say, and other optimizations that use
	// optimize calls would be dangerous on setters.
	c.AssertNoPureCalls(
		"class Foo {\n" +
			"  foo() {}\n" +
			"}\n" +
			"x.foo();\n")
}

// Java: PureFunctionIdentifierTest.testCallCache()
func (c *PureFunctionIdentifierTest) TestCallCache() {
	source :=
		"var valueFn = function() {};\n" +
			"goog.reflect.cache(externObj, \"foo\", valueFn)\n"
	c.AssertPureCallsMarkedWithPostcondition(
		source,
		[]string{"goog.reflect.cache"},
		CompilerTestCase_PostconditionFunc(func(compiler *Compiler) {
			lastRoot := compiler.GetRoot().GetLastChild()
			call := c.FindQualifiedNameNode("goog.reflect.cache", lastRoot).GetParent()
			AssertThat(call.IsNoSideEffectsCall()).IsTrue()
			AssertThat(call.MayMutateGlobalStateOrThrow()).IsFalse()
		}))
}

// Java: PureFunctionIdentifierTest.testCallCache_withKeyFn()
func (c *PureFunctionIdentifierTest) TestCallCache_withKeyFn() {
	source :=
		"var valueFn = function(v) { return v };\n" +
			"var keyFn = function(v) { return v };\n" +
			"goog.reflect.cache(externObj, \"foo\", valueFn, keyFn)\n"
	c.AssertPureCallsMarkedWithPostcondition(
		source,
		[]string{"goog.reflect.cache"},
		CompilerTestCase_PostconditionFunc(func(compiler *Compiler) {
			lastRoot := compiler.GetRoot().GetLastChild()
			call := c.FindQualifiedNameNode("goog.reflect.cache", lastRoot).GetParent()
			AssertThat(call.IsNoSideEffectsCall()).IsTrue()
			AssertThat(call.MayMutateGlobalStateOrThrow()).IsFalse()
		}))
}

// Java: PureFunctionIdentifierTest.testCallCache_anonymousFn()
func (c *PureFunctionIdentifierTest) TestCallCache_anonymousFn() {
	source := "goog.reflect.cache(externObj, \"foo\", function(v) { return v })"
	c.AssertPureCallsMarkedWithPostcondition(
		source,
		[]string{"goog.reflect.cache"},
		CompilerTestCase_PostconditionFunc(func(compiler *Compiler) {
			lastRoot := compiler.GetRoot().GetLastChild()
			call := c.FindQualifiedNameNode("goog.reflect.cache", lastRoot).GetParent()
			AssertThat(call.IsNoSideEffectsCall()).IsTrue()
			AssertThat(call.MayMutateGlobalStateOrThrow()).IsFalse()
		}))
}

// Java: PureFunctionIdentifierTest.testCallCache_anonymousFn_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestCallCache_anonymousFn_hasSideEffects() {
	source :=
		"var x = 0;\n" +
			"goog.reflect.cache(externObj, \"foo\", function(v) { return (x+=1) })\n"
	c.AssertNoPureCallsWithPostcondition(
		source,
		CompilerTestCase_PostconditionFunc(func(compiler *Compiler) {
			lastRoot := compiler.GetRoot().GetLastChild()
			call := c.FindQualifiedNameNode("goog.reflect.cache", lastRoot).GetParent()
			AssertThat(call.IsNoSideEffectsCall()).IsFalse()
			AssertThat(call.MayMutateGlobalStateOrThrow()).IsTrue()
		}))
}

// Java: PureFunctionIdentifierTest.testCallCache_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestCallCache_hasSideEffects() {
	source :=
		"var x = 0;\n" +
			"var valueFn = function() { return (x+=1); };\n" +
			"goog.reflect.cache(externObj, \"foo\", valueFn)\n"
	c.AssertNoPureCallsWithPostcondition(
		source,
		CompilerTestCase_PostconditionFunc(func(compiler *Compiler) {
			lastRoot := compiler.GetRoot().GetLastChild()
			call := c.FindQualifiedNameNode("goog.reflect.cache", lastRoot).GetParent()
			AssertThat(call.IsNoSideEffectsCall()).IsFalse()
			AssertThat(call.MayMutateGlobalStateOrThrow()).IsTrue()
		}))
}

// Java: PureFunctionIdentifierTest.testCallCache_withKeyFn_hasSideEffects()
func (c *PureFunctionIdentifierTest) TestCallCache_withKeyFn_hasSideEffects() {
	source :=
		"var x = 0;\n" +
			"var keyFn = function(v) { return (x+=1) };\n" +
			"var valueFn = function(v) { return v };\n" +
			"goog.reflect.cache(externObj, \"foo\", valueFn, keyFn)\n"
	c.AssertNoPureCallsWithPostcondition(
		source,
		CompilerTestCase_PostconditionFunc(func(compiler *Compiler) {
			lastRoot := compiler.GetRoot().GetLastChild()
			call := c.FindQualifiedNameNode("goog.reflect.cache", lastRoot).GetParent()
			AssertThat(call.IsNoSideEffectsCall()).IsFalse()
			AssertThat(call.MayMutateGlobalStateOrThrow()).IsTrue()
		}))
}

// Java: PureFunctionIdentifierTest.testCallCache_propagatesSideEffects()
func (c *PureFunctionIdentifierTest) TestCallCache_propagatesSideEffects() {
	source :=
		"var valueFn = function(x) { return x * 2; };\n" +
			"var helper = function(x) { return goog.reflect.cache(externObj, x, valueFn); };\n" +
			"helper(10);\n"
	c.AssertPureCallsMarkedWithPostcondition(
		source,
		[]string{"goog.reflect.cache", "helper"},
		CompilerTestCase_PostconditionFunc(func(compiler *Compiler) {
			lastRoot := compiler.GetRoot().GetLastChild()
			cacheCall := c.FindQualifiedNameNode("goog.reflect.cache", lastRoot).GetParent()
			AssertThat(cacheCall.IsNoSideEffectsCall()).IsTrue()
			AssertThat(cacheCall.MayMutateGlobalStateOrThrow()).IsFalse()

			helperNodes := c.FindQualifiedNameNodes("helper", lastRoot)
			helperCall := helperNodes[len(helperNodes)-1].GetParent() // Iterables.getLast
			AssertThat(helperCall.IsNoSideEffectsCall()).IsTrue()
			AssertThat(helperCall.MayMutateGlobalStateOrThrow()).IsFalse()
		}))
}

// Java: PureFunctionIdentifierTest.testDynamicImport()
func (c *PureFunctionIdentifierTest) TestDynamicImport() {
	c.IgnoreWarnings(DiagnosticGroups_MODULE_LOAD)
	c.AssertNoPureCalls("import('./module.js')")
	c.AssertNoPureCalls("function doImport() { import('./module.js'); } doImport();")
}

// Java: PureFunctionIdentifierTest.testParenthesizedExpression()
func (c *PureFunctionIdentifierTest) TestParenthesizedExpression() {
	c.AssertPureCallsMarked(
		"const namespace = {};\n"+
			"namespace.noSideEffects = function(x) { return 1; };\n"+
			"// NOTE: `IRFactory` will unwrap `(0, callee)(42)`, so we need to use a non-number\n"+
			"// to preserve the parentheses for the test.\n"+
			"('', namespace.noSideEffects)(42);\n",
		[]string{"('', namespace.noSideEffects)"})
}

// Java: PureFunctionIdentifierTest.assertCallableExpressionPure(boolean, String)
func (c *PureFunctionIdentifierTest) AssertCallableExpressionPure(purity bool, expression string) {
	expression = "(" + expression + ")"
	directInvocation := expression + "()"
	aliasInvocation := "var $testFn = " + expression + "; $testFn()"

	if purity {
		c.AssertPureCallsMarked(directInvocation, []string{expression})
		c.AssertPureCallsMarked(aliasInvocation, []string{"$testFn"})
	} else {
		c.AssertNoPureCalls(directInvocation)
		c.AssertNoPureCalls(aliasInvocation)
	}
}

// Java: PureFunctionIdentifierTest.assertNoPureCalls(String)
func (c *PureFunctionIdentifierTest) AssertNoPureCalls(source string) {
	c.AssertPureCallsMarkedWithPostcondition(source, []string{}, nil)
}

// Java: PureFunctionIdentifierTest.assertNoPureCalls(String, Postcondition)
func (c *PureFunctionIdentifierTest) AssertNoPureCallsWithPostcondition(source string, post CompilerTestCase_Postcondition) {
	c.AssertPureCallsMarkedWithPostcondition(source, []string{}, post)
}

// Java: PureFunctionIdentifierTest.assertPureCallsMarked(String, List<String>)
func (c *PureFunctionIdentifierTest) AssertPureCallsMarked(source string, expected []string) {
	c.AssertPureCallsMarkedWithPostcondition(source, expected, nil)
}

// Java: PureFunctionIdentifierTest.assertPureCallsMarked(String, List<String>, @Nullable Postcondition)
func (c *PureFunctionIdentifierTest) AssertPureCallsMarkedWithPostcondition(source string, expected []string, post CompilerTestCase_Postcondition) {
	c.TestSame(
		c.Srcs(source),
		c.Postcondition(
			func(compiler *Compiler) {
				AssertThat(c.noSideEffectCalls).
					ComparingElementsUsing(JSCompCorrespondences_EQUALITY_WHEN_PARSED_AS_EXPRESSION).
					ContainsExactlyElementsIn(expected)
				if post != nil {
					post.Verify(compiler)
				}
			}))
}

// Java: PureFunctionIdentifierTest.getProcessor(Compiler)
func (c *PureFunctionIdentifierTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPureFunctionIdentifierTest_NoSideEffectCallEnumerator(c, compiler)
}

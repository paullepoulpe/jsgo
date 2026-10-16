// Port of com/google/javascript/jscomp/DevirtualizeMethodsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/DevirtualizeMethodsTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"strings"
	"testing"
)

type DevirtualizeMethodsTest struct {
	CompilerTestCaseBase
}

// Java: DevirtualizeMethodsTest()
func NewDevirtualizeMethodsTest() *DevirtualizeMethodsTest {
	c := &DevirtualizeMethodsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, DevirtualizeMethodsTest_EXTERNAL_SYMBOLS)
	return c
}

func TestDevirtualizeMethodsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewDevirtualizeMethodsTest() })
}

var DevirtualizeMethodsTest_EXTERNAL_SYMBOLS = CompilerTestCase_DEFAULT_EXTERNS + "var extern;extern.externalMethod"

// JUnitBefore runs the @Before methods: CompilerTestCase.setUp() (the
// superclass one) first, then customSetUp().
func (c *DevirtualizeMethodsTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: DevirtualizeMethodsTest.customSetUp()
func (c *DevirtualizeMethodsTest) CustomSetUp() {
	c.EnableNormalize() // Required for `OptimizeCalls`.
	c.DisableTypeCheck()
}

/** Combine source strings using ';' as the separator. */
// Java: DevirtualizeMethodsTest.semicolonJoin(String...)
func DevirtualizeMethodsTest_SemicolonJoin(parts ...string) string {
	return strings.Join(parts, ";")
}

// Java: DevirtualizeMethodsTest.testRewritePrototypeMethodsWithCorrectColors()
func (c *DevirtualizeMethodsTest) TestRewritePrototypeMethodsWithCorrectColors() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	input :=
		"/** @constructor */\n" +
			"function A() { this.x = 3; }\n" +
			"A: A;\n" +
			"/** @return {number} */\n" +
			"A.prototype.foo = function() { return this.x; };\n" +
			"/** @param {number} p\n" +
			"    @return {number} */\n" +
			"A.prototype.bar = function(p) { return this.x; };\n" +
			"A.prototype.baz = function() {};\n" +
			"var o = new A();\n" +
			"FOO_RESULT: o.foo();\n" +
			"BAR_RESULT: o.bar(2);\n" +
			"BAZ_RESULT: o.baz()\n"
	expected :=
		"/** @constructor */\n" +
			"function A(){ this.x = 3; }\n" +
			"A: A;\n" +
			"var JSCompiler_StaticMethods_foo =\n" +
			"  function(JSCompiler_StaticMethods_foo$self) {\n" +
			"    return JSCompiler_StaticMethods_foo$self.x\n" +
			"  };\n" +
			"var JSCompiler_StaticMethods_bar =\n" +
			"  function(JSCompiler_StaticMethods_bar$self, p) {\n" +
			"    return JSCompiler_StaticMethods_bar$self.x\n" +
			"  };\n" +
			"var JSCompiler_StaticMethods_baz =\n" +
			"  function(JSCompiler_StaticMethods_baz$self) {};\n" +
			"var o = new A();\n" +
			"FOO_RESULT: JSCompiler_StaticMethods_foo(o);\n" +
			"BAR_RESULT: JSCompiler_StaticMethods_bar(o, 2);\n" +
			"BAZ_RESULT: JSCompiler_StaticMethods_baz(o)\n"

	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()

	c.Test(input, expected)
	c.CheckColorOfRewrittenMethods()
}

// Java: DevirtualizeMethodsTest.checkColorOfRewrittenMethods()
func (c *DevirtualizeMethodsTest) CheckColorOfRewrittenMethods() {
	fooType := c.GetLabelledExpression("FOO_RESULT").GetFirstChild().GetColor()
	barType := c.GetLabelledExpression("BAR_RESULT").GetFirstChild().GetColor()
	bazType := c.GetLabelledExpression("BAZ_RESULT").GetFirstChild().GetColor()
	fooResultType := c.GetLabelledExpression("FOO_RESULT").GetColor()
	barResultType := c.GetLabelledExpression("BAR_RESULT").GetColor()
	bazResultType := c.GetLabelledExpression("BAZ_RESULT").GetColor()

	AssertThat(fooResultType).IsEqualTo(colors.StandardColors_NUMBER)
	AssertThat(barResultType).IsEqualTo(colors.StandardColors_NUMBER)
	AssertThat(bazResultType).IsEqualTo(colors.StandardColors_NULL_OR_VOID)

	AssertThat(fooType).IsEqualTo(colors.StandardColors_TOP_OBJECT)
	AssertThat(barType).IsEqualTo(colors.StandardColors_TOP_OBJECT)
	AssertThat(bazType).IsEqualTo(colors.StandardColors_TOP_OBJECT)
}

// Java: DevirtualizeMethodsTest.getLabelledExpression(String)
func (c *DevirtualizeMethodsTest) GetLabelledExpression(label string) *rhino.Node {
	root := c.GetLastCompiler().GetJsRoot()

	return jutil.CheckNotNullMsg1(
		DevirtualizeMethodsTest_GetLabelledExpressionIfPresent(label, root),
		"Could not find statement matching label %s",
		label)
}

// Java: DevirtualizeMethodsTest.getLabelledExpressionIfPresent(String, Node)
func DevirtualizeMethodsTest_GetLabelledExpressionIfPresent(label string, root *rhino.Node) *rhino.Node {
	if root.IsLabel() && root.GetFirstChild().GetString() == label {
		labelledBlock := root.GetSecondChild()
		jutil.CheckStateMsg1(labelledBlock.IsBlock(), "%s", labelledBlock)
		jutil.CheckStateMsg1(
			labelledBlock.HasOneChild() && labelledBlock.GetOnlyChild().IsExprResult(),
			"Unexpected children of BLOCK %s",
			labelledBlock)
		return labelledBlock.GetOnlyChild().GetOnlyChild()
	}
	for child := root.GetFirstChild(); child != nil; child = child.GetNext() {
		possibleMatch := DevirtualizeMethodsTest_GetLabelledExpressionIfPresent(label, child)
		if possibleMatch != nil {
			return possibleMatch
		}
	}
	return nil
}

// Java: DevirtualizeMethodsTest.testRewriteChained()
func (c *DevirtualizeMethodsTest) TestRewriteChained() {
	source :=
		"A.prototype.foo = function(){return this.b};\n" +
			"B.prototype.bar = function(){};\n" +
			"o.foo().bar()\n"

	expected :=
		"var JSCompiler_StaticMethods_foo =\n" +
			"function(JSCompiler_StaticMethods_foo$self) {\n" +
			"  return JSCompiler_StaticMethods_foo$self.b\n" +
			"};\n" +
			"var JSCompiler_StaticMethods_bar =\n" +
			"function(JSCompiler_StaticMethods_bar$self) {\n" +
			"};\n" +
			"JSCompiler_StaticMethods_bar(JSCompiler_StaticMethods_foo(o))\n"
	c.Test(source, expected)
}

/** Inputs for declaration used as an r-value tests. */
const DevirtualizeMethodsTest_NoRewriteDeclarationUsedAsRValue_DECL = "a.prototype.foo = function() {}"
const DevirtualizeMethodsTest_NoRewriteDeclarationUsedAsRValue_CALL = "o.foo()"

// Java: DevirtualizeMethodsTest.testRewriteDeclIsExpressionStatement()
func (c *DevirtualizeMethodsTest) TestRewriteDeclIsExpressionStatement() {
	c.Test(
		DevirtualizeMethodsTest_SemicolonJoin(DevirtualizeMethodsTest_NoRewriteDeclarationUsedAsRValue_DECL, DevirtualizeMethodsTest_NoRewriteDeclarationUsedAsRValue_CALL),
		"var JSCompiler_StaticMethods_foo =\n"+
			"function(JSCompiler_StaticMethods_foo$self) {};\n"+
			"JSCompiler_StaticMethods_foo(o)\n")
}

// Java: DevirtualizeMethodsTest.testNoRewriteDeclUsedAsAssignmentRhs()
func (c *DevirtualizeMethodsTest) TestNoRewriteDeclUsedAsAssignmentRhs() {
	c.TestSame(
		DevirtualizeMethodsTest_SemicolonJoin(
			"var c = "+DevirtualizeMethodsTest_NoRewriteDeclarationUsedAsRValue_DECL,
			DevirtualizeMethodsTest_NoRewriteDeclarationUsedAsRValue_CALL))
}

// Java: DevirtualizeMethodsTest.testNoRewriteDeclUsedAsCallArgument()
func (c *DevirtualizeMethodsTest) TestNoRewriteDeclUsedAsCallArgument() {
	c.TestSame(
		DevirtualizeMethodsTest_SemicolonJoin(
			"f("+DevirtualizeMethodsTest_NoRewriteDeclarationUsedAsRValue_DECL+")",
			DevirtualizeMethodsTest_NoRewriteDeclarationUsedAsRValue_CALL))
}

// Java: DevirtualizeMethodsTest.testRewrite_ifDefined_unconditionally()
func (c *DevirtualizeMethodsTest) TestRewrite_ifDefined_unconditionally() {
	c.Test(
		"function a(){}\n"+
			"a.prototype.foo = function() {return this.x};\n"+
			"var o = new a;\n"+
			"o.foo()\n",
		"function a(){}\n"+
			"var JSCompiler_StaticMethods_foo =\n"+
			"function(JSCompiler_StaticMethods_foo$self) {\n"+
			"  return JSCompiler_StaticMethods_foo$self.x\n"+
			"};\n"+
			"var o = new a;\n"+
			"JSCompiler_StaticMethods_foo(o);\n")
}

// Java: DevirtualizeMethodsTest.noRewrite_forPropertyNamesAccessedReflectively()
func (c *DevirtualizeMethodsTest) TestNoRewrite_forPropertyNamesAccessedReflectively() {
	c.Test(
		c.Externs("function use() {}"),
		c.Srcs(
			"class C { m() {} n() {} }\n"+
				"const c = new C();\n"+
				"c.m();\n"+
				"c.n();\n"+
				"// this call should prevent devirtualizing m() but not n()\n"+
				"use(C.prototype, goog.reflect.objectProperty('m', C.prototype));\n"),
		c.Expected(
			"var JSCompiler_StaticMethods_n = function(JSCompiler_StaticMethods_n$self) {};\n"+
				"class C { m() {} }\n"+
				"const c = new C();\n"+
				"c.m();\n"+
				"JSCompiler_StaticMethods_n(c);\n"+
				"use(C.prototype, goog.reflect.objectProperty('m', C.prototype));\n"))
}

// Java: DevirtualizeMethodsTest.testNoRewriteIfDefinitionSiteBetween(String, String)
func (c *DevirtualizeMethodsTest) TestNoRewriteIfDefinitionSiteBetween(prefix string, suffix string) {
	c.TestSame(
		jutil.StringReplace(jutil.StringReplace("function a(){}\n"+
			"_PREFIX_a.prototype.foo = function() {return this.x}_SUFFIX_;\n"+
			"var o = new a;\n"+
			"o.foo()\n", "_PREFIX_", prefix),
			"_SUFFIX_", suffix))
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedIn_ifScope()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedIn_ifScope() {
	c.TestNoRewriteIfDefinitionSiteBetween("if (true) ", "")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedIn_loopScope()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedIn_loopScope() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestNoRewriteIfDefinitionSiteBetween("while (true) ", "")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedIn_switchScope()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedIn_switchScope() {
	c.TestNoRewriteIfDefinitionSiteBetween("switch (true) { case true: ", "; }")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedIn_functionScope()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedIn_functionScope() {
	c.TestNoRewriteIfDefinitionSiteBetween("function f() { ", "; }")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedIn_arrowFunctionScope()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedIn_arrowFunctionScope() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestNoRewriteIfDefinitionSiteBetween("() => ", "")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedIn_blockScope()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedIn_blockScope() {
	// Some declarations are block scoped in ES6 and so might have different values. This could make
	// multiple definitions with identical node structure behave differently.
	c.TestNoRewriteIfDefinitionSiteBetween("{ ", "; }")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedBy_andExpression_withLiteral()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedBy_andExpression_withLiteral() {
	// TODO(nickreid): This may be unnecessarily restrictive. We could probably rely on the
	// definitions being identical or not to filter this case.
	c.TestNoRewriteIfDefinitionSiteBetween("", " && function() {}")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedBy_andExpression_withReference()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedBy_andExpression_withReference() {
	c.TestNoRewriteIfDefinitionSiteBetween("", " && bar")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedBy_orExpression_withLiteral()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedBy_orExpression_withLiteral() {
	// TODO(nickreid): This may be unnecessarily restrictive. We could probably rely on the
	// definitions being identical or not to filter this case.
	c.TestNoRewriteIfDefinitionSiteBetween("", " || function() {}")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedBy_orExpression_withReference()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedBy_orExpression_withReference() {
	c.TestNoRewriteIfDefinitionSiteBetween("", " || bar")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedBy_ternaryExpression_withLiteral()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedBy_ternaryExpression_withLiteral() {
	// TODO(nickreid): This may be unnecessarily restrictive. We could probably rely on the
	// definitions being identical or not to filter this case.

	// Make the functions look different, just to be safe.
	c.TestNoRewriteIfDefinitionSiteBetween("", " ? function() { this.a; } : function() { this.b; }")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedBy_ternaryExpression_withReference()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedBy_ternaryExpression_withReference() {
	c.TestNoRewriteIfDefinitionSiteBetween("", " ? function() { } : bar")
}

// Java: DevirtualizeMethodsTest.testRewritePreservesFunctionKind(String)
func (c *DevirtualizeMethodsTest) TestRewritePreservesFunctionKind(fnKeyword string) {
	c.Test(
		c.Srcs(
			jutil.StringReplace("function a(){}\n"+
				"a.prototype.foo = FN_KEYWORD() { return 0; };\n"+
				"\n"+
				"(new a()).foo();\n", "FN_KEYWORD", fnKeyword)),
		c.Expected(
			jutil.StringReplace("function a(){}\n"+
				"var JSCompiler_StaticMethods_foo =\n"+
				"    FN_KEYWORD(JSCompiler_StaticMethods_foo$self) { return 0; };\n"+
				"\n"+
				"JSCompiler_StaticMethods_foo(new a());\n", "FN_KEYWORD", fnKeyword)))
}

// Java: DevirtualizeMethodsTest.testRewrite_preservesAsync()
func (c *DevirtualizeMethodsTest) TestRewrite_preservesAsync() {
	c.TestRewritePreservesFunctionKind("async function")
}

// Java: DevirtualizeMethodsTest.testRewrite_preservesAsyncGenerator()
func (c *DevirtualizeMethodsTest) TestRewrite_preservesAsyncGenerator() {
	c.TestRewritePreservesFunctionKind("async function*")
}

// Java: DevirtualizeMethodsTest.testRewrite_preservesGenerator()
func (c *DevirtualizeMethodsTest) TestRewrite_preservesGenerator() {
	c.TestRewritePreservesFunctionKind("function*")
}

// Java: DevirtualizeMethodsTest.testRewrite_replacesMultipleThisRefreences()
func (c *DevirtualizeMethodsTest) TestRewrite_replacesMultipleThisRefreences() {
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  a() {\n"+
				"     alert(this);\n"+
				"     alert(this, this);\n"+
				"  }\n"+
				"}\n"+
				"new Foo().a();\n"),
		c.Expected(
			"var JSCompiler_StaticMethods_a = function(JSCompiler_StaticMethods_a$self) {\n"+
				"  alert(JSCompiler_StaticMethods_a$self);\n"+
				"  alert(JSCompiler_StaticMethods_a$self, JSCompiler_StaticMethods_a$self);\n"+
				"};\n"+
				"\n"+
				"class Foo { }\n"+
				"JSCompiler_StaticMethods_a(new Foo());\n"))
}

// Java: DevirtualizeMethodsTest.testRewrite_defaultParams_usesThisReference()
func (c *DevirtualizeMethodsTest) TestRewrite_defaultParams_usesThisReference() {
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  bar(y = this) { return y; }\n"+
				"}\n"+
				"\n"+
				"// We need at least one normal call to trigger rewriting.\n"+
				"x.bar();\n"),
		c.Expected(
			"var JSCompiler_StaticMethods_bar = function(\n"+
				"    JSCompiler_StaticMethods_bar$self,\n"+
				"    y = JSCompiler_StaticMethods_bar$self) {\n"+
				"  return y;\n"+
				"};\n"+
				"class Foo { }\n"+
				"\n"+
				"JSCompiler_StaticMethods_bar(x);\n"))
}

// Java: DevirtualizeMethodsTest.testNoRewrite_defaultParams_usesSuperReference()
func (c *DevirtualizeMethodsTest) TestNoRewrite_defaultParams_usesSuperReference() {
	c.TestSame(
		"class Foo {\n" +
			"  bar(y = super.toString()) { return y; }\n" +
			"}\n" +
			"\n" +
			"// We need at least one normal call to trigger rewriting.\n" +
			"x.bar();\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifDefinedByArrow()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifDefinedByArrow() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame(
		"function a(){};\n" +
			"a.prototype.foo = () => 5;\n" +
			"\n" +
			"(new a()).foo();\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_namespaceFunctions()
func (c *DevirtualizeMethodsTest) TestNoRewrite_namespaceFunctions() {
	source := "function a(){}; a.foo = function() {return this.x}; a.foo()"
	c.TestSame(source)
}

// Java: DevirtualizeMethodsTest.testRewrite_ifMultipleIdenticalDefinitions()
func (c *DevirtualizeMethodsTest) TestRewrite_ifMultipleIdenticalDefinitions() {
	c.Test(
		"function A(){};\n"+
			"A.prototype.getFoo = function() { return 1; };\n"+
			"\n"+
			"function B(){};\n"+
			"B.prototype.getFoo = function() { return 1; };\n"+
			"\n"+
			"x.getFoo();\n",
		"function A() {};\n"+
			"var JSCompiler_StaticMethods_getFoo =\n"+
			"    function(JSCompiler_StaticMethods_getFoo$self) { return 1; };\n"+
			"\n"+
			"function B(){};\n"+
			"B.prototype.getFoo = function() { return 1 }; // Dead definition.\n"+
			"\n"+
			"JSCompiler_StaticMethods_getFoo(x);\n")
}

// Java: DevirtualizeMethodsTest.testRewrite_ifMultipleIdenticalDefinitions_ensuresDefinitionBeforeInvocations()
func (c *DevirtualizeMethodsTest) TestRewrite_ifMultipleIdenticalDefinitions_ensuresDefinitionBeforeInvocations() {
	c.Test(
		"function A() {};\n"+
			"A.prototype.getFoo = function() { return 1; };\n"+
			"\n"+
			"x.getFoo();\n"+
			"\n"+
			"function B() {};\n"+
			"B.prototype.getFoo = function() { return 1; };\n"+
			"\n"+
			"y.getFoo();\n",
		"function A() {};\n"+
			"var JSCompiler_StaticMethods_getFoo =\n"+
			"    function(JSCompiler_StaticMethods_getFoo$self) { return 1; };\n"+
			"\n"+
			"JSCompiler_StaticMethods_getFoo(x);\n"+
			"\n"+
			"function B() {};\n"+
			"B.prototype.getFoo=function() { return 1; }; // Dead definition.\n"+
			"\n"+
			"JSCompiler_StaticMethods_getFoo(y);\n")
}

// Java: DevirtualizeMethodsTest.testRewrite_ifMultipleIdenticalDefinitions_withThis()
func (c *DevirtualizeMethodsTest) TestRewrite_ifMultipleIdenticalDefinitions_withThis() {
	c.Test(
		"function A() {};\n"+
			"A.prototype.getFoo = function() { return this._foo + 1; };\n"+
			"\n"+
			"function B() {};\n"+
			"B.prototype.getFoo = function() { return this._foo + 1; };\n"+
			"\n"+
			"x.getFoo();\n",
		"function A() {};\n"+
			"var JSCompiler_StaticMethods_getFoo =\n"+
			"    function(JSCompiler_StaticMethods_getFoo$self) {\n"+
			"      return JSCompiler_StaticMethods_getFoo$self._foo + 1\n"+
			"    };\n"+
			"\n"+
			"function B() {};\n"+
			"B.prototype.getFoo = function() { return this._foo + 1 }; // Dead definition.\n"+
			"\n"+
			"JSCompiler_StaticMethods_getFoo(x);\n")
}

// Java: DevirtualizeMethodsTest.testRewrite_ifMultipleIdenticalDefinitions_withLocalNames()
func (c *DevirtualizeMethodsTest) TestRewrite_ifMultipleIdenticalDefinitions_withLocalNames() {
	// This case is included for completeness. `Normalization` is a prerequisite for this pass so
	// the naming conflict is resolved before devirtualization even begins. The change in names
	// invalidates devirtualization by making the definition subtrees unequal.
	c.Test(
		"function A() {};\n"+
			"A.prototype.getFoo = function f() { return f.prop; };\n"+
			"\n"+
			"function B() {};\n"+
			"B.prototype.getFoo = function f() { return f.prop; };\n"+
			"\n"+
			"x.getFoo();\n",
		"function A() {};\n"+
			"A.prototype.getFoo = function f() { return f.prop; };\n"+
			"\n"+
			"function B() {};\n"+
			"B.prototype.getFoo = function f$jscomp$1() { return f$jscomp$1.prop; };\n"+
			"\n"+
			"x.getFoo();\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_ifMultipleDistinctDefinitions()
func (c *DevirtualizeMethodsTest) TestNoRewrite_ifMultipleDistinctDefinitions() {
	c.TestSame(
		"function A(){}; A.prototype.getFoo = function() { return 1; };\n" +
			"function B(){}; B.prototype.getFoo = function() { return 2; };\n" +
			"var x = Math.random() ? new A() : new B();\n" +
			"alert(x.getFoo());\n")
}

/** Inputs for object literal tests. */
const DevirtualizeMethodsTest_NoRewritePrototypeObjectLiteralsTestInput_REGULAR = "b.prototype.foo = function() { return 1; }"
const DevirtualizeMethodsTest_NoRewritePrototypeObjectLiteralsTestInput_OBJ_LIT = "a.prototype = {foo : function() { return 2; }}"
const DevirtualizeMethodsTest_NoRewritePrototypeObjectLiteralsTestInput_CALL = "o.foo()"

// Java: DevirtualizeMethodsTest.testRewritePrototypeNoObjectLiterals()
func (c *DevirtualizeMethodsTest) TestRewritePrototypeNoObjectLiterals() {
	c.Test(
		DevirtualizeMethodsTest_SemicolonJoin(
			DevirtualizeMethodsTest_NoRewritePrototypeObjectLiteralsTestInput_REGULAR,
			DevirtualizeMethodsTest_NoRewritePrototypeObjectLiteralsTestInput_CALL),
		"var JSCompiler_StaticMethods_foo =\n"+
			"function(JSCompiler_StaticMethods_foo$self) { return 1; };\n"+
			"JSCompiler_StaticMethods_foo(o)\n")
}

// Java: DevirtualizeMethodsTest.testRewrite_definedUsingProtoObjectLit()
func (c *DevirtualizeMethodsTest) TestRewrite_definedUsingProtoObjectLit() {
	c.Test(
		DevirtualizeMethodsTest_SemicolonJoin(
			DevirtualizeMethodsTest_NoRewritePrototypeObjectLiteralsTestInput_OBJ_LIT,
			DevirtualizeMethodsTest_NoRewritePrototypeObjectLiteralsTestInput_CALL),
		"var JSCompiler_StaticMethods_foo = function(JSCompiler_StaticMethods_foo$self) {\n"+
			"  return 2;\n"+
			"};\n"+
			"a.prototype={};\n"+
			"\n"+
			"JSCompiler_StaticMethods_foo(o)\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_multipleDefinitions_definedUsingProtoObjectLit_definedUsingGetProp()
func (c *DevirtualizeMethodsTest) TestNoRewrite_multipleDefinitions_definedUsingProtoObjectLit_definedUsingGetProp() {
	c.TestSame(
		DevirtualizeMethodsTest_SemicolonJoin(
			DevirtualizeMethodsTest_NoRewritePrototypeObjectLiteralsTestInput_OBJ_LIT,
			DevirtualizeMethodsTest_NoRewritePrototypeObjectLiteralsTestInput_REGULAR,
			DevirtualizeMethodsTest_NoRewritePrototypeObjectLiteralsTestInput_CALL))
}

// Java: DevirtualizeMethodsTest.testNoRewrite_noDefinition()
func (c *DevirtualizeMethodsTest) TestNoRewrite_noDefinition() {
	c.TestSame("a.externalMethod()")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_externMethod()
func (c *DevirtualizeMethodsTest) TestNoRewrite_externMethod() {
	c.TestSame(
		c.Externs("A.prototype.externalMethod = function(){};"), //
		c.Srcs("o.externalMethod()"))
}

// Java: DevirtualizeMethodsTest.testNoRewrite_exportedMethod_viaCodingConvention()
func (c *DevirtualizeMethodsTest) TestNoRewrite_exportedMethod_viaCodingConvention() {
	// no rewriting without call; regardless of leading underscore
	c.TestSame("a.prototype.foo = function() {};")
	c.TestSame("a.prototype._foo = function() {};")

	// renames as expected
	c.Test(
		"function a() {} a.prototype.foo = function() {}; let o = new a; o.foo();",
		"function a() {}\n"+
			"var JSCompiler_StaticMethods_foo =\n"+
			"function(JSCompiler_StaticMethods_foo$self) {};\n"+
			"let o = new a;\n"+
			"JSCompiler_StaticMethods_foo(o);\n")

	// no renaming, as leading _ indicates exported symbol
	c.TestSame("function a() {} a.prototype._foo = function() {}; let o = new a; o._foo();")
}

// Java: DevirtualizeMethodsTest.testRewriteNoVarArgs()
func (c *DevirtualizeMethodsTest) TestRewriteNoVarArgs() {
	source :=
		"function a(){}\n" +
			"a.prototype.foo = function(args) {return args};\n" +
			"var o = new a;\n" +
			"o.foo()\n"

	expected :=
		"function a(){}\n" +
			"var JSCompiler_StaticMethods_foo =\n" +
			"  function(JSCompiler_StaticMethods_foo$self, args) {return args};\n" +
			"var o = new a;\n" +
			"JSCompiler_StaticMethods_foo(o)\n"

	c.Test(source, expected)
}

// Java: DevirtualizeMethodsTest.testRewrite_argInsideOptionalChainingCall()
func (c *DevirtualizeMethodsTest) TestRewrite_argInsideOptionalChainingCall() {
	source :=
		"function a(){}\n" +
			"a.prototype.foo = function(args) {return args};\n" +
			"var o = new a;\n" +
			"a?.(o.foo())\n"

	expected :=
		"function a(){}\n" +
			"var JSCompiler_StaticMethods_foo =\n" +
			"  function(JSCompiler_StaticMethods_foo$self, args) {return args};\n" +
			"var o = new a;\n" +
			"a?.(JSCompiler_StaticMethods_foo(o))\n"

	c.Test(source, expected)
}

// Java: DevirtualizeMethodsTest.testRewrite_lhsOfOptionalChainingCall()
func (c *DevirtualizeMethodsTest) TestRewrite_lhsOfOptionalChainingCall() {
	source :=
		"function a(){}\n" +
			"a.prototype.foo = function(args) {return args};\n" +
			"var o = new a;\n" +
			"o.foo()?.a\n"

	expected :=
		"function a(){}\n" +
			"var JSCompiler_StaticMethods_foo =\n" +
			"  function(JSCompiler_StaticMethods_foo$self, args) {return args};\n" +
			"var o = new a;\n" +
			"JSCompiler_StaticMethods_foo(o)?.a\n"

	c.Test(source, expected)
}

// Java: DevirtualizeMethodsTest.testNoRewriteVarArgs()
func (c *DevirtualizeMethodsTest) TestNoRewriteVarArgs() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	source :=
		"function a(){}\n" +
			"a.prototype.foo = function(var_args) {return arguments};\n" +
			"var o = new a;\n" +
			"o.foo()\n"
	c.TestSame(source)
}

/** Inputs for invalidating reference tests. */
const DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE = "function a() {}\n" +
	"a.prototype.foo = function() {return this.x};\n" +
	"var o = new a;\n"

	// Java: DevirtualizeMethodsTest.testRewrite_callReference()
func (c *DevirtualizeMethodsTest) TestRewrite_callReference() {
	expected :=
		"function a(){}\n" +
			"var JSCompiler_StaticMethods_foo =\n" +
			"function(JSCompiler_StaticMethods_foo$self) {\n" +
			"  return JSCompiler_StaticMethods_foo$self.x\n" +
			"};\n" +
			"var o = new a;\n" +
			"JSCompiler_StaticMethods_foo(o);\n"

	c.Test(DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE+"o.foo()", expected)
}

// Java: DevirtualizeMethodsTest.testNoRewrite_noReferences()
func (c *DevirtualizeMethodsTest) TestNoRewrite_noReferences() {
	c.TestSame(DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE)
}

// Java: DevirtualizeMethodsTest.testNoRewrite_nonCallReference_viaGetprop()
func (c *DevirtualizeMethodsTest) TestNoRewrite_nonCallReference_viaGetprop() {
	c.TestSame(
		DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE +
			"o.foo(); // We need at least one normal call to trigger rewriting.\n" +
			"o.foo;\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_nonCallReference_viaDestructuring()
func (c *DevirtualizeMethodsTest) TestNoRewrite_nonCallReference_viaDestructuring() {
	c.TestSame(
		DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE +
			"o.foo(); // We need at least one normal call to trigger rewriting.\n" +
			"const {foo: x} = o;\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_nonCallReference_viaGetprop_usingFnCall()
func (c *DevirtualizeMethodsTest) TestNoRewrite_nonCallReference_viaGetprop_usingFnCall() {
	// TODO(nickreid): Add rewriting support for this.
	c.TestSame(
		DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE +
			"o.foo(); // We need at least one normal call to trigger rewriting.\n" +
			"o.foo.call(null);\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_nonCallReference_viaGetprop_usingFnApply()
func (c *DevirtualizeMethodsTest) TestNoRewrite_nonCallReference_viaGetprop_usingFnApply() {
	// TODO(nickreid): Add rewriting support for this.
	c.TestSame(
		DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE +
			"o.foo(); // We need at least one normal call to trigger rewriting.\n" +
			"o.foo.apply(null);\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_nonCallReference_viaGetprop_asArgument()
func (c *DevirtualizeMethodsTest) TestNoRewrite_nonCallReference_viaGetprop_asArgument() {
	c.TestSame(
		DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE +
			"o.foo(); // We need at least one normal call to trigger rewriting.\n" +
			"bar(o.foo, null);\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_nonCallReference_viaTaggedTemplateString()
func (c *DevirtualizeMethodsTest) TestNoRewrite_nonCallReference_viaTaggedTemplateString() {
	// TODO(nickreid): Add rewriting support for this.
	c.TestSame(
		DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE +
			"o.foo(); // We need at least one normal call to trigger rewriting.\n" +
			"o.foo`Hello World!`;\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_optChain()
func (c *DevirtualizeMethodsTest) TestNoRewrite_optChain() {
	c.TestSame(
		DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE +
			"o.foo(); // We need at least one normal call to trigger rewriting.\n" +
			"o.foo?.();\n")
}

// Java: DevirtualizeMethodsTest.testNoRewrite_nonCallReference_viaNew()
func (c *DevirtualizeMethodsTest) TestNoRewrite_nonCallReference_viaNew() {
	// TODO(nickreid): Add rewriting support for this.
	c.TestSame(
		DevirtualizeMethodsTest_NoRewriteNonCallReferenceTestInput_BASE +
			"o.foo();\n" +
			"new o.foo();\n")
}

/** Inputs for nested definition tests. */
const DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_PREFIX = "a.prototype.foo = function() {"
const DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_SUFFIX = "o.foo()"
const DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_INNER = "a.prototype.bar = function() {}; o.bar()"
const DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_EXPECTED_PREFIX = "var JSCompiler_StaticMethods_foo=\n" +
	"function(JSCompiler_StaticMethods_foo$self){\n"
const DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_EXPECTED_SUFFIX = "JSCompiler_StaticMethods_foo(o)"

// Java: DevirtualizeMethodsTest.testRewriteNoNestedFunction()
func (c *DevirtualizeMethodsTest) TestRewriteNoNestedFunction() {
	c.Test(
		DevirtualizeMethodsTest_SemicolonJoin(
			DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_PREFIX+"}",
			DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_SUFFIX,
			DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_INNER),
		DevirtualizeMethodsTest_SemicolonJoin(
			DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_EXPECTED_PREFIX+"}",
			DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_EXPECTED_SUFFIX,
			"var JSCompiler_StaticMethods_bar=\n"+
				"function(JSCompiler_StaticMethods_bar$self){}\n",
			"JSCompiler_StaticMethods_bar(o)"))
}

// Java: DevirtualizeMethodsTest.testNoRewriteNestedFunction()
func (c *DevirtualizeMethodsTest) TestNoRewriteNestedFunction() {
	c.Test(
		DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_PREFIX+
			DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_INNER+
			"};"+
			DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_SUFFIX,
		DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_EXPECTED_PREFIX+
			DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_INNER+
			"};"+
			DevirtualizeMethodsTest_NoRewriteNestedFunctionTestInput_EXPECTED_SUFFIX)
}

// Java: DevirtualizeMethodsTest.testRewrite_definedUsingGetProp_withArgs_callUsingGetProp()
func (c *DevirtualizeMethodsTest) TestRewrite_definedUsingGetProp_withArgs_callUsingGetProp() {
	source :=
		"function a(){}\n" +
			"a.prototype.foo = function(args) {return args};\n" +
			"var o = new a;\n" +
			"o.foo()\n"
	expected :=
		"function a(){}\n" +
			"var JSCompiler_StaticMethods_foo =\n" +
			"  function(JSCompiler_StaticMethods_foo$self, args) {return args};\n" +
			"var o = new a;\n" +
			"JSCompiler_StaticMethods_foo(o)\n"
	c.Test(source, expected)
}

// Java: DevirtualizeMethodsTest.testNoRewrite_optChainGetProp()
func (c *DevirtualizeMethodsTest) TestNoRewrite_optChainGetProp() {
	source :=
		"function a(){}\n" +
			"a.prototype.foo = function(args) {return args};\n" +
			"var o = new a;\n" +
			"o?.foo()\n"
	c.TestSame(source)
}

// Java: DevirtualizeMethodsTest.testRewrite_definedUsingGetElem_withArgs_callUsingGetProp()
func (c *DevirtualizeMethodsTest) TestRewrite_definedUsingGetElem_withArgs_callUsingGetProp() {
	source :=
		"function a(){}\n" +
			"a.prototype['foo'] = function(args) {return args};\n" +
			"var o = new a;\n" +
			"o.foo()\n"
	c.TestSame(source)
}

// Java: DevirtualizeMethodsTest.testNoRewrite_definedUsingGetProp_withArgs_noCall_bracketAccess()
func (c *DevirtualizeMethodsTest) TestNoRewrite_definedUsingGetProp_withArgs_noCall_bracketAccess() {
	source :=
		"function a(){}\n" +
			"a.prototype.foo = function(args) {return args};\n" +
			"var o = new a;\n" +
			"o['foo']\n"
	c.TestSame(source)
}

// Java: DevirtualizeMethodsTest.testNoRewrite_optChainGetElemAccess()
func (c *DevirtualizeMethodsTest) TestNoRewrite_optChainGetElemAccess() {
	source :=
		"function a(){}\n" +
			"a.prototype.foo = function(args) {return args};\n" +
			"var o = new a;\n" +
			"o?.['foo']\n"
	c.TestSame(source)
}

// Java: DevirtualizeMethodsTest.testNoRewrite_definedUsingGetElem_withArgs_noCall_bracketAccess()
func (c *DevirtualizeMethodsTest) TestNoRewrite_definedUsingGetElem_withArgs_noCall_bracketAccess() {
	source :=
		"function a(){}\n" +
			"a.prototype['foo'] = function(args) {return args};\n" +
			"var o = new a;\n" +
			"o['foo']\n"
	c.TestSame(source)
}

// Java: DevirtualizeMethodsTest.testRewrite_definedInExpression_thatCreatesScope_reportsScopeAsDeleted()
func (c *DevirtualizeMethodsTest) TestRewrite_definedInExpression_thatCreatesScope_reportsScopeAsDeleted() {
	c.Test(
		"(function() {}).prototype.foo = function() {extern();};\n"+
			"// A call is needed to trigger rewriting.\n"+
			"a.foo();\n",
		"var JSCompiler_StaticMethods_foo = function(JSCompiler_StaticMethods_foo$self) {\n"+
			"  extern();\n"+
			"};\n"+
			"JSCompiler_StaticMethods_foo(a);\n")
}

// Java: DevirtualizeMethodsTest.testRewrite_definedInExpression_withSideEffects()
func (c *DevirtualizeMethodsTest) TestRewrite_definedInExpression_withSideEffects() {
	// TODO(nickreid): Expect this not to be a rewrite (or confirm it's safe).
	c.Test(
		"extern().prototype.foo = function() { };\n"+
			"// A call is needed to trigger rewriting.\n"+
			"a.foo()\n",
		// We think this is risky because if `extern()` had side effects they'd be eliminated.
		"var JSCompiler_StaticMethods_foo=function(JSCompiler_StaticMethods_foo$self){};\n"+
			"JSCompiler_StaticMethods_foo(a)\n")
}

// Java: DevirtualizeMethodsTest.testRewrite_definedInExpression_withInvocation()
func (c *DevirtualizeMethodsTest) TestRewrite_definedInExpression_withInvocation() {
	c.Test(
		c.Srcs(
			"(\n"+
				"  class { bar() { } },\n"+
				"  a.bar()\n"+
				");\n"),
		c.Expected(
			"var JSCompiler_StaticMethods_bar = function(JSCompiler_StaticMethods_bar$self) {\n"+
				"};\n"+
				"(\n"+
				"  class { },\n"+
				"  JSCompiler_StaticMethods_bar(a)\n"+
				");\n"))
}

// Java: DevirtualizeMethodsTest.testNoRewrite_definedUsingStringKey_inPrototypeLiteral_usingSuper()
func (c *DevirtualizeMethodsTest) TestNoRewrite_definedUsingStringKey_inPrototypeLiteral_usingSuper() {
	// TODO(b/120452418): Add rewriting support for this.
	c.TestSame(
		c.Srcs(
			"function a(){}\n" +
				"a.prototype = {\n" +
				"// Don't use `super.bar()` because that might effect the test.\n" +
				"  foo() { return super.x; }\n" +
				"};\n" +
				"\n" +
				"// We need at least one normal call to trigger rewriting.\n" +
				"x.foo()\n"))
}

// Java: DevirtualizeMethodsTest.testRewrite_definedUsingClassMember_prototypeMethod()
func (c *DevirtualizeMethodsTest) TestRewrite_definedUsingClassMember_prototypeMethod() {
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  bar() { return 5; }\n"+
				"}\n"+
				"\n"+
				"x.bar();\n"),
		c.Expected(
			"var JSCompiler_StaticMethods_bar = function(JSCompiler_StaticMethods_bar$self) {\n"+
				"  return 5;\n"+
				"};\n"+
				"class Foo { }\n"+
				"\n"+
				"JSCompiler_StaticMethods_bar(x);\n"))
}

// Java: DevirtualizeMethodsTest.testRewrite_definedUsingClassMember_prototypeMethod_usingThis()
func (c *DevirtualizeMethodsTest) TestRewrite_definedUsingClassMember_prototypeMethod_usingThis() {
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  bar() { return this; }\n"+
				"}\n"+
				"\n"+
				"// We need at least one normal call to trigger rewriting.\n"+
				"x.bar();\n"),
		c.Expected(
			"var JSCompiler_StaticMethods_bar = function(JSCompiler_StaticMethods_bar$self) {\n"+
				"    return JSCompiler_StaticMethods_bar$self;\n"+
				"};\n"+
				"class Foo { }\n"+
				"\n"+
				"JSCompiler_StaticMethods_bar(x);\n"))
}

// Java: DevirtualizeMethodsTest.testNoRewrite_definedUsingClassMember_prototypeMethod_usingSuper()
func (c *DevirtualizeMethodsTest) TestNoRewrite_definedUsingClassMember_prototypeMethod_usingSuper() {
	// TODO(b/120452418): Add rewriting support for this.
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"// Don't use `super.bar()` because that might effect the test.\n" +
				"  bar() { return super.x; }\n" +
				"}\n" +
				"\n" +
				"// We need at least one normal call to trigger rewriting.\n" +
				"x.bar();\n"))
}

// Java: DevirtualizeMethodsTest.testNoRewrite_definedUsingClassMember_constructor()
func (c *DevirtualizeMethodsTest) TestNoRewrite_definedUsingClassMember_constructor() {
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"  constructor() { }\n" +
				"}\n" +
				"\n" +
				"// We need at least one normal call to trigger rewriting.\n" +
				"x.constructor();\n"))
}

// Java: DevirtualizeMethodsTest.testRewrite_definedUsingClassMember_staticMethod()
func (c *DevirtualizeMethodsTest) TestRewrite_definedUsingClassMember_staticMethod() {
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  static bar() { return 5; }\n"+
				"}\n"+
				"\n"+
				"// We need at least one normal call to trigger rewriting.\n"+
				"x.bar();\n"),
		c.Expected(
			"var JSCompiler_StaticMethods_bar = function(JSCompiler_StaticMethods_bar$self) {\n"+
				"  return 5;\n"+
				"};\n"+
				"class Foo { }\n"+
				"\n"+
				"JSCompiler_StaticMethods_bar(x);\n"))
}

// Java: DevirtualizeMethodsTest.testRewrite_definedUsingAssignment_staticMethod_onClass()
func (c *DevirtualizeMethodsTest) TestRewrite_definedUsingAssignment_staticMethod_onClass() {
	c.Test(
		c.Srcs(
			"class Foo { }\n"+
				"Foo.bar = function() { return 5; }\n"+
				"\n"+
				"// We need at least one normal call to trigger rewriting.\n"+
				"x.bar();\n"),
		c.Expected(
			"class Foo { }\n"+
				"\n"+
				"var JSCompiler_StaticMethods_bar = function(JSCompiler_StaticMethods_bar$self) {\n"+
				"  return 5;\n"+
				"};\n"+
				"\n"+
				"JSCompiler_StaticMethods_bar(x);\n"))
}

// Java: DevirtualizeMethodsTest.testRewrite_definedUsingAssignment_staticMethod_onFunction()
func (c *DevirtualizeMethodsTest) TestRewrite_definedUsingAssignment_staticMethod_onFunction() {
	for _, annotation := range []string{"/** @constructor */", "/** @interface */", "/** @record*/"} {
		c.TestRewrite_definedUsingAssignment_staticMethod_onFunctionWithAnnotation(annotation)
	}
}

// Java: DevirtualizeMethodsTest.testRewrite_definedUsingAssignment_staticMethod_onFunction(String)
func (c *DevirtualizeMethodsTest) TestRewrite_definedUsingAssignment_staticMethod_onFunctionWithAnnotation(annotation string) {
	c.DisableCompareJsDoc() // multistage compilation simplifies jsdoc
	c.Test(
		c.Srcs(
			annotation+
				"function Foo() { }\n"+
				"Foo.bar = function() { return 5; }\n"+
				"\n"+
				"// We need at least one normal call to trigger rewriting.\n"+
				"x.bar();\n"),
		c.Expected(
			"function Foo() { }\n"+
				"\n"+
				"var JSCompiler_StaticMethods_bar = function(JSCompiler_StaticMethods_bar$self) {\n"+
				"  return 5;\n"+
				"};\n"+
				"\n"+
				"JSCompiler_StaticMethods_bar(x);\n"))
}

// Java: DevirtualizeMethodsTest.testRewrite_definedUsingClassMember_staticMethod_usingThis()
func (c *DevirtualizeMethodsTest) TestRewrite_definedUsingClassMember_staticMethod_usingThis() {
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  static bar() { return this; }\n"+
				"}\n"+
				"\n"+
				"// We need at least one normal call to trigger rewriting.\n"+
				"x.bar();\n"),
		c.Expected(
			"var JSCompiler_StaticMethods_bar = function(JSCompiler_StaticMethods_bar$self) {\n"+
				"  return JSCompiler_StaticMethods_bar$self;\n"+
				"};\n"+
				"class Foo { }\n"+
				"\n"+
				"JSCompiler_StaticMethods_bar(x);\n"))
}

// Java: DevirtualizeMethodsTest.testNoRewrite_definedUsingClassMember_staticMethod_usingSuper()
func (c *DevirtualizeMethodsTest) TestNoRewrite_definedUsingClassMember_staticMethod_usingSuper() {
	// TODO(b/120452418): Add rewriting support for this.
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"// Don't use `super.bar()` because that might effect the test.\n" +
				"  static bar() { return super.x; }\n" +
				"}\n" +
				"\n" +
				"// We need at least one normal call to trigger rewriting.\n" +
				"x.bar();\n"))
}

// Java: DevirtualizeMethodsTest.testRewrite_callWithSuperReceiver_passesThis()
func (c *DevirtualizeMethodsTest) TestRewrite_callWithSuperReceiver_passesThis() {
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"// Don't use `super.bar()` because that might effect the test.\n"+
				"  bar() { return super.qux(); }\n"+
				"}\n"+
				"\n"+
				"class Tig {\n"+
				"  qux() { return 5; }\n"+
				"}\n"),
		c.Expected(
			"class Foo {\n"+
				"  bar() { return JSCompiler_StaticMethods_qux(this); }\n"+
				"}\n"+
				"\n"+
				"var JSCompiler_StaticMethods_qux = function(JSCompiler_StaticMethods_qux$self) {\n"+
				"  return 5;\n"+
				"};\n"+
				"class Tig { }\n"))
}

// Java: DevirtualizeMethodsTest.testNoRewrite_inClassWithLocalName_asClassName()
func (c *DevirtualizeMethodsTest) TestNoRewrite_inClassWithLocalName_asClassName() {
	c.TestSame(
		c.Srcs(
			"const Qux = class Foo {\n" +
				"// TODO(nickreid): Add rewriting support for this so long as the local name is not\n" +
				"// referenced.\n" +
				"  bar() { return 5; }\n" +
				"}\n" +
				"\n" +
				"// We need at least one normal call to trigger rewriting.\n" +
				"x.bar();\n"))
}

// Java: DevirtualizeMethodsTest.testPropDefinition_notOnQualifiedName_doesNotCrash()
func (c *DevirtualizeMethodsTest) TestPropDefinition_notOnQualifiedName_doesNotCrash() {
	c.TestSame(
		c.Srcs(
			"class Qux { }\n" +
				"\n" +
				"new Qux().prop = function() { }\n"))
}

// Java: DevirtualizeMethodsTest.testRewriteDeclWithConstJSDoc()
func (c *DevirtualizeMethodsTest) TestRewriteDeclWithConstJSDoc() {
	c.Test(
		"class C {\n"+
			"  /** @const */ foo() {}\n"+
			"}\n"+
			"o.foo();\n",
		"/** @const */\n"+
			"var JSCompiler_StaticMethods_foo =\n"+
			"  function(JSCompiler_StaticMethods_foo$self) {};\n"+
			"class C {}\n"+
			"JSCompiler_StaticMethods_foo(o)\n")
}

// Java: DevirtualizeMethodsTest.testNoRewriteGet1()
func (c *DevirtualizeMethodsTest) TestNoRewriteGet1() {
	// Getters and setter require special handling.
	c.TestSame("function a(){}; a.prototype = {get foo(){return f}}; var o = new a; o.foo()")
}

// Java: DevirtualizeMethodsTest.testNoRewriteGet2()
func (c *DevirtualizeMethodsTest) TestNoRewriteGet2() {
	// Getters and setter require special handling.
	c.TestSame("function a(){}; a.prototype = {get foo(){return 1}}; var o = new a; o.foo")
}

// Java: DevirtualizeMethodsTest.testNoRewriteSet1()
func (c *DevirtualizeMethodsTest) TestNoRewriteSet1() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	// Getters and setter require special handling.
	source := "function a(){}; a.prototype = {set foo(a){}}; var o = new a; o.foo()"
	c.TestSame(source)
}

// Java: DevirtualizeMethodsTest.testNoRewriteSet2()
func (c *DevirtualizeMethodsTest) TestNoRewriteSet2() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	// Getters and setter require special handling.
	source := "function a(){}; a.prototype = {set foo(a){}}; var o = new a; o.foo = 1"
	c.TestSame(source)
}

// Java: DevirtualizeMethodsTest.testNoRewrite_notImplementedMethod()
func (c *DevirtualizeMethodsTest) TestNoRewrite_notImplementedMethod() {
	c.TestSame("function a(){}; var o = new a; o.foo()")
}

// Java: DevirtualizeMethodsTest.testWrapper()
func (c *DevirtualizeMethodsTest) TestWrapper() {
	c.TestSame("(function() {})()")
}

// Java: DevirtualizeMethodsTest.testRewrite_nestedFunction_hasThisBoundCorrectly()
func (c *DevirtualizeMethodsTest) TestRewrite_nestedFunction_hasThisBoundCorrectly() {
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  bar() {\n"+
				"    return function() { return this; };\n"+
				"  }\n"+
				"}\n"+
				"\n"+
				"// We need at least one normal call to trigger rewriting.\n"+
				"x.bar();\n"),
		c.Expected(
			"var JSCompiler_StaticMethods_bar = function(JSCompiler_StaticMethods_bar$self) {\n"+
				"  return function() { return this; };\n"+
				"};\n"+
				"class Foo { }\n"+
				"\n"+
				"JSCompiler_StaticMethods_bar(x);\n"))
}

// Java: DevirtualizeMethodsTest.testRewrite_nestedArrow_hasThisBoundCorrectly()
func (c *DevirtualizeMethodsTest) TestRewrite_nestedArrow_hasThisBoundCorrectly() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  bar() {\n"+
				"    return () => this;\n"+
				"  }\n"+
				"}\n"+
				"\n"+
				"// We need at least one normal call to trigger rewriting.\n"+
				"x.bar();\n"),
		c.Expected(
			"var JSCompiler_StaticMethods_bar = function(JSCompiler_StaticMethods_bar$self) {\n"+
				"  return () => JSCompiler_StaticMethods_bar$self;\n"+
				"};\n"+
				"class Foo { }\n"+
				"\n"+
				"JSCompiler_StaticMethods_bar(x);\n"))
}

// Java: DevirtualizeMethodsTest.testExistenceOfAGetter_preventsDevirtualization()
func (c *DevirtualizeMethodsTest) TestExistenceOfAGetter_preventsDevirtualization() {
	c.DeclareAccessor("foo", AccessorSummary_PropertyAccessKind_GETTER_ONLY)

	// Imagine the getter returned a function.
	c.TestSame(
		"class Foo {\n" +
			"  foo() {}\n" +
			"}\n" +
			"x.foo();\n")
}

// Java: DevirtualizeMethodsTest.testExistenceOfASetter_preventsDevirtualization()
func (c *DevirtualizeMethodsTest) TestExistenceOfASetter_preventsDevirtualization() {
	c.DeclareAccessor("foo", AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	// This doesn't actually seem like a risk but it's hard to say, and other optimizations that use
	// optimize calls would be dangerous on setters.
	c.TestSame(
		"class Foo {\n" +
			"  foo() {}\n" +
			"}\n" +
			"x.foo();\n")
}

// Java: DevirtualizeMethodsTest.testThisProperty()
func (c *DevirtualizeMethodsTest) TestThisProperty() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame(
		"class Foo {\n" +
			"  constructor() {\n" +
			"    this.a = function b() { return 5; };\n" +
			"    this.plus1 = (arg) => arg + 1;\n" +
			"    this.tmp = this.plus1(1);\n" +
			"  }\n" +
			"}\n" +
			"console.log(new Foo().a());\n")
}

// Java: DevirtualizeMethodsTest.testNonStaticClassFieldNoRHS()
func (c *DevirtualizeMethodsTest) TestNonStaticClassFieldNoRHS() {
	c.TestSame(
		"class Foo {\n" +
			"  a;\n" +
			"}\n" +
			"console.log(new Foo().a);\n")
}

// Java: DevirtualizeMethodsTest.testNonStaticClassFieldNonFunction()
func (c *DevirtualizeMethodsTest) TestNonStaticClassFieldNonFunction() {
	c.TestSame(
		"class Foo {\n" +
			"  a = 2;\n" +
			"}\n" +
			"console.log(new Foo().a);\n")
}

// Java: DevirtualizeMethodsTest.testNonStaticClassFieldFunction()
func (c *DevirtualizeMethodsTest) TestNonStaticClassFieldFunction() {
	c.TestSame(
		"class Foo {\n" +
			"  a = function x() { return 5; };\n" +
			"}\n" +
			"console.log(new Foo().a);\n")
}

// Java: DevirtualizeMethodsTest.testStaticClassFieldNoRHS()
func (c *DevirtualizeMethodsTest) TestStaticClassFieldNoRHS() {
	c.TestSame(
		"class Foo {\n" +
			"  static a;\n" +
			"}\n" +
			"console.log(Foo.a);\n")
}

// Java: DevirtualizeMethodsTest.testStaticClassFieldNonFunction()
func (c *DevirtualizeMethodsTest) TestStaticClassFieldNonFunction() {
	c.TestSame(
		"class Foo {\n" +
			"  static a = 2;\n" +
			"}\n" +
			"console.log(Foo.a);\n")
}

// Java: DevirtualizeMethodsTest.testStaticClassFieldFunction()
func (c *DevirtualizeMethodsTest) TestStaticClassFieldFunction() {
	c.TestSame(
		"class Foo {\n" +
			"  static a = function x() { return 5; };\n" +
			"}\n" +
			"console.log(Foo.a);\n")
}

const DevirtualizeMethodsTest_ModuleTestInput_DEFINITION = "a.prototype.foo = function() {}"
const DevirtualizeMethodsTest_ModuleTestInput_USE = "x.foo()"
const DevirtualizeMethodsTest_ModuleTestInput_REWRITTEN_DEFINITION = "var JSCompiler_StaticMethods_foo=\n" +
	"function(JSCompiler_StaticMethods_foo$self){}\n"
const DevirtualizeMethodsTest_ModuleTestInput_REWRITTEN_USE = "JSCompiler_StaticMethods_foo(x)"

// Java: DevirtualizeMethodsTest.testRewriteSameModule1()
func (c *DevirtualizeMethodsTest) TestRewriteSameModule1() {
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk(DevirtualizeMethodsTest_SemicolonJoin(DevirtualizeMethodsTest_ModuleTestInput_DEFINITION, DevirtualizeMethodsTest_ModuleTestInput_USE)).
			// m2
			AddChunk("").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			DevirtualizeMethodsTest_SemicolonJoin(DevirtualizeMethodsTest_ModuleTestInput_REWRITTEN_DEFINITION, DevirtualizeMethodsTest_ModuleTestInput_REWRITTEN_USE),
			// m2
			""))
}

// Java: DevirtualizeMethodsTest.testRewriteSameModule2()
func (c *DevirtualizeMethodsTest) TestRewriteSameModule2() {
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("").
			// m2
			AddChunk(DevirtualizeMethodsTest_SemicolonJoin(DevirtualizeMethodsTest_ModuleTestInput_DEFINITION, DevirtualizeMethodsTest_ModuleTestInput_USE)).
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			DevirtualizeMethodsTest_SemicolonJoin(DevirtualizeMethodsTest_ModuleTestInput_REWRITTEN_DEFINITION, DevirtualizeMethodsTest_ModuleTestInput_REWRITTEN_USE)))
}

// Java: DevirtualizeMethodsTest.testRewriteSameModule3()
func (c *DevirtualizeMethodsTest) TestRewriteSameModule3() {
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk(DevirtualizeMethodsTest_SemicolonJoin(DevirtualizeMethodsTest_ModuleTestInput_USE, DevirtualizeMethodsTest_ModuleTestInput_DEFINITION)).
			// m2
			AddChunk("").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			DevirtualizeMethodsTest_SemicolonJoin(DevirtualizeMethodsTest_ModuleTestInput_REWRITTEN_USE, DevirtualizeMethodsTest_ModuleTestInput_REWRITTEN_DEFINITION),
			// m2
			""))
}

// Java: DevirtualizeMethodsTest.testRewrite_definitionModule_beforeUseModule()
func (c *DevirtualizeMethodsTest) TestRewrite_definitionModule_beforeUseModule() {
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk(DevirtualizeMethodsTest_ModuleTestInput_DEFINITION).
			// m2
			AddChunk(DevirtualizeMethodsTest_ModuleTestInput_USE).
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			DevirtualizeMethodsTest_ModuleTestInput_REWRITTEN_DEFINITION,
			// m2
			DevirtualizeMethodsTest_ModuleTestInput_REWRITTEN_USE))
}

// Java: DevirtualizeMethodsTest.testNoRewrite_definitionModule_afterUseModule()
func (c *DevirtualizeMethodsTest) TestNoRewrite_definitionModule_afterUseModule() {
	chunks :=
		JSChunkGraphBuilder_ForStar().
			AddChunk(DevirtualizeMethodsTest_ModuleTestInput_USE).
			AddChunk(DevirtualizeMethodsTest_ModuleTestInput_DEFINITION).
			Build()

	c.TestSame(c.Srcs(chunks))
}

// Java: DevirtualizeMethodsTest.getProcessor(Compiler)
func (c *DevirtualizeMethodsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return OptimizeCalls_Builder_().
		SetCompiler(compiler).
		SetConsiderExterns(false).
		AddPass(NewDevirtualizeMethods(compiler)).
		Build()
}

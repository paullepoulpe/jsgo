// Port of com/google/javascript/jscomp/GenerateExportsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/GenerateExportsTest.java)

/*
 * Copyright 2008 The Closure Compiler Authors.
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

// Generate exports unit test.
type GenerateExportsTest struct {
	CompilerTestCaseBase

	allowNonGlobalExports  bool
	exportSymbolFunction   *string
	exportPropertyFunction *string
}

const generateExportsTest_EXTERNS = `/** @const */ var goog = {};
goog.exportSymbol = function(a, b, c) {};
goog.exportProperty = function(a, b, c) {};
`

// Java: GenerateExportsTest()
func NewGenerateExportsTest() *GenerateExportsTest {
	c := &GenerateExportsTest{allowNonGlobalExports: true}
	c.InitCompilerTestCaseBaseWithExterns(c, generateExportsTest_EXTERNS)
	return c
}

func TestGenerateExportsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewGenerateExportsTest() })
}

// Java: GenerateExportsTest.getProcessor(Compiler)
func (c *GenerateExportsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewGenerateExports(
		compiler, c.allowNonGlobalExports, c.exportSymbolFunction, c.exportPropertyFunction)
}

// Java: GenerateExportsTest.setUp()
func (c *GenerateExportsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.allowNonGlobalExports = true
	c.exportSymbolFunction = jutil.StringPtr("goog.exportSymbol")
	c.exportPropertyFunction = jutil.StringPtr("goog.exportProperty")
	// TODO(b/76025401): since this pass sometimes runs after typechecking,
	// verify it correctly propagates type information.
	// enableTypeCheck();
	// enableTypeInfoValidation();
}

// Java: GenerateExportsTest.testExportSymbol()
func (c *GenerateExportsTest) TestExportSymbol() {
	c.Test(
		"/** @export */function foo() {}",
		`/** @export */ function foo(){}
goog.exportSymbol("foo",foo)
`)
}

// Java: GenerateExportsTest.testExportSymbolAndProperties()
func (c *GenerateExportsTest) TestExportSymbolAndProperties() {
	c.Test(
		`/** @export */function foo() {}
/** @export */foo.prototype.bar = function() {}
`,
		`/** @export */function foo(){}
goog.exportSymbol("foo",foo);
/** @export */foo.prototype.bar=function(){};
goog.exportProperty(foo.prototype,"bar",foo.prototype.bar)
`)
}

// Java: GenerateExportsTest.testExportPrototypeProperty()
func (c *GenerateExportsTest) TestExportPrototypeProperty() {
	c.Test(
		`function Foo() {}
/** @export */ Foo.prototype.bar = function() {};
`,
		`function Foo() {}
/** @export */ Foo.prototype.bar = function(){};
goog.exportProperty(Foo.prototype, 'bar', Foo.prototype.bar);
`)
}

// Java: GenerateExportsTest.testExportSymbolAndConstantProperties()
func (c *GenerateExportsTest) TestExportSymbolAndConstantProperties() {
	c.Test(
		`/** @export */function foo() {}
/** @export */foo.BAR = 5;
`,
		`/** @export */function foo(){}
goog.exportSymbol("foo",foo);
/** @export */foo.BAR=5;
goog.exportProperty(foo,"BAR",foo.BAR)
`)
}

// Java: GenerateExportsTest.testExportVars()
func (c *GenerateExportsTest) TestExportVars() {
	c.Test(
		"/** @export */var FOO = 5",
		`/** @export */var FOO=5;
goog.exportSymbol('FOO',FOO)
`)
}

// Java: GenerateExportsTest.testExportLet()
func (c *GenerateExportsTest) TestExportLet() {
	c.Test(
		"/** @export */let FOO = 5",
		`/** @export */let FOO = 5;
goog.exportSymbol('FOO', FOO)
`)
}

// Java: GenerateExportsTest.testExportConst()
func (c *GenerateExportsTest) TestExportConst() {
	c.Test(
		"/** @export */const FOO = 5",
		`/** @export */const FOO = 5;
goog.exportSymbol('FOO', FOO)
`)
}

// Java: GenerateExportsTest.testExportEs6ArrowFunction()
func (c *GenerateExportsTest) TestExportEs6ArrowFunction() {
	c.Test(
		"/** @export */var fn = ()=>{};",
		`/** @export */var fn = ()=>{};
goog.exportSymbol('fn', fn)
`)
}

// Java: GenerateExportsTest.testExportPublicField()
func (c *GenerateExportsTest) TestExportPublicField() {
	c.AllowExternsChanges()
	code := "class Foo { /** @export */ FIELD; }"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.FIELD"))

	code = "class Bar { /** @export */ MSG = 'str'; }"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.MSG"))
}

// Java: GenerateExportsTest.testExportAnonymousClassPublicField()
func (c *GenerateExportsTest) TestExportAnonymousClassPublicField() {
	c.AllowExternsChanges()
	code := "foo(class { /** @export */ x = 4; });"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.x"))
}

// Java: GenerateExportsTest.testExportDeclaredClassPrototypeProperty()
func (c *GenerateExportsTest) TestExportDeclaredClassPrototypeProperty() {
	c.Test(
		"class Foo {} /** @export */ Foo.prototype.x = 3;",
		`class Foo {}
/** @export */
Foo.prototype.x = 3;
goog.exportProperty(Foo.prototype, 'x', Foo.prototype.x);
`)
}

// Java: GenerateExportsTest.testExportConstructorProperty()
func (c *GenerateExportsTest) TestExportConstructorProperty() {
	c.AllowExternsChanges()
	code := `class Foo {
  constructor() {
    /** @export */
    this.x = 3;
  }
}
`
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.x;"))
}

// Java: GenerateExportsTest.testNoExport()
func (c *GenerateExportsTest) TestNoExport() {
	c.Test("var FOO = 5", "var FOO=5")
}

// Java: GenerateExportsTest.testNestedVarAssign()
func (c *GenerateExportsTest) TestNestedVarAssign() {
	c.allowNonGlobalExports = false
	c.TestError(
		`var BAR;
/** @export */ var FOO = BAR = 5
`,
		FindExportableNodes_NON_GLOBAL_ERROR)

	c.allowNonGlobalExports = true
	c.TestError(
		`var BAR;
/** @export */ var FOO = BAR = 5
`,
		FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
}

// Java: GenerateExportsTest.testNestedAssign()
func (c *GenerateExportsTest) TestNestedAssign() {
	c.allowNonGlobalExports = false
	c.TestError(
		`var BAR;var FOO = {};
/** @export */FOO.test = BAR = 5
`,
		FindExportableNodes_NON_GLOBAL_ERROR)

	c.allowNonGlobalExports = true
	c.TestError(
		`var BAR;
var FOO = {};
/** @export */FOO.test = BAR = 5
`,
		FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
}

// Java: GenerateExportsTest.testNonGlobalScopeExport1()
func (c *GenerateExportsTest) TestNonGlobalScopeExport1() {
	c.allowNonGlobalExports = false
	c.TestError("(function() { /** @export */var FOO = 5 })()", FindExportableNodes_NON_GLOBAL_ERROR)

	c.allowNonGlobalExports = true
	c.TestError(
		"(function() { /** @export */var FOO = 5 })()",
		FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
}

// Java: GenerateExportsTest.testNonGlobalScopeExport2()
func (c *GenerateExportsTest) TestNonGlobalScopeExport2() {
	c.allowNonGlobalExports = false
	c.TestError("var x = {/** @export */ A:function() {}}", FindExportableNodes_NON_GLOBAL_ERROR)
}

// Java: GenerateExportsTest.testExportClass()
func (c *GenerateExportsTest) TestExportClass() {
	c.Test(
		"/** @export */ function G() {} foo();",
		"/** @export */ function G() {} goog.exportSymbol('G', G); foo();")
}

// Java: GenerateExportsTest.testExportClassMember()
func (c *GenerateExportsTest) TestExportClassMember() {
	c.Test(
		`/** @export */ function F() {}
/** @export */ F.prototype.method = function() {};
`,
		`/** @export */ function F() {}
goog.exportSymbol('F', F);
/** @export */ F.prototype.method = function() {};
goog.exportProperty(F.prototype, 'method', F.prototype.method);
`)
}

// Java: GenerateExportsTest.testExportEs6ClassSymbol()
func (c *GenerateExportsTest) TestExportEs6ClassSymbol() {
	c.Test(
		"/** @export */ class G {} foo();",
		"/** @export */ class G {} goog.exportSymbol('G', G); foo();")

	c.Test(
		"/** @export */ G = class {}; foo();",
		"/** @export */ G = class {}; goog.exportSymbol('G', G); foo();")
}

// Java: GenerateExportsTest.testExportEs6ClassProperty()
func (c *GenerateExportsTest) TestExportEs6ClassProperty() {
	c.Test(
		`/** @export */ G = class {};
/** @export */ G.foo = class {};
`,
		`/** @export */ G = class {};
goog.exportSymbol('G', G);
/** @export */ G.foo = class {};
goog.exportProperty(G, 'foo', G.foo)
`)

	c.Test(
		`G = class {};
/** @export */ G.prototype.foo = class {};
`,
		`G = class {};
/** @export */ G.prototype.foo = class {};
goog.exportProperty(G.prototype, 'foo', G.prototype.foo)
`)
}

// Java: GenerateExportsTest.testExportEs6ClassMembers()
func (c *GenerateExportsTest) TestExportEs6ClassMembers() {
	c.Test(
		`/** @export */
class G {
  /** @export */ method() {}
}
`,
		`/** @export */ class G { /** @export */ method() {} }
goog.exportSymbol('G', G);
goog.exportProperty(G.prototype, 'method', G.prototype.method);
`)

	c.Test(
		`/** @export */
class G {
  /** @export */ static method() {}
}
`,
		`/** @export */ class G { /** @export */ static method() {} }
goog.exportSymbol('G', G);
goog.exportProperty(G, 'method', G.method);
`)
}

// Java: GenerateExportsTest.testExportEs6ClassMembersWithSameName()
func (c *GenerateExportsTest) TestExportEs6ClassMembersWithSameName() {
	// Regression test for b/113617023, where we were silently dropping exports of different member
	// functions with the same name.
	c.Test(
		`/** @export */
class G {
  /** @export */ method() {}
}
/** @export */
class H {
  /** @export */ method() {}
}
`,
		`/** @export */
class G {
  /** @export */ method() {}
}
goog.exportSymbol('G', G);
goog.exportProperty(G.prototype, 'method', G.prototype.method);
/** @export */
class H {
  /** @export */ method() {}
}
goog.exportSymbol('H', H);
goog.exportProperty(H.prototype, 'method', H.prototype.method);
`)
}

// Java: GenerateExportsTest.testExportEs6ClassMembersOnSameClass()
func (c *GenerateExportsTest) TestExportEs6ClassMembersOnSameClass() {
	c.Test(
		`class G {
  /** @export */ methodA() {}
  /** @export */ methodB() {}
}
`,
		`class G {
  /** @export */ methodA() {}
  /** @export */ methodB() {}
}
goog.exportProperty(G.prototype, 'methodB', G.prototype.methodB);
goog.exportProperty(G.prototype, 'methodA', G.prototype.methodA);
`)
}

// Java: GenerateExportsTest.testExportEs6ClassMemberInObjectLitProperty()
func (c *GenerateExportsTest) TestExportEs6ClassMemberInObjectLitProperty() {
	c.Test(
		`const obj = {
  prop: class G {
    /** @export */ method() {}
  }
};
`,
		`const obj = {
  prop: class G {
    /** @export */ method() {}
  }
};
goog.exportProperty(obj.prop.prototype, 'method', obj.prop.prototype.method);
`)
}

// Java: GenerateExportsTest.testCannotExportEs6ClassMemberInClassExpressionAssignedToGetElem()
func (c *GenerateExportsTest) TestCannotExportEs6ClassMemberInClassExpressionAssignedToGetElem() {
	c.TestError(
		`const ns = {};
ns['G'] = class G {
  /** @export */ method() {}
};
`,
		FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
}

// Java: GenerateExportsTest.testCannotExportEs6ClassMemberInClassExpressionWithoutLValue()
func (c *GenerateExportsTest) TestCannotExportEs6ClassMemberInClassExpressionWithoutLValue() {
	c.TestError(
		`use(class G {
  /** @export */ method() {}
});
`,
		FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
}

// Java: GenerateExportsTest.testCannotExportEs6ClassMemberInClassExpressionInReturnExpression()
func (c *GenerateExportsTest) TestCannotExportEs6ClassMemberInClassExpressionInReturnExpression() {
	c.TestError(
		`function f() {
  return class {
    /** @export */ method() {}
  }
};
`,
		FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
}

// Java: GenerateExportsTest.testCannotExportEs6StaticBlockInClassExpressionInReturnExpression()
func (c *GenerateExportsTest) TestCannotExportEs6StaticBlockInClassExpressionInReturnExpression() {
	c.TestError(
		`function f() {
  return class {
    /** @export */ static {}
  }
};
`,
		FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
}

// Java: GenerateExportsTest.testCannotExportStaticBlock()
func (c *GenerateExportsTest) TestCannotExportStaticBlock() {
	c.TestError(
		`class G {
/** @export */ static {}
}
`,
		FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
}

// Java: GenerateExportsTest.testExportPublicStaticField()
func (c *GenerateExportsTest) TestExportPublicStaticField() {
	c.Test(
		`class G {
  /** @export */ static x = 3;
}
`,
		"class G { /** @export */ static x = 3; } goog.exportSymbol('G.x', G.x);")

	c.Test(
		`const G = class {
  /** @export */ static x;
};
`,
		"const G = class { /** @export */ static x; }; goog.exportSymbol('G.x', G.x);")

	c.Test(
		`const G = class localOnlyName {
  /** @export */ static x;
};
`,
		`const G = class localOnlyName {
  /** @export */ static x;
};
goog.exportSymbol('G.x', G.x);
`)

	c.Test(
		`let ns = {};
ns = class Bar {
  /** @export */ static x;
};
`,
		`let ns = {};
ns = class Bar {
  /** @export */ static x;
};
goog.exportSymbol('ns.x', ns.x);
`)

	c.Test(
		`const ns = {};
ns.G = class {
  /** @export */ static x;
};
`,
		`const ns = {};
ns.G = class {
  /** @export */ static x;
};
goog.exportSymbol('ns.G.x', ns.G.x);
`)

	c.Test(
		`class G {
  /** @export */ static getRand = Math.random;
}
`,
		`class G {
  /** @export */
  static getRand = Math.random;
}
goog.exportSymbol('G.getRand', G.getRand);
`)
}

// Java: GenerateExportsTest.testCannotExportStaticFieldInAnonymousClass()
func (c *GenerateExportsTest) TestCannotExportStaticFieldInAnonymousClass() {
	c.TestError(
		"foo(class { /** @export */ static x = 4; });",
		FindExportableNodes_EXPORT_ANNOTATION_NOT_ALLOWED)
}

// Java: GenerateExportsTest.testOldStyleStaticField()
func (c *GenerateExportsTest) TestOldStyleStaticField() {
	c.Test(
		"class Foo {} /** @export */ Foo.x = 3;",
		"class Foo {} /** @export */ Foo.x = 3; goog.exportSymbol('Foo.x', Foo.x);")
}

// Java: GenerateExportsTest.testExportClassMemberInStaticBlock()
func (c *GenerateExportsTest) TestExportClassMemberInStaticBlock() {
	c.AllowExternsChanges()
	code := "class G { static { /** @export */ this.foo = 1;} }"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.foo;"))
}

// Java: GenerateExportsTest.testClassMemberExportDoesntConflictInStaticBlock()
func (c *GenerateExportsTest) TestClassMemberExportDoesntConflictInStaticBlock() {
	c.AllowExternsChanges()
	code := `/** @export */
class G {
  static {
// TODO(b/235871861): Either generate an export for this property or report an error for
// it.
    /** @export */ this.foo=1;
  }
}
`
	result := `/** @export */
class G {
  static {
    /** @export */ this.foo=1;}
  }
goog.exportSymbol('G', G);
`
	c.Test(code, result)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.foo; var goog;"))
}

// Java: GenerateExportsTest.testGoogScopeFunctionOutput()
func (c *GenerateExportsTest) TestGoogScopeFunctionOutput() {
	c.Test(
		"/** @export */ $jscomp.scope.foo = /** @export */ function() {}",
		`/** @export */ $jscomp.scope.foo = /** @export */ function() {};
goog.exportSymbol('$jscomp.scope.foo', $jscomp.scope.foo);
`)
}

// Java: GenerateExportsTest.testGoogScopeClassOutput()
func (c *GenerateExportsTest) TestGoogScopeClassOutput() {
	c.Test(
		"/** @export */ $jscomp.scope.foo = /** @export */ class {}",
		`/** @export */ $jscomp.scope.foo = /** @export */ class {};
goog.exportSymbol('$jscomp.scope.foo', $jscomp.scope.foo);
`)
}

// Java: GenerateExportsTest.testExportSubclass()
func (c *GenerateExportsTest) TestExportSubclass() {
	c.Test(
		`var goog = {}; function F() {}
/** @export */ function G() {} goog.inherits(G, F);
`,
		`var goog = {}; function F() {}
/** @export */ function G() {} goog.inherits(G, F); goog.exportSymbol('G', G);
`)
}

// Java: GenerateExportsTest.testExportEnum()
func (c *GenerateExportsTest) TestExportEnum() {
	// TODO(johnlenz): Issue 310, should the values also be externed?
	c.Test(
		"/** @enum {string}\n @export */ var E = {A:1, B:2};",
		"/** @enum {string}\n @export */ var E = {A:1, B:2}; goog.exportSymbol('E', E);")
}

// Java: GenerateExportsTest.testExportObjectLit1()
func (c *GenerateExportsTest) TestExportObjectLit1() {
	c.AllowExternsChanges()
	code := "var E = {/** @export */ A:1, B:2};"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.A;"))
}

// Java: GenerateExportsTest.testExportObjectLit2()
func (c *GenerateExportsTest) TestExportObjectLit2() {
	c.AllowExternsChanges()
	code := "var E = {/** @export */ get A() { return 1 }, B:2};"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.A;"))
}

// Java: GenerateExportsTest.testExportObjectLit3()
func (c *GenerateExportsTest) TestExportObjectLit3() {
	c.AllowExternsChanges()
	code := "var E = {/** @export */ set A(v) {}, B:2};"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.A;"))
}

// Java: GenerateExportsTest.testExportObjectLit4()
func (c *GenerateExportsTest) TestExportObjectLit4() {
	c.AllowExternsChanges()
	code := "var E = {/** @export */ A:function() {}, B:2};"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.A;"))
}

// Java: GenerateExportsTest.testExportObjectLit5()
func (c *GenerateExportsTest) TestExportObjectLit5() {
	c.AllowExternsChanges()
	code := "var E = {/** @export */ A() {}, B:2};"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.A;"))
}

// Java: GenerateExportsTest.testExportClassMember1()
func (c *GenerateExportsTest) TestExportClassMember1() {
	c.AllowExternsChanges()
	code := "var E = function() { /** @export */ this.foo = 1; };"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.foo;"))
}

// Java: GenerateExportsTest.testMemberExportDoesntConflict()
func (c *GenerateExportsTest) TestMemberExportDoesntConflict() {
	c.AllowExternsChanges()
	code := `var foo = function() { /** @export */ this.foo = 1; };
/** @export */ foo.method = function(){};
`
	result := `var foo = function() { /** @export */ this.foo = 1; };
/** @export */ foo.method = function(){};
goog.exportSymbol('foo.method', foo.method);
`
	c.Test(code, result)
	c.DisableCompareAsTree()
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.foo; var goog;"))
}

// Java: GenerateExportsTest.testExportClassMemberStub()
func (c *GenerateExportsTest) TestExportClassMemberStub() {
	c.AllowExternsChanges()
	code := "var E = function() { /** @export */ this.foo; };"
	c.TestSame(code)
	c.TestExternChanges(c.Srcs(code), c.Expected("Object.prototype.foo;"))
}

// Java: GenerateExportsTest.testExportExprResultProperty()
func (c *GenerateExportsTest) TestExportExprResultProperty() {
	c.allowNonGlobalExports = false
	c.TestSame(
		`/** @record */ function Foo() {}
/** @export {number} */ Foo.prototype.myprop;
`)
}

// Java: GenerateExportsTest.testRequiresExportSymbolDefinition()
func (c *GenerateExportsTest) TestRequiresExportSymbolDefinition() {
	c.Test(
		c.Externs(""),
		c.Srcs("/** @export */ function Foo() {}"),
		c.Error(GenerateExports_MISSING_GOOG_FOR_EXPORT))
}

// Java: GenerateExportsTest.testRequiresExportPropertyDefinition()
func (c *GenerateExportsTest) TestRequiresExportPropertyDefinition() {
	c.Test(
		c.Externs(""),
		c.Srcs("var a = {}; /** @export */ a.b = function() {}"),
		c.Error(GenerateExports_MISSING_GOOG_FOR_EXPORT))
}

// Java: GenerateExportsTest.testRequiresExportSymbolConvention_ifUsesExport()
func (c *GenerateExportsTest) TestRequiresExportSymbolConvention_ifUsesExport() {
	c.exportSymbolFunction = nil

	c.Test(
		c.Srcs("/** @export */ function Foo() {}"), c.Error(GenerateExports_MISSING_EXPORT_CONVENTION))

	c.exportSymbolFunction = jutil.StringPtr("exportSymbol")
	c.exportPropertyFunction = nil

	c.Test(
		c.Srcs("/** @export */ function Foo() {}"), c.Error(GenerateExports_MISSING_EXPORT_CONVENTION))
}

// Java: GenerateExportsTest.testDoesNotRequireExportSymbolConvention_ifNoExports()
func (c *GenerateExportsTest) TestDoesNotRequireExportSymbolConvention_ifNoExports() {
	c.exportSymbolFunction = nil

	c.TestSame(c.Srcs("function Foo() {}"))
}

// Java: GenerateExportsTest.testFindExportableNodes()
func (c *GenerateExportsTest) TestFindExportableNodes() {
	compiler := NewCompiler()
	root := compiler.ParseTestCode("/** @export */ function Foo() {}")
	pass := NewFindExportableNodes(compiler, true)
	NodeTraversal_Traverse(compiler, root, pass)
	AssertThat(pass.GetExports()).IsNotEmpty()
}

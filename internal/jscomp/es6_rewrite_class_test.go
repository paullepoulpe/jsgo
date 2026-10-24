// Port of com/google/javascript/jscomp/Es6RewriteClassTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/Es6RewriteClassTest.java)

/*
 * Copyright 2014 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type Es6RewriteClassTest struct {
	CompilerTestCaseBase
	es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation
}

// Java: Es6RewriteClassTest()
func NewEs6RewriteClassTest() *Es6RewriteClassTest {
	c := &Es6RewriteClassTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, es6RewriteClassTest_EXTERNS_BASE)
	return c
}

func TestEs6RewriteClassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteClassTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *Es6RewriteClassTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: Es6RewriteClassTest.es6RewriteClassTest_EXTERNS_BASE
var es6RewriteClassTest_EXTERNS_BASE = NewTestExternsBuilder().
	AddArguments().
	AddArray().
	AddFunction().
	AddConsole().
	AddAlert().
	AddObject().
	AddClosureExterns().
	AddJSCompLibraries().
	Build()

// Java: Es6RewriteClassTest.customSetUp()
func (c *Es6RewriteClassTest) CustomSetUp() {
	c.EnableNormalize()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
	c.es6SubclassTranspilation = CompilerOptions_Es6SubclassTranspilation_CONCISE_UNSAFE
}

// Java: Es6RewriteClassTest.getProcessor(Compiler)
func (c *Es6RewriteClassTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory("injectTranspilationRuntimeLibraries", func(c AbstractCompiler) CompilerPass { return NewInjectTranspilationRuntimeLibraries(c) }))
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("es6NormalizeClasses", func(c AbstractCompiler) CompilerPass { return NewEs6NormalizeClasses(c) }))
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("es6ConvertSuper", func(c AbstractCompiler) CompilerPass { return NewEs6ConvertSuper(c) }))
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"es6RewriteClass", func(AbstractCompiler) CompilerPass { return NewEs6RewriteClass(compiler, c.es6SubclassTranspilation) }))
	return optimizer
}

// Java: Es6RewriteClassTest.testSimpleClassStatement()
func (c *Es6RewriteClassTest) TestSimpleClassStatement() {
	c.Test("class C { }", "/** @constructor */ let C = function() {};")

	// `C` from `let C = function() {};`
	cName := c.GetNodeWithName(c.GetLastCompiler().GetJsRoot(), "C")
	AssertNode(cName).HasToken(rhino.Token_NAME)
	AssertNode(cName).HasColorThat().IsConstructor()

	// function() {}
	function := cName.GetOnlyChild()
	AssertNode(function).HasToken(rhino.Token_FUNCTION)
	AssertNode(function).HasColorThat().IsEqualTo(cName.GetColor())
}

// Java: Es6RewriteClassTest.testClassStatementWithConstructor()
func (c *Es6RewriteClassTest) TestClassStatementWithConstructor() {
	c.Test("class C { constructor() {} }", "/** @constructor */ let C = function() {};")

	// `C` from `let C = function() {};`
	cName := c.GetNodeWithName(c.GetLastCompiler().GetJsRoot(), "C")
	AssertNode(cName).HasToken(rhino.Token_NAME)
	AssertNode(cName).HasColorThat().IsConstructor()

	// function() {}
	function := cName.GetOnlyChild()
	AssertNode(function).HasToken(rhino.Token_FUNCTION)
	AssertNode(function).HasColorThat().IsEqualTo(cName.GetColor())
}

// Java: Es6RewriteClassTest.testClassStatementWithMethod()
func (c *Es6RewriteClassTest) TestClassStatementWithMethod() {
	c.Test(
		"class C { method() {}; }",
		`/** @constructor */
let C = function() {};
C.prototype.method = function() {};
`)

	c_ := c.GetNodeWithName(c.GetLastCompiler().GetJsRoot(), "C")

	// `C.prototype.method`
	cPrototypeMethod := c.GetNodeWithName(c.GetLastCompiler().GetJsRoot(), "method")

	// `C.prototype`
	cPrototype := cPrototypeMethod.GetFirstChild()
	AssertNode(cPrototype).MatchesQualifiedName("C.prototype")
	AssertThat(c_.GetColor().GetPrototypes()).ContainsExactly(cPrototype.GetColor())

	// `C`
	cName := cPrototype.GetFirstChild()
	AssertNode(cName).MatchesQualifiedName("C")
	AssertNode(cName).HasColorThat().IsEqualTo(c_.GetColor())
	AssertNode(cName).HasColorThat().IsConstructor()

	// `function() {}`
	function := cPrototypeMethod.GetNext()
	AssertNode(function).HasToken(rhino.Token_FUNCTION)
	AssertNode(function).HasColorThat().IsEqualTo(cPrototypeMethod.GetColor())
}

// Java: Es6RewriteClassTest.testEs5MethodsAreEnumerable()
func (c *Es6RewriteClassTest) TestEs5MethodsAreEnumerable() {
	// Won't fix (b/573609265): natively, class methods are non-enumerable, but the transpiled code
	// assigns them to the prototype and to the constructor, which makes them enumerable. For
	// example, `Object.keys(C.prototype)` is `["m"]` instead of `[]`. We decided to keep this
	// output, since defining every method as non-enumerable would cost code size.
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { m() {} static s() {} }",
		`/** @constructor */
let C = function() {};
C.prototype.m = function() {};
C.s = function() {};
`)
}

// Java: Es6RewriteClassTest.testClassStatementWithConstructorReferencingThis()
func (c *Es6RewriteClassTest) TestClassStatementWithConstructorReferencingThis() {
	c.Test(
		"class C { constructor(a) { this.a = a; } }",
		"/** @constructor */ let C = function(a) { this.a = a; };")
}

// Java: Es6RewriteClassTest.testClassStatementWithConstructorAndMethod()
func (c *Es6RewriteClassTest) TestClassStatementWithConstructorAndMethod() {
	c.Test(
		"class C { constructor() {} foo() {} }",
		`/** @constructor */
let C = function() {};
C.prototype.foo = function() {};
`)
}

// Java: Es6RewriteClassTest.testClassStatementWithConstructorAndMethods()
func (c *Es6RewriteClassTest) TestClassStatementWithConstructorAndMethods() {
	c.Test(
		"class C { constructor() {}; foo() {}; bar() {} }",
		`/** @constructor */
let C = function() {};
C.prototype.foo = function() {};
C.prototype.bar = function() {};
`)
}

// Java: Es6RewriteClassTest.testClassStatementWithTwoMethods()
func (c *Es6RewriteClassTest) TestClassStatementWithTwoMethods() {
	c.Test(
		"class C { foo() {}; bar() {} }",
		`/** @constructor */
let C = function() {};
C.prototype.foo = function() {};
C.prototype.bar = function() {};
`)
}

// Java: Es6RewriteClassTest.testClassStatementWithNonEmptyConstructorAndMethods()
func (c *Es6RewriteClassTest) TestClassStatementWithNonEmptyConstructorAndMethods() {
	c.Test(
		`class C {
  constructor(a) { this.a = a; }

  foo() { console.log(this.a); }

  bar() { alert(this.a); }
}
`,
		`/** @constructor */
let C = function(a) { this.a = a; };
C.prototype.foo = function() { console.log(this.a); };
C.prototype.bar = function() { alert(this.a); };
`)
}

// Java: Es6RewriteClassTest.testClassStatementsInDifferentBlockScopes()
func (c *Es6RewriteClassTest) TestClassStatementsInDifferentBlockScopes() {
	c.Test(
		`if (true) {
  class Foo{}
} else {
  class FooBar{}
}
`,
		`if (true) {
   /** @constructor */
   let Foo = function() {};
} else {
   /** @constructor */
   let FooBar = function() {};
}
`)
}

// Java: Es6RewriteClassTest.testAnonymousWithSuper()
func (c *Es6RewriteClassTest) TestAnonymousWithSuper() {
	c.EnableClosurePass()
	c.IgnoreWarnings(FunctionTypeBuilder_RESOLVED_TAG_EMPTY, TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	c.Test(
		`goog.forwardDeclare('D')
f(class extends D { f() { super.g() } })
`,
		`goog.forwardDeclare('D')
/** @constructor
 */
const CLASS_DECL$0 = function() {
  return D.apply(this, arguments) || this;
};
$jscomp.inherits(CLASS_DECL$0, D);
CLASS_DECL$0.prototype.f = function() { D.prototype.g.call(this); };
f(CLASS_DECL$0)
`)
}

// Java: Es6RewriteClassTest.testAnonymousWithSuper_reflectConstruct()
func (c *Es6RewriteClassTest) TestAnonymousWithSuper_reflectConstruct() {
	c.es6SubclassTranspilation = CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT
	c.EnableClosurePass()
	c.IgnoreWarnings(FunctionTypeBuilder_RESOLVED_TAG_EMPTY, TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	c.Test(
		`goog.forwardDeclare('D')
f(class extends D { f() { super.g() } })
`,
		`goog.forwardDeclare('D')
/** @constructor
 */
const CLASS_DECL$0 = function() {
  return $jscomp.construct(D, arguments, this.constructor);
};
$jscomp.inherits(CLASS_DECL$0, D);
CLASS_DECL$0.prototype.f = function() { D.prototype.g.call(this); };
f(CLASS_DECL$0)
`)
}

// Java: Es6RewriteClassTest.testComplexBaseWithConditional()
func (c *Es6RewriteClassTest) TestComplexBaseWithConditional() {
	c.Test(
		`/** @define {boolean} */
const COND = goog.define('COND', false);
class Base1 {}
class Base2 extends Base1 {}
/** @type {typeof Base1} */
const ActualBase = COND ? Base1 : Base2;
class Sub extends ActualBase {}
`,
		`/** @define {!JSDocSerializer_placeholder_type} */
const COND = goog.define('COND', false);
/** @constructor */ let Base1 = function() {
};
/**
 * @constructor
 */
let Base2 = function() {
  Base1.apply(this, arguments);
};
$jscomp.inherits(Base2, Base1);
const ActualBase = COND ? Base1 : Base2;
/**
 * @constructor
 */
let Sub = function() {
  ActualBase.apply(this, arguments);
};
$jscomp.inherits(Sub, ActualBase);
`)
}

// Java: Es6RewriteClassTest.testClassWithNoJsDoc()
func (c *Es6RewriteClassTest) TestClassWithNoJsDoc() {
	c.Test("class C { }", "/** @constructor */ let C = function() { };")
}

// Java: Es6RewriteClassTest.testClassWithSuppressJsDoc()
func (c *Es6RewriteClassTest) TestClassWithSuppressJsDoc() {
	c.Test(
		"/** @suppress {partialAlias} */ class C { }",
		"/** @constructor @suppress {partialAlias} */ let C = function() {};")
}

// Java: Es6RewriteClassTest.testInterfaceWithJsDoc()
func (c *Es6RewriteClassTest) TestInterfaceWithJsDoc() {
	c.EnableClosurePass()
	c.Test(
		`goog.forwardDeclare('X');
goog.forwardDeclare('Y');
/**
 * Converts Xs to Ys.
 * @interface
 */
class Converter {
  /**
   * @param {X} x
   * @return {Y}
   */
  convert(x) {}
}
`,
		`goog.forwardDeclare('X');
goog.forwardDeclare('Y');
/**
 * Converts Xs to Ys.
 * @interface
 */
let Converter = function() { };

Converter.prototype.convert = function(x) {};
`)
}

// Java: Es6RewriteClassTest.testNoTypeForParam()
func (c *Es6RewriteClassTest) TestNoTypeForParam() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		`class MdMenu {
  /**
   * @param c
   */
  set classList(c) {}
}
`,
		`/** @constructor */
let MdMenu=function(){};
$jscomp.global.Object.defineProperties(
    MdMenu.prototype,
    {
      classList: {
        configurable:true,
        enumerable:true,
        set:function(c){}
      }
    })
`)
}

// Java: Es6RewriteClassTest.testParamNameMismatch()
func (c *Es6RewriteClassTest) TestParamNameMismatch() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.IgnoreWarnings(FunctionTypeBuilder_INEXISTENT_PARAM)

	c.Test(
		`class MdMenu {
/**
* @param {boolean} classes
*/
set classList(c) {}
}
`,
		`/** @constructor */
let MdMenu=function(){};
$jscomp.global.Object.defineProperties(
    MdMenu.prototype,
    {
      classList: {
        configurable:true,
        enumerable:true,
       set:function(c){}
     }
   })
`)
}

// Java: Es6RewriteClassTest.testParamNameMismatchAndNoParamType()
func (c *Es6RewriteClassTest) TestParamNameMismatchAndNoParamType() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.IgnoreWarnings(FunctionTypeBuilder_INEXISTENT_PARAM)

	c.Test(
		`class MdMenu {
/**
* @param classes
*/
set classList(c) {}
}
`,
		`/** @constructor */
let MdMenu=function(){};
$jscomp.global.Object.defineProperties(
    MdMenu.prototype,
    {
      classList: {
        configurable:true,
        enumerable:true,
        set:function(c){}
      }
   })
`)
}

// Java: Es6RewriteClassTest.testRecordWithJsDoc()
func (c *Es6RewriteClassTest) TestRecordWithJsDoc() {
	c.EnableClosurePass()
	c.Test(
		`goog.forwardDeclare('X');
goog.forwardDeclare('Y');
/**
 * @record
 */
class Converter {
  /**
   * @param {X} x
   * @return {Y}
   */
  convert(x) {}
}
`,
		`goog.forwardDeclare('X');
goog.forwardDeclare('Y');
/**
 * @interface
 */
let Converter = function() { };

Converter.prototype.convert = function(x) {};
`)
}

// Java: Es6RewriteClassTest.testMemberWithJsDoc()
func (c *Es6RewriteClassTest) TestMemberWithJsDoc() {
	c.Test(
		"class C { /** @param {boolean} b */ foo(b) {} }",
		`/**
 * @constructor
 */
let C = function() {};

C.prototype.foo = function(b) {};
`)
}

// Java: Es6RewriteClassTest.testClassStatementInsideIf()
func (c *Es6RewriteClassTest) TestClassStatementInsideIf() {
	c.Test("if (foo) { class C { } }", "if (foo) { /** @constructor */ let C = function() {}; }")
}

// Java: Es6RewriteClassTest.testClassStatementInsideBlocklessIf()
func (c *Es6RewriteClassTest) TestClassStatementInsideBlocklessIf() {
	c.Test("if (foo) class C {}", "if (foo) { /** @constructor */ let C = function() {}; }")
}

// Class expressions that are the RHS of a 'var' statement.
//
// Java: Es6RewriteClassTest.testClassExpressionInVar()
func (c *Es6RewriteClassTest) TestClassExpressionInVar() {
	c.Test("var C = class { }", "/** @constructor */ var C = function() {}")
}

// Java: Es6RewriteClassTest.testClassExpressionInVarWithMethod()
func (c *Es6RewriteClassTest) TestClassExpressionInVarWithMethod() {
	c.Test(
		"var C = class { foo() {} }",
		`/** @constructor */ var C = function() {}

C.prototype.foo = function() {}
`)
}

// Java: Es6RewriteClassTest.testNamedClassExpressionInVar()
func (c *Es6RewriteClassTest) TestNamedClassExpressionInVar() {
	c.Test(
		"var C = class C { }",
		`/** @constructor */
var C = function() {};
`)
}

// Java: Es6RewriteClassTest.testNamedClassExpressionInVarWithMethod()
func (c *Es6RewriteClassTest) TestNamedClassExpressionInVarWithMethod() {
	c.Test(
		"var C = class C { foo() {} }",
		`/** @constructor */
var C = function() {};
C.prototype.foo = function() {};
`)
}

// Class expressions that are the RHS of an assignment.
//
// Java: Es6RewriteClassTest.testClassExpressionInAssignment()
func (c *Es6RewriteClassTest) TestClassExpressionInAssignment() {
	c.Test(
		`/** @const */ goog.example = {};
goog.example.C = class { }
`,
		`/** @const */ goog.example = {};
/** @constructor */ goog.example.C = function() {}
`)
}

// Java: Es6RewriteClassTest.testClassExpressionInAssignmentWithMethod()
func (c *Es6RewriteClassTest) TestClassExpressionInAssignmentWithMethod() {
	c.Test(
		`/** @const */ goog.example = {};
goog.example.C = class { foo() {} }
`,
		`/** @const */ goog.example = {};
/** @constructor */ goog.example.C = function() {}
goog.example.C.prototype.foo = function() {};
`)
}

// Java: Es6RewriteClassTest.testClassExpressionInAssignment_getElem()
func (c *Es6RewriteClassTest) TestClassExpressionInAssignment_getElem() {
	c.Test(
		"window['MediaSource'] = class {};",
		`/** @constructor */
const CLASS_DECL$0 = function() {};
window['MediaSource'] = CLASS_DECL$0;
`)
}

// Java: Es6RewriteClassTest.testClassExpressionImmediatelyInstantiated()
func (c *Es6RewriteClassTest) TestClassExpressionImmediatelyInstantiated() {
	c.Test(
		"var C = new (class {})();",
		`/** @constructor */
const CLASS_DECL$0=function(){};
var C=new CLASS_DECL$0
`)
}

// Java: Es6RewriteClassTest.testClassExpressionAssignedToComplexLhs()
func (c *Es6RewriteClassTest) TestClassExpressionAssignedToComplexLhs() {
	c.Test(
		"(condition ? obj1 : obj2).prop = class C { };",
		`/** @constructor */
const CLASS_DECL$0 = function(){};
(condition ? obj1 : obj2).prop = CLASS_DECL$0;
`)
}

// Java: Es6RewriteClassTest.testClassExpression_cannotConvert()
func (c *Es6RewriteClassTest) TestClassExpression_cannotConvert() {
	c.Test(
		"var C = new (foo || (foo = class { }))();",
		`var C = new (foo || (foo = (() => {
  /** @constructor */
  const CLASS_DECL$0 = function() {};
  return CLASS_DECL$0;
})()))();
`)
}

// Java: Es6RewriteClassTest.testClassExpressionInDefaultParamValue_canConvert()
func (c *Es6RewriteClassTest) TestClassExpressionInDefaultParamValue_canConvert() {
	c.Test(
		"function bar(foo = class { }) {};",
		`function bar(foo = (() => {
  /** @constructor */ const CLASS_DECL$0 = function() {};
  return CLASS_DECL$0; })()) {};
`)
}

// Java: Es6RewriteClassTest.testExtendsWithImplicitConstructor()
func (c *Es6RewriteClassTest) TestExtendsWithImplicitConstructor() {
	c.Test(
		"class D {} class C extends D {}",
		`/** @constructor */
let D = function() {};
/** @constructor
 */
let C = function() { D.apply(this, arguments); };
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testExtendsWithImplicitConstructorNoTypeInfo()
func (c *Es6RewriteClassTest) TestExtendsWithImplicitConstructorNoTypeInfo() {
	c.DisableTypeCheck()
	c.DisableTypeInfoValidation()
	c.Test(
		c.Srcs("class D {} class C extends D {}"),
		c.Expected(
			"/** @constructor */\nlet D = function() {};\n/** @constructor\n */\nlet C = function() {\n// even without type information we recognize `arguments` and avoid generating\n// `[...arguments]`\n  D.apply(this, arguments);\n};\n$jscomp.inherits(C, D);\n"))
}

// Java: Es6RewriteClassTest.testExtendsWithImplicitConstructor_knownSuperclass_doesNotEs6SubclassTranspilation()
func (c *Es6RewriteClassTest) TestExtendsWithImplicitConstructor_knownSuperclass_doesNotEs6SubclassTranspilation() {
	c.es6SubclassTranspilation = CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT
	c.Test(
		"class D {} class C extends D {}",
		`/** @constructor */
let D = function() {};
/** @constructor
 */
let C = function() { D.apply(this, arguments); };
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testExtendsWithSpreadsPassedToSuper()
func (c *Es6RewriteClassTest) TestExtendsWithSpreadsPassedToSuper() {
	c.Test(
		`class D {
  /** @param {...string|number} args*/
  constructor(...args) {
    /** @type {!Arguments} */
    this.args = arguments;
  }
}
class C extends D {
  /**
   * @param {!Array<number>} numbers
   * @param {!Array<string>} strings
   */
  constructor(numbers, strings) {
    super(...numbers, ...strings);
  }
}
`,
		`/**
 * @constructor
 */
let D = function(...args) {
  this.args = arguments;
};
/**
 * @constructor
 */
let C = function(numbers, strings) {
  D.apply(this, [...numbers, ...strings]);
};
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testExtendsWithExplicitSuperCall()
func (c *Es6RewriteClassTest) TestExtendsWithExplicitSuperCall() {
	c.Test(
		"class D {} class C extends D { constructor() { super(); } }",
		`/** @constructor */
let D = function() {};
/** @constructor */
let C = function() {
  D.call(this);
}
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testExtendsWithExplicitSuperCallWithArguments()
func (c *Es6RewriteClassTest) TestExtendsWithExplicitSuperCallWithArguments() {
	c.Test(
		`class D { constructor(strArg) {} }
class C extends D { constructor(word) { super(word); } }
`,
		`/** @constructor */
let D = function(strArg) {};
/** @constructor */
let C = function(word) {
  D.call(this, word);
}
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testDerivedConstructorCompletesWithUndefinedBeforeSuperCall()
func (c *Es6RewriteClassTest) TestDerivedConstructorCompletesWithUndefinedBeforeSuperCall() {
	// Won't fix (b/573571176): natively, `new RetUndef()`, `new EarlyReturn(true)` and
	// `new MaybeObj(false)` throw a ReferenceError, because the constructor completes with
	// `undefined` before calling `super()`. The transpiled constructors are plain functions, so
	// `new` returns an instance on which `Base` never ran. We decided not to fix this, since the
	// compiler isn't required to preserve errors in code that is already broken.
	// The type checker reports these returns as inconsistent with the constructor type.
	c.DisableTypeCheck()
	c.DisableTypeInfoValidation()
	c.Test(
		`class Base { constructor() {} }
class RetUndef extends Base { constructor() { return undefined; } }
class EarlyReturn extends Base { constructor(w) { if (w) return; super(); } }
class MaybeObj extends Base { constructor(w) { if (w) return {}; } }
`,
		`/** @constructor */
let Base = function() {};
/** @constructor */
let RetUndef = function() {
  return undefined;
};
$jscomp.inherits(RetUndef, Base);
/** @constructor */
let EarlyReturn = function(w) {
  if (w) {
    return;
  }
  Base.call(this);
};
$jscomp.inherits(EarlyReturn, Base);
/** @constructor */
let MaybeObj = function(w$jscomp$1) {
  if (w$jscomp$1) {
    return {};
  }
};
$jscomp.inherits(MaybeObj, Base);
`)
}

// Java: Es6RewriteClassTest.testExtends_nonGlobalClass()
func (c *Es6RewriteClassTest) TestExtends_nonGlobalClass() {
	c.Test(
		c.Srcs(
			`function f() {
  class C {}
  class D extends C {}
}
`),
		c.Expected(
			`function f() {
  /** @constructor */
  let C = function() {};
  /** @constructor */
  let D = function() {
    return C.apply(this, arguments) || this;
  };
  $jscomp.inherits(D, C);
}
`))
}

// Java: Es6RewriteClassTest.testExtends_nonGlobalClass_usingReflectConstruct()
func (c *Es6RewriteClassTest) TestExtends_nonGlobalClass_usingReflectConstruct() {
	c.es6SubclassTranspilation = CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT
	c.Test(
		c.Srcs(
			`function f() {
  class C {}
  class D extends C {}
}
`),
		c.Expected(
			`function f() {
  /** @constructor */
  let C = function() {};
  /** @constructor */
  let D = function() {
    return $jscomp.construct(C, arguments, this.constructor);
  };
  $jscomp.inherits(D, C);
}
`))
}

// Java: Es6RewriteClassTest.testExtendsExternsClass()
func (c *Es6RewriteClassTest) TestExtendsExternsClass() {
	c.Test(
		c.Externs(es6RewriteClassTest_EXTERNS_BASE, "const ns = {}; /** @constructor */ ns.D = function() {};"),
		c.Srcs("class C extends ns.D { }"),
		c.Expected(
			`/** @constructor
 */
let C = function() {
 return ns.D.apply(this, arguments) || this;
};
$jscomp.inherits(C, ns.D);
`))
}

// Java: Es6RewriteClassTest.testExtendForwardDeclaredClass()
func (c *Es6RewriteClassTest) TestExtendForwardDeclaredClass() {
	c.EnableClosurePass()
	c.TestWarning(
		`goog.forwardDeclare('ns.D');
class C extends ns.D { }
`,
		TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
}

// Java: Es6RewriteClassTest.testImplementForwardDeclaredInterface()
func (c *Es6RewriteClassTest) TestImplementForwardDeclaredInterface() {
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	c.EnableClosurePass()
	c.TestWarning(
		`goog.forwardDeclare('ns.D');
/** @implements {ns.D} */
class C { }
`,
		FunctionTypeBuilder_RESOLVED_TAG_EMPTY)
}

// Java: Es6RewriteClassTest.testExtendForwardDeclaredInterface()
func (c *Es6RewriteClassTest) TestExtendForwardDeclaredInterface() {
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	c.EnableClosurePass()
	c.TestWarning(
		`goog.forwardDeclare('ns.D');
/** @interface @extends {ns.D} */ class C { }
`,
		FunctionTypeBuilder_RESOLVED_TAG_EMPTY)
}

// Java: Es6RewriteClassTest.testExtendNonNativeErrorWithAutogeneratedConstructor()
func (c *Es6RewriteClassTest) TestExtendNonNativeErrorWithAutogeneratedConstructor() {
	c.Test(
		`class Error {
  /** @param {string} msg */
  constructor(msg) {
    /** @const */ this.message = msg;
  }
}
class C extends Error {}
`, // autogenerated constructor
		`/** @constructor */
let Error = function(msg) {
  /** @const */ this.message = msg;
};
/** @constructor
 */
let C = function() { Error.apply(this, arguments); };
$jscomp.inherits(C, Error);
`)
}

// Java: Es6RewriteClassTest.testExtendNonNativeErrorWithExplicitSuperCall()
func (c *Es6RewriteClassTest) TestExtendNonNativeErrorWithExplicitSuperCall() {
	c.Test(
		`class Error {
  /** @param {string} msg */
  constructor(msg) {
    /** @const */ this.message = msg;
  }
}
class C extends Error {
  constructor() {
    super('C error'); // explicit super() call
  }
}
`,
		`/** @constructor */
let Error = function(msg) {
  /** @const */ this.message = msg;
};
/** @constructor
 */
let C = function() { Error.call(this, 'C error'); };
$jscomp.inherits(C, Error);
`)
}

// Java: Es6RewriteClassTest.testExtendNativeErrorWithAutogeneratedConstructor()
func (c *Es6RewriteClassTest) TestExtendNativeErrorWithAutogeneratedConstructor() {
	c.IgnoreWarnings(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
	c.Test(
		"class C extends Error {}", // autogenerated constructor
		`/** @constructor
 */
let C = function() {
  var $jscomp$tmp$error$m1146332801$1;
  $jscomp$tmp$error$m1146332801$1 = Error.apply(this, arguments),
      this.message = $jscomp$tmp$error$m1146332801$1.message,
      ('stack' in $jscomp$tmp$error$m1146332801$1) && (this.stack = $jscomp$tmp$error$m1146332801$1.stack),
      this;
};
$jscomp.inherits(C, Error);
`)
}

// Java: Es6RewriteClassTest.testExtendNativeAggregateErrorWithAutogeneratedConstructor()
func (c *Es6RewriteClassTest) TestExtendNativeAggregateErrorWithAutogeneratedConstructor() {
	c.IgnoreWarnings(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
	c.Test(
		"class C extends AggregateError {}", // autogenerated constructor
		`/** @constructor
 */
let C = function() {
  var $jscomp$tmp$error$m1146332801$1;
  $jscomp$tmp$error$m1146332801$1 = AggregateError.apply(this, arguments),
      this.message = $jscomp$tmp$error$m1146332801$1.message,
      ('stack' in $jscomp$tmp$error$m1146332801$1) && (this.stack = $jscomp$tmp$error$m1146332801$1.stack),
      this;
};
$jscomp.inherits(C, AggregateError);
`)
}

// Java: Es6RewriteClassTest.testExtendNativeErrorExplicitSuperCall()
func (c *Es6RewriteClassTest) TestExtendNativeErrorExplicitSuperCall() {
	c.IgnoreWarnings(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
	c.Test(
		`class C extends Error {
  constructor() {
    var self = super('C error') || this; // explicit super() call in an expression
  }
}
`,
		`/** @constructor
 */
let C = function() {
  var $jscomp$tmp$error$m1146332801$1;
  var self =
      ($jscomp$tmp$error$m1146332801$1 = Error.call(this, 'C error'),
          this.message = $jscomp$tmp$error$m1146332801$1.message,
          ('stack' in $jscomp$tmp$error$m1146332801$1) && (this.stack = $jscomp$tmp$error$m1146332801$1.stack),
          this) ||
      this;
};
$jscomp.inherits(C, Error);
`)
}

// Java: Es6RewriteClassTest.testExtendNativeAggregateErrorExplicitSuperCall()
func (c *Es6RewriteClassTest) TestExtendNativeAggregateErrorExplicitSuperCall() {
	c.IgnoreWarnings(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
	c.Test(
		`class C extends AggregateError {
  constructor() {
    var self = super([new Error('msg')]) || this; // explicit super() call in an
// expression
  }
}
`,
		`/** @constructor
 */
let C = function() {
  var $jscomp$tmp$error$m1146332801$1;
  var self =
      ($jscomp$tmp$error$m1146332801$1 = AggregateError.call(this, [new Error('msg')]),
          this.message = $jscomp$tmp$error$m1146332801$1.message,
          ('stack' in $jscomp$tmp$error$m1146332801$1) && (this.stack = $jscomp$tmp$error$m1146332801$1.stack),
          this) ||
      this;
};
$jscomp.inherits(C, AggregateError);
`)
}

// Java: Es6RewriteClassTest.testDynamicExtends_callWithoutJSDoc()
func (c *Es6RewriteClassTest) TestDynamicExtends_callWithoutJSDoc() {
	c.TestWarning("class C extends foo() {}", TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC)
}

// Java: Es6RewriteClassTest.testDynamicExtends_callWithJSDoc()
func (c *Es6RewriteClassTest) TestDynamicExtends_callWithJSDoc() {
	c.Test(
		"/** @extends {Object} */ class C extends foo() {}",
		`const CLASS_EXTENDS$0 = foo();
/** @constructor */
let C = function() {
  return CLASS_EXTENDS$0.apply(this, arguments) || this;
};
$jscomp.inherits(C, CLASS_EXTENDS$0);
`)
}

// Java: Es6RewriteClassTest.testDynamicExtends_functionExpression_suppressCheckTypes()
func (c *Es6RewriteClassTest) TestDynamicExtends_functionExpression_suppressCheckTypes() {
	c.Test(
		"/** @suppress {checkTypes} */ class C extends function(){} {}",
		`const CLASS_EXTENDS$0 = function() {};
/** @constructor */
let C = function() {
  CLASS_EXTENDS$0.apply(this, arguments);
};
$jscomp.inherits(C, CLASS_EXTENDS$0);
`)
}

// Java: Es6RewriteClassTest.testDynamicExtends_logicalExpressionWithJSDoc()
func (c *Es6RewriteClassTest) TestDynamicExtends_logicalExpressionWithJSDoc() {
	c.Test(
		"class A {} class B {} /** @extends {Object} */ class C extends (foo ? A : B) {}",
		`/** @constructor */
let A = function() {};
/** @constructor */
let B = function() {};
const CLASS_EXTENDS$0 = foo ? A : B;

/** @constructor */
let C = function() {
  CLASS_EXTENDS$0.apply(this, arguments);
};
$jscomp.inherits(C, CLASS_EXTENDS$0);
`)
}

// Java: Es6RewriteClassTest.testExtendsInterface()
func (c *Es6RewriteClassTest) TestExtendsInterface() {
	c.Test(
		`/** @interface */
class D {
  f() {}
}
/** @interface */
class C extends D {
  g() {}
}
`,
		`/** @interface */
let D = function() {};
D.prototype.f = function() {};
/** @interface */
let C = function() {};
$jscomp.inherits(C,D);
C.prototype.g = function() {};
`)
}

// Java: Es6RewriteClassTest.testExtendsRecord()
func (c *Es6RewriteClassTest) TestExtendsRecord() {
	c.Test(
		`/** @record */
class D {
  f() {}
}
/** @record */
class C extends D {
  g() {}
}
`,
		`/** @interface */
let D = function() {};
D.prototype.f = function() {};
/**
 * @interface */
let C = function() {};
$jscomp.inherits(C,D);
C.prototype.g = function() {};
`)
}

// Java: Es6RewriteClassTest.testImplementsInterface()
func (c *Es6RewriteClassTest) TestImplementsInterface() {
	c.Test(
		`/** @interface */
class D {
  f() {}
}
/** @implements {D} */
class C {
  f() {console.log('hi');}
}
`,
		`/** @interface */
let D = function() {};
D.prototype.f = function() {};
/** @constructor */
let C = function() {};
C.prototype.f = function() {console.log('hi');};
`)
}

// Java: Es6RewriteClassTest.testSuperCall()
func (c *Es6RewriteClassTest) TestSuperCall() {
	c.Test(
		`class D {}
class C extends D {
  constructor() { SUPER: super(); }
}
`,
		`/** @constructor */
let D = function() {};
/** @constructor */
let C = function() {
  SUPER: { D.call(this); }
}
$jscomp.inherits(C, D);
`)

	dType := c.GetNodeWithName(c.GetLastCompiler().GetJsRoot(), "D").GetColor()
	cType := c.GetNodeWithName(c.GetLastCompiler().GetJsRoot(), "C").GetColor()

	// D.call(this);
	callNode := c.GetNodeMatchingLabel(c.GetLastCompiler().GetJsRoot(), "SUPER").GetOnlyChild().GetOnlyChild()
	AssertNode(callNode).HasToken(rhino.Token_CALL)
	AssertNode(callNode).HasColorThat().IsEqualTo(colors.StandardColors_NULL_OR_VOID)

	// D.call
	callee := callNode.GetFirstChild()
	AssertNode(callee).MatchesQualifiedName("D.call")
	AssertNode(callee).HasColorThat().IsEqualTo(colors.StandardColors_TOP_OBJECT)

	// D
	superDotGReplacement := callee.GetFirstChild()
	AssertNode(superDotGReplacement).MatchesQualifiedName("D")
	AssertNode(superDotGReplacement).HasColorThat().IsEqualTo(dType)

	// `this` node from `D.call(this)`
	thisNode := callee.GetNext()
	AssertNode(thisNode).HasToken(rhino.Token_THIS)
	AssertThat(cType.GetInstanceColors()).ContainsExactly(thisNode.GetColor())
}

// Java: Es6RewriteClassTest.testSuperCall_withArgument()
func (c *Es6RewriteClassTest) TestSuperCall_withArgument() {
	c.IgnoreWarnings(TypeCheck_WRONG_ARGUMENT_COUNT)
	c.Test(
		"class D {} class C extends D { constructor(string) { super(string); } }",
		`/** @constructor */
let D = function() {}
/** @constructor */
let C = function(string) {
  D.call(this,string);
}
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testSuperCall_withArgumentAndSubsequentCode()
func (c *Es6RewriteClassTest) TestSuperCall_withArgumentAndSubsequentCode() {
	c.IgnoreWarnings(TypeCheck_WRONG_ARGUMENT_COUNT)
	c.Test(
		"class D {} class C extends D { constructor(strArg, n) { super(strArg); this.n = n; } }",
		`/** @constructor */
let D = function() {}
/** @constructor */
let C = function(strArg, n) {
  D.call(this,strArg);
  this.n = n;
}
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testSuperMethodCall()
func (c *Es6RewriteClassTest) TestSuperMethodCall() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class D {}
class C extends D {
  constructor() { }
  foo() { return super.foo(); }
}
`,
		`/** @constructor */
let D = function() {}
/** @constructor */
let C = function() { }
$jscomp.inherits(C, D);
C.prototype.foo = function() {
  return D.prototype.foo.call(this);
}
`)
}

// Java: Es6RewriteClassTest.testSuperMethodCall_withArgument()
func (c *Es6RewriteClassTest) TestSuperMethodCall_withArgument() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class D {}
class C extends D {
  constructor() {}
  foo(bar) { return super.foo(bar); }
}
`,
		`/** @constructor */
let D = function() {}
/** @constructor */
let C = function() {};
$jscomp.inherits(C, D);
C.prototype.foo = function(bar) {
  return D.prototype.foo.call(this, bar);
}
`)
}

// Java: Es6RewriteClassTest.testSuperCall_inSubclassInsideMethod()
func (c *Es6RewriteClassTest) TestSuperCall_inSubclassInsideMethod() {
	c.Test(
		"class C { method() { class D extends C { constructor() { super(); }}}}",
		`/** @constructor */
let C = function() {}
C.prototype.method = function() {
  /** @constructor */
  let D = function() {
    C.call(this);
  }
  $jscomp.inherits(D, C);
};
`)
}

// Java: Es6RewriteClassTest.testSuperKnownNotToChangeThis()
func (c *Es6RewriteClassTest) TestSuperKnownNotToChangeThis() {
	c.Test(
		`class D {
  /** @param {string} word */
  constructor(word) {
    this.str = word;
    /** @type {?} */ this.n;
    return; // Empty return should not trigger this-changing behavior.
  }
}
class C extends D {
  /**
   * @param {string} strArg
   * @param {number} n
   */
  constructor(strArg, n) {
// This is nuts, but confirms that super() used in an expression works.
    super(strArg).n = n;
// Also confirm that an existing empty return is handled correctly.
    return;
  }
}
`,
		"/**\n * @constructor\n */\nlet D = function(word) {\n  this.str = word;\n  this.n;\n  return;\n}\n/**\n * @constructor\n */\nlet C = function(strArg, n) {\n  (D.call(this,strArg), this).n = n; // super() returns `this`.\n  return;\n}\n$jscomp.inherits(C, D);\n")
}

// Java: Es6RewriteClassTest.testSuperMightChangeThis()
func (c *Es6RewriteClassTest) TestSuperMightChangeThis() {
	// Class D is unknown, so we must assume its constructor could change `this`.
	c.EnableClosurePass()
	c.IgnoreWarnings(FunctionTypeBuilder_RESOLVED_TAG_EMPTY)
	c.Test(
		`goog.forwardDeclare('D');
class C extends D {
  constructor(strArg, n) {
// This is nuts, but confirms that super() used in an expression works.
    super(strArg).n = n;
// Also confirm that an existing empty return is handled correctly.
    return;
  }
}
`,
		`goog.forwardDeclare('D');
/** @constructor */
let C = function(strArg, n) {
  var $jscomp$super$this$m1146332801$0;
  ($jscomp$super$this$m1146332801$0 = D.call(this,strArg) || this).n = n;
  return $jscomp$super$this$m1146332801$0; // Duplicate because of existing return
// statement.
  return $jscomp$super$this$m1146332801$0;
}
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testAlternativeSuperCalls_withExplicitSuperclassCtor()
func (c *Es6RewriteClassTest) TestAlternativeSuperCalls_withExplicitSuperclassCtor() {
	c.Test(
		`class D {
  /** @param {string} nameStr */
  constructor(nameStr) {
    this.name = nameStr;
  }
}
class C extends D {
  /** @param {string} strArg
   * @param {number} n */
  constructor(strArg, n) {
    if (n >= 0) {
      super('positive: ' + strArg);
    } else {
      super('negative: ' + strArg);
    }
    this.n = n;
  }
}
`,
		`/** @constructor */
let D = function(nameStr) {
  this.name = nameStr;
}
/** @constructor */
let C = function(strArg, n) {
  if (n >= 0) {
    D.call(this, 'positive: ' + strArg);
  } else {
    D.call(this, 'negative: ' + strArg);
  }
  this.n = n;
}
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testAlternativeSuperCalls_withUnknkownSuperclass()
func (c *Es6RewriteClassTest) TestAlternativeSuperCalls_withUnknkownSuperclass() {
	// Class being extended is unknown, so we must assume super() could change the value of `this`.
	c.EnableClosurePass()
	c.IgnoreWarnings(FunctionTypeBuilder_RESOLVED_TAG_EMPTY)
	c.Test(
		`goog.forwardDeclare('D');
class C extends D {
  /** @param {string} strArg
   * @param {number} n */
  constructor(strArg, n) {
    if (n >= 0) {
      super('positive: ' + strArg);
    } else {
      super('negative: ' + strArg);
    }
    this.n = n;
  }
}
`,
		`goog.forwardDeclare('D');
/** @constructor */
let C = function(strArg, n) {
  var $jscomp$super$this$m1146332801$0;
  if (n >= 0) {
    $jscomp$super$this$m1146332801$0 = D.call(this, 'positive: ' + strArg) || this;
  } else {
    $jscomp$super$this$m1146332801$0 = D.call(this, 'negative: ' + strArg) || this;
  }
  $jscomp$super$this$m1146332801$0.n = n;
  return $jscomp$super$this$m1146332801$0;
}
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testAlternativeSuperCalls_withUnknownSuperclass_es6SubclassTranspilation()
func (c *Es6RewriteClassTest) TestAlternativeSuperCalls_withUnknownSuperclass_es6SubclassTranspilation() {
	// Class being extended is unknown, so we must assume super() could change the value of `this`.
	c.EnableClosurePass()
	c.IgnoreWarnings(FunctionTypeBuilder_RESOLVED_TAG_EMPTY)
	c.es6SubclassTranspilation = CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT
	c.Test(
		`goog.forwardDeclare('D');
class C extends D {
  /** @param {string} strArg
   * @param {number} n */
  constructor(strArg, n) {
    if (n >= 0) {
      super('positive: ' + strArg);
    } else {
      super('negative: ' + strArg);
    }
    this.n = n;
  }
}
`,
		`goog.forwardDeclare('D');
/** @constructor */
let C = function(strArg, n) {
  var $jscomp$super$this$m1146332801$0;
  if (n >= 0) {
    $jscomp$super$this$m1146332801$0 =
        $jscomp.construct(D, ["positive: " + strArg], this.constructor);
  } else {
    $jscomp$super$this$m1146332801$0 =
        $jscomp.construct(D, ["negative: " + strArg], this.constructor);
  }
  $jscomp$super$this$m1146332801$0.n = n;
  return $jscomp$super$this$m1146332801$0;
}
$jscomp.inherits(C, D);
`)
}

// Java: Es6RewriteClassTest.testPrimitiveConstructorReturns()
func (c *Es6RewriteClassTest) TestPrimitiveConstructorReturns() {
	// A base constructor ignores a primitive return value, but the transpiled function returns it,
	// so `Base.call(this, 0) || this` in a subclass can evaluate to `42`. Each such return is
	// reported as an error. A `return` in a nested function belongs to that function, so it isn't
	// reported.
	c.IgnoreWarnings(TypeValidator_TYPE_MISMATCH_WARNING)
	c.Test(
		c.Srcs(
			`class Base {
  constructor(v) {
    if (v === 1) return 42;
    if (v === 2) return null;
    this.f = function() { return -1; };
  }
}
class Sub extends Base {
  constructor() {
    super(0);
    this.g = function() { return 's'; };
  }
}
`),
		c.Error(Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN),
		c.Error(Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN))
	// A derived constructor that returns a primitive, before or after `super()`, throws a
	// TypeError, so it is reported as an error.
	c.Test(
		c.Srcs(
			`class Base {}
class NoSuper extends Base {
  constructor(w) { if (w) return 42; super(); }
}
class AfterSuper extends Base {
  constructor() { super(); return 's'; }
}
class ReturnsNull extends Base {
  constructor() { super(); return null; }
}
`),
		c.Error(Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN),
		c.Error(Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN),
		c.Error(Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN))
}

// Java: Es6RewriteClassTest.testUndefinedConstructorReturnsAreKept()
func (c *Es6RewriteClassTest) TestUndefinedConstructorReturnsAreKept() {
	// Both base and derived constructors may return `undefined`.
	c.IgnoreWarnings(TypeValidator_TYPE_MISMATCH_WARNING)
	c.Test(
		`class Base {
  constructor(v) { if (v) return undefined; }
}
class Sub extends Base {
  constructor(w) { super(0); if (w) return void 0; }
}
`,
		`/** @constructor */
let Base = function(v) {
  if (v) {
    return undefined;
  }
};
/** @constructor */
let Sub = function(w) {
  Base.call(this, 0);
  if (w) {
    return void 0;
  }
};
$jscomp.inherits(Sub, Base);
`)
}

// Java: Es6RewriteClassTest.testComputedSuper()
func (c *Es6RewriteClassTest) TestComputedSuper() {
	c.Test(
		`/** @unrestricted */
class Foo {
  ['m']() { return 1; }
}

/** @unrestricted */
class Bar extends Foo {
  ['m']() {
    RETURN: return super['m']() + 1;
  }
}
`,
		`/** @constructor */
let Foo = function() {};
Foo.prototype['m'] = function() { return 1; };
/** @constructor */
let Bar = function() { Foo.apply(this, arguments); };
$jscomp.inherits(Bar, Foo);
Bar.prototype['m'] = function () {
  RETURN: { return Foo.prototype['m'].call(this) + 1; }
};
`)

	fooType := c.GetNodeWithName(c.GetLastCompiler().GetJsRoot(), "Foo").GetColor()
	barType := c.GetNodeWithName(c.GetLastCompiler().GetJsRoot(), "Bar").GetColor()

	// Foo.prototype['m'].call(this)
	callNode := c.GetNodeMatchingLabel(c.GetLastCompiler().GetJsRoot(), "RETURN").
		GetFirstFirstChild().
		GetFirstChild()
	AssertNode(callNode).HasToken(rhino.Token_CALL)
	AssertNode(callNode).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// Foo.prototype['m'].call
	callee := callNode.GetFirstChild()
	AssertNode(callee).HasToken(rhino.Token_GETPROP)
	AssertNode(callee).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// Foo.prototype['m']
	property := callee.GetFirstChild()
	AssertNode(property).HasToken(rhino.Token_GETELEM)
	AssertNode(callee).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)

	// Foo.prototype
	prototype := property.GetFirstChild()
	AssertNode(prototype).MatchesQualifiedName("Foo.prototype").HasOriginalName("super")
	AssertThat(fooType.GetPrototypes()).ContainsExactly(prototype.GetColor())

	// Foo
	superDotGReplacement := prototype.GetFirstChild()
	AssertNode(superDotGReplacement).MatchesQualifiedName("Foo")
	AssertNode(superDotGReplacement).HasColorThat().IsEqualTo(fooType)

	// `this` node from `Foo.prototype['m'].call(this)`
	thisNode := callee.GetNext()
	AssertNode(thisNode).HasToken(rhino.Token_THIS)
	AssertThat(barType.GetInstanceColors()).ContainsExactly(thisNode.GetColor())
}

// Java: Es6RewriteClassTest.testSuperMethodInGetter()
func (c *Es6RewriteClassTest) TestSuperMethodInGetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		`class Base {
  method() {
    return 5;
  }
}

class Subclass extends Base {
  constructor() {
    super();
  }

  get x() {
    return super.method();
  }
}
`,
		`/** @constructor */
let Base = function() {};
Base.prototype.method = function() { return 5; };

/** @constructor */
let Subclass = function() { Base.call(this); };

$jscomp.inherits(Subclass, Base);
$jscomp.global.Object.defineProperties(Subclass.prototype, {
  x: {
    configurable:true,
    enumerable:true,
    get: function() { return Base.prototype.method.call(this); },
  }
});
`)
}

// Java: Es6RewriteClassTest.testOverrideOfGetterFromInterface()
func (c *Es6RewriteClassTest) TestOverrideOfGetterFromInterface() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		`/** @interface */
class Int {
  /** @return {number} */
  get x() {}
}

/** @implements {Int} */
class C {
  /** @override @return {number} */
  get x() {}
}
`,
		`/** @interface */
let Int = function() {};
$jscomp.global.Object.defineProperties(Int.prototype, {
  x: {
    configurable:true,
    enumerable:true,
    get: function() {},
  }
});

/** @constructor */
let C = function() {};

$jscomp.global.Object.defineProperties(C.prototype, {
  x: {
    configurable:true,
    enumerable:true,
    get: function() {},
  }
});
`)
}

// Java: Es6RewriteClassTest.testSuperMethodInSetter()
func (c *Es6RewriteClassTest) TestSuperMethodInSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		`class Base {
  constructor() {}
  method() {
    var x = 5;
  }
}

class Subclass extends Base {
  constructor() {
    super();
  }

  set x(num) {
    super.method();
  }
}
`,
		`/** @constructor */
let Base = function() {};
Base.prototype.method = function() { var x = 5; };

/** @constructor */
let Subclass = function() { Base.call(this); };

$jscomp.inherits(Subclass, Base);
$jscomp.global.Object.defineProperties(Subclass.prototype, {
  x: {
    configurable:true,
    enumerable:true,
    set: function(num) { Base.prototype.method.call(this); }
  }
});
`)
}

// Java: Es6RewriteClassTest.testExtendNativeClass_withExplicitConstructor()
func (c *Es6RewriteClassTest) TestExtendNativeClass_withExplicitConstructor() {
	// Function and other native classes cannot be correctly extended in transpiled form.
	// Test both explicit and automatically generated constructors.
	c.Test(
		`class FooPromise extends Promise {
  /** @param {string} msg */
  constructor(callbackArg, msg) {
    super(callbackArg);
    this.msg = msg;
  }
}
`,
		`/**
 * @constructor
 */
let FooPromise = function(callbackArg, msg) {
  var $jscomp$super$this$m1146332801$0;
  $jscomp$super$this$m1146332801$0 = $jscomp.construct(Promise, [callbackArg], this.constructor)
  $jscomp$super$this$m1146332801$0.msg = msg;
  return $jscomp$super$this$m1146332801$0;
}
$jscomp.inherits(FooPromise, Promise);

`)
}

// Java: Es6RewriteClassTest.testExtendNativeClass_withExplicitConstructor_withInnerFunction()
func (c *Es6RewriteClassTest) TestExtendNativeClass_withExplicitConstructor_withInnerFunction() {
	c.Test(
		`class FooPromise extends Promise {
  /** @param {string} msg */
  constructor(callbackArg, msg) {
    function inner() {} // hoisted, normalized
    super(callbackArg);
    this.msg = msg;
  }
}
`,
		`/**
 * @constructor
 */
let FooPromise = function(callbackArg, msg) {
// stays hoisted
  function inner() {}
// declaration created after function declaration
  var $jscomp$super$this$m1146332801$0;
  $jscomp$super$this$m1146332801$0 =
      $jscomp.construct(Promise, [callbackArg], this.constructor);
  $jscomp$super$this$m1146332801$0.msg = msg;
  return $jscomp$super$this$m1146332801$0;
}
$jscomp.inherits(FooPromise, Promise);
`)
}

// Java: Es6RewriteClassTest.testExtendNativeClass_withImplicitConstructor()
func (c *Es6RewriteClassTest) TestExtendNativeClass_withImplicitConstructor() {
	c.Test(
		"class FooPromise extends Promise {}",
		`/**
 * @constructor
 */
let FooPromise = function() {
  return $jscomp.construct(Promise, arguments, this.constructor)
}
$jscomp.inherits(FooPromise, Promise);
`)
}

// Java: Es6RewriteClassTest.testExtendObject_replaceSuperCallWithThis()
func (c *Es6RewriteClassTest) TestExtendObject_replaceSuperCallWithThis() {
	// Object can be correctly extended in transpiled form, but we don't want or need to call
	// the `Object()` constructor in place of `super()`. Just replace `super()` with `this` instead.
	// Test both explicit and automatically generated constructors.
	c.Test(
		`class Foo extends Object {
  /** @param {string} msg */
  constructor(msg) {
    super();
    this.msg = msg;
  }
}
`,
		`/**
 * @constructor
 */
let Foo = function(msg) {
  this; // super() replaced with its return value
  this.msg = msg;
};
$jscomp.inherits(Foo, Object);
`)
}

// Java: Es6RewriteClassTest.testExtendObject_withSideEffectfulSuperArgs()
func (c *Es6RewriteClassTest) TestExtendObject_withSideEffectfulSuperArgs() {
	c.Test(
		`class Foo extends Object {
  constructor() {
    super(foo());
  }
}
`,
		`/**
 * @constructor
 */
let Foo = function() {
  foo(), this;
};
$jscomp.inherits(Foo, Object);
`)
}

// Java: Es6RewriteClassTest.testExtendObject_withMultipleSideEffectfulSuperArgs()
func (c *Es6RewriteClassTest) TestExtendObject_withMultipleSideEffectfulSuperArgs() {
	c.IgnoreWarnings(TypeCheck_WRONG_ARGUMENT_COUNT)
	c.Test(
		`class Foo extends Object {
  constructor() {
    super(foo(), 1, bar());
  }
}
`,
		`/**
 * @constructor
 */
let Foo = function() {
  foo(), bar(), this;
};
$jscomp.inherits(Foo, Object);
`)
}

// Java: Es6RewriteClassTest.testExtendObject_withSpreadSideEffectfulSuperArgs()
func (c *Es6RewriteClassTest) TestExtendObject_withSpreadSideEffectfulSuperArgs() {
	c.Test(
		`class Foo extends Object {
  constructor() {
    super(...foo());
  }
}
`,
		`/**
 * @constructor
 */
let Foo = function() {
  foo(), this;
};
$jscomp.inherits(Foo, Object);
`)
}

// Java: Es6RewriteClassTest.testExtendObject_withImplicitConstructor()
func (c *Es6RewriteClassTest) TestExtendObject_withImplicitConstructor() {
	c.Test(
		"class Foo extends Object {}",
		`/**
 * @constructor
 */
let Foo = function() {
  this; // super.apply(this, arguments) replaced with its return value
};
$jscomp.inherits(Foo, Object);
`)
}

// Java: Es6RewriteClassTest.testExtendNonNativeObject_withSuperCall()
func (c *Es6RewriteClassTest) TestExtendNonNativeObject_withSuperCall() {
	// No special handling when Object is redefined.
	c.IgnoreWarnings(TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS, FunctionTypeBuilder_TYPE_REDEFINITION)
	c.Test(
		c.Externs("/** @const */ var $jscomp = {}; $jscomp.inherits = function (a, b) {};"),
		c.Srcs(
			`class Object {}
class Foo extends Object {
  /** @param {string} msg */
  constructor(msg) {
    super();
    this.msg = msg;
  }
}
`),
		c.Expected(
			`/**
 * @constructor
 */
let Object = function() {
};
/**
 * @constructor
 */
let Foo = function(msg) {
  Object.call(this);
  this.msg = msg;
};
$jscomp.inherits(Foo, Object);
`))
}

// Java: Es6RewriteClassTest.testExtendNonNativeObject_withImplicitConstructor()
func (c *Es6RewriteClassTest) TestExtendNonNativeObject_withImplicitConstructor() {
	c.IgnoreWarnings(TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS, FunctionTypeBuilder_TYPE_REDEFINITION)
	c.Test(
		c.Externs("/** @const */ var $jscomp = {}; $jscomp.inherits = function (a, b) {};"),
		c.Srcs(
			`class Object {}
class Foo extends Object {}
`), // autogenerated constructor
		c.Expected(
			`/**
 * @constructor
 */
let Object = function() {
};
/**
 * @constructor
 */
let Foo = function() {
  Object.apply(this, arguments); // all arguments passed on to super()
};
$jscomp.inherits(Foo, Object);
`))
}

// Java: Es6RewriteClassTest.testMultiNameClass_inVarDeclaration()
func (c *Es6RewriteClassTest) TestMultiNameClass_inVarDeclaration() {
	c.Test(
		"var F = class G {}",
		`/** @constructor */
var F = function() {};
`)
}

// Java: Es6RewriteClassTest.testMultiNameClass_inAssignment()
func (c *Es6RewriteClassTest) TestMultiNameClass_inAssignment() {
	c.Test(
		"F = class G {}",
		`/** @constructor */
F = function() {};
`)
}

// Java: Es6RewriteClassTest.testClassNested()
func (c *Es6RewriteClassTest) TestClassNested() {
	c.Test(
		"class C { f() { class D {} } }",
		`/** @constructor */
let C = function() {};
C.prototype.f = function() {
  /** @constructor */
  let D = function() {}
};
`)
}

// Java: Es6RewriteClassTest.testClassNestedWithSuperclass()
func (c *Es6RewriteClassTest) TestClassNestedWithSuperclass() {
	c.Test(
		"class C { f() { class D extends C {} } }",
		`/** @constructor */
let C = function() {};
C.prototype.f = function() {
  /** @constructor */
  let D = function() {
    C.apply(this, arguments);
  };
  $jscomp.inherits(D, C);
};
`)
}

// Java: Es6RewriteClassTest.testSuperGetProp()
func (c *Es6RewriteClassTest) TestSuperGetProp() {
	c.Test(
		"class D { d() {} } class C extends D { f() {var i = super.d;} }",
		`/** @constructor */
let D = function() {};
D.prototype.d = function() {};
/** @constructor */
let C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.prototype.f = function() {
  var i = D.prototype.d;
};
`)
}

// Java: Es6RewriteClassTest.testSuperGetElem()
func (c *Es6RewriteClassTest) TestSuperGetElem() {
	c.Test(
		`/** @unrestricted */
class D { ['d']() {} }
/** @unrestricted */
class C extends D { f() {var i = super['d'];} }
`,
		`/** @constructor */
let D = function() {};
D.prototype['d'] = function() {};
/** @constructor */
let C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.prototype.f = function() {
  var i = D.prototype['d'];
};
`)
}

// Java: Es6RewriteClassTest.testSuperGetPropInStaticMethod()
func (c *Es6RewriteClassTest) TestSuperGetPropInStaticMethod() {
	// NOTE: super.d refers to a *static* property D.d (which is not defined), rather than the
	// instance propery D.prototype.d, which is defined.
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class D { d() {} }
class C extends D { static f() {var i = super.d;} }
`,
		`/** @constructor */
let D = function() {};
D.prototype.d = function() {};
/** @constructor */
let C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.f = function() {
  var i = D.d;
};
`)
}

// Java: Es6RewriteClassTest.testSuperGetElemInStaticMethod()
func (c *Es6RewriteClassTest) TestSuperGetElemInStaticMethod() {
	c.Test(
		`/** @unrestricted */
class D { ['d']() {}}
/** @unrestricted */
class C extends D { static f() {var i = super['d'];} }
`,
		`/** @constructor */
let D = function() {};
D.prototype['d'] = function() {};
/**
 * @constructor */
let C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.f = function() {
  var i = D['d'];
};
`)
}

// Java: Es6RewriteClassTest.testSuperGetProp_returnUndefinedProperty()
func (c *Es6RewriteClassTest) TestSuperGetProp_returnUndefinedProperty() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class D {}
class C extends D { f() {return super.s;} }
`,
		`/** @constructor */
let D = function() {};
/** @constructor */
let C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.prototype.f = function() {
  return D.prototype.s;
};
`)
}

// Java: Es6RewriteClassTest.testSuperGetProp_useUndefinedPropInCall()
func (c *Es6RewriteClassTest) TestSuperGetProp_useUndefinedPropInCall() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class D {}
class C extends D { f() { m(super.s);} }
`,
		`/** @constructor */
let D = function() {};
/**  @constructor */
let C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.prototype.f = function() {
  m(D.prototype.s);
};
`)
}

// Java: Es6RewriteClassTest.testSuperGetProp_dereferenceUndefinedProperty()
func (c *Es6RewriteClassTest) TestSuperGetProp_dereferenceUndefinedProperty() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class D {}
class C extends D { foo() { return super.m.foo();} }
`,
		`/** @constructor */
let D = function() {};
/** @constructor */
let C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.prototype.foo = function() {
  return D.prototype.m.foo();
};
`)
}

// Java: Es6RewriteClassTest.testSuperGet_dereferenceUndefinedPropertyInStaticMethod()
func (c *Es6RewriteClassTest) TestSuperGet_dereferenceUndefinedPropertyInStaticMethod() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class D {}
class C extends D { static foo() { return super.m.foo();} }
`,
		`/** @constructor */
let D = function() {};
/** @constructor */
let C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.foo = function() {
  return D.m.foo();
};
`)
}

// Java: Es6RewriteClassTest.testStaticThis()
func (c *Es6RewriteClassTest) TestStaticThis() {
	c.Test(
		"class F { static f() { return this; } }",
		`/** @constructor */ let F = function() {}
/** @this {?} */ F.f = function() { return this; };
`)
}

// Java: Es6RewriteClassTest.testStaticMethods()
func (c *Es6RewriteClassTest) TestStaticMethods() {
	c.Test(
		"class C { static foo() {} }",
		"/** @constructor */ let C = function() {}; C.foo = function() {};")
}

// Java: Es6RewriteClassTest.testStaticMethods_withSameNamedInstanceMethod()
func (c *Es6RewriteClassTest) TestStaticMethods_withSameNamedInstanceMethod() {
	c.Test(
		"class C { static foo() {}; foo() {} }",
		`/** @constructor */
let C = function() {};

C.foo = function() {};

C.prototype.foo = function() {};
`)
}

// Java: Es6RewriteClassTest.testStaticMethods_withCallInsideInstanceMethod()
func (c *Es6RewriteClassTest) TestStaticMethods_withCallInsideInstanceMethod() {
	c.Test(
		"class C { static foo() {}; bar() { C.foo(); } }",
		`/** @constructor */
let C = function() {};

C.foo = function() {};

C.prototype.bar = function() { C.foo(); };
`)
}

// Java: Es6RewriteClassTest.testStaticInheritance()
func (c *Es6RewriteClassTest) TestStaticInheritance() {
	c.Test(
		`class D {
  static f() {}
}
class C extends D { constructor() {} }
C.f();
`,
		`/** @constructor */
let D = function() {};
D.f = function () {};
/** @constructor */
let C = function() {};
$jscomp.inherits(C, D);
C.f();
`)
}

// Java: Es6RewriteClassTest.testStaticInheritance_withSameNamedInstanceMethod()
func (c *Es6RewriteClassTest) TestStaticInheritance_withSameNamedInstanceMethod() {
	c.Test(
		`class D {
  static f() {}
}
class C extends D {
  constructor() {}
  f() {}
}
C.f();
`,
		`/** @constructor */
let D = function() {};
D.f = function() {};
/** @constructor */
let C = function() { };
$jscomp.inherits(C, D);
C.prototype.f = function() {};
C.f();
`)
}

// Java: Es6RewriteClassTest.testStaticInheritance_overrideStaticMethod()
func (c *Es6RewriteClassTest) TestStaticInheritance_overrideStaticMethod() {
	c.Test(
		`class D {
  static f() {}
}
class C extends D {
  constructor() {}
  static f() {}
  g() {}
}
`,
		`/** @constructor */
let D = function() {};
D.f = function() {};
/** @constructor */
let C = function() { };
$jscomp.inherits(C, D);
C.f = function() {};
C.prototype.g = function() {};
`)
}

// Java: Es6RewriteClassTest.testInheritFromExterns()
func (c *Es6RewriteClassTest) TestInheritFromExterns() {
	c.Test(
		c.Externs(
			es6RewriteClassTest_EXTERNS_BASE+
				`/** @constructor */ function ExternsClass() {}
ExternsClass.m = function() {};
`),
		c.Srcs("class CodeClass extends ExternsClass {}"),
		c.Expected(
			`/** @constructor */
let CodeClass = function() {
  return ExternsClass.apply(this, arguments) || this;
};
$jscomp.inherits(CodeClass,ExternsClass)
`))
}

// Java: Es6RewriteClassTest.testInheritFromExterns_withReflectConstruct()
func (c *Es6RewriteClassTest) TestInheritFromExterns_withReflectConstruct() {
	c.es6SubclassTranspilation = CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT
	c.Test(
		c.Externs(
			es6RewriteClassTest_EXTERNS_BASE+
				`/** @constructor */ function ExternsClass() {}
ExternsClass.m = function() {};
`),
		c.Srcs("class CodeClass extends ExternsClass {}"),
		c.Expected(
			`/** @constructor */
let CodeClass = function() {
  return $jscomp.construct(ExternsClass, arguments, this.constructor);
};
$jscomp.inherits(CodeClass,ExternsClass)
`))
}

// Make sure we don't crash on this code.
// https://github.com/google/closure-compiler/issues/752
//
// Java: Es6RewriteClassTest.testGithub752()
func (c *Es6RewriteClassTest) TestGithub752() {
	c.Test(
		"function f() { var a = b = class {};}",
		`function f() {
  /** @constructor */
  const CLASS_DECL$0 = function() {};
  var a = b = CLASS_DECL$0;
}
`)
}

// Java: Es6RewriteClassTest.testGithub752b()
func (c *Es6RewriteClassTest) TestGithub752b() {
	c.Test(
		"var ns = {}; function f() { var self = ns.Child = class {};}",
		`var ns = {};
function f() {
  /** @constructor */
  const CLASS_DECL$0 = function() {};
  var self = ns.Child = CLASS_DECL$0
}
`)
}

// Getters and setters are supported, both in object literals and in classes, but only if the
// output language is ES5.
//
// Java: Es6RewriteClassTest.testEs5GettersAndSettersClasses_withSingleGetter()
func (c *Es6RewriteClassTest) TestEs5GettersAndSettersClasses_withSingleGetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { get value() { return 0; } }",
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    get: function() {
      return 0;
    }
  }
});
`)
}

// Java: Es6RewriteClassTest.testEs5GettersAndSettersClasses_withSingleSetter()
func (c *Es6RewriteClassTest) TestEs5GettersAndSettersClasses_withSingleSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { set value(val) { var internalVal = val; } }",
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    set: function(val) {
      var internalVal = val;
    }
  }
});
`)
}

// Java: Es6RewriteClassTest.testEs5GettersAndSettersClasses_withGetterSetterPair()
func (c *Es6RewriteClassTest) TestEs5GettersAndSettersClasses_withGetterSetterPair() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`/** @unrestricted */
class C {
  set value(val) {
    this.internalVal = val;
  }
  get value() {
    return this.internalVal;
  }
}
`,
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    set: function(val) {
      this.internalVal = val;
    },
    get: function() {
      return this.internalVal;
    }
  }
});
`)
}

// Java: Es6RewriteClassTest.testEs5GettersAndSettersClasses_protoAccessor()
func (c *Es6RewriteClassTest) TestEs5GettersAndSettersClasses_protoAccessor() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`class C {
  get __proto__() {
    return null;
  }
  set __proto__(v) {}
}
`,
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  ['__proto__']: {
    configurable: true,
    enumerable: true,
    get: function() {
      return null;
    },
    set: function(v) {}
  }
});
`)
}

// Java: Es6RewriteClassTest.testEs5GettersAndSettersClasses_withTwoGetters()
func (c *Es6RewriteClassTest) TestEs5GettersAndSettersClasses_withTwoGetters() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`class C {
  get alwaysTwo() {
    return 2;
  }

  get alwaysThree() {
    return 3;
  }
}
`,
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  alwaysTwo: {
    configurable: true,
    enumerable: true,
    get: function() {
      return 2;
    }
  },
  alwaysThree: {
    configurable: true,
    enumerable: true,
    get: function() {
      return 3;
    }
  },
});
`)
}

// Java: Es6RewriteClassTest.testEs5GettersAndSettersOnClassesWithClassSideInheritance()
func (c *Es6RewriteClassTest) TestEs5GettersAndSettersOnClassesWithClassSideInheritance() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { static get value() {} }  class D extends C { static get value() {} }",
		`/** @constructor */
let C = function() {};
/** @nocollapse */
C.value;
$jscomp.global.Object.defineProperties(C, {
  value: {
    configurable: true,
    enumerable: true,
    get: function() {}
  }
});
/** @constructor */
let D = function() {
  C.apply(this, arguments);
};
/** @nocollapse */
D.value;
$jscomp.inherits(D, C);
$jscomp.global.Object.defineProperties(D, {
  value: {
    configurable: true,
    enumerable: true,
    get: function() {}
  }
});
`)
}

// Check that the types from the getter/setter are copied to the declaration on the prototype.
//
// Java: Es6RewriteClassTest.testEs5GettersAndSettersClassesWithTypes_withSingleGetter()
func (c *Es6RewriteClassTest) TestEs5GettersAndSettersClassesWithTypes_withSingleGetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { /** @return {number} */ get value() { return 0; } }",
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    get: function() {
      return 0;
    }
  }
});
`)
}

// Java: Es6RewriteClassTest.testEs5GettersAndSettersClassesWithTypes_withSingleSetter()
func (c *Es6RewriteClassTest) TestEs5GettersAndSettersClassesWithTypes_withSingleSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { /** @param {string} v */ set value(v) { } }",
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    set: function(v) {}
  }
});
`)
}

// Java: Es6RewriteClassTest.testEs5GettersAndSettersClassesWithTypes_withConflictingGetterSetterType()
func (c *Es6RewriteClassTest) TestEs5GettersAndSettersClassesWithTypes_withConflictingGetterSetterType() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.TestWarning(
		`class C {
  /** @return {string} */
  get value() { }

  /** @param {number} v */
  set value(v) { }
}
`,
		TypeCheck_CONFLICTING_GETTER_SETTER_TYPE)

	// Also verify what the actual output is
	c.DisableTypeCheck()
	c.DisableTypeInfoValidation()

	c.Test(
		`class C {
  /** @return {string} */
  get value() { }

  /** @param {number} v */
  set value(v) { }
}
`,
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    get: function() {},
    set:function(v){}
  }
});
`)
}

// Java: Es6RewriteClassTest.testEs5GetterWithExport()
func (c *Es6RewriteClassTest) TestEs5GetterWithExport() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { /** @export @return {string} */ get foo() {} }",
		`/**
 * @constructor
 */
let C = function() {}

$jscomp.global.Object.defineProperties(C.prototype, {
  foo: {
    configurable: true,
    enumerable: true,
    get: function() {},
  }
});
`)
}

// Java: Es6RewriteClassTest.testEs5SetterWithExport()
func (c *Es6RewriteClassTest) TestEs5SetterWithExport() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { /** @export @param {string} x */ set foo(x) {} }",
		`/**
 * @constructor
 */
let C = function() {}

$jscomp.global.Object.defineProperties(C.prototype, {
  foo: {
    configurable: true,
    enumerable: true,
    set: function(x) {},
  }
});
`)
}

// b/20536614
//
// Java: Es6RewriteClassTest.testStaticGetter()
func (c *Es6RewriteClassTest) TestStaticGetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { static get foo() {} }",
		`/** @constructor */
let C = function() {};
/** @nocollapse */
C.foo;
$jscomp.global.Object.defineProperties(C, {
  foo: {
    configurable: true,
    enumerable: true,
    get: function() {}
  }
})
`)
}

// Java: Es6RewriteClassTest.testStaticGetter_withClassSideInheritance()
func (c *Es6RewriteClassTest) TestStaticGetter_withClassSideInheritance() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`class C { static get foo() {} }
class Sub extends C {}
`,
		`/** @constructor */
let C = function() {};
/** @nocollapse */
C.foo;
$jscomp.global.Object.defineProperties(C, {
  foo: {
    configurable: true,
    enumerable: true,
    get: function() {}
  }
})

/** @constructor */
let Sub = function() {
  C.apply(this, arguments);
};
$jscomp.inherits(Sub, C)
`)
}

// Java: Es6RewriteClassTest.testStaticSetter()
func (c *Es6RewriteClassTest) TestStaticSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { static set foo(x) {} }",
		`/** @constructor */
let C = function() {};
/** @nocollapse */
C.foo;
$jscomp.global.Object.defineProperties(C, {
  foo: {
    configurable: true,
    enumerable: true,
    set: function(x) {}
  }
});
`)
}

// Java: Es6RewriteClassTest.testClassStaticComputedProps_withSetter()
func (c *Es6RewriteClassTest) TestClassStaticComputedProps_withSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"/** @unrestricted */ class C { static set [se()](val) {}}",
		`var COMP_FIELD$0 = se();
/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperty(C, COMP_FIELD$0, {
  configurable: true,
  enumerable: true,
  set: function(val) {},
});
`)
}

// Java: Es6RewriteClassTest.testClassStaticComputedProps_withGetter()
func (c *Es6RewriteClassTest) TestClassStaticComputedProps_withGetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"/** @unrestricted */ class C { static get [se()]() {}}",
		`var COMP_FIELD$0 = se();
/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperty(C, COMP_FIELD$0, {
  configurable: true,
  enumerable: true,
  get: function() {},
});
`)
}

// Java: Es6RewriteClassTest.testClassStaticComputedProps_withGetterSetterPair()
func (c *Es6RewriteClassTest) TestClassStaticComputedProps_withGetterSetterPair() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"/** @unrestricted */ class C { static get [se()]() {} static set [se()](val) {}}",
		`var COMP_FIELD$0 = se();
var COMP_FIELD$1 = se();
/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperty(C, COMP_FIELD$0, {
  configurable: true,
  enumerable: true,
  get: function() {},
});
$jscomp.global.Object.defineProperty(C, COMP_FIELD$1, {
  configurable: true,
  enumerable: true,
  set: function(val) {},
});
`)
}

// Java: Es6RewriteClassTest.testSuperDotPropWithoutExtends()
func (c *Es6RewriteClassTest) TestSuperDotPropWithoutExtends() {
	c.TestError(
		`class C {
  foo() { return super.x; }
}
`,
		TranspilationUtil_CANNOT_CONVERT_YET)
}

// Java: Es6RewriteClassTest.testClassComputedPropGetterAndSetter()
func (c *Es6RewriteClassTest) TestClassComputedPropGetterAndSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`/** @unrestricted */
class C {
  /** @return {boolean} */
  get [foo]() {}
  /** @param {boolean} val */
  set [foo](val) {}
}
`,
		`var COMP_FIELD$0 = foo;
var COMP_FIELD$1 = foo;
/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperty(C.prototype, COMP_FIELD$0, {
  configurable: true,
  enumerable: true,
  get: function() {},
});
$jscomp.global.Object.defineProperty(C.prototype, COMP_FIELD$1, {
  configurable: true,
  enumerable: true,
  set: function(val) {},
});
`)
}

// Java: Es6RewriteClassTest.testClassComputedPropGetterAndSetter_withConflictingGetterSetterType()
func (c *Es6RewriteClassTest) TestClassComputedPropGetterAndSetter_withConflictingGetterSetterType() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.TestNoWarning(
		`/** @unrestricted */
class C {
  /** @return {boolean} */
  get [foo]() {}
  /** @param {string} val */
  set [foo](val) {}
}
`)
}

// Java: Es6RewriteClassTest.testClassComputedPropGetterAndSetter_mixedWithOtherMembers()
func (c *Es6RewriteClassTest) TestClassComputedPropGetterAndSetter_mixedWithOtherMembers() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		c.Externs(
			NewTestExternsBuilder().
				AddObject().
				AddReflect().
				AddJSCompLibraries().
				AddExtra(
					`function one() {}
function two() {}
function three() {}
`).
				Build()),
		c.Srcs(
			`/** @unrestricted */
class C {
  get [one()]() { return true; }
  get a() {}
  [two()]() {}
  get b() {}
  c() {}
  set [three()](val) { return true; }
}
`),
		c.Expected(
			`var COMP_FIELD$0 = one();
var COMP_FIELD$1 = two();
var COMP_FIELD$2 = three();
/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperty(C.prototype, COMP_FIELD$0, {
  configurable: true,
  enumerable: true,
  get: function() {
    return true;
  },
});
$jscomp.global.Object.defineProperties(C.prototype, {
  a: {configurable: true, enumerable: true, get: function() {}},
});
C.prototype[COMP_FIELD$1] = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  b: {configurable: true, enumerable: true, get: function() {}},
});
C.prototype.c = function() {};
$jscomp.global.Object.defineProperty(C.prototype, COMP_FIELD$2, {
  configurable: true,
  enumerable: true,
  set: function(val) {
    return true;
  },
});
`))
}

// Java: Es6RewriteClassTest.testGetterFollowedByMethodWithSameName()
func (c *Es6RewriteClassTest) TestGetterFollowedByMethodWithSameName() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`/** @unrestricted */
class C {
  get clean() { return 1; }
  clean(x) { return 2; }
}
`,
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  clean: {
    configurable: true,
    enumerable: true,
    get: function() { return 1; }
  }
});
C.prototype.clean = function(x) { return 2; };
`)
}

// Java: Es6RewriteClassTest.testGetterFollowedByComputedGetterWithSameName()
func (c *Es6RewriteClassTest) TestGetterFollowedByComputedGetterWithSameName() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`/** @unrestricted */
class C {
  get clean() { return 1; }
  get ['clean']() { return 2; }
}
`,
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  clean: {
    configurable: true,
    enumerable: true,
    get: function() { return 1; }
  }
});
$jscomp.global.Object.defineProperty(C.prototype, "clean", {
  configurable: true,
  enumerable: true,
  get: function() { return 2; }
});
`)
}

// Java: Es6RewriteClassTest.testMethodFollowedByGetterWithSameName()
func (c *Es6RewriteClassTest) TestMethodFollowedByGetterWithSameName() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`/** @unrestricted @suppress {checkTypes} */
class C {
  clean(x) { return 1; }
  get clean() { return 2; }
}
`,
		`/** @constructor */
let C = function() {};
C.prototype.clean = function(x) { return 1; };
$jscomp.global.Object.defineProperties(C.prototype, {
  clean: {
    configurable: true,
    enumerable: true,
    get: function() { return 2; }
  }
});
`)
}

// Java: Es6RewriteClassTest.testDuplicateGettersWithSameName()
func (c *Es6RewriteClassTest) TestDuplicateGettersWithSameName() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`/** @unrestricted */
class C {
  get clean() { return 1; }
  get clean() { return 2; }
}
`,
		`/** @constructor */
let C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  clean: {
    configurable: true,
    enumerable: true,
    get: function() { return 1; }
  }
});
$jscomp.global.Object.defineProperties(C.prototype, {
  clean: {
    configurable: true,
    enumerable: true,
    get: function() { return 2; }
  }
});
`)
}

// Java: Es6RewriteClassTest.testStaticGetterFollowedByStaticMethodWithSameName()
func (c *Es6RewriteClassTest) TestStaticGetterFollowedByStaticMethodWithSameName() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`/** @unrestricted */
class C {
  static get clean() { return 1; }
  static clean(x) { return 2; }
}
`,
		`/** @constructor */
let C = function() {};
/** @nocollapse */
C.clean;
$jscomp.global.Object.defineProperties(C, {
  clean: {
    configurable: true,
    enumerable: true,
    get: function() { return 1; }
  }
});
C.clean = function(x) { return 2; };
`)
}

// Java: Es6RewriteClassTest.testStaticGetterFollowedByStaticComputedGetterWithSameName()
func (c *Es6RewriteClassTest) TestStaticGetterFollowedByStaticComputedGetterWithSameName() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`/** @unrestricted */
class C {
  static get clean() { return 1; }
  static get ['clean']() { return 2; }
}
`,
		`/** @constructor */
let C = function() {};
/** @nocollapse */
C.clean;
$jscomp.global.Object.defineProperties(C, {
  clean: {
    configurable: true,
    enumerable: true,
    get: function() { return 1; }
  }
});
$jscomp.global.Object.defineProperty(C, "clean", {
  configurable: true,
  enumerable: true,
  get: function() { return 2; }
});
`)
}

// Java: Es6RewriteClassTest.testComputedPropClass()
func (c *Es6RewriteClassTest) TestComputedPropClass() {
	c.Test(
		"/** @unrestricted */ class C { [foo]() { alert(1); } }",
		`var COMP_FIELD$0 = foo;
/** @constructor */
let C = function() {
};
C.prototype[COMP_FIELD$0] = function() {
  alert(1);
};
`)
}

// Java: Es6RewriteClassTest.testStaticComputedPropClass()
func (c *Es6RewriteClassTest) TestStaticComputedPropClass() {
	c.Test(
		"/** @unrestricted */ class C { static [foo]() { alert(2); } }",
		`var COMP_FIELD$0 = foo;
/** @constructor */
let C = function() {
};
C[COMP_FIELD$0] = function() {
  alert(2);
};
`)
}

// Java: Es6RewriteClassTest.testComputedPropGeneratorMethods()
func (c *Es6RewriteClassTest) TestComputedPropGeneratorMethods() {
	c.Test(
		"/** @unrestricted */ class C { *[foo]() { yield 1; } }",
		`var COMP_FIELD$0 = foo;
/** @constructor */
let C = function() {
};
C.prototype[COMP_FIELD$0] = function*() {
  yield 1;
};
`)
}

// Java: Es6RewriteClassTest.testStaticComputedPropGeneratorMethods()
func (c *Es6RewriteClassTest) TestStaticComputedPropGeneratorMethods() {
	c.Test(
		"/** @unrestricted */ class C { static *[foo]() { yield 2; } }",
		`var COMP_FIELD$0 = foo;
/** @constructor */
let C = function() {
};
C[COMP_FIELD$0] = function*() {
  yield 2;
};
`)
}

// Java: Es6RewriteClassTest.testClassGenerator()
func (c *Es6RewriteClassTest) TestClassGenerator() {
	c.Test(
		"class C { *foo() { yield 1; } }",
		`/** @constructor */
let C = function() {};
C.prototype.foo = function*() { yield 1;};
`)
}

// Java: Es6RewriteClassTest.testClassInsideCast()
func (c *Es6RewriteClassTest) TestClassInsideCast() {
	c.Test(
		"const Foo = /** @type {?} */ (class {});",
		"/** @constructor */ const Foo = function() {};")
}

// Java: Es6RewriteClassTest.testClassWithSuperclassInsideCast()
func (c *Es6RewriteClassTest) TestClassWithSuperclassInsideCast() {
	c.Test(
		`function mixin(baseClass) {
  return /** @type {?} */ (class extends baseClass {
    constructor() {
      super();
    }
  });
}
`,
		`function mixin(baseClass){
  /**
   * @constructor
   */
  const CLASS_DECL$0 = function(){
    return baseClass.call(this) || this
  };
  $jscomp.inherits(CLASS_DECL$0, baseClass);
  return CLASS_DECL$0;
}
`)
}

// Tests that we have reasonable source information on the transpiled nodes.
//
// <p>We do this by comparing the source information of nodes on the transpiled AST with nodes on
// a different parsed-but-not-transpiled AST.
//
// Java: Es6RewriteClassTest.testSimpleClassStatement_hasCorrectSourceInfo()
func (c *Es6RewriteClassTest) TestSimpleClassStatement_hasCorrectSourceInfo() {
	source := "class C { }"
	expected := "/** @constructor */ let C = function() {};"

	asts := c.TestAndReturnAsts(source, expected)
	sourceRoot := asts.sourceRoot
	expectedRoot := asts.expectedRoot

	// Get nodes from the original, pre-transpiled AST
	// `C` from `class C { }`
	sourceCName := c.GetNodeWithName(sourceRoot, "C")
	sourceClass := sourceCName.GetParent()

	// Multistage compilation does not preserve length in the actual output, so remove it from the
	// expected output as well.
	Es6RewriteClassTest_RemoveLengthFromSubtree(sourceRoot)

	// `C` from `let C = function() {};` matches `C` from `class C { }`
	expectedCName := c.GetNodeWithName(expectedRoot, "C")
	AssertNode(expectedCName).MatchesQualifiedName("C").HasEqualSourceInfoTo(sourceCName)

	// function() {}
	expectedFunction := expectedCName.GetOnlyChild()
	AssertNode(expectedFunction).HasToken(rhino.Token_FUNCTION).HasEqualSourceInfoTo(sourceClass)
}

// Java: Es6RewriteClassTest.testClassStatementWithConstructor_hasCorrectSourceInfo()
func (c *Es6RewriteClassTest) TestClassStatementWithConstructor_hasCorrectSourceInfo() {
	source := "class C { constructor() {} }"
	expected := "/** @constructor */ let C = function() {};"

	asts := c.TestAndReturnAsts(source, expected)
	sourceRoot := asts.sourceRoot
	expectedRoot := asts.expectedRoot

	// Get nodes from the original, pre-transpiled AST
	// `C` from `class C {`
	sourceCName := c.GetNodeWithName(sourceRoot, "C")
	// `constructor() {}`
	sourceConstructorFunction := c.GetNodeWithName(sourceRoot, "constructor").GetOnlyChild()
	AssertNode(sourceConstructorFunction).HasToken(rhino.Token_FUNCTION)

	// Multistage compilation does not preserve length in the actual output, so remove it from the
	// expected output as well.
	Es6RewriteClassTest_RemoveLengthFromSubtree(sourceRoot)

	// `C` from `let C = function() {};` matches `C` from `class C {`
	expectedCName := c.GetNodeWithName(expectedRoot, "C")
	AssertNode(expectedCName).MatchesQualifiedName("C").HasEqualSourceInfoTo(sourceCName)

	// `function() {}` matches `constructor() {}`
	expectedFunction := expectedCName.GetOnlyChild()
	AssertNode(expectedFunction).
		HasToken(rhino.Token_FUNCTION).
		HasEqualSourceInfoTo(sourceConstructorFunction)
}

// Java: Es6RewriteClassTest.testClassStatementWithMethod_hasCorrectSourceInfo()
func (c *Es6RewriteClassTest) TestClassStatementWithMethod_hasCorrectSourceInfo() {
	source := "class C { method() {}; }"
	expected := `/** @constructor */
let C = function() {};
C.prototype.method = function() {};
`

	asts := c.TestAndReturnAsts(source, expected)
	sourceRoot := asts.sourceRoot
	expectedRoot := asts.expectedRoot

	// Get nodes from the original, pre-transpiled AST
	// The MEMBER_FUNCTION_DEF for `method`
	sourceMethodMemberDef := c.GetNodeWithName(sourceRoot, "method")
	// The FUNCTION node for `method() {}`
	sourceMethodFunction := sourceMethodMemberDef.GetOnlyChild()

	// Multistage compilation does not preserve length in the actual output, so remove it from the
	// expected output as well.
	Es6RewriteClassTest_RemoveLengthFromSubtree(sourceRoot)

	// `C.prototype.method` has source info matching `method`
	cPrototypeMethod := c.GetNodeWithName(expectedRoot, "method")
	AssertNode(cPrototypeMethod).HasEqualSourceInfoTo(sourceMethodMemberDef)

	// `C.prototype` has source info matching `method`
	cPrototype := cPrototypeMethod.GetFirstChild()
	AssertNode(cPrototype).
		MatchesQualifiedName("C.prototype").
		HasEqualSourceInfoTo(sourceMethodMemberDef)

	// `C` has source info matching `method`
	cName := cPrototype.GetFirstChild()
	AssertNode(cName).MatchesQualifiedName("C").HasEqualSourceInfoTo(sourceMethodMemberDef)

	// `function() {}` has source info matching `method() {}`
	AssertNode(cPrototypeMethod.GetNext()).
		HasToken(rhino.Token_FUNCTION).
		HasEqualSourceInfoTo(sourceMethodFunction)
}

// Java: Es6RewriteClassTest.testSuperCall_hasCorrectSourceInfo()
func (c *Es6RewriteClassTest) TestSuperCall_hasCorrectSourceInfo() {
	source := `class D {}
class C extends D {
  constructor() {
    SUPER: super();
  }
}
`
	expected := `/** @constructor */
let D = function() {};
/** @constructor */
let C = function() {
  SUPER: { D.call(this); }
}
$jscomp.inherits(C, D);
`

	asts := c.TestAndReturnAsts(source, expected)
	sourceRoot := asts.sourceRoot
	expectedRoot := asts.expectedRoot

	// Get nodes from the original, pre-transpiled AST
	sourceSuperCall := c.GetNodeMatchingLabel(sourceRoot, "SUPER").GetOnlyChild()
	AssertNode(sourceSuperCall).HasToken(rhino.Token_CALL)
	sourceSuper := sourceSuperCall.GetFirstChild()
	AssertNode(sourceSuper).HasToken(rhino.Token_SUPER)

	// Multistage compilation does not preserve length in the actual output, so remove it from the
	// expected output as well.
	Es6RewriteClassTest_RemoveLengthFromSubtree(sourceRoot)

	// D.call(this); has the position and length of `super()`
	callNode := c.GetNodeMatchingLabel(expectedRoot, "SUPER").GetOnlyChild().GetOnlyChild()
	AssertNode(callNode).HasToken(rhino.Token_CALL).HasEqualSourceInfoTo(sourceSuperCall)

	// D.call has the position and length of `super`
	callee := callNode.GetFirstChild()
	AssertNode(callee).MatchesQualifiedName("D.call").HasEqualSourceInfoTo(sourceSuper)

	// D has the position and length of `super`
	superDotGReplacement := callee.GetFirstChild()
	AssertNode(superDotGReplacement).MatchesQualifiedName("D").HasEqualSourceInfoTo(sourceSuper)

	// `this` node from `D.call(this)` has the position and length of `super`
	thisNode := callee.GetNext()
	AssertNode(thisNode).HasToken(rhino.Token_THIS).HasEqualSourceInfoTo(sourceSuper)
}

// Java: Es6RewriteClassTest.testComputedSuper_hasCorrectSourceInfo()
func (c *Es6RewriteClassTest) TestComputedSuper_hasCorrectSourceInfo() {
	source := `/** @unrestricted */
class Foo {
  ['m']() { return 1; }
}

/** @unrestricted */
class Bar extends Foo {
  ['m']() {
    RETURN: return super['m']();
  }
}
`
	expected := `/** @constructor */
let Foo = function() {};
Foo.prototype['m'] = function() { return 1; };
/** @constructor */
let Bar = function() { Foo.apply(this, arguments); };
$jscomp.inherits(Bar, Foo);
Bar.prototype['m'] = function () {
  RETURN: { return Foo.prototype['m'].call(this); }
};
`

	asts := c.TestAndReturnAsts(source, expected)
	sourceRoot := asts.sourceRoot
	expectedRoot := asts.expectedRoot

	// Get nodes from the original, pre-transpiled AST
	sourceSuperCall := c.GetNodeMatchingLabel(sourceRoot, "RETURN").GetOnlyChild()
	AssertNode(sourceSuperCall).HasToken(rhino.Token_CALL)
	sourceSuperGet := sourceSuperCall.GetFirstChild()
	AssertNode(sourceSuperGet).HasToken(rhino.Token_GETELEM)
	sourceSuper := sourceSuperGet.GetFirstChild()
	AssertNode(sourceSuper).HasToken(rhino.Token_SUPER)

	// Multistage compilation does not preserve length in the actual output, so remove it from the
	// expected output as well.
	Es6RewriteClassTest_RemoveLengthFromSubtree(sourceRoot)

	// Foo.prototype['m'].call(this) matches source info of `super['m']()`
	callNode := c.GetNodeMatchingLabel(expectedRoot, "RETURN").GetOnlyChild().GetOnlyChild()
	AssertNode(callNode).HasToken(rhino.Token_CALL).HasEqualSourceInfoTo(sourceSuperCall)

	// Foo.prototype['m'].call has the position and length of `super['m']'`
	callee := callNode.GetFirstChild()
	AssertNode(callee).HasToken(rhino.Token_GETPROP).HasEqualSourceInfoTo(sourceSuperGet)

	// Foo.prototype['m'] has the position and length of `super['m']'`
	property := callee.GetFirstChild()
	AssertNode(property).HasToken(rhino.Token_GETELEM).HasEqualSourceInfoTo(sourceSuperGet)

	// Foo.prototype has the position and length of `super`
	prototype := property.GetFirstChild()
	AssertNode(prototype).MatchesQualifiedName("Foo.prototype").HasEqualSourceInfoTo(sourceSuper)

	// Foo has the position and length of `super`
	superDotGReplacement := prototype.GetFirstChild()
	AssertNode(superDotGReplacement).MatchesQualifiedName("Foo").HasEqualSourceInfoTo(sourceSuper)

	// `this` node from `Foo.prototype['m'].call(this)` has position and length of `super['m']`
	thisNode := callee.GetNext()
	AssertNode(thisNode).HasToken(rhino.Token_THIS).HasEqualSourceInfoTo(sourceSuperGet)
}

// Java: Es6RewriteClassTest.testStringKeyGetterAndSetter()
func (c *Es6RewriteClassTest) TestStringKeyGetterAndSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	// Can't add a property to a struct
	c.DisableTypeCheck()
	c.DisableTypeInfoValidation()

	c.Test(
		`class MyClass {
  /** @param {number} v */
  set 'setter'(v) {}
  /** @return {number} */
  get 'getter'() { return 1; }
}
`,
		`/** @constructor */
let MyClass=function(){};

$jscomp.global.Object.defineProperties(
    MyClass.prototype,
    {
      'setter': {
        configurable:true,
        enumerable:true,
        set:function(v){ }
      },
      'getter': {
        configurable:true,
        enumerable:true,
        get:function(){ return 1; }
      },
    })
`)
}

// Java: Es6RewriteClassTest.testNumericGetterAndSetter()
func (c *Es6RewriteClassTest) TestNumericGetterAndSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	// Can't add a property to a struct
	c.DisableTypeCheck()
	c.DisableTypeInfoValidation()

	c.Test(
		`class MyClass {
  /** @param {number} v */
  set 1(v) {}
  /** @return {number} */
  get 1.5() { return 1; }
}
`,
		`/** @constructor */
let MyClass=function(){};
$jscomp.global.Object.defineProperties(
    MyClass.prototype,
    {
      1: {
        configurable:true,
        enumerable:true,
        set:function(v){ }
      },
      1.5: {
        configurable:true,
        enumerable:true,
        get:function(){ return 1; }
      },
    })
`)
}

// GetNodeMatchingLabel returns the first node (preorder) in the given AST
// labeled as `label`
//
// Java: Es6RewriteClassTest.getNodeMatchingLabel(Node, String)
func (c *Es6RewriteClassTest) GetNodeMatchingLabel(root *rhino.Node, label string) *rhino.Node {
	if root.IsLabel() && root.GetFirstChild().GetString() == label {
		// "FOO: return foo;"
		//   becomes
		// LABEL
		//   LABEL_NAME FOO
		//   RETURN
		//   NAME foo
		return root.GetSecondChild()
	}
	for child := root.GetFirstChild(); child != nil; child = child.GetNext() {
		result := c.GetNodeMatchingLabel(child, label)
		if result != nil {
			return result
		}
	}
	return nil
}

// GetNodeWithName returns the first node (preorder) in the given AST that
// matches the given qualified name
//
// Java: Es6RewriteClassTest.getNodeWithName(Node, String)
func (c *Es6RewriteClassTest) GetNodeWithName(root *rhino.Node, name string) *rhino.Node {
	switch root.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_NAME, rhino.Token_STRINGLIT, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_STRING_KEY:
		if root.GetString() == name {
			return root
		}
	default:
	}

	for child := root.GetFirstChild(); child != nil; child = child.GetNext() {
		result := c.GetNodeWithName(child, name)
		if result != nil {
			return result
		}
	}
	return nil
}

// TestAndReturnAsts tests that the transpiled source matches the expected
// code and then returns the AST roots for the source and the transpiled
// source respectively.
//
// Java: Es6RewriteClassTest.testAndReturnAsts(String, String)
func (c *Es6RewriteClassTest) TestAndReturnAsts(source string, expected string) *es6RewriteClassTest_AstPair {
	originalRoot := c.CreateCompiler().ParseTestCode(source)

	c.Test(source, expected)

	return newEs6RewriteClassTest_AstPair(originalRoot, c.GetLastCompiler().GetJsRoot())
}

// Java: private static class Es6RewriteClassTest.AstPair
type es6RewriteClassTest_AstPair struct {
	sourceRoot *rhino.Node

	expectedRoot *rhino.Node
}

// Java: Es6RewriteClassTest.AstPair(Node, Node)
func newEs6RewriteClassTest_AstPair(sourceRoot *rhino.Node, expectedRoot *rhino.Node) *es6RewriteClassTest_AstPair {
	return &es6RewriteClassTest_AstPair{sourceRoot: sourceRoot, expectedRoot: expectedRoot}
}

// Es6RewriteClassTest_RemoveLengthFromSubtree sets n.getLength() to 0 for
// every node in the given subtree
//
// Java: Es6RewriteClassTest.removeLengthFromSubtree(Node)
func Es6RewriteClassTest_RemoveLengthFromSubtree(root *rhino.Node) {
	root.SetLength(0)
	for child := root.GetFirstChild(); child != nil; child = child.GetNext() {
		Es6RewriteClassTest_RemoveLengthFromSubtree(child)
	}
}

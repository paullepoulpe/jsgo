// Port of com/google/javascript/jscomp/CheckAccessControlsOldSyntaxTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/CheckAccessControlsOldSyntaxTest.java)

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
)

// CheckAccessControlsOldSyntaxTest tests CheckAccessControls.
//
// This file has a fork, CheckAccessControlsTest, because nearly all cases
// using legacy ES5-sytle classes should also be tested with the more modern
// `class` keyword. If a case using `@constructor`, `@interface`, or
// `@record` is added to this suite, a similar case should be added there
// under the same name using `class`.
type CheckAccessControlsOldSyntaxTest struct {
	CompilerTestCaseBase
}

// Java: CheckAccessControlsOldSyntaxTest()
func NewCheckAccessControlsOldSyntaxTest() *CheckAccessControlsOldSyntaxTest {
	c := &CheckAccessControlsOldSyntaxTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTypeTestCase_DEFAULT_EXTERNS)
	return c
}

func TestCheckAccessControlsOldSyntaxTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckAccessControlsOldSyntaxTest() })
}

// Java: CheckAccessControlsOldSyntaxTest.setUp()
func (c *CheckAccessControlsOldSyntaxTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableParseTypeInfo()
	c.EnableClosurePass()
	c.EnableRewriteClosureCode()
	c.EnableCreateModuleMap()
}

// Java: CheckAccessControlsOldSyntaxTest.getProcessor(Compiler)
func (c *CheckAccessControlsOldSyntaxTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckAccessControls(compiler)
}

// Java: CheckAccessControlsOldSyntaxTest.getOptions()
func (c *CheckAccessControlsOldSyntaxTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_ACCESS_CONTROLS, CheckLevel_ERROR)
	options.SetWarningLevel(DiagnosticGroups_DEPRECATED, CheckLevel_ERROR)
	options.SetWarningLevel(DiagnosticGroups_CONSTANT_PROPERTY, CheckLevel_ERROR)
	return options
}

// Java: CheckAccessControlsOldSyntaxTest.testDepName(String, String)
func (c *CheckAccessControlsOldSyntaxTest) TestDepName(js string, errorMessage string) {
	c.TestError(js, CheckAccessControls_DEPRECATED_NAME_REASON, errorMessage)
}

// Java: CheckAccessControlsOldSyntaxTest.testDepProp(String, String)
func (c *CheckAccessControlsOldSyntaxTest) TestDepProp(js string, errorMessage string) {
	c.TestError(js, CheckAccessControls_DEPRECATED_PROP_REASON, errorMessage)
}

// Java: CheckAccessControlsOldSyntaxTest.testDepClass(String, String)
func (c *CheckAccessControlsOldSyntaxTest) TestDepClass(js string, errorMessage string) {
	c.TestError(js, CheckAccessControls_DEPRECATED_CLASS_REASON, errorMessage)
}

// Java: CheckAccessControlsOldSyntaxTest.testDeprecatedFunctionNoReason()
func (c *CheckAccessControlsOldSyntaxTest) TestDeprecatedFunctionNoReason() {
	c.TestError("/** @deprecated */ function f() {} function g() { f(); }", CheckAccessControls_DEPRECATED_NAME)
}

// Java: CheckAccessControlsOldSyntaxTest.testDeprecatedFunctionNoReason_googModule()
func (c *CheckAccessControlsOldSyntaxTest) TestDeprecatedFunctionNoReason_googModule() {
	c.DisableRewriteClosureCode()
	// Remove this line once closure rewriting is after typechecking
	c.TestError(
		"goog.module('m'); /** @deprecated */ function f() {} function g() { f(); }",
		CheckAccessControls_DEPRECATED_NAME)
}

// Java: CheckAccessControlsOldSyntaxTest.testDeprecatedFunction()
func (c *CheckAccessControlsOldSyntaxTest) TestDeprecatedFunction() {
	c.TestDepName(
		"/** @deprecated Some Reason */ function f() {} function g() { f(); }",
		"Variable f has been deprecated: Some Reason")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningOnDeprecatedConstVariable()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningOnDeprecatedConstVariable() {
	c.TestDepName(
		"/** @deprecated Another reason */ var f = 4; function g() { alert(f); }",
		"Variable f has been deprecated: Another reason")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningOnDeprecatedConstVariableWithConst()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningOnDeprecatedConstVariableWithConst() {
	c.TestDepName(
		"/** @deprecated Another reason */ const f = 4; function g() { alert(f); }",
		"Variable f has been deprecated: Another reason")
}

// Java: CheckAccessControlsOldSyntaxTest.testThatNumbersArentDeprecated()
func (c *CheckAccessControlsOldSyntaxTest) TestThatNumbersArentDeprecated() {
	c.TestSame("/** @deprecated */ var f = 4; var h = 3; function g() { alert(h); }")
}

// Java: CheckAccessControlsOldSyntaxTest.testDeprecatedFunctionVariable()
func (c *CheckAccessControlsOldSyntaxTest) TestDeprecatedFunctionVariable() {
	c.TestDepName(
		"/** @deprecated I like g... */ var f = function() {}; function g() { f(); }",
		"Variable f has been deprecated: I like g...")
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningInGlobalScope()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningInGlobalScope() {
	c.TestSame(
		`goog.makeSingleton = function(x) {};
/** @deprecated */ function f() {} goog.makeSingleton(f);
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningInGlobalScopeForCall()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningInGlobalScopeForCall() {
	c.TestDepName(
		"/** @deprecated Some global scope */ function f() {} f();",
		"Variable f has been deprecated: Some global scope")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningInGoogModuleScopeForCall()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningInGoogModuleScopeForCall() {
	c.DisableRewriteClosureCode()
	// Remove this line once Closure rewriting is after typechecking.
	c.TestDepName(
		"goog.module('m'); /** @deprecated Some module scope */ function f() {} f();",
		"Variable f has been deprecated: Some module scope")
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningInGoogModuleScopeWithFileoverviewForCall()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningInGoogModuleScopeWithFileoverviewForCall() {
	c.DisableRewriteClosureCode()
	// Remove this line once Closure rewriting is after typechecking.
	c.TestDepName(
		`/** @deprecated @fileoverview */ goog.module('m'); /** @deprecated Some module scope */
 function f() {} f();
`,
		"Variable f has been deprecated: Some module scope")
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningInDeprecatedFunction()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningInDeprecatedFunction() {
	c.TestSame("/** @deprecated */ function f() {} /** @deprecated */ function g() { f(); }")
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningInDeprecatedMethod()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningInDeprecatedMethod() {
	c.TestSame("/** @deprecated */ function f() {} var obj = {/** @deprecated */ g() { f(); }};")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningInNormalMethod()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningInNormalMethod() {
	c.TestDepName(
		"/** @deprecated Msg */ function f() {} var obj = {g() { f(); }};",
		"Variable f has been deprecated: Msg")
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningInDeprecatedComputedMethod()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningInDeprecatedComputedMethod() {
	c.TestSame("/** @deprecated */ function f() {} var obj = {/** @deprecated */ ['g']() { f(); }};")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningInNormalComputedMethod()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningInNormalComputedMethod() {
	c.TestDepName(
		"/** @deprecated Msg */ function f() {} var obj = {['g']() { f(); }};",
		"Variable f has been deprecated: Msg")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningInNormalClass()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningInNormalClass() {
	c.TestDepName(
		`/** @deprecated FooBar */ function f() {}
/** @constructor */  var Foo = function() {};
Foo.prototype.bar = function() { f(); }
`,
		"Variable f has been deprecated: FooBar")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForProperty1()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForProperty1() {
	c.TestDepProp(
		`/** @constructor */ function Foo() {}
/** @deprecated A property is bad */ Foo.prototype.bar = 3;
Foo.prototype.baz = function() { alert((new Foo()).bar); };
`,
		"Property bar of type Foo has been deprecated: A property is bad")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForProperty2()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForProperty2() {
	c.TestDepProp(
		`/** @constructor */ function Foo() {}
/** @deprecated Zee prop, it is deprecated! */ Foo.prototype.bar = 3;
Foo.prototype.baz = function() { alert(this.bar); };
`,
		"Property bar of type Foo has been deprecated: Zee prop, it is deprecated!")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForDeprecatedClass()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForDeprecatedClass() {
	c.TestDepClass(
		`/** @constructor 
* @deprecated Use the class 'Bar' */ function Foo() {}
function f() { new Foo(); }
`,
		"Class Foo has been deprecated: Use the class 'Bar'")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForDeprecatedClassNoReason()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForDeprecatedClassNoReason() {
	c.TestError(
		`/** @constructor 
* @deprecated */ function Foo() {}
function f() { new Foo(); }
`,
		CheckAccessControls_DEPRECATED_CLASS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningForDeprecatedClassInstance()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningForDeprecatedClassInstance() {
	c.TestSame(
		`/** @constructor 
* @deprecated */ function Foo() {}
/** @param {Foo} x */ function f(x) { return x; }
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningForDeprecatedSuperClass()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningForDeprecatedSuperClass() {
	c.TestNoWarning(
		`/** @constructor @deprecated Superclass to the rescue! */
function Foo() {}
/** @constructor * @extends {Foo} */
function SubFoo() {}
function f() { new SubFoo(); }
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningForDeprecatedSuperClassOnNamespace()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningForDeprecatedSuperClassOnNamespace() {
	c.TestNoWarning(
		`/**
 * @constructor
 * @deprecated Its only weakness is Kryptoclass
*/
 function Foo() {}
/** @const */ var namespace = {};
/** @constructor
 * @extends {Foo} */
namespace.SubFoo = function() {};
function f() { new namespace.SubFoo(); }
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForPrototypeProperty()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForPrototypeProperty() {
	js := `/** @constructor */ function Foo() {}
/** @deprecated It is now in production, use that model... */
Foo.prototype.bar = 3;
Foo.prototype.baz = function() { alert(Foo.prototype.bar); };
`
	c.TestDepProp(
		js,
		`Property bar of type Foo.prototype has been deprecated: It is now in production, use that model...
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningForNumbers()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningForNumbers() {
	c.TestSame(
		`/** @constructor */ function Foo() {}
/** @deprecated */ Foo.prototype.bar = 3;
Foo.prototype.baz = function() { alert(3); };
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForMethod1()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForMethod1() {
	c.TestDepProp(
		`/** @constructor */ function Foo() {}
/** @deprecated There is a madness to this method */
Foo.prototype.bar = function() {};
Foo.prototype.baz = function() { this.bar(); };
`,
		"Property bar of type Foo has been deprecated: There is a madness to this method")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForMethod2()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForMethod2() {
	c.TestDepProp(
		`/** @constructor */ function Foo() {}
/** @deprecated Stop the ringing! */ Foo.prototype.bar;
Foo.prototype.baz = function() { this.bar(); };
`,
		"Property bar of type Foo has been deprecated: Stop the ringing!")
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningInDeprecatedClass()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningInDeprecatedClass() {
	c.TestSame(
		`/** @deprecated */ function f() {}
/** @constructor 
 * @deprecated */
var Foo = function() {};
Foo.prototype.bar = function() { f(); }
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningOnDeclaration()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningOnDeclaration() {
	c.TestSame(
		`/** @constructor */ function F() {
  /**
   * @type {number}
   * @deprecated Use something else.
   */
  this.code;
}
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningInDeprecatedClass2()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningInDeprecatedClass2() {
	c.TestSame(
		`/** @deprecated */ function f() {}
/** @constructor 
 * @deprecated */
var Foo = function() {};
Foo.bar = function() { f(); }
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningInDeprecatedStaticMethod()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningInDeprecatedStaticMethod() {
	c.TestSame(
		`/** @deprecated */ function f() {}
/** @constructor */
var Foo = function() {};
/** @deprecated */ Foo.bar = function() { f(); }
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningInStaticMethod()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningInStaticMethod() {
	c.TestDepName(
		`/** @deprecated crazy! */ function f() {}
/** @constructor */
var Foo = function() {};
Foo.bar = function() { f(); }
`,
		"Variable f has been deprecated: crazy!")
}

// Java: CheckAccessControlsOldSyntaxTest.testDeprecatedObjLitKey()
func (c *CheckAccessControlsOldSyntaxTest) TestDeprecatedObjLitKey() {
	c.TestDepProp(
		`/** @const */ var f = {};
/** @deprecated It is literally not used anymore */ f.foo = 3;
function g() { return f.foo; }
`,
		"Property foo of type f has been deprecated: It is literally not used anymore")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForSubclassMethod()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForSubclassMethod() {
	c.TestDepProp(
		`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() {};
/** @constructor 
 * @extends {Foo} */ function SubFoo() {}
/** @deprecated I have a parent class! */ SubFoo.prototype.bar = function() {};
function f() { (new SubFoo()).bar(); };
`,
		"Property bar of type SubFoo has been deprecated: I have a parent class!")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForSuperClassWithDeprecatedSubclassMethod()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForSuperClassWithDeprecatedSubclassMethod() {
	c.TestSame(
		`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() {};
/** @constructor 
 * @extends {Foo} */ function SubFoo() {}
/** @deprecated 
 * @override */ SubFoo.prototype.bar =
function() {};
function f() { (new Foo()).bar(); };
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForSuperclassMethod()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForSuperclassMethod() {
	c.TestDepProp(
		`/** @constructor */ function Foo() {}
/** @deprecated I have a child class! */ Foo.prototype.bar = function() {};
/** @constructor 
 * @extends {Foo} */ function SubFoo() {}
SubFoo.prototype.bar = function() {};
function f() { (new SubFoo()).bar(); };
`,
		"Property bar of type SubFoo has been deprecated: I have a child class!")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForSuperclassMethod2()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForSuperclassMethod2() {
	c.TestDepProp(
		`/** @constructor */ function Foo() {}
/** @deprecated I have another child class... 
* @protected */
Foo.prototype.bar = function() {};
/** @constructor 
 * @extends {Foo} */ function SubFoo() {}
/** @protected */SubFoo.prototype.bar = function() {};
function f() { (new SubFoo()).bar(); };
`,
		"Property bar of type SubFoo has been deprecated: I have another child class...")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForFunctionProp()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForFunctionProp() {
	c.TestDepProp(
		`/** @deprecated Don't call me... */ Function.prototype.deprecatedMethod = function() {};
(function() {}).deprecatedMethod();
`,
		"Property deprecatedMethod of type function has been deprecated: Don't call me...")
}

// Java: CheckAccessControlsOldSyntaxTest.testWarningForDeprecatedClassInGlobalScope()
func (c *CheckAccessControlsOldSyntaxTest) TestWarningForDeprecatedClassInGlobalScope() {
	c.TestDepClass(
		`/** @constructor 
 * @deprecated I'm a very worldly object! */ var Foo = function() {};
new Foo();
`,
		"Class Foo has been deprecated: I'm a very worldly object!")
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningForPrototypeCopying()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningForPrototypeCopying() {
	c.TestSame(
		`/** @constructor */ var Foo = function() {};
Foo.prototype.bar = function() {};
/** @deprecated */ Foo.prototype.baz = Foo.prototype.bar;
(new Foo()).bar();
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoWarningOnDeprecatedPrototype()
func (c *CheckAccessControlsOldSyntaxTest) TestNoWarningOnDeprecatedPrototype() {
	// This used to cause an NPE.
	c.TestSame(
		`/** @constructor */ var Foo = function() {};
/** @deprecated */ Foo.prototype = {};
Foo.prototype.bar = function() {};
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testPrivateAccessForNames()
func (c *CheckAccessControlsOldSyntaxTest) TestPrivateAccessForNames() {
	c.TestSame("/** @private */ function foo_() {}; foo_();")
	c.TestError(c.Srcs("/** @private */ function foo_() {};", "foo_();"), CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testPrivateAccessForNames2()
func (c *CheckAccessControlsOldSyntaxTest) TestPrivateAccessForNames2() {
	// Not private by convention
	c.TestSame("function foo_() {}; foo_();")
	c.TestSame(c.Srcs("function foo_() {};", "foo_();"))
}

// Java: CheckAccessControlsOldSyntaxTest.testPrivateAccessForProperties1()
func (c *CheckAccessControlsOldSyntaxTest) TestPrivateAccessForProperties1() {
	c.TestSame(
		`/** @constructor */ function Foo() {}
/** @private */ Foo.prototype.bar_ = function() {};
Foo.prototype.baz = function() { this.bar_(); }; (new Foo).bar_();
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testPrivateAccessForProperties2()
func (c *CheckAccessControlsOldSyntaxTest) TestPrivateAccessForProperties2() {
	c.TestSame(
		c.Srcs(
			"/** @constructor */ function Foo() {}",
			`/** @private */ Foo.prototype.bar_ = function() {};
Foo.prototype.baz = function() { this.bar_(); }; (new Foo).bar_();
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testPrivateAccessForProperties3()
func (c *CheckAccessControlsOldSyntaxTest) TestPrivateAccessForProperties3() {
	// Even though baz is "part of the Foo class" the access is disallowed since it's
	// not in the same file.
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @private */ Foo.prototype.bar_ = function() {}; (new Foo).bar_();
`,
			"Foo.prototype.baz = function() { this.bar_(); };"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testPrivateAccessForProperties4()
func (c *CheckAccessControlsOldSyntaxTest) TestPrivateAccessForProperties4() {
	c.TestSame(
		`/** @constructor */ function Foo() {}
/** @private */ Foo.prototype.bar_ = function() {};
Foo.prototype['baz'] = function() { (new Foo()).bar_(); };
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testPrivateAccessForProperties5()
func (c *CheckAccessControlsOldSyntaxTest) TestPrivateAccessForProperties5() {
	c.TestError(
		c.Srcs(
			`/** @constructor */
function Parent () {
  /** @private */
  this.prop = 'foo';
};
`,
			`/**
 * @constructor
 * @extends {Parent}
 */
function Child() {
  this.prop = 'asdf';
}
Child.prototype = new Parent();
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS).
			WithMessage("Access to private property prop of Parent not allowed here."))
}

// Java: CheckAccessControlsOldSyntaxTest.testBug133902968()
func (c *CheckAccessControlsOldSyntaxTest) TestBug133902968() {
	c.Test(
		c.Externs(
			`/** @interface */
function CSSProperties() {}
/** @type {string} */
CSSProperties.prototype.fontStyle;

/** @struct @interface @extends {CSSProperties} */
function CSSStyleDeclaration() {}

function alert(s) {}
`),
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("apackage", "afile.js"),
				`/**
 * @fileoverview
 * @package
 */

/** @param {!CSSStyleDeclaration} style */
function f(style) {
  style.fontStyle = 'normal';
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("anotherpackage", "anotherfile.js"),
				`/** @param {!CSSStyleDeclaration} style */
function g(style) {
  alert(style.fontStyle);
}
`)))
}

// Java: CheckAccessControlsOldSyntaxTest.testPrivatePropAccess_inSameFile_throughDestructuring()
func (c *CheckAccessControlsOldSyntaxTest) TestPrivatePropAccess_inSameFile_throughDestructuring() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { }

/** @private */
Foo.prototype.bar_ = function() { };

function f(/** !Foo */ x) {
  const {bar_: bar} = x;
}
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testPrivateAccessForProperties6()
func (c *CheckAccessControlsOldSyntaxTest) TestPrivateAccessForProperties6() {
	c.TestError(
		c.Srcs(
			`goog.provide('x.y.z.Parent');

/** @constructor */
x.y.z.Parent = function() {
  /** @private */
  this.prop = 'foo';
};
`,
			`goog.require('x.y.z.Parent');

/**
 * @constructor
 * @extends {x.y.z.Parent}
 */
function Child() {
  this.prop = 'asdf';
}
Child.prototype = new x.y.z.Parent();
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS).
			WithMessage("Access to private property prop of x.y.z.Parent not allowed here."))
}

// Java: CheckAccessControlsOldSyntaxTest.testPrivateAccess_googModule()
func (c *CheckAccessControlsOldSyntaxTest) TestPrivateAccess_googModule() {
	c.TestError(
		c.Srcs(
			`goog.module('example.One');
/** @constructor */ function One() {}
/** @private */ One.prototype.m = function() {};
exports = One;
`,
			`goog.module('example.two');
var One = goog.require('example.One');
(new One()).m();
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS).
			WithMessage("Access to private property m of One not allowed here."))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties1()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties1() {
	c.TestError(
		c.Srcs(
			"/** @constructor */ function Foo() {} (new Foo).bar_();",
			`/** @private */ Foo.prototype.bar_ = function() {};
Foo.prototype.baz = function() { this.bar_(); };
`),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties2()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties2() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @private */ Foo.prototype.bar_ = function() {};
Foo.prototype.baz = function() { this.bar_(); };
`,
			"(new Foo).bar_();"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties3()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties3() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @private */ Foo.prototype.bar_ = function() {};
`,
			"/** @constructor */ function OtherFoo() { (new Foo).bar_(); }"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties4()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties4() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @private */ Foo.prototype.bar_ = function() {};
`,
			`/** @constructor 
 * @extends {Foo} */
function SubFoo() { this.bar_(); }
`),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties5()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties5() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @private */ Foo.prototype.bar_ = function() {};
`,
			`/** @constructor 
 * @extends {Foo} */
function SubFoo() {};
SubFoo.prototype.baz = function() { this.bar_(); }
`),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties6()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties6() {
	// Overriding a private property with a non-private property
	// in a different file causes problems.
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() {}

/** @private */
Foo.prototype.bar_ = function() {};
`,
			`/**
 * @constructor
 * @extends {Foo}
 */
function SubFoo() {};

SubFoo.prototype.bar_ = function() {};
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties6a()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties6a() {
	// Same as above, except with namespaced constructors
	c.TestError(
		c.Srcs(
			`/** @const */ var ns = {};
/** @constructor */ ns.Foo = function() {};
/** @private */ ns.Foo.prototype.bar_ = function() {};
`,
			`/** @constructor 
 * @extends {ns.Foo} */
ns.SubFoo = function() {};
ns.SubFoo.prototype.bar_ = function() {};
`),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties7()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties7() {
	// It's OK to override a private property with a non-private property
	// in the same file, but you'll get yelled at when you try to use it.
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @private */ Foo.prototype.bar_ = function() {};
/** @constructor 
 * @extends {Foo} */
function SubFoo() {};
SubFoo.prototype.bar_ = function() {};
`,
			"SubFoo.prototype.baz = function() { this.bar_(); }"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties8()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties8() {
	c.TestError(
		c.Srcs(
			"/** @constructor */ function Foo() { /** @private */ this.bar_ = 3; }",
			`/** @constructor 
 * @extends {Foo} */
function SubFoo() { /** @private */ this.bar_ = 3; };
`),
		CheckAccessControls_PRIVATE_OVERRIDE)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties9()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties9() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
Foo.prototype = {
/** @private */ bar_: 3
}
`,
			"new Foo().bar_;"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties10()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties10() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
Foo.prototype = {
/** @private */ bar_: function() {}
}
`,
			"new Foo().bar_();"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties11()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties11() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
Foo.prototype = {
/** @private */ get bar_() { return 1; }
}
`,
			"var a = new Foo().bar_;"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForProperties12()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForProperties12() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
Foo.prototype = {
/** @private */ set bar_(x) { this.barValue = x; }
}
`,
			"new Foo().bar_ = 1;"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivatePropAccess_inDifferentFile_throughDestructuring()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivatePropAccess_inDifferentFile_throughDestructuring() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { }

/** @private */
Foo.prototype.bar_ = function() { };
`,
			`function f(/** !Foo */ x) {
  const {bar_: bar} = x;
}
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivateAccessForNamespaces()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivateAccessForNamespaces() {
	c.TestError(
		c.Srcs(
			`/** @const */ var foo = {};
/** @private */ foo.bar_ = function() {};
`,
			"foo.bar_();"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties1()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties1() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @protected */ Foo.prototype.bar = function() {};
(new Foo).bar();
`,
			"Foo.prototype.baz = function() { this.bar(); };"))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties2()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties2() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @protected */ Foo.prototype.bar = function() {};
(new Foo).bar();
`,
			`/** @constructor 
 * @extends {Foo} */
function SubFoo() { this.bar(); }
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties3()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties3() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @protected */ Foo.prototype.bar = function() {};
(new Foo).bar();
`,
			`/** @constructor 
 * @extends {Foo} */
function SubFoo() { }
SubFoo.baz = function() { (new Foo).bar(); }
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties4()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties4() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @protected */ Foo.bar = function() {};
`,
			`/** @constructor 
 * @extends {Foo} */
function SubFoo() { Foo.bar(); }
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties5()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties5() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @protected */ Foo.prototype.bar = function() {};
(new Foo).bar();
`,
			`/** @constructor 
 * @extends {Foo} */
var SubFoo = function() { this.bar(); }
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties6()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties6() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ goog.Foo = function() {};
/** @protected */ goog.Foo.prototype.bar = function() {};
`,
			`/** @constructor 
 * @extends {goog.Foo} */
goog.SubFoo = function() { this.bar(); };
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties7()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties7() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ var Foo = function() {};
Foo.prototype = { /** @protected */ bar: function() {} }
`,
			`/** @constructor 
 * @extends {Foo} */
var SubFoo = function() { this.bar(); };
SubFoo.prototype = { moo: function() { this.bar(); }};
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties8()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties8() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ var Foo = function() {};
Foo.prototype = { /** @protected */ bar: function() {} }
`,
			`/** @constructor 
 * @extends {Foo} */
var SubFoo = function() {};
SubFoo.prototype = { get moo() { this.bar(); }};
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties9()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties9() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ var Foo = function() {};
Foo.prototype = { /** @protected */ bar: function() {} }
`,
			`/** @constructor 
 * @extends {Foo} */
var SubFoo = function() {};
SubFoo.prototype = { set moo(val) { this.x = this.bar(); }};
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties10()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties10() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.js",
				`/** @constructor */ var Foo = function() {};
/** @protected */ Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				"sub_foo.js",
				`/** @constructor @extends {Foo} */
var SubFoo = function() {};
(function() {
SubFoo.prototype.baz = function() { this.bar(); }
})();
`)))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties11()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties11() {
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				"foo.js",
				`goog.provide('Foo');
/** @interface */ Foo = function() {};
/** @protected */ Foo.prop = {};
`),
			SourceFile_FromCode(
				"bar.js",
				`goog.require('Foo');
/** @constructor @implements {Foo} */
function Bar() { Foo.prop; };
`)))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties12()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties12() {
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				"a.js",
				`goog.provide('A');
/** @constructor */
var A = function() {
  /**
   * @type {?String}
   * @protected
   */
  this.prop;
}
`),
			SourceFile_FromCode(
				"b.js",
				`goog.require('A');
/**
 * @constructor
 * @extends {A}
 */
var B = function() {
  this.prop.length;
  this.prop.length;
};
`)))
}

// FYI: Java warns for the b1.method access in c.js.
// Instead of following that in NTI, we chose to follow the behavior of
// the old JSCompiler type checker, to make migration easier.
// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties13()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties13() {
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				"a.js",
				`goog.provide('A');

/** @constructor */
var A = function() {}

/** @protected */
A.prototype.method = function() {};
`),
			SourceFile_FromCode(
				"b1.js",
				`goog.require('A');
goog.provide('B1');

/**
 * @constructor
 * @extends {A}
 */
var B1 = function() {};

/** @override */
B1.prototype.method = function() {};
`),
			SourceFile_FromCode(
				"b2.js",
				`goog.require('A');
goog.provide('B2');

/**
 * @constructor
 * @extends {A}
 */
var B2 = function() {};

/** @override */
B2.prototype.method = function() {};
`),
			SourceFile_FromCode(
				"c.js",
				`goog.require('B1');
goog.require('B2');

/**
 * @param {!B1} b1
 * @constructor
 * @extends {B2}
 */
var C = function(b1) {
  var x = b1.method();
};
`)))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties14()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties14() {
	// access in member function
	c.Test(
		c.Srcs(
			`/** @constructor */ var Foo = function() {};
Foo.prototype = { /** @protected */ bar: function() {} }
`,
			`/** @constructor @extends {Foo} */
var OtherFoo = function() { this.bar(); };
OtherFoo.prototype = { moo() { new Foo().bar(); }};
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties15()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties15() {
	// access in computed member function
	c.Test(
		c.Srcs(
			`/** @constructor */ var Foo = function() {};
Foo.prototype = { /** @protected */ bar: function() {} }
`,
			`/** @constructor @extends {Foo} */
var OtherFoo = function() { this['bar'](); };
OtherFoo.prototype = { ['bar']() { new Foo().bar(); }};
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessForProperties16()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessForProperties16() {
	// access in nested arrow function
	c.Test(
		c.Srcs(
			`/** @constructor */ var Foo = function() {};
/** @protected */ Foo.prototype.bar = function() {};
`,
			`/** @constructor @extends {Foo} */
var OtherFoo = function() { var f = () => this.bar(); };
OtherFoo.prototype.baz = function() { return () => this.bar(); };
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedPropAccess_inDifferentFile_inSubclass_throughDestructuring()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedPropAccess_inDifferentFile_inSubclass_throughDestructuring() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { }

/** @protected */
Foo.prototype.bar = function() { }
`,
			`/** @constructor @extends {Foo} */
function SubFoo() { }

SubFoo.prototype.method = function(/** !Foo */ x) {
  const {bar: bar} = x;
};
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccess_forOverriddenProperty_elsewhereInSubclassFile()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccess_forOverriddenProperty_elsewhereInSubclassFile() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { }

/** @protected */
Foo.prototype.bar = function() { };
`,
			`/**
 * @constructor
 * @extends {Foo}
 */
function Bar() { }

/** @override */
Bar.prototype.bar = function() { };

// TODO(b/113705099): This should be legal.
(new Bar()).bar();
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessThroughNestedFunction()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessThroughNestedFunction() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { }

/** @protected */
Foo.prototype.bar = function() { };
`,
			`/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {
  function f(/** !Foo */ foo) {
    foo.bar();
  }
}
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessThroughNestedEs5Class()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessThroughNestedEs5Class() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { }

/** @protected */
Foo.prototype.bar = function() { };
`,
			`/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {
  /** @constructor */
  var Nested = function() { }

  /** @param {!Foo} foo */
  Nested.prototype.qux = function(foo) {
    foo.bar();
  }
}
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testProtectedAccessThroughNestedEs6Class()
func (c *CheckAccessControlsOldSyntaxTest) TestProtectedAccessThroughNestedEs6Class() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { }

/** @protected */
Foo.prototype.bar = function() { };
`,
			`/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {
  class Nested {
    /** @param {!Foo} foo */
    qux(foo) {
      foo.bar();
    }
  }
}
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccessForProperties1()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccessForProperties1() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @protected */ Foo.prototype.bar = function() {};
`,
			"(new Foo).bar();"),
		CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccessForProperties2()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccessForProperties2() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @protected */ Foo.prototype.bar = function() {};
`,
			"/** @constructor */ function OtherFoo() { (new Foo).bar(); }"),
		CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccessForProperties3()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccessForProperties3() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @constructor 
 * @extends {Foo} */
function SubFoo() {}
/** @protected */ SubFoo.prototype.bar = function() {};
`,
			`/** @constructor 
 * @extends {Foo} */
function SubberFoo() { (new SubFoo).bar(); }
`),
		CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccessForProperties4()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccessForProperties4() {
	c.TestError(
		c.Srcs(
			"/** @constructor */ function Foo() { (new SubFoo).bar(); } ",
			`/** @constructor 
 * @extends {Foo} */
function SubFoo() {}
/** @protected */ SubFoo.prototype.bar = function() {};
`),
		CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccessForProperties5()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccessForProperties5() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ goog.Foo = function() {};
/** @protected */ goog.Foo.prototype.bar = function() {};
`,
			`/** @constructor */
goog.NotASubFoo = function() { (new goog.Foo).bar(); };
`),
		CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccessForProperties6()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccessForProperties6() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
Foo.prototype = {
/** @protected */ bar: 3
}
`,
			"new Foo().bar;"),
		CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccessForProperties7()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccessForProperties7() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
Foo.prototype = {
/** @protected */ bar: function() {}
}
`,
			"new Foo().bar();"),
		CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccessForProperties8()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccessForProperties8() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ var Foo = function() {};
Foo.prototype = { /** @protected */ bar: function() {} }
`,
			`/** @constructor */
var OtherFoo = function() { this.bar(); };
OtherFoo.prototype = { moo() { new Foo().bar(); }};
`),
		CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccessForProperties9()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccessForProperties9() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ var Foo = function() {};
Foo.prototype = { /** @protected */ bar: function() {} }
`,
			`/** @constructor */
var OtherFoo = function() { this['bar'](); };
OtherFoo.prototype = { ['bar']() { new Foo().bar(); }};
`),
		CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedPropAccess_inDifferentFile_throughDestructuring()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedPropAccess_inDifferentFile_throughDestructuring() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { }

/** @protected */
Foo.prototype.bar = function() { }
`,
			`function f(/** !Foo */ x) {
  const {bar: bar} = x;
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccess_forInheritedProperty_elsewhereInSubclassFile()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccess_forInheritedProperty_elsewhereInSubclassFile() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { }

/** @protected */
Foo.prototype.bar = function() { };
`,
			`/**
 * @constructor
 * @extends {Foo}
 */
function Bar() { }

(new Bar()).bar();
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoProtectedAccessForPropertiesWithNoRhs()
func (c *CheckAccessControlsOldSyntaxTest) TestNoProtectedAccessForPropertiesWithNoRhs() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @protected */ Foo.prototype.x;
`,
			`/** @constructor @extends {Foo} */ function Bar() {}
/** @protected */ Bar.prototype.x;
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackagePrivateAccessForNames()
func (c *CheckAccessControlsOldSyntaxTest) TestPackagePrivateAccessForNames() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @package */
var name = 'foo';
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"name;")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackagePrivateAccessForNames_googModule()
func (c *CheckAccessControlsOldSyntaxTest) TestPackagePrivateAccessForNames_googModule() {
	// TODO(b/133450410): This should be a visibility violation.
	c.DisableRewriteClosureCode()
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`goog.module('Foo');
/** @package */
var name = 'foo';
exports = name;
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				"goog.module('client'); const name = goog.require('Foo'); name;")))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackagePrivateAccessForNames_esModule()
func (c *CheckAccessControlsOldSyntaxTest) TestPackagePrivateAccessForNames_esModule() {
	// TODO(b/133450410): This should be a visibility violation.
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @package */
var name = 'foo';
export {name};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				"import {name} from '/foo/bar.js'; name;")))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackagePrivateAccessForProperties1()
func (c *CheckAccessControlsOldSyntaxTest) TestPackagePrivateAccessForProperties1() {
	c.TestSame(
		`/** @constructor */ function Foo() {}
/** @package */ Foo.prototype.bar = function() {};
Foo.prototype.baz = function() { this.bar(); }; (new Foo).bar();
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testPackagePrivateAccessForProperties2()
func (c *CheckAccessControlsOldSyntaxTest) TestPackagePrivateAccessForProperties2() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"), "/** @constructor */ function Foo() {}"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`/** @package */ Foo.prototype.bar = function() {};
Foo.prototype.baz = function() { this.bar(); }; (new Foo).bar();
`)))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackagePrivateAccessForProperties3()
func (c *CheckAccessControlsOldSyntaxTest) TestPackagePrivateAccessForProperties3() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */ function Foo() {}
/** @package */ Foo.prototype.bar = function() {}; (new Foo).bar();
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "baz.js"),
				"Foo.prototype.baz = function() { this.bar(); };")))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackagePrivateAccessForProperties4()
func (c *CheckAccessControlsOldSyntaxTest) TestPackagePrivateAccessForProperties4() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */ function Foo() {}
/** @package */ Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "baz.js"),
				"Foo.prototype['baz'] = function() { (new Foo()).bar(); };")))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackagePrivateAccessForProperties5()
func (c *CheckAccessControlsOldSyntaxTest) TestPackagePrivateAccessForProperties5() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */
function Parent () {
  /** @package */
  this.prop = 'foo';
};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`/**
 * @constructor
 * @extends {Parent}
 */
function Child() {
  this.prop = 'asdf';
}
Child.prototype = new Parent();
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackagePropAccess_inSamePackage_throughDestructuring()
func (c *CheckAccessControlsOldSyntaxTest) TestPackagePropAccess_inSamePackage_throughDestructuring() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */
function Foo() { }

/** @package */
Foo.prototype.bar = function() { };
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "quux.js"),
				`function f(/** !Foo */ x) {
  const {bar: bar} = x;
}
`)))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPackagePropAccess_inDifferentPackage_throughDestructuring()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPackagePropAccess_inDifferentPackage_throughDestructuring() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */
function Foo() { }

/** @package */
Foo.prototype.bar = function() { };
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`function f(/** !Foo */ x) {
  const {bar: bar} = x;
}
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPackagePrivateAccessForProperties1()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPackagePrivateAccessForProperties1() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				"/** @constructor */ function Foo() {} (new Foo).bar();"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`/** @package */ Foo.prototype.bar = function() {};
Foo.prototype.baz = function() { this.bar(); };
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPackagePrivateAccessForProperties2()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPackagePrivateAccessForProperties2() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */ function Foo() {}
/** @package */ Foo.prototype.bar = function() {};
Foo.prototype.baz = function() { this.bar(); };
`),
			SourceFile_FromCode(Compiler_JoinPathParts("baz", "quux.js"), "(new Foo).bar();")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPackagePrivateAccessForProperties3()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPackagePrivateAccessForProperties3() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */ function Foo() {}
/** @package */ Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				"/** @constructor */ function OtherFoo() { (new Foo).bar(); }")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPackagePrivateAccessForProperties4()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPackagePrivateAccessForProperties4() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */ function Foo() {}
/** @package */ Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`/** @constructor 
 * @extends {Foo} */
function SubFoo() { this.bar(); }
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPackagePrivateAccessForNamespaces()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPackagePrivateAccessForNamespaces() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @const */ var foo = {};
/** @package */ foo.bar = function() {};
`),
			SourceFile_FromCode(Compiler_JoinPathParts("baz", "quux.js"), "foo.bar();")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPackagePrivateAccessForProperties5()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPackagePrivateAccessForProperties5() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */ function Foo() {}
/** @package */ Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`/** @constructor 
 * @extends {Foo} */
function SubFoo() {};
SubFoo.prototype.baz = function() { this.bar(); }
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPackagePrivateAccessForProperties6()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPackagePrivateAccessForProperties6() {
	// Overriding a private property with a non-package-private property
	// in a different file causes problems.
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */ function Foo() {}
/** @package */ Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`/** @constructor 
 * @extends {Foo} */
function SubFoo() {};
SubFoo.prototype.bar = function() {};
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPackagePrivateAccessForProperties7()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPackagePrivateAccessForProperties7() {
	// It's OK to override a package-private property with a
	// non-package-private property in the same file, but you'll get
	// yelled at when you try to use it.
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @constructor */ function Foo() {}
/** @package */ Foo.prototype.bar = function() {};
/** @constructor 
 * @extends {Foo} */
function SubFoo() {};
SubFoo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				"SubFoo.prototype.baz = function() { this.bar(); }")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_OneFile()
func (c *CheckAccessControlsOldSyntaxTest) TestOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_OneFile() {
	c.TestError(
		`/**
* @fileoverview
* @package
*/
/** @struct @constructor */
Foo = function() {};
/** @private */
Foo.prototype.privateMethod_ = function() {};
/** @struct @constructor @extends {Foo} */
Bar = function() {};
/** @override */
Bar.prototype.privateMethod_ = function() {};
`,
		CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY)
}

// Java: CheckAccessControlsOldSyntaxTest.testOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_GoogModule()
func (c *CheckAccessControlsOldSyntaxTest) TestOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_GoogModule() {
	c.TestError(
		`/**
 * @fileoverview
 * @package
 */
goog.module('mod')
/** @struct @constructor */
Foo = function() {};
/** @private */
Foo.prototype.privateMethod_ = function() {};
/** @struct @constructor @extends {Foo} */
Bar = function() {};
/** @override */
Bar.prototype.privateMethod_ = function() {};
`,
		CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY)
}

// Java: CheckAccessControlsOldSyntaxTest.testOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_esModule()
func (c *CheckAccessControlsOldSyntaxTest) TestOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_esModule() {
	c.TestError(
		`/**
 * @fileoverview
 * @package
 */
/** @struct @constructor */
var Foo = function() {};
/** @private */
Foo.prototype.privateMethod_ = function() {};
/** @struct @constructor @extends {Foo} */
var Bar = function() {};
/** @override */
Bar.prototype.privateMethod_ = function() {};
export {Foo, Bar};
`,
		CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY)
}

// Java: CheckAccessControlsOldSyntaxTest.testNamespacedFunctionDoesNotNeedVisibilityRedeclInFileWithFileOverviewVisibility()
func (c *CheckAccessControlsOldSyntaxTest) TestNamespacedFunctionDoesNotNeedVisibilityRedeclInFileWithFileOverviewVisibility() {
	c.TestSame(
		`/**
 * @fileoverview
 * @package
 */
/** @return {string} */
foo.bar = function() { return 'asdf'; };
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_TwoFiles()
func (c *CheckAccessControlsOldSyntaxTest) TestOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_TwoFiles() {
	c.TestError(
		c.Srcs(
			`/** @struct @constructor */
Foo = function() {};
/** @protected */
Foo.prototype.protectedMethod = function() {};
`,
			`  /**
* @fileoverview 
* @package
*/
/** @struct @constructor @extends {Foo} */
Bar = function() {};
/** @override */
Bar.prototype.protectedMethod = function() {};
`),
		CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY)
}

// Java: CheckAccessControlsOldSyntaxTest.testOverrideWithoutVisibilityRedeclInFileWithNoFileOverviewOk()
func (c *CheckAccessControlsOldSyntaxTest) TestOverrideWithoutVisibilityRedeclInFileWithNoFileOverviewOk() {
	c.TestSame(
		`/** @struct @constructor */
Foo = function() {};
/** @private */
Foo.prototype.privateMethod_ = function() {};
/** @struct @constructor @extends {Foo} */
Bar = function() {};
/** @override */
Bar.prototype.privateMethod_ = function() {};
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testOverrideWithoutVisibilityRedeclInFileWithNoFileOverviewVisibilityOk()
func (c *CheckAccessControlsOldSyntaxTest) TestOverrideWithoutVisibilityRedeclInFileWithNoFileOverviewVisibilityOk() {
	c.TestSame(
		`/**
  * @fileoverview
  */
/** @struct @constructor */
Foo = function() {};
/** @private */
Foo.prototype.privateMethod_ = function() {};
/** @struct @constructor @extends {Foo} */
Bar = function() {};
/** @override */
Bar.prototype.privateMethod_ = function() {};
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testOverrideWithVisibilityRedeclInFileWithFileOverviewVisibilityOk_OneFile()
func (c *CheckAccessControlsOldSyntaxTest) TestOverrideWithVisibilityRedeclInFileWithFileOverviewVisibilityOk_OneFile() {
	c.TestSame(
		`/**
  * @fileoverview
  * @package
  */
/** @struct @constructor */
Foo = function() {};
/** @private */
Foo.prototype.privateMethod_ = function() {};
/** @struct @constructor @extends {Foo} */
Bar = function() {};
/** @override @private */
Bar.prototype.privateMethod_ = function() {};
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testOverrideWithVisibilityRedeclInFileWithFileOverviewVisibilityOk_TwoFiles()
func (c *CheckAccessControlsOldSyntaxTest) TestOverrideWithVisibilityRedeclInFileWithFileOverviewVisibilityOk_TwoFiles() {
	c.TestSame(
		c.Srcs(
			`/** @struct @constructor */
Foo = function() {};
/** @protected */
Foo.prototype.protectedMethod = function() {};
`,
			`  /**
* @fileoverview
* @package
*/
/** @struct @constructor @extends {Foo} */
Bar = function() {};
/** @override @protected */
Bar.prototype.protectedMethod = function() {};
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testPublicFileOverviewVisibilityDoesNotApplyToNameWithExplicitPackageVisibility()
func (c *CheckAccessControlsOldSyntaxTest) TestPublicFileOverviewVisibilityDoesNotApplyToNameWithExplicitPackageVisibility() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @public
 */
/** @constructor @package */ function Foo() {};
`),
			SourceFile_FromCode(Compiler_JoinPathParts("baz", "quux.js"), "new Foo();")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackageFileOverviewVisibilityDoesNotApplyToNameWithExplicitPublicVisibility()
func (c *CheckAccessControlsOldSyntaxTest) TestPackageFileOverviewVisibilityDoesNotApplyToNameWithExplicitPublicVisibility() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
/** @constructor @public */ function Foo() {};
`),
			SourceFile_FromCode(Compiler_JoinPathParts("baz", "quux.js"), "new Foo();")))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility()
func (c *CheckAccessControlsOldSyntaxTest) TestPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
/** @constructor */
var Foo = function() {};
`),
			SourceFile_FromCode(Compiler_JoinPathParts("baz", "quux.js"), "new Foo();")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility_googModule()
func (c *CheckAccessControlsOldSyntaxTest) TestPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility_googModule() {
	// TODO(b/133450410): Requiring 'Foo' should be a visibility violation.
	c.DisableRewriteClosureCode()
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
goog.module('Foo');
/** @constructor */
var Foo = function() {};
exports = Foo;
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				"goog.module('client'); const Foo = goog.require('Foo'); new Foo();")))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility_esModule()
func (c *CheckAccessControlsOldSyntaxTest) TestPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility_esModule() {
	// TODO(b/133450410): Importing 'Foo' should be a visibility violation.
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
/** @constructor */
var Foo = function() {};
export {Foo};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				"import {Foo} from '/foo/bar.js'; new Foo();")))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackageFileOverviewVisibilityDoesNotApplyToPropertyWithExplicitPublicVisibility()
func (c *CheckAccessControlsOldSyntaxTest) TestPackageFileOverviewVisibilityDoesNotApplyToPropertyWithExplicitPublicVisibility() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
/** @constructor */
Foo = function() {};
/** @public */
Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`var foo = new Foo();
foo.bar();
`)))
}

// Java: CheckAccessControlsOldSyntaxTest.testFileoverviewVisibilityDoesNotApplyToGoogProvidedNamespace1()
func (c *CheckAccessControlsOldSyntaxTest) TestFileoverviewVisibilityDoesNotApplyToGoogProvidedNamespace1() {
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode("foo.js", "goog.provide('foo');"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**

  * @fileoverview

  * @package

  */

goog.provide('foo.bar');
`),
			SourceFile_FromCode("bar.js", "goog.require('foo')")))
}

// Java: CheckAccessControlsOldSyntaxTest.testFileoverviewVisibilityDoesNotApplyToGoogProvidedNamespace2()
func (c *CheckAccessControlsOldSyntaxTest) TestFileoverviewVisibilityDoesNotApplyToGoogProvidedNamespace2() {
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
goog.provide('foo.bar');
`),
			SourceFile_FromCode("foo.js", "goog.provide('foo');"),
			SourceFile_FromCode(
				"bar.js",
				`goog.require('foo');
var x = foo;
`)))
}

// Java: CheckAccessControlsOldSyntaxTest.testFileoverviewVisibilityDoesNotApplyToGoogProvidedNamespace3()
func (c *CheckAccessControlsOldSyntaxTest) TestFileoverviewVisibilityDoesNotApplyToGoogProvidedNamespace3() {
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
goog.provide('one.two');
one.two.three = function(){};
`),
			SourceFile_FromCode(
				"baz.js",
				`goog.require('one.two');
var x = one.two;
`)))
}

// Java: CheckAccessControlsOldSyntaxTest.testFileoverviewVisibilityDoesNotApplyToGoogProvidedNamespace4()
func (c *CheckAccessControlsOldSyntaxTest) TestFileoverviewVisibilityDoesNotApplyToGoogProvidedNamespace4() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
goog.provide('one.two');
one.two.three = function(){};
`),
			SourceFile_FromCode(
				"baz.js",
				`goog.require('one.two');
var x = one.two.three();
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testPublicFileOverviewVisibilityDoesNotApplyToPropertyWithExplicitPackageVisibility()
func (c *CheckAccessControlsOldSyntaxTest) TestPublicFileOverviewVisibilityDoesNotApplyToPropertyWithExplicitPackageVisibility() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @public
 */
/** @constructor */
Foo = function() {};
/** @package */
Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`var foo = new Foo();
foo.bar();
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testPublicFileOverviewVisibilityAppliesToPropertyWithoutExplicitVisibility()
func (c *CheckAccessControlsOldSyntaxTest) TestPublicFileOverviewVisibilityAppliesToPropertyWithoutExplicitVisibility() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @public
 */
/** @constructor */
Foo = function() {};
Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`var foo = new Foo();
foo.bar();
`)))
}

// Java: CheckAccessControlsOldSyntaxTest.testPackageFileOverviewVisibilityAppliesToPropertyWithoutExplicitVisibility()
func (c *CheckAccessControlsOldSyntaxTest) TestPackageFileOverviewVisibilityAppliesToPropertyWithoutExplicitVisibility() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
/** @constructor */
Foo = function() {};
Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`var foo = new Foo();
foo.bar();
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testFileOverviewVisibilityComesFromDeclarationFileNotUseFile()
func (c *CheckAccessControlsOldSyntaxTest) TestFileOverviewVisibilityComesFromDeclarationFileNotUseFile() {
	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
/** @constructor */
Foo = function() {};
Foo.prototype.bar = function() {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`/**
 * @fileoverview
 * @public
 */
var foo = new Foo();
foo.bar();
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoExceptionsWithBadConstructors1()
func (c *CheckAccessControlsOldSyntaxTest) TestNoExceptionsWithBadConstructors1() {
	c.TestSame(
		c.Srcs(
			`function Foo() { (new SubFoo).bar(); }
/** @constructor */ function SubFoo() {}
/** @protected */ SubFoo.prototype.bar = function() {};
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoExceptionsWithBadConstructors2()
func (c *CheckAccessControlsOldSyntaxTest) TestNoExceptionsWithBadConstructors2() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() {};
/** @constructor */
function SubFoo() {}
/** @protected */
SubFoo.prototype.bar = function() { (new Foo).bar(); };
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testGoodOverrideOfProtectedProperty()
func (c *CheckAccessControlsOldSyntaxTest) TestGoodOverrideOfProtectedProperty() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() { }
/** @protected */ Foo.prototype.bar = function() {};
`,
			`/** @constructor 
 * @extends {Foo} */
function SubFoo() {}
/** @inheritDoc */ SubFoo.prototype.bar = function() {};
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testBadOverrideOfProtectedProperty()
func (c *CheckAccessControlsOldSyntaxTest) TestBadOverrideOfProtectedProperty() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() { }
/** @protected */ Foo.prototype.bar = function() {};
`,
			`/** @constructor 
 * @extends {Foo} */
function SubFoo() {}
/** @private */ SubFoo.prototype.bar = function() {};
`),
		CheckAccessControls_VISIBILITY_MISMATCH)
}

// Java: CheckAccessControlsOldSyntaxTest.testBadOverrideOfPrivateProperty()
func (c *CheckAccessControlsOldSyntaxTest) TestBadOverrideOfPrivateProperty() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() { }
/** @private */ Foo.prototype.bar = function() {};
`,
			`/** @constructor 
 * @extends {Foo} */
function SubFoo() {}
/** @protected */ SubFoo.prototype.bar = function() {};
`),
		CheckAccessControls_PRIVATE_OVERRIDE)
}

// Java: CheckAccessControlsOldSyntaxTest.testAccessOfStaticMethodOnPrivateConstructor()
func (c *CheckAccessControlsOldSyntaxTest) TestAccessOfStaticMethodOnPrivateConstructor() {
	c.TestSame(
		c.Srcs(
			`/** @constructor 
 * @private */ function Foo() { }
Foo.create = function() { return new Foo(); };
`,
			"Foo.create()"))
}

// Java: CheckAccessControlsOldSyntaxTest.testAccessOfStaticMethodOnPrivateQualifiedConstructor()
func (c *CheckAccessControlsOldSyntaxTest) TestAccessOfStaticMethodOnPrivateQualifiedConstructor() {
	c.TestSame(
		c.Srcs(
			`/** @constructor 
 * @private */ goog.Foo = function() { };
goog.Foo.create = function() { return new goog.Foo(); };
`,
			"goog.Foo.create()"))
}

// Java: CheckAccessControlsOldSyntaxTest.testInstanceofOfPrivateConstructor()
func (c *CheckAccessControlsOldSyntaxTest) TestInstanceofOfPrivateConstructor() {
	c.TestSame(
		c.Srcs(
			`/** @constructor 
 * @private */ goog.Foo = function() { };
goog.Foo.create = function() { return new goog.Foo(); };
`,
			"goog instanceof goog.Foo"))
}

// Java: CheckAccessControlsOldSyntaxTest.testOkAssignmentOfDeprecatedProperty()
func (c *CheckAccessControlsOldSyntaxTest) TestOkAssignmentOfDeprecatedProperty() {
	c.TestSame(
		`/** @constructor */ function Foo() {
 /** @deprecated */ this.bar = 3;
}
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testBadReadOfDeprecatedProperty()
func (c *CheckAccessControlsOldSyntaxTest) TestBadReadOfDeprecatedProperty() {
	c.TestDepProp(
		`/** @constructor */ function Foo() {
 /** @deprecated GRR */ this.bar = 3;
  this.baz = this.bar;
}
`,
		"Property bar of type Foo has been deprecated: GRR")
}

// Java: CheckAccessControlsOldSyntaxTest.testAutoboxedDeprecatedProperty()
func (c *CheckAccessControlsOldSyntaxTest) TestAutoboxedDeprecatedProperty() {
	c.TestError(
		c.Externs(CompilerTestCase_DEFAULT_EXTERNS),
		c.Srcs("/** @deprecated %s */ String.prototype.prop; function f() { return 'x'.prop; }"),
		c.Error(CheckAccessControls_DEPRECATED_PROP_REASON))
}

// Java: CheckAccessControlsOldSyntaxTest.testAutoboxedPrivateProperty()
func (c *CheckAccessControlsOldSyntaxTest) TestAutoboxedPrivateProperty() {
	c.TestError(
		c.Externs(CompilerTestCase_DEFAULT_EXTERNS+"/** @private */ String.prototype.prop;"),
		c.Srcs("function f() { return 'x'.prop; }"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsOldSyntaxTest.testNullableDeprecatedProperty()
func (c *CheckAccessControlsOldSyntaxTest) TestNullableDeprecatedProperty() {
	c.TestError(
		`/** @constructor */ function Foo() {}
/** @deprecated */ Foo.prototype.length;
/** @param {?Foo} x */ function f(x) { return x.length; }
`,
		CheckAccessControls_DEPRECATED_PROP)
}

// Java: CheckAccessControlsOldSyntaxTest.testNullablePrivateProperty()
func (c *CheckAccessControlsOldSyntaxTest) TestNullablePrivateProperty() {
	c.TestError(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @private */ Foo.prototype.length;
`,
			"/** @param {?Foo} x */ function f(x) { return x.length; }"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivatePropertyByConvention1()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivatePropertyByConvention1() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() {}
/** @type {number} */ Foo.prototype.length_;
`,
			"/** @param {?Foo} x */ function f(x) { return x.length_; }\n"))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoPrivatePropertyByConvention2()
func (c *CheckAccessControlsOldSyntaxTest) TestNoPrivatePropertyByConvention2() {
	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() {
  /** @type {number} */ this.length_ = 1;
}
/** @type {number} */ Foo.prototype.length_;
`,
			"/** @param {Foo} x */ function f(x) { return x.length_; }\n"))
}

// Java: CheckAccessControlsOldSyntaxTest.testDeclarationAndConventionConflict1()
func (c *CheckAccessControlsOldSyntaxTest) TestDeclarationAndConventionConflict1() {
	// "private" convention is no longer recognized
	c.TestSame("/** @constructor */ function Foo() {} /** @protected */ Foo.prototype.length_;")
}

// Java: CheckAccessControlsOldSyntaxTest.testDeclarationAndConventionConflict2()
func (c *CheckAccessControlsOldSyntaxTest) TestDeclarationAndConventionConflict2() {
	// "private" convention is no longer recognized
	c.TestSame("/** @constructor */ function Foo() {  /** @protected */ this.length_ = 1;\n}\n")
}

// Java: CheckAccessControlsOldSyntaxTest.testDeclarationAndConventionConflict3()
func (c *CheckAccessControlsOldSyntaxTest) TestDeclarationAndConventionConflict3() {
	// "private" convention is no longer recognized
	c.TestSame("/** @protected */ var Foo_;")
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyByJsdoc_initialAssignmentOk()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyByJsdoc_initialAssignmentOk() {
	c.TestNoWarning(
		`/** @constructor */ function A() {
  /** @const */ this.bar = 3;
}
/** @constructor */ function B() {
  /** @const */ this.bar = 3;
}
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyByConvention_initialAssignmentOk()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyByConvention_initialAssignmentOk() {
	c.TestNoWarning(
		`/** @constructor */ function A() {
  this.BAR = 3;
}
/** @constructor */ function B() {
  this.BAR = 3;
}
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyByJsdoc_reassignmentWarns()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyByJsdoc_reassignmentWarns() {
	c.TestError(
		`/** @constructor */ function A() {
/** @const */ this.bar = 3;}
/** @constructor */ function B() {
/** @const */ this.bar = 3;this.bar += 4;}
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyReassigned_crossModuleWithCollidingNames()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyReassigned_crossModuleWithCollidingNames() {
	c.DisableRewriteClosureCode()
	c.TestNoWarning(
		c.Srcs(
			`goog.module('mod1');
/** @constructor */
function A() {
  /** @const */
  this.bar = 3;
}
`,
			`goog.module('mod2');
/** @constructor */
function A() {
  /** @const */
  this.bar = 3;
}
`))
}

// Java: CheckAccessControlsOldSyntaxTest.testNoConstantPropertyByConvention()
func (c *CheckAccessControlsOldSyntaxTest) TestNoConstantPropertyByConvention() {
	c.TestSame(
		`/** @constructor */ function A() {
this.BAR = 3;}
/** @constructor */ function B() {
this.BAR = 3;this.BAR += 4;}
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyByJsdocOnPrototype_reassignmentWarns()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyByJsdocOnPrototype_reassignmentWarns() {
	c.TestError(
		`/** @constructor */ function Foo() {}
/** @const */ Foo.prototype.prop = 2;
var foo = new Foo();
foo.prop = 3;
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyByConventionOnPrototype_notEnforced()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyByConventionOnPrototype_notEnforced() {
	c.TestSame(
		`/** @constructor */ function Foo() {}
Foo.prototype.PROP = 2;
var foo = new Foo();
foo.PROP = 3;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnConstNamespaceByAssignment_initialAssignmentOk()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnConstNamespaceByAssignment_initialAssignmentOk() {
	c.TestNoWarning(
		`/** @const */ var o = {};
/** @const */ o.x = 1;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnConstNamespaceByAssignment_reassignmentWarns()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnConstNamespaceByAssignment_reassignmentWarns() {
	c.TestError(
		`/** @const */ var o = {};
/** @const */ o.x = 1;
o.x = 2;
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnNamespaceByAssignment_reassignmentWarns()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnNamespaceByAssignment_reassignmentWarns() {
	c.TestError(
		`var o = {};
/** @const */ o.x = 1;
o.x = 2;
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testNamespaceConstantPropertyOnNamespace_assigningSeparateNamespacesOk()
func (c *CheckAccessControlsOldSyntaxTest) TestNamespaceConstantPropertyOnNamespace_assigningSeparateNamespacesOk() {
	c.TestSame(
		`/** @const */ var o = {};
/** @const */ o.x = 1;
/** @const */ var o2 = {};
/** @const */ o2.x = 1;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnObjectLiteralByLiteralKey_initialAssignmentOk()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnObjectLiteralByLiteralKey_initialAssignmentOk() {
	c.TestNoWarning(
		`var o = {
  /** @const */ x: 1
 };
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnObjectLiteralByLiteralKey_reassignmentWarns()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnObjectLiteralByLiteralKey_reassignmentWarns() {
	c.TestError(
		`var o = {
  /** @const */ x: 1
 };
o.x = 2;
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnConstObjectLiteralByLiteralKey_reassignmentWarns()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnConstObjectLiteralByLiteralKey_reassignmentWarns() {
	c.TestError(
		`/** @const */
var o = {
  /** @const */ x: 1
 };
o.x = 2;
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnObjectLiteralByConvention_notEnforced()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnObjectLiteralByConvention_notEnforced() {
	c.TestSame(
		`var o = {
  XYZ: 1
 };
o.XYZ = 2;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnConstObjectLiteralByConvention_notEnforced()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnConstObjectLiteralByConvention_notEnforced() {
	c.TestSame(
		`/** @const */ var o = {
  XYZ: 1
};
o.XYZ = 2;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty4()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty4() {
	c.TestError(
		`/** @constructor */ function cat(name) {}
/** @const */ cat.test = 1;
cat.test *= 2;
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty4b()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty4b() {
	c.TestSame(
		`/** @constructor */ function cat(name) {}
cat.TEST = 1;
cat.TEST *= 2;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty5()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty5() {
	c.TestError(
		`/** @constructor */ function Foo() { this.prop = 1;}
/** @const */ Foo.prototype.prop;
Foo.prototype.prop = 2
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty6()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty6() {
	c.TestError(
		`/** @constructor */ function Foo() { this.prop = 1;}
/** @const */ Foo.prototype.prop = 2;
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty7()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty7() {
	c.TestSame(
		`/** @constructor */ function Foo() {}
Foo.prototype.bar_ = function() {};
/** @constructor 
 * @extends {Foo} */
function SubFoo() {};
/** @const */ /** @override */ SubFoo.prototype.bar_ = function() {};
SubFoo.prototype.baz = function() { this.bar_(); }
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty8()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty8() {
	c.TestSame(
		`/** @const */ var o = { /** @const */ x: 1 };
var y = o.x;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty9()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty9() {
	c.TestSame(
		`/** @constructor */ function A() {
/** @const */ this.bar = 3;}
/** @constructor */ function B() {
this.bar = 4;}
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty10a()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty10a() {
	c.TestSame(
		`/** @constructor */ function Foo() { this.prop = 1;}
/** @const */ Foo.prototype.prop;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty10b()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty10b() {
	c.TestSame(
		`/** @constructor */ function Foo() { this.PROP = 1;}
Foo.prototype.PROP;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty11()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty11() {
	c.TestError(
		`/** @constructor */ function Foo() {}
/** @const */ Foo.prototype.bar;
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() { this.bar = 5; this.bar = 6; }
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty12()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty12() {
	c.TestSame(
		`/** @constructor */ function Foo() {}
/** @const */ Foo.prototype.bar;
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() { this.bar = 5; }
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo2() { this.bar = 5; }
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty13()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty13() {
	c.TestError(
		`/** @constructor */ function Foo() {}
/** @const */ Foo.prototype.bar;
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() { this.bar = 5; }
/**
 * @constructor
 * @extends {SubFoo}
 */ function SubSubFoo() { this.bar = 5; }
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty14()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty14() {
	c.TestError(
		`/** @constructor */ function Foo() {
/** @const */ this.bar = 3; delete this.bar; }
`,
		CheckAccessControls_CONST_PROPERTY_DELETED)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty_recordType()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty_recordType() {
	c.Test(
		c.Srcs(
			`/** @record */
function Foo() { }
/** @const {number} */
Foo.prototype.bar;

const /** !Foo */ x = {
  bar: 9,
};
x.bar = 0;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE).
			WithMessageContaining("unknown location due to structural typing"))
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty_fromExternsOrIjs_duplicateExternOk()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty_fromExternsOrIjs_duplicateExternOk() {
	c.TestSame(
		c.Externs(
			`/** @constructor */ function Foo() {}
/** @const */ Foo.prototype.PROP;
/** @const */ Foo.prototype.PROP;
`),
		c.Srcs(""))
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty_fromExternsOrIjs()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty_fromExternsOrIjs() {
	c.Test(
		c.Externs(
			`/** @constructor */ function Foo() {}
/** @const */ Foo.prototype.PROP;
`),
		c.Srcs(
			`var f = new Foo();
f.PROP = 1;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE).WithMessageContaining("at externs:2:"))
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty15()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty15() {
	c.TestSame(
		`/** @constructor */ function Foo() {};
Foo.CONST = 100;
/** @type {Foo} */
var foo = new Foo();
/** @type {number} */
foo.CONST = Foo.CONST;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty_convention_not_enforced1()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty_convention_not_enforced1() {
	c.TestSame(
		`/** @constructor */ function Foo() { this.CONST = 100; };
/** @type {Foo} */
var foo = new Foo();
/** @type {number} */
foo.CONST = 0;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty_convention_not_enforced2()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty_convention_not_enforced2() {
	c.TestSame(
		`/** @constructor */ function Foo() {};
Foo.prototype.CONST = 100;
/** @type {Foo} */
var foo = new Foo();
/** @type {number} */
foo.CONST = 0;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty_convention_not_enforced3()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty_convention_not_enforced3() {
	c.TestSame(
		`/** @constructor */ function Bar() {this.CONST = 100;};
/** @constructor 
 @extends {Bar} */ function Foo() {};
/** @type {Foo} */
var foo = new Foo();
/** @type {number} */
foo.CONST = 0;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty16()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty16() {
	c.TestSame(
		`/** @constructor */ function Foo() {};
Foo.CONST = 100;
/** @constructor */ function Bar() {};
Bar.CONST = 100;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty17()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty17() {
	c.TestSame(
		`function Foo() {};
Foo.CONST = 100;
function Bar() {};
Bar.CONST = 100;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty18()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty18() {
	c.TestSame(
		`/** @param {string} a */
function Foo(a) {};
Foo.CONST = 100;
/** @param {string} a */
function Bar(a) {};
Bar.CONST = 100;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantProperty19()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantProperty19() {
	c.TestSame(
		`/** @param {string} a */
function Foo(a) {};
Foo.CONST = 100;
/** @param {number} a */
function Bar(a) {};
Bar.CONST = 100;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testNamespaceOnStructuralInterface()
func (c *CheckAccessControlsOldSyntaxTest) TestNamespaceOnStructuralInterface() {
	// NOTE: The FooStatic pattern is not uncommon in TS externs.  If additional sub-namespaces are
	// defined underneath, tsickle generates code that looks like this.  It is important that the
	// moment.unitOfTime declaration be treated as a declaration so that it doesn't just look like
	// a mutation of the readonly property that it just defined.
	c.TestNoWarning(
		`/** @externs */
/** @record */
function MomentStatic() {}
/** @return {number} */
MomentStatic.prototype.now = function() {};
/** @const */
moment.unitOfTime = {};
/** @const {!MomentStatic} */
var moment;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnStructuralInterfaceByJsdoc_initializeFromObjLitOk()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnStructuralInterfaceByJsdoc_initializeFromObjLitOk() {
	c.TestNoWarning(
		`/** @record */
function Foo() {
  /** @const {number} */ this.bar;
}
var /** !Foo */ foo = {bar: 1};
var /** !Foo */ baz = {bar: 2};
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnStructuralInterfaceByJsdoc_reassignmentWarns()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnStructuralInterfaceByJsdoc_reassignmentWarns() {
	c.TestError(
		`/** @record */
function Foo() {
  /** @const {number} */ this.bar;
}
var /** !Foo */ foo = {bar: 1};
foo.bar = 2;
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnStructuralInterfaceByJsdoc_assignmentOnOpaqueObjectWarns()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnStructuralInterfaceByJsdoc_assignmentOnOpaqueObjectWarns() {
	c.TestError(
		`/** @record */
function Foo() {
  /** @const {number} */ this.bar;
}
function f(/** !Foo */ foo) {
  foo.bar = 2;
}
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testConstantPropertyOnStructuralInterfaceByConvention_noReassignmentWarns()
func (c *CheckAccessControlsOldSyntaxTest) TestConstantPropertyOnStructuralInterfaceByConvention_noReassignmentWarns() {
	c.TestSame(
		`/** @record */
function Foo() {
  /** @type {number} */ this.BAR;
}
var /** !Foo */ foo = {BAR: 1};
foo.BAR = 2;
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testFunctionWithNewType_canReturnFinalClass()
func (c *CheckAccessControlsOldSyntaxTest) TestFunctionWithNewType_canReturnFinalClass() {
	c.TestNoWarning(
		`goog.forwardDeclare('Parent');
/** @type {function(new: Parent): !Parent} */
const PatchedParent = function Parent() { return /** @type {?} */ (0); }
`)
	c.TestNoWarning(
		`/** @constructor @final */
const Parent = function() {}
/** @type {function(new: Parent): !Parent} */
const PatchedParent = function Parent() { return /** @type {?} */ (0); }
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testFinalClassCannotBeSubclassed()
func (c *CheckAccessControlsOldSyntaxTest) TestFinalClassCannotBeSubclassed() {
	c.TestError(
		`/**
 * @constructor
 * @final
 */ var Foo = function() {};
/**
 * @constructor
 * @extends {Foo}*
 */ var Bar = function() {};
`,
		CheckAccessControls_EXTEND_FINAL_CLASS)
	c.TestError(
		`/**
 * @constructor
 * @final
 */ function Foo() {};
/**
 * @constructor
 * @extends {Foo}*
 */ function Bar() {};
`,
		CheckAccessControls_EXTEND_FINAL_CLASS)
	c.TestSame(
		`/**
 * @constructor
 * @const
 */ var Foo = function() {};
/**
 * @constructor
 * @extends {Foo}
 */ var Bar = function() {};
`)
}

// Java: CheckAccessControlsOldSyntaxTest.testCircularPrototypeLink()
func (c *CheckAccessControlsOldSyntaxTest) TestCircularPrototypeLink() {
	// NOTE: this does yield a useful warning, except we don't check for it in this test:
	//      WARNING - Cycle detected in inheritance chain of type Foo
	// This warning already has a test: TypeCheckTest::testPrototypeLoop.
	c.TestError(
		`/** @constructor @extends {Foo} */ function Foo() {}
/** @const */ Foo.prop = 1;
Foo.prop = 2;
`,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE)
}

// Java: CheckAccessControlsOldSyntaxTest.testCheckFinalClass_prototypeMethodMarkedCtor_butNotOwnerFunction_doesNotCrash()
func (c *CheckAccessControlsOldSyntaxTest) TestCheckFinalClass_prototypeMethodMarkedCtor_butNotOwnerFunction_doesNotCrash() {
	// Covers an edge case reported in b/129361702.
	c.TestSame(
		`function Foo() {}

Foo.prototype = {
  /** @constructor */
  init: function() { }
};
`)
}

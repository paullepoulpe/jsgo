// Port of com/google/javascript/jscomp/Es6RewriteArrowFunctionTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/Es6RewriteArrowFunctionTest.java)

/*
 * Copyright 2015 The Closure Compiler Authors.
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

import "testing"

// Es6RewriteArrowFunctionTest has the test cases for the transpilation pass
// that converts arrow functions.
type Es6RewriteArrowFunctionTest struct {
	CompilerTestCaseBase

	languageOut CompilerOptions_LanguageMode
}

// Java: Es6RewriteArrowFunctionTest()
func NewEs6RewriteArrowFunctionTest() *Es6RewriteArrowFunctionTest {
	c := &Es6RewriteArrowFunctionTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_MINIMAL_EXTERNS)
	return c
}

func TestEs6RewriteArrowFunctionTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteArrowFunctionTest() })
}

// Java: Es6RewriteArrowFunctionTest.setUp()
func (c *Es6RewriteArrowFunctionTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT3

	c.EnableNormalize()
	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
	c.SetGenericNameReplacements(
		transpileTests_ImmutableMapOf(
			"$jscomp$this$UID", "$jscomp$this$", "$jscomp$arguments$UID", "$jscomp$arguments$"))
}

// Java: Es6RewriteArrowFunctionTest.getOptions()
func (c *Es6RewriteArrowFunctionTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetLanguageOut(c.languageOut)
	return options
}

// Java: Es6RewriteArrowFunctionTest.getProcessor(Compiler)
func (c *Es6RewriteArrowFunctionTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs6RewriteArrowFunction(compiler)
}

// Java: Es6RewriteArrowFunctionTest.testAssigningArrowToVariable_BlockBody()
func (c *Es6RewriteArrowFunctionTest) TestAssigningArrowToVariable_BlockBody() {
	c.Test("var f = x => { return x+1; };", "var f = function(x) { return x+1; };")
}

// Java: Es6RewriteArrowFunctionTest.testAssigningArrowToVariable_ExpressionBody()
func (c *Es6RewriteArrowFunctionTest) TestAssigningArrowToVariable_ExpressionBody() {
	c.Test("var f = x => x+1;", "var f = function(x) { return x+1; };")
}

// Java: Es6RewriteArrowFunctionTest.testPassingArrowToMethod_ExpressionBody()
func (c *Es6RewriteArrowFunctionTest) TestPassingArrowToMethod_ExpressionBody() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddArray().Build()),
		c.Srcs("var odds = [1,2,3,4].filter((n) => n%2 == 1);"),
		c.Expected("var odds = [1,2,3,4].filter(function(n) { return n%2 == 1; });"))
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInArrow_ExpressionBody()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInArrow_ExpressionBody() {
	c.Test(
		"var f = () => this;",
		"const $jscomp$this$UID$0 = this; var f = function() { return $jscomp$this$UID$0; };")
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInArrow_BlockBody()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInArrow_BlockBody() {
	c.Test(
		c.Externs(
			`window.init = function() { };
window.doThings = function() { };
window.done = function() { };
`),
		c.Srcs(
			`var f = x => {
  this.init();
  this.doThings();
  this.done();
};
`),
		c.Expected(
			`const $jscomp$this$UID$0 = this;
var f = function(x) {
  $jscomp$this$UID$0.init();
  $jscomp$this$UID$0.doThings();
  $jscomp$this$UID$0.done();
};
`))
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInArrowPlacesAliasAboveContainingStatement()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInArrowPlacesAliasAboveContainingStatement() {
	// We use `switch` here because it's a very complex kind of statement.
	c.Test(
		"switch(a) { case b: (() => { this; })(); }",
		`const $jscomp$this$UID$0 = this;
switch(a) {
  case b:
    (function() { $jscomp$this$UID$0; })();
}
`)
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInMultipleArrowsPlacesOneAliasAboveContainingStatement()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInMultipleArrowsPlacesOneAliasAboveContainingStatement() {
	// We use `switch` here because it's a very complex kind of statement.
	c.Test(
		`switch(a) {
  case b:
    (() => { this; })();
  case c:
    (() => { this; })();
}
`,
		`const $jscomp$this$UID$0 = this;
switch(a) {
  case b:
    (function() { $jscomp$this$UID$0; })();
  case c:
    (function() { $jscomp$this$UID$0; })();
}
`)
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInMultipleArrowsPlacesOneAliasAboveAllContainingStatements()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInMultipleArrowsPlacesOneAliasAboveAllContainingStatements() {
	// We use `switch` here because it's a very complex kind of statement.
	c.Test(
		`switch(a) {
  case b:
    (() => { this; })();
}
switch (c) {
  case d:
    (() => { this; })();
}
`,
		`const $jscomp$this$UID$0 = this;
switch(a) {
  case b:
    (function() { $jscomp$this$UID$0; })();
}
switch (c) {
  case d:
    (function() { $jscomp$this$UID$0; })();
}
`)
}

// Java: Es6RewriteArrowFunctionTest.testCapturingEnclosingFunctionArgumentsInArrow()
func (c *Es6RewriteArrowFunctionTest) TestCapturingEnclosingFunctionArgumentsInArrow() {
	c.Test(
		`function f() {
  var x = () => arguments;
}
`,
		`function f() {
  const $jscomp$arguments$UID$1 = arguments;
  var x = function() { return $jscomp$arguments$UID$1; };
}
`)
}

// Java: Es6RewriteArrowFunctionTest.testAssigningArrowToObjectLiteralField_ExpressionBody()
func (c *Es6RewriteArrowFunctionTest) TestAssigningArrowToObjectLiteralField_ExpressionBody() {
	c.Test("var obj = { f: () => 'bar' };", "var obj = { f: function() { return 'bar'; } };")
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInArrowFromClassMethod()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInArrowFromClassMethod() {
	// TODO(b/76024335): Enable these validations and checks.
	// We need to test classes the type-checker doesn't understand class syntax and fails before the
	// test even runs.
	c.DisableTypeInfoValidation()
	c.DisableTypeCheck()

	c.Test(
		`class C {
  constructor() {
    this.counter = 0;
  }

  init() {
    document.onclick = () => this.logClick();
  }

  logClick() {
     this.counter++;
  }
}
`,
		`class C {
  constructor() {
    this.counter = 0;
  }

  init() {
    const $jscomp$this$UID$2 = this;
    document.onclick = function() {return $jscomp$this$UID$2.logClick()}
  }

  logClick() {
     this.counter++;
  }
}
`)
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInArrowFromClassConstructorWithSuperCall()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInArrowFromClassConstructorWithSuperCall() {
	// TODO(b/76024335): Enable these validations and checks.
	// We need to test super, but super only makes sense in the context of a class, but
	// the type-checker doesn't understand class syntax and fails before the test even runs.
	c.DisableTypeInfoValidation()
	c.DisableTypeCheck()

	c.Test(
		`class B {
  constructor(x) {
    this.x = x;
  }
}
class C extends B {
  constructor(x, y) {
    console.log('statement before super');
    super(x);
    this.wrappedXGetter = () => this.x;
    this.y = y;
    this.wrappedYGetter = () => this.y;
  }
}
`,
		"class B {\n  constructor(x) {\n    this.x = x;\n  }\n}\nclass C extends B {\n  constructor(x$jscomp$1, y) {\n    console.log('statement before super');\n    super(x$jscomp$1);\n    const $jscomp$this$UID$2 = this; // Must not use `this` before super() call.\n    this.wrappedXGetter = function() { return $jscomp$this$UID$2.x; };\n    this.y = y;\n    this.wrappedYGetter = function() { return $jscomp$this$UID$2.y; };\n  }\n}\n")
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInArrowFromClassConstructorWithMultipleSuperCallPaths()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInArrowFromClassConstructorWithMultipleSuperCallPaths() {
	// TODO(b/76024335): Enable these validations and checks.
	// We need to test super, but super only makes sense in the context of a class, but
	// the type-checker doesn't understand class syntax and fails before the test even runs.
	c.DisableTypeInfoValidation()
	c.DisableTypeCheck()

	c.Test(
		`class B {
  constructor(x) {
    this.x = x;
  }
}
class C extends B {
  constructor(x, y) {
    if (x < 1) {
      super(x);
    } else {
      super(-x);
    }
    this.wrappedXGetter = () => this.x;
    this.y = y;
    this.wrappedYGetter = () => this.y;
  }
}
`,
		"class B {\n  constructor(x) {\n    this.x = x;\n  }\n}\nclass C extends B {\n  constructor(x$jscomp$1, y) {\n    if (x$jscomp$1 < 1) {\n      super(x$jscomp$1);\n    } else {\n      super(-x$jscomp$1);\n    }\n    const $jscomp$this$UID$2 = this; // Must not use `this` before super() call.\n    this.wrappedXGetter = function() { return $jscomp$this$UID$2.x; };\n    this.y = y;\n    this.wrappedYGetter = function() { return $jscomp$this$UID$2.y; };\n  }\n}\n")
}

// Java: Es6RewriteArrowFunctionTest.testMultipleArrowsInSameFreeScope()
func (c *Es6RewriteArrowFunctionTest) TestMultipleArrowsInSameFreeScope() {
	c.Test(
		"var a1 = x => x+1; var a2 = x => x-1;",
		`var a1 = function(x) { return x+1; };
var a2 = function(x$jscomp$1) { return x$jscomp$1-1; };
`)
}

// Java: Es6RewriteArrowFunctionTest.testMultipleArrowsInSameFunctionScope()
func (c *Es6RewriteArrowFunctionTest) TestMultipleArrowsInSameFunctionScope() {
	c.Test(
		"function f() { var a1 = x => x+1; var a2 = x => x-1; }",
		`function f() {
  var a1 = function(x) { return x+1; };
  var a2 = function(x$jscomp$1) { return x$jscomp$1-1; };
}
`)
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInMultipleArrowsInSameFunctionScope()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInMultipleArrowsInSameFunctionScope() {
	c.Test(
		`({
  x: 0,
  y: 'a',
  f: function() {
    var a1 = () => this.x;
    var a2 = () => this.y;
  },
})
`,
		`({
  x: 0,
  y: 'a',
  f: function() {
    const $jscomp$this$UID$1 = this;
    var a1 = function() { return $jscomp$this$UID$1.x; };
    var a2 = function() { return $jscomp$this$UID$1.y; };
  },
})
`)
}

// Java: Es6RewriteArrowFunctionTest.testGeneratedVariableDeclarations_placedAfterFunctionDeclarations()
func (c *Es6RewriteArrowFunctionTest) TestGeneratedVariableDeclarations_placedAfterFunctionDeclarations() {
	c.Test(
		`({
  x: 0,
  y: 'a',
  f: function() {
    function foo() { this.x;}
    var a2 = () => this.y;
  },
})
`,
		`({
  x: 0,
  y: 'a',
  f: function() {
// stays hoisted
    function foo() { this.x; }
// variable declarations placed after function declarations
    const $jscomp$this$UID$1 = this;
    var a2 = function() { return $jscomp$this$UID$1.y; };
  },
})
`)
}

// Java: Es6RewriteArrowFunctionTest.testPassingMultipleArrowsInSameFreeScopeAsMethodParams()
func (c *Es6RewriteArrowFunctionTest) TestPassingMultipleArrowsInSameFreeScopeAsMethodParams() {
	c.Test(
		c.Externs(
			CompilerTestCase_MINIMAL_EXTERNS+
				`/**
 * @template Y
 * @param {function(T):Y} mapper
 * @return {!Array<Y>}
 */
Array.prototype.map = function(mapper) { };
`),
		c.Srcs("var a = [1,2,3,4]; var b = a.map(x => x+1).map(x => x*x);"),
		c.Expected(
			`var a = [1,2,3,4];
var b = a.map(function(x) { return x+1; }).map(function(x$jscomp$1) { return x$jscomp$1*x$jscomp$1; });
`))
}

// Java: Es6RewriteArrowFunctionTest.testMultipleArrowsInSameFunctionScopeAsMethodParams()
func (c *Es6RewriteArrowFunctionTest) TestMultipleArrowsInSameFunctionScopeAsMethodParams() {
	c.Test(
		c.Externs(
			CompilerTestCase_MINIMAL_EXTERNS+
				`/**
 * @template Y
 * @param {function(T):Y} mapper
 * @return {!Array<Y>}
 */
Array.prototype.map = function(mapper) { };
`),
		c.Srcs(
			`function f() {
  var a = [1,2,3,4];
  var b = a.map(x => x+1).map(x => x*x);
}
`),
		c.Expected(
			`function f() {
  var a = [1,2,3,4];
  var b = a.map(function(x) { return x+1; }).map(function(x$jscomp$1) { return x$jscomp$1*x$jscomp$1; });
}
`))
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInArrowFromNestedScopes()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInArrowFromNestedScopes() {
	c.Test(
		`var outer = {
  x: null,

  f: function() {
     var a1 = () => this.x;
     var inner = {
       y: null,

       f: function() {
         var a2 = () => this.y;
       }
     };
  }
}
`,
		`var outer = {
  x: null,

  f: function() {
     const $jscomp$this$UID$1 = this;
     var a1 = function() { return $jscomp$this$UID$1.x; }
     var inner = {
       y: null,

       f: function() {
         const $jscomp$this$UID$2 = this;
         var a2 = function() { return $jscomp$this$UID$2.y; }
       }
     };
  }
}
`)
}

// Java: Es6RewriteArrowFunctionTest.testCapturingThisInArrowWithNestedConstructor()
func (c *Es6RewriteArrowFunctionTest) TestCapturingThisInArrowWithNestedConstructor() {
	c.Test(
		`({
  f: null,

  g: function() {
    var setup = () => {
      /** @constructor @struct */
      function Foo() { this.x = 5; }

      this.f = new Foo;
    };
  },
})
`,
		`({
  f: null,

  g: function() {
    const $jscomp$this$UID$1 = this;
    var setup = function() {
      /** @constructor */
      function Foo() { this.x = 5; }

      $jscomp$this$UID$1.f = new Foo;
    };
  },
})
`)
}

// Java: Es6RewriteArrowFunctionTest.testNestingArrow()
func (c *Es6RewriteArrowFunctionTest) TestNestingArrow() {
	c.Test(
		c.Externs(""),
		c.Srcs(`var f = x =>
 y => x+y;`),
		c.Expected("var f = function(x) {return function(y) { return x+y; }; };"))
}

// Java: Es6RewriteArrowFunctionTest.testNestingArrowsCapturingThis()
func (c *Es6RewriteArrowFunctionTest) TestNestingArrowsCapturingThis() {
	c.Test(
		c.Externs("window.foo = function() { };"),
		c.Srcs("var f = (x => { var g = (y => { this.foo(); }) });"),
		c.Expected(
			`const $jscomp$this$UID$0 = this;
var f = function(x) {
  var g = function(y) {
    $jscomp$this$UID$0.foo();
  }
}
`))
}

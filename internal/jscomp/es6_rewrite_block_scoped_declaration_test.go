// Port of com/google/javascript/jscomp/Es6RewriteBlockScopedDeclarationTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/Es6RewriteBlockScopedDeclarationTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type Es6RewriteBlockScopedDeclarationTest struct {
	CompilerTestCaseBase
	runVarCheck bool
}

// Java: Es6RewriteBlockScopedDeclarationTest()
func NewEs6RewriteBlockScopedDeclarationTest() *Es6RewriteBlockScopedDeclarationTest {
	c := &Es6RewriteBlockScopedDeclarationTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_DEFAULT_EXTERNS)
	return c
}

func TestEs6RewriteBlockScopedDeclarationTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteBlockScopedDeclarationTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *Es6RewriteBlockScopedDeclarationTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: Es6RewriteBlockScopedDeclarationTest.getOptions()
func (c *Es6RewriteBlockScopedDeclarationTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_OFF)
	return options
}

// Java: Es6RewriteBlockScopedDeclarationTest.customSetUp()
func (c *Es6RewriteBlockScopedDeclarationTest) CustomSetUp() {
	c.AllowExternsChanges()
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableNormalize()
	c.EnableMultistageCompilation()
}

// Java: Es6RewriteBlockScopedDeclarationTest.getProcessor(Compiler)
func (c *Es6RewriteBlockScopedDeclarationTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		if c.runVarCheck {
			// Adds a Synthetic definition of $jscomp$lookupPolyfilledValue
			NewVarCheck(compiler).Process(externs, root)
		}
		NewEs6RewriteBlockScopedDeclaration(compiler).Process(externs, root)
	})
}

// Java: Es6RewriteBlockScopedDeclarationTest.loopClosureTest(Sources, Expected)
func (c *Es6RewriteBlockScopedDeclarationTest) LoopClosureTest(srcs CompilerTestCase_Sources, originalExpected *CompilerTestCase_Expected) {
	c.LoopClosureTestWithExterns(c.Externs(""), srcs, originalExpected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.loopClosureTest(Externs, Sources, Expected)
func (c *Es6RewriteBlockScopedDeclarationTest) LoopClosureTestWithExterns(externs *CompilerTestCase_Externs, srcs CompilerTestCase_Sources, originalExpected *CompilerTestCase_Expected) {
	loop := jutil.NewLinkedHashMap[string, string]()
	loop.Put("LOOP", "$jscomp$loop$")
	// change the "LOOP" in expected with "$jscomp$loop$m123..456" name
	modifiedExpected := c.Expected(
		UnitTestUtils_UpdateGenericVarNamesInExpectedFiles(
			srcs.(*CompilerTestCase_FlatSources), originalExpected, loop))
	param := jutil.NewLinkedHashMap[string, string]()
	param.Put("$PARAM", "$jscomp$loop_param$")
	// change the "$PARAM" in expected with "$jscomp$loop_param$m123..456" parameter name
	modifiedExpected =
		c.Expected(
			UnitTestUtils_UpdateGenericVarNamesInExpectedFiles(
				srcs.(*CompilerTestCase_FlatSources),
				modifiedExpected,
				param))
	thisArgs := jutil.NewLinkedHashMap[string, string]()
	thisArgs.Put("$THIS", "$jscomp$this$")
	thisArgs.Put("$ARGUMENTS", "$jscomp$arguments$")
	// change the "$THIS" and "$ARGUMENTS" in expected with "$jscomp$this$" and "$jscomp$arguments$"
	modifiedExpected =
		c.Expected(
			UnitTestUtils_UpdateGenericVarNamesInExpectedFiles(
				srcs.(*CompilerTestCase_FlatSources),
				modifiedExpected,
				thisArgs))
	c.Test(externs, srcs, modifiedExpected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testVarNameCollisionWithExterns()
func (c *Es6RewriteBlockScopedDeclarationTest) TestVarNameCollisionWithExterns() {
	externs := c.Externs("var url;")
	srcs := c.Srcs("export {}; { const url = ''; alert(url);}")
	expected := c.Expected("export {}; { /** @const */ var url$jscomp$1 = ''; alert(url$jscomp$1); }")
	c.Test(externs, srcs, expected)

	externs = c.Externs("var url;" + CompilerTestCase_DEFAULT_EXTERNS)
	srcs = c.Srcs("goog.module('main'); { const url = ''; alert(url);}")
	expected =
		c.Expected(
			"goog.module('main'); { /** @const */ var url$jscomp$1 = ''; alert(url$jscomp$1); }")
	c.Test(externs, srcs, expected)

	c.Test(
		c.Externs("var url;"),
		c.Srcs(
			`export {};
function foo() {
  const url = '';
  alert(url);
}
`),
		c.Expected(
			`export {};
function foo() {
  /** @const */ var url$jscomp$1 = '';
  alert(url$jscomp$1);
}
`))

	c.Test(
		c.Externs("var url;"+CompilerTestCase_DEFAULT_EXTERNS),
		c.Srcs(
			`goog.module('main');
function foo() {
  const url = '';
  alert(url);
}
`),
		c.Expected(
			`goog.module('main');
function foo() {
  /** @const */
  var url$jscomp$1 = '';
  alert(url$jscomp$1);
}
`))

	c.Test(
		c.Externs("var url;"+CompilerTestCase_DEFAULT_EXTERNS),
		c.Srcs(
			`function foo() {
  const url = '';
  alert(url);
}
`),
		c.Expected(
			`function foo() {
  /** @const */
  var url$jscomp$1 = '';
  alert(url$jscomp$1);
}
`))
}

// Java: Es6RewriteBlockScopedDeclarationTest.testSimple()
func (c *Es6RewriteBlockScopedDeclarationTest) TestSimple() {
	c.Test("let x = 3;", "var x = 3;")
	c.Test("const x = 3;", "/** @const */ var x = 3;")
	c.Test("const x = 1, y = 2;", "/** @const */ var x = 1; /** @const */ var y = 2;")
	c.Test("const a = 0; a;", "/** @const */ var a = 0; a;")
	c.Test(
		"if (a) { let x; }", //
		"if (a) { var x; }")
	c.Test("function f() { const x = 3; }", "function f() { /** @const */ var x = 3; }")
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLetShadowing()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLetShadowing() {
	srcs := c.Srcs(
		`function f() {
  var x = 1;
  if (a) {
    let x = 2;
    x = function() { return x; };
  }
  return x;
}
`)
	expected := c.Expected(
		`function f() {
  var x = 1;
  if (a) {
    var x$jscomp$1 = 2;
    x$jscomp$1 = function() { return x$jscomp$1; };
  }
  return x;
}
`)
	c.Test(srcs, expected)

	srcs =
		c.Srcs(
			`function f() {
  const x = 3;
  if (true) {
    let x;
  }
}
`)
	expected =
		c.Expected(
			`function f() {
  /** @const */ var x = 3;
  if (true) {
    var x$jscomp$1;
  }
}
`)
	c.Test(srcs, expected)

	srcs =
		c.Srcs(
			`function f() {
  var x = 1;
  if (a) {
    var g = function() { return x; };
    let x = 2;
    return g();
  }
}
`)
	expected =
		c.Expected(
			`function f() {
  var x = 1;
  if (a) {
    var g = function() { return x$jscomp$1; };
    var x$jscomp$1 = 2;
    return g();
  }
}
`)
	c.Test(srcs, expected)

	srcs =
		c.Srcs(
			`var x = 2;
function f() {
  x = 1;
  if (a) {
    let x = 2;
  }
}
`)
	expected =
		c.Expected(
			`var x = 2;
function f() {
  x = 1;
  if (a) {
    var x$jscomp$1 = 2;
  }
}
`)
	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLetShadowingUndeclaredWithVarCheck()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLetShadowingUndeclaredWithVarCheck() {
	c.runVarCheck = true
	srcs := c.Srcs(
		`function use(x) {}
function f() {
  {
    let inner$jscomp$1 = 2;
  }
  use(inner)
}
`)
	expected := c.Expected(
		`function use(x) {}
function f() {
  {
    var inner$jscomp$1 = 2;
  }
  use(inner)
}
`)
	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLetShadowingTranspileOnly()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLetShadowingTranspileOnly() {
	srcs := c.Srcs(
		`function use(x) {}
function f() {
  {
    let inner = 2;
  }
  use(inner)
}
`)
	expected := c.Expected(
		"function use(x) {}\nfunction f() {\n  {\n// Notice that without VarCheck running MakeDeclaredNamesUnique fails to\n// rename this variable to avoid conflicting with the undeclared global variable\n// used after this block in `use(inner)`.\n// This probably indicates a flaw in MakeDeclaredNamesUnique when an undeclared\n// variable comes after the scope of a local with the same name.\n// However, fixing this is not a high priority because:\n// 1. VarCheck always runs unless we are in transpile-only mode.\n// 2. Undeclared variables are already a known hazard with lots of warnings in the\n//    compiler to push you away from having them in your code.\n    var inner = 2;\n  }\n  use(inner)\n}\n")
	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLetShadowingWithMultivariateDeclaration()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLetShadowingWithMultivariateDeclaration() {
	srcs := c.Srcs("var x, y; for (let x, y;;) {}")
	// normalized var decls outside the for loop
	expected := c.Expected("var x; var y; var x$jscomp$1 = void 0; var y$jscomp$1 = void 0; for (;;) {}")
	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testNonUniqueLet()
func (c *Es6RewriteBlockScopedDeclarationTest) TestNonUniqueLet() {
	srcs := c.Srcs(
		`function f() {
  var x = 1;
  if (a) {
    let x = 2;
    assert(x === 2);
  }
  if (b) {
    let x;
    assert(x === void 0);
  }
  assert(x === 1);
}
`)
	expected := c.Expected(
		`function f() {
  var x = 1;
  if (a) {
    var x$jscomp$1 = 2;
    assert(x$jscomp$1 === 2);
  }
  if (b) {
    var x$jscomp$2;
    assert(x$jscomp$2 === void 0);
  }
  assert(x === 1);
}
`)
	c.Test(srcs, expected)

	srcs =
		c.Srcs(
			`function f() {
  if (a) {
    let x = 2;
    assert(x === 2);
    if (b) {
      let x;
      assert(x === void 0);
    }
  }
}
`)
	expected =
		c.Expected(
			`function f() {
  if (a) {
    var x = 2;
    assert(x === 2);
    if (b) {
      var x$jscomp$1;
      assert(x$jscomp$1 === void 0);
    }
  }
}
`)
	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testRenameConflict()
func (c *Es6RewriteBlockScopedDeclarationTest) TestRenameConflict() {
	srcs := c.Srcs(
		`function f() {
  let x = 1;
  let x$0 = 2;
  {
    let x = 3;
  }
}
`)
	expected := c.Expected(
		`function f() {
  var x = 1;
  var x$0 = 2;
  {
    var x$jscomp$1 = 3;
  }
}
`)
	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testVanillaFor_letInitializer_getsNormalized()
func (c *Es6RewriteBlockScopedDeclarationTest) TestVanillaFor_letInitializer_getsNormalized() {
	srcs := c.Srcs(
		`function f() {
  return function foo() {
     for (let i = 0; i < 10; i++) {
     }
  };
}
`)

	expected := c.Expected(
		`function f() {
  return function foo() {
     var i = 0; // normalized
     for (; i < 10; i++) {
     }
  };
}
`)

	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testVanillaFor_constInitializer_getsNormalized()
func (c *Es6RewriteBlockScopedDeclarationTest) TestVanillaFor_constInitializer_getsNormalized() {
	srcs := c.Srcs(
		`function f() {
  return function foo() {
// const initializer (allowed syntactically) can not be updated in the incr.
     for (const i = 0; i < 10; /* no incr */) {
     }
  };
}
`)

	expected := c.Expected(
		`function f() {
  return function foo() {
     var i = 0; // normalized (outside the for) and no @const annotation after
// rewriting
     for (; i < 10; /* no incr */) {
     }
  };
}
`)

	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testForLetInitializer_declarationList_getsNormalized()
func (c *Es6RewriteBlockScopedDeclarationTest) TestForLetInitializer_declarationList_getsNormalized() {
	srcs := c.Srcs(
		`function f() {
  return function foo() {
     for (let i = 0, y = 0; i < 10; i++) {
     }
  };
}
`)

	expected := c.Expected(
		`function f() {
  return function foo() {
     var i = 0; var y=0; // normalized
     for (; i < 10; i++) {
     }
  };
}
`)

	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testsForVarInitializer_staysNormalized()
func (c *Es6RewriteBlockScopedDeclarationTest) TestsForVarInitializer_staysNormalized() {
	srcs := c.Srcs(
		`function f() {
  return function foo() {
     var i = 0;
     for ( ;i < 10; i++) { // originally normalized
     }
  };
}
`)

	c.TestSame(srcs)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testsForVarInitializer_unchanged()
func (c *Es6RewriteBlockScopedDeclarationTest) TestsForVarInitializer_unchanged() {
	srcs := c.Srcs(
		`function f() {
  return function foo() {
     for (var i = 0; i < 10; i++) { // originally unnormalized
     }
  };
}
`)
	expected := c.Expected(
		`function f() {
  return function foo() {
     var i = 0; // normalized; no rewriting change in this pass
     for (; i < 10; i++) {
     }
  };
}
`)
	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testForLoop()
func (c *Es6RewriteBlockScopedDeclarationTest) TestForLoop() {
	srcs := c.Srcs(
		`function use(x) {}
function f() {
  const y = 0;
  for (let x = 0; x < 10; x++) {
    const y = x * 2;
    const z = y;
  }
  use(y);
}
`)
	expected := c.Expected(
		`function use(x) {}
function f() {
  /** @const */ var y = 0;
  var x$jscomp$1 = 0;
  for (; x$jscomp$1 < 10; x$jscomp$1++) {
    /** @const */ var y$jscomp$1 = x$jscomp$1 * 2;
    /** @const */ var z = y$jscomp$1;
  }
  use(y);
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`for (let i in [0, 1]) {
  let f = function() {
    let i = 0;
    if (true) {
      let i = 1;
    }
  }
}
`)
	expected =
		c.Expected(
			`var i;
for (i in [0, 1]) {
  var f = function() {
    var i$jscomp$1 = 0;
    if (true) {
      var i$jscomp$2 = 1;
    }
  }
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs = c.Srcs("for (let i = 0;;) { let i; }")
	expected = c.Expected("var i = 0; for (;;) { var i$jscomp$1 = void 0; }")
	c.LoopClosureTest(srcs, expected)

	srcs = c.Srcs("for (let i = 0;;) {} let i;")
	expected = c.Expected("var i$jscomp$1 = 0; for (;;) {} var i;")
	c.LoopClosureTest(srcs, expected)

	c.Test(
		`for (var x in y) {
  /** @type {number} */
  let i;
}
`,
		`var x;
for (x in y) {
  var i = void 0;
}
`)

	srcs =
		c.Srcs(
			`for (const i in [0, 1]) {
  let f = function() {
    let i = 0;
    if (true) {
      let i = 1;
    }
  }
}
`)
	expected =
		c.Expected(
			`var i;
for (i in [0, 1]) {
  var f = function() {
    var i$jscomp$1 = 0;
    if (true) {
      var i$jscomp$2 = 1;
    }
  }
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLoopClosure()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLoopClosure() {
	srcs := c.Srcs(
		`const arr = [];
for (let i = 0; i < 10; i++) {
  arr.push(function() { return i; });
}
`)
	expected := c.Expected(
		`/** @const */ var arr = [];
var LOOP$0 = {};
LOOP$0.i = 0;
for (; LOOP$0.i < 10;
    LOOP$0 = {
      i: LOOP$0.i
    },
    LOOP$0.i++) {
  arr.push((function(LOOP$0$PARAM$1) {
      return function() { return LOOP$0$PARAM$1.i; };
  })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`const arr = [];
for (let i = 0; i < 10; i++) {
  let y = i;
  arr.push(function() { return y; });
}
`)
	expected =
		c.Expected(
			`/** @const */ var arr = [];
var LOOP$0 = {};
var i = 0;
for (; i < 10;
     LOOP$0 = {y: void 0},
     i++) {
  LOOP$0.y = i;
  arr.push((function(LOOP$0$PARAM$1) {
      return function() { return LOOP$0$PARAM$1.y; };
  })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`const arr = [];
while (true) {
  let i = 0;
  arr.push(function() { return i; });
}
`)
	expected =
		c.Expected(
			`/** @const */ var arr = [];
var LOOP$0 = {}
for (; true; LOOP$0 = {i:void 0}) {
  LOOP$0.i = 0;
  arr.push(
      (function(LOOP$0$PARAM$1) {
        return function() {
          return LOOP$0$PARAM$1.i;
        };
       })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`const arr = [];
for (let i = 0; i < 10; i++) {
  let y = i;
  arr.push(function() { return y + i; });
}
`)
	expected =
		c.Expected(
			`/** @const */ var arr = [];
var LOOP$0 = {};
LOOP$0.i = 0;
for (; LOOP$0.i < 10;
        LOOP$0 = {y: void 0, i: LOOP$0.i},
        LOOP$0.i++) {
  LOOP$0.y = LOOP$0.i;
  arr.push((function(LOOP$0$PARAM$1) {
    return function() {
      return LOOP$0$PARAM$1.y
        + LOOP$0$PARAM$1.i;
    };
  }(LOOP$0)));
}
`)
	c.LoopClosureTest(srcs, expected)

	// Renamed inner i
	srcs =
		c.Srcs(
			`const arr = [];
let x = 0
for (let i = 0; i < 10; i++) {
  let i = x + 1;
  arr.push(function() { return i + i; });
  x++;
}
`)
	expected =
		c.Expected(
			`/** @const */ var arr = [];
var x = 0
var LOOP$0 = {};
var i = 0;
for (; i < 10;
    LOOP$0 = {
      i$jscomp$1: void 0
    },
    i++) {
  LOOP$0.i$jscomp$1 = x + 1;
  arr.push((function(LOOP$0$PARAM$1) {
      return function() {
          return LOOP$0$PARAM$1.i$jscomp$1 + LOOP$0$PARAM$1.i$jscomp$1;
      };
  }(LOOP$0)));
  x++;
}
`)
	c.LoopClosureTest(srcs, expected)

	// Renamed, but both closures reference the inner i
	srcs =
		c.Srcs(
			`const arr = [];
let x = 0
for (let i = 0; i < 10; i++) {
  arr.push(function() { return i + i; });
  let i = x + 1;
  arr.push(function() { return i + i; });
  x++;
}
`)
	expected =
		c.Expected(
			`/** @const */ var arr = [];
var x = 0
var LOOP$0 = {};
var i = 0;
for (; i < 10;
    LOOP$0 = {
      i$jscomp$1: void 0
    }, i++) {
  arr.push((function(LOOP$0$PARAM$1) {
      return function() {
          return LOOP$0$PARAM$1.i$jscomp$1 + LOOP$0$PARAM$1.i$jscomp$1;
      };
  }(LOOP$0)));
  LOOP$0.i$jscomp$1 = x + 1;
  arr.push((function(LOOP$0$PARAM$2) {
      return function() {
          return LOOP$0$PARAM$2.i$jscomp$1 + LOOP$0$PARAM$2.i$jscomp$1;
      };
  }(LOOP$0)));
  x++;
}
`)
	c.LoopClosureTest(srcs, expected)

	// Renamed distinct captured variables
	srcs =
		c.Srcs(
			`for (let i = 0; i < 10; i++) {
  if (true) {
    let i = x - 1;
    arr.push(function() { return i + i; });
  }
  let i = x + 1;
  arr.push(function() { return i + i; });
  x++;
}
`)
	expected =
		c.Expected(
			`var LOOP$0 = {};
var i = 0;
for (; i < 10;
    LOOP$0 = {
      i$jscomp$2: void 0,
      i$jscomp$1: void 0
    }, i++) {
  if (true) {
    LOOP$0.i$jscomp$2 = x - 1;
    arr.push(
      (function(LOOP$0$PARAM$1) {
        return function() {
          return LOOP$0$PARAM$1.i$jscomp$2 + LOOP$0$PARAM$1.i$jscomp$2;
        };
      })(LOOP$0));
  }
  LOOP$0.i$jscomp$1 = x + 1;
  arr.push((function(LOOP$0$PARAM$2) {
      return function() {
         return LOOP$0$PARAM$2.i$jscomp$1 + LOOP$0$PARAM$2.i$jscomp$1;
      };
  })(LOOP$0));
  x++;
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs = c.Srcs("for (;;) { /** @type {number} */ let x = 3; var f = function() { return x; } }")
	expected =
		c.Expected(
			`var LOOP$0 = {};
for (;;LOOP$0 = {x: void 0}) {
  LOOP$0.x = 3;
  var f = function(LOOP$0$PARAM$1) {
    return function() { return LOOP$0$PARAM$1.x; }
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs = c.Srcs("for (;;) { let /** number */ x = 3; var f = function() { return x; } }")
	expected =
		c.Expected(
			`var LOOP$0 = {};
for (;; LOOP$0 = { x: void 0 }) {
  LOOP$0.x = 3;
  var f = function(LOOP$0$PARAM$1) {
    return function() { return LOOP$0$PARAM$1.x; }
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	// Preserve constancy
	srcs = c.Srcs("for (;;) { const /** number */ x = 3; var f = function() { return x; } }")
	expected =
		c.Expected(
			`var LOOP$0 = {};
for (;; LOOP$0 = { x: void 0 }) {
  /** @const */ LOOP$0.x = 3;
  var f = function(LOOP$0$PARAM$1) {
    return function() { return LOOP$0$PARAM$1.x}
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`for (;;) { let /** number */ x = 3, /** number */ y = 4;
var f = function() { return x + y; } }
`)
	expected =
		c.Expected(
			`var LOOP$0 = {};
for (;;LOOP$0 = {x: void 0, y: void 0}) {
  LOOP$0.x = 3;
  LOOP$0.y = 4;
  var f = function(LOOP$0$PARAM$1) {
    return function() { return LOOP$0$PARAM$1.x + LOOP$0$PARAM$1.y; };
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	// Preserve constancy on declaration lists
	srcs =
		c.Srcs(
			`for (;;) {
  const /** number */ x = 3, /** number */ y = 4;
  var f = function() { return x + y; };
}
`)
	expected =
		c.Expected(
			`var LOOP$0 = {};
for (;;LOOP$0 = { x: void 0, y: void 0 }) {
  /** @const */ LOOP$0.x = 3;
  /** @const */ LOOP$0.y = 4;
  var f = function(LOOP$0$PARAM$1) {
    return function() {
        return LOOP$0$PARAM$1.x + LOOP$0$PARAM$1.y;
    }
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	// No-op, vars don't need transpilation
	c.TestSame("for (;;) { var x = 3; var f = function() { return x; } }")

	srcs =
		c.Srcs(
			`var i;
for (i = 0;;) {
  let x = 0;
  var f = function() { x; };
}
`)
	expected =
		c.Expected(
			`var i;
i = 0;
var LOOP$0={};
for(;; LOOP$0 = { x: void 0 }) {
  LOOP$0.x = 0;
  var f = (function(LOOP$0$PARAM$1) {
    return function() { LOOP$0$PARAM$1.x; };
  })(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`for (foo();;) {
  let x = 0;
  var f = function() { x; };
}
`)
	expected =
		c.Expected(
			`foo();
var LOOP$0={};
for(;; LOOP$0 = { x: void 0 }) {
  LOOP$0.x = 0;
  var f = (function(LOOP$0$PARAM$1) {
    return function() { LOOP$0$PARAM$1.x; };
  })(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`for (function foo() {};;) {
  let x = 0;
  var f = function() { x; };
}
`)
	expected =
		c.Expected(
			`var LOOP$0={};
(function foo() {});
for(;; LOOP$0 = { x: void 0 }) {
  LOOP$0.x = 0;
  var f = (function(LOOP$0$PARAM$1) {
    return function() { LOOP$0$PARAM$1.x; };
  })(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`for (;;) {
  let x;
  foo(function() { return x; });
  x = 5;
}
`)
	expected =
		c.Expected(
			`var LOOP$0 = {};
for(;; LOOP$0 = {x: void 0}) {
  LOOP$0.x = void 0;
  foo(function(LOOP$0$PARAM$1) {
    return function() {
      return LOOP$0$PARAM$1.x;
    };
  }(LOOP$0));
  LOOP$0.x=5;
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLoopClosureWithContinue()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLoopClosureWithContinue() {
	// We must add a labeled block and convert continue statements to breaks to ensure that the
	// loop variable gets updated on every loop iteration of all loop types other than vanilla for.
	// For vanilla for(;;) loops we place the loop variable update in the loop update expression at
	// the top of the loop.

	// for-in case
	srcs := c.Srcs(
		`const obj = {a: 1, b: 2, c: 3, skipMe: 4};
const arr = [];
for (const p in obj) {
  if (p == 'skipMe') {
    continue;
  }
  arr.push(function() { return obj[p]; });
}
`)
	expected := c.Expected(
		`/** @const */ var obj = {a: 1, b: 2, c: 3, skipMe: 4};
/** @const */ var arr = [];
var LOOP$0 = {};
var p
for (p in obj) {
  LOOP$0 = {
    p: LOOP$0.p
  };
  LOOP$0.p = p;
  if (LOOP$0.p == 'skipMe') {
    continue;
  }
  arr.push(
      (function(LOOP$0$PARAM$1) {
        return function() {
          return obj[LOOP$0$PARAM$1.p];
        };
      })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`const values = ['a', 'b', 'c', 'skipMe'];
const arr = [];
while (values.length > 0) {
  const v = values.shift();
  if (v == 'skipMe') {
    continue;
  }
  arr.push(function() { return v; });
}
`)
	expected =
		c.Expected(
			`/** @const */ var values = ['a', 'b', 'c', 'skipMe'];
/** @const */ var arr = [];
var LOOP$0 = {};
for (; values.length > 0; LOOP$0 = {v:void 0}) {
  /** @const */ LOOP$0.v = values.shift();
  if (LOOP$0.v == 'skipMe') {
    continue
  }
  arr.push(
    (function(LOOP$0$PARAM$1) {
      return function() {
        return LOOP$0$PARAM$1.v;
      };
    })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`const values = ['a', 'b', 'c', 'skipMe'];
const arr = [];
do {
  const v = values.shift();
  if (v == 'skipMe') {
    continue;
  }
  arr.push(function() { return v; });
} while (values.length > 0);
`)
	expected =
		c.Expected(
			`/** @const */ var values = ['a', 'b', 'c', 'skipMe'];
/** @const */ var arr = [];
var LOOP$0 = {};
do {
  LOOP$0 = {
    v: void 0
  };
  /** @const */ LOOP$0.v = values.shift();
  if (LOOP$0.v == 'skipMe') {
    continue
  }
  arr.push(
      (function(LOOP$0$PARAM$1) {
        return function() {
          return LOOP$0$PARAM$1.v;
        };
      })(LOOP$0));
} while (values.length > 0);
`)
	c.LoopClosureTest(srcs, expected)

	// labeled continue case - TODO: b/197349249 Delete after b/197349249 is fixed
	srcs =
		c.Srcs(
			`const values = ['a', 'b', 'c', 'skipMe'];
const arr = [];
LABEL: while (values.length > 0) {
  const v = values.shift();
  if (v == 'skipMe') {
    continue LABEL;
  }
  arr.push(function() { return v; });
}
`)
	expected =
		c.Expected(
			`/** @const */ var values = ['a', 'b', 'c', 'skipMe'];
/** @const */ var arr = [];
var LOOP$0 = {};
LABEL: for (; values.length > 0; LOOP$0 = {v:void 0}) {
  /** @const */ LOOP$0.v = values.shift();
  if (LOOP$0.v == 'skipMe') {
    continue LABEL;
  }
  arr.push(
      (function(LOOP$0$PARAM$1) {
        return function() {
          return LOOP$0$PARAM$1.v;
        };
      })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)

	// nested labeled continue case - TODO: b/197349249 Delete after b/197349249 is fixed
	srcs =
		c.Srcs(
			`const values = ['abc', 'def', 'ghi', 'jkl'];
const words = [];
const letters = [];
OUTER: while (values.length > 0) {
  const v = values.shift();
  let i = 0;
  while (i < v.length) {
    const c = v.charAt(i);
    if (c == 'a') {
      i++;
      continue;
    } else if (c == 'i') {
      continue OUTER;
    } else {
      letters.push(function() { return c; });
    }
    i++;
  }
  words.push(function() { return v; });
}
`)
	expected =
		c.Expected(
			`/** @const */ var values = ['abc', 'def', 'ghi', 'jkl'];
/** @const */ var words = [];
/** @const */ var letters = [];
var LOOP$1 = {};
OUTER: for (; values.length > 0; LOOP$1 = {v:void 0}) {
  /** @const */ LOOP$1.v = values.shift();
  var i = 0;
  var LOOP$0 = {};
  for (; i < LOOP$1.v.length; LOOP$0 = {c:void 0}) {
      /** @const */ LOOP$0.c = LOOP$1.v.charAt(i);
      if (LOOP$0.c == 'a') {
        i++;
        continue;
      } else if (LOOP$0.c == 'i') {
        continue OUTER;
      } else {
        letters.push(
            (function(LOOP$0$PARAM$2) {
              return function() {
                return LOOP$0$PARAM$2.c;
              };
            })(LOOP$0));
      }
      i++;
  }
  words.push(
      (function(LOOP$1$PARAM$3) {
        return function() { return LOOP$1$PARAM$3.v; };
      })(LOOP$1));
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLoopClosureCommaInBody()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLoopClosureCommaInBody() {
	srcs := c.Srcs(
		`const foobar = [];
let j = 0;
for (let i = 0; i < 10; i++) {
  let i, j = 0;
  foobar.push(function() { return i + j; });
}
`)
	expected := c.Expected(
		`/** @const */
var foobar = [];
var j = 0;
var LOOP$0 = {};
var i = 0;
for (; i < 10;
    LOOP$0 = {
      i$jscomp$1: void 0,
      j$jscomp$1: void 0
    },
    i++) {
  LOOP$0.i$jscomp$1 = void 0;
  LOOP$0.j$jscomp$1 = 0;
  foobar.push(
      function(LOOP$0$PARAM$1) {
        return function() {
          return LOOP$0$PARAM$1.i$jscomp$1 + LOOP$0$PARAM$1.j$jscomp$1;
        };
      }(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLoopClosureCommaInIncrement()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLoopClosureCommaInIncrement() {
	srcs := c.Srcs(
		`const arr = [];
let j = 0;
for (let i = 0; i < 10; i++, j++) {
  arr.push(function() { return i + j; });
}
`)
	expected := c.Expected(
		`/** @const */ var arr = [];
var j = 0;
var LOOP$0 = {};
LOOP$0.i = 0;
for (; LOOP$0.i < 10;
    LOOP$0 = {
      i:LOOP$0.i
    }, (LOOP$0.i++, j++)) {
  arr.push((function(LOOP$0$PARAM$1) {
    return function() {
      return LOOP$0$PARAM$1.i + j;
    };
  })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLoopClosureCommaInInitializerAndIncrement()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLoopClosureCommaInInitializerAndIncrement() {
	srcs := c.Srcs(
		`const arr = [];
for (let i = 0, j = 0; i < 10; i++, j++) {
  arr.push(function() { return i + j; });
}
`)
	expected := c.Expected(
		`/** @const */ var arr = [];
var LOOP$0 = {};
LOOP$0.i = 0;
LOOP$0.j = 0;
for (; LOOP$0.i < 10;
    LOOP$0 = {
      i: LOOP$0.i,
      j: LOOP$0.j
    },
    (LOOP$0.i++, LOOP$0.j++)) {
  arr.push((function(LOOP$0$PARAM$1) {
      return function() {
         return LOOP$0$PARAM$1.i + LOOP$0$PARAM$1.j;
    };
  })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`const arr = [];
for (let i = 0, j = 0; i < 10; i++, j++) {
  arr.push(function() { return j; });
}
`)
	expected =
		c.Expected(
			`/** @const */ var arr = [];
var LOOP$0 = {};
var i = 0;
LOOP$0.j = 0;
for (; i < 10;
    LOOP$0 = {j :
    LOOP$0.j},
    (i++, LOOP$0.j++)) {
  arr.push((function(LOOP$0$PARAM$1) {
      return function() { return LOOP$0$PARAM$1.j; };
  })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLoopClosureMutated()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLoopClosureMutated() {
	srcs := c.Srcs(
		`const arr = [];
for (let i = 0; i < 10; i++) {
  arr.push(function() { return ++i; });
}
`)
	expected := c.Expected(
		`/** @const */ var arr = [];
var LOOP$0 = {};
LOOP$0.i = 0;
for (; LOOP$0.i < 10;
    LOOP$0 = {i: LOOP$0.i},
    LOOP$0.i++) {
  arr.push((function(LOOP$0$PARAM$1) {
// Is this right?
      return function() {
          return ++LOOP$0$PARAM$1.i;
      };
  }(LOOP$0)));
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`const arr = [];
for (let i = 0; i < 10; i++) {
  arr.push(function() { return i; });
  i += 100;
}
`)
	expected =
		c.Expected(
			`/** @const */ var arr = [];
var LOOP$0 = {};
LOOP$0.i = 0;
for (; LOOP$0.i < 10;
    LOOP$0 = {i: LOOP$0.i},
    LOOP$0.i++) {
  arr.push((function(LOOP$0$PARAM$1) {
      return function() {
          return LOOP$0$PARAM$1.i;
      };
  }(LOOP$0)));
  LOOP$0.i = LOOP$0.i + 100;
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLoopClosureWithNestedInnerFunctions()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLoopClosureWithNestedInnerFunctions() {
	srcs := c.Srcs(
		`for (let i = 0; i < 10; i++) {
  later(function(ctr) {
    (function() { return use(i); })();
  });
}
`)
	expected := c.Expected(
		`var LOOP$0 = {};
LOOP$0.i = 0;
for (; LOOP$0.i < 10;
    LOOP$0 = { i: LOOP$0.i },
    LOOP$0.i++) {
  later((function(LOOP$0$PARAM$1) {
    return function(ctr) {
      (function() { return use(LOOP$0$PARAM$1.i); })();
    };
  })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`for (let i = 0; i < 10; i++) {
  var f = function() {
    return function() {
      return i;
    };
  };
}
`)
	expected =
		c.Expected(
			`var LOOP$0 = {};
LOOP$0.i = 0;
for (; LOOP$0.i < 10;
    LOOP$0 = {i: LOOP$0.i},
    LOOP$0.i++) {
  var f = function(LOOP$0$PARAM$1) {
    return function() {
      return function() {
        return LOOP$0$PARAM$1.i;
      };
    };
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`use(function() {
  later(function(ctr) {
    for (let i = 0; i < 10; i++) {
      (function() { return use(i); })();
    }
  });
});
`)
	expected =
		c.Expected(
			`use(function() {
  later(function(ctr) {
    var LOOP$0 = {};
    LOOP$0.i = 0;
    for (; LOOP$0.i < 10;
        LOOP$0 = { i: LOOP$0.i },
        LOOP$0.i++) {
      (function(LOOP$0$PARAM$1) {
        return function() { return use(LOOP$0$PARAM$1.i); }
      })(LOOP$0)();
    }
  });
});
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testNestedLoop()
func (c *Es6RewriteBlockScopedDeclarationTest) TestNestedLoop() {
	srcs := c.Srcs(
		`function f() {
  let array = [];
  for (let i = 0; i < 10; i++) {
    for (let j = 0; j < 10; j++) {
      array.push(function() { return j++ + i++; });
      array.push(function() { return j++ + i++; });
    }
  }
}
`)
	expected := c.Expected(
		`function f() {
  var array = [];
  var LOOP$1 = {};
  LOOP$1.i = 0;
  for (; LOOP$1.i < 10;
      LOOP$1 = {
        i: LOOP$1.i
      },
      LOOP$1.i++) {
    var LOOP$0 = {};
    LOOP$0.j = 0;
    for (; LOOP$0.j < 10;
        LOOP$0 = {
          j: LOOP$0.j
        },
        LOOP$0.j++) {
      array.push((function(LOOP$0$PARAM$2, LOOP$1$PARAM$3) {
          return function() {
              return LOOP$0$PARAM$2.j++ + LOOP$1$PARAM$3.i++;
          };
      }(LOOP$0, LOOP$1)));
      array.push((function(LOOP$0$PARAM$4, LOOP$1$PARAM$5) {
          return function() {
              return LOOP$0$PARAM$4.j++ + LOOP$1$PARAM$5.i++;
          };
      }(LOOP$0, LOOP$1)));
    }
  }
}
`)
	c.LoopClosureTest(srcs, expected)

	// Renamed inner i
	srcs =
		c.Srcs(
			`function f() {
  let array = [];
  for (let i = 0; i < 10; i++) {
    array.push(function() { return i++ + i++; });
    for (let i = 0; i < 10; i++) {
      array.push(function() { return i++ + i++; });
    }
  }
}
`)
	expected =
		c.Expected(
			`function f() {
  var array = [];
  var LOOP$0 = {};
  LOOP$0.i = 0;
  for (; LOOP$0.i < 10;
      LOOP$0 = { i: LOOP$0.i },
      LOOP$0.i++) {
    array.push(
        (function(LOOP$0$PARAM$2) {
          return function() {
            return LOOP$0$PARAM$2.i++ + LOOP$0$PARAM$2.i++;
          };
        }(LOOP$0)));
    var LOOP$1 = {};
    LOOP$1.i$jscomp$1 = 0;
    for (; LOOP$1.i$jscomp$1 < 10;
        LOOP$1 = { i$jscomp$1: LOOP$1.i$jscomp$1 },
        LOOP$1.i$jscomp$1++) {
      array.push(
          (function(LOOP$1$PARAM$3) {
            return function() {
              return LOOP$1$PARAM$3.i$jscomp$1++ + LOOP$1$PARAM$3.i$jscomp$1++;
            };
          }(LOOP$1)));
    }
  }
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLabeledLoop()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLabeledLoop() {
	srcs := c.Srcs(
		`label1:
label2:
for (let x = 1;;) {
  let f = function() {
    return x;
  }
}
`)
	expected := c.Expected(
		`var LOOP$0 = {};
LOOP$0.x = 1;
label1:
label2:
for (;; LOOP$0 = { x: LOOP$0.x }) {
  var f = function(LOOP$0$PARAM$1) {
    return function() {
      return LOOP$0$PARAM$1.x;
    }
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testForIn()
func (c *Es6RewriteBlockScopedDeclarationTest) TestForIn() {
	srcs := c.Srcs(
		`const arr = [];
for (let i in [0, 1]) {
  arr.push(function() { return i; });
}
`)
	expected := c.Expected(
		`/** @const */ var arr = [];
var LOOP$0 = {};
var i;
for (i in [0, 1]) {
  LOOP$0 = {
    i: LOOP$0.i
  };
  LOOP$0.i = i;
  arr.push((function(LOOP$0$PARAM$1) {
      return function() { return LOOP$0$PARAM$1.i; };
  })(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			"for (;;) {\n  let a = getArray();\n  f = function() {\n// gets normalized first to `var x; for (x in use(a))` because of\n// `enableNormalize()` before the let/const rewriting pass runs.\n    for (var x in use(a)) {\n      f(a);\n      a.push(x);\n      return x;\n    }\n  }\n}\n")
	expected =
		c.Expected(
			`var LOOP$0 = {};
for (;; LOOP$0 = { a: void 0 }) {
  LOOP$0.a = getArray();
  f = (function(LOOP$0$PARAM$1) {
    return function() {
      var x;
      for (x in use(LOOP$0$PARAM$1.a)) {
        f(LOOP$0$PARAM$1.a);
        LOOP$0$PARAM$1.a.push(x);
        return x;
      }
    };
  }(LOOP$0));
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testDoWhileForInCapturedLet()
func (c *Es6RewriteBlockScopedDeclarationTest) TestDoWhileForInCapturedLet() {
	srcs := c.Srcs(
		`const arr = [];
do {
  let special = 99;
  const obj = { 0: 0, 1: 1, 2: special, 3: 3, 4: 4, 5: 5 };
  for (let i in obj) {
    i = Number(i);
    arr.push(function() { return i++; });
    arr.push(function() { return i + special; });
  }
} while (false);
`)
	expected := c.Expected(
		`/** @const */ var arr = [];
var LOOP$1 = {};
do {
  LOOP$1 = {
    special: void 0
  };
  LOOP$1.special = 99;
  /** @const */
  var obj =
    {
      0: 0,
      1: 1,
      2: LOOP$1.special,
      3: 3,
      4: 4,
      5: 5
    };
  var LOOP$0 = {};
  var i;
  for (i in obj) {
    LOOP$0 = {
      i: LOOP$0.i
    };
    LOOP$0.i = i
    LOOP$0.i = Number(LOOP$0.i);
    arr.push(
        (function(LOOP$0$PARAM$2) {
          return function() {
            return LOOP$0$PARAM$2.i++;
          };
        }(LOOP$0)));
    arr.push(
        (function(LOOP$0$PARAM$3, LOOP$1$PARAM$4) {
          return function() {
            return LOOP$0$PARAM$3.i + LOOP$1$PARAM$4.special;
          };
        }(LOOP$0, LOOP$1)));
  }
} while (false);
`)
	c.LoopClosureTest(srcs, expected)
}

// https://github.com/google/closure-compiler/issues/1124
//
// Java: Es6RewriteBlockScopedDeclarationTest.testFunctionsInLoop()
func (c *Es6RewriteBlockScopedDeclarationTest) TestFunctionsInLoop() {
	srcs := c.Srcs(
		`while (true) {
  let x = null;
  var f = function() {
    x();
  }
}
`)
	expected := c.Expected(
		`var LOOP$0 = {};
for (; true; LOOP$0 = {x:void 0}) {
  LOOP$0.x = null;
  var f = function(LOOP$0$PARAM$1) {
    return function() {
      (0,LOOP$0$PARAM$1.x)();
    };
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`while (true) {
  let x = null;
  let f = function() {
    x();
  }
}
`)
	expected =
		c.Expected(
			`var LOOP$0 = {};
for (; true; LOOP$0 = {x:void 0}) {
  LOOP$0.x = null;
  var f = function(LOOP$0$PARAM$1) {
    return function() {
      (0,LOOP$0$PARAM$1.x)();
    };
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`while (true) {
  let x = null;
  (function() {
    x();
  })();
}
`)
	expected =
		c.Expected(
			`var LOOP$0 = {};
for (; true; LOOP$0 = {x:void 0}) {
  LOOP$0.x = null;
  (function(LOOP$0$PARAM$1) {
    return function() {
      (0,LOOP$0$PARAM$1.x)();
    };
  })(LOOP$0)();
}
`)
	c.LoopClosureTest(srcs, expected)

	// NOTE: if the function itself is closed over then it must be included in the loop object
	srcs =
		c.Srcs(
			`while (true) {
  let x = 1;
  let f = function(i) {
    if (i < x) f(i + 1);
  };
  use(f);
}
`)
	expected =
		c.Expected(
			`var LOOP$0 = {};
for (; true; LOOP$0 = {x:void 0, f:void 0}) {
  LOOP$0.x = 1;
  LOOP$0.f = (function(LOOP$0$PARAM$1) {
    return function(i) {
      if (i < LOOP$0$PARAM$1.x) {
        (0,LOOP$0$PARAM$1.f)(i + 1);
      }
    };
  })(LOOP$0);
  use(LOOP$0.f);
}
`)
	c.LoopClosureTest(srcs, expected)
}

// https://github.com/google/closure-compiler/issues/1557
//
// Java: Es6RewriteBlockScopedDeclarationTest.testNormalizeDeclarations()
func (c *Es6RewriteBlockScopedDeclarationTest) TestNormalizeDeclarations() {
	srcs := c.Srcs(
		`while(true) {
  let x, y;
  let f = function() {
    x = 1;
    y = 2;
  }
}
`)
	expected := c.Expected(
		`var LOOP$0 = {};
for (; true; LOOP$0 = {x:void 0, y:void 0}) {
  LOOP$0.x = void 0;
  LOOP$0.y = void 0;
  var f = function(LOOP$0$PARAM$1) {
    return function() {
      LOOP$0$PARAM$1.x = 1;
      LOOP$0$PARAM$1.y = 2;
    }
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)

	srcs =
		c.Srcs(
			`while(true) {
  let x, y;
  let f = function() {
    y = 2;
    x = 1;
  }
}
`)
	expected =
		c.Expected(
			`var LOOP$0 = {};
for (; true; LOOP$0 = {y:void 0, x:void 0}) {
  LOOP$0.x = void 0;
  LOOP$0.y = void 0;
  var f = function(LOOP$0$PARAM$1) {
    return function() {
      LOOP$0$PARAM$1.y = 2;
      LOOP$0$PARAM$1.x = 1;
    }
  }(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testTypeAnnotationsOnLetConst()
func (c *Es6RewriteBlockScopedDeclarationTest) TestTypeAnnotationsOnLetConst() {
	mismatch := c.Warning(TypeValidator_TYPE_MISMATCH_WARNING)

	c.Test(c.Srcs("/** @type {number} */ let x = 5; x = 'str';"), mismatch)
	c.Test(c.Srcs("let /** number */ x = 5; x = 'str';"), mismatch)
	c.Test(c.Srcs("let /** @type {number} */ x = 5; x = 'str';"), mismatch)

	c.Test(c.Srcs("/** @type {number} */ const x = 'str';"), mismatch)
	c.Test(c.Srcs("const /** number */ x = 'str';"), mismatch)
	c.Test(c.Srcs("const /** @type {number} */ x = 'str';"), mismatch)
	c.Test(c.Srcs("const /** @type {string} */ x = 3, /** @type {number} */ y = 3;"), mismatch)
	c.Test(c.Srcs("const /** @type {string} */ x = 'str', /** @type {string} */ y = 3;"), mismatch)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testDoWhileForOfCapturedLetAnnotated()
func (c *Es6RewriteBlockScopedDeclarationTest) TestDoWhileForOfCapturedLetAnnotated() {
	c.Test(
		`while (true) {
  /** @type {number} */ let x = 5;
  (function() { x++; })();
  x = 7;
}
`,
		(*string)(nil))

	c.Test(
		`for (/** @type {number} */ let x = 5;;) {
  (function() { x++; })();
  x = 7;
}
`,
		(*string)(nil))

	// TODO(sdh): NTI does not detect the type mismatch in the transpiled code,
	// since the LOOP$0 object does not have its type inferred until after
	// the mistyped assignment.
	c.Test(
		c.Srcs(
			`while (true) {
  /** @type {number} */ let x = 5;
  (function() { x++; })();
  x = 'str';
}
`),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))

	c.Test(
		c.Srcs(
			`for (/** @type {number} */ let x = 5;;) {
  (function() { x++; })();
  x = 'str';
}
`),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLetForInitializers()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLetForInitializers() {
	srcs := c.Srcs(
		`{
  let l = [];
  for (var vx = 1, vy = 2, vz = 3; vx < 10; vx++) {
    let lx = vx, ly = vy, lz = vz;
    l.push(function() { return [ lx, ly, lz ]; });
  }
}
`)
	expected := c.Expected(
		`{
  var l = [];
  var vx = 1; var vy = 2; var vz = 3;
  var LOOP$0 = {};
  for (; vx < 10;
      LOOP$0 = {
        lx: void 0,
        ly: void 0,
        lz: void 0
      },
      vx++){
    LOOP$0.lx = vx;
    LOOP$0.ly = vy;
    LOOP$0.lz = vz;
    l.push(function(LOOP$0$PARAM$1) {
        return function() {
            return [
                LOOP$0$PARAM$1.lx,
                LOOP$0$PARAM$1.ly,
                LOOP$0$PARAM$1.lz ];
        };
    }(LOOP$0));
  }
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testCatch()
func (c *Es6RewriteBlockScopedDeclarationTest) TestCatch() {
	srcs := c.Srcs("function f(e) { try {} catch (e) { throw e; } }")
	expected := c.Expected("function f(e) { try {} catch (e$jscomp$1) { throw e$jscomp$1; } }")
	c.Test(srcs, expected)

	srcs =
		c.Srcs(
			`function f(e) {
  try {
    let f = function(e) {
      try {} catch (e) { e++; }
    }
  } catch (e) { e--; }
}
`)
	expected =
		c.Expected(
			`function f(e) {
  try {
    var f$jscomp$1 = function(e$jscomp$1) {
      try {} catch (e$jscomp$2) { e$jscomp$2++; }
    }
  } catch (e$jscomp$3) { e$jscomp$3--; }
}
`)
	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testReferenceToLoopScopedLetInObjectGetterAndSetter()
func (c *Es6RewriteBlockScopedDeclarationTest) TestReferenceToLoopScopedLetInObjectGetterAndSetter() {
	srcs :=
		c.Srcs(
			`for (let i = 0; i < 2; i++) {
   let bar = 42;
   let a = {
     get foo() {
      return bar
     },
     set foo(x) {
      use(bar);
     },
     prop: bar
   };
   bar = 43;
   use(a);
}
`)
	expected :=
		c.Expected(
			`var LOOP$0 = {};
var i = 0;
for (; i < 2;
    LOOP$0 = { bar: void 0 }, i++) {
  LOOP$0.bar = 42;
// Note that we wrap the entire object literal in an IIFE, because that's simpler
// than trying to individually wrap the getter and setter methods defined in it.
  var a =
      (function(LOOP$0$PARAM$2) {
        return {
          get foo() {
            return LOOP$0$PARAM$2.bar;
          },
          set foo(x) {
            use(LOOP$0$PARAM$2.bar);
          },
          prop: LOOP$0$PARAM$2.bar
      };
  })(LOOP$0);
  LOOP$0.bar = 43;
  use(a);
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testReferenceToMultipleLoopScopedLetConstInObjectWithGetter()
func (c *Es6RewriteBlockScopedDeclarationTest) TestReferenceToMultipleLoopScopedLetConstInObjectWithGetter() {
	srcs :=
		c.Srcs(
			`for (let i = 0; i < 2; i++) {
   let bar = 42;
   const baz = 43;
   let a = {
     get foo() {
      return bar + baz
     },
   };
}
`)
	expected :=
		c.Expected(
			`var LOOP$0 = {};
var i = 0;
for (; i < 2;
    LOOP$0 = { bar: void 0, baz: void 0 },
    i++) {
  LOOP$0.bar = 42;
  /** @const */
  LOOP$0.baz = 43;
  var a = (function(LOOP$0$PARAM$2) {
   return {
    get foo() {
       return LOOP$0$PARAM$2.bar +
           LOOP$0$PARAM$2.baz;
    }
  };
 })(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testReferenceToMultipleLoopScopedLetsInObjectWithSetter()
func (c *Es6RewriteBlockScopedDeclarationTest) TestReferenceToMultipleLoopScopedLetsInObjectWithSetter() {
	srcs :=
		c.Srcs(
			`for (let i = 0; i < 2; i++) {
   let bar = 42;
   let baz = 43;
   let a = {
     set foo(x =  bar) {
      return x + baz
     },
   };
}
`)
	expected :=
		c.Expected(
			`var LOOP$0 = {};
var i = 0;
for (; i < 2;
    LOOP$0 = {
       bar: void 0,
       baz: void 0
    },
    i++) {
  LOOP$0.bar = 42;
  LOOP$0.baz = 43;
  var a =
      (function(LOOP$0$PARAM$2) {
        return {
          set foo(x = LOOP$0$PARAM$2.bar) {
            return x + LOOP$0$PARAM$2.baz;
          }
        };
      })(LOOP$0);
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testUniqueIdAcrossMultipleFiles()
func (c *Es6RewriteBlockScopedDeclarationTest) TestUniqueIdAcrossMultipleFiles() {
	srcs := c.Srcs(
		`for (const i in [0, 1]) {
  let f = function() {
    let i = 0;
    if (true) {
      let i = 1;
    }
  }
}
`,
		`for (const b in [0, 1]) {
  let f = function() {
    let b = 0;
    if (true) {
      let b = 1;
    }
  }
}
`)
	expected := c.Expected(
		`var i;
for (i in [0, 1]) {
  var f = function() {
    var i$jscomp$1 = 0;
    if (true) {
      var i$jscomp$2 = 1;
    }
  }
}
`,
		`var b
for (b in[0, 1]) {
  var f$jscomp$1 = function() {
    var b$jscomp$1 = 0;
    if (true) {
      var b$jscomp$2 = 1;
    }
  };
}
`)
	c.Test(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testLoopClosure_overlappingLoopObjectNames()
func (c *Es6RewriteBlockScopedDeclarationTest) TestLoopClosure_overlappingLoopObjectNames() {
	srcs :=
		c.Srcs(
			`const arr = [];
for (let i0 = 0; i0 < 10; i0++) {
  arr.push(function() { return i0; });
  for (let i1 = 0; i1 < 10; i1++) {
    arr.push(function() { return i1; });
    for (let i2 = 0; i2 < 10; i2++) {
      arr.push(function() { return i2; });
      for (let i3 = 0; i3 < 10; i3++) {
        arr.push(function() { return i3; });
        for (let i4 = 0; i4 < 10; i4++) {
          arr.push(function() { return i4; });
          for (let i5 = 0; i5 < 10; i5++) {
            arr.push(function() { return i5; });
            for (let i6 = 0; i6 < 10; i6++) {
              arr.push(function() { return i6; });
              for (let i7 = 0; i7 < 10; i7++) {
                arr.push(function() { return i7; });
                for (let i8 = 0; i8 < 10; i8++) {
                  arr.push(function() { return i8; });
                  for (let i9 = 0; i9 < 10; i9++) {
                    arr.push(function() { return i9; });
                    for (let i10 = 0; i10 < 10; i10++) {
                      arr.push(function() { return i1 + i10; });
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`)
	expected :=
		c.Expected(
			`/** @const */ var arr = [];
var LOOP$0 = {};
LOOP$0.i0 = 0;
for (; LOOP$0.i0 < 10;
    LOOP$0 = {i0: LOOP$0.i0},
    LOOP$0.i0++) {
  arr.push((function(LOOP$0$PARAM$11) {
    return function() { return LOOP$0$PARAM$11.i0; };
  })(LOOP$0));
  var LOOP$1 = {};
  LOOP$1.i1 = 0;
  for (; LOOP$1.i1 < 10;
      LOOP$1 = {i1: LOOP$1.i1},
      LOOP$1.i1++) {
    arr.push((function(LOOP$1$PARAM$12) {
      return function() { return LOOP$1$PARAM$12.i1; };
    })(LOOP$1));
    var LOOP$2 = {};
    LOOP$2.i2 = 0;
    for (; LOOP$2.i2 < 10;
        LOOP$2 = {i2: LOOP$2.i2},
        LOOP$2.i2++) {
      arr.push((function(LOOP$2$PARAM$13) {
        return function() { return LOOP$2$PARAM$13.i2; };
      })(LOOP$2));
      var LOOP$3 = {};
      LOOP$3.i3 = 0;
      for (; LOOP$3.i3 < 10;
          LOOP$3 = {i3: LOOP$3.i3},
          LOOP$3.i3++) {
        arr.push((function(LOOP$3$PARAM$14) {
          return function() { return LOOP$3$PARAM$14.i3; };
        })(LOOP$3));
        var LOOP$4 = {};
        LOOP$4.i4 = 0;
        for (; LOOP$4.i4 < 10;
            LOOP$4 = {i4: LOOP$4.i4},
            LOOP$4.i4++) {
          arr.push((function(LOOP$4$PARAM$15) {
            return function() { return LOOP$4$PARAM$15.i4; };
          })(LOOP$4));
          var LOOP$5 = {};
          LOOP$5.i5 = 0;
          for (; LOOP$5.i5 < 10;
              LOOP$5 = {i5: LOOP$5.i5},
              LOOP$5.i5++) {
            arr.push((function(LOOP$5$PARAM$16) {
              return function() { return LOOP$5$PARAM$16.i5; };
            })(LOOP$5));
            var LOOP$6 = {};
            LOOP$6.i6 = 0;
            for (; LOOP$6.i6 < 10;
                LOOP$6 = {i6: LOOP$6.i6},
                LOOP$6.i6++) {
              arr.push((function(LOOP$6$PARAM$17) {
                return function() { return LOOP$6$PARAM$17.i6; };
              })(LOOP$6));
              var LOOP$7 = {};
              LOOP$7.i7 = 0;
              for (; LOOP$7.i7 < 10;
                  LOOP$7 = {i7: LOOP$7.i7},
                  LOOP$7.i7++) {
                arr.push((function(LOOP$7$PARAM$18) {
                  return function() { return LOOP$7$PARAM$18.i7; };
                })(LOOP$7));
                var LOOP$8 = {};
                LOOP$8.i8 = 0;
                for (; LOOP$8.i8 < 10;
                    LOOP$8 = {i8: LOOP$8.i8},
                    LOOP$8.i8++) {
                  arr.push((function(LOOP$8$PARAM$19) {
                    return function() { return LOOP$8$PARAM$19.i8; };
                  })(LOOP$8));
                  var LOOP$9 = {};
                  LOOP$9.i9 = 0;
                  for (; LOOP$9.i9 < 10;
                      LOOP$9 = {i9: LOOP$9.i9},
                      LOOP$9.i9++) {
                    arr.push((function(LOOP$9$PARAM$20) {
                      return function() { return LOOP$9$PARAM$20.i9; };
                    })(LOOP$9));
                    var LOOP$10 = {};
                    LOOP$10.i10 = 0;
                    for (; LOOP$10.i10 < 10;
                        LOOP$10 = {i10: LOOP$10.i10},
                        LOOP$10.i10++) {
                      arr.push((function(LOOP$1$PARAM$21, LOOP$10$PARAM$22) {
                        return function() { return LOOP$1$PARAM$21.i1 + LOOP$10$PARAM$22.i10; };
                      })(LOOP$1, LOOP$10));
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`)
	c.LoopClosureTest(srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testReferenceToLoopScopedLetInObjectWithThisAndArguments()
func (c *Es6RewriteBlockScopedDeclarationTest) TestReferenceToLoopScopedLetInObjectWithThisAndArguments() {
	srcs :=
		c.Srcs(
			`function h(userInput) {
  var out = [];
  for (let i = 0; i < 3; i++) {
    out.push({
      get idx() { return i; },
      ctx: this.policy,
      raw: arguments[0]
    });
  }
  return out;
}
`)
	expected :=
		c.Expected(
			`function h(userInput) {
  /** @const */
  var $ARGUMENTS$1 = arguments;
  /** @const */
  var $THIS$1 = this;
  var out = [];
  var LOOP$0 = {};
  LOOP$0.i = 0;
  for (; LOOP$0.i < 3;
      LOOP$0 = { i: LOOP$0.i },
      LOOP$0.i++) {
    out.push((function(LOOP$0$PARAM$2) {
      return {
        get idx() {
          return LOOP$0$PARAM$2.i;
        },
        ctx: $THIS$1.policy,
        raw: $ARGUMENTS$1[0]
      };
    })(LOOP$0));
  }
  return out;
}
`)
	c.LoopClosureTestWithExterns(c.Externs("var policy;"), srcs, expected)
}

// Java: Es6RewriteBlockScopedDeclarationTest.testTopLevelLet_shadowedInBareBlock()
func (c *Es6RewriteBlockScopedDeclarationTest) TestTopLevelLet_shadowedInBareBlock() {
	c.Test(
		c.Srcs("let a = 1; { let a = 2; a; } a;"),
		c.Expected("var a = 1; { var a$jscomp$1 = 2; a$jscomp$1; } a;"))
}

// Port of com/google/javascript/jscomp/Es6RewriteDestructuringTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/Es6RewriteDestructuringTest.java)

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

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteDestructuringTest has the test cases for Es6RewriteDestructuring.
type Es6RewriteDestructuringTest struct {
	CompilerTestCaseBase

	destructuringRewriteMode Es6RewriteDestructuring_ObjectDestructuringRewriteMode
}

// Java: Es6RewriteDestructuringTest.EXTERNS_BASE
var Es6RewriteDestructuringTest_EXTERNS_BASE = NewTestExternsBuilder().AddObject().AddArguments().AddString().AddJSCompLibraries().Build()

// Java: Es6RewriteDestructuringTest()
func NewEs6RewriteDestructuringTest() *Es6RewriteDestructuringTest {
	c := &Es6RewriteDestructuringTest{}
	c.destructuringRewriteMode = Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS
	c.InitCompilerTestCaseBaseWithExterns(c, Es6RewriteDestructuringTest_EXTERNS_BASE)
	return c
}

func TestEs6RewriteDestructuringTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteDestructuringTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *Es6RewriteDestructuringTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: Es6RewriteDestructuringTest.getProcessor(Compiler)
func (c *Es6RewriteDestructuringTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"injectTranspilationRuntimeLibraries", func(abstractCompiler AbstractCompiler) CompilerPass {
				return NewInjectTranspilationRuntimeLibraries(abstractCompiler)
			}))
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"es6RewriteDestructuring",
			func(abstractCompiler AbstractCompiler) CompilerPass {
				return NewEs6RewriteDestructuring_Builder(abstractCompiler).
					SetDestructuringRewriteMode(c.destructuringRewriteMode).
					Build()
			}))
	return optimizer
}

// Java: Es6RewriteDestructuringTest.customSetUp()
func (c *Es6RewriteDestructuringTest) CustomSetUp() {
	c.EnableNormalize()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2018)
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()

	// there are a lot of 'property x never defined on ?' warnings caused by object destructuring
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
}

// Java: Es6RewriteDestructuringTest.tearDown()
func (c *Es6RewriteDestructuringTest) TearDown() {
	c.CompilerTestCaseBase.TearDown()
	c.destructuringRewriteMode = Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS
}

// Java: Es6RewriteDestructuringTest.getOptions()
func (c *Es6RewriteDestructuringTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	return options
}

// Java: Es6RewriteDestructuringTest.testForLoopInitializer_let_accessedInClosure()
func (c *Es6RewriteDestructuringTest) TestForLoopInitializer_let_accessedInClosure() {
	c.Test(
		`const aClosures = [];
const bClosures = [];
for (let a = 0, [b, c] = [a, 3]; b < c; a--, b++) {
  aClosures.push(() => a);
  bClosures.push(() => b);
}
`,
		`const aClosures = [];
const bClosures = [];
for (let a = 0, b, c,
         $jscomp$destructuring$var0$unused = ($jscomp$destructuring$var1 => {
           var $jscomp$destructuring$var2 =
               (0, $jscomp.makeIterator)($jscomp$destructuring$var1);
           b = $jscomp$destructuring$var2.next().value;
           c = $jscomp$destructuring$var2.next().value;
           return $jscomp$destructuring$var1;
         })([a, 3]);
     b < c; a--, b++) {
  aClosures.push(() => {
    return a;
  });
  bClosures.push(() => {
    return b;
  });
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuring_forLoopInitializer_doesNotCrash()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuring_forLoopInitializer_doesNotCrash() {
	c.Test(
		`function isDeclaredInLoop(path) {
  for (let {
      parentPath,
      key
    } = path;
    parentPath; ({
      parentPath,
      key
    } = parentPath)) {
    return isDeclaredInLoop(parentPath);
  }
  return false;
}
`,
		`function isDeclaredInLoop(path) {
  for (let parentPath, key,
       $jscomp$destructuring$var0$unused = ($jscomp$destructuring$var1 => {
         var $jscomp$destructuring$var2 = $jscomp$destructuring$var1;
         parentPath = $jscomp$destructuring$var2.parentPath;
         key = $jscomp$destructuring$var2.key;
         return $jscomp$destructuring$var1;
       })(path);
       parentPath; ($jscomp$destructuring$var3 => {
         var $jscomp$destructuring$var4 = $jscomp$destructuring$var3;
         parentPath = $jscomp$destructuring$var4.parentPath;
         key = $jscomp$destructuring$var4.key;
         return $jscomp$destructuring$var3;
       })(parentPath)) {
    return isDeclaredInLoop(parentPath);
  }
  return false;
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuring()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuring() {
	c.Test(
		"var {a: b, c: d} = foo();",
		`var b; var d;
var $jscomp$destructuring$var0 = foo();
b = $jscomp$destructuring$var0.a;
d = $jscomp$destructuring$var0.c;
`)

	c.Test(
		"var {a,b} = foo();",
		`var a; var b;
var $jscomp$destructuring$var0 = foo();
a = $jscomp$destructuring$var0.a;
b = $jscomp$destructuring$var0.b;
`)

	c.Test(
		"let {a,b} = foo();",
		`var $jscomp$destructuring$var0 = foo();
let a = $jscomp$destructuring$var0.a;
let b = $jscomp$destructuring$var0.b;
`)

	c.Test(
		"const {a,b} = foo();",
		`/** @const */ var $jscomp$destructuring$var0 = foo();
const a = $jscomp$destructuring$var0.a;
const b = $jscomp$destructuring$var0.b;
`)

	c.Test(
		"var x; ({a: x} = foo());",
		`var x;
var $jscomp$destructuring$var0 = foo();
x = $jscomp$destructuring$var0.a;
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuringWithInitializer()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuringWithInitializer() {
	c.Test(
		"var {a : b = 'default'} = foo();",
		`var b;
var $jscomp$destructuring$var0 = foo();
b = ($jscomp$destructuring$var0.a === void 0) ?
    'default' :
    $jscomp$destructuring$var0.a
`)

	c.Test(
		"var {a = 'default'} = foo();",
		`var a;
var $jscomp$destructuring$var0 = foo();
a = ($jscomp$destructuring$var0.a === void 0) ?
    'default' :
    $jscomp$destructuring$var0.a
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuring_inVanillaForInitialize_var()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuring_inVanillaForInitialize_var() {
	c.Test(
		`A: B: for (
    var z = 1, {x: {a: b = 'default'}} = foo(), x = 0, {c} = bar();
    true;
    c++) {
}
`,
		`var z = 1;
var b;
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.x;
b = $jscomp$destructuring$var1.a === void 0 ?
    'default' : $jscomp$destructuring$var1.a;

var x = 0;
var c;
var $jscomp$destructuring$var2 = bar();
c = $jscomp$destructuring$var2.c;
A: B: for (; true; c++) {}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuring_inVanillaForInitialize_const()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuring_inVanillaForInitialize_const() {
	c.Test(
		`A: B: for (
    const z = 1, {x: {a: b = 'default'}} = foo(), x = 0, {c} = bar();
    true;
    c++) {
}
`,
		`{
  const z = 1;
  /** @const */ var $jscomp$destructuring$var0 = foo();
  /** @const */ var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.x;
  const b = $jscomp$destructuring$var1.a === void 0 ?
      'default' : $jscomp$destructuring$var1.a;
  {
    const x = 0;
    /** @const */ var $jscomp$destructuring$var2 = bar();
    const c = $jscomp$destructuring$var2.c;
    A: B: for (; true; c++) {}
  }
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuring_inVanillaForInitialize_let()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuring_inVanillaForInitialize_let() {
	c.Test(
		`A: B: for (
    let z = 1, {x: {a: b = 'default'}} = foo(), x = 0, {c} = bar();
    true;
    c++) {
}
`,
		`A: B: for (let z = 1, b,
           $jscomp$destructuring$var0$unused = ($jscomp$destructuring$var1 => {
             var $jscomp$destructuring$var2 = $jscomp$destructuring$var1;
             var $jscomp$destructuring$var3 = $jscomp$destructuring$var2.x;
             b = $jscomp$destructuring$var3.a === void 0 ?
                 'default' :
                 $jscomp$destructuring$var3.a;
             return $jscomp$destructuring$var1;
           })(foo()),
           x = 0, c,
           $jscomp$destructuring$var4$unused = ($jscomp$destructuring$var5 => {
             var $jscomp$destructuring$var6 = $jscomp$destructuring$var5;
             c = $jscomp$destructuring$var6.c;
             return $jscomp$destructuring$var5;
           })(bar());
           true; c++) {}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuringNested()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuringNested() {
	c.Test(
		"var {a: {b}} = foo();",
		`var b;
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.a;
b = $jscomp$destructuring$var1.b
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuringComputedProps()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuringComputedProps() {
	c.Test(
		"var {[a]: b} = foo();",
		"var b; var $jscomp$destructuring$var0 = foo(); b = $jscomp$destructuring$var0[a];")

	c.Test(
		"({[a]: b} = foo());",
		"var $jscomp$destructuring$var0 = foo(); b = $jscomp$destructuring$var0[a];")

	c.Test(
		"var {[foo()]: x = 5} = {};",
		`var x;
var $jscomp$destructuring$var0 = {};
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0[foo()];
x = $jscomp$destructuring$var1 === void 0 ?
    5 : $jscomp$destructuring$var1
`)

	c.Test(
		"function f({['KEY']: x}) {}",
		`function f($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0
  x = $jscomp$destructuring$var1['KEY']
}
`)
}

// Java: Es6RewriteDestructuringTest.testGithubIssue2189()
func (c *Es6RewriteDestructuringTest) TestGithubIssue2189() {
	c.SetExpectParseWarningsInThisTest()
	c.Test(
		`/**
 * @param {string} a
 * @param {{b: ?<?>}} __1
 */
function x(a, { b }) {}
`,
		`function x(a, $jscomp$destructuring$var0) {
  var b;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  b=$jscomp$destructuring$var1.b;
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuringStrangeProperties()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuringStrangeProperties() {
	c.Test(
		"var {5: b} = foo();",
		"var b; var $jscomp$destructuring$var0 = foo(); b = $jscomp$destructuring$var0['5']")

	c.Test(
		"var {0.1: b} = foo();",
		`var b;
var $jscomp$destructuring$var0 = foo();
b = $jscomp$destructuring$var0['0.1']
`)

	c.Test(
		"var {'str': b} = foo();",
		`var b;
var $jscomp$destructuring$var0 = foo();
b = $jscomp$destructuring$var0['str']
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuringFunction()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuringFunction() {
	c.Test(
		"function f({a: b}) {}",
		`function f($jscomp$destructuring$var0) {
  var b;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0
  b = $jscomp$destructuring$var1.a
}
`)

	c.Test(
		"function f({a}) {}",
		`function f($jscomp$destructuring$var0) {
  var a;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0
  a = $jscomp$destructuring$var1.a
}
`)

	c.Test(
		"function f({k: {subkey : a}}) {}",
		`function f($jscomp$destructuring$var0) {
  var a;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0
  var $jscomp$destructuring$var2 = $jscomp$destructuring$var1.k;
  a = $jscomp$destructuring$var2.subkey;
}
`)

	c.Test(
		"function f({k: [x, y, z]}) {}",
		`function f($jscomp$destructuring$var0) {
  var x; var y; var z;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0
  var $jscomp$destructuring$var2 =
      (0, $jscomp.makeIterator)($jscomp$destructuring$var1.k);
  x = $jscomp$destructuring$var2.next().value;
  y = $jscomp$destructuring$var2.next().value;
  z = $jscomp$destructuring$var2.next().value;
}
`)

	c.Test(
		"function f({key: x = 5}) {}",
		`function f($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0
  x = $jscomp$destructuring$var1.key === void 0 ?
      5 : $jscomp$destructuring$var1.key
}
`)

	c.Test(
		"function f({[key]: x = 5}) {}",
		`function f($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0
  var $jscomp$destructuring$var2 = $jscomp$destructuring$var1[key]
  x = $jscomp$destructuring$var2 === void 0 ?
      5 : $jscomp$destructuring$var2
}
`)

	c.Test(
		"function f({x = 5}) {}",
		`function f($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0
  x = $jscomp$destructuring$var1.x === void 0 ?
      5 : $jscomp$destructuring$var1.x
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuringFunctionBadJsdoc()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuringFunctionBadJsdoc() {
	// see https://github.com/google/closure-compiler/issues/3175
	c.SetExpectParseWarningsInThisTest() // intentionally pass bad JSDoc
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`/**
 * @param {{foo: string[]}} obj
 * @param {string} id
 */
function f({foo}, id) {}
`,
		`function f($jscomp$destructuring$var0, id) {
  var foo;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  foo = $jscomp$destructuring$var1.foo;
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuringFunctionJsDoc()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuringFunctionJsDoc() {
	c.Test(
		"function f(/** {x: number, y: number} */ {x, y}) {}",
		`function f($jscomp$destructuring$var0) {
  var x; var y;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  x = $jscomp$destructuring$var1.x;
  y = $jscomp$destructuring$var1.y;
}
`)
}

// Java: Es6RewriteDestructuringTest.testDefaultParametersDestructuring()
func (c *Es6RewriteDestructuringTest) TestDefaultParametersDestructuring() {
	c.Test(
		"function f({a,b} = foo()) {}",
		`function f($jscomp$destructuring$var0){
  var a; var b;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0=== void 0 ?
    foo() : $jscomp$destructuring$var0;
  a = $jscomp$destructuring$var1.a;
  b = $jscomp$destructuring$var1.b;
}
`)
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuring()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuring() {
	c.Test(
		"var [x,y] = z();",
		`var x;
var y;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(z());
x = $jscomp$destructuring$var0.next().value;
y = $jscomp$destructuring$var0.next().value;
`)

	c.Test(
		`var x,y;
[x,y] = z();
`,
		`var x;
var y;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(z());
x = $jscomp$destructuring$var0.next().value;
y = $jscomp$destructuring$var0.next().value;
`)

	c.Test(
		`var [a,b] = c();
var [x,y] = z();
`,
		`var a;
var b;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(c());
a = $jscomp$destructuring$var0.next().value;
b = $jscomp$destructuring$var0.next().value;
var x;
var y;
var $jscomp$destructuring$var1 = (0, $jscomp.makeIterator)(z());
x = $jscomp$destructuring$var1.next().value;
y = $jscomp$destructuring$var1.next().value;
`)
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuringDefaultValues()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuringDefaultValues() {
	c.Test(
		"var a; [a=1] = b();",
		`var a;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(b())
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.next().value
a = ($jscomp$destructuring$var1 === void 0) ?
    1 :
    $jscomp$destructuring$var1;
`)

	c.Test(
		"var [a=1] = b();",
		`var a;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(b())
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.next().value
a = ($jscomp$destructuring$var1 === void 0) ?
    1 :
    $jscomp$destructuring$var1;
`)

	c.Test(
		"var [a, b=1, c] = d();",
		`var a;
var b;
var c;
var $jscomp$destructuring$var0=(0, $jscomp.makeIterator)(d());
a = $jscomp$destructuring$var0.next().value;
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.next().value;
b = ($jscomp$destructuring$var1 === void 0) ?
    1 :
    $jscomp$destructuring$var1;
c=$jscomp$destructuring$var0.next().value
`)

	c.Test(
		c.Srcs("var a; [[a] = ['b']] = [];"),
		c.Expected(
			`var a;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)([]);
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.next().value;
var $jscomp$destructuring$var2 = (0, $jscomp.makeIterator)(
    $jscomp$destructuring$var1 === void 0
        ? ['b']
        : $jscomp$destructuring$var1);
a = $jscomp$destructuring$var2.next().value
`))
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuringParam()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuringParam() {
	c.Test(
		"function f([x,y]) { use(x); use(y); }",
		`function f($jscomp$destructuring$var0) {
  var x; var y;
  var $jscomp$destructuring$var1 = (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
  x = $jscomp$destructuring$var1.next().value;
  y = $jscomp$destructuring$var1.next().value;
  use(x);
  use(y);
}
`)

	c.Test(
		"function f([x, , y]) { use(x); use(y); }",
		`function f($jscomp$destructuring$var0) {
  var x; var y;
  var $jscomp$destructuring$var1 = (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
  x = $jscomp$destructuring$var1.next().value;
  $jscomp$destructuring$var1.next();
  y = $jscomp$destructuring$var1.next().value;
  use(x);
  use(y);
}
`)
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuring_inVanillaForInitializer_var()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuring_inVanillaForInitializer_var() {
	c.Test(
		"A: B: for (var z = 1, [a = 'default'] = foo(), x = 0, [c] = bar(); true; c++) { }",
		`var z = 1;
var a;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(foo());
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.next().value;
a = $jscomp$destructuring$var1 === void 0 ?
    'default' : $jscomp$destructuring$var1;

var x = 0;
var c;
var $jscomp$destructuring$var2 = (0, $jscomp.makeIterator)(bar());
c = $jscomp$destructuring$var2.next().value;

A: B: for (; true; c++) {}
`)
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuring_inVanillaForInitializer_const()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuring_inVanillaForInitializer_const() {
	c.Test(
		"A: B: for (const z = 1, [a = 'default'] = foo(), x = 0, [c] = bar(); true; c++) { }",
		`{
  const z = 1;
  var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(foo());
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.next().value;
  const a = $jscomp$destructuring$var1 === void 0 ?
      'default' : $jscomp$destructuring$var1;
  {
    const x = 0;
    var $jscomp$destructuring$var2 = (0, $jscomp.makeIterator)(bar());
    const c = $jscomp$destructuring$var2.next().value;
    A: B: for (; true; c++) {}
  }
}
`)
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuring_inVanillaForInitializer_let()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuring_inVanillaForInitializer_let() {

	c.Test(
		"A: B: for (let z = 1, [a = 'default'] = foo(), x = 0, [c] = bar(); true; c++) { }",
		`A: B: for (let z = 1, a,
           $jscomp$destructuring$var0$unused = ($jscomp$destructuring$var1 => {
             var $jscomp$destructuring$var2 =
                 (0, $jscomp.makeIterator)($jscomp$destructuring$var1);
             var $jscomp$destructuring$var3 =
                 $jscomp$destructuring$var2.next().value;
             a = $jscomp$destructuring$var3 === void 0 ?
                 'default' :
                 $jscomp$destructuring$var3;
             return $jscomp$destructuring$var1;
           })(foo()),
           x = 0, c,
           $jscomp$destructuring$var4$unused = ($jscomp$destructuring$var5 => {
             var $jscomp$destructuring$var6 =
                 (0, $jscomp.makeIterator)($jscomp$destructuring$var5);
             c = $jscomp$destructuring$var6.next().value;
             return $jscomp$destructuring$var5;
           })(bar());
           true; c++) {}
`)
}

// Java: Es6RewriteDestructuringTest.testSimpleDestructuring_constAnnotationAdded()
func (c *Es6RewriteDestructuringTest) TestSimpleDestructuring_constAnnotationAdded() {

	c.Test(
		"function foo(...{length}) {}",
		// pull destructuring out of params first.
		`function foo(...$jscomp$destructuring$var0) {
  var length;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  length = $jscomp$destructuring$var1.length;
}
`)
	script := c.GetLastCompiler().GetJsRoot().GetFirstChild()
	jutil.CheckState(script.IsScript())
	func_ := script.GetFirstChild()
	jutil.CheckState(func_.IsFunction())
	block := func_.GetLastChild()
	jutil.CheckState(block.IsBlock())
	// `var $jscomp$destructuring$var1 =  ...`
	var1 := block.GetSecondChild().GetFirstChild()
	jutil.CheckState(var1.IsName())
	AssertThat(var1.GetString()).IsEqualTo("$jscomp$destructuring$var1")
	AssertThat(var1.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuringRest()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuringRest() {
	c.Test(
		"let [one, ...others] = f();",
		`var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(f());
let one = $jscomp$destructuring$var0.next().value;
let others = (0, $jscomp.arrayFromIterator)($jscomp$destructuring$var0);
`)

	c.Test(
		"function f([first, ...rest]) {}",
		`function f($jscomp$destructuring$var0) {
  var first; var rest;
  var $jscomp$destructuring$var1 = (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
  first = $jscomp$destructuring$var1.next().value;
  rest = (0, $jscomp.arrayFromIterator)($jscomp$destructuring$var1);
}
`)
}

// Java: Es6RewriteDestructuringTest.testRestParamDestructuring()
func (c *Es6RewriteDestructuringTest) TestRestParamDestructuring() {
	c.Test(
		"function f(first, ...[re, st, ...{length: num_left}]) {}",
		`function f(first, ...$jscomp$destructuring$var0) {
  var re; var st; var num_left;
  var $jscomp$destructuring$var1 = (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
  re = $jscomp$destructuring$var1.next().value;
  st = $jscomp$destructuring$var1.next().value;
  var $jscomp$destructuring$var2 = (0, $jscomp.arrayFromIterator)($jscomp$destructuring$var1);
  num_left = $jscomp$destructuring$var2.length;
}
`)
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuringMixedRest()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuringMixedRest() {
	c.Test(
		c.Srcs("let [first, ...[re, st, ...{length: num_left}]] = f();"),
		c.Expected(
			`var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(f());
let first = $jscomp$destructuring$var0.next().value;
var $jscomp$destructuring$var1 =
    (0, $jscomp.makeIterator)(
        (0, $jscomp.arrayFromIterator)($jscomp$destructuring$var0));
let re = $jscomp$destructuring$var1.next().value;
let st = $jscomp$destructuring$var1.next().value;
var $jscomp$destructuring$var2 =
    (0, $jscomp.arrayFromIterator)($jscomp$destructuring$var1);
let num_left = $jscomp$destructuring$var2.length;
`))
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuringArguments()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuringArguments() {
	c.Test(
		"function f() { var [x, y] = arguments; }",
		`function f() {
  var x;
  var y;
  var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(arguments);
  x = $jscomp$destructuring$var0.next().value;
  y = $jscomp$destructuring$var0.next().value;
}
`)
}

// Java: Es6RewriteDestructuringTest.testMixedDestructuring()
func (c *Es6RewriteDestructuringTest) TestMixedDestructuring() {
	c.Test(
		"var [a,{b,c}] = foo();",
		`var a;
var b;
var c;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(foo());
a = $jscomp$destructuring$var0.next().value;
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.next().value;
b = $jscomp$destructuring$var1.b;
c = $jscomp$destructuring$var1.c
`)

	c.Test(
		"var {a,b:[c,d]} = foo();",
		`var a;
var c;
var d;
var $jscomp$destructuring$var0 = foo();
a = $jscomp$destructuring$var0.a;
var $jscomp$destructuring$var1 = (0, $jscomp.makeIterator)($jscomp$destructuring$var0.b);
c = $jscomp$destructuring$var1.next().value;
d = $jscomp$destructuring$var1.next().value
`)
}

// Java: Es6RewriteDestructuringTest.testDestructuringForOf()
func (c *Es6RewriteDestructuringTest) TestDestructuringForOf() {
	c.Test(
		"for ({x} of y) { console.log(x); }",
		`var $jscomp$destructuring$var0;
for ($jscomp$destructuring$var0 of y) {
   var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
   x = $jscomp$destructuring$var1.x;
   console.log(x);
}
`)
}

// Java: Es6RewriteDestructuringTest.testDestructuringForOfWithShadowing()
func (c *Es6RewriteDestructuringTest) TestDestructuringForOfWithShadowing() {
	c.Test(
		c.Srcs("for (const [x] of []) { const y = 0; }"),
		c.Expected(
			`for (const $jscomp$destructuring$var0 of []) {
  var $jscomp$destructuring$var1 =
      (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
  const x = $jscomp$destructuring$var1.next().value;
  {
    const y = 0;
  }
}
`))
}

// Java: Es6RewriteDestructuringTest.testDestructuringForInWithShadowing()
func (c *Es6RewriteDestructuringTest) TestDestructuringForInWithShadowing() {
	c.Test(
		c.Srcs("for (const [x] in {}) { const y = 0; }"),
		c.Expected(
			`for (const $jscomp$destructuring$var0 in {}) {
  var $jscomp$destructuring$var1 =
      (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
  const x = $jscomp$destructuring$var1.next().value;
  {
    const y = 0;
  }
}
`))
}

// Java: Es6RewriteDestructuringTest.testDefaultValueInObjectPattern()
func (c *Es6RewriteDestructuringTest) TestDefaultValueInObjectPattern() {
	c.Test(
		"function f({x = a()}, y = b()) {}",
		`function f($jscomp$destructuring$var0, y) {
  var x;
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
x = $jscomp$destructuring$var1.x === void 0
       ? a() : $jscomp$destructuring$var1.x;
y = y === void 0 ? b() : y
}
`)
}

// Java: Es6RewriteDestructuringTest.testDefaultParameters()
func (c *Es6RewriteDestructuringTest) TestDefaultParameters() {
	c.EnableTypeCheck()

	c.Test(
		"function f(/** ? */ zero, /** ?= */ one = 1, /** ?= */ two = 2) {}; f(1); f(1,2,3);",
		`function f(zero, one, two) {
  one = (one === void 0) ? 1 : one;
  two = (two === void 0) ? 2 : two;
};
f(1); f(1,2,3);
`)

	c.Test(
		c.Srcs("function f(/** ? */ zero, /** ?= */ one = 1, /** ?= */ two = 2) {}; f();"),
		c.Expected(
			`function f(zero, one, two) {
  one = (one === void 0) ? 1 : one;
  two = (two === void 0) ? 2 : two;
}; f();
`),
		c.Warning(TypeCheck_WRONG_ARGUMENT_COUNT))
}

// Java: Es6RewriteDestructuringTest.testDefaultAndRestParameters()
func (c *Es6RewriteDestructuringTest) TestDefaultAndRestParameters() {
	c.Test(
		"function f(zero, one = 1, ...two) {}",
		`function f(zero, one, ...two) {
  one = (one === void 0) ? 1 : one;
}
`)

	c.Test(
		"function f(/** number= */ x = 5) {}",
		`function f(x) {
  x = (x === void 0) ? 5 : x;
}
`)
}

// Java: Es6RewriteDestructuringTest.testDefaultUndefinedParameters()
func (c *Es6RewriteDestructuringTest) TestDefaultUndefinedParameters() {
	c.EnableTypeCheck()

	c.Test("function f(zero, one=undefined) {}", "function f(zero, one) {}")

	c.Test("function f(zero, one=void 42) {}", "function f(zero, one) {}")

	c.Test("function f(zero, one=void(42)) {}", "function f(zero, one) {}")

	c.Test("function f(zero, one=void '\\x42') {}", "function f(zero, one) {}")

	c.Test(
		"function f(zero, one='undefined') {}",
		"function f(zero, one) {   one = (one === void 0) ? 'undefined' : one; }")

	c.Test(
		"function f(zero, one=void g()) {}",
		"function f(zero, one) {   one = (one === void 0) ? void g() : one; }")
}

// Java: Es6RewriteDestructuringTest.testCatch()
func (c *Es6RewriteDestructuringTest) TestCatch() {
	c.Test(
		"try {} catch ({message}) {}",
		`try {} catch($jscomp$destructuring$var0) {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  let message = $jscomp$destructuring$var1.message
}
`)
}

// Java: Es6RewriteDestructuringTest.testTypeCheck()
func (c *Es6RewriteDestructuringTest) TestTypeCheck() {
	c.EnableTypeCheck()

	c.Test(
		"/** @param {{x: number}} obj */ function f({x}) {}",
		`function f($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  x = $jscomp$destructuring$var1.x;
}
`)

	c.Test(
		c.Srcs(
			`/** @param {{x: number}} obj */
function f({x}) {}
f({ x: 'str'});
`),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))

	c.Test(
		`/** @param {{x: number}} obj */
var f = function({x}) {}
`,
		`var f = function($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  x = $jscomp$destructuring$var1.x;
}
`)

	c.Test(
		`/** @param {{x: number}} obj */
f = function({x}) {}
`,
		`f = function($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  x = $jscomp$destructuring$var1.x;
}
`)

	c.Test(
		`/** @param {{x: number}} obj */
ns.f = function({x}) {}
`,
		`ns.f = function($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  x = $jscomp$destructuring$var1.x;
}
`)

	c.Test(
		"ns.f = function({x} = {x: 0}) {};",
		`ns.f = function($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 =
      $jscomp$destructuring$var0 === void 0 ? {x:0} : $jscomp$destructuring$var0;
  x = $jscomp$destructuring$var1.x
};
`)

	c.Test(
		`/** @param {{x: number}=} obj */
ns.f = function({x} = {x: 0}) {};
`,
		`ns.f = function($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 =
      $jscomp$destructuring$var0=== void 0 ? {x:0} : $jscomp$destructuring$var0;
  x = $jscomp$destructuring$var1.x
};
`)
}

// Java: Es6RewriteDestructuringTest.testDestructuringPatternInExterns()
func (c *Es6RewriteDestructuringTest) TestDestructuringPatternInExterns() {
	c.EnableTypeCheck()
	c.AllowExternsChanges()

	c.TestSame(
		c.Externs(
			Es6RewriteDestructuringTest_EXTERNS_BASE+
				`/** @constructor */
function Foo() {}

Foo.prototype.bar = function({a}) {};
`),
		c.Srcs("(new Foo).bar({b: 0});"))
	// TODO(sdh): figure out what's going on here
}

// Java: Es6RewriteDestructuringTest.testTypeCheck_inlineAnnotations()
func (c *Es6RewriteDestructuringTest) TestTypeCheck_inlineAnnotations() {
	c.EnableTypeCheck()

	c.Test(
		"function f(/** {x: number} */ {x}) {}",
		`function f($jscomp$destructuring$var0) {
  var x;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  x = $jscomp$destructuring$var1.x;
}
`)

	c.Test(
		c.Srcs(
			`function f(/** {x: number} */ {x}) {}
f({ x: 'str'});
`),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))
}

// Java: Es6RewriteDestructuringTest.testDestructuringArrayNotInExprResult()
func (c *Es6RewriteDestructuringTest) TestDestructuringArrayNotInExprResult() {
	c.Test(
		c.Srcs(
			`var x, a, b;
x = ([a,b] = [1,2])
`),
		c.Expected(
			`var x;
var a;
var b;
x = ($jscomp$destructuring$var0 => {
  var $jscomp$destructuring$var1 =
      (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
  a = $jscomp$destructuring$var1.next().value;
  b = $jscomp$destructuring$var1.next().value;
  return $jscomp$destructuring$var0;
})([1, 2]);
`))

	c.Test(
		c.Srcs(
			`var foo = function () {
var x, a, b;
x = ([a,b] = [1,2]);
}
foo();
`),
		c.Expected(
			`var foo = function() {
  var x;
  var a;
  var b;
  x = ($jscomp$destructuring$var0 => {
    var $jscomp$destructuring$var1 =
        (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
    a = $jscomp$destructuring$var1.next().value;
    b = $jscomp$destructuring$var1.next().value;
    return $jscomp$destructuring$var0;
  })([1, 2]);
};
foo();
`))

	c.Test(
		c.Srcs(
			`var prefix;
for (;;[, prefix] = /** @type {!Array<string>} */ (/\.?([^.]+)$/.exec(prefix))){
}
`),
		c.Expected(
			`var prefix;
for (;; ($jscomp$destructuring$var0 => {
      var $jscomp$destructuring$var1 =
          (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
      $jscomp$destructuring$var1.next();
      prefix = $jscomp$destructuring$var1.next().value;
      return $jscomp$destructuring$var0;
    })(/\.?([^.]+)$/.exec(prefix))) {
}
`))

	c.Test(
		c.Srcs(
			`var prefix;
for (;;[, prefix] = /** @type {!Array<string>} */ (/\.?([^.]+)$/.exec(prefix))){
   console.log(prefix);
}
`),
		c.Expected(
			`var prefix;
for (;; ($jscomp$destructuring$var0 => {
      var $jscomp$destructuring$var1 =
          (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
      $jscomp$destructuring$var1.next();
      prefix = $jscomp$destructuring$var1.next().value;
      return $jscomp$destructuring$var0;
    })(/\.?([^.]+)$/.exec(prefix))) {
  console.log(prefix);
}
`))

	c.Test(
		c.Srcs(
			`var x = 1;
for (; x < 3; [x,] = [3,4]){
   console.log(x);
}
`),
		c.Expected(
			`var x = 1;
for (; x < 3; (($jscomp$destructuring$var0) => {
   var $jscomp$destructuring$var1 =
       (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
   x = $jscomp$destructuring$var1.next().value;
   return $jscomp$destructuring$var0;
 })([3, 4])){
console.log(x);
}
`))
}

// Java: Es6RewriteDestructuringTest.testDestructuringObjectNotInExprResult()
func (c *Es6RewriteDestructuringTest) TestDestructuringObjectNotInExprResult() {
	c.Test(
		"var x = ({a: b, c: d} = foo());",
		`var x = ($jscomp$destructuring$var0 => {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  b = $jscomp$destructuring$var1.a;
  d = $jscomp$destructuring$var1.c;
  return $jscomp$destructuring$var0;
})(foo());
`)

	c.Test(
		"var x = ({a: b, c: d} = foo());",
		`var x = ($jscomp$destructuring$var0 => {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  b = $jscomp$destructuring$var1.a;
  d = $jscomp$destructuring$var1.c;
  return $jscomp$destructuring$var0;
})(foo());
`)

	c.Test(
		"var x; var y = ({a: x} = foo());",
		`var x;
var y = ($jscomp$destructuring$var0 => {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  x = $jscomp$destructuring$var1.a;
  return $jscomp$destructuring$var0;
})(foo());
`)

	c.Test(
		"var x; var y = (() => {return {a,b} = foo();})();",
		`var x;
var y = (() => {
  return ($jscomp$destructuring$var0 => {
    var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
    a = $jscomp$destructuring$var1.a;
    b = $jscomp$destructuring$var1.b;
    return $jscomp$destructuring$var0;
  })(foo());
})();
`)
}

// Java: Es6RewriteDestructuringTest.testNestedDestructuring()
func (c *Es6RewriteDestructuringTest) TestNestedDestructuring() {
	c.Test(
		c.Srcs("var [[x]] = [[1]];"),
		c.Expected(
			`var x;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)([[1]]);
var $jscomp$destructuring$var1 =
(0, $jscomp.makeIterator)($jscomp$destructuring$var0.next().value);
x = $jscomp$destructuring$var1.next().value;
`))

	c.Test(
		c.Srcs("var [[x,y],[z]] = [[1,2],[3]];"),
		c.Expected(
			`var x;
var y;
var z;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)([[1,2],[3]]);
var $jscomp$destructuring$var1 =
(0, $jscomp.makeIterator)($jscomp$destructuring$var0.next().value);
x = $jscomp$destructuring$var1.next().value;
y = $jscomp$destructuring$var1.next().value;
var $jscomp$destructuring$var2 =
(0, $jscomp.makeIterator)($jscomp$destructuring$var0.next().value);
z = $jscomp$destructuring$var2.next().value;
`))

	c.Test(
		c.Srcs("var [[x,y],z] = [[1,2],3];"),
		c.Expected(
			`var x;
var y;
var z;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)([[1,2],3]);
var $jscomp$destructuring$var1 =
(0, $jscomp.makeIterator)($jscomp$destructuring$var0.next().value);
x = $jscomp$destructuring$var1.next().value;
y = $jscomp$destructuring$var1.next().value;
z = $jscomp$destructuring$var0.next().value;
`))
}

// Java: Es6RewriteDestructuringTest.testTryCatch()
func (c *Es6RewriteDestructuringTest) TestTryCatch() {
	c.Test(
		`var a = 1;
try {
  throw [];
} catch ([x]) {}
`,
		`var a = 1;
try {
  throw [];
} catch ($jscomp$destructuring$var0) {
   var $jscomp$destructuring$var1 = (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
   let x = $jscomp$destructuring$var1.next().value;
}
`)

	c.Test(
		`var a = 1;
try {
  throw [[]];
} catch ([[x]]) {}
`,
		`var a = 1;
try {
  throw [[]];
} catch ($jscomp$destructuring$var0) {
   var $jscomp$destructuring$var1 = (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
   var $jscomp$destructuring$var2 =
(0, $jscomp.makeIterator)($jscomp$destructuring$var1.next().value);
   let x = $jscomp$destructuring$var2.next().value;
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestDecl()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestDecl() {
	c.Test(
		"var {a: b, c: d, ...rest} = foo();",
		`var b; var d; var rest
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
b = $jscomp$destructuring$var0.a;
d = $jscomp$destructuring$var0.c;
rest = (delete $jscomp$destructuring$var1.a,
            delete $jscomp$destructuring$var1.c,
            $jscomp$destructuring$var1);
`)

	c.Test(
		"const {a: b, c: d, ...rest} = foo();",
		`/** @const */ var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
const b = $jscomp$destructuring$var0.a;
const d = $jscomp$destructuring$var0.c;
const rest = (delete $jscomp$destructuring$var1.a,
              delete $jscomp$destructuring$var1.c,
              $jscomp$destructuring$var1);
`)

	c.Test(
		"let {a: b, c: d, ...rest} = foo();",
		`var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
let b = $jscomp$destructuring$var0.a;
let d = $jscomp$destructuring$var0.c;
let rest = (delete $jscomp$destructuring$var1.a,
            delete $jscomp$destructuring$var1.c,
            $jscomp$destructuring$var1);
`)

	c.Test(
		"var pre = foo(); var {a: b, c: d, ...rest} = foo();",
		`var pre = foo();
var b; var d; var rest;
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
b = $jscomp$destructuring$var0.a;
d = $jscomp$destructuring$var0.c;
rest = (delete $jscomp$destructuring$var1.a,
            delete $jscomp$destructuring$var1.c,
            $jscomp$destructuring$var1);
`)

	c.Test(
		"var {a: b, c: d, ...rest} = foo(); var post = foo();",
		`var b; var d; var rest;
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
b = $jscomp$destructuring$var0.a;
d = $jscomp$destructuring$var0.c;
rest = (delete $jscomp$destructuring$var1.a,
            delete $jscomp$destructuring$var1.c,
            $jscomp$destructuring$var1);
var post = foo();
`)

	c.Test(
		"var pre = foo(); var {a: b, c: d, ...rest} = foo(); var post = foo();",
		`var pre = foo();
var b; var d; var rest;
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
b = $jscomp$destructuring$var0.a;
d = $jscomp$destructuring$var0.c;
rest = (delete $jscomp$destructuring$var1.a,
            delete $jscomp$destructuring$var1.c,
            $jscomp$destructuring$var1);
var post = foo();
`)

	c.Test(
		"var {a: b1, c: d1, ...rest1} = foo(); var {a: b2, c: d2, ...rest2} = foo();",
		`var b1; var d1; var rest1;
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
b1 = $jscomp$destructuring$var0.a;
d1 = $jscomp$destructuring$var0.c;
rest1 = (delete $jscomp$destructuring$var1.a,
             delete $jscomp$destructuring$var1.c,
             $jscomp$destructuring$var1);
var b2; var d2; var rest2;
var $jscomp$destructuring$var2 = foo();
var $jscomp$destructuring$var3 = Object.assign({}, $jscomp$destructuring$var2);
b2 = $jscomp$destructuring$var2.a;
d2 = $jscomp$destructuring$var2.c;
rest2 = (delete $jscomp$destructuring$var3.a,
             delete $jscomp$destructuring$var3.c,
             $jscomp$destructuring$var3);
`)

	c.Test(
		"var {...rest} = foo();",
		`var rest;
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
rest = ($jscomp$destructuring$var1);
`)

	c.Test(
		"const {...rest} = foo();",
		`/** @const */ var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
const rest = ($jscomp$destructuring$var1);
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestAssignStatement()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestAssignStatement() {
	c.Test(
		"var b,d,rest; ({a: b, c: d, ...rest} = foo());",
		`var b;
var d;
var rest;
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
b = $jscomp$destructuring$var0.a;
d = $jscomp$destructuring$var0.c;
rest = (delete $jscomp$destructuring$var1.a,
            delete $jscomp$destructuring$var1.c,
            $jscomp$destructuring$var1);
`)

	c.Test(
		"var b,d,rest,pre; pre = foo(), {a: b, c: d, ...rest} = foo();",
		`var b;
var d;
var rest;
var pre;
pre = foo(), ($jscomp$destructuring$var0 => {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 =
      Object.assign({}, $jscomp$destructuring$var1);
  b = $jscomp$destructuring$var1.a;
  d = $jscomp$destructuring$var1.c;
  rest =
      (delete $jscomp$destructuring$var2.a, delete $jscomp$destructuring$var2.c,
       $jscomp$destructuring$var2);
  return $jscomp$destructuring$var0;
})(foo());
`)

	c.Test(
		"var b,d,rest,post; ({a: b, c: d, ...rest} = foo()), post = foo();",
		`var b;
var d;
var rest;
var post;
($jscomp$destructuring$var0 => {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 =
      Object.assign({}, $jscomp$destructuring$var1);
  b = $jscomp$destructuring$var1.a;
  d = $jscomp$destructuring$var1.c;
  rest =
      (delete $jscomp$destructuring$var2.a, delete $jscomp$destructuring$var2.c,
       $jscomp$destructuring$var2);
  return $jscomp$destructuring$var0;
})(foo()),
    post = foo();
`)

	c.Test(
		"var b,d,rest,pre,post; pre = foo(), {a: b, c: d, ...rest} = foo(), post = foo();",
		`var b;
var d;
var rest;
var pre;
var post;
pre = foo(), ($jscomp$destructuring$var0 => {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 =
      Object.assign({}, $jscomp$destructuring$var1);
  b = $jscomp$destructuring$var1.a;
  d = $jscomp$destructuring$var1.c;
  rest =
      (delete $jscomp$destructuring$var2.a, delete $jscomp$destructuring$var2.c,
       $jscomp$destructuring$var2);
  return $jscomp$destructuring$var0;
})(foo()),
    post = foo();
`)

	c.Test(
		`var b1,d1,rest1,b2,d2,rest2;
({a: b1, c: d1, ...rest1} = foo(),
 {a: b2, c: d2, ...rest2} = foo());
`,
		`var b1;
var d1;
var rest1;
var b2;
var d2;
var rest2;
($jscomp$destructuring$var0 => {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 =
      Object.assign({}, $jscomp$destructuring$var1);
  b1 = $jscomp$destructuring$var1.a;
  d1 = $jscomp$destructuring$var1.c;
  rest1 =
      (delete $jscomp$destructuring$var2.a, delete $jscomp$destructuring$var2.c,
       $jscomp$destructuring$var2);
  return $jscomp$destructuring$var0;
})(foo()),
    ($jscomp$destructuring$var3 => {
      var $jscomp$destructuring$var4 = $jscomp$destructuring$var3;
      var $jscomp$destructuring$var5 =
          Object.assign({}, $jscomp$destructuring$var4);
      b2 = $jscomp$destructuring$var4.a;
      d2 = $jscomp$destructuring$var4.c;
      rest2 =
          (delete $jscomp$destructuring$var5.a,
           delete $jscomp$destructuring$var5.c, $jscomp$destructuring$var5);
      return $jscomp$destructuring$var3;
    })(foo());
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestAssignStatement_assignedToProperty()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestAssignStatement_assignedToProperty() {
	c.Test(
		"var a, obj = {}; ({a, ...obj.rest} = foo());",
		`var a;
var obj = {};
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
a = $jscomp$destructuring$var0.a;
obj.rest = (delete $jscomp$destructuring$var1.a, $jscomp$destructuring$var1);
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternInAssignExpression_inGeneratorUsingYield()
func (c *Es6RewriteDestructuringTest) TestObjectPatternInAssignExpression_inGeneratorUsingYield() {
	// Test usage of yield in the RHS of the assign. We need to make it an argument to the arrow
	// function, as we cannot reference 'yield' within the arrow function body itself.
	c.Test(
		`function *g(o) {
  return yield zzz(), ({[a]: o.a} = foo(yield b(), yield c())), yield zzz();
}
`,
		`function* g(o) {
  return yield zzz(), ($jscomp$destructuring$var0 => {
    var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
    o.a = $jscomp$destructuring$var1[a];
    return $jscomp$destructuring$var0;
  })(foo(yield b(), yield c())), yield zzz();
}
`)
	// A yield in the LHS is explicitly unsupported, for now. We could also provide it as an
	// argument, but it's more complicated than just providing the entire RHS, and we haven't seen
	// any actual instances of someone trying to put yield here.
	// Additionally, the compiler actually refuses to parse yield in a destructuring pattern default
	// value, e.g. `function *g() { const {x = yield 0} = {}; }`. The edge case of someone
	// using yield in either a computed key, or as the assignment target, seems very very unlikely.
	ex :=
		AssertThrows[jutil.IRuntimeException](
			func() {
				c.Test(
					`function *g(o) {
  return ({[yield 0]: o[yield a]} = foo());
}
`,
					"")
			})
	AssertThat(ex.GetMessage()).
		Contains(
			"Cannot transpile yet: destructuring assignment referencing await or yield in lhs")
}

// Java: Es6RewriteDestructuringTest.testDestructuringAssign_notInExprResult_cannotDecomposeToStatement()
func (c *Es6RewriteDestructuringTest) TestDestructuringAssign_notInExprResult_cannotDecomposeToStatement() {
	// Test a case where we cannot trivially extract the destructuring expression to a statement.
	// Instead we wrap it in an arrow function.
	c.Test(
		`for (let
  a = 0,
  b = ([z] = y);;
) {}
`,
		`for (let a = 0, b = ($jscomp$destructuring$var0 => {
                  var $jscomp$destructuring$var1 =
                      (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
                  z = $jscomp$destructuring$var1.next().value;
                  return $jscomp$destructuring$var0;
                })(y);
     ;) {
}
`)

	c.Test(
		`function *gen() {
  for (let
    a = 0,
    b = ([z] = yield y);;
  ) {}
}
`,
		`function* gen() {
  for (let a = 0, b = ($jscomp$destructuring$var0 => {
                    var $jscomp$destructuring$var1 =
                        (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
                    z = $jscomp$destructuring$var1.next().value;
                    return $jscomp$destructuring$var0;
                  })(yield y);
       ;) {
  }
}
`)

	c.Test(
		`async function f() {
  for (let a = 0, b = ([z] = await y());;) {}
}
`,
		`async function f() {
  for (let a = 0, b = ($jscomp$destructuring$var0 => {
                    var $jscomp$destructuring$var1 =
                        (0, $jscomp.makeIterator)($jscomp$destructuring$var0);
                    z = $jscomp$destructuring$var1.next().value;
                    return $jscomp$destructuring$var0;
                  })(await y());
       ;) {
  }
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestAssignExpr()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestAssignExpr() {
	c.Test(
		"var x,b,d,rest; x = ({a: b, c: d, ...rest} = foo());",
		`var x;
var b;
var d;
var rest;
x = ($jscomp$destructuring$var0 => {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 =
      Object.assign({}, $jscomp$destructuring$var1);
  b = $jscomp$destructuring$var1.a;
  d = $jscomp$destructuring$var1.c;
  rest =
      (delete $jscomp$destructuring$var2.a, delete $jscomp$destructuring$var2.c,
       $jscomp$destructuring$var2);
  return $jscomp$destructuring$var0;
})(foo());
`)

	c.Test(
		"var x,b,d,rest; baz({a: b, c: d, ...rest} = foo());",
		`var x;
var b;
var d;
var rest;
baz(($jscomp$destructuring$var0 => {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 =
      Object.assign({}, $jscomp$destructuring$var1);
  b = $jscomp$destructuring$var1.a;
  d = $jscomp$destructuring$var1.c;
  rest =
      (delete $jscomp$destructuring$var2.a, delete $jscomp$destructuring$var2.c,
       $jscomp$destructuring$var2);
  return $jscomp$destructuring$var0;
})(foo()));
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestForOf()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestForOf() {
	c.Test(
		"for ({a: b, c: d, ...rest} of foo()) { console.log(rest.z); }",
		`var $jscomp$destructuring$var0;
for ($jscomp$destructuring$var0 of foo()) {
    var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
    var $jscomp$destructuring$var2 = Object.assign({}, $jscomp$destructuring$var1);
    b = $jscomp$destructuring$var1.a;
    d = $jscomp$destructuring$var1.c;
    rest = (delete $jscomp$destructuring$var2.a,
            delete $jscomp$destructuring$var2.c,
            $jscomp$destructuring$var2);
    console.log(rest.z);
}
`)

	c.Test(
		"for (var {a: b, c: d, ...rest} of foo()) { console.log(rest.z); }",
		`var b;
var d;
var rest;
var $jscomp$destructuring$var0;
for ($jscomp$destructuring$var0 of foo()) {
    var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
    var $jscomp$destructuring$var2 = Object.assign({}, $jscomp$destructuring$var1);
    b = $jscomp$destructuring$var1.a;
    d = $jscomp$destructuring$var1.c;
    rest = (delete $jscomp$destructuring$var2.a,
            delete $jscomp$destructuring$var2.c,
            $jscomp$destructuring$var2);
    console.log(rest.z);
}
`)

	c.Test(
		"for (let {a: b, c: d, ...rest} of foo()) { console.log(rest.z); }",
		`for (let $jscomp$destructuring$var0 of foo()) {
    var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
    var $jscomp$destructuring$var2 = Object.assign({}, $jscomp$destructuring$var1);
    let b = $jscomp$destructuring$var1.a;
    let d = $jscomp$destructuring$var1.c;
    let rest = (delete $jscomp$destructuring$var2.a,
            delete $jscomp$destructuring$var2.c,
            $jscomp$destructuring$var2);
    {
      console.log(rest.z);
    }
}
`)

	c.Test(
		"for (const {a: b, c: d, ...rest} of foo()) { console.log(rest.z); }",
		`for (const $jscomp$destructuring$var0 of foo()) {
    /** @const */ var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
    var $jscomp$destructuring$var2 = Object.assign({}, $jscomp$destructuring$var1);
    const b = $jscomp$destructuring$var1.a;
    const d = $jscomp$destructuring$var1.c;
    const rest = (delete $jscomp$destructuring$var2.a,
            delete $jscomp$destructuring$var2.c,
            $jscomp$destructuring$var2);
    {
      console.log(rest.z);
    }
}
`)

	c.Test(
		"for (var {a: b, [baz()]: d, ...rest} of foo()) { console.log(rest.z); }",
		`var b;
var d;
var rest;
var $jscomp$destructuring$var0;
for ($jscomp$destructuring$var0 of foo()) {
    var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
    var $jscomp$destructuring$var2 = Object.assign({}, $jscomp$destructuring$var1);
    b = $jscomp$destructuring$var1.a;
    var $jscomp$destructuring$var3 = baz();
    d = $jscomp$destructuring$var1[$jscomp$destructuring$var3];
    rest = (delete $jscomp$destructuring$var2.a,
            delete $jscomp$destructuring$var2[$jscomp$destructuring$var3],
            $jscomp$destructuring$var2);
    console.log(rest.z);
}
`)

	c.Test(
		"for (var {a: b, [baz()]: d = 1, ...rest} of foo()) { console.log(rest.z); }",
		`var b;
var d;
var rest;
var $jscomp$destructuring$var0;
for ($jscomp$destructuring$var0 of foo()) {
    var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
    var $jscomp$destructuring$var2 = Object.assign({}, $jscomp$destructuring$var1);
    b = $jscomp$destructuring$var1.a;
    var $jscomp$destructuring$var3 = baz();
    var $jscomp$destructuring$var4 =
        $jscomp$destructuring$var1[$jscomp$destructuring$var3];
    d = $jscomp$destructuring$var4=== void 0 ? 1 : $jscomp$destructuring$var4;
    rest = (delete $jscomp$destructuring$var2.a,
            delete $jscomp$destructuring$var2[$jscomp$destructuring$var3],
            $jscomp$destructuring$var2);
    console.log(rest.z);
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestAndComputedPropertyName()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestAndComputedPropertyName() {
	c.Test(
		"var {a: b = 3, [bar()]: d, [baz()]: e, ...rest} = foo();",
		`var b; var d; var e; var rest;
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({},$jscomp$destructuring$var0);
b = $jscomp$destructuring$var0.a=== void 0 ? 3 : $jscomp$destructuring$var0.a;
var $jscomp$destructuring$var2 = bar();
d = $jscomp$destructuring$var0[$jscomp$destructuring$var2];
var $jscomp$destructuring$var3 = baz();
e = $jscomp$destructuring$var0[$jscomp$destructuring$var3];
rest = (delete $jscomp$destructuring$var1.a,
            delete $jscomp$destructuring$var1[$jscomp$destructuring$var2],
            delete $jscomp$destructuring$var1[$jscomp$destructuring$var3],
            $jscomp$destructuring$var1);
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestAndDefaults()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestAndDefaults() {
	c.Test(
		"var {a = 3, ...rest} = foo();",
		`var a; var rest;
var $jscomp$destructuring$var0 = foo();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
a = $jscomp$destructuring$var0.a=== void 0 ? 3 : $jscomp$destructuring$var0.a;
rest = (delete $jscomp$destructuring$var1.a,
            $jscomp$destructuring$var1);
`)

	c.Test(
		"var {[bar()]:a = 3, 'b c':b = 12, ...rest} = foo();",
		`var a; var b; var rest;
var $jscomp$destructuring$var0=foo();
var $jscomp$destructuring$var1 = Object.assign({},$jscomp$destructuring$var0);
var $jscomp$destructuring$var2 = bar();
var $jscomp$destructuring$var3 =
    $jscomp$destructuring$var0[$jscomp$destructuring$var2];
a = $jscomp$destructuring$var3=== void 0 ? 3 : $jscomp$destructuring$var3;
b = $jscomp$destructuring$var0["b c"]=== void 0
    ? 12 : $jscomp$destructuring$var0["b c"];
rest=(delete $jscomp$destructuring$var1[$jscomp$destructuring$var2],
          delete $jscomp$destructuring$var1["b c"],
          $jscomp$destructuring$var1);
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestInCatch()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestInCatch() {
	c.Test(
		"try {} catch ({first, second, ...rest}) { console.log(rest.z); }",
		`try {}
catch ($jscomp$destructuring$var0) {
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 = Object.assign({}, $jscomp$destructuring$var1);
  let first = $jscomp$destructuring$var1.first;
  let second = $jscomp$destructuring$var1.second;
  let rest = (delete $jscomp$destructuring$var2.first,
              delete $jscomp$destructuring$var2.second,
              $jscomp$destructuring$var2);
  console.log(rest.z);
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestAssignReturn()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestAssignReturn() {
	c.Test(
		"function f() { return {x:a, ...rest} = foo(); }",
		`function f() {
  return ($jscomp$destructuring$var0 => {
    var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
    var $jscomp$destructuring$var2 =
        Object.assign({}, $jscomp$destructuring$var1);
    a = $jscomp$destructuring$var1.x;
    rest = (delete $jscomp$destructuring$var2.x, $jscomp$destructuring$var2);
    return $jscomp$destructuring$var0;
  })(foo());
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestParamList()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestParamList() {
	c.Test(
		"function f({x = a(), ...rest}, y=b()) { console.log(y); }",
		`function f($jscomp$destructuring$var0,y) {
  var x; var rest;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 = Object.assign({},$jscomp$destructuring$var1);
  x = $jscomp$destructuring$var1.x === void 0
      ? a() : $jscomp$destructuring$var1.x;
  rest= (delete $jscomp$destructuring$var2.x,
             $jscomp$destructuring$var2);
  y = y=== void 0 ? b() : y;
  console.log(y)
}
`)

	c.Test(
		"function f({x = a(), ...rest}={}, y=b()) { console.log(y); }",
		`function f($jscomp$destructuring$var0,y) {
  var x; var rest;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0=== void 0
      ? {} : $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 = Object.assign({},$jscomp$destructuring$var1);
  x = $jscomp$destructuring$var1.x=== void 0
      ? a() : $jscomp$destructuring$var1.x;
  rest= (delete $jscomp$destructuring$var2.x,
             $jscomp$destructuring$var2);
  y = y=== void 0 ? b() : y;
  console.log(y)
}
`)
}

// Java: Es6RewriteDestructuringTest.testObjectPatternWithRestArrowParamList()
func (c *Es6RewriteDestructuringTest) TestObjectPatternWithRestArrowParamList() {
	c.Test(
		"var f = ({x = a(), ...rest}, y=b()) => { console.log(y); };",
		`var f = ($jscomp$destructuring$var0,y) => {
  var x; var rest;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 = Object.assign({},$jscomp$destructuring$var1);
  x = $jscomp$destructuring$var1.x=== void 0
      ? a() : $jscomp$destructuring$var1.x;
  rest = (delete $jscomp$destructuring$var2.x,
              $jscomp$destructuring$var2);
  y = y=== void 0 ? b() : y;
  console.log(y)
}
`)

	c.Test(
		"var f = ({x = a(), ...rest}={}, y=b()) => { console.log(y); };",
		`var f = ($jscomp$destructuring$var0,y) => {
  var x; var rest;
  var $jscomp$destructuring$var1 = $jscomp$destructuring$var0=== void 0
      ? {} : $jscomp$destructuring$var0;
  var $jscomp$destructuring$var2 = Object.assign({},$jscomp$destructuring$var1);
  x = $jscomp$destructuring$var1.x=== void 0
      ? a() : $jscomp$destructuring$var1.x;
  rest= (delete $jscomp$destructuring$var2.x,
             $jscomp$destructuring$var2);
  y = y=== void 0 ? b() : y;
  console.log(y)
}
`)
}

// Java: Es6RewriteDestructuringTest.testAllRewriteMode()
func (c *Es6RewriteDestructuringTest) TestAllRewriteMode() {
	c.destructuringRewriteMode = Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS

	c.Test(
		"var {a} = foo();",
		`var a; var $jscomp$destructuring$var0 = foo();
a = $jscomp$destructuring$var0.a;
`)

	c.Test(
		"var {a} = foo(); var {...b} = bar();",
		`var a;
var $jscomp$destructuring$var0 = foo();
a = $jscomp$destructuring$var0.a;
var b;
var $jscomp$destructuring$var1 = bar();
var $jscomp$destructuring$var2 = Object.assign({}, $jscomp$destructuring$var1);
b = $jscomp$destructuring$var2;
`)

	c.Test(
		"var {[foo0()]: {[foo1()]: a, ...r}, [foo2()]: { [foo3()]: b}, [foo4()]: c } = bar();",
		`var a; var r; var b; var c;
var $jscomp$destructuring$var0 = bar();
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0[foo0()];
var $jscomp$destructuring$var2 = Object.assign({},$jscomp$destructuring$var1);
var $jscomp$destructuring$var3 = foo1();
a = $jscomp$destructuring$var1[$jscomp$destructuring$var3];
r = (delete $jscomp$destructuring$var2[$jscomp$destructuring$var3],
         $jscomp$destructuring$var2);
var $jscomp$destructuring$var4 = $jscomp$destructuring$var0[foo2()];
b = $jscomp$destructuring$var4[foo3()];
c = $jscomp$destructuring$var0[foo4()]
`)

	c.Test(
		"var [a] = foo(); var [{b}] = foo(); var [{...c}] = foo();",
		`var a;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(foo());
a = $jscomp$destructuring$var0.next().value;
// var [{b}] = foo();
var b;
var $jscomp$destructuring$var1 = (0, $jscomp.makeIterator)(foo());
var $jscomp$destructuring$var2 = $jscomp$destructuring$var1.next().value;
b = $jscomp$destructuring$var2.b;
// var [{...c}] = foo();
var c;
var $jscomp$destructuring$var3 = (0, $jscomp.makeIterator)(foo());
var $jscomp$destructuring$var4 = $jscomp$destructuring$var3.next().value;
var $jscomp$destructuring$var5 = Object.assign({},$jscomp$destructuring$var4);
c = $jscomp$destructuring$var5
`)
}

// Java: Es6RewriteDestructuringTest.testOnlyRestRewriteMode()
func (c *Es6RewriteDestructuringTest) TestOnlyRestRewriteMode() {
	c.destructuringRewriteMode = Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST

	c.Test("var {a} = foo();", "var a; ({a} = foo());")

	c.Test(
		"var {a} = foo(); var {...b} = bar();",
		`var a; ({a} = foo());
var b;
var $jscomp$destructuring$var0 = bar();
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
b = $jscomp$destructuring$var1;
`)

	// test that object patterns are rewritten if they have a rest property nested within
	c.Test(
		"var {[foo0()]: {[foo1()]: a, ...r}, [foo2()]: { [foo3()]: b}, [foo4()]: c } = bar();",
		`var a; var r; var b; var c;
var $jscomp$destructuring$var0 = bar();
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0[foo0()];
var $jscomp$destructuring$var2 = Object.assign({},$jscomp$destructuring$var1);
var $jscomp$destructuring$var3 = foo1();
a = $jscomp$destructuring$var1[$jscomp$destructuring$var3];
r = (delete $jscomp$destructuring$var2[$jscomp$destructuring$var3],
         $jscomp$destructuring$var2);
var $jscomp$destructuring$var4 = $jscomp$destructuring$var0[foo2()];
b = $jscomp$destructuring$var4[foo3()];
c = $jscomp$destructuring$var0[foo4()]
`)

	c.Test(
		"var [a] = foo(); var [{b}] = foo(); var [{...c}] = foo();",
		`var a;
[a] = foo();
// var [{b}] = foo();
var b;
[{b}] = foo();
// var [{...c}] = foo();
var c;
var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(foo());
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0.next().value;
var $jscomp$destructuring$var2 = Object.assign({},$jscomp$destructuring$var1);
c = $jscomp$destructuring$var2
`)
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuring_getsCorrectTypes()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuring_getsCorrectTypes() {

	c.Test(
		c.Srcs(
			`function takesIterable(/** !Iterable<number> */ iterableVar) {
  const [a] = iterableVar;
}
`),
		c.Expected(
			`function takesIterable(iterableVar) {
  var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(iterableVar);
  const a = $jscomp$destructuring$var0.next().value;
}
`))

	lastCompiler := c.GetLastCompiler()

	// `$jscomp$destructuring$var0` is an Iterator<number>
	jscompDestructuringVar0 :=
		c.GetNodeMatchingQName(lastCompiler.GetJsRoot(), "$jscomp$destructuring$var0")
	AssertNode(jscompDestructuringVar0).
		HasColorThat().
		IsEqualTo(lastCompiler.GetColorRegistry().Get(colors.StandardColors_ITERATOR_ID))

	// `a` is a number
	aName := c.GetNodeMatchingQName(lastCompiler.GetJsRoot(), "a")
	AssertNode(aName).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)

	// `$jscomp$destructuring$var0.next().value` is a number
	destructuringVarNextDotValue := aName.GetOnlyChild()
	AssertNode(destructuringVarNextDotValue).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)

	// `$jscomp$destructuring$var0.next()` is an IIterableResult<number>
	destructuringVarNextCall := destructuringVarNextDotValue.GetFirstChild()
	AssertNode(destructuringVarNextCall).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
}

// Java: Es6RewriteDestructuringTest.testArrayDestructuringRest_getsCorrectTypes()
func (c *Es6RewriteDestructuringTest) TestArrayDestructuringRest_getsCorrectTypes() {
	c.Test(
		c.Srcs(
			`function takesIterable(/** !Iterable<number> */ iterableVar) {
  const [a, ...rest] = iterableVar;
}
`),
		c.Expected(
			`function takesIterable(iterableVar) {
  var $jscomp$destructuring$var0 = (0, $jscomp.makeIterator)(iterableVar);
  const a = $jscomp$destructuring$var0.next().value;
  const rest = (0, $jscomp.arrayFromIterator)($jscomp$destructuring$var0);
}
`))

	lastCompiler := c.GetLastCompiler()

	jscompArrayFromIterator :=
		c.GetNodeMatchingQName(lastCompiler.GetJsRoot(), "$jscomp.arrayFromIterator")
	AssertNode(jscompArrayFromIterator).HasColorThat().IsEqualTo(colors.StandardColors_TOP_OBJECT)

	jscompArrayFromIteratorCall := jscompArrayFromIterator.GetParent()
	AssertNode(jscompArrayFromIteratorCall).
		HasColorThat().
		IsEqualTo(lastCompiler.GetColorRegistry().Get(colors.StandardColors_ARRAY_ID))

	// `rest` is Array<number>
	restName := c.GetNodeMatchingQName(lastCompiler.GetJsRoot(), "rest")
	AssertNode(restName).
		HasColorThat().
		IsEqualTo(lastCompiler.GetColorRegistry().Get(colors.StandardColors_ARRAY_ID))
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuring_getsCorrectTypes()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuring_getsCorrectTypes() {
	c.Test(
		`const obj = {a: 3, b: 'string', c: null};
const {a} = obj;
`,
		`const obj = {a: 3, b: 'string', c: null};
/** @const */ var $jscomp$destructuring$var0=obj;
const a = $jscomp$destructuring$var0.a;
`)

	jsRoot := c.GetLastCompiler().GetJsRoot()

	aName := c.GetNodeMatchingQName(jsRoot, "a")
	AssertNode(aName).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)

	objType := c.GetNodeMatchingQName(jsRoot, "obj").GetColor()
	// `$jscomp$destructuring$var0` has the same type as `obj`
	jscompDestructuringVar0Name := c.GetNodeMatchingQName(jsRoot, "$jscomp$destructuring$var0")
	AssertNode(jscompDestructuringVar0Name).HasColorThat().IsEqualTo(objType)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuringDefaultValue_getsCorrectTypes()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuringDefaultValue_getsCorrectTypes() {
	c.Test(
		`const obj = {a: 3, b: 'string', c: null};
const {a = 4} = obj;
`,
		`const obj = {a: 3, b: 'string', c: null};
/** @const */ var $jscomp$destructuring$var0=obj;
const a = $jscomp$destructuring$var0.a === void 0
    ? 4: $jscomp$destructuring$var0.a;
`)

	jsRoot := c.GetLastCompiler().GetJsRoot()

	aName := c.GetNodeMatchingQName(jsRoot, "a")
	AssertNode(aName).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)

	objType := c.GetNodeMatchingQName(jsRoot, "obj").GetColor()

	// `$jscomp$destructuring$var0` has the same type as `obj`
	colorSet := jutil.NewObjectLinkedHashSet[*colors.Color]()
	for _, n := range c.GetAllNodesMatchingQName(jsRoot, "$jscomp$destructuring$var0") {
		colorSet.Add(n.GetColor())
	}
	AssertThat(colorSet.ToSlice()).
		ContainsExactly(objType)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuringComputedPropWithDefault_getsCorrectTypes()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuringComputedPropWithDefault_getsCorrectTypes() {
	c.Test(
		`const /** !Object<string, number> */ obj = {['a']: 3};
const {['a']: a = 4} = obj;
`,
		`const obj = {['a']: 3};
/** @const */ var $jscomp$destructuring$var0=obj;
var $jscomp$destructuring$var1 = $jscomp$destructuring$var0['a'];
const a = $jscomp$destructuring$var1 === void 0
    ? 4: $jscomp$destructuring$var1;
`)

	jsRoot := c.GetLastCompiler().GetJsRoot()

	aName := c.GetNodeMatchingQName(jsRoot, "a")
	AssertNode(aName).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)

	// `$jscomp$destructuring$var0` has the same type as `obj`
	jscompDestructuringVar0Name := c.GetNodeMatchingQName(jsRoot, "$jscomp$destructuring$var0")

	objType := jscompDestructuringVar0Name.GetOnlyChild().GetColor()
	AssertThat(objType).IsEqualTo(colors.StandardColors_TOP_OBJECT)
	AssertNode(jscompDestructuringVar0Name).HasColorThat().IsEqualTo(objType)

	// `$jscomp$destructuring$var1` is typed as `number` (this is probably less important!)
	jscompDestructuringVar1Name := c.GetNodeMatchingQName(jsRoot, "$jscomp$destructuring$var0")
	AssertNode(jscompDestructuringVar1Name).HasColorThat().IsEqualTo(objType)
}

// Java: Es6RewriteDestructuringTest.testObjectDestructuringRest_typesRestAsObject()
func (c *Es6RewriteDestructuringTest) TestObjectDestructuringRest_typesRestAsObject() {
	c.Test(
		`const obj = {a: 3, b: 'string', c: null};
const {...rest} = obj;
`,
		`const obj = {a: 3, b: 'string', c: null};
/** @const */ var $jscomp$destructuring$var0=obj;
var $jscomp$destructuring$var1 = Object.assign({}, $jscomp$destructuring$var0);
const rest = $jscomp$destructuring$var1;
`)

	jsRoot := c.GetLastCompiler().GetJsRoot()

	objType := c.GetNodeMatchingQName(jsRoot, "obj").GetColor()

	// `$jscomp$destructuring$var0` has the same type as `obj`
	jscompDestructuringVar0Name := c.GetNodeMatchingQName(jsRoot, "$jscomp$destructuring$var0")
	AssertNode(jscompDestructuringVar0Name).HasColorThat().IsEqualTo(objType)

	// `$jscomp$destructuring$var1` has the same type as `obj`
	jscompDestructuringVar1Name := c.GetNodeMatchingQName(jsRoot, "$jscomp$destructuring$var1")
	AssertNode(jscompDestructuringVar1Name).HasColorThat().IsEqualTo(objType)

	// TODO(b/128355893) Do better inferrence. For now we just consider `rest` an `Object` rather
	// than trying to figure out what properties it gets.
	restName := c.GetNodeMatchingQName(jsRoot, "rest")
	AssertNode(restName).HasColorThat().IsEqualTo(colors.StandardColors_TOP_OBJECT)
}

// Returns a list of all nodes in the given AST that matches the given qualified name
// Java: Es6RewriteDestructuringTest.getAllNodesMatchingQName(Node, String)
func (c *Es6RewriteDestructuringTest) GetAllNodesMatchingQName(root *rhino.Node, qname string) []*rhino.Node {
	builder := []*rhino.Node{}
	c.AddAllNodesMatchingQNameHelper(root, qname, &builder)
	return builder
}

// Java: Es6RewriteDestructuringTest.addAllNodesMatchingQNameHelper(Node, String, ImmutableList.Builder<Node>)
func (c *Es6RewriteDestructuringTest) AddAllNodesMatchingQNameHelper(root *rhino.Node, qname string, nodesSoFar *[]*rhino.Node) {
	if root.MatchesQualifiedName(qname) {
		*nodesSoFar = append(*nodesSoFar, root)
	}
	for child := root.GetFirstChild(); child != nil; child = child.GetNext() {
		c.AddAllNodesMatchingQNameHelper(child, qname, nodesSoFar)
	}
}

// Returns the first node (preorder) in the given AST that matches the given qualified name
// Java: Es6RewriteDestructuringTest.getNodeMatchingQName(Node, String)
func (c *Es6RewriteDestructuringTest) GetNodeMatchingQName(root *rhino.Node, qname string) *rhino.Node {
	if root.MatchesQualifiedName(qname) {
		return root
	}
	for child := root.GetFirstChild(); child != nil; child = child.GetNext() {
		result := c.GetNodeMatchingQName(child, qname)
		if result != nil {
			return result
		}
	}
	return nil
}

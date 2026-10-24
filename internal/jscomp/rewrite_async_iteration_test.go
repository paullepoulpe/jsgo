// Port of com/google/javascript/jscomp/RewriteAsyncIterationTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/RewriteAsyncIterationTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteAsyncIterationTest has the test cases for RewriteAsyncIteration.
type RewriteAsyncIterationTest struct {
	CompilerTestCaseBase
}

// Java: RewriteAsyncIterationTest()
func NewRewriteAsyncIterationTest() *RewriteAsyncIterationTest {
	c := &RewriteAsyncIterationTest{}
	c.InitCompilerTestCaseBaseWithExterns(c,
		NewTestExternsBuilder().
			AddAsyncIterable().
			AddArray().
			AddMath().
			AddArguments().
			AddObject().
			AddJSCompLibraries().
			Build())
	return c
}

func TestRewriteAsyncIterationTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteAsyncIterationTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then enableTypeCheckBeforePass().
func (c *RewriteAsyncIterationTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.EnableTypeCheckBeforePass()
}

// Java: RewriteAsyncIterationTest.enableTypeCheckBeforePass()
func (c *RewriteAsyncIterationTest) EnableTypeCheckBeforePass() {
	c.EnableNormalize()
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.AllowExternsChanges()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: RewriteAsyncIterationTest.getOptions()
func (c *RewriteAsyncIterationTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2017)
	return options
}

// Java: RewriteAsyncIterationTest.getProcessor(Compiler)
func (c *RewriteAsyncIterationTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"injectTranspilationRuntimeLibraries", func(abstractCompiler AbstractCompiler) CompilerPass {
				return NewInjectTranspilationRuntimeLibraries(abstractCompiler)
			}))
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory("rewriteAsyncIteration", func(abstractCompiler AbstractCompiler) CompilerPass {
			return RewriteAsyncIteration_Create(abstractCompiler)
		}))
	return optimizer
}

// Java: RewriteAsyncIterationTest.getGlobalColor(ColorId)
func (c *RewriteAsyncIterationTest) GetGlobalColor(colorId colors.ColorId) *colors.Color {
	return c.GetLastCompiler().GetColorRegistry().Get(colorId)
}

// Java: RewriteAsyncIterationTest.getJSCompAsyncGeneratorWrapperClassColor()
func (c *RewriteAsyncIterationTest) GetJSCompAsyncGeneratorWrapperClassColor() *colors.Color {
	return CodeSubTree_FindFirstNode(
		c.GetLastCompiler().GetExternsRoot(),
		func(n rhino.Node) bool { return n.MatchesQualifiedName("$jscomp.AsyncGeneratorWrapper") }).
		GetNext().
		GetColor()
}

// Java: RewriteAsyncIterationTest.getJSCompAsyncGeneratorWrapperInstanceColor()
func (c *RewriteAsyncIterationTest) GetJSCompAsyncGeneratorWrapperInstanceColor() *colors.Color {
	return colors.Color_CreateUnion(c.GetJSCompAsyncGeneratorWrapperClassColor().GetInstanceColors())
}

// Java: RewriteAsyncIterationTest.testForAwaitWithThrow()
func (c *RewriteAsyncIterationTest) TestForAwaitWithThrow() {
	c.Test(
		`async function test() {
    for await (const i of source()) {
      if (i === 2) {
        throw new Error('');
      }
    }
}
`,
		`async function test() {
  var $jscomp$forAwait$errResult0;
  var $jscomp$forAwait$tempResult0;
  var $jscomp$forAwait$retFn0;
  try {
    var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(source());
    for (;;) {
      $jscomp$forAwait$tempResult0 = void 0;
      $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
      if ($jscomp$forAwait$tempResult0.done) {
        break;
      }
      const i = $jscomp$forAwait$tempResult0.value;
      {
        if (i === 2) {
          throw new Error("");
        }
      }
    }
  } catch ($jscomp$forAwait$catchErrParam0) {
      $jscomp$forAwait$errResult0 = {error:$jscomp$forAwait$catchErrParam0};
  } finally {
    try {
      if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done && ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
        await $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0);
      }
    }
    finally {
      if ($jscomp$forAwait$errResult0) {
        throw $jscomp$forAwait$errResult0.error;
      }
    }
  }
}
`)
}

// Java: RewriteAsyncIterationTest.testBug173319540()
func (c *RewriteAsyncIterationTest) TestBug173319540() {
	c.Test(
		c.Srcs(
			`let key;
let value;
window.onload = async function() {
  for await ([key,value] of window["unknownAsyncIterable"]) {
    alert(key,value);
  }
}
`),
		c.Expected(
			`let key;
let value;
window.onload = async function() {
  var $jscomp$forAwait$errResult0;
  var $jscomp$forAwait$tempResult0;
  var $jscomp$forAwait$retFn0;
  try {
var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(window["unknownAsyncIterable"]);
    for (;;) {
      $jscomp$forAwait$tempResult0 = void 0;
      $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
      if ($jscomp$forAwait$tempResult0.done) {
        break;
      }
      [key, value] = $jscomp$forAwait$tempResult0.value;
      { alert(key, value);}
    }
  } catch ($jscomp$forAwait$catchErrParam0) {
    $jscomp$forAwait$errResult0 = {error:$jscomp$forAwait$catchErrParam0};
  } finally {
    try {
      if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done &&
        ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
        await $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0);
      }
    } finally {
      if ($jscomp$forAwait$errResult0) {
        throw $jscomp$forAwait$errResult0.error;
      }
    }
  }
};

`))

	c.Test(
		c.Srcs(
			`window.onload = async function() {
  for await (const [key,value] of window["unknownAsyncIterable"]) {
    alert(key,value);
  }
}
`),
		c.Expected(
			`window.onload = async function() {
  var $jscomp$forAwait$errResult0;
  var $jscomp$forAwait$tempResult0;
  var $jscomp$forAwait$retFn0;
  try {
var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(window["unknownAsyncIterable"]);
    for (;;) {
      $jscomp$forAwait$tempResult0 = void 0;
      $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
      if ($jscomp$forAwait$tempResult0.done) {
        break;
      }
      const [key, value$jscomp$3] = $jscomp$forAwait$tempResult0.value;
      {
        alert(key, value$jscomp$3);
      }
    }
  } catch ($jscomp$forAwait$catchErrParam0) {
    $jscomp$forAwait$errResult0 = { error:$jscomp$forAwait$catchErrParam0 };
  } finally {
    try {
      if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done &&
        ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
        await $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0);
      }
    } finally {
      if ($jscomp$forAwait$errResult0) {
        throw $jscomp$forAwait$errResult0.error;
      }
    }
  }
};

`))
}

// Java: RewriteAsyncIterationTest.testAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestAsyncGenerator() {
	c.Test(
		"async function* baz() { foo() }",
		`function baz() {
  return new $jscomp.AsyncGeneratorWrapper((function*() {
    foo();
  })());
}
`)

	bazSubTree := CodeSubTree_FindFunctionDefinition(c.GetLastCompiler(), "baz")
	wrapper := bazSubTree.FindMatchingQNameReferences("$jscomp.AsyncGeneratorWrapper")[0]
	newExpr := wrapper.GetParent()
	innerGeneratorCall := newExpr.GetSecondChild()

	AssertNode(wrapper).HasColorThat().IsEqualTo(c.GetJSCompAsyncGeneratorWrapperClassColor())
	AssertNode(newExpr).HasColorThat().IsEqualTo(c.GetJSCompAsyncGeneratorWrapperInstanceColor())
	AssertNode(innerGeneratorCall).
		HasColorThat().
		IsEqualTo(c.GetGlobalColor(colors.StandardColors_GENERATOR_ID))
}

// Java: RewriteAsyncIterationTest.testAwaitInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestAwaitInAsyncGenerator() {
	c.Test(
		`let /** function(): !Promise<number> */ foo;
/** @return {!AsyncGenerator<undefined>} */
async function* baz() { await foo() }
`,
		`let foo;
function baz() {
  return new $jscomp.AsyncGeneratorWrapper((function*() {
    yield new $jscomp.AsyncGeneratorWrapper$ActionRecord(
      $jscomp.AsyncGeneratorWrapper$ActionEnum.AWAIT_VALUE, foo());
  })());
}
`)

	bazSubTree := CodeSubTree_FindFunctionDefinition(c.GetLastCompiler(), "baz")
	wrapper := bazSubTree.FindMatchingQNameReferences("$jscomp.AsyncGeneratorWrapper")[0]
	newExpr := wrapper.GetParent()
	innerGeneratorCall := newExpr.GetSecondChild()

	AssertNode(wrapper).HasColorThat().IsEqualTo(c.GetJSCompAsyncGeneratorWrapperClassColor())
	AssertNode(newExpr).HasColorThat().IsEqualTo(c.GetJSCompAsyncGeneratorWrapperInstanceColor())
	AssertNode(innerGeneratorCall).
		HasColorThat().
		IsEqualTo(c.GetGlobalColor(colors.StandardColors_GENERATOR_ID))

	c.Test(
		`let /** function(): !Promise<number> */ foo;
async function* baz() { bar = await foo() }
`,
		`let foo;
function baz() {
  return new $jscomp.AsyncGeneratorWrapper((function*() {
    bar = yield new $jscomp.AsyncGeneratorWrapper$ActionRecord(
      $jscomp.AsyncGeneratorWrapper$ActionEnum.AWAIT_VALUE, foo());
  })());
}
`)

	bar :=
		CodeSubTree_FindFunctionDefinition(c.GetLastCompiler(), "baz").FindMatchingQNameReferences("bar")[0]

	AssertNode(bar).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

// Java: RewriteAsyncIterationTest.testYieldInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestYieldInAsyncGenerator() {
	c.Test(
		`/** @return {!AsyncGenerator<number>} */
async function* baz() { yield 2+2 }
`,
		`function baz() {
  return new $jscomp.AsyncGeneratorWrapper((function*() {
    yield new $jscomp.AsyncGeneratorWrapper$ActionRecord(
      $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE, 2+2);
  })());
}
`)

	bazSubTree := CodeSubTree_FindFunctionDefinition(c.GetLastCompiler(), "baz")
	wrapper := bazSubTree.FindMatchingQNameReferences("$jscomp.AsyncGeneratorWrapper")[0]
	newExpr := wrapper.GetParent()
	innerGeneratorCall := newExpr.GetSecondChild()

	AssertNode(wrapper).HasColorThat().IsEqualTo(c.GetJSCompAsyncGeneratorWrapperClassColor())
	AssertNode(newExpr).HasColorThat().IsEqualTo(c.GetJSCompAsyncGeneratorWrapperInstanceColor())
	AssertNode(innerGeneratorCall).
		HasColorThat().
		IsEqualTo(c.GetGlobalColor(colors.StandardColors_GENERATOR_ID))

	c.Test(
		`/** @return {!AsyncGenerator<number>} */
async function* baz() { bar = yield 2+2 }
`,
		`function baz() {
  return new $jscomp.AsyncGeneratorWrapper((function*() {
    bar = yield new $jscomp.AsyncGeneratorWrapper$ActionRecord(
      $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE, 2+2);
  })());
}
`)

	bar :=
		CodeSubTree_FindFunctionDefinition(c.GetLastCompiler(), "baz").FindMatchingQNameReferences("bar")[0]

	// The generator yields numbers but yield expressions should always be "?" as next accepts "?"
	AssertNode(bar).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
}

// Java: RewriteAsyncIterationTest.testYieldAllInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestYieldAllInAsyncGenerator() {
	c.Test(
		"async function* baz() { yield* foo() }",
		`function baz() {
  return new $jscomp.AsyncGeneratorWrapper((function*() {
    yield new $jscomp.AsyncGeneratorWrapper$ActionRecord(
      $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_STAR, foo());
  })());
}
`)

	c.Test(
		"async function* baz() { bar = yield* foo() }",
		`function baz() {
  return new $jscomp.AsyncGeneratorWrapper((function*() {
    bar = yield new $jscomp.AsyncGeneratorWrapper$ActionRecord(
      $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_STAR, foo());
  })());
}
`)
}

// Java: RewriteAsyncIterationTest.testComplexAsyncGeneratorStatements()
func (c *RewriteAsyncIterationTest) TestComplexAsyncGeneratorStatements() {
	c.Test(
		"async function* baz() { yield* (await foo()); }",
		`function baz() {
  return new $jscomp.AsyncGeneratorWrapper((function*() {
    yield new $jscomp.AsyncGeneratorWrapper$ActionRecord(
      $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_STAR,
      yield new $jscomp.AsyncGeneratorWrapper$ActionRecord(
          $jscomp.AsyncGeneratorWrapper$ActionEnum.AWAIT_VALUE,
          foo()));
  })());
}
`)
}

// Java: RewriteAsyncIterationTest.testThisInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestThisInAsyncGenerator() {
	c.Test(
		"async function* baz() { yield this; }",
		`function baz() {
  const $jscomp$asyncIter$this$m1146332801$0 = this;
  return new $jscomp.AsyncGeneratorWrapper((function*() {
    yield new $jscomp.AsyncGeneratorWrapper$ActionRecord(
      $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE, $jscomp$asyncIter$this$m1146332801$0);
  })());
}
`)
}

// Java: RewriteAsyncIterationTest.testThisInAsyncGeneratorNestedInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestThisInAsyncGeneratorNestedInAsyncGenerator() {
	c.Test(
		`async function* baz(outerT = this) {
  return async function*(innerT = this) {
    yield innerT || this;
  }
}
`,
		"function baz(outerT = this) {\n  return new $jscomp.AsyncGeneratorWrapper(\n      (function*() {\n        return new $jscomp.AsyncGeneratorWrapper$ActionRecord(\n            $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE,\n// `this` in parameter list shouldn't be aliased\n            function(innerT = this) {\n              const $jscomp$asyncIter$this$m1146332801$1 = this;\n              return new $jscomp.AsyncGeneratorWrapper(\n                  (function*() {\n                    yield new $jscomp.AsyncGeneratorWrapper$ActionRecord(\n                        $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE,\n// `this` in body should be aliased\n                        innerT || $jscomp$asyncIter$this$m1146332801$1);\n                  })());\n            });\n      })());\n}\n")
}

// Java: RewriteAsyncIterationTest.testThisInArrowNestedInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestThisInArrowNestedInAsyncGenerator() {
	c.Test(
		"async function* baz() {\n// both instances of `this` musts be changed to aliases\n  return (t = this) => t || this;\n}\n",
		`function baz() {
  const $jscomp$asyncIter$this$m1146332801$0 = this;
  return new $jscomp.AsyncGeneratorWrapper((function*() {
    return new $jscomp.AsyncGeneratorWrapper$ActionRecord(
        $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE,
        (t = $jscomp$asyncIter$this$m1146332801$0) => {
         return t || $jscomp$asyncIter$this$m1146332801$0});
  })());
}
`)
}

// Java: RewriteAsyncIterationTest.testThisInFunctionNestedInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestThisInFunctionNestedInAsyncGenerator() {
	c.Test(
		"async function* baz() {  return function() { return this; }; }",
		`function baz() {
  return new $jscomp.AsyncGeneratorWrapper(
      (function*() {
        return new $jscomp.AsyncGeneratorWrapper$ActionRecord(
            $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE,
            function() { return this; });
      })());
}
`)

	c.Test(
		"async function* baz() {  return () => this; }",
		`function baz() {
  const $jscomp$asyncIter$this$m1146332801$0 = this;
  return new $jscomp.AsyncGeneratorWrapper(
      (function*() {
        return new $jscomp.AsyncGeneratorWrapper$ActionRecord(
            $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE,
            () => {
               return $jscomp$asyncIter$this$m1146332801$0;
             });
      })());
}
`)
}

// Java: RewriteAsyncIterationTest.testInnerSuperReferenceInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestInnerSuperReferenceInAsyncGenerator() {
	c.Test(
		`class A {
  m() {
    return this;
  }
}
class X extends A {
  async *m() {
    const tmp = super.m;
    return tmp.call(null);
  }
}
`,
		`class A {
  m() {
    return this;
  }
}
class X extends A {
  m() {
    const $jscomp$asyncIter$super$get$m = () => {
            return super.m;
    };
    return new $jscomp.AsyncGeneratorWrapper(
        function* () {
          const tmp = $jscomp$asyncIter$super$get$m();
          return new $jscomp.AsyncGeneratorWrapper$ActionRecord(
              $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE,
              tmp.call(null));
        }());
  }
}
`)
}

// Java: RewriteAsyncIterationTest.testInnerSuperCallInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestInnerSuperCallInAsyncGenerator() {
	c.Test(
		`class A {
  m() {
    return this;
  }
}
class X extends A {
  async *m() {
    return super.m();
  }
}
`,
		`class A {
  m() {
    return this;
  }
}
class X extends A {
  m() {
    const $jscomp$asyncIter$this$m1146332801$1 = this;
    const $jscomp$asyncIter$super$get$m = () => {
         return super.m;
    };
    return new $jscomp.AsyncGeneratorWrapper(
        function* () {
          return new $jscomp.AsyncGeneratorWrapper$ActionRecord(
              $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE,
              $jscomp$asyncIter$super$get$m().call($jscomp$asyncIter$this$m1146332801$1));
        }());
  }
}
`)
}

// Java: RewriteAsyncIterationTest.testInnerSuperTaggedTemplateLiteralInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestInnerSuperTaggedTemplateLiteralInAsyncGenerator() {
	c.Test(
		"class A {\n  m(strings, ...values) {\n    return this;\n  }\n}\nclass X extends A {\n  async *m() {\n    return super.m`<p>hello</p>`;\n  }\n}\n",
		"class A {\n  m(strings, ...values) {\n    return this;\n  }\n}\nclass X extends A {\n  m() {\n    const $jscomp$asyncIter$this$m1146332801$1 = this;\n    const $jscomp$asyncIter$super$get$m = () => {\n         return super.m;\n    };\n    return new $jscomp.AsyncGeneratorWrapper(\n        function* () {\n          return new $jscomp.AsyncGeneratorWrapper$ActionRecord(\n              $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE,\n              $jscomp$asyncIter$super$get$m().bind($jscomp$asyncIter$this$m1146332801$1)`<p>hello</p>`);\n        }());\n  }\n}\n")
}

// Java: RewriteAsyncIterationTest.testCannotConvertSuperGetElemInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestCannotConvertSuperGetElemInAsyncGenerator() {
	// The rewriting gets partially done before we notice and report that we cannot convert
	// the code. The partially done code is invalid, so we must disable AST validation to see the
	// error message. (AST validation is not enabled in normal execution, just developer mode.)
	c.DisableAstValidation()
	c.TestError(
		`class A {
  m() {
    return this;
  }
}
class X extends A {
  async *m() {
    const tmp = super['m'];
    return tmp.call(null);
  }
}
`,
		RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN,
		RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN.Format(
			"super only allowed with getprop (like super.foo(), not super['foo']())"))
}

// Java: RewriteAsyncIterationTest.testInnerArrowFunctionUsingArguments()
func (c *RewriteAsyncIterationTest) TestInnerArrowFunctionUsingArguments() {
	c.Test(
		`class X {
  async *m() {
    return new Promise((resolve, reject) => {
      resolve(arguments);
    });
  }
}
`,
		`class X {
  m() {
    const $jscomp$asyncIter$arguments = arguments;
    return new $jscomp.AsyncGeneratorWrapper(
        function* () {
          return new $jscomp.AsyncGeneratorWrapper$ActionRecord(
              $jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE,
              new Promise((resolve, reject) => {
                resolve($jscomp$asyncIter$arguments);
              }));
        }());
  }
}
`)
}

// Java: RewriteAsyncIterationTest.testForAwaitOfDeclarations()
func (c *RewriteAsyncIterationTest) TestForAwaitOfDeclarations() {
	c.Test(
		`/** @type {number|undefined} */
let a;
async function abc() { for await (a of foo()) { bar(); } }
/** @return {!AsyncGenerator<number>} */
function foo() {}
`,
		`let a;
async function abc() {
  var $jscomp$forAwait$errResult0;
  var $jscomp$forAwait$tempResult0;
  var $jscomp$forAwait$retFn0;
  try {
    var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(foo());
    for (;;) {
      $jscomp$forAwait$tempResult0 = void 0;
      $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
      if ($jscomp$forAwait$tempResult0.done) {
        break;
      }
      a = $jscomp$forAwait$tempResult0.value;
      {
        bar();
      }
     }
  } catch ($jscomp$forAwait$catchErrParam0) {
    $jscomp$forAwait$errResult0 = {error:$jscomp$forAwait$catchErrParam0};
  } finally {
    try {
      if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done && ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
        await $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0);
      }
    } finally {
      if ($jscomp$forAwait$errResult0) {
        throw $jscomp$forAwait$errResult0.error;
      }
     }
   }
}
function foo() {
}
`)

	abcFunction := CodeSubTree_FindFunctionDefinition(c.GetLastCompiler(), "abc").GetRootNode()
	firstTry :=
		abcFunction.
			GetLastChild(). // block
			GetLastChild()  // try
	forNode :=
		firstTry.
			GetFirstChild(). // block
			GetSecondChild() // for
	AssertNode(forNode).HasToken(rhino.Token_FOR)
	tempIterator0 :=
		firstTry.
			GetFirstFirstChild(). // var
			GetFirstChild()       // name
	AssertNode(tempIterator0).HasToken(rhino.Token_NAME)
	// Find (0, $jscomp.makeAsyncIterator)(foo())
	makeAsyncIteratorCall := tempIterator0.GetFirstChild()
	block := forNode.GetLastChild()
	AssertNode(block).IsBlock()
	// Find $jscomp$forAwait$tempResult0 = void 0;
	resetTempResult0 := block.GetFirstFirstChild()
	AssertNode(resetTempResult0).IsAssign()
	// Find $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
	tempResult0 := block.GetSecondChild().GetFirstChild()
	AssertNode(tempResult0).IsAssign()
	await := tempResult0.GetFirstChild().GetNext()
	nextCall := await.GetFirstChild()
	done :=
		block.
			GetChildAtIndex(2). // if
			GetFirstChild()
	value :=
		block.
			GetChildAtIndex(3). // exprResult
			GetFirstChild().    // assign
			GetLastChild()      // getprop
	AssertNode(tempIterator0).
		HasColorThat().
		IsEqualTo(c.GetGlobalColor(colors.StandardColors_ASYNC_ITERATOR_ITERABLE_ID))
	AssertNode(makeAsyncIteratorCall).
		HasColorThat().
		IsEqualTo(c.GetGlobalColor(colors.StandardColors_ASYNC_ITERATOR_ITERABLE_ID))
	AssertNode(nextCall).HasColorThat().IsEqualTo(c.GetGlobalColor(colors.StandardColors_PROMISE_ID))
	AssertNode(done).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
	AssertNode(value).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)

	c.Test(
		"async function abc() { for await (var a of foo()) { bar(); } }",
		`async function abc() {
 var a$jscomp$3;
 var $jscomp$forAwait$errResult0;
 var $jscomp$forAwait$tempResult0;
 var $jscomp$forAwait$retFn0;
 try {
   var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(foo());
   for (;;) {
     $jscomp$forAwait$tempResult0 = void 0;
     $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
     if ($jscomp$forAwait$tempResult0.done) {
       break;
     }
     a$jscomp$3 = $jscomp$forAwait$tempResult0.value;
     {
       bar();
     }
    }
 } catch ($jscomp$forAwait$catchErrParam0) {
   $jscomp$forAwait$errResult0 = {error:$jscomp$forAwait$catchErrParam0};
 } finally {
   try {
     if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done && ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
       await $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0);
     }
   } finally {
     if ($jscomp$forAwait$errResult0) {
       throw $jscomp$forAwait$errResult0.error;
     }
    }
  }
}
`)

	c.Test(
		"async function abc() { for await (let a of foo()) { bar(); } }",
		`async function abc() {
 var $jscomp$forAwait$errResult0;
 var $jscomp$forAwait$tempResult0;
 var $jscomp$forAwait$retFn0;
 try {
   var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(foo());
   for (;;) {
     $jscomp$forAwait$tempResult0 = void 0;
     $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
     if ($jscomp$forAwait$tempResult0.done) {
       break;
     }
     let a$jscomp$3 = $jscomp$forAwait$tempResult0.value;
     {
       bar();
     }
    }
 } catch ($jscomp$forAwait$catchErrParam0) {
   $jscomp$forAwait$errResult0 = {error:$jscomp$forAwait$catchErrParam0};
 } finally {
   try {
     if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done && ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
       await $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0);
     }
   } finally {
     if ($jscomp$forAwait$errResult0) {
       throw $jscomp$forAwait$errResult0.error;
     }
    }
  }
}
`)

	c.Test(
		"async function abc() { for await (const a of foo()) { bar(); } }",
		`async function abc() {
 var $jscomp$forAwait$errResult0;
 var $jscomp$forAwait$tempResult0;
 var $jscomp$forAwait$retFn0;
 try {
   var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(foo());
   for (;;) {
     $jscomp$forAwait$tempResult0 = void 0;
     $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
     if ($jscomp$forAwait$tempResult0.done) {
       break;
     }
     const a$jscomp$3 = $jscomp$forAwait$tempResult0.value;
     {
       bar();
     }
    }
 } catch ($jscomp$forAwait$catchErrParam0) {
   $jscomp$forAwait$errResult0 = {error:$jscomp$forAwait$catchErrParam0};
 } finally {
   try {
     if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done && ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
       await $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0);
     }
   } finally {
     if ($jscomp$forAwait$errResult0) {
       throw $jscomp$forAwait$errResult0.error;
     }
    }
  }
}
`)
}

// Java: RewriteAsyncIterationTest.testForAwaitOfInAsyncArrow()
func (c *RewriteAsyncIterationTest) TestForAwaitOfInAsyncArrow() {
	c.Test(
		"async () => { for await (let a of foo()) { bar(); } }",
		`async() => {
  var $jscomp$forAwait$errResult0;
  var $jscomp$forAwait$tempResult0;
  var $jscomp$forAwait$retFn0;
  try {
var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(foo());
    for (;;) {
      $jscomp$forAwait$tempResult0 = void 0;
      $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
      if ($jscomp$forAwait$tempResult0.done) {
        break;
      }
      let a$jscomp$3 = $jscomp$forAwait$tempResult0.value;
      {
        bar();
      }
    }
  } catch ($jscomp$forAwait$catchErrParam0) {
    $jscomp$forAwait$errResult0 = {error:$jscomp$forAwait$catchErrParam0};
  } finally {
    try {
      if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done && ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
        await $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0);
      }
    } finally {
      if ($jscomp$forAwait$errResult0) {
        throw $jscomp$forAwait$errResult0.error;
      }
    }
  }
};
`)
}

// Java: RewriteAsyncIterationTest.testLabelledForAwaitOf()
func (c *RewriteAsyncIterationTest) TestLabelledForAwaitOf() {
	c.Test(
		`async () => {
  label:
  for await (let a of foo()) {
    bar();
  }
}
`,
		`async() => {
  var $jscomp$forAwait$errResult0;
  var $jscomp$forAwait$tempResult0;
  var $jscomp$forAwait$retFn0;
  try {
// rewriting does not lose the label with the for await of statement
  var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(foo());
    label: for (;;) {
      $jscomp$forAwait$tempResult0 = void 0;
      $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
      if ($jscomp$forAwait$tempResult0.done) {
        break;
      }
      let a$jscomp$3 = $jscomp$forAwait$tempResult0.value;
      {
        bar();
      }
    }
  } catch ($jscomp$forAwait$catchErrParam0) {
    $jscomp$forAwait$errResult0 = {error:$jscomp$forAwait$catchErrParam0};
  } finally {
    try {
      if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done && ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
        await $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0);
      }
    } finally {
      if ($jscomp$forAwait$errResult0) {
        throw $jscomp$forAwait$errResult0.error;
      }
    }
  }
};
`)
}

// Java: RewriteAsyncIterationTest.testForAwaitOfInAsyncGenerator()
func (c *RewriteAsyncIterationTest) TestForAwaitOfInAsyncGenerator() {
	c.Test(
		`async function* foo() {
  for await (let val of bar()) {
    yield* val;
  }
}
`,
		`function foo() {
  return new $jscomp.AsyncGeneratorWrapper(function*() {
    var $jscomp$forAwait$errResult0;
    var $jscomp$forAwait$tempResult0;
    var $jscomp$forAwait$retFn0;
    try {
var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(bar());
      for (;;) {
        $jscomp$forAwait$tempResult0 = void 0;
        $jscomp$forAwait$tempResult0 = yield new $jscomp.AsyncGeneratorWrapper$ActionRecord($jscomp.AsyncGeneratorWrapper$ActionEnum.AWAIT_VALUE, $jscomp$forAwait$tempIterator0.next());
        if ($jscomp$forAwait$tempResult0.done) {
          break;
        }
        let val = $jscomp$forAwait$tempResult0.value;
        {
          yield new $jscomp.AsyncGeneratorWrapper$ActionRecord($jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_STAR, val);
        }
      }
    } catch ($jscomp$forAwait$catchErrParam0) {
      $jscomp$forAwait$errResult0 = {error:$jscomp$forAwait$catchErrParam0};
    } finally {
      try {
        if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done && ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
          yield new $jscomp.AsyncGeneratorWrapper$ActionRecord($jscomp.AsyncGeneratorWrapper$ActionEnum.AWAIT_VALUE, $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0));
        }
      } finally {
        if ($jscomp$forAwait$errResult0) {
          throw $jscomp$forAwait$errResult0.error;
        }
      }
    }
  }());
}
`)
}

// Java: RewriteAsyncIterationTest.testForAwaitResetsTempResultBeforeAwaitNext()
func (c *RewriteAsyncIterationTest) TestForAwaitResetsTempResultBeforeAwaitNext() {
	c.Test(
		`async function run() {
  for await (const x of source()) {
    use(x);
  }
}
`,
		`async function run() {
  var $jscomp$forAwait$errResult0;
  var $jscomp$forAwait$tempResult0;
  var $jscomp$forAwait$retFn0;
  try {
    var $jscomp$forAwait$tempIterator0 = (0, $jscomp.makeAsyncIterator)(source());
    for (;;) {
      $jscomp$forAwait$tempResult0 = void 0;
      $jscomp$forAwait$tempResult0 = await $jscomp$forAwait$tempIterator0.next();
      if ($jscomp$forAwait$tempResult0.done) {
        break;
      }
      const x = $jscomp$forAwait$tempResult0.value;
      {
        use(x);
      }
    }
  } catch ($jscomp$forAwait$catchErrParam0) {
    $jscomp$forAwait$errResult0 = {error:$jscomp$forAwait$catchErrParam0};
  } finally {
    try {
      if ($jscomp$forAwait$tempResult0 && !$jscomp$forAwait$tempResult0.done && ($jscomp$forAwait$retFn0 = $jscomp$forAwait$tempIterator0.return)) {
        await $jscomp$forAwait$retFn0.call($jscomp$forAwait$tempIterator0);
      }
    } finally {
      if ($jscomp$forAwait$errResult0) {
        throw $jscomp$forAwait$errResult0.error;
      }
    }
  }
}
`)
}

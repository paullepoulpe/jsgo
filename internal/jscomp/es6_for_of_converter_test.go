// Port of com/google/javascript/jscomp/Es6ForOfConverterTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/Es6ForOfConverterTest.java)

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

// Es6ForOfConverterTest has the test cases for the for-of transpilation pass.
type Es6ForOfConverterTest struct {
	CompilerTestCaseBase
}

// Java: Es6ForOfConverterTest.EXTERNS_BASE
var Es6ForOfConverterTest_EXTERNS_BASE = NewTestExternsBuilder().AddArguments().AddConsole().AddJSCompLibraries().Build()

// Java: Es6ForOfConverterTest()
func NewEs6ForOfConverterTest() *Es6ForOfConverterTest {
	c := &Es6ForOfConverterTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, Es6ForOfConverterTest_EXTERNS_BASE)
	return c
}

func TestEs6ForOfConverterTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6ForOfConverterTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *Es6ForOfConverterTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// transpileTests_ImmutableMapOf is ImmutableMap.of(k1, v1, k2, v2, ...) for
// the String maps of the transpilation tests.
//
// Go-only: test-helper.
func transpileTests_ImmutableMapOf(keysAndValues ...string) *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	for i := 0; i < len(keysAndValues); i += 2 {
		jutil.CheckArgumentMsg1(!m.ContainsKey(keysAndValues[i]), "duplicate key: %s", keysAndValues[i])
		m.Put(keysAndValues[i], keysAndValues[i+1])
	}
	return m
}

// Java: Es6ForOfConverterTest.customSetUp()
func (c *Es6ForOfConverterTest) CustomSetUp() {
	c.EnableNormalize()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
	c.SetGenericNameReplacements(
		transpileTests_ImmutableMapOf(
			"KEY", "$jscomp$key$",
			"ITER", "$jscomp$iter$"))
}

// Java: Es6ForOfConverterTest.getProcessor(Compiler)
func (c *Es6ForOfConverterTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"injectTranspilationRuntimeLibraries", func(abstractCompiler AbstractCompiler) CompilerPass {
				return NewInjectTranspilationRuntimeLibraries(abstractCompiler)
			}))
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("es6ForOfConverter", func(abstractCompiler AbstractCompiler) CompilerPass { return NewEs6ForOfConverter(abstractCompiler) }))
	return optimizer
}

// Java: Es6ForOfConverterTest.testForOfLoop()
func (c *Es6ForOfConverterTest) TestForOfLoop() {
	// With array literal and declaring new bound variable.
	c.Test(
		"for (var i of [1,2,3]) { console.log(i); }",
		`var i;
var ITER$0 = (0, $jscomp.makeIterator)([1,2,3]);
var KEY$1$i = ITER$0.next();
try {
  for (;
      !KEY$1$i.done; KEY$1$i = ITER$0.next()) {
     i = KEY$1$i.value;
    {
      console.log(i);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$i);
}
`)

	// With simple assign instead of var declaration in bound variable.
	c.Test(
		"for (i of [1,2,3]) { console.log(i); }",
		`var ITER$0 = (0, $jscomp.makeIterator)([1,2,3])
var KEY$1$i = ITER$0.next();
try {
  for (;
      !KEY$1$i.done; KEY$1$i = ITER$0.next()) {
    i = KEY$1$i.value;
    {
      console.log(i);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$i);
}
`)

	// With name instead of array literal.
	c.Test(
		"for (var i of arr) { console.log(i); }",
		`var i;
var ITER$0 = (0, $jscomp.makeIterator)(arr)
var KEY$1$i = ITER$0.next();
try {
  for (;
      !KEY$1$i.done; KEY$1$i = ITER$0.next()) {
     i = KEY$1$i.value;
    {
      console.log(i);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$i);
}
`)

	// for of with const initializer
	c.Test(
		"for (const i of [1,2,3]) { console.log(i); }",
		`var ITER$0 = (0, $jscomp.makeIterator)([1,2,3]);
var KEY$1$i = ITER$0.next();
try {
  for (;
      !KEY$1$i.done; KEY$1$i = ITER$0.next()) {
    const i = KEY$1$i.value;
    {
      console.log(i);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$i);
}
`)

	// multiple for-of loops with the const initializer name
	c.Test(
		"for (const i of [1,2,3]) { console.log(i); } for (const i of [4,5,6]) { console.log(i); }",
		`var ITER$0 = (0, $jscomp.makeIterator)([1,2,3]);
var KEY$1$i = ITER$0.next();
try {
  for (;
      !KEY$1$i.done; KEY$1$i = ITER$0.next()) {
    const i = KEY$1$i.value;
    {
      console.log(i);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$i);
}
var ITER$2 = (0, $jscomp.makeIterator)([4, 5, 6]);
var KEY$3$i$jscomp$1 = ITER$2.next();
try {
  for (;
      !KEY$3$i$jscomp$1.done; KEY$3$i$jscomp$1 = ITER$2.next()) {
    const i$jscomp$1 = KEY$3$i$jscomp$1.value;
    {
      console.log(i$jscomp$1);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$2, KEY$3$i$jscomp$1);
}
`)

	// With empty loop body.
	c.Test(
		"for (var i of [1,2,3]);",
		`var i;
var ITER$0 = (0, $jscomp.makeIterator)([1,2,3])
var KEY$1$i = ITER$0.next();
try {
  for (;
      !KEY$1$i.done; KEY$1$i = ITER$0.next()) {
     i = KEY$1$i.value;
    {}
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$i);
}
`)

	// With no block in for loop body.
	c.Test(
		"for (var i of [1,2,3]) console.log(i);",
		`var i;
var ITER$0 = (0, $jscomp.makeIterator)([1,2,3]);
var KEY$1$i = ITER$0.next();
try {
  for (;
      !KEY$1$i.done; KEY$1$i = ITER$0.next()) {
     i = KEY$1$i.value;
    {
      console.log(i);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$i);
}
`)

	// Iteration var shadows an outer var ()
	c.Test(
		"var i = 'outer'; for (let i of [1, 2, 3]) { alert(i); } alert(i);",
		`var i = 'outer';
var ITER$0 = (0, $jscomp.makeIterator)([1,2,3])
var KEY$1$i$jscomp$1 = ITER$0.next();
try {
  for (;
      !KEY$1$i$jscomp$1.done; KEY$1$i$jscomp$1 = ITER$0.next()) {
    let i$jscomp$1 = KEY$1$i$jscomp$1.value;
    {
      alert(i$jscomp$1);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$i$jscomp$1);
}
alert(i);
`)
}

// Java: Es6ForOfConverterTest.testConstnessPreservedInNewDeclarations()
func (c *Es6ForOfConverterTest) TestConstnessPreservedInNewDeclarations() {
	c.Test(
		"for (let CID of [1, 2, 3]) { alert(CID); }",
		`var ITER$0 = (0, $jscomp.makeIterator)([1,2,3])
var KEY$1$CID = ITER$0.next();
try {
  for (;
      !KEY$1$CID.done; KEY$1$CID = ITER$0.next()) {
    let CID = KEY$1$CID.value;
    {
      alert(CID);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$CID);
}
`)
	script := c.GetLastCompiler().GetJsRoot().GetOnlyChild()
	jutil.CheckStateMsg1(script.IsScript(), "%s", script.GetToken())
	tryFinally := script.GetLastChild()
	jutil.CheckStateMsg1(tryFinally.IsTry(), "%s", tryFinally.GetToken())
	tryBlock := tryFinally.GetFirstChild()
	jutil.CheckStateMsg1(tryBlock.IsBlock(), "%s", tryBlock.GetToken())
	forLoop := tryBlock.GetFirstChild()
	jutil.CheckStateMsg1(forLoop.IsVanillaFor(), "%s", forLoop.GetToken())
	forBody := forLoop.GetLastChild()
	jutil.CheckState(forBody.IsBlock())
	declaration := forBody.GetFirstChild()
	jutil.CheckState(declaration.IsLet())
	name := declaration.GetFirstChild()
	jutil.CheckState(name.GetString() == "CID")
	jutil.CheckState(name.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME))
}

// Java: Es6ForOfConverterTest.testForOfRedeclaredVar()
func (c *Es6ForOfConverterTest) TestForOfRedeclaredVar() {
	c.Test(
		`for (let x of []) {
  let x = 0;
}
`,
		`var ITER$0=(0, $jscomp.makeIterator)([]);
var KEY$1$x=ITER$0.next();
try {
  for(;
      !KEY$1$x.done; KEY$1$x=ITER$0.next()) {
    let x = KEY$1$x.value;
    {
      let x$jscomp$1 = 0;
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$x);
}
`)
}

// Java: Es6ForOfConverterTest.testForOfJSDoc()
func (c *Es6ForOfConverterTest) TestForOfJSDoc() {
	c.Test(
		"for (/** @type {string} */ let x of []) {}",
		`var ITER$0=(0, $jscomp.makeIterator)([]);
var KEY$1$x=ITER$0.next();
try {
  for(;
      !KEY$1$x.done;KEY$1$x=ITER$0.next()) {
    let x = KEY$1$x.value;
    {}
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$x);
}
`)
	c.Test(
		"for (/** @type {string} */ x of []) {}",
		`var ITER$0=(0, $jscomp.makeIterator)([]);
var KEY$1$x=ITER$0.next();
try {
  for(;
      !KEY$1$x.done;KEY$1$x=ITER$0.next()) {
    x = KEY$1$x.value;
    {}
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$x);
}
`)
}

// Java: Es6ForOfConverterTest.testForOfOnNonIterable()
func (c *Es6ForOfConverterTest) TestForOfOnNonIterable() {
	c.TestWarning(
		`var arrayLike = {
  0: 'x',
  1: 'y',
  length: 2,
};
for (var x of arrayLike) {}
`,
		TypeValidator_TYPE_MISMATCH_WARNING)
}

// Java: Es6ForOfConverterTest.testLabelForOf()
func (c *Es6ForOfConverterTest) TestLabelForOf() {
	// Tests if iterator variables come before a single label
	c.Test(
		"a: for(var i of [1,2]){console.log(i)}",
		`var i;
var ITER$0 = (0, $jscomp.makeIterator)([1,2]);
var KEY$1$i = ITER$0.next();
try {
  a: for (;
      !KEY$1$i.done; KEY$1$i = ITER$0.next()) {
     i = KEY$1$i.value;
    {
      console.log(i);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$i);
}
`)
	// Test if the iterator variables come before two labels
	c.Test(
		"a: b: for(var x of [1,2]){console.log(x)}",
		`var x;
var ITER$0 = (0, $jscomp.makeIterator)([1,2]);
var KEY$1$x = ITER$0.next();
try {
  a: b: for(;
      !KEY$1$x.done; KEY$1$x = ITER$0.next()) {
     x = KEY$1$x.value;
    {
      console.log(x);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$x);
}
`)
}

// Java: Es6ForOfConverterTest.testForOfWithQualifiedNameInitializer()
func (c *Es6ForOfConverterTest) TestForOfWithQualifiedNameInitializer() {
	c.Test(
		"var obj = {a: 0}; for (obj.a of [1,2,3]) { console.log(obj.a); }",
		`var obj = {a: 0};
var ITER$0 = (0, $jscomp.makeIterator)([1,2,3])
var KEY$1$a = ITER$0.next();
try {
  for (;
      !KEY$1$a.done; KEY$1$a = ITER$0.next()) {
    obj.a = KEY$1$a.value;
    {
      console.log(obj.a);
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$a);
}
`)
}

// Java: Es6ForOfConverterTest.testForOfWithComplexInitializer()
func (c *Es6ForOfConverterTest) TestForOfWithComplexInitializer() {
	c.Test(
		"function f() { return {}; } for (f()['x' + 1] of [1,2,3]) {}",
		`function f() { return {}; }
var ITER$0 = (0, $jscomp.makeIterator)([1,2,3]);
var KEY$1$a = ITER$0.next();
try {
  for (;
      !KEY$1$a.done; KEY$1$a = ITER$0.next()) {
    f()['x' + 1] = KEY$1$a.value;
    {}
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$a);
}
`)
}

// Java: Es6ForOfConverterTest.testForLetOfWithoutExterns()
func (c *Es6ForOfConverterTest) TestForLetOfWithoutExterns() {
	c.Test(
		// add only minimal runtime library stubs to prevent AstFactory crash
		c.Externs(
			`/** @const */
var $jscomp = {};
/**
  * @param {?} iterable
  * @return {!Iterator<T>}
  * @template T
  */
$jscomp.makeIterator = function(iterable) {};
/**
  * @param {?} iterator
  * @param {?} iterResult
  * @return {void}
  */
$jscomp.iteratorClose = function(iterator, iterResult) {};
`),
		c.Srcs("for (let x of [1, 2, 3]) {}"),
		c.Expected(
			`var ITER$0 = (0, $jscomp.makeIterator)([1,2,3]);
var KEY$1$x = ITER$0.next();
try {
  for (;
      !KEY$1$x.done; KEY$1$x = ITER$0.next()) {
    let x = KEY$1$x.value;
    {}
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$x);
}
`))
}

// Java: Es6ForOfConverterTest.testForOfBreak()
func (c *Es6ForOfConverterTest) TestForOfBreak() {
	c.Test(
		"for (var x of [1, 2, 3]) { if (x === 2) break; }",
		`var x;
var ITER$0 = (0, $jscomp.makeIterator)([1, 2, 3]);
var KEY$1$x = ITER$0.next();
try {
  for (;
      !KEY$1$x.done; KEY$1$x = ITER$0.next()) {
    x = KEY$1$x.value;
    {
      if (x === 2) {
        break;
      }
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$x);
}
`)
}

// Java: Es6ForOfConverterTest.testForOfReturn()
func (c *Es6ForOfConverterTest) TestForOfReturn() {
	c.Test(
		"function f() { for (var x of [1, 2, 3]) { if (x === 2) return x; } }",
		`function f() {
  var x;
  var ITER$0 = (0, $jscomp.makeIterator)([1, 2, 3]);
  var KEY$1$x = ITER$0.next();
  try {
    for (;
        !KEY$1$x.done; KEY$1$x = ITER$0.next()) {
      x = KEY$1$x.value;
      {
        if (x === 2) {
          return x;
        }
      }
    }
  } finally {
    $jscomp.iteratorClose(ITER$0, KEY$1$x);
  }
}
`)
}

// Java: Es6ForOfConverterTest.testForOfThrow()
func (c *Es6ForOfConverterTest) TestForOfThrow() {
	c.Test(
		"for (var x of [1, 2, 3]) { if (x === 2) throw new Error(); }",
		`var x;
var ITER$0 = (0, $jscomp.makeIterator)([1, 2, 3]);
var KEY$1$x = ITER$0.next();
try {
  for (;
      !KEY$1$x.done; KEY$1$x = ITER$0.next()) {
    x = KEY$1$x.value;
    {
      if (x === 2) {
        throw new Error();
      }
    }
  }
} finally {
  $jscomp.iteratorClose(ITER$0, KEY$1$x);
}
`)
}

// Port of com/google/javascript/jscomp/Es6ExtractClassesTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/Es6ExtractClassesTest.java)

/*
 * Copyright 2016 The Closure Compiler Authors.
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

type Es6ExtractClassesTest struct {
	CompilerTestCaseBase
}

// Java: Es6ExtractClassesTest()
func NewEs6ExtractClassesTest() *Es6ExtractClassesTest {
	c := &Es6ExtractClassesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestEs6ExtractClassesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6ExtractClassesTest() })
}

// Java: Es6ExtractClassesTest.getProcessor(Compiler)
func (c *Es6ExtractClassesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs6NormalizeClasses(compiler)
}

// Java: Es6ExtractClassesTest.setUp()
func (c *Es6ExtractClassesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)

	c.EnableTypeInfoValidation()
	c.EnableRewriteClosureCode()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
}

// Java: Es6ExtractClassesTest.testExtractionFromCall()
func (c *Es6ExtractClassesTest) TestExtractionFromCall() {
	c.Test(
		"f(class{});",
		`const CLASS_DECL$0 = class {};
f(CLASS_DECL$0);
`)
}

// Java: Es6ExtractClassesTest.testSelfReference1()
func (c *Es6ExtractClassesTest) TestSelfReference1() {
	c.Test(
		`var Outer = class Inner {
  constructor() {
    alert(Inner);
  }
};
`,
		`var Outer = class {
  constructor() {
    alert(Outer);
  }
};
`)

	c.Test(
		`let Outer = class Inner {
  constructor() {
    alert(Inner);
  }
};
`,
		`let Outer = class {
  constructor() {
    alert(Outer);
  }
};
`)

	c.Test(
		`const Outer = class Inner {
  constructor() {
    alert(Inner);
  }
};
`,
		`const Outer = class {
  constructor() {
    alert(Outer);
  }
};
`)
}

// Java: Es6ExtractClassesTest.testSelfReference2()
func (c *Es6ExtractClassesTest) TestSelfReference2() {
	c.Test(
		`alert(class C {
  constructor() {
    alert(C);
  }
});
`,
		`const CLASS_DECL$0 = class {
  constructor() {
    alert(CLASS_DECL$0);
  }
};
alert(CLASS_DECL$0)
`)
}

// Java: Es6ExtractClassesTest.testSelfReference3()
func (c *Es6ExtractClassesTest) TestSelfReference3() {
	c.Test(
		`alert(class C {
  m1() {
    class C {};
    alert(C);
  }
  m2() {
    alert(C);
  }
});
`,
		`const CLASS_DECL$0 = class {
  m1() {
    class C {};
    alert(C);
  }
  m2() {
    alert(CLASS_DECL$0);
  }
};
alert(CLASS_DECL$0)
`)
}

// Java: Es6ExtractClassesTest.testSelfReference_googModule()
func (c *Es6ExtractClassesTest) TestSelfReference_googModule() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddClosureExterns().Build()),
		c.Srcs(
			`goog.module('example');
exports = class Inner { constructor() { alert(Inner); } };
`),
		c.Expected(
			`/** @const */
var module$exports$example = class {
  constructor() {
    alert(module$exports$example);
  }
};
`))
}

// Java: Es6ExtractClassesTest.testSelfReference_qualifiedName()
func (c *Es6ExtractClassesTest) TestSelfReference_qualifiedName() {
	c.Test(
		`const outer = {};
/** @const */
outer.qual = {};
outer.qual.Name = class Inner {
  constructor() {
    alert(Inner);
  }
};
`,
		`const outer = {};
/** @const */
outer.qual = {};
outer.qual.Name = class {
  constructor() {
    alert(outer.qual.Name);
  }
};
`)
}

// Java: Es6ExtractClassesTest.testConstAssignment()
func (c *Es6ExtractClassesTest) TestConstAssignment() {
	c.Test(
		"var foo = bar(class {});",
		`const CLASS_DECL$0 = class {};
var foo = bar(CLASS_DECL$0);
`)
}

// Java: Es6ExtractClassesTest.testLetAssignment()
func (c *Es6ExtractClassesTest) TestLetAssignment() {
	c.Test(
		"let foo = bar(class {});",
		`const CLASS_DECL$0 = class {};
let foo = bar(CLASS_DECL$0);
`)
}

// Java: Es6ExtractClassesTest.testVarAssignment()
func (c *Es6ExtractClassesTest) TestVarAssignment() {
	c.Test(
		"var foo = bar(class {});",
		`const CLASS_DECL$0 = class {};
var foo = bar(CLASS_DECL$0);
`)
}

// Java: Es6ExtractClassesTest.testJSDocOnVar()
func (c *Es6ExtractClassesTest) TestJSDocOnVar() {
	c.Test(
		"/** @unrestricted */ var foo = class bar {};",
		`var foo = class {};
`)
}

// Java: Es6ExtractClassesTest.testFilenameContainsAt()
func (c *Es6ExtractClassesTest) TestFilenameContainsAt() {
	c.Test(
		c.Srcs(SourceFile_FromCode("unusual@name", "alert(class {});")),
		c.Expected(
			SourceFile_FromCode(
				"unusual@name",
				`const CLASS_DECL$0 = class{};
alert(CLASS_DECL$0);
`)))
}

// Java: Es6ExtractClassesTest.testFilenameContainsPlus()
func (c *Es6ExtractClassesTest) TestFilenameContainsPlus() {
	c.Test(
		c.Srcs(SourceFile_FromCode("+some/+path/file", "alert(class {});")),
		c.Expected(
			SourceFile_FromCode(
				"+path/file",
				`const CLASS_DECL$0 = class{};
alert(CLASS_DECL$0);
`)))
}

// Java: Es6ExtractClassesTest.testConditionalBlocksExtractionFromCall()
func (c *Es6ExtractClassesTest) TestConditionalBlocksExtractionFromCall() {
	c.Test(
		"maybeTrue() && f(class{});",
		`maybeTrue() && f((() => {
  const CLASS_DECL$0 = class {};
  return CLASS_DECL$0;
})());
`)
}

// Java: Es6ExtractClassesTest.testExtractionFromArrayLiteral()
func (c *Es6ExtractClassesTest) TestExtractionFromArrayLiteral() {
	c.Test(
		"var c = [class C {}];",
		`const CLASS_DECL$0 = class {};
var c = [CLASS_DECL$0];
`)
}

// Java: Es6ExtractClassesTest.testTernaryOperatorBlocksExtraction()
func (c *Es6ExtractClassesTest) TestTernaryOperatorBlocksExtraction() {
	c.Test(
		"var c = maybeTrue() ? class A {} : anotherExpr",
		`var c = maybeTrue() ? (() => {
  const CLASS_DECL$0 = class {};
  return CLASS_DECL$0;
})() : anotherExpr;
`)
	c.Test(
		"var c = maybeTrue() ? anotherExpr : class B {}",
		`var c = maybeTrue() ? anotherExpr : (() => {
  const CLASS_DECL$0 = class {};
  return CLASS_DECL$0;
})();
`)
}

// Java: Es6ExtractClassesTest.testNormalisedArrowFunction()
func (c *Es6ExtractClassesTest) TestNormalisedArrowFunction() {
	c.Test(
		"(() => { return class A {};})();",
		"(() => { const CLASS_DECL$0 = class {}; return CLASS_DECL$0;})()")

	c.Test(
		"function foo(x = () => {return class A {}; }) {} use(foo());",
		`function foo(
  x = () => {
    const CLASS_DECL$0 = class {};
    return CLASS_DECL$0;
  })
 {}
use(foo());
`)
}

// Tests that class that is a default parameter value are converted into normalized arrow
// functions and successfully extracted
//
// Java: Es6ExtractClassesTest.testClassAsDefaultParamValue()
func (c *Es6ExtractClassesTest) TestClassAsDefaultParamValue() {
	c.Test(
		`function foo(x = class A {}) {}
use(foo());
`,
		`function foo(x =
  (() => { const CLASS_DECL$0 = class {};
           return CLASS_DECL$0;
         }
  )()) {}
use(foo());
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: Es6ExtractClassesTest.classExtractedInNormalizedArrowNested()
func (c *Es6ExtractClassesTest) TestClassExtractedInNormalizedArrowNested() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddClosureExterns().Build()),
		c.Srcs(
			`goog.module('some');
exports.some = ((outer) => {
  ((c) => { return class extends c{}}
)});
`),
		c.Expected(
			`/** @const */ var module$exports$some = {};
/** @const */ module$exports$some.some = outer => {
c => {
  const CLASS_DECL$0 = class extends c {};
  return CLASS_DECL$0;
  };
};
`))
}

// Java: Es6ExtractClassesTest.testCannotExtract()
func (c *Es6ExtractClassesTest) TestCannotExtract() {
	c.Test(
		"let x; while(x = class A{}) {use(x);}",
		`let x;
while(x = (() => { const CLASS_DECL$0 = class {}; return CLASS_DECL$0;})())
{use(x);}
`)

	c.Test(
		`/** @type {number} */ var x = 0;
function f(x, y) {}
while(f(x = 2, class Foo { bar() {} })){}
`,
		`var x = 0;
function f(x, y) {}
while(
  f(x = 2,
    (() => { const CLASS_DECL$0 = class { bar() {} };
          return CLASS_DECL$0;
         })()
)){}
`)
}

// Java: Es6ExtractClassesTest.testClassesHandledByEs6ToEs3Converter()
func (c *Es6ExtractClassesTest) TestClassesHandledByEs6ToEs3Converter() {
	c.TestSame("class C{}")
	c.TestSame("var c = class {};")
}

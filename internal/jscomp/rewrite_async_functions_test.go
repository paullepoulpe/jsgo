// Port of com/google/javascript/jscomp/RewriteAsyncFunctionsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RewriteAsyncFunctionsTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// RewriteAsyncFunctionsTest has the test cases for RewriteAsyncFunctions.
type RewriteAsyncFunctionsTest struct {
	CompilerTestCaseBase
}

// Java: RewriteAsyncFunctionsTest.EXTERNS_BASE
var RewriteAsyncFunctionsTest_EXTERNS_BASE = NewTestExternsBuilder().AddArguments().AddJSCompLibraries().Build()

// Java: RewriteAsyncFunctionsTest.REPLACEMENTS_MAP
var RewriteAsyncFunctionsTest_REPLACEMENTS_MAP = transpileTests_ImmutableMapOf(
	"ASYNC_THIS",
	"$jscomp$async$this$",
	"ASYNC_ARGUMENTS",
	"$jscomp$async$arguments$",
	"ASYNC_NEW_TARGET",
	"$jscomp$async$new$target$",
	"ASYNC_SUPER_GET",
	"$jscomp$async$super$get$")

// Java: RewriteAsyncFunctionsTest()
func NewRewriteAsyncFunctionsTest() *RewriteAsyncFunctionsTest {
	c := &RewriteAsyncFunctionsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, RewriteAsyncFunctionsTest_EXTERNS_BASE)
	return c
}

func TestRewriteAsyncFunctionsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteAsyncFunctionsTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *RewriteAsyncFunctionsTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: RewriteAsyncFunctionsTest.customSetUp()
func (c *RewriteAsyncFunctionsTest) CustomSetUp() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.EnableTypeCheck()
	c.EnableNormalize()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
	c.SetGenericNameReplacements(RewriteAsyncFunctionsTest_REPLACEMENTS_MAP)
}

// Java: RewriteAsyncFunctionsTest.getProcessor(Compiler)
func (c *RewriteAsyncFunctionsTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"injectTranspilationRuntimeLibraries", func(abstractCompiler AbstractCompiler) CompilerPass {
				return NewInjectTranspilationRuntimeLibraries(abstractCompiler)
			}))
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory("rewriteAsyncFunctions", func(abstractCompiler AbstractCompiler) CompilerPass {
			return RewriteAsyncFunctions_Create(abstractCompiler)
		}))
	return optimizer
}

// Java: RewriteAsyncFunctionsTest.getGlobalColor(ColorId)
func (c *RewriteAsyncFunctionsTest) GetGlobalColor(colorId colors.ColorId) *colors.Color {
	return c.GetLastCompiler().GetColorRegistry().Get(colorId)
}

// Java: RewriteAsyncFunctionsTest.getGlobalInstanceColor(String)
func (c *RewriteAsyncFunctionsTest) GetGlobalInstanceColor(globalClassName string) *colors.Color {
	return colors.Color_CreateUnionSlice(
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), globalClassName).
			GetRootNode().
			GetColor().
			GetInstanceColors())
}

// Java: RewriteAsyncFunctionsTest.testDefaultParameterUsingThis()
func (c *RewriteAsyncFunctionsTest) TestDefaultParameterUsingThis() {
	c.Test(
		`class X {
  /**
   * @param {number} a
   */
  constructor(a) {
    /** @const */ this.a = a;
  }
  /**
   * @param {number} b
   * @return {!Promise<number>}
   */
  async m(b = this.a) {
      return this.a + b;
  }
}
`,
		`class X {
  constructor(a) {
    /** @const */ this.a = a;
  }
  m(b = this.a) { // this in parameter default value doesn't get changed
    const ASYNC_THIS$3 = this;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
            return ASYNC_THIS$3.a + b;
        });
  }
}
`)

	classXInstanceType := c.GetGlobalInstanceColor("X")

	thisAliasNameReferences :=
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "X").
			FindMethodDefinition("m").
			FindMatchingQNameReferences("$jscomp$async$this$m1146332801$3")
	AssertThat(thisAliasNameReferences).HasSize(2)

	// const ASYNC_THIS$3 = this;
	// confirm that `this` and `ASYNC_THIS$3` nodes have the right types in
	// declaration
	aliasDeclarationReference := thisAliasNameReferences[0]
	AssertNode(aliasDeclarationReference).HasColorThat().IsEqualTo(classXInstanceType)
	thisNode := aliasDeclarationReference.GetOnlyChild()
	AssertNode(thisNode).IsThis().HasColorThat().IsEqualTo(classXInstanceType)

	// make sure the single reference to ASYNC_THIS$3 has the right type
	AssertNode(thisAliasNameReferences[1]).HasColorThat().IsEqualTo(classXInstanceType)
}

// Java: RewriteAsyncFunctionsTest.testInnerArrowFunctionUsingThis()
func (c *RewriteAsyncFunctionsTest) TestInnerArrowFunctionUsingThis() {
	c.Test(
		`class X {
  async m() {
    return new Promise((resolve, reject) => {
      return this;
    });
  }
}
`,
		`class X {
  m() {
    const ASYNC_THIS$1 = this;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return new Promise((resolve, reject) => {
            return ASYNC_THIS$1;
          });
        });
  }
}
`)

	classXInstanceType := c.GetGlobalInstanceColor("X")

	thisAliasNameReferences :=
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "X").
			FindMethodDefinition("m").
			FindMatchingQNameReferences("$jscomp$async$this$m1146332801$1")
	AssertThat(thisAliasNameReferences).HasSize(2)

	// const ASYNC_THIS$1 = this;
	// confirm that `this` and `ASYNC_THIS$1` nodes have the right types in
	// declaration
	aliasDeclarationReference := thisAliasNameReferences[0]
	AssertNode(aliasDeclarationReference).HasColorThat().IsEqualTo(classXInstanceType)
	thisNode := aliasDeclarationReference.GetOnlyChild()
	AssertNode(thisNode).IsThis().HasColorThat().IsEqualTo(classXInstanceType)

	// make sure the single reference to ASYNC_THIS$1 has the right type
	AssertNode(thisAliasNameReferences[1]).HasColorThat().IsEqualTo(classXInstanceType)
}

// Java: RewriteAsyncFunctionsTest.testInnerSuperCall()
func (c *RewriteAsyncFunctionsTest) TestInnerSuperCall() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddPromise().AddJSCompLibraries().Build()),
		c.Srcs(
			`class A {
  m() {
    return Promise.resolve(this);
  }
}
class X extends A {
  async m() {
    return super.m();
  }
}
`),
		c.Expected(
			`class A {
  m() {
    return Promise.resolve(this);
  }
}
class X extends A {
  m() {
    const ASYNC_THIS$3 = this;
    const ASYNC_SUPER_GET$5$m = () => {
      return super.m;
    };
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return ASYNC_SUPER_GET$5$m().call(ASYNC_THIS$3);
        });
  }
}
`))

	classAInstanceType := c.GetGlobalInstanceColor("A")
	// type of A.prototype.m
	classAPropertyMType :=
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").
			FindMethodDefinition("m").
			GetRootNode().
			GetColor()

	classXMethodMDefinition :=
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "X").FindMethodDefinition("m")

	// Check type information on wrapper function for `super.m`
	superMethodWrapperNameNodes :=
		classXMethodMDefinition.FindMatchingQNameReferences(
			"$jscomp$async$super$get$m1146332801$5$m")
	// one declaration and one reference
	AssertThat(superMethodWrapperNameNodes).HasSize(2)

	// first name node is declaration
	// const ASYNC_SUPER_GET$m = () => super.m;
	wrapperDeclarationNameNode := superMethodWrapperNameNodes[0]
	wrapperArrowFunction := wrapperDeclarationNameNode.GetOnlyChild()
	// optimization colors don't track function signatures
	AssertNode(wrapperArrowFunction).
		IsArrowFunction().
		HasColorThat().
		IsEqualTo(colors.StandardColors_TOP_OBJECT)
	// wrapper function variable has type matching the function itself
	wrapperArrowColor := wrapperArrowFunction.GetColor()
	AssertNode(wrapperDeclarationNameNode).HasColorThat().IsEqualTo(wrapperArrowColor)

	// get `super.m` from `() => `super.m`
	superDotM := wrapperArrowFunction.GetLastChild().GetFirstFirstChild()
	AssertNode(superDotM).
		MatchesQualifiedName("super.m").
		HasColorThat().
		IsEqualTo(classAPropertyMType)
	superNode := superDotM.GetFirstChild()
	AssertNode(superNode).IsSuper().HasColorThat().IsEqualTo(classAInstanceType)

	// second name node is reference
	// return ASYNC_SUPER_GET$m().call(ASYNC_THIS$1);
	wrapperReferenceNameNode := superMethodWrapperNameNodes[1]
	// optimization colors don't track function signatures
	AssertNode(wrapperArrowFunction).HasColorThat().IsEqualTo(colors.StandardColors_TOP_OBJECT)
	// `ASYNC_SUPER_GET$m()`
	wrapperCallNode := wrapperReferenceNameNode.GetParent()
	AssertNode(wrapperCallNode).IsCall().HasColorThat().IsEqualTo(classAPropertyMType)

	// `ASYNC_SUPER_GET$m().call(ASYNC_THIS$1)`
	methodCallNode := wrapperCallNode.GetGrandparent()
	// optimization colors don't track .call types
	AssertNode(methodCallNode).IsCall().HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
}

// Java: RewriteAsyncFunctionsTest.testInnerSuperTaggedTemplateLiteral()
func (c *RewriteAsyncFunctionsTest) TestInnerSuperTaggedTemplateLiteral() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddPromise().AddJSCompLibraries().Build()),
		c.Srcs(
			"class A {\n  m(strings, ...values) {\n    return Promise.resolve(this);\n  }\n}\nclass X extends A {\n  async m() {\n    return super.m`<p>hello</p>`;\n  }\n}\n"),
		c.Expected(
			"class A {\n  m(strings, ...values) {\n    return Promise.resolve(this);\n  }\n}\nclass X extends A {\n  m() {\n    const ASYNC_THIS$3 = this;\n    const ASYNC_SUPER_GET$5$m = () => {\n      return super.m;\n    };\n    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(\n        function* () {\n          return ASYNC_SUPER_GET$5$m().bind(ASYNC_THIS$3)`<p>hello</p>`;\n        });\n  }\n}\n"))
}

// Java: RewriteAsyncFunctionsTest.testInnerSuperReference()
func (c *RewriteAsyncFunctionsTest) TestInnerSuperReference() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddFunction().AddJSCompLibraries().Build()),
		c.Srcs(
			`class A {
  m() {
    return this;
  }
}
class X extends A {
  async m() {
    const tmp = super.m;
    return tmp.call(null);
  }
}
`),
		c.Expected(
			"class A {\n  m() {\n    return this;\n  }\n}\nclass X extends A {\n  m() {\n    const ASYNC_SUPER_GET$5$m = () => {\n      return super.m;\n    };\n    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(\n        function* () {\n          const tmp = ASYNC_SUPER_GET$5$m();\n// type of tmp will indicate it requires `this` be provided, but will allow null.\n          return tmp.call(null);\n        });\n  }\n}\n"))

	// type of A.prototype.m
	classAPropertyMType :=
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "A").
			FindMethodDefinition("m").
			GetRootNode().
			GetColor()
	classAInstanceType := c.GetGlobalInstanceColor("A")

	classXMethodMDefinition :=
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "X").FindMethodDefinition("m")

	// Check type information on wrapper function for `super.m`
	superMethodWrapperNameNodes :=
		classXMethodMDefinition.FindMatchingQNameReferences(
			"$jscomp$async$super$get$m1146332801$5$m")
	// one declaration and one reference
	AssertThat(superMethodWrapperNameNodes).HasSize(2)

	// first name node is declaration
	// const ASYNC_SUPER_GET$m = () => super.m;
	wrapperDeclarationNameNode := superMethodWrapperNameNodes[0]
	// arrow function has a Color representing a object
	wrapperArrowFunction := wrapperDeclarationNameNode.GetOnlyChild()
	AssertNode(wrapperArrowFunction).
		IsArrowFunction().
		HasColorThat().
		IsEqualTo(colors.StandardColors_TOP_OBJECT)
	// wrapper function variable has type matching the function itself
	wrapperArrowColor := wrapperArrowFunction.GetColor()
	AssertNode(wrapperDeclarationNameNode).HasColorThat().IsEqualTo(wrapperArrowColor)

	// get `super.m` from `() => `super.m`
	superDotM := wrapperArrowFunction.GetLastChild().GetFirstFirstChild()
	AssertNode(superDotM).
		MatchesQualifiedName("super.m").
		HasColorThat().
		IsEqualTo(classAPropertyMType)
	superNode := superDotM.GetFirstChild()
	AssertNode(superNode).HasColorThat().IsEqualTo(classAInstanceType)

	// second name node is reference
	// const tmp = ASYNC_SUPER_GET$m();
	wrapperReferenceNameNode := superMethodWrapperNameNodes[1]
	// optimization colors don't track function signatures
	AssertNode(wrapperReferenceNameNode).HasColorThat().IsEqualTo(colors.StandardColors_TOP_OBJECT)
	// `ASYNC_SUPER_GET$m()`
	wrapperCallNode := wrapperReferenceNameNode.GetParent()
	AssertNode(wrapperCallNode).IsCall().HasColorThat().IsEqualTo(classAPropertyMType)
}

// Java: RewriteAsyncFunctionsTest.testMultipleSuperAccessesInAsyncFunction_havingNonIdenticalUnknownTypes()
func (c *RewriteAsyncFunctionsTest) TestMultipleSuperAccessesInAsyncFunction_havingNonIdenticalUnknownTypes() {
	c.Test(
		"class UpdatingElement {\n  getUpdateComplete() {\n  }\n}\n\nclass TextFieldBase extends UpdatingElement {\n  async _getUpdateComplete() {\n    if (super.getUpdateComplete) { // `?` type\n      await super.getUpdateComplete(); // `??` type\n    }\n  }\n}\n",
		`class UpdatingElement {
  getUpdateComplete() {
  }
}
class TextFieldBase extends UpdatingElement {
  _getUpdateComplete() {
    const ASYNC_THIS$3 = this;
    const ASYNC_SUPER_GET$5$getUpdateComplete = () => {
      return super.getUpdateComplete;
    };
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(function*() {
      if (ASYNC_SUPER_GET$5$getUpdateComplete()) {
        yield ASYNC_SUPER_GET$5$getUpdateComplete().call(ASYNC_THIS$3);
      }
    });
  }
}
`)
}

// Java: RewriteAsyncFunctionsTest.testNestedArrowFunctionUsingThis()
func (c *RewriteAsyncFunctionsTest) TestNestedArrowFunctionUsingThis() {
	c.Test(
		`class X {
  m() {
    return async () => (() => this);
  }
}
`,
		`class X {
  m() {
    return () => {
      const ASYNC_THIS$3 = this;
      return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
          function* () {
            return () => {
              return ASYNC_THIS$3;
            };
          })
    }
  }
}
`)
}

// Java: RewriteAsyncFunctionsTest.testInnerArrowFunctionUsingArguments()
func (c *RewriteAsyncFunctionsTest) TestInnerArrowFunctionUsingArguments() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddArguments().AddJSCompLibraries().Build()),
		c.Srcs(
			`class X {
  async m() {
    return new Promise((resolve, reject) => {
      return arguments;
    });
  }
}
`),
		c.Expected(
			`class X {
  m() {
    const ASYNC_ARGUMENTS$1 = arguments;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return new Promise((resolve, reject) => {
            return ASYNC_ARGUMENTS$1
          });
        });
  }
}
`))

	argumentsAliasRefs :=
		CodeSubTree_FindClassDefinition(c.GetLastCompiler(), "X").
			FindMethodDefinition("m").
			FindMatchingQNameReferences("$jscomp$async$arguments$m1146332801$1")
	AssertThat(argumentsAliasRefs).HasSize(2) // one declaration and 1 use

	argumentsColor := c.GetGlobalColor(colors.StandardColors_ARGUMENTS_ID)

	// declaration reference
	// const ASYNC_ARGUMENTS$1 = arguments;
	argumentsAliasDeclaration := argumentsAliasRefs[0]
	argumentsValue := argumentsAliasDeclaration.GetOnlyChild()
	AssertNode(argumentsValue).
		MatchesQualifiedName("arguments").
		HasColorThat().
		IsEqualTo(argumentsColor)
	AssertNode(argumentsAliasDeclaration).
		MatchesQualifiedName("$jscomp$async$arguments$m1146332801$1").
		HasColorThat().
		IsEqualTo(argumentsColor)

	// usage reference
	// return ASYNC_ARGUMENTS$1;
	argumentsAliasUsage := argumentsAliasRefs[1]
	AssertNode(argumentsAliasUsage).
		MatchesQualifiedName("$jscomp$async$arguments$m1146332801$1").
		HasColorThat().
		IsEqualTo(argumentsColor)
}

// Java: RewriteAsyncFunctionsTest.testAwaitReplacement()
func (c *RewriteAsyncFunctionsTest) TestAwaitReplacement() {
	c.Test(
		"async function foo(promise) { return await promise; }",
		`function foo(promise) {
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return yield promise;
        });
}
`)
}

// Java: RewriteAsyncFunctionsTest.testArgumentsReplacement_topLevelCode()
func (c *RewriteAsyncFunctionsTest) TestArgumentsReplacement_topLevelCode() {
	c.TestSame("arguments;")
}

// Java: RewriteAsyncFunctionsTest.testArgumentsReplacement_normalFunction()
func (c *RewriteAsyncFunctionsTest) TestArgumentsReplacement_normalFunction() {
	c.TestSame("function f(a, b, ...rest) { return arguments.length; }")
}

// Java: RewriteAsyncFunctionsTest.testArgumentsReplacement_asyncFunction()
func (c *RewriteAsyncFunctionsTest) TestArgumentsReplacement_asyncFunction() {
	c.Test(
		"async function f(a, b, ...rest) { return arguments.length; }",
		`function f(a, b, ...rest) {
  const ASYNC_ARGUMENTS$1 = arguments;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return ASYNC_ARGUMENTS$1.length; // arguments replaced
        });
}
`)
}

// Java: RewriteAsyncFunctionsTest.testArgumentsReplacement_asyncClosure()
func (c *RewriteAsyncFunctionsTest) TestArgumentsReplacement_asyncClosure() {
	c.Test(
		`function outer() {
  /**
   * @param {...?} varArgs
   * @return {!Promise<number>}
   */
  async function f(varArgs) { return arguments.length; }
  return f(arguments)
}
`,
		`function outer() {
  function f(varArgs) {
    const ASYNC_ARGUMENTS$3 = arguments;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return ASYNC_ARGUMENTS$3.length; // arguments replaced
        });
  }
  return f(arguments) // unchanged
}
`)
}

// Java: RewriteAsyncFunctionsTest.testArgumentsReplacement_normalClosureInAsync()
func (c *RewriteAsyncFunctionsTest) TestArgumentsReplacement_normalClosureInAsync() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddFunction().AddJSCompLibraries().Build()),
		c.Srcs(
			`async function a() {
  function inner() {
    return arguments.length;
  }
  return inner.apply(undefined, arguments); // this should get replaced
}
`),
		c.Expected(
			`function a() {
  const ASYNC_ARGUMENTS$1 = arguments;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          function inner() {
            return arguments.length; // unchanged
          }
          return inner.apply(undefined, ASYNC_ARGUMENTS$1);
        });
}
`))
}

// Java: RewriteAsyncFunctionsTest.testClassMethod()
func (c *RewriteAsyncFunctionsTest) TestClassMethod() {
	c.Test(
		`class A {
  /**
   * @param {number} x
   */
  constructor(x) {
    /** @type {number} */ this.x = x;
  }
  async f() {
    return this.x;
  }
}
`,
		`class A {
  constructor(x) {
    this.x = x;
  }
  f() {
    const ASYNC_THIS$3 = this;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function *() {
          return ASYNC_THIS$3.x; // this replaced
        });
  }
}
`)
}

// Java: RewriteAsyncFunctionsTest.testAsyncClassMethodWithAsyncArrow()
func (c *RewriteAsyncFunctionsTest) TestAsyncClassMethodWithAsyncArrow() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddConsole().AddJSCompLibraries().Build()),
		c.Srcs(
			`class A {
  async f() {
    let g = async () => { console.log(this, arguments); };
    g();
  }
}
`),
		c.Expected(
			`class A {
  f() {
    const ASYNC_THIS$1 = this;
    const ASYNC_ARGUMENTS$1 = arguments;
      return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
          function *() {
            let g = () => {
              return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
                  function *() {
                    console.log(ASYNC_THIS$1, ASYNC_ARGUMENTS$1);
                  });
            };
            g();
          });
  }
}
`))
}

// Java: RewriteAsyncFunctionsTest.testNonAsyncClassMethodWithAsyncArrow()
func (c *RewriteAsyncFunctionsTest) TestNonAsyncClassMethodWithAsyncArrow() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddConsole().AddJSCompLibraries().Build()),
		c.Srcs(
			`class A {
  f() {
    let g = async () => { console.log(this, arguments); };
    g();
  }
}
`),
		c.Expected(
			`class A {
  f() {
    let g = () => {
      const ASYNC_THIS$3 = this;
      const ASYNC_ARGUMENTS$3 = arguments;
      return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
          function *() {
            console.log(ASYNC_THIS$3, ASYNC_ARGUMENTS$3);
          });
    };
    g();
  }
}
`))
}

// Java: RewriteAsyncFunctionsTest.testArrowFunctionExpressionBody()
func (c *RewriteAsyncFunctionsTest) TestArrowFunctionExpressionBody() {
	c.Test(
		"let f = async () => 1;",
		`let f = () => {
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return 1;
        });
}
`)
}

// Java: RewriteAsyncFunctionsTest.testGlobalScopeArrowFunctionRefersToThis()
func (c *RewriteAsyncFunctionsTest) TestGlobalScopeArrowFunctionRefersToThis() {
	c.Test(
		"let f = async () => this;",
		`let f = () => {
    const ASYNC_THIS$1 = this;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return ASYNC_THIS$1;
        });
}
`)
}

// Java: RewriteAsyncFunctionsTest.testGlobalScopeAsyncArrowFunctionDefaultParamValueRefersToThis()
func (c *RewriteAsyncFunctionsTest) TestGlobalScopeAsyncArrowFunctionDefaultParamValueRefersToThis() {
	c.Test(
		"let f = async (t = this) => t;",
		`let f = (t = this) => {
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return t;
        });
}
`)
}

// Java: RewriteAsyncFunctionsTest.testNestedAsyncArrowFunctionDefaultParamValueRefersToThis()
func (c *RewriteAsyncFunctionsTest) TestNestedAsyncArrowFunctionDefaultParamValueRefersToThis() {
	c.Test(
		`let f = async function(outerT = this) {
  return async (t = this) => t;
};
`,
		"let f = function(outerT = this) {\n  const ASYNC_THIS$1 = this;\n  return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(\n      function* () {\n// `this` is aliased here\n        return (t = ASYNC_THIS$1) => {\n          return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(\n              function* () {\n                return t;\n              });\n        };\n      });\n};\n")
}

// Java: RewriteAsyncFunctionsTest.testSuperPropAssignment()
func (c *RewriteAsyncFunctionsTest) TestSuperPropAssignment() {
	c.TestError(
		`class A {}
class B extends A {
  async f() {
    super.x = 5;
  }
}
`,
		TranspilationUtil_CANNOT_CONVERT_YET,
		"Transpilation of 'assignment to super property' is not yet implemented.")
}

// Java: RewriteAsyncFunctionsTest.testAsyncArrowUsingNewTargetInConstructor()
func (c *RewriteAsyncFunctionsTest) TestAsyncArrowUsingNewTargetInConstructor() {
	c.Test(
		`class X {
  constructor() {
    const f = async () => {
      return new.target;
    };
  }
}
`,
		`class X {
  constructor() {
    const f = () => {
      const ASYNC_NEW_TARGET$3 = new.target;
      return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
          function* () {
            return ASYNC_NEW_TARGET$3;
          });
    };
  }
}
`)
}

// Java: RewriteAsyncFunctionsTest.testAsyncArrowUsingNewTargetInFunction()
func (c *RewriteAsyncFunctionsTest) TestAsyncArrowUsingNewTargetInFunction() {
	c.Test(
		`function Foo() {
  const f = async () => {
    return new.target;
  };
}
`,
		`function Foo() {
  const f = () => {
    const ASYNC_NEW_TARGET$3 = new.target;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return ASYNC_NEW_TARGET$3;
        });
  };
}
`)
}

// Java: RewriteAsyncFunctionsTest.testNestedArrowFunctionUsingNewTarget()
func (c *RewriteAsyncFunctionsTest) TestNestedArrowFunctionUsingNewTarget() {
	c.Test(
		`function Foo() {
  const f = async () => (() => new.target);
}
`,
		`function Foo() {
  const f = () => {
    const ASYNC_NEW_TARGET$3 = new.target;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return () => {
            return ASYNC_NEW_TARGET$3;
          };
        });
  };
}
`)
}

// Java: RewriteAsyncFunctionsTest.testAsyncFunctionUsingNewTarget()
func (c *RewriteAsyncFunctionsTest) TestAsyncFunctionUsingNewTarget() {
	c.Test(
		`async function f() {
  return new.target;
}
`,
		`function f() {
  const ASYNC_NEW_TARGET$1 = new.target;
  return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
      function* () {
        return ASYNC_NEW_TARGET$1;
      });
}
`)
}

// Java: RewriteAsyncFunctionsTest.testAsyncClassMethodWithAsyncArrowUsingNewTarget()
func (c *RewriteAsyncFunctionsTest) TestAsyncClassMethodWithAsyncArrowUsingNewTarget() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddConsole().AddJSCompLibraries().Build()),
		c.Srcs(
			`class A {
  async f() {
    let g = async () => { console.log(new.target); };
    g();
  }
}
`),
		c.Expected(
			`class A {
  f() {
    const ASYNC_NEW_TARGET$1 = new.target;
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          let g = () => {
            return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
                function* () {
                  console.log(ASYNC_NEW_TARGET$1);
                });
          };
          g();
        });
  }
}
`))
}

// Java: RewriteAsyncFunctionsTest.testNonAsyncClassMethodWithAsyncArrowUsingNewTarget()
func (c *RewriteAsyncFunctionsTest) TestNonAsyncClassMethodWithAsyncArrowUsingNewTarget() {
	c.Test(
		c.Externs(NewTestExternsBuilder().AddConsole().AddJSCompLibraries().Build()),
		c.Srcs(
			`class A {
  f() {
    let g = async () => { console.log(new.target); };
    g();
  }
}
`),
		c.Expected(
			`class A {
  f() {
    let g = () => {
      const ASYNC_NEW_TARGET$3 = new.target;
      return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
          function* () {
            console.log(ASYNC_NEW_TARGET$3);
          });
    };
    g();
  }
}
`))
}

// Java: RewriteAsyncFunctionsTest.testDefaultParamArrowUsingThisArgumentsNewTarget()
func (c *RewriteAsyncFunctionsTest) TestDefaultParamArrowUsingThisArgumentsNewTarget() {
	c.Test(
		`async function f(a = () => this, b = () => arguments, c = () => new.target) {
  return a();
}
`,
		`function f(a = () => {
  return this;
}, b = () => {
  return arguments;
}, c = () => {
  return new.target;
}) {
  return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
      function* () {
        return a();
      });
}
`)
}

// Java: RewriteAsyncFunctionsTest.testDefaultParamArrowUsingSuperProp()
func (c *RewriteAsyncFunctionsTest) TestDefaultParamArrowUsingSuperProp() {
	c.Test(
		`class A {
  prop() {
    return 1;
  }
}
class B extends A {
  async m(a = () => super.prop()) {
    return a();
  }
}
`,
		`class A {
  prop() {
    return 1;
  }
}
class B extends A {
  m(a = () => {
    return super.prop();
  }) {
    return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(
        function* () {
          return a();
        });
  }
}
`)
}

// Java: RewriteAsyncFunctionsTest.testDefaultParamAsyncArrowUsingThis()
func (c *RewriteAsyncFunctionsTest) TestDefaultParamAsyncArrowUsingThis() {
	c.Test(
		`async function f(a = async () => this) {
  return a();
}
`,
		`function f(a = () => {
  const ASYNC_THIS$3 = this;
  return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(function*() {
    return ASYNC_THIS$3;
  });
}) {
  return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(function*() {
    return a();
  });
}
`)
}

// Java: RewriteAsyncFunctionsTest.testDefaultParamArrowInNestedAsyncArrowUsingThis()
func (c *RewriteAsyncFunctionsTest) TestDefaultParamArrowInNestedAsyncArrowUsingThis() {
	c.Test(
		`async function f() {
  const g = async (a = () => this) => a();
  return g();
}
`,
		`function f() {
  const ASYNC_THIS$1 = this;
  return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(function*() {
    const g = (a = () => {
      return ASYNC_THIS$1;
    }) => {
      return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(function*() {
        return a();
      });
    };
    return g();
  });
}
`)
}

// Java: RewriteAsyncFunctionsTest.testDefaultParamArrowInNonAsyncFunctionInAsyncFunction()
func (c *RewriteAsyncFunctionsTest) TestDefaultParamArrowInNonAsyncFunctionInAsyncFunction() {
	c.Test(
		`async function f() {
  function g(a = () => this) {
    return a();
  }
  return g();
}
`,
		`function f() {
  return (0, $jscomp.asyncExecutePromiseGeneratorFunction)(function*() {
    function g(a = () => {
      return this;
    }) {
      return a();
    }
    return g();
  });
}
`)
}

// Java: RewriteAsyncFunctionsTest.testObjectLiteralAsyncMethodSuperCall()
func (c *RewriteAsyncFunctionsTest) TestObjectLiteralAsyncMethodSuperCall() {
	// TODO(b/572972833): Object-literal async method super.m() passes null enclosing class to
	// AstFactory, causing NPE.
	e :=
		AssertThrows[jutil.IRuntimeException](
			func() {
				c.TestSame(
					`const obj = {
  async m() {
    return super.toString();
  }
};
`)
			})
	// Throwables.getRootCause(e)
	var root jutil.Throwable = e
	for root.GetCause() != nil {
		root = root.GetCause()
	}
	IsInstanceOf[jutil.INullPointerException](AssertThat(root))
	AssertThat(root.GetMessage()).Contains("Node.isClass()")
}

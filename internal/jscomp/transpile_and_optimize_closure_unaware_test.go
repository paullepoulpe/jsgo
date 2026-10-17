// Port of com/google/javascript/jscomp/TranspileAndOptimizeClosureUnawareTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/TranspileAndOptimizeClosureUnawareTest.java)

/*
 * Copyright 2025 The Closure Compiler Authors.
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
	"os"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type TranspileAndOptimizeClosureUnawareTest struct {
	CompilerTestCaseBase

	inputCount   int
	outputCount  int
	extraOptions func(*CompilerOptions)
	mode         NestedCompilerRunner_Mode
}

func NewTranspileAndOptimizeClosureUnawareTest() *TranspileAndOptimizeClosureUnawareTest {
	c := &TranspileAndOptimizeClosureUnawareTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestTranspileAndOptimizeClosureUnawareTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTranspileAndOptimizeClosureUnawareTest() })
}

// Java: TranspileAndOptimizeClosureUnawareTest.getProcessor(Compiler)
func (c *TranspileAndOptimizeClosureUnawareTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewTranspileAndOptimizeClosureUnaware(compiler, c.mode)
}

// SetUp runs CompilerTestCase.setUp and then the @Before method setup().
func (c *TranspileAndOptimizeClosureUnawareTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.setup()
}

// Java: TranspileAndOptimizeClosureUnawareTest.setup()
func (c *TranspileAndOptimizeClosureUnawareTest) setup() {
	c.EnableMultistageCompilation()
	c.DisableCompareJsDoc()
	c.inputCount = 0
	c.outputCount = 0
	c.extraOptions = func(options *CompilerOptions) {}
	c.mode = NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE
}

// Java: TranspileAndOptimizeClosureUnawareTest.getOptions()
func (c *TranspileAndOptimizeClosureUnawareTest) GetOptions() *CompilerOptions {
	baseOptions := c.CompilerTestCaseBase.GetOptions()
	c.extraOptions(baseOptions)
	return baseOptions
}

// Java: TranspileAndOptimizeClosureUnawareTest.createPrettyPrinter(Node)
func (c *TranspileAndOptimizeClosureUnawareTest) createPrettyPrinter(n *rhino.Node) string {
	return NewCodePrinter_Builder(n).
		SetCompilerOptions(c.GetLastCompiler().GetOptions()).
		SetPrettyPrint(true).
		Build()
}

// Java: TranspileAndOptimizeClosureUnawareTest.testEmptyClosureUnawareFn_unchanged()
func (c *TranspileAndOptimizeClosureUnawareTest) TestEmptyClosureUnawareFn_unchanged() {
	c.Test(c.Srcs(c.closureUnaware("")), c.Expected(c.expectedClosureUnaware("")))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testFoldConstants()
func (c *TranspileAndOptimizeClosureUnawareTest) TestFoldConstants() {
	c.Test(
		c.Srcs(c.closureUnaware("console.log(1 + 2);")),
		c.Expected(c.expectedClosureUnaware("console.log(3);")))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testFoldConstants_multipleFiles()
func (c *TranspileAndOptimizeClosureUnawareTest) TestFoldConstants_multipleFiles() {
	c.Test(
		c.Srcs(
			c.closureUnaware("console.log(1 + 2);"),
			c.closureUnaware("console.log(11 + 22);"),
			c.closureUnaware("console.log(111 + 222);")),
		c.Expected(
			c.expectedClosureUnaware("console.log(3);"),
			c.expectedClosureUnaware("console.log(33);"),
			c.expectedClosureUnaware("console.log(333);")))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testFoldConstants_multipleFiles_ignoresNonClosureUnaware()
func (c *TranspileAndOptimizeClosureUnawareTest) TestFoldConstants_multipleFiles_ignoresNonClosureUnaware() {
	c.Test(
		c.Srcs(
			c.closureUnaware("console.log(1 + 2);"),
			`goog.module('some.other.file');
console.log(11 + 22);
`),
		c.Expected(
			c.expectedClosureUnaware("console.log(3);"),
			`goog.module('some.other.file');
console.log(11 + 22);
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testFoldConstants_multipleClosureUnawareBlocksInFile()
func (c *TranspileAndOptimizeClosureUnawareTest) TestFoldConstants_multipleClosureUnawareBlocksInFile() {
	c.Test(
		c.Srcs(
			c.closureUnaware(
				"console.log(1 + 2);", "console.log(11 + 22);", "console.log(111 + 222);")),
		c.Expected(
			c.expectedClosureUnaware("console.log(3);", "console.log(33);", "console.log(333);")))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testFoldConstants_regressionTest_crashInNodeUtilAddFeatureToScript()
func (c *TranspileAndOptimizeClosureUnawareTest) TestFoldConstants_regressionTest_crashInNodeUtilAddFeatureToScript() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)

	c.Test(
		c.closureUnaware(
			`return function f() {
  function referenceX() { return x; }

  return 1;

  // Regression test for a crash in PeepholeRemoveDeadCode.
  let x = 0;
  referenceX();
  x++;
}
`),
		c.expectedClosureUnaware(
			`return function() {
  return 1;
}
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testInlineLocalObject()
func (c *TranspileAndOptimizeClosureUnawareTest) TestInlineLocalObject() {
	c.Test(
		c.Srcs(
			c.closureUnaware(
				`const obj = {x: 1, unused: 2};
return obj.x;
`)),
		c.Expected(
			c.expectedClosureUnaware(
				`return 1;
`)))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testDoNotRenameOrFlattenEscapedObjects()
func (c *TranspileAndOptimizeClosureUnawareTest) TestDoNotRenameOrFlattenEscapedObjects() {
	c.Test(
		c.Srcs(
			c.closureUnaware(
				`const obj = {x: 1, y: 2};
foo(obj);
return obj.x;
`)),
		c.Expected(
			c.expectedClosureUnaware(
				`const a = {x: 1, y: 2};
foo(a);
return a.x;
`)))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testAssumeAllPropertyReadsHaveSideEffects()
func (c *TranspileAndOptimizeClosureUnawareTest) TestAssumeAllPropertyReadsHaveSideEffects() {
	c.Test(
		c.Srcs(
			c.closureUnaware(
				`return function(obj) {
  const x = obj.x;
  // Assume that reading obj.y might have other side effects & affect obj.x.
  const y = obj.y;
  console.log(x);
}
`)),
		c.Expected(
			c.expectedClosureUnaware(
				`return function(a) {
  const b = a.x;
  a.y;
  console.log(b);
}
`)))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testAssumeAllPropertyWritesHaveSideEffects()
func (c *TranspileAndOptimizeClosureUnawareTest) TestAssumeAllPropertyWritesHaveSideEffects() {
	c.Test(
		c.closureUnaware(
			`return function(obj) {
  const x = obj.x;
  // Assume that writing to obj.y might have other side effects & affect obj.x.
  obj.y = 0;
  console.log(x);
}
`),
		c.expectedClosureUnaware(
			`return function(a) {
  const b = a.x;
  a.y = 0;
  console.log(b);
}
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testNoClosureAssertsRemoval()
func (c *TranspileAndOptimizeClosureUnawareTest) TestNoClosureAssertsRemoval() {
	c.Test(
		c.closureUnaware("goog.asserts.assert(true);"),
		c.expectedClosureUnaware("goog.asserts.assert(!0);"))
}

// Java: TranspileAndOptimizeClosureUnawareTest.doesntRenameExternRefs_shadowedLocally()
func (c *TranspileAndOptimizeClosureUnawareTest) TestDoesntRenameExternRefs_shadowedLocally() {
	c.Test(
		c.closureUnaware(
			`return function() {
  const maybeExternal = () => typeof External === 'undefined' ? null : External;
  globalThis.foo = function() {
    const External = maybeExternal();
    console.log(External);
  };
};
`),
		c.expectedClosureUnaware(
			// Regression test: the compiler used to incorrectly obfuscate 'typeof External'
			// to 'typeof a'.
			`return function() {
  globalThis.foo = function() {
  const a = typeof External === "undefined" ? null : External;
  console.log(a);
  };
};
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testAssumeRegexpMayHaveGlobalSideEffects()
func (c *TranspileAndOptimizeClosureUnawareTest) TestAssumeRegexpMayHaveGlobalSideEffects() {
	// References to the global RegExp object make regular expression calls unoptimizable. In the
	// main AST the CheckRegExp pass looks for these references. But in @closureUnaware compilation,
	// we cannot assume that the CheckRegexpPass will see all possible Regexp usages.
	// Hence we cannot remove the seemingly unused `/abc/.match(str)` call.
	// Note that we still remove completely unused regex literals.
	c.Test(c.closureUnaware("/abc/.match(str); /unused/"), c.expectedClosureUnaware("/abc/.match(str);"))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testDontAssumeToStringAndValueOfArePure()
func (c *TranspileAndOptimizeClosureUnawareTest) TestDontAssumeToStringAndValueOfArePure() {
	c.TestSame(c.closureUnaware("x.toString(); x.valueOf();"))
}

// Java: TranspileAndOptimizeClosureUnawareTest.transpilesExponentialOperator_ifTargetingES2015()
func (c *TranspileAndOptimizeClosureUnawareTest) TestTranspilesExponentialOperator_ifTargetingES2015() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.Test(
		c.Srcs(c.closureUnaware("console.log(x ** y);")),
		c.Expected(c.expectedClosureUnaware("console.log(Math.pow(x, y));")))
}

// Java: TranspileAndOptimizeClosureUnawareTest.doesNotTranspileExponentialOperator_ifTargetingES2017()
func (c *TranspileAndOptimizeClosureUnawareTest) TestDoesNotTranspileExponentialOperator_ifTargetingES2017() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2017)
	c.Test(
		c.Srcs(c.closureUnaware("console.log(x ** y);")),
		c.Expected(c.expectedClosureUnaware("console.log(x ** y);")))
}

// Java: TranspileAndOptimizeClosureUnawareTest.transpileOnlyMode_transpilesExponentialOperator_ifTargettingES2015()
func (c *TranspileAndOptimizeClosureUnawareTest) TestTranspileOnlyMode_transpilesExponentialOperator_ifTargettingES2015() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.mode = NestedCompilerRunner_Mode_TRANSPILE_ONLY
	c.DisableCompareJsDoc()
	c.Test(
		c.Srcs(c.closureUnaware("console.log(x ** y);")),
		c.Expected(c.expectedClosureUnaware("console.log(Math.pow(x, y));")))
}

// Java: TranspileAndOptimizeClosureUnawareTest.transpileOnlyMode_doesNotTranspileExponentialOperator_ifTargetingES2017()
func (c *TranspileAndOptimizeClosureUnawareTest) TestTranspileOnlyMode_doesNotTranspileExponentialOperator_ifTargetingES2017() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2017)
	c.mode = NestedCompilerRunner_Mode_TRANSPILE_ONLY
	c.DisableCompareJsDoc()
	c.Test(
		c.Srcs(c.closureUnaware("console.log(x ** y);")),
		c.Expected(c.expectedClosureUnaware("console.log(x ** y);")))
}

// Java: TranspileAndOptimizeClosureUnawareTest.transpileOnlyMode_doesNotRunConstantFolding()
func (c *TranspileAndOptimizeClosureUnawareTest) TestTranspileOnlyMode_doesNotRunConstantFolding() {
	c.mode = NestedCompilerRunner_Mode_TRANSPILE_ONLY
	c.DisableCompareJsDoc()
	c.Test(
		c.Srcs(c.closureUnaware("console.log(1 + 2);")),
		c.Expected(c.expectedClosureUnaware("console.log(1 + 2);")))
}

// Java: TranspileAndOptimizeClosureUnawareTest.classInheritance_injectsRequiredRuntimeLibraries()
func (c *TranspileAndOptimizeClosureUnawareTest) TestClassInheritance_injectsRequiredRuntimeLibraries() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		c.Srcs(
			`/** @fileoverview @closureUnaware */
goog.module('test');
/** @closureUnaware */
(function() {
  class Parent {
    method() {
      return 10;
    }
  }
  class Child extends Parent {}
  return new Child();
}).call(globalThis);
`),
		// NOTE: we could make transpilation smart enough to detect that we don't need the
		// $jscomp.construct call here. Currently it only attempts to analyze superclasses
		// defined in the global scope, and @closureUnaware code is never global. We could
		// easily make the pass look up the definition of 'Parent' in the function-local scope, but
		// decided not to as of Jan 2026, as it isn't clearly worth the added complexity.
		c.Expected(
			`/** @fileoverview @closureUnaware */
goog.module("test");
/** @closureUnaware */
(function(c, d) {
  /** @constructor */
  var a = function() {};
  a.prototype.method = function() {
    return 10;
  };
  /** @constructor */
  var b = function() {
    return c(a, arguments, this.constructor);
  };
  d(b, a);
  return new b();
}).call(globalThis, $jscomp.construct, $jscomp.inherits);
`))

	AssertThat(c.GetLastCompiler().GetRuntimeJsLibManager().GetInjectedLibraries()).
		ContainsExactly("es6/util/construct", "es6/util/inherits", "es6/util/arrayfromiterable")
}

// Java: TranspileAndOptimizeClosureUnawareTest.runtimeLibraryInjection_injectsOnlyWhatsNeededPerScript()
func (c *TranspileAndOptimizeClosureUnawareTest) TestRuntimeLibraryInjection_injectsOnlyWhatsNeededPerScript() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		c.Srcs(
			c.closureUnaware("return {...obj};"),
			c.closureUnaware("return [...arr];"),
			c.closureUnaware(
				`class Parent {}
class Child extends Parent {}
return new Child();
`)),
		c.Expected(
			// Note that Object.assign doesn't need to be injected because it's a polyfill, not a
			// runtime library field.
			`/** @fileoverview @closureUnaware */
goog.module("test0");
/** @closureUnaware */
(function() {
  return Object.assign({}, obj);
}).call(undefined);
`,
			`/** @fileoverview @closureUnaware */
goog.module("test1");
/** @closureUnaware */
(function(a) {
  return [].concat(a(arr));
}).call(undefined, $jscomp.arrayFromIterable);
`,
			`/** @fileoverview @closureUnaware */
goog.module("test2");
/** @closureUnaware */
(function(a, d) {
  var b = function() {}, c = function() {
    return a(b, arguments, this.constructor);
  };
  d(c, b);
  return new c();
}).call(undefined, $jscomp.construct, $jscomp.inherits);
`))

	AssertThat(c.GetLastCompiler().GetRuntimeJsLibManager().GetInjectedLibraries()).
		ContainsExactly(
			"es6/util/construct",
			"es6/util/inherits",
			"es6/util/arrayfromiterable",
			"es6/object/assign")
}

// Java: TranspileAndOptimizeClosureUnawareTest.classInheritance_usesReflectConstructForUnknownSuperclass()
func (c *TranspileAndOptimizeClosureUnawareTest) TestClassInheritance_usesReflectConstructForUnknownSuperclass() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		c.Srcs(
			`/** @fileoverview @closureUnaware */
goog.module('test');
/** @closureUnaware */
(function() {
  class Child extends UnknownParent {}
  return new Child();
}).call(globalThis);
`),
		c.Expected(
			`/** @fileoverview @closureUnaware */
goog.module("test");
/** @closureUnaware */
(function(b, c) {
  /** @constructor */
  var a = function() {
    return b(UnknownParent, arguments, this.constructor);
  };
  c(a, UnknownParent);
  return new a();
}).call(globalThis, $jscomp.construct, $jscomp.inherits);
`))

	AssertThat(c.GetLastCompiler().GetRuntimeJsLibManager().GetInjectedLibraries()).
		ContainsExactly("es6/util/construct", "es6/util/inherits", "es6/util/arrayfromiterable")
}

// Java: TranspileAndOptimizeClosureUnawareTest.runtimeFieldInjection_injectsAfterExistingParameters()
func (c *TranspileAndOptimizeClosureUnawareTest) TestRuntimeFieldInjection_injectsAfterExistingParameters() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		c.Srcs(
			`/** @fileoverview @closureUnaware */
goog.module('test');

/** @closureUnaware */
(function(foo, bar, rest) {
  return [foo, bar, ...rest];
}).call(undefined, 1, 2, [3, 4, 5]);
`),
		c.Expected(
			`/** @fileoverview @closureUnaware */
goog.module("test");
/** @closureUnaware */
(function(a, b, c, d) {
  return [a, b].concat(d(c));
}).call(undefined, 1, 2, [3, 4, 5], $jscomp.arrayFromIterable);
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.runtimeFieldInjection_incompatibleWithRestArg()
func (c *TranspileAndOptimizeClosureUnawareTest) TestRuntimeFieldInjection_incompatibleWithRestArg() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.mode = NestedCompilerRunner_Mode_TRANSPILE_ONLY

	// To make injecting runtime library parameters simpler, assume that closure unaware functions
	// never have rest parameters.
	parts :=
		c.Srcs(
			`/** @fileoverview @closureUnaware */
goog.module('test');

/** @closureUnaware */
(function(...rest) {
  return async function foo() {};
}).call(undefined, 1, 2);
`)
	ex := AssertThrows[jutil.IIllegalStateException](func() { c.Test(parts) })
	AssertThat(ex.GetMessage()).
		Contains("rest parameters not allowed in closureUnaware function")
}

// Java: TranspileAndOptimizeClosureUnawareTest.doesNotAutomaticallyInjectPolyfills()
func (c *TranspileAndOptimizeClosureUnawareTest) TestDoesNotAutomaticallyInjectPolyfills() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2019)
	// Use `globalThis` inside @closureUnaware code. The compiler will typically detect that
	// `globalThis` is defined in ES2020 so is not present in the ES2019 output, and then add
	// a polyfill for `globalThis`.
	// But we do not support this automatic polyfill detection in @closureUnaware code, as it's
	// more difficult to statically analyze. (Particularly in the case of polyfilled methods).
	c.TestSame(
		c.Srcs(
			`/** @fileoverview @closureUnaware */
goog.module('test');
/** @closureUnaware */
(function() {
  return globalThis.FOO;
}).call(globalThis);
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.taggedTemplateLitTranspilation()
func (c *TranspileAndOptimizeClosureUnawareTest) TestTaggedTemplateLitTranspilation() {
	c.extraOptions = func(options *CompilerOptions) { options.SetGeneratePseudoNames(true) }
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		c.Srcs(
			c.closureUnaware(
				`function tag(strings, ...values) {
  return {strings, values};
}
tag`+"`"+`a${42}b`+"`"+`;
tag`+"`"+`a${42}b`+"`"+`; // identical template lit
tag`+"`"+`a${-42}b`+"`"+`; // identical template strings, different substitution content
`)),
		c.Expected(
			`/** @fileoverview @closureUnaware */
goog.module('test0');
/** @closureUnaware */
(function($$jscomp_getRestArguments$$, $$jscomp_createTemplateTagFirstArg$$) {
  function $tag$$($strings$$) {
    var $values$$ = $$jscomp_getRestArguments$$.apply(1, arguments);
    return {strings:$strings$$, values:$values$$};
  }
  /** @noinline */
  var $$jscomp$templatelit$m854121055$0$$ =
          $$jscomp_createTemplateTagFirstArg$$(["a", "b"]),
      $$jscomp$templatelit$m854121055$1$$ =
          $$jscomp_createTemplateTagFirstArg$$(["a", "b"]),
      $$jscomp$templatelit$m854121055$2$$ =
          $$jscomp_createTemplateTagFirstArg$$(["a", "b"]);
  $tag$$($$jscomp$templatelit$m854121055$0$$, 42);
  $tag$$($$jscomp$templatelit$m854121055$1$$, 42);
  $tag$$($$jscomp$templatelit$m854121055$2$$, -42);
}).call(undefined, $jscomp.getRestArguments, $jscomp.createTemplateTagFirstArg);
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.bigintLiteral_passedThroughUntranspiled()
func (c *TranspileAndOptimizeClosureUnawareTest) TestBigintLiteral_passedThroughUntranspiled() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(c.closureUnaware("return 123n + 456n;"), c.expectedClosureUnaware("return 579n;"))
}

// Java: TranspileAndOptimizeClosureUnawareTest.untranspilableRegexFeature_passedThroughUntranspiled()
func (c *TranspileAndOptimizeClosureUnawareTest) TestUntranspilableRegexFeature_passedThroughUntranspiled() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2015)

	// ES2018 feature: https://github.com/tc39/proposal-regexp-dotall-flag
	// ES2022 feature: https://github.com/tc39/proposal-regexp-match-indices
	c.TestSame(c.closureUnaware("return /a/sd"))
}

// Java: TranspileAndOptimizeClosureUnawareTest.computedGetterSetter_inEs5Out_passedThroughUntranspiled()
func (c *TranspileAndOptimizeClosureUnawareTest) TestComputedGetterSetter_inEs5Out_passedThroughUntranspiled() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	// JSCompiler does not support transpiling computed getters/setters in object literals down to
	// ES5. For closureUnaware, verify they are passed through unchanged.
	c.TestSame(
		c.closureUnaware(
			`return {
  get [x()]() {
    return 1;
  },
  set [x()](a) {}
};
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.generator_cannotConvertCaseStatementContainingYield_errors()
func (c *TranspileAndOptimizeClosureUnawareTest) TestGenerator_cannotConvertCaseStatementContainingYield_errors() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.TestError(
		c.closureUnaware("function *f(b, i) {switch (i) { case yield: return b; }} return f;"),
		TranspilationUtil_CANNOT_CONVERT_YET)
}

// Java: TranspileAndOptimizeClosureUnawareTest.loadFile(Path)
func transpileAndOptimizeClosureUnawareTest_loadFile(path *jutil.Path) string {
	data, err := os.ReadFile(path.String())
	if err != nil {
		panic(jutil.NewAssertionErrorWithMessageAndCause("Failed to load debug log at "+path.String(), jutil.AsThrowable(err)))
	}
	// Files.lines(path).collect(joining("\n")): the lines without their terminators.
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	return strings.Join(jutil.StringSplitLiteral(strings.TrimSuffix(text, "\n"), "\n"), "\n")
}

// Java: TranspileAndOptimizeClosureUnawareTest.listFiles(Path)
func (c *TranspileAndOptimizeClosureUnawareTest) listFiles(dir *jutil.Path) []*jutil.Path {
	entries, err := os.ReadDir(dir.String())
	if err != nil {
		panic(jutil.NewAssertionErrorWithCause(jutil.AsThrowable(err)))
	}
	var files []*jutil.Path
	for _, e := range entries {
		files = append(files, dir.ResolveString(e.Name()))
	}
	return files
}

// Java: TranspileAndOptimizeClosureUnawareTest.testFoldConstants_directCallInsteadOfDotCall()
func (c *TranspileAndOptimizeClosureUnawareTest) TestFoldConstants_directCallInsteadOfDotCall() {
	c.Test(
		c.Srcs(
			// Don't use the "c.closureUnaware()" wrapper because it always uses .call - we want to
			// specifically test a regular direct call.
			`/** @fileoverview @closureUnaware */
goog.module('test');
/** @closureUnaware */
(function() {
  console.log(1 + 2);
})();
`),
		c.Expected(
			`/** @fileoverview @closureUnaware */
goog.module('test');
/** @closureUnaware */
(function() {
  console.log(3);
})();
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testSupportsClosureUnawareTag_simpleOptimizations()
func (c *TranspileAndOptimizeClosureUnawareTest) TestSupportsClosureUnawareTag_simpleOptimizations() {
	c.Test(
		c.Srcs(
			`/**
 * @fileoverview
 * @closureUnaware {SIMPLE}
 */
goog.module('test');
/** @closureUnaware */
(function() {
  console.log(1 + 2);
}).call(undefined);
`),
		c.Expected(
			`/**
 * @fileoverview
 * @closureUnaware {SIMPLE}
 */
goog.module('test');
/** @closureUnaware */
(function() {
  console.log(3);
}).call(undefined);
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.testSupportsClosureUnawareTag_whitespaceOnly()
func (c *TranspileAndOptimizeClosureUnawareTest) TestSupportsClosureUnawareTag_whitespaceOnly() {
	c.TestSame(
		c.Srcs(
			`/**
 * @fileoverview
 * @closureUnaware {WHITESPACE}
 */
goog.module('test');
/** @closureUnaware */
(function() {
  console.log(1 + 2);
}).call(undefined);
`))
}

// Java: TranspileAndOptimizeClosureUnawareTest.closureUnaware(String...)
func (c *TranspileAndOptimizeClosureUnawareTest) closureUnaware(closureUnaware ...string) string {
	idx := c.inputCount
	c.inputCount++
	return transpileAndOptimizeClosureUnawareTest_buildCode(idx, closureUnaware...)
}

// Java: TranspileAndOptimizeClosureUnawareTest.expectedClosureUnaware(String...)
func (c *TranspileAndOptimizeClosureUnawareTest) expectedClosureUnaware(closureUnaware ...string) string {
	idx := c.outputCount
	c.outputCount++
	return transpileAndOptimizeClosureUnawareTest_buildCode(idx, closureUnaware...)
}

// Java: TranspileAndOptimizeClosureUnawareTest.buildCode(int, String...)
func transpileAndOptimizeClosureUnawareTest_buildCode(idx int, closureUnaware ...string) string {
	prefix :=
		jutil.StringFormat(
			`/** @fileoverview @closureUnaware */
goog.module('test%d');
`,
			idx)
	output := jutil.NewStringBuilder().Append(prefix)
	for _, block := range closureUnaware {
		output.Append(
			jutil.StringFormat(
				`/** @closureUnaware */
(function() {
  %s
}).call(undefined);
`,
				block))
	}
	return output.String()
}

// TODO: b/421971366 - add a greater variety of test cases.

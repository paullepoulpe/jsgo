// Port of com/google/javascript/jscomp/ManageClosureUnawareCodeTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ManageClosureUnawareCodeTest.java)

/*
 * Copyright 2024 The Closure Compiler Authors.
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

// ManageClosureUnawareCodeTest holds tests for ManageClosureUnawareCode,
// specifically the wrapping part.
type ManageClosureUnawareCodeTest struct {
	CompilerTestCaseBase

	runUnwrapPass          bool
	gatheredShadowNodeRoot rhino.Node
	// nil is Optional.empty().
	languageInOverride *CompilerOptions_LanguageMode
}

// Java: ManageClosureUnawareCodeTest()
func NewManageClosureUnawareCodeTest() *ManageClosureUnawareCodeTest {
	c := &ManageClosureUnawareCodeTest{
		runUnwrapPass:          true,
		gatheredShadowNodeRoot: rhino.IR_Root(),
		languageInOverride:     nil,
	}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestManageClosureUnawareCodeTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewManageClosureUnawareCodeTest() })
}

// Java: ManageClosureUnawareCodeTest.setUp()
func (c *ManageClosureUnawareCodeTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.AllowExternsChanges()
	// Wrapping and unwrapping closureUnaware code does not need to record AST changes, as it is
	// conceptually really part of parsing/pre-code-printing. (We just use the CompilerPass API &
	// CompilerTestCase class for convenience.)
	c.DisableValidateAstChangeMarking()
	c.DisableCompareJsDoc()
}

// Java: ManageClosureUnawareCodeTest.tearDown()
func (c *ManageClosureUnawareCodeTest) TearDown() {
	c.CompilerTestCaseBase.TearDown()
	c.runUnwrapPass = true
	c.gatheredShadowNodeRoot = rhino.IR_Root()
}

// Java: ManageClosureUnawareCodeTest.getOptions()
func (c *ManageClosureUnawareCodeTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	if c.languageInOverride != nil {
		options.SetLanguageIn(*c.languageInOverride)
	}
	return options
}

// Java: ManageClosureUnawareCodeTest.getProcessor(Compiler)
func (c *ManageClosureUnawareCodeTest) GetProcessor(compiler *Compiler) CompilerPass {
	phaseopt := NewPhaseOptimizer(compiler, nil)
	passes := jutil.NewArrayList[*PassFactory]()
	passes.Add(
		PassFactory_Builder_().
			SetName("gatherShadowNodes").
			SetInternalFactory(
				func(x AbstractCompiler) CompilerPass {
					return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
						NodeUtil_VisitPreOrder(
							root,
							NodeUtil_VisitorFunc(func(node rhino.Node) {
								shadow := node.GetClosureUnawareShadow()
								if !shadow.IsNil() {
									if !c.runUnwrapPass {
										node.SetClosureUnawareShadow(rhino.NilNode)
									}
									c.gatheredShadowNodeRoot.AddChildToBack(shadow.CloneTree())
								}
							}))
					})
				}).
			Build())

	if c.runUnwrapPass {
		passes.Add(
			PassFactory_Builder_().
				SetName("unwrapClosureUnawareCode").
				SetInternalFactory(func(x AbstractCompiler) CompilerPass { return ManageClosureUnawareCode_Unwrap(x) }).
				Build())
	}
	phaseopt.Consume(passes.Slice())
	return phaseopt
}

// Java: ManageClosureUnawareCodeTest.doTest(String, String, List<String>)
func (c *ManageClosureUnawareCodeTest) DoTest(js string, expectedWrapped string, expectedShadowNodeContents []string) {
	c.runUnwrapPass = false // Validate that only wrapping results in the expected wrapped contents.
	c.Test(js, expectedWrapped)
	// expectedShadowNodeContents.stream().map(s -> parseExpectedJs(s).detach())
	//     .collect(IR::root, Node::addChildToBack, Node::addChildToBack)
	expectedShadowNodeRoot := rhino.IR_Root()
	for _, s := range expectedShadowNodeContents {
		expectedShadowNodeRoot.AddChildToBack(c.ParseExpectedJs(s).Detach())
	}
	AssertNode(c.gatheredShadowNodeRoot).
		UsingSerializer(
			func(n rhino.Node) string {
				return NewCodePrinter_Builder(n).
					SetCompilerOptions(c.GetLastCompiler().GetOptions()).
					SetPrettyPrint(true).
					Build()
			}).
		IsEqualIncludingJsDocTo(expectedShadowNodeRoot)

	// now test with unwrapping enabled so it is a no-op
	c.runUnwrapPass = true
	jsWithoutClosureUnware := jutil.StringReplace(js, "@closureUnaware", "")
	c.Test(js, jsWithoutClosureUnware)
}

// Java: ManageClosureUnawareCodeTest.testDirectLoad()
func (c *ManageClosureUnawareCodeTest) TestDirectLoad() {
	c.DoTest(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
/** @closureUnaware */
(function() {
    window['foo'] = 5;
  }).call(globalThis);
`,
		`/** @fileoverview @closureUnaware */
goog.module('foo.bar.baz_raw');
$jscomp_wrap_closure_unaware_code.call(globalThis)
`,
		[]string{
			`$jscomp_sink_closure_unaware_impl(
function() {
  window['foo'] = 5;
});
`})
}

// Java: ManageClosureUnawareCodeTest.testDirectLoadWithRequireAndExports()
func (c *ManageClosureUnawareCodeTest) TestDirectLoadWithRequireAndExports() {
	c.DoTest(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
goog.require('foo.bar');
const {a} = goog.require('foo.baz');
/** @closureUnaware */
(function() {
  window['foo'] = 5;
}).call(globalThis);
exports = globalThis['foo'];
`,
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
goog.require('foo.bar');
const {a} = goog.require('foo.baz');
$jscomp_wrap_closure_unaware_code.call(globalThis)
exports = globalThis['foo'];
`,
		[]string{
			`$jscomp_sink_closure_unaware_impl(
  function() {
    window['foo'] = 5;
  });
`})
}

// Java: ManageClosureUnawareCodeTest.testConditionalLoad()
func (c *ManageClosureUnawareCodeTest) TestConditionalLoad() {
	c.DoTest(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (!window['foo']) {
  /** @closureUnaware */
  (function() {
    window['foo'] = 5;
  }).call(globalThis);
}
`,
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (!window['foo']) {
  $jscomp_wrap_closure_unaware_code.call(globalThis)
}
`,
		[]string{
			`$jscomp_sink_closure_unaware_impl(
  function() {
    window['foo'] = 5;
  });
`})
}

// Java: ManageClosureUnawareCodeTest.testDebugSrcLoad()
func (c *ManageClosureUnawareCodeTest) TestDebugSrcLoad() {
	c.DoTest(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (goog.DEBUG) {
  /** @closureUnaware */
  (function() {
    window['foo'] = 5;
  }).call(globalThis);
} else {
  /** @closureUnaware */
  (function() {
     window['foo'] = 10;
  }).call(globalThis);
}
`,
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (goog.DEBUG) {
  $jscomp_wrap_closure_unaware_code.call(globalThis)
} else {
  $jscomp_wrap_closure_unaware_code.call(globalThis)
}
`,
		[]string{
			`$jscomp_sink_closure_unaware_impl(
  function() {
    window['foo'] = 5;
  });
`,
			`$jscomp_sink_closure_unaware_impl(
  function() {
    window['foo'] = 10;
  });
`})
}

// Java: ManageClosureUnawareCodeTest.testConditionalAndDebugSrcLoad()
func (c *ManageClosureUnawareCodeTest) TestConditionalAndDebugSrcLoad() {
	c.DoTest(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (!window['foo']) {
  if (goog.DEBUG) {
    /** @closureUnaware */
    (function() {
      window['foo'] = 5;
    }).call(globalThis);
  } else {
    /** @closureUnaware */
    (function() {
       window['foo'] = 10;
    }).call(globalThis);
  }
}
`,
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (!window['foo']) {
  if (goog.DEBUG) {
    $jscomp_wrap_closure_unaware_code.call(globalThis)
  } else {
    $jscomp_wrap_closure_unaware_code.call(globalThis)
  }
}
`,
		[]string{
			`$jscomp_sink_closure_unaware_impl(
  function() {
    window['foo'] = 5;
  });
`,
			`$jscomp_sink_closure_unaware_impl(
  function() {
    window['foo'] = 10;
  });
`})
}

// Java: ManageClosureUnawareCodeTest.testDirectLoad_nestedChangeScopes()
func (c *ManageClosureUnawareCodeTest) TestDirectLoad_nestedChangeScopes() {
	c.DoTest(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (goog.DEBUG) {
  /** @closureUnaware */
  (function() {
    function bar() {
      window['foo'] = 5;
    }
    bar();
  }).call(globalThis);
}
`,
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (goog.DEBUG) {
  $jscomp_wrap_closure_unaware_code.call(globalThis)
}
`,
		[]string{
			`$jscomp_sink_closure_unaware_impl(
  function() {
    function bar() { window['foo'] = 5;} bar();
  });
`})
}

// Java: ManageClosureUnawareCodeTest.testDirectLoad_nestedGlobalThisIIFEIsNotRewritten()
func (c *ManageClosureUnawareCodeTest) TestDirectLoad_nestedGlobalThisIIFEIsNotRewritten() {
	c.DoTest(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (goog.DEBUG) {
  /** @closureUnaware */
  (function() {
    (function() {
      window['foo'] = 10;
    }).call(globalThis);
  }).call(globalThis);
}
`,
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (goog.DEBUG) {
  $jscomp_wrap_closure_unaware_code.call(globalThis)
}
`,
		[]string{
			`$jscomp_sink_closure_unaware_impl(
  function() {
   (function() {
      window['foo'] = 10;
    }).call(globalThis);
  });
`})
}

// Java: ManageClosureUnawareCodeTest.testErrorsOnWrapping_invalidAnnotation_statement()
func (c *ManageClosureUnawareCodeTest) TestErrorsOnWrapping_invalidAnnotation_statement() {
	c.Test(
		c.Srcs(
			`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
/** @closureUnaware */
exports.bar = 10;
`),
		c.Error(RhinoErrorReporter_PARSE_ERROR).
			WithMessageContaining("@closureUnaware root must be function"))
}

// Java: ManageClosureUnawareCodeTest.testErrorsForScriptsWithoutFileoverviewClosureUnawareAnnotation()
func (c *ManageClosureUnawareCodeTest) TestErrorsForScriptsWithoutFileoverviewClosureUnawareAnnotation() {
	c.Test(
		c.Srcs(
			`/**
 * @fileoverview
 */
goog.module('foo.bar.baz_raw');
/** @closureUnaware */
(function() {
  window['foo'] = 10;
}).call(globalThis);
`),
		c.Error(RhinoErrorReporter_PARSE_ERROR).
			WithMessageContaining("must be in a file annotated with @closureUnaware"))
}

// Java: ManageClosureUnawareCodeTest.testNoErrorIfMissingNodeLevelAnnotation()
func (c *ManageClosureUnawareCodeTest) TestNoErrorIfMissingNodeLevelAnnotation() {
	// It is an error to have a file-level @closureUnaware annotation, but no annotation on some
	// other node to indicate the start of the closure-unaware sub-AST.
	c.TestSame(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
$jscomp_wrap_closure_unaware_code(5)
`)
}

// Java: ManageClosureUnawareCodeTest.testNoErrorsForConfusingNestedAnnotations()
func (c *ManageClosureUnawareCodeTest) TestNoErrorsForConfusingNestedAnnotations() {
	c.DoTest(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
(/** @closureUnaware */ function() {
// This nested usage of the @closureUnaware annotation is not a parse error, but it
// could be confusing if a human author expected there to be multiple shadowed ASTs /
// sections of closure-aware and closure-unaware code.
// This test-case is just codifying that it is always the top-most AST node that is
// converted into a shadow AST node, and that there won't be nested shadow ASTs.
  (/** @closureUnaware */ function() {
    window['foo'] = 10;
  }).call(globalThis);
}).call(globalThis);
`,
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
$jscomp_wrap_closure_unaware_code.call(globalThis)
`,
		[]string{
			`$jscomp_sink_closure_unaware_impl(
  function() {
    // Note that there isn't a JSDoc closureUnaware annotation on the inner function, as
    // it was never created during parsing of the AST.
    (function() {
      window['foo'] = 10;
    }).call(globalThis);
  });
`})
}

// Java: ManageClosureUnawareCodeTest.testAllowsSpecifyingAnnotationOnIIFE()
func (c *ManageClosureUnawareCodeTest) TestAllowsSpecifyingAnnotationOnIIFE() {
	c.DoTest(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
/** @closureUnaware */
(function() {
  window['foo'] = 5;
}).call(globalThis);
`,
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
$jscomp_wrap_closure_unaware_code.call(globalThis)
`,
		[]string{
			`$jscomp_sink_closure_unaware_impl(
  function() {
    window['foo'] = 5;
  });
`})
}

// Java: ManageClosureUnawareCodeTest.testReparseWithInvalidJSDocTags()
func (c *ManageClosureUnawareCodeTest) TestReparseWithInvalidJSDocTags() {
	// @dependency is not a valid JSDoc tag, and for normal code we would issue a parser error.
	// However, for closure-unaware code we don't care at all what is in jsdoc comments.
	c.DoTest(
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
/** @closureUnaware */
(function() {
  /** @dependency */
  window['foo'] = 5;
}).call(globalThis);
`,
		`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
$jscomp_wrap_closure_unaware_code.call(globalThis)
`,
		[]string{
			`$jscomp_sink_closure_unaware_impl(
  function() {
    window['foo'] = 5;
  });
`})
}

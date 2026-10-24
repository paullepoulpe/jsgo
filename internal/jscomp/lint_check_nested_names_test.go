// Port of com/google/javascript/jscomp/lint/CheckNestedNamesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckNestedNamesTest.java)

/*
 * Copyright 2020 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/parsing"
)

// Test case for {@link CheckNestedNames}.
type CheckNestedNamesTest struct {
	CompilerTestCaseBase
}

func NewCheckNestedNamesTest() *CheckNestedNamesTest {
	c := &CheckNestedNamesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckNestedNamesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckNestedNamesTest() })
}

var checkNestedNamesTest_CLOSURE_EXTERNS = NewTestExternsBuilder().AddClosureExterns().Build()

// Java: CheckNestedNamesTest.setUp()
func (c *CheckNestedNamesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableCreateModuleMap()
	c.EnableTypeCheck()
}

// Java: CheckNestedNamesTest.getProcessor(Compiler)
func (c *CheckNestedNamesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckNestedNames(compiler)
}

// Java: CheckNestedNamesTest.getOptions()
func (c *CheckNestedNamesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE)
	options.SetWarningLevel(DiagnosticGroups_ANALYZER_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckNestedNamesTest.testNoWarningOnNonStaticPropertyAssignment()
func (c *CheckNestedNamesTest) TestNoWarningOnNonStaticPropertyAssignment() {
	c.NoWarning("goog.module('a'); class C {}; C.prototype.E = class {};")
	c.NoWarning("goog.module('a'); /** @interface */ C = class {}; C.prototype.E = class {};")
	c.DisableTypeCheck()
	// type checking reports inexistent .prototype access on enum, typedef
	c.NoWarning("goog.module('a'); /** @enum */ C = {}; C.prototype.E = class {};")
	c.NoWarning("goog.module('a'); /** @typedef {?} */ let T; T.prototype.E = class {};")
}

// Java: CheckNestedNamesTest.testNoWarningInScriptScope()
func (c *CheckNestedNamesTest) TestNoWarningInScriptScope() {
	c.NoWarning("class C {}; C.E = class {};")
	c.NoWarning("/** @enum */ let C = {}; C.E = class {};")
	c.NoWarning("/** @interface */ let C = function() {}; C.E = class {};")
	c.NoWarning("/** @typedef {?} */ const C = {}; C.E = class {};")
}

// Java: CheckNestedNamesTest.testNoWarningInFunctionScope()
func (c *CheckNestedNamesTest) TestNoWarningInFunctionScope() {
	c.NoWarning("goog.module('a'); function foo() { class C {}; /** @enum */ C.E = {};} foo();")
	c.NoWarning("goog.module('a'); function foo() { let C = {}; /** @enum */ C.E = {};} foo();")
	c.NoWarning(
		`goog.module('a');
/** @param {!Function} x */
function foo(x) {
  /** @enum */ x.E = {};
}
function y() {}
foo(y);
`)
}

// Java: CheckNestedNamesTest.testNestedClass()
func (c *CheckNestedNamesTest) TestNestedClass() {
	c.NestedNameWarning("goog.module('a'); class C {}; C.C = class {};")
	c.NestedNameWarning("goog.module('a'); let obj = {}; obj.C = class {};")
	c.NestedNameWarning("goog.module('a'); function F() {}; F.C = class {};")
	c.NestedNameWarning("goog.module('a'); let F = function() {}; F.E = class {};")
	c.NestedNameWarning("goog.module('a'); let F = function() {}; /** something */ F.E = class {};")
}

// Java: CheckNestedNamesTest.testNestedNames_noRHS()
func (c *CheckNestedNamesTest) TestNestedNames_noRHS() {
	c.NoWarning("goog.module('a'); class C {}; C.E;")
	c.DisableTypeCheck()
	// interfaces and enums that do not have rhs get reported in type checking; don't report here.
	c.NoWarning("goog.module('a'); class C {}; /** @interface */ C.E;")
	c.NoWarning("goog.module('a'); class C {}; /** @enum {string} */ C.E;")
	c.EnableTypeCheck()
	c.NestedNameWarning("goog.module('a'); class C {}; /** @typedef {{a:number}} */ C.T;")
}

// Java: CheckNestedNamesTest.testNestedEnum()
func (c *CheckNestedNamesTest) TestNestedEnum() {
	c.NestedNameWarning("goog.module('a'); let obj = {}; /** @enum */ obj.E = {};")
	c.NestedNameWarning("goog.module('a'); class C {}; /** @enum */ C.E = {};")
	c.NestedNameWarning("goog.module('a'); function F() {}; /** @enum */ F.E = {};")
	c.NestedNameWarning("goog.module('a'); let F = function() {}; /** @enum */ F.E = {};")
}

// Java: CheckNestedNamesTest.testNestedInterfaces()
func (c *CheckNestedNamesTest) TestNestedInterfaces() {
	c.NestedNameWarning("goog.module('a'); let obj = {}; /** @interface */ obj.I = class {};")
	c.NestedNameWarning("goog.module('a'); class C {}; /** @interface */ C.I = class {};")
	c.NestedNameWarning("goog.module('a'); function F() {}; /** @interface */ F.I = class {};")
	c.NestedNameWarning("goog.module('a'); let F = function() {}; /** @interface */ F.I = class {};")
}

// Java: CheckNestedNamesTest.testNoAnnotations()
func (c *CheckNestedNamesTest) TestNoAnnotations() {
	// tests `X.Y = ...` where `X.Y` is missing the @typedef/@interface/@enum annotation
	c.NestedNameWarning("goog.module('a'); class C {}; class D {}; /** @const */ C.C = D;")
	c.NestedNameWarning("goog.module('a'); class C {}; C.C = class {};")
	c.NestedNameWarning(
		`goog.module('a');
class C {};
/** @constructor */
function D() {};
/** @const */ C.C = D;
`)
	c.NoWarning("goog.module('a'); class C {}; function D() {}; /** @const */ C.C = D;")
	c.NestedNameWarning("goog.module('a'); class C {}; /** @enum */ E = {}; C.E = E;")
	c.NestedNameWarning("goog.module('a'); class C {}; /** @interface */ I = class {}; C.I = I;")
	// For the lhs to be a valid typedef, it must be annotated with `@const`.
	c.NoWarning("goog.module('a'); class C {}; /** @typedef {number} */ let T; C.T = T;")
	// The `@const` makes JSCompiler infer the lhs as a typedef, and hence we warn on it.
	c.NestedNameWarning(
		"goog.module('a'); class C {}; /** @typedef {number} */ let T; /** @const */ C.T = T;")
}

// Java: CheckNestedNamesTest.testNestingOnExternNames()
func (c *CheckNestedNamesTest) TestNestingOnExternNames() {
	c.TestWarning(
		c.Externs(NewTestExternsBuilder().AddClosureExterns().AddArray().Build()),
		c.Srcs("goog.module('a');  Array.C = class {};"),
		c.Warning(CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE))
}

// Java: CheckNestedNamesTest.testNestingOnGlobalNames()
func (c *CheckNestedNamesTest) TestNestingOnGlobalNames() {
	c.TestWarning(
		c.Externs(
			"let GlobalFn = function() {};"+NewTestExternsBuilder().AddClosureExterns().Build()),
		c.Srcs("goog.module('a');  GlobalFn.C = class {};"),
		c.Warning(CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE))
}

// Java: CheckNestedNamesTest.testErrorMessageText()
func (c *CheckNestedNamesTest) TestErrorMessageText() {
	c.TestWarning(
		c.Externs(NewTestExternsBuilder().AddClosureExterns().Build()),
		c.Srcs("goog.module('a'); class C {}; C.E = class {};"),
		c.Warning(CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE).WithMessageContaining(
			"A nested class is created on the name `C`."+
				" Fix this linter finding by converting the module-level static property"+
				" assignment on `C` into a module-level flat name (i.e. change `C.prop ="+
				" ...` into `C_prop = ...`. You can (if required) export this flat name"+
				" using named exports (`exports.C_prop = C_prop`).",
		))
}

// Java: CheckNestedNamesTest.testNoWarningOnNestedValues()
func (c *CheckNestedNamesTest) TestNoWarningOnNestedValues() {
	c.NoWarning("goog.module('a'); class C {}; C.E = '';")
	c.NoWarning("goog.module('a'); class SomeClass {} class C {}; C.E = new SomeClass();")
	c.NoWarning("goog.module('a'); /** @typedef {?} */ let C; C.E = {};")
	c.NoWarning("goog.module('a'); /** @interface {?} */ let C = class {}; C.E = {};")
	c.NoWarning("goog.module('a'); let C = {}; C.E = {};")
}

// Java: CheckNestedNamesTest.testNoWarningOnNamedExports()
func (c *CheckNestedNamesTest) TestNoWarningOnNamedExports() {
	c.NoWarning(
		"goog.module('a'); /** @enum {string} */ exports.SomeEnum = { ENUM_CONSTANT: 'value' };")
}

// Java: CheckNestedNamesTest.noWarning(String)
func (c *CheckNestedNamesTest) NoWarning(js string) {
	c.TestNoWarning(c.Externs(checkNestedNamesTest_CLOSURE_EXTERNS), c.Srcs(js))
}

// Java: CheckNestedNamesTest.nestedNameWarning(String)
func (c *CheckNestedNamesTest) NestedNameWarning(js string) {
	c.TestWarning(c.Externs(checkNestedNamesTest_CLOSURE_EXTERNS), c.Srcs(js), c.Warning(CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE))
}

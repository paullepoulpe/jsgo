// Port of com/google/javascript/jscomp/ProcessClosurePrimitivesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ProcessClosurePrimitivesTest.java)

/*
 * Copyright 2006 The Closure Compiler Authors.
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

// ProcessClosurePrimitivesTest holds tests for ProcessClosurePrimitives.
type ProcessClosurePrimitivesTest struct {
	CompilerTestCaseBase
}

func NewProcessClosurePrimitivesTest() *ProcessClosurePrimitivesTest {
	c := &ProcessClosurePrimitivesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestProcessClosurePrimitivesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewProcessClosurePrimitivesTest() })
}

// Java: ProcessClosurePrimitivesTest.getOptions()
func (c *ProcessClosurePrimitivesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	c.EnableCreateModuleMap()

	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// Java: ProcessClosurePrimitivesTest.getProcessor(Compiler)
func (c *ProcessClosurePrimitivesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewProcessClosurePrimitives(compiler)
}

const processClosurePrimitivesTest_METHOD_FORMAT = "function Foo() {} Foo.prototype.method = function() { %s };"

const processClosurePrimitivesTest_FOO_INHERITS = "goog.inherits(Foo, BaseFoo);"

// Java: ProcessClosurePrimitivesTest.testProvide_weakNamespaceError()
func (c *ProcessClosurePrimitivesTest) TestProvide_weakNamespaceError() {
	c.Test(
		c.Srcs("goog.provide('foo'); /** @type {Object} */ var foo = {};"),
		c.Warning(ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE))
	c.Test(
		c.Srcs("goog.provide('foo'); /** @type {!Object} */ var foo = {};"),
		c.Warning(ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE))
	c.Test(
		c.Srcs("goog.provide('foo.bar'); /** @type {Object} */ foo.bar = {};"),
		c.Warning(ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE))
	c.Test(
		c.Srcs("goog.provide('foo.bar'); /** @type {!Object} */ foo.bar = {};"),
		c.Warning(ProcessClosurePrimitives_WEAK_NAMESPACE_TYPE))

	c.TestNoWarning("goog.provide('foo'); /** @type {Object<string>} */ var foo = {};")

	c.TestNoWarning(c.Externs("/** @type {!Object} */ var foo = {};"), c.Srcs("goog.provide('foo');"))
}

// Java: ProcessClosurePrimitivesTest.testProvideValidObjectType()
func (c *ProcessClosurePrimitivesTest) TestProvideValidObjectType() {
	c.TestNoWarning(c.Srcs("goog.provide('foo'); /** @type {Object<string>} */ var foo = {};"))
}

// Java: ProcessClosurePrimitivesTest.testProvidedDeclaredFunctionError()
func (c *ProcessClosurePrimitivesTest) TestProvidedDeclaredFunctionError() {
	c.TestError("goog.provide('foo'); function foo(){}", ProcessClosurePrimitives_FUNCTION_NAMESPACE_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testProvidedDeclaredFunctionError_ignoredInExterns()
func (c *ProcessClosurePrimitivesTest) TestProvidedDeclaredFunctionError_ignoredInExterns() {
	c.TestNoWarning(c.Externs(" function foo(){}"), c.Srcs("goog.provide('foo');"))
}

// Java: ProcessClosurePrimitivesTest.testProvidedDeclaredClassError()
func (c *ProcessClosurePrimitivesTest) TestProvidedDeclaredClassError() {
	c.TestError("goog.provide('foo'); class foo {}", ProcessClosurePrimitives_CLASS_NAMESPACE_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testProvidedDeclaredClassError_ignoredInExterns()
func (c *ProcessClosurePrimitivesTest) TestProvidedDeclaredClassError_ignoredInExterns() {
	c.TestNoWarning(c.Externs("class foo {}"), c.Srcs("goog.provide('foo');"))
}

// Java: ProcessClosurePrimitivesTest.testAddDependency()
func (c *ProcessClosurePrimitivesTest) TestAddDependency() {
	c.Test("goog.addDependency('x.js', ['A', 'B'], []);", "0")

	compiler := c.GetLastCompiler()
	AssertThat(compiler.GetTypeRegistry().IsForwardDeclaredType("A")).IsFalse()
	AssertThat(compiler.GetTypeRegistry().IsForwardDeclaredType("B")).IsFalse()
	AssertThat(compiler.GetTypeRegistry().IsForwardDeclaredType("C")).IsFalse()
}

// Java: ProcessClosurePrimitivesTest.testSetCssNameMappingByShortHand()
func (c *ProcessClosurePrimitivesTest) TestSetCssNameMappingByShortHand() {
	c.TestError(
		"goog.setCssNameMapping({shortHandFirst, shortHandSecond});",
		ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testSetCssNameMappingByTemplate()
func (c *ProcessClosurePrimitivesTest) TestSetCssNameMappingByTemplate() {
	c.TestError(
		"goog.setCssNameMapping({foo: `bar`});", ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR)
	c.TestError(
		"goog.setCssNameMapping({foo: `${vari}bar`});",
		ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testSetCssNameMappingNonStringValueReturnsError()
func (c *ProcessClosurePrimitivesTest) TestSetCssNameMappingNonStringValueReturnsError() {
	// Make sure the argument is an object literal.
	c.TestError("var BAR = {foo:'bar'}; goog.setCssNameMapping(BAR);", ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR)
	c.TestError("goog.setCssNameMapping([]);", ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR)
	c.TestError("goog.setCssNameMapping(false);", ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR)
	c.TestError("goog.setCssNameMapping(null);", ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR)
	c.TestError("goog.setCssNameMapping(undefined);", ProcessClosurePrimitives_EXPECTED_OBJECTLIT_ERROR)

	// Make sure all values of the object literal are string literals.
	c.TestError(
		"var BAR = 'bar'; goog.setCssNameMapping({foo:BAR});",
		ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR)
	c.TestError("goog.setCssNameMapping({foo:6});", ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR)
	c.TestError(
		"goog.setCssNameMapping({foo:false});", ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR)
	c.TestError(
		"goog.setCssNameMapping({foo:null});", ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR)
	c.TestError(
		"goog.setCssNameMapping({foo:undefined});",
		ProcessClosurePrimitives_NON_STRING_PASSED_TO_SET_CSS_NAME_MAPPING_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testForwardDeclarations()
func (c *ProcessClosurePrimitivesTest) TestForwardDeclarations() {
	c.TestSame("goog.forwardDeclare('A.B')")

	compiler := c.GetLastCompiler()
	AssertThat(compiler.GetTypeRegistry().IsForwardDeclaredType("A.B")).IsTrue()
	AssertThat(compiler.GetTypeRegistry().IsForwardDeclaredType("C.D")).IsFalse()

	c.TestSame("goog.module('mod'); goog.forwardDeclare('A.B');")

	compiler = c.GetLastCompiler()
	AssertThat(compiler.GetTypeRegistry().IsForwardDeclaredType("A.B")).IsTrue()

	c.TestSame(
		c.Srcs("goog.provide('A.B');", "goog.module('mod'); const B = goog.forwardDeclare('A.B');"))

	compiler = c.GetLastCompiler()
	// This is a valid forward declaration, but for historical reasons, does not actually forward
	// declare the type 'A.B'.
	AssertThat(compiler.GetTypeRegistry().IsForwardDeclaredType("A.B")).IsFalse()
}

// Java: ProcessClosurePrimitivesTest.testInvalidForwardDeclarations()
func (c *ProcessClosurePrimitivesTest) TestInvalidForwardDeclarations() {
	c.TestError(
		"const B = goog.forwardDeclare('A.B');",
		ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR)
	c.TestError("goog.forwardDeclare();", ProcessClosurePrimitives_INVALID_FORWARD_DECLARE)

	c.TestError(
		"goog.forwardDeclare('A.B', 'C.D');", ProcessClosurePrimitives_INVALID_FORWARD_DECLARE)

	c.TestError("goog.forwardDeclare(`template`);", ProcessClosurePrimitives_INVALID_FORWARD_DECLARE)
	c.TestError(
		"goog.forwardDeclare(`${template}Sub`);", ProcessClosurePrimitives_INVALID_FORWARD_DECLARE)
}

// Java: ProcessClosurePrimitivesTest.testSetCssNameMappingValidity()
func (c *ProcessClosurePrimitivesTest) TestSetCssNameMappingValidity() {
	// Make sure that the keys don't have -'s
	c.TestWarning("goog.setCssNameMapping({'a': 'b', 'a-a': 'c'})", ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP)

	// In full mode, we check that map(a-b)=map(a)-map(b)
	c.TestWarning(
		"goog.setCssNameMapping({'a': 'b', 'a-a': 'c'}, 'BY_WHOLE')", ProcessClosurePrimitives_INVALID_CSS_RENAMING_MAP)

	// Unknown mapping type
	c.TestError("goog.setCssNameMapping({foo:'bar'}, 'UNKNOWN');", ProcessClosurePrimitives_INVALID_STYLE_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidAddDependency()
func (c *ProcessClosurePrimitivesTest) TestInvalidAddDependency() {
	c.TestError(
		"goog.provide('a.b'); var x = x || goog.addDependency('a.b', ['foo'], ['bar']);",
		ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR)
	c.TestError(
		"goog.provide('a.b'); x = goog.addDependency('a.b', ['foo'], ['bar']);",
		ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR)
	c.TestError(
		"goog.provide('a.b'); function f() { goog.addDependency('a.b', ['foo'], ['bar']); }",
		ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testValidGoogMethod()
func (c *ProcessClosurePrimitivesTest) TestValidGoogMethod() {
	c.TestSame("function f() { goog.isDef('a.b'); }")
	c.TestSame("function f() { goog.inherits(a, b); }")
	c.TestSame("function f() { goog.exportSymbol(a, b); }")
	c.TestSame("function f() { goog.setCssNameMapping({}); }")
	c.TestSame("x || goog.isDef('a.b');")
	c.TestSame("x || goog.inherits(a, b);")
	c.TestSame("x || goog.exportSymbol(a, b);")
	c.TestSame("x || void 0")
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase1()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase1() {
	c.TestError(
		"var Foo = function() {};"+processClosurePrimitivesTest_FOO_INHERITS+"Foo.base(this, 'method');", ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase2()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase2() {
	c.TestError(
		jutil.StringReplace(`function Foo() {}
FOO_INHERITS
Foo.method = function() {
  Foo.base(this, 'method');
};
`, "FOO_INHERITS", processClosurePrimitivesTest_FOO_INHERITS),
		ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase3()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase3() {
	c.TestError(jutil.StringFormat(processClosurePrimitivesTest_FOO_INHERITS+processClosurePrimitivesTest_METHOD_FORMAT, "Foo.base();"), ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase4()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase4() {
	c.TestError(
		jutil.StringFormat(processClosurePrimitivesTest_FOO_INHERITS+processClosurePrimitivesTest_METHOD_FORMAT, "Foo.base(this, 'bar');"), ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase5()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase5() {
	c.TestError(
		jutil.StringFormat(processClosurePrimitivesTest_FOO_INHERITS+processClosurePrimitivesTest_METHOD_FORMAT, "Foo.base('foo', 'method');"),
		ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase7()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase7() {
	c.TestError("function Foo() { Foo.base(this); };"+processClosurePrimitivesTest_FOO_INHERITS, ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase8()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase8() {
	c.TestError("var Foo = function() { Foo.base(this); };"+processClosurePrimitivesTest_FOO_INHERITS, ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase9()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase9() {
	c.TestError(
		"var goog = {}; goog.Foo = function() { goog.Foo.base(this); };"+processClosurePrimitivesTest_FOO_INHERITS,
		ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase10()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase10() {
	c.TestError("function Foo() { Foo.base(this); }"+processClosurePrimitivesTest_FOO_INHERITS, ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase11()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase11() {
	c.TestError("function Foo() { Foo.base(this, 'method'); }"+processClosurePrimitivesTest_FOO_INHERITS, ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase12()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase12() {
	c.TestError("function Foo() { Foo.base(this, 1, 2); }"+processClosurePrimitivesTest_FOO_INHERITS, ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase13()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase13() {
	c.TestError(
		`function Bar(){ Bar.base(this, 'constructor'); }
goog.inherits(Bar, Goo);
function Foo(){ Bar.base(this, 'constructor'); }
`+
			processClosurePrimitivesTest_FOO_INHERITS,
		ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidGoogBase14()
func (c *ProcessClosurePrimitivesTest) TestInvalidGoogBase14() {
	c.TestError("class Foo extends BaseFoo { constructor() { Foo.base(this); } }", ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testInvalidGoogBase14b()
func (c *ProcessClosurePrimitivesTest) TestInvalidGoogBase14b() {
	c.TestError(
		"class Foo extends BaseFoo { method() { Foo.base(this, 'method'); } }", ProcessClosurePrimitives_BASE_CLASS_ERROR)
}

// Java: ProcessClosurePrimitivesTest.testValidBase1()
func (c *ProcessClosurePrimitivesTest) TestValidBase1() {
	c.Test(
		processClosurePrimitivesTest_FOO_INHERITS+jutil.StringFormat(processClosurePrimitivesTest_METHOD_FORMAT, "Foo.base(this, 'method');"),
		processClosurePrimitivesTest_FOO_INHERITS+jutil.StringFormat(processClosurePrimitivesTest_METHOD_FORMAT, "Foo.superClass_.method.call(this)"))
}

// Java: ProcessClosurePrimitivesTest.testValidBase2()
func (c *ProcessClosurePrimitivesTest) TestValidBase2() {
	c.Test(
		processClosurePrimitivesTest_FOO_INHERITS+jutil.StringFormat(processClosurePrimitivesTest_METHOD_FORMAT, "Foo.base(this, 'method', 1, 2);"),
		processClosurePrimitivesTest_FOO_INHERITS+jutil.StringFormat(processClosurePrimitivesTest_METHOD_FORMAT, "Foo.superClass_.method.call(this, 1, 2)"))
}

// Java: ProcessClosurePrimitivesTest.testValidBase3()
func (c *ProcessClosurePrimitivesTest) TestValidBase3() {
	c.Test(
		processClosurePrimitivesTest_FOO_INHERITS+jutil.StringFormat(processClosurePrimitivesTest_METHOD_FORMAT, "return Foo.base(this, 'method');"),
		processClosurePrimitivesTest_FOO_INHERITS+jutil.StringFormat(processClosurePrimitivesTest_METHOD_FORMAT, "return Foo.superClass_.method.call(this)"))
}

// Java: ProcessClosurePrimitivesTest.testValidBase4()
func (c *ProcessClosurePrimitivesTest) TestValidBase4() {
	c.Test(
		"function Foo() { Foo.base(this, 'constructor', 1, 2); }"+processClosurePrimitivesTest_FOO_INHERITS,
		"function Foo() { BaseFoo.call(this, 1, 2); } "+processClosurePrimitivesTest_FOO_INHERITS)
}

// Java: ProcessClosurePrimitivesTest.testValidBase5()
func (c *ProcessClosurePrimitivesTest) TestValidBase5() {
	c.Test(
		"var Foo = function() { Foo.base(this, 'constructor', 1); };"+processClosurePrimitivesTest_FOO_INHERITS,
		"var Foo = function() { BaseFoo.call(this, 1); }; "+processClosurePrimitivesTest_FOO_INHERITS)
}

// Java: ProcessClosurePrimitivesTest.testValidBase6()
func (c *ProcessClosurePrimitivesTest) TestValidBase6() {
	c.Test(
		`var goog = {}; goog.Foo = function() {
goog.Foo.base(this, 'constructor'); };
goog.inherits(goog.Foo, goog.BaseFoo);
`,
		`var goog = {}; goog.Foo = function() { goog.BaseFoo.call(this); };
goog.inherits(goog.Foo, goog.BaseFoo);
`)
}

// Java: ProcessClosurePrimitivesTest.testValidBase7()
func (c *ProcessClosurePrimitivesTest) TestValidBase7() {
	// No goog.inherits, so this is probably a different 'base' function.
	c.TestSame(
		`var a = function() {
  a.base(this, 'constructor');
};
`)
}

// Java: ProcessClosurePrimitivesTest.testValidBase_exportsAssignmentsBeforeGoogInherits()
func (c *ProcessClosurePrimitivesTest) TestValidBase_exportsAssignmentsBeforeGoogInherits() {
	c.Test(
		`goog.module('my.Foo');
class Bar {}
function Foo() { Foo.base(this, 'constructor', 1, 2); }
exports.Foo = Foo;
exports.Bar = Bar;
`+
			processClosurePrimitivesTest_FOO_INHERITS,
		`goog.module('my.Foo');
class Bar {}
function Foo() { BaseFoo.call(this, 1, 2); }
exports.Foo = Foo;
exports.Bar = Bar;
`+
			processClosurePrimitivesTest_FOO_INHERITS)
}

// Java: ProcessClosurePrimitivesTest.testInvalidBase_nonAliasLinesBeforeGoogInherits()
func (c *ProcessClosurePrimitivesTest) TestInvalidBase_nonAliasLinesBeforeGoogInherits() {
	c.TestSame(
		`goog.module('my.Foo');
function Foo() { Foo.base(this, 'constructor', 1, 2); }
alert(0);
alert(1);
alert(2);
` +
			processClosurePrimitivesTest_FOO_INHERITS)
}

// Java: ProcessClosurePrimitivesTest.testValidBase_googProvide_googRequiredInOtherModule()
func (c *ProcessClosurePrimitivesTest) TestValidBase_googProvide_googRequiredInOtherModule() {
	c.Test(
		c.Srcs(
			`goog.provide('my.Foo');
/** @constructor */
my.Foo = function() {}
goog.inherits(my.Foo, BaseFoo);
`,
			`goog.module('test');
const Foo = goog.require('my.Foo');
Foo.prototype.method = function() {
  Foo.base(this, 'method');
};
`),
		c.Error(ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR))
}

// Java: ProcessClosurePrimitivesTest.testValidBase_googModule_googRequiredInOtherModule()
func (c *ProcessClosurePrimitivesTest) TestValidBase_googModule_googRequiredInOtherModule() {
	c.Test(
		c.Srcs(
			`goog.module('FooModule');
/** @constructor */ function Foo() {}
goog.inherits(Foo, BaseFoo);
exports = {Foo};
`,
			`goog.module('test');
const {Foo: FooRequired} = goog.require('FooModule');
FooRequired.prototype.method = function() {
  FooRequired.base(this, 'method');
};
`),
		c.Error(ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR))
}

// Java: ProcessClosurePrimitivesTest.testValidPrimitiveCalls()
func (c *ProcessClosurePrimitivesTest) TestValidPrimitiveCalls() {
	c.TestNoWarning(
		`goog.module('c');
goog.forwardDeclare('A.b');
`)
	c.TestNoWarning(
		`goog.module('d');
goog.addDependency('C.D');
`)
}

// Java: ProcessClosurePrimitivesTest.testInvalidGoogWeakUsage()
func (c *ProcessClosurePrimitivesTest) TestInvalidGoogWeakUsage() {
	// Error: the argument must be a name.
	c.TestError("goog.weakUsage(1);", ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL)

	// Error: there must be exactly one argument.
	c.TestError("a = 1; b = 1; goog.weakUsage(a, b);", ProcessClosurePrimitives_INVALID_GOOG_WEAK_USAGE_CALL)

	// Regular name is allowed.
	c.TestSame("a = 1; goog.weakUsage(a);")

	// Qualified name is allowed.
	c.TestSame("a = {b:1}; goog.weakUsage(a.b);")
}

// Java: ProcessClosurePrimitivesTest.testOtherBaseCall()
func (c *ProcessClosurePrimitivesTest) TestOtherBaseCall() {
	c.TestSame("class Foo extends BaseFoo { method() { baz.base('arg'); } }")
}

// Java: ProcessClosurePrimitivesTest.testRenameFunction_withOneStringLit_isOk()
func (c *ProcessClosurePrimitivesTest) TestRenameFunction_withOneStringLit_isOk() {
	c.TestSame("const p = JSCompiler_renameProperty('a')")
}

// Java: ProcessClosurePrimitivesTest.testRenameFunction_withOneStringLit_andAnotherArg_isOk()
func (c *ProcessClosurePrimitivesTest) TestRenameFunction_withOneStringLit_andAnotherArg_isOk() {
	c.TestSame("const p = JSCompiler_renameProperty('a', 0)")
}

// Java: ProcessClosurePrimitivesTest.testRenameFunction_withZeroArgs_isReported()
func (c *ProcessClosurePrimitivesTest) TestRenameFunction_withZeroArgs_isReported() {
	c.Test(
		c.Srcs("const p = JSCompiler_renameProperty()"),
		c.Error(ProcessClosurePrimitives_INVALID_RENAME_FUNCTION).WithMessageContaining("1 or 2 arguments"),
		c.Error(ProcessClosurePrimitives_INVALID_RENAME_FUNCTION).WithMessageContaining("string literal"))
}

// Java: ProcessClosurePrimitivesTest.testRenameFunction_withThreeArgs_isReported()
func (c *ProcessClosurePrimitivesTest) TestRenameFunction_withThreeArgs_isReported() {
	c.Test(
		c.Srcs("const p = JSCompiler_renameProperty(1, 2, 3)"),
		c.Error(ProcessClosurePrimitives_INVALID_RENAME_FUNCTION).WithMessageContaining("1 or 2 arguments"),
		c.Error(ProcessClosurePrimitives_INVALID_RENAME_FUNCTION).WithMessageContaining("string literal"))
}

// Java: ProcessClosurePrimitivesTest.testRenameFunction_withNonStringArg_isReported()
func (c *ProcessClosurePrimitivesTest) TestRenameFunction_withNonStringArg_isReported() {
	c.Test(
		c.Srcs("const p = JSCompiler_renameProperty(0)"),
		c.Error(ProcessClosurePrimitives_INVALID_RENAME_FUNCTION).WithMessageContaining("string literal"))
}

// Java: ProcessClosurePrimitivesTest.testInvalidRenameFunction_withPropertyRefInFirstArg_isReported()
func (c *ProcessClosurePrimitivesTest) TestInvalidRenameFunction_withPropertyRefInFirstArg_isReported() {
	c.Test(
		c.Srcs("const p = JSCompiler_renameProperty('a.b')"),
		c.Error(ProcessClosurePrimitives_INVALID_RENAME_FUNCTION).WithMessageContaining("property path"))
}

// Java: ProcessClosurePrimitivesTest.testPreprocessorSymbolTable()
func (c *ProcessClosurePrimitivesTest) TestPreprocessorSymbolTable() {
	root := rhino.IR_Root()
	table := NewPreprocessorSymbolTable(root)
	n := rhino.IR_Name("test")
	table.AddReferenceWithName(n, "test")
	AssertThat(table.GetReferences(table.GetSlot("test").(rhino.SimpleSlot))).IsNotNull()
}

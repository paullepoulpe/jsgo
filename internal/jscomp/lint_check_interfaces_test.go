// Port of com/google/javascript/jscomp/lint/CheckInterfacesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/lint/CheckInterfacesTest.java)

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
)

// Test case for {@link CheckInterfaces}.
type CheckInterfacesTest struct {
	CompilerTestCaseBase
}

func NewCheckInterfacesTest() *CheckInterfacesTest {
	c := &CheckInterfacesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckInterfacesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckInterfacesTest() })
}

// Java: CheckInterfacesTest.getProcessor(Compiler)
func (c *CheckInterfacesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckInterfaces(compiler)
}

// Java: CheckInterfacesTest.getOptions()
func (c *CheckInterfacesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckInterfacesTest.testInterfaceArgs()
func (c *CheckInterfacesTest) TestInterfaceArgs() {
	c.TestSame(
		"/** @interface */ function A(x) {}",
		CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS)
	c.TestSame(
		"/** @interface */ class C { constructor(x) {} }",
		CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS)
	c.TestSame(
		`var ns = {};

/** @interface */

ns.SomeInterface = function(x) {};
`,
		CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS)
	c.TestSame(
		`var ns = {};

/** @interface */

ns.SomeInterface = class { constructor(x) {}};
`,
		CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS)
}

// Java: CheckInterfacesTest.testInterfaceArgs_withES6Modules()
func (c *CheckInterfacesTest) TestInterfaceArgs_withES6Modules() {
	c.TestSame(
		"export /** @interface */ function A(x) {}",
		CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS)
	c.TestSame(
		"export /** @interface */ class C { constructor(x) {} }",
		CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS)
}

// Java: CheckInterfacesTest.testInterfaceConstructorWithFieldDeclarations_noWarning()
func (c *CheckInterfacesTest) TestInterfaceConstructorWithFieldDeclarations_noWarning() {
	c.TestSame("/** @interface */ function A() { /** @type {string} */ this.foo; }")
	c.TestSame("/** @interface */ class C { constructor() { /** @type {string} */ this.foo; }}")
	c.TestSame(
		`var ns = {};
/** @interface */
ns.SomeInterface = function() { /** @type {string} */ this.foo; };
`)
	c.TestSame(
		`var ns = {};
/** @interface */
ns.SomeInterface = class { constructor() { /** @type {string} */ this.foo; }; }
`)
}

// Java: CheckInterfacesTest.testInterfaceMethodNotEmpty()
func (c *CheckInterfacesTest) TestInterfaceMethodNotEmpty() {
	c.TestSame(
		`/** @interface */
class C {
  constructor() {}
  A() { this.foo; }
}
`,
		CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY)
}

// Java: CheckInterfacesTest.testInterfaceWithFieldDeclarations_withES6Modules()
func (c *CheckInterfacesTest) TestInterfaceWithFieldDeclarations_withES6Modules() {
	c.TestSame("export /** @interface */ function A() { /** @type {string} */ this.foo; }")
	c.TestSame(
		"export /** @interface */ class C { constructor() { /** @type {string} */ this.foo; } }")
}

// Java: CheckInterfacesTest.testInterfaceClass_callToSuperInConstructorNotAllowed()
func (c *CheckInterfacesTest) TestInterfaceClass_callToSuperInConstructorNotAllowed() {
	c.DisableAstValidation()
	// silences ASTValidator for `super()` used without `extends`
	c.TestWarning(
		`class D {} /** @interface 
 @extends {D} */ class C { constructor() { super(); } }`,
		CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE)
	c.TestWarning(
		`class D {} /** @interface 
 @extends {D} */ class C { constructor() { super(x); } }`,
		CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE)
	c.TestWarning(
		`class D {} /** @record 
 @extends {D} */ class C { constructor() { super(); } }`,
		CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE)
	c.TestWarning(
		`class D {} /** @record 
 @extends {D} */ class C { constructor() { super(x); } }`,
		CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE)
}

// Java: CheckInterfacesTest.testInterfaceClass_extendsKeywordNotAllowed()
func (c *CheckInterfacesTest) TestInterfaceClass_extendsKeywordNotAllowed() {
	c.Test(
		c.Srcs("class D {} /** @interface */ class C extends D { constructor() { super(); } }"),
		c.Warning(CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS),
		c.Warning(CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE))
	c.Test(
		c.Srcs("class D {} /** @interface */ class C extends D { constructor() { super(x); } }"),
		c.Warning(CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS),
		c.Warning(CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE))
	c.Test(
		c.Srcs("class D {} /** @record */ class C extends D { constructor() { super(); } }"),
		c.Warning(CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS),
		c.Warning(CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE))
	c.Test(
		c.Srcs("class D {} /** @record */ class C extends D { constructor() { super(x); } }"),
		c.Warning(CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS),
		c.Warning(CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE))
}

// Java: CheckInterfacesTest.testInterfaceComputedProperties()
func (c *CheckInterfacesTest) TestInterfaceComputedProperties() {
	c.TestSame(
		"/** @interface */  class C { ['f']() { return 1; }}",
		CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY)
}

// Java: CheckInterfacesTest.testInterfaceGetters()
func (c *CheckInterfacesTest) TestInterfaceGetters() {
	c.TestSame(
		"/** @interface */  class C { get One() { return 1; }}",
		CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY)
}

// Java: CheckInterfacesTest.testInterfaceSetters()
func (c *CheckInterfacesTest) TestInterfaceSetters() {
	c.TestSame(
		`/** @interface */
class C {
  set One(x) {
   this.one = x;
  }
}
`,
		CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY)
}

// Java: CheckInterfacesTest.testInterfaceEs6ClassDeclaration_havingStaticMethod()
func (c *CheckInterfacesTest) TestInterfaceEs6ClassDeclaration_havingStaticMethod() {
	c.TestSame(
		c.Srcs(
			`/** @interface */
class I {
  constructor() {}
  static foo() {}
}
I.foo();
`),
		c.Warning(CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS).WithMessageContaining(
			"Consider pulling out the static method into a flat name as I_foo"))
}

// Java: CheckInterfacesTest.testInterfaceEs6ClassAssignment_havingStaticMethod()
func (c *CheckInterfacesTest) TestInterfaceEs6ClassAssignment_havingStaticMethod() {
	c.TestSame(
		c.Srcs(
			`/** @interface */
let I = class {
  constructor() {}
  static foo() {}
}
I.foo();
`),
		c.Warning(CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS).WithMessageContaining(
			"Consider pulling out the static method into a flat name as I_foo"))
}

// Java: CheckInterfacesTest.testInterfaceEs6ClassAssignment_havingStaticMethod2()
func (c *CheckInterfacesTest) TestInterfaceEs6ClassAssignment_havingStaticMethod2() {
	c.TestSame(
		c.Srcs(
			`let C;
/** @interface */
C.I = class {
  constructor() {}
  static foo() {}
}
C.I.foo();
`),
		c.Warning(CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS).WithMessageContaining(
			"Consider pulling out the static method into a flat name as C.I_foo"))
}

// Java: CheckInterfacesTest.testInterfaceWithFieldDeclarations_missingJDoc()
func (c *CheckInterfacesTest) TestInterfaceWithFieldDeclarations_missingJDoc() {
	c.TestSame(
		"/** @interface */ function A() { this.foo; }",
		CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT)
	c.TestSame(
		"/** @interface */ class C { constructor() { this.foo; }}",
		CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT)
	// in ES6 modules
	c.TestSame(
		"export /** @interface */ function A() {this.foo; }",
		CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT)
}

// Java: CheckInterfacesTest.testRecordWithFieldDeclarations()
func (c *CheckInterfacesTest) TestRecordWithFieldDeclarations() {
	c.TestSame(
		`/** @record */
function R() {
  /** @type {string} */
  this.foo;

  /** @type {number} */
  this.bar;
}
`)
}

// Java: CheckInterfacesTest.testRecordWithFieldDeclarationsMissingJSDoc()
func (c *CheckInterfacesTest) TestRecordWithFieldDeclarationsMissingJSDoc() {
	c.TestSame(
		`/** @record */
function R() {
  // This should have a JSDoc.
  this.noJSDoc;

}
`,
		CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT)
}

// Java: CheckInterfacesTest.testRecordWithOtherContents()
func (c *CheckInterfacesTest) TestRecordWithOtherContents() {
	c.TestSame(
		`/** @record */
function R() {
  /** @type {string} */
  let foo = '';
}
`,
		CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE)
}

// Java: CheckInterfacesTest.testInterfaceWithOtherContents()
func (c *CheckInterfacesTest) TestInterfaceWithOtherContents() {
	c.TestSame(
		`/** @interface */
function R() {
  /** @type {string} */
  let foo = '';
}
`,
		CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE)
}

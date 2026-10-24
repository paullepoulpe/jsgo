// Port of com/google/javascript/jscomp/SymbolTableTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/SymbolTableTest.java)

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"slices"
	"strconv"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// SymbolTableTest is an integration test for symbol table creation.
//
// Java: public final class SymbolTableTest
type SymbolTableTest struct {
	t *testing.T

	options *CompilerOptions
}

// Java: SymbolTableTest()
func NewSymbolTableTest() *SymbolTableTest {
	return &SymbolTableTest{}
}

func TestSymbolTableTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewSymbolTableTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *SymbolTableTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore runs the @Before method setUp.
func (c *SymbolTableTest) JUnitBefore() { c.SetUp() }

// JUnitAfter: the class has no @After method.
func (c *SymbolTableTest) JUnitAfter() {}

// Java: SymbolTableTest.setUp()
func (c *SymbolTableTest) SetUp() {

	c.options = NewCompilerOptions()
	c.options.SetCodingConvention(NewClosureCodingConvention())
	CompilationLevel_SIMPLE_OPTIMIZATIONS.SetOptionsForCompilationLevel(c.options)
	WarningLevel_VERBOSE.SetOptionsForWarningLevel(c.options)
	c.options.SetChecksOnly(true)
	c.options.SetPreserveDetailedSourceInfo(true)
	c.options.SetPreserveClosurePrimitives(true)
	c.options.SetContinueAfterErrors(true)
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE)
	c.options.SetEnableModuleRewriting(false)
}

// Java: SymbolTableTest.testFunctionInCall()
func (c *SymbolTableTest) TestFunctionInCall() {
	c.CreateSymbolTable(
		`class Y { constructor(fn) {} }
class X extends Y {
  constructor() {
    super(function() {});
  }
}
`)
}

// Java: SymbolTableTest.testGlobalVar()
func (c *SymbolTableTest) TestGlobalVar() {
	table := c.CreateSymbolTable("/** @type {number} */ var x = 5;")
	AssertThat(c.GetGlobalVar(table, "y")).IsNull()
	AssertThat(c.GetGlobalVar(table, "x")).IsNotNull()
	AssertThat(c.GetGlobalVar(table, "x").GetType().String()).IsEqualTo("number")

	// 2 == sizeof({x, *global*})
	AssertThat(c.GetVars(table)).HasSize(2)
}

// Java: SymbolTableTest.testGlobalThisReferences()
func (c *SymbolTableTest) TestGlobalThisReferences() {
	table :=
		c.CreateSymbolTable("var x = this; function f() { return this + this + this; }")

	global := c.GetGlobalVar(table, "*global*")
	AssertThat(global).IsNotNull()

	refs := table.GetReferenceList(global)
	AssertThat(refs).HasSize(1)
}

// Java: SymbolTableTest.testGlobalThisReferences2()
func (c *SymbolTableTest) TestGlobalThisReferences2() {
	// Make sure the global this is declared, even if it isn't referenced.
	table := c.CreateSymbolTable("")

	global := c.GetGlobalVar(table, "*global*")
	AssertThat(global).IsNotNull()

	refs := table.GetReferenceList(global)
	AssertThat(refs).IsEmpty()
}

// Java: SymbolTableTest.testGlobalThisReferences3()
func (c *SymbolTableTest) TestGlobalThisReferences3() {
	table := c.CreateSymbolTable("this.foo = {}; this.foo.bar = {};")

	global := c.GetGlobalVar(table, "*global*")
	AssertThat(global).IsNotNull()

	refs := table.GetReferenceList(global)
	AssertThat(refs).HasSize(2)
}

// Java: SymbolTableTest.testSourceInfoForProvidedSymbol()
func (c *SymbolTableTest) TestSourceInfoForProvidedSymbol() {
	table := c.CreateSymbolTable("goog.provide('foo.bar.Baz'); foo.bar.Baz = class {};")

	foo := c.GetGlobalVar(table, "foo")
	AssertThat(foo).IsNotNull()
	// goog.provide doesn't define foo. Instead foo is declared by the firts occurence in actual
	// code.
	AssertNode(table.GetReferenceList(foo)[0].GetNode()).HasCharno(29)

	fooBar := c.GetGlobalVar(table, "foo.bar")
	AssertThat(fooBar).IsNotNull()
	AssertNode(table.GetReferenceList(fooBar)[0].GetNode()).HasCharno(33)

	fooBarBaz := c.GetGlobalVar(table, "foo.bar.Baz")
	AssertThat(fooBarBaz).IsNotNull()
	AssertNode(table.GetReferenceList(fooBarBaz)[0].GetNode()).HasCharno(37)
}

// Java: SymbolTableTest.testGlobalThisPropertyReferences()
func (c *SymbolTableTest) TestGlobalThisPropertyReferences() {
	table := c.CreateSymbolTable("/** @constructor */ function Foo() {} this.Foo;")

	foo := c.GetGlobalVar(table, "Foo")
	AssertThat(foo).IsNotNull()

	refs := table.GetReferenceList(foo)
	AssertThat(refs).HasSize(2)
}

// Java: SymbolTableTest.testGlobalVarReferences()
func (c *SymbolTableTest) TestGlobalVarReferences() {
	table := c.CreateSymbolTable("/** @type {number} */ var x = 5; x = 6;")
	x := c.GetGlobalVar(table, "x")
	refs := table.GetReferenceList(x)

	AssertThat(refs).HasSize(2)
	AssertThat(refs[0]).IsEqualTo(x.GetDeclaration())
	AssertThat(refs[0].GetNode().GetParent().GetToken()).IsEqualTo(rhino.Token_VAR)
	AssertThat(refs[1].GetNode().GetParent().GetToken()).IsEqualTo(rhino.Token_ASSIGN)
}

// Java: SymbolTableTest.testLocalVarReferences()
func (c *SymbolTableTest) TestLocalVarReferences() {
	table := c.CreateSymbolTable("function f(x) { return x; }")
	x := c.GetLocalVar(table, "x")
	refs := table.GetReferenceList(x)

	AssertThat(refs).HasSize(2)
	AssertThat(refs[0]).IsEqualTo(x.GetDeclaration())
	AssertThat(refs[0].GetNode().GetParent().GetToken()).IsEqualTo(rhino.Token_PARAM_LIST)
	AssertThat(refs[1].GetNode().GetParent().GetToken()).IsEqualTo(rhino.Token_RETURN)
}

// Java: SymbolTableTest.testLocalThisReferences()
func (c *SymbolTableTest) TestLocalThisReferences() {
	table :=
		c.CreateSymbolTable("/** @constructor */ function F() { this.foo = 3; this.bar = 5; }")

	f := c.GetGlobalVar(table, "F")
	AssertThat(f).IsNotNull()

	t := table.GetParameterInFunction(f, "this")
	AssertThat(t).IsNotNull()

	refs := table.GetReferenceList(t)
	AssertThat(refs).HasSize(2)
}

// Java: SymbolTableTest.testLocalThisReferences3()
func (c *SymbolTableTest) TestLocalThisReferences3() {
	table := c.CreateSymbolTable("/** @constructor */ function F() {}")

	baz := c.GetGlobalVar(table, "F")
	AssertThat(baz).IsNotNull()

	t := table.GetParameterInFunction(baz, "this")
	AssertThat(t).IsNull()
}

// Java: SymbolTableTest.testObjectLiteral()
func (c *SymbolTableTest) TestObjectLiteral() {
	table := c.CreateSymbolTable("var obj = {foo: 0};")

	foo := c.GetGlobalVar(table, "obj.foo")
	AssertThat(foo.GetName()).IsEqualTo("foo")
}

// Java: SymbolTableTest.testObjectLiteralQuoted()
func (c *SymbolTableTest) TestObjectLiteralQuoted() {
	table := c.CreateSymbolTable("var obj = {'foo': 0};")

	// Quoted keys are not symbols.
	AssertThat(c.GetGlobalVar(table, "obj.foo")).IsNull()
	// Only obj and *global* are symbols.
	AssertThat(c.GetVars(table)).HasSize(2)
}

// Java: SymbolTableTest.testObjectLiteralWithMemberFunction()
func (c *SymbolTableTest) TestObjectLiteralWithMemberFunction() {
	input := "var obj = { fn() {} }; obj.fn();"
	table := c.CreateSymbolTable(input)

	objFn := c.GetGlobalVar(table, "obj.fn")
	AssertThat(objFn).IsNotNull()
	references := table.GetReferenceList(objFn)
	AssertThat(references).HasSize(2)

	// The declaration node corresponds to "fn", not "fn() {}", in the source info.
	declaration := objFn.GetDeclarationNode()
	AssertThat(declaration.GetCharno()).IsEqualTo(12)
	AssertThat(declaration.GetLength()).IsEqualTo(2)
}

// Java: SymbolTableTest.testNamespacedReferences()
func (c *SymbolTableTest) TestNamespacedReferences() {
	// Because the type of 'my' is anonymous, we build its properties into
	// the global scope.
	table :=
		c.CreateSymbolTable(
			`var my = {};
my.dom = {};
my.dom.DomHelper = function(){};
`)

	my := c.GetGlobalVar(table, "my")
	AssertThat(my).IsNotNull()
	AssertThat(table.GetReferences(my)).HasSize(3)

	myDom := c.GetGlobalVar(table, "my.dom")
	AssertThat(myDom).IsNotNull()
	AssertThat(table.GetReferences(myDom)).HasSize(2)

	myDomHelper := c.GetGlobalVar(table, "my.dom.DomHelper")
	AssertThat(myDomHelper).IsNotNull()
	AssertThat(table.GetReferences(myDomHelper)).HasSize(1)
}

// Java: SymbolTableTest.testIncompleteNamespacedReferences()
func (c *SymbolTableTest) TestIncompleteNamespacedReferences() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */
my.dom.DomHelper = function(){};
var y = my.dom.DomHelper;
`)
	my := c.GetGlobalVar(table, "my")
	AssertThat(my).IsNotNull()
	AssertThat(table.GetReferenceList(my)).HasSize(2)

	myDom := c.GetGlobalVar(table, "my.dom")
	AssertThat(myDom).IsNotNull()
	AssertThat(table.GetReferenceList(myDom)).HasSize(2)

	myDomHelper := c.GetGlobalVar(table, "my.dom.DomHelper")
	AssertThat(myDomHelper).IsNotNull()
	AssertThat(table.GetReferences(myDomHelper)).HasSize(2)
}

// Java: SymbolTableTest.testGlobalRichObjectReference()
func (c *SymbolTableTest) TestGlobalRichObjectReference() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */
function A(){};
/** @type {?A} */ A.prototype.b;
/** @type {A} */ var a = new A();
function g() {
  return a.b ? 'x' : 'y';
}
(function() {
  var x; if (x) { x = a.b.b; } else { x = a.b.c; }
  return x;
})();
`)

	ab := c.GetGlobalVar(table, "a.b")
	AssertThat(ab).IsNull()

	propB := c.GetGlobalVar(table, "A.prototype.b")
	AssertThat(propB).IsNotNull()
	AssertThat(table.GetReferenceList(propB)).HasSize(5)
}

// Java: SymbolTableTest.testRemovalOfNamespacedReferencesOfProperties()
func (c *SymbolTableTest) TestRemovalOfNamespacedReferencesOfProperties() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ var DomHelper = function(){};
/** method */ DomHelper.method = function() {};
`)

	domHelper := c.GetGlobalVar(table, "DomHelper")
	AssertThat(domHelper).IsNotNull()

	domHelperNamespacedMethod := c.GetGlobalVar(table, "DomHelper.method")
	AssertThat(domHelperNamespacedMethod.GetName()).IsEqualTo("method")

	domHelperMethod := domHelper.GetPropertyScope().GetSlot("method")
	AssertThat(domHelperMethod).IsNotNull()
}

// Java: SymbolTableTest.testGoogScopeReferences()
func (c *SymbolTableTest) TestGoogScopeReferences() {
	table :=
		c.CreateSymbolTableWithDefaultExterns(
			// goog.scope is defined in the default externs, among other Closure methods
			"goog.scope(function() {});")

	googScope := c.GetGlobalVar(table, "goog.scope")
	AssertThat(googScope).IsNotNull()
	AssertThat(table.GetReferences(googScope)).HasSize(2)
}

// Java: SymbolTableTest.testGoogRequireReferences()
func (c *SymbolTableTest) TestGoogRequireReferences() {
	table :=
		c.CreateSymbolTableWithDefaultExterns(
			`goog.provide('goog.dom');
goog.require('goog.dom');
`)
	googRequire := c.GetGlobalVar(table, "goog.require")
	AssertThat(googRequire).IsNotNull()

	AssertThat(table.GetReferences(googRequire)).HasSize(2)
}

// Java: SymbolTableTest.testNamespaceReferencesInGoogRequire()
func (c *SymbolTableTest) TestNamespaceReferencesInGoogRequire() {
	table :=
		c.CreateSymbolTable(
			`goog.provide('my.dom');
goog.require('my.dom');
`)
	googRequire := c.GetGlobalVar(table, "ns$my.dom")
	AssertThat(googRequire).IsNotNull()

	AssertThat(table.GetReferences(googRequire)).HasSize(2)
}

// Java: SymbolTableTest.exerciseGoogRequireReferenceCorrectly(String)
func (c *SymbolTableTest) ExerciseGoogRequireReferenceCorrectly(inputProvidingNamespace string) {
	table :=
		c.CreateSymbolTableFromManySources(
			inputProvidingNamespace,
			`goog.provide('module.two');
goog.require('module.one');
goog.requireType('module.one');
goog.forwardDeclare('module.one');
`,
			`goog.module('module.three');
const one = goog.require('module.one');
const oneType = goog.requireType('module.one');
const oneForward = goog.forwardDeclare('module.one');
`,
			`goog.module('module.four');
goog.module.declareLegacyNamespace();
const one = goog.require('module.one');
const oneType = goog.requireType('module.one');
const oneForward = goog.forwardDeclare('module.one');
`)
	namespace := c.GetGlobalVar(table, "ns$module.one")
	AssertThat(namespace).IsNotNull()
	// 1 ref from declaration and then 2 refs in each of 3 files.
	AssertThat(table.GetReferences(namespace)).HasSize(10)
}

// Java: SymbolTableTest.testGoogProvideReferenced()
func (c *SymbolTableTest) TestGoogProvideReferenced() {
	c.ExerciseGoogRequireReferenceCorrectly("goog.provide('module.one');")
}

// Java: SymbolTableTest.testGoogModuleReferenced()
func (c *SymbolTableTest) TestGoogModuleReferenced() {
	c.ExerciseGoogRequireReferenceCorrectly("goog.module('module.one');")
}

// Java: SymbolTableTest.testGoogLegacyModuleReferenced()
func (c *SymbolTableTest) TestGoogLegacyModuleReferenced() {
	c.ExerciseGoogRequireReferenceCorrectly(
		`goog.module('module.one');
goog.module.declareLegacyNamespace();
`)
}

// Java: SymbolTableTest.verifySymbolReferencedInSecondFile(String, String, String)
func (c *SymbolTableTest) VerifySymbolReferencedInSecondFile(firstFile string, secondFile string, symbolName string) {
	table := c.CreateSymbolTableFromManySources(firstFile, secondFile)
	symbol :=
		symbolTableTest_findFirst(table.GetAllSymbols(), func(s *SymbolTable_Symbol) bool {
			return *s.GetSourceFileName() == "file1.js" && s.GetName() == symbolName
		})
	for ref := range table.GetReferences(symbol) {
		if *ref.GetNode().GetSourceFileName() == "file2.js" {
			return
		}
	}
	Fail("Did not find references in file2.js of symbol " + symbol.String())
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_moduleDefaultExportClass()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_moduleDefaultExportClass() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.module('some.Foo');
class Foo {}
exports = Foo;
`,
		`goog.module('some.bar');
const Foo = goog.require('some.Foo');
`,
		"Foo")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_moduleDefaultExportConstant()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_moduleDefaultExportConstant() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.module('some.constant');
const constant = 42;
exports = constant;
`,
		`goog.module('some.bar');
const constant = goog.require('some.constant');
`,
		// This is a tricky one. `const constant` require from the second file is a reference to the
		// `exports` symbol in the first file.
		"exports")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_moduleIndividualExports()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_moduleIndividualExports() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.module('some.foo');
exports.one = 1;
`,
		`goog.module('some.bar');
const {one} = goog.require('some.foo');
`,
		"one")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_requireType()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_requireType() {
	table :=
		c.CreateSymbolTableFromManySources(
			`goog.module('some.bar');
const {one} = goog.requireType('some.foo');
`,
			`goog.module('some.foo');
exports.one = 1;
`)
	symbol :=
		symbolTableTest_findFirst(table.GetAllSymbols(), func(s *SymbolTable_Symbol) bool { return *s.GetSourceFileName() == "file2.js" && s.GetName() == "one" })
	for ref := range table.GetReferences(symbol) {
		if *ref.GetNode().GetSourceFileName() == "file1.js" {
			return
		}
	}
	Fail("Did not find references in file2.js of symbol " + symbol.String())
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_provideDefaultExportClass()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_provideDefaultExportClass() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.provide('some.Foo');
some.Foo = class {}
`,
		`goog.module('some.bar');
const Foo = goog.require('some.Foo');
`,
		"some.Foo")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_provideDefaultExportPrimivite()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_provideDefaultExportPrimivite() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.provide('some.constant');
some.constant = 123;
`,
		`goog.module('some.bar');
const constant = goog.require('some.constant');
`,
		"some.constant")
}

// Java: SymbolTableTest.testNestedGoogProvides()
func (c *SymbolTableTest) TestNestedGoogProvides() {
	table :=
		c.CreateSymbolTableFromManySources(
			`goog.provide('a.b.c.d');
goog.provide('a.b.c.d.Foo');
goog.provide('a.b.c.d.Foo.Bar');
a.b.c.d.Foo = class {};
a.b.c.d.Foo.Bar = class {};
`)
	AssertThat(c.GetGlobalVar(table, "a.b.c.d.Foo")).IsNotNull()
	AssertThat(c.GetGlobalVar(table, "a.b.c.d.Foo.Bar")).IsNotNull()
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_provideIndividualExports()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_provideIndividualExports() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.provide('some.foo');
some.foo.one = 1;
`,
		`goog.module('some.bar');
const {one} = goog.require('some.foo');
`,
		"one")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_property_moduleImportsModule()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_property_moduleImportsModule() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.module('some.foo');
exports.one = 1;
`,
		`goog.module('some.bar');
const foo = goog.require('some.foo');
foo.one;
`,
		"one")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_moduleEarlyLexicographically()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_moduleEarlyLexicographically() {
	// Tests bug where imported module is named 'a' which is lexigraphically earlier than 'exports'
	// so 'a' becomes the symbol on which exported module properties are defined. The issue that
	// both 'a' and 'exports' have the same object type and we don't know which one is the true
	// object that declares it. So we just treat 'exports' as special name.
	c.VerifySymbolReferencedInSecondFile(
		`goog.module('some.foo');
const one = 1
exports.one = one;
`,
		`goog.module('some.bar');
const a = goog.require('some.foo');
a.one;
`,
		"one")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_property_moduleImportsProvide()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_property_moduleImportsProvide() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.provide('some.foo');
some.foo.one = 1;
`,
		`goog.module('some.bar');
const foo = goog.require('some.foo');
foo.one;
`,
		"one")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_property_provideImportsProvide()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_property_provideImportsProvide() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.provide('some.foo');
some.foo.one = 1;
`,
		`goog.provide('some.bar');
goog.require('some.foo');
some.foo.one;
`,
		"one")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_class()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_class() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.module('some.foo');
exports.Foo = class { doFoo() {} };
`,
		`goog.module('some.bar');
const {Foo} = goog.require('some.foo');
new Foo().doFoo();
`,
		"doFoo")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_typedef_module()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_typedef_module() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.module('some.foo');
/** @typedef {string} */ exports.StringType;
`,
		`goog.module('some.bar');
const foo = goog.require('some.foo');
const /** !foo.StringType */ k = '123';
`,
		"StringType")
}

// Java: SymbolTableTest.testGoogRequiredSymbolsConnectedToDefinitions_typedef_provide()
func (c *SymbolTableTest) TestGoogRequiredSymbolsConnectedToDefinitions_typedef_provide() {
	c.VerifySymbolReferencedInSecondFile(
		`goog.provide('some.foo');
/** @typedef {string} */ some.foo.StringType;
`,
		`goog.module('some.bar');
const foo = goog.require('some.foo');
const /** !foo.StringType */ k = '123';
`,
		"StringType")
}

// Java: SymbolTableTest.testClassPropertiesDisableModuleRewriting()
func (c *SymbolTableTest) TestClassPropertiesDisableModuleRewriting() {
	table :=
		c.CreateSymbolTableFromManySources(
			`goog.module('foo');
class Person {
  constructor() {
    /** @type {string} */
    this.lastName;
  }
}
exports = {Person};
`,
			`goog.module('bar');
const {Person} = goog.require('foo');
let /** !Person */ p;
p.lastName;
`)
	lastName :=
		symbolTableTest_findFirst(table.GetAllSymbols(), func(s *SymbolTable_Symbol) bool { return s.GetName() == "lastName" })
	AssertThat(table.GetReferences(lastName)).HasSize(2)
}

// Java: SymbolTableTest.testEnums()
func (c *SymbolTableTest) TestEnums() {
	table :=
		c.CreateSymbolTableFromManySources(
			`goog.module('foo');
/** @enum {number} */
const Color = {RED: 1};
exports.Color = Color;
Color.RED;
`,
			`goog.module('bar');
const foo = goog.require('foo');
foo.Color.RED;
`,
			`goog.module('baz');
const {Color} = goog.require('foo');
Color.RED;
`)
	red :=
		symbolTableTest_findFirst(table.GetAllSymbols(), func(s *SymbolTable_Symbol) bool { return s.GetName() == "RED" })
	// Make sure that RED belongs to the scope of Color that is defined in the first file and not
	// third. Because the third file also has "Color" enum that has the same type.
	AssertThat(table.GetScope(red).GetSymbolForScope().GetSourceFileName()).IsEqualTo("file1.js")
	AssertThat(table.GetReferences(red)).HasSize(4)
}

// Java: SymbolTableTest.testEnumsWithDirectExport()
func (c *SymbolTableTest) TestEnumsWithDirectExport() {
	table :=
		c.CreateSymbolTableFromManySources(
			`goog.module('foo');
/** @enum {number} */
exports.Color = {RED: 1} // exports.Color = Color;
`,
			`goog.module('bar');
const foo = goog.require('foo');
foo.Color.RED;
`,
			`goog.module('baz');
const {Color} = goog.require('foo');
Color.RED;
`)
	red :=
		symbolTableTest_findFirst(table.GetAllSymbols(), func(s *SymbolTable_Symbol) bool { return s.GetName() == "RED" })
	AssertThat(table.GetScope(red).GetSymbolForScope().GetSourceFileName()).IsEqualTo("file1.js")
	AssertThat(table.GetReferences(red)).HasSize(3)
}

// Java: SymbolTableTest.testEnumsProvidedAsNamespace()
func (c *SymbolTableTest) TestEnumsProvidedAsNamespace() {
	table :=
		c.CreateSymbolTableFromManySources(
			`goog.provide('foo.bar.Color');
/** @enum {number} */
foo.bar.Color = {RED: 1};
const /** !foo.bar.Color */ color =
  foo.bar.Color.RED;
`)
	red :=
		symbolTableTest_findFirst(table.GetAllSymbols(), func(s *SymbolTable_Symbol) bool { return s.GetName() == "RED" })
	AssertThat(table.GetScope(red).GetSymbolForScope().GetSourceFileName()).IsEqualTo("file1.js")
	AssertThat(table.GetReferences(red)).HasSize(2)
	AssertThat(table.GetReferences(c.GetGlobalVar(table, "foo.bar.Color"))).HasSize(3)
}

// Java: SymbolTableTest.testEnumsProvidedAsMemberOfNamespace()
func (c *SymbolTableTest) TestEnumsProvidedAsMemberOfNamespace() {
	table :=
		c.CreateSymbolTableFromManySources(
			`goog.provide('foo.bar');
/** @enum {number} */
foo.bar.Color = {RED: 1};
const /** !foo.bar.Color */ color =
  foo.bar.Color.RED;
`)
	red :=
		symbolTableTest_findFirst(table.GetAllSymbols(), func(s *SymbolTable_Symbol) bool { return s.GetName() == "RED" })
	// Make sure that RED belongs to the scope of Color that is defined in the first file and not
	// third. Because the third file also has "Color" enum that has the same type.
	AssertThat(table.GetScope(red).GetSymbolForScope().GetSourceFileName()).IsEqualTo("file1.js")
	AssertThat(table.GetReferences(red)).HasSize(2)
	colorEnum :=
		symbolTableTest_findFirst(table.GetAllSymbols(), func(s *SymbolTable_Symbol) bool { return s.GetName() == "Color" })
	AssertThat(table.GetReferences(colorEnum)).HasSize(3)
}

// Java: SymbolTableTest.testGlobalVarInExterns()
func (c *SymbolTableTest) TestGlobalVarInExterns() {
	table :=
		c.CreateSymbolTableWithExterns("customExternFn(1);", "function customExternFn(customExternArg) {}")
	fn := c.GetGlobalVar(table, "customExternFn")
	refs := table.GetReferenceList(fn)
	AssertThat(refs).HasSize(3)

	scope := table.GetEnclosingScope(refs[0].GetNode())
	AssertThat(scope.IsGlobalScope()).IsTrue()
	AssertThat(table.GetSymbolForScope(scope).GetName()).IsEqualTo(SymbolTable_GLOBAL_THIS)
}

// Java: SymbolTableTest.testLocalVarInExterns()
func (c *SymbolTableTest) TestLocalVarInExterns() {
	table := c.CreateSymbolTableWithExterns("", "function customExternFn(customExternArg) {}")
	arg := c.GetLocalVar(table, "customExternArg")
	refs := table.GetReferenceList(arg)
	AssertThat(refs).HasSize(1)

	fn := c.GetGlobalVar(table, "customExternFn")
	scope := table.GetEnclosingScope(refs[0].GetNode())
	AssertThat(scope.IsGlobalScope()).IsFalse()
	AssertThat(table.GetSymbolForScope(scope)).IsEqualTo(fn)
}

// Java: SymbolTableTest.testSymbolsForType()
func (c *SymbolTableTest) TestSymbolsForType() {
	table :=
		c.CreateSymbolTableWithDefaultExterns(
			`function random() { return 1; }
/** @constructor */ function Foo() {}
/** @constructor */ function Bar() {}
var x = random() ? new Foo() : new Bar();
`)

	x := c.GetGlobalVar(table, "x")
	foo := c.GetGlobalVar(table, "Foo")
	bar := c.GetGlobalVar(table, "Bar")
	fooPrototype := c.GetGlobalVar(table, "Foo.prototype")
	fn := c.GetGlobalVar(table, "Function")
	AssertThat(table.GetAllSymbolsForTypeOf(x)).ContainsExactly(foo, bar)
	AssertThat(table.GetAllSymbolsForTypeOf(foo)).ContainsExactly(fn)
	AssertThat(table.GetAllSymbolsForTypeOf(fooPrototype)).ContainsExactly(foo)
	AssertThat(table.GetSymbolDeclaredBy(foo.GetType().ToMaybeFunctionType())).IsEqualTo(foo)
}

// Java: SymbolTableTest.testStaticMethodReferences()
func (c *SymbolTableTest) TestStaticMethodReferences() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ var DomHelper = function(){};
/** method */ DomHelper.method = function() {};
function f() { var x = DomHelper; x.method() + x.method(); }
`)

	method := c.GetGlobalVar(table, "DomHelper").GetPropertyScope().GetSlot("method")
	AssertThat(table.GetReferences(method)).HasSize(3)
}

// Java: SymbolTableTest.testMethodReferences()
func (c *SymbolTableTest) TestMethodReferences() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ var DomHelper = function(){};
/** method */ DomHelper.prototype.method = function() {};
function f() {
  (new DomHelper()).method(); (new DomHelper()).method(); };
`)

	method := c.GetGlobalVar(table, "DomHelper.prototype.method")
	AssertThat(table.GetReferences(method)).HasSize(3)
}

// Java: SymbolTableTest.testMethodReferencesNullableClass()
func (c *SymbolTableTest) TestMethodReferencesNullableClass() {
	table :=
		c.CreateSymbolTable(
			`class DomHelper { method() {} }
let /** ?DomHelper */ helper;
helper.method();
`)

	method := c.GetGlobalVar(table, "DomHelper.prototype.method")
	AssertThat(table.GetReferences(method)).HasSize(2)
}

// Java: SymbolTableTest.testMethodReferencesNullableInterface()
func (c *SymbolTableTest) TestMethodReferencesNullableInterface() {
	table :=
		c.CreateSymbolTable(
			`/** @interface */ class DomHelper { method() {} }
let /** ?DomHelper */ helper;
helper.method();
`)

	method := c.GetGlobalVar(table, "DomHelper.prototype.method")
	AssertThat(table.GetReferences(method)).HasSize(2)
}

// Java: SymbolTableTest.testMethodReferencesNullableRecord()
func (c *SymbolTableTest) TestMethodReferencesNullableRecord() {
	table :=
		c.CreateSymbolTable(
			`/** @record */ function DomHelper() {}
DomHelper.prototype.method = function() {};
let /** ?DomHelper */ helper;
helper.method();
`)

	method := c.GetGlobalVar(table, "DomHelper.prototype.method")
	AssertThat(table.GetReferences(method)).HasSize(2)
}

// Java: SymbolTableTest.testRecordInheritance()
func (c *SymbolTableTest) TestRecordInheritance() {
	table :=
		c.CreateSymbolTable(
			`/** @record */ function DomHelper() {}
DomHelper.prototype.method = function() {};

/** @record @extends {DomHelper} */ function DomHelper2() {}
DomHelper2.prototype.method2 = function() {};

let /** ?DomHelper2 */ helper;
helper.method();
`)

	method := c.GetGlobalVar(table, "DomHelper.prototype.method")
	AssertThat(table.GetReferences(method)).HasSize(2)
}

// Java: SymbolTableTest.testUnionType()
func (c *SymbolTableTest) TestUnionType() {
	table :=
		c.CreateSymbolTable(
			`class First { method() {} }
class Second { method() {} }
let /** !Second|!First */ obj;
obj.method();
`)

	AssertThat(table.GetReferences(c.GetGlobalVar(table, "First.prototype.method"))).HasSize(2)
	AssertThat(table.GetReferences(c.GetGlobalVar(table, "Second.prototype.method"))).HasSize(2)
}

// Java: SymbolTableTest.testMultipleInterfaceInheritance()
func (c *SymbolTableTest) TestMultipleInterfaceInheritance() {
	table :=
		c.CreateSymbolTable(
			`/** @interface */
class First { doFirst() {} }

/** @interface */
class Second { doSecond() {} }

/** @interface @extends {First} @extends {Second} */
class Third {}

let /** !Third */ obj;
obj.doFirst();
obj.doSecond();
`)

	AssertThat(table.GetReferences(c.GetGlobalVar(table, "First.prototype.doFirst"))).HasSize(2)
	AssertThat(table.GetReferences(c.GetGlobalVar(table, "Second.prototype.doSecond"))).HasSize(2)
}

// Java: SymbolTableTest.testSuperClassMethodReferences()
func (c *SymbolTableTest) TestSuperClassMethodReferences() {
	table :=
		c.CreateSymbolTable(
			`var goog = {};
goog.inherits = function(a, b) {};
/** @constructor */ var A = function(){};
/** method */ A.prototype.method = function() {};
/**
 * @constructor
 * @extends {A}
 */
var B = function(){};
goog.inherits(B, A);
/** method */ B.prototype.method = function() {
  B.superClass_.method();
};
`)

	methodA := c.GetGlobalVar(table, "A.prototype.method")
	AssertThat(table.GetReferences(methodA)).HasSize(2)
}

// Java: SymbolTableTest.testMethodReferencesMissingTypeInfo()
func (c *SymbolTableTest) TestMethodReferencesMissingTypeInfo() {
	table :=
		c.CreateSymbolTable(
			`/**
 * @constructor
 * @extends {Missing}
 */ var DomHelper = function(){};
/** method */ DomHelper.prototype.method = function() {
  this.method();
};
function f() {
  (new DomHelper()).method();
};
`)

	method := c.GetGlobalVar(table, "DomHelper.prototype.method")
	AssertThat(table.GetReferences(method)).HasSize(3)
}

// Java: SymbolTableTest.testFieldReferencesMissingTypeInfo()
func (c *SymbolTableTest) TestFieldReferencesMissingTypeInfo() {
	table :=
		c.CreateSymbolTable(
			`/**
 * @constructor
 * @extends {Missing}
 */ var DomHelper = function(){ this.prop = 1; };
/** @type {number} */ DomHelper.prototype.prop = 2;
function f() {
  return (new DomHelper()).prop;
};
`)

	prop := c.GetGlobalVar(table, "DomHelper.prototype.prop")
	AssertThat(table.GetReferenceList(prop)).HasSize(3)

	AssertThat(c.GetLocalVar(table, "this.prop")).IsNull()
}

// Java: SymbolTableTest.testFieldReferences()
func (c *SymbolTableTest) TestFieldReferences() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ var DomHelper = function(){
  /** @type {number} */ this.field = 3;
};
function f() {
  return (new DomHelper()).field + (new DomHelper()).field; };
`)

	field := c.GetGlobalVar(table, "DomHelper.prototype.field")
	AssertThat(table.GetReferences(field)).HasSize(3)
}

// Java: SymbolTableTest.testUndeclaredFieldReferences()
func (c *SymbolTableTest) TestUndeclaredFieldReferences() {
	// We do not currently create symbol table entries for undeclared fields,
	// but this may change in the future.
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ var DomHelper = function(){};
DomHelper.prototype.method = function() {
  this.field = 3;
  return x.field;
}
`)

	field := c.GetGlobalVar(table, "DomHelper.prototype.field")
	AssertThat(field).IsNull()
}

// Java: SymbolTableTest.testPrototypeReferences()
func (c *SymbolTableTest) TestPrototypeReferences() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ function DomHelper() {}
DomHelper.prototype.method = function() {};
`)
	prototype := c.GetGlobalVar(table, "DomHelper.prototype")
	AssertThat(prototype).IsNotNull()

	refs := table.GetReferenceList(prototype)

	// One of the refs is implicit in the declaration of the function.
	AssertWithMessage("%s", jutil.ToJavaString(refs)).That(refs).HasSize(2)
}

// Java: SymbolTableTest.testPrototypeReferences2()
func (c *SymbolTableTest) TestPrototypeReferences2() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */
function Snork() {}
Snork.prototype.baz = 3;
`)
	prototype := c.GetGlobalVar(table, "Snork.prototype")
	AssertThat(prototype).IsNotNull()

	refs := table.GetReferenceList(prototype)
	AssertThat(refs).HasSize(2)
}

// Java: SymbolTableTest.testPrototypeReferences3()
func (c *SymbolTableTest) TestPrototypeReferences3() {
	table := c.CreateSymbolTable("/** @constructor */ function Foo() {}")
	fooPrototype := c.GetGlobalVar(table, "Foo.prototype")
	AssertThat(fooPrototype).IsNotNull()

	refs := table.GetReferenceList(fooPrototype)
	AssertThat(refs).HasSize(1)
	AssertThat(refs[0].GetNode().GetToken()).IsEqualTo(rhino.Token_NAME)

	// Make sure that the ctor and its prototype are declared at the
	// same node.
	AssertThat(table.GetReferenceList(c.GetGlobalVar(table, "Foo"))[0].GetNode()).
		IsEqualTo(refs[0].GetNode())
}

// Java: SymbolTableTest.testPrototypeReferences4()
func (c *SymbolTableTest) TestPrototypeReferences4() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ function Foo() {}
Foo.prototype = {bar: 3}
`)
	fooPrototype := c.GetGlobalVar(table, "Foo.prototype")
	AssertThat(fooPrototype).IsNotNull()

	refs := slices.Collect(table.GetReferences(fooPrototype))
	AssertThat(refs).HasSize(1)
	AssertThat(refs[0].GetNode().GetToken()).IsEqualTo(rhino.Token_GETPROP)
	AssertThat(refs[0].GetNode().GetQualifiedName()).IsEqualTo("Foo.prototype")
}

// Java: SymbolTableTest.testPrototypeReferences5()
func (c *SymbolTableTest) TestPrototypeReferences5() {
	table :=
		c.CreateSymbolTable("var goog = {}; /** @constructor */ goog.Foo = function() {};")
	fooPrototype := c.GetGlobalVar(table, "goog.Foo.prototype")
	AssertThat(fooPrototype).IsNotNull()

	refs := table.GetReferenceList(fooPrototype)
	AssertThat(refs).HasSize(1)
	AssertThat(refs[0].GetNode().GetToken()).IsEqualTo(rhino.Token_GETPROP)

	// Make sure that the ctor and its prototype are declared at the
	// same node.
	AssertThat(table.GetReferenceList(c.GetGlobalVar(table, "goog.Foo"))[0].GetNode()).
		IsEqualTo(refs[0].GetNode())
}

// Java: SymbolTableTest.testPrototypeReferences_es6Class()
func (c *SymbolTableTest) TestPrototypeReferences_es6Class() {
	table := c.CreateSymbolTable("class DomHelper { method() {} }")
	prototype := c.GetGlobalVar(table, "DomHelper.prototype")
	AssertThat(prototype).IsNotNull()
	refs := table.GetReferenceList(prototype)

	// The class declaration creates an implicit .prototype reference.
	AssertWithMessage("%s", jutil.ToJavaString(refs)).That(refs).HasSize(1)
	AssertNode(refs[0].GetNode().GetParent()).HasToken(rhino.Token_CLASS)
}

// Java: SymbolTableTest.testReferencesInJSDocType()
func (c *SymbolTableTest) TestReferencesInJSDocType() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ function Foo() {}
/** @type {Foo} */ var x;
/** @param {Foo} x */ function f(x) {}
/** @return {function(): Foo} */ function g() {}
/**
 * @constructor
 * @extends {Foo}
 */ function Sub() {}
`)
	foo := c.GetGlobalVar(table, "Foo")
	AssertThat(foo).IsNotNull()

	refs := table.GetReferenceList(foo)
	AssertThat(refs).HasSize(5)

	AssertThat(refs[0].GetNode().GetLineno()).IsEqualTo(1)
	AssertThat(refs[0].GetNode().GetCharno()).IsEqualTo(29)
	AssertThat(refs[0].GetNode().GetLength()).IsEqualTo(3)

	AssertThat(refs[1].GetNode().GetLineno()).IsEqualTo(2)
	AssertThat(refs[1].GetNode().GetCharno()).IsEqualTo(11)

	AssertThat(refs[2].GetNode().GetLineno()).IsEqualTo(3)
	AssertThat(refs[2].GetNode().GetCharno()).IsEqualTo(12)

	AssertThat(refs[3].GetNode().GetLineno()).IsEqualTo(4)
	AssertThat(refs[3].GetNode().GetCharno()).IsEqualTo(25)

	AssertThat(refs[4].GetNode().GetLineno()).IsEqualTo(7)
	AssertThat(refs[4].GetNode().GetCharno()).IsEqualTo(13)
}

// Java: SymbolTableTest.testReferencesInJSDocType2()
func (c *SymbolTableTest) TestReferencesInJSDocType2() {
	table :=
		c.CreateSymbolTableWithDefaultExterns("/** @param {string} x */ function f(x) {}")
	str := c.GetGlobalVar(table, "String")
	AssertThat(str).IsNotNull()

	refs := table.GetReferenceList(str)

	// We're going to pick up a lot of references from the externs,
	// so it's not meaningful to check the number of references.
	// We really want to make sure that all the references are in the externs,
	// except the last one.
	AssertThat(len(refs)).IsGreaterThan(1)

	last := len(refs) - 1
	for i := 0; i < len(refs); i++ {
		ref := refs[i]
		AssertThat(ref.GetNode().IsFromExterns()).IsEqualTo(i != last)
		if !ref.GetNode().IsFromExterns() {
			AssertThat(ref.GetNode().GetSourceFileName()).IsEqualTo("in1")
		}
	}
}

// Java: SymbolTableTest.testDottedReferencesInJSDocType()
func (c *SymbolTableTest) TestDottedReferencesInJSDocType() {
	table :=
		c.CreateSymbolTable(
			`var goog = {};
/** @constructor */ goog.Foo = function() {}
/** @type {goog.Foo} */ var x;
/** @param {goog.Foo} x */ function f(x) {}
/** @return {function(): goog.Foo} */ function g() {}
/**
 * @constructor
 * @extends {goog.Foo}
 */ function Sub() {}
`)
	foo := c.GetGlobalVar(table, "goog.Foo")
	AssertThat(foo).IsNotNull()

	refs := table.GetReferenceList(foo)
	AssertThat(refs).HasSize(5)

	AssertThat(refs[0].GetNode().GetLineno()).IsEqualTo(2)
	AssertThat(refs[0].GetNode().GetCharno()).IsEqualTo(25)
	AssertThat(refs[0].GetNode().GetLength()).IsEqualTo(3)

	AssertThat(refs[1].GetNode().GetLineno()).IsEqualTo(3)
	AssertThat(refs[1].GetNode().GetCharno()).IsEqualTo(16)

	AssertThat(refs[2].GetNode().GetLineno()).IsEqualTo(4)
	AssertThat(refs[2].GetNode().GetCharno()).IsEqualTo(17)

	AssertThat(refs[3].GetNode().GetLineno()).IsEqualTo(5)
	AssertThat(refs[3].GetNode().GetCharno()).IsEqualTo(30)

	AssertThat(refs[4].GetNode().GetLineno()).IsEqualTo(8)
	AssertThat(refs[4].GetNode().GetCharno()).IsEqualTo(18)
}

// Java: SymbolTableTest.testReferencesInJSDocName()
func (c *SymbolTableTest) TestReferencesInJSDocName() {
	code := "/** @param {Object} x */ function f(x) {}"
	table := c.CreateSymbolTable(code)
	x := c.GetLocalVar(table, "x")
	AssertThat(x).IsNotNull()

	refs := table.GetReferenceList(x)
	AssertThat(refs).HasSize(2)

	AssertThat(refs[0].GetNode().GetCharno()).IsEqualTo(jutil.StringIndexOf(code, "x) {"))
	AssertThat(refs[1].GetNode().GetCharno()).IsEqualTo(jutil.StringIndexOf(code, "x */"))
	AssertThat(refs[0].GetNode().GetSourceFileName()).IsEqualTo("in1")
}

// Java: SymbolTableTest.testLocalQualifiedNamesInLocalScopes()
func (c *SymbolTableTest) TestLocalQualifiedNamesInLocalScopes() {
	table := c.CreateSymbolTable("function f() { var x = {}; x.number = 3; }")
	xNumber := c.GetLocalVar(table, "x.number")
	AssertThat(xNumber).IsNotNull()
	AssertThat(table.GetScope(xNumber).IsGlobalScope()).IsFalse()

	AssertThat(xNumber.GetType().String()).IsEqualTo("number")
}

// Java: SymbolTableTest.testNaturalSymbolOrdering()
func (c *SymbolTableTest) TestNaturalSymbolOrdering() {
	table :=
		c.CreateSymbolTable(
			`/** @const */ var a = {};
/** @const */ a.b = {};
/** @param {number} x */ function f(x) {}
`)
	a := c.GetGlobalVar(table, "a")
	ab := c.GetGlobalVar(table, "a.b")
	f := c.GetGlobalVar(table, "f")
	x := c.GetLocalVar(table, "x")
	ordering := table.GetNaturalSymbolOrdering()
	c.AssertSymmetricOrdering(ordering, a, ab)
	c.AssertSymmetricOrdering(ordering, a, f)
	c.AssertSymmetricOrdering(ordering, f, ab)
	c.AssertSymmetricOrdering(ordering, f, x)
}

// Java: SymbolTableTest.testDeclarationDisagreement()
func (c *SymbolTableTest) TestDeclarationDisagreement() {
	table :=
		c.CreateSymbolTable(
			`/** @const */ var goog = goog || {};
/** @param {!Function} x */
goog.addSingletonGetter2 = function(x) {};
/** Wakka wakka wakka */
goog.addSingletonGetter = goog.addSingletonGetter2;
/** @param {!Function} x */
goog.addSingletonGetter = function(x) {};
`)

	method := c.GetGlobalVar(table, "goog.addSingletonGetter")
	refs := table.GetReferenceList(method)
	AssertThat(refs).HasSize(2)

	// Note that the declaration should show up second.
	AssertThat(method.GetDeclaration().GetNode().GetLineno()).IsEqualTo(7)
	AssertThat(refs[1].GetNode().GetLineno()).IsEqualTo(5)
}

// Java: SymbolTableTest.testMultipleExtends()
func (c *SymbolTableTest) TestMultipleExtends() {
	table :=
		c.CreateSymbolTable(
			`/** @const */ var goog = goog || {};
goog.inherits = function(x, y) {};
/** @constructor */
goog.A = function() { this.fieldA = this.constructor; };
/** @constructor */ goog.A.FooA = function() {};
/** @return {void} */ goog.A.prototype.methodA = function() {};
/**
 * @constructor
 * @extends {goog.A}
 */
goog.B = function() { this.fieldB = this.constructor; };
goog.inherits(goog.B, goog.A);
/** @return {void} */ goog.B.prototype.methodB = function() {};
/**
 * @constructor
 * @extends {goog.A}
 */
goog.B2 = function() { this.fieldB = this.constructor; };
goog.inherits(goog.B2, goog.A);
/** @constructor */ goog.B2.FooB = function() {};
/** @return {void} */ goog.B2.prototype.methodB = function() {};
/**
 * @constructor
 * @extends {goog.B}
 */
goog.C = function() { this.fieldC = this.constructor; };
goog.inherits(goog.C, goog.B);
/** @constructor */ goog.C.FooC = function() {};
/** @return {void} */ goog.C.prototype.methodC = function() {};
`)

	bCtor := c.GetGlobalVar(table, "goog.B.prototype.constructor")
	AssertThat(bCtor).IsNotNull()

	bRefs := table.GetReferenceList(bCtor)
	AssertThat(bRefs).HasSize(2)
	AssertThat(bCtor.GetDeclaration().GetNode().GetLineno()).IsEqualTo(11)

	cCtor := c.GetGlobalVar(table, "goog.C.prototype.constructor")
	AssertThat(cCtor).IsNotNull()

	cRefs := table.GetReferenceList(cCtor)
	AssertThat(cRefs).HasSize(2)
	AssertThat(cCtor.GetDeclaration().GetNode().GetLineno()).IsEqualTo(26)
}

// Java: SymbolTableTest.testJSDocAssociationWithBadNamespace()
func (c *SymbolTableTest) TestJSDocAssociationWithBadNamespace() {
	table :=
		c.CreateSymbolTable(
			// Notice that the declaration for "goog" is missing.
			// We want to recover anyway and print out what we know
			// about goog.Foo.
			"/** @constructor */ goog.Foo = function(){};")

	foo := c.GetGlobalVar(table, "goog.Foo")
	AssertThat(foo).IsNotNull()

	info := foo.GetJSDocInfo()
	AssertThat(info).IsNotNull()
	AssertThat(info.IsConstructor()).IsTrue()
}

// Java: SymbolTableTest.testMissingConstructorTag()
func (c *SymbolTableTest) TestMissingConstructorTag() {
	table :=
		c.CreateSymbolTable(
			`function F() {
  this.field1 = 3;
}
F.prototype.method1 = function() {
  this.field1 = 5;
};
(new F()).method1();
`)

	// Because the constructor tag is missing, this is going
	// to be missing a lot of inference.
	AssertThat(c.GetGlobalVar(table, "F.prototype.field1")).IsNull()

	sym := c.GetGlobalVar(table, "F.prototype.method1")
	AssertThat(table.GetReferenceList(sym)).HasSize(1)
}

// Java: SymbolTableTest.testTypeCheckingOff()
func (c *SymbolTableTest) TestTypeCheckingOff() {
	c.options = NewCompilerOptions()

	// Turning type-checking off is even worse than not annotating anything.
	table :=
		c.CreateSymbolTable(
			`/** @contstructor */
function F() {
  this.field1 = 3;
}
F.prototype.method1 = function() {
  this.field1 = 5;
};
(new F()).method1();
`)
	AssertThat(c.GetGlobalVar(table, "F.prototype.field1")).IsNull()
	AssertThat(c.GetGlobalVar(table, "F.prototype.method1")).IsNull()

	sym := c.GetGlobalVar(table, "F")
	AssertThat(table.GetReferenceList(sym)).HasSize(3)
}

// Java: SymbolTableTest.testSuperClassReference()
func (c *SymbolTableTest) TestSuperClassReference() {
	table :=
		c.CreateSymbolTable(
			`  var a = {b: {}};
/** @constructor */
a.b.BaseClass = function() {};
a.b.BaseClass.prototype.doSomething = function() {
  alert('hi');
};
/**
 * @constructor
 * @extends {a.b.BaseClass}
 */
a.b.DerivedClass = function() {};
goog.inherits(a.b.DerivedClass, a.b.BaseClass);
/** @override */
a.b.DerivedClass.prototype.doSomething = function() {
  a.b.DerivedClass.superClass_.doSomething();
};
`)

	bad := c.GetGlobalVar(table, "a.b.DerivedClass.superClass_.doSomething")
	AssertThat(bad).IsNull()

	good := c.GetGlobalVar(table, "a.b.BaseClass.prototype.doSomething")
	AssertThat(good).IsNotNull()

	refs := table.GetReferenceList(good)
	AssertThat(refs).HasSize(2)
	AssertThat(refs[1].GetNode().GetQualifiedName()).
		IsEqualTo("a.b.DerivedClass.superClass_.doSomething")
}

// Java: SymbolTableTest.testInnerEnum()
func (c *SymbolTableTest) TestInnerEnum() {
	table :=
		c.CreateSymbolTable(
			`var goog = {}; goog.ui = {};
  /** @constructor */
goog.ui.Zippy = function() {};
/** @enum {string} */
goog.ui.Zippy.EventType = { TOGGLE: 'toggle' };
`)

	eventType := c.GetGlobalVar(table, "goog.ui.Zippy.EventType")
	AssertThat(eventType).IsNotNull()
	AssertThat(eventType.GetType().IsEnumType()).IsTrue()

	toggle := c.GetGlobalVar(table, "goog.ui.Zippy.EventType.TOGGLE")
	AssertThat(toggle).IsNotNull()
}

// Java: SymbolTableTest.testMethodInAnonObject1()
func (c *SymbolTableTest) TestMethodInAnonObject1() {
	table := c.CreateSymbolTable("var a = {}; a.b = {}; a.b.c = function() {};")
	a := c.GetGlobalVar(table, "a")
	ab := c.GetGlobalVar(table, "a.b")
	abc := c.GetGlobalVar(table, "a.b.c")

	AssertThat(abc).IsNotNull()
	AssertThat(table.GetReferenceList(abc)).HasSize(1)

	AssertThat(a.GetType().String()).IsEqualTo("{b: {c: function(): undefined}}")
	AssertThat(ab.GetType().String()).IsEqualTo("{c: function(): undefined}")
	AssertThat(abc.GetType().String()).IsEqualTo("function(): undefined")
}

// Java: SymbolTableTest.testMethodInAnonObject2()
func (c *SymbolTableTest) TestMethodInAnonObject2() {
	table := c.CreateSymbolTable("var a = {b: {c: function() {}}};")
	a := c.GetGlobalVar(table, "a")
	ab := c.GetGlobalVar(table, "a.b")
	abc := c.GetGlobalVar(table, "a.b.c")

	AssertThat(abc).IsNotNull()
	AssertThat(table.GetReferenceList(abc)).HasSize(1)

	AssertThat(a.GetType().String()).IsEqualTo("{b: {c: function(): undefined}}")
	AssertThat(ab.GetType().String()).IsEqualTo("{c: function(): undefined}")
	AssertThat(abc.GetType().String()).IsEqualTo("function(): undefined")
}

// Java: SymbolTableTest.testJSDocOnlySymbol()
func (c *SymbolTableTest) TestJSDocOnlySymbol() {
	table :=
		c.CreateSymbolTable(
			`/**
 * @param {number} x
 * @param y
 */
var a;
`)
	x := c.GetDocVar(table, "x")
	AssertThat(x).IsNotNull()
	AssertThat(x.GetType().String()).IsEqualTo("number")
	AssertThat(table.GetReferenceList(x)).HasSize(1)

	y := c.GetDocVar(table, "y")
	AssertThat(x).IsNotNull()
	AssertThat(y.GetType()).IsNull()
	AssertThat(table.GetReferenceList(y)).HasSize(1)
}

// Java: SymbolTableTest.testNamespaceDefinitionOrder()
func (c *SymbolTableTest) TestNamespaceDefinitionOrder() {
	// Sometimes, weird things can happen where the files appear in
	// a strange order. We need to make sure we're robust against this.
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ goog.dom.Foo = function() {};
/** @const */ goog.dom = {};
/** @const */ var goog = {};
`)

	goog := c.GetGlobalVar(table, "goog")
	dom := c.GetGlobalVar(table, "goog.dom")
	foo := c.GetGlobalVar(table, "goog.dom.Foo")

	AssertThat(goog).IsNotNull()
	AssertThat(dom).IsNotNull()
	AssertThat(foo).IsNotNull()

	AssertThat(goog.GetPropertyScope().GetSlot("dom")).IsEqualTo(dom)
	AssertThat(dom.GetPropertyScope().GetSlot("Foo")).IsEqualTo(foo)
}

// Java: SymbolTableTest.testConstructorAlias()
func (c *SymbolTableTest) TestConstructorAlias() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ var Foo = function() {};
/** desc */ Foo.prototype.bar = function() {};
/** @constructor */ var FooAlias = Foo;
/** desc */ FooAlias.prototype.baz = function() {};
`)

	foo := c.GetGlobalVar(table, "Foo")
	fooAlias := c.GetGlobalVar(table, "FooAlias")
	bar := c.GetGlobalVar(table, "Foo.prototype.bar")
	baz := c.GetGlobalVar(table, "Foo.prototype.baz")
	bazAlias := c.GetGlobalVar(table, "FooAlias.prototype.baz")

	AssertThat(foo).IsNotNull()
	AssertThat(fooAlias).IsNotNull()
	AssertThat(bar).IsNotNull()
	AssertThat(baz).IsEqualTo(bazAlias)

	barScope := table.GetSymbolForScope(table.GetScope(bar))
	AssertThat(barScope).IsNotNull()

	bazScope := table.GetSymbolForScope(table.GetScope(baz))
	AssertThat(bazScope).IsNotNull()

	fooPrototype := foo.GetPropertyScope().GetSlot("prototype")
	AssertThat(fooPrototype).IsNotNull()

	AssertThat(barScope).IsEqualTo(fooPrototype)
	AssertThat(bazScope).IsEqualTo(fooPrototype)
}

// Java: SymbolTableTest.testSymbolForScopeOfNatives()
func (c *SymbolTableTest) TestSymbolForScopeOfNatives() {
	table := c.CreateSymbolTableWithDefaultExterns("")

	// From the externs.
	sliceArg := c.GetLocalVar(table, "sliceArg")
	AssertThat(sliceArg).IsNotNull()

	scope := table.GetSymbolForScope(table.GetScope(sliceArg))
	AssertThat(scope).IsNotNull()
	AssertThat(c.GetGlobalVar(table, "String.prototype.slice")).IsEqualTo(scope)

	proto := c.GetGlobalVar(table, "String.prototype")
	AssertThat(proto.GetDeclaration().GetNode().GetSourceFileName()).IsEqualTo("externs1")
}

// Java: SymbolTableTest.testJSDocNameVisibility()
func (c *SymbolTableTest) TestJSDocNameVisibility() {
	table :=
		c.CreateSymbolTable(
			`/** @public */ var foo;
/** @protected */ var bar;
/** @package */ var baz;
/** @private */ var quux;
var xyzzy;
`)

	AssertThat(c.GetGlobalVar(table, "foo").GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PUBLIC)
	AssertThat(c.GetGlobalVar(table, "bar").GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PROTECTED)
	AssertThat(c.GetGlobalVar(table, "baz").GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PACKAGE)
	AssertThat(c.GetGlobalVar(table, "quux").GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PRIVATE)
	AssertThat(c.GetGlobalVar(table, "xyzzy").GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	AssertThat(c.GetGlobalVar(table, "xyzzy").GetJSDocInfo()).IsNull()
}

// Java: SymbolTableTest.testJSDocNameVisibilityWithFileOverviewVisibility()
func (c *SymbolTableTest) TestJSDocNameVisibilityWithFileOverviewVisibility() {
	table :=
		c.CreateSymbolTable(
			`/** @fileoverview
 @package */
/** @public */ var foo;
/** @protected */ var bar;
/** @package */ var baz;
/** @private */ var quux;
var xyzzy;
`)
	AssertThat(c.GetGlobalVar(table, "foo").GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PUBLIC)
	AssertThat(c.GetGlobalVar(table, "bar").GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PROTECTED)
	AssertThat(c.GetGlobalVar(table, "baz").GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PACKAGE)
	AssertThat(c.GetGlobalVar(table, "quux").GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PRIVATE)
	AssertThat(c.GetGlobalVar(table, "xyzzy").GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PACKAGE)
	AssertThat(c.GetGlobalVar(table, "xyzzy").GetJSDocInfo()).IsNull()
}

// Java: SymbolTableTest.testJSDocPropertyVisibility()
func (c *SymbolTableTest) TestJSDocPropertyVisibility() {
	table :=
		c.CreateSymbolTable(
			`/** @constructor */ var Foo = function() {};
/** @public */ Foo.prototype.bar;
/** @protected */ Foo.prototype.baz;
/** @package */ Foo.prototype.quux;
/** @private */ Foo.prototype.xyzzy;
Foo.prototype.plugh;
/** @constructor @extends {Foo} */ var SubFoo = function() {};
/** @override */ SubFoo.prototype.bar = function() {};
/** @override */ SubFoo.prototype.baz = function() {};
/** @override */ SubFoo.prototype.quux = function() {};
/** @override */ SubFoo.prototype.xyzzy = function() {};
/** @override */ SubFoo.prototype.plugh = function() {};
`)

	AssertThat(c.GetGlobalVar(table, "Foo.prototype.bar").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_PUBLIC)
	AssertThat(c.GetGlobalVar(table, "Foo.prototype.baz").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_PROTECTED)
	AssertThat(c.GetGlobalVar(table, "Foo.prototype.quux").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_PACKAGE)
	AssertThat(c.GetGlobalVar(table, "Foo.prototype.xyzzy").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_PRIVATE)
	AssertThat(c.GetGlobalVar(table, "Foo.prototype.plugh").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	AssertThat(c.GetGlobalVar(table, "Foo.prototype.plugh").GetJSDocInfo()).IsNull()

	AssertThat(c.GetGlobalVar(table, "SubFoo.prototype.bar").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	AssertThat(c.GetGlobalVar(table, "SubFoo.prototype.baz").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	AssertThat(c.GetGlobalVar(table, "SubFoo.prototype.quux").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	AssertThat(c.GetGlobalVar(table, "SubFoo.prototype.xyzzy").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	AssertThat(c.GetGlobalVar(table, "SubFoo.prototype.plugh").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
}

// Java: SymbolTableTest.testJSDocPropertyVisibilityWithFileOverviewVisibility()
func (c *SymbolTableTest) TestJSDocPropertyVisibilityWithFileOverviewVisibility() {
	table :=
		c.CreateSymbolTable(
			`/** @fileoverview
 @package */
/** @constructor */ var Foo = function() {};
/** @public */ Foo.prototype.bar;
/** @protected */ Foo.prototype.baz;
/** @package */ Foo.prototype.quux;
/** @private */ Foo.prototype.xyzzy;
Foo.prototype.plugh;
/** @constructor @extends {Foo} */ var SubFoo = function() {};
/** @override @public */ SubFoo.prototype.bar = function() {};
/** @override @protected */ SubFoo.prototype.baz = function() {};
/** @override @package */ SubFoo.prototype.quux = function() {};
/** @override @private */ SubFoo.prototype.xyzzy = function() {};
/** @override */ SubFoo.prototype.plugh = function() {};
`)

	AssertThat(c.GetGlobalVar(table, "Foo.prototype.bar").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_PUBLIC)
	AssertThat(c.GetGlobalVar(table, "Foo.prototype.baz").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_PROTECTED)
	AssertThat(c.GetGlobalVar(table, "Foo.prototype.quux").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_PACKAGE)
	AssertThat(c.GetGlobalVar(table, "Foo.prototype.xyzzy").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_PRIVATE)
	AssertThat(c.GetGlobalVar(table, "Foo.prototype.plugh").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_PACKAGE)
	AssertThat(c.GetGlobalVar(table, "Foo.prototype.plugh").GetJSDocInfo()).IsNull()

	AssertThat(c.GetGlobalVar(table, "SubFoo.prototype.bar").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	AssertThat(c.GetGlobalVar(table, "SubFoo.prototype.baz").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	AssertThat(c.GetGlobalVar(table, "SubFoo.prototype.quux").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	AssertThat(c.GetGlobalVar(table, "SubFoo.prototype.xyzzy").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	AssertThat(c.GetGlobalVar(table, "SubFoo.prototype.plugh").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
}

// Java: SymbolTableTest.testPrototypeSymbolEqualityForTwoPathsToSamePrototype()
func (c *SymbolTableTest) TestPrototypeSymbolEqualityForTwoPathsToSamePrototype() {
	input :=
		`/**
* An employer.
*
* @param {String} name name of employer.
* @param {String} address address of employer.
* @constructor
*/
function Employer(name, address) {
this.name = name;
this.address = address;
}

/**
* @return {String} information about an employer.
*/
Employer.prototype.getInfo = function() {
return this.name + '.' + this.address;
};
`
	table := c.CreateSymbolTable(input)
	employer := c.GetGlobalVar(table, "Employer")
	AssertThat(employer).IsNotNull()
	propertyScope := employer.GetPropertyScope()
	AssertThat(propertyScope).IsNotNull()
	prototypeOfEmployer := propertyScope.GetQualifiedSlot("prototype")
	AssertThat(prototypeOfEmployer).IsNotNull()

	employerPrototype := c.GetGlobalVar(table, "Employer.prototype")
	AssertThat(employerPrototype).IsNotNull()

	rhinotesting.NewEqualsTester().AddEqualityGroup(employerPrototype, prototypeOfEmployer).TestEquals(c.t)
}

// Java: SymbolTableTest.testRestParameter()
func (c *SymbolTableTest) TestRestParameter() {
	input := "function f(...x) {} f(1, 2, 3);"

	table := c.CreateSymbolTable(input)

	f := c.GetGlobalVar(table, "f")
	AssertThat(f).IsNotNull()

	x := table.GetParameterInFunction(f, "x")
	AssertThat(x).IsNotNull()
}

// Java: SymbolTableTest.testGetEnclosingScope_Global()
func (c *SymbolTableTest) TestGetEnclosingScope_Global() {
	table := c.CreateSymbolTable("const baz = 1;")

	baz := c.GetGlobalVar(table, "baz")
	bazScope := table.GetEnclosingScope(baz.GetDeclarationNode())

	AssertThat(bazScope).IsNotNull()
	AssertThat(bazScope.IsGlobalScope()).IsTrue()
}

// Java: SymbolTableTest.testGetEnclosingScope_GlobalNested()
func (c *SymbolTableTest) TestGetEnclosingScope_GlobalNested() {
	table := c.CreateSymbolTable("{ const baz = 1; }")

	baz := c.GetLocalVar(table, "baz")
	bazScope := table.GetEnclosingScope(baz.GetDeclarationNode())

	AssertThat(bazScope).IsNotNull()
	AssertThat(bazScope.IsBlockScope()).IsTrue()
	AssertThat(bazScope.GetParentScope().IsGlobalScope()).IsTrue()
}

// Java: SymbolTableTest.testGetEnclosingScope_FunctionNested()
func (c *SymbolTableTest) TestGetEnclosingScope_FunctionNested() {
	table := c.CreateSymbolTable("function foo() { { const baz = 1; } }")

	baz := c.GetLocalVar(table, "baz")
	bazScope := table.GetEnclosingScope(baz.GetDeclarationNode())

	AssertThat(bazScope).IsNotNull()
	AssertThat(bazScope.IsBlockScope()).IsTrue()
	AssertThat(bazScope.GetParentScope().IsBlockScope()).IsTrue()
	foo := c.GetGlobalVar(table, "foo").GetDeclarationNode().GetParent()
	AssertThat(bazScope.GetParentScope().GetParentScope().GetRootNode()).IsEqualTo(foo)
}

// Java: SymbolTableTest.testGetEnclosingFunctionScope_GlobalScopeNoNesting()
func (c *SymbolTableTest) TestGetEnclosingFunctionScope_GlobalScopeNoNesting() {
	table := c.CreateSymbolTable("const baz = 1;")

	baz := c.GetGlobalVar(table, "baz")
	bazFunctionScope := table.GetEnclosingFunctionScope(baz.GetDeclarationNode())

	AssertThat(bazFunctionScope).IsNotNull()
	AssertThat(bazFunctionScope.IsGlobalScope()).IsTrue()
}

// Java: SymbolTableTest.testGetEnclosingFunctionScope_GlobalScopeNested()
func (c *SymbolTableTest) TestGetEnclosingFunctionScope_GlobalScopeNested() {
	table := c.CreateSymbolTable("{ const baz = 1; }")

	baz := c.GetLocalVar(table, "baz")
	bazFunctionScope := table.GetEnclosingFunctionScope(baz.GetDeclarationNode())

	AssertThat(bazFunctionScope).IsNotNull()
	AssertThat(bazFunctionScope.IsGlobalScope()).IsTrue()
}

// Java: SymbolTableTest.testGetEnclosingFunctionScope_FunctionNoNestedScopes()
func (c *SymbolTableTest) TestGetEnclosingFunctionScope_FunctionNoNestedScopes() {
	table := c.CreateSymbolTable("function foo() { const baz = 1; }")

	baz := c.GetLocalVar(table, "baz")
	bazFunctionScope := table.GetEnclosingFunctionScope(baz.GetDeclarationNode())

	AssertThat(bazFunctionScope).IsNotNull()
	foo := c.GetGlobalVar(table, "foo").GetDeclarationNode().GetParent()
	AssertThat(bazFunctionScope.GetRootNode()).IsEqualTo(foo)
}

// Java: SymbolTableTest.testGetEnclosingFunctionScope_FunctionNested()
func (c *SymbolTableTest) TestGetEnclosingFunctionScope_FunctionNested() {
	table := c.CreateSymbolTable("function foo() { { { const baz = 1; } } }")

	baz := c.GetLocalVar(table, "baz")
	bazFunctionScope := table.GetEnclosingFunctionScope(baz.GetDeclarationNode())

	AssertThat(bazFunctionScope).IsNotNull()
	foo := c.GetGlobalVar(table, "foo").GetDeclarationNode().GetParent()
	AssertThat(bazFunctionScope.GetRootNode()).IsEqualTo(foo)
}

// Java: SymbolTableTest.testGetEnclosingFunctionScope_FunctionInsideFunction()
func (c *SymbolTableTest) TestGetEnclosingFunctionScope_FunctionInsideFunction() {
	table := c.CreateSymbolTable("function foo() { function baz() {} }")

	baz := c.GetLocalVar(table, "baz")
	bazFunctionScope := table.GetEnclosingFunctionScope(baz.GetDeclarationNode())

	AssertThat(bazFunctionScope).IsNotNull()
	foo := c.GetGlobalVar(table, "foo").GetDeclarationNode().GetParent()
	AssertThat(bazFunctionScope.GetRootNode()).IsEqualTo(foo)
}

// Java: SymbolTableTest.testSymbolSuperclassStaticInheritance()
func (c *SymbolTableTest) TestSymbolSuperclassStaticInheritance() {
	// set this option so that typechecking sees untranspiled classes.
	table :=
		c.CreateSymbolTable(
			`class Bar { static staticMethod() {} }
class Foo extends Bar {
  act() { Foo.staticMethod(); }
}
`)

	foo := c.GetGlobalVar(table, "Foo")
	bar := c.GetGlobalVar(table, "Bar")

	barType := jutil.CheckNotNull(bar.GetType().ToMaybeFunctionType())
	fooType := jutil.CheckNotNull(foo.GetType().ToMaybeFunctionType())

	superclassCtorType := fooType.GetSuperClassConstructor()

	// TypeSubject.isEqualTo is JSType.equals.
	AssertThat(superclassCtorType).IsEqualTo(barType)
	superSymbol := table.GetSymbolDeclaredBy(superclassCtorType)
	AssertThat(superSymbol).IsEqualTo(bar)

	barStaticMethod := c.GetGlobalVar(table, "Bar.staticMethod")
	barStaticMethodRefs := table.GetReferenceList(barStaticMethod)
	AssertThat(barStaticMethodRefs).HasSize(2)
	AssertThat(barStaticMethodRefs[0]).IsEqualTo(barStaticMethod.GetDeclaration())
	// We recognize that the call to Foo.staticMethod() is actually a call to Bar.staticMethod().
	AssertNode(barStaticMethodRefs[1].GetNode()).MatchesQualifiedName("Foo.staticMethod")
}

// Java: SymbolTableTest.testSuperKeywords()
func (c *SymbolTableTest) TestSuperKeywords() {
	table :=
		c.CreateSymbolTable(
			`class Parent { doFoo() {} }
class Child extends Parent {
  constructor() { super(); }
  doFoo() { super.doFoo(); }
}
`)

	parentClass := c.GetGlobalVar(table, "Parent")
	numberOfSuperReferences :=
		symbolTableTest_count(table.GetReferenceList(parentClass), func(ref *SymbolTable_Reference) bool { return ref.GetNode().IsSuper() })
	// TODO(b/112359882): change number of refs to 2 once Es6ConvertSuper pass is moved after type
	// checks.
	AssertThat(numberOfSuperReferences).IsEqualTo(1)
}

// Java: SymbolTableTest.testDuplicatedWindowExternsMerged()
func (c *SymbolTableTest) TestDuplicatedWindowExternsMerged() {
	// Add window so that it triggers logic that defines all global externs on window.
	// See DeclaredGlobalExternsOnWindow.java pass.
	externs :=
		`/** @externs */
var window
/** @type {number} */ var foo;
`
	mainCode :=
		`foo = 2;
window.foo = 1;
`
	table := c.CreateSymbolTableWithExterns(mainCode, externs)

	refsPerFile := jutil.NewJavaHashMap[string, int]()
	for _, reference := range table.GetReferenceList(c.GetGlobalVar(table, "foo")) {
		if !reference.GetNode().IsIndexable() {
			continue
		}
		file := reference.GetSourceFile().GetName()
		refsPerFile.Put(file, refsPerFile.GetOrDefault(file, 0)+1)
	}
	AssertThat(refsPerFile).ContainsExactlyEntriesIn(symbolTableTest_mapOf("in1", 2, "externs1", 1))
}

// Java: SymbolTableTest.testDuplicatedWindowExternsMergedWithWindowPrototype()
func (c *SymbolTableTest) TestDuplicatedWindowExternsMergedWithWindowPrototype() {
	// Add window so that it triggers logic that defines all global externs on window.
	// See DeclaredGlobalExternsOnWindow.java pass.
	// Also add Window class as it affects symbols (e.g. window.foo is set on Window.prototype.foo
	// instead).
	externs :=
		`/** @externs */
/** @constructor */ function Window() {}
/** @type {!Window} */ var window;
/** @type {number} */ var foo;
`
	mainCode :=
		`foo = 2;
window.foo = 1;
`
	table := c.CreateSymbolTableWithExterns(mainCode, externs)

	refsPerFile := jutil.NewJavaHashMap[string, int]()
	for _, reference := range table.GetReferenceList(c.GetGlobalVar(table, "foo")) {
		file := reference.GetSourceFile().GetName()
		refsPerFile.Put(file, refsPerFile.GetOrDefault(file, 0)+1)
	}
	AssertThat(refsPerFile).ContainsExactlyEntriesIn(symbolTableTest_mapOf("in1", 2, "externs1", 1))
}

// Java: SymbolTableTest.testTypedefWithProperties()
func (c *SymbolTableTest) TestTypedefWithProperties() {
	table :=
		c.CreateSymbolTable(
			`/** @typedef {{prop: string}} */ let TypedefWithProperties;
let /** ?TypedefWithProperties */ obj;
obj.prop;
`)

	prop := c.GetGlobalVar(table, "TypedefWithProperties.prop")
	AssertThat(table.GetReferences(prop)).HasSize(2)
}

// Java: SymbolTableTest.assertSymmetricOrdering(Ordering<Symbol>, Symbol, Symbol)
func (c *SymbolTableTest) AssertSymmetricOrdering(ordering func(a, b *SymbolTable_Symbol) int, first *SymbolTable_Symbol, second *SymbolTable_Symbol) {
	AssertThat(ordering(first, first)).IsEqualTo(0)
	AssertThat(ordering(second, second)).IsEqualTo(0)
	AssertThat(ordering(first, second)).IsLessThan(0)
	AssertThat(ordering(second, first)).IsGreaterThan(0)
}

// Java: SymbolTableTest.getGlobalVar(SymbolTable, String)
func (c *SymbolTableTest) GetGlobalVar(table *SymbolTable, name string) *SymbolTable_Symbol {
	return table.GetGlobalScope().GetQualifiedSlot(name)
}

// Java: SymbolTableTest.getDocVar(SymbolTable, String)
func (c *SymbolTableTest) GetDocVar(table *SymbolTable, name string) *SymbolTable_Symbol {
	for _, sym := range table.GetAllSymbols() {
		if sym.IsDocOnlyParameter() && sym.GetName() == name {
			return sym
		}
	}
	return nil
}

// Java: SymbolTableTest.getLocalVar(SymbolTable, String)
func (c *SymbolTableTest) GetLocalVar(table *SymbolTable, name string) *SymbolTable_Symbol {
	for scope := range table.GetAllScopes() {
		if !scope.IsGlobalScope() &&
			scope.IsLexicalScope() &&
			scope.GetQualifiedSlot(name) != nil {
			return scope.GetQualifiedSlot(name)
		}
	}
	return nil
}

// GetVars returns all non-extern vars.
//
// Java: SymbolTableTest.getVars(SymbolTable)
func (c *SymbolTableTest) GetVars(table *SymbolTable) []*SymbolTable_Symbol {
	result := []*SymbolTable_Symbol{}
	for _, symbol := range table.GetAllSymbols() {
		if symbol.GetDeclaration() != nil && !symbol.GetDeclaration().GetNode().IsFromExterns() {
			result = append(result, symbol)
		}
	}
	return result
}

// Java: SymbolTableTest.createSymbolTableWithDefaultExterns(String)
func (c *SymbolTableTest) CreateSymbolTableWithDefaultExterns(input string) *SymbolTable {
	return c.CreateSymbolTableWithExterns(input, CompilerTypeTestCase_DEFAULT_EXTERNS)
}

// Java: SymbolTableTest.createSymbolTable(String)
func (c *SymbolTableTest) CreateSymbolTable(input string) *SymbolTable {
	return c.CreateSymbolTableWithExterns(input, "")
}

// Java: SymbolTableTest.createSymbolTable(String, String)
func (c *SymbolTableTest) CreateSymbolTableWithExterns(input string, externsCode string) *SymbolTable {
	inputs := []*SourceFile{SourceFile_FromCode("in1", input)}
	externs := []*SourceFile{SourceFile_FromCode("externs1", externsCode)}

	compiler := NewCompilerWithErrorManager(NewBlackHoleErrorManager())
	compiler.Compile(externs, inputs, c.options)
	return c.AssertSymbolTableValid(compiler.BuildKnownSymbolTable())
}

// Java: SymbolTableTest.createSymbolTableFromManySources(String...)
func (c *SymbolTableTest) CreateSymbolTableFromManySources(inputs ...string) *SymbolTable {
	sources := []*SourceFile{}
	for i := 0; i < len(inputs); i++ {
		sources = append(sources, SourceFile_FromCode("file"+strconv.Itoa(i+1)+".js", inputs[i]))
	}
	externs := []*SourceFile{}

	compiler := NewCompilerWithErrorManager(NewBlackHoleErrorManager())
	compiler.Compile(externs, sources, c.options)
	return c.AssertSymbolTableValid(compiler.BuildKnownSymbolTable())
}

// AssertSymbolTableValid asserts that the symbol table meets some
// invariants. Returns the same table for easy chaining.
//
// Java: SymbolTableTest.assertSymbolTableValid(SymbolTable)
func (c *SymbolTableTest) AssertSymbolTableValid(table *SymbolTable) *SymbolTable {
	allSymbols := jutil.NewJavaHashSet[*SymbolTable_Symbol]()
	for _, sym := range table.GetAllSymbols() {
		allSymbols.Add(sym)
	}
	for _, sym := range table.GetAllSymbols() {
		// Make sure that grabbing the symbol's scope and looking it up
		// again produces the same symbol.
		AssertThat(table.GetScope(sym).GetQualifiedSlot(sym.GetName())).IsEqualTo(sym)

		for ref := range table.GetReferences(sym) {
			// Make sure that the symbol and reference are mutually linked.
			AssertThat(ref.GetSymbol()).IsEqualTo(sym)
		}

		symbolForScope := table.GetSymbolForScope(table.GetScope(sym))
		if symbolForScope != nil {
			AssertWithMessage("The %s has a scope that shouldn't exist; a zombie scope.", sym).
				That(allSymbols).
				Contains(symbolForScope)
		}
	}

	// Make sure that the global "this" is declared at the first input root.
	global := c.GetGlobalVar(table, SymbolTable_GLOBAL_THIS)
	AssertThat(global).IsNotNull()
	AssertThat(global.GetDeclaration()).IsNotNull()
	AssertThat(global.GetDeclaration().GetNode().GetToken()).IsEqualTo(rhino.Token_SCRIPT)

	globalRefs := table.GetReferenceList(global)

	// The main reference list should never contain the synthetic declaration
	// for the global root.
	AssertThat(globalRefs).DoesNotContain(global.GetDeclaration())

	return table
}

// symbolTableTest_mapOf is ImmutableMap.of(k1, v1, k2, v2) for the
// expected side of containsExactly.
//
// Go-only: language-adapter — Truth's MapSubject.containsExactly(k, v, ...).
func symbolTableTest_mapOf(k1 string, v1 int, k2 string, v2 int) *jutil.LinkedHashMap[string, int] {
	m := jutil.NewLinkedHashMap[string, int]()
	m.Put(k1, v1)
	m.Put(k2, v2)
	return m
}

// symbolTableTest_findFirst is list.stream().filter(pred).findFirst().get().
//
// Go-only: language-adapter — a Java stream pipeline.
func symbolTableTest_findFirst[T any](list []T, pred func(T) bool) T {
	for _, e := range list {
		if pred(e) {
			return e
		}
	}
	panic(jutil.NewNoSuchElementExceptionMsg("No value present"))
}

// symbolTableTest_count is list.stream().filter(pred).count().
//
// Go-only: language-adapter — a Java stream pipeline.
func symbolTableTest_count[T any](list []T, pred func(T) bool) int64 {
	var n int64
	for _, e := range list {
		if pred(e) {
			n++
		}
	}
	return n
}

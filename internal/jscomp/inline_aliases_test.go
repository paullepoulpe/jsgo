// Port of com/google/javascript/jscomp/InlineAliasesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/InlineAliasesTest.java)

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
)

// InlineAliasesTest: unit tests for InlineAndCollapseProperties.InlineAliases.
type InlineAliasesTest struct {
	CompilerTestCaseBase
}

// Java: InlineAliasesTest()
func NewInlineAliasesTest() *InlineAliasesTest {
	c := &InlineAliasesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestInlineAliasesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInlineAliasesTest() })
}

// SetUp runs CompilerTestCase's @Before method, then InlineAliasesTest's
// (JUnit runs superclass @Before methods first).
func (c *InlineAliasesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.CustomSetUp()
}

// Java: InlineAliasesTest.getProcessor(Compiler)
func (c *InlineAliasesTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("es6NormalizeClasses", func(comp AbstractCompiler) CompilerPass { return NewEs6NormalizeClasses(comp) }))
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"inlineAndCollapseProperties",
			func(comp AbstractCompiler) CompilerPass {
				return InlineAndCollapseProperties_Builder_(compiler).
					SetPropertyCollapseLevel(CompilerOptions_PropertyCollapseLevel_NONE).
					Build()
			}))
	return optimizer
}

// Java: InlineAliasesTest.customSetUp()
func (c *InlineAliasesTest) CustomSetUp() {
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
}

// Java: InlineAliasesTest.getOptions()
func (c *InlineAliasesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
	return options
}

// Returns the number of times the pass should be run before results are verified.
//
// <p>This pass is not idempotent so we only run it once.
// Java: InlineAliasesTest.testSimpleAliasInJSDoc_isUnchanged()
func (c *InlineAliasesTest) TestSimpleAliasInJSDoc_isUnchanged() {
	c.TestSame("/** @constructor */ function Foo(){} const alias = Foo; /** @type {alias} */ var x;")
	c.TestSame(
		`var ns={};
/** @constructor */ function Foo(){};
/** @const */ ns.alias = Foo;
/** @type {ns.alias} */ var x;
`)
	c.TestSame(
		`/** @const */
var ns={};
/** @constructor */ function Foo(){};
Foo.Subfoo = class {};
/** @const */ ns.alias = Foo;
/** @type {ns.alias.Subfoo} */ var x;
`)
}

// Java: InlineAliasesTest.testSimpleAliasInCode()
func (c *InlineAliasesTest) TestSimpleAliasInCode() {
	c.Test(
		`/** @constructor */ function Foo(){};
Foo.Subfoo = class {};
var /** @const */ alias = Foo; var x = new alias;
`,
		`/** @constructor */ function Foo(){};
Foo.Subfoo = class {};var /** @const */ alias = Foo; var x = new Foo;
`)
	c.Test(
		`var ns={}; /** @constructor */ function Foo(){};
Foo.Subfoo = class {};
/** @const */ ns.alias = Foo; var x = new ns.alias;
`,
		`var ns={}; /** @constructor */ function Foo(){};
Foo.Subfoo = class {};
/** @const */ ns.alias = Foo; var x = new Foo;
`)
	c.Test(
		`var ns={}; /** @constructor */ function Foo(){};
Foo.Subfoo = class {};
/** @const */ ns.alias = Foo; var x = new ns.alias.Subfoo;
`,
		`var ns={}; /** @constructor */ function Foo(){};
Foo.Subfoo = class {};
/** @const */ ns.alias = Foo; var x = new Foo.Subfoo;
`)
}

// Java: InlineAliasesTest.testAliasQualifiedName()
func (c *InlineAliasesTest) TestAliasQualifiedName() {
	c.TestSame(
		`/** @const */
var ns = {};
ns.Foo = function(){};
ns.Foo.Subfoo = class {};
/** @const */ ns.alias = ns.Foo;
/** @type {ns.alias.Subfoo} */ var x;
`)
	c.Test(
		`var ns = {};
ns.Foo = function(){};
ns.Foo.Subfoo = class {};
/** @const */ ns.alias = ns.Foo;
var x = new ns.alias.Subfoo;
`,
		`var ns = {};
ns.Foo = function(){};
ns.Foo.Subfoo = class {};
/** @const */ ns.alias = ns.Foo;
var x = new ns.Foo.Subfoo;
`)
}

// Java: InlineAliasesTest.testHoistedAliasesInCode()
func (c *InlineAliasesTest) TestHoistedAliasesInCode() {
	// Unqualified
	c.Test(
		`function Foo(){};
function Bar(){ var x = alias; };
var /** @const */ alias = Foo;
`,
		`function Foo(){};
function Bar(){ var x = Foo; };
var /** @const */ alias = Foo;
`)
	// Qualified
	c.Test(
		`var ns = {};
ns.Foo = function(){};
function Bar(){ var x = ns.alias; };
/** @const */ ns.alias = ns.Foo;
`,
		`var ns = {};
ns.Foo = function(){};
function Bar(){ var x = ns.Foo; };
/** @const */ ns.alias = ns.Foo;
`)
}

// Java: InlineAliasesTest.testAliasCycleError()
func (c *InlineAliasesTest) TestAliasCycleError() {
	c.TestError(
		`/** @const */ var x = y;
/** @const */ var y = x;
`,
		InlineAndCollapseProperties_ALIAS_CYCLE)
}

// Java: InlineAliasesTest.testTransitiveAliases()
func (c *InlineAliasesTest) TestTransitiveAliases() {
	c.Test(
		`/** @const */ var ns = {};
/** @constructor */ ns.Foo = function() {};
/** @constructor */ ns.Foo.Bar = function() {};
var /** @const */ alias = ns.Foo;
var /** @const */ alias2 = alias.Bar;
var x = new alias2
`,
		`/** @const */ var ns = {};
/** @constructor */ ns.Foo = function() {};
/** @constructor */ ns.Foo.Bar = function() {};
var /** @const */ alias = ns.Foo;
var /** @const */ alias2 = ns.Foo.Bar;
// Note: in order to replace "alias2" with "ns.Foo.Bar", we would either have to do
// multiple traversals of the AST in InlineAliases, or mark alias2 as an alias of
// ns.Foo.Bar in the GlobalNamespace after replacing "alias2 = alias.Bar" with
// "alias2 = ns.Foo.Bar"
var x = new alias2;
`)
}

// Java: InlineAliasesTest.testAliasChains()
func (c *InlineAliasesTest) TestAliasChains() {
	// Unqualified
	c.Test(
		`/** @constructor */ var Foo = function() {};
var /** @const */ alias1 = Foo;
var /** @const */ alias2 = alias1;
var x = new alias2
`,
		`/** @constructor */ var Foo = function() {};
var /** @const */ alias1 = Foo;
var /** @const */ alias2 = Foo;
var x = new Foo;
`)
	// Qualified
	c.Test(
		`/** @const */ var ns = {};
/** @constructor */ ns.Foo = function() {};
var /** @const */ alias1 = ns.Foo;
var /** @const */ alias2 = alias1;
var x = new alias2
`,
		`/** @const */ var ns = {};
/** @constructor */ ns.Foo = function() {};
var /** @const */ alias1 = ns.Foo;
var /** @const */ alias2 = ns.Foo;
var x = new ns.Foo;
`)
}

// Java: InlineAliasesTest.testAliasedEnums()
func (c *InlineAliasesTest) TestAliasedEnums() {
	c.Test(
		"/** @enum {number} */ var E = { A : 1 }; var /** @const */ alias = E.A; alias;",
		"/** @enum {number} */ var E = { A : 1 }; var /** @const */ alias = E.A; E.A;")
}

// Java: InlineAliasesTest.testIncorrectConstAnnotationDoesntCrash()
func (c *InlineAliasesTest) TestIncorrectConstAnnotationDoesntCrash() {
	c.TestSame("var x = 0; var /** @const */ alias = x; alias = 5; use(alias);")
	c.TestSame("var x = 0; var ns={}; /** @const */ ns.alias = x; ns.alias = 5; use(ns.alias);")
}

// Java: InlineAliasesTest.testRedefinedAliasesNotRenamed()
func (c *InlineAliasesTest) TestRedefinedAliasesNotRenamed() {
	c.TestSame("var x = 0; var /** @const */ alias = x; x = 5; use(alias);")
}

// Java: InlineAliasesTest.testDefinesAreNotInlined()
func (c *InlineAliasesTest) TestDefinesAreNotInlined() {
	c.TestSame("var ns = {Foo: 0}; var /** @define {number} */ alias = ns.Foo; use(alias);")
}

// Java: InlineAliasesTest.testConstWithTypesAreNotInlined()
func (c *InlineAliasesTest) TestConstWithTypesAreNotInlined() {
	c.TestSame(
		`var /** @type {number} */ n = 5
var /** @const {number} */ alias = n;
var x = use(alias)
`)
}

// Java: InlineAliasesTest.testShadowedAliasesNotRenamed()
func (c *InlineAliasesTest) TestShadowedAliasesNotRenamed() {
	c.TestSame(
		`var ns = {};
ns.Foo = function(){};
var /** @const */ alias = ns.Foo;
function f(alias) {
  var x = alias
}
`)
	c.TestSame(
		`var ns = {};
ns.Foo = function(){};
var /** @const */ alias = ns.Foo;
function f() {
  var /** @const */ alias = 5;
  var x = alias
}
`)
	c.TestSame(
		`/** @const */
var x = y;
function f() {
  var x = 123;
  function g() {
    return x;
  }
}
`)
}

// Java: InlineAliasesTest.testShadowedAliasesNotRenamed_withBlockScope()
func (c *InlineAliasesTest) TestShadowedAliasesNotRenamed_withBlockScope() {
	c.TestSame(
		`var ns = {};
ns.Foo = function(){};
var /** @const */ alias = ns.Foo;
if (true) {
  const alias = 5;
  var x = alias
}
`)
	c.TestSame(
		`/** @const */
var x = y;
if (true) {
  const x = 123;
  function g() {
    return x;
  }
}
`)
}

// Java: InlineAliasesTest.testES6VarAliasClassDeclarationWithNew()
func (c *InlineAliasesTest) TestES6VarAliasClassDeclarationWithNew() {
	c.Test(
		"class Foo{}; var /** @const */ alias = Foo; var x = new alias;",
		"class Foo{}; var /** @const */ alias = Foo; var x = new Foo;")
}

// Java: InlineAliasesTest.testES6VarAliasClassDeclarationWithoutNew()
func (c *InlineAliasesTest) TestES6VarAliasClassDeclarationWithoutNew() {
	c.Test(
		"class Foo{}; var /** @const */ alias = Foo; var x = alias;",
		"class Foo{}; var /** @const */ alias = Foo; var x = Foo;")
}

// Java: InlineAliasesTest.testNoInlineAliasesInsideClassConstructor()
func (c *InlineAliasesTest) TestNoInlineAliasesInsideClassConstructor() {
	c.TestSame(
		`class Foo {
 constructor(x) {
     /** @const */
     this.x = class {};
     var /** @const */ alias1 = this.x;
     var /** @const */ alias2 = alias1;
     var z = new alias2;
 }
}
`)
}

// Java: InlineAliasesTest.testArrayDestructuringVarAssign()
func (c *InlineAliasesTest) TestArrayDestructuringVarAssign() {
	c.Test(
		`var Foo = class {};
var /** @const */ A = Foo;
var a = [5, A];
var [one, two] = a;
`,
		`var Foo = class {};
var /** @const */ A = Foo;
var a = [5, Foo];
var [one, two] = a;
`)
}

// Java: InlineAliasesTest.testArrayDestructuringFromFunction()
func (c *InlineAliasesTest) TestArrayDestructuringFromFunction() {
	c.Test(
		`var Foo = class {};
var /** @const */ A = Foo;
function f() {
  return [A, 3];
}
var a, b;
[a, b] = f();
`,
		`var Foo = class {};
var /** @const */ A = Foo;
function f() {
  return [Foo, 3];
}
var a, b;
[a, b] = f();
`)
}

// Java: InlineAliasesTest.testArrayDestructuringSwapIsNotInlined()
func (c *InlineAliasesTest) TestArrayDestructuringSwapIsNotInlined() {
	c.TestSame(
		`var Foo = class {};
var /** @const */ A = Foo;
var temp = 3;
[A, temp] = [temp, A];
`)
}

// Java: InlineAliasesTest.testArrayDestructuringSwapIsNotInlinedWithClassDeclaration()
func (c *InlineAliasesTest) TestArrayDestructuringSwapIsNotInlinedWithClassDeclaration() {
	c.TestSame(
		`class Foo {};
var /** @const */ A = Foo;
var temp = 3;
[A, temp] = [temp, A];
`)
}

// Java: InlineAliasesTest.testArrayDestructuringAndRedefinedAliasesNotRenamed()
func (c *InlineAliasesTest) TestArrayDestructuringAndRedefinedAliasesNotRenamed() {
	c.TestSame("var x = 0; var /** @const */ alias = x; [x] = [5]; use(alias);")
}

// Java: InlineAliasesTest.testArrayDestructuringTwoVarsAndRedefinedAliasesNotRenamed()
func (c *InlineAliasesTest) TestArrayDestructuringTwoVarsAndRedefinedAliasesNotRenamed() {
	c.TestSame(
		`var x = 0;
var /** @const */ alias = x;
var y = 5;
[x] = [y];
use(alias);
`)
}

// Java: InlineAliasesTest.testObjectDestructuringBasicAssign()
func (c *InlineAliasesTest) TestObjectDestructuringBasicAssign() {
	c.Test(
		`var Foo = class {};
var /** @const */ A = Foo;
var o = {p: A, q: 5};
var {p, q} = o;
`,
		`var Foo = class {};
var /** @const */ A = Foo;
var o = {p: Foo, q: 5};
var {p, q} = o;
`)
}

// Java: InlineAliasesTest.testObjectDestructuringAssignWithoutDeclaration()
func (c *InlineAliasesTest) TestObjectDestructuringAssignWithoutDeclaration() {
	c.Test(
		`var Foo = class {};
var /** @const */ A = Foo;
({a, b} = {a: A, b: A});
`,
		`var Foo = class {};
var /** @const */ A = Foo;
({a, b} = {a: Foo, b: Foo});
`)
}

// Java: InlineAliasesTest.testObjectDestructuringAssignNewVarNames()
func (c *InlineAliasesTest) TestObjectDestructuringAssignNewVarNames() {
	c.Test(
		`var Foo = class {};
var /** @const */ A = Foo;
var o = {p: A, q: true};
var {p: newName1, q: newName2} = o;
`,
		`var Foo = class {};
var /** @const */ A = Foo;
var o = {p: Foo, q: true};
var {p: newName1, q: newName2} = o;
`)
}

// Java: InlineAliasesTest.testObjectDestructuringDefaultVals()
func (c *InlineAliasesTest) TestObjectDestructuringDefaultVals() {
	c.IgnoreWarnings(DiagnosticGroups_MISSING_PROPERTIES)
	c.Test(
		`var Foo = class {};
var /** @const */ A = Foo;
var {a = A, b = A} = {a: 13};
`,
		`var Foo = class {};
var /** @const */ A = Foo;
var {a = Foo, b = Foo} = {a: 13};
`)
}

// Java: InlineAliasesTest.testArrayDestructuringWithParameter()
func (c *InlineAliasesTest) TestArrayDestructuringWithParameter() {
	c.Test(
		`var Foo = class {};
var /** @const */ A = Foo;
function f([name, val]) {
   alert(name, val);
}
f([A, A]);
`,
		`var Foo = class {};
var /** @const */ A = Foo;
function f([name, val]) {
   alert(name, val);
}
f([Foo, Foo]);
`)
}

// Java: InlineAliasesTest.testObjectDestructuringWithParameters()
func (c *InlineAliasesTest) TestObjectDestructuringWithParameters() {
	c.Test(
		`var Foo = class {};
var /** @const */ A = Foo;
function g({
   name: n,
   val: v
}) {
   alert(n, v);
}
g({
   name: A,
   val: A
});
`,
		`var Foo = class {};
var /** @const */ A = Foo;
function g({
   name: n,
   val: v
}) {
   alert(n, v);
}
g({
   name: Foo,
   val: Foo
});
`)
}

// Java: InlineAliasesTest.testObjectDestructuringWithParametersAndStyleShortcut()
func (c *InlineAliasesTest) TestObjectDestructuringWithParametersAndStyleShortcut() {
	c.Test(
		`var Foo = class {};
var /** @const */ A = Foo;
function h({
   name,
   val
}) {
   alert(name, val);
}
h({name: A, val: A});
`,
		`var Foo = class {};
var /** @const */ A = Foo;
function h({
   name,
   val
}) {
   alert(name, val);
}
h({name: Foo, val: Foo});
`)
}

// Java: InlineAliasesTest.testSimpleConstAliasInCode()
func (c *InlineAliasesTest) TestSimpleConstAliasInCode() {
	c.Test(
		"/** @constructor */ function Foo(){}; const alias = Foo; var x = new alias;",
		"/** @constructor */ function Foo(){}; const alias = Foo; var x = new Foo;")
}

// Java: InlineAliasesTest.testSimpleLetAliasInCode()
func (c *InlineAliasesTest) TestSimpleLetAliasInCode() {
	c.Test(
		"/** @constructor */ function Foo(){}; let /** @const */ alias = Foo; var x = new alias;",
		"/** @constructor */ function Foo(){}; let /** @const */ alias = Foo; var x = new Foo;")
}

// Java: InlineAliasesTest.testClassExtendsAlias1()
func (c *InlineAliasesTest) TestClassExtendsAlias1() {
	c.Test(
		"class Foo{} const alias = Foo; class Bar extends alias {}",
		"class Foo {} const alias = Foo; class Bar extends Foo {}")
}

// Java: InlineAliasesTest.testClassExtendsAlias2()
func (c *InlineAliasesTest) TestClassExtendsAlias2() {
	c.Test(
		"var ns = {}; ns.Foo = class {}; const alias = ns.Foo; class Bar extends alias {}",
		"var ns = {}; ns.Foo = class {}; const alias = ns.Foo; class Bar extends ns.Foo {}")
}

// Java: InlineAliasesTest.testBlockScopedAlias()
func (c *InlineAliasesTest) TestBlockScopedAlias() {
	c.TestSame("function Foo() {} if (true) { const alias = Foo; alias; }")
}

// Java: InlineAliasesTest.testVarAliasDeclaredInBlockScope()
func (c *InlineAliasesTest) TestVarAliasDeclaredInBlockScope() {
	c.TestSame("function Foo() {} { var /** @const */ alias = Foo; alias; }")
}

// Java: InlineAliasesTest.testDontInlineEscapedQnameProperty()
func (c *InlineAliasesTest) TestDontInlineEscapedQnameProperty() {
	c.TestSame(
		c.Externs("function use(obj) {}"),
		c.Srcs(
			`/** @const */
var ns = {};
ns.foo = 3;
const alias = ns.foo;
use(ns);
// "ns" escapes and we don't know if the value of "ns.foo" has also changed, so
// we cannot replace "alias" with "ns.foo".
alert(alias);
`))
}

// Java: InlineAliasesTest.testDoInlineEscapedConstructorProperty()
func (c *InlineAliasesTest) TestDoInlineEscapedConstructorProperty() {
	// TODO(b/80429954): this is unsafe. The call to use(Foobar) could have changed the value of
	// Foobar.foo
	c.Test(
		c.Externs("function use(obj) {}"),
		c.Srcs(
			`/** @constructor */
function Foobar() {}
Foobar.foo = 3;
const alias = Foobar.foo;
use(Foobar);
alert(alias);
`),
		c.Expected(
			`/** @constructor */
function Foobar() {}
Foobar.foo = 3;
const alias = Foobar.foo;
use(Foobar);
alert(Foobar.foo);
`))
}

// Java: InlineAliasesTest.testForwardedExport()
func (c *InlineAliasesTest) TestForwardedExport() {
	c.TestSame(
		`const proto = {};
/** @const */
proto.google = {};
/** @const */
proto.google.type = {};
proto.google.type.Date = class {};
const alias = proto;
function f() {
  const d = new alias.google.type.Date();
  const proto = 0;
}
`)
}

// Java: InlineAliasesTest.testForwardedExportNested()
func (c *InlineAliasesTest) TestForwardedExportNested() {
	c.TestSame(
		`const proto = {};
/** @const */
proto.google = {};
/** @const */
proto.google.type = {};
proto.google.type.Date = class {};
const alias = proto.google;
function f() {
  const d = new alias.type.Date();
  const proto = 0;
}
`)
}

// Java: InlineAliasesTest.testQualifiedNameSetViaUnaryDecrementNotInlined()
func (c *InlineAliasesTest) TestQualifiedNameSetViaUnaryDecrementNotInlined() {
	c.TestSame(
		`const a = {b: 0, c: 0};
const v1 = a.b;
a.b--;
const v2 = a.b;
a.b--;
use(v1 + v2);
`)
}

// Java: InlineAliasesTest.testQualifiedNameSetViaUnaryIncrementNotInlined()
func (c *InlineAliasesTest) TestQualifiedNameSetViaUnaryIncrementNotInlined() {
	c.TestSame(
		`const a = {b: 0};
const v1 = a.b;
a.b++;
const v2 = a.b;
a.b++;
use(v1 + v2);
`)
}

// Java: InlineAliasesTest.testAliasOfStubDeclaration()
func (c *InlineAliasesTest) TestAliasOfStubDeclaration() {
	c.Test(
		`const a = {};
var stubDeclaration;
/** @const */
a.b = stubDeclaration;
alert(a.b);
`,
		`const a = {};
var stubDeclaration;
/** @const */
a.b = stubDeclaration;
alert(stubDeclaration);
`)
}

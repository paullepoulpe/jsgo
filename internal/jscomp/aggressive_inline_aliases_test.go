// Port of com/google/javascript/jscomp/AggressiveInlineAliasesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/AggressiveInlineAliasesTest.java)

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

// AggressiveInlineAliasesTest: tests for the aggressive inlining part of
// InlineAndCollapseProperties.
type AggressiveInlineAliasesTest struct {
	CompilerTestCaseBase
}

// Java: AggressiveInlineAliasesTest.EXTERNS
const AggressiveInlineAliasesTest_EXTERNS = `var window
function alert(s) {}
function parseInt(s) {}
function String() {};
var arguments
`

// Java: AggressiveInlineAliasesTest()
func NewAggressiveInlineAliasesTest() *AggressiveInlineAliasesTest {
	c := &AggressiveInlineAliasesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, AggressiveInlineAliasesTest_EXTERNS)
	return c
}

func TestAggressiveInlineAliasesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewAggressiveInlineAliasesTest() })
}

// SetUp runs CompilerTestCase's @Before method, then
// AggressiveInlineAliasesTest's (JUnit runs superclass @Before methods
// first).
func (c *AggressiveInlineAliasesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.CustomSetUp()
}

// Java: AggressiveInlineAliasesTest.getProcessor(Compiler)
func (c *AggressiveInlineAliasesTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("es6NormalizeClasses", func(comp AbstractCompiler) CompilerPass { return NewEs6NormalizeClasses(comp) }))
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"inlineAndCollapseProperties",
			func(comp AbstractCompiler) CompilerPass {
				return InlineAndCollapseProperties_Builder_(compiler).
					SetPropertyCollapseLevel(CompilerOptions_PropertyCollapseLevel_ALL).
					SetChunkOutputType(CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE).
					SetHaveModulesBeenRewritten(false).
					SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER).
					TestAggressiveInliningOnly(c.ValidateGlobalNamespace).
					Build()
			}))
	return optimizer
}

// Java: AggressiveInlineAliasesTest.customSetUp()
func (c *AggressiveInlineAliasesTest) CustomSetUp() {
	c.EnableNormalize()
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
}

// Java: AggressiveInlineAliasesTest.testHasOwnProperty()
func (c *AggressiveInlineAliasesTest) TestHasOwnProperty() {
	c.Test(
		"var a = {'b': 1, 'c': 1}; var alias = a;    alert(alias.hasOwnProperty('c'));", //
		"var a = {'b': 1, 'c': 1}; var alias = null; alert(a.hasOwnProperty('c'));")
}

// Java: AggressiveInlineAliasesTest.test_b19179602()
func (c *AggressiveInlineAliasesTest) Test_b19179602() {
	c.Test(
		`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
/** @constructor */
function f() {
  while (true) {
    var b = a.b;
    alert(b.staticProp);
  }
}
`,
		`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
/** @constructor */
function f() {
  for(; true; ) {
    var b = null; // replaced with null
    alert(a.b.staticProp); // value inlined
  }
}
`)
	c.Test(
		`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
/** @constructor */
function f() {
  while (true) {
    var b = a.b;
    alert(b?.staticProp);
  }
}
`,
		`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
/** @constructor */
function f() {
  for(; true; ) {
    var b = null; // replaced with null
    alert(a.b?.staticProp); // value inlined
  }
}
`)
	c.TestSame(
		`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
/** @constructor */
function f() {
  for (; true;) {
    var b = a?.b; // not inlined when optional chain
    alert(b?.staticProp);
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.test_b19179602_declareOutsideLoop()
func (c *AggressiveInlineAliasesTest) Test_b19179602_declareOutsideLoop() {
	c.Test(
		`var a = {};
a.b = function() {};
a.b.staticProp = 5;
function f() {
  var b = a.b;
  while (true) {
    alert(b.staticProp);
  }
}
`,
		`var a = {};
a.b = function() {};
a.b.staticProp = 5;
function f() {
  var b = null; // replaced with null
  for (;true;) // converted to for-loop
    alert(a.b.staticProp); // inlined
}
`)
	c.Test(
		`var a = {};
a.b = function() {};
a.b.staticProp = 5;
function f() {
  var b = a.b;
  while (true) {
    alert(b?.staticProp); // reference with opt chain
  }
}
`,
		`var a = {};
a.b = function() {};
a.b.staticProp = 5;
function f() {
  var b = null; // replaced with null
  for (;true;) // converted to for-loop
    alert(a.b?.staticProp); // inlined
}
`)
	c.TestSame(
		"var a = {};\n" +
			"a.b = function() {};\n" +
			"a.b.staticProp = 5;\n" +
			"function f() {\n" +
			"// Not an alias for `a.b` when it's an optional chain.\n" +
			"  var b = a?.b;\n" +
			"  for (; true;) {\n" +
			"    alert(b?.staticProp);\n" +
			"  }\n" +
			"}\n")
}

// Java: AggressiveInlineAliasesTest.testCtorAliasedMultipleTimesNoWarning1()
func (c *AggressiveInlineAliasesTest) TestCtorAliasedMultipleTimesNoWarning1() {
	// We can't inline the alias, but there are no unsafe property accesses.
	c.TestSame(
		`/** @constructor */
function a() {}
function f() {
  var alias = a;
  use(alias.prototype); // The prototype is not collapsible.
  alias.apply(); // Can't collapse externs properties.
  alias = function() {}
}
`)
}

// Java: AggressiveInlineAliasesTest.testCtorAliasedMultipleTimesNoWarning2()
func (c *AggressiveInlineAliasesTest) TestCtorAliasedMultipleTimesNoWarning2() {
	c.TestSame(
		`/** @constructor */
function a() {}
/** @nocollapse */
a.staticProp = 5; // Explicitly forbid collapsing a.staticProp.
function f() {
  var alias = a;
  use(alias.staticProp); // Safe because we don't collapse a.staticProp
  alias = function() {}
}
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalCtorAliasedMultipleTimes()
func (c *AggressiveInlineAliasesTest) TestGlobalCtorAliasedMultipleTimes() {
	// TODO(lharker): Also warn for unsafe global ctor aliasing
	c.TestSame(
		`/** @constructor */
function a() {}
a.staticProp = 5;
var alias = a;
use(alias.staticProp); // Unsafe because a.staticProp becomes a$staticProp.
alias = function() {}
`)
}

// Java: AggressiveInlineAliasesTest.testCtorAliasedMultipleTimesWarning1()
func (c *AggressiveInlineAliasesTest) TestCtorAliasedMultipleTimesWarning1() {
	c.TestSame(
		`/** @constructor */
function a() {}
a.staticProp = 5;
function f() {
  var alias = a;
  use(alias.staticProp); // Unsafe because a.staticProp becomes a$staticProp.
  alias = function() {}
}
`,
		InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING)
}

// Java: AggressiveInlineAliasesTest.testCtorAliasedMultipleTimesWarning2()
func (c *AggressiveInlineAliasesTest) TestCtorAliasedMultipleTimesWarning2() {
	c.TestSame(
		`/** @constructor */
function a() {}
a.b = {};
a.b.staticProp = 5;
function f() {
  var alias = a;
  use(alias.b.staticProp); // Unsafe because a.b.staticProp becomes a$b$staticProp.
  alias = function() {}
}
`,
		InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING)
	c.TestSame(
		`/** @constructor */
function a() {}
a.b = {};
a.b.staticProp = 5;
function f() {
  var alias = a;
// Make sure a property reference via optional chain is still
// recognized as unsafe.
  use(alias?.b.staticProp); // Unsafe because a.b.staticProp becomes a$b$staticProp.
  alias = function() {}
}
`,
		InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING)
}

// Java: AggressiveInlineAliasesTest.testCtorAliasedMultipleTimesWarning3()
func (c *AggressiveInlineAliasesTest) TestCtorAliasedMultipleTimesWarning3() {
	c.TestSame(
		`/** @constructor */
function a() {}
a.staticProp = 5;
function f(alias) {
  alias();
  alias = a;
  use(alias.staticProp); // Unsafe because a.staticProp becomes a$staticProp.
}
`,
		InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING)
}

// Java: AggressiveInlineAliasesTest.testCtorAliasedMultipleTimesWarning4()
func (c *AggressiveInlineAliasesTest) TestCtorAliasedMultipleTimesWarning4() {
	c.TestWarning(
		`/** @constructor */
function a() {}
a.staticProp = 5;
function f() {
  if (true) {
    var alias = a;
    use(alias.staticProp);
  } else {
    alias = {staticProp: 34};
    use(alias.staticProp);
  }
}
`,
		InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING)
}

// Java: AggressiveInlineAliasesTest.testAliasingOfReassignedProperty1()
func (c *AggressiveInlineAliasesTest) TestAliasingOfReassignedProperty1() {
	c.TestSame("var obj = {foo: 3}; var foo = obj.foo; obj.foo = 42; alert(foo);")
	c.TestSame("var obj = {foo: 3}; var foo = obj?.foo; obj.foo = 42; alert(foo);")
}

// Java: AggressiveInlineAliasesTest.testAliasingOfReassignedProperty2()
func (c *AggressiveInlineAliasesTest) TestAliasingOfReassignedProperty2() {
	c.TestSame("var obj = {foo: 3}; var foo = obj.foo; obj = {}; alert(foo);")
}

// Java: AggressiveInlineAliasesTest.testAliasingOfReassignedProperty3()
func (c *AggressiveInlineAliasesTest) TestAliasingOfReassignedProperty3() {
	c.TestSame("var obj = {foo: {bar: 3}}; var bar = obj.foo.bar; obj.foo = {}; alert(bar);")
	c.TestSame("var obj = {foo: {bar: 3}}; var bar = obj.foo?.bar; obj.foo = {}; alert(bar);")
	c.TestSame("var obj = {foo: {bar: 3}}; var bar = obj?.foo.bar; obj.foo = {}; alert(bar);")
}

// Java: AggressiveInlineAliasesTest.testAliasingOfReassignedProperty4()
func (c *AggressiveInlineAliasesTest) TestAliasingOfReassignedProperty4() {
	// Note: it should be safe to inline aliases for properties of ns below, even though ns
	// has multiple definitions. Not inlining "foo" to "ns.ctor.foo" actually causes bad code later,
	// because CollapseProperties unsafely collapses aliased ctor properties.
	c.Test(
		`var ns = {};
/** @constructor */ ns.ctor = function() {};
ns.ctor.foo = 3;
var foo = ns.ctor.foo;
ns = ns || {}; // safe reinitialization of ns.
alert(foo);
`,
		`var ns = {};
/** @constructor */ ns.ctor = function() {};
ns.ctor.foo = 3;
var foo = null;
ns = ns || {};
alert(ns.ctor.foo);
`)
	c.TestSame(
		`var ns = {};
/** @constructor */ ns.ctor = function() {};
ns.ctor.foo = 3;
// optional chain means this isn't a real alias
var foo = ns.ctor?.foo;
ns = ns || {}; // safe reinitialization of ns.
alert(foo);
`)
}

// Java: AggressiveInlineAliasesTest.testAliasingOfReassignedProperty5()
func (c *AggressiveInlineAliasesTest) TestAliasingOfReassignedProperty5() {
	c.Test(
		"var obj = {foo: {bar: 3}}; var bar = obj.foo.bar; var obj; alert(bar);",
		"var obj = {foo: {bar: 3}}; var bar =        null;          alert(obj.foo.bar);")
}

// Java: AggressiveInlineAliasesTest.testAddPropertyToChildFuncOfUncollapsibleObjectInLocalScope()
func (c *AggressiveInlineAliasesTest) TestAddPropertyToChildFuncOfUncollapsibleObjectInLocalScope() {
	c.Test(
		`var a = {};
a.b = function () {};
a.b.x = 0;
var c = a;
(function() { a.b.y = 1; })();
a.b.x;
a.b.y;
`,
		`var a = {};
a.b = function() {};
a.b.x = 0;
var c = null;
(function() { a.b.y = 1; })();
a.b.x;
a.b.y;
`)
}

// Java: AggressiveInlineAliasesTest.testAddPropertyToChildOfUncollapsibleCtorInLocalScope()
func (c *AggressiveInlineAliasesTest) TestAddPropertyToChildOfUncollapsibleCtorInLocalScope() {
	c.Test(
		`var a = function() {};
a.b = { x: 0 };
var c = a;
(function() { a.b.y = 0; })();
a.b.y;
`,
		`var a = function() {};
a.b = { x: 0 };
var c = null;
(function() { a.b.y = 0; })();
a.b.y;
`)
}

// Java: AggressiveInlineAliasesTest.testAddPropertyToChildOfUncollapsibleFunctionInLocalScope()
func (c *AggressiveInlineAliasesTest) TestAddPropertyToChildOfUncollapsibleFunctionInLocalScope() {
	c.Test(
		`function a() {} a.b = { x: 0 };
var c = a;
(function() { a.b.y = 0; })();
a.b.y;
`,
		`function a() {} a.b = { x: 0 };
var c = null;
(function() { a.b.y = 0; })();
a.b.y;
`)
}

// Java: AggressiveInlineAliasesTest.testAddPropertyToChildTypeOfUncollapsibleObjectInLocalScope()
func (c *AggressiveInlineAliasesTest) TestAddPropertyToChildTypeOfUncollapsibleObjectInLocalScope() {
	c.Test(
		`var a = {};
a.b = function () {};
a.b.x = 0;
var c = a;
(function() { a.b.y = 1; })();
a.b.x;
a.b.y;
`,
		`var a = {};
a.b = function() {};
a.b.x = 0;
var c = null;
(function() { a.b.y = 1; })();
a.b.x;
a.b.y;
`)
}

// Java: AggressiveInlineAliasesTest.testAddPropertyToUncollapsibleCtorInLocalScopeDepth1()
func (c *AggressiveInlineAliasesTest) TestAddPropertyToUncollapsibleCtorInLocalScopeDepth1() {
	c.Test(
		`var a = function() {};
var c = a;
(function() { a.b = 0;
})();
a.b;
`,
		`var a = function() {};
var c = null;
(function() { a.b = 0 })();
a.b;
`)
}

// Java: AggressiveInlineAliasesTest.testAddPropertyToUncollapsibleCtorInLocalScopeDepth2()
func (c *AggressiveInlineAliasesTest) TestAddPropertyToUncollapsibleCtorInLocalScopeDepth2() {
	c.Test(
		`var a = {};
a.b = function () {};
var d = a.b;
(function() { a.b.c = 0; })();
a.b.c;
`,
		`var a = {};
a.b = function() {};
var d = null;
(function() { a.b.c = 0; })();
a.b.c;
`)
}

// Java: AggressiveInlineAliasesTest.testAddPropertyToUncollapsibleNamedCtorInLocalScopeDepth1()
func (c *AggressiveInlineAliasesTest) TestAddPropertyToUncollapsibleNamedCtorInLocalScopeDepth1() {
	c.Test(
		`function a() {} var a$b;
var c = a;
(function() { a$b = 0;
})();
a$b;
`,
		`function a() {} var a$b;
var c = null;
(function() { a$b = 0;
})();
a$b;
`)
}

// Java: AggressiveInlineAliasesTest.testAddPropertyToUncollapsibleObjectInLocalScopeDepth1()
func (c *AggressiveInlineAliasesTest) TestAddPropertyToUncollapsibleObjectInLocalScopeDepth1() {
	c.Test(
		`var a = {};
var c = a;
use(c);
(function() { a.b = 0;
})();
a.b;
`,
		`var a = {};
var c = null;
use(a);
(function() { a.b = 0;
})();
a.b;
`)
}

// Java: AggressiveInlineAliasesTest.testAddPropertyToUncollapsibleObjectInLocalScopeDepth2()
func (c *AggressiveInlineAliasesTest) TestAddPropertyToUncollapsibleObjectInLocalScopeDepth2() {
	c.Test(
		`var a = {};
a.b = {};
var d = a.b;
use(d);
(function() { a.b.c = 0; })();
a.b.c;
`,
		`var a = {};
a.b = {};
var d = null;
use(a.b);
(function() { a.b.c = 0; })();
a.b.c;
`)
	c.TestSame(
		`var a = {};
a.b = {};
var d = a?.b; // optional chain makes this no longer an alias
use(d);
(function() { a.b.c = 0; })();
a.b.c;
`)
}

// Java: AggressiveInlineAliasesTest.testAliasCreatedForClassDepth1_1()
func (c *AggressiveInlineAliasesTest) TestAliasCreatedForClassDepth1_1() {
	c.Test(
		`var a = {};
a.b = function() {};
var c = a;
c.b = 0;
a.b != c.b;
`,
		`var a = {};
a.b = function() {};
var c = null;
a.b = 0;
a.b != a.b;
`)
	c.Test(
		`var a = {};
a.b = function() {};
var c = a;
c.b = 0;
a.b == c.b;
`,
		`var a = {};
a.b = function() {};
var c = null;
a.b = 0;
a.b == a.b;
`)
}

// Java: AggressiveInlineAliasesTest.testAliasCreatedForCtorDepth1_1()
func (c *AggressiveInlineAliasesTest) TestAliasCreatedForCtorDepth1_1() {
	c.Test(
		`var a = function() {};
a.b = 1;
var c = a;
c.b = 2;
a.b == c.b;
`,
		`var a = function() {};
a.b = 1;
var c = null;
a.b = 2;
a.b == a.b;
`)
	c.Test(
		`var a = function() {};
a.b = 1;
var c = a;
c.b = 2;
a?.b == c?.b;
`,
		`var a = function() {};
a.b = 1;
var c = null;
a.b = 2;
a?.b == a?.b;
`)
}

// Java: AggressiveInlineAliasesTest.testAliasCreatedForCtorDepth2()
func (c *AggressiveInlineAliasesTest) TestAliasCreatedForCtorDepth2() {
	c.Test(
		`var a = {};
a.b = function() {};
a.b.c = 1;
var d = a.b;
a.b.c == d.c;
`,
		`var a = {};
a.b = function() {};
a.b.c = 1;
var d = null;
a.b.c == a.b.c
`)
	c.TestSame(
		`var a = {};
a.b = function() {};
a.b.c = 1;
var d = a?.b; // optional chain prevents this being a true alias
a.b.c == d?.c;
`)
}

// Java: AggressiveInlineAliasesTest.testAliasCreatedForFunctionDepth2()
func (c *AggressiveInlineAliasesTest) TestAliasCreatedForFunctionDepth2() {
	c.Test(
		`var a = {};
a.b = function() {};
a.b.c = 1;
var d = a.b;
a.b.c != d.c;
`,
		`var a = {};
a.b = function() {};
a.b.c = 1;
var d = null;
a.b.c != a.b.c
`)
	c.Test(
		`var a = {};
a.b = function() {};
a.b.c = 1;
var d = a.b;
a.b.c == d.c;
`,
		`var a = {};
a.b = function() {};
a.b.c = 1;
var d = null;
a.b.c == a.b.c
`)
}

// Java: AggressiveInlineAliasesTest.testAliasCreatedForObjectDepth1_1()
func (c *AggressiveInlineAliasesTest) TestAliasCreatedForObjectDepth1_1() {
	c.Test(
		`var a = { b: 0 }
var c = a;
c.b = 1;
a.b == c.b;
`,
		`var a = { b: 0 }
var c = null;
a.b = 1;
a.b == a.b;
`)
	c.Test(
		`var a = { b: 0 }
var c = a;
c.b = 1;
a.b == c.b;
use(c);
`,
		`var a = { b: 0 }
var c = null;
a.b = 1;
a.b == a.b;
use(a);
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAliasCreatedAfterVarDeclaration1()
func (c *AggressiveInlineAliasesTest) TestLocalAliasCreatedAfterVarDeclaration1() {
	c.Test(
		`var a = { b : 3 };
function f() {
  var tmp;
  if (true) {
    tmp = a;
    use(tmp);
  }
}
`,
		`var a = { b : 3 };
function f() {
  var tmp;
  if (true) {
    tmp = null;
    use(a);
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAliasCreatedAfterVarDeclaration2()
func (c *AggressiveInlineAliasesTest) TestLocalAliasCreatedAfterVarDeclaration2() {
	c.Test(
		`var a = { b : 3 };
function f() {
  var tmp;
  if (true) {
    tmp = a;
    use(tmp);
  }
}
`,
		`var a = { b : 3 };
function f() {
  var tmp;
  if (true) {
    tmp = null;
    use(a);
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAliasCreatedAfterVarDeclaration3()
func (c *AggressiveInlineAliasesTest) TestLocalAliasCreatedAfterVarDeclaration3() {
	c.TestSame(
		`var a = { b : 3 };
function f() {
  var tmp;
  if (true) {
    tmp = a;
  }
  use(tmp);
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalCtorAliasCreatedAfterVarDeclaration1()
func (c *AggressiveInlineAliasesTest) TestLocalCtorAliasCreatedAfterVarDeclaration1() {
	c.Test(
		`/** @constructor @struct */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  var tmp;
  tmp = Main;
  tmp.doSomething(5);
}
`,
		`/** @constructor */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  var tmp;
  tmp = null;
  Main.doSomething(5);
}
`)
	c.Test(
		`/** @constructor @struct */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  var tmp;
  tmp = Main;
  tmp?.doSomething(5);
}
`,
		`/** @constructor */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  var tmp;
  tmp = null;
  Main?.doSomething(5);
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalCtorAliasCreatedAfterVarDeclaration2()
func (c *AggressiveInlineAliasesTest) TestLocalCtorAliasCreatedAfterVarDeclaration2() {
	c.Test(
		`/** @constructor @struct */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  tmp = Main;
  tmp.doSomething(5);
  if (true) {
    tmp.doSomething(6);
  }
  var tmp;
}
`,
		`/** @constructor */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  tmp = null;
  Main.doSomething(5);
  if (true) {
    Main.doSomething(6);
  }
  var tmp;
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalCtorAliasCreatedAfterVarDeclaration3()
func (c *AggressiveInlineAliasesTest) TestLocalCtorAliasCreatedAfterVarDeclaration3() {
	c.Test(
		`/** @constructor @struct */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  if (tmp = Main) {
    tmp.doSomething(6);
  }
  var tmp;
}
`,
		`/** @constructor */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  if (tmp = Main) { // Don't set "tmp = null" because that would change control flow.
    Main.doSomething(6);
  }
  var tmp;
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalCtorAliasCreatedAfterVarDeclaration4()
func (c *AggressiveInlineAliasesTest) TestLocalCtorAliasCreatedAfterVarDeclaration4() {
	c.Test(
		`/** @constructor @struct */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  var a = tmp = Main;
  tmp.doSomething(6);
  a.doSomething(7);
  var tmp;
}
`,
		`/** @constructor */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  var a = tmp = Main; // Don't set "tmp = null" because that would mess up a's value.
  Main.doSomething(6);
  a.doSomething(7); // Main doesn't get inlined here, which makes collapsing unsafe.
  var tmp;
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalCtorAliasAssignedInLoop1()
func (c *AggressiveInlineAliasesTest) TestLocalCtorAliasAssignedInLoop1() {
	c.TestWarning(
		`/** @constructor @struct */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  var tmp;
  for (let i = 0; i < n(); i++) {
    tmp = Main;
    tmp.doSomething(5);
    use(tmp);
  }
  use(tmp);
  use(tmp.doSomething);
}
`,
		InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING)
}

// Java: AggressiveInlineAliasesTest.testLocalCtorAliasAssignedInLoop2()
func (c *AggressiveInlineAliasesTest) TestLocalCtorAliasAssignedInLoop2() {
	// Test when the alias is assigned in a loop after being used.
	c.TestWarning(
		`/** @constructor @struct */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  var tmp;
  for (let i = 0; i < n(); i++) {
    use(tmp); // Don't inline tmp here since it changes between loop iterations.
    tmp = Main;
  }
  use(tmp);
  use(tmp.doSomething);
}
`,
		InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING)
}

// Java: AggressiveInlineAliasesTest.testLocalCtorAliasAssignedInSwitchCase()
func (c *AggressiveInlineAliasesTest) TestLocalCtorAliasAssignedInSwitchCase() {
	// This mimics how the async generator polyfill behaves.
	c.Test(
		`/** @constructor @struct */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  function g() {
    for (;;) {
      switch(state) {
        case 0:
          tmp1 = Main;
          tmp2 = tmp1.doSomething;
          state = 1;
          return;
         case 1:
           return tmp2.call(tmp1, 3);
      }
    }
  }
  var tmp1;
  var tmp2;
  var state = 0;
  g();
  return g();
}
`,
		`/** @constructor */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  function g() {
    for (;;) {
      switch(state) {
        case 0:
          tmp1 = Main;
          tmp2 = Main.doSomething;
          state = 1;
          return;
         case 1:
           return tmp2.call(tmp1, 3);
      }
    }
  }
  var tmp1;
  var tmp2;
  var state = 0;
  g();
  return g();
}
`)
}

// Java: AggressiveInlineAliasesTest.testCodeGeneratedByGoogModule()
func (c *AggressiveInlineAliasesTest) TestCodeGeneratedByGoogModule() {
	c.Test(
		`var $jscomp = {}
$jscomp.scope = {};
$jscomp.scope.Foo = function() {};
var exports = $jscomp.scope.Foo;
exports.staticprop = { A: 1 };
var y = exports.staticprop.A;
`,
		`var $jscomp = {}
$jscomp.scope = {};
$jscomp.scope.Foo = function() {};
var exports = null;
$jscomp.scope.Foo.staticprop = { A: 1 };
var y = null;
`)
	c.Test(
		`var $jscomp = {}
$jscomp.scope = {};
$jscomp.scope.Foo = function() {};
$jscomp.scope.Foo.staticprop = { A: 1 };
var exports = $jscomp.scope.Foo;
var y = exports.staticprop.A;
`,
		`var $jscomp = {}
$jscomp.scope = {};
$jscomp.scope.Foo = function() {};
$jscomp.scope.Foo.staticprop = { A: 1 };
var exports = null;
var y = null;
`)
}

// Java: AggressiveInlineAliasesTest.testCollapse()
func (c *AggressiveInlineAliasesTest) TestCollapse() {
	c.Test(
		`var a = {}
a.b = {};
var c = a.b;
`,
		`var a = {}
a.b = {};
var c = null
`)
	c.Test(
		`var a = {}
a.b = {};
var c = a.b;
use(c);
`,
		`var a = {}
a.b = {};
var c = null;
use(a.b)
`)
	c.TestSame(
		`var a = {};
a.b;
`)
}

// Java: AggressiveInlineAliasesTest.testCollapsePropertiesOfClass1()
func (c *AggressiveInlineAliasesTest) TestCollapsePropertiesOfClass1() {
	c.Test(
		`/** @constructor */ var namespace = function() {};
goog.inherits(namespace, Object);
namespace.includeExtraParam = true;
/** @enum { number } */
namespace.Param = { param1: 1, param2: 2 };
if (namespace.includeExtraParam) namespace.Param.optParam = 3;
/** @constructor */
function f() {
  var Param = namespace.Param;
  log(namespace.Param.optParam);
  log(Param.optParam);
}
`,
		`/** @constructor */ var namespace = function() {};
goog.inherits(namespace,Object);
namespace.includeExtraParam = true;
/** @enum {JSDocSerializer_placeholder_type} */
namespace.Param = { param1: 1,param2: 2 };
if(namespace.includeExtraParam) namespace.Param.optParam = 3;
/** @constructor */
function f() {
  var Param = null;
  log(namespace.Param.optParam);
  log(namespace.Param.optParam);
}
`)
}

// Java: AggressiveInlineAliasesTest.testCollapsePropertiesOfClass2()
func (c *AggressiveInlineAliasesTest) TestCollapsePropertiesOfClass2() {
	c.Test(
		`var goog = goog || {}
goog.addSingletonGetter = function(cls) {};
var a = {};
a.b = function() {};
goog.addSingletonGetter(a.b);
a.b.prototype.get = function(key) {};
a.b.c = function() {};
a.b.c.XXX = new a.b.c;
function f() {
  var x = a.b.getInstance();
  var Key = a.b.c;
  x.get(Key.XXX);
}
`,
		`var goog = goog || {}
goog.addSingletonGetter = function(cls) {};
var a = {};
a.b = function() {};
goog.addSingletonGetter(a.b);
a.b.prototype.get = function(key) {};
a.b.c = function() {};
a.b.c.XXX = new a.b.c;
function f() {
  var x = a.b.getInstance();
  var Key = null;
  x.get(a.b.c.XXX);
}
`)
}

// Java: AggressiveInlineAliasesTest.testCommaOperator()
func (c *AggressiveInlineAliasesTest) TestCommaOperator() {
	c.Test(
		`var ns = {}
ns.Foo = {};
var Baz = {};
Baz.Foo = ns.Foo;
Baz.Foo.bar = 10, 123;
`,
		`var ns = {}
ns.Foo = {};
var Baz = {};
Baz.Foo = null;
ns.Foo.bar = 10, 123;
`)
}

// Java: AggressiveInlineAliasesTest.testDontCrashCtorAliasWithEnum()
func (c *AggressiveInlineAliasesTest) TestDontCrashCtorAliasWithEnum() {
	c.Test(
		`var ns = {};
/** @constructor */ ns.Foo = function () {};
var Bar = ns.Foo;
/** @const @enum */
Bar.prop = { A: 1 };
`,
		`var ns = {};
/** @constructor */ ns.Foo = function() {};
var Bar = null;
/** @const @enum {JSDocSerializer_placeholder_type} */
ns.Foo.prop = { A: 1 }
`)
}

// Java: AggressiveInlineAliasesTest.testFunctionAlias2()
func (c *AggressiveInlineAliasesTest) TestFunctionAlias2() {
	c.Test(
		`var a = {}
a.b = {};
a.b.c = function() {};
a.b.d = a.b.c;
use(a.b.d);
`,
		`var a = {}
a.b = {};
a.b.c = function() {};
a.b.d = null;
use(a.b.c);
`)
	c.TestSame(
		`var a = {}
a.b = {};
a.b.c = function() {};
a.b.d = a.b.c;
// reference to alias via optional chain stops inlining
use(a.b?.d);
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalAliasOfAncestor()
func (c *AggressiveInlineAliasesTest) TestGlobalAliasOfAncestor() {
	c.Test(
		`var a = { b: { c: 5 } }
function f() {
  var x = a.b;
  f(x.c);
}
`,
		`var a = { b: { c: 5 } }
function f() {
  var x = null;
  f(a.b.c);
}
`)
	c.Test(
		`var a = { b: { c: 5 } }
function f() {
  var x = a.b;
  f(x?.c);
}
`,
		`var a = { b: { c: 5 } }
function f() {
  var x = null;
  f(a.b?.c);
}
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalES5ClassVarWithInnerNameDotPropReference()
func (c *AggressiveInlineAliasesTest) TestGlobalES5ClassVarWithInnerNameDotPropReference() {
	c.Test(
		`/** @constructor */
var GlobalName = function InnerName() {
    InnerName.staticMethod(); // replace InnerName with GlobalName
};
GlobalName.staticMethod = function() {
  console.log('staticMethod');
};
`,
		`/** @constructor */
var GlobalName = function InnerName() {
    GlobalName.staticMethod();
};
GlobalName.staticMethod = function() {
  console.log('staticMethod');
};
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalClassVarWithInnerNameDotPropReference()
func (c *AggressiveInlineAliasesTest) TestGlobalClassVarWithInnerNameDotPropReference() {
	c.Test(
		`const GlobalName = class InnerName {
  method() {
    return InnerName.staticMethod(); // replace InnerName with GlobalName
  }
  static staticMethod() {
    console.log('staticMethod');
  }
};
`,
		`const GlobalName = class {
  method() {
    return GlobalName.staticMethod();
  }
  static staticMethod() {
    console.log("staticMethod");
  }
};
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalClassPropWithInnerNameDotPropReference()
func (c *AggressiveInlineAliasesTest) TestGlobalClassPropWithInnerNameDotPropReference() {
	c.Test(
		`const GlobalName = {};
GlobalName.prop = class InnerName {
  method() {
    return InnerName.staticMethod(); // replace InnerName with GlobalName.prop
  }
  static staticMethod() {
    console.log('staticMethod');
  }
};
`,
		`const GlobalName = {};
GlobalName.prop = class {
  method() {
    return GlobalName.prop.staticMethod();
  }
  static staticMethod() {
    console.log("staticMethod");
  }
};
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalClassVarSetTwiceWithInnerNameDotPropReference()
func (c *AggressiveInlineAliasesTest) TestGlobalClassVarSetTwiceWithInnerNameDotPropReference() {
	c.Test(
		`let GlobalName = class InnerName {
  method() {
    return InnerName.staticMethod(); // replace InnerName with GlobalName
  }
  static staticMethod() {
    console.log('staticMethod');
  }
};
// second assignment prevents inlining
GlobalName = function SecondValue() {}
`,
		`let GlobalName = class {
  method() {
    return GlobalName.staticMethod();
  }
  static staticMethod() {
    console.log("staticMethod");
  }
};
GlobalName = function SecondValue() {};
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalClassPropWithInnerNameInstanceOfReference()
func (c *AggressiveInlineAliasesTest) TestGlobalClassPropWithInnerNameInstanceOfReference() {
	c.Test(
		`const GlobalName = {};
GlobalName.prop = class InnerName {
  method() {
    return this instanceof InnerName;
  }
};
`,
		`const GlobalName = {};
GlobalName.prop = class {
  method() {
    return this instanceof GlobalName.prop;
  }
};
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalAliasWithProperties1()
func (c *AggressiveInlineAliasesTest) TestGlobalAliasWithProperties1() {
	c.Test(
		`var ns = {}
ns.Foo = function() {};
/** @enum { number } */ ns.Foo.EventType = { A: 1, B: 2 };
ns.Bar = ns.Foo;
var x = function() { use(ns.Bar.EventType.A) };
use(x);
`,
		`var ns = {}
ns.Foo = function() {};
/** @enum {JSDocSerializer_placeholder_type} */ ns.Foo.EventType = { A: 1, B: 2 };
ns.Bar = null;
var x = function() { use(ns.Foo.EventType.A) };
use(x);
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalAliasWithProperties2()
func (c *AggressiveInlineAliasesTest) TestGlobalAliasWithProperties2() {
	c.Test(
		`var ns = {}
ns.Foo = function() {};
/** @enum { number } */ ns.Foo.EventType = { A: 1, B: 2 };
ns.Bar = ns.Foo;
/** @const */ ns.Bar.EventType = ns.Foo.EventType;
var x = function() { use(ns.Bar.EventType.A) };
use(x)
`,
		`var ns = {}
ns.Foo = function() {};
/** @enum {JSDocSerializer_placeholder_type} */ ns.Foo.EventType = { A: 1, B: 2 };
ns.Bar = null;
/** @const */ ns.Foo.EventType = ns.Foo.EventType;
var x = function() { use(ns.Foo.EventType.A) };
use(x)
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalAliasWithProperties3()
func (c *AggressiveInlineAliasesTest) TestGlobalAliasWithProperties3() {
	c.Test(
		`var ns = {}
ns.Foo = function() {};
/** @enum { number } */ ns.Foo.EventType = { A: 1, B: 2 };
ns.Bar = ns.Foo;
/** @const */ ns.Bar.Other = { X: 1, Y: 2 };
var x = function() { use(ns.Bar.Other.X) };
use(x)
`,
		`var ns = {}
ns.Foo = function() {};
/** @enum {JSDocSerializer_placeholder_type} */ ns.Foo.EventType = { A: 1, B: 2 };
ns.Bar = null;
/** @const */ ns.Foo.Other = { X: 1, Y: 2 };
var x = function() { use(ns.Foo.Other.X) };
use(x)
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalAliasWithPropertiesAndOptChainReference()
func (c *AggressiveInlineAliasesTest) TestGlobalAliasWithPropertiesAndOptChainReference() {
	c.Test(
		`var ns = {}
ns.Foo = function() {};
/** @enum { number } */ ns.Foo.EventType = { A: 1, B: 2 };
ns.Bar = ns.Foo;
/** @const */ ns.Bar.Other = { X: 1, Y: 2 };
var x = function() { use(ns.Bar?.Other.X) };
use(x)
`,
		`var ns = {}
ns.Foo = function() {};
/** @enum {JSDocSerializer_placeholder_type} */ ns.Foo.EventType = { A: 1, B: 2 };
ns.Bar = null;
/** @const */ ns.Foo.Other = { X: 1, Y: 2 };
var x = function() { use(ns.Foo?.Other.X) };
use(x)
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalAliasWithPropertiesAsNestedObjectLits()
func (c *AggressiveInlineAliasesTest) TestGlobalAliasWithPropertiesAsNestedObjectLits() {
	c.Test(
		`var ns = {}
ns.Foo = function() {};
ns.Bar = ns.Foo;
/** @enum { number } */ ns.Bar.Other = { X: {Y: 1}};
var x = function() { use(ns.Bar.Other.X.Y) };
use(x)
`,
		`var ns = {}
ns.Foo = function() {};
ns.Bar = null;
/** @enum {JSDocSerializer_placeholder_type} */ ns.Foo.Other = { X: {Y: 1}};
var x = function() { use(ns.Foo.Other.X.Y) };
use(x)
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalWriteToNonAncestor()
func (c *AggressiveInlineAliasesTest) TestGlobalWriteToNonAncestor() {
	c.Test(
		`var a = { b: 3 }
function f() {
  var x = a;
  f(a.b);
}
a.b = 5;
`,
		`var a = { b: 3 }
function f() {
  var x = null;
  f(a.b);
}
a.b = 5;
`)
}

// Java: AggressiveInlineAliasesTest.testInlineCtorInObjLit()
func (c *AggressiveInlineAliasesTest) TestInlineCtorInObjLit() {
	c.Test(
		`function Foo() {}
var Bar = Foo;
var objlit = { 'prop' : Bar };
`,
		`function Foo() {}
var Bar = null;
var objlit = { 'prop': Foo };
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAlias1()
func (c *AggressiveInlineAliasesTest) TestLocalAlias1() {
	c.Test(
		`var a = { b: 3 }
function f() {
  var x = a;
  f(x.b);
}
`,
		`var a = { b: 3 }
function f() {
  var x = null;
  f(a.b);
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAlias2()
func (c *AggressiveInlineAliasesTest) TestLocalAlias2() {
	c.Test(
		`var a = { b: 3, c: 4 }
function f() { var x = a;
f(x.b);
f(x.c);
}
`,
		`var a = { b: 3, c: 4 }
function f() { var x = null;
f(a.b);
f(a.c);
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAlias3()
func (c *AggressiveInlineAliasesTest) TestLocalAlias3() {
	c.Test(
		`var a = { b: 3, c: { d: 5 } }
function f() {
  var x = a;
  f(x.b);
  f(x.c);
  f(x.c?.d);
}
`,
		`var a = { b: 3, c: { d: 5 } }
function f() {
  var x = null;
  f(a.b);
  f(a.c);
  f(a.c?.d)
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAlias4()
func (c *AggressiveInlineAliasesTest) TestLocalAlias4() {
	c.Test(
		`var a = { b: 3 }
var c = { d: 5 };
function f() {
  var x = a;
  var y = c;
  f(x.b);
  f(y.d);
}
`,
		`var a = { b: 3 }
var c = { d: 5 };
function f() {
  var x = null;
  var y = null;
  f(a.b);
  f(c.d)
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAlias5()
func (c *AggressiveInlineAliasesTest) TestLocalAlias5() {
	c.Test(
		`var a = { b: { c: 5 } }
function f() {
  var x = a;
  var y = x.b;
  f(a.b.c);
  f(y.c);
}
`,
		`var a = { b: { c: 5 } }
function f() {
  var x = null;
  var y = null;
  f(a.b.c);
  f(a.b.c);
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAlias6()
func (c *AggressiveInlineAliasesTest) TestLocalAlias6() {
	c.Test(
		`var a = { b: 3 }
function f() {
  var x = a;
  if (x.b) f(x.b);
}
`,
		`var a = { b: 3 }
function f() {
  var x = null;
  if (a.b) f(a.b);
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAlias7()
func (c *AggressiveInlineAliasesTest) TestLocalAlias7() {
	c.Test(
		`var a = { b: { c: 5 } }
function f() {
  var x = a.b;
  f(x.c);
}
`,
		`var a = { b: { c: 5 } }
function f() {
  var x = null;
  f(a.b.c);
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAlias8()
func (c *AggressiveInlineAliasesTest) TestLocalAlias8() {
	c.TestSame(
		`var a = { b: 3 };
function f() { if (true) { var x = a; f(x.b); } x = { b : 4}; }
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAliasOfEnumWithInstanceofCheck()
func (c *AggressiveInlineAliasesTest) TestLocalAliasOfEnumWithInstanceofCheck() {
	c.Test(
		`/** @constructor */ var Enums = function() {};
/** @enum { number } */
Enums.Fruit = { APPLE: 1, BANANA: 2 };
/** @constructor */ function foo(f) {
if (f instanceof Enums) { alert('what?'); return; }
var Fruit = Enums.Fruit;
if (f == Fruit.APPLE) alert('apple');
if (f == Fruit.BANANA) alert('banana'); }
`,
		`/** @constructor */ var Enums = function() {};
/** @enum {JSDocSerializer_placeholder_type} */
Enums.Fruit = { APPLE: 1,BANANA: 2 };
/** @constructor */ function foo(f) {
if (f instanceof Enums) { alert('what?'); return; }
var Fruit = null;
if(f == Enums.Fruit.APPLE) alert('apple');
if(f == Enums.Fruit.BANANA) alert('banana'); }
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAliasOfFunction()
func (c *AggressiveInlineAliasesTest) TestLocalAliasOfFunction() {
	c.Test(
		`var a = function() {}
a.b = 5;
function f() { var x = a.b;
f(x);
}
`,
		`var a = function() {}
a.b = 5;
function f() { var x = null;
f(a.b);
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalWriteToNonAncestor()
func (c *AggressiveInlineAliasesTest) TestLocalWriteToNonAncestor() {
	c.Test(
		`var a = { b: 3 }
function f() {
  a.b = 5;
  var x = a;
  f(a.b);
}
`,
		`var a = { b: 3 }
function f() {
  a.b = 5;
  var x = null;
  f(a.b);
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAliasInChainedAssignment()
func (c *AggressiveInlineAliasesTest) TestLocalAliasInChainedAssignment() {
	c.TestSame("var a = { b: 3 }; function f() { var c; var d = c = a; a.b; d.b; }")
}

// Java: AggressiveInlineAliasesTest.testMisusedConstructorTag()
func (c *AggressiveInlineAliasesTest) TestMisusedConstructorTag() {
	c.Test(
		`var a = {}
var d = a;
a.b = function() {};
a.b.c = 0;
a.b.c;
`,
		`var a = {}
var d = null;
a.b = function() {};
a.b.c = 0;
a.b.c;
`)
}

// Java: AggressiveInlineAliasesTest.testNoInlineGetpropIntoCall()
func (c *AggressiveInlineAliasesTest) TestNoInlineGetpropIntoCall() {
	c.Test(
		`var b = x
function f() {
  var a = b;
  a();
}
`,
		`var b = x
function f() {
  var a = null;
  b()
}
`)
}

// Java: AggressiveInlineAliasesTest.testObjLitAssignmentDepth3()
func (c *AggressiveInlineAliasesTest) TestObjLitAssignmentDepth3() {
	c.Test(
		`var a = {}
a.b = {};
a.b.c = { d: 1, e: 2 };
var f = a.b.c.d;
`,
		`var a = {}
a.b = {};
a.b.c = { d: 1,e: 2 };
var f = null;
`)
	c.Test(
		`var a = {}
a.b = {};
a.b.c = { d: 1, e: 2 };
var f = a.b.c.d;
use(f);
`,
		`var a = {}
a.b = {};
a.b.c = { d: 1,e: 2 };
var f = null;
use(a.b.c.d);
`)
	c.Test(
		`var a = {}
a.b = {};
a.b.c = { d: 1, e: 2 };
var f = a.b.c.d;
var g = a.b.c.e;
`,
		`var a = {}
a.b = {};
a.b.c = { d: 1, e: 2 };
var f = null;
var g = null;
`)
	c.Test(
		`var a = {}
a.b = {};
a.b.c = { d: 1, e: 2 };
var f = a.b.c.d;
var g = a.b.c.e;
use(f, g);
`,
		`var a = {}
a.b = {};
a.b.c = { d: 1, e: 2 };
var f = null;
var g = null;
use(a.b.c.d, a.b.c.e);
`)
	c.TestSame(
		`var a = {}
a.b = {};
a.b.c = { d: 1, e: 2 };
var f = null;
var g = null;
`)
}

// Java: AggressiveInlineAliasesTest.testObjLitAssignmentDepth4()
func (c *AggressiveInlineAliasesTest) TestObjLitAssignmentDepth4() {
	c.Test(
		`var a = {}
a.b = {};
a.b.c = {};
a.b.c.d = { e: 1, f: 2 };
var g = a.b.c.d.e;
`,
		`var a = {}
a.b = {};
a.b.c = {};
a.b.c.d = { e: 1, f: 2 };
var g = null;
`)
	c.Test(
		`var a = {}
a.b = {};
a.b.c = {};
a.b.c.d = { e: 1, f: 2 };
var g = a.b.c.d.e;
use(g);
`,
		`var a = {}
a.b = {};
a.b.c = {};
a.b.c.d = { e: 1, f: 2 };
var g = null;
use(a.b.c.d.e);
`)
}

// Java: AggressiveInlineAliasesTest.testObjLitDeclaration()
func (c *AggressiveInlineAliasesTest) TestObjLitDeclaration() {
	c.Test(
		`var a = { b: {}, c: {} }
var d = a.b;
var e = a.c;
`,
		`var a = { b: {}, c: {} }
var d = null;
var e = null;
`)
	c.Test(
		`var a = { b: {}, c: {} }
var d = a.b;
var e = a.c;
use(d, e);
`,
		`var a = { b: {}, c: {} }
var d = null;
var e = null;
use(a.b,a.c);
`)
	c.Test(
		`var a = { b: {}, c: {} }
var d = a.b;
var e = a.c;
`,
		`var a = { b: {}, c: {} }
var d = null;
var e = null;
`)
	c.Test(
		`var a = { b: {}, c: {} }
var d = a.b;
var e = a.c;
use(d, e);
`,
		`var a = { b: {}, c: {} }
var d = null;
var e = null;
use(a.b, a.c);
`)
}

// Java: AggressiveInlineAliasesTest.testObjLitDeclarationUsedInSameVarList()
func (c *AggressiveInlineAliasesTest) TestObjLitDeclarationUsedInSameVarList() {
	c.Test(
		`var a = { b: {}, c: {} }
var d = a.b;
var e = a.c;
use(d, e);
`,
		`var a = { b: {}, c: {} }
var d = null;
var e = null;
use(a.b, a.c)
`)
}

// Java: AggressiveInlineAliasesTest.testObjLitDeclarationWithGet2()
func (c *AggressiveInlineAliasesTest) TestObjLitDeclarationWithGet2() {
	c.Test(
		`var a = { b: {}, get c() {} }
var d = a.b;
var e = a.c;
`,
		`var a = { b: {}, get c() {} }
var d = null;
var e = a.c;
`)
	c.Test(
		`var a = { b: {}, get c() {} }
var d = a.b;
var e = a.c;
use(d);
`,
		`var a = { b: {}, get c() {} }
var d = null;
var e = a.c;
use(a.b);
`)
}

// Java: AggressiveInlineAliasesTest.testObjLitDeclarationWithSet2()
func (c *AggressiveInlineAliasesTest) TestObjLitDeclarationWithSet2() {
	c.Test(
		`var a = { b: {}, set c(a) {} }
var d = a.b;
var e = a.c
`,
		`var a = { b: {}, set c(a$jscomp$1) {} }
var d = null;
var e = a.c;
`)
	c.Test(
		`var a = { b: {}, set c(a) {} }
var d = a.b;
var e = a.c;
use(d);
`,
		`var a = { b: {}, set c(a$jscomp$1) {} }
var d = null;
var e = a.c;
use(a.b);
`)
}

// Java: AggressiveInlineAliasesTest.testPropertyOfChildFuncOfUncollapsibleObjectDepth1()
func (c *AggressiveInlineAliasesTest) TestPropertyOfChildFuncOfUncollapsibleObjectDepth1() {
	c.Test(
		`var a = {}
var c = a;
a.b = function () {};
a.b.x = 0;
a.b.x;
`,
		`var a = {}
var c = null;
a.b = function() {};
a.b.x = 0;
a.b.x;
`)
	c.Test(
		`var a = {}
var c = a;
a.b = function () {};
a.b.x = 0;
a.b.x;
use(c);
`,
		`var a = {}
var c = null;
a.b = function() {};
a.b.x = 0;
a.b.x;
use(a);
`)
}

// Java: AggressiveInlineAliasesTest.testPropertyOfChildFuncOfUncollapsibleObjectDepth2()
func (c *AggressiveInlineAliasesTest) TestPropertyOfChildFuncOfUncollapsibleObjectDepth2() {
	c.Test(
		`var a = {}
a.b = {};
var c = a.b;
a.b.c = function () {};
a.b.c.x = 0;
a.b.c.x;
`,
		`var a = {}
a.b = {};
var c = null;
a.b.c = function() {};
a.b.c.x = 0;
a.b.c.x
`)
	c.Test(
		`var a = {}
a.b = {};
var c = a.b;
a.b.c = function () {};
a.b.c.x = 0;
a.b.c.x;
use(c);
`,
		`var a = {}
a.b = {};
var c = null;
a.b.c = function() {};
a.b.c.x = 0;
a.b.c.x;
use(a.b)
`)
}

// Java: AggressiveInlineAliasesTest.testTypeDefAlias1()
func (c *AggressiveInlineAliasesTest) TestTypeDefAlias1() {
	c.Test(
		`var D = function() {}
D.L = function() {};
/** @type { D.L } */ D.L.A = new D.L();
/** @const */ var M = {};
/** @typedef { D.L } */ M.L = D.L;
use(M.L.A);
`,
		`var D = function() {}
D.L = function() {};
 D.L.A = new D.L;
/** @const */ var M = {};
M.L = null;
use(D.L.A);
`)
}

// Java: AggressiveInlineAliasesTest.testGlobalAliasWithConst()
func (c *AggressiveInlineAliasesTest) TestGlobalAliasWithConst() {
	c.Test("const a = {}; a.b = {}; const c = a.b;", "const a = {}; a.b = {}; const c = null")
	c.Test(
		"const a = {}; a.b = {}; const c = a.b; use(c);",
		"const a = {}; a.b = {}; const c = null; use(a.b)")
	c.TestSame("const a = {}; a.b;")
}

// Java: AggressiveInlineAliasesTest.testGlobalAliasWithLet1()
func (c *AggressiveInlineAliasesTest) TestGlobalAliasWithLet1() {
	c.Test("let a = {}; a.b = {}; let c = a.b;", "let a = {}; a.b = {}; let c = null")
	c.Test(
		"let a = {}; a.b = {}; let c = a.b; use(c);",
		"let a = {}; a.b = {}; let c = null; use(a.b)")
	c.TestSame("let a = {}; a.b;")
	c.Test(
		"let a = {}; if (true) { a.b = 5; } let c = a.b; use(c);",
		"let a = {}; if (true) { a.b = 5; } let c = null; use(a.b);")
}

// Java: AggressiveInlineAliasesTest.testGlobalAliasWithLet2()
func (c *AggressiveInlineAliasesTest) TestGlobalAliasWithLet2() {
	// Inlining ns is unsafe because ns2 may or may not be equal to ns.
	c.TestSame(
		`let ns = {};
ns.foo = 'bar';
let ns2;
if (true) {
  ns2 = ns;
}
use(ns2.foo);
`)
	// In this case, it would be safe to inline ns, as long as all subsequent references
	// to ns2 are inside the if block, but the algorithm is not complex enough to know that.
	c.TestSame(
		`let ns = {};
ns.foo = 'bar';
let ns2;
if (true) {
  ns2 = ns;
  use(ns2.foo);
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAliasWithLet1()
func (c *AggressiveInlineAliasesTest) TestLocalAliasWithLet1() {
	c.Test(
		`var a = {};
a.b = {};
function f() {
  if (true) {
    let c = a.b;
    alert(c);
  }
}
`,
		`var a = {};
a.b = {};
function f() {
  if (true) {
    let c = null;
    alert(a.b);
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.testLocalAliasWithLet2()
func (c *AggressiveInlineAliasesTest) TestLocalAliasWithLet2() {
	c.Test(
		"var a = {}; a.b = {}; if (true) { let c = a.b;  alert(c); }",
		"var a = {}; a.b = {}; if (true) { let c = null;  alert(a.b); }")
}

// Java: AggressiveInlineAliasesTest.testLocalAliasWithLet3()
func (c *AggressiveInlineAliasesTest) TestLocalAliasWithLet3() {
	c.Test(
		"let ns = {a: 1}; { let y = ns; use(y.a); }",
		"let ns = {a: 1}; { let y = null; use(ns.a); }")
}

// Java: AggressiveInlineAliasesTest.testLocalAliasInsideClass()
func (c *AggressiveInlineAliasesTest) TestLocalAliasInsideClass() {
	c.Test(
		"var a = {x: 5}; class A { fn() { var b = a; use(b.x); } }",
		"var a = {x: 5}; class A { fn() { var b = null; use(a.x); } }")
}

// Java: AggressiveInlineAliasesTest.testGlobalClassAlias1()
func (c *AggressiveInlineAliasesTest) TestGlobalClassAlias1() {
	c.Test(
		"class A {} A.foo = 5; const B = A; use(B.foo);",
		"class A {} A.foo = 5; const B = null; use(A.foo)")
}

// Java: AggressiveInlineAliasesTest.testGlobalClassAlias2()
func (c *AggressiveInlineAliasesTest) TestGlobalClassAlias2() {
	c.Test(
		"class A { static fn() {} } const B = A; B.fn();",
		"class A { static fn() {} } const B = null; A.fn();")
}

// Java: AggressiveInlineAliasesTest.testGlobalClassAlias3()
func (c *AggressiveInlineAliasesTest) TestGlobalClassAlias3() {
	c.Test(
		"class A {} const B =    A; B.prototype.fn = () =>          5   ;",
		"class A {} const B = null; A.prototype.fn = () => { return 5; };")
}

// Java: AggressiveInlineAliasesTest.testDestructuringAlias1()
func (c *AggressiveInlineAliasesTest) TestDestructuringAlias1() {
	// CollapseProperties backs off on destructuring, so it's okay not to inline here.
	c.Test(
		"var a = {x: 5}; var [b] = [a]; use(b.x);", //
		"var a = {x: 5}; var b; [b] = [a]; use(b.x);")
}

// Java: AggressiveInlineAliasesTest.testDestructuringAlias2()
func (c *AggressiveInlineAliasesTest) TestDestructuringAlias2() {
	c.TestSame("var a = {x: 5}; const {b} = {b: a}; use(b.x);")
}

// Java: AggressiveInlineAliasesTest.testDestructuringArrayAlias()
func (c *AggressiveInlineAliasesTest) TestDestructuringArrayAlias() {
	c.TestSame("var a = [5, 6]; const [b, c] = a; use(b);")
}

// Java: AggressiveInlineAliasesTest.testObjectRest_restingFromANamespace_isNotInlinable()
func (c *AggressiveInlineAliasesTest) TestObjectRest_restingFromANamespace_isNotInlinable() {
	// TODO(nickreid): These might actually be inlinable.
	c.TestSame("var a = {x: 5, y: 6}; const {...b} = a; use(b.y);")
	c.TestSame("var a = {x: 5, y: 6}; const {x, ...b} = a; use(b.y);")
}

// Java: AggressiveInlineAliasesTest.testObjectSpread_spreadingInNamespaceDef_preventsInliningItsProps()
func (c *AggressiveInlineAliasesTest) TestObjectSpread_spreadingInNamespaceDef_preventsInliningItsProps() {
	c.TestSame("var a = {x: 5, y: 6}; var b = {...a}; use(b.z);")
	c.TestSame("var a = {x: 5, y: 6, z: 7}; var b = {z: -7, ...a}; use(b.z);")
	c.TestSame("var a = {x: 5, y: 6, z: 7}; var b = Object.assign({}, {z: -7}, a); use(b.z);")
	c.TestSame("var a = {x: 5, y: 6, z: 7}; var b = {...a, z: -7}; use(b.z);")
	c.TestSame("var a = {x: 5, y: 6, z: 7}; var b = Object.assign({}, a, {z: -7}); use(b.z);")
}

// Java: AggressiveInlineAliasesTest.testDefaultParamAlias()
func (c *AggressiveInlineAliasesTest) TestDefaultParamAlias() {
	c.Test(
		"var a = {b: 5}; var b = a; function f(x=b) { alert(x.b); }",
		"var a = {b: 5}; var b = null; function f(x=a) { alert(x.b); }")
}

// Java: AggressiveInlineAliasesTest.testComputedPropertyNames()
func (c *AggressiveInlineAliasesTest) TestComputedPropertyNames() {
	// We don't support computed properties.
	c.TestSame(
		`var foo = {['ba' + 'r']: {}};
var foobar = foo.bar;
foobar.baz = 5;
use(foo.bar.baz);
`)
}

// Java: AggressiveInlineAliasesTest.testAliasInTemplateString()
func (c *AggressiveInlineAliasesTest) TestAliasInTemplateString() {
	c.Test(
		"const a = {b: 5}; const c = a; alert(`${c.b}`);",
		"const a = {b: 5}; const c = null; alert(`${a.b}`);")
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_method()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_method() {
	c.Test(
		"class A { static s() {} } class B extends A {} const C = B; C.s();",
		"class A { static s() {} } class B extends A {} const C = null; A.s();")
	c.Test(
		"class A { static s() {} } class B extends A {} B.s();",
		"class A { static s() {} } class B extends A {} A.s();")
	c.Test(
		"class A {} A.s = function() {}; class B extends A {} B.s();",
		"class A {} A.s = function() {}; class B extends A {} A.s();")
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_methodsFromMultipleClasses()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_methodsFromMultipleClasses() {
	// C is an alias of B
	c.Test(
		`class A { static a() {} }
class B extends A { static b() {} }
const C = B;
C.a();
C.b()
`,
		`class A { static a() {} }
class B extends A { static b() {} }
const C = null;
A.a();
B.b()
`)
	// C is a subclass of A and B
	c.Test(
		`class A { static a() {} }
class B extends A { static b() {} }
class C extends B {}
C.a();
C.b()
`,
		`class A { static a() {} }
class B extends A { static b() {} }
class C extends B {}
A.a();
B.b()
`)
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_methodWithNoCollapse()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_methodWithNoCollapse() {
	// back off on replacing `B.s` -> `A.s` if A.s is not collapsible for two reasons:
	//  1. the main reason we do this replacing is to make collapsing safer
	//  2. people may use @nocollapse to avoid breaking static `this` refs, and if that were the
	//     case then inlining would also break those refs.
	c.TestSame("class A { /** @nocollapse */ static s() {} } class B extends A {} B.s();")
}

// Java: AggressiveInlineAliasesTest.testChainedClassStaticInheritance_methodWithNoCollapse()
func (c *AggressiveInlineAliasesTest) TestChainedClassStaticInheritance_methodWithNoCollapse() {
	// verify we also don't replace C.s with B.s, and inherit the 'non-collapsibility' from 'A.s'
	c.TestSame(
		`class A {/** @nocollapse */ static s() {}}
class B extends A {}
class C extends B {}
C.s();
`)
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_propertyAlias()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_propertyAlias() {
	c.TestSame("class A {} A.staticProp = 6; class B extends A {} let b = new B;")
	c.Test(
		"class A {} A.staticProp = 6; class B extends A {} use(B.staticProp);",
		"class A {} A.staticProp = 6; class B extends A {} use(A.staticProp);")
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_classExpression()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_classExpression() {
	c.Test(
		"var A = class {}; A.staticProp = 6; var B = class extends A {}; use(B.staticProp);",
		"var A = class {}; A.staticProp = 6; var B = class extends A {}; use(A.staticProp);")
	c.Test(
		"let A = class {}; A.staticProp = 6; let B = class extends A {}; use(B.staticProp);",
		"let A = class {}; A.staticProp = 6; let B = class extends A {}; use(A.staticProp);")
	c.Test(
		"const A = class {}; A.staticProp = 6; const B = class extends A {}; use(B.staticProp);",
		"const A = class {}; A.staticProp = 6; const B = class extends A {}; use(A.staticProp);")
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritence_lateInitializedExpression()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritence_lateInitializedExpression() {
	c.TestSame(
		"var A; A = class {}; A.staticProp = 6; var B = class extends A {}; use(B.staticProp);")
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_propertyWithSubproperty()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_propertyWithSubproperty() {
	c.Test(
		"class A {} A.ns = {foo: 'bar'}; class B extends A {} use(B.ns.foo);",
		"class A {} A.ns = {foo: 'bar'}; class B extends A {} use(A.ns.foo);")
	c.Test(
		"class A {} A.ns = {foo: 'bar'}; class B extends A {} B.ns.foo = 'baz'; use(B.ns.foo);",
		"class A {} A.ns = {foo: 'bar'}; class B extends A {} A.ns.foo = 'baz'; use(A.ns.foo);")
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_propertyWithShadowing()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_propertyWithShadowing() {
	c.TestSame("class A {} A.staticProp = 6; class B extends A {} B.staticProp = 7;")
	// Here, B.staticProp is a different property from A.staticProp, so don't rewrite.
	c.TestSame(
		"class A {} A.staticProp = 6; class B extends A {} B.staticProp = 7; use(B.staticProp);")
	// At the time use() is called, B.staticProp is still the same as A.staticProp, so we
	// *could* rewrite it. But instead we back off because of the shadowing afterwards.
	c.TestSame(
		"class A {} A.staticProp = 6; class B extends A {} use(B.staticProp); B.staticProp = 7;")
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_propertyMultiple()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_propertyMultiple() {
	c.Test(
		"class A {} A.foo = 5; A.bar = 6; class B extends A {} use(B.foo); use(B.bar);",
		"class A {} A.foo = 5; A.bar = 6; class B extends A {} use(A.foo); use(A.bar);")
	c.TestSame("class A {} A.foo = {bar: 5}; A.baz = 6; class B extends A {} B.baz = 7;")
	c.Test(
		"class A {} A.foo = {}; A.baz = 6; class B extends A {}  B.foo.bar = 5; B.baz = 7;",
		"class A {} A.foo = {}; A.baz = 6; class B extends A {} A.foo.bar = 5; B.baz = 7;")
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_property_chainedSubclasses()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_property_chainedSubclasses() {
	c.Test(
		"class A {} A.foo = 5; class B extends A {} class C extends B {} use(C.foo);",
		"class A {} A.foo = 5; class B extends A {} class C extends B {} use(A.foo);")
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_namespacedClass()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_namespacedClass() {
	c.Test(
		`var ns1 = {}, ns2 = {};
ns1.A = class {};
ns1.A.staticProp = {foo: 'bar'};
ns2.B = class extends ns1.A {}
use(ns2.B.staticProp.bar);
`,
		`var ns1 = {};
var ns2 = {};
ns1.A = class {};
ns1.A.staticProp = {foo: 'bar'};
ns2.B = class extends ns1.A {}
use(ns1.A.staticProp.bar);
`)
	c.Test(
		`var ns1 = {}, ns2 = {};
ns1.A = class {};
ns1.A.staticProp = {foo: 'bar'};
ns2.B = class extends ns1.A {}
use(ns2.B.staticProp?.bar);
`,
		`var ns1 = {};
var ns2 = {};
ns1.A = class {};
ns1.A.staticProp = {foo: 'bar'};
ns2.B = class extends ns1.A {}
use(ns1.A.staticProp?.bar);
`)
	c.TestSame(
		`var ns1 = {};
var ns2 = {};
ns1.A = class {};
ns1.A.staticProp = {foo: 'bar'};
ns2.B = class extends ns1.A {}
// ns2.B isn't an alias for ns1.A, so it isn't inlined.
// ns2.B?.staticProp isn't an alias, because it's an optional chain.
use(ns2.B?.staticProp.bar);
`)
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_es5Class()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_es5Class() {
	// ES6 classes do not inherit static properties of ES5 class constructors.
	c.TestSame(
		`/** @constructor */
function A() {}
A.staticProp = 5;
class B extends A {}
use(B.staticProp);
`)
	// undefined
}

// Java: AggressiveInlineAliasesTest.testClassStaticInheritance_cantDetermineSuperclass()
func (c *AggressiveInlineAliasesTest) TestClassStaticInheritance_cantDetermineSuperclass() {
	// Currently we only inline inherited properties when the extends clause contains a simple or
	// qualified name.
	c.Test(
		`class A {}
A.foo = 5;
class B {}
B.foo = 6;
function getSuperclass() { return A; }
class C extends getSuperclass() {}
use(C.foo);
`,
		`class A {}
A.foo = 5;
class B {}
B.foo = 6;
function getSuperclass() {
  return A;
}
const CLASS_EXTENDS$0 = getSuperclass();
class C extends CLASS_EXTENDS$0 {}
use(C.foo);
`)
}

// Java: AggressiveInlineAliasesTest.testAliasInsideGenerator()
func (c *AggressiveInlineAliasesTest) TestAliasInsideGenerator() {
	c.Test(
		"const a = {b: 5}; const c = a;    function *f() { yield c.b; }",
		"const a = {b: 5}; const c = null; function *f() { yield a.b; }")
}

// Java: AggressiveInlineAliasesTest.testAliasInsideModuleScope()
func (c *AggressiveInlineAliasesTest) TestAliasInsideModuleScope() {
	// CollapseProperties currently only handles global variables, so we don't handle aliasing in
	// module bodies here.
	c.TestSame("const a = {b: 5}; const c = a; export default function() {};")
}

// Java: AggressiveInlineAliasesTest.testAliasForSuperclassNamespace()
func (c *AggressiveInlineAliasesTest) TestAliasForSuperclassNamespace() {
	c.Test(
		`var ns = {};
class Foo {}
ns.clazz = Foo;
var Bar = class extends ns.clazz.Baz {}
`,
		`var ns = {};
class Foo {}
ns.clazz = null;
var Bar = class extends Foo.Baz {}
`)
	c.Test(
		`var ns = {};
class Foo {}
Foo.Builder = class {}
ns.clazz = Foo;
var Bar = class extends ns.clazz.Builder {}
`,
		`var ns = {};
class Foo {}
Foo.Builder = class {}
ns.clazz = null;
var Bar = class extends Foo.Builder {}
`)
}

// Java: AggressiveInlineAliasesTest.testAliasForSuperclass_withStaticInheritance()
func (c *AggressiveInlineAliasesTest) TestAliasForSuperclass_withStaticInheritance() {
	c.Test(
		`var ns = {};
class Foo {}
Foo.baz = 3;
ns.clazz = Foo;
var Bar = class extends ns.clazz {}
use(Bar.baz);
`,
		`var ns = {};
class Foo {}
Foo.baz = 3;
ns.clazz = null;
var Bar = class extends Foo {}
use(Foo.baz);
`)
}

// Java: AggressiveInlineAliasesTest.testStaticInheritance_superclassIsStaticProperty()
func (c *AggressiveInlineAliasesTest) TestStaticInheritance_superclassIsStaticProperty() {
	c.Test(
		`class Foo {}
Foo.Builder = class {}
Foo.Builder.baz = 3;
var Bar = class extends Foo.Builder {}
use(Bar.baz);
`,
		`class Foo {}
Foo.Builder = class {}
Foo.Builder.baz = 3;
var Bar = class extends Foo.Builder {}
use(Foo.Builder.baz);
`)
}

// Java: AggressiveInlineAliasesTest.testAliasForSuperclassNamespace_withStaticInheritance()
func (c *AggressiveInlineAliasesTest) TestAliasForSuperclassNamespace_withStaticInheritance() {
	c.Test(
		`var ns = {};
class Foo {}
Foo.Builder = class {}
Foo.Builder.baz = 3;
ns.clazz = Foo;
var Bar = class extends ns.clazz.Builder {}
use(Bar.baz);
`,
		`var ns = {};
class Foo {}
Foo.Builder = class {}
Foo.Builder.baz = 3;
ns.clazz = null;
var Bar = class extends Foo.Builder {}
use(Foo.Builder.baz);
`)
}

// Java: AggressiveInlineAliasesTest.testGithubIssue2754()
func (c *AggressiveInlineAliasesTest) TestGithubIssue2754() {
	c.TestSame(
		`var ns = {};
/** @constructor */
ns.Bean = function() {}

ns.Bean.x = function(a=null) {
  if (a == null || a === undefined) {
    a = ns.Bean;
  }
  return new a();
}
`)
}

// Java: AggressiveInlineAliasesTest.testParameterDefaultIsAlias()
func (c *AggressiveInlineAliasesTest) TestParameterDefaultIsAlias() {
	c.TestSame(
		`/** @constructor */
function a() {}
a.staticProp = 5;

function f(param = a) {
  return new param();
}
`)
}

// Java: AggressiveInlineAliasesTest.testInliningPropertyAliasBeforeNamespace_withConstructorProperty()
func (c *AggressiveInlineAliasesTest) TestInliningPropertyAliasBeforeNamespace_withConstructorProperty() {
	c.Test(
		`var prop = 1;
/** @constructor */
var Foo = function() {}

Foo.prop = prop;

var aliasOfFoo = Foo;
alert(aliasOfFoo.prop);
`,
		`var prop = 1;
/** @constructor */
var Foo = function() {}

Foo.prop = null;

var aliasOfFoo = null;
alert(prop);
`)
}

// Java: AggressiveInlineAliasesTest.testInliningPropertyAliasBeforeNamespace_withNonConstructorProperty()
func (c *AggressiveInlineAliasesTest) TestInliningPropertyAliasBeforeNamespace_withNonConstructorProperty() {
	// NOTE: the only difference between this input and the above input is that the
	// above has annotated Foo with @constructor. The @constructor case used to generate bad output
	c.Test(
		`var prop = 1;
/** @const */
var Foo = {}

Foo.prop = prop;

var aliasOfFoo = Foo;
alert(aliasOfFoo.prop);
`,
		`var prop = 1;
/** @const */
var Foo = {};

Foo.prop = null;

var aliasOfFoo = null;
alert(prop);
`)
}

// Java: AggressiveInlineAliasesTest.testInliningPropertyAlias_onEscapedConstructor()
func (c *AggressiveInlineAliasesTest) TestInliningPropertyAlias_onEscapedConstructor() {
	// This test demonstrates somewhat unsafe behavior in this pass that ends up working around
	// CollapseProperties behavior.
	//
	c.Test(
		`var prop = 1;
/** @constructor */
var Foo = function() {}

Foo.prop = prop;

/** @constructor */
function Bar() {}
Bar.aliasOfFoo = Foo;
use(Bar);
alert(Bar.aliasOfFoo.prop);
`,
		`var prop = 1;
/** @constructor */
var Foo = function() {}

Foo.prop = prop;

/** @constructor */
function Bar() {}
Bar.aliasOfFoo = Foo; // don't remove this initialization
use(Bar); // this might write to Bar.aliasOfFoo.prop
alert(prop);
`)
	// but set this to prop anyway to fix CollapseProperties
}

// Java: AggressiveInlineAliasesTest.testAliasChain()
func (c *AggressiveInlineAliasesTest) TestAliasChain() {
	c.Test(
		`var goog = {};
goog.DEBUG = false;
var foo = {};
var global_DEBUG = goog.DEBUG;
foo.DEBUG = global_DEBUG;
alert(foo.DEBUG);
`,
		`var goog = {};
goog.DEBUG = false;
var foo = {};
var global_DEBUG = null;
foo.DEBUG = null;
alert(goog.DEBUG);
`)
	c.Test(
		`var goog = {};
goog.DEBUG = false;
var foo = {};
var global_DEBUG = goog.DEBUG;
foo.DEBUG = global_DEBUG;
alert(foo?.DEBUG);
`,
		`var goog = {};
goog.DEBUG = false;
var foo = {};
var global_DEBUG = null;
// Does not get inlined because of the optional chain reference
foo.DEBUG = goog.DEBUG;
alert(foo?.DEBUG);
`)
}

// Java: AggressiveInlineAliasesTest.testTranspiledEs6StaticMethods_withNoCollapse()
func (c *AggressiveInlineAliasesTest) TestTranspiledEs6StaticMethods_withNoCollapse() {
	// This is what transpiled ES6 class statics look like.
	// We don't replace "Child.f = Parent.f" with "Child.f = null" because of the @nocollapse
	c.DisableCompareJsDoc()
	// multistage compilation erases the @extends
	c.TestSame(
		`/** @constructor */ var Parent = function() {};
/** @nocollapse */ Parent.f = function() {};
/** @constructor @extends {Parent} @param {...?} var_args  */
var Child = function(var_args) {
  Parent.apply(this, arguments);
}
$jscomp.inherits(Child, Parent);
/** @nocollapse */ Child.f = Parent.f;
Child.prototype.g = function() { return this.f(); }
`)
}

// Java: AggressiveInlineAliasesTest.testLoopInAliasChainWithTypedefConstructorProperty()
func (c *AggressiveInlineAliasesTest) TestLoopInAliasChainWithTypedefConstructorProperty() {
	// This kind of code can get produced by module exports rewriting and was causing a crash in
	// AggressiveInlineAliases.
	c.DisableCompareJsDoc()
	// multistage compilation erases the @typedef
	c.TestSame(
		`/** @constructor */ var Item = function() {};
/** @typedef {number} */ Item.Models;
Item.Models = Item.Models;
`)
}

// Java: AggressiveInlineAliasesTest.testDontInlinePropertiesOnEscapedNamespace()
func (c *AggressiveInlineAliasesTest) TestDontInlinePropertiesOnEscapedNamespace() {
	c.TestSame(
		c.Externs("function use(obj) {}"),
		c.Srcs(
			`function Foo() {}
Foo.Bar = {};
Foo.Bar.baz = {A: 1, B: 2};

var $jscomp$destructuring$var1 = Foo.Bar;
// This call could potentially have changed the value of Foo.Bar, so don't replace
// $jscomp$destructuring$var1.baz with Foo.Bar.baz
use(Foo);
var baz = $jscomp$destructuring$var1.baz;
use(baz.A);
`))
}

// Java: AggressiveInlineAliasesTest.testInlinePropertiesOnEscapedNamespace_withDeclaredType()
func (c *AggressiveInlineAliasesTest) TestInlinePropertiesOnEscapedNamespace_withDeclaredType() {
	c.Test(
		c.Externs("function use(obj) {}"),
		c.Srcs(
			`/** @constructor */
function Foo() {}
/** @constructor */
Foo.Bar = function() {};
/** @enum {number} */
Foo.Bar.baz = {A: 1, B: 2};

var $jscomp$destructuring$var1 = Foo.Bar;
// This call could potentially have changed the value of Foo.Bar.
use(Foo);
var baz = $jscomp$destructuring$var1.baz;
use(baz.A);
`),
		c.Expected(
			`/** @constructor */
function Foo() {}
/** @constructor */
Foo.Bar = function() {};
/** @enum {JSDocSerializer_placeholder_type} */
Foo.Bar.baz = {A: 1, B: 2};

var $jscomp$destructuring$var1 = null;
use(Foo);
var baz = null;
// If we didn't unsafely replace baz.A with Foo.Bar.baz.A, this reference would
// break after CollapseProperties runs because Foo.Bar.baz -> Foo$Bar$baz
// So although inlining this is technically unsafe, because use(Foo) could have
// changed the value of Foo.Bar, it actually fixes a breakage caused by
// CollapseProperties.
use(Foo.Bar.baz.A);
`))
}

// Java: AggressiveInlineAliasesTest.testDontInlinePropertiesOnNamespace_withNoCollapse()
func (c *AggressiveInlineAliasesTest) TestDontInlinePropertiesOnNamespace_withNoCollapse() {
	c.DisableCompareJsDoc()
	// multistage compilation removes the @enum type
	c.TestSame(
		c.Externs("function use(obj) {}"),
		c.Srcs(
			`/** @constructor */
function Foo() {}
/** @constructor @nocollapse */
Foo.Bar = function() {};
/** @enum {number} @nocollapse */
Foo.Bar.baz = {A: 1, B: 2};

var $jscomp$destructuring$var1 = Foo.Bar;
use(Foo);
var baz = $jscomp$destructuring$var1.baz;
use(baz.A);
`))
}

// Java: AggressiveInlineAliasesTest.testNestedAssignWithAlias()
func (c *AggressiveInlineAliasesTest) TestNestedAssignWithAlias() {
	c.Test(
		`var ns = {a: {}};
var Letters = { B: 'b'};
ns.c = ns.a.b = Letters.B;
use(ns.c);
use(ns.a.b);
`,
		`var ns = {a: {}};
var Letters = {B: 'b'};
// test that we handle nested assigns correctly
ns.c = null;
use(Letters.B);
use(Letters.B);
`)
}

// Java: AggressiveInlineAliasesTest.testCommaExpression()
func (c *AggressiveInlineAliasesTest) TestCommaExpression() {
	c.Test(
		`/** @const */ var exports = {};
/** @const @enum {string} */ var Letters = {
  A: 'a',
  B: 'b'};
exports.A = Letters.A, exports.B = Letters.B;
use(exports.B);
`,
		"/** @const */ var exports = {};\n"+
			"/** @const @enum {JSDocSerializer_placeholder_type} */ var Letters = {\n"+
			"  A: 'a',\n"+
			"  B: 'b'};\n"+
			"// this used to become\n"+
			"//   exports.A = null, Letters.B = null;\n"+
			"// breaking all references to Letters.B\n"+
			"// GlobalNamespace treats 'exports.B' as both a read and write of\n"+
			"// 'exports.B', and AggressiveInlineAliases used to replace all reads of exports.B with\n"+
			"// `Letters.B` without verifying that the read was not also a write.\n"+
			"Letters.A, Letters.B;\n"+
			"use(Letters.B);\n")
}

// Java: AggressiveInlineAliasesTest.testDontInlineDestructuredAliasProp_quoted()
func (c *AggressiveInlineAliasesTest) TestDontInlineDestructuredAliasProp_quoted() {
	c.TestSame("var a = {x: 2}; var b = {}; b['y'] = a.x; const {'y': y} = b; use(y);")
}

// Java: AggressiveInlineAliasesTest.testInlineDestructuredAliasProp()
func (c *AggressiveInlineAliasesTest) TestInlineDestructuredAliasProp() {
	c.Test(
		"var a = {x: 2}; var b = {}; b.y =  a.x; const {y} =    b; use( y );",
		"var a = {x: 2}; var b = {}; b.y = null; const  y  = null; use(a.x);")
}

// Java: AggressiveInlineAliasesTest.testInlineDestructuredAliasPropWithKeyBefore()
func (c *AggressiveInlineAliasesTest) TestInlineDestructuredAliasPropWithKeyBefore() {
	c.Test(
		"var a = {x: 2}; var b = {z: 3}; b.y = a.x; b.z = 4; const {z, y} = b; use(y + z);",
		`var a = {x: 2};
var b = {z: 3};
b.y = null;
b.z = 4;
const z = b.z;
const y = null;
use(a.x + z);
`)
}

// Java: AggressiveInlineAliasesTest.testInlineDestructuredAliasPropWithKeyAfter()
func (c *AggressiveInlineAliasesTest) TestInlineDestructuredAliasPropWithKeyAfter() {
	c.Test(
		"var a = {x: 2}; var b = {z: 3}; b.y = a.x; b.z = 4; const {y, z} = b; use(y + z);",
		`var a = {x: 2};
var b = {z: 3};
b.y = null;
b.z = 4;
const y = null;
const z = b.z;
use(a.x + z);
`)
}

// Java: AggressiveInlineAliasesTest.testInlineDestructuredAliasPropWithKeyBeforeAndAfter()
func (c *AggressiveInlineAliasesTest) TestInlineDestructuredAliasPropWithKeyBeforeAndAfter() {
	c.Test(
		`var a = {x: 2};
var b = {z: 3};
b.y = a.x;
b.z = 4; // add second assign so that this won't get inlined
const {x, y, z} = b;
use(y + z);
`,
		`var a = {x: 2};
var b = {z: 3};
b.y = null;
b.z = 4;
const x = b.x;
const y = null;
const z = b.z;
use(a.x + z);
`)
}

// Java: AggressiveInlineAliasesTest.testDestructuredPropAccessInAssignWithKeyBefore()
func (c *AggressiveInlineAliasesTest) TestDestructuredPropAccessInAssignWithKeyBefore() {
	c.TestSame(
		`var a = {x: 2};
var b = {};
b.y = a.x;
var obj = {};
({missing: obj.foo, y: obj.foo} = b);
use(obj.foo);
`)
}

// Java: AggressiveInlineAliasesTest.testDestructuredPropAccessInAssignWithKeyAfter()
func (c *AggressiveInlineAliasesTest) TestDestructuredPropAccessInAssignWithKeyAfter() {
	c.TestSame(
		`var a = {x: 2};
var b = {};
b.y = a.x;
var obj = {};
({y: obj.foo, missing: obj.foo} = b);
use(obj.foo);
`)
}

// Java: AggressiveInlineAliasesTest.testDestructuredPropAccessInDeclarationWithDefault()
func (c *AggressiveInlineAliasesTest) TestDestructuredPropAccessInDeclarationWithDefault() {
	c.TestSame(
		`var a = {x: {}};
var b = {};
b.y = a.x;
const {y = 0} = b;
use(y);
`)
}

// Java: AggressiveInlineAliasesTest.testDestructuringPropertyOnAliasedNamespace()
func (c *AggressiveInlineAliasesTest) TestDestructuringPropertyOnAliasedNamespace() {
	// We can inline a part of a getprop chain on the rhs of a destructuring pattern:
	//   replace 'alias -> a.b' in 'const {c} = alias.Enum;'
	c.Test(
		`const a = {};
/** @const */ a.b = {};
/** @enum {string} */ a.b.Enum = {c: 'c'};

const alias = a.b;
function f() { const {c} = alias.Enum; use(c); }
`,
		`const a = {};
/** @const */ a.b = {};
/** @enum {JSDocSerializer_placeholder_type} */ a.b.Enum = {c: 'c'};

const alias = null;
function f() { const c = null; use(a.b.Enum.c); }
`)
}

// Java: AggressiveInlineAliasesTest.testReplaceSuperGetPropInStaticMethod()
func (c *AggressiveInlineAliasesTest) TestReplaceSuperGetPropInStaticMethod() {
	c.Test(
		`class Foo {
  static m() {}
}
class Bar extends Foo {
  static m() {
    super.m();
  }
}
`,
		`class Foo {
  static m() {}
}
class Bar extends Foo {
  static m() {
    Foo.m();
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.testReplaceSuperInArrowInStaticMethod()
func (c *AggressiveInlineAliasesTest) TestReplaceSuperInArrowInStaticMethod() {
	c.Test(
		`class Foo {
  static m() {}
}
class Bar extends Foo {
  static m() {
    return () => super.m();
  }
}
`,
		`class Foo {
  static m() {}
}
class Bar extends Foo {
  static m() {
    return () => {
      return Foo.m();
    };
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.testReplaceSuperInStaticMethodWithQualifiedNameSuperclass()
func (c *AggressiveInlineAliasesTest) TestReplaceSuperInStaticMethodWithQualifiedNameSuperclass() {
	c.Test(
		`const a = {b:{}};
/** @const */
a.b.Foo = class {
  static m() {}
};
class Bar extends a.b.Foo {
  static m() {
    super.m();
  }
}
`,
		`const a = {b:{}};
/** @const */
a.b.Foo = class {
  static m() {}
};
class Bar extends a.b.Foo {
  static m() {
    a.b.Foo.m();
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.testDontReplaceSuperInObjectLiteralMethod()
func (c *AggressiveInlineAliasesTest) TestDontReplaceSuperInObjectLiteralMethod() {
	c.TestSame("var obj = {m() { super.n(); } };")
}

// Java: AggressiveInlineAliasesTest.testDontReplaceSuperInObjectLitFnInStaticClassMethod()
func (c *AggressiveInlineAliasesTest) TestDontReplaceSuperInObjectLitFnInStaticClassMethod() {
	c.TestSame(
		"class Foo { static m() {} }\n" +
			"// `super.n` refers to a different object than `Foo.n`\n" +
			"class Bar extends Foo { static m() { return {m() { super.n(); }}; } }\n")
}

// Java: AggressiveInlineAliasesTest.testDestructuredClassAlias()
func (c *AggressiveInlineAliasesTest) TestDestructuredClassAlias() {
	c.Test(
		`const ns = {};
ns.Foo = class {};
ns.Foo.STR = '';
const {Foo} = ns;
foo(Foo.STR);
`,
		`const ns = {};
ns.Foo = class {};
ns.Foo.STR = '';
const Foo = null;
foo(ns.Foo.STR);
`)
}

// Java: AggressiveInlineAliasesTest.testReplaceChainedSuperRefInStaticMethod()
func (c *AggressiveInlineAliasesTest) TestReplaceChainedSuperRefInStaticMethod() {
	c.Test(
		`class Foo {
  static m() {}
}
class Bar extends Foo {}
class Baz extends Bar {
  static m() {
    super.m();
  }
}
`,
		`class Foo {
  static m() {}
}
class Bar extends Foo {}
class Baz extends Bar {
  static m() {
    Foo.m();
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.testReplaceSuperInStaticMethodWithNonQnameSuperclass()
func (c *AggressiveInlineAliasesTest) TestReplaceSuperInStaticMethodWithNonQnameSuperclass() {
	c.Test(
		`class Bar extends getClass() {
  static m() {
    super.m();
  }
}
`,
		`const CLASS_EXTENDS$0 = getClass();
class Bar extends CLASS_EXTENDS$0 {
  static m() {
    CLASS_EXTENDS$0.m();
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.testReplaceSuperInStaticMethodWithNonQnameGetPropSuperclass()
func (c *AggressiveInlineAliasesTest) TestReplaceSuperInStaticMethodWithNonQnameGetPropSuperclass() {
	c.Test(
		`class Bar extends getClasses().Foo {
  static m() {
    super.m();
  }
}
`,
		`const CLASS_EXTENDS$0 = getClasses().Foo;
class Bar extends CLASS_EXTENDS$0 {
  static m() {
    CLASS_EXTENDS$0.m();
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.testReplaceSuperGetElemInStaticMethod()
func (c *AggressiveInlineAliasesTest) TestReplaceSuperGetElemInStaticMethod() {
	// while CollapseProperties won't collapse Foo['m'], replacing `super` enables collapsing Bar.m
	// (since CollapseProperties cannot collapse methods using super) and helps code size.
	c.Test(
		`class Foo {
  static 'm'() {}
}
class Bar extends Foo {
  static m() {
    super['m']();
  }
}
`,
		`class Foo {
  static ["m"]() {}
}
class Bar extends Foo {
  static m() {
    Foo["m"]();
  }
}
`)
}

// Java: AggressiveInlineAliasesTest.testDontReplaceSuperInClassPrototypeMethod()
func (c *AggressiveInlineAliasesTest) TestDontReplaceSuperInClassPrototypeMethod() {
	c.TestSame(
		`class Foo {
  m() {}
}
class Bar extends Foo {
  m() {
    super.m();
  }
}
`)
}

// ValidateGlobalNamespace: to ensure that as we modify the AST, the GlobalNamespace stays up-to-date, we do a consistency
// check after every unit test.
//
// <p>This check compares the names in the global namespace in the pass with a freshly-created
// global namespace.
//
// Java: AggressiveInlineAliasesTest.validateGlobalNamespace(GlobalNamespace)
func (c *AggressiveInlineAliasesTest) ValidateGlobalNamespace(passGlobalNamespace *GlobalNamespace) {
	expectedGlobalNamespace := NewGlobalNamespace(c.GetLastCompiler(), c.GetLastCompiler().GetJsRoot())
	// GlobalNamespace (understandably) does not override equals. It would be silly to put it in
	// a datastructure. Neither does GlobalNamespace.Name (which probably could?)
	// So to compare equality: we verify that
	//  1. the two namespaces have the same qualified names, bar extern names
	//  2. each name has the same number of references in both namespaces
	//  3. each name has the same computed `Inlinability`
	for _, expectedName := range expectedGlobalNamespace.GetNameForest() {
		if expectedName.InExterns() {
			continue
		}
		fullName := expectedName.GetFullName()
		actualName, _ := passGlobalNamespace.GetSlot(expectedName.GetFullName()).(*GlobalNamespace_Name)
		AssertWithMessage(fullName).That(actualName).IsNotNull()
		AssertWithMessage(fullName).
			That(actualName.GetAliasingGets()).
			IsEqualTo(expectedName.GetAliasingGets())
		AssertWithMessage(fullName).
			That(actualName.GetSubclassingGets()).
			IsEqualTo(expectedName.GetSubclassingGets())
		AssertWithMessage(fullName).
			That(actualName.GetLocalSets()).
			IsEqualTo(expectedName.GetLocalSets())
		AssertWithMessage(fullName).
			That(actualName.GetGlobalSets()).
			IsEqualTo(expectedName.GetGlobalSets())
		AssertWithMessage(fullName).
			That(actualName.GetDeleteProps()).
			IsEqualTo(expectedName.GetDeleteProps())
		AssertWithMessage(fullName).
			That(actualName.GetCallGets()).
			IsEqualTo(expectedName.GetCallGets())
		AssertWithMessage("%s: canCollapseOrInline()", fullName).
			That(actualName.CanCollapseOrInline()).
			IsEqualTo(expectedName.CanCollapseOrInline())
		AssertWithMessage("%s: canCollapseOrInlineChildNames()", fullName).
			That(actualName.CanCollapseOrInlineChildNames()).
			IsEqualTo(expectedName.CanCollapseOrInlineChildNames())
	}
	// Verify that no names in the actual name forest are not present in the expected name forest
	for _, actualName := range passGlobalNamespace.GetNameForest() {
		actualFullName := actualName.GetFullName()
		AssertWithMessage(actualFullName).
			That(expectedGlobalNamespace.GetSlot(actualFullName)).
			IsNotNull()
	}
}

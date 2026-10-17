// Port of com/google/javascript/jscomp/InlineAndCollapsePropertiesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/InlineAndCollapsePropertiesTest.java)

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
)

// InlineAndCollapsePropertiesTest: tests for InlineAndCollapseProperties.
type InlineAndCollapsePropertiesTest struct {
	CompilerTestCaseBase
}

// Java: InlineAndCollapsePropertiesTest.EXTERNS
const InlineAndCollapsePropertiesTest_EXTERNS = `var window;
function alert(s) {}
function parseInt(s) {}
/** @constructor */ function String() {};
var arguments
`

// Java: InlineAndCollapsePropertiesTest()
func NewInlineAndCollapsePropertiesTest() *InlineAndCollapsePropertiesTest {
	c := &InlineAndCollapsePropertiesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, InlineAndCollapsePropertiesTest_EXTERNS)
	return c
}

func TestInlineAndCollapsePropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInlineAndCollapsePropertiesTest() })
}

// Java: InlineAndCollapsePropertiesTest.getProcessor(Compiler)
func (c *InlineAndCollapsePropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("es6NormalizeClasses", func(comp AbstractCompiler) CompilerPass { return NewEs6NormalizeClasses(comp) }))
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"inlineAndCollapseProperties",
			func(comp AbstractCompiler) CompilerPass {
				return InlineAndCollapseProperties_Builder_(comp).
					SetPropertyCollapseLevel(CompilerOptions_PropertyCollapseLevel_ALL).
					SetChunkOutputType(CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE).
					SetHaveModulesBeenRewritten(false).
					SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER).
					Build()
			}))
	return optimizer
}

// Java: InlineAndCollapsePropertiesTest.setUp()
func (c *InlineAndCollapsePropertiesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.DisableCompareJsDoc()
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
}

// Java: InlineAndCollapsePropertiesTest.testTs52OutputChange()
func (c *InlineAndCollapsePropertiesTest) TestTs52OutputChange() {
	c.Test(
		c.Srcs(
			`var alias;
var module$exports$C = class {
  method1() {
    return alias.staticPropOnC;
  }
  method2() {
    return alias.staticPropOnAlias;
  }
};
alias = module$exports$C;
(() => {
  alias.staticPropOnAlias = 1;
})();
module$exports$C.staticPropOnC = 2;
`),
		c.Expected(
			`var alias;
var module$exports$C = class {
  method1() {
    return module$exports$C$staticPropOnC;
  }
  method2() {
    return module$exports$C$staticPropOnAlias;
  }
};
var module$exports$C$staticPropOnAlias;
alias = null;
(() => {
  module$exports$C$staticPropOnAlias = 1;
})();
var module$exports$C$staticPropOnC = 2;
`))
}

// Java: InlineAndCollapsePropertiesTest.testTs52OutputChangeVariableReassignment()
func (c *InlineAndCollapsePropertiesTest) TestTs52OutputChangeVariableReassignment() {
	// This test case was created by executing
	// AdvancedOptimizationsIntegrationTest#testTSVariableReassignmentAndAliasingDueToDecoration
	// with `options.setPrintSourceAfterEachPass(true)`.
	// The input source here is how the source code looks at the point just before
	// InlineAndCollapseProperties runs in that test case except for the addition of one line of
	// JSDoc to make the unit test recognize a function declaration as a class declaration.
	c.Test(
		c.Externs(
			[]*SourceFile{
				NewTestExternsBuilder().
					AddObject().
					AddConsole().
					AddClosureExterns().
					AddExtra(
						// simulate "const tslib_1 = goog.require('tslib');",
						`var tslib_1 = {
  __decorate: function(decorators, clazz) {}
};
`).
					BuildExternsFile("externs.js")}),
		c.Srcs(
			`var module$exports$main = {};
var module$contents$main_module =
    module$contents$main_module || {id: 'main.ts'};
var module$contents$main_Foo_1;
function module$contents$main_noopDecorator(arg) {
  return arg;
}
// This JSDoc annotation makes this unit test recognize the declaration as a
// class definition.
/** @constructor */
var i0$classdecl$var0 = function() {};
i0$classdecl$var0.foo = function() {
  console.log('Hello');
};
i0$classdecl$var0.prototype.bar = function() {
  module$contents$main_Foo_1.foo();
  console.log('ID: ' + module$contents$main_Foo_1.ID + '');
};
var module$contents$main_Foo = module$contents$main_Foo_1 = i0$classdecl$var0;
module$contents$main_Foo.ID = 'original';
module$contents$main_Foo.ID2 = module$contents$main_Foo_1.ID;
(function() {
module$contents$main_Foo_1.foo();
console.log('ID: ' + module$contents$main_Foo_1.ID + '');
})();
module$contents$main_Foo = module$contents$main_Foo_1 = tslib_1.__decorate(
    [module$contents$main_noopDecorator], module$contents$main_Foo);
if (false) {
  module$contents$main_Foo.ID;
  module$contents$main_Foo.ID2;
}
(new module$contents$main_Foo()).bar();
`),
		c.Expected(
			`var module$contents$main_module =
    module$contents$main_module || {id: 'main.ts'};
var module$contents$main_Foo_1;
function module$contents$main_noopDecorator(arg) {
  return arg;
}
/** @constructor */
var i0$classdecl$var0 = function() {};
// TODO : b/299055739 - This property collapse breaks the output code.
var i0$classdecl$var0$foo = function() {
// "i0$classdecl$var0.foo = function() {", // this is the correct line
  console.log('Hello');
};
i0$classdecl$var0.prototype.bar = function() {
// This reference to foo() is broken.
  module$contents$main_Foo_1.foo();
  console.log('ID: ' + module$contents$main_Foo_1.ID + '');
};
var module$contents$main_Foo = module$contents$main_Foo_1 = i0$classdecl$var0;
module$contents$main_Foo.ID = 'original';
module$contents$main_Foo.ID2 = module$contents$main_Foo_1.ID;
(function() {
module$contents$main_Foo_1.foo();
console.log('ID: ' + module$contents$main_Foo_1.ID + '');
})();
module$contents$main_Foo = module$contents$main_Foo_1 = tslib_1.__decorate(
    [module$contents$main_noopDecorator], module$contents$main_Foo);
if (false) {
  module$contents$main_Foo.ID;
  module$contents$main_Foo.ID2;
}
(new module$contents$main_Foo()).bar();
`))
}

// Java: InlineAndCollapsePropertiesTest.testDoNotCollapseDeletedProperty()
func (c *InlineAndCollapsePropertiesTest) TestDoNotCollapseDeletedProperty() {
	c.TestSame(
		c.Srcs(
			`const global = window;
delete global.HTMLElement;
global.HTMLElement = (class {});
`))
}

// Java: InlineAndCollapsePropertiesTest.testConstObjRefInTemplateLiteralComputedPropKey()
func (c *InlineAndCollapsePropertiesTest) TestConstObjRefInTemplateLiteralComputedPropKey() {
	c.Test(
		c.Srcs(
			"var module$name = {}\n"+
				"module$name.cssClasses = {\n"+
				"  CLASS_A: 'class-a',\n"+
				"};\n"+
				"\n"+
				"module$name.oldCssClassesMap = {\n"+
				"  [`${module$name.cssClasses.CLASS_A}`]: 'old-class-a',\n"+
				"};\n"),
		c.Expected(
			"var module$name$cssClasses$CLASS_A = 'class-a';\n"+
				"var module$name$oldCssClassesMap = {\n"+
				"  [`${module$name$cssClasses$CLASS_A}`]:'old-class-a'\n"+
				"};\n"))
}

// Java: InlineAndCollapsePropertiesTest.testObjLitSpread()
func (c *InlineAndCollapsePropertiesTest) TestObjLitSpread() {
	c.Test(
		`const other = {bar: 'some' };
let ns = {}
ns.foo = { ...other };
use(ns.foo.bar);
`,
		`const other = {bar: 'some' };
var ns$foo = { ...other };
use(ns$foo.bar);
`)
}

// same behavior with transpiled version of above test
// Java: InlineAndCollapsePropertiesTest.testObjAssign()
func (c *InlineAndCollapsePropertiesTest) TestObjAssign() {
	c.Test(
		`const other = {bar: 'some' };
let ns = {}
ns.foo = Object.assign({}, other);
use(ns.foo.bar);
`,
		`const other = {bar: 'some' };
var ns$foo = Object.assign({}, other);
use(ns$foo.bar);
`)
}

// Java: InlineAndCollapsePropertiesTest.testObjLitSpread_twoSpreads()
func (c *InlineAndCollapsePropertiesTest) TestObjLitSpread_twoSpreads() {
	c.Test(
		`const other = {bar: 'bar' };
const another = {baz: 'baz' };
let ns = {};
ns.foo = { ...other, ...another };
use(ns.foo.bar, ns.foo.baz);
`,
		`const other = {bar: 'bar' };
const another = {baz : 'baz'};
var ns$foo = { ...other, ...another };
use(ns$foo.bar, ns$foo.baz);
`)
}

// same behavior with transpiled version of above test
// Java: InlineAndCollapsePropertiesTest.testObjAssign_twoSpreads()
func (c *InlineAndCollapsePropertiesTest) TestObjAssign_twoSpreads() {
	c.Test(
		`const other = {bar: 'bar'};
const another = {baz: 'baz'};
let ns = {};
ns.foo = Object.assign({}, other, another);
use(ns.foo.bar, ns.foo.baz);
`,
		`const other = {bar: 'bar'};
const another = {baz : 'baz'};
var ns$foo = Object.assign({}, other, another);
use(ns$foo.bar, ns$foo.baz);
`)
}

// Java: InlineAndCollapsePropertiesTest.testObjLitSpread_withNormalPropAfter()
func (c *InlineAndCollapsePropertiesTest) TestObjLitSpread_withNormalPropAfter() {
	c.Test(
		`const other = {bar: 'some' };
let ns = {}
ns.foo = { ...other, prop : 0};
use(ns.foo.bar, ns.foo.prop);
`,
		`const other = {bar: 'some' };
var ns$foo$prop = 0;
var ns$foo = { ...other};
use(ns$foo.bar, ns$foo$prop);
`)
}

// Java: InlineAndCollapsePropertiesTest.testObjAssign_withNormalPropAfter()
func (c *InlineAndCollapsePropertiesTest) TestObjAssign_withNormalPropAfter() {
	c.Test(
		`const other = {bar: 'some' };
let ns = {}
ns.foo = Object.assign({}, other, {prop : 0});
use(ns.foo.bar, ns.foo.prop);
`,
		`const other = {bar: 'some' };
var ns$foo = Object.assign({}, other, {prop:0});
use(ns$foo.bar, ns$foo.prop);
`)
	// both properties not collapsed.
}

// Java: InlineAndCollapsePropertiesTest.testObjLitSpread_chained()
func (c *InlineAndCollapsePropertiesTest) TestObjLitSpread_chained() {
	c.Test(
		`const other = {bar: 'some' };
let ns = {...other}
let ns2 = { ...ns };
use(ns2.bar);
`,
		`const other = {bar: 'some' };
let ns = { ...other };
let ns2 = { ...ns };
use(ns2.bar);
`)
}

// same behavior with transpiled version of above test
// Java: InlineAndCollapsePropertiesTest.testObjAssign_chained()
func (c *InlineAndCollapsePropertiesTest) TestObjAssign_chained() {
	c.Test(
		`const other = {bar: 'some' };
let ns = Object.assign({}, other);
let ns2 = Object.assign({}, ns);
use(ns2.bar);
`,
		`const other = {bar: 'some' };
let ns = Object.assign({}, other);
let ns2 = Object.assign({}, ns);
use(ns2.bar);
`)
}

// Java: InlineAndCollapsePropertiesTest.testObjLitSpread_chainedWithGetProp()
func (c *InlineAndCollapsePropertiesTest) TestObjLitSpread_chainedWithGetProp() {
	c.Test(
		`const other = {bar: 'some' };
let ns = {...other};
let ns2 = {};
ns2.foo = { ...ns };
use(ns2.foo.bar);
`,
		`const other = {bar: 'some' };
let ns = { ...other };
var ns2$foo = { ...ns };
use(ns2$foo.bar);
`)
}

// same behavior with transpiled version of above test
// Java: InlineAndCollapsePropertiesTest.testObjAssign_chainedWithGetProp()
func (c *InlineAndCollapsePropertiesTest) TestObjAssign_chainedWithGetProp() {
	c.Test(
		`const other = {bar: 'some' };
let ns = Object.assign({}, other);
let ns2 = {};
ns2.foo = Object.assign({}, ns);
use(ns2.foo.bar);
`,
		`const other = {bar: 'some' };
let ns = Object.assign({}, other);
var ns2$foo = Object.assign({}, ns);
use(ns2$foo.bar);
`)
}

// Java: InlineAndCollapsePropertiesTest.testGitHubIssue3733()
func (c *InlineAndCollapsePropertiesTest) TestGitHubIssue3733() {
	c.TestSame(
		c.Srcs(
			`const X = {Y: 1};

function fn(a) {
  if (a) {
// Before issue #3733 was fixed GlobalNamespace failed to see this reference
// as creating an alias for X due to a switch statement that failed to check
// for the RETURN node type, so X.Y was incorrectly collapsed.
    return a ? X : {};
  }
}

console.log(fn(true).Y);
`))
}

// Java: InlineAndCollapsePropertiesTest.testCollapse()
func (c *InlineAndCollapsePropertiesTest) TestCollapse() {
	c.Test(
		"var a = {}; a.b = {}; var c = a.b;", //
		"var c = null")
	c.Test(
		"var a = {}; a.b = {}; var c = a.b; use(c);", //
		"var a$b = {}; var c = null; use(a$b);")
	c.TestSame("var a = {}; /** @nocollapse */ a.b;")
}

// Java: InlineAndCollapsePropertiesTest.testCollapseKeepsSourceInfoForAliases()
func (c *InlineAndCollapsePropertiesTest) TestCollapseKeepsSourceInfoForAliases() {
	c.Test(
		`var a = {};
a.use = function(arg) {};
a.b = {};
var c = a.b;
a.use(c);
`,
		`var a$use = function (arg) { };
var a$b = {};
var c = null;
a$use(a$b);
`)
	scriptNode := c.GetLastCompiler().
		GetRoot().
		GetLastChild(). // sources root
		GetOnlyChild()
	// only one source file
	AssertNode(scriptNode).IsScript().HasXChildren(4)
	statement1 := scriptNode.GetFirstChild()
	statement2 := statement1.GetNext()
	statement3 := statement2.GetNext()
	statement4 := statement3.GetNext()
	scriptName := scriptNode.GetSourceFileName()
	// Original line 2: `a.use = function(arg) {};`
	// Compiled line 1: `var a$use = function(arg) {};`
	AssertNode(statement1).
		IsVar().
		HasSourceFileName(scriptName).
		HasLineno(2).
		HasCharno(0).
		HasOneChildThat(). // `a$use = function(arg) {};
		IsName("a$use").
		HasSourceFileName(scriptName).
		HasLineno(2).
		HasCharno(0).      // `a.use = ` was the beginning of the line
		HasOneChildThat(). // `function(arg) {}`
		IsFunction().
		HasSourceFileName(scriptName).
		HasLineno(2).
		HasCharno(8)
	// `a.use = ` 8 chars
	// Original line 3: `a.b = {};`
	// Compiled line 2: `var a$b = {};`
	AssertNode(statement2).
		IsVar().
		HasSourceFileName(scriptName).
		HasLineno(3).
		HasCharno(0).
		HasOneChildThat(). // a$b = {}
		IsName("a$b").
		HasSourceFileName(scriptName).
		HasLineno(3).
		HasCharno(0).
		HasOneChildThat(). // {}
		IsObjectLit().
		HasSourceFileName(scriptName).
		HasLineno(3).
		HasCharno(6)
	// Original line 4: `var c = a.b;`
	// Compiled line 3: `var c = null;`
	AssertNode(statement3).
		IsVar().
		HasSourceFileName(scriptName).
		HasLineno(4).
		HasCharno(0).
		HasOneChildThat(). // c = a.b
		IsName("c").
		HasSourceFileName(scriptName).
		HasLineno(4).
		HasCharno(4).
		HasOneChildThat(). // null
		IsNullNode().
		HasSourceFileName(scriptName).
		HasLineno(4).
		HasCharno(10)
	// Original line 5: `a.use(c);`
	// Compiled line 4: `a$use(a$b);`
	callNodeSubject := AssertNode(statement4).
		IsExprResult().
		HasSourceFileName(scriptName).
		HasLineno(5).
		HasCharno(0).
		HasOneChildThat(). // a.use(c)
		IsCall().
		HasSourceFileName(scriptName).
		HasLineno(5).
		HasCharno(0)
	// `a$use` from `a.use(c);`
	callNodeSubject.
		HasFirstChildThat().
		IsName("a$use").
		HasSourceFileName(scriptName).
		HasLineno(5).
		HasCharno(2)
	// source info points to the `use` part of `a.use`
	// Original line 5: `c` in `a.use(c)`
	// Compiled line 4: `a$b` in `a$use(a$b)`
	callNodeSubject.
		HasSecondChildThat().
		IsName("a$b").
		HasSourceFileName(scriptName).
		HasLineno(5).
		HasCharno(6)
	// `a.use(` is 6 chars
}

// Java: InlineAndCollapsePropertiesTest.testOptChainPreventsInlineAndCollapse()
func (c *InlineAndCollapsePropertiesTest) TestOptChainPreventsInlineAndCollapse() {
	c.TestSame(
		`var a = {};
a.b = {};
// c is not really an alias due to optional chain,
// and optional chain prevents collapse of a.b.
var c = a?.b;
use(c);
`)
	c.Test(
		`var a = {};
a.b = {};
var b = a.b; // can be inlined and collapsed
b.c = {};
var c = b?.c; // opt chain prevents both
use(c);
`,
		`var a$b = {};
var b = null;
a$b.c = {};
var c = a$b?.c;
use(c);
`)
	c.Test(
		`const a = {};
a.b = {};
a.b.c = {};
const {b} = a; // can be inlined and collapsed
const c = b?.c; // opt chain prevents both
use(c);
`,
		`var a$b = {};
a$b.c = {};
const b = null;
const c = a$b?.c;
use(c);
`)
	c.Test(
		`const ns = {};
ns.Y = {};
ns.Y.prop = 3;
const {Y} = ns;
use(Y);
`,
		`var ns$Y = {};
ns$Y.prop = 3;
const Y = null;
use(ns$Y);
`)
	c.Test(
		`var a = {};
a.b = {};
a.b.c = {};
var c = a.b?.c;
use(c);
`,
		`var a$b = {};
a$b.c = {};
var c = a$b?.c;
use(c);
`)
	c.Test(
		`const ns = {};
/** @constructor */
ns.Y = function() {};
ns.Y.prop = 3;
const {Y} = ns;
use(Y);
`,
		`/** @constructor */
var ns$Y = function() {};
var ns$Y$prop = 3;
const Y = null;
use(ns$Y);
`)
	c.Test(
		`const ns = {};
/** @constructor */
ns.Y = function() {};
ns.Y.prop = 3;
const {Y} = ns;
use(Y?.prop);
`,
		`/** @constructor */
var ns$Y = function() {};
var ns$Y$prop = 3;
const Y = null;
use(ns$Y?.prop);
`)
}

// Java: InlineAndCollapsePropertiesTest.testObjLitDeclaration()
func (c *InlineAndCollapsePropertiesTest) TestObjLitDeclaration() {
	c.Test(
		"var a = {b: {}, c: {}}; var d = a.b; var e = a.c", //
		"var d = null; var e = null;")
	c.Test(
		"var a = {b: {}, c: {}}; var d = a.b; var e = a.c; use(d, e);", //
		"var a$b = {}; var a$c = {}; var d = null; var e = null; use(a$b, a$c);")
	c.Test(
		"var a = {b: {}, /** @nocollapse */ c: {}}; var d = a.b; var e = a.c", //
		"var a = {/** @nocollapse */ c: {}}; var d = null; var e = null;")
	c.Test(
		"var a = {b: {}, /** @nocollapse */ c: {}}; var d = a.b; var e = a.c; use(d, e);", //
		"var a$b = {};var a = {/** @nocollapse */ c: {}};var d = null;var e = null;use(a$b, a.c);")
}

// Java: InlineAndCollapsePropertiesTest.testObjLitDeclarationWithGet1()
func (c *InlineAndCollapsePropertiesTest) TestObjLitDeclarationWithGet1() {
	c.TestSame("var a = {get b(){}};")
}

// Java: InlineAndCollapsePropertiesTest.testObjLitDeclarationWithGet2()
func (c *InlineAndCollapsePropertiesTest) TestObjLitDeclarationWithGet2() {
	c.Test(
		"var a = {b: {}, get c(){}}; var d = a.b; var e = a.c;", //
		"var a = {get c() {}}; var d=null; var e = a.c")
	c.Test(
		"var a = {b: {}, get c(){}}; var d = a.b; var e = a.c; use(d);", //
		"var a$b = {};var a = {get c(){}};var d = null; var e = a.c; use(a$b)")
}

// Java: InlineAndCollapsePropertiesTest.testObjLitDeclarationWithGet3()
func (c *InlineAndCollapsePropertiesTest) TestObjLitDeclarationWithGet3() {
	c.Test(
		"var a = {b: {get c() { return 3; }}};", //
		"var a$b = {get c() { return 3; }};")
}

// Java: InlineAndCollapsePropertiesTest.testObjLitDeclarationWithSet1()
func (c *InlineAndCollapsePropertiesTest) TestObjLitDeclarationWithSet1() {
	c.TestSame("var a = {set b(a){}};")
}

// Java: InlineAndCollapsePropertiesTest.testObjLitDeclarationWithSet2()
func (c *InlineAndCollapsePropertiesTest) TestObjLitDeclarationWithSet2() {
	c.Test(
		"var a = {b: {}, set c(a){}}; var d = a.b; var e = a.c", //
		"var a = {set c(a$jscomp$1){}}; var d=null; var e = a.c")
	c.Test(
		"var a = {b: {}, set c(a){}}; var d = a.b; var e = a.c; use(d);", //
		"var a$b = {}; var a = {set c(a$jscomp$1){}}; var d=null; var e = a.c; use(a$b);")
}

// Java: InlineAndCollapsePropertiesTest.testObjLitDeclarationWithSet3()
func (c *InlineAndCollapsePropertiesTest) TestObjLitDeclarationWithSet3() {
	c.Test(
		"var a = {b: {set c(d) {}}};", //
		"var a$b = {set c(d) {}};")
}

// Java: InlineAndCollapsePropertiesTest.testObjLitDeclarationWithGetAndSet1()
func (c *InlineAndCollapsePropertiesTest) TestObjLitDeclarationWithGetAndSet1() {
	c.Test(
		"var a = {b: {get c() { return 3; },set c(d) {}}};", //
		"var a$b = {get c() { return 3; },set c(d) {}};")
}

// Java: InlineAndCollapsePropertiesTest.testObjLitAssignmentDepth3()
func (c *InlineAndCollapsePropertiesTest) TestObjLitAssignmentDepth3() {
	c.Test(
		"var a = {}; a.b = {}; a.b.c = {d: 1, e: 2}; var f = a.b.c.d", //
		"var a$b$c$d = 1; var a$b$c$e = 2; var f = null")
	c.Test(
		"var a = {}; a.b = {}; a.b.c = {d: 1, e: 2}; var f = a.b.c.d; use(f);", //
		"var a$b$c$d = 1; var a$b$c$e = 2; var f = null; use(a$b$c$d);")
	c.Test(
		`var a = {}; a.b = {}; a.b.c = {d: 1, /** @nocollapse */ e: 2};
var f = a.b.c.d; var g = a.b.c.e
`,
		"var a$b$c$d = 1; var a$b$c = {/** @nocollapse */ e: 2}; var f = null; var g = null;")
	c.Test(
		`var a = {}; a.b = {}; a.b.c = {d: 1, /** @nocollapse */ e: 2};
var f = a.b.c.d; var g = a.b.c.e; use(f, g);
`,
		`var a$b$c$d = 1; var a$b$c = { /** @nocollapse */ e: 2};
var f = null; var g = null; use(a$b$c$d, a$b$c.e);
`)
	c.TestSame(
		`var a = {}; /** @nocollapse*/ a.b = {};
a.b.c = {d: 1, e: 2};
var f = null; var g = null;
`)
}

// Java: InlineAndCollapsePropertiesTest.testObjLitAssignmentDepth4()
func (c *InlineAndCollapsePropertiesTest) TestObjLitAssignmentDepth4() {
	c.Test(
		"var a = {}; a.b = {}; a.b.c = {}; a.b.c.d = {e: 1, f: 2}; var g = a.b.c.d.e;", //
		"var a$b$c$d$e = 1; var a$b$c$d$f = 2; var g = null;")
	c.Test(
		"var a = {}; a.b = {}; a.b.c = {}; a.b.c.d = {e: 1, f: 2}; var g = a.b.c.d.e; use(g);", //
		"var a$b$c$d$e = 1; var a$b$c$d$f = 2; var g = null; use(a$b$c$d$e);")
}

// Java: InlineAndCollapsePropertiesTest.testHasOwnProperty()
func (c *InlineAndCollapsePropertiesTest) TestHasOwnProperty() {
	c.Test(
		"var a = {'b': 1, 'c': 1}; var alias = a;    alert(alias.hasOwnProperty('c'));", //
		"var a = {'b': 1, 'c': 1}; var alias = null; alert(a.hasOwnProperty('c'));")
}

// Java: InlineAndCollapsePropertiesTest.testAliasCreatedForObjectDepth1_1()
func (c *InlineAndCollapsePropertiesTest) TestAliasCreatedForObjectDepth1_1() {
	// An object's properties are not collapsed if the object is referenced
	// in a such a way that an alias is created for it, if that alias is used.
	c.Test(
		"var a = {b: 0}; var c = a; c.b = 1; a.b == c.b;", //
		"var a$b = 0; var c = null; a$b = 1; a$b == a$b;")
	c.Test(
		"var a = {b: 0}; var c =    a; c.b = 1; a?.b == c?.b;",
		"var a = {b: 0}; var c = null; a.b = 1; a?.b == a?.b;")
	c.Test(
		"var a = {b: 0}; var c = a; c.b = 1; a.b == c.b; use(c);", //
		"var a={b:0}; var c=null; a.b=1; a.b == a.b; use(a);")
}

// Java: InlineAndCollapsePropertiesTest.testAliasCreatedForObjectDepth1_2()
func (c *InlineAndCollapsePropertiesTest) TestAliasCreatedForObjectDepth1_2() {
	c.TestSame("var a = {b: 0}; f(a); a.b;")
}

// Java: InlineAndCollapsePropertiesTest.testAliasCreatedForObjectDepth1_3()
func (c *InlineAndCollapsePropertiesTest) TestAliasCreatedForObjectDepth1_3() {
	c.TestSame("var a = {b: 0}; new f(a); a.b;")
}

// Java: InlineAndCollapsePropertiesTest.testMisusedConstructorTag()
func (c *InlineAndCollapsePropertiesTest) TestMisusedConstructorTag() {
	c.Test(
		`var a = {}; var d = a; a.b = function() {};
/** @constructor */ a.b.c = 0; a.b.c;
`,
		"var d=null; var a$b=function(){}; /** @constructor */ var a$b$c=0; a$b$c;")
}

// Java: InlineAndCollapsePropertiesTest.testAliasCreatedForCtorDepth1_1()
func (c *InlineAndCollapsePropertiesTest) TestAliasCreatedForCtorDepth1_1() {
	// A constructor's properties *are* collapsed even if the function is
	// referenced in a such a way that an alias is created for it,
	// since a function with custom properties is considered a class and its
	// non-prototype properties are considered static methods and variables.
	// People don't typically iterate through static members of a class or
	// refer to them using an alias for the class name.
	c.Test(
		`/** @constructor */ var a = function(){}; a.b = 1;
var c = a; c.b = 2; a.b == c.b;
`,
		`/** @constructor */ var a = function(){}; var a$b = 1;
var c = null; a$b = 2; a$b == a$b;
`)
	c.Test(
		`/** @constructor */ var a = function(){}; a.b = 1;
var c = a; c.b = 2; a?.b == c?.b;
`,
		`/** @constructor */ var a = function(){}; var a$b = 1;
var c = null; a$b = 2; a?.b == a?.b;
`)
	// Sometimes we want to prevent static members of a constructor from
	// being collapsed.
	c.Test(
		`/** @constructor */ var a = function(){};
/** @nocollapse */ a.b = 1; var c = a; c.b = 2; a.b == c.b;
`,
		`/** @constructor */ var a = function(){};
/** @nocollapse */ a.b = 1; var c = null; a.b = 2; a.b == a.b;
`)
	c.Test(
		`/** @constructor */ var a = function(){};
/** @nocollapse */ a.b = 1; var c = a; c.b = 2; a?.b == c?.b;
`,
		`/** @constructor */ var a = function(){};
/** @nocollapse */ a.b = 1; var c = null; a.b = 2; a?.b == a?.b;
`)
}

// Java: InlineAndCollapsePropertiesTest.testAliasCreatedForFunctionDepth2()
func (c *InlineAndCollapsePropertiesTest) TestAliasCreatedForFunctionDepth2() {
	c.Test(
		"var a = {}; a.b = function() {}; a.b.c = 1; var d = a.b; a.b.c != d.c;", //
		"var a$b = function() {}; var a$b$c = 1; var d = null; a$b$c != a$b$c;")
	c.TestSame("var a = {}; a.b = function() {}; a.b.c = 1; var d = a?.b; a.b?.c != d?.c;")
	c.Test(
		`var a = {}; a.b = function() {}; /** @nocollapse */ a.b.c = 1;
var d = a.b; a.b.c == d.c;
`,
		"var a$b = function() {}; /** @nocollapse */ a$b.c = 1; var d = null; a$b.c == a$b.c;")
}

// Java: InlineAndCollapsePropertiesTest.testAliasCreatedForCtorDepth2()
func (c *InlineAndCollapsePropertiesTest) TestAliasCreatedForCtorDepth2() {
	c.Test(
		`var a = {}; /** @constructor */ a.b = function() {}; a.b.c = 1; var d = a.b;
a.b.c == d.c;
`,
		`/** @constructor */ var a$b = function() {}; var a$b$c = 1; var d = null;
a$b$c == a$b$c;
`)
	c.TestSame(
		c.Srcs(
			`var a = {}; /** @constructor */ a.b = function() {}; a.b.c = 1; var d = a?.b;
a.b?.c == d?.c;
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
	c.Test(
		`var a = {}; /** @constructor */ a.b = function() {};
/** @nocollapse */ a.b.c = 1; var d = a.b;
a.b.c == d.c;
`,
		`/** @constructor */ var a$b = function() {};
/** @nocollapse */ a$b.c = 1; var d = null;
a$b.c == a$b.c;
`)
}

// Java: InlineAndCollapsePropertiesTest.testAliasCreatedForClassDepth1_1()
func (c *InlineAndCollapsePropertiesTest) TestAliasCreatedForClassDepth1_1() {
	c.Test(
		`var a = {}; /** @constructor */ a.b = function(){};
var c = a; c.b = 0; a.b == c.b;
`,
		`/** @constructor */ var a$b = function(){};
var c = null; a$b = 0; a$b == a$b;
`)
	c.Test(
		c.Srcs(
			"var a = {}; /** @constructor */ a.b = function(){};\n"+
				"// `a?.b` and `c?.b` are aliasing gets\n"+
				"var c = a; c.b = 0; a?.b == c?.b;\n"),
		c.Expected(
			"var a = {}; /** @constructor */ a.b = function(){};\n"+
				"// `a?.b` and `c?.b` are aliasing gets\n"+
				"var c = null; a.b = 0; a?.b == a?.b;\n"),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
	c.TestSame(
		`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
var c = 1; c = a; c.b = 0; a.b == c.b;
`,
		InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING)
	c.TestSame(
		`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
var c = 1; c = a; c.b = 0; a?.b == c?.b;
`,
		InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING)
	c.Test(
		`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
var c = a; c.b = 0; a.b == c.b;
`,
		`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
var c = null; a.b = 0; a.b == a.b;
`)
	c.Test(
		c.Srcs(
			`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
var c = a; c.b = 0; a?.b == c?.b;
`),
		c.Expected(
			`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
var c = null; a.b = 0; a?.b == a?.b;
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
	c.Test(
		`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
var c = a; c.b = 0; a.b == c.b; use(c);
`,
		`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
var c = null; a.b = 0; a.b == a.b; use(a);
`,
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
	c.Test(
		`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
var c = a; c.b = 0; a?.b == c?.b; use(c);
`,
		`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
var c = null; a.b = 0; a?.b == a?.b; use(a);
`,
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: InlineAndCollapsePropertiesTest.testObjLitDeclarationUsedInSameVarList()
func (c *InlineAndCollapsePropertiesTest) TestObjLitDeclarationUsedInSameVarList() {
	// The collapsed properties must be defined in the same place in the var list
	// where they were originally defined (and not, for example, at the end).
	c.Test(
		"var a = {b: {}, c: {}}; var d = a.b; var e = a.c; use(d, e);", //
		"var a$b = {}; var a$c = {}; var d = null; var e = null; use(a$b, a$c);")
	c.Test(
		"var a = {b: {}, /** @nocollapse */ c: {}}; var d = a.b; var e = a.c; use(d, e);", //
		"var a$b = {};var a = {/** @nocollapse */ c: {}};var d = null;var e = null;use(a$b, a.c);")
}

// Java: InlineAndCollapsePropertiesTest.testAddPropertyToUncollapsibleObjectInLocalScopeDepth1()
func (c *InlineAndCollapsePropertiesTest) TestAddPropertyToUncollapsibleObjectInLocalScopeDepth1() {
	c.Test(
		"var a = {}; var c = a; use(c); (function() {a.b = 0;})(); a.b;", //
		"var a={}; var c=null; use(a); (function(){ a.b = 0; })(); a.b;")
	c.Test(
		"var a = {}; var c = 1; c = a; (function() {a.b = 0;})(); a.b;", //
		"var a = {}; var c=1; c = a; (function(){a.b = 0;})(); a.b;")
}

// Java: InlineAndCollapsePropertiesTest.testAddPropertyToUncollapsibleNamedCtorInLocalScopeDepth1()
func (c *InlineAndCollapsePropertiesTest) TestAddPropertyToUncollapsibleNamedCtorInLocalScopeDepth1() {
	c.Test(
		`/** @constructor */ function a() {} var a$b; var c = a;
(function() {a$b = 0;})(); a$b;
`,
		`/** @constructor */ function a() {} var a$b; var c = null;
(function() {a$b = 0;})(); a$b;
`)
}

// Java: InlineAndCollapsePropertiesTest.testAddPropertyToUncollapsibleCtorInLocalScopeDepth1()
func (c *InlineAndCollapsePropertiesTest) TestAddPropertyToUncollapsibleCtorInLocalScopeDepth1() {
	c.Test(
		`/** @constructor */ var a = function() {}; var c = a;
(function() {a.b = 0;})(); a.b;
`,
		`/** @constructor */ var a = function() {}; var a$b;
var c = null; (function() {a$b = 0;})(); a$b;
`)
}

// Java: InlineAndCollapsePropertiesTest.testAddPropertyToUncollapsibleObjectInLocalScopeDepth2()
func (c *InlineAndCollapsePropertiesTest) TestAddPropertyToUncollapsibleObjectInLocalScopeDepth2() {
	c.Test(
		`var a = {}; a.b = {}; var d = a.b; use(d);
(function() {a.b.c = 0;})(); a.b.c;
`,
		`var a$b = {}; var d = null; use(a$b);
(function() {a$b.c = 0;})(); a$b.c;
`)
}

// Java: InlineAndCollapsePropertiesTest.testAddPropertyToUncollapsibleCtorInLocalScopeDepth2()
func (c *InlineAndCollapsePropertiesTest) TestAddPropertyToUncollapsibleCtorInLocalScopeDepth2() {
	c.Test(
		`var a = {}; /** @constructor */ a.b = function (){}; var d = a.b;
(function() {a.b.c = 0;})(); a.b.c;
`,
		`/** @constructor */ var a$b = function (){}; var a$b$c; var d = null;
(function() {a$b$c = 0;})(); a$b$c;
`)
}

// Java: InlineAndCollapsePropertiesTest.testPropertyOfChildFuncOfUncollapsibleObjectDepth1()
func (c *InlineAndCollapsePropertiesTest) TestPropertyOfChildFuncOfUncollapsibleObjectDepth1() {
	c.Test(
		"var a = {}; var c = a; a.b = function (){}; a.b.x = 0; a.b.x;", //
		"var c = null; var a$b=function() {}; var a$b$x = 0; a$b$x;")
	c.Test(
		"var a = {}; var c = a; a.b = function (){}; a.b.x = 0; a.b.x; use(c);", //
		"var a = {}; var c = null; a.b=function() {}; a.b.x = 0; a.b.x; use(a);")
}

// Java: InlineAndCollapsePropertiesTest.testPropertyOfChildFuncOfUncollapsibleObjectDepth2()
func (c *InlineAndCollapsePropertiesTest) TestPropertyOfChildFuncOfUncollapsibleObjectDepth2() {
	c.Test(
		"var a = {}; a.b = {}; var c = a.b; a.b.c = function (){}; a.b.c.x = 0; a.b.c.x;", //
		"var c=null; var a$b$c = function(){}; var a$b$c$x = 0; a$b$c$x")
	c.Test(
		`var a = {}; a.b = {}; var c = a.b; a.b.c = function (){}; a.b.c.x = 0; a.b.c.x;
 use(c);
`, //
		"var a$b = {}; var c=null; a$b.c = function(){}; a$b.c.x=0; a$b.c.x; use(a$b);")
}

// Java: InlineAndCollapsePropertiesTest.testAddPropertyToChildFuncOfUncollapsibleObjectInLocalScope()
func (c *InlineAndCollapsePropertiesTest) TestAddPropertyToChildFuncOfUncollapsibleObjectInLocalScope() {
	c.Test(
		`var a = {}; a.b = function (){}; a.b.x = 0;
var c = a; (function() {a.b.y = 1;})(); a.b.x; a.b.y;
`,
		`var a$b=function() {}; var a$b$y; var a$b$x = 0; var c=null;
(function() { a$b$y=1; })(); a$b$x; a$b$y
`)
}

// Java: InlineAndCollapsePropertiesTest.testAddPropertyToChildTypeOfUncollapsibleObjectInLocalScope()
func (c *InlineAndCollapsePropertiesTest) TestAddPropertyToChildTypeOfUncollapsibleObjectInLocalScope() {
	c.Test(
		`var a = {};
/** @constructor */
a.b = function (){};
a.b.x = 0;
var c = a;
(function() {a.b.y = 1;})();
a.b.x;
a.b.y;
`,
		`/** @constructor */
var a$b = function (){};
var a$b$y;
var a$b$x = 0;
var c = null;
(function() {a$b$y = 1;})();
a$b$x;
a$b$y;
`)
	c.Test(
		`var a = {};
/** @constructor */
a.b = function (){};
a.b.x = 0;
var c = a;
(function() {a.b.y = 1;})();
a.b.x;
a.b.y;
use(c);
`,
		`var a = {};
/** @constructor */
a.b = function (){};
a.b.x = 0;
var c = null;
(function() {a.b.y = 1;})();
a.b.x;
a.b.y;
use(a);
`,
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: InlineAndCollapsePropertiesTest.testAddPropertyToChildOfUncollapsibleFunctionInLocalScope()
func (c *InlineAndCollapsePropertiesTest) TestAddPropertyToChildOfUncollapsibleFunctionInLocalScope() {
	c.Test(
		"function a() {} a.b = {x: 0}; var c = a; (function() {a.b.y = 0;})(); a.b.y;", //
		"function a() {} var a$b$x=0; var a$b$y; var c=null; (function(){a$b$y=0})(); a$b$y")
}

// Java: InlineAndCollapsePropertiesTest.testAddPropertyToChildOfUncollapsibleCtorInLocalScope()
func (c *InlineAndCollapsePropertiesTest) TestAddPropertyToChildOfUncollapsibleCtorInLocalScope() {
	c.Test(
		`/** @constructor */ var a = function() {}; a.b = {x: 0}; var c = a;
(function() {a.b.y = 0;})(); a.b.y;
`,
		`/** @constructor */ var a = function() {}; var a$b$x = 0; var a$b$y; var c = null;
(function() {a$b$y = 0;})(); a$b$y;
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasCreatedAfterVarDeclaration1()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasCreatedAfterVarDeclaration1() {
	c.Test(
		`var a = {b: 3};
function f() {
  var tmp;
  tmp = a;
  use(tmp.b);
}
`,
		`var a$b = 3
function f() {
  var tmp;
  tmp = null;
  use(a$b);
}
`)
	c.Test(
		`var a = {b: 3};
function f() {
  var tmp;
  tmp = a;
  use?.(tmp?.b);
}
`,
		`var a = {b: 3};
function f() {
  var tmp;
  tmp = null;
  use?.(a?.b);
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasCreatedAfterVarDeclaration2()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasCreatedAfterVarDeclaration2() {
	c.Test(
		`var a = {b: 3}
function f() {
  var tmp;
  if (true) {
    tmp = a;
    use(tmp.b);
  }
}
`,
		`var a$b = 3;
function f() {
  var tmp;
  if (true) {
    tmp = null;
    use(a$b);
  }
}
`)
	c.Test(
		`var a = {b: 3}
function f() {
  var tmp;
  if (true) {
    tmp = a;
    use?.(tmp?.b);
  }
}
`,
		`var a = {b: 3}
function f() {
  var tmp;
  if (true) {
    tmp = null;
    use?.(a?.b);
  }
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasCreatedAfterVarDeclaration3()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasCreatedAfterVarDeclaration3() {
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

// Java: InlineAndCollapsePropertiesTest.testPartialLocalCtorAlias()
func (c *InlineAndCollapsePropertiesTest) TestPartialLocalCtorAlias() {
	c.TestWarning(
		`/** @constructor */ var Main = function() {};
Main.doSomething = function(i) {}
function f() {
  var tmp;
  if (g()) {
    use(tmp.doSomething);
    tmp = Main;
    tmp.doSomething(5);
  }
  use(tmp.doSomething); // can't inline this use of tmp.
}
`,
		InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING)
}

// Java: InlineAndCollapsePropertiesTest.testFunctionAlias2()
func (c *InlineAndCollapsePropertiesTest) TestFunctionAlias2() {
	c.Test(
		"var a = {}; a.b = {}; a.b.c = function(){}; a.b.d = a.b.c;use(a.b.d)", //
		"var a$b$c = function(){}; var a$b$d = null;use(a$b$c);")
	c.Test(
		"var a = {}; a.b = {}; a.b.c = function(){}; a.b.d = a.b?.c; use?.(a.b?.d);", //
		"var         a$b = {}; a$b.c = function(){}; a$b.d = a$b?.c; use?.(a$b?.d);")
}

// Java: InlineAndCollapsePropertiesTest.testLocalAlias1()
func (c *InlineAndCollapsePropertiesTest) TestLocalAlias1() {
	c.Test(
		"var   a = {b: 3}; function f() { var x =    a; f(x.b); }", //
		"var a$b =     3 ; function f() { var x = null; f(a$b); }")
	c.Test(
		"var   a = {b: 3}; function f() { var x =    a; f(x?.b); }",
		"var   a = {b: 3}; function f() { var x = null; f(a?.b); }")
}

// Java: InlineAndCollapsePropertiesTest.testLocalAlias2()
func (c *InlineAndCollapsePropertiesTest) TestLocalAlias2() {
	c.Test(
		"var a   = {b: 3,       c : 4}; function f() { var x =    a; f(x.b); f(x.c);}",
		"var a$b =     3; var a$c = 4 ; function f() { var x = null; f(a$b); f(a$c);}")
	c.Test(
		"var a = {b: 3, c: 4}; function f() { var x =    a; f?.(x?.b); f?.(x?.c);}",
		"var a = {b: 3, c: 4}; function f() { var x = null; f?.(a?.b); f?.(a?.c);}")
}

// Java: InlineAndCollapsePropertiesTest.testLocalAlias3()
func (c *InlineAndCollapsePropertiesTest) TestLocalAlias3() {
	c.Test(
		`var a = {b: 3, c: {d: 5}};
function f() { var x = a; f(x.b); f(x.c); f(x.c.d); }
`,
		`var a$b = 3; var a$c = {d: 5};
function f() { var x = null; f(a$b); f(a$c); f(a$c.d);}
`)
	c.Test(
		`var a = {b: 3, c: {d: 5}};
function f() { var x =    a; f?.(x?.b); f?.(x?.c); f?.(x.c?.d);}
`,
		`var a = {b: 3, c: {d:5}};
function f() { var x = null; f?.(a?.b); f?.(a?.c); f?.(a.c?.d);}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAlias4()
func (c *InlineAndCollapsePropertiesTest) TestLocalAlias4() {
	c.Test(
		`var a = {b: 3}; var c = {d: 5};
function f() { var x = a; var y = c; f(x.b); f(y.d); }
`,
		`var a$b = 3; var c$d = 5;
function f() { var x = null; var y = null; f(a$b); f(c$d);}
`)
	c.Test(
		`var a = {b: 3}; var c = {d: 5};
function f() { var x =    a; var y =    c; f?.(x?.b); f?.(y?.d); }
`,
		`var a = {b: 3}; var c = {d: 5};
function f() { var x = null; var y = null; f?.(a?.b); f?.(c?.d); }
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAlias5()
func (c *InlineAndCollapsePropertiesTest) TestLocalAlias5() {
	c.Test(
		`var a = {b: {c: 5}};
function f() { var x = a; var y = x.b; f(a.b.c); f(y.c); }
`,
		`var a$b$c = 5;
function f() { var x = null; var y = null; f(a$b$c); f(a$b$c);}
`)
	c.Test(
		`var a = {b: {c: 5}};
function f() { var x =    a; var y = x?.b; f?.(a.b?.c); f?.(y?.c); }
`,
		`var a = {b: {c: 5}};
function f() { var x = null; var y = a?.b; f?.(a.b?.c); f?.(y?.c); }
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAlias6()
func (c *InlineAndCollapsePropertiesTest) TestLocalAlias6() {
	c.Test(
		"var a = {b: 3}; function f() { var x = a; if (x.b) { f(x.b); } }", //
		"var a$b = 3; function f() { var x = null; if (a$b) { f(a$b); } }")
	c.Test(
		"var a = {b: 3}; function f() { var x =    a; if (x?.b) { f?.(x?.b); } }", //
		"var a = {b: 3}; function f() { var x = null; if (a?.b) { f?.(a?.b); } }")
}

// Java: InlineAndCollapsePropertiesTest.testLocalAlias7()
func (c *InlineAndCollapsePropertiesTest) TestLocalAlias7() {
	c.Test(
		"var a = {b: {c: 5}}; function f() { var x = a.b; f(x.c); }", //
		"var a$b$c = 5; function f() { var x = null; f(a$b$c); }")
	c.TestSame("var a = {b: {c: 5}}; function f() { var x = a?.b; f?.(x?.c); }")
}

// Java: InlineAndCollapsePropertiesTest.testLocalAlias8()
func (c *InlineAndCollapsePropertiesTest) TestLocalAlias8() {
	c.TestSame("var a = { b: 3 }; function f() { if (true) { var x = a; f(x.b); } x = { b : 4}; }")
}

// Java: InlineAndCollapsePropertiesTest.testGlobalWriteToAncestor()
func (c *InlineAndCollapsePropertiesTest) TestGlobalWriteToAncestor() {
	c.TestSame("var a = {b: 3}; function f() { var x = a; f(a.b); } a = 5;")
}

// Java: InlineAndCollapsePropertiesTest.testGlobalWriteToNonAncestor()
func (c *InlineAndCollapsePropertiesTest) TestGlobalWriteToNonAncestor() {
	c.Test(
		"var a = {b: 3}; function f() { var x =    a; f(a.b); } a.b = 5;", //
		"var    a$b = 3; function f() { var x = null; f(a$b); } a$b = 5;")
	c.Test(
		"var a = {b: 3}; function f() { var x =    a; f(a?.b); } a.b = 5;",
		"var a = {b: 3}; function f() { var x = null; f(a?.b); } a.b = 5;")
}

// Java: InlineAndCollapsePropertiesTest.testLocalWriteToAncestor()
func (c *InlineAndCollapsePropertiesTest) TestLocalWriteToAncestor() {
	c.TestSame("var a = {b: 3}; function f() { a = 5; var x = a; f(a.b); } ")
}

// Java: InlineAndCollapsePropertiesTest.testLocalWriteToNonAncestor()
func (c *InlineAndCollapsePropertiesTest) TestLocalWriteToNonAncestor() {
	c.Test(
		`var a = {b: 3};
function f() { a.b = 5; var x = a; f(a.b); }
`,
		"var a$b = 3; function f() { a$b = 5; var x = null; f(a$b); } ")
	c.Test(
		`var a = {b: 3};
function f() { a.b = 5; var x =    a; f?.(a?.b); }
`,
		`var a = {b: 3};
function f() { a.b = 5; var x = null; f?.(a?.b); }
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasOfAncestor()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasOfAncestor() {
	c.TestSame(
		`var a = {b: {c: 5}}; function g() { f(a); }
function f() { var x = a.b; f(x.c); }
`)
}

// Java: InlineAndCollapsePropertiesTest.testGlobalAliasOfAncestor()
func (c *InlineAndCollapsePropertiesTest) TestGlobalAliasOfAncestor() {
	c.Test(
		"var a = {b: {c: 5}}; function f() { var x = a.b; f(x.c); }", //
		"var a$b$c=5; function f() {var x=null; f(a$b$c); }")
	c.TestSame("var a = {b: {c: 5}}; function f() { var x = a?.b; f?.(x?.c); }")
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasOfOtherName()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasOfOtherName() {
	c.TestSame(
		`var foo = function() { return {b: 3}; };
var a = foo(); a.b = 5;
function f() { var x = a.b; f(x); }
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasOfFunction()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasOfFunction() {
	c.Test(
		`var a = function() {}; a.b = 5;
function f() { var x = a.b; f(x); }
`,
		`var a = function() {}; var a$b = 5;
function f() { var x = null; f(a$b); }
`)
}

// Java: InlineAndCollapsePropertiesTest.testNoInlineGetpropIntoCall()
func (c *InlineAndCollapsePropertiesTest) TestNoInlineGetpropIntoCall() {
	c.Test(
		"var b = x; function f() { var a = b; a(); }", //
		"var b = x; function f() { var a = null; b(); }")
	c.Test(
		"var b = {}; b.c = x; function f() { var a = b.c; a(); }", //
		"var b$c = x; function f() { var a = null; b$c(); }")
}

// Java: InlineAndCollapsePropertiesTest.testCommaOperator()
func (c *InlineAndCollapsePropertiesTest) TestCommaOperator() {
	c.Test(
		"var a = {}; a.b = function() {}, a.b();", //
		"var a$b; a$b=function() {}, a$b();")
	c.Test(
		`var ns = {};
ns.Foo = {};
var Baz = {};
Baz.Foo = ns.Foo;
(Baz.Foo.bar = 10, 123);
`,
		`var Baz$Foo=null;
var ns$Foo$bar;
(ns$Foo$bar = 10, 123);
`)
	c.Test(
		`var ns = {};
ns.Foo = {};
var Baz = {};
Baz.Foo = ns?.Foo;
(Baz.Foo.bar = 10, 123);
`,
		`var ns = {};
ns.Foo = {};
var Baz$Foo = ns?.Foo;
(Baz$Foo.bar = 10, 123);
`)
	c.Test(
		`var ns = {};
ns.Foo = {};
var Baz = {};
Baz.Foo = ns.Foo;
function f() { (Baz.Foo.bar = 10, 123); }
`,
		`var ns$Foo$bar;
var Baz$Foo=null;
function f() { (ns$Foo$bar = 10, 123); }
`)
	c.Test(
		`var ns = {};
ns.Foo = {};
var Baz = {};
Baz.Foo = ns?.Foo;
function f() { (Baz.Foo.bar = 10, 123); }
`,
		`var ns = {};
ns.Foo = {};
var Baz$Foo = ns?.Foo;
function f() { (Baz$Foo.bar = 10, 123); }
`)
}

// Java: InlineAndCollapsePropertiesTest.testTypeDefAlias1()
func (c *InlineAndCollapsePropertiesTest) TestTypeDefAlias1() {
	c.Test(
		`/** @constructor */ var D = function() {};
/** @constructor */ D.L = function() {};
/** @type {D.L} */ D.L.A = new D.L();

/** @const */ var M = {};
/** @typedef {D.L} */ M.L = D.L;

use(M.L.A);
`,
		`/** @constructor */ var D = function() {};
/** @constructor */ var D$L = function() {};
/** @type {D.L} */ var D$L$A = new D$L();
/** @typedef {D.L} */ var M$L = null
use(D$L$A);
`)
	c.Test(
		`/** @constructor */ var D = function() {};
/** @constructor */ D.L = function() {};
/** @type {D.L} */ D.L.A = new D.L();

/** @const */ var M = {};
/** @typedef {D.L} */ M.L = D.L;

use?.(M.L?.A);
`,
		`/** @constructor */ var D = function() {};
/** @constructor */ var D$L = function() {};
/** @type {D.L} */ var D$L$A = new D$L();
/** @typedef {D.L} */ var M$L = null
// TODO(b/148237949): collapse above breaks this reference
use?.(D$L?.A);
`)
}

// Java: InlineAndCollapsePropertiesTest.testGlobalAliasWithProperties1()
func (c *InlineAndCollapsePropertiesTest) TestGlobalAliasWithProperties1() {
	c.Test(
		`var ns = {};
/** @constructor */ ns.Foo = function() {};
/** @enum {number} */ ns.Foo.EventType = {A:1, B:2};
/** @constructor */ ns.Bar = ns.Foo;
var x = function() {use(ns.Bar.EventType.A)};
use(x);
`,
		`/** @constructor */ var ns$Foo = function(){};
var ns$Foo$EventType$A = 1;
var ns$Foo$EventType$B = 2;
/** @constructor */ var ns$Bar = null;
var x = function(){use(ns$Foo$EventType$A)};
use(x);
`)
	c.TestSame(
		c.Srcs(
			`var ns = {};
/** @constructor */ ns.Foo = function() {};
/** @enum {number} */ ns.Foo.EventType = {A:1, B:2};
/** @constructor */ ns.Bar = ns?.Foo;
var x = function() {use?.(ns.Bar.EventType?.A)};
use?.(x);
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: InlineAndCollapsePropertiesTest.testGlobalAliasWithProperties2()
func (c *InlineAndCollapsePropertiesTest) TestGlobalAliasWithProperties2() {
	// Reassignment of properties was necessary to prevent invalid code in
	// previous iterations of this optimization.  Verify we don't break
	// code like this.  Now it causes a back-off of the collapsing because
	// the value is assigned more than once.
	c.Test(
		`var ns = {};
/** @constructor */ ns.Foo = function() {};
/** @enum {number} */ ns.Foo.EventType = {A:1, B:2};
/** @constructor */ ns.Bar = ns.Foo;
/** @enum {number} */ ns.Bar.EventType = ns.Foo.EventType;
var x = function() {use(ns.Bar.EventType.A)};
use(x)
`,
		`/** @constructor */ var ns$Foo = function(){};
/** @enum {number} */ var ns$Foo$EventType = {A:1, B:2};
/** @constructor */ var ns$Bar = null;
/** @enum {number} */ ns$Foo$EventType = ns$Foo$EventType;
var x = function(){use(ns$Foo$EventType.A)};
use(x);
`)
	c.TestSame(
		c.Srcs(
			`var ns = {};
/** @constructor */ ns.Foo = function() {};
/** @enum {number} */ ns.Foo.EventType = {A:1, B:2};
/** @constructor */ ns.Bar = ns?.Foo;
/** @enum {number} */ ns.Bar.EventType = ns.Foo?.EventType;
var x = function() {use?.(ns.Bar.EventType?.A)};
use?.(x)
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: InlineAndCollapsePropertiesTest.testGlobalAliasWithProperties3()
func (c *InlineAndCollapsePropertiesTest) TestGlobalAliasWithProperties3() {
	c.Test(
		`var ns = {};
/** @constructor */ ns.Foo = function() {};
/** @enum {number} */ ns.Foo.EventType = {A:1, B:2};
/** @constructor */ ns.Bar = ns.Foo;
/** @enum {number} */ ns.Bar.Other = {X:1, Y:2};
var x = function() {use(ns.Bar.Other.X)};
use(x)
`,
		`/** @constructor */ var ns$Foo=function(){};
var ns$Foo$EventType$A=1;
var ns$Foo$EventType$B=2;
/** @constructor */ var ns$Bar=null;
var ns$Foo$Other$X=1;
var ns$Foo$Other$Y=2;
var x=function(){use(ns$Foo$Other$X)};
use(x)
`)
	c.TestSame(
		c.Srcs(
			`var ns = {};
/** @constructor */ ns.Foo = function() {};
/** @enum {number} */ ns.Foo.EventType = {A:1, B:2};
/** @constructor */ ns.Bar = ns?.Foo;
/** @enum {number} */ ns.Bar.Other = {X:1, Y:2};
var x = function() {use?.(ns.Bar.Other?.X)};
use?.(x)
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: InlineAndCollapsePropertiesTest.testGlobalAliasWithProperties4()
func (c *InlineAndCollapsePropertiesTest) TestGlobalAliasWithProperties4() {
	c.TestSame(
		`var nullFunction = function(){};
var blob = {};
blob.init = nullFunction;
use(blob)
`)
}

// Java: InlineAndCollapsePropertiesTest.testGlobalAlias_propertyOnExternedConstructor_isNotChanged()
func (c *InlineAndCollapsePropertiesTest) TestGlobalAlias_propertyOnExternedConstructor_isNotChanged() {
	c.TestSame(
		c.Externs("/** @constructor */ var blob = function() {}"),
		c.Srcs(
			`var nullFunction = function(){};
blob.init = nullFunction;
use(blob.init)
`))
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasOfEnumWithInstanceofCheck()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasOfEnumWithInstanceofCheck() {
	c.Test(
		`/** @constructor */
var Enums = function() {
};

/** @enum {number} */
Enums.Fruit = {
 APPLE: 1,
 BANANA: 2,
};

function foo(f) {
 if (f instanceof Enums) { alert('what?'); return; }

 var Fruit = Enums.Fruit;
 if (f == Fruit.APPLE) alert('apple');
 if (f == Fruit.BANANA) alert('banana');
}
`,
		`/** @constructor */
var Enums = function() {};
var Enums$Fruit$APPLE = 1;
var Enums$Fruit$BANANA = 2;
function foo(f) {
 if (f instanceof Enums) { alert('what?'); return; }
 var Fruit = null;
 if (f == Enums$Fruit$APPLE) alert('apple');
 if (f == Enums$Fruit$BANANA) alert('banana');
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testCollapsePropertiesOfClass1()
func (c *InlineAndCollapsePropertiesTest) TestCollapsePropertiesOfClass1() {
	c.Test(
		`/** @constructor */
var namespace = function() {};
goog.inherits(namespace, Object);

namespace.includeExtraParam = true;

/** @enum {number} */
namespace.Param = {
  param1: 1,
  param2: 2
};

if (namespace.includeExtraParam) {
  namespace.Param.optParam = 3;
}

function f() {
  var Param = namespace.Param;
  log(namespace.Param.optParam);
  log(Param.optParam);
}
`,
		`/** @constructor */
var namespace = function() {};
goog.inherits(namespace, Object);
var namespace$includeExtraParam = true;
var namespace$Param$param1 = 1;
var namespace$Param$param2 = 2;
if (namespace$includeExtraParam) {
  var namespace$Param$optParam = 3;
}
function f() {
  var Param = null;
  log(namespace$Param$optParam);
  log(namespace$Param$optParam);
}
`)
	// TODO(b/148237949): CollapseProperties breaks several optional chain references here
	c.Test(
		`/** @constructor */
var namespace = function() {};
goog.inherits(namespace, Object);

namespace.includeExtraParam = true;

/** @enum {number} */
namespace.Param = {
  param1: 1,
  param2: 2
};

if (namespace?.includeExtraParam) {
  namespace.Param.optParam = 3;
}

function f() {
  var Param = namespace?.Param;
  log(namespace.Param?.optParam);
  log(Param?.optParam);
}
`,
		`/** @constructor */
var namespace = function() {};
goog.inherits(namespace, Object);
var namespace$includeExtraParam = true;
var namespace$Param$param1 = 1;
var namespace$Param$param2 = 2;
/** @enum {number} */
var namespace$Param = {
  param1: namespace$Param$param1,
  param2: namespace$Param$param2
};
if (namespace?.includeExtraParam) { // broken
  var namespace$Param$optParam = 3;
}
function f() {
  var Param = namespace?.Param; // broken
  log(namespace$Param?.optParam); // broken
  log(Param?.optParam); // broken
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testCollapsePropertiesOfClass2()
func (c *InlineAndCollapsePropertiesTest) TestCollapsePropertiesOfClass2() {
	c.Test(
		`var goog = goog || {};
goog.addSingletonGetter = function(cls) {};

var a = {};

/** @constructor */
a.b = function() {};
goog.addSingletonGetter(a.b);
a.b.prototype.get = function(key) {};

/** @constructor */
a.b.c = function() {};
a.b.c.XXX = new a.b.c();

function f() {
  var x = a.b.getInstance();
  var Key = a.b.c;
  x.get(Key.XXX);
}
`,
		`var goog = goog || {};
var goog$addSingletonGetter = function(cls) {};
/** @constructor */
var a$b = function() {};
goog$addSingletonGetter(a$b);
a$b.prototype.get = function(key) {};
/** @constructor */
var a$b$c = function() {};
var a$b$c$XXX = new a$b$c();

function f() {
  var x = a$b.getInstance();
  var Key = null;
  x.get(a$b$c$XXX);
}
`)
}

// Java: InlineAndCollapsePropertiesTest.test_b19179602()
func (c *InlineAndCollapsePropertiesTest) Test_b19179602() {
	// Note that this only collapses a.b.staticProp because a.b is a constructor.
	// Normally AggressiveInlineAliases would not inline "b" inside the loop.
	c.Test(
		`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
function f() {
  while (true) {
// b is declared inside a loop, so it is reassigned multiple times
    var b = a.b;
    alert(b.staticProp);
  }
}
`,
		`/** @constructor */ var a$b = function() {};
var a$b$staticProp = 5;
function f() {
  while (true) {
    var b = null;
    alert(a$b$staticProp);
  }
}
`)
	c.Test(
		`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
function f() {
  while (true) {
// b is declared inside a loop, so it is reassigned multiple times
    var b = a.b;
    alert(b?.staticProp);
  }
}
`,
		`/** @constructor */ var a$b = function() {};
var a$b$staticProp = 5;
function f() {
  while (true) {
    var b = null;
// TODO(bradfordcsmith): This reference is broken by collapsing above.
    alert(a$b?.staticProp);
  }
}
`)
}

// Java: InlineAndCollapsePropertiesTest.test_b19179602_declareOutsideLoop()
func (c *InlineAndCollapsePropertiesTest) Test_b19179602_declareOutsideLoop() {
	c.Test(
		`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
function f() {
// b is declared outside the loop
  var b = a.b;
  while (true) {
    alert(b.staticProp);
  }
}
`,
		`/** @constructor */
var a$b = function() {};
var a$b$staticProp = 5;

function f() {
  var b = null;
  while (true) {
    alert(a$b$staticProp);
  }
}
`)
	c.Test(
		`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
function f() {
// b is declared outside the loop
  var b = a.b;
  while (true) {
    alert(b?.staticProp);
  }
}
`,
		`/** @constructor */
var a$b = function() {};
// TODO(b/148237949): Collapsing this breaks the optional chain reference below.
var a$b$staticProp = 5;

function f() {
  var b = null;
  while (true) {
    alert(a$b?.staticProp);
  }
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testCtorManyAssignmentsDontInlineWarn()
func (c *InlineAndCollapsePropertiesTest) TestCtorManyAssignmentsDontInlineWarn() {
	c.Test(
		`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
function f(y, z) {
  var x = a.b;
  if (y) {
    x = z;
  }
  return x.staticProp;
}
`,
		`/** @constructor */ var a$b = function() {};
var a$b$staticProp = 5;
function f(y, z) {
  var x = a$b;
  if (y) {
    x = z;
  }
  return x.staticProp;
}
`,
		c.Warning(InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING))
}

// Java: InlineAndCollapsePropertiesTest.testCollapseOfGoogModule()
func (c *InlineAndCollapsePropertiesTest) TestCollapseOfGoogModule() {
	c.Test(
		`var goog = goog || {};
/** @constructor */
goog.module = function(id) {};
goog.module.get = function(id) {};
goog.module.getInternal_ = function(id) {};
goog.module = goog.module || {};
`,
		`var goog = goog || {};
var goog$module = function(id) {};
// TODO(b/389129315): This should be goog$module$get.
goog$module.get = function(id) {};
goog$module.getInternal_ = function(id) {};
goog$module = goog$module || {};
`)
}

// Java: InlineAndCollapsePropertiesTest.testCodeGeneratedByGoogModule()
func (c *InlineAndCollapsePropertiesTest) TestCodeGeneratedByGoogModule() {
	// The static property is added to the exports object
	c.Test(
		`var $jscomp = {};
$jscomp.scope = {};
/** @constructor */
$jscomp.scope.Foo = function() {};
var exports = $jscomp.scope.Foo;
exports.staticprop = {A:1};
var y = exports.staticprop.A;
`,
		`/** @constructor */
var $jscomp$scope$Foo = function() {}
var exports = null;
var $jscomp$scope$Foo$staticprop$A = 1;
var y = null;
`)
	// The static property is added to the constructor
	c.Test(
		`var $jscomp = {};
$jscomp.scope = {};
/** @constructor */
$jscomp.scope.Foo = function() {};
$jscomp.scope.Foo.staticprop = {A:1};
var exports = $jscomp.scope.Foo;
var y = exports.staticprop.A;
`,
		`/** @constructor */
var $jscomp$scope$Foo = function() {}
var $jscomp$scope$Foo$staticprop$A = 1;
var exports = null;
var y = null;
`)
}

// Java: InlineAndCollapsePropertiesTest.test_b269515361_codeGeneratedByGoogRequireDynamicForEs5()
func (c *InlineAndCollapsePropertiesTest) Test_b269515361_codeGeneratedByGoogRequireDynamicForEs5() {
	// Case 1
	// The original code of this test case is:
	// ```
	// const {LateLoadTs} = await goog.requireDynamic('a.b.c');
	// new LateLoadTs().render();
	// ```
	// This can be reproduced by actual compilation to es5.
	//
	// This case is for testing the following line, in which the alias is used in dot get property.
	// LateLoadTs$jscomp$1 = $jscomp$destructuring$var22.LateLoadTs
	c.Test(
		`var module$exports$path$lateload_ts =
  {
    LateLoadTs:class LateLoadTs {
      render() {}
  }
};

var func = function() {
  var $jscomp$destructuring$var22;
  var LateLoadTs$jscomp$1;
  return $jscomp$asyncExecutePromiseGeneratorProgram(
    function($jscomp$generator$context){
      if ($jscomp$generator$context.nextAddress == 1) {
        return $jscomp$generator$context.yield(goog.importHandler_('WS8L6d'), 2);
      }
      $jscomp$destructuring$var22 = module$exports$path$lateload_ts;
      LateLoadTs$jscomp$1 = $jscomp$destructuring$var22.LateLoadTs;
      (new LateLoadTs$jscomp$1()).render();
      $jscomp$generator$context.jumpToEnd();
    }
  );
};
`,
		`const CLASS_DECL$0 = class {
  render() {
  }
};
var module$exports$path$lateload_ts$LateLoadTs = CLASS_DECL$0;
var func = function() {
  var $jscomp$destructuring$var22;
  var LateLoadTs$jscomp$1;
  return $jscomp$asyncExecutePromiseGeneratorProgram(function($jscomp$generator$context) {
    if ($jscomp$generator$context.nextAddress == 1) {
      return $jscomp$generator$context.yield(goog.importHandler_("WS8L6d"), 2);
    }
    $jscomp$destructuring$var22 = null;
    LateLoadTs$jscomp$1 = module$exports$path$lateload_ts$LateLoadTs;
    (new LateLoadTs$jscomp$1()).render();
    $jscomp$generator$context.jumpToEnd();
  });
};
`)
	// Case 2
	// This case is for testing assigning value to the exported object.
	//
	// The original code of this test case is:
	// ```
	// const x = await goog.requireDynamic('a.b.c');
	// x.LateLoadTs.num = 0;
	// console.log(x.LateLoadTs);
	// ```
	// This can be reproduced by actual compilation to es5.
	c.Test(
		`var module$exports$path$lateload_ts =
  {
    LateLoadTs:class LateLoadTs { }
};

var func = function() {
  var x;
  return $jscomp$asyncExecutePromiseGeneratorProgram(
    function($jscomp$generator$context){
      if ($jscomp$generator$context.nextAddress == 1) {
        return $jscomp$generator$context.yield(goog.importHandler_('WS8L6d'), 2);
      }
      x = module$exports$path$lateload_ts;
      x.LateLoadTs.num = 0;
      console.log(x.LateLoadTs);
      $jscomp$generator$context.jumpToEnd();
    }
  );
};
`,
		`const CLASS_DECL$0 = class {
};
var module$exports$path$lateload_ts$LateLoadTs = CLASS_DECL$0;
var func = function() {
  var x;
  return $jscomp$asyncExecutePromiseGeneratorProgram(function($jscomp$generator$context) {
    if ($jscomp$generator$context.nextAddress == 1) {
      return $jscomp$generator$context.yield(goog.importHandler_("WS8L6d"), 2);
    }
    x = null;
    module$exports$path$lateload_ts$LateLoadTs.num = 0;
    console.log(module$exports$path$lateload_ts$LateLoadTs);
    $jscomp$generator$context.jumpToEnd();
  });
};
`)
	// Case 3
	// The original code of this test case is:
	// ```
	// const {LateLoadTs} = await goog.requireDynamic('a.b.c');
	// new LateLoadTs().render();
	// ```
	// This can be reproduced by actual compilation to es5.
	//
	// This case is for testing non-exported module, i.e., the global name is not prefixed with
	// 'module$exports$'. Test should fail.
	c.Test(
		`var not_module$exports$path$lateload_ts =
  {
    LateLoadTs:class LateLoadTs {
      render() {}
  }
};

var func = function() {
  var $jscomp$destructuring$var22;
  var LateLoadTs$jscomp$1;
  return $jscomp$asyncExecutePromiseGeneratorProgram(
    function($jscomp$generator$context){
      if ($jscomp$generator$context.nextAddress == 1) {
        return $jscomp$generator$context.yield(goog.importHandler_('WS8L6d'), 2);
      }
      $jscomp$destructuring$var22 = not_module$exports$path$lateload_ts;
      LateLoadTs$jscomp$1 = $jscomp$destructuring$var22.LateLoadTs;
      (new LateLoadTs$jscomp$1()).render();
      $jscomp$generator$context.jumpToEnd();
    }
  );
};
`,
		`const CLASS_DECL$0 = class {
  render() {
  }
};
var not_module$exports$path$lateload_ts = {LateLoadTs:CLASS_DECL$0};
var func = function() {
  var $jscomp$destructuring$var22;
  var LateLoadTs$jscomp$1;
  return $jscomp$asyncExecutePromiseGeneratorProgram(function($jscomp$generator$context) {
    if ($jscomp$generator$context.nextAddress == 1) {
      return $jscomp$generator$context.yield(goog.importHandler_("WS8L6d"), 2);
    }
    $jscomp$destructuring$var22 = not_module$exports$path$lateload_ts;
    LateLoadTs$jscomp$1 = $jscomp$destructuring$var22.LateLoadTs;
    (new LateLoadTs$jscomp$1()).render();
    $jscomp$generator$context.jumpToEnd();
  });
};
`)
	// Case 4
	// This is a manually configured example for this test.
	//
	// This case is for testing the case where the alias is not only used in property access/
	// destructuring.
	c.Test(
		`var _module$exports$path$lateload_ts =
  {
    LateLoadTs:class LateLoadTs {
      render() {}
  }
};

var func = function() {
  var $jscomp$destructuring$var22;
  var LateLoadTs$jscomp$1;
  return $jscomp$asyncExecutePromiseGeneratorProgram(
    function($jscomp$generator$context){
      if ($jscomp$generator$context.nextAddress == 1) {
        return $jscomp$generator$context.yield(goog.importHandler_('WS8L6d'), 2);
      }
      $jscomp$destructuring$var22 = _module$exports$path$lateload_ts;
      $jscomp$destructuring$var22 = 'foo';
      LateLoadTs$jscomp$1 = $jscomp$destructuring$var22.LateLoadTs;
      (new LateLoadTs$jscomp$1()).render();
      $jscomp$generator$context.jumpToEnd();
    }
  );
};
`,
		`const CLASS_DECL$0 = class {
  render() {
  }
};
var _module$exports$path$lateload_ts = {LateLoadTs:CLASS_DECL$0};
var func = function() {
  var $jscomp$destructuring$var22;
  var LateLoadTs$jscomp$1;
  return $jscomp$asyncExecutePromiseGeneratorProgram(function($jscomp$generator$context) {
    if ($jscomp$generator$context.nextAddress == 1) {
      return $jscomp$generator$context.yield(goog.importHandler_("WS8L6d"), 2);
    }
    $jscomp$destructuring$var22 = _module$exports$path$lateload_ts;
    $jscomp$destructuring$var22 = "foo";
    LateLoadTs$jscomp$1 = $jscomp$destructuring$var22.LateLoadTs;
    (new LateLoadTs$jscomp$1()).render();
    $jscomp$generator$context.jumpToEnd();
  });
};
`)
	// Case 5
	// This is a manually configured example for this test.
	//
	// This case is for testing the case where the alias is initialized in declaration and then
	// reassigned a value.
	c.Test(
		`var module$exports$path$lateload_ts =
  {
    LateLoadTs:class LateLoadTs {
      render() {}
  }
};

var func = function() {
  var $jscomp$destructuring$var22 = 'foo';
  var LateLoadTs$jscomp$1;
  return $jscomp$asyncExecutePromiseGeneratorProgram(
    function($jscomp$generator$context){
      if ($jscomp$generator$context.nextAddress == 1) {
        return $jscomp$generator$context.yield(goog.importHandler_('WS8L6d'), 2);
      }
      $jscomp$destructuring$var22 = module$exports$path$lateload_ts;
      LateLoadTs$jscomp$1 = $jscomp$destructuring$var22.LateLoadTs;
      (new LateLoadTs$jscomp$1()).render();
      $jscomp$generator$context.jumpToEnd();
    }
  );
};
`,
		`const CLASS_DECL$0 = class {
  render() {
  }
};
var module$exports$path$lateload_ts = {LateLoadTs:CLASS_DECL$0};
var func = function() {
  var $jscomp$destructuring$var22 = "foo";
  var LateLoadTs$jscomp$1;
  return $jscomp$asyncExecutePromiseGeneratorProgram(function($jscomp$generator$context) {
    if ($jscomp$generator$context.nextAddress == 1) {
      return $jscomp$generator$context.yield(goog.importHandler_("WS8L6d"), 2);
    }
    $jscomp$destructuring$var22 = module$exports$path$lateload_ts;
    LateLoadTs$jscomp$1 = $jscomp$destructuring$var22.LateLoadTs;
    (new LateLoadTs$jscomp$1()).render();
    $jscomp$generator$context.jumpToEnd();
  });
};
`)
}

// Java: InlineAndCollapsePropertiesTest.testInlineCtorInObjLit()
func (c *InlineAndCollapsePropertiesTest) TestInlineCtorInObjLit() {
	c.Test(
		`/** @constructor */
function Foo() {}

/** @constructor */
var Bar = Foo;

var objlit = {
  'prop' : Bar
};
`,
		`/** @constructor */
function Foo() {}
/** @constructor */
var Bar = null;
var objlit$prop = Foo;
`)
}

// Java: InlineAndCollapsePropertiesTest.testNoCollapseExportedNode()
func (c *InlineAndCollapsePropertiesTest) TestNoCollapseExportedNode() {
	c.Test(
		"var x = {}; x.y = {}; var dontExportMe = x.y; use(dontExportMe);", //
		"var x$y = {}; var dontExportMe = null; use(x$y);")
	c.Test(
		"var x = {}; x.y = {}; var _exportMe = x.y;", //
		"var x$y = {}; var _exportMe = x$y;")
}

// Java: InlineAndCollapsePropertiesTest.testDontCrashCtorAliasWithEnum()
func (c *InlineAndCollapsePropertiesTest) TestDontCrashCtorAliasWithEnum() {
	c.Test(
		`var ns = {};
/** @constructor */
ns.Foo = function () {};
var Bar = ns.Foo;
/** @const @enum */
Bar.prop = { A: 1 };
`,
		`/** @constructor */
var ns$Foo = function(){};
var Bar = null;
var ns$Foo$prop$A = 1
`)
}

// Java: InlineAndCollapsePropertiesTest.testDontCrashNamespaceAliasAcrossScopes()
func (c *InlineAndCollapsePropertiesTest) TestDontCrashNamespaceAliasAcrossScopes() {
	c.Test(
		`var ns = {};
ns.VALUE = 0.01;
function f() {
    var constants = ns;
    (function() {
       var x = constants.VALUE;
    })();
}
`,
		(*string)(nil))
}

// Java: InlineAndCollapsePropertiesTest.testGlobalAliasWithLet()
func (c *InlineAndCollapsePropertiesTest) TestGlobalAliasWithLet() {
	c.Test(
		"let ns = {}; ns.foo = 'bar'; let ns2 = ns; use(ns2.foo);", //
		"var ns$foo = 'bar'; let ns2 = null; use(ns$foo);")
	c.Test(
		"let ns = {}; ns.foo = 'bar'; let ns2 =   ns; use?.(ns2?.foo);", //
		"let ns = {}; ns.foo = 'bar'; let ns2 = null; use?.( ns?.foo);")
}

// Java: InlineAndCollapsePropertiesTest.testGlobalAliasForLet2()
func (c *InlineAndCollapsePropertiesTest) TestGlobalAliasForLet2() {
	// We don't do any sort of branch prediction, so can't collapse here.
	c.TestSame(
		`let ns = {};
ns.foo = 'bar';
let ns2;
if (true) {
  ns2 = ns;
}
use(ns2.foo);
`)
}

// Java: InlineAndCollapsePropertiesTest.testGlobalAliasWithLet3()
func (c *InlineAndCollapsePropertiesTest) TestGlobalAliasWithLet3() {
	// Back off since in general we don't check that ns2 is only ever accessed within the if block.
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

// Java: InlineAndCollapsePropertiesTest.testLocalAliasWithLet1()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasWithLet1() {
	c.Test(
		"let ns = {}; ns.foo = 'bar'; function f() { let ns2 =   ns; use(ns2.foo); }",
		"var          ns$foo = 'bar'; function f() { let ns2 = null; use( ns$foo); }")
	c.Test(
		"let ns = {}; ns.foo = 'bar'; function f() { let ns2 =   ns; use?.(ns2?.foo); }",
		"let ns = {}; ns.foo = 'bar'; function f() { let ns2 = null; use?.( ns?.foo); }")
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasWithLet2()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasWithLet2() {
	c.Test(
		`let ns = {};
ns.foo = 'bar';
function f() {
  if (true) {
    let ns2 = ns;
    use(ns2.foo);
  }
}
`,
		`var ns$foo = 'bar';
function f() {
  if (true) {
    let ns2 = null;
    use(ns$foo);
  }
}
`)
	c.Test(
		`let ns = {};
ns.foo = 'bar';
function f() {
  if (true) {
    let ns2 = ns;
    use?.(ns2?.foo);
  }
}
`,
		`let ns = {};
ns.foo = 'bar';
function f() {
  if (true) {
    let ns2 = null;
    use?.(ns?.foo);
  }
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasWithLet3()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasWithLet3() {
	c.Test(
		`let ns = {};
ns.foo = 'bar';
if (true) {
  let ns2 = ns;
  use(ns2.foo);
}
`,
		`var ns$foo = 'bar';
if (true) {
  let ns2 = null;
  use(ns$foo);
}
`)
	c.Test(
		`let ns = {};
ns.foo = 'bar';
if (true) {
  let ns2 = ns;
  use?.(ns2?.foo);
}
`,
		`let ns = {};
ns.foo = 'bar';
if (true) {
  let ns2 = null;
  use?.(ns?.foo);
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasWithLet4()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasWithLet4() {
	c.Test(
		`let ns = {};
ns.foo = 'bar';
if (true) {
  let baz = ns.foo;
  use(baz);
}
`,
		`var ns$foo = 'bar';
if (true) {
  let baz = null;
  use(ns$foo);
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasWithLet5()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasWithLet5() {
	// For local variables (VAR, CONST, or LET) we only handle cases where the alias is a variable,
	// not a property.
	c.TestSame(
		`let ns = {};
ns.foo = 'bar';
if (true) {
  let ns2 = {};
  ns2.baz = ns;
  use(ns2.baz.foo);
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testGlobalAliasWithConst()
func (c *InlineAndCollapsePropertiesTest) TestGlobalAliasWithConst() {
	c.Test(
		"const ns = {}; ns.foo = 'bar'; const ns2 = ns; use(ns2.foo);", //
		"var ns$foo = 'bar'; const ns2 = null; use(ns$foo);")
	c.Test(
		"const ns = {}; ns.foo = 'bar'; const ns2 =   ns; use?.(ns2?.foo);", //
		"const ns = {}; ns.foo = 'bar'; const ns2 = null; use?.( ns?.foo);")
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasWithConst1()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasWithConst1() {
	c.Test(
		`const ns = {};
ns.foo = 'bar';
function f() {
  const ns2 = ns;
  use(ns2.foo);
}
`,
		`var ns$foo = 'bar';
function f() {
  const ns2 = null;
  use(ns$foo);
}
`)
	c.Test(
		`const ns = {};
ns.foo = 'bar';
function f() {
  const ns2 = ns;
  use?.(ns2?.foo);
}
`,
		`const ns = {};
ns.foo = 'bar';
function f() {
  const ns2 = null;
  use?.(ns?.foo);
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasWithConst2()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasWithConst2() {
	c.Test(
		`const ns = {};
ns.foo = 'bar';
function f() {
  if (true) {
    const ns2 = ns;
    use(ns2.foo);
  }
}
`,
		`var ns$foo = 'bar';
function f() {
  if (true) {
    const ns2 = null;
    use(ns$foo);
  }
}
`)
	c.Test(
		`const ns = {};
ns.foo = 'bar';
function f() {
  if (true) {
    const ns2 = ns;
    use?.(ns2?.foo);
  }
}
`,
		`const ns = {};
ns.foo = 'bar';
function f() {
  if (true) {
    const ns2 = null;
    use?.(ns?.foo);
  }
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasWithConst3()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasWithConst3() {
	c.Test(
		`const ns = {};
ns.foo = 'bar';
if (true) {
  const ns2 = ns;
  use(ns2.foo);
}
`,
		`var ns$foo = 'bar';
if (true) {
  const ns2 = null;
  use(ns$foo);
}
`)
	c.Test(
		`const ns = {};
ns.foo = 'bar';
if (true) {
  const ns2 = ns;
  use?.(ns2?.foo);
}
`,
		`const ns = {};
ns.foo = 'bar';
if (true) {
  const ns2 = null;
  use?.(ns?.foo);
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasWithConst4()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasWithConst4() {
	c.Test(
		`const ns = {};
ns.foo = 'bar';
if (true) {
  const baz = ns.foo;
  use(baz);
}
`,
		`var ns$foo = 'bar';
if (true) {
  const baz = null;
  use(ns$foo);
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasWithConst5()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasWithConst5() {
	// For local variables (VAR, CONST, or LET) we only handle cases where the alias is a variable,
	// not a property.
	c.TestSame(
		`const ns = {};
ns.foo = 'bar';
if (true) {
  const ns2 = {};
  ns2.baz = ns;
  use(ns2.baz.foo);
}
`)
}

// Java: InlineAndCollapsePropertiesTest.testLocalAliasInsideClass()
func (c *InlineAndCollapsePropertiesTest) TestLocalAliasInsideClass() {
	c.Test(
		"var a = {b: 5}; class A { fn() { var c = a; use(a.b); } }", //
		"var a$b = 5; class A { fn() { var c = null; use(a$b); } }")
	c.Test(
		"var a = {b: 5}; class A { fn() { var c =    a; use?.(a?.b); } }", //
		"var a = {b: 5}; class A { fn() { var c = null; use?.(a?.b); } }")
}

// Java: InlineAndCollapsePropertiesTest.testEs6ClassStaticProperties()
func (c *InlineAndCollapsePropertiesTest) TestEs6ClassStaticProperties() {
	// Collapsing static properties (A.foo and A.useFoo in this case) is known to be unsafe.
	c.Test(
		c.Srcs(
			`class A {
  static useFoo() {
    alert(this.foo);
  }
}
A.foo = 'bar';
const B = A;
B.foo = 'baz';
B.useFoo();
`),
		c.Expected(
			`var A$useFoo = function() { alert(this.foo); };
class A {}
var A$foo = 'bar';
const B = null;
A$foo = 'baz';
A$useFoo();
`),
		c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
	// Adding @nocollapse makes this safe.
	c.Test(
		`class A {
  /** @nocollapse */
  static useFoo() {
    alert(this.foo);
  }
}
/** @nocollapse */
A.foo = 'bar';
const B = A;
B.foo = 'baz';
B.useFoo();
`,
		`class A {
  /** @nocollapse */
  static useFoo() {
    alert(this.foo);
  }
}
/** @nocollapse */
A.foo = 'bar';
const B = null;
A.foo = 'baz';
A.useFoo();
`)
}

// Java: InlineAndCollapsePropertiesTest.testEs6ClassStaticProperties_asStaticField()
func (c *InlineAndCollapsePropertiesTest) TestEs6ClassStaticProperties_asStaticField() {
	// Collapsing static properties (A.foo and A.useFoo in this case) is known to be unsafe.
	c.Test(
		c.Srcs(
			`class A {
  static useFoo() {
    alert(this.foo);
  }
  static foo = 'bar';
}
const B = A;
B.foo = 'baz';
B.useFoo();
`),
		c.Expected(
			`var A$useFoo = function() {
  alert(this.foo);
};
var A$foo;
class A {}
A$foo = "bar";
const B = null;
A$foo = "baz";
A$useFoo();
`),
		c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
	// Adding @nocollapse makes this safe.
	c.Test(
		`class A {
  /** @nocollapse */
  static useFoo() {
    alert(this.foo);
  }
  /** @nocollapse */
  static foo = 'bar';
}
const B = A;
B.foo = 'baz';
B.useFoo();
`,
		`class A {
  /** @nocollapse */
  static useFoo() {
    alert(this.foo);
  }
  /** @nocollapse */
  static foo;
}
A.foo = "bar";
const B = null;
A.foo = "baz";
A.useFoo();
`)
}

// Java: InlineAndCollapsePropertiesTest.testClassStaticInheritance_method()
func (c *InlineAndCollapsePropertiesTest) TestClassStaticInheritance_method() {
	c.Test(
		"class A { static s() {} } class B extends A {} const C = B;    C.s();", //
		"var A$s = function() {}; class A {} class B extends A {} const C = null; A$s();")
	c.Test(
		"class A { static s() {} } class B extends A {} B.s();", //
		"var A$s = function() {}; class A {} class B extends A {} A$s();")
	c.Test(
		"class A {}     A.s = function() {}; class B extends A {} B.s();", //
		"class A {} var A$s = function() {}; class B extends A {} A$s();")
}

// Java: InlineAndCollapsePropertiesTest.testClassStaticInheritance_propertyAlias()
func (c *InlineAndCollapsePropertiesTest) TestClassStaticInheritance_propertyAlias() {
	c.Test(
		"class A {}     A.staticProp = 6; class B extends A {} let b = new B;", //
		"class A {} var A$staticProp = 6; class B extends A {} let b = new B;")
	c.Test(
		"class A {}     A.staticProp = 6; class B extends A {} use(B.staticProp);", //
		"class A {} var A$staticProp = 6; class B extends A {} use(A$staticProp);")
	c.Test(
		"class A {}     A.staticProp = 6; class B extends A {} use(B?.staticProp);",
		// TODO(b/148237949): collapse breaks `B?.staticProp` reference
		"class A {} var A$staticProp = 6; class B extends A {} use(B?.staticProp);")
}

// Java: InlineAndCollapsePropertiesTest.testClassStaticInheritance_propertyWithSubproperty()
func (c *InlineAndCollapsePropertiesTest) TestClassStaticInheritance_propertyWithSubproperty() {
	c.Test(
		"class A {}     A.ns = {foo: 'bar'}; class B extends A {} use(B.ns.foo);", //
		"class A {} var A$ns$foo = 'bar'; class B extends A {} use(A$ns$foo);")
	c.Test(
		"class A {}     A.ns = {foo: 'bar'}; class B extends A {} use(B.ns?.foo);", //
		"class A {} var A$ns = {foo: 'bar'}; class B extends A {} use(A$ns?.foo);")
	c.Test(
		"class A {} A.ns = {}; class B extends A {}     B.ns.foo = 'baz'; use(B.ns.foo);",
		"class A {}            class B extends A {} var A$ns$foo = 'baz'; use(A$ns$foo);")
	c.Test(
		"class A {}     A.ns = {}; class B extends A {} B.ns.foo = 'baz'; use(B.ns?.foo);",
		"class A {} var A$ns = {}; class B extends A {} A$ns.foo = 'baz'; use(A$ns?.foo);")
}

// Java: InlineAndCollapsePropertiesTest.testClassStaticInheritance_propertyWithShadowing()
func (c *InlineAndCollapsePropertiesTest) TestClassStaticInheritance_propertyWithShadowing() {
	c.Test(
		"class A {}     A.staticProp = 6; class B extends A {}     B.staticProp = 7;", //
		"class A {} var A$staticProp = 6; class B extends A {} var B$staticProp = 7;")
	// At the time use() is called, B.staticProp is still the same as A.staticProp, but we back off
	// rewriting it because of the shadowing afterwards. This makes CollapseProperties unsafely
	// collapse in this case - the same issue occurs when transpiling down to ES5.
	c.Test(
		"class A {}     A.foo = 6; class B extends A {} use(B.foo); B.foo = 7;", //
		"class A {} var A$foo = 6; class B extends A {} use(B$foo); var B$foo = 7;")
	c.Test(
		"class A {}     A.foo = 6; class B extends A {} use(B?.foo); B.foo = 7;",
		// TODO(b/148237949): collapse properties breaks `B?.foo` reference
		"class A {} var A$foo = 6; class B extends A {} use(B?.foo); var B$foo = 7;")
}

// Java: InlineAndCollapsePropertiesTest.testClassStaticInheritance_propertyWithSuperNocollapseSubclassCollapse()
func (c *InlineAndCollapsePropertiesTest) TestClassStaticInheritance_propertyWithSuperNocollapseSubclassCollapse() {
	c.Test(
		c.Srcs(
			`class A {}
/** @nocollapse */
A.foo = 5;
class B extends A {}
B.foo = 6;
use(B.foo);
`),
		c.Expected(
			`class A {}
/** @nocollapse */
A.foo = 5;
class B extends A {}
var B$foo = 6;
use(B$foo);
`))
}

// Java: InlineAndCollapsePropertiesTest.testClassStaticInheritance_cantDetermineSuperclass()
func (c *InlineAndCollapsePropertiesTest) TestClassStaticInheritance_cantDetermineSuperclass() {
	// Here A.foo and B.foo are unsafely collapsed because getSuperclass() creates an alias for them
	c.Test(
		`class A {}
A.foo = 5;
class B {}
B.foo = 6;
function getSuperclass() { return 1 < 2 ? A : B; }
class C extends getSuperclass() {}
use(C.foo);
`,
		`class A {
}
var A$foo = 5;
class B {
}
var B$foo = 6;
function getSuperclass() {
  return 1 < 2 ? A : B;
}
const CLASS_EXTENDS$0 = getSuperclass();
class C extends CLASS_EXTENDS$0 {
}
use(C.foo);
`)
}

// Java: InlineAndCollapsePropertiesTest.testTypeScriptDecoratedClass()
func (c *InlineAndCollapsePropertiesTest) TestTypeScriptDecoratedClass() {
	// TypeScript 5.8 emits this for decorated classes.
	c.Test(
		`class A {
  static getId() {
    return A.ID;
  }
}
A.ID = 'a';
tslib.__decorate([], A);
if (false) {
  /** @const {string} */ A.ID;
}
console.log(A.getId());
`,
		`var A$getId = function() {
  return A$ID;
};
class A {
}
var A$ID = "a";
tslib.__decorate([], A);
if (false) {
  A$ID;
}
console.log(A$getId());
`)
}

// Java: InlineAndCollapsePropertiesTest.testAliasForSuperclassNamespace()
func (c *InlineAndCollapsePropertiesTest) TestAliasForSuperclassNamespace() {
	c.Test(
		`var ns = {};
class Foo {}
ns.clazz = Foo;
var Bar = class extends ns.clazz.Baz {}
`,
		`class Foo {}
var ns$clazz = null;
var Bar = class extends Foo.Baz {}
`)
	c.Test(
		`var ns = {};
class Foo {}
Foo.Builder = class {}
ns.clazz = Foo;
var Bar = class extends ns.clazz.Builder {}
`,
		`class Foo {}
var Foo$Builder = class {}
var ns$clazz = null;
var Bar = class extends Foo$Builder {}
`)
}

// Java: InlineAndCollapsePropertiesTest.testAliasForSuperclassNamespace_withStaticInheritance()
func (c *InlineAndCollapsePropertiesTest) TestAliasForSuperclassNamespace_withStaticInheritance() {
	c.Test(
		`var ns = {};
class Foo {}
Foo.Builder = class {}
Foo.Builder.baz = 3;
ns.clazz = Foo;
var Bar = class extends ns.clazz.Builder {}
use(Bar.baz);
`,
		`class Foo {}
var Foo$Builder = class {}
var Foo$Builder$baz = 3;
var ns$clazz = null;
var Bar = class extends Foo$Builder {}
use(Foo$Builder$baz);
`)
	c.Test(
		`var ns = {};
class Foo {}
Foo.Builder = class {}
Foo.Builder.baz = 3;
ns.clazz = Foo;
var Bar = class extends ns.clazz?.Builder {}
use(Bar?.baz);
`,
		`class Foo {
}
var Foo$Builder = class {
};
var Foo$Builder$baz = 3;
var ns$clazz = null;
// TODO(b/148237979): Collapse above breaks this reference.
const CLASS_EXTENDS$0 = Foo?.Builder;
var Bar = class extends CLASS_EXTENDS$0 {
};
// TODO(b/148237979): Collapse above breaks this reference.
use(Bar?.baz);
`)
}

// Java: InlineAndCollapsePropertiesTest.testDestructuringAlias1()
func (c *InlineAndCollapsePropertiesTest) TestDestructuringAlias1() {
	c.TestSame("var a = { x: 5 }; var [b] = [a]; use(b.x);")
}

// Java: InlineAndCollapsePropertiesTest.testDestructuringAlias2()
func (c *InlineAndCollapsePropertiesTest) TestDestructuringAlias2() {
	c.TestSame("var a = { x: 5 }; var {b} = {b: a}; use(b.x);")
}

// Java: InlineAndCollapsePropertiesTest.testDestructingAliasWithConstructor()
func (c *InlineAndCollapsePropertiesTest) TestDestructingAliasWithConstructor() {
	c.Test(
		`var ns = {};
/** @constructor */
ns.ctor = function() {}
const {ctor} = ns;
let c = new ctor;
`,
		`/** @constructor */
var ns$ctor = function() {}
const ctor = null;
let c = new ns$ctor;
`)
}

// Java: InlineAndCollapsePropertiesTest.testDestructuringThatUsedToCrash()
func (c *InlineAndCollapsePropertiesTest) TestDestructuringThatUsedToCrash() {
	c.Test(
		`var CakeFlavors = {
  CARROT: 1,
  TIRAMISU: 2,
};
class UglyCake {}
const Cake = UglyCake;
/** @const */
UglyCake.Flavors = CakeFlavors;
const {Flavors: {CARROT, TIRAMISU}} = Cake;
alert(CARROT, TIRAMISU);
`,
		`var CakeFlavors$CARROT = 1;
var CakeFlavors$TIRAMISU = 2;
class UglyCake {}
const Cake = null;
/** @const */
var UglyCake$Flavors = null;
const destructuring$m1146332801$0 = null;
const CARROT = null;
const TIRAMISU = null;
alert(CakeFlavors$CARROT, CakeFlavors$TIRAMISU);
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineAndCollapsePropertiesTest.namespaceInDestructuringPattern()
func (c *InlineAndCollapsePropertiesTest) TestNamespaceInDestructuringPattern() {
	c.TestSame(
		`const ns = {};
ns.x = 1;
ns.y = 2;
let {x, y} = ns;
x = 4; // enforce that we can't inline x -> ns.x because it's set multiple times
use(x + y);
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineAndCollapsePropertiesTest.inlineDestructuringPatternConstructorWithProperty()
func (c *InlineAndCollapsePropertiesTest) TestInlineDestructuringPatternConstructorWithProperty() {
	c.Test(
		`const ns = {};
/** @constructor */
ns.Y = function() {};
ns.Y.prop = 3;
const {Y} = ns;
use(Y.prop);
`,
		`/** @constructor */
var ns$Y = function() {};
var ns$Y$prop = 3;
const Y = null;
use(ns$Y$prop);
`)
	c.Test(
		`const ns = {};
/** @constructor */
ns.Y = function() {};
ns.Y.prop = 3;
const {Y} = ns;
use(Y?.prop);
`,
		`/** @constructor */
var ns$Y = function() {};
// TODO(bradfordcsmith): This collapse breaks the optional chain
//     reference below.
var ns$Y$prop = 3;
const Y = null;
use(ns$Y?.prop);
`)
	// Only `const` destructuring aliases have special handling by AggressiveInlineAliases
	c.TestSame(
		`const ns = {};
/** @constructor */
ns.Y = function() {};
ns.Y.prop = 3;
let {Y} = ns;
use(Y.prop);
`,
		InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING)
}

// TODO(b/538122483): InlineAndCollapseProperties: handle destructuring patterns without creating
// extra variables or losing side effects
// Java: InlineAndCollapsePropertiesTest.testDestructuringQualifiedNameMultipleKeysEvaluatesOnce()
func (c *InlineAndCollapsePropertiesTest) TestDestructuringQualifiedNameMultipleKeysEvaluatesOnce() {
	c.Test(
		`const {a, b} = this.validated;
use(a, b);
`,
		`const a = this.validated.a;
const b = this.validated.b;
use(a, b);
`)
}

// Java: InlineAndCollapsePropertiesTest.testDestructuringQualifiedNameSingleKeyEvaluatesDirectly()
func (c *InlineAndCollapsePropertiesTest) TestDestructuringQualifiedNameSingleKeyEvaluatesDirectly() {
	c.Test(
		`const {a} = this.validated;
use(a);
`,
		`const a = this.validated.a;
use(a);
`)
}

// TODO(b/538122483): InlineAndCollapseProperties: handle destructuring patterns without creating
// extra variables or losing side effects
// Java: InlineAndCollapsePropertiesTest.testEmptyDestructuringPatternNotStripped()
func (c *InlineAndCollapsePropertiesTest) TestEmptyDestructuringPatternNotStripped() {
	c.Test("const {} = this.validated;", "")
	c.Test("const {} = obj;", "")
}

// TODO(b/538122483): InlineAndCollapseProperties: handle destructuring patterns without creating
// extra variables or losing side effects
// Java: InlineAndCollapsePropertiesTest.testDestructuringQualifiedNameWithCollapsibleNamespace()
func (c *InlineAndCollapsePropertiesTest) TestDestructuringQualifiedNameWithCollapsibleNamespace() {
	c.Test(
		`const ns = {
  sub: {
    a: 1,
    b: 2,
  },
};
const {a, b} = ns.sub;
use(a, b);
`,
		`var ns$sub$a = 1;
var ns$sub$b = 2;
const a = null;
const b = null;
use(ns$sub$a, ns$sub$b);
`)
}

// Java: InlineAndCollapsePropertiesTest.testDefaultParamAlias1()
func (c *InlineAndCollapsePropertiesTest) TestDefaultParamAlias1() {
	c.Test(
		"var a = {b: 5}; var b = a; function f(x=b) { alert(x.b); }", //
		"var a = {b: 5}; var b = null; function f(x=a) { alert(x.b); }")
}

// Java: InlineAndCollapsePropertiesTest.testDefaultParamAlias2()
func (c *InlineAndCollapsePropertiesTest) TestDefaultParamAlias2() {
	c.Test(
		"var a = {b: {c: 5}}; var b = a; function f(x=b.b) { alert(x.c); }", //
		"var a$b = {c: 5}; var b = null; function f(x=a$b) { alert(x.c); }")
	c.Test(
		"var a = {b: {c: 5}}; var b =    a; function f(x=b?.b) { alert(x?.c); }", //
		"var a = {b: {c: 5}}; var b = null; function f(x=a?.b) { alert(x?.c); }")
}

// Java: InlineAndCollapsePropertiesTest.testAliasPropertyOnUnsafelyRedefinedNamespace()
func (c *InlineAndCollapsePropertiesTest) TestAliasPropertyOnUnsafelyRedefinedNamespace() {
	c.TestSame("var obj = {foo: 3}; var foo = obj.foo; obj = {}; alert(foo);")
}

// Java: InlineAndCollapsePropertiesTest.testAliasPropertyOnSafelyRedefinedNamespace()
func (c *InlineAndCollapsePropertiesTest) TestAliasPropertyOnSafelyRedefinedNamespace() {
	// non-constructor property doesn't get collapsed
	c.Test(
		"var obj = {foo: 3}; var foo = obj.foo; obj = obj || {}; alert(foo);", //
		"var obj = {foo: 3}; var foo = null   ; obj = obj || {}; alert(obj.foo);")
	// constructor property does get collapsed
	c.Test(
		`var ns = {};
/** @constructor */ ns.ctor = function() {};
ns.ctor.foo = 3;
var foo = ns.ctor.foo;
ns = ns || {};
alert(foo);
`,
		`var ns = {};
/** @constructor */ var ns$ctor = function() {};
var ns$ctor$foo = 3;
var foo = null;
ns = ns || {};
alert(ns$ctor$foo);
`)
	// NOTE(lharker): this mirrors current code in Closure library
	c.Test(
		`var goog = {};
goog.module = function() {};
/** @constructor */ goog.module.ModuleManager = function() {};
goog.module.ModuleManager.getInstance = function() {};
goog.module = goog.module || {};
var ModuleManager = goog.module.ModuleManager;
alert(ModuleManager.getInstance());
`,
		`var goog$module = function() {};
/** @constructor */ var goog$module$ModuleManager = function() {};
var goog$module$ModuleManager$getInstance = function() {};
goog$module = goog$module || {};
var ModuleManager = null;
alert(goog$module$ModuleManager$getInstance());
`)
}

// Java: InlineAndCollapsePropertiesTest.testClassInObjectLiteral()
func (c *InlineAndCollapsePropertiesTest) TestClassInObjectLiteral() {
	c.Test(
		`var obj = {
  foo: class {}
};
`,
		`const CLASS_DECL$0 = class {
};
var obj$foo = CLASS_DECL$0;
`)
	// TODO(lharker): this is unsafe, obj$foo.foo is undefined now that A$foo is collapsed
	c.Test(
		`class A {}
A.foo = 3;
var obj = {foo: class extends A {}};
use(obj.foo.foo);
`,
		`class A {
}
var A$foo = 3;
const CLASS_DECL$0 = class extends A {
};
var obj$foo = CLASS_DECL$0;
use(obj$foo.foo);
`)
}

// Java: InlineAndCollapsePropertiesTest.testToStringValueOfPropertiesInObjectLiteral()
func (c *InlineAndCollapsePropertiesTest) TestToStringValueOfPropertiesInObjectLiteral() {
	// toString/valueOf are not collapsed because they are properties implicitly used as part of the
	// JS language
	c.TestSame(
		`let z = {
  toString() { return 'toString'; },
  valueOf() { return 'valueOf'; }
};
var zAsString = z + "";
`)
}

// Java: InlineAndCollapsePropertiesTest.testToStringValueOfNonMethodPropertiesInObjectLiteral()
func (c *InlineAndCollapsePropertiesTest) TestToStringValueOfNonMethodPropertiesInObjectLiteral() {
	c.TestSame(
		`let z = {
  toString: function() { return 'a'; },
  valueOf: function() { return 3; },
};
`)
}

// Java: InlineAndCollapsePropertiesTest.testToStringValueOfPropertiesInObjectLiteralAssignmentDepth()
func (c *InlineAndCollapsePropertiesTest) TestToStringValueOfPropertiesInObjectLiteralAssignmentDepth() {
	c.Test(
		`var a = {}; a.b = {}; a.b.c = {};
a.b.c.d = {
  toString() { return 'a'; },
  valueOf: function() { return 3; },
};
`,
		`var a$b$c$d = {
  toString() { return 'a'; },
  valueOf: function() { return 3; },
};
`)
}

// Java: InlineAndCollapsePropertiesTest.testLoopInAliasChainOfTypedefConstructorProperty()
func (c *InlineAndCollapsePropertiesTest) TestLoopInAliasChainOfTypedefConstructorProperty() {
	c.Test(
		`/** @constructor */ var Item = function() {};
/** @typedef {number} */ Item.Models;
Item.Models = Item.Models;
`,
		`/** @constructor */ var Item = function() {};
Item$Models;
var Item$Models = Item$Models;
`)
	c.Test(
		`/** @constructor */ var Item = function() {};
/** @typedef {number} */ Item.Models;
Item.Models = Item?.Models;
`,
		`/** @constructor */ var Item = function() {};
Item$Models;
// TODO(b/148237949): Collapsing breaks this reference.
var Item$Models = Item?.Models;
`)
}

// Java: InlineAndCollapsePropertiesTest.testDontInlinePropertiesOnEscapedNamespace()
func (c *InlineAndCollapsePropertiesTest) TestDontInlinePropertiesOnEscapedNamespace() {
	c.TestSame(
		c.Externs("function use(obj) {}"),
		c.Srcs(
			`function Foo() {}
Foo.Bar = {};
Foo.Bar.baz = {A: 1, B: 2};

var $jscomp$destructuring$var1 = Foo.Bar;
// This call could potentially have changed the value of Foo.Bar, so don't
// collapse Foo.Bar.baz.A/Foo.Bar.baz.B
use(Foo);
var baz = $jscomp$destructuring$var1.baz;
use(baz.A);
`))
}

// Java: InlineAndCollapsePropertiesTest.testInliningPropertiesOnEscapedNamespace_withDeclaredType()
func (c *InlineAndCollapsePropertiesTest) TestInliningPropertiesOnEscapedNamespace_withDeclaredType() {
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
use(Foo);
var baz = $jscomp$destructuring$var1.baz;
use(baz.A);
`),
		c.Expected(
			`/** @constructor */
function Foo() {}
/** @constructor */
var Foo$Bar = function() {};
var Foo$Bar$baz$A = 1;
var Foo$Bar$baz$B = 2;

var $jscomp$destructuring$var1 = null;
// This call could potentially read Foo.Bar, which will now be broken.
// AggressiveInlineAliases/CollapseProperties intentionally generate unsafe code
// in this case. The main motivation is to collapse long namespaces generated by
// goog.provide or goog.module.declareLegacyNamespace()
// https://developers.google.com/closure/compiler/docs/limitations#implications-of-object-property-flattening
use(Foo);
var baz = null;
use(Foo$Bar$baz$A);
`))
}

// Java: InlineAndCollapsePropertiesTest.testDontInlinePropertiesOnNamespace_withNoCollapse()
func (c *InlineAndCollapsePropertiesTest) TestDontInlinePropertiesOnNamespace_withNoCollapse() {
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
// This call could potentially have changed the value of Foo.Bar, and without
// @nocollapse on Foo.Bar and Foo.Bar.baz the compiler would generate unsafe code.
use(Foo);
var baz = $jscomp$destructuring$var1.baz;
use(baz.A);
`))
}

// Java: InlineAndCollapsePropertiesTest.testInlinePropertyOnAliasedConstructor()
func (c *InlineAndCollapsePropertiesTest) TestInlinePropertyOnAliasedConstructor() {
	// TODO(b/117905881): as long as we unsafely collapse Foo.prop -> Foo$prop, there's no way to
	// safely rewrite the "use(ns.alias.prop)" call in the compiler
	// Either
	// a) we replace "ns.alias.prop" with "Foo$prop" or
	// b) we leave "ns.alias.prop" as is.
	// Option a) can break code if the "use(ns)" call actually reassigns ns.alias (which is not
	// declared const).
	// Option b) breaks code if "use(ns)" does NOT reassign ns.alias, since now
	// ns.alias.prop is Foo.prop, which is undefined.
	// but we might want to consider movign to (b) instead of (a), since CollapseProperties will
	// also break case (a)
	c.Test(
		c.Externs("function use(ctor) {}"),
		c.Srcs(
			`/** @constructor */
function Foo() {}
Foo.prop = 2;
var ns = {};
ns.alias = Foo;
use(ns);
use(ns.alias.prop);
`),
		c.Expected(
			`/** @constructor */
function Foo() {}
var Foo$prop = 2;
var ns = {};
ns.alias = Foo;
use(ns);
use(ns.alias.prop);
`))
}

// Java: InlineAndCollapsePropertiesTest.testClassStaticMemberAccessedWithSuper()
func (c *InlineAndCollapsePropertiesTest) TestClassStaticMemberAccessedWithSuper() {
	c.Test(
		`class Bar {
  static double(n) {
    return n*2
  }
}
class Baz extends Bar {
  static quadruple(n) {
    return 2 * super.double(n);
  }

  static val1;
  static {
    Baz.val1 = super.double(5);
  }

  static val2 = super.double(6);
}
`,
		`var Bar$double = function(n) {
  return n * 2;
};
class Bar {}
var Baz$quadruple = function(n$jscomp$1) {
  return 2 * Bar$double(n$jscomp$1);
};
var Baz$val1;
var Baz$val2;
class Baz extends Bar {}
{
  Baz$val1 = Bar$double(5);
}
Baz$val2 = Bar$double(6);
`)
}

// Java: InlineAndCollapsePropertiesTest.testClassStaticMemberAccessedWithSuperAndThis()
func (c *InlineAndCollapsePropertiesTest) TestClassStaticMemberAccessedWithSuperAndThis() {
	src := "class Bar {\n" +
		"  static get name() {\n" +
		"    return 'Bar';\n" +
		"  }\n" +
		"  /** @nocollapse */\n" +
		"  static getClassname() {\n" +
		"    return `${this.name} class`;\n" +
		"  }\n" +
		"}\n" +
		"class Baz extends Bar {\n" +
		"  static get name() {\n" +
		"    return 'Baz';\n" +
		"  }\n" +
		"  static get classname() {\n" +
		"    return `${super.getClassname()} - is a subclass`;\n" +
		"  }\n" +
		"}\n"
	c.SetAssumeStaticInheritanceIsNotUsed(false)
	c.TestSame(src)
	c.SetAssumeStaticInheritanceIsNotUsed(true)
	c.Test(
		src,
		"class Bar {\n"+
			"  static get name() {\n"+
			"    return \"Bar\";\n"+
			"  }\n"+
			"  /** @nocollapse */\n"+
			"  static getClassname() {\n"+
			"    return `${this.name} class`;\n"+
			"  }\n"+
			"}\n"+
			"class Baz extends Bar {\n"+
			"  static get name() {\n"+
			"    return \"Baz\";\n"+
			"  }\n"+
			"  static get classname() {\n"+
			"    return `${Bar.getClassname()} - is a subclass`;\n"+
			"  }\n"+
			"}\n")
}

// Java: InlineAndCollapsePropertiesTest.testClassStaticMemberAccessedWithThis()
func (c *InlineAndCollapsePropertiesTest) TestClassStaticMemberAccessedWithThis() {
	c.Test(
		`class Baz extends Bar {
  static double(n) {
    return 2 * n;
  }

  static val1;
  static {
    this.val1 = this.double(1);
  }

  static val2;
  static {
    Baz.val2 = Baz.double(2);
  }

  static val3 = this.double(3);

  static val4 = Baz.double(4);
}
`,
		`var Baz$double = function(n) {
  return 2 * n;
};
var Baz$val1;
var Baz$val2;
var Baz$val3;
var Baz$val4;
class Baz extends Bar {}
{
  Baz$val1 = Baz$double(1);
}
{
  Baz$val2 = Baz$double(2);
}
Baz$val3 = Baz$double(3);
Baz$val4 = Baz$double(4);
`)
}

// Java: InlineAndCollapsePropertiesTest.testClassStaticGetterAccessedWithSuper()
func (c *InlineAndCollapsePropertiesTest) TestClassStaticGetterAccessedWithSuper() {
	src := `function JSCompiler_renameProperty(str, strContainerType) {
  return a;
}

class Parent {
  static getName() {
    return 'Parent';
  }
  static get greeting() {
    return 'Hello ' + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return 'Child';
  }
  static msg = super.greeting;  // 'Hello Child'
}
`
	// non-strict
	c.Test(
		src,
		`function JSCompiler_renameProperty(str, strContainerType) {
  return a;
}
var Parent$getName = function() {
  return "Parent";
};
class Parent {
  static get greeting() {
    return "Hello " + this.getName();
  }
}
var Child$getName = function() {
  return "Child";
};
var Child$msg;
class Child extends Parent {}
Child$msg = Parent.greeting;
`)
	// strict
	c.SetAssumeStaticInheritanceIsNotUsed(false)
	c.Test(
		src,
		`function JSCompiler_renameProperty(str, strContainerType) {
  return a;
}
class Parent {
  static getName() {
    return "Parent";
  }
  static get greeting() {
    return "Hello " + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return "Child";
  }
  static msg;
  static STATIC_INIT$0() {
    Child.msg = super.greeting;
  }
}
Child.STATIC_INIT$0();
`)
}

// Java: InlineAndCollapsePropertiesTest.testClassExpressionWithInnerClassName()
func (c *InlineAndCollapsePropertiesTest) TestClassExpressionWithInnerClassName() {
	c.Test(
		`const OuterName = class InnerName {
  static sf1 = 1;
  static sf2 = InnerName.sf1;
  static sf3 = this.sf2;

  static {
    Object.defineProperties(InnerName.prototype, {bar: {value: 1}});
    InnerName.sf2++;
    this.sf3++;
  }
};
alert(OuterName.sf3);
`,
		`var OuterName$sf1;
var OuterName$sf2;
var OuterName$sf3;
const OuterName = class {};
OuterName$sf1 = 1;
OuterName$sf2 = OuterName$sf1;
OuterName$sf3 = OuterName$sf2;
{
  Object.defineProperties(OuterName.prototype, {bar:{value:1}});
  OuterName$sf2++;
  OuterName$sf3++;
}
alert(OuterName$sf3);
`)
}

// Java: InlineAndCollapsePropertiesTest.testCommaCallAliasing1()
func (c *InlineAndCollapsePropertiesTest) TestCommaCallAliasing1() {
	c.Test(
		`var xid = function(a) {};
xid.internal_ = function() {};
(0,xid)('abc');
`,
		`var xid = function(a) {};
// The indirect call below does not prevent this property collapse
var xid$internal_ = function() {};
(0,xid)('abc');
`)
}

// Java: InlineAndCollapsePropertiesTest.testCommaCallAliasing2()
func (c *InlineAndCollapsePropertiesTest) TestCommaCallAliasing2() {
	c.Test(
		`var xid = function(a) {};
xid.internal_ = function() {};
fn((0,xid)('abc'));
`,
		`var xid = function(a) {};
// The indirect call below does not prevent this property collapse
var xid$internal_ = function() {};
fn((0,xid)('abc'));
`)
}

// Java: InlineAndCollapsePropertiesTest.testCommaCallAliasing3()
func (c *InlineAndCollapsePropertiesTest) TestCommaCallAliasing3() {
	c.Test(
		`var xid = function(a) {};
xid.internal_ = function() {};
throw ((0, xid), 0);
`,
		`var xid = function(a) {};
// The indirect call below does not prevent this property collapse
var xid$internal_ = function() {};
throw ((0, xid), 0);
`)
}

// Java: InlineAndCollapsePropertiesTest.testHookNewAliasing1()
func (c *InlineAndCollapsePropertiesTest) TestHookNewAliasing1() {
	c.TestSame(
		`var xid = function(a) {};
xid.internal_ = function() {};
// The indirect call prevents this property collapse
var xx = new Thing(x ? xid : abc).method();
`)
}

// Java: InlineAndCollapsePropertiesTest.testHookGetProp1()
func (c *InlineAndCollapsePropertiesTest) TestHookGetProp1() {
	c.TestSame(
		`var xid = function(a) {};
xid.internal_ = function() {};
// The indirect call prevents this property collapse
var xx = (x ? xid : abc).method();
`)
}

// Java: InlineAndCollapsePropertiesTest.testHookGetProp2()
func (c *InlineAndCollapsePropertiesTest) TestHookGetProp2() {
	c.Test(
		`var xid = function(a) {};
xid.internal_ = function() {};
var xx = ((x ? xid : abc) - def).toExponential(5);
`,
		`var xid = function(a) {};
// The indirect call below does not prevent this property collapse
var xid$internal_ = function() {};
var xx = ((x ? xid : abc) - def).toExponential(5);
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineAndCollapsePropertiesTest.conditionallyLoadedChunk_withExistenceCheckInEarlierChunk()
func (c *InlineAndCollapsePropertiesTest) TestConditionallyLoadedChunk_withExistenceCheckInEarlierChunk() {
	// Regression test for b/461276859.
	chunks := JSChunkGraphBuilder_ForChain().
		// m1, always loaded
		AddChunk(
			`const mods = {};

function f() {
  if (mods.fn) {
    mods.fn();
  }
}
`).
		// m2, conditionally loaded after m1
		AddChunk("mods.fn = () => {};").
		Build()
	c.Test(
		c.Srcs(chunks[0], chunks[1]),
		c.Expected(
			`var mods$fn;
function f() {
  if (mods$fn) {
    mods$fn();
  }
}
`,
			"mods$fn = () => {};"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineAndCollapsePropertiesTest.conditionallyLoadedChunk_withExistenceCheckInEarlierChunk_constDecl()
func (c *InlineAndCollapsePropertiesTest) TestConditionallyLoadedChunk_withExistenceCheckInEarlierChunk_constDecl() {
	// Regression test for b/461276859.
	chunks := JSChunkGraphBuilder_ForChain().
		// m1, always loaded
		AddChunk(
			`function f() {
  if (typeof mods !== 'undefined') {
    mods.fn();
  }
}
`).
		// m2, conditionally loaded after m1
		AddChunk("const mods = {}; mods.fn = () => {};").
		Build()
	c.Test(
		c.Srcs(chunks[0], chunks[1]),
		c.Expected(
			// We don't need to add a stub declaration of 'mods' as this pass isn't collapsing it.
			// The input code is responsible for avoiding unsafe 'mods' cross-chunk accesses.
			`var mods$fn;
function f() {
  if (typeof mods !== 'undefined') {
    mods$fn();
  }
}
`,
			`const mods = {};
mods$fn = () => {};
`))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineAndCollapsePropertiesTest.conditionallyLoadedChunk_withExistenceCheckInEarlierChunk_usesLeastCommonAncestor()
func (c *InlineAndCollapsePropertiesTest) TestConditionallyLoadedChunk_withExistenceCheckInEarlierChunk_usesLeastCommonAncestor() {
	// Regression test for b/461276859.
	chunks := JSChunkGraphBuilder_ForTree().
		// base chunk, always loaded
		AddChunk(
			`const mods = {};
`).
		// c1, conditionally loaded. does not depend on c2.
		AddChunk(
			`function f() {
  if (mods.fn) {
    mods.fn();
  }
}
`).
		// c2, conditionally loaded. depends on c1.
		AddChunk("mods.fn = () => {};").
		Build()
	c.Test(
		c.Srcs(chunks[0], chunks[1], chunks[2]),
		c.Expected(
			"var mods$fn;",
			`function f() {
  if (mods$fn) {
    mods$fn();
  }
}
`,
			"mods$fn = () => {};"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineAndCollapsePropertiesTest.conditionallyLoadedChunk_withExistenceCheckInEarlierChunk_andChildName()
func (c *InlineAndCollapsePropertiesTest) TestConditionallyLoadedChunk_withExistenceCheckInEarlierChunk_andChildName() {
	// Regression test for b/461276859.
	chunks := JSChunkGraphBuilder_ForChain().
		// m1, always loaded
		AddChunk(
			`const mods = {};
function f() {
  if (mods.a && mods.a.crossChunk) {
    mods.a.crossChunk();
  }
}
`).
		// m2, conditionally loaded after m1
		AddChunk(
			`mods.a = {};
mods.a.crossChunk = () => {};
mods.a.local = () => {};
`).
		Build()
	c.Test(
		c.Srcs(chunks[0], chunks[1]),
		c.Expected(
			`var mods$a$crossChunk;
var mods$a;
function f() {
  if (mods$a && mods$a$crossChunk) {
    mods$a$crossChunk();
  }
}
`,
			`mods$a = {};
mods$a$crossChunk = () => {};
// Verify that we wait to declare mods$a$local until this chunk, as it's not referenced
// elsewhere.
var mods$a$local = () => {};
`))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineAndCollapsePropertiesTest.conditionallyLoadedChunk_unsafeToInlineVar()
func (c *InlineAndCollapsePropertiesTest) TestConditionallyLoadedChunk_unsafeToInlineVar() {
	chunks := JSChunkGraphBuilder_ForChain().
		// m1, always loaded
		AddChunk(
			`var modFactory;
const registry = {};
const factories = {};
factories.one = function() {};

function build() {
  const mod = modFactory;
  if (mod) mod();
}
`).
		// m2, conditionally loaded after m1
		AddChunk("modFactory = factories.one; modFactory();").
		Build()
	c.Test(
		c.Srcs(chunks[0], chunks[1]),
		c.Expected(
			`var modFactory;
var factories$one = function() {};

function build() {
  const mod = modFactory;
  if (mod) mod();
}
`,
			"modFactory = factories$one; factories$one();"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineAndCollapsePropertiesTest.conditionallyLoadedChunk_unsafeToInlineQname()
func (c *InlineAndCollapsePropertiesTest) TestConditionallyLoadedChunk_unsafeToInlineQname() {
	chunks := JSChunkGraphBuilder_ForChain().
		// m1, always loaded
		AddChunk(
			`const registry = {};
const factories = {};
factories.one = function() {};

function build() {
  const mod = registry.mod;
  if (mod) mod();
}
`).
		// m2, conditionally loaded after m1
		AddChunk("registry.mod = factories.one; registry.mod();").
		Build()
	c.Test(
		c.Srcs(chunks[0], chunks[1]),
		c.Expected(
			`var registry$mod;
var factories$one = function() {};

function build() {
  const mod = null;
  if (registry$mod) registry$mod();
}
`,
			"registry$mod = factories$one; factories$one();"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineAndCollapsePropertiesTest.conditionallyLoadedChunk_safeToInline()
func (c *InlineAndCollapsePropertiesTest) TestConditionallyLoadedChunk_safeToInline() {
	chunks := JSChunkGraphBuilder_ForChain().
		// m1, always loaded
		AddChunk(
			`const registry = {};
const factories = {};
factories.one = function() {};
registry.mod = factories.one;
`).
		// m2, conditionally loaded after m1
		// Even though this chunk is conditionally loaded, it depends on m1, so it's safe to
		// inline the value of registry.mod because if this chunk executes, m1 must have
		// executed first.
		AddChunk(
			`function build() {
  const mod = registry.mod;
  if (mod) mod();
}
`).
		Build()
	c.Test(
		c.Srcs(chunks[0], chunks[1]),
		c.Expected(
			// m1
			`var factories$one = function() {};
var registry$mod = null;
`,
			// m2
			`function build() {
  const mod = null;
  if (factories$one) factories$one();
}
`))
}

// Java: InlineAndCollapsePropertiesTest.testDestructuringGlobalNameExtractor()
func (c *InlineAndCollapsePropertiesTest) TestDestructuringGlobalNameExtractor() {
	c.Test("var ns = {}; ns.foo = 1; const {foo} = ns;", "var ns$foo = 1; const foo = null;")
}

// Java: InlineAndCollapsePropertiesTest.testConcretizeStaticInheritanceForInlining()
func (c *InlineAndCollapsePropertiesTest) TestConcretizeStaticInheritanceForInlining() {
	c.Test(
		"class A {} class B extends A {} B.foo = 1;",
		"class A {} class B extends A {} var B$foo = 1;")
}

// Java: InlineAndCollapsePropertiesTest.testStaticFieldOnClassInAliasedNamespace()
func (c *InlineAndCollapsePropertiesTest) TestStaticFieldOnClassInAliasedNamespace() {
	// TODO(b/571958823): Inlining alias crashes on static MEMBER_FIELD_DEF.
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2022)
	e :=
		AssertThrows[jutil.IRuntimeException](
			func() {
				c.TestSame(
					`const b = {};
const a = b;
a.Foo = class {
  static bar = 1;
};
alert(b.Foo.bar);
`)
			})
	// Throwables.getRootCause(e)
	var root jutil.Throwable = e
	for root.GetCause() != nil {
		root = root.GetCause()
	}
	IsInstanceOf[*jutil.IllegalStateException](AssertThat(root))
	AssertThat(root.GetMessage()).Contains("unexpected node: MEMBER_FIELD_DEF")
}

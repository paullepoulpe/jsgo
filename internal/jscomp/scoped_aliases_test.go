// Port of com/google/javascript/jscomp/ScopedAliasesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ScopedAliasesTest.java)

/*
 * Copyright 2010 The Closure Compiler Authors.
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

// Tests for ScopedAliases
type ScopedAliasesTest struct {
	CompilerTestCaseBase

	invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling
}

const scopedAliasesTest_GOOG_SCOPE_START_BLOCK = "goog.scope(function() {"
const scopedAliasesTest_GOOG_SCOPE_END_BLOCK = "});"

var scopedAliasesTest_EXTERNS = CompilerTestCase_MINIMAL_EXTERNS + "var window;"

// Java: ScopedAliasesTest()
func NewScopedAliasesTest() *ScopedAliasesTest {
	c := &ScopedAliasesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, scopedAliasesTest_EXTERNS)
	return c
}

func TestScopedAliasesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewScopedAliasesTest() })
}

// Java: ScopedAliasesTest.setUp()
func (c *ScopedAliasesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.DisableTypeCheck()
	c.EnableRunTypeCheckAfterProcessing()
	c.EnableCreateModuleMap()
	c.invalidModuleGetHandling = ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME
	c.SetGenericNameReplacements(jutil.NewImmutableMapBuilder[string, string]().Put("SCOPED_ALIASES", "jscomp$scopedAliases$").BuildOrThrow())
}

// Java: ScopedAliasesTest.testScoped(String, String, Postcondition...)
func (c *ScopedAliasesTest) testScoped(code string, expected string, postconditions ...CompilerTestCase_Postcondition) {
	parts := []any{
		c.Srcs(scopedAliasesTest_GOOG_SCOPE_START_BLOCK + code + scopedAliasesTest_GOOG_SCOPE_END_BLOCK),
		c.Expected(expected),
	}
	for _, p := range postconditions {
		parts = append(parts, p)
	}
	c.Test(parts...)
}

// Java: ScopedAliasesTest.testScopedNoChanges(String, String, Postcondition...)
func (c *ScopedAliasesTest) testScopedNoChanges(aliases string, code string, postconditions ...CompilerTestCase_Postcondition) {
	c.testScoped(aliases+code, code, postconditions...)
}

var scopedAliasesTest_VERIFY_TYPES CompilerTestCase_Postcondition = CompilerTestCase_PostconditionFunc(
	func(compiler *Compiler) {
		NewScopedAliasesTest_TypeVerifyingPass(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	})

// Java: ScopedAliasesTest.testTypes(String, String)
func (c *ScopedAliasesTest) testTypes(aliases string, code string) {
	c.testScopedNoChanges(aliases, code, scopedAliasesTest_VERIFY_TYPES)
}

// Java: ScopedAliasesTest.testScopedError(String, DiagnosticType)
func (c *ScopedAliasesTest) testScopedError(code string, expectedError *DiagnosticType) {
	c.TestError("goog.scope(function() {"+code+"});", expectedError)
}

// Java: ScopedAliasesTest.getProcessor(Compiler)
func (c *ScopedAliasesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewScopedAliases(compiler, ScopedAliasesOptions{
		ModuleMetadataMap:        compiler.GetModuleMetadataMap(),
		InvalidModuleGetHandling: c.invalidModuleGetHandling,
	})
}

// ScopedAliasesTest_TypeVerifyingPass is ScopedAliasesTest.TypeVerifyingPass.
type ScopedAliasesTest_TypeVerifyingPass struct {
	compiler    *Compiler
	actualTypes []rhino.Node
}

// Java: ScopedAliasesTest.TypeVerifyingPass(Compiler)
func NewScopedAliasesTest_TypeVerifyingPass(compiler *Compiler) *ScopedAliasesTest_TypeVerifyingPass {
	return &ScopedAliasesTest_TypeVerifyingPass{compiler: compiler}
}

// Java: ScopedAliasesTest.TypeVerifyingPass.process(Node, Node)
func (p *ScopedAliasesTest_TypeVerifyingPass) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, p)
}

// Java: ScopedAliasesTest.TypeVerifyingPass.shouldTraverse(NodeTraversal, Node, Node)
func (p *ScopedAliasesTest_TypeVerifyingPass) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// Java: ScopedAliasesTest.TypeVerifyingPass.visit(NodeTraversal, Node, Node)
func (p *ScopedAliasesTest_TypeVerifyingPass) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	info := n.GetJSDocInfo()
	if info != nil {
		typeNodes := info.GetTypeNodes()
		if len(typeNodes) != 0 {
			if p.actualTypes != nil {
				expectedTypes := append([]rhino.Node{}, info.GetTypeNodes()...)
				AssertWithMessage("Wrong number of JsDoc types").
					That(len(p.actualTypes)).
					IsEqualTo(len(expectedTypes))
				for i := 0; i < len(expectedTypes); i++ {
					AssertNode(p.actualTypes[i]).IsEqualTo(expectedTypes[i])
				}
			} else {
				p.actualTypes = append([]rhino.Node{}, info.GetTypeNodes()...)
			}
		}
	}
}

// Java: ScopedAliasesTest.testLet()
func (c *ScopedAliasesTest) TestLet() {
	c.testScoped("let d = goog.dom; d.createElement(DIV);", "goog.dom.createElement(DIV)")
}

// Java: ScopedAliasesTest.testConst()
func (c *ScopedAliasesTest) TestConst() {
	c.testScoped("const d = goog.dom; d.createElement(DIV);", "goog.dom.createElement(DIV)")
}

// Java: ScopedAliasesTest.testOneLevel()
func (c *ScopedAliasesTest) TestOneLevel() {
	c.testScoped(
		"var g = goog;g.dom.createElement(g.dom.TagName.DIV);",
		"goog.dom.createElement(goog.dom.TagName.DIV);")
}

// Java: ScopedAliasesTest.testTwoLevel()
func (c *ScopedAliasesTest) TestTwoLevel() {
	c.testScoped(
		"var d = goog.dom;d.createElement(d.TagName.DIV);",
		"goog.dom.createElement(goog.dom.TagName.DIV);")
}

// Java: ScopedAliasesTest.testSourceInfo()
func (c *ScopedAliasesTest) TestSourceInfo() {
	c.testScoped(
		`var d = dom;
var e = event;
alert(e.EventType.MOUSEUP);
alert(d.TagName.DIV);
`,
		"alert(event.EventType.MOUSEUP); alert(dom.TagName.DIV);")
	root := c.GetLastCompiler().GetRoot()
	dom := c.FindQualifiedNameNode("dom", root)
	event := c.FindQualifiedNameNode("event", root)
	AssertWithMessage("Dom line should be after event line.").
		That(dom.GetLineno()).
		IsGreaterThan(event.GetLineno())
}

// Java: ScopedAliasesTest.testTransitive()
func (c *ScopedAliasesTest) TestTransitive() {
	c.testScoped(
		"var d = goog.dom;var DIV = d.TagName.DIV;d.createElement(DIV);",
		"goog.dom.createElement(goog.dom.TagName.DIV);")
}

// Java: ScopedAliasesTest.testTransitiveInSameVar()
func (c *ScopedAliasesTest) TestTransitiveInSameVar() {
	c.testScoped(
		"var d = goog.dom, DIV = d.TagName.DIV;d.createElement(DIV);",
		"goog.dom.createElement(goog.dom.TagName.DIV);")
}

// Java: ScopedAliasesTest.testMultipleTransitive()
func (c *ScopedAliasesTest) TestMultipleTransitive() {
	c.testScoped(
		`var g=goog;var d=g.dom;var t=d.TagName;var DIV=t.DIV;
d.createElement(DIV);
`,
		"goog.dom.createElement(goog.dom.TagName.DIV);")
}

// Java: ScopedAliasesTest.testFourLevel()
func (c *ScopedAliasesTest) TestFourLevel() {
	c.testScoped(
		"var DIV = goog.dom.TagName.DIV;goog.dom.createElement(DIV);",
		"goog.dom.createElement(goog.dom.TagName.DIV);")
}

// Java: ScopedAliasesTest.testWorksInClosures()
func (c *ScopedAliasesTest) TestWorksInClosures() {
	c.testScoped(
		`var DIV = goog.dom.TagName.DIV;
goog.x = function() {goog.dom.createElement(DIV);};
`,
		"goog.x = function() {goog.dom.createElement(goog.dom.TagName.DIV);};")
}

// Java: ScopedAliasesTest.testOverridden()
func (c *ScopedAliasesTest) TestOverridden() {
	// Test that the alias doesn't get unaliased when it's overridden by a
	// parameter.
	c.testScopedNoChanges("var g = goog;", "goog.x = function(g) {g.z()};")
	// Same for a local.
	c.testScopedNoChanges("var g = goog;", "goog.x = function() {var g = {}; g.z()};")
}

// Java: ScopedAliasesTest.testTwoScopes()
func (c *ScopedAliasesTest) TestTwoScopes() {
	c.Test(
		`goog.scope(function() {var g = goog;g.method()});
goog.scope(function() {g.method();});
`,
		"goog.method();g.method();")
}

// Java: ScopedAliasesTest.testTwoSymbolsInTwoScopes()
func (c *ScopedAliasesTest) TestTwoSymbolsInTwoScopes() {
	c.Test(
		`var goog = {};
goog.scope(function() { var g = goog; g.Foo = function() {}; });
goog.scope(function() {
  var Foo = goog.Foo; goog.bar = function() { return new Foo(); };
});
`,
		`var goog = {};
goog.Foo = function() {};
goog.bar = function() { return new goog.Foo(); };
`)
}

// Java: ScopedAliasesTest.testAliasOfSymbolInGoogScope()
func (c *ScopedAliasesTest) TestAliasOfSymbolInGoogScope() {
	c.Test(
		`var goog = {};
goog.scope(function() {
  var g = goog;
  g.Foo = function() {};
  var Foo = g.Foo;
  Foo.prototype.bar = function() {};
});
`,
		`var goog = {}; goog.Foo = function() {};
goog.Foo.prototype.bar = function() {};
`)
}

// Java: ScopedAliasesTest.testScopedFunctionReturnThis()
func (c *ScopedAliasesTest) TestScopedFunctionReturnThis() {
	c.Test(
		`goog.scope(function() {
  var g = goog; g.f = function() { return this; };
});
`,
		"goog.f = function() { return this; };")
}

// Java: ScopedAliasesTest.testScopedFunctionAssignsToVar()
func (c *ScopedAliasesTest) TestScopedFunctionAssignsToVar() {
	c.Test(
		`goog.scope(function() {
  var g = goog; g.f = function(x) { x = 3; return x; };
});
`,
		"goog.f = function(x) { x = 3; return x; };")
}

// Java: ScopedAliasesTest.testScopedFunctionThrows()
func (c *ScopedAliasesTest) TestScopedFunctionThrows() {
	c.Test(
		`goog.scope(function() {
  var g = goog; g.f = function() { throw 'error'; };
});
`,
		"goog.f = function() { throw 'error'; };")
}

// Java: ScopedAliasesTest.testPropertiesNotChanged()
func (c *ScopedAliasesTest) TestPropertiesNotChanged() {
	c.testScopedNoChanges("var x = goog.dom;", "y.x();")
}

// Java: ScopedAliasesTest.testShadowedVar()
func (c *ScopedAliasesTest) TestShadowedVar() {
	c.Test(
		`var Popup = {};
var OtherPopup = {};
goog.scope(function() {
  var Popup = OtherPopup;
  Popup.newMethod = function() { return new Popup(); };
});
`,
		`var Popup = {};
var OtherPopup = {};
OtherPopup.newMethod = function() { return new OtherPopup(); };
`)
}

// Java: ScopedAliasesTest.testShadowedScopedVar()
func (c *ScopedAliasesTest) TestShadowedScopedVar() {
	c.Test(
		`var goog = {};
goog.bar = {};
goog.scope(function() {
  var bar = goog.bar;
  // This is bogus, because when the aliases are expanded, goog will shadow goog.bar.
  bar.newMethod = function(goog) { return goog + bar; };
});
`,
		`var goog={};
goog.bar={};
goog.bar.newMethod=function(goog$jscomp$1){return goog$jscomp$1 + goog.bar}
`)
}

// Java: ScopedAliasesTest.testShadowedScopedVarTwoScopes()
func (c *ScopedAliasesTest) TestShadowedScopedVarTwoScopes() {
	c.Test(
		`var goog = {};
goog.bar = {};
goog.scope(function() {
  var bar = goog.bar;
  bar.newMethod = function(goog, a) { return bar + a; };
});
goog.scope(function() {
  var bar = goog.bar;
  bar.newMethod2 = function(goog, b) { return bar + b; };
});
`,
		`var goog={};
goog.bar={};
goog.bar.newMethod=function(goog$jscomp$1, a){return goog.bar + a};
goog.bar.newMethod2=function(goog$jscomp$1, b){return goog.bar + b};
`)
}

// Java: ScopedAliasesTest.testIconImage2()
func (c *ScopedAliasesTest) TestIconImage2() {
	c.testScoped(
		`var getIconImage = (function(){
 let _iconsImage;
 function getIconImage(){
   _iconsImage = 3;
   return _iconsImage;
 }
})
`,
		`var $jscomp$scope$m1146332801$0$getIconImage = (function(){
 let _iconsImage;
 function getIconImage(){
   _iconsImage = 3;
   return _iconsImage;
 }
})
`)
}

// Java: ScopedAliasesTest.testIconImage()
func (c *ScopedAliasesTest) TestIconImage() {
	c.testScoped(
		`let _iconsImage;
function getIconImage(){
  _iconsImage = 3;
  return _iconsImage;
}
`,
		`var $jscomp$scope$m1146332801$1$getIconImage = function() {
  $jscomp$scope$m1146332801$0$_iconsImage = 3;
  return $jscomp$scope$m1146332801$0$_iconsImage;
};
var $jscomp$scope$m1146332801$0$_iconsImage;
`)
}

// Java: ScopedAliasesTest.testFunctionDeclarationInScope()
func (c *ScopedAliasesTest) TestFunctionDeclarationInScope() {
	c.testScoped("var foo = function() {};", "var $jscomp$scope$m1146332801$0$foo = function() {};")
}

// Java: ScopedAliasesTest.testFunctionDeclarationInScope_letConst()
func (c *ScopedAliasesTest) TestFunctionDeclarationInScope_letConst() {
	c.testScoped(
		"var baz = goog.bar; let foo = function() {return baz;};",
		"var $jscomp$scope$m1146332801$0$foo = function() {return goog.bar;};")
	c.testScoped(
		"var baz = goog.bar; const foo = function() {return baz;};",
		"var $jscomp$scope$m1146332801$0$foo = function() {return goog.bar;};")
}

// Java: ScopedAliasesTest.testLetConstShadowing()
func (c *ScopedAliasesTest) TestLetConstShadowing() {
	c.testScoped(
		`var foo = goog.bar; var f = function() {
let foo = baz; return foo;};
`,
		`var $jscomp$scope$m1146332801$0$f = function() {
let foo = baz; return foo;};
`)
	c.testScoped(
		`var foo = goog.bar; var f = function() {
const foo = baz; return foo;};
`,
		`var $jscomp$scope$m1146332801$0$f = function() {
const foo = baz; return foo;};
`)
}

// Java: ScopedAliasesTest.testYieldExpression()
func (c *ScopedAliasesTest) TestYieldExpression() {
	c.testScoped(
		"var foo = goog.bar; var f = function*() {yield foo;};",
		"var $jscomp$scope$m1146332801$0$f = function*() {yield goog.bar;};")
}

// Java: ScopedAliasesTest.testDestructuringError()
func (c *ScopedAliasesTest) TestDestructuringError() {
	c.testScopedError("var [x] = [1];", ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL)
}

// Java: ScopedAliasesTest.testObjectDescructuringError1()
func (c *ScopedAliasesTest) TestObjectDescructuringError1() {
	c.testScopedError("var {x} = {x: 1};", ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL)
}

// Java: ScopedAliasesTest.testObjectDescructuringError2()
func (c *ScopedAliasesTest) TestObjectDescructuringError2() {
	c.testScopedError("var {x: y} = {x: 1};", ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL)
}

// Java: ScopedAliasesTest.testNonTopLevelDestructuring()
func (c *ScopedAliasesTest) TestNonTopLevelDestructuring() {
	c.testScoped(
		"var f = function() {var [x, y] = [1, 2];};",
		"var $jscomp$scope$m1146332801$0$f = function() {var [x, y] = [1, 2];};")
}

// Java: ScopedAliasesTest.testArrowFunction()
func (c *ScopedAliasesTest) TestArrowFunction() {
	c.testScoped(
		"var foo = goog.bar; var v = (x => x + foo);",
		"var $jscomp$scope$m1146332801$0$v = (x => x + goog.bar)")
}

// Java: ScopedAliasesTest.testClassDefinition1()
func (c *ScopedAliasesTest) TestClassDefinition1() {
	c.testScoped("class Foo {}", "var $jscomp$scope$m1146332801$0$Foo=class{}")
}

// Java: ScopedAliasesTest.testClassDefinition2()
func (c *ScopedAliasesTest) TestClassDefinition2() {
	c.testScoped(
		`var bar = goog.bar;
class Foo { constructor() { this.x = bar; }}
`,
		"var $jscomp$scope$m1146332801$0$Foo = class { constructor() { this.x = goog.bar; } }")
}

// Java: ScopedAliasesTest.testClassDefinition3()
func (c *ScopedAliasesTest) TestClassDefinition3() {
	c.testScoped(
		`var bar = {};
bar.Foo = class {};
`,
		`var $jscomp$scope$m1146332801$0$bar = {}; $jscomp$scope$m1146332801$0$bar.Foo = class
 {}
`)
}

// Java: ScopedAliasesTest.testClassDefinition_letConst()
func (c *ScopedAliasesTest) TestClassDefinition_letConst() {
	c.testScoped(
		`let bar = {};
bar.Foo = class {};
`,
		`var $jscomp$scope$m1146332801$0$bar = {}; $jscomp$scope$m1146332801$0$bar.Foo = class
 {}
`)
	c.testScoped(
		`const bar = {};
bar.Foo = class {};
`,
		`var $jscomp$scope$m1146332801$0$bar = {}; $jscomp$scope$m1146332801$0$bar.Foo = class
 {}
`)
}

// Java: ScopedAliasesTest.testDefaultParameter()
func (c *ScopedAliasesTest) TestDefaultParameter() {
	c.testScoped(
		"var foo = goog.bar; var f = function(y=foo) {};",
		"var $jscomp$scope$m1146332801$0$f = function(y=goog.bar) {};")
}

// Java: ScopedAliasesTest.testObjectLiteral()
func (c *ScopedAliasesTest) TestObjectLiteral() {
	c.testScoped(
		`var Foo = goog.Foo;
goog.x = {
  /** @param {Foo} foo */
  y: function(foo) { }
};
`,
		`goog.x = {
  /** @param {goog.Foo} foo */
  y: function(foo) {}
};
`)

	c.testScoped(
		`var Foo = goog.Foo;
goog.x = {
  y: /** @param {Foo} foo */ function(foo) {}
};
`,
		`goog.x = {
  y: /** @param {goog.Foo} foo */ function(foo) {}
};
`)

	c.testScoped(
		`var Foo = goog.Foo;
goog.x = {
  y: /** @type {function(Foo)} */ (function(foo) {})
};
`,
		`goog.x = {
  y: /** @type {function(goog.Foo)} */ (function(foo) {})
};
`)
}

// Java: ScopedAliasesTest.testObjectLiteralShorthand()
func (c *ScopedAliasesTest) TestObjectLiteralShorthand() {
	c.testScoped(
		"var bar = goog.bar; var Foo = {bar};",
		"var $jscomp$scope$m1146332801$0$Foo={bar: goog.bar};")
}

// Java: ScopedAliasesTest.testObjectLiteralMethods()
func (c *ScopedAliasesTest) TestObjectLiteralMethods() {
	c.testScoped(
		"var foo = goog.bar; var obj = {toString() {return foo}};",
		"var $jscomp$scope$m1146332801$0$obj = {toString() {return goog.bar}};")
}

// Java: ScopedAliasesTest.testObjectLiteralComputedPropertyNames()
func (c *ScopedAliasesTest) TestObjectLiteralComputedPropertyNames() {
	c.testScoped(
		"var foo = goog.bar; var obj = {[(() => foo)()]: baz};",
		"var $jscomp$scope$m1146332801$0$obj = {[(() => goog.bar)()]:baz};")
	c.testScoped(
		"var foo = goog.bar; var obj = {[x => x + foo]: baz};",
		"var $jscomp$scope$m1146332801$0$obj = {[x => x + goog.bar]:baz};")
}

// Java: ScopedAliasesTest.testJsDocNotIgnored()
func (c *ScopedAliasesTest) TestJsDocNotIgnored() {
	c.EnableTypeCheck()

	externs := CompilerTestCase_MINIMAL_EXTERNS +
		`/** @const */ var ns = {};
/** @constructor */
ns.Foo;

var goog;
/** @param {function()} fn */
goog.scope = function(fn) {}
`

	js := `goog.scope(function() {
  var Foo = ns.Foo;
  var x = {
    /** @param {Foo} foo */ y: function(foo) {}
  };
  x.y('');
});
`
	c.Test(c.Externs(externs), c.Srcs(js), c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))

	js =
		`goog.scope(function() {
  var Foo = ns.Foo;
  var x = {
    y: /** @param {Foo} foo */ function(foo) {}
  };
  x.y('');
});
`
	c.Test(c.Externs(externs), c.Srcs(js), c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))
}

// Java: ScopedAliasesTest.testUsingObjectLiteralToEscapeScoping()
func (c *ScopedAliasesTest) TestUsingObjectLiteralToEscapeScoping() {
	// There are many ways to shoot yourself in the foot with goog.scope
	// and make the compiler generate bad code. We generally don't care.
	//
	// We only try to protect against accidental mis-use, not deliberate
	// mis-use.
	c.Test(
		`var goog = {};
goog.bar = {};
goog.scope(function() {
  var bar = goog.bar;
  var baz = goog.bar.baz;
  goog.foo = function() {
    goog.bar = {baz: 3};
    return baz;
  };
});
`,
		`var goog = {};
goog.bar = {};
goog.foo = function(){
  goog.bar = {baz:3};
  return goog.bar.baz;
};
`)
}

// Java: ScopedAliasesTest.testJsDocType()
func (c *ScopedAliasesTest) TestJsDocType() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @type {goog.Timer} */ types.actual;
/** @type {goog.Timer} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testJsDocParameter()
func (c *ScopedAliasesTest) TestJsDocParameter() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @param {goog.Timer} a */ types.actual;
/** @param {goog.Timer} a */ types.expected;
`)
}

// Java: ScopedAliasesTest.testJsDocExtends()
func (c *ScopedAliasesTest) TestJsDocExtends() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @extends {goog.Timer} */ types.actual;
/** @extends {goog.Timer} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testJsDocImplements()
func (c *ScopedAliasesTest) TestJsDocImplements() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @implements {goog.Timer} */ types.actual;
/** @implements {goog.Timer} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testJsDocEnum()
func (c *ScopedAliasesTest) TestJsDocEnum() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @enum {goog.Timer} */ types.actual;
/** @enum {goog.Timer} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testJsDocReturn()
func (c *ScopedAliasesTest) TestJsDocReturn() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @return {goog.Timer} */ types.actual;
/** @return {goog.Timer} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testJsDocThis()
func (c *ScopedAliasesTest) TestJsDocThis() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @this {goog.Timer} */ types.actual;
/** @this {goog.Timer} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testJsDocThrows()
func (c *ScopedAliasesTest) TestJsDocThrows() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @throws {goog.Timer} */ types.actual;
/** @throws {goog.Timer} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testJsDocSubType()
func (c *ScopedAliasesTest) TestJsDocSubType() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @type {goog.Timer.Enum} */ types.actual;
/** @type {goog.Timer.Enum} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testJsDocTypedef()
func (c *ScopedAliasesTest) TestJsDocTypedef() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @typedef {goog.Timer} */ types.actual;
/** @typedef {goog.Timer} */ types.expected;
`)

	c.testScoped(
		`/** @typedef {string} */ var s;
/** @type {s} */ var t;
`,
		`/** @typedef {string} */ var $jscomp$scope$m1146332801$0$s;
/** @type {$jscomp$scope$m1146332801$0$s} */ var $jscomp$scope$m1146332801$1$t;
`)

	c.testScoped(
		`/** @typedef {string} */ let s;
/** @type {s} */ var t;
`,
		`/** @typedef {string} */ var $jscomp$scope$m1146332801$0$s;
/** @type {$jscomp$scope$m1146332801$0$s} */ var $jscomp$scope$m1146332801$1$t;
`)
}

// Java: ScopedAliasesTest.testJsDocRecord()
func (c *ScopedAliasesTest) TestJsDocRecord() {
	c.EnableTypeCheck()
	c.Test(
		`/** @const */ var ns = {};
goog.scope(function () {
  var x = goog.Timer;
  /** @type {{x: string}} */ ns.y = {'goog.Timer': 'x'};
});
`,
		`/** @const */ var ns = {};
/** @type {{x: string}} */ ns.y = {'goog.Timer': 'x'};
`,
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))
}

// Java: ScopedAliasesTest.testArrayJsDoc()
func (c *ScopedAliasesTest) TestArrayJsDoc() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @type {Array.<goog.Timer>} */ types.actual;
/** @type {Array.<goog.Timer>} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testObjectJsDoc()
func (c *ScopedAliasesTest) TestObjectJsDoc() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @type {{someKey: goog.Timer}} */ types.actual;
/** @type {{someKey: goog.Timer}} */ types.expected;
`)
	c.testTypes(
		"var x = goog.Timer;",
		`/** @type {{x: number}} */ types.actual;
/** @type {{x: number}} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testObjectJsDoc2()
func (c *ScopedAliasesTest) TestObjectJsDoc2() {
	c.testTypes(
		"var x = goog$Timer;",
		`/** @type {{someKey: goog$Timer}} */ types.actual;
/** @type {{someKey: goog$Timer}} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testUnionJsDoc()
func (c *ScopedAliasesTest) TestUnionJsDoc() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @type {goog.Timer|Object} */ types.actual;
/** @type {goog.Timer|Object} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testFunctionJsDoc()
func (c *ScopedAliasesTest) TestFunctionJsDoc() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @type {function(goog.Timer) : void} */ types.actual;
/** @type {function(goog.Timer) : void} */ types.expected;
`)
	c.testTypes(
		"var x = goog.Timer;",
		`/** @type {function() : goog.Timer} */ types.actual;
/** @type {function() : goog.Timer} */ types.expected;
`)
}

// Java: ScopedAliasesTest.testForwardJsDoc()
func (c *ScopedAliasesTest) TestForwardJsDoc() {
	c.testScoped(
		`/**
 * @constructor
 */
foo.Foo = function() {};
/** @param {Foo.Bar} x */ foo.Foo.actual = function(x) {3};
var Foo = foo.Foo;
/** @constructor */ Foo.Bar = function() {};
/** @param {foo.Foo.Bar} x */ foo.Foo.expected = function(x) {};
`,
		`/**
 * @constructor
 */
foo.Foo = function() {};
/** @param {foo.Foo.Bar} x */ foo.Foo.actual = function(x) {3};
/** @constructor */ foo.Foo.Bar = function() {};
/** @param {foo.Foo.Bar} x */ foo.Foo.expected = function(x) {};
`,
		scopedAliasesTest_VERIFY_TYPES)
}

// Java: ScopedAliasesTest.testTestTypes()
func (c *ScopedAliasesTest) TestTestTypes() {
	jutil.Try(func() {
		c.testTypes(
			"var x = goog.Timer;",
			`/** @type {function() : x} */ types.actual;
/** @type {function() : wrong.wrong} */ types.expected;
`)
		panic(jutil.NewJavaLangError("Test types should fail here."))
	}, func(expected *jutil.AssertionError) {
	})
}

// Java: ScopedAliasesTest.testNullType()
func (c *ScopedAliasesTest) TestNullType() {
	c.testTypes(
		"var x = goog.Timer;",
		`/** @param draggable */ types.actual;
/** @param draggable */ types.expected;
`)
}

// Java: ScopedAliasesTest.testJSDocCopiedForFunctions()
func (c *ScopedAliasesTest) TestJSDocCopiedForFunctions() {
	c.testScoped(
		"/** @export */ function Foo() {}",
		"/** @export */ var $jscomp$scope$m1146332801$0$Foo =/** @export */ function() {}")
}

// Java: ScopedAliasesTest.testJSDocCopiedForClasses()
func (c *ScopedAliasesTest) TestJSDocCopiedForClasses() {
	c.testScoped(
		"/** @export */ class Foo {}",
		"/** @export */ var $jscomp$scope$m1146332801$0$Foo = /** @export */ class {}")
}

// Java: ScopedAliasesTest.testIssue772()
func (c *ScopedAliasesTest) TestIssue772() {
	c.testTypes(
		`var b = a.b;
var c = b.c;
`,
		`/** @param {a.b.c.MyType} x */ types.actual;
/** @param {a.b.c.MyType} x */ types.expected;
`)
}

// Java: ScopedAliasesTest.testInlineJsDoc()
func (c *ScopedAliasesTest) TestInlineJsDoc() {
	c.EnableTypeCheck()
	c.Test(
		c.Srcs(
			`/** @const */ var ns = {};
/** @constructor */ ns.A = function() {};
goog.scope(function() {
  /** @const */ var A = ns.A;
  var /** @type {?A} */ b = null;
});
`),
		c.Expected(
			`/** @const */ var ns = {};
/** @constructor */ ns.A = function() {};
/** @type {?ns.A} */ var $jscomp$scope$m1146332801$0$b = null;
`),
		scopedAliasesTest_VERIFY_TYPES)
}

// Java: ScopedAliasesTest.testInlineReturn()
func (c *ScopedAliasesTest) TestInlineReturn() {
	c.EnableTypeCheck()
	c.Test(
		c.Srcs(
			`/** @const */ var ns = {};
/** @constructor */ ns.A = function() {};
goog.scope(function() {
  /** @const */ var A = ns.A;
  function /** ?A */ b() { return null; }
});
`),
		c.Expected(
			`/** @const */ var ns = {};
/** @constructor */ ns.A = function() {};
// TODO(moz): See if we can avoid generating duplicate @return's
/** @return {?ns.A} */ var $jscomp$scope$m1146332801$0$b =
    /** @return {?ns.A} */ function() { return null; };
`),
		scopedAliasesTest_VERIFY_TYPES)
}

// Java: ScopedAliasesTest.testInlineParam()
func (c *ScopedAliasesTest) TestInlineParam() {
	c.EnableTypeCheck()
	c.Test(
		c.Srcs(
			`/** @const */ var ns = {};
/** @constructor */ ns.A = function() {};
goog.scope(function() {
  /** @const */ var A = ns.A;
  function b(/** ?A */ bee) {}
});
`),
		c.Expected(
			`/** @const */ var ns = {};
/** @constructor */ ns.A = function() {};
var $jscomp$scope$m1146332801$0$b = function(/** ?ns.A */ bee) {};
`),
		scopedAliasesTest_VERIFY_TYPES)
}

// Java: ScopedAliasesTest.testScopedThis()
func (c *ScopedAliasesTest) TestScopedThis() {
	c.testScopedError("this.y = 10;", ScopedAliases_GOOG_SCOPE_REFERENCES_THIS)
	c.testScopedError("var x = this;", ScopedAliases_GOOG_SCOPE_REFERENCES_THIS)
	c.testScopedError("fn(this);", ScopedAliases_GOOG_SCOPE_REFERENCES_THIS)
}

// Java: ScopedAliasesTest.testAliasRedefinition()
func (c *ScopedAliasesTest) TestAliasRedefinition() {
	c.testScopedError("var x = goog.dom; x = goog.events;", ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED)
}

// Java: ScopedAliasesTest.testAliasNonRedefinition()
func (c *ScopedAliasesTest) TestAliasNonRedefinition() {
	c.Test("var y = {}; goog.scope(function() { goog.dom = y; });", "var y = {}; goog.dom = y;")
}

// Java: ScopedAliasesTest.testCtorAlias()
func (c *ScopedAliasesTest) TestCtorAlias() {
	c.Test(
		`var x = {y: {}};
goog.scope(function() {
  var y = x.y;
  y.ClassA = function() { this.b = new ClassB(); };
  y.ClassB = function() {};
  var ClassB = y.ClassB;
});
`,
		`var x = {y: {}};
x.y.ClassA = function() { this.b = new x.y.ClassB(); };
x.y.ClassB = function() { };
`)
}

// Java: ScopedAliasesTest.testAliasCycle()
func (c *ScopedAliasesTest) TestAliasCycle() {
	c.TestError(
		`var x = {y: {}};
goog.scope(function() {
  var y = z.x;
  var z = y.x;
  y.ClassA = function() {};
  z.ClassB = function() {};
});
`,
		ScopedAliases_GOOG_SCOPE_ALIAS_CYCLE)
}

// Java: ScopedAliasesTest.testScopedReturn()
func (c *ScopedAliasesTest) TestScopedReturn() {
	c.testScopedError("return;", ScopedAliases_GOOG_SCOPE_USES_RETURN)
	c.testScopedError("var x = goog.dom; return;", ScopedAliases_GOOG_SCOPE_USES_RETURN)
}

// Java: ScopedAliasesTest.testScopedThrow()
func (c *ScopedAliasesTest) TestScopedThrow() {
	c.testScopedError("throw 'error';", ScopedAliases_GOOG_SCOPE_USES_THROW)
}

// Java: ScopedAliasesTest.testUsedImproperly()
func (c *ScopedAliasesTest) TestUsedImproperly() {
	c.TestError("var x = goog.scope(function() {});", ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE)
	c.TestError(
		"var f = function() { goog.scope(function() {}); }",
		ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE)
}

// Java: ScopedAliasesTest.testScopeCallInIf()
func (c *ScopedAliasesTest) TestScopeCallInIf() {
	c.Test("if (true) { goog.scope(function() {});}", "if (true) {}")
	c.Test("if (true) { goog.scope(function()  { var x = foo; });}", "if (true) { }")
	c.Test(
		"if (true) { goog.scope(function()  { var x = foo; console.log(x); });}",
		"if (true) { console.log(foo); }")
}

// Java: ScopedAliasesTest.testBadParameters()
func (c *ScopedAliasesTest) TestBadParameters() {
	c.TestError("goog.scope()", ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS)
	c.TestError("goog.scope(10)", ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS)
	c.TestError("goog.scope(function() {}, 10)", ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS)
	c.TestError("goog.scope(function z() {})", ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS)
	c.TestError("goog.scope(function(a, b, c) {})", ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS)
}

// Java: ScopedAliasesTest.testNonAliasLocal()
func (c *ScopedAliasesTest) TestNonAliasLocal() {
	c.testScopedError("for (var k in { a: 1, b: 2 }) {}", ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL)
	c.testScopedError("for (var k of [1, 2, 3]) {}", ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL)
}

// Java: ScopedAliasesTest.testInvalidVariableInScope()
func (c *ScopedAliasesTest) TestInvalidVariableInScope() {
	c.testScopedError("if (true) { function f() {}}", ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE)
	c.testScopedError("for (;;) { function f() {}}", ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE)
}

// Java: ScopedAliasesTest.testWithCatch1()
func (c *ScopedAliasesTest) TestWithCatch1() {
	c.testScoped(
		"var x = foo(); try { } catch (e) {}",
		"var $jscomp$scope$m1146332801$0$x = foo(); try { } catch (e) {}")
}

// Java: ScopedAliasesTest.testWithCatch2()
func (c *ScopedAliasesTest) TestWithCatch2() {
	c.testScoped(
		"try { } catch (e) {var x = foo();}",
		"try { } catch (e) {var $jscomp$scope$m1146332801$0$x = foo();}")
}

// Java: ScopedAliasesTest.testVariablesInCatchBlock()
func (c *ScopedAliasesTest) TestVariablesInCatchBlock() {
	c.testScopedNoChanges("", "try {} catch (e) {}")
	c.testScopedNoChanges("", "try {} catch (e) { let x = foo }")
	c.testScopedNoChanges("", "try {} catch (e) { const x = foo }")
}

// Java: ScopedAliasesTest.testLetConstInBlock()
func (c *ScopedAliasesTest) TestLetConstInBlock() {
	c.testScopedNoChanges("", "if (true) {let x = foo;}")
	c.testScopedNoChanges("", "if (true) {const x = foo;}")
}

// Java: ScopedAliasesTest.testHoistedAliases()
func (c *ScopedAliasesTest) TestHoistedAliases() {
	c.testScoped("if (true) { var x = foo;}", "if (true) {}")
	c.testScoped("if (true) { var x = foo; console.log(x); }", "if (true) { console.log(foo); }")
}

// Java: ScopedAliasesTest.testOkAliasLocal()
func (c *ScopedAliasesTest) TestOkAliasLocal() {
	c.testScoped("var x = 10;", "var $jscomp$scope$m1146332801$0$x = 10")
	c.testScoped("var x = goog['dom'];", "var $jscomp$scope$m1146332801$0$x = goog['dom']")
	c.testScoped(
		"var x = 10, y = 9;",
		"var $jscomp$scope$m1146332801$0$x = 10; var $jscomp$scope$m1146332801$1$y = 9;")
	c.testScoped(
		"var x = 10, y = 9; goog.getX = function () { return x + y; }",
		`var $jscomp$scope$m1146332801$0$x = 10; var $jscomp$scope$m1146332801$1$y = 9;
goog.getX = function () {
    return $jscomp$scope$m1146332801$0$x + $jscomp$scope$m1146332801$1$y; }
`)
}

// Java: ScopedAliasesTest.testOkAliasLocal_letConst()
func (c *ScopedAliasesTest) TestOkAliasLocal_letConst() {
	c.testScoped("let x = 10;", "var $jscomp$scope$m1146332801$0$x = 10")
	c.testScoped("const x = 10;", "var $jscomp$scope$m1146332801$0$x = 10")
}

// Java: ScopedAliasesTest.testHoistedFunctionDeclaration()
func (c *ScopedAliasesTest) TestHoistedFunctionDeclaration() {
	c.testScoped(
		" g(f); function f() {} ",
		` var $jscomp$scope$m1146332801$0$f = function () {};
g($jscomp$scope$m1146332801$0$f); 
`)
}

// Java: ScopedAliasesTest.testAliasReassign()
func (c *ScopedAliasesTest) TestAliasReassign() {
	c.testScopedError("var x = 3; x = 5;", ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED)
}

// Java: ScopedAliasesTest.testMultipleLocals()
func (c *ScopedAliasesTest) TestMultipleLocals() {
	c.Test(
		`goog.scope(function () { var x = 3; });
goog.scope(function () { var x = 4; });
`,
		"var $jscomp$scope$m1146332801$0$x = 3; var $jscomp$scope$m1146332801$1$x = 4")
}

// Java: ScopedAliasesTest.testIssue1103a()
func (c *ScopedAliasesTest) TestIssue1103a() {
	c.Test(
		`goog.scope(function () {
  var a;
  foo.bar = function () { a = 1; };
});
`,
		`var $jscomp$scope$m1146332801$0$a;
foo.bar = function () {
  $jscomp$scope$m1146332801$0$a = 1;
}
`)
}

// Java: ScopedAliasesTest.testIssue1103b()
func (c *ScopedAliasesTest) TestIssue1103b() {
	c.Test(
		`goog.scope(function () {
  var a = foo, b, c = 1;
});
`,
		"var $jscomp$scope$m1146332801$0$b; var $jscomp$scope$m1146332801$1$c=1")
}

// Java: ScopedAliasesTest.testIssue1103c()
func (c *ScopedAliasesTest) TestIssue1103c() {
	c.Test(
		`goog.scope(function () {
  /** @type {number} */ var a;
});
`,
		"/** @type {number} */ var $jscomp$scope$m1146332801$0$a;")
}

// Java: ScopedAliasesTest.testIssue1144()
func (c *ScopedAliasesTest) TestIssue1144() {
	c.Test(
		`var ns = {};
ns.sub = {};
/** @constructor */ ns.sub.C = function () {};
goog.scope(function () {
  var sub = ns.sub;
  /** @type {sub.C} */
  var x = null;
});
`,
		`var ns = {};
ns.sub = {};
/** @constructor */ ns.sub.C = function () {};
/** @type {ns.sub.C} */
var $jscomp$scope$m1146332801$0$x = null;
`)
}

// https://github.com/google/closure-compiler/issues/2211
// Java: ScopedAliasesTest.testIssue2211()
func (c *ScopedAliasesTest) TestIssue2211() {
	c.Test(
		"var ns = {};\n"+
			"var y = 1;\n"+
			"goog.scope(function () {\n"+
			"  ns.fact = function y(n) {\n"+
			"    return n == 1 ? 1 : n * y(n - 1);\n"+
			"  };\n"+
			"});\n",
		"var ns = {};\n"+
			"var y = 1;\n"+
			"ns.fact = function y$SCOPED_ALIASES$0(n) {\n"+
			"  return n == 1 ? 1 : n * y$SCOPED_ALIASES$0(n - 1);\n"+
			"};\n")
}

// https://github.com/google/closure-compiler/issues/2211
// Java: ScopedAliasesTest.testIssue2211b()
func (c *ScopedAliasesTest) TestIssue2211b() {
	c.Test(
		"var ns = {};\n"+
			"var y = 1;\n"+
			"goog.scope(function () {\n"+
			"  function x(y) {}\n"+
			"  ns.fact = function y(n) {\n"+
			"    return n == 1 ? 1 : n * y(n - 1);\n"+
			"  };\n"+
			"});\n",
		"var ns = {};\n"+
			"var y = 1;\n"+
			"var $jscomp$scope$m1146332801$0$x = function (y) {};\n"+
			"ns.fact = function y$SCOPED_ALIASES$1(n) {\n"+
			"  return n == 1 ? 1 : n * y$SCOPED_ALIASES$1(n - 1);\n"+
			"};\n")
}

// https://github.com/google/closure-compiler/issues/2211
// Java: ScopedAliasesTest.testIssue2211c()
func (c *ScopedAliasesTest) TestIssue2211c() {
	c.testScoped(
		"foo(() => {\n"+
			"  const y = function y() {\n"+
			"    use(y);\n"+
			"  };\n"+
			"});\n",
		"foo(() => {\n"+
			"  const y = function y$SCOPED_ALIASES$0() {\n"+
			"    use(y$SCOPED_ALIASES$0);\n"+
			"  };\n"+
			"});\n")
}

// Java: ScopedAliasesTest.testGoogModuleGet1()
func (c *ScopedAliasesTest) TestGoogModuleGet1() {
	c.Test(
		`goog.provide('provided');
goog.provide('other.name.Foo');
goog.scope(function() {
  var Foo = goog.module.get('other.name.Foo')
  /** @type {!Foo} */
  provided.f = new Foo;
});
`,
		`goog.provide('provided');
goog.provide('other.name.Foo');
/** @type {!other.name.Foo} */
provided.f = new (goog.module.get('other.name.Foo'));
`)
}

// Java: ScopedAliasesTest.testGoogModuleGet2()
func (c *ScopedAliasesTest) TestGoogModuleGet2() {
	c.Test(
		`goog.provide('foo.baz');
goog.provide('other.thing');

goog.scope(function() {

const a = goog.module.get('other.thing');
const b = a.b;
foo.baz = b

}); // goog.scope
`,
		`goog.provide('foo.baz');
goog.provide('other.thing');
foo.baz = goog.module.get('other.thing').b;
`)
}

// Java: ScopedAliasesTest.testGoogModuleGet3()
func (c *ScopedAliasesTest) TestGoogModuleGet3() {
	c.Test(
		`goog.provide('foo.baz');
goog.provide('other.thing');

goog.scope(function() {

const a = goog.module.get('other.thing').b;
foo.baz = a

}); // goog.scope
`,
		`goog.provide('foo.baz');
goog.provide('other.thing');
foo.baz = goog.module.get('other.thing').b;
`)
}

// Java: ScopedAliasesTest.testGoogModuleGet_ofModule()
func (c *ScopedAliasesTest) TestGoogModuleGet_ofModule() {
	c.Test(
		c.Srcs(
			"goog.module('other.name.Foo');",
			`goog.provide('provided');
goog.scope(function() {
  var Foo = goog.module.get('other.name.Foo')
  /** @type {!Foo} */
  provided.f = new Foo;
});
`),
		c.Expected(
			"goog.module('other.name.Foo');",
			`goog.provide('provided');
/** @type {!other.name.Foo} */
provided.f = new (goog.module.get('other.name.Foo'));
`))
}

// Java: ScopedAliasesTest.testGoogModuleGet_missing_noDeletion()
func (c *ScopedAliasesTest) TestGoogModuleGet_missing_noDeletion() {
	c.invalidModuleGetHandling = ScopedAliases_InvalidModuleGetHandling_PRESERVE
	c.Test(
		`goog.provide('provided');

goog.scope(function() {
  var Foo = goog.module.get('missing.name.Foo');
  /** @type {!Foo} */
  provided.f = new Foo;
});
`,
		`goog.provide('provided');

/** @type {!missing.name.Foo} */
provided.f = new (goog.module.get('missing.name.Foo'));
`)
}

// Java: ScopedAliasesTest.testGoogModuleGet_missing()
func (c *ScopedAliasesTest) TestGoogModuleGet_missing() {
	c.Test(
		`goog.provide('provided');

goog.scope(function() {
  var Foo = goog.module.get('missing.name.Foo')
  /** @type {!Foo} */
  provided.f = new Foo;
});
`,
		`goog.provide('provided');

/** @type {!Foo$$jscomp$missingAlias$module$testcode} */
provided.f = new Foo$$jscomp$missingAlias$module$testcode();
`)
}

// Java: ScopedAliasesTest.testGoogModuleGet_missingUsedInRhs()
func (c *ScopedAliasesTest) TestGoogModuleGet_missingUsedInRhs() {
	c.Test(
		`goog.provide('provided');

goog.scope(function() {
  var Foo = goog.module.get('missing.name.Foo')
  provided.Bar = class extends Foo {};
});
`,
		`goog.provide('provided');

provided.Bar = class extends Foo$$jscomp$missingAlias$module$testcode {};
`)
}

// Java: ScopedAliasesTest.testGoogModuleGet_missingRecursive()
func (c *ScopedAliasesTest) TestGoogModuleGet_missingRecursive() {
	c.Test(
		`goog.provide('provided');

goog.scope(function() {
  var Foo = goog.module.get('missing.name.Foo')
  /** @type {!Foo.Bar} */
  provided.f = new Foo.Bar;
});
`,
		`goog.provide('provided');

/** @type {!Foo$$jscomp$missingAlias$module$testcode.Bar} */
provided.f = new Foo$$jscomp$missingAlias$module$testcode.Bar();
`)
}

// Java: ScopedAliasesTest.testGoogModuleGet_missingProperty()
func (c *ScopedAliasesTest) TestGoogModuleGet_missingProperty() {
	c.Test(
		`goog.provide('provided');

goog.scope(function() {
  var Bar = goog.module.get('missing.name.Foo').Bar
  /** @type {!Bar} */
  provided.f = new Bar;
});
`,
		`goog.provide('provided');

/** @type {!Bar$$jscomp$missingAlias$module$testcode} */
provided.f = new Bar$$jscomp$missingAlias$module$testcode();
`)
}

// Java: ScopedAliasesTest.testGoogModuleGet_aliasOfDeletedModuleGetAlias()
func (c *ScopedAliasesTest) TestGoogModuleGet_aliasOfDeletedModuleGetAlias() {
	c.Test(
		`goog.provide('provided');

goog.scope(function() {
  var Foo = goog.module.get('missing.name.Foo');
  var Bar = Foo.Bar
  /** @type {!Bar} */
  provided.f = new Bar;
});
`,
		`goog.provide('provided');

/** @type {!Foo$$jscomp$missingAlias$module$testcode.Bar} */
provided.f = new Foo$$jscomp$missingAlias$module$testcode.Bar();
`)
}

// Java: ScopedAliasesTest.testObjectPattern()
func (c *ScopedAliasesTest) TestObjectPattern() {
	c.testScopedNoChanges("", "{foo: ({bar}) => baz};")
}

// Java: ScopedAliasesTest.testTypeCheck()
func (c *ScopedAliasesTest) TestTypeCheck() {
	c.EnableTypeCheck()

	c.Test(
		`goog.scope(function () {
  /** @constructor */ function F() {}
  /** @return {F} */ function createFoo() { return 1; }
});
`,
		`/** @return {$jscomp$scope$m1146332801$0$F} */
var $jscomp$scope$m1146332801$1$createFoo =
/** @return {$jscomp$scope$m1146332801$0$F} */
function() { return 1; };
/** @constructor */ var $jscomp$scope$m1146332801$0$F =
/** @constructor */ function() { };
`,
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))
}

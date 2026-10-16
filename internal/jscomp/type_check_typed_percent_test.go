// Port of com/google/javascript/jscomp/TypeCheckTypedPercentTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/TypeCheckTypedPercentTest.java)

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

// TypeCheckTypedPercentTest tests the type checker.
type TypeCheckTypedPercentTest struct {
	typeCheckTests_noFixture
}

// Java: TypeCheckTypedPercentTest()
func NewTypeCheckTypedPercentTest() *TypeCheckTypedPercentTest {
	c := &TypeCheckTypedPercentTest{}
	return c
}

func TestTypeCheckTypedPercentTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckTypedPercentTest() })
}

// Java: TypeCheckTypedPercentTest.testGetTypedPercent1()
func (c *TypeCheckTypedPercentTest) TestGetTypedPercent1() {
	js := `var id = function(x) { return x; }
var id2 = function(x) { return id(x); }
`
	AssertThat(c.GetTypedPercent(js)).IsWithin(0.1).
		Of(50.0)
}

// Java: TypeCheckTypedPercentTest.testGetTypedPercent2()
func (c *TypeCheckTypedPercentTest) TestGetTypedPercent2() {
	js := "var x = {}; x.y = 1;"
	AssertThat(c.GetTypedPercent(js)).IsWithin(0.1).
		Of(100.0)
}

// Java: TypeCheckTypedPercentTest.testGetTypedPercent3()
func (c *TypeCheckTypedPercentTest) TestGetTypedPercent3() {
	js := "var f = function(x) { x.a = x.b; }"
	AssertThat(c.GetTypedPercent(js)).IsWithin(0.1).
		Of(25.0)
}

// Java: TypeCheckTypedPercentTest.testGetTypedPercent4()
func (c *TypeCheckTypedPercentTest) TestGetTypedPercent4() {
	js := `var n = {};
/** @constructor */ n.T = function() {};
/** @type {n.T} */ var x = new n.T();
`
	AssertThat(c.GetTypedPercent(js)).IsWithin(0.1).
		Of(100.0)
}

// Java: TypeCheckTypedPercentTest.testGetTypedPercent5()
func (c *TypeCheckTypedPercentTest) TestGetTypedPercent5() {
	js := "/** @enum {number} */ keys = {A: 1,B: 2,C: 3};"
	AssertThat(c.GetTypedPercent(js)).IsWithin(0.1).
		Of(100.0)
}

// Java: TypeCheckTypedPercentTest.testGetTypedPercent6()
func (c *TypeCheckTypedPercentTest) TestGetTypedPercent6() {
	js := "a = {TRUE: 1, FALSE: 0};"
	AssertThat(c.GetTypedPercent(js)).IsWithin(0.1).
		Of(100.0)
}

// Java: TypeCheckTypedPercentTest.testResolvingNamedTypes()
func (c *TypeCheckTypedPercentTest) TestResolvingNamedTypes() {
	externs := NewTestExternsBuilder().AddObject().Build()
	js := `/** @constructor */
var Foo = function() {}
/** @param {number} a */
Foo.prototype.foo = function(a) {
  return this.baz().toString();
};
/** @return {Baz} */
Foo.prototype.baz = function() { return new Baz(); };
/** @constructor
  * @extends Foo */
var Bar = function() {};
/** @constructor */
var Baz = function() {};
`
	AssertThat(c.GetTypedPercentWithExterns(externs, js)).IsWithin(0.1).
		Of(100.0)
}

// Java: TypeCheckTypedPercentTest.testGetTypedPercent_constAndLet()
func (c *TypeCheckTypedPercentTest) TestGetTypedPercent_constAndLet() {
	js := `const id = function(x) { return x; }
let id2 = function(x) { return id(x); }
`
	AssertThat(c.GetTypedPercent(js)).IsWithin(0.1).
		Of(50.0)
}

// Java: TypeCheckTypedPercentTest.getTypedPercent(String)
func (c *TypeCheckTypedPercentTest) GetTypedPercent(js string) float64 {
	return c.GetTypedPercentWithExterns("", js)
}

// Java: TypeCheckTypedPercentTest.getTypedPercentWithExterns(String, String)
func (c *TypeCheckTypedPercentTest) GetTypedPercentWithExterns(externs string, js string) float64 {
	var compiler *Compiler = NewCompiler()
	var registry *rhino.JSTypeRegistry = compiler.GetTypeRegistry()
	var jsRoot *rhino.Node = rhino.IR_Root(compiler.ParseTestCode(js))
	var externsRoot *rhino.Node = rhino.IR_Root(compiler.ParseTestCode(externs))
	rhino.IR_Root(externsRoot, jsRoot)
	var t *TypeCheck = NewTypeCheck(compiler, NewSemanticReverseAbstractInterpreter(registry), registry)
	t.ProcessForTesting(externsRoot, jsRoot)
	AssertThat(compiler.GetErrors()).IsEmpty()
	return t.GetTypedPercent()
}

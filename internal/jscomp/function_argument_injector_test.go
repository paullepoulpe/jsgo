// Port of com/google/javascript/jscomp/FunctionArgumentInjectorTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/FunctionArgumentInjectorTest.java)

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// FunctionArgumentInjectorTest: tests for the static methods in
// FunctionArgumentInjector.
type FunctionArgumentInjectorTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	compiler                 *Compiler
	functionArgumentInjector *FunctionArgumentInjector
}

var FunctionArgumentInjectorTest_EMPTY_STRING_SET = jutil.NewLinkedHashSet[string]()

// Java: FunctionArgumentInjectorTest()
func NewFunctionArgumentInjectorTest() *FunctionArgumentInjectorTest {
	return &FunctionArgumentInjectorTest{}
}

func TestFunctionArgumentInjectorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewFunctionArgumentInjectorTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *FunctionArgumentInjectorTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore runs the @Before method setUp().
func (c *FunctionArgumentInjectorTest) JUnitBefore() { c.SetUp() }

// JUnitAfter: the class has no @After method.
func (c *FunctionArgumentInjectorTest) JUnitAfter() {}

// Java: FunctionArgumentInjectorTest.setUp()
func (c *FunctionArgumentInjectorTest) SetUp() {
	c.compiler = NewCompiler()
	options := NewCompilerOptions()

	c.compiler.InitOptions(options)
	c.functionArgumentInjector = NewFunctionArgumentInjector(c.compiler.GetAstAnalyzer())
}

// functionArgumentInjectorTest_mapOf is ImmutableMap.of(k, v).
//
// Go-only: test-helper — ImmutableMap is a jutil.LinkedHashMap in the port.
func functionArgumentInjectorTest_mapOf(k string, v rhino.Node) *jutil.LinkedHashMap[string, rhino.Node] {
	return jutil.NewImmutableMapBuilder[string, rhino.Node]().Put(k, v).BuildOrThrow()
}

// functionArgumentInjectorTest_mapOf2 is ImmutableMap.of(k1, v1, k2, v2).
//
// Go-only: test-helper — ImmutableMap is a jutil.LinkedHashMap in the port.
func functionArgumentInjectorTest_mapOf2(k1 string, v1 rhino.Node, k2 string, v2 rhino.Node) *jutil.LinkedHashMap[string, rhino.Node] {
	return jutil.NewImmutableMapBuilder[string, rhino.Node]().Put(k1, v1).Put(k2, v2).BuildOrThrow()
}

// Java: FunctionArgumentInjectorTest.testInject0()
func (c *FunctionArgumentInjectorTest) TestInject0() {
	result :=
		c.functionArgumentInjector.Inject(
			c.compiler,
			NodeUtil_GetFunctionBody(c.ParseFunction("function f(x) { alert(x); }")),
			rhino.NilNode,
			functionArgumentInjectorTest_mapOf("x", c.Parse("null").GetFirstFirstChild()))
	AssertNode(result).IsEqualTo(NodeUtil_GetFunctionBody(c.ParseFunction("function f(x) { alert(null); }")))
}

// Java: FunctionArgumentInjectorTest.testInject1_thisRef()
func (c *FunctionArgumentInjectorTest) TestInject1_thisRef() {
	result := c.functionArgumentInjector.Inject(
		c.compiler,
		NodeUtil_GetFunctionBody(c.ParseFunction("function f() { alert(this); }")),
		rhino.NilNode,
		functionArgumentInjectorTest_mapOf("this", c.Parse("null").GetFirstFirstChild()))
	AssertNode(result).IsEqualTo(NodeUtil_GetFunctionBody(c.ParseFunction("function f() { alert(null); }")))
}

// Java: FunctionArgumentInjectorTest.testInject_thisRef_classField_instance()
func (c *FunctionArgumentInjectorTest) TestInject_thisRef_classField_instance() {
	result := c.functionArgumentInjector.Inject(
		c.compiler,
		NodeUtil_GetFunctionBody(c.ParseFunction("function f() { return class { field = this; }; }")),
		rhino.NilNode,
		functionArgumentInjectorTest_mapOf("this", c.Parse("null").GetFirstFirstChild()))

	AssertNode(result).
		IsEqualTo(
			NodeUtil_GetFunctionBody(c.ParseFunction("function f() { return class { field = this; }; }")))
}

// Java: FunctionArgumentInjectorTest.testInject_thisRef_classField_static()
func (c *FunctionArgumentInjectorTest) TestInject_thisRef_classField_static() {
	result := c.functionArgumentInjector.Inject(
		c.compiler,
		NodeUtil_GetFunctionBody(
			c.ParseFunction("function f() { return class { static field = this; }; }")),
		rhino.NilNode,
		functionArgumentInjectorTest_mapOf("this", c.Parse("null").GetFirstFirstChild()))

	AssertNode(result).
		IsEqualTo(
			NodeUtil_GetFunctionBody(
				c.ParseFunction("function f() { return class { static field = this; }; }")))
}

// Java: FunctionArgumentInjectorTest.testInject_thisRef_classField_computed()
func (c *FunctionArgumentInjectorTest) TestInject_thisRef_classField_computed() {
	result := c.functionArgumentInjector.Inject(
		c.compiler,
		NodeUtil_GetFunctionBody(c.ParseFunction("function f() { return class { [this.k] = this.v; }; }")),
		rhino.NilNode,
		functionArgumentInjectorTest_mapOf("this", c.Parse("x").GetFirstFirstChild()))

	AssertNode(result).
		IsEqualTo(
			NodeUtil_GetFunctionBody(c.ParseFunction("function f() { return class { [x.k] = this.v; }; }")))
}

// Java: FunctionArgumentInjectorTest.testInject_thisRef_classField_staticComputed()
func (c *FunctionArgumentInjectorTest) TestInject_thisRef_classField_staticComputed() {
	result := c.functionArgumentInjector.Inject(
		c.compiler,
		NodeUtil_GetFunctionBody(
			c.ParseFunction("function f() { return class { static [this.k] = this.v; }; }")),
		rhino.NilNode,
		functionArgumentInjectorTest_mapOf("this", c.Parse("x").GetFirstFirstChild()))

	AssertNode(result).
		IsEqualTo(
			NodeUtil_GetFunctionBody(
				c.ParseFunction("function f() { return class { static [x.k] = this.v; }; }")))
}

// Java: FunctionArgumentInjectorTest.testInject_thisRef_classStaticBlock()
func (c *FunctionArgumentInjectorTest) TestInject_thisRef_classStaticBlock() {
	result := c.functionArgumentInjector.Inject(
		c.compiler,
		NodeUtil_GetFunctionBody(
			c.ParseFunction("function f() { return class { static { alert(this); } }; }")),
		rhino.NilNode,
		functionArgumentInjectorTest_mapOf("this", c.Parse("null").GetFirstFirstChild()))

	AssertNode(result).
		IsEqualTo(
			NodeUtil_GetFunctionBody(
				c.ParseFunction("function f() { return class { static { alert(this); } }; }")))
}

// Java: FunctionArgumentInjectorTest.testInject_thisRef_classField_paramSubstitution()
func (c *FunctionArgumentInjectorTest) TestInject_thisRef_classField_paramSubstitution() {
	result := c.functionArgumentInjector.Inject(
		c.compiler,
		NodeUtil_GetFunctionBody(
			c.ParseFunction("function f(val) { return class { field = val; [val] = this; }; }")),
		rhino.NilNode,
		functionArgumentInjectorTest_mapOf2(
			"val", c.Parse("123").GetFirstFirstChild(),
			"this", c.Parse("null").GetFirstFirstChild()))

	AssertNode(result).
		IsEqualTo(
			NodeUtil_GetFunctionBody(
				c.ParseFunction("function f(val) { return class { field = 123; [123] = this; }; }")))
}

// TODO(johnlenz): Add more unit tests for "inject"

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters0()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters0() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a){ return a; }"))).
		IsEmpty()
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters1()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters1() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a){ return a==0; }"))).
		IsEmpty()
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters2()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters2() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(c.ParseFunction("function f(a){ b=a }"))).
		IsEmpty()
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters3()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters3() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(c.ParseFunction("function f(a){ a=0 }"))).
		ContainsExactly("a")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters4()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters4() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a,b){ a=0;b=0 }"))).
		ContainsExactly("a", "b")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters5()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters5() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a,b){ a; if (a) b=0 }"))).
		ContainsExactly("b")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters6()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters6() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a,b){ function f(){ a;b; } }"))).
		ContainsExactly("a", "b")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters7()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters7() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a,b){ a; function f(){ b; } }"))).
		ContainsExactly("b")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters8()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters8() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a,b){ a; function f(){ function g() { b; } } }"))).
		ContainsExactly("b")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters9()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters9() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a,b){ (function(){ a;b; }) }"))).
		ContainsExactly("a", "b")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters10()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters10() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a,b){ a; (function (){ b; }) }"))).
		ContainsExactly("b")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters11()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters11() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a,b){ a; (function(){ (function () { b; }) }) }"))).
		ContainsExactly("b")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters12()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters12() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a){ { let a = 1; } }"))).
		IsEmpty()
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters13()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters13() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a){ { const a = 1; } }"))).
		IsEmpty()
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters14()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters14() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a){ for (a in []) {} }"))).
		ContainsExactly("a")
}

// Note: This is technically incorrect. The parameter a is shadowed, not modified. However, this
// will just cause the function inliner to do a little bit of unnecessary work; it will not
// result in incorrect output.
// Java: FunctionArgumentInjectorTest.testFindModifiedParameters15()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters15() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a){ for (const a in []) {} }"))).
		IsEmpty()
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters16()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters16() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a){ for (a of []) {} }"))).
		ContainsExactly("a")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters17()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters17() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a){ [a] = [2]; }"))).
		ContainsExactly("a")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters18()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters18() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(a){ var [a] = [2]; }"))).
		ContainsExactly("a")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters19()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters19() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(p){ const stash = (p = 42); }"))).
		ContainsExactly("p")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters20()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters20() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction("function f(p){ let stash = (p = 42); }"))).
		ContainsExactly("p")
}

// Java: FunctionArgumentInjectorTest.testFindModifiedParameters_classFields()
func (c *FunctionArgumentInjectorTest) TestFindModifiedParameters_classFields() {
	AssertThat(
		c.functionArgumentInjector.FindModifiedParameters(
			c.ParseFunction(
				`function f(a, b, c, d, e, g) {
  return class {
    x = a;
    [b] = c;
    static y = d;
    static [d] = e;
    uninitialized;
    [g];
  };
}
`))).
		ContainsExactly("a", "c")
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps1()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps1() {
	// Parameters with side-effects must be executed
	// even if they aren't referenced.
	c.TestNeededTemps("function foo(a,b){}; foo(goo(),goo());", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps1_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps1_optChain() {
	// Parameters with side-effects must be executed even if they aren't referenced.
	c.TestNeededTemps("function foo(a,b){}; foo?.(goo(),goo());", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps2()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps2() {
	// Unreferenced parameters without side-effects can be ignored.
	c.TestNeededTemps("function foo(a,b){}; foo(1,2);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForModifiedParams1()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForModifiedParams1() {
	// The param is not modified so no need for temps.
	c.TestNeededTemps("function foo(a,b,c){a;b;c;}; foo(x,x,x);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForModifiedParams2()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForModifiedParams2() {
	// The first param is modified, so we add a temp for it.
	c.TestNeededTemps("function foo(a,b,c){a;b;c;}; foo(x=2,x,x);", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForModifiedParams3()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForModifiedParams3() {
	// The second param is modified, with first param using the same var, so we add temps for both.
	c.TestNeededTemps("function foo(a,b,c){a;b;c;}; foo(x,x=2,x);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForModifiedParams4()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForModifiedParams4() {
	// The second param is modified, so we only add temp for both.
	c.TestNeededTemps("function foo(a,b,c){a;b;c;}; foo(y,x=2,x);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForModifiedParams5()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForModifiedParams5() {
	// The second param is modified with a more complex statement, we detect it and add temps for
	// first two params.
	c.TestNeededTemps(
		"function foo(a,b,c){a;b;c;}; foo(x,(y=2,x=1,z=6),x);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForModifiedParams6()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForModifiedParams6() {
	// The second param is modified with a function call, we detect it and add temps for first two
	// params.
	c.TestNeededTemps(
		"function foo(a,b,c){a;b;c;}; foo(x,(function(){y=2,x=1})(),x);",
		"foo",
		jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForModifiedParams7()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForModifiedParams7() {
	// The second param is modified with addition, so we add temps for both.
	c.TestNeededTemps(
		"function foo(a,b,c){a;b;c;}; foo(x,x+=2,x);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForModifiedParams8()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForModifiedParams8() {
	// The second param is modified with increment, with first param using the same var, so we add
	// temps for both.
	c.TestNeededTemps("function foo(a,b,c){a;b;c;}; foo(x,x++,x);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForModifiedParams9()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForModifiedParams9() {
	// The second param is used in addition but not modified, no temps needed.
	c.TestNeededTemps("function foo(a,b,c){a;b;c;}; foo(x,x+2,x);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForModifiedParams10()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForModifiedParams10() {
	// The second param is used in addition, the third param is modified, temp all refs.
	c.TestNeededTemps(
		"function foo(a,b,c){a;b;c;}; foo(x,x+2,x=5);", "foo", jutil.NewLinkedHashSet("a", "b", "c"))
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForCallParams1()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForCallParams1() {
	// The second param is a function call so temp it and previous args with names or calls
	c.TestNeededTemps(
		"function foo(a,b,c){a;b;c;}; var x; function incX(){return ++x}; foo(x,incX(),x);",
		"foo",
		jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testAddingTempsForCallParams2()
func (c *FunctionArgumentInjectorTest) TestAddingTempsForCallParams2() {
	// The second param is a function call so temp it.
	c.TestNeededTemps(
		"function foo(a,b,c){a;b;c;}; var x; function incX(){return ++x}; foo(4,incX(),x);",
		"foo",
		// TODO: b/309593967 The first param is a number, so we can skip creating a temp for it.
		jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps2_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps2_optChain() {
	// Unreferenced parameters without side-effects can be ignored.
	c.TestNeededTemps("function foo(a,b){}; foo?.(1,2);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps3()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps3() {
	// Referenced parameters without side-effects
	// don't need temps.
	c.TestNeededTemps("function foo(a,b){a;b;}; foo(x,y);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps3_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps3_optChain() {
	// Referenced parameters without side-effects don't need temps.
	c.TestNeededTemps("function foo(a,b){a;b;}; foo?.(x,y);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps4()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps4() {
	// Parameters referenced after side-effect must be assigned to temps.
	// Since b is getting hoisted as a temp, we must also hoist a.
	c.TestNeededTemps("function foo(a,b){a;goo();b;}; foo(x,y);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testBodyHasConditionalCode()
func (c *FunctionArgumentInjectorTest) TestBodyHasConditionalCode() {
	c.TestNeededTemps("function foo(a,b){a;goo?.();b;}; foo(x,y);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)

	c.TestNeededTemps("function foo(a,b){a; p&&q; b;}; foo(x,y);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)

	c.TestNeededTemps("function foo(a,b){a; p?q:r; b;}; foo(x,y);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps4_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps4_optChain() {
	// Parameters referenced after side-effect must be assigned to temps.
	// Any param before that param must also be assigned to a temp.
	c.TestNeededTemps("function foo(a,b){a;goo();b;}; foo?.(x,y);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps5()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps5() {
	// Parameters referenced after out-of-scope side-effect must
	// be assigned to temps.
	c.TestNeededTemps("function foo(a,b){x = b; y = a;}; foo(x,y);", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps5_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps5_optChain() {
	// Parameters referenced after out-of-scope side-effect must be assigned to temps.
	c.TestNeededTemps("function foo(a,b){x = b; y = a;}; foo?.(x,y);", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps6()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps6() {
	// Parameter referenced after a out-of-scope side-effect must
	// be assigned to a temp.
	c.TestNeededTemps("function foo(a){x++;a;}; foo(x);", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps6_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps6_optChain() {
	// Parameter referenced after a out-of-scope side-effect must be assigned to a temp.
	c.TestNeededTemps("function foo(a){x++;a;}; foo?.(x);", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps7()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps7() {
	// No temp needed after local side-effects.
	c.TestNeededTemps("function foo(a){var c; c = 0; a;}; foo(x);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)

	c.TestNeededTemps("function foo(a){let c; c = 0; a;}; foo(x);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)

	c.TestNeededTemps("function foo(a){const c = 0; a;}; foo(x);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps8()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps8() {
	// Temp needed for side-effects to object using local name.
	c.TestNeededTemps(
		"function foo(a){var c = {}; c.goo=0; a;}; foo(x);", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps8_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps8_optChain() {
	// Temp needed for side-effects to object using local name.
	c.TestNeededTemps(
		"function foo(a){var c = {}; c.goo=0; a;}; foo?.(x);", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps9()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps9() {
	// Parameters referenced in a loop with side-effects must
	// be assigned to temps.
	c.TestNeededTemps(
		"function foo(a,b){while(true){a;goo();b;}}; foo(x,y);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps9_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps9_optChain() {
	// Parameters referenced in a loop with side-effects must be assigned to temps.
	c.TestNeededTemps(
		"function foo(a,b){while(true){a;goo();b;}}; foo?.(x,y);",
		"foo",
		jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps10()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps10() {
	// No temps for parameters referenced in a loop with no side-effects.
	c.TestNeededTemps(
		"function foo(a,b){while(true){a;true;b;}}; foo(x,y);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps11()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps11() {
	// Parameters referenced in a loop with side-effects must
	// be assigned to temps.
	c.TestNeededTemps(
		"function foo(a,b){do{a;b;}while(goo());}; foo(x,y);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps11_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps11_optChain() {
	// Parameters referenced in a loop with side-effects must be assigned to temps.
	c.TestNeededTemps(
		"function foo(a,b){do{a;b;}while(goo());}; foo?.(x,y);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps12()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps12() {
	// Parameters referenced in a loop with side-effects must
	// be assigned to temps.
	c.TestNeededTemps(
		"function foo(a,b){for(;;){a;b;goo();}}; foo(x,y);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps13()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps13() {
	// Parameters referenced in a inner loop without side-effects must
	// be assigned to temps if the outer loop has side-effects.
	c.TestNeededTemps(
		"function foo(a,b){for(;;){for(;;){a;b;}goo();}}; foo(x,y);",
		"foo",
		jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps13_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps13_optChain() {
	// Parameters referenced in a inner loop without side-effects must be assigned to temps if the
	// outer loop has side-effects.
	c.TestNeededTemps(
		"function foo(a,b){for(;;){for(;;){a;b;}goo();}}; foo?.(x,y);",
		"foo",
		jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps14()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps14() {
	// Parameters referenced in a loop must
	// be assigned to temps.
	c.TestNeededTemps(
		"function foo(a,b){goo();for(;;){a;b;}}; foo(x,y);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps20()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps20() {
	// A long string referenced more than once should have a temp.
	c.TestNeededTemps("function foo(a){a;a;}; foo(\"blah blah\");", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps21()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps21() {
	// A short string referenced once should not have a temp.
	c.TestNeededTemps("function foo(a){a;a;}; foo(\"\");", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps22()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps22() {
	// A object literal not referenced.
	c.TestNeededTemps("function foo(a){}; foo({x:1});", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
	// A object literal referenced after side-effect, should have a temp.
	c.TestNeededTemps("function foo(a){alert('foo');a;}; foo({x:1});", "foo", jutil.NewLinkedHashSet("a"))
	// A object literal referenced after side-effect, should have a temp.
	c.TestNeededTemps(
		"function foo(a,b){b;a;}; foo({x:1},alert('foo'));", "foo", jutil.NewLinkedHashSet("a", "b"))
	// A object literal, referenced more than once, should have a temp.
	c.TestNeededTemps("function foo(a){a;a;}; foo({x:1});", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps22b()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps22b() {
	// A object literal not referenced.
	c.TestNeededTemps(
		"function foo(a){a(this)}; foo.call(f(),g());", "foo", jutil.NewLinkedHashSet("a", "this"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps22b_optChain()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps22b_optChain() {
	// A object literal not referenced.
	c.TestNeededTemps(
		"function foo(a){a(this)}; foo?.call(f(),g());", "foo", jutil.NewLinkedHashSet("a", "this"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps23()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps23() {
	// A array literal, not referenced.
	c.TestNeededTemps("function foo(a){}; foo([1,2]);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
	// A array literal, referenced once after side-effect, should have a temp.
	c.TestNeededTemps("function foo(a){alert('foo');a;}; foo([1,2]);", "foo", jutil.NewLinkedHashSet("a"))
	// A array literal, referenced more than once, should have a temp.
	c.TestNeededTemps("function foo(a){a;a;}; foo([1,2]);", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps24()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps24() {
	// A regex literal, not referenced.
	c.TestNeededTemps("function foo(a){}; foo(/mac/);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
	// A regex literal, referenced once after side-effect, should have a temp.
	c.TestNeededTemps("function foo(a){alert('foo');a;}; foo(/mac/);", "foo", jutil.NewLinkedHashSet("a"))
	// A regex literal, referenced more than once, should have a temp.
	c.TestNeededTemps("function foo(a){a;a;}; foo(/mac/);", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps25()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps25() {
	// A side-effect-less constructor, not referenced.
	c.TestNeededTemps("function foo(a){}; foo(new Date());", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
	// A side-effect-less constructor, referenced once after sideeffect, should have a temp.
	c.TestNeededTemps(
		"function foo(a){alert('foo');a;}; foo(new Date());", "foo", jutil.NewLinkedHashSet("a"))
	// A side-effect-less constructor, referenced more than once, should have
	// a temp.
	c.TestNeededTemps("function foo(a){a;a;}; foo(new Date());", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps26()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps26() {
	// A constructor, not referenced.
	c.TestNeededTemps("function foo(a){}; foo(new Bar());", "foo", jutil.NewLinkedHashSet("a"))
	// A constructor, referenced once after a sideeffect, should have a temp.
	c.TestNeededTemps(
		"function foo(a){alert('foo');a;}; foo(new Bar());", "foo", jutil.NewLinkedHashSet("a"))
	// A constructor, referenced more than once, should have a temp.
	c.TestNeededTemps("function foo(a){a;a;}; foo(new Bar());", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps27()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps27() {
	// Ensure the correct parameter is given a temp, when there is
	// a this value in the call. Since `goo()` is a call and needs to be evaluated in a temp before
	// inlining, we also evaluate the previous args in temporaries.
	c.TestNeededTemps(
		"function foo(a,b,c){}; foo.call(this,1,goo(),2);",
		"foo",
		jutil.NewLinkedHashSet("this", "a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps28()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps28() {
	// true/false are don't need temps
	c.TestNeededTemps("function foo(a){a;a;}; foo(true);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps29()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps29() {
	// true/false are don't need temps
	c.TestNeededTemps("function foo(a){a;a;}; foo(false);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps30()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps30() {
	// true/false are don't need temps
	c.TestNeededTemps("function foo(a){a;a;}; foo(!0);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps31()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps31() {
	// true/false are don't need temps
	c.TestNeededTemps("function foo(a){a;a;}; foo(!1);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps32()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps32() {
	// void 0 doesn't need a temp
	c.TestNeededTemps("function foo(a){a;a;}; foo(void 0);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps33()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps33() {
	// doesn't need a temp
	c.TestNeededTemps("function foo(a){return a;}; foo(new X);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps_argMutatesState()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps_argMutatesState() {
	// `y=z` changes state. We create a temp for it and all previous args that it could affect.
	c.TestNeededTemps(
		"function foo(a, b, c){return a+b+c;}; foo(x,y+1,y=z);",
		"foo",
		jutil.NewLinkedHashSet("a", "b", "c"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTemps_argMutatesState2()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTemps_argMutatesState2() {
	// `y=z` happens before `y+1` is read. Safe to inject a temp for evaluating it.
	// Consequently, we also inject a temp for previous arg `a` that it could affect
	c.TestNeededTemps(
		"function foo(a, b, c){return a+b+c;}; foo(x,y=z,y+1);", "foo", jutil.NewLinkedHashSet("a", "b"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTempsInLoops()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTempsInLoops() {
	// A mutable parameter referenced in loop needs a
	// temporary.
	c.TestNeededTemps("function foo(a){for(;;)a;}; foo(new Bar());", "foo", jutil.NewLinkedHashSet("a"))

	c.TestNeededTemps("function foo(a){while(true)a;}; foo(new Bar());", "foo", jutil.NewLinkedHashSet("a"))

	c.TestNeededTemps(
		"function foo(a){do{a;}while(true)}; foo(new Bar());", "foo", jutil.NewLinkedHashSet("a"))
}

// Java: FunctionArgumentInjectorTest.nullishCoalesce()
func (c *FunctionArgumentInjectorTest) TestNullishCoalesce() {
	c.TestNeededTemps(
		"function foo(...args) {return args ?? x;} foo(1, 2);", "foo", jutil.NewLinkedHashSet("args"))
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTempsRestParam1()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTempsRestParam1() {
	c.TestNeededTemps("function foo(...args) {return args;} foo(1, 2);", "foo", FunctionArgumentInjectorTest_EMPTY_STRING_SET)
}

// Java: FunctionArgumentInjectorTest.testGatherCallArgumentsNeedingTempsRestParam2()
func (c *FunctionArgumentInjectorTest) TestGatherCallArgumentsNeedingTempsRestParam2() {
	// Since args is hoisted as a temp, we must also hoist the previous argument a
	c.TestNeededTemps(
		"function foo(x, ...args) {return args;} foo(1, 2);", "foo", jutil.NewLinkedHashSet("x", "args"))
}

// Java: FunctionArgumentInjectorTest.testArgMapWithRestParam1()
func (c *FunctionArgumentInjectorTest) TestArgMapWithRestParam1() {
	c.AssertArgMapHasKeys(
		"function foo(...args){return args;} foo(1, 2);", "foo", jutil.NewLinkedHashSet("this", "args"))
}

// Java: FunctionArgumentInjectorTest.testArgMapWithRestParam2()
func (c *FunctionArgumentInjectorTest) TestArgMapWithRestParam2() {
	c.AssertArgMapHasKeys(
		"function foo(...args){return args;} foo();", "foo", jutil.NewLinkedHashSet("this", "args"))
}

// Java: FunctionArgumentInjectorTest.assertArgMapHasKeys(String, String, ImmutableSet<String>)
func (c *FunctionArgumentInjectorTest) AssertArgMapHasKeys(code string, fnName string, expectedKeys *jutil.LinkedHashSet[string]) {
	n := c.Parse(code)
	fn := FunctionArgumentInjectorTest_FindFunction(n, fnName)
	AssertThat(fn).IsNotNull()
	call := FunctionArgumentInjectorTest_FindCall(n, fnName)
	AssertThat(call).IsNotNull()
	actualMap := c.GetAndValidateFunctionCallParameterMap(fn, call)
	AssertThat(jutil.NewLinkedHashSet(actualMap.KeySlice()...)).IsEqualTo(expectedKeys) // keySet()
}

// Java: FunctionArgumentInjectorTest.testNeededTemps(String, String, ImmutableSet<String>)
func (c *FunctionArgumentInjectorTest) TestNeededTemps(code string, fnName string, expectedTemps *jutil.LinkedHashSet[string]) {
	n := c.Parse(code)
	fn := FunctionArgumentInjectorTest_FindFunction(n, fnName)
	AssertThat(fn).IsNotNull()
	call := FunctionArgumentInjectorTest_FindCall(n, fnName)
	AssertThat(call).IsNotNull()
	args := c.GetAndValidateFunctionCallParameterMap(fn, call)

	actualTemps :=
		c.functionArgumentInjector.GatherCallArgumentsNeedingTemps(
			c.compiler, fn, args, jutil.NewLinkedHashSet[string](), NewClosureCodingConvention())

	AssertThat(actualTemps).IsEqualTo(expectedTemps)
}

// Java: FunctionArgumentInjectorTest.getAndValidateFunctionCallParameterMap(Node, Node)
func (c *FunctionArgumentInjectorTest) GetAndValidateFunctionCallParameterMap(
	fn rhino.Node, call rhino.Node) *jutil.LinkedHashMap[string, *FunctionArgumentInjector_ParamArgPair] {
	map_ := c.functionArgumentInjector.GetFunctionCallParameterMap(fn, call, FunctionArgumentInjectorTest_GetNameSupplier())
	// Verify that all nodes in the map have source info, so they are valid to add to the AST
	for name, value := range map_.All() {
		node := value.Arg()

		NewSourceInfoCheck(c.compiler).SetCheckSubTree(node)
		AssertWithMessage("errors for name: %s", name).That(c.compiler.GetErrors()).IsEmpty()
	}

	return map_
}

// Java: FunctionArgumentInjectorTest.getNameSupplier()
func FunctionArgumentInjectorTest_GetNameSupplier() func() string {
	i := 0
	return func() string {
		v := jutil.ToJavaString(i) // String.valueOf(i++)
		i++
		return v
	}
}

// Java: FunctionArgumentInjectorTest.findCall(Node, String)
func FunctionArgumentInjectorTest_FindCall(n rhino.Node, name string) rhino.Node {
	if NodeUtil_IsNormalOrOptChainCall(n) {
		var callee rhino.Node
		if NodeUtil_IsNormalOrOptChainGetProp(n.GetFirstChild()) {
			callee = n.GetFirstFirstChild()
			// Only "call" is supported at this point.
			jutil.CheckArgument(callee.GetParent().GetString() == "call")
		} else {
			callee = n.GetFirstChild()
		}

		if callee.IsName() && callee.GetString() == name {
			return n
		}
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		result := FunctionArgumentInjectorTest_FindCall(c, name)
		if !result.IsNil() {
			return result
		}
	}

	return rhino.NilNode
}

// Java: FunctionArgumentInjectorTest.findFunction(Node, String)
func FunctionArgumentInjectorTest_FindFunction(n rhino.Node, name string) rhino.Node {
	if n.IsFunction() {
		if n.GetFirstChild().GetString() == name {
			return n
		}
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		result := FunctionArgumentInjectorTest_FindFunction(c, name)
		if !result.IsNil() {
			return result
		}
	}

	return rhino.NilNode
}

// Java: FunctionArgumentInjectorTest.parseFunction(String)
func (c *FunctionArgumentInjectorTest) ParseFunction(js string) rhino.Node {
	return c.Parse(js).GetFirstChild()
}

// Java: FunctionArgumentInjectorTest.parse(String)
func (c *FunctionArgumentInjectorTest) Parse(js string) rhino.Node {
	n := c.compiler.ParseTestCode(js)
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	return n
}

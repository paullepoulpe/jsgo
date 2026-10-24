// Port of com/google/javascript/jscomp/OptimizeCallsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/OptimizeCallsTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"iter"
	"testing"
)

type OptimizeCallsTest struct {
	CompilerTestCaseBase

	// Will be assigned with the references collected from the most recent pass.
	references *OptimizeCalls_ReferenceMap

	// Whether to consider externs during the next collection. Must be explicitly set.
	considerExterns *bool
}

// Java: OptimizeCallsTest()
func NewOptimizeCallsTest() *OptimizeCallsTest {
	c := &OptimizeCallsTest{}
	c.InitCompilerTestCaseBase(c)
	c.considerExterns = nil
	return c
}

func TestOptimizeCallsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewOptimizeCallsTest() })
}

// Java: OptimizeCallsTest.setUp()
func (c *OptimizeCallsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize() // Required for `OptimizeCalls`.
	// Even if we're ignoring externs we need to know their names to do that.
	c.EnableGatherExternProperties()
}

// Java: OptimizeCallsTest.getProcessor(Compiler)
func (c *OptimizeCallsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return OptimizeCalls_Builder_().
		SetCompiler(compiler).
		SetConsiderExterns(*c.considerExterns). // Boolean unboxing: null throws.
		AddPass(OptimizeCalls_CallGraphCompilerPassFunc(func(externs rhino.Node, root rhino.Node, references *OptimizeCalls_ReferenceMap) {
			c.references = references
		})).
		Build()
}

// setConsiderExterns is `considerExterns = b` (a
// Boolean field).
//
// Go-only: test-helper — a *bool needs an addressable value.
func (c *OptimizeCallsTest) setConsiderExterns(b bool) { c.considerExterns = &b }

// optimizeCallsTest_Entry is a Map.Entry<String, ArrayList<Node>> of a
// ReferenceMap's references.
//
// Go-only: test-helper — Java iterates the entry set directly.
type optimizeCallsTest_Entry struct {
	key   string
	value *jutil.ArrayList[rhino.Node]
}

// optimizeCallsTest_entries collects ReferenceMap.getNameReferences() /
// getPropReferences() (an Iterable of entries) for assertThat.
func optimizeCallsTest_entries(refs iter.Seq2[string, *jutil.ArrayList[rhino.Node]]) []optimizeCallsTest_Entry {
	var out []optimizeCallsTest_Entry
	for k, v := range refs {
		out = append(out, optimizeCallsTest_Entry{key: k, value: v})
	}
	return out
}

// optimizeCallsTest_copyOf is ImmutableMap.copyOf(Iterable<Entry>).
func optimizeCallsTest_copyOf(refs iter.Seq2[string, *jutil.ArrayList[rhino.Node]]) *jutil.LinkedHashMap[string, *jutil.ArrayList[rhino.Node]] {
	m := jutil.NewImmutableMapBuilder[string, *jutil.ArrayList[rhino.Node]]()
	for k, v := range refs {
		m.Put(k, v)
	}
	return m.BuildOrThrow()
}

var OptimizeCallsTest_KEY_EQUALITY = Correspondence_Transforming(
	func(e optimizeCallsTest_Entry) string { return e.key }, "has key")

var OptimizeCallsTest_HAS_TOKEN = Correspondence_Transforming((rhino.Node).GetToken, "has token")

// Java: OptimizeCallsTest.testReferenceCollection_ignoringExterns()
func (c *OptimizeCallsTest) TestReferenceCollection_ignoringExterns() {
	c.setConsiderExterns(false)

	c.Test(
		c.Externs(
			"function externFoo(dummyParam1, dummyParam2) {\n"+
				"  deadRef1;\n"+
				"}\n"+
				"\n"+
				"var externSpace = {};\n"+
				"\n"+
				"externSpace.externProp = {\n"+
				"  externChildProp: 0,\n"+
				"\n"+
				"  externChildMethod(dummyParam3) { }\n"+
				"};\n"+
				"\n"+
				"class ExternBar { }\n"),
		c.Srcs(
			"function foo(param1, param2) {\n"+
				"  var liveRef1;\n"+
				"  new ExternBar();\n"+
				"}\n"+
				"\n"+
				"var space = {};\n"+
				"\n"+
				"space.prop = {\n"+
				"  childProp: foo,\n"+
				"\n"+
				"  childMethod(param3, externFoo) { }\n"+
				"};\n"+
				"\n"+
				"class Bar { }\n"+
				"\n"+
				"new Bar()\n"+
				"space.prop.childMethod(externSpace);\n"))

	AssertThat(optimizeCallsTest_entries(c.references.GetNameReferences())).
		ComparingElementsUsing(OptimizeCallsTest_KEY_EQUALITY).
		// Only global names should be collected.
		ContainsExactly("Bar", "foo", "space")

	AssertThat(optimizeCallsTest_entries(c.references.GetPropReferences())).
		ComparingElementsUsing(OptimizeCallsTest_KEY_EQUALITY).
		ContainsExactly("childMethod", "childProp", "prop")
}

// Java: OptimizeCallsTest.testReferenceCollection_consideringExterns()
func (c *OptimizeCallsTest) TestReferenceCollection_consideringExterns() {
	c.setConsiderExterns(true)

	c.Test(
		c.Externs(
			"function externFoo(dummyParam1, dummyParam2) {\n"+
				"  deadRef1;\n"+
				"}\n"+
				"\n"+
				"var externSpace = {};\n"+
				"\n"+
				"externSpace.externProp = {\n"+
				"  externChildProp: 0,\n"+
				"\n"+
				"  externChildMethod(dummyParam3) { }\n"+
				"};\n"+
				"\n"+
				"class ExternBar { }\n"),
		c.Srcs(
			"function foo(param1, param2) {\n"+
				"  liveRef1;\n"+
				"  new ExternBar();\n"+
				"}\n"+
				"\n"+
				"var space = {};\n"+
				"\n"+
				"space.prop = {\n"+
				"  childProp: foo,\n"+
				"\n"+
				"  childMethod(param3, externFoo) { }\n"+
				"};\n"+
				"\n"+
				"class Bar { }\n"+
				"\n"+
				"new Bar()\n"+
				"space.prop.childMethod(externSpace);\n"))

	AssertThat(optimizeCallsTest_entries(c.references.GetNameReferences())).
		ComparingElementsUsing(OptimizeCallsTest_KEY_EQUALITY).
		// Only global names should be collected.
		ContainsExactly("Bar", "ExternBar", "externFoo", "externSpace", "foo", "space")

	AssertThat(optimizeCallsTest_entries(c.references.GetPropReferences())).
		ComparingElementsUsing(OptimizeCallsTest_KEY_EQUALITY).
		ContainsExactly(
			"childMethod",
			"childProp",
			"externChildMethod",
			"externChildProp",
			"externProp",
			"prop")
}

// Java: OptimizeCallsTest.testReferenceCollection_doesNotCollectTheEmptyName()
func (c *OptimizeCallsTest) TestReferenceCollection_doesNotCollectTheEmptyName() {
	c.setConsiderExterns(true)

	c.Test(c.Srcs("var f = function() { }"))

	AssertThat(optimizeCallsTest_entries(c.references.GetNameReferences())).
		ComparingElementsUsing(OptimizeCallsTest_KEY_EQUALITY).
		// Anonymous functions have the empty string as a name.
		DoesNotContain("")
}

// Java: OptimizeCallsTest.testReferenceCollection_findsSuperCalls_simpleNames()
func (c *OptimizeCallsTest) TestReferenceCollection_findsSuperCalls_simpleNames() {
	c.setConsiderExterns(true)

	c.Test(
		c.Srcs(
			"class SuperClass {\n" +
				"  constructor() {}\n" +
				"}\n" +
				"class SubClass extends SuperClass {\n" +
				"  constructor() {\n" +
				"    super();\n" +
				"  }\n" +
				"}\n" +
				"new SuperClass();\n" +
				"new SubClass();\n"))

	nameToRefs :=
		optimizeCallsTest_copyOf(c.references.GetNameReferences())
	AssertThat(nameToRefs.KeySlice()).ContainsExactly("SuperClass", "SubClass")

	AssertThat(nameToRefs.Get("SuperClass").Slice()).
		ComparingElementsUsing(OptimizeCallsTest_HAS_TOKEN).
		ContainsExactly(rhino.Token_NAME, rhino.Token_NAME, rhino.Token_SUPER, rhino.Token_NAME).
		InOrder()

	AssertThat(nameToRefs.Get("SubClass").Slice()).
		ComparingElementsUsing(OptimizeCallsTest_HAS_TOKEN).
		ContainsExactly(rhino.Token_NAME, rhino.Token_NAME).
		InOrder()
}

// Java: OptimizeCallsTest.testReferenceCollection_findsSuperCalls_qualifiedNames()
func (c *OptimizeCallsTest) TestReferenceCollection_findsSuperCalls_qualifiedNames() {
	c.setConsiderExterns(true)

	c.Test(
		c.Srcs(
			"const ns = {};\n" +
				"ns.SuperClass = class {\n" +
				"  constructor() {}\n" +
				"}\n" +
				"ns.SubClass = class extends ns.SuperClass {\n" +
				"  constructor() {\n" +
				"    super();\n" +
				"  }\n" +
				"}\n" +
				"new ns.SuperClass();\n" +
				"new ns.SubClass();\n"))

	nameToRefs :=
		optimizeCallsTest_copyOf(c.references.GetPropReferences())
	AssertThat(nameToRefs.KeySlice()).ContainsExactly("SuperClass", "SubClass", "constructor")

	AssertThat(nameToRefs.Get("SuperClass").Slice()).
		ComparingElementsUsing(OptimizeCallsTest_HAS_TOKEN).
		ContainsExactly(rhino.Token_GETPROP, rhino.Token_GETPROP, rhino.Token_SUPER, rhino.Token_GETPROP).
		InOrder()

	AssertThat(nameToRefs.Get("SubClass").Slice()).
		ComparingElementsUsing(OptimizeCallsTest_HAS_TOKEN).
		ContainsExactly(rhino.Token_GETPROP, rhino.Token_GETPROP).
		InOrder()
}

// Java: OptimizeCallsTest.testReferenceCollection_findsClassFields()
func (c *OptimizeCallsTest) TestReferenceCollection_findsClassFields() {
	c.setConsiderExterns(true)

	c.Test(
		c.Srcs(
			"class C {\n" +
				"  a = 2;\n" +
				"  b;\n" +
				"  ['c'] = 'hi';\n" +
				"  'd' = 5;\n" +
				"  1 = 2;\n" +
				"  e = function() { return 1; };\n" +
				"  ['f'] = function z() { return 2; };\n" +
				"}\n" +
				"new C();\n"))

	AssertThat(optimizeCallsTest_entries(c.references.GetNameReferences())).
		ComparingElementsUsing(OptimizeCallsTest_KEY_EQUALITY).
		ContainsExactly("C")

	AssertThat(optimizeCallsTest_entries(c.references.GetPropReferences())).
		ComparingElementsUsing(OptimizeCallsTest_KEY_EQUALITY).
		ContainsExactly("a", "b", "e")

	nameToRefs :=
		optimizeCallsTest_copyOf(c.references.GetPropReferences())
	AssertThat(nameToRefs.KeySlice()).ContainsExactly("a", "b", "e")

	AssertThat(nameToRefs.Get("a").Slice()).
		ComparingElementsUsing(OptimizeCallsTest_HAS_TOKEN).
		ContainsExactly(rhino.Token_MEMBER_FIELD_DEF).
		InOrder()

	AssertThat(nameToRefs.Get("b").Slice()).
		ComparingElementsUsing(OptimizeCallsTest_HAS_TOKEN).
		ContainsExactly(rhino.Token_MEMBER_FIELD_DEF).
		InOrder()

	AssertThat(nameToRefs.Get("e").Slice()).
		ComparingElementsUsing(OptimizeCallsTest_HAS_TOKEN).
		ContainsExactly(rhino.Token_MEMBER_FIELD_DEF).
		InOrder()
}

// Java: OptimizeCallsTest.testReferenceCollection_reflectedProps()
func (c *OptimizeCallsTest) TestReferenceCollection_reflectedProps() {
	c.setConsiderExterns(false)

	c.Test(
		c.Srcs(
			"class C {\n" +
				"  m() {\n" +
				"    console.log('hello');\n" +
				"  }\n" +
				"}\n" +
				"C.prototype[goog.reflect.objectProperty('m', C.prototype)]();\n"))

	nameToRefs :=
		optimizeCallsTest_copyOf(c.references.GetPropReferences())
	AssertThat(nameToRefs.Get("m").Slice()).
		ComparingElementsUsing(OptimizeCallsTest_HAS_TOKEN).
		ContainsExactly(rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_CALL).
		InOrder()
}

// Java: OptimizeCallsTest.testIsAllowedReference_taggedTemplateLiteral()
func (c *OptimizeCallsTest) TestIsAllowedReference_taggedTemplateLiteral() {
	c.setConsiderExterns(false)

	c.Test(
		c.Srcs(
			"var obj = {\n" +
				"  tag() {}\n" +
				"};\n" +
				"obj.tag`template`;\n"))

	nameToRefs :=
		optimizeCallsTest_copyOf(c.references.GetNameReferences())
	objRefs := nameToRefs.Get("obj")
	AssertThat(objRefs).IsNotNull()
	hasDisallowedRef := false
	for ref := range objRefs.All() {
		if !OptimizeCalls_IsAllowedReference(ref) {
			hasDisallowedRef = true
		}
	}
	AssertThat(hasDisallowedRef).IsTrue()
}

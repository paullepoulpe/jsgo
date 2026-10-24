// Port of com/google/javascript/jscomp/serialization/SerializeTypesToPointersTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/serialization/SerializeTypesToPointersTest.java)

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"strconv"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// The Java @Test methods are not named test*; their Go names get a "Test"
// prefix so that RunJUnit4 runs them.
//
// Java: public final class SerializeTypesToPointersTest
type SerializeTypesToPointersTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	compiler          *Compiler
	stringPoolBuilder *StringPool_Builder
}

// Java: SerializeTypesToPointersTest()
func NewSerializeTypesToPointersTest() *SerializeTypesToPointersTest {
	return &SerializeTypesToPointersTest{}
}

func TestSerializeTypesToPointersTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewSerializeTypesToPointersTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *SerializeTypesToPointersTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore runs the @Before method.
func (c *SerializeTypesToPointersTest) JUnitBefore() { c.SetUp() }

// JUnitAfter: the class has no @After method.
func (c *SerializeTypesToPointersTest) JUnitAfter() {}

// Java: SerializeTypesToPointersTest.setUp()
func (c *SerializeTypesToPointersTest) SetUp() {
	options := NewCompilerOptions()
	options.SetCheckTypes(true)
	c.compiler = NewCompiler()
	c.compiler.InitOptions(options)
	c.stringPoolBuilder = StringPool_Builder_()
}

// Java: SerializeTypesToPointersTest.outputsTypePointerForClass()
func (c *SerializeTypesToPointersTest) TestOutputsTypePointerForClass() {
	src := c.parseAndTypecheckFiles("class Foo {}")
	var fooCtorType rhino.JSType = c.getGlobalType("Foo").ToObjectType().GetConstructor()

	serializer :=
		SerializeTypesToPointers_Create(
			c.compiler,
			c.stringPoolBuilder,
			SerializationOptions_Builder_().SetIncludeDebugInfo(true).Build())
	serializer.GatherTypesOnAst(src)

	AssertThat(serializeTypesToPointersTest_get(serializer.GetTypePointersByJstype(), fooCtorType)).IsNotNull()

	fooCtorPointer := *serializeTypesToPointersTest_get(serializer.GetTypePointersByJstype(), fooCtorType)
	AssertThat(
		serializer.
			GetTypePool().
			GetType(TypePointers_TrimOffset(fooCtorPointer)).
			GetObject().
			GetMarkedConstructor()).
		IsTrue()
}

// Java: SerializeTypesToPointersTest.serializesPropertiesReferencedInSources()
func (c *SerializeTypesToPointersTest) TestSerializesPropertiesReferencedInSources() {
	src := c.parseAndTypecheckFiles("class Foo { serializeMe() {} } Foo.prototype;")
	var fooPrototypeType rhino.JSType = c.getGlobalType("Foo").ToObjectType().GetImplicitPrototype()

	serializer :=
		SerializeTypesToPointers_Create(
			c.compiler,
			c.stringPoolBuilder,
			SerializationOptions_Builder_().SetIncludeDebugInfo(true).Build())
	serializer.GatherTypesOnAst(src)

	AssertThat(serializeTypesToPointersTest_get(serializer.GetTypePointersByJstype(), fooPrototypeType)).IsNotNull()
	fooPrototypePointer := *serializeTypesToPointersTest_get(serializer.GetTypePointersByJstype(), fooPrototypeType)
	AssertThat(
		serializer.
			GetTypePool().
			GetType(TypePointers_TrimOffset(fooPrototypePointer)).
			GetObject().
			GetOwnPropertyList()).
		IsEqualTo(c.findAllInStringPool("serializeMe"))
}

// Java: SerializeTypesToPointersTest.doesNotSerializePropertiesOnlyReferencedInTypeSummary()
func (c *SerializeTypesToPointersTest) TestDoesNotSerializePropertiesOnlyReferencedInTypeSummary() {
	root :=
		c.parseAndTypecheckFiles(
			"/** @typeSummary */ class Foo { serializeMe() {} andMe() {} doNotSerializeMe() {} }",
			"new Foo().serializeMe(); const otherObj = {andMe: 0}; Foo.prototype;")
	var fooPrototypeType rhino.JSType = c.getGlobalType("Foo").ToObjectType().GetImplicitPrototype()

	serializer :=
		SerializeTypesToPointers_Create(
			c.compiler,
			c.stringPoolBuilder,
			SerializationOptions_Builder_().SetIncludeDebugInfo(true).Build())
	serializer.GatherTypesOnAst(root)

	AssertThat(serializeTypesToPointersTest_get(serializer.GetTypePointersByJstype(), fooPrototypeType)).IsNotNull()
	fooPrototypePointer := *serializeTypesToPointersTest_get(serializer.GetTypePointersByJstype(), fooPrototypeType)
	AssertThat(
		serializer.
			GetTypePool().
			GetType(TypePointers_TrimOffset(fooPrototypePointer)).
			GetObject().
			GetOwnPropertyList()).
		IsEqualTo(c.findAllInStringPool("andMe", "serializeMe"))
}

// Java: SerializeTypesToPointersTest.serializePropertiesInNonTypeSummaryExterns()
func (c *SerializeTypesToPointersTest) TestSerializePropertiesInNonTypeSummaryExterns() {
	root :=
		c.parseAndTypecheckFiles(
			"/** @typeSummary @externs */ class Foo { serializeMe() {} doNotSerializeMe() {} }",
			"/** @externs */\n"+
				"class Bar { serializeMe() {} }\n"+
				"/** @type {string} */\n"+
				"Foo.prototype.andMe;\n")
	var fooPrototypeType rhino.JSType = c.getGlobalType("Foo").ToObjectType().GetImplicitPrototype()

	serializer :=
		SerializeTypesToPointers_Create(
			c.compiler,
			c.stringPoolBuilder,
			SerializationOptions_Builder_().SetIncludeDebugInfo(true).Build())
	serializer.GatherTypesOnAst(root)

	AssertThat(serializeTypesToPointersTest_get(serializer.GetTypePointersByJstype(), fooPrototypeType)).IsNotNull()
	fooPrototypePointer := *serializeTypesToPointersTest_get(serializer.GetTypePointersByJstype(), fooPrototypeType)
	AssertThat(
		serializer.
			GetTypePool().
			GetType(TypePointers_TrimOffset(fooPrototypePointer)).
			GetObject().
			GetOwnPropertyList()).
		IsEqualTo(c.findAllInStringPool("andMe", "serializeMe"))
}

// Java: SerializeTypesToPointersTest.getGlobalType(String)
func (c *SerializeTypesToPointersTest) getGlobalType(typeName string) rhino.JSType {
	return c.compiler.GetTypeRegistry().GetGlobalType(typeName)
}

// Java: SerializeTypesToPointersTest.parseAndTypecheckFiles(String...)
func (c *SerializeTypesToPointersTest) parseAndTypecheckFiles(files ...string) rhino.Node {
	root := rhino.IR_Root()
	index := 0
	for _, file := range files {
		root.AddChildToBack(c.compiler.ParseSyntheticCode("test_"+strconv.Itoa(index), file))
		index++
	}
	AssertThat(c.compiler.GetErrors()).IsEmpty()
	rhino.IR_Root( /* externs */ rhino.IR_Root() /* js */, root) // make this a valid AST
	NewTypeCheck(
		c.compiler,
		NewClosureReverseAbstractInterpreter(c.compiler.GetTypeRegistry()),
		c.compiler.GetTypeRegistry()).
		ProcessForTesting(rhino.NilNode, root)
	return root
}

// Java: SerializeTypesToPointersTest.findInStringPool(String)
func (c *SerializeTypesToPointersTest) findInStringPool(str string) int {
	return c.stringPoolBuilder.Put(str)
}

// The list holds the Integers a proto's repeated int32 field holds (int32
// in Go).
//
// Java: SerializeTypesToPointersTest.findAllInStringPool(String...)
func (c *SerializeTypesToPointersTest) findAllInStringPool(str ...string) []int32 {
	out := make([]int32, 0, len(str))
	for _, s := range str {
		out = append(out, int32(c.findInStringPool(s)))
	}
	return out
}

// serializeTypesToPointersTest_get is Map.get on the type pointer map: a nil
// pointer stands for Java's null for a missing key.
//
// Go-only: test-helper — LinkedIdentityHashMap.Get returns the zero value for
// a missing key.
func serializeTypesToPointersTest_get(m *base.LinkedIdentityHashMap[rhino.JSType, int], key rhino.JSType) *int {
	var found *int
	m.ForEach(func(k rhino.JSType, v int) {
		if found == nil && k == key {
			found = &v
		}
	})
	return found
}

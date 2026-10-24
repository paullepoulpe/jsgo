// Port of com/google/javascript/jscomp/serialization/JSTypeReconserializerTest.java @ closure-compiler v20261015

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// The Java @Test methods not named test* get a "Test" prefix so that
// RunJUnit4 runs them.
//
// Java: public final class JSTypeReconserializerTest extends CompilerTestCase
type JSTypeReconserializerTest struct {
	CompilerTestCaseBase

	// individual test cases may override this
	typesToForwardDeclare   *jutil.LinkedHashSet[string]
	shouldSerializeProperty func(string) bool

	typePool          *TypePool
	stringPoolBuilder *StringPool_Builder
	labelToPointer    *jutil.LinkedHashMap[string, int]
}

// jsTypeReconserializerTest_field is an ObjectTypeProto field descriptor
// (ObjectTypeProto.getDescriptor().findFieldByName(...)) that ProtoTruth's
// ignoringFieldDescriptors(...) leaves out of a comparison.
//
// Go-only: language-adapter — ProtoTruth's field descriptors: the ignored
// fields are cleared on both sides before comparing with equals. Every
// field the test ignores is an ObjectTypeProto field, and ObjectTypeProto
// only occurs as TypeProto's object.
type jsTypeReconserializerTest_field int

const (
	jsTypeReconserializerTest_uuid jsTypeReconserializerTest_field = iota
	jsTypeReconserializerTest_ownProperty
	jsTypeReconserializerTest_isInvalidating
	jsTypeReconserializerTest_prototype
	jsTypeReconserializerTest_instanceType
)

// Proto fields commonly ignored in tests because hardcoding their values is brittle
var (
	// Java: JSTypeReconserializerTest.OBJECT_UUID
	JSTypeReconserializerTest_OBJECT_UUID = jsTypeReconserializerTest_uuid
	// Java: JSTypeReconserializerTest.OBJECT_PROPERTIES
	JSTypeReconserializerTest_OBJECT_PROPERTIES = jsTypeReconserializerTest_ownProperty
	// Java: JSTypeReconserializerTest.IS_INVALIDATING
	JSTypeReconserializerTest_IS_INVALIDATING = jsTypeReconserializerTest_isInvalidating

	// Java: JSTypeReconserializerTest.BRITTLE_TYPE_FIELDS
	JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS = []jsTypeReconserializerTest_field{
		JSTypeReconserializerTest_OBJECT_UUID, JSTypeReconserializerTest_OBJECT_PROPERTIES}
)

// Java: JSTypeReconserializerTest()
func NewJSTypeReconserializerTest() *JSTypeReconserializerTest {
	c := &JSTypeReconserializerTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestJSTypeReconserializerTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewJSTypeReconserializerTest() })
}

// Java: JSTypeReconserializerTest.setUp()
func (c *JSTypeReconserializerTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.typesToForwardDeclare = jutil.NewLinkedHashSet[string]()
	c.shouldSerializeProperty = func(propertyName string) bool { return true }
	c.EnableSourceInformationAnnotator()
}

// jsTypeReconserializerTest_labelFinder is the anonymous
// NodeTraversal.AbstractPostOrderCallback of getProcessor.
type jsTypeReconserializerTest_labelFinder struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer      *JSTypeReconserializerTest
	serializer *JSTypeReconserializer
}

// Java: JSTypeReconserializerTest.getProcessor(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (f *jsTypeReconserializerTest_labelFinder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.GetJSType() != nil {
		pointer := f.serializer.SerializeType(n.GetJSType())

		if n.GetParent().IsExprResult() && n.GetGrandparent().IsLabel() {
			labelName := n.GetParent().GetPrevious().GetString()
			f.outer.labelToPointer.Put(labelName, pointer)
		}
	}
}

// Java: JSTypeReconserializerTest.getProcessor(Compiler)
func (c *JSTypeReconserializerTest) GetProcessor(compiler *Compiler) CompilerPass {
	c.stringPoolBuilder = StringPool_Builder_()
	c.labelToPointer = jutil.NewLinkedHashMap[string, int]()

	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		serializer :=
			JSTypeReconserializer_Create(
				compiler.GetTypeRegistry(),
				NewInvalidatingTypes_Builder(compiler.GetTypeRegistry()).
					AddAllTypeMismatches(compiler.GetTypeMismatches()).
					Build(),
				c.stringPoolBuilder,
				c.shouldSerializeProperty,
				SerializationOptions_Builder_().
					SetIncludeDebugInfo(true).
					SetRunValidation(true).
					Build())

		NodeTraversal_TraverseRoots(
			compiler,
			&jsTypeReconserializerTest_labelFinder{outer: c, serializer: serializer},
			externs,
			root)

		c.typePool = serializer.GenerateTypePool()
	})
}

// Java: JSTypeReconserializerTest.createCompiler()
func (c *JSTypeReconserializerTest) CreateCompiler() *Compiler {
	compiler := c.CompilerTestCaseBase.CreateCompiler()
	for typeName := range c.typesToForwardDeclare.All() {
		compiler.ForwardDeclareType(typeName)
	}
	return compiler
}

// Java: JSTypeReconserializerTest.testNativeTypesAreNotSerialized()
func (c *JSTypeReconserializerTest) TestNativeTypesAreNotSerialized() {
	jsTypeReconserializerTest_assertThat(c.compileToTypes("const /** bigint */ x = 1n;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
	jsTypeReconserializerTest_assertThat(c.compileToTypes("const /** boolean */ x = false;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
	jsTypeReconserializerTest_assertThat(c.compileToTypes("const /** number */ x = 5;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
	jsTypeReconserializerTest_assertThat(c.compileToTypes("const /** string */ x = '';")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
	jsTypeReconserializerTest_assertThat(c.compileToTypes("const /** symbol */ x = Symbol('x');")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
	jsTypeReconserializerTest_assertThat(c.compileToTypes("const /** null */ x = null;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
	jsTypeReconserializerTest_assertThat(c.compileToTypes("const /** undefined */ x = void 0;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
	jsTypeReconserializerTest_assertThat(c.compileToTypes("const /** ? */ x = /** @type {?} */ (0);")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
	jsTypeReconserializerTest_assertThat(c.compileToTypes("const /** * */ x = /** @type {?} */ (0);")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
}

// Java: JSTypeReconserializerTest.testNativeTypes_serializesFunctionToTopFunction()
func (c *JSTypeReconserializerTest) TestNativeTypes_serializesFunctionToTopFunction() {
	_ = c.compileToTypePool(
		"function foo() {} /** @type {!Function} */ var goo; FOO: foo; GOO: goo;")
	fooPointer := int(c.pointerForLabel("FOO"))
	gooPointer := int(c.pointerForLabel("GOO"))
	AssertThat(TypePointers_IsAxiomatic(fooPointer)).IsTrue()
	AssertThat(TypePointers_IsAxiomatic(gooPointer)).IsTrue()
	AssertThat(TypePointers_OFFSET_TO_AXIOMATIC_COLOR[fooPointer]).
		IsEqualTo(colors.StandardColors_TOP_FUNCTION)
	AssertThat(TypePointers_OFFSET_TO_AXIOMATIC_COLOR[gooPointer]).
		IsEqualTo(colors.StandardColors_TOP_FUNCTION)
}

// Java: JSTypeReconserializerTest.testNativeTypes_serializesGbigint()
func (c *JSTypeReconserializerTest) TestNativeTypes_serializesGbigint() {
	_ = c.compileToTypePool(
		"/** @interface */ function gbigint() {} /** @type {!gbigint} */ var gb; GB: gb;")
	gbPointer := int(c.pointerForLabel("GB"))
	AssertThat(TypePointers_IsAxiomatic(gbPointer)).IsTrue()
	AssertThat(TypePointers_OFFSET_TO_AXIOMATIC_COLOR[gbPointer]).
		IsEqualTo(colors.StandardColors_GBIGINT)
}

// Java: JSTypeReconserializerTest.testObjectTypeSerializationFormat()
func (c *JSTypeReconserializerTest) TestObjectTypeSerializationFormat() {
	typePool :=
		c.compileToTypes(
			`class Foo { m() {} }
new Foo().m();

FOO: new Foo();
FOO_PROTOTYPE: Foo.prototype;
`)

	jsTypeReconserializerTest_assertThat(typePool).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_OBJECT_UUID).
		ContainsAtLeast(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(
					ObjectTypeProto_NewBuilder().
						AddPrototype(c.pointerForLabel("FOO_PROTOTYPE")).
						AddInstanceType(c.pointerForLabel("FOO")).
						SetMarkedConstructor(true).
						AddOwnProperty(c.findInStringPool("prototype"))).
				Build(),
			TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build(),
			TypeProto_NewBuilder().
				SetObjectFromBuilder(
					ObjectTypeProto_NewBuilder().
						AddOwnProperty(c.findInStringPool("constructor")).
						AddOwnProperty(c.findInStringPool("m"))).
				Build(),
			TypeProto_NewBuilder().
				SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetIsInvalidating(true)).
				Build())
}

// Java: JSTypeReconserializerTest.testObjectTypeSerializationFormat_canAvoidSerializingProperty()
func (c *JSTypeReconserializerTest) TestObjectTypeSerializationFormat_canAvoidSerializingProperty() {
	propertiesToSerialize := jutil.NewLinkedHashSet("constructor", "serializeMe")
	c.shouldSerializeProperty = propertiesToSerialize.Contains

	typePool :=
		c.compileToTypes(
			`class Foo {
  serializeMe() {}
  doNotSerializeMe() {}
}
new Foo().serializeMe();
new Foo().doNotSerializeMe();
`)

	jsTypeReconserializerTest_assertThat(typePool).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_OBJECT_UUID).
		Contains(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(
					ObjectTypeProto_NewBuilder().
						AddOwnProperty(c.findInStringPool("constructor")).
						AddOwnProperty(c.findInStringPool("serializeMe"))).
				// the "doNotSerializeMe" property is not present
				Build())
}

// Java: JSTypeReconserializerTest.testDisambiguationEdges_pointFromInstanceToPrototype()
func (c *JSTypeReconserializerTest) TestDisambiguationEdges_pointFromInstanceToPrototype() {
	typePool := c.compileToTypePool("class Foo { m() {} } new Foo().m(); FOO: new Foo();")

	jsTypeReconserializerTest_assertThat(c.getNonPrimitiveSupertypesFor(typePool, "FOO")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		Contains(TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build())
}

// Java: JSTypeReconserializerTest.testDisambiguationEdges_pointFromInterfaceToPrototype()
func (c *JSTypeReconserializerTest) TestDisambiguationEdges_pointFromInterfaceToPrototype() {
	typePool :=
		c.compileToTypePool(
			`/** @interface */ class IFoo { m() {} }
let /** !IFoo */ x;
IFOO: x;
`)

	jsTypeReconserializerTest_assertThat(c.getNonPrimitiveSupertypesFor(typePool, "IFOO")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		Contains(TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build())
}

// Java: JSTypeReconserializerTest.testDisambiguationEdges_pointFromInstanceToInterface()
func (c *JSTypeReconserializerTest) TestDisambiguationEdges_pointFromInstanceToInterface() {
	typePool :=
		c.compileToTypePool(
			`/** @interface */ class IFoo {}
/** @implements {IFoo} */ class Foo {}
FOO: new Foo();
`)

	jsTypeReconserializerTest_assertThat(c.getNonPrimitiveSupertypesFor(typePool, "FOO")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		IgnoringFieldDescriptors(
			jsTypeReconserializerTest_prototype,
			jsTypeReconserializerTest_instanceType).
		Contains(TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build())
}

// Java: JSTypeReconserializerTest.testDisambiguationEdges_pointFromSubctorToSuperctor()
func (c *JSTypeReconserializerTest) TestDisambiguationEdges_pointFromSubctorToSuperctor() {
	typePool := c.compileToTypePool("class Foo {} class Bar extends Foo {} BAR_CTOR: Bar;")

	jsTypeReconserializerTest_assertThat(c.getNonPrimitiveSupertypesFor(typePool, "BAR_CTOR")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		IgnoringFieldDescriptors(
			jsTypeReconserializerTest_prototype,
			jsTypeReconserializerTest_instanceType).
		ContainsExactly(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetMarkedConstructor(true)).
				Build())
}

// Java: JSTypeReconserializerTest.testDisambiguationEdges_dont_pointFromPrototypeToInstance()
func (c *JSTypeReconserializerTest) TestDisambiguationEdges_dont_pointFromPrototypeToInstance() {
	typePool :=
		c.compileToTypePool("class Foo { m() {} } new Foo().m(); FOO_PROTOTYPE: Foo.prototype;")

	AssertThat(c.getNonPrimitiveSupertypesFor(typePool, "FOO_PROTOTYPE")).IsEmpty()
}

// Java: JSTypeReconserializerTest.testMultipleAnonymousTypesDoesntCrash()
func (c *JSTypeReconserializerTest) TestMultipleAnonymousTypesDoesntCrash() {
	c.compileToTypes("(() => 5); (() => 'str');")
}

// Java: JSTypeReconserializerTest.testUnion_proxyLikeTypes_dontPreventUnionCollapsing()
func (c *JSTypeReconserializerTest) TestUnion_proxyLikeTypes_dontPreventUnionCollapsing() {
	typePool :=
		c.compileToTypes(
			`/** @enum {boolean|number} */
const Foo = {};
/** @enum {string|number|!Foo} */
const Bar = {};

let /** symbol|!Bar|string */ x;
`)

	unionsContainingUnions := []*UnionTypeProto{}
	for _, t := range typePool {
		if !t.HasUnion() {
			continue
		}
		u := t.GetUnion()
		for _, p := range u.GetUnionMemberList() {
			if !TypePointers_IsAxiomatic(int(p)) && typePool[TypePointers_TrimOffset(int(p))].HasUnion() {
				unionsContainingUnions = append(unionsContainingUnions, u)
				break
			}
		}
	}
	AssertThat(unionsContainingUnions).IsEmpty()

	AssertThat(typePool).
		Contains(
			TypeProto_NewBuilder().
				SetUnion(
					UnionTypeProto_NewBuilder().
						AddUnionMember(PrimitiveType_BOOLEAN_TYPE.GetNumber()).
						AddUnionMember(PrimitiveType_STRING_TYPE.GetNumber()).
						AddUnionMember(PrimitiveType_NUMBER_TYPE.GetNumber()).
						AddUnionMember(PrimitiveType_SYMBOL_TYPE.GetNumber()).
						Build()).
				Build())
}

// Java: JSTypeReconserializerTest.testSerializedInstanceTypeHasClassName()
func (c *JSTypeReconserializerTest) TestSerializedInstanceTypeHasClassName() {
	jsTypeReconserializerTest_assertThat(c.compileToTypes("/** @constructor */ function Foo() {} new Foo;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		Contains(TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build())
}

// Java: JSTypeReconserializerTest.testEnumElementSerializesToPrimitive()
func (c *JSTypeReconserializerTest) TestEnumElementSerializesToPrimitive() {
	// Serialize E but treat E.A as a number
	jsTypeReconserializerTest_assertThat(c.compileToTypes("/** @enum {number} */ const E = {A: 0, B: 1}; E.A;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		Contains(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetPropertiesKeepOriginalName(true)).
				Build())
}

// Java: JSTypeReconserializerTest.testBoxableScalarTypesKeepOriginalName()
func (c *JSTypeReconserializerTest) TestBoxableScalarTypesKeepOriginalName() {
	// Serialize E but treat E.A as a number
	jsTypeReconserializerTest_assertThat(c.compileToTypes("/** @enum {number} */ const E = {A: 0, B: 1}; E.A;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		Contains(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetPropertiesKeepOriginalName(true)).
				Build())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.marksClosureAssertions()
func (c *JSTypeReconserializerTest) TestMarksClosureAssertions() {
	jsTypeReconserializerTest_assertThat(
		c.compileToTypes(
			`/** @closurePrimitive {asserts.truthy} */
function assert(x) { if (!x) { throw new Error(); }}

/** @closurePrimitive {asserts.fail} */
function fail() { throw new Error(); }
`)).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsAtLeast(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(
					ObjectTypeProto_NewBuilder().SetClosureAssert(true).SetIsInvalidating(true)).
				Build(),
			TypeProto_NewBuilder().
				SetObjectFromBuilder(
					// ClosurePrimitive.ASSERTS_FAIL is not a removable Closure assertion
					ObjectTypeProto_NewBuilder().SetClosureAssert(false).SetIsInvalidating(true)).
				Build())
}

// Java: JSTypeReconserializerTest.testSerializeObjectLiteralTypesAsInvalidating()
func (c *JSTypeReconserializerTest) TestSerializeObjectLiteralTypesAsInvalidating() {
	typePool := c.compileToTypes("const /** {x: string} */ obj = {x: ''};")
	jsTypeReconserializerTest_assertThat(typePool).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_OBJECT_UUID).
		Contains(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(
					ObjectTypeProto_NewBuilder().
						SetIsInvalidating(true).
						AddOwnProperty(c.findInStringPool("x"))).
				Build())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.serializesTemplateTypesAsUnknown()
func (c *JSTypeReconserializerTest) TestSerializesTemplateTypesAsUnknown() {
	typePool :=
		c.compileToTypes(
			`/** @interface @template T */
class Foo {
// Add in this reference to ` + "`T`" + ` so that it is seen by the serializer.
  constructor() { /** @type {T} */ this.x; }
}

/** @type {!Foo<string>} */ let foo;
FOO: foo;
FOO_PROTOTYPE: Foo.prototype;
`)
	jsTypeReconserializerTest_assertThat(typePool).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsAtLeast(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(
					ObjectTypeProto_NewBuilder().
						AddPrototype(c.pointerForLabel("FOO_PROTOTYPE")).
						AddInstanceType(c.pointerForLabel("FOO"))).
				Build(),
			TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build(),
			TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build())
	// Verify no additional objects were serialized to represent the template type: we only have
	// the well-known types and the three TypeProtos tested for above.
	AssertThat(typePool).HasSize(len(c.nativeObjects()) + 3)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.reconcile_literalFunctions_toNativeType()
func (c *JSTypeReconserializerTest) TestReconcile_literalFunctions_toNativeType() {
	output :=
		c.compileToTypes(
			`const ns = {};

ns.f = x => x;
ns.f.a = 0;

ns.f = (/** number */ x) => x;
ns.f.b = 1;

ns.f = (/** string */ x) => x;
ns.f.c = 2;
`)

	jsTypeReconserializerTest_assertThat(output).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_OBJECT_UUID).
		IgnoringFieldDescriptors(jsTypeReconserializerTest_prototype).
		ContainsExactlyElementsIn(
			append(append([]*TypeProto{}, c.nativeObjects()...),
				TypeProto_NewBuilder().
					SetObjectFromBuilder(
						ObjectTypeProto_NewBuilder().
							SetIsInvalidating(true).
							AddOwnProperty(c.findInStringPool("f"))).
					Build()))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.discards_differentTemplatizationsOfSameRawType()
func (c *JSTypeReconserializerTest) TestDiscards_differentTemplatizationsOfSameRawType() {
	typePool :=
		c.compileToTypes(
			`/** @interface @template T */
class Foo {}
var /** !Foo<string> */ fooString;
var /** !Foo<number> */ fooNumber;

FOO: fooString;
FOO_PROTOTYPE: Foo.prototype;
`)

	jsTypeReconserializerTest_assertThat(typePool).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsAtLeast(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(
					ObjectTypeProto_NewBuilder().
						AddPrototype(c.pointerForLabel("FOO_PROTOTYPE")).
						AddInstanceType(c.pointerForLabel("FOO"))).
				Build(),
			TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build(),
			TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.doesntCrashOnFunctionWithNonObjectThisType()
func (c *JSTypeReconserializerTest) TestDoesntCrashOnFunctionWithNonObjectThisType() {
	typePool :=
		c.compileToTypes(
			`var /** function(new: AllType) */ x = function() {};
/** @typedef {*} */ let AllType;
`)

	jsTypeReconserializerTest_assertThat(typePool).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		Contains(
			TypeProto_NewBuilder().
				SetObject(
					ObjectTypeProto_NewBuilder().
						SetIsInvalidating(true).
						AddPrototype(PrimitiveType_UNKNOWN_TYPE.GetNumber()).
						AddInstanceType(PrimitiveType_UNKNOWN_TYPE.GetNumber()).
						SetMarkedConstructor(true).
						Build()).
				Build())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.supertypeListForInterface_includesIllegalExtendedTypes()
func (c *JSTypeReconserializerTest) TestSupertypeListForInterface_includesIllegalExtendedTypes() {
	source :=
		`class Foo {}
/** @interface */
class Bar {}
/** @typedef {{x: string}} */
let RecordType;
/** @interface
// [JSC_CONFLICTING_EXTENDS_TYPE]: interfaces can only extend interfaces
 * @extends {Foo}
 * @extends {Bar}
 * @extends {RecordType}
 * @suppress {checkTypes}
*/
class Baz {}

let /** !Baz */ baz;
BAZ: baz;
`
	typePool := c.compileToTypePool(source)

	jsTypeReconserializerTest_assertThat(c.getNonPrimitiveSupertypesFor(typePool, "BAZ")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactly(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetIsInvalidating(true)).
				Build(),
			TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build(),
			TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build(),
			TypeProto_NewBuilder().SetObject(ObjectTypeProto_GetDefaultInstance()).Build())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.doesntSerializeNullOrUndefinedUnion()
func (c *JSTypeReconserializerTest) TestDoesntSerializeNullOrUndefinedUnion() {
	jsTypeReconserializerTest_assertThat(c.compileToTypes("var /** null|undefined */ x;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.doesntSerializeForwardDeclaredType()
func (c *JSTypeReconserializerTest) TestDoesntSerializeForwardDeclaredType() {
	c.typesToForwardDeclare = jutil.NewLinkedHashSet("forward.declared.TypeProto")
	jsTypeReconserializerTest_assertThat(c.compileToTypes("var /** !forward.declared.TypeProto */ x;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.doesntSerializeMissingType()
func (c *JSTypeReconserializerTest) TestDoesntSerializeMissingType() {
	c.IgnoreWarnings(DiagnosticGroups_CHECK_TYPES)
	jsTypeReconserializerTest_assertThat(c.compileToTypes("var /** !missing.TypeProto */ x;")).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		ContainsExactlyElementsIn(c.nativeObjects())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.serializesSourceOfSyntheticCodeDoesntCrash()
func (c *JSTypeReconserializerTest) TestSerializesSourceOfSyntheticCodeDoesntCrash() {
	c.EnsureLibraryInjected("base")
	c.compileToTypePoolWithExterns(NewTestExternsBuilder().AddMath().AddObject().Build(), "let s = 'hello';")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.reconcile_setInstanceColors()
func (c *JSTypeReconserializerTest) TestReconcile_setInstanceColors() {
	// When
	typePool :=
		c.compileToTypes(
			`class Foo { }
class Bar { }

let /** function(new:Foo) */ x;
let /** function(new:Bar) */ y;

FOO: new Foo();
BAR: new Bar();
`)

	// Then
	jsTypeReconserializerTest_assertThat(typePool).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_IS_INVALIDATING).
		Contains(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(
					ObjectTypeProto_NewBuilder().
						SetMarkedConstructor(true).
						AddInstanceType(c.pointerForLabel("FOO")).
						AddInstanceType(c.pointerForLabel("BAR")).
						AddPrototype(PrimitiveType_UNKNOWN_TYPE.GetNumber())).
				Build())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.reconcile_addPrototypes_noActualTestcase()
func (c *JSTypeReconserializerTest) TestReconcile_addPrototypes_noActualTestcase() {
	// TODO(b/185519307): Come up with a situation where this can happen.
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.reconcile_setOwnProperties()
func (c *JSTypeReconserializerTest) TestReconcile_setOwnProperties() {
	// When
	typePool :=
		c.compileToTypes(
			`var Foo = class { a() { } }
Foo = class { b() { } }
`)

	// Then
	jsTypeReconserializerTest_assertThat(typePool).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_OBJECT_UUID).
		Contains(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(
					ObjectTypeProto_NewBuilder().
						AddOwnProperty(c.findInStringPool("a")).
						AddOwnProperty(c.findInStringPool("b")).
						AddOwnProperty(c.findInStringPool("constructor"))).
				Build())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.reconcile_setClosureAssert()
func (c *JSTypeReconserializerTest) TestReconcile_setClosureAssert() {
	// When
	typePool :=
		c.compileToTypes(
			`var /** ? */ assertFoo;

assertFoo = function(/** ? */ a) { }

/** @closurePrimitive {asserts.truthy} */
assertFoo = function(/** ? */ a) { }

assertFoo = function(/** ? */ a) { }
`)

	// Then
	jsTypeReconserializerTest_assertThat(typePool).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_IS_INVALIDATING).
		Contains(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetClosureAssert(true)).
				Build())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.reconcile_setConstructor_noActualTestcase()
func (c *JSTypeReconserializerTest) TestReconcile_setConstructor_noActualTestcase() {
	// TODO(b/185519307): Come up with a situation where this can happen.
	// Being a constructor is part of the ColorId
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.reconcile_setInvalidating()
func (c *JSTypeReconserializerTest) TestReconcile_setInvalidating() {
	// When
	typePool :=
		c.compileToTypes(
			`var /** ? */ Foo;

Foo = function() { };

Foo = class { };
`)

	// Then
	jsTypeReconserializerTest_assertThat(typePool).
		IgnoringFieldDescriptors(JSTypeReconserializerTest_BRITTLE_TYPE_FIELDS...).
		Contains(
			TypeProto_NewBuilder().
				SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetIsInvalidating(true)).
				Build())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.reconcile_setPropertiesKeepOriginalName_noActualTestcase()
func (c *JSTypeReconserializerTest) TestReconcile_setPropertiesKeepOriginalName_noActualTestcase() {
	// TODO(b/185519307): Come up with a situation where this can happen.
	// Enum types have a unique reference names format which is part of the ColorId
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: JSTypeReconserializerTest.debugInfo_mismatchLocations()
func (c *JSTypeReconserializerTest) TestDebugInfo_mismatchLocations() {
	typePool :=
		c.compileToTypePool(
			`class A { }
class B { }
class C { }

/**
 * @suppress {checkTypes}
 * @type {(!A|!B)}
 */
const x = new C();

A_CTOR: A;
B_CTOR: B;
C_CTOR: C;
`)

	aCtor :=
		c.typePool.GetType(TypePointers_TrimOffset(int(c.pointerForLabel("A_CTOR")))).GetObject()
	bCtor :=
		c.typePool.GetType(TypePointers_TrimOffset(int(c.pointerForLabel("B_CTOR")))).GetObject()
	cCtor :=
		c.typePool.GetType(TypePointers_TrimOffset(int(c.pointerForLabel("C_CTOR")))).GetObject()

	AssertThat(typePool.GetDebugInfo().GetMismatchList()).
		Contains(
			TypePool_DebugInfo_Mismatch_NewBuilder().
				SetSourceRef("testcode:9:10").
				AddInvolvedColor(c.pointerForLabel("A_CTOR")).
				AddInvolvedColor(aCtor.GetInstanceTypeList()[0]).
				AddInvolvedColor(aCtor.GetPrototypeList()[0]).
				AddInvolvedColor(c.pointerForLabel("B_CTOR")).
				AddInvolvedColor(bCtor.GetInstanceTypeList()[0]).
				AddInvolvedColor(bCtor.GetPrototypeList()[0]).
				AddInvolvedColor(c.pointerForLabel("C_CTOR")).
				AddInvolvedColor(cCtor.GetInstanceTypeList()[0]).
				AddInvolvedColor(cCtor.GetPrototypeList()[0]).
				Build())
}

// Java: JSTypeReconserializerTest.compileToTypes(String)
func (c *JSTypeReconserializerTest) compileToTypes(source string) []*TypeProto {
	return c.compileToTypePool(source).GetTypeList()
}

// Java: JSTypeReconserializerTest.compileToTypePool(String)
func (c *JSTypeReconserializerTest) compileToTypePool(source string) *TypePool {
	return c.compileToTypePoolWithExterns(CompilerTestCase_DEFAULT_EXTERNS, source)
}

// Java: JSTypeReconserializerTest.compileToTypePool(String, String)
func (c *JSTypeReconserializerTest) compileToTypePoolWithExterns(externs string, source string) *TypePool {
	c.TestNoWarning(c.Externs(externs), c.Srcs(source))
	return c.typePool
}

// Go-only: language-adapter — returns the Java int as the int32 the proto
// builders take.
//
// Java: JSTypeReconserializerTest.findInStringPool(String)
func (c *JSTypeReconserializerTest) findInStringPool(str string) int32 {
	return int32(c.stringPoolBuilder.Put(str))
}

// Go-only: language-adapter — returns the Java int as the int32 the proto
// builders take.
//
// Java: JSTypeReconserializerTest.pointerForLabel(String)
func (c *JSTypeReconserializerTest) pointerForLabel(label string) int32 {
	jutil.CheckStateMsg1(c.labelToPointer.ContainsKey(label), "Cannot find label %s", label)
	return int32(c.labelToPointer.Get(label))
}

// Returns the types that are serialized for every compilation, even given an empty source
//
// Java: JSTypeReconserializerTest.nativeObjects()
func (c *JSTypeReconserializerTest) nativeObjects() []*TypeProto {
	return append([]*TypeProto{}, c.compileToTypes("")...)
}

// Accepts a label name in source code, and finds its primitive supertypes.
//
// Java: JSTypeReconserializerTest.getNonPrimitiveSupertypesFor(TypePool, String)
func (c *JSTypeReconserializerTest) getNonPrimitiveSupertypesFor(typePool *TypePool, labelName string) []*TypeProto {
	supertypes := []*TypeProto{}
	typeOffset := c.pointerForLabel(labelName)

	for _, edge := range typePool.GetDisambiguationEdgesList() {
		if edge.GetSubtype() != typeOffset {
			continue
		}
		if TypePointers_IsAxiomatic(int(edge.GetSupertype())) {
			continue // Skip axiomatic supertpyes as they don't have a TypeProto
		}
		supertypes = append(supertypes, typePool.GetType(TypePointers_TrimOffset(int(edge.GetSupertype()))))
	}
	return supertypes
}

// jsTypeReconserializerTest_protoSubject is ProtoTruth's
// IterableOfProtosSubject for a List<TypeProto>.
//
// Go-only: language-adapter — see jsTypeReconserializerTest_field.
type jsTypeReconserializerTest_protoSubject struct {
	actual  []*TypeProto
	ignored []jsTypeReconserializerTest_field
}

// jsTypeReconserializerTest_assertThat is ProtoTruth.assertThat(Iterable<TypeProto>).
func jsTypeReconserializerTest_assertThat(actual []*TypeProto) *jsTypeReconserializerTest_protoSubject {
	return &jsTypeReconserializerTest_protoSubject{actual: actual}
}

// IgnoringFieldDescriptors is ignoringFieldDescriptors(FieldDescriptor...):
// the ignored fields accumulate.
func (s *jsTypeReconserializerTest_protoSubject) IgnoringFieldDescriptors(
	fields ...jsTypeReconserializerTest_field) *jsTypeReconserializerTest_protoSubject {
	return &jsTypeReconserializerTest_protoSubject{
		actual: s.actual, ignored: append(append([]jsTypeReconserializerTest_field{}, s.ignored...), fields...)}
}

func (s *jsTypeReconserializerTest_protoSubject) clearAll(types []*TypeProto) []*TypeProto {
	out := []*TypeProto{}
	for _, t := range types {
		out = append(out, s.clear(t))
	}
	return out
}

// clear returns type_ with the ignored ObjectTypeProto fields cleared.
func (s *jsTypeReconserializerTest_protoSubject) clear(type_ *TypeProto) *TypeProto {
	if !type_.HasObject() || len(s.ignored) == 0 {
		return type_
	}
	object := type_.GetObject().ToBuilder()
	for _, f := range s.ignored {
		switch f {
		case jsTypeReconserializerTest_uuid:
			object.ClearUuid()
		case jsTypeReconserializerTest_ownProperty:
			object.ClearOwnProperty()
		case jsTypeReconserializerTest_isInvalidating:
			object.ClearIsInvalidating()
		case jsTypeReconserializerTest_prototype:
			object.ClearPrototype()
		case jsTypeReconserializerTest_instanceType:
			object.ClearInstanceType()
		}
	}
	return type_.ToBuilder().SetObject(object.Build()).Build()
}

func (s *jsTypeReconserializerTest_protoSubject) Contains(expected *TypeProto) {
	AssertThat(s.clearAll(s.actual)).Contains(s.clear(expected))
}

func (s *jsTypeReconserializerTest_protoSubject) ContainsAtLeast(expected ...*TypeProto) {
	AssertThat(s.clearAll(s.actual)).ContainsAtLeastElementsIn(s.clearAll(expected))
}

func (s *jsTypeReconserializerTest_protoSubject) ContainsExactly(expected ...*TypeProto) {
	AssertThat(s.clearAll(s.actual)).ContainsExactlyElementsIn(s.clearAll(expected))
}

func (s *jsTypeReconserializerTest_protoSubject) ContainsExactlyElementsIn(expected []*TypeProto) {
	AssertThat(s.clearAll(s.actual)).ContainsExactlyElementsIn(s.clearAll(expected))
}

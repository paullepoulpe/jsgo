// Port of com/google/javascript/jscomp/TypeTransformationTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/TypeTransformationTest.java)

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// Java: public final class TypeTransformationTest extends CompilerTypeTestCase
type TypeTransformationTest struct {
	CompilerTypeTestCaseBase

	typeVars *jutil.LinkedHashMap[string, rhino.JSType]
	nameVars *jutil.LinkedHashMap[string, string]
}

// Java: private static JSType TypeTransformationTest.recordTypeTest
var typeTransformationTest_recordTypeTest rhino.JSType

// Java: private static JSType TypeTransformationTest.nestedRecordTypeTest
var typeTransformationTest_nestedRecordTypeTest rhino.JSType

// Java: private static JSType TypeTransformationTest.asynchRecord
var typeTransformationTest_asynchRecord rhino.JSType

// Java: static final String TypeTransformationTest.EXTRA_TYPE_DEFS
const TypeTransformationTest_EXTRA_TYPE_DEFS = "/** @type {number} */ Array.prototype.length;\n" +
	"/** @typedef {!Array<?>} */ var ArrayAlias;\n" +
	"\n" +
	"/** @constructor */\n" +
	"function Bar() {}\n" +
	"\n" +
	"/** @type {number} */\n" +
	"var n = 10;\n"

// Java: TypeTransformationTest()
func NewTypeTransformationTest() *TypeTransformationTest {
	c := &TypeTransformationTest{}
	c.InitCompilerTypeTestCaseBase(c)
	return c
}

func TestTypeTransformationTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeTransformationTest() })
}

// Java: TypeTransformationTest.setUp()
func (c *TypeTransformationTest) SetUp() {
	c.CompilerTypeTestCaseBase.SetUp()
	c.errorReporter = rhinotesting.NewTestErrorReporter()
	c.initRecordTypeTests()
	c.typeVars =
		jutil.NewImmutableMapBuilder[string, rhino.JSType]().
			Put("S", c.GetNativeStringType()).
			Put("N", c.GetNativeNumberType()).
			Put("B", c.GetNativeBooleanType()).
			Put("BOT", c.GetNativeNoType()).
			Put("TOP", c.GetNativeAllType()).
			Put("UNK", c.GetNativeUnknownType()).
			Put("CHKUNK", c.GetNativeCheckedUnknownType()).
			Put("SO", c.GetNativeStringObjectType()).
			Put("NO", c.GetNativeNumberObjectType()).
			Put("BO", c.GetNativeBooleanObjectType()).
			Put("NULL", c.GetNativeNullType()).
			Put("OBJ", c.GetNativeObjectType()).
			Put("UNDEF", c.GetNativeVoidType()).
			Put("ARR", c.GetNativeArrayType()).
			Put("ARRNUM", c.type_(c.GetNativeArrayType(), c.GetNativeNumberType())).
			Put("REC", typeTransformationTest_recordTypeTest).
			Put("NESTEDREC", typeTransformationTest_nestedRecordTypeTest).
			Put("ASYNCH", typeTransformationTest_asynchRecord).
			BuildOrThrow()
	c.nameVars =
		jutil.NewImmutableMapBuilder[string, string]().
			Put("s", "string").
			Put("n", "number").
			Put("b", "boolean").
			BuildOrThrow()
}

// Java: TypeTransformationTest.testTransformationWithValidBasicTypePredicate()
func (c *TypeTransformationTest) TestTransformationWithValidBasicTypePredicate() {
	c.testTTL(c.GetNativeNumberType(), "'number'")
}

// Java: TypeTransformationTest.testTransformationWithBasicTypePredicateWithInvalidTypename()
func (c *TypeTransformationTest) TestTransformationWithBasicTypePredicateWithInvalidTypename() {
	c.testTTL(c.GetNativeUnknownType(), "'foo'", "Reference to an unknown type name foo")
}

// Java: TypeTransformationTest.testTransformationWithSingleTypeVar()
func (c *TypeTransformationTest) TestTransformationWithSingleTypeVar() {
	c.testTTL(c.GetNativeStringType(), "S")
}

// Java: TypeTransformationTest.testTransformationWithMultipleTypeVars()
func (c *TypeTransformationTest) TestTransformationWithMultipleTypeVars() {
	c.testTTL(c.GetNativeStringType(), "S")
	c.testTTL(c.GetNativeNumberType(), "N")
}

// Java: TypeTransformationTest.testTransformationWithValidUnionTypeOnlyVars()
func (c *TypeTransformationTest) TestTransformationWithValidUnionTypeOnlyVars() {
	c.testTTL(c.union(c.GetNativeNumberType(), c.GetNativeStringType()), "union(N, S)")
}

// Java: TypeTransformationTest.testTransformationWithValidUnionTypeOnlyTypePredicates()
func (c *TypeTransformationTest) TestTransformationWithValidUnionTypeOnlyTypePredicates() {
	c.testTTL(c.union(c.GetNativeNumberType(), c.GetNativeStringType()), "union('number', 'string')")
}

// Java: TypeTransformationTest.testTransformationWithValidUnionTypeMixed()
func (c *TypeTransformationTest) TestTransformationWithValidUnionTypeMixed() {
	c.testTTL(c.union(c.GetNativeNumberType(), c.GetNativeStringType()), "union(S, 'number')")
}

// Java: TypeTransformationTest.testTransformationWithUnknownParameter()
func (c *TypeTransformationTest) TestTransformationWithUnknownParameter() {
	// Returns ? because foo is not defined
	c.testTTL(
		c.GetNativeUnknownType(),
		"union(foo, 'number')",
		"Reference to an unknown type variable foo")
}

// Java: TypeTransformationTest.testTransformationWithUnknownParameter2()
func (c *TypeTransformationTest) TestTransformationWithUnknownParameter2() {
	// Returns ? because foo is not defined
	c.testTTL(c.GetNativeUnknownType(), "union(N, 'foo')", "Reference to an unknown type name foo")
}

// Java: TypeTransformationTest.testTransformationWithNestedUnionInFirstParameter()
func (c *TypeTransformationTest) TestTransformationWithNestedUnionInFirstParameter() {
	c.testTTL(
		c.union(c.GetNativeNumberType(), c.GetNativeNullType(), c.GetNativeStringType()),
		"union(union(N, 'null'), S)")
}

// Java: TypeTransformationTest.testTransformationWithNestedUnionInSecondParameter()
func (c *TypeTransformationTest) TestTransformationWithNestedUnionInSecondParameter() {
	c.testTTL(
		c.union(c.GetNativeNumberType(), c.GetNativeNullType(), c.GetNativeStringType()),
		"union(N, union('null', S))")
}

// Java: TypeTransformationTest.testTransformationWithRepeatedTypePredicate()
func (c *TypeTransformationTest) TestTransformationWithRepeatedTypePredicate() {
	c.testTTL(c.GetNativeNumberType(), "union('number', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithUndefinedTypeVar()
func (c *TypeTransformationTest) TestTransformationWithUndefinedTypeVar() {
	c.testTTL(c.GetNativeUnknownType(), "foo", "Reference to an unknown type variable foo")
}

// Java: TypeTransformationTest.testTransformationWithTrueEqtypeConditional()
func (c *TypeTransformationTest) TestTransformationWithTrueEqtypeConditional() {
	c.testTTL(c.GetNativeStringType(), "cond(eq(N, N), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseEqtypeConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseEqtypeConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(eq(N, S), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueSubtypeConditional()
func (c *TypeTransformationTest) TestTransformationWithTrueSubtypeConditional() {
	c.testTTL(c.GetNativeStringType(), "cond( sub('Number', 'Object'), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseSubtypeConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseSubtypeConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond( sub('Number', 'String'), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueStreqConditional()
func (c *TypeTransformationTest) TestTransformationWithTrueStreqConditional() {
	c.testTTL(c.GetNativeStringType(), "cond(streq(n, 'number'), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueStreqConditional2()
func (c *TypeTransformationTest) TestTransformationWithTrueStreqConditional2() {
	c.testTTL(c.GetNativeStringType(), "cond(streq(n, n), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueStreqConditional3()
func (c *TypeTransformationTest) TestTransformationWithTrueStreqConditional3() {
	c.testTTL(c.GetNativeStringType(), "cond(streq('number', 'number'), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseStreqConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseStreqConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(streq('number', 'foo'), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseStreqConditional2()
func (c *TypeTransformationTest) TestTransformationWithFalseStreqConditional2() {
	c.testTTL(c.GetNativeNumberType(), "cond(streq(n, 'foo'), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseStreqConditional3()
func (c *TypeTransformationTest) TestTransformationWithFalseStreqConditional3() {
	c.testTTL(c.GetNativeNumberType(), "cond(streq(n, s), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithInvalidEqConditional()
func (c *TypeTransformationTest) TestTransformationWithInvalidEqConditional() {
	// It returns number since a failing boolean expression returns false
	c.testTTL(
		c.GetNativeNumberType(),
		"cond(eq(foo, S), 'string', 'number')",
		"Reference to an unknown type variable foo")
}

// Java: TypeTransformationTest.testTransformationWithInvalidStreqConditional()
func (c *TypeTransformationTest) TestTransformationWithInvalidStreqConditional() {
	// It returns number since a failing boolean expression returns false
	c.testTTL(
		c.GetNativeNumberType(),
		"cond(streq(foo, s), 'string', 'number')",
		"Reference to an unknown string variable foo")
}

// Java: TypeTransformationTest.testTransformationWithNestedExpressionInBooleanFirstParam()
func (c *TypeTransformationTest) TestTransformationWithNestedExpressionInBooleanFirstParam() {
	c.testTTL(
		c.GetNativeStringType(),
		"cond( eq( cond(eq(N, N), 'string', 'number'), 'string'),'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithNestedExpressionInBooleanSecondParam()
func (c *TypeTransformationTest) TestTransformationWithNestedExpressionInBooleanSecondParam() {
	c.testTTL(
		c.GetNativeStringType(),
		"cond( eq( 'string', cond(eq(N, N), 'string', 'number')),'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithNestedExpressionInIfBranch()
func (c *TypeTransformationTest) TestTransformationWithNestedExpressionInIfBranch() {
	c.testTTL(
		c.GetNativeStringObjectType(),
		"cond(eq(N, N), cond(eq(N, S), 'string', 'String'), 'number')")
}

// Java: TypeTransformationTest.testTransformationWithNestedExpressionInElseBranch()
func (c *TypeTransformationTest) TestTransformationWithNestedExpressionInElseBranch() {
	c.testTTL(
		c.GetNativeStringObjectType(),
		"cond(eq(N, S), 'number', cond(eq(N, S), 'string', 'String'))")
}

// Java: TypeTransformationTest.testTransformationWithMapunionMappingEverythingToString()
func (c *TypeTransformationTest) TestTransformationWithMapunionMappingEverythingToString() {
	c.testTTL(c.GetNativeStringType(), "mapunion(union(S, N), (x) => S)")
}

// Java: TypeTransformationTest.testTransformationWithMapunionIdentity()
func (c *TypeTransformationTest) TestTransformationWithMapunionIdentity() {
	c.testTTL(c.union(c.GetNativeNumberType(), c.GetNativeStringType()), "mapunion(union(N, S), (x) => x)")
}

// Java: TypeTransformationTest.testTransformationWithMapunionWithUnionEvaluatedToANonUnion()
func (c *TypeTransformationTest) TestTransformationWithMapunionWithUnionEvaluatedToANonUnion() {
	c.testTTL(c.GetNativeNumberType(), "mapunion(union(N, 'number'), (x) => x)")
}

// Java: TypeTransformationTest.testTransformationWithMapunionFilterWithOnlyString()
func (c *TypeTransformationTest) TestTransformationWithMapunionFilterWithOnlyString() {
	c.testTTL(c.GetNativeStringType(), "mapunion(union(S, B, N), (x) => cond(eq(x, S), x, BOT))")
}

// Java: TypeTransformationTest.testTransformationWithMapunionOnSingletonStringToNumber()
func (c *TypeTransformationTest) TestTransformationWithMapunionOnSingletonStringToNumber() {
	c.testTTL(c.GetNativeNumberType(), "mapunion(S, (x) => cond(eq(x, S), N, BOT))")
}

// Java: TypeTransformationTest.testTransformationWithNestedUnionInMapunionFilterString()
func (c *TypeTransformationTest) TestTransformationWithNestedUnionInMapunionFilterString() {
	c.testTTL(
		c.union(c.GetNativeNumberType(), c.GetNativeBooleanType()),
		"mapunion(union(union(S, B), union(N, S)), (x) => cond(eq(x, S), BOT, x))")
}

// Java: TypeTransformationTest.testTransformationWithNestedMapunionInMapFunctionBody()
func (c *TypeTransformationTest) TestTransformationWithNestedMapunionInMapFunctionBody() {
	c.testTTL(
		c.GetNativeStringType(),
		"mapunion(union(S, B),\n"+
			"(x) => mapunion(union(S, N),\n"+
			"(y) => cond(eq(x, y), x, BOT)))\n")
}

// Java: TypeTransformationTest.testTransformationWithObjectUseCase()
func (c *TypeTransformationTest) TestTransformationWithObjectUseCase() {
	c.testTTL(
		c.GetNativeObjectType(),
		"mapunion(\n"+
			"union(S, N, B, NULL, UNDEF, ARR),\n"+
			"(x) =>\n"+
			"cond(eq(x, S), SO,\n"+
			"cond(eq(x, N), NO,\n"+
			"cond(eq(x, B), BO,\n"+
			"cond(eq(x, NULL), OBJ,\n"+
			"cond(eq(x, UNDEF), OBJ,\n"+
			"x ))))))\n")
}

// none() is evaluated to bottom in TTL expressions, but if the overall expression evaluates
// to bottom, we return unknown to the context.
// Java: TypeTransformationTest.testTransformatioWithNoneType()
func (c *TypeTransformationTest) TestTransformatioWithNoneType() {
	c.testTTL(c.GetNativeUnknownType(), "none()")
}

// The conditional is true, so we follow the THEN branch, and the bottom result is returned
// to the context as unknown.
// Java: TypeTransformationTest.testTransformatioWithNoneTypeInConditional()
func (c *TypeTransformationTest) TestTransformatioWithNoneTypeInConditional() {
	c.testTTL(c.GetNativeUnknownType(), "cond(eq(BOT, none()), none(), N)")
}

// Java: TypeTransformationTest.testTransformatioWithNoneTypeInMapunionFilterString()
func (c *TypeTransformationTest) TestTransformatioWithNoneTypeInMapunionFilterString() {
	c.testTTL(c.GetNativeStringType(), "mapunion(union(S, B, N), (x) => cond(eq(x, S), x, none()))")
}

// Java: TypeTransformationTest.testTransformatioWithAllType()
func (c *TypeTransformationTest) TestTransformatioWithAllType() {
	c.testTTL(c.GetNativeAllType(), "all()")
}

// Java: TypeTransformationTest.testTransformatioWithAllTypeInConditional()
func (c *TypeTransformationTest) TestTransformatioWithAllTypeInConditional() {
	c.testTTL(c.GetNativeAllType(), "cond(eq(TOP, all()), all(), N)")
}

// Java: TypeTransformationTest.testTransformatioWithAllTypeMixUnion()
func (c *TypeTransformationTest) TestTransformatioWithAllTypeMixUnion() {
	c.testTTL(c.GetNativeAllType(), "mapunion(union(S, B, N), (x) => cond(eq(x, S), x, all()))")
}

// Java: TypeTransformationTest.testTransformatioWithUnknownType()
func (c *TypeTransformationTest) TestTransformatioWithUnknownType() {
	c.testTTL(c.GetNativeUnknownType(), "unknown()")
}

// Java: TypeTransformationTest.testTransformatioWithUnknownTypeInConditional()
func (c *TypeTransformationTest) TestTransformatioWithUnknownTypeInConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(eq(UNK, unknown()), N, S)")
}

// Java: TypeTransformationTest.testTransformatioWithUnknownTypeInMapunionStringToUnknown()
func (c *TypeTransformationTest) TestTransformatioWithUnknownTypeInMapunionStringToUnknown() {
	c.testTTL(
		c.GetNativeUnknownType(), "mapunion(union(S, B, N), (x) => cond(eq(x, S), x, unknown()))")
}

// Java: TypeTransformationTest.testTransformationWithTemplatizedType()
func (c *TypeTransformationTest) TestTransformationWithTemplatizedType() {
	c.testTTL(c.type_(c.GetNativeArrayType(), c.GetNativeNumberType()), "type('Array', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTemplatizedType2()
func (c *TypeTransformationTest) TestTransformationWithTemplatizedType2() {
	c.testTTL(c.type_(c.GetNativeArrayType(), c.GetNativeNumberType()), "type(ARR, 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTemplatizedType3()
func (c *TypeTransformationTest) TestTransformationWithTemplatizedType3() {
	c.testTTL(c.type_(c.GetNativeArrayType(), c.GetNativeNumberType()), "type(ARR, N)")
}

// Java: TypeTransformationTest.testTransformationWithTemplatizedTypeInvalidBaseType()
func (c *TypeTransformationTest) TestTransformationWithTemplatizedTypeInvalidBaseType() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"type('string', 'number')",
		"The type string cannot be templatized")
}

// Java: TypeTransformationTest.testTransformationWithTemplatizedTypeInvalidBaseType2()
func (c *TypeTransformationTest) TestTransformationWithTemplatizedTypeInvalidBaseType2() {
	c.testTTL(c.GetNativeUnknownType(), "type(S, 'number')", "The type string cannot be templatized")
}

// Java: TypeTransformationTest.testTransformationWithTemplatizedTypeInvalidBaseType3()
func (c *TypeTransformationTest) TestTransformationWithTemplatizedTypeInvalidBaseType3() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"type('ArrayAlias', number)",
		"The type Array<?> cannot be templatized")
}

// Java: TypeTransformationTest.testTransformationWithRawTypeOf()
func (c *TypeTransformationTest) TestTransformationWithRawTypeOf() {
	c.testTTL(c.GetNativeArrayType(), "rawTypeOf(type('Array', 'number'))")
}

// Java: TypeTransformationTest.testTransformationWithRawTypeOf2()
func (c *TypeTransformationTest) TestTransformationWithRawTypeOf2() {
	c.testTTL(c.GetNativeArrayType(), "rawTypeOf(ARRNUM)")
}

// Java: TypeTransformationTest.testTransformationWithNestedRawTypeOf()
func (c *TypeTransformationTest) TestTransformationWithNestedRawTypeOf() {
	c.testTTL(c.GetNativeArrayType(), "rawTypeOf(type('Array', rawTypeOf(ARRNUM)))")
}

// Java: TypeTransformationTest.testTransformationWithInvalidRawTypeOf()
func (c *TypeTransformationTest) TestTransformationWithInvalidRawTypeOf() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"rawTypeOf(N)",
		"Expected templatized type in rawTypeOf found number")
}

// Java: TypeTransformationTest.testTransformationWithTemplateTypeOf()
func (c *TypeTransformationTest) TestTransformationWithTemplateTypeOf() {
	c.testTTL(c.GetNativeNumberType(), "templateTypeOf(type('Array', 'number'), 0)")
}

// Java: TypeTransformationTest.testTransformationWithTemplateTypeOf2()
func (c *TypeTransformationTest) TestTransformationWithTemplateTypeOf2() {
	c.testTTL(c.GetNativeNumberType(), "templateTypeOf(ARRNUM, 0)")
}

// Java: TypeTransformationTest.testTransformationWithNestedTemplateTypeOf()
func (c *TypeTransformationTest) TestTransformationWithNestedTemplateTypeOf() {
	c.testTTL(
		c.GetNativeNumberType(),
		"templateTypeOf(templateTypeOf(type('Array', type('Array', 'number')), 0), 0)")
}

// Java: TypeTransformationTest.testTransformationWithInvalidTypeTemplateTypeOf()
func (c *TypeTransformationTest) TestTransformationWithInvalidTypeTemplateTypeOf() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"templateTypeOf(N, 0)",
		"Expected templatized type in templateTypeOf found number")
}

// Java: TypeTransformationTest.testTransformationWithInvalidIndexTemplateTypeOf()
func (c *TypeTransformationTest) TestTransformationWithInvalidIndexTemplateTypeOf() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"templateTypeOf(ARRNUM, 2)",
		"Index out of bounds in templateTypeOf: expected a number less than 1, found 2")
}

// Java: TypeTransformationTest.testTransformationWithRecordType()
func (c *TypeTransformationTest) TestTransformationWithRecordType() {
	c.testTTL(c.recordWithProperty("x", c.GetNativeNumberType()), "record({x:'number'})")
}

// Java: TypeTransformationTest.testTransformationWithRecordType2()
func (c *TypeTransformationTest) TestTransformationWithRecordType2() {
	c.testTTL(c.recordWithProperty("0", c.GetNativeNumberType()), "record({0:'number'})")
}

// Java: TypeTransformationTest.testTransformationWithRecordTypeMultipleProperties()
func (c *TypeTransformationTest) TestTransformationWithRecordTypeMultipleProperties() {
	c.testTTL(
		c.recordWithProperties2("x", c.GetNativeNumberType(), "y", c.GetNativeStringType()),
		"record({x:'number', y:S})")
}

// Java: TypeTransformationTest.testTransformationWithNestedRecordType()
func (c *TypeTransformationTest) TestTransformationWithNestedRecordType() {
	c.testTTL(
		c.recordWithProperties2("x", c.recordWithProperty("z", c.GetNativeBooleanType()), "y", c.GetNativeStringType()),
		"record({x:record({z:B}), y:S})")
}

// Java: TypeTransformationTest.testTransformationWithNestedRecordType2()
func (c *TypeTransformationTest) TestTransformationWithNestedRecordType2() {
	c.testTTL(c.recordWithProperty("x", c.GetNativeNumberType()), "record(record({x:N}))")
}

// Java: TypeTransformationTest.testTransformationWithEmptyRecordType()
func (c *TypeTransformationTest) TestTransformationWithEmptyRecordType() {
	c.testTTL(c.record(), "record({})")
}

// Java: TypeTransformationTest.testTransformationWithMergeRecord()
func (c *TypeTransformationTest) TestTransformationWithMergeRecord() {
	c.testTTL(
		c.recordWithProperties3("x", c.GetNativeNumberType(), "y", c.GetNativeStringType(), "z", c.GetNativeBooleanType()),
		"record({x:N}, {y:S}, {z:B})")
}

// Java: TypeTransformationTest.testTransformationWithMergeDuplicatedRecords()
func (c *TypeTransformationTest) TestTransformationWithMergeDuplicatedRecords() {
	c.testTTL(c.recordWithProperty("x", c.GetNativeNumberType()), "record({x:N}, {x:N}, {x:N})")
}

// Java: TypeTransformationTest.testTransformationWithMergeRecordTypeWithEmpty()
func (c *TypeTransformationTest) TestTransformationWithMergeRecordTypeWithEmpty() {
	c.testTTL(c.recordWithProperty("x", c.GetNativeNumberType()), "record({x:N}, {})")
}

// Java: TypeTransformationTest.testTransformationWithInvalidRecordType()
func (c *TypeTransformationTest) TestTransformationWithInvalidRecordType() {
	c.testTTL(c.GetNativeUnknownType(), "record(N)", "Expected a record type, found number")
}

// Java: TypeTransformationTest.testTransformationWithInvalidMergeRecordType()
func (c *TypeTransformationTest) TestTransformationWithInvalidMergeRecordType() {
	c.testTTL(c.GetNativeUnknownType(), "record({x:N}, N)", "Expected a record type, found number")
}

// Java: TypeTransformationTest.testTransformationWithTTLTypeTransformationInFirstParamMapunion()
func (c *TypeTransformationTest) TestTransformationWithTTLTypeTransformationInFirstParamMapunion() {
	c.testTTL(
		c.union(c.GetNativeNumberType(), c.GetNativeStringType()),
		"mapunion(templateTypeOf(type(ARR, union(N, S)), 0), (x) => x)")
}

// Java: TypeTransformationTest.testTransformationWithInvalidNestedMapunion()
func (c *TypeTransformationTest) TestTransformationWithInvalidNestedMapunion() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"mapunion(union(S, B),\n"+
			"(x) => mapunion(union(S, N),\n"+
			"(x) => cond(eq(x, x), x, BOT)))\n",
		"The variable x is already defined")
}

// Java: TypeTransformationTest.testTransformationWithTTLRecordWithReference()
func (c *TypeTransformationTest) TestTransformationWithTTLRecordWithReference() {
	c.testTTL(
		c.recordWithProperties3(
			"number",
			c.GetNativeNumberType(),
			"string",
			c.GetNativeStringType(),
			"boolean",
			c.GetNativeBooleanType()),
		"record({[n]:N, [s]:S, [b]:B})")
}

// Java: TypeTransformationTest.testTransformationWithTTLRecordWithInvalidReference()
func (c *TypeTransformationTest) TestTransformationWithTTLRecordWithInvalidReference() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"record({[Foo]:N})",
		"Expected a record type, found ?",
		"Reference to an unknown name variable Foo")
}

// Java: TypeTransformationTest.testTransformationWithMaprecordMappingEverythingToString()
func (c *TypeTransformationTest) TestTransformationWithMaprecordMappingEverythingToString() {
	// {n:number, s:string, b:boolean}
	// is transformed to
	// {n:string, s:string, b:string}
	c.testTTL(
		c.recordWithProperties3("n", c.GetNativeStringType(), "s", c.GetNativeStringType(), "b", c.GetNativeStringType()),
		"maprecord(REC, (k, v) => record({[k]:S}))")
}

// Java: TypeTransformationTest.testTransformationWithMaprecordIdentity()
func (c *TypeTransformationTest) TestTransformationWithMaprecordIdentity() {
	// {n:number, s:string, b:boolean} remains the same
	c.testTTL(typeTransformationTest_recordTypeTest, "maprecord(REC, (k, v) => record({[k]:v}))")
}

// Java: TypeTransformationTest.testTransformationWithMaprecordDeleteEverything()
func (c *TypeTransformationTest) TestTransformationWithMaprecordDeleteEverything() {
	// {n:number, s:string, b:boolean}
	// is transformed to object type
	c.testTTL(c.GetNativeObjectType(), "maprecord(REC, (k, v) => BOT)")
}

// Java: TypeTransformationTest.testTransformationWithInvalidMaprecord()
func (c *TypeTransformationTest) TestTransformationWithInvalidMaprecord() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"maprecord(REC, (k, v) => 'number')",
		"The body of a maprecord function must evaluate to a record type or a no type, found number")
}

// Java: TypeTransformationTest.testTransformationWithMaprecordFilterWithOnlyString()
func (c *TypeTransformationTest) TestTransformationWithMaprecordFilterWithOnlyString() {
	// {n:number, s:string, b:boolean}
	// is transformed to
	// {s:string}
	c.testTTL(
		c.recordWithProperty("s", c.GetNativeStringType()),
		"maprecord(REC, (k, v) => cond(eq(v, S), record({[k]:v}), BOT))")
}

// Java: TypeTransformationTest.testTransformationWithInvalidMaprecordFirstParam()
func (c *TypeTransformationTest) TestTransformationWithInvalidMaprecordFirstParam() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"maprecord(N, (k, v) => BOT)",
		"The first parameter of a maprecord must be a record type, found number")
}

// Java: TypeTransformationTest.testTransformationWithObjectInMaprecord()
func (c *TypeTransformationTest) TestTransformationWithObjectInMaprecord() {
	c.testTTL(c.GetNativeObjectType(), "maprecord(OBJ, (k, v) => record({[k]:v}))")
}

// Java: TypeTransformationTest.testTransformationWithUnionInMaprecord()
func (c *TypeTransformationTest) TestTransformationWithUnionInMaprecord() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"maprecord(union(record({n:N}), S), (k, v) => record({[k]:v}))",
		"The first parameter of a maprecord must be a record type, found (string|{n: number})")
}

// Java: TypeTransformationTest.testTransformationWithUnionOfRecordsInMaprecord()
func (c *TypeTransformationTest) TestTransformationWithUnionOfRecordsInMaprecord() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"maprecord(union(record({n:N}), record({s:S})),  (k, v) => record({[k]:v}))",
		"The first parameter of a maprecord must be a record type, found ({n: number}|{s: string})")
}

// Java: TypeTransformationTest.testTransformationWithNestedRecordInMaprecordFilterOneLevelString()
func (c *TypeTransformationTest) TestTransformationWithNestedRecordInMaprecordFilterOneLevelString() {
	// {s:string, r:{s:string, b:boolean}}
	// is transformed to
	// {r:{s:string, b:boolean}}
	c.testTTL(
		c.recordWithProperty("r", c.recordWithProperties2("s", c.GetNativeStringType(), "b", c.GetNativeBooleanType())),
		"maprecord(NESTEDREC, (k, v) => cond(eq(v, S), BOT, record({[k]:v})))")
}

// Java: TypeTransformationTest.testTransformationWithNestedRecordInMaprecordFilterTwoLevelsString()
func (c *TypeTransformationTest) TestTransformationWithNestedRecordInMaprecordFilterTwoLevelsString() {
	// {s:string, r:{s:string, b:boolean}}
	// is transformed to
	// {r:{b:boolean}}
	c.testTTL(
		c.recordWithProperty("r", c.recordWithProperty("b", c.GetNativeBooleanType())),
		"maprecord(NESTEDREC,\n"+
			"    (k1, v1) =>\n"+
			"        cond(sub(v1, 'Object'),\n"+
			"            maprecord(v1, (k2, v2) =>\n"+
			"                cond(eq(v2, S), BOT, record({[k1]:record({[k2]:v2})}))),\n"+
			"            cond(eq(v1, S), BOT, record({[k1]:v1}))))\n")
}

// Java: TypeTransformationTest.testTransformationWithNestedIdentityOneLevel()
func (c *TypeTransformationTest) TestTransformationWithNestedIdentityOneLevel() {
	// {r:{n:number, s:string}}
	c.testTTL(
		c.recordWithProperty("r", c.recordWithProperties2("b", c.GetNativeBooleanType(), "s", c.GetNativeStringType())),
		"maprecord(record({r:record({b:B, s:S})}),\n"+
			"    (k1, v1) =>\n"+
			"        maprecord(v1,\n"+
			"            (k2, v2) =>\n"+
			"                record({[k1]:record({[k2]:v2})})))\n")
}

// Java: TypeTransformationTest.testTransformationWithNestedIdentityOneLevel2()
func (c *TypeTransformationTest) TestTransformationWithNestedIdentityOneLevel2() {
	// {r:{n:number, s:string}}
	c.testTTL(
		c.recordWithProperty("r", c.recordWithProperties2("b", c.GetNativeBooleanType(), "s", c.GetNativeStringType())),
		"maprecord(record({r:record({b:B, s:S})}),\n"+
			"    (k1, v1) =>\n"+
			"        record({[k1]:\n"+
			"            maprecord(v1,\n"+
			"                (k2, v2) => record({[k2]:v2}))}))\n")
}

// Java: TypeTransformationTest.testTransformationWithNestedIdentityTwoLevels()
func (c *TypeTransformationTest) TestTransformationWithNestedIdentityTwoLevels() {
	// {r:{r2:{n:number, s:string}}}
	c.testTTL(
		c.recordWithProperty("r", c.recordWithProperty("r2", c.recordWithProperties2("n", c.GetNativeNumberType(), "s", c.GetNativeStringType()))),
		"maprecord(record({r:record({r2:record({n:N, s:S})})}),\n"+
			"    (k1, v1) =>\n"+
			"        maprecord(v1,\n"+
			"            (k2, v2) =>\n"+
			"                maprecord(v2,\n"+
			"                    (k3, v3) =>\n"+
			"                        record({[k1]:\n"+
			"                            record({[k2]:\n"+
			"                                record({[k3]:v3})})}))))\n")
}

// Java: TypeTransformationTest.testTransformationWithNestedIdentityTwoLevels2()
func (c *TypeTransformationTest) TestTransformationWithNestedIdentityTwoLevels2() {
	// {r:{r2:{n:number, s:string}}}
	c.testTTL(
		c.recordWithProperty("r", c.recordWithProperty("r2", c.recordWithProperties2("n", c.GetNativeNumberType(), "s", c.GetNativeStringType()))),
		"maprecord(record({r:record({r2:record({n:N, s:S})})}),\n"+
			"    (k1, v1) =>\n"+
			"        record({[k1]:\n"+
			"            maprecord(v1,\n"+
			"                (k2, v2) =>\n"+
			"                    record({[k2]:\n"+
			"                        maprecord(v2,\n"+
			"                            (k3, v3) =>\n"+
			"                                record({[k3]:v3}))}))}))\n")
}

// Java: TypeTransformationTest.testTransformationWithNestedIdentityThreeLevels()
func (c *TypeTransformationTest) TestTransformationWithNestedIdentityThreeLevels() {
	// {r:{r2:{r3:{n:number, s:string}}}}
	c.testTTL(
		c.recordWithProperty(
			"r",
			c.recordWithProperty(
				"r2",
				c.recordWithProperty("r3", c.recordWithProperties2("n", c.GetNativeNumberType(), "s", c.GetNativeStringType())))),
		"maprecord(record({r:record({r2:record({r3:record({n:N, s:S})})})}),\n"+
			"    (k1, v1) =>\n"+
			"        maprecord(v1,\n"+
			"            (k2, v2) =>\n"+
			"                maprecord(v2,\n"+
			"                    (k3, v3) =>\n"+
			"                        maprecord(v3,\n"+
			"                            (k4, v4) =>\n"+
			"                                record({[k1]:\n"+
			"                                    record({[k2]:\n"+
			"                                        record({[k3]:\n"+
			"                                            record({[k4]:v4})})})})))))\n")
}

// Java: TypeTransformationTest.testTransformationWithNestedIdentityThreeLevels2()
func (c *TypeTransformationTest) TestTransformationWithNestedIdentityThreeLevels2() {
	// {r:{r2:{r3:{n:number, s:string}}}}
	c.testTTL(
		c.recordWithProperty(
			"r",
			c.recordWithProperty(
				"r2",
				c.recordWithProperty("r3", c.recordWithProperties2("n", c.GetNativeNumberType(), "s", c.GetNativeStringType())))),
		"maprecord(record({r:record({r2:record({r3:record({n:N, s:S})})})}),\n"+
			"    (k1, v1) =>\n"+
			"        record({[k1]:\n"+
			"            maprecord(v1,\n"+
			"                (k2, v2) =>\n"+
			"                    record({[k2]:\n"+
			"                        maprecord(v2,\n"+
			"                            (k3, v3) =>\n"+
			"                                record({[k3]:\n"+
			"                                    maprecord(v3,\n"+
			"                                        (k4, v4) =>\n"+
			"                                            record({[k4]:v4}))}))}))}))\n")
}

// Java: TypeTransformationTest.testTransformationWithNestedRecordDeleteLevelTwoAndThree()
func (c *TypeTransformationTest) TestTransformationWithNestedRecordDeleteLevelTwoAndThree() {
	// {r:{r2:{r3:{n:number, s:string}}}}
	// is transformed into
	// {r:{n:number, s:string}}
	c.testTTL(
		c.recordWithProperty("r", c.recordWithProperties2("n", c.GetNativeNumberType(), "s", c.GetNativeStringType())),
		"maprecord(record({r:record({r2:record({r3:record({n:N, s:S})})})}),\n"+
			"    (k1, v1) =>\n"+
			"        maprecord(v1,\n"+
			"            (k2, v2) =>\n"+
			"                maprecord(v2,\n"+
			"                    (k3, v3) =>\n"+
			"                        maprecord(v3,\n"+
			"                            (k4, v4) =>\n"+
			"                                record({[k1]:\n"+
			"                                    record({[k4]:v4})})))))\n")
}

// Java: TypeTransformationTest.testTransformationWithNestedRecordDeleteLevelTwoAndThree2()
func (c *TypeTransformationTest) TestTransformationWithNestedRecordDeleteLevelTwoAndThree2() {
	// {r:{r2:{r3:{n:number, s:string}}}}
	// is transformed into
	// {r:{n:number, s:string}}
	c.testTTL(
		c.recordWithProperty("r", c.recordWithProperties2("n", c.GetNativeNumberType(), "s", c.GetNativeStringType())),
		"maprecord(record({r:record({r2:record({r3:record({n:N, s:S})})})}),\n"+
			"    (k1, v1) =>\n"+
			"        record({[k1]:\n"+
			"            maprecord(v1,\n"+
			"                (k2, v2) =>\n"+
			"                    maprecord(v2,\n"+
			"                        (k3, v3) =>\n"+
			"                            maprecord(v3,\n"+
			"                                (k4, v4) =>\n"+
			"                                    record({[k4]:v4}))))}))\n")
}

// Java: TypeTransformationTest.testTransformationWithNestedRecordCollapsePropertiesToRecord()
func (c *TypeTransformationTest) TestTransformationWithNestedRecordCollapsePropertiesToRecord() {
	// {a:Array, b:{n:number}}
	// is transformed to
	// {foo:{n:number}}
	c.testTTL(
		c.recordWithProperty("foo", c.recordWithProperty("n", c.GetNativeNumberType())),
		"maprecord(record({a:ARR, b:record({n:N})}), (k, v) => record({foo:v}))")
}

// Java: TypeTransformationTest.testTransformationWithNestedRecordCollapsePropertiesToType()
func (c *TypeTransformationTest) TestTransformationWithNestedRecordCollapsePropertiesToType() {
	// {a:{n:number}, b:Array}
	// is transformed to
	// {foo:Array}
	c.testTTL(
		c.recordWithProperty("foo", c.GetNativeArrayType()),
		"maprecord(record({a:record({n:N}), b:ARR}), (k, v) => record({foo:v}))")
}

// Java: TypeTransformationTest.testTransformationWithNestedRecordCollapsePropertiesJoinRecords()
func (c *TypeTransformationTest) TestTransformationWithNestedRecordCollapsePropertiesJoinRecords() {
	// {a:{n:number}, b:{s:Array}}
	// is transformed to
	// {foo:{n:number, s:Array}}
	c.testTTL(
		c.recordWithProperty("foo", c.recordWithProperties2("n", c.GetNativeNumberType(), "s", c.GetNativeArrayType())),
		"maprecord(record({a:record({n:N}), b:record({s:ARR})}), (k, v) => record({foo:v}))")
}

// Java: TypeTransformationTest.testTransformationWithNestedRecordCollapsePropertiesJoinRecords2()
func (c *TypeTransformationTest) TestTransformationWithNestedRecordCollapsePropertiesJoinRecords2() {
	// {a:{n:number, {x:number}}, b:{s:Array, {y:number}}}
	// is transformed to
	// {foo:{n:number, s:Array, r:{x:number, y:number}}}
	c.testTTL(
		c.recordWithProperty(
			"foo",
			c.recordWithProperties3(
				"n",
				c.GetNativeNumberType(),
				"s",
				c.GetNativeArrayType(),
				"r",
				c.recordWithProperties2("x", c.GetNativeNumberType(), "y", c.GetNativeNumberType()))),
		"maprecord(\n"+
			"    record({a:record({n:N, r:record({x:N})}), b:record({s:ARR, r:record({y:N})})}),\n"+
			"    (k, v) => record({foo:v}))\n")
}

// Java: TypeTransformationTest.testTransformationWithAsynchUseCase()
func (c *TypeTransformationTest) TestTransformationWithAsynchUseCase() {
	// TODO(lpino): Use the type Promise instead of Array
	// {service:Array<number>}
	// is transformed to
	// {service:number}
	c.testTTL(
		c.recordWithProperty("service", c.GetNativeNumberType()),
		"cond(\n"+
			"    sub(ASYNCH, 'Object'),\n"+
			"    maprecord(\n"+
			"        ASYNCH,\n"+
			"        (k, v) =>\n"+
			"            cond(\n"+
			"                eq(rawTypeOf(v), 'Array'),\n"+
			"                record({[k]:templateTypeOf(v, 0)}),\n"+
			"                record({[k]:'undefined'}))),\n"+
			"    ASYNCH)\n")
}

// Java: TypeTransformationTest.testTransformationWithInvalidNestedMaprecord()
func (c *TypeTransformationTest) TestTransformationWithInvalidNestedMaprecord() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"maprecord(NESTEDREC, (k, v) => maprecord(v, (k, v) => BOT))",
		"The variable k is already defined")
}

// Java: TypeTransformationTest.testTransformationWithMaprecordAndStringEquivalence()
func (c *TypeTransformationTest) TestTransformationWithMaprecordAndStringEquivalence() {
	c.testTTL(
		c.recordWithProperties2("bool", c.GetNativeNumberType(), "str", c.GetNativeStringType()),
		"maprecord(record({bool:B, str:S}),\n"+
			"(k, v) => record({[k]:cond(streq(k, 'bool'), N, v)}))\n")
}

// Java: TypeTransformationTest.testTransformationWithTypeOfVar()
func (c *TypeTransformationTest) TestTransformationWithTypeOfVar() {
	c.testTTL(c.GetNativeNumberType(), "typeOfVar('n')")
}

// Java: TypeTransformationTest.testTransformationWithUnknownTypeOfVar()
func (c *TypeTransformationTest) TestTransformationWithUnknownTypeOfVar() {
	c.testTTL(c.GetNativeUnknownType(), "typeOfVar('foo')", "Variable foo is undefined in the scope")
}

// Java: TypeTransformationTest.testTransformationWithTrueIsConstructorConditional()
func (c *TypeTransformationTest) TestTransformationWithTrueIsConstructorConditional() {
	c.testTTL(c.GetNativeStringType(), "cond(isCtor(typeOfVar('Bar')), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseIsConstructorConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseIsConstructorConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(isCtor(N), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueIsTemplatizedConditional()
func (c *TypeTransformationTest) TestTransformationWithTrueIsTemplatizedConditional() {
	c.testTTL(c.GetNativeStringType(), "cond(isTemplatized(type(ARR, N)), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseIsTemplatizedConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseIsTemplatizedConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(isTemplatized(ARR), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueIsRecordConditional()
func (c *TypeTransformationTest) TestTransformationWithTrueIsRecordConditional() {
	c.testTTL(c.GetNativeStringType(), "cond(isRecord(REC), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseIsRecordConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseIsRecordConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(isRecord(N), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueIsDefinedConditional()
func (c *TypeTransformationTest) TestTransformationWithTrueIsDefinedConditional() {
	c.testTTL(c.GetNativeStringType(), "cond(isDefined(N), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseIsDefinedConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseIsDefinedConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(isDefined(Foo), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueIsUnknownConditional()
func (c *TypeTransformationTest) TestTransformationWithTrueIsUnknownConditional() {
	c.testTTL(c.GetNativeStringType(), "cond(isUnknown(UNK), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueIsUnknownConditional2()
func (c *TypeTransformationTest) TestTransformationWithTrueIsUnknownConditional2() {
	c.testTTL(c.GetNativeStringType(), "cond(isUnknown(CHKUNK), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseIsUnknownConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseIsUnknownConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(isUnknown(N), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueAndConditional()
func (c *TypeTransformationTest) TestTransformationWithTrueAndConditional() {
	c.testTTL(c.GetNativeStringType(), "cond(isDefined(N) && isDefined(N), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseAndConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseAndConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(isDefined(N) && isDefined(Foo), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseAndConditional2()
func (c *TypeTransformationTest) TestTransformationWithFalseAndConditional2() {
	c.testTTL(c.GetNativeNumberType(), "cond(isDefined(Foo) && isDefined(N), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseAndConditional3()
func (c *TypeTransformationTest) TestTransformationWithFalseAndConditional3() {
	c.testTTL(c.GetNativeNumberType(), "cond(isDefined(Foo) && isDefined(Foo), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueOrConditional()
func (c *TypeTransformationTest) TestTransformationWithTrueOrConditional() {
	c.testTTL(c.GetNativeStringType(), "cond(isDefined(N) || isDefined(N), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueOrConditional2()
func (c *TypeTransformationTest) TestTransformationWithTrueOrConditional2() {
	c.testTTL(c.GetNativeStringType(), "cond(isDefined(Foo) || isDefined(N), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueOrConditional3()
func (c *TypeTransformationTest) TestTransformationWithTrueOrConditional3() {
	c.testTTL(c.GetNativeStringType(), "cond(isDefined(N) || isDefined(Foo), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseOrConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseOrConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(isDefined(Foo) || isDefined(Foo), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithTrueNotConditional3()
func (c *TypeTransformationTest) TestTransformationWithTrueNotConditional3() {
	c.testTTL(c.GetNativeStringType(), "cond(!isDefined(Foo), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithFalseNotConditional()
func (c *TypeTransformationTest) TestTransformationWithFalseNotConditional() {
	c.testTTL(c.GetNativeNumberType(), "cond(!isDefined(N), 'string', 'number')")
}

// Java: TypeTransformationTest.testTransformationWithInstanceOf()
func (c *TypeTransformationTest) TestTransformationWithInstanceOf() {
	c.testTTL(c.GetNativeNumberObjectType(), "instanceOf(typeOfVar('Number'))")
}

// Java: TypeTransformationTest.testTransformationWithInvalidInstanceOf()
func (c *TypeTransformationTest) TestTransformationWithInvalidInstanceOf() {
	c.testTTL(c.GetNativeUnknownType(), "instanceOf(N)", "Expected a constructor type, found number")
}

// Java: TypeTransformationTest.testTransformationWithInvalidInstanceOf2()
func (c *TypeTransformationTest) TestTransformationWithInvalidInstanceOf2() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"instanceOf(foo)",
		"Expected a constructor type, found Unknown",
		"Reference to an unknown type variable foo")
}

// Java: TypeTransformationTest.testTransformationWithTypeExpr()
func (c *TypeTransformationTest) TestTransformationWithTypeExpr() {
	c.testTTL(c.GetNativeNumberType(), "typeExpr('number')")
}

// Java: TypeTransformationTest.testParserWithTTLNativeTypeExprUnion()
func (c *TypeTransformationTest) TestParserWithTTLNativeTypeExprUnion() {
	c.testTTL(c.union(c.GetNativeNumberType(), c.GetNativeBooleanType()), "typeExpr('number|boolean')")
}

// Java: TypeTransformationTest.testParserWithTTLNativeTypeExprRecord()
func (c *TypeTransformationTest) TestParserWithTTLNativeTypeExprRecord() {
	c.testTTL(
		c.recordWithProperties2("foo", c.GetNativeNumberType(), "bar", c.GetNativeBooleanType()),
		"typeExpr('{foo:number, bar:boolean}')")
}

// Java: TypeTransformationTest.testParserWithTTLNativeTypeExprNullable()
func (c *TypeTransformationTest) TestParserWithTTLNativeTypeExprNullable() {
	c.testTTL(c.union(c.GetNativeNumberType(), c.GetNativeNullType()), "typeExpr('?number')")
}

// Java: TypeTransformationTest.testParserWithTTLNativeTypeExprNonNullable()
func (c *TypeTransformationTest) TestParserWithTTLNativeTypeExprNonNullable() {
	c.testTTL(c.GetNativeNumberType(), "typeExpr('!number')")
}

// Java: TypeTransformationTest.testTransformationPrintType()
func (c *TypeTransformationTest) TestTransformationPrintType() {
	c.testTTL(c.GetNativeNumberType(), "printType('Test message: ', N)")
}

// Java: TypeTransformationTest.testTransformationPrintType2()
func (c *TypeTransformationTest) TestTransformationPrintType2() {
	c.testTTL(typeTransformationTest_recordTypeTest, "printType('Test message: ', REC)")
}

// Java: TypeTransformationTest.testTransformationPropType()
func (c *TypeTransformationTest) TestTransformationPropType() {
	c.testTTL(c.GetNativeNumberType(), "propType('a', record({a:N, b:record({x:B})}))")
}

// Java: TypeTransformationTest.testTransformationPropType2()
func (c *TypeTransformationTest) TestTransformationPropType2() {
	c.testTTL(c.recordWithProperty("x", c.GetNativeBooleanType()), "propType('b', record({a:N, b:record({x:B})}))")
}

// Java: TypeTransformationTest.testTransformationPropTypeNotFound()
func (c *TypeTransformationTest) TestTransformationPropTypeNotFound() {
	c.testTTL(c.GetNativeUnknownType(), "propType('c', record({a:N, b:record({x:B})}))")
}

// Java: TypeTransformationTest.testTransformationPropTypeInvalid()
func (c *TypeTransformationTest) TestTransformationPropTypeInvalid() {
	c.testTTL(c.GetNativeUnknownType(), "propType('c', N)", "Expected object type, found number")
}

// Java: TypeTransformationTest.testTransformationInstanceObjectToRecord()
func (c *TypeTransformationTest) TestTransformationInstanceObjectToRecord() {
	c.testTTL(c.GetNativeObjectType(), "record(type(OBJ, N))")
}

// Java: TypeTransformationTest.testTransformationInstanceObjectToRecord2()
func (c *TypeTransformationTest) TestTransformationInstanceObjectToRecord2() {
	// TODO(bradfordcsmith): Define Array.prototype.length using externs instead.
	c.GetNativeArrayType().
		DefineDeclaredProperty("length", c.GetNativeNumberType() /* propertyNode= */, nil)
	c.testTTL(c.recordWithProperty("length", c.GetNativeNumberType()), "record(type(ARR, N))")
}

// Java: TypeTransformationTest.testTransformationInstanceObjectToRecordInvalid()
func (c *TypeTransformationTest) TestTransformationInstanceObjectToRecordInvalid() {
	c.testTTL(
		c.GetNativeUnknownType(),
		"record(union(OBJ, NULL))",
		"Expected a record type, found (Object|null)")
}

// Java: TypeTransformationTest.initRecordTypeTests()
func (c *TypeTransformationTest) initRecordTypeTests() {
	// {n:number, s:string, b:boolean}
	typeTransformationTest_recordTypeTest =
		c.recordWithProperties3("n", c.GetNativeNumberType(), "s", c.GetNativeStringType(), "b", c.GetNativeBooleanType())
	// {n:number, r:{s:string, b:boolean}}
	typeTransformationTest_nestedRecordTypeTest =
		c.recordWithProperties2(
			"s",
			c.GetNativeStringType(),
			"r",
			c.recordWithProperties2("s", c.GetNativeStringType(), "b", c.GetNativeBooleanType()))
	// {service:Array<number>}
	typeTransformationTest_asynchRecord = c.recordWithProperty("service", c.type_(c.GetNativeArrayType(), c.GetNativeNumberType()))
}

// Java: TypeTransformationTest.union(JSType...)
func (c *TypeTransformationTest) union(variants ...rhino.JSType) rhino.JSType {
	type_ := c.CreateUnionType(variants...)
	AssertThat(type_.IsUnionType()).IsTrue()
	return type_
}

// Java: TypeTransformationTest.type(ObjectType, JSType...)
func (c *TypeTransformationTest) type_(baseType rhino.ObjectType, templatizedTypes ...rhino.JSType) rhino.JSType {
	return c.CreateTemplatizedType(baseType, templatizedTypes...)
}

// Java: TypeTransformationTest.record()
func (c *TypeTransformationTest) record() rhino.JSType {
	return c.recordFromMap(jutil.NewLinkedHashMap[string, rhino.JSType]())
}

// Java: TypeTransformationTest.record(String, JSType)
func (c *TypeTransformationTest) recordWithProperty(p1 string, t1 rhino.JSType) rhino.JSType {
	return c.recordFromMap(jutil.NewImmutableMapBuilder[string, rhino.JSType]().Put(p1, t1).BuildOrThrow())
}

// Java: TypeTransformationTest.record(String, JSType, String, JSType)
func (c *TypeTransformationTest) recordWithProperties2(p1 string, t1 rhino.JSType, p2 string, t2 rhino.JSType) rhino.JSType {
	return c.recordFromMap(jutil.NewImmutableMapBuilder[string, rhino.JSType]().Put(p1, t1).Put(p2, t2).BuildOrThrow())
}

// Java: TypeTransformationTest.record(String, JSType, String, JSType, String, JSType)
func (c *TypeTransformationTest) recordWithProperties3(p1 string, t1 rhino.JSType, p2 string, t2 rhino.JSType, p3 string, t3 rhino.JSType) rhino.JSType {
	return c.recordFromMap(
		jutil.NewImmutableMapBuilder[string, rhino.JSType]().Put(p1, t1).Put(p2, t2).Put(p3, t3).BuildOrThrow())
}

// Java: TypeTransformationTest.record(ImmutableMap<String, JSType>)
func (c *TypeTransformationTest) recordFromMap(props *jutil.LinkedHashMap[string, rhino.JSType]) rhino.JSType {
	builder := c.CreateRecordTypeBuilder()
	for key, value := range props.All() {
		builder.AddProperty(key, value, nil)
	}
	return builder.Build()
}

// Java: TypeTransformationTest.testTTL(JSType, String, String...)
func (c *TypeTransformationTest) testTTL(expectedType rhino.JSType, ttlExp string, expectedWarnings ...string) {
	closer := c.compiler.GetTypeRegistry().GetResolver().OpenForDefinition()
	defer closer.Close()
	ttlParser :=
		parsing.NewTypeTransformationParser(
			ttlExp, SourceFile_FromCode("[testcode]", ttlExp), c.errorReporter, 0, 0)
	// Run the test if the parsing was successful
	if ttlParser.ParseTypeTransformation() {
		ast := ttlParser.GetTypeTransformationAst()
		// Create the scope using the extra definitions
		extraTypeDefs := c.compiler.ParseTestCode(TypeTransformationTest_EXTRA_TYPE_DEFS)
		// TypedScopeCreator.createScope returns a TypedScope (asserted through any so that
		// this compiles whatever the declared Go result type is).
		var scopeAny any = NewTypedScopeCreator(c.compiler).
			CreateScope(rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(extraTypeDefs)), nil)
		scope := scopeAny.(*TypedScope)
		// Evaluate the type transformation
		typeTransformation := NewTypeTransformation(c.compiler, scope)
		resultType := typeTransformation.EvalWithNameVars(ast, c.typeVars, c.nameVars)
		c.CheckReportedWarningsHelper(expectedWarnings)
		c.AssertTypeEquals(expectedType, resultType)
	}
}

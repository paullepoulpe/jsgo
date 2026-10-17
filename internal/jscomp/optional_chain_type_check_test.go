// Port of com/google/javascript/jscomp/OptionalChainTypeCheckTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/OptionalChainTypeCheckTest.java)

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Tests types for Optional chain nodes.
//
// The Java class runs with JUnit's Enclosed runner: each nested class is a
// test class of its own, run by a Go test function of its own. The
// Parameterized ones run their test methods once per element of their
// cases(), as subtests named "[<index>]" like JUnit's default parameter
// names (optionalChainTypeCheckTest_runParameterized).

// OptionalChainTypeCheckTest_OptChainTestCase is the AutoValue
// OptChainTestCase: the generated AutoValue class is folded into it.
//
// Go-only: language-adapter — an Optional<String> property is a *string
// (nil is Optional.empty()); get() on an empty one dereferences nil.
type OptionalChainTypeCheckTest_OptChainTestCase struct {
	forObjType         *string
	withExpr           string
	assignedTo         *string
	mustReport         *string
	withPropReturnType *string
}

// Java: OptionalChainTypeCheckTest.OptChainTestCase.forObjType()
func (c *OptionalChainTypeCheckTest_OptChainTestCase) ForObjType() *string {
	return c.forObjType
}

// Java: OptionalChainTypeCheckTest.OptChainTestCase.withExpr()
func (c *OptionalChainTypeCheckTest_OptChainTestCase) WithExpr() string {
	return c.withExpr
}

// Java: OptionalChainTypeCheckTest.OptChainTestCase.assignedTo()
func (c *OptionalChainTypeCheckTest_OptChainTestCase) AssignedTo() *string {
	return c.assignedTo
}

// Java: OptionalChainTypeCheckTest.OptChainTestCase.mustReport()
func (c *OptionalChainTypeCheckTest_OptChainTestCase) MustReport() *string {
	return c.mustReport
}

// WithPropReturnType is only used in testing CALL.
//
// Java: OptionalChainTypeCheckTest.OptChainTestCase.withPropReturnType()
func (c *OptionalChainTypeCheckTest_OptChainTestCase) WithPropReturnType() *string {
	return c.withPropReturnType
}

// Java: OptionalChainTypeCheckTest.OptChainTestCase.builder()
func OptionalChainTypeCheckTest_OptChainTestCase_Builder_() *OptionalChainTypeCheckTest_OptChainTestCase_Builder {
	return &OptionalChainTypeCheckTest_OptChainTestCase_Builder{}
}

// OptionalChainTypeCheckTest_OptChainTestCase_Builder is the
// @AutoValue.Builder of OptChainTestCase, with the generated
// AutoValue_OptionalChainTypeCheckTest_OptChainTestCase.Builder folded in.
type OptionalChainTypeCheckTest_OptChainTestCase_Builder struct {
	forObjType         *string
	withExpr           *string
	assignedTo         *string
	mustReport         *string
	withPropReturnType *string
}

// Java: OptionalChainTypeCheckTest.OptChainTestCase.Builder.forObjType(String)
func (b *OptionalChainTypeCheckTest_OptChainTestCase_Builder) ForObjType(objType string) *OptionalChainTypeCheckTest_OptChainTestCase_Builder {
	b.forObjType = &objType
	return b
}

// Java: OptionalChainTypeCheckTest.OptChainTestCase.Builder.withExpr(String)
func (b *OptionalChainTypeCheckTest_OptChainTestCase_Builder) WithExpr(expr string) *OptionalChainTypeCheckTest_OptChainTestCase_Builder {
	b.withExpr = &expr
	return b
}

// Java: OptionalChainTypeCheckTest.OptChainTestCase.Builder.assignedTo(String)
func (b *OptionalChainTypeCheckTest_OptChainTestCase_Builder) AssignedTo(type_ string) *OptionalChainTypeCheckTest_OptChainTestCase_Builder {
	b.assignedTo = &type_
	return b
}

// Java: OptionalChainTypeCheckTest.OptChainTestCase.Builder.mustReport(String)
func (b *OptionalChainTypeCheckTest_OptChainTestCase_Builder) MustReport(error_ string) *OptionalChainTypeCheckTest_OptChainTestCase_Builder {
	b.mustReport = &error_
	return b
}

// WithPropReturnType sets the return type of the function property, if any
// on this object type.
//
// Java: OptionalChainTypeCheckTest.OptChainTestCase.Builder.withPropReturnType(String)
func (b *OptionalChainTypeCheckTest_OptChainTestCase_Builder) WithPropReturnType(type_ string) *OptionalChainTypeCheckTest_OptChainTestCase_Builder {
	b.withPropReturnType = &type_
	return b
}

// Build is the generated AutoValue build(): withExpr is the only required
// property.
//
// Java: OptionalChainTypeCheckTest.OptChainTestCase.Builder.build()
func (b *OptionalChainTypeCheckTest_OptChainTestCase_Builder) Build() *OptionalChainTypeCheckTest_OptChainTestCase {
	if b.withExpr == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties: withExpr"))
	}
	return &OptionalChainTypeCheckTest_OptChainTestCase{
		forObjType:         b.forObjType,
		withExpr:           *b.withExpr,
		assignedTo:         b.assignedTo,
		mustReport:         b.mustReport,
		withPropReturnType: b.withPropReturnType,
	}
}

// OptionalChainTypeCheckTest_CreateTestString constructs the test case
// string from a test case object. e.g. `/** @type {({b:number})} *\/ var a;
// var x; x =a?.[b];`
//
// Java: OptionalChainTypeCheckTest.createTestString(OptChainTestCase)
func OptionalChainTypeCheckTest_CreateTestString(testCase *OptionalChainTypeCheckTest_OptChainTestCase) string {
	objType := *testCase.ForObjType()
	expr := testCase.WithExpr()
	expectType := *testCase.AssignedTo()
	return jutil.StringFormat(`/** @type {(%s)} */
var a;
/** @type {(%s)} */
var x;
x = %s;
`, objType, expectType, expr)
}

// optionalChainTypeCheckTest_runParameterized is JUnit's Parameterized
// runner: for each parameter, every test method of the class runs on a
// fresh instance whose @Parameter field newTest sets.
//
// Go-only: test-helper — the Parameterized runner.
func optionalChainTypeCheckTest_runParameterized(
	t *testing.T,
	cases []*OptionalChainTypeCheckTest_OptChainTestCase,
	newTest func(testCase *OptionalChainTypeCheckTest_OptChainTestCase) JUnitTestCase) {
	t.Helper()
	for i, testCase := range cases {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			RunJUnit4(t, func() JUnitTestCase { return newTest(testCase) })
		})
	}
}

// OptionalChainTypeCheckTest_OptChainGetElemTests is the Parameterized
// nested class OptChainGetElemTests.
type OptionalChainTypeCheckTest_OptChainGetElemTests struct {
	typeCheckTests_noFixture
	testCase *OptionalChainTypeCheckTest_OptChainTestCase
}

// Java: OptionalChainTypeCheckTest.OptChainGetElemTests()
func NewOptionalChainTypeCheckTest_OptChainGetElemTests() *OptionalChainTypeCheckTest_OptChainGetElemTests {
	return &OptionalChainTypeCheckTest_OptChainGetElemTests{}
}

func TestOptionalChainTypeCheckTest_OptChainGetElemTests(t *testing.T) {
	optionalChainTypeCheckTest_runParameterized(t, OptionalChainTypeCheckTest_OptChainGetElemTests_Cases(),
		func(testCase *OptionalChainTypeCheckTest_OptChainTestCase) JUnitTestCase {
			c := NewOptionalChainTypeCheckTest_OptChainGetElemTests()
			c.testCase = testCase
			return c
		})
}

// OptionalChainTypeCheckTest_OptChainGetElemTests_Cases generates test
// cases for testing `?.[]` OPTCHAIN_GETPROP expressions first followed by
// normal `[]` GETPROP expressions. Each generated test case is executed by
// the `test()` method.
//
// Java: OptionalChainTypeCheckTest.OptChainGetElemTests.cases()
func OptionalChainTypeCheckTest_OptChainGetElemTests_Cases() []*OptionalChainTypeCheckTest_OptChainTestCase {
	return []*OptionalChainTypeCheckTest_OptChainTestCase{
		// Inferred type for [] access on an Object is always unknown. This is consistent
		// for optional as well as normal [] access.
		// Then the unknown value can be assigned to the variable no matter what its type is.
		// Note the intentionally different types on the property and the assigned-to variable
		// which illustrate this.
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{prop: number}").
			WithExpr("a[index]").
			AssignedTo("string").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{prop: number}").
			WithExpr("a?.[index]").
			AssignedTo("string").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{prop: number}").
			WithExpr("a['index'];").
			AssignedTo("string").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{prop: number}").
			WithExpr("a?.['index'];").
			AssignedTo("string").
			Build(),
		// missing prop
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{b:number}").
			WithExpr("a?.c[b]").
			AssignedTo("?").
			MustReport("Property c never defined on a").
			Build(),
		// Tests that Arrays and Objects are inferred correctly by OPTCHAIN_GETELEM
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("Array<number>").
			WithExpr("a?.[0];").
			AssignedTo("number|undefined").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("Array<number|string>").
			WithExpr("a?.[0];").
			AssignedTo("number|string|undefined").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("Array<number|?string>").
			WithExpr("a?.[0];").
			AssignedTo("number|?string|undefined").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("Array<?number>").
			WithExpr("a?.[0];").
			AssignedTo("?number|undefined").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("Object<string, number>").
			WithExpr("a?.['b'];").
			AssignedTo("number|undefined").
			Build(),
		// type mismatch on assignments using `?.[]` on `Array` and `Object` are reported
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("Array<number>").
			WithExpr("a?.[0];").
			AssignedTo("string").
			MustReport(`assignment
found   : (number|undefined)
required: string
`).
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("Array<?number>").
			WithExpr("a?.[0];").
			AssignedTo("number|undefined").
			MustReport(`assignment
found   : (null|number|undefined)
required: (number|undefined)
`).
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("null").
			WithExpr("a?.['b'];").
			AssignedTo("undefined").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{b:number}").
			WithExpr("a?.[c];").
			AssignedTo("?").
			Build(),
		// testing non-optional exprs
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{b:(?|number)}").
			WithExpr("a[b];").
			AssignedTo("(?|number)").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{b:number}").
			WithExpr("a[b];").
			AssignedTo("number").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("?{b:number}").
			WithExpr("a[b];").
			AssignedTo("number").
			Build(),
		// error cases - when using unguarded null|undefined LHS non-optionally
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("null").
			WithExpr("a[b];").
			AssignedTo("?").
			MustReport(`only arrays or objects can be accessed
found   : null
required: Object
`).
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("null").
			WithExpr("a.b?.[c];").
			AssignedTo("?").
			MustReport(`No properties on this expression
found   : null
required: Object
`).
			Build(),
	}
}

// Java: OptionalChainTypeCheckTest.OptChainGetElemTests.test()
func (c *OptionalChainTypeCheckTest_OptChainGetElemTests) Test() {
	js := OptionalChainTypeCheckTest_CreateTestString(c.testCase)
	if c.testCase.MustReport() == nil {
		TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(js).Run()
	} else {
		TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(js).AddDiagnostic(*c.testCase.MustReport()).Run()
	}
}

// OptionalChainTypeCheckTest_OptChainGetPropTests is the Parameterized
// nested class OptChainGetPropTests.
type OptionalChainTypeCheckTest_OptChainGetPropTests struct {
	typeCheckTests_noFixture
	testCase *OptionalChainTypeCheckTest_OptChainTestCase
}

// Java: OptionalChainTypeCheckTest.OptChainGetPropTests()
func NewOptionalChainTypeCheckTest_OptChainGetPropTests() *OptionalChainTypeCheckTest_OptChainGetPropTests {
	return &OptionalChainTypeCheckTest_OptChainGetPropTests{}
}

func TestOptionalChainTypeCheckTest_OptChainGetPropTests(t *testing.T) {
	optionalChainTypeCheckTest_runParameterized(t, OptionalChainTypeCheckTest_OptChainGetPropTests_Cases(),
		func(testCase *OptionalChainTypeCheckTest_OptChainTestCase) JUnitTestCase {
			c := NewOptionalChainTypeCheckTest_OptChainGetPropTests()
			c.testCase = testCase
			return c
		})
}

// Java: OptionalChainTypeCheckTest.OptChainGetPropTests.cases()
func OptionalChainTypeCheckTest_OptChainGetPropTests_Cases() []*OptionalChainTypeCheckTest_OptChainTestCase {
	return []*OptionalChainTypeCheckTest_OptChainTestCase{
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{b:number}").
			WithExpr("a?.b;").
			AssignedTo("number").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{b:?number}").
			WithExpr("a?.b;").
			AssignedTo("?number").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("?{b:number}").
			WithExpr("a?.b;").
			AssignedTo("undefined|number").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("?{b:?number}").
			WithExpr("a?.b;").
			AssignedTo("undefined|(?number)").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("null").
			WithExpr("a?.b;").
			AssignedTo("undefined").
			Build(),
		// missing property
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{b:number}").
			WithExpr("a?.c;").
			AssignedTo("?").
			MustReport("Property c never defined on a").
			Build(),
		// type mismatch reported
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{b:number}").
			WithExpr("a?.b;").
			AssignedTo("string").
			MustReport(`assignment
found   : number
required: string
`).
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("?{b:number}").
			WithExpr("a?.b;").
			AssignedTo("undefined|string").
			MustReport(`assignment
found   : (number|undefined)
required: (string|undefined)
`).
			Build(),
		// normal GETPROP
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("null").
			WithExpr("a.alert();").
			AssignedTo("?").
			MustReport(`No properties on this expression
found   : null
required: Object
`).
			Build(),
	}
}

// Java: OptionalChainTypeCheckTest.OptChainGetPropTests.test()
func (c *OptionalChainTypeCheckTest_OptChainGetPropTests) Test() {
	js := OptionalChainTypeCheckTest_CreateTestString(c.testCase)
	if c.testCase.MustReport() == nil {
		TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(js).Run()
	} else {
		TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(js).AddDiagnostic(*c.testCase.MustReport()).Run()
	}
}

// OptionalChainTypeCheckTest_OptChainAccessOnDictTests is the
// Parameterized nested class OptChainAccessOnDictTests.
type OptionalChainTypeCheckTest_OptChainAccessOnDictTests struct {
	TypeCheckTestCaseBase
	testCase *OptionalChainTypeCheckTest_OptChainTestCase
}

// Java: OptionalChainTypeCheckTest.OptChainAccessOnDictTests()
func NewOptionalChainTypeCheckTest_OptChainAccessOnDictTests() *OptionalChainTypeCheckTest_OptChainAccessOnDictTests {
	c := &OptionalChainTypeCheckTest_OptChainAccessOnDictTests{}
	c.InitTypeCheckTestCaseBase(c)
	return c
}

func TestOptionalChainTypeCheckTest_OptChainAccessOnDictTests(t *testing.T) {
	optionalChainTypeCheckTest_runParameterized(t, OptionalChainTypeCheckTest_OptChainAccessOnDictTests_Cases(),
		func(testCase *OptionalChainTypeCheckTest_OptChainTestCase) JUnitTestCase {
			c := NewOptionalChainTypeCheckTest_OptChainAccessOnDictTests()
			c.testCase = testCase
			return c
		})
}

// Java: OptionalChainTypeCheckTest.OptChainAccessOnDictTests.cases()
func OptionalChainTypeCheckTest_OptChainAccessOnDictTests_Cases() []*OptionalChainTypeCheckTest_OptChainTestCase {
	return []*OptionalChainTypeCheckTest_OptChainTestCase{
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			WithExpr("x?.prop;").
			MustReport("Cannot do '?.' access on a dict").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			WithExpr("x.prop;").
			MustReport("Cannot do '.' access on a dict").
			Build(),
	}
}

// Java: OptionalChainTypeCheckTest.OptChainAccessOnDictTests.testOptChainGetProp_accessOnDict1()
func (c *OptionalChainTypeCheckTest_OptChainAccessOnDictTests) TestOptChainGetProp_accessOnDict1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(
			strings.ReplaceAll(`/**
 * @constructor
 * @dict
 */
function Dict1(){ this['prop'] = 123; }
/** @param{Dict1} x */
function takesDict(x) {
  return WITH_EXPR
}
`, "WITH_EXPR", c.testCase.WithExpr())).
		AddDiagnostic(*c.testCase.MustReport()).
		Run()
}

// Java: OptionalChainTypeCheckTest.OptChainAccessOnDictTests.testOptChainGetPropDict2()
func (c *OptionalChainTypeCheckTest_OptChainAccessOnDictTests) TestOptChainGetPropDict2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(
			strings.ReplaceAll(`/**
 * @constructor
 * @dict
 */
function Dict1(){ this['prop'] = 123; }
/**
 * @constructor
 * @extends {Dict1}
 */
function Dict1kid(){ this['prop'] = 123; }
/** @param{Dict1kid} x */
function takesDict(x) {
  return WITH_EXPR
 }
`, "WITH_EXPR", c.testCase.WithExpr())).
		AddDiagnostic(*c.testCase.MustReport()).
		Run()
}

// Java: OptionalChainTypeCheckTest.OptChainAccessOnDictTests.testOptChainGetPropDict_accessingDictOrNonDict()
func (c *OptionalChainTypeCheckTest_OptChainAccessOnDictTests) TestOptChainGetPropDict_accessingDictOrNonDict() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(
			strings.ReplaceAll(`/**
 * @constructor
 * @dict
 */
function Dict1() { this['prop'] = 123; }
/** @constructor */
function NonDict() { this.prop = 321; }
/** @param{(NonDict|Dict1)} x */
function takesDict(x) {
  return WITH_EXPR
}
`, "WITH_EXPR", c.testCase.WithExpr())).
		AddDiagnostic(*c.testCase.MustReport()).
		Run()
}

// Java: OptionalChainTypeCheckTest.OptChainAccessOnDictTests.testOptChainGetProp_accessingDictOrStruct()
func (c *OptionalChainTypeCheckTest_OptChainAccessOnDictTests) TestOptChainGetProp_accessingDictOrStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(
			strings.ReplaceAll(`/**
 * @constructor
 * @dict
 */
function Dict1() { this['prop'] = 123; }
/**
 * @constructor
 * @struct
 */
function Struct1() { this.prop = 123; }
/** @param{(Struct1|Dict1)} x */
function takesNothing(x) {
  return WITH_EXPR
}
`, "WITH_EXPR", c.testCase.WithExpr())).
		AddDiagnostic(*c.testCase.MustReport()).
		Run()
}

// OptionalChainTypeCheckTest_OptChainCallTests is the Parameterized nested
// class OptChainCallTests.
type OptionalChainTypeCheckTest_OptChainCallTests struct {
	typeCheckTests_noFixture
	testCase *OptionalChainTypeCheckTest_OptChainTestCase
}

// Java: OptionalChainTypeCheckTest.OptChainCallTests()
func NewOptionalChainTypeCheckTest_OptChainCallTests() *OptionalChainTypeCheckTest_OptChainCallTests {
	return &OptionalChainTypeCheckTest_OptChainCallTests{}
}

func TestOptionalChainTypeCheckTest_OptChainCallTests(t *testing.T) {
	optionalChainTypeCheckTest_runParameterized(t, OptionalChainTypeCheckTest_OptChainCallTests_Cases(),
		func(testCase *OptionalChainTypeCheckTest_OptChainTestCase) JUnitTestCase {
			c := NewOptionalChainTypeCheckTest_OptChainCallTests()
			c.testCase = testCase
			return c
		})
}

// Java: OptionalChainTypeCheckTest.OptChainCallTests.cases()
func OptionalChainTypeCheckTest_OptChainCallTests_Cases() []*OptionalChainTypeCheckTest_OptChainTestCase {
	return []*OptionalChainTypeCheckTest_OptChainTestCase{
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			WithPropReturnType("string").
			WithExpr("a?.prop()").
			AssignedTo("string").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			WithPropReturnType("?string").
			WithExpr("a?.prop()").
			AssignedTo("?string").
			Build(),
		// testing null objType
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("null").
			WithExpr("a?.prop()").
			AssignedTo("undefined").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("null").
			WithExpr("a?.()").
			AssignedTo("undefined").
			Build(),
		// missing prop
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("{b:number}").
			WithExpr("a?.c()").
			AssignedTo("?").
			MustReport("Property c never defined on a").
			Build(),
		// mismatch
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			WithPropReturnType("string").
			WithExpr("a?.prop()").
			AssignedTo("number").
			MustReport(`assignment
found   : string
required: number
`).
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			WithPropReturnType("?string").
			WithExpr("a?.prop()").
			AssignedTo("string").
			MustReport(`assignment
found   : (null|string)
required: string
`).
			Build(),
		// non-optional tests
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			WithPropReturnType("string").
			WithExpr("a.prop()").
			AssignedTo("string").
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			WithPropReturnType("?string").
			WithExpr("a.prop()").
			AssignedTo("?string").
			Build(),
		// type mismatch
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			WithPropReturnType("string").
			WithExpr("a.prop()").
			AssignedTo("number").
			MustReport(`assignment
found   : string
required: number
`).
			Build(),
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			WithPropReturnType("?string").
			WithExpr("a.prop()").
			AssignedTo("string").
			MustReport(`assignment
found   : (null|string)
required: string
`).
			Build(),
		// We don't report an error when a null or undefined type is called as a function
		OptionalChainTypeCheckTest_OptChainTestCase_Builder_().
			ForObjType("null").
			WithExpr("a()").
			AssignedTo("?").
			Build(),
	}
}

// Java: OptionalChainTypeCheckTest.OptChainCallTests.createOptChainCallTestString(OptChainTestCase)
func OptionalChainTypeCheckTest_OptChainCallTests_CreateOptChainCallTestString(testCase *OptionalChainTypeCheckTest_OptChainTestCase) string {
	var js strings.Builder
	js.WriteString("let a = function() {};")
	if testCase.WithPropReturnType() != nil {
		js.WriteString("/** @type {function():")
		js.WriteString(*testCase.WithPropReturnType())
		js.WriteString(" } */")
	}
	js.WriteString("  a.prop = function() {};")
	js.WriteString("/**@type {")
	js.WriteString(*testCase.AssignedTo())
	js.WriteString("} */var x;")
	js.WriteString("x = ")
	js.WriteString(testCase.WithExpr())
	return js.String()
}

// Java: OptionalChainTypeCheckTest.OptChainCallTests.testOptChainCallExpressions()
func (c *OptionalChainTypeCheckTest_OptChainCallTests) TestOptChainCallExpressions() {
	js := OptionalChainTypeCheckTest_OptChainCallTests_CreateOptChainCallTestString(c.testCase)
	if c.testCase.MustReport() == nil {
		TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(js).Run()
	} else {
		TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(js).AddDiagnostic(*c.testCase.MustReport()).Run()
	}
}

// OptionalChainTypeCheckTest_OptChainTestsNonParameterized is the nested
// class OptChainTestsNonParameterized.
type OptionalChainTypeCheckTest_OptChainTestsNonParameterized struct {
	TypeCheckTestCaseBase
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized()
func NewOptionalChainTypeCheckTest_OptChainTestsNonParameterized() *OptionalChainTypeCheckTest_OptChainTestsNonParameterized {
	c := &OptionalChainTypeCheckTest_OptChainTestsNonParameterized{}
	c.InitTypeCheckTestCaseBase(c)
	return c
}

func TestOptionalChainTypeCheckTest_OptChainTestsNonParameterized(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewOptionalChainTypeCheckTest_OptChainTestsNonParameterized() })
}

// TestOptChainGetElemExpressions_nonNullObject confirms that
// OPTCHAIN_GETELEM nodes are inferred as unknown type.
//
// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testOptChainGetElemExpressions_nonNullObject()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestOptChainGetElemExpressions_nonNullObject() {
	js := "/** @type {({b:number})} */ var a; var x; x = a?.[b]"

	var script *rhino.Node = c.ParseAndTypeCheck(js)
	var assign *rhino.Node = script.GetLastChild().GetFirstChild()
	AssertThat(assign.IsAssign()).IsTrue()

	var xName *rhino.Node = assign.GetFirstChild()
	AssertThat(xName.IsName()).IsTrue()
	AssertThat(xName.GetJSType().IsUnknownType()).IsTrue()

	var getElem *rhino.Node = assign.GetSecondChild()
	AssertThat(getElem.IsOptChainGetElem() || getElem.IsGetElem()).IsTrue()
	AssertThat(getElem.GetJSType().IsUnknownType()).IsTrue()
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testOptChainGetElemOnStruct1()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestOptChainGetElemOnStruct1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(`/**
 * @constructor
 * @struct
 */
function AStruct(){ this.b = 123; }
/** @param{AStruct} x */
function takesStruct(x) {
  var a = x;
  return a?.[b];
}
`).
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testOptChainGetElemOnStruct2()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestOptChainGetElemOnStruct2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(`/**
 * @constructor
 * @struct
 */
function Struct1(){ this.b = 123; }
/**
 * @constructor
 * @extends {Struct1}
 */
function Struct1kid(){ this.b = 123; }
/** @param{Struct1kid} a */
function takesStruct2(a) {
  return a?.[b];
}
`).
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testGetProp_propAccessOnVoidFunction()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestGetProp_propAccessOnVoidFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource("/** @return {void}*/function foo(){foo().bar;}").
		AddDiagnostic(`No properties on this expression
found   : undefined
required: Object
`).
		Run()
}

// TestOptChainGetProp_propAccessOnVoidFunction_guarded: optionally
// accessing a prop on NULL_TYPE or VOID_TYPE is not an error.
//
// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testOptChainGetProp_propAccessOnVoidFunction_guarded()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestOptChainGetProp_propAccessOnVoidFunction_guarded() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {void}*/function foo(){foo()?.bar;}").Run()
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testOptChainGetProp_propAccessedAfterObjectInitialized()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestOptChainGetProp_propAccessedAfterObjectInitialized() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(`/** @constructor */
function Foo() { /** @type {?Object} */ this.x = null; }
Foo.prototype.initX = function() { this.x = {foo: 1}; };
Foo.prototype.bar = function() {
  if (this.x == null) { this.initX(); alert(this.x?.foo); }
};
`).
		Run()
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testAssignToUnknown_noError()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestAssignToUnknown_noError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(`let a = 4;
/** @type {?} */ let b;
a =b;
`).
		Run()
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testOptChainCall_numberCalled()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestOptChainCall_numberCalled() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("3?.();").AddDiagnostic("number expressions are not callable").Run()
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testCall_paramMismatch()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestCall_paramMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource("/** @param {!Number} foo*/function bar(foo){ bar('abc'); }").
		AddDiagnostic(`actual parameter 1 of bar does not match formal parameter
found   : string
required: Number
`).
		Run()
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testOptChainCall_paramMismatch()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestOptChainCall_paramMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource("/** @param {!Number} foo*/function bar(foo){ bar?.('abc'); }").
		AddDiagnostic(`actual parameter 1 of bar does not match formal parameter
found   : string
required: Number
`).
		Run()
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testCall3()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestCall3() {
	// We are checking that an unresolved named type can successfully
	// meet with a functional type to produce a callable type.
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(`/** @type {Function|undefined} */var opt_f;
/** @type {some.unknown.type} */var f1;
var f2 = opt_f || f1;
f2();
`).
		AddDiagnostic("Bad type annotation. Unknown type some.unknown.type").
		Run()
}

// Java: OptionalChainTypeCheckTest.OptChainTestsNonParameterized.testOptChainCall3()
func (c *OptionalChainTypeCheckTest_OptChainTestsNonParameterized) TestOptChainCall3() {
	// We are checking that an unresolved named type can successfully
	// meet with a functional type to produce a callable type.
	TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(`/** @type {Function|undefined} */var opt_f;
/** @type {some.unknown.type} */var f1;
var f2 = opt_f || f1;
f2();
`).
		AddDiagnostic("Bad type annotation. Unknown type some.unknown.type").
		Run()
}

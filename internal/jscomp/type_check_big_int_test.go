// Port of com/google/javascript/jscomp/TypeCheckBigIntTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/TypeCheckBigIntTest.java)

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
)

// TypeCheckBigIntTest tests the type checker.
type TypeCheckBigIntTest struct {
	typeCheckTests_noFixture
}

// Java: TypeCheckBigIntTest()
func NewTypeCheckBigIntTest() *TypeCheckBigIntTest {
	c := &TypeCheckBigIntTest{}
	return c
}

func TestTypeCheckBigIntTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckBigIntTest() })
}

// Java: TypeCheckBigIntTest.testTypeofBigInt()
func (c *TypeCheckBigIntTest) TestTypeofBigInt() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {bigint|number} i
 * @return {boolean}
 */
function foo(i) {
  return typeof i === 'bigint';
}
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntArgument()
func (c *TypeCheckBigIntTest) TestBigIntArgument() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("BigInt(1)").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("BigInt({})").
		AddDiagnostic(`actual parameter 1 of BigInt does not match formal parameter
found   : {}
required: (bigint|number|string)
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_increment()
func (c *TypeCheckBigIntTest) TestBigIntOperators_increment() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; x++;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!BigInt} */ var x; x++;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */var x; x++;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */var x; x++;").
		AddDiagnostic(`increment/decrement
found   : (bigint|string)
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_decrement()
func (c *TypeCheckBigIntTest) TestBigIntOperators_decrement() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; x--;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!BigInt} */ var x; x--;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */var x; x--;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */var x; x--;").
		AddDiagnostic(`increment/decrement
found   : (bigint|string)
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_logicalNot()
func (c *TypeCheckBigIntTest) TestBigIntOperators_logicalNot() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; !x;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!BigInt} */ var x; !x;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */var x; !x;").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_bitwiseNot()
func (c *TypeCheckBigIntTest) TestBigIntOperators_bitwiseNot() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; ~x;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!BigInt} */ var x; ~x;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {?} */var x; ~x;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */var x; ~x;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */var x; ~x;").
		AddDiagnostic(`bitwise NOT
found   : (bigint|string)
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntValueOperators_unaryPlusIsForbidden()
func (c *TypeCheckBigIntTest) TestBigIntValueOperators_unaryPlusIsForbidden() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1n; +x;").
		AddDiagnostic("unary operator + cannot be applied to bigint").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntObjectOperators_unaryPlusIsForbidden()
func (c *TypeCheckBigIntTest) TestBigIntObjectOperators_unaryPlusIsForbidden() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!BigInt} */ var x; +x;").
		AddDiagnostic("unary operator + cannot be applied to BigInt").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntUnionOperators_unaryPlusIsForbidden()
func (c *TypeCheckBigIntTest) TestBigIntUnionOperators_unaryPlusIsForbidden() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x; +x;").
		AddDiagnostic("unary operator + cannot be applied to (bigint|number)").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntEnumOperators_unaryPlusIsForbidden()
func (c *TypeCheckBigIntTest) TestBigIntEnumOperators_unaryPlusIsForbidden() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @enum {bigint} */ const BIGINTS = {ONE: 1n, TWO: 2n}; +BIGINTS.ONE;").
		AddDiagnostic("unary operator + cannot be applied to BIGINTS<bigint>").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_unaryMinus()
func (c *TypeCheckBigIntTest) TestBigIntOperators_unaryMinus() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; -x;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!BigInt} */ var x; -x;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {?} */var x; -x;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */var x; -x;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */var x; -x;").
		AddDiagnostic(`unary minus operator
found   : (bigint|string)
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_binaryOperationWithSelf()
func (c *TypeCheckBigIntTest) TestBigIntOperators_binaryOperationWithSelf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 1n; x * y").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_assignOpWithSelf()
func (c *TypeCheckBigIntTest) TestBigIntOperators_assignOpWithSelf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 1n; x *= y").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_binaryBitwiseOperationWithSelf()
func (c *TypeCheckBigIntTest) TestBigIntOperators_binaryBitwiseOperationWithSelf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 1n; x | y").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_additionWithString()
func (c *TypeCheckBigIntTest) TestBigIntOperators_additionWithString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 'str'; x + y").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_assignAddWithString()
func (c *TypeCheckBigIntTest) TestBigIntOperators_assignAddWithString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 'str'; x += y").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_binaryOperationWithString()
func (c *TypeCheckBigIntTest) TestBigIntOperators_binaryOperationWithString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 'str'; x * y").
		AddDiagnostic("operator * cannot be applied to bigint and string").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_assignOpWithString()
func (c *TypeCheckBigIntTest) TestBigIntOperators_assignOpWithString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 'str'; x *= y").
		AddDiagnostic("operator *= cannot be applied to bigint and string").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_binaryOperationWithUnknown()
func (c *TypeCheckBigIntTest) TestBigIntOperators_binaryOperationWithUnknown() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1n; /** @type {?} */var y; x * y").
		AddDiagnostic("operator * cannot be applied to bigint and ?").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_assignOpWithUnknown()
func (c *TypeCheckBigIntTest) TestBigIntOperators_assignOpWithUnknown() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1n; /** @type {?} */var y; x *= y").
		AddDiagnostic("operator *= cannot be applied to bigint and ?").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_binaryOperationWithNumber()
func (c *TypeCheckBigIntTest) TestBigIntOperators_binaryOperationWithNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 1; x * y").
		AddDiagnostic("operator * cannot be applied to bigint and number").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_assignOpWithNumber()
func (c *TypeCheckBigIntTest) TestBigIntOperators_assignOpWithNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 1; x *= y").
		AddDiagnostic("operator *= cannot be applied to bigint and number").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_binaryBitwiseOperationWithNumber()
func (c *TypeCheckBigIntTest) TestBigIntOperators_binaryBitwiseOperationWithNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 1; x | y").
		AddDiagnostic("operator | cannot be applied to bigint and number").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntLeftShift()
func (c *TypeCheckBigIntTest) TestBigIntLeftShift() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("1n << 2n").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("let x = 1n; x <<= 2n").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntRightShift()
func (c *TypeCheckBigIntTest) TestBigIntRightShift() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("2n >> 1n").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("let x = 2n; x >>= 1n").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_unsignedRightShift()
func (c *TypeCheckBigIntTest) TestBigIntOperators_unsignedRightShift() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; x >>> x;").
		AddDiagnostic("operator >>> cannot be applied to bigint and bigint").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOperators_assignUnsignedRightShift()
func (c *TypeCheckBigIntTest) TestBigIntOperators_assignUnsignedRightShift() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("let x = 1n; x >>>= x;").
		AddDiagnostic("operator >>>= cannot be applied to bigint and bigint").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrNumberOperators_binaryOperationWithSelf()
func (c *TypeCheckBigIntTest) TestBigIntOrNumberOperators_binaryOperationWithSelf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x; x * x;").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrNumberOperators_assignOpWithSelf()
func (c *TypeCheckBigIntTest) TestBigIntOrNumberOperators_assignOpWithSelf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x; x *= x;").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrNumberOperators_binaryBitwiseOperationWithSelf()
func (c *TypeCheckBigIntTest) TestBigIntOrNumberOperators_binaryBitwiseOperationWithSelf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x; x | x;").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrNumberOperators_binaryOperationWithBigInt()
func (c *TypeCheckBigIntTest) TestBigIntOrNumberOperators_binaryOperationWithBigInt() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1n; /** @type {bigint|number} */ var y; x * y;").
		AddDiagnostic("operator * cannot be applied to bigint and (bigint|number)").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrNumberOperators_assignOpWithBigInt()
func (c *TypeCheckBigIntTest) TestBigIntOrNumberOperators_assignOpWithBigInt() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1n; /** @type {bigint|number} */ var y; x *= y;").
		AddDiagnostic("operator *= cannot be applied to bigint and (bigint|number)").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrNumberOperators_binaryBitwiseOperationWithBigInt()
func (c *TypeCheckBigIntTest) TestBigIntOrNumberOperators_binaryBitwiseOperationWithBigInt() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1n; /** @type {bigint|number} */ var y; x | y;").
		AddDiagnostic("operator | cannot be applied to bigint and (bigint|number)").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrNumberOperators_unsignedRightShift()
func (c *TypeCheckBigIntTest) TestBigIntOrNumberOperators_unsignedRightShift() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x = 1n; x >>> x;").
		AddDiagnostic("operator >>> cannot be applied to (bigint|number) and (bigint|number)").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrNumberOperators_assignUnsignedRightShift()
func (c *TypeCheckBigIntTest) TestBigIntOrNumberOperators_assignUnsignedRightShift() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x = 1n; x >>>= x;").
		AddDiagnostic("operator >>>= cannot be applied to (bigint|number) and (bigint|number)").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrOtherOperators_binaryOperationWithBigInt()
func (c *TypeCheckBigIntTest) TestBigIntOrOtherOperators_binaryOperationWithBigInt() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */ var x; var y = 1n; x * y;").
		AddDiagnostic("operator * cannot be applied to (bigint|string) and bigint").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrOtherOperators_assignOpWithBigInt()
func (c *TypeCheckBigIntTest) TestBigIntOrOtherOperators_assignOpWithBigInt() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */ var x; var y = 1n; x *= y;").
		AddDiagnostic("operator *= cannot be applied to (bigint|string) and bigint").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrOtherOperators_binaryOperationWithBigIntOrNumber()
func (c *TypeCheckBigIntTest) TestBigIntOrOtherOperators_binaryOperationWithBigIntOrNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */ var x; /** @type {bigint|number} */ var y; x * y;").
		AddDiagnostic("operator * cannot be applied to (bigint|string) and (bigint|number)").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrOtherOperators_assignOpWithBigIntOrNumber()
func (c *TypeCheckBigIntTest) TestBigIntOrOtherOperators_assignOpWithBigIntOrNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */ var x; /** @type {bigint|number} */ var y; x *= y;").
		AddDiagnostic("operator *= cannot be applied to (bigint|string) and (bigint|number)").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntEnumOperators_binaryOperationWithSelf()
func (c *TypeCheckBigIntTest) TestBigIntEnumOperators_binaryOperationWithSelf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @enum {bigint} */ const BIGINTS = {ONE: 1n, TWO: 2n}; BIGINTS.ONE * BIGINTS.TWO;").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntEnumOperators_assignOpWithSelf()
func (c *TypeCheckBigIntTest) TestBigIntEnumOperators_assignOpWithSelf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @enum {bigint} */ const BIGINTS = {ONE: 1n, TWO: 2n}; BIGINTS.ONE *= BIGINTS.TWO;").
		Run()
}

// Java: TypeCheckBigIntTest.testValidBigIntComparisons()
func (c *TypeCheckBigIntTest) TestValidBigIntComparisons() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1n; var y = 2n; x < y").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1n; /** @type {!BigInt} */ var y; x < y").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1n; var y = 2; x < y").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1n; /** @type {?} */ var y; x < y").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntComparisonWithString()
func (c *TypeCheckBigIntTest) TestBigIntComparisonWithString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const x = 1n; const y = 'asdf'; x < y;").
		AddDiagnostic(`right side of numeric comparison
found   : string
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckBigIntTest.testValidBigIntObjectComparisons()
func (c *TypeCheckBigIntTest) TestValidBigIntObjectComparisons() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!BigInt} */ var x; var y = 2n; x < y").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!BigInt} */ var x; /** @type {!BigInt} */ var y; x < y").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!BigInt} */ var x; var y = 2; x < y").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!BigInt} */ var x; /** @type {?} */ var y; x < y").
		Run()
}

// Java: TypeCheckBigIntTest.testValidBigIntOrNumberComparisons()
func (c *TypeCheckBigIntTest) TestValidBigIntOrNumberComparisons() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x; /** @type {bigint|number} */ var y; x < y;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x; var y = 2; x < y;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x; var y = 2n; x < y;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x; /** @type {!BigInt} */ var y; x < y;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x; /** @type {?} */ var y; x < y").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrNumberComparisonWithString()
func (c *TypeCheckBigIntTest) TestBigIntOrNumberComparisonWithString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|number} */ var x; 'asdf' < x;").
		AddDiagnostic(`left side of numeric comparison
found   : string
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckBigIntTest.testValidBigIntOrOtherComparisons()
func (c *TypeCheckBigIntTest) TestValidBigIntOrOtherComparisons() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */ var x; /** @type {bigint|string} */ var y; x < y;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */ var x; /** @type {?} */ var y; x < y").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrOtherComparisonWithBigint()
func (c *TypeCheckBigIntTest) TestBigIntOrOtherComparisonWithBigint() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */ var x; var y = 2n; x < y;").
		AddDiagnostic(`left side of numeric comparison
found   : (bigint|string)
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrOtherComparisonWithNumber()
func (c *TypeCheckBigIntTest) TestBigIntOrOtherComparisonWithNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint|string} */ var x; var y = 2; x < y;").
		AddDiagnostic(`left side of numeric comparison
found   : (bigint|string)
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntObjectIndex()
func (c *TypeCheckBigIntTest) TestBigIntObjectIndex() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var obj = {};
/** @type {bigint} */ var b;
/** @type {bigint|number} */ var bn;
obj[b] = 1;
obj[bn] = 3;
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntArrayIndex()
func (c *TypeCheckBigIntTest) TestBigIntArrayIndex() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var arr = []; /** @type {bigint} */ var b; arr[b];").
		AddDiagnostic(`restricted index type
found   : bigint
required: number
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntConstructorWithNew()
func (c *TypeCheckBigIntTest) TestBigIntConstructorWithNew() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("new BigInt(1)").
		AddDiagnostic("cannot instantiate non-constructor, found type: function(new:BigInt, (bigint|number|string)): bigint").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntOrNumberArrayIndex()
func (c *TypeCheckBigIntTest) TestBigIntOrNumberArrayIndex() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var arr = []; /** @type {bigint|number} */ var bn; arr[bn];").
		AddDiagnostic(`restricted index type
found   : (bigint|number)
required: number
`).
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntAsComputedPropForObjects()
func (c *TypeCheckBigIntTest) TestBigIntAsComputedPropForObjects() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {bigint} */ var x; ({[x]: 'value', 123n() {}});").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntAsComputedPropForClasses()
func (c *TypeCheckBigIntTest) TestBigIntAsComputedPropForClasses() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @unrestricted */ class C { 123n() {} }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { 123n() {} }").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntAsComputedPropForStructClasses()
func (c *TypeCheckBigIntTest) TestBigIntAsComputedPropForStructClasses() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class C { 123n() {} }").
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckBigIntTest.testBigIntLiteralProperty()
func (c *TypeCheckBigIntTest) TestBigIntLiteralProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddBigInt().Build()).
		AddSource("(1n).toString()").
		Run()
}

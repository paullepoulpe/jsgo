// Port of com/google/javascript/jscomp/PeepholeFoldConstantsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/PeepholeFoldConstantsTest.java)

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PeepholeFoldConstantsTest is the tests for PeepholeFoldConstants in
// isolation. Tests for the interaction of multiple peephole passes are in
// PeepholeIntegrationTest.
type PeepholeFoldConstantsTest struct {
	CompilerTestCaseBase

	late              bool
	useTypes          bool
	numRepetitions    int
	assumeGettersPure bool
}

// Java: PeepholeFoldConstantsTest()
func NewPeepholeFoldConstantsTest() *PeepholeFoldConstantsTest {
	c := &PeepholeFoldConstantsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c,
		NewTestExternsBuilder().
			AddArray().
			AddIterable().
			AddObject().
			AddUndefined().
			AddFunction().
			AddString().
			AddInfinity().
			AddNaN().
			Build())
	c.useTypes = true
	c.assumeGettersPure = false
	return c
}

func TestPeepholeFoldConstantsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPeepholeFoldConstantsTest() })
}

var peepholeFoldConstantsTest_LITERAL_OPERANDS = []string{
	"null",
	"undefined",
	"void 0",
	"true",
	"false",
	"!0",
	"!1",
	"0",
	"1",
	"''",
	"'123'",
	"'abc'",
	"'def'",
	"NaN",
	"Infinity",
	// TODO(nicksantos): Add more literals
	"-Infinity",
	// "({})",
	// "[]"
	// "[0]",
	// "Object",
	// "(function() {})"
}

// Java: PeepholeFoldConstantsTest.setUp()
func (c *PeepholeFoldConstantsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.DisableTypeCheck()
	c.late = false
	c.useTypes = true
	// Reduce this to 1 if we get better expression evaluators.
	c.numRepetitions = 2
	c.EnableNormalize()
	c.DisableCompareJsDoc()
}

// Java: PeepholeFoldConstantsTest.getProcessor(Compiler)
func (c *PeepholeFoldConstantsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPeepholeOptimizationsPass(
		compiler, c.GetName(), NewPeepholeFoldConstants(c.late, c.useTypes))
}

// Java: PeepholeFoldConstantsTest.getNumRepetitions()
func (c *PeepholeFoldConstantsTest) GetNumRepetitions() int {
	return c.numRepetitions
}

// Java: PeepholeFoldConstantsTest.getOptions()
func (c *PeepholeFoldConstantsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetAssumeGettersArePure(c.assumeGettersPure)
	return options
}

// Java: PeepholeFoldConstantsTest.testUndefinedComparison1()
func (c *PeepholeFoldConstantsTest) TestUndefinedComparison1() {
	c.Test("undefined == undefined", "true")
	c.Test("undefined == null", "true")
	c.Test("undefined == void 0", "true")
	c.Test("undefined == 0", "false")
	c.Test("undefined == 1", "false")
	c.Test("undefined == 'hi'", "false")
	c.Test("undefined == true", "false")
	c.Test("undefined == false", "false")
	c.Test("undefined === undefined", "true")
	c.Test("undefined === null", "false")
	c.Test("undefined === void 0", "true")
	c.TestSame("undefined == this")
	c.TestSame("undefined == x")
	c.Test("undefined != undefined", "false")
	c.Test("undefined != null", "false")
	c.Test("undefined != void 0", "false")
	c.Test("undefined != 0", "true")
	c.Test("undefined != 1", "true")
	c.Test("undefined != 'hi'", "true")
	c.Test("undefined != true", "true")
	c.Test("undefined != false", "true")
	c.Test("undefined !== undefined", "false")
	c.Test("undefined !== void 0", "false")
	c.Test("undefined !== null", "true")
	c.TestSame("undefined != this")
	c.TestSame("undefined != x")
	c.Test("undefined < undefined", "false")
	c.Test("undefined > undefined", "false")
	c.Test("undefined >= undefined", "false")
	c.Test("undefined <= undefined", "false")
	c.Test("0 < undefined", "false")
	c.Test("true > undefined", "false")
	c.Test("'hi' >= undefined", "false")
	c.Test("null <= undefined", "false")
	c.Test("undefined < 0", "false")
	c.Test("undefined > true", "false")
	c.Test("undefined >= 'hi'", "false")
	c.Test("undefined <= null", "false")
	c.Test("null == undefined", "true")
	c.Test("0 == undefined", "false")
	c.Test("1 == undefined", "false")
	c.Test("'hi' == undefined", "false")
	c.Test("true == undefined", "false")
	c.Test("false == undefined", "false")
	c.Test("null === undefined", "false")
	c.Test("void 0 === undefined", "true")
	c.Test("undefined == NaN", "false")
	c.Test("NaN == undefined", "false")
	c.Test("undefined == Infinity", "false")
	c.Test("Infinity == undefined", "false")
	c.Test("undefined == -Infinity", "false")
	c.Test("-Infinity == undefined", "false")
	c.Test("({}) == undefined", "false")
	c.Test("undefined == ({})", "false")
	c.Test("([]) == undefined", "false")
	c.Test("undefined == ([])", "false")
	c.Test("(/a/g) == undefined", "false")
	c.Test("undefined == (/a/g)", "false")
	c.Test("(function(){}) == undefined", "false")
	c.Test("undefined == (function(){})", "false")
	c.Test("undefined != NaN", "true")
	c.Test("NaN != undefined", "true")
	c.Test("undefined != Infinity", "true")
	c.Test("Infinity != undefined", "true")
	c.Test("undefined != -Infinity", "true")
	c.Test("-Infinity != undefined", "true")
	c.Test("({}) != undefined", "true")
	c.Test("undefined != ({})", "true")
	c.Test("([]) != undefined", "true")
	c.Test("undefined != ([])", "true")
	c.Test("(/a/g) != undefined", "true")
	c.Test("undefined != (/a/g)", "true")
	c.Test("(function(){}) != undefined", "true")
	c.Test("undefined != (function(){})", "true")
	c.TestSame("this == undefined")
	c.TestSame("x == undefined")
}

// Java: PeepholeFoldConstantsTest.testUndefinedComparison2()
func (c *PeepholeFoldConstantsTest) TestUndefinedComparison2() {
	c.Test("\"123\" !== void 0", "true")
	c.Test("\"123\" === void 0", "false")
	c.Test("void 0 !== \"123\"", "true")
	c.Test("void 0 === \"123\"", "false")
}

// Java: PeepholeFoldConstantsTest.testUndefinedComparison3()
func (c *PeepholeFoldConstantsTest) TestUndefinedComparison3() {
	c.Test("\"123\" !== undefined", "true")
	c.Test("\"123\" === undefined", "false")
	c.Test("undefined !== \"123\"", "true")
	c.Test("undefined === \"123\"", "false")
}

// Java: PeepholeFoldConstantsTest.testUndefinedComparison4()
func (c *PeepholeFoldConstantsTest) TestUndefinedComparison4() {
	c.Test("1 !== void 0", "true")
	c.Test("1 === void 0", "false")
	c.Test("null !== void 0", "true")
	c.Test("null === void 0", "false")
	c.Test("undefined !== void 0", "false")
	c.Test("undefined === void 0", "true")
}

// Java: PeepholeFoldConstantsTest.testNullComparison1()
func (c *PeepholeFoldConstantsTest) TestNullComparison1() {
	c.Test("null == undefined", "true")
	c.Test("null == null", "true")
	c.Test("null == void 0", "true")
	c.Test("null == 0", "false")
	c.Test("null == 1", "false")
	c.Test("null == 0n", "false")
	c.Test("null == 1n", "false")
	c.Test("null == 'hi'", "false")
	c.Test("null == true", "false")
	c.Test("null == false", "false")
	c.Test("null === undefined", "false")
	c.Test("null === null", "true")
	c.Test("null === void 0", "false")
	c.TestSame("null === x")
	c.TestSame("null == this")
	c.TestSame("null == x")
	c.Test("null != undefined", "false")
	c.Test("null != null", "false")
	c.Test("null != void 0", "false")
	c.Test("null != 0", "true")
	c.Test("null != 1", "true")
	c.Test("null != 0n", "true")
	c.Test("null != 1n", "true")
	c.Test("null != 'hi'", "true")
	c.Test("null != true", "true")
	c.Test("null != false", "true")
	c.Test("null !== undefined", "true")
	c.Test("null !== void 0", "true")
	c.Test("null !== null", "false")
	c.TestSame("null != this")
	c.TestSame("null != x")
	c.Test("null < null", "false")
	c.Test("null > null", "false")
	c.Test("null >= null", "true")
	c.Test("null <= null", "true")
	c.Test("0 < null", "false")
	c.Test("0 > null", "false")
	c.Test("0 >= null", "true")
	c.Test("0n < null", "false")
	c.Test("0n > null", "false")
	c.Test("0n >= null", "true")
	c.Test("true > null", "true")
	c.Test("'hi' < null", "false")
	c.Test("'hi' >= null", "false")
	c.Test("null <= null", "true")
	c.Test("null < 0", "false")
	c.Test("null < 0n", "false")
	c.Test("null > true", "false")
	c.Test("null < 'hi'", "false")
	c.Test("null >= 'hi'", "false")
	c.Test("null <= null", "true")
	c.Test("null == null", "true")
	c.Test("0 == null", "false")
	c.Test("1 == null", "false")
	c.Test("'hi' == null", "false")
	c.Test("true == null", "false")
	c.Test("false == null", "false")
	c.Test("null === null", "true")
	c.Test("void 0 === null", "false")
	c.Test("null == NaN", "false")
	c.Test("NaN == null", "false")
	c.Test("null == Infinity", "false")
	c.Test("Infinity == null", "false")
	c.Test("null == -Infinity", "false")
	c.Test("-Infinity == null", "false")
	c.Test("({}) == null", "false")
	c.Test("null == ({})", "false")
	c.Test("([]) == null", "false")
	c.Test("null == ([])", "false")
	c.Test("(/a/g) == null", "false")
	c.Test("null == (/a/g)", "false")
	c.Test("(function(){}) == null", "false")
	c.Test("null == (function(){})", "false")
	c.Test("null != NaN", "true")
	c.Test("NaN != null", "true")
	c.Test("null != Infinity", "true")
	c.Test("Infinity != null", "true")
	c.Test("null != -Infinity", "true")
	c.Test("-Infinity != null", "true")
	c.Test("({}) != null", "true")
	c.Test("null != ({})", "true")
	c.Test("([]) != null", "true")
	c.Test("null != ([])", "true")
	c.Test("(/a/g) != null", "true")
	c.Test("null != (/a/g)", "true")
	c.Test("(function(){}) != null", "true")
	c.Test("null != (function(){})", "true")
	c.TestSame("({a:f()}) == null")
	c.TestSame("null == ({a:f()})")
	c.TestSame("([f()]) == null")
	c.TestSame("null == ([f()])")
	c.TestSame("this == null")
	c.TestSame("x == null")
}

// Java: PeepholeFoldConstantsTest.testBooleanBooleanComparison()
func (c *PeepholeFoldConstantsTest) TestBooleanBooleanComparison() {
	c.TestSame("!x == !y")
	c.TestSame("!x < !y")
	c.TestSame("!x !== !y")
	c.TestSame("!x == !x")
	// foldable
	c.TestSame("!x < !x")
	// foldable
	c.TestSame("!x !== !x")
	// foldable
}

// Java: PeepholeFoldConstantsTest.testBooleanNumberComparison()
func (c *PeepholeFoldConstantsTest) TestBooleanNumberComparison() {
	c.TestSame("!x == +y")
	c.TestSame("!x <= +y")
	c.Test("!x !== +y", "true")
}

// Java: PeepholeFoldConstantsTest.testNumberBooleanComparison()
func (c *PeepholeFoldConstantsTest) TestNumberBooleanComparison() {
	c.TestSame("+x == !y")
	c.TestSame("+x <= !y")
	c.Test("+x === !y", "false")
}

// Java: PeepholeFoldConstantsTest.testBooleanStringComparison()
func (c *PeepholeFoldConstantsTest) TestBooleanStringComparison() {
	c.TestSame("!x == '' + y")
	c.TestSame("!x <= '' + y")
	c.Test("!x !== '' + y", "true")
}

// Java: PeepholeFoldConstantsTest.testStringBooleanComparison()
func (c *PeepholeFoldConstantsTest) TestStringBooleanComparison() {
	c.TestSame("'' + x == !y")
	c.TestSame("'' + x <= !y")
	c.Test("'' + x === !y", "false")
}

// Java: PeepholeFoldConstantsTest.testNumberNumberComparison()
func (c *PeepholeFoldConstantsTest) TestNumberNumberComparison() {
	c.Test("1 > 1", "false")
	c.Test("2 == 3", "false")
	c.Test("3.6 === 3.6", "true")
	c.TestSame("+x > +y")
	c.TestSame("+x == +y")
	c.TestSame("+x === +y")
	c.TestSame("+x == +x")
	c.TestSame("+x === +x")
	c.TestSame("+x > +x")
	// foldable
}

// Java: PeepholeFoldConstantsTest.testStringStringComparison()
func (c *PeepholeFoldConstantsTest) TestStringStringComparison() {
	c.Test("'a' < 'b'", "true")
	c.Test("'a' <= 'b'", "true")
	c.Test("'a' > 'b'", "false")
	c.Test("'a' >= 'b'", "false")
	c.Test("+'a' < +'b'", "false")
	c.TestSame("typeof a < 'a'")
	c.TestSame("'a' >= typeof a")
	c.Test("typeof a < typeof a", "false")
	c.Test("typeof a >= typeof a", "true")
	c.Test("typeof 3 > typeof 4", "false")
	c.Test("typeof function() {} < typeof function() {}", "false")
	c.Test("'a' == 'a'", "true")
	c.Test("'b' != 'a'", "true")
	c.TestSame("'undefined' == typeof a")
	c.TestSame("typeof a != 'number'")
	c.TestSame("'undefined' == typeof a")
	c.TestSame("'undefined' == typeof a")
	c.Test("typeof a == typeof a", "true")
	c.Test("'a' === 'a'", "true")
	c.Test("'b' !== 'a'", "true")
	c.Test("typeof a === typeof a", "true")
	c.Test("typeof a !== typeof a", "false")
	c.TestSame("'' + x <= '' + y")
	c.TestSame("'' + x != '' + y")
	c.TestSame("'' + x === '' + y")
	c.TestSame("'' + x <= '' + x")
	// potentially foldable
	c.TestSame("'' + x != '' + x")
	// potentially foldable
	c.TestSame("'' + x === '' + x")
	// potentially foldable
}

// Java: PeepholeFoldConstantsTest.testNumberStringComparison()
func (c *PeepholeFoldConstantsTest) TestNumberStringComparison() {
	c.Test("1 < '2'", "true")
	c.Test("2 > '1'", "true")
	c.Test("123 > '34'", "true")
	c.Test("NaN >= 'NaN'", "false")
	c.Test("1 == '2'", "false")
	c.Test("1 != '1'", "false")
	c.Test("NaN == 'NaN'", "false")
	c.Test("1 === '1'", "false")
	c.Test("1 !== '1'", "true")
	c.TestSame("+x > '' + y")
	c.TestSame("+x == '' + y")
	c.Test("+x !== '' + y", "true")
}

// Java: PeepholeFoldConstantsTest.testStringNumberComparison()
func (c *PeepholeFoldConstantsTest) TestStringNumberComparison() {
	c.Test("'1' < 2", "true")
	c.Test("'2' > 1", "true")
	c.Test("'123' > 34", "true")
	c.Test("'NaN' < NaN", "false")
	c.Test("'1' == 2", "false")
	c.Test("'1' != 1", "false")
	c.Test("'NaN' == NaN", "false")
	c.Test("'1' === 1", "false")
	c.Test("'1' !== 1", "true")
	c.TestSame("'' + x < +y")
	c.TestSame("'' + x == +y")
	c.Test("'' + x === +y", "false")
}

// Java: PeepholeFoldConstantsTest.testBigIntNumberComparison()
func (c *PeepholeFoldConstantsTest) TestBigIntNumberComparison() {
	c.Test("1n < 2", "true")
	c.Test("1n > 2", "false")
	c.Test("1n == 1", "true")
	c.Test("1n == 2", "false")
	// comparing with decimals is allowed
	c.Test("1n < 1.1", "true")
	c.Test("1n < 1.9", "true")
	c.Test("1n < 0.9", "false")
	c.Test("-1n < -1.1", "false")
	c.Test("-1n < -1.9", "false")
	c.Test("-1n < -0.9", "true")
	c.Test("1n > 1.1", "false")
	c.Test("1n > 0.9", "true")
	c.Test("-1n > -1.1", "true")
	c.Test("-1n > -0.9", "false")
	// Don't fold unsafely large numbers because there might be floating-point error
	var maxSafeInt int64 = 9007199254740991
	c.Test("0n > "+jutil.LongToString(maxSafeInt), "false")
	c.Test("0n < "+jutil.LongToString(maxSafeInt), "true")
	c.Test("0n > "+jutil.LongToString(-maxSafeInt), "true")
	c.Test("0n < "+jutil.LongToString(-maxSafeInt), "false")
	c.TestSame("0n > " + jutil.LongToString(maxSafeInt+1))
	c.TestSame("0n < " + jutil.LongToString(maxSafeInt+1))
	c.TestSame("0n > " + jutil.LongToString(-(maxSafeInt + 1)))
	c.TestSame("0n < " + jutil.LongToString(-(maxSafeInt + 1)))
	// comparing with Infinity is allowed
	c.Test("1n < Infinity", "true")
	c.Test("1n > Infinity", "false")
	c.Test("1n < -Infinity", "false")
	c.Test("1n > -Infinity", "true")
	// null is interpreted as 0 when comparing with bigint
	c.Test("1n < null", "false")
	c.Test("1n > null", "true")
}

// Java: PeepholeFoldConstantsTest.testBigIntStringComparison()
func (c *PeepholeFoldConstantsTest) TestBigIntStringComparison() {
	c.Test("1n < '2'", "true")
	c.Test("1n <= '2'", "true")
	c.Test("1n > '2'", "false")
	c.Test("1n >= '2'", "false")
	c.Test("1n < '1'", "false")
	c.Test("1n <= '1'", "true")
	c.Test("1n > '1'", "false")
	c.Test("1n >= '1'", "true")
	c.Test("1n == '1'", "true")
	c.Test("1n != '1'", "false")
	c.Test("2n > '1'", "true")
	c.Test("123n > '34'", "true")
	c.Test("1n == '2'", "false")
	c.Test("1n === '1'", "false")
	c.Test("1n !== '1'", "true")

	c.Test("10n == '  10  '", "true")
	c.Test("10n < '  20  '", "true")
	c.Test("16n == '0x10'", "true")
	c.Test("2n == '0b10'", "true")
	c.Test("8n == '0o10'", "true")
	c.Test("-1n == '-1'", "true")
	c.Test("1n == '+1'", "true")

	// Invalid BigInt strings (StringToBigInt returns undefined) evaluate to false for all
	// relational comparisons
	c.Test("1n < 'foo'", "false")
	c.Test("1n > 'foo'", "false")
	c.Test("1n <= 'foo'", "false")
	c.Test("1n >= 'foo'", "false")
	c.Test("1n == 'foo'", "false")
	c.Test("1n != 'foo'", "true")

	c.Test("1n < '1e2'", "false")
	c.Test("1n > '1e2'", "false")
	c.Test("1n <= '1e2'", "false")
	c.Test("1n >= '1e2'", "false")
	c.Test("1n == '1e2'", "false")
	c.Test("1n != '1e2'", "true")

	c.Test("1n < '1.5'", "false")
	c.Test("1n > '1.5'", "false")
	c.Test("1n <= '1.5'", "false")
	c.Test("1n >= '1.5'", "false")
	c.Test("1n == '1.5'", "false")
	c.Test("1n != '1.5'", "true")

	c.Test("1n < 'Infinity'", "false")
	c.Test("1n > 'Infinity'", "false")
	c.Test("1n <= 'Infinity'", "false")
	c.Test("1n >= 'Infinity'", "false")

	c.Test("1n < '-Infinity'", "false")
	c.Test("1n > '-Infinity'", "false")
	c.Test("1n <= '-Infinity'", "false")
	c.Test("1n >= '-Infinity'", "false")
}

// Java: PeepholeFoldConstantsTest.testStringBigIntComparison()
func (c *PeepholeFoldConstantsTest) TestStringBigIntComparison() {
	c.Test("'2' < 1n", "false")
	c.Test("'2' <= 1n", "false")
	c.Test("'2' > 1n", "true")
	c.Test("'2' >= 1n", "true")
	c.Test("'1' < 1n", "false")
	c.Test("'1' <= 1n", "true")
	c.Test("'1' > 1n", "false")
	c.Test("'1' >= 1n", "true")
	c.Test("'1' == 1n", "true")
	c.Test("'1' != 1n", "false")
	c.Test("'1' < 2n", "true")
	c.Test("'123' > 34n", "true")
	c.Test("'1' == 2n", "false")
	c.Test("'1' === 1n", "false")
	c.Test("'1' !== 1n", "true")

	c.Test("'  10  ' == 10n", "true")
	c.Test("'  20  ' > 10n", "true")
	c.Test("'0x10' == 16n", "true")
	c.Test("'0b10' == 2n", "true")
	c.Test("'0o10' == 8n", "true")
	c.Test("'-1' == -1n", "true")
	c.Test("'+1' == 1n", "true")

	c.Test("'foo' < 1n", "false")
	c.Test("'foo' > 1n", "false")
	c.Test("'foo' <= 1n", "false")
	c.Test("'foo' >= 1n", "false")
	c.Test("'foo' == 1n", "false")
	c.Test("'foo' != 1n", "true")

	c.Test("'1e2' < 1n", "false")
	c.Test("'1e2' > 1n", "false")
	c.Test("'1e2' <= 1n", "false")
	c.Test("'1e2' >= 1n", "false")

	c.Test("'1.5' < 1n", "false")
	c.Test("'1.5' > 1n", "false")
	c.Test("'1.5' <= 1n", "false")
	c.Test("'1.5' >= 1n", "false")
}

// Java: PeepholeFoldConstantsTest.testBigIntEqualityWithNullOrUndefined()
func (c *PeepholeFoldConstantsTest) TestBigIntEqualityWithNullOrUndefined() {
	c.Test("null == 1n", "false")
	c.Test("undefined == 1n", "false")
	c.Test("1n == null", "false")
	c.Test("1n == undefined", "false")
	c.Test("null != 1n", "true")
	c.Test("undefined != 1n", "true")
	c.Test("1n != null", "true")
	c.Test("1n != undefined", "true")
}

// Java: PeepholeFoldConstantsTest.testBigIntEqualityWithVariables()
func (c *PeepholeFoldConstantsTest) TestBigIntEqualityWithVariables() {
	c.TestSame("x == 1n")
	c.TestSame("1n == y")
	c.TestSame("x != 1n")
	c.TestSame("1n != y")
}

// Java: PeepholeFoldConstantsTest.testNaNComparison()
func (c *PeepholeFoldConstantsTest) TestNaNComparison() {
	c.Test("NaN < 1", "false")
	c.Test("NaN <= 1", "false")
	c.Test("NaN > 1", "false")
	c.Test("NaN >= 1", "false")
	c.Test("NaN < 1n", "false")
	c.Test("NaN <= 1n", "false")
	c.Test("NaN > 1n", "false")
	c.Test("NaN >= 1n", "false")
	c.Test("NaN < NaN", "false")
	c.Test("NaN >= NaN", "false")
	c.Test("NaN == NaN", "false")
	c.Test("NaN === NaN", "false")
	c.Test("NaN < null", "false")
	c.Test("null >= NaN", "false")
	c.Test("NaN == null", "false")
	c.Test("null != NaN", "true")
	c.Test("null === NaN", "false")
	c.Test("NaN < undefined", "false")
	c.Test("undefined >= NaN", "false")
	c.Test("NaN == undefined", "false")
	c.Test("undefined != NaN", "true")
	c.Test("undefined === NaN", "false")
	c.TestSame("NaN < x")
	c.TestSame("x >= NaN")
	c.TestSame("NaN == x")
	c.TestSame("x != NaN")
	c.Test("NaN === x", "false")
	c.Test("x !== NaN", "true")
	c.TestSame("NaN == foo()")
}

// Java: PeepholeFoldConstantsTest.testObjectComparison1()
func (c *PeepholeFoldConstantsTest) TestObjectComparison1() {
	c.Test("!new Date()", "false")
	c.Test("!!new Date()", "true")
	c.Test("new Date() == null", "false")
	c.Test("new Date() == undefined", "false")
	c.Test("new Date() != null", "true")
	c.Test("new Date() != undefined", "true")
	c.Test("null == new Date()", "false")
	c.Test("undefined == new Date()", "false")
	c.Test("null != new Date()", "true")
	c.Test("undefined != new Date()", "true")
}

// Java: PeepholeFoldConstantsTest.testUnaryOps()
func (c *PeepholeFoldConstantsTest) TestUnaryOps() {
	// Running on just changed code results in an exception on only the first invocation. Don't
	// repeat because it confuses the exception verification.
	c.numRepetitions = 1
	// These cases are handled by PeepholeRemoveDeadCode.
	c.TestSame("!foo()")
	c.TestSame("~foo()")
	c.TestSame("-foo()")
	// These cases are handled here.
	c.Test("a=!true", "a=false")
	c.Test("a=!10", "a=false")
	c.Test("a=!false", "a=true")
	c.TestSame("a=!foo()")
	c.Test("a=-0", "a=-0.0")
	c.Test("a=-(0)", "a=-0.0")
	c.TestSame("a=-Infinity")
	c.Test("a=-NaN", "a=NaN")
	c.TestSame("a=-foo()")
	c.Test("a=~~0", "a=0")
	c.Test("a=~~10", "a=10")
	c.Test("a=~-7", "a=6")
	c.Test("a=+true", "a=1")
	c.Test("a=+10", "a=10")
	c.Test("a=+false", "a=0")
	c.TestSame("a=+foo()")
	c.TestSame("a=+f")
	c.Test("a=+(f?true:false)", "a=+(f?1:0)")
	// TODO(johnlenz): foldable
	c.Test("a=+0", "a=0")
	c.Test("a=+Infinity", "a=Infinity")
	c.Test("a=+NaN", "a=NaN")
	c.Test("a=+-7", "a=-7")
	c.Test("a=+.5", "a=.5")
	c.Test("a=~0xffffffff", "a=0")
	c.Test("a=~~0xffffffff", "a=-1")
	c.TestSame("a=~.5", PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND)
}

// Java: PeepholeFoldConstantsTest.testUnaryOpsWithBigInt()
func (c *PeepholeFoldConstantsTest) TestUnaryOpsWithBigInt() {
	c.Test("-(1n)", "-1n")
	c.Test("- -1n", "1n")
	c.Test("!1n", "false")
	c.Test("~0n", "-1n")
}

// Java: PeepholeFoldConstantsTest.testUnaryOpsStringCompare()
func (c *PeepholeFoldConstantsTest) TestUnaryOpsStringCompare() {
	c.TestSame("a = -1")
	c.Test("a = ~0", "a = -1")
	c.Test("a = ~1", "a = -2")
	c.Test("a = ~101", "a = -102")
}

// Java: PeepholeFoldConstantsTest.testFoldLogicalOp()
func (c *PeepholeFoldConstantsTest) TestFoldLogicalOp() {
	c.Test("x = true && x", "x = x")
	c.Test("x = [foo()] && x", "x = ([foo()],x)")
	c.Test("x = false && x", "x = false")
	c.Test("x = true || x", "x = true")
	c.Test("x = false || x", "x = x")
	c.Test("x = 0 && x", "x = 0")
	c.Test("x = 3 || x", "x = 3")
	c.Test("x = 0n && x", "x = 0n")
	c.Test("x = 3n || x", "x = 3n")
	c.Test("x = false || 0", "x = 0")
	// unfoldable, because the right-side may be the result
	c.Test("a = x && true", "a=x && true")
	c.Test("a = x && false", "a=x && false")
	c.Test("a = x || 3", "a=x || 3")
	c.Test("a = x || false", "a=x || false")
	c.Test("a = b ? c : x || false", "a=b ? c:x || false")
	c.Test("a = b ? x || false : c", "a=b ? x || false:c")
	c.Test("a = b ? c : x && true", "a=b ? c:x && true")
	c.Test("a = b ? x && true : c", "a=b ? x && true:c")
	// folded, but not here.
	c.TestSame("a = x || false ? b : c")
	c.TestSame("a = x && true ? b : c")
	c.Test("x = foo() || true || bar()", "x = foo() || true")
	c.Test("x = foo() || true && bar()", "x = foo() || bar()")
	c.Test("x = foo() || false && bar()", "x = foo() || false")
	c.Test("x = foo() && false && bar()", "x = foo() && false")
	c.Test("x = foo() && false || bar()", "x = (foo() && false,bar())")
	c.Test("x = foo() || false || bar()", "x = foo() || bar()")
	c.Test("x = foo() && true && bar()", "x = foo() && bar()")
	c.Test("x = foo() || true || bar()", "x = foo() || true")
	c.Test("x = foo() && false && bar()", "x = foo() && false")
	c.Test("x = foo() && 0 && bar()", "x = foo() && 0")
	c.Test("x = foo() && 1 && bar()", "x = foo() && bar()")
	c.Test("x = foo() || 0 || bar()", "x = foo() || bar()")
	c.Test("x = foo() || 1 || bar()", "x = foo() || 1")
	c.Test("x = foo() && 0n && bar()", "x = foo() && 0n")
	c.Test("x = foo() && 1n && bar()", "x = foo() && bar()")
	c.Test("x = foo() || 0n || bar()", "x = foo() || bar()")
	c.Test("x = foo() || 1n || bar()", "x = foo() || 1n")
	c.TestSame("x = foo() || bar() || baz()")
	c.TestSame("x = foo() && bar() && baz()")
	c.Test("0 || b()", "b()")
	c.Test("1 && b()", "b()")
	c.Test("a() && (1 && b())", "a() && b()")
	c.Test("(a() && 1) && b()", "a() && b()")
	c.Test("(x || '') || y;", "x || y")
	c.Test("false || (x || '');", "x || ''")
	c.Test("(x && 1) && y;", "x && y")
	c.Test("true && (x && 1);", "x && 1")
	// Really not foldable, because it would change the type of the
	// expression if foo() returns something truthy but not true.
	// Cf. FoldConstants.tryFoldAndOr().
	// An example would be if foo() is 1 (truthy) and bar() is 0 (falsey):
	// (1 && true) || 0 == true
	// 1 || 0 == 1, but true =/= 1
	c.TestSame("x = foo() && true || bar()")
	c.TestSame("foo() && true || bar()")
}

// Java: PeepholeFoldConstantsTest.testFoldLogicalOp2()
func (c *PeepholeFoldConstantsTest) TestFoldLogicalOp2() {
	c.Test("x = function(){} && x", "x = x")
	c.Test("x = true && function(){}", "x = function(){}")
	c.Test("x = [(function(){alert(x)})()] && x", "x = ([(function(){alert(x)})()],x)")
}

// Java: PeepholeFoldConstantsTest.testFoldNullishCoalesce()
func (c *PeepholeFoldConstantsTest) TestFoldNullishCoalesce() {
	// fold if left is null/undefined
	c.Test("null ?? 1", "1")
	c.Test("undefined ?? false", "false")
	c.Test("(a(), null) ?? 1", "(a(), null, 1)")
	c.Test("x = [foo()] ?? x", "x = [foo()]")
	// short circuit on all non nullish LHS
	c.Test("x = false ?? x", "x = false")
	c.Test("x = true ?? x", "x = true")
	c.Test("x = 0 ?? x", "x = 0")
	c.Test("x = 3 ?? x", "x = 3")
	// unfoldable, because the right-side may be the result
	c.TestSame("a = x ?? true")
	c.TestSame("a = x ?? false")
	c.TestSame("a = x ?? 3")
	c.TestSame("a = b ? c : x ?? false")
	c.TestSame("a = b ? x ?? false : c")
	// folded, but not here.
	c.TestSame("a = x ?? false ? b : c")
	c.TestSame("a = x ?? true ? b : c")
	c.TestSame("x = foo() ?? true ?? bar()")
	c.Test("x = foo() ?? (true && bar())", "x = foo() ?? bar()")
	c.TestSame("x = (foo() || false) ?? bar()")
	c.Test("a() ?? (1 ?? b())", "a() ?? 1")
	c.Test("(a() ?? 1) ?? b()", "a() ?? 1 ?? b()")
}

// Java: PeepholeFoldConstantsTest.testFoldOptChain()
func (c *PeepholeFoldConstantsTest) TestFoldOptChain() {
	// can't fold when optional part may execute
	c.TestSame("a = x?.y")
	c.TestSame("a = x?.()")
	// fold args of optional call
	c.Test("x = foo() ?. (true && bar())", "x = foo() ?.(bar())")
	c.Test("a() ?. (1 ?? b())", "a() ?. (1)")
	c.Test("({a})?.a.b.c.d()?.x.y.z", "a.b.c.d()?.x.y.z")
	// potential optimization
	c.TestSame("x = undefined?.y")
	// `x = void 0;`
}

// Java: PeepholeFoldConstantsTest.testFoldOptChain_nonNullReceiver()
func (c *PeepholeFoldConstantsTest) TestFoldOptChain_nonNullReceiver() {
	// Constant folding on literals with optional chaining
	c.Test("x = 'hello'?.length", "x = 5")
	c.Test("x = ''?.length", "x = 0")
	c.Test("x = [1, 2, 3]?.length", "x = 3")
	c.Test("x = []?.length", "x = 0")
	c.Test("x = ([10, 20])?.[0]", "x = 10")
	c.Test("x = ([10, 20])?.[1]", "x = 20")
	c.Test("x = 'abcdef'?.[1]", "x = 'b'")
	c.Test("x = ({a: 1})?.a", "x = 1")
	c.Test("x = ({a: 1})?.['a']", "x = 1")

	// Guard cases / unknown receivers must remain untouched
	c.TestSame("x = str?.length")
	c.TestSame("x = arr?.length")
	c.TestSame("x = obj?.a")
	c.TestSame("x = obj?.[key]")
	c.TestSame("x = fn?.(arg)")
	c.TestSame("x = [foo()]?.length") // side-effecting array
}

// Java: PeepholeFoldConstantsTest.testFoldOptChain_nullishReceiver()
func (c *PeepholeFoldConstantsTest) TestFoldOptChain_nullishReceiver() {
	// Current compiler behavior in PeepholeFoldConstants (unfolded / preserved)
	c.TestSame("x = null?.y")
	c.TestSame("x = undefined?.y")
	c.TestSame("x = (void 0)?.y")
	c.TestSame("x = null?.[0]")
	c.TestSame("x = undefined?.[foo()]")
	c.TestSame("x = null?.(foo())")
	c.TestSame("x = (foo(), null)?.y")
}

// Java: PeepholeFoldConstantsTest.testFoldNullishCoalesce_advanced()
func (c *PeepholeFoldConstantsTest) TestFoldNullishCoalesce_advanced() {
	// Basic nullish LHS folding
	c.Test("null ?? 1", "1")
	c.Test("undefined ?? false", "false")
	c.Test("(void 0) ?? 'default'", "'default'")
	c.Test("(a(), null) ?? 1", "(a(), null, 1)")

	// Non-nullish LHS folding
	c.Test("x = false ?? x", "x = false")
	c.Test("x = true ?? x", "x = true")
	c.Test("x = 0 ?? x", "x = 0")
	c.Test("x = '' ?? x", "x = ''")
	c.Test("x = 'hello' ?? x", "x = 'hello'")
	c.Test("x = 3 ?? x", "x = 3")
	c.Test("x = [1, 2] ?? x", "x = [1, 2]")
	c.Test("x = ({a: 1}) ?? x", "x = ({a: 1})")
	c.Test("x = [foo()] ?? x", "x = [foo()]")

	// Guard cases: RHS cannot be folded if LHS is unknown
	c.TestSame("a = x ?? true")
	c.TestSame("a = x ?? false")
	c.TestSame("a = x ?? 3")
	c.TestSame("a = x ?? null")
	c.TestSame("a = (foo() ? null : 1) ?? 2")
}

// Java: PeepholeFoldConstantsTest.testFoldLogicalAssignments()
func (c *PeepholeFoldConstantsTest) TestFoldLogicalAssignments() {
	// Under normalization, logical assignments normalize to short-circuiting conditional assigns
	c.TestSame("x || (x = y)")
	c.TestSame("x && (x = y)")
	c.TestSame("x ?? (x = y)")
	c.TestSame("x || (x = false)")
	c.TestSame("x && (x = true)")
	c.TestSame("x ?? (x = null)")
	c.TestSame("obj[foo()] || (obj[foo()] = y)")
}

// Java: PeepholeFoldConstantsTest.testBatchC_templateLiteralsAndSpread()
func (c *PeepholeFoldConstantsTest) TestBatchC_templateLiteralsAndSpread() {
	// OPP-012: Template literal substitutions and guards
	c.Test("`a${'b'}c`", "`abc`")
	c.Test("`a${123}c`", "`a123c`")
	c.Test("`a${true}c`", "`atruec`")
	c.Test("`a${false}c`", "`afalsec`")
	c.Test("`${''}a${''}`", "`a`")
	c.TestSame("tag`a${'b'}c`")
	c.TestSame("`a${foo()}c`")
	c.TestSame("`a${x}c`")
	c.TestSame("`$${''}{`")

	// OPP-013: Object spread constant flattening and guards
	c.Test("x = {...{}}", "x = {}")
	c.Test("x = {a, ...{}, b}", "x = {a, b}")
	c.Test("x = {...{a, b}, c, ...{d, e}}", "x = {a, b, c, d, e}")
	c.Test("x = {...{...{a}, b}, c}", "x = {a, b, c}")
	c.TestSame("x = {...obj}")
	c.TestSame("x = {...foo()}")
	c.TestSame("x = {...{get a() { return 1; }}}")
	c.TestSame("x = {...{set a(v) { }}}")

	// OPP-014: Array spread constant flattening and guards
	c.Test("x = [...[]]", "x = []")
	c.Test("x = [0, ...[], 1]", "x = [0, 1]")
	c.Test("x = [...[0, 1], 2, ...[3, 4]]", "x = [0, 1, 2, 3, 4]")
	c.Test("x = [...[...[0], 1], 2]", "x = [0, 1, 2]")
	c.Test("foo(...[0], 1)", "foo(0, 1)")
	c.Test("foo([...[...[0], 1], 2])", "foo([0, 1, 2])")
	c.TestSame("x = [...iter]")
	c.TestSame("foo(...iter)")
	c.TestSame("x = [...{}]")
	c.TestSame("x = {...[]}")
}

// Java: PeepholeFoldConstantsTest.testBatchE_arithmeticNeutralElements()
func (c *PeepholeFoldConstantsTest) TestBatchE_arithmeticNeutralElements() {
	// OPP-020: Arithmetic Neutral Elements & Strength Reductions
	// Literal double negation
	c.Test("x = -(-4)", "x = 4")
	c.Test("x = -(-4n)", "x = 4n")

	// Guard cases: expressions that must not fold or are preserved
	c.TestSame("x = foo() ** 0")     // side-effectful base
	c.TestSame("x = bigIntVal ** 1") // cannot mix BigInt with number 1
	c.TestSame("x = (a + b) ** 1")   // non-constant exponentiation
	c.TestSame("x = (a + b) ** 0")
	c.TestSame("x = -(-(a + b))")
	c.TestSame("x = a % 1")
}

// Java: PeepholeFoldConstantsTest.testBatchE_bitwiseNeutralAndAbsorptionIdentities()
func (c *PeepholeFoldConstantsTest) TestBatchE_bitwiseNeutralAndAbsorptionIdentities() {
	// OPP-021: Bitwise Neutral, Idempotent & Absorption Identities
	// Fold constant bitwise operations
	c.Test("x = 1.5 | 0", "x = 1")
	c.Test("x = 4294967295 | 0", "x = -1")
	c.Test("x = -1 & 0", "x = 0")
	c.Test("x = 0 & -1", "x = 0")
	c.Test("x = ~0", "x = -1")
	c.Test("x = ~-7", "x = 6")

	// Guard cases: non-constant / side-effectful expressions
	c.TestSame("x = foo() ^ foo()") // side effects cannot be dropped
	c.TestSame("x = foo() & 0")     // side effects cannot be dropped
	c.TestSame("x = (a | b) | 0")
	c.TestSame("x = (a & b) & -1")
	c.TestSame("x = (a ^ b) ^ 0")
	c.TestSame("x = (a >> b) >> 0")
	c.TestSame("x = (a >>> b) >>> 0")
	c.TestSame("x = ~~ (a | b)")
	c.TestSame("x = (a + b) ^ (a + b)")
	c.TestSame("x = (a + b) & 0")
	c.TestSame("x = y | 0")
}

// Java: PeepholeFoldConstantsTest.testFoldBitwiseOp()
func (c *PeepholeFoldConstantsTest) TestFoldBitwiseOp() {
	c.Test("x = 1 & 1", "x = 1")
	c.Test("x = 1 & 2", "x = 0")
	c.Test("x = 3 & 1", "x = 1")
	c.Test("x = 3 & 3", "x = 3")
	c.Test("x = 1 | 1", "x = 1")
	c.Test("x = 1 | 2", "x = 3")
	c.Test("x = 3 | 1", "x = 3")
	c.Test("x = 3 | 3", "x = 3")
	c.Test("x = 1 ^ 1", "x = 0")
	c.Test("x = 1 ^ 2", "x = 3")
	c.Test("x = 3 ^ 1", "x = 2")
	c.Test("x = 3 ^ 3", "x = 0")
	c.Test("x = -1 & 0", "x = 0")
	c.Test("x = 0 & -1", "x = 0")
	c.Test("x = 1 & 4", "x = 0")
	c.Test("x = 2 & 3", "x = 2")
	// make sure we fold only when we are supposed to -- not when doing so would
	// lose information or when it is performed on nonsensical arguments.
	c.Test("x = 1 & 1.1", "x = 1")
	c.Test("x = 1.1 & 1", "x = 1")
	c.Test("x = 1 & 3000000000", "x = 0")
	c.Test("x = 3000000000 & 1", "x = 0")
	// Try some cases with | as well
	c.Test("x = 1 | 4", "x = 5")
	c.Test("x = 1 | 3", "x = 3")
	c.Test("x = 1 | 1.1", "x = 1")
	c.TestSame("x = 1 | 3E9")
	// these cases look strange because bitwise OR converts unsigned numbers to be signed
	c.Test("x = 1 | 3000000001", "x = -1294967295")
	c.Test("x = 4294967295 | 0", "x = -1")
}

// Java: PeepholeFoldConstantsTest.testFoldBitwiseOp2()
func (c *PeepholeFoldConstantsTest) TestFoldBitwiseOp2() {
	c.Test("x = y & 1 & 1", "x = y & 1")
	c.Test("x = y & 1 & 2", "x = y & 0")
	c.Test("x = y & 3 & 1", "x = y & 1")
	c.Test("x = 3 & y & 1", "x = y & 1")
	c.Test("x = y & 3 & 3", "x = y & 3")
	c.Test("x = 3 & y & 3", "x = y & 3")
	c.Test("x = y | 1 | 1", "x = y | 1")
	c.Test("x = y | 1 | 2", "x = y | 3")
	c.Test("x = y | 3 | 1", "x = y | 3")
	c.Test("x = 3 | y | 1", "x = y | 3")
	c.Test("x = y | 3 | 3", "x = y | 3")
	c.Test("x = 3 | y | 3", "x = y | 3")
	c.Test("x = y ^ 1 ^ 1", "x = y ^ 0")
	c.Test("x = y ^ 1 ^ 2", "x = y ^ 3")
	c.Test("x = y ^ 3 ^ 1", "x = y ^ 2")
	c.Test("x = 3 ^ y ^ 1", "x = y ^ 2")
	c.Test("x = y ^ 3 ^ 3", "x = y ^ 0")
	c.Test("x = 3 ^ y ^ 3", "x = y ^ 0")
	c.Test("x = Infinity | NaN", "x=0")
	c.Test("x = 12 | NaN", "x=12")
}

// Java: PeepholeFoldConstantsTest.testFoldBitwiseOpWithBigInt()
func (c *PeepholeFoldConstantsTest) TestFoldBitwiseOpWithBigInt() {
	c.Test("x = 1n & 1n", "x = 1n")
	c.Test("x = 1n & 2n", "x = 0n")
	c.Test("x = 3n & 1n", "x = 1n")
	c.Test("x = 3n & 3n", "x = 3n")
	c.Test("x = 1n | 1n", "x = 1n")
	c.Test("x = 1n | 2n", "x = 3n")
	c.Test("x = 1n | 3n", "x = 3n")
	c.Test("x = 3n | 1n", "x = 3n")
	c.Test("x = 3n | 3n", "x = 3n")
	c.Test("x = 1n | 4n", "x = 5n")
	c.Test("x = 1n ^ 1n", "x = 0n")
	c.Test("x = 1n ^ 2n", "x = 3n")
	c.Test("x = 3n ^ 1n", "x = 2n")
	c.Test("x = 3n ^ 3n", "x = 0n")
	c.Test("x = -1n & 0n", "x = 0n")
	c.Test("x = 0n & -1n", "x = 0n")
	c.Test("x = 1n & 4n", "x = 0n")
	c.Test("x = 2n & 3n", "x = 2n")
	c.Test("x = 1n & 3000000000n", "x = 0n")
	c.Test("x = 3000000000n & 1n", "x = 0n")
	// bitwise OR does not affect the sign of a bigint
	c.Test("x = 1n | 3000000001n", "x = 3000000001n")
	c.Test("x = 4294967295n | 0n", "x = 4294967295n")
	c.Test("x = y & 1n & 1n", "x = y & 1n")
	c.Test("x = y & 1n & 2n", "x = y & 0n")
	c.Test("x = y & 3n & 1n", "x = y & 1n")
	c.Test("x = 3n & y & 1n", "x = y & 1n")
	c.Test("x = y & 3n & 3n", "x = y & 3n")
	c.Test("x = 3n & y & 3n", "x = y & 3n")
	c.Test("x = y | 1n | 1n", "x = y | 1n")
	c.Test("x = y | 1n | 2n", "x = y | 3n")
	c.Test("x = y | 3n | 1n", "x = y | 3n")
	c.Test("x = 3n | y | 1n", "x = y | 3n")
	c.Test("x = y | 3n | 3n", "x = y | 3n")
	c.Test("x = 3n | y | 3n", "x = y | 3n")
	c.Test("x = y ^ 1n ^ 1n", "x = y ^ 0n")
	c.Test("x = y ^ 1n ^ 2n", "x = y ^ 3n")
	c.Test("x = y ^ 3n ^ 1n", "x = y ^ 2n")
	c.Test("x = 3n ^ y ^ 1n", "x = y ^ 2n")
	c.Test("x = y ^ 3n ^ 3n", "x = y ^ 0n")
	c.Test("x = 3n ^ y ^ 3n", "x = y ^ 0n")
}

// Java: PeepholeFoldConstantsTest.testFoldingMixTypesLate()
func (c *PeepholeFoldConstantsTest) TestFoldingMixTypesLate() {
	c.late = true
	c.DisableNormalize()
	c.Test("x = x + '2'", "x+='2'")
	c.Test("x = +x + +'2'", "x = +x + 2")
	c.Test("x = x - '2'", "x-=2")
	c.Test("x = x ^ '2'", "x^=2")
	c.Test("x = '2' ^ x", "x^=2")
	c.Test("x = '2' & x", "x&=2")
	c.Test("x = '2' | x", "x|=2")
	c.Test("x = '2' | y", "x=2|y")
	c.Test("x = y | '2'", "x=y|2")
	c.Test("x = y | (a && '2')", "x=y|(a&&2)")
	c.Test("x = y | (a,'2')", "x=y|(a,2)")
	c.Test("x = y | (a?'1':'2')", "x=y|(a?1:2)")
	c.Test("x = y | ('x'?'1':'2')", "x=y|('x'?1:2)")
}

// Java: PeepholeFoldConstantsTest.testFoldingMixTypesEarly()
func (c *PeepholeFoldConstantsTest) TestFoldingMixTypesEarly() {
	c.late = false
	c.TestSame("x = x + '2'")
	c.Test("x = +x + +'2'", "x = +x + 2")
	c.Test("x = x - '2'", "x = x - 2")
	c.Test("x = x ^ '2'", "x = x ^ 2")
	c.Test("x = '2' ^ x", "x = 2 ^ x")
	c.Test("x = '2' & x", "x = 2 & x")
	c.Test("x = '2' | x", "x = 2 | x")
	c.Test("x = '2' | y", "x=2|y")
	c.Test("x = y | '2'", "x=y|2")
	c.Test("x = y | (a && '2')", "x=y|(a&&2)")
	c.Test("x = y | (a,'2')", "x=y|(a,2)")
	c.Test("x = y | (a?'1':'2')", "x=y|(a?1:2)")
	c.Test("x = y | ('x'?'1':'2')", "x=y|('x'?1:2)")
}

// Java: PeepholeFoldConstantsTest.testFoldingAdd1()
func (c *PeepholeFoldConstantsTest) TestFoldingAdd1() {
	c.Test("x = null + true", "x=1")
	c.TestSame("x = a + true")
	c.Test("x = '' + {}", "x = '[object Object]'")
	c.Test("x = [] + {}", "x = '[object Object]'")
	c.Test("x = {} + []", "x = '[object Object]'")
	c.Test("x = {} + ''", "x = '[object Object]'")
}

// Java: PeepholeFoldConstantsTest.testFoldingAdd2()
func (c *PeepholeFoldConstantsTest) TestFoldingAdd2() {
	c.Test("x = false + []", "x='false'")
	c.Test("x = [] + true", "x='true'")
	c.Test("NaN + []", "'NaN'")
}

// Java: PeepholeFoldConstantsTest.testFoldBitwiseOpStringCompare()
func (c *PeepholeFoldConstantsTest) TestFoldBitwiseOpStringCompare() {
	c.Test("x = -1 | 0", "x = -1")
}

// Java: PeepholeFoldConstantsTest.testFoldBitShifts()
func (c *PeepholeFoldConstantsTest) TestFoldBitShifts() {
	// Running on just changed code results in an exception on only the first invocation. Don't
	// repeat because it confuses the exception verification.
	c.numRepetitions = 1
	c.Test("x = 1 << 0", "x = 1")
	c.Test("x = -1 << 0", "x = -1")
	c.Test("x = 1 << 1", "x = 2")
	c.Test("x = 3 << 1", "x = 6")
	c.Test("x = 1 << 8", "x = 256")
	c.Test("x = 1 >> 0", "x = 1")
	c.Test("x = -1 >> 0", "x = -1")
	c.Test("x = 1 >> 1", "x = 0")
	c.Test("x = 2 >> 1", "x = 1")
	c.Test("x = 5 >> 1", "x = 2")
	c.Test("x = 127 >> 3", "x = 15")
	c.Test("x = 3 >> 1", "x = 1")
	c.Test("x = 3 >> 2", "x = 0")
	c.Test("x = 10 >> 1", "x = 5")
	c.Test("x = 10 >> 2", "x = 2")
	c.Test("x = 10 >> 5", "x = 0")
	c.Test("x = 10 >>> 1", "x = 5")
	c.Test("x = 10 >>> 2", "x = 2")
	c.Test("x = 10 >>> 5", "x = 0")
	c.Test("x = -1 >>> 1", "x = 2147483647")
	// 0x7fffffff
	c.Test("x = -1 >>> 0", "x = 4294967295")
	// 0xffffffff
	c.Test("x = -2 >>> 0", "x = 4294967294")
	// 0xfffffffe
	c.Test("x = 0x90000000 >>> 28", "x = 9")
	c.Test("x = 0xffffffff << 0", "x = -1")
	c.Test("x = 0xffffffff << 4", "x = -16")
	c.TestSame("1 << 32")
	c.TestSame("1 << -1")
	c.TestSame("1 >> 32")
	c.TestSame("1.5 << 0", PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND)
	c.TestSame("1 << .5", PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND)
	c.TestSame("1.5 >>> 0", PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND)
	c.TestSame("1 >>> .5", PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND)
	c.TestSame("1.5 >> 0", PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND)
	c.TestSame("1 >> .5", PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND)
}

// Java: PeepholeFoldConstantsTest.testFoldBitShiftsStringCompare()
func (c *PeepholeFoldConstantsTest) TestFoldBitShiftsStringCompare() {
	c.Test("x = -1 << 1", "x = -2")
	c.Test("x = -1 << 8", "x = -256")
	c.Test("x = -1 >> 1", "x = -1")
	c.Test("x = -2 >> 1", "x = -1")
	c.Test("x = -1 >> 0", "x = -1")
}

// Java: PeepholeFoldConstantsTest.testStringAdd()
func (c *PeepholeFoldConstantsTest) TestStringAdd() {
	c.Test("x = 'a' + 'bc'", "x = 'abc'")
	c.Test("x = 'a' + 5", "x = 'a5'")
	c.Test("x = 5 + 'a'", "x = '5a'")
	c.Test("x = 'a' + 5n", "x = 'a5'")
	c.Test("x = 5n + 'a'", "x = '5a'")
	c.Test("x = 'a' + ''", "x = 'a'")
	c.Test("x = 'a' + foo()", "x = 'a'+foo()")
	c.Test("x = foo() + 'a' + 'b'", "x = foo()+'ab'")
	c.Test("x = (foo() + 'a') + 'b'", "x = foo()+'ab'")
	// believe it!
	c.Test("x = foo() + 'a' + 'b' + 'cd' + bar()", "x = foo()+'abcd'+bar()")
	c.Test("x = foo() + 2 + 'b'", "x = foo()+2+\"b\"")
	// don't fold!
	c.Test("x = foo() + 'a' + 2", "x = foo()+\"a2\"")
	c.Test("x = '' + null", "x = 'null'")
	c.Test("x = true + '' + false", "x = 'truefalse'")
	c.Test("x = '' + []", "x = ''")
	c.Test("x = foo() + 'a' + 1 + 1", "x = foo() + 'a11'")
	c.Test("x = 1 + 1 + 'a'", "x = '2a'")
	c.Test("x = 1 + 1 + 'a'", "x = '2a'")
	c.Test("x = 'a' + (1 + 1)", "x = 'a2'")
	c.Test("x = '_' + p1 + '_' + ('' + p2)", "x = '_' + p1 + '_' + p2")
	c.Test("x = 'a' + ('_' + 1 + 1)", "x = 'a_11'")
	c.Test("x = 'a' + ('_' + 1) + 1", "x = 'a_11'")
	c.Test("x = 1 + (p1 + '_') + ('' + p2)", "x = 1 + (p1 + '_') + p2")
	c.Test("x = 1 + p1 + '_' + ('' + p2)", "x = 1 + p1 + '_' + p2")
	c.Test("x = 1 + 'a' + p1", "x = '1a' + p1")
	c.Test("x = (p1 + (p2 + 'a')) + 'b'", "x = (p1 + (p2 + 'ab'))")
	c.Test("'a' + ('b' + p1) + 1", "'ab' + p1 + 1")
	c.Test("x = 'a' + ('b' + p1 + 'c')", "x = 'ab' + (p1 + 'c')")
	c.TestSame("x = 'a' + (4 + p1 + 'a')")
	c.TestSame("x = p1 / 3 + 4")
	c.TestSame("foo() + 3 + 'a' + foo()")
	c.TestSame("x = 'a' + ('b' + p1 + p2)")
	c.TestSame("x = 1 + ('a' + p1)")
	c.TestSame("x = p1 + '' + p2")
	c.TestSame("x = 'a' + (1 + p1)")
	c.TestSame("x = (p2 + 'a') + (1 + p1)")
	c.TestSame("x = (p2 + 'a') + (1 + p1 + p2)")
	c.TestSame("x = (p2 + 'a') + (1 + (p1 + p2))")
}

// Java: PeepholeFoldConstantsTest.testStringAdd_identity()
func (c *PeepholeFoldConstantsTest) TestStringAdd_identity() {
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()
	c.FoldStringTypes("x + ''", "x")
	c.FoldStringTypes("'' + x", "x")
}

// Java: PeepholeFoldConstantsTest.testIssue821()
func (c *PeepholeFoldConstantsTest) TestIssue821() {
	c.TestSame("var a =(Math.random()>0.5? '1' : 2 ) + 3 + 4;")
	c.TestSame("var a = ((Math.random() ? 0 : 1) || (Math.random()>0.5? '1' : 2 )) + 3 + 4;")
}

// Java: PeepholeFoldConstantsTest.testFoldConstructor()
func (c *PeepholeFoldConstantsTest) TestFoldConstructor() {
	c.Test("x = this[new String('a')]", "x = this['a']")
	c.Test("x = ob[new String(12)]", "x = ob['12']")
	c.Test("x = ob[new String(false)]", "x = ob['false']")
	c.Test("x = ob[new String(null)]", "x = ob['null']")
	c.Test("x = 'a' + new String('b')", "x = 'ab'")
	c.Test("x = 'a' + new String(23)", "x = 'a23'")
	c.Test("x = 2 + new String(1)", "x = '21'")
	c.TestSame("x = ob[new String(a)]")
	c.TestSame("x = new String('a')")
	c.TestSame("x = (new String('a'))[3]")
}

// Java: PeepholeFoldConstantsTest.testFoldArithmetic()
func (c *PeepholeFoldConstantsTest) TestFoldArithmetic() {
	c.Test("x = 10 + 20", "x = 30")
	c.Test("x = 2 / 4", "x = 0.5")
	c.Test("x = 2.25 & 3", "x = 2")
	c.TestSame("z = x & y")
	c.TestSame("x = y & 5")
	c.TestSame("x = 1 / 0")
	c.Test("x = 3 % 2", "x = 1")
	c.Test("x = 3 % -2", "x = 1")
	c.Test("x = -1 % 3", "x = -1")
	c.TestSame("x = 1 % 0")
	c.Test("x = 2 ** 3", "x = 8")
	c.Test("x = 2 ** -3", "x = 0.125")
	c.TestSame("x = 2 ** 55")
	// backs off folding because 2 ** 55 is too large
	c.TestSame("x = 3 ** -1")
	// backs off because 3**-1 is shorter than 0.3333333333333333
}

// Java: PeepholeFoldConstantsTest.testFoldArithmetic2()
func (c *PeepholeFoldConstantsTest) TestFoldArithmetic2() {
	c.TestSame("x = y + 10 + 20")
	c.TestSame("x = y / 2 / 4")
	c.Test("x = y & 2.25 & 3", "x = y & 2")
	c.TestSame("x = y & 2.25 & z & 3")
	c.TestSame("z = x &y")
	c.TestSame("x = y & 5")
	c.Test("x = y + (z & 24 & 60 & 60 & 1000)", "x = y + (z & 8)")
}

// Java: PeepholeFoldConstantsTest.testFoldArithmetic3()
func (c *PeepholeFoldConstantsTest) TestFoldArithmetic3() {
	c.Test("x = null | undefined", "x = 0")
	c.Test("x = null | 1", "x = 1")
	c.Test("x = (null - 1)| 2", "x = -1")
	c.Test("x = (null + 1) | 2", "x = 3")
	c.Test("x = null ** 0", "x = 1")
	c.Test("x = (-0) ** 3", "x = -0")
}

// Java: PeepholeFoldConstantsTest.testFoldArithmeticInfinity()
func (c *PeepholeFoldConstantsTest) TestFoldArithmeticInfinity() {
	c.Test("x=-Infinity-2", "x=-Infinity")
	c.Test("x=Infinity-2", "x=Infinity")
	c.Test("x=Infinity*5", "x=Infinity")
	c.Test("x = Infinity ** 2", "x = Infinity")
	c.Test("x = Infinity ** -2", "x = 0")
}

// Java: PeepholeFoldConstantsTest.testFoldArithmeticStringComp()
func (c *PeepholeFoldConstantsTest) TestFoldArithmeticStringComp() {
	c.Test("x = 10 - 20", "x = -10")
}

// Java: PeepholeFoldConstantsTest.testNoFoldArithmeticWithSideEffects()
func (c *PeepholeFoldConstantsTest) TestNoFoldArithmeticWithSideEffects() {
	// can't fold this to "x = y & 6.75" because you can't remove the "sideEffects()" call
	c.TestSame("x = y & 2.25 & (sideEffects(), 3)")
}

// Java: PeepholeFoldConstantsTest.testFoldBigIntArithmetic()
func (c *PeepholeFoldConstantsTest) TestFoldBigIntArithmetic() {
	c.Test("x = 1n + 2n", "x = 3n")
	c.Test("x = 1n - 2n", "x = -1n")
	c.Test("x = 2n * 3n", "x = 6n")
	c.Test("x = 6n / 2n", "x = 3n")
	c.Test("x = 3n % 2n", "x = 1n")
	c.Test("x = 2n ** 3n", "x = 8n")
	// The compiler is not designed to fold expressions with an exponent > 2147483647
	c.Test("x = 1n ** 2147483647n", "x = 1n")
	c.TestSame("x = 1n ** 2147483648n")
	c.Test("x = y & 2n & 3n", "x = y & 2n")
	// TODO b/361826515: Optimize associative bigint operations
	c.TestSame("x = y * 2n * z * 3n")
	c.TestSame("x = y + 2n + z + 3n")
}

// Java: PeepholeFoldConstantsTest.testNoFoldBigIntArithmeticWithSideEffects()
func (c *PeepholeFoldConstantsTest) TestNoFoldBigIntArithmeticWithSideEffects() {
	// can't fold this to "x = y * 6.75" because you can't remove the "sideEffects()" call
	c.TestSame("x = y * 2n * (sideEffects(), 3n)")
}

// Java: PeepholeFoldConstantsTest.testFoldComparison()
func (c *PeepholeFoldConstantsTest) TestFoldComparison() {
	c.Test("x = 0 == 0", "x = true")
	c.Test("x = 1 == 2", "x = false")
	c.Test("x = 0n == 0n", "x = true")
	c.Test("x = 1n == 2n", "x = false")
	c.Test("x = 'abc' == 'def'", "x = false")
	c.Test("x = 'abc' == 'abc'", "x = true")
	c.Test("x = \"\" == ''", "x = true")
	c.Test("x = foo() == bar()", "x = foo()==bar()")
	c.Test("x = 1 != 0", "x = true")
	c.Test("x = 1 != 1", "x = false")
	c.Test("x = 1n != 0n", "x = true")
	c.Test("x = 1n != 1n", "x = false")
	c.Test("x = 'abc' != 'def'", "x = true")
	c.Test("x = 'a' != 'a'", "x = false")
	c.Test("x = 1 < 20", "x = true")
	c.Test("x = 3 < 3", "x = false")
	c.Test("x = 10 > 1.0", "x = true")
	c.Test("x = 10 > 10.25", "x = false")
	c.Test("x = 1 <= 1", "x = true")
	c.Test("x = 1 <= 0", "x = false")
	c.Test("x = 0 >= 0", "x = true")
	c.Test("x = -1 >= 9", "x = false")
	c.Test("x = 1n < 20n", "x = true")
	c.Test("x = 3n < 3n", "x = false")
	c.Test("x = 10n > 1n", "x = true")
	c.Test("x = 10n > 10n", "x = false")
	c.Test("x = 1n <= 1n", "x = true")
	c.Test("x = 1n <= 0n", "x = false")
	c.Test("x = 0n >= 0n", "x = true")
	c.Test("x = -1n >= 9n", "x = false")
	c.TestSame("x = y == y")
	c.Test("x = y < y", "x = false")
	c.Test("x = y > y", "x = false")
	c.Test("x = true == true", "x = true")
	c.Test("x = false == false", "x = true")
	c.Test("x = false == null", "x = false")
	c.Test("x = false == true", "x = false")
	c.Test("x = true == null", "x = false")
	c.Test("0 == 0", "true")
	c.Test("1 == 2", "false")
	c.Test("0n == 0n", "true")
	c.Test("1n == 2n", "false")
	c.Test("'abc' == 'def'", "false")
	c.Test("'abc' == 'abc'", "true")
	c.Test("\"\" == ''", "true")
	c.TestSame("foo() == bar()")
	c.Test("1 != 0", "true")
	c.Test("1 != 1", "false")
	c.Test("1n != 0n", "true")
	c.Test("1n != 1n", "false")
	c.Test("'abc' != 'def'", "true")
	c.Test("'a' != 'a'", "false")
	c.Test("1 < 20", "true")
	c.Test("3 < 3", "false")
	c.Test("10 > 1.0", "true")
	c.Test("10 > 10.25", "false")
	c.TestSame("x == x")
	c.Test("x < x", "false")
	c.Test("x > x", "false")
	c.Test("1 <= 1", "true")
	c.Test("1 <= 0", "false")
	c.Test("0 >= 0", "true")
	c.Test("-1 >= 9", "false")
	c.Test("1n < 20n", "true")
	c.Test("3n < 3n", "false")
	c.Test("10n > 1n", "true")
	c.Test("10n > 10n", "false")
	c.Test("1n <= 1n", "true")
	c.Test("1n <= 0n", "false")
	c.Test("0n >= 0n", "true")
	c.Test("-1n >= 9n", "false")
	c.Test("true == true", "true")
	c.Test("false == null", "false")
	c.Test("false == true", "false")
	c.Test("true == null", "false")
}

// ===, !== comparison tests
// Java: PeepholeFoldConstantsTest.testFoldComparison2()
func (c *PeepholeFoldConstantsTest) TestFoldComparison2() {
	c.Test("x = 0 === 0", "x = true")
	c.Test("x = 1 === 2", "x = false")
	c.Test("x = 0n === 0n", "x = true")
	c.Test("x = 1n === 2n", "x = false")
	c.Test("x = 'abc' === 'def'", "x = false")
	c.Test("x = 'abc' === 'abc'", "x = true")
	c.Test("x = \"\" === ''", "x = true")
	c.Test("x = foo() === bar()", "x = foo()===bar()")
	c.Test("x = 1 !== 0", "x = true")
	c.Test("x = 1 !== 1", "x = false")
	c.Test("x = 1n !== 0n", "x = true")
	c.Test("x = 1n !== 1n", "x = false")
	c.Test("x = 'abc' !== 'def'", "x = true")
	c.Test("x = 'a' !== 'a'", "x = false")
	c.Test("x = y === y", "x = y===y")
	c.Test("x = true === true", "x = true")
	c.Test("x = false === false", "x = true")
	c.Test("x = false === null", "x = false")
	c.Test("x = false === true", "x = false")
	c.Test("x = true === null", "x = false")
	c.Test("0 === 0", "true")
	c.Test("1 === 2", "false")
	c.Test("0n === 0n", "true")
	c.Test("1n === 2n", "false")
	c.Test("'abc' === 'def'", "false")
	c.Test("'abc' === 'abc'", "true")
	c.Test("\"\" === ''", "true")
	c.TestSame("foo() === bar()")
	c.Test("1 === '1'", "false")
	c.Test("1 === true", "false")
	c.Test("1 !== '1'", "true")
	c.Test("1 !== true", "true")
	c.Test("1 !== 0", "true")
	c.Test("'abc' !== 'def'", "true")
	c.Test("'a' !== 'a'", "false")
	c.TestSame("x === x")
	c.Test("true === true", "true")
	c.Test("false === null", "false")
	c.Test("false === true", "false")
	c.Test("true === null", "false")
}

// Java: PeepholeFoldConstantsTest.testFoldComparison3()
func (c *PeepholeFoldConstantsTest) TestFoldComparison3() {
	c.Test("x = !1 == !0", "x = false")
	c.Test("x = !0 == !0", "x = true")
	c.Test("x = !1 == !1", "x = true")
	c.Test("x = !1 == null", "x = false")
	c.Test("x = !1 == !0", "x = false")
	c.Test("x = !0 == null", "x = false")
	c.Test("!0 == !0", "true")
	c.Test("!1 == null", "false")
	c.Test("!1 == !0", "false")
	c.Test("!0 == null", "false")
	c.Test("x = !0 === !0", "x = true")
	c.Test("x = !1 === !1", "x = true")
	c.Test("x = !1 === null", "x = false")
	c.Test("x = !1 === !0", "x = false")
	c.Test("x = !0 === null", "x = false")
	c.Test("!0 === !0", "true")
	c.Test("!1 === null", "false")
	c.Test("!1 === !0", "false")
	c.Test("!0 === null", "false")
}

// Java: PeepholeFoldConstantsTest.testFoldComparison4()
func (c *PeepholeFoldConstantsTest) TestFoldComparison4() {
	c.TestSame("[] == false")
	// true
	c.TestSame("[] == true")
	// false
	c.TestSame("[0] == false")
	// true
	c.TestSame("[0] == true")
	// false
	c.TestSame("[1] == false")
	// false
	c.TestSame("[1] == true")
	// true
	c.TestSame("({}) == false")
	// false
	c.TestSame("({}) == true")
	// true
}

// Java: PeepholeFoldConstantsTest.testFoldGetElem1()
func (c *PeepholeFoldConstantsTest) TestFoldGetElem1() {
	// Running on just changed code results in an exception on only the first invocation. Don't
	// repeat because it confuses the exception verification.
	c.numRepetitions = 1

	c.Test("x = [,10][0]", "x = void 0")
	c.Test("x = [10, 20][0]", "x = 10")
	c.Test("x = [10, 20][1]", "x = 20")

	c.TestSame("x = [10, 20][0.5]", PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR)
	c.Test("x = [10, 20][-1]", "x = void 0;")
	c.Test("x = [10, 20][2]", "x = void 0;")

	// TODO(b/538156673): Fix PeepholeFoldConstants to not fold out-of-bounds GETELEM on arrays with
	// side-effect elements
	c.TestSame("x = [foo(), 0][1]")
	c.Test("x = [0, foo()][1]", "x = foo()")
	c.TestSame("x = [0, foo()][0]")
	c.TestSame("x = [foo(), 0][-1]")
	c.TestSame("x = [0, foo()][-1]")
	c.TestSame("x = [foo(), 0][2]")
	c.TestSame("x = [0, foo()][2]")
	c.TestSame("for([1][0] in {});")
}

// Optional versions of the above `testFoldGetElem1` tests
// Java: PeepholeFoldConstantsTest.testFoldOptChainGetElem1()
func (c *PeepholeFoldConstantsTest) TestFoldOptChainGetElem1() {
	c.numRepetitions = 1
	c.Test("x = [,10]?.[0]", "x = void 0")
	c.Test("x = [10, 20]?.[0]", "x = 10")
	c.Test("x = [10, 20]?.[1]", "x = 20")

	c.TestSame("x = [10, 20]?.[0.5]", PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR)
	c.Test("x = [10, 20]?.[-1]", "x = void 0;")
	c.Test("x = [10, 20]?.[2]", "x = void 0;")

	c.TestSame("x = [foo(), 0]?.[1]")
	c.Test("x = [0, foo()]?.[1]", "x = foo()")
	c.TestSame("x = [0, foo()]?.[0]")
	c.TestSame("x = [foo(), 0]?.[-1]")
	c.TestSame("x = [0, foo()]?.[-1]")
	c.TestSame("x = [foo(), 0]?.[2]")
	c.TestSame("x = [0, foo()]?.[2]")
}

// Java: PeepholeFoldConstantsTest.testFoldGetElem2()
func (c *PeepholeFoldConstantsTest) TestFoldGetElem2() {
	// Running on just changed code results in an exception on only the first invocation. Don't
	// repeat because it confuses the exception verification.
	c.numRepetitions = 1
	c.Test("x = 'string'[5]", "x = 'g'")
	c.Test("x = 'string'[0]", "x = 's'")
	c.Test("x = 's'[0]", "x = 's'")
	c.TestSame("x = '\U0001f4a9'[0]")
	c.TestSame("x = 'string'[0.5]", PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR)
	c.Test("x = 'string'[-1]", "x = void 0;")
	c.Test("x = 'string'[6]", "x = void 0;")
}

// Optional versions of the above `testFoldGetElem2` tests
// Java: PeepholeFoldConstantsTest.testFoldOptChainGetElem2()
func (c *PeepholeFoldConstantsTest) TestFoldOptChainGetElem2() {
	// Running on just changed code results in an exception on only the first invocation. Don't
	// repeat because it confuses the exception verification.
	c.numRepetitions = 1
	c.Test("x = 'string'?.[5]", "x = 'g'")
	c.Test("x = 'string'?.[0]", "x = 's'")
	c.Test("x = 's'?.[0]", "x = 's'")
	c.TestSame("x = '\U0001f4a9'?.[0]")
	c.TestSame("x = 'string'?.[0.5]", PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR)
	c.Test("x = 'string'?.[-1]", "x = void 0;")
	c.Test("x = 'string'?.[6]", "x = void 0;")
}

// Java: PeepholeFoldConstantsTest.testFoldArrayLitSpreadGetElem()
func (c *PeepholeFoldConstantsTest) TestFoldArrayLitSpreadGetElem() {
	c.numRepetitions = 1
	c.Test("x = [...[0]][0]", "x = 0;")
	c.Test("x = [0, 1, ...[2, 3, 4]][3]", "x = 3;")
	c.Test("x = [...[0, 1], 2, ...[3, 4]][3]", "x = 3;")
	c.Test("x = [...[...[0, 1], 2, 3], 4][0]", "x = 0")
	c.Test("x = [...[...[0, 1], 2, 3], 4][3]", "x = 3")
	c.Test(c.Srcs("x = [...[]][100]"), c.Expected("x = void 0;"))
	c.Test(c.Srcs("x = [...[0]][100]"), c.Expected("x = void 0;"))
}

// Optional versions of the above `testFoldArrayLitSpreadGetElem` tests
// Java: PeepholeFoldConstantsTest.testFoldArrayLitSpreadOptChainGetElem()
func (c *PeepholeFoldConstantsTest) TestFoldArrayLitSpreadOptChainGetElem() {
	c.numRepetitions = 1
	c.Test("x = [...[0]]?.[0]", "x = 0;")
	c.Test("x = [0, 1, ...[2, 3, 4]]?.[3]", "x = 3;")
	c.Test("x = [...[0, 1], 2, ...[3, 4]]?.[3]", "x = 3;")
	c.Test("x = [...[...[0, 1], 2, 3], 4]?.[0]", "x = 0")
	c.Test("x = [...[...[0, 1], 2, 3], 4]?.[3]", "x = 3")
	c.Test(c.Srcs("x = [...[]]?.[100]"), c.Expected("x = void 0;"))
	c.Test(c.Srcs("x = [...[0]]?.[100]"), c.Expected("x = void 0;"))
}

// Java: PeepholeFoldConstantsTest.testDontFoldNonLiteralSpreadGetElem()
func (c *PeepholeFoldConstantsTest) TestDontFoldNonLiteralSpreadGetElem() {
	c.TestSame("x = [...iter][0];")
	c.TestSame("x = [0, 1, ...iter][2];")
	//  `...iter` could have side effects, so don't replace `x` with `0`
	c.TestSame("x = [0, 1, ...iter][0];")
}

// Java: PeepholeFoldConstantsTest.testFoldArraySpread()
func (c *PeepholeFoldConstantsTest) TestFoldArraySpread() {
	c.numRepetitions = 1
	c.Test("x = [...[]]", "x = []")
	c.Test("x = [0, ...[], 1]", "x = [0, 1]")
	c.Test("x = [...[0, 1], 2, ...[3, 4]]", "x = [0, 1, 2, 3, 4]")
	c.Test("x = [...[...[0], 1], 2]", "x = [0, 1, 2]")
	c.TestSame("[...[x]] = arr")
	c.Test("foo([...[...[0], 1], 2])", "foo([0, 1, 2])")
	c.TestSame("x = [1, ...[,], 2];")
	c.TestSame("x = [0, ...[,,], 3];")
	c.TestSame("x = [...[,]];")
}

// Java: PeepholeFoldConstantsTest.testFoldArrayLitSpreadInArg()
func (c *PeepholeFoldConstantsTest) TestFoldArrayLitSpreadInArg() {
	c.Test("foo(...[0], 1)", "foo(0, 1)")
	c.TestSame("foo(...[,]);")
	c.TestSame("foo(...[,,,,\"foo\"], 1)")
	c.TestSame("foo(...(false ? [0] : [1]))") // other opts need to fold the ternery first
}

// Java: PeepholeFoldConstantsTest.testFoldObjectLitSpreadGetProp()
func (c *PeepholeFoldConstantsTest) TestFoldObjectLitSpreadGetProp() {
	c.numRepetitions = 1
	c.Test("x = {...{a}}.a", "x = a;")
	c.Test("x = {a, b, ...{c, d, e}}.d", "x = d;")
	c.Test("x = {...{a, b}, c, ...{d, e}}.d", "x = d;")
	c.Test("x = {...{...{a, b}, c, d}, e}.a", "x = a")
	c.Test("x = {...{...{a, b}, c, d}, e}.d", "x = d")
}

// Java: PeepholeFoldConstantsTest.testDontFoldNonLiteralObjectSpreadGetProp_gettersImpure()
func (c *PeepholeFoldConstantsTest) TestDontFoldNonLiteralObjectSpreadGetProp_gettersImpure() {
	c.assumeGettersPure = false
	c.TestSame("x = {...obj}.a;")
	c.TestSame("x = {a, ...obj, c}.a;")
	c.TestSame("x = {a, ...obj, c}.c;")
}

// Java: PeepholeFoldConstantsTest.testDontFoldNonLiteralObjectSpreadGetProp_assumeGettersPure()
func (c *PeepholeFoldConstantsTest) TestDontFoldNonLiteralObjectSpreadGetProp_assumeGettersPure() {
	c.assumeGettersPure = true
	c.TestSame("x = {...obj}.a;")
	c.TestSame("x = {a, ...obj, c}.a;")
	c.Test("x = {a, ...obj, c}.c;", "x = c;")
	// We assume object spread has no side-effects.
}

// Java: PeepholeFoldConstantsTest.testFoldObjectSpread()
func (c *PeepholeFoldConstantsTest) TestFoldObjectSpread() {
	c.numRepetitions = 1
	c.Test("x = {...{}}", "x = {}")
	c.Test("x = {a, ...{}, b}", "x = {a, b}")
	c.Test("x = {...{a, b}, c, ...{d, e}}", "x = {a, b, c, d, e}")
	c.Test("x = {...{...{a}, b}, c}", "x = {a, b, c}")
	c.TestSame("({...{x}} = obj)")
}

// Java: PeepholeFoldConstantsTest.testDontFoldObjectSpread_withGetterSetter()
func (c *PeepholeFoldConstantsTest) TestDontFoldObjectSpread_withGetterSetter() {
	c.numRepetitions = 1
	// At runtime, object spread converts getters to regular data properties, and omits setters.
	// We back off on folding spread with getters & setters.
	// In theory, we could make the compiler smart enough to turn
	//   `{...{get a() { return 1; }}}` into `{a: 1}`, or `{...{set a(v) { }}}` into `{}`, but
	// it doesn't seem worth the complexity right now.
	c.TestSame("x = {...{get a() { return 1; }}}")
	c.TestSame("x = {...{set a(v) { }}}")
	// Test case with side-effects. If in the future we implement spread folding for getters/setters
	// we would need to preserve evaluation of the console.log.
	c.TestSame("x = {...{get a() { console.log('hi'); return 1; }}}")
}

// Java: PeepholeFoldConstantsTest.testDontFoldObjectSpread_withComputedGetterSetter()
func (c *PeepholeFoldConstantsTest) TestDontFoldObjectSpread_withComputedGetterSetter() {
	c.numRepetitions = 1
	c.TestSame("x = {...{get [a]() { return 1; }}}")
	c.TestSame("x = {...{set [a](v) { }}}")
}

// Java: PeepholeFoldConstantsTest.testDontFoldMixedObjectAndArraySpread()
func (c *PeepholeFoldConstantsTest) TestDontFoldMixedObjectAndArraySpread() {
	c.numRepetitions = 1
	c.TestSame("x = [...{}]")
	c.TestSame("x = {...[]}")
	c.Test("x = [a, ...[...{}]]", "x = [a, ...{}]")
	c.Test("x = {a, ...{...[]}}", "x = {a, ...[]}")
}

// Java: PeepholeFoldConstantsTest.testFoldComplex()
func (c *PeepholeFoldConstantsTest) TestFoldComplex() {
	c.Test("x = (3 / 1.0) + (1 * 2)", "x = 5")
	c.Test("x = (1 == 1.0) && foo() && true", "x = foo()&&true")
	c.Test("x = 'abc' + 5 + 10", "x = \"abc510\"")
}

// Java: PeepholeFoldConstantsTest.testFoldLeft()
func (c *PeepholeFoldConstantsTest) TestFoldLeft() {
	c.TestSame("(+x - 1) + 2")
	// not yet
	c.Test("(+x & 1) & 2", "+x & 0")
}

// Java: PeepholeFoldConstantsTest.testFoldArrayLength()
func (c *PeepholeFoldConstantsTest) TestFoldArrayLength() {
	// Can fold
	c.Test("x = [].length", "x = 0")
	c.Test("x = [1,2,3].length", "x = 3")
	c.Test("x = [a,b].length", "x = 2")

	// Not handled yet
	c.Test("x = [,,1].length", "x = 3")

	// Cannot fold
	c.Test("x = [foo(), 0].length", "x = [foo(),0].length")
	c.TestSame("x = y.length")
	c.TestSame("[1, 2].length = 0;")
	c.Test("[1, 2].length >>= 1;", "[1, 2].length = 1;")
	c.TestSame("[1, 2].length++;")
	c.TestSame("++[1, 2].length;")
	c.TestSame("[1, 2].length--;")
	c.TestSame("--[1, 2].length;")
	c.TestSame("x = [...'abc'].length")
	c.TestSame("x = [1, ...'ab', 3].length")
	c.TestSame("x = [...`abc`].length")
}

// Java: PeepholeFoldConstantsTest.testFoldStringLength()
func (c *PeepholeFoldConstantsTest) TestFoldStringLength() {
	// Can fold basic strings.
	c.Test("x = ''.length", "x = 0")
	c.Test("x = '123'.length", "x = 3")

	// Test Unicode escapes are accounted for.
	c.Test("x = '123\u01dc'.length", "x = 4")

	// Cannot fold when length is an lvalue
	c.TestSame("\"a\".length = 1;")
	c.Test("\"a\".length >>= 1;", "\"a\".length = 0;")
	c.TestSame("\"a\".length++;")
	c.TestSame("++\"a\".length;")
	c.TestSame("\"a\".length--;")
	c.TestSame("--\"a\".length;")
}

// Java: PeepholeFoldConstantsTest.testFoldTypeof()
func (c *PeepholeFoldConstantsTest) TestFoldTypeof() {
	c.Test("x = typeof 1", "x = \"number\"")
	c.Test("x = typeof 'foo'", "x = \"string\"")
	c.Test("x = typeof true", "x = \"boolean\"")
	c.Test("x = typeof false", "x = \"boolean\"")
	c.Test("x = typeof null", "x = \"object\"")
	c.Test("x = typeof undefined", "x = \"undefined\"")
	c.Test("x = typeof void 0", "x = \"undefined\"")
	c.Test("x = typeof []", "x = \"object\"")
	c.Test("x = typeof [1]", "x = \"object\"")
	c.Test("x = typeof [1,[]]", "x = \"object\"")
	c.Test("x = typeof {}", "x = \"object\"")
	c.Test("x = typeof function() {}", "x = 'function'")
	c.TestSame("x = typeof[1,[foo()]]")
	c.TestSame("x = typeof{bathwater:baby()}")
}

// Java: PeepholeFoldConstantsTest.testFoldInstanceOf()
func (c *PeepholeFoldConstantsTest) TestFoldInstanceOf() {
	// Non object types are never instances of anything.
	c.Test("64 instanceof Object", "false")
	c.Test("64 instanceof Number", "false")
	c.Test("'' instanceof Object", "false")
	c.Test("'' instanceof String", "false")
	c.Test("true instanceof Object", "false")
	c.Test("true instanceof Boolean", "false")
	c.Test("!0 instanceof Object", "false")
	c.Test("!0 instanceof Boolean", "false")
	c.Test("false instanceof Object", "false")
	c.Test("null instanceof Object", "false")
	c.Test("undefined instanceof Object", "false")
	c.Test("NaN instanceof Object", "false")
	c.Test("Infinity instanceof Object", "false")

	// Untagged template literals always evaluate to strings.
	c.Test("`` instanceof Object", "false")
	c.Test("`${[]}` instanceof Object", "false")
	c.Test("`${function(){}}` instanceof Object", "false")
	// Back off when substitutions have side-effects.
	c.TestSame("`${(console.log(0), [])}` instanceof Object")

	// Tagged template literals may evaluate to a non-string, though. Would require type information
	// to fold.
	c.TestSame("tag`${function(){}}` instanceof Object")

	// Array and object literals are known to be objects.
	c.Test("[] instanceof Object", "true")
	c.Test("({}) instanceof Object", "true")

	// These cases is foldable, but no handled currently.
	c.TestSame("new Foo() instanceof Object")
	// These would require type information to fold.
	c.TestSame("[] instanceof Foo")
	c.TestSame("({}) instanceof Foo")

	c.Test("(function() {}) instanceof Object", "true")

	// An unknown value should never be folded.
	c.TestSame("x instanceof Foo")
}

// Java: PeepholeFoldConstantsTest.testDivision()
func (c *PeepholeFoldConstantsTest) TestDivision() {
	// Make sure the 1/3 does not expand to 0.333333
	c.TestSame("print(1/3)")
	// Decimal form is preferable to fraction form when strings are the
	// same length.
	c.Test("print(1/2)", "print(0.5)")
}

// Java: PeepholeFoldConstantsTest.testAssignOpsLate()
func (c *PeepholeFoldConstantsTest) TestAssignOpsLate() {
	c.late = true
	c.DisableNormalize()
	c.Test("x=x+y", "x+=y")
	c.TestSame("x=y+x")
	c.Test("x=x*y", "x*=y")
	c.Test("x=y*x", "x*=y")
	c.Test("x.y=x.y+z", "x.y+=z")
	c.TestSame("next().x = next().x + 1")

	c.Test("x=x-y", "x-=y")
	c.TestSame("x=y-x")
	c.Test("x=x|y", "x|=y")
	c.Test("x=y|x", "x|=y")
	c.Test("x=x|y|z", "x|=y|z")
	c.TestSame("x=x&&y&&z")
	c.Test("x=x*y", "x*=y")
	c.Test("x=y*x", "x*=y")
	c.TestSame("x=foo()*x")
	c.TestSame("x=(x=5)*x")
	c.TestSame("x=(y++)*x")
	c.TestSame("x=foo()|x")
	c.TestSame("x=foo()&x")
	c.TestSame("x=foo()^x")
	c.Test("x=x*foo()", "x*=foo()")
	c.Test("x=x|foo()", "x|=foo()")
	c.Test("x=x&foo()", "x&=foo()")
	c.Test("x=x^foo()", "x^=foo()")
	c.Test("x=x+foo()", "x+=foo()")
	c.Test("x=x**y", "x**=y")
	c.TestSame("x=y**x")
	c.Test("x.y=x.y+z", "x.y+=z")
	c.TestSame("next().x = next().x + 1")
	// This is OK, really.
	c.Test("({a:1}).a = ({a:1}).a + 1", "({a:1}).a = 2")
}

// Java: PeepholeFoldConstantsTest.testAssignOpsEarly()
func (c *PeepholeFoldConstantsTest) TestAssignOpsEarly() {
	c.late = false
	c.TestSame("x=x+y")
	c.TestSame("x=y+x")
	c.TestSame("x=x*y")
	c.TestSame("x=y*x")
	c.TestSame("x.y=x.y+z")
	c.TestSame("next().x = next().x + 1")
	c.TestSame("x=x-y")
	c.TestSame("x=y-x")
	c.TestSame("x=x|y")
	c.TestSame("x=y|x")
	c.TestSame("x=x*y")
	c.TestSame("x=y*x")
	c.TestSame("x=x**y")
	c.TestSame("x=y**2")
	c.TestSame("x.y=x.y+z")
	c.TestSame("next().x = next().x + 1")
	// This is OK, really.
	c.Test("({a:1}).a = ({a:1}).a + 1", "({a:1}).a = 2")
}

// Java: PeepholeFoldConstantsTest.testUnfoldAssignOpsLate()
func (c *PeepholeFoldConstantsTest) TestUnfoldAssignOpsLate() {
	c.late = true
	c.DisableNormalize()
	c.TestSame("x+=y")
	c.TestSame("x*=y")
	c.TestSame("x.y+=z")
	c.TestSame("x-=y")
	c.TestSame("x|=y")
	c.TestSame("x*=y")
	c.TestSame("x**=y")
	c.TestSame("x.y+=z")
	c.TestSame("\"a\".length >>= 1;")
	c.TestSame("[1, 2].length >>= 1;")
}

// Java: PeepholeFoldConstantsTest.testUnfoldAssignOpsEarly()
func (c *PeepholeFoldConstantsTest) TestUnfoldAssignOpsEarly() {
	c.late = false
	c.Test("x+=y", "x=x+y")
	c.Test("x*=y", "x=x*y")
	c.Test("x.y+=z", "x.y=x.y+z")
	c.Test("x-=y", "x=x-y")
	c.Test("x|=y", "x=x|y")
	c.Test("x*=y", "x=x*y")
	c.Test("x**=y", "x=x**y")
	c.Test("x.y+=z", "x.y=x.y+z")
}

// Java: PeepholeFoldConstantsTest.testFoldAdd1()
func (c *PeepholeFoldConstantsTest) TestFoldAdd1() {
	c.Test("x=false+1", "x=1")
	c.Test("x=true+1", "x=2")
	c.Test("x=1+false", "x=1")
	c.Test("x=1+true", "x=2")
}

// Java: PeepholeFoldConstantsTest.testFoldLiteralNames()
func (c *PeepholeFoldConstantsTest) TestFoldLiteralNames() {
	c.Test("NaN == NaN", "false")
	c.Test("Infinity == Infinity", "true")
	c.Test("Infinity == NaN", "false")
	c.Test("undefined == NaN", "false")
	c.Test("undefined == Infinity", "false")
	c.Test("Infinity >= Infinity", "true")
	c.Test("NaN >= NaN", "false")
}

// Java: PeepholeFoldConstantsTest.testFoldLiteralsTypeMismatches()
func (c *PeepholeFoldConstantsTest) TestFoldLiteralsTypeMismatches() {
	c.Test("true == true", "true")
	c.Test("true == false", "false")
	c.Test("true == null", "false")
	c.Test("false == null", "false")
	// relational operators convert its operands
	c.Test("null <= null", "true")
	// 0 = 0
	c.Test("null >= null", "true")
	c.Test("null > null", "false")
	c.Test("null < null", "false")
	c.Test("false >= null", "true")
	// 0 = 0
	c.Test("false <= null", "true")
	c.Test("false > null", "false")
	c.Test("false < null", "false")
	c.Test("true >= null", "true")
	// 1 > 0
	c.Test("true <= null", "false")
	c.Test("true > null", "true")
	c.Test("true < null", "false")
	c.Test("true >= false", "true")
	// 1 > 0
	c.Test("true <= false", "false")
	c.Test("true > false", "true")
	c.Test("true < false", "false")
}

// Java: PeepholeFoldConstantsTest.testFoldLeftChildConcat()
func (c *PeepholeFoldConstantsTest) TestFoldLeftChildConcat() {
	c.TestSame("x +5 + \"1\"")
	c.Test("x+\"5\" + \"1\"", "x + \"51\"")
	// test("\"a\"+(c+\"b\")", "\"a\"+c+\"b\"");
	c.Test("\"a\"+(\"b\"+c)", "\"ab\"+c")
}

// Java: PeepholeFoldConstantsTest.testFoldLeftChildOp()
func (c *PeepholeFoldConstantsTest) TestFoldLeftChildOp() {
	c.Test("x & Infinity & 2", "x & 0")
	c.TestSame("x - Infinity - 2")
	// want "x-Infinity"
	c.TestSame("x - 1 + Infinity")
	c.TestSame("x - 2 + 1")
	c.TestSame("x - 2 + 3")
	c.TestSame("1 + x - 2 + 1")
	c.TestSame("1 + x - 2 + 3")
	c.TestSame("1 + x - 2 + 3 - 1")
	c.TestSame("f(x)-0")
	c.Test("x-0-0", "x-0")
	c.TestSame("x+2-2+2")
	c.TestSame("x+2-2+2-2")
	c.TestSame("x-2+2")
	c.TestSame("x-2+2-2")
	c.TestSame("x-2+2-2+2")
	c.TestSame("1+x-0-NaN")
	c.TestSame("1+f(x)-0-NaN")
	c.TestSame("1+x-0+NaN")
	c.TestSame("1+f(x)-0+NaN")
	c.TestSame("1+x+NaN")
	// unfoldable
	c.TestSame("x+2-2")
	// unfoldable
	c.TestSame("x+2")
	// nothing to do
	c.TestSame("x-2")
	// nothing to do
}

// Java: PeepholeFoldConstantsTest.testFoldSimpleArithmeticOp()
func (c *PeepholeFoldConstantsTest) TestFoldSimpleArithmeticOp() {
	c.TestSame("x|NaN")
	c.TestSame("NaN/y")
	c.TestSame("f(x)-0")
	c.TestSame("f(x)|1")
	c.TestSame("1|f(x)")
	c.TestSame("0+a+b")
	c.TestSame("0-a-b")
	c.TestSame("a+b-0")
	c.TestSame("(1+x)|NaN")
	c.TestSame("(1+f(x))|NaN")
	// don't fold side-effects
}

// Java: PeepholeFoldConstantsTest.testFoldLiteralsAsNumbers()
func (c *PeepholeFoldConstantsTest) TestFoldLiteralsAsNumbers() {
	c.Test("x/'12'", "x/12")
	c.Test("x/('12'+'6')", "x/126")
	c.Test("true*x", "1*x")
	c.Test("x/false", "x/0")
	// should we add an error check? :)
}

// Java: PeepholeFoldConstantsTest.testNotFoldBackToTrueFalse()
func (c *PeepholeFoldConstantsTest) TestNotFoldBackToTrueFalse() {
	c.late = false
	c.Test("!0", "true")
	c.Test("!1", "false")
	c.Test("!3", "false")
	c.late = true
	c.DisableNormalize()
	c.TestSame("!0")
	c.TestSame("!1")
	c.Test("!3", "false")
	c.TestSame("false")
	c.TestSame("true")
}

// Java: PeepholeFoldConstantsTest.testFoldBangConstants()
func (c *PeepholeFoldConstantsTest) TestFoldBangConstants() {
	c.Test("1 + !0", "2")
	c.Test("1 + !1", "1")
	c.Test("'a ' + !1", "'a false'")
	c.Test("'a ' + !0", "'a true'")
}

// Java: PeepholeFoldConstantsTest.testFoldMixed()
func (c *PeepholeFoldConstantsTest) TestFoldMixed() {
	c.Test("''+[1]", "'1'")
	c.Test("false+[]", "\"false\"")
}

// Java: PeepholeFoldConstantsTest.testFoldVoid()
func (c *PeepholeFoldConstantsTest) TestFoldVoid() {
	c.TestSame("void 0")
	c.Test("void 1", "void 0")
	c.Test("void x", "void 0")
	c.TestSame("void x()")
}

// Java: PeepholeFoldConstantsTest.testObjectLiteral()
func (c *PeepholeFoldConstantsTest) TestObjectLiteral() {
	c.Test("(!{})", "false")
	c.Test("(!{a:1})", "false")
	c.TestSame("(!{a:foo()})")
	c.TestSame("(!{'a':foo()})")
}

// Java: PeepholeFoldConstantsTest.testArrayLiteral()
func (c *PeepholeFoldConstantsTest) TestArrayLiteral() {
	c.Test("(![])", "false")
	c.Test("(![1])", "false")
	c.Test("(![a])", "false")
	c.TestSame("(![foo()])")
}

// Java: PeepholeFoldConstantsTest.testIssue601()
func (c *PeepholeFoldConstantsTest) TestIssue601() {
	c.TestSame("'\\v' == 'v'")
	c.TestSame("'v' == '\\v'")
	c.TestSame("'\\u000B' == '\\v'")
}

// Java: PeepholeFoldConstantsTest.testFoldObjectLiteralRef1()
func (c *PeepholeFoldConstantsTest) TestFoldObjectLiteralRef1() {
	// Leave extra side-effects in place
	c.TestSame("var x = ({a:foo(),b:bar()}).a")
	c.TestSame("var x = ({a:1,b:bar()}).a")
	c.TestSame("function f() { return {b:foo(), a:2}.a; }")

	// on the LHS the object act as a temporary leave it in place.
	c.TestSame("({a:x}).a = 1")
	c.Test("({a:x}).a += 1", "({a:x}).a = x + 1")
	c.TestSame("({a:x}).a ++")
	c.TestSame("({a:x}).a --")

	// Getters should not be inlined.
	c.TestSame("({get a() {return this}}).a")
	c.TestSame("({get a() {return this}})?.a")

	// Except, if we can see that the getter function never references 'this'.
	c.Test("({get a() {return 0}}).a", "(function() {return 0})()")
	c.Test("({get a() {return 0}})?.a", "(function() {return 0})()")
	c.Test("({get a() {return 0}})?.a.b", "(function() {return 0})().b")

	// It's okay to inline functions, as long as they're not immediately called.
	// (For tests where they are immediately called, see testFoldObjectLiteral_X)
	c.Test("({a:function(){return this}}).a", "(function(){return this})")
	c.Test("({a:function(){return this}})?.a", "(function(){return this})")

	// It's also okay to inline functions that are immediately called, so long as we know for
	// sure the function doesn't reference 'this'.
	c.Test("({a:function(){return 0}}).a()", "(function(){return 0})()")
	c.Test("({a:function(){return 0}})?.a()", "(function(){return 0})()")
	c.Test("({a:function(){return 0}})?.a().b", "(function(){return 0})().b")

	// Don't inline setters.
	c.TestSame("({set a(b) {return this}}).a")
	c.TestSame("({set a(b) {this._a = b}}).a")

	// Don't inline if there are side-effects.
	c.TestSame("({[foo()]: 1,   a: 0}).a")
	c.TestSame("({['x']: foo(), a: 0}).a")
	c.TestSame("({x: foo(),     a: 0}).a")

	// Leave unknown props alone, the might be on the prototype
	c.TestSame("({}).a")

	// setters by themselves don't provide a definition
	c.TestSame("({}).a")
	c.TestSame("({set a(b) {}}).a")
	// sets don't hide other definitions.
	c.Test("({a:1,set a(b) {}}).a", "1")

	// get is transformed to a call (gets don't have self referential names)
	c.Test("({get a() {}}).a", "(function (){})()")
	// sets don't hide other definitions.
	c.Test("({get a() {},set a(b) {}}).a", "(function (){})()")

	// a function remains a function not a call.
	c.Test("var x = ({a:function(){return 1}}).a", "var x = function(){return 1}")

	c.Test("var x = ({a:1}).a", "var x = 1")
	c.Test("var x = ({a:1, a:2}).a", "var x = 2")
	c.Test("var x = ({a:1, a:foo()}).a", "var x = foo()")
	c.Test("var x = ({a:foo()}).a", "var x = foo()")

	c.Test("function f() { return {a:1, b:2}.a; }", "function f() { return 1; }")

	// GETELEM is handled the same way.
	c.Test("var x = ({'a':1})['a']", "var x = 1")

	// try folding string computed properties
	c.Test("var a = {['a']:x}['a']", "var a = x")
	c.Test("var a = {['a']:x}?.['a']", "var a = x")
	c.Test("var a = {['a']:x}?.['a'].b", "var a = x.b")
	c.Test("var a = {a: {b: 1}}?.a?.b", "var a = 1")
	c.Test("var a = {a: null}?.a?.b", "var a = null?.b")
	c.Test("var a = {['a']: {b: 1}}?.['a']?.b", "var a = 1")
	c.Test("var a = {['a']: null}?.['a']?.b", "var a = null?.b")

	c.Test("var a = { get ['a']() { return 1; }}['a']", "var a = function() { return 1; }();")
	c.Test("var a = {'a': x, ['a']: y}['a']", "var a = y;")
	c.TestSame("var a = {['foo']: x}.a;")
	// Note: it may be useful to fold symbols in the future.
	c.TestSame("var y = Symbol(); var a = {[y]: 3}[y];")

	// We can fold member functions sometimes.
	//
	// <p>Even though they're different from fn expressions and arrow fns, extracting them only
	// causes programs that would have thrown errors to change behaviour.
	c.Test("var x = {a() { 1; }}.a;", "var x = function() { 1; };")
	// Notice `a` isn't invoked, so beahviour didn't change.
	c.Test("var x = {a() { return this; }}.a;", "var x = function() { return this; };")
	// `super` is invisibly captures the object that declared the method so we can't fold.
	c.TestSame("var x = {a() { return super.a; }}.a;")
	c.Test("var x = {a: 1, a() { 2; }}.a;", "var x = function() { 2; };")
	c.Test("var x = {a() {}, a: 1}.a;", "var x = 1;")
	c.TestSame("var x = {a() {}}.b")
	// Don't fold non-computed setters.
	c.TestSame("var x = ({ set a(v) { sideEffect() } }).a;")
	// Don't fold computed setters.
	c.TestSame("var x = ({ set ['a'](v) { sideEffect() } }).a;")
	// Fold computed getters.
	c.Test("var x = ({ get ['a']() { return 1; } }).a;", "var x = function() { return 1; }();")
	c.Test(
		"var x = ({ get ['a']() { return 1; }, set ['a'](v) {} }).a;",
		"var x = function() { return 1; }();")
	c.Test(
		"var x = ({ set ['a'](v) {}, get ['a']() { return 1; } }).a;",
		"var x = function() { return 1; }();")
}

// Java: PeepholeFoldConstantsTest.testFoldObjectLiteralRef2()
func (c *PeepholeFoldConstantsTest) TestFoldObjectLiteralRef2() {
	c.late = false
	c.Test("({a:x}).a += 1", "({a:x}).a = x + 1")
	c.late = true
	c.DisableNormalize()
	c.TestSame("({a:x}).a += 1")
}

// Regression test for https://github.com/google/closure-compiler/issues/2873
// It would be incorrect to fold this to "x();" because the 'this' value inside the function
// will be the global object, instead of the object {a:x} as it should be.
// Java: PeepholeFoldConstantsTest.testFoldObjectLiteral_methodCall_nonLiteralFn()
func (c *PeepholeFoldConstantsTest) TestFoldObjectLiteral_methodCall_nonLiteralFn() {
	c.TestSame("({a:x}).a()")
	c.TestSame("({a:x})?.a()")
	c.TestSame("({a:x})?.a()?.b")
}

// Java: PeepholeFoldConstantsTest.testFoldObjectLiteral_freeMethodCall()
func (c *PeepholeFoldConstantsTest) TestFoldObjectLiteral_freeMethodCall() {
	c.Test("({a() { return 1; }}).a()", "(function() { return 1; })()")
	c.Test("({a() { return 1; }})?.a()", "(function() { return 1; })()")
	// grandparent of optional chaining AST continues the chain
	c.Test("({a() { return 1; }})?.a().b", "(function() { return 1; })().b")
	c.Test("({a() { return 1; }})?.a().b.c?.d", "(function() { return 1; })().b.c?.d")
}

// Java: PeepholeFoldConstantsTest.testFoldObjectLiteral_freeArrowCall_usingEnclosingThis_late()
func (c *PeepholeFoldConstantsTest) TestFoldObjectLiteral_freeArrowCall_usingEnclosingThis_late() {
	c.late = true
	c.DisableNormalize()
	c.Test("({a: () => this }).a()", "(() => this)()")
	c.Test("({a: () => this })?.a()", "(() => this)()")
}

// Java: PeepholeFoldConstantsTest.testFoldObjectLiteral_unfreeMethodCall_dueToThis()
func (c *PeepholeFoldConstantsTest) TestFoldObjectLiteral_unfreeMethodCall_dueToThis() {
	c.TestSame("({a() { return this; }}).a()")
	c.TestSame("({a() { return this; }})?.a()")
}

// Java: PeepholeFoldConstantsTest.testFoldObjectLiteral_unfreeMethodCall_dueToSuper()
func (c *PeepholeFoldConstantsTest) TestFoldObjectLiteral_unfreeMethodCall_dueToSuper() {
	c.TestSame("({a() { return super.toString(); }}).a()")
	c.TestSame("({a() { return super.toString(); }})?.a()")
}

// Java: PeepholeFoldConstantsTest.testFoldObjectLiteral_paramToInvocation()
func (c *PeepholeFoldConstantsTest) TestFoldObjectLiteral_paramToInvocation() {
	c.Test("console.log({a: 1}.a)", "console.log(1)")
	c.Test("console.log({a: 1}?.a)", "console.log(1)")
}

// Java: PeepholeFoldConstantsTest.testIEString()
func (c *PeepholeFoldConstantsTest) TestIEString() {
	c.TestSame("!+'\\v1'")
}

// Java: PeepholeFoldConstantsTest.testIssue522()
func (c *PeepholeFoldConstantsTest) TestIssue522() {
	c.TestSame("[][1] = 1;")
}

// Java: PeepholeFoldConstantsTest.testTypeBasedFoldConstant()
func (c *PeepholeFoldConstantsTest) TestTypeBasedFoldConstant() {
	c.EnableTypeCheck()
	c.TestSame("function f(/** number */ x) { x + 1 + 1 + x; }")
	c.TestSame("function f(/** boolean */ x) { x + 1 + 1 + x; }")
	c.TestSame("function f(/** null */ x) { var y = true > x; }")
	c.TestSame("function f(/** null */ x) { var y = null > x; }")
	c.TestSame("function f(/** string */ x) { x + 1 + 1 + x; }")
	c.useTypes = false
	c.TestSame("function f(/** number */ x) { x + 1 + 1 + x; }")
}

// Java: PeepholeFoldConstantsTest.testColorBasedFoldConstant()
func (c *PeepholeFoldConstantsTest) TestColorBasedFoldConstant() {
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()
	c.TestSame("function f(/** number */ x) { x + 1 + 1 + x; }")
	c.TestSame("function f(/** boolean */ x) { x + 1 + 1 + x; }")
	c.TestSame("function f(/** null */ x) { var y = true > x; }")
	c.TestSame("function f(/** null */ x) { var y = null > x; }")
	c.TestSame("function f(/** string */ x) { x + 1 + 1 + x; }")
	c.useTypes = false
	c.TestSame("function f(/** number */ x) { x + 1 + 1 + x; }")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: PeepholeFoldConstantsTest.foldDefineProperties()
func (c *PeepholeFoldConstantsTest) TestFoldDefineProperties() {
	c.Test("Object.defineProperties({}, {})", "({})")
	c.Test("Object.defineProperties(a, {})", "a")
	c.TestSame("Object.defineProperties(a, {anything:1})")
}

// Java: PeepholeFoldConstantsTest.testES6Features()
func (c *PeepholeFoldConstantsTest) TestES6Features() {
	c.Test("var x = {[undefined != true] : 1};", "var x = {[true] : 1};")
	c.Test("let x = false && y;", "let x = false;")
	c.Test("const x = null == undefined", "const x = true")
	c.Test("var [a, , b] = [false+1, true+1, ![]]", "var a; var b; [a, , b] = [1, 2, false]")
	c.Test("var x = () =>  true || x;", "var x = () => { return true; }")
	c.Test(
		"function foo(x = (1 !== void 0), y) {return x+y;}",
		"function foo(x = true, y) {return x+y;}")
	c.Test(
		"class Foo {\n"+
			"  constructor() {this.x = null <= null;}\n"+
			"}\n",
		"class Foo {\n"+
			"  constructor() {this.x = true;}\n"+
			"}\n")
	c.Test("function foo() {return `${false && y}`}", "function foo() {return `false`}")
}

// Java: PeepholeFoldConstantsTest.testES6Features_late()
func (c *PeepholeFoldConstantsTest) TestES6Features_late() {
	c.late = true
	c.DisableNormalize()
	c.Test("var [a, , b] = [false+1, true+1, ![]]", "var [a, , b] = [1, 2, false]")
	c.Test("var x = () =>  true || x;", "var x = () => true;")
}

// Java: PeepholeFoldConstantsTest.testClassField()
func (c *PeepholeFoldConstantsTest) TestClassField() {
	c.Test(
		"class Foo {\n"+
			"  x = null <= null;\n"+
			"}\n",
		"class Foo {\n"+
			"  x = true;\n"+
			"}\n")
}

// Java: PeepholeFoldConstantsTest.testInvertibleOperators()
func (c *PeepholeFoldConstantsTest) TestInvertibleOperators() {
	inverses := jutil.NewLinkedHashMap[string, string]()
	inverses.Put("==", "!=")
	inverses.Put("===", "!==")
	inverses.Put("<=", ">")
	inverses.Put("<", ">=")
	inverses.Put(">=", "<")
	inverses.Put(">", "<=")
	inverses.Put("!=", "==")
	inverses.Put("!==", "===")
	comparators := jutil.NewLinkedHashSet("<=", "<", ">=", ">")
	equalitors := jutil.NewLinkedHashSet("==", "===")
	uncomparables := jutil.NewLinkedHashSet("undefined", "void 0")
	operators := slices.Collect(inverses.Values())
	for iOperandA := 0; iOperandA < len(peepholeFoldConstantsTest_LITERAL_OPERANDS); iOperandA++ {
		for iOperandB := 0; iOperandB < len(peepholeFoldConstantsTest_LITERAL_OPERANDS); iOperandB++ {
			for iOp := 0; iOp < len(operators); iOp++ {
				a := peepholeFoldConstantsTest_LITERAL_OPERANDS[iOperandA]
				b := peepholeFoldConstantsTest_LITERAL_OPERANDS[iOperandB]
				op := operators[iOp]
				inverse := inverses.Get(op)

				// Test invertability.
				if comparators.Contains(op) {
					if uncomparables.Contains(a) ||
						uncomparables.Contains(b) ||
						(a == "null" && NodeUtil_GetStringNumberValue(b) == nil) {
						c.AssertSameResults(PeepholeFoldConstantsTest_Join(a, op, b), "false")
						c.AssertSameResults(PeepholeFoldConstantsTest_Join(a, inverse, b), "false")
					}
				} else if a == b && equalitors.Contains(op) {
					if a == "NaN" || a == "Infinity" || a == "-Infinity" {
						expected := "true"
						if a == "NaN" {
							expected = "false"
						}
						c.Test(PeepholeFoldConstantsTest_Join(a, op, b), expected)
					} else {
						c.AssertSameResults(PeepholeFoldConstantsTest_Join(a, op, b), "true")
						c.AssertSameResults(PeepholeFoldConstantsTest_Join(a, inverse, b), "false")
					}
				} else {
					c.AssertNotSameResults(PeepholeFoldConstantsTest_Join(a, op, b), PeepholeFoldConstantsTest_Join(a, inverse, b))
				}
			}
		}
	}
}

// Java: PeepholeFoldConstantsTest.testCommutativeOperators()
func (c *PeepholeFoldConstantsTest) TestCommutativeOperators() {
	c.late = true
	c.DisableNormalize()
	operators :=
		[]string{"==", "!=", "===", "!==", "*", "|", "&", "^"}
	for _, a := range peepholeFoldConstantsTest_LITERAL_OPERANDS {
		for _, b := range peepholeFoldConstantsTest_LITERAL_OPERANDS {
			for _, op := range operators {
				// Test commutativity.
				c.AssertSameResults(PeepholeFoldConstantsTest_Join(a, op, b), PeepholeFoldConstantsTest_Join(b, op, a))
			}
		}
	}
}

// Java: PeepholeFoldConstantsTest.testConvertToNumberNegativeInf()
func (c *PeepholeFoldConstantsTest) TestConvertToNumberNegativeInf() {
	c.TestSame("var x = 3 & (r ? Infinity : -Infinity);")
}

// Java: PeepholeFoldConstantsTest.testAlgebraicIdentities()
func (c *PeepholeFoldConstantsTest) TestAlgebraicIdentities() {
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()
	c.FoldNumericTypes("x+0", "x")
	c.FoldNumericTypes("0+x", "x")
	c.FoldNumericTypes("x+0+0+x+x+0", "x+x+x")
	c.FoldNumericTypes("x-0", "x")
	c.FoldNumericTypes("x-0-0-0", "x")
	// 'x-0' is numeric even if x isn't
	c.Test("var x='a'; x-0-0", "var x='a';x-0")
	c.FoldNumericTypes("0-x", "-x")
	c.Test("for (var i = 0; i < 5; i++) var x = 0 + i * 1", "var i = 0; for(; i < 5; i++) var x=i")
	c.FoldNumericTypes("x*1", "x")
	c.FoldNumericTypes("1*x", "x")
	// can't optimize these without a non-NaN prover
	c.TestSame("x*0")
	c.TestSame("0*x")
	c.TestSame("0/x")
	c.FoldNumericTypes("x/1", "x")
}

// Java: PeepholeFoldConstantsTest.testBigIntAlgebraicIdentities()
func (c *PeepholeFoldConstantsTest) TestBigIntAlgebraicIdentities() {
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()
	c.FoldBigIntTypes("x+0n", "x")
	c.FoldBigIntTypes("0n+x", "x")
	c.FoldBigIntTypes("x+0n+0n+x+x+0n", "x+x+x")
	c.FoldBigIntTypes("x-0n", "x")
	c.FoldBigIntTypes("0n-x", "-x")
	c.FoldBigIntTypes("x-0n-0n-0n", "x")
	c.FoldBigIntTypes("x*1n", "x")
	c.FoldBigIntTypes("1n*x", "x")
	c.FoldBigIntTypes("x*1n*1n*x*x*1n", "x*x*x")
	c.FoldBigIntTypes("x/1n", "x")
	c.FoldBigIntTypes("x/0n", "x/0n")
	c.Test(
		"for (var i = 0n; i < 5n; i++) var x = 0n + i * 1n",
		"var i = 0n; for(; i < 5n; i++) var x=i")
	c.Test(
		"for (var i = 0n; i % 2n === 0n; i++) var x = i % 2n",
		"var i = 0n; for (; i % 2n === 0n; i++) var x = i % 2n")
	c.Test("(doSomething(),0n)*1n", "(doSomething(),0n)")
	c.Test("1n*(doSomething(),0n)", "(doSomething(),0n)")
	c.IgnoreWarnings(DiagnosticGroups_CHECK_TYPES)
	c.TestSame("(0n,doSomething())*1n")
	c.TestSame("1n*(0n,doSomething())")
}

// Java: PeepholeFoldConstantsTest.testAssociativeFoldConstantsWithVariables()
func (c *PeepholeFoldConstantsTest) TestAssociativeFoldConstantsWithVariables() {
	// MUL and ADD should not fold
	c.TestSame("alert(x * 12 * 20);")
	c.TestSame("alert(12 * x * 20);")
	c.TestSame("alert(x + 12 + 20);")
	c.TestSame("alert(12 + x +  20);")
	c.Test("alert(x & 12 & 20);", "alert(x & 4);")
	c.Test("alert(12 & x & 20);", "alert(x & 4);")
}

// Java: PeepholeFoldConstantsTest.testTemplateLiteralConcat()
func (c *PeepholeFoldConstantsTest) TestTemplateLiteralConcat() {
	// join at variables
	c.Test(
		"const a = `${x}` + `${x}`;\n",
		"const a = `${x}${x}`;\n")
	c.Test(
		"const a = `${x}${y}` + `${x}${y}`;\n",
		"const a = `${x}${y}${x}${y}`;\n")

	// join at variable and string
	c.Test(
		"const a = `${x}` + `a${x}`;\n",
		"const a =`${x}a${x}`;\n")
	c.Test(
		"const a = `${x}` + ` ${x}`;\n",
		"const a =`${x} ${x}`;\n")

	// join at string and variable
	c.Test(
		"const a = `a${x}b` + `${x}c`;\n",
		"const a = `a${x}b${x}c`;\n")
	c.Test(
		"const a = `a${x} ` + `${x}b`;\n",
		"const a = `a${x} ${x}b`;\n")

	// join at strings
	c.Test(
		"const x = 'X';\n"+
			"console.log(`a${x}b` + `c${x}d`);\n",
		"const x = 'X';\n"+
			"console.log(`a${x}bc${x}d`);\n")

	// complex joins
	c.Test(
		"const a = `${foo() + bar()}` + `${baz()}`;\n",
		"const a =`${foo() + bar()}${baz()}`;\n")
	c.Test(
		"console.log(`<h1>${url}</h1>` + `<p>The URL is ${url}.</p>`);\n",
		"console.log(`<h1>${url}</h1><p>The URL is ${url}.</p>`);\n")
	c.Test(
		"console.log(`${(() => {return url;})()}</h1>` + `<p>The URL is ${url}.</p>`);\n",
		"console.log(`${(() => {return url;})()}</h1><p>The URL is ${url}.</p>`);\n")

	// don't fold tagged template literals
	c.TestSame(
		"const a = foo`${b}` + `${b}`;\n")
	c.TestSame(
		"const a = foo`${b}` + bar`${b}`;\n")
	// don't fold when raw strings meet at unescaped $ and {
	c.TestSame("`${x}$` + `{y}${z}`")
	c.TestSame("`${x}$` + `{y}`")
	c.Test("`${x}\\$` + `{y}`", "`${x}\\${y}`")
	c.TestSame("`${x}\\\\$` + `{y}`")
}

// Java: PeepholeFoldConstantsTest.testFoldTemplateLiteralSubstitutions()
func (c *PeepholeFoldConstantsTest) TestFoldTemplateLiteralSubstitutions() {
	// Basic substitution
	c.Test("`a${'b'}c`", "`abc`")
	// Number substitution
	c.Test("`a${123}c`", "`a123c`")
	// Boolean substitution
	c.Test("`a${true}c`", "`atruec`")
	c.Test("`a${false}c`", "`afalsec`")

	// Multiple substitutions
	c.Test("`a${'b'}c${'d'}e`", "`abcde`")
	// Adjacent substitutions
	c.Test("`${'a'}${'b'}`", "`ab`")
	// Empty strings
	c.Test("`${''}a${''}`", "`a`")

	// Corner cases where it should NOT fold
	c.TestSame("`a${'`'}c`")
	c.TestSame("`a${'$'}c`")
	c.TestSame("`a${'\\\\'}c`")
	c.TestSame("`a${'\\\\n'}c`")
	c.TestSame("`a${'\\\\r'}c`")
	c.TestSame("`${x} ${'$'}${'{'}`") // specifically tests that $ and { are not folded together
	c.Test("`$${''}{${''}`", "`$${''}{`")
	c.TestSame("`$${''}{foo}`")
	c.Test("`\\$${''}{foo}`", "`\\${foo}`")
	// These should not fold because \0 followed by a digit is a syntax error in untagged templates.
	c.TestSame("`\\0${'77'}`")
	c.TestSame("`\\0${77}`")
	c.Test("`\\0${'a'}`", "`\\0a` ")
	c.Test("`\\\\0${'7'}`", "`\\\\07` ")
	c.Test("`\\0${''}${'7'}`", "`\\0${'7'}`")

	// Tagged template literals should NOT fold
	c.TestSame("tag`a${'b'}c`")

	// Non-constant substitutions should NOT fold
	c.TestSame("`a${x}c`")
	c.TestSame("`a${foo()}c`")
}

// Java: PeepholeFoldConstantsTest.testFoldAddTemplateLiteralAndConstant()
func (c *PeepholeFoldConstantsTest) TestFoldAddTemplateLiteralAndConstant() {
	// String merge
	c.Test("`a${x}` + 'b'", "`a${x}b` ")
	c.Test("'a' + `b${x}`", "`ab${x}` ")
	// Number merge
	c.Test("`a${x}` + 1", "`a${x}1` ")
	c.Test("1 + `a${x}`", "`1a${x}` ")
	// Boolean merge
	c.Test("`a${x}` + true", "`a${x}true` ")
	c.Test("true + `a${x}`", "`truea${x}` ")
	// BigInt merge
	c.Test("`a${x}` + 10n", "`a${x}10` ")
	// Null/Undefined merge
	c.Test("`a${x}` + null", "`a${x}null` ")
	c.Test("`a${x}` + undefined", "`a${x}undefined` ")
	// Nested templates
	c.Test("`a${x}` + `b` ", "`a${x}b` ")
	c.Test("`a` + `b${x}` ", "`ab${x}` ")
	// Restricted octal escape
	c.TestSame("`\\0` + 7")
	c.TestSame("`\\0` + '7'")
	c.Test("`\\0a` + '7'", "`\\0a7` ")
	// Safety - special chars should NOT fold
	c.TestSame("`a${x}` + '$'")
	c.TestSame("`a${x}` + '{'")
	c.TestSame("`a${x}` + '`'")
	c.TestSame("`a${x}` + '\\\\'")
	c.TestSame("`a${x}` + '\\\\n'")
	// Compound addition
	c.Test("`a` + `b${x}` + `c` + 'd'", "`ab${x}cd` ")
}

// Java: PeepholeFoldConstantsTest.testFoldAddTemplateLiterals_validateRawAndCookedStringForSpecialChars_andValidateChildren()
func (c *PeepholeFoldConstantsTest) TestFoldAddTemplateLiterals_validateRawAndCookedStringForSpecialChars_andValidateChildren() {
	// This is a bit tricky because there are two layers of escaping.
	// Java sees "\\n" as "(\\)n" so it becomes "\n" for the javascript compiler.
	// The compiler then interprets this as the special character "\n"
	c.Test("var x = `${a}\\n` + `\\t${b}`", "var x = `${a}\\n\\t${b}`")
	// getJsRoot() returns the ROOT node for inputs only.
	// The first child is the SCRIPT node containing the input code.
	script := c.GetLastCompiler().GetJsRoot().GetFirstChild()
	var_ := script.GetFirstChild()
	name := var_.GetFirstChild()
	templateLit := name.GetFirstChild()
	AssertThat(templateLit.IsTemplateLit()).IsTrue()
	// TEMPLATELIT should have 5 total children
	AssertThat(templateLit.GetChildCount()).IsEqualTo(5)
	// Index 0: TEMPLATELIT_STRING ""
	child0 := templateLit.GetChildAtIndex(0)
	AssertThat(child0.IsTemplateLitString()).IsTrue()
	AssertThat(child0.GetCookedString()).IsEmpty()
	AssertThat(child0.GetRawString()).IsEmpty()
	// Index 1: SUB a
	child1 := templateLit.GetChildAtIndex(1)
	AssertThat(child1.IsTemplateLitSub()).IsTrue()
	AssertThat(child1.GetOnlyChild().IsName()).IsTrue()
	AssertThat(child1.GetOnlyChild().GetString()).IsEqualTo("a")
	// Index 2: TEMPLATELIT_STRING "\n\t"
	child2 := templateLit.GetChildAtIndex(2)
	AssertThat(child2.IsTemplateLitString()).IsTrue()
	AssertThat(child2.GetCookedString()).IsEqualTo("\n\t")
	AssertThat(child2.GetRawString()).IsEqualTo("\\n\\t")
	// Index 3: SUB b
	child3 := templateLit.GetChildAtIndex(3)
	AssertThat(child3.IsTemplateLitSub()).IsTrue()
	AssertThat(child3.GetOnlyChild().IsName()).IsTrue()
	AssertThat(child3.GetOnlyChild().GetString()).IsEqualTo("b")
	// Index 4: TEMPLATELIT_STRING ""
	child4 := templateLit.GetChildAtIndex(4)
	AssertThat(child4.IsTemplateLitString()).IsTrue()
	AssertThat(child4.GetCookedString()).IsEmpty()
	AssertThat(child4.GetRawString()).IsEmpty()
}

// Java: PeepholeFoldConstantsTest.testFoldAddTemplateLiterals_validateRawAndCookedStringForLiteralBackslash_andValidateChildren()
func (c *PeepholeFoldConstantsTest) TestFoldAddTemplateLiterals_validateRawAndCookedStringForLiteralBackslash_andValidateChildren() {
	// This is a bit tricky because there are two layers of escaping.
	// Java sees "\\\\n" as "(\\)(\\)n" so it becomes "\\n" for the javascript compiler.
	// The compiler then interprets this as "(\\)n" so it becomes "\" followed by
	// "n" NOT the special character "\n".
	c.Test(
		"var x = `${foo()}\\\\` + `n\\t${()=> {return b;}}`",
		"var x = `${foo()}\\\\n\\t${()=> {return b;}}`")
	// getJsRoot() returns the ROOT node for inputs only.
	// The first child is the SCRIPT node containing the input code.
	script := c.GetLastCompiler().GetJsRoot().GetFirstChild()
	var_ := script.GetFirstChild()
	name := var_.GetFirstChild()
	templateLit := name.GetFirstChild()
	AssertThat(templateLit.IsTemplateLit()).IsTrue()
	// TEMPLATELIT should have 5 children:
	AssertThat(templateLit.GetChildCount()).IsEqualTo(5)
	// Index 0: TEMPLATELIT_STRING ""
	child0 := templateLit.GetChildAtIndex(0)
	AssertThat(child0.IsTemplateLitString()).IsTrue()
	AssertThat(child0.GetCookedString()).IsEmpty()
	AssertThat(child0.GetRawString()).IsEmpty()
	// Index 1: SUB foo()
	child1 := templateLit.GetChildAtIndex(1)
	AssertThat(child1.IsTemplateLitSub()).IsTrue()
	AssertThat(child1.GetOnlyChild().IsCall()).IsTrue()
	AssertThat(child1.GetOnlyChild().GetOnlyChild().GetString()).IsEqualTo("foo")
	// Index 2: TEMPLATELIT_STRING "\n\t"
	child2 := templateLit.GetChildAtIndex(2)
	AssertThat(child2.IsTemplateLitString()).IsTrue()
	AssertThat(child2.GetCookedString()).IsEqualTo("\\n\t")
	AssertThat(child2.GetRawString()).IsEqualTo("\\\\n\\t")
	// Index 3: SUB ()=> {return b;}
	child3 := templateLit.GetChildAtIndex(3)
	AssertThat(child3.IsTemplateLitSub()).IsTrue()
	AssertThat(child3.GetOnlyChild().IsFunction()).IsTrue()
	AssertThat(child3.GetOnlyChild().GetFirstChild().GetString()).IsEmpty()
	// Index 4: TEMPLATELIT_STRING ""
	child4 := templateLit.GetChildAtIndex(4)
	AssertThat(child4.IsTemplateLitString()).IsTrue()
	AssertThat(child4.GetCookedString()).IsEmpty()
	AssertThat(child4.GetRawString()).IsEmpty()
}

// Java: PeepholeFoldConstantsTest.foldBigIntTypes(String, String)
func (c *PeepholeFoldConstantsTest) FoldBigIntTypes(js string, expected string) {
	c.Test(
		"function f(/** @type {bigint} */ x) { "+js+" }",
		"function f(/** @type {bigint} */ x) { "+expected+" }")
}

// Java: PeepholeFoldConstantsTest.foldNumericTypes(String, String)
func (c *PeepholeFoldConstantsTest) FoldNumericTypes(js string, expected string) {
	c.Test(
		"function f(/** @type {number} */ x) { "+js+" }",
		"function f(/** @type {number} */ x) { "+expected+" }")
}

// Java: PeepholeFoldConstantsTest.foldStringTypes(String, String)
func (c *PeepholeFoldConstantsTest) FoldStringTypes(js string, expected string) {
	c.Test(
		"function f(/** @type {string} */ x) { "+js+" }",
		"function f(/** @type {string} */ x) { "+expected+" }")
}

// Java: PeepholeFoldConstantsTest.join(String, String, String)
func PeepholeFoldConstantsTest_Join(operandA string, op string, operandB string) string {
	return operandA + " " + op + " " + operandB
}

// Java: PeepholeFoldConstantsTest.assertSameResults(String, String)
func (c *PeepholeFoldConstantsTest) AssertSameResults(exprA string, exprB string) {
	AssertWithMessage("Expressions did not fold the same\nexprA: %s\nexprB: %s", exprA, exprB).
		That(c.Process(exprB)).
		IsEqualTo(c.Process(exprA))
}

// Java: PeepholeFoldConstantsTest.assertNotSameResults(String, String)
func (c *PeepholeFoldConstantsTest) AssertNotSameResults(exprA string, exprB string) {
	AssertWithMessage("Expressions folded the same\nexprA: %s\nexprB: %s", exprA, exprB).
		That(c.Process(exprA) == c.Process(exprB)).
		IsFalse()
}

// Java: PeepholeFoldConstantsTest.process(String)
func (c *PeepholeFoldConstantsTest) Process(js string) string {
	return c.PrintHelper(js, true)
}

// Java: PeepholeFoldConstantsTest.printHelper(String, boolean)
func (c *PeepholeFoldConstantsTest) PrintHelper(js string, runProcessor bool) string {
	compiler := c.CreateCompiler()
	options := c.GetOptions()
	compiler.Init(
		[]*SourceFile{},
		[]*SourceFile{SourceFile_FromCode("testcode", js)},
		options)
	root := compiler.ParseInputs()
	compiler.SetAccessorSummary(AccessorSummary_Create(jutil.NewLinkedHashMap[string, AccessorSummary_PropertyAccessKind]()))
	errors := compiler.GetErrors()
	errorStrings := make([]string, 0, len(errors))
	for _, e := range errors {
		errorStrings = append(errorStrings, jutil.ToJavaString(e))
	}
	AssertWithMessage(
		"Unexpected parse error(s): %s\nEXPR: %s",
		strings.Join(errorStrings, "\n"), js).
		That(root).
		IsNotNull()
	externsRoot := root.GetFirstChild()
	mainRoot := externsRoot.GetNext()
	if runProcessor {
		c.GetProcessor(compiler).Process(externsRoot, mainRoot)
	}
	return compiler.ToSourceOfNode(mainRoot)
}

// Java: PeepholeFoldConstantsTest.testFoldOptChain_nonNullReceiver_chainContinues()
func (c *PeepholeFoldConstantsTest) TestFoldOptChain_nonNullReceiver_chainContinues() {
	// TODO(b/571220152): folding an opt-chain segment start orphans the continuation
	e := AssertThrows[jutil.IRuntimeException](func() { c.TestSame("x = [1, 2, 3]?.length.toString()") })
	// Throwables.getRootCause(e)
	var root jutil.Throwable = e
	for root.GetCause() != nil && root.GetCause() != root {
		root = root.GetCause()
	}
	_, isIllegalState := root.(jutil.IIllegalStateException)
	AssertThat(isIllegalState).IsTrue()
	AssertThat(root.GetMessage()).
		Contains("Start of optional chain node OPTCHAIN_GETPROP is not marked as the start")
}

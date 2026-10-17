// Port of com/google/javascript/jscomp/PeepholeReplaceKnownMethodsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/PeepholeReplaceKnownMethodsTest.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PeepholeReplaceKnownMethodsTest is the unit test for
// PeepholeReplaceKnownMethods.
type PeepholeReplaceKnownMethodsTest struct {
	CompilerTestCaseBase

	late     bool
	useTypes bool
}

// Java: PeepholeReplaceKnownMethodsTest()
func NewPeepholeReplaceKnownMethodsTest() *PeepholeReplaceKnownMethodsTest {
	c := &PeepholeReplaceKnownMethodsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c,
		CompilerTestCase_MINIMAL_EXTERNS+
			"/** @type {function(this: Array, ...*): !Array<?>} */ function returnArrayType() {}\n"+
			"/** @type {function(this: Array, ...*): !Array<?>|string} */ function returnUnionType(){}\n"+
			"/** @constructor */ function Foo(){}\n"+
			"/** @type {function(this: Foo, ...*): !Foo} */ Foo.prototype.concat\n"+
			"var obj = new Foo();\n"+
			"/**\n"+
			" * @param {...T} var_args\n"+
			" * @return {!Array<T>}\n"+
			" * @template T\n"+
			" */\n"+
			"Array.of = function(var_args) {};\n")
	c.late = true
	c.useTypes = true
	return c
}

func TestPeepholeReplaceKnownMethodsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPeepholeReplaceKnownMethodsTest() })
}

// Java: PeepholeReplaceKnownMethodsTest.setUp()
func (c *PeepholeReplaceKnownMethodsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.late = true
	c.useTypes = true
	c.DisableTypeCheck()
	c.EnableNormalize()
}

// Java: PeepholeReplaceKnownMethodsTest.getProcessor(Compiler)
func (c *PeepholeReplaceKnownMethodsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPeepholeOptimizationsPass(
		compiler, c.GetName(), NewPeepholeReplaceKnownMethods(c.late, c.useTypes))
}

// Java: PeepholeReplaceKnownMethodsTest.testStringIndexOf()
func (c *PeepholeReplaceKnownMethodsTest) TestStringIndexOf() {
	c.Fold("x = 'abcdef'.indexOf('g')", "x = -1")
	c.Fold("x = 'abcdef'.indexOf('b')", "x = 1")
	c.Fold("x = 'abcdefbe'.indexOf('b', 2)", "x = 6")
	c.Fold("x = 'abcdef'.indexOf('bcd')", "x = 1")
	c.Fold("x = 'abcdefsdfasdfbcdassd'.indexOf('bcd', 4)", "x = 13")
	c.Fold("x = 'abcdef'.lastIndexOf('b')", "x = 1")
	c.Fold("x = 'abcdefbe'.lastIndexOf('b')", "x = 6")
	c.Fold("x = 'abcdefbe'.lastIndexOf('b', 5)", "x = 1")
	// Both elements must be strings. Don't do anything if either one is not
	// string.
	c.Fold("x = 'abc1def'.indexOf(1)", "x = 3")
	c.Fold("x = 'abcNaNdef'.indexOf(NaN)", "x = 3")
	c.Fold("x = 'abcundefineddef'.indexOf(undefined)", "x = 3")
	c.Fold("x = 'abcnulldef'.indexOf(null)", "x = 3")
	c.Fold("x = 'abctruedef'.indexOf(true)", "x = 3")
	// The following test case fails with JSC_PARSE_ERROR. Hence omitted.
	// foldSame("x = 1.indexOf('bcd');");
	c.FoldSame("x = NaN.indexOf('bcd')")
	c.FoldSame("x = undefined.indexOf('bcd')")
	c.FoldSame("x = null.indexOf('bcd')")
	c.FoldSame("x = true.indexOf('bcd')")
	c.FoldSame("x = false.indexOf('bcd')")
	// Avoid dealing with regex or other types.
	c.FoldSame("x = 'abcdef'.indexOf(/b./)")
	c.FoldSame("x = 'abcdef'.indexOf({a:2})")
	c.FoldSame("x = 'abcdef'.indexOf([1,2])")
	// Template Strings
	c.FoldSame("x = `abcdef`.indexOf('b')")
	c.FoldSame("x = `Hello ${name}`.indexOf('a')")
	c.FoldSame("x = tag `Hello ${name}`.indexOf('a')")
}

// Java: PeepholeReplaceKnownMethodsTest.testStringJoinAddSparse()
func (c *PeepholeReplaceKnownMethodsTest) TestStringJoinAddSparse() {
	c.Fold("x = [,,'a'].join(',')", "x = ',,a'")
}

// Java: PeepholeReplaceKnownMethodsTest.testNoStringJoin()
func (c *PeepholeReplaceKnownMethodsTest) TestNoStringJoin() {
	c.FoldSame("x = [].join(',',2)")
	c.FoldSame("x = [].join(f)")
}

// Java: PeepholeReplaceKnownMethodsTest.testStringJoinAdd()
func (c *PeepholeReplaceKnownMethodsTest) TestStringJoinAdd() {
	c.Fold("x = ['a', 'b', 'c'].join('')", "x = \"abc\"")
	c.Fold("x = [].join(',')", "x = \"\"")
	c.Fold("x = ['a'].join(',')", "x = \"a\"")
	c.Fold("x = ['a', 'b', 'c'].join(',')", "x = \"a,b,c\"")
	c.Fold("x = ['a', foo, 'b', 'c'].join(',')", "x = [\"a\",foo,\"b,c\"].join()")
	c.Fold("x = [foo, 'a', 'b', 'c'].join(',')", "x = [foo,\"a,b,c\"].join()")
	c.Fold("x = ['a', 'b', 'c', foo].join(',')", "x = [\"a,b,c\",foo].join()")
	// Works with numbers
	c.Fold("x = ['a=', 5].join('')", "x = \"a=5\"")
	c.Fold("x = ['a', '5'].join(7)", "x = \"a75\"")
	// Works on boolean
	c.Fold("x = ['a=', false].join('')", "x = \"a=false\"")
	c.Fold("x = ['a', '5'].join(true)", "x = \"atrue5\"")
	c.Fold("x = ['a', '5'].join(false)", "x = \"afalse5\"")
	// Only optimize if it's a size win.
	c.Fold(
		"x = ['a', '5', 'c'].join('a very very very long chain')",
		"x = [\"a\",\"5\",\"c\"].join(\"a very very very long chain\")")
	// Template strings
	c.Fold("x = [`a`, `b`, `c`].join(``)", "x = 'abc'")
	c.Fold("x = [`a`, `b`, `c`].join('')", "x = 'abc'")
	// TODO(user): Its possible to fold this better.
	c.FoldSame("x = ['', foo].join('-')")
	c.FoldSame("x = ['', foo, ''].join()")
	c.Fold(
		"x = ['', '', foo, ''].join(',')", //
		"x = [ ','  , foo, ''].join()")
	c.Fold(
		"x = ['', '', foo, '', ''].join(',')", //
		"x = [ ',',   foo,  ','].join()")
	c.Fold(
		"x = ['', '', foo, '', '', bar].join(',')", //
		"x = [ ',',   foo,  ',',   bar].join()")
	c.Fold(
		"x = [1,2,3].join('abcdef')", //
		"x = '1abcdef2abcdef3'")
	c.Fold("x = [1,2].join()", "x = '1,2'")
	c.Fold("x = [null,undefined,''].join(',')", "x = ',,'")
	c.Fold("x = [null,undefined,0].join(',')", "x = ',,0'")
	// This can be folded but we don't currently.
	c.FoldSame("x = [[1,2],[3,4]].join()")
	// would like: "x = '1,2,3,4'"
}

// Java: PeepholeReplaceKnownMethodsTest.testStringJoinAdd_b1992789()
func (c *PeepholeReplaceKnownMethodsTest) TestStringJoinAdd_b1992789() {
	c.Fold("x = ['a'].join('')", "x = \"a\"")
	c.FoldSame("x = [foo()].join('')")
	c.FoldSame("[foo()].join('')")
	c.Fold("[null].join('')", "''")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldStringSubstr()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldStringSubstr() {
	c.Fold("x = 'abcde'.substr(0,2)", "x = 'ab'")
	c.Fold("x = 'abcde'.substr(1,2)", "x = 'bc'")
	c.Fold("x = 'abcde'.substr(2)", "x = 'cde'")
	// we should be leaving negative indexes alone for now
	c.FoldSame("x = 'abcde'.substr(-1)")
	c.FoldSame("x = 'abcde'.substr(1, -2)")
	c.FoldSame("x = 'abcde'.substr(1, 2, 3)")
	c.FoldSame("x = 'a'.substr(0, 2)")
	// Template strings
	c.FoldSame("x = `abcdef`.substr(0,2)")
	c.FoldSame("x = `abc ${xyz} def`.substr(0,2)")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldStringReplace()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldStringReplace() {
	c.Fold("'c'.replace('c','x')", "'x'")
	c.Fold("'ac'.replace('c','x')", "'ax'")
	c.Fold("'ca'.replace('c','x')", "'xa'")
	c.Fold("'ac'.replace('c','xxx')", "'axxx'")
	c.Fold("'ca'.replace('c','xxx')", "'xxxa'")
	// only one instance replaced
	c.Fold("'acaca'.replace('c','x')", "'axaca'")
	c.Fold("'ab'.replace('','x')", "'xab'")
	c.FoldSame("'acaca'.replace(/c/,'x')")
	// this will affect the global RegExp props
	c.FoldSame("'acaca'.replace(/c/g,'x')")
	// this will affect the global RegExp props
	// not a literal
	c.FoldSame("x.replace('x','c')")
	c.FoldSame("'Xyz'.replace('Xyz', '$$')")
	// would fold to '$'
	c.FoldSame("'PreXyzPost'.replace('Xyz', '$&')")
	// would fold to 'PreXyzPost'
	c.FoldSame("'PreXyzPost'.replace('Xyz', '$`')")
	// would fold to 'PrePrePost'
	c.FoldSame("'PreXyzPost'.replace('Xyz', '$\\'')")
	// would fold to  'PrePostPost'
	c.FoldSame("'PreXyzPostXyz'.replace('Xyz', '$\\'')")
	// would fold to 'PrePostXyzPostXyz'
	c.FoldSame("'123'.replace('2', '$`')")
	// would fold to '113'
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldStringReplaceAll()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldStringReplaceAll() {
	c.Fold("x = 'abcde'.replaceAll('bcd','c')", "x = 'ace'")
	c.Fold("x = 'abcde'.replaceAll('c','xxx')", "x = 'abxxxde'")
	c.Fold("x = 'abcde'.replaceAll('xxx','c')", "x = 'abcde'")
	c.Fold("'ab'.replaceAll('','x')", "'xaxbx'")
	c.Fold("x = 'c_c_c'.replaceAll('c','x')", "x = 'x_x_x'")
	c.FoldSame("x = 'acaca'.replaceAll(/c/,'x')")
	// this should throw
	c.FoldSame("x = 'acaca'.replaceAll(/c/g,'x')")
	// this will affect the global RegExp props
	// not a literal
	c.FoldSame("x.replaceAll('x','c')")
	c.FoldSame("'Xyz'.replaceAll('Xyz', '$$')")
	// would fold to '$'
	c.FoldSame("'PreXyzPost'.replaceAll('Xyz', '$&')")
	// would fold to 'PreXyzPost'
	c.FoldSame("'PreXyzPost'.replaceAll('Xyz', '$`')")
	// would fold to 'PrePrePost'
	c.FoldSame("'PreXyzPost'.replaceAll('Xyz', '$\\'')")
	// would fold to  'PrePostPost'
	c.FoldSame("'PreXyzPostXyz'.replaceAll('Xyz', '$\\'')")
	// would fold to 'PrePostXyzPost'
	c.FoldSame("'123'.replaceAll('2', '$`')")
	// would fold to '113'
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldStringSubstring()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldStringSubstring() {
	c.Fold("x = 'abcde'.substring(0,2)", "x = 'ab'")
	c.Fold("x = 'abcde'.substring(1,2)", "x = 'b'")
	c.Fold("x = 'abcde'.substring(2)", "x = 'cde'")
	// we should be leaving negative, out-of-bound, and inverted indices alone for now
	c.FoldSame("x = 'abcde'.substring(-1)")
	c.FoldSame("x = 'abcde'.substring(1, -2)")
	c.FoldSame("x = 'abcde'.substring(1, 2, 3)")
	c.FoldSame("x = 'abcde'.substring(2, 0)")
	c.FoldSame("x = 'a'.substring(0, 2)")
	// Template strings
	c.FoldSame("x = `abcdef`.substring(0,2)")
	c.FoldSame("x = `abcdef ${abc}`.substring(0,2)")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldStringSlice()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldStringSlice() {
	c.Fold("x = 'abcde'.slice(0,2)", "x = 'ab'")
	c.Fold("x = 'abcde'.slice(1,2)", "x = 'b'")
	c.Fold("x = 'abcde'.slice(2)", "x = 'cde'")
	// we should be leaving negative, out-of-bound, and inverted indices alone for now
	c.FoldSame("x = 'abcde'.slice(-1)")
	c.FoldSame("x = 'abcde'.slice(1, -2)")
	c.FoldSame("x = 'abcde'.slice(1, 2, 3)")
	c.FoldSame("x = 'abcde'.slice(2, 0)")
	c.FoldSame("x = 'a'.slice(0, 2)")
	// Template strings
	c.FoldSame("x = `abcdef`.slice(0,2)")
	c.FoldSame("x = `abcdef ${abc}`.slice(0,2)")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldStringCharAt()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldStringCharAt() {
	c.Fold("x = 'abcde'.charAt(0)", "x = 'a'")
	c.Fold("x = 'abcde'.charAt(1)", "x = 'b'")
	c.Fold("x = 'abcde'.charAt(2)", "x = 'c'")
	c.Fold("x = 'abcde'.charAt(3)", "x = 'd'")
	c.Fold("x = 'abcde'.charAt(4)", "x = 'e'")
	c.FoldSame("x = 'abcde'.charAt(5)")
	// or x = ''
	c.FoldSame("x = 'abcde'.charAt(-1)")
	// or x = ''
	c.FoldSame("x = 'abcde'.charAt(y)")
	c.FoldSame("x = 'abcde'.charAt()")
	// or x = 'a'
	c.FoldSame("x = 'abcde'.charAt(0, ++z)")
	// or (++z, 'a')
	c.FoldSame("x = 'abcde'.charAt(null)")
	// or x = 'a'
	c.FoldSame("x = 'abcde'.charAt(true)")
	// or x = 'b'
	c.Fold(jutil.StringFromUTF16([]uint16{0x0078, 0x0020, 0x003d, 0x0020, 0x0027, 0x005c, 0x0075, 0x0064, 0x0038, 0x0033, 0x0034, 0xdd1e, 0x0027, 0x002e, 0x0063, 0x0068, 0x0061, 0x0072, 0x0041, 0x0074, 0x0028, 0x0030, 0x0029}), "x = '\\ud834'")
	c.Fold(jutil.StringFromUTF16([]uint16{0x0078, 0x0020, 0x003d, 0x0020, 0x0027, 0x005c, 0x0075, 0x0064, 0x0038, 0x0033, 0x0034, 0xdd1e, 0x0027, 0x002e, 0x0063, 0x0068, 0x0061, 0x0072, 0x0041, 0x0074, 0x0028, 0x0031, 0x0029}), "x = '\\udd1e'")
	// Template strings
	c.FoldSame("x = `abcdef`.charAt(0)")
	c.FoldSame("x = `abcdef ${abc}`.charAt(0)")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldStringCharCodeAt()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldStringCharCodeAt() {
	c.Fold("x = 'abcde'.charCodeAt(0)", "x = 97")
	c.Fold("x = 'abcde'.charCodeAt(1)", "x = 98")
	c.Fold("x = 'abcde'.charCodeAt(2)", "x = 99")
	c.Fold("x = 'abcde'.charCodeAt(3)", "x = 100")
	c.Fold("x = 'abcde'.charCodeAt(4)", "x = 101")
	c.FoldSame("x = 'abcde'.charCodeAt(5)")
	// or x = (0/0)
	c.FoldSame("x = 'abcde'.charCodeAt(-1)")
	// or x = (0/0)
	c.FoldSame("x = 'abcde'.charCodeAt(y)")
	c.FoldSame("x = 'abcde'.charCodeAt()")
	// or x = 97
	c.FoldSame("x = 'abcde'.charCodeAt(0, ++z)")
	// or (++z, 97)
	c.FoldSame("x = 'abcde'.charCodeAt(null)")
	// or x = 97
	c.FoldSame("x = 'abcde'.charCodeAt(true)")
	// or x = 98
	c.Fold(jutil.StringFromUTF16([]uint16{0x0078, 0x0020, 0x003d, 0x0020, 0x0027, 0x005c, 0x0075, 0x0064, 0x0038, 0x0033, 0x0034, 0xdd1e, 0x0027, 0x002e, 0x0063, 0x0068, 0x0061, 0x0072, 0x0043, 0x006f, 0x0064, 0x0065, 0x0041, 0x0074, 0x0028, 0x0030, 0x0029}), "x = 55348")
	c.Fold(jutil.StringFromUTF16([]uint16{0x0078, 0x0020, 0x003d, 0x0020, 0x0027, 0x005c, 0x0075, 0x0064, 0x0038, 0x0033, 0x0034, 0xdd1e, 0x0027, 0x002e, 0x0063, 0x0068, 0x0061, 0x0072, 0x0043, 0x006f, 0x0064, 0x0065, 0x0041, 0x0074, 0x0028, 0x0031, 0x0029}), "x = 56606")
	// Template strings
	c.FoldSame("x = `abcdef`.charCodeAt(0)")
	c.FoldSame("x = `abcdef ${abc}`.charCodeAt(0)")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldStringSplit()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldStringSplit() {
	c.late = false
	c.Fold("x = 'abcde'.split('foo')", "x = ['abcde']")
	c.Fold("x = 'abcde'.split()", "x = ['abcde']")
	c.Fold("x = 'abcde'.split(null)", "x = ['abcde']")
	c.Fold("x = 'a b c d e'.split(' ')", "x = ['a','b','c','d','e']")
	c.Fold("x = 'a b c d e'.split(' ', 0)", "x = []")
	c.Fold("x = 'abcde'.split('cd')", "x = ['ab','e']")
	c.Fold("x = 'a b c d e'.split(' ', 1)", "x = ['a']")
	c.Fold("x = 'a b c d e'.split(' ', 3)", "x = ['a','b','c']")
	c.Fold("x = 'a b c d e'.split(null, 1)", "x = ['a b c d e']")
	c.Fold("x = 'aaaaa'.split('a')", "x = ['', '', '', '', '', '']")
	c.Fold("x = 'xyx'.split('x')", "x = ['', 'y', '']")
	// Empty separator
	c.Fold("x = 'abcde'.split('')", "x = ['a','b','c','d','e']")
	c.Fold("x = 'abcde'.split('', 3)", "x = ['a','b','c']")
	// Empty separator AND empty string
	c.Fold("x = ''.split('')", "x = []")
	// Separator equals string
	c.Fold("x = 'aaa'.split('aaa')", "x = ['','']")
	c.Fold("x = ' '.split(' ')", "x = ['','']")
	c.FoldSame("x = 'abcde'.split(/ /)")
	c.FoldSame("x = 'abcde'.split(' ', -1)")
	// Template strings
	c.FoldSame("x = `abcdef`.split()")
	c.FoldSame("x = `abcdef ${abc}`.split()")
	c.late = true
	c.FoldSame("x = 'a b c d e'.split(' ')")
}

// Java: PeepholeReplaceKnownMethodsTest.testJoinBug()
func (c *PeepholeReplaceKnownMethodsTest) TestJoinBug() {
	c.Fold("var x = [].join();", "var x = '';")
	c.FoldSame("var x = [x].join();")
	c.FoldSame("var x = [x,y].join();")
	c.FoldSame("var x = [x,y,z].join();")
	c.FoldSame(
		"shape['matrix'] = [\n" +
			"    Number(headingCos2).toFixed(4),\n" +
			"    Number(-headingSin2).toFixed(4),\n" +
			"    Number(headingSin2 * yScale).toFixed(4),\n" +
			"    Number(headingCos2 * yScale).toFixed(4),\n" +
			"    0,\n" +
			"    0\n" +
			"  ].join()\n")
}

// Java: PeepholeReplaceKnownMethodsTest.testJoinSpread1()
func (c *PeepholeReplaceKnownMethodsTest) TestJoinSpread1() {
	c.FoldSame("var x = [...foo].join('');")
	c.FoldSame("var x = [...someMap.keys()].join('');")
	c.FoldSame("var x = [foo, ...bar].join('');")
	c.FoldSame("var x = [...foo, bar].join('');")
	c.FoldSame("var x = [...foo, 'bar'].join('');")
	c.FoldSame("var x = ['1', ...'2', '3'].join('');")
	c.FoldSame("var x = ['1', ...['2'], '3'].join('');")
}

// Java: PeepholeReplaceKnownMethodsTest.testJoinSpread2()
func (c *PeepholeReplaceKnownMethodsTest) TestJoinSpread2() {
	c.Fold("var x = [...foo].join(',');", "var x = [...foo].join();")
	c.Fold("var x = [...someMap.keys()].join(',');", "var x = [...someMap.keys()].join();")
	c.Fold("var x = [foo, ...bar].join(',');", "var x = [foo, ...bar].join();")
	c.Fold("var x = [...foo, bar].join(',');", "var x = [...foo, bar].join();")
	c.Fold("var x = [...foo, 'bar'].join(',');", "var x = [...foo, 'bar'].join();")
	c.Fold("var x = ['1', ...'2', '3'].join(',');", "var x = ['1', ...'2', '3'].join();")
	c.Fold("var x = ['1', ...['2'], '3'].join(',');", "var x = ['1', ...['2'], '3'].join();")
}

// Java: PeepholeReplaceKnownMethodsTest.testToUpper()
func (c *PeepholeReplaceKnownMethodsTest) TestToUpper() {
	c.Fold("'a'.toUpperCase()", "'A'")
	c.Fold("'A'.toUpperCase()", "'A'")
	c.Fold("'aBcDe'.toUpperCase()", "'ABCDE'")
	c.FoldSame("`abc`.toUpperCase()")
	c.FoldSame("`a ${bc}`.toUpperCase()")
	// Make sure things aren't totally broken for non-ASCII strings, non-exhaustive.
	//
	// <p>This includes things like:
	//
	// <ul>
	// <li>graphemes with multiple code-points
	// <li>graphemes represented by multiple graphemes in other cases
	// <li>graphemes whose case changes are not round-trippable
	// <li>graphemes that change case in a position sentitive way
	// </ul>
	c.Fold("'I'.toUpperCase()", "'I'")
	c.Fold("'i'.toUpperCase()", "'I'")
	c.Fold("'\u0130'.toUpperCase()", "'\u0130'")
	c.Fold("'\u0131'.toUpperCase()", "'I'")
	c.Fold("'I\u0307'.toUpperCase()", "'I\u0307'")
	c.Fold("'\u00df'.toUpperCase()", "'SS'")
	c.Fold("'SS'.toUpperCase()", "'SS'")
	c.Fold("'\u03c3'.toUpperCase()", "'\u03a3'")
	c.Fold("'\u03c3\u03c2'.toUpperCase()", "'\u03a3\u03a3'")
}

// Java: PeepholeReplaceKnownMethodsTest.testToLower()
func (c *PeepholeReplaceKnownMethodsTest) TestToLower() {
	c.Fold("'A'.toLowerCase()", "'a'")
	c.Fold("'a'.toLowerCase()", "'a'")
	c.Fold("'aBcDe'.toLowerCase()", "'abcde'")
	c.FoldSame("`ABC`.toLowerCase()")
	c.FoldSame("`A ${BC}`.toLowerCase()")
	// Make sure things aren't totally broken for non-ASCII strings, non-exhaustive.
	//
	// <p>This includes things like:
	//
	// <ul>
	// <li>graphemes with multiple code-points
	// <li>graphemes with multiple representations
	// <li>graphemes represented by multiple graphemes in other cases
	// <li>graphemes whose case changes are not round-trippable
	// <li>graphemes that change case in a position sentitive way
	// </ul>
	c.Fold("'I'.toLowerCase()", "'i'")
	c.Fold("'i'.toLowerCase()", "'i'")
	c.Fold("'\u0130'.toLowerCase()", "'i\u0307'")
	c.Fold("'\u0131'.toLowerCase()", "'\u0131'")
	c.Fold("'I\u0307'.toLowerCase()", "'i\u0307'")
	c.Fold("'\u00df'.toLowerCase()", "'\u00df'")
	c.Fold("'SS'.toLowerCase()", "'ss'")
	c.Fold("'\u03a3'.toLowerCase()", "'\u03c3'")
	c.Fold("'\u03a3\u03a3'.toLowerCase()", "'\u03c3\u03c2'")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctionsBug()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctionsBug() {
	c.FoldSame("Math[0]()")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_abs()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_abs() {
	c.FoldSame("Math.abs(Math.random())")
	c.Fold("Math.abs('-1')", "1")
	c.Fold("Math.abs(-2)", "2")
	c.Fold("Math.abs(null)", "0")
	c.Fold("Math.abs('')", "0")
	c.Fold("Math.abs([])", "0")
	c.Fold("Math.abs([2])", "2")
	c.Fold("Math.abs([1,2])", "NaN")
	c.Fold("Math.abs({})", "NaN")
	c.Fold("Math.abs('string');", "NaN")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_imul()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_imul() {
	c.FoldSame("Math.imul(Math.random(),2)")
	c.Fold("Math.imul(-1,1)", "-1")
	c.Fold("Math.imul(2,2)", "4")
	c.Fold("Math.imul(2)", "0")
	c.Fold("Math.imul(2,3,5)", "6")
	c.Fold("Math.imul(0xfffffffe, 5)", "-10")
	c.Fold("Math.imul(0xffffffff, 5)", "-5")
	c.Fold("Math.imul(0xfffffffffffff34f, 0xfffffffffff342)", "13369344")
	c.Fold("Math.imul(0xfffffffffffff34f, -0xfffffffffff342)", "-13369344")
	c.Fold("Math.imul(NaN, 2)", "0")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_ceil()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_ceil() {
	c.FoldSame("Math.ceil(Math.random())")
	c.Fold("Math.ceil(1)", "1")
	c.Fold("Math.ceil(1.5)", "2")
	c.Fold("Math.ceil(1.3)", "2")
	c.Fold("Math.ceil(-1.3)", "-1")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_floor()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_floor() {
	c.FoldSame("Math.floor(Math.random())")
	c.Fold("Math.floor(1)", "1")
	c.Fold("Math.floor(1.5)", "1")
	c.Fold("Math.floor(1.3)", "1")
	c.Fold("Math.floor(-1.3)", "-2")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_fround()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_fround() {
	c.FoldSame("Math.fround(Math.random())")
	c.Fold("Math.fround(NaN)", "NaN")
	c.Fold("Math.fround(Infinity)", "Infinity")
	c.Fold("Math.fround(1)", "1")
	c.Fold("Math.fround(0)", "0")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_fround_j2cl()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_fround_j2cl() {
	c.FoldSame("Math.fround(1.2)")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_round()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_round() {
	c.FoldSame("Math.round(Math.random())")
	c.Fold("Math.round(NaN)", "NaN")
	c.Fold("Math.round(3.5)", "4")
	c.Fold("Math.round(-3.5)", "-3")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_sign()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_sign() {
	c.FoldSame("Math.sign(Math.random())")
	c.Fold("Math.sign(NaN)", "NaN")
	c.Fold("Math.sign(3.5)", "1")
	c.Fold("Math.sign(-3.5)", "-1")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_trunc()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_trunc() {
	c.FoldSame("Math.trunc(Math.random())")
	c.Fold("Math.sign(NaN)", "NaN")
	c.Fold("Math.trunc(3.5)", "3")
	c.Fold("Math.trunc(-3.5)", "-3")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_clz32()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_clz32() {
	c.Fold("Math.clz32(0)", "32")
	var x int32 = 1 // Java int
	for i := int32(31); i >= 0; i-- {
		c.Fold("Math.clz32("+jutil.ToJavaString(x)+")", ""+jutil.ToJavaString(i))
		c.Fold("Math.clz32("+jutil.ToJavaString(2*x-1)+")", ""+jutil.ToJavaString(i))
		x *= 2
	}
	c.Fold("Math.clz32('52')", "26")
	c.Fold("Math.clz32([52])", "26")
	c.Fold("Math.clz32([52, 53])", "32")
	// Overflow cases
	c.Fold("Math.clz32(0x100000000)", "32")
	c.Fold("Math.clz32(0x100000001)", "31")
	// NaN -> 0
	c.Fold("Math.clz32(NaN)", "32")
	c.Fold("Math.clz32('foo')", "32")
	c.Fold("Math.clz32(Infinity)", "32")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_max()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_max() {
	c.FoldSame("Math.max(Math.random(), 1)")
	c.Fold("Math.max()", "-Infinity")
	c.Fold("Math.max(0)", "0")
	c.Fold("Math.max(0, 1)", "1")
	c.Fold("Math.max(0, 1, -1, 200)", "200")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_min()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_min() {
	c.FoldSame("Math.min(Math.random(), 1)")
	c.Fold("Math.min()", "Infinity")
	c.Fold("Math.min(3)", "3")
	c.Fold("Math.min(0, 1)", "0")
	c.Fold("Math.min(0, 1, -1, 200)", "-1")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldMathFunctions_pow()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldMathFunctions_pow() {
	c.Fold("Math.pow(1, 2)", "1")
	c.Fold("Math.pow(2, 0)", "1")
	c.Fold("Math.pow(2, 2)", "4")
	c.Fold("Math.pow(2, 32)", "4294967296")
	c.Fold("Math.pow(Infinity, 0)", "1")
	c.Fold("Math.pow(Infinity, 1)", "Infinity")
	c.Fold("Math.pow('a', 33)", "NaN")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldNumberFunctions_isSafeInteger()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldNumberFunctions_isSafeInteger() {
	c.Fold("Number.isSafeInteger(1)", "true")
	c.Fold("Number.isSafeInteger(1.5)", "false")
	c.Fold("Number.isSafeInteger(9007199254740991)", "true")
	c.Fold("Number.isSafeInteger(9007199254740992)", "false")
	c.Fold("Number.isSafeInteger(-9007199254740991)", "true")
	c.Fold("Number.isSafeInteger(-9007199254740992)", "false")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldNumberFunctions_isFinite()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldNumberFunctions_isFinite() {
	c.Fold("Number.isFinite(1)", "true")
	c.Fold("Number.isFinite(1.5)", "true")
	c.Fold("Number.isFinite(NaN)", "false")
	c.Fold("Number.isFinite(Infinity)", "false")
	c.Fold("Number.isFinite(-Infinity)", "false")
	c.FoldSame("Number.isFinite('a')")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldNumberFunctions_isNaN()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldNumberFunctions_isNaN() {
	c.Fold("Number.isNaN(1)", "false")
	c.Fold("Number.isNaN(1.5)", "false")
	c.Fold("Number.isNaN(NaN)", "true")
	c.FoldSame("Number.isNaN('a')")
	// unknown function may have side effects
	c.FoldSame("Number.isNaN(+(void unknown()))")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldParseNumbers()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldParseNumbers() {
	// Template Strings
	c.FoldSame("x = parseInt(`123`)")
	c.FoldSame("x = parseInt(` 123`)")
	c.FoldSame("x = parseInt(`12 ${a}`)")
	c.FoldSame("x = parseFloat(`1.23`)")
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Fold("x = parseInt('123')", "x = 123")
	c.Fold("x = parseInt(' 123')", "x = 123")
	c.Fold("x = parseInt('123', 10)", "x = 123")
	c.Fold("x = parseInt('0xA')", "x = 10")
	c.Fold("x = parseInt('0xA', 16)", "x = 10")
	c.Fold("x = parseInt('07', 8)", "x = 7")
	c.Fold("x = parseInt('08')", "x = 8")
	c.Fold("x = parseInt('0')", "x = 0")
	c.Fold("x = parseInt('-0')", "x = -0")
	c.Fold("x = parseFloat('0')", "x = 0")
	c.Fold("x = parseFloat('1.23')", "x = 1.23")
	c.Fold("x = parseFloat('-1.23')", "x = -1.23")
	c.Fold("x = parseFloat('1.2300')", "x = 1.23")
	c.Fold("x = parseFloat(' 0.3333')", "x = 0.3333")
	c.Fold("x = parseFloat('0100')", "x = 100")
	c.Fold("x = parseFloat('0100.000')", "x = 100")
	// Mozilla Dev Center test cases
	c.Fold("x = parseInt(' 0xF', 16)", "x = 15")
	c.Fold("x = parseInt(' F', 16)", "x = 15")
	c.Fold("x = parseInt('17', 8)", "x = 15")
	c.Fold("x = parseInt('015', 10)", "x = 15")
	c.Fold("x = parseInt('1111', 2)", "x = 15")
	c.Fold("x = parseInt('12', 13)", "x = 15")
	c.Fold("x = parseInt(15.99, 10)", "x = 15")
	c.Fold("x = parseInt(-15.99, 10)", "x = -15")
	// Java's Integer.parseInt("-15.99", 10) throws an exception, because of the decimal point.
	c.FoldSame("x = parseInt('-15.99', 10)")
	c.Fold("x = parseFloat('3.14')", "x = 3.14")
	c.Fold("x = parseFloat(3.14)", "x = 3.14")
	c.Fold("x = parseFloat(-3.14)", "x = -3.14")
	c.Fold("x = parseFloat('-3.14')", "x = -3.14")
	c.Fold("x = parseFloat('-0')", "x = -0")
	// Valid calls - unable to fold
	c.FoldSame("x = parseInt('FXX123', 16)")
	c.FoldSame("x = parseInt('15*3', 10)")
	c.FoldSame("x = parseInt('15e2', 10)")
	c.FoldSame("x = parseInt('15px', 10)")
	c.FoldSame("x = parseInt('-0x08')")
	c.FoldSame("x = parseInt('1', -1)")
	c.FoldSame("x = parseFloat('3.14more non-digit characters')")
	c.FoldSame("x = parseFloat('314e-2')")
	c.FoldSame("x = parseFloat('0.0314E+2')")
	c.FoldSame("x = parseFloat('3.333333333333333333333333')")
	// Invalid calls
	c.FoldSame("x = parseInt('0xa', 10)")
	c.FoldSame("x = parseInt('')")
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.FoldSame("x = parseInt('08')")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldParseOctalNumbers()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldParseOctalNumbers() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Fold("x = parseInt('021', 8)", "x = 17")
	c.Fold("x = parseInt('-021', 8)", "x = -17")
}

// Java: PeepholeReplaceKnownMethodsTest.testReplaceWithCharAt()
func (c *PeepholeReplaceKnownMethodsTest) TestReplaceWithCharAt() {
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()
	c.FoldStringTyped("a.substring(0, 1)", "a.charAt(0)")
	c.FoldSameStringTyped("a.substring(-4, -3)")
	c.FoldSameStringTyped("a.substring(i, j + 1)")
	c.FoldSameStringTyped("a.substring(i, i + 1)")
	c.FoldSameStringTyped("a.substring(1, 2, 3)")
	c.FoldSameStringTyped("a.substring()")
	c.FoldSameStringTyped("a.substring(1)")
	c.FoldSameStringTyped("a.substring(1, 3, 4)")
	c.FoldSameStringTyped("a.substring(-1, 3)")
	c.FoldSameStringTyped("a.substring(2, 1)")
	c.FoldSameStringTyped("a.substring(3, 1)")
	c.FoldStringTyped("a.slice(4, 5)", "a.charAt(4)")
	c.FoldSameStringTyped("a.slice(-2, -1)")
	c.FoldStringTyped("var /** number */ i; a.slice(0, 1)", "var /** number */ i; a.charAt(0)")
	c.FoldSameStringTyped("a.slice(i, j + 1)")
	c.FoldSameStringTyped("a.slice(i, i + 1)")
	c.FoldSameStringTyped("a.slice(1, 2, 3)")
	c.FoldSameStringTyped("a.slice()")
	c.FoldSameStringTyped("a.slice(1)")
	c.FoldSameStringTyped("a.slice(1, 3, 4)")
	c.FoldSameStringTyped("a.slice(-1, 3)")
	c.FoldSameStringTyped("a.slice(2, 1)")
	c.FoldSameStringTyped("a.slice(3, 1)")
	c.FoldStringTyped("a.substr(0, 1)", "a.charAt(0)")
	c.FoldStringTyped("a.substr(2, 1)", "a.charAt(2)")
	c.FoldSameStringTyped("a.substr(-2, 1)")
	c.FoldSameStringTyped("a.substr(bar(), 1)")
	c.FoldSameStringTyped("''.substr(bar(), 1)")
	c.FoldSameStringTyped("a.substr(2, 1, 3)")
	c.FoldSameStringTyped("a.substr(1, 2, 3)")
	c.FoldSameStringTyped("a.substr()")
	c.FoldSameStringTyped("a.substr(1)")
	c.FoldSameStringTyped("a.substr(1, 2)")
	c.FoldSameStringTyped("a.substr(1, 2, 3)")
	c.EnableTypeCheck()
	c.FoldSame("function f(/** ? */ a) { a.substring(0, 1); }")
	c.FoldSame("function f(/** ? */ a) { a.substr(0, 1); }")
	c.FoldSame(
		"/** @constructor */ function A() {};\n" +
			"A.prototype.substring = function(begin$jscomp$1, end$jscomp$1) {};\n" +
			"function f(/** !A */ a) { a.substring(0, 1); }\n")
	c.FoldSame(
		"/** @constructor */ function A() {};\n" +
			"A.prototype.slice = function(begin$jscomp$1, end$jscomp$1) {};\n" +
			"function f(/** !A */ a) { a.slice(0, 1); }\n")
	c.useTypes = false
	c.FoldSameStringTyped("a.substring(0, 1)")
	c.FoldSameStringTyped("a.substr(0, 1)")
	c.FoldSameStringTyped("''.substring(i, i + 1)")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldConcatChaining()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldConcatChaining() {
	c.EnableTypeCheck()
	c.Fold("[1,2].concat(1).concat(2,['abc']).concat('abc')", "[1,2].concat(1,2,['abc'],'abc')")
	c.Fold("[].concat(['abc']).concat(1).concat([2,3])", "['abc'].concat(1,[2,3])")
	// cannot fold concat based on type information
	c.FoldSame("returnArrayType().concat(returnArrayType()).concat(1).concat(2)")
	c.FoldSame("returnArrayType().concat(returnUnionType()).concat(1).concat(2)")
	c.Fold(
		"[1,2,1].concat(1).concat(returnArrayType()).concat(2)",
		"[1,2,1].concat(1).concat(returnArrayType(),2)")
	c.Fold(
		"[1].concat(1).concat(2).concat(returnArrayType())",
		"[1].concat(1,2).concat(returnArrayType())")
	c.FoldSame("[].concat(1).concat(returnArrayType())")
	c.FoldSame("obj.concat([1,2]).concat(1)")
}

// Java: PeepholeReplaceKnownMethodsTest.testRemoveArrayLiteralFromFrontOfConcat()
func (c *PeepholeReplaceKnownMethodsTest) TestRemoveArrayLiteralFromFrontOfConcat() {
	c.EnableTypeCheck()
	c.Fold("[].concat([1,2,3],1)", "[1,2,3].concat(1)")
	c.FoldSame("[1,2,3].concat(returnArrayType())")
	// Call method with the same name as Array.prototype.concat
	c.FoldSame("obj.concat([1,2,3])")
	c.FoldSame("[].concat(1,[1,2,3])")
	c.FoldSame("[].concat(1)")
	c.Fold("[].concat([1])", "[1].concat()")
	// Chained folding of empty array lit
	c.Fold("[].concat([], [1,2,3], [4])", "[1,2,3].concat([4])")
	c.Fold("[].concat([]).concat([1]).concat([2,3])", "[1].concat([2,3])")
	// Cannot fold based on type information
	c.FoldSame("[].concat(returnArrayType(),1)")
	c.FoldSame("[].concat(returnArrayType())")
	c.FoldSame("[].concat(returnUnionType())")
}

// Java: PeepholeReplaceKnownMethodsTest.testArrayOfSpread()
func (c *PeepholeReplaceKnownMethodsTest) TestArrayOfSpread() {
	c.Fold("x = Array.of(...['a', 'b', 'c'])", "x = [...['a', 'b', 'c']]")
	c.Fold("x = Array.of(...['a', 'b', 'c',])", "x = [...['a', 'b', 'c']]")
	c.Fold("x = Array.of(...['a'], ...['b', 'c'])", "x = [...['a'], ...['b', 'c']]")
	c.Fold("x = Array.of('a', ...['b', 'c'])", "x = ['a', ...['b', 'c']]")
	c.Fold("x = Array.of('a', ...['b', 'c'])", "x = ['a', ...['b', 'c']]")
}

// Java: PeepholeReplaceKnownMethodsTest.testArrayOfNoSpread()
func (c *PeepholeReplaceKnownMethodsTest) TestArrayOfNoSpread() {
	c.Fold("x = Array.of('a', 'b', 'c')", "x = ['a', 'b', 'c']")
	c.Fold("x = Array.of('a', ['b', 'c'])", "x = ['a', ['b', 'c']]")
	c.Fold("x = Array.of('a', ['b', 'c'],)", "x = ['a', ['b', 'c']]")
}

// Java: PeepholeReplaceKnownMethodsTest.testArrayOfNoArgs()
func (c *PeepholeReplaceKnownMethodsTest) TestArrayOfNoArgs() {
	c.Fold("x = Array.of()", "x = []")
}

// Java: PeepholeReplaceKnownMethodsTest.testArrayOfNoChange()
func (c *PeepholeReplaceKnownMethodsTest) TestArrayOfNoChange() {
	c.FoldSame("x = Array.of.apply(window, ['a', 'b', 'c'])")
	c.FoldSame("x = ['a', 'b', 'c']")
	c.FoldSame("x = [Array.of, 'a', 'b', 'c']")
}

// Java: PeepholeReplaceKnownMethodsTest.testFoldArrayBug()
func (c *PeepholeReplaceKnownMethodsTest) TestFoldArrayBug() {
	c.FoldSame("Array[123]()")
}

// Java: PeepholeReplaceKnownMethodsTest.foldSame(String)
func (c *PeepholeReplaceKnownMethodsTest) FoldSame(js string) {
	c.TestSame(js)
}

// Java: PeepholeReplaceKnownMethodsTest.fold(String, String)
func (c *PeepholeReplaceKnownMethodsTest) Fold(js string, expected string) {
	c.Test(js, expected)
}

// Java: PeepholeReplaceKnownMethodsTest.foldSameStringTyped(String)
func (c *PeepholeReplaceKnownMethodsTest) FoldSameStringTyped(js string) {
	c.FoldStringTyped(js, js)
}

// Java: PeepholeReplaceKnownMethodsTest.foldStringTyped(String, String)
func (c *PeepholeReplaceKnownMethodsTest) FoldStringTyped(js string, expected string) {
	c.Test(
		"function f(/** string */ a) {"+js+"}",
		"function f(/** string */ a) {"+expected+"}")
}

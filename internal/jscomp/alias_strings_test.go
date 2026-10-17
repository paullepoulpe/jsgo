// Port of com/google/javascript/jscomp/AliasStringsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/AliasStringsTest.java)

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"testing"
)

type AliasStringsTest struct {
	CompilerTestCaseBase

	hashReduction bool

	aliasStringsMode CompilerOptions_AliasStringsMode
}

const AliasStringsTest_EXTERNS = "alert"

// Java: AliasStringsTest()
func NewAliasStringsTest() *AliasStringsTest {
	c := &AliasStringsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, AliasStringsTest_EXTERNS)
	c.hashReduction = false
	c.aliasStringsMode = CompilerOptions_AliasStringsMode_ALL
	return c
}

func TestAliasStringsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewAliasStringsTest() })
}

// Java: AliasStringsTest.getProcessor(Compiler)
func (c *AliasStringsTest) GetProcessor(compiler *Compiler) CompilerPass {
	pass := NewAliasStrings(compiler, compiler.GetChunkGraph(), false, c.aliasStringsMode)
	if c.hashReduction {
		pass.unitTestHashReductionMask = 0
	}
	return pass
}

// Java: AliasStringsTest.setUp()
func (c *AliasStringsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.aliasStringsMode = CompilerOptions_AliasStringsMode_ALL
}

// Java: AliasStringsTest.testTemplateLiteral()
func (c *AliasStringsTest) TestTemplateLiteral() {
	// TODO(bradfordcsmith): Maybe implement using aliases in template literals?
	c.Test(
		"const A = 'aliasable string';\n"+
			"const B = 'aliasable string';\n"+
			"const AB = `${A}aliasable string${B}`;\n",
		"var $$S_aliasable$20string = 'aliasable string';\n"+
			"const A = $$S_aliasable$20string;\n"+
			"const B = $$S_aliasable$20string;\n"+
			"const AB = `${A}aliasable string${B}`\n")
}

// Java: AliasStringsTest.testAliasAggressively()
func (c *AliasStringsTest) TestAliasAggressively() {
	c.TestSame("function f() { return 'aliasable string'; }")

	c.aliasStringsMode = CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE
	c.Test(
		"function f() { return 'aliasable string'; }",
		"var $$S_aliasable$20string = 'aliasable string';\n"+
			"function f() { return $$S_aliasable$20string; }\n")
}

// Java: AliasStringsTest.testProtectedMessage()
func (c *AliasStringsTest) TestProtectedMessage() {
	c.Test(
		jutil.StringReplace("const A = 'aliasable string';\n"+
			"const B = 'aliasable string';\n"+
			"var MSG_A =\n"+
			"    DEFINE_MSG_CALLEE(\n"+
			"        {\n"+
			"          \"key\":    \"MSG_A\",\n"+
			"          \"msg_text\":\"aliasable string\",\n"+
			"        });\n"+
			"\n", "DEFINE_MSG_CALLEE", ReplaceMessagesConstants_DEFINE_MSG_CALLEE),
		// msg_text's string is left unchanged instead of using an alias,
		// because `ReplaceMessages` needs the literal string here.
		jutil.StringReplace("var $$S_aliasable$20string = 'aliasable string';\n"+
			"const A = $$S_aliasable$20string;\n"+
			"const B = $$S_aliasable$20string;\n"+
			"var MSG_A =\n"+
			"    DEFINE_MSG_CALLEE(\n"+
			"        {\n"+
			"          \"key\":    \"MSG_A\",\n"+
			"          \"msg_text\":\"aliasable string\",\n"+
			"        });\n"+
			"\n", "DEFINE_MSG_CALLEE", ReplaceMessagesConstants_DEFINE_MSG_CALLEE))
}

// Java: AliasStringsTest.testLongStableAlias()
func (c *AliasStringsTest) TestLongStableAlias() {
	// Check long strings get a hash code

	c.Test(
		"a='Antidisestablishmentarianism';\n"+
			"b='Antidisestablishmentarianism';\n",
		"var $$S_Antidisestablishment_e428eaa9=\n"+
			"  'Antidisestablishmentarianism';\n"+
			"a=$$S_Antidisestablishment_e428eaa9;\n"+
			"b=$$S_Antidisestablishment_e428eaa9\n")

	// Check that small changes give different hash codes

	c.Test(
		"a='AntidisestablishmentarianIsm';\n"+
			"b='AntidisestablishmentarianIsm';\n",
		"var $$S_Antidisestablishment_e4287289=\n"+
			"  'AntidisestablishmentarianIsm';\n"+
			"a=$$S_Antidisestablishment_e4287289;\n"+
			"b=$$S_Antidisestablishment_e4287289\n")

	// TODO(user): check that hash code collisions are handled.
}

// Java: AliasStringsTest.testLongStableAliasHashCollision()
func (c *AliasStringsTest) TestLongStableAliasHashCollision() {
	c.hashReduction = true

	// Check that hash code collisions generate different alias
	// variable names

	c.Test(
		"f('Antidisestablishmentarianism');\n"+
			"f('Antidisestablishmentarianism');\n"+
			"f('Antidisestablishmentarianismo');\n"+
			"f('Antidisestablishmentarianismo');\n",
		"var $$S_Antidisestablishment_0=\n"+
			"  'Antidisestablishmentarianism';\n"+
			"var $$S_Antidisestablishment_0_1=\n"+
			"  'Antidisestablishmentarianismo';\n"+
			"f($$S_Antidisestablishment_0);\n"+
			"f($$S_Antidisestablishment_0);\n"+
			"f($$S_Antidisestablishment_0_1);\n"+
			"f($$S_Antidisestablishment_0_1);\n")
}

// Java: AliasStringsTest.testStringsThatAreGlobalVarValues()
func (c *AliasStringsTest) TestStringsThatAreGlobalVarValues() {

	c.TestSame("var foo='foo'; var bar='';")

	// Regular array
	c.TestSame("var foo=['foo','bar'];")

	// Nested array
	c.TestSame("var foo=['foo',['bar']];")

	// Same string is in a global array and a local in a function
	c.TestSame("var foo=['foo', 'bar'];function bar() {return 'foo';}")

	// Regular object literal
	c.TestSame("var foo={'foo': 'bar'};")

	// Nested object literal
	c.TestSame("var foo={'foo': {'bar': 'baz'}};")

	// Same string is in a global object literal (as key) and local in a
	// function
	c.TestSame("var foo={'foo': 'bar'};function bar() {return 'foo';}")

	// Same string is in a global object literal (as value) and local in a
	// function
	c.TestSame("var foo={'foo': 'foo'};function bar() {return 'foo';}")
}

// Java: AliasStringsTest.testStringsInChunks()
func (c *AliasStringsTest) TestStringsInChunks() {

	// Aliases must be placed in the correct chunk. The alias for
	// '------adios------' must be lifted from m2 and m3 and go in the
	// common parent chunks m1

	chunks :=
		JSChunkGraphBuilder_ForBush().
			AddChunk(
				"function f(a) { alert('ffffffffffffffffffff' + 'ffffffffffffffffffff' + a); }\n" +
					"function g() { alert('ciaociaociaociaociao'); }\n").
			AddChunk(
				"f('---------hi---------');f('bye');function h(a) { alert('hhhhhhhhhhhhhhhhhhhh'\n" +
					" + 'hhhhhhhhhhhhhhhhhhhh' + a); }\n").
			AddChunk(
				"f('---------hi---------');h('ciaociaociaociaociao' +\n" +
					" '--------adios-------');(function() { alert('zzzzzzzzzzzzzzzzzzzz' +\n" +
					" 'zzzzzzzzzzzzzzzzzzzz'); })();\n").
			AddChunk(
				"f('---------hi---------'); alert('--------adios-------');\n" +
					"h('-------peaches------'); h('-------peaches------');\n").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"var $$S_ciaociaociaociaociao = 'ciaociaociaociaociao';\n"+
				"var $$S_ffffffffffffffffffff = 'ffffffffffffffffffff';\n"+
				"function f(a) {\n"+
				"  alert($$S_ffffffffffffffffffff + $$S_ffffffffffffffffffff + a);\n"+
				"}\n"+
				"function g() { alert($$S_ciaociaociaociaociao); }\n",
			// m2
			"var $$S_$2d$2d$2d$2d$2d$2d$2d$2d$2dhi$2d$2d$2d$2d$2d$2d$2d$2d$2d\n"+
				" = '---------hi---------';\n"+
				"var $$S_$2d$2d$2d$2d$2d$2d$2d$2d_adios$2d$2d$2d$2d$2d$2d$2d\n"+
				" = '--------adios-------';\n"+
				"var $$S_hhhhhhhhhhhhhhhhhhhh = 'hhhhhhhhhhhhhhhhhhhh';\n"+
				"f($$S_$2d$2d$2d$2d$2d$2d$2d$2d$2dhi$2d$2d$2d$2d$2d$2d$2d$2d$2d);\n"+
				"f('bye');\n"+
				"function h(a) {\n"+
				"  alert($$S_hhhhhhhhhhhhhhhhhhhh + $$S_hhhhhhhhhhhhhhhhhhhh + a);\n"+
				"}\n",
			// m3
			"var $$S_zzzzzzzzzzzzzzzzzzzz = 'zzzzzzzzzzzzzzzzzzzz';\n"+
				"f($$S_$2d$2d$2d$2d$2d$2d$2d$2d$2dhi$2d$2d$2d$2d$2d$2d$2d$2d$2d);\n"+
				"h($$S_ciaociaociaociaociao + \n"+
				"$$S_$2d$2d$2d$2d$2d$2d$2d$2d_adios$2d$2d$2d$2d$2d$2d$2d);\n"+
				"(function() { alert($$S_zzzzzzzzzzzzzzzzzzzz + $$S_zzzzzzzzzzzzzzzzzzzz)\n"+
				" })();\n",
			// m4
			"var $$S_$2d$2d$2d$2d$2d$2d$2dpeaches$2d$2d$2d$2d$2d$2d\n"+
				" = '-------peaches------';\n"+
				"f($$S_$2d$2d$2d$2d$2d$2d$2d$2d$2dhi$2d$2d$2d$2d$2d$2d$2d$2d$2d);\n"+
				"alert($$S_$2d$2d$2d$2d$2d$2d$2d$2d_adios$2d$2d$2d$2d$2d$2d$2d);\n"+
				"h($$S_$2d$2d$2d$2d$2d$2d$2dpeaches$2d$2d$2d$2d$2d$2d);\n"+
				"h($$S_$2d$2d$2d$2d$2d$2d$2dpeaches$2d$2d$2d$2d$2d$2d);\n"))
}

// Java: AliasStringsTest.testStringsInChunks2()
func (c *AliasStringsTest) TestStringsInChunks2() {

	// Aliases must be placed in the correct chunk. The alias for
	// '------adios------' must be lifted from m2 and m3 and go in the
	// common parent chunk m1

	chunks :=
		JSChunkGraphBuilder_ForBush().
			AddChunk("function g() { alert('ciaociaociaociaociao'); }").
			AddChunk(
				"function h(a) {\n" +
					"  alert('hhhhhhhhhhhhhhhhhhh:' + a);\n" +
					"  alert('hhhhhhhhhhhhhhhhhhh:' + a);\n" +
					"}\n").
			AddChunk("h('ciaociaociaociaociao' + 'adios');").
			AddChunk("g();").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"var $$S_ciaociaociaociaociao = 'ciaociaociaociaociao';\n"+
				"function g() { alert($$S_ciaociaociaociaociao); }\n",
			// m2
			"var $$S_hhhhhhhhhhhhhhhhhhh$3a = 'hhhhhhhhhhhhhhhhhhh:';\n"+
				"function h(a) {\n"+
				"  alert($$S_hhhhhhhhhhhhhhhhhhh$3a + a);\n"+
				"  alert($$S_hhhhhhhhhhhhhhhhhhh$3a + a);\n"+
				"}\n",
			// m3
			"h($$S_ciaociaociaociaociao + 'adios');",
			// m4
			"g();"))
}

// Java: AliasStringsTest.testAliasInCommonChunkInclusive()
func (c *AliasStringsTest) TestAliasInCommonChunkInclusive() {

	chunks :=
		JSChunkGraphBuilder_ForBush().
			AddChunk("").
			AddChunk("function g() { alert('ciaociaociaociaociao'); }").
			AddChunk("h('ciaociaociaociaociao' + 'adios');").
			AddChunk("g();").
			Build()

	// The "ciao" string is used in m1 and m2.
	// Since m2 depends on m1, we should create the alias there and not force it into m0.
	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m0
			"",
			// m1
			"var $$S_ciaociaociaociaociao = 'ciaociaociaociaociao';\n"+
				"function g() { alert($$S_ciaociaociaociaociao); }\n",
			// m2
			"h($$S_ciaociaociaociaociao + 'adios');",
			// m3
			"g();"))
}

// Java: AliasStringsTest.testEmptyChunks()
func (c *AliasStringsTest) TestEmptyChunks() {
	chunks :=
		JSChunkGraphBuilder_ForStar().
			AddChunk("").
			AddChunk("function foo() { f('goodgoodgoodgoodgood') }").
			AddChunk("function foo() { f('goodgoodgoodgoodgood') }").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m0
			"var $$S_goodgoodgoodgoodgood='goodgoodgoodgoodgood'",
			// m1
			"function foo() {f($$S_goodgoodgoodgoodgood)}",
			// m2
			"function foo() {f($$S_goodgoodgoodgoodgood)}"))
}

// Java: AliasStringsTest.testOnlyAliasLargeStrings()
func (c *AliasStringsTest) TestOnlyAliasLargeStrings() {
	c.aliasStringsMode = CompilerOptions_AliasStringsMode_LARGE

	c.Test(
		"const A = 'non aliasable string with length <= 100 characters';\n"+
			"const B = 'non aliasable string with length <= 100 characters';\n"+
			"const C = 'aliasable large string largestringlargestringlargestringlargestringlargestringlargestringlargestring!';\n"+
			"const D = 'aliasable large string largestringlargestringlargestringlargestringlargestringlargestringlargestring!';\n",
		"var $$S_aliasable$20large$20stri_6c7cf169 = 'aliasable large string largestringlargestringlargestringlargestringlargestringlargestringlargestring!';\n"+
			"const A = 'non aliasable string with length <= 100 characters';\n"+
			"const B = 'non aliasable string with length <= 100 characters';\n"+
			"const C = $$S_aliasable$20large$20stri_6c7cf169;\n"+
			"const D = $$S_aliasable$20large$20stri_6c7cf169;\n")
}

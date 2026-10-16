// Port of com/google/javascript/jscomp/ExtractPrototypeMemberDeclarationsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/ExtractPrototypeMemberDeclarationsTest.java)

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"testing"
)

type ExtractPrototypeMemberDeclarationsTest struct {
	CompilerTestCaseBase

	pattern ExtractPrototypeMemberDeclarations_Pattern
}

const ExtractPrototypeMemberDeclarationsTest_TMP = "JSCompiler_prototypeAlias"

// Java: ExtractPrototypeMemberDeclarationsTest()
func NewExtractPrototypeMemberDeclarationsTest() *ExtractPrototypeMemberDeclarationsTest {
	c := &ExtractPrototypeMemberDeclarationsTest{}
	c.InitCompilerTestCaseBase(c)
	c.pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP
	return c
}

func TestExtractPrototypeMemberDeclarationsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewExtractPrototypeMemberDeclarationsTest() })
}

// Java: ExtractPrototypeMemberDeclarationsTest.setUp()
func (c *ExtractPrototypeMemberDeclarationsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP
}

// Java: ExtractPrototypeMemberDeclarationsTest.getProcessor(Compiler)
func (c *ExtractPrototypeMemberDeclarationsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewExtractPrototypeMemberDeclarations(compiler, c.pattern)
}

// Java: ExtractPrototypeMemberDeclarationsTest.testNotEnoughPrototypeToExtract()
func (c *ExtractPrototypeMemberDeclarationsTest) TestNotEnoughPrototypeToExtract() {
	// switch statement with stuff after "return"
	for i := 0; i < 7; i++ {
		c.TestSame(c.GeneratePrototypeDeclarations("x", i))
	}
}

// Java: ExtractPrototypeMemberDeclarationsTest.testExtractingSingleClassPrototype()
func (c *ExtractPrototypeMemberDeclarationsTest) TestExtractingSingleClassPrototype() {
	c.Extract(
		c.GeneratePrototypeDeclarations("x", 7),
		c.LoadPrototype("x")+c.GenerateExtractedDeclarations(7))
}

// Java: ExtractPrototypeMemberDeclarationsTest.testNoOpOnEs6Class()
func (c *ExtractPrototypeMemberDeclarationsTest) TestNoOpOnEs6Class() {
	c.TestSame("class Example { method1() {} method2() {} }")
	c.TestSame("export class Example { method1() {} method2() {} }")
}

// Java: ExtractPrototypeMemberDeclarationsTest.testClassDefinedInBlock()
func (c *ExtractPrototypeMemberDeclarationsTest) TestClassDefinedInBlock() {
	c.Test(
		jutil.StringReplace("{\n"+
			"PROTYPE_DECLARATIONS\n"+
			"}\n", "PROTYPE_DECLARATIONS", c.GeneratePrototypeDeclarations("x", 7)),
		jutil.StringReplace(jutil.StringReplace("var ExtractPrototypeMemberDeclarationsTest_TMP;\n"+
			"{\n"+
			"ExtractPrototypeMemberDeclarationsTest_TMP = x.prototype;\n"+
			"EXTRACTED_DECLARATIONS\n"+
			"}\n", "TMP", ExtractPrototypeMemberDeclarationsTest_TMP), "EXTRACTED_DECLARATIONS", c.GenerateExtractedDeclarations(7)))
}

// Java: ExtractPrototypeMemberDeclarationsTest.testClassDefinedInFunction()
func (c *ExtractPrototypeMemberDeclarationsTest) TestClassDefinedInFunction() {
	c.TestSame(
		jutil.StringReplace("function f() {\n"+
			"PROTYPE_DECLARATIONS\n"+
			"}\n", "PROTYPE_DECLARATIONS", c.GeneratePrototypeDeclarations("x", 7)))
}

/** Currently, this does not run on classes defined in ES6 modules. */
// TODO(tbreisacher): Make this work on modules. The initial 'var' needs to go into a non-module
// script node.
// Java: ExtractPrototypeMemberDeclarationsTest.testEs6Module()
func (c *ExtractPrototypeMemberDeclarationsTest) TestEs6Module() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	importStatement := "import {someValue} from './another_module.js';"
	c.TestSame(importStatement + c.GeneratePrototypeDeclarations("x", 17))
}

// Java: ExtractPrototypeMemberDeclarationsTest.testExtractingTwoClassPrototype()
func (c *ExtractPrototypeMemberDeclarationsTest) TestExtractingTwoClassPrototype() {
	c.Extract(
		c.GeneratePrototypeDeclarations("x", 6)+c.GeneratePrototypeDeclarations("y", 6),
		c.LoadPrototype("x")+
			c.GenerateExtractedDeclarations(6)+
			c.LoadPrototype("y")+
			c.GenerateExtractedDeclarations(6))
}

// Java: ExtractPrototypeMemberDeclarationsTest.testExtractingTwoClassPrototypeInDifferentBlocks()
func (c *ExtractPrototypeMemberDeclarationsTest) TestExtractingTwoClassPrototypeInDifferentBlocks() {
	c.Extract(
		c.GeneratePrototypeDeclarations("x", 6)+
			"if (foo()) {"+
			c.GeneratePrototypeDeclarations("y", 6)+
			"}",
		c.LoadPrototype("x")+
			c.GenerateExtractedDeclarations(6)+
			"if (foo()) {"+
			c.LoadPrototype("y")+
			c.GenerateExtractedDeclarations(6)+
			"}")
}

// Java: ExtractPrototypeMemberDeclarationsTest.testNoMemberDeclarations()
func (c *ExtractPrototypeMemberDeclarationsTest) TestNoMemberDeclarations() {
	c.TestSame(
		"x.prototype = {}; x.prototype = {}; x.prototype = {};\n" +
			"x.prototype = {}; x.prototype = {}; x.prototype = {};\n" +
			"x.prototype = {}; x.prototype = {}; x.prototype = {};\n")
}

// Java: ExtractPrototypeMemberDeclarationsTest.testExtractingPrototypeWithQName()
func (c *ExtractPrototypeMemberDeclarationsTest) TestExtractingPrototypeWithQName() {
	c.Extract(
		c.GeneratePrototypeDeclarations("com.google.javascript.jscomp.x", 7),
		c.LoadPrototype("com.google.javascript.jscomp.x")+c.GenerateExtractedDeclarations(7))
}

// Java: ExtractPrototypeMemberDeclarationsTest.testInterweaved()
func (c *ExtractPrototypeMemberDeclarationsTest) TestInterweaved() {
	c.TestSame(
		"x.prototype.a=1; y.prototype.a=1;\n" +
			"x.prototype.b=1; y.prototype.b=1;\n" +
			"x.prototype.c=1; y.prototype.c=1;\n" +
			"x.prototype.d=1; y.prototype.d=1;\n" +
			"x.prototype.e=1; y.prototype.e=1;\n" +
			"x.prototype.f=1; y.prototype.f=1;\n")
}

// Java: ExtractPrototypeMemberDeclarationsTest.testExtractingPrototypeWithNestedMembers()
func (c *ExtractPrototypeMemberDeclarationsTest) TestExtractingPrototypeWithNestedMembers() {
	c.Extract(
		"x.prototype.y.a = 1;\n"+
			"x.prototype.y.b = 1;\n"+
			"x.prototype.y.c = 1;\n"+
			"x.prototype.y.d = 1;\n"+
			"x.prototype.y.e = 1;\n"+
			"x.prototype.y.f = 1;\n"+
			"x.prototype.y.g = 1;\n",
		c.LoadPrototype("x")+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".y.a = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".y.b = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".y.c = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".y.d = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".y.e = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".y.f = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".y.g = 1;")
}

// Java: ExtractPrototypeMemberDeclarationsTest.testWithDevirtualization()
func (c *ExtractPrototypeMemberDeclarationsTest) TestWithDevirtualization() {
	c.Extract(
		"x.prototype.a = 1;\n"+
			"x.prototype.b = 1;\n"+
			"function devirtualize1() { }\n"+
			"x.prototype.c = 1;\n"+
			"x.prototype.d = 1;\n"+
			"x.prototype.e = 1;\n"+
			"x.prototype.f = 1;\n"+
			"x.prototype.g = 1;\n",
		c.LoadPrototype("x")+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".a = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".b = 1;\n"+
			"function devirtualize1() { }"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".c = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".d = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".e = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".f = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".g = 1;")

	c.Extract(
		"x.prototype.a = 1;\n"+
			"x.prototype.b = 1;\n"+
			"function devirtualize1() { }\n"+
			"x.prototype.c = 1;\n"+
			"x.prototype.d = 1;\n"+
			"function devirtualize2() { }\n"+
			"x.prototype.e = 1;\n"+
			"x.prototype.f = 1;\n"+
			"function devirtualize3() { }\n"+
			"x.prototype.g = 1;\n",
		c.LoadPrototype("x")+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".a = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".b = 1;\n"+
			"function devirtualize1() { }"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".c = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".d = 1;\n"+
			"function devirtualize2() { }"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".e = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".f = 1;\n"+
			"function devirtualize3() { }"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".g = 1;")
}

// Java: ExtractPrototypeMemberDeclarationsTest.testAnonSimple()
func (c *ExtractPrototypeMemberDeclarationsTest) TestAnonSimple() {
	c.pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE

	c.Extract(
		c.GeneratePrototypeDeclarations("x", 3),
		c.GenerateExtractedDeclarations(3)+c.LoadPrototype("x"))

	c.TestSame(c.GeneratePrototypeDeclarations("x", 1))
	c.TestSame(c.GeneratePrototypeDeclarations("x", 2))

	c.Extract(
		c.GeneratePrototypeDeclarations("x", 7),
		c.GenerateExtractedDeclarations(7)+c.LoadPrototype("x"))
}

// Java: ExtractPrototypeMemberDeclarationsTest.testAnonWithDevirtualization()
func (c *ExtractPrototypeMemberDeclarationsTest) TestAnonWithDevirtualization() {
	c.pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE

	c.Extract(
		"x.prototype.a = 1;\n"+
			"x.prototype.b = 1;\n"+
			"function devirtualize() { }\n"+
			"x.prototype.c = 1;\n",
		"(function("+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			") {"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".a = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".b = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".c = 1;"+
			c.LoadPrototype("x")+
			"function devirtualize() { }")

	c.Extract(
		"x.prototype.a = 1;\n"+
			"function devirtualize1() { }\n"+
			"x.prototype.b = 1;\n"+
			"function devirtualize2() { }\n"+
			"x.prototype.c = 1;\n"+
			"function devirtualize3() { }\n",
		"(function("+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			") {"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".a = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".b = 1;"+
			ExtractPrototypeMemberDeclarationsTest_TMP+
			".c = 1;"+
			c.LoadPrototype("x")+
			"function devirtualize1() { }\n"+
			"function devirtualize2() { }\n"+
			"function devirtualize3() { }")
}

// Java: ExtractPrototypeMemberDeclarationsTest.testAnonWithSideFx()
func (c *ExtractPrototypeMemberDeclarationsTest) TestAnonWithSideFx() {
	c.pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE
	c.TestSame(
		"function foo() {};\n" +
			"foo.prototype.a1 = 1;\n" +
			"bar();;\n" +
			"foo.prototype.a2 = 2;\n" +
			"bar();;\n" +
			"foo.prototype.a3 = 3;\n" +
			"bar();;\n" +
			"foo.prototype.a4 = 4;\n" +
			"bar();;\n" +
			"foo.prototype.a5 = 5;\n" +
			"bar();;\n" +
			"foo.prototype.a6 = 6;\n" +
			"bar();;\n" +
			"foo.prototype.a7 = 7;\n" +
			"bar();\n")
}

// Java: ExtractPrototypeMemberDeclarationsTest.testNotEnoughPrototypeToExtractInChunk()
func (c *ExtractPrototypeMemberDeclarationsTest) TestNotEnoughPrototypeToExtractInChunk() {
	c.pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP
	for i := 0; i < 7; i++ {
		chunks :=
			JSChunkGraphBuilder_ForStar().
				AddChunk(c.GeneratePrototypeDeclarations("x", i)).
				AddChunk(c.GeneratePrototypeDeclarations("y", i)).
				Build()
		c.TestSame(c.Srcs(chunks))
	}
}

// Java: ExtractPrototypeMemberDeclarationsTest.testExtractingSingleClassPrototypeInChunk()
func (c *ExtractPrototypeMemberDeclarationsTest) TestExtractingSingleClassPrototypeInChunk() {
	c.pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP
	chunks :=
		JSChunkGraphBuilder_ForStar().
			AddChunk(c.GeneratePrototypeDeclarations("x", 7)).
			AddChunk(c.GeneratePrototypeDeclarations("y", 7)).
			Build()

	builderX := jutil.NewStringBuilder()
	builderY := jutil.NewStringBuilder()
	xTmp := ExtractPrototypeMemberDeclarationsTest_TMP + "0"
	yTmp := ExtractPrototypeMemberDeclarationsTest_TMP + "1"
	fileIndex0, fileIndex1 := 0, 1
	builderX.
		Append(fmt.Sprintf("var %s; %s = x.prototype;", xTmp, xTmp)).
		Append(c.GenerateExtractedDeclarationsWithFileIndex(7, &fileIndex0))
	builderY.
		Append(fmt.Sprintf("var %s; %s = y.prototype;", yTmp, yTmp)).
		Append(c.GenerateExtractedDeclarationsWithFileIndex(7, &fileIndex1))

	expectedChunks :=
		JSChunkGraphBuilder_ForStar().
			AddChunk(builderX.String()).
			AddChunk(builderY.String()).
			Build()

	c.Test(c.Srcs(chunks), c.Expected(expectedChunks))
}

// Java: ExtractPrototypeMemberDeclarationsTest.loadPrototype(String)
func (c *ExtractPrototypeMemberDeclarationsTest) LoadPrototype(qName string) string {
	if c.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP {
		return ExtractPrototypeMemberDeclarationsTest_TMP + " = " + qName + ".prototype;"
	} else {
		return "})(" + qName + ".prototype);"
	}
}

// Java: ExtractPrototypeMemberDeclarationsTest.extract(String, String)
func (c *ExtractPrototypeMemberDeclarationsTest) Extract(src string, expected string) {
	if c.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP {
		c.Test(src, "var "+ExtractPrototypeMemberDeclarationsTest_TMP+";"+expected)
	} else {
		c.Test(src, expected)
	}
}

// Java: ExtractPrototypeMemberDeclarationsTest.generatePrototypeDeclarations(String, int)
func (c *ExtractPrototypeMemberDeclarationsTest) GeneratePrototypeDeclarations(className string, num int) string {
	builder := jutil.NewStringBuilder()
	for i := 0; i < num; i++ {
		member := uint16('a' + i)
		builder.Append(c.GeneratePrototypeDeclaration(className, jutil.StringFromChar(member), jutil.StringFromChar(member)))
	}
	return builder.String()
}

// Java: ExtractPrototypeMemberDeclarationsTest.generatePrototypeDeclaration(String, String, String)
func (c *ExtractPrototypeMemberDeclarationsTest) GeneratePrototypeDeclaration(className string, member string, value string) string {
	return className + ".prototype." + member + " = " + value + ";"
}

// Java: ExtractPrototypeMemberDeclarationsTest.generateExtractedDeclarations(int)
func (c *ExtractPrototypeMemberDeclarationsTest) GenerateExtractedDeclarations(num int) string {
	return c.GenerateExtractedDeclarationsWithFileIndex(num, nil)
}

// Java: ExtractPrototypeMemberDeclarationsTest.generateExtractedDeclarations(int, Integer)
func (c *ExtractPrototypeMemberDeclarationsTest) GenerateExtractedDeclarationsWithFileIndex(num int, fileIndex *int) string {
	builder := jutil.NewStringBuilder()

	alias := ExtractPrototypeMemberDeclarationsTest_TMP
	if c.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE {
		builder.Append("(function(").Append(ExtractPrototypeMemberDeclarationsTest_TMP).Append("){")
	} else if c.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP && fileIndex != nil {
		alias += jutil.ToJavaString(*fileIndex)
	}

	for i := 0; i < num; i++ {
		member := uint16('a' + i)
		builder.Append(ExtractPrototypeMemberDeclarationsTest_GenerateExtractedDeclaration(alias, jutil.StringFromChar(member), jutil.StringFromChar(member)))
	}
	return builder.String()
}

// Java: ExtractPrototypeMemberDeclarationsTest.generateExtractedDeclaration(String, String, String)
func ExtractPrototypeMemberDeclarationsTest_GenerateExtractedDeclaration(alias string, member string, value string) string {
	return alias + "." + member + " = " + value + ";"
}

// Port of test/com/google/javascript/jscomp/SourceMapTest.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// Go-only: test-helper — SourceMapTest extends SourceMapTestCase, which
// lives in the importable package sourcemaptesting (it is shared with
// sourcemap's tests), and that package imports jscomp: so this test is in
// the external test package jscomp_test.
package jscomp_test

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
	"github.com/paullepoulpe/jsgo/internal/sourcemap/sourcemaptesting"
)

// Java: public final class SourceMapTest extends SourceMapTestCase
type SourceMapTest struct {
	sourcemaptesting.SourceMapTestCaseBase

	mappings []jscomp.SourceMap_LocationMapping
	// Go-only: language-adapter — an ImmutableMap.Builder whose keys are
	// never repeated here (buildOrThrow would throw on a repeat).
	inputMaps *jutil.LinkedHashMap[string, *jscomp.SourceMapInput]
}

// Java: SourceMapTest()
func NewSourceMapTest() *SourceMapTest {
	c := &SourceMapTest{}
	c.InitSourceMapTestCaseBase(c, nil)
	return c
}

func TestSourceMapTest(t *testing.T) {
	jscomp.RunJUnit4(t, func() jscomp.JUnitTestCase { return NewSourceMapTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *SourceMapTest) JUnitSetT(t *testing.T) { c.SetT(t) }

// JUnitBefore runs the @Before method (it overrides SourceMapTestCase's).
func (c *SourceMapTest) JUnitBefore() { c.SetUp() }

// JUnitAfter: the class has no @After method.
func (c *SourceMapTest) JUnitAfter() {}

// sourceMapTest_lambdaMapping is the lambda (location) -> "mapped/" + location.
type sourceMapTest_lambdaMapping struct {
	_ byte // Go-only: language-adapter — distinct identity
}

func (*sourceMapTest_lambdaMapping) Map(location string) (string, bool) {
	return "mapped/" + location, true
}

// Java: SourceMapTest.testPrefixReplacement1()
func (c *SourceMapTest) TestPrefixReplacement1() {
	// mapping can be used to remove a prefix
	c.mappings = []jscomp.SourceMap_LocationMapping{jscomp.NewSourceMap_PrefixLocationMapping("pre/", "")}
	c.checkSourceMap2(
		"alert(1);",
		jscomp.Compiler_JoinPathParts("pre", "file1"),
		"alert(2);",
		jscomp.Compiler_JoinPathParts("pre", "file2"),
		`{
"version":3,
"file":"testcode",
"lineCount":1,
"mappings":"A,aAAAA,KAAA,CAAM,CAAN,C,CCAAA,KAAA,CAAM,CAAN;",
"sources":["file1","file2"],
"names":["alert"]
}
`)
}

// Java: SourceMapTest.testPrefixReplacement2()
func (c *SourceMapTest) TestPrefixReplacement2() {
	// mapping can be used to replace a prefix
	c.mappings = []jscomp.SourceMap_LocationMapping{jscomp.NewSourceMap_PrefixLocationMapping("pre/file", "src")}
	c.checkSourceMap2(
		"alert(1);",
		jscomp.Compiler_JoinPathParts("pre", "file1"),
		"alert(2);",
		"pre/file2",
		`{
"version":3,
"file":"testcode",
"lineCount":1,
"mappings":"A,aAAAA,KAAA,CAAM,CAAN,C,CCAAA,KAAA,CAAM,CAAN;",
"sources":["src1","src2"],
"names":["alert"]
}
`)
}

// Java: SourceMapTest.testPrefixReplacement3()
func (c *SourceMapTest) TestPrefixReplacement3() {
	// multiple mappings can be applied
	c.mappings =
		[]jscomp.SourceMap_LocationMapping{
			jscomp.NewSourceMap_PrefixLocationMapping("file1", "x"),
			jscomp.NewSourceMap_PrefixLocationMapping("file2", "y"),
		}
	c.checkSourceMap2(
		"alert(1);",
		"file1",
		"alert(2);",
		"file2",
		`{
"version":3,
"file":"testcode",
"lineCount":1,
"mappings":"A,aAAAA,KAAA,CAAM,CAAN,C,CCAAA,KAAA,CAAM,CAAN;",
"sources":["x","y"],
"names":["alert"]
}
`)
}

// Java: SourceMapTest.testPrefixReplacement4()
func (c *SourceMapTest) TestPrefixReplacement4() {
	// first match wins
	c.mappings =
		[]jscomp.SourceMap_LocationMapping{
			jscomp.NewSourceMap_PrefixLocationMapping("file1", "x"),
			jscomp.NewSourceMap_PrefixLocationMapping("file", "y"),
		}
	c.checkSourceMap2(
		"alert(1);",
		"file1",
		"alert(2);",
		"file2",
		`{
"version":3,
"file":"testcode",
"lineCount":1,
"mappings":"A,aAAAA,KAAA,CAAM,CAAN,C,CCAAA,KAAA,CAAM,CAAN;",
"sources":["x","y2"],
"names":["alert"]
}
`)
}

// Java: SourceMapTest.testLambdaReplacement()
func (c *SourceMapTest) TestLambdaReplacement() {
	c.mappings = []jscomp.SourceMap_LocationMapping{&sourceMapTest_lambdaMapping{}}
	c.checkSourceMap2(
		"alert(1);",
		"file1",
		"alert(2);",
		"file2",
		`{
"version":3,
"file":"mapped/testcode",
"lineCount":1,
"mappings":"A,aAAAA,KAAA,CAAM,CAAN,C,CCAAA,KAAA,CAAM,CAAN;",
"sources":["mapped/file1","mapped/file2"],
"names":["alert"]
}
`)
}

// This is taken from SourceMapJsLangTest. That test can't run under J2CL because it
// uses a parameterized runner but we do want to test this basic behavior under J2CL.
//
// Java: SourceMapTest.testRepeatedCompilation()
func (c *SourceMapTest) TestRepeatedCompilation() {
	// Run compiler twice feeding sourcemaps from the first run as input to the second run.
	// This way we ensure that compiler works fine with its own sourcemaps and doesn't lose
	// important information.
	fileContent := "function foo() {} alert(foo());"
	fileName := "foo.js"
	firstCompilation := c.Compile(fileContent, fileName)
	newFileName := fileName + ".compiled"
	c.inputMaps.Put(
		newFileName,
		jscomp.NewSourceMapInput(
			jscomp.SourceFile_FromCode("sourcemap", firstCompilation.SourceMapFileContent)))

	secondCompilation := c.Compile(firstCompilation.GeneratedSource, newFileName)
	c.CheckOfInput(
		fileName,
		fileContent,
		secondCompilation.GeneratedSource,
		secondCompilation.SourceMapFileContent)
}

// Tests that when a generated intermediate file contains a source map, we
// only map the source positions that the intermediate map specifies. If a
// mapping is sourceless, it is omitted from the final source map.
//
// Java: SourceMapTest.testIntermediateFilesOmitUnmappedCode()
func (c *SourceMapTest) TestIntermediateFilesOmitUnmappedCode() {
	// This map has been massaged to have a sourceless mapping covering the `alert` identifier.
	// Pretend it was synthetically generated during transforms.
	generatedFile := "generated.tsx.js"
	generatedSource := "'use strict';function foo(){}alert(foo());"
	generatedSourceMap :=
		`{
"version":3,
"file":"testcode",
"lineCount":1,
"mappings":"A,aAAAA,QAASA,IAAG,EAAG,E,KAAG,CAAMA,GAAA,EAAN;",
"sources":["foo.js"],
"names":["foo"]
}
`

	c.inputMaps.Put(
		generatedFile, jscomp.NewSourceMapInput(jscomp.SourceFile_FromCode("sourcemap", generatedSourceMap)))

	compilation := c.Compile(generatedSource, generatedFile)
	jscomp.AssertThat(compilation.SourceMapFileContent).
		IsEqualTo(
			`{
"version":3,
"file":"testcode",
"lineCount":1,
"mappings":"A,aAAAA,QAASA,IAAG,EAAG,E,MAASA,GAAAA;",
"sources":["foo.js"],
"names":["foo"]
}
`)
}

// Tests that empty "default" source maps, like those provided by Gulp, are
// ignored and do not cause the compiler to treat the file as an intermediate
// source file.
//
// Java: SourceMapTest.testEmptySourceMapsAreIgnored()
func (c *SourceMapTest) TestEmptySourceMapsAreIgnored() {
	// https://github.com/gulp-sourcemaps/gulp-sourcemaps/blob/91b94954/src/init/index.js#L116-L124
	file := "input.js"
	code := "inputCode"
	emptySourceMap :=
		`{
"version":3,
"names":[],
"mappings":"",
"sources":["foo/bar/baz.js"],
"sourcesContent":["inputCode"],
"file":"foo/bar/baz.js"
}
`

	c.inputMaps.Put(file, jscomp.NewSourceMapInput(jscomp.SourceFile_FromCode("sourcemap", emptySourceMap)))

	compilation := c.Compile(code, file)
	jscomp.AssertThat(compilation.SourceMapFileContent).
		IsEqualTo(
			`{
"version":3,
"file":"testcode",
"lineCount":1,
"mappings":"A,aAAAA;",
"sources":["input.js"],
"names":["inputCode"]
}
`)
}

// Java: SourceMapTest.getCompilerOptions()
func (c *SourceMapTest) GetCompilerOptions() *jscomp.CompilerOptions {
	options := c.SourceMapTestCaseBase.GetCompilerOptions()
	if c.mappings != nil {
		options.SetSourceMapLocationMappings(c.mappings)
	}

	if !c.inputMaps.IsEmpty() {
		options.SetApplyInputSourceMaps(true)
		options.SetInputSourceMaps(c.inputMaps)
	}
	return options
}

// Java: SourceMapTest.setUp()
func (c *SourceMapTest) SetUp() {
	c.SourceMapTestCaseBase.SetUp()
	c.inputMaps = jutil.NewLinkedHashMap[string, *jscomp.SourceMapInput]()
}

// Java: SourceMapTest.checkSourceMap2(String, String, String, String, String)
func (c *SourceMapTest) checkSourceMap2(js1 string, file1 string, js2 string, file2 string, expectedMap string) {
	result := c.CompileOfTwo(js1, file1, &js2, &file2)
	jscomp.AssertThat(result.SourceMapFileContent).IsEqualTo(expectedMap)
	jscomp.AssertThat(c.GetSourceMap(result)).IsEqualTo(result.SourceMapFileContent)
}

// Java: SourceMapTest.getSourceMapFormat()
func (c *SourceMapTest) GetSourceMapFormat() jscomp.SourceMap_Format {
	return jscomp.SourceMap_Format_V3
}

// Java: SourceMapTest.getSourceMapConsumer()
func (c *SourceMapTest) GetSourceMapConsumer() sourcemap.SourceMapConsumer {
	return sourcemap.NewSourceMapConsumerV3()
}

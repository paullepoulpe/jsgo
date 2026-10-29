// Port of test/com/google/debugging/sourcemap/SourceMapGeneratorV3Test.java @ closure-compiler v20261021

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

package sourcemap_test

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
	"github.com/paullepoulpe/jsgo/internal/sourcemap/sourcemaptesting"
)

// Java: public final class SourceMapGeneratorV3Test extends SourceMapTestCase
type SourceMapGeneratorV3Test struct {
	sourcemaptesting.SourceMapTestCaseBase
}

// newSourceMapGeneratorV3Test is a fresh SourceMapGeneratorV3Test after its
// @Before setUp(), on the Go test t.
func newSourceMapGeneratorV3Test(t *testing.T) *SourceMapGeneratorV3Test {
	c := &SourceMapGeneratorV3Test{}
	c.InitSourceMapTestCaseBase(c, t)
	c.SetUp()
	return c
}

// Java: SourceMapGeneratorV3Test.getSourceMapConsumer()
func (c *SourceMapGeneratorV3Test) GetSourceMapConsumer() sourcemap.SourceMapConsumer {
	return sourcemap.NewSourceMapConsumerV3()
}

// Java: SourceMapGeneratorV3Test.getSourceMapFormat()
func (c *SourceMapGeneratorV3Test) GetSourceMapFormat() jscomp.SourceMap_Format {
	return jscomp.SourceMap_Format_V3
}

func getEncodedFileName() string {
	// File.separatorChar is '/' on the platforms the port supports.
	return "c:\\myfile.js"
}

func TestBasicMapping1(t *testing.T) {
	newSourceMapGeneratorV3Test(t).CompileAndCheck("function __BASIC__() { }")
}

func TestBasicMappingGoldenOutput(t *testing.T) {
	// Empty source map test
	newSourceMapGeneratorV3Test(t).CheckSourceMap(
		"function __BASIC__() { }",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings("A,aAAAA,QAASA,UAAS,EAAG;").
			setSources("testcode").
			setNames("__BASIC__").
			build())
}

func TestBasicMapping2(t *testing.T) {
	newSourceMapGeneratorV3Test(t).CompileAndCheck("function __BASIC__(__PARAM1__) {}")
}

func TestLiteralMappings(t *testing.T) {
	newSourceMapGeneratorV3Test(t).CompileAndCheck("function __BASIC__(__PARAM1__, __PARAM2__) { var __VAR__ = '__STR__'; }")
}

func TestLiteralMappingsGoldenOutput(t *testing.T) {
	// Empty source map test
	newSourceMapGeneratorV3Test(t).CheckSourceMap(
		"function __BASIC__(__PARAM1__, __PARAM2__) { var __VAR__ = '__STR__'; }",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings("A,aAAAA,QAASA,UAAS,CAACC,UAAD,CAAaC,UAAb,CAAyB,CAAE,IAAIC,QAAU,SAAhB;").
			setSources("testcode").
			setNames("__BASIC__", "__PARAM1__", "__PARAM2__", "__VAR__").
			build())
}

func TestMultilineMapping(t *testing.T) {
	newSourceMapGeneratorV3Test(t).CompileAndCheck(
		"function __BASIC__(__PARAM1__, __PARAM2__) {\n" +
			"var __VAR__ = '__STR__';\n" +
			"var __ANO__ = \"__STR2__\";\n" +
			"}\n")
}

func TestMultilineMapping2(t *testing.T) {
	newSourceMapGeneratorV3Test(t).CompileAndCheck(
		"function __BASIC__(__PARAM1__, __PARAM2__) {\n" +
			"var __VAR__ = 1;\n" +
			"var __ANO__ = 2;\n" +
			"}\n")
}

func TestMultiFunctionMapping(t *testing.T) {
	newSourceMapGeneratorV3Test(t).CompileAndCheck(
		"function __BASIC__(__PARAM1__, __PARAM2__) {\n" +
			"var __VAR__ = '__STR__';\n" +
			"var __ANO__ = \"__STR2__\";\n" +
			"}\n" +
			"function __BASIC2__(__PARAM3__, __PARAM4__) {\n" +
			"var __VAR2__ = '__STR2__';\n" +
			"var __ANO2__ = \"__STR3__\";\n" +
			"}\n")
}

func TestGoldenOutput0(t *testing.T) {
	// Empty source map test
	newSourceMapGeneratorV3Test(t).CheckSourceMap(
		"",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings("A;").
			setSources().
			setNames().
			build())
}

func TestGoldenOutput0a(t *testing.T) {
	tc := newSourceMapGeneratorV3Test(t)
	tc.CheckSourceMap(
		"a;",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings("A,aAAAA;").
			setSources("testcode").
			setNames("a").
			build())

	tc.SourceMapIncludeSourcesContent = true

	tc.CheckSourceMap(
		"a;",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings("A,aAAAA;").
			setSources("testcode").
			setSourcesContent("a;").
			setNames("a").
			build())
}

func TestGoldenOutput1(t *testing.T) {
	tc := newSourceMapGeneratorV3Test(t)
	tc.DetailLevel = jscomp.SourceMap_DetailLevel_ALL

	tc.CheckSourceMap(
		"function f(foo, bar) { foo = foo + bar + 2; return foo; }",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings(
				"A,aAAAA,QAASA,EAAC,CAACC,GAAD,CAAMC,GAAN,"+
					"CAAW,CAAED,GAAA,CAAMA,GAAN,CAAYC,GAAZ,CAAkB,CAAG,"+
					"OAAOD,IAA9B;").
			setSources("testcode").
			setNames("f", "foo", "bar").
			build())

	tc.DetailLevel = jscomp.SourceMap_DetailLevel_SYMBOLS

	tc.CheckSourceMap(
		"function f(foo, bar) { foo = foo + bar + 2; return foo; }",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings(
				"A,aAAAA,QAASA,EAATA,CAAWC,GAAXD,CAAgBE,"+
					"GAAhBF,EAAuBC,GAAvBD,CAA6BC,GAA7BD,CAAmCE,GAAnCF,"+
					"SAAmDC,IAAnDD;").
			setSources("testcode").
			setNames("f", "foo", "bar").
			build())

	tc.SourceMapIncludeSourcesContent = true

	tc.CheckSourceMap(
		"function f(foo, bar) { foo = foo + bar + 2; return foo; }",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings(
				"A,aAAAA,QAASA,EAATA,CAAWC,GAAXD,CAAgBE,"+
					"GAAhBF,EAAuBC,GAAvBD,CAA6BC,GAA7BD,CAAmCE,GAAnCF,"+
					"SAAmDC,IAAnDD;").
			setSources("testcode").
			setSourcesContent("function f(foo, bar) { foo = foo + bar + 2; return foo; }").
			setNames("f", "foo", "bar").
			build())
}

func TestGoldenOutput2(t *testing.T) {
	newSourceMapGeneratorV3Test(t).CheckSourceMap(
		"function f(foo, bar) {\r\n\n\n\nfoo = foo + bar + foo;"+"\nreturn foo;\n}",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings(
				"A,aAAAA,QAASA,EAAC,CAACC,GAAD,CAAMC,GAAN,"+
					"CAAW,CAIrBD,GAAA,CAAMA,GAAN,CAAYC,GAAZ,CAAkBD,"+
					"GAClB,OAAOA,IALc;").
			setSources("testcode").
			setNames("f", "foo", "bar").
			build())
}

func TestGoldenOutput3(t *testing.T) {
	newSourceMapGeneratorV3Test(t).CheckSourceMapWithFileName(
		"c:\\myfile.js",
		"foo;",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings("A,aAAAA;").
			setSources(getEncodedFileName()).
			setNames("foo").
			build())
}

func TestGoldenOutput4(t *testing.T) {
	newSourceMapGeneratorV3Test(t).CheckSourceMapWithFileName(
		"c:\\myfile.js",
		"foo;   boo;   goo;",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(1).
			setMappings("A,aAAAA,GAAOC,IAAOC;").
			setSources(getEncodedFileName()).
			setNames("foo", "boo", "goo").
			build())
}

func goldenOutput5Source() string {
	var sb strings.Builder
	sb.WriteString("/** @preserve\n * this is a test.\n */\n")
	sb.WriteString("console.log(a + 'this is a really long line that will force the")
	sb.WriteString(" mapping to span multiple lines 123456789 123456789")
	for i := 0; i < 24; i++ {
		sb.WriteString(" 123456789 123456789 123456789 123456789 123456789")
	}
	sb.WriteString("' + c + d + e);")
	return sb.String()
}

func TestGoldenOutput5(t *testing.T) {
	tc := newSourceMapGeneratorV3Test(t)
	tc.DetailLevel = jscomp.SourceMap_DetailLevel_ALL

	tc.CheckSourceMapWithFileName(
		"c:\\myfile.js",
		goldenOutput5Source(),
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(6).
			setMappings(
				"A;;;;aAGAA,OAAQC,CAAAA,GAAR,CAAYC,CAAZ,CAAgB,mxCAAhB;AAAsyCC,CAAtyC,CAA0yCC,CAA1yC,CAA8yCC,CAA9yC;").
			setSources(getEncodedFileName()).
			setNames("console", "log", "a", "c", "d", "e").
			build())

	tc.DetailLevel = jscomp.SourceMap_DetailLevel_SYMBOLS

	tc.CheckSourceMapWithFileName(
		"c:\\myfile.js",
		goldenOutput5Source(),
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(6).
			setMappings(
				"A;;;;aAGAA,OAAQC,CAAAA,GAAR,CAAYC,CAAZ;AAAsyCC,CAAtyC,CAA0yCC,CAA1yC,CAA8yCC,CAA9yC;").
			setSources(getEncodedFileName()).
			setNames("console", "log", "a", "c", "d", "e").
			build())
}

func TestGoldenOutput_semanticLineBreaks(t *testing.T) {
	// Mapping code segments that contained newlines caused b/121282798. It
	// appears to be related to template literals because they can contain
	// newlines which are meaningful within a single expression, and are thus
	// preserved.
	newSourceMapGeneratorV3Test(t).CheckSourceMapWithFileName(
		"c:\\myfile.js",
		"var myMultilineTemplate = `Item ${a}\n"+
			"Item ${b}\n"+
			"Item ${c}`;\n",
		testJsonBuilder_Create().
			setVersion(3).
			setFile("testcode").
			setLineCount(3).
			setMappings("A,aAAA,IAAIA,oBAAuB,QAAOC,CAAP;OACpBC,CADoB;OAEpBC,CAFoB;").
			setSources(getEncodedFileName()).
			setNames("myMultilineTemplate", "a", "b", "c").
			build())
}

func TestBasicDeterminism(t *testing.T) {
	tc := newSourceMapGeneratorV3Test(t)
	result1 := tc.CompileOfTwo("file1", "foo;", jutil.StringPtr("file2"), jutil.StringPtr("bar;"))
	result2 := tc.CompileOfTwo("file2", "foo;", jutil.StringPtr("file1"), jutil.StringPtr("bar;"))

	map1 := tc.GetSourceMap(result1)
	map2 := tc.GetSourceMap(result2)

	// Assert that the files section of the maps are the same. The actual
	// entries will differ, so we cannot do a simple full comparison.

	// Line 5 has the file information.
	files1 := strings.Split(map1, "\n")[4]
	files2 := strings.Split(map2, "\n")[4]

	if files2 != files1 {
		t.Errorf("files2 = %q, want %q", files2, files1)
	}
}

func TestWriteMetaMap(t *testing.T) {
	var out strings.Builder
	name := "./app.js"
	appSections := []*sourcemap.SourceMapSection{
		sourcemap.SourceMapSection_ForURL("src1", 0, 0),
		sourcemap.SourceMapSection_ForURL("src2", 100, 10),
		sourcemap.SourceMapSection_ForURL("src3", 150, 5),
	}

	generator := sourcemap.NewSourceMapGeneratorV3()
	generator.AppendIndexMapTo(&out, name, appSections)

	want := "{\n" +
		"\"version\":3,\n" +
		"\"file\":\"./app.js\",\n" +
		"\"sections\":[\n" +
		"{\n" +
		"\"offset\":{\n" +
		"\"line\":0,\n" +
		"\"column\":0\n" +
		"},\n" +
		"\"url\":\"src1\"\n" +
		"},\n" +
		"{\n" +
		"\"offset\":{\n" +
		"\"line\":100,\n" +
		"\"column\":10\n" +
		"},\n" +
		"\"url\":\"src2\"\n" +
		"},\n" +
		"{\n" +
		"\"offset\":{\n" +
		"\"line\":150,\n" +
		"\"column\":5\n" +
		"},\n" +
		"\"url\":\"src3\"\n" +
		"}\n" +
		"]\n" +
		"}\n"
	if out.String() != want {
		t.Errorf("got:\n%s\nwant:\n%s", out.String(), want)
	}
}

func getEmptyMapFor(name string) string {
	var out strings.Builder
	generator := sourcemap.NewSourceMapGeneratorV3()
	generator.AppendTo(&out, &name)
	return out.String()
}

func TestWriteMetaMap2(t *testing.T) {
	var out strings.Builder
	name := "./app.js"
	appSections := []*sourcemap.SourceMapSection{
		// Map and URLs can be mixed.
		sourcemap.SourceMapSection_ForMap(getEmptyMapFor("./part.js"), 0, 0),
		sourcemap.SourceMapSection_ForURL("src2", 100, 10),
	}

	generator := sourcemap.NewSourceMapGeneratorV3()
	generator.AppendIndexMapTo(&out, name, appSections)

	want := "{\n" +
		"\"version\":3,\n" +
		"\"file\":\"./app.js\",\n" +
		"\"sections\":[\n" +
		"{\n" +
		"\"offset\":{\n" +
		"\"line\":0,\n" +
		"\"column\":0\n" +
		"},\n" +
		"\"map\":{\n" +
		"\"version\":3,\n" +
		"\"file\":\"./part.js\",\n" +
		"\"lineCount\":1,\n" +
		"\"mappings\":\";\",\n" +
		"\"sources\":[],\n" +
		"\"names\":[]\n" +
		"}\n" +
		"\n" +
		"},\n" +
		"{\n" +
		"\"offset\":{\n" +
		"\"line\":100,\n" +
		"\"column\":10\n" +
		"},\n" +
		"\"url\":\"src2\"\n" +
		"}\n" +
		"]\n" +
		"}\n"
	if out.String() != want {
		t.Errorf("got:\n%s\nwant:\n%s", out.String(), want)
	}
}

func TestParseSourceMetaMap(t *testing.T) {
	tc := newSourceMapGeneratorV3Test(t)
	const INPUT1 = "file1"
	const INPUT2 = "file2"
	inputs := []sourcemaptesting.SourceMapTestInput{
		{Name: INPUT1, Code: "var __FOO__ = 1;"},
		{Name: INPUT2, Code: "var __BAR__ = 2;"},
	}
	result1 := tc.Compile(inputs[0].Code, INPUT1)
	result2 := tc.Compile(inputs[1].Code, INPUT2)

	const MAP1 = "map1"
	const MAP2 = "map2"
	maps := map[string]string{
		MAP1: result1.SourceMapFileContent,
		MAP2: result2.SourceMapFileContent,
	}

	var sections []*sourcemap.SourceMapSection

	var output strings.Builder
	offset := appendAndCount(&output, result1.GeneratedSource)
	sections = append(sections, sourcemap.SourceMapSection_ForURL(MAP1, 0, 0))
	output.WriteString(result2.GeneratedSource)
	sections = append(sections, sourcemap.SourceMapSection_ForURL(MAP2, offset.GetLine(), offset.GetColumn()))

	generator := sourcemap.NewSourceMapGeneratorV3()
	var mapContents strings.Builder
	generator.AppendIndexMapTo(&mapContents, "out.js", sections)

	tc.CheckWithSupplier(
		inputs,
		output.String(),
		mapContents.String(),
		mapSupplierFunc(func(url string) *string {
			if v, ok := maps[url]; ok {
				return &v
			}
			return nil
		}))
}

type mapSupplierFunc func(url string) *string

func (f mapSupplierFunc) GetSourceMap(url string) *string { return f(url) }

func TestSourceMapMerging(t *testing.T) {
	tc := newSourceMapGeneratorV3Test(t)
	const INPUT1 = "file1"
	const INPUT2 = "file2"
	inputs := []sourcemaptesting.SourceMapTestInput{
		{Name: INPUT1, Code: "var __FOO__ = 1;"},
		{Name: INPUT2, Code: "var __BAR__ = 2;"},
	}
	result1 := tc.Compile(inputs[0].Code, INPUT1)
	result2 := tc.Compile(inputs[1].Code, INPUT2)

	var output strings.Builder
	offset := appendAndCount(&output, result1.GeneratedSource)
	output.WriteString(result2.GeneratedSource)

	generator := sourcemap.NewSourceMapGeneratorV3()

	generator.MergeMapSection(0, 0, result1.SourceMapFileContent)
	generator.MergeMapSection(offset.GetLine(), offset.GetColumn(), result2.SourceMapFileContent)

	var mapContents strings.Builder
	generator.AppendTo(&mapContents, jutil.StringPtr("out.js"))

	tc.Check(inputs, output.String(), mapContents.String())
}

func TestSourceMerging1(t *testing.T) {
	// construct first source map
	sourceMapGeneratorV3 := sourcemap.NewSourceMapGeneratorV3()
	sourceMapGeneratorV3.AddMapping(
		jutil.StringPtr("sourceName1"),
		jutil.StringPtr("symbolName1"),
		sourcemap.NewFilePosition(1, 1),
		sourcemap.NewFilePosition(2, 2),
		sourcemap.NewFilePosition(3, 3))
	var w strings.Builder
	sourceMapGeneratorV3.AppendTo(&w, jutil.StringPtr("foo.js.sourcemap"))
	sourceMap1 := w.String()

	// construct second source map
	sourceMapGeneratorV3 = sourcemap.NewSourceMapGeneratorV3()
	sourceMapGeneratorV3.AddMapping(
		jutil.StringPtr("sourceName2"),
		jutil.StringPtr("symbolName2"),
		sourcemap.NewFilePosition(1, 4),
		sourcemap.NewFilePosition(2, 5),
		sourcemap.NewFilePosition(3, 6))
	w.Reset()
	sourceMapGeneratorV3.AppendTo(&w, jutil.StringPtr("bar.js.sourcemap"))
	sourceMap2 := w.String()

	// combine the two
	sourceMapGeneratorV3 = sourcemap.NewSourceMapGeneratorV3()
	sourceMapGeneratorV3.MergeMapSection(0, 0, sourceMap1)
	sourceMapGeneratorV3.MergeMapSection(4, 0, sourceMap2)
	sourceMapGeneratorV3.AppendTo(&w, jutil.StringPtr("foobar.js.sourcemap"))
	combinedMap := w.String()

	// Construct the same map manually
	sourceMapGeneratorV3 = sourcemap.NewSourceMapGeneratorV3()
	sourceMapGeneratorV3.AddMapping(
		jutil.StringPtr("sourceName1"),
		jutil.StringPtr("symbolName1"),
		sourcemap.NewFilePosition(1, 1),
		sourcemap.NewFilePosition(2, 2),
		sourcemap.NewFilePosition(3, 3))
	sourceMapGeneratorV3.AddMapping(
		jutil.StringPtr("sourceName2"),
		jutil.StringPtr("symbolName2"),
		sourcemap.NewFilePosition(1, 4),
		sourcemap.NewFilePosition(6, 5),
		sourcemap.NewFilePosition(7, 6))
	sourceMapGeneratorV3.AppendTo(&w, jutil.StringPtr("foobar.js.sourcemap"))
	manuallyCombined := w.String()
	// TODO(b/62591567): these should be equal, but combinedMap is missing the
	// last mapping.
	if combinedMap == manuallyCombined {
		t.Errorf("combinedMap unexpectedly equals manuallyCombined")
	}
}

func TestSourceMapExtensions(t *testing.T) {
	// generating the json
	mapper := sourcemap.NewSourceMapGeneratorV3()
	mapper.AddExtension("x_google_foo", jutil.NewJsonObject())
	mapper.AddExtension("x_google_test", parseJsonObject("{\"number\" : 1}"))
	mapper.AddExtension("x_google_array", jutil.NewJsonArray())
	mapper.AddExtension("x_google_int", 2)
	mapper.AddExtension("x_google_str", "Some text")

	mapper.RemoveExtension("x_google_foo")
	var out strings.Builder
	mapper.AppendTo(&out, jutil.StringPtr("out.js"))

	if !mapper.HasExtension("x_google_test") {
		t.Error("x_google_test missing")
	}

	// reading & checking the extension properties
	sourceMap := parseJsonObject(out.String())

	if sourceMap.Has("x_google_foo") || sourceMap.Has("google_test") {
		t.Errorf("unexpected extension in %s", out.String())
	}
	if sourceMap.Get("x_google_test").GetAsJsonObject().Get("number").GetAsInt() != 1 {
		t.Errorf("x_google_test.number != 1")
	}
	if !sourceMap.Get("x_google_array").GetAsJsonArray().IsEmpty() {
		t.Errorf("x_google_array not empty")
	}
	if sourceMap.Get("x_google_int").GetAsInt() != 2 {
		t.Errorf("x_google_int != 2")
	}
	if sourceMap.Get("x_google_str").GetAsString() != "Some text" {
		t.Errorf("x_google_str != Some text")
	}
}

type extensionMergeActionFunc func(extensionKey string, currentValue, newValue any) any

func (f extensionMergeActionFunc) Merge(extensionKey string, currentValue, newValue any) any {
	return f(extensionKey, currentValue, newValue)
}

func TestSourceMapMergeExtensions(t *testing.T) {
	mapper := sourcemap.NewSourceMapGeneratorV3()
	mapper.MergeMapSection(
		0,
		0,
		"{\n"+
			"\"version\":3,\n"+
			"\"file\":\"testcode\",\n"+
			"\"lineCount\":1,\n"+
			"\"mappings\":\"AAAAA,a,QAASA,UAAS,EAAG;\",\n"+
			"\"sources\":[\"testcode\"],\n"+
			"\"names\":[\"__BASIC__\"],\n"+
			"\"x_company_foo\":2\n"+
			"}\n")

	if mapper.HasExtension("x_company_foo") {
		t.Error("x_company_foo merged")
	}

	mapper.AddExtension("x_company_baz", 2)

	mapper.MergeMapSectionWithMergeAction(
		0,
		0,
		"{\n"+
			"\"version\":3,\n"+
			"\"file\":\"testcode2\",\n"+
			"\"lineCount\":0,\n"+
			"\"mappings\":\"\",\n"+
			"\"sources\":[\"testcode2\"],\n"+
			"\"names\":[],\n"+
			"\"x_company_baz\":3,\n"+
			"\"x_company_bar\":false\n"+
			"}\n",
		extensionMergeActionFunc(func(extensionKey string, currentValue, newValue any) any {
			return currentValue.(int) + int(newValue.(*jutil.JsonPrimitive).GetAsInt())
		}))

	if got := mapper.GetExtension("x_company_baz"); got != 5 {
		t.Errorf("x_company_baz = %v", got)
	}
	if mapper.GetExtension("x_company_bar").(*jutil.JsonPrimitive).GetAsBoolean() {
		t.Error("x_company_bar is true")
	}
}

func TestSourceRoot(t *testing.T) {
	mapper := sourcemap.NewSourceMapGeneratorV3()

	// checking absence of sourceRoot
	var out strings.Builder
	mapper.AppendTo(&out, jutil.StringPtr("out.js"))
	mapping := parseJsonObject(out.String())

	if mapping.Get("version").GetAsInt() != 3 {
		t.Error("version != 3")
	}
	if mapping.Has("sourceRoot") {
		t.Error("unexpected sourceRoot")
	}

	out.Reset()
	mapper.SetSourceRoot("")
	mapper.AppendTo(&out, jutil.StringPtr("out2.js"))
	mapping = parseJsonObject(out.String())

	if mapping.Has("sourceRoot") {
		t.Error("unexpected empty sourceRoot")
	}

	// checking sourceRoot
	out.Reset()
	mapper.SetSourceRoot("http://url/path")
	mapper.AppendTo(&out, jutil.StringPtr("out3.js"))
	mapping = parseJsonObject(out.String())

	if mapping.Get("sourceRoot").GetAsString() != "http://url/path" {
		t.Errorf("sourceRoot = %s", mapping.Get("sourceRoot"))
	}
}

func count(js string) *sourcemap.FilePosition {
	line := 0
	column := 0
	for _, c := range jutil.StringToUTF16(js) {
		if c == '\n' {
			line++
			column = 0
		} else {
			column++
		}
	}
	return sourcemap.NewFilePosition(line, column)
}

func appendAndCount(out *strings.Builder, js string) *sourcemap.FilePosition {
	out.WriteString(js)
	return count(js)
}

func parseJsonObject(json string) *jutil.JsonObject {
	return jutil.NewGson().FromJsonObject(json)
}

// Port of com/google/javascript/jscomp/CommandLineRunnerTest.java @ closure-compiler v20261021

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

package cli

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/testjunit"
)

// CommandLineRunnerTest tests CommandLineRunner.
type CommandLineRunnerTest struct {
	t *testing.T

	temporaryFolder *commandLineRunnerTest_TemporaryFolder

	lastCompiler          *jscomp.Compiler   // nil is null
	lastCommandLineRunner *CommandLineRunner // nil is null
	exitCodes             []int
	outReader             *bytes.Buffer
	errReader             *bytes.Buffer
	filenames             map[int]string

	// If set, this will be appended to the end of the args list.
	// For testing args parsing.
	lastArg *string // nil is null

	// If set to true, uses comparison by string instead of by AST.
	useStringComparison bool

	useChunks commandLineRunnerTest_ChunkPattern

	args []string

	externs []*jscomp.SourceFile
}

func TestCommandLineRunnerTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &CommandLineRunnerTest{} })
}

// JUnitSetT is part of the JUnit 4 runner (junit4_test.go).
func (c *CommandLineRunnerTest) JUnitSetT(t *testing.T) {
	c.t = t
	// @Rule public final TemporaryFolder temporaryFolder = new TemporaryFolder();
	c.temporaryFolder = &commandLineRunnerTest_TemporaryFolder{t: t}
}

// JUnitBefore runs @Before setUp.
func (c *CommandLineRunnerTest) JUnitBefore() { c.setUp() }

// JUnitAfter runs the @After methods (none).
func (c *CommandLineRunnerTest) JUnitAfter() {}

// commandLineRunnerTest_TemporaryFolder is JUnit's TemporaryFolder rule over
// the Go test's temporary directory.
//
// Go-only: test-helper — the JUnit 4 rule.
type commandLineRunnerTest_TemporaryFolder struct {
	t    *testing.T
	root string
}

func (f *commandLineRunnerTest_TemporaryFolder) getRoot() string {
	if f.root == "" {
		f.root = f.t.TempDir()
	}
	return f.root
}

// NewFile is TemporaryFolder.newFile(String): creates the empty file.
func (f *commandLineRunnerTest_TemporaryFolder) NewFile(name string) *jutil.File {
	p := filepath.Join(f.getRoot(), name)
	if err := os.WriteFile(p, nil, 0o666); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	return jutil.NewFile(p)
}

// NewFolder is TemporaryFolder.newFolder(String).
func (f *commandLineRunnerTest_TemporaryFolder) NewFolder(name string) *jutil.File {
	p := filepath.Join(f.getRoot(), name)
	if err := os.MkdirAll(p, 0o777); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	return jutil.NewFile(p)
}

type commandLineRunnerTest_ChunkPattern int32

const (
	commandLineRunnerTest_ChunkPattern_NONE commandLineRunnerTest_ChunkPattern = iota
	commandLineRunnerTest_ChunkPattern_CHAIN
	commandLineRunnerTest_ChunkPattern_STAR
)

// Externs for the test
var commandLineRunnerTest_DEFAULT_EXTERNS = []*jscomp.SourceFile{
	jscomp.SourceFile_FromCode(
		"externs",
		`var arguments;
/**
 * @constructor
 * @param {...*} var_args
 * @throws {Error}
 */
function Function(var_args) {}
/**
 * @param {...*} var_args
 * @return {*}
 */
Function.prototype.call = function(var_args) {};
/**
 * @constructor
 * @param {...*} var_args
 * @return {!Array}
 */
function Array(var_args) {}
/**
 * @param {*=} opt_begin
 * @param {*=} opt_end
 * @return {!Array}
 * @this {Object}
 */
Array.prototype.slice = function(opt_begin, opt_end) {};
/** @constructor */ function Window() {}
/** @type {string} */ Window.prototype.name;
/** @type {Window} */ var window;
/** @constructor */ function Element() {}
Element.prototype.offsetWidth;
/** @nosideeffects */ function noSideEffects() {}
/** @param {...*} x */ function alert(x) {}
function Symbol() {}
`),
}

// Java: CommandLineRunnerTest.setUp()
func (c *CommandLineRunnerTest) setUp() {
	c.externs = commandLineRunnerTest_DEFAULT_EXTERNS
	c.filenames = map[int]string{}
	c.lastCompiler = nil
	c.lastArg = nil
	c.outReader = &bytes.Buffer{}
	c.errReader = &bytes.Buffer{}
	c.useStringComparison = false
	c.useChunks = commandLineRunnerTest_ChunkPattern_NONE
	c.args = c.args[:0]
	c.exitCodes = []int{}
}

// argsAdd is args.add(String).
func (c *CommandLineRunnerTest) argsAdd(arg string) { c.args = append(c.args, arg) }

func (c *CommandLineRunnerTest) TestStage1ErrorExitStatus() {
	// Create an input file
	srcFile := c.temporaryFolder.NewFile("input.js")
	c.writeFile(
		srcFile,
		// Intentionally incorrect type to generate a compiler error
		`/** @type {undefined} */
const x = 1;
`)
	// Create a path for the stage 1 output
	stage1Save := c.temporaryFolder.NewFile("stage1.save")
	commonFlags := []string{"--jscomp_error=checkTypes", "--js", srcFile.String()}
	// Run the compiler to generate the stage 1 save file
	stage1Flags := c.createStringList(
		commonFlags,
		[]string{
			"--filename_to_save_to",
			stage1Save.String(),
			"--segment_of_compilation_to_run",
			"CHECKS",
		})
	runner := NewCommandLineRunnerWithOutAndErr(
		c.stringListToArray(stage1Flags), jutil.NewPrintStream(c.outReader), jutil.NewPrintStream(c.errReader))
	// Expect an exit status of 1 because there is one error.
	testjunit.AssertThat(runner.DoRun()).IsEqualTo(1)
	testjunit.AssertThat(c.outReader.String()).IsEmpty()
}

func (c *CommandLineRunnerTest) Test3StageCompile() {

	// Create a message bundle to use
	msgBundle := c.temporaryFolder.NewFile("messages.xtb")
	c.writeFile(
		msgBundle,
		`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE translationbundle SYSTEM "translationbundle.dtd">
<translationbundle lang="es">
<translation id="6289482750305328564">hola</translation>
</translationbundle>
`)
	// Create test c.externs with a definition for goog.GetMsg().
	externsFile := c.temporaryFolder.NewFile("externs.js")
	c.writeFile(
		externsFile,
		`/**
 * @fileoverview test externs
 * @externs
 */
var goog = {};
/**
 * @nosideeffects
 * @param {string} msg
 * @param {Object=} placeholderReplacements
 * @param {Object=} options
 * @return {string}
 */
goog.getMsg = function(msg, placeholderReplacements, options) {};
`)
	// Create an input file
	srcFile := c.temporaryFolder.NewFile("input.js")
	c.writeFile(
		srcFile,
		`/** @desc greeting */
const MSG_HELLO = goog.getMsg('hello');
console.log(MSG_HELLO);
`)
	// Create a path for the stage 1 output
	stage1Save := c.temporaryFolder.NewFile("stage1.save")
	commonFlags := []string{
		"--compilation_level=ADVANCED_OPTIMIZATIONS",
		"--source_map_include_content",
		"--translations_file",
		msgBundle.String(),
		"--externs",
		externsFile.String(),
		"--js",
		srcFile.String()}
	// Run the compiler to generate the stage 1 save file
	stage1Flags := c.createStringList(
		commonFlags,
		[]string{
			"--filename_to_save_to",
			stage1Save.String(),
			"--segment_of_compilation_to_run",
			"CHECKS",
		})
	c.verifyFlagsAreIncompatibleWithChecksOnly(stage1Flags)
	runner := NewCommandLineRunnerWithOutAndErr(
		c.stringListToArray(stage1Flags), jutil.NewPrintStream(c.outReader), jutil.NewPrintStream(c.errReader))
	testjunit.AssertThat(runner.DoRun()).IsEqualTo(0)
	testjunit.AssertThat(c.outReader.String()).IsEmpty()
	testjunit.AssertThat(runner.GetCompiler().ToSource()).
		IsEqualTo("const MSG_HELLO=goog.getMsg(\"hello\");console.log(MSG_HELLO);")
	// Create a path for the stage 2 output
	stage2Save := c.temporaryFolder.NewFile("stage2.save")
	// run the compiler to generate the stage 2 save file
	stage2Flags := c.createStringList(
		commonFlags,
		[]string{
			"--filename_to_restore_from",
			stage1Save.String(),
			"--filename_to_save_to",
			stage2Save.String(),
			"--segment_of_compilation_to_run",
			"OPTIMIZATIONS",
		})
	c.verifyFlagsAreIncompatibleWithChecksOnly(stage2Flags)
	runner = NewCommandLineRunner(c.stringListToArray(stage2Flags))
	testjunit.AssertThat(runner.DoRun()).IsEqualTo(0)
	// During stage 2 the message is wrapped in a function call to protect it from mangling by
	// optimizations.
	testjunit.AssertThat(runner.GetCompiler().ToSource()).
		IsEqualTo(
			c.concatStrings(
				"console.log(",
				"__jscomp_define_msg__({\"key\":\"MSG_HELLO\",\"msg_text\":\"hello\"})",
				");"))
	// Create a path for the final output
	compiledFile := c.temporaryFolder.NewFile("compiled.js")
	// Create a path for the output source map
	sourceMapFile := c.temporaryFolder.NewFile("compiled.sourcemap")
	// run the compiler to generate the final output
	stage3Flags := c.createStringList(
		commonFlags,
		[]string{
			"--filename_to_restore_from",
			stage2Save.String(),
			"--segment_of_compilation_to_run",
			"FINALIZATIONS",
			"--js_output_file",
			compiledFile.String(),
			"--create_source_map",
			sourceMapFile.String(),
		})
	c.verifyFlagsAreIncompatibleWithChecksOnly(stage3Flags)
	runner = NewCommandLineRunner(c.stringListToArray(stage3Flags))
	testjunit.AssertThat(runner.DoRun()).IsEqualTo(0)
	// During stage 3 the message is actually replaced and the output written to the compiled
	// output file.
	compiledJs := jutil.Files_ReadStringUTF8(compiledFile.ToPath())
	testjunit.AssertThat(compiledJs).IsEqualTo(`console.log("hola");
`)
	expectedSourceMap := jutil.NewJsonObject()
	expectedSourceMap.Add("version", jutil.NewJsonPrimitiveNumber(jutil.Integer(3)))
	expectedSourceMap.AddPropertyString("file", jutil.StringPtr(compiledFile.GetAbsolutePath()))
	expectedSourceMap.Add("lineCount", jutil.NewJsonPrimitiveNumber(jutil.Integer(1)))
	expectedSourceMap.AddPropertyString("mappings", jutil.StringPtr("AAEAA,OAAQC,CAAAA,GAAR,CADkBC,MAClB;"))
	expectedSourceMap.Add("sources", c.newJsonArrayOfStrings(srcFile.GetAbsolutePath()))
	expectedSourceMap.Add(
		"sourcesContent", c.newJsonArrayOfStrings(jutil.Files_ReadStringUTF8(srcFile.ToPath())))
	expectedSourceMap.Add("names", c.newJsonArrayOfStrings("console", "log", "MSG_HELLO"))
	sourceMapText := jutil.Files_ReadStringUTF8(sourceMapFile.ToPath())
	actualSourceMap := c.getJsonObjectFromJson(sourceMapText)
	testjunit.AssertThat(actualSourceMap).IsEqualTo(expectedSourceMap)
}

func (c *CommandLineRunnerTest) TestUnknownDiagnosticGroupOnCommandLine() {
	c.argsAdd("--jscomp_error=unknownDiagnosticGroup")
	c.exitCodes = c.exitCodes[:0]
	c.compile([]string{"alert(1);"})
	// An exit code of -1 indicates that the command line options were bad
	testjunit.AssertThat(c.exitCodes).ContainsExactly(-1)
	testjunit.AssertThat(c.errReader.String()).
		Contains("Unknown diagnostic group: 'unknownDiagnosticGroup'")
}

func (c *CommandLineRunnerTest) TestUnknownAnnotation() {
	c.argsAdd("--warning_level=VERBOSE")
	c.test("/** @unknownTag */ function f() {}", jscomp.RhinoErrorReporter_BAD_JSDOC_ANNOTATION)
	c.argsAdd("--extra_annotation_name=unknownTag")
	c.testSame("/** @unknownTag */ function f() {}")
}

// See b/26884264
func (c *CommandLineRunnerTest) TestForOfTypecheck() {
	c.argsAdd("--jscomp_error=checkTypes")
	c.argsAdd("--language_in=ES6_STRICT")
	c.argsAdd("--language_out=ES3")
	c.externs =
		[]*jscomp.SourceFile{
			NewTestExternsBuilder().AddArray().AddArguments().BuildExternsFile("externs")}
	c.test(
		`class Cat {meow() {}}
class Dog {}

/** @type {!Array<!Dog>} */
var dogs = [];

for (var dog of dogs) {
  dog.meow(); // type error
}
`,
		jscomp.TypeCheck_INEXISTENT_PROPERTY)
}

func (c *CommandLineRunnerTest) TestWarningGuardOrdering1() {
	c.argsAdd("--jscomp_error=globalThis")
	c.argsAdd("--jscomp_off=globalThis")
	c.testSame("function f() { this.a = 3; }")
}

func (c *CommandLineRunnerTest) TestWarningGuardOrdering2() {
	c.argsAdd("--jscomp_off=globalThis")
	c.argsAdd("--jscomp_error=globalThis")
	c.test("function f() { this.a = 3; }", jscomp.CheckGlobalThis_GLOBAL_THIS)
}

func (c *CommandLineRunnerTest) TestWarningGuardOrdering3() {
	c.argsAdd("--jscomp_warning=globalThis")
	c.argsAdd("--jscomp_off=globalThis")
	c.testSame("function f() { this.a = 3; }")
}

func (c *CommandLineRunnerTest) TestWarningGuardOrdering4() {
	c.argsAdd("--jscomp_off=globalThis")
	c.argsAdd("--jscomp_warning=globalThis")
	c.test("function f() { this.a = 3; }", jscomp.CheckGlobalThis_GLOBAL_THIS)
}

func (c *CommandLineRunnerTest) TestWarningGuardWildcard1() {
	c.argsAdd("--jscomp_warning=*")
	c.test("/** @public */function f() { this.a = 3; }", jscomp.CheckGlobalThis_GLOBAL_THIS)
}

func (c *CommandLineRunnerTest) TestWarningGuardWildcardOrdering() {
	c.argsAdd("--jscomp_warning=*")
	c.argsAdd("--jscomp_off=globalThis")
	c.testSame("/** @public */function f() { this.a = 3; }")
}

func (c *CommandLineRunnerTest) TestWarningGuardHideWarningsFor1() {
	c.argsAdd("--jscomp_warning=globalThis")
	c.argsAdd("--hide_warnings_for=foo/bar")
	c.setFilename(0, "foo/bar.baz")
	c.testSame("function f() { this.a = 3; }")
}

func (c *CommandLineRunnerTest) TestWarningGuardHideWarningsFor2() {
	c.argsAdd("--jscomp_warning=globalThis")
	c.argsAdd("--hide_warnings_for=bar/baz")
	c.setFilename(0, "foo/bar.baz")
	c.test("function f() { this.a = 3; }", jscomp.CheckGlobalThis_GLOBAL_THIS)
}

func (c *CommandLineRunnerTest) TestSimpleModeLeavesUnusedParams() {
	c.argsAdd("--compilation_level=SIMPLE_OPTIMIZATIONS")
	c.testSame("window.f = function(a) {};")
}

func (c *CommandLineRunnerTest) TestAdvancedModeRemovesUnusedParams() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.test("window.f = function(a) {};", "window.a = function() {};")
}

func (c *CommandLineRunnerTest) TestCheckGlobalThisOffByDefault() {
	c.testSame("function f() { this.a = 3; }")
}

func (c *CommandLineRunnerTest) TestCheckGlobalThisOnWithAdvancedMode() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.test("function f() { this.a = 3; }", jscomp.CheckGlobalThis_GLOBAL_THIS)
}

func (c *CommandLineRunnerTest) TestCheckGlobalThisOnWithAdvanced() {
	c.argsAdd("-O=ADVANCED")
	c.test("function f() { this.a = 3; }", jscomp.CheckGlobalThis_GLOBAL_THIS)
}

func (c *CommandLineRunnerTest) TestCheckGlobalThisOnWithErrorFlag() {
	c.argsAdd("--jscomp_error=globalThis")
	c.test("function f() { this.a = 3; }", jscomp.CheckGlobalThis_GLOBAL_THIS)
}

func (c *CommandLineRunnerTest) TestCheckGlobalThisOff() {
	c.argsAdd("--warning_level=VERBOSE")
	c.argsAdd("--jscomp_off=globalThis")
	c.testSame("function f() { this.a = 3; }")
}

func (c *CommandLineRunnerTest) TestTypeCheckingOffByDefault() {
	c.test("function f(x) { return x; } f();", "function f(a) { return a; } f();")
}

func (c *CommandLineRunnerTest) TestReflectedMethods() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.argsAdd("--jscomp_off=checkVars")
	c.test(
		`/** @constructor */
function Foo() {}
Foo.prototype.handle = function(x, y) {
  alert(y);
};
var x = goog.reflect.object(Foo, {handle: 1});
for (var i in x) {
  x[i].call(x);
}
window['Foo'] = Foo;
`,
		`function a() {}
a.prototype.a = function(e, d) {
  alert(d);
};
var b = goog.c.b(a, {a: 1}), c;
for (c in b) {
  b[c].call(b);
}
window.Foo = a;
`)
}

func (c *CommandLineRunnerTest) TestInlineVariables() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.argsAdd("--jscomp_off=checkVars")
	// Verify local var "val" in method "bar" is not inlined over the "inc"
	// method call (which has side-effects) but "c" is inlined (which can't be
	// modified by the call).
	c.test(
		`/** @constructor */ function F() { this.a = 0; }
F.prototype.inc = function() { this.a++; return 10; };
F.prototype.bar = function() {
  var c = 3; var val = this.inc(); this.a += val + c;
};
window['f'] = new F();
window['f']['inc'] = window['f'].inc;
window['f']['bar'] = window['f'].bar;
use(window['f'].a)
`,
		`function a(){ this.a = 0; }
a.prototype.b = function(){ this.a++; return 10; };
a.prototype.c = function(){ var b=this.b(); this.a += b + 3; };
window.f = new a;
window.f.inc = window.f.b;
window.f.bar = window.f.c;
use(window.f.a);
`)
}

func (c *CommandLineRunnerTest) TestTypedAdvanced() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.argsAdd("--jscomp_warning=checkTypes")
	c.test(
		`/** @constructor */
function Foo() {}
Foo.prototype.handle1 = function(x, y) { alert(y); };
/** @constructor */
function Bar() {}
Bar.prototype.handle1 = function(x, y) {};
new Foo().handle1(1, 2);
new Bar().handle1(1, 2);
`,
		"alert(2)")
}

func (c *CommandLineRunnerTest) TestTypedDisabledAdvanced() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.argsAdd("--use_types_for_optimization=false")
	c.test(
		`/** @constructor */
function Foo() {}
Foo.prototype.handle1 = function(x, y) { alert(y); };
/** @constructor */
function Bar() {}
Bar.prototype.handle1 = function(x, y) {};
new Foo().handle1(1, 2);
new Bar().handle1(1, 2);
`,
		`function a() {}
a.prototype.a = function(c) { alert(c); };
function b() {}
b.prototype.a = function() {};
(new a).a(2);
(new b).a(2);
`)
}

func (c *CommandLineRunnerTest) TestTypeCheckingOnWithVerbose() {
	c.argsAdd("--warning_level=VERBOSE")
	c.test("function f(x) { return x; } f();", jscomp.TypeCheck_WRONG_ARGUMENT_COUNT)
}

func (c *CommandLineRunnerTest) TestTypeCheckingOnWithWVerbose() {
	c.argsAdd("-W=VERBOSE")
	c.test("function f(x) { return x; } f();", jscomp.TypeCheck_WRONG_ARGUMENT_COUNT)
}

func (c *CommandLineRunnerTest) TestTypeParsingOffByDefault() {
	c.testSame("/** @return {number */ function f(a) { return a; }")
}

func (c *CommandLineRunnerTest) TestTypeParsingOnWithVerbose() {
	c.argsAdd("--warning_level=VERBOSE")
	c.test("/** @return {number */ function f(a) { return a; }", jscomp.RhinoErrorReporter_TYPE_PARSE_ERROR)
	c.test(
		"/** @return {n} */ function f(a) { return a; }",
		jscomp.RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}

func (c *CommandLineRunnerTest) TestTypeCheckOverride1() {
	c.argsAdd("--warning_level=VERBOSE")
	c.argsAdd("--jscomp_off=checkTypes")
	c.testSame("var x = x || {}; x.f = function() {}; x.f(3);")
}

func (c *CommandLineRunnerTest) TestTypeCheckOverride2() {
	c.argsAdd("--warning_level=DEFAULT")
	c.testSame("var x = x || {}; x.f = function() {}; x.f(3);")
	c.argsAdd("--jscomp_warning=checkTypes")
	c.test("var x = x || {}; x.f = function() {}; x.f(3);", jscomp.TypeCheck_WRONG_ARGUMENT_COUNT)
}

func (c *CommandLineRunnerTest) TestCheckSymbolsOffForDefault() {
	c.argsAdd("--warning_level=DEFAULT")
	c.test("x = 3; var y; var y;", "x=3; var y;")
}

func (c *CommandLineRunnerTest) TestCheckSymbolsOnForVerbose() {
	c.argsAdd("--jscomp_error=checkVars")
	c.argsAdd("--warning_level=VERBOSE")
	c.test("x = 3;", jscomp.VarCheck_UNDEFINED_VAR_ERROR)
	c.test("var y; var y;", jscomp.VarCheck_VAR_MULTIPLY_DECLARED_ERROR)
}

func (c *CommandLineRunnerTest) TestCheckSymbolsOverrideForVerbose() {
	c.argsAdd("--warning_level=VERBOSE")
	c.argsAdd("--jscomp_off=undefinedVars")
	c.testSame("x = 3;")
}

func (c *CommandLineRunnerTest) TestCheckSymbolsOverrideForQuiet() {
	c.argsAdd("--warning_level=QUIET")
	c.argsAdd("--jscomp_error=undefinedVars")
	c.test("x = 3;", jscomp.VarCheck_UNDEFINED_VAR_ERROR)
}

func (c *CommandLineRunnerTest) TestCheckUndefinedProperties1() {
	c.argsAdd("--warning_level=VERBOSE")
	c.argsAdd("--jscomp_error=missingProperties")
	c.test("var x = {}; var y = x.bar;", jscomp.TypeCheck_INEXISTENT_PROPERTY)
}

func (c *CommandLineRunnerTest) TestCheckUndefinedProperties3() {
	c.argsAdd("--warning_level=VERBOSE")
	c.test("function f() {var x = {}; var y = x.bar;}", jscomp.TypeCheck_INEXISTENT_PROPERTY)
}

func (c *CommandLineRunnerTest) TestDuplicateParams() {
	c.test("function f(a, a) {}", jscomp.RhinoErrorReporter_DUPLICATE_PARAM)
	testjunit.AssertThat(c.lastCompiler.HasHaltingErrors()).IsTrue()
}

func (c *CommandLineRunnerTest) TestDefineFlag() {
	c.argsAdd("--define=FOO")
	c.argsAdd("--define=\"BAR=5\"")
	c.argsAdd("--D")
	c.argsAdd("CCC")
	c.argsAdd("-D")
	c.argsAdd("DDD")
	c.test(
		`/** @define {boolean} */ var FOO = false;
/** @define {number} */ var BAR = 3;
/** @define {boolean} */ var CCC = false;
/** @define {boolean} */ var DDD = false;
`,
		"var FOO = !0, BAR = 5, CCC = !0, DDD = !0;")
}

func (c *CommandLineRunnerTest) TestDefineFlag2() {
	c.argsAdd("--define=FOO='x\"'")
	c.test("/** @define {string} */ var FOO = \"a\";", "var FOO = \"x\\\"\";")
}

func (c *CommandLineRunnerTest) TestDefineFlag3() {
	c.argsAdd("--define=FOO=\"x'\"")
	c.test("/** @define {string} */ var FOO = \"a\";", "var FOO = \"x'\";")
}

func (c *CommandLineRunnerTest) TestGoogFeatureSetYearIsNotDefinedWhenBrowserFeaturesetYearFlagIsNotSupplied() {
	c.testSame("var x = 3") // input does not matter
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetDefineReplacements()).
		DoesNotContainKey("goog.FEATURESET_YEAR")
}

/**
 * Test that browser_featureset_year flag can not be used in conjunction with --language_out flag.
 */
func (c *CommandLineRunnerTest) TestBrowserFeaturesetYearFlag_usedWithLanguageOutFlag() {
	c.argsAdd("--browser_featureset_year=2019")
	c.argsAdd("--language_out=ECMASCRIPT_2020")
	e := testjunit.AssertThrows[*jscomp.FlagUsageException](func() { c.compileString("", c.args) })
	testjunit.AssertThat(e.GetMessage()).
		Contains("ERROR - both flags `--browser_featureset_year` and `--language_out` specified.")
}

/** Test that browser_featureset_year flag overrides the default goog.FEATURESET_YEAR define */
func (c *CommandLineRunnerTest) TestBrowserFeaturesetYearFlagDefinesGoogFeaturesetYear() {
	c.argsAdd("--browser_featureset_year=2019")
	original := `/** @define {number} */
goog.FEATURESET_YEAR = goog.define('goog.FEATURESET_YEAR', 2012);
`
	expected := "goog.FEATURESET_YEAR=2019"
	c.test(original, expected)
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetDefineReplacements()).
		ContainsKey("goog.FEATURESET_YEAR")
	n := c.lastCompiler.GetOptions().GetDefineReplacements().Get("goog.FEATURESET_YEAR")
	testjunit.AssertThat(n.GetDouble()).IsEqualTo(2019.0)
}

/** Minimum valid browser featureset year is 2012 */
func (c *CommandLineRunnerTest) TestInvalidBrowserFeaturesetYearFlagGeneratesError1() {
	c.argsAdd("--browser_featureset_year=2011")
	e := testjunit.AssertThrows[*jscomp.FlagUsageException](func() { c.compileString("", c.args) })
	testjunit.AssertThat(e.GetMessage()).
		Matches(
			"Illegal browser_featureset_year=2011. We support values 2012, or 2018..\\d{4} only")
}

/** Giving a browser featureset year between 2012 and 2018 reports error */
func (c *CommandLineRunnerTest) TestInvalidBrowserFeaturesetYearFlagGeneratesError2() {
	c.argsAdd("--browser_featureset_year=2015")
	e := testjunit.AssertThrows[*jscomp.FlagUsageException](func() { c.compileString("", c.args) })
	testjunit.AssertThat(e.GetMessage()).
		Matches(
			"Illegal browser_featureset_year=2015. We support values 2012, or 2018..\\d{4} only")
}

/** Giving a future year as the browser featureset year reports error */
func (c *CommandLineRunnerTest) TestInvalidBrowserFeaturesetYearFlagGeneratesError3() {
	currentYear := time.Now().Year() // Java: Calendar.getInstance().get(Calendar.YEAR)
	higherThanCurrentYear := currentYear + 1
	c.argsAdd("--browser_featureset_year=" + strconv.Itoa(higherThanCurrentYear))
	e := testjunit.AssertThrows[*jscomp.FlagUsageException](func() { c.compileString("", c.args) })
	testjunit.AssertThat(e.GetMessage()).
		Matches(
			"Illegal browser_featureset_year=" +
				strconv.Itoa(higherThanCurrentYear) +
				". We support values 2012, or 2018..\\d{4} only")
}

/** Check that --browser_featureset_year = 2012 correctly sets the language out */
func (c *CommandLineRunnerTest) TestBrowserFeatureSetYearSetsLanguageOut1() {
	c.argsAdd("--browser_featureset_year=2012")
	original := `/** @define {number} */
goog.FEATURESET_YEAR = goog.define('goog.FEATURESET_YEAR', 2012);
`
	expected := "goog.FEATURESET_YEAR=2012"
	c.test(original, expected)
	/* Browser's year is not expected to match output language's year
	   Flag value --browser_featureset_year=2012 corresponds to output ECMASCRIPT5_STRICT */
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetOutputFeatureSet()).
		IsEqualTo(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT.ToFeatureSet())
}

/** Check that --browser_featureset_year = 2019 correctly sets the language out */
func (c *CommandLineRunnerTest) TestBrowserFeatureSetYearSetsLanguageOut2() {
	c.argsAdd("--browser_featureset_year=2019")
	original := `/** @define {number} */
goog.FEATURESET_YEAR = goog.define('goog.FEATURESET_YEAR', 2012);
`
	expected := "goog.FEATURESET_YEAR=2019"
	c.test(original, expected)
	/* Browser's year is not expected to match output language's year
	   Flag value --browser_featureset_year=2019 corresponds to output ECMASCRIPT_2017 */
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetOutputFeatureSet()).
		IsEqualTo(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT_2017.ToFeatureSet())
}

/** Check that --browser_featureset_year = 2018 correctly sets the language out */
func (c *CommandLineRunnerTest) TestBrowserFeatureSetYearSetsLanguageOut3() {
	c.argsAdd("--browser_featureset_year=2018")
	original := `/** @define {number} */
goog.FEATURESET_YEAR = goog.define('goog.FEATURESET_YEAR', 2012);
`
	expected := "goog.FEATURESET_YEAR=2018"
	c.test(original, expected)
	/* Browser's year is not expected to match output language's year
	   Flag value --browser_featureset_year=2018 corresponds to output ECMASCRIPT_2016 */
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetOutputFeatureSet()).
		IsEqualTo(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT_2016.ToFeatureSet())
}

func (c *CommandLineRunnerTest) TestScriptStrictModeNoWarning() {
	c.test("'use strict';", "")
	c.test("'no use strict';", jscomp.CheckSideEffects_USELESS_CODE_ERROR)
}

func (c *CommandLineRunnerTest) TestFunctionStrictModeNoWarning() {
	c.test("function f() {'use strict';}", "function f() {}")
	c.test("function f() {'no use strict';}", jscomp.CheckSideEffects_USELESS_CODE_ERROR)
}

func (c *CommandLineRunnerTest) TestQuietMode() {
	c.argsAdd("--warning_level=DEFAULT")
	c.test(`/** @const 
 * @const */ var x;`, jscomp.RhinoErrorReporter_PARSE_ERROR)
	c.argsAdd("--warning_level=QUIET")
	c.testSame(`/** @const 
 * @const */ var x;`)
}

func (c *CommandLineRunnerTest) TestProcessClosurePrimitives() {
	c.test("var goog = {}; goog.provide('goog.dom');", "var goog = {dom:{}};")
	c.argsAdd("--process_closure_primitives=false")
	c.testSame("var goog = {}; goog.provide('goog.dom');")
}

func (c *CommandLineRunnerTest) TestGetMsgWiring() {
	c.test(
		`var goog = {}; goog.getMsg = function(x) { return x; };
/** @desc A real foo. */ var MSG_FOO = goog.getMsg('foo');
`,
		`var goog={getMsg:function(a){return a}}, MSG_FOO=goog.getMsg('foo');
`)
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.test(
		`var goog = {}; goog.getMsg = function(x) { return x; };
/** @desc A real foo. */ var MSG_FOO = goog.getMsg('foo');
window['foo'] = MSG_FOO;
`,
		"window.foo = 'foo';")
}

func (c *CommandLineRunnerTest) TestGetMsgWiringNoWarnings() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.test("/** @desc A bad foo. */ var MSG_FOO = 1;", "")
}

func (c *CommandLineRunnerTest) TestCssNameWiring() {
	c.test(
		`var goog = {}; goog.getCssName = function() {};
goog.setCssNameMapping = function() {};
goog.setCssNameMapping({'goog': 'a', 'button': 'b'});
var a = goog.getCssName('goog-button');
var b = goog.getCssName('css-button');
var c = goog.getCssName('goog-menu');
var d = goog.getCssName('css-menu');
`,
		`var goog = { getCssName: function() {},
             setCssNameMapping: function() {} },
    a = 'a-b',
    b = 'css-b',
    c = 'a-menu',
    d = 'css-menu';
`)
}

func (c *CommandLineRunnerTest) TestIssue70a() {
	c.argsAdd("--language_in=ECMASCRIPT5")
	c.test("function foo({}) {}", jscomp.RhinoErrorReporter_LANGUAGE_FEATURE)
}

func (c *CommandLineRunnerTest) TestIssue70b() {
	c.argsAdd("--language_in=ECMASCRIPT5")
	c.test("function foo([]) {}", jscomp.RhinoErrorReporter_LANGUAGE_FEATURE)
}

func (c *CommandLineRunnerTest) TestIssue81() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.argsAdd("--jscomp_off=checkVars")
	c.useStringComparison = true
	c.test("eval('1'); var x = eval; x('2');", "eval(\"1\");(0,eval)(\"2\");")
}

func (c *CommandLineRunnerTest) TestIssue115() {
	c.argsAdd("--compilation_level=SIMPLE_OPTIMIZATIONS")
	c.argsAdd("--jscomp_off=es5Strict")
	c.argsAdd("--strict_mode_input=false")
	c.argsAdd("--warning_level=VERBOSE")
	c.test(
		`function f() {
  var arguments = Array.prototype.slice.call(arguments, 0);
  return arguments[0];
}
`,
		`function f() {
  arguments = Array.prototype.slice.call(arguments, 0);
  return arguments[0];
}
`)
}

func (c *CommandLineRunnerTest) TestIssue297() {
	c.argsAdd("--compilation_level=SIMPLE_OPTIMIZATIONS")
	c.test(
		`function f(p) {
  var x;
  return ((x=p.id) && (x=parseInt(x.substr(1)))) && x>0;
}
`,
		`function f(b) {
  var a;
  return ((a=b.id) && (a=parseInt(a.substr(1)))) && a>0;
}
`)
}

func (c *CommandLineRunnerTest) TestHiddenSideEffect() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.argsAdd("--jscomp_off=checkVars")
	c.test("element.offsetWidth;", "element.offsetWidth", jscomp.CheckSideEffects_USELESS_CODE_ERROR)
}

func (c *CommandLineRunnerTest) TestIssue504() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.test(
		"void function() { alert('hi'); }();", //
		"alert('hi');")
}

func (c *CommandLineRunnerTest) TestIssue601() {
	c.argsAdd("--compilation_level=WHITESPACE_ONLY")
	c.test(
		"function f() { return '\\v' == 'v'; } window['f'] = f;",
		"function f(){return'\\v'=='v'}window['f']=f")
}

func (c *CommandLineRunnerTest) TestIssue601b() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.test(
		"function f() { return '\\v' == 'v'; } window['f'] = f;",
		"window.f=function(){return'\\v'=='v'}")
}

func (c *CommandLineRunnerTest) TestIssue601c() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.test(
		"function f() { return '\\u000B' == 'v'; } window['f'] = f;",
		"window.f=function(){return'\\u000B'=='v'}")
}

func (c *CommandLineRunnerTest) TestIssue846() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.testSame("try { new Function('this is an error'); } catch(a) { alert('x'); }")
}

func (c *CommandLineRunnerTest) TestSideEffectIntegration() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.test(
		`/** @constructor */
var Foo = function() {};
Foo.prototype.blah = function() {
  Foo.bar_(this)
};
Foo.bar_ = function(f) {
  f.x = 5;
};
var y = new Foo();
Foo.bar_({});
// We used to strip this too due to bad side-effect propagation.
y.blah();
alert(y);
`,
		"var a = new function(){}; a.a = 5; alert(a);")
}

func (c *CommandLineRunnerTest) TestDebugFlag1() {
	c.argsAdd("--compilation_level=SIMPLE_OPTIMIZATIONS")
	c.argsAdd("--debug=false")
	c.test("function foo(a) {}", "function foo(a) {}")
}

func (c *CommandLineRunnerTest) TestDebugFlag2() {
	c.argsAdd("--compilation_level=SIMPLE_OPTIMIZATIONS")
	c.argsAdd("--debug=true")
	c.test("function foo(a) {alert(a)}", "function foo($a$$) {alert($a$$)}")
}

func (c *CommandLineRunnerTest) TestDebugFlag3() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.argsAdd("--warning_level=QUIET")
	c.argsAdd("--debug=false")
	c.test(
		`function Foo() {}
Foo.x = 1;
function f() {throw new Foo().x;} f();
`,
		"throw (new function() {}).a;")
}

func (c *CommandLineRunnerTest) TestDebugFlag4() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.argsAdd("--warning_level=QUIET")
	c.argsAdd("--debug=true")
	c.test(
		`function Foo() {}
Foo.x = 1;
function f() {throw new Foo().x;} f();
`,
		"throw (new function() {}).$x$;")
}

func (c *CommandLineRunnerTest) TestBooleanFlag1() {
	c.argsAdd("--compilation_level=SIMPLE_OPTIMIZATIONS")
	c.argsAdd("--debug")
	c.test("function foo(a) {alert(a)}", "function foo($a$$) {alert($a$$)}")
}

func (c *CommandLineRunnerTest) TestBooleanFlag2() {
	c.argsAdd("--debug")
	c.argsAdd("--compilation_level=SIMPLE_OPTIMIZATIONS")
	c.test("function foo(a) {alert(a)}", "function foo($a$$) {alert($a$$)}")
}

func (c *CommandLineRunnerTest) TestHelpFlag() {
	c.argsAdd("--help")
	runner := c.createCommandLineRunner([]string{"function f() {}"})
	testjunit.AssertThat(runner.ShouldRunCompiler()).IsFalse()
	testjunit.AssertThat(runner.HasErrors()).IsFalse()
	output := c.outReader.String()
	testjunit.AssertThat(output).Contains(" --help ")
	testjunit.AssertThat(output).Contains(" --version ")
}

func (c *CommandLineRunnerTest) TestHoistedFunction1() {
	c.argsAdd("--jscomp_off=es5Strict")
	c.argsAdd("-W=VERBOSE")
	c.test("if (true) { f(); function f() {} }", jscomp.VariableReferenceCheck_EARLY_REFERENCE)
}

func (c *CommandLineRunnerTest) TestHoistedFunction2() {
	c.argsAdd("--language_out=STABLE")
	c.test("if (window) { f(); function f() {} }", "if (window) { var f = function() {}; f(); }")
}

func (c *CommandLineRunnerTest) TestExternsLifting1() {
	code := "/** @externs */ function f() {}"
	c.test([]string{code}, []string{})
	// NOTE: c.externs are [{SyntheticVarsDeclar}, c.externs, input0].
	//  - The first is added by VarCheck
	//  - The second is the default c.externs
	//  - The third is the source input that we're ensuring is lifted into c.externs
	testjunit.AssertThat(c.lastCompiler.GetExternsForTesting()).HasSize(3)
	extern := c.lastCompiler.GetExternsForTesting()[2]
	testjunit.AssertThat(extern.GetChunk()).IsNull()
	testjunit.AssertThat(extern.IsExtern()).IsTrue()
	testjunit.AssertThat(extern.GetCode()).IsEqualTo(code)
	testjunit.AssertThat(c.lastCompiler.GetInputsForTesting()).HasSize(1)
	input := c.lastCompiler.GetInputsForTesting()[0]
	testjunit.AssertThat(input.GetChunk()).IsNotNull()
	testjunit.AssertThat(input.IsExtern()).IsFalse()
	testjunit.AssertThat(input.GetCode()).IsEmpty()
}

func (c *CommandLineRunnerTest) TestExternsLifting2() {
	c.argsAdd("--warning_level=VERBOSE")
	c.test(
		[]string{"/** @externs */ function f() {}", "f(3);"},
		[]string{"f(3);"},
		jscomp.TypeCheck_WRONG_ARGUMENT_COUNT)
}

func (c *CommandLineRunnerTest) TestSourceSortingOff() {
	c.argsAdd("--compilation_level=WHITESPACE_ONLY")
	c.testSame([]string{"goog.require('beer');", "goog.provide('beer');"})
}

func (c *CommandLineRunnerTest) TestSourceSortingOn() {
	c.test(
		[]string{"goog.require('beer');", "goog.provide('beer');"},
		[]string{"var beer = {};", ""})
}

func (c *CommandLineRunnerTest) TestSourceSortingOn2() {
	c.test(
		[]string{
			"goog.provide('a');", "goog.require('a');",
		},
		[]string{"var a={};", ""})
}

func (c *CommandLineRunnerTest) TestSourceSortingOn3() {
	c.argsAdd("--dependency_mode=PRUNE_LEGACY")
	c.argsAdd("--language_in=ECMASCRIPT5")
	c.test(
		[]string{
			`goog.addDependency('sym', [], []);
var x = 3;`,
			"/** This is base.js @provideGoog */ var COMPILED = false;",
		},
		[]string{"var COMPILED = !1;", "var x = 3;"})
}

func (c *CommandLineRunnerTest) TestSourceSortingCircularDeps1() {
	c.argsAdd("--dependency_mode=PRUNE_LEGACY")
	c.argsAdd("--language_in=ECMASCRIPT5")
	c.test(
		[]string{
			"goog.provide('gin'); goog.require('tonic'); var gin = {};",
			"goog.provide('tonic'); goog.require('gin'); var tonic = {};",
			"goog.require('gin'); goog.require('tonic');",
		},
		jscomp.CheckClosureImports_LATE_PROVIDE_ERROR)
}

func (c *CommandLineRunnerTest) TestSourceSortingCircularDeps2() {
	c.argsAdd("--dependency_mode=PRUNE_LEGACY")
	c.argsAdd("--language_in=ECMASCRIPT5")
	c.test(
		[]string{
			"goog.provide('roses.lime.juice');",
			"goog.provide('gin'); goog.require('tonic'); var gin = {};",
			"goog.provide('tonic'); goog.require('gin'); var tonic = {};",
			"goog.require('gin'); goog.require('tonic');",
			"goog.provide('gimlet'); goog.require('gin'); goog.require('roses.lime.juice');",
		},
		jscomp.CheckClosureImports_LATE_PROVIDE_ERROR)
}

func (c *CommandLineRunnerTest) TestSourcePruningOn1() {
	c.argsAdd("--dependency_mode=PRUNE_LEGACY")
	c.argsAdd("--language_in=ECMASCRIPT5")
	c.test(
		[]string{
			"goog.require('beer');", "goog.provide('beer');", "goog.provide('scotch'); var x = 3;",
		},
		[]string{"var beer = {};", ""})
}

func (c *CommandLineRunnerTest) TestSourcePruningOn2() {
	c.argsAdd("--entry_point=goog:guinness")
	c.test(
		[]string{
			`goog.provide('guinness');
goog.require('beer');`,
			"goog.provide('beer');",
			"goog.provide('scotch'); var x = 3;",
		},
		[]string{"var beer = {};", "var guinness = {};"})
}

func (c *CommandLineRunnerTest) TestSourcePruningOn3() {
	c.argsAdd("--entry_point=goog:scotch")
	c.test(
		[]string{
			`goog.provide('guinness');
goog.require('beer');`,
			"goog.provide('beer');",
			"goog.provide('scotch'); var x = 3;",
		},
		[]string{
			"var scotch = {}, x = 3;",
		})
}

func (c *CommandLineRunnerTest) TestSourcePruningOn4() {
	c.argsAdd("--entry_point=goog:scotch")
	c.argsAdd("--entry_point=goog:beer")
	c.test(
		[]string{
			`goog.provide('guinness');
goog.require('beer');`,
			"goog.provide('beer');",
			"goog.provide('scotch'); var x = 3;",
		},
		[]string{"var scotch = {}, x = 3;", "var beer = {};"})
}

func (c *CommandLineRunnerTest) TestSourcePruningOn5() {
	c.argsAdd("--entry_point=goog:shiraz")
	c.test(
		[]string{
			`goog.provide('guinness');
goog.require('beer');`,
			"goog.provide('beer');",
			"goog.provide('scotch'); var x = 3;",
		},
		jscomp.Compiler_MISSING_ENTRY_ERROR)
}

func (c *CommandLineRunnerTest) TestSourcePruningOn6() {
	c.argsAdd("--entry_point=goog:scotch")
	c.test(
		[]string{
			"goog.require('beer');", "goog.provide('beer');", "goog.provide('scotch'); var x = 3;",
		},
		[]string{"var beer = {};", "", "var scotch = {}, x = 3;"})
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetDependencyOptions()).
		IsEqualTo(
			jscomp.DependencyOptions_PruneLegacyForEntryPoints(
				[]*jscomp.ModuleIdentifier{jscomp.ModuleIdentifier_ForClosure("scotch")}))
}

func (c *CommandLineRunnerTest) TestSourcePruningOn7() {
	c.argsAdd("--dependency_mode=PRUNE_LEGACY")
	c.test(
		[]string{
			"/** This is base.js @provideGoog */ var COMPILED = false;",
		},
		[]string{
			"var COMPILED = !1;",
		})
}

func (c *CommandLineRunnerTest) TestSourcePruningOn8() {
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--entry_point=goog:scotch")
	c.argsAdd("--warning_level=VERBOSE")
	c.test(
		[]string{
			`/** @externs */
var externVar;
` +
				NewTestExternsBuilder().AddClosureExterns().Build(),
			`goog.provide('scotch');
var x = externVar;
`,
		},
		[]string{
			"var scotch = {}, x = externVar;",
		})
}

func (c *CommandLineRunnerTest) TestChunkEntryPoint() {
	c.useChunks = commandLineRunnerTest_ChunkPattern_STAR
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--entry_point=goog:m1:a")
	c.test(
		[]string{"goog.provide('a');", "goog.provide('b');"},
		// Check that 'b' was stripped out, and 'a' was moved to the second
		// module (m1).
		[]string{"", "var a = {};"})
}

func (c *CommandLineRunnerTest) TestNoCompile() {
	c.argsAdd("--warning_level=VERBOSE")
	c.test(
		[]string{
			`/** @nocompile */
goog.provide('x');
var dupeVar;
`,
			"var dupeVar;",
		},
		[]string{"var dupeVar;"})
}

func (c *CommandLineRunnerTest) TestDependencySortingWhitespaceMode() {
	c.argsAdd("--dependency_mode=PRUNE_LEGACY")
	c.argsAdd("--compilation_level=WHITESPACE_ONLY")
	c.test(
		[]string{
			"goog.require('beer');",
			`goog.provide('beer');
goog.require('hops');`,
			"goog.provide('hops');",
		},
		[]string{
			"goog.provide('hops');",
			`goog.provide('beer');
goog.require('hops');`,
			"goog.require('beer');",
		})
}

func (c *CommandLineRunnerTest) TestOnlyClosureDependenciesEmptyEntryPoints() {
	// Prevents this from trying to load c.externs.zip
	c.argsAdd("--env=CUSTOM")
	c.argsAdd("--dependency_mode=PRUNE")
	runner := c.createCommandLineRunner([]string{})
	testjunit.AssertThat(runner.HasErrors()).IsTrue()
	testjunit.AssertThat(runner.ShouldRunCompiler()).IsFalse()
	testjunit.AssertThat(c.errReader.String()).
		Contains("--dependency_mode=PRUNE requires --entry_point.")
}

func (c *CommandLineRunnerTest) TestOnlyClosureDependenciesOneEntryPoint() {
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--entry_point=goog:beer")
	c.test(
		[]string{
			"goog.require('beer'); var beerRequired = 1;",
			`goog.provide('beer');
goog.require('hops');
var beerProvided = 1;`,
			"goog.provide('hops'); var hopsProvided = 1;",
			"goog.provide('scotch'); var scotchProvided = 1;",
			`goog.require('scotch');
var includeFileWithoutProvides = 1;`,
			"/** This is base.js @provideGoog */ var COMPILED = false;",
		},
		[]string{
			"var COMPILED = !1;",
			"var hops = {}, hopsProvided = 1;",
			"var beer = {}, beerProvided = 1;",
		})
}

func (c *CommandLineRunnerTest) TestOnlyClosureDependenciesOneEntryPoint_pruneAllowNoEntryPoints() {
	c.argsAdd("--dependency_mode=PRUNE_ALLOW_NO_ENTRY_POINTS")
	c.argsAdd("--entry_point=goog:beer")
	c.test(
		[]string{
			"goog.require('beer'); var beerRequired = 1;",
			`goog.provide('beer');
goog.require('hops');
var beerProvided = 1;`,
			"goog.provide('hops'); var hopsProvided = 1;",
			"goog.provide('scotch'); var scotchProvided = 1;",
			`goog.require('scotch');
var includeFileWithoutProvides = 1;`,
			"/** This is base.js @provideGoog */ var COMPILED = false;",
		},
		[]string{
			"var COMPILED = !1;",
			"var hops = {}, hopsProvided = 1;",
			"var beer = {}, beerProvided = 1;",
		})
}

func (c *CommandLineRunnerTest) TestPruneAllowNoEntryPoints_noEntryPointsPrunesEverything() {
	c.argsAdd("--dependency_mode=PRUNE_ALLOW_NO_ENTRY_POINTS")
	c.test(
		[]string{
			"goog.require('beer'); var beerRequired = 1;",
			`goog.provide('beer');
goog.require('hops');
var beerProvided = 1;`,
			"goog.provide('hops'); var hopsProvided = 1;",
			"goog.provide('scotch'); var scotchProvided = 1;",
			`goog.require('scotch');
var includeFileWithoutProvides = 1;`,
		},
		[]string{})
}

func (c *CommandLineRunnerTest) TestSourceMapExpansion1() {
	c.argsAdd("--js_output_file")
	c.argsAdd("/path/to/out.js")
	c.argsAdd("--create_source_map=%outname%.map")
	c.testSame("var x = 3;")
	testjunit.AssertThat(c.lastCommandLineRunner.ExpandSourceMapPath(c.lastCompiler.GetOptions(), nil)).
		IsEqualTo("/path/to/out.js.map")
}

func (c *CommandLineRunnerTest) TestSourceMapExpansion2() {
	c.useChunks = commandLineRunnerTest_ChunkPattern_CHAIN
	c.argsAdd("--create_source_map=%outname%.map")
	c.argsAdd("--chunk_output_path_prefix=foo")
	c.testSame([]string{"var x = 3;", "var y = 5;"})
	testjunit.AssertThat(c.lastCommandLineRunner.ExpandSourceMapPath(c.lastCompiler.GetOptions(), nil)).
		IsEqualTo("foo.map")
}

func (c *CommandLineRunnerTest) TestSourceMapExpansion3() {
	c.useChunks = commandLineRunnerTest_ChunkPattern_CHAIN
	c.argsAdd("--create_source_map=%outname%.map")
	c.argsAdd("--chunk_output_path_prefix=foo_")
	c.testSame([]string{"var x = 3;", "var y = 5;"})
	testjunit.AssertThat(
		c.lastCommandLineRunner.ExpandSourceMapPath(
			c.lastCompiler.GetOptions(), c.lastCompiler.GetChunkGraph().GetRootChunk())).
		IsEqualTo("foo_m0.js.map")
}

func (c *CommandLineRunnerTest) TestInvalidSourceMapPattern() {
	c.useChunks = commandLineRunnerTest_ChunkPattern_CHAIN
	c.argsAdd("--create_source_map=out.map")
	c.argsAdd("--chunk_output_path_prefix=foo_")
	c.test(
		[]string{"var x = 3;", "var y = 5;"},
		AbstractCommandLineRunner_INVALID_CHUNK_SOURCEMAP_PATTERN)
}

func (c *CommandLineRunnerTest) TestSourceMapFormat1() {
	c.argsAdd("--js_output_file")
	c.argsAdd("/path/to/out.js")
	c.testSame("var x = 3;")
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetSourceMapFormat()).IsEqualTo(jscomp.SourceMap_Format_DEFAULT)
}

func (c *CommandLineRunnerTest) TestSourceMapFormat2() {
	c.argsAdd("--js_output_file")
	c.argsAdd("/path/to/out.js")
	c.argsAdd("--source_map_format=V3")
	c.testSame("var x = 3;")
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetSourceMapFormat()).IsEqualTo(jscomp.SourceMap_Format_V3)
}

func (c *CommandLineRunnerTest) TestSourceMapLocationsTranslations1() {
	c.argsAdd("--js_output_file")
	c.argsAdd("/path/to/out.js")
	c.argsAdd("--create_source_map=%outname%.map")
	c.argsAdd("--source_map_location_mapping=foo/|http://bar")
	c.testSame("var x = 3;")
	mappings := c.lastCompiler.GetOptions().GetSourceMapLocationMappings()
	testjunit.AssertThat(mappings).
		ContainsExactly(jscomp.NewSourceMap_PrefixLocationMapping("foo/", "http://bar"))
}

func (c *CommandLineRunnerTest) TestSourceMapLocationsTranslations2() {
	c.argsAdd("--js_output_file")
	c.argsAdd("/path/to/out.js")
	c.argsAdd("--create_source_map=%outname%.map")
	c.argsAdd("--source_map_location_mapping=foo/|http://bar")
	c.argsAdd("--source_map_location_mapping=xxx/|http://yyy")
	c.testSame("var x = 3;")
	mappings := c.lastCompiler.GetOptions().GetSourceMapLocationMappings()
	testjunit.AssertThat(mappings).
		ContainsExactly(
			jscomp.NewSourceMap_PrefixLocationMapping("foo/", "http://bar"),
			jscomp.NewSourceMap_PrefixLocationMapping("xxx/", "http://yyy"))
}

func (c *CommandLineRunnerTest) TestSourceMapLocationsTranslations3() {
	// Prevents this from trying to load c.externs.zip
	c.argsAdd("--env=CUSTOM")
	c.argsAdd("--js_output_file")
	c.argsAdd("/path/to/out.js")
	c.argsAdd("--create_source_map=%outname%.map")
	c.argsAdd("--source_map_location_mapping=foo/")
	runner := c.createCommandLineRunner([]string{})
	testjunit.AssertThat(runner.ShouldRunCompiler()).IsFalse()
	testjunit.AssertThat(c.errReader.String()).
		Contains("Bad value for --source_map_location_mapping")
}

func (c *CommandLineRunnerTest) TestInputOneZip() {
	zip1Contents := jutil.NewLinkedHashMap[string, string]()
	zip1Contents.Put("run.js", "console.log(\"Hello World\");")
	zipFile1 := c.createZipFile(zip1Contents)
	c.compileFiles("console.log(\"Hello World\");", zipFile1)
}

func (c *CommandLineRunnerTest) TestInputMultipleZips() {
	zip1Contents := jutil.NewLinkedHashMap[string, string]()
	zip1Contents.Put("run.js", "console.log(\"Hello World\");")
	zipFile1 := c.createZipFile(zip1Contents)
	zip2Contents := jutil.NewLinkedHashMap[string, string]()
	zip2Contents.Put("run1.js", "window.alert(\"Hi Browser\");")
	zipFile2 := c.createZipFile(zip2Contents)
	c.compileFiles("console.log(\"Hello World\");window.alert(\"Hi Browser\");", zipFile1, zipFile2)
}

func (c *CommandLineRunnerTest) TestInputMultipleContents() {
	zip1Contents := jutil.NewLinkedHashMap[string, string]()
	zip1Contents.Put("a.js", "console.log(\"File A\");")
	zip1Contents.Put("b.js", "console.log(\"File B\");")
	zip1Contents.Put("c.js", "console.log(\"File C\");")
	zipFile1 := c.createZipFile(zip1Contents)
	c.compileFiles(
		"console.log(\"File A\");console.log(\"File B\");console.log(\"File C\");", zipFile1)
}

func (c *CommandLineRunnerTest) TestInputMultipleFiles() {
	zip1Contents := jutil.NewLinkedHashMap[string, string]()
	zip1Contents.Put("run.js", "console.log(\"Hello World\");")
	zipFile1 := c.createZipFile(zip1Contents)
	jsFile1 := c.createJsFile("testjsfile", "var a;")
	zip2Contents := jutil.NewLinkedHashMap[string, string]()
	zip2Contents.Put("run1.js", "window.alert(\"Hi Browser\");")
	zipFile2 := c.createZipFile(zip2Contents)
	c.compileFiles(
		"console.log(\"Hello World\");var a;window.alert(\"Hi Browser\");",
		zipFile1,
		jsFile1,
		zipFile2)
}

func (c *CommandLineRunnerTest) TestInputMultipleJsFilesWithOneJsFlag() {
	// Test that file order is preserved with --js test3.js test2.js test1.js
	jsFile1 := c.createJsFile("test1", "var a;")
	jsFile2 := c.createJsFile("test2", "var b;")
	jsFile3 := c.createJsFile("test3", "var c;")
	c.compileJsFiles("var c;var b;var a;", jsFile3, jsFile2, jsFile1)
}

func (c *CommandLineRunnerTest) TestGlobJs1() {
	jsFile1 := c.createJsFile("test1", "var a;")
	jsFile2 := c.createJsFile("test2", "var b;")
	// Move test2 to the same directory as test1, also make the filename of test2
	// lexicographically larger than test1
	testjunit.AssertThat(
		jutil.NewFile(jsFile2.GetValue()).
			RenameTo(
				jutil.NewFile(
					jutil.NewFile(jsFile1.GetValue()).GetParentFile().String() +
						jutil.File_separator +
						"utest2.js"))).
		IsTrue()
	glob := *jutil.NewFile(jsFile1.GetValue()).GetParent() + jutil.File_separator + "**.js"
	c.compileFiles("var a;var b;", NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob))
}

func (c *CommandLineRunnerTest) TestGlobJs2() {
	jsFile1 := c.createJsFile("test1", "var a;")
	jsFile2 := c.createJsFile("test2", "var b;")
	testjunit.AssertThat(
		jutil.NewFile(jsFile2.GetValue()).
			RenameTo(
				jutil.NewFile(
					jutil.NewFile(jsFile1.GetValue()).GetParentFile().String() +
						jutil.File_separator +
						"utest2.js"))).
		IsTrue()
	glob := *jutil.NewFile(jsFile1.GetValue()).GetParent() + jutil.File_separator + "*test*.js"
	c.compileFiles("var a;var b;", NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob))
}

func (c *CommandLineRunnerTest) TestGlobJs3() {
	jsFile1 := c.createJsFile("test1", "var a;")
	jsFile2 := c.createJsFile("test2", "var b;")
	testjunit.AssertThat(
		jutil.NewFile(jsFile2.GetValue()).
			RenameTo(
				jutil.NewFile(
					jutil.NewFile(jsFile1.GetValue()).GetParentFile().String() +
						jutil.File_separator +
						"test2.js"))).
		IsTrue()
	// Make sure test2.js is excluded from the inputs when the exclusion
	// comes after the inclusion
	glob1 := *jutil.NewFile(jsFile1.GetValue()).GetParent() + jutil.File_separator + "**.js"
	glob2 := "!" + *jutil.NewFile(jsFile1.GetValue()).GetParent() + jutil.File_separator + "**test2.js"
	c.compileFiles(
		"var a;", NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob1), NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob2))
}

func (c *CommandLineRunnerTest) TestGlobJs4() {
	jsFile1 := c.createJsFile("test1", "var a;")
	jsFile2 := c.createJsFile("test2", "var b;")
	testjunit.AssertThat(
		jutil.NewFile(jsFile2.GetValue()).
			RenameTo(
				jutil.NewFile(
					jutil.NewFile(jsFile1.GetValue()).GetParentFile().String() +
						jutil.File_separator +
						"test2.js"))).
		IsTrue()
	// Make sure test2.js is excluded from the inputs when the exclusion
	// comes before the inclusion
	glob1 := "!" + *jutil.NewFile(jsFile1.GetValue()).GetParent() + jutil.File_separator + "**test2.js"
	glob2 := *jutil.NewFile(jsFile1.GetValue()).GetParent() + jutil.File_separator + "**.js"
	c.compileFiles(
		"var a;", NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob1), NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob2))
}

func (c *CommandLineRunnerTest) TestGlobJs5() {
	jsFile1 := c.createJsFile("test1", "var a;")
	jsFile2 := c.createJsFile("test2", "var b;")
	temp1 := jutil.NewFile(c.t.TempDir())
	temp2 := jutil.NewFile(c.t.TempDir())
	jscompTempDir := jutil.NewFile(jsFile1.GetValue()).GetParentFile()
	newTemp1 := jutil.NewFile(jscompTempDir.String() + jutil.File_separator + "temp1")
	newTemp2 := jutil.NewFile(jscompTempDir.String() + jutil.File_separator + "temp2")
	testjunit.AssertThat(temp1.RenameTo(newTemp1)).IsTrue()
	testjunit.AssertThat(temp2.RenameTo(newTemp2)).IsTrue()
	jutil.NewFile(jsFile1.GetValue()).RenameTo(jutil.NewFile(newTemp1.String() + jutil.File_separator + "test1.js"))
	jutil.NewFile(jsFile2.GetValue()).RenameTo(jutil.NewFile(newTemp2.String() + jutil.File_separator + "test2.js"))
	// Test multiple segments with glob patterns, like /foo/bar/**/*.js
	glob := jscompTempDir.String() + jutil.File_separator + "**" + jutil.File_separator + "*.js"
	c.compileFiles("var a;var b;", NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob))
}

// @Ignore // TODO(tbreisacher): Re-enable this test when we drop Ant.
func (c *CommandLineRunnerTest) testGlobJs6() {
	jsFile1 := c.createJsFile("test1", "var a;")
	jsFile2 := c.createJsFile("test2", "var b;")
	ignoredJs := jutil.NewFile("." + jutil.File_separator + "ignored.js")
	if ignoredJs.IsDirectory() {
		for _, f := range ignoredJs.ListFiles() {
			f.Delete()
		}
	}
	ignoredJs.Delete()
	testjunit.AssertThat(jutil.NewFile(jsFile2.GetValue()).RenameTo(ignoredJs)).IsTrue()
	// Make sure patterns like "!**\./ignored**.js" work
	glob1 := "!**\\." + jutil.File_separator + "ignored**.js"
	glob2 := *jutil.NewFile(jsFile1.GetValue()).GetParent() + jutil.File_separator + "**.js"
	c.compileFiles(
		"var a;", NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob1), NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob2))
	ignoredJs.Delete()
}

// @Ignore // TODO(tbreisacher): Re-enable this test when we drop Ant.
func (c *CommandLineRunnerTest) testGlobJs7() {
	jsFile1 := c.createJsFile("test1", "var a;")
	jsFile2 := c.createJsFile("test2", "var b;")
	takenJs := jutil.NewFile("." + jutil.File_separator + "globTestTaken.js")
	ignoredJs := jutil.NewFile("." + jutil.File_separator + "globTestIgnored.js")
	if takenJs.IsDirectory() {
		for _, f := range takenJs.ListFiles() {
			f.Delete()
		}
	}
	takenJs.Delete()
	if ignoredJs.IsDirectory() {
		for _, f := range ignoredJs.ListFiles() {
			f.Delete()
		}
	}
	ignoredJs.Delete()
	testjunit.AssertThat(jutil.NewFile(jsFile1.GetValue()).RenameTo(takenJs)).IsTrue()
	testjunit.AssertThat(jutil.NewFile(jsFile2.GetValue()).RenameTo(ignoredJs)).IsTrue()
	// Make sure that relative paths like "!**ignored.js" work with absolute paths.
	glob1 := takenJs.GetParentFile().GetAbsolutePath() + jutil.File_separator + "**Taken.js"
	glob2 := "!**Ignored.js"
	c.compileFiles(
		"var a;", NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob1), NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, glob2))
	takenJs.Delete()
	ignoredJs.Delete()
}

func (c *CommandLineRunnerTest) TestSymlink() {
	jsFile1 := c.createJsFile("test1", "var a;")
	symlink1 := jutil.NewFile(c.t.TempDir()).ToPath().ResolveString("symlink1")
	jscompTempDir := jutil.Path_Of(jsFile1.GetValue()).GetParent()
	if err := os.Symlink(jscompTempDir.String(), symlink1.String()); err != nil { // Java: Files.createSymbolicLink
		panic(jutil.NewIOException(err.Error()))
	}
	c.compileFiles("var a;", NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, symlink1.String()))
}

func (c *CommandLineRunnerTest) TestSourceMapInputs() {
	c.argsAdd("--js_output_file")
	c.argsAdd("/path/to/out.js")
	c.argsAdd("--source_map_input=input1|input1.sourcemap")
	c.argsAdd("--source_map_input=input2|input2.sourcemap")
	c.testSame("var x = 3;")
	inputMaps := c.lastCompiler.GetOptions().GetInputSourceMaps()
	testjunit.AssertThat(inputMaps).HasSize(2)
	testjunit.AssertThat(inputMaps.Get("input1").GetOriginalPath()).IsEqualTo("input1.sourcemap")
	testjunit.AssertThat(inputMaps.Get("input2").GetOriginalPath()).IsEqualTo("input2.sourcemap")
}

func (c *CommandLineRunnerTest) TestChunkWrapperBaseNameExpansion() {
	c.useChunks = commandLineRunnerTest_ChunkPattern_CHAIN
	c.argsAdd("--chunk_wrapper=m0:%s // %basename%")
	c.testSame([]string{"var x = 3;", "var y = 4;"})
	builder := jutil.NewStringWriter()
	licenseTracker := jscomp.NewCompiler_ScriptNodeLicensesOnlyTracker(c.lastCompiler)
	chunk := c.lastCompiler.GetChunkGraph().GetRootChunk()
	filename := c.lastCommandLineRunner.GetChunkOutputFileName(chunk)
	c.lastCommandLineRunner.WriteChunkOutput(filename, builder, licenseTracker, chunk)
	testjunit.AssertThat(builder.String()).IsEqualTo(`var x=3; // m0.js
`)
}

func (c *CommandLineRunnerTest) TestChunkWrapperExpansion() {
	c.useChunks = commandLineRunnerTest_ChunkPattern_CHAIN
	c.argsAdd("--chunk_wrapper=m0:%output%%n%//# SourceMappingUrl=%basename%.map")
	c.testSame([]string{"var x = 3;", "var y = 4;"})
	builder := jutil.NewStringWriter()
	licenseTracker := jscomp.NewCompiler_ScriptNodeLicensesOnlyTracker(c.lastCompiler)
	chunk := c.lastCompiler.GetChunkGraph().GetRootChunk()
	filename := c.lastCommandLineRunner.GetChunkOutputFileName(chunk)
	c.lastCommandLineRunner.WriteChunkOutput(filename, builder, licenseTracker, chunk)
	testjunit.AssertThat(builder.String()).IsEqualTo(`var x=3;
//# SourceMappingUrl=m0.js.map
`)
}

func (c *CommandLineRunnerTest) TestMultistageCompilation() {
	saveFile := commandLineRunnerTest_createTempFile("serialized", "state")
	inputString := "[{\"src\": \"alert('foo');\", \"path\":\"foo.js\"}]"
	c.argsAdd("--json_streams=BOTH")
	c.argsAdd("--chunk=foo--bar.baz:1")
	// Perform stage1
	stage1Args := append([]string{}, c.args...)
	stage1Args = append(stage1Args, "--filename_to_save_to="+saveFile.GetAbsolutePath())
	stage1Args = append(stage1Args, "--segment_of_compilation_to_run=CHECKS")
	c.compileString(inputString, stage1Args)
	// Perform stage2
	stage2Args := append([]string{}, c.args...)
	stage2Args = append(stage2Args, "--filename_to_restore_from="+saveFile.GetAbsolutePath())
	stage2Args = append(stage2Args, "--segment_of_compilation_to_run=OPTIMIZATIONS")
	multistageOutput := c.compileString(inputString, stage2Args)
	// Perform single stage compilation
	singleStageOutput := c.compileString(inputString, c.args)
	testjunit.AssertThat(multistageOutput).IsEqualTo(singleStageOutput)
}

func (c *CommandLineRunnerTest) TestCharSetExpansion() {
	c.testSame("")
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetOutputCharset()).IsEqualTo(jutil.StandardCharsets_US_ASCII)
	c.argsAdd("--charset=UTF-8")
	c.testSame("")
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetOutputCharset()).IsEqualTo(jutil.StandardCharsets_UTF_8)
}

func (c *CommandLineRunnerTest) TestChainChunkManifest() {
	c.useChunks = commandLineRunnerTest_ChunkPattern_CHAIN
	c.testSame([]string{"var x = 3;", "var y = 5;", "var z = 7;", "var a = 9;"})
	builder := jutil.NewStringWriter()
	c.lastCommandLineRunner.PrintChunkGraphManifestOrBundleTo(
		c.lastCompiler.GetChunkGraph(), builder, true)
	testjunit.AssertThat(builder.String()).
		IsEqualTo(
			`{m0}
i0.js

{m1:m0}
i1.js

{m2:m1}
i2.js

{m3:m2}
i3.js

{$weak$:m0,m1,m2,m3}
`)
}

func (c *CommandLineRunnerTest) TestStarChunkManifest() {
	c.useChunks = commandLineRunnerTest_ChunkPattern_STAR
	c.testSame([]string{"var x = 3;", "var y = 5;", "var z = 7;", "var a = 9;"})
	builder := jutil.NewStringWriter()
	c.lastCommandLineRunner.PrintChunkGraphManifestOrBundleTo(
		c.lastCompiler.GetChunkGraph(), builder, true)
	testjunit.AssertThat(builder.String()).
		IsEqualTo(
			`{m0}
i0.js

{m1:m0}
i1.js

{m2:m0}
i2.js

{m3:m0}
i3.js

{$weak$:m0,m1,m2,m3}
`)
}

func (c *CommandLineRunnerTest) TestOutputChunkGraphJson() {
	c.useChunks = commandLineRunnerTest_ChunkPattern_STAR
	c.testSame([]string{"var x = 3;", "var y = 5;", "var z = 7;", "var a = 9;"})
	builder := jutil.NewStringWriter()
	c.lastCommandLineRunner.PrintChunkGraphJsonTo(builder)
	testjunit.AssertThat(builder.String()).Contains("transitive-dependencies")
}

func (c *CommandLineRunnerTest) TestVersionFlag_firstArg() {
	c.argsAdd("--version")
	runner := c.createCommandLineRunner([]string{"function f() {}"})
	testjunit.AssertThat(runner.ShouldRunCompiler()).IsTrue()
	testjunit.AssertThat(runner.HasErrors()).IsFalse()
	runner.DoRun()
	testjunit.AssertThat(c.outReader.String()).
		IsEqualTo(runner.GetVersionText() + `
`)
}

func (c *CommandLineRunnerTest) TestVersionFlag_lastArg() {
	c.lastArg = jutil.StringPtr("--version")
	runner := c.createCommandLineRunner([]string{"function f() {}"})
	testjunit.AssertThat(runner.ShouldRunCompiler()).IsTrue()
	testjunit.AssertThat(runner.HasErrors()).IsFalse()
	runner.DoRun()
	testjunit.AssertThat(c.outReader.String()).
		IsEqualTo(runner.GetVersionText() + `
`)
}

func (c *CommandLineRunnerTest) TestPrintAstFlag() {
	c.argsAdd("--print_ast=true")
	c.testSame("")
	testjunit.AssertThat(c.outReader.String()).
		IsEqualTo(
			`digraph AST {
  node [color=lightblue2, style=filled];
  node0 [label="ROOT"];
  node1 [label="SCRIPT"];
  node0 -> node1 [weight=1];
  node1 -> RETURN [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
  node0 -> node1 [label="UNCOND", fontcolor="red", weight=0.01, color="red"];
}

`)
}

func (c *CommandLineRunnerTest) TestSyntheticExterns() {
	c.externs =
		[]*jscomp.SourceFile{jscomp.SourceFile_FromCode("externs", "function Symbol() {}; myVar.property;")}
	c.test(
		"var theirVar = {}; var myVar = {}; var yourVar = {};",
		jscomp.VarCheck_UNDEFINED_EXTERN_VAR_ERROR)
	c.argsAdd("--jscomp_off=externsValidation")
	c.argsAdd("--warning_level=VERBOSE")
	c.test(
		"var theirVar = {}; var myVar = {}; var yourVar = {};",
		"var theirVar={},myVar={},yourVar={};")
	c.argsAdd("--jscomp_off=externsValidation")
	c.argsAdd("--jscomp_error=checkVars")
	c.argsAdd("--warning_level=VERBOSE")
	c.test(
		"var theirVar = {}; var myVar = {}; var myVar = {};", jscomp.VarCheck_VAR_MULTIPLY_DECLARED_ERROR)
}

func (c *CommandLineRunnerTest) TestGoogAssertStripping() {
	c.argsAdd("--compilation_level=ADVANCED_OPTIMIZATIONS")
	c.argsAdd("--jscomp_off=checkVars")
	c.test("goog.asserts.assert(false)", "")
	c.argsAdd("--debug")
	c.test("goog.asserts.assert(false)", "goog.$asserts$.$assert$(!1)")
}

func (c *CommandLineRunnerTest) TestMissingReturnCheckOnWithVerbose() {
	c.argsAdd("--warning_level=VERBOSE")
	c.test(
		"/** @return {number} */ function f() {f()} f();",
		jscomp.CheckMissingReturn_MISSING_RETURN_STATEMENT)
}

func (c *CommandLineRunnerTest) TestChecksOnlySkipsOptimizations() {
	c.argsAdd("--checks_only")
	c.test("var foo = 1 + 1;", "var foo = 1 + 1;")
}

func (c *CommandLineRunnerTest) TestChecksOnlyWithParseError() {
	c.argsAdd("--compilation_level=WHITESPACE_ONLY")
	c.argsAdd("--checks_only")
	c.test("val foo = 1;", jscomp.RhinoErrorReporter_PARSE_ERROR)
}

func (c *CommandLineRunnerTest) TestChecksOnlyWithWarning() {
	c.argsAdd("--checks_only")
	c.argsAdd("--warning_level=VERBOSE")
	c.test("/** @deprecated */function foo() {}; foo();", jscomp.CheckAccessControls_DEPRECATED_NAME)
}

func (c *CommandLineRunnerTest) TestGenerateExports() {
	c.argsAdd("--generate_exports=true")
	c.test(
		"var goog; /** @export */ foo.prototype.x = function() {};",
		`var goog; foo.prototype.x=function(){};
goog.exportProperty(foo.prototype,"x",foo.prototype.x);
`)
}

func (c *CommandLineRunnerTest) TestDepreciationWithVerbose() {
	c.argsAdd("--warning_level=VERBOSE")
	c.test("/** @deprecated */ function f() {}; f()", jscomp.CheckAccessControls_DEPRECATED_NAME)
}

func (c *CommandLineRunnerTest) TestTwoParseErrors() {
	// If parse errors are reported in different files, make
	// sure all of them are reported.
	compiler := c.compile([]string{"var a b;", "var b c;"})
	testjunit.AssertThat(compiler.GetErrors()).HasSize(2)
}

func (c *CommandLineRunnerTest) TestES3() {
	c.argsAdd("--language_in=ECMASCRIPT3")
	c.argsAdd("--language_out=ECMASCRIPT3")
	c.argsAdd("--strict_mode_input=false")
	c.useStringComparison = true
	c.test("var x = f.function", "var x=f[\"function\"];", jscomp.RhinoErrorReporter_INVALID_ES3_PROP_NAME)
	c.testSame("var let;")
}

func (c *CommandLineRunnerTest) TestES3plusStrictModeChecks() {
	c.argsAdd("--language_in=ECMASCRIPT3")
	c.argsAdd("--language_out=ECMASCRIPT3")
	c.useStringComparison = true
	c.test("var x = f.function", "var x=f[\"function\"];", jscomp.RhinoErrorReporter_INVALID_ES3_PROP_NAME)
	c.test("var let", jscomp.RhinoErrorReporter_PARSE_ERROR)
	c.test("function f(x) { delete x; }", jscomp.StrictModeCheck_DELETE_VARIABLE)
}

func (c *CommandLineRunnerTest) TestES6NotTranspiledByDefault() {
	c.testSame("var x = class {};")
	c.argsAdd("--language_out=STABLE")
	c.test("var x = class {};", "var x = function() {};")
}

func (c *CommandLineRunnerTest) TestES5ChecksByDefault() {
	c.test("var x = 3; delete x;", jscomp.StrictModeCheck_DELETE_VARIABLE)
}

func (c *CommandLineRunnerTest) TestES5ChecksInVerbose() {
	c.argsAdd("--warning_level=VERBOSE")
	c.test("function f(x) { delete x; }", jscomp.StrictModeCheck_DELETE_VARIABLE)
}

func (c *CommandLineRunnerTest) TestES5() {
	c.argsAdd("--language_in=ECMASCRIPT5")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.argsAdd("--strict_mode_input=false")
	c.test("var x = f.function", "var x = f.function")
	c.test("var let", "var let")
}

func (c *CommandLineRunnerTest) TestES5Strict() {
	c.argsAdd("--language_in=ECMASCRIPT5_STRICT")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.argsAdd("--emit_use_strict=true")
	c.test("var x = f.function", "'use strict';var x = f.function")
	c.test("var let", jscomp.RhinoErrorReporter_PARSE_ERROR)
	c.test("function f(x) { delete x; }", jscomp.StrictModeCheck_DELETE_VARIABLE)
}

func (c *CommandLineRunnerTest) TestES5StrictUseStrict() {
	c.argsAdd("--language_in=ECMASCRIPT5_STRICT")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.argsAdd("--emit_use_strict=true")
	compiler := c.compile([]string{"var x = f.function"})
	outputSource := compiler.ToSource()
	testjunit.AssertThat(outputSource).StartsWith("'use strict'")
}

func (c *CommandLineRunnerTest) TestES5StrictUseStrictMultipleInputs() {
	c.argsAdd("--language_in=ECMASCRIPT5_STRICT")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.argsAdd("--emit_use_strict=true")
	compiler := c.compile([]string{"var x = f.function", "var y = f.function", "var z = f.function"})
	outputSource := compiler.ToSource()
	testjunit.AssertThat(outputSource).StartsWith("'use strict'")
	testjunit.AssertThat(jutil.StringSubstringFrom(outputSource, 13)).DoesNotContain("'use strict'")
}

func (c *CommandLineRunnerTest) TestEs5OutputRunsClassCodeInSloppyMode() {
	// Won't fix (b/573608886): class code is strict mode code, so natively `this.ro = 2` throws a
	// TypeError, because `ro` only has a getter. With --language_out=ECMASCRIPT5, unlike
	// ECMASCRIPT5_STRICT, the output has no 'use strict', so the transpiled class code runs in
	// sloppy mode and the write is silently ignored. We decided not to change this, since the
	// compiler isn't required to preserve errors in code that is already broken.
	c.argsAdd("--language_in=ECMASCRIPT_2015")
	c.argsAdd("--language_out=ECMASCRIPT5")
	compiler :=
		c.compile(
			[]string{
				"class C { get ro() { return 1; } write() { this.ro = 2; } } new C().write();",
			})
	outputSource := compiler.ToSource()
	testjunit.AssertThat(outputSource).DoesNotContain("use strict")
	testjunit.AssertThat(outputSource).Contains("this.ro=2")
}

func (c *CommandLineRunnerTest) TestWithKeywordWithEs5ChecksOff() {
	c.argsAdd("--jscomp_off=es5Strict")
	c.argsAdd("--strict_mode_input=false")
	c.testSame("var x = {}; with (x) {}")
}

func (c *CommandLineRunnerTest) TestIsolationMode() {
	c.argsAdd("--isolation_mode=IIFE")
	c.testSame("window.x = \"123\";")
	output := c.outReader.String()
	testjunit.AssertThat(output).IsEqualTo(`(function(){window.x="123";}).call(this);
`)
}

func (c *CommandLineRunnerTest) TestNoSrCFilesWithManifest() {
	c.argsAdd("--env=CUSTOM")
	c.argsAdd("--output_manifest=test.MF")
	runner := c.createCommandLineRunner([]string{})
	jutil.Try(func() {
		runner.DoRun()
		testjunit.AssertWithMessage("Expected flag usage exception").Fail()
	}, func(e jscomp.IFlagUsageException) {
		testjunit.AssertThat(e.GetMessage()).
			IsEqualTo(
				"Bad --js flag. Manifest files cannot be generated when the input is from stdin.")
	})
}

func (c *CommandLineRunnerTest) TestProcessCJS() {
	c.useStringComparison = true
	c.argsAdd("--process_common_js_modules")
	c.argsAdd("--entry_point=foo/bar")
	c.setFilename(0, "foo/bar.js")
	expected := "var module$foo$bar={default:{}};module$foo$bar.default.test=1;"
	c.test("exports.test = 1", expected)
	testjunit.AssertThat(c.outReader.String()).IsEqualTo(expected + `
`)
}

func (c *CommandLineRunnerTest) TestProcessCJSWithChunkOutput() {
	c.argsAdd("--process_common_js_modules")
	c.argsAdd("--entry_point=foo/bar")
	c.argsAdd("--chunk=auto")
	c.setFilename(0, "foo/bar.js")
	c.test("exports.test = 1", "var module$foo$bar={default: {}}; module$foo$bar.default.test = 1;")
	// With chunk=auto no direct output is created.
	testjunit.AssertThat(c.outReader.String()).IsEmpty()
}

func (c *CommandLineRunnerTest) TestSimpleJsonFileInclusionCompiles() {
	c.argsAdd("--module_resolution=NODE")
	c.argsAdd("--compilation_level=ADVANCED")
	c.setFilename(0, "index.js")
	c.setFilename(1, "package.json")
	c.test(
		[]string{
			"const /** string */ typeError = 0;",
			`{
  "name": "test"
}
`,
		},
		jscomp.TypeValidator_TYPE_MISMATCH_WARNING)
}

/**
 * closure requires mixed with cjs, raised in https://github.com/google/closure-compiler/pull/630
 * https://gist.github.com/sayrer/c4c4ce0c1748573f863e
 */
func (c *CommandLineRunnerTest) TestProcessCJSWithClosureRequires() {
	c.argsAdd("--process_common_js_modules")
	c.argsAdd("--entry_point=app")
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--module_resolution=NODE")
	c.setFilename(0, "base.js")
	c.setFilename(1, "array.js")
	c.setFilename(2, "Baz.js")
	c.setFilename(3, "app.js")
	c.test(
		[]string{
			`/** @provideGoog */
/** @const */ var goog = goog || {};
var COMPILED = false;
goog.provide = function (arg) {};
goog.require = function (arg) {};
`,
			"goog.provide('goog.array');",
			`goog.require('goog.array');
function Baz() {}
Baz.prototype = {
  baz: function() {
    return goog.array.last(['asdf','asd','baz']);
  },
  bar: function () {
    return 4 + 4;
  }
};
module.exports = Baz;
`,
			`var Baz = require('./Baz');
var baz = new Baz();
console.log(baz.baz());
console.log(baz.bar());
`,
		},
		[]string{
			`var goog=goog||{},COMPILED=!1;
goog.provide=function(a){};goog.require=function(a){};
`,
			"goog.array={};",
			`var module$Baz = {/** @constructor */ default: function (){} };
module$Baz.default.prototype={
  baz:function(){return goog.array.last(['asdf','asd','baz'])},
  bar:function(){return 8}
};
`,
			`var Baz = module$Baz.default,
    baz = new module$Baz.default();
console.log(baz.baz());
console.log(baz.bar());
`,
		})
}

/**
 * closure requires mixed with cjs, raised in https://github.com/google/closure-compiler/pull/630
 * https://gist.github.com/sayrer/c4c4ce0c1748573f863e
 */
func (c *CommandLineRunnerTest) TestProcessCJSWithClosureRequires2() {
	c.argsAdd("--process_common_js_modules")
	c.argsAdd("--dependency_mode=PRUNE_LEGACY")
	c.argsAdd("--entry_point=app")
	c.argsAdd("--module_resolution=NODE")
	c.setFilename(0, "base.js")
	c.setFilename(1, "array.js")
	c.setFilename(2, "Baz.js")
	c.setFilename(3, "app.js")
	c.test(
		[]string{
			`/** @provideGoog */
/** @const */ var goog = goog || {};
var COMPILED = false;
goog.provide = function (arg) {};
goog.require = function (arg) {};
`,
			"goog.provide('goog.array');",
			`goog.require('goog.array');
function Baz() {}
Baz.prototype = {
  baz: function() {
    return goog.array.last(['asdf','asd','baz']);
  },
  bar: function () {
    return 4 + 4;
  }
};
module.exports = Baz;
`,
			`var Baz = require('./Baz');
var baz = new Baz();
console.log(baz.baz());
console.log(baz.bar());
`,
		},
		[]string{
			`var goog=goog||{},COMPILED=!1;
goog.provide=function(a){};goog.require=function(a){};
`,
			"goog.array={};",
			`var module$Baz = {default: function (){}};
module$Baz.default.prototype={
  baz:function(){return goog.array.last(["asdf","asd","baz"])},
  bar:function(){return 8}
};
`,
			`var Baz = module$Baz.default,
    baz = new module$Baz.default();
console.log(baz.baz());
console.log(baz.bar());
`,
		})
}

func (c *CommandLineRunnerTest) TestProcessCJSWithES6Export() {
	c.argsAdd("--process_common_js_modules")
	c.argsAdd("--entry_point=app")
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.argsAdd("--module_resolution=NODE")
	c.setFilename(0, "foo.js")
	c.setFilename(1, "app.js")
	c.test(
		[]string{
			`export default class Foo {
  bar() { console.log('bar'); }
}
`,
			`var FooBar = require('./foo').default;
var baz = new FooBar();
console.log(baz.bar());
`,
		},
		[]string{
			`var Foo$$module$foo=function(){};
Foo$$module$foo.prototype.bar=function(){console.log("bar")};
var module$foo={};
/** @const */ module$foo.default=Foo$$module$foo;
`,
			`var FooBar = Foo$$module$foo,
    baz = new Foo$$module$foo();
console.log(baz.bar());
`,
		})
}

func (c *CommandLineRunnerTest) TestES6ImportOfCJS() {
	c.argsAdd("--process_common_js_modules")
	c.argsAdd("--entry_point=app")
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.argsAdd("--module_resolution=NODE")
	c.setFilename(0, "foo.js")
	c.setFilename(1, "app.js")
	c.test(
		[]string{
			`/** @constructor */ function Foo () {}
Foo.prototype.bar = function() { console.log('bar'); };
module.exports = Foo;
`,
			`import * as FooBar from './foo';
var baz = new FooBar();
console.log(baz.bar());
`,
		},
		[]string{
			`/** @const */ var module$foo = {/** @constructor */ default: function(){} };
module$foo.default.prototype.bar=function(){console.log('bar')};
`,
			`var baz$$module$app = new module$foo();
console.log(baz$$module$app.bar());
/** @const */ var module$app = {};
`,
		})
}

func (c *CommandLineRunnerTest) TestES6ImportOfFileWithoutImportsOrExports() {
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--entry_point='./app.js'")
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.setFilename(0, "foo.js")
	c.setFilename(1, "app.js")
	c.test(
		[]string{
			`function foo() { alert('foo'); }
foo();
`,
			"import './foo.js';",
		},
		[]string{
			`function foo$$module$foo(){ alert('foo'); }
foo$$module$foo();
/** @const */ var module$foo = {}
`,
			"/** @const */ var module$app = {};",
		})
}

func (c *CommandLineRunnerTest) TestES6ImportOfFileWithImportsButNoExports() {
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--entry_point='./app.js'")
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.setFilename(0, "message.js")
	c.setFilename(1, "foo.js")
	c.setFilename(2, "app.js")
	c.test(
		[]string{
			"export default 'message';",
			`import message from './message.js';
function foo() { alert(message); }
foo();`,
			"import './foo.js';",
		},
		[]string{
			`var $jscompDefaultExport$$module$message = 'message', module$message = {};
/** @const */ module$message.default = $jscompDefaultExport$$module$message;
`,
			`function foo$$module$foo(){
  alert($jscompDefaultExport$$module$message);
}
foo$$module$foo();
/** @const */ var module$foo = {};
`,
			"/** @const */ var module$app = {};",
		})
}

func (c *CommandLineRunnerTest) TestCommonJSRequireOfFileWithoutExports() {
	c.argsAdd("--process_common_js_modules")
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--entry_point='./app.js'")
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--module_resolution=NODE")
	c.setFilename(0, "foo.js")
	c.setFilename(1, "app.js")
	c.test(
		[]string{
			`function foo() { alert('foo'); }
foo();
`,
			"require('./foo');",
		},
		[]string{
			`/** @const */ var module$foo = {/** @const */ default: {}};
function foo$$module$foo(){ alert('foo'); }
foo$$module$foo();
`,
			`'use strict';

`,
		})
}

/** override the order of the entries that the module loader should look for */
func (c *CommandLineRunnerTest) TestProcessCJSWithPackageJsonBrowserField() {
	c.useStringComparison = true
	c.argsAdd("--process_common_js_modules")
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--entry_point=app")
	c.argsAdd("--module_resolution=NODE")
	c.argsAdd("--package_json_entry_names=browser,main")
	c.setFilename(0, "app.js")
	c.setFilename(1, "node_modules/foo/package.json")
	c.setFilename(2, "node_modules/foo/browser.js")
	c.test(
		[]string{
			"var Foo = require('foo');",
			"{\"browser\":\"browser.js\",\"name\":\"foo\"}",
			`function Foo() {}
Foo.prototype = {
  bar: function () {
    return 4 + 4;
  }
};
module.exports = Foo;
`,
		},
		[]string{
			"var module$node_modules$foo$browser={default:function(){}};",
			"module$node_modules$foo$browser.default.prototype={bar:function(){return 8}};",
			"var Foo=module$node_modules$foo$browser.default;",
		})
}

func (c *CommandLineRunnerTest) TestFormattingSingleQuote() {
	c.testSame("var x = '';")
	testjunit.AssertThat(c.lastCompiler.ToSource()).IsEqualTo("var x=\"\";")
	c.argsAdd("--formatting=SINGLE_QUOTES")
	c.testSame("var x = '';")
	testjunit.AssertThat(c.lastCompiler.ToSource()).IsEqualTo("var x='';")
}

func (c *CommandLineRunnerTest) TestChunkJSON() {
	c.argsAdd("--process_common_js_modules")
	c.argsAdd("--entry_point=foo/bar")
	c.argsAdd("--output_chunk_dependencies=test.json")
	c.setFilename(0, "foo/bar.js")
	c.test(
		"module.exports = {foo: 1};",
		`/** @const */ var module$foo$bar = {/** @const */ default: {}};
module$foo$bar.default.foo = 1;
`)
}

func (c *CommandLineRunnerTest) TestOutputSameAsInput() {
	c.argsAdd("--js_output_file=" + c.getFilename(0))
	c.test("", AbstractCommandLineRunner_OUTPUT_SAME_AS_INPUT_ERROR)
}

func (c *CommandLineRunnerTest) TestOutputWrapperFlag() {
	// if the output wrapper flag is specified without a valid output marker,
	// ensure that the compiler displays an error and exits.
	// See github issue 123
	c.argsAdd("--output_wrapper=output")
	runner := c.createCommandLineRunner([]string{"function f() {}"})
	testjunit.AssertThat(runner.ShouldRunCompiler()).IsFalse()
	testjunit.AssertThat(runner.HasErrors()).IsTrue()
}

func (c *CommandLineRunnerTest) TestJsonStreamInputFlag() {
	inputString := "[{\"src\": \"alert('foo');\", \"path\":\"foo.js\"}]"
	c.argsAdd("--json_streams=IN")
	runner := NewCommandLineRunnerWithStreams(
		append([]string{}, c.args...),
		bytes.NewReader(jutil.StringGetBytesUTF8(inputString)),
		jutil.NewPrintStream(c.outReader),
		jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		jutil.Throwable_PrintStackTrace(e, jutil.System_Err())
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})
	output := runner.GetCompiler().ToSource()
	testjunit.AssertThat(output).IsEqualTo("alert(\"foo\");")
}

func (c *CommandLineRunnerTest) TestJsonStreamOutputFlag() {
	inputString := "alert('foo');"
	c.argsAdd("--json_streams=OUT")
	runner := NewCommandLineRunnerWithStreams(
		append([]string{}, c.args...),
		bytes.NewReader(jutil.StringGetBytesUTF8(inputString)),
		jutil.NewPrintStream(c.outReader),
		jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		jutil.Throwable_PrintStackTrace(e, jutil.System_Err())
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})
	output := c.outReader.String()
	testjunit.AssertThat(output).
		IsEqualTo(
			"[{\"src\":\"alert(\\\"foo\\\");\\n\",\"path\":\"compiled.js\",\"source_map\":\"{\\n\\\"version\\\":3,\\n\\\"file\\\":\\\"compiled.js\\\",\\n\\\"lineCount\\\":1,\\n\\\"mappings\\\":\\\"AAAAA,KAAA,CAAM,KAAN;\\\",\\n\\\"sources\\\":[\\\"stdin\\\"],\\n\\\"names\\\":[\\\"alert\\\"]\\n}\\n\"}]")
}

func (c *CommandLineRunnerTest) TestJsonStreamBothFlag() {
	inputString := "[{\"src\": \"alert('foo');\", \"path\":\"foo.js\"}]"
	c.argsAdd("--json_streams=BOTH")
	c.argsAdd("--js_output_file=bar.js")
	runner := NewCommandLineRunnerWithStreams(
		append([]string{}, c.args...),
		bytes.NewReader(jutil.StringGetBytesUTF8(inputString)),
		jutil.NewPrintStream(c.outReader),
		jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		jutil.Throwable_PrintStackTrace(e, jutil.System_Err())
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})
	output := c.outReader.String()
	testjunit.AssertThat(output).
		IsEqualTo(
			"[{\"src\":\"alert(\\\"foo\\\");\\n\",\"path\":\"bar.js\",\"source_map\":\"{\\n\\\"version\\\":3,\\n\\\"file\\\":\\\"bar.js\\\",\\n\\\"lineCount\\\":1,\\n\\\"mappings\\\":\\\"AAAAA,KAAA,CAAM,KAAN;\\\",\\n\\\"sources\\\":[\\\"foo.js\\\"],\\n\\\"names\\\":[\\\"alert\\\"]\\n}\\n\"}]")
}

func (c *CommandLineRunnerTest) TestJsonStreamSourceMap() {
	inputSourceMap := `{
"version":3,
"file":"one.out.js",
"lineCount":1,
"mappings":"AAAAA,QAASA,IAAG,CAACC,CAAD,CAAI,CACdC,OAAAF,IAAA,CAAYC,CAAZ,CADc,CAGhBD,GAAA,CAAI,QAAJ;",
"sources":["one.js"],
"names":["log","a","console"]
}`
	inputSourceMap = strings.ReplaceAll(inputSourceMap, "\"", "\\\"")
	inputString := strings.ReplaceAll("[{\"src\": \"function log(a){console.log(a)}log(\\\"one.js\\\");\", \"path\":\"one.out.js\", \"sourceMap\": \"INPUT_SOURCE_MAP\" }]", "INPUT_SOURCE_MAP", inputSourceMap)
	c.argsAdd("--json_streams=BOTH")
	c.argsAdd("--js_output_file=bar.js")
	c.argsAdd("--apply_input_source_maps")
	runner := NewCommandLineRunnerWithStreams(
		append([]string{}, c.args...),
		bytes.NewReader(jutil.StringGetBytesUTF8(inputString)),
		jutil.NewPrintStream(c.outReader),
		jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		jutil.Throwable_PrintStackTrace(e, jutil.System_Err())
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})
	output := c.outReader.String()
	testjunit.AssertThat(output).
		IsEqualTo(
			"[{\"src\":\"function log(a){console.log(a)}log(\\\"one.js\\\");\\n\",\"path\":\"bar.js\",\"source_map\":\"{\\n\\\"version\\\":3,\\n\\\"file\\\":\\\"bar.js\\\",\\n\\\"lineCount\\\":1,\\n\\\"mappings\\\":\\\"AAAAA,QAASA,IAAG,CAACC,CAAD,CAAI,CACdC,OAAAF,CAAAA,GAAAE,CAAYD,CAAZC,CADc,CAGhBF,GAAAA,CAAI,QAAJA;\\\",\\n\\\"sources\\\":[\\\"one.js\\\"],\\n\\\"names\\\":[\\\"log\\\",\\\"a\\\",\\\"console\\\"]\\n}\\n\"}]")
}

func (c *CommandLineRunnerTest) TestJsonStreamSourceMapUnderscore() {
	inputSourceMap := `{
"version":3,
"file":"one.out.js",
"lineCount":1,
"mappings":"AAAAA,QAASA,IAAG,CAACC,CAAD,CAAI,CACdC,OAAAF,IAAA,CAAYC,CAAZ,CADc,CAGhBD,GAAA,CAAI,QAAJ;",
"sources":["one.js"],
"names":["log","a","console"]
}`
	inputSourceMap = strings.ReplaceAll(inputSourceMap, "\"", "\\\"")
	inputString := // input JSON with `source_map` field instead of `sourceMap`
		strings.ReplaceAll("[{\"src\": \"function log(a){console.log(a)}log(\\\"one.js\\\");\", \"path\":\"one.out.js\", \"source_map\": \"INPUT_SOURCE_MAP\" }]", "INPUT_SOURCE_MAP", inputSourceMap)
	c.argsAdd("--json_streams=BOTH")
	c.argsAdd("--js_output_file=bar.js")
	c.argsAdd("--apply_input_source_maps")
	runner := NewCommandLineRunnerWithStreams(
		append([]string{}, c.args...),
		bytes.NewReader(jutil.StringGetBytesUTF8(inputString)),
		jutil.NewPrintStream(c.outReader),
		jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		jutil.Throwable_PrintStackTrace(e, jutil.System_Err())
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})
	output := c.outReader.String()
	testjunit.AssertThat(output).
		IsEqualTo(
			"[{\"src\":\"function log(a){console.log(a)}log(\\\"one.js\\\");\\n\",\"path\":\"bar.js\",\"source_map\":\"{\\n\\\"version\\\":3,\\n\\\"file\\\":\\\"bar.js\\\",\\n\\\"lineCount\\\":1,\\n\\\"mappings\\\":\\\"AAAAA,QAASA,IAAG,CAACC,CAAD,CAAI,CACdC,OAAAF,CAAAA,GAAAE,CAAYD,CAAZC,CADc,CAGhBF,GAAAA,CAAI,QAAJA;\\\",\\n\\\"sources\\\":[\\\"one.js\\\"],\\n\\\"names\\\":[\\\"log\\\",\\\"a\\\",\\\"console\\\"]\\n}\\n\"}]")
}

func (c *CommandLineRunnerTest) TestJsonStreamAllowsAnyChunkName() {
	inputString := "[{\"src\": \"alert('foo');\", \"path\":\"foo.js\"}]"
	c.argsAdd("--json_streams=BOTH")
	c.argsAdd("--chunk=foo/bar/baz:1")
	runner := NewCommandLineRunnerWithStreams(
		append([]string{}, c.args...),
		bytes.NewReader(jutil.StringGetBytesUTF8(inputString)),
		jutil.NewPrintStream(c.outReader),
		jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		jutil.Throwable_PrintStackTrace(e, jutil.System_Err())
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})
	output := c.outReader.String()
	testjunit.AssertThat(output).
		IsEqualTo(
			"[{\"src\":\"alert(\\\"foo\\\");\\n\",\"path\":\"./foo/bar/baz.js\",\"source_map\":\"{\\n\\\"version\\\":3,\\n\\\"file\\\":\\\"./foo/bar/baz.js\\\",\\n\\\"lineCount\\\":1,\\n\\\"mappings\\\":\\\"AAAAA,KAAA,CAAM,KAAN;\\\",\\n\\\"sources\\\":[\\\"foo.js\\\"],\\n\\\"names\\\":[\\\"alert\\\"]\\n}\\n\"}]")
}

func (c *CommandLineRunnerTest) TestOutputModuleNaming() {
	inputString := "[{\"src\": \"alert('foo');\", \"path\":\"foo.js\"}]"
	c.argsAdd("--json_streams=BOTH")
	c.argsAdd("--chunk=foo--bar.baz:1")
	runner := NewCommandLineRunnerWithStreams(
		append([]string{}, c.args...),
		bytes.NewReader(jutil.StringGetBytesUTF8(inputString)),
		jutil.NewPrintStream(c.outReader),
		jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		jutil.Throwable_PrintStackTrace(e, jutil.System_Err())
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})
	output := c.outReader.String()
	testjunit.AssertThat(output).
		IsEqualTo(
			"[{\"src\":\"alert(\\\"foo\\\");\\n\",\"path\":\"./foo--bar.baz.js\",\"source_map\":\"{\\n\\\"version\\\":3,\\n\\\"file\\\":\\\"./foo--bar.baz.js\\\",\\n\\\"lineCount\\\":1,\\n\\\"mappings\\\":\\\"AAAAA,KAAA,CAAM,KAAN;\\\",\\n\\\"sources\\\":[\\\"foo.js\\\"],\\n\\\"names\\\":[\\\"alert\\\"]\\n}\\n\"}]")
}

func (c *CommandLineRunnerTest) TestAssumeFunctionWrapper() {
	c.argsAdd("--compilation_level=SIMPLE_OPTIMIZATIONS")
	c.argsAdd("--assume_function_wrapper")
	// remove used vars enabled
	c.test("var someName = function(a) {};", "")
	// function inlining enabled.
	c.test("var someName = function() {return 'hi'};alert(someName())", "alert('hi')")
	// renaming enabled, and collapse anonymous functions enabled
	c.test(
		"var someName = function() {return 'hi'};alert(someName);alert(someName)",
		"function a() {return 'hi'}alert(a);alert(a)")
}

func (c *CommandLineRunnerTest) TestWebpackModuleIds() {
	inputString := `[
  {"src": "__webpack_require__(2);", "path":"foo.js", "webpackId": "1"},
  {"src": "console.log('bar');", "path":"bar.js", "webpackId": "2"}
]
`
	c.argsAdd("--json_streams=BOTH")
	c.argsAdd("--module_resolution=WEBPACK")
	c.argsAdd("--process_common_js_modules")
	c.argsAdd("--entry_point=foo.js")
	c.argsAdd("--dependency_mode=PRUNE")
	c.argsAdd("--js_output_file=out.js")
	runner := NewCommandLineRunnerWithStreams(
		append([]string{}, c.args...),
		bytes.NewReader(jutil.StringGetBytesUTF8(inputString)),
		jutil.NewPrintStream(c.outReader),
		jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	runner.DoRun()
	output := c.outReader.String()
	testjunit.AssertThat(output).
		IsEqualTo(
			"[{\"src\":\"var module$bar={default:{}};console.log(\\\"bar\\\");var module$foo={default:{}};\\n\",\"path\":\"out.js\",\"source_map\":\"{\\n\\\"version\\\":3,\\n\\\"file\\\":\\\"out.js\\\",\\n\\\"lineCount\\\":1,\\n\\\"mappings\\\":\\\"AAAA,IAAA,WAAA,CAAA,QAAA,EAAA,CAAAA,QAAQC,CAAAA,GAAR,CAAY,KAAZ,C,CCAA,IAAA,WAAA,CAAA,QAAA,EAAA;\\\",\\n\\\"sources\\\":[\\\"bar.js\\\",\\\"foo.js\\\"],\\n\\\"names\\\":[\\\"console\\\",\\\"log\\\"]\\n}\\n\"}]")
}

func (c *CommandLineRunnerTest) TestInstrumentCodeByLine() {
	c.argsAdd("--instrument_for_coverage_option=LINE")
	c.test(
		"function foo(){ const answerToAll = 42; }",
		`(function(a){a.window||(a.window=a,a.window.top=a)})(typeof self!=="undefined"?self:globalThis);var __jscov=window.top.__jscov||
(window.top.__jscov={fileNames:[],instrumentedLines:[],executedLines:[]}),
JSCompiler_lcov_data_input0=[];
__jscov.executedLines.push(JSCompiler_lcov_data_input0);
__jscov.instrumentedLines.push("01");
__jscov.fileNames.push("input0");
function foo(){JSCompiler_lcov_data_input0[0]=!0}
`)
}

func (c *CommandLineRunnerTest) TestInstrumentCodeByBranch() {
	c.argsAdd("--instrument_for_coverage_option=BRANCH")
	c.test(
		"function foo(){ const answerToAll = 42; }",
		`(function(a){a.window||(a.window=a,a.window.top=a)})(typeof self!=="undefined"?self:globalThis);
var __jscov=window.top.__jscov||(window.top.__jscov=
{fileNames:[],branchPresent:[],branchesInLine:[],branchesTaken:[]});
function foo(){}
`)
}

func (c *CommandLineRunnerTest) TestInstrumentCodeMappingFileNotSet() {
	c.argsAdd("--instrument_for_coverage_option=PRODUCTION")
	e := testjunit.AssertThrows[*jscomp.FlagUsageException](func() { c.compileString("", c.args) })
	testjunit.AssertThat(e.GetMessage()).
		IsEqualTo(
			"Expected --instrument_mapping_report to be set when --instrument_for_coverage_option is set to Production")
}

func (c *CommandLineRunnerTest) TestInstrumentCodeProductionArgNotSet() {
	c.argsAdd("--instrument_mapping_report=someFile.txt")
	e := testjunit.AssertThrows[*jscomp.FlagUsageException](func() { c.compileString("", c.args) })
	testjunit.AssertThat(e.GetMessage()).
		IsEqualTo(
			"Expected --instrument_for_coverage_option to be passed with PRODUCTION when --instrument_mapping_report is set")
}

func (c *CommandLineRunnerTest) TestInstrumentCodeProductionArrayArgNotSet() {
	c.argsAdd("--instrument_for_coverage_option=PRODUCTION")
	c.argsAdd("--instrument_mapping_report=someFile.txt")
	e := testjunit.AssertThrows[*jscomp.FlagUsageException](func() { c.compileString("", c.args) })
	testjunit.AssertThat(e.GetMessage()).
		IsEqualTo(
			"Expected --production_instrumentation_array_name to be set when --instrument_for_coverage_option is set to Production")
}

func (c *CommandLineRunnerTest) TestBundleOutput_bundlesGoogModule() {
	// Create a path for the final output
	bundledFile := c.temporaryFolder.NewFile("bundle.js")
	jsFile1 := c.createJsFile("test1", "var a;")
	jsFile2 := c.createJsFile("test2", "goog.module('foo'); var b;")
	c.argsAdd("--compilation_level=BUNDLE")
	c.argsAdd("--dependency_mode=NONE")
	c.argsAdd("--js_output_file")
	c.argsAdd(bundledFile.GetPath())
	c.compileJsFiles("", jsFile1, jsFile2)
	bundledJs := jutil.Files_ReadStringUTF8(bundledFile.ToPath())
	expected := strings.ReplaceAll(strings.ReplaceAll(`//VALUE_1
var a;
//VALUE_2
goog.loadModule(function(exports) {'use strict';goog.module('foo'); var b;
;return exports;});

`, "VALUE_1", jsFile1.GetValue()), "VALUE_2", jsFile2.GetValue())
	testjunit.AssertThat(bundledJs).IsEqualTo(expected)
}

func (c *CommandLineRunnerTest) TestBundleOutput_ignoresSyntaxErrors() {
	// Verify that if bundling, the compiler doesn't run a full parse and thus doesn't report
	// syntax errors.

	// Create a path for the final output
	bundledFile := c.temporaryFolder.NewFile("bundle.js")
	jsFile := c.createJsFile("test1", "var a; syntax error!")
	c.argsAdd("--compilation_level=BUNDLE")
	c.argsAdd("--dependency_mode=NONE")
	c.argsAdd("--js_output_file")
	c.argsAdd(bundledFile.GetPath())
	c.compileJsFiles("", jsFile)
	bundledJs := jutil.Files_ReadStringUTF8(bundledFile.ToPath())
	expected := strings.ReplaceAll(`//VALUE
var a; syntax error!
`, "VALUE", jsFile.GetValue())
	testjunit.AssertThat(bundledJs).IsEqualTo(expected)
}

// @Rule public TemporaryFolder folder = new TemporaryFolder() is
// c.temporaryFolder here.
func (c *CommandLineRunnerTest) TestInstrumentCodeProductionCreatesInstrumentationMapping() {
	tempFolderPath := c.temporaryFolder.getRoot()
	filePath := tempFolderPath + "\\someFile.txt"
	c.argsAdd("--instrument_for_coverage_option=PRODUCTION")
	c.argsAdd("--instrument_mapping_report=" + filePath)
	c.argsAdd("--production_instrumentation_array_name=ist_arr")
	c.argsAdd("--language_out=NO_TRANSPILE")
	c.argsAdd("--formatting=PRETTY_PRINT")
	source := `function foo() {
   console.log('Hello');
}
`
	expected := `var ist_arr = [];
function foo() {
  ist_arr.push('C');
  console.log('Hello');
}
`
	c.externs =
		[]*jscomp.SourceFile{
			NewTestExternsBuilder().
				AddArray().
				AddAlert().
				AddExtra("let ist_arr;").
				BuildExternsFile("externs")}
	c.test(source, expected)
	variableMap := jutil.NewFile(filePath)
	variableMapFile := commandLineRunnerTest_readLines(variableMap)
	testjunit.AssertThat(variableMapFile).
		ContainsExactly(
			" FileNames:[\"input0\"]",
			" FunctionNames:[\"foo\"]",
			" Types:[\"FUNCTION\"]",
			"C:AAACA").
		InOrder()
}

func (c *CommandLineRunnerTest) TestEscapeDollarInTemplateLiteralInOutput() {
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--language_out=ECMASCRIPT6")
	c.test("let Foo; const x = `${Foo}`;", "let Foo; const x = `${Foo}`;")
	c.test("const x = `\\${Foo}`;", "const x = '\\${Foo}'")
	c.test("let Foo; const x = `${Foo}\\${Foo}`;", "let Foo; const x = `${Foo}\\${Foo}`;")
	c.test("let Foo; const x = `\\${Foo}${Foo}`;", "let Foo; const x = `\\${Foo}${Foo}`;")
}

func (c *CommandLineRunnerTest) TestEscapeDollarInTemplateLiteralEs5Output() {
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.test(
		"let Foo; const x = `${Foo}`;",
		`var $jscomp=$jscomp||{};$jscomp.scope={};
$jscomp.createTemplateTagFirstArg=function(a){
  return $jscomp.createTemplateTagFirstArgWithRaw(a,a)
};
$jscomp.createTemplateTagFirstArgWithRaw=function(a,b){
  a.raw=b;
  Object.freeze && (Object.freeze(a), Object.freeze(b));
  return a
}
var Foo,x=""+Foo
`)
	c.test(
		"const x = `\\${Foo}`;",
		`var $jscomp=$jscomp||{};$jscomp.scope={};
$jscomp.createTemplateTagFirstArg=function(a){
  return $jscomp.createTemplateTagFirstArgWithRaw(a,a)
};
$jscomp.createTemplateTagFirstArgWithRaw=function(a,b){
  a.raw=b;
  Object.freeze && (Object.freeze(a), Object.freeze(b));
  return a
}
var x="${Foo}"
`)
	c.test(
		"let Foo; const x = `${Foo}\\${Foo}`;",
		`var $jscomp=$jscomp||{};$jscomp.scope={};
$jscomp.createTemplateTagFirstArg=function(a){
  return $jscomp.createTemplateTagFirstArgWithRaw(a,a)
};
$jscomp.createTemplateTagFirstArgWithRaw=function(a,b){
  a.raw=b;
  Object.freeze && (Object.freeze(a), Object.freeze(b));
  return a
}
var Foo,x=Foo+"${Foo}"
`)
	c.test(
		"let Foo; const x = `\\${Foo}${Foo}`;",
		`var $jscomp=$jscomp||{};$jscomp.scope={};
$jscomp.createTemplateTagFirstArg=function(a){
  return $jscomp.createTemplateTagFirstArgWithRaw(a,a)
};
$jscomp.createTemplateTagFirstArgWithRaw=function(a,b){
  a.raw=b;
  Object.freeze && (Object.freeze(a), Object.freeze(b));
  return a
}
var Foo,x="${Foo}"+Foo
`)
}

/** windows shells can add extra quotes to an argument */
func (c *CommandLineRunnerTest) TestWarningGuardQuotedValue() {
	c.argsAdd("--jscomp_error='\"*\"'")
	c.argsAdd("--jscomp_warning=\"'*'\"")
	c.argsAdd("--jscomp_off='\"*\"'")
	c.testSame("alert('hello world')")
}

func (c *CommandLineRunnerTest) TestOutputSourceMapContainsSourcesContent() {
	inputString := "[{\"src\": \"\", \"path\": \"base.js\"}," +
		"{\"src\": \"alert('foo');\", \"path\":\"foo.js\"}]"
	c.argsAdd("--json_streams=BOTH")
	c.argsAdd("--chunk=m1:1")
	c.argsAdd("--chunk=m2:1:m1")
	c.argsAdd("--source_map_include_content")
	runner := NewCommandLineRunnerWithStreams(
		append([]string{}, c.args...),
		bytes.NewReader(jutil.StringGetBytesUTF8(inputString)),
		jutil.NewPrintStream(c.outReader),
		jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		jutil.Throwable_PrintStackTrace(e, jutil.System_Err())
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})
	output := c.outReader.String()
	testjunit.AssertThat(output).
		IsEqualTo(
			"[{\"src\":\"\\n\",\"path\":\"./m1.js\",\"source_map\":\"{\\n\\\"version\\\":3,\\n\\\"file\\\":\\\"./m1.js\\\",\\n\\\"lineCount\\\":1,\\n\\\"mappings\\\":\\\";\\\",\\n\\\"sources\\\":[],\\n\\\"names\\\":[]\\n}\\n\"},{\"src\":\"alert(\\\"foo\\\");\\n\",\"path\":\"./m2.js\",\"source_map\":\"{\\n\\\"version\\\":3,\\n\\\"file\\\":\\\"./m2.js\\\",\\n\\\"lineCount\\\":1,\\n\\\"mappings\\\":\\\"AAAAA,KAAA,CAAM,KAAN;\\\",\\n\\\"sources\\\":[\\\"foo.js\\\"],\\n\\\"sourcesContent\\\":[\\\"alert('foo');\\\"],\\n\\\"names\\\":[\\\"alert\\\"]\\n}\\n\"}]")
}

func (c *CommandLineRunnerTest) TestOptionalCatch() {
	c.argsAdd("--language_in=ECMASCRIPT_2019")
	c.argsAdd("--language_out=ECMASCRIPT_2018")
	c.test("try { x(); } catch {}", "try{x()}catch(a){}")
}

func (c *CommandLineRunnerTest) TestChunkOutputFiles() {
	inDir := jutil.NewFile(c.t.TempDir())
	outDir := jutil.NewFile(c.t.TempDir())
	inputFile1 := jutil.NewFileWithParent(inDir, "input1.js")
	inputSource1 := `var x=1;
`
	c.writeFile(inputFile1, inputSource1) // Java: Files.asCharSink(inputFile1, UTF_8).write(inputSource1)
	inputFile2 := jutil.NewFileWithParent(inDir, "input2.js")
	inputSource2 := `var y=2;
`
	c.writeFile(inputFile2, inputSource2) // Java: Files.asCharSink(inputFile2, UTF_8).write(inputSource2)
	outputFile1 := jutil.NewFileWithParent(outDir, "a.js")
	outputFile2 := jutil.NewFileWithParent(outDir, "b.js")
	weakFile := jutil.NewFileWithParent(outDir, jscomp.JSChunk_WEAK_CHUNK_NAME+".js")
	c.argsAdd("--chunk_output_path_prefix")
	c.argsAdd(outDir.String() + "/")
	c.argsAdd("--chunk=a:1")
	c.argsAdd("--chunk=b:1")
	c.argsAdd("--js")
	c.argsAdd(inputFile1.String())
	c.argsAdd("--js")
	c.argsAdd(inputFile2.String())
	runner := NewCommandLineRunnerWithOutAndErr(
		append([]string{}, c.args...), jutil.NewPrintStream(c.outReader), jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		jutil.Throwable_PrintStackTrace(e, jutil.System_Err())
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})

	testjunit.AssertThat(jutil.Files_AsCharSource(outputFile1, jutil.StandardCharsets_UTF_8).Read()).IsEqualTo(inputSource1)
	testjunit.AssertThat(jutil.Files_AsCharSource(outputFile2, jutil.StandardCharsets_UTF_8).Read()).IsEqualTo(inputSource2)
	testjunit.AssertThat(weakFile.Exists()).IsFalse()
}

func (c *CommandLineRunnerTest) TestAssumeStaticInheritanceIsNotUsed() {
	c.testSame("")
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetAssumeStaticInheritanceIsNotUsed()).IsTrue()
	c.argsAdd("--assume_static_inheritance_is_not_used=false")
	c.testSame("")
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetAssumeStaticInheritanceIsNotUsed()).IsFalse()
}

func (c *CommandLineRunnerTest) TestRewriteAndIsolatePolyfills() {
	c.testSame("")
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetRewritePolyfills()).IsTrue()
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetIsolatePolyfills()).IsFalse()
	c.argsAdd("--isolate_polyfills=true")
	c.testSame("")
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetRewritePolyfills()).IsTrue()
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetIsolatePolyfills()).IsTrue()
}

func (c *CommandLineRunnerTest) TestCrossChunkCodeMotionNoStubMethods() {
	c.testSame("")
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetCrossChunkCodeMotionNoStubMethods()).IsFalse()
	c.argsAdd("--assume_no_prototype_method_enumeration=true")
	c.testSame("")
	testjunit.AssertThat(c.lastCompiler.GetOptions().GetCrossChunkCodeMotionNoStubMethods()).IsTrue()
}

func (c *CommandLineRunnerTest) TestParamModification1() {
	c.argsAdd("--compilation_level=ADVANCED")
	c.test(
		`function substr (value, begin, end) {
  return value.slice(begin, end)
}
window.bug = function (s, i) {
  return substr(s, i, i = 5);
}
`,
		"window.a=function(b,c){return b.slice(c,5);};")
}

func (c *CommandLineRunnerTest) TestParamModification2() {
	c.argsAdd("--compilation_level=ADVANCED")
	c.test(
		`function substr (value, begin, end) {
  return value.slice(begin, end)
}
window.bug = function (s, i) {
  return substr(s, i, (s='',i=5));
}
`,
		"window.a=function(b,c){return b.slice(c,5)}")
}

func (c *CommandLineRunnerTest) TestParamModification3() {
	c.argsAdd("--compilation_level=SIMPLE")
	c.test(
		`function substr (value, begin, end) {
  return value.slice(begin, end)
}
window.bug = function (s, i) {
  return substr(s, i, (s='',i=5));
}
`,
		`function substr(a,b,c){
return a.slice(b,c)}
window.bug=function(a,b){return substr(a,b,5)}
`)
}

func (c *CommandLineRunnerTest) TestParamModification4() {
	c.argsAdd("--compilation_level=ADVANCED")
	c.test(
		`function substr (value, begin, end, a, b, c) {
  return value.slice(begin, end, a, b, c)
}
window.bug = function (s, i) {
  return substr(s, i, i=5, i, i=7, i);
}
`,
		"window.a=function(c,b){var d=b,e=b=5,f=b,g=b=7;return c.slice(d,e,f,g,b)}")
}

func (c *CommandLineRunnerTest) TestParamModification5() {
	c.argsAdd("--compilation_level=ADVANCED")
	c.test(
		`function substr (value, begin, end) {
  return value.slice(begin, end)
}
window.bug = function (a, b, c) {
  return substr(a, b+1, b=c);
}
`,
		"window.a=function(b,c,d){return b.slice(c+1,d)}")
}

func (c *CommandLineRunnerTest) TestTranspileOnlyModePolyfillInjection() {
	c.argsAdd("--compilation_level=TRANSPILE_ONLY")
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.test(
		`const arr = [1, 2, 3];
const found = arr.find((element) => element > 10);
`,
		`var $jscomp = $jscomp || {};
$jscomp.scope = {};
$jscomp.findInternal = function(array, callback, thisArg) {
  if (array instanceof String) array = String(array);
  var len = array.length;
  for (var i = 0; i < len; i++) {
    var value = array[i];
    if (callback.call(thisArg, value, i, array)) return {
        i: i, v: value
      }
  }
  return {
    i: -1, v: void 0
  }
};
$jscomp.ASSUME_ES2020 = false;
$jscomp.ASSUME_ES6 = $jscomp.ASSUME_ES2020;
$jscomp.ASSUME_ES5 = $jscomp.ASSUME_ES6;
$jscomp.ASSUME_NO_NATIVE_MAP = false;
$jscomp.ASSUME_NO_NATIVE_SET = false;
$jscomp.ISOLATE_POLYFILLS = false;
$jscomp.FORCE_POLYFILL_PROMISE = false;
$jscomp.FORCE_POLYFILL_PROMISE_WHEN_NO_UNHANDLED_REJECTION = false;
$jscomp.INSTRUMENT_ASYNC_CONTEXT = true;
$jscomp.defineProperty =
    $jscomp.ASSUME_ES5 || typeof Object.defineProperties == 'function' ?
    Object.defineProperty :
    function(target, property, descriptor) {
      if (target == Array.prototype || target == Object.prototype)
        return target;
      target[property] = descriptor.value;
      return target
    };
$jscomp.getGlobal = function(passedInThis) {
  var possibleGlobals = [
    'object' == typeof globalThis && globalThis, passedInThis,
    'object' == typeof window && window, 'object' == typeof self && self,
    'object' == typeof global && global
  ];
  for (var i = 0; i < possibleGlobals.length; ++i) {
    var maybeGlobal = possibleGlobals[i];
    if (maybeGlobal && maybeGlobal['Math'] == Math) return maybeGlobal
  }
  return {
    valueOf: function() {
      throw new Error('Cannot find global object');
    }
  }.valueOf()
};
$jscomp.global = $jscomp.ASSUME_ES2020 ? globalThis : $jscomp.getGlobal(this);
$jscomp.IS_SYMBOL_NATIVE =
    typeof Symbol === 'function' && typeof Symbol('x') === 'symbol';
$jscomp.TRUST_ES6_POLYFILLS =
    !$jscomp.ISOLATE_POLYFILLS || $jscomp.IS_SYMBOL_NATIVE;
$jscomp.polyfills = {};
$jscomp.propertyToPolyfillSymbol = {};
$jscomp.POLYFILL_PREFIX = '$jscp$';
var $jscomp$lookupPolyfilledValue = function(
    target, property, isOptionalAccess) {
  if (isOptionalAccess && target == null) return undefined;
  var obfuscatedName = $jscomp.propertyToPolyfillSymbol[property];
  if (obfuscatedName == null) return target[property];
  var polyfill = target[obfuscatedName];
  return polyfill !== undefined ? polyfill : target[property]
};
$jscomp.TYPED_ARRAY_CLASSES = function() {
  var classes = [
    'Int8', 'Uint8', 'Uint8Clamped', 'Int16', 'Uint16', 'Int32', 'Uint32',
    'Float32', 'Float64'
  ];
  if ($jscomp.global.BigInt64Array) {
    classes.push('BigInt64');
    classes.push('BigUint64')
  }
  return classes
}();
$jscomp.polyfillTypedArrayMethod = function(
    methodName, polyfill, fromLang, toLang) {
  if (!polyfill) return;
  for (var i = 0; i < $jscomp.TYPED_ARRAY_CLASSES.length; i++) {
    var target =
        $jscomp.TYPED_ARRAY_CLASSES[i] + 'Array.prototype.' + methodName;
    if ($jscomp.ISOLATE_POLYFILLS)
      $jscomp.polyfillIsolated(target, polyfill, fromLang, toLang);
    else
      $jscomp.polyfillUnisolated(target, polyfill, fromLang, toLang)
  }
};
$jscomp.polyfill = function(target, polyfill, fromLang, toLang) {
  if (!polyfill) return;
  if ($jscomp.ISOLATE_POLYFILLS)
    $jscomp.polyfillIsolated(target, polyfill, fromLang, toLang);
  else
    $jscomp.polyfillUnisolated(target, polyfill, fromLang, toLang)
};
$jscomp.polyfillUnisolated = function(target, polyfill, fromLang, toLang) {
  var obj = $jscomp.global;
  var split = target.split('.');
  for (var i = 0; i < split.length - 1; i++) {
    var key = split[i];
    if (!(key in obj)) return;
    obj = obj[key]
  }
  var property = split[split.length - 1];
  var orig = obj[property];
  var impl = polyfill(orig);
  if (impl == orig || impl == null) return;
  $jscomp.defineProperty(
      obj, property, {configurable: true, writable: true, value: impl})
};
$jscomp.polyfillIsolated = function(target, polyfill, fromLang, toLang) {
  var split = target.split('.');
  var isSimpleName = split.length === 1;
  var root = split[0];
  if (!isSimpleName && root in $jscomp.polyfills)
    var ownerObject = $jscomp.polyfills;
  else
    ownerObject = $jscomp.global;
  for (var i = 0; i < split.length - 1; i++) {
    var key = split[i];
    if (!(key in ownerObject)) return;
    ownerObject = ownerObject[key]
  }
  var property = split[split.length - 1];
  var nativeImpl = $jscomp.IS_SYMBOL_NATIVE && fromLang === 'es6' ?
      ownerObject[property] :
      null;
  var impl = polyfill(nativeImpl);
  if (impl == null) return;
  if (isSimpleName)
    $jscomp.defineProperty(
        $jscomp.polyfills, property,
        {configurable: true, writable: true, value: impl});
  else if (impl !== nativeImpl) {
    if ($jscomp.propertyToPolyfillSymbol[property] === undefined) {
      var BIN_ID = Math.random() * 1E9 >>> 0;
      $jscomp.propertyToPolyfillSymbol[property] = $jscomp.IS_SYMBOL_NATIVE ?
          $jscomp.global['Symbol'](property) :
          $jscomp.POLYFILL_PREFIX + BIN_ID + '$' + property
    }
    var obfuscatedName = $jscomp.propertyToPolyfillSymbol[property];
    $jscomp.defineProperty(
        ownerObject, obfuscatedName,
        {configurable: true, writable: true, value: impl})
  }
};
$jscomp.polyfill('Array.prototype.find', function(orig) {
  if (orig) return orig;
  var polyfill = function(callback, opt_thisArg) {
    return $jscomp.findInternal(this, callback, opt_thisArg).v
  };
  return polyfill
}, 'es6', 'es3');
var arr = [1, 2, 3];
var found = arr.find(function(element) {
  return element > 10
});
`)
}

func (c *CommandLineRunnerTest) TestTranspileOnlyModeDoesNotDoOptimizations() {
	c.argsAdd("--compilation_level=TRANSPILE_ONLY")
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.test(
		"const x = () => { return 1; }; const y = x();",
		`var x = function() {
  return 1;
};
var y = x();
`)
}

func (c *CommandLineRunnerTest) TestTranspileOnlyModePreservesComments() {
	c.argsAdd("--compilation_level=TRANSPILE_ONLY")
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.test(
		`// This is a comment
const x = () => {};`,
		`// This is a comment
var x = function() {};
`)
}

func (c *CommandLineRunnerTest) TestTranspileOnlyModePreservesTypeAnnotations() {
	c.argsAdd("--compilation_level=TRANSPILE_ONLY")
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.test(
		"/** @type {number} */ const x = 1;",
		`/** @type {number} */
var x = 1;
`)
}

func (c *CommandLineRunnerTest) TestTranspileOnlyModeSyntaxError() {
	c.argsAdd("--compilation_level=TRANSPILE_ONLY")
	c.argsAdd("--language_in=ECMASCRIPT6")
	c.argsAdd("--language_out=ECMASCRIPT5")
	c.test("const x = {", jscomp.RhinoErrorReporter_PARSE_ERROR)
}

func (c *CommandLineRunnerTest) TestExpectedDiagnostics_matchErrorAndWarning() {
	c.argsAdd("--jscomp_error=checkTypes")
	c.argsAdd("--jscomp_warning=uselessCode")
	actualErrors := c.testExpectedDiagnostics(
		"/** @type {string} */ var x = 1; var y = 1; y;",
		[]string{
			`input0:1:30: ERROR - \[JSC_TYPE_MISMATCH\] initializing variable
found   : number
required: string`,
			"input0:1:44: WARNING - \\[JSC_USELESS_CODE\\] Suspicious code. This code lacks" +
				" side-effects. Is there a bug?",
		})
	testjunit.AssertThat(actualErrors).IsEmpty()
}

func (c *CommandLineRunnerTest) TestExpectedDiagnostics_partialMatch() {
	c.argsAdd("--jscomp_error=checkTypes")
	actualErrors := c.testExpectedDiagnostics(
		"/** @type {string} */ var x = 1;", []string{"JSC_TYPE_MISMATCH"})
	testjunit.AssertThat(actualErrors).IsEmpty()
}

func (c *CommandLineRunnerTest) TestExpectedDiagnostics_regexMatch() {
	c.argsAdd("--jscomp_error=checkTypes")
	actualErrors := c.testExpectedDiagnostics(
		"/** @type {string} */ var x = 1;",
		[]string{"input0:.* ERROR - \\[JSC_TYPE_MISMATCH\\] .*"})
	testjunit.AssertThat(actualErrors).IsEmpty()
}

func (c *CommandLineRunnerTest) TestExpectedDiagnostics_unmatchedExpectation() {
	c.argsAdd("--jscomp_error=checkTypes")
	actualErrors := c.testExpectedDiagnostics(
		"/** @type {string} */ var x = 1;", []string{".*JSC_OTHER_ERROR.*"})
	testjunit.AssertThat(actualErrors).HasSize(2)
	AssertError(actualErrors[0]).
		HasType(AbstractCommandLineRunner_EXPECTED_DIAGNOSTIC_NOT_FOUND)
	AssertError(actualErrors[0]).
		HasMessage("Expected diagnostic not found: .*JSC_OTHER_ERROR.*")
	AssertError(actualErrors[1]).HasType(jscomp.TypeValidator_TYPE_MISMATCH_WARNING)
}

func (c *CommandLineRunnerTest) TestExpectedDiagnostics_ambiguousMatch() {
	c.argsAdd("--jscomp_error=checkTypes")
	actualErrors := c.testExpectedDiagnostics(
		"/** @type {string} */ var x = 1;",
		[]string{".*JSC_TYPE_MISMATCH.*", ".*initializing variable.*"})
	testjunit.AssertThat(actualErrors).HasSize(1)
	AssertError(actualErrors[0]).HasType(AbstractCommandLineRunner_AMBIGUOUS_EXPECTATION)
	AssertError(actualErrors[0]).
		HasMessage(
			"Multiple expected diagnostics matched the error: \"input0:1:30: ERROR -" +
				` [JSC_TYPE_MISMATCH] initializing variable
` +
				`found   : number
` +
				`required: string
` +
				"\". Matches: \".*JSC_TYPE_MISMATCH.*\"")
}

func (c *CommandLineRunnerTest) TestExpectedDiagnostics_identicalExpectations() {
	c.argsAdd("--jscomp_error=checkTypes")
	actualErrors := c.testExpectedDiagnostics(
		"/** @type {string} */ var x = 1; /** @type {number} */ var y = 'a';",
		[]string{".*JSC_TYPE_MISMATCH.*", ".*JSC_TYPE_MISMATCH.*"})
	testjunit.AssertThat(actualErrors).IsEmpty()
}

func (c *CommandLineRunnerTest) TestExpectedDiagnostics_singleRegexMultipleActuals() {
	c.argsAdd("--jscomp_error=checkTypes")
	actualErrors := c.testExpectedDiagnostics(
		"/** @type {string} */ var x = 1; /** @type {number} */ var y = 'a';",
		[]string{".*JSC_TYPE_MISMATCH.*"})
	testjunit.AssertThat(actualErrors).HasSize(1)
	AssertError(actualErrors[0]).HasType(jscomp.TypeValidator_TYPE_MISMATCH_WARNING)
	AssertError(actualErrors[0]).
		HasMessage(`initializing variable
found   : string
required: number`)
}

// Java: CommandLineRunnerTest.newJsonArrayOfStrings(String...)
func (c *CommandLineRunnerTest) newJsonArrayOfStrings(strings ...string) *jutil.JsonArray {
	a := jutil.NewJsonArray()
	for _, s := range strings {
		a.AddString(jutil.StringPtr(s))
	}
	return a
}

// Java: CommandLineRunnerTest.getJsonObjectFromJson(String)
func (c *CommandLineRunnerTest) getJsonObjectFromJson(json string) *jutil.JsonObject {
	return jutil.NewGson().FromJsonObject(json)
}

// The given flags should be incompatible with `--checks_only`.
//
// Java: CommandLineRunnerTest.verifyFlagsAreIncompatibleWithChecksOnly(ImmutableList<String>)
func (c *CommandLineRunnerTest) verifyFlagsAreIncompatibleWithChecksOnly(flags []string) {
	additionalFlag := "--checks_only"
	c.verifyFlagConflictIsReported(flags, additionalFlag)
}

// Java: CommandLineRunnerTest.verifyFlagConflictIsReported(ImmutableList<String>, String)
func (c *CommandLineRunnerTest) verifyFlagConflictIsReported(flags []string, additionalFlag string) {
	combinedFlags := append(append([]string{}, flags...), additionalFlag)
	checksOnlyRunner := NewCommandLineRunner(c.stringListToArray(combinedFlags))
	testjunit.AssertThrows[*jscomp.FlagUsageException](func() { checksOnlyRunner.DoRun() })
}

// Java: CommandLineRunnerTest.stringListToArray(ImmutableList<String>)
func (c *CommandLineRunnerTest) stringListToArray(stringList []string) []string {
	return append([]string{}, stringList...)
}

// Java: CommandLineRunnerTest.createStringList(Iterable<String>, String[])
func (c *CommandLineRunnerTest) createStringList(someStrings []string, additionalStrings []string) []string {
	return append(append([]string{}, someStrings...), additionalStrings...)
}

// Java: CommandLineRunnerTest.writeFile(File, String)
func (c *CommandLineRunnerTest) writeFile(file *jutil.File, content string) {
	if err := os.WriteFile(file.GetPath(), jutil.StringGetBytesUTF8(content), 0o666); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
}

/* Helper functions */

// test dispatches to the Java overloads of test by the arguments' types:
// test(String, String), test(String[], String[]),
// test(String[], String[], DiagnosticType), test(String, DiagnosticType),
// test(String, String, DiagnosticType) and test(String[], DiagnosticType).
//
// Go-only: test-helper — overloads.
func (c *CommandLineRunnerTest) test(original any, rest ...any) {
	asArray := func(v any) []string {
		switch x := v.(type) {
		case string:
			return []string{x}
		case []string:
			return x
		}
		panic("bad argument")
	}
	switch len(rest) {
	case 1:
		if w, ok := rest[0].(*jscomp.DiagnosticType); ok {
			c.testWarning(asArray(original), w)
			return
		}
		c.testWithWarning(asArray(original), asArray(rest[0]), nil)
	case 2:
		c.testWithWarning(asArray(original), asArray(rest[0]), rest[1].(*jscomp.DiagnosticType))
	default:
		panic("bad arguments")
	}
}

// Java: CommandLineRunnerTest.testSame(String) and testSame(String[])
func (c *CommandLineRunnerTest) testSame(original any) {
	c.test(original, original)
}

// testWithWarning asserts that when compiling with the given compiler
// options, original is transformed into compiled. If warning is non-null,
// we will also check if the given warning type was emitted.
//
// Java: CommandLineRunnerTest.test(String[], String[], DiagnosticType)
func (c *CommandLineRunnerTest) testWithWarning(original []string, compiled []string, warning *jscomp.DiagnosticType) {
	c.exitCodes = c.exitCodes[:0]
	compiler := c.compile(original)

	testjunit.AssertThat(c.exitCodes).HasSize(1)
	if c.exitCodes[0] != 0 {
		panic(jutil.NewAssertionError(
			"Got nonzero exit code " +
				strconv.Itoa(c.exitCodes[0]) +
				"\nContents of err printstream:\n" +
				c.errReader.String()))
	}

	if warning == nil {
		testjunit.AssertWithMessage(
			"Expected no warnings or errors"+"\nErrors: \n%s\nWarnings: \n%s",
			commandLineRunnerTest_join(compiler.GetErrors()),
			commandLineRunnerTest_join(compiler.GetWarnings())).
			That(len(compiler.GetErrors()) + len(compiler.GetWarnings())).
			IsEqualTo(0)
	} else {
		testjunit.AssertThat(compiler.GetWarnings()).HasSize(1)
		testjunit.AssertThat(compiler.GetWarnings()[0].Type()).IsEqualTo(warning)
	}

	root := compiler.GetRoot().GetLastChild()
	if c.useStringComparison {
		testjunit.AssertThat(compiler.ToSource()).IsEqualTo(strings.Join(compiled, ""))
	} else {
		expectedRoot := c.parse(compiled)
		AssertNode(root).UsingSerializer(compiler.ToSourceOfNode).IsEqualTo(expectedRoot)
	}
}

func commandLineRunnerTest_join(errors []*jscomp.JSError) string {
	parts := make([]string, len(errors))
	for i, e := range errors {
		parts[i] = e.String()
	}
	return strings.Join(parts, "\n")
}

// testWarning asserts that when compiling, there is an error or warning.
//
// Java: CommandLineRunnerTest.test(String[], DiagnosticType)
func (c *CommandLineRunnerTest) testWarning(original []string, warning *jscomp.DiagnosticType) {
	compiler := c.compile(original)
	testjunit.AssertWithMessage(
		"Expected exactly one warning or error "+"\nErrors: \n%s\nWarnings: \n%s",
		commandLineRunnerTest_join(compiler.GetErrors()),
		commandLineRunnerTest_join(compiler.GetWarnings())).
		That(len(compiler.GetErrors()) + len(compiler.GetWarnings())).
		IsEqualTo(1)

	testjunit.AssertThat(c.exitCodes).IsNotEmpty()
	lastExitCode := c.exitCodes[len(c.exitCodes)-1]

	if len(compiler.GetErrors()) != 0 {
		testjunit.AssertThat(compiler.GetErrors()).HasSize(1)
		testjunit.AssertThat(compiler.GetErrors()[0].Type()).IsEqualTo(warning)
		testjunit.AssertWithMessage("Expected exit code of 1.  Contents of err printstream:\n%s", c.errReader.String()).
			That(lastExitCode).
			IsEqualTo(1)
	} else {
		testjunit.AssertThat(compiler.GetWarnings()).HasSize(1)
		testjunit.AssertThat(compiler.GetWarnings()[0].Type()).IsEqualTo(warning)
		testjunit.AssertWithMessage("Expected exit code of 0.  Contents of err printstream:\n%s", c.errReader.String()).
			That(lastExitCode).
			IsEqualTo(0)
	}
}

// Java: CommandLineRunnerTest.createCommandLineRunner(String[])
func (c *CommandLineRunnerTest) createCommandLineRunner(original []string) *CommandLineRunner {
	for i := 0; i < len(original); i++ {
		c.argsAdd("--js")
		c.argsAdd("/path/to/input" + strconv.Itoa(i) + ".js")
		if c.useChunks == commandLineRunnerTest_ChunkPattern_CHAIN {
			c.argsAdd("--chunk")
			dep := ""
			if i > 0 {
				dep = ":m" + strconv.Itoa(i-1)
			}
			c.argsAdd("m" + strconv.Itoa(i) + ":1" + dep)
		} else if c.useChunks == commandLineRunnerTest_ChunkPattern_STAR {
			c.argsAdd("--chunk")
			dep := ""
			if i > 0 {
				dep = ":m0"
			}
			c.argsAdd("m" + strconv.Itoa(i) + ":1" + dep)
		}
	}

	if c.lastArg != nil {
		c.argsAdd(*c.lastArg)
	}

	argStrings := append([]string{}, c.args...)
	return NewCommandLineRunnerWithOutAndErr(
		argStrings, jutil.NewPrintStream(c.outReader), jutil.NewPrintStream(c.errReader))
}

// Java: CommandLineRunnerTest.createZipFile(Map<String, String>)
func (c *CommandLineRunnerTest) createZipFile(entryContentsByName *jutil.LinkedHashMap[string, string]) *AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType] {
	dir, err := os.MkdirTemp("", "jscomp")
	if err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	tempZipFile, err := os.CreateTemp(dir, "testdata*.js.zip")
	if err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	zipOutputStream := zip.NewWriter(tempZipFile)
	for name, content := range entryContentsByName.All() {
		w, err := zipOutputStream.Create(name)
		if err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
		w.Write(jutil.StringGetBytesUTF8(content))
	}
	zipOutputStream.Close()
	tempZipFile.Close()

	return NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS_ZIP, tempZipFile.Name())
}

// Java: CommandLineRunnerTest.createJsFile(String, String)
func (c *CommandLineRunnerTest) createJsFile(filename string, fileContent string) *AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType] {
	dir, err := os.MkdirTemp("", "jscomp")
	if err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	tempJsFile, err := os.CreateTemp(dir, filename+"*.js")
	if err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	tempJsFile.Write(jutil.StringGetBytesUTF8(fileContent))
	tempJsFile.Close()

	return NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_JS, tempJsFile.Name())
}

// compileFiles is a helper for compiling from zip and js files and checking
// output string.
//
// Java: CommandLineRunnerTest.compileFiles(String, FlagEntry<JsSourceType>...)
func (c *CommandLineRunnerTest) compileFiles(expectedOutput string, entries ...*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]) {
	c.setupFlags(entries...)
	c.compileArgs(expectedOutput, nil)
}

// Java: CommandLineRunnerTest.setupFlags(FlagEntry<JsSourceType>...)
func (c *CommandLineRunnerTest) setupFlags(entries ...*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]) {
	for _, entry := range entries {
		c.argsAdd("--" + entry.GetFlag().FlagName() + "=" + entry.GetValue())
	}
}

// compileJsFiles is a helper for compiling js files and checking output
// string, using a single --js flag.
//
// Java: CommandLineRunnerTest.compileJsFiles(String, FlagEntry<JsSourceType>...)
func (c *CommandLineRunnerTest) compileJsFiles(expectedOutput string, entries ...*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]) {
	c.argsAdd("--js")
	for _, entry := range entries {
		c.argsAdd(entry.GetValue())
	}
	c.compileArgs(expectedOutput, nil)
}

// Java: CommandLineRunnerTest.compileArgs(String, DiagnosticType)
func (c *CommandLineRunnerTest) compileArgs(expectedOutput string, expectedError *jscomp.DiagnosticType) {
	argStrings := append([]string{}, c.args...)

	runner := NewCommandLineRunnerWithOutAndErr(argStrings, jutil.NewPrintStream(c.outReader), jutil.NewPrintStream(c.errReader))
	c.lastCompiler = runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})
	compiler := runner.GetCompiler()
	output := compiler.ToSource()
	if expectedError == nil {
		testjunit.AssertThat(compiler.GetErrors()).IsEmpty()
		testjunit.AssertThat(compiler.GetWarnings()).IsEmpty()
		testjunit.AssertThat(output).IsEqualTo(expectedOutput)
	} else {
		testjunit.AssertThat(compiler.GetErrors()).HasSize(1)
		AssertError(compiler.GetErrors()[0]).HasType(expectedError)
	}
	c.lastCommandLineRunner = runner
}

// compileString is compile(String, List<String>).
//
// Java: CommandLineRunnerTest.compile(String, List<String>)
func (c *CommandLineRunnerTest) compileString(inputString string, args []string) string {
	outputStream := &bytes.Buffer{}
	errorStream := &bytes.Buffer{}
	runner := NewCommandLineRunnerWithStreams(
		append([]string{}, args...),
		bytes.NewReader(jutil.StringGetBytesUTF8(inputString)),
		jutil.NewPrintStream(outputStream),
		jutil.NewPrintStream(errorStream))

	runner.GetCompiler()
	jutil.Try(func() {
		runner.DoRun()
	}, func(e jutil.IIOException) {
		testjunit.AssertWithMessage("Unexpected exception %s", e).Fail()
	})
	return jutil.NewStringFromUTF8Bytes(outputStream.Bytes())
}

// Java: CommandLineRunnerTest.compile(String[])
func (c *CommandLineRunnerTest) compile(original []string) *jscomp.Compiler {
	runner := c.createCommandLineRunner(original)
	if !runner.ShouldRunCompiler() {
		testjunit.AssertThat(runner.HasErrors()).IsTrue()
		testjunit.AssertWithMessage(jutil.NewStringFromUTF8Bytes(c.errReader.Bytes())).Fail()
	}
	var inputsSupplier func() []*jscomp.SourceFile
	var chunksSupplier func() []*jscomp.JSChunk

	switch c.useChunks {
	case commandLineRunnerTest_ChunkPattern_NONE:
		var inputs []*jscomp.SourceFile
		for i := 0; i < len(original); i++ {
			inputs = append(inputs, jscomp.SourceFile_FromCode(c.getFilename(i), original[i]))
		}
		inputsSupplier = func() []*jscomp.SourceFile { return inputs }
	case commandLineRunnerTest_ChunkPattern_STAR:
		chunks := JSChunkGraphBuilder_ForStar().AddChunks(original).Build()
		chunksSupplier = func() []*jscomp.JSChunk { return chunks }
	case commandLineRunnerTest_ChunkPattern_CHAIN:
		chunks := JSChunkGraphBuilder_ForChain().AddChunks(original).Build()
		chunksSupplier = func() []*jscomp.JSChunk { return chunks }
	}

	externs := c.externs
	runner.EnableTestMode(
		func() []*jscomp.SourceFile { return externs },
		inputsSupplier,
		chunksSupplier,
		func(exitCode int) {
			c.exitCodes = append(c.exitCodes, exitCode)
		})
	runner.Run()
	c.lastCompiler = runner.GetCompiler()
	c.lastCommandLineRunner = runner
	return c.lastCompiler
}

// Java: CommandLineRunnerTest.testExpectedDiagnostics(String, String[])
func (c *CommandLineRunnerTest) testExpectedDiagnostics(input string, expectedDiagnostics []string) []*jscomp.JSError {
	for _, diag := range expectedDiagnostics {
		c.argsAdd("--expected_diagnostics=" + diag)
	}
	c.compile([]string{input})

	compiler := c.lastCompiler
	return append(append([]*jscomp.JSError{}, compiler.GetErrors()...), compiler.GetWarnings()...)
}

// Java: CommandLineRunnerTest.parse(String[])
func (c *CommandLineRunnerTest) parse(original []string) rhino.Node {
	argStrings := append([]string{}, c.args...)
	runner := NewCommandLineRunner(argStrings)
	compiler := runner.CreateCompiler()
	var inputs []*jscomp.SourceFile
	for i := 0; i < len(original); i++ {
		inputs = append(inputs, jscomp.SourceFile_FromCode(c.getFilename(i), original[i]))
	}
	options := runner.CreateOptions()
	compiler.Init(c.externs, inputs, options)
	all := compiler.ParseInputs()
	testjunit.AssertThat(compiler.GetErrors()).IsEmpty()
	rhino.CheckNotNull(all)
	return all.GetLastChild()
}

// Java: CommandLineRunnerTest.setFilename(int, String)
func (c *CommandLineRunnerTest) setFilename(i int, filename string) {
	c.filenames[i] = filename
}

// Java: CommandLineRunnerTest.getFilename(int)
func (c *CommandLineRunnerTest) getFilename(i int) string {
	if len(c.filenames) == 0 {
		return "input" + strconv.Itoa(i)
	}
	name, ok := c.filenames[i]
	jutil.CheckState(ok && name != "")
	return name
}

// Java: CommandLineRunnerTest.concatStrings(String...)
func (c *CommandLineRunnerTest) concatStrings(strs ...string) string {
	return strings.Join(strs, "")
}

// commandLineRunnerTest_createTempFile is File.createTempFile(prefix, suffix).
func commandLineRunnerTest_createTempFile(prefix, suffix string) *jutil.File {
	f, err := os.CreateTemp("", prefix+"*"+suffix)
	if err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
	f.Close()
	return jutil.NewFile(f.Name())
}

// commandLineRunnerTest_readLines is Guava Files.readLines(File, UTF_8).
func commandLineRunnerTest_readLines(f *jutil.File) []string {
	content := jutil.Files_ReadStringUTF8(f.ToPath())
	lines := strings.Split(strings.ReplaceAll(strings.ReplaceAll(content, "\r\n", "\n"), "\r", "\n"), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

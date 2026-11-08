// Port of test/com/google/javascript/jscomp/integration/ES2022IntegrationTest.java @ closure-compiler v20261104
// Port of test/com/google/javascript/jscomp/integration/IntegrationTestCase.java @ closure-compiler v20261104

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

/*
 * Copyright 2012 The Closure Compiler Authors.
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
	"runtime/debug"
	"strings"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: IntegrationTestCase.defaultExternsBuilderwithoutClosure()
func es2022IntegrationTestCase_defaultExternsBuilderWithoutClosure() *TestExternsBuilder {
	return NewTestExternsBuilder().
		AddArguments().
		AddString().
		AddObject().
		AddReflect().
		AddFunction().
		AddIterable().
		AddPromise().
		AddArray().
		AddAlert().
		AddMap().
		AddExtra(
			`/**
 * @const
 */
var Math = {};
/**
 * @param {?} n1
 * @param {?} n2
 * @return {number}
 * @nosideeffects
 */
Math.pow = function(n1, n2) {};
Math.random = function() {}
var isNaN;
var Infinity;
/**
 * @constructor
 * @extends {Array<string>}
 */
var ITemplateArray = function() {};
/** @constructor */
var Set;
/** @constructor */ function Window() {}
/** @type {string} */ Window.prototype.name;
/** @type {string} */ Window.prototype.offsetWidth;
/** @type {Window} */ var window;

/** @nosideeffects */ function noSideEffects() {}

/**
 * @constructor
 * @nosideeffects
 */
function Widget() {}
/** @modifies {this} */ Widget.prototype.go = function() {};
/** @return {string} */ var widgetToken = function() {};

/**
 * @constructor
 * @return {number}
 * @param {*=} opt_n
 */
function Number(opt_n) {}

/**
 * @constructor
 * @return {boolean}
 * @param {*=} opt_b
 */
function Boolean(opt_b) {}

/**
 * @constructor
 * @return {!TypeError}
 * @param {*=} message
 * @param {*=} fileNameOrOptions
 * @param {*=} line
 */
function TypeError(message, fileNameOrOptions, line) {}
/**
 * @constructor
 * @param {*=} message
 * @param {*=} fileNameOrOptions
 * @param {*=} line
 * @return {!Error}
 * @nosideeffects
 */
function Error(message, fileNameOrOptions, line) {}

/** @constructor */
var HTMLElement = function() {};
`)
}

// Java: IntegrationTestCase.DEFAULT_EXTERNS
var es2022IntegrationTestCase_defaultExterns = sync.OnceValue(func() []*SourceFile {
	return []*SourceFile{
		es2022IntegrationTestCase_defaultExternsBuilderWithoutClosure().
			AddClosureExterns().
			BuildExternsFile("externs"),
	}
})

// es2022IntegrationTestCase ports IntegrationTestCase for ES2022IntegrationTest.
//
// Java: IntegrationTestCase
type es2022IntegrationTestCase struct {
	junitT              *testing.T
	externs             []*SourceFile
	lastCompiler        *Compiler
	inputFileNamePrefix string
	inputFileNameSuffix string
}

func (c *es2022IntegrationTestCase) JUnitSetT(t *testing.T) { c.junitT = t }
func (c *es2022IntegrationTestCase) JUnitBefore()           { c.SetUp() }
func (c *es2022IntegrationTestCase) JUnitAfter()            {}

// Java: IntegrationTestCase.setUp()
func (c *es2022IntegrationTestCase) SetUp() {
	c.externs = es2022IntegrationTestCase_defaultExterns()
	c.lastCompiler = nil
	c.inputFileNamePrefix = "i"
	c.inputFileNameSuffix = ".js"
}

// Java: IntegrationTestCase.testSame(CompilerOptions, String)
func (c *es2022IntegrationTestCase) TestSame(options *CompilerOptions, original string) {
	c.TestSameSources(options, []string{original})
}

// Java: IntegrationTestCase.testSame(CompilerOptions, String[])
func (c *es2022IntegrationTestCase) TestSameSources(options *CompilerOptions, original []string) {
	c.TestSources(options, original, original)
}

// Java: IntegrationTestCase.test(CompilerOptions, String, String)
func (c *es2022IntegrationTestCase) Test(options *CompilerOptions, original string, compiled string) {
	c.TestSources(options, []string{original}, []string{compiled})
}

// Java: IntegrationTestCase.test(CompilerOptions, String[], String[])
func (c *es2022IntegrationTestCase) TestSources(options *CompilerOptions, original []string, compiled []string) {
	compiler := c.CompileSources(options, original)

	root := compiler.GetRoot().GetLastChild()

	errors := compiler.GetErrors()
	warnings := compiler.GetWarnings()
	var errStrs, warnStrs []string
	for _, e := range errors {
		errStrs = append(errStrs, e.String())
	}
	for _, w := range warnings {
		warnStrs = append(warnStrs, w.String())
	}
	AssertWithMessage(
		"Expected no warnings or errors\nErrors:\n%s\nWarnings:\n%s\n",
		strings.Join(errStrs, "\n"),
		strings.Join(warnStrs, "\n")).
		That(len(errors) + len(warnings)).
		IsEqualTo(0)

	if compiled != nil {
		expectedRoot := c.ParseExpectedCode(compiled, options)
		AssertNode(root).UsingSerializer(compiler.ToSourceOfNode).IsEqualTo(expectedRoot)
	}
}

// Java: IntegrationTestCase.test(CompilerOptions, String, DiagnosticGroup)
func (c *es2022IntegrationTestCase) TestWarning(options *CompilerOptions, original string, warning *DiagnosticGroup) {
	c.TestSourcesWithDiagnosticGroup(options, []string{original}, nil, warning)
}

// Java: IntegrationTestCase.test(CompilerOptions, String[], String[], DiagnosticGroup[])
func (c *es2022IntegrationTestCase) TestSourcesWithDiagnosticGroups(
	options *CompilerOptions, original []string, compiled []string, warnings []*DiagnosticGroup,
) {
	compiler := c.CompileSources(options, original)
	c.CheckUnexpectedErrorsOrWarnings(compiler, len(warnings))

	if compiled != nil {
		root := compiler.GetRoot().GetLastChild()
		expectedRoot := c.ParseExpectedCode(compiled, options)
		AssertNode(root).UsingSerializer(compiler.ToSourceOfNode).IsEqualTo(expectedRoot)
	}

	var all []*JSError
	all = append(all, compiler.GetErrors()...)
	all = append(all, compiler.GetWarnings()...)
	AssertThat(all).
		ComparingElementsUsing(JSCompCorrespondences_OWNING_DIAGNOSTIC_GROUP).
		ContainsExactlyElementsIn(warnings)
}

// Java: IntegrationTestCase.test(CompilerOptions, String[], String[], DiagnosticGroup)
func (c *es2022IntegrationTestCase) TestSourcesWithDiagnosticGroup(
	options *CompilerOptions, original []string, compiled []string, warnings *DiagnosticGroup,
) {
	compiler := c.CompileSources(options, original)
	c.CheckUnexpectedErrorsOrWarnings(compiler, 1)

	if compiled != nil {
		root := compiler.GetRoot().GetLastChild()
		expectedRoot := c.ParseExpectedCode(compiled, options)
		AssertNode(root).UsingSerializer(compiler.ToSourceOfNode).IsEqualTo(expectedRoot)
	}

	var diagnostic *DiagnosticType
	if compiler.GetErrorCount() == 1 {
		diagnostic = compiler.GetErrors()[0].Type()
	} else {
		diagnostic = compiler.GetWarnings()[0].Type()
	}
	AssertWithMessage(
		"Error not in expected diagnostic group. Error: %s\nExpected group: %s",
		diagnostic.Key(), warnings).
		That(warnings.MatchesType(diagnostic)).
		IsTrue()
}

// Java: IntegrationTestCase.checkUnexpectedErrorsOrWarnings(Compiler, int)
func (c *es2022IntegrationTestCase) CheckUnexpectedErrorsOrWarnings(compiler *Compiler, expected int) {
	actual := len(compiler.GetErrors()) + len(compiler.GetWarnings())
	if actual != expected {
		var msg strings.Builder
		for _, err := range compiler.GetErrors() {
			msg.WriteString("Error:")
			msg.WriteString(err.String())
			msg.WriteByte('\n')
		}
		for _, err := range compiler.GetWarnings() {
			msg.WriteString("Warning:")
			msg.WriteString(err.String())
			msg.WriteByte('\n')
		}
		AssertWithMessage("Unexpected warnings or errors.\n %s", msg.String()).
			That(actual).
			IsEqualTo(expected)
	}
}

// Java: IntegrationTestCase.compile(CompilerOptions, String)
func (c *es2022IntegrationTestCase) Compile(options *CompilerOptions, original string) *Compiler {
	return c.CompileSources(options, []string{original})
}

// Java: IntegrationTestCase.compile(CompilerOptions, String[])
func (c *es2022IntegrationTestCase) CompileSources(options *CompilerOptions, original []string) *Compiler {
	return c.CompileChunks(
		options,
		JSChunkGraphBuilder_ForChain().
			AddChunks(original).
			SetFilenameFormat(c.inputFileNamePrefix+"%s"+c.inputFileNameSuffix).
			Build())
}

// Java: IntegrationTestCase.compile(CompilerOptions, ImmutableList<JSChunk>)
func (c *es2022IntegrationTestCase) CompileChunks(options *CompilerOptions, chunks []*JSChunk) *Compiler {
	compiler := c.CreateCompilerWithErrorManager(NewBlackHoleErrorManager())
	c.lastCompiler = compiler
	_ = compiler.CompileChunks(c.externs, chunks, options)
	return compiler
}

// Java: IntegrationTestCase.createCompiler(ErrorManager)
func (c *es2022IntegrationTestCase) CreateCompilerWithErrorManager(errorManager ErrorManager) *Compiler {
	return NewCompilerWithErrorManager(errorManager)
}

// Java: IntegrationTestCase.createCompiler()
func (c *es2022IntegrationTestCase) CreateCompiler() *Compiler {
	return NewCompiler()
}

// Java: IntegrationTestCase.testNoWarnings(CompilerOptions, String)
func (c *es2022IntegrationTestCase) TestNoWarnings(options *CompilerOptions, code string) {
	c.TestNoWarningsSources(options, []string{code})
}

// Java: IntegrationTestCase.testNoWarnings(CompilerOptions, String[])
func (c *es2022IntegrationTestCase) TestNoWarningsSources(options *CompilerOptions, sources []string) {
	compiler := c.CompileSources(options, sources)
	AssertThat(compiler.GetErrors()).IsEmpty()
	AssertThat(compiler.GetWarnings()).IsEmpty()
}

// Java: IntegrationTestCase.parseExpectedCode(String[], CompilerOptions)
func (c *es2022IntegrationTestCase) ParseExpectedCode(original []string, options *CompilerOptions) rhino.Node {
	oldProcessCommonJsModules := options.GetProcessCommonJSModules()
	options.SetProcessCommonJSModules(false)
	expectedRoot := c.Parse(original, options)
	options.SetProcessCommonJSModules(oldProcessCommonJsModules)
	return expectedRoot
}

// Java: IntegrationTestCase.parse(String[], CompilerOptions)
func (c *es2022IntegrationTestCase) Parse(original []string, options *CompilerOptions) rhino.Node {
	compiler := c.CreateCompiler()
	inputs := make([]*SourceFile, 0, len(original))
	for i := 0; i < len(original); i++ {
		inputs = append(inputs, SourceFile_FromCode(fmt.Sprintf("%s%d%s", c.inputFileNamePrefix, i, c.inputFileNameSuffix), original[i]))
	}
	compiler.Init(c.externs, inputs, options)
	c.CheckUnexpectedErrorsOrWarnings(compiler, 0)
	compiler.Parse()
	c.CheckUnexpectedErrorsOrWarnings(compiler, 0)

	return compiler.GetRoot().GetLastChild()
}

// ES2022IntegrationTest holds integration tests related to ES2022 features like class fields.
//
// Java: ES2022IntegrationTest
type ES2022IntegrationTest struct {
	es2022IntegrationTestCase
}

// Java: ES2022IntegrationTest()
func NewES2022IntegrationTest() *ES2022IntegrationTest {
	return &ES2022IntegrationTest{}
}

func TestES2022IntegrationTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewES2022IntegrationTest() })
}

// CreateCompilerOptions creates a CompilerOptions object with google coding conventions.
//
// Java: ES2022IntegrationTest.createCompilerOptions()
func (c *ES2022IntegrationTest) CreateCompilerOptions() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	options.SetDevMode(CompilerOptions_DevMode_EVERY_PASS)
	options.SetCodingConvention(NewGoogleCodingConvention())
	return options
}

// Java: ES2022IntegrationTest.checksOnlyCompilerOptions()
func (c *ES2022IntegrationTest) checksOnlyCompilerOptions() *CompilerOptions {
	options := c.CreateCompilerOptions()
	options.SetChecksOnly(true)
	WarningLevel_VERBOSE.SetOptionsForWarningLevel(options)

	return options
}

// Java: ES2022IntegrationTest.fullyOptimizedCompilerOptions()
func (c *ES2022IntegrationTest) fullyOptimizedCompilerOptions() *CompilerOptions {
	options := c.CreateCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetTypeBasedOptimizationOptions(options)
	WarningLevel_VERBOSE.SetOptionsForWarningLevel(options)

	return options
}

// Java: ES2022IntegrationTest.publicClassFields_supportedInChecksOnlyMode()
func (c *ES2022IntegrationTest) TestPublicClassFields_supportedInChecksOnlyMode() {
	options := c.checksOnlyCompilerOptions()
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	c.TestNoWarnings(
		options,
		`class MyClass {
  /** @type {number} */
  x = 2;
  y;
}
console.log(new MyClass().x);
`)
}

// Java: ES2022IntegrationTest.publicClassFields_supportedInChecksOnlyMode2()
func (c *ES2022IntegrationTest) TestPublicClassFields_supportedInChecksOnlyMode2() {
	options := c.checksOnlyCompilerOptions()
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	c.TestNoWarnings(
		options,
		`class MyClass {
  x = '';
  y;
}
console.log(new MyClass().x);
`)
}

// Java: ES2022IntegrationTest.publicClassFields_supportedInChecksOnlyMode3()
func (c *ES2022IntegrationTest) TestPublicClassFields_supportedInChecksOnlyMode3() {
	options := c.checksOnlyCompilerOptions()

	c.TestSourcesWithDiagnosticGroups(
		options,
		[]string{
			`class MyClass {
  /** @type {string} */
  x = 2;
}
`,
		},
		nil,
		[]*DiagnosticGroup{DiagnosticGroups_CHECK_TYPES})
}

// Java: ES2022IntegrationTest.computedPublicClassFields_supportedInChecksOnlyMode()
func (c *ES2022IntegrationTest) TestComputedPublicClassFields_supportedInChecksOnlyMode() {
	options := c.checksOnlyCompilerOptions()
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	c.TestNoWarnings(
		options,
		`/** @dict */
class MyClass {
  [3 + 4] = 5;
  [6];
  'x' = 2;
}
console.log(new MyClass()[6]);
`)
}

// Java: ES2022IntegrationTest.computedPublicClassFields_supportedInChecksOnlyMode2()
func (c *ES2022IntegrationTest) TestComputedPublicClassFields_supportedInChecksOnlyMode2() {
	options := c.checksOnlyCompilerOptions()
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	c.TestNoWarnings(
		options,
		`/** @dict */
class MyClass {
  ['x'] = 5;
}
console.log(new MyClass()['x']);
`)
}

// deprecated warnings aren't given on computed fields:
// users should be careful about tags on computed fields
//
// Java: ES2022IntegrationTest.computedPublicClassFields_supportedInChecksOnlyMode3()
func (c *ES2022IntegrationTest) TestComputedPublicClassFields_supportedInChecksOnlyMode3() {
	options := c.checksOnlyCompilerOptions()

	c.TestNoWarnings(
		options,
		`/** @unrestricted */
class MyClass {
  /** @deprecated */
  ['x'] = 5;
  baz() { return this['x']; }
}
`)
}

// Java: ES2022IntegrationTest.computedPublicClassFields_supportedInChecksOnlyMode4()
func (c *ES2022IntegrationTest) TestComputedPublicClassFields_supportedInChecksOnlyMode4() {
	options := c.checksOnlyCompilerOptions()

	// test will give JSC_ILLEGAL_PROPERTY_ACCESS error because @dict or @restricted is missing
	c.TestSourcesWithDiagnosticGroups(
		options,
		[]string{
			`class MyClass {
  [3 + 4] = 5;
}
`,
		},
		nil,
		[]*DiagnosticGroup{DiagnosticGroups_CHECK_TYPES})
}

// Java: ES2022IntegrationTest.publicMixedClassFields_supportedInChecksOnlyMode()
func (c *ES2022IntegrationTest) TestPublicMixedClassFields_supportedInChecksOnlyMode() {
	options := c.checksOnlyCompilerOptions()

	c.TestNoWarnings(
		options,
		`/** @unrestricted */
class MyClass {
  a = 2;
  ['b'] = 'hi';
  'c' = 5;
  2 = 4;
  d;
  ['e'];
}
`)
}

// Java: ES2022IntegrationTest.publicClassFields_supportedInOptimizationsMode()
func (c *ES2022IntegrationTest) TestPublicClassFields_supportedInOptimizationsMode() {
	options := c.fullyOptimizedCompilerOptions()
	options.SetPrintSourceAfterEachPass(true)

	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	c.Test(
		options,
		`class C {
  f1 = 1;
static f2 = 3;
  m1() {return this.f1}
}
console.log(new C().f1);
console.log(C.f2);
console.log(new C().m1());
`,
		`console.log(1);
console.log(3);
console.log(1);
`)
}

// Java: ES2022IntegrationTest.publicClassFields_supportedInOptimizationsMode1()
func (c *ES2022IntegrationTest) TestPublicClassFields_supportedInOptimizationsMode1() {
	options := c.fullyOptimizedCompilerOptions()
	options.SetPrintSourceAfterEachPass(true)

	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	c.Test(
		options,
		`class MyClass {
  /** @type {number} */
  f1 = 2;
  /** @type {string} */
  f2 = 'hi';
  f3 = function() { return this.f1 };
  m1() { return this.f2; }
}
console.log(new MyClass().f1);
console.log(new MyClass().f2);
console.log(new MyClass().f3());
console.log(new MyClass().m1());
`,
		`console.log(2);
console.log('hi');
console.log(2);
console.log('hi');
`)
}

// Java: ES2022IntegrationTest.publicClassFields_supportedInOptimizationsMode2()
func (c *ES2022IntegrationTest) TestPublicClassFields_supportedInOptimizationsMode2() {
	options := c.fullyOptimizedCompilerOptions()
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	c.Test(
		options,
		`/** @dict */
class MyClass {
  ['f1'] = 2;
  'f2' = 'hi';
  2 = 4;
  ['m1']() { return this['f1']; }
}
console.log(new MyClass()['f1']);
console.log(new MyClass()[2]);
console.log(new MyClass()['m1']());
`,
		`class a {
  f1 = 2;
  f2= 'hi';
  2 = 4;
  m1() { return this.f1; }
}
console.log((new a).f1);
console.log((new a)[2]);
console.log((new a).m1());
`)
}

// Java: ES2022IntegrationTest.publicClassFields_supportedInOptimizationsMode3()
func (c *ES2022IntegrationTest) TestPublicClassFields_supportedInOptimizationsMode3() {
	options := c.fullyOptimizedCompilerOptions()
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	c.Test(
		options,
		`/** @unrestricted */
class MyClass {
  f1 = 1;
  ['a'] = 'hi';
  ['f3'] = function() { return this.f1; };
}
console.log(new MyClass().f1);
console.log(new MyClass()['a']);
console.log(new MyClass()['f3']());
`,
		`class a {
  b = 1;
  a = 'hi';
  f3 = function() { return this.b; };
}
console.log(1);
console.log(new a().a);
console.log(new a().f3());
`)
}

// Java: ES2022IntegrationTest.publicClassFields_supportedInOptimizationsMode4()
func (c *ES2022IntegrationTest) TestPublicClassFields_supportedInOptimizationsMode4() {
	options := c.fullyOptimizedCompilerOptions()

	c.Test(
		options,
		`/** @unrestricted */
class MyClass {
  static f1 = alert(2);
  ['f2'] = 'hi';
  'f3' = 5;
  4 = 4;
  f5;
  ['f6'] = alert(1);
}
`,
		`alert(2);
`)
}

// Java: ES2022IntegrationTest.publicClassFields_newTargetInFieldInitializers()
func (c *ES2022IntegrationTest) TestPublicClassFields_newTargetInFieldInitializers() {
	// `new.target` is not supported in field initializers, where it is always `undefined`.
	options := c.CreateCompilerOptions()
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}
	src := `class C {
  inst = new.target;
  static stat = new.target;
}
`

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.TestSourcesWithDiagnosticGroups(
		options,
		[]string{src},
		nil,
		[]*DiagnosticGroup{DiagnosticGroups_PARSING, DiagnosticGroups_PARSING})

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2022)
	c.TestSourcesWithDiagnosticGroups(
		options,
		[]string{src},
		nil,
		[]*DiagnosticGroup{DiagnosticGroups_PARSING, DiagnosticGroups_PARSING})

	// Also when the class is inside a function.
	srcInFunction := `function f() {
  class C {
    inst = new.target;
    static stat = new.target;
  }
  return C;
}
`

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.TestSourcesWithDiagnosticGroups(
		options,
		[]string{srcInFunction},
		nil,
		[]*DiagnosticGroup{DiagnosticGroups_PARSING, DiagnosticGroups_PARSING})

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2022)
	c.TestSourcesWithDiagnosticGroups(
		options,
		[]string{srcInFunction},
		nil,
		[]*DiagnosticGroup{DiagnosticGroups_PARSING, DiagnosticGroups_PARSING})
}

// Java: ES2022IntegrationTest.publicClassFields_derivedConstructorReturnsObjectWithoutSuperCall()
func (c *ES2022IntegrationTest) TestPublicClassFields_derivedConstructorReturnsObjectWithoutSuperCall() {
	// The fields of `Sub` are never initialized, since its constructor never calls `super()`.
	// Like TypeScript's TS2377, `CheckSuper` reports it, although the constructor returns a value.
	options := c.CreateCompilerOptions()
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	src := `class Base {}
class Sub extends Base {
  y = 2;
  constructor() {
    return {tag: 'ret'};
  }
}
console.log(new Sub().tag);
`
	compiler := c.Compile(options, src)
	c.CheckUnexpectedErrorsOrWarnings(compiler, 1)
	AssertThat(compiler.GetErrors()[0].Type().Key()).
		IsEqualTo("JSC_RETURN_WITHOUT_CALL_TO_SUPER")
}

// Java: ES2022IntegrationTest.computedFieldExecutionOrderAndDeadAssignmentElimination()
func (c *ES2022IntegrationTest) TestComputedFieldExecutionOrderAndDeadAssignmentElimination() {
	options := c.fullyOptimizedCompilerOptions()

	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddConsole().
			AddExtra("var window;").
			BuildExternsFile("externs"),
	}

	c.Test(
		options,
		`window.test = function() {
  var x = 0;
  /** @unrestricted */
  class MyClass {
    static f1 = x;
    static[(x = 1)] = 1;  // (x = 1) executes before assigning 'static f1 = x'
  }
  console.log(MyClass.f1);  // prints 1
};
`,
		`window.b = function() {
  class a{
    static a;
    static [1];
  }
  a.a = 1;
  a[1] = 1;
  console.log(a.a);  // prints 1
}
`)
}

// Java: ES2022IntegrationTest.computedMethodExecutionOrderAndDeadAssignmentElimination()
func (c *ES2022IntegrationTest) TestComputedMethodExecutionOrderAndDeadAssignmentElimination() {
	options := c.fullyOptimizedCompilerOptions()

	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddConsole().
			AddExtra("var window;").
			BuildExternsFile("externs"),
	}

	c.Test(
		options,
		`window.test = function() {
  var x = 0;
  /** @unrestricted */
  class MyClass {
    static f1 = x;
    static[(x = 1)]() {};  // (x = 1) executes before assigning 'static f1 = x'
  }
  console.log(MyClass.f1);  // prints 1
};
`,
		`window.b = function() {
  var a = 0, d = a = 1;
  class c {
    static a;
    static[d]() {}
  }
  c.a = a;
  console.log(c.a)  // prints 1
}
`)
}

// Java: ES2022IntegrationTest.publicClassFields_derivedConstructorParamDefaultReadsThis()
func (c *ES2022IntegrationTest) TestPublicClassFields_derivedConstructorParamDefaultReadsThis() {
	// In a derived class `this` is unbound until super() returns, and parameter defaults run before
	// the body, so `z = this.b` throws a ReferenceError natively instead of reading the field.
	// `CheckSuper` reports it.
	options := c.CreateCompilerOptions()
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)

	src := `class Parent {}
class Child extends Parent {
  b = 2;
  constructor(z = this.b) {
    super();
    this.z = z;
  }
}
`
	compiler := c.Compile(options, src)
	c.CheckUnexpectedErrorsOrWarnings(compiler, 1)
	AssertThat(compiler.GetErrors()[0].Type().Key()).IsEqualTo("JSC_THIS_BEFORE_SUPER")
}

// Java: ES2022IntegrationTest.testEs6RewriteClassExtendsExpression()
func (c *ES2022IntegrationTest) TestEs6RewriteClassExtendsExpression() {
	options := c.CreateCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	options.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_OFF)
	options.SetPropertyRenaming(PropertyRenamingPolicy_OFF)
	options.SetGeneratePseudoNames(true)

	src := `/** @unrestricted */
class __PRIVATE_WebChannelConnection extends class __PRIVATE_RestConnection {
  constructor(e) {
    this.databaseInfo = e, this.databaseId = e.databaseId;
  }
}
{
  constructor(e) {
    super(e), this.forceLongPolling = e.forceLongPolling,
              this.autoDetectLongPolling = e.autoDetectLongPolling,
              this.useFetchStreams = e.useFetchStreams,
              this.longPollingOptions = e.longPollingOptions;
    console.log('test');
  }
}
`
	expected := `const $$jscomp$classExtends$98447280$0$$ = class {};
class __PRIVATE_WebChannelConnection extends $$jscomp$classExtends$98447280$0$$ {
  constructor() {
    super();
    console.log("test");
  }
}
`

	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSTABLE)
	options.SetLanguageOut(CompilerOptions_LanguageMode_UNSTABLE)
	c.Test(options, src, expected)

	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSTABLE)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2019)
	c.Test(options, src, expected)
}

// Java: ES2022IntegrationTest.privateClassFields_supportedInChecksOnlyMode()
func (c *ES2022IntegrationTest) TestPrivateClassFields_supportedInChecksOnlyMode() {
	options := c.checksOnlyCompilerOptions()
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	c.TestNoWarnings(
		options,
		`class MyClass {
  #x = 2;
  #y;
  getX() {
    return this.#x;
  }
}
console.log(new MyClass().getX());
`)
}

// Java: ES2022IntegrationTest.privateClassFields_supportedInOptimizationsMode()
func (c *ES2022IntegrationTest) TestPrivateClassFields_supportedInOptimizationsMode() {
	options := c.CreateCompilerOptions()
	// Avoid prepending polyfill implementation code (e.g. WeakMap) into the expected test output.
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	options.SetVariableRenaming(VariableRenamingPolicy_OFF)
	// Declare $jscomp in externs to satisfy VarCheck when runtime libraries are not injected.
	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddConsole().
			AddExtra("/** @const */ var $jscomp = {};").
			BuildExternsFile("externs"),
	}

	src := `class MyClass {
  #x = 1;
  inc() {
    ++this.#x;
  }
  getX() {
    return this.#x;
  }
}

const o1 = new MyClass();
console.log(o1.getX());
o1.inc();
console.log(o1.getX());

const o2 = new MyClass();
o1.inc();
o1.inc();
console.log(o2.getX());
`

	// TODO(b/236744850): When languageOut is ECMASCRIPT_2022, private fields should be emitted
	// natively instead of being transpiled to $jscomp.PrivateMap (or modified to public fields).
	expected := `const $jscomp$privateMap$98447280$0 = new $jscomp.PrivateMap();
class MyClass {
  constructor() {
    var $jscomp$priv$98447280$1 = Object.create(null);
    $jscomp$privateMap$98447280$0.set(this, $jscomp$priv$98447280$1);
    $jscomp$priv$98447280$1.x = 1;
  }
  inc() {
    ++$jscomp$privateMap$98447280$0.get(this).x;
  }
  getX() {
    return $jscomp$privateMap$98447280$0.get(this).x;
  }
}
const o1 = new MyClass();
console.log(o1.getX());
o1.inc();
console.log(o1.getX());
const o2 = new MyClass();
o1.inc();
o1.inc();
console.log(o2.getX());
`

	// TODO(b/236744850): Update expected output for ES2022 once private fields are emitted
	// natively.
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2022)
	c.Test(options, src, expected)

	// Transpile ES2022 private class fields to ES2021-compatible JavaScript ($jscomp.PrivateMap).
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	c.Test(options, src, expected)
}

// Java: ES2022IntegrationTest.privateClassFields_nativeWithTranspiledPublicFields_knownBug()
func (c *ES2022IntegrationTest) TestPrivateClassFields_nativeWithTranspiledPublicFields_knownBug() {
	// Expected behavior: with public fields lowered and private fields kept native, field
	// initializers still run in declaration order. So `#b` reads `a` after `a = 1` has run, and
	// `new C().getB()` returns 1, as in the input.
	options := c.CreateCompilerOptions()
	options.SetTranspilePublicClassFields(true)
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}

	// TODO(b/236744850): Lowering only the public fields changes the initializer order.
	// Current behavior: `a = 1` becomes `this.a = 1` in the constructor, but `#b = this.a` stays a
	//   field, and fields are initialized before the constructor body runs.
	// Why it's a bug: `#b` reads `this.a` before `this.a = 1` runs, so `new C().getB()` returns
	//   undefined instead of 1.
	// In a full build: only when a tool calls `setTranspilePublicClassFields(true)`. No flag sets
	//   it, and no output language lowers public fields but keeps private ones.
	c.Test(
		options,
		`class C {
  a = 1;
  #b = this.a;
  getB() { return this.#b; }
}
console.log(new C().getB());
`,
		`class C {
  constructor() {
    this.a = 1;
  }
  #b = this.a;
  getB() { return this.#b; }
}
console.log(new C().getB());
`)
}

// Java: ES2022IntegrationTest.privateClassMembers_fullIntegrationTest()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_fullIntegrationTest() {
	options := c.CreateCompilerOptions()
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	options.SetVariableRenaming(VariableRenamingPolicy_OFF)
	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddConsole().
			AddExtra("/** @const */ var $jscomp = {};").
			BuildExternsFile("externs"),
	}

	src := `class C {
  x = 1;
  get #foo() {
    return this.x;
  }
  set #foo(x) {
    this.x = x;
  }
  #y = this.x + this.#foo;
  #z;
  #method(p) {
    return this.#foo + this.#y + p;
  }
  callAll() {
    return (
        this.#method(1) +
        ++this.#foo +
        C.#sMethod()
    );
  }
  static brandChecks(x) {
    if (#y in x) return 'C instance';
    if (#s2 in x) return 'C';
    return false;
  }
  canRewrite() {
    try {
      this.#foo = 1;
      return true;
    } catch (e) {
      return false;
    }
  }
  static s1 = 1;
  static #s2 = this.s1;
  static #sMethod() {
    return this.s1 + this.#s2;
  }
}
`

	expected := `const $jscomp$privateMap$98447280$0 = new $jscomp.PrivateMap();
const $jscomp$priv$proto$98447280$2 = Object.create(null, {
  foo: {
    get: function() {
      return this.$self.x;
    },
    set: function(x) {
      this.$self.x = x;
    }
  },
  method: {
    value: function(p) {
      return $jscomp$privateMap$98447280$0.get(this).foo + $jscomp$privateMap$98447280$0.get(this).y + p;
    }
  }
});
const $jscomp$staticPrivateMap$98447280$1 = new $jscomp.PrivateMap();
class C {
  constructor() {
    var $jscomp$priv$98447280$3 = Object.create($jscomp$priv$proto$98447280$2);
    $jscomp$priv$98447280$3.$self = this;
    $jscomp$privateMap$98447280$0.set(this, $jscomp$priv$98447280$3);
    this.x = 1;
    $jscomp$priv$98447280$3.y = this.x + $jscomp$privateMap$98447280$0.get(this).foo;
    $jscomp$priv$98447280$3.z = void 0;
  }
  callAll() {
    return $jscomp$privateMap$98447280$0.get(this).method.call(this, 1)
        + ++$jscomp$privateMap$98447280$0.get(this).foo
        + $jscomp$staticPrivateMap$98447280$1.get(C).sMethod.call(C);
  }
  static brandChecks(x) {
    if ($jscomp$privateMap$98447280$0.has(x)) return 'C instance';
    if ($jscomp$staticPrivateMap$98447280$1.has(x)) return 'C';
    return false;
  }
  canRewrite() {
    try {
      $jscomp$privateMap$98447280$0.get(this).foo = 1;
      return true;
    } catch (e) {
      return false;
    }
  }
  static $jscomp$staticInit$98447280$5() {
    var $jscomp$priv$98447280$4 = Object.create(null);
    $jscomp$staticPrivateMap$98447280$1.set(C, $jscomp$priv$98447280$4);
    $jscomp$priv$98447280$4.sMethod = function() {
      return this.s1 + $jscomp$staticPrivateMap$98447280$1.get(this).s2;
    };
    C.s1 = 1;
    $jscomp$priv$98447280$4.s2 = C.s1;
  }
}
C.$jscomp$staticInit$98447280$5();
`

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	c.Test(options, src, expected)
}

// Java: ES2022IntegrationTest.privateClassMembers_subclassAndCrossInstanceAccess()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_subclassAndCrossInstanceAccess() {
	options := c.CreateCompilerOptions()
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	options.SetVariableRenaming(VariableRenamingPolicy_OFF)
	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddConsole().
			AddExtra("/** @const */ var $jscomp = {};").
			BuildExternsFile("externs"),
	}

	src := `class Base {
  #secret = 42;
  #getSecret() { return this.#secret; }
  compareSecret(other) {
    return this.#getSecret() === other.#getSecret();
  }
}
class Derived extends Base {
  #childSecret = 100;
  getChildSecret() {
    return this.#childSecret;
  }
}
`

	expected := `const $jscomp$privateMap$98447280$0 = new $jscomp.PrivateMap();
const $jscomp$priv$proto$98447280$1 = Object.create(null, {
  getSecret: {
    value: function() {
      return $jscomp$privateMap$98447280$0.get(this).secret;
    }
  }
});
class Base {
  constructor() {
    var $jscomp$priv$98447280$2 = Object.create($jscomp$priv$proto$98447280$1);
    $jscomp$privateMap$98447280$0.set(this, $jscomp$priv$98447280$2);
    $jscomp$priv$98447280$2.secret = 42;
  }
  compareSecret(other) {
    return $jscomp$privateMap$98447280$0.get(this).getSecret.call(this)
        === $jscomp$privateMap$98447280$0.get(other).getSecret.call(other);
  }
}
const $jscomp$privateMap$98447280$3 = new $jscomp.PrivateMap();
class Derived extends Base {
  constructor() {
    super(...arguments);
    var $jscomp$priv$98447280$4 = Object.create(null);
    $jscomp$privateMap$98447280$3.set(this, $jscomp$priv$98447280$4);
    $jscomp$priv$98447280$4.childSecret = 100;
  }
  getChildSecret() {
    return $jscomp$privateMap$98447280$3.get(this).childSecret;
  }
}
`

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	c.Test(options, src, expected)
}

// Java: ES2022IntegrationTest.privateClassMembers_sideEffectingReceivers()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_sideEffectingReceivers() {
	options := c.CreateCompilerOptions()
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	options.SetVariableRenaming(VariableRenamingPolicy_OFF)
	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddConsole().
			AddExtra("/** @const */ var $jscomp = {};").
			BuildExternsFile("externs"),
	}

	src := `class Target {
  #data = 0;
  #update(val) {
    this.#data = val;
    return this.#data;
  }
  run(factory) {
    return factory().#update(10);
  }
}
`

	expected := `const $jscomp$privateMap$98447280$0 = new $jscomp.PrivateMap();
const $jscomp$priv$proto$98447280$2 = Object.create(null, {
  update: {
    value: function(val) {
      $jscomp$privateMap$98447280$0.get(this).data = val;
      return $jscomp$privateMap$98447280$0.get(this).data;
    }
  }
});
class Target {
  constructor() {
    var $jscomp$priv$98447280$3 = Object.create($jscomp$priv$proto$98447280$2);
    $jscomp$privateMap$98447280$0.set(this, $jscomp$priv$98447280$3);
    $jscomp$priv$98447280$3.data = 0;
  }
  run(factory) {
    var $jscomp$tmp$98447280$1;
    return ($jscomp$tmp$98447280$1 = factory(),
            $jscomp$privateMap$98447280$0.get($jscomp$tmp$98447280$1).update.call($jscomp$tmp$98447280$1, 10));
  }
}
`

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	c.Test(options, src, expected)
}

// Java: ES2022IntegrationTest.privateClassMembers_accessorCompoundAssignAndIncDec()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_accessorCompoundAssignAndIncDec() {
	options := c.CreateCompilerOptions()
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	options.SetVariableRenaming(VariableRenamingPolicy_OFF)
	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddConsole().
			AddExtra("/** @const */ var $jscomp = {};").
			BuildExternsFile("externs"),
	}

	src := `class Counter {
  #count = 0;
  get #val() { return this.#count; }
  set #val(v) { this.#count = v; }
  step() {
    this.#val += 5;
    ++this.#val;
    return this.#val++;
  }
}
`

	expected := `const $jscomp$privateMap$98447280$0 = new $jscomp.PrivateMap();
const $jscomp$priv$proto$98447280$1 = Object.create(null, {
  val: {
    get: function() {
      return $jscomp$privateMap$98447280$0.get(this.$self).count;
    },
    set: function(v) {
      $jscomp$privateMap$98447280$0.get(this.$self).count = v;
    }
  }
});
class Counter {
  constructor() {
    var $jscomp$priv$98447280$2 = Object.create($jscomp$priv$proto$98447280$1);
    $jscomp$priv$98447280$2.$self = this;
    $jscomp$privateMap$98447280$0.set(this, $jscomp$priv$98447280$2);
    $jscomp$priv$98447280$2.count = 0;
  }
  step() {
    $jscomp$privateMap$98447280$0.get(this).val += 5;
    ++$jscomp$privateMap$98447280$0.get(this).val;
    return $jscomp$privateMap$98447280$0.get(this).val++;
  }
}
`

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	c.Test(options, src, expected)
}

// Java: ES2022IntegrationTest.privateClassMembers_staticMethodAndFieldOrder()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_staticMethodAndFieldOrder() {
	options := c.CreateCompilerOptions()
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	options.SetVariableRenaming(VariableRenamingPolicy_OFF)
	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddConsole().
			AddExtra("/** @const */ var $jscomp = {};").
			BuildExternsFile("externs"),
	}

	src := `class C {
  static #a = 1;
  static #b() { return this.#a + 1; }
  static #c = this.#b() + 1;
  static getC() { return this.#c; }
}
`

	expected := `const $jscomp$staticPrivateMap$98447280$0 = new $jscomp.PrivateMap();
class C {
  static getC() {
    return $jscomp$staticPrivateMap$98447280$0.get(this).c;
  }
  static $jscomp$staticInit$98447280$2() {
    var $jscomp$priv$98447280$1 = Object.create(null);
    $jscomp$staticPrivateMap$98447280$0.set(C, $jscomp$priv$98447280$1);
    $jscomp$priv$98447280$1.b = function() {
      return $jscomp$staticPrivateMap$98447280$0.get(this).a + 1;
    };
    $jscomp$priv$98447280$1.a = 1;
    $jscomp$priv$98447280$1.c = $jscomp$staticPrivateMap$98447280$0.get(C).b.call(C) + 1;
  }
}
C.$jscomp$staticInit$98447280$2();
`

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	c.Test(options, src, expected)
}

// Java: ES2022IntegrationTest.privateClassMembers_superReadWithConstReflectExterns_knownBug()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_superReadWithConstReflectExterns_knownBug() {
	// Expected behavior: `super.p` in a private method is lowered to
	// `Reflect.get(Base.prototype, "p", this)`, with the `Reflect` name marked constant. The
	// externs declare `@const Reflect`, and the compiler requires every reference to a constant
	// name to be marked constant.
	// In a full build: only `DevMode.EVERY_PASS` checks this; the output is correct either way.
	//
	// TODO(b/236744850): `super.p` in a private method crashes the compiler with `@const Reflect`.
	// Current behavior: the `Reflect` name in the lowered `Reflect.get(...)` isn't marked constant,
	//   so the validity check after `es6NormalizeClasses` fails.
	// Why it's a bug: the compiler crashes on valid code when its validity checks are on, as in
	//   this test (`DevMode.EVERY_PASS`).
	options := c.CreateCompilerOptions()
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	options.SetVariableRenaming(VariableRenamingPolicy_OFF)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddReflect().
			AddExtra("/** @const */ var $jscomp = {};").
			BuildExternsFile("externs"),
	}
	src := `class Base {}
class Sub extends Base {
  #m() {
    return super.p;
  }
}
`

	ex := AssertThrows[jutil.IRuntimeException](func() { c.Compile(options, src) })
	AssertThat(ex.GetMessage()).Contains("Validity check failed for es6NormalizeClasses")
	AssertThat(ex.GetCause().GetCause().GetMessage()).
		Contains("Should be const but not marked as const")
}

// Java: ES2022IntegrationTest.privateClassMembers_staticPrivateAccessorWithPropertyRenaming_knownBug()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_staticPrivateAccessorWithPropertyRenaming_knownBug() {
	// Expected behavior: property renaming renames the key of a static private accessor and its
	// reads the same way, so `C.read()` still returns "sa". Renaming must not change what a
	// property read returns.
	//
	// TODO(b/236744850): Property renaming breaks reads of a static private accessor.
	// Current behavior: the getter is defined under the quoted key "sa", which property renaming
	//   doesn't rename, but the inlined read of `C.#sa` is renamed.
	// Why it's a bug: `console.log(C.read())` logs `undefined` instead of "sa".
	options := c.CreateCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	options.SetVariableRenaming(VariableRenamingPolicy_OFF)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddConsole().
			AddObject().
			AddExtra(
				"/** @const */ var $jscomp = {};",
				"/** @constructor */ $jscomp.PrivateMap = function() {};",
				"$jscomp.PrivateMap.prototype.get = function(o) {};",
				"$jscomp.PrivateMap.prototype.set = function(o, v) {};").
			BuildExternsFile("externs"),
	}

	c.Test(
		options,
		`class C {
  static get #sa() { return 'sa'; }
  static read() { return C.#sa; }
}
console.log(C.read());
`,
		`const $jscomp$staticPrivateMap$98447280$0 = new $jscomp.PrivateMap();
class C {}
(function() {
  var $jscomp$priv$98447280$1 = Object.create(null);
  $jscomp$staticPrivateMap$98447280$0.set(C, $jscomp$priv$98447280$1);
  Object.defineProperty($jscomp$priv$98447280$1, "sa", {
    get: function() {
      return "sa";
    }
  });
})();
console.log($jscomp$staticPrivateMap$98447280$0.get(C).a);
`)
}

// Java: ES2022IntegrationTest.privateClassMembers_optionalChainingIntegrationTest()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_optionalChainingIntegrationTest() {
	options := c.CreateCompilerOptions()
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	options.SetVariableRenaming(VariableRenamingPolicy_OFF)
	c.externs = []*SourceFile{
		NewTestExternsBuilder().
			AddConsole().
			AddExtra("/** @const */ var $jscomp = {};").
			BuildExternsFile("externs"),
	}

	src := `class Foo {
  #field = 100;
  #fnField = (x) => x * 2;
  #nullFn = null;
  #method(val) {
    return this.#field + val;
  }
  get #getter() {
    return this.#field * 2;
  }

  getField(target) {
    return target?.#field;
  }
  callMethod(target, val) {
    return target?.#method(val);
  }
  callFnField(target, val) {
    return target?.#fnField?.(val);
  }
  callNullFn(target, val) {
    return target?.#nullFn?.(val);
  }
  getGetter(target) {
    return target?.#getter;
  }

  static #staticField = 500;
  static #staticMethod(val) {
    return this.#staticField + val;
  }
  static get #staticGetter() {
    return this.#staticField * 3;
  }

  static getStaticField(target) {
    return target?.#staticField;
  }
  static callStaticMethod(target, val) {
    return target?.#staticMethod(val);
  }
  static getStaticGetter(target) {
    return target?.#staticGetter;
  }
}

const f = new Foo();
console.log(f.getField(f));
console.log(f.getField(null));
try {
  f.getField({});
} catch (e) {
  console.log(e instanceof TypeError);
}
console.log(f.callMethod(f, 50));
console.log(f.callMethod(null, 50));
try {
  f.callMethod({}, 50);
} catch (e) {
  console.log(e instanceof TypeError);
}
console.log(f.callFnField(f, 10));
console.log(f.callFnField(null, 10));
console.log(f.callNullFn(f, 10));
console.log(f.getGetter(f));
console.log(f.getGetter(null));
console.log(Foo.getStaticField(Foo));
console.log(Foo.getStaticField(null));
console.log(Foo.callStaticMethod(Foo, 100));
console.log(Foo.callStaticMethod(null, 100));
console.log(Foo.getStaticGetter(Foo));
console.log(Foo.getStaticGetter(null));
`

	expected := `const $jscomp$privateMap$98447280$0 = new $jscomp.PrivateMap();
const $jscomp$priv$proto$98447280$2 = Object.create(null, {
  method: {
    value: function(val) {
      return $jscomp$privateMap$98447280$0.get(this).field + val;
    }
  },
  getter: {
    get: function() {
      return $jscomp$privateMap$98447280$0.get(this.$self).field * 2;
    }
  }
});
const $jscomp$staticPrivateMap$98447280$1 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    var $jscomp$priv$98447280$3 = Object.create($jscomp$priv$proto$98447280$2);
    $jscomp$priv$98447280$3.$self = this;
    $jscomp$privateMap$98447280$0.set(this, $jscomp$priv$98447280$3);
    $jscomp$priv$98447280$3.field = 100;
    $jscomp$priv$98447280$3.fnField = x => x * 2;
    $jscomp$priv$98447280$3.nullFn = null;
  }
  getField(target) {
    return target == null ? void 0 : $jscomp$privateMap$98447280$0.get(target).field;
  }
  callMethod(target, val) {
    return target == null ? void 0 : $jscomp$privateMap$98447280$0.get(target).method.call(target, val);
  }
  callFnField(target, val) {
    return target == null ? void 0 : $jscomp$privateMap$98447280$0.get(target).fnField?.call(target, val);
  }
  callNullFn(target, val) {
    return target == null ? void 0 : $jscomp$privateMap$98447280$0.get(target).nullFn?.call(target, val);
  }
  getGetter(target) {
    return target == null ? void 0 : $jscomp$privateMap$98447280$0.get(target).getter;
  }
  static getStaticField(target) {
    return target == null ? void 0 : $jscomp$staticPrivateMap$98447280$1.get(target).staticField;
  }
  static callStaticMethod(target, val) {
    return target == null ? void 0 : $jscomp$staticPrivateMap$98447280$1.get(target).staticMethod.call(target, val);
  }
  static getStaticGetter(target) {
    return target == null ? void 0 : $jscomp$staticPrivateMap$98447280$1.get(target).staticGetter;
  }
  static $jscomp$staticInit$98447280$5() {
    var $jscomp$priv$98447280$4 = Object.create(null);
    $jscomp$priv$98447280$4.$self = Foo;
    $jscomp$staticPrivateMap$98447280$1.set(Foo, $jscomp$priv$98447280$4);
    $jscomp$priv$98447280$4.staticMethod = function(val) {
      return $jscomp$staticPrivateMap$98447280$1.get(this).staticField + val;
    };
    Object.defineProperty($jscomp$priv$98447280$4, "staticGetter", {
      get: function() {
        return $jscomp$staticPrivateMap$98447280$1.get(this.$self).staticField * 3;
      }
    });
    $jscomp$priv$98447280$4.staticField = 500;
  }
}
Foo.$jscomp$staticInit$98447280$5();
const f = new Foo();
console.log(f.getField(f));
console.log(f.getField(null));
try {
  f.getField({});
} catch (e) {
  console.log(e instanceof TypeError);
}
console.log(f.callMethod(f, 50));
console.log(f.callMethod(null, 50));
try {
  f.callMethod({}, 50);
} catch (e) {
  console.log(e instanceof TypeError);
}
console.log(f.callFnField(f, 10));
console.log(f.callFnField(null, 10));
console.log(f.callNullFn(f, 10));
console.log(f.getGetter(f));
console.log(f.getGetter(null));
console.log(Foo.getStaticField(Foo));
console.log(Foo.getStaticField(null));
console.log(Foo.callStaticMethod(Foo, 100));
console.log(Foo.callStaticMethod(null, 100));
console.log(Foo.getStaticGetter(Foo));
console.log(Foo.getStaticGetter(null));
`

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	c.Test(options, src, expected)
}

// Java: ES2022IntegrationTest.privateClassMembers_brandCheckWithTypeAndSymbolChecks_knownBug()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_brandCheckWithTypeAndSymbolChecks_knownBug() {
	// Expected behavior: no diagnostics with `setCheckSymbols(true)` or `setCheckTypes(true)`. In
	// `#x in this`, `#x` is a private name of `C`, not a variable, and a brand check is valid on
	// any object, including a class instance, which the type checker treats as a struct.
	js := "class C { #x; has() { return #x in this; } }"
	checkSymbols := c.CreateCompilerOptions()
	checkSymbols.SetCheckSymbols(true)
	checkTypes := c.CreateCompilerOptions()
	checkTypes.SetCheckTypes(true)
	// TODO(b/236744850): `VarCheck` reports the private name `#x` as an undeclared variable.
	// Current behavior: `VarCheck` reports the error `JSC_UNDEFINED_VARIABLE` ("variable #x is
	//   undeclared").
	// Why it's a bug: the error is false, because `#x` is a private name declared in `C`, not a
	//   variable. Valid code fails to compile.
	c.TestWarning(checkSymbols, js, DiagnosticGroups_UNDEFINED_VARIABLES)
	// TODO(b/236744850): `TypeCheck` reports a brand check on a struct.
	// Current behavior: `TypeCheck` reports the warning `JSC_IN_USED_WITH_STRUCT` ("Cannot use
	//   the IN operator with structs"), because class instances are structs.
	// Why it's a bug: the warning is false, because `#x in this` only checks whether `this` has
	//   the private field `#x`, which is fine on a struct. In strict builds, the warning is an
	//   error, so valid code fails to compile.
	c.TestWarning(checkTypes, js, DiagnosticGroups_CHECK_TYPES)
}

// Java: ES2022IntegrationTest.privateClassMembers_weakMapPolyfillInjection_wontFix()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_weakMapPolyfillInjection_wontFix() {
	// Won't fix (b/236744850): The `WeakMap` polyfill is injected for output with native classes.
	// Native behavior: every engine with native classes also has a native `WeakMap`, which backs
	// `$jscomp.PrivateMap`, so only output without native classes needs the polyfill.
	options := c.CreateCompilerOptions()
	options.SetRewritePolyfills(true)
	src := "class C { #x = 1; }"

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	AssertThat(c.Compile(options, src).GetRuntimeJsLibManager().GetInjectedLibraries()).
		ContainsAtLeast("es6/privatemap", "es6/weakmap")

	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2021)
	injected := c.Compile(options, src).GetRuntimeJsLibManager().GetInjectedLibraries()
	AssertThat(injected).Contains("es6/privatemap")
	// Current behavior: `privatemap.js` requires `es6/weakmap`, so the polyfill is injected
	//   whenever `es6/privatemap` is, including for this ES2021 output.
	// Why it's a bug: ES2021 output only runs on engines with a native `WeakMap`, so the polyfill
	//   just makes the output bigger.
	// Why we won't fix it: it only costs size (about 2 KB minified), not behavior. Without the
	//   `require`, the compiler would need its own code to inject the polyfill for ES5 output;
	//   if it does so only when `rewritePolyfills` is on, ES5 builds with it off would have to
	//   provide `WeakMap` themselves.
	AssertThat(injected).Contains("es6/weakmap")
}

// Java: ES2022IntegrationTest.privateClassMembers_nativeOutputWithAdvancedOptimizations_knownBug()
func (c *ES2022IntegrationTest) TestPrivateClassMembers_nativeOutputWithAdvancedOptimizations_knownBug() {
	// Expected behavior: when private elements are not lowered, ADVANCED optimizations keep the
	// private members in `C` and rename them to other private names. Private names are valid only
	// inside the body of their class.
	options := c.CreateCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	c.externs = []*SourceFile{NewTestExternsBuilder().AddConsole().BuildExternsFile("externs")}
	js := `class C {
  static #s = 1;
  #x = 2;
  #m() { return C.#s; }
  run() { return this.#x + this.#m(); }
}
console.log(new C().run());
`
	expected := `class a {
  static #b = 1;
  #c = 2;
  #a() { return a.#b; }
  d() { return this.#c + this.#a(); }
}
console.log(new a().d());
`
	// TODO(b/236744850): ADVANCED optimizations crash on a static private field.
	// Current behavior: `CollapseProperties` tries to move `static #s = 1` out of `C` into a
	//   global variable, and fails a `checkState` because `1` is not a function.
	//   `test(options, js, expected)` throws the `IllegalStateException`.
	// Why it's a bug: the compiler crashes on valid code.
	var recovered any
	var stack string
	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = r
				stack = string(debug.Stack())
			}
		}()
		c.Test(options, js, expected)
	}()
	AssertWithMessage("Expected IllegalStateException").That(recovered != nil).IsTrue()
	IsInstanceOf[jutil.IIllegalStateException](AssertThat(recovered))
	// The failed `checkState` has no message, so check that it is the one in
	// `CheckForReceiverAffectedByCollapse`.
	AssertThat(strings.ToLower(stack)).Contains("checkforreceiveraffectedbycollapse")
}

func (c *ES2022IntegrationTest) es2021OptimizedWithoutTranspilationOptions() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetDevMode(CompilerOptions_DevMode_EVERY_PASS)
	options.SetCodingConvention(NewGoogleCodingConvention())
	WarningLevel_VERBOSE.SetOptionsForWarningLevel(options)
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetTypeBasedOptimizationOptions(options)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	return options
}

func (c *ES2022IntegrationTest) es2021FullyOptimizedOptions() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetDevMode(CompilerOptions_DevMode_EVERY_PASS)
	options.SetCodingConvention(NewGoogleCodingConvention())
	WarningLevel_VERBOSE.SetOptionsForWarningLevel(options)
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetTypeBasedOptimizationOptions(options)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	return options
}

// Java: ES2021IntegrationTest.logicalAssignmentStaticClassPropertyNotTranspiledOutput()
func (c *ES2022IntegrationTest) TestLogicalAssignmentStaticClassPropertyNotTranspiledOutput() {
	options := c.es2021OptimizedWithoutTranspilationOptions()

	c.Test(
		options,
		`class MyClass {}
MyClass.myField = 1;
alert(() => {
  return (MyClass.myField ??= {});
});
alert(MyClass.myField);
`,
		`var a = 1;
alert(() => a ??= {});
alert(a);
`)
}

// Java: ES2021IntegrationTest.logicalAssignmentStaticClassPropertyTranspiledOutput()
func (c *ES2022IntegrationTest) TestLogicalAssignmentStaticClassPropertyTranspiledOutput() {
	options := c.es2021FullyOptimizedOptions()

	c.Test(
		options,
		`class MyClass {}
MyClass.myField = 1;
alert(() => {
  return (MyClass.myField ??= {});
});
alert(MyClass.myField);
`,
		`var a = 1;
alert(() => a ?? (a = {}));
alert(a);
`)
}

// Java: ES2021IntegrationTest.logicalAssignmentNamespacedStaticClassPropertyTranspiledOutput()
func (c *ES2022IntegrationTest) TestLogicalAssignmentNamespacedStaticClassPropertyTranspiledOutput() {
	options := c.es2021FullyOptimizedOptions()

	c.Test(
		options,
		`const MyNs = {};
MyNs.MyClass = class {};
MyNs.MyClass.myField = 1;
alert(() => {
  return (MyNs.MyClass.myField ??= {});
});
alert(MyNs.MyClass.myField);
`,
		`var a = 1;
alert(() => a ?? (a = {}));
alert(a);
`)
}

// Java: ES2021IntegrationTest.logicalAssignmentStaticClassPropertyGlobalConditionalOutput_knownBug()
func (c *ES2022IntegrationTest) TestLogicalAssignmentStaticClassPropertyGlobalConditionalOutput_knownBug() {
	options := c.es2021FullyOptimizedOptions()

	// TODO: b/476221453 - `(c = a).a ?? (c.a = 3)` should be collapsed, preventing `MyClass.b`
	// from being constant-inlined as `2` in `alert(2)`.
	c.Test(
		options,
		`class MyClass {}
MyClass.a = 1;
MyClass.b = 2;
alert(MyClass.a ||= (MyClass.b ??= 3));
alert(MyClass.b);
`,
		`class a {}
var b = 1;
let c;
alert(b || (b = (c = a).a ?? (c.a = 3)));
alert(2);
`)
}

// Java: ES2021IntegrationTest.logicalAssignmentStaticClassPropertyInLoopNotTranspiledOutput()
func (c *ES2022IntegrationTest) TestLogicalAssignmentStaticClassPropertyInLoopNotTranspiledOutput() {
	options := c.es2021OptimizedWithoutTranspilationOptions()

	c.Test(
		options,
		`class MyClass {}
MyClass.myField = 1;
alert(() => {
  while (MyClass.myField ??= {}) {
    break;
  }
});
alert(MyClass.myField);
`,
		`var a = 1;
alert(() => {
  for (; a ??= {};) {
    break;
  }
});
alert(a);
`)
}

// Java: ES2021IntegrationTest.logicalAssignmentStaticClassPropertyInForDoWhileForOfTranspiledOutput()
func (c *ES2022IntegrationTest) TestLogicalAssignmentStaticClassPropertyInForDoWhileForOfTranspiledOutput() {
	options := c.es2021FullyOptimizedOptions()

	c.Test(
		options,
		`class MyClass {}
MyClass.a = 1;
MyClass.b = 2;
MyClass.items = [1];
alert(() => {
  for (; MyClass.a ??= 10; MyClass.b ||= 20) {
    if (window.stop) break;
  }
  do {
    if (window.stop) break;
  } while (MyClass.a &&= 30);
  for (const x of (MyClass.items ??= [])) {
    alert(x);
  }
});
alert([MyClass.a, MyClass.b, MyClass.items]);
`,
		`var b = 1, c = 2, d = [1];
alert(() => {
  for (; (b ?? (b = 10)) && !window.a; c || (c = 20)) {}
  do {
    if (window.a) break;
  } while (b && (b = 30));
  for (let e of d ?? (d = [])) {
    alert(e);
  }
});
alert([b, c, d]);
`)
}

func (c *ES2022IntegrationTest) integrationTestCreateCompilerOptions() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	options.SetDevMode(CompilerOptions_DevMode_EVERY_PASS)
	options.SetCodingConvention(NewGoogleCodingConvention())
	options.renamePrefixNamespaceAssumeCrossChunkNames = true
	options.SetAssumeGettersArePure(false)
	options.SetPrettyPrint(true)
	return options
}

// Java: IntegrationTest.testExtractPrototypeMemberDeclarationsAndAmbiguateProperties()
func (c *ES2022IntegrationTest) TestExtractPrototypeMemberDeclarationsAndAmbiguateProperties() {
	options := c.integrationTestCreateCompilerOptions()
	options.SetCheckTypes(true)
	options.SetExtractPrototypeMemberDeclarations(true)
	options.SetAmbiguateProperties(true)
	options.SetPropertyRenaming(PropertyRenamingPolicy_ALL_UNQUOTED)
	options.SetVariableRenaming(VariableRenamingPolicy_ALL)

	var code strings.Builder
	code.WriteString("/** @constructor */ function Foo() {}\n")
	for i := 0; i < 7; i++ {
		fmt.Fprintf(&code, "Foo.prototype.fooM%d = function() {};\n", i)
	}
	code.WriteString("/** @constructor */ function Bar() {}\n")
	for i := 0; i < 7; i++ {
		fmt.Fprintf(&code, "Bar.prototype.barM%d = function() {};\n", i)
	}

	// TODO(b/575626948): AmbiguateProperties should ambiguate Foo and Bar prototype methods to the
	// same property names (a..g) instead of distinct names (a..g and h..n).
	var expected strings.Builder
	expected.WriteString("var a; function b() {} a = b.prototype;")
	for _, m := range []string{"h", "i", "j", "k", "l", "m", "n"} {
		fmt.Fprintf(&expected, "a.%s=function(){};", m)
	}
	expected.WriteString("function c() {} a = c.prototype;")
	for _, m := range []string{"a", "b", "c", "d", "e", "f", "g"} {
		fmt.Fprintf(&expected, "a.%s=function(){};", m)
	}

	c.Test(options, code.String(), expected.String())
}

// Java: IntegrationTest.testGenerateIjsWithHoistedExternsInChunks()
func (c *ES2022IntegrationTest) TestGenerateIjsWithHoistedExternsInChunks() {
	options := c.integrationTestCreateCompilerOptions()
	options.SetIncrementalChecks(CompilerOptions_IncrementalCheckMode_GENERATE_IJS)

	compiler := c.CompileSources(
		options,
		[]string{
			"/** @externs */ /** @constructor */ function MyExternType() {}",
			"/** @type {!MyExternType} */ var x = new MyExternType();",
		})
	c.CheckUnexpectedErrorsOrWarnings(compiler, 0)

	chunks := compiler.GetChunks()
	AssertThat(compiler.ToSourceOfJSChunk(chunks[0])).
		IsEqualTo(
			"/** @fileoverview @typeSummary */\n" +
				"/** @externs */ \n" +
				"/** @constructor */ \n" +
				"function MyExternType() {\n" +
				"}\n" +
				";")
	AssertThat(compiler.ToSourceOfJSChunk(chunks[1])).
		IsEqualTo(
			"/** @fileoverview @typeSummary */\n" +
				"/** @type {!MyExternType} */ \n" +
				"var x;\n")
}

// Java: IntegrationTest.testUnusedTestOnlyOverride_compileToZeroDependsOnDeclarationOrder()
func (c *ES2022IntegrationTest) TestUnusedTestOnlyOverride_compileToZeroDependsOnDeclarationOrder() {
	options := c.integrationTestCreateCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetDefineToBooleanLiteral("DEBUG", false)

	// When `let debugOverride;` precedes `isDebugMode()` and `assertInRealm()`, RemoveUnusedCode
	// strips `TEST_ONLY`, InlineVariables folds `debugOverride` to `void 0`, `assertInRealm`
	// simplifies to an empty function, and the entire bundle compiles to zero bytes.
	c.Test(
		options,
		`/** @define {boolean} */
const DEBUG = true;
let debugOverride;
function isDebugMode() {
  return debugOverride ?? DEBUG;
}
function assertInRealm(name) {
  if (isDebugMode() && !window['inRealm']) {
    throw new Error(`+"`Assertion failed for ${name}: must run inside child realm.`"+`);
  }
}
const TEST_ONLY = {
  setDebug(val) {
    debugOverride = val;
  },
};
assertInRealm('ext-a');
assertInRealm('ext-b');
`,
		"")

	// TODO(b/581117790): Moving `let debugOverride;` below `assertInRealm()` (above `TEST_ONLY`)
	// deadlocks InlineFunctions and InlineVariables even though `TEST_ONLY` is completely stripped
	// by RemoveUnusedCode and `debugOverride` has zero remaining assignments: multi-call
	// `assertInRealm` is too large for InlineFunctions to inline before its body shrinks, while
	// InlineVariables refuses to fold `debugOverride` to `void 0` because ReferenceCollector
	// visits `assertInRealm` (declared above `let debugOverride;`) before `let debugOverride;`.
	c.Test(
		options,
		`/** @define {boolean} */
const DEBUG = true;
function isDebugMode() {
  return debugOverride ?? DEBUG;
}
function assertInRealm(name) {
  if (isDebugMode() && !window['inRealm']) {
    throw new Error(`+"`Assertion failed for ${name}: must run inside child realm.`"+`);
  }
}
let debugOverride;
const TEST_ONLY = {
  setDebug(val) {
    debugOverride = val;
  },
};
assertInRealm('ext-a');
assertInRealm('ext-b');
`,
		`function a(b) {
  if ((c ?? !1) && !window.inRealm) {
    throw Error(`+"`Assertion failed for ${b}: must run inside child realm.`"+`);
  }
}
let c;
a("ext-a");
a("ext-b");
`)
}

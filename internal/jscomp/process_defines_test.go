// Port of com/google/javascript/jscomp/ProcessDefinesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ProcessDefinesTest.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type processDefinesTest_override struct {
	name  string
	value rhino.Node
}

type processDefinesTest_overrides []processDefinesTest_override

func (o *processDefinesTest_overrides) Put(name string, value rhino.Node) {
	*o = append(*o, processDefinesTest_override{name: name, value: value})
}

func (o *processDefinesTest_overrides) Clear() {
	*o = (*o)[:0]
}

type ProcessDefinesTest struct {
	CompilerTestCaseBase

	overrides               processDefinesTest_overrides
	namespace               *GlobalNamespace
	mode                    ProcessDefines_Mode
	recognizeClosureDefines bool

	enableJ2clPasses                   bool
	enableDefineWithoutGoogDefineCheck bool

	outputFeatureSet *parser.FeatureSet
}

// Java: ProcessDefinesTest()
func NewProcessDefinesTest() *ProcessDefinesTest {
	c := &ProcessDefinesTest{
		recognizeClosureDefines: true,
		enableJ2clPasses:        false,
		outputFeatureSet:        parser.FeatureSet_ES2017,
	}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_DEFAULT_EXTERNS+"var externMethod;")
	return c
}

func TestProcessDefinesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewProcessDefinesTest() })
}

// Java: ProcessDefinesTest.setUp()
func (c *ProcessDefinesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.overrides.Clear()
	c.mode = ProcessDefines_Mode_CHECK_AND_OPTIMIZE
	c.enableDefineWithoutGoogDefineCheck = false

	// ProcessDefines emits warnings if the user tries to re-define a constant,
	// but the constant is not defined anywhere in the binary.
	c.AllowSourcelessWarnings()
}

// Java: ProcessDefinesTest.getProcessor(Compiler)
func (c *ProcessDefinesTest) GetProcessor(compiler *Compiler) CompilerPass {
	if c.enableJ2clPasses {
		J2clSourceFileChecker_MarkToRunJ2clPasses(compiler)
	}
	return c.NewProcessDefinesWithInjectedNamespace(compiler)
}

var processDefinesTest_ZONE_INPUT_PATTERN = jutil.Pattern_Compile(".*/packages/zone.js/.*\\.js")

// Java: ProcessDefinesTest.getOptions()
func (c *ProcessDefinesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()

	options.SetOutputFeatureSet(c.outputFeatureSet)
	options.SetEnableZonesDefineName(jutil.StringPtr("javascript.angular2.ENABLE_ZONES"))
	options.SetZoneInputPattern(processDefinesTest_ZONE_INPUT_PATTERN)
	if c.enableDefineWithoutGoogDefineCheck {
		options.SetWarningLevel(DiagnosticGroups_DEFINE_WITHOUT_GOOG_DEFINE, CheckLevel_ERROR)
	}
	return options
}

// testWithPrefix is a helper for tests that expects definitions to remain
// unchanged, such that definitions+js is converted to definitions+expected.
//
// Java: ProcessDefinesTest.testWithPrefix(String, String, String)
func (c *ProcessDefinesTest) testWithPrefix(definitions string, js string, expected string) {
	c.Test(definitions+js, definitions+expected)
}

// ProcessDefinesTest_ProcessDefinesWithInjectedNamespace is
// ProcessDefinesTest.ProcessDefinesWithInjectedNamespace (an inner class).
type ProcessDefinesTest_ProcessDefinesWithInjectedNamespace struct {
	outer    *ProcessDefinesTest
	compiler *Compiler
}

// Java: ProcessDefinesTest.ProcessDefinesWithInjectedNamespace(Compiler)
func (c *ProcessDefinesTest) NewProcessDefinesWithInjectedNamespace(compiler *Compiler) *ProcessDefinesTest_ProcessDefinesWithInjectedNamespace {
	return &ProcessDefinesTest_ProcessDefinesWithInjectedNamespace{outer: c, compiler: compiler}
}

// Java: ProcessDefinesTest.ProcessDefinesWithInjectedNamespace.process(Node, Node)
func (p *ProcessDefinesTest_ProcessDefinesWithInjectedNamespace) Process(externs rhino.Node, js rhino.Node) {
	c := p.outer
	c.namespace = NewGlobalNamespaceWithRoots(p.compiler, externs, js)
	var replacements OrderedStringMap[rhino.Node]
	for _, ov := range c.overrides {
		replacements.Put(ov.name, ov.value)
	}
	NewProcessDefines(p.compiler, ProcessDefinesOptions{
		Replacements:            replacements,
		Mode:                    c.mode,
		NamespaceSupplier:       func() *GlobalNamespace { return c.namespace },
		RecognizeClosureDefines: c.recognizeClosureDefines,
		EnableZonesDefineName:   p.compiler.GetOptions().GetEnableZonesDefineName(),
		ZoneInputPattern:        p.compiler.GetOptions().GetZoneInputPattern(),
	}).Process(externs, js)
}

type ProcessDefines_Builder struct {
	compiler AbstractCompiler
	opts     ProcessDefinesOptions
}

func NewProcessDefines_Builder(compiler AbstractCompiler) *ProcessDefines_Builder {
	return &ProcessDefines_Builder{
		compiler: compiler,
		opts:     ProcessDefinesOptions{RecognizeClosureDefines: true},
	}
}

func (b *ProcessDefines_Builder) PutReplacement(name string, value rhino.Node) *ProcessDefines_Builder {
	b.opts.Replacements.Put(name, value)
	return b
}

func (b *ProcessDefines_Builder) PutReplacements(replacementValues jutil.MapSource[string, rhino.Node]) *ProcessDefines_Builder {
	for k, v := range replacementValues.All() {
		b.opts.Replacements.Put(k, v)
	}
	return b
}

func (b *ProcessDefines_Builder) SetMode(x ProcessDefines_Mode) *ProcessDefines_Builder {
	b.opts.Mode = x
	return b
}

func (b *ProcessDefines_Builder) InjectNamespace(namespaceSupplier func() *GlobalNamespace) *ProcessDefines_Builder {
	b.opts.NamespaceSupplier = namespaceSupplier
	return b
}

func (b *ProcessDefines_Builder) SetUnknownDefinesToIgnore(unknownDefinesToIgnore []string) *ProcessDefines_Builder {
	b.opts.UnknownDefinesToIgnore = unknownDefinesToIgnore
	return b
}

func (b *ProcessDefines_Builder) SetRecognizeClosureDefines(recognizeClosureDefines bool) *ProcessDefines_Builder {
	b.opts.RecognizeClosureDefines = recognizeClosureDefines
	return b
}

func (b *ProcessDefines_Builder) SetEnableZonesDefineName(enableZonesDefineName *string) *ProcessDefines_Builder {
	b.opts.EnableZonesDefineName = enableZonesDefineName
	return b
}

func (b *ProcessDefines_Builder) SetZoneInputPattern(zoneInputPattern *jutil.Pattern) *ProcessDefines_Builder {
	b.opts.ZoneInputPattern = zoneInputPattern
	return b
}

func (b *ProcessDefines_Builder) Build() *ProcessDefines {
	var copied OrderedStringMap[rhino.Node]
	for _, e := range b.opts.Replacements.Entries() {
		copied.Put(e.Key, e.Value)
	}
	opts := b.opts
	opts.Replacements = copied
	return NewProcessDefines(b.compiler, opts)
}

// Java: ProcessDefinesTest.testBasicDefine1()
func (c *ProcessDefinesTest) TestBasicDefine1() {
	c.Test("/** @define {boolean} */ var DEF = true", "/** @define {boolean} */ var DEF=true")
}

// Java: ProcessDefinesTest.testBasicDefine2()
func (c *ProcessDefinesTest) TestBasicDefine2() {
	c.Test("/** @define {string} */ var DEF = 'a'", "/** @define {string} */ var DEF=\"a\"")
}

// Java: ProcessDefinesTest.testBasicDefine3()
func (c *ProcessDefinesTest) TestBasicDefine3() {
	c.Test("/** @define {number} */ var DEF = 0", "/** @define {number} */ var DEF=0")
}

// Java: ProcessDefinesTest.testDefineBadType()
func (c *ProcessDefinesTest) TestDefineBadType() {
	c.Test(
		c.Srcs("/** @define {Object} */ var DEF = {}"),
		c.Error(ProcessDefines_INVALID_DEFINE_TYPE),
		c.Error(ProcessDefines_INVALID_DEFINE_VALUE))
}

// Java: ProcessDefinesTest.testDefineBadType_nullishCoalesce()
func (c *ProcessDefinesTest) TestDefineBadType_nullishCoalesce() {
	c.Test(c.Srcs("/** @define {string} */ var DEF = 'a' ?? 'b'"))
	c.TestError("/** @define {string} */ var DEF = 'a' ?? null", ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testChecksOnlyProducesErrors()
func (c *ProcessDefinesTest) TestChecksOnlyProducesErrors() {
	c.mode = ProcessDefines_Mode_CHECK
	c.Test(
		c.Srcs("/** @define {Object} */ var DEF = {}"),
		c.Error(ProcessDefines_INVALID_DEFINE_TYPE),
		c.Error(ProcessDefines_INVALID_DEFINE_VALUE))
}

// Java: ProcessDefinesTest.testUnknownDefineWarning()
func (c *ProcessDefinesTest) TestUnknownDefineWarning() {
	c.mode = ProcessDefines_Mode_OPTIMIZE
	c.overrides.Put("a.B", rhino.NewNode(rhino.Token_TRUE))
	c.Test("var a = {};", "var a = {};", c.Warning(ProcessDefines_UNKNOWN_DEFINE_WARNING))
}

// Java: ProcessDefinesTest.testDefineWithBadValue1()
func (c *ProcessDefinesTest) TestDefineWithBadValue1() {
	c.TestError(
		"/** @define {boolean} */ var DEF = new Boolean(true);",
		ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testDefineWithBadValue2()
func (c *ProcessDefinesTest) TestDefineWithBadValue2() {
	c.TestError("/** @define {string} */ var DEF = 'x' + y;", ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testDefineWithBadValue3()
func (c *ProcessDefinesTest) TestDefineWithBadValue3() {
	// alias is not const
	c.TestError(
		"let x = 'x'; /** @define {string} */ var DEF = x;", ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testDefineWithBadValue4()
func (c *ProcessDefinesTest) TestDefineWithBadValue4() {
	c.TestError("/** @define {string} */ var DEF = null;", ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testDefineWithBadValue5()
func (c *ProcessDefinesTest) TestDefineWithBadValue5() {
	c.TestError("/** @define {string} */ var DEF = undefined;", ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testDefineWithBadValue6()
func (c *ProcessDefinesTest) TestDefineWithBadValue6() {
	c.TestError("/** @define {string} */ var DEF = NaN;", ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testGoogDefineWithBadValue()
func (c *ProcessDefinesTest) TestGoogDefineWithBadValue() {
	c.TestError(
		"/** @define {boolean} */ var DEF = goog.define('DEF_XYZ', new Boolean(true));",
		ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testDefineWithLet()
func (c *ProcessDefinesTest) TestDefineWithLet() {
	c.TestError(
		"/** @define {boolean} */ let DEF = new Boolean(true);",
		ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testDefineInExterns()
func (c *ProcessDefinesTest) TestDefineInExterns() {
	c.TestSame(c.Externs(CompilerTestCase_DEFAULT_EXTERNS+"/** @define {boolean} */ var EXTERN_DEF;"), c.Srcs(""))
}

// Java: ProcessDefinesTest.testDefineInExternsPlusUsage()
func (c *ProcessDefinesTest) TestDefineInExternsPlusUsage() {
	c.TestSame(
		c.Externs(CompilerTestCase_DEFAULT_EXTERNS+"/** @define {boolean} */ var EXTERN_DEF;"),
		c.Srcs("/** @define {boolean} */ var DEF = EXTERN_DEF"))
}

// Java: ProcessDefinesTest.testNonDefineInExternsPlusUsage()
func (c *ProcessDefinesTest) TestNonDefineInExternsPlusUsage() {
	c.TestError(
		c.Externs(CompilerTestCase_DEFAULT_EXTERNS+"/** @const {boolean} */ var EXTERN_NON_DEF;"),
		c.Srcs("/** @define {boolean} */ var DEF = EXTERN_NON_DEF"),
		c.Error(ProcessDefines_INVALID_DEFINE_VALUE))
}

// Java: ProcessDefinesTest.testDefineCompiledInExterns()
func (c *ProcessDefinesTest) TestDefineCompiledInExterns() {
	c.TestSame(c.Externs(CompilerTestCase_DEFAULT_EXTERNS+"/** @define {boolean} */ var COMPILED;"), c.Srcs(""))
}

// Java: ProcessDefinesTest.testDefineWithDependentValue()
func (c *ProcessDefinesTest) TestDefineWithDependentValue() {
	c.Test(
		`/** @define {boolean} */ var BASE = false;
/** @define {boolean} */ var DEF = !BASE;
`,
		`/** @define {boolean} */ var BASE = false;
/** @define {boolean} */ var DEF = !BASE
`)
	c.Test(
		`var a = {};
/** @define {boolean} */ a.BASE = false;
/** @define {boolean} */ a.DEF = !a.BASE;
`,
		`var a={};
/** @define {boolean} */ a.BASE = false;
/** @define {boolean} */ a.DEF = !a.BASE
`)
}

// Java: ProcessDefinesTest.testDefineWithInvalidDependentValue()
func (c *ProcessDefinesTest) TestDefineWithInvalidDependentValue() {
	c.TestError(
		`var BASE = false;
/** @define {boolean} */ var DEF = !BASE;
`,
		ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testOverriding1()
func (c *ProcessDefinesTest) TestOverriding1() {
	c.overrides.Put("DEF_OVERRIDE_TO_TRUE", rhino.NewNode(rhino.Token_TRUE))
	c.overrides.Put("DEF_OVERRIDE_TO_FALSE", rhino.NewNode(rhino.Token_FALSE))
	c.Test(
		`/** @define {boolean} */ var DEF_OVERRIDE_TO_TRUE = false;
/** @define {boolean} */ var DEF_OVERRIDE_TO_FALSE = true
`,
		`/** @define {boolean} */ var DEF_OVERRIDE_TO_TRUE = true;
/** @define {boolean} */var DEF_OVERRIDE_TO_FALSE=false
`)
}

// Java: ProcessDefinesTest.testOverriding2()
func (c *ProcessDefinesTest) TestOverriding2() {
	c.overrides.Put("DEF_OVERRIDE_TO_TRUE", rhino.NewNode(rhino.Token_TRUE))
	normalConst := "var DEF_OVERRIDE_TO_FALSE=true;"
	c.testWithPrefix(
		normalConst,
		"/** @define {boolean} */ var DEF_OVERRIDE_TO_TRUE = false",
		"/** @define {boolean} */ var DEF_OVERRIDE_TO_TRUE = true")
}

// Java: ProcessDefinesTest.testOverriding3()
func (c *ProcessDefinesTest) TestOverriding3() {
	c.overrides.Put("DEF_OVERRIDE_TO_TRUE", rhino.NewNode(rhino.Token_TRUE))
	c.Test(
		"/** @define {boolean} */ var DEF_OVERRIDE_TO_TRUE = true;",
		"/** @define {boolean} */ var DEF_OVERRIDE_TO_TRUE = true")
}

// Java: ProcessDefinesTest.testOverridingString0()
func (c *ProcessDefinesTest) TestOverridingString0() {
	c.Test(
		"/** @define {string} */ var DEF_OVERRIDE_STRING = 'x';",
		"/** @define {string} */ var DEF_OVERRIDE_STRING=\"x\"")
}

// Java: ProcessDefinesTest.testOverridingString1()
func (c *ProcessDefinesTest) TestOverridingString1() {
	c.Test(
		"/** @define {string} */ var DEF_OVERRIDE_STRING = 'x' + 'y';",
		"/** @define {string} */ var DEF_OVERRIDE_STRING=\"x\" + \"y\"")
}

// Java: ProcessDefinesTest.testOverridingString2()
func (c *ProcessDefinesTest) TestOverridingString2() {
	c.overrides.Put("DEF_OVERRIDE_STRING", rhino.Node_NewString("foo"))
	c.Test(
		"/** @define {string} */ var DEF_OVERRIDE_STRING = 'x';",
		"/** @define {string} */ var DEF_OVERRIDE_STRING=\"foo\"")
}

// Java: ProcessDefinesTest.testOverridingString3()
func (c *ProcessDefinesTest) TestOverridingString3() {
	c.overrides.Put("DEF_OVERRIDE_STRING", rhino.Node_NewString("foo"))
	c.Test(
		"/** @define {string} */ var DEF_OVERRIDE_STRING = 'x' + 'y';",
		"/** @define {string} */ var DEF_OVERRIDE_STRING=\"foo\"")
}

// Java: ProcessDefinesTest.testMisspelledOverride()
func (c *ProcessDefinesTest) TestMisspelledOverride() {
	c.overrides.Put("DEF_BAD_OVERIDE", rhino.NewNode(rhino.Token_TRUE)) // NOTYPO: Intentional misspelling.
	c.Test(
		"/** @define {boolean} */ var DEF_BAD_OVERRIDE = true",
		"/** @define {boolean} */ var DEF_BAD_OVERRIDE = true",
		c.Warning(ProcessDefines_UNKNOWN_DEFINE_WARNING))
}

// Java: ProcessDefinesTest.testCompiledIsKnownDefine()
func (c *ProcessDefinesTest) TestCompiledIsKnownDefine() {
	c.overrides.Put("COMPILED", rhino.NewNode(rhino.Token_TRUE))
	c.TestSame("")
}

// Java: ProcessDefinesTest.testSimpleReassign1()
func (c *ProcessDefinesTest) TestSimpleReassign1() {
	c.Test(
		c.Srcs(
			`/** @define {boolean} */ var DEF = false;
DEF = true;
`),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testSimpleReassign2()
func (c *ProcessDefinesTest) TestSimpleReassign2() {
	c.Test(
		c.Srcs(
			`/** @define {number|boolean} */ var DEF=false;
DEF=true;
DEF=3
`),
		c.Error(ProcessDefines_NON_CONST_DEFINE),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testSimpleReassign3()
func (c *ProcessDefinesTest) TestSimpleReassign3() {
	c.Test(
		c.Srcs(
			`/** @define {boolean} */ var DEF = false;
var x;
x = DEF = true;
`),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testDefineAssignedToSimpleAlias()
func (c *ProcessDefinesTest) TestDefineAssignedToSimpleAlias() {
	c.TestSame(
		`const x = true;
const ALIAS = x;
/** @define {boolean} */ const DEF2 = ALIAS;
`)
}

// Java: ProcessDefinesTest.testDefineAssignedToNonConstAlias()
func (c *ProcessDefinesTest) TestDefineAssignedToNonConstAlias() {
	c.TestError(
		`let X = true;
X = false;
/** @define {boolean} */ const DEF2 = X;
`,
		ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testDefineAssignedToEnumAlias()
func (c *ProcessDefinesTest) TestDefineAssignedToEnumAlias() {
	c.TestError(
		`/** @enum {string} */ const E = {A: 'a'};
/** @define {string} */ const DEF2 = E.A;
`,
		// TODO(sdh): It would be nice if this worked, but doesn't seem worth implementing.
		ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testDefineAssignedToDefineAlias()
func (c *ProcessDefinesTest) TestDefineAssignedToDefineAlias() {
	c.overrides.Put("DEF2", rhino.NewNode(rhino.Token_TRUE))
	c.Test(
		`/** @define {boolean} */ const DEF1 = false;
const ALIAS = DEF1;
/** @define {boolean} */ const DEF2 = ALIAS;
`,
		`/** @define {boolean} */ const DEF1 = false;
const ALIAS = DEF1;
/** @define {boolean} */ const DEF2 = true;
`)
}

// Java: ProcessDefinesTest.testDefineAssignedToQualifiedNameAlias()
func (c *ProcessDefinesTest) TestDefineAssignedToQualifiedNameAlias() {
	c.overrides.Put("DEF1", rhino.NewNode(rhino.Token_TRUE))
	c.Test(
		`const ns = {};
/** @define {boolean} */ const DEF1 = false;
/** @const */ ns.ALIAS = DEF1;
/** @define {boolean} */ const DEF2 = ns.ALIAS;
`,
		`const ns = {};
/** @define {boolean} */ const DEF1 = true;
/** @const */ ns.ALIAS = DEF1;
/** @define {boolean} */ const DEF2 = ns.ALIAS;
`)
}

// Java: ProcessDefinesTest.testDefineAssignedToNonconstDefineAlias()
func (c *ProcessDefinesTest) TestDefineAssignedToNonconstDefineAlias() {
	c.TestError(
		`/** @define {boolean} */ const DEF1 = false;
var ALIAS = DEF1;
/** @define {boolean} */ const DEF2 = ALIAS;
`,
		ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testDefineAssignedToNonconstQualifiedNameAlias()
func (c *ProcessDefinesTest) TestDefineAssignedToNonconstQualifiedNameAlias() {
	c.TestError(
		`const ns = {};
/** @define {boolean} */ const DEF1 = false;
ns.ALIAS = DEF1;
/** @define {boolean} */ const DEF2 = ns.ALIAS;
`,
		ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testAssignBeforeDeclaration_var()
func (c *ProcessDefinesTest) TestAssignBeforeDeclaration_var() {
	c.Test(
		c.Srcs("DEF=false;var b=false,/** @define {boolean} */DEF=true,c=false"),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testAssignBeforeDeclaration_const()
func (c *ProcessDefinesTest) TestAssignBeforeDeclaration_const() {
	c.Test(
		c.Srcs("DEF=false;const b=false,/** @define {boolean} */DEF=true,c=false"),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testAssignBeforeDeclaration_var_withOverride()
func (c *ProcessDefinesTest) TestAssignBeforeDeclaration_var_withOverride() {
	c.overrides.Put("DEF_OVERRIDE_TO_TRUE", rhino.NewNode(rhino.Token_TRUE))
	c.Test(
		c.Srcs(
			`DEF_OVERRIDE_TO_TRUE = 3;
/** @define {boolean|number} */ var DEF_OVERRIDE_TO_TRUE = false;
`),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testEmptyDeclaration()
func (c *ProcessDefinesTest) TestEmptyDeclaration() {
	c.TestError("/** @define {boolean} */ var DEF;", ProcessDefines_INVALID_DEFINE_VALUE)
}

// Java: ProcessDefinesTest.testReassignAfterCall()
func (c *ProcessDefinesTest) TestReassignAfterCall() {
	c.TestError(
		"/** @define {boolean} */var DEF=true;externMethod();DEF=false",
		ProcessDefines_NON_CONST_DEFINE)
}

// Java: ProcessDefinesTest.testReassignAfterRef()
func (c *ProcessDefinesTest) TestReassignAfterRef() {
	c.TestError(
		"/** @define {boolean} */var DEF=true;var x = DEF;DEF=false",
		ProcessDefines_NON_CONST_DEFINE)
}

// Java: ProcessDefinesTest.testReassignWithExpr()
func (c *ProcessDefinesTest) TestReassignWithExpr() {
	c.Test(
		c.Srcs("/** @define {boolean} */var DEF=true;var x;DEF=x=false"),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testReassignAfterRefInConditional()
func (c *ProcessDefinesTest) TestReassignAfterRefInConditional() {
	c.TestError(
		"/** @define {boolean} */var DEF=true; if (false) {var x=DEF} DEF=false;",
		ProcessDefines_NON_CONST_DEFINE)
}

// Java: ProcessDefinesTest.testAssignInFunctionScope()
func (c *ProcessDefinesTest) TestAssignInFunctionScope() {
	c.TestError(
		"/** @define {boolean} */var DEF=true;function foo() {DEF=false};",
		ProcessDefines_NON_CONST_DEFINE)
}

// Java: ProcessDefinesTest.testDeclareInFunctionScope()
func (c *ProcessDefinesTest) TestDeclareInFunctionScope() {
	c.TestError(
		"function foo() {/** @define {boolean} */var DEF=true;};",
		ProcessDefines_INVALID_DEFINE_LOCATION)
}

// Java: ProcessDefinesTest.testDeclareInFunctionScope_withOtherSet()
func (c *ProcessDefinesTest) TestDeclareInFunctionScope_withOtherSet() {
	c.Test(
		c.Srcs(
			`var DEF = 0;
function foo() {
  /** @define {boolean} */ DEF=true;
};
`),
		c.Error(ProcessDefines_INVALID_DEFINE_LOCATION))
}

// Java: ProcessDefinesTest.testDeclareInBlockScope()
func (c *ProcessDefinesTest) TestDeclareInBlockScope() {
	c.TestError(
		"{ /** @define {boolean} */ const DEF=true; };", //
		ProcessDefines_INVALID_DEFINE_LOCATION)
}

// Java: ProcessDefinesTest.testDeclareInClassStaticBlock()
func (c *ProcessDefinesTest) TestDeclareInClassStaticBlock() {
	c.TestError(
		"class C {static {/** @define {boolean} */ const DEF=true;}}",
		ProcessDefines_INVALID_DEFINE_LOCATION)
}

// Java: ProcessDefinesTest.testDefineAssignmentInLoop()
func (c *ProcessDefinesTest) TestDefineAssignmentInLoop() {
	c.TestError(
		"/** @define {boolean} */var DEF=true;var x=0;while (x) {DEF=false;}",
		ProcessDefines_NON_CONST_DEFINE)
}

// Java: ProcessDefinesTest.testWithNoDefines()
func (c *ProcessDefinesTest) TestWithNoDefines() {
	c.TestSame("var DEF=true;var x={};x.foo={}")
}

// Java: ProcessDefinesTest.testNamespacedDefine1()
func (c *ProcessDefinesTest) TestNamespacedDefine1() {
	c.Test(
		c.Srcs("var a = {}; /** @define {boolean} */ a.B = false; a.B = true;"),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testNamespacedDefine2a()
func (c *ProcessDefinesTest) TestNamespacedDefine2a() {
	c.overrides.Put("a.B", rhino.NewNode(rhino.Token_TRUE))
	c.Test(
		"var a = {}; /** @define {boolean} */ a.B = false;",
		"var a = {}; /** @define {boolean} */ a.B = true;")
}

// Java: ProcessDefinesTest.testNamespacedDefine2b()
func (c *ProcessDefinesTest) TestNamespacedDefine2b() {
	c.overrides.Put("a.B", rhino.NewNode(rhino.Token_TRUE))
	c.TestError(
		"var a = { /** @define {boolean} */ B : false };", //
		ProcessDefines_INVALID_DEFINE_LOCATION)
}

// Java: ProcessDefinesTest.testNamespacedDefine2c()
func (c *ProcessDefinesTest) TestNamespacedDefine2c() {
	c.overrides.Put("a.B", rhino.NewNode(rhino.Token_TRUE))
	c.TestError(
		"var a = { /** @define {boolean} */ get B() { return false } };",
		ProcessDefines_INVALID_DEFINE_LOCATION)
}

// Java: ProcessDefinesTest.testNamespacedDefine3()
func (c *ProcessDefinesTest) TestNamespacedDefine3() {
	c.overrides.Put("a.B", rhino.NewNode(rhino.Token_TRUE))
	c.Test("var a = {};", "var a = {};", c.Warning(ProcessDefines_UNKNOWN_DEFINE_WARNING))
}

// Java: ProcessDefinesTest.testNamespacedDefine4()
func (c *ProcessDefinesTest) TestNamespacedDefine4() {
	c.overrides.Put("a.B", rhino.NewNode(rhino.Token_TRUE))
	c.Test(
		"var a = {}; /** @define {boolean} */ a.B = false;",
		"var a = {}; /** @define {boolean} */ a.B = true;")
}

// Java: ProcessDefinesTest.testGoogDefine_notOverridden()
func (c *ProcessDefinesTest) TestGoogDefine_notOverridden() {
	c.Test(
		"/** @define {boolean} */ const B = goog.define('a.B', false);",
		"/** @define {boolean} */ const B = false;")
}

// Java: ProcessDefinesTest.testGoogDefine_overridden()
func (c *ProcessDefinesTest) TestGoogDefine_overridden() {
	c.overrides.Put("a.B", rhino.NewNode(rhino.Token_TRUE))
	c.Test(
		"/** @define {boolean} */ const B = goog.define('a.B', false);",
		"/** @define {boolean} */ const B = true;")
}

// Java: ProcessDefinesTest.testGoogDefineAllowedFormats_notOverridden()
func (c *ProcessDefinesTest) TestGoogDefineAllowedFormats_notOverridden() {
	jsdoc := "/** @define {number} */\n"
	c.Test(jsdoc+"var name = goog.define('name', 1);", jsdoc+"var name = 1")
	c.Test(jsdoc+"var name = goog.define('otherName', 1);", jsdoc+"var name = 1")
	c.Test(jsdoc+"const name = goog.define('name', 1);", jsdoc+"const name = 1")
	c.Test(
		"const ns = {};\n"+jsdoc+"ns.name = goog.define('ns.name', 1);",
		"const ns = {};\n"+jsdoc+"ns.name = 1;")
}

// Java: ProcessDefinesTest.testGoogDefine_invalidCallFormats()
func (c *ProcessDefinesTest) TestGoogDefine_invalidCallFormats() {
	c.mode = ProcessDefines_Mode_CHECK

	jsdoc := "/** @define {number} */\n"
	c.TestError("const name = goog.define('name', 1);", ProcessDefines_MISSING_DEFINE_ANNOTATION)
	c.TestError(
		"const name = {};\n"+jsdoc+"name.two = goog.define('name.2', 1);",
		ProcessDefines_INVALID_DEFINE_NAME_ERROR)
	c.TestError(jsdoc+"const x = goog.define();", ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR)
	c.TestError(
		jsdoc+"const value = goog.define('value');", ClosurePrimitiveErrors_NULL_ARGUMENT_ERROR)
	c.TestError(
		jsdoc+"const five = goog.define(5);", ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR)
	c.TestError(
		jsdoc+"const five = goog.define('FOO', 5, 6);",
		ClosurePrimitiveErrors_TOO_MANY_ARGUMENTS_ERROR)

	c.TestError(
		jsdoc+"const templateName = goog.define(`templateName`, 1);",
		ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR)
	c.TestError(
		jsdoc+"const templateName = goog.define(`${template}Name`, 1);",
		ClosurePrimitiveErrors_INVALID_ARGUMENT_ERROR)
}

// Java: ProcessDefinesTest.testDefineWithoutGoogDefine_disabledByDefault()
func (c *ProcessDefinesTest) TestDefineWithoutGoogDefine_disabledByDefault() {
	c.mode = ProcessDefines_Mode_CHECK
	c.TestSame("/** @define {boolean} */ const MY_FLAG = true;")
}

// Java: ProcessDefinesTest.testDefineWithoutGoogDefine()
func (c *ProcessDefinesTest) TestDefineWithoutGoogDefine() {
	c.mode = ProcessDefines_Mode_CHECK
	c.enableDefineWithoutGoogDefineCheck = true
	c.TestError(
		"/** @define {boolean} */ const MY_FLAG = true;",
		ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE)
}

// Java: ProcessDefinesTest.testDefineWithoutGoogDefine_compiledInBaseJs()
func (c *ProcessDefinesTest) TestDefineWithoutGoogDefine_compiledInBaseJs() {
	c.mode = ProcessDefines_Mode_CHECK
	c.enableDefineWithoutGoogDefineCheck = true
	c.TestSame(
		c.Srcs(SourceFile_FromCode("base.js", "/** @define {boolean} */ var COMPILED = false;")))
}

// Java: ProcessDefinesTest.testDefineWithoutGoogDefine_jscompOptOut()
func (c *ProcessDefinesTest) TestDefineWithoutGoogDefine_jscompOptOut() {
	c.mode = ProcessDefines_Mode_CHECK
	c.enableDefineWithoutGoogDefineCheck = true
	c.TestSame("var $jscomp = {}; /** @define {boolean} */ $jscomp.FOO = true;")
}

// Java: ProcessDefinesTest.testDefineWithoutGoogDefine_suppressed()
func (c *ProcessDefinesTest) TestDefineWithoutGoogDefine_suppressed() {
	c.mode = ProcessDefines_Mode_CHECK
	c.enableDefineWithoutGoogDefineCheck = true
	c.TestSame(
		"/** @fileoverview @suppress {defineWithoutGoogDefine} */ /** @define {boolean} */ const" +
			" MY_FLAG = true;")
	c.TestSame("/** @suppress {defineWithoutGoogDefine} @define {boolean} */ const MY_FLAG = true;")
}

// Java: ProcessDefinesTest.testGoogDefine_invalidCallLocation()
func (c *ProcessDefinesTest) TestGoogDefine_invalidCallLocation() {
	c.mode = ProcessDefines_Mode_CHECK

	c.Test(c.Srcs("goog.define('name', 1);"), c.Error(ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT))
	c.Test(
		c.Srcs("/** @define {number} */ goog.define('name', 1);"),
		c.Error(ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT),
		c.Error(ProcessDefines_INVALID_DEFINE_LOCATION))
	c.TestError(
		"var x = x || goog.define('goog.DEBUG', true);",
		ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT)
	c.TestError(
		"function f() { const debug = goog.define('goog.DEBUG', true); }",
		ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR)
}

// Java: ProcessDefinesTest.testGoogDefine_enabledByRecognizeClosureDefines()
func (c *ProcessDefinesTest) TestGoogDefine_enabledByRecognizeClosureDefines() {
	c.mode = ProcessDefines_Mode_CHECK
	c.Test(c.Srcs("goog.define('name', 1);"), c.Error(ProcessDefines_DEFINE_CALL_WITHOUT_ASSIGNMENT))

	c.recognizeClosureDefines = false
	c.TestSame("goog.define('name', 1);")
}

// Java: ProcessDefinesTest.testOverrideAfterAlias()
func (c *ProcessDefinesTest) TestOverrideAfterAlias() {
	c.TestError(
		"var x; /** @define {boolean} */var DEF=true; x=DEF; DEF=false;",
		ProcessDefines_NON_CONST_DEFINE)
}

// Java: ProcessDefinesTest.testBasicConstDeclaration()
func (c *ProcessDefinesTest) TestBasicConstDeclaration() {
	c.Test("/** @define {boolean} */ const DEF = true", "/** @define {boolean} */ const DEF=true")
	c.Test("/** @define {string} */ const DEF = 'a'", "/** @define {string} */ const DEF=\"a\"")
	c.Test("/** @define {number} */ const DEF = 0", "/** @define {number} */ const DEF=0")
}

// Java: ProcessDefinesTest.testConstOverriding1()
func (c *ProcessDefinesTest) TestConstOverriding1() {
	c.overrides.Put("DEF_OVERRIDE_TO_TRUE", rhino.NewNode(rhino.Token_TRUE))
	c.Test(
		"/** @define {boolean} */ const DEF_OVERRIDE_TO_TRUE = false;",
		"/** @define {boolean} */ const DEF_OVERRIDE_TO_TRUE = true;")
}

// Java: ProcessDefinesTest.testConstOverriding2()
func (c *ProcessDefinesTest) TestConstOverriding2() {
	c.Test(
		"/** @define {string} */ const DEF_OVERRIDE_STRING = 'x';",
		"/** @define {string} */ const DEF_OVERRIDE_STRING=\"x\"")
}

// Java: ProcessDefinesTest.testConstProducesUnknownDefineWarning()
func (c *ProcessDefinesTest) TestConstProducesUnknownDefineWarning() {
	c.mode = ProcessDefines_Mode_OPTIMIZE
	c.overrides.Put("a.B", rhino.NewNode(rhino.Token_TRUE))
	c.Test("const a = {};", "const a = {};", c.Warning(ProcessDefines_UNKNOWN_DEFINE_WARNING))
}

// Java: ProcessDefinesTest.testSimpleConstReassign()
func (c *ProcessDefinesTest) TestSimpleConstReassign() {
	c.TestError(
		`/** @define {boolean} */ const DEF = false;
DEF = true;
`,
		ProcessDefines_NON_CONST_DEFINE)
}

// Java: ProcessDefinesTest.testRedeclaration_twoGoogDefine_differentLocalNames()
func (c *ProcessDefinesTest) TestRedeclaration_twoGoogDefine_differentLocalNames() {
	c.Test(
		c.Srcs(
			`/** @define {boolean} */
const A = goog.define('a.B', false);

/** @define {boolean} */
const B = goog.define('a.B', false);
`),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testRedeclaration_oneGoogDefine_varWithGoogDefineName()
func (c *ProcessDefinesTest) TestRedeclaration_oneGoogDefine_varWithGoogDefineName() {
	c.Test(
		c.Srcs(
			`/** @define {boolean} */
const A = goog.define('B', false);

const B = false;
`),
		c.Expected(
			`/** @define {boolean} */
const A = false;

const B = false;
`))
}

// Java: ProcessDefinesTest.testRedeclaration_oneGoogDefine_varWithSameLocalName()
func (c *ProcessDefinesTest) TestRedeclaration_oneGoogDefine_varWithSameLocalName() {
	c.Test(
		c.Srcs(
			`/** @define {boolean} */
var A = goog.define('B', false);

var A = false;
`),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testRedeclaration_oneGoogDefine_oneAtDefine()
func (c *ProcessDefinesTest) TestRedeclaration_oneGoogDefine_oneAtDefine() {
	c.Test(
		c.Srcs(
			`/** @define {boolean} */
const A = goog.define('B', false);

/** @define {boolean} */
const B = false;
`),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testRedeclaration_oneAtDefine_varWithSameName()
func (c *ProcessDefinesTest) TestRedeclaration_oneAtDefine_varWithSameName() {
	c.Test(
		c.Srcs(
			`/** @define {boolean} */
var A = false;

var A = false;
`),
		c.Error(ProcessDefines_NON_CONST_DEFINE))
}

// Java: ProcessDefinesTest.testClosureDefineValues_replacements()
func (c *ProcessDefinesTest) TestClosureDefineValues_replacements() {
	c.mode = ProcessDefines_Mode_CHECK_AND_OPTIMIZE
	c.Test(
		`var CLOSURE_DEFINES = {'FOO': 'closureDefault'};
/** @define {string} */ const FOO = 'original';
`,
		`var CLOSURE_DEFINES = {'FOO': 'closureDefault'};
/** @define {string} */ const FOO = 'closureDefault';
`)

	c.Test(
		`var CLOSURE_DEFINES = CLOSURE_DEFINES || {};
CLOSURE_DEFINES['FOO'] = 'closureDefault';
/** @define {string} */ const FOO = 'original';
`,
		`var CLOSURE_DEFINES = CLOSURE_DEFINES || {};
CLOSURE_DEFINES['FOO'] = 'closureDefault';
/** @define {string} */ const FOO = 'closureDefault';
`)

	c.Test(
		`var CLOSURE_DEFINES = {'FOO': true};
/** @define {boolean} */ const FOO = false;
`,
		`var CLOSURE_DEFINES = {'FOO': true};
/** @define {boolean} */ const FOO = true;
`)

	c.Test(
		`var CLOSURE_DEFINES = {'FOO': false};
/** @define {boolean} */ const FOO = false;
`,
		`var CLOSURE_DEFINES = {'FOO': false};
/** @define {boolean} */ const FOO = false;
`)

	c.Test(
		`var CLOSURE_DEFINES = {'FOO': 1};
/** @define {number} */ const FOO = 2;
`,
		`var CLOSURE_DEFINES = {'FOO': 1};
/** @define {number} */ const FOO = 1;
`)

	c.Test(
		`var CLOSURE_DEFINES = {'FOO': 0xABCD};
/** @define {number} */ const FOO = 2;
`,
		`var CLOSURE_DEFINES = {'FOO': 0xABCD};
/** @define {number} */ const FOO = 0xABCD;
`)
}

// Java: ProcessDefinesTest.testClosureDefineValues_duplicateKey()
func (c *ProcessDefinesTest) TestClosureDefineValues_duplicateKey() {
	c.mode = ProcessDefines_Mode_CHECK_AND_OPTIMIZE
	c.Test(
		c.Srcs(
			`var CLOSURE_DEFINES = CLOSURE_DEFINES || {};
CLOSURE_DEFINES['FOO'] = 'firstVersionIgnored';
CLOSURE_DEFINES['FOO'] = 'closureDefault';
/** @define {string} */ const FOO = 'original';
`),
		c.Error(ProcessDefines_CLOSURE_DEFINES_MULTIPLE))

	c.Test(
		c.Srcs(
			`var CLOSURE_DEFINES = {'FOO': 'firstVersionIgnored'};
CLOSURE_DEFINES['FOO'] = 'closureDefault';
/** @define {string} */ const FOO = 'original';
`),
		c.Error(ProcessDefines_CLOSURE_DEFINES_MULTIPLE))
}

// Java: ProcessDefinesTest.testClosureDefineValues_namespacedReplacement()
func (c *ProcessDefinesTest) TestClosureDefineValues_namespacedReplacement() {
	c.Test(
		`var CLOSURE_DEFINES = {'a.b': 'closureDefault'};
const a = {};
/** @define {string} */ a.b = 'original';
`,
		`var CLOSURE_DEFINES = {'a.b': 'closureDefault'};
const a = {};
/** @define {string} */ a.b = 'closureDefault';
`)
}

// Java: ProcessDefinesTest.testClosureDefineValues_replacementWithGoogDefine()
func (c *ProcessDefinesTest) TestClosureDefineValues_replacementWithGoogDefine() {
	c.mode = ProcessDefines_Mode_CHECK_AND_OPTIMIZE
	c.Test(
		`var CLOSURE_DEFINES = {'FOO': 'closureDefault'};
/** @define {string} */ const f = goog.define('FOO', 'original');
`,
		`var CLOSURE_DEFINES = {'FOO': 'closureDefault'};
/** @define {string} */ const f = 'closureDefault';
`)
}

// Java: ProcessDefinesTest.testClosureDefineValues_replacementWithOverriddenDefine()
func (c *ProcessDefinesTest) TestClosureDefineValues_replacementWithOverriddenDefine() {
	// Command-line flag takes precedence over CLOSURE_DEFINES.
	c.mode = ProcessDefines_Mode_CHECK_AND_OPTIMIZE
	c.overrides.Put("FOO", rhino.IR_String("override"))
	c.Test(
		`var CLOSURE_DEFINES = {'FOO': 'closureDefault'};
/** @define {string} */ const f = goog.define('FOO', 'original');
`,
		`var CLOSURE_DEFINES = {'FOO': 'closureDefault'};
/** @define {string} */ const f = 'override';
`)
}

// Java: ProcessDefinesTest.testClosureDefineValues_checkOnlyDoesntModifyAst()
func (c *ProcessDefinesTest) TestClosureDefineValues_checkOnlyDoesntModifyAst() {
	c.mode = ProcessDefines_Mode_CHECK
	c.TestSame(
		`var CLOSURE_DEFINES = {'FOO': 'string'};
/** @define {string} */ const f = goog.define('FOO', 'tmp');
`)
}

// Java: ProcessDefinesTest.testClosureDefines_unknownDefineErrors()
func (c *ProcessDefinesTest) TestClosureDefines_unknownDefineErrors() {
	c.mode = ProcessDefines_Mode_OPTIMIZE
	c.Test(c.Srcs("var CLOSURE_DEFINES = {'FOO': 0};"), c.Warning(ProcessDefines_UNKNOWN_DEFINE_WARNING))
	c.Test(c.Srcs("CLOSURE_DEFINES['FOO'] = 0;"), c.Warning(ProcessDefines_UNKNOWN_DEFINE_WARNING))
}

// Java: ProcessDefinesTest.testClosureDefines_mustBeGlobal()
func (c *ProcessDefinesTest) TestClosureDefines_mustBeGlobal() {
	c.mode = ProcessDefines_Mode_CHECK
	c.Test(
		c.Srcs(
			`(function() { var CLOSURE_DEFINES = {'FOO': 0}; })();
/** @define {number} */ const FOO = 1;
`),
		c.Error(ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR))

	c.Test(
		c.Srcs(
			`if (cond) {
  var CLOSURE_DEFINES = {'FOO': 0};
}
/** @define {number} */ const FOO = 1;
`),
		c.Error(ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR))

	c.Test(
		c.Srcs(
			`if (cond) {
  CLOSURE_DEFINES['FOO'] = 0;
}
/** @define {number} */ const FOO = 1;
`),
		c.Error(ProcessDefines_NON_GLOBAL_CLOSURE_DEFINES_ERROR))
}

// Java: ProcessDefinesTest.testClosureDefines_valuesErrors()
func (c *ProcessDefinesTest) TestClosureDefines_valuesErrors() {
	c.mode = ProcessDefines_Mode_CHECK
	c.TestError("var CLOSURE_DEFINES = {'FOO': a};", ProcessDefines_CLOSURE_DEFINES_ERROR)
	c.TestError("var CLOSURE_DEFINES = {'FOO': 0+1};", ProcessDefines_CLOSURE_DEFINES_ERROR)
	c.TestError(
		"var CLOSURE_DEFINES = {'FOO': 'value' + 'value'};", ProcessDefines_CLOSURE_DEFINES_ERROR)
	c.TestError("var CLOSURE_DEFINES = {'FOO': !true};", ProcessDefines_CLOSURE_DEFINES_ERROR)
	c.TestError("var CLOSURE_DEFINES = {'FOO': -true};", ProcessDefines_CLOSURE_DEFINES_ERROR)

	c.TestError("var CLOSURE_DEFINES = {SHORTHAND};", ProcessDefines_CLOSURE_DEFINES_ERROR)
	c.TestError(
		"var CLOSURE_DEFINES = {'TEMPLATE': `template`};", ProcessDefines_CLOSURE_DEFINES_ERROR)
	c.TestError(
		"var CLOSURE_DEFINES = {'TEMPLATE': `${template}Sub`};",
		ProcessDefines_CLOSURE_DEFINES_ERROR)

	c.TestError("CLOSURE_DEFINES[notStringLiteral] = 42;", ProcessDefines_CLOSURE_DEFINES_ERROR)
	c.TestError("CLOSURE_DEFINES['FOO'] = a;", ProcessDefines_CLOSURE_DEFINES_ERROR)
}

// Java: ProcessDefinesTest.testClosureDefinesErrors_enabledByRecognizeClosureDefines()
func (c *ProcessDefinesTest) TestClosureDefinesErrors_enabledByRecognizeClosureDefines() {
	c.mode = ProcessDefines_Mode_CHECK

	c.TestError("var CLOSURE_DEFINES = {'FOO': a};", ProcessDefines_CLOSURE_DEFINES_ERROR)

	c.recognizeClosureDefines = false
	c.TestSame("var CLOSURE_DEFINES = {'FOO': a};")
}

// Java: ProcessDefinesTest.testGenerateGlobalAliases()
func (c *ProcessDefinesTest) TestGenerateGlobalAliases() {
	c.enableJ2clPasses = true

	c.Test(
		`var a = {};
/** @define {number} */ a.b = goog.define('a.b', 1);
`,
		`var a = {};
/** @define {number} */ a.b = 1;
var jscomp$defines$a$b = a.b;
`)
	c.Test(
		`var a = {};
/** @define {number} */ a.b = goog.define('c.d', 1);
`,
		`var a = {};
/** @define {number} */ a.b = 1;
var jscomp$defines$c$d = a.b;
`)
	c.Test(
		"/** @define {number} */ var a = goog.define('c.d', 1);",
		`/** @define {number} */ var a = 1;
var jscomp$defines$c$d = a;
`)
	c.Test(
		`/** @define {number} */ var a = goog.define('a', 1);
/** @define {number} */ var b = goog.define('b', a);
`,
		`/** @define {number} */ var a = 1;
var jscomp$defines$a = a;
/** @define {number} */ var b = a;
var jscomp$defines$b = b;
`)
}

// Java: ProcessDefinesTest.testEs2017OrGreaterFailsWhenCompilingZone()
func (c *ProcessDefinesTest) TestEs2017OrGreaterFailsWhenCompilingZone() {
	c.outputFeatureSet = parser.FeatureSet_ES2017

	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				"third_party/javascript/angular2/rc/packages/zone.js/lib/zone.closure.js",
				`/** @define {boolean} */
const ENABLE_ZONES = goog.define('javascript.angular2.ENABLE_ZONES', true)
`)),
		ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT)
}

// Java: ProcessDefinesTest.testEs2017OrGreaterAllowedWhenCompilingZoneWithItDisabled()
func (c *ProcessDefinesTest) TestEs2017OrGreaterAllowedWhenCompilingZoneWithItDisabled() {
	c.outputFeatureSet = parser.FeatureSet_ES2017

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"defines.js", "var CLOSURE_DEFINES = {'javascript.angular2.ENABLE_ZONES': false};"),
			SourceFile_FromCode(
				"third_party/javascript/angular2/rc/packages/zone.js/lib/zone.closure.js",
				`/** @define {boolean} */
const ENABLE_ZONES = goog.define('javascript.angular2.ENABLE_ZONES', true)
`)))
}

// Java: ProcessDefinesTest.testEs2017OrGreaterFailsWhenCompilingZoneWithItExplicitlyEnabled()
func (c *ProcessDefinesTest) TestEs2017OrGreaterFailsWhenCompilingZoneWithItExplicitlyEnabled() {
	c.outputFeatureSet = parser.FeatureSet_ES2017

	c.TestError(
		c.Srcs(
			SourceFile_FromCode(
				"defines.js", "var CLOSURE_DEFINES = {'javascript.angular2.ENABLE_ZONES': true};"),
			SourceFile_FromCode(
				"third_party/javascript/angular2/rc/packages/zone.js/lib/zone.closure.js",
				`/** @define {boolean} */
const ENABLE_ZONES = goog.define('javascript.angular2.ENABLE_ZONES', true)
`)),
		ProcessDefines_ZONE_NOT_SUPPORTED_WITH_NATIVE_ASYNC_AWAIT)
}

// Java: ProcessDefinesTest.testZoneAllowedWhenCompilingLessThanEs2017()
func (c *ProcessDefinesTest) TestZoneAllowedWhenCompilingLessThanEs2017() {
	c.outputFeatureSet = parser.FeatureSet_ES2016

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"third_party/javascript/angular2/rc/packages/zone.js/lib/zone.closure.js", "")))
}

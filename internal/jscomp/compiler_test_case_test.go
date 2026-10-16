// Port of com/google/javascript/jscomp/CompilerTestCase.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/CompilerTestCase.java)

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

// How a Java CompilerTestCase subclass is ported (see also junit4_test.go):
//
//	public final class FooTest extends CompilerTestCase {
//	  public FooTest() { super(DEFAULT_EXTERNS); }
//	  @Override @Before public void setUp() throws Exception { super.setUp(); enableNormalize(); }
//	  @Override protected CompilerPass getProcessor(Compiler compiler) { return new Foo(compiler); }
//	  @Test public void testBar() { test("a", "b"); testSame(srcs("c")); }
//	}
//
// becomes
//
//	type FooTest struct{ CompilerTestCaseBase }
//
//	func NewFooTest() *FooTest {
//		c := &FooTest{}
//		c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_DEFAULT_EXTERNS)
//		return c
//	}
//	func (c *FooTest) SetUp() { c.CompilerTestCaseBase.SetUp(); c.EnableNormalize() }
//	func (c *FooTest) GetProcessor(compiler *Compiler) CompilerPass { return NewFoo(compiler) }
//	func (c *FooTest) TestBar() { c.Test("a", "b"); c.TestSame(c.Srcs("c")) }
//
//	func TestFooTest(t *testing.T) { RunJUnit4(t, func() JUnitTestCase { return NewFooTest() }) }
//
// Java's overloaded test(...), testSame(...), testWarning(...),
// testError(...), testNoWarning(...), testExternChanges(...), srcs(...),
// expected(...) and externs(...) keep their plain names and take ...any:
// the overload is chosen at run time from the argument types, exactly as
// javac would choose it statically (Go-only: language-adapter). The static
// factories are package functions CompilerTestCase_Srcs, ..._Expected_,
// ..._Externs_, ..._Warning, ..._Error, ..._Postcondition_ (a trailing _
// where a nested type, e.g. CompilerTestCase_Expected, has the plain name);
// they are also methods without the suffix (c.Srcs(...), c.Expected(...),
// c.Externs(...), c.Warning(t), c.Error(t), c.Postcondition(func)) because
// Java test code calls them unqualified (Go-only: language-adapter).

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// CompilerTestCase is the base class for testing JS compiler classes that
// change the node tree of a compiled JS input.
//
// Pulls in shared functionality from different test cases. Also supports
// node tree comparison for input and output (instead of string comparison),
// which makes it easier to write tests b/c you don't have to get the syntax
// exactly correct to the spacing.
//
// The interface declares the abstract and overridable methods; subclasses
// embed CompilerTestCaseBase (PORTING.md §4).
type CompilerTestCase interface {
	JUnitTestCase

	// GetCompilerTestCaseBase gives access to the base class.
	GetCompilerTestCaseBase() *CompilerTestCaseBase

	// Java: CompilerTestCase.getName()
	GetName() string
	// Java: CompilerTestCase.setUp()
	SetUp()
	// Java: CompilerTestCase.tearDown()
	TearDown()
	// GetProcessor gets the compiler pass instance to use for a test.
	//
	// Java: CompilerTestCase.getProcessor(Compiler)
	GetProcessor(compiler *Compiler) CompilerPass
	// Java: CompilerTestCase.getOptions()
	GetOptions() *CompilerOptions
	// Java: CompilerTestCase.getCodingConvention()
	GetCodingConvention() CodingConvention
	// Java: CompilerTestCase.getNumRepetitions()
	GetNumRepetitions() int
	// Java: CompilerTestCase.getLastCompiler()
	GetLastCompiler() *Compiler
	// Java: CompilerTestCase.test(...) (every overload)
	Test(args ...any)
	// Java: CompilerTestCase.testSame(...) (every overload)
	TestSame(args ...any)
	// Java: CompilerTestCase.testError(...) (every overload)
	TestError(args ...any)
	// Java: CompilerTestCase.testWarning(...) (every overload)
	TestWarning(args ...any)
	// Java: CompilerTestCase.testNoWarning(...) (every overload)
	TestNoWarning(args ...any)
	// Java: CompilerTestCase.testExternChanges(...) (every overload)
	TestExternChanges(args ...any)
	// Java: CompilerTestCase.testInternal(Externs, Sources, Expected, List<Diagnostic>, List<Postcondition>)
	TestInternal(externs *CompilerTestCase_Externs, inputs CompilerTestCase_Sources, expected *CompilerTestCase_Expected, diagnostics []*CompilerTestCase_Diagnostic, postconditions []CompilerTestCase_Postcondition)
	// Java: CompilerTestCase.createAndInitializeCompiler(Externs, Sources)
	CreateAndInitializeCompiler(externs *CompilerTestCase_Externs, inputs CompilerTestCase_Sources) *Compiler
	// Java: CompilerTestCase.parseExpectedJs(String) / parseExpectedJs(Expected)
	ParseExpectedJs(expected any) *rhino.Node
	// Java: CompilerTestCase.createCompiler()
	CreateCompiler() *Compiler
}

// CompilerTestCase_GENERATED_SRC_NAME is the file name of generated sources.
// The file name is included in the AST display string attached to each
// node. Consistently using the same name for src or externs generated files
// avoids spurious differences in the diffs displayed when a test case fails.
// TODO(bradfordcsmith): "testcode" and "externs" are magical strings.
// Other testing code assumes these names are used and expects them.
const CompilerTestCase_GENERATED_SRC_NAME = "testcode"

// CompilerTestCase_GENERATED_EXTERNS_NAME is the file name of generated
// externs.
const CompilerTestCase_GENERATED_EXTERNS_NAME = "externs"

// CompilerTestCaseBase holds the fields and base implementations of
// CompilerTestCase.
type CompilerTestCaseBase struct {
	self CompilerTestCase

	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	// Default externs for the test. Individual calls to test() may provide
	// their own externs instead.
	defaultExternsInputs []*SourceFile

	// Libraries to inject before typechecking
	librariesToInject *jutil.LinkedHashSet[string]

	// Whether to compare input and output as trees instead of strings
	compareAsTree bool

	// Whether to parse type info from JSDoc comments
	parseTypeInfo bool

	// Whether to take JSDoc into account when comparing ASTs.
	compareJsDoc bool

	// Whether we check warnings without source information.
	allowSourcelessWarnings bool

	// True iff closure pass runs before pass being tested.
	closurePassEnabled bool

	// Whether the closure pass is run on the expected JS.
	closurePassEnabledForExpected bool

	// Whether to rewrite commons js modules before the test is run.
	processCommonJsModules bool

	// Whether to rewrite Closure code before the test is run.
	rewriteClosureCode bool

	rewriteModulesAfterTypechecking bool

	// Whether to rewrite Closure code before the test is run.
	rewriteClosureProvides bool

	// If true, run type checking together with the pass being tested. A
	// separate flag controls whether type checking runs before or after the
	// pass.
	typeCheckEnabled bool

	// If true, run the ConvertTypesToColors pass, which removes references to
	// JSTypes and converts JSDoc to its simplified optimization form.
	replaceTypesWithColors bool

	// If true performs the test using multistage compilation.
	multistageCompilation bool

	// Whether to test the compiler pass before the type check.
	runTypeCheckAfterProcessing bool

	// Whether to scan externs for property names.
	gatherExternPropertiesEnabled bool

	// Whether to run ModuleMapCreator after parsing.
	createModuleMap bool

	// Whether the Normalize pass runs before pass being tested and whether
	// the expected JS strings should be normalized.
	normalizeEnabled bool

	// If Normalize is enabled, then run it on the expected output also
	// before comparing.
	normalizeExpectedOutputEnabled bool

	polymerPass bool

	// Whether ES module rewriting runs before the pass being tested.
	rewriteEsModulesEnabled bool

	// Whether the transpilation passes run before the pass being tested.
	transpileEnabled bool

	// Whether we run InferConsts before checking.
	inferConsts bool

	// Whether we run CheckAccessControls after the pass being tested.
	checkAccessControls bool

	// Whether to check that changed scopes are marked as changed
	checkAstChangeMarking bool

	// Whether we expect parse warnings in the current test.
	expectParseWarningsInThisTest bool

	// Whether the PureFunctionIdentifier pass runs before the pass being
	// tested
	computeSideEffects bool

	// Whether the SourceInformationAnnotator pass runs after parsing
	annotateSourceInfo bool

	// The set of accessors declared to exist "somewhere" in the test program.
	//
	// When this field is populated, automatic getter and setter collection is
	// disabled.
	declaredAccessors *jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind]

	// The most recently used Compiler instance.
	lastCompiler *Compiler

	// Whether to accept ES6, ES5 or ES3 source.
	acceptedLanguage CompilerOptions_LanguageMode

	languageOut CompilerOptions_LanguageMode

	browserFeaturesetYear *int

	// How to interpret ES6 module imports
	moduleResolutionMode ModuleLoader_ResolutionMode

	// How to parse JS Documentation.
	parseJsDocDocumentation parsing.Config_JsDocParsing

	// Compiler option.
	assumeStaticInheritanceIsNotUsed bool

	// Whether externs changes should be allowed for this pass.
	allowExternsChanges bool

	// Whether the AST should be validated.
	astValidationEnabled bool

	// Whether we should verify that type information is present for every
	// AST node that should have it, and not present for those that shouldn't
	// have it.
	typeInfoValidationEnabled bool

	// Whether we should verify that for every SCRIPT node, all features
	// present in its AST are also present in its FeatureSet property.
	scriptFeatureValidationEnabled bool

	ignoredWarnings *jutil.ObjectLinkedHashSet[*DiagnosticType]

	webpackModulesById *jutil.LinkedHashMap[string, string]

	genericNameReplacements *jutil.LinkedHashMap[string, string]

	// Whether setUp has run.
	setUpRan bool

	// Enable debug logging for this test run
	debugLoggingEnabled bool
}

// compilerTestCase_LINE_JOINER is Joiner.on('\n').
func compilerTestCase_LINE_JOINER[T any](parts []T) string {
	strs := make([]string, len(parts))
	for i, p := range parts {
		strs[i] = jutil.ToJavaString(p)
	}
	return strings.Join(strs, "\n")
}

// CompilerTestCase_VAR_CHECK_EXTERNS: these externs are inserted by VarCheck
// and should not be input as externs for tests. It is provided here for
// tests that assert on output externs in cases where pulling in the extra
// externs is unavoidable.
var CompilerTestCase_VAR_CHECK_EXTERNS = compilerTestCase_varCheckExterns()

func compilerTestCase_varCheckExterns() string {
	var sb strings.Builder
	for s := range VarCheck_REQUIRED_SYMBOLS.All() {
		sb.WriteString("var " + s + ";\n")
	}
	return sb.String()
}

// CompilerTestCase_MINIMAL_EXTERNS is a minimal set of externs, consisting
// of only those needed for the typechecker not to blow up.
var CompilerTestCase_MINIMAL_EXTERNS = NewTestExternsBuilder().
	AddArray().
	AddIterable().
	AddObject().
	AddUndefined().
	AddFunction().
	AddString().
	Build()

// CompilerTestCase_DEFAULT_EXTERNS is a default set of externs for testing.
var CompilerTestCase_DEFAULT_EXTERNS = NewTestExternsBuilder().
	AddArray().
	AddIterable().
	AddObject().
	AddUndefined().
	AddFunction().
	AddString().
	AddPromise().
	AddClosureExterns().
	AddITemplateArray().
	Build()

// InitCompilerTestCaseBaseWithExterns constructs a test.
//
// externs: Externs JS as a string
//
// Java: CompilerTestCase(String)
func (b *CompilerTestCaseBase) InitCompilerTestCaseBaseWithExterns(self CompilerTestCase, externs string) {
	b.self = self
	b.declaredAccessors = jutil.NewLinkedHashMap[string, AccessorSummary_PropertyAccessKind]()
	b.ignoredWarnings = jutil.NewObjectLinkedHashSet[*DiagnosticType]()
	b.webpackModulesById = jutil.NewLinkedHashMap[string, string]()
	b.genericNameReplacements = jutil.NewLinkedHashMap[string, string]()
	b.rewriteModulesAfterTypechecking = true
	b.defaultExternsInputs = []*SourceFile{SourceFile_FromCode(CompilerTestCase_GENERATED_EXTERNS_NAME, externs)}
	b.librariesToInject = jutil.NewLinkedHashSet[string]()
}

// InitCompilerTestCaseBase constructs a test. Uses AST comparison and no
// externs.
//
// Java: CompilerTestCase()
func (b *CompilerTestCaseBase) InitCompilerTestCaseBase(self CompilerTestCase) {
	b.InitCompilerTestCaseBaseWithExterns(self, "")
}

// GetCompilerTestCaseBase gives access to the base class.
func (b *CompilerTestCaseBase) GetCompilerTestCaseBase() *CompilerTestCaseBase { return b }

// JUnitSetT is part of JUnitTestCase.
func (b *CompilerTestCaseBase) JUnitSetT(t *testing.T) { b.t = t }

// T is the Go test running the current test method (Go-only: test-helper),
// e.g. for c.T().Skip(...) in place of @Ignore.
func (b *CompilerTestCaseBase) T() *testing.T { return b.t }

// JUnitBefore runs the @Before method setUp().
func (b *CompilerTestCaseBase) JUnitBefore() { b.self.SetUp() }

// JUnitAfter runs tearDown(), which subclasses override with an @After
// method.
func (b *CompilerTestCaseBase) JUnitAfter() { b.self.TearDown() }

// Java: CompilerTestCase.getName()
func (b *CompilerTestCaseBase) GetName() string {
	// Java: this.getClass().getSimpleName()
	name := fmt.Sprintf("%T", b.self)
	if i := strings.LastIndexByte(name, '.'); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// SetUp is overridden here so that we can easily find all classes that
// override.
//
// Java: CompilerTestCase.setUp()
func (b *CompilerTestCaseBase) SetUp() {
	// TODO(sdh): Initialize *all* the options here, but first we must ensure no subclass
	// is changing them in the constructor, rather than in their own setUp method.
	b.acceptedLanguage = CompilerOptions_LanguageMode_UNSUPPORTED
	b.moduleResolutionMode = ModuleLoader_ResolutionMode_BROWSER
	b.parseJsDocDocumentation = parsing.Config_JsDocParsing_TYPES_ONLY
	b.allowExternsChanges = false
	b.allowSourcelessWarnings = false
	b.astValidationEnabled = true
	b.typeInfoValidationEnabled = false
	b.scriptFeatureValidationEnabled = true
	b.checkAccessControls = false
	b.checkAstChangeMarking = true
	b.closurePassEnabled = false
	b.closurePassEnabledForExpected = false
	b.compareAsTree = true
	b.compareJsDoc = true
	b.computeSideEffects = false
	b.annotateSourceInfo = false
	b.expectParseWarningsInThisTest = false
	b.gatherExternPropertiesEnabled = false
	b.inferConsts = false
	b.languageOut = CompilerOptions_LanguageMode_NO_TRANSPILE
	b.browserFeaturesetYear = nil
	b.multistageCompilation = false
	b.normalizeEnabled = false
	b.normalizeExpectedOutputEnabled = false
	b.parseTypeInfo = false
	b.polymerPass = false
	b.processCommonJsModules = false
	b.rewriteClosureCode = false
	// default to true: this doesn't have any effect unless rewriteClosureCode is enabled. The
	// default matches options.setBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf.
	b.rewriteModulesAfterTypechecking = true
	b.runTypeCheckAfterProcessing = false
	b.rewriteEsModulesEnabled = false
	b.transpileEnabled = false
	b.typeCheckEnabled = false
	b.assumeStaticInheritanceIsNotUsed = true

	b.setUpRan = true
}

// Java: CompilerTestCase.tearDown()
func (b *CompilerTestCaseBase) TearDown() {
	b.setUpRan = false
}

// Java: CompilerTestCase.enableDebugLogging(boolean)
func (b *CompilerTestCaseBase) EnableDebugLogging(debugLoggingEnabled bool) {
	b.debugLoggingEnabled = debugLoggingEnabled
}

// GetOptions gets the compiler options to use for this test. Use
// getProcessor to determine what passes should be run.
//
// Java: CompilerTestCase.getOptions()
func (b *CompilerTestCaseBase) GetOptions() *CompilerOptions {
	options := NewCompilerOptions()

	options.SetLanguageIn(b.acceptedLanguage)
	options.SetEmitUseStrict(false)
	options.SetLanguageOut(b.languageOut)
	if b.browserFeaturesetYear != nil {
		// Must be set after languageOut, because it must override languageOut.
		options.SetBrowserFeaturesetYear(*b.browserFeaturesetYear)
	}
	options.SetModuleResolutionMode(b.moduleResolutionMode)
	options.SetParseJsDocDocumentation(b.parseJsDocDocumentation)
	options.SetAssumeStaticInheritanceIsNotUsed(b.assumeStaticInheritanceIsNotUsed)
	options.SetPreserveTypeAnnotations(true)
	options.SetAssumeGettersArePure(false) // Default to the complex case.

	// This doesn't affect whether checkSymbols is run--it just affects
	// whether variable warnings are filtered.
	options.SetCheckSymbols(true)

	options.SetWarningLevel(DiagnosticGroups_INVALID_CASTS, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_MISPLACED_MSG_ANNOTATION, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_MISSING_PROPERTIES, CheckLevel_WARNING)
	if !b.ignoredWarnings.IsEmpty() {
		options.SetWarningLevel(
			NewDiagnosticGroup(b.ignoredWarnings.ToSlice()...), CheckLevel_OFF)
	}
	options.SetCodingConvention(b.self.GetCodingConvention())
	polymerVersion := int32(1)
	options.SetPolymerVersion(&polymerVersion)
	if b.debugLoggingEnabled {
		CompilerTestCaseUtils_SetDebugLogDirectoryOn(options)
	}
	options.SetRuntimeLibraryMode(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS)

	return options
}

// Java: CompilerTestCase.getCodingConvention()
func (b *CompilerTestCaseBase) GetCodingConvention() CodingConvention {
	return NewGoogleCodingConvention()
}

// checkSetUpRan is Java's checkState(this.setUpRan, "Attempted to configure
// before running setUp().").
func (b *CompilerTestCaseBase) checkSetUpRan() {
	jutil.CheckStateMsg(b.setUpRan, "Attempted to configure before running setUp().")
}

// EnableParseTypeInfo enables parsing type info from JSDoc comments. This
// sets the compiler option, but does not actually run the type checking
// pass.
//
// Java: CompilerTestCase.enableParseTypeInfo()
func (b *CompilerTestCaseBase) EnableParseTypeInfo() {
	b.checkSetUpRan()
	b.parseTypeInfo = true
}

// DisableCompareJsDoc turns off taking JSDoc into account when comparing
// ASTs.
//
// Java: CompilerTestCase.disableCompareJsDoc()
func (b *CompilerTestCaseBase) DisableCompareJsDoc() {
	b.checkSetUpRan()
	b.compareJsDoc = false
}

// EnableRunTypeCheckAfterProcessing moves type checking to occur after the
// processor, instead of before.
//
// Java: CompilerTestCase.enableRunTypeCheckAfterProcessing()
func (b *CompilerTestCaseBase) EnableRunTypeCheckAfterProcessing() {
	b.checkSetUpRan()
	b.runTypeCheckAfterProcessing = true
}

// GetNumRepetitions returns the number of times the pass should be run
// before results are verified.
//
// Java: CompilerTestCase.getNumRepetitions()
func (b *CompilerTestCaseBase) GetNumRepetitions() int {
	return 1
}

// IgnoreWarnings adds the given DiagnosticTypes (or, Java's other overload,
// DiagnosticGroups) to the set of warnings to ignore.
//
// Java: CompilerTestCase.ignoreWarnings(DiagnosticType...)
// Java: CompilerTestCase.ignoreWarnings(DiagnosticGroup...)
func (b *CompilerTestCaseBase) IgnoreWarnings(warnings ...any) {
	b.checkSetUpRan()
	for _, w := range warnings {
		switch w := w.(type) {
		case *DiagnosticType:
			b.ignoredWarnings.Add(w)
		case *DiagnosticGroup:
			b.ignoredWarnings.AddAll(w.GetTypes().All())
		default:
			panic(jutil.NewIllegalArgumentExceptionMsg(fmt.Sprintf("ignoreWarnings(%T)", w)))
		}
	}
}

// AllowSourcelessWarnings: expect warnings without source information.
//
// Java: CompilerTestCase.allowSourcelessWarnings()
func (b *CompilerTestCaseBase) AllowSourcelessWarnings() {
	b.checkSetUpRan()
	b.allowSourcelessWarnings = true
}

// GetLastCompiler is the most recently used JSComp instance.
//
// Java: CompilerTestCase.getLastCompiler()
func (b *CompilerTestCaseBase) GetLastCompiler() *Compiler {
	return b.lastCompiler
}

// SetAcceptedLanguage sets what language to allow in source parsing. Also
// sets the output language.
//
// Java: CompilerTestCase.setAcceptedLanguage(LanguageMode)
func (b *CompilerTestCaseBase) SetAcceptedLanguage(lang CompilerOptions_LanguageMode) {
	b.checkSetUpRan()
	b.SetLanguage(lang, lang)
}

// SetLanguage sets the input and output language modes..
//
// Java: CompilerTestCase.setLanguage(LanguageMode, LanguageMode)
func (b *CompilerTestCaseBase) SetLanguage(langIn, langOut CompilerOptions_LanguageMode) {
	b.acceptedLanguage = langIn
	b.SetLanguageOut(langOut)
}

// Java: CompilerTestCase.setLanguageOut(LanguageMode)
func (b *CompilerTestCaseBase) SetLanguageOut(langOut CompilerOptions_LanguageMode) {
	b.checkSetUpRan()
	b.languageOut = langOut
}

// SetBrowserFeaturesetYear takes the Java @Nullable Integer as *int.
//
// Java: CompilerTestCase.setBrowserFeaturesetYear(Integer)
func (b *CompilerTestCaseBase) SetBrowserFeaturesetYear(year *int) {
	b.checkSetUpRan()
	b.browserFeaturesetYear = year
}

// Java: CompilerTestCase.setModuleResolutionMode(ResolutionMode)
func (b *CompilerTestCaseBase) SetModuleResolutionMode(moduleResolutionMode ModuleLoader_ResolutionMode) {
	b.checkSetUpRan()
	b.moduleResolutionMode = moduleResolutionMode
}

// Java: CompilerTestCase.setJsDocumentation(JsDocParsing)
func (b *CompilerTestCaseBase) SetJsDocumentation(parseJsDocDocumentation parsing.Config_JsDocParsing) {
	b.checkSetUpRan()
	b.parseJsDocDocumentation = parseJsDocDocumentation
}

// Java: CompilerTestCase.setAssumeStaticInheritanceIsNotUsed(boolean)
func (b *CompilerTestCaseBase) SetAssumeStaticInheritanceIsNotUsed(assumeStaticInheritanceIsNotUsed bool) {
	b.checkSetUpRan()
	b.assumeStaticInheritanceIsNotUsed = assumeStaticInheritanceIsNotUsed
}

// EnableInferConsts sets whether to run InferConsts before passes
//
// Java: CompilerTestCase.enableInferConsts()
func (b *CompilerTestCaseBase) EnableInferConsts() {
	b.checkSetUpRan()
	b.inferConsts = true
}

// EnableCheckAccessControls enables running CheckAccessControls after the
// pass being tested (and checking types).
//
// Java: CompilerTestCase.enableCheckAccessControls()
func (b *CompilerTestCaseBase) EnableCheckAccessControls() {
	b.checkSetUpRan()
	b.checkAccessControls = true
}

// AllowExternsChanges allows externs to change.
//
// Java: CompilerTestCase.allowExternsChanges()
func (b *CompilerTestCaseBase) AllowExternsChanges() {
	b.checkSetUpRan()
	b.allowExternsChanges = true
}

// Java: CompilerTestCase.enablePolymerPass()
func (b *CompilerTestCaseBase) EnablePolymerPass() {
	b.checkSetUpRan()
	b.polymerPass = true
}

// EnableTypeCheck performs type checking before running the test pass. This
// will check for type errors and annotate nodes with type information.
//
// See TypeCheck.
//
// Java: CompilerTestCase.enableTypeCheck()
func (b *CompilerTestCaseBase) EnableTypeCheck() {
	b.checkSetUpRan()
	b.typeCheckEnabled = true
	b.EnableCreateModuleMap()
}

// DisableTypeCheck: do not run type checking before running the test pass.
//
// See TypeCheck.
//
// Java: CompilerTestCase.disableTypeCheck()
func (b *CompilerTestCaseBase) DisableTypeCheck() {
	b.checkSetUpRan()
	b.typeCheckEnabled = false
}

// ReplaceTypesWithColors converts any inferred JSTypes attached to Nodes to
// Colors, also deleting any JSTypes or JSDocInfo references.
//
// Java: CompilerTestCase.replaceTypesWithColors()
func (b *CompilerTestCaseBase) ReplaceTypesWithColors() {
	b.checkSetUpRan()
	b.replaceTypesWithColors = true
}

// Java: CompilerTestCase.enableCreateModuleMap()
func (b *CompilerTestCaseBase) EnableCreateModuleMap() {
	b.checkSetUpRan()
	b.createModuleMap = true
}

// EnableMultistageCompilation runs using multistage compilation.
//
// Java: CompilerTestCase.enableMultistageCompilation()
func (b *CompilerTestCaseBase) EnableMultistageCompilation() {
	b.checkSetUpRan()
	b.multistageCompilation = true
}

// DisableMultistageCompilation runs using singlestage compilation.
//
// Java: CompilerTestCase.disableMultistageCompilation()
func (b *CompilerTestCaseBase) DisableMultistageCompilation() {
	b.checkSetUpRan()
	b.multistageCompilation = false
}

// DisableValidateAstChangeMarking disables validating AST change marking.
//
// Java: CompilerTestCase.disableValidateAstChangeMarking()
func (b *CompilerTestCaseBase) DisableValidateAstChangeMarking() {
	b.checkSetUpRan()
	b.checkAstChangeMarking = false
}

// EnableClosurePass processes closure library primitives.
//
// Java: CompilerTestCase.enableClosurePass()
func (b *CompilerTestCaseBase) EnableClosurePass() {
	b.checkSetUpRan()
	b.closurePassEnabled = true
	b.EnableCreateModuleMap() // ProcessClosurePrimitives looks at the ModuleMetadataMap
}

// Java: CompilerTestCase.enableClosurePassForExpected()
func (b *CompilerTestCaseBase) EnableClosurePassForExpected() {
	b.checkSetUpRan()
	b.closurePassEnabledForExpected = true
}

// EnableProcessCommonJsModules rewrites CommonJS modules before the test
// run.
//
// Java: CompilerTestCase.enableProcessCommonJsModules()
func (b *CompilerTestCaseBase) EnableProcessCommonJsModules() {
	b.checkSetUpRan()
	b.processCommonJsModules = true
}

// EnableRewriteClosureCode rewrites Closure code before the test is run.
//
// Java: CompilerTestCase.enableRewriteClosureCode()
func (b *CompilerTestCaseBase) EnableRewriteClosureCode() {
	b.checkSetUpRan()
	b.rewriteClosureCode = true
	b.EnableCreateModuleMap()
}

// DisableRewriteClosureCode: don't rewrite Closure code before the test is
// run.
//
// Java: CompilerTestCase.disableRewriteClosureCode()
func (b *CompilerTestCaseBase) DisableRewriteClosureCode() {
	b.checkSetUpRan()
	b.rewriteClosureCode = false
}

// Java: CompilerTestCase.enableRewriteModulesAfterTypechecking()
func (b *CompilerTestCaseBase) EnableRewriteModulesAfterTypechecking() {
	b.checkSetUpRan()
	b.EnableRewriteClosureCode()
	b.rewriteModulesAfterTypechecking = true
}

// Java: CompilerTestCase.disableRewriteModulesAfterTypechecking()
func (b *CompilerTestCaseBase) DisableRewriteModulesAfterTypechecking() {
	b.checkSetUpRan()
	b.rewriteModulesAfterTypechecking = false
}

// EnableRewriteClosureProvides rewrites goog.provides
//
// Java: CompilerTestCase.enableRewriteClosureProvides()
func (b *CompilerTestCaseBase) EnableRewriteClosureProvides() {
	b.checkSetUpRan()
	b.rewriteClosureProvides = true
	b.EnableCreateModuleMap()
}

// EnableNormalize performs AST normalization before running the test pass
//
// This option also enables multistage compilation to increase coverage for
// multistage builds (Normalize is stage 2, so normalized tests must cover
// stage 2 passes). Disable via disableMultistageCompilation().
//
// See Normalize.
//
// Java: CompilerTestCase.enableNormalize()
func (b *CompilerTestCaseBase) EnableNormalize() {
	b.checkSetUpRan()
	b.normalizeEnabled = true
	b.EnableMultistageCompilation()
}

// EnableNormalizeExpectedOutput: if normalization is enabled for test
// inputs, then also run normalization on the expected outputs.
//
// Enabling this can cause the results of test failures to be confusing,
// since the expected output written in the test case is modified before it
// is compared with the test results. So, you should avoid using this when
// you can. However, there are some ways that Normalization modifies the AST
// (e.g. adding properties to nodes) that cannot be represented in input JS
// source text but will affect whether the AST comparison succeeds. This
// option is intended to allow Normalization to make those changes on the
// expected output code.
//
// See EnableNormalize.
//
// Java: CompilerTestCase.enableNormalizeExpectedOutput()
func (b *CompilerTestCaseBase) EnableNormalizeExpectedOutput() {
	b.checkSetUpRan()
	jutil.CheckStateMsg(b.normalizeEnabled, "Enabled normalize on output, but not input.")
	b.normalizeExpectedOutputEnabled = true
}

// EnableTranspile performs AST transpilation before running the test pass
// and also on the expected output.
//
// This also has the side effect of forcing normalization of the expected
// output if you also enable normalization.
//
// Although this method is not deprecated, it is only used in a few very
// special cases. The test behavior this creates can be confusing,
// especially since it modifies the expected output, so that diffs will not
// necessarily reflect the expected output text written directly into the
// test cases. Adding new usages of this is probably a bad idea.
//
// Java: CompilerTestCase.enableTranspile()
func (b *CompilerTestCaseBase) EnableTranspile() {
	b.checkSetUpRan()
	b.transpileEnabled = true
}

// DisableNormalize: don't perform AST normalization before running the test
// pass.
//
// See Normalize.
//
// Java: CompilerTestCase.disableNormalize()
func (b *CompilerTestCaseBase) DisableNormalize() {
	b.checkSetUpRan()
	b.normalizeEnabled = false
}

// EnableRewriteEsModules performs ES module transpilation before running the
// test pass.
//
// Java: CompilerTestCase.enableRewriteEsModules()
func (b *CompilerTestCaseBase) EnableRewriteEsModules() {
	b.checkSetUpRan()
	b.rewriteEsModulesEnabled = true
}

// EnableComputeSideEffects runs the PureFunctionIdentifier pass before
// running the test pass.
//
// See PureFunctionIdentifier.
//
// Java: CompilerTestCase.enableComputeSideEffects()
func (b *CompilerTestCaseBase) EnableComputeSideEffects() {
	b.checkSetUpRan()
	b.computeSideEffects = true
}

// DisableComputeSideEffects disables runng the PureFunctionIdentifier pass
// before running the test pass.
//
// See PureFunctionIdentifier.
//
// Java: CompilerTestCase.disableComputeSideEffects()
func (b *CompilerTestCaseBase) DisableComputeSideEffects() {
	b.checkSetUpRan()
	b.computeSideEffects = false
}

// Java: CompilerTestCase.enableSourceInformationAnnotator()
func (b *CompilerTestCaseBase) EnableSourceInformationAnnotator() {
	b.checkSetUpRan()
	b.annotateSourceInfo = true
}

// DeclareAccessor declares an accessor as being present "somewhere" in a the
// test program.
//
// Calling this method also disabled automatic getter / setter collection.
//
// See GatherGetterAndSetterProperties.
//
// Java: CompilerTestCase.declareAccessor(String, PropertyAccessKind)
func (b *CompilerTestCaseBase) DeclareAccessor(name string, kind AccessorSummary_PropertyAccessKind) {
	b.checkSetUpRan()
	jutil.CheckStateMsg(kind.HasGetterOrSetter(), "Kind should be a getter and/or setter.")
	b.declaredAccessors.Put(name, kind)
}

// EnableGatherExternProperties scans externs for properties that should not
// be renamed.
//
// Java: CompilerTestCase.enableGatherExternProperties()
func (b *CompilerTestCaseBase) EnableGatherExternProperties() {
	b.checkSetUpRan()
	b.gatherExternPropertiesEnabled = true
}

// Java: CompilerTestCase.getGatheredExternProperties()
func (b *CompilerTestCaseBase) GetGatheredExternProperties() *jutil.LinkedHashSet[string] {
	jutil.CheckStateMsg(b.gatherExternPropertiesEnabled, "Must enable gatherExternProperties")
	return b.lastCompiler.GetExternProperties()
}

// DisableAstValidation disables validating the AST after each run of the
// pass.
//
// Java: CompilerTestCase.disableAstValidation()
func (b *CompilerTestCaseBase) DisableAstValidation() {
	b.checkSetUpRan()
	b.astValidationEnabled = false
}

// EnableTypeInfoValidation enables validating type information in the AST
// after each run of the pass.
//
// Java: CompilerTestCase.enableTypeInfoValidation()
func (b *CompilerTestCaseBase) EnableTypeInfoValidation() {
	b.checkSetUpRan()
	b.typeInfoValidationEnabled = true
}

// DisableTypeInfoValidation disables validating type information in the AST
// after each run of the pass.
//
// Java: CompilerTestCase.disableTypeInfoValidation()
func (b *CompilerTestCaseBase) DisableTypeInfoValidation() {
	b.checkSetUpRan()
	b.typeInfoValidationEnabled = false
}

// EnableScriptFeatureValidation enables validating script featuresets after
// each run of the pass.
//
// Java: CompilerTestCase.enableScriptFeatureValidation()
func (b *CompilerTestCaseBase) EnableScriptFeatureValidation() {
	b.checkSetUpRan()
	b.scriptFeatureValidationEnabled = true
}

// DisableScriptFeatureValidation disables validating script featuresets
// after each run of the pass.
//
// Java: CompilerTestCase.disableScriptFeatureValidation()
func (b *CompilerTestCaseBase) DisableScriptFeatureValidation() {
	b.checkSetUpRan()
	b.scriptFeatureValidationEnabled = false
}

// DisableCompareAsTree disables comparing the expected output as a tree or
// string. 99% of the time you want to compare as a tree. There are a few
// special cases where you don't, like if you want to test the code printing
// of "unnatural" syntax trees. For example,
//
//	IF
//	  IF
//	    STATEMENT
//	ELSE
//	  STATEMENT
//
// Java: CompilerTestCase.disableCompareAsTree()
func (b *CompilerTestCaseBase) DisableCompareAsTree() {
	b.checkSetUpRan()
	b.compareAsTree = false
}

// SetExpectParseWarningsInThisTest sets whether we should ignore parse
// warnings for the current test method.
//
// Java: CompilerTestCase.setExpectParseWarningsInThisTest()
func (b *CompilerTestCaseBase) SetExpectParseWarningsInThisTest() {
	b.checkSetUpRan()
	b.expectParseWarningsInThisTest = true
}

// Java: CompilerTestCase.setWebpackModulesById(Map<String, String>)
func (b *CompilerTestCaseBase) SetWebpackModulesById(webpackModulesById *jutil.LinkedHashMap[string, string]) {
	b.webpackModulesById.Clear()
	for k, v := range webpackModulesById.All() {
		b.webpackModulesById.Put(k, v)
	}
}

// SetGenericNameReplacements sets a map of replacement generic names to the
// prefix of the original generated one.
//
// If set, used in conjunction with
// UnitTestUtils.updateGenericVarNamesInExpectedFiles on expected output.
//
// Note: Requires that sources are a FlatSources.
//
// Java: CompilerTestCase.setGenericNameReplacements(Map<String, String>)
func (b *CompilerTestCaseBase) SetGenericNameReplacements(genericNameReplacements *jutil.LinkedHashMap[string, string]) {
	b.genericNameReplacements = jutil.NewLinkedHashMapFrom(genericNameReplacements)
}

// Java: CompilerTestCase.disableGenericNameReplacements()
func (b *CompilerTestCaseBase) DisableGenericNameReplacements() {
	b.genericNameReplacements = jutil.NewLinkedHashMap[string, string]()
}

// compilerTestCase_createTypeCheck returns a newly created TypeCheck.
//
// Java: CompilerTestCase.createTypeCheck(Compiler)
func compilerTestCase_createTypeCheck(compiler *Compiler) *TypeCheck {
	var rai ReverseAbstractInterpreter = NewSemanticReverseAbstractInterpreter(compiler.GetTypeRegistry())
	typeChecker := NewTypeCheck(compiler, rai, compiler.GetTypeRegistry())
	compiler.SetTypeCheckingHasRun(true)
	return typeChecker
}

// EnsureLibraryInjected ensures the given library is injected before
// typechecking
//
// Java: CompilerTestCase.ensureLibraryInjected(String)
func (b *CompilerTestCaseBase) EnsureLibraryInjected(resourceName string) {
	b.librariesToInject.Add(resourceName)
}

// Test is every overload of Java's test(...), chosen from the argument
// types:
//
//	test(String js, String expected)
//	    Verifies that the compiler pass's JS output matches the expected
//	    output.
//	test(String js, String expected, Diagnostic diagnostic)
//	    Verifies that the compiler pass's JS output matches the expected
//	    output and (optionally) that an expected warning is issued. Or, if
//	    an error is expected, this method just verifies that the error is
//	    encountered.
//	test(TestPart... parts)
//
// Java: CompilerTestCase.test(String, String)
// Java: CompilerTestCase.test(String, String, Diagnostic)
// Java: CompilerTestCase.test(TestPart...)
func (b *CompilerTestCaseBase) Test(args ...any) {
	switch {
	case compilerTestCase_argsAre(args, "string", "string"):
		js, expected := args[0].(string), args[1].(string)
		b.self.Test(CompilerTestCase_Srcs(js), CompilerTestCase_Expected_(expected))
	case compilerTestCase_argsAre(args, "string", "string", "Diagnostic"):
		js, expected, diagnostic := args[0].(string), args[1].(string), args[2].(*CompilerTestCase_Diagnostic)
		b.self.Test(CompilerTestCase_Externs_(b.defaultExternsInputs), CompilerTestCase_Srcs(js), CompilerTestCase_Expected_(expected), diagnostic)
	default:
		b.testInternalParts(compilerTestCase_testParts("test", args))
	}
}

// TestError is every overload of Java's testError(...):
//
//	testError(String js, DiagnosticType error)
//	    Verifies that the compiler generates the given error for the given
//	    input.
//	testError(String js, DiagnosticType error, String description)
//	    Verifies that the compiler generates the given error and description
//	    for the given input.
//	testError(Sources srcs, Diagnostic error)
//	testError(Externs externs, Sources srcs, Diagnostic error)
//	testError(Sources js, DiagnosticType error)
//
// Java: CompilerTestCase.testError(String, DiagnosticType)
// Java: CompilerTestCase.testError(String, DiagnosticType, String)
// Java: CompilerTestCase.testError(Sources, Diagnostic)
// Java: CompilerTestCase.testError(Externs, Sources, Diagnostic)
// Java: CompilerTestCase.testError(Sources, DiagnosticType)
func (b *CompilerTestCaseBase) TestError(args ...any) {
	switch {
	case compilerTestCase_argsAre(args, "string", "DiagnosticType"):
		js, error_ := args[0].(string), args[1].(*DiagnosticType)
		b.self.Test(CompilerTestCase_Srcs(js), CompilerTestCase_Error(error_))
	case compilerTestCase_argsAre(args, "string", "DiagnosticType", "string"):
		js, error_, description := args[0].(string), args[1].(*DiagnosticType), args[2].(string)
		AssertThat(error_).IsNotNull()
		b.self.Test(CompilerTestCase_Srcs(js), CompilerTestCase_Error(error_).WithMessage(description))
	case compilerTestCase_argsAre(args, "Sources", "Diagnostic"):
		srcs, error_ := args[0].(CompilerTestCase_Sources), args[1].(*CompilerTestCase_Diagnostic)
		AssertThat(error_.level).IsEqualTo(CheckLevel_ERROR)
		b.self.Test(srcs, error_)
	case compilerTestCase_argsAre(args, "Externs", "Sources", "Diagnostic"):
		externs, srcs, error_ := args[0].(*CompilerTestCase_Externs), args[1].(CompilerTestCase_Sources), args[2].(*CompilerTestCase_Diagnostic)
		AssertThat(error_.level).IsEqualTo(CheckLevel_ERROR)
		b.self.Test(externs, srcs, error_)
	case compilerTestCase_argsAre(args, "Sources", "DiagnosticType"):
		js, error_ := args[0].(CompilerTestCase_Sources), args[1].(*DiagnosticType)
		AssertThat(error_).IsNotNull()
		b.self.Test(js, CompilerTestCase_Error(error_))
	default:
		compilerTestCase_noOverload("testError", args)
	}
}

// TestWarning is every overload of Java's testWarning(...):
//
//	testWarning(String js, DiagnosticType warning)
//	    Verifies that the compiler generates the given warning for the given
//	    input.
//	testWarning(Sources srcs, Diagnostic warning)
//	testWarning(Externs externs, Sources srcs, Diagnostic warning)
//	testWarning(Sources srcs, DiagnosticType warning)
//	testWarning(String js, DiagnosticType warning, String description)
//	    Verifies that the compiler generates the given warning and
//	    description for the given input.
//	testWarning(Externs externs, Sources srcs, DiagnosticType warning, String description)
//
// Java: CompilerTestCase.testWarning(String, DiagnosticType)
// Java: CompilerTestCase.testWarning(Sources, Diagnostic)
// Java: CompilerTestCase.testWarning(Externs, Sources, Diagnostic)
// Java: CompilerTestCase.testWarning(Sources, DiagnosticType)
// Java: CompilerTestCase.testWarning(String, DiagnosticType, String)
// Java: CompilerTestCase.testWarning(Externs, Sources, DiagnosticType, String)
func (b *CompilerTestCaseBase) TestWarning(args ...any) {
	switch {
	case compilerTestCase_argsAre(args, "string", "DiagnosticType"):
		js, warning := args[0].(string), args[1].(*DiagnosticType)
		AssertThat(warning).IsNotNull()
		b.self.Test(CompilerTestCase_Srcs(js), CompilerTestCase_Warning(warning))
	case compilerTestCase_argsAre(args, "Sources", "Diagnostic"):
		srcs, warning := args[0].(CompilerTestCase_Sources), args[1].(*CompilerTestCase_Diagnostic)
		AssertThat(warning.level).IsEqualTo(CheckLevel_WARNING)
		b.self.Test(srcs, warning)
	case compilerTestCase_argsAre(args, "Externs", "Sources", "Diagnostic"):
		externs, srcs, warning := args[0].(*CompilerTestCase_Externs), args[1].(CompilerTestCase_Sources), args[2].(*CompilerTestCase_Diagnostic)
		AssertThat(warning.level).IsEqualTo(CheckLevel_WARNING)
		b.self.Test(externs, srcs, warning)
	case compilerTestCase_argsAre(args, "Sources", "DiagnosticType"):
		srcs, warning := args[0].(CompilerTestCase_Sources), args[1].(*DiagnosticType)
		AssertThat(warning).IsNotNull()
		b.self.Test(srcs, CompilerTestCase_Warning(warning))
	case compilerTestCase_argsAre(args, "string", "DiagnosticType", "string"):
		js, warning, description := args[0].(string), args[1].(*DiagnosticType), args[2].(string)
		AssertThat(warning).IsNotNull()
		b.self.Test(CompilerTestCase_Srcs(js), CompilerTestCase_Warning(warning).WithMessage(jutil.StringTrim(description)))
	case compilerTestCase_argsAre(args, "Externs", "Sources", "DiagnosticType", "string"):
		externs, srcs, warning, description := args[0].(*CompilerTestCase_Externs), args[1].(CompilerTestCase_Sources), args[2].(*DiagnosticType), args[3].(string)
		AssertThat(warning).IsNotNull()
		b.self.Test(externs, srcs, CompilerTestCase_Warning(warning).WithMessage(jutil.StringTrim(description)))
	default:
		compilerTestCase_noOverload("testWarning", args)
	}
}

// TestNoWarning is every overload of Java's testNoWarning(...): verifies
// that the compiler generates no warnings for the given input.
//
// Java: CompilerTestCase.testNoWarning(String)
// Java: CompilerTestCase.testNoWarning(Sources)
// Java: CompilerTestCase.testNoWarning(Externs, Sources)
func (b *CompilerTestCaseBase) TestNoWarning(args ...any) {
	switch {
	case compilerTestCase_argsAre(args, "string"):
		b.self.Test(CompilerTestCase_Srcs(args[0].(string)))
	case compilerTestCase_argsAre(args, "Sources"):
		b.self.Test(args[0].(CompilerTestCase_Sources))
	case compilerTestCase_argsAre(args, "Externs", "Sources"):
		b.self.Test(args[0].(*CompilerTestCase_Externs), args[1].(CompilerTestCase_Sources))
	default:
		compilerTestCase_noOverload("testNoWarning", args)
	}
}

// TestInternal runs one test configuration.
//
// Java: CompilerTestCase.testInternal(Externs, Sources, Expected, List<Diagnostic>, List<Postcondition>)
func (b *CompilerTestCaseBase) TestInternal(
	externs *CompilerTestCase_Externs,
	inputs CompilerTestCase_Sources,
	expected *CompilerTestCase_Expected,
	diagnostics []*CompilerTestCase_Diagnostic,
	postconditions []CompilerTestCase_Postcondition) {

	compiler := b.self.CreateAndInitializeCompiler(externs, inputs)
	b.lastCompiler = compiler

	if b.typeCheckEnabled {
		rhinotesting.BaseJSTypeTestCase_AddNativeProperties(compiler.GetTypeRegistry())
	}

	b.testInternalWithCompiler(compiler, externs, inputs, expected, diagnostics, postconditions)
}

// CreateAndInitializeCompiler creates a fresh Compiler object initialized
// with the given externs and inputs.
//
// Important: This method does not affect the return value of
// `getLastCompiler()`.
//
// Java: CompilerTestCase.createAndInitializeCompiler(Externs, Sources)
func (b *CompilerTestCaseBase) CreateAndInitializeCompiler(externs *CompilerTestCase_Externs, inputs CompilerTestCase_Sources) *Compiler {
	compiler := b.self.CreateCompiler()

	options := b.self.GetOptions()

	if flatSources, ok := inputs.(*CompilerTestCase_FlatSources); ok {
		// TODO(bradfordcsmith): Why do we set this only for the non-chunk case?
		//     I extracted this method from testInternal().
		options.SetCheckTypes(b.parseTypeInfo || b.typeCheckEnabled)
		compiler.Init(externs.externs, flatSources.sources, options)
	} else {
		compiler.InitChunks(externs.externs, inputs.(*CompilerTestCase_ChunkSources).chunks, b.self.GetOptions())
	}
	return compiler
}

// Java: CompilerTestCase.maybeCreateSources(String, String)
func compilerTestCase_maybeCreateSources(name string, srcText *string) []*SourceFile {
	if srcText != nil {
		return []*SourceFile{SourceFile_FromCode(name, *srcText)}
	}
	return nil
}

// CompilerTestCase_CreateSources names each source name+index.
//
// Java: CompilerTestCase.createSources(String, String...)
// Java: CompilerTestCase.createSources(String, List<String>)
func CompilerTestCase_CreateSources(name string, sources []string) []*SourceFile {
	if sources == nil {
		return nil
	}
	expectedSources := []*SourceFile{}
	for i := 0; i < len(sources); i++ {
		expectedSources = append(expectedSources, SourceFile_FromCode(name+fmt.Sprint(i), sources[i]))
	}
	return expectedSources
}

// TestSame is every overload of Java's testSame(...):
//
//	testSame(String js)
//	    Verifies that the compiler pass's JS output is the same as its
//	    input.
//	testSame(String js, DiagnosticType warning)
//	    Verifies that the compiler pass's JS output is the same as its input
//	    and (optionally) that an expected warning is issued.
//	testSame(TestPart... parts)
//
// Java: CompilerTestCase.testSame(String)
// Java: CompilerTestCase.testSame(String, DiagnosticType)
// Java: CompilerTestCase.testSame(TestPart...)
func (b *CompilerTestCaseBase) TestSame(args ...any) {
	switch {
	case compilerTestCase_argsAre(args, "string"):
		js := args[0].(string)
		b.self.Test(CompilerTestCase_Srcs(js), CompilerTestCase_Expected_(js))
	case compilerTestCase_argsAre(args, "string", "DiagnosticType"):
		js, warning := args[0].(string), args[1].(*DiagnosticType)
		b.self.Test(CompilerTestCase_Srcs(js), CompilerTestCase_Expected_(js), CompilerTestCase_Warning(warning))
	default:
		parts := compilerTestCase_testParts("testSame", args)
		b.testInternalParts(append(parts, compilerTestCase_EXPECTED_SAME))
	}
}

// testInternalWithCompiler verifies that the compiler pass's JS output
// matches the expected output and (optionally) that an expected warning is
// issued. Or, if an error is expected, this method just verifies that the
// error is encountered.
//
// compiler: A compiler that has been initialized via Compiler.init;
// inputsObj: Input source files; expectedObj: Expected outputs;
// diagnostics: Expected warning/error diagnostics.
//
// Java: CompilerTestCase.testInternal(Compiler, Externs, Sources, Expected, List<Diagnostic>, List<Postcondition>)
func (b *CompilerTestCaseBase) testInternalWithCompiler(
	compiler *Compiler,
	externs *CompilerTestCase_Externs,
	inputsObj CompilerTestCase_Sources, // TODO remove this parameter
	expectedObj *CompilerTestCase_Expected,
	diagnostics []*CompilerTestCase_Diagnostic,
	postconditions []CompilerTestCase_Postcondition) {
	genericNameMapping := jutil.NewLinkedHashMap[string, string]()
	if flatSources, ok := inputsObj.(*CompilerTestCase_FlatSources); !b.genericNameReplacements.IsEmpty() &&
		expectedObj != nil &&
		expectedObj.expected != nil &&
		ok {
		expectedObj = CompilerTestCase_Expected_(
			UnitTestUtils_UpdateGenericVarNamesInExpectedFilesWithMapping(
				flatSources, expectedObj, b.genericNameReplacements, genericNameMapping))
	}

	var inputs []*SourceFile
	if flatSources, ok := inputsObj.(*CompilerTestCase_FlatSources); ok {
		inputs = flatSources.sources
	}
	var expected []*SourceFile
	if expectedObj != nil {
		expected = expectedObj.expected
	}
	var expectedErrors, expectedWarnings []*CompilerTestCase_Diagnostic
	for _, d := range diagnostics {
		if d.level == CheckLevel_ERROR {
			expectedErrors = append(expectedErrors, d)
		}
	}
	for _, d := range diagnostics {
		if d.level == CheckLevel_WARNING {
			expectedWarnings = append(expectedWarnings, d)
		}
	}
	jutil.CheckStateMsg(
		len(expectedErrors) == 0 || len(expectedWarnings) == 0,
		"Cannot expect both errors and warnings.")
	jutil.CheckStateMsg(
		len(expectedErrors) == 0 || expected == nil,
		"Cannot expect both errors and compiled output.")
	jutil.CheckStateMsg(b.setUpRan, "CompilerTestCase.setUp not run: call super.setUp() from overrides.")
	recentChange := NewRecentChange()
	compiler.GetChangeTracker().AddChangeHandler(recentChange)

	root := compiler.ParseInputs()

	if root == nil {
		// Might be an expected parse error.
		AssertWithMessage("parse errors").
			That(compiler.GetErrors()).
			ComparingElementsUsing(compilerTestCase_DIAGNOSTIC_CORRESPONDENCE).
			ContainsExactlyElementsIn(expectedErrors)
		for _, postcondition := range postconditions {
			postcondition.Verify(compiler)
		}
		return
	}
	if !b.expectParseWarningsInThisTest {
		AssertWithMessage("Unexpected parser warning(s)").That(compiler.GetWarnings()).IsEmpty()
	} else {
		AssertThat(compiler.GetWarningCount()).IsGreaterThan(0)
	}
	externsRoot := root.GetFirstChild()
	mainRoot := root.GetLastChild()

	if b.annotateSourceInfo {
		NodeTraversal_Traverse(compiler, mainRoot, SourceInformationAnnotator_Create())
	}

	if b.createModuleMap {
		NewGatherModuleMetadata(
			compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER).
			Process(externsRoot, mainRoot)
		NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).
			Process(externsRoot, mainRoot)
	}

	if b.astValidationEnabled {
		// NOTE: We do not enable type validation here, because type information never exists
		// immediately after parsing.
		NewAstValidatorWithValidateScriptFeatures(compiler, b.scriptFeatureValidationEnabled).ValidateRoot(root)
	}

	// Save the tree for later comparison.
	rootClone := root.CloneTree()
	externsRootClone := rootClone.GetFirstChild()
	mainRootClone := rootClone.GetLastChild()

	numRepetitions := b.self.GetNumRepetitions()
	errorManagers := make([]ErrorManager, numRepetitions)
	aggregateWarnings := []*JSError{}
	hasCodeChanged := false

	for i := 0; i < numRepetitions; i++ {
		if compiler.GetErrorCount() == 0 {
			errorManagers[i] = NewBlackHoleErrorManager()
			compiler.SetErrorManager(errorManagers[i])

			if b.polymerPass && i == 0 {
				recentChange.Reset()
				NewPolymerPass(compiler).Process(externsRoot, mainRoot)
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			}

			if b.rewriteClosureCode && i == 0 {
				NewCheckClosureImports(compiler, compiler.GetModuleMetadataMap()).
					Process(externsRoot, mainRoot)
				ScopedAliases_Builder_(compiler).Build().Process(externsRoot, mainRoot)
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			}

			if b.rewriteClosureCode && !b.rewriteModulesAfterTypechecking && i == 0 {
				NewClosureRewriteModule(compiler, nil, nil).Process(externsRoot, mainRoot)
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			}

			// Only run process closure primitives once, if asked.
			if b.closurePassEnabled && i == 0 {
				recentChange.Reset()
				NewProcessClosurePrimitives(compiler).Process(externsRoot, mainRoot)
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			}

			// Only run process closure primitives once, if asked.
			if b.processCommonJsModules && i == 0 {
				recentChange.Reset()
				NewProcessCommonJSModules(compiler).Process(externsRoot, mainRoot)
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			}

			if (b.rewriteEsModulesEnabled || b.transpileEnabled) && i == 0 {
				// If we're transpiling down to ES5, then that includes ES module transpilation,
				// but that pass isn't part of transpileToEs5(), because it happens earlier in
				// the compilation process than most transpilation passes.
				recentChange.Reset()
				compilerTestCase_rewriteEsModules(compiler, externsRoot, mainRoot)
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			}

			injectLibrariesFromTypedAsts := b.multistageCompilation || b.replaceTypesWithColors
			if !b.librariesToInject.IsEmpty() && i == 0 && !injectLibrariesFromTypedAsts {
				recentChange.Reset()
				for resourceName := range b.librariesToInject.All() {
					compiler.GetRuntimeJsLibManager().EnsureLibraryInjected(resourceName, true)
				}
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			}

			// Only run the type checking pass once, if asked.
			// Running it twice can cause unpredictable behavior because duplicate
			// objects for the same type are created, and the type system
			// uses reference equality to compare many types.
			if !b.runTypeCheckAfterProcessing && b.typeCheckEnabled && i == 0 {
				check := compilerTestCase_createTypeCheck(compiler)
				check.ProcessForTesting(externsRoot, mainRoot)
			}

			if b.rewriteClosureCode && b.rewriteModulesAfterTypechecking && i == 0 {
				NewClosureRewriteModule(compiler, nil, compiler.GetTopScope()).
					Process(externsRoot, mainRoot)
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			}

			if b.inferConsts && i == 0 {
				NewInferConsts(compiler).Process(externsRoot, mainRoot)
			}

			if b.gatherExternPropertiesEnabled && i == 0 {
				NewGatherExternProperties(compiler, GatherExternProperties_Mode_CHECK_AND_OPTIMIZE).
					Process(externsRoot, mainRoot)
			}

			if i == 0 {
				if b.multistageCompilation {
					if inputs != nil {
						recentChange.Reset()
						compiler = CompilerTestCaseUtils_MultistageSerializeAndDeserialize(
							b.self, compiler, externs.externs, inputs, recentChange)
						root = compiler.GetRoot()
						externsRoot = compiler.GetExternsRoot()
						mainRoot = compiler.GetJsRoot()
						b.lastCompiler = compiler
						hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
					}
				} else if b.replaceTypesWithColors {
					recentChange.Reset()
					NewRemoveCastNodes(compiler).Process(externsRoot, mainRoot)
					NewConvertTypesToColors(
						compiler, SerializationOptions_Builder_().SetIncludeDebugInfo(true).Build()).
						Process(externsRoot, mainRoot)

					compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC)
					hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
				}

				if !b.librariesToInject.IsEmpty() && injectLibrariesFromTypedAsts {
					recentChange.Reset()
					for resourceName := range b.librariesToInject.All() {
						compiler.GetRuntimeJsLibManager().EnsureLibraryInjected(resourceName, true)
					}
					hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
				}
			}

			// Only rewrite provides once, if asked.
			if b.rewriteClosureProvides && i == 0 {
				recentChange.Reset()
				NewProcessClosureProvidesAndRequires(compiler, false).Process(externsRoot, mainRoot)
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			}

			if b.transpileEnabled && i == 0 {
				recentChange.Reset()
				b.transpileToEs5(compiler, externsRoot, mainRoot)
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			} else if b.normalizeEnabled && i == 0 {
				// An explicit normalize pass is necessary since it wasn't done as part of transpilation.
				compilerTestCase_normalizeActualCode(compiler, externsRoot, mainRoot)
			}

			b.updateAccessorSummary(compiler, mainRoot)

			if b.computeSideEffects && i == 0 {
				recentChange.Reset()
				mark := NewPureFunctionIdentifier_Driver(compiler)
				mark.Process(externsRoot, mainRoot)
				hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			}

			recentChange.Reset()

			var changeVerifier *ChangeVerifier
			if b.checkAstChangeMarking {
				changeVerifier = NewChangeVerifier(compiler)
				changeVerifier.Snapshot(mainRoot)
			}

			// Call "beforePass" as some passes ask for the index of the current pass being run and
			// expect it to be >= 0.
			compiler.BeforePass(b.self.GetName())
			b.self.GetProcessor(compiler).Process(externsRoot, mainRoot)

			if b.checkAstChangeMarking {
				// TODO(johnlenz): add support for multiple passes in getProcessor so that we can
				// check the AST marking after each pass runs.
				// Verify that changes to the AST are properly marked on the AST.
				changeVerifier.CheckRecordedChanges(mainRoot)
			}

			b.verifyAccessorSummary(compiler, mainRoot)

			if b.runTypeCheckAfterProcessing && b.typeCheckEnabled && i == 0 {
				check := compilerTestCase_createTypeCheck(compiler)
				check.ProcessForTesting(externsRoot, mainRoot)
			}

			// Transpilation passes are allowed to leave the AST in a bad state when there is a halting
			// error.
			if b.astValidationEnabled && !compiler.HasHaltingErrors() {
				var typeValidationMode AstValidator_TypeInfoValidation
				if b.typeInfoValidationEnabled {
					if compiler.HasOptimizationColors() {
						typeValidationMode = AstValidator_TypeInfoValidation_COLOR
					} else {
						typeValidationMode = AstValidator_TypeInfoValidation_JSTYPE
					}
				} else {
					typeValidationMode = AstValidator_TypeInfoValidation_NONE
				}
				NewAstValidatorWithValidateScriptFeatures(compiler, b.scriptFeatureValidationEnabled).
					SetTypeValidationMode(typeValidationMode).
					ValidateRoot(root)
			}
			NewSourceInfoCheck(compiler).Process(externsRoot, mainRoot)

			if b.checkAccessControls {
				NewCheckAccessControls(compiler).Process(externsRoot, mainRoot)
			}

			hasCodeChanged = hasCodeChanged || recentChange.HasCodeChanged()
			aggregateWarnings = append(aggregateWarnings, compiler.GetWarnings()...)

			if b.normalizeEnabled {
				verifyDeclaredConstants := true
				NewValidityCheck_VerifyConstants(compiler, verifyDeclaredConstants).
					Process(externsRoot, mainRoot)
			}
		}
	}

	if len(expectedErrors) == 0 {
		AssertThat(compiler.GetErrors()).IsEmpty()

		b.validateWarnings(
			compiler, expectedWarnings, errorManagers, aggregateWarnings, numRepetitions)

		// If we ran normalize on the AST, we must also run normalize on the clone before checking for
		// changes.
		if b.normalizeEnabled {
			hasTypecheckingRun := compiler.HasTypeCheckingRun()
			// we don't run type inference over the clone of the AST, so need to mark that in the
			// compiler or Normalize will crash due to lack of inferred types on the clone AST nodes.
			compiler.SetTypeCheckingHasRun(false)
			compilerTestCase_normalizeActualCode(compiler, externsRootClone, mainRootClone)
			compiler.SetTypeCheckingHasRun(hasTypecheckingRun)
		}

		if b.checkAstChangeMarking || !b.allowExternsChanges {
			b.validateCodeChangeReporting(
				mainRoot, externsRoot, mainRootClone, externsRootClone, compiler, hasCodeChanged)
		}

		if expected != nil {
			expectedRoot := b.parseExpectedJsList(expected)
			if b.compareAsTree {
				b.compareExpectedToActualAsTree(mainRoot, expectedRoot, compiler, genericNameMapping)
			} else {
				b.compareExpectedToActualAsStrings(mainRoot, compiler, expected)
			}
		}

		b.validateNormalizationInvariants(root, compiler)
	} else {
		AssertWithMessage("compile errors").
			That(compiler.GetErrors()).
			ComparingElementsUsing(compilerTestCase_DIAGNOSTIC_CORRESPONDENCE).
			ContainsExactlyElementsIn(expectedErrors)
		for _, error_ := range compiler.GetErrors() {
			b.validateSourceLocation(error_)
			AssertWithMessage("Some placeholders in the error message were not replaced").
				That(error_.Description()).
				DoesNotContainMatch("\\{\\d\\}")
		}
	}
	for _, postcondition := range postconditions {
		postcondition.Verify(compiler)
	}
}

// Java: CompilerTestCase.rewriteEsModules(AbstractCompiler, Node, Node)
func compilerTestCase_rewriteEsModules(compiler AbstractCompiler, externsRoot *rhino.Node, codeRoot *rhino.Node) {
	options := compiler.GetOptions()
	factories := NewPassListBuilder(options)

	gatherModuleMetadata := NewGatherModuleMetadata(
		compiler, options.GetProcessCommonJSModules(), options.GetModuleResolutionMode())
	factories.MaybeAdd(
		PassFactory_Builder_().
			SetName(PassNames_GATHER_MODULE_METADATA).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(func(x1 AbstractCompiler) CompilerPass { return gatherModuleMetadata }).
			Build())
	factories.MaybeAdd(
		PassFactory_Builder_().
			SetName(PassNames_CREATE_MODULE_MAP).
			SetRunInFixedPointLoop(true).
			SetInternalFactory(
				func(x AbstractCompiler) CompilerPass {
					return NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap())
				}).
			Build())
	TranspilationPasses_AddEs6ModulePass(
		factories, NewPreprocessorSymbolTable_CachedInstanceFactory())
	for _, factory := range factories.Build() {
		factory.Create(compiler).Process(externsRoot, codeRoot)
	}
}

// Java: CompilerTestCase.transpileToEs5(AbstractCompiler, Node, Node)
func (b *CompilerTestCaseBase) transpileToEs5(compiler AbstractCompiler, externsRoot *rhino.Node, codeRoot *rhino.Node) {
	options := compiler.GetOptions()
	factories := NewPassListBuilder(options)

	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	TranspilationPasses_AddTranspilationRuntimeLibraries(factories)
	TranspilationPasses_AddRewritePolyfillPass(factories)
	// Transpilation requires normalization.
	factories.MaybeAdd(
		PassFactory_Builder_().
			SetName(PassNames_NORMALIZE).
			SetInternalFactory(func(abstractCompiler AbstractCompiler) CompilerPass {
				return Normalize_Builder_(abstractCompiler).Build()
			}).
			Build())
	TranspilationPasses_AddTranspilationPasses(factories, options)
	// We need to put back the original variable names where possible once transpilation is
	// complete. This matches the behavior in DefaultPassConfig. See comments there for further
	// explanation.
	factories.MaybeAdd(
		PassFactory_Builder_().
			SetName("invertContextualRenaming").
			SetInternalFactory(MakeDeclaredNamesUnique_GetContextualRenameInverter).
			Build())
	for _, factory := range factories.Build() {
		factory.Create(compiler).Process(externsRoot, codeRoot)
	}
}

// Java: CompilerTestCase.validateWarnings(Compiler, List<Diagnostic>, ErrorManager[], List<JSError>, int)
func (b *CompilerTestCaseBase) validateWarnings(
	compiler *Compiler,
	expectedWarnings []*CompilerTestCase_Diagnostic,
	errorManagers []ErrorManager,
	aggregateWarnings []*JSError,
	numRepetitions int) {
	formatter := NewLightweightMessageFormatter(compiler)
	if len(expectedWarnings) == 0 {
		formatted := make([]string, len(aggregateWarnings))
		for i, w := range aggregateWarnings {
			formatted[i] = formatter.FormatWarning(w)
		}
		AssertWithMessage(
			"aggregate warnings: %s",
			strings.Join(formatted, "\n")).
			That(aggregateWarnings).
			IsEmpty()
	} else {
		AssertWithMessage(
			"There should be %s warnings, repeated %s time(s). Warnings: \n%s",
			len(expectedWarnings), numRepetitions, compilerTestCase_LINE_JOINER(aggregateWarnings)).
			That(len(aggregateWarnings)).
			IsEqualTo(numRepetitions * len(expectedWarnings))
		for i := 0; i < numRepetitions; i++ {
			AssertWithMessage("compile warnings from repetition %s", (i + 1)).
				That(errorManagers[i].GetWarnings()).
				ComparingElementsUsing(compilerTestCase_DIAGNOSTIC_CORRESPONDENCE).
				ContainsExactlyElementsIn(expectedWarnings)
			for _, warning := range errorManagers[i].GetWarnings() {
				b.validateSourceLocation(warning)
			}
		}
	}
}

// Java: CompilerTestCase.validateCodeChangeReporting(Node, Node, Node, Node, Compiler, boolean)
func (b *CompilerTestCaseBase) validateCodeChangeReporting(
	mainRoot *rhino.Node,
	externsRoot *rhino.Node,
	mainRootClone *rhino.Node,
	externsRootClone *rhino.Node,
	compiler *Compiler,
	wasCodeChangeReported bool) {
	// NOTE: We have to use DEEP_NO_SHADOW here because the code change reporting logic is,
	// deliberately, unaware of shadows ASTs.
	codeChange := !mainRootClone.IsEquivalentToWithModes(
		mainRoot,
		rhino.Node_RecursionMode_DEEP_NO_SHADOW,
		rhino.Node_TypeComparison_IGNORE,
		rhino.Node_JsDocComparison_IGNORE,
		rhino.Node_SideEffectComparison_COMPARE)
	externsChange := !externsRootClone.IsEquivalentToWithModes(
		externsRoot,
		rhino.Node_RecursionMode_DEEP_NO_SHADOW,
		rhino.Node_TypeComparison_IGNORE,
		rhino.Node_JsDocComparison_IGNORE,
		rhino.Node_SideEffectComparison_COMPARE)

	// Generally, externs should not be changed by the compiler passes.
	if externsChange && !b.allowExternsChanges {
		NodeSubject_AssertNode(externsRoot).
			UsingSerializer(compilerTestCase_createPrettyPrinter(compiler)).
			IsEqualTo(externsRootClone)
	}
	if b.checkAstChangeMarking {
		if !codeChange && !externsChange {
			AssertWithMessage("compiler.reportCodeChange() was called even though nothing changed").
				That(wasCodeChangeReported).
				IsFalse()
		} else {
			AssertWithMessage(
				"compiler.reportCodeChange() should have been called."+"\nOriginal: %s\nNew: %s",
				mainRootClone.ToStringTree(), mainRoot.ToStringTree()).
				That(wasCodeChangeReported).
				IsTrue()
		}
	}
}

// Java: CompilerTestCase.compareExpectedToActualAsTree(Node, Node, Compiler, Map<String, String>)
func (b *CompilerTestCaseBase) compareExpectedToActualAsTree(
	mainRoot *rhino.Node, expectedRoot *rhino.Node, compiler *Compiler, genericNameMapping *jutil.LinkedHashMap[string, string]) {
	if b.compareJsDoc {
		NodeSubject_AssertNode(mainRoot).
			UsingSerializer(compilerTestCase_createPrettyPrinter(compiler)).
			WithGenericNameReplacements(jutil.NewLinkedHashMapFrom(genericNameMapping)).
			IsEqualIncludingJsDocTo(expectedRoot)
	} else {
		NodeSubject_AssertNode(mainRoot).
			UsingSerializer(compilerTestCase_createPrettyPrinter(compiler)).
			WithGenericNameReplacements(jutil.NewLinkedHashMapFrom(genericNameMapping)).
			IsEqualTo(expectedRoot)
	}
}

// compilerTestCase_spacesBeforeNewline is the regex " +\n".
var compilerTestCase_spacesBeforeNewline = regexp.MustCompile(" +\n")

// Java: CompilerTestCase.compareExpectedToActualAsStrings(Node, Compiler, List<SourceFile>)
func (b *CompilerTestCaseBase) compareExpectedToActualAsStrings(
	mainRoot *rhino.Node, compiler *Compiler, expected []*SourceFile) {
	expectedSources := make([]string, len(expected))
	for i := 0; i < len(expected); i++ {
		expectedSources[i] = expected[i].GetCode()
	}
	AssertThat(compilerTestCase_spacesBeforeNewline.ReplaceAllString(compiler.ToSourceOfNode(mainRoot), "\n")).
		IsEqualTo(strings.Join(expectedSources, ""))
}

// Java: CompilerTestCase.validateNormalizationInvariants(Node, Compiler)
func (b *CompilerTestCaseBase) validateNormalizationInvariants(root *rhino.Node, compiler *Compiler) {
	mainRoot := root.GetLastChild()

	// Verify normalization is not invalidated.
	normalizeCheckRootClone := root.CloneTree()
	normalizeCheckExternsRootClone := normalizeCheckRootClone.GetFirstChild()
	normalizeCheckMainRootClone := normalizeCheckRootClone.GetLastChild()

	NodeSubject_AssertNode(normalizeCheckMainRootClone).
		UsingSerializer(compilerTestCase_createPrettyPrinter(compiler)).
		IsEqualTo(mainRoot)

	// TODO(johnlenz): enable this for most test cases.
	// Currently, this invalidates test for while-loops, for-loop
	// initializers, and other naming.  However, a set of code
	// (Closure primitive rewrites, etc) runs before the Normalize pass,
	// so this can't be force on everywhere.
	if b.normalizeEnabled {
		Normalize_Builder_(compiler).
			AssertOnChange(true).
			Build().
			Process(normalizeCheckExternsRootClone, normalizeCheckMainRootClone)

		NodeSubject_AssertNode(normalizeCheckMainRootClone).
			UsingSerializer(compilerTestCase_createPrettyPrinter(compiler)).
			IsEqualTo(mainRoot)
	}
}

// Java: CompilerTestCase.validateSourceLocation(JSError)
func (b *CompilerTestCaseBase) validateSourceLocation(jserror *JSError) {
	// Make sure that source information is always provided.
	if !b.allowSourcelessWarnings {
		sourceName := jserror.SourceName()
		AssertWithMessage("Missing source file name in warning: %s", jserror).
			That(sourceName != nil && *sourceName != "").
			IsTrue()
		scriptNode := b.lastCompiler.GetScriptNode(*sourceName)
		AssertWithMessage("No SCRIPT node found for warning: %s", jserror).
			That(scriptNode).
			IsNotNull()
		AssertWithMessage("Missing line number in warning: %s", jserror).
			That(jserror.Lineno() != -1).
			IsTrue()
		AssertWithMessage("Missing char number in warning: %s", jserror).
			That(jserror.Charno() != -1).
			IsTrue()
	}
}

// Java: CompilerTestCase.normalizeActualCode(Compiler, Node, Node)
func compilerTestCase_normalizeActualCode(compiler *Compiler, externsRoot *rhino.Node, mainRoot *rhino.Node) {
	normalize := Normalize_CreateNormalizeForOptimizations(compiler)
	normalize.Process(externsRoot, mainRoot)
}

// Java: CompilerTestCase.updateAccessorSummary(Compiler, Node)
func (b *CompilerTestCaseBase) updateAccessorSummary(compiler *Compiler, mainRoot *rhino.Node) {
	union := jutil.NewLinkedHashMapFrom(GatherGetterAndSetterProperties_Gather(compiler, mainRoot.GetParent()))
	for k, v := range b.declaredAccessors.All() {
		union.Put(k, v)
	}
	compiler.SetAccessorSummary(AccessorSummary_Create(union))
}

// Java: CompilerTestCase.verifyAccessorSummary(Compiler, Node)
func (b *CompilerTestCaseBase) verifyAccessorSummary(compiler *Compiler, mainRoot *rhino.Node) {
	if compiler.GetAccessorSummary() == nil {
		return
	}

	AssertWithMessage(
		"Pass created new getters/setters without calling" +
			" GatherGetterAndSetterProperties.update").
		That(compiler.GetAccessorSummary().GetAccessors()).
		ContainsAtLeastEntriesIn(
			GatherGetterAndSetterProperties_Gather(compiler, mainRoot.GetParent()))
}

// ParseExpectedJs parses the expected JS (a String, or an Expected) and
// returns the root of the parse tree.
//
// Java: CompilerTestCase.parseExpectedJs(String)
// Java: CompilerTestCase.parseExpectedJs(Expected)
func (b *CompilerTestCaseBase) ParseExpectedJs(expected any) *rhino.Node {
	switch expected := expected.(type) {
	case string:
		return b.parseExpectedJsList([]*SourceFile{SourceFile_FromCode(CompilerTestCase_GENERATED_SRC_NAME, expected)})
	case *CompilerTestCase_Expected:
		return b.parseExpectedJsList(expected.expected)
	}
	compilerTestCase_noOverload("parseExpectedJs", []any{expected})
	return nil
}

// parseExpectedJsList parses expected JS inputs and returns the root of the
// parse tree.
//
// Java: CompilerTestCase.parseExpectedJs(List<SourceFile>)
func (b *CompilerTestCaseBase) parseExpectedJsList(inputs []*SourceFile) *rhino.Node {
	compiler := b.self.CreateCompiler()

	compiler.Init(b.defaultExternsInputs, inputs, b.self.GetOptions())
	root := compiler.ParseInputs()
	AssertWithMessage("Unexpected parse error(s): %s", compilerTestCase_LINE_JOINER(compiler.GetErrors())).
		That(root).
		IsNotNull()
	externsRoot := root.GetFirstChild()
	mainRoot := externsRoot.GetNext()

	// Normalize CAST nodes from the expected AST if asked.
	if b.replaceTypesWithColors || b.multistageCompilation {
		NewRemoveCastNodes(compiler).Process(externsRoot, mainRoot)
	}

	if b.closurePassEnabled && b.closurePassEnabledForExpected && !compiler.HasErrors() {
		NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).
			Process(externsRoot, mainRoot)
		NewProcessClosurePrimitives(compiler).Process(externsRoot, mainRoot)
	}

	if b.rewriteClosureCode {
		NewClosureRewriteModule(compiler, nil, nil).Process(externsRoot, mainRoot)
		ScopedAliases_Builder_(compiler).Build().Process(externsRoot, mainRoot)
	}

	if b.rewriteClosureProvides && b.closurePassEnabledForExpected && !compiler.HasErrors() {
		NewProcessClosureProvidesAndRequires(compiler, false).Process(externsRoot, mainRoot)
	}

	if b.transpileEnabled && !compiler.HasErrors() {
		compilerTestCase_rewriteEsModules(compiler, externsRoot, mainRoot)
		// TODO(bradfordcsmith): Notice that transpileToEs5() will perform normalization, if that
		// is enabled, without checking normalizeExpectedOutputEnabled.
		// This behavior is documented in enableTranspile(), which is used only in a very few
		// special cases.
		b.transpileToEs5(compiler, externsRoot, mainRoot)
	} else if b.normalizeEnabled && b.normalizeExpectedOutputEnabled && !compiler.HasErrors() {
		// An explicit normalize pass is necessary since it wasn't done as part of transpilation.
		compilerTestCase_normalizeActualCode(compiler, externsRoot, mainRoot)
	}
	return mainRoot
}

// TestExternChanges is every overload of Java's testExternChanges(...):
//
//	testExternChanges(Sources input, Expected expectedExtern, Diagnostic... warnings)
//	testExternChanges(Externs externs, Sources inputs, Expected expectedExtern, Diagnostic... warnings)
//
// Java: CompilerTestCase.testExternChanges(Sources, Expected, Diagnostic...)
// Java: CompilerTestCase.testExternChanges(Externs, Sources, Expected, Diagnostic...)
func (b *CompilerTestCaseBase) TestExternChanges(args ...any) {
	if len(args) >= 1 {
		if externs, ok := args[0].(*CompilerTestCase_Externs); ok {
			inputs := args[1].(CompilerTestCase_Sources)
			expectedExtern := args[2].(*CompilerTestCase_Expected)
			warnings := compilerTestCase_diagnostics("testExternChanges", args[3:])
			compiler := b.self.CreateCompiler()
			options := b.self.GetOptions()
			if flatSources, ok := inputs.(*CompilerTestCase_FlatSources); ok {
				compiler.Init(externs.externs, flatSources.sources, options)
			} else {
				compiler.InitChunks(externs.externs, inputs.(*CompilerTestCase_ChunkSources).chunks, b.self.GetOptions())
			}
			b.testExternChangesInternal(compiler, expectedExtern, warnings)
			return
		}
	}
	if len(args) < 2 {
		compilerTestCase_noOverload("testExternChanges", args)
	}
	b.self.TestExternChanges(append([]any{CompilerTestCase_Externs_("")}, args...)...)
}

// Java: CompilerTestCase.testExternChangesInternal(Compiler, Expected, Diagnostic...)
func (b *CompilerTestCaseBase) testExternChangesInternal(
	compiler *Compiler, expectedExterns *CompilerTestCase_Expected, warnings []*CompilerTestCase_Diagnostic) {
	compiler.ParseInputs()

	if b.createModuleMap {
		NewGatherModuleMetadata(
			compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER).
			Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
		NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).
			Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	}
	AssertThat(compiler.GetErrors()).IsEmpty()

	externsAndJs := compiler.GetRoot()
	root := externsAndJs.GetLastChild()

	externs := externsAndJs.GetFirstChild()

	expectedRoot := b.self.ParseExpectedJs(expectedExterns)
	AssertThat(compiler.GetErrors()).IsEmpty()

	// Call "beforePass" as some passes ask for the index of the current pass being run and expect
	// it to be >= 0.
	compiler.BeforePass(b.self.GetName())
	(b.self.GetProcessor(compiler)).Process(externs, root)

	if b.compareAsTree {
		// Ignore and remove empty externs, so that if we start with an empty extern and only add
		// to the synthetic externs, we can still enable compareAsTree.
		if externs.HasMoreThanOneChild() {
			for c := externs.GetFirstChild(); c != nil; c = c.GetNext() {
				if !c.HasChildren() {
					c.Detach()
				}
			}
		}

		// Expected output parsed without implied block.
		jutil.CheckState(externs.IsRoot())

		expectedRoot.Detach()

		if b.compareJsDoc {
			NodeSubject_AssertNode(externs).
				UsingSerializer(compilerTestCase_createPrettyPrinter(compiler)).
				IsEqualIncludingJsDocTo(expectedRoot)
		} else {
			NodeSubject_AssertNode(externs).UsingSerializer(compilerTestCase_createPrettyPrinter(compiler)).IsEqualTo(expectedRoot)
		}
	} else {
		externsCode := compiler.ToSourceOfNode(externs)
		expectedCode := compiler.ToSourceOfNode(expectedRoot)
		AssertThat(externsCode).IsEqualTo(expectedCode)
	}

	if warnings != nil {
		var warningBuilder strings.Builder
		for _, actualWarning := range compiler.GetWarnings() {
			warningBuilder.WriteString(actualWarning.Description())
			warningBuilder.WriteString("\n")
		}
		warningMessage := warningBuilder.String()
		AssertWithMessage("There should be %s warnings. %s", len(warnings), warningMessage).
			That(compiler.GetWarningCount()).
			IsEqualTo(len(warnings))
		for i := 0; i < len(warnings); i++ {
			warning := warnings[i].diagnostic
			AssertWithMessage(warningMessage).
				That(compiler.GetWarnings()[i].Type()).
				IsEqualTo(warning)
		}
	}
}

// Java: CompilerTestCase.createCompiler()
func (b *CompilerTestCaseBase) CreateCompiler() *Compiler {
	compiler := NewCompiler()
	compiler.SetAllowableFeatures(b.acceptedLanguage.ToFeatureSet())
	if !b.webpackModulesById.IsEmpty() {
		compiler.InitWebpackMap(jutil.NewLinkedHashMapFrom(b.webpackModulesById))
	}
	compiler.SetPreferRegexParser(false)
	return compiler
}

// FindQualifiedNameNode finds the first matching qualified name node in
// post-traversal order.
//
// Java: CompilerTestCase.findQualifiedNameNode(String, Node)
func (b *CompilerTestCaseBase) FindQualifiedNameNode(name string, root *rhino.Node) *rhino.Node {
	return b.FindQualifiedNameNodes(name, root)[0]
}

// FindQualifiedNameNodes finds all the matching qualified name nodes in
// post-traversal order.
//
// Java: CompilerTestCase.findQualifiedNameNodes(String, Node)
func (b *CompilerTestCaseBase) FindQualifiedNameNodes(name string, root *rhino.Node) []*rhino.Node {
	matches := []*rhino.Node{}
	NodeUtil_VisitPostOrder(
		root,
		NodeUtil_VisitorFunc(func(n *rhino.Node) {
			if n.MatchesQualifiedName(name) {
				matches = append(matches, n)
			}
		}))
	return matches
}

// CompilerTestCase_FindDefinition returns a node that defines the given
// qualified name. The returned node should have a valid type defined on it,
// if type checking occurred. Will return null if no definition is found.
// Only nodes in the top-level script are traversed.
//
// Java: CompilerTestCase.findDefinition(Compiler, String)
func CompilerTestCase_FindDefinition(compiler *Compiler, name string) *rhino.Node {
	for node := compiler.GetJsRoot().GetFirstFirstChild(); node != nil; node = node.GetNext() {
		switch node.GetToken() {
		case rhino.Token_FUNCTION:
			if name == node.GetFirstChild().GetString() {
				return node
			}
		case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
			if name == node.GetFirstChild().GetString() {
				return node.GetFirstChild()
			}
		case rhino.Token_EXPR_RESULT:
			if node.GetFirstChild().IsAssign() &&
				node.GetFirstFirstChild().MatchesQualifiedName(name) {
				return node.GetFirstFirstChild()
			}
		default:
		}
	}
	return nil
}

// CompilerTestCase_Srcs is every overload of Java's srcs(...):
//
//	srcs(String srcText)
//	srcs(String... srcTexts)
//	srcs(List<SourceFile> files)
//	srcs(SourceFile... files)
//	srcs(JSChunk... chunks)
//
// A Java null String is a nil *string.
//
// Java: CompilerTestCase.srcs(String)
// Java: CompilerTestCase.srcs(String...)
// Java: CompilerTestCase.srcs(List<SourceFile>)
// Java: CompilerTestCase.srcs(SourceFile...)
// Java: CompilerTestCase.srcs(JSChunk...)
func CompilerTestCase_Srcs(args ...any) CompilerTestCase_Sources {
	switch {
	case compilerTestCase_argsAre(args, "string"):
		s := args[0].(string)
		return newCompilerTestCase_FlatSources(compilerTestCase_maybeCreateSources(CompilerTestCase_GENERATED_SRC_NAME, &s))
	case compilerTestCase_argsAre(args, "*string"):
		return newCompilerTestCase_FlatSources(compilerTestCase_maybeCreateSources(CompilerTestCase_GENERATED_SRC_NAME, args[0].(*string)))
	case compilerTestCase_argsAre(args, "[]string"):
		return newCompilerTestCase_FlatSources(CompilerTestCase_CreateSources(CompilerTestCase_GENERATED_SRC_NAME, args[0].([]string)))
	case compilerTestCase_argsAre(args, "[]*SourceFile"):
		return newCompilerTestCase_FlatSources(args[0].([]*SourceFile))
	case compilerTestCase_argsAre(args, "[]*JSChunk"):
		return newCompilerTestCase_ChunkSources(args[0].([]*JSChunk))
	}
	if strs, ok := compilerTestCase_allOf[string](args); ok {
		return newCompilerTestCase_FlatSources(CompilerTestCase_CreateSources(CompilerTestCase_GENERATED_SRC_NAME, strs))
	}
	if files, ok := compilerTestCase_allOf[*SourceFile](args); ok {
		return newCompilerTestCase_FlatSources(files)
	}
	if chunks, ok := compilerTestCase_allOf[*JSChunk](args); ok {
		return newCompilerTestCase_ChunkSources(chunks)
	}
	compilerTestCase_noOverload("srcs", args)
	return nil
}

// CompilerTestCase_Expected_ is every overload of Java's expected(...):
//
//	expected(String srcText)
//	expected(String... srcTexts)
//	expected(List<SourceFile> files)
//	expected(SourceFile... files)
//	expected(JSChunk[] chunks)
//
// Java: CompilerTestCase.expected(String)
// Java: CompilerTestCase.expected(String...)
// Java: CompilerTestCase.expected(List<SourceFile>)
// Java: CompilerTestCase.expected(SourceFile...)
// Java: CompilerTestCase.expected(JSChunk[])
func CompilerTestCase_Expected_(args ...any) *CompilerTestCase_Expected {
	switch {
	case compilerTestCase_argsAre(args, "string"):
		s := args[0].(string)
		return newCompilerTestCase_Expected(compilerTestCase_maybeCreateSources(CompilerTestCase_GENERATED_SRC_NAME, &s))
	case compilerTestCase_argsAre(args, "*string"):
		return newCompilerTestCase_Expected(compilerTestCase_maybeCreateSources(CompilerTestCase_GENERATED_SRC_NAME, args[0].(*string)))
	case compilerTestCase_argsAre(args, "[]string"):
		return newCompilerTestCase_Expected(CompilerTestCase_CreateSources(CompilerTestCase_GENERATED_SRC_NAME, args[0].([]string)))
	case compilerTestCase_argsAre(args, "[]*SourceFile"):
		return newCompilerTestCase_Expected(args[0].([]*SourceFile))
	case compilerTestCase_argsAre(args, "[]*JSChunk"):
		return compilerTestCase_expectedOfChunks(args[0].([]*JSChunk))
	}
	if strs, ok := compilerTestCase_allOf[string](args); ok {
		return newCompilerTestCase_Expected(CompilerTestCase_CreateSources(CompilerTestCase_GENERATED_SRC_NAME, strs))
	}
	if files, ok := compilerTestCase_allOf[*SourceFile](args); ok {
		return newCompilerTestCase_Expected(files)
	}
	compilerTestCase_noOverload("expected", args)
	return nil
}

// Java: CompilerTestCase.expected(JSChunk[])
func compilerTestCase_expectedOfChunks(chunks []*JSChunk) *CompilerTestCase_Expected {
	expectedSrcs := []string{}
	for _, chunk := range chunks {
		for _, input := range chunk.GetInputs() {
			expectedSrcs = append(expectedSrcs, input.GetSourceFile().GetCode())
		}
	}
	return CompilerTestCase_Expected_(expectedSrcs)
}

// CompilerTestCase_Externs_ is every overload of Java's externs(...):
//
//	externs(String externSrc)
//	externs(String... srcTexts)
//	externs(SourceFile... externs)
//	    Copy SourceFile objects to prevent the externs bit from polluting
//	    tests.
//	externs(List<SourceFile> files)
//
// Java: CompilerTestCase.externs(String)
// Java: CompilerTestCase.externs(String...)
// Java: CompilerTestCase.externs(SourceFile...)
// Java: CompilerTestCase.externs(List<SourceFile>)
func CompilerTestCase_Externs_(args ...any) *CompilerTestCase_Externs {
	switch {
	case compilerTestCase_argsAre(args, "string"):
		s := args[0].(string)
		return newCompilerTestCase_Externs(compilerTestCase_maybeCreateSources(CompilerTestCase_GENERATED_EXTERNS_NAME, &s))
	case compilerTestCase_argsAre(args, "*string"):
		return newCompilerTestCase_Externs(compilerTestCase_maybeCreateSources(CompilerTestCase_GENERATED_EXTERNS_NAME, args[0].(*string)))
	case compilerTestCase_argsAre(args, "[]string"):
		return newCompilerTestCase_Externs(CompilerTestCase_CreateSources(CompilerTestCase_GENERATED_EXTERNS_NAME, args[0].([]string)))
	case compilerTestCase_argsAre(args, "[]*SourceFile"):
		return newCompilerTestCase_Externs(args[0].([]*SourceFile))
	}
	if strs, ok := compilerTestCase_allOf[string](args); ok {
		return newCompilerTestCase_Externs(CompilerTestCase_CreateSources(CompilerTestCase_GENERATED_EXTERNS_NAME, strs))
	}
	if externs, ok := compilerTestCase_allOf[*SourceFile](args); ok {
		// Copy SourceFile objects to prevent the externs bit from polluting tests.
		copies := make([]*SourceFile, len(externs))
		for i, f := range externs {
			copies[i] = SourceFile_FromCode(f.GetName(), f.GetCode())
		}
		return newCompilerTestCase_Externs(copies)
	}
	compilerTestCase_noOverload("externs", args)
	return nil
}

// CompilerTestCase_Warning is warning(DiagnosticType).
//
// Java: CompilerTestCase.warning(DiagnosticType)
func CompilerTestCase_Warning(type_ *DiagnosticType) *CompilerTestCase_Diagnostic {
	return newCompilerTestCase_WarningDiagnostic(type_)
}

// CompilerTestCase_Error is error(DiagnosticType).
//
// Java: CompilerTestCase.error(DiagnosticType)
func CompilerTestCase_Error(type_ *DiagnosticType) *CompilerTestCase_Diagnostic {
	return newCompilerTestCase_ErrorDiagnostic(type_)
}

// CompilerTestCase_Postcondition_ is postcondition(Postcondition); a Java
// lambda is a CompilerTestCase_PostconditionFunc.
//
// Java: CompilerTestCase.postcondition(Postcondition)
func CompilerTestCase_Postcondition_(postcondition CompilerTestCase_Postcondition) CompilerTestCase_Postcondition {
	return postcondition
}

// Srcs is CompilerTestCase_Srcs, called unqualified from a subclass
// (Go-only: language-adapter).
func (b *CompilerTestCaseBase) Srcs(args ...any) CompilerTestCase_Sources {
	return CompilerTestCase_Srcs(args...)
}

// Expected is CompilerTestCase_Expected_, called unqualified from a subclass
// (Go-only: language-adapter).
func (b *CompilerTestCaseBase) Expected(args ...any) *CompilerTestCase_Expected {
	return CompilerTestCase_Expected_(args...)
}

// Externs is CompilerTestCase_Externs_, called unqualified from a subclass
// (Go-only: language-adapter).
func (b *CompilerTestCaseBase) Externs(args ...any) *CompilerTestCase_Externs {
	return CompilerTestCase_Externs_(args...)
}

// Warning is CompilerTestCase_Warning, called unqualified from a subclass
// (Go-only: language-adapter).
func (b *CompilerTestCaseBase) Warning(type_ *DiagnosticType) *CompilerTestCase_Diagnostic {
	return CompilerTestCase_Warning(type_)
}

// Error is CompilerTestCase_Error, called unqualified from a subclass
// (Go-only: language-adapter).
func (b *CompilerTestCaseBase) Error(type_ *DiagnosticType) *CompilerTestCase_Diagnostic {
	return CompilerTestCase_Error(type_)
}

// Postcondition is CompilerTestCase_Postcondition_ for a Java lambda, called
// unqualified from a subclass (Go-only: language-adapter).
func (b *CompilerTestCaseBase) Postcondition(verify func(compiler *Compiler)) CompilerTestCase_Postcondition {
	return CompilerTestCase_PostconditionFunc(verify)
}

// testInternalParts is testInternal(Iterable<TestPart>).
//
// Java: CompilerTestCase.testInternal(Iterable<TestPart>)
func (b *CompilerTestCaseBase) testInternalParts(parts []CompilerTestCase_TestPart) {
	// TODO(johnlenz): make "ignore" and "nothing" explicit.
	var externs *CompilerTestCase_Externs
	var srcs CompilerTestCase_Sources
	var expected *CompilerTestCase_Expected
	diagnostics := []*CompilerTestCase_Diagnostic{}
	postconditions := []CompilerTestCase_Postcondition{}
	for _, part := range parts {
		switch part := part.(type) {
		case *CompilerTestCase_Externs:
			jutil.CheckState(externs == nil)
			externs = part
		case CompilerTestCase_Sources:
			jutil.CheckState(srcs == nil)
			srcs = part
		case *CompilerTestCase_Expected:
			jutil.CheckState(expected == nil)
			expected = part
		case *CompilerTestCase_Diagnostic:
			diagnostics = append(diagnostics, part)
		case CompilerTestCase_Postcondition:
			postconditions = append(postconditions, part)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(fmt.Sprintf("unexpected %T", part)))
		}
	}
	if compilerTestCase_EXPECTED_SAME == expected {
		expected = compilerTestCase_fromSources(srcs)
	}
	if externs == nil {
		externs = CompilerTestCase_Externs_(b.defaultExternsInputs)
	}
	b.self.TestInternal(externs, srcs, expected, diagnostics, postconditions)
}

// Java: CompilerTestCase.fromSources(Sources)
func compilerTestCase_fromSources(srcs CompilerTestCase_Sources) *CompilerTestCase_Expected {
	switch srcs := srcs.(type) {
	case *CompilerTestCase_FlatSources:
		return CompilerTestCase_Expected_(srcs.sources)
	case *CompilerTestCase_ChunkSources:
		return CompilerTestCase_Expected_(srcs.chunks)
	}
	panic(jutil.NewIllegalStateExceptionMsg("unexpected"))
}

// CompilerTestCase_TestPart is the marker interface for configuration
// parameters of a test invocation. This is essentially a closed union type
// consisting of four subtypes: (1) Externs, (2) Expected, (3) Sources, and
// (4) Diagnostic. Sharing a common marker interface, allows specifying any
// combination of these types to the 'test' method, and makes it very clear
// what function each parameter serves (since the first three can all be
// represented as simple strings). Moreover, it reduces the combinatorial
// explosion of different ways to express the parts (a single string, a list
// of SourceFiles, a list of chunks, a DiagnosticType, a DiagnosticType with
// an expected message, etc).
//
// Note that this API is intended to be a medium-term temporary API while
// developing and migrating to a more fluent assertion style.
//
// TODO(johnlenz): make this a abstract class with a private constructor
//
// Java: CompilerTestCase.TestPart
type CompilerTestCase_TestPart interface {
	isCompilerTestCase_TestPart()
}

// CompilerTestCase_Expected is the expected output of a test.
//
// Java: CompilerTestCase.Expected
type CompilerTestCase_Expected struct {
	expected []*SourceFile
}

func (*CompilerTestCase_Expected) isCompilerTestCase_TestPart() {}

// Java: CompilerTestCase.Expected(List<SourceFile>)
func newCompilerTestCase_Expected(files []*SourceFile) *CompilerTestCase_Expected {
	return &CompilerTestCase_Expected{expected: files}
}

// compilerTestCase_EXPECTED_SAME marks testSame: the expected output is the
// input.
//
// TODO(sdh): make a shorter function to get this - e.g. same(), expectSame(), sameOutput() ?
// TODO(sdh): also make an ignoreOutput() and noOutput() ?
var compilerTestCase_EXPECTED_SAME = newCompilerTestCase_Expected(nil)

// CompilerTestCase_Sources is the input of a test.
//
// Java: CompilerTestCase.Sources
type CompilerTestCase_Sources interface {
	CompilerTestCase_TestPart
	isCompilerTestCase_Sources()
}

// CompilerTestCase_FlatSources are sources in a single chunk.
//
// Java: CompilerTestCase.FlatSources
type CompilerTestCase_FlatSources struct {
	sources []*SourceFile
}

func (*CompilerTestCase_FlatSources) isCompilerTestCase_TestPart() {}
func (*CompilerTestCase_FlatSources) isCompilerTestCase_Sources()  {}

// Java: CompilerTestCase.FlatSources(List<SourceFile>)
func newCompilerTestCase_FlatSources(files []*SourceFile) *CompilerTestCase_FlatSources {
	// ImmutableList.copyOf(files) throws on null.
	jutil.CheckNotNull(files)
	return &CompilerTestCase_FlatSources{sources: append([]*SourceFile(nil), files...)}
}

// CompilerTestCase_ChunkSources are sources in chunks.
//
// Java: CompilerTestCase.ChunkSources
type CompilerTestCase_ChunkSources struct {
	chunks []*JSChunk
}

func (*CompilerTestCase_ChunkSources) isCompilerTestCase_TestPart() {}
func (*CompilerTestCase_ChunkSources) isCompilerTestCase_Sources()  {}

// Java: CompilerTestCase.ChunkSources(JSChunk[])
func newCompilerTestCase_ChunkSources(chunks []*JSChunk) *CompilerTestCase_ChunkSources {
	return &CompilerTestCase_ChunkSources{chunks: append([]*JSChunk(nil), chunks...)}
}

// CompilerTestCase_Externs are the externs of a test.
//
// Java: CompilerTestCase.Externs
type CompilerTestCase_Externs struct {
	externs []*SourceFile
}

func (*CompilerTestCase_Externs) isCompilerTestCase_TestPart() {}

// Java: CompilerTestCase.Externs(List<SourceFile>)
func newCompilerTestCase_Externs(files []*SourceFile) *CompilerTestCase_Externs {
	e := &CompilerTestCase_Externs{externs: files}
	if files != nil {
		for _, s := range files {
			s.SetKind(rhino.StaticSourceFile_SourceKind_EXTERN)
		}
	}
	return e
}

// CompilerTestCase_Diagnostic is an expected error or warning.
//
// Java's private subclasses ErrorDiagnostic and WarningDiagnostic only fix
// the level in their constructors; they are folded into this struct
// (Go-only: language-adapter).
//
// Java: CompilerTestCase.Diagnostic
type CompilerTestCase_Diagnostic struct {
	level            CheckLevel
	diagnostic       *DiagnosticType
	messagePredicate *compilerTestCase_NamedPredicate
	line             int
	charno           int
	length           int
}

func (*CompilerTestCase_Diagnostic) isCompilerTestCase_TestPart() {}

// Java: CompilerTestCase.Diagnostic(CheckLevel, DiagnosticType, NamedPredicate<String>)
func newCompilerTestCase_Diagnostic(
	level CheckLevel,
	diagnostic *DiagnosticType,
	messagePredicate *compilerTestCase_NamedPredicate) *CompilerTestCase_Diagnostic {
	return newCompilerTestCase_DiagnosticWithLocation(level, diagnostic, messagePredicate, -1, -1, -1)
}

// Java: CompilerTestCase.Diagnostic(CheckLevel, DiagnosticType, NamedPredicate<String>, int, int, int)
func newCompilerTestCase_DiagnosticWithLocation(
	level CheckLevel,
	diagnostic *DiagnosticType,
	messagePredicate *compilerTestCase_NamedPredicate,
	line int,
	charno int,
	length int) *CompilerTestCase_Diagnostic {
	return &CompilerTestCase_Diagnostic{
		level:            level,
		diagnostic:       diagnostic,
		messagePredicate: messagePredicate,
		line:             line,
		charno:           charno,
		length:           length,
	}
}

// formatDiff returns a description of the mismatch, or nil if error
// matches.
//
// Java: CompilerTestCase.Diagnostic.formatDiff(JSError)
func (d *CompilerTestCase_Diagnostic) formatDiff(error_ *JSError) *string {
	if !jutil.ObjectsEquals(d.diagnostic, error_.Type()) {
		return compilerTestCase_strPtr(
			fmt.Sprintf(
				"diagnostic type <%s> did not match <%s>", error_.Type().key, d.diagnostic.key))
	}
	if d.messagePredicate != nil && !d.messagePredicate.Apply(error_.Description()) {
		return compilerTestCase_strPtr(
			fmt.Sprintf("message <%s> was not <%s>", error_.Description(), d.messagePredicate))
	}
	if d.line != -1 && error_.Lineno() != d.line {
		return compilerTestCase_strPtr(fmt.Sprintf("line <%d> did not match <%d>", error_.Lineno(), d.line))
	}
	if d.charno != -1 && error_.Charno() != d.charno {
		return compilerTestCase_strPtr(fmt.Sprintf("charno <%d> did not match <%d>", error_.Charno(), d.charno))
	}
	if d.length != -1 && error_.Length() != d.length {
		return compilerTestCase_strPtr(fmt.Sprintf("length <%d> did not match <%d>", error_.Length(), d.length))
	}
	return nil
}

func compilerTestCase_strPtr(s string) *string { return &s }

// Java: CompilerTestCase.Diagnostic.withMessage(String)
func (d *CompilerTestCase_Diagnostic) WithMessage(expectedRaw string) *CompilerTestCase_Diagnostic {
	expected := jutil.StringTrim(expectedRaw)
	jutil.CheckState(d.messagePredicate == nil)
	return newCompilerTestCase_Diagnostic(
		d.level,
		d.diagnostic,
		compilerTestCase_NamedPredicate_Of(func(actual string) bool { return jutil.StringTrim(actual) == expected }, "\""+expected+"\""))
}

// Java: CompilerTestCase.Diagnostic.withMessageContaining(String)
func (d *CompilerTestCase_Diagnostic) WithMessageContaining(substring string) *CompilerTestCase_Diagnostic {
	jutil.CheckState(d.messagePredicate == nil)
	return newCompilerTestCase_Diagnostic(
		d.level,
		d.diagnostic,
		compilerTestCase_NamedPredicate_Of(
			func(message string) bool { return strings.Contains(message, substring) }, "containing \""+substring+"\""))
}

// Java: CompilerTestCase.Diagnostic.withLocation(int, int, int)
func (d *CompilerTestCase_Diagnostic) WithLocation(line int, charno int, length int) *CompilerTestCase_Diagnostic {
	return newCompilerTestCase_DiagnosticWithLocation(d.level, d.diagnostic, d.messagePredicate, line, charno, length)
}

// Java: CompilerTestCase.Diagnostic.toString()
func (d *CompilerTestCase_Diagnostic) String() string {
	s := d.diagnostic.key
	if d.messagePredicate != nil {
		s += " with message " + d.messagePredicate.String()
	}
	return s
}

// Java: CompilerTestCase.ErrorDiagnostic(DiagnosticType)
func newCompilerTestCase_ErrorDiagnostic(diagnostic *DiagnosticType) *CompilerTestCase_Diagnostic {
	return newCompilerTestCase_Diagnostic(CheckLevel_ERROR, diagnostic, nil)
}

// Java: CompilerTestCase.WarningDiagnostic(DiagnosticType)
func newCompilerTestCase_WarningDiagnostic(diagnostic *DiagnosticType) *CompilerTestCase_Diagnostic {
	return newCompilerTestCase_Diagnostic(CheckLevel_WARNING, diagnostic, nil)
}

// Java: CompilerTestCase.DIAGNOSTIC_CORRESPONDENCE
var compilerTestCase_DIAGNOSTIC_CORRESPONDENCE = FormattingDiffsUsing(
	Correspondence_From(
		func(actual *JSError, expected *CompilerTestCase_Diagnostic) bool {
			return expected.formatDiff(actual) == nil
		},
		"is a JSError matching"),
	func(actual *JSError, expected *CompilerTestCase_Diagnostic) string {
		return *expected.formatDiff(actual)
	})

// compilerTestCase_NamedPredicate is a Predicate<String> with a name.
//
// Java: CompilerTestCase.NamedPredicate
type compilerTestCase_NamedPredicate struct {
	delegate func(string) bool
	name     string
}

// Java: CompilerTestCase.NamedPredicate.of(Predicate<T>, String)
func compilerTestCase_NamedPredicate_Of(delegate func(string) bool, name string) *compilerTestCase_NamedPredicate {
	return newCompilerTestCase_NamedPredicate(delegate, name)
}

// Java: CompilerTestCase.NamedPredicate(Predicate<T>, String)
func newCompilerTestCase_NamedPredicate(delegate func(string) bool, name string) *compilerTestCase_NamedPredicate {
	return &compilerTestCase_NamedPredicate{delegate: delegate, name: name}
}

// Java: CompilerTestCase.NamedPredicate.apply(T)
func (p *compilerTestCase_NamedPredicate) Apply(arg string) bool {
	return p.delegate(arg)
}

// Java: CompilerTestCase.NamedPredicate.toString()
func (p *compilerTestCase_NamedPredicate) String() string {
	return p.name
}

// CompilerTestCase_Postcondition is checked on the compiler after the test.
//
// Java: CompilerTestCase.Postcondition
type CompilerTestCase_Postcondition interface {
	CompilerTestCase_TestPart
	// Java: CompilerTestCase.Postcondition.verify(Compiler)
	Verify(compiler *Compiler)
}

// CompilerTestCase_PostconditionFunc is a Java lambda implementing
// Postcondition (Go-only: language-adapter).
type CompilerTestCase_PostconditionFunc func(compiler *Compiler)

func (CompilerTestCase_PostconditionFunc) isCompilerTestCase_TestPart() {}

// Verify calls f.
func (f CompilerTestCase_PostconditionFunc) Verify(compiler *Compiler) { f(compiler) }

// Java: CompilerTestCase.createPrettyPrinter(Compiler)
func compilerTestCase_createPrettyPrinter(compiler *Compiler) func(*rhino.Node) string {
	return func(n *rhino.Node) string {
		return NewCodePrinter_Builder(n).
			SetCompilerOptions(compiler.GetOptions()).
			SetPrettyPrint(true).
			Build()
	}
}

// Java: CompilerTestCase.makePassFactory(String, Function<AbstractCompiler, CompilerPass>)
func CompilerTestCase_MakePassFactory(name string, pass func(AbstractCompiler) CompilerPass) *PassFactory {
	return PassFactory_Builder_().SetName(name).SetInternalFactory(pass).Build()
}

// --- Go-only: language-adapter — run-time overload resolution.

// compilerTestCase_argsAre reports whether args match the Java parameter
// kinds: "string", "*string" (a nullable String), "[]string",
// "DiagnosticType", "Diagnostic", "Sources", "Externs", "Expected",
// "[]*SourceFile", "[]*JSChunk".
func compilerTestCase_argsAre(args []any, kinds ...string) bool {
	if len(args) != len(kinds) {
		return false
	}
	for i, k := range kinds {
		var ok bool
		switch k {
		case "string":
			_, ok = args[i].(string)
		case "*string":
			_, ok = args[i].(*string)
		case "[]string":
			_, ok = args[i].([]string)
		case "DiagnosticType":
			_, ok = args[i].(*DiagnosticType)
		case "Diagnostic":
			_, ok = args[i].(*CompilerTestCase_Diagnostic)
		case "Sources":
			_, ok = args[i].(CompilerTestCase_Sources)
		case "Externs":
			_, ok = args[i].(*CompilerTestCase_Externs)
		case "Expected":
			_, ok = args[i].(*CompilerTestCase_Expected)
		case "[]*SourceFile":
			_, ok = args[i].([]*SourceFile)
		case "[]*JSChunk":
			_, ok = args[i].([]*JSChunk)
		default:
			panic("unknown kind " + k)
		}
		if !ok {
			return false
		}
	}
	return true
}

// compilerTestCase_allOf matches a Java varargs parameter T...: every
// argument is a T (and there is at least one, or none of any kind).
func compilerTestCase_allOf[T any](args []any) ([]T, bool) {
	out := make([]T, len(args))
	for i, a := range args {
		t, ok := a.(T)
		if !ok {
			return nil, false
		}
		out[i] = t
	}
	return out, true
}

// compilerTestCase_testParts matches a TestPart... parameter.
func compilerTestCase_testParts(method string, args []any) []CompilerTestCase_TestPart {
	parts, ok := compilerTestCase_allOf[CompilerTestCase_TestPart](args)
	if !ok {
		compilerTestCase_noOverload(method, args)
	}
	return parts
}

// compilerTestCase_diagnostics matches a Diagnostic... parameter.
func compilerTestCase_diagnostics(method string, args []any) []*CompilerTestCase_Diagnostic {
	ds, ok := compilerTestCase_allOf[*CompilerTestCase_Diagnostic](args)
	if !ok {
		compilerTestCase_noOverload(method, args)
	}
	return ds
}

// compilerTestCase_noOverload is a call that javac would have rejected.
func compilerTestCase_noOverload(method string, args []any) {
	types := make([]string, len(args))
	for i, a := range args {
		types[i] = fmt.Sprintf("%T", a)
	}
	panic(fmt.Sprintf("CompilerTestCase.%s has no overload for (%s)", method, strings.Join(types, ", ")))
}

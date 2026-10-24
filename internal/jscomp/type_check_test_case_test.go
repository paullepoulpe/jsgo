// Port of com/google/javascript/jscomp/TypeCheckTestCase.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/TypeCheckTestCase.java)

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypeCheckTestCase is the base class of the TypeCheck tests.
type TypeCheckTestCase interface {
	CompilerTypeTestCase

	// GetTypeCheckTestCaseBase gives access to the base class.
	GetTypeCheckTestCaseBase() *TypeCheckTestCaseBase
}

// TypeCheckTestCaseBase holds TypeCheckTestCase's base implementations.
type TypeCheckTestCaseBase struct {
	CompilerTypeTestCaseBase
}

// InitTypeCheckTestCaseBase is the Go form of the implicit super().
//
// Java: TypeCheckTestCase()
func (b *TypeCheckTestCaseBase) InitTypeCheckTestCaseBase(self TypeCheckTestCase) {
	b.InitCompilerTypeTestCaseBase(self)
}

// GetTypeCheckTestCaseBase gives access to the base class.
func (b *TypeCheckTestCaseBase) GetTypeCheckTestCaseBase() *TypeCheckTestCaseBase {
	return b
}

// Java: TypeCheckTestCase.setUp()
func (b *TypeCheckTestCaseBase) SetUp() {
	b.CompilerTypeTestCaseBase.SetUp()
	// Enable missing override checks that are disabled by default.
	b.compiler.GetOptions().SetWarningLevel(DiagnosticGroups_MISSING_OVERRIDE, CheckLevel_WARNING)
	b.compiler.
		GetOptions().
		SetWarningLevel(DiagnosticGroups_STRICT_MISSING_PROPERTIES, CheckLevel_WARNING)
	b.compiler.
		GetOptions().
		SetWarningLevel(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS, CheckLevel_WARNING)
}

// Java: TypeCheckTestCase.getInstanceType(Node)
func TypeCheckTestCase_GetInstanceType(js1Node *rhino.Node) rhino.ObjectType {
	type_ := js1Node.GetFirstChild().GetJSType()
	AssertThat(type_).IsNotNull()
	IsInstanceOf[rhino.FunctionType](AssertThat(type_))
	functionType := type_.(rhino.FunctionType)
	AssertThat(functionType.IsConstructor()).IsTrue()
	return functionType.GetInstanceType()
}

// Java: TypeCheckTestCase.checkObjectType(ObjectType, String, JSType)
func (b *TypeCheckTestCaseBase) CheckObjectType(objectType rhino.ObjectType, propertyName string, expectedType rhino.JSType) {
	AssertWithMessage(
		"Expected %s to have property %s", jutil.ToJavaString(objectType.GetReferenceName()), propertyName).
		That(objectType.HasProperty(propertyName)).
		IsTrue()
	b.AssertTypeEqualsMsg(
		"Expected "+
			jutil.ToJavaString(objectType.GetReferenceName())+
			"'s property "+
			propertyName+
			" to have type "+
			jutil.ToJavaString(expectedType),
		expectedType,
		objectType.GetPropertyType(propertyName))
}

// NewTestLegacy creates a new TypeTestBuilder based on the instance state in
// this class.
//
// Deprecated: prefer to statically import TypeTestBuilder.newTest, to make
// tests more hermetic.
//
// Java: TypeCheckTestCase.newTestLegacy()
func (b *TypeCheckTestCaseBase) NewTestLegacy() *TypeCheckTestCase_TypeTestBuilder {
	return NewTypeCheckTestCase_TypeTestBuilder(b.compiler)
}

// TypeCheckTestCase_TypeTestBuilder builds and runs one type-checking test.
type TypeCheckTestCase_TypeTestBuilder struct {
	sources                []*SourceFile
	externs                []string
	includeDefaultExterns  bool
	diagnosticTypes        []*DiagnosticType
	diagnosticDescriptions []string
	diagnosticsAreErrors   bool
	reportUnknownTypes     bool
	suppress               []*DiagnosticGroup
	hasRun                 bool
	compiler               *Compiler
}

// Java: TypeCheckTestCase.TypeTestBuilder.newTest()
func TypeCheckTestCase_TypeTestBuilder_NewTest() *TypeCheckTestCase_TypeTestBuilder {
	compiler := NewCompiler()
	options := CompilerTypeTestCase_DefaultOptions()
	// Enable missing override checks that are disabled by default.
	options.SetWarningLevel(DiagnosticGroups_MISSING_OVERRIDE, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_STRICT_MISSING_PROPERTIES, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS, CheckLevel_WARNING)
	compiler.InitOptions(options)

	compiler.MarkFeatureNotAllowed(parser.FeatureSet_Feature_MODULES)
	return NewTypeCheckTestCase_TypeTestBuilder(compiler)
}

// Java: TypeCheckTestCase.TypeTestBuilder(Compiler)
func NewTypeCheckTestCase_TypeTestBuilder(compiler *Compiler) *TypeCheckTestCase_TypeTestBuilder {
	return &TypeCheckTestCase_TypeTestBuilder{compiler: compiler}
}

// Java: TypeCheckTestCase.TypeTestBuilder.addSource(String)
func (b *TypeCheckTestCase_TypeTestBuilder) AddSource(code string) *TypeCheckTestCase_TypeTestBuilder {
	return b.AddSourceWithFilename("testcode"+strconv.Itoa(len(b.sources)), code)
}

// Java: TypeCheckTestCase.TypeTestBuilder.addSource(String, String)
func (b *TypeCheckTestCase_TypeTestBuilder) AddSourceWithFilename(filename string, code string) *TypeCheckTestCase_TypeTestBuilder {
	b.sources = append(b.sources, SourceFile_FromCode(filename, code))
	return b
}

// Java: TypeCheckTestCase.TypeTestBuilder.addExterns(String)
func (b *TypeCheckTestCase_TypeTestBuilder) AddExterns(externs string) *TypeCheckTestCase_TypeTestBuilder {
	b.externs = append(b.externs, externs)
	return b
}

// Java: TypeCheckTestCase.TypeTestBuilder.includeDefaultExterns()
func (b *TypeCheckTestCase_TypeTestBuilder) IncludeDefaultExterns() *TypeCheckTestCase_TypeTestBuilder {
	b.includeDefaultExterns = true
	return b
}

// Java: TypeCheckTestCase.TypeTestBuilder.diagnosticsAreErrors()
func (b *TypeCheckTestCase_TypeTestBuilder) DiagnosticsAreErrors() *TypeCheckTestCase_TypeTestBuilder {
	b.diagnosticsAreErrors = true
	return b
}

// Java: TypeCheckTestCase.TypeTestBuilder.enableReportUnknownTypes()
func (b *TypeCheckTestCase_TypeTestBuilder) EnableReportUnknownTypes() *TypeCheckTestCase_TypeTestBuilder {
	b.reportUnknownTypes = true
	return b
}

// Java: TypeCheckTestCase.TypeTestBuilder.suppress(DiagnosticGroup)
func (b *TypeCheckTestCase_TypeTestBuilder) Suppress(group *DiagnosticGroup) *TypeCheckTestCase_TypeTestBuilder {
	b.suppress = append(b.suppress, group)
	return b
}

// Java: TypeCheckTestCase.TypeTestBuilder.addDiagnostic(DiagnosticType)
func (b *TypeCheckTestCase_TypeTestBuilder) AddDiagnosticType(x *DiagnosticType) *TypeCheckTestCase_TypeTestBuilder {
	b.diagnosticTypes = append(b.diagnosticTypes, x)
	return b
}

// Java: TypeCheckTestCase.TypeTestBuilder.addDiagnostic(String)
func (b *TypeCheckTestCase_TypeTestBuilder) AddDiagnostic(x string) *TypeCheckTestCase_TypeTestBuilder {
	b.diagnosticDescriptions = append(b.diagnosticDescriptions, jutil.StringTrim(x))
	return b
}

// Java: TypeCheckTestCase.TypeTestBuilder.run()
func (b *TypeCheckTestCase_TypeTestBuilder) Run() {
	jutil.CheckStateMsg(len(b.sources) != 0, "Must provide a source")
	jutil.CheckStateMsg(!b.hasRun, "Cannot run the same test twice")
	b.hasRun = true
	if b.includeDefaultExterns {
		b.externs = append([]string{CompilerTypeTestCase_DEFAULT_EXTERNS}, b.externs...)
	}

	allExterns := strings.Join(b.externs, "\n")

	var diagnostics []any
	var correspondence *Correspondence
	if len(b.diagnosticTypes) != 0 {
		for _, x := range b.diagnosticTypes {
			diagnostics = append(diagnostics, x)
		}
		correspondence = JSCompCorrespondences_DIAGNOSTIC_EQUALITY
	} else {
		for _, x := range b.diagnosticDescriptions {
			diagnostics = append(diagnostics, x)
		}
		correspondence = JSCompCorrespondences_DESCRIPTION_EQUALITY
	}
	for _, suppress := range b.suppress {
		b.compiler.GetOptions().SetWarningLevel(suppress, CheckLevel_OFF)
	}

	TypeCheckTestCase_ParseAndTypeCheckWithScopeWithCompiler(b.compiler, allExterns, b.sources, b.reportUnknownTypes)

	var assertedErrors []*JSError
	var emptyErrors []*JSError
	if b.diagnosticsAreErrors {
		assertedErrors = b.compiler.GetErrors()
		emptyErrors = b.compiler.GetWarnings()
	} else {
		assertedErrors = b.compiler.GetWarnings()
		emptyErrors = b.compiler.GetErrors()
	}

	AssertThat(assertedErrors).
		ComparingElementsUsing(correspondence).
		ContainsExactlyElementsIn(diagnostics).
		InOrder()
	AssertThat(emptyErrors).IsEmpty()
}

// ParseAndTypeCheck parses and type checks the JavaScript code.
//
// Java: TypeCheckTestCase.parseAndTypeCheck(String)
func (b *TypeCheckTestCaseBase) ParseAndTypeCheck(js string) *rhino.Node {
	return b.ParseAndTypeCheckWithExterns("", js)
}

// Java: TypeCheckTestCase.parseAndTypeCheck(String, String)
func (b *TypeCheckTestCaseBase) ParseAndTypeCheckWithExterns(externs string, js string) *rhino.Node {
	return b.ParseAndTypeCheckWithScopeWithSourceNameExtension(externs, js, "").root
}

// ParseAndTypeCheckWithScope parses and type checks the JavaScript code and
// returns the TypedScope used whilst type checking.
//
// Java: TypeCheckTestCase.parseAndTypeCheckWithScope(String)
func (b *TypeCheckTestCaseBase) ParseAndTypeCheckWithScope(js string) *TypeCheckTestCase_TypeCheckResult {
	return b.ParseAndTypeCheckWithScopeWithSourceNameExtension("", js, "")
}

// ParseAndTypeCheckWithScopeWithExterns parses and type checks the
// JavaScript code and returns the TypedScope used whilst type checking.
//
// Java: TypeCheckTestCase.parseAndTypeCheckWithScope(String, String)
func (b *TypeCheckTestCaseBase) ParseAndTypeCheckWithScopeWithExterns(externs string, js string) *TypeCheckTestCase_TypeCheckResult {
	return b.ParseAndTypeCheckWithScopeWithSourceNameExtension(externs, js, "")
}

// Java: TypeCheckTestCase.parseAndTypeCheckWithScope(String, String, String)
func (b *TypeCheckTestCaseBase) ParseAndTypeCheckWithScopeWithSourceNameExtension(
	externs string, js string, sourceNameExtension string) *TypeCheckTestCase_TypeCheckResult {
	return TypeCheckTestCase_ParseAndTypeCheckWithScopeWithCompiler(
		b.compiler,
		externs,
		[]*SourceFile{SourceFile_FromCode("[testcode]"+sourceNameExtension, js)},
		false)
}

// Java: TypeCheckTestCase.parseAndTypeCheckWithScope(Compiler, String, List<SourceFile>, boolean)
func TypeCheckTestCase_ParseAndTypeCheckWithScopeWithCompiler(
	compiler *Compiler, externs string, sources []*SourceFile, reportUnknownTypes bool) *TypeCheckTestCase_TypeCheckResult {
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("[externs]", externs)},
		sources,
		compiler.GetOptions())
	compiler.Parse()

	jsNode := compiler.GetJsRoot()
	externsNode := compiler.GetExternsRoot()
	NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).Process(externsNode, jsNode)
	NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).Process(externsNode, jsNode)
	NewInferConsts(compiler).Process(externsNode, jsNode)

	AssertWithMessage("Regarding errors:").That(compiler.GetErrors()).IsEmpty()

	s :=
		TypeCheckTestCase_MakeTypeCheckWithCompiler(compiler, compiler.GetTypeRegistry()).
			ReportUnknownTypes(reportUnknownTypes).
			ProcessForTesting(externsNode, jsNode)

	NewAstValidator(compiler).
		SetTypeValidationMode(AstValidator_TypeInfoValidation_JSTYPE).
		Process(externsNode, jsNode)

	return NewTypeCheckTestCase_TypeCheckResult(jsNode.GetFirstChild(), s)
}

// TypeCheck: n is a valid statement node, SCRIPT, or ROOT node.
//
// Java: TypeCheckTestCase.typeCheck(Node)
func (b *TypeCheckTestCaseBase) TypeCheck(n *rhino.Node) *rhino.Node {
	var jsRoot *rhino.Node
	if n.IsRoot() {
		// This is fine as is.
		jsRoot = n
	} else if n.IsScript() {
		jsRoot = rhino.IR_Root(n)
	} else {
		script := rhino.IR_ScriptOfNodes(n)
		jsRoot = rhino.IR_Root(script)
		script.SetInputId(rhino.NewInputId("test"))
	}
	externsNode := rhino.IR_Root()
	externAndJsRoot := rhino.IR_Root(externsNode)
	externAndJsRoot.AddChildToBack(jsRoot)

	b.MakeTypeCheck().ProcessForTesting(externsNode, externAndJsRoot.GetSecondChild())
	return n
}

// Java: TypeCheckTestCase.makeTypeCheck()
func (b *TypeCheckTestCaseBase) MakeTypeCheck() *TypeCheck {
	return TypeCheckTestCase_MakeTypeCheckWithCompiler(b.compiler, b.registry)
}

// Java: TypeCheckTestCase.makeTypeCheck(Compiler, JSTypeRegistry)
func TypeCheckTestCase_MakeTypeCheckWithCompiler(compiler *Compiler, registry *rhino.JSTypeRegistry) *TypeCheck {
	return NewTypeCheck(compiler, NewSemanticReverseAbstractInterpreter(registry), registry)
}

// Java: TypeCheckTestCase.suppressMissingProperty(String...)
func (b *TypeCheckTestCaseBase) SuppressMissingProperty(props ...string) string {
	result := "function dummy(x) { "
	for _, prop := range props {
		result += "x." + prop + " = 3;"
	}
	return result + "}"
}

// Java: TypeCheckTestCase.suppressMissingPropertyFor(String, String...)
func (b *TypeCheckTestCaseBase) SuppressMissingPropertyFor(type_ string, props ...string) string {
	result := "function dummy(x) { "
	for _, prop := range props {
		result += type_ + ".prototype." + prop + " = 3;"
	}
	return result + "}"
}

// Java: TypeCheckTestCase.assertHasXMorePropertiesThanNativeObject(ObjectType, int)
func (b *TypeCheckTestCaseBase) AssertHasXMorePropertiesThanNativeObject(
	instanceType rhino.ObjectType, numExtraProperties int) {
	AssertThat(instanceType.GetPropertiesCount()).
		IsEqualTo(b.GetNativeObjectPropertiesCount() + numExtraProperties)
}

// Java: TypeCheckTestCase.getNativeObjectPropertiesCount()
func (b *TypeCheckTestCaseBase) GetNativeObjectPropertiesCount() int {
	return b.GetNativeObjectType().GetPropertiesCount()
}

// TypeCheckTestCase_TypeCheckResult is the root and top scope of a
// type-checked program.
type TypeCheckTestCase_TypeCheckResult struct {
	root  *rhino.Node
	scope *TypedScope
}

// Java: TypeCheckTestCase.TypeCheckResult(Node, TypedScope)
func NewTypeCheckTestCase_TypeCheckResult(root *rhino.Node, scope *TypedScope) *TypeCheckTestCase_TypeCheckResult {
	return &TypeCheckTestCase_TypeCheckResult{root: root, scope: scope}
}

// typeCheckTests_noFixture gives a ported test class without a superclass,
// @Before or @After methods its JUnitTestCase hooks.
//
// Go-only: test-helper — the JUnit runner's view of a plain test class.
type typeCheckTests_noFixture struct {
	t *testing.T
}

// JUnitSetT is part of JUnitTestCase.
func (b *typeCheckTests_noFixture) JUnitSetT(t *testing.T) { b.t = t }

// T is the Go test running the current test method.
func (b *typeCheckTests_noFixture) T() *testing.T { return b.t }

// JUnitBefore runs the class's @Before methods: none.
func (b *typeCheckTests_noFixture) JUnitBefore() {}

// JUnitAfter runs the class's @After methods: none.
func (b *typeCheckTests_noFixture) JUnitAfter() {}

// typeCheckTests_join is Guava's Joiner.on(separator).join(parts): the
// parts' toString() joined by separator.
//
// Go-only: test-helper — Guava Joiner.
func typeCheckTests_join[T any](parts []T, separator string) string {
	strs := make([]string, len(parts))
	for i, p := range parts {
		strs[i] = jutil.ToJavaString(p)
	}
	return strings.Join(strs, separator)
}

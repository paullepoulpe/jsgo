// Port of com/google/javascript/jscomp/CompilerTypeTestCase.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CompilerTypeTestCase.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// CompilerTypeTestCase is mostly used by passes testing TypeCheck.
//
// The interface declares the overridable methods; subclasses embed
// CompilerTypeTestCaseBase (PORTING.md §4). A subclass is run with
// RunJUnit4 like a CompilerTestCase (see compiler_test_case_test.go).
type CompilerTypeTestCase interface {
	JUnitTestCase

	// GetCompilerTypeTestCaseBase gives access to the base class.
	GetCompilerTypeTestCaseBase() *CompilerTypeTestCaseBase

	// Java: CompilerTypeTestCase.getDefaultOptions()
	GetDefaultOptions() *CompilerOptions
	// Java: CompilerTypeTestCase.getCodingConvention()
	GetCodingConvention() CodingConvention
	// Java: CompilerTypeTestCase.setUp()
	SetUp()
	// Java: CompilerTypeTestCase.validateWarningsAndErrors()
	ValidateWarningsAndErrors()
	// Java: CompilerTypeTestCase.initializeNewCompiler(CompilerOptions)
	InitializeNewCompiler(options *CompilerOptions)
}

// Java: CompilerTypeTestCase.CLOSURE_DEFS
const CompilerTypeTestCase_CLOSURE_DEFS = `goog.inherits = function(x, y) {};
/** @type {!Function} */ goog.abstractMethod = function() {};
goog.isFunction = function(x) {};
goog.isObject = function(x) {};
/** @const */ goog.array = {};
// simplified ArrayLike definition
/**
 * @typedef {Array|{length: number}}
 */
goog.array.ArrayLike;
/**
 * @param {Array<T>|{length:number}} arr
 * @param {function(this:S, T, number, goog.array.ArrayLike):boolean} f
 * @param {S=} obj
 * @return {!Array<T>}
 * @template T,S
 */
// return empty array to satisfy return type
goog.array.filter = function(arr, f, obj){ return []; };
goog.asserts = {};
/** @return {*} */ goog.asserts.assert = function(obj, msg = undefined) { return obj; };
goog.loadModule = function(mod) {};
`

// CompilerTypeTestCase_DEFAULT_EXTERNS is a default set of externs for
// testing.
//
// TODO(bradfordcsmith): Replace this with externs built by
// TestExternsBuilder.
//
// Java: CompilerTypeTestCase.DEFAULT_EXTERNS
var CompilerTypeTestCase_DEFAULT_EXTERNS = CompilerTestCase_DEFAULT_EXTERNS

// CompilerTypeTestCaseBase holds the fields and base implementations of
// CompilerTypeTestCase.
type CompilerTypeTestCaseBase struct {
	self CompilerTypeTestCase

	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	compiler      *Compiler
	registry      *rhino.JSTypeRegistry
	errorReporter *rhinotesting.TestErrorReporter
}

// InitCompilerTypeTestCaseBase is the Go form of the implicit super().
//
// Java: CompilerTypeTestCase()
func (b *CompilerTypeTestCaseBase) InitCompilerTypeTestCaseBase(self CompilerTypeTestCase) {
	b.self = self
}

// GetCompilerTypeTestCaseBase gives access to the base class.
func (b *CompilerTypeTestCaseBase) GetCompilerTypeTestCaseBase() *CompilerTypeTestCaseBase {
	return b
}

// JUnitSetT is part of JUnitTestCase.
func (b *CompilerTypeTestCaseBase) JUnitSetT(t *testing.T) { b.t = t }

// T is the Go test running the current test method (Go-only: test-helper).
func (b *CompilerTypeTestCaseBase) T() *testing.T { return b.t }

// JUnitBefore runs the @Before method setUp().
func (b *CompilerTypeTestCaseBase) JUnitBefore() { b.self.SetUp() }

// JUnitAfter runs the @After method validateWarningsAndErrors().
func (b *CompilerTypeTestCaseBase) JUnitAfter() { b.self.ValidateWarningsAndErrors() }

// Java: CompilerTypeTestCase.defaultOptions()
func CompilerTypeTestCase_DefaultOptions() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetCodingConvention(NewGoogleCodingConvention())
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	options.SetWarningLevel(DiagnosticGroups_MISSING_PROPERTIES, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_MISPLACED_TYPE_ANNOTATION, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_INVALID_CASTS, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_JSDOC_MISSING_TYPE, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_BOUNDED_GENERICS, CheckLevel_WARNING)
	return options
}

// Java: CompilerTypeTestCase.getDefaultOptions()
func (b *CompilerTypeTestCaseBase) GetDefaultOptions() *CompilerOptions {
	return CompilerTypeTestCase_DefaultOptions()
}

// Java: CompilerTypeTestCase.getCodingConvention()
func (b *CompilerTypeTestCaseBase) GetCodingConvention() CodingConvention {
	return NewGoogleCodingConvention()
}

// Java: CompilerTypeTestCase.checkReportedWarningsHelper(String[])
func (b *CompilerTypeTestCaseBase) CheckReportedWarningsHelper(expected []string) {
	if expected == nil {
		expected = []string{}
	}

	AssertWithMessage("Regarding warnings:").
		That(b.compiler.GetWarnings()).
		ComparingElementsUsing(JSCompCorrespondences_DESCRIPTION_EQUALITY).
		ContainsExactlyElementsIn(expected).
		InOrder()
}

// Java: CompilerTypeTestCase.setUp()
func (b *CompilerTypeTestCaseBase) SetUp() {
	b.errorReporter = rhinotesting.NewTestErrorReporter()
	b.self.InitializeNewCompiler(b.self.GetDefaultOptions())
}

// Java: CompilerTypeTestCase.validateWarningsAndErrors()
func (b *CompilerTypeTestCaseBase) ValidateWarningsAndErrors() {
	b.errorReporter.VerifyHasEncounteredAllWarningsAndErrors(assertionTB{})
}

// Java: CompilerTypeTestCase.initializeNewCompiler(CompilerOptions)
func (b *CompilerTypeTestCaseBase) InitializeNewCompiler(options *CompilerOptions) {
	b.compiler = NewCompiler()
	b.compiler.InitOptions(options)
	b.compiler.MarkFeatureNotAllowed(parser.FeatureSet_Feature_MODULES)
	b.registry = b.compiler.GetTypeRegistry()
}

// Java: CompilerTypeTestCase.createUnionType(JSType...)
func (b *CompilerTypeTestCaseBase) CreateUnionType(variants ...rhino.JSType) rhino.JSType {
	return b.registry.CreateUnionType(variants...)
}

// Java: CompilerTypeTestCase.createRecordTypeBuilder()
func (b *CompilerTypeTestCaseBase) CreateRecordTypeBuilder() *rhino.RecordTypeBuilder {
	return rhino.NewRecordTypeBuilder(b.registry)
}

// Java: CompilerTypeTestCase.createNullableType(JSType)
func (b *CompilerTypeTestCaseBase) CreateNullableType(type_ rhino.JSType) rhino.JSType {
	return b.registry.CreateNullableType(type_)
}

// Java: CompilerTypeTestCase.createOptionalType(JSType)
func (b *CompilerTypeTestCaseBase) CreateOptionalType(type_ rhino.JSType) rhino.JSType {
	return b.registry.CreateOptionalType(type_)
}

// Java: CompilerTypeTestCase.createTemplatizedType(ObjectType, ImmutableList<JSType>)
func (b *CompilerTypeTestCaseBase) CreateTemplatizedTypeOfList(
	baseType rhino.ObjectType, templatizedTypes []rhino.JSType) *rhino.TemplatizedType {
	return b.registry.CreateTemplatizedTypeOfList(baseType, templatizedTypes)
}

// Java: CompilerTypeTestCase.createTemplatizedType(ObjectType, JSType...)
func (b *CompilerTypeTestCaseBase) CreateTemplatizedType(baseType rhino.ObjectType, templatizedType ...rhino.JSType) *rhino.TemplatizedType {
	return b.CreateTemplatizedTypeOfList(baseType, append([]rhino.JSType(nil), templatizedType...))
}

// AssertTypeEqualsNode asserts that a Node representing a type expression
// resolves to the correct JSType.
//
// Java: CompilerTypeTestCase.assertTypeEquals(JSType, Node)
func (b *CompilerTypeTestCaseBase) AssertTypeEqualsNode(expected rhino.JSType, actual rhino.Node) {
	b.AssertTypeEqualsExpr(expected, rhino.NewJSTypeExpression(actual, "<BaseJSTypeTestCase.java>"))
}

// AssertTypeEqualsExpr asserts that a a type expression resolves to the
// correct JSType.
//
// Java: CompilerTypeTestCase.assertTypeEquals(JSType, JSTypeExpression)
func (b *CompilerTypeTestCaseBase) AssertTypeEqualsExpr(expected rhino.JSType, actual *rhino.JSTypeExpression) {
	b.AssertTypeEquals(expected, b.Resolve(actual))
}

// Java: CompilerTypeTestCase.assertTypeEquals(JSType, JSType)
func (b *CompilerTypeTestCaseBase) AssertTypeEquals(a rhino.JSType, bType rhino.JSType) {
	TypeSubject_AssertType(bType).IsEqualTo(a)
}

// Java: CompilerTypeTestCase.assertTypeEquals(String, JSType, JSType)
func (b *CompilerTypeTestCaseBase) AssertTypeEqualsMsg(msg string, a rhino.JSType, bType rhino.JSType) {
	rhinotesting.NewTruth(assertionTB{}).AssertWithMessage("%s", msg).AboutTypes().That(bType).IsEqualTo(a)
}

// Resolve resolves a type expression, expecting the given warnings.
//
// Java: CompilerTypeTestCase.resolve(JSTypeExpression, String...)
func (b *CompilerTypeTestCaseBase) Resolve(n *rhino.JSTypeExpression, warnings ...string) rhino.JSType {
	b.errorReporter.ExpectAllWarnings(warnings...)
	return n.Evaluate(nil, b.registry)
}

// Java: CompilerTypeTestCase.getNativeNoObjectType()
func (b *CompilerTypeTestCaseBase) GetNativeNoObjectType() rhino.ObjectType {
	return b.GetNativeObjectTypeOf(rhino.JSTypeNative_NO_OBJECT_TYPE)
}

// Java: CompilerTypeTestCase.getNativeArrayType()
func (b *CompilerTypeTestCaseBase) GetNativeArrayType() rhino.ObjectType {
	return b.GetNativeObjectTypeOf(rhino.JSTypeNative_ARRAY_TYPE)
}

// Java: CompilerTypeTestCase.getNativeReadonlyArrayType()
func (b *CompilerTypeTestCaseBase) GetNativeReadonlyArrayType() rhino.ObjectType {
	return b.GetNativeObjectTypeOf(rhino.JSTypeNative_READONLY_ARRAY_TYPE)
}

// Java: CompilerTypeTestCase.getNativeStringObjectType()
func (b *CompilerTypeTestCaseBase) GetNativeStringObjectType() rhino.ObjectType {
	return b.GetNativeObjectTypeOf(rhino.JSTypeNative_STRING_OBJECT_TYPE)
}

// Java: CompilerTypeTestCase.getNativeNumberObjectType()
func (b *CompilerTypeTestCaseBase) GetNativeNumberObjectType() rhino.ObjectType {
	return b.GetNativeObjectTypeOf(rhino.JSTypeNative_NUMBER_OBJECT_TYPE)
}

// Java: CompilerTypeTestCase.getNativeBooleanObjectType()
func (b *CompilerTypeTestCaseBase) GetNativeBooleanObjectType() rhino.ObjectType {
	return b.GetNativeObjectTypeOf(rhino.JSTypeNative_BOOLEAN_OBJECT_TYPE)
}

// Java: CompilerTypeTestCase.getNativeNoType()
func (b *CompilerTypeTestCaseBase) GetNativeNoType() rhino.ObjectType {
	return b.GetNativeObjectTypeOf(rhino.JSTypeNative_NO_TYPE)
}

// Java: CompilerTypeTestCase.getNativeUnknownType()
func (b *CompilerTypeTestCaseBase) GetNativeUnknownType() rhino.ObjectType {
	return b.GetNativeObjectTypeOf(rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: CompilerTypeTestCase.getNativeCheckedUnknownType()
func (b *CompilerTypeTestCaseBase) GetNativeCheckedUnknownType() rhino.ObjectType {
	return b.GetNativeObjectTypeOf(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
}

// Java: CompilerTypeTestCase.getNativeObjectType()
func (b *CompilerTypeTestCaseBase) GetNativeObjectType() rhino.ObjectType {
	return b.GetNativeObjectTypeOf(rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: CompilerTypeTestCase.getNativeObjectType(JSTypeNative)
func (b *CompilerTypeTestCaseBase) GetNativeObjectTypeOf(jsTypeNative rhino.JSTypeNative) rhino.ObjectType {
	return b.registry.GetNativeObjectType(jsTypeNative)
}

// Java: CompilerTypeTestCase.getNativeObjectConstructorType()
func (b *CompilerTypeTestCaseBase) GetNativeObjectConstructorType() rhino.FunctionType {
	return b.GetNativeFunctionTypeOf(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE)
}

// Java: CompilerTypeTestCase.getNativeArrayConstructorType()
func (b *CompilerTypeTestCaseBase) GetNativeArrayConstructorType() rhino.FunctionType {
	return b.GetNativeFunctionTypeOf(rhino.JSTypeNative_ARRAY_FUNCTION_TYPE)
}

// Java: CompilerTypeTestCase.getNativeBooleanObjectConstructorType()
func (b *CompilerTypeTestCaseBase) GetNativeBooleanObjectConstructorType() rhino.FunctionType {
	return b.GetNativeFunctionTypeOf(rhino.JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE)
}

// Java: CompilerTypeTestCase.getNativeNumberObjectConstructorType()
func (b *CompilerTypeTestCaseBase) GetNativeNumberObjectConstructorType() rhino.FunctionType {
	return b.GetNativeFunctionTypeOf(rhino.JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE)
}

// Java: CompilerTypeTestCase.getNativeStringObjectConstructorType()
func (b *CompilerTypeTestCaseBase) GetNativeStringObjectConstructorType() rhino.FunctionType {
	return b.GetNativeFunctionTypeOf(rhino.JSTypeNative_STRING_OBJECT_FUNCTION_TYPE)
}

// Java: CompilerTypeTestCase.getNativeDateConstructorType()
func (b *CompilerTypeTestCaseBase) GetNativeDateConstructorType() rhino.FunctionType {
	return b.GetNativeFunctionTypeOf(rhino.JSTypeNative_DATE_FUNCTION_TYPE)
}

// Java: CompilerTypeTestCase.getNativeRegexpConstructorType()
func (b *CompilerTypeTestCaseBase) GetNativeRegexpConstructorType() rhino.FunctionType {
	return b.GetNativeFunctionTypeOf(rhino.JSTypeNative_REGEXP_FUNCTION_TYPE)
}

// Java: CompilerTypeTestCase.getNativeFunctionType()
func (b *CompilerTypeTestCaseBase) GetNativeFunctionType() rhino.FunctionType {
	return b.GetNativeFunctionTypeOf(rhino.JSTypeNative_FUNCTION_TYPE)
}

// Java: CompilerTypeTestCase.getNativeFunctionType(JSTypeNative)
func (b *CompilerTypeTestCaseBase) GetNativeFunctionTypeOf(jsTypeNative rhino.JSTypeNative) rhino.FunctionType {
	return b.registry.GetNativeFunctionType(jsTypeNative)
}

// Java: CompilerTypeTestCase.getNativeVoidType()
func (b *CompilerTypeTestCaseBase) GetNativeVoidType() rhino.JSType {
	return b.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
}

// Java: CompilerTypeTestCase.getNativeNullType()
func (b *CompilerTypeTestCaseBase) GetNativeNullType() rhino.JSType {
	return b.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
}

// Java: CompilerTypeTestCase.getNativeNullVoidType()
func (b *CompilerTypeTestCaseBase) GetNativeNullVoidType() rhino.JSType {
	return b.GetNativeType(rhino.JSTypeNative_NULL_VOID)
}

// Java: CompilerTypeTestCase.getNativeNumberType()
func (b *CompilerTypeTestCaseBase) GetNativeNumberType() rhino.JSType {
	return b.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: CompilerTypeTestCase.getNativeBooleanType()
func (b *CompilerTypeTestCaseBase) GetNativeBooleanType() rhino.JSType {
	return b.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
}

// Java: CompilerTypeTestCase.getNativeStringType()
func (b *CompilerTypeTestCaseBase) GetNativeStringType() rhino.JSType {
	return b.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
}

// Java: CompilerTypeTestCase.getNativeObjectNumberStringBooleanType()
func (b *CompilerTypeTestCaseBase) GetNativeObjectNumberStringBooleanType() rhino.JSType {
	return b.registry.CreateUnionTypeOfNatives(rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_BOOLEAN_TYPE)
}

// Java: CompilerTypeTestCase.getNativeNumberStringBooleanType()
func (b *CompilerTypeTestCaseBase) GetNativeNumberStringBooleanType() rhino.JSType {
	return b.GetNativeType(rhino.JSTypeNative_NUMBER_STRING_BOOLEAN)
}

// Java: CompilerTypeTestCase.getNativeObjectNumberStringBooleanSymbolType()
func (b *CompilerTypeTestCaseBase) GetNativeObjectNumberStringBooleanSymbolType() rhino.JSType {
	return b.registry.CreateUnionTypeOfNatives(
		rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_BOOLEAN_TYPE, rhino.JSTypeNative_SYMBOL_TYPE)
}

// Java: CompilerTypeTestCase.getNativeValueTypes()
func (b *CompilerTypeTestCaseBase) GetNativeValueTypes() rhino.JSType {
	return b.GetNativeType(rhino.JSTypeNative_VALUE_TYPES)
}

// Java: CompilerTypeTestCase.getNativeAllType()
func (b *CompilerTypeTestCaseBase) GetNativeAllType() rhino.JSType {
	return b.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
}

// Java: CompilerTypeTestCase.getNativeType(JSTypeNative)
func (b *CompilerTypeTestCaseBase) GetNativeType(jsTypeNative rhino.JSTypeNative) rhino.JSType {
	return b.registry.GetNativeType(jsTypeNative)
}

// Port of com/google/javascript/jscomp/CheckClosureImportsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/CheckClosureImportsTest.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Tests for CheckClosureImports.
type CheckClosureImportsTest struct {
	CompilerTestCaseBase

	moduleType ModuleMetadataMap_ModuleType
	// use the static ModuleMetadata fields rather than rebuilding the
	// ModuleMetadataMap
	useDefaultModuleMetadata bool
}

const checkClosureImportsTest_PROVIDES_SYMBOL_PATH = "symbol.js"

var checkClosureImportsTest_PROVIDES_SYMBOL_SRC = SourceFile_FromCode(checkClosureImportsTest_PROVIDES_SYMBOL_PATH, "goog.provide('symbol');")

const checkClosureImportsTest_ES_MODULE_PATH = "es_module.js"

var checkClosureImportsTest_ES_MODULE_SRC = SourceFile_FromCode(checkClosureImportsTest_ES_MODULE_PATH, "goog.declareModuleId('es.module');\nexport {};")

const checkClosureImportsTest_TEST_CODE_PATH = "testcode"

var checkClosureImportsTest_PROVIDES_SYMBOL_METADATA = ModuleMetadataMap_ModuleMetadata_Builder_().
	AddGoogNamespace("symbol").
	ModuleType(ModuleMetadataMap_ModuleType_GOOG_PROVIDE).
	UsesClosure(true).
	IsTestOnly(false).
	Build()

var checkClosureImportsTest_ES_MODULE_METADATA = ModuleMetadataMap_ModuleMetadata_Builder_().
	AddGoogNamespace("es.module").
	ModuleType(ModuleMetadataMap_ModuleType_ES6_MODULE).
	UsesClosure(true).
	IsTestOnly(false).
	Build()

var checkClosureImportsTest_EXTERN_METADATA = ModuleMetadataMap_ModuleMetadata_Builder_().
	ModuleType(ModuleMetadataMap_ModuleType_SCRIPT).
	UsesClosure(false).
	IsTestOnly(false).
	Build()

// Java: CheckClosureImportsTest.makeTestFile(String)
func checkClosureImportsTest_makeTestFile(text string) *SourceFile {
	return SourceFile_FromCode(checkClosureImportsTest_TEST_CODE_PATH, text)
}

func NewCheckClosureImportsTest() *CheckClosureImportsTest {
	c := &CheckClosureImportsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckClosureImportsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckClosureImportsTest() })
}

// JUnitBefore runs the superclass @Before (setUp) then reset().
func (c *CheckClosureImportsTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.Reset()
}

// Java: CheckClosureImportsTest.reset()
func (c *CheckClosureImportsTest) Reset() {
	c.moduleType = ModuleMetadataMap_ModuleType_SCRIPT
	c.useDefaultModuleMetadata = true
}

// Java: CheckClosureImportsTest.getOptions()
func (c *CheckClosureImportsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_ERROR)
	return options
}

// Java: CheckClosureImportsTest.getProcessor(Compiler)
func (c *CheckClosureImportsTest) GetProcessor(compiler *Compiler) CompilerPass {
	testMetadata :=
		ModuleMetadataMap_ModuleMetadata_Builder_().
			AddGoogNamespace("test").
			ModuleType(c.moduleType).
			UsesClosure(true).
			IsTestOnly(false).
			Build()

	var metadata *ModuleMetadataMap
	if c.useDefaultModuleMetadata {
		byPath := jutil.NewLinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]()
		byPath.Put(checkClosureImportsTest_PROVIDES_SYMBOL_PATH, checkClosureImportsTest_PROVIDES_SYMBOL_METADATA)
		byPath.Put(checkClosureImportsTest_ES_MODULE_PATH, checkClosureImportsTest_ES_MODULE_METADATA)
		byPath.Put(checkClosureImportsTest_TEST_CODE_PATH, testMetadata)
		byPath.Put("externs", checkClosureImportsTest_EXTERN_METADATA)
		byNamespace := jutil.NewLinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]()
		byNamespace.Put("symbol", checkClosureImportsTest_PROVIDES_SYMBOL_METADATA)
		byNamespace.Put("es.module", checkClosureImportsTest_ES_MODULE_METADATA)
		byNamespace.Put("test", testMetadata)
		metadata = NewModuleMetadataMap(byPath, byNamespace)
	} else {
		metadata = compiler.GetModuleMetadataMap()
	}

	return NewCheckClosureImports(compiler, metadata)
}

// TestCommonCase is a test common to goog.require, goog.forwardDeclare, and
// goog.requireType.
//
// Go-only: test-helper — the two Java overloads testCommonCase(String) and
// testCommonCase(String, DiagnosticType) are one variadic helper; it is
// unexported-by-signature (it takes parameters) so RunJUnit4 skips it.
//
// Java: CheckClosureImportsTest.testCommonCase(String)
// Java: CheckClosureImportsTest.testCommonCase(String, DiagnosticType)
func (c *CheckClosureImportsTest) TestCommonCase(source string, error_ ...*DiagnosticType) {
	if len(error_) == 0 {
		c.Test(c.Srcs(checkClosureImportsTest_PROVIDES_SYMBOL_SRC, checkClosureImportsTest_makeTestFile(jutil.StringReplace(source, "<import>", "goog.require"))))
		c.Test(
			c.Srcs(checkClosureImportsTest_PROVIDES_SYMBOL_SRC, checkClosureImportsTest_makeTestFile(jutil.StringReplace(source, "<import>", "goog.forwardDeclare"))))
		c.Test(c.Srcs(checkClosureImportsTest_PROVIDES_SYMBOL_SRC, checkClosureImportsTest_makeTestFile(jutil.StringReplace(source, "<import>", "goog.requireType"))))
		return
	}
	c.Test(
		c.Srcs(checkClosureImportsTest_PROVIDES_SYMBOL_SRC, checkClosureImportsTest_makeTestFile(jutil.StringReplace(source, "<import>", "goog.require"))),
		c.Error(error_[0]))
	c.Test(
		c.Srcs(checkClosureImportsTest_PROVIDES_SYMBOL_SRC, checkClosureImportsTest_makeTestFile(jutil.StringReplace(source, "<import>", "goog.forwardDeclare"))),
		c.Error(error_[0]))
	c.Test(
		c.Srcs(checkClosureImportsTest_PROVIDES_SYMBOL_SRC, checkClosureImportsTest_makeTestFile(jutil.StringReplace(source, "<import>", "goog.requireType"))),
		c.Error(error_[0]))
}

// Java: CheckClosureImportsTest.referenceMissingSymbolIsError()
func (c *CheckClosureImportsTest) Test_referenceMissingSymbolIsError() {
	c.TestCommonCase("<import>('symbol');")
	c.TestSame("() => goog.module.get('symbol');")

	// Not an error for goog.forwardDeclare in scripts.
	c.TestError("goog.require('dne');", ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE)
	c.TestError("goog.requireType('dne');", ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE)
	c.TestError("goog.requireDynamic('dne');", ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE)
	c.TestError("() => goog.module.get('dne');", ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE)
}

// Java: CheckClosureImportsTest.importInGlobalScopeIsOk()
func (c *CheckClosureImportsTest) Test_importInGlobalScopeIsOk() {
	c.TestCommonCase("<import>('symbol');")
}

// Java: CheckClosureImportsTest.importInGoogProvideIsOk()
func (c *CheckClosureImportsTest) Test_importInGoogProvideIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_PROVIDE

	c.TestCommonCase(
		`goog.provide('test');
<import>('symbol');
`)
}

// Java: CheckClosureImportsTest.importInGoogModuleScopeIsOk()
func (c *CheckClosureImportsTest) Test_importInGoogModuleScopeIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.TestCommonCase(
		`goog.module('test');
<import>('symbol');
`)
}

// Java: CheckClosureImportsTest.importInGoogLoadModuleScopeIsOk()
func (c *CheckClosureImportsTest) Test_importInGoogLoadModuleScopeIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.TestCommonCase(
		`goog.loadModule(function(exports) {
  goog.module('test');
  <import>('symbol');
});
`)
}

// Java: CheckClosureImportsTest.importInEsModuleScopeIsOk()
func (c *CheckClosureImportsTest) Test_importInEsModuleScopeIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE

	c.TestCommonCase(
		`<import>('symbol');
export {};
`)
}

// Java: CheckClosureImportsTest.importInBlockScopeInGlobalHoistScopeIsOk()
func (c *CheckClosureImportsTest) Test_importInBlockScopeInGlobalHoistScopeIsOk() {
	c.TestCommonCase("{ <import>('symbol'); }")
}

// Java: CheckClosureImportsTest.importInBlockScopeInModuleHoistScopeIsError()
func (c *CheckClosureImportsTest) Test_importInBlockScopeInModuleHoistScopeIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE
	c.TestCommonCase(
		`goog.module('test');
{ <import>('symbol'); }
`,
		ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR)

	c.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE

	c.TestCommonCase(
		`{ <import>('symbol'); }
export {};
`,
		ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR)
}

// Java: CheckClosureImportsTest.importInFunctionScopeIsError()
func (c *CheckClosureImportsTest) Test_importInFunctionScopeIsError() {
	c.TestCommonCase("() => <import>('symbol');", ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR)
}

// Java: CheckClosureImportsTest.importInSingleDeclarationIsOk()
func (c *CheckClosureImportsTest) Test_importInSingleDeclarationIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.TestCommonCase(
		`goog.module('test');
var symbol = <import>('symbol');
`)

	c.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE

	c.TestCommonCase(
		`const symbol = <import>('symbol');
export {};
`)
}

// Java: CheckClosureImportsTest.mustMatchCapitalization()
func (c *CheckClosureImportsTest) Test_mustMatchCapitalization() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.TestCommonCase(
		`goog.module('test');
const Symbol = <import>('symbol');
`,
		ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION)
}

// Java: CheckClosureImportsTest.canHaveDifferentName()
func (c *CheckClosureImportsTest) Test_canHaveDifferentName() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.TestCommonCase(
		`goog.module('test');
const localSymbol = <import>('symbol');
`)
}

// Java: CheckClosureImportsTest.importInMultiDeclarationisError()
func (c *CheckClosureImportsTest) Test_importInMultiDeclarationisError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.TestCommonCase(
		`goog.module('test');
var before, symbol = goog.require('symbol');
`,
		CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION)

	c.TestCommonCase(
		`goog.module('test');
var symbol = goog.require('symbol'), after;
`,
		CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION)

	c.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE

	c.TestCommonCase(
		`const before = 0, symbol = goog.require('symbol');
export {};
`,
		CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION)

	c.TestCommonCase(
		`var symbol = goog.require('symbol'), after = 0;
export {};
`,
		CheckClosureImports_ONE_CLOSURE_IMPORT_PER_DECLARATION)
}

// Java: CheckClosureImportsTest.googRequireDeclarationInEsModuleMustBeConst()
func (c *CheckClosureImportsTest) Test_googRequireDeclarationInEsModuleMustBeConst() {
	c.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE

	c.TestSame(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`const symbol = goog.require('symbol');
export {};
`)))

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`let symbol = goog.require('symbol');
export {};
`)),
		c.Error(CheckClosureImports_LET_CLOSURE_IMPORT))

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`var symbol = goog.require('symbol');
export {};
`)),
		c.Error(CheckClosureImports_LHS_OF_CLOSURE_IMPORT_MUST_BE_CONST_IN_ES_MODULE))
}

// Java: CheckClosureImportsTest.googRequireDeclarationCannotBeLet()
func (c *CheckClosureImportsTest) Test_googRequireDeclarationCannotBeLet() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
let symbol = goog.require('symbol');
`)),
		c.Error(CheckClosureImports_LET_CLOSURE_IMPORT))
}

// Java: CheckClosureImportsTest.googRequireTypeDeclarationCannotBeLet()
func (c *CheckClosureImportsTest) Test_googRequireTypeDeclarationCannotBeLet() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
let symbol = goog.requireType('symbol');
`)),
		c.Error(CheckClosureImports_LET_CLOSURE_IMPORT))
}

// Java: CheckClosureImportsTest.googRequireSimpleObjectDestructuringIsOk()
func (c *CheckClosureImportsTest) Test_googRequireSimpleObjectDestructuringIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {property0, property1: p1} = goog.require('symbol');
`)))
}

// Java: CheckClosureImportsTest.googRequireTypeSimpleObjectDestructuringIsOk()
func (c *CheckClosureImportsTest) Test_googRequireTypeSimpleObjectDestructuringIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {property0, property1: p1} = goog.requireType('symbol');
`)))
}

// Java: CheckClosureImportsTest.googForwardDeclareSimpleObjectDestructuringIsError()
func (c *CheckClosureImportsTest) Test_googForwardDeclareSimpleObjectDestructuringIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {property0} = goog.forwardDeclare('symbol');
`)),
		c.Error(CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING))
}

// Java: CheckClosureImportsTest.googRequireComplexObjectDestructuringIsError()
func (c *CheckClosureImportsTest) Test_googRequireComplexObjectDestructuringIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {property0: {a}} = goog.require('symbol');
`)),
		c.Error(CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING))
}

// Java: CheckClosureImportsTest.googRequireTypeComplexObjectDestructuringIsError()
func (c *CheckClosureImportsTest) Test_googRequireTypeComplexObjectDestructuringIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {property0: {a}} = goog.requireType('symbol');
`)),
		c.Error(CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING))
}

// Java: CheckClosureImportsTest.googRequireSimpleArrayDestructuringIsError()
func (c *CheckClosureImportsTest) Test_googRequireSimpleArrayDestructuringIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const [e0] = goog.require('symbol');
`)),
		c.Error(CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING))
}

// Java: CheckClosureImportsTest.googRequireTypeSimpleArrayDestructuringIsError()
func (c *CheckClosureImportsTest) Test_googRequireTypeSimpleArrayDestructuringIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const [e0] = goog.requireType('symbol');
`)),
		c.Error(CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING))
}

// Java: CheckClosureImportsTest.googRequireDefaultValueDestructuringIsError()
func (c *CheckClosureImportsTest) Test_googRequireDefaultValueDestructuringIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {property0 = 1} = goog.require('symbol');
`)),
		c.Error(CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING))
}

// Java: CheckClosureImportsTest.googRequireTypeDefaultValueDestructuringIsError()
func (c *CheckClosureImportsTest) Test_googRequireTypeDefaultValueDestructuringIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {property0 = 1} = goog.requireType('symbol');
`)),
		c.Error(CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING))
}

// Java: CheckClosureImportsTest.googRequireComputedPropertyDestructuringIsError()
func (c *CheckClosureImportsTest) Test_googRequireComputedPropertyDestructuringIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {[property0]: a} = goog.require('symbol');
`)),
		c.Error(CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING))
}

// Java: CheckClosureImportsTest.googRequireObjectRestDestructuringIsError()
func (c *CheckClosureImportsTest) Test_googRequireObjectRestDestructuringIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {...rest} = goog.require('symbol');
`)),
		c.Error(CheckClosureImports_INVALID_CLOSURE_IMPORT_DESTRUCTURING))
}

// Java: CheckClosureImportsTest.googForwardDeclareNestedObjectDestructuringIsError()
func (c *CheckClosureImportsTest) Test_googForwardDeclareNestedObjectDestructuringIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {property0: {a}} = goog.forwardDeclare('symbol');
`)),
		c.Error(CheckClosureImports_NO_CLOSURE_IMPORT_DESTRUCTURING))
}

// Java: CheckClosureImportsTest.googRequireLateIsError()
func (c *CheckClosureImportsTest) Test_googRequireLateIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(c.Srcs(checkClosureImportsTest_PROVIDES_SYMBOL_SRC, checkClosureImportsTest_makeTestFile("goog.require('symbol');")))
	c.Test(
		c.Srcs(checkClosureImportsTest_makeTestFile("goog.require('symbol');"), checkClosureImportsTest_PROVIDES_SYMBOL_SRC),
		c.Error(CheckClosureImports_LATE_PROVIDE_ERROR))
}

// Java: CheckClosureImportsTest.googRequireTypeLateIsOk()
func (c *CheckClosureImportsTest) Test_googRequireTypeLateIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(c.Srcs(checkClosureImportsTest_PROVIDES_SYMBOL_SRC, checkClosureImportsTest_makeTestFile("goog.requireType('symbol');")))

	c.Test(c.Srcs(checkClosureImportsTest_makeTestFile("goog.requireType('symbol');"), checkClosureImportsTest_PROVIDES_SYMBOL_SRC))
}

// Java: CheckClosureImportsTest.googForwardDeclareLateIsError()
func (c *CheckClosureImportsTest) Test_googForwardDeclareLateIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(c.Srcs(checkClosureImportsTest_PROVIDES_SYMBOL_SRC, checkClosureImportsTest_makeTestFile("goog.forwardDeclare('symbol');")))

	c.Test(c.Srcs(checkClosureImportsTest_makeTestFile("goog.forwardDeclare('symbol');"), checkClosureImportsTest_PROVIDES_SYMBOL_SRC))
}

// Java: CheckClosureImportsTest.moduleGetInModuleScopeIsError()
func (c *CheckClosureImportsTest) Test_moduleGetInModuleScopeIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
goog.module.get('symbol');
`)),
		c.Error(ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET))

	c.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module.get('symbol');
export {};
`)),
		c.Error(ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET))
}

// Java: CheckClosureImportsTest.moduleGetInGlobalScopeIsOK()
func (c *CheckClosureImportsTest) Test_moduleGetInGlobalScopeIsOK() {
	c.moduleType = ModuleMetadataMap_ModuleType_SCRIPT

	// only global variable aliasing goog.module.get results are banned.
	c.Test(c.Srcs(checkClosureImportsTest_ES_MODULE_SRC, checkClosureImportsTest_makeTestFile("goog.module.get('es.module');")))
}

// Java: CheckClosureImportsTest.moduleGetInFunctionScopeIsOk()
func (c *CheckClosureImportsTest) Test_moduleGetInFunctionScopeIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
() => goog.module.get('symbol');
`)))

	c.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`() => goog.module.get('symbol');
export {};
`)))

	c.moduleType = ModuleMetadataMap_ModuleType_SCRIPT

	c.Test(c.Srcs(checkClosureImportsTest_PROVIDES_SYMBOL_SRC, checkClosureImportsTest_makeTestFile("() => goog.module.get('symbol');")))
}

// Java: CheckClosureImportsTest.moduleGetFillInForwardDeclareIsOk()
func (c *CheckClosureImportsTest) Test_moduleGetFillInForwardDeclareIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
let symbol = goog.forwardDeclare('symbol');
() => symbol = goog.module.get('symbol');
`)))
}

// Java: CheckClosureImportsTest.moduleGetAssignmentToNamespace()
func (c *CheckClosureImportsTest) Test_moduleGetAssignmentToNamespace() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_PROVIDE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.provide('a.b.c');
goog.require('symbol');
a.b.c = goog.module.get('symbol');
`)))
}

// Java: CheckClosureImportsTest.moduleGetAssignmentToName()
func (c *CheckClosureImportsTest) Test_moduleGetAssignmentToName() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_PROVIDE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.provide('test');
var sym = goog.module.get('symbol');
`)),
		c.Error(ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE))

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.provide('test');
var sym = goog.module.get('symbol').sym;
`)),
		c.Error(ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE))

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.provide('test');
var sym = goog.module.get('symbol')?.sym;
`)),
		c.Error(ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE))

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.provide('test');
const {sym} = goog.module.get('symbol');
`)),
		c.Error(ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE))

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.provide('test');
var sym
sym = goog.module.get('symbol').sym;
`)),
		c.Error(ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE))

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.provide('test');
let sym; ({sym} = goog.module.get('symbol'));
`)),
		c.Error(ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE))
}

// Java: CheckClosureImportsTest.moduleGetOtherAssignmentDeclareIsError()
func (c *CheckClosureImportsTest) Test_moduleGetOtherAssignmentDeclareIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
let symbol;
() => symbol = goog.module.get('symbol');
`)),
		c.Error(ClosureRewriteModule_INVALID_GET_ALIAS))
}

// Java: CheckClosureImportsTest.moduleGetFillInWrongNamespaceIsError()
func (c *CheckClosureImportsTest) Test_moduleGetFillInWrongNamespaceIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
let symbol = goog.forwardDeclare('symbol');
() => symbol = goog.module.get('test');
`)),
		c.Error(ClosureRewriteModule_INVALID_GET_ALIAS))
}

// Java: CheckClosureImportsTest.googForwardDeclareForGlobalInScriptIsOk()
func (c *CheckClosureImportsTest) Test_googForwardDeclareForGlobalInScriptIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_SCRIPT

	c.Test(c.Srcs(checkClosureImportsTest_makeTestFile("goog.forwardDeclare('MyGlobal');")))
}

// Java: CheckClosureImportsTest.googForwardDeclareForGlobalInProvideIsOk()
func (c *CheckClosureImportsTest) Test_googForwardDeclareForGlobalInProvideIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_PROVIDE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_makeTestFile(
				`goog.provide('test');
goog.forwardDeclare('MyGlobal');
`)))
}

// Java: CheckClosureImportsTest.googForwardDeclareForGlobalNotAliasedInGoogModuleIsOk()
func (c *CheckClosureImportsTest) Test_googForwardDeclareForGlobalNotAliasedInGoogModuleIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
goog.forwardDeclare('MyGlobal');
`)))
}

// Java: CheckClosureImportsTest.googForwardDeclareForGlobalAliasedInGoogModuleIsError()
func (c *CheckClosureImportsTest) Test_googForwardDeclareForGlobalAliasedInGoogModuleIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const MyGlobal = goog.forwardDeclare('MyGlobal');
`)),
		c.Error(ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE))
}

// Java: CheckClosureImportsTest.googForwardDeclareForGlobalNotAliasedInEsModuleIsOk()
func (c *CheckClosureImportsTest) Test_googForwardDeclareForGlobalNotAliasedInEsModuleIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_makeTestFile(
				`export {};
goog.forwardDeclare('MyGlobal');
`)))
}

// Java: CheckClosureImportsTest.googForwardDeclareForGlobalAliasedInEsModuleIsError()
func (c *CheckClosureImportsTest) Test_googForwardDeclareForGlobalAliasedInEsModuleIsError() {
	c.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_makeTestFile(
				`export {};
const MyGlobal = goog.forwardDeclare('MyGlobal');
`)),
		c.Error(ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE))
}

// Java: CheckClosureImportsTest.googRequireCodeReferenceAliasIsOk()
func (c *CheckClosureImportsTest) Test_googRequireCodeReferenceAliasIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.TestSame(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {alias} = goog.require('symbol');
alias;
`)))
}

// Java: CheckClosureImportsTest.googForwardDeclareCodeReferenceAliasIsOk()
func (c *CheckClosureImportsTest) Test_googForwardDeclareCodeReferenceAliasIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.TestSame(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const alias = goog.forwardDeclare('symbol');
alias;
`)))
}

// Java: CheckClosureImportsTest.googRequireTypeTypeReferenceAliasIsOk()
func (c *CheckClosureImportsTest) Test_googRequireTypeTypeReferenceAliasIsOk() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
const {Type} = goog.requireType('symbol');
/** !Type */ let x;
`)))
}

// Java: CheckClosureImportsTest.googRequireBetweenEsModulesIsWarning()
func (c *CheckClosureImportsTest) Test_googRequireBetweenEsModulesIsWarning() {
	c.moduleType = ModuleMetadataMap_ModuleType_ES6_MODULE

	c.Test(
		c.Srcs(
			checkClosureImportsTest_ES_MODULE_SRC,
			checkClosureImportsTest_makeTestFile(
				`const {Type} = goog.require('es.module');
export {};
`)),
		c.Warning(Es6RewriteModules_SHOULD_IMPORT_ES6_MODULE))
}

// Java: CheckClosureImportsTest.testIgnoreArguments()
func (c *CheckClosureImportsTest) TestIgnoreArguments() {
	c.moduleType = ModuleMetadataMap_ModuleType_GOOG_MODULE
	c.Test(
		c.Srcs(
			checkClosureImportsTest_PROVIDES_SYMBOL_SRC,
			checkClosureImportsTest_makeTestFile(
				`goog.module('test');
function fn() {
  return arguments;
}
`)))
}

// Java: CheckClosureImportsTest.crossChunkGoogRequire_ofGoogModule_isError()
func (c *CheckClosureImportsTest) Test_crossChunkGoogRequire_ofGoogModule_isError() {
	c.useDefaultModuleMetadata = false
	c.EnableCreateModuleMap()

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("// base").
				AddChunk("goog.module('test');").
				AddChunk("goog.require('test');").
				Build()),
		c.Warning(CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR))
}

// Java: CheckClosureImportsTest.crossChunkGoogRequire_ofBundledGoogModule_isError()
func (c *CheckClosureImportsTest) Test_crossChunkGoogRequire_ofBundledGoogModule_isError() {
	c.useDefaultModuleMetadata = false
	c.EnableCreateModuleMap()

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("// base").
				AddChunk(
					`goog.loadModule(function(exports) {
goog.module('test');
return exports;
});
`).
				AddChunk("goog.require('test');").
				Build()),
		c.Warning(CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR))
}

// Java: CheckClosureImportsTest.crossChunkGoogRequire_ofGoogProvide_isError()
func (c *CheckClosureImportsTest) Test_crossChunkGoogRequire_ofGoogProvide_isError() {
	c.useDefaultModuleMetadata = false
	c.EnableCreateModuleMap()

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("// base").
				AddChunk("goog.provide('test');").
				AddChunk("goog.require('test');").
				Build()),
		c.Warning(CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR))
}

// Java: CheckClosureImportsTest.crossChunkGoogRequireType_isOk()
func (c *CheckClosureImportsTest) Test_crossChunkGoogRequireType_isOk() {
	c.useDefaultModuleMetadata = false
	c.EnableCreateModuleMap()

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("// base").
				AddChunk("goog.module('test');").
				AddChunk("goog.requireType('test');").
				Build()))
}

// Java: CheckClosureImportsTest.crossChunkGoogForwardDeclare_isOk()
func (c *CheckClosureImportsTest) Test_crossChunkGoogForwardDeclare_isOk() {
	c.useDefaultModuleMetadata = false
	c.EnableCreateModuleMap()

	c.TestNoWarning(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("// base").
				AddChunk("goog.module('test');").
				AddChunk("goog.forwardDeclare('test');").
				Build()))
}

// Java: CheckClosureImportsTest.crossChunkGoogRequire_withExplicitChunkDependency_isOk()
func (c *CheckClosureImportsTest) Test_crossChunkGoogRequire_withExplicitChunkDependency_isOk() {
	c.useDefaultModuleMetadata = false
	c.EnableCreateModuleMap()

	c.TestNoWarning(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("goog.module('test');").
				AddChunk("goog.require('test');").
				Build()))
}

// Port of com/google/javascript/jscomp/RewriteDynamicImportsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RewriteDynamicImportsTest.java)

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type RewriteDynamicImportsTest struct {
	CompilerTestCaseBase

	dynamicImportAlias *string
	language           *CompilerOptions_LanguageMode
	languageIn         *CompilerOptions_LanguageMode
	chunkOutputType    CompilerOptions_ChunkOutputType
}

// Java: RewriteDynamicImportsTest()
func NewRewriteDynamicImportsTest() *RewriteDynamicImportsTest {
	c := &RewriteDynamicImportsTest{}
	c.dynamicImportAlias = jutil.StringPtr("imprt_")
	c.language = nil
	c.languageIn = nil
	c.chunkOutputType = CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE
	c.InitCompilerTestCaseBaseWithExterns(c, NewTestExternsBuilder().AddPromise().Build())
	return c
}

func TestRewriteDynamicImportsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteDynamicImportsTest() })
}

// compilerOptionsLanguageModePtr is a @Nullable LanguageMode value.
//
// Go-only: test-helper.
func compilerOptionsLanguageModePtr(m CompilerOptions_LanguageMode) *CompilerOptions_LanguageMode {
	return &m
}

// Java: RewriteDynamicImportsTest.setUp()
func (c *RewriteDynamicImportsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableCreateModuleMap()
	c.EnableTypeInfoValidation()
}

// Java: RewriteDynamicImportsTest.getOptions()
func (c *RewriteDynamicImportsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetPrettyPrint(true)
	if c.language != nil {
		options.SetLanguage(*c.language)
	}
	if c.languageIn != nil {
		options.SetLanguageIn(*c.languageIn)
	}
	return options
}

// Java: RewriteDynamicImportsTest.getProcessor(Compiler)
func (c *RewriteDynamicImportsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		// NOTE: we cannot just use enableTypeCheck(), because we need to be able
		// to retrieve globalTypedScope for use by Es6RewriteModules.
		var rai ReverseAbstractInterpreter = NewSemanticReverseAbstractInterpreter(compiler.GetTypeRegistry())
		globalTypedScope :=
			jutil.CheckNotNull(
				NewTypeCheck(compiler, rai, compiler.GetTypeRegistry()).
					ProcessForTesting(externs, root))
		compiler.SetTypeCheckingHasRun(true)
		// We need to make sure modules are rewritten, because RewriteDynamicImports
		// expects to be able to see the module variables.
		NewEs6RewriteModules(
			compiler,
			compiler.GetModuleMetadataMap(),
			compiler.GetModuleMap(),
			/* preprocessorSymbolTable= */ nil,
			globalTypedScope).
			Process(externs, root)

		NewRewriteDynamicImports(compiler, c.dynamicImportAlias, c.chunkOutputType).
			Process(externs, root)
	})
}

// Java: RewriteDynamicImportsTest.externalImportWithoutAlias()
func (c *RewriteDynamicImportsTest) TestExternalImportWithoutAlias() {
	c.dynamicImportAlias = nil
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.TestSame("import('./external.js')")
}

// Java: RewriteDynamicImportsTest.externalImportWithAlias()
func (c *RewriteDynamicImportsTest) TestExternalImportWithAlias() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		c.Externs("/** @param {string} a @return {!Promise<?>} */ function imprt_(a) {}"),
		c.Srcs("import('./external.js')"),
		c.Expected("imprt_('./external.js')"))
}

// Java: RewriteDynamicImportsTest.internalImportSameChunkWithoutAlias_unused()
func (c *RewriteDynamicImportsTest) TestInternalImportSameChunkWithoutAlias_unused() {
	c.dynamicImportAlias = nil
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))
	actualChunk0.AddSourceFile(SourceFile_FromCode("i1.js", "import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode("i1.js", "Promise.resolve();")})

	c.Test(c.Srcs([]*JSChunk{actualChunk0}), expectedSrcs)
}

// Java: RewriteDynamicImportsTest.internalImportSameChunkWithoutAlias_used()
func (c *RewriteDynamicImportsTest) TestInternalImportSameChunkWithoutAlias_used() {
	c.dynamicImportAlias = nil
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))
	actualChunk0.AddSourceFile(SourceFile_FromCode("i1.js", "const ns = import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode("i1.js", "const ns = Promise.resolve(module$i0);")})

	c.Test(c.Srcs([]*JSChunk{actualChunk0}), expectedSrcs)
}

// Java: RewriteDynamicImportsTest.internalImportSameChunkWithAlias()
func (c *RewriteDynamicImportsTest) TestInternalImportSameChunkWithAlias() {
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))
	actualChunk0.AddSourceFile(SourceFile_FromCode("i1.js", "const ns = import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode("i1.js", "const ns = Promise.resolve(module$i0);")})

	c.Test(
		c.Externs(NewTestExternsBuilder().AddObject().AddFunction().AddPromise().Build()),
		c.Srcs([]*JSChunk{actualChunk0}),
		expectedSrcs)
}

// Java: RewriteDynamicImportsTest.internalImportDifferentChunksWithoutAlias_unused()
func (c *RewriteDynamicImportsTest) TestInternalImportDifferentChunksWithoutAlias_unused() {
	c.dynamicImportAlias = nil
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode("i1.js", "import('./chunk0.js');")})

	c.Test(c.Srcs([]*JSChunk{actualChunk0, actualChunk1}), expectedSrcs)
}

// Java: RewriteDynamicImportsTest.internalImportDifferentChunksWithAlias_unused()
func (c *RewriteDynamicImportsTest) TestInternalImportDifferentChunksWithAlias_unused() {
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode("i1.js", "imprt_('./chunk0.js');")})

	c.Test(c.Srcs([]*JSChunk{actualChunk0, actualChunk1}), expectedSrcs)
}

// Java: RewriteDynamicImportsTest.internalImportDifferentChunksWithoutAlias_used()
func (c *RewriteDynamicImportsTest) TestInternalImportDifferentChunksWithoutAlias_used() {
	c.dynamicImportAlias = nil
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "const nsPromise = import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode(
					"i1.js",
					`const nsPromise =
    import('./chunk0.js').then(() => module$i0);
`)})

	c.Test(c.Srcs([]*JSChunk{actualChunk0, actualChunk1}), expectedSrcs)
}

// Java: RewriteDynamicImportsTest.internalImportDifferentChunksWithAlias_used()
func (c *RewriteDynamicImportsTest) TestInternalImportDifferentChunksWithAlias_used() {
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "const nsPromise = import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode(
					"i1.js", "const nsPromise = imprt_('./chunk0.js').then(() => module$i0);")})

	c.Test(c.Srcs([]*JSChunk{actualChunk0, actualChunk1}), expectedSrcs)
}

// Java: RewriteDynamicImportsTest.qualifiedNameAlias()
func (c *RewriteDynamicImportsTest) TestQualifiedNameAlias() {
	c.dynamicImportAlias = jutil.StringPtr("ns.imprt_")
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		`/** @const */
const ns = {};
/** @const */
ns.imprt_ = function(path) {};
import('./other.js');
`,
		`/** @const */
const ns = {};
/** @const */
ns.imprt_ = function(path) {};
ns.imprt_('./other.js');
`)
}

// Java: RewriteDynamicImportsTest.qualifiedNameAliasExtern()
func (c *RewriteDynamicImportsTest) TestQualifiedNameAliasExtern() {
	c.dynamicImportAlias = jutil.StringPtr("ns.imprt_")
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		c.Externs("const ns = {}; /** @const */ ns.imprt_ = function(path) {};"),
		c.Srcs("import('./other.js');"),
		c.Expected("ns.imprt_('./other.js');"))
}

// Java: RewriteDynamicImportsTest.internalImportPathError()
func (c *RewriteDynamicImportsTest) TestInternalImportPathError() {
	c.dynamicImportAlias = nil
	actualChunk0 := NewJSChunk("folder1/chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("../folder2/chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "import('./i0.js');"))

	c.TestError(c.Srcs([]*JSChunk{actualChunk0, actualChunk1}), ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH)
}

// Java: RewriteDynamicImportsTest.lowerLanguageWithoutAlias()
func (c *RewriteDynamicImportsTest) TestLowerLanguageWithoutAlias() {
	c.dynamicImportAlias = nil
	c.language = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.languageIn = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "import('./i0.js');"))

	c.TestWarning(c.Srcs([]*JSChunk{actualChunk0, actualChunk1}), RewriteDynamicImports_DYNAMIC_IMPORT_ALIASING_REQUIRED)
}

// Java: RewriteDynamicImportsTest.outputModulesExternInjected()
func (c *RewriteDynamicImportsTest) TestOutputModulesExternInjected() {
	c.chunkOutputType = CompilerOptions_ChunkOutputType_ES_MODULES
	c.dynamicImportAlias = nil
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "const nsPromise = import('./i0.js');"))

	c.TestExternChanges(
		c.Srcs(actualChunk0, actualChunk1),
		c.Expected("function "+ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN+"(importCallback) {}"))
}

// Java: RewriteDynamicImportsTest.outputModulesInternalImportSameChunk_unused()
func (c *RewriteDynamicImportsTest) TestOutputModulesInternalImportSameChunk_unused() {
	c.AllowExternsChanges()
	c.chunkOutputType = CompilerOptions_ChunkOutputType_ES_MODULES
	c.dynamicImportAlias = nil
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))
	actualChunk0.AddSourceFile(SourceFile_FromCode("i1.js", "import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode("i1.js", "Promise.resolve();")})

	c.Test(c.Srcs([]*JSChunk{actualChunk0}), expectedSrcs)
}

// Java: RewriteDynamicImportsTest.outputModulesInternalImportSameChunk_used()
func (c *RewriteDynamicImportsTest) TestOutputModulesInternalImportSameChunk_used() {
	c.AllowExternsChanges()
	c.chunkOutputType = CompilerOptions_ChunkOutputType_ES_MODULES
	c.dynamicImportAlias = nil
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))
	actualChunk0.AddSourceFile(SourceFile_FromCode("i1.js", "const ns = import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode("i1.js", "const ns = Promise.resolve(module$i0);")})

	c.Test(c.Srcs([]*JSChunk{actualChunk0}), expectedSrcs)
}

// Java: RewriteDynamicImportsTest.outputModulesInternalImportDifferentChunks_unused()
func (c *RewriteDynamicImportsTest) TestOutputModulesInternalImportDifferentChunks_unused() {
	c.AllowExternsChanges()
	c.chunkOutputType = CompilerOptions_ChunkOutputType_ES_MODULES
	c.dynamicImportAlias = nil
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode("i1.js", "import('./chunk0.js');")})

	c.Test(c.Srcs([]*JSChunk{actualChunk0, actualChunk1}), expectedSrcs)
}

// Java: RewriteDynamicImportsTest.outputModulesInternalImportDifferentChunks_used()
func (c *RewriteDynamicImportsTest) TestOutputModulesInternalImportDifferentChunks_used() {
	c.AllowExternsChanges()
	c.chunkOutputType = CompilerOptions_ChunkOutputType_ES_MODULES
	c.dynamicImportAlias = nil
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "const nsPromise = import('./i0.js');"))

	expectedSrcs :=
		newCompilerTestCase_Expected(
			[]*SourceFile{
				SourceFile_FromCode(
					"i0.js",
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`),
				SourceFile_FromCode(
					"i1.js",
					strings.ReplaceAll(`const nsPromise =
    import('./chunk0.js')
        .then(DYNAMIC_IMPORT_CALLBACK_FN(() => module$i0));
`, "DYNAMIC_IMPORT_CALLBACK_FN", ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN))})

	c.Test(c.Srcs([]*JSChunk{actualChunk0, actualChunk1}), expectedSrcs)
}

// Java: RewriteDynamicImportsTest.aliasExternInjectedSimple()
func (c *RewriteDynamicImportsTest) TestAliasExternInjectedSimple() {
	c.language = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.languageIn = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.dynamicImportAlias = jutil.StringPtr("import_")
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "const nsPromise = import('./i0.js');"))

	c.TestExternChanges(
		c.Externs(""), c.Srcs(actualChunk0, actualChunk1), c.Expected("function import_(specifier) {}"))
}

// Java: RewriteDynamicImportsTest.aliasExternInjectedSimpleImport()
func (c *RewriteDynamicImportsTest) TestAliasExternInjectedSimpleImport() {
	c.language = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.languageIn = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.dynamicImportAlias = jutil.StringPtr("import")

	// "import" is not a valid JS identifier name and will not parse.
	// Instead, we have to manually check the AST for the expected name.
	compiler := c.CreateCompiler()
	options := c.GetOptions()
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode(CompilerTestCase_GENERATED_EXTERNS_NAME, "")},
		[]*SourceFile{
			SourceFile_FromCode(
				CompilerTestCase_GENERATED_SRC_NAME, "var external = 'url'; const nsPromise = import(external);")},
		options)
	compiler.ParseInputs()
	NewGatherModuleMetadata(compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER).
		Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).
		Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	AssertThat(compiler.GetErrors()).IsEmpty()
	externsAndJs := compiler.GetRoot()
	externs := externsAndJs.GetFirstChild()
	root := externsAndJs.GetLastChild()
	compiler.BeforePass(c.GetName())
	c.GetProcessor(compiler).Process(externs, root)

	syntheticExterns := externs.GetFirstChild()
	AssertThat(syntheticExterns).IsNotNull()
	AssertThat(syntheticExterns.IsScript()).IsTrue()
	injectedAlias := syntheticExterns.GetFirstChild()
	AssertThat(injectedAlias).IsNotNull()
	AssertThat(injectedAlias.IsFunction()).IsTrue()
	AssertThat(injectedAlias.GetFirstChild().IsName()).IsTrue()
	AssertThat(injectedAlias.GetFirstChild().GetString()).IsEqualTo("import")
}

// Java: RewriteDynamicImportsTest.aliasExternInjectedSimpleAlreadyDefined()
func (c *RewriteDynamicImportsTest) TestAliasExternInjectedSimpleAlreadyDefined() {
	c.language = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.languageIn = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.dynamicImportAlias = jutil.StringPtr("import_")
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "const nsPromise = import('./i0.js');"))

	c.TestExternChanges(
		c.Externs("function import_() {}"),
		c.Srcs(actualChunk0, actualChunk1),
		c.Expected("function import_() {}"))
}

// Java: RewriteDynamicImportsTest.aliasExternInjectedQualified()
func (c *RewriteDynamicImportsTest) TestAliasExternInjectedQualified() {
	c.language = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.languageIn = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.dynamicImportAlias = jutil.StringPtr("foo.import")
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "const nsPromise = import('./i0.js');"))

	c.TestExternChanges(
		c.Externs(""),
		c.Srcs(actualChunk0, actualChunk1),
		c.Expected("var foo = {}; foo.import = function(specifier) {};"))
}

// Java: RewriteDynamicImportsTest.aliasExternInjectedQualifiedAlreadyDefinedComplete()
func (c *RewriteDynamicImportsTest) TestAliasExternInjectedQualifiedAlreadyDefinedComplete() {
	c.language = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.languageIn = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.dynamicImportAlias = jutil.StringPtr("foo.import")
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "const nsPromise = import('./i0.js');"))

	c.TestExternChanges(
		c.Externs("var foo = {}; foo.import = function(specifier) {};"),
		c.Srcs(actualChunk0, actualChunk1),
		c.Expected("var foo = {}; foo.import = function(specifier) {};"))
}

// Java: RewriteDynamicImportsTest.aliasExternInjectedQualifiedAlreadyDefinedPartial()
func (c *RewriteDynamicImportsTest) TestAliasExternInjectedQualifiedAlreadyDefinedPartial() {
	c.language = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.languageIn = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.dynamicImportAlias = jutil.StringPtr("foo.import")
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "const nsPromise = import('./i0.js');"))

	c.TestExternChanges(
		c.Externs("var foo = {};"), c.Srcs(actualChunk0, actualChunk1), c.Expected("var foo = {};"))
}

// Java: RewriteDynamicImportsTest.invalidAliasSimpleName()
func (c *RewriteDynamicImportsTest) TestInvalidAliasSimpleName() {
	c.dynamicImportAlias = jutil.StringPtr("fo'oo")
	c.language = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.languageIn = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "import('./i0.js');"))

	c.TestError(c.Srcs([]*JSChunk{actualChunk0, actualChunk1}), RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS)
}

// Java: RewriteDynamicImportsTest.invalidAliasQualifiedName()
func (c *RewriteDynamicImportsTest) TestInvalidAliasQualifiedName() {
	c.dynamicImportAlias = jutil.StringPtr("foo.bar['chunk']")
	c.language = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.languageIn = compilerOptionsLanguageModePtr(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	actualChunk0 := NewJSChunk("chunk0")
	actualChunk0.AddSourceFile(SourceFile_FromCode("i0.js", "const a = 1; export default a;"))

	actualChunk1 := NewJSChunk("chunk1")
	actualChunk1.AddSourceFile(SourceFile_FromCode("i1.js", "import('./i0.js');"))

	c.TestError(c.Srcs([]*JSChunk{actualChunk0, actualChunk1}), RewriteDynamicImports_DYNAMIC_IMPORT_INVALID_ALIAS)
}

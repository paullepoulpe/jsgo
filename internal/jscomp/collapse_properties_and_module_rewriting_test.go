// Port of com/google/javascript/jscomp/CollapsePropertiesAndModuleRewritingTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CollapsePropertiesAndModuleRewritingTest.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CollapsePropertiesAndModuleRewritingTest: tests for dynamic import
// handling that requires cooperation between ModuleRewriting and
// CollapseProperties.
type CollapsePropertiesAndModuleRewritingTest struct {
	CompilerTestCaseBase

	collapseLevel   CompilerOptions_PropertyCollapseLevel
	chunkOutputType CompilerOptions_ChunkOutputType
}

// Java: CollapsePropertiesAndModuleRewritingTest()
func NewCollapsePropertiesAndModuleRewritingTest() *CollapsePropertiesAndModuleRewritingTest {
	c := &CollapsePropertiesAndModuleRewritingTest{}
	c.InitCompilerTestCaseBase(c)
	c.collapseLevel = CompilerOptions_PropertyCollapseLevel_ALL
	c.chunkOutputType = CompilerOptions_ChunkOutputType_ES_MODULES
	return c
}

func TestCollapsePropertiesAndModuleRewritingTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCollapsePropertiesAndModuleRewritingTest() })
}

// SetUp runs CompilerTestCase's @Before method, then
// CollapsePropertiesAndModuleRewritingTest's (JUnit runs superclass @Before
// methods first).
func (c *CollapsePropertiesAndModuleRewritingTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.CustomSetUp()
}

// Java: CollapsePropertiesAndModuleRewritingTest.getProcessor(Compiler)
func (c *CollapsePropertiesAndModuleRewritingTest) GetProcessor(compiler *Compiler) CompilerPass {
	options := compiler.GetOptions()
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		factories := NewPassListBuilder(options)
		gatherModuleMetadata :=
			NewGatherModuleMetadata(compiler, false, options.GetModuleResolutionMode())
		factories.MaybeAdd(NewLoopablePassFactory(
			PassNames_GATHER_MODULE_METADATA,
			func(x AbstractCompiler) CompilerPass { return gatherModuleMetadata },
		))
		factories.MaybeAdd(NewLoopablePassFactory(
			PassNames_CREATE_MODULE_MAP,
			func(x AbstractCompiler) CompilerPass {
				return NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap())
			},
		))
		TranspilationPasses_AddEs6ModulePass(
			factories, NewPreprocessorSymbolTable_CachedInstanceFactory())
		factories.MaybeAdd(NewPassFactory(
			"REWRITE_DYNAMIC_IMPORT",
			func(x AbstractCompiler) CompilerPass {
				return NewRewriteDynamicImports(compiler, nil, CompilerOptions_ChunkOutputType_ES_MODULES)
			},
		))
		factories.MaybeAdd(NewPassFactory(
			PassNames_ES6_NORMALIZE_CLASSES,
			func(comp AbstractCompiler) CompilerPass { return NewEs6NormalizeClasses(comp) },
		))
		factories.MaybeAdd(NewLoopablePassFactory(
			PassNames_COLLAPSE_PROPERTIES,
			func(x AbstractCompiler) CompilerPass {
				return NewInlineAndCollapseProperties(compiler, InlineAndCollapsePropertiesOptions{
					PropertyCollapseLevel:    c.collapseLevel,
					ChunkOutputType:          c.chunkOutputType,
					HaveModulesBeenRewritten: true,
					ModuleResolutionMode:     options.GetModuleResolutionMode(),
				})
			},
		))

		for _, factory := range factories.Build() {
			factory.Create(compiler).Process(externs, root)
		}
	})
}

// Java: CollapsePropertiesAndModuleRewritingTest.customSetUp()
func (c *CollapsePropertiesAndModuleRewritingTest) CustomSetUp() {
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
}

// Java: CollapsePropertiesAndModuleRewritingTest.testModuleDynamicImport()
func (c *CollapsePropertiesAndModuleRewritingTest) TestModuleDynamicImport() {
	c.AllowExternsChanges()
	// TODO(bradfordcsmith): aggressive alias inlining prevents this example from working with
	// PropertyCollapseLevel.ALL
	c.collapseLevel = CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT

	inputChunks := []*JSChunk{NewJSChunk("entry"), NewJSChunk("mod1")}
	inputChunks[0].AddSourceFile(
		SourceFile_FromCode(
			"entry.js", "import('./mod1.js').then((ns) => console.log(ns.Foo.bar()));"))
	inputChunks[1].AddSourceFile(
		SourceFile_FromCode(
			"mod1.js",
			`export class Foo {
  static bar() { return 'bar'; }
}
`))
	inputChunks[1].AddDependency(inputChunks[0])

	expectedChunks := []*JSChunk{NewJSChunk("entry"), NewJSChunk("mod1")}
	expectedChunks[0].AddSourceFile(
		SourceFile_FromCode(
			"entry.js",
			`import('./mod1.js')
    .then(jscomp$DynamicImportCallback(() => module$mod1 ))
    .then(ns => console.log(ns.Foo.bar()));
`))
	expectedChunks[1].AddSourceFile(
		SourceFile_FromCode(
			"mod1.js",
			`class Foo$$module$mod1 {
  static bar() { return 'bar'; }
}
/** @const */ var module$mod1={};
/** @const */ module$mod1.Foo = Foo$$module$mod1;
`))
	expectedChunks[1].AddDependency(expectedChunks[0])

	c.Test(c.Srcs(inputChunks), c.Expected(expectedChunks))
}

// Java: CollapsePropertiesAndModuleRewritingTest.testModuleDynamicImportCommonJs()
func (c *CollapsePropertiesAndModuleRewritingTest) TestModuleDynamicImportCommonJs() {
	c.EnableProcessCommonJsModules()
	c.collapseLevel = CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT
	c.chunkOutputType = CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE
	c.SetModuleResolutionMode(ModuleLoader_ResolutionMode_WEBPACK)
	c.SetWebpackModulesById(
		jutil.NewImmutableMapBuilder[string, string]().
			Put("1", "entry.js").
			Put("2", "mod1.js").
			BuildOrThrow())

	inputs := []*SourceFile{}
	inputs = append(inputs,
		SourceFile_FromCode(
			"entry.js",
			`__webpack_require__.e(2).then(
    function() { return __webpack_require__(2);})
`))
	inputs = append(inputs, SourceFile_FromCode("mod1.js", "module.exports = 123;"))

	expected := []*SourceFile{}
	expected = append(expected,
		SourceFile_FromCode(
			"entry.js",
			"__webpack_require__.e(2).then(function() { return module$mod1.default;})"))
	expected = append(expected,
		SourceFile_FromCode(
			"mod1.js",
			"/** @const */ var module$mod1={}; /** @const */ module$mod1.default = 123;"))

	c.Test(c.Srcs(inputs), c.Expected(expected))
}

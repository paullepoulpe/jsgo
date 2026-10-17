// Port of com/google/javascript/jscomp/ConvertChunksToESModulesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ConvertChunksToESModulesTest.java)

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

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ConvertChunksToESModulesTest holds unit tests for ConvertChunksToESModules
type ConvertChunksToESModulesTest struct {
	CompilerTestCaseBase
}

func NewConvertChunksToESModulesTest() *ConvertChunksToESModulesTest {
	c := &ConvertChunksToESModulesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestConvertChunksToESModulesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewConvertChunksToESModulesTest() })
}

// Java: ConvertChunksToESModulesTest.getProcessor(Compiler)
func (c *ConvertChunksToESModulesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewConvertChunksToESModules(compiler)
}

// Java: ConvertChunksToESModulesTest.testVarDeclarations_acrossChunks()
func (c *ConvertChunksToESModulesTest) TestVarDeclarations_acrossChunks() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = 1;").
				AddChunk("a").
				Build()),
		c.Expected(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = 1; export {a}").
				AddChunk("import {a} from './m0.js'; a").
				Build()))

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = 1, b = 2, c = 3;").
				AddChunk("a;c;").
				Build()),
		c.Expected(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = 1, b = 2, c = 3; export {a, c};").
				AddChunk("import {a, c} from './m0.js'; a; c;").
				Build()))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = 1, b = 2, c = 3;").
				AddChunk("b;c;").
				Build()),
		c.Expected(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = 1, b = 2, c = 3; export {b,c};").
				AddChunk("import {b, c} from './m0.js'; b;c;").
				Build()))
}

// Java: ConvertChunksToESModulesTest.testMultipleInputsPerChunk()
func (c *ConvertChunksToESModulesTest) TestMultipleInputsPerChunk() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	original := JSChunkGraphBuilder_ForStar().
		AddChunk("var a = 1;").
		AddChunk("a").
		Build()

	original[0].AddSourceFile(SourceFile_FromCode("m0-1", "console.log(a)"))

	expected := JSChunkGraphBuilder_ForStar().
		AddChunk("var a = 1; console.log(a); export {a}").
		AddChunk("import {a} from './m0.js'; a").
		Build()

	expected[0].AddSourceFile(SourceFile_FromCode("m0-1", ""))

	c.Test(c.Srcs(original), c.Expected(expected))
}

// Java: ConvertChunksToESModulesTest.testImportPathReferenceAbsolute()
func (c *ConvertChunksToESModulesTest) TestImportPathReferenceAbsolute() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	original := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1;", "/js/m0").
		AddChunkWithName("a", "/js/m1").
		Build()

	expected := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1; export {a}", "/js/m0").
		AddChunkWithName("import {a} from './m0.js'; a", "/js/m1").
		Build()

	c.Test(c.Srcs(original), c.Expected(expected))
}

// Java: ConvertChunksToESModulesTest.testImportPathReferenceAbsoluteWithRelative1()
func (c *ConvertChunksToESModulesTest) TestImportPathReferenceAbsoluteWithRelative1() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	original := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1;", "other/m0").
		AddChunkWithName("a", "/js/m1").
		Build()

	c.TestError(
		c.Srcs(original),
		c.Error(ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH).
			WithMessage(
				"Unable to compute relative import path from \"/js/m1.js\" to \"other/m0.js\""))
}

// Java: ConvertChunksToESModulesTest.testImportPathReferenceAbsoluteWithRelative2()
func (c *ConvertChunksToESModulesTest) TestImportPathReferenceAbsoluteWithRelative2() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	original := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1;", "/other/m0").
		AddChunkWithName("a", "js/m1").
		Build()

	c.TestError(
		c.Srcs(original),
		c.Error(ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH).
			WithMessage(
				"Unable to compute relative import path from \"js/m1.js\" to \"/other/m0.js\""))
}

// Java: ConvertChunksToESModulesTest.testImportPathAmbiguous()
func (c *ConvertChunksToESModulesTest) TestImportPathAmbiguous() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	original := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1;", "js/m0").
		AddChunkWithName("a", "js/m1").
		Build()

	expected := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1; export {a}", "js/m0").
		AddChunkWithName("import {a} from './m0.js'; a", "js/m1").
		Build()

	c.Test(c.Srcs(original), c.Expected(expected))
}

// Java: ConvertChunksToESModulesTest.testImportPathMixedDepth1()
func (c *ConvertChunksToESModulesTest) TestImportPathMixedDepth1() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	original := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1;", "js/m0").
		AddChunkWithName("a", "m1").
		Build()

	c.TestError(
		c.Srcs(original),
		c.Error(ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH).
			WithMessage("Unable to compute relative import path from \"m1.js\" to \"js/m0.js\""))
}

// Java: ConvertChunksToESModulesTest.testImportPathMixedDepth2()
func (c *ConvertChunksToESModulesTest) TestImportPathMixedDepth2() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	original := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1;", "m0").
		AddChunkWithName("a", "js/m1").
		Build()

	expected := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1; export {a}", "m0").
		AddChunkWithName("import {a} from '../m0.js'; a", "js/m1").
		Build()

	c.Test(c.Srcs(original), c.Expected(expected))
}

// Java: ConvertChunksToESModulesTest.testImportPathMixedDepth3()
func (c *ConvertChunksToESModulesTest) TestImportPathMixedDepth3() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	original := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1;", "js/other/path/one/m0").
		AddChunkWithName("a", "external/path/m1").
		Build()

	expected := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1; export {a}", "js/other/path/one/m0").
		AddChunkWithName(
			"import {a} from '../../js/other/path/one/m0.js'; a", "external/path/m1").
		Build()

	c.Test(c.Srcs(original), c.Expected(expected))
}

// Java: ConvertChunksToESModulesTest.testImportPathParentAboveRoot()
func (c *ConvertChunksToESModulesTest) TestImportPathParentAboveRoot() {
	original := JSChunkGraphBuilder_ForStar().
		AddChunkWithName("var a = 1;", "js/m0").
		AddChunkWithName("a", "../node_modules/m1").
		Build()

	c.TestError(
		c.Srcs(original),
		c.Error(ConvertChunksToESModules_UNABLE_TO_COMPUTE_RELATIVE_PATH).
			WithMessage(
				"Unable to compute relative import path from \"../node_modules/m1.js\" to"+
					" \"js/m0.js\""))
}

// Java: ConvertChunksToESModulesTest.testForcedESModuleSemantics()
func (c *ConvertChunksToESModulesTest) TestForcedESModuleSemantics() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = 1;").
				AddChunk("var b = 1;").
				Build()),
		c.Expected(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = 1; export {};").
				AddChunk("import './m0.js'; var b = 1;").
				Build()))
}

// Java: ConvertChunksToESModulesTest.testAssignToImport()
func (c *ConvertChunksToESModulesTest) TestAssignToImport() {
	c.TestError(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = 1;").
				AddChunk("a = 2;").
				Build()),
		c.Error(ConvertChunksToESModules_ASSIGNMENT_TO_IMPORT).
			WithMessage(
				"Imported symbol \"a\" in chunk \"m1.js\" cannot be assigned (defined in"+
					" \"m0.js\")"))
}

// Java: ConvertChunksToESModulesTest.testChunkDependenciesCreateImportStatementsForSideEffects()
func (c *ConvertChunksToESModulesTest) TestChunkDependenciesCreateImportStatementsForSideEffects() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("window.a = true").
				AddChunk("console.log(window.a) // should be true").
				Build()),
		c.Expected(
			JSChunkGraphBuilder_ForStar().
				AddChunk("window.a = true; export {};").
				AddChunk("import './m0.js'; console.log(window.a);").
				Build()))
}

// Java: ConvertChunksToESModulesTest.testDynamicImportRewriting1()
func (c *ConvertChunksToESModulesTest) TestDynamicImportRewriting1() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					jutil.StringReplace(`import('./m1.js')
    .then(DYNAMIC_IMPORT_CALLBACK_FN(function () { return module$i1; }))
    .then(function (ns) { console.log(ns.default); });
`, "DYNAMIC_IMPORT_CALLBACK_FN", ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN)).
				AddChunk(
					`const a$$module$i1 = 1;
var $jscompDefaultExport$$module$i1 = a$$module$i1;
/** @const */ var module$i1 = {};
/** @const */ module$i1.default = $jscompDefaultExport$$module$i1;
`).
				Build()),
		c.Expected(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					`import('./m1.js')
    .then(function ($) { return $.module$i1; })
    .then(function (ns) { console.log(ns.default); });
export {};
`).
				AddChunk(
					`import './m0.js';
const a$$module$i1 = 1;
var $jscompDefaultExport$$module$i1 = a$$module$i1;
/** @const */ var module$i1 = {};
/** @const */ module$i1.default = $jscompDefaultExport$$module$i1;
export {module$i1};
`).
				Build()))
}

// Java: ConvertChunksToESModulesTest.testDynamicImportRewriting2()
func (c *ConvertChunksToESModulesTest) TestDynamicImportRewriting2() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`).
				AddChunk(
					jutil.StringReplace(`import('./m0.js')
    .then(DYNAMIC_IMPORT_CALLBACK_FN(() => module$i0))
    .then((ns) => console.log(ns.default));
`, "DYNAMIC_IMPORT_CALLBACK_FN", ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN)).
				Build()),
		c.Expected(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					`const a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
export {module$i0};
`).
				AddChunk(
					`import './m0.js';
import('./m0.js')
    .then(($) => $.module$i0)
    .then((ns) => console.log(ns.default));
`).
				Build()))
}

// Java: ConvertChunksToESModulesTest.testDynamicImportWithoutCallback()
func (c *ConvertChunksToESModulesTest) TestDynamicImportWithoutCallback() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					`var a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`).
				AddChunk(
					`import('./m0.js')
    .then((ns) => console.log(ns.default));
`).
				Build()),
		c.Expected(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					`var a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
export {};
`).
				AddChunk(
					`import './m0.js';
import('./m0.js')
    .then((ns) => console.log(ns.default));
`).
				Build()))
}

// Java: ConvertChunksToESModulesTest.testDynamicImportRewritingErrors1()
func (c *ConvertChunksToESModulesTest) TestDynamicImportRewritingErrors1() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.TestError(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					`var a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`).
				AddChunk(
					jutil.StringReplace(`import('./m0.js')
    .then(DYNAMIC_IMPORT_CALLBACK_FN(() => {}))
    .then((ns) => console.log(ns.default));
`, "DYNAMIC_IMPORT_CALLBACK_FN", ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN)).
				Build()),
		c.Error(ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK).
			WithMessage(
				"Dynamic import callback encountered wih an invalid format. "+
					"Unable to find valid namespace reference."))
}

// Java: ConvertChunksToESModulesTest.testDynamicImportRewritingErrors2()
func (c *ConvertChunksToESModulesTest) TestDynamicImportRewritingErrors2() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.TestError(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					`var a$$module$i0 = 1;
var $jscompDefaultExport$$module$i0 = a$$module$i0;
/** @const */ var module$i0 = {};
/** @const */ module$i0.default = $jscompDefaultExport$$module$i0;
`).
				AddChunk(
					jutil.StringReplace(`import('./m0.js')
    .then(DYNAMIC_IMPORT_CALLBACK_FN(true))
    .then((ns) => console.log(ns.default));
`, "DYNAMIC_IMPORT_CALLBACK_FN", ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN)).
				Build()),
		c.Error(ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK).
			WithMessage(
				"Dynamic import callback encountered wih an invalid format. "+
					"Unable to find valid callback function."))
}

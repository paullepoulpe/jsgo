// Port of com/google/javascript/jscomp/RewriteGoogJsImportsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/RewriteGoogJsImportsTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Tests for RewriteGoogJsImports that involve rewriting. CheckGoogJsImportTest
// has the link tests.
type RewriteGoogJsImportsTest struct {
	CompilerTestCaseBase
}

// JsFileRegexParser determines if this file is base.js by looking at the
// first comment of the file.
var rewriteGoogJsImportsTest_BASE = SourceFile_FromCode(
	"/closure/base.js",
	`/**
 * @fileoverview
 * @provideGoog
 */
`)

var rewriteGoogJsImportsTest_GOOG = SourceFile_FromCode(
	"/closure/goog.js",
	`export const require = goog.require;
export function foo() {}
export class MyClass {}
export const constant = 0;
`)

func NewRewriteGoogJsImportsTest() *RewriteGoogJsImportsTest {
	c := &RewriteGoogJsImportsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteGoogJsImportsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteGoogJsImportsTest() })
}

// Java: RewriteGoogJsImportsTest.getProcessor(Compiler)
func (c *RewriteGoogJsImportsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		gmm :=
			NewGatherModuleMetadata(
				compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER)
		gmm.Process(externs, root)
		mmc := NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap())
		mmc.Process(externs, root)
		NewRewriteGoogJsImports(compiler, RewriteGoogJsImports_Mode_LINT_AND_REWRITE, compiler.GetModuleMap()).
			Process(externs, root)
	})
}

// Java: RewriteGoogJsImportsTest.testBaseAndGoogUntouched()
func (c *RewriteGoogJsImportsTest) TestBaseAndGoogUntouched() {
	c.TestSame(c.Srcs(rewriteGoogJsImportsTest_BASE, rewriteGoogJsImportsTest_GOOG))
}

// Java: RewriteGoogJsImportsTest.testIfCannotDetectGoogJsThenGlobalizesAll()
func (c *RewriteGoogJsImportsTest) TestIfCannotDetectGoogJsThenGlobalizesAll() {
	testcode := SourceFile_FromCode(
		"testcode", "import * as goog from './closure/goog.js'; use(goog.bad);")

	expected := SourceFile_FromCode("testcode", "import './closure/goog.js'; use(goog.bad);")

	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)

	// No base.js = no detecting goog.js
	c.Test(c.Srcs(rewriteGoogJsImportsTest_GOOG, testcode), c.Expected(rewriteGoogJsImportsTest_GOOG, expected))

	// No goog.js
	c.Test(c.Srcs(rewriteGoogJsImportsTest_BASE, testcode), c.Expected(rewriteGoogJsImportsTest_BASE, expected))
	c.Test(c.Srcs(testcode), c.Expected(expected))

	// Linting still happens.
	c.TestError("import * as notgoog from './goog.js';", RewriteGoogJsImports_GOOG_JS_IMPORT_MUST_BE_GOOG_STAR)
}

// Java: RewriteGoogJsImportsTest.testImportStar()
func (c *RewriteGoogJsImportsTest) TestImportStar() {
	c.Test(
		c.Srcs(
			rewriteGoogJsImportsTest_BASE,
			rewriteGoogJsImportsTest_GOOG,
			SourceFile_FromCode(
				"testcode",
				`import * as goog from './closure/goog.js';
use(goog.require, goog.foo, goog.MyClass, goog.constant);
`)),
		c.Expected(
			rewriteGoogJsImportsTest_BASE,
			rewriteGoogJsImportsTest_GOOG,
			SourceFile_FromCode(
				"testcode",
				`import './closure/goog.js';
use(goog.require, goog.foo, goog.MyClass, goog.constant);
`)))
	AssertThat(
		append([]string{}, c.GetLastCompiler().
			GetModuleMap().
			GetModuleByPath(c.GetLastCompiler().GetModuleLoader().Resolve("testcode")).
			BoundNames().
			KeySlice()...)).
		IsEmpty()
}

// Java: RewriteGoogJsImportsTest.testGoogAndBaseInExterns()
func (c *RewriteGoogJsImportsTest) TestGoogAndBaseInExterns() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)

	c.Test(
		c.Externs(rewriteGoogJsImportsTest_BASE, rewriteGoogJsImportsTest_GOOG),
		c.Srcs(
			SourceFile_FromCode(
				"testcode",
				`import * as goog from './closure/goog.js';
use(goog.require, goog.foo, goog.MyClass, goog.constant);
`)),
		c.Expected(
			SourceFile_FromCode(
				"testcode",
				`import './closure/goog.js';
use(goog.require, goog.foo, goog.MyClass, goog.constant);
`)))
}

// Java: RewriteGoogJsImportsTest.testKnownBadPropertyAccess()
func (c *RewriteGoogJsImportsTest) TestKnownBadPropertyAccess() {
	c.Test(
		c.Srcs(
			rewriteGoogJsImportsTest_BASE,
			rewriteGoogJsImportsTest_GOOG,
			SourceFile_FromCode(
				"testcode",
				`import * as goog from './closure/goog.js';
use(goog.require, goog.bad);
`)),
		c.Expected(
			rewriteGoogJsImportsTest_BASE,
			rewriteGoogJsImportsTest_GOOG,
			SourceFile_FromCode(
				"testcode",
				`import * as $goog from './closure/goog.js';
use(goog.require, $goog.bad);
`)))
	AssertThat(
		append([]string{}, c.GetLastCompiler().
			GetModuleMap().
			GetModuleByPath(c.GetLastCompiler().GetModuleLoader().Resolve("testcode")).
			BoundNames().
			KeySlice()...)).
		ContainsExactly("$goog")
}

// Java: RewriteGoogJsImportsTest.testReexportGoog()
func (c *RewriteGoogJsImportsTest) TestReexportGoog() {
	c.TestError(
		c.Srcs(
			rewriteGoogJsImportsTest_BASE,
			rewriteGoogJsImportsTest_GOOG,
			SourceFile_FromCode(
				"testcode",
				`import * as goog from './closure/goog.js';
export {goog};
`)),
		c.Error(RewriteGoogJsImports_GOOG_JS_REEXPORTED))

	c.TestError(
		c.Srcs(
			rewriteGoogJsImportsTest_BASE,
			rewriteGoogJsImportsTest_GOOG,
			SourceFile_FromCode(
				"testcode",
				`import * as goog from './closure/goog.js';
export default goog;
`)),
		c.Error(RewriteGoogJsImports_GOOG_JS_REEXPORTED))

	c.TestError(
		c.Srcs(rewriteGoogJsImportsTest_BASE, rewriteGoogJsImportsTest_GOOG, SourceFile_FromCode("testcode", "export * from './closure/goog.js';")),
		c.Error(RewriteGoogJsImports_GOOG_JS_REEXPORTED))

	c.TestError(
		c.Srcs(
			rewriteGoogJsImportsTest_BASE,
			rewriteGoogJsImportsTest_GOOG,
			SourceFile_FromCode("testcode", "export {require} from './closure/goog.js';")),
		c.Error(RewriteGoogJsImports_GOOG_JS_REEXPORTED))
}

// Java: RewriteGoogJsImportsTest.testImportMeta()
func (c *RewriteGoogJsImportsTest) TestImportMeta() {
	c.Test(
		c.Srcs(
			rewriteGoogJsImportsTest_BASE,
			rewriteGoogJsImportsTest_GOOG,
			SourceFile_FromCode(
				"testcode",
				`import * as goog from './closure/goog.js';
use(import.meta);
`)),
		c.Expected(
			rewriteGoogJsImportsTest_BASE,
			rewriteGoogJsImportsTest_GOOG,
			SourceFile_FromCode(
				"testcode",
				`import './closure/goog.js';
use(import.meta);
`)))
}

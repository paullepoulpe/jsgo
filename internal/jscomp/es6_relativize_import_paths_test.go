// Port of com/google/javascript/jscomp/Es6RelativizeImportPathsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/Es6RelativizeImportPathsTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"testing"
)

type Es6RelativizeImportPathsTest struct {
	CompilerTestCaseBase

	prefixReplacements   *jutil.LinkedHashMap[string, string]
	moduleResolutionMode ModuleLoader_ResolutionMode
	moduleRoots          []string
}

func NewEs6RelativizeImportPathsTest() *Es6RelativizeImportPathsTest {
	c := &Es6RelativizeImportPathsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestEs6RelativizeImportPathsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RelativizeImportPathsTest() })
}

// Java: Es6RelativizeImportPathsTest.setUp()
func (c *Es6RelativizeImportPathsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.prefixReplacements = esModulesTest_mapOf()
	c.moduleResolutionMode = ModuleLoader_ResolutionMode_BROWSER
	c.moduleRoots = []string{}
}

// Java: Es6RelativizeImportPathsTest.getOptions()
func (c *Es6RelativizeImportPathsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetBrowserResolverPrefixReplacements(c.prefixReplacements)
	options.SetModuleResolutionMode(c.moduleResolutionMode)
	options.SetModuleRoots(c.moduleRoots)
	options.SetPathEscaper(ModuleLoader_PathEscaper_CANONICALIZE_ONLY)
	return options
}

// Java: Es6RelativizeImportPathsTest.getProcessor(Compiler)
func (c *Es6RelativizeImportPathsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs6RelativizeImportPaths(compiler)
}

// Java: Es6RelativizeImportPathsTest.testLocalRelativePathStayTheSame()
func (c *Es6RelativizeImportPathsTest) TestLocalRelativePathStayTheSame() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.TestSame(c.Srcs(SourceFile_FromCode("/file.js", "import './foo.js';")))

	c.TestSame(c.Srcs(SourceFile_FromCode("/nested/file.js", "import './foo.js';")))

	c.TestSame(c.Srcs(SourceFile_FromCode("/file.js", "import './nested/foo.js';")))

	c.TestSame(c.Srcs(SourceFile_FromCode("http://google.com/file.js", "import './foo.js';")))
}

// Java: Es6RelativizeImportPathsTest.testRelativeParentPathStayTheSame()
func (c *Es6RelativizeImportPathsTest) TestRelativeParentPathStayTheSame() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.TestSame(c.Srcs(SourceFile_FromCode("/nested/file.js", "import '../foo.js';")))

	c.TestSame(c.Srcs(SourceFile_FromCode("/nested/file.js", "import '../nested/foo.js';")))

	c.TestSame(c.Srcs(SourceFile_FromCode("/p0/p1/file.js", "import '../p2/foo.js';")))

	c.TestSame(c.Srcs(SourceFile_FromCode("/p0/p1/file.js", "import '../../foo.js';")))
	c.TestSame(c.Srcs(SourceFile_FromCode("http://google.com/nested/file.js", "import '../foo.js';")))
}

// Java: Es6RelativizeImportPathsTest.testAbsolutePathsStayTheSame()
func (c *Es6RelativizeImportPathsTest) TestAbsolutePathsStayTheSame() {
	// Absolute means that with a scheme, not those starting with "/"
	c.IgnoreWarnings(ModuleLoader_INVALID_MODULE_PATH)
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"http://google.com/file.js", "import 'http://google.com/other/file.js';")))

	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"http://golang.org/file.js", "import 'http://google.com/other/file.js';")))

	c.TestSame(c.Srcs(SourceFile_FromCode("/file.js", "import 'http://google.com/other/file.js';")))
}

// Java: Es6RelativizeImportPathsTest.testRootedPathsBecomeRelative()
func (c *Es6RelativizeImportPathsTest) TestRootedPathsBecomeRelative() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		c.Srcs(SourceFile_FromCode("/file.js", "import '/foo.js';")),
		c.Expected(SourceFile_FromCode("/file.js", "import './foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("/nested/file.js", "import '/foo.js';")),
		c.Expected(SourceFile_FromCode("/nested/file.js", "import '../foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("/nested/file.js", "import '/path/foo.js';")),
		c.Expected(SourceFile_FromCode("/nested/file.js", "import '../path/foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("/file.js", "import '/nested/foo.js';")),
		c.Expected(SourceFile_FromCode("/file.js", "import './nested/foo.js';")))
}

// Java: Es6RelativizeImportPathsTest.testNonStandardResolutionBecomesRelative()
func (c *Es6RelativizeImportPathsTest) TestNonStandardResolutionBecomesRelative() {
	c.moduleResolutionMode = ModuleLoader_ResolutionMode_BROWSER_WITH_TRANSFORMED_PREFIXES
	c.prefixReplacements = esModulesTest_mapOf("raw0/", "/mapped0/", "raw1/", "/mapped1/")

	c.Test(
		c.Srcs(SourceFile_FromCode("/file.js", "import 'raw0/foo.js';")),
		c.Expected(SourceFile_FromCode("/file.js", "import './mapped0/foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("/nested/file.js", "import 'raw0/foo.js';")),
		c.Expected(SourceFile_FromCode("/nested/file.js", "import '../mapped0/foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("/file.js", "import 'raw0/bar/foo.js';")),
		c.Expected(SourceFile_FromCode("/file.js", "import './mapped0/bar/foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("/file.js", "import 'raw1/foo.js';")),
		c.Expected(SourceFile_FromCode("/file.js", "import './mapped1/foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("/mapped1/file.js", "import 'raw1/foo.js';")),
		c.Expected(SourceFile_FromCode("/mapped1/file.js", "import './foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("http://google.com/file.js", "import 'raw1/foo.js';")),
		c.Expected(SourceFile_FromCode("http://google.com/file.js", "import './mapped1/foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("http://google.com/nested/file.js", "import 'raw1/foo.js';")),
		c.Expected(SourceFile_FromCode("http://google.com/file.js", "import '../mapped1/foo.js';")))
}

// Java: Es6RelativizeImportPathsTest.testModuleRootsBecomeRelative()
func (c *Es6RelativizeImportPathsTest) TestModuleRootsBecomeRelative() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.moduleRoots = []string{"root1", "root2"}

	c.Test(
		c.Srcs(SourceFile_FromCode("/file.js", "import '/root1/foo.js';")),
		c.Expected(SourceFile_FromCode("/file.js", "import './foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("/nested/file.js", "import '/root2/foo.js';")),
		c.Expected(SourceFile_FromCode("/nested/file.js", "import '../foo.js';")))

	c.Test(
		c.Srcs(SourceFile_FromCode("/file.js", "import '/root1/bar/foo.js';")),
		c.Expected(SourceFile_FromCode("/file.js", "import './bar/foo.js';")))
}

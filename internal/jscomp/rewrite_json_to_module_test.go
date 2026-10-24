// Port of com/google/javascript/jscomp/RewriteJsonToModuleTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RewriteJsonToModuleTest.java)

/*
 * Copyright 2016 The Closure Compiler Authors.
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

// RewriteJsonToModuleTest holds unit tests for RewriteJsonToModule
type RewriteJsonToModuleTest struct {
	CompilerTestCaseBase
}

func NewRewriteJsonToModuleTest() *RewriteJsonToModuleTest {
	c := &RewriteJsonToModuleTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteJsonToModuleTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteJsonToModuleTest() })
}

// rewriteJsonToModuleTest_noopPass is the anonymous CompilerPass of
// getProcessor.
type rewriteJsonToModuleTest_noopPass struct {
	// Go-only: language-adapter — distinct identity
	_ byte
}

// Java: RewriteJsonToModuleTest.getProcessor(Compiler).new CompilerPass() {...}.process(Node, Node)
func (p *rewriteJsonToModuleTest_noopPass) Process(externs rhino.Node, root rhino.Node) {
	// No-op, RewriteJsonToModule handling is done directly after parsing.
}

// Java: RewriteJsonToModuleTest.getProcessor(Compiler)
func (c *RewriteJsonToModuleTest) GetProcessor(compiler *Compiler) CompilerPass {
	return &rewriteJsonToModuleTest_noopPass{}
}

// Java: RewriteJsonToModuleTest.getOptions()
func (c *RewriteJsonToModuleTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	// Trigger module processing after parsing.
	options.SetProcessCommonJSModules(true)
	options.SetModuleResolutionMode(ModuleLoader_ResolutionMode_NODE)
	options.SetPackageJsonEntryNames([]string{"browser", "main"})
	return options
}

// rewriteJsonToModuleTest_mapOf is ImmutableMap.of(k1, v1, ...), the
// expected entries of a MapSubject.containsExactly(k1, v1, ...).
//
// Go-only: test-helper — Truth's MapSubject.containsExactly(Object...)
// becomes ContainsExactlyEntriesIn(map).
func rewriteJsonToModuleTest_mapOf(kvs ...string) *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	for i := 0; i < len(kvs); i += 2 {
		m.Put(kvs[i], kvs[i+1])
	}
	return m
}

// Java: RewriteJsonToModuleTest.testJsonFile()
func (c *RewriteJsonToModuleTest) TestJsonFile() {
	c.Test(
		c.Srcs(SourceFile_FromCode("/test.json", "{ 'foo': 'bar'}")),
		c.Expected(
			`/** @fileoverview @suppress {undefinedVars} */
goog.provide('module$test_json'); var module$test_json = { 'foo': 'bar'};
`))

	AssertThat(c.GetLastCompiler().GetModuleLoader().GetPackageJsonMainEntries()).IsEmpty()
}

// Java: RewriteJsonToModuleTest.testPackageJsonFile()
func (c *RewriteJsonToModuleTest) TestPackageJsonFile() {
	c.Test(
		c.Srcs(SourceFile_FromCode("/package.json", "{ 'main': 'foo/bar/baz.js'}")),
		c.Expected(
			`/** @fileoverview @suppress {undefinedVars} */
goog.provide('module$package_json')
var module$package_json = {'main': 'foo/bar/baz.js'};
`))

	AssertThat(c.GetLastCompiler().GetModuleLoader().GetPackageJsonMainEntries()).HasSize(1)
	AssertThat(c.GetLastCompiler().GetModuleLoader().GetPackageJsonMainEntries()).
		ContainsKey("/package.json")
	AssertThat(c.GetLastCompiler().GetModuleLoader().GetPackageJsonMainEntries()).
		ContainsEntry("/package.json", "/foo/bar/baz.js")
}

// Java: RewriteJsonToModuleTest.testPackageJsonWithoutMain()
func (c *RewriteJsonToModuleTest) TestPackageJsonWithoutMain() {
	c.Test(
		c.Srcs(SourceFile_FromCode("/package.json", "{'other': { 'main': 'foo/bar/baz.js'}}")),
		c.Expected(
			`/** @fileoverview @suppress {undefinedVars} */
goog.provide('module$package_json')
var module$package_json = {'other': { 'main': 'foo/bar/baz.js'}};
`))

	AssertThat(c.GetLastCompiler().GetModuleLoader().GetPackageJsonMainEntries()).IsEmpty()
}

// Java: RewriteJsonToModuleTest.testPackageJsonFileBrowserField()
func (c *RewriteJsonToModuleTest) TestPackageJsonFileBrowserField() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"/package.json", "{ 'main': 'foo/bar/baz.js', 'browser': 'browser/foo.js' }")),
		c.Expected(
			`/** @fileoverview @suppress {undefinedVars} */
goog.provide('module$package_json')
var module$package_json = {
  'main': 'foo/bar/baz.js',
  'browser': 'browser/foo.js'
};
`))

	AssertThat(c.GetLastCompiler().GetModuleLoader().GetPackageJsonMainEntries()).
		ContainsExactlyEntriesIn(rewriteJsonToModuleTest_mapOf("/package.json", "/browser/foo.js"))
}

// Java: RewriteJsonToModuleTest.testPackageJsonFileBrowserFieldAdvancedUsage()
func (c *RewriteJsonToModuleTest) TestPackageJsonFileBrowserFieldAdvancedUsage() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"/package.json",
				`{ 'main': 'foo/bar/baz.js',
  'browser': { 'dont/include.js': false,
               'foo/bar/baz.js': 'replaced/main.js',
               'override/relative.js': './with/this.js',
               'override/explicitly.js': 'with/other.js'} }
`)),
		c.Expected(
			`/** @fileoverview @suppress {undefinedVars} */
goog.provide('module$package_json')
var module$package_json = {
  'main': 'foo/bar/baz.js',
  'browser': {
    'dont/include.js': false,
    'foo/bar/baz.js': 'replaced/main.js',
    'override/relative.js': './with/this.js',
    'override/explicitly.js': 'with/other.js'
  }
};
`))

	packageJsonMainEntries :=
		c.GetLastCompiler().GetModuleLoader().GetPackageJsonMainEntries()
	AssertThat(packageJsonMainEntries).HasSize(5)
	AssertThat(packageJsonMainEntries).ContainsEntry("/package.json", "/foo/bar/baz.js")
	AssertThat(packageJsonMainEntries).ContainsEntry("/foo/bar/baz.js", "/replaced/main.js")
	// NodeModuleResolver knows how to normalize this entry's value
	AssertThat(packageJsonMainEntries).ContainsEntry("/override/relative.js", "/./with/this.js")
	AssertThat(packageJsonMainEntries).
		ContainsEntry("/dont/include.js", ModuleLoader_JSC_BROWSER_SKIPLISTED_MARKER)
	AssertThat(packageJsonMainEntries).ContainsEntry("/override/explicitly.js", "/with/other.js")
}

// Java: RewriteJsonToModuleTest.testPackageJsonBrowserFieldAdvancedUsageGH2625()
func (c *RewriteJsonToModuleTest) TestPackageJsonBrowserFieldAdvancedUsageGH2625() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"/package.json",
				`{ 'main': 'foo/bar/baz.js',
  'browser': { './a/b.js': './c/d.js',
               './server.js': 'client.js'} }
`)),
		c.Expected(
			`/** @fileoverview @suppress {undefinedVars} */
goog.provide('module$package_json')
var module$package_json = {
  'main': 'foo/bar/baz.js',
  'browser': {
    './a/b.js': './c/d.js',
    './server.js': 'client.js'
  }
};
`))

	packageJsonMainEntries :=
		c.GetLastCompiler().GetModuleLoader().GetPackageJsonMainEntries()
	AssertThat(packageJsonMainEntries).
		ContainsExactlyEntriesIn(rewriteJsonToModuleTest_mapOf(
			"/package.json", "/foo/bar/baz.js",

			// Test that we have normalized the key, value is normalized by NodeModuleResolver
			"/a/b.js", "/./c/d.js",
			"/server.js", "/client.js"))
}

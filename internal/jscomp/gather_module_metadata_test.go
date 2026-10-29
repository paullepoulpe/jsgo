// Port of com/google/javascript/jscomp/GatherModuleMetadataTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/GatherModuleMetadataTest.java)

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
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type GatherModuleMetadataTest struct {
	CompilerTestCaseBase

	rewriteScriptsToModules bool
	sortOnly                bool
	entryPoints             []*ModuleIdentifier
}

func NewGatherModuleMetadataTest() *GatherModuleMetadataTest {
	c := &GatherModuleMetadataTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestGatherModuleMetadataTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewGatherModuleMetadataTest() })
}

// gatherModuleMetadataTest_list is a multiset as the Iterable that Truth's
// assertThat(Iterable) gets.
func gatherModuleMetadataTest_list(m *StringCounts) []string {
	return m.ToSlice()
}

// gatherModuleMetadataTest_keys is a map's keySet() as the Iterable that
// Truth's assertThat(Iterable) gets.
//
// Go-only: test-helper — an empty key set is an empty collection, not null.
func gatherModuleMetadataTest_keys(m *OrderedStringMap[*ModuleMetadataMap_ModuleMetadata]) []string {
	return slices.AppendSeq([]string{}, m.Keys())
}

// Java: GatherModuleMetadataTest.setUp()
func (c *GatherModuleMetadataTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.entryPoints = []*ModuleIdentifier{}
	c.rewriteScriptsToModules = false
	c.sortOnly = false
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
}

// Java: GatherModuleMetadataTest.getOptions()
func (c *GatherModuleMetadataTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	if len(c.entryPoints) != 0 {
		jutil.CheckStateMsg(!c.sortOnly, "sortOnly must be false if entry points are provided.")
		options.SetDependencyOptions(DependencyOptions_PruneForEntryPoints(c.entryPoints))
	} else if c.sortOnly {
		options.SetDependencyOptions(DependencyOptions_SortOnly())
	}
	return options
}

// Java: GatherModuleMetadataTest.getProcessor(Compiler)
func (c *GatherModuleMetadataTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		if c.rewriteScriptsToModules {
			NewEs6RewriteScriptsToModules(compiler).Process(externs, root)
		}
		NewGatherModuleMetadata(compiler /* processCommonJsModules= */, true, ModuleLoader_ResolutionMode_BROWSER).
			Process(externs, root)
	})
}

// Java: GatherModuleMetadataTest.metadataMap()
func (c *GatherModuleMetadataTest) MetadataMap() *ModuleMetadataMap {
	return c.GetLastCompiler().GetModuleMetadataMap()
}

// Java: GatherModuleMetadataTest.testGoogProvide()
func (c *GatherModuleMetadataTest) TestGoogProvide() {
	c.TestSame("goog.provide('my.provide');")
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("my.provide")
	AssertThat(c.MetadataMap().GetModulesByPath().ContainsKey("testcode")).IsTrue()
	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.provide")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.provide")
	AssertThat(m.IsGoogProvide()).IsTrue()
}

// Java: GatherModuleMetadataTest.testGoogProvideWithGoogDeclaredInOtherFile()
func (c *GatherModuleMetadataTest) TestGoogProvideWithGoogDeclaredInOtherFile() {
	// Closure's base.js declare the global goog. It should be ignored when scanning the provide'd
	// file. Only local variables named goog should cause the pass to back off.
	c.TestSame(c.Srcs("var goog;", "goog.provide('my.provide');"))
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("my.provide")
	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.provide")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.provide")
	AssertThat(m.IsGoogProvide()).IsTrue()
}

// Java: GatherModuleMetadataTest.testLocalGoogIsIgnored()
func (c *GatherModuleMetadataTest) TestLocalGoogIsIgnored() {
	c.TestSame("function bar(goog) { goog.provide('my.provide'); }")
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).IsEmpty()
	AssertThat(c.MetadataMap().GetModulesByPath().Get("testcode").UsesClosure()).IsFalse()
}

// Java: GatherModuleMetadataTest.testMultipleGoogProvide()
func (c *GatherModuleMetadataTest) TestMultipleGoogProvide() {
	c.TestSame("goog.provide('my.first.provide'); goog.provide('my.second.provide');")
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).
		ContainsExactly("my.first.provide", "my.second.provide")
	AssertThat(c.MetadataMap().GetModulesByPath().ContainsKey("testcode")).IsTrue()
	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.first.provide")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.first.provide", "my.second.provide")
	AssertThat(m.IsGoogProvide()).IsTrue()

	m = c.MetadataMap().GetModulesByGoogNamespace().Get("my.second.provide")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.first.provide", "my.second.provide")
	AssertThat(m.IsGoogProvide()).IsTrue()

	m = c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.first.provide", "my.second.provide")
	AssertThat(m.IsGoogProvide()).IsTrue()
}

// Java: GatherModuleMetadataTest.testProvideNamespaceValidation()
func (c *GatherModuleMetadataTest) TestProvideNamespaceValidation() {
	c.Test(c.Srcs("goog.provide('');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))
	c.Test(c.Srcs("goog.provide(' ');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))
	c.Test(c.Srcs("goog.provide('a..b');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))

	c.Test(c.Srcs("goog.provide('ā');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))
	c.TestSame(c.Srcs("goog.provide('a');"))

	c.TestSame(c.Srcs("goog.provide('a.class');"))
	c.Test(c.Srcs("goog.provide('class.a');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.TestError(c.Srcs("goog.provide('a.class');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))
	c.TestError(c.Srcs("goog.provide('class.a');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))
}

// Java: GatherModuleMetadataTest.testGoogModule()
func (c *GatherModuleMetadataTest) TestGoogModule() {
	c.TestSame("goog.module('my.module');")
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("my.module")
	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.module")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.module")
	AssertThat(m.IsGoogProvide()).IsFalse()
	AssertThat(m.IsGoogModule()).IsTrue()
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.IsLegacyGoogModule()).IsFalse()
}

// Java: GatherModuleMetadataTest.testModuleIdValidation()
func (c *GatherModuleMetadataTest) TestModuleIdValidation() {
	c.Test(c.Srcs("goog.module();"), c.Error(ClosureRewriteModule_INVALID_MODULE_ID_ARG))
	c.Test(c.Srcs("goog.module('');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))
	c.Test(c.Srcs("goog.module(' ');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))
	c.Test(c.Srcs("goog.module('a..b');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))
	c.Test(c.Srcs("goog.module('a. .b');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))
	c.Test(c.Srcs("goog.module('a.-.b');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))

	c.Test(c.Srcs("goog.module('0');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))
	c.Test(c.Srcs("goog.module('ā');"), c.Error(GatherModuleMetadata_INVALID_NAMESPACE_OR_MODULE_ID))

	c.TestSame(c.Srcs("goog.module('a');"))
	c.TestSame(c.Srcs("goog.module('a0');"))
	c.TestSame(c.Srcs("goog.module('$');"))
}

// Java: GatherModuleMetadataTest.testGoogModuleWithDefaultExport()
func (c *GatherModuleMetadataTest) TestGoogModuleWithDefaultExport() {
	// exports = 0; on its own is CommonJS!
	c.TestSame("goog.module('my.module'); exports = 0;")
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("my.module")
	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.module")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.module")
	AssertThat(m.IsGoogProvide()).IsFalse()
	AssertThat(m.IsGoogModule()).IsTrue()
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.IsLegacyGoogModule()).IsFalse()
}

// Java: GatherModuleMetadataTest.testLegacyGoogModule()
func (c *GatherModuleMetadataTest) TestLegacyGoogModule() {
	c.TestSame("goog.module('my.module'); goog.module.declareLegacyNamespace();")
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("my.module")
	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.module")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.module")
	AssertThat(m.IsGoogProvide()).IsFalse()
	AssertThat(m.IsGoogModule()).IsTrue()
	AssertThat(m.IsNonLegacyGoogModule()).IsFalse()
	AssertThat(m.IsLegacyGoogModule()).IsTrue()
}

// Java: GatherModuleMetadataTest.testLoadModule()
func (c *GatherModuleMetadataTest) TestLoadModule() {
	c.TestSame(
		`goog.loadModule(function(exports) {
  goog.module('my.module');
  return exports;
});
`)

	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("my.module")

	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.module")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.module")
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.Path()).IsNull()

	m = c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).IsEmpty()
	AssertThat(m.IsNonProvideScript()).IsTrue()
}

// Java: GatherModuleMetadataTest.testLoadModuleLegacyNamespace()
func (c *GatherModuleMetadataTest) TestLoadModuleLegacyNamespace() {
	c.TestSame(
		`goog.loadModule(function(exports) {
  goog.module('my.module');
  goog.module.declareLegacyNamespace();
  return exports;
});
`)

	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("my.module")

	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.module")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.module")
	AssertThat(m.IsLegacyGoogModule()).IsTrue()
	AssertThat(m.Path()).IsNull()

	m = c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).IsEmpty()
	AssertThat(m.IsNonProvideScript()).IsTrue()
}

// Java: GatherModuleMetadataTest.testLoadModuleUseStrict()
func (c *GatherModuleMetadataTest) TestLoadModuleUseStrict() {
	c.TestSame(
		`goog.loadModule(function(exports) {
  'use strict';
  goog.module('with.strict');
  return exports;
});
`)

	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("with.strict")

	m := c.MetadataMap().GetModulesByGoogNamespace().Get("with.strict")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("with.strict")
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.Path()).IsNull()

	m = c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).IsEmpty()
	AssertThat(m.IsNonProvideScript()).IsTrue()
}

// Java: GatherModuleMetadataTest.testMultipleGoogModuleCallsInLoadModule()
func (c *GatherModuleMetadataTest) TestMultipleGoogModuleCallsInLoadModule() {
	c.TestSame(
		`goog.loadModule(function(exports) {
  goog.module('multiple.calls.c0');
  goog.module('multiple.calls.c1');
  return exports;
});
`)

	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).
		ContainsExactly("multiple.calls.c0", "multiple.calls.c1")

	m := c.MetadataMap().GetModulesByGoogNamespace().Get("multiple.calls.c0")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("multiple.calls.c0", "multiple.calls.c1")
	AssertThat(c.MetadataMap().GetModulesByGoogNamespace().Get("multiple.calls.c1")).
		IsSameInstanceAs(m)
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.Path()).IsNull()

	m = c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).IsEmpty()
	AssertThat(m.IsNonProvideScript()).IsTrue()
}

// Java: GatherModuleMetadataTest.testMultipleGoogLoadModules()
func (c *GatherModuleMetadataTest) TestMultipleGoogLoadModules() {
	c.TestSame(
		`goog.loadModule(function(exports) {
  goog.module('multiple.calls.c0');
  return exports;
});

goog.loadModule(function(exports) {
  goog.module('multiple.calls.c1');
  return exports;
});
`)

	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).
		ContainsExactly("multiple.calls.c0", "multiple.calls.c1")

	m := c.MetadataMap().GetModulesByGoogNamespace().Get("multiple.calls.c0")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("multiple.calls.c0")
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.Path()).IsNull()

	m = c.MetadataMap().GetModulesByGoogNamespace().Get("multiple.calls.c1")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("multiple.calls.c1")
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.Path()).IsNull()

	m = c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).IsEmpty()
	AssertThat(m.IsNonProvideScript()).IsTrue()
}

// Java: GatherModuleMetadataTest.testBundleGoogLoadModuleAndProvides()
func (c *GatherModuleMetadataTest) TestBundleGoogLoadModuleAndProvides() {
	c.TestSame(
		`goog.provide('some.provide');

goog.provide('some.other.provide');

goog.loadModule(function(exports) {
  goog.module('multiple.calls.c0');
  return exports;
});

goog.loadModule(function(exports) {
  goog.module('multiple.calls.c1');
  return exports;
});
`)

	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).
		ContainsExactly(
			"some.provide", "some.other.provide", "multiple.calls.c0", "multiple.calls.c1")

	m := c.MetadataMap().GetModulesByGoogNamespace().Get("multiple.calls.c0")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("multiple.calls.c0")
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.Path()).IsNull()
	AssertThat(m.UsesClosure()).IsTrue()

	m = c.MetadataMap().GetModulesByGoogNamespace().Get("multiple.calls.c1")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("multiple.calls.c1")
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.Path()).IsNull()
	AssertThat(m.UsesClosure()).IsTrue()

	m = c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("some.provide", "some.other.provide")
	AssertThat(m.IsGoogProvide()).IsTrue()
	AssertThat(m.UsesClosure()).IsTrue()
}

// Java: GatherModuleMetadataTest.testBundleGoogLoadModuleAndProvidesWithGoogDefined()
func (c *GatherModuleMetadataTest) TestBundleGoogLoadModuleAndProvidesWithGoogDefined() {
	c.TestSame(
		`/** @provideGoog */
var goog = {};

goog.provide('some.provide');

goog.provide('some.other.provide');

goog.loadModule(function(exports) {
  goog.module('multiple.calls.c0');
  return exports;
});

goog.loadModule(function(exports) {
  goog.module('multiple.calls.c1');
  return exports;
});
`)

	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).
		ContainsExactly(
			"some.provide", "some.other.provide", "multiple.calls.c0", "multiple.calls.c1")

	m := c.MetadataMap().GetModulesByGoogNamespace().Get("multiple.calls.c0")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("multiple.calls.c0")
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.Path()).IsNull()
	AssertThat(m.UsesClosure()).IsFalse()

	m = c.MetadataMap().GetModulesByGoogNamespace().Get("multiple.calls.c1")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("multiple.calls.c1")
	AssertThat(m.IsNonLegacyGoogModule()).IsTrue()
	AssertThat(m.Path()).IsNull()
	AssertThat(m.UsesClosure()).IsFalse()

	m = c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("some.provide", "some.other.provide")
	AssertThat(m.IsGoogProvide()).IsTrue()
	AssertThat(m.UsesClosure()).IsFalse()
}

// Java: GatherModuleMetadataTest.testEs6Module()
func (c *GatherModuleMetadataTest) TestEs6Module() {
	c.TestSame("export var x;")
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).IsEmpty()
	AssertThat(c.MetadataMap().GetModulesByPath().ContainsKey("testcode")).IsTrue()
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).IsEmpty()
	AssertThat(m.IsEs6Module()).IsTrue()
}

// Java: GatherModuleMetadataTest.testEs6ModuleDeclareModuleId()
func (c *GatherModuleMetadataTest) TestEs6ModuleDeclareModuleId() {
	c.TestSame("export var x; goog.declareModuleId('my.module');")
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("my.module")
	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.module")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.module")
	AssertThat(m.IsEs6Module()).IsTrue()
	AssertThat(m.IsGoogModule()).IsFalse()
}

// Java: GatherModuleMetadataTest.testEs6ModuleDeclareModuleIdImportedGoog()
func (c *GatherModuleMetadataTest) TestEs6ModuleDeclareModuleIdImportedGoog() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode("goog.js", ""),
			SourceFile_FromCode(
				"testcode",
				`import * as goog from './goog.js';
export var x;
goog.declareModuleId('my.module');
`)))
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("my.module")
	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.module")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.module")
	AssertThat(m.IsEs6Module()).IsTrue()
	AssertThat(m.IsGoogModule()).IsFalse()
}

// Java: GatherModuleMetadataTest.testCommonJsModule()
func (c *GatherModuleMetadataTest) TestCommonJsModule() {
	c.TestSame("exports = 0;")
	AssertThat(c.MetadataMap().GetModulesByGoogNamespace().IsEmpty()).IsTrue()
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(m.IsCommonJs()).IsTrue()
}

// Java: GatherModuleMetadataTest.testDuplicateProvides()
func (c *GatherModuleMetadataTest) TestDuplicateProvides() {
	c.TestError(
		c.Srcs("goog.provide('duplciated');", "goog.provide('duplciated');"),
		ClosurePrimitiveErrors_DUPLICATE_NAMESPACE)
}

// Java: GatherModuleMetadataTest.testDuplicateProvidesInSameFile()
func (c *GatherModuleMetadataTest) TestDuplicateProvidesInSameFile() {
	c.TestError(
		"goog.provide('duplciated');\ngoog.provide('duplciated');",
		ClosurePrimitiveErrors_DUPLICATE_NAMESPACE)
}

// Java: GatherModuleMetadataTest.testDuplicateProvideAndGoogModule()
func (c *GatherModuleMetadataTest) TestDuplicateProvideAndGoogModule() {
	c.TestError(
		c.Srcs("goog.provide('duplciated');", "goog.module('duplciated');"),
		ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE)
	c.TestError(
		c.Srcs("goog.module('duplciated');", "goog.provide('duplciated');"),
		ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE)
}

// Java: GatherModuleMetadataTest.testDuplicateProvideAndEs6Module()
func (c *GatherModuleMetadataTest) TestDuplicateProvideAndEs6Module() {
	c.TestError(
		c.Srcs("goog.provide('duplciated');", "export {}; goog.declareModuleId('duplciated');"),
		ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE)
	c.TestError(
		c.Srcs("export {}; goog.declareModuleId('duplciated');", "goog.provide('duplciated');"),
		ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE)
}

// Java: GatherModuleMetadataTest.testDuplicateGoogModules()
func (c *GatherModuleMetadataTest) TestDuplicateGoogModules() {
	c.TestError(
		c.Srcs("goog.module('duplciated');", "goog.module('duplciated');"),
		ClosurePrimitiveErrors_DUPLICATE_MODULE)
}

// Java: GatherModuleMetadataTest.testDuplicateGoogAndEs6Module()
func (c *GatherModuleMetadataTest) TestDuplicateGoogAndEs6Module() {
	c.TestError(
		c.Srcs("goog.module('duplciated');", "export {}; goog.declareModuleId('duplciated');"),
		ClosurePrimitiveErrors_DUPLICATE_MODULE)
	c.TestError(
		c.Srcs("export {}; goog.declareModuleId('duplciated');", "goog.module('duplciated');"),
		ClosurePrimitiveErrors_DUPLICATE_MODULE)
}

// Java: GatherModuleMetadataTest.testDuplicatEs6Modules()
func (c *GatherModuleMetadataTest) TestDuplicatEs6Modules() {
	c.TestError(
		c.Srcs(
			"export {}; goog.declareModuleId('duplciated');",
			"export {}; goog.declareModuleId('duplciated');"),
		ClosurePrimitiveErrors_DUPLICATE_MODULE)
	c.TestError(
		c.Srcs(
			"export {}; goog.declareModuleId('duplciated');",
			"export {}; goog.declareModuleId('duplciated');"),
		ClosurePrimitiveErrors_DUPLICATE_MODULE)
}

// Java: GatherModuleMetadataTest.testDuplicateModuleWarningsIncludeFileName()
func (c *GatherModuleMetadataTest) TestDuplicateModuleWarningsIncludeFileName() {
	// duplicates in the same file
	c.Test(
		c.Srcs("goog.provide('duplicated'); goog.provide('duplicated')"),
		c.Error(ClosurePrimitiveErrors_DUPLICATE_NAMESPACE).WithMessageContaining("testcode"))
	// duplicate of provide in earlier file
	c.Test(
		c.Srcs("goog.provide('duplicated');", "goog.module('duplicated')"),
		c.Error(ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE).
			WithMessageContaining("testcode0"))
	// duplicate of module in earlier file
	c.Test(
		c.Srcs("goog.module('duplicated');", "goog.module('duplicated')"),
		c.Error(ClosurePrimitiveErrors_DUPLICATE_MODULE).WithMessageContaining("testcode0"))
}

// Java: GatherModuleMetadataTest.testUsesGlobalClosure()
func (c *GatherModuleMetadataTest) TestUsesGlobalClosure() {
	c.TestSame("goog.isArray(foo);")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(m.UsesClosure()).IsTrue()
}

// Java: GatherModuleMetadataTest.testUsesGlobalClosureNoFunctionCall()
func (c *GatherModuleMetadataTest) TestUsesGlobalClosureNoFunctionCall() {
	c.TestSame("var b = goog.isArray;")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(m.UsesClosure()).IsTrue()
}

// Java: GatherModuleMetadataTest.testLocalGoogIsNotClosure()
func (c *GatherModuleMetadataTest) TestLocalGoogIsNotClosure() {
	c.TestSame("function bar() { var goog; goog.isArray(foo); }")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(m.UsesClosure()).IsFalse()
}

// Java: GatherModuleMetadataTest.testGoogInSameScriptGoogIsNotClosure()
func (c *GatherModuleMetadataTest) TestGoogInSameScriptGoogIsNotClosure() {
	c.TestSame("/** @provideGoog */ var goog = {}; goog.isArray(foo);")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(m.UsesClosure()).IsFalse()
}

// Java: GatherModuleMetadataTest.testGoogInOtherScriptGoogIsClosure()
func (c *GatherModuleMetadataTest) TestGoogInOtherScriptGoogIsClosure() {
	c.TestSame(c.Srcs("/** @provideGoog */ var goog = {};", "goog.isArray(foo);"))
	m := c.MetadataMap().GetModulesByPath().Get("testcode1")
	AssertThat(m.UsesClosure()).IsTrue()
}

// Java: GatherModuleMetadataTest.testImportedGoogIsClosure()
func (c *GatherModuleMetadataTest) TestImportedGoogIsClosure() {
	c.TestSame("import * as goog from '/goog.js'; goog.isArray(foo);")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(m.UsesClosure()).IsTrue()
}

// Java: GatherModuleMetadataTest.testRequireType()
func (c *GatherModuleMetadataTest) TestRequireType() {
	c.TestSame("goog.requireType('my.Type');")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.WeaklyRequiredGoogNamespaces())).ContainsExactly("my.Type")
}

// Java: GatherModuleMetadataTest.testRequireDynamic()
func (c *GatherModuleMetadataTest) TestRequireDynamic() {
	c.TestSame("async function test() {await goog.requireDynamic('my.Type');}")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.DynamicallyRequiredGoogNamespaces())).ContainsExactly("my.Type")
}

// Java: GatherModuleMetadataTest.testRequireDynamicWithIllegalArg()
func (c *GatherModuleMetadataTest) TestRequireDynamicWithIllegalArg() {
	c.TestError(
		"async function test() {await goog.requireDynamic('my.Type','extra.Type');}",
		GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC)
	c.TestError(
		"async function test() {await goog.requireDynamic(42);}",
		GatherModuleMetadata_INVALID_REQUIRE_DYNAMIC)
}

// Java: GatherModuleMetadataTest.testMaybeRequire()
func (c *GatherModuleMetadataTest) TestMaybeRequire() {
	c.TestSame("goog.maybeRequireFrameworkInternalOnlyDoNotCallOrElse('my.Type')")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.MaybeRequiredGoogNamespaces())).ContainsExactly("my.Type")
}

// Java: GatherModuleMetadataTest.testMaybeRequireWithIllegalArg()
func (c *GatherModuleMetadataTest) TestMaybeRequireWithIllegalArg() {
	c.TestError(
		"goog.maybeRequireFrameworkInternalOnlyDoNotCallOrElse('my.Type','extra.Type');",
		GatherModuleMetadata_INVALID_MAYBE_REQUIRE)
	c.TestError(
		"goog.maybeRequireFrameworkInternalOnlyDoNotCallOrElse(42);",
		GatherModuleMetadata_INVALID_MAYBE_REQUIRE)
}

// Java: GatherModuleMetadataTest.testRequiredClosureNamespaces()
func (c *GatherModuleMetadataTest) TestRequiredClosureNamespaces() {
	c.TestSame("goog.require('my.Type');")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.StronglyRequiredGoogNamespaces())).ContainsExactly("my.Type")
}

// Java: GatherModuleMetadataTest.testImport()
func (c *GatherModuleMetadataTest) TestImport() {
	c.TestSame("import '@spec!';")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.Es6ImportSpecifiers())).ContainsExactly("@spec!")
}

// Java: GatherModuleMetadataTest.testExport()
func (c *GatherModuleMetadataTest) TestExport() {
	c.TestSame("export { name } from '@spec!';")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.Es6ImportSpecifiers())).ContainsExactly("@spec!")
}

// Java: GatherModuleMetadataTest.testImportOrder()
func (c *GatherModuleMetadataTest) TestImportOrder() {
	c.TestSame("import 'first'; export { name } from 'second'; import 'third';")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.Es6ImportSpecifiers())).ContainsExactly("first", "second", "third")
}

// Java: GatherModuleMetadataTest.testSetTestOnly()
func (c *GatherModuleMetadataTest) TestSetTestOnly() {
	c.TestSame("goog.setTestOnly();")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(m.IsTestOnly()).IsTrue()
}

// Java: GatherModuleMetadataTest.testSetTestOnlyWithStringArg()
func (c *GatherModuleMetadataTest) TestSetTestOnlyWithStringArg() {
	c.TestSame("goog.setTestOnly('string');")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(m.IsTestOnly()).IsTrue()
}

// Java: GatherModuleMetadataTest.testSetTestOnlyWithExtraArg()
func (c *GatherModuleMetadataTest) TestSetTestOnlyWithExtraArg() {
	c.TestError("goog.setTestOnly('string', 'string');", GatherModuleMetadata_INVALID_SET_TEST_ONLY)
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(m.IsTestOnly()).IsFalse()
}

// Java: GatherModuleMetadataTest.testSetTestOnlyWithInvalidArg()
func (c *GatherModuleMetadataTest) TestSetTestOnlyWithInvalidArg() {
	c.TestError("goog.setTestOnly(0);", GatherModuleMetadata_INVALID_SET_TEST_ONLY)
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(m.IsTestOnly()).IsFalse()
}

// Java: GatherModuleMetadataTest.testGatherFromExterns()
func (c *GatherModuleMetadataTest) TestGatherFromExterns() {
	// js_lib will put data in externs for .i.js files.
	c.Test(c.Externs("export var x; goog.declareModuleId('my.module');"), c.Srcs(""))
	AssertThat(gatherModuleMetadataTest_keys(c.MetadataMap().GetModulesByGoogNamespace())).ContainsExactly("my.module")
	m := c.MetadataMap().GetModulesByGoogNamespace().Get("my.module")
	AssertThat(gatherModuleMetadataTest_list(m.GoogNamespaces())).ContainsExactly("my.module")
	AssertThat(m.IsEs6Module()).IsTrue()
	AssertThat(m.IsGoogModule()).IsFalse()
}

// Java: GatherModuleMetadataTest.testImportedScript()
func (c *GatherModuleMetadataTest) TestImportedScript() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode("imported.js", "console.log('lol');"),
			SourceFile_FromCode("notimported.js", "console.log('lol');"),
			SourceFile_FromCode("module.js", "import './imported.js';")))

	AssertThat(c.MetadataMap().GetModulesByPath().Get("imported.js").ModuleType()).
		IsEqualTo(ModuleMetadataMap_ModuleType_SCRIPT)
	AssertThat(c.MetadataMap().GetModulesByPath().Get("notimported.js").ModuleType()).
		IsEqualTo(ModuleMetadataMap_ModuleType_SCRIPT)
}

// Java: GatherModuleMetadataTest.testImportedScriptWithScriptsToModules()
func (c *GatherModuleMetadataTest) TestImportedScriptWithScriptsToModules() {
	// Default dependency options should still mark imported files as ES modules.
	c.rewriteScriptsToModules = true

	c.Test(
		c.Srcs(
			SourceFile_FromCode("imported.js", "console.log('lol');"),
			SourceFile_FromCode("notimported.js", "console.log('lol');"),
			SourceFile_FromCode("module.js", "import './imported.js';")))

	AssertThat(c.MetadataMap().GetModulesByPath().Get("imported.js").ModuleType()).
		IsEqualTo(ModuleMetadataMap_ModuleType_ES6_MODULE)
	AssertThat(c.MetadataMap().GetModulesByPath().Get("notimported.js").ModuleType()).
		IsEqualTo(ModuleMetadataMap_ModuleType_SCRIPT)
}

// Java: GatherModuleMetadataTest.testImportedScriptWithEntryPoint()
func (c *GatherModuleMetadataTest) TestImportedScriptWithEntryPoint() {
	c.rewriteScriptsToModules = true
	c.entryPoints = []*ModuleIdentifier{ModuleIdentifier_ForFile("module.js")}

	c.Test(
		c.Srcs(
			SourceFile_FromCode("imported.js", "console.log('lol');"),
			SourceFile_FromCode("notimported.js", "console.log('lol');"),
			SourceFile_FromCode("module.js", "import './imported.js';")))

	AssertThat(c.MetadataMap().GetModulesByPath().Get("imported.js").ModuleType()).
		IsEqualTo(ModuleMetadataMap_ModuleType_ES6_MODULE)
	// Pruned
	AssertThat(c.MetadataMap().GetModulesByPath().ContainsKey("notimported.js")).IsFalse()
}

// Java: GatherModuleMetadataTest.testImportedScriptWithSortOnly()
func (c *GatherModuleMetadataTest) TestImportedScriptWithSortOnly() {
	c.rewriteScriptsToModules = true
	c.sortOnly = true

	c.Test(
		c.Srcs(
			SourceFile_FromCode("imported.js", "console.log('lol');"),
			SourceFile_FromCode("notimported.js", "console.log('lol');"),
			SourceFile_FromCode("module.js", "import './imported.js';")))

	AssertThat(c.MetadataMap().GetModulesByPath().Get("imported.js").ModuleType()).
		IsEqualTo(ModuleMetadataMap_ModuleType_ES6_MODULE)
	AssertThat(c.MetadataMap().GetModulesByPath().Get("notimported.js").ModuleType()).
		IsEqualTo(ModuleMetadataMap_ModuleType_SCRIPT)
}

// Java: GatherModuleMetadataTest.testDynamicImport()
func (c *GatherModuleMetadataTest) TestDynamicImport() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode("imported.js", "export default function() {};"),
			SourceFile_FromCode("script.js", "import('./imported.js')")))

	m := c.MetadataMap().GetModulesByPath().Get("script.js")
	AssertThat(m.IsNonProvideScript()).IsTrue()
	AssertThat(gatherModuleMetadataTest_list(m.Es6ImportSpecifiers())).ContainsExactly("./imported.js")
}

// Java: GatherModuleMetadataTest.testReadToggle()
func (c *GatherModuleMetadataTest) TestReadToggle() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo$2etoggles.js",
				`goog.module('foo$2etoggles');
exports.TOGGLE_foo = goog.readToggleInternalDoNotCallDirectly('foo');
exports.TOGGLE_bar = goog.readToggleInternalDoNotCallDirectly('bar');
exports.TOGGLE_b_a_z = goog.readToggleInternalDoNotCallDirectly('b_a_z');
`)))
	m := c.MetadataMap().GetModulesByPath().Get("foo$2etoggles.js")
	AssertThat(gatherModuleMetadataTest_list(m.ReadToggles())).IsEmpty()
}

// Java: GatherModuleMetadataTest.testToggleModuleImportDestructured()
func (c *GatherModuleMetadataTest) TestToggleModuleImportDestructured() {
	c.TestSame("const {TOGGLE_foo, TOGGLE_b_a_z} = goog.require('foo$2etoggles');")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.ReadToggles())).ContainsExactly("foo", "b_a_z")
}

// Java: GatherModuleMetadataTest.testToggleModuleImportDestructuredAliased()
func (c *GatherModuleMetadataTest) TestToggleModuleImportDestructuredAliased() {
	c.TestSame("const {TOGGLE_bar: baz} = goog.require('foo$2etoggles');")
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.ReadToggles())).ContainsExactly("bar")
}

// Java: GatherModuleMetadataTest.testToggleModuleImportAsModule()
func (c *GatherModuleMetadataTest) TestToggleModuleImportAsModule() {
	c.TestSame(
		`const toggles = goog.require('foo$2etoggles');
function foo() {
  console.log(toggles.TOGGLE_foo);
  console.log(toggles.TOGGLE_b_a_z);
}
`)
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.ReadToggles())).ContainsExactly("foo", "b_a_z")
}

// Java: GatherModuleMetadataTest.testToggleModuleGoogModuleGet()
func (c *GatherModuleMetadataTest) TestToggleModuleGoogModuleGet() {
	c.TestSame(
		`goog.provide('foo.soy.gencode');
goog.require('foo$2etoggles');
function foo() {
  console.log(goog.module.get('foo$2etoggles').TOGGLE_foo);
  console.log(goog.module.get('bar').TOGGLE_bar);
}
`)
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.ReadToggles())).ContainsExactly("foo")
}

// Java: GatherModuleMetadataTest.testToggleModulegoogModuleGetInvalidUsage()
func (c *GatherModuleMetadataTest) TestToggleModulegoogModuleGetInvalidUsage() {
	c.Test(
		c.Srcs(
			`goog.provide('foo.soy.gencode');
goog.require('foo$2etoggles');
const toggles = goog.module.get('foo$2etoggles');
`),
		c.Error(GatherModuleMetadata_INVALID_TOGGLE_USAGE))
}

// Java: GatherModuleMetadataTest.testToggleModuleImportAsSideEffect()
func (c *GatherModuleMetadataTest) TestToggleModuleImportAsSideEffect() {
	c.Test(
		c.Srcs("goog.require('foo$2etoggles');"), //
		c.Error(GatherModuleMetadata_INVALID_TOGGLE_USAGE))
}

// Java: GatherModuleMetadataTest.testToggleModuleImportDestructuredWithInvalidName()
func (c *GatherModuleMetadataTest) TestToggleModuleImportDestructuredWithInvalidName() {
	c.Test(
		c.Srcs("const {foo} = goog.require('foo$2etoggles');"),
		c.Error(GatherModuleMetadata_INVALID_TOGGLE_USAGE))
}

// Java: GatherModuleMetadataTest.testToggleModuleImportAsModuleWithInvalidPropertyName()
func (c *GatherModuleMetadataTest) TestToggleModuleImportAsModuleWithInvalidPropertyName() {
	c.Test(
		c.Srcs(
			`const foo = goog.require('foo$2etoggles');
console.log(foo.bar);
`),
		c.Error(GatherModuleMetadata_INVALID_TOGGLE_USAGE))
}

// Java: GatherModuleMetadataTest.testToggleModuleInvalidModuleUsage()
func (c *GatherModuleMetadataTest) TestToggleModuleInvalidModuleUsage() {
	c.Test(
		c.Srcs(
			`const foo = goog.require('foo$2etoggles');
console.log(foo);
`),
		c.Error(GatherModuleMetadata_INVALID_TOGGLE_USAGE))
}

// Java: GatherModuleMetadataTest.testToggleModuleIgnoreShadowedUsage()
func (c *GatherModuleMetadataTest) TestToggleModuleIgnoreShadowedUsage() {
	c.TestSame(
		`const toggles = goog.require('foo$2etoggles');
function foo() {
  const toggles = {};
  console.log(toggles.foo)
  console.log(toggles.TOGGLE_bar)
}
`)
	m := c.MetadataMap().GetModulesByPath().Get("testcode")
	AssertThat(gatherModuleMetadataTest_list(m.ReadToggles())).IsEmpty()
}

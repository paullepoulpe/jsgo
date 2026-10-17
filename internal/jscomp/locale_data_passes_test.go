// Port of com/google/javascript/jscomp/LocaleDataPassesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/LocaleDataPassesTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// LocaleDataPassesTest is the test which checks that replacer works
// correctly.
type LocaleDataPassesTest struct {
	CompilerTestCaseBase

	// Messages returned from fake bundle, keyed by `JsMessage.id`.
	testMode LocaleDataPassesTest_TestMode
	locale   string
}

// LocaleDataPassesTest_TestMode indicates which part of the replacement
// we're currently testing
type LocaleDataPassesTest_TestMode int32

const (
	LocaleDataPassesTest_TestMode_PROTECT_DATA LocaleDataPassesTest_TestMode = iota
	// Test replacement of the protected function call form with the final message values.
	LocaleDataPassesTest_TestMode_REPLACE_PROTECTED_DATA
)

// Java: LocaleDataPassesTest()
func NewLocaleDataPassesTest() *LocaleDataPassesTest {
	c := &LocaleDataPassesTest{testMode: LocaleDataPassesTest_TestMode_PROTECT_DATA, locale: "en"}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestLocaleDataPassesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewLocaleDataPassesTest() })
}

// Java: LocaleDataPassesTest.setUp()
func (c *LocaleDataPassesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// locale prework and replacement always occurs after goog.provide and module rewriting
	// so enable this rewriting.
	c.EnableRewriteClosureCode()
	c.EnableRewriteClosureProvides()
	c.EnableClosurePassForExpected()
	c.EnableNormalize()

	c.testMode = LocaleDataPassesTest_TestMode_PROTECT_DATA
	c.locale = "en"
}

// Java: LocaleDataPassesTest.getOptions()
func (c *LocaleDataPassesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	locale := c.locale
	options.SetLocale(&locale)
	return options
}

// localeDataPassesTest_protectPass is the anonymous PROTECT_DATA CompilerPass.
type localeDataPassesTest_protectPass struct {
	compiler *Compiler
}

func (p *localeDataPassesTest_protectPass) Process(externs *rhino.Node, root *rhino.Node) {
	extract := NewLocaleDataPasses_ProtectGoogLocale(p.compiler)
	extract.Process(externs, root)
}

// localeDataPassesTest_replacePass is the anonymous REPLACE_PROTECTED_DATA
// CompilerPass.
type localeDataPassesTest_replacePass struct {
	compiler *Compiler
}

func (p *localeDataPassesTest_replacePass) Process(externs *rhino.Node, root *rhino.Node) {
	extract := NewLocaleDataPasses_ProtectGoogLocale(p.compiler)
	extract.Process(externs, root)
	subs :=
		NewLocaleDataPasses_LocaleSubstitutions(
			p.compiler, p.compiler.GetOptions().GetLocale())
	subs.Process(externs, root)
}

// Java: LocaleDataPassesTest.getProcessor(Compiler)
func (c *LocaleDataPassesTest) GetProcessor(compiler *Compiler) CompilerPass {
	switch c.testMode {
	case LocaleDataPassesTest_TestMode_PROTECT_DATA:
		return &localeDataPassesTest_protectPass{compiler: compiler}
	default: // REPLACE_PROTECTED_DATA
		return &localeDataPassesTest_replacePass{compiler: compiler}
	}
}

// Java: static class LocaleDataPassesTest.LocaleResult
type LocaleDataPassesTest_LocaleResult struct {
	locale string
	result *CompilerTestCase_Expected
}

// Java: LocaleDataPassesTest.LocaleResult(String, String)
func NewLocaleDataPassesTest_LocaleResult(locale string, result string) *LocaleDataPassesTest_LocaleResult {
	return &LocaleDataPassesTest_LocaleResult{locale: locale, result: CompilerTestCase_Expected_(result)}
}

// Java: LocaleDataPassesTest.LocaleResult(String, Expected)
func NewLocaleDataPassesTest_LocaleResultOfExpected(locale string, result *CompilerTestCase_Expected) *LocaleDataPassesTest_LocaleResult {
	return &LocaleDataPassesTest_LocaleResult{locale: locale, result: result}
}

// multiTest is the primary test method to use in this file.
//
// originalJs is the original, input JS code; protectedJs is what the code
// should look like after message definitions are is protected from
// mangling during optimizations, but before they are replaced with the
// localized message data; allExpected is what the code should look like
// after full replacement with localized messages has been done.
//
// Java: LocaleDataPassesTest.multiTest(Sources, Expected, LocaleResult...)
func (c *LocaleDataPassesTest) multiTestSources(originalJs CompilerTestCase_Sources, protectedJs *CompilerTestCase_Expected, allExpected ...*LocaleDataPassesTest_LocaleResult) {
	// The PROTECT_DATA mode needs to add externs for the protection functions.
	c.AllowExternsChanges()
	c.testMode = LocaleDataPassesTest_TestMode_PROTECT_DATA
	c.Test(originalJs, protectedJs)
	for _, expected := range allExpected {
		c.testMode = LocaleDataPassesTest_TestMode_REPLACE_PROTECTED_DATA
		c.locale = expected.locale
		c.Test(originalJs, expected.result)
	}
}

// Java: LocaleDataPassesTest.multiTest(String, String, LocaleResult...)
func (c *LocaleDataPassesTest) multiTest(originalJs string, protectedJs string, allExpected ...*LocaleDataPassesTest_LocaleResult) {
	c.multiTestSources(c.Srcs(originalJs), c.Expected(protectedJs), allExpected...)
}

// Java: LocaleDataPassesTest.testBaseJsGoogLocaleRef()
func (c *LocaleDataPassesTest) TestBaseJsGoogLocaleRef() {
	// We're confirming that there won't be any error reported for the use of `goog.LOCALE`.
	c.multiTest( //
		"/**\n"+
			" * @fileoverview\n"+
			" * @provideGoog // no @localeFile, but base.js has this special annotation\n"+
			" */\n"+
			"goog.provide('some.Obj');\n"+
			"goog.LOCALE = 'en';\n"+
			"console.log(goog.LOCALE);\n",
		"/**\n"+
			" * @fileoverview\n"+
			" * @provideGoog // no @localeFile, but base.js has this special annotation\n"+
			" */\n"+
			"goog.provide('some.Obj');\n"+
			"goog.LOCALE = __JSC_LOCALE__;\n"+
			"console.log(__JSC_LOCALE__);\n",
		NewLocaleDataPassesTest_LocaleResult(
			"es_ES",
			"/**\n"+
				" * @fileoverview\n"+
				" * @provideGoog // no @localeFile, but base.js has this special annotation\n"+
				" */\n"+
				"goog.provide('some.Obj');\n"+
				"goog.LOCALE = 'es_ES';\n"+
				"console.log('es_ES');\n"))
}

// Port of com/google/javascript/jscomp/ReportUntranspilableFeaturesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ReportUntranspilableFeaturesTest.java)

/*
 * Copyright 2014 The Closure Compiler Authors.
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
)

type ReportUntranspilableFeaturesTest struct {
	CompilerTestCaseBase
	languageOut           CompilerOptions_LanguageMode
	browserFeaturesetYear *CompilerOptions_BrowserFeaturesetYear
}

// Java: ReportUntranspilableFeaturesTest()
func NewReportUntranspilableFeaturesTest() *ReportUntranspilableFeaturesTest {
	c := &ReportUntranspilableFeaturesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestReportUntranspilableFeaturesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewReportUntranspilableFeaturesTest() })
}

// Java: ReportUntranspilableFeaturesTest.setUp()
func (c *ReportUntranspilableFeaturesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT3
	c.browserFeaturesetYear = nil
	c.DisableTypeCheck()
	c.EnableRunTypeCheckAfterProcessing()
}

// Java: ReportUntranspilableFeaturesTest.getProcessor(Compiler)
func (c *ReportUntranspilableFeaturesTest) GetProcessor(compiler *Compiler) CompilerPass {
	var outputFeatures *parser.FeatureSet
	if c.browserFeaturesetYear != nil {
		outputFeatures = c.browserFeaturesetYear.GetFeatureSet()
	} else {
		outputFeatures = c.languageOut.ToFeatureSet()
	}
	return PeepholeTranspilationsPass_Create(
		compiler,
		[]AbstractPeepholeTranspilation{
			NewReportUntranspilableFeatures(
				compiler,
				c.browserFeaturesetYear,
				outputFeatures)})
}

// Java: ReportUntranspilableFeaturesTest.testEs2018RegexFlagS()
func (c *ReportUntranspilableFeaturesTest) TestEs2018RegexFlagS() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2018
	c.TestSame("const a = /asdf/;")
	c.TestSame("const a = /asdf/g;")
	c.TestSame("const a = /asdf/s;")
	c.TestSame("const a = /asdf/gs;")

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2017
	c.TestSame("const a = /asdf/;")
	c.TestSame("const a = /asdf/g;")
	c.TestError(
		"const a = /asdf/s;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		"Cannot convert feature \"RegExp flag 's'\" to targeted output language. Feature requires at minimum ECMASCRIPT_2018. Consider targeting a more modern output.")
	c.TestError(
		"const a = /asdf/gs;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		"Cannot convert feature \"RegExp flag 's'\" to targeted output language. Feature requires at minimum ECMASCRIPT_2018. Consider targeting a more modern output.")
}

// Java: ReportUntranspilableFeaturesTest.testInvalidRegExpReportsWarning()
func (c *ReportUntranspilableFeaturesTest) TestInvalidRegExpReportsWarning() {
	c.TestWarning("const a = /([0-9a-zA-Z_\\-]{20,}/", CheckRegExp_MALFORMED_REGEXP)
}

// Java: ReportUntranspilableFeaturesTest.testEs2018RegexLookbehind()
func (c *ReportUntranspilableFeaturesTest) TestEs2018RegexLookbehind() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2018
	c.TestSame("const a = /(?<=asdf)/;")
	c.TestSame("const a = /(?<!asdf)/;")

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2017
	c.TestSame("const a = /(?=asdf)/;")  // Lookaheads are fine
	c.TestSame("const a = /(?!asdf)/g;") // Lookaheads are fine
	c.TestError(
		"const a = /(?<=asdf)/;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		"Cannot convert feature \"RegExp Lookbehind\" to targeted output language. Feature requires at minimum ECMASCRIPT_2018. Consider targeting a more modern output.")
	c.TestError(
		"const a = /(?<!asdf)/;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		"Cannot convert feature \"RegExp Lookbehind\" to targeted output language. Feature requires at minimum ECMASCRIPT_2018. Consider targeting a more modern output.")
}

// Java: ReportUntranspilableFeaturesTest.testEs2018RegexUnicodePropertyEscape()
func (c *ReportUntranspilableFeaturesTest) TestEs2018RegexUnicodePropertyEscape() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2018
	c.TestSame("const a = /\\p{Script=Greek}/u;")
	c.TestSame("const a = /\\P{Script=Greek}/u;")
	c.TestSame("const a = /\\p{Script=Greek}/;") // Without u flag, /\p/ is same as /p/
	c.TestSame("const a = /\\P{Script=Greek}/;") // Without u flag, /\p/ is same as /p/

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2017
	c.TestSame("const a = /\\p{Script=Greek}/;") // Without u flag, /\p/ is same as /p/
	c.TestSame("const a = /\\P{Script=Greek}/;") // Without u flag, /\p/ is same as /p/
	c.TestError(
		"const a = /\\p{Script=Greek}/u;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		"Cannot convert feature \"RegExp unicode property escape\" to targeted output language. Feature requires at minimum ECMASCRIPT_2018. Consider targeting a more modern output.")
	c.TestError(
		"const a = /\\P{Script=Greek}/u;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		"Cannot convert feature \"RegExp unicode property escape\" to targeted output language. Feature requires at minimum ECMASCRIPT_2018. Consider targeting a more modern output.")
}

// Java: ReportUntranspilableFeaturesTest.testRegExpConstructorCalls()
func (c *ReportUntranspilableFeaturesTest) TestRegExpConstructorCalls() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2017
	// TODO(bradfordcsmith): report errors from RegExp in this form
	c.TestSame("const a = new RegExp('asdf', 'gs');")
}

// Java: ReportUntranspilableFeaturesTest.testEs2018RegexNamedCaptureGroups()
func (c *ReportUntranspilableFeaturesTest) TestEs2018RegexNamedCaptureGroups() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2018
	c.TestSame("const a = /(?<name>)/u;")

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2017
	c.TestError(
		"const a = /(?<name>)/;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		"Cannot convert feature \"RegExp named groups\" to targeted output language. Feature requires at minimum ECMASCRIPT_2018. Consider targeting a more modern output.")
	c.TestError(
		"const a = /(?<$var>).*/u;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		"Cannot convert feature \"RegExp named groups\" to targeted output language. Feature requires at minimum ECMASCRIPT_2018. Consider targeting a more modern output.")
	// test valid regex with '<' or '>' that is not named capture group
	c.TestSame("const a = /(<name>)/;")
	c.TestSame("const a = /(>.>)/u;")
}

// Java: ReportUntranspilableFeaturesTest.testEs2018RegexNamedCaptureGroupsBackReferencing()
func (c *ReportUntranspilableFeaturesTest) TestEs2018RegexNamedCaptureGroupsBackReferencing() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2018
	c.TestSame("const a = /^(?<half>.*).\\k<half>$/u;")

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2017
	c.TestError(
		"const a = /^(?<half>.*).\\k<half>$/u;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		"Cannot convert feature \"RegExp named groups\" to targeted output language. Feature requires at minimum ECMASCRIPT_2018. Consider targeting a more modern output.")
}

// Java: ReportUntranspilableFeaturesTest.testEs2018RegexNamedCaptureGroupsBackReferencing_usingBrowserFeaturesetYear()
func (c *ReportUntranspilableFeaturesTest) TestEs2018RegexNamedCaptureGroupsBackReferencing_usingBrowserFeaturesetYear() {
	c.browserFeaturesetYear = CompilerOptions_BrowserFeaturesetYear_YEAR_2021.Ptr()
	c.TestSame("const a = /^(?<half>.*).\\k<half>$/u;")

	c.browserFeaturesetYear = CompilerOptions_BrowserFeaturesetYear_YEAR_2020.Ptr()
	c.TestError(
		"const a = /^(?<half>.*).\\k<half>$/u;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		`Cannot convert feature "RegExp named groups" to targeted output language. Feature requires at minimum browser featureset year 2021. Consider targeting a more modern output.
Current browser featureset year: 2020`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ReportUntranspilableFeaturesTest.reportErrorWithBigIntLiteralTranspilation()
func (c *ReportUntranspilableFeaturesTest) TestReportErrorWithBigIntLiteralTranspilation() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2021
	c.TestSame("1234n")

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT3
	c.TestError("1234n", ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ReportUntranspilableFeaturesTest.reportErrorWithBigIntLiteralTranspilation_browserFeaturesetYear()
func (c *ReportUntranspilableFeaturesTest) TestReportErrorWithBigIntLiteralTranspilation_browserFeaturesetYear() {
	c.browserFeaturesetYear = CompilerOptions_BrowserFeaturesetYear_YEAR_2021.Ptr()
	c.TestSame("1234n")

	c.browserFeaturesetYear = CompilerOptions_BrowserFeaturesetYear_YEAR_2020.Ptr()
	c.TestError(
		"1234n",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		`Cannot convert feature "bigint" to targeted output language. Feature requires at minimum browser featureset year 2021. Consider targeting a more modern output.
Current browser featureset year: 2020`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ReportUntranspilableFeaturesTest.noErrorWithBigIntConstructorTranspilation()
func (c *ReportUntranspilableFeaturesTest) TestNoErrorWithBigIntConstructorTranspilation() {
	// Do not report an error for use of the `BigInt()` method, just the literal form.
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT3
	c.TestSame("BigInt(1234)")
}

// Java: ReportUntranspilableFeaturesTest.testEs2022RegexFlagD()
func (c *ReportUntranspilableFeaturesTest) TestEs2022RegexFlagD() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.TestSame("const a = /^(?<half>.*).\\k<half>$/d;")

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2021
	c.TestError(
		"const a = /^(?<half>.*).\\k<half>$/d;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		"Cannot convert feature \"RegExp flag 'd'\" to targeted output language. Feature requires at minimum ECMASCRIPT_2022. Consider targeting a more modern output.")
}

// Java: ReportUntranspilableFeaturesTest.testEs2022RegexFlagD_usingBrowserFeaturesetYear()
func (c *ReportUntranspilableFeaturesTest) TestEs2022RegexFlagD_usingBrowserFeaturesetYear() {
	c.browserFeaturesetYear = CompilerOptions_BrowserFeaturesetYear_YEAR_2022.Ptr()
	c.TestError(
		"const a = /^(?<half>.*).\\k<half>$/d;",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
		`Cannot convert feature "RegExp flag 'd'" to targeted output language. Feature requires at minimum browser featureset year 2023. Consider targeting a more modern output.
Current browser featureset year: 2022`)
}

// Java: ReportUntranspilableFeaturesTest.testEs2022TopLevelAwait()
func (c *ReportUntranspilableFeaturesTest) TestEs2022TopLevelAwait() {
	testTopLevelAwaitUnsupported := func() {
		c.TestError(
			"export const x = await 5;",
			ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT,
			`The feature "Top-level await" is currently unsupported for transpilation.
`)
	}

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	testTopLevelAwaitUnsupported()

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2022
	testTopLevelAwaitUnsupported()

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2021
	testTopLevelAwaitUnsupported()

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT3 // Reset

	c.browserFeaturesetYear = CompilerOptions_BrowserFeaturesetYear_YEAR_2023.Ptr()
	testTopLevelAwaitUnsupported()

	c.browserFeaturesetYear = CompilerOptions_BrowserFeaturesetYear_YEAR_2022.Ptr()
	testTopLevelAwaitUnsupported()
}

// Java: ReportUntranspilableFeaturesTest.testEs2022TopLevelForAwaitOf()
func (c *ReportUntranspilableFeaturesTest) TestEs2022TopLevelForAwaitOf() {
	testTopLevelAwaitUnsupported := func() {
		c.TestError(
			"export const x = 1; for await (let a of b) foo();",
			ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT,
			`The feature "Top-level await" is currently unsupported for transpilation.
`)
	}

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	testTopLevelAwaitUnsupported()

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2022
	testTopLevelAwaitUnsupported()

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2021
	testTopLevelAwaitUnsupported()

	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT3 // Reset

	c.browserFeaturesetYear = CompilerOptions_BrowserFeaturesetYear_YEAR_2023.Ptr()
	testTopLevelAwaitUnsupported()

	c.browserFeaturesetYear = CompilerOptions_BrowserFeaturesetYear_YEAR_2022.Ptr()
	testTopLevelAwaitUnsupported()
}

// Java: ReportUntranspilableFeaturesTest.testCannotConvertEs5GetterToEs3()
func (c *ReportUntranspilableFeaturesTest) TestCannotConvertEs5GetterToEs3() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT3
	c.TestError("var o = {get x() {}};", ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
}

// Java: ReportUntranspilableFeaturesTest.testCannotConvertEs5SetterToEs3()
func (c *ReportUntranspilableFeaturesTest) TestCannotConvertEs5SetterToEs3() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT3
	c.TestError("var o = {set x(val) {}};", ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
}

// Java: ReportUntranspilableFeaturesTest.testAllowsEs5GetterSetterWithEs5Out()
func (c *ReportUntranspilableFeaturesTest) TestAllowsEs5GetterSetterWithEs5Out() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT5
	c.TestNoWarning(
		`var o = {
  get x() {},
  set x(val) {}
};
`)
}

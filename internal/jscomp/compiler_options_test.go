// Port of test/com/google/javascript/jscomp/CompilerOptionsTest.java @ closure-compiler v20261015

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"reflect"
	"regexp"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func compilerOptionsTest_assertFeatureSet(t *testing.T, got *parser.FeatureSet, want *parser.FeatureSet) {
	t.Helper()
	if !got.Equals(want) {
		t.Errorf("getOutputFeatureSet: got %v, want %v", got, want)
	}
}

func TestCompilerOptions_BrowserFeaturesetYearOptionSetsLanguageOut(t *testing.T) {
	options := NewCompilerOptions()
	options.SetBrowserFeaturesetYear(2012)
	compilerOptionsTest_assertFeatureSet(t, options.GetOutputFeatureSet(),
		CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT.ToFeatureSet())

	options.SetBrowserFeaturesetYear(2019)
	compilerOptionsTest_assertFeatureSet(t, options.GetOutputFeatureSet(),
		CompilerOptions_LanguageMode_ECMASCRIPT_2017.ToFeatureSet())

	options.SetBrowserFeaturesetYear(2020)
	compilerOptionsTest_assertFeatureSet(t, options.GetOutputFeatureSet(), parser.FeatureSet_BROWSER_2020)

	options.SetBrowserFeaturesetYear(2021)
	compilerOptionsTest_assertFeatureSet(t, options.GetOutputFeatureSet(), parser.FeatureSet_BROWSER_2021)

	options.SetBrowserFeaturesetYear(2022)
	compilerOptionsTest_assertFeatureSet(t, options.GetOutputFeatureSet(), parser.FeatureSet_BROWSER_2022)

	options.SetBrowserFeaturesetYear(2023)
	compilerOptionsTest_assertFeatureSet(t, options.GetOutputFeatureSet(), parser.FeatureSet_BROWSER_2023)

	options.SetBrowserFeaturesetYear(2024)
	compilerOptionsTest_assertFeatureSet(t, options.GetOutputFeatureSet(), parser.FeatureSet_BROWSER_2024)

	options.SetBrowserFeaturesetYear(2025)
	compilerOptionsTest_assertFeatureSet(t, options.GetOutputFeatureSet(), parser.FeatureSet_BROWSER_2025)

	options.SetBrowserFeaturesetYear(2026)
	compilerOptionsTest_assertFeatureSet(t, options.GetOutputFeatureSet(), parser.FeatureSet_BROWSER_2026)
}

func compilerOptionsTest_defineToken(options *CompilerOptions, name string) rhino.Token {
	return options.GetDefineReplacements().Get(name).GetToken()
}

func TestCompilerOptions_BrowserFeaturesetYearOptionSetsAssumeES5(t *testing.T) {
	options := NewCompilerOptions()
	options.SetBrowserFeaturesetYear(2012)
	assertEqual(t, "2012", compilerOptionsTest_defineToken(options, "$jscomp.ASSUME_ES5"), rhino.Token_FALSE)
	options.SetBrowserFeaturesetYear(2019)
	assertEqual(t, "2019", compilerOptionsTest_defineToken(options, "$jscomp.ASSUME_ES5"), rhino.Token_TRUE)
}

func TestCompilerOptions_BrowserFeaturesetYearOptionSetsAssumeES6(t *testing.T) {
	options := NewCompilerOptions()
	options.SetBrowserFeaturesetYear(2012)
	assertEqual(t, "2012", compilerOptionsTest_defineToken(options, "$jscomp.ASSUME_ES6"), rhino.Token_FALSE)
	options.SetBrowserFeaturesetYear(2018)
	assertEqual(t, "2018", compilerOptionsTest_defineToken(options, "$jscomp.ASSUME_ES6"), rhino.Token_TRUE)
	options.SetBrowserFeaturesetYear(2020)
	assertEqual(t, "2020", compilerOptionsTest_defineToken(options, "$jscomp.ASSUME_ES2020"), rhino.Token_FALSE)
	options.SetBrowserFeaturesetYear(2021)
	assertEqual(t, "2021", compilerOptionsTest_defineToken(options, "$jscomp.ASSUME_ES2020"), rhino.Token_TRUE)
}

func TestCompilerOptions_MinimumBrowserFeatureSetYearRequiredFor(t *testing.T) {
	assertEqual(t, "GETTER", *CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor(parser.FeatureSet_Feature_GETTER),
		CompilerOptions_BrowserFeaturesetYear_YEAR_2012)
	assertEqual(t, "CLASSES", *CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor(parser.FeatureSet_Feature_CLASSES),
		CompilerOptions_BrowserFeaturesetYear_YEAR_2018)
	assertEqual(t, "REGEXP_UNICODE_PROPERTY_ESCAPE",
		*CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor(parser.FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE),
		CompilerOptions_BrowserFeaturesetYear_YEAR_2021)
}

func TestCompilerOptions_MinimumBrowserFeatureSetYearRequiredFor_returnsUnspecifiedIfUnsupported(t *testing.T) {
	// Newer features, in particular anything in ES_NEXT, may not be part of a browser featureset
	// year yet.
	if got := CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor(parser.FeatureSet_Feature_ES_NEXT_RUNTIME); got != nil {
		t.Errorf("got %v, want nil", *got)
	}
}

func TestCompilerOptions_Defines(t *testing.T) {
	options := NewCompilerOptions()
	options.SetDefineToBooleanLiteral("trueVar", true)
	options.SetDefineToBooleanLiteral("falseVar", false)
	options.SetDefineToNumberLiteral("threeVar", 3)
	options.SetDefineToStringLiteral("strVar", "str")

	actual := options.GetDefineReplacements()
	compilerOptionsTest_assertEquivalent(t, rhino.NewNode(rhino.Token_TRUE), actual.Get("trueVar"))
	compilerOptionsTest_assertEquivalent(t, rhino.NewNode(rhino.Token_FALSE), actual.Get("falseVar"))
	compilerOptionsTest_assertEquivalent(t, rhino.Node_NewNumber(3), actual.Get("threeVar"))
	compilerOptionsTest_assertEquivalent(t, rhino.Node_NewString("str"), actual.Get("strVar"))
}

func compilerOptionsTest_assertEquivalent(t *testing.T, a rhino.Node, b rhino.Node) {
	t.Helper()
	if !a.IsEquivalentTo(b) {
		t.Errorf("%v is not equivalent to %v", a, b)
	}
}

func TestCompilerOptions_LanguageModeFromString(t *testing.T) {
	from := func(s string) *CompilerOptions_LanguageMode { return CompilerOptions_LanguageMode_FromString(&s) }
	assertEqual(t, "ECMASCRIPT3", from("ECMASCRIPT3"), CompilerOptions_LanguageMode_ECMASCRIPT3.Ptr())
	// Whitespace should be trimmed, characters converted to uppercase and leading 'ES' replaced
	// with 'ECMASCRIPT'.
	assertEqual(t, "es3", from("  es3  "), CompilerOptions_LanguageMode_ECMASCRIPT3.Ptr())
	assertEqual(t, "junk", from("junk"), nil)

	assertEqual(t, "ECMASCRIPT_2020", from("ECMASCRIPT_2020"), CompilerOptions_LanguageMode_ECMASCRIPT_2020.Ptr())
	// Whitespace should be trimmed, characters converted to uppercase and leading 'ES' replaced
	// with 'ECMASCRIPT'.
	assertEqual(t, "es_2020", from("  es_2020  "), CompilerOptions_LanguageMode_ECMASCRIPT_2020.Ptr())
	assertEqual(t, "es2020", from("  es2020  "), nil) // generates invalid "ECMASCRIPT2020"
	assertEqual(t, "junk", from("junk"), nil)
}

func TestCompilerOptions_EmitUseStrictWorksInEs3(t *testing.T) {
	options := NewCompilerOptions()
	options.SetEmitUseStrict(true)

	if !options.ShouldEmitUseStrict() {
		t.Errorf("shouldEmitUseStrict() is false")
	}
}

func TestCompilerOptions_RemoveRegexFromPath(t *testing.T) {
	options := NewCompilerOptions()
	pattern := options.GetConformanceRemoveRegexFromPath()
	// Go-only: test-helper — pattern.matcher(s).replaceAll("") for this ^-anchored pattern
	// (one match at most), which uses no construct whose RE2 meaning differs from Java's.
	re := regexp.MustCompile(pattern.Pattern())
	replaceAll := func(s string) string { return re.ReplaceAllString(s, "") }
	cases := [][2]string{
		{"/blaze-out/k8-fastbin/bin/some/path", "some/path"},
		{"blaze-out/k8-fastbin/bin/some/path", "some/path"},
		{"google3/blaze-out/k8-fastbin/bin/some/path", "some/path"},
		{"google3/blaze-out/k8-fastbin/genfiles/some/path", "some/path"},
		{"something/google3/blaze-out/k8-fastbin/some/path", "blaze-out/k8-fastbin/some/path"},
		{"/something/google3/blaze-out/k8-fastbin/bin/some/path", "some/path"},
		{"google3/some/path", "some/path"},
		{"google3/foo/blaze-out/some/path", "foo/blaze-out/some/path"},
		{"google3/blaze-out/foo/blaze-out/some/path", "blaze-out/foo/blaze-out/some/path"},
		{"bazel-out/k8-fastbin/bin/some/path", "some/path"},
		{"bazel-out/k8-fastbin/genfiles/some/path", "some/path"},
	}
	for _, c := range cases {
		assertEqual(t, c[0], replaceAll(c[0]), c[1])
	}
}

func TestCompilerOptions_AllInstanceFieldsArePrivate(t *testing.T) {
	// Go-only: test-helper — the Go form of "all instance fields are private" is that no
	// field of the struct is exported.
	typ := reflect.TypeOf(CompilerOptions{})
	nonPrivateInstanceFields := []string{}
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).IsExported() {
			nonPrivateInstanceFields = append(nonPrivateInstanceFields, typ.Field(i).Name)
		}
	}
	if len(nonPrivateInstanceFields) != 0 {
		t.Errorf("All instance fields in CompilerOptions must be private: %v", nonPrivateInstanceFields)
	}
}

// Java: CompilerOptionsTest.testTranspilePublicClassFields()
func TestCompilerOptions_TranspilePublicClassFields(t *testing.T) {
	options := NewCompilerOptions()
	options.SetBrowserFeaturesetYear(2023)
	if options.GetTranspilePublicClassFields() {
		t.Error("getTranspilePublicClassFields: got true, want false")
	}
	if !options.GetOutputFeatureSet().Has(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS) {
		t.Error("getOutputFeatureSet().has(PUBLIC_CLASS_FIELDS): got false, want true")
	}

	options.SetTranspilePublicClassFields(true)
	if !options.GetTranspilePublicClassFields() {
		t.Error("getTranspilePublicClassFields: got false, want true")
	}
	if options.GetOutputFeatureSet().Has(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS) {
		t.Error("getOutputFeatureSet().has(PUBLIC_CLASS_FIELDS): got true, want false")
	}

	options.SetTranspilePublicClassFields(false)
	if options.GetTranspilePublicClassFields() {
		t.Error("getTranspilePublicClassFields: got true, want false")
	}
	if !options.GetOutputFeatureSet().Has(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS) {
		t.Error("getOutputFeatureSet().has(PUBLIC_CLASS_FIELDS): got false, want true")
	}
}

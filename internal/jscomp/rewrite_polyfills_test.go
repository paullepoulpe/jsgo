// Port of com/google/javascript/jscomp/RewritePolyfillsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/RewritePolyfillsTest.java)

/*
 * Copyright 2015 The Closure Compiler Authors.
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

type RewritePolyfillsTest struct {
	CompilerTestCaseBase
	injectableLibraries      *jutil.LinkedHashMap[string, string]
	injectBeforePass         *jutil.LinkedHashSet[string]
	polyfillTable            *jutil.ArrayList[string]
	isolatePolyfills         bool
	injectPolyfills          bool
	injectPolyfillsNewerThan *CompilerOptions_LanguageMode
}

// Java: RewritePolyfillsTest()
func NewRewritePolyfillsTest() *RewritePolyfillsTest {
	c := &RewritePolyfillsTest{}
	c.injectableLibraries = jutil.NewLinkedHashMap[string, string]()
	c.injectBeforePass = jutil.NewLinkedHashSet[string]()
	c.polyfillTable = jutil.NewArrayList[string]()
	c.isolatePolyfills = false
	c.injectPolyfills = true
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewritePolyfillsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewritePolyfillsTest() })
}

// Java: RewritePolyfillsTest.rewritePolyfillsTest_ES_2020
const rewritePolyfillsTest_ES_2020 = CompilerOptions_LanguageMode_ECMASCRIPT_2020

// Java: RewritePolyfillsTest.rewritePolyfillsTest_ES6
const rewritePolyfillsTest_ES6 = CompilerOptions_LanguageMode_ECMASCRIPT_2015

// Java: RewritePolyfillsTest.rewritePolyfillsTest_ES5
const rewritePolyfillsTest_ES5 = CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT

// Java: RewritePolyfillsTest.rewritePolyfillsTest_ES3
const rewritePolyfillsTest_ES3 = CompilerOptions_LanguageMode_ECMASCRIPT3

// Java: RewritePolyfillsTest.addLibrary(String, String, String, String)
func (c *RewritePolyfillsTest) AddLibrary(name string, from string, to string, library *string) {
	if library != nil {
		c.injectableLibraries.Put(
			*library,
			"$jscomp.polyfill('"+name+"', function() {}, '"+from+"', '"+to+"');\n")
	}
	// Go-only: test-helper — String.format("%s %s %s %s", ...) with
	// Strings.nullToEmpty(library).
	nullToEmpty := ""
	if library != nil {
		nullToEmpty = *library
	}
	c.polyfillTable.Add(name + " " + from + " " + to + " " + nullToEmpty)
}

// Java: RewritePolyfillsTest.setUp()
func (c *RewritePolyfillsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.injectableLibraries.Clear()
	c.injectBeforePass.Clear()
	c.polyfillTable.Clear()
	c.injectPolyfillsNewerThan = nil
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
}

// Java: RewritePolyfillsTest.getProcessor(Compiler)
func (c *RewritePolyfillsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRewritePolyfillsWithRuntimeJsLibManager(
		compiler,
		c.CreateRuntimeJsLibManager(compiler),
		PolyfillUsageFinder_Polyfills_FromTable(strings.Join(c.polyfillTable.ToSlice(), "\n")),
		c.injectPolyfills,
		c.isolatePolyfills,
		c.injectPolyfillsNewerThan)
}

// Java: RewritePolyfillsTest.getOptions()
func (c *RewritePolyfillsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MISSING_POLYFILL, CheckLevel_WARNING)
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
	return options
}

// Java: RewritePolyfillsTest.createRuntimeJsLibManager(Compiler)
func (c *RewritePolyfillsTest) CreateRuntimeJsLibManager(compiler *Compiler) *RuntimeJsLibManager {
	runtimeLibs :=
		RuntimeJsLibManager_Create(
			RuntimeJsLibManager_RuntimeLibraryMode_INJECT,
			// stub out the resource parsing
			func(resource string, path string) rhino.Node {
				return compiler.ParseTestCode(c.injectableLibraries.Get(resource))
			},
			compiler.GetChangeTracker(),
			func() rhino.Node { return compiler.GetNodeForCodeInsertion(nil) })
	for toInject := range c.injectBeforePass.All() {
		runtimeLibs.EnsureLibraryInjected(toInject /* force= */, false)
	}
	return runtimeLibs
}

// Java: RewritePolyfillsTest.addLibraries(String, String[])
func (c *RewritePolyfillsTest) AddLibraries(code string, libraries []string) string {
	expected := jutil.NewStringBuilder()
	for _, library := range libraries {
		expected.Append(c.injectableLibraries.Get(library))
	}
	expected.Append(code)
	return expected.String()
}

// Java: RewritePolyfillsTest.testDoesNotInject(String)
func (c *RewritePolyfillsTest) TestDoesNotInject(code string) {
	c.TestInjects(code) // empty list of injections
}

// Java: RewritePolyfillsTest.testInjects(String, String...)
func (c *RewritePolyfillsTest) TestInjects(code string, libraries ...string) {
	c.Test(code, c.AddLibraries(code, libraries))
}

// Java: RewritePolyfillsTest.testInjects(String, DiagnosticType, String...)
func (c *RewritePolyfillsTest) TestInjectsWithWarning(code string, warning *DiagnosticType, libraries ...string) {
	c.Test(code, c.AddLibraries(code, libraries), c.Warning(warning))
}

// Java: RewritePolyfillsTest.testEmpty()
func (c *RewritePolyfillsTest) TestEmpty() {
	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestInjects("")
}

// Java: RewritePolyfillsTest.testClassesInjected()
func (c *RewritePolyfillsTest) TestClassesInjected() {
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))
	c.AddLibrary("Set", "es6", "es3", jutil.StringPtr("es6/set"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestInjects("var m = new Map();", "es6/map")
	c.TestInjects("var m = new goog.global.Map();", "es6/map")
	c.TestInjects("var Map = goog.global.Map; new Map();", "es6/map")
	c.TestInjects("var m = new window.Map();", "es6/map")

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES3)
	c.TestInjects("var s = new Set();", "es6/set")
}

// Java: RewritePolyfillsTest.testLibrariesOnlyInjectedOnce()
func (c *RewritePolyfillsTest) TestLibrariesOnlyInjectedOnce() {
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestInjects("var m = new Map(); m = new Map();", "es6/map")
}

// Java: RewritePolyfillsTest.testClassesNotInjectedIfSufficientLanguageOut()
func (c *RewritePolyfillsTest) TestClassesNotInjectedIfSufficientLanguageOut() {
	c.AddLibrary("Proxy", "es6", "es6", nil)
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))
	c.AddLibrary("Set", "es6", "es3", jutil.StringPtr("es6/set"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES6)
	c.TestDoesNotInject("new Proxy();")
	c.TestDoesNotInject("var m = new Map();")
	c.TestDoesNotInject("new Set();")
}

// Java: RewritePolyfillsTest.testClassesNotInjectedIfDeclaredInScope()
func (c *RewritePolyfillsTest) TestClassesNotInjectedIfDeclaredInScope() {
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestDoesNotInject("/** @constructor */ var Map = function() {}; new Map();")
}

// Java: RewritePolyfillsTest.testClassesWarnIfInsufficientLanguageOut()
func (c *RewritePolyfillsTest) TestClassesWarnIfInsufficientLanguageOut() {
	c.AddLibrary("Proxy", "es6", "es6", nil)
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestInjectsWithWarning("new Proxy();", RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR)

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES3)
	c.TestInjectsWithWarning("new Map();", RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR, "es6/map")
}

// Java: RewritePolyfillsTest.testGlobalThisInjected()
func (c *RewritePolyfillsTest) TestGlobalThisInjected() {
	c.AddLibrary("globalThis", "es_2020", "es3", jutil.StringPtr("es6/globalThis"))
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.TestInjects("globalThis.String", "es6/globalThis")
	c.TestInjects("var m = new globalThis.Map();", "es6/globalThis", "es6/map")
}

// Java: RewritePolyfillsTest.testGlobalThisInjectedNotInjectedGivenSufficientLanguageOut()
func (c *RewritePolyfillsTest) TestGlobalThisInjectedNotInjectedGivenSufficientLanguageOut() {
	c.AddLibrary("globalThis", "es_2020", "es3", jutil.StringPtr("es6/globalThis"))
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.SetLanguage(rewritePolyfillsTest_ES_2020, rewritePolyfillsTest_ES_2020)
	c.TestDoesNotInject("globalThis.String")
	c.TestDoesNotInject("var m = new globalThis.Map();")
}

// Java: RewritePolyfillsTest.testStaticMethodsInjected()
func (c *RewritePolyfillsTest) TestStaticMethodsInjected() {
	c.AddLibrary("Math.clz32", "es6", "es5", jutil.StringPtr("es6/math/clz32"))
	c.AddLibrary("Array.of", "es6", "es3", jutil.StringPtr("es6/array/of"))
	c.AddLibrary("Object.keys", "es5", "es3", jutil.StringPtr("es5/object/keys"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestInjects("Math.clz32(x);", "es6/math/clz32")

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES3)
	c.TestInjects("Array.of(x);", "es6/array/of")

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES3)
	c.TestInjects("Object.keys(x);", "es5/object/keys")
}

// Java: RewritePolyfillsTest.testStaticMethodsNotInjectedIfSufficientLanguageOut()
func (c *RewritePolyfillsTest) TestStaticMethodsNotInjectedIfSufficientLanguageOut() {
	c.AddLibrary("Array.from", "es6", "es6", nil)
	c.AddLibrary("Math.clz32", "es6", "es5", jutil.StringPtr("es6/math/clz32"))
	c.AddLibrary("Array.of", "es6", "es3", jutil.StringPtr("es6/array/of"))
	c.AddLibrary("Object.keys", "es5", "es3", jutil.StringPtr("es5/object/keys"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES6)
	c.TestDoesNotInject("Array.from(x);")
	c.TestDoesNotInject("Math.clz32(x);")
	c.TestDoesNotInject("Array.of(x);")

	c.SetLanguage(rewritePolyfillsTest_ES5, rewritePolyfillsTest_ES5)
	c.TestDoesNotInject("Object.keys(x);")
}

// Java: RewritePolyfillsTest.testStaticMethodsNotInjectedIfDeclaredInScope()
func (c *RewritePolyfillsTest) TestStaticMethodsNotInjectedIfDeclaredInScope() {
	c.AddLibrary("Math.clz32", "es6", "es5", jutil.StringPtr("es6/math/clz32"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestDoesNotInject("var Math = {clz32: function() {}}; Math.clz32(x);")
}

// Java: RewritePolyfillsTest.testStaticMethodsWarnIfInsufficientLanguageOut()
func (c *RewritePolyfillsTest) TestStaticMethodsWarnIfInsufficientLanguageOut() {
	c.AddLibrary("Array.from", "es6", "es6", nil)
	c.AddLibrary("Math.clz32", "es6", "es5", jutil.StringPtr("es6/math/clz32"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestInjectsWithWarning("Array.from(x);", RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR)

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES3)
	c.TestInjectsWithWarning(
		"Math.clz32(x);", RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR, "es6/math/clz32")
}

// Java: RewritePolyfillsTest.testStaticMethodsNotInstalledIfGuardedByIf()
func (c *RewritePolyfillsTest) TestStaticMethodsNotInstalledIfGuardedByIf() {
	c.AddLibrary("Array.of", "es6", "es5", jutil.StringPtr("es6/array/of"))

	c.TestDoesNotInject("if (Array.of) { Array.of(); } else { Array.of(); }")
	c.TestDoesNotInject("if (x || Array.of) { Array.of(x); }")
	c.TestDoesNotInject("if (x && Array.of) { Array.of(x); }")
	c.TestDoesNotInject("if (!Array.of) { Array.of(x); }")
	c.TestDoesNotInject("if (Array.of != 'x') { Array.of(x); }")
	c.TestDoesNotInject("if (Array.of !== 'x') { Array.of(x); }")
	c.TestDoesNotInject("if (Array.of == 'x') { Array.of(x); }")
	c.TestDoesNotInject("if (Array.of === 'x') { Array.of(x); }")
	c.TestDoesNotInject("if (typeof Array.of == 'function') { Array.of(); }")
}

// Java: RewritePolyfillsTest.testStaticMethodsNotInstalledIfGuardedByLogicalOperator()
func (c *RewritePolyfillsTest) TestStaticMethodsNotInstalledIfGuardedByLogicalOperator() {
	c.AddLibrary("Array.of", "es6", "es5", jutil.StringPtr("es6/array/of"))

	c.TestDoesNotInject("Array.of && Array.of();")
	c.TestDoesNotInject("!Array.of || Array.of();")
	// NOTE: needs to be first argument to actually guard.
	c.TestInjects("x && Array.of;", "es6/array/of")
	// NOTE: || is not safe by itself.
	c.TestInjects("Array.of || Array.of();", "es6/array/of")
}

// Java: RewritePolyfillsTest.testStaticMethodsNotInstalledIfGuardedByNullishCoalesce()
func (c *RewritePolyfillsTest) TestStaticMethodsNotInstalledIfGuardedByNullishCoalesce() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_UNSTABLE)
	c.AddLibrary("Array.of", "es_unstable", "es5", jutil.StringPtr("es6/array/of"))

	c.TestDoesNotInject("!Array.of ?? Array.of();")
	// NOTE: ?? is not safe by itself.
	c.SetLanguage(rewritePolyfillsTest_ES_2020, rewritePolyfillsTest_ES5)
	c.TestInjects("Array.of ?? Array.of();", "es6/array/of")
}

// Java: RewritePolyfillsTest.testStaticMethodsNotInstalledIfGuardedByHook()
func (c *RewritePolyfillsTest) TestStaticMethodsNotInstalledIfGuardedByHook() {
	c.AddLibrary("Array.of", "es6", "es5", jutil.StringPtr("es6/array/of"))

	c.TestDoesNotInject("var x = Array.of ? y : function(z) { Array.of(z); };")
	c.TestDoesNotInject("var x = Array.of ? function(y) { Array.of(y); } : z;")
	c.TestDoesNotInject("typeof Array.of ? Array.of(x) : Array.of(y);")
	c.TestDoesNotInject("String(Array.of) == 'foo' ? Array.of(x) : Array.of(y);")
	c.TestDoesNotInject("Array.of instanceof Function ? Array.of(x) : Array.of(y);")
	// NOTE: needs to be first argument to actually guard.
	c.TestInjects("x ? (Array.of ? y : z) : Array.of(x)", "es6/array/of")
}

// Java: RewritePolyfillsTest.testStaticMethodsNotInstalledIfGuardedByAbruptReturn()
func (c *RewritePolyfillsTest) TestStaticMethodsNotInstalledIfGuardedByAbruptReturn() {
	c.AddLibrary("Array.of", "es6", "es5", jutil.StringPtr("es6/array/of"))

	c.TestDoesNotInject("if (!Array.of) throw 'x'; Array.of();")
	c.TestDoesNotInject("!function() { if (!Array.of) return; Array.of(); }")
	c.TestDoesNotInject("if (!Array.of) { throw 'x'; } Array.of();")
	// NOTE: abrupt return must be a sibling conditional's direct child
	c.TestInjects("{ if (!Array.of) throw 'x'; } Array.of();", "es6/array/of")
	c.TestInjects("!function() { { if (!Array.of) return; } Array.of(); }()", "es6/array/of")
	c.TestInjects("if (!Array.of) { { throw 'x'; } } Array.of();", "es6/array/of")
	c.TestInjects("{ if (!Array.of) throw 'x'; throw 'x'; } Array.of();", "es6/array/of")
	c.TestInjects("{ if (Array.of) throw 'x'; { throw 'x'; } } Array.of();", "es6/array/of")
	c.TestInjects(
		"if (unrelated) { if (Array.of) throw 'x'; { throw 'x'; } } else Array.of();",
		"es6/array/of")
}

// Java: RewritePolyfillsTest.testPrototypeMethodsInjected()
func (c *RewritePolyfillsTest) TestPrototypeMethodsInjected() {
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))
	c.AddLibrary("Array.prototype.fill", "es6", "es3", jutil.StringPtr("es6/array/fill"))
	c.AddLibrary("Array.prototype.forEach", "es5", "es3", jutil.StringPtr("es5/array/foreach"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestInjects("x.endsWith(y);", "es6/string/endswith")
	c.TestInjects("x.fill();", "es6/array/fill")

	c.SetLanguage(rewritePolyfillsTest_ES5, rewritePolyfillsTest_ES3)
	c.TestInjects("x.forEach(y);", "es5/array/foreach")
}

// Java: RewritePolyfillsTest.testPrototypeMethodsNotInjectedIfSufficientLanguageOut()
func (c *RewritePolyfillsTest) TestPrototypeMethodsNotInjectedIfSufficientLanguageOut() {
	c.AddLibrary("String.prototype.normalize", "es6", "es6", nil)
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))
	c.AddLibrary("Array.prototype.fill", "es6", "es3", jutil.StringPtr("es6/array/fill"))
	c.AddLibrary("Array.prototype.forEach", "es5", "es3", jutil.StringPtr("es5/array/foreach"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES6)
	c.TestDoesNotInject("x.normalize();")
	c.TestDoesNotInject("x.endsWith();")
	c.TestDoesNotInject("x.fill(y);")

	c.SetLanguage(rewritePolyfillsTest_ES5, rewritePolyfillsTest_ES5)
	c.TestDoesNotInject("x.forEach();")
}

// Java: RewritePolyfillsTest.testPrototypeMethodsDontWarnIfInsufficientVersionNumber()
func (c *RewritePolyfillsTest) TestPrototypeMethodsDontWarnIfInsufficientVersionNumber() {
	c.AddLibrary("String.prototype.normalize", "es6", "es6", nil)
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))
	c.AddLibrary("Array.prototype.fill", "es6", "es3", jutil.StringPtr("es6/array/fill"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES3)
	c.TestDoesNotInject("x.normalize();")
	c.TestInjects("x.endsWith();", "es6/string/endswith")
	c.TestInjects("x.fill(y);", "es6/array/fill")
}

// Java: RewritePolyfillsTest.testMultiplePrototypeMethodsWithSameName()
func (c *RewritePolyfillsTest) TestMultiplePrototypeMethodsWithSameName() {
	c.AddLibrary("Array.prototype.includes", "es6", "es3", jutil.StringPtr("es6/array/includes"))
	c.AddLibrary("String.prototype.includes", "es5", "es3", jutil.StringPtr("es5/string/includes"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestInjects("x.includes();", "es6/array/includes")

	c.SetLanguage(rewritePolyfillsTest_ES5, rewritePolyfillsTest_ES3)
	c.TestInjects("x.includes();", "es6/array/includes", "es5/string/includes")
}

// Java: RewritePolyfillsTest.testPrototypeMethodsInstalledIfStaticMethodShadowed()
func (c *RewritePolyfillsTest) TestPrototypeMethodsInstalledIfStaticMethodShadowed() {
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestInjects(
		`var string = {}; string.endsWith = function() {};
string.foo = function(string) { return string.endsWith('x'); };
`,
		"es6/string/endswith")
}

// NOTE(sdh): it's not clear what makes the most sense here.  At one point we
// took care to avoid installing these, but it may make sense to instead leave
// this distinction to a type-based optimization.  As such, I've simplified the
// logic to no longer look at variables' scope and instead just skiplist known
// symbols like goog.string and goog.array.
//
// Java: RewritePolyfillsTest.testPrototypeMethodsInstalledIfActuallyStatic()
func (c *RewritePolyfillsTest) TestPrototypeMethodsInstalledIfActuallyStatic() {
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.TestInjects(
		"var string = {}; string.endsWith = function() {}; string.endsWith('x');",
		"es6/string/endswith")
	c.TestInjects(
		"var string = {endsWith: function() {}}; string.endsWith('x');", "es6/string/endswith")
	c.TestInjects(
		`var string = {}; string.endsWith = function() {};
string.foo = function() { return string.endsWith('x'); };
`,
		"es6/string/endswith")
}

// Java: RewritePolyfillsTest.testPrototypeMethodsNotInstalledIfGuardedByIf()
func (c *RewritePolyfillsTest) TestPrototypeMethodsNotInstalledIfGuardedByIf() {
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.TestDoesNotInject("if (String.prototype.endsWith) { x.endsWith(); } else { y.endsWith(); }")
	c.TestDoesNotInject("if (x || String.prototype.endsWith) { x.endsWith(); }")
	c.TestDoesNotInject("if (x && String.prototype.endsWith) { x.endsWith(); }")
	c.TestDoesNotInject("if (!String.prototype.endsWith) { x.endsWith(); }")
	c.TestDoesNotInject("if (String.prototype.endsWith != 'x') { x.endsWith(); }")
	c.TestDoesNotInject("if (String.prototype.endsWith !== 'x') { x.endsWith(); }")
	c.TestDoesNotInject("if (String.prototype.endsWith == 'x') { x.endsWith(); }")
	c.TestDoesNotInject("if (String.prototype.endsWith === 'x') { x.endsWith(); }")
	c.TestDoesNotInject("if (typeof String.prototype.endsWith == 'function') { x.endsWith(); }")
}

// Java: RewritePolyfillsTest.testPrototypeMethodsNotInstalledIfGuardedByLogicalOperator()
func (c *RewritePolyfillsTest) TestPrototypeMethodsNotInstalledIfGuardedByLogicalOperator() {
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.TestDoesNotInject("String.prototype.endsWith && x.endsWith();")
	c.TestDoesNotInject("!String.prototype.endsWith || x.endsWith();")
	c.TestDoesNotInject("x.endsWith && x.endsWith();")
	// NOTE: needs to be first argument to actually guard.
	c.TestInjects("x && x.endsWith;", "es6/string/endswith")
	// NOTE: || is not safe by itself.
	c.TestInjects("String.prototype.endsWith || x.endsWith();", "es6/string/endswith")
}

// Java: RewritePolyfillsTest.testPrototypeMethodsNotInstalledIfGuardedByNullishCoalesce()
func (c *RewritePolyfillsTest) TestPrototypeMethodsNotInstalledIfGuardedByNullishCoalesce() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_UNSTABLE)
	c.AddLibrary("String.prototype.endsWith", "es_unstable", "es5", jutil.StringPtr("es6/string/endswith"))

	c.TestDoesNotInject("!String.prototype.endsWith ?? x.endsWith();")
	// NOTE: ?? is not safe by itself.
	c.SetLanguage(rewritePolyfillsTest_ES_2020, rewritePolyfillsTest_ES5)
	c.TestInjects("String.prototype.endsWith ?? x.endsWith();", "es6/string/endswith")
}

// Java: RewritePolyfillsTest.testPrototypeMethodsNotInstalledIfGuardedByHook()
func (c *RewritePolyfillsTest) TestPrototypeMethodsNotInstalledIfGuardedByHook() {
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.TestDoesNotInject("var x = String.prototype.endsWith ? y : function(z) { z.endsWith(); };")
	c.TestDoesNotInject("var x = String.prototype.endsWith ? function(y) { y.endsWith(); } : z;")
	c.TestDoesNotInject("typeof String.prototype.endsWith ? x.endsWith() : y.endsWith();")
	c.TestDoesNotInject("String(x.endsWith) == 'foo' ? x.endsWith() : y.endsWith();")
	c.TestDoesNotInject("Boolean(x.endsWith) ? x.endsWith() : y.endsWith();")
	// NOTE: needs to be first argument to actually guard.
	c.TestInjects("x ? (String.prototype.endsWith ? y : z) : x.endsWith()", "es6/string/endswith")
}

// Java: RewritePolyfillsTest.testPrototypeMethodsNotInstalledIfGuardedByAbruptReturn()
func (c *RewritePolyfillsTest) TestPrototypeMethodsNotInstalledIfGuardedByAbruptReturn() {
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.TestDoesNotInject("if (!x.endsWith) throw 'x'; y.endsWith();")
	c.TestDoesNotInject("if (!x.endsWith) { throw 'x'; } y.endsWith();")
	c.TestDoesNotInject("!function() { if (!x.endsWith) return; y.endsWith(); }()")
	// NOTE: abrupt return must be a sibling conditional's direct child
	c.TestInjects("{ if (!x.endsWith) throw 'x'; } y.endsWith();", "es6/string/endswith")
	c.TestInjects("if (!x.endsWith) { { throw 'x'; } } y.endsWith();", "es6/string/endswith")
	c.TestInjects(
		"!function() { { if (!x.endsWith) return; } y.endsWith(); }()", "es6/string/endswith")
	c.TestInjects("{ if (!x.endsWith) throw 'x'; throw 'x'; } y.endsWith();", "es6/string/endswith")
	c.TestInjects(
		"{ if (x.endsWith) throw 'x'; { throw 'x'; } } x.endsWith();", "es6/string/endswith")
	c.TestInjects(
		"if (unrelated) { if (x.endsWith) throw 'x'; { throw 'x'; } } else x.endsWith();",
		"es6/string/endswith")
}

// Java: RewritePolyfillsTest.testCleansUpUnnecessaryPolyfills()
func (c *RewritePolyfillsTest) TestCleansUpUnnecessaryPolyfills() {
	// Put two polyfill statements in the same library.
	c.injectableLibraries.Put(
		"es6/set",
		`$jscomp.polyfill('Set', function() {}, 'es6', 'es3');
$jscomp.polyfill('Map', function() {}, 'es5', 'es3');
`)
	c.polyfillTable.Add("Set es6 es3 es6/set")

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.Test(
		"var set = new Set();",
		`$jscomp.polyfill('Set', function() {}, 'es6', 'es3');
var set = new Set();
`)

	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES3)
	c.Test(
		"var set = new Set();",
		`$jscomp.polyfill('Set', function() {}, 'es6', 'es3');
$jscomp.polyfill('Map', function() {}, 'es5', 'es3');
var set = new Set();
`)
}

// Java: RewritePolyfillsTest.testRegexFeatureSetException()
func (c *RewritePolyfillsTest) TestRegexFeatureSetException() {
	// Tests the special casing done to handle the fact that Safari has lagged behind on
	// implementation of some regex features, which messes with the way polyfills are linked
	// with the feature set of the specification that added them.

	// Promise.prototype.finally was added in ES_2018
	c.AddLibrary("Promise.prototype.finally", "es_2018", "es3", jutil.StringPtr("es6/promise/finally"))

	// input is assumed to contain features from the rewritePolyfillsTest_ES_2020 spec
	c.SetAcceptedLanguage(rewritePolyfillsTest_ES_2020)

	// The output is required to work for the latest browsers as of early 2020
	// This will set an output `FeatureSet` that does not include some regex features that are
	// part of ES_2018, so ES_2018 will not be "supported", and could cause a polyfill
	// that is part of ES_2018 to be included. We have special casing to avoid this.
	year := 2020
	c.SetBrowserFeaturesetYear(&year)

	// the change marking check will report an error because:
	// 1. We add the polyfill code to the SCRIPT node
	// 2. We then strip out the polyfill
	// 3. We mark the SCRIPT node as changed, but it ended up the same as it started.
	c.DisableValidateAstChangeMarking()

	// The polyfill is not included even though the output feature set is not a super set of ES_2018
	c.TestSame("Promise.resolve(1).finally(console.log('done'));")
}

// Java: RewritePolyfillsTest.testCleansUpUnnecessaryPreviouslyInjectedPolyfills()
func (c *RewritePolyfillsTest) TestCleansUpUnnecessaryPreviouslyInjectedPolyfills() {
	// Put two polyfill statements in the same library.
	c.injectableLibraries.Put(
		"es6/set",
		`$jscomp.polyfill('Set', function() {}, 'es6', 'es3');
// pretend Map isn't needed for rewritePolyfillsTest_ES5
$jscomp.polyfill('Map', function() {}, 'es5', 'es3');
`)
	c.polyfillTable.Add("Set es6 es3 es6/set")

	// simulate injection of Map by a prior-run pass
	c.injectBeforePass.Add("es6/set")
	c.SetLanguage(rewritePolyfillsTest_ES6, rewritePolyfillsTest_ES5)
	c.Test(
		"var set = new Set();",
		` // Map gets removed even though not added by RewritePolyfills
$jscomp.polyfill('Set', function() {}, 'es6', 'es3');
var set = new Set();
`)
}

// Java: RewritePolyfillsTest.testAddsPolyfillMethodToExterns()
func (c *RewritePolyfillsTest) TestAddsPolyfillMethodToExterns() {
	c.isolatePolyfills = true
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.TestExternChanges(c.Srcs(""), c.Expected("var $jscomp$lookupPolyfilledValue"))
}

// Java: RewritePolyfillsTest.testNoCodeChangesIfInjectionDisabled()
func (c *RewritePolyfillsTest) TestNoCodeChangesIfInjectionDisabled() {
	c.isolatePolyfills = true
	c.injectPolyfills = false
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.TestExternChanges(c.Srcs(""), c.Expected("var $jscomp$lookupPolyfilledValue"))

	c.AllowExternsChanges()
	c.TestSame("'x'.endsWith('y');")
}

// Java: RewritePolyfillsTest.testForceInject_es5_addsES6AndES8()
func (c *RewritePolyfillsTest) TestForceInject_es5_addsES6AndES8() {
	c.injectPolyfillsNewerThan = CompilerOptions_LanguageMode_ECMASCRIPT5.Ptr()
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))
	c.AddLibrary("Object.values", "es8", "es3", jutil.StringPtr("es6/object/values"))

	c.TestInjects("", "es6/string/endswith", "es6/object/values")
}

// Java: RewritePolyfillsTest.testForceInject_es2015_addsES8Polyfill()
func (c *RewritePolyfillsTest) TestForceInject_es2015_addsES8Polyfill() {
	c.injectPolyfillsNewerThan = CompilerOptions_LanguageMode_ECMASCRIPT5.Ptr()
	c.AddLibrary("Object.values", "es8", "es3", jutil.StringPtr("es6/object/values"))

	c.TestInjects("", "es6/object/values")
}

// Java: RewritePolyfillsTest.testForceInject_es2015_skipsEs2015Polyfills()
func (c *RewritePolyfillsTest) TestForceInject_es2015_skipsEs2015Polyfills() {
	c.injectPolyfillsNewerThan = CompilerOptions_LanguageMode_ECMASCRIPT_2015.Ptr()
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.TestDoesNotInject("")
}

// Java: RewritePolyfillsTest.testIteratorPolyfill()
func (c *RewritePolyfillsTest) TestIteratorPolyfill() {
	c.AddLibrary("Iterator", "es_next", "es3", jutil.StringPtr("es6/iterator"))
	c.AddLibrary("Iterator.from", "es_next", "es3", jutil.StringPtr("es6/iterator/from"))

	c.SetLanguage(rewritePolyfillsTest_ES_2020, rewritePolyfillsTest_ES3)
	c.TestInjects("Iterator;", "es6/iterator")
	c.TestInjects("class MyIter extends Iterator {}", "es6/iterator")
	c.TestInjects("Iterator.from(x);", "es6/iterator", "es6/iterator/from")
	c.TestDoesNotInject("")
}

// Port of com/google/javascript/jscomp/IsolatePolyfillsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/IsolatePolyfillsTest.java)

/*
 * Copyright 2020 The Closure Compiler Authors.
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

type IsolatePolyfillsTest struct {
	CompilerTestCaseBase
	polyfillTable            []string
	polyfillsToInject        *jutil.LinkedHashSet[string]
	enablePropertyFlattening bool
}

// Java: IsolatePolyfillsTest()
func NewIsolatePolyfillsTest() *IsolatePolyfillsTest {
	c := &IsolatePolyfillsTest{}
	c.polyfillsToInject = jutil.NewLinkedHashSet[string]()
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestIsolatePolyfillsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewIsolatePolyfillsTest() })
}

// Java: IsolatePolyfillsTest.isolatePolyfillsTest_ES_2020
const isolatePolyfillsTest_ES_2020 = CompilerOptions_LanguageMode_ECMASCRIPT_2020

// Java: IsolatePolyfillsTest.isolatePolyfillsTest_ES6
const isolatePolyfillsTest_ES6 = CompilerOptions_LanguageMode_ECMASCRIPT_2015

// Java: IsolatePolyfillsTest.isolatePolyfillsTest_ES5
const isolatePolyfillsTest_ES5 = CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT

// Java: IsolatePolyfillsTest.isolatePolyfillsTest_ES3
const isolatePolyfillsTest_ES3 = CompilerOptions_LanguageMode_ECMASCRIPT3

// Java: IsolatePolyfillsTest.addLibrary(String, String, String, String)
func (c *IsolatePolyfillsTest) AddLibrary(name string, from string, to string, library *string) {
	// Go-only: test-helper — String.format("%s %s %s %s", ...) with
	// Strings.nullToEmpty(library).
	nullToEmpty := ""
	if library != nil {
		nullToEmpty = *library
	}
	c.polyfillTable = append(c.polyfillTable, name+" "+from+" "+to+" "+nullToEmpty)
	c.polyfillsToInject.Add(name)
}

// Java: IsolatePolyfillsTest.setUp()
func (c *IsolatePolyfillsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.polyfillTable = c.polyfillTable[:0]
	c.polyfillsToInject.Clear()
	c.AllowExternsChanges()
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
}

// Java: IsolatePolyfillsTest.getProcessor(Compiler)
func (c *IsolatePolyfillsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		// Synthetic definition of $jscomp$lookupPolyfilledValue
		compiler.
			GetSynthesizedExternsInput().
			GetAstRoot(compiler).
			AddChildToBack(rhino.IR_Var(rhino.IR_Name("$jscomp$lookupPolyfilledValue")))
		NewIsolatePolyfillsWithPolyfills(compiler, PolyfillUsageFinder_Polyfills_FromTable(strings.Join(c.polyfillTable, "\n"))).
			Process(externs, root)
	})
}

// Java: IsolatePolyfillsTest.getOptions()
func (c *IsolatePolyfillsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	if c.enablePropertyFlattening {
		options.SetCollapsePropertiesLevel(CompilerOptions_PropertyCollapseLevel_ALL)
	} else {
		options.SetCollapsePropertiesLevel(CompilerOptions_PropertyCollapseLevel_NONE)
	}
	return options
}

// BuildPolyfillJs builds a string with all definitions of the polyfills
//
// Java: IsolatePolyfillsTest.buildPolyfillJs()
func (c *IsolatePolyfillsTest) BuildPolyfillJs() string {
	if c.polyfillsToInject.IsEmpty() {
		return ""
	}
	jscompPolyfillName := "$jscomp.polyfill"
	if c.enablePropertyFlattening {
		jscompPolyfillName = "$jscomp$polyfill"
	}

	var syntheticCode strings.Builder
	syntheticCode.WriteString("var $jscomp = {};\n")

	for polyfill := range c.polyfillsToInject.All() {
		// $jscomp.polyfill('syntheticName');
		syntheticCode.WriteString(jscompPolyfillName)
		syntheticCode.WriteString("('")
		syntheticCode.WriteString(polyfill)
		syntheticCode.WriteString("');\n")
	}
	return syntheticCode.String()
}

// Java: IsolatePolyfillsTest.testWithPolyfills(String, String)
func (c *IsolatePolyfillsTest) TestWithPolyfills(input string, expected string) {
	c.Test(c.BuildPolyfillJs()+input, c.BuildPolyfillJs()+expected)
}

// Java: IsolatePolyfillsTest.testSameWithPolyfills(String)
func (c *IsolatePolyfillsTest) TestSameWithPolyfills(input string) {
	c.TestSame(c.BuildPolyfillJs() + input)
}

// Java: IsolatePolyfillsTest.testEmpty()
func (c *IsolatePolyfillsTest) TestEmpty() {
	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestSameWithPolyfills("")
}

// Java: IsolatePolyfillsTest.testOptChainCall_doesNotCrash()
func (c *IsolatePolyfillsTest) TestOptChainCall_doesNotCrash() {
	c.AddLibrary("String.prototype.includes", "es6", "es5", jutil.StringPtr("es6/string/includes"))
	c.AddLibrary("Array.prototype.includes", "es6", "es5", jutil.StringPtr("es6/array/includes"))
	c.Test(
		c.BuildPolyfillJs()+"if (a.b()?.includes()) {}",
		`var $jscomp$polyfillTmp;
`+
			c.BuildPolyfillJs()+
			`if (($jscomp$polyfillTmp = a.b(),
    $jscomp$lookupPolyfilledValue($jscomp$polyfillTmp, 'includes', true))?.
        call($jscomp$polyfillTmp)) {}
`)
}

// Java: IsolatePolyfillsTest.testOptChainGetProp_doesNotCrash()
func (c *IsolatePolyfillsTest) TestOptChainGetProp_doesNotCrash() {
	c.AddLibrary("String.prototype.includes", "es6", "es5", jutil.StringPtr("es6/string/includes"))
	c.AddLibrary("Array.prototype.includes", "es6", "es5", jutil.StringPtr("es6/array/includes"))
	c.TestWithPolyfills(
		"if (a.b?.includes()) {}",
		"if ($jscomp$lookupPolyfilledValue(a.b, 'includes', true)?.call(a.b)) {}")
}

// Java: IsolatePolyfillsTest.testClassesAreIsolated()
func (c *IsolatePolyfillsTest) TestClassesAreIsolated() {
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills("var m = new Map();", "var m = new $jscomp.polyfills['Map']();")
	c.TestWithPolyfills("var m = new window.Map();", "var m = new $jscomp.polyfills['Map']();")
	c.TestWithPolyfills("var m = new goog.global.Map();", "var m = new $jscomp.polyfills['Map']();")
}

// Java: IsolatePolyfillsTest.testClassesAreNotIsolatedUnlessPolyfillInjected()
func (c *IsolatePolyfillsTest) TestClassesAreNotIsolatedUnlessPolyfillInjected() {
	// Model a case where there is no actual $jscomp.polyfill('Map' call.
	// (Possibly because RemoveUnusedCode deleted it.)
	// Replacing "new Map()" with "new $jscomp.polyfills['Map']()" would result in an undefined
	// reference at runtime, so just leave "new Map()" as is.

	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))
	c.polyfillsToInject.Remove("Map")

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestSame("var m = new Map();")
	c.TestSame("var m = new window.Map();")
	c.TestSame("var m = new goog.global.Map();")
}

// Java: IsolatePolyfillsTest.testClassesGuardedByIfAreIsolated()
func (c *IsolatePolyfillsTest) TestClassesGuardedByIfAreIsolated() {
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills(
		"if (Map) { var m = new Map(); }",
		"if ($jscomp.polyfills['Map']) { var m = new $jscomp.polyfills['Map'](); }")
}

// Java: IsolatePolyfillsTest.testClassesAccessedViaBrackets_notIsolated()
func (c *IsolatePolyfillsTest) TestClassesAccessedViaBrackets_notIsolated() {
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	// We could support literal accesses like these but are currently choosing not to.
	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestSame("var m = new window['Map']();")
	c.TestSame("var m = new goog.global['Map']();")
}

// Java: IsolatePolyfillsTest.testGlobalThisIsolated()
func (c *IsolatePolyfillsTest) TestGlobalThisIsolated() {
	c.AddLibrary("globalThis", "es_2020", "es3", jutil.StringPtr("es6/globalthis"))
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills("alert(globalThis);", "alert($jscomp.polyfills['globalThis']);")

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES_2020)
	c.TestSameWithPolyfills("var m = new globalThis.Map();")
}

// Java: IsolatePolyfillsTest.testClassOnGlobalThisIsolated()
func (c *IsolatePolyfillsTest) TestClassOnGlobalThisIsolated() {
	c.AddLibrary("globalThis", "es_2020", "es3", jutil.StringPtr("es6/globalthis"))
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills("var m = new globalThis.Map();", "var m = new $jscomp.polyfills['Map']();")

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES6)
	c.TestWithPolyfills(
		"var m = new globalThis.Map();", "var m = new $jscomp.polyfills['globalThis'].Map();")
}

// Java: IsolatePolyfillsTest.testEnablingPropertyFlatteningUsesFlattenedJscompPolyfills()
func (c *IsolatePolyfillsTest) TestEnablingPropertyFlatteningUsesFlattenedJscompPolyfills() {
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.enablePropertyFlattening = true
	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills("var m = new Map();", "var m = new $jscomp$polyfills['Map']();")
	c.TestWithPolyfills("var m = new window.Map();", "var m = new $jscomp$polyfills['Map']();")
	c.TestWithPolyfills("var m = new goog.global.Map();", "var m = new $jscomp$polyfills['Map']();")
}

// Java: IsolatePolyfillsTest.testMultipleUsagesOfClassAllIsolated()
func (c *IsolatePolyfillsTest) TestMultipleUsagesOfClassAllIsolated() {
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills(
		"var m = new Map(); m = new Map();",
		"var m = new $jscomp.polyfills['Map'](); m = new $jscomp.polyfills['Map']();")
}

// Java: IsolatePolyfillsTest.testClassesGivenSufficientLanguageOutNotIsolated()
func (c *IsolatePolyfillsTest) TestClassesGivenSufficientLanguageOutNotIsolated() {
	c.AddLibrary("Proxy", "es6", "es6", nil)
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))
	c.AddLibrary("Set", "es6", "es3", jutil.StringPtr("es6/set"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES6)
	c.TestSameWithPolyfills("new Proxy();")
	c.TestSameWithPolyfills("var m = new Map();")
	c.TestSameWithPolyfills("new Set();")
}

// Java: IsolatePolyfillsTest.testClassesDeclaredInScopeNotIsolated()
func (c *IsolatePolyfillsTest) TestClassesDeclaredInScopeNotIsolated() {
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestSameWithPolyfills("/** @constructor */ var Map = function() {}; new Map();")
}

// Java: IsolatePolyfillsTest.testStaticMethodsIsolated()
func (c *IsolatePolyfillsTest) TestStaticMethodsIsolated() {
	c.AddLibrary("Array.of", "es6", "es3", jutil.StringPtr("es6/array/of"))
	c.AddLibrary("Object.keys", "es5", "es3", jutil.StringPtr("es5/object/keys"))
	c.AddLibrary("Math.clz32", "es6", "es5", jutil.StringPtr("es6/math/clz32"))

	// NOTE: Storing an entry for `Array.of` in the $jscomp.polyfills dictionary instead of using
	// `$jscomp$lookupPolyfilledValue` is a potential optimization, but not necessary for
	// correctness.
	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills("Array.of(x);", "$jscomp$lookupPolyfilledValue(Array, 'of').call(Array, x);")
	c.TestWithPolyfills(
		"Math.clz32(x);", "$jscomp$lookupPolyfilledValue(Math, 'clz32').call(Math, x);")

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES3)
	c.TestWithPolyfills(
		"Object.keys(x);", "$jscomp$lookupPolyfilledValue(Object, 'keys').call(Object, x);")
}

// Java: IsolatePolyfillsTest.testStaticMethodsWithSufficientLanguageOutNotIsolated()
func (c *IsolatePolyfillsTest) TestStaticMethodsWithSufficientLanguageOutNotIsolated() {
	c.AddLibrary("Array.from", "es6", "es6", nil)
	c.AddLibrary("Math.clz32", "es6", "es5", jutil.StringPtr("es6/math/clz32"))
	c.AddLibrary("Array.of", "es6", "es3", jutil.StringPtr("es6/array/of"))
	c.AddLibrary("Object.keys", "es5", "es3", jutil.StringPtr("es5/object/keys"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES6)
	c.TestSameWithPolyfills("Array.from(x);")
	c.TestSameWithPolyfills("Math.clz32(x);")
	c.TestSameWithPolyfills("Array.of(x);")

	c.SetLanguage(isolatePolyfillsTest_ES5, isolatePolyfillsTest_ES5)
	c.TestSameWithPolyfills("Object.keys(x);")
}

// Java: IsolatePolyfillsTest.testStaticMethodsDeclaredInScopeNotIsolated()
func (c *IsolatePolyfillsTest) TestStaticMethodsDeclaredInScopeNotIsolated() {
	c.AddLibrary("Math.clz32", "es6", "es5", jutil.StringPtr("es6/math/clz32"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestSameWithPolyfills("var Math = {clz32: function() {}}; Math.clz32(x);")
}

// Java: IsolatePolyfillsTest.testStaticMethodsGuardedByIfStillIsolated()
func (c *IsolatePolyfillsTest) TestStaticMethodsGuardedByIfStillIsolated() {
	c.AddLibrary("Array.of", "es6", "es5", jutil.StringPtr("es6/array/of"))

	c.TestWithPolyfills(
		"if (Array.of) { Array.of(); } else { Array.of(); }",
		`if ($jscomp$lookupPolyfilledValue(Array, 'of')) {
$jscomp$lookupPolyfilledValue(Array, 'of').call(Array);
} else {
  $jscomp$lookupPolyfilledValue(Array, 'of').call(Array);
}
`)
}

// Java: IsolatePolyfillsTest.testStaticMethodIsolatedOnPolyfilledNameIsolated()
func (c *IsolatePolyfillsTest) TestStaticMethodIsolatedOnPolyfilledNameIsolated() {
	c.AddLibrary("Promise", "es6", "es3", jutil.StringPtr("es6/promise/promise"))
	c.AddLibrary("Promise.allSettled", "es_2020", "es3", jutil.StringPtr("es6/promise/allSettled"))

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES3)
	c.TestWithPolyfills(
		"Promise.allSettled([p1, p2]);",
		`$jscomp$lookupPolyfilledValue($jscomp.polyfills['Promise'], 'allSettled')
   .call($jscomp.polyfills['Promise'], [p1, p2]);
`)

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES6)
	c.TestWithPolyfills(
		"Promise.allSettled([p1, p2]);",
		"$jscomp$lookupPolyfilledValue(Promise, 'allSettled').call(Promise, [p1, p2]);")
}

// Java: IsolatePolyfillsTest.testSinglePrototypeMethodIsIsolated()
func (c *IsolatePolyfillsTest) TestSinglePrototypeMethodIsIsolated() {
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills("x.endsWith(y);", "$jscomp$lookupPolyfilledValue(x, 'endsWith').call(x, y);")

	c.Test(
		c.BuildPolyfillJs()+"x().endsWith(y);",
		`var $jscomp$polyfillTmp;
`+
			c.BuildPolyfillJs()+
			`($jscomp$polyfillTmp = x(),
  $jscomp$lookupPolyfilledValue($jscomp$polyfillTmp, 'endsWith'))
      .call($jscomp$polyfillTmp, y);
`)
}

// Java: IsolatePolyfillsTest.testMultiplePrototypeMethodsAllIsolated()
func (c *IsolatePolyfillsTest) TestMultiplePrototypeMethodsAllIsolated() {
	c.AddLibrary("String.prototype.startsWith", "es6", "es5", jutil.StringPtr("es6/string/startswith"))
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills(
		"x.endsWith(y) || x.startsWith(y);",
		`$jscomp$lookupPolyfilledValue(x, 'endsWith').call(x, y) ||
  $jscomp$lookupPolyfilledValue(x, 'startsWith').call(x, y);
`)

	c.Test(
		c.BuildPolyfillJs()+"x().endsWith(y) || x().startsWith(y);",
		"var $jscomp$polyfillTmp;"+
			c.BuildPolyfillJs()+
			`($jscomp$polyfillTmp = x(),
    $jscomp$lookupPolyfilledValue($jscomp$polyfillTmp, 'endsWith'))
  .call($jscomp$polyfillTmp, y) ||
 ($jscomp$polyfillTmp = x(),
      $jscomp$lookupPolyfilledValue($jscomp$polyfillTmp, 'startsWith'))
    .call($jscomp$polyfillTmp, y);
`)
}

// Java: IsolatePolyfillsTest.testMethodMatchingMultiplePolyfillsIsolatedOnlyOnce()
func (c *IsolatePolyfillsTest) TestMethodMatchingMultiplePolyfillsIsolatedOnlyOnce() {
	c.AddLibrary("String.prototype.includes", "es6", "es5", jutil.StringPtr("es6/string/includes"))
	c.AddLibrary("Array.prototype.includes", "es6", "es5", jutil.StringPtr("es6/array/includes"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills("x.includes(y);", "$jscomp$lookupPolyfilledValue(x, 'includes').call(x, y);")
}

// Java: IsolatePolyfillsTest.testMethodsNotIsolatedUnlessPolyfillInjected()
func (c *IsolatePolyfillsTest) TestMethodsNotIsolatedUnlessPolyfillInjected() {
	c.AddLibrary("String.prototype.includes", "es6", "es5", jutil.StringPtr("es6/string/includes"))
	c.AddLibrary("String.prototype.endsWith", "es6", "es5", jutil.StringPtr("es6/string/endswith"))
	c.polyfillsToInject.Remove("String.prototype.endsWith")

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)
	c.TestWithPolyfills("x.endsWith(y);", "x.endsWith(y);")
	c.TestWithPolyfills(
		"x.includes(y) && x.endsWith(z);",
		"$jscomp$lookupPolyfilledValue(x, 'includes').call(x, y) && x.endsWith(z);")
}

// Java: IsolatePolyfillsTest.testLvaluePolyfillUsageNotIsolated()
func (c *IsolatePolyfillsTest) TestLvaluePolyfillUsageNotIsolated() {
	c.AddLibrary("String.prototype.includes", "es6", "es5", jutil.StringPtr("es6/string/includes"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)

	c.TestSameWithPolyfills("x.includes = true;")
}

// Java: IsolatePolyfillsTest.testLvaluePolyfillUsage_doesntCauseBackoffForOtherUsages()
func (c *IsolatePolyfillsTest) TestLvaluePolyfillUsage_doesntCauseBackoffForOtherUsages() {
	c.AddLibrary("String.prototype.includes", "es6", "es5", jutil.StringPtr("es6/string/includes"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)

	c.TestWithPolyfills(
		"x.includes = true; y.includes(z)",
		"x.includes = true; $jscomp$lookupPolyfilledValue(y, 'includes').call(y, z)")
}

// Java: IsolatePolyfillsTest.testLvaluePolyfillUsageInNestedAssignNotIsolated()
func (c *IsolatePolyfillsTest) TestLvaluePolyfillUsageInNestedAssignNotIsolated() {
	c.AddLibrary("String.prototype.includes", "es6", "es5", jutil.StringPtr("es6/string/includes"))

	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)

	c.TestSameWithPolyfills("y = x.includes = true")
	c.TestSameWithPolyfills("y = x.includes = z.includes = true")
}

// Java: IsolatePolyfillsTest.testPolyfilledMethodOnPolyfilledClassAreBothIsolated()
func (c *IsolatePolyfillsTest) TestPolyfilledMethodOnPolyfilledClassAreBothIsolated() {
	c.AddLibrary("Promise", "es6", "es3", jutil.StringPtr("es6/promise/promise"))
	c.AddLibrary("Promise.prototype.finally", "es9", "es3", jutil.StringPtr("es6/promise/finally"))

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES3)

	c.TestWithPolyfills("p.finally(cb);", "$jscomp$lookupPolyfilledValue(p, 'finally').call(p, cb);")

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES6)
	c.TestWithPolyfills("p.finally(cb);", "$jscomp$lookupPolyfilledValue(p, 'finally').call(p, cb);")
}

// Java: IsolatePolyfillsTest.testNonPolyfilledMethodOnPolyfilledClassNotIsolated()
func (c *IsolatePolyfillsTest) TestNonPolyfilledMethodOnPolyfilledClassNotIsolated() {
	c.AddLibrary("Promise", "es6", "es3", jutil.StringPtr("es6/promise/promise"))

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES3)
	c.TestWithPolyfills("Promise.all(p1, p2);", "$jscomp.polyfills['Promise'].all(p1, p2);")
}

// Java: IsolatePolyfillsTest.testExternForJSCompLookupPolyfillDeleted()
func (c *IsolatePolyfillsTest) TestExternForJSCompLookupPolyfillDeleted() {
	c.DisableCompareAsTree()
	// The getProcessor() method injects $jscomp$lookupPolyfilledValue into the input externs.
	c.TestExternChanges(c.Externs(""), c.Srcs("x.includes(y);"), c.Expected(""))
}

// Java: IsolatePolyfillsTest.testJscompLookupPolyfillDeletedIfNotUsed()
func (c *IsolatePolyfillsTest) TestJscompLookupPolyfillDeletedIfNotUsed() {
	c.AddLibrary("Promise", "es6", "es3", jutil.StringPtr("es6/promise/promise"))
	c.AddLibrary("Promise.prototype.finally", "es9", "es3", jutil.StringPtr("es6/promise/promise"))

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES3)
	c.TestWithPolyfills("var $jscomp$lookupPolyfilledValue = function() {};", "")
	c.TestWithPolyfills(
		`var $jscomp$lookupPolyfilledValue = function() {};
new Promise();
`,
		"new $jscomp.polyfills['Promise']();")
}

// Java: IsolatePolyfillsTest.testJscompLookupPolyfillKeptIfUsed()
func (c *IsolatePolyfillsTest) TestJscompLookupPolyfillKeptIfUsed() {
	c.AddLibrary("Promise", "es6", "es3", jutil.StringPtr("es6/promise/promise"))
	c.AddLibrary("Promise.prototype.finally", "es9", "es3", jutil.StringPtr("es6/promise/promise"))

	c.SetLanguage(isolatePolyfillsTest_ES_2020, isolatePolyfillsTest_ES3)
	c.Test(
		c.BuildPolyfillJs()+
			`var $jscomp$lookupPolyfilledValue = function() {};
new Promise().finally(cb);
`,
		`var $jscomp$polyfillTmp;
`+
			c.BuildPolyfillJs()+
			`var $jscomp$lookupPolyfilledValue = function() {};
($jscomp$polyfillTmp = new $jscomp.polyfills['Promise'],
    $jscomp$lookupPolyfilledValue($jscomp$polyfillTmp, 'finally'))
  .call($jscomp$polyfillTmp, cb);
`)
}

// Java: IsolatePolyfillsTest.testPolyfilTemp_multipleFiles_insertedInFirstFile()
func (c *IsolatePolyfillsTest) TestPolyfilTemp_multipleFiles_insertedInFirstFile() {
	// Regression test for failure to report code change when adding "var $jscomp$polyfillTmp;"
	// Don't call addLibrary() for this test: it reports a code change regardless & make the
	// test incorrectly pass.
	c.polyfillTable = append(c.polyfillTable, "String.prototype.endsWith es6 es5 es6/string/endswith")
	c.SetLanguage(isolatePolyfillsTest_ES6, isolatePolyfillsTest_ES5)

	c.Test(
		c.Srcs("$jscomp.polyfill('String.prototype.endsWith', 'es6', 'es5');", "x().endsWith(y);"),
		c.Expected(
			`var $jscomp$polyfillTmp;
$jscomp.polyfill('String.prototype.endsWith', 'es6', 'es5');
`,
			`($jscomp$polyfillTmp = x(),
  $jscomp$lookupPolyfilledValue($jscomp$polyfillTmp, 'endsWith'))
      .call($jscomp$polyfillTmp, y);
`))
}

// Java: IsolatePolyfillsTest.testOptChainMethodCallOnPlainGetProp()
func (c *IsolatePolyfillsTest) TestOptChainMethodCallOnPlainGetProp() {
	// TODO(b/571890155): silent miscompile for x.includes?.(y)
	c.AddLibrary("String.prototype.includes", "es6", "es5", jutil.StringPtr("es6/string/includes"))
	c.TestWithPolyfills(
		"x.includes?.(y);", "$jscomp$lookupPolyfilledValue(x, 'includes').call?.(x, y);")
}

// Java: IsolatePolyfillsTest.testOptChainNonCallContinuation()
func (c *IsolatePolyfillsTest) TestOptChainNonCallContinuation() {
	// TODO(b/571890155): semantic bug caught by AstValidator.validateFirstNodeOfOptChain
	c.AddLibrary("String.prototype.includes", "es6", "es5", jutil.StringPtr("es6/string/includes"))
	e :=
		AssertThrows[jutil.IRuntimeException](func() { c.TestSameWithPolyfills("x?.includes.bind(x);") })
	// Throwables.getRootCause(e)
	var root jutil.Throwable = e
	for root.GetCause() != nil {
		root = root.GetCause()
	}
	IsInstanceOf[jutil.IIllegalStateException](AssertThat(root))
	AssertThat(root.GetMessage()).
		Contains("Start of optional chain node OPTCHAIN_GETPROP is not marked as the start")
}

// Java: IsolatePolyfillsTest.testOptChainOnGlobalPolyfill()
func (c *IsolatePolyfillsTest) TestOptChainOnGlobalPolyfill() {
	// TODO(b/571890155): replaces window?.Map with $jscomp.polyfills['Map'] instead of preserving ?.
	c.AddLibrary("Map", "es6", "es5", jutil.StringPtr("es6/map"))
	c.TestWithPolyfills("alert(window?.Map);", "alert($jscomp.polyfills['Map']);")
}

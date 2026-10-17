// Port of com/google/javascript/jscomp/CollapsePropertiesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/CollapsePropertiesTest.java)

/*
 * Copyright 2006 The Closure Compiler Authors.
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

// CollapsePropertiesTest tests InlineAndCollapseProperties.CollapseProperties.
type CollapsePropertiesTest struct {
	CompilerTestCaseBase

	propertyCollapseLevel CompilerOptions_PropertyCollapseLevel
}

// Java: CollapsePropertiesTest.EXTERNS
const CollapsePropertiesTest_EXTERNS = `var window;
function alert(s) {}
function parseInt(s) {}
/** @constructor */ function String() {};
var arguments
`

// Java: CollapsePropertiesTest()
func NewCollapsePropertiesTest() *CollapsePropertiesTest {
	c := &CollapsePropertiesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CollapsePropertiesTest_EXTERNS)
	c.propertyCollapseLevel = CompilerOptions_PropertyCollapseLevel_ALL
	return c
}

func TestCollapsePropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCollapsePropertiesTest() })
}

// SetUp runs CompilerTestCase's @Before method, then CollapsePropertiesTest's
// (JUnit runs superclass @Before methods first).
func (c *CollapsePropertiesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.CustomSetUp()
}

// Java: CollapsePropertiesTest.getProcessor(Compiler)
func (c *CollapsePropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("es6NormalizeClasses", func(comp AbstractCompiler) CompilerPass { return NewEs6NormalizeClasses(comp) }))
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"inlineAndCollapseProperties",
			func(comp AbstractCompiler) CompilerPass {
				return InlineAndCollapseProperties_Builder_(comp).
					SetPropertyCollapseLevel(c.propertyCollapseLevel).
					SetChunkOutputType(CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE).
					SetHaveModulesBeenRewritten(false).
					SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER).
					Build()
			}))
	return optimizer
}

// Java: CollapsePropertiesTest.customSetUp()
func (c *CollapsePropertiesTest) CustomSetUp() {
	c.EnableNormalize()
	c.DisableCompareJsDoc()
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
}

// Java: CollapsePropertiesTest.setupModuleExportsOnly()
func (c *CollapsePropertiesTest) SetupModuleExportsOnly() {
	c.EnableProcessCommonJsModules()
	c.EnableTranspile()
	c.propertyCollapseLevel = CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: CollapsePropertiesTest.doNotCollapseAStaticPropertyUsedInAGetter()
func (c *CollapsePropertiesTest) TestDoNotCollapseAStaticPropertyUsedInAGetter() {
	c.Test(
		c.Srcs(
			`class C {
  /** @type {number} */
  static get p() { return this.p_; }
}
// declaration and initialization of private static field
/** @private {number} */
C.p_ = 1;
// reference via getter
alert(C.p);
`),
		c.Expected(
			`class C {
  /** @type {number} */
  static get p() { return this.p_; }
}
// TODO(b/130682799): C.p_ should not be collapsed
/** @private {number} */
var C$p_ = 1;
alert(C.p);
`))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: CollapsePropertiesTest.doNotCollapseAStaticPropertyUsedInASetter()
func (c *CollapsePropertiesTest) TestDoNotCollapseAStaticPropertyUsedInASetter() {
	c.Test(
		c.Srcs(
			`class C {
  /** @param {number} v */
  static set p(v) { this.p_ = v; }
}
// declaration and initialization of private static field
/** @private {number} */
C.p_ = 0;
// changing private static field through setter property
C.p = 1;
alert(C.p_);
`),
		c.Expected(
			`class C {
  /** @param {number} v */
  static set p(v) { this.p_ = v; }
}
// TODO(b/130682799): C.p_ should not be collapsed
/** @private {number} */
var C$p_ = 0;
C.p = 1;
alert(C$p_);
`))
	c.Test(
		c.Srcs(
			`class C {
  /** @param {number} v */
  static set p(v) { this.p_ = v; }
}
// declaration and initialization of private static field
/** @private {number} */
C.p_ = 0;
// changing private static field through setter property
C.p = 1;
alert?.(C?.p_);
`),
		c.Expected(
			`class C {
  /** @param {number} v */
  static set p(v) { this.p_ = v; }
}
// TODO(b/130682799): C.p_ should not be collapsed
/** @private {number} */
var C$p_ = 0;
C.p = 1;
alert?.(C?.p_);
`))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: CollapsePropertiesTest.doNotCollapseAStaticPropertyReadInAStaticMethod()
func (c *CollapsePropertiesTest) TestDoNotCollapseAStaticPropertyReadInAStaticMethod() {
	c.TestWarning(
		`class C {
  /** @return {number} */
  static getP() { return this.p_; }
}
// declaration and initialization of private static field
/** @private {number} */
C.p_ = 1;
// reference via method
alert(C.getP());
`,
		// TODO(b/117437011): should recognize type of `this` in a static method
		InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: CollapsePropertiesTest.doNotCollapseAStaticPropertyAssignedInAStaticMethod()
func (c *CollapsePropertiesTest) TestDoNotCollapseAStaticPropertyAssignedInAStaticMethod() {
	c.TestWarning(
		`class C {
  /** @param {number} v */
  static setP(v) { this.p_ = v; }
}
// declaration and initialization of private static field
/** @private {number} */
C.p_ = 0;
// changing private static field through method
C.setP(1);
alert(C.p_);
`,
		// TODO(b/117437011): should recognize type of `this` in a static method
		InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: CollapsePropertiesTest.doNotCollapseAPropertyReadByAnObjectLiteralMethod()
func (c *CollapsePropertiesTest) TestDoNotCollapseAPropertyReadByAnObjectLiteralMethod() {
	c.Test(
		c.Srcs(
			`const obj = {
  /** @private {number} */
  p_: 0,
  /**
   * @this {{p_: number}}
   * @return {number}
   */
  getP() { return this.p_; }
};
alert(obj.getP());
`),
		c.Expected(
			`var obj$p_ = 0;
var obj$getP = function() { return this.p_; };
alert(obj$getP());
`))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: CollapsePropertiesTest.doNotCollapseAPropertyAssignedInAnObjectLiteralMethod()
func (c *CollapsePropertiesTest) TestDoNotCollapseAPropertyAssignedInAnObjectLiteralMethod() {
	c.Test(
		c.Srcs(
			`const obj = {
  /** @private {number} */
  p_: 0,
  /**
   * @this {{p_: number}}
   * @param {number} v
   */
  setP(v) { this.p_ = v; }
}
obj.setP(1);
alert(obj.p_);
`),
		c.Expected(
			`var obj$p_ = 0;
var obj$setP = function(v) { this.p_ = v; };
obj$setP(1);
alert(obj$p_);
`))
	c.TestSame(
		`const obj = {
  /** @private {number} */
  p_: 0,
  /**
   * @this {{p_: number}}
   * @param {number} v
   */
  setP(v) { this.p_ = v; }
}
obj.setP?.(1);
alert(obj?.p_);
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: CollapsePropertiesTest.doNotCollapseAPropertyUsedInAnObjectLiteralGetter()
func (c *CollapsePropertiesTest) TestDoNotCollapseAPropertyUsedInAnObjectLiteralGetter() {
	c.Test(
		c.Srcs(
			`const obj = {
  /** @private {number} */
  p_: 0,
  /** @type {number} */
  get p() { return this.p_; }
}
alert(obj.p);
`),
		c.Expected(
			`var obj$p_ = 0;
const obj = {
  /** @type {number} */
  get p() { return this.p_; }
}
alert(obj.p);
`))
	c.TestSame(
		`const obj = {
  /** @private {number} */
  p_: 0,
  /** @type {number} */
  get p() { return this?.p_; }
}
alert(obj?.p);
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: CollapsePropertiesTest.doNotCollapseAPropertyUsedInAnObjectLiteralSetter()
func (c *CollapsePropertiesTest) TestDoNotCollapseAPropertyUsedInAnObjectLiteralSetter() {
	c.Test(
		c.Srcs(
			`const obj = {
  /** @private {number} */
  p_: 0,
  /** @param {number} v */
  set p(v) { this.p_ = v; }
}
obj.p = 1;
alert(obj.p_);
`),
		c.Expected(
			`var obj$p_ = 0;
const obj = {
  /** @param {number} v */
  set p(v) { this.p_ = v; }
}
obj.p = 1;
alert(obj$p_);
`))
	c.TestSame(
		`const obj = {
  /** @private {number} */
  p_: 0,
  /** @param {number} v */
  set p(v) { this.p_ = v; }
}
obj.p = 1;
alert?.(obj?.p_);
`)
}

// Java: CollapsePropertiesTest.testMultiLevelCollapse()
func (c *CollapsePropertiesTest) TestMultiLevelCollapse() {
	c.Test(
		c.Srcs("var a = {}; a.b = {}; a.b.c = {}; var d = 1; d = a.b.c;"),
		c.Expected("var                   a$b$c = {}; var d = 1; d = a$b$c;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {}; a.b.c = {}; var d = 1; d = a.b?.c;"),
		c.Expected("var         a$b = {}; a$b.c = {}; var d = 1; d = a$b?.c;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {}; /** @nocollapse */ a.b.c = {}; var d = 1; d = a.b.c;"),
		c.Expected("var a$b = {}; /** @nocollapse */ a$b.c = {}; var d = 1; d = a$b.c;"))
}

// Java: CollapsePropertiesTest.testDecrement()
func (c *CollapsePropertiesTest) TestDecrement() {
	c.Test(c.Srcs("var a = {}; a.b = 5; a.b--; a.b = 5"), c.Expected("var a$b = 5; a$b--; a$b = 5"))
}

// Java: CollapsePropertiesTest.testIncrement()
func (c *CollapsePropertiesTest) TestIncrement() {
	c.Test(c.Srcs("var a = {}; a.b = 5; a.b++; a.b = 5"), c.Expected("var a$b = 5; a$b++; a$b = 5"))
}

// Java: CollapsePropertiesTest.testObjLitDeclarationWithGet1()
func (c *CollapsePropertiesTest) TestObjLitDeclarationWithGet1() {
	c.TestSame("var a = {get b(){}};")
}

// Java: CollapsePropertiesTest.testObjLitDeclarationWithGet3()
func (c *CollapsePropertiesTest) TestObjLitDeclarationWithGet3() {
	c.Test(
		c.Srcs("var a = {b: {get c() { return 3; }}};"),
		c.Expected("var a$b = {get c() { return 3; }};"))
}

// Java: CollapsePropertiesTest.testObjLitDeclarationWithSet1()
func (c *CollapsePropertiesTest) TestObjLitDeclarationWithSet1() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame("var a = {set b(a){}};")
}

// Java: CollapsePropertiesTest.testObjLitDeclarationWithSet3()
func (c *CollapsePropertiesTest) TestObjLitDeclarationWithSet3() {
	c.Test(c.Srcs("var a = {b: {set c(d) {}}};"), c.Expected("var a$b = {set c(d) {}};"))
}

// Java: CollapsePropertiesTest.testObjLitDeclarationWithUsedSetter()
func (c *CollapsePropertiesTest) TestObjLitDeclarationWithUsedSetter() {
	c.TestSame("var a = {set b(c) {}}; a.b = 4;")
}

// Java: CollapsePropertiesTest.testObjLitDeclarationDoesntCollapsePropertiesOnGetter()
func (c *CollapsePropertiesTest) TestObjLitDeclarationDoesntCollapsePropertiesOnGetter() {
	c.Test(
		`var a = {
  get b() { return class {}; },
  set b(c) {}
};
a.b = class {};
a.b.c = 4;
`,
		`var a = {
  get b() {
    const CLASS_DECL$0 = class {};
    return CLASS_DECL$0;
  },
  set b(c) {}
};
a.b = class {};
a.b.c = 4;
`)
}

// Java: CollapsePropertiesTest.testObjLitDeclarationWithGetAndSet1()
func (c *CollapsePropertiesTest) TestObjLitDeclarationWithGetAndSet1() {
	c.Test(
		c.Srcs("var a = {b: {get c() { return 3; },set c(d) {}}};"),
		c.Expected("var a$b = {get c() { return 3; },set c(d) {}};"))
}

// Java: CollapsePropertiesTest.testObjLitAssignmentDepth1()
func (c *CollapsePropertiesTest) TestObjLitAssignmentDepth1() {
	c.Test(
		c.Srcs("var a = {b: {}, c: {}}; var d = 1; var e = 1; d = a.b; e = a.c"),
		c.Expected("var a$b = {}; var a$c = {}; var d = 1; var e = 1; d = a$b; e = a$c"))
	c.TestSame("var a = {b: {}, c: {}}; var d = 1; var e = 1; d = a?.b; e = a?.c")
	c.Test(
		c.Srcs(
			`var a = {b: {}, /** @nocollapse */ c: {}};
var d = 1;
d = a.b;
var e = 1;
e = a.c;
`),
		c.Expected(
			`var a$b = {};
var a = { /** @nocollapse */c: {}};
var d = 1;
d = a$b;
var e = 1;
e = a.c;
`))
	c.TestSame("var a = {b: {}, /** @nocollapse */ c: {}}; var d = 1; d = a?.b; var e = 1; e = a?.c")
}

// Java: CollapsePropertiesTest.testObjLitAssignmentDepth2()
func (c *CollapsePropertiesTest) TestObjLitAssignmentDepth2() {
	c.Test(
		c.Srcs("var a = {}; a.b = {c: {}, d: {}}; var e = 1; e = a.b.c; var f = 1; f = a.b.d"),
		c.Expected("var a$b$c = {}; var a$b$d = {}; var e = 1; e = a$b$c; var f = 1; f = a$b$d;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {c: {}, d: {}}; var e = 1; e = a.b?.c; var f = 1; f = a.b?.d"),
		c.Expected(
			"var         a$b = {c: {}, d: {}}; var e = 1; e = a$b?.c; var f = 1; f = a$b?.d;"))
	c.Test(
		c.Srcs(
			`var a = {}; a.b = {c: {}, /** @nocollapse */ d: {}}; var e = 1; e = a.b.c;
var f = 1; f = a.b.d
`),
		c.Expected(
			`var a$b$c = {}; var a$b = { /** @nocollapse */ d: {}}; var e = 1; e = a$b$c;
var f = 1; f = a$b.d;
`))
	c.Test(
		c.Srcs(
			`var a = {}; a.b = {c: {}, /** @nocollapse */ d: {}}; var e = 1; e = a.b?.c;
var f = 1; f = a.b?.d
`),
		c.Expected(
			`var         a$b = {c: {}, /** @nocollapse */ d: {}}; var e = 1; e = a$b?.c;
var f = 1; f = a$b?.d;
`))
}

// Java: CollapsePropertiesTest.testGlobalObjectDeclaredToPreserveItsPreviousValue1()
func (c *CollapsePropertiesTest) TestGlobalObjectDeclaredToPreserveItsPreviousValue1() {
	c.Test(c.Srcs("var a = a ? a : {}; a.c = 1;"), c.Expected("var a = a ? a : {}; var a$c = 1;"))
	c.Test(
		c.Srcs("var a = a ? a : {}; /** @nocollapse */ a.c = 1;"),
		c.Expected("var a = a ? a : {}; /** @nocollapse */ a.c = 1;"))
}

// Java: CollapsePropertiesTest.testGlobalObjectDeclaredToPreserveItsPreviousValue2()
func (c *CollapsePropertiesTest) TestGlobalObjectDeclaredToPreserveItsPreviousValue2() {
	c.Test(c.Srcs("var a = a || {}; a.c = 1;"), c.Expected("var a = a || {}; var a$c = 1;"))
	c.TestSame("var a = a || {}; /** @nocollapse */ a.c = 1;")
}

// Java: CollapsePropertiesTest.testGlobalObjectDeclaredToPreserveItsPreviousValue3()
func (c *CollapsePropertiesTest) TestGlobalObjectDeclaredToPreserveItsPreviousValue3() {
	c.Test(
		c.Srcs("var a = a || {get b() {}}; a.c = 1;"),
		c.Expected("var a = a || {get b() {}}; var a$c = 1;"))
	c.TestSame("var a = a || {get b() {}}; /** @nocollapse */ a.c = 1;")
}

// Java: CollapsePropertiesTest.testGlobalObjectNameInBooleanExpressionDepth1_1()
func (c *CollapsePropertiesTest) TestGlobalObjectNameInBooleanExpressionDepth1_1() {
	c.Test(
		c.Srcs("var a = {b: 0}; a.c = 1; if (a) x();"),
		c.Expected("var a$b = 0; var a = {}; var a$c = 1; if (a) x();"))
	c.Test(
		c.Srcs("var a = {/** @nocollapse */ b: 0}; a.c = 1; if (a) x();"),
		c.Expected("var a = {/** @nocollapse */ b: 0}; var a$c = 1; if (a) x();"))
	c.Test(
		c.Srcs("var a = {b: 0}; /** @nocollapse */ a.c = 1; if (a) x();"),
		c.Expected("var a$b = 0; var a = {}; /** @nocollapse */ a.c = 1; if (a) x();"))
}

// Java: CollapsePropertiesTest.testGlobalObjectNameInBooleanExpressionDepth1_2()
func (c *CollapsePropertiesTest) TestGlobalObjectNameInBooleanExpressionDepth1_2() {
	c.Test(
		c.Srcs("var a = {b: 0}; a.c = 1; if (!(a && a.c)) x();"),
		c.Expected("var a$b = 0; var a = {}; var a$c = 1; if (!(a && a$c)) x();"))
	c.TestSame("var a = {b: 0}; a.c = 1; if (!(a && a?.c)) x?.();")
	c.TestSame("var a = {/** @nocollapse */ b: 0}; a.c = 1; if (!(a && a?.c)) x?.();")
	c.Test(
		c.Srcs("var a = {b: 0}; /** @nocollapse */ a.c = 1; if (!(a && a.c)) x();"),
		c.Expected("var a$b = 0; var a = {}; /** @nocollapse */ a.c = 1; if (!(a && a.c)) x();"))
	c.TestSame("var a = {b: 0}; /** @nocollapse */ a.c = 1; if (!(a && a?.c)) x?.();")
}

// Java: CollapsePropertiesTest.testGlobalObjectNameInBooleanExpressionDepth1_3()
func (c *CollapsePropertiesTest) TestGlobalObjectNameInBooleanExpressionDepth1_3() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs("var a = {b: 0}; a.c = 1; while (a || a.c) x();"),
		c.Expected("var a$b = 0; var a = {}; var a$c = 1; while (a || a$c) x();"))
	c.TestSame("var a = {b: 0}; a.c = 1; while (a || a?.c) x?.();")
	c.Test(
		c.Srcs("var a = {/** @nocollapse */ b: 0}; a.c = 1; while (a || a.c) x();"),
		c.Expected("var a = {/** @nocollapse */ b: 0}; var a$c = 1; while (a || a$c) x();"))
	c.TestSame("var a = {/** @nocollapse */ b: 0}; a.c = 1; while (a || a?.c) x?.();")
	c.Test(
		c.Srcs("var a = {b: 0}; /** @nocollapse */ a.c = 1; while (a || a.c) x();"),
		c.Expected("var a$b = 0; var a = {}; /** @nocollapse */ a.c = 1; while (a || a.c) x();"))
	c.TestSame("var a = {b: 0}; /** @nocollapse */ a.c = 1; while (a || a?.c) x?.();")
}

// Java: CollapsePropertiesTest.testGlobalObjectNameInBooleanExpressionDepth1_4()
func (c *CollapsePropertiesTest) TestGlobalObjectNameInBooleanExpressionDepth1_4() {
	c.TestSame("var a = {}; a.c = 1; var d = a || {}; a.c;")
	c.TestSame("var a = {}; /** @nocollapse */ a.c = 1; var d = a || {}; a.c;")
}

// Java: CollapsePropertiesTest.testGlobalObjectNameInBooleanExpressionDepth1_5()
func (c *CollapsePropertiesTest) TestGlobalObjectNameInBooleanExpressionDepth1_5() {
	c.TestSame("var a = {}; a.c = 1; var d = a.c || a; a.c;")
	c.TestSame("var a = {}; /** @nocollapse */ a.c = 1; var d = a.c || a; a.c;")
}

// Java: CollapsePropertiesTest.testGlobalObjectNameInBooleanExpressionDepth1_6()
func (c *CollapsePropertiesTest) TestGlobalObjectNameInBooleanExpressionDepth1_6() {
	c.Test(
		c.Srcs("var a = {b: 0}; a.c = 1; var d = !(a.c || a); a.c;"),
		c.Expected("var a$b = 0; var a = {}; var a$c = 1; var d = !(a$c || a); a$c;"))
	c.Test(
		c.Srcs("var a = {/** @nocollapse */ b: 0}; a.c = 1; var d = !(a.c || a); a.c;"),
		c.Expected("var a = {/** @nocollapse */ b: 0}; var a$c = 1; var d = !(a$c || a); a$c;"))
	c.Test(
		c.Srcs("var a = {b: 0}; /** @nocollapse */ a.c = 1; var d = !(a.c || a); a.c;"),
		c.Expected("var a$b = 0; var a = {}; /** @nocollapse */ a.c = 1; var d = !(a.c || a); a.c;"))
}

// Java: CollapsePropertiesTest.testGlobalObjectNameInBooleanExpressionDepth2()
func (c *CollapsePropertiesTest) TestGlobalObjectNameInBooleanExpressionDepth2() {
	c.Test(
		c.Srcs("var a = {b: {}}; a.b.c = 1; if (a.b) x(a.b.c);"),
		c.Expected("var a$b = {}; var a$b$c = 1; if (a$b) x(a$b$c);"))
	c.TestSame("var a = {b: {}}; a.b.c = 1; if (a?.b) x?.(a.b?.c);")
	c.TestSame("var a = {/** @nocollapse */ b: {}}; a.b.c = 1; if (a.b) x(a.b.c);")
	c.Test(
		c.Srcs("var a = {b: {}}; /** @nocollapse */ a.b.c = 1; if (a.b) x(a.b.c);"),
		c.Expected("var a$b = {}; /** @nocollapse */ a$b.c = 1; if (a$b) x(a$b.c);"))
	c.TestSame("var a = {b: {}}; /** @nocollapse */ a.b.c = 1; if (a?.b) x?.(a.b?.c);")
}

// Java: CollapsePropertiesTest.testGlobalObjectNameInBooleanExpressionDepth3()
func (c *CollapsePropertiesTest) TestGlobalObjectNameInBooleanExpressionDepth3() {
	// TODO(user): Make CollapseProperties even more aggressive so that
	// a$b.z gets collapsed. Right now, it doesn't get collapsed because the
	// expression (a.b && a.b.c) could return a.b. But since it returns a.b iff
	// a.b *is* safely collapsible, the Boolean logic should be smart enough to
	// only consider the right side of the && as aliasing.
	c.Test(
		`var a = {}; a.b = {}; /** @constructor */ a.b.c = function(){};
 a.b.z = 1; var d = a.b && a.b.c;
`,
		`var a$b = {}; /** @constructor */ a$b.c = function(){};
 a$b.z = 1; var d = a$b && a$b.c;
`,
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: CollapsePropertiesTest.testGlobalFunctionNameInBooleanExpressionDepth1()
func (c *CollapsePropertiesTest) TestGlobalFunctionNameInBooleanExpressionDepth1() {
	c.Test(
		c.Srcs("function a() {}     a.c = 1; if (a) x(a.c);"),
		c.Expected("function a() {} var a$c = 1; if (a) x(a$c);"))
	c.TestSame("function a() {} a.c = 1; if (a) x?.(a?.c);")
	c.Test(
		c.Srcs("function a() {} /** @nocollapse */ a.c = 1; if (a) x(a.c);"),
		c.Expected("function a() {} /** @nocollapse */ a.c = 1; if (a) x(a.c);"))
}

// Java: CollapsePropertiesTest.testGlobalFunctionNameInBooleanExpressionDepth2()
func (c *CollapsePropertiesTest) TestGlobalFunctionNameInBooleanExpressionDepth2() {
	c.Test(
		c.Srcs("var a = {b: function(){}}; a.b.c = 1; if (a.b) x(a.b.c);"),
		c.Expected("var a$b = function(){}; var a$b$c = 1; if (a$b) x(a$b$c);"))
	c.TestSame("var a = {b: function(){}}; a.b.c = 1; if (a?.b) x?.(a.b?.c);")
	c.TestSame(
		`var a = {/** @nocollapse */ b: function(){}}; a.b.c = 1;
if (a.b) x(a.b.c);
`)
	c.Test(
		c.Srcs(
			`var a = {b: function(){}}; /** @nocollapse */ a.b.c = 1;
if (a.b) x(a.b.c);
`),
		c.Expected("var a$b = function(){}; /** @nocollapse */ a$b.c = 1; if (a$b) x(a$b.c);"))
	c.TestSame(
		`var a = {b: function(){}};
/** @nocollapse */ a.b.c = 1;
if (a?.b) x?.(a.b?.c);
`)
}

// Java: CollapsePropertiesTest.testDontCollapseObjectLiteralVarDeclarationInsideLoop()
func (c *CollapsePropertiesTest) TestDontCollapseObjectLiteralVarDeclarationInsideLoop() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	// See https://github.com/google/closure-compiler/issues/3050
	// Another solution to the issue would be to explicitly initialize obj.val to undefined at the
	// start of the loop, but that requires some more refactoring of CollapseProperties.
	c.TestSame(
		`for (var i = 0; i < 2; i++) {
  var obj = {};
  if (i == 0) {
    obj.val = 1;
  }
  alert(obj.val);
}
`)
}

// Java: CollapsePropertiesTest.testDontCollapseObjectLiteralPropertyDeclarationInsideLoop()
func (c *CollapsePropertiesTest) TestDontCollapseObjectLiteralPropertyDeclarationInsideLoop() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	// we can collapse `obj.x` but not `obj.x.val`
	c.Test(
		c.Srcs(
			`var obj = {};
for (var i = 0; i < 2; i++) {
  obj.x = {};
  if (i == 0) {
    obj.x.val = 1;
  }
  alert(obj.x.val);
}
`),
		c.Expected(
			`for (var i = 0; i < 2; i++) {
  var obj$x = {};
  if (i == 0) {
    obj$x.val = 1;
  }
  alert(obj$x.val);
}
`))
}

// Java: CollapsePropertiesTest.testDontCollapseConstructorDeclarationInsideLoop()
func (c *CollapsePropertiesTest) TestDontCollapseConstructorDeclarationInsideLoop() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame(
		`for (var i = 0; i < 2; i++) {
  /** @constructor */
  var Foo = function () {}
  if (i == 0) {
    Foo.val = 1;
  }
  alert(Foo.val);
}
`)
}

// Java: CollapsePropertiesTest.testDoCollapsePropertiesDeclaredInsideLoop()
func (c *CollapsePropertiesTest) TestDoCollapsePropertiesDeclaredInsideLoop() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	// It's okay that this property is declared inside a loop as long as the object it's on is not.
	c.Test(
		c.Srcs(
			`var obj = {};
for (var i = 0; i < 2; i++) {
  if (i == 0) {
    obj.val = 1;
  }
  alert(obj.val);
}
`),
		c.Expected(
			`for (var i = 0; i < 2; i++) {
  if (i == 0) {
    var obj$val = 1;
  }
  alert(obj$val);
}
`))
	c.TestSame(
		`var obj = {};
for (var i = 0; i < 2; i++) {
  if (i == 0) {
    obj.val = 1;
  }
  alert?.(obj?.val);
}
`)
}

// Java: CollapsePropertiesTest.testDontCollapseAliasedNamespace_1()
func (c *CollapsePropertiesTest) TestDontCollapseAliasedNamespace_1() {
	c.TestSame("var a = {b: 0}; f(a); a.b;")
}

// Java: CollapsePropertiesTest.testDontCollapseAliasedNamespace_2()
func (c *CollapsePropertiesTest) TestDontCollapseAliasedNamespace_2() {
	c.TestSame("var a = {b: 0}; new f(a); a.b;")
}

// Java: CollapsePropertiesTest.testAliasCreatedForObjectDepth2_1()
func (c *CollapsePropertiesTest) TestAliasCreatedForObjectDepth2_1() {
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0}; var d = 1; d = a.b; a.b.c == d.c;"),
		c.Expected("var a$b = {c: 0}; var d = 1; d = a$b; a$b.c == d.c;"))
	c.TestSame("var a = {}; a.b = {c: 0}; var d = 1; d = a?.b; a.b?.c == d?.c;")
	c.Test(
		c.Srcs(
			`var a = {}; /** @nocollapse */ a.b = {c: 0}; var d = 1; d = a.b;
a.b.c == d.c;
`),
		c.Expected("var a = {}; /** @nocollapse */ a.b = {c: 0}; var d = 1; d = a.b; a.b.c == d.c;"))
}

// Java: CollapsePropertiesTest.testAliasCreatedForObjectDepth2_2()
func (c *CollapsePropertiesTest) TestAliasCreatedForObjectDepth2_2() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0}; for (var p in a.b) { e(a.b[p]); }"),
		c.Expected("var a$b = {c: 0}; for (var p in a$b) { e(a$b[p]); }"))
	c.TestSame("var a = {}; a.b = {c: 0}; for (var p in a?.b) { e?.(a.b?.[p]); }")
}

// Java: CollapsePropertiesTest.testEnumDepth1()
func (c *CollapsePropertiesTest) TestEnumDepth1() {
	c.Test(c.Srcs("/** @enum */ var a = {b: 0, c: 1};"), c.Expected("var a$b = 0; var a$c = 1;"))
	c.Test(
		c.Srcs("/** @enum */ var a = { /** @nocollapse */ b: 0, c: 1};"),
		c.Expected("var a$c = 1; /** @enum */ var a = { /** @nocollapse */ b: 0};"))
}

// Java: CollapsePropertiesTest.testEnumDepth2()
func (c *CollapsePropertiesTest) TestEnumDepth2() {
	c.Test(
		c.Srcs("var a = {}; /** @enum */ a.b = {c: 0, d: 1};"),
		c.Expected("var a$b$c = 0; var a$b$d = 1;"))
	c.TestSame("var a = {}; /** @nocollapse @enum */ a.b = {c: 0, d: 1};")
}

// Java: CollapsePropertiesTest.testAliasCreatedForEnumDepth1_1()
func (c *CollapsePropertiesTest) TestAliasCreatedForEnumDepth1_1() {
	// An enum's values are always collapsed, even if the enum object is
	// referenced in a such a way that an alias is created for it.
	// Unless an enum property has @nocollapse
	c.Test(
		c.Srcs("/** @enum */ var a = {b: 0}; var c = 1; c = a; c.b = 1; a.b != c.b;"),
		c.Expected(
			"var a$b = 0; /** @enum */ var a = {b: a$b}; var c = 1; c = a; c.b = 1; a$b != c.b;"))
	c.Test(
		c.Srcs(
			`/** @enum */ var a = {b:   0};
var c = 1;
c = a;
c.b = 1;
a?.b != c?.b;
`),
		c.Expected(
			`var a$b = 0;
/** @enum */ var a = {b: a$b};
var c = 1;
c = a;
c.b = 1;
a?.b != c?.b
`))
	c.Test(
		c.Srcs(
			`/** @enum */ var a = { /** @nocollapse */ b: 0};
var c = 1;
c = a;
c.b = 1;
a.b == c.b
`),
		c.Expected(
			`/** @enum */ var a = { /** @nocollapse */ b: 0};
var c = 1;
c = a;
c.b = 1;
a.b == c.b
`))
}

// Java: CollapsePropertiesTest.testAliasCreatedForEnumDepth1_2()
func (c *CollapsePropertiesTest) TestAliasCreatedForEnumDepth1_2() {
	c.Test(
		c.Srcs("/** @enum */ var a = {b: 0}; f(a); a.b;"),
		c.Expected("var a$b = 0; /** @enum */ var a = {b: a$b}; f(a); a$b;"))
}

// Java: CollapsePropertiesTest.testAliasCreatedForEnumDepth1_3()
func (c *CollapsePropertiesTest) TestAliasCreatedForEnumDepth1_3() {
	c.Test(
		c.Srcs("/** @enum */ var a = {b: 0}; new f(a); a.b;"),
		c.Expected("var a$b = 0; /** @enum */ var a = {b: a$b}; new f(a); a$b;"))
}

// Java: CollapsePropertiesTest.testAliasCreatedForEnumDepth1_4()
func (c *CollapsePropertiesTest) TestAliasCreatedForEnumDepth1_4() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs("/** @enum */ var a = {b: 0}; for (var p in a) { f(a[p]); }"),
		c.Expected("var a$b = 0; /** @enum */ var a = {b: a$b}; for (var p in a) { f(a[p]); }"))
}

// Java: CollapsePropertiesTest.testAliasCreatedForEnumDepth2_1()
func (c *CollapsePropertiesTest) TestAliasCreatedForEnumDepth2_1() {
	c.Test(
		c.Srcs(
			`var a = {}; /** @enum */ a.b = {c: 0};
var d = 1; d = a.b; d.c = 1; a.b.c != d.c;
`),
		c.Expected(
			`var a$b$c = 0; /** @enum */ var a$b = {c: a$b$c};
var d = 1; d = a$b; d.c = 1; a$b$c != d.c;
`))
	c.TestSame(
		c.Srcs(
			`var a = {}; /** @enum */ a.b = {c: 0};
var d = 1; d = a?.b; d.c = 1; a.b?.c != d?.c;
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
	c.TestSame(
		`var a = {}; /** @nocollapse @enum */ a.b = {c: 0};
var d = 1; d = a.b; d.c = 1; a.b.c == d.c;
`)
	c.TestSame(
		c.Srcs(
			`var a = {}; /** @nocollapse @enum */ a.b = {c: 0};
var d = 1; d = a?.b; d.c = 1; a.b?.c == d?.c;
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
	c.Test(
		c.Srcs(
			`var a = {}; /** @enum */ a.b = {/** @nocollapse */ c: 0};
var d = 1; d = a.b; d.c = 1; a.b.c == d.c;
`),
		c.Expected(
			`/** @enum */ var a$b = { /** @nocollapse */ c: 0};
var d = 1; d = a$b; d.c = 1; a$b.c == d.c;
`))
	c.TestSame(
		c.Srcs(
			`var a = {}; /** @enum */ a.b = {/** @nocollapse */ c: 0};
var d = 1; d = a?.b; d.c = 1; a.b?.c == d?.c;
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: CollapsePropertiesTest.testAliasCreatedForEnumDepth2_2()
func (c *CollapsePropertiesTest) TestAliasCreatedForEnumDepth2_2() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs(
			`var a = {}; /** @enum */ a.b = {c: 0};
for (var p in a.b) { f(a.b[p]); }
`),
		c.Expected(
			`var a$b$c = 0; /** @enum */ var a$b = {c: a$b$c};
for (var p in a$b) { f(a$b[p]); }
`))
	c.TestSame(
		c.Srcs(
			`var a = {}; /** @enum */ a.b = {c: 0};
for (var p in a?.b) { f?.(a.b?.[p]); }
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: CollapsePropertiesTest.testDontCollapseEnumWhenParentNamespaceAliased()
func (c *CollapsePropertiesTest) TestDontCollapseEnumWhenParentNamespaceAliased() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame(
		`var a = {}; var d = 1; d = a; /** @enum */ a.b = {c: 0};
for (var p in a.b) { f(a.b[p]); }
`,
		InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING)
}

// Java: CollapsePropertiesTest.testAliasCreatedForEnumOfObjects()
func (c *CollapsePropertiesTest) TestAliasCreatedForEnumOfObjects() {
	c.Test(
		c.Srcs(
			`var a = {};
/** @enum {Object} */ a.b = {c: {d: 1}};
a.b.c;
searchEnum(a.b);
`),
		c.Expected(
			`var a$b$c = {d: 1};
/** @enum {Object} */ var a$b = {c: a$b$c};
a$b$c;
searchEnum(a$b)
`))
	c.TestSame(
		c.Srcs(
			`var a = {};
/** @enum {Object} */
a.b = {c: {d: 1}};
a.b?.c;
searchEnum(a?.b);
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: CollapsePropertiesTest.testAliasCreatedForEnumOfObjects2()
func (c *CollapsePropertiesTest) TestAliasCreatedForEnumOfObjects2() {
	c.Test(
		c.Srcs(
			`var a = {};
/** @enum {Object} */ a.b = {c: {d: 1}}; a.b.c.d;
searchEnum(a.b);
`),
		c.Expected(
			`var a$b$c = {d: 1}; /** @enum {Object} */ var a$b = {c: a$b$c}; a$b$c.d;
searchEnum(a$b)
`))
	c.TestSame(
		c.Srcs(
			`var a = {};
/** @enum {Object} */ a.b = {c: {d: 1}};
a.b.c?.d;
searchEnum?.(a?.b);
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: CollapsePropertiesTest.testAliasCreatedForPropertyOfEnumOfObjects()
func (c *CollapsePropertiesTest) TestAliasCreatedForPropertyOfEnumOfObjects() {
	c.Test(
		c.Srcs(
			`var a = {};
/** @enum {Object} */ a.b = {c: {d: 1}}; a.b.c;
searchEnum(a.b.c);
`),
		c.Expected("var a$b$c = {d: 1}; a$b$c; searchEnum(a$b$c);"))
	c.Test(
		c.Srcs(
			`var a = {};
/** @enum {Object} */ a.b = {c: {d: 1}}; a.b?.c;
searchEnum?.(a.b?.c);
`),
		c.Expected(
			`var a$b$c = {d:1};
/** @enum {!Object} */ var a$b = {c:a$b$c};
a$b?.c;
searchEnum?.(a$b?.c);
`))
}

// Java: CollapsePropertiesTest.testAliasCreatedForPropertyOfEnumOfObjects2()
func (c *CollapsePropertiesTest) TestAliasCreatedForPropertyOfEnumOfObjects2() {
	c.Test(
		c.Srcs(
			`var a = {};
/** @enum {Object} */ a.b = {c: {d: 1}}; a.b.c.d;
searchEnum(a.b.c);
`),
		c.Expected("var a$b$c = {d: 1}; a$b$c.d; searchEnum(a$b$c);"))
	c.Test(
		c.Srcs(
			`var a = {};
/** @enum {Object} */
a.b = {c: {d: 1}};
a.b.c?.d;
searchEnum?.(a.b?.c);
`),
		c.Expected(
			`var a$b$c = {d:1};
/** @enum {!Object} */
var a$b = {c:a$b$c};
a$b$c?.d;
searchEnum?.(a$b?.c);
`))
}

// Java: CollapsePropertiesTest.testMisusedEnumTag()
func (c *CollapsePropertiesTest) TestMisusedEnumTag() {
	c.TestSame("var a = {}; var d = 1; d = a; a.b = function() {}; /** @enum */ a.b.c = 0; a.b.c;")
}

// Java: CollapsePropertiesTest.testAliasCreatedForFunctionDepth1_1()
func (c *CollapsePropertiesTest) TestAliasCreatedForFunctionDepth1_1() {
	c.TestSame("var a = function(){}; a.b = 1; var c = 1; c = a; c.b = 2; a.b != c.b;")
}

// Java: CollapsePropertiesTest.testAliasCreatedForFunctionDepth1_2()
func (c *CollapsePropertiesTest) TestAliasCreatedForFunctionDepth1_2() {
	c.TestSame("var a = function(){}; a.b = 1; f(a); a.b;")
}

// Java: CollapsePropertiesTest.testAliasCreatedForCtorDepth1_2()
func (c *CollapsePropertiesTest) TestAliasCreatedForCtorDepth1_2() {
	c.Test(
		c.Srcs("/** @constructor */ var a = function(){}; a.b = 1; f(a); a.b;"),
		c.Expected("/** @constructor */ var a = function(){}; var a$b = 1; f(a); a$b;"))
	c.TestSame("/** @constructor */ var a = function(){}; /** @nocollapse */ a.b = 1; f(a); a.b;")
}

// Java: CollapsePropertiesTest.testAliasCreatedForFunctionDepth1_3()
func (c *CollapsePropertiesTest) TestAliasCreatedForFunctionDepth1_3() {
	c.TestSame("var a = function(){}; a.b = 1; new f(a); a.b;")
}

// Java: CollapsePropertiesTest.testAliasCreatedForCtorDepth1_3()
func (c *CollapsePropertiesTest) TestAliasCreatedForCtorDepth1_3() {
	c.Test(
		c.Srcs("/** @constructor */ var a = function(){}; a.b = 1; new f(a); a.b;"),
		c.Expected("/** @constructor */ var a = function(){}; var a$b = 1; new f(a); a$b;"))
	c.TestSame(
		`/** @constructor */ var a = function(){};
/** @nocollapse */ a.b = 1; new f(a); a.b;
`)
}

// Java: CollapsePropertiesTest.testAliasCreatedForClassDepth1_2()
func (c *CollapsePropertiesTest) TestAliasCreatedForClassDepth1_2() {
	c.TestSame(
		"var a = {}; /** @constructor */ a.b = function(){}; f(a); a.b;",
		InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING)
}

// Java: CollapsePropertiesTest.testAliasCreatedForClassDepth1_3()
func (c *CollapsePropertiesTest) TestAliasCreatedForClassDepth1_3() {
	c.TestSame(
		"var a = {}; /** @constructor */ a.b = function(){}; new f(a); a.b;",
		InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING)
}

// Java: CollapsePropertiesTest.testAliasCreatedForClassDepth2_1()
func (c *CollapsePropertiesTest) TestAliasCreatedForClassDepth2_1() {
	c.Test(
		c.Srcs(
			`var a = {};
a.b = {};
/** @constructor */
a.b.c = function(){};
a.b.c;
`),
		c.Expected(
			`/** @constructor */
var a$b$c = function(){};
a$b$c;
`))
	c.Test(
		c.Srcs(
			`var a = {};
a.b = {};
 /** @constructor @nocollapse */
 a.b.c = function(){};
a.b.c;
`),
		c.Expected(
			`var a$b = {};
/** @constructor @nocollapse */
a$b.c = function(){};
a$b.c;
`))
}

// Java: CollapsePropertiesTest.testDontCollapseClassIfParentNamespaceAliased()
func (c *CollapsePropertiesTest) TestDontCollapseClassIfParentNamespaceAliased() {
	c.Test(
		`var a = {};
a.b = {};
/** @constructor */
a.b.c = function(){};
var d = 1;
d = a.b;
a.b.c == d.c;
`,
		`var a$b = {};
/** @constructor */
a$b.c = function(){};
var d = 1;
d = a$b;
a$b.c == d.c;
`,
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
	c.TestSame(
		c.Srcs(
			`var a = {};
a.b = {};
/** @constructor */
a.b.c = function(){};
var d = 1;
d = a?.b;
a.b?.c == d?.c;
`),
		// 2 warnings (for `a` and for `a.b`)
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: CollapsePropertiesTest.testDontCollapseIfParentNamespaceAliased2()
func (c *CollapsePropertiesTest) TestDontCollapseIfParentNamespaceAliased2() {
	c.Test(
		"var a = {}; a.b = {}; /** @constructor */ a.b.c = function(){}; f(a.b); a.b.c;",
		"var a$b = {}; /** @constructor */ a$b.c = function(){}; f(a$b); a$b.c;",
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
	c.TestSame(
		c.Srcs("var a = {}; a.b = {}; /** @constructor */ a.b.c = function(){}; f?.(a?.b); a.b?.c;"),
		// 2 warnings (`a` and `a.b`)
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: CollapsePropertiesTest.testDontCollapseIfParentNamespaceAliased3()
func (c *CollapsePropertiesTest) TestDontCollapseIfParentNamespaceAliased3() {
	c.Test(
		"var a = {}; a.b = {}; /** @constructor */ a.b.c = function(){}; new f(a.b); a.b.c;",
		"var a$b = {}; /** @constructor */ a$b.c = function(){}; new f(a$b); a$b.c;",
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
	c.TestSame(
		c.Srcs(
			`var a = {};
a.b = {};
/** @constructor */
a.b.c = function(){};
new f(a?.b); a.b?.c;
`),
		// 2 warnings (`a` and `a.b`)
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: CollapsePropertiesTest.testAliasCreatedForClassProperty()
func (c *CollapsePropertiesTest) TestAliasCreatedForClassProperty() {
	c.Test(
		c.Srcs(
			`var a = {};
/** @constructor */ a.b = function(){};
a.b.c = {d:3};
new f(a.b.c);
a.b.c.d
`),
		c.Expected(
			`/** @constructor */ var a$b = function(){};
var a$b$c = {d:3};
new f(a$b$c);
a$b$c.d
`))
	c.Test(
		c.Srcs(
			`var a = {};
/** @constructor */
a.b = function(){};
a.b.c = {d:3};
new f(a.b?.c);
a.b.c?.d;
`),
		c.Expected(
			`/** @constructor */
var a$b = function(){};
var a$b$c = {d:3};
// TODO(b/148237949): Collapsing breaks this reference
new f(a$b?.c);
a$b$c?.d;
`))
	c.TestSame(
		`var a = {}; /** @constructor @nocollapse */ a.b = function(){};
a.b.c = {d: 3}; new f(a.b.c); a.b.c.d;
`)
	c.Test(
		c.Srcs(
			`var a = {};
/** @constructor */
a.b = function(){};
/** @nocollapse */
 a.b.c = {d: 3};
 new f(a.b.c);
 a.b.c.d;
`),
		c.Expected(
			`/** @constructor */
 var a$b = function(){};
 /** @nocollapse */
 a$b.c = {d:3};
 new f(a$b.c);
a$b.c.d;
`))
}

// Java: CollapsePropertiesTest.testNestedObjLit()
func (c *CollapsePropertiesTest) TestNestedObjLit() {
	c.Test(
		c.Srcs("var a = {}; a.b = {f: 0, c: {d: 1}}; var e = 1; e = a.b.c.d"),
		c.Expected("var a$b$f = 0; var a$b$c$d = 1; var e = 1; e = a$b$c$d;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {f : 0,         c : {d: 1}}; var e = 1; e = a.b.c?.d"),
		c.Expected("var            a$b$f = 0; var a$b$c = {d: 1} ; var e = 1; e = a$b$c?.d;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {f: 0, /** @nocollapse */ c: {d: 1}}; var e = 1; e = a.b.c.d"),
		c.Expected(
			"var a$b$f = 0; var a$b ={ /** @nocollapse */ c: { d: 1 }}; var e = 1; e = a$b.c.d;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {f: 0, c: {/** @nocollapse */ d: 1}}; var e = 1; e = a.b.c.d"),
		c.Expected("var a$b$f = 0; var a$b$c = { /** @nocollapse */ d: 1}; var e = 1; e = a$b$c.d;"))
}

// Java: CollapsePropertiesTest.testPropGetInsideAnObjLit()
func (c *CollapsePropertiesTest) TestPropGetInsideAnObjLit() {
	c.Test(
		c.Srcs("var x = {}; x.y = 1; var a = {}; a.b = {c: x.y}"),
		c.Expected("var x$y = 1; var a$b$c = x$y;"))
	c.Test(
		c.Srcs("var x = {}; x.y = 1; var a = {}; a.b = {c : x?.y}"),
		c.Expected("var x = {}; x.y = 1; var            a$b$c = x?.y;"))
	c.Test(
		c.Srcs("var x = {}; /** @nocollapse */ x.y = 1; var a = {}; a.b = {c: x.y}"),
		c.Expected("var x = {}; /** @nocollapse */ x.y = 1; var a$b$c = x.y;"))
	c.Test(
		c.Srcs("var x = {}; x.y = 1; var a = {}; a.b = { /** @nocollapse */ c: x.y}"),
		c.Expected("var x$y = 1; var a$b = { /** @nocollapse */ c: x$y};"))
	c.Test(
		c.Srcs("var x = {}; x.y = 1; var a = {}; a.b = { /** @nocollapse */ c: x?.y}"),
		c.Expected("var x = {}; x.y = 1; var         a$b = { /** @nocollapse */ c: x?.y};"))
	c.TestSame(
		`var x = {}; /** @nocollapse */ x.y = 1; var a = {};
/** @nocollapse */ a.b = {c: x.y}
`)
}

// Java: CollapsePropertiesTest.testObjLitWithQuotedKeyThatDoesNotGetRead()
func (c *CollapsePropertiesTest) TestObjLitWithQuotedKeyThatDoesNotGetRead() {
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0, 'd': 1}; var e = 1; e = a.b.c;"),
		c.Expected("var a$b$c = 0; var a$b$d = 1; var e = 1; e = a$b$c;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0, 'd': 1}; var e = 1; e = a.b?.c;"),
		c.Expected("var         a$b = {c: 0, 'd': 1}; var e = 1; e = a$b?.c;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0, /** @nocollapse */ 'd': 1}; var e = 1; e = a.b.c;"),
		c.Expected("var a$b$c = 0; var a$b = {/** @nocollapse */ 'd': 1}; var e = 1; e = a$b$c;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0, /** @nocollapse */ 'd': 1}; var e = 1; e = a.b?.c;"),
		c.Expected("var         a$b = {c: 0, /** @nocollapse */ 'd': 1}; var e = 1; e = a$b?.c;"))
}

// Java: CollapsePropertiesTest.testObjLitWithQuotedKeyThatGetsRead()
func (c *CollapsePropertiesTest) TestObjLitWithQuotedKeyThatGetsRead() {
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0, 'd': 1}; var e = a.b['d'];"),
		c.Expected("var a$b = {c: 0, 'd': 1}; var e = a$b['d'];"))
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0, /** @nocollapse */ 'd': 1}; var e = a.b['d'];"),
		c.Expected("var a$b = {c: 0, /** @nocollapse */ 'd': 1}; var e = a$b['d'];"))
}

// Java: CollapsePropertiesTest.testObjLitWithQuotedKeyThatDoesNotGetReadComputed()
func (c *CollapsePropertiesTest) TestObjLitWithQuotedKeyThatDoesNotGetReadComputed() {
	// quoted/computed does not get read
	c.Test(
		c.Srcs("    var a = {}; a.b = {c: 0, ['d']: 1}; var e = 1; e = a.b.c;"),
		c.Expected("var a$b$c = 0; var a$b = {['d']:1}; var e = 1; e = a$b$c"))
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0, ['d']: 1}; var e = 1; e = a.b?.c;"),
		c.Expected("var         a$b = {c: 0, ['d']: 1}; var e = 1; e = a$b?.c"))
	// quoted/computed gets read
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0, ['d']: 1}; var e = a.b['d'];"),
		c.Expected("var a$b = {c: 0, ['d']: 1}; var e = a$b['d'];"))
	// key collision
	c.Test(
		c.Srcs("    var a = {}; a.b = {c : 0, ['c']: 1}; var e = a.b.c;"),
		c.Expected("var a$b$c = 0;  var a$b = {['c']:1};  var e =  null;"))
}

// Java: CollapsePropertiesTest.testFunctionWithQuotedPropertyThatDoesNotGetRead()
func (c *CollapsePropertiesTest) TestFunctionWithQuotedPropertyThatDoesNotGetRead() {
	c.Test(
		c.Srcs("var a = {}; a.b = function() {}; a.b['d'] = 1;"),
		c.Expected("var a$b = function() {}; a$b['d'] = 1;"))
	c.Test(
		c.Srcs("var a = {}; /** @nocollapse */ a.b = function() {}; a.b['d'] = 1;"),
		c.Expected("var a = {}; /** @nocollapse */ a.b = function() {}; a.b['d'] = 1;"))
	c.Test(
		c.Srcs("var a = {}; a.b = function() {}; /** @nocollapse */ a.b['d'] = 1;"),
		c.Expected("var a$b = function() {}; /** @nocollapse */ a$b['d'] = 1;"))
}

// Java: CollapsePropertiesTest.testFunctionWithQuotedPropertyThatGetsRead()
func (c *CollapsePropertiesTest) TestFunctionWithQuotedPropertyThatGetsRead() {
	c.Test(
		c.Srcs("var a = {}; a.b = function() {}; a.b['d'] = 1; f(a.b['d']);"),
		c.Expected("var a$b = function() {}; a$b['d'] = 1; f(a$b['d']);"))
	c.TestSame(
		`var a = {}; /** @nocollapse */ a.b = function() {};
a.b['d'] = 1; f(a.b['d']);
`)
	c.Test(
		c.Srcs("var a = {}; a.b = function() {}; /** @nocollapse */ a.b['d'] = 1; f(a.b['d']);"),
		c.Expected("var a$b = function() {}; /** @nocollapse */ a$b['d'] = 1; f(a$b['d']);"))
}

// Java: CollapsePropertiesTest.testObjLitAssignedToMultipleNames1()
func (c *CollapsePropertiesTest) TestObjLitAssignedToMultipleNames1() {
	// An object literal that's assigned to multiple names isn't collapsed.
	c.TestSame("var a = b = {c: 0, d: 1}; var e = a.c; var f = b.d;")
	c.TestSame(
		`var a = b = {c: 0, /** @nocollapse */ d: 1}; var e = a.c;
var f = b.d;
`)
}

// Java: CollapsePropertiesTest.testObjLitAssignedToMultipleNames2()
func (c *CollapsePropertiesTest) TestObjLitAssignedToMultipleNames2() {
	c.TestSame("a = b = {c: 0, d: 1}; var e = a.c; var f = b.d;")
}

// Java: CollapsePropertiesTest.testObjLitRedefinedInGlobalScope()
func (c *CollapsePropertiesTest) TestObjLitRedefinedInGlobalScope() {
	c.TestSame("a = {b: 0}; a = {c: 1}; var d = a.b; var e = a.c;")
}

// Java: CollapsePropertiesTest.testObjLitRedefinedInLocalScope()
func (c *CollapsePropertiesTest) TestObjLitRedefinedInLocalScope() {
	c.Test(
		c.Srcs("var a = {}; a.b = {c: 0}; function d() { a.b = {c: 1}; } e(a.b.c);"),
		c.Expected("var a$b = {c: 0}; function d() { a$b = {c: 1}; } e(a$b.c);"))
	c.TestSame("var a = {};/** @nocollapse */ a.b = {c: 0}; function d() { a.b = {c: 1};} e(a.b.c);")
	// redefinition with @nocollapse
	c.Test(
		c.Srcs(
			`var a = {}; a.b = {c: 0};
function d() { a.b = {/** @nocollapse */ c: 1}; } e(a.b.c);
`),
		c.Expected("var a$b = {c: 0}; function d() { a$b = {/** @nocollapse */ c: 1}; } e(a$b.c);"))
}

// Java: CollapsePropertiesTest.testObjLitAssignedInTernaryExpression1()
func (c *CollapsePropertiesTest) TestObjLitAssignedInTernaryExpression1() {
	c.TestSame("a = x ? {b: 0} : d; var c = a.b;")
}

// Java: CollapsePropertiesTest.testObjLitAssignedInTernaryExpression2()
func (c *CollapsePropertiesTest) TestObjLitAssignedInTernaryExpression2() {
	c.TestSame("a = x ? {b: 0} : {b: 1}; var c = a.b;")
}

// Java: CollapsePropertiesTest.testGlobalVarSetToObjLitConditionally1()
func (c *CollapsePropertiesTest) TestGlobalVarSetToObjLitConditionally1() {
	c.TestSame("var a; if (x) a = {b: 0}; var c = x ? a.b : 0;")
}

// Java: CollapsePropertiesTest.testGlobalVarSetToObjLitConditionally1b()
func (c *CollapsePropertiesTest) TestGlobalVarSetToObjLitConditionally1b() {
	c.Test(
		c.Srcs("if (x) var a = {b: 0}; var c = x ? a.b : 0;"),
		c.Expected("if (x) var a$b = 0; var c = x ? a$b : 0;"))
	c.TestSame("if (x) var a = { /** @nocollapse */ b: 0}; var c = x ? a.b : 0;")
}

// Java: CollapsePropertiesTest.testGlobalVarSetToObjLitConditionally2()
func (c *CollapsePropertiesTest) TestGlobalVarSetToObjLitConditionally2() {
	c.Test(
		c.Srcs("if (x) var a = {b: 0}; var c = 1; c = a.b; var d = a.c;"),
		c.Expected("if (x){ var a$b = 0; var a = {}; }var c = 1; c = a$b; var d = a.c;"))
	c.TestSame("if (x) var a = {b: 0}; var c = 1; c = a?.b; var d = a?.c;")
	c.TestSame("if (x) var a = {/** @nocollapse */ b: 0}; var c = 1; c = a.b; var d = a.c;")
}

// Java: CollapsePropertiesTest.testGlobalVarSetToObjLitConditionally3()
func (c *CollapsePropertiesTest) TestGlobalVarSetToObjLitConditionally3() {
	c.TestSame("var a; if (x) a = {b: 0}; else a = {b: 1}; var c = a.b;")
	c.TestSame(
		`var a; if (x) a = {b: 0}; else a = {/** @nocollapse */ b: 1};
var c = a.b;
`)
}

// Java: CollapsePropertiesTest.testObjectPropertySetToObjLitConditionally()
func (c *CollapsePropertiesTest) TestObjectPropertySetToObjLitConditionally() {
	c.Test(
		c.Srcs("var a = {}; if (x) a.b = {c: 0}; var d = a.b ? a.b.c : 0;"),
		c.Expected("if (x){ var a$b$c = 0; var a$b = {} } var d = a$b ? a$b$c : 0;"))
	c.Test(
		c.Srcs("var a = {}; if (x) a.b = {/** @nocollapse */ c: 0}; var d = a.b ? a.b.c : 0;"),
		c.Expected("if (x){ var a$b = {/** @nocollapse */ c: 0};} var d = a$b ? a$b.c : 0;"))
}

// Java: CollapsePropertiesTest.testFunctionPropertySetToObjLitConditionally()
func (c *CollapsePropertiesTest) TestFunctionPropertySetToObjLitConditionally() {
	c.Test(
		c.Srcs("function a() {} if (x) a.b = {c: 0}; var d = a.b ? a.b.c : 0;"),
		c.Expected(
			`function a() {} if (x){ var a$b$c = 0; var a$b = {} }
var d = a$b ? a$b$c : 0;
`))
	c.TestSame(
		`function a() {} if (x) /** @nocollapse */ a.b = {c: 0};
var d = a.b ? a.b.c : 0;
`)
	c.Test(
		c.Srcs(
			`function a() {}
if (x) a.b = {/** @nocollapse */ c: 0};
var d = a.b ? a.b.c : 0;
`),
		c.Expected(
			`function a() {}
if (x) {
  var a$b = {/** @nocollapse */ c: 0};
}
var d = a$b ? a$b.c : 0
`))
}

// Java: CollapsePropertiesTest.testPrototypePropertySetToAnObjectLiteral()
func (c *CollapsePropertiesTest) TestPrototypePropertySetToAnObjectLiteral() {
	c.Test(
		c.Srcs("var a = {b: function(){}}; a.b.prototype.c = {d: 0};"),
		c.Expected("var a$b = function(){}; a$b.prototype.c = {d: 0};"))
	c.TestSame(
		`var a = {/** @nocollapse */ b: function(){}};
a.b.prototype.c = {d: 0};
`)
}

// Java: CollapsePropertiesTest.testObjectPropertyResetInLocalScope()
func (c *CollapsePropertiesTest) TestObjectPropertyResetInLocalScope() {
	c.Test(
		c.Srcs("var z = {}; z.a = 0; function f() {z.a = 5; return z.a}"),
		c.Expected("var z$a = 0; function f() {z$a = 5; return z$a}"))
	c.TestSame("var z = {}; z.a = 0; function f() {z.a = 5; return z?.a}")
	c.TestSame(
		`var z = {}; z.a = 0;
function f() { /** @nocollapse */ z.a = 5; return z.a}
`)
	c.TestSame(
		`var z = {}; /** @nocollapse */ z.a = 0;
function f() {z.a = 5; return z.a}
`)
}

// Java: CollapsePropertiesTest.testFunctionPropertyResetInLocalScope()
func (c *CollapsePropertiesTest) TestFunctionPropertyResetInLocalScope() {
	c.Test(
		c.Srcs("function z() {}     z.a = 0; function f() {z.a = 5; return z.a}"),
		c.Expected("function z() {} var z$a = 0; function f() {z$a = 5; return z$a}"))
	c.TestSame("function z() {} z.a = 0; function f() {z.a = 5; return z?.a}")
	c.TestSame(
		`function z() {} /** @nocollapse */ z.a = 0;
function f() {z.a = 5; return z.a}
`)
	c.TestSame(
		`function z() {} z.a = 0;
function f() { /** @nocollapse */ z.a = 5; return z.a}
`)
}

// Java: CollapsePropertiesTest.testNamespaceResetInGlobalScope1()
func (c *CollapsePropertiesTest) TestNamespaceResetInGlobalScope1() {
	c.Test(
		"var a = {}; /** @constructor */ a.b = function() {}; a = {};",
		"var a = {}; /** @constructor */ var a$b = function() {}; a = {};",
		c.Warning(InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING))
	c.TestSame(
		`var a = {}; /** @constructor @nocollapse */a.b = function() {};
a = {};
`,
		InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING)
}

// Java: CollapsePropertiesTest.testNamespaceResetInGlobalScope2()
func (c *CollapsePropertiesTest) TestNamespaceResetInGlobalScope2() {
	c.Test(
		"var a = {}; a = {}; /** @constructor */ a.b = function() {};",
		"var a = {}; a = {}; /** @constructor */ var a$b = function() {};",
		c.Warning(InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING))
	c.TestSame(
		"var a = {}; a = {}; /** @constructor @nocollapse */a.b = function() {};",
		InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING)
}

// Java: CollapsePropertiesTest.testNamespaceResetInGlobalScope3()
func (c *CollapsePropertiesTest) TestNamespaceResetInGlobalScope3() {
	c.Test(
		c.Srcs("var a = {}; /** @constructor */ a.b = function() {}; a = a || {};"),
		c.Expected("var a = {}; /** @constructor */ var a$b = function() {}; a = a || {};"))
	c.TestSame("var a = {}; /** @constructor @nocollapse */ a.b = function() {}; a = a || {};")
}

// Java: CollapsePropertiesTest.testNamespaceResetInGlobalScope4()
func (c *CollapsePropertiesTest) TestNamespaceResetInGlobalScope4() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs("var a = {}; /** @constructor */ a.b = function() {}; var a = a || {};"),
		c.Expected("var a = {}; /** @constructor */ var a$b = function() {}; var a = a || {};"))
	c.TestSame("var a = {}; /** @constructor @nocollapse */a.b = function() {}; var a = a || {};")
}

// Java: CollapsePropertiesTest.testNamespaceResetInLocalScope1()
func (c *CollapsePropertiesTest) TestNamespaceResetInLocalScope1() {
	c.Test(
		"var a = {}; /** @constructor */ a.b = function() {}; function f() { a = {}; }",
		"var a = {}; /** @constructor */ var a$b = function() {}; function f() { a = {}; }",
		c.Warning(InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING))
	c.TestSame(
		`var a = {}; /** @constructor @nocollapse */a.b = function() {};
 function f() { a = {}; }
`,
		InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING)
}

// Java: CollapsePropertiesTest.testNamespaceResetInLocalScope2()
func (c *CollapsePropertiesTest) TestNamespaceResetInLocalScope2() {
	c.Test(
		"var a = {}; function f() { a = {}; } /** @constructor */ a.b = function() {};",
		"var a = {}; function f() { a = {}; } /** @constructor */ var a$b = function() {};",
		c.Warning(InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING))
	c.TestSame(
		`var a = {}; function f() { a = {}; }
 /** @constructor @nocollapse */a.b = function() {};
`,
		InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING)
}

// Java: CollapsePropertiesTest.testNamespaceDefinedInLocalScope()
func (c *CollapsePropertiesTest) TestNamespaceDefinedInLocalScope() {
	c.Test(
		c.Srcs(
			`var a = {};
(function() { a.b = {}; })();
/** @constructor */ a.b.c = function() {};
`),
		c.Expected(
			`var a$b;
(function() { a$b = {}; })();
/** @constructor */ var a$b$c = function() {}
`))
	c.Test(
		c.Srcs(
			`var a = {}; (function() { /** @nocollapse */ a.b = {}; })();
 /** @constructor */a.b.c = function() {};
`),
		c.Expected(
			`var a = {}; (function() { /** @nocollapse */ a.b = {}; })();
 /** @constructor */ var a$b$c = function() {};
`))
	c.Test(
		c.Srcs(
			`var a = {}; (function() { a.b = {}; })();
 /** @constructor @nocollapse */a.b.c = function() {};
`),
		c.Expected(
			`var a$b; (function() { a$b = {}; })();
/** @constructor @nocollapse */ a$b.c = function() {};
`))
}

// Java: CollapsePropertiesTest.testAddPropertyToObjectInLocalScopeDepth1()
func (c *CollapsePropertiesTest) TestAddPropertyToObjectInLocalScopeDepth1() {
	c.Test(
		c.Srcs("var a = {b: 0}; function f() { a.c = 5; return a.c; }"),
		c.Expected("var a$b = 0; var a$c; function f() { a$c = 5; return a$c; }"))
	c.TestSame("var a = {b: 0}; function f() { a.c = 5; return a?.c; }")
}

// Java: CollapsePropertiesTest.testAddPropertyToObjectInLocalScopeDepth2()
func (c *CollapsePropertiesTest) TestAddPropertyToObjectInLocalScopeDepth2() {
	c.Test(
		c.Srcs("var a = {}; a.b = {}; (function() {a.b.c = 0;})(); x = a.b.c;"),
		c.Expected("var a$b$c; (function() {a$b$c = 0;})(); x = a$b$c;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {}; (function() {a.b.c = 0;})(); x = a.b?.c;"),
		c.Expected("var         a$b = {}; (function() {a$b.c = 0;})(); x = a$b?.c;"))
}

// Java: CollapsePropertiesTest.testAddPropertyToFunctionInLocalScopeDepth1()
func (c *CollapsePropertiesTest) TestAddPropertyToFunctionInLocalScopeDepth1() {
	c.Test(
		c.Srcs("function a() {} function f() { a.c = 5; return a.c; }"),
		c.Expected("function a() {} var a$c; function f() { a$c = 5; return a$c; }"))
	c.TestSame("function a() {} function f() { a.c = 5; return a?.c; }")
}

// Java: CollapsePropertiesTest.testAddPropertyToFunctionInLocalScopeDepth2()
func (c *CollapsePropertiesTest) TestAddPropertyToFunctionInLocalScopeDepth2() {
	c.Test(
		c.Srcs("var a = {}; a.b = function() {}; function f() {a.b.c = 0;}"),
		c.Expected("var a$b = function() {}; var a$b$c; function f() {a$b$c = 0;}"))
}

// Java: CollapsePropertiesTest.testAddPropertyToUncollapsibleFunctionInLocalScopeDepth1()
func (c *CollapsePropertiesTest) TestAddPropertyToUncollapsibleFunctionInLocalScopeDepth1() {
	c.TestSame("function a() {} var c = 1; c = a; (function() {a.b = 0;})(); a.b;")
}

// Java: CollapsePropertiesTest.testAddPropertyToUncollapsibleFunctionInLocalScopeDepth2()
func (c *CollapsePropertiesTest) TestAddPropertyToUncollapsibleFunctionInLocalScopeDepth2() {
	c.Test(
		c.Srcs(
			`var a = {}; a.b = function (){}; var d = 1; d = a.b;
(function() {a.b.c = 0;})(); a.b.c;
`),
		c.Expected(
			`var a$b = function (){}; var d = 1; d = a$b;
(function() {a$b.c = 0;})(); a$b.c;
`))
	c.TestSame(
		`var a = {}; a.b = function (){}; var d = 1; d = a?.b;
(function() {a.b.c = 0;})(); a.b?.c;
`)
}

// Java: CollapsePropertiesTest.testResetObjectPropertyInLocalScope()
func (c *CollapsePropertiesTest) TestResetObjectPropertyInLocalScope() {
	c.Test(
		c.Srcs("var a = {b: 0}; a.c = 1; function f() { a.c = 5; }"),
		c.Expected("var a$b = 0; var a$c = 1; function f() { a$c = 5; }"))
}

// Java: CollapsePropertiesTest.testResetFunctionPropertyInLocalScope()
func (c *CollapsePropertiesTest) TestResetFunctionPropertyInLocalScope() {
	c.Test(
		c.Srcs("function a() {}; a.c = 1; function f() { a.c = 5; }"),
		c.Expected("function a() {}; var a$c = 1; function f() { a$c = 5; }"))
}

// Java: CollapsePropertiesTest.testGlobalNameReferencedInLocalScopeBeforeDefined1()
func (c *CollapsePropertiesTest) TestGlobalNameReferencedInLocalScopeBeforeDefined1() {
	// Because referencing global names earlier in the source code than they're
	// defined is such a common practice, we collapse them even though a runtime
	// exception could result (in the off-chance that the function gets called
	// before the alias variable is defined).
	c.Test(
		c.Srcs("var a = {b: 0}; function f() { a.c = 5; } a.c = 1;"),
		c.Expected("var a$b = 0; function f() { a$c = 5; } var a$c = 1;"))
}

// Java: CollapsePropertiesTest.testGlobalNameReferencedInLocalScopeBeforeDefined2()
func (c *CollapsePropertiesTest) TestGlobalNameReferencedInLocalScopeBeforeDefined2() {
	c.Test(
		c.Srcs("var a = {b: 0}; function f() { return a.c; } a.c = 1;"),
		c.Expected("var a$b = 0; function f() { return a$c; } var a$c = 1;"))
	c.TestSame("var a = {b: 0}; function f() { return a?.c; } a.c = 1;")
}

// Java: CollapsePropertiesTest.testTwiceDefinedGlobalNameDepth1_1()
func (c *CollapsePropertiesTest) TestTwiceDefinedGlobalNameDepth1_1() {
	c.TestSame(
		`var a = {}; function f() { a.b(); }
a = function() {}; a.b = function() {};
`)
}

// Java: CollapsePropertiesTest.testTwiceDefinedGlobalNameDepth1_2()
func (c *CollapsePropertiesTest) TestTwiceDefinedGlobalNameDepth1_2() {
	c.TestSame(
		`var a = {}; /** @constructor */ a = function() {};
a.b = {}; a.b.c = 0; function f() { a.b.d = 1; }
`)
}

// Java: CollapsePropertiesTest.testTwiceDefinedGlobalNameDepth2()
func (c *CollapsePropertiesTest) TestTwiceDefinedGlobalNameDepth2() {
	c.Test(
		c.Srcs(
			`var a = {}; a.b = {}; function f() { a.b.c(); }
a.b = function() {}; a.b.c = function() {};
`),
		c.Expected(
			`var a$b = {}; function f() { a$b.c(); }
a$b = function() {}; a$b.c = function() {};
`))
}

// Java: CollapsePropertiesTest.testFunctionCallDepth1()
func (c *CollapsePropertiesTest) TestFunctionCallDepth1() {
	c.Test(
		c.Srcs("var a = {}; a.b = function(){}; var c = a.b();"),
		c.Expected("var a$b = function(){}; var c = a$b()"))
}

// Java: CollapsePropertiesTest.testFunctionCallDepth2()
func (c *CollapsePropertiesTest) TestFunctionCallDepth2() {
	c.Test(
		c.Srcs("var a = {}; a.b = {}; a.b.c = function(){}; a.b.c();"),
		c.Expected("var a$b$c = function(){}; a$b$c();"))
}

// Java: CollapsePropertiesTest.testFunctionAlias1()
func (c *CollapsePropertiesTest) TestFunctionAlias1() {
	c.Test(
		c.Srcs("var a = {}; a.b = {}; a.b.c = function(){}; a.b.d = a.b.c;a.b.d=null"),
		c.Expected("var a$b$c = function(){}; var a$b$d = a$b$c;a$b$d=null;"))
	c.Test(
		c.Srcs("var a = {}; a.b = {}; a.b.c = function(){}; a.b.d = a.b?.c; a.b.d=null"),
		c.Expected("var         a$b = {}; a$b.c = function(){}; a$b.d = a$b?.c; a$b.d=null;"))
}

// Java: CollapsePropertiesTest.testCallToRedefinedFunction()
func (c *CollapsePropertiesTest) TestCallToRedefinedFunction() {
	c.Test(
		c.Srcs("var a = {}; a.b = function(){}; a.b = function(){}; a.b();"),
		c.Expected("var a$b = function(){}; a$b = function(){}; a$b();"))
}

// Java: CollapsePropertiesTest.testCollapsePrototypeName()
func (c *CollapsePropertiesTest) TestCollapsePrototypeName() {
	c.Test(
		c.Srcs(
			`var a = {}; a.b = {}; a.b.c = function(){};
a.b.c.prototype.d = function(){}; (new a.b.c()).d();
`),
		c.Expected(
			`var a$b$c = function(){}; a$b$c.prototype.d = function(){};
new a$b$c().d();
`))
	c.Test(
		c.Srcs(
			`var a = {}; a.b = {}; a.b.c = function(){};
a.b.c.prototype.d = function(){}; (new a.b.c()).d?.();
`),
		c.Expected(
			`var a$b$c = function(){}; a$b$c.prototype.d = function(){};
(new a$b$c).d?.();
`))
}

// Java: CollapsePropertiesTest.testReferencedPrototypeProperty()
func (c *CollapsePropertiesTest) TestReferencedPrototypeProperty() {
	c.Test(
		c.Srcs(
			`var a = {b: {}}; a.b.c = function(){}; a.b.c.prototype.d = {};
e = a.b.c.prototype.d;
`),
		c.Expected(
			`var a$b$c = function(){}; a$b$c.prototype.d = {};
e = a$b$c.prototype.d;
`))
}

// Java: CollapsePropertiesTest.testSetStaticAndPrototypePropertiesOnFunction()
func (c *CollapsePropertiesTest) TestSetStaticAndPrototypePropertiesOnFunction() {
	c.Test(
		c.Srcs("var a = {}; a.b = function(){}; a.b.prototype.d = 0; a.b.c = 1;"),
		c.Expected("var a$b = function(){}; a$b.prototype.d = 0; var a$b$c = 1;"))
}

// Java: CollapsePropertiesTest.testReadUndefinedPropertyDepth1()
func (c *CollapsePropertiesTest) TestReadUndefinedPropertyDepth1() {
	c.Test(c.Srcs("var a = {b: 0}; var c = a.d;"), c.Expected("var a$b = 0; var a = {}; var c = a.d;"))
}

// Java: CollapsePropertiesTest.testReadUndefinedPropertyDepth2()
func (c *CollapsePropertiesTest) TestReadUndefinedPropertyDepth2() {
	c.Test(
		c.Srcs("var a = {b: {c: 0}}; f(a.b.c); f(a.b.d);"),
		c.Expected("var a$b$c = 0; var a$b = {}; f(a$b$c); f(a$b.d);"))
	c.Test(
		c.Srcs("var a = {b : {c: 0}}; f?.(a.b?.c); f?.(a.b?.d);"),
		c.Expected("var    a$b = {c: 0} ; f?.(a$b?.c); f?.(a$b?.d);"))
}

// Java: CollapsePropertiesTest.testCallUndefinedMethodOnObjLitDepth1()
func (c *CollapsePropertiesTest) TestCallUndefinedMethodOnObjLitDepth1() {
	c.Test(c.Srcs("var a = {b: 0}; a.c();"), c.Expected("var a$b = 0; var a = {}; a.c();"))
}

// Java: CollapsePropertiesTest.testCallUndefinedMethodOnObjLitDepth2()
func (c *CollapsePropertiesTest) TestCallUndefinedMethodOnObjLitDepth2() {
	c.Test(
		c.Srcs("var a = {b: {}}; a.b.c = function() {}; a.b.c(); a.b.d();"),
		c.Expected("var a$b = {}; var a$b$c = function() {}; a$b$c(); a$b.d();"))
}

// Java: CollapsePropertiesTest.testPropertiesOfAnUndefinedVar()
func (c *CollapsePropertiesTest) TestPropertiesOfAnUndefinedVar() {
	c.TestSame("a.document = d; f(a.document.innerHTML);")
}

// Java: CollapsePropertiesTest.testPropertyOfAnObjectThatIsNeitherFunctionNorObjLit()
func (c *CollapsePropertiesTest) TestPropertyOfAnObjectThatIsNeitherFunctionNorObjLit() {
	c.TestSame("var a = window; a.document = d; f(a.document)")
}

// Java: CollapsePropertiesTest.testStaticFunctionReferencingThis1()
func (c *CollapsePropertiesTest) TestStaticFunctionReferencingThis1() {
	// Note: Google's JavaScript Style Guide says to avoid using the 'this'
	// keyword in a static function.
	c.Test(
		"var a = {}; a.b = function() {this.c};",
		"var a$b = function() {this.c}; ",
		c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
}

// Java: CollapsePropertiesTest.testStaticFunctionReferencingThis2()
func (c *CollapsePropertiesTest) TestStaticFunctionReferencingThis2() {
	// This gives no warning, because `this` is in a scope whose name is not
	// getting collapsed.
	c.Test(
		c.Srcs("var a = {}; a.b = function() { return function(){ return this; }; };"),
		c.Expected("var a$b = function() { return function(){ return this; }; };"))
}

// Java: CollapsePropertiesTest.testPropAssignedToFunctionReferencingThis()
func (c *CollapsePropertiesTest) TestPropAssignedToFunctionReferencingThis() {
	// This gives no warning, because `this` is in a scope whose name is not getting collapsed.
	c.TestNoWarning("var a = {}; var b = function() { this.c;}; a.b=b;")
}

// Java: CollapsePropertiesTest.testStaticFunctionReferencingThis3()
func (c *CollapsePropertiesTest) TestStaticFunctionReferencingThis3() {
	c.Test(
		"var a = {b: function() {this.c}};",
		"var a$b = function() { this.c };",
		c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
}

// Java: CollapsePropertiesTest.testStaticFunctionReferencingThis4()
func (c *CollapsePropertiesTest) TestStaticFunctionReferencingThis4() {
	c.Test(
		c.Srcs("var a = {/** @this {Element} */ b: function() {this.c}};"),
		c.Expected("var a$b = function() { this.c };"))
}

// Java: CollapsePropertiesTest.testPrototypeMethodReferencingThis()
func (c *CollapsePropertiesTest) TestPrototypeMethodReferencingThis() {
	c.TestSame("var A = function(){}; A.prototype = {b: function() {this.c}};")
}

// Java: CollapsePropertiesTest.testConstructorReferencingThis()
func (c *CollapsePropertiesTest) TestConstructorReferencingThis() {
	c.Test(
		c.Srcs(
			`var a = {};
/** @constructor */ a.b = function() { this.a = 3; };
`),
		c.Expected("/** @constructor */ var a$b = function() { this.a = 3; };"))
}

// Java: CollapsePropertiesTest.testRecordReferencingThis()
func (c *CollapsePropertiesTest) TestRecordReferencingThis() {
	c.Test(
		c.Srcs(
			`/** @const */ var a = {};
/** @record */ a.b = function() { /** @type {string} */ this.a; };
`),
		c.Expected("/** @record */ var a$b = function() { /** @type {string} */ this.a; };"))
}

// Java: CollapsePropertiesTest.testSafeReferenceOfThis()
func (c *CollapsePropertiesTest) TestSafeReferenceOfThis() {
	c.Test(
		c.Srcs("var a = {}; /** @this {Object} */ a.b = function() { this.a = 3; };"),
		c.Expected(" /** @this {Object} */ var a$b = function() { this.a = 3; };"))
}

// Java: CollapsePropertiesTest.testGlobalFunctionReferenceOfThis()
func (c *CollapsePropertiesTest) TestGlobalFunctionReferenceOfThis() {
	c.TestSame("var a = function() { this.a = 3; };")
}

// Java: CollapsePropertiesTest.testFunctionGivenTwoNames()
func (c *CollapsePropertiesTest) TestFunctionGivenTwoNames() {
	// It's okay to collapse f's properties because g is not added to the
	// global scope as an alias for f. (Try it in your browser.)
	c.Test(
		c.Srcs("var f = function g() {}; f.a = 1; h(f.a);"),
		c.Expected("var f = function g() {}; var f$a = 1; h(f$a);"))
	c.TestSame("var f = function g() {}; f.a = 1; h?.(f?.a);")
}

// Java: CollapsePropertiesTest.testObjLitWithUsedNumericKey()
func (c *CollapsePropertiesTest) TestObjLitWithUsedNumericKey() {
	c.TestSame("a = {40: {}, c: {}}; var d = a[40]; var e = a.c;")
}

// Java: CollapsePropertiesTest.testObjLitWithUnusedNumericKey()
func (c *CollapsePropertiesTest) TestObjLitWithUnusedNumericKey() {
	c.Test(
		c.Srcs("var a = {40: {},       c : {}}; var e = 1; e = a.c;"),
		c.Expected("var    a$1 = {}; var a$c = {} ; var e = 1; e = a$c"))
	// incorrect
	c.TestSame("var a = {40: {}, c: {}}; var e = 1; e =  a?.c;")
}

// Java: CollapsePropertiesTest.testObjLitWithNonIdentifierKeys()
func (c *CollapsePropertiesTest) TestObjLitWithNonIdentifierKeys() {
	c.TestSame("a = {' ': 0, ',': 1}; var c = a[' '];")
	c.TestSame(
		`var FOO = {
  'bar': {
    'baz,qux': {
      'beep': 'xxxxx',
    },
  }
};
alert(FOO);
`)
}

// Java: CollapsePropertiesTest.testChainedAssignments1()
func (c *CollapsePropertiesTest) TestChainedAssignments1() {
	c.Test(c.Srcs("var x = {}; x.y = a = 0;"), c.Expected("var x$y = a = 0;"))
}

// Java: CollapsePropertiesTest.testChainedAssignments2()
func (c *CollapsePropertiesTest) TestChainedAssignments2() {
	c.Test(c.Srcs("var x = {}; x.y = a = b = c();"), c.Expected("var x$y = a = b = c();"))
}

// Java: CollapsePropertiesTest.testChainedAssignments3()
func (c *CollapsePropertiesTest) TestChainedAssignments3() {
	c.Test(c.Srcs("var x = {y: 1}; a = b = x.y;"), c.Expected("var x$y = 1; a = b = x$y;"))
	c.TestSame("var x = {y: 1}; a = b = x?.y;")
}

// Java: CollapsePropertiesTest.testChainedAssignments4()
func (c *CollapsePropertiesTest) TestChainedAssignments4() {
	c.TestSame("var x = {}; a = b = x.y;")
}

// Java: CollapsePropertiesTest.testChainedAssignments5()
func (c *CollapsePropertiesTest) TestChainedAssignments5() {
	c.Test(c.Srcs("var x = {}; a = x.y = 0;"), c.Expected("var x$y; a = x$y = 0;"))
}

// Java: CollapsePropertiesTest.testChainedAssignments6()
func (c *CollapsePropertiesTest) TestChainedAssignments6() {
	c.Test(c.Srcs("var x = {}; a = x.y = b = c();"), c.Expected("var x$y; a = x$y = b = c();"))
}

// Java: CollapsePropertiesTest.testChainedAssignments7()
func (c *CollapsePropertiesTest) TestChainedAssignments7() {
	c.Test(
		"var x = {}; a = x.y = {}; /** @constructor */ x.y.z = function() {};",
		"var x$y; a = x$y = {}; /** @constructor */ x$y.z = function() {};",
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: CollapsePropertiesTest.testChainedVarAssignments1()
func (c *CollapsePropertiesTest) TestChainedVarAssignments1() {
	c.Test(c.Srcs("var x = {y: 1}; var a = x.y = 0;"), c.Expected("var x$y = 1; var a = x$y = 0;"))
}

// Java: CollapsePropertiesTest.testChainedVarAssignments2()
func (c *CollapsePropertiesTest) TestChainedVarAssignments2() {
	c.Test(
		c.Srcs("var x = {y: 1}; var a = x.y = b = 0;"),
		c.Expected("var x$y = 1; var a = x$y = b = 0;"))
}

// Java: CollapsePropertiesTest.testChainedVarAssignments3()
func (c *CollapsePropertiesTest) TestChainedVarAssignments3() {
	c.Test(
		c.Srcs("var x = {y: {z: 1}}; var b = 0; var a = x.y.z = 1; var c = 2;"),
		c.Expected("var x$y$z = 1; var b = 0; var a = x$y$z = 1; var c = 2;"))
}

// Java: CollapsePropertiesTest.testChainedVarAssignments4()
func (c *CollapsePropertiesTest) TestChainedVarAssignments4() {
	c.Test(c.Srcs("var x = {}; var a = b = x.y = 0;"), c.Expected("var x$y; var a = b = x$y = 0;"))
}

// Java: CollapsePropertiesTest.testChainedVarAssignments5()
func (c *CollapsePropertiesTest) TestChainedVarAssignments5() {
	c.Test(
		c.Srcs("var x = {y: {}}; var a = b = x.y.z = 0;"),
		c.Expected("var x$y$z; var a = b = x$y$z = 0;"))
}

// Java: CollapsePropertiesTest.testChainedVarAssignments6()
func (c *CollapsePropertiesTest) TestChainedVarAssignments6() {
	c.TestSame("var a = x = 0; var x;")
}

// Java: CollapsePropertiesTest.testChainedVarAssignments7()
func (c *CollapsePropertiesTest) TestChainedVarAssignments7() {
	c.TestSame("x = {}; var a = x.y = 0; var x;")
}

// Java: CollapsePropertiesTest.testPeerAndSubpropertyOfUncollapsibleProperty()
func (c *CollapsePropertiesTest) TestPeerAndSubpropertyOfUncollapsibleProperty() {
	c.Test(
		c.Srcs(
			`var x = {}; var a = x.y = 0; x.w = 1; x.y.z = 2;
b = x.w; c = x.y.z;
`),
		c.Expected(
			`var x$y; var a = x$y = 0; var x$w = 1; x$y.z = 2;
b = x$w; c = x$y.z;
`))
	c.TestSame(
		`var x = {}; var a = x.y = 0; x.w = 1; x.y.z = 2;
b = x?.w; c = x.y?.z;
`)
}

// Java: CollapsePropertiesTest.testComplexAssignmentAfterInitialAssignment()
func (c *CollapsePropertiesTest) TestComplexAssignmentAfterInitialAssignment() {
	c.Test(
		c.Srcs("var d = {}; d.e = {}; d.e.f = 0; a = b = d.e.f = 1;"),
		c.Expected("var d$e$f = 0; a = b = d$e$f = 1;"))
}

// Java: CollapsePropertiesTest.testRenamePrefixOfUncollapsibleProperty()
func (c *CollapsePropertiesTest) TestRenamePrefixOfUncollapsibleProperty() {
	c.Test(
		c.Srcs("var d = {}; d.e = {}; a = b = d.e.f = 0;"),
		c.Expected("var d$e$f; a = b = d$e$f = 0;"))
}

// Java: CollapsePropertiesTest.testNewOperator()
func (c *CollapsePropertiesTest) TestNewOperator() {
	// Using the new operator on a name doesn't prevent its (static) properties
	// from getting collapsed.
	c.Test(
		c.Srcs("var a = {}; a.b = function() {}; a.b.c = 1; var d = new a.b();"),
		c.Expected("var a$b = function() {}; var a$b$c = 1; var d = new a$b();"))
}

// Java: CollapsePropertiesTest.testMethodCall()
func (c *CollapsePropertiesTest) TestMethodCall() {
	c.Test(
		c.Srcs("var a = {}; a.b = function() {}; var d = a.b();"),
		c.Expected("var a$b = function() {}; var d = a$b();"))
}

// Java: CollapsePropertiesTest.testObjLitDefinedInLocalScopeIsLeftAlone()
func (c *CollapsePropertiesTest) TestObjLitDefinedInLocalScopeIsLeftAlone() {
	c.Test(
		c.Srcs(
			`var a = {}; a.b = function() {};
a.b.prototype.f_ = function() {
  var x = { p: '', q: '', r: ''}; var y = x.q;
};
`),
		c.Expected(
			`var a$b = function() {};
a$b.prototype.f_ = function() {
  var x = { p: '', q: '', r: ''}; var y = x.q;
};
`))
}

// Java: CollapsePropertiesTest.testPropertiesOnBothSidesOfAssignment()
func (c *CollapsePropertiesTest) TestPropertiesOnBothSidesOfAssignment() {
	// This verifies that replacements are done in the right order. Collapsing
	// the l-value in an assignment affects the parse tree immediately above
	// the r-value, so we update all rvalues before any lvalues.
	c.Test(
		c.Srcs("var a = {b: 0}; a.c = a.b;a.c = null"),
		c.Expected("var a$b = 0; var a$c = a$b;a$c = null"))
	c.TestSame("var a = {b: 0}; a.c = a?.b; a.c = null")
}

// Java: CollapsePropertiesTest.testCallOnUndefinedProperty()
func (c *CollapsePropertiesTest) TestCallOnUndefinedProperty() {
	// The "inherits" property is not explicitly defined on a.b anywhere, but
	// it is accessed as though it certainly exists (it is called), so we infer
	// that it must be an uncollapsible property that has come into existence
	// some other way.
	c.Test(
		c.Srcs("var a = {}; a.b = function(){}; a.b.inherits(x);"),
		c.Expected("var a$b = function(){}; a$b.inherits(x);"))
}

// Java: CollapsePropertiesTest.testGetPropOnUndefinedProperty()
func (c *CollapsePropertiesTest) TestGetPropOnUndefinedProperty() {
	// The "superClass_" property is not explicitly defined on a.b anywhere,
	// but it is accessed as though it certainly exists (a subproperty of it
	// is accessed), so we infer that it must be an uncollapsible property that
	// has come into existence some other way.
	c.Test(
		c.Srcs(
			`var a = {b: function(){}}; a.b.prototype.c =
function() { a.b.superClass_.c.call(this); }
`),
		c.Expected(
			`var a$b = function(){}; a$b.prototype.c =
function() { a$b.superClass_.c.call(this); }
`))
}

// Java: CollapsePropertiesTest.testNonWellformedAlias1()
func (c *CollapsePropertiesTest) TestNonWellformedAlias1() {
	c.TestSame("var a = {b: 3}; function f() { f(x); var x = a; f(x.b); }")
}

// Java: CollapsePropertiesTest.testNonWellformedAlias2()
func (c *CollapsePropertiesTest) TestNonWellformedAlias2() {
	c.TestSame(
		`var a = {b: 3};
function f() { if (false) { var x = a; f(x.b); } f(x); }
`)
}

// Java: CollapsePropertiesTest.testInlineAliasWithModifications()
func (c *CollapsePropertiesTest) TestInlineAliasWithModifications() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame("var x = 10; function f() { var y = x; x++; alert(y)} ")
	c.TestSame("var x = 10; function f() { var y = x; x+=1; alert(y)} ")
	c.Test(
		c.Srcs("var x = {}; x.x = 10; function f() {var y=x.x; x.x++; alert(y)}"),
		c.Expected("var x$x = 10; function f() {var y=x$x; x$x++; alert(y)}"))
	c.DisableNormalize()
	c.Test(
		c.Srcs("var x = {}; x.x = 10; function f() {var y=x.x; x.x+=1; alert(y)}"),
		c.Expected("var x$x = 10; function f() {var y=x$x; x$x+=1; alert(y)}"))
}

// Java: CollapsePropertiesTest.testDoNotCollapsePropertyOnExternType()
func (c *CollapsePropertiesTest) TestDoNotCollapsePropertyOnExternType() {
	c.TestSame("String.myFunc = function() {}; String.myFunc()")
}

// Java: CollapsePropertiesTest.testBug1704733()
func (c *CollapsePropertiesTest) TestBug1704733() {
	prelude := `function protect(x) { return x; }
function O() {}
protect(O).m1 = function() {};
protect(O).m2 = function() {};
protect(O).m3 = function() {};
`
	c.TestSame(prelude + "alert(O.m1); alert(O.m2()); alert(!O.m3);")
}

// Java: CollapsePropertiesTest.testBug1956277()
func (c *CollapsePropertiesTest) TestBug1956277() {
	c.Test(c.Srcs("var CONST = {}; CONST.URL = 3;"), c.Expected("var CONST$URL = 3;"))
}

// Java: CollapsePropertiesTest.testBug1974371()
func (c *CollapsePropertiesTest) TestBug1974371() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs("/** @enum {Object} */ var Foo = {A: {c: 2}, B: {c: 3}}; for (var key in Foo) {}"),
		c.Expected(
			`var Foo$A = {c: 2}; var Foo$B = {c: 3};
/** @enum {Object} */ var Foo = {A: Foo$A, B: Foo$B};
for (var key in Foo) {}
`))
}

// Java: CollapsePropertiesTest.testHasOwnProperty()
func (c *CollapsePropertiesTest) TestHasOwnProperty() {
	c.TestSame("var a = {b: 3}; if (a.hasOwnProperty(foo)) { alert('ok'); }")
	c.TestSame("var a = {b: 3}; if (a.hasOwnProperty(foo)) { alert('ok'); } a?.b;")
}

// Java: CollapsePropertiesTest.testHasOwnPropertyOnNonGlobalName()
func (c *CollapsePropertiesTest) TestHasOwnPropertyOnNonGlobalName() {
	c.TestSame(
		`/** @constructor */
function A() {
  this.foo = {a: 1, b: 1};
}
A.prototype.bar = function(prop) {
  return this.foo.hasOwnProperty(prop);
}
`)
	c.TestSame("var a = {'b': {'c': 1}}; a['b'].hasOwnProperty('c');")
	c.TestSame("var a = {b: 3}; if (Object.prototype.hasOwnProperty.call(a, 'b')) { alert('ok'); }")
}

// Java: CollapsePropertiesTest.testHasOwnPropertyNested()
func (c *CollapsePropertiesTest) TestHasOwnPropertyNested() {
	c.Test(
		c.Srcs("var a = {b: {c: 3}}; if (a.b.hasOwnProperty('c')) { alert('ok'); }"),
		c.Expected("var a$b =   {c: 3};  if (a$b.hasOwnProperty('c')) { alert('ok'); }"))
	c.Test(
		c.Srcs("var a = {b: {c: 3}}; if (a.b.hasOwnProperty('c')) { alert('ok'); } a.b.c;"),
		c.Expected("var a$b =   {c: 3};  if (a$b.hasOwnProperty('c')) { alert('ok'); } a$b.c;"))
	c.Test(
		c.Srcs("var a = {b: {c: 3}}; if (a.b.hasOwnProperty('c')) { alert('ok'); } a.b?.c;"),
		c.Expected("var a$b =   {c: 3};  if (a$b.hasOwnProperty('c')) { alert('ok'); } a$b?.c;"))
	c.Test(
		c.Srcs("var a = {}; a.b = function(p) { log(a.b.c.hasOwnProperty(p)); }; a.b.c = {};"),
		c.Expected("var a$b = function(p) { log(a$b$c.hasOwnProperty(p)); }; var a$b$c = {};"))
	c.Test(
		c.Srcs("var a = {}; a.b = function(p) { log(a.b.c?.hasOwnProperty(p)); };     a.b.c = {};"),
		c.Expected(
			"var         a$b = function(p) { log(a$b$c?.hasOwnProperty(p)); }; var a$b$c = {};"))
}

// Java: CollapsePropertiesTest.testHasOwnPropertyMultiple()
func (c *CollapsePropertiesTest) TestHasOwnPropertyMultiple() {
	c.TestSame("var a = {b: 3, c: 4, d: 5}; if (a.hasOwnProperty(prop)) { alert('ok'); }")
}

// Java: CollapsePropertiesTest.testObjectStaticMethodsPreventCollapsing()
func (c *CollapsePropertiesTest) TestObjectStaticMethodsPreventCollapsing() {
	c.TestSame("var a = {b: 3}; alert(Object.getOwnPropertyDescriptor(a, 'b'));")
	c.TestSame("var a = {b: 3}; alert(Object.getOwnPropertyDescriptors(a));")
	c.TestSame("var a = {b: 3}; alert(Object.getOwnPropertyNames(a));")
	c.TestSame("var a = {b: 3, [Symbol('c')]: 4}; alert(Object.getOwnPropertySymbols(a));")
}

// Java: CollapsePropertiesTest.testEnumOfObjects1()
func (c *CollapsePropertiesTest) TestEnumOfObjects1() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs(
			`/** @enum {Object} */
var Foo = {A: {c: 2}, B: {c: 3}};
for (var key in Foo.A) {}
`),
		c.Expected(
			`var Foo$A = {c: 2};
var Foo$B$c = 3;
for (var key in Foo$A) {}
`))
	c.Test(
		c.Srcs(
			`/** @enum {Object} */
var Foo = {A: {c: 2}, B: {c: 3}};
for (var key in Foo?.A) {}
`),
		c.Expected(
			`var Foo$A = {c: 2};
var Foo$B = {c: 3};
/** @enum {!Object} */
var Foo = {A:Foo$A, B:Foo$B};
for (var key in Foo?.A) {}
`))
}

// Java: CollapsePropertiesTest.testEnumOfObjects2()
func (c *CollapsePropertiesTest) TestEnumOfObjects2() {
	c.Test(
		c.Srcs(
			`/** @enum {Object} */ var Foo = {A: {c: 2}, B: {c: 3}};
foo(Foo.A.c);
`),
		c.Expected("var Foo$A$c = 2; var Foo$B$c = 3; foo(Foo$A$c);"))
	c.Test(
		c.Srcs(
			`/** @enum {Object} */
var Foo = {A: {c: 2}, B: {c: 3}};
foo?.(Foo.A?.c);
`),
		c.Expected(
			`var Foo$A = {c: 2};
var Foo$B$c = 3;
foo?.(Foo$A?.c);
`))
}

// Java: CollapsePropertiesTest.testEnumOfObjects3()
func (c *CollapsePropertiesTest) TestEnumOfObjects3() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs(
			`var x = {c: 2};
var y = {c: 3};
/** @enum {Object} */
var Foo = {A: x, B: y};
for (var key in Foo) {}
`),
		c.Expected(
			`var x = {c: 2};
var y = {c: 3};
var Foo$A = x;
var Foo$B = y;
/** @enum {Object} */
var Foo = {A: Foo$A, B: Foo$B};
for (var key in Foo) {}
`))
}

// Java: CollapsePropertiesTest.testEnumOfObjects4()
func (c *CollapsePropertiesTest) TestEnumOfObjects4() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	// Note that this produces bad code, but that's OK, because
	// checkConsts will yell at you for reassigning an enum value.
	// (enum values have to be constant).
	c.Test(
		c.Srcs(
			`/** @enum {Object} */ var Foo = {A: {c: 2}, B: {c: 3}};
for (var key in Foo) {} Foo.A = 3; alert(Foo.A);
`),
		c.Expected(
			`var Foo$A = {c: 2}; var Foo$B = {c: 3};
/** @enum {Object} */ var Foo = {A: Foo$A, B: Foo$B};
for (var key in Foo) {} Foo$A = 3; alert(Foo$A);
`))
	c.Test(
		c.Srcs(
			`/** @enum {Object} */ var Foo = {A: {c: 2}, B: {c: 3}};
for (var key in Foo) {} Foo.A = 3; alert?.(Foo?.A);
`),
		c.Expected(
			`var Foo$A = {c: 2}; var Foo$B = {c: 3};
/** @enum {Object} */ var Foo = {A: Foo$A, B: Foo$B};
var key;
for (key in Foo) {} Foo$A = 3; alert?.(Foo?.A);
`))
}

// Java: CollapsePropertiesTest.testObjectOfObjects1()
func (c *CollapsePropertiesTest) TestObjectOfObjects1() {
	// Basically the same as testEnumOfObjects4, but without the
	// constant enum values.
	c.TestSame(
		"var Foo = {a: {c: 2}, b: {c: 3}}; var key; for (key in Foo) {} Foo.a = 3; alert(Foo.a);")
}

// Java: CollapsePropertiesTest.testReferenceInAnonymousObject0()
func (c *CollapsePropertiesTest) TestReferenceInAnonymousObject0() {
	c.Test(
		c.Srcs(
			`var a = {};
a.b = function(){};
a.b.prototype.c = function(){};
var d = a.b.prototype.c;
`),
		c.Expected(
			`var a$b = function(){};
a$b.prototype.c = function(){};
var d = a$b.prototype.c;
`))
}

// Java: CollapsePropertiesTest.testReferenceInAnonymousObject1()
func (c *CollapsePropertiesTest) TestReferenceInAnonymousObject1() {
	c.Test(
		c.Srcs(
			`var a = {};
a.b = function(){};
var d = a.b.prototype.c;
`),
		c.Expected(
			`var a$b = function(){};
var d = a$b.prototype.c;
`))
}

// Java: CollapsePropertiesTest.testReferenceInAnonymousObject2()
func (c *CollapsePropertiesTest) TestReferenceInAnonymousObject2() {
	c.Test(
		c.Srcs(
			`var a = {};
a.b = function(){};
a.b.prototype.c = function(){};
var d = {c: a.b.prototype.c};
`),
		c.Expected(
			`var a$b = function(){};
a$b.prototype.c = function(){};
var d$c = a$b.prototype.c;
`))
}

// Java: CollapsePropertiesTest.testReferenceInAnonymousObject3()
func (c *CollapsePropertiesTest) TestReferenceInAnonymousObject3() {
	c.Test(
		c.Srcs(
			`function CreateClass(a$jscomp$1) {}
var a = {};
a.b = function(){};
a.b.prototype.c = function(){};
a.d = CreateClass({c: a.b.prototype.c});
`),
		c.Expected(
			`function CreateClass(a$jscomp$1) {}
var a$b = function(){};
a$b.prototype.c = function(){};
var a$d = CreateClass({c: a$b.prototype.c});
`))
}

// Java: CollapsePropertiesTest.testReferenceInAnonymousObject4()
func (c *CollapsePropertiesTest) TestReferenceInAnonymousObject4() {
	c.Test(
		c.Srcs(
			`function CreateClass(a) {}
var a = {};
a.b = CreateClass({c: function() {}});
a.d = CreateClass({c: a.b.c});
`),
		c.Expected(
			`function CreateClass(a$jscomp$1) {}
var a$b = CreateClass({c: function() {}});
var a$d = CreateClass({c: a$b.c});
`))
}

// Java: CollapsePropertiesTest.testReferenceInAnonymousObject5()
func (c *CollapsePropertiesTest) TestReferenceInAnonymousObject5() {
	c.Test(
		c.Srcs(
			`function CreateClass(a) {}
var a = {};
a.b = CreateClass({c: function() {}});
a.d = CreateClass({c: a.b.prototype.c});
`),
		c.Expected(
			`function CreateClass(a$jscomp$1) {}
var a$b = CreateClass({c: function() {}});
var a$d = CreateClass({c: a$b.prototype.c});
`))
}

// Java: CollapsePropertiesTest.testCrashInNestedAssign()
func (c *CollapsePropertiesTest) TestCrashInNestedAssign() {
	c.Test(
		c.Srcs("var a = {}; if (a.b = function() {}) a.b();"),
		c.Expected("var a$b; if (a$b=function() {}) { a$b(); }"))
}

// Java: CollapsePropertiesTest.testTwinReferenceCancelsChildCollapsing()
func (c *CollapsePropertiesTest) TestTwinReferenceCancelsChildCollapsing() {
	c.Test(
		c.Srcs("var a = {}; if (a.b = function() {}) { a.b.c = 3; a.b(a.b.c); }"),
		c.Expected("var a$b; if (a$b = function() {}) { a$b.c = 3; a$b(a$b.c); }"))
}

// Java: CollapsePropertiesTest.testPropWithDollarSign()
func (c *CollapsePropertiesTest) TestPropWithDollarSign() {
	c.Test(c.Srcs("var a = {$: 3}"), c.Expected("var a$$0 = 3;"))
}

// Java: CollapsePropertiesTest.testPropWithDollarSign2()
func (c *CollapsePropertiesTest) TestPropWithDollarSign2() {
	c.Test(c.Srcs("var a = {$: function(){}}"), c.Expected("var a$$0 = function(){};"))
}

// Java: CollapsePropertiesTest.testPropWithDollarSign3()
func (c *CollapsePropertiesTest) TestPropWithDollarSign3() {
	c.Test(
		c.Srcs("var a = {b: {c: 3}, b$c: function(){}}"),
		c.Expected("var a$b$c = 3; var a$b$0c = function(){};"))
}

// Java: CollapsePropertiesTest.testPropWithDollarSign4()
func (c *CollapsePropertiesTest) TestPropWithDollarSign4() {
	c.Test(c.Srcs("var a = {$$: {$$$: 3}};"), c.Expected("var a$$0$0$$0$0$0 = 3;"))
}

// Java: CollapsePropertiesTest.testPropWithDollarSign5()
func (c *CollapsePropertiesTest) TestPropWithDollarSign5() {
	c.Test(
		c.Srcs("var a = {b: {$0c: true}, b$0c: false};"),
		c.Expected("var a$b$$00c = true; var a$b$00c = false;"))
}

// Java: CollapsePropertiesTest.testConstKey()
func (c *CollapsePropertiesTest) TestConstKey() {
	c.Test(c.Srcs("var foo = {A: 3};"), c.Expected("var foo$A = 3;"))
}

// Java: CollapsePropertiesTest.testPropertyOnGlobalCtor()
func (c *CollapsePropertiesTest) TestPropertyOnGlobalCtor() {
	c.Test(
		c.Srcs("/** @constructor */ function Map() {} Map.foo = 3; Map;"),
		c.Expected("/** @constructor */ function Map() {} var Map$foo = 3; Map;"))
}

// Java: CollapsePropertiesTest.testPropertyOnGlobalInterface()
func (c *CollapsePropertiesTest) TestPropertyOnGlobalInterface() {
	c.Test(
		c.Srcs("/** @interface */ function Map() {} Map.foo = 3; Map;"),
		c.Expected("/** @interface */ function Map() {} var Map$foo = 3; Map;"))
}

// Java: CollapsePropertiesTest.testPropertyOnGlobalFunction()
func (c *CollapsePropertiesTest) TestPropertyOnGlobalFunction() {
	c.TestSame("function Map() {} Map.foo = 3; alert(Map);")
}

// Java: CollapsePropertiesTest.testIssue389()
func (c *CollapsePropertiesTest) TestIssue389() {
	c.TestSame(
		`function alias() {}
var dojo = {};
dojo.gfx = {};
dojo.declare = function() {};
/** @constructor */
dojo.gfx.Shape = function() {};
dojo.gfx.Shape = dojo.declare('dojo.gfx.Shape');
alias(dojo);
`,
		InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING)
}

// Java: CollapsePropertiesTest.testAliasedTopLevelName()
func (c *CollapsePropertiesTest) TestAliasedTopLevelName() {
	c.TestSame(
		`function alias() {}
var dojo = {};
dojo.gfx = {};
dojo.declare = function() {};
dojo.gfx.Shape = {SQUARE: 2};
dojo.gfx.Shape = dojo.declare('dojo.gfx.Shape');
alias(dojo);
alias(dojo$gfx$Shape$SQUARE);
`)
}

// Java: CollapsePropertiesTest.testAliasedTopLevelEnum()
func (c *CollapsePropertiesTest) TestAliasedTopLevelEnum() {
	c.TestSame(
		`function alias() {}
var dojo = {};
dojo.gfx = {};
dojo.declare = function() {};
/** @enum {number} */
dojo.gfx.Shape = {SQUARE: 2};
dojo.gfx.Shape = dojo.declare('dojo.gfx.Shape');
alias(dojo);
alias(dojo.gfx.Shape.SQUARE);
`,
		InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING)
}

// Java: CollapsePropertiesTest.testAssignFunctionBeforeDefinition()
func (c *CollapsePropertiesTest) TestAssignFunctionBeforeDefinition() {
	c.TestSame(
		`f = function() {};
var f = null;
`)
}

// Java: CollapsePropertiesTest.testObjectLitBeforeDefinition()
func (c *CollapsePropertiesTest) TestObjectLitBeforeDefinition() {
	c.TestSame(
		`a = {b: 3};
var a = null;
this.c = a.b;
`)
}

// Java: CollapsePropertiesTest.testTypedef1()
func (c *CollapsePropertiesTest) TestTypedef1() {
	c.Test(
		c.Srcs(
			`var foo = {};
/** @typedef {number} */ foo.Baz;
`),
		c.Expected("var foo = {}; var foo$Baz;"))
}

// Java: CollapsePropertiesTest.testTypedef2()
func (c *CollapsePropertiesTest) TestTypedef2() {
	c.Test(
		c.Srcs(
			`var foo = {};
/** @typedef {number} */ foo.Bar.Baz;
foo.Bar = function() {};
`),
		c.Expected("var foo$Bar$Baz; var foo$Bar = function(){};"))
}

// Java: CollapsePropertiesTest.testDelete1()
func (c *CollapsePropertiesTest) TestDelete1() {
	c.TestSame(
		`var foo = {};
foo.bar = 3;
delete foo.bar;
`)
}

// Java: CollapsePropertiesTest.testDelete2()
func (c *CollapsePropertiesTest) TestDelete2() {
	c.Test(
		c.Srcs(
			`var foo = {};
foo.bar = 3;
foo.baz = 3;
delete foo.bar;
`),
		c.Expected(
			`var foo = {};
foo.bar = 3;
var foo$baz = 3;
delete foo.bar;
`))
	c.TestSame(
		`var foo = {};
foo.bar = 3;
foo.baz = 3;
delete foo?.bar;
`)
}

// Java: CollapsePropertiesTest.testDelete3()
func (c *CollapsePropertiesTest) TestDelete3() {
	c.TestSame(
		`var foo = {bar: 3};
delete foo.bar;
`)
}

// Java: CollapsePropertiesTest.testDelete4()
func (c *CollapsePropertiesTest) TestDelete4() {
	c.Test(
		c.Srcs(
			`var foo = {bar: 3, baz: 3};
delete foo.bar;
`),
		c.Expected("var foo$baz=3;var foo={bar:3};delete foo.bar"))
	c.TestSame(
		`var foo = {bar: 3, baz: 3};
delete foo?.bar;
`)
}

// Java: CollapsePropertiesTest.testDelete5()
func (c *CollapsePropertiesTest) TestDelete5() {
	c.Test(
		c.Srcs(
			`var x = {};
x.foo = {};
x.foo.bar = 3;
delete x.foo.bar;
`),
		c.Expected(
			`var x$foo = {};
x$foo.bar = 3;
delete x$foo.bar;
`))
}

// Java: CollapsePropertiesTest.testDelete6()
func (c *CollapsePropertiesTest) TestDelete6() {
	c.Test(
		c.Srcs(
			`var x = {};
x.foo = {};
x.foo.bar = 3;
x.foo.baz = 3;
delete x.foo.bar;
`),
		c.Expected(
			`var x$foo = {};
x$foo.bar = 3;
var x$foo$baz = 3;
delete x$foo.bar;
`))
	c.Test(
		c.Srcs(
			`var x = {};
x.foo = {};
x.foo.bar = 3;
x.foo.baz = 3;
delete x.foo?.bar;
`),
		c.Expected(
			`var x$foo = {};
x$foo.bar = 3;
x$foo.baz = 3;
delete x$foo?.bar;
`))
}

// Java: CollapsePropertiesTest.testDelete7()
func (c *CollapsePropertiesTest) TestDelete7() {
	c.Test(
		c.Srcs(
			`var x = {};
x.foo = {bar: 3};
delete x.foo.bar;
`),
		c.Expected(
			`var x$foo = {bar: 3};
delete x$foo.bar;
`))
}

// Java: CollapsePropertiesTest.testDelete8()
func (c *CollapsePropertiesTest) TestDelete8() {
	c.Test(
		c.Srcs(
			`var x = {};
x.foo = {bar: 3, baz: 3};
delete x.foo.bar;
`),
		c.Expected(
			`var x$foo$baz = 3; var x$foo = {bar: 3};
delete x$foo.bar;
`))
	c.Test(
		c.Srcs(
			`var x = {};
x.foo = {bar: 3, baz: 3};
delete x.foo?.bar;
`),
		c.Expected(
			`var x$foo = {bar:3, baz:3};
delete x$foo?.bar;
`))
}

// Java: CollapsePropertiesTest.testDelete9()
func (c *CollapsePropertiesTest) TestDelete9() {
	c.TestSame(
		`var x = {};
x.foo = {};
x.foo.bar = 3;
delete x.foo;
`)
}

// Java: CollapsePropertiesTest.testDelete10()
func (c *CollapsePropertiesTest) TestDelete10() {
	c.TestSame(
		`var x = {};
x.foo = {bar: 3};
delete x.foo;
`)
}

// Java: CollapsePropertiesTest.testDelete11()
func (c *CollapsePropertiesTest) TestDelete11() {
	// Constructors are always collapsed.
	c.Test(
		`var x = {};
x.foo = {};
/** @constructor */ x.foo.Bar = function() {};
delete x.foo;
`,
		`var x = {};
x.foo = {};
/** @constructor */ var x$foo$Bar = function() {};
delete x.foo;
`,
		c.Warning(InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING))
	c.TestSame(
		c.Srcs(
			`var x = {};
x.foo = {};
/** @constructor */
x.foo.Bar = function() {};
delete x?.foo;
`),
		c.Warning(InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING))
}

// Java: CollapsePropertiesTest.testDelete11WithOptionalChain()
func (c *CollapsePropertiesTest) TestDelete11WithOptionalChain() {
	c.TestWarning(
		`var x = {};
x.foo = {};
/** @constructor */ x.foo.Bar = function() {};
// It's weird, but it is allowed to delete an optional chain.
delete x?.foo;
`,
		// The optional chain effectively creates an alias for `x`,
		// so we get a partial namespace warning instead of a namespase
		// redefined warning.
		InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING)
	c.TestWarning(
		`var x = {};
x.foo = {};
/** @constructor */ x.foo.Bar = function() {};
delete x.foo;
// The optional chain starts late enough that it doesn't generate
// a partial namespace warning.
x.foo.Bar?.baz
`,
		InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING)
}

// Java: CollapsePropertiesTest.testPreserveConstructorDoc()
func (c *CollapsePropertiesTest) TestPreserveConstructorDoc() {
	c.Test(
		c.Srcs("var foo = {}; /** @constructor */ foo.bar = function() {}"),
		c.Expected("/** @constructor */ var foo$bar = function() {}"))
}

// Java: CollapsePropertiesTest.testTypeDefAlias2()
func (c *CollapsePropertiesTest) TestTypeDefAlias2() {
	// TODO(johnlenz): make CollapseProperties safer around aliases of
	// functions and object literals.  Currently, this pass trades correctness
	// for code size.  We should able to create a safer compromise by teaching
	// the pass about goog.inherits and similiar calls.
	c.Test(
		c.Srcs(
			`/** @constructor */ var D = function() {};
/** @constructor */ D.L = function() {};
/** @type {D.L} */ D.L.A = new D.L();

/** @const */ var M = {};
if (random) { /** @typedef {D.L} */ M.L = D.L; }

use(M.L);
use(M.L.A);
`),
		c.Expected(
			`/** @constructor */ var D = function() {};
/** @constructor */ var D$L = function() {};
/** @type {D.L} */ var D$L$A = new D$L();
if (random) { /** @typedef {D.L} */ var M$L = D$L; }
use(M$L);
use(M$L.A);
`))
	c.Test(
		c.Srcs(
			`/** @constructor */ var D = function() {};
/** @constructor */ D.L = function() {};
/** @type {D.L} */ D.L.A = new D.L();

/** @const */ var M = {};
if (random) { /** @typedef {D.L} */ M.L = D?.L; }

use?.(M?.L);
use?.(M.L?.A);
`),
		c.Expected(
			`/** @constructor */ var D = function() {};
/** @constructor */ var D$L = function() {};
/** @type {D.L} */ var D$L$A = new D$L();
/** @const */ var M = {};
if (random) { /** @typedef {D.L} */ M.L = D?.L; }
// Similar to above, these optional chain references are broken by collapsing.
use?.(M?.L);
use?.(M.L?.A);
`))
}

// Java: CollapsePropertiesTest.testGlobalCatch()
func (c *CollapsePropertiesTest) TestGlobalCatch() {
	c.TestSame(
		`try {
  throw Error();
} catch (e) {
  console.log(e.name)
}
`)
}

// Java: CollapsePropertiesTest.testCtorManyAssignmentsDontInlineDontWarn()
func (c *CollapsePropertiesTest) TestCtorManyAssignmentsDontInlineDontWarn() {
	c.Test(
		c.Srcs(
			`var a = {};
/** @constructor */ a.b = function() {};
a.b.staticProp = 5;
function f(y, z) {
  var x = a.b;
  if (y) {
    x = z;
  }
  return new x();
}
`),
		c.Expected(
			`/** @constructor */
var a$b = function() {};
var a$b$staticProp = 5;
function f(y, z) {
  var x = a$b;
  if (y) {
    x = z;
  }
  return new x();
}
`))
}

// Java: CollapsePropertiesTest.testExpressionResultReferenceWontPreventCollapse()
func (c *CollapsePropertiesTest) TestExpressionResultReferenceWontPreventCollapse() {
	c.Test(
		c.Srcs(
			`var ns = {};
ns.Outer = {};

ns.Outer;
ns.Outer.Inner = function() {}
`),
		c.Expected(
			`var ns$Outer={};
ns$Outer;
var ns$Outer$Inner=function(){};
`))
}

// Java: CollapsePropertiesTest.testNoCollapseWithInvalidEnums()
func (c *CollapsePropertiesTest) TestNoCollapseWithInvalidEnums() {
	c.Test(
		c.Srcs(
			`/** @enum { { a: { b: number}} } */
var e = { KEY1: { a: { /** @nocollapse */ b: 123}},
  KEY2: { a: { b: 456}}
}
`),
		c.Expected("var e$KEY1$a={/** @nocollapse */ b:123}; var e$KEY2$a$b=456;"))
	c.Test(
		c.Srcs(
			`/** @enum */ var e = { A: 1, B: 2 };
/** @type {{ c: { d: number } }} */ e.name1 = {
  c: { /** @nocollapse */ d: 123 } };
`),
		c.Expected(
			"var e$A=1; var e$B=2; var e$name1$c={/** @nocollapse */ /** @nocollapse */ d:123};"))
	c.Test(
		c.Srcs(
			`/** @enum */ var e = { A: 1, B: 2};
/** @nocollapse */ e.foo = { bar: true };
`),
		c.Expected(
			`var e$A=1;
var e$B=2;
/** @enum */ var e = {};
/** @nocollapse */ e.foo = { bar: true }
`))
}

// Java: CollapsePropertiesTest.testDontCrashNamespaceAliasAcrossScopes()
func (c *CollapsePropertiesTest) TestDontCrashNamespaceAliasAcrossScopes() {
	c.Test(
		`var ns = {};
ns.VALUE = 0.01;
function f() {
    var constants = ns;
    (function() {
       var x = constants.VALUE;
    })();
}
`,
		(*string)(nil))
}

// Java: CollapsePropertiesTest.testCollapsedNameAlreadyTaken()
func (c *CollapsePropertiesTest) TestCollapsedNameAlreadyTaken() {
	c.Test(
		c.Srcs(
			`/** @constructor */ function Funny$Name(){};
function Funny(){};
Funny.Name = 5;
var x = new Funny$Name();
`),
		c.Expected(
			`/** @constructor */ function Funny$Name(){};
function Funny(){};
var Funny$Name$1 = 5;
var x = new Funny$Name();
`))
	c.Test(
		c.Srcs("var ns$x = 0; var ns$x$0 = 1; var ns = {}; ns.x = 8;"),
		c.Expected("var ns$x = 0; var ns$x$0 = 1; var ns$x$1 = 8;"))
	c.Test(
		c.Srcs("var ns$x = 0; var ns$x$0 = 1; var ns$x$1 = 2; var ns = {}; ns.x = 8;"),
		c.Expected("var ns$x = 0; var ns$x$0 = 1; var ns$x$1 = 2; var ns$x$2 = 8;"))
	c.Test(
		c.Srcs("var ns$x = {}; ns$x.y = 2; var ns = {}; ns.x = {}; ns.x.y = 8;"),
		c.Expected("var ns$x$y = 2; var ns$x$1$y = 8;"))
	c.Test(
		c.Srcs(
			`var ns$x = {};
ns$x.y = 1;
var ns = {};
ns.x$ = {};
ns.x$.y = 2;
ns.x = {};
ns.x.y = 3
`),
		c.Expected(
			`var ns$x$y = 1;
var ns$x$0$y = 2;
var ns$x$1$y = 3;
`))
}

// New Es6 Feature Tests - Some do not pass yet.
// Java: CollapsePropertiesTest.testArrowFunctionProperties()
func (c *CollapsePropertiesTest) TestArrowFunctionProperties() {
	// Add property to arrow function in local scope
	c.Test(
		c.Srcs("function a() {} () => { a.c = 5; return a.c; }"),
		c.Expected("function a() {} var a$c; () => { a$c = 5; return a$c; }"))
	c.TestSame("function a() {} () => { a.c = 5; return a?.c; }")
	// Reassign function property
	c.Test(
		c.Srcs("function a() {}; a.c = 1; () => { a.c = 5; }"),
		c.Expected("function a() {}; var a$c = 1; () => { a$c = 5; }"))
	// Arrow function assignment
	c.Test(
		c.Srcs("var a = () => {}; function f() { a.c = 5; }"),
		c.Expected("var a = () => {}; var a$c; function f() { a$c = 5; }"))
}

// Java: CollapsePropertiesTest.testDestructuredProperiesObjectLit()
func (c *CollapsePropertiesTest) TestDestructuredProperiesObjectLit() {
	// Using destructuring shorthand
	c.Test(
		c.Srcs("var { a, b } = { a: {}, b: {} }; a.a = 5; var c = a.a; "),
		c.Expected("var a; var b; ({a:a,b:b}={a:{},b:{}});a.a=5;var c=a.a"))
	// Without destructuring shorthand
	c.Test(
		c.Srcs("var { a:a, b:b } = { a: {}, b: {} }; a.a = 5; var c = a.a; "),
		c.Expected("var a; var b; ({a:a,b:b}={a:{},b:{}});a.a=5;var c=a.a"))
	// Test with greater depth
	c.Test(
		c.Srcs("var { a, b } = { a: {}, b: {} }; a.a.a = 5; var c = a.a; var d = c.a;"),
		c.Expected("var a; var b; ({a:a,b:b}={a:{},b:{}});a.a.a=5;var c=a.a;var d=c.a"))
}

// Java: CollapsePropertiesTest.testConditionalDestructuringAssignment()
func (c *CollapsePropertiesTest) TestConditionalDestructuringAssignment() {
	c.TestSame(
		`const X = { a: 1, b: 2 };
// Destructuring assignment effectively creates an alias for
// X and its properties.
// Make sure conditional assignment doesn't hide this.
const {a, b} = false || X;
console.log(a, b);
`)
}

// Java: CollapsePropertiesTest.testDestructuredClassAlias()
func (c *CollapsePropertiesTest) TestDestructuredClassAlias() {
	c.Test(
		`const ns = {}
ns.Foo = class {};
ns.Foo.STR = '';
const {Foo} = ns;
foo(Foo.STR);
`,
		`var ns$Foo = class {};
var ns$Foo$STR = '';
const Foo = null;
foo(ns$Foo$STR);
`)
}

// Java: CollapsePropertiesTest.testCanCollapseSinglePropertyInObjectPattern()
func (c *CollapsePropertiesTest) TestCanCollapseSinglePropertyInObjectPattern() {
	// Aggressive alias inlining replaces the reference to `y` with `x.y`
	// and sets `y` to `null`.
	c.Test(
		"const x   = {y: 1}; const {y} = x   ; use(  y);", //
		"var   x$y =     1 ; const  y  = null; use(x$y);")
}

// Java: CollapsePropertiesTest.testCanCollapseSinglePropertyInObjectPatternWithDefaultValue()
func (c *CollapsePropertiesTest) TestCanCollapseSinglePropertyInObjectPatternWithDefaultValue() {
	c.TestSame("const x = {y: 1}; const {y = 0} = x; use(y);")
}

// Java: CollapsePropertiesTest.testCollapsePropertyInNestedObjectPattern()
func (c *CollapsePropertiesTest) TestCollapsePropertyInNestedObjectPattern() {
	c.Test(
		"const x = {y: {z: 1}}; const {y: {z}} = x; use(z);",
		"var x$y$z = 1; const destructuring$m1146332801$0 = null; const z = null; use(x$y$z);")
}

// Java: CollapsePropertiesTest.testCanCollapseSinglePropertyInObjectPatternAssign()
func (c *CollapsePropertiesTest) TestCanCollapseSinglePropertyInObjectPatternAssign() {
	c.TestSame("const x = {y: 1}; var y; ({y} = x); use(y);")
}

// Java: CollapsePropertiesTest.testCanCollapseSinglePropertyInObjectPatternInForLoopClosure()
func (c *CollapsePropertiesTest) TestCanCollapseSinglePropertyInObjectPatternInForLoopClosure() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame("const x = {y: 1}; for (const {y} = x; true;) { use(() => y); }")
}

// Java: CollapsePropertiesTest.testPropertyInArray()
func (c *CollapsePropertiesTest) TestPropertyInArray() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame("var a, b = [{}, {}]; a.foo = 5; b.bar = 6;")
	c.Test(
		c.Srcs("var a = {}; a.b = 5; var c, d = [6, a.b]"),
		c.Expected("var a$b = 5; var c, d = [6, a$b];"))
}

// Java: CollapsePropertiesTest.testCollapsePropertySetInPattern()
func (c *CollapsePropertiesTest) TestCollapsePropertySetInPattern() {
	// TODO(b/120303257): collapse lvalues in destructuring patterns. We delayed implementing this
	// because it's uncommon to have properties as lvalues in destructuring patterns.
	c.Test(
		c.Srcs("var a = {}; a.b = {}; [a.b.c, a.b.d] = [1, 2];"),
		c.Expected("var a$b = {}; [a$b.c, a$b.d] = [1, 2];"))
	c.Test(
		c.Srcs("var a = {}; a.b = {}; ({x: a.b.c, y: a.b.d} = {});"),
		c.Expected("var a$b = {}; ({x: a$b.c, y: a$b.d} = {});"))
}

// Java: CollapsePropertiesTest.testDontCollapsePropertiesOfSpreadNamespace()
func (c *CollapsePropertiesTest) TestDontCollapsePropertiesOfSpreadNamespace() {
	// Notice we can still collapse the parent namespace, `a`.
	c.Test(
		c.Srcs("var a = {b: {c: 0}}; var d = {...a.b}; use(a.b.c);"),
		c.Expected("var a$b = {c: 0};var d = {...a$b}; use(a$b.c);"))
	// test transpiled version of spread.
	c.Test(
		"var a = {b: {c: 0}}; var d = Object.assign({}, a.b); use(a.b.c);",
		"var a$b = {c: 0}; var d = Object.assign({}, a$b); use(a$b.c);")
}

// Java: CollapsePropertiesTest.testDontCollapsePropertiesOfRestedNamespace()
func (c *CollapsePropertiesTest) TestDontCollapsePropertiesOfRestedNamespace() {
	c.Test(
		c.Srcs("var a = {b: {c: 0}}; var {...d} = a.b; use(a.b.c);"),
		c.Expected("var a$b = {c: 0}; var d; ({...d} = a$b); use(a$b.c);"))
	c.Test(
		"var a = {b: {c: 0}}; var {...d} = a?.b; use?.(a.b?.c);",
		"var a = {b: {c: 0}}; var d; ({...d} = a?.b); use?.(a.b?.c);")
	// "a.b.c" is not aliased by the REST in this case.
	c.Test(
		c.Srcs("var a = {b: {c: 0}}; var {d: {...d}} = a.b; use(a.b.c);"),
		c.Expected("var a$b = {c: 0}; var d; ({d: {...d}} = a$b); use(a$b.c);"))
	c.Test(
		"var a = {b: {c: 0}}; var {d: {...d}} = a?.b; use?.(a.b?.c);",
		"var a = {b: {c: 0}}; var d; ({d: {...d}} = a?.b); use?.(a.b?.c);")
}

// Java: CollapsePropertiesTest.testCollapsePropertiesWhenDeclaredInObjectLitWithSpread()
func (c *CollapsePropertiesTest) TestCollapsePropertiesWhenDeclaredInObjectLitWithSpread() {
	c.TestSame("var a = {b: 0, ...c}; use(a.b);")
	// transpiled spread
	c.TestSame("var a = Object.assign({}, {b: 0}, c); use(a.b);")
	c.TestSame("var a = {b: 0, ...c}; use?.(a?.b);")
	c.TestSame("var a = {...c}; use?.(a?.b);")
	c.Test(
		c.Srcs("    var a = {...c, b: 0};        use(a.b);"),
		c.Expected("var a$b = 0; var a = {...c}; use(a$b);"))
	// test transpiled - does not collapse properties after spread.
	c.TestSame("var a = Object.assign({}, c, {b: 0}); use(a.b);")
	c.TestSame(
		`var a = {...c, b: 0};
use?.(a?.b);
`)
}

// Java: CollapsePropertiesTest.testComputedPropertyNames()
func (c *CollapsePropertiesTest) TestComputedPropertyNames() {
	// Computed property in object literal. This following test code is bad style - it does not
	// follow the assumptions of the pass and thus produces the following output.
	c.Test(
		c.Srcs(
			`var a = {
  ['val' + ++i]: i,
  ['val' + ++i]: i
};
a.val1;
`),
		c.Expected(
			`var a = {
  ['val' + ++i]: i,
  ['val' + ++i]: i
};
var a$val1;
`))
	c.Test(
		c.Srcs("var a = { ['val']: i, ['val']: i }; a.val;"),
		c.Expected("var a = { ['val']: i, ['val']: i }; var a$val;"))
	// Computed property method name in class
	c.TestSame(
		`class Bar {
  constructor(){}
  ['f'+'oo']() {
    return 1
  }
}
var bar = new Bar()
bar.foo();
`)
	// Computed property method name in class - no concatenation
	c.TestSame(
		`class Bar {
  constructor(){}
  ['foo']() {
    return 1
  }
}
var bar = new Bar()
bar.foo();
`)
}

// Java: CollapsePropertiesTest.testClassGetSetMembers()
func (c *CollapsePropertiesTest) TestClassGetSetMembers() {
	// Get and set methods
	c.TestSame(
		`class Bar {
  constructor(x) {
    this.x = x;
  }
  get foo() {
    return this.x;
  }
  set foo(xIn) {
    x = xIn;
  }
}
var bar = new Bar(1);
bar.foo;
bar.foo = 2;
`)
}

// Java: CollapsePropertiesTest.testClassNonStaticMembers()
func (c *CollapsePropertiesTest) TestClassNonStaticMembers() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	// Call class method inside class scope
	c.TestSame(
		`function getA() {};
class Bar {
  constructor(){}
  getA() {
    return 1;
  }
  getB(x) {
    this.getA();
  }
}
`)
	// Call class method outside class scope
	c.TestSame(
		`class Bar {
  constructor(){}
  getB(x) {}
}
var too;
var too = new Bar();
too.getB(too);
`)
	// Non-static method
	c.TestSame(
		`class Bar {
  constructor(x){
    this.x = x;
  }
  double() {
    return x*2;
  }
}
var too;
var too = new Bar();
too.double(1);
`)
}

// Java: CollapsePropertiesTest.testClassDeclarationWithStaticMembers()
func (c *CollapsePropertiesTest) TestClassDeclarationWithStaticMembers() {
	c.Test(
		c.Srcs("class Bar { static double(n) { return n*2; } } Bar.double(1);"),
		c.Expected("var Bar$double = function(n) { return n * 2; }; class Bar {} Bar$double(1);"))
}

// Java: CollapsePropertiesTest.testClassAssignmentWithStaticMembers()
func (c *CollapsePropertiesTest) TestClassAssignmentWithStaticMembers() {
	c.Test(
		c.Srcs(
			`const ns = {};
ns.Bar = class {
  static double(n) { return n*2; }
}
ns.Bar.double(1)
`),
		c.Expected(
			`var ns$Bar$double = function(n) { return n * 2; }
var ns$Bar = class {}
ns$Bar$double(1);
`))
}

// Java: CollapsePropertiesTest.testClassWithMultipleStaticMembers()
func (c *CollapsePropertiesTest) TestClassWithMultipleStaticMembers() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs(
			`class Bar {
  static double(n) {
    return n*2
  }
  static triple(n) {
    return n*3;
  }
}
Bar.double(1);
Bar.triple(2);
`),
		c.Expected(
			`var Bar$double = function(n) { return n * 2; }
var Bar$triple = function(n) { return n * 3; }
class Bar {}
Bar$double(1);
Bar$triple(2);
`))
}

// Java: CollapsePropertiesTest.testClassStaticMembersWithAdditionalUndeclaredProperty()
func (c *CollapsePropertiesTest) TestClassStaticMembersWithAdditionalUndeclaredProperty() {
	c.Test(
		c.Srcs(
			`class Bar {
  static double(n) {
    return n*2
  }
}
use(Bar.double.trouble);
`),
		c.Expected(
			`var Bar$double = function(n) { return n * 2; }
class Bar {}
use(Bar$double.trouble);
`))
}

// Java: CollapsePropertiesTest.testClassStaticMembersWithAdditionalPropertySetLocally()
func (c *CollapsePropertiesTest) TestClassStaticMembersWithAdditionalPropertySetLocally() {
	c.Test(
		c.Srcs(
			`class Bar {
  static double(n) {
    return n*2
  }
}
(function() { Bar.double.trouble = -1; })();
use(Bar.double.trouble);
`),
		c.Expected(
			`var Bar$double = function(n) { return n * 2; }
class Bar {}
var Bar$double$trouble;
(function() { Bar$double$trouble = -1; })();
use(Bar$double$trouble);
`))
	c.Test(
		c.Srcs(
			`class Bar {
  static double(n) {
    return n*2
  }
}
(function() { Bar.double.trouble = -1; })();
use?.(Bar.double?.trouble);
`),
		c.Expected(
			`var Bar$double = function(n) { return n * 2; };
class Bar {}
(function() { Bar$double.trouble = -1; })();
use?.(Bar$double?.trouble);
`))
}

// Java: CollapsePropertiesTest.testClassStaticMembersWithAdditionalPropertySetGlobally()
func (c *CollapsePropertiesTest) TestClassStaticMembersWithAdditionalPropertySetGlobally() {
	c.Test(
		c.Srcs(
			`class Bar {
  static double(n) {
    return n*2
  }
}
Bar.double.trouble = -1;
use(Bar.double.trouble);
`),
		c.Expected(
			`var Bar$double = function(n) { return n * 2; }
class Bar {}
var Bar$double$trouble = -1;
use(Bar$double$trouble);
`))
	c.Test(
		c.Srcs(
			`class Bar {
  static double(n) {
    return n*2
  }
}
Bar.double.trouble = -1;
use?.(Bar.double?.trouble);
`),
		c.Expected(
			`var Bar$double = function(n) {
  return n * 2;
}
class Bar {
}
Bar$double.trouble = -1;
use?.(Bar$double?.trouble);
`))
}

// Java: CollapsePropertiesTest.testClassDeclarationWithEscapedStaticMember_collapsesMemberButNotMemberProp()
func (c *CollapsePropertiesTest) TestClassDeclarationWithEscapedStaticMember_collapsesMemberButNotMemberProp() {
	c.Test(
		c.Srcs(
			`class Bar {
  static double(n) {
    return n*2
  }
}
Bar.double.trouble = 4;
use(Bar.double);
use(Bar.double.trouble);
`),
		c.Expected(
			`var Bar$double = function(n) { return n * 2; }
class Bar {}
// don't collapse this property because Bar$double is escaped
Bar$double.trouble = 4;
use(Bar$double);
use(Bar$double.trouble);
`))
	c.Test(
		c.Srcs(
			`class Bar {
  static double(n) {
    return n*2
  }
}
Bar.double.trouble = 4;
use?.(Bar?.double);
use?.(Bar.double?.trouble);
`),
		c.Expected(
			"var Bar$double = function(n) {\n"+
				"  return n * 2;\n"+
				"};\n"+
				"class Bar {\n"+
				"}\n"+
				"Bar$double.trouble = 4;\n"+
				"// TODO(b/148237949): collapsing breaks `Bar?.double` reference\n"+
				"use?.(Bar?.double);\n"+
				"use?.(Bar$double?.trouble);\n"))
}

// Java: CollapsePropertiesTest.testClassDeclarationWithStaticMembersWithNoCollapse()
func (c *CollapsePropertiesTest) TestClassDeclarationWithStaticMembersWithNoCollapse() {
	c.TestSame(
		`class Bar {
  /** @nocollapse */
  static double(n) {
    return n*2
  }
}
Bar.double(1);
`)
}

// Java: CollapsePropertiesTest.testClassStaticMemberWithUnreferencedInnerName()
func (c *CollapsePropertiesTest) TestClassStaticMemberWithUnreferencedInnerName() {
	c.Test(
		c.Srcs(
			`const Bar = class BarInternal {
  static baz(n) {
    use(n);
  }
}
Bar.baz();
`),
		c.Expected(
			`var Bar$baz = function(n) {
  use(n);
};
const Bar = class {};
Bar$baz();
`))
}

// Java: CollapsePropertiesTest.testDontCollapseClassStaticMemberReferencingInnerName()
func (c *CollapsePropertiesTest) TestDontCollapseClassStaticMemberReferencingInnerName() {
	c.Test(
		`const Bar = class BarInternal {
  static baz(n) {
    use(BarInternal);
  }
}
Bar.baz();
`,
		`var Bar$baz = function(n) {
  use(Bar);
};
const Bar = class {};
Bar$baz();
`)
}

// Java: CollapsePropertiesTest.testDontCollapseLateAddedClassStaticMemberReferencedByInnerName()
func (c *CollapsePropertiesTest) TestDontCollapseLateAddedClassStaticMemberReferencedByInnerName() {
	c.Test(
		c.Srcs(
			`const Bar = class BarInternal {
  method() {
    return BarInternal.Enum.E1; // ref via inner name
  }
}
/** @enum {number} */
Bar.Enum = { E1: 1 }; // added after class definition
`),
		c.Expected(
			`const Bar = class {
  method() {
    return Bar$Enum$E1;
  }
};
var Bar$Enum$E1 = 1;
`))
}

// Java: CollapsePropertiesTest.testDontCollapseClassStaticMemberReferencingInnerNameInNestedFunction()
func (c *CollapsePropertiesTest) TestDontCollapseClassStaticMemberReferencingInnerNameInNestedFunction() {
	c.Test(
		`const Bar = class BarInternal {
  static baz(n) {
    return function() { use(BarInternal); };
  }
}
Bar.baz();
`,
		`var Bar$baz = function(n) {
  return function() {
    use(Bar);
  };
};
const Bar = class {};
Bar$baz();
`)
}

// Java: CollapsePropertiesTest.testClassStaticMemberUsingSuperNotCollapsed()
func (c *CollapsePropertiesTest) TestClassStaticMemberUsingSuperNotCollapsed() {
	c.Test(
		`class Baz extends Bar() {
  /** @nocollapse */
  static quadruple(n) {
    return 2 * super.double(n);
  }
}
`,
		`const CLASS_EXTENDS$0 = Bar();
class Baz extends CLASS_EXTENDS$0 {
  /** @nocollapse */
  static quadruple(n) {
    return 2 * CLASS_EXTENDS$0.double(n);
  }
}
`)
}

// Java: CollapsePropertiesTest.testClassStaticMemberUsingSuperInArrowFnNotCollapsed()
func (c *CollapsePropertiesTest) TestClassStaticMemberUsingSuperInArrowFnNotCollapsed() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		`class Baz extends Bar() {
  /** @nocollapse */
  static quadruple(n) {
    return () => 2 * super.double(n);
 }
}
`,
		`const CLASS_EXTENDS$0 = Bar();
class Baz extends CLASS_EXTENDS$0 {
  /** @nocollapse */
  static quadruple(n) {
    return () => {
      return 2 * CLASS_EXTENDS$0.double(n);
    };
  }
}
`)
}

// Java: CollapsePropertiesTest.testClassStaticMemberWithInnerClassUsingSuperIsCollapsed()
func (c *CollapsePropertiesTest) TestClassStaticMemberWithInnerClassUsingSuperIsCollapsed() {
	c.Test(
		c.Srcs(
			`class Bar extends fn() {
  static m() {
    class Inner extends fn() {
      static n() { super.n(); }  // does not block collapsing Bar.m
    }
  }
}
`),
		c.Expected(
			`const CLASS_EXTENDS$0 = fn();
var Bar$m = function() {
  const CLASS_EXTENDS$1 = fn();
  class Inner extends CLASS_EXTENDS$1 {
    static n() {
      CLASS_EXTENDS$1.n();
    }
  }
};
class Bar extends CLASS_EXTENDS$0 {}
`))
}

// Java: CollapsePropertiesTest.testEs6StaticMemberAddedAfterDefinition()
func (c *CollapsePropertiesTest) TestEs6StaticMemberAddedAfterDefinition() {
	c.Test(
		c.Srcs(
			`class Bar {}
Bar.double = function(n) {
  return n*2
}
Bar.double(1);
`),
		c.Expected(
			`class Bar {}
var Bar$double = function(n) {
  return n*2
}
Bar$double(1);
`))
}

// Java: CollapsePropertiesTest.testEs6StaticMemberOnEscapedClassIsCollapsed()
func (c *CollapsePropertiesTest) TestEs6StaticMemberOnEscapedClassIsCollapsed() {
	c.Test(
		c.Srcs("class Bar { static m() {} } use(Bar);"),
		c.Expected("var Bar$m = function() {}; class Bar {} use(Bar);"))
}

// Java: CollapsePropertiesTest.testClassStaticProperties()
func (c *CollapsePropertiesTest) TestClassStaticProperties() {
	c.Test(
		c.Srcs("class A {} A.foo = 'bar'; use(A.foo);"),
		c.Expected("class A {} var A$foo = 'bar'; use(A$foo);"))
	c.Test(
		c.Srcs("class A {}     A.foo = 'bar'; use?.(A?.foo);"),
		c.Expected("class A {} var A$foo = 'bar'; use?.(A?.foo);"))
	// Collapsing A.foo is known to be unsafe.
	c.Test(
		c.Srcs("class A { static useFoo() { alert(this.foo); } } A.foo = 'bar'; A.useFoo();"),
		c.Expected(
			`var A$useFoo = function() { alert(this.foo); };
class A {}
var A$foo = 'bar';
A$useFoo();
`),
		c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
	c.Test(
		c.Srcs(
			`class A {
  static useFoo() {
    alert(this.foo);
  }
}
/** @nocollapse */
A.foo = 'bar';
A.useFoo();
`),
		c.Expected(
			`var A$useFoo = function() { alert(this.foo); };
class A {}
/** @nocollapse */
A.foo = 'bar';
A$useFoo();
`),
		c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
}

// Java: CollapsePropertiesTest.testClassStaticAndPrototypePropWithSameName()
func (c *CollapsePropertiesTest) TestClassStaticAndPrototypePropWithSameName() {
	c.Test(
		c.Srcs("const ns = {}; ns.C = class { x() {} }; ns.C.x = 3;"),
		c.Expected("var ns$C = class { x() {} }; var ns$C$x = 3;"))
}

// Java: CollapsePropertiesTest.testClassStaticProperties_locallyDeclared1()
func (c *CollapsePropertiesTest) TestClassStaticProperties_locallyDeclared1() {
	c.Test(
		c.Srcs(
			`class A {}
function addStaticPropToA() {
  A.staticProp = 5;
}
if (A.staticProp) {
  use(A.staticProp);
}
`),
		c.Expected(
			`class A {}
var A$staticProp;
function addStaticPropToA() {
  A$staticProp = 5;
}
if (A$staticProp) {
  use(A$staticProp);
}
`))
	c.Test(
		c.Srcs(
			`class A {}
function addStaticPropToA() {
  A.staticProp = 5;
}
if (A?.staticProp) {
  use?.(A?.staticProp);
}
`),
		c.Expected(
			"class A {}\n"+
				"var A$staticProp;\n"+
				"function addStaticPropToA() {\n"+
				"  A$staticProp = 5;\n"+
				"}\n"+
				"// TODO(b/148237949): collapsing breaks `A?.staticProp` reference\n"+
				"if (A?.staticProp) {\n"+
				"  use?.(A?.staticProp);\n"+
				"}\n"))
}

// Java: CollapsePropertiesTest.testClassStaticProperties_locallyDeclared2()
func (c *CollapsePropertiesTest) TestClassStaticProperties_locallyDeclared2() {
	c.Test(
		c.Srcs(
			`const A = class {}
function addStaticPropToA() {
  A.staticProp = 5;
}
if (A.staticProp) {
  use(A.staticProp);
}
`),
		c.Expected(
			`const A = class {}
var A$staticProp;
function addStaticPropToA() {
  A$staticProp = 5;
}
if (A$staticProp) {
  use(A$staticProp);
}
`))
	c.Test(
		c.Srcs(
			`const A = class {}
function addStaticPropToA() {
  A.staticProp = 5;
}
if (A?.staticProp) {
  use?.(A?.staticProp);
}
`),
		c.Expected(
			"const A = class {}\n"+
				"var A$staticProp;\n"+
				"function addStaticPropToA() {\n"+
				"  A$staticProp = 5;\n"+
				"}\n"+
				"// TODO(b/148237949): collapsing breaks `A?.staticProp` reference\n"+
				"if (A?.staticProp) {\n"+
				"  use?.(A?.staticProp);\n"+
				"}\n"))
}

// Java: CollapsePropertiesTest.testEs6ClassStaticInheritance()
func (c *CollapsePropertiesTest) TestEs6ClassStaticInheritance() {
	c.Test(
		c.Srcs("class A {} A.foo = 5; use(A.foo); class B extends A {}"),
		c.Expected("class A {} var A$foo = 5; use(A$foo); class B extends A {}"))
	c.Test(
		c.Srcs("class A {}     A.foo = 5; use?.(A?.foo); class B extends A {}"),
		c.Expected("class A {} var A$foo = 5; use?.(A?.foo); class B extends A {}"))
	// We potentially collapse unsafely when the subclass accesses a static property on its
	// superclass. However, AggressiveInlineAliases tries to rewrite inherited accesses to make this
	// collapsing safe. See InlineAndCollapsePropertiesTests for examples.
	c.Test(
		c.Srcs("    class A {}     A.foo = 5; use(A.foo); class B extends A {} use(B.foo);"),
		c.Expected("class A {} var A$foo = 5; use(A$foo); class B extends A {} use(A$foo);"))
	c.Test(
		c.Srcs("class A {}     A.foo = 5; use?.(A?.foo); class B extends A {} use?.(B?.foo);"),
		c.Expected("class A {} var A$foo = 5; use?.(A?.foo); class B extends A {} use?.(B?.foo);"))
	c.Test(
		c.Srcs(
			`class A {}
A.foo = 5;
class B extends A {}
use(B.foo);
B.foo = 6;
use(B.foo)
`),
		c.Expected(
			`class A {}
var A$foo = 5;
class B extends A {}
use(B$foo);
var B$foo = 6;
use(B$foo)
`))
	c.Test(
		c.Srcs(
			`class A {}
A.foo = 5;
class B extends A {}
use(B?.foo);
B.foo = 6;
use(B?.foo)
`),
		c.Expected(
			`class A {}
var A$foo = 5;
class B extends A {}
use(B?.foo);
var B$foo = 6;
use(B?.foo)
`))
	c.TestSame(
		`class A {}
/** @nocollapse */
A.foo = 5;
class B extends A {}
use(B.foo);
/** @nocollapse */
B.foo = 6;
use(B.foo);
`)
}

// Java: CollapsePropertiesTest.testEs6ClassExtendsChildClass()
func (c *CollapsePropertiesTest) TestEs6ClassExtendsChildClass() {
	c.Test(
		c.Srcs(
			`class Thing {}
Thing.Builder = class {};
class Subthing {}
Subthing.Builder = class extends Thing.Builder {}
`),
		c.Expected(
			`class Thing {}
var Thing$Builder = class {};
class Subthing {}
var Subthing$Builder = class extends Thing$Builder {}
`))
}

// Java: CollapsePropertiesTest.testSuperExtern()
func (c *CollapsePropertiesTest) TestSuperExtern() {
	c.TestSame(
		`class Foo {
  constructor(){
    this.x = x;
  }
  getX() {
    return this.x;
  }
}
class Bar extends Foo {
  constructor(x, y) {
    super(x);
    this.y = y;
  }
  getX() {
    return super.getX() + this.y;
  }
}
let too = new Bar();
too.getX();
`)
}

// Java: CollapsePropertiesTest.testPropertyMethodAssignment_receiverAffectedByCollapse()
func (c *CollapsePropertiesTest) TestPropertyMethodAssignment_receiverAffectedByCollapse() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	// ES5 version
	c.Test(
		`var foo = {
  bar: 1,
  myFunc: function myFunc() {
    return this.bar
  }
};
foo.myFunc();
`,
		`var foo$bar = 1;
var foo$myFunc = function myFunc() {
  return this.bar
};
foo$myFunc();
`,
		c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
	c.Test(
		c.Srcs(
			`var foo = {
  bar: 1,
  myFunc: function myFunc() {
    function inner() {
      return this.bar;
    }
  }
};
foo.myFunc();
`),
		c.Expected(
			`var foo$bar = 1;
var foo$myFunc = function myFunc() {
  function inner() {
    return this.bar;
  }
};
foo$myFunc();
`))
	// ES6 version
	c.Test(
		`var foo = {
  bar: 1,
  myFunc() {
    return this.bar;
  }
};
foo.myFunc();
`,
		`var foo$bar = 1;
var foo$myFunc = function() {
  return this.bar
};
foo$myFunc();
`,
		c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
	// "this" is lexically scoped in arrow functions so collapsing is safe.
	c.Test(
		c.Srcs(
			`var foo = {
  myFunc: () => {
    return this;
  }
};
foo.myFunc();
`),
		c.Expected(
			`var foo$myFunc = () => {
  return this
};
foo$myFunc();
`))
	// but not in inner functions
	c.Test(
		`var foo = {
  myFunc: function() {
    return (() => this);
  }
};
foo.myFunc();
`,
		`var foo$myFunc = function() {
  return (() => this);
};
foo$myFunc();
`,
		c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
	// `super` is also a receiver reference
	// TODO(b/122665204): Make this case pass rather than crash.
	AssertThrows[jutil.IRuntimeException](func() {
		c.Test(
			`var foo = {
  bar: 1,
  myFunc() {
    return super.bar;
  }
};
foo.myFunc();
`,
			`var foo$bar = 1;
var foo$myFunc = function() {
  return super.bar // Unclear what goes here.
};
foo$myFunc();
`,
			c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
	})
	// check references in param lists
	c.Test(
		`var foo = {
  bar: 1,
  myFunc(x = this) {
  }
};
foo.myFunc();
`,
		`var foo$bar = 1;
var foo$myFunc = function(x = this) {
};
foo$myFunc();
`,
		c.Warning(InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE))
}

// Java: CollapsePropertiesTest.testPropertyMethodAssignment_noThis()
func (c *CollapsePropertiesTest) TestPropertyMethodAssignment_noThis() {
	// ES5 Version
	c.Test(
		c.Srcs(
			`var foo = {
  bar: 1,
  myFunc: function myFunc() {
    return 5
  }
};
foo.myFunc();
`),
		c.Expected(
			`var foo$bar = 1;
var foo$myFunc = function myFunc() {
    return 5;
};
foo$myFunc();
`))
	// ES6 version
	c.Test(
		c.Srcs(
			`var foo = {
  bar: 1,
  myFunc() {
    return 5;
  }
};
foo.myFunc();
`),
		c.Expected(
			`var foo$bar = 1;
var foo$myFunc = function() {
    return 5;
};
foo$myFunc();
`))
}

// Java: CollapsePropertiesTest.testMethodPropertyShorthand()
func (c *CollapsePropertiesTest) TestMethodPropertyShorthand() {
	c.Test(
		c.Srcs(
			`var foo = {
  bar: 1,
   myFunc() {
    return 2
  }
};
foo.myFunc();
`),
		c.Expected(
			`var foo$bar = 1;
var foo$myFunc = function() {
    return 2;
};
foo$myFunc();
`))
}

// Java: CollapsePropertiesTest.testLetConstObjectAssignmentProperties()
func (c *CollapsePropertiesTest) TestLetConstObjectAssignmentProperties() {
	// All qualified names - even for variables that are initially declared as LETS and CONSTS -
	// are being declared as VAR statements, but this is correct because we are only
	// collapsing for global names.
	c.Test(
		c.Srcs("let a = {}; a.b = {}; a.b.c = {}; let d = 1; d = a.b.c;"),
		c.Expected("var a$b$c = {}; let d = 1; d = a$b$c;"))
	c.Test(
		c.Srcs("let a = {}; a.b = {}; a.b.c = {}; let d = 1; d = a.b?.c;"),
		c.Expected("var         a$b = {}; a$b.c = {}; let d = 1; d = a$b?.c;"))
	c.Test(c.Srcs("let a = {}; if(1)  { a.b = 1; }"), c.Expected("if(1) { var a$b = 1; }"))
	c.TestSame("if(1) { let a = {}; a.b = 1; }")
	c.Test(
		c.Srcs("var a = {}; a.b = 1; if(1) { let a = {}; a.b = 2; }"),
		c.Expected("var a$b = 1; if(1) { let a$jscomp$1 = {}; a$jscomp$1.b = 2; }"))
}

// Java: CollapsePropertiesTest.testTemplateStrings()
func (c *CollapsePropertiesTest) TestTemplateStrings() {
	c.Test(
		c.Srcs(
			"var a = {};\n"+
				"a.b = 'foo';\n"+
				"var c = `Hi ${a.b}`;\n"),
		c.Expected(
			"var a$b = 'foo';\n"+
				"var c = `Hi ${a$b}`;\n"))
}

// Java: CollapsePropertiesTest.testDoesNotCollapseInEs6ModuleScope()
func (c *CollapsePropertiesTest) TestDoesNotCollapseInEs6ModuleScope() {
	c.TestSame("var a = {}; a.b = {}; a.b.c = 5; export default function() {};")
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test(
		c.Srcs("import * as a from './a.js'; let b = a.b;", "var a = {}; a.b = 5;"),
		c.Expected("import * as a$jscomp$1 from './a.js'; let b = a$jscomp$1.b;", "var a$b = 5;"))
}

// Java: CollapsePropertiesTest.testDefaultParameters()
func (c *CollapsePropertiesTest) TestDefaultParameters() {
	c.TestSame("var a = {b: 5}; function f(x=a) { alert(x.b); }")
	c.TestSame("var a = {b: 5}; function f(x=a) { alert?.(x?.b); }")
	c.Test(
		c.Srcs("var a   = {b: {c: 5}}; function f(x=a.b) { alert(x.c); }"),
		c.Expected("var a$b =     {c: 5} ; function f(x=a$b) { alert(x.c); }"))
	c.TestSame("var a = {b: {c: 5}}; function f(x=a?.b) { alert?.(x?.c); }")
	c.Test(
		c.Srcs("var a = {b: 5}; function f(x=a.b) { alert(x); }"),
		c.Expected("var a$b = 5; function f(x=a$b) { alert(x); }"))
	c.TestSame("var a = {b: 5}; function f(x=a?.b) { alert?.(x); }")
}

// Java: CollapsePropertiesTest.testModuleExportsBasicCommonJs()
func (c *CollapsePropertiesTest) TestModuleExportsBasicCommonJs() {
	c.SetupModuleExportsOnly()
	inputs := []*SourceFile{}
	inputs = append(inputs, SourceFile_FromCode("mod1.js", "module.exports = 123;"))
	inputs = append(inputs, SourceFile_FromCode("entry.js", "var mod = require('./mod1.js'); alert(mod);"))
	expected := []*SourceFile{}
	expected = append(expected,
		SourceFile_FromCode(
			"mod1.js",
			"/** @const */ var module$mod1={}; /** @const */ var module$mod1$default = 123;"))
	expected = append(expected,
		SourceFile_FromCode(
			"entry.js", "var mod = module$mod1$default; alert(module$mod1$default);"))
	c.Test(c.Srcs(inputs), c.Expected(expected))
}

// Java: CollapsePropertiesTest.testModuleExportsBasicEsm()
func (c *CollapsePropertiesTest) TestModuleExportsBasicEsm() {
	c.SetupModuleExportsOnly()
	inputs := []*SourceFile{}
	inputs = append(inputs, SourceFile_FromCode("mod1.js", "export default 123; export var bar = 'bar'"))
	inputs = append(inputs,
		SourceFile_FromCode(
			"entry.js", "import mod, {bar} from './mod1.js'; alert(mod); alert(bar)"))
	expected := []*SourceFile{}
	expected = append(expected,
		SourceFile_FromCode(
			"mod1.js",
			`var $jscompDefaultExport$$module$mod1=123;
var bar$$module$mod1 = 'bar';
/** @const */ var module$mod1={};
/** @const */ var module$mod1$bar = bar$$module$mod1
/** @const */ var module$mod1$default = $jscompDefaultExport$$module$mod1;
`))
	expected = append(expected,
		SourceFile_FromCode(
			"entry.js", "alert($jscompDefaultExport$$module$mod1); alert(bar$$module$mod1);"))
	c.Test(c.Srcs(inputs), c.Expected(expected))
}

// Java: CollapsePropertiesTest.testMutableModuleExportsBasicEsm()
func (c *CollapsePropertiesTest) TestMutableModuleExportsBasicEsm() {
	c.SetupModuleExportsOnly()
	inputs := []*SourceFile{}
	inputs = append(inputs,
		SourceFile_FromCode(
			"mod1.js",
			`export default 123;
export var bar = 'ba';
function f() { bar += 'r'; }
`))
	inputs = append(inputs,
		SourceFile_FromCode(
			"entry.js", "import mod, {bar} from './mod1.js'; alert(mod); alert(bar)"))
	expected := []*SourceFile{}
	expected = append(expected,
		SourceFile_FromCode(
			"mod1.js",
			`var $jscompDefaultExport$$module$mod1=123;
var bar$$module$mod1 = 'ba';
function f$$module$mod1() { bar$$module$mod1 += 'r'; }
/** @const */ var module$mod1= {
  /** @return {?} */ get bar() { return bar$$module$mod1; },
};
/** @const */ var module$mod1$default = $jscompDefaultExport$$module$mod1;
`))
	expected = append(expected,
		SourceFile_FromCode(
			"entry.js",
			`alert($jscompDefaultExport$$module$mod1);
alert(bar$$module$mod1);
`))
	c.Test(c.Srcs(inputs), c.Expected(expected))
}

// Java: CollapsePropertiesTest.testModuleExportsObjectCommonJs()
func (c *CollapsePropertiesTest) TestModuleExportsObjectCommonJs() {
	c.SetupModuleExportsOnly()
	inputs := []*SourceFile{}
	inputs = append(inputs,
		SourceFile_FromCode(
			"mod1.js",
			`var foo ={};
foo.bar = {};
foo.bar.baz = 123;
module.exports = foo;
`))
	inputs = append(inputs,
		SourceFile_FromCode("entry.js", "var mod = require('./mod1.js'); alert(mod.bar.baz);"))
	expected := []*SourceFile{}
	expected = append(expected,
		SourceFile_FromCode(
			"mod1.js",
			`/** @const */ var module$mod1={};
 /** @const */ var module$mod1$default = {};
module$mod1$default.bar = {};
module$mod1$default.bar.baz = 123;
`))
	expected = append(expected,
		SourceFile_FromCode(
			"entry.js", "var mod = module$mod1$default; alert(module$mod1$default.bar.baz);"))
	c.Test(c.Srcs(inputs), c.Expected(expected))
}

// Java: CollapsePropertiesTest.testModuleExportsObjectEsm()
func (c *CollapsePropertiesTest) TestModuleExportsObjectEsm() {
	c.SetupModuleExportsOnly()
	inputs := []*SourceFile{}
	inputs = append(inputs,
		SourceFile_FromCode(
			"mod1.js",
			`var foo ={};
foo.bar = {};
foo.bar.baz = 123;
export default foo;
`))
	inputs = append(inputs, SourceFile_FromCode("entry.js", "import mod from './mod1.js'; alert(mod.bar.baz);"))
	expected := []*SourceFile{}
	expected = append(expected,
		SourceFile_FromCode(
			"mod1.js",
			`var foo$$module$mod1 = {};
foo$$module$mod1.bar = {};
foo$$module$mod1.bar.baz = 123;
var $jscompDefaultExport$$module$mod1 = foo$$module$mod1;
/** @const */ var module$mod1={};
/** @const */ var module$mod1$default = $jscompDefaultExport$$module$mod1;
`))
	expected = append(expected,
		SourceFile_FromCode("entry.js", "alert($jscompDefaultExport$$module$mod1.bar.baz);"))
	c.Test(c.Srcs(inputs), c.Expected(expected))
}

// Java: CollapsePropertiesTest.testModuleExportsObjectSubPropertyCommonJs()
func (c *CollapsePropertiesTest) TestModuleExportsObjectSubPropertyCommonJs() {
	c.SetupModuleExportsOnly()
	inputs := []*SourceFile{}
	inputs = append(inputs,
		SourceFile_FromCode(
			"mod1.js",
			`var foo ={};
module.exports = foo;
module.exports.bar = 'bar';
`))
	inputs = append(inputs, SourceFile_FromCode("entry.js", "var mod = require('./mod1.js'); alert(mod.bar);"))
	expected := []*SourceFile{}
	expected = append(expected,
		SourceFile_FromCode(
			"mod1.js",
			`/** @const */ var module$mod1={};
 /** @const */ var module$mod1$default = {};
module$mod1$default.bar = 'bar';
`))
	expected = append(expected,
		SourceFile_FromCode(
			"entry.js", "var mod = module$mod1$default; alert(module$mod1$default.bar);"))
	c.Test(c.Srcs(inputs), c.Expected(expected))
}

// Java: CollapsePropertiesTest.testOrExpression()
func (c *CollapsePropertiesTest) TestOrExpression() {
	// left branch, read declared prop
	c.TestSame("var a = {b: 1} || x; var t = a.b;")
	// left branch, read undeclared prop
	c.TestSame("var a = {b: 1} || x; var t = a.c;")
	// left branch, write to declared prop
	c.TestSame("var a = {b: 1} || x; a.b = 2;")
	// left branch, write to undeclared prop
	// TODO(tjgq): Consider also collapsing a.c in this case.
	c.TestSame("var a = {b: 1} || x; a.c = 2;")
	// right branch, read declared prop
	// aggressive alias inlining still works
	c.Test(
		"var a = x || {b: 1}; var t =  a.b; use(  t);", //
		"var a = x || {b: 1}; var t = null; use(a.b);")
	// right branch, read undeclared prop
	c.TestSame("var a = x || {b: 1}; var t = a.c;")
	// right branch, write to declared prop
	c.TestSame("var a = x || {b: 1}; a.b = 2;")
	// right branch, write to undeclared prop
	c.Test(c.Srcs("var a = x || {b: 1}; a.c = 2;"), c.Expected("var a = x || {b: 1}; var a$c = 2;"))
}

// Java: CollapsePropertiesTest.testTernaryExpression()
func (c *CollapsePropertiesTest) TestTernaryExpression() {
	// left branch, read declared prop
	// aggressive alias inlining still works
	c.Test(
		"var a = p ? {b: 1} : x; var t =  a.b; use(  t);", //
		"var a = p ? {b: 1} : x; var t = null; use(a.b);")
	// left branch, read undeclared prop
	c.TestSame("var a = p ? {b: 1} : x; var t = a.c;")
	// left branch, write to declared prop
	c.TestSame("var a = p ? {b: 1} : x; a.b = 2;")
	// left branch, write to undeclared prop
	c.Test(
		c.Srcs("var a = p ? {b: 1} : x; a.c = 2;"), c.Expected("var a = p ? {b: 1} : x; var a$c = 2;"))
	// right branch, read declared prop
	// aggressive alias inlining still works
	c.Test(
		"var a = p ? x : {b: 1}; var t =  a.b; use(  t);", //
		"var a = p ? x : {b: 1}; var t = null; use(a.b);")
	// right branch, read undeclared prop
	c.TestSame("var a = p ? x : {b: 1}; var t = a.c;")
	// right branch, write to declared prop
	c.TestSame("var a = p ? x : {b: 1}; a.b = 2;")
	// right branch, write to undeclared prop
	c.Test(
		c.Srcs("var a = p ? x : {b: 1}; a.c = 2;"), c.Expected("var a = p ? x : {b: 1}; var a$c = 2;"))
	// right branch, use nested prop, alias is inlined but prop is not collapsed.
	c.Test(
		c.Srcs("    var a = p ? x : {b: { insideB: 1 }}; var t = a.b.insideB; use(          t);"),
		c.Expected("var a = p ? x : {b: { insideB: 1 }}; var t = null       ; use(a.b.insideB);"))
}

// Java: CollapsePropertiesTest.testModuleExportsOnly()
func (c *CollapsePropertiesTest) TestModuleExportsOnly() {
	c.SetupModuleExportsOnly()
	inputs := []*SourceFile{}
	inputs = append(inputs,
		SourceFile_FromCode(
			"entry.js",
			`var mod = require('./mod1.js');
alert(mod);
var a = {};
a.b = {};
/** @constructor */
a.b.c = function(){};
`))
	inputs = append(inputs, SourceFile_FromCode("mod1.js", "module.exports = 123;"))
	expected := []*SourceFile{}
	expected = append(expected,
		SourceFile_FromCode(
			"entry.js",
			`var mod = module$mod1$default;
alert(module$mod1$default);
var a = {};
a.b = {};
/** @constructor */
a.b.c = function(){};
`))
	expected = append(expected,
		SourceFile_FromCode(
			"mod1.js",
			`/** @const */ var module$mod1 = {};
/** @const */ var module$mod1$default = 123;
`))
	c.Test(c.Srcs(inputs), c.Expected(expected))
}

// Port of com/google/javascript/jscomp/disambiguate/DisambiguatePropertiesTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/disambiguate/DisambiguatePropertiesTest.java)

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DisambiguatePropertiesTest tests DisambiguateProperties.
//
// The Java @Test methods not named test* get a "Test" prefix in Go so that
// RunJUnit4 runs them.
//
// Java: public final class DisambiguatePropertiesTest extends CompilerTestCase
type DisambiguatePropertiesTest struct {
	CompilerTestCaseBase

	propertiesThatMustDisambiguate *jutil.LinkedHashSet[string]
}

// Java: DisambiguatePropertiesTest.PROP_DEFINER_DEFINITION
const disambiguatePropertiesTest_PROP_DEFINER_DEFINITION = `/** @const */ var goog = {};
/** @const */ goog.reflect = {};

/** @return {string} */
goog.reflect.objectProperty = function(prop, obj) { return ''; };
`

// Java: DisambiguatePropertiesTest()
func NewDisambiguatePropertiesTest() *DisambiguatePropertiesTest {
	c := &DisambiguatePropertiesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, "")
	c.propertiesThatMustDisambiguate = jutil.NewLinkedHashSet[string]()
	return c
}

func TestDisambiguatePropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewDisambiguatePropertiesTest() })
}

// Java: DisambiguatePropertiesTest.setUp()
func (c *DisambiguatePropertiesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()
	c.EnableDebugLogging(true)
}

// Java: DisambiguatePropertiesTest.getProcessor(Compiler)
func (c *DisambiguatePropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewDisambiguateProperties(compiler, c.propertiesThatMustDisambiguate)
}

// Java: DisambiguatePropertiesTest.getOptions()
func (c *DisambiguatePropertiesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.AddWarningsGuard(newDisambiguatePropertiesTest_SilenceNoiseGuard())
	return options
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_byExtends_fromInterface()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_byExtends_fromInterface() {
	c.Test(
		c.Srcs(
			`/** @interface */
class IFoo {
  x() { }
}

/**
 * @interface
 * @extends {IFoo}
 */
class IFoo2 {
  x() { }
}

class Other {
  x() { }
}
`),
		c.Expected(
			`/** @interface */
class IFoo {
  JSC$3_x() { }
}

/**
 * @interface
 * @extends {IFoo}
 */
class IFoo2 {
  JSC$3_x() { }
}

class Other {
  JSC$5_x() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_byExtends_fromClass()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_byExtends_fromClass() {
	c.Test(
		c.Srcs(
			`class Foo {
  y() { }
}

/** @extends {Foo} */
class Foo2 {
  y() { }
}

class Other {
  y() { }
}
`),
		c.Expected(
			`class Foo {
  JSC$1_y() { }
}

/** @extends {Foo} */
class Foo2 {
  JSC$1_y() { }
}

class Other {
  JSC$5_y() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_byExtends_viaPrototype()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_byExtends_viaPrototype() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { };
/** @type {number} */
Foo.z = 0;

class Foo2 extends Foo {
  static z() { }
}

class Other {
  z() { }
}
`),
		c.Expected(
			`/** @constructor */
function Foo() { };
/** @type {number} */
Foo.JSC$1_z = 0;

class Foo2 extends Foo {
  static JSC$1_z() { }
}

class Other {
  JSC$3_z() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_byImplements_fromInterface()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_byImplements_fromInterface() {
	c.Test(
		c.Srcs(
			`/** @interface */
class IFoo {
  a() { }
}

/** @implements {IFoo} */
class Foo {
  a() { }
}

class Other {
  a() { }
}
`),
		c.Expected(
			`/** @interface */
class IFoo {
  JSC$3_a() { }
}

/** @implements {IFoo} */
class Foo {
  JSC$3_a() { }
}

class Other {
  JSC$5_a() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_byCommonDescendent()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_byCommonDescendent() {
	c.Test(
		c.Srcs(
			`/** @interface */
class IFoo0 {
  b() { }
  c() { }
}

/** @interface */
class IFoo1 {
  b() { }
  d() { }
}

/**
 * @implements {IFoo0}
 * @implements {IFoo1}
 */
class Foo {
  b() { }
}

class Other {
  b() { }
  c() { }
  d() { }
}
`),
		c.Expected(
			`/** @interface */
class IFoo0 {
  JSC$5_b() { }
  JSC$1_c() { }
}

/** @interface */
class IFoo1 {
  JSC$5_b() { }
  JSC$3_d() { }
}

/**
 * @implements {IFoo0}
 * @implements {IFoo1}
 */
class Foo {
  JSC$5_b() { }
}

class Other {
  JSC$7_b() { }
  JSC$7_c() { }
  JSC$7_d() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_betweenDistantAncestors()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_betweenDistantAncestors() {
	c.Test(
		c.Srcs(
			`class Foo0 {
  f() { }
}

class Foo1 extends Foo0 { }
class Foo2 extends Foo1 { }
class Foo3 extends Foo2 { }
class Foo4 extends Foo3 { }

class Foo5 extends Foo4 {
  f() { }
}

class Other {
  f() { }
}
`),
		c.Expected(
			`class Foo0 {
  JSC$1_f() { }
}

class Foo1 extends Foo0 { }
class Foo2 extends Foo1 { }
class Foo3 extends Foo2 { }
class Foo4 extends Foo3 { }

class Foo5 extends Foo4 {
  JSC$1_f() { }
}

class Other {
  JSC$9_f() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_byUnions_iffUnionUsesProperty()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_byUnions_iffUnionUsesProperty() {
	c.Test(
		c.Srcs(
			`class Foo0 {
  a() { }
  b() { }
}

class Foo1 {
  a() { }
  b() { }
}

function mix(/** (!Foo0|!Foo1) */ x) {
  x.a();
}

class Other {
  a() { }
  b() { }
}
`),
		c.Expected(
			`class Foo0 {
  JSC$5_a() { }
  JSC$1_b() { }
}

class Foo1 {
  JSC$5_a() { }
  JSC$3_b() { }
}

function mix(/** (!Foo0|!Foo1) */ x) {
  x.JSC$5_a();
}

class Other {
  JSC$6_a() { }
  JSC$6_b() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_betweenUnionAncestors_andTypesMismatchedWithTheUnion()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_betweenUnionAncestors_andTypesMismatchedWithTheUnion() {
	c.TestSame(
		c.Srcs(
			`class Foo0 {
  a() { }
}
class Foo1 extends Foo0 { }
class Foo2 extends Foo0 { }

class Bar {
  a() { }
}

/** @type {(!Foo1|!Foo2)} @suppress {checkTypes} */
const x = new Bar();

class Other {
  a() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_acrossTemplateSpecializations()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_acrossTemplateSpecializations() {
	c.Test(
		c.Srcs(
			`/** @template T */
class Foo {
  h() { }
}

const /** !Foo<string> */ a = new Foo();
a.h();
const /** !Foo<?> */ b = new Foo();
b.h();
const /** !Foo */ c = new Foo();
c.h();

class Other {
  h() { }
}
`),
		c.Expected(
			`/** @template T */
class Foo {
  JSC$1_h() { }
}

const /** !Foo<string> */ a = new Foo();
a.JSC$1_h();
const /** !Foo<?> */ b = new Foo();
b.JSC$1_h();
const /** !Foo */ c = new Foo();
c.JSC$1_h();

class Other {
  JSC$4_h() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_betweenEnumsAndTheirValues()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_betweenEnumsAndTheirValues() {
	c.Test(
		c.Srcs(
			`class Foo {
  p() { }
}

/** @enum {!Foo} */
const Bar = {
  ONE: new Foo(),
};

Bar.ONE.p

class Other {
  p() { }
}
`),
		c.Expected(
			`class Foo {
  JSC$1_p() { }
}

/** @enum {!Foo} */
const Bar = {
  ONE: new Foo(),
};

Bar.ONE.JSC$1_p

class Other {
  JSC$5_p() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreConflated_acrossNullableAndVoidableReferences()
func (c *DisambiguatePropertiesTest) TestPropertiesAreConflated_acrossNullableAndVoidableReferences() {
	c.Test(
		c.Srcs(
			`class Foo {
  g() { }
}

function use(/** (null|undefined|!Foo) */ x) {
  x.g()
}

class Other {
  g() { }
}
`),
		c.Expected(
			`class Foo {
  JSC$1_g() { }
}

function use(/** (null|undefined|!Foo) */ x) {
  x.JSC$1_g()
}

class Other {
  JSC$4_g() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreDisambiguated_betweenContructorAndInstanceTypes()
func (c *DisambiguatePropertiesTest) TestPropertiesAreDisambiguated_betweenContructorAndInstanceTypes() {
	c.Test(
		c.Srcs(
			`class Foo {
  static g() { }
  g() { }
}
`),
		c.Expected(
			`class Foo {
  static JSC$1_g() { }
  JSC$2_g() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreDisambiguated_betweenAncestorTypes_ifCommonChildHasNoPropUse()
func (c *DisambiguatePropertiesTest) TestPropertiesAreDisambiguated_betweenAncestorTypes_ifCommonChildHasNoPropUse() {
	c.Test(
		c.Srcs(
			`/** @interface */
class IFoo {
  t() { }
}
class Foo {
  t() { }
}
/** @implements {IFoo} */
class SubFoo extends Foo { }

class Other {
  t() { }
}
`),
		c.Expected(
			`/** @interface */
class IFoo {
  JSC$3_t() { }
}
class Foo {
  JSC$3_t() { }
}
/** @implements {IFoo} */
class SubFoo extends Foo { }

class Other {
  JSC$6_t() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreInvalidated_byToplikeAndBottomlikeTypes()
func (c *DisambiguatePropertiesTest) TestPropertiesAreInvalidated_byToplikeAndBottomlikeTypes() {
	annotations :=
		[]string{"?", "*", "!null", "(?|undefined)", "null", "undefined", "(null|undefined)"}

	for _, annotation := range annotations {
		c.Test(
			c.Srcs(strings.ReplaceAll(
				`class Foo0 {
  k() { }
  t() { }
}
class Foo1 {
  k() { }
  t() { }
}

function mix(/** ANNOTATION */ x) {
  x.k();
}
`, "ANNOTATION", annotation)),
			c.Expected(strings.ReplaceAll(
				`class Foo0 {
  k() { }
  JSC$1_t() { }
}
class Foo1 {
  k() { }
  JSC$3_t() { }
}

function mix(/** ANNOTATION */ x) {
  x.k();
}
`, "ANNOTATION", annotation)))
	}
}

// Java: DisambiguatePropertiesTest.propertiesAreInvalidated_onObjectLiteralTypes()
func (c *DisambiguatePropertiesTest) TestPropertiesAreInvalidated_onObjectLiteralTypes() {
	c.TestSame(
		c.Srcs(
			`const z = {ab: 0};

class Other {
  ab() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreDisambiguated_acrossStructuralTypeMatches_iffMatchUsed()
func (c *DisambiguatePropertiesTest) TestPropertiesAreDisambiguated_acrossStructuralTypeMatches_iffMatchUsed() {
	c.TestSame(
		c.Srcs(
			`/** @record */
class Foo0 {
  a() { }
}
/** @record */
class Foo1 extends Foo0 {
  a() { }
  b() { }
}
class Foo2 {
  a() { }
  b() { }
}

const /** !Foo0 */ x = new Foo2();

class Other {
 a() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreInvalidated_byAnonymousFunctionTypes()
func (c *DisambiguatePropertiesTest) TestPropertiesAreInvalidated_byAnonymousFunctionTypes() {
	c.Test(
		c.Srcs(
			`class Foo {
  static g() { }
  static j() { }
}

function use(/** function() */ x) {
  x.g;
}

class Other {
  j() { }
}
`),
		c.Expected(
			`class Foo {
  static g() { }
  static JSC$1_j() { }
}

function use(/** function() */ x) {
  x.g;
}

class Other {
  JSC$3_j() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreDisambiguated_acrossTypeMismatches()
func (c *DisambiguatePropertiesTest) TestPropertiesAreDisambiguated_acrossTypeMismatches() {
	c.TestSame(
		c.Srcs(
			`class Foo0 {
  a() { }
}
class Foo1 extends Foo0 {
  a() { }
  b() { }
}
class Foo2 {
  a() { }
  b() { }
}

const /** !Foo0 */ x = new Foo2();

class Other {
 a() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreNotRenamed_fromExternsClusters()
func (c *DisambiguatePropertiesTest) TestPropertiesAreNotRenamed_fromExternsClusters() {
	c.Test(
		c.Externs(
			`class ExternFoo {
  v() { }
}
`),
		c.Srcs(
			`class Foo extends ExternFoo {
  v() { }
  p() { }
}

class Other {
  v() { }
  p() { }
}
`),
		c.Expected(
			`class Foo extends ExternFoo {
  v() { }
  JSC$3_p() { }
}

class Other {
  JSC$5_v() { }
  JSC$5_p() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreNotRenamed_ifTheyHaveASingleCluster()
func (c *DisambiguatePropertiesTest) TestPropertiesAreNotRenamed_ifTheyHaveASingleCluster() {
	c.Test(
		c.Srcs(
			`class Foo0 {
  w() { }
  x() { }
}
class Foo1 extends Foo0 {
  w() { }
  y() { }
}
class Foo2 {
  w() { }
  z() { }
}

function mix(/** (!Foo1|!Foo2) */ x) {
  x.w();
}

class Other {
  x() { }
  y() { }
  z() { }
}
`),
		c.Expected(
			`class Foo0 {
  w() { }
  JSC$1_x() { }
}
class Foo1 extends Foo0 {
  w() { }
  JSC$3_y() { }
}
class Foo2 {
  w() { }
  JSC$5_z() { }
}

function mix(/** (!Foo1|!Foo2) */ x) {
  x.w();
}

class Other {
  JSC$8_x() { }
  JSC$8_y() { }
  JSC$8_z() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreInvalidated_givenMissingPropertyError()
func (c *DisambiguatePropertiesTest) TestPropertiesAreInvalidated_givenMissingPropertyError() {
	c.Test(
		c.Srcs(
			`class Foo {
  x() { }
  y() { }
  z() { }
}
class Bar { }

function mix(/** (!Foo|!Bar) */ fooBar, /** !Bar */ bar) {
// Invalidate bar.y() but not fooBar.x(). This mirrors the Closure typechecker,
// which thinks (!Foo|!Bar) may have an 'x' property and so allows the access.
  fooBar.w();
  fooBar.x();
  bar.y();
}

class Other {
  w() { }
  x() { }
  y() { }
  z() { }
}
`),
		c.Expected(
			`class Foo {
  JSC$4_x() { }
  y() { }
  JSC$1_z() { }
}
class Bar { }

function mix(/** (!Foo|!Bar) */ fooBar, /** !Bar */ bar) {
// x is disambiguated while w and y are invalidated.
  fooBar.w();
  fooBar.JSC$4_x();
  bar.y();
}

class Other {
  w() { }
  JSC$6_x() { }
  y() { }
  JSC$6_z() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesAreInvalidated_ifUsedOnType_butNotDeclaredOnAncestor()
func (c *DisambiguatePropertiesTest) TestPropertiesAreInvalidated_ifUsedOnType_butNotDeclaredOnAncestor() {
	c.Test(
		c.Srcs(
			`/** @interface */
class IFoo {
  a() { }
}

/** @implements {IFoo} */
class Foo {
  b() { }
}

new Foo().a;
new Foo().b;
new Foo().c;

class Other {
  a() { }
  b() { }
  c() { }
}
`),
		c.Expected(
			`/** @interface */
class IFoo {
  JSC$1_a() { }
}

/** @implements {IFoo} */
class Foo {
  JSC$3_b() { }
}

new Foo().JSC$1_a;
new Foo().JSC$3_b;
new Foo().c;

class Other {
  JSC$6_a() { }
  JSC$6_b() { }
  c() { }
}
`))
}

// Java: DisambiguatePropertiesTest.propertiesReferenced_throughReflectorFunctions_areRenamed()
func (c *DisambiguatePropertiesTest) TestPropertiesReferenced_throughReflectorFunctions_areRenamed() {
	c.Test(
		c.Externs(disambiguatePropertiesTest_PROP_DEFINER_DEFINITION),
		c.Srcs(
			`class Foo {
  m() { }
}

class Other {
  m() { }
}

goog.reflect.objectProperty('m', Foo.prototype);
`),
		c.Expected(
			`class Foo {
  JSC$3_m() { }
}

class Other {
  JSC$5_m() { }
}

goog.reflect.objectProperty('JSC$3_m', Foo.prototype);
`))
}

// Java: DisambiguatePropertiesTest.errorReported_forInvalidation_ofSpecifiedPropNames_oncePerName()
func (c *DisambiguatePropertiesTest) TestErrorReported_forInvalidation_ofSpecifiedPropNames_oncePerName() {
	c.AllowSourcelessWarnings()
	c.propertiesThatMustDisambiguate = jutil.NewLinkedHashSet[string](
		"invalid0", //
		"invalid1")

	c.Test(
		c.Srcs(
			`class Foo {
  invalid0() { }
  invalid1() { }
}

function use(/** ? */ x) {
  x.invalid0
  x.invalid1
  x.invalid0
  x.invalid1
}
`),
		c.Error(DisambiguateProperties_PROPERTY_INVALIDATION),
		c.Error(DisambiguateProperties_PROPERTY_INVALIDATION))
}

// Java: DisambiguatePropertiesTest.invalidatingSubtype_doesNotInvalidatePropertyOnlyReferencedOnSupertype()
func (c *DisambiguatePropertiesTest) TestInvalidatingSubtype_doesNotInvalidatePropertyOnlyReferencedOnSupertype() {
	// TODO(b/135045845): track mismatches through subtypes/supertypes
	c.Test(
		c.Srcs(
			`class FooParent {
  parent() {}
}

class FooChild extends FooParent {
  child() {}
}

class BarParent {
  parent() {}
}

class BarChild extends BarParent {
  child() {}
}

/**
 * @suppress {checkTypes} intentional type error
 * @type {!FooChild}
 */
const fooChild = '';
`),
		c.Expected(
			`class FooParent {
  JSC$1_parent() {}
}

class FooChild extends FooParent {
  child() {}
}

class BarParent {
  JSC$5_parent() {}
}

class BarChild extends BarParent {
  child() {}
}

/**
 * @suppress {checkTypes} intentional type error
 * @type {!FooChild}
 */
const fooChild = '';
`))
}

// Java: DisambiguatePropertiesTest.classField()
func (c *DisambiguatePropertiesTest) TestClassField() {
	c.Test(
		c.Srcs(
			`class Foo {
  a;
  static a = 2;
}

class Bar {
  a;
  static a = 3;
}
`),
		c.Expected(
			`class Foo {
  JSC$1_a;
  static JSC$2_a = 2;
}

class Bar {
  JSC$3_a;
  static JSC$4_a = 3;
}
`))
}

// Java: DisambiguatePropertiesTest.classComputedField()
func (c *DisambiguatePropertiesTest) TestClassComputedField() {
	c.TestSame(
		c.Srcs(
			`class Foo {
  ['a'];
  'b' = 2;
  1 = 'hello';
  static ['a'] = 2;
  static 'b' = 2;
  static 1 = 'hello';
}

class Bar {
  ['a'];
  'b' = 2;
  1 = 'hello';
  static ['a'] = 2;
  static 'b' = 2;
  static 1 = 'hello';
}
`))
}

// Java: DisambiguatePropertiesTest.classMixedFields()
func (c *DisambiguatePropertiesTest) TestClassMixedFields() {
	c.Test(
		c.Srcs(
			`class Foo {
  a;
  static a = 2;
  ['a'];
  static ['a'] = 2;
}

class Bar {
  a;
  static a = 2;
  ['a'];
  static ['a'] = 2;
}
`),
		c.Expected(
			`class Foo {
  JSC$1_a;
  static JSC$2_a = 2;
  ['a'];
  static ['a'] = 2;
}

class Bar {
  JSC$3_a;
  static JSC$4_a = 2;
  ['a'];
  static ['a'] = 2;
}
`))
}

// Java: DisambiguatePropertiesTest.classStaticBlock()
func (c *DisambiguatePropertiesTest) TestClassStaticBlock() {
	c.Test(
		c.Srcs(
			`class Foo {
  static {
    this.x = 1;
  }
}
Foo.x;
class Bar {
  static {
    this.x = 2;
  }
}
`),
		c.Expected(
			`class Foo {
  static {
    this.JSC$1_x = 1;
  }
}
Foo.JSC$1_x;
class Bar {
  static {
    this.JSC$2_x = 2;
  }
}
`))
	c.Test(
		c.Srcs(
			`class Foo {
  static x = 0;
  x = 0;
  static {
    this.x = 1;
  }
}
Foo.x;
let f = new Foo();
f.x;
class Bar {
  static x = 0;
  x = 0;
  static {
    this.x = 2;
  }
}
`),
		c.Expected(
			`class Foo {
  static JSC$1_x = 0;
  JSC$2_x = 0;
  static {
    this.JSC$1_x = 1;
  }
}
Foo.JSC$1_x;
let f = new Foo();
f.JSC$2_x;
class Bar {
  static JSC$3_x = 0;
  JSC$4_x = 0;
  static {
    this.JSC$3_x = 2;
  }
}
`))
}

// Java: DisambiguatePropertiesTest.classStaticBlock_method()
func (c *DisambiguatePropertiesTest) TestClassStaticBlock_method() {
	c.Test(
		c.Srcs(
			`class Foo {
  m() { }
}

class Baz {
  static {
    let f = new Foo();
    f.m();
  }
}
class Bar {
  m() { }
}
`),
		c.Expected(
			`class Foo {
  JSC$1_m() { }
}

class Baz {
  static {
    let f = new Foo();
    f.JSC$1_m();
  }
}
class Bar {
  JSC$5_m() { }
}
`))

	c.Test(
		c.Srcs(
			`class Foo {
  static m() { }
}

class Baz {
  static m() { }
  static {
    Foo.m();
  }
}
class Bar {
  static m() { }
}
`),
		c.Expected(
			`class Foo {
  static JSC$1_m() { }
}

class Baz {
  static JSC$2_m() { }
  static {
    Foo.JSC$1_m();
  }
}
class Bar {
  static JSC$3_m() { }
}
`))
}

// Java: DisambiguatePropertiesTest.testClassStaticBlock_function()
func (c *DisambiguatePropertiesTest) TestClassStaticBlock_function() {
	c.Test(
		c.Srcs(
			`/** @constructor */
function Foo() { };

class Foo2 extends Foo {
  static z() { }
  static {
    Foo.z = 0;
  }
}

class Other {
  z() { }
}
`),
		c.Expected(
			`/** @constructor */
function Foo() { };

class Foo2 extends Foo {
  static JSC$1_z() { }
  static {
    Foo.JSC$1_z = 0;
  }
}

class Other {
  JSC$3_z() { }
}
`))
}

// Java: DisambiguatePropertiesTest.generateDiagnosticFiles()
func (c *DisambiguatePropertiesTest) GenerateDiagnosticFiles() {
	c.Test(
		c.Srcs(
			`class Foo0 {
  w() { }
  x() { }
}
class Foo1 extends Foo0 {
  w() { }
  y() { }
}
class Foo2 {
  w() { }
  z() { }
}
class Bar {}
function use(/** (!Foo0|!Foo2) */ x, /** !Bar */ y) {
  x.x();
  y.y();
}
`))
}

// Java: private static final class DisambiguatePropertiesTest.SilenceNoiseGuard extends WarningsGuard
type disambiguatePropertiesTest_SilenceNoiseGuard struct {
	WarningsGuardBase
}

// Java: DisambiguatePropertiesTest.SilenceNoiseGuard.RELEVANT_DIAGNOSTICS
var disambiguatePropertiesTest_SilenceNoiseGuard_RELEVANT_DIAGNOSTICS = jutil.NewLinkedHashSet[*DiagnosticType](DisambiguateProperties_PROPERTY_INVALIDATION)

func newDisambiguatePropertiesTest_SilenceNoiseGuard() *disambiguatePropertiesTest_SilenceNoiseGuard {
	return &disambiguatePropertiesTest_SilenceNoiseGuard{}
}

// Java: DisambiguatePropertiesTest.SilenceNoiseGuard.getPriority()
func (g *disambiguatePropertiesTest_SilenceNoiseGuard) GetPriority() int {
	return WarningsGuard_Priority_MAX.GetValue()
}

// Java: DisambiguatePropertiesTest.SilenceNoiseGuard.level(JSError)
func (g *disambiguatePropertiesTest_SilenceNoiseGuard) Level(error *JSError) *CheckLevel {
	if strings.Contains(error.Description(), "Parse") {
		return nil
	} else if disambiguatePropertiesTest_SilenceNoiseGuard_RELEVANT_DIAGNOSTICS.Contains(error.Type()) {
		return nil
	} else {
		return CheckLevel_OFF.Ptr()
	}
}

// Java: DisambiguatePropertiesTest.loadFile(Path)
func disambiguatePropertiesTest_loadFile(path string) string {
	return serializeTypedAstPassTest_loadFile(path)
}

// Java: DisambiguatePropertiesTest.debugLogFiles()
func (c *DisambiguatePropertiesTest) debugLogFiles() []string {
	dir := filepath.Join(
		c.GetLastCompiler().GetOptions().GetDebugLogDirectory().String(),
		"DisambiguateProperties")
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(jutil.NewAssertionErrorWithMessageAndCause(err.Error(), jutil.NewIOException(err.Error())))
	}
	// Go-only: language-adapter — Files.list's order is unspecified; sorted.
	var files []string
	for _, e := range entries {
		files = append(files, filepath.Join(dir, e.Name()))
	}
	sort.Strings(files)
	return files
}

// Java: DisambiguatePropertiesTest.assertValidJson(String)
func disambiguatePropertiesTest_assertValidJson(src string) {
	serializeTypedAstPassTest_assertValidJson(src)
}

// Port of com/google/javascript/jscomp/ClosureCodingConventionTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ClosureCodingConventionTest.java)

/*
 * Copyright 2007 The Closure Compiler Authors.
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

// ClosureCodingConventionTest is the test class for GoogleCodingConvention.
type ClosureCodingConventionTest struct {
	t    *testing.T
	conv *ClosureCodingConvention
}

func (c *ClosureCodingConventionTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *ClosureCodingConventionTest) JUnitBefore()           {}
func (c *ClosureCodingConventionTest) JUnitAfter()            {}

// Java: ClosureCodingConventionTest()
func NewClosureCodingConventionTest() *ClosureCodingConventionTest {
	c := &ClosureCodingConventionTest{}
	c.conv = NewClosureCodingConvention()
	return c
}

func TestClosureCodingConventionTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewClosureCodingConventionTest() })
}

// Java: ClosureCodingConventionTest.testVarAndOptionalParams()
func (c *ClosureCodingConventionTest) TestVarAndOptionalParams() {
	args :=
		rhino.NewNodeWith2Children(rhino.Token_PARAM_LIST, rhino.Node_NewStringWithToken(rhino.Token_NAME, "a"), rhino.Node_NewStringWithToken(rhino.Token_NAME, "b"))
	optArgs :=
		rhino.NewNodeWith2Children(rhino.Token_PARAM_LIST, rhino.Node_NewStringWithToken(rhino.Token_NAME, "opt_a"), rhino.Node_NewStringWithToken(rhino.Token_NAME, "opt_b"))
	AssertThat(c.conv.IsVarArgsParameter(args.GetFirstChild())).IsFalse()
	AssertThat(c.conv.IsVarArgsParameter(args.GetLastChild())).IsFalse()
	AssertThat(c.conv.IsVarArgsParameter(optArgs.GetFirstChild())).IsFalse()
	AssertThat(c.conv.IsVarArgsParameter(optArgs.GetLastChild())).IsFalse()
	AssertThat(c.conv.IsOptionalParameter(args.GetFirstChild())).IsFalse()
	AssertThat(c.conv.IsOptionalParameter(args.GetLastChild())).IsFalse()
	AssertThat(c.conv.IsOptionalParameter(optArgs.GetFirstChild())).IsFalse()
	AssertThat(c.conv.IsOptionalParameter(optArgs.GetLastChild())).IsFalse()
}

// Java: ClosureCodingConventionTest.testInlineName()
func (c *ClosureCodingConventionTest) TestInlineName() {
	AssertThat(c.conv.IsConstant("a")).IsFalse()
	AssertThat(c.conv.IsConstant("XYZ123_")).IsFalse()
	AssertThat(c.conv.IsConstant("ABC")).IsFalse()
	AssertThat(c.conv.IsConstant("ABCdef")).IsFalse()
	AssertThat(c.conv.IsConstant("aBC")).IsFalse()
	AssertThat(c.conv.IsConstant("A")).IsFalse()
	AssertThat(c.conv.IsConstant("_XYZ123")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$XYZ123_")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$ABC_DEF")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$A")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$a")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$ABCdef")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$aBC")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$")).IsFalse()
	AssertThat(c.conv.IsConstant("$")).IsFalse()
}

// Java: ClosureCodingConventionTest.testExportedName()
func (c *ClosureCodingConventionTest) TestExportedName() {
	AssertThat(c.conv.IsExported("_a")).IsFalse()
	AssertThat(c.conv.IsExported("_a_")).IsFalse()
	AssertThat(c.conv.IsExported("a")).IsFalse()
	AssertThat(c.conv.IsExportedWithLocal("$super", false)).IsFalse()
	AssertThat(c.conv.IsExportedWithLocal("$super", true)).IsTrue()
	AssertThat(c.conv.IsExported("$super")).IsTrue()
}

// Java: ClosureCodingConventionTest.testEnumKey()
func (c *ClosureCodingConventionTest) TestEnumKey() {
	AssertThat(c.conv.IsValidEnumKey("A")).IsTrue()
	AssertThat(c.conv.IsValidEnumKey("123")).IsTrue()
	AssertThat(c.conv.IsValidEnumKey("FOO_BAR")).IsTrue()
	AssertThat(c.conv.IsValidEnumKey("a")).IsTrue()
	AssertThat(c.conv.IsValidEnumKey("someKeyInCamelCase")).IsTrue()
	AssertThat(c.conv.IsValidEnumKey("_FOO_BAR")).IsTrue()
}

// Java: ClosureCodingConventionTest.testInheritanceDetection1()
func (c *ClosureCodingConventionTest) TestInheritanceDetection1() {
	c.AssertNotClassDefining("goog.foo(A, B);")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection2()
func (c *ClosureCodingConventionTest) TestInheritanceDetection2() {
	c.AssertDefinesClasses("goog.inherits(A, B);", "A", "B")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection3()
func (c *ClosureCodingConventionTest) TestInheritanceDetection3() {
	c.AssertNotClassDefining("A.inherits(B);")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection4()
func (c *ClosureCodingConventionTest) TestInheritanceDetection4() {
	c.AssertDefinesClasses("goog.inherits(goog.A, goog.B);", "goog.A", "goog.B")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection5()
func (c *ClosureCodingConventionTest) TestInheritanceDetection5() {
	c.AssertNotClassDefining("goog.A.inherits(goog.B);")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection6()
func (c *ClosureCodingConventionTest) TestInheritanceDetection6() {
	c.AssertNotClassDefining("A.inherits(this.B);")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection7()
func (c *ClosureCodingConventionTest) TestInheritanceDetection7() {
	c.AssertNotClassDefining("this.A.inherits(B);")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection8()
func (c *ClosureCodingConventionTest) TestInheritanceDetection8() {
	c.AssertDefinesClasses("goog.inherits(A, B, C);", "A", "B")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection9()
func (c *ClosureCodingConventionTest) TestInheritanceDetection9() {
	c.AssertNotClassDefining("A.mixin(B.prototype);")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection11()
func (c *ClosureCodingConventionTest) TestInheritanceDetection11() {
	c.AssertNotClassDefining("A.mixin(B)")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection15()
func (c *ClosureCodingConventionTest) TestInheritanceDetection15() {
	c.AssertDefinesClasses("$jscomp.inherits(A, B)", "A", "B")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection16()
func (c *ClosureCodingConventionTest) TestInheritanceDetection16() {
	c.AssertDefinesClasses("$jscomp$inherits(A, B)", "A", "B")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection17()
func (c *ClosureCodingConventionTest) TestInheritanceDetection17() {
	c.AssertDefinesClasses(
		"ValueType.mixin(A, B, 5, goog.reflect.objectProperty('foo', A))", "A", "B")
}

// Java: ClosureCodingConventionTest.testInheritanceDetection18()
func (c *ClosureCodingConventionTest) TestInheritanceDetection18() {
	c.AssertNotClassDefining("ValueType.mixin(A, 5, goog.reflect.objectProperty('foo', A))")
}

// Java: ClosureCodingConventionTest.testInheritanceDetectionPostCollapseProperties()
func (c *ClosureCodingConventionTest) TestInheritanceDetectionPostCollapseProperties() {
	c.AssertDefinesClasses("goog$inherits(A, B);", "A", "B")
	c.AssertNotClassDefining("goog$inherits(A);")
}

// Java: ClosureCodingConventionTest.testObjectLiteralCast()
func (c *ClosureCodingConventionTest) TestObjectLiteralCast() {
	c.AssertNotObjectLiteralCast("goog.reflect.object();")
	c.AssertNotObjectLiteralCast("goog.reflect.object(A);")
	c.AssertNotObjectLiteralCast("goog.reflect.object(1, {});")
	c.AssertObjectLiteralCast("goog.reflect.object(A, {});")
	c.AssertNotObjectLiteralCast("$jscomp.reflectObject();")
	c.AssertNotObjectLiteralCast("$jscomp.reflectObject(A);")
	c.AssertNotObjectLiteralCast("$jscomp.reflectObject(1, {});")
	c.AssertObjectLiteralCast("$jscomp.reflectObject(A, {});")
}

// Java: ClosureCodingConventionTest.testFunctionBind()
func (c *ClosureCodingConventionTest) TestFunctionBind() {
	c.AssertNotFunctionBind("goog.bind()") // invalid bind
	c.AssertFunctionBind("goog.bind(f)")
	c.AssertFunctionBind("goog.bind(f, obj)")
	c.AssertFunctionBind("goog.bind(f, obj, p1)")
	c.AssertNotFunctionBind("goog$bind()") // invalid bind
	c.AssertFunctionBind("goog$bind(f)")
	c.AssertFunctionBind("goog$bind(f, obj)")
	c.AssertFunctionBind("goog$bind(f, obj, p1)")
	c.AssertNotFunctionBind("goog.partial()") // invalid bind
	c.AssertFunctionBind("goog.partial(f)")
	c.AssertFunctionBind("goog.partial(f, obj)")
	c.AssertFunctionBind("goog.partial(f, obj, p1)")
	c.AssertNotFunctionBind("goog$partial()") // invalid bind
	c.AssertFunctionBind("goog$partial(f)")
	c.AssertFunctionBind("goog$partial(f, obj)")
	c.AssertFunctionBind("goog$partial(f, obj, p1)")
	c.AssertFunctionBind("(function(){}).bind()")
	c.AssertFunctionBind("(function(){}).bind(obj)")
	c.AssertFunctionBind("(function(){}).bind(obj, p1)")
	c.AssertNotFunctionBind("Function.prototype.bind.call()")
	c.AssertFunctionBind("Function.prototype.bind.call(obj)")
	c.AssertFunctionBind("Function.prototype.bind.call(obj, p1)")
}

// Java: ClosureCodingConventionTest.testRequire()
func (c *ClosureCodingConventionTest) TestRequire() {
	c.AssertRequire("goog.require('foo')")
	c.AssertNotRequire("goog.require(foo)")
	c.AssertNotRequire("goog.require()")
	c.AssertNotRequire("foo()")
}

// Java: ClosureCodingConventionTest.testApplySubclassRelationship()
func (c *ClosureCodingConventionTest) TestApplySubclassRelationship() {
	registry := rhino.NewJSTypeRegistry(nil)
	func() {
		closer := registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		nodeA := rhino.NewNode(rhino.Token_FUNCTION)
		a := "A"
		ctorA :=
			registry.CreateConstructorType(
				&a, nodeA, registry.CreateParameters(), nil, nil, false)
		nodeB := rhino.NewNode(rhino.Token_FUNCTION)
		b := "B"
		ctorB :=
			registry.CreateConstructorType(
				&b, nodeB, registry.CreateParameters(), nil, nil, false)
		c.conv.ApplySubclassRelationship(
			rhino.NewNominalTypeBuilder(ctorA, ctorA.GetInstanceType()),
			rhino.NewNominalTypeBuilder(ctorB, ctorB.GetInstanceType()),
			CodingConvention_SubclassType_INHERITS)
		AssertThat(ctorB.GetPrototype().HasOwnProperty("constructor")).IsTrue()
		AssertThat(ctorB.GetPrototype().GetPropertyNode("constructor")).IsEqualTo(nodeB)
		AssertThat(ctorB.HasOwnProperty("superClass_")).IsTrue()
		AssertThat(ctorB.GetPropertyNode("superClass_")).IsEqualTo(nodeB)
	}()
}

// Java: ClosureCodingConventionTest.testDescribeCachingCall()
func (c *ClosureCodingConventionTest) TestDescribeCachingCall() {
	c.AssertCachingCall("goog.reflect.cache(obj, 10, function() {})")
	c.AssertCachingCall("goog.reflect.cache(obj, 10, function() {}, function() {})")
	c.AssertCachingCall("goog$reflect$cache(obj, 10, function() {})")
	c.AssertCachingCall("goog$reflect$cache(obj, 10, function() {}, function() {})")
	c.AssertNotCachingCall("goog.reflect.cache()")
	c.AssertNotCachingCall("goog.reflect.cache(obj)")
	c.AssertNotCachingCall("goog.reflect.cache(obj, 10)")
	c.AssertNotCachingCall("foo.cache(obj, 10, function() {}, function() {})")
}

// Java: ClosureCodingConventionTest.assertFunctionBind(String)
func (c *ClosureCodingConventionTest) AssertFunctionBind(code string) {
	n := c.ParseTestCode(code)
	AssertThat(c.conv.DescribeFunctionBind(n.GetFirstChild(), false)).IsNotNull()
}

// Java: ClosureCodingConventionTest.assertNotFunctionBind(String)
func (c *ClosureCodingConventionTest) AssertNotFunctionBind(code string) {
	n := c.ParseTestCode(code)
	AssertThat(c.conv.DescribeFunctionBind(n.GetFirstChild(), false)).IsNull()
}

// Java: ClosureCodingConventionTest.assertRequire(String)
func (c *ClosureCodingConventionTest) AssertRequire(code string) {
	n := c.ParseTestCode(code)
	AssertThat(c.conv.ExtractClassNameIfRequire(n.GetFirstChild(), n)).IsNotNull()
}

// Java: ClosureCodingConventionTest.assertNotRequire(String)
func (c *ClosureCodingConventionTest) AssertNotRequire(code string) {
	n := c.ParseTestCode(code)
	AssertThat(c.conv.ExtractClassNameIfRequire(n.GetFirstChild(), n)).IsNull()
}

// Java: ClosureCodingConventionTest.assertNotObjectLiteralCast(String)
func (c *ClosureCodingConventionTest) AssertNotObjectLiteralCast(code string) {
	n := c.ParseTestCode(code)
	AssertThat(c.conv.GetObjectLiteralCast(n.GetFirstChild())).IsNull()
}

// Java: ClosureCodingConventionTest.assertObjectLiteralCast(String)
func (c *ClosureCodingConventionTest) AssertObjectLiteralCast(code string) {
	n := c.ParseTestCode(code)
	AssertThat(c.conv.GetObjectLiteralCast(n.GetFirstChild())).IsNotNull()
}

// Java: ClosureCodingConventionTest.assertNotClassDefining(String)
func (c *ClosureCodingConventionTest) AssertNotClassDefining(code string) {
	n := c.ParseTestCode(code)
	AssertThat(c.conv.GetClassesDefinedByCall(n.GetFirstChild())).IsNull()
}

// Java: ClosureCodingConventionTest.assertDefinesClasses(String, String, String)
func (c *ClosureCodingConventionTest) AssertDefinesClasses(code string, subclassName string, superclassName string) {
	n := c.ParseTestCode(code)
	classes := c.conv.GetClassesDefinedByCall(n.GetFirstChild())
	AssertThat(classes).IsNotNull()
	AssertThat(classes.subclassName).IsEqualTo(subclassName)
	AssertThat(classes.superclassName).IsEqualTo(superclassName)
}

// Java: ClosureCodingConventionTest.assertCachingCall(String)
func (c *ClosureCodingConventionTest) AssertCachingCall(code string) {
	node := c.ParseTestCode(code)
	AssertThat(c.conv.DescribeCachingCall(node.GetFirstChild())).IsNotNull()
}

// Java: ClosureCodingConventionTest.assertNotCachingCall(String)
func (c *ClosureCodingConventionTest) AssertNotCachingCall(code string) {
	node := c.ParseTestCode(code)
	AssertThat(c.conv.DescribeCachingCall(node.GetFirstChild())).IsNull()
}

// Java: ClosureCodingConventionTest.parseTestCode(String)
func (c *ClosureCodingConventionTest) ParseTestCode(code string) rhino.Node {
	compiler := NewCompiler()
	return compiler.ParseTestCode(code).GetFirstChild()
}

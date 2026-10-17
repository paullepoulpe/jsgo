// Port of com/google/javascript/jscomp/GoogleCodingConventionTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/GoogleCodingConventionTest.java)

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

// GoogleCodingConventionTest is the test class for GoogleCodingConvention.
type GoogleCodingConventionTest struct {
	t    *testing.T
	conv *GoogleCodingConvention
}

func (c *GoogleCodingConventionTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *GoogleCodingConventionTest) JUnitBefore()           {}
func (c *GoogleCodingConventionTest) JUnitAfter()            {}

// Java: GoogleCodingConventionTest()
func NewGoogleCodingConventionTest() *GoogleCodingConventionTest {
	c := &GoogleCodingConventionTest{}
	c.conv = NewGoogleCodingConvention()
	return c
}

func TestGoogleCodingConventionTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewGoogleCodingConventionTest() })
}

// Java: GoogleCodingConventionTest.testVarAndOptionalParams()
func (c *GoogleCodingConventionTest) TestVarAndOptionalParams() {
	args := rhino.IR_ParamList(rhino.IR_Name("a"), rhino.IR_Name("b"))
	optArgs := rhino.IR_ParamList(rhino.IR_Name("opt_a"), rhino.IR_Name("opt_b"))
	rest := rhino.IR_ParamList(rhino.IR_IterRest(rhino.IR_Name("more")))
	AssertThat(c.conv.IsVarArgsParameter(args.GetFirstChild())).IsFalse()
	AssertThat(c.conv.IsVarArgsParameter(args.GetLastChild())).IsFalse()
	AssertThat(c.conv.IsVarArgsParameter(optArgs.GetFirstChild())).IsFalse()
	AssertThat(c.conv.IsVarArgsParameter(optArgs.GetLastChild())).IsFalse()
	AssertThat(c.conv.IsOptionalParameter(args.GetFirstChild())).IsFalse()
	AssertThat(c.conv.IsOptionalParameter(args.GetLastChild())).IsFalse()
	AssertThat(c.conv.IsOptionalParameter(optArgs.GetFirstChild())).IsTrue()
	AssertThat(c.conv.IsOptionalParameter(optArgs.GetLastChild())).IsTrue()
	AssertThat(c.conv.IsVarArgsParameter(rest.GetLastChild())).IsTrue()
	AssertThat(c.conv.IsOptionalParameter(rest.GetFirstChild())).IsFalse()
}

// Java: GoogleCodingConventionTest.testIsConstant()
func (c *GoogleCodingConventionTest) TestIsConstant() {
	AssertThat(c.conv.IsConstant("a")).IsFalse()
	AssertThat(c.conv.IsConstant("XYZ123_")).IsTrue()
	AssertThat(c.conv.IsConstant("ABC")).IsTrue()
	AssertThat(c.conv.IsConstant("ABCdef")).IsFalse()
	AssertThat(c.conv.IsConstant("aBC")).IsFalse()
	AssertThat(c.conv.IsConstant("A")).IsFalse()
	AssertThat(c.conv.IsConstant("_XYZ123")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$XYZ123_")).IsTrue()
	AssertThat(c.conv.IsConstant("a$b$ABC_DEF")).IsTrue()
	AssertThat(c.conv.IsConstant("a$b$A")).IsTrue()
	AssertThat(c.conv.IsConstant("a$b$a")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$ABCdef")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$aBC")).IsFalse()
	AssertThat(c.conv.IsConstant("a$b$")).IsFalse()
	AssertThat(c.conv.IsConstant("$")).IsFalse()
	AssertThat(c.conv.IsConstant("$A")).IsTrue()
	AssertThat(c.conv.IsConstant("$a")).IsFalse()
}

// Java: GoogleCodingConventionTest.testIsConstantKey()
func (c *GoogleCodingConventionTest) TestIsConstantKey() {
	AssertThat(c.conv.IsConstantKey("a")).IsFalse()
	AssertThat(c.conv.IsConstantKey("XYZ123_")).IsTrue()
	AssertThat(c.conv.IsConstantKey("ABC")).IsTrue()
	AssertThat(c.conv.IsConstantKey("ABCdef")).IsFalse()
	AssertThat(c.conv.IsConstantKey("aBC")).IsFalse()
	AssertThat(c.conv.IsConstantKey("A")).IsTrue()
	AssertThat(c.conv.IsConstantKey("_XYZ123")).IsFalse()
	AssertThat(c.conv.IsConstantKey("a$b$ABC")).IsFalse()
	AssertThat(c.conv.IsConstantKey("$")).IsFalse()
}

// Java: GoogleCodingConventionTest.testExportedName()
func (c *GoogleCodingConventionTest) TestExportedName() {
	AssertThat(c.conv.IsExported("_a")).IsTrue()
	AssertThat(c.conv.IsExported("_a_")).IsTrue()
	AssertThat(c.conv.IsExported("a")).IsFalse()
	AssertThat(c.conv.IsExportedWithLocal("$super", false)).IsFalse()
	AssertThat(c.conv.IsExportedWithLocal("$super", true)).IsTrue()
	AssertThat(c.conv.IsExported("$super")).IsTrue()
}

// Java: GoogleCodingConventionTest.testEnumKey()
func (c *GoogleCodingConventionTest) TestEnumKey() {
	AssertThat(c.conv.IsValidEnumKey("A")).IsTrue()
	AssertThat(c.conv.IsValidEnumKey("123")).IsTrue()
	AssertThat(c.conv.IsValidEnumKey("FOO_BAR")).IsTrue()
	AssertThat(c.conv.IsValidEnumKey("a")).IsFalse()
	AssertThat(c.conv.IsValidEnumKey("someKeyInCamelCase")).IsFalse()
	AssertThat(c.conv.IsValidEnumKey("_FOO_BAR")).IsFalse()
}

// Java: GoogleCodingConventionTest.testInheritanceDetection1()
func (c *GoogleCodingConventionTest) TestInheritanceDetection1() {
	c.AssertNotClassDefining("goog.foo(A, B);")
}

// Java: GoogleCodingConventionTest.testInheritanceDetection2()
func (c *GoogleCodingConventionTest) TestInheritanceDetection2() {
	c.AssertDefinesClasses("goog.inherits(A, B);", "A", "B")
}

// Java: GoogleCodingConventionTest.testInheritanceDetection3()
func (c *GoogleCodingConventionTest) TestInheritanceDetection3() {
	c.AssertNotClassDefining("A.inherits(B);")
}

// Java: GoogleCodingConventionTest.testInheritanceDetection4()
func (c *GoogleCodingConventionTest) TestInheritanceDetection4() {
	c.AssertDefinesClasses("goog.inherits(goog.A, goog.B);", "goog.A", "goog.B")
}

// Java: GoogleCodingConventionTest.testInheritanceDetection5()
func (c *GoogleCodingConventionTest) TestInheritanceDetection5() {
	c.AssertNotClassDefining("goog.A.inherits(goog.B);")
}

// Java: GoogleCodingConventionTest.testInheritanceDetection6()
func (c *GoogleCodingConventionTest) TestInheritanceDetection6() {
	c.AssertNotClassDefining("A.inherits(this.B);")
}

// Java: GoogleCodingConventionTest.testInheritanceDetection7()
func (c *GoogleCodingConventionTest) TestInheritanceDetection7() {
	c.AssertNotClassDefining("this.A.inherits(B);")
}

// Java: GoogleCodingConventionTest.testInheritanceDetection8()
func (c *GoogleCodingConventionTest) TestInheritanceDetection8() {
	c.AssertDefinesClasses("goog.inherits(A, B, C);", "A", "B")
}

// Java: GoogleCodingConventionTest.testInheritanceDetection9()
func (c *GoogleCodingConventionTest) TestInheritanceDetection9() {
	c.AssertNotClassDefining("A.mixin(B.prototype);")
}

// Java: GoogleCodingConventionTest.testInheritanceDetectionPostCollapseProperties()
func (c *GoogleCodingConventionTest) TestInheritanceDetectionPostCollapseProperties() {
	c.AssertDefinesClasses("goog$inherits(A, B);", "A", "B")
	c.AssertNotClassDefining("goog$inherits(A);")
}

// Java: GoogleCodingConventionTest.testPackageNames()
func (c *GoogleCodingConventionTest) TestPackageNames() {
	c.AssertPackageName("foo.js", "")
	c.AssertPackageName("foo/bar.js", "foo")
	c.AssertPackageName("foo/bar/baz.js", "foo/bar")
	c.AssertPackageName("foo/bar/baz/quux.js", "foo/bar/baz")
	c.AssertPackageName("foo/test/bar.js", "foo")
	c.AssertPackageName("foo/tests/bar.js", "foo")
	c.AssertPackageName("foo/testing/bar.js", "foo")
	c.AssertPackageName("foo/jstest/bar.js", "foo/jstest")
	c.AssertPackageName("foo/bar/test/baz.js", "foo/bar")
	c.AssertPackageName("foo/bar/tests/baz.js", "foo/bar")
	c.AssertPackageName("foo/bar/testing/baz.js", "foo/bar")
	c.AssertPackageName("foo/bar/testament/baz.js", "foo/bar/testament")
	c.AssertPackageName("foo/test/bar/baz.js", "foo/test/bar")
	c.AssertPackageName("foo/bar/baz/test/quux.js", "foo/bar/baz")
	c.AssertPackageName("foo/bar/baz/tests/quux.js", "foo/bar/baz")
	c.AssertPackageName("foo/bar/baz/testing/quux.js", "foo/bar/baz")
	c.AssertPackageName("foo/bar/baz/unittests/quux.js", "foo/bar/baz/unittests")
	c.AssertPackageName("foo/bar/test/baz/quux.js", "foo/bar/test/baz")
	c.AssertPackageName("foo/test/bar/baz/quux.js", "foo/test/bar/baz")
	c.AssertPackageName("bazel-out/host/genfiles/bar/baz/quux.js", "bar/baz")
	c.AssertPackageName("bazel-out/host/genfiles/foo/test/bar.js", "foo")
	c.AssertPackageName("bazel-out/host/bin/bar/baz/quux.js", "bar/baz")
	c.AssertPackageName("bazel-out/host/bin/foo/test/bar.js", "foo")
}

// Java: GoogleCodingConventionTest.assertPackageName(String, String)
func (c *GoogleCodingConventionTest) AssertPackageName(filename string, expectedPackageName string) {
	var sourceFile rhino.StaticSourceFile = SourceFile_FromCode(filename, "")
	AssertThat(c.conv.GetPackageName(sourceFile)).IsEqualTo(expectedPackageName)
}

// Java: GoogleCodingConventionTest.assertNotClassDefining(String)
func (c *GoogleCodingConventionTest) AssertNotClassDefining(code string) {
	n := c.ParseTestCode(code)
	AssertThat(c.conv.GetClassesDefinedByCall(n.GetFirstChild())).IsNull()
}

// Java: GoogleCodingConventionTest.assertDefinesClasses(String, String, String)
func (c *GoogleCodingConventionTest) AssertDefinesClasses(code string, subclassName string, superclassName string) {
	n := c.ParseTestCode(code)
	classes := c.conv.GetClassesDefinedByCall(n.GetFirstChild())
	AssertThat(classes).IsNotNull()
	AssertThat(classes.subclassName).IsEqualTo(subclassName)
	AssertThat(classes.superclassName).IsEqualTo(superclassName)
}

// Java: GoogleCodingConventionTest.parseTestCode(String)
func (c *GoogleCodingConventionTest) ParseTestCode(code string) *rhino.Node {
	compiler := NewCompiler()
	return compiler.ParseTestCode(code).GetFirstChild()
}

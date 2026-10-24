// Port of com/google/javascript/jscomp/J2clPassTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/J2clPassTest.java)

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
	"testing"
)

// J2clPassTest is the tests for J2clPass.
type J2clPassTest struct {
	CompilerTestCaseBase
}

// Java: J2clPassTest()
func NewJ2clPassTest() *J2clPassTest {
	c := &J2clPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestJ2clPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewJ2clPassTest() })
}

// JUnitBefore runs the @Before methods: the superclass's setUp(), then
// customSetUp().
func (c *J2clPassTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: J2clPassTest.customSetUp()
func (c *J2clPassTest) CustomSetUp() {
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableNormalize()
}

// Java: J2clPassTest.getProcessor(Compiler)
func (c *J2clPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewJ2clPass(compiler)
}

// Java: J2clPassTest.createCompiler()
func (c *J2clPassTest) CreateCompiler() *Compiler {
	compiler := c.CompilerTestCaseBase.CreateCompiler()
	J2clSourceFileChecker_MarkToRunJ2clPasses(compiler)
	return compiler
}

// Java: J2clPassTest.testQualifiedInlines_arrays()
func (c *J2clPassTest) TestQualifiedInlines_arrays() {
	// Function definitions and calls are qualified globals.
	declarations :=
		`class Arrays {
  static $create() { return 1; }
  static $init() { return 2; }
  static $instanceIsOfType() { return 3; }
  static $castTo() { return 4; }
  static $stampType() { return 5; }
}
`

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Arrays.impl.java.js",
				declarations+`alert(Arrays.$create());
alert(Arrays.$init());
alert(Arrays.$instanceIsOfType());
alert(Arrays.$castTo());
alert(Arrays.$stampType());
`)),
		c.Expected(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Arrays.impl.java.js",
				declarations+`alert(1);
alert(2);
alert(3);
alert(4);
alert(5);
`)))
}

// Java: J2clPassTest.testQualifiedInlines_casts()
func (c *J2clPassTest) TestQualifiedInlines_casts() {
	// Function definitions and calls are qualified globals.
	declarations :=
		`class Casts {
  static $to() { return 1; }
}
`

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Casts.impl.java.js",
				declarations+`alert(Casts.$to());
`)),
		c.Expected(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Casts.impl.java.js",
				declarations+`alert(1);
`)))
}

// Java: J2clPassTest.testQualifiedInlines_markImplementor()
func (c *J2clPassTest) TestQualifiedInlines_markImplementor() {
	// Function definitions and calls are qualified globals.
	declarations :=
		`class FooInterface {
  static $markImplementor(classDef) {
    classDef.$implements__FooInterface = true;
  }
}
`

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"name/doesnt/matter/Foo.impl.java.js",
				declarations+`/** @constructor */
var Foo = function() {};
FooInterface.$markImplementor(Foo);
`)),
		c.Expected(
			SourceFile_FromCode(
				"name/doesnt/matter/Foo.impl.java.js",
				declarations+`/** @constructor */
var Foo = function() {};
{Foo.$implements__FooInterface = true;}
`)))
	scriptNode := c.GetLastCompiler().GetJsRoot().GetLastChild()
	AssertNode(scriptNode).IsScript()
	inlinedBlock := scriptNode.GetLastChild()
	varStatement := inlinedBlock.GetPrevious()

	declaredFooNameNodeSubject :=
		AssertNode(varStatement).IsVar().HasOneChildThat().IsName("Foo")
	declaredColorSubject := declaredFooNameNodeSubject.HasColorThat()

	inlinedFooNameNodeSubject :=
		AssertNode(inlinedBlock).
			IsBlock().
			HasOneChildThat().
			IsExprResult().
			HasOneChildThat().
			IsAssign().
			HasFirstChildThat().
			IsGetProp().
			HasOneChildThat().
			IsName("Foo")
	inlinedColorSubject := inlinedFooNameNodeSubject.HasColorThat()

	inlinedColorSubject.IsEqualToSubject(declaredColorSubject)
}

// Java: J2clPassTest.testRenamedQualifierStillInlines_arrays()
func (c *J2clPassTest) TestRenamedQualifierStillInlines_arrays() {
	// Function definitions and calls are qualified globals.
	declarations :=
		`var $jscomp = {};
$jscomp.scope = {};
$jscomp.scope.Arrays = class {;
  static $create() { return 1; }
  static $init() { return 2; }
  static $instanceIsOfType() { return 3; }
  static $castTo() { return 4; }
}
`

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Arrays.impl.java.js",
				declarations+`alert($jscomp.scope.Arrays.$create());
alert($jscomp.scope.Arrays.$init());
alert($jscomp.scope.Arrays.$instanceIsOfType());
alert($jscomp.scope.Arrays.$castTo());
`)),
		c.Expected(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Arrays.impl.java.js",
				declarations+`alert(1);
alert(2);
alert(3);
alert(4);
`)))
}

// Java: J2clPassTest.testRenamedQualifierStillInlines_casts()
func (c *J2clPassTest) TestRenamedQualifierStillInlines_casts() {
	// Function definitions and calls are qualified globals.
	declarations :=
		`var $jscomp = {};
$jscomp.scope = {};
$jscomp.scope.Casts = class {
  static $to() { return 1; }
}
`

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Casts.impl.java.js",
				declarations+`alert($jscomp.scope.Casts.$to());
`)),
		c.Expected(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Casts.impl.java.js",
				declarations+`alert(1);
`)))
}

// Java: J2clPassTest.testRenamedQualifierStillInlines_markImplementor()
func (c *J2clPassTest) TestRenamedQualifierStillInlines_markImplementor() {
	// Function definitions and calls are qualified globals.
	declarations :=
		`var $jscomp = {};
$jscomp.scope = {};
$jscomp.scope.FooInterface = class {
  static $markImplementor(classDef) {
    classDef.$implements__FooInterface = true;
  }
}
`

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"name/doesnt/matter/Foo.impl.java.js",
				declarations+`$jscomp.scope.Foo = function() {};
$jscomp.scope.FooInterface.$markImplementor($jscomp.scope.Foo);
`)),
		c.Expected(
			SourceFile_FromCode(
				"name/doesnt/matter/Foo.impl.java.js",
				declarations+`$jscomp.scope.Foo = function() {};
{$jscomp.scope.Foo.$implements__FooInterface = true;}
`)))
}

// Java: J2clPassTest.testUnexpectedFunctionDoesntInline()
func (c *J2clPassTest) TestUnexpectedFunctionDoesntInline() {
	// Arrays functions.
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Arrays.impl.java.js",
				`var Arrays = class {
  static fooBar() { return 4; }
}

alert(Arrays.fooBar());
`)))

	// Casts functions.
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Casts.impl.java.js",
				`var Casts = class {
  static fooBar() { return 4; }
}

alert(Casts.fooBar());
`)))

	// No applicable for $markImplementor() inlining since it is not limited to just certain class
	// files and so there are no specific files in which "other" functions should be ignored.
}

// Java: J2clPassTest.testUnqualifiedDoesntInline()
func (c *J2clPassTest) TestUnqualifiedDoesntInline() {
	// Arrays functions.
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Arrays.impl.java.js",
				`var $create = function() { return 1; }
var $init = function() { return 2; }
var $instanceIsOfType = function() { return 3; }
var $castTo = function() { return 4; }

alert($create());
alert($init());
alert($instanceIsOfType());
alert($castTo());
`)))

	// Casts functions.
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Casts.impl.java.js",
				`var to = function() { return 1; }

alert(to());
`)))

	// Interface $markImplementor() functions.
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"name/doesnt/matter/Foo.impl.java.js",
				`var $markImplementor = function(classDef) {
  classDef.$implements__FooInterface = true;
}

var Foo = function() {};
$markImplementor(Foo);
`)))
}

// Java: J2clPassTest.testWrongFileNameDoesntInline()
func (c *J2clPassTest) TestWrongFileNameDoesntInline() {
	// Arrays functions.
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"Arrays2.impl.java.js",
				`var Arrays = function() {};
Arrays.$create = function() { return 1; }
Arrays.$init = function() { return 2; }
Arrays.$instanceIsOfType = function() { return 3; }
Arrays.$castTo = function() { return 4; }

alert(Arrays.$create());
alert(Arrays.$init());
alert(Arrays.$instanceIsOfType());
alert(Arrays.$castTo());
`)))

	// Casts functions.
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"Casts2.impl.java.js",
				`var Casts = function() {};
Casts.to = function() { return 1; }

alert(Casts.to());
`)))

	// No applicable for $markImplementor() inlining since it is not limited to just certain class
	// files.
}

// Java: J2clPassTest.testMarksChanges()
func (c *J2clPassTest) TestMarksChanges() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Casts.impl.java.js",
				`var Casts = function() {};
Casts.$to = function(instance) { return instance; }

alert(Casts.$to(function(a) { return a; }));
`)),
		c.Expected(
			SourceFile_FromCode(
				"j2cl/transpiler/vmbootstrap/Casts.impl.java.js",
				`var Casts = function() {};
Casts.$to = function(instance) { return instance; }

alert(function(a) { return a; });
`)))
}

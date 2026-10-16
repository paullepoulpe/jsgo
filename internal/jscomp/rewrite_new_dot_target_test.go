// Port of com/google/javascript/jscomp/RewriteNewDotTargetTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/RewriteNewDotTargetTest.java)

/*
 * Copyright 2019 The Closure Compiler Authors.
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

type RewriteNewDotTargetTest struct {
	CompilerTestCaseBase
}

// Java: RewriteNewDotTargetTest()
func NewRewriteNewDotTargetTest() *RewriteNewDotTargetTest {
	c := &RewriteNewDotTargetTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteNewDotTargetTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteNewDotTargetTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then enableTypeCheckBeforePass().
func (c *RewriteNewDotTargetTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.EnableTypeCheckBeforePass()
}

// Java: RewriteNewDotTargetTest.enableTypeCheckBeforePass()
func (c *RewriteNewDotTargetTest) EnableTypeCheckBeforePass() {
	c.EnableNormalize()
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: RewriteNewDotTargetTest.getProcessor(Compiler)
func (c *RewriteNewDotTargetTest) GetProcessor(compiler *Compiler) CompilerPass {
	return PeepholeTranspilationsPass_Create(
		compiler, []AbstractPeepholeTranspilation{NewRewriteNewDotTarget(compiler)})
}

// Java: RewriteNewDotTargetTest.getOptions()
func (c *RewriteNewDotTargetTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSTABLE)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT)
	return options
}

// Java: RewriteNewDotTargetTest.testNewTargetInNonConstructorMethod()
func (c *RewriteNewDotTargetTest) TestNewTargetInNonConstructorMethod() {
	c.TestError(
		`/** @constructor */
function Foo() { new.target; } // not an ES6 class constructor
`,
		TranspilationUtil_CANNOT_CONVERT_YET)

	c.TestError(
		"function foo() { new.target; }", // not a constructor at all
		TranspilationUtil_CANNOT_CONVERT_YET)

	c.TestError(
		`class Foo {
  method() { // not a constructor
    new.target;
  }
}
`,
		TranspilationUtil_CANNOT_CONVERT_YET)
}

// Java: RewriteNewDotTargetTest.testNewTargetInEs6Constructor()
func (c *RewriteNewDotTargetTest) TestNewTargetInEs6Constructor() {
	c.Test(
		`class Foo {
  constructor() {
    new.target;
    () => new.target; // works in arrow functions, too
  }
}
`,
		`class Foo {
  constructor() {
    this.constructor;
    () => { return this.constructor; }; // works in arrow functions, too
  }
}
`)
}

// Java: RewriteNewDotTargetTest.testNewTargetInEs6Constructor_superCtorReturnsObject()
func (c *RewriteNewDotTargetTest) TestNewTargetInEs6Constructor_superCtorReturnsObject() {
	// TODO(b/157140030): This is an unexpected behaviour change.
	c.Test(
		`/** @constructor */
function Base() {
  return {};
}

class Child extends Base {
  constructor() {
    super();
    new.target; // Child
  }
}
`,
		`/** @constructor */
function Base() {
  return {};
}

class Child extends Base {
  constructor() {
    super();
    this.constructor; // Object
  }
}
`)
}

// Java: RewriteNewDotTargetTest.testNewTargetInEs6Constructor_beforeSuper()
func (c *RewriteNewDotTargetTest) TestNewTargetInEs6Constructor_beforeSuper() {
	c.Test(
		`class Foo extends Object {
  constructor() {
    new.target;
    super();
  }
}
`,
		"class Foo extends Object {\n  constructor() {\n    this.constructor; // `this` before `super`\n    super();\n  }\n}\n")
}

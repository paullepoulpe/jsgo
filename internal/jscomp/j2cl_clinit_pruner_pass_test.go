// Port of com/google/javascript/jscomp/J2clClinitPrunerPassTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/J2clClinitPrunerPassTest.java)

/*
 * Copyright 2016 The Closure Compiler Authors.
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

// J2clClinitPrunerPassTest tests J2clClinitPrunerPass.
type J2clClinitPrunerPassTest struct {
	CompilerTestCaseBase
}

// Java: J2clClinitPrunerPassTest()
func NewJ2clClinitPrunerPassTest() *J2clClinitPrunerPassTest {
	c := &J2clClinitPrunerPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestJ2clClinitPrunerPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewJ2clClinitPrunerPassTest() })
}

// Java: J2clClinitPrunerPassTest.getProcessor(Compiler)
func (c *J2clClinitPrunerPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewJ2clClinitPrunerPass(
		compiler, compiler.GetChangeTracker().GetChangedScopeNodesForPass("J2clClinitPrunerPass"))
}

// Java: J2clClinitPrunerPassTest.createCompiler()
func (c *J2clClinitPrunerPassTest) CreateCompiler() *Compiler {
	compiler := c.CompilerTestCaseBase.CreateCompiler()
	J2clSourceFileChecker_MarkToRunJ2clPasses(compiler)
	return compiler
}

// Java: J2clClinitPrunerPassTest.setUp()
func (c *J2clClinitPrunerPassTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()

	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableComputeSideEffects()
}

// Java: J2clClinitPrunerPassTest.testRemoveDuplicates()
func (c *J2clClinitPrunerPassTest) TestRemoveDuplicates() {
	c.Test(
		`var someClass = {};
someClass.$clinit = function() {}
someClass.someOtherFunction = function() {
  someClass.$clinit();
  someClass.$clinit();
};
`,
		`var someClass = {};
someClass.$clinit = function() {}
someClass.someOtherFunction = function() {
  someClass.$clinit();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRemoveDuplicates_commaExpressions()
func (c *J2clClinitPrunerPassTest) TestRemoveDuplicates_commaExpressions() {
	c.Test(
		`var someClass = {};
someClass.$clinit = function() {}
someClass.someOtherFunction = function() {
  (someClass.$clinit(), true);
  (someClass.$clinit(), true);
};
`,
		`var someClass = {};
someClass.$clinit = function() {}
someClass.someOtherFunction = function() {
  (someClass.$clinit(), true);
  (void 0, true);
};
`)
}

// Java: J2clClinitPrunerPassTest.testRemoveDuplicates_controlBlocks()
func (c *J2clClinitPrunerPassTest) TestRemoveDuplicates_controlBlocks() {
	c.Test(
		`var someClass = {};
someClass.$clinit = function() {}
someClass.someOtherFunction = function() {
  someClass.$clinit();
  if (true) {
    someClass.$clinit();
    while(true) {
      someClass.$clinit();
    }
  } else {
    someClass.$clinit();
  }
  var a = (someClass.$clinit(), true) ? (someClass.$clinit(), void 0) : void 0;
  var b = function() { someClass.$clinit(); };
  var c = function c() { someClass.$clinit(); };
  var d = function*() { someClass.$clinit(); };
  var e = async function() { someClass.$clinit(); };
  var f = () => { someClass.$clinit(); };
  [].forEach(function() { someClass.$clinit(); });
  someClass.$clinit();
};
`,
		`var someClass = {};
someClass.$clinit = function() {}
someClass.someOtherFunction = function() {
  someClass.$clinit();
  if (true) {
    while(true) {
    }
  } else {
  }
  var a = (void 0, true) ? (void 0, void 0) : void 0;
  var b = function() {};
  var c = function c() {};
  var d = function*() {};
  var e = async function() {};
  var f = () => {};
  [].forEach(function() {});
};
`)
}

// Java: J2clClinitPrunerPassTest.testRemoveDuplicates_selfRemoval()
func (c *J2clClinitPrunerPassTest) TestRemoveDuplicates_selfRemoval() {
	c.Test(
		`var someClass = {};
someClass.$clinit = function() {
  someClass.$clinit();
};
`,
		`var someClass = {};
someClass.$clinit = function() {}
`)

	c.Test(
		`function someClass$$0clinit() {
  someClass$$0clinit();
}
`,
		"function someClass$$0clinit() {}")
}

// Java: J2clClinitPrunerPassTest.testRemoveDuplicates_jumpFunctionDeclarations()
func (c *J2clClinitPrunerPassTest) TestRemoveDuplicates_jumpFunctionDeclarations() {
	c.Test(
		`var someClass = {};
someClass.$clinit = function() {}
someClass.someOtherFunction = function() {
  myFunc();
  someClass.$clinit();
  function myFunc() {
    someClass.$clinit();
    someClass.$clinit();
  }
};
`,
		"var someClass = {};\nsomeClass.$clinit = function() {}\nsomeClass.someOtherFunction = function() {\n  myFunc();\n// Control flow analysis doesn't understand that this was already called in `myFunc`.\n  someClass.$clinit();\n  function myFunc() {\n    someClass.$clinit();\n  }\n};\n")
}

// Java: J2clClinitPrunerPassTest.testRemoveDuplicates_jumpDefaults()
func (c *J2clClinitPrunerPassTest) TestRemoveDuplicates_jumpDefaults() {
	c.Test(
		`var someClass = {};
someClass.$clinit = function() {}
someClass.someOtherFunction = function(p=someClass.$clinit()) {
  var { x = someClass.$clinit() } = someClass.$clinit();
  someClass.$clinit();
};
`,
		`var someClass = {};
someClass.$clinit = function() {}
someClass.someOtherFunction = function(p=someClass.$clinit()) {
  var { x = someClass.$clinit() } = someClass.$clinit();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRemoveDuplicates_avoidControlBlocks()
func (c *J2clClinitPrunerPassTest) TestRemoveDuplicates_avoidControlBlocks() {
	c.TestSame(
		`var someClass = {};
someClass.$clinit = function() {}
someClass.anotherMethod = function() {
  (false && someClass.$clinit());
  (true || someClass.$clinit());
  if (true) {
    someClass.$clinit();
  } else {
    someClass.$clinit();
  }
  while(false) {
    someClass.$clinit();
  }
  for(;false;) {
    someClass.$clinit();
  }
  try {
    someClass.$clinit();
  } catch(e) {
    someClass.$clinit();
  }
  switch(2) {
    case 1: someClass.$clinit(); break;
    case 2: break;
    default: someClass.$clinit();
  }
  var a = true ? (someClass.$clinit(), void 0) : void 0;
  var b = function() { someClass.$clinit(); }
  var c = function*() { someClass.$clinit(); }
  var d = async function() { someClass.$clinit(); }
  var e = () => { someClass.$clinit(); }
  someClass.$clinit();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRemoveDuplicates_avoidRemovalAcrossScriptRoots()
func (c *J2clClinitPrunerPassTest) TestRemoveDuplicates_avoidRemovalAcrossScriptRoots() {
	c.TestSame(
		newCompilerTestCase_FlatSources(
			[]*SourceFile{
				SourceFile_FromCode(
					"file1",
					`var someClass = {};
someClass.$clinit = function() {}
someClass.$clinit();
`),
				SourceFile_FromCode("file2", "someClass.$clinit();")}))
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_returnCtor()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_returnCtor() {
	c.Test(
		`var Foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  Foo.$clinit();
  return new Foo();
};
`,
		`var Foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  return new Foo();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_returnCtor_es6()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_returnCtor_es6() {
	c.T().Skip("@Ignore(\"Improve look ahead pruner to handle ES6 classes\")")
	c.Test(
		`class Foo {
  constructor(){
    Foo.$clinit();
  }
  static ctor() {
    Foo.$clinit();
    return new Foo();
  }
};
`,
		`class Foo {
  constructor(){
    Foo.$clinit();
  }
  static ctor() {
    return new Foo();
  }
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_returnCall()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_returnCall() {
	c.Test(
		`var foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  Foo.$clinit();
  return foo();
};
`,
		`var foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  return foo();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_returnCall_es6()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_returnCall_es6() {
	c.T().Skip("@Ignore(\"Improve look ahead pruner to handle ES6 classes\")")
	c.Test(
		`class Foo {
  static bar(){
    Foo.$clinit();
  }
  static zoo() {
    Foo.$clinit();
    return Foo.bar();
  }
};
`,
		`class Foo {
  static bar(){
    Foo.$clinit();
  }
  static zoo() {
    return bar();
  }
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_exprResult()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_exprResult() {
	c.Test(
		`var foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  Foo.$clinit();
  foo();
};
`,
		`var foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  foo();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_var()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_var() {
	c.Test(
		`var foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  Foo.$clinit();
  var x = foo();
};
`,
		`var foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  var x = foo();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_let()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_let() {
	c.Test(
		`var foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  Foo.$clinit();
  let x = foo();
};
`,
		`var foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  let x = foo();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_const()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_const() {
	c.Test(
		`var foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  Foo.$clinit();
  const x = foo();
};
`,
		`var foo = function() {
  Foo.$clinit();
};
Foo.ctor = function() {
  const x = foo();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_literalArgs()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_literalArgs() {
	c.Test(
		`var Foo = function(a) {
  Foo.$clinit();
};
Foo.ctor = function() {
  Foo.$clinit();
  return new Foo(1);
};
`,
		`var Foo = function(a) {
  Foo.$clinit();
};
Foo.ctor = function() {
  return new Foo(1);
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_paramArgs()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_paramArgs() {
	c.Test(
		`var Foo = function(a, b) {
  Foo.$clinit();
};
Foo.ctor = function(a) {
  Foo.$clinit();
  return new Foo(a, 1);
};
`,
		`var Foo = function(a, b) {
  Foo.$clinit();
};
Foo.ctor = function(a) {
  return new Foo(a, 1);
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_unsafeArgs()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_unsafeArgs() {
	c.TestSame(
		`var Foo = function(a) {
  Foo.$clinit();
};
Foo.STATIC_VAR = null;
Foo.ctor = function() {
  Foo.$clinit();
  return new Foo(Foo.STATIC_VAR);
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_otherClinit()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_otherClinit() {
	c.TestSame(
		`var Foo = function() {
  Foo1.$clinit();
};
Foo.ctor = function() {
  Foo.$clinit();
  return new Foo();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_clinitNotFirstStatement()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_clinitNotFirstStatement() {
	c.TestSame(
		`var Foo = function() {
  var x = 1;
  Foo.$clinit();
};
Foo.ctor = function() {
  Foo.$clinit();
  return new Foo();
};
`)
}

// Java: J2clClinitPrunerPassTest.testRedundantClinit_recursiveCall()
func (c *J2clClinitPrunerPassTest) TestRedundantClinit_recursiveCall() {
	c.TestSame(
		`var foo = function() {
  Foo1.$clinit();
  foo();
};
`)
}

// Java: J2clClinitPrunerPassTest.testFoldClinit_es5()
func (c *J2clClinitPrunerPassTest) TestFoldClinit_es5() {
	c.Test(
		`var someClass = {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
};
`,
		`var someClass = {};
someClass.$clinit = function() {};
`)
}

// Java: J2clClinitPrunerPassTest.testFoldClinit_es6()
func (c *J2clClinitPrunerPassTest) TestFoldClinit_es6() {
	c.Test(
		`class someClass {
  static $clinit() {
    someClass.$clinit = function() {};
  }
}
`,
		`class someClass {
  static $clinit() {}
}
`)
}

// Java: J2clClinitPrunerPassTest.testFoldClinit_classHierarchy_es5()
func (c *J2clClinitPrunerPassTest) TestFoldClinit_classHierarchy_es5() {
	c.Test(
		`var someClass = {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
};
var someChildClass = {};
someChildClass.$clinit = function() {
  someChildClass.$clinit = function() {};
  someClass.$clinit();
};
someChildClass.someFunction = function() {
  someChildClass.$clinit();
  someClass.$clinit();
};
`,
		`var someClass = {};
someClass.$clinit = function() {};
var someChildClass = {};
someChildClass.$clinit = function() {};
someChildClass.someFunction = function() {};
`)
}

// Java: J2clClinitPrunerPassTest.testFoldClinit_classHierarchy_es6()
func (c *J2clClinitPrunerPassTest) TestFoldClinit_classHierarchy_es6() {
	c.Test(
		`class someClass {
  static $clinit() {
    someClass.$clinit = function() {};
  }
}
class someChildClass {
  static $clinit() {
    someChildClass.$clinit = function() {};
    someClass.$clinit();
  }

  static someFunction() {
    someChildClass.$clinit();
    someClass.$clinit();
  }
}
`,
		`class someClass {
  static $clinit() {}
}
class someChildClass {
  static $clinit() {}

  static someFunction() {}
}
`)
}

// Java: J2clClinitPrunerPassTest.testFoldClinit_classHierarchyNonEmpty_es5()
func (c *J2clClinitPrunerPassTest) TestFoldClinit_classHierarchyNonEmpty_es5() {
	c.Test(
		`var someClass = {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
  somefn();
};
var someChildClass = {};
someChildClass.$clinit = function() {
  someChildClass.$clinit = function() {};
  someClass.$clinit();
};
someChildClass.someFunction = function() {
  someChildClass.$clinit();
};
`,
		`var someClass = {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
  somefn();
};
var someChildClass = {};
someChildClass.$clinit = function() {
  someClass.$clinit();
};
someChildClass.someFunction = function() {
  someClass.$clinit();
};
`)
}

// Java: J2clClinitPrunerPassTest.testFoldClinit_classHierarchyNonEmpty_es6()
func (c *J2clClinitPrunerPassTest) TestFoldClinit_classHierarchyNonEmpty_es6() {
	c.Test(
		`class someClass {
  static $clinit() {
    someClass.$clinit = function() {};
    somefn();
  }
}
class someChildClass {
  static $clinit() {
    someChildClass.$clinit = function() {};
    someClass.$clinit();
  }

  static someFunction() {
    someChildClass.$clinit();
  }
}
`,
		`class someClass {
  static $clinit() {
    someClass.$clinit = function() {};
    somefn();
  }
}
class someChildClass {
  static $clinit() {
    someClass.$clinit();
  }

  static someFunction() {
    someClass.$clinit();
  }
}
`)
}

// Java: J2clClinitPrunerPassTest.testFoldClinit_invalidCandidates_es5()
func (c *J2clClinitPrunerPassTest) TestFoldClinit_invalidCandidates_es5() {
	c.TestSame(
		`var someClass = /** @constructor */ function() {};
someClass.foo = function() {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
  someClass.foo();
};
`)
	c.TestSame(
		`var someClass = {}, otherClass = {};
someClass.$clinit = function() {
  otherClass.$clinit = function() {};
};
`)
	c.TestSame(
		`var someClass = {};
someClass.$notClinit = function() {
  someClass.$notClinit = function() {};
};
`)
}

// Java: J2clClinitPrunerPassTest.testFoldClinit_invalidCandidates_es6()
func (c *J2clClinitPrunerPassTest) TestFoldClinit_invalidCandidates_es6() {
	c.TestSame(
		`class someClass {
  static foo() {}
  static $clinit() {
    someClass.$clinit = function() {};
    someClass.foo();
  }
}
`)
	c.TestSame(
		`class someClass {
  static $clinit() {
    otherClass.$clinit = function() {};
  }
}
class otherClass {}
`)
	c.TestSame(
		`class someClass {
  static $notClinit() {
    someClass.$notClinit = function() {};
  }
}
`)
}

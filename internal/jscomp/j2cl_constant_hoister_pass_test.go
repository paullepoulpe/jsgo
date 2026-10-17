// Port of com/google/javascript/jscomp/J2clConstantHoisterPassTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/J2clConstantHoisterPassTest.java)

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

// J2clConstantHoisterPassTest tests J2clConstantHoisterPass.
type J2clConstantHoisterPassTest struct {
	CompilerTestCaseBase
}

// Java: J2clConstantHoisterPassTest()
func NewJ2clConstantHoisterPassTest() *J2clConstantHoisterPassTest {
	c := &J2clConstantHoisterPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestJ2clConstantHoisterPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewJ2clConstantHoisterPassTest() })
}

// Java: J2clConstantHoisterPassTest.getProcessor(Compiler)
func (c *J2clConstantHoisterPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewJ2clConstantHoisterPass(compiler)
}

// Java: J2clConstantHoisterPassTest.createCompiler()
func (c *J2clConstantHoisterPassTest) CreateCompiler() *Compiler {
	compiler := c.CompilerTestCaseBase.CreateCompiler()
	J2clSourceFileChecker_MarkToRunJ2clPasses(compiler)
	return compiler
}

// Java: J2clConstantHoisterPassTest.testHoistClinitConstantAssignments()
func (c *J2clConstantHoisterPassTest) TestHoistClinitConstantAssignments() {
	c.Test(
		`var someClass = /** @constructor */ function() {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
  someClass$foo = true;
  someClass$bar = 'hey';
  someClass$buzz = function() { return someClass$bar; };
};
var someClass$foo = false;
var someClass$bar = null;
var someClass$buzz = null;
`,
		`var someClass = /** @constructor */ function() {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
};
var someClass$foo = true;
var someClass$bar = 'hey';
var someClass$buzz = function(){ return someClass$bar; };
`)
}

// Java: J2clConstantHoisterPassTest.testHoistClinitConstantAssignments_assignmentAfterInit()
func (c *J2clConstantHoisterPassTest) TestHoistClinitConstantAssignments_assignmentAfterInit() {
	c.Test(
		`var someClass = /** @constructor */ function() {};
var someClass$foo = false;
var someClass$bar = null;
var someClass$buzz = null;
someClass.$clinit = function() {
  someClass.$clinit = function() {};
  someClass$foo = true;
  someClass$bar = 'hey';
  someClass$buzz = function() { return someClass$bar; };
};
`,
		`var someClass = /** @constructor */ function() {};
var someClass$foo = true;
var someClass$bar = 'hey';
var someClass$buzz = function(){ return someClass$bar; };
someClass.$clinit = function() {
  someClass.$clinit = function() {};
};
`)
}

// Java: J2clConstantHoisterPassTest.testHoistClinitConstantAssignments_stubInit()
func (c *J2clConstantHoisterPassTest) TestHoistClinitConstantAssignments_stubInit() {
	c.Test(
		`var someClass = /** @constructor */ function() {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
  someClass$foo = true;
  someClass$bar = 'hey';
  someClass$buzz = function() { return someClass$bar; };
};
var someClass$foo = false;
var someClass$bar;
var someClass$buzz;
`,
		`var someClass = /** @constructor */ function() {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
};
var someClass$foo = true;
var someClass$bar = 'hey';
var someClass$buzz = function(){ return someClass$bar; };
`)
}

// Java: J2clConstantHoisterPassTest.testHoistClinitConstantAssignments_avoidUnsafe()
func (c *J2clConstantHoisterPassTest) TestHoistClinitConstantAssignments_avoidUnsafe() {
	c.TestSame(
		`var someClass = /** @constructor */ function() {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
  someClass$foo1 = true // More than two assignments.
  someClass$foo2 = true // More than two assignments.
  someClass$bar = true // Also initialized by another method.
  someClass$zoo = true // Not initialized in declaration phase.
  someClass$hoo = someClass$zoo // Not literal value.
};
someClass.$otherMethod = function() {
  someClass$foo1 |= true // Compound assignment to make it less trivial to detect.
  someClass$foo2++ // Compound assignment to make it less trivial to detect.
  someClass$bar = true;
};
var someClass$foo1 = false;
var someClass$foo2 = false;
var someClass$hoo = false;
`)
}

// Java: J2clConstantHoisterPassTest.testHoistClinitConstantAssignments_devirtualized()
func (c *J2clConstantHoisterPassTest) TestHoistClinitConstantAssignments_devirtualized() {
	c.Test(
		`var someClass = /** @constructor */ function() {};
var someClass$$0clinit = function() {
  someClass$$0clinit = function() {};
  someClass$shouldHoist = true;
  someClass$shouldNotHoist = new Object();
};
var someClass$shouldHoist = false;
var someClass$shouldNotHoist = null;
`,
		`var someClass = /** @constructor */ function() {};
var someClass$$0clinit = function() {
  someClass$$0clinit = function() {};
  someClass$shouldNotHoist = new Object();
};
var someClass$shouldHoist = true;
var someClass$shouldNotHoist = null;
`)
}

// Java: J2clConstantHoisterPassTest.testHoistClinitConstantAssignments_avoidUnsafeFunc()
func (c *J2clConstantHoisterPassTest) TestHoistClinitConstantAssignments_avoidUnsafeFunc() {
	c.TestSame(
		`var someClass = /** @constructor */ function() {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
  var x = 10;
  someClass$foo = function() { return x; } // Depends on the scope of the clinit.
};
var someClass$foo = null;
`)
	c.TestSame(
		`var someClass = /** @constructor */ function() {};
someClass.$clinit = function() {
  someClass.$clinit = function() {};
  var x = 10;
  someClass$foo = function() { return 1; } // x is assumed to be used
};
var someClass$foo = null;
`)
}

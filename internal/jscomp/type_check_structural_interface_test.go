// Port of com/google/javascript/jscomp/TypeCheckStructuralInterfaceTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/TypeCheckStructuralInterfaceTest.java)

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
)

// TypeCheckStructuralInterfaceTest tests the type checker.
type TypeCheckStructuralInterfaceTest struct {
	typeCheckTests_noFixture
}

// Java: TypeCheckStructuralInterfaceTest()
func NewTypeCheckStructuralInterfaceTest() *TypeCheckStructuralInterfaceTest {
	c := &TypeCheckStructuralInterfaceTest{}
	return c
}

func TestTypeCheckStructuralInterfaceTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckStructuralInterfaceTest() })
}

var TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD = `/** @record */
function I5() {}
/** @type {I5} */
I5.prototype.next;

/** @constructor */
function C6() {}
/** @type {C6} */
C6.prototype.next;

/** @constructor */
function C6_1() {}
/** @type {C6} */
C6_1.prototype.next;

/** @constructor */
function C6_2() {}
/** @type {C6_1} */
C6_2.prototype.next;

/** @constructor */
function C6_3() {}
/** @type {C6_2} */
C6_3.prototype.next;

/** @constructor */
function C6_4() {}
/** @type {C6_3} */
C6_4.prototype.next;

/** @constructor */
function C6_5() {}
/** @type {C6_4} */
C6_5.prototype.next;

/** @constructor */
function C5() {}
/** @type {C6_5} */
C5.prototype.next;
`

var TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_NONMATCHING_CHAIN = `/** @record */
function I5() {}
/** @type {I5} */
I5.prototype.next;

/** @constructor */
function C6() {}
/** @type {number} */
C6.prototype.next;

/** @constructor */
function C6_1() {}
/** @type {C6} */
C6_1.prototype.next;

/** @constructor */
function C6_2() {}
/** @type {C6_1} */
C6_2.prototype.next;

/** @constructor */
function C6_3() {}
/** @type {C6_2} */
C6_3.prototype.next;

/** @constructor */
function C6_4() {}
/** @type {C6_3} */
C6_4.prototype.next;

/** @constructor */
function C6_5() {}
/** @type {C6_4} */
C6_5.prototype.next;

/** @interface */
function C5() {}
/** @type {C6_5} */
C5.prototype.next;
`

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function Interface1() {}
/** @type {number} */
Interface1.prototype.length;

/** @constructor */
function C1() {}
/** @type {number} */
C1.prototype.length;
`).
		AddSource(`/** @type{Interface1} */
var obj1;
/** @type{C1} */
var obj2 = new C1();
obj1 = obj2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching2()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function Interface1() {}
/** @type {number} */
Interface1.prototype.length;

/** @constructor */
function C1() {}
/** @type {number} */
C1.prototype.length;
`).
		AddSource(`/** @type{Interface1} */
var obj1;
var obj2 = new C1();
obj1 = obj2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching3()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I1() {}

/** @record */
function I2() {}
`).
		AddSource(`/** @type {I1} */
var i1;
/** @type {I2} */
var i2;
i1 = i2;
i2 = i1;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching4_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching4_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I1() {}

/** @record */
function I2() {}
`).
		AddSource(`/** @type {I1} */
var i1;
/** @type {I2} */
var i2;
i2 = i1;
i1 = i2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching5_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching5_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I1() {}

/** @interface */
function I3() {}
/** @type {number} */
I3.prototype.length;
`).
		AddSource(`/** @type {I1} */
var i1;
/** @type {I3} */
var i3;
i1 = i3;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching7_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching7_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I1() {}

/** @constructor */
function C1() {}
/** @type {number} */
C1.prototype.length;
`).
		AddSource(`/** @type {I1} */
var i1;/** @type {C1} */
var c1;
i1 = c1;   // no warning
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching9()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C1() {}
/** @type {number} */
C1.prototype.length;

/** @constructor */
function C2() {}
/** @type {number} */
C2.prototype.length;
`).
		AddSource(`/** @type {C1} */
var c1;/** @type {C2} */
var c2;
c1 = c2;
`).
		AddDiagnostic(`assignment
found   : (C2|null)
required: (C1|null)
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching11_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching11_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */
function I3() {}
/** @type {number} */
I3.prototype.length;

/**
 * @record
 * @extends I3
 */
function I4() {}
/** @type {boolean} */
I4.prototype.prop;

/** @constructor */
function C4() {}
/** @type {number} */
C4.prototype.length;
/** @type {boolean} */
C4.prototype.prop;
`).
		AddSource(`/** @type {I4} */
var i4;/** @type {C4} */
var c4;
i4 = c4;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching13()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
   * @record
   */
  function I5() {}
  /** @type {I5} */
  I5.prototype.next;

  /**
   * @interface
   */
  function C5() {}
  /** @type {C5} */
  C5.prototype.next;
`).
		AddSource(`/** @type {I5} */
var i5;/** @type {C5} */
var c5;
i5 = c5;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching13_2()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching13_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
 * @record
 */
function I5() {}
/** @type {I5} */
I5.prototype.next;

/**
 * @record
 */
function C5() {}
/** @type {C5} */
C5.prototype.next;
`).
		AddSource(`/** @type {I5} */
var i5;/** @type {C5} */
var c5;
i5 = c5;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching13_3()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching13_3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
 * @interface
 */
function I5() {}
/** @type {I5} */
I5.prototype.next;

/**
 * @record
 */
function C5() {}
/** @type {C5} */
C5.prototype.next;
`).
		AddSource(`/** @type {I5} */
var i5;/** @type {C5} */
var c5;
i5 = c5;
`).
		AddDiagnostic(`assignment
found   : (C5|null)
required: (I5|null)
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching15()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I5() {}
/** @type {I5} */
I5.prototype.next;

/** @constructor */
function C6() {}
/** @type {C6} */
C6.prototype.next;

/** @constructor */
function C5() {}
/** @type {C6} */
C5.prototype.next;
`).
		AddSource(`/** @type {I5} */
var i5;/** @type {C5} */
var c5;
i5 = c5;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching16_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching16_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD).
		AddSource(`/** @type {I5} */
var i5;/** @type {C5} */
var c5;
i5 = c5;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching17_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching17_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD).
		AddSource(`/** @type {C5} */
var c5;
/**
 * @param {I5} i5
 */
function f(i5) {}

f(c5);
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching18_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching18_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD).
		AddSource(`/** @type {I5} */
var i5;/** @type {C5} */
var c5;
i5.next = c5;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching19()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching19() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_NONMATCHING_CHAIN).
		AddSource(`/** @type {I5} */
var i5;
/** @type {C5} */
var c5;
i5 = c5;
`).
		AddDiagnostic(`assignment
found   : (C5|null)
required: (I5|null)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching20()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching20() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_NONMATCHING_CHAIN).
		AddSource(`/** @type {C5} */
var c5;
/**
 * @param {I5} i5
 */
function f(i5) {}

f(c5);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (C5|null)
required: (I5|null)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching21()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching21() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_NONMATCHING_CHAIN).
		AddSource(`/** @type {I5} */
var i5;
/** @type {C5} */
var c5;
i5.next = c5;
`).
		AddDiagnostic(`assignment to property next of I5
found   : (C5|null)
required: (I5|null)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching22_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching22_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching23()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching23() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_NONMATCHING_CHAIN + `/** @record */
function I7() {}
/** @type{function(): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		AddDiagnostic(`assignment
found   : (C7|null)
required: (I7|null)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching24_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching24_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(C5): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(I5): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching26_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching26_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(C5, C5, I5): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(I5, C5): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching29_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching29_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching30_1_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching30_1_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(this:I5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:C5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching30_2_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching30_2_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching30_3_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching30_3_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */ function I5() {}
/** @constructor @implements {I5} */ function C5() {}
/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching30_3_2()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching30_3_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */ function I5() {}
/** @constructor @implements {I5} */ function C5() {}
/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching30_3_3()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching30_3_3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */ function I5() {}
/** @constructor */ function C5() {}
/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		AddDiagnostic(`assignment
found   : (C7|null)
required: (I7|null)
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching30_3_4()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching30_3_4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */ function I5() {}
/** @constructor */ function C5() {}
/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching30_4_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching30_4_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */ function I5() {}
/** @constructor @implements {I5} */ function C5() {}
/** @record */
function I7() {}
/** @type{function(this:I5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:C5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching30_4_2()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching30_4_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */ function I5() {}
/** @constructor */ function C5() {}
/** @record */
function I7() {}
/** @type{function(this:I5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:C5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		AddDiagnostic(`assignment
found   : (C7|null)
required: (I7|null)
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching31_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching31_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching32_2()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching32_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {{prop: I7, prop2: C7}}*/
var r1;
/** @type {{prop: C7, prop2: C7}} */
var r2;
r1 = r2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching33_3()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching33_3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {{prop: I7, prop2: C7}}*/
var r1;
/** @type {{prop: C7, prop2: C7, prop3: C7}} */
var r2;
r1 = r2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching36_2()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching36_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {{fun: function(C7):I7, prop: {prop: I7}}} */
 var com1;
/** @type {{fun: function(I7):C7, prop: {prop: C7}}} */
var com2;

com1 = com2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching36_3()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching36_3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {{fun: function(C7):I7, prop: {prop: I7}}} */
 var com1;
/** @type {{fun: function(I7):C7, prop: {prop: C7}}} */
var com2;

com1 = com2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching37()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching37() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_NONMATCHING_CHAIN + `/** @record */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {{fun: function(C7):I7, prop: {prop: I7}}} */
var com1;
/** @type {{fun: function(I7):C7, prop: {prop: C7}}} */
var com2;

com1 = com2;
`).
		AddDiagnostic(`assignment
found   : {
  fun: function((I7|null)): (C7|null),
  prop: {prop: (C7|null)}
}
required: {
  fun: function((C7|null)): (I7|null),
  prop: {prop: (I7|null)}
}
missing : []
mismatch: [fun,prop]
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching_forObjectLiterals_39()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching_forObjectLiterals_39() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {}
/** @type {number} */
I2.prototype.length;
`).
		AddSource(`/** @type {I2} */
var o1 = {length : 'test'};
`).
		AddDiagnostic(`initializing variable
found   : {length: string}
required: (I2|null)
missing : []
mismatch: [length]
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching_forObjectLiterals_prototypeProp_matching()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching_forObjectLiterals_prototypeProp_matching() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {}
/** @type {number} */
I2.prototype.length;
`).
		AddSource(`/** @type {I2} */
var o1 = {length : 123};
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching_forObjectLiterals_instanceProp_matching()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching_forObjectLiterals_instanceProp_matching() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {
  /** @type {number} */
  this.length;
}
`).
		AddSource(`/** @type {!I2} */
var o1 = {length : 123};
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching_forObjectLiterals_prototypeProp_missing()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching_forObjectLiterals_prototypeProp_missing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {}
/** @type {number} */
I2.prototype.length;
`).
		AddSource(`/** @type {!I2} */
var o1 = {};
`).
		AddDiagnostic(`initializing variable
found   : {}
required: I2
missing : [length]
mismatch: []
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching_forObjectLiterals_instanceProp_missing()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching_forObjectLiterals_instanceProp_missing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {
  /** @type {number} */
  this.length;
}
`).
		AddSource(`/** @type {!I2} */
var o1 = {};
`).
		AddDiagnostic(`initializing variable
found   : {}
required: I2
missing : [length]
mismatch: []
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching_forObjectLiterals_prototypeProp_mismatch()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching_forObjectLiterals_prototypeProp_mismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {}
/** @type {number} */
I2.prototype.length;
`).
		AddSource(`/** @type {!I2} */
var o1 = {length: null};
`).
		AddDiagnostic(`initializing variable
found   : {length: null}
required: I2
missing : []
mismatch: [length]
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching_forObjectLiterals_instanceProp_mismatch()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching_forObjectLiterals_instanceProp_mismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {
  /** @type {number} */
  this.length;
}
`).
		AddSource(`/** @type {!I2} */
var o1 = {length: null};
`).
		AddDiagnostic(`initializing variable
found   : {length: null}
required: I2
missing : []
mismatch: [length]
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching_forObjectLiterals_41()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching_forObjectLiterals_41() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {}
/** @type {number} */
I2.prototype.length;
`).
		AddSource(`/** @type {I2} */
var o1 = {length : 123};
/** @type {I2} */
var i;
i = o1;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching41_forObjectLiterals_41_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching41_forObjectLiterals_41_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {}
/** @type {number} */
I2.prototype.length;
`).
		AddSource(`/** @type {I2} */
var o1 = {length : 123};
/** @type {I2} */
var i;
i = o1;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching_forObjectLiterals_42()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching_forObjectLiterals_42() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {}
/** @type {number} */
I2.prototype.length;
`).
		AddSource(`/** @type {{length: number}} */
var o1 = {length : 123};
/** @type {I2} */
var i;
i = o1;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching_forObjectLiterals_43()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching_forObjectLiterals_43() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {}
/** @type {number} */
I2.prototype.length;
`).
		AddSource(`var o1 = {length : 123};
/** @type {I2} */
var i;
i = o1;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching44()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching44() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */ function I() {}
/** @type {!Function} */ I.prototype.removeEventListener;
/** @type {!Function} */ I.prototype.addEventListener;
/** @constructor */ function C() {}
/** @type {!Function} */ C.prototype.addEventListener;
`).
		AddSource(`/** @param {C|I} x */
function f(x) { x.addEventListener(); }
f(new C());
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching45()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching45() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @record
 * @template X
 */
function I() {}
/** @constructor */
function C() {}
var /** !I */ i = new C;
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching46()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching46() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
function I2() {}
/**
 * @interface
 * @extends {I2}
 */
function I3() {}
/**
 * @record
 * @extends {I3}
 */
function I4() {}
/** @type {I4} */
var i4;
/** @type {I2} */
var i2;
i4 = i2;
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching47()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching47() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */
function I2() {}
/**
 * @interface
 * @extends {I2}
 */
function I3() {}
/**
 * @record
 * @extends {I3}
 */
function I4() {}
`).
		AddSource(`/** @type {I4} */
var i4;
/** @type {I2} */
var i2;
i4 = i2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching48()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching48() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("").
		AddSource(`/** @interface */
function I2() {}
/**
 * @record
 * @extends {I2}
 */
function I3() {}
/** @type {I3} */
var i3;
/** @type {I2} */
var i2;
i3 = i2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching49()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching49() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */
function I2() {}
/**
 * @record
 * @extends {I2}
 */
function I3() {}
`).
		AddSource(`/** @type {I3} */
var i3;
/** @type {I2} */
var i2;
i3 = i2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching49_2()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching49_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @record */
function I2() {}
/**
 * @record
 * @extends {I2}
 */
function I3() {}
`).
		AddSource(`/** @type {I3} */
var i3;
/** @type {I2} */
var i2;
i3 = i2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching50()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching50() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */
function I2() {}
/**
 * @record
 * @extends {I2}
 */
function I3() {}
`).
		AddSource(`/** @type {I3} */
var i3;
/** @type {{length : number}} */
var r = {length: 123};
i3 = r;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching1_1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching1_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */
function Interface1() {}
/** @type {number} */
Interface1.prototype.length;

/** @constructor */
function C1() {}
/** @type {number} */
C1.prototype.length;
`).
		AddSource(`/** @type{Interface1} */
var obj1;
/** @type{C1} */
var obj2 = new C1();
obj1 = obj2;
`).
		AddDiagnostic(`assignment
found   : (C1|null)
required: (Interface1|null)
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching22_2()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching22_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckStructuralInterfaceTest_EXTERNS_FOR_LONG_MATCHING_CHAIN_RECORD + `/** @interface */
function I7() {}
/** @type{function(): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		AddDiagnostic(`assignment
found   : (C7|null)
required: (I7|null)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceMatching30_3()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceMatching30_3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */ function I5() {}
/** @constructor @implements {I5} */ function C5() {}
/** @interface */
function I7() {}
/** @type{function(this:C5, C5, C5, I5=): I5} */
I7.prototype.getElement = function(){};

/** @constructor */
function C7() {}
/** @type{function(this:I5, I5, C5=, I5=): C5} */
C7.prototype.getElement = function(){};
`).
		AddSource(`/** @type {I7} */
var i7;
/** @type {C7} */
var c7;

i7 = c7;
`).
		AddDiagnostic(`assignment
found   : (C7|null)
required: (I7|null)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceWithOptionalProperty()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceWithOptionalProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function Rec() {}
/** @type {string} */ Rec.prototype.str;
/** @type {(number|undefined)} */ Rec.prototype.opt_num;

/** @constructor */ function Foo() {}
Foo.prototype.str = 'foo';

var /** !Rec */ x = new Foo;
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceWithUnknownProperty()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceWithUnknownProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function Rec() {}
/** @type {string} */ Rec.prototype.str;
/** @type {?} */ Rec.prototype.unknown;

/** @constructor */ function Foo() {}
Foo.prototype.str = 'foo';

var /** !Rec */ x = new Foo;
`).
		AddDiagnostic(`initializing variable
found   : Foo
required: Rec
missing : [unknown]
mismatch: []
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceWithOptionalUnknownProperty()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceWithOptionalUnknownProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function Rec() {}
/** @type {string} */ Rec.prototype.str;
/** @type {?|undefined} */ Rec.prototype.opt_unknown;

/** @constructor */ function Foo() {}
Foo.prototype.str = 'foo';

var /** !Rec */ x = new Foo;
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceWithTopProperty()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceWithTopProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function Rec() {}
/** @type {string} */ Rec.prototype.str;
/** @type {*} */ Rec.prototype.top;

/** @constructor */ function Foo() {}
Foo.prototype.str = 'foo';

var /** !Rec */ x = new Foo;
`).
		AddDiagnostic(`initializing variable
found   : Foo
required: Rec
missing : [top]
mismatch: []
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfaceCycleDoesntCrash()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfaceCycleDoesntCrash() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**  @record */ function Foo() {};
/**  @return {MutableFoo} */ Foo.prototype.toMutable;
/**  @record */ function MutableFoo() {};
/**  @param {Foo} from */ MutableFoo.prototype.copyFrom;

/**  @record */ function Bar() {};
/**  @return {MutableBar} */ Bar.prototype.toMutable;
/**  @record */ function MutableBar() {};
/**  @param {Bar} from */ MutableBar.prototype.copyFrom;

/** @constructor @implements {MutableBar} */ function MutableBarImpl() {};
/** @override */ MutableBarImpl.prototype.copyFrom = function(from) {};
/** @constructor  @implements {MutableFoo} */ function MutableFooImpl() {};
/** @override */ MutableFooImpl.prototype.copyFrom = function(from) {};
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfacesMatchOwnProperties1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfacesMatchOwnProperties1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function WithProp() {}
/** @type {number} */ WithProp.prototype.prop;

/** @constructor */
function Foo() {
  /** @type {number} */ this.prop = 5;
}
var /** !WithProp */ wp = new Foo;
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfacesMatchOwnProperties2()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfacesMatchOwnProperties2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function WithProp() {}
/** @type {number} */ WithProp.prototype.prop;

/** @constructor */
function Foo() {
  /** @type {number} */ this.oops = 5;
}
var /** !WithProp */ wp = new Foo;
`).
		AddDiagnostic(`initializing variable
found   : Foo
required: WithProp
missing : [prop]
mismatch: []
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfacesMatchOwnProperties3()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfacesMatchOwnProperties3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function WithProp() {}
/** @type {number} */ WithProp.prototype.prop;

/** @constructor */
function Foo() {
  /** @type {string} */ this.prop = 'str';
}
var /** !WithProp */ wp = new Foo;
`).
		AddDiagnostic(`initializing variable
found   : Foo
required: WithProp
missing : []
mismatch: [prop]
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfacesMatchFunctionNamespace1()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfacesMatchFunctionNamespace1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function WithProp() {}
/** @type {number} */ WithProp.prototype.prop;

var ns = function() {};
/** @type {number} */ ns.prop;
var /** !WithProp */ wp = ns;
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfacesMatchFunctionNamespace2()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfacesMatchFunctionNamespace2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function WithProp() {}
/** @type {number} */ WithProp.prototype.prop;

var ns = function() {};
/** @type {number} */ ns.oops;
var /** !WithProp */ wp = ns;
`).
		AddDiagnostic(`initializing variable
found   : function(): undefined
required: WithProp
missing : [prop]
mismatch: []
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testStructuralInterfacesMatchFunctionNamespace3()
func (c *TypeCheckStructuralInterfaceTest) TestStructuralInterfacesMatchFunctionNamespace3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function WithProp() {}
/** @type {number} */ WithProp.prototype.prop;

var ns = function() {};
/** @type {string} */ ns.prop;
var /** !WithProp */ wp = ns;
`).
		AddDiagnostic(`initializing variable
found   : function(): undefined
required: WithProp
missing : []
mismatch: [prop]
`).
		Run()
}

// Java: TypeCheckStructuralInterfaceTest.testRecursiveTemplatizedStructuralInterface()
func (c *TypeCheckStructuralInterfaceTest) TestRecursiveTemplatizedStructuralInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @record
 * @template T
 */
var Rec = function() { };
/** @type {!Rec<T>} */
Rec.prototype.p;

/**
 * @constructor @implements {Rec<U>}
 * @template U
 */
var Foo = function() {};
/** @override */
Foo.prototype.p = new Foo;
`).
		Run()
}

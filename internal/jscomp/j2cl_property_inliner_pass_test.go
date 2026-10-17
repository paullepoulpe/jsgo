// Port of com/google/javascript/jscomp/J2clPropertyInlinerPassTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/J2clPropertyInlinerPassTest.java)

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

// J2clPropertyInlinerPassTest tests J2clPropertyInlinerPass.
type J2clPropertyInlinerPassTest struct {
	CompilerTestCaseBase
}

// Java: J2clPropertyInlinerPassTest()
func NewJ2clPropertyInlinerPassTest() *J2clPropertyInlinerPassTest {
	c := &J2clPropertyInlinerPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestJ2clPropertyInlinerPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewJ2clPropertyInlinerPassTest() })
}

// Java: J2clPropertyInlinerPassTest.setUp()
func (c *J2clPropertyInlinerPassTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize() // Inlining will fail if normalization hasn't happened yet.
}

// Java: J2clPropertyInlinerPassTest.getProcessor(Compiler)
func (c *J2clPropertyInlinerPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewJ2clPropertyInlinerPass(compiler)
}

// Java: J2clPropertyInlinerPassTest.createCompiler()
func (c *J2clPropertyInlinerPassTest) CreateCompiler() *Compiler {
	compiler := c.CompilerTestCaseBase.CreateCompiler()
	J2clSourceFileChecker_MarkToRunJ2clPasses(compiler)
	return compiler
}

// Java: J2clPropertyInlinerPassTest.testDoesntChange(Sources)
func (c *J2clPropertyInlinerPassTest) TestDoesntChange(js CompilerTestCase_Sources) {
	c.TestSame(js)
}

// Java: J2clPropertyInlinerPassTest.testInlineCollapsedProp()
func (c *J2clPropertyInlinerPassTest) TestInlineCollapsedProp() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
var A$$0clinit = function() {
  A$$0x = 2;
};
Object.defineProperties(A, {x :{
  configurable:true,
  enumerable:true,
  get:function() {
    return A$$0clinit(), A$$0x;
  }
}});
var A$$0x = null;
var x = A.x;
`)),
		c.Expected(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
var A$$0clinit = function() {
  A$$0x = 2;
};
var A$$0x = null;
var x = (A$$0clinit(), A$$0x);
`)))
}

// Java: J2clPropertyInlinerPassTest.testNoInlineNonJ2clProps()
func (c *J2clPropertyInlinerPassTest) TestNoInlineNonJ2clProps() {
	c.TestDoesntChange(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
var globalValue = null;
var initGlobalValue = function() {
  globalValue = 2;
};
Object.defineProperties(A, {x :{
  configurable:true,
  enumerable:true,
  get:function() {
    return initGlobalValue, globalValue;
  }
}});
var x = A.x;
`)))
}

// Java: J2clPropertyInlinerPassTest.testNoInlineNonJ2clPropsValue()
func (c *J2clPropertyInlinerPassTest) TestNoInlineNonJ2clPropsValue() {
	c.TestDoesntChange(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
var A$$0clinit = function() {
  A$$0x = 2;
};
Object.defineProperties(A, {x :{
  configurable:true,
  enumerable:true,
  value: 2
}});
var A$$0x = null;
var x = A.x;
`)))
}

// Java: J2clPropertyInlinerPassTest.testNoStripDefineProperties()
func (c *J2clPropertyInlinerPassTest) TestNoStripDefineProperties() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$j2cl_prop = 2;
};
Object.defineProperties(A, {
  j2cl_prop: {
    configurable: true,
    enumerable: true,
    get: function() {
      return A.$clinit(), A.$j2cl_prop;
    },
    set: function(value) {
      A.$clinit(), A.$j2cl_prop = value;
    }
  },
  non_j2cl_prop: {
    configurable: true,
    enumerable: true,
    get: function() {
      return 55;
    },
    set: function(v) {
      console.log(v);
    }
  },
});
`)),
		c.Expected(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$j2cl_prop = 2;
};
Object.defineProperties(A, {
  non_j2cl_prop: {
    configurable: true,
    enumerable: true,
    get: function() {
      return 55;
    },
    set: function(v) {
      console.log(v);
    }
  },
});
`)))
}

// Java: J2clPropertyInlinerPassTest.testInlineDefinePropertiesGetter()
func (c *J2clPropertyInlinerPassTest) TestInlineDefinePropertiesGetter() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
Object.defineProperties(A, {x: {
  configurable:true,
  enumerable:true,
  get:function() {
    return A.$clinit(), A.$x;
  },
  set: function(value) {
    A.$clinit(), A.$x = value;
  }
}});
A.$x = 3;
var xx = A.x;
`)),
		c.Expected(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
A.$x = 3;
var xx = (A.$clinit(), A.$x);
`)))
}

// Java: J2clPropertyInlinerPassTest.testInlineDefinePropertiesSetter()
func (c *J2clPropertyInlinerPassTest) TestInlineDefinePropertiesSetter() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
Object.defineProperties(A, {x: {
  configurable:true,
  enumerable:true,
  get:function() {
    return A.$clinit(), A.$x;
  },
  set: function(value) {
    A.$clinit(), A.$x = value;
  }
}});
A.$x = 3;
A.x = 10;
`)),
		c.Expected(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
A.$x = 3;
{(A.$clinit(), A.$x = 10);}
`)))
}

// Java: J2clPropertyInlinerPassTest.testInlineGettersInQualifier()
func (c *J2clPropertyInlinerPassTest) TestInlineGettersInQualifier() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = {y: 2};
};
Object.defineProperties(A, {x: {
  configurable:true,
  enumerable:true,
  get:function() {
    return A.$clinit(), A.$x;
  },
  set: function(value) {
    A.$clinit(), A.$x = value;
  }
}});
A.$x = null;
var xy = A.x.y;
`)),
		c.Expected(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = {y: 2};
};
A.$x = null;
var xy = (A.$clinit(), A.$x).y;
`)))
}

// Java: J2clPropertyInlinerPassTest.testInlineGettersInQualifierNoSetter()
func (c *J2clPropertyInlinerPassTest) TestInlineGettersInQualifierNoSetter() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = {y: 2};
};
Object.defineProperties(A, {x: {
  configurable:true,
  enumerable:true,
  get:function() {
    return A.$clinit(), A.$x;
  },
}});
A.$x = null;
var xy = A.x.y;
`)),
		c.Expected(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = {y: 2};
};
A.$x = null;
var xy = (A.$clinit(), A.$x).y;
`)))
}

// Java: J2clPropertyInlinerPassTest.testNoInlineCompoundAssignment()
func (c *J2clPropertyInlinerPassTest) TestNoInlineCompoundAssignment() {
	c.TestDoesntChange(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
Object.defineProperties(A, {x: {
  configurable:true,
  enumerable:true,
  get:function() {
    return A.$clinit(), A.$x;
  },
  set: function(value) {
    A.$clinit(), A.$x = value;
  }
}});
A.$x = 3;
A.x += 5;
`)))
}

// Java: J2clPropertyInlinerPassTest.testNoInlineIncrementGetter()
func (c *J2clPropertyInlinerPassTest) TestNoInlineIncrementGetter() {
	// Test ++
	c.TestDoesntChange(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
Object.defineProperties(A, {x: {
  configurable:true,
  enumerable:true,
  get:function() {
    return A.$clinit(), A.$x;
  },
  set: function(value) {
    A.$clinit(), A.$x = value;
  }
}});
A.$x = 3;
A.x--;
`)))

	// Test --
	c.TestDoesntChange(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
Object.defineProperties(A, {x: {
  configurable:true,
  enumerable:true,
  get:function() {
    return A.$clinit(), A.$x;
  },
  set: function(value) {
    A.$clinit(), A.$x = value;
  }
}});
A.$x = 3;
A.x++;
`)))
}

// Java: J2clPropertyInlinerPassTest.testNoInlineSetterOnlyGetter()
func (c *J2clPropertyInlinerPassTest) TestNoInlineSetterOnlyGetter() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = {y: 2};
};
Object.defineProperties(A, {x: {
  configurable:true,
  enumerable:true,
  get:function() {
    return A.$clinit(), A.$x;
  },
}});
A.$x = null;
A.x = null;
`)),
		c.Expected(
			SourceFile_FromCode(
				"someFile.js",
				`var A = function() {};
A.$clinit = function() {
  A.$x = {y: 2};
};
A.$x = null;
A.x = null;
`)))
}

// Java: J2clPropertyInlinerPassTest.testInlineEs6Getter()
func (c *J2clPropertyInlinerPassTest) TestInlineEs6Getter() {
	c.Test(
		`class A {
  static $clinit() {
    A.$clinit = function() {};
    A.$x = 2;
  }
  static get x() {
    return A.$clinit(), A.$x
  }
  static set x(value) {
    A.$clinit(), A.$x = value;
  }
}
A.$x = 3;
var xx = A.x
`,
		`class A {
  static $clinit() {
    A.$clinit = function() {};
    A.$x = 2;
  }
}
A.$x = 3;
var xx = (A.$clinit(), A.$x);
`)
}

// Java: J2clPropertyInlinerPassTest.testInlineEs6GetterNoSetter()
func (c *J2clPropertyInlinerPassTest) TestInlineEs6GetterNoSetter() {
	c.Test(
		`class A {
  static $clinit() {
    A.$clinit = function() {};
    A.$x = 2;
  }
  static get x() {
    return A.$clinit(), A.$x
  }
}
A.$x = 3;
var xx = A.x
`,
		`class A {
  static $clinit() {
    A.$clinit = function() {};
    A.$x = 2;
  }
}
A.$x = 3;
var xx = (A.$clinit(), A.$x);
`)
}

// Java: J2clPropertyInlinerPassTest.testInlineEs6Setter()
func (c *J2clPropertyInlinerPassTest) TestInlineEs6Setter() {
	c.Test(
		`class A {
  static $clinit() {
    A.$clinit = function() {};
    A.$x = 2;
  }
  static get x() {
    return A.$clinit(), A.$x
  }
  static set x(value) {
    A.$clinit(), A.$x = value;
  }
}
A.$x = 3;
A.x = 5;
`,
		`class A {
  static $clinit() {
    A.$clinit = function() {};
    A.$x = 2;
  }
}
A.$x = 3;
{(A.$clinit(), A.$x = 5)}
`)
}

// Java: J2clPropertyInlinerPassTest.testInlineEs6GetterSetter_multiple()
func (c *J2clPropertyInlinerPassTest) TestInlineEs6GetterSetter_multiple() {
	c.Test(
		`class A {
  static $clinit() {
    A.$clinit = function() {};
    A.$x = 2;
    A.$y = 3;
  }
  static get x() {
    return A.$clinit(), A.$x
  }
  static set y(value) {
    A.$clinit(), A.$y = value;
  }
  static get y() {
    return A.$clinit(), A.$y
  }
  static set x(value) {
    A.$clinit(), A.$x = value;
  }
}
var xx = A.x;
var yy = A.y
A.x = 5;
A.y = 5;
`,
		`class A {
  static $clinit() {
    A.$clinit = function() {};
    A.$x = 2;
    A.$y = 3;
  }
}
var xx = (A.$clinit(), A.$x);
var yy = (A.$clinit(), A.$y)
{(A.$clinit(), A.$x = 5)}
{(A.$clinit(), A.$y = 5)}
`)
}

// Java: J2clPropertyInlinerPassTest.testOptChainReadLeavesAccessorRemoved()
func (c *J2clPropertyInlinerPassTest) TestOptChainReadLeavesAccessorRemoved() {
	// TODO(b/571560495): Object.defineProperties accessor removed while A?.x read remains uninlined
	c.Test(
		`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
Object.defineProperties(A, {x: {
  configurable:true,
  enumerable:true,
  get:function() {
    return A.$clinit(), A.$x;
  },
  set: function(value) {
    A.$clinit(), A.$x = value;
  }
}});
A.$x = 3;
var r = A?.x;
`,
		`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
A.$x = 3;
var r = A?.x;
`)
}

// Java: J2clPropertyInlinerPassTest.testOptChainReadLeavesEs6AccessorRemoved()
func (c *J2clPropertyInlinerPassTest) TestOptChainReadLeavesEs6AccessorRemoved() {
	// TODO(b/571560495): ES6 class static getter/setter removed while A?.x read remains uninlined
	c.Test(
		`class A {
  static $clinit() {
    A.$clinit = function() {};
    A.$x = 2;
  }
  static get x() {
    return A.$clinit(), A.$x
  }
  static set x(value) {
    A.$clinit(), A.$x = value;
  }
}
A.$x = 3;
var xx = A?.x;
`,
		`class A {
  static $clinit() {
    A.$clinit = function() {};
    A.$x = 2;
  }
}
A.$x = 3;
var xx = A?.x;
`)
}

// Java: J2clPropertyInlinerPassTest.testOptChainReadMixedWithPlainRead()
func (c *J2clPropertyInlinerPassTest) TestOptChainReadMixedWithPlainRead() {
	// TODO(b/571560495): Object.defineProperties accessor removed while A?.x read remains uninlined
	c.Test(
		`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
Object.defineProperties(A, {x: {
  configurable:true,
  enumerable:true,
  get:function() {
    return A.$clinit(), A.$x;
  },
  set: function(value) {
    A.$clinit(), A.$x = value;
  }
}});
A.$x = 3;
var r = A?.x;
var s = A.x;
`,
		`var A = function() {};
A.$clinit = function() {
  A.$x = 2;
};
A.$x = 3;
var r = A?.x;
var s = (A.$clinit(), A.$x);
`)
}

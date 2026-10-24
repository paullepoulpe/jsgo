// Port of com/google/javascript/jscomp/ChromePassTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ChromePassTest.java)

/*
 * Copyright 2014 The Closure Compiler Authors.
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

// ChromePassTest tests ChromePass.
type ChromePassTest struct {
	CompilerTestCaseBase
}

// Java: ChromePassTest()
func NewChromePassTest() *ChromePassTest {
	c := &ChromePassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestChromePassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewChromePassTest() })
}

// Java: ChromePassTest.getProcessor(Compiler)
func (c *ChromePassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewChromePass(compiler)
}

// Java: ChromePassTest.testCrDefineCreatesObjectsForQualifiedName()
func (c *ChromePassTest) TestCrDefineCreatesObjectsForQualifiedName() {
	c.Test(
		"cr.define('my.namespace.name', function() {\n"+
			"  return {};\n"+
			"});\n",
		"var my = my || {};\n"+
			"my.namespace = my.namespace || {};\n"+
			"my.namespace.name = my.namespace.name || {};\n"+
			"cr.define('my.namespace.name', function() {\n"+
			"  return {};\n"+
			"});\n")
}

// Java: ChromePassTest.testChromePassIgnoresModules()
func (c *ChromePassTest) TestChromePassIgnoresModules() {
	c.TestSame("export var x;")
}

// Java: ChromePassTest.testCrDefineAssignsExportedFunctionByQualifiedName()
func (c *ChromePassTest) TestCrDefineAssignsExportedFunctionByQualifiedName() {
	c.Test(
		"cr.define('namespace', function() {\n"+
			"  function internalStaticMethod() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  return {\n"+
			"    externalStaticMethod: internalStaticMethod\n"+
			"  };\n"+
			"});\n",
		"var namespace = namespace || {};\n"+
			"cr.define('namespace', function() {\n"+
			"  namespace.externalStaticMethod = function internalStaticMethod() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  return {\n"+
			"    externalStaticMethod: namespace.externalStaticMethod\n"+
			"  };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineCopiesJSDocForExportedFunction()
func (c *ChromePassTest) TestCrDefineCopiesJSDocForExportedFunction() {
	c.Test(
		"cr.define('namespace', function() {\n"+
			"  /** I'm function's JSDoc */\n"+
			"  function internalStaticMethod() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  return {\n"+
			"    externalStaticMethod: internalStaticMethod\n"+
			"  };\n"+
			"});\n",
		"var namespace = namespace || {};\n"+
			"cr.define('namespace', function() {\n"+
			"  /** I'm function's JSDoc */\n"+
			"  namespace.externalStaticMethod = function internalStaticMethod() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  return {\n"+
			"    externalStaticMethod: namespace.externalStaticMethod\n"+
			"  };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineReassignsExportedVarByQualifiedName()
func (c *ChromePassTest) TestCrDefineReassignsExportedVarByQualifiedName() {
	c.Test(
		"cr.define('namespace', function() {\n"+
			"  var internalStaticMethod = function() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  return {\n"+
			"    externalStaticMethod: internalStaticMethod\n"+
			"  };\n"+
			"});\n",
		"var namespace = namespace || {};\n"+
			"cr.define('namespace', function() {\n"+
			"  namespace.externalStaticMethod = function() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  return {\n"+
			"    externalStaticMethod: namespace.externalStaticMethod\n"+
			"  };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineExportsVarsWithoutAssignment()
func (c *ChromePassTest) TestCrDefineExportsVarsWithoutAssignment() {
	c.Test(
		"cr.define('namespace', function() {\n"+
			"  var a;\n"+
			"  return {\n"+
			"    a: a\n"+
			"  };\n"+
			"});\n",
		"var namespace = namespace || {};\n"+
			"cr.define('namespace', function() {\n"+
			"  namespace.a;\n"+
			"  return {\n"+
			"    a: namespace.a\n"+
			"  };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineExportsVarsWithoutAssignmentWithJSDoc()
func (c *ChromePassTest) TestCrDefineExportsVarsWithoutAssignmentWithJSDoc() {
	c.Test(
		"cr.define('namespace', function() {\n"+
			"  /** @type {number} */\n"+
			"  var a;\n"+
			"  return {\n"+
			"    a: a\n"+
			"  };\n"+
			"});\n",
		"var namespace = namespace || {};\n"+
			"cr.define('namespace', function() {\n"+
			"  /** @type {number} */\n"+
			"  namespace.a;\n"+
			"  return {\n"+
			"    a: namespace.a\n"+
			"  };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineCopiesJSDocForExportedVariable()
func (c *ChromePassTest) TestCrDefineCopiesJSDocForExportedVariable() {
	c.Test(
		"cr.define('namespace', function() {\n"+
			"  /** I'm function's JSDoc */\n"+
			"  var internalStaticMethod = function() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  return {\n"+
			"    externalStaticMethod: internalStaticMethod\n"+
			"  };\n"+
			"});\n",
		"var namespace = namespace || {};\n"+
			"cr.define('namespace', function() {\n"+
			"  /** I'm function's JSDoc */\n"+
			"  namespace.externalStaticMethod = function() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  return {\n"+
			"    externalStaticMethod: namespace.externalStaticMethod\n"+
			"  };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineDoesNothingWithNonExportedFunction()
func (c *ChromePassTest) TestCrDefineDoesNothingWithNonExportedFunction() {
	c.Test(
		"cr.define('namespace', function() {\n"+
			"  function internalStaticMethod() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  return {};\n"+
			"});\n",
		"var namespace = namespace || {};\n"+
			"cr.define('namespace', function() {\n"+
			"  function internalStaticMethod() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  return {};\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineDoesNothingWithNonExportedVar()
func (c *ChromePassTest) TestCrDefineDoesNothingWithNonExportedVar() {
	c.Test(
		"cr.define('namespace', function() {\n"+
			"  var a;\n"+
			"  var b;\n"+
			"  return {\n"+
			"    a: a\n"+
			"  };\n"+
			"});\n",
		"var namespace = namespace || {};\n"+
			"cr.define('namespace', function() {\n"+
			"  namespace.a;\n"+
			"  var b;\n"+
			"  return {\n"+
			"    a: namespace.a\n"+
			"  };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineDoesNothingWithExportedNotAName()
func (c *ChromePassTest) TestCrDefineDoesNothingWithExportedNotAName() {
	c.Test(
		"cr.define('namespace', function() {\n"+
			"  return {\n"+
			"    a: 42\n"+
			"  };\n"+
			"});\n",
		"var namespace = namespace || {};\n"+
			"cr.define('namespace', function() {\n"+
			"  return {\n"+
			"    a: 42\n"+
			"  };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineChangesReferenceToExportedFunction()
func (c *ChromePassTest) TestCrDefineChangesReferenceToExportedFunction() {
	c.Test(
		"cr.define('namespace', function() {\n"+
			"  function internalStaticMethod() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  function letsUseIt() {\n"+
			"    internalStaticMethod();\n"+
			"  }\n"+
			"  return {\n"+
			"    externalStaticMethod: internalStaticMethod\n"+
			"  };\n"+
			"});\n",
		"var namespace = namespace || {};\n"+
			"cr.define('namespace', function() {\n"+
			"  namespace.externalStaticMethod = function internalStaticMethod() {\n"+
			"    alert(42);\n"+
			"  }\n"+
			"  function letsUseIt() {\n"+
			"    namespace.externalStaticMethod();\n"+
			"  }\n"+
			"  return {\n"+
			"    externalStaticMethod: namespace.externalStaticMethod\n"+
			"  };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineConstEnum()
func (c *ChromePassTest) TestCrDefineConstEnum() {
	c.Test(
		"cr.define('foo', function() {\n"+
			"  /**\n"+
			"   * @enum {string}\n"+
			"   */\n"+
			"  const DangerType = {\n"+
			"    NOT_DANGEROUS: 'NOT_DANGEROUS',\n"+
			"    DANGEROUS: 'DANGEROUS',\n"+
			"  };\n"+
			"\n"+
			"  return {\n"+
			"    DangerType: DangerType,\n"+
			"  };\n"+
			"});\n",
		"var foo = foo || {};\n"+
			"cr.define('foo', function() {\n"+
			"  /** @enum {string} */\n"+
			"  foo.DangerType = {\n"+
			"    NOT_DANGEROUS:'NOT_DANGEROUS',\n"+
			"    DANGEROUS:'DANGEROUS',\n"+
			"  };\n"+
			"\n"+
			"  return {\n"+
			"    DangerType: foo.DangerType\n"+
			"  }\n"+
			"})\n")
}

// Java: ChromePassTest.testCrDefineLetEnum()
func (c *ChromePassTest) TestCrDefineLetEnum() {
	c.Test(
		"cr.define('foo', function() {\n"+
			"  /**\n"+
			"   * @enum {string}\n"+
			"   */\n"+
			"  let DangerType = {\n"+
			"    NOT_DANGEROUS: 'NOT_DANGEROUS',\n"+
			"    DANGEROUS: 'DANGEROUS',\n"+
			"  };\n"+
			"\n"+
			"  return {\n"+
			"    DangerType: DangerType,\n"+
			"  };\n"+
			"});\n",
		"var foo = foo || {};\n"+
			"cr.define('foo', function() {\n"+
			"  /** @enum {string} */\n"+
			"  foo.DangerType = {\n"+
			"    NOT_DANGEROUS:'NOT_DANGEROUS',\n"+
			"    DANGEROUS:'DANGEROUS',\n"+
			"  };\n"+
			"\n"+
			"  return {\n"+
			"    DangerType: foo.DangerType\n"+
			"  }\n"+
			"})\n")
}

// Java: ChromePassTest.testCrDefineWrongNumberOfArguments()
func (c *ChromePassTest) TestCrDefineWrongNumberOfArguments() {
	c.TestError(
		"cr.define('namespace', function() { return {}; }, 'invalid argument')\n",
		ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS)
}

// Java: ChromePassTest.testCrDefineInvalidFirstArgument()
func (c *ChromePassTest) TestCrDefineInvalidFirstArgument() {
	c.TestError(
		"cr.define(42, function() { return {}; })\n", ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT)
}

// Java: ChromePassTest.testCrDefineInvalidSecondArgument()
func (c *ChromePassTest) TestCrDefineInvalidSecondArgument() {
	c.TestError("cr.define('namespace', 42)\n", ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT)
}

// Java: ChromePassTest.testCrDefineInvalidReturnInFunction()
func (c *ChromePassTest) TestCrDefineInvalidReturnInFunction() {
	c.TestError(
		"cr.define('namespace', function() {})\n", ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION)
}

// Java: ChromePassTest.testObjectDefinePropertyWithoutArgsDoesntThrow()
func (c *ChromePassTest) TestObjectDefinePropertyWithoutArgsDoesntThrow() {
	c.Test("Object.defineProperty();", "Object.defineProperty();")
	c.Test("Object.defineProperty(a);", "Object.defineProperty(a);")
	c.Test("Object.defineProperty(a, 'b');", "Object.defineProperty(a, 'b');")
}

// Java: ChromePassTest.testObjectDefinePropertyDefinesUnquotedProperty()
func (c *ChromePassTest) TestObjectDefinePropertyDefinesUnquotedProperty() {
	c.Test(
		"Object.defineProperty(a.b, 'c', {});",
		"Object.defineProperty(a.b, 'c', {});\n"+
			"/** @type {?} */\n"+
			"a.b.c;\n")
}

// Java: ChromePassTest.testCrDefinePropertyTooFewArguments()
func (c *ChromePassTest) TestCrDefinePropertyTooFewArguments() {
	c.TestError("cr.defineProperty();\n", ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS)
	c.TestError("cr.defineProperty(a.prototype);\n", ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS)
}

// Java: ChromePassTest.testCrDefinePropertyDefinesUnquotedPropertyWithStringTypeForPropertyKindAttr()
func (c *ChromePassTest) TestCrDefinePropertyDefinesUnquotedPropertyWithStringTypeForPropertyKindAttr() {
	c.Test(
		"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.ATTR);",
		"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.ATTR);\n"+
			"/** @type {string} */\n"+
			"a.prototype.c;\n")
}

// Java: ChromePassTest.testCrDefinePropertyDefinesUnquotedPropertyWithBooleanTypeForPropertyKindBoolAttr()
func (c *ChromePassTest) TestCrDefinePropertyDefinesUnquotedPropertyWithBooleanTypeForPropertyKindBoolAttr() {
	c.Test(
		"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.BOOL_ATTR);",
		"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.BOOL_ATTR);\n"+
			"/** @type {boolean} */\n"+
			"a.prototype.c;\n")
}

// Java: ChromePassTest.testCrDefinePropertyDefinesUnquotedPropertyWithAnyTypeForPropertyKindJs()
func (c *ChromePassTest) TestCrDefinePropertyDefinesUnquotedPropertyWithAnyTypeForPropertyKindJs() {
	c.Test(
		"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.JS);",
		"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.JS);\n"+
			"/** @type {?} */\n"+
			"a.prototype.c;\n")
}

// Java: ChromePassTest.testCrDefinePropertyDefinesUnquotedPropertyWithTypeInfoForPropertyKindJs()
func (c *ChromePassTest) TestCrDefinePropertyDefinesUnquotedPropertyWithTypeInfoForPropertyKindJs() {
	c.Test(
		"/** @type {!Object} */\n"+
			"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.JS);\n",
		"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.JS);\n"+
			"// But gets moved here.\n"+
			"/** @type {!Object} */\n"+
			"a.prototype.c;\n")
}

// Java: ChromePassTest.testCrDefinePropertyDefinesUnquotedPropertyIgnoringJsDocWhenBoolAttrIsPresent()
func (c *ChromePassTest) TestCrDefinePropertyDefinesUnquotedPropertyIgnoringJsDocWhenBoolAttrIsPresent() {
	c.Test(
		"/** @type {!Object} */\n"+
			"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.BOOL_ATTR);\n",
		"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.BOOL_ATTR);\n"+
			"// @type is now {boolean}. Earlier, manually-specified @type is ignored.\n"+
			"/** @type {boolean} */\n"+
			"a.prototype.c;\n")
}

// Java: ChromePassTest.testCrDefinePropertyDefinesUnquotedPropertyIgnoringJsDocWhenAttrIsPresent()
func (c *ChromePassTest) TestCrDefinePropertyDefinesUnquotedPropertyIgnoringJsDocWhenAttrIsPresent() {
	c.Test(
		"/** @type {!Array} */\n"+
			"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.ATTR);\n",
		"cr.defineProperty(a.prototype, 'c', cr.PropertyKind.ATTR);\n"+
			"// @type is now {string}. Earlier, manually-specified @type is ignored.\n"+
			"/** @type {string} */\n"+
			"a.prototype.c;\n")
}

// Java: ChromePassTest.testCrDefinePropertyCalledWithouthThirdArgumentMeansCrPropertyKindJs()
func (c *ChromePassTest) TestCrDefinePropertyCalledWithouthThirdArgumentMeansCrPropertyKindJs() {
	c.Test(
		"cr.defineProperty(a.prototype, 'c');",
		"cr.defineProperty(a.prototype, 'c');\n"+
			"/** @type {?} */\n"+
			"a.prototype.c;\n")
}

// Java: ChromePassTest.testCrDefinePropertyDefinesUnquotedPropertyOnPrototypeWhenFunctionIsPassed()
func (c *ChromePassTest) TestCrDefinePropertyDefinesUnquotedPropertyOnPrototypeWhenFunctionIsPassed() {
	c.Test(
		"cr.defineProperty(a, 'c', cr.PropertyKind.JS);",
		"cr.defineProperty(a, 'c', cr.PropertyKind.JS);\n"+
			"/** @type {?} */\n"+
			"a.prototype.c;\n")
}

// Java: ChromePassTest.testCrDefinePropertyInvalidPropertyKind()
func (c *ChromePassTest) TestCrDefinePropertyInvalidPropertyKind() {
	c.TestError(
		"cr.defineProperty(a.b, 'c', cr.PropertyKind.INEXISTENT_KIND);",
		ChromePass_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND)
}

// Java: ChromePassTest.testCrDefineCreatesEveryObjectOnlyOnce()
func (c *ChromePassTest) TestCrDefineCreatesEveryObjectOnlyOnce() {
	c.Test(
		"cr.define('a.b.c.d', function() {\n"+
			"  return {};\n"+
			"});\n"+
			"cr.define('a.b.e.f', function() {\n"+
			"  return {};\n"+
			"});\n",
		"var a = a || {};\n"+
			"a.b = a.b || {};\n"+
			"a.b.c = a.b.c || {};\n"+
			"a.b.c.d = a.b.c.d || {};\n"+
			"cr.define('a.b.c.d', function() {\n"+
			"  return {};\n"+
			"});\n"+
			"a.b.e = a.b.e || {};\n"+
			"a.b.e.f = a.b.e.f || {};\n"+
			"cr.define('a.b.e.f', function() {\n"+
			"  return {};\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineDoesntRedefineCrVar()
func (c *ChromePassTest) TestCrDefineDoesntRedefineCrVar() {
	c.Test(
		"cr.define('cr.ui', function() {\n"+
			"  return {};\n"+
			"});\n",
		"cr.ui = cr.ui || {};\n"+
			"cr.define('cr.ui', function() {\n"+
			"  return {};\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineFunction()
func (c *ChromePassTest) TestCrDefineFunction() {
	c.Test(
		"cr.define('settings', function() {\n"+
			"  var x = 0;\n"+
			"  function C() {}\n"+
			"  return { C: C };\n"+
			"});\n",
		"var settings = settings || {};\n"+
			"cr.define('settings', function() {\n"+
			"  var x = 0;\n"+
			"  settings.C = function C() {};\n"+
			"  return { C: settings.C };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineClassStatement()
func (c *ChromePassTest) TestCrDefineClassStatement() {
	c.Test(
		"cr.define('settings', function() {\n"+
			"  var x = 0;\n"+
			"  class C {}\n"+
			"  return { C: C };\n"+
			"});\n",
		"var settings = settings || {};\n"+
			"cr.define('settings', function() {\n"+
			"  var x = 0;\n"+
			"  settings.C = class {}\n"+
			"  return { C: settings.C };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineClassExpression()
func (c *ChromePassTest) TestCrDefineClassExpression() {
	c.Test(
		"cr.define('settings', function() {\n"+
			"  var x = 0;\n"+
			"  var C = class {}\n"+
			"  return { C: C };\n"+
			"});\n",
		"var settings = settings || {};\n"+
			"cr.define('settings', function() {\n"+
			"  var x = 0;\n"+
			"  settings.C = class {}\n"+
			"  return { C: settings.C };\n"+
			"});\n")
}

// Java: ChromePassTest.testCrDefineClassWithInternalSelfReference()
func (c *ChromePassTest) TestCrDefineClassWithInternalSelfReference() {
	c.Test(
		"cr.define('settings', function() {\n"+
			"  var x = 0;\n"+
			"  class C {\n"+
			"    static create() { return new C; }\n"+
			"  }\n"+
			"  return { C: C };\n"+
			"});\n",
		"var settings = settings || {};\n"+
			"cr.define('settings', function() {\n"+
			"  var x = 0;\n"+
			"  settings.C = class {\n"+
			"    static create() { return new settings.C; }\n"+
			"  }\n"+
			"  return { C: settings.C };\n"+
			"});\n")
}

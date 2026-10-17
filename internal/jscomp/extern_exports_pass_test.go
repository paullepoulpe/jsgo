// Port of com/google/javascript/jscomp/ExternExportsPassTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/ExternExportsPassTest.java)

/*
 * Copyright 2009 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Tests for ExternExportsPass.
type ExternExportsPassTest struct {
	CompilerTestCaseBase
}

func NewExternExportsPassTest() *ExternExportsPassTest {
	c := &ExternExportsPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestExternExportsPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewExternExportsPassTest() })
}

// Java: ExternExportsPassTest.setUp()
func (c *ExternExportsPassTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
}

// Java: ExternExportsPassTest.getOptions()
func (c *ExternExportsPassTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetExternExportsPath(jutil.StringPtr("exports.js"))
	// Check types so we can make sure our exported externs have type
	// information.
	options.SetCheckSymbols(true)
	return options
}

// Java: ExternExportsPassTest.getProcessor(Compiler)
func (c *ExternExportsPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewExternExportsPass(compiler)
}

// Java: ExternExportsPassTest.compileAndCheck(String, String)
func (c *ExternExportsPassTest) compileAndCheck(js string, expected string) {
	c.compileAndCheckWithExterns(CompilerTestCase_MINIMAL_EXTERNS, js, expected)
}

// Java: ExternExportsPassTest.compileAndCheck(String, String, String)
func (c *ExternExportsPassTest) compileAndCheckWithExterns(externs string, js string, expected string) {
	c.compileAndExportExternsWithConsumer(
		js,
		externs,
		func(generatedExterns *string) {
			fileoverview :=
				`/**
 * @fileoverview Generated externs.
 * @externs
 */
`
			// NOTE(sdh): The type checker just produces {?}.
			// For now we will not worry about this distinction and just
			// normalize it.
			normalized := jutil.StringReplace(*generatedExterns, "?=", "?")

			AssertThat(normalized).IsEqualTo(fileoverview + expected)
		})
}

// compileAndExportExterns compiles the passed in JavaScript and returns the
// new externs exported by the this pass. js is the source to be compiled.
//
// Java: ExternExportsPassTest.compileAndExportExterns(String)
func (c *ExternExportsPassTest) compileAndExportExterns(js string) {
	c.compileAndExportExternsWithExterns(js, CompilerTestCase_MINIMAL_EXTERNS)
}

// compileAndExportExternsWithExterns compiles the passed in JavaScript with
// the passed in externs and returns the new externs exported by the this
// pass. js is the source to be compiled; externs the externs the js source
// needs.
//
// Java: ExternExportsPassTest.compileAndExportExterns(String, String)
func (c *ExternExportsPassTest) compileAndExportExternsWithExterns(js string, externs string) {
	c.compileAndExportExternsWithConsumer(js, externs, nil)
}

// compileAndExportExternsWithConsumer compiles the passed in JavaScript with
// the passed in externs and returns the new externs exported by the this
// pass. js is the source to be compiled; externs the externs the js source
// needs; consumer the consumer for the externs generated from js.
//
// Java: ExternExportsPassTest.compileAndExportExterns(String, String, Consumer<String>)
func (c *ExternExportsPassTest) compileAndExportExternsWithConsumer(
	js string, externs string, consumer func(*string)) {
	js =
		`/** @const */ var goog = {};
goog.exportSymbol = function(a, b) {};
goog.exportProperty = function(a, b, c) {};
` + js

	c.Test(
		c.Externs(externs),
		c.Srcs(js),
		CompilerTestCase_PostconditionFunc(
			func(compiler *Compiler) {
				if consumer != nil {
					consumer(compiler.GetResult().externExport)
				}
			}))
}

// Java: ExternExportsPassTest.testExportSymbol()
func (c *ExternExportsPassTest) TestExportSymbol() {
	c.compileAndCheck(
		`/** @const */ var a = {};
/** @const */ a.b = {};
a.b.c = function(d, e, f) {};
var x;
// Ensure we don't have a recurrence of a bug that caused us to consider the last
// assignment containing 'a.b.c' to be its definition, even if it was on the rhs.
x = a.b.c;
goog.exportSymbol('foobar', a.b.c)
`,
		`/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @return {undefined}
 */
var foobar = function(d, e, f) {
};
`)
}

// Java: ExternExportsPassTest.exportClassAsSymbolAndMethodAsProperty()
func (c *ExternExportsPassTest) Test_exportClassAsSymbolAndMethodAsProperty() {
	c.compileAndCheck(
		`/** @const */ var a = {};
/** @const */ a.b = {};
a.b.c = class {
  constructor(d, e, f) {
    /** @export {number} */ this.thisProp = 1;
  }
  /**
   * @param {string} a
   * @return {number}
   */
  method(a) {}

  static staticMethod() {}
};
goog.exportSymbol('foobar', a.b.c)
goog.exportProperty(a.b.c.prototype, 'exportedMethod', a.b.c.prototype.method)
goog.exportProperty(a.b.c, 'exportedStaticMethod', a.b.c.staticMethod)
`,
		// TODO(b/123352214): `@this {!a.b.c}` should be renamed to `foobar`
		`/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @constructor
 */
var foobar = function(d, e, f) {
};
/**
 * @return {undefined}
 * @this {(typeof a.b.c)}
 */
foobar.exportedStaticMethod = function() {
};
/**
 * @param {string} a
 * @return {number}
 * @this {!a.b.c}
 */
foobar.prototype.exportedMethod = function(a) {
};
foobar.prototype.thisProp;
`)
}

// Java: ExternExportsPassTest.noAtThisIsGeneratedForClassMethodWhenExportedClassNameMatches()
func (c *ExternExportsPassTest) Test_noAtThisIsGeneratedForClassMethodWhenExportedClassNameMatches() {
	c.compileAndCheck(
		`class Foo {
  /**
   * @param {string} a
   * @return {number}
   */
  method(a) {}
};
goog.exportSymbol('Foo', Foo)
goog.exportProperty(Foo.prototype, 'method', Foo.prototype.method)
`,
		// @this not generated on the method because exported and internal class names are the same
		`/**
 * @constructor
 */
var Foo = function() {
};
/**
 * @param {string} a
 * @return {number}
 */
Foo.prototype.method = function(a) {
};
`)
}

// Java: ExternExportsPassTest.exportClassHierarchyAsSymbols()
func (c *ExternExportsPassTest) Test_exportClassHierarchyAsSymbols() {
	c.compileAndCheck(
		`/** @const */ var a = {};
/** @const */ a.b = {};
a.b.c = class {
  constructor(d, e, f) {}
};
goog.exportSymbol('foobar', a.b.c)
a.b.d = class extends a.b.c {
  constructor(d, e, f) {}
};
goog.exportSymbol('bazboff', a.b.d)
`,
		// NOTE: foobar and bazboff end up sorted in ASCII order. This is expected.
		// TODO(b/123352214): @extends {a.b.c} should be updated
		`/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @extends {a.b.c}
 * @constructor
 */
var bazboff = function(d, e, f) {
};
/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @constructor
 */
var foobar = function(d, e, f) {
};
`)
}

// Java: ExternExportsPassTest.testInterface()
func (c *ExternExportsPassTest) TestInterface() {
	c.compileAndCheck(
		`/** @interface */ function Iface() {};
goog.exportSymbol('Iface', Iface)
`,
		`/**
 * @interface
 */
var Iface = function() {
};
`)
}

// Java: ExternExportsPassTest.exportInterfaceDefinedWithClass()
func (c *ExternExportsPassTest) Test_exportInterfaceDefinedWithClass() {
	c.compileAndCheck(
		`/** @interface */ class Iface {};
goog.exportSymbol('Iface', Iface)
`,
		`/**
 * @interface
 */
var Iface = function() {
};
`)
}

// Java: ExternExportsPassTest.exportInterfaceExtendedWithClass()
func (c *ExternExportsPassTest) Test_exportInterfaceExtendedWithClass() {
	c.compileAndCheck(
		`/** @interface */ class Iface {};
 goog.exportSymbol('Iface', Iface)
/** @interface */ class ExtendedIface extends Iface {};
 goog.exportSymbol('ExtendedIface', ExtendedIface)
`,
		// order of the interfaces is reversed, but that doesn't matter
		`/**
 * @extends {Iface}
 * @interface
 */
var ExtendedIface = function() {
};
/**
 * @interface
 */
var Iface = function() {
};
`)
}

// Java: ExternExportsPassTest.exportInterfaceDefinedWithClassThatExtendsUnexportedInterface()
func (c *ExternExportsPassTest) Test_exportInterfaceDefinedWithClassThatExtendsUnexportedInterface() {
	c.compileAndCheck(
		`/** @interface */ class Iface {}; // not exported
/** @interface */ class ExtendedIface extends Iface {};
 goog.exportSymbol('ExtendedIface', ExtendedIface)
`,
		`/**
 * @extends {Iface}
 * @interface
 */
var ExtendedIface = function() {
};
`)
}

// Java: ExternExportsPassTest.testRecord()
func (c *ExternExportsPassTest) TestRecord() {
	c.compileAndCheck(
		"/** @record */ function Iface() {}; goog.exportSymbol('Iface', Iface)",
		`/**
 * @record
 */
var Iface = function() {
};
`)
}

// Java: ExternExportsPassTest.exportRecordDefinedWithClass()
func (c *ExternExportsPassTest) Test_exportRecordDefinedWithClass() {
	c.compileAndCheck(
		`/** @record */ class Iface {};
goog.exportSymbol('Iface', Iface)
`,
		`/**
 * @record
 */
var Iface = function() {
};
`)
}

// Java: ExternExportsPassTest.exportRecordExtendedWithClass()
func (c *ExternExportsPassTest) Test_exportRecordExtendedWithClass() {
	c.compileAndCheck(
		`/** @record */ class Iface {};
 goog.exportSymbol('Iface', Iface)
/** @record */ class ExtendedIface extends Iface {};
 goog.exportSymbol('ExtendedIface', ExtendedIface)
`,
		// order of the interfaces is reversed, but that doesn't matter
		`/**
 * @extends {Iface}
 * @record
 */
var ExtendedIface = function() {
};
/**
 * @record
 */
var Iface = function() {
};
`)
}

// Java: ExternExportsPassTest.testExportSymbolDefinedInVar()
func (c *ExternExportsPassTest) TestExportSymbolDefinedInVar() {
	c.compileAndCheck(
		"var a = function(d, e, f) {}; goog.exportSymbol('foobar', a)",
		`/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @return {undefined}
 */
var foobar = function(d, e, f) {
};
`)
}

// Java: ExternExportsPassTest.exportClassDefinedWithConst()
func (c *ExternExportsPassTest) Test_exportClassDefinedWithConst() {
	c.compileAndCheck(
		`const a = class {
  constructor(d, e, f) {}
};
goog.exportSymbol('foobar', a)
`,
		// const becomes a var because we only generate var definitions for now
		`/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @constructor
 */
var foobar = function(d, e, f) {
};
`)
}

// Java: ExternExportsPassTest.exportClassDefinedWithLet()
func (c *ExternExportsPassTest) Test_exportClassDefinedWithLet() {
	c.compileAndCheck(
		`let a = class {
  constructor(d, e, f) {}
};
goog.exportSymbol('foobar', a)
`,
		// The let becomes a var because we only generate var definitions for now
		`/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @constructor
 */
var foobar = function(d, e, f) {
};
`)
}

// Java: ExternExportsPassTest.testExportProperty()
func (c *ExternExportsPassTest) TestExportProperty() {
	c.compileAndCheck(
		`/** @const */ var a = {};
/** @const */ a.b = {};
a.b.c = function(d, e, f) {};
goog.exportProperty(a.b, 'cprop', a.b.c)
`,
		`var a;
a.b;
/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @return {undefined}
 */
a.b.cprop = function(d, e, f) {
};
`)
}

// Java: ExternExportsPassTest.exportClassAsNamespaceProperty()
func (c *ExternExportsPassTest) Test_exportClassAsNamespaceProperty() {
	c.compileAndCheck(
		`const a = {};
/** @const */ a.b = {};
a.b.c = class {
  constructor(d, e, f) {}
};
goog.exportProperty(a.b, 'cprop', a.b.c)
`,
		`/**
 * @const
 * @suppress {const,duplicate}
 */
var a = {};
/**
 * @const
 * @suppress {const,duplicate}
 */
a.b = {};
/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @constructor
 */
a.b.cprop = function(d, e, f) {
};
`)
}

// Java: ExternExportsPassTest.exportQNameFunctionAsSymbolPlusPropertiesOnIt()
func (c *ExternExportsPassTest) Test_exportQNameFunctionAsSymbolPlusPropertiesOnIt() {
	c.compileAndCheck(
		`/** @const */ var a = {};
a.b = function(p1) {};
/** @constructor */
a.b.c = function(d, e, f) {};
/** @return {number} */
a.b.c.prototype.method = function(g, h, i) {};
goog.exportSymbol('a.b', a.b);
goog.exportProperty(a.b, 'c', a.b.c);
goog.exportProperty(a.b.c.prototype, 'exportedMethod', a.b.c.prototype.method);
`,
		`var a;
/**
 * @param {?} p1
 * @return {undefined}
 */
a.b = function(p1) {
};
/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @constructor
 */
a.b.c = function(d, e, f) {
};
/**
 * @param {?} g
 * @param {?} h
 * @param {?} i
 * @return {number}
 */
a.b.c.prototype.exportedMethod = function(g, h, i) {
};
`)
}

// Java: ExternExportsPassTest.exportQNameClassAsSymbolPlusPropertiesOnIt()
func (c *ExternExportsPassTest) Test_exportQNameClassAsSymbolPlusPropertiesOnIt() {
	c.compileAndCheck(
		`const a = {};
a.b = class {
  constructor(p1) {
  }
};
a.b.c = function(d, e, f) {};
a.b.prototype.c = function(g, h, i) {};
goog.exportSymbol('a.b', a.b);
goog.exportProperty(a.b, 'c', a.b.c);
goog.exportProperty(a.b.prototype, 'c', a.b.prototype.c);
`,
		`/**
 * @const
 * @suppress {const,duplicate}
 */
var a = {};
/**
 * @param {?} p1
 * @constructor
 */
a.b = function(p1) {
};
/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @return {undefined}
 */
a.b.c = function(d, e, f) {
};
/**
 * @param {?} g
 * @param {?} h
 * @param {?} i
 * @return {undefined}
 */
a.b.prototype.c = function(g, h, i) {
};
`)
}

// Java: ExternExportsPassTest.symbolsRenamedInExportAreAlsoRenamedInExportedChildQNames()
func (c *ExternExportsPassTest) Test_symbolsRenamedInExportAreAlsoRenamedInExportedChildQNames() {
	c.compileAndCheck(
		`/** @const */ var a = {};
a.b = function(p1) {};
a.b.c = function(d, e, f) {};
a.b.prototype.c = function(g, h, i) {};
goog.exportSymbol('hello', a);
goog.exportProperty(a.b, 'c', a.b.c);
goog.exportProperty(a.b.prototype, 'c', a.b.prototype.c);
`,
		`/** @type {{b: function(?): undefined}} */
var hello = {};
hello.b;
/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @return {undefined}
 */
hello.b.c = function(d, e, f) {
};
/**
 * @param {?} g
 * @param {?} h
 * @param {?} i
 * @return {undefined}
 */
hello.b.prototype.c = function(g, h, i) {
};
`)
}

// Java: ExternExportsPassTest.exportingQnameAsSimpleNameAffectsLaterPropertyExports()
func (c *ExternExportsPassTest) Test_exportingQnameAsSimpleNameAffectsLaterPropertyExports() {
	c.compileAndCheck(
		`/** @const */ var a = {};
a.b = function(p1) {};
a.b.c = function(d, e, f) {};
a.b.prototype.c = function(g, h, i) {};
goog.exportSymbol('prefix', a.b);
goog.exportProperty(a.b, 'c', a.b.c);
`,
		`/**
 * @param {?} p1
 * @return {undefined}
 */
var prefix = function(p1) {
};
/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @return {undefined}
 */
prefix.c = function(d, e, f) {
};
`)
}

// Java: ExternExportsPassTest.testExportNonStaticSymbol()
func (c *ExternExportsPassTest) TestExportNonStaticSymbol() {
	c.compileAndCheck(
		`/** @const */ var a = {};
/** @const */ a.b = {};
/** @const */ var d = {};
a.b.c = d;
goog.exportSymbol('foobar', a.b.c)
`,
		"var foobar;\n")
}

// Java: ExternExportsPassTest.testExportNonStaticSymbol2()
func (c *ExternExportsPassTest) TestExportNonStaticSymbol2() {
	c.compileAndCheck(
		`/** @const */ var a = {};
/** @const */ a.b = {};
var d = function() {};
a.b.c = d;
goog.exportSymbol('foobar', a.b.c())
`,
		"var foobar;\n")
}

// Java: ExternExportsPassTest.testExportNonexistentProperty()
func (c *ExternExportsPassTest) TestExportNonexistentProperty() {
	c.compileAndCheck(
		`/** @fileoverview @suppress {missingProperties} */
/** @const */ var a = {};
/** @const */ a.b = {};
a.b.c = function(d, e, f) {};
goog.exportProperty(a.b, 'none', a.b.none)
`,
		`var a;
a.b;
a.b.none;
`)
}

// Java: ExternExportsPassTest.testExportSymbolWithTypeAnnotation()
func (c *ExternExportsPassTest) TestExportSymbolWithTypeAnnotation() {

	c.compileAndCheck(
		`var internalName;
/**
 * @param {string} param1
 * @param {number} param2
 * @return {string}
 */
internalName = function(param1, param2) {
return param1 + param2;
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {string} param1
 * @param {number} param2
 * @return {string}
 */
var externalName = function(param1, param2) {
};
`)
}

// Java: ExternExportsPassTest.testExportSymbolWithTemplateAnnotation()
func (c *ExternExportsPassTest) TestExportSymbolWithTemplateAnnotation() {

	c.compileAndCheck(
		`var internalName;
/**
 * @param {T} param1
 * @return {T}
 * @template T
 */
internalName = function(param1) {
return param1;
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {T} param1
 * @return {T}
 * @template T
 */
var externalName = function(param1) {
};
`)
}

// Java: ExternExportsPassTest.exportSubClassWithoutConstructor()
func (c *ExternExportsPassTest) Test_exportSubClassWithoutConstructor() {
	c.compileAndCheck(
		`class SuperClass {
  /**
   * @param {number} num
   */
   constructor(num) {
   }
}
goog.exportSymbol('SuperClass', SuperClass);
class SubClass extends SuperClass {
}
goog.exportSymbol('SubClass', SubClass);
`,
		// The "a" parameter name is automatically generated
		`/**
 * @param {number} a
 * @extends {SuperClass}
 * @constructor
 */
var SubClass = function(a) {
};
/**
 * @param {number} num
 * @constructor
 */
var SuperClass = function(num) {
};
`)
}

// Java: ExternExportsPassTest.exportClassWithTemplateAnnotation()
func (c *ExternExportsPassTest) Test_exportClassWithTemplateAnnotation() {

	c.compileAndCheck(
		`var internalName;
/**
 * @template T
 */
internalName = class {
  /**
   * @param {T} param1
   */
  constructor(param1) {
    /** @const */
    this.data = param1;
  }
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {T} param1
 * @constructor
 * @template T
 */
var externalName = function(param1) {
};
`)
}

// Java: ExternExportsPassTest.testExportSymbolWithMultipleTemplateAnnotation()
func (c *ExternExportsPassTest) TestExportSymbolWithMultipleTemplateAnnotation() {

	c.compileAndCheck(
		`var internalName;

/**
 * @param {K} param1
 * @return {V}
 * @template K,V
 */
internalName = function(param1) {
  return /** @type {?} */ (param1);
};
goog.exportSymbol('externalName', internalName);
`,
		`/**
 * @param {K} param1
 * @return {V}
 * @template K,V
 */
var externalName = function(param1) {
};
`)
}

// Java: ExternExportsPassTest.testExportSymbolWithoutTypeCheck()
func (c *ExternExportsPassTest) TestExportSymbolWithoutTypeCheck() {
	// ExternExportsPass should not emit annotations
	// if there is no type information available.
	c.DisableTypeCheck()

	c.compileAndCheck(
		`var internalName;

/**
 * @param {string} param1
 * @param {number} param2
 * @return {string}
 */
internalName = function(param1, param2) {
return param1 + param2;
};
goog.exportSymbol('externalName', internalName)
`,
		`var externalName = function(param1, param2) {
};
`)
}

// Java: ExternExportsPassTest.testExportSymbolWithConstructor()
func (c *ExternExportsPassTest) TestExportSymbolWithConstructor() {
	c.compileAndCheck(
		`var internalName;

/**
 * @constructor
 */
internalName = function() {
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @constructor
 */
var externalName = function() {
};
`)
}

// Java: ExternExportsPassTest.exportClass()
func (c *ExternExportsPassTest) Test_exportClass() {
	c.compileAndCheck(
		`var internalName;

internalName = class {
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @constructor
 */
var externalName = function() {
};
`)
}

// Java: ExternExportsPassTest.testNonNullTypes()
func (c *ExternExportsPassTest) TestNonNullTypes() {
	c.compileAndCheck(
		`/**
 * @constructor
 */
function Foo() {}
goog.exportSymbol('Foo', Foo);
/**
 * @param {!Foo} x
 * @return {!Foo}
 */
Foo.f = function(x) { return x; };
goog.exportProperty(Foo, 'f', Foo.f);
`,
		`/**
 * @constructor
 */
var Foo = function() {
};
/**
 * @param {!Foo} x
 * @return {!Foo}
 */
Foo.f = function(x) {
};
`)
}

// Java: ExternExportsPassTest.testExportSymbolWithConstructorWithoutTypeCheck()
func (c *ExternExportsPassTest) TestExportSymbolWithConstructorWithoutTypeCheck() {
	// For now, skipping type checking should prevent generating
	// annotations of any kind, so, e.g., @constructor is not preserved.
	// This is probably not ideal, but since JSDocInfo for functions is attached
	// to JSTypes and not Nodes (and no JSTypes are created when checkTypes
	// is false), we don't really have a choice.

	c.DisableTypeCheck()

	c.compileAndCheck(
		`var internalName;
/**
 * @constructor
 */
internalName = function() {
};
goog.exportSymbol('externalName', internalName)
`,
		`var externalName = function() {
};
`)
}

// Java: ExternExportsPassTest.exportClassWithoutTypeCheck()
func (c *ExternExportsPassTest) Test_exportClassWithoutTypeCheck() {
	c.T().Skip("@Ignore: (b/141729691): this test fails, but unlikely to fix as this feature is deprecated")
	// For now, skipping type checking should prevent generating
	// annotations of any kind, so, e.g., @constructor is not preserved.
	// This is probably not ideal, but since JSDocInfo for functions is attached
	// to JSTypes and not Nodes (and no JSTypes are created when checkTypes
	// is false), we don't really have a choice.

	c.DisableTypeCheck()

	c.compileAndCheck(
		`var internalName;
internalName = class {
};
goog.exportSymbol('externalName', internalName)
`,
		`var externalName = function() {
};
`)
}

// x.Y is present in the generated externs but lacks the @constructor annotation.
//
// Java: ExternExportsPassTest.testExportPrototypePropsWithoutConstructor()
func (c *ExternExportsPassTest) TestExportPrototypePropsWithoutConstructor() {
	c.compileAndCheck(
		`/** @constructor */
x.Y = function() {};
x.Y.prototype.z = function() {};
goog.exportProperty(x.Y.prototype, 'z', x.Y.prototype.z);
`,
		`var x;
x.Y;
/**
 * @return {undefined}
 */
x.Y.prototype.z = function() {
};
`)
}

// x.Y is present in the generated externs but lacks the @constructor annotation.
//
// Java: ExternExportsPassTest.exportMethodButNotTheClass()
func (c *ExternExportsPassTest) Test_exportMethodButNotTheClass() {
	c.compileAndCheck(
		`x.Y = class {
  z() {};
};
goog.exportProperty(x.Y.prototype, 'z', x.Y.prototype.z);
`,
		`var x;
x.Y;
/**
 * @return {undefined}
 */
x.Y.prototype.z = function() {
};
`)
}

// Java: ExternExportsPassTest.testExportFunctionWithOptionalArguments1()
func (c *ExternExportsPassTest) TestExportFunctionWithOptionalArguments1() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {number=} a
 */
internalName = function(a) {
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {number=} a
 * @return {undefined}
 */
var externalName = function(a) {
};
`)
}

// Java: ExternExportsPassTest.testExportFunctionWithOptionalArguments2()
func (c *ExternExportsPassTest) TestExportFunctionWithOptionalArguments2() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {number=} a
 */
internalName = function(a) {
  return /** @type {?} */ (6);
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {number=} a
 * @return {?}
 */
var externalName = function(a) {
};
`)
}

// Java: ExternExportsPassTest.testExportFunctionWithOptionalArguments3()
func (c *ExternExportsPassTest) TestExportFunctionWithOptionalArguments3() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {number=} a
 */
internalName = function(a) {
  return /** @type {?} */ (a);
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {number=} a
 * @return {?}
 */
var externalName = function(a) {
};
`)
}

// Java: ExternExportsPassTest.testExportFunctionWithVariableArguments()
func (c *ExternExportsPassTest) TestExportFunctionWithVariableArguments() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {...number} a
 * @return {number}
 */
internalName = function(a) {
  return 6;
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {...number} a
 * @return {number}
 */
var externalName = function(a) {
};
`)
}

// Java: ExternExportsPassTest.exportArrowFunction()
func (c *ExternExportsPassTest) Test_exportArrowFunction() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {number} a
 * @return {number}
 */
internalName = (a) => a;
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {number} a
 * @return {number}
 */
var externalName = function(a) {
};
`)
}

// Java: ExternExportsPassTest.exportGeneratorFunction()
func (c *ExternExportsPassTest) Test_exportGeneratorFunction() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {number} a
 * @return {!Iterator<number>}
 */
internalName = function *(a) { yield a; };
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {number} a
 * @return {!Iterator<number,?,?>}
 */
var externalName = function(a) {
};
`)
}

// Java: ExternExportsPassTest.exportAsyncGeneratorFunction()
func (c *ExternExportsPassTest) Test_exportAsyncGeneratorFunction() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {number} a
 * @return {!AsyncGenerator<number>}
 */
internalName = async function *(a) { yield a; };
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {number} a
 * @return {!AsyncGenerator<number,?,?>}
 */
var externalName = function(a) {
};
`)
}

// Java: ExternExportsPassTest.exportAsyncFunction()
func (c *ExternExportsPassTest) Test_exportAsyncFunction() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {number} a
 * @return {!Promise<number>}
 */
internalName = async function(a) { return a; };
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {number} a
 * @return {!Promise<number>}
 */
var externalName = function(a) {
};
`)
}

// Java: ExternExportsPassTest.exportFunctionWithDefaultParameter()
func (c *ExternExportsPassTest) Test_exportFunctionWithDefaultParameter() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {number=} a
 */
internalName = function(a = 1) {
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {number=} a
 * @return {undefined}
 */
var externalName = function(a) {
};
`)
}

// Java: ExternExportsPassTest.exportFunctionWithRestParameter()
func (c *ExternExportsPassTest) Test_exportFunctionWithRestParameter() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {...number} numbers
 */
internalName = function(...numbers) {
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {...number} numbers
 * @return {undefined}
 */
var externalName = function(numbers) {
};
`)
}

// Java: ExternExportsPassTest.exportFunctionWithPatternParameters()
func (c *ExternExportsPassTest) Test_exportFunctionWithPatternParameters() {
	c.compileAndCheck(
		`var internalName;

/**
 * @param {{p: number, q: string}} options
 * @param {number} a
 * @param {!Array<number>} moreOptions
 */
internalName = function({p, q}, a, [x, y]) {
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {{p: number, q: string}} b
 * @param {number} a
 * @param {!Array<number>} c
 * @return {undefined}
 */
var externalName = function(b, a, c) {
};
`)
}

// Java: ExternExportsPassTest.exportFunctionWithInlineDeclaredPatternParameters()
func (c *ExternExportsPassTest) Test_exportFunctionWithInlineDeclaredPatternParameters() {
	c.compileAndCheck(
		`var internalName;

internalName = function(
    /** !Array<string> */ [p, q],
    /** number */ a, // generated names start with 'a', make sure we don't conflict
    /** !{x: number, y: string} */ {x, y}) {
};
goog.exportSymbol('externalName', internalName)
`,
		// Param names are "b, a, c" because name 'a' was taken by a named parameter
		`/**
 * @param {!Array<string>} b
 * @param {number} a
 * @param {{x: number, y: string}} c
 * @return {undefined}
 */
var externalName = function(b, a, c) {
};
`)
}

// Java: ExternExportsPassTest.testExportEnum()
func (c *ExternExportsPassTest) TestExportEnum() {
	// We don't care what the values of the object properties are.
	// They're ignored by the type checker, and even if they weren't, it'd
	// be incomputable to get them correct in all cases
	// (think complex objects).
	c.compileAndCheck(
		`/**
 * @enum {string}
 * @export
 */
var E = {A:'a', B:'b'};
goog.exportSymbol('E', E);
`,
		`/** @enum {string} */
var E = {A:1, B:2};
`)
}

// Java: ExternExportsPassTest.testExportWithReferenceToEnum()
func (c *ExternExportsPassTest) TestExportWithReferenceToEnum() {
	js := `/**
 * @enum {number}
 * @export
 */
var E = {A:1, B:2};
goog.exportSymbol('E', E);

/**
 * @param {!E} e
 * @export
 */
function f(e) {}
goog.exportSymbol('f', f);
`
	expected := `/** @enum {number} */
var E = {A:1, B:2};
/**
 * @param {E} e
 * @return {undefined}
 */
var f = function(e) {
};
`

	// NOTE: The type should print {E} for the @param, but is not.
	c.compileAndCheck(js, jutil.StringReplace(expected, "{E}", "{number}"))
}

// Java: ExternExportsPassTest.exportDontEmitPrototypePathPrefixForClassMethod()
func (c *ExternExportsPassTest) Test_exportDontEmitPrototypePathPrefixForClassMethod() {
	c.compileAndCheck(
		`var Foo = class {
  /**
   * @return {number}
   */
  m() {return 6;};
};
goog.exportSymbol('Foo', Foo);
goog.exportProperty(Foo.prototype, 'm', Foo.prototype.m);
`,
		`/**
 * @constructor
 */
var Foo = function() {
};
/**
 * @return {number}
 */
Foo.prototype.m = function() {
};
`)
}

// Java: ExternExportsPassTest.useExportsAsExternsWithClass()
func (c *ExternExportsPassTest) Test_useExportsAsExternsWithClass() {
	librarySource := `var InternalName = class {
  /**
   * @param {number} n
   */
  constructor(n) {
  }
};
goog.exportSymbol('ExternalName', InternalName)
`

	clientSource := `var foo = new ExternalName(6);
/**
 * @param {ExternalName} x
 */
var bar = function(x) {};
`

	c.compileAndExportExternsWithConsumer(
		librarySource,
		CompilerTestCase_MINIMAL_EXTERNS,
		func(generatedExterns *string) {
			c.compileAndExportExternsWithExterns(clientSource, CompilerTestCase_MINIMAL_EXTERNS+jutil.ToJavaString(generatedExterns))
		})
}

// Java: ExternExportsPassTest.testDontWarnOnExportFunctionWithUnknownReturnType()
func (c *ExternExportsPassTest) TestDontWarnOnExportFunctionWithUnknownReturnType() {
	librarySource := `var InternalName = function() {
  return 6;
};
goog.exportSymbol('ExternalName', InternalName)
`

	c.compileAndExportExterns(librarySource)
}

// Java: ExternExportsPassTest.testDontWarnOnExportConstructorWithUnknownReturnType()
func (c *ExternExportsPassTest) TestDontWarnOnExportConstructorWithUnknownReturnType() {
	librarySource := `/**
 * @constructor
 */
var InternalName = function() {
};
goog.exportSymbol('ExternalName', InternalName)
`

	c.compileAndExportExterns(librarySource)
}

// Java: ExternExportsPassTest.testTypedef()
func (c *ExternExportsPassTest) TestTypedef() {
	c.compileAndCheck(
		`/** @typedef {{x: number, y: number}} */ var Coord;
/**
 * @param {Coord} a
 * @export
 */
var fn = function(a) {};
goog.exportSymbol('fn', fn);
`,
		`/**
 * @param {{x: number, y: number}} a
 * @return {undefined}
 */
var fn = function(a) {
};
`)
}

// Java: ExternExportsPassTest.testExportParamWithNull()
func (c *ExternExportsPassTest) TestExportParamWithNull() {
	c.compileAndCheck(
		`/** @param {string|null=} d */
var f = function(d) {};
goog.exportSymbol('foobar', f)
`,
		`/**
 * @param {(null|string)=} d
 * @return {undefined}
 */
var foobar = function(d) {
};
`)
}

// Java: ExternExportsPassTest.testExportConstructor()
func (c *ExternExportsPassTest) TestExportConstructor() {
	c.compileAndCheck(
		"/** @constructor */ var a = function() {}; goog.exportSymbol('foobar', a)",
		`/**
 * @constructor
 */
var foobar = function() {
};
`)
}

// Java: ExternExportsPassTest.exportEs5ClassHierarchy()
func (c *ExternExportsPassTest) Test_exportEs5ClassHierarchy() {
	c.compileAndCheck(
		`/** @constructor */
function SuperClass() {}
goog.exportSymbol('Foo', SuperClass);

/**
 * @constructor
 * @extends {SuperClass}
 */
function SubClass() {}
goog.exportSymbol('Bar', SubClass);
`,
		// TODO(b/123352214): SuperClass should be called Foo in exported @extends annotation
		`/**
 * @extends {SuperClass}
 * @constructor
 */
var Bar = function() {
};
/**
 * @constructor
 */
var Foo = function() {
};
`)
}

// Java: ExternExportsPassTest.testExportLocalPropertyInConstructor()
func (c *ExternExportsPassTest) TestExportLocalPropertyInConstructor() {
	c.compileAndCheck(
		"/** @constructor */function F() { /** @export */ this.x = 5;} goog.exportSymbol('F', F);",
		`/**
 * @constructor
 */
var F = function() {
};
F.prototype.x;
`)
}

// Java: ExternExportsPassTest.testExportLocalPropertyInConstructor2()
func (c *ExternExportsPassTest) TestExportLocalPropertyInConstructor2() {
	c.compileAndCheck(
		`/** @constructor */function F() { /** @export */ this.x = 5;}
goog.exportSymbol('F', F);
goog.exportProperty(F.prototype, 'x', F.prototype.x);
`,
		`/**
 * @constructor
 */
var F = function() {
};
F.prototype.x;
`)
}

// Java: ExternExportsPassTest.testExportLocalPropertyInConstructor3()
func (c *ExternExportsPassTest) TestExportLocalPropertyInConstructor3() {
	c.compileAndCheck(
		"/** @constructor */function F() { /** @export */ this.x;} goog.exportSymbol('F', F);",
		`/**
 * @constructor
 */
var F = function() {
};
F.prototype.x;
`)
}

// Java: ExternExportsPassTest.testExportLocalPropertyInConstructor4()
func (c *ExternExportsPassTest) TestExportLocalPropertyInConstructor4() {
	c.compileAndCheck(
		`/** @constructor */
function F() { /** @export */ this.x = function(/** string */ x){};}
goog.exportSymbol('F', F);
`,
		`/**
 * @constructor
 */
var F = function() {
};
F.prototype.x;
`)
}

// Java: ExternExportsPassTest.testExportLocalPropertyNotInConstructor()
func (c *ExternExportsPassTest) TestExportLocalPropertyNotInConstructor() {
	c.compileAndCheck(
		"/** @this {?} */ function f() { /** @export */ this.x = 5;} goog.exportSymbol('f', f);",
		`/**
 * @return {undefined}
 */
var f = function() {
};
`)
}

// Java: ExternExportsPassTest.testExportParamWithSymbolDefinedInFunction()
func (c *ExternExportsPassTest) TestExportParamWithSymbolDefinedInFunction() {
	c.compileAndCheck(
		`var id = function() {return /** @type {?} */ ('id')};
var ft = function() {
  var id;
  return 1;
};
goog.exportSymbol('id', id);
`,
		`/**
 * @return {?}
 */
var id = function() {
};
`)
}

// Java: ExternExportsPassTest.testExportSymbolWithFunctionDefinedAsFunction()
func (c *ExternExportsPassTest) TestExportSymbolWithFunctionDefinedAsFunction() {

	c.compileAndCheck(
		`/**
 * @param {string} param1
 * @return {string}
 */
function internalName(param1) {
  return param1
};
goog.exportSymbol('externalName', internalName)
`,
		`/**
 * @param {string} param1
 * @return {string}
 */
var externalName = function(param1) {
};
`)
}

// Java: ExternExportsPassTest.testExportSymbolWithFunctionAlias()
func (c *ExternExportsPassTest) TestExportSymbolWithFunctionAlias() {

	c.compileAndCheck(
		`/**
 * @param {string} param1
 */
var y = function(param1) {
};
/**
 * @param {string} param1
 * @param {string} param2
 */
var x = function y(param1, param2) {
};
goog.exportSymbol('externalName', y)
`,
		`/**
 * @param {string} param1
 * @return {undefined}
 */
var externalName = function(param1) {
};
`)
}

// Java: ExternExportsPassTest.testNamespaceDefinitionInExterns()
func (c *ExternExportsPassTest) TestNamespaceDefinitionInExterns() {
	c.compileAndCheck(
		`/** @const */
var ns = {};
/** @const */
ns.subns = {};
/** @constructor */
ns.subns.Foo = function() {};
goog.exportSymbol('ns.subns.Foo', ns.subns.Foo);
`,
		`/**
 * @const
 * @suppress {const,duplicate}
 */
var ns = {};
/**
 * @const
 * @suppress {const,duplicate}
 */
ns.subns = {};
/**
 * @constructor
 */
ns.subns.Foo = function() {
};
`)
}

// Java: ExternExportsPassTest.testNullabilityInFunctionTypes()
func (c *ExternExportsPassTest) TestNullabilityInFunctionTypes() {
	c.compileAndCheck(
		`/**
 * @param {function(Object)} takesNullable
 * @param {function(!Object)} takesNonNullable
 */
function x(takesNullable, takesNonNullable) {}
goog.exportSymbol('x', x);
`,
		`/**
 * @param {function((Object|null)): ?} takesNullable
 * @param {function(!Object): ?} takesNonNullable
 * @return {undefined}
 */
var x = function(takesNullable, takesNonNullable) {
};
`)
}

// Java: ExternExportsPassTest.testNullabilityInRecordTypes()
func (c *ExternExportsPassTest) TestNullabilityInRecordTypes() {
	c.compileAndCheck(
		`/** @typedef {{ nonNullable: !Object, nullable: Object }} */
var foo;
/** @param {foo} record */
function x(record) {}
goog.exportSymbol('x', x);
`,
		`/**
 * @param {{nonNullable: !Object, nullable: (Object|null)}} record
 * @return {undefined}
 */
var x = function(record) {
};
`)
}

// Java: ExternExportsPassTest.testDontWarnOnExportFunctionWithUnknownParameterTypes()
func (c *ExternExportsPassTest) TestDontWarnOnExportFunctionWithUnknownParameterTypes() {
	/* This source is missing types for the b and c parameters */
	librarySource := `/**
 * @param {number} a
 * @return {number}
 */
var InternalName = function(a,b,c) {
  return 6;
};
goog.exportSymbol('ExternalName', InternalName)
`

	c.compileAndExportExterns(librarySource)
}

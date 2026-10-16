// Port of com/google/javascript/jscomp/Es6TranspilationIntegrationTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/Es6TranspilationIntegrationTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

// Es6TranspilationIntegrationTest tests the transpilation passes together, as
// TranspilationPasses.addTranspilationPasses adds them.
type Es6TranspilationIntegrationTest struct {
	CompilerTestCaseBase
}

// Java: Es6TranspilationIntegrationTest()
func NewEs6TranspilationIntegrationTest() *Es6TranspilationIntegrationTest {
	c := &Es6TranspilationIntegrationTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, Es6TranspilationIntegrationTest_GetDefaultExterns())
	return c
}

func TestEs6TranspilationIntegrationTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6TranspilationIntegrationTest() })
}

// Java: Es6TranspilationIntegrationTest.getDefaultExterns()
func Es6TranspilationIntegrationTest_GetDefaultExterns() string {
	return Es6TranspilationIntegrationTest_GetDefaultExternsBuilder().Build()
}

// Java: Es6TranspilationIntegrationTest.getDefaultExternsBuilder()
func Es6TranspilationIntegrationTest_GetDefaultExternsBuilder() *TestExternsBuilder {
	return NewTestExternsBuilder().
		AddAsyncIterable().
		AddArray().
		AddArguments().
		AddObject().
		AddMath().
		AddExtra(
			// stubs of runtime libraries
			`/** @const */
var $jscomp = {};
$jscomp.generator = {};
$jscomp.generator.createGenerator = function() {};
/** @constructor */
$jscomp.generator.Context = function() {};
/** @constructor */
$jscomp.generator.Context.PropertyIterator = function() {};
$jscomp.asyncExecutePromiseGeneratorFunction = function(program) {};
`)
}

// Java: Es6TranspilationIntegrationTest.setUp()
func (c *Es6TranspilationIntegrationTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2016)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	// We will do normalization as part of transpilation in getProcessor()
	c.DisableNormalize()
	c.DisableTypeCheck()
	c.DisableCompareJsDoc() // optimization passes see simplified JSDoc.
	// Normalization does renaming on the externs (parameter names and maybe other things) and the
	// pass that inverts renaming ignores the externs, leaving them changed.
	c.AllowExternsChanges()
	// Our getProcessor() returns a fake "pass" that actually runs several passes,
	// including one that undoes some of the changes that were made.
	// this confuses the logic for validating AST change marking.
	// That logic is really only valid when testing a single, real pass.
	c.DisableValidateAstChangeMarking()
	replacements := jutil.NewLinkedHashMap[string, string]()
	for k, v := range Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS.All() {
		replacements.Put(k, v)
	}
	jutil.CheckArgumentMsg1(!replacements.ContainsKey("KEY"), "duplicate key: %s", "KEY") // buildOrThrow()
	replacements.Put("KEY", "$jscomp$key$")
	c.SetGenericNameReplacements(replacements)
}

// Java: Es6TranspilationIntegrationTest.getProcessor(Compiler)
func (c *Es6TranspilationIntegrationTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)

	compilerOptions := compiler.GetOptions()
	passes := NewPassListBuilder(compilerOptions)

	passes.MaybeAdd(
		PassFactory_Builder_().
			SetName("es6InjectRuntimeLibraries").
			SetInternalFactory(func(abstractCompiler AbstractCompiler) CompilerPass {
				return NewInjectTranspilationRuntimeLibraries(abstractCompiler)
			}).
			Build())

	passes.MaybeAdd(
		PassFactory_Builder_().
			SetName("rewritePolyfills").
			SetInternalFactory(
				func(ac AbstractCompiler) CompilerPass {
					return NewRewritePolyfills(
						ac /* injectPolyfills= */, true /* isolatePolyfills= */, false, nil)
				}).
			Build())

	passes.MaybeAdd(
		PassFactory_Builder_().
			SetName("convertTypesToColors").
			SetInternalFactory(
				func(ac AbstractCompiler) CompilerPass {
					return NewConvertTypesToColors(
						ac, SerializationOptions_Builder_().SetIncludeDebugInfo(true).Build())
				}).
			Build())

	passes.MaybeAdd(
		PassFactory_Builder_().
			SetName(PassNames_NORMALIZE).
			SetInternalFactory(func(abstractCompiler AbstractCompiler) CompilerPass {
				return Normalize_Builder_(abstractCompiler).Build()
			}).
			Build())
	TranspilationPasses_AddTranspilationPasses(passes, compilerOptions)
	// Since we're testing the transpile-only case, we need to put back the original variable names
	// where possible once transpilation is complete. This matches the behavior in
	// DefaultPassConfig. See comments there for further explanation.
	passes.MaybeAdd(
		PassFactory_Builder_().
			SetName("invertContextualRenaming").
			SetInternalFactory(MakeDeclaredNamesUnique_GetContextualRenameInverter).
			Build())
	optimizer.Consume(passes.Build())

	return optimizer
}

// Java: actually.rewriteUniqueIdAndTest(Sources, Expected)
func (c *Es6TranspilationIntegrationTest) RewriteUniqueIdAndTest(srcs CompilerTestCase_Sources, originalExpected *CompilerTestCase_Expected) {
	externs := c.Externs("")
	c.Test(externs, srcs, originalExpected)
}

// Java: actually.testObjectLiteralStringKeysWithNoValue()
func (c *Es6TranspilationIntegrationTest) TestObjectLiteralStringKeysWithNoValue() {
	c.Test("var x = {a, b};", "var x = {a: a, b: b};")
	AssertThat(c.GetLastCompiler().GetRuntimeJsLibManager().GetInjectedLibraries()).IsEmpty()
}

// Java: actually.testSpreadLibInjection()
func (c *Es6TranspilationIntegrationTest) TestSpreadLibInjection() {
	c.Test("var x = [...a];", "var x=[].concat((0, $jscomp.arrayFromIterable)(a))")
	AssertThat(c.GetLastCompiler().GetRuntimeJsLibManager().GetInjectedLibraries()).
		ContainsExactly("es6/util/arrayfromiterable")
}

// Java: actually.testObjectLiteralMemberFunctionDef()
func (c *Es6TranspilationIntegrationTest) TestObjectLiteralMemberFunctionDef() {
	c.Test(
		"var x = {/** @return {number} */ a() { return 0; } };",
		"var x = {/** @return {number} */ a: function() { return 0; } };")
	AssertThat(c.GetLastCompiler().GetRuntimeJsLibManager().GetInjectedLibraries()).IsEmpty()
}

// Java: actually.testClassStatement()
func (c *Es6TranspilationIntegrationTest) TestClassStatement() {
	c.Test("class C { }", "/** @constructor */ var C = function() {};")
	c.Test("class C { constructor() {} }", "/** @constructor */ var C = function() {};")
	c.Test(
		"class C { method() {}; }",
		`/** @constructor */
var C = function() {};
C.prototype.method = function() {};
`)
	c.Test(
		"class C { constructor(a) { this.a = a; } }",
		"/** @constructor */ var C = function(a) { this.a = a; };")

	c.Test(
		"class C { constructor() {} foo() {} }",
		`/** @constructor */
var C = function() {};
C.prototype.foo = function() {};
`)

	c.Test(
		"class C { constructor() {}; foo() {}; bar() {} }",
		`/** @constructor */
var C = function() {};
C.prototype.foo = function() {};
C.prototype.bar = function() {};
`)

	c.Test(
		"class C { foo() {}; bar() {} }",
		`/** @constructor */
var C = function() {};
C.prototype.foo = function() {};
C.prototype.bar = function() {};
`)

	c.Test(
		`class C {
  constructor(a) { this.a = a; }

  foo() { console.log(this.a); }

  bar() { alert(this.a); }
}
`,
		`/** @constructor */
var C = function(a) { this.a = a; };
C.prototype.foo = function() { console.log(this.a); };
C.prototype.bar = function() { alert(this.a); };
`)

	c.RewriteUniqueIdAndTest(
		c.Srcs(
			`if (true) {
   class Foo{}
} else {
   class Foo{}
}
`),
		c.Expected(
			`if (true) {
    /** @constructor */
    var Foo = function() {};
} else {
    /** @constructor */
    var Foo$jscomp$1 = function() {};
}
`))
}

// Java: actually.testAnonymousSuper()
func (c *Es6TranspilationIntegrationTest) TestAnonymousSuper() {
	c.Test(
		"f(class extends D { f() { super.g() } })",
		`/** @const @constructor */
var CLASS_DECL$0 = function() {
  return D.apply(this, arguments) || this;
};
$jscomp.inherits(CLASS_DECL$0, D);
CLASS_DECL$0.prototype.f = function() {
  D.prototype.g.call(this);
};
f(CLASS_DECL$0);
`)
}

// Java: actually.testNewTarget()
func (c *Es6TranspilationIntegrationTest) TestNewTarget() {
	c.TestError("function Foo() { new.target; }", TranspilationUtil_CANNOT_CONVERT_YET)
	c.TestError("class Example { foo() { new.target; } }", TranspilationUtil_CANNOT_CONVERT_YET)
}

// Java: actually.testClassWithJsDoc()
func (c *Es6TranspilationIntegrationTest) TestClassWithJsDoc() {
	c.Test("class C { }", "/** @constructor */ var C = function() { };")

	c.Test(
		"/** @deprecated */ class C { }", "/** @constructor @deprecated */ var C = function() {};")

	c.Test(
		"/**              @dict */ class C              { }", //
		"/** @constructor @dict */ var   C = function() { };")

	c.Test(
		"/** @template T */ class C { }", "/** @constructor @template T */ var C = function() {};")

	c.Test("/** @final */ class C { }", "/** @constructor @final */ var C = function() {};")

	c.Test("/** @private */ class C { }", "/** @constructor @private */ var C = function() {};")
}

// Java: actually.testInterfaceWithJsDoc()
func (c *Es6TranspilationIntegrationTest) TestInterfaceWithJsDoc() {
	c.Test(
		`/**
 * Converts Xs to Ys.
 * @interface
 */
class Converter {
  /**
   * @param {X} x
   * @return {Y}
   */
  convert(x) {}
}
`,
		`/**
 * Converts Xs to Ys.
 * @interface
 */
var Converter = function() { };

/**
 * @param {X} x
 * @return {Y}
 */
Converter.prototype.convert = function(x) {};
`)
}

// Java: actually.testRecordWithJsDoc()
func (c *Es6TranspilationIntegrationTest) TestRecordWithJsDoc() {
	c.Test(
		`/**
 * @record
 */
class Converter {
  /**
   * @param {X} x
   * @return {Y}
   */
  convert(x) {}
}
`,
		`/**
 * @record
 */
var Converter = function() { };

/**
 * @param {X} x
 * @return {Y}
 */
Converter.prototype.convert = function(x) {};
`)
}

// Java: actually.testMemberWithJsDoc()
func (c *Es6TranspilationIntegrationTest) TestMemberWithJsDoc() {
	c.Test(
		"class C { /** @param {boolean} b */ foo(b) {} }",
		`/**
 * @constructor
 */
var C = function() {};

/** @param {boolean} b */
C.prototype.foo = function(b) {};
`)
}

// Java: actually.testClassStatementInsideIf()
func (c *Es6TranspilationIntegrationTest) TestClassStatementInsideIf() {
	c.Test("if (foo) { class C { } }", "if (foo) { /** @constructor */ var C = function() {}; }")

	c.Test("if (foo) class C {}", "if (foo) { /** @constructor */ var C = function() {}; }")
}

// Class expressions that are the RHS of a 'var' statement.
// Java: actually.testClassExpressionInVar()
func (c *Es6TranspilationIntegrationTest) TestClassExpressionInVar() {
	c.Test("var C = class { }", "/** @constructor */ var C = function() {}")

	c.Test(
		"var C = class { foo() {} }",
		`/** @constructor */ var C = function() {}

C.prototype.foo = function() {}
`)

	c.Test(
		"var C = class C { }",
		`/** @constructor */
var C = function() {};
`)

	c.Test(
		"var C = class C { foo() {} }",
		`/** @constructor */
var C = function() {};
C.prototype.foo = function() {};
`)
}

// Class expressions that are the RHS of an assignment.
// Java: actually.testClassExpressionInAssignment()
func (c *Es6TranspilationIntegrationTest) TestClassExpressionInAssignment() {
	c.Test("goog.example.C = class { }", "/** @constructor */ goog.example.C = function() {}")

	c.Test(
		"goog.example.C = class { foo() {} }",
		`/** @constructor */ goog.example.C = function() {}
goog.example.C.prototype.foo = function() {};
`)
}

// Java: actually.testClassExpressionInAssignment_getElem()
func (c *Es6TranspilationIntegrationTest) TestClassExpressionInAssignment_getElem() {
	c.Test(
		"window['MediaSource'] = class {};",
		`/** @const @constructor */
var CLASS_DECL$0 = function() {
};
window["MediaSource"] = CLASS_DECL$0;
`)
}

// Java: actually.testClassExpression()
func (c *Es6TranspilationIntegrationTest) TestClassExpression() {
	c.Test(
		"var C = new (class {})();",
		`/** @const @constructor */
var CLASS_DECL$0 = function() {
};
var C = new CLASS_DECL$0();
`)
	c.Test(
		"(condition ? obj1 : obj2).prop = class C { };",
		`/** @const @constructor */
var CLASS_DECL$0 = function() {
};
(condition ? obj1 : obj2).prop = CLASS_DECL$0;
`)
}

// Java: actually.testAbstractClass()
func (c *Es6TranspilationIntegrationTest) TestAbstractClass() {
	c.EnableTypeCheck()

	c.Test(
		"/** @abstract */ class Foo {} var x = new Foo();",
		"/** @abstract @constructor */ var Foo = function() {}; var x = new Foo();",
		c.Warning(TypeCheck_INSTANTIATE_ABSTRACT_CLASS))
}

// Java: actually.testClassExpression_cannotConvert()
func (c *Es6TranspilationIntegrationTest) TestClassExpression_cannotConvert() {
	c.Test(
		"var C = new (foo || (foo = class { }))();",
		`var C = new (foo || (foo = function() {
  /** @const @constructor */
  var CLASS_DECL$0 = function() {
  };
  return CLASS_DECL$0;
}()))();
`)
}

// Java: actually.testExtends()
func (c *Es6TranspilationIntegrationTest) TestExtends() {
	c.Test(
		"class D {} class C extends D {}",
		`/** @constructor */
var D = function() {};
/** @constructor
 * @extends {D}
 */
var C = function() { D.apply(this, arguments); };
$jscomp.inherits(C, D);
`)
	AssertThat(c.GetLastCompiler().GetRuntimeJsLibManager().GetInjectedLibraries()).
		ContainsExactly("es6/util/inherits", "es6/util/construct", "es6/util/arrayfromiterable")

	c.Test(
		"class D {} class C extends D { constructor() { super(); } }",
		`/** @constructor */
var D = function() {};
/** @constructor @extends {D} */
var C = function() {
  D.call(this);
}
$jscomp.inherits(C, D);
`)

	c.Test(
		"class D {} class C extends D { constructor(str) { super(str); } }",
		`/** @constructor */
var D = function() {};
/** @constructor @extends {D} */
var C = function(str) {
  D.call(this, str);
}
$jscomp.inherits(C, D);
`)

	c.Test(
		"class C extends ns.D { }",
		`/** @constructor
 * @extends {ns.D}
 */
var C = function() {
 return ns.D.apply(this, arguments) || this;
};
$jscomp.inherits(C, ns.D);
`)
}

// Java: actually.testExtendNonNativeError()
func (c *Es6TranspilationIntegrationTest) TestExtendNonNativeError() {
	c.Test(
		`class Error {
  /** @param {string} msg */
  constructor(msg) {
    /** @const */ this.message = msg;
  }
}
class C extends Error {}
`, // autogenerated constructor
		`/** @constructor
 */
var Error = function(msg) {
  /** @const */ this.message = msg;
};
/** @constructor
 * @extends {Error}
 */
var C = function() { Error.apply(this, arguments); };
$jscomp.inherits(C, Error);
`)
	c.Test(
		`class Error {
  /** @param {string} msg */
  constructor(msg) {
    /** @const */ this.message = msg;
  }
}
class C extends Error {
  constructor() {
    super('C error'); // explicit super() call
  }
}
`,
		`/** @constructor
 */
var Error = function(msg) {
  /** @const */ this.message = msg;
};
/** @constructor
 * @extends {Error}
 */
var C = function() { Error.call(this, 'C error'); };
$jscomp.inherits(C, Error);
`)
}

// Java: actually.testExtendNativeError()
func (c *Es6TranspilationIntegrationTest) TestExtendNativeError() {
	c.Test(
		"class C extends Error {}", // autogenerated constructor
		`/** @constructor
 * @extends {Error}
 */
var C = function() {
  var $jscomp$tmp$error$m1146332801$1;
  $jscomp$tmp$error$m1146332801$1 = Error.apply(this, arguments),
      this.message = $jscomp$tmp$error$m1146332801$1.message,
      ('stack' in $jscomp$tmp$error$m1146332801$1) && (this.stack = $jscomp$tmp$error$m1146332801$1.stack),
      this;
};
$jscomp.inherits(C, Error);
`)
	c.Test(
		`class C extends Error {
  constructor() {
    var self = super('C error') || this; // explicit super() call in an expression
  }
}
`,
		`/** @constructor
 * @extends {Error}
 */
var C = function() {
  var $jscomp$tmp$error$m1146332801$1;
  var self =
      ($jscomp$tmp$error$m1146332801$1 = Error.call(this, 'C error'),
          this.message = $jscomp$tmp$error$m1146332801$1.message,
          ('stack' in $jscomp$tmp$error$m1146332801$1) && (this.stack = $jscomp$tmp$error$m1146332801$1.stack),
          this)
      || this;
};
$jscomp.inherits(C, Error);
`)
}

// Java: actually.testDynamicExtends()
func (c *Es6TranspilationIntegrationTest) TestDynamicExtends() {
	c.Test(
		"class C extends foo() {}",
		`/** @const */
var CLASS_EXTENDS$0 = foo();
/** @constructor */
var C = function() {
  return CLASS_EXTENDS$0.apply(this, arguments) || this;
};
$jscomp.inherits(C, CLASS_EXTENDS$0);
`)

	c.Test(
		"class C extends function(){} {}",
		`/** @const */
var CLASS_EXTENDS$0 = function() {
};
/** @constructor */
var C = function() {
  CLASS_EXTENDS$0.apply(this, arguments);
};
$jscomp.inherits(C, CLASS_EXTENDS$0);
`)
}

// Java: actually.testExtendsInterface()
func (c *Es6TranspilationIntegrationTest) TestExtendsInterface() {
	c.Test(
		`/** @interface */
class D {
  f() {}
}
/** @interface */
class C extends D {
  g() {}
}
`,
		`/** @interface */
var D = function() {};
D.prototype.f = function() {};
/**
 * @interface
 * @extends{D} */
var C = function() {};
$jscomp.inherits(C, D);
C.prototype.g = function() {};
`)
}

// Java: actually.testExtendsRecord()
func (c *Es6TranspilationIntegrationTest) TestExtendsRecord() {
	c.Test(
		`/** @record */
class D {
  f() {}
}
/** @record */
class C extends D {
  g() {}
}
`,
		`/** @record */
var D = function() {};
D.prototype.f = function() {};
/**
 * @record
 * @extends{D} */
var C = function() {};
$jscomp.inherits(C, D);
C.prototype.g = function() {};
`)
}

// Java: actually.testImplementsInterface()
func (c *Es6TranspilationIntegrationTest) TestImplementsInterface() {
	c.Test(
		`/** @interface */
class D {
  f() {}
}
/** @implements {D} */
class C {
  f() {console.log('hi');}
}
`,
		`/** @interface */
var D = function() {};
D.prototype.f = function() {};
/** @constructor @implements{D} */
var C = function() {};
C.prototype.f = function() {console.log('hi');};
`)
}

// Java: actually.testSuperCall()
func (c *Es6TranspilationIntegrationTest) TestSuperCall() {
	c.Test(
		"class D {} class C extends D { constructor() { super(); } }",
		`/** @constructor */
var D = function() {};
/** @constructor @extends {D} */
var C = function() {
  D.call(this);
}
$jscomp.inherits(C, D);
`)

	c.Test(
		"class D {} class C extends D { constructor(str) { super(str); } }",
		`/** @constructor */
var D = function() {}
/** @constructor @extends {D} */
var C = function(str) {
  D.call(this,str);
}
$jscomp.inherits(C, D);
`)

	c.Test(
		"class D {} class C extends D { constructor(str, n) { super(str); this.n = n; } }",
		`/** @constructor */
var D = function() {}
/** @constructor @extends {D} */
var C = function(str, n) {
  D.call(this,str);
  this.n = n;
}
$jscomp.inherits(C, D);
`)

	c.Test(
		`class D {}
class C extends D {
  constructor() { }
  foo() { return super.foo(); }
}
`,
		`/** @constructor */
var D = function() {}
/** @constructor @extends {D} */
var C = function() { }
$jscomp.inherits(C, D);
C.prototype.foo = function() {
  return D.prototype.foo.call(this);
}
`)

	c.Test(
		`class D {}
class C extends D {
  constructor() {}
  foo(bar) { return super.foo(bar); }
}
`,
		`/** @constructor */
var D = function() {}
/** @constructor @extends {D} */
var C = function() {};
$jscomp.inherits(C, D);
C.prototype.foo = function(bar) {
  return D.prototype.foo.call(this, bar);
}
`)

	c.Test(
		"class C { method() { class D extends C { constructor() { super(); }}}}",
		`/** @constructor */
var C = function() {}
C.prototype.method = function() {
  /** @constructor @extends{C} */
  var D = function() {
    C.call(this);
  }
  $jscomp.inherits(D, C);
};
`)
}

// Java: actually.testSuperKnownNotToChangeThis()
func (c *Es6TranspilationIntegrationTest) TestSuperKnownNotToChangeThis() {
	c.Test(
		`class D {
  /** @param {string} str */
  constructor(str) {
    this.str = str;
    return; // Empty return should not trigger this-changing behavior.
  }
}
class C extends D {
  /**
   * @param {string} str
   * @param {number} n
   */
  constructor(str, n) {
// This is nuts, but confirms that super() used in an expression works.
    super(str).n = n;
// Also confirm that an existing empty return is handled correctly.
    return;
  }
}
`,
		"/**\n * @constructor\n */\nvar D = function(str) {\n  this.str = str;\n  return;\n}\n/**\n * @constructor @extends {D}\n */\nvar C = function(str, n) {\n  (D.call(this,str), this).n = n; // super() returns `this`.\n  return;\n}\n$jscomp.inherits(C, D);\n")
}

// Java: actually.testSuperMightChangeThis()
func (c *Es6TranspilationIntegrationTest) TestSuperMightChangeThis() {
	// Class D is unknown, so we must assume its constructor could change `this`.
	c.Test(
		`class C extends D {
  constructor(str, n) {
// This is nuts, but confirms that super() used in an expression works.
    super(str).n = n;
// Also confirm that an existing empty return is handled correctly.
    return;
  }
}
`,
		`/** @constructor @extends {D} */
var C = function(str, n) {
  var $jscomp$super$this$m1146332801$0;
  ($jscomp$super$this$m1146332801$0 = D.call(this,str) || this).n = n;
  return $jscomp$super$this$m1146332801$0; // Duplicate because of existing return
// statement.
  return $jscomp$super$this$m1146332801$0;
}
$jscomp.inherits(C, D);
`)
}

// Java: actually.testAlternativeSuperCalls()
func (c *Es6TranspilationIntegrationTest) TestAlternativeSuperCalls() {
	c.Test(
		`class D {
  /** @param {string} name */
  constructor(name) {
    this.name = name;
  }
}
class C extends D {
  /** @param {string} str
   * @param {number} n */
  constructor(str, n) {
    if (n >= 0) {
      super('positive: ' + str);
    } else {
      super('negative: ' + str);
    }
    this.n = n;
  }
}
`,
		`/** @constructor */
var D = function(name) {
  this.name = name;
}
/** @constructor @extends {D} */
var C = function(str, n) {
  if (n >= 0) {
    D.call(this, 'positive: ' + str);
  } else {
    D.call(this, 'negative: ' + str);
  }
  this.n = n;
}
$jscomp.inherits(C, D);
`)

	// Class being extended is unknown, so we must assume super() could change the value of `this`.
	c.Test(
		`class C extends D {
  /** @param {string} str
   * @param {number} n */
  constructor(str, n) {
    if (n >= 0) {
      super('positive: ' + str);
    } else {
      super('negative: ' + str);
    }
    this.n = n;
  }
}
`,
		`/** @constructor @extends {D} */
var C = function(str, n) {
  var $jscomp$super$this$m1146332801$0;
  if (n >= 0) {
    $jscomp$super$this$m1146332801$0 = D.call(this, 'positive: ' + str) || this;
  } else {
    $jscomp$super$this$m1146332801$0 = D.call(this, 'negative: ' + str) || this;
  }
  $jscomp$super$this$m1146332801$0.n = n;
  return $jscomp$super$this$m1146332801$0;
}
$jscomp.inherits(C, D);
`)
}

// Java: actually.testComputedSuper()
func (c *Es6TranspilationIntegrationTest) TestComputedSuper() {
	c.Test(
		`class Foo {
  ['m']() { return 1; }
}

class Bar extends Foo {
  ['m']() {
    return super['m']() + 1;
  }
}
`,
		`/** @constructor */
var Foo = function() {};
Foo.prototype['m'] = function() { return 1; };
/** @constructor @extends {Foo} */
var Bar = function() { Foo.apply(this, arguments); };
$jscomp.inherits(Bar, Foo);
Bar.prototype['m'] = function () { return Foo.prototype['m'].call(this) + 1; };
`)
}

// Java: actually.testSuperMethodInGetter()
func (c *Es6TranspilationIntegrationTest) TestSuperMethodInGetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		`class Base {
  method() {
    return 5;
  }
}

class Subclass extends Base {
  constructor() {
    super();
  }

  get x() {
    return super.method();
  }
}
`,
		`/** @constructor */
var Base = function() {};
Base.prototype.method = function() { return 5; };

/** @constructor @extends {Base} */
var Subclass = function() { Base.call(this); };

$jscomp.inherits(Subclass, Base);
$jscomp.global.Object.defineProperties(Subclass.prototype, {
  x: {
    configurable:true,
    enumerable:true,
    get: function() { return Base.prototype.method.call(this); },
  }
});
`)
}

// Java: actually.testSuperMethodInSetter()
func (c *Es6TranspilationIntegrationTest) TestSuperMethodInSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		`class Base {
  method() {
    this._x = 5;
  }
}

class Subclass extends Base {
  constructor() {
    super();
  }

  set x(value) {
    super.method();
  }
}
`,
		`/** @constructor */
var Base = function() {};
Base.prototype.method = function() { this._x = 5; };

/** @constructor @extends {Base} */
var Subclass = function() { Base.call(this); };

$jscomp.inherits(Subclass, Base);
$jscomp.global.Object.defineProperties(Subclass.prototype, {
  x: {
    configurable:true,
    enumerable:true,
    set: function(value) { Base.prototype.method.call(this); },
  }
});
`)
}

// Java: actually.testExtendNativeClass()
func (c *Es6TranspilationIntegrationTest) TestExtendNativeClass() {
	c.Test(
		`class FooPromise extends Promise {
  /** @param {string} msg */
// explicit constructor
  constructor(callback, msg) {
    super(callback);
    this.msg = msg;
  }
}
`,
		`/**
 * @constructor
 * @extends {Promise}
 */
var FooPromise = function(callback, msg) {
  var $jscomp$super$this$m1146332801$0;
  $jscomp$super$this$m1146332801$0 = $jscomp.construct(Promise, [callback], this.constructor)
  $jscomp$super$this$m1146332801$0.msg = msg;
  return $jscomp$super$this$m1146332801$0;
}
$jscomp.inherits(FooPromise, Promise);

`)

	c.Test(
		// automatically generated constructor
		"class FooPromise extends Promise {}",
		`/**
 * @constructor
 * @extends {Promise}
 */
var FooPromise = function() {
  return $jscomp.construct(Promise, arguments, this.constructor)
}
$jscomp.inherits(FooPromise, Promise);
`)
}

// Java: actually.testExtendObject()
func (c *Es6TranspilationIntegrationTest) TestExtendObject() {
	// Object can be correctly extended in transpiled form, but we don't want or need to call
	// the `Object()` constructor in place of `super()`. Just replace `super()` with `this` instead.
	// Test both explicit and automatically generated constructors.
	c.Test(
		`class Foo extends Object {
  /** @param {string} msg */
  constructor(msg) {
    super();
    this.msg = msg;
  }
}
`,
		`/**
 * @constructor @extends {Object}
 */
var Foo = function(msg) {
  this; // super() replaced with its return value
  this.msg = msg;
};
$jscomp.inherits(Foo, Object);
`)
	c.Test(
		"class Foo extends Object {}",
		`/**
 * @constructor @extends {Object}
 */
var Foo = function() {
  this; // super.apply(this, arguments) replaced with its return value
};
$jscomp.inherits(Foo, Object);
`)
}

// Java: actually.testExtendNonNativeObject()
func (c *Es6TranspilationIntegrationTest) TestExtendNonNativeObject() {
	// No special handling when Object is redefined.
	customExterns :=
		c.Externs(
			NewTestExternsBuilder().
				AddAsyncIterable().
				AddArray().
				AddArguments().
				// TODO(sdh): If a normal Object extern is found, then this test fails because
				// the pass adds extra handling for super() possibly changing 'this'.  Does setting
				// the externs to "" to prevent this defeat the purpose of this test?  It became
				// necessary as a result of adding "/** @constructor */ function Object() {}" to
				// the externs used by this test.
				// .addObject()
				AddMath().
				AddExtra(
					// stubs of runtime libraries
					`/** @const */
var $jscomp = {};
$jscomp.generator = {};
$jscomp.generator.createGenerator = function() {};
/** @constructor */
$jscomp.generator.Context = function() {};
/** @constructor */
$jscomp.generator.Context.PropertyIterator = function() {};
$jscomp.asyncExecutePromiseGeneratorFunction = function(program) {};
`).
				Build())
	c.Test(
		customExterns,
		c.Srcs(
			`class Object {}
class Foo extends Object {
  /** @param {string} msg */
  constructor(msg) {
    super();
    this.msg = msg;
  }
}
`),
		c.Expected(
			`/**
 * @constructor
 */
var Object = function() {
};
/**
 * @constructor @extends {Object}
 */
var Foo = function(msg) {
  Object.call(this);
  this.msg = msg;
};
$jscomp.inherits(Foo, Object);
`))
	c.Test(
		customExterns,
		c.Srcs(
			`class Object {}
class Foo extends Object {}
`), // autogenerated constructor
		c.Expected(
			`/**
 * @constructor
 */
var Object = function() {
};
/**
 * @constructor @extends {Object}
 */
var Foo = function() {
  Object.apply(this, arguments); // all arguments passed on to super()
};
$jscomp.inherits(Foo, Object);
`))
}

// Java: actually.testMultiNameClass()
func (c *Es6TranspilationIntegrationTest) TestMultiNameClass() {
	c.Test(
		"var F = class G {}",
		`/** @constructor */
var F = function() {};
`)

	c.Test(
		"F = class G {}",
		`/** @constructor */
F = function() {};
`)
}

// Java: actually.testOutputLevelES3_compilerFeatureSetIsUpdated()
func (c *Es6TranspilationIntegrationTest) TestOutputLevelES3_compilerFeatureSetIsUpdated() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.Test(
		"class C { f() { class D {} } }",
		`/** @constructor */
var C = function() {};
C.prototype.f = function() {
  /** @constructor */
  var D = function() {}
};
`)
	// The compiler feature set gets updated to ES3.
	AssertThat(c.GetLastCompiler().GetAllowableFeatures()).IsEqualTo(parser.FeatureSet_ES3)
}

// Java: actually.testOutputLevelES3_classGettersSettersAreReported()
func (c *Es6TranspilationIntegrationTest) TestOutputLevelES3_classGettersSettersAreReported() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.TestError(
		"class C { get x() { return 1; }}",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
	c.TestError(
		"class C { set x(value) {}}", ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
}

// Java: actually.testES5FeatureTrailingCommaIsRemovedUnconditionally()
func (c *Es6TranspilationIntegrationTest) TestES5FeatureTrailingCommaIsRemovedUnconditionally() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	// trailing comma is removed
	c.Test("let obj = {a: 1, b: 2,};", "var obj = {a: 1, b: 2};")
	// also removed from the featureset
	AssertThat(c.GetLastCompiler().GetAllowableFeatures()).IsEqualTo(parser.FeatureSet_ES3)

	// also removed for ES5 output
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	// trailing comma is removed unconditionally regardless of output level
	c.Test("let obj = {a: 1, b: 2,};", "var obj = {a: 1, b: 2};")
	// also removed from the featureset
	AssertThat(c.GetLastCompiler().GetAllowableFeatures().ContainsFeature(parser.FeatureSet_Feature_TRAILING_COMMA)).IsFalse()
}

// Java: actually.testES5FeatureMultiLineStringContinuationIsRemovedUnconditionally()
func (c *Es6TranspilationIntegrationTest) TestES5FeatureMultiLineStringContinuationIsRemovedUnconditionally() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	// string continuation is removed
	c.Test(`let obj = 'a\
b';`, "var obj = 'ab';")
	// also removed from the featureset
	AssertThat(c.GetLastCompiler().GetAllowableFeatures().ContainsFeature(parser.FeatureSet_Feature_STRING_CONTINUATION)).
		IsFalse()

	// also removed for ES5 output
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	// string continuation is removed
	c.Test(`let obj = 'a\
b';`, "var obj = 'ab';")
	// also removed from the featureset
	AssertThat(c.GetLastCompiler().GetAllowableFeatures().ContainsFeature(parser.FeatureSet_Feature_STRING_CONTINUATION)).
		IsFalse()
}

// Java: actually.testClassNested()
func (c *Es6TranspilationIntegrationTest) TestClassNested() {
	c.Test(
		"class C { f() { class D {} } }",
		`/** @constructor */
var C = function() {};
C.prototype.f = function() {
  /** @constructor */
  var D = function() {}
};
`)

	c.Test(
		"class C { f() { class D extends C {} } }",
		`/** @constructor */
var C = function() {};
C.prototype.f = function() {
  /**
 * @constructor
 * @extends{C} */
  var D = function() {
    C.apply(this, arguments);
  };
  $jscomp.inherits(D, C);
};
`)
}

// Java: actually.testSuperGet()
func (c *Es6TranspilationIntegrationTest) TestSuperGet() {
	c.Test(
		"class D { d() {} } class C extends D { f() {var i = super.d;} }",
		`/** @constructor */
var D = function() {};
D.prototype.d = function() {};
/**
 * @constructor
 * @extends{D} */
var C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.prototype.f = function() {
  var i = D.prototype.d;
};
`)

	c.Test(
		"class D { ['d']() {} } class C extends D { f() {var i = super['d'];} }",
		`/** @constructor */
var D = function() {};
D.prototype['d'] = function() {};
/**
 * @constructor
 * @extends{D} */
var C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.prototype.f = function() {
  var i = D.prototype['d'];
};
`)

	c.Test(
		"class D { d() {}} class C extends D { static f() {var i = super.d;} }",
		`/** @constructor */
var D = function() {};
D.prototype.d = function() {};
/**
 * @constructor
 * @extends{D} */
var C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.f = function() {
  var i = D.d;
};
`)

	c.Test(
		"class D { ['d']() {}} class C extends D { static f() {var i = super['d'];} }",
		`/** @constructor */
var D = function() {};
D.prototype['d'] = function() {};
/**
 * @constructor
 * @extends{D} */
var C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.f = function() {
  var i = D['d'];
};
`)

	c.Test(
		"class D {} class C extends D { f() {return super.s;} }",
		`/** @constructor */
var D = function() {};
/**
 * @constructor
 * @extends{D} */
var C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.prototype.f = function() {
  return D.prototype.s;
};
`)

	c.Test(
		"class D {} class C extends D { f() { m(super.s);} }",
		`/** @constructor */
var D = function() {};
/**
 * @constructor
 * @extends{D} */
var C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.prototype.f = function() {
  m(D.prototype.s);
};
`)

	c.Test(
		"class D {} class C extends D { foo() { return super.m.foo();} }",
		`/** @constructor */
var D = function() {};
/**
 * @constructor
 * @extends{D} */
var C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.prototype.foo = function() {
  return D.prototype.m.foo();
};
`)

	c.Test(
		"class D {} class C extends D { static foo() { return super.m.foo();} }",
		`/** @constructor */
var D = function() {};
/**
 * @constructor
 * @extends{D} */
var C = function() {
  D.apply(this, arguments);
};
$jscomp.inherits(C, D);
C.foo = function() {
  return D.m.foo();
};
`)
}

// Java: actually.testSuperAccessToGettersAndSetters()
func (c *Es6TranspilationIntegrationTest) TestSuperAccessToGettersAndSetters() {
	// Getters cannot be transpiled to ES3
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		`class Base {
  get g() { return 'base'; }
  set g(v) { alert('base.prototype.g = ' + v); }
}
class Sub extends Base {
  get g() { return super.g + '-sub'; }
}
`,
		`/** @constructor */
var Base = function() {};
$jscomp.global.Object.defineProperties(Base.prototype, {
  g: {
    configurable: true,
    enumerable: true,
    get: function() {
      return 'base';
    },
    set: function(v) {
      alert('base.prototype.g = ' + v);
    }
  }
});
/** @constructor */
var Sub = function() {
  Base.apply(this, arguments);
};
$jscomp.inherits(Sub, Base);
$jscomp.global.Object.defineProperties(Sub.prototype, {
  g: {
    configurable: true,
    enumerable: true,
    get: function() {
      return Reflect.get(
                 Base.prototype, JSCompiler_renameProperty('g', Base), this) +
          '-sub';
    }
  }
});
`)

	c.TestError(
		`class Base {
  get g() { return 'base'; }
  set g(v) { alert('base.prototype.g = ' + v); }
}
class Sub extends Base {
  get g() { return super.g + '-sub'; }
  set g(v) { super.g = v + '-sub'; }
}
`,
		TranspilationUtil_CANNOT_CONVERT_YET)
}

// Java: actually.testStaticThis()
func (c *Es6TranspilationIntegrationTest) TestStaticThis() {
	c.Test(
		"class F { static f() { return this; } }",
		`/** @constructor */ var F = function() {}
/** @this {?} */ F.f = function() { return this; };
`)
}

// Java: actually.testStaticMethods()
func (c *Es6TranspilationIntegrationTest) TestStaticMethods() {
	c.Test(
		"class C { static foo() {} }",
		"/** @constructor */ var C = function() {}; C.foo = function() {};")

	c.Test(
		"class C { static foo() {}; foo() {} }",
		`/** @constructor */
var C = function() {};

C.foo = function() {};

C.prototype.foo = function() {};
`)

	c.Test(
		"class C { static foo() {}; bar() { C.foo(); } }",
		`/** @constructor */
var C = function() {};

C.foo = function() {};

C.prototype.bar = function() { C.foo(); };
`)
}

// Java: actually.testStaticInheritance()
func (c *Es6TranspilationIntegrationTest) TestStaticInheritance() {

	c.Test(
		`class D {
  static f() {}
}
class C extends D { constructor() {} }
C.f();
`,
		`/** @constructor */
var D = function() {};
D.f = function () {};
/** @constructor @extends{D} */
var C = function() {};
$jscomp.inherits(C, D);
C.f();
`)

	c.Test(
		`class D {
  static f() {}
}
class C extends D {
  constructor() {}
  f() {}
}
C.f();
`,
		`/** @constructor */
var D = function() {};
D.f = function() {};
/** @constructor @extends{D} */
var C = function() { };
$jscomp.inherits(C, D);
C.prototype.f = function() {};
C.f();
`)

	c.Test(
		`class D {
  static f() {}
}
class C extends D {
  constructor() {}
  static f() {}
  g() {}
}
`,
		`/** @constructor */
var D = function() {};
D.f = function() {};
/** @constructor @extends{D} */
var C = function() { };
$jscomp.inherits(C, D);
C.f = function() {};
C.prototype.g = function() {};
`)
}

// Java: actually.testInheritFromExterns()
func (c *Es6TranspilationIntegrationTest) TestInheritFromExterns() {
	c.Test(
		c.Externs(
			Es6TranspilationIntegrationTest_GetDefaultExternsBuilder().
				AddExtra(
					`/** @constructor */ function ExternsClass() {}
ExternsClass.m = function() {};
`).
				Build()),
		c.Srcs("class CodeClass extends ExternsClass {}"),
		c.Expected(
			`/** @constructor
 * @extends {ExternsClass}
 */
var CodeClass = function() {
  return ExternsClass.apply(this,arguments) || this;
};
$jscomp.inherits(CodeClass,ExternsClass)
`))
}

// Java: actually.testGithub752()
func (c *Es6TranspilationIntegrationTest) TestGithub752() {
	c.Test(
		"function f() { var a = b = class {};}",
		`function f() {
  /** @const @constructor */
  var CLASS_DECL$0 = function() {
  };
  var a = b = CLASS_DECL$0;
}
`)

	c.Test(
		"var ns = {}; function f() { var self = ns.Child = class {};}",
		`var ns = {};
function f() {
  /** @const @constructor */
  var CLASS_DECL$0 = function() {
  };
  var self = ns.Child = CLASS_DECL$0;
}
`)
}

// Java: actually.testInvalidClassUse()
func (c *Es6TranspilationIntegrationTest) TestInvalidClassUse() {
	externsFile :=
		NewTestExternsBuilder().
			AddArguments().
			AddFunction().
			AddJSCompLibraries().
			BuildExternsFile("externs.js")
	// Normalization changes some externs declarations,
	c.AllowExternsChanges()
	c.EnableTypeCheck()

	c.Test(
		c.Externs(externsFile),
		c.Srcs(
			`/** @constructor */
function Foo() {}
Foo.prototype.f = function() {};
class Sub extends Foo {}
(new Sub).f();
`),
		c.Expected(
			`/** @constructor */
function Foo() {}
Foo.prototype.f = function() {};
/**
 * @constructor
 * @extends {Foo}
 */
var Sub=function() { Foo.apply(this, arguments); }
$jscomp.inherits(Sub, Foo);
(new Sub).f();
`))

	c.Test(
		c.Externs(externsFile),
		c.Srcs(
			`/** @constructor @struct */
function Foo() {}
Foo.f = function() {};
class Sub extends Foo {}
Sub.f();
`),
		c.Expected(
			`/** @constructor @struct */
function Foo() {}
Foo.f = function() {};
/** @constructor
 * @extends {Foo}
 */
var Sub = function() { Foo.apply(this, arguments); };
$jscomp.inherits(Sub, Foo);
Sub.f();
`))

	c.Test(
		c.Externs(externsFile),
		c.Srcs(
			`/** @constructor */
function Foo() {}
Foo.f = function() {};
class Sub extends Foo {}
`),
		c.Expected(
			`/** @constructor */
function Foo() {}
Foo.f = function() {};
/** @constructor
 * @extends {Foo}
 */
var Sub = function() { Foo.apply(this, arguments); };
$jscomp.inherits(Sub, Foo);
`))
}

// Getters and setters are supported, both in object literals and in classes, but only if the
// output language is ES5.
// Java: actually.testEs5GettersAndSettersClasses()
func (c *Es6TranspilationIntegrationTest) TestEs5GettersAndSettersClasses() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		"class C { get value() { return 0; } }",
		`/** @constructor */
var C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    get: function() {
      return 0;
    }
  }
});
`)

	c.Test(
		"class C { set value(val) { this.internalVal = val; } }",
		`/** @constructor */
var C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    set: function(val) {
      this.internalVal = val;
    }
  }
});
`)

	c.Test(
		`class C {
  set value(val) {
    this.internalVal = val;
  }
  get value() {
    return this.internalVal;
  }
}
`,
		`/** @constructor */
var C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    set: function(val) {
      this.internalVal = val;
    },
    get: function() {
      return this.internalVal;
    }
  }
});
`)

	c.Test(
		`class C {
  get alwaysTwo() {
    return 2;
  }

  get alwaysThree() {
    return 3;
  }
}
`,
		`/** @constructor */
var C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  alwaysTwo: {
    configurable: true,
    enumerable: true,
    get: function() {
      return 2;
    }
  },
  alwaysThree: {
    configurable: true,
    enumerable: true,
    get: function() {
      return 3;
    }
  },
});
`)
}

// Java: actually.testEs5GettersAndSettersOnClassesWithClassSideInheritance()
func (c *Es6TranspilationIntegrationTest) TestEs5GettersAndSettersOnClassesWithClassSideInheritance() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { static get value() {} }  class D extends C { static get value() {} }",
		`/** @constructor */
var C = function() {};
/** @nocollapse */
C.value;
$jscomp.global.Object.defineProperties(C, {
  value: {
    configurable: true,
    enumerable: true,
    get: function() {}
  }
});
/** @constructor
 * @extends {C}
 */
var D = function() {
  C.apply(this,arguments);
};
/** @nocollapse */
D.value;
$jscomp.inherits(D, C);
$jscomp.global.Object.defineProperties(D, {
  value: {
    configurable: true,
    enumerable: true,
    get: function() {}
  }
});
`)
}

// Check that the types from the getter/setter are copied to the declaration on the prototype.
// Java: actually.testEs5GettersAndSettersClassesWithTypes()
func (c *Es6TranspilationIntegrationTest) TestEs5GettersAndSettersClassesWithTypes() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		"class C { /** @return {number} */ get value() { return 0; } }",
		`/** @constructor */
var C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    /**
     * @return {number}
     */
    get: function() {
      return 0;
    }
  }
});
`)

	c.Test(
		"class C { /** @param {string} v */ set value(v) { } }",
		`/** @constructor */
var C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    /**
     * @param {string} v
     */
    set: function(v) {}
  }
});
`)
}

// Java: actually.testClassEs5GetterSetterIncorrectTypes()
func (c *Es6TranspilationIntegrationTest) TestClassEs5GetterSetterIncorrectTypes() {
	// TODO(b/144721663): these should cause some sort of type warning
	c.EnableTypeCheck()
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	// Using @type instead of @return on a getter.
	externsFile :=
		NewTestExternsBuilder().
			AddArguments().
			AddFunction().
			AddJSCompLibraries().
			AddObject().
			BuildExternsFile("externs.js")
	// Normalization changes some externs declarations,
	c.AllowExternsChanges()

	c.Test(
		c.Externs(externsFile),
		c.Srcs("class C { /** @type {string} */ get value() { } }"),
		c.Expected(
			`/** @constructor */
var C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    /** @type {string} */
    get: function() {}
  }
});
`))

	// Using @type instead of @param on a setter.
	c.Test(
		c.Externs(externsFile),
		c.Srcs("class C { /** @type {string} */ set value(v) { } }"),
		c.Expected(
			`/** @constructor */
var C = function() {};
$jscomp.global.Object.defineProperties(C.prototype, {
  value: {
    configurable: true,
    enumerable: true,
    /** @type {string} */
    set: function(v) {}
  }
});
`))
}

// @bug 20536614
// Java: actually.testStaticGetterSetter()
func (c *Es6TranspilationIntegrationTest) TestStaticGetterSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.Test(
		"class C { static get foo() {} }",
		`/** @constructor */
var C = function() {};
/** @nocollapse */
C.foo;
$jscomp.global.Object.defineProperties(C, {
  foo: {
    configurable: true,
    enumerable: true,
    get: function() {}
  }
})
`)

	c.Test(
		`class C { static get foo() {} }
class Sub extends C {}
`,
		`/** @constructor */
var C = function() {};
/** @nocollapse */
C.foo;
$jscomp.global.Object.defineProperties(C, {
  foo: {
    configurable: true,
    enumerable: true,
    get: function() {}
  }
})

/** @constructor
 * @extends {C}
 */
var Sub = function() {
  C.apply(this, arguments);
};
$jscomp.inherits(Sub, C)
`)
}

// Java: actually.testStaticSetter()
func (c *Es6TranspilationIntegrationTest) TestStaticSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"class C { static set foo(x) {} }",
		`/** @constructor */
var C = function() {};
/** @nocollapse */
C.foo;
$jscomp.global.Object.defineProperties(C, {
  foo: {
    configurable: true,
    enumerable: true,
    set: function(x) {}
  }
});
`)
}

// Java: actually.testInitSymbol()
func (c *Es6TranspilationIntegrationTest) TestInitSymbol() {
	// Include the extern for Symbol, because we're testing renaming of local variables with the
	// same name.
	externsFileWithSymbol :=
		NewTestExternsBuilder().
			AddIterable(). // there's no method for just Symbol, but this will get it
			BuildExternsFile("externs.js")

	c.Test(
		c.Externs(externsFileWithSymbol), //
		c.Srcs("let a = alert(Symbol.thimble);"),
		c.Expected("var a = alert(Symbol.thimble)"))
	AssertThat(c.GetLastCompiler().GetRuntimeJsLibManager().GetInjectedLibraries()).
		ContainsExactly("es6/symbol")

	c.Test(
		c.Externs(externsFileWithSymbol), //
		c.Srcs("let a = alert(Symbol.iterator);"),
		c.Expected("var a = alert(Symbol.iterator)"))
	AssertThat(c.GetLastCompiler().GetRuntimeJsLibManager().GetInjectedLibraries()).
		ContainsExactly("es6/symbol")

	c.Test(
		c.Externs(externsFileWithSymbol),
		c.Srcs(
			`function f() {
  let x = 1;
  let y = Symbol('nimble');
}
`),
		c.Expected(
			`function f() {
  var x = 1;
  var y = Symbol('nimble');
}
`))
	externs := c.Externs(externsFileWithSymbol)
	srcs :=
		c.Srcs(
			`function f() {
  if (true) {
     let Symbol = function() {};
  }
// This Symbol is the global one
  alert(Symbol.ism)
}
`)
	expected :=
		c.Expected(
			`function f() {
  if (true) {
// normalization renames the local Symbol to be different from the global Symbol
     var Symbol$jscomp$0 = function() {};
  }
  alert(Symbol.ism)
}
`)
	c.Test(externs, srcs, expected)

	externs = c.Externs(externsFileWithSymbol)
	srcs =
		c.Srcs(
			`function f() {
  if (true) {
    let Symbol = function() {};
// This Symbol is the local definition. There's no use of the global Symbol in
// this function.
    alert(Symbol.ism)
  }
}
`)
	expected =
		c.Expected(
			`function f() {
  if (true) {
// The local definition of Symbol doesn't have to be renamed, because there's
// no usage of the global Symbol to conflict with it.
    var Symbol = function() {};
    alert(Symbol.ism)
  }
}
`)
	c.Test(externs, srcs, expected)
	// No $jscomp.initSymbol in externs
	c.TestSame(c.Externs("alert(Symbol.thimble);"), c.Srcs(""))
}

// Java: actually.testInitSymbolIterator()
func (c *Es6TranspilationIntegrationTest) TestInitSymbolIterator() {
	c.Test(
		"var x = {[Symbol.iterator]: function() { return this; }};",
		`var $jscomp$compprop0 = {};
var x = ($jscomp$compprop0[Symbol.iterator] = function() {return this;},
         $jscomp$compprop0)
`)
}

// ES5 getters and setters should report an error if the languageOut is ES3.
// Java: actually.testEs5GettersAndSetters_es3()
func (c *Es6TranspilationIntegrationTest) TestEs5GettersAndSetters_es3() {
	c.TestError(
		"let x = { get y() {} };", ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
	c.TestError(
		"let x = { set y(value) {} };",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
}

// ES5 getters and setters on object literals should be left alone if the languageOut is ES5.
// Java: actually.testEs5GettersAndSettersObjLit_es5()
func (c *Es6TranspilationIntegrationTest) TestEs5GettersAndSettersObjLit_es5() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.TestSame("var x = { get y() {} };")
	c.TestSame("var x = { set y(value) {} };")
}

// Java: actually.testForOfLoop()
func (c *Es6TranspilationIntegrationTest) TestForOfLoop() {
	// Iteration var shadows an outer var ()
	c.Test(
		"var i = 'outer'; for (let i of [1, 2, 3]) { alert(i); } alert(i);",
		"var i = 'outer';\nvar $jscomp$iter$0 = (0, $jscomp.makeIterator)([1, 2, 3]);\n// Normalize runs before for-of rewriting. Therefore, first Normalize renames the\n// `let i` to `let i$jscomp$1` to avoid conficting it with outer `i`. Then, the\n// for-of rewriting prepends the unique ID `$jscomp$key$m123..456$0` to its declared\n// name as it does to all for-of loop keys.\nvar KEY$0$i$jscomp$1 = $jscomp$iter$0.next();\nfor (; !KEY$0$i$jscomp$1.done; KEY$0$i$jscomp$1 = $jscomp$iter$0.next()) {\n  var i$jscomp$1 = KEY$0$i$jscomp$1.value;\n  {\n    alert(i$jscomp$1);\n  }\n}\nalert(i);\n")
}

// Java: actually.testForOfWithConstInitiliazer()
func (c *Es6TranspilationIntegrationTest) TestForOfWithConstInitiliazer() {
	c.EnableNormalize()

	c.Test(
		"for(const i of [1,2]) {i;}",
		"var $jscomp$iter$0 = (0, $jscomp.makeIterator)([1, 2]);\n// Normalize runs before for-of rewriting. Normalize does not rename the `const i`\n// if there is no other conflicting `i` declaration. Then, the  for-of rewriting\n// prepends `$jscomp$key$m123..456$0` to its declared name as it does to all for-of\n// loop keys.\nvar KEY$0$i = $jscomp$iter$0.next();\nfor (; !KEY$0$i.done; KEY$0$i = $jscomp$iter$0.next()) {\n  /** @const */\n  var i = KEY$0$i.value; // marked as const name\n  {\n    i; // marked as const name\n  }\n}\n")
}

// Java: actually.testMultipleForOfWithSameInitializerName()
func (c *Es6TranspilationIntegrationTest) TestMultipleForOfWithSameInitializerName() {
	c.EnableNormalize()
	c.Test(
		`function* inorder1(t) {
    for (var x of []) {
      yield x;
    }
    for (var x of []) {
      yield x;
    }
}
`,
		`function inorder1(t) {
  var x;
  var $jscomp$iter$0;
  var KEY$0$x; // key for first for-of loop
  var $jscomp$iter$1;
  var KEY$1$x;
  return $jscomp.generator.createGenerator(inorder1, function($jscomp$generator$context$m1146332801$2) {
    switch($jscomp$generator$context$m1146332801$2.getNextAddressJsc()) {
      case 1:
        $jscomp$iter$0 = (0, $jscomp.makeIterator)([]);
        KEY$0$x = $jscomp$iter$0.next();
      case 2:
        if (!!KEY$0$x.done) {
          $jscomp$generator$context$m1146332801$2.jumpTo(4);
          break;
        }
        x = KEY$0$x.value;
        return $jscomp$generator$context$m1146332801$2.yield(x, 3);
      case 3:
        KEY$0$x = $jscomp$iter$0.next();
        $jscomp$generator$context$m1146332801$2.jumpTo(2);
        break;
      case 4:
        $jscomp$iter$1 = (0, $jscomp.makeIterator)([]);
        KEY$1$x = $jscomp$iter$1.next();
      case 6:
        if (!!KEY$1$x.done) {
          $jscomp$generator$context$m1146332801$2.jumpTo(0);
          break;
        }
        x = KEY$1$x.value;
        return $jscomp$generator$context$m1146332801$2.yield(x, 7);
      case 7:
        KEY$1$x = $jscomp$iter$1.next();
        $jscomp$generator$context$m1146332801$2.jumpTo(6);
        break;
    }
  });}
`)
}

// Java: actually.testForOfRedeclaredVar()
func (c *Es6TranspilationIntegrationTest) TestForOfRedeclaredVar() {
	c.Test(
		`for (let x of []) {
  let x = 0;
}
`,
		`var $jscomp$iter$0 = (0, $jscomp.makeIterator)([]);
var KEY$0$x = $jscomp$iter$0.next();
for (; !KEY$0$x.done; KEY$0$x = $jscomp$iter$0.next()) {
  var x = KEY$0$x.value;
  {
    var x$jscomp$1 = 0;
  }
}
`)
}

// Java: actually.testArgumentsEscaped()
func (c *Es6TranspilationIntegrationTest) TestArgumentsEscaped() {
	c.TestSame(
		`function f() {
  return g(arguments);
}
`)
}

// Java: actually.testMethodInObject()
func (c *Es6TranspilationIntegrationTest) TestMethodInObject() {
	c.Test(
		"var obj = {           f() {alert(1); } };", //
		"var obj = { f: function() {alert(1); } };")

	c.Test(
		"var obj = {           f() { alert(1); },    x };", //
		"var obj = { f: function() { alert(1); }, x: x };")
}

// Java: actually.testComputedPropertiesWithMethod()
func (c *Es6TranspilationIntegrationTest) TestComputedPropertiesWithMethod() {
	c.Test(
		"var obj = { ['f' + 1]: 1, m() {}, ['g' + 1]: 1, };",
		`var $jscomp$compprop0 = {};
var obj = ($jscomp$compprop0['f' + 1] = 1,
  ($jscomp$compprop0.m = function() {},
     ($jscomp$compprop0['g' + 1] = 1, $jscomp$compprop0)));
`)
}

// Java: actually.testComputedProperties()
func (c *Es6TranspilationIntegrationTest) TestComputedProperties() {
	c.Test(
		"var obj = { ['f' + 1] : 1, ['g' + 1] : 1 };",
		`var $jscomp$compprop0 = {};
var obj = ($jscomp$compprop0['f' + 1] = 1,
  ($jscomp$compprop0['g' + 1] = 1, $jscomp$compprop0));
`)

	c.Test(
		"var obj = { ['f'] : 1};",
		`var $jscomp$compprop0 = {};
var obj = ($jscomp$compprop0['f'] = 1,
  $jscomp$compprop0);
`)

	c.Test(
		"var o = { ['f'] : 1}; var p = { ['g'] : 1};",
		`var $jscomp$compprop0 = {};
var o = ($jscomp$compprop0['f'] = 1,
  $jscomp$compprop0);
var $jscomp$compprop1 = {};
var p = ($jscomp$compprop1['g'] = 1,
  $jscomp$compprop1);
`)

	c.Test(
		"({['f' + 1] : 1})",
		`var $jscomp$compprop0 = {};
($jscomp$compprop0['f' + 1] = 1,
  $jscomp$compprop0)
`)

	c.Test(
		"({'a' : 2, ['f' + 1] : 1})",
		`var $jscomp$compprop0 = {};
($jscomp$compprop0['a'] = 2,
  ($jscomp$compprop0['f' + 1] = 1, $jscomp$compprop0));
`)

	c.Test(
		"({['f' + 1] : 1, 'a' : 2})",
		`var $jscomp$compprop0 = {};
($jscomp$compprop0['f' + 1] = 1,
  ($jscomp$compprop0['a'] = 2, $jscomp$compprop0));
`)

	c.Test(
		"({'a' : 1, ['f' + 1] : 1, 'b' : 1})", //
		`var $jscomp$compprop0 = {};
($jscomp$compprop0['a'] = 1,
  ($jscomp$compprop0['f' + 1] = 1,
    ($jscomp$compprop0['b'] = 1,
      $jscomp$compprop0)));
`)

	c.Test(
		"({'a' : x++, ['f' + x++] : 1, 'b' : x++})",
		`var $jscomp$compprop0 = {};
($jscomp$compprop0['a'] = x++, ($jscomp$compprop0['f' + x++] = 1,
  ($jscomp$compprop0['b'] = x++, $jscomp$compprop0)))
`)

	c.Test(
		"({a : x++, ['f' + x++] : 1, b : x++})",
		`var $jscomp$compprop0 = {};
($jscomp$compprop0.a = x++, ($jscomp$compprop0['f' + x++] = 1,
  ($jscomp$compprop0.b = x++, $jscomp$compprop0)))
`)

	c.Test(
		"({a, ['f' + 1] : 1})",
		`var $jscomp$compprop0 = {};
  ($jscomp$compprop0.a = a, ($jscomp$compprop0['f' + 1] = 1, $jscomp$compprop0))
`)

	c.Test(
		"({['f' + 1] : 1, a})",
		`var $jscomp$compprop0 = {};
  ($jscomp$compprop0['f' + 1] = 1, ($jscomp$compprop0.a = a, $jscomp$compprop0))
`)

	c.Test(
		"var obj = { [foo]() {}}",
		`var $jscomp$compprop0 = {};
var obj = ($jscomp$compprop0[foo] = function(){}, $jscomp$compprop0)
`)
}

// Java: actually.testComputedPropGetterSetter()
func (c *Es6TranspilationIntegrationTest) TestComputedPropGetterSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.TestSame("var obj = {get latest () {return undefined;}}")
	c.TestSame("var obj = {set latest (str) {}}")
	c.Test(
		"var obj = {'a' : 2, get l () {return null;}, ['f' + 1] : 1}",
		`var $jscomp$compprop0 = {get l () {return null;}};
var obj = ($jscomp$compprop0['a'] = 2,
  ($jscomp$compprop0['f' + 1] = 1, $jscomp$compprop0));
`)
	c.Test(
		"var obj = {['a' + 'b'] : 2, set l (str) {}}",
		`var $jscomp$compprop0 = {set l (str) {}};
var obj = ($jscomp$compprop0['a' + 'b'] = 2, $jscomp$compprop0);
`)
}

// Java: actually.testComputedPropClass()
func (c *Es6TranspilationIntegrationTest) TestComputedPropClass() {
	c.Test(
		"class C { [foo]() { alert(1); } }",
		`var COMP_FIELD$0 = foo;
/** @constructor */
var C = function() {
};
C.prototype[COMP_FIELD$0] = function() {
  alert(1);
};
`)

	c.Test(
		"class C { static [foo]() { alert(2); } }",
		`var COMP_FIELD$0 = foo;
/** @constructor */
var C = function() {
};
C[COMP_FIELD$0] = function() {
  alert(2);
};
`)
}

// Java: actually.testComputedPropCannotConvert()
func (c *Es6TranspilationIntegrationTest) TestComputedPropCannotConvert() {
	c.TestError(
		"var o = { get [foo]() {}}", ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
	c.TestError(
		"var o = { set [foo](val) {}}",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
}

// Java: actually.testNoComputedProperties()
func (c *Es6TranspilationIntegrationTest) TestNoComputedProperties() {
	c.TestSame("({'a' : 1})")
	c.TestSame("({'a' : 1, f : 1, b : 1})")
}

// Java: actually.testUntaggedTemplateLiteral()
func (c *Es6TranspilationIntegrationTest) TestUntaggedTemplateLiteral() {
	c.Test("``", "''")
	c.Test("`\"`", "'\\\"'")
	c.Test("`'`", "\"'\"")
	c.Test("`\\``", "'`'")
	c.Test("`\\\"`", "'\\\"'")
	c.Test("`\\\\\"`", "'\\\\\\\"'")
	c.Test("`\"\\\\`", "'\"\\\\'")
	c.Test("`$$`", "'$$'")
	c.Test("`$$$`", "'$$$'")
	c.Test("`\\$$$`", "'$$$'")
	c.Test("`hello`", "'hello'")
	c.Test("`hello\nworld`", "'hello\\nworld'")
	c.Test("`hello\rworld`", "'hello\\nworld'")
	c.Test("`hello\r\nworld`", "'hello\\nworld'")
	c.Test("`hello\n\nworld`", "'hello\\n\\nworld'")
	c.Test("`hello\\r\\nworld`", "'hello\\r\\nworld'")
	c.Test("`${world}`", "'' + world")
	c.Test("`hello ${world}`", "'hello ' + world")
	c.Test("`${hello} world`", "hello + ' world'")
	c.Test("`${hello}${world}`", "'' + hello + world")
	c.Test("`${a} b ${c} d ${e}`", "a + ' b ' + c + ' d ' + e")
	c.Test("`hello ${a + b}`", "'hello ' + (a + b)")
	c.Test("`hello ${a, b, c}`", "'hello ' + (a, b, c)")
	c.Test("`hello ${a ? b : c}${a * b}`", "'hello ' + (a ? b : c) + (a * b)")
}

// Java: actually.testUnicodeEscapes()
func (c *Es6TranspilationIntegrationTest) TestUnicodeEscapes() {
	c.Test(
		"var \\u{73} = '\\u{2603}'", // ☃
		"var       s = '☃'")
	c.Test(
		"var \\u{63} = '\\u{1f42a}'", // 🐪
		"var       c = '🐪'")
	c.Test(
		"var str = `begin\\u{2026}end`", //
		"var str = 'begin\\u2026end'")
}

// Java: actually.testObjectLiteralShorthand()
func (c *Es6TranspilationIntegrationTest) TestObjectLiteralShorthand() {
	c.RewriteUniqueIdAndTest(
		c.Srcs(
			`function f() {
  var x = 1;
  if (a) {
    let x = 2;
    return {x};
  }
  return x;
}
`),
		c.Expected(
			`function f() {
  var x = 1;
  if (a) {
    var x$jscomp$0 = 2;
    return {x: x$jscomp$0};
  }
  return x;
}
`))

	c.RewriteUniqueIdAndTest(
		c.Srcs(
			`function f(a) {
  var {x} = a;
  if (a) {
    let x = 2;
    return x;
  }
  return x;
}
`),
		c.Expected(
			`function f(a) {
  var x;
  var $jscomp$destructuring$var0 = a;
  x = $jscomp$destructuring$var0.x;
  if (a) {
    var x$jscomp$0 = 2;
    return x$jscomp$0;
  }
  return x;
}
`))

	// Note: if the inner `let` declaration is defined as a destructuring assignment
	// then the test would fail because Es6RewriteBlockScopeDeclaration does not even
	// look at destructuring declarations, expecting them to already have been
	// rewritten, and this test does not include that pass.
}

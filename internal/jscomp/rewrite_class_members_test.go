// Port of com/google/javascript/jscomp/RewriteClassMembersTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RewriteClassMembersTest.java)

/*
 * Copyright 2021 The Closure Compiler Authors.
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

type RewriteClassMembersTest struct {
	CompilerTestCaseBase
}

// Java: RewriteClassMembersTest()
func NewRewriteClassMembersTest() *RewriteClassMembersTest {
	c := &RewriteClassMembersTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteClassMembersTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteClassMembersTest() })
}

// Java: RewriteClassMembersTest.setUp()
func (c *RewriteClassMembersTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
}

// Java: RewriteClassMembersTest.getProcessor(Compiler)
func (c *RewriteClassMembersTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs6NormalizeClasses(compiler)
}

// Java: RewriteClassMembersTest.withOptions()
func (c *RewriteClassMembersTest) WithOptions() *RewriteClassMembersTest_Options {
	return NewRewriteClassMembersTest_Options(true, CompilerOptions_LanguageMode_ECMASCRIPT_2021)
}

// RewriteClassMembersTest_Options is the record
// Options(assumeStaticInheritanceIsNotUsed, languageOut).
type RewriteClassMembersTest_Options struct {
	assumeStaticInheritanceIsNotUsed bool
	languageOut                      CompilerOptions_LanguageMode
}

// Java: RewriteClassMembersTest.Options(boolean, LanguageMode)
func NewRewriteClassMembersTest_Options(assumeStaticInheritanceIsNotUsed bool, languageOut CompilerOptions_LanguageMode) *RewriteClassMembersTest_Options {
	return &RewriteClassMembersTest_Options{
		assumeStaticInheritanceIsNotUsed: assumeStaticInheritanceIsNotUsed,
		languageOut:                      languageOut,
	}
}

// Java: RewriteClassMembersTest.Options.assumeStaticInheritanceIsNotUsed()
func (o *RewriteClassMembersTest_Options) AssumeStaticInheritanceIsNotUsed() bool {
	return o.assumeStaticInheritanceIsNotUsed
}

// Java: RewriteClassMembersTest.Options.languageOut()
func (o *RewriteClassMembersTest_Options) LanguageOut() CompilerOptions_LanguageMode {
	return o.languageOut
}

// Java: RewriteClassMembersTest.Options.useStaticInheritance()
func (o *RewriteClassMembersTest_Options) UseStaticInheritance() *RewriteClassMembersTest_Options {
	return NewRewriteClassMembersTest_Options(false, o.LanguageOut())
}

// Java: RewriteClassMembersTest.Options.useEs2022LanguageOut()
func (o *RewriteClassMembersTest_Options) UseEs2022LanguageOut() *RewriteClassMembersTest_Options {
	return NewRewriteClassMembersTest_Options(o.AssumeStaticInheritanceIsNotUsed(), CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
}

// Java: RewriteClassMembersTest.Options.useEs5LanguageOut()
func (o *RewriteClassMembersTest_Options) UseEs5LanguageOut() *RewriteClassMembersTest_Options {
	return NewRewriteClassMembersTest_Options(o.AssumeStaticInheritanceIsNotUsed(), CompilerOptions_LanguageMode_ECMASCRIPT5)
}

// Java: RewriteClassMembersTest.test(String, String)
// Java: RewriteClassMembersTest.test(Options, String, String)
func (c *RewriteClassMembersTest) Test(args ...any) {
	if compilerTestCase_argsAre(args, "string", "string") {
		input, expected := args[0].(string), args[1].(string)
		c.Test(c.WithOptions(), input, expected)
		return
	}
	if len(args) == 3 {
		if options, ok := args[0].(*RewriteClassMembersTest_Options); ok {
			input, expected := args[1].(string), args[2].(string)
			c.SetAssumeStaticInheritanceIsNotUsed(options.AssumeStaticInheritanceIsNotUsed())
			c.SetLanguageOut(options.LanguageOut())

			c.CompilerTestCaseBase.Test(input, expected)
			return
		}
	}
	// Go-only: language-adapter — the test overloads this class does not override.
	c.CompilerTestCaseBase.Test(args...)
}

// Java: RewriteClassMembersTest.testError(String, DiagnosticType)
// Java: RewriteClassMembersTest.testError(Options, String, DiagnosticType)
func (c *RewriteClassMembersTest) TestError(args ...any) {
	if len(args) == 2 {
		if input, ok := args[0].(string); ok {
			if error_, ok := args[1].(*DiagnosticType); ok {
				c.TestError(c.WithOptions(), input, error_)
				return
			}
		}
	}
	if len(args) == 3 {
		if options, ok := args[0].(*RewriteClassMembersTest_Options); ok {
			input, error_ := args[1].(string), args[2].(*DiagnosticType)
			c.SetAssumeStaticInheritanceIsNotUsed(options.AssumeStaticInheritanceIsNotUsed())
			c.SetLanguageOut(options.LanguageOut())

			c.CompilerTestCaseBase.TestError(input, error_)
			return
		}
	}
	// Go-only: language-adapter — the testError overloads this class does not override.
	c.CompilerTestCaseBase.TestError(args...)
}

// Java: RewriteClassMembersTest.testSame(String)
// Java: RewriteClassMembersTest.testSame(Options, String)
func (c *RewriteClassMembersTest) TestSame(args ...any) {
	if compilerTestCase_argsAre(args, "string") {
		src := args[0].(string)
		c.Test(c.WithOptions(), src, src)
		return
	}
	if len(args) == 2 {
		if options, ok := args[0].(*RewriteClassMembersTest_Options); ok {
			src := args[1].(string)
			c.Test(options, src, src)
			return
		}
	}
	// Go-only: language-adapter — the testSame overloads this class does not override.
	c.CompilerTestCaseBase.TestSame(args...)
}

// Java: RewriteClassMembersTest.testClassStaticBlock_superRef()
func (c *RewriteClassMembersTest) TestClassStaticBlock_superRef() {
	c.Test(
		`class B {
  static y = 3;
}
class C extends B {
  static {
    let x = super.y;
  }
}
`,
		`class B {}
B.y = 3;
class C extends B {}
{
  // TODO (tflo): Reflect.get(B, 'y') is the technically correct way.
  let x = B.y;
}
`)
}

// Java: RewriteClassMembersTest.testPrivateField()
func (c *RewriteClassMembersTest) TestPrivateField() {
	c.Test(
		`class Foo {
  #field;
  #field_initialized = 1;
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.field = void 0;
    PRIVATE$1.field_initialized = 1;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateFieldWithIdentifierInitializer()
func (c *RewriteClassMembersTest) TestPrivateFieldWithIdentifierInitializer() {
	c.Test(
		`const bar = {};
class Foo {
  #field = bar;
}
`,
		`const bar = {};
const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.field = bar;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateField_interleavedWithPublicFields()
func (c *RewriteClassMembersTest) TestPrivateField_interleavedWithPublicFields() {
	c.Test(
		`class Foo {
  f1 = 1;
  #pf2 = this.f1 + 2;
  f3 = this.#pf2 + 3;
  #pf4 = this.f3 + this.#pf2;
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    this.f1 = 1;
    PRIVATE$1.pf2 = this.f1 + 2;
    this.f3 = PRIVATE_MAP$0.get(this).pf2 + 3;
    PRIVATE$1.pf4 = this.f3 + PRIVATE_MAP$0.get(this).pf2;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticField()
func (c *RewriteClassMembersTest) TestPrivateStaticField() {
	c.Test(
		`class Foo {
  static #static_field;
  static #static_field_initialized = 2;
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.static_field = void 0;
    PRIVATE$1.static_field_initialized = 2;
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateAndStaticFieldMultipleClasses()
func (c *RewriteClassMembersTest) TestPrivateAndStaticFieldMultipleClasses() {
	c.Test(
		`class Foo {
  #field;
  #field_initialized = 1;
  static #static_field;
  static #static_field_initialized = 2;
}
class Bar {
  #field;
  #field_initialized = 1;
  static #static_field;
  static #static_field_initialized = 2;
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const STATIC_PRIVATE_MAP$1 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    PRIVATE$2.field = void 0;
    PRIVATE$2.field_initialized = 1;
  }
  static STATIC_INIT$4() {
    const PRIVATE$3 = Object.create(null);
    STATIC_PRIVATE_MAP$1.set(Foo, PRIVATE$3);
    PRIVATE$3.static_field = void 0;
    PRIVATE$3.static_field_initialized = 2;
  }
}
Foo.STATIC_INIT$4();
const PRIVATE_MAP$5 = new $jscomp.PrivateMap();
const STATIC_PRIVATE_MAP$6 = new $jscomp.PrivateMap();
class Bar {
  constructor() {
    const PRIVATE$7 = Object.create(null);
    PRIVATE_MAP$5.set(this, PRIVATE$7);
    PRIVATE$7.field = void 0;
    PRIVATE$7.field_initialized = 1;
  }
  static STATIC_INIT$9() {
    const PRIVATE$8 = Object.create(null);
    STATIC_PRIVATE_MAP$6.set(Bar, PRIVATE$8);
    PRIVATE$8.static_field = void 0;
    PRIVATE$8.static_field_initialized = 2;
  }
}
Bar.STATIC_INIT$9();
`)
}

// Java: RewriteClassMembersTest.testPrivateMethod()
func (c *RewriteClassMembersTest) TestPrivateMethod() {
	c.Test(
		`class Foo {
  #method() {}
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticMethod()
func (c *RewriteClassMembersTest) TestPrivateStaticMethod() {
	c.Test(
		`class Foo {
  static #staticMethod() { return 42; }
  static callStatic() {
    return Foo.#staticMethod();
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static callStatic() {
    return STATIC_PRIVATE_MAP$0.get(Foo).staticMethod.call(Foo);
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.staticMethod = function() {
      return 42;
    };
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticMethodCallThis()
func (c *RewriteClassMembersTest) TestPrivateStaticMethodCallThis() {
	c.Test(
		`class Foo {
  static #staticMethod() { return 42; }
  static #helper() { return this.#staticMethod(); }
  static callStatic() {
    return this.#helper();
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static callStatic() {
    return STATIC_PRIVATE_MAP$0.get(this).helper.call(this);
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.staticMethod = function() {
      return 42;
    };
    PRIVATE$1.helper = function() {
      return STATIC_PRIVATE_MAP$0.get(this).staticMethod.call(this);
    };
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticMethodAndFieldOrder()
func (c *RewriteClassMembersTest) TestPrivateStaticMethodAndFieldOrder() {
	c.Test(
		`class C {
  static #a = 1;
  static #b() { return this.#a + 1; }
  static #c = this.#b() + 1;
  static getC() { return this.#c; }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class C {
  static getC() {
    return STATIC_PRIVATE_MAP$0.get(this).c;
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(C, PRIVATE$1);
    PRIVATE$1.b = function() {
      return STATIC_PRIVATE_MAP$0.get(this).a + 1;
    };
    PRIVATE$1.a = 1;
    PRIVATE$1.c = STATIC_PRIVATE_MAP$0.get(C).b.call(C) + 1;
  }
}
C.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateAndStaticMethod()
func (c *RewriteClassMembersTest) TestPrivateAndStaticMethod() {
	c.Test(
		`class Foo {
  #instanceMethod() { return 1; }
  static #staticMethod() { return 2; }
  callInstance() { return this.#instanceMethod(); }
  static callStatic() { return Foo.#staticMethod(); }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$2 = Object.create(null, {
  instanceMethod: {
    value: function() {
      return 1;
    }
  }
});
const STATIC_PRIVATE_MAP$1 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$3 = Object.create(PRIVATE_PROTO$2);
    PRIVATE_MAP$0.set(this, PRIVATE$3);
  }
  callInstance() {
    return PRIVATE_MAP$0.get(this).instanceMethod.call(this);
  }
  static callStatic() {
    return STATIC_PRIVATE_MAP$1.get(Foo).staticMethod.call(Foo);
  }
  static STATIC_INIT$5() {
    const PRIVATE$4 = Object.create(null);
    STATIC_PRIVATE_MAP$1.set(Foo, PRIVATE$4);
    PRIVATE$4.staticMethod = function() {
      return 2;
    };
  }
}
Foo.STATIC_INIT$5();
`)
}

// Java: RewriteClassMembersTest.testPrivateGetter()
func (c *RewriteClassMembersTest) TestPrivateGetter() {
	c.Test(
		`class Foo {
  get #prop() { return 3; }
  getValue() {
    return this.#prop;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  prop: {
    get: function() {
      return 3;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  getValue() {
    return PRIVATE_MAP$0.get(this).prop;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticGetter()
func (c *RewriteClassMembersTest) TestPrivateStaticGetter() {
	c.Test(
		`class Foo {
  static get #prop() { return 4; }
  static getValue() {
    return Foo.#prop;
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static getValue() {
    return STATIC_PRIVATE_MAP$0.get(Foo).prop;
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    Object.defineProperty(PRIVATE$1, "prop", {
      get: function() {
        return 4;
      }
    });
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateGetterWithThis()
func (c *RewriteClassMembersTest) TestPrivateGetterWithThis() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class Foo {
  get #prop() { return this.x; }
  getValue() {
    return this.#prop;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  prop: {
    get: function() {
      return this.$self.x;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE$2.$self = this;
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  getValue() {
    return PRIVATE_MAP$0.get(this).prop;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticGetterWithThis()
func (c *RewriteClassMembersTest) TestPrivateStaticGetterWithThis() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class Foo {
  static get #prop() { return this.x; }
  static getValue() {
    return Foo.#prop;
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static getValue() {
    return STATIC_PRIVATE_MAP$0.get(Foo).prop;
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE$1.$self = Foo;
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    Object.defineProperty(PRIVATE$1, "prop", {
      get: function() {
        return this.$self.x;
      }
    });
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateSetter()
func (c *RewriteClassMembersTest) TestPrivateSetter() {
	c.Test(
		`class Foo {
  set #prop(val) {}
  setValue(val) {
    this.#prop = val;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  prop: {
    set: function(val) {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  setValue(val$jscomp$1) {
    PRIVATE_MAP$0.get(this).prop = val$jscomp$1;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateSetterWithThis()
func (c *RewriteClassMembersTest) TestPrivateSetterWithThis() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class Foo {
  x = 0;
  set #prop(val) { this.x = val; }
  setValue(val) {
    this.#prop = val;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  prop: {
    set: function(val) {
      this.$self.x = val;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE$2.$self = this;
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    this.x = 0;
  }
  setValue(val$jscomp$1) {
    PRIVATE_MAP$0.get(this).prop = val$jscomp$1;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticSetter()
func (c *RewriteClassMembersTest) TestPrivateStaticSetter() {
	c.Test(
		`class Foo {
  static set #prop(val) {}
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    Object.defineProperty(PRIVATE$1, "prop", {
      set: function(val) {}
    });
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateIdInOperator()
func (c *RewriteClassMembersTest) TestPrivateIdInOperator() {
	c.Test(
		`class Foo {
  #field;
  brandCheck(x) { return #field in x; }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.field = void 0;
  }
  brandCheck(x) {
    return PRIVATE_MAP$0.has(x);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateIdInOperator_static()
func (c *RewriteClassMembersTest) TestPrivateIdInOperator_static() {
	c.Test(
		`class Foo {
  static #staticField;
  brandCheck(x) { return #staticField in x; }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  brandCheck(x) {
    return STATIC_PRIVATE_MAP$0.has(x);
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.staticField = void 0;
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateIdInOperator_methodAndAccessor()
func (c *RewriteClassMembersTest) TestPrivateIdInOperator_methodAndAccessor() {
	c.Test(
		`class Foo {
  #method() {}
  get #prop() { return 1; }
  brandCheck(x) { return #method in x && #prop in x; }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {}
  },
  prop: {
    get: function() {
      return 1;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  brandCheck(x) {
    return PRIVATE_MAP$0.has(x) && PRIVATE_MAP$0.has(x);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateDestructuringAssignment()
func (c *RewriteClassMembersTest) TestPrivateDestructuringAssignment() {
	c.IgnoreWarnings(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY)
	c.Test(
		`class Foo {
  #x;
  assignFromObj(obj) {
    ({ val: this.#x } = obj);
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.x = void 0;
  }
  assignFromObj(obj) {
    ({ val: PRIVATE_MAP$0.get(this).x } = obj);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodWithComplexParameters()
func (c *RewriteClassMembersTest) TestPrivateMethodWithComplexParameters() {
	c.Test(
		`class Foo {
  #method(a = 1, ...rest) { return a + rest.length; }
  callMethod() { return this.#method(10, 20, 30); }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function(a = 1, ...rest) {
      return a + rest.length;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  callMethod() {
    return PRIVATE_MAP$0.get(this).method.call(this, 10, 20, 30);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateUpdateAndCompoundAssignment()
func (c *RewriteClassMembersTest) TestPrivateUpdateAndCompoundAssignment() {
	c.Test(
		`class Foo {
  #x = 1;
  increment() {
    this.#x++;
    this.#x += 10;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.x = 1;
  }
  increment() {
    PRIVATE_MAP$0.get(this).x++;
    PRIVATE_MAP$0.get(this).x += 10;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateFieldAccess()
func (c *RewriteClassMembersTest) TestPrivateFieldAccess() {
	c.Test(
		`class Foo {
  #x = 1;
  getX(otherObj1) {
    return otherObj1.#x;
  }
  setX(otherObj2, v) {
    otherObj2.#x = v;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.x = 1;
  }
  getX(otherObj1) {
    return PRIVATE_MAP$0.get(otherObj1).x;
  }
  setX(otherObj2, v) {
    PRIVATE_MAP$0.get(otherObj2).x = v;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateFieldAccess_direct()
func (c *RewriteClassMembersTest) TestPrivateFieldAccess_direct() {
	c.Test(
		`class Foo {
  #x = 1;
  getX() {
    return this.#x;
  }
  setX(v) {
    this.#x = v;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.x = 1;
  }
  getX() {
    return PRIVATE_MAP$0.get(this).x;
  }
  setX(v) {
    PRIVATE_MAP$0.get(this).x = v;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodAccess()
func (c *RewriteClassMembersTest) TestPrivateMethodAccess() {
	c.Test(
		`class Foo {
  #method(paramX) { return paramX + 1; }
  callMethod(otherObj, argX) {
    return otherObj.#method(argX);
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function(paramX) {
      return paramX + 1;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  callMethod(otherObj, argX) {
    return PRIVATE_MAP$0.get(otherObj).method.call(otherObj, argX);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodAccess_direct()
func (c *RewriteClassMembersTest) TestPrivateMethodAccess_direct() {
	c.Test(
		`class Foo {
  #method(paramX) { return paramX + 1; }
  callMethod(argX) {
    return this.#method(argX);
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function(paramX) {
      return paramX + 1;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  callMethod(argX) {
    return PRIVATE_MAP$0.get(this).method.call(this, argX);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateGetterSetterAccess()
func (c *RewriteClassMembersTest) TestPrivateGetterSetterAccess() {
	c.Test(
		`class Foo {
  #propVal = 1;
  get #prop() { return this.#propVal; }
  set #prop(paramV) { this.#propVal = paramV; }
  access(otherObj, argV) {
    otherObj.#prop = otherObj.#prop + argV;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  prop: {
    get: function() {
      return PRIVATE_MAP$0.get(this.$self).propVal;
    },
    set: function(paramV) {
      PRIVATE_MAP$0.get(this.$self).propVal = paramV;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE$2.$self = this;
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    PRIVATE$2.propVal = 1;
  }
  access(otherObj, argV) {
    PRIVATE_MAP$0.get(otherObj).prop = PRIVATE_MAP$0.get(otherObj).prop + argV;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateGetterSetterAccess_direct()
func (c *RewriteClassMembersTest) TestPrivateGetterSetterAccess_direct() {
	c.Test(
		`class Foo {
  #propVal = 1;
  get #prop() { return this.#propVal; }
  set #prop(paramV) { this.#propVal = paramV; }
  access(argV) {
    this.#prop = this.#prop + argV;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  prop: {
    get: function() {
      return PRIVATE_MAP$0.get(this.$self).propVal;
    },
    set: function(paramV) {
      PRIVATE_MAP$0.get(this.$self).propVal = paramV;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE$2.$self = this;
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    PRIVATE$2.propVal = 1;
  }
  access(argV) {
    PRIVATE_MAP$0.get(this).prop = PRIVATE_MAP$0.get(this).prop + argV;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticMemberAccess()
func (c *RewriteClassMembersTest) TestPrivateStaticMemberAccess() {
	c.Test(
		`class Foo {
  static #staticField = 1;
  static #staticPropVal = 3;
  static #staticMethod() { return 2; }
  static get #staticProp() { return this.#staticPropVal; }
  static set #staticProp(paramV) { this.#staticPropVal = paramV; }
  static access(argV) {
    Foo.#staticProp = Foo.#staticField + Foo.#staticMethod() + Foo.#staticProp + argV;
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static access(argV) {
    STATIC_PRIVATE_MAP$0.get(Foo).staticProp =
        STATIC_PRIVATE_MAP$0.get(Foo).staticField
            + STATIC_PRIVATE_MAP$0.get(Foo).staticMethod.call(Foo)
            + STATIC_PRIVATE_MAP$0.get(Foo).staticProp
            + argV;
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE$1.$self = Foo;
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.staticMethod = function() {
      return 2;
    };
    Object.defineProperty(PRIVATE$1, "staticProp", {
      get: function() {
        return STATIC_PRIVATE_MAP$0.get(this.$self).staticPropVal;
      },
      set: function(paramV) {
        STATIC_PRIVATE_MAP$0.get(this.$self).staticPropVal = paramV;
      }
    });
    PRIVATE$1.staticField = 1;
    PRIVATE$1.staticPropVal = 3;
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticMemberAccess_this()
func (c *RewriteClassMembersTest) TestPrivateStaticMemberAccess_this() {
	c.Test(
		`class Foo {
  static #staticField = 1;
  static #staticPropVal = 3;
  static #staticMethod() { return 2; }
  static get #staticProp() { return this.#staticPropVal; }
  static set #staticProp(paramV) { this.#staticPropVal = paramV; }
  static access(argV) {
    this.#staticProp = this.#staticField + this.#staticMethod() + this.#staticProp + argV;
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static access(argV) {
    STATIC_PRIVATE_MAP$0.get(this).staticProp =
        STATIC_PRIVATE_MAP$0.get(this).staticField
            + STATIC_PRIVATE_MAP$0.get(this).staticMethod.call(this)
            + STATIC_PRIVATE_MAP$0.get(this).staticProp
            + argV;
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE$1.$self = Foo;
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.staticMethod = function() {
      return 2;
    };
    Object.defineProperty(PRIVATE$1, "staticProp", {
      get: function() {
        return STATIC_PRIVATE_MAP$0.get(this.$self).staticPropVal;
      },
      set: function(paramV) {
        STATIC_PRIVATE_MAP$0.get(this.$self).staticPropVal = paramV;
      }
    });
    PRIVATE$1.staticField = 1;
    PRIVATE$1.staticPropVal = 3;
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining() {
	c.Test(
		`class Foo {
  #field = 2;
  access(otherObj) {
    return otherObj?.#field;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.field = 2;
  }
  access(otherObj) {
    return otherObj == null ? void 0 : PRIVATE_MAP$0.get(otherObj).field;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_direct()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_direct() {
	c.Test(
		`class Foo {
  #field = 2;
  access() {
    return this?.#field;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.field = 2;
  }
  access() {
    return this == null ? void 0 : PRIVATE_MAP$0.get(this).field;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodAccessWithOptionalChain()
func (c *RewriteClassMembersTest) TestPrivateMethodAccessWithOptionalChain() {
	c.Test(
		`class Foo {
  #method() {}
  access(otherObj) {
    return otherObj?.#method;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  access(otherObj) {
    return otherObj == null ? void 0 : PRIVATE_MAP$0.get(otherObj).method;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodCallWithOptionalChain()
func (c *RewriteClassMembersTest) TestPrivateMethodCallWithOptionalChain() {
	c.Test(
		`class Foo {
  #method() {}
  access(otherObj) {
    return otherObj?.#method();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  access(otherObj) {
    return otherObj == null ? void 0 : PRIVATE_MAP$0.get(otherObj).method.call(otherObj);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodCallWithOptionalChain_sideEffectingReceiver()
func (c *RewriteClassMembersTest) TestPrivateMethodCallWithOptionalChain_sideEffectingReceiver() {
	c.Test(
		`class Foo {
  #method() {}
  access() {
    return getObj()?.#method();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$2 = Object.create(null, {
  method: {
    value: function() {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$3 = Object.create(PRIVATE_PROTO$2);
    PRIVATE_MAP$0.set(this, PRIVATE$3);
  }
  access() {
    let TMP$1;
    return (TMP$1 = getObj(), TMP$1 == null ? void 0 : PRIVATE_MAP$0.get(TMP$1).method.call(TMP$1));
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_sideEffectingReceiverField()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_sideEffectingReceiverField() {
	c.Test(
		`class Foo {
  #field = 10;
  access() {
    return getObj()?.#field;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    PRIVATE$2.field = 10;
  }
  access() {
    let TMP$1;
    return (TMP$1 = getObj(), TMP$1 == null ? void 0 : PRIVATE_MAP$0.get(TMP$1).field);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_trailingSegment()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_trailingSegment() {
	c.Test(
		`class Foo {
  #field = { prop: 42 };
  access(obj) {
    return obj?.#field.prop;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.field = { prop: 42 };
  }
  access(obj) {
    return obj == null ? void 0 : PRIVATE_MAP$0.get(obj).field.prop;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_continuationChain()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_continuationChain() {
	c.Test(
		`class Foo {
  #field = 10;
  access(obj) {
    return obj?.a.#field;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    PRIVATE$2.field = 10;
  }
  access(obj) {
    let TMP$1;
    return (TMP$1 = obj?.a, TMP$1 == null ? void 0 : PRIVATE_MAP$0.get(TMP$1).field);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_outerClassFieldFromInner()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_outerClassFieldFromInner() {
	c.Test(
		`class Outer {
  #x = 1;
  createInner() {
    return class Inner {
      readOuter(o) { return o?.#x; }
    };
  }
}
`,
		`const PRIVATE_MAP$1 = new $jscomp.PrivateMap();
class Outer {
  constructor() {
    const PRIVATE$2 = Object.create(null);
    PRIVATE_MAP$1.set(this, PRIVATE$2);
    PRIVATE$2.x = 1;
  }
  createInner() {
    const CLASS_DECL$0 = class {
      readOuter(o) {
        return o == null ? void 0 : PRIVATE_MAP$1.get(o).x;
      }
    };
    return CLASS_DECL$0;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_instantiatedCalls()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_instantiatedCalls() {
	c.Test(
		`class Foo {
  #field = 42;
  #method() { return 100; }
  getField(obj) {
    return obj?.#field;
  }
  getMethod(obj) {
    return obj?.#method;
  }
}
const foo = new Foo();
const validField = foo.getField(foo);
const nullField = foo.getField(null);
const validMethod = foo.getMethod(foo);
const nullMethod = foo.getMethod(null);
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {
      return 100;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    PRIVATE$2.field = 42;
  }
  getField(obj) {
    return obj == null ? void 0 : PRIVATE_MAP$0.get(obj).field;
  }
  getMethod(obj$jscomp$1) {
    return obj$jscomp$1 == null ? void 0 : PRIVATE_MAP$0.get(obj$jscomp$1).method;
  }
}
const foo = new Foo();
const validField = foo.getField(foo);
const nullField = foo.getField(null);
const validMethod = foo.getMethod(foo);
const nullMethod = foo.getMethod(null);
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticFieldWithOptionalChain()
func (c *RewriteClassMembersTest) TestPrivateStaticFieldWithOptionalChain() {
	c.Test(
		`class Foo {
  static #staticField = 10;
  static getStaticField(cls) {
    return cls?.#staticField;
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static getStaticField(cls) {
    return cls == null ? void 0 : STATIC_PRIVATE_MAP$0.get(cls).staticField;
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.staticField = 10;
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticMethodWithOptionalChain()
func (c *RewriteClassMembersTest) TestPrivateStaticMethodWithOptionalChain() {
	c.Test(
		`class Foo {
  static #staticMethod() { return 20; }
  static callStaticMethod(cls) {
    return cls?.#staticMethod();
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static callStaticMethod(cls) {
    return cls == null ? void 0 : STATIC_PRIVATE_MAP$0.get(cls).staticMethod.call(cls);
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.staticMethod = function() {
      return 20;
    };
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateGetterWithOptionalChain()
func (c *RewriteClassMembersTest) TestPrivateGetterWithOptionalChain() {
	c.Test(
		`class Foo {
  get #getter() { return this; }
  accessGetter(obj) {
    return obj?.#getter;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  getter: {
    get: function() { return this.$self; }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE$2.$self = this;
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  accessGetter(obj) {
    return obj == null ? void 0 : PRIVATE_MAP$0.get(obj).getter;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticGetterWithOptionalChain()
func (c *RewriteClassMembersTest) TestPrivateStaticGetterWithOptionalChain() {
	c.Test(
		`class Foo {
  static get #staticGetter() { return this; }
  static accessStaticGetter(cls) {
    return cls?.#staticGetter;
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static accessStaticGetter(cls) {
    return cls == null ? void 0 : STATIC_PRIVATE_MAP$0.get(cls).staticGetter;
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE$1.$self = Foo;
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    Object.defineProperty(PRIVATE$1, "staticGetter", {
      get: function() { return this.$self; }
    });
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticMethodWithOptionalMethodCall()
func (c *RewriteClassMembersTest) TestPrivateStaticMethodWithOptionalMethodCall() {
	c.Test(
		`class Foo {
  static #staticMethod(a) { return a; }
  static call(a) {
    return Foo.#staticMethod?.(a);
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static call(a$jscomp$1) {
    return STATIC_PRIVATE_MAP$0.get(Foo).staticMethod?.call(Foo, a$jscomp$1);
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.staticMethod = function(a) {
      return a;
    };
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticMethodWithBothReceiverAndMethodOptional()
func (c *RewriteClassMembersTest) TestPrivateStaticMethodWithBothReceiverAndMethodOptional() {
	c.Test(
		`class Foo {
  static #staticMethod(a) { return a; }
  static call(a) {
    return Foo?.#staticMethod?.(a);
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static call(a$jscomp$1) {
    return Foo == null ? void 0 : STATIC_PRIVATE_MAP$0.get(Foo).staticMethod?.call(Foo, a$jscomp$1);
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.staticMethod = function(a) {
      return a;
    };
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_optionalMethodCall()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_optionalMethodCall() {
	c.Test(
		`class Foo {
  #method(a) { return a; }
  call(obj, a) {
    return obj.#method?.(a);
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function(a) { return a; }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  call(obj, a$jscomp$1) {
    return PRIVATE_MAP$0.get(obj).method?.call(obj, a$jscomp$1);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_bothReceiverAndMethodOptional()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_bothReceiverAndMethodOptional() {
	c.Test(
		`class Foo {
  #method(a) { return a; }
  call(obj, a) {
    return obj?.#method?.(a);
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function(a) { return a; }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  call(obj, a$jscomp$1) {
    return obj == null ? void 0 : PRIVATE_MAP$0.get(obj).method?.call(obj, a$jscomp$1);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_bothOptional_chainContinues()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_bothOptional_chainContinues() {
	// TODO(b/571219920): trailing segment left outside nullish hook
	e := AssertThrows[jutil.IRuntimeException](func() {
		c.TestSame(
			`class Foo {
  #method(a) { return {x: a}; }
  call(obj, a) {
    return obj?.#method?.(a).x;
  }
}
`)
	})
	// Throwables.getRootCause(e)
	var root jutil.Throwable = e
	for root.GetCause() != nil {
		root = root.GetCause()
	}
	IsInstanceOf[*jutil.IllegalStateException](AssertThat(root))
	AssertThat(root.GetMessage()).
		Contains("Start of optional chain node OPTCHAIN_GETPROP is not marked as the start.")
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_fieldHoldingFunction()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_fieldHoldingFunction() {
	c.Test(
		`class Foo {
  #field = (a) => a;
  call(a) {
    return this.#field?.(a);
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.field = (a) => {
      return a;
    };
  }
  call(a$jscomp$1) {
    return PRIVATE_MAP$0.get(this).field?.call(this, a$jscomp$1);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_staticFieldHoldingArrowFunction()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_staticFieldHoldingArrowFunction() {
	c.Test(
		`class Foo {
  static #field = (a) => a;
  static call(a) {
    return Foo.#field?.(a);
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static call(a$jscomp$1) {
    return STATIC_PRIVATE_MAP$0.get(Foo).field?.call(Foo, a$jscomp$1);
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.field = (a) => {
      return a;
    };
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_fieldHoldingFunctionExpressionWithThis()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_fieldHoldingFunctionExpressionWithThis() {
	c.Test(
		`class Foo {
  x = 1;
  #field = function() { return this.x; };
  call() {
    return this.#field?.();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    this.x = 1;
    PRIVATE$1.field = function() {
      return this.x;
    };
  }
  call() {
    return PRIVATE_MAP$0.get(this).field?.call(this);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_staticFieldHoldingFunctionExpressionWithThis()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_staticFieldHoldingFunctionExpressionWithThis() {
	c.Test(
		`class Foo {
  static x = 1;
  static #field = function() { return this.x; };
  static call() {
    return Foo.#field?.();
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static call() {
    return STATIC_PRIVATE_MAP$0.get(Foo).field?.call(Foo);
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    Foo.x = 1;
    PRIVATE$1.field = function() {
      return this.x;
    };
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_chainedAccess()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_chainedAccess() {
	c.Test(
		`class Foo {
  #field = { prop: 123 };
  getProp(obj) {
    return obj?.#field?.prop;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.field = { prop: 123 };
  }
  getProp(obj) {
    return (obj == null ? void 0 : PRIVATE_MAP$0.get(obj).field)?.prop;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateOptionalChaining_chainedMethodCall()
func (c *RewriteClassMembersTest) TestPrivateOptionalChaining_chainedMethodCall() {
	c.Test(
		`class Foo {
  #method() {}
  callMethod(obj) {
    return obj.foo?.#method();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$2 = Object.create(null, {
  method: {
    value: function() {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$3 = Object.create(PRIVATE_PROTO$2);
    PRIVATE_MAP$0.set(this, PRIVATE$3);
  }
  callMethod(obj) {
    let TMP$1;
    return (TMP$1 = obj.foo, TMP$1 == null ? void 0 : PRIVATE_MAP$0.get(TMP$1).method.call(TMP$1));
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateInheritanceAndSuper()
func (c *RewriteClassMembersTest) TestPrivateInheritanceAndSuper() {
	c.Test(
		`class Base {
  constructor() {
    this.baseProp = 1;
  }
}
class Sub extends Base {
  #field = 2;
  #method() { return 3; }
  constructor() {
    super();
    this.#field = this.baseProp + this.#method();
  }
}
`,
		`class Base {
  constructor() {
    this.baseProp = 1;
  }
}
const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {
      return 3;
    }
  }
});
class Sub extends Base {
  constructor() {
    super();
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    PRIVATE$2.field = 2;
    PRIVATE_MAP$0.get(this).field = this.baseProp + PRIVATE_MAP$0.get(this).method.call(this);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateInheritanceWithSuperclassHavingPrivates()
func (c *RewriteClassMembersTest) TestPrivateInheritanceWithSuperclassHavingPrivates() {
	c.Test(
		`class Base {
  #baseField = 1;
  #baseMethod() { return this.#baseField; }
  getBase() { return this.#baseMethod(); }
}
class Sub extends Base {
  #subField = 2;
  #subMethod() { return this.#subField; }
  getSub() { return this.#subMethod(); }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  baseMethod: {
    value: function() {
      return PRIVATE_MAP$0.get(this).baseField;
    }
  }
});
class Base {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    PRIVATE$2.baseField = 1;
  }
  getBase() {
    return PRIVATE_MAP$0.get(this).baseMethod.call(this);
  }
}
const PRIVATE_MAP$3 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$4 = Object.create(null, {
  subMethod: {
    value: function() {
      return PRIVATE_MAP$3.get(this).subField;
    }
  }
});
class Sub extends Base {
  constructor() {
    super(...arguments);
    const PRIVATE$5 = Object.create(PRIVATE_PROTO$4);
    PRIVATE_MAP$3.set(this, PRIVATE$5);
    PRIVATE$5.subField = 2;
  }
  getSub() {
    return PRIVATE_MAP$3.get(this).subMethod.call(this);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodWithSuperReference()
func (c *RewriteClassMembersTest) TestPrivateMethodWithSuperReference() {
	c.Test(
		`class Base {
  baseProp = 5;
  baseMethod() { return 10; }
}
class Sub extends Base {
  #privateMethod() {
    return super.baseMethod() + super.baseProp;
  }
  callPrivate() {
    return this.#privateMethod();
  }
}
`,
		`class Base {
  constructor() {
    this.baseProp = 5;
  }
  baseMethod() { return 10; }
}
const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  privateMethod: {
    value: function() {
      return Base.prototype.baseMethod.call(this) + Reflect.get(Base.prototype, JSCompiler_renameProperty("baseProp", Base), this);
    }
  }
});
class Sub extends Base {
  constructor() {
    super(...arguments);
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  callPrivate() {
    return PRIVATE_MAP$0.get(this).privateMethod.call(this);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodValueRead()
func (c *RewriteClassMembersTest) TestPrivateMethodValueRead() {
	c.Test(
		`class Foo {
  #method() {}
  getMethod() {
    return this.#method;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  getMethod() {
    return PRIVATE_MAP$0.get(this).method;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodCallingPrivateMethod()
func (c *RewriteClassMembersTest) TestPrivateMethodCallingPrivateMethod() {
	c.Test(
		`class Foo {
  #a() { return 1; }
  #b() { return this.#a() + 2; }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  a: {
    value: function() {
      return 1;
    }
  },
  b: {
    value: function() {
      return PRIVATE_MAP$0.get(this).a.call(this) + 2;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodCallingPublicMethod()
func (c *RewriteClassMembersTest) TestPrivateMethodCallingPublicMethod() {
	c.Test(
		`class Foo {
  a() { return 1; }
  #b() { return this.a() + 2; }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  b: {
    value: function() {
      return this.a() + 2;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  a() { return 1; }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodWithNestedClassUsingSuper()
func (c *RewriteClassMembersTest) TestPrivateMethodWithNestedClassUsingSuper() {
	c.Test(
		`class A {
  m() {}
}
class B {
  #priv() {
    return class C extends A {
      m() {
        super.m();
      }
    }
  }
}
`,
		`class A {
  m() {}
}
const PRIVATE_MAP$1 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$2 = Object.create(null, {
  priv: {
    value: function() {
      const CLASS_DECL$0 = class extends A {
        m() {
          super.m();
        }
      };
      return CLASS_DECL$0;
    }
  }
});
class B {
  constructor() {
    const PRIVATE$3 = Object.create(PRIVATE_PROTO$2);
    PRIVATE_MAP$1.set(this, PRIVATE$3);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodWithSuper()
func (c *RewriteClassMembersTest) TestPrivateMethodWithSuper() {
	c.Test(
		`class Base {
  foo() { return 10; }
}
class Child extends Base {
  #privateMethod() {
    return super.foo() + 1;
  }
  callPrivate() {
    return this.#privateMethod();
  }
}
`,
		`class Base {
  foo() { return 10; }
}
const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  privateMethod: {
    value: function() {
      return Base.prototype.foo.call(this) + 1;
    }
  }
});
class Child extends Base {
  constructor() {
    super(...arguments);
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  callPrivate() {
    return PRIVATE_MAP$0.get(this).privateMethod.call(this);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticMethodWithSuper()
func (c *RewriteClassMembersTest) TestPrivateStaticMethodWithSuper() {
	c.Test(
		`class Base {
  static staticFoo() { return 20; }
}
class Child extends Base {
  static #staticPrivateMethod() {
    return super.staticFoo() + 5;
  }
  static callStaticPrivate() {
    return Child.#staticPrivateMethod();
  }
}
`,
		`class Base {
  static staticFoo() { return 20; }
}
const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Child extends Base {
  static callStaticPrivate() {
    return STATIC_PRIVATE_MAP$0.get(Child).staticPrivateMethod.call(Child);
  }
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Child, PRIVATE$1);
    PRIVATE$1.staticPrivateMethod = function() {
      return Base.staticFoo() + 5;
    };
  }
}
Child.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodSideEffectingReceiver()
func (c *RewriteClassMembersTest) TestPrivateMethodSideEffectingReceiver() {
	c.Test(
		`class Foo {
  #method() {}
  callMethod(getObj) {
    return getObj().#method();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$2 = Object.create(null, {
  method: {
    value: function() {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$3 = Object.create(PRIVATE_PROTO$2);
    PRIVATE_MAP$0.set(this, PRIVATE$3);
  }
  callMethod(getObj) {
    let $jscomp$tmp$m1146332801$1;
    return ($jscomp$tmp$m1146332801$1 = getObj(), PRIVATE_MAP$0.get($jscomp$tmp$m1146332801$1).method.call($jscomp$tmp$m1146332801$1));
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodPropertyAccessReceiver()
func (c *RewriteClassMembersTest) TestPrivateMethodPropertyAccessReceiver() {
	c.Test(
		`class Foo {
  #method() {}
  static callOnProp(holder) {
    return holder['prop'].#method();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$2 = Object.create(null, {
  method: {
    value: function() {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$3 = Object.create(PRIVATE_PROTO$2);
    PRIVATE_MAP$0.set(this, PRIVATE$3);
  }
  static callOnProp(holder) {
    let $jscomp$tmp$m1146332801$1;
    return ($jscomp$tmp$m1146332801$1 = holder['prop'],
        PRIVATE_MAP$0.get($jscomp$tmp$m1146332801$1).method.call($jscomp$tmp$m1146332801$1));
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodAsyncAndGenerator()
func (c *RewriteClassMembersTest) TestPrivateMethodAsyncAndGenerator() {
	c.Test(
		`class Foo {
  async #asyncMethod() {}
  *#genMethod() {}
  async *#asyncGenMethod() {}
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  asyncMethod: {
    value: async function() {}
  },
  genMethod: {
    value: function*() {}
  },
  asyncGenMethod: {
    value: async function*() {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodCalledOnOtherInstance()
func (c *RewriteClassMembersTest) TestPrivateMethodCalledOnOtherInstance() {
	c.Test(
		`class Foo {
  #method() { return 42; }
  callOther(other) {
    return other.#method();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {
      return 42;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  callOther(other) {
    return PRIVATE_MAP$0.get(other).method.call(other);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodMultiplePrivateMethodsAndFields()
func (c *RewriteClassMembersTest) TestPrivateMethodMultiplePrivateMethodsAndFields() {
	c.Test(
		`class Foo {
  #field = 10;
  #getValue() { return this.#field; }
  publicMethod() {
    return this.#getValue();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  getValue: {
    value: function() {
      return PRIVATE_MAP$0.get(this).field;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    PRIVATE$2.field = 10;
  }
  publicMethod() {
    return PRIVATE_MAP$0.get(this).getValue.call(this);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodWriteUnsupported()
func (c *RewriteClassMembersTest) TestPrivateMethodWriteUnsupported() {
	c.TestError(
		`class Foo {
  #method() {}
  write() {
    this.#method = 1;
  }
}
`,
		Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT)
}

// Java: RewriteClassMembersTest.testPrivateAccessorCompoundAssignment()
func (c *RewriteClassMembersTest) TestPrivateAccessorCompoundAssignment() {
	c.Test(
		`class Foo {
  get #prop() { return 1; }
  set #prop(v) {}
  update() {
    this.#prop += 5;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  prop: {
    get: function() {
      return 1;
    },
    set: function(v) {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  update() {
    PRIVATE_MAP$0.get(this).prop += 5;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateAccessorIncDec()
func (c *RewriteClassMembersTest) TestPrivateAccessorIncDec() {
	c.Test(
		`class Foo {
  get #prop() { return 1; }
  set #prop(v) {}
  update() {
    ++this.#prop;
    this.#prop++;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  prop: {
    get: function() {
      return 1;
    },
    set: function(v) {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  update() {
    ++PRIVATE_MAP$0.get(this).prop;
    PRIVATE_MAP$0.get(this).prop++;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateAccessorCompoundAssignmentSideEffectingReceiver()
func (c *RewriteClassMembersTest) TestPrivateAccessorCompoundAssignmentSideEffectingReceiver() {
	c.Test(
		`class Foo {
  get #prop() { return 1; }
  set #prop(v) {}
  update() {
    getObj().#prop += 5;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  prop: {
    get: function() {
      return 1;
    },
    set: function(v) {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  update() {
    PRIVATE_MAP$0.get(getObj()).prop += 5;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodExplicitCall()
func (c *RewriteClassMembersTest) TestPrivateMethodExplicitCall() {
	c.IgnoreWarnings(TypeCheck_INEXISTENT_PROPERTY)
	c.Test(
		`class Foo {
  #method() { return this.val; }
  test(other) {
    return this.#method.call(other);
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {
      return this.val;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  test(other) {
    return PRIVATE_MAP$0.get(this).method.call(other);
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodWithSuperGetterAndSetter()
func (c *RewriteClassMembersTest) TestPrivateMethodWithSuperGetterAndSetter() {
	c.Test(
		`class S {
  get g() { return 'g S'; }
  set s(v) {}
}
class C extends S {
  get g() { return 'g C'; }
  #checkG(v) {
    super.s = v;
    return this.g + ' ' + super.g;
  }
}
`,
		`class S {
  get g() { return 'g S'; }
  set s(v) {}
}
const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  checkG: {
    value: function(v$jscomp$1) {
      Reflect.set(S.prototype, JSCompiler_renameProperty("s", S), v$jscomp$1, this);
      return this.g + " " + Reflect.get(S.prototype, JSCompiler_renameProperty("g", S), this);
    }
  }
});
class C extends S {
  constructor() {
    super(...arguments);
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  get g() { return 'g C'; }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateGetterWithSuper()
func (c *RewriteClassMembersTest) TestPrivateGetterWithSuper() {
	c.Test(
		`class Base {
  get baseProp() { return 'base'; }
}
class Child extends Base {
  get #foo() {
    return super.baseProp + ' child';
  }
  getFoo() {
    return this.#foo;
  }
}
`,
		`class Base {
  get baseProp() { return 'base'; }
}
const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  foo: {
    get: function() {
      return Reflect.get(Base.prototype, JSCompiler_renameProperty("baseProp", Base), this.$self) + " child";
    }
  }
});
class Child extends Base {
  constructor() {
    super(...arguments);
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE$2.$self = this;
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  getFoo() {
    return PRIVATE_MAP$0.get(this).foo;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateSetterWithSuperProperty()
func (c *RewriteClassMembersTest) TestPrivateSetterWithSuperProperty() {
	c.Test(
		`class Base {
  set baseProp(v) {}
}
class Child extends Base {
  set #foo(v) {
    super.baseProp = v;
  }
  setFoo(v) {
    this.#foo = v;
  }
}
`,
		`class Base {
  set baseProp(v) {}
}
const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  foo: {
    set: function(v$jscomp$1) {
      Reflect.set(Base.prototype, JSCompiler_renameProperty("baseProp", Base), v$jscomp$1, this.$self);
    }
  }
});
class Child extends Base {
  constructor() {
    super(...arguments);
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE$2.$self = this;
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  setFoo(v$jscomp$2) {
    PRIVATE_MAP$0.get(this).foo = v$jscomp$2;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodSuperInClassExpression()
func (c *RewriteClassMembersTest) TestPrivateMethodSuperInClassExpression() {
	c.Test(
		`class C extends (class { m() { return 'x'; } }) {
  m() { return 'y'; }
  #pm() { return super.m(); }
}
`,
		`const CLASS_EXTENDS$0 = (() => {
  const CLASS_DECL$1 = class {
    m() {
      return "x";
    }
  };
  return CLASS_DECL$1;
})();
const PRIVATE_MAP$2 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$3 = Object.create(null, {
  pm: {
    value: function() {
      return CLASS_EXTENDS$0.prototype.m.call(this);
    }
  }
});
class C extends CLASS_EXTENDS$0 {
  constructor() {
    super(...arguments);
    const PRIVATE$4 = Object.create(PRIVATE_PROTO$3);
    PRIVATE_MAP$2.set(this, PRIVATE$4);
  }
  m() {
    return "y";
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateSetterWithSuper()
func (c *RewriteClassMembersTest) TestPrivateSetterWithSuper() {
	c.Test(
		`class Base {
  setVal(v) {}
}
class Sub extends Base {
  set #prop(v) {
    super.setVal(v);
  }
  update(v) {
    this.#prop = v;
  }
}
`,
		`class Base {
  setVal(v) {}
}
const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  prop: {
    set: function(v$jscomp$1) {
      Base.prototype.setVal.call(this.$self, v$jscomp$1);
    }
  }
});
class Sub extends Base {
  constructor() {
    super(...arguments);
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE$2.$self = this;
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  update(v$jscomp$2) {
    PRIVATE_MAP$0.get(this).prop = v$jscomp$2;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateFieldWithSuperArrowFunction()
func (c *RewriteClassMembersTest) TestPrivateFieldWithSuperArrowFunction() {
	c.Test(
		`class Sup {
  m() {}
}
class Sub extends Sup {
  #x = () => super.m();
}
`,
		`class Sup {
  m() {}
}
const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Sub extends Sup {
  constructor() {
    super(...arguments);
    const PRIVATE$1 = Object.create(null);
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.x = () => {
      return super.m();
    };
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateNestedClassScoping()
func (c *RewriteClassMembersTest) TestPrivateNestedClassScoping() {
	c.Test(
		`class Outer {
  #x = 1;
  createInner() {
    return class Inner {
      readOuter(o) { return o.#x; }
    };
  }
}
`,
		`const PRIVATE_MAP$1 = new $jscomp.PrivateMap();
class Outer {
  constructor() {
    const PRIVATE$2 = Object.create(null);
    PRIVATE_MAP$1.set(this, PRIVATE$2);
    PRIVATE$2.x = 1;
  }
  createInner() {
    const CLASS_DECL$0 = class {
      readOuter(o) {
        return PRIVATE_MAP$1.get(o).x;
      }
    };
    return CLASS_DECL$0;
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMemberInConstructorArrowFunction()
func (c *RewriteClassMembersTest) TestPrivateMemberInConstructorArrowFunction() {
	c.Test(
		`class Foo {
  #field = 10;
  #method() { return 20; }
  constructor() {
    const getVal = () => this.#field + this.#method();
    this.val = getVal();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {
      return 20;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
    PRIVATE$2.field = 10;
    const getVal = () => {
      return PRIVATE_MAP$0.get(this).field + PRIVATE_MAP$0.get(this).method.call(this);
    };
    this.val = getVal();
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateMethodCalledInsideArrowFunction()
func (c *RewriteClassMembersTest) TestPrivateMethodCalledInsideArrowFunction() {
	c.Test(
		`class Foo {
  #method() { return 42; }
  getArrow() {
    return () => this.#method();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  method: {
    value: function() {
      return 42;
    }
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  getArrow() {
    return () => {
      return PRIVATE_MAP$0.get(this).method.call(this);
    };
  }
}
`)
}

// Java: RewriteClassMembersTest.testNestedClassAccessingOuterPrivateMethod()
func (c *RewriteClassMembersTest) TestNestedClassAccessingOuterPrivateMethod() {
	c.Test(
		`class Outer {
  #foo() { return 1; }
  method() {
    class Inner {
      #bar() { return 2; }
      use(outer) {
        return outer.#foo() + this.#bar();
      }
    }
    return new Inner();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$4 = Object.create(null, {
  foo: {
    value: function() {
      return 1;
    }
  }
});
class Outer {
  constructor() {
    const PRIVATE$5 = Object.create(PRIVATE_PROTO$4);
    PRIVATE_MAP$0.set(this, PRIVATE$5);
  }
  method() {
    const PRIVATE_MAP$1 = new $jscomp.PrivateMap();
    const PRIVATE_PROTO$2 = Object.create(null, {
      bar: {
        value: function() {
          return 2;
        }
      }
    });
    class Inner {
      constructor() {
        const PRIVATE$3 = Object.create(PRIVATE_PROTO$2);
        PRIVATE_MAP$1.set(this, PRIVATE$3);
      }
      use(outer) {
        return PRIVATE_MAP$0.get(outer).foo.call(outer) + PRIVATE_MAP$1.get(this).bar.call(this);
      }
    }
    return new Inner();
  }
}
`)
}

// Java: RewriteClassMembersTest.testNestedClassShadowingPrivateMemberName()
func (c *RewriteClassMembersTest) TestNestedClassShadowingPrivateMemberName() {
	c.Test(
		`class Outer {
  #foo() { return 1; }
  method() {
    class Inner {
      #foo() { return 2; }
      use(outer) {
        return outer.#foo() + this.#foo();
      }
    }
    return new Inner();
  }
}
`,
		`const PRIVATE_MAP$3 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$4 = Object.create(null, {
  foo: {
    value: function() {
      return 1;
    }
  }
});
class Outer {
  constructor() {
    const PRIVATE$5 = Object.create(PRIVATE_PROTO$4);
    PRIVATE_MAP$3.set(this, PRIVATE$5);
  }
  method() {
    const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
    const PRIVATE_PROTO$1 = Object.create(null, {
      foo: {
        value: function() {
          return 2;
        }
      }
    });
    class Inner {
      constructor() {
        const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
        PRIVATE_MAP$0.set(this, PRIVATE$2);
      }
      use(outer) {
        return PRIVATE_MAP$0.get(outer).foo.call(outer) + PRIVATE_MAP$0.get(this).foo.call(this);
      }
    }
    return new Inner();
  }
}
`)
}

// Java: RewriteClassMembersTest.testPrivateStaticMemberAccessInStaticBlock()
func (c *RewriteClassMembersTest) TestPrivateStaticMemberAccessInStaticBlock() {
	c.Test(
		`class Foo {
  static #x = 42;
  static #foo() { return this.#x; }
  static y;
  static {
    this.y = this.#foo();
  }
}
`,
		`const STATIC_PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  static STATIC_INIT$2() {
    const PRIVATE$1 = Object.create(null);
    STATIC_PRIVATE_MAP$0.set(Foo, PRIVATE$1);
    PRIVATE$1.foo = function() {
      return STATIC_PRIVATE_MAP$0.get(this).x;
    };
    PRIVATE$1.x = 42;
    Foo.y = void 0;
    {
      Foo.y = STATIC_PRIVATE_MAP$0.get(Foo).foo.call(Foo);
    }
  }
}
Foo.STATIC_INIT$2();
`)
}

// Java: RewriteClassMembersTest.testNestedClassInsideStaticBlock()
func (c *RewriteClassMembersTest) TestNestedClassInsideStaticBlock() {
	c.Test(
		`class Outer {
  #foo() { return 1; }
  static innerInstance;
  static {
    class Inner {
      #bar() { return 2; }
      use(outer) {
        return outer.#foo() + this.#bar();
      }
    }
    Outer.innerInstance = new Inner();
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$4 = Object.create(null, {
  foo: {
    value: function() {
      return 1;
    }
  }
});
class Outer {
  constructor() {
    const PRIVATE$5 = Object.create(PRIVATE_PROTO$4);
    PRIVATE_MAP$0.set(this, PRIVATE$5);
  }
}
Outer.innerInstance = void 0;
{
  const PRIVATE_MAP$1 = new $jscomp.PrivateMap();
  const PRIVATE_PROTO$2 = Object.create(null, {
    bar: {
      value: function() {
        return 2;
      }
    }
  });
  class Inner {
    constructor() {
      const PRIVATE$3 = Object.create(PRIVATE_PROTO$2);
      PRIVATE_MAP$1.set(this, PRIVATE$3);
    }
    use(outer) {
      return PRIVATE_MAP$0.get(outer).foo.call(outer) + PRIVATE_MAP$1.get(this).bar.call(this);
    }
  }
  Outer.innerInstance = new Inner();
}
`)
}

// Java: RewriteClassMembersTest.testClassStaticBlock_superRef_onClassWithNameSpace()
func (c *RewriteClassMembersTest) TestClassStaticBlock_superRef_onClassWithNameSpace() {
	c.Test(
		`const ns = {};
ns.B = class {
  static y = 3;
};
class C extends ns.B {
  static {
    let x = super.y;
  }
}
`,
		`const ns = {};
ns.B = class {};
ns.B.y = 3;
class C extends ns.B {}
{
  let x = ns.B.y;
}
`)
}

// Java: RewriteClassMembersTest.testClassStaticBlock_thisRef()
func (c *RewriteClassMembersTest) TestClassStaticBlock_thisRef() {
	src := `class C {
  static {
    C.x = 2
    const y = this.x
  }
}
`

	c.Test(
		src,
		`class C {}
{
  C.x = 2;
  const y = C.x;
}
`)

	c.Test(
		c.WithOptions().UseStaticInheritance(),
		src,
		`class C {}
{
  C.x = 2;
  const y = C.x;
}
`)
}

// Java: RewriteClassMembersTest.testClassStaticBlock_varInStaticBlock()
func (c *RewriteClassMembersTest) TestClassStaticBlock_varInStaticBlock() {
	c.Test(
		`var z = 1
class C {
  static {
    let x = 2
    var z = 3;
  }
}
`,
		`var z = 1;
class C {}
{
  let x = 2;
  var z$jscomp$1 = 3;
}
`)
}

// Java: RewriteClassMembersTest.testClassStaticBlock_classExpression()
func (c *RewriteClassMembersTest) TestClassStaticBlock_classExpression() {
	c.Test(
		`let C = class {
  static prop = 5;
};
let D = class extends C {
  static {
    this.prop = 10;
  }
};
`,
		`let C = class {};
C.prop = 5;
let D = class extends C {};
{
  D.prop = 10;
}
`)
}

// Java: RewriteClassMembersTest.testClassStaticBlock_multipleClassesInLet()
func (c *RewriteClassMembersTest) TestClassStaticBlock_multipleClassesInLet() {
	c.Test(
		`let C = class {
  static prop = 5;
},
D = class extends C {
  static {
    this.prop = 10;
  }
}
`,
		`let C = class {};
C.prop = 5;
let D = class extends C {};
{
  D.prop = 10;
}
`)
}

// Java: RewriteClassMembersTest.testClassStaticBlock_fieldAndBlock()
func (c *RewriteClassMembersTest) TestClassStaticBlock_fieldAndBlock() {
	c.Test(
		`class C {
  static f;
  static {
    C.f = 1;
  }
}
`,
		`class C {}
C.f = void 0;
{
  C.f = 1;
}
`)
}

// Java: RewriteClassMembersTest.testMultipleStaticBlocks()
func (c *RewriteClassMembersTest) TestMultipleStaticBlocks() {
	c.Test(
		`var z = 1
/** @unrestricted */
class C {
  static x = 2;
  static {
    z = z + this.x;
  }
  static [z] = 3;
  static w = 5;
  static {
    z = z + this.w;
  }
}
`,
		`var z = 1;
var COMP_FIELD$0 = z;
class C {}
C.x = 2;
{
  z = z + C.x;
}
C[COMP_FIELD$0] = 3;
C.w = 5;
{
  z = z + C.w;
}
`)
}

// Java: RewriteClassMembersTest.testThisInNonStaticPublicField()
func (c *RewriteClassMembersTest) TestThisInNonStaticPublicField() {
	src := `class A {
  /** @suppress {partialAlias} */
  b = 'word';
  c = this.b;
}
`
	c.Test(
		src,
		`class A {
  constructor() {
    /** @suppress {partialAlias} */
    this.b = 'word';
    this.c = this.b;
  }
}
`)
	c.TestSame(c.WithOptions().UseEs2022LanguageOut(), src)

	c.Test(
		`let obj = { bar() { return 9; } };
class D {
  e = obj;
  f = this.e.bar() * 4;
}
`,
		`let obj = { bar() { return 9; } };
class D {
  constructor() {
    this.e = obj;
    this.f = this.e.bar() * 4;
  }
}
`)

	c.Test(
		`class Foo {
  y = 'apple';
  x = () => { return this.y + ' and banana'; };
}
`,
		`class Foo {
  constructor() {
    this.y = 'apple';
    this.x = () => { return this.y + ' and banana'; };
  }
}
`)

	c.Test(
		`class Bar {
  x = () => { this.method(); };
  method() {}
}
`,
		`class Bar {
  constructor() {
    this.x = () => { this.method(); };
  }
  method() {}
}
`)
}

// Java: RewriteClassMembersTest.testSuperInNonStaticPublicField()
func (c *RewriteClassMembersTest) TestSuperInNonStaticPublicField() {
	src := `class Foo {
  x() {
    return 3;
  }
}
class Bar extends Foo {
  y = 1 + super.x();
}
`
	c.Test(
		src,
		`class Foo {
  x() {
    return 3;
  }
}
class Bar extends Foo {
  constructor() {
    super(...arguments);
    this.y = 1 + super.x();
  }
}
`)
	c.TestSame(c.WithOptions().UseEs2022LanguageOut(), src)
}

// Java: RewriteClassMembersTest.testThisInStaticField_otherFieldRef()
func (c *RewriteClassMembersTest) TestThisInStaticField_otherFieldRef() {
	src := `class C {
  static x = 1;
  static y = this.x + 1;
}
`

	c.Test(
		src,
		`class C {}
C.x = 1;
C.y = C.x + 1;
`)

	c.Test(
		c.WithOptions().UseStaticInheritance(),
		src,
		`class C {}
C.x = 1;
C.y = C.x + 1;
`)

	c.Test(
		c.WithOptions().UseEs2022LanguageOut(),
		src,
		`class C {
  static x;
  static y;
}
C.x = 1;
C.y = C.x + 1;
`)
}

// Java: RewriteClassMembersTest.testThisInStaticField_thisInArrowFunction()
func (c *RewriteClassMembersTest) TestThisInStaticField_thisInArrowFunction() {
	src := `class C {
  static x = 2;
  static y = () => this.x;
}
`

	c.Test(
		src,
		`class C {}
C.x = 2;
C.y = () => {
  // Note: This is the correct behavior.
  return C.x;
};
`)

	c.Test(
		c.WithOptions().UseStaticInheritance(),
		src,
		`class C {}
C.x = 2;
C.y = () => {
  return C.x;
};
`)

	c.Test(
		c.WithOptions().UseEs2022LanguageOut(),
		src,
		`class C {
  static x;
  static y;
}
C.x = 2;
C.y = () => {
  return C.x;
};
`)
}

// Java: RewriteClassMembersTest.testThisInStaticField_staticMethodCall()
func (c *RewriteClassMembersTest) TestThisInStaticField_staticMethodCall() {
	c.Test(
		`class F {
  static a = 'there';
  static b = this.c() + this.a;
  static c() { return 'hi'; }
}
`,
		`class F {
  static c() {
    return "hi";
  }
}
F.a = "there";
F.b = F.c() + F.a;
`)
}

// Java: RewriteClassMembersTest.testThisInStaticGetter()
func (c *RewriteClassMembersTest) TestThisInStaticGetter() {
	c.TestSame(
		`let x = 1;
class Child {
  static getX() {
    return x;
  }
  static get prop() {
    return this.getX();
  }
}
`)
}

// Java: RewriteClassMembersTest.testThisInStaticSetter()
func (c *RewriteClassMembersTest) TestThisInStaticSetter() {
	c.TestSame(
		`let x = 1;
class Child {
  static setX(newX) {
    x = newX;
  }
  static set prop(p) {
    this.setX(p);
  }
}
`)
}

// Java: RewriteClassMembersTest.testSuperInStaticField()
func (c *RewriteClassMembersTest) TestSuperInStaticField() {
	c.Test(
		`class Foo {
  static x(a, b) {
    return a + b;
  }
  static y(c, d) {
    return c - d;
  }
}
class Bar extends Foo {
  static z = () => super.x(1, 2) + 12 + super.y(3, 4);
}
`,
		`class Foo {
  static x(a, b) {
    return a + b;
  }
  static y(c, d) {
    return c - d;
  }
}
class Bar extends Foo {}
Bar.z = () => {
  return Foo.x(1, 2) + 12 + Foo.y(3, 4);
};
`)

	c.Test(
		`const ns = {};
ns.Foo = class {
  static x(a, b) {
    return 5;
  }
}
class Bar extends ns.Foo {
  static z = () => super.x(1, 2);
}
`,
		`const ns = {};
ns.Foo = class {
  static x(a, b) {
    return 5;
  }
};
class Bar extends ns.Foo {}
Bar.z = () => {
  return ns.Foo.x(1, 2);
};
`)

	c.Test(
		`class Bar {
  static a = { method1() {} };
  static b = { method2() { super.method1(); } };
}
`,
		`class Bar {}
Bar.a = {method1() {}};
Bar.b = {method2() {
  super.method1();
}};
`)

	c.Test(
		`class Parent {
  static get parentGetter() {
    return {val: 1};
  }
}
class Child extends Parent {
  static val = super.parentGetter.val;
}
`,
		`class Parent {
  static get parentGetter() {
    return {val:1};
  }
}
class Child extends Parent {}
Child.val = Parent.parentGetter.val;
`)
}

// Java: RewriteClassMembersTest.testSuperInStaticField_superInArrowFunction()
func (c *RewriteClassMembersTest) TestSuperInStaticField_superInArrowFunction() {
	src := `class Parent {
  static x = 1;
}
class Child extends Parent {
  static y = () => super.x;
}
`

	c.Test(
		src,
		`class Parent {}
Parent.x = 1;
class Child extends Parent {}
Child.y = () => {
  return Parent.x;
};
`)

	c.Test(
		c.WithOptions().UseStaticInheritance(),
		src,
		`class Parent {}
Parent.x = 1;
class Child extends Parent {
  // static_init method stays because of the super reference.
  static STATIC_INIT$0() {
    Child.y = () => {
      return super.x;
    };
  }
}
Child.STATIC_INIT$0();
`)
}

// Java: RewriteClassMembersTest.testSuperReferencesStaticGetter()
func (c *RewriteClassMembersTest) TestSuperReferencesStaticGetter() {
	src := `class Parent {
  static getName() {
    return 'Parent';
  }
  static get greeting() {
    return 'Hello ' + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return 'Child';
  }
  static msg = super.greeting;  // 'Hello Child'
}
`

	// non-strict
	c.Test(
		src,
		`class Parent {
  static getName() {
    return "Parent";
  }
  static get greeting() {
    return "Hello " + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return "Child";
  }
}
Child.msg = Parent.greeting;
`)

	// strict
	c.Test(
		c.WithOptions().UseStaticInheritance(),
		src,
		`class Parent {
  static getName() {
    return "Parent";
  }
  static get greeting() {
    return "Hello " + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return "Child";
  }
  // static_init method stays because of the super reference.
  static STATIC_INIT$0() {
    Child.msg = super.greeting;
  }
}
Child.STATIC_INIT$0();
`)
}

// Java: RewriteClassMembersTest.testSuperReferencesStaticGetter_viaElementAccess()
func (c *RewriteClassMembersTest) TestSuperReferencesStaticGetter_viaElementAccess() {
	src := `/** @unrestricted */
class Parent {
  static getName() {
    return 'Parent';
  }
  static get ['greeting']() {
    return 'Hello ' + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return 'Child';
  }
  static msg = super['greeting'];  // 'Hello Child'
}
`

	// non-strict
	c.Test(
		src,
		`class Parent {
  static getName() {
    return "Parent";
  }
  static get ["greeting"]() {
    return "Hello " + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return "Child";
  }
}
Child.msg = Parent["greeting"];
`)

	// strict
	c.Test(
		c.WithOptions().UseStaticInheritance(),
		src,
		`class Parent {
  static getName() {
    return "Parent";
  }
  static get ["greeting"]() {
    return "Hello " + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return "Child";
  }
  // static_init method stays because of the super reference.
  static STATIC_INIT$0() {
    Child.msg = super["greeting"];
  }
}
Child.STATIC_INIT$0();
`)
}

// Java: RewriteClassMembersTest.testSuperInStaticFieldObjectSpread()
func (c *RewriteClassMembersTest) TestSuperInStaticFieldObjectSpread() {
	c.Test(
		`class Base {
  static X = {a: 1};
}

class Child extends Base {
  /** @type {!Object} */
  static Y = {...super.X, b: 2};
}
`,
		`class Base {}
Base.X = {a:1};
class Child extends Base {}
Child.Y = {...Base.X, b:2};
`)

	c.Test(
		`const ns = {};
ns.Base = class {
  static X = {a: 1};
};

class Child extends ns.Base {
  /** @type {!Object} */
  static Y = {...super.X, b: 2};
}
`,
		`const ns = {};
ns.Base = class {};
ns.Base.X = {a:1};
class Child extends ns.Base {}
Child.Y = {...ns.Base.X, b:2};
`)
}

// Java: RewriteClassMembersTest.testSuperInStaticBlock()
func (c *RewriteClassMembersTest) TestSuperInStaticBlock() {
	c.Test(
		`class Parent {
  static getGreeting() {
    return 'Hello';
  }
}
class Child extends Parent {
  static {
    alert(super.getGreeting());
  }
}
`,
		`class Parent {
  static getGreeting() {
    return "Hello";
  }
}
class Child extends Parent {}
{
  alert(Parent.getGreeting());
}
`)
}

// Java: RewriteClassMembersTest.testSuperInStaticBlock_strictSuperRewrite()
func (c *RewriteClassMembersTest) TestSuperInStaticBlock_strictSuperRewrite() {
	source := `class Parent {
  static getName() {
    return 'Parent';
  }
  static getGreeting() {
    return 'Hello ' + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return 'Child';
  }
  static {
    alert(super.getGreeting());  // Alerts: 'Hello Child'
  }
}
`

	// Test both conditions that trigger strict super rewrite.

	c.Test(
		c.WithOptions().UseStaticInheritance(),
		source,
		`class Parent {
  static getName() {
    return "Parent";
  }
  static getGreeting() {
    return "Hello " + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return "Child";
  }
  // static_init method stays because of the super reference.
  static STATIC_INIT$0() {
    {
      alert(super.getGreeting());
    }
  }
}
Child.STATIC_INIT$0();
`)

	c.Test(
		c.WithOptions().UseEs5LanguageOut(),
		source,
		`class Parent {
  static getName() {
    return "Parent";
  }
  static getGreeting() {
    return "Hello " + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return "Child";
  }
  // static_init method stays because of the super reference.
  static STATIC_INIT$0() {
    {
      alert(super.getGreeting());
    }
  }
}
Child.STATIC_INIT$0();
`)
}

// Java: RewriteClassMembersTest.testSuperInStaticMethod()
func (c *RewriteClassMembersTest) TestSuperInStaticMethod() {
	c.Test(
		`class Parent {
  static getGreeting() {
    return 'Hello';
  }
}
class Child extends Parent {
  static sayHello() {
    alert(super.getGreeting());
  }
}
Child.sayHello();
`,
		`class Parent {
  static getGreeting() {
    return "Hello";
  }
}
class Child extends Parent {
  static sayHello() {
    alert(Parent.getGreeting());
  }
}
Child.sayHello();
`)
}

// Java: RewriteClassMembersTest.testSuperInStaticMethod_strictSuperCallRewrite()
func (c *RewriteClassMembersTest) TestSuperInStaticMethod_strictSuperCallRewrite() {
	source := `class Parent {
  static getName() {
    return 'Parent';
  }
  static getGreeting() {
    return 'Hello ' + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return 'Child';
  }
  static sayHello() {
    alert(super.getGreeting());  // Alerts: 'Hello Child'
  }
}
Child.sayHello();
`

	// Test both conditions that trigger strict super rewrite.

	c.Test(
		c.WithOptions().UseStaticInheritance(),
		source,
		`class Parent {
  static getName() {
    return "Parent";
  }
  static getGreeting() {
    return "Hello " + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return "Child";
  }
  static sayHello() {
    alert(super.getGreeting());
  }
}
Child.sayHello();
`)

	c.Test(
		c.WithOptions().UseEs5LanguageOut(),
		source,
		`class Parent {
  static getName() {
    return "Parent";
  }
  static getGreeting() {
    return "Hello " + this.getName();
  }
}
class Child extends Parent {
  static getName() {
    return "Child";
  }
  static sayHello() {
    alert(super.getGreeting());
  }
}
Child.sayHello();
`)
}

// Java: RewriteClassMembersTest.testSuperInComputedStaticMethod()
func (c *RewriteClassMembersTest) TestSuperInComputedStaticMethod() {
	c.Test(
		`class Parent {
  static getGreeting() {
    return 'Hello';
  }
}
/** @unrestricted */
class Child extends Parent {
  static ['sayHello']() {
    alert(super.getGreeting());
  }
}
Child['sayHello']();
`,
		`class Parent {
  static getGreeting() {
    return "Hello";
  }
}
class Child extends Parent {
  static ["sayHello"]() {
    alert(Parent.getGreeting());
  }
}
Child["sayHello"]();
`)
}

// Java: RewriteClassMembersTest.testSuperInStaticGetter()
func (c *RewriteClassMembersTest) TestSuperInStaticGetter() {
	c.Test(
		`class Parent {
  static getVal() {
    return 1;
  }
}
class Child extends Parent {
  static get childProp() {
    return super.getVal();
  }
}
`,
		`class Parent {
  static getVal() {
    return 1;
  }
}
class Child extends Parent {
  static get childProp() {
    return Parent.getVal();
  }
}
`)
}

// Java: RewriteClassMembersTest.testSuperInStaticSetter()
func (c *RewriteClassMembersTest) TestSuperInStaticSetter() {
	c.Test(
		`class Parent {
  static getVal() {
    return 1;
  }
}
class Child extends Parent {
  static set childProp(x) {
    alert(super.getVal());
  }
}
`,
		`class Parent {
  static getVal() {
    return 1;
  }
}
class Child extends Parent {
  static set childProp(x) {
    alert(Parent.getVal());
  }
}
`)
}

// Java: RewriteClassMembersTest.testComputedPropInNonStaticField()
func (c *RewriteClassMembersTest) TestComputedPropInNonStaticField() {
	src := `/** @unrestricted */
class C {
  [x+=1];
  [x+=2] = 3;
}
`
	c.Test(
		src,
		`var COMP_FIELD$0 = x = x + 1;
var COMP_FIELD$1 = x = x + 2;
class C {
  constructor() {
    this[COMP_FIELD$0] = void 0;
    this[COMP_FIELD$1] = 3;
  }
}
`)
	c.Test(
		c.WithOptions().UseEs2022LanguageOut(),
		src,
		`var COMP_FIELD$0 = x = x + 1;
var COMP_FIELD$1 = x = x + 2;
class C {
  [COMP_FIELD$0];
  [COMP_FIELD$1] = 3;
}
`)

	c.Test(
		`/** @unrestricted */
class C {
  [1] = 1;
  /** @suppress {partialAlias} */
  [2] = this[1];
}
`,
		`class C {
  constructor() {
    this[1] = 1;
    /** @suppress {partialAlias} */
    this[2] = this[1];
  }
}
`)

	c.Test(
		`/** @unrestricted */
let c = class C {
  static [1] = 2;
  [2] = C[1]
}
`,
		`let c = class {
  constructor() {
    this[2] = c[1];
  }
};
c[1] = 2;
`)

	c.Test(
		`foo(/** @unrestricted */ class C {
  static [1] = 2;
  [2] = C[1]
})
`,
		`const CLASS_DECL$0 = class {
  constructor() {
    this[2] = CLASS_DECL$0[1];
  }
};
CLASS_DECL$0[1] = 2;
foo(CLASS_DECL$0);
`)

	c.Test(
		`let c = class {
  x = 1
  y = this.x
}
/** @unrestricted */
class B {
  [1] = 2;
  [2] = this[1]
}
`,
		`let c = class {
  constructor() {
    this.x = 1;
    this.y = this.x;
  }
};
class B {
  constructor() {
    this[1] = 2;
    this[2] = this[1];
  }
}
`)

	c.TestSame(
		`class Clazz {
  [Symbol.toPrimitive]() {
    return 42;
  }
}
`)
}

// Java: RewriteClassMembersTest.testProtoField_instance()
func (c *RewriteClassMembersTest) TestProtoField_instance() {
	src :=
		`/** @unrestricted */
class C {
  __proto__ = {a: 1};
  ['__proto__'] = {b: 2};
  '__proto__' = {c: 3};
}
`
	c.Test(
		src,
		`class C {
  constructor() {
    $jscomp.global.Object.defineProperty(this, "__proto__", {
      configurable: true,
      enumerable: true,
      writable: true,
      value: {a: 1}
    });
    $jscomp.global.Object.defineProperty(this, "__proto__", {
      configurable: true,
      enumerable: true,
      writable: true,
      value: {b: 2}
    });
    $jscomp.global.Object.defineProperty(this, "__proto__", {
      configurable: true,
      enumerable: true,
      writable: true,
      value: {c: 3}
    });
  }
}
`)
	c.Test(
		c.WithOptions().UseEs2022LanguageOut(),
		src,
		`class C {
  __proto__ = {a: 1};
  ["__proto__"] = {b: 2};
  ["__proto__"] = {c: 3};
}
`)
}

// Java: RewriteClassMembersTest.testProtoField_uninitializedInstance()
func (c *RewriteClassMembersTest) TestProtoField_uninitializedInstance() {
	src :=
		`/** @unrestricted */
class C {
  __proto__;
  ['__proto__'];
}
`
	c.Test(
		src,
		`class C {
  constructor() {
    $jscomp.global.Object.defineProperty(this, "__proto__", {
      configurable: true,
      enumerable: true,
      writable: true,
      value: void 0
    });
    $jscomp.global.Object.defineProperty(this, "__proto__", {
      configurable: true,
      enumerable: true,
      writable: true,
      value: void 0
    });
  }
}
`)
	c.Test(
		c.WithOptions().UseEs2022LanguageOut(),
		src,
		`class C {
  __proto__;
  ["__proto__"];
}
`)
}

// Java: RewriteClassMembersTest.testProtoField_static()
func (c *RewriteClassMembersTest) TestProtoField_static() {
	src :=
		`/** @unrestricted */
class C {
  static __proto__ = {a: 1};
  static ['__proto__'] = {b: 2};
}
`
	c.Test(
		src,
		`class C {}
$jscomp.global.Object.defineProperty(C, "__proto__", {
  configurable: true,
  enumerable: true,
  writable: true,
  value: {a: 1}
});
$jscomp.global.Object.defineProperty(C, "__proto__", {
  configurable: true,
  enumerable: true,
  writable: true,
  value: {b: 2}
});
`)
	c.Test(
		c.WithOptions().UseEs2022LanguageOut(),
		src,
		`class C {
  static __proto__;
  static ["__proto__"];
}
$jscomp.global.Object.defineProperty(C, "__proto__", {
  configurable: true,
  enumerable: true,
  writable: true,
  value: {a: 1}
});
$jscomp.global.Object.defineProperty(C, "__proto__", {
  configurable: true,
  enumerable: true,
  writable: true,
  value: {b: 2}
});
`)
}

// Java: RewriteClassMembersTest.testProtoField_uninitializedStatic()
func (c *RewriteClassMembersTest) TestProtoField_uninitializedStatic() {
	src :=
		`/** @unrestricted */
class C {
  static __proto__;
  static ['__proto__'];
}
`
	c.Test(
		src,
		`class C {}
$jscomp.global.Object.defineProperty(C, "__proto__", {
  configurable: true,
  enumerable: true,
  writable: true,
  value: void 0
});
$jscomp.global.Object.defineProperty(C, "__proto__", {
  configurable: true,
  enumerable: true,
  writable: true,
  value: void 0
});
`)
	c.Test(
		c.WithOptions().UseEs2022LanguageOut(),
		src,
		`class C {
  static __proto__;
  static ["__proto__"];
}
`)
}

// Java: RewriteClassMembersTest.testProtoField_interleavedWithNormalFields()
func (c *RewriteClassMembersTest) TestProtoField_interleavedWithNormalFields() {
	c.Test(
		`class C {
  x = 1;
  __proto__ = 2;
  y = 3;
  static a = 4;
  static __proto__ = 5;
  static b = 6;
}
`,
		`class C {
  constructor() {
    this.x = 1;
    $jscomp.global.Object.defineProperty(this, "__proto__", {
      configurable: true,
      enumerable: true,
      writable: true,
      value: 2
    });
    this.y = 3;
  }
}
C.a = 4;
$jscomp.global.Object.defineProperty(C, "__proto__", {
  configurable: true,
  enumerable: true,
  writable: true,
  value: 5
});
C.b = 6;
`)
}

// Java: RewriteClassMembersTest.testComputedPropInStaticField()
func (c *RewriteClassMembersTest) TestComputedPropInStaticField() {
	src := `/** @unrestricted */
class C {
  static ['x'];
  static ['y'] = 2;
}
`
	c.Test(
		src,
		`class C {}
C["x"] = void 0;
C["y"] = 2;
`)
	c.Test(
		c.WithOptions().UseEs2022LanguageOut(),
		src,
		`class C {
  static ["x"];
  static ["y"];
}
C["y"] = 2;
`)

	c.Test(
		`/** @unrestricted */
class C {
  static [1] = 1;
  static [2] = this[1];
}
`,
		`class C {}
C[1] = 1;
C[2] = C[1];
`)

	c.Test(
		`/** @unrestricted */
const C = class {
  static [1] = 1;
  static [2] = this[1];
}
`,
		`const C = class {};
C[1] = 1;
C[2] = C[1];
`)

	c.Test(
		`/** @unrestricted */
const C = class InnerC {
  static [1] = 1;
  static [2] = this[1];
  static [3] = InnerC[2];
}
`,
		`const C = class {};
C[1] = 1;
C[2] = C[1];
C[3] = C[2];
`)

	c.Test(
		`/** @unrestricted */
let c = class C {
  static [1] = 2;
  static [2] = C[1]
}
`,
		`let c = class {};
c[1] = 2;
c[2] = c[1];
`)

	c.Test(
		`foo(/** @unrestricted */ class C {
  static [1] = 2;
  static [2] = C[1]
})
`,
		`const CLASS_DECL$0 = class {};
CLASS_DECL$0[1] = 2;
CLASS_DECL$0[2] = CLASS_DECL$0[1];
foo(CLASS_DECL$0);
`)

	c.Test(
		`foo(/** @unrestricted */ class {
  static [1] = 1
})
`,
		`const CLASS_DECL$0 = class {};
CLASS_DECL$0[1] = 1;
foo(CLASS_DECL$0);
`)

	c.TestSame(
		`class Clazz {
  static [Symbol.hasInstance](x) {
    return false;
  }
}
`)
}

// Java: RewriteClassMembersTest.testSideEffectsInComputedField()
func (c *RewriteClassMembersTest) TestSideEffectsInComputedField() {
	c.Test(
		`function bar() {
  this.x = 3;
  /** @unrestricted */
  class Foo {
    y;
    [this.x] = 2;
  }
}
`,
		`function bar() {
  this.x = 3;
  var COMP_FIELD$0 = this.x;
  class Foo {
    constructor() {
      this.y = void 0;
      this[COMP_FIELD$0] = 2;
    }
  }
}
`)

	c.Test(
		`class E {
  y() { return 1; }
}
class F extends E {
  x() {
    return /** @unrestricted */ class {
      [super.y()] = 4;
    }
  }
}
`,
		`class E {
  y() {
    return 1;
  }
}
class F extends E {
  x() {
    var COMP_FIELD$1 = super.y();
    const CLASS_DECL$0 = class {
      constructor() {
        this[COMP_FIELD$1] = 4;
      }
    };
    return CLASS_DECL$0;
  }
}
`)

	c.Test(
		`function bar(num) {}
/** @unrestricted */
class Foo {
  [bar(1)] = 'a';
  static b = bar(3);
  static [bar(2)] = bar(4);
}
`,
		`function bar(num) {}
var COMP_FIELD$0 = bar(1);
var COMP_FIELD$1 = bar(2);
class Foo {
  constructor() {
    this[COMP_FIELD$0] = "a";
  }
}
Foo.b = bar(3);
Foo[COMP_FIELD$1] = bar(4);
`)

	c.Test(
		`let x = 'hello';
/** @unrestricted */ class Foo {
  static n = (x=5);
  static [x] = 'world';
}
`,
		`let x = "hello";
var COMP_FIELD$0 = x;
class Foo {}
Foo.n = x = 5;
Foo[COMP_FIELD$0] = "world";
`)

	c.Test(
		`function foo(num) {}
/** @unrestricted */
class Baz {
  ['f' + foo(1)];
  static x = foo(6);
  ['m' + foo(2)]() {};
  static [foo(3)] = foo(7);
  [foo(4)] = 2;
  get [foo(5)]() {}
}
`,
		`function foo(num) {}
var COMP_FIELD$0 = "f" + foo(1);
var COMP_FIELD$1 = "m" + foo(2);
var COMP_FIELD$2 = foo(3);
var COMP_FIELD$3 = foo(4);
var COMP_FIELD$4 = foo(5);
class Baz {
  constructor() {
    this[COMP_FIELD$0] = void 0;
    this[COMP_FIELD$3] = 2;
  }
  [COMP_FIELD$1]() {}
  get [COMP_FIELD$4]() {}
}
Baz.x = foo(6);
Baz[COMP_FIELD$2] = foo(7);
`)
}

// Java: RewriteClassMembersTest.testClassStaticBlocksNoFieldAssign()
func (c *RewriteClassMembersTest) TestClassStaticBlocksNoFieldAssign() {
	c.Test(
		`class C {
  static {
  }
}
`,
		`class C {
}
`)

	c.Test(
		`class C {
  static {
    let x = 2
    const y = x
  }
}
`,
		`class C {}
{
  let x = 2;
  const y = x;
}
`)

	c.Test(
		`class C {
  static {
    let x = 2
    const y = x
    let z;
    if (x - y == 0) {z = 1} else {z = 2}
    while (x - z > 10) {z++;}
    for (;;) {break;}
  }
}
`,
		`class C {}
{
  let x = 2;
  const y = x;
  let z;
  if (x - y == 0) {
    z = 1;
  } else {
    z = 2;
  }
  for (; x - z > 10;) {
    z++;
  }
  for (;;) {
    break;
  }
}
`)

	c.Test(
		`class C {
  static {
    let x = 2
  }
  static {
    const y = x
  }
}
`,
		`class C {}
{
  let x = 2;
}
{
  const y = x;
}
`)

	c.Test(
		`class C {
  static {
    let x = 2
  }
  static {
    const y = x
  }
}
class D {
  static {
    let z = 1
  }
}
`,
		`class C {}
{
  let x = 2;
}
{
  const y = x;
}
class D {}
{
  let z = 1;
}
`)

	c.Test(
		`class C {
  static {
    let x = function () {return 1;}
    const y = () => {return 2;}
    function a() {return 3;}
    let z = (() => {return 4;})();
  }
}
`,
		`class C {}
{
  function a() {
    return 3;
  }
  let x = function() {
    return 1;
  };
  const y = () => {
    return 2;
  };
  let z = (() => {
    return 4;
  })();
}
`)

	c.Test(
		`class C {
  static {
    C.x = 2
    const y = C.x;
  }
}
`,
		`class C {}
{
  C.x = 2;
  const y = C.x;
}
`)

	c.Test(
		`class Foo {
  static {
    let x = 5;
    class Bar {
      static {
        let x = 'str';
      }
    }
  }
}
`,
		`class Foo {}
{
  let x = 5;
  class Bar {}
  {
    let x$jscomp$1 = "str";
  }
}
`)
}

// Java: RewriteClassMembersTest.testStaticNoncomputed()
func (c *RewriteClassMembersTest) TestStaticNoncomputed() {
	c.Test(
		`class C {
  static x = 2
}
`,
		`class C {}
C.x = 2;
`)

	src := `class C {
  static x;
}
`
	c.Test(
		src,
		`class C {}
C.x = void 0;
`)
	c.TestSame(c.WithOptions().UseEs2022LanguageOut(), src)

	src =
		`class C {
  static x = 2
  static y = 'hi'
  static z;
}
`
	c.Test(
		src,
		`class C {}
C.x = 2;
C.y = "hi";
C.z = void 0;
`)
	c.Test(
		c.WithOptions().UseEs2022LanguageOut(),
		src,
		`class C {
  static x;
  static y;
  static z;
}
C.x = 2;
C.y = "hi";
`)

	c.Test(
		`class C {
  static x = 2
  static y = 3
}
class D {
  static z = 1
}
`,
		`class C {}
C.x = 2;
C.y = 3;
class D {}
D.z = 1;
`)

	c.Test(
		`class C {
  static w = function () {return 1;};
  static x = () => {return 2;};
  static y = (function a() {return 3;})();
  static z = (() => {return 4;})();
}
`,
		`class C {}
C.w = function() {
  return 1;
};
C.x = () => {
  return 2;
};
C.y = function a() {
  return 3;
}();
C.z = (() => {
  return 4;
})();
`)

	c.Test(
		`class C {
  static x = 2
  static y = C.x
}
`,
		`class C {}
C.x = 2;
C.y = C.x;
`)

	c.Test(
		`class C {
  static x = 2
  static {let y = C.x}
}
`,
		`class C {}
C.x = 2;
{
  let y = C.x;
}
`)
}

// Java: RewriteClassMembersTest.testInstanceNoncomputedWithNonemptyConstructor()
func (c *RewriteClassMembersTest) TestInstanceNoncomputedWithNonemptyConstructor() {
	c.Test(
		`class C extends Object {
  x = 1;
  z = 3;
  constructor() {
    super();
    this.y = 2;
  }
}
`,
		`class C extends Object{
  constructor() {
    super();
    this.x = 1
    this.z = 3
    this.y = 2;
  }
}
`)

	c.Test(
		`class C {
  x;
  constructor() {
    this.y = 2;
  }
}
`,
		`class C {
  constructor() {
    this.x = void 0;
    this.y = 2;
  }
}
`)

	c.Test(
		`class C {
  x = 1
  y = 2
  constructor() {
    this.z = 3;
  }
}
`,
		`class C {
  constructor() {
    this.x = 1;
    this.y = 2;
    this.z = 3;
  }
}
`)

	c.Test(
		`class C {
  x = 1
  y = 2
  constructor() {
    alert(3);
    this.z = 4;
  }
}
`,
		`class C {
  constructor() {
    this.x = 1;
    this.y = 2;
    alert(3);
    this.z = 4;
  }
}
`)

	c.Test(
		`class C {
  x = 1
  constructor() {
    alert(3);
    this.z = 4;
  }
  y = 2
}
`,
		`class C {
  constructor() {
    this.x = 1;
    this.y = 2;
    alert(3);
    this.z = 4;
  }
}
`)

	c.Test(
		`class C {
  x = 1
  constructor() {
    alert(3);
    this.z = 4;
  }
  y = 2
}
class D {
  a = 5;
  constructor() { this.b = 6;}
}
`,
		`class C {
  constructor() {
    this.x = 1;
    this.y = 2;
    alert(3);
    this.z = 4;
  }
}
class D {
constructor() {
  this.a = 5;
  this.b = 6
}
}
`)
}

// Java: RewriteClassMembersTest.testInstanceComputedWithNonemptyConstructorAndSuper()
func (c *RewriteClassMembersTest) TestInstanceComputedWithNonemptyConstructorAndSuper() {
	src := `class A { constructor() { alert(1); } }
/** @unrestricted */ class C extends A {
  ['x'] = 1;
  constructor() {
    super();
    this['y'] = 2;
    this['z'] = 3;
  }
}
`
	c.Test(
		src,
		`class A { constructor() { alert(1); } }
class C extends A {
  constructor() {
    super()
    this['x'] = 1
    this['y'] = 2;
    this['z'] = 3;
  }
}
`)
	c.Test(
		c.WithOptions().UseEs2022LanguageOut(),
		src,
		`class A {
  constructor() {
    alert(1);
  }
}
class C extends A {
  ["x"] = 1;
  constructor() {
    super();
    this["y"] = 2;
    this["z"] = 3;
  }
}
`)
}

// Java: RewriteClassMembersTest.testInstanceNoncomputedWithNonemptyConstructorAndSuper()
func (c *RewriteClassMembersTest) TestInstanceNoncomputedWithNonemptyConstructorAndSuper() {
	c.Test(
		`class A { constructor() { alert(1); } }
class C extends A {
  x = 1;
  constructor() {
    super()
    this.y = 2;
  }
}
`,
		`class A { constructor() { alert(1); } }
class C extends A {
  constructor() {
    super()
    this.x = 1
    this.y = 2;
  }
}
`)

	c.Test(
		`class A { constructor() { this.x = 1; } }
class C extends A {
  y;
  constructor() {
    super()
    alert(3);
    this.z = 4;
  }
}
`,
		`class A {
  constructor() {
    this.x = 1;
  }
}
class C extends A {
  constructor() {
    super();
    this.y = void 0;
    alert(3);
    this.z = 4;
  }
}
`)

	c.Test(
		`class A { constructor() { this.x = 1; } }
class C extends A {
  y;
  constructor() {
    alert(3);
    super()
    this.z = 4;
  }
}
`,
		`class A {
  constructor() {
    this.x = 1;
  }
}
class C extends A {
  constructor() {
    alert(3);
    super();
    this.y = void 0;
    this.z = 4;
  }
}
`)
}

// Java: RewriteClassMembersTest.testNonComputedInstanceWithEmptyConstructor()
func (c *RewriteClassMembersTest) TestNonComputedInstanceWithEmptyConstructor() {
	c.Test(
		`class C {
  x = 2;
  constructor() {}
}
`,
		`class C {
  constructor() {
    this.x = 2;
  }
}
`)

	c.Test(
		`class C {
  x;
  constructor() {}
}
`,
		`class C {
  constructor() {
    this.x = void 0;
  }
}
`)

	c.Test(
		`class C {
  x = 2
  y = 'hi'
  z;
  constructor() {}
}
`,
		`class C {
  constructor() {
    this.x = 2;
    this.y = "hi";
    this.z = void 0;
  }
}
`)

	c.Test(
		`class C {
  x = 1
  constructor() {
  }
  y = 2
}
`,
		`class C {
  constructor() {
    this.x = 1;
    this.y = 2;
  }
}
`)

	c.Test(
		`class C {
  x = 1
  constructor() {
  }
  y = 2
}
class D {
  a = 5;
  constructor() {}
}
`,
		`class C {
  constructor() {
    this.x = 1;
    this.y = 2;
  }
}
class D {
constructor() {
  this.a = 5;
}
}
`)

	c.Test(
		`class C {
  w = function () {return 1;};
  x = () => {return 2;};
  y = (function a() {return 3;})();
  z = (() => {return 4;})();
  constructor() {}
}
`,
		`class C {
  constructor() {
    this.w = function () {return 1;};
    this.x = () => {return 2;};
    this.y = (function a() {return 3;})();
    this.z = (() => {return 4;})();
  }
}
`)

	c.Test(
		`class C {
  static x = 2
  constructor() {}
  y = C.x
}
`,
		`class C {
  constructor() {
    this.y = C.x;
  }
}
C.x = 2;
`)
}

// Java: RewriteClassMembersTest.testInstanceNoncomputedNoConstructor()
func (c *RewriteClassMembersTest) TestInstanceNoncomputedNoConstructor() {
	c.Test(
		`class C {
  x = 2;
}
`,
		`class C {
  constructor() {this.x=2;}
}
`)

	c.Test(
		`class C {
  x;
}
`,
		`class C {
  constructor() {
    this.x = void 0;
  }
}
`)

	c.Test(
		`class C {
  x = 2
  y = 'hi'
  z;
}
`,
		`class C {
  constructor() {
    this.x = 2;
    this.y = "hi";
    this.z = void 0;
  }
}
`)
	c.Test(
		`class C {
  foo() {}
  x = 1;
}
`,
		`class C {
  constructor() {this.x = 1;}
  foo() {}
}
`)

	c.Test(
		`class C {
  static x = 2
  y = C.x
}
`,
		`class C {
  constructor() {
    this.y = C.x;
  }
}
C.x = 2;
`)

	c.Test(
		`class C {
  w = function () {return 1;};
  x = () => {return 2;};
  y = (function a() {return 3;})();
  z = (() => {return 4;})();
}
`,
		`class C {
  constructor() {
    this.w = function () {return 1;};
    this.x = () => {return 2;};
    this.y = (function a() {return 3;})();
    this.z = (() => {return 4;})();
  }
}
`)
}

// Java: RewriteClassMembersTest.testInstanceNonComputedNoConstructorWithSuperclass()
func (c *RewriteClassMembersTest) TestInstanceNonComputedNoConstructorWithSuperclass() {
	c.Test(
		`class B {}
class C extends B {x = 1;}
`,
		`class B {}
class C extends B {
  constructor() {
    super(...arguments);
    this.x = 1;
  }
}
`)
	c.Test(
		`class B {constructor() {}; y = 2;}
class C extends B {x = 1;}
`,
		`class B {constructor() {this.y = 2}}
class C extends B {
  constructor() {
    super(...arguments);
    this.x = 1;
  }
}
`)
	c.Test(
		`class B {constructor(a, b) {}; y = 2;}
class C extends B {x = 1;}
`,
		`class B {constructor(a, b) {this.y = 2}}
class C extends B {
  constructor() {
    super(...arguments);
    this.x = 1;
  }
}
`)
}

// Java: RewriteClassMembersTest.testClassExpressionsStaticBlocks()
func (c *RewriteClassMembersTest) TestClassExpressionsStaticBlocks() {
	c.Test(
		`let c = class C {
  static {
    C.y = 2;
    let x = C.y
  }
}
`,
		`let c = class {};
{
  c.y = 2;
  let x = c.y;
}
`)

	c.Test(
		`foo(class C {
  static {
    C.y = 2;
    let x = C.y
  }
})
`,
		`foo((() => {
  const CLASS_DECL$0 = class {};
  {
    CLASS_DECL$0.y = 2;
    let x = CLASS_DECL$0.y;
  }
  return CLASS_DECL$0;
})());
`)

	c.Test(
		`class A { static b = {}; }
foo(A.b.c = class C {
  static {
    C.y = 2;
    let x = C.y
  }
})
`,
		`class A {}
A.b = {};
foo(A.b.c = (() => {
  const CLASS_DECL$0 = class {};
  {
    CLASS_DECL$0.y = 2;
    let x = CLASS_DECL$0.y;
  }
  return CLASS_DECL$0;
})());
`)
}

// Java: RewriteClassMembersTest.testNonClassDeclarationsStaticBlocks()
func (c *RewriteClassMembersTest) TestNonClassDeclarationsStaticBlocks() {
	c.Test(
		`let c = class {
  static {
    let x = 1
  }
}
`,
		`let c = class {};
{
  let x = 1;
}
`)

	c.Test(
		`class A {}
A.c = class {
  static {
    let x = 1
  }
}
`,
		`class A {}
A.c = class {};
{
  let x = 1;
}
`)

	c.Test(
		`class A {}
A[1] = class {
  static {
    let x = 1
  }
}
`,
		`class A {}
A[1] = (() => {
  const CLASS_DECL$0 = class {};
  {
    let x = 1;
  }
  return CLASS_DECL$0;
})();
`)
}

// Java: RewriteClassMembersTest.testNonClassDeclarationsStaticNoncomputedFields()
func (c *RewriteClassMembersTest) TestNonClassDeclarationsStaticNoncomputedFields() {
	c.Test(
		`let c = class {
  static x = 1
}
`,
		`let c = class {};
c.x = 1;
`)

	c.Test(
		`class A {}
A.c = class {
  static x = 1
}
`,
		`class A {}
A.c = class {};
A.c.x = 1;
`)

	c.Test(
		`class A {}
A[1] = class {
  static x = 1
}
`,
		`class A {}
const CLASS_DECL$0 = class {};
CLASS_DECL$0.x = 1;
A[1] = CLASS_DECL$0;
`)

	c.Test(
		`let c = class C {
  static y = 2;
  static x = C.y
}
`,
		`let c = class {};
c.y = 2;
c.x = c.y;
`)

	c.Test(
		`foo(class C {
  static y = 2;
  static x = C.y
})
`,
		`const CLASS_DECL$0 = class {};
CLASS_DECL$0.y = 2;
CLASS_DECL$0.x = CLASS_DECL$0.y;
foo(CLASS_DECL$0);
`)

	c.Test(
		`foo(class C {
  static y = 2;
  x = C.y
})
`,
		`const CLASS_DECL$0 = class {
  constructor() {
    this.x = CLASS_DECL$0.y;
  }
};
CLASS_DECL$0.y = 2;
foo(CLASS_DECL$0);
`)
}

// Java: RewriteClassMembersTest.testNonClassDeclarationsInstanceNoncomputedFields()
func (c *RewriteClassMembersTest) TestNonClassDeclarationsInstanceNoncomputedFields() {
	c.Test(
		`let c = class {
  y = 2;
}
`,
		`let c = class {
  constructor() {
    this.y = 2;
  }
}
`)

	c.Test(
		`let c = class C {
  y = 2;
}
`,
		`let c = class {
  constructor() {
    this.y = 2;
  }
};
`)

	c.Test(
		`class A {}
A.c = class {
  y = 2;
}
`,
		`class A {}
A.c = class {
  constructor() {
    this.y = 2;
  }
}
`)

	c.Test(
		`A[1] = class {
  y = 2;
}
`,
		`const CLASS_DECL$0 = class {
  constructor() {
    this.y = 2;
  }
};
A[1] = CLASS_DECL$0;
`)

	c.Test(
		`let c = class C {
  y = 2;
}
`,
		`let c = class {
  constructor() {
    this.y = 2;
  }
};
`)

	c.Test(
		`class A {}
A.c = class C {
  y = 2;
}
`,
		`class A {}
A.c = class {
  constructor() {
    this.y = 2;
  }
};

`)

	c.Test(
		`A[1] = class C {
  y = 2;
}
`,
		`const CLASS_DECL$0 = class {
  constructor() {
    this.y = 2;
  }
};
A[1] = CLASS_DECL$0;
`)

	c.Test(
		`foo(class C {
  y = 2;
})
`,
		`const CLASS_DECL$0 = class {
  constructor() {
    this.y = 2;
  }
};
foo(CLASS_DECL$0);
`)
}

// Java: RewriteClassMembersTest.testConstuctorAndStaticFieldDontConflict()
func (c *RewriteClassMembersTest) TestConstuctorAndStaticFieldDontConflict() {
	c.Test(
		`let x = 2;
class C {
  static y = x
  constructor(x) {}
}
`,
		`let x = 2;
class C {
  constructor(x$jscomp$1) {}
}
C.y = x;
`)
}

// Java: RewriteClassMembersTest.testInstanceInitializerShadowsConstructorDeclaration()
func (c *RewriteClassMembersTest) TestInstanceInitializerShadowsConstructorDeclaration() {
	c.Test(
		`let x = 2;
class C {
  y = x;
  constructor(x) {}
}
`,
		`let x = 2;
class C {
  constructor(x$jscomp$1) {
    this.y = x;
  }
}
`)

	c.Test(
		`let x = 2;
class C {
  y = x;
  constructor() { let x; }
}
`,
		`let x = 2;
class C {
  constructor() {
    this.y = x;
    let x$jscomp$1;
  }
}
`)

	c.Test(
		`let x = 2;
class C {
  y = x
  constructor() { {var x;} }
}
`,
		`let x = 2;
class C {
  constructor() {
    this.y = x;
    {
     var x$jscomp$1;
    }
  }
}
`)

	c.Test(
		`function f() { return 4; }
class C {
  y = f();
  constructor() {function f() { return 'str'; }}
}
`,
		`function f() {
  return 4;
}
class C {
  constructor() {
    function f$jscomp$1() {
      return 'str';
    }
    this.y = f();
  }
}
`)

	c.Test(
		`class Foo {
  constructor(x) {}
  y = (x) => x;
}
`,
		`class Foo {
  constructor(x) {
    this.y = x$jscomp$1 => {
      return x$jscomp$1;
    };
  }
}
`)

	c.Test(
		`let x = 2;
class C {
  y = (x) => x;
  constructor(x) {}
}
`,
		`let x = 2;
class C {
  constructor(x$jscomp$2) {
    this.y = x$jscomp$1 => {
      return x$jscomp$1;
    };
  }
}
`)
}

// Java: RewriteClassMembersTest.testInstanceInitializerDoesntShadowConstructorDeclaration()
func (c *RewriteClassMembersTest) TestInstanceInitializerDoesntShadowConstructorDeclaration() {
	c.Test(
		`let x = 2;
class C {
  y = x;
  constructor() { {let x;} }
}
`,
		`let x = 2;
class C {
  constructor() {
    this.y = x;
    {let x$jscomp$1;}
  }
}
`)

	c.Test(
		`let x = 2;
class C {
  y = x
  constructor() {() => { let x; };}
}
`,
		`let x = 2;
class C {
  constructor() {
    this.y = x;
    () => { let x$jscomp$1; };
  }
}
`)

	c.Test(
		`let x = 2;
class C {
  y = x
  constructor() {(x) => 3;}
}
`,
		`let x = 2;
class C {
  constructor() {
    this.y = x;
    (x$jscomp$1) => { return 3; };
  }
}
`)
}

// Java: RewriteClassMembersTest.testInstanceFieldInitializersDontBleedOut()
func (c *RewriteClassMembersTest) TestInstanceFieldInitializersDontBleedOut() {
	c.Test(
		`class C {
  y = z
  method() { x; }
  constructor(x) {}
}
`,
		`class C {
  method() { x; }
  constructor(x) {
    this.y = z;
  }
}
`)
}

// Java: RewriteClassMembersTest.testNestedClassesWithShadowingInstanceFields()
func (c *RewriteClassMembersTest) TestNestedClassesWithShadowingInstanceFields() {
	c.Test(
		`let x = 2;
class C {
  y = () => {
    class Foo { z = x }
  };
  constructor(x) {}
}
`,
		`let x = 2;
class C {
  constructor(x$jscomp$1) {
    this.y = () => {
      class Foo {
        constructor() {
          this.z = x;
        }
      }
    };
  }
}
`)
}

// Added when fixing transpilation of real-world code that passed a class expression to a
// constructor call.
//
// Java: RewriteClassMembersTest.testPublicFieldsInClassExpressionInNew()
func (c *RewriteClassMembersTest) TestPublicFieldsInClassExpressionInNew() {
	c.Test(
		`let foo = new (
    class Bar {
      x;
      static y;
    }
)();
`,
		`const CLASS_DECL$0 = class {
  constructor() {
    this.x = void 0;
  }
};
CLASS_DECL$0.y = void 0;
let foo = new CLASS_DECL$0();
`)
}

// Java: RewriteClassMembersTest.testNonClassDeclarationsFunctionArgs()
func (c *RewriteClassMembersTest) TestNonClassDeclarationsFunctionArgs() {
	c.Test(
		`A[foo()] = class {
  static x;
}
`,
		`A[foo()] = (() => {
  const CLASS_DECL$0 = class {};
  CLASS_DECL$0.x = void 0;
  return CLASS_DECL$0;
})();
`)

	c.Test(
		`foo(c = class {
  static x;
})
`,
		`const CLASS_DECL$0 = class {};
CLASS_DECL$0.x = void 0;
foo(c = CLASS_DECL$0);
`)

	c.Test(
		`function foo(c = class {
  static x;
}) {}
`,
		`function foo(c = (() => {
  const CLASS_DECL$0 = class {};
  CLASS_DECL$0.x = void 0;
  return CLASS_DECL$0;
})()) {}
`)
}

// Java: RewriteClassMembersTest.testAnonymousClassExpression()
func (c *RewriteClassMembersTest) TestAnonymousClassExpression() {
	c.Test(
		`function foo() {
  return class {
    y;
    static x;
  }
}
`,
		`function foo() {
  const CLASS_DECL$0 = class {
    constructor() {
      this.y = void 0;
    }
  };
  CLASS_DECL$0.x = void 0;
  return CLASS_DECL$0;
}
`)

	c.Test(
		`foo(class {
  y = 2;
})
`,
		`const CLASS_DECL$0 = class {
  constructor() {
    this.y = 2;
  }
};
foo(CLASS_DECL$0);
`)

	c.Test(
		`foo(class {
  static x = 1;
})
`,
		`const CLASS_DECL$0 = class {};
CLASS_DECL$0.x = 1;
foo(CLASS_DECL$0);
`)
}

// Java: RewriteClassMembersTest.testNestedSuperCallWithClassField_inCommaStatement()
func (c *RewriteClassMembersTest) TestNestedSuperCallWithClassField_inCommaStatement() {
	// We want to support this code pattern in case the compiler is asked to transpile classes that
	// were already minified by some other minifier, as that often leads to combining super()
	// calls with other code in the constructor.
	c.Test(
		`class A { constructor(...args) {} }
class B extends A {
  prop;
  constructor() {
    (super(1, 2, 3), bar());
  }
}
`,
		`class A { constructor(...args) {} }
class B extends A {
  constructor() {
    var JSCompiler_inline_result$jscomp$0 = super(1, 2, 3);
    this.prop = void 0;
    (JSCompiler_inline_result$jscomp$0,  bar());
  }
}
`)
}

// Java: RewriteClassMembersTest.testSuperCallInComplexExpression()
func (c *RewriteClassMembersTest) TestSuperCallInComplexExpression() {
	c.Test(
		`class A { constructor(...args) {} }
class B extends A {
  prop = 0;
  constructor(a) {
    foo(bar(), super(1, 2, 3), baz());
  }
}
`,
		`class A { constructor(...args) {} }
class B extends A {
  constructor(a) {
    var JSCompiler_temp_const$jscomp$1 = foo;
    var JSCompiler_temp_const$jscomp$0 = bar();
    var JSCompiler_inline_result$jscomp$2 = super(1, 2, 3);
    this.prop = 0;
    JSCompiler_temp_const$jscomp$1(
        JSCompiler_temp_const$jscomp$0, JSCompiler_inline_result$jscomp$2, baz());
  }
}
`)
}

// Java: RewriteClassMembersTest.testSuperCallInComplexExpression_cannotDecompose()
func (c *RewriteClassMembersTest) TestSuperCallInComplexExpression_cannotDecompose() {
	ex := AssertThrows[jutil.IRuntimeException](func() {
		c.Test(
			`class A { constructor(...args) {} }
class B extends A {
  prop = 0;
  constructor(a) {
    for (const x = foo(), y = super();;) {}
  }
}
`,
			"")
	})
	AssertThat(ex.GetMessage()).
		Contains(
			"Cannot decompose super() call in a class with class fields. Move super() call to the" +
				" root of the constructor.")
}

// Java: RewriteClassMembersTest.testConditionalSuperCalls()
func (c *RewriteClassMembersTest) TestConditionalSuperCalls() {
	// See: BUG-1 in go/tsjs-private-elements-bugs
	//
	// Compiler Rejection on Branching super() with Fields (Won't fix now):
	// Under ES2022, derived class constructors are allowed to call super() conditionally in
	// separate branches (e.g. if (cond) super(-1) else super(1)). Class field initializations are
	// defined by the specification to run immediately after super() returns and binds `this`.
	//
	// Current Compiler Behavior:
	// Es6NormalizeClasses assumes that constructors with fields have at most one super() call at
	// the root of the constructor body. When addTemporaryInsertionPoint() detects more than one
	// super() call, it throws an IllegalStateException:
	// "classes with public fields must have only one super() call at the constructor root".
	//
	// Reviewer Decision & Potential Future Fix:
	// This will not be fixed now because the compiler error message is clear and descriptive. We
	// will re-evaluate if necessary later if user requests arise. If revisited later, field and
	// private initializations can be extracted into an arrow function assigned to a local variable
	// and invoked immediately after each super() call.
	ex := AssertThrows[jutil.IRuntimeException](func() {
		c.Test(
			"class A { constructor(...args) {} }\nclass B extends A {\n  prop = 0;\n  #x = 1;\n  constructor(a) {\n    if (a < 0) {\n      // In ES2022, calling super() in branches is valid and binds `this`.\n      // Compiler rejects multiple super() calls in branches at compile time.\n      super(-1);\n    } else {\n      super(1);\n    }\n  }\n}\n",
			"")
	})
	AssertThat(ex.GetMessage()).
		Contains(
			"classes with public fields must have only one super() call at the constructor root")
}

// Java: RewriteClassMembersTest.testForwardPrivateFieldAccess_evaluatesToUndefined()
func (c *RewriteClassMembersTest) TestForwardPrivateFieldAccess_evaluatesToUndefined() {
	// TODO(b/236744850): Fix in a future CL by deferring property definition on the shadow instance
	// record until initialization and wrapping forward reads with $jscomp.checkPrivateGet.
	// See: BUG-2 in go/tsjs-private-elements-bugs
	//
	// Forward Private Field Access Evaluates to undefined instead of throwing TypeError:
	// ES2022 Specification:
	// Under ES2022 § 10.2.1.2 (InitializeInstanceElements), private field slots are installed
	// sequentially as each field definition executes. Reading a private field whose initializer
	// has not yet run must throw a TypeError at runtime because its private slot is uninitialized.
	//
	// Current Compiler Behavior:
	// Closure Compiler allocates a single shared record object (PRIVATE$1 = Object.create(null))
	// and registers the instance in the PrivateMap upon constructor entry:
	//   PRIVATE_MAP$0.set(this, PRIVATE$1);
	// When `#first = this.#later` executes, `#later` accesses the unassigned property on
	// PRIVATE$1:
	//   PRIVATE$1.first = PRIVATE_MAP$0.get(this).later;
	// Because in JavaScript reading an unassigned property on an object returns undefined,
	// this.#first silently evaluates to undefined instead of throwing a TypeError.
	c.Test(
		"class Foo {\n  // BUG: Evaluating `this.#later` here should throw a runtime TypeError per\n  // ES2022 § 10.2.1.2 because `#later` has not yet been initialized.\n  // In Closure Compiler, it silently evaluates to `undefined`.\n  #first = this.#later;\n  #later = 1;\n  getFirst() {\n    return this.#first;\n  }\n}\n",
		"const PRIVATE_MAP$0 = new $jscomp.PrivateMap();\nclass Foo {\n  constructor() {\n    const PRIVATE$1 = Object.create(null);\n    PRIVATE_MAP$0.set(this, PRIVATE$1);\n    // BUG: `later` property is unassigned on PRIVATE$1, so this evaluates to undefined\n    // instead of throwing a runtime TypeError.\n    PRIVATE$1.first = PRIVATE_MAP$0.get(this).later;\n    PRIVATE$1.later = 1;\n  }\n  getFirst() {\n    return PRIVATE_MAP$0.get(this).first;\n  }\n}\n")
}

// Java: RewriteClassMembersTest.testAbruptConstructorExecution_leaksBrandedInstance()
func (c *RewriteClassMembersTest) TestAbruptConstructorExecution_leaksBrandedInstance() {
	// TODO(b/236744850): Fix in a future CL by updating field brand checks (#prop in x) to check
	// property existence on the shadow record (key in privateMap.get(x)) using non-throwing lookup.
	// See: BUG-3 in go/tsjs-private-elements-bugs
	//
	// Abrupt Constructor Execution Leaks Fully-Branded Instance:
	// ES2022 Specification:
	// If an error is thrown during instance field initialization, constructor execution aborts.
	// Any subsequent private fields that were never reached must NOT be branded on the leaked
	// instance, so brand checks (`#later in leaked`) must return false.
	//
	// Current Compiler Behavior:
	// Closure Compiler brands the instance at the very beginning of constructor execution by
	// setting the instance in PRIVATE_MAP$0 before any field initializers run:
	//   PRIVATE_MAP$0.set(this, PRIVATE$1);
	// Even though `#later` throws an exception, `this` is already present in PRIVATE_MAP$0.
	// Brand check `#later in x` transpiles to `PRIVATE_MAP$0.has(x)`, which returns true for
	// the partially-constructed leaked instance.
	c.Test(
		"let leaked;\nclass Foo {\n  #first = (leaked = this, 1);\n  // Constructor aborts abruptly when evaluating this throw expression.\n  #later = (() => { throw new Error('fail'); })();\n  static isLaterBranded(x) {\n    // BUG: Under ES2022, `#later in leaked` should return false because construction\n    // aborted before #later was reached. In Closure Compiler, it returns true.\n    return #later in x;\n  }\n}\n\n// Code that triggers the bug at runtime:\ntry {\n  // 1. Attempting construction throws while initializing #later, leaking `this`.\n  new Foo();\n} catch (e) {\n  // Expected: constructor aborted due to error in #later initializer.\n}\n\n// 2. Evaluate brand check on the leaked partially-constructed instance:\n// - ES2022 Spec Requirement: Foo.isLaterBranded(leaked) === false (never reached #later).\n// - Current Compiler Behavior: Evaluates to `true` because PRIVATE_MAP$0.set(this, ...)\n//   ran at constructor entry before #later threw.\nconst leakedIsBranded = Foo.isLaterBranded(leaked);\n",
		`let leaked;
const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
class Foo {
  constructor() {
    const PRIVATE$1 = Object.create(null);
    // BUG: The instance is branded in PRIVATE_MAP$0 before field initializers run.
    // When #later throws below, leaked instance remains branded.
    PRIVATE_MAP$0.set(this, PRIVATE$1);
    PRIVATE$1.first = (leaked = this, 1);
    PRIVATE$1.later = (() => { throw new Error('fail'); })();
  }
  static isLaterBranded(x) {
    // BUG: Returns true on the leaked instance because PRIVATE_MAP$0.has(x) is true.
    return PRIVATE_MAP$0.has(x);
  }
}

try {
  new Foo();
} catch (e) {
}
const leakedIsBranded = Foo.isLaterBranded(leaked);
`)
}

// Java: RewriteClassMembersTest.testSelfIdentifierCollision_overwritesInstancePointer()
func (c *RewriteClassMembersTest) TestSelfIdentifierCollision_overwritesInstancePointer() {
	// TODO(b/236744850): Fix in a future CL by disambiguating private member property names on the
	// shadow record with a class-unique prefix (and dedicated prefix for $self).
	// See: BUG-4 in go/tsjs-private-elements-bugs
	//
	// Private Identifier #$self Name Collision with Compiler Internal Pointer:
	// ES2022 Specification:
	// `#$self` is a valid ECMAScript private identifier (PrivateIdentifier starts with # followed
	// by IdentifierName, where $ is a valid identifier character).
	//
	// Current Compiler Behavior:
	// In downlevel transpilation, private accessors/methods access the receiver instance via an
	// internal property named `$self` attached to the private record:
	//   PRIVATE$2.$self = this;
	// When a user declares a private field named `#$self`, the compiler strips the `#` and emits:
	//   PRIVATE$2.$self = 1;
	// which directly overwrites the internal instance back-pointer with `1`.
	// Subsequently, invoking the getter attempts to read `PRIVATE_MAP$0.get(this.$self).$self`,
	// which evaluates to `PRIVATE_MAP$0.get(1).$self` and throws a runtime TypeError:
	// "Cannot read properties of undefined".
	c.Test(
		`class Foo {
  #$self = 1;
  get #val() {
    return this.#$self;
  }
  getVal() {
    // BUG: Invoking getVal() throws a runtime TypeError ("Cannot read properties of undefined").
    // The user field #$self overwrites the compiler's internal $self receiver pointer.
    return this.#val;
  }
}
`,
		"const PRIVATE_MAP$0 = new $jscomp.PrivateMap();\nconst PRIVATE_PROTO$1 = Object.create(null, {\n  val: {\n    get: function() {\n      // BUG: this.$self was overwritten with 1 below; PRIVATE_MAP$0.get(1) returns undefined,\n      // and accessing .$self on undefined throws a runtime TypeError.\n      return PRIVATE_MAP$0.get(this.$self).$self;\n    }\n  }\n});\nclass Foo {\n  constructor() {\n    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);\n    PRIVATE$2.$self = this;\n    PRIVATE_MAP$0.set(this, PRIVATE$2);\n    // BUG: User field assignment `#$self = 1` overwrites `PRIVATE$2.$self = this` set above.\n    PRIVATE$2.$self = 1;\n  }\n  getVal() {\n    return PRIVATE_MAP$0.get(this).val;\n  }\n}\n")
}

// Java: RewriteClassMembersTest.testPrivateMethodAssignment_reportsCompileTimeError()
func (c *RewriteClassMembersTest) TestPrivateMethodAssignment_reportsCompileTimeError() {
	// See: BUG-5 in go/tsjs-private-elements-bugs
	//
	// Private Method Assignment Triggers Compile-Time Error (Won't fix now):
	// ES2022 Specification:
	// Under ES2022 § 13.15.2, assigning to a private method (e.g. `this.#connect = null`) is
	// syntactically valid code in ECMAScript grammar. At runtime, evaluating the assignment must
	// throw a TypeError: "Cannot assign to private method #connect".
	//
	// Current Compiler Behavior & Reviewer Decision:
	// Es6NormalizeClasses treats assignment to a private method as a compile-time build failure:
	//   compiler.report(JSError.make(getPropNode, ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT, ...));
	// This will not be fixed now because catching illegal assignments at compile time intentionally
	// shifts detectable bugs left to build time for Google codebases. We will re-evaluate if
	// strict runtime parity is required later.
	c.TestError(
		`class Service {
  #connect() {}
  tamper() {
    // BUG: Under ES2022 § 13.15.2, assigning to a private method is valid syntax that
    // should throw a TypeError at runtime when tamper() executes.
    // Instead, Closure Compiler raises a compile-time error:
    // JSC_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT.
    this.#connect = null;
  }
}
`,
		Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT)
}

// Java: RewriteClassMembersTest.testSetterOnlyPrivatePropertyRead_evaluatesToUndefined()
func (c *RewriteClassMembersTest) TestSetterOnlyPrivatePropertyRead_evaluatesToUndefined() {
	// TODO(b/236744850): Fix in a future CL by defining a throwing getter stub on the prototype
	// record (or reporting a compile-time diagnostic).
	// See: BUG-5 in go/tsjs-private-elements-bugs
	//
	// Reading Setter-Only Private Member Evaluates to undefined instead of TypeError:
	// ES2022 Specification:
	// Under ES2022 § 13.3.7.1, reading a private member that only defines a setter (and no getter)
	// must throw a runtime TypeError: "Cannot read private member #val without a getter".
	//
	// Current Compiler Behavior:
	// Private property reads rewrite to `PRIVATE_MAP$0.get(this).val`.
	// Because PRIVATE_PROTO$1 defines a setter on `val` but no getter:
	//   val: { set: function(v) {} }
	// reading `PRIVATE_MAP$0.get(this).val` returns undefined without throwing any error.
	c.Test(
		`class Foo {
  set #val(v) {}
  getVal() {
    // BUG: Under ES2022 § 13.3.7.1, reading a private property that has only a setter
    // must throw a runtime TypeError ("Cannot read private member without a getter").
    // Instead, reading this.#val evaluates to undefined.
    return this.#val;
  }
}
`,
		`const PRIVATE_MAP$0 = new $jscomp.PrivateMap();
const PRIVATE_PROTO$1 = Object.create(null, {
  val: {
    set: function(v) {}
  }
});
class Foo {
  constructor() {
    const PRIVATE$2 = Object.create(PRIVATE_PROTO$1);
    PRIVATE_MAP$0.set(this, PRIVATE$2);
  }
  getVal() {
    // BUG: PRIVATE_PROTO$1 has no getter, so reading .val evaluates to undefined
    // instead of throwing a runtime TypeError.
    return PRIVATE_MAP$0.get(this).val;
  }
}
`)
}

// Java: RewriteClassMembersTest.testBaseConstructorParamDefaultReadsField()
func (c *RewriteClassMembersTest) TestBaseConstructorParamDefaultReadsField() {
	// A base class initializes its fields before it binds the constructor's parameters, so
	// `y = this.a` reads 1. The parameter default moves into an assignment that runs after the
	// fields are initialized.
	c.Test(
		`class C {
  a = 1;
  constructor(x, y = this.a) {
    this.x = x;
    this.y = y;
  }
}
`,
		`class C {
  constructor(x, y) {
    this.a = 1;
    y = y === void 0 ? this.a : y;
    this.x = x;
    this.y = y;
  }
}
`)
	// `this` in a nested non-arrow function is not the instance being constructed.
	c.Test(
		`class D {
  a = 1;
  constructor(f = function() { return this; }) {
    this.f = f;
  }
}
`,
		`class D {
  constructor(f = function() { return this; }) {
    this.a = 1;
    this.f = f;
  }
}
`)
	// `super.hasOwnProperty('b')` passes `this` as the receiver, so the parameter default moves
	// after the fields too, and sees the field `b`.
	c.Test(
		`class E {
  b = 2;
  constructor(z = super.hasOwnProperty('b')) {
    this.z = z;
  }
}
`,
		`class E {
  constructor(z) {
    this.b = 2;
    z = z === void 0 ? super.hasOwnProperty('b') : z;
    this.z = z;
  }
}
`)
	// Destructuring parameters move into `let` declarations after the fields.
	c.Test(
		`class F {
  x = 1;
  constructor({a, b = this.x}) {
    this.a = a;
    this.b = b;
  }
}
`,
		`class F {
  constructor(TMP$0) {
    this.x = 1;
    let {a, b = this.x} = TMP$0;
    this.a = a;
    this.b = b;
  }
}
`)
	c.Test(
		`class G {
  x = {a: 1, b: 2};
  constructor({a, b} = this.x, ...[c = this.x]) {
    this.a = a;
    this.b = b;
    this.c = c;
  }
}
`,
		`class G {
  constructor(TMP$0, ...TMP$1) {
    this.x = {a: 1, b: 2};
    let {a, b} = TMP$0 === void 0 ? this.x : TMP$0;
    let [c = this.x] = TMP$1;
    this.a = a;
    this.b = b;
    this.c = c;
  }
}
`)
	// A local declaration in the constructor body that shadows a name referenced by a parameter
	// default is renamed by Normalize, so moving the default into the body still references the
	// outer variable.
	c.Test(
		`const x = 1;
class H {
  a = 2;
  constructor(y = x + this.a) {
    const x = 3;
    this.y = y + x;
  }
}
`,
		`const x = 1;
class H {
  constructor(y) {
    this.a = 2;
    y = y === void 0 ? x + this.a : y;
    const x$jscomp$1 = 3;
    this.y = y + x$jscomp$1;
  }
}
`)
}

// Java: RewriteClassMembersTest.testBaseClassParameterDefaultsOrder()
func (c *RewriteClassMembersTest) TestBaseClassParameterDefaultsOrder() {
	// b/572052294: A parameter default reads `this`, so all the defaults move into assignments
	// after the instance fields are initialized, as in a native base class constructor.
	c.Test(
		`class A {
  x = alert(1);
  constructor(a = alert(2), b = this.x) {
    alert(a, b);
  }
}
`,
		`class A {
  constructor(a, b) {
    this.x = alert(1);
    a = a === void 0 ? alert(2) : a;
    b = b === void 0 ? this.x : b;
    alert(a, b);
  }
}
`)
	// TODO(b/572052294): Defaults that don't read `this` run before the fields, so `alert(2)` runs
	// before `alert(1)`.
	c.Test(
		`class A {
  x = alert(1);
  constructor(a = alert(2)) {
    alert(a);
  }
}
`,
		`class A {
  constructor(a = alert(2)) {
    this.x = alert(1);
    alert(a);
  }
}
`)
}

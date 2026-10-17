// Port of com/google/javascript/jscomp/Es6RewriteClassExtendsExpressionsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/Es6RewriteClassExtendsExpressionsTest.java)

/*
 * Copyright 2018 The Closure Compiler Authors.
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

type Es6RewriteClassExtendsExpressionsTest struct {
	CompilerTestCaseBase
}

// Java: Es6RewriteClassExtendsExpressionsTest()
func NewEs6RewriteClassExtendsExpressionsTest() *Es6RewriteClassExtendsExpressionsTest {
	c := &Es6RewriteClassExtendsExpressionsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestEs6RewriteClassExtendsExpressionsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewEs6RewriteClassExtendsExpressionsTest() })
}

// Java: Es6RewriteClassExtendsExpressionsTest.getProcessor(Compiler)
func (c *Es6RewriteClassExtendsExpressionsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewEs6NormalizeClasses(compiler)
}

// Java: Es6RewriteClassExtendsExpressionsTest.setUp()
func (c *Es6RewriteClassExtendsExpressionsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.EnableTypeCheck()
	c.IgnoreWarnings(DiagnosticGroups_CHECK_TYPES)
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
	c.SetGenericNameReplacements(Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testBasic()
func (c *Es6RewriteClassExtendsExpressionsTest) TestBasic() {
	c.Test(
		"const foo = {'bar': Object}; /** @extends {Object} */ class Foo extends foo['bar'] {}",
		`const foo = {'bar': Object};
const CLASS_EXTENDS$0 = foo['bar'];
class Foo extends CLASS_EXTENDS$0 {}
`)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testName()
func (c *Es6RewriteClassExtendsExpressionsTest) TestName() {
	c.TestSame("class Foo extends Object {}")
}

// Java: Es6RewriteClassExtendsExpressionsTest.testGetProp()
func (c *Es6RewriteClassExtendsExpressionsTest) TestGetProp() {
	c.TestSame("const foo = { bar: Object}; class Foo extends foo.bar {}")
}

// Java: Es6RewriteClassExtendsExpressionsTest.testMixinFunction()
func (c *Es6RewriteClassExtendsExpressionsTest) TestMixinFunction() {
	c.Test(
		`/** @return {function(new:Object)} */
function mixObject(Superclass) {
  return class extends Superclass {
    bar() { return 'bar'; }
  };
}
class Foo {}
/** @extends {Object} */
class Bar extends mixObject(Foo) {}
`,
		`function mixObject(Superclass) {
  const CLASS_DECL$0 = class extends Superclass {
    bar() { return 'bar'; }
  };
  return CLASS_DECL$0;
}
class Foo {}
const CLASS_EXTENDS$1 = mixObject(Foo);
class Bar extends CLASS_EXTENDS$1 {}
`)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testClassExpressions()
func (c *Es6RewriteClassExtendsExpressionsTest) TestClassExpressions() {
	c.Test(
		`const foo = { bar: Object};
function baz(arg) {}
baz(class extends foo.bar {});
`,
		`const foo = {bar: Object};
function baz(arg) {}
const CLASS_DECL$0 = class extends foo.bar {};
baz(CLASS_DECL$0);
`)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testClassExtendsClassTranspilation()
func (c *Es6RewriteClassExtendsExpressionsTest) TestClassExtendsClassTranspilation() {
	code := `class __TEST_WebChannelConnection extends class __TEST_RestConnection {
  constructor(e) {
    this.databaseInfo = e, this.databaseId = e.databaseId;
  }
} {
  constructor(e) {
    super(e), this.forceLongPolling = e.forceLongPolling, this.autoDetectLongPolling = e.autoDetectLongPolling, this.useFetchStreams = e.useFetchStreams, this.longPollingOptions = e.longPollingOptions;
    console.log('test');
  }
}
`
	expectedCode := // The transpilation is somewhat unnecessarily complicated - we don't need to wrap the
		// extended class in an arrow function - but this is an artifact of the current
		// implementation.
		`const CLASS_EXTENDS$0 = (() => {
  const CLASS_DECL$1 = class {
    constructor(e) {
      this.databaseInfo = e, this.databaseId = e.databaseId;
    }
  };
  return CLASS_DECL$1;
})();
class __TEST_WebChannelConnection extends CLASS_EXTENDS$0 {
  constructor(e) {
    super(e), this.forceLongPolling = e.forceLongPolling, this.autoDetectLongPolling = e.autoDetectLongPolling, this.useFetchStreams = e.useFetchStreams, this.longPollingOptions = e.longPollingOptions;
    console.log('test');
  }
}
`
	c.Test(code, expectedCode)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testVarDeclaration()
func (c *Es6RewriteClassExtendsExpressionsTest) TestVarDeclaration() {
	c.Test(
		`const foo = {'bar': Object};
/** @extends {Object} */ var Foo = class extends foo['bar'] {};
`,
		`const foo = {'bar': Object};
const CLASS_EXTENDS$0 = foo['bar'];
var Foo = class extends CLASS_EXTENDS$0 {};
`)

	c.Test(
		`const foo = {'bar': Object};
/** @extends {Object} */ let Foo = class extends foo['bar'] {};
`,
		`const foo = {'bar': Object};
const CLASS_EXTENDS$0 = foo['bar'];
let Foo = class extends CLASS_EXTENDS$0 {};
`)

	c.Test(
		`const foo = {'bar': Object};
/** @extends {Object} */ const Foo = class extends foo['bar'] {};
`,
		`const foo = {'bar': Object};
const CLASS_EXTENDS$0 = foo['bar'];
const Foo = class extends CLASS_EXTENDS$0 {};
`)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testDeclarationInForLoop()
func (c *Es6RewriteClassExtendsExpressionsTest) TestDeclarationInForLoop() {
	c.Test(
		`const foo = {'bar': Object};
for (let Foo = /** @extends {Object} */ class extends foo['bar'] {}, i = 0;
    i < 1;
    i++)
  {}
`,
		`const foo = {'bar': Object};
// use an iife since parent of the let isn't a block where we can add statements
// TODO(bradfordcsmith): Would it be better (smaller output code) to just declare
//     the temporary variable in the block containing the for-loop?
for (let Foo = (() => {
    const CLASS_EXTENDS$0 = foo['bar'];
    const CLASS_DECL$1 = class extends CLASS_EXTENDS$0 {};
    return CLASS_DECL$1;
  })(), i = 0; i < 1; i++) {}
`)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testAssign()
func (c *Es6RewriteClassExtendsExpressionsTest) TestAssign() {
	c.Test(
		`const foo = {'bar': Object};
var Foo;
/** @extends {Object} */
Foo = class extends foo['bar'] {};
`,
		`const foo = {'bar': Object};
var Foo;
const CLASS_EXTENDS$0 = foo['bar'];
Foo = class extends CLASS_EXTENDS$0 {};
`)

	c.Test(
		`const foo = {'bar': Object};
/** @extends {Object} */
foo.baz = class extends foo['bar'] {};
`,
		`const foo = {'bar': Object};
const CLASS_EXTENDS$0 = foo['bar'];
foo.baz = class extends CLASS_EXTENDS$0 {};
`)

	c.Test(
		`const foo = {'bar': Object};
/** @extends {Object} */
foo.foo = foo.baz = class extends foo['bar'] {};
`,
		`const foo = {'bar': Object};
foo.foo = foo.baz = (() => {
  const CLASS_EXTENDS$0 = foo['bar'];
  const CLASS_DECL$1 = class extends CLASS_EXTENDS$0 {};
  return CLASS_DECL$1;
})();
`)

	c.Test(
		`const foo = {'bar': Object};
/** @extends {Object} */
foo['baz'] = class extends foo['bar'] {};
`,
		`const foo = {'bar': Object};
const CLASS_EXTENDS$0 = foo['bar'];
const CLASS_DECL$1 = class extends CLASS_EXTENDS$0 {};
foo['baz'] = CLASS_DECL$1;
`)

	c.Test(
		`const foo = {'bar': Object, baz: {}};
/** @extends {Object} */
foo.baz['foo'] = class extends foo['bar'] {};
`,
		`const foo = {'bar': Object, baz: {}};
const CLASS_EXTENDS$0 = foo['bar'];
const CLASS_DECL$1 = class extends CLASS_EXTENDS$0 {};
foo.baz['foo'] = CLASS_DECL$1;
`)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testAssignWithSideEffects()
func (c *Es6RewriteClassExtendsExpressionsTest) TestAssignWithSideEffects() {
	c.Test(
		`let baz = 1;
const foo = {'bar': Object};
/** @extends {Object} */
foo[baz++] = class extends foo['bar'] {};
`,
		`let baz = 1;
const foo = {'bar': Object};
// Use an IIFE to preserve execution order when expressions have side effects.
foo[baz++] = (() => {
  const CLASS_EXTENDS$0 = foo['bar'];
  const CLASS_DECL$1 = class extends CLASS_EXTENDS$0 {};
  return CLASS_DECL$1;
})();
`)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testMultipleVarDeclaration()
func (c *Es6RewriteClassExtendsExpressionsTest) TestMultipleVarDeclaration() {
	c.Test(
		`const foo = { 'bar': Object};
function mayHaveSideEffects() {}
var baz = mayHaveSideEffects(),
    Foo = /** @extends {Object} */ class extends foo['bar'] {};
`,
		`const foo = { 'bar': Object};
function mayHaveSideEffects() {}
var baz = mayHaveSideEffects(), Foo = (() => {
  const CLASS_EXTENDS$0 = foo['bar'];
  const CLASS_DECL$1 = class extends CLASS_EXTENDS$0 {};
  return CLASS_DECL$1;
})();
`)

	c.Test(
		`const foo = {'bar': Object};
const Foo = /** @extends {Object} */ class extends foo['bar'] {}, baz = false;
`,
		`const foo = {'bar': Object};
const CLASS_EXTENDS$0 = foo['bar'];
const Foo = class extends CLASS_EXTENDS$0 {}, baz = false;
`)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testDynamicExtends_logicalExpressionWithJSDoc()
func (c *Es6RewriteClassExtendsExpressionsTest) TestDynamicExtends_logicalExpressionWithJSDoc() {
	c.Test(
		`class A {}
class B {}
/** @extends {Object} */
class C extends (foo ? A : B) {}
`,
		`class A {}
class B {}
const CLASS_EXTENDS$0 = foo ? A : B;
class C extends CLASS_EXTENDS$0 {}
`)
}

// Java: Es6RewriteClassExtendsExpressionsTest.testNestedClassInExtends_withAnotherExtends()
func (c *Es6RewriteClassExtendsExpressionsTest) TestNestedClassInExtends_withAnotherExtends() {
	// Regression test for b/483475275.
	// Disable typechecking: this was seen in real-world code not written for the Closure Compiler,
	// and the typechecker will emit an error. That's fine, we just want to test transpilation.
	c.DisableTypeCheck()
	c.DisableTypeInfoValidation()

	c.Test(
		`new (class extends (class extends boo() {}) {});
`,
		`new ((() => {
  const CLASS_EXTENDS$0 = (() => {
    const CLASS_EXTENDS$2 = boo();
    const CLASS_DECL$3 = class extends CLASS_EXTENDS$2 {};
    return CLASS_DECL$3;
  })();
  const CLASS_DECL$1 = class extends CLASS_EXTENDS$0 {};
  return CLASS_DECL$1;
})())();
`)
}

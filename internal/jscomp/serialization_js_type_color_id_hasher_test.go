// Port of com/google/javascript/jscomp/serialization/JSTypeColorIdHasherTest.java @ closure-compiler v20261015

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
	"fmt"
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// The Java @Test methods are not named test*; their Go names get a "Test"
// prefix so that RunJUnit4 runs them. The class is @RunWith(Parameterized):
// every test method runs once per element of cases(), as subtests named
// "[<index>]" like JUnit's default parameter names.
//
// Java: public final class JSTypeColorIdHasherTest extends CompilerTestCase
type JSTypeColorIdHasherTest struct {
	CompilerTestCaseBase

	hasher           *JSTypeColorIdHasher
	labelToColorId   *jutil.LinkedHashMap[string, colors.ColorId]
	colorIdToJSTypes *jutil.LinkedHashMultimap[colors.ColorId, rhino.JSType] // Useful for debugging.

	// Java: @Parameterized.Parameter
	rewriteClosureModules bool
}

// Java: JSTypeColorIdHasherTest.CLOSURE_GLOBALS
const jsTypeColorIdHasherTest_CLOSURE_GLOBALS = `var goog = {};
goog.loadModule = function(def) {};
goog.module = function(name) {};
goog.provide = function(id) {};
`

// Java: JSTypeColorIdHasherTest.cases()
func JSTypeColorIdHasherTest_Cases() []bool {
	return []bool{true, false}
}

// Java: JSTypeColorIdHasherTest()
func NewJSTypeColorIdHasherTest() *JSTypeColorIdHasherTest {
	c := &JSTypeColorIdHasherTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, jsTypeColorIdHasherTest_CLOSURE_GLOBALS)
	return c
}

func TestJSTypeColorIdHasherTest(t *testing.T) {
	for i, rewriteClosureModules := range JSTypeColorIdHasherTest_Cases() {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			RunJUnit4(t, func() JUnitTestCase {
				c := NewJSTypeColorIdHasherTest()
				c.rewriteClosureModules = rewriteClosureModules
				return c
			})
		})
	}
}

// Java: JSTypeColorIdHasherTest.setUp()
func (c *JSTypeColorIdHasherTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableCreateModuleMap()
	c.EnableTypeCheck()

	if c.rewriteClosureModules {
		c.EnableRewriteClosureCode()
		c.DisableRewriteModulesAfterTypechecking()
	}
}

// Java: JSTypeColorIdHasherTest.LabeledTypeFinder
type jsTypeColorIdHasherTest_LabeledTypeFinder struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *JSTypeColorIdHasherTest
}

// Java: JSTypeColorIdHasherTest.LabeledTypeFinder.visit(NodeTraversal, Node, Node)
func (f *jsTypeColorIdHasherTest_LabeledTypeFinder) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !n.IsLabel() {
		return
	}

	label := n.GetFirstChild().GetString()
	type_ := n.GetSecondChild().GetFirstChild().GetJSType().ToMaybeObjectType()
	id := f.outer.hasher.HashObjectType(type_)

	AssertThat(f.outer.labelToColorId).DoesNotContainKey(label)
	f.outer.labelToColorId.Put(label, id)
	f.outer.colorIdToJSTypes.Put(id, type_)
}

// Java: JSTypeColorIdHasherTest.labelIdSet()
func (c *JSTypeColorIdHasherTest) labelIdSet() *jutil.LinkedHashSet[colors.ColorId] {
	return jutil.NewLinkedHashSet(slices.Collect(c.labelToColorId.Values())...)
}

// Java: JSTypeColorIdHasherTest.getProcessor(Compiler)
func (c *JSTypeColorIdHasherTest) GetProcessor(compiler *Compiler) CompilerPass {
	c.hasher = NewJSTypeColorIdHasher(compiler.GetTypeRegistry())
	c.labelToColorId = jutil.NewLinkedHashMap[string, colors.ColorId]()
	c.colorIdToJSTypes = jutil.NewLinkedHashMultimap[colors.ColorId, rhino.JSType]()

	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		NodeTraversal_TraverseRoots(compiler, &jsTypeColorIdHasherTest_LabeledTypeFinder{outer: c}, externs, root)
	})
}

// Java: JSTypeColorIdHasherTest.haveDifferentIds_ctorInstancePrototype()
func (c *JSTypeColorIdHasherTest) TestHaveDifferentIds_ctorInstancePrototype() {
	c.TestSame(
		`class Foo { }

CTOR: Foo;
INST: new Foo();
PROTO: Foo.prototype;
`)

	AssertThat(c.labelToColorId.Values()).ContainsNoDuplicates()
	AssertThat(c.labelToColorId.Keys()).ContainsExactly("CTOR", "INST", "PROTO")
}

// Java: JSTypeColorIdHasherTest.haveDifferentIds_qnamesWithSameEnding()
func (c *JSTypeColorIdHasherTest) TestHaveDifferentIds_qnamesWithSameEnding() {
	c.TestSame(
		`Foo = class { }
var a = {}; a.Foo = class { }
var b = {}; b.Foo = class { }

FOO: Foo;
A_FOO: a.Foo;
B_FOO: b.Foo;
`)

	AssertThat(c.labelToColorId.Values()).ContainsNoDuplicates()
	AssertThat(c.labelToColorId.Keys()).ContainsExactly("FOO", "A_FOO", "B_FOO")
}

// Java: JSTypeColorIdHasherTest.haveDifferentIds_classWithSameName_inDifferentGoogModules()
func (c *JSTypeColorIdHasherTest) TestHaveDifferentIds_classWithSameName_inDifferentGoogModules() {
	c.TestSame(
		c.Srcs(
			`goog.module('a');
class Foo { }
A_FOO: Foo;
`,
			`goog.module('b');
class Foo { }
B_FOO: Foo;
`))

	AssertThat(c.labelToColorId.Values()).ContainsNoDuplicates()
	AssertThat(c.labelToColorId.Keys()).ContainsExactly("A_FOO", "B_FOO")
}

// Java: JSTypeColorIdHasherTest.haveDifferentIds_classWithSameName_inDifferentGoogModules_inLoadModuleFile()
func (c *JSTypeColorIdHasherTest) TestHaveDifferentIds_classWithSameName_inDifferentGoogModules_inLoadModuleFile() {
	c.TestSame(
		`goog.loadModule(function(exports) {
  goog.module('a');

  class Foo { }
  A_FOO: Foo;

 return exports;
});

goog.loadModule(function(exports) {
  goog.module('b');

  class Foo { }
  B_FOO: Foo;

 return exports;
});
`)

	AssertThat(c.labelToColorId.Values()).ContainsNoDuplicates()
	AssertThat(c.labelToColorId.Keys()).ContainsExactly("A_FOO", "B_FOO")
}

// Java: JSTypeColorIdHasherTest.haveDifferentIds_classWithSameName_insideAndOutsideGoogModule()
func (c *JSTypeColorIdHasherTest) TestHaveDifferentIds_classWithSameName_insideAndOutsideGoogModule() {
	c.TestSame(
		c.Srcs(
			`class Foo { }
FOO: Foo;
`,
			`goog.module('b');
class Foo { }
B_FOO: Foo;
`))

	AssertThat(c.labelToColorId.Values()).ContainsNoDuplicates()
	AssertThat(c.labelToColorId.Keys()).ContainsExactly("FOO", "B_FOO")
}

// Java: JSTypeColorIdHasherTest.haveDifferentIds_anonymousTypes_withDifferentPropertyNames()
func (c *JSTypeColorIdHasherTest) TestHaveDifferentIds_anonymousTypes_withDifferentPropertyNames() {
	c.TestSame(
		`function withC() { };
/** @type {number} */
withC.c = 0;

WITH_A: ({a: 0});
WITH_B: ({b: 0});
WITH_C: withC;
`)

	AssertThat(c.labelToColorId.Values()).ContainsNoDuplicates()
	AssertThat(c.labelToColorId.Keys()).ContainsExactly("WITH_A", "WITH_B", "WITH_C")
}

// Java: JSTypeColorIdHasherTest.haveSameIds_anonymousTypes_withSamePropertyNames_andDiffererntPropertyTypes()
func (c *JSTypeColorIdHasherTest) TestHaveSameIds_anonymousTypes_withSamePropertyNames_andDiffererntPropertyTypes() {
	c.TestSame(
		`function withC() { };
/** @type {number} */
withC.c = 0;

WITH_A: ({a: null});
WITH_B: ({b: ''});
WITH_C: withC;
`)

	AssertThat(c.labelIdSet()).HasSize(3)
	AssertThat(c.labelToColorId.Keys()).ContainsExactly("WITH_A", "WITH_B", "WITH_C")
}

// TODO(b/185519307): We should probably differentiate the scoped types from the unscoped one, but
// it's hard to check the scoping after rewriting.
//
// Java: JSTypeColorIdHasherTest.haveSameIds_nominalTypes_fromDifferentScopes_fromSameGoogModule()
func (c *JSTypeColorIdHasherTest) TestHaveSameIds_nominalTypes_fromDifferentScopes_fromSameGoogModule() {
	c.TestSame(
		`goog.module('a');

class Foo { }
OUTER_FOO: Foo;

function scope0() {
  class Foo { }
  FOO_0: Foo;
}

function scope1() {
  class Foo { }
  FOO_1: Foo;
}

{
  class Foo { }
  FOO_BLOCK: Foo;
}
`)

	// TODO(b/185519307): With rewriting enabled, the reference name of OUTER_FOO is changed.
	expectedIdCount := 1
	if c.rewriteClosureModules {
		expectedIdCount = 2
	}

	AssertThat(c.labelIdSet()).HasSize(expectedIdCount)
	AssertThat(c.labelToColorId.Keys()).
		ContainsExactly("OUTER_FOO", "FOO_0", "FOO_1", "FOO_BLOCK")
}

// Java: JSTypeColorIdHasherTest.haveRecompileStableIds_nominalTypes()
func (c *JSTypeColorIdHasherTest) TestHaveRecompileStableIds_nominalTypes() {
	code :=
		`class Foo { }
CLASS: Foo;
INST: new Foo();
PROTO: Foo.prototype;

function namedFunction() { }
FUNC: namedFunction;

/** @enum */ const Enum = {};
ENUM: Enum;
`

	c.TestSame(code)
	idsFromFirstCompile := c.labelToColorId

	c.TestSame(code)
	idsFromSecondCompile := c.labelToColorId

	AssertThat(idsFromFirstCompile).IsNotEmpty()
	AssertThat(idsFromFirstCompile).IsNotSameInstanceAs(idsFromSecondCompile)
	AssertThat(idsFromFirstCompile).IsEqualTo(idsFromSecondCompile)
}

// Java: JSTypeColorIdHasherTest.haveRecompileStableIds_anonymousTypes()
func (c *JSTypeColorIdHasherTest) TestHaveRecompileStableIds_anonymousTypes() {
	code :=
		`LAMBDA: (function() { });
DOC_LAMBDA: /** @type {function()} */ (function() {});

RECORD: ({a: 0});
DOC_RECORD: /** @type {{a: (number|undefined)}} */ ({});
`

	c.TestSame(code)
	idsFromFirstCompile := c.labelToColorId

	c.TestSame(code)
	idsFromSecondCompile := c.labelToColorId

	AssertThat(idsFromFirstCompile).IsNotEmpty()
	AssertThat(idsFromFirstCompile).IsNotSameInstanceAs(idsFromSecondCompile)
	AssertThat(idsFromFirstCompile).IsEqualTo(idsFromSecondCompile)
}

// Java: JSTypeColorIdHasherTest.haveHeaderDerivableIds_nominalTypes()
func (c *JSTypeColorIdHasherTest) TestHaveHeaderDerivableIds_nominalTypes() {
	c.TestSame(
		`goog.module('a');

class Foo {
  /** @return {number} */
  someMethod() {
    return this.toString().length;
  }
}
CLASS: Foo;
INST: new Foo();
PROTO: Foo.prototype;

  /** @return {number} */
function namedFunction() {
  return this.toString().length;
}
FUNC: namedFunction;

/** @enum */ const Enum = {
  A: 0,
};
ENUM: Enum;
`)
	idsFromSource := c.labelToColorId

	c.TestSame(
		`goog.loadModule(function(exports) {
  goog.module('a');

  class Foo {
    /** @return {number} */
    someMethod() {
    }
  }
  CLASS: Foo;
  INST: new Foo();
  PROTO: Foo.prototype;

    /** @return {number} */
  function namedFunction() {
  }
  FUNC: namedFunction;

  /** @enum */ const Enum = {
    A: 0,
  };
  ENUM: Enum;

  return exports;
});
`)
	idsFromHeader := c.labelToColorId

	AssertThat(idsFromSource).IsNotEmpty()
	AssertThat(idsFromSource).IsNotSameInstanceAs(idsFromHeader)
	AssertThat(idsFromSource).IsEqualTo(idsFromHeader)
}

// Java: JSTypeColorIdHasherTest.haveHeaderDerivableIds_anonymousTypes()
func (c *JSTypeColorIdHasherTest) TestHaveHeaderDerivableIds_anonymousTypes() {
	c.TestSame(
		`goog.module('a');

LAMBDA: (function() { });
DOC_LAMBDA: /** @type {function()} */ (function() {});

RECORD: ({a: 0});
DOC_RECORD: /** @type {{a: (number|undefined)}} */ ({});
`)
	idsFromSource := c.labelToColorId

	c.TestSame(
		`goog.loadModule(function(exports) {
  goog.module('a');

  LAMBDA: (function() { });
  DOC_LAMBDA: /** @type {function()} */ (function() {});

  RECORD: ({a: 0});
  DOC_RECORD: /** @type {{a: (number|undefined)}} */ ({});

  return exports;
});
`)
	idsFromHeader := c.labelToColorId

	AssertThat(idsFromSource).IsNotEmpty()
	AssertThat(idsFromSource).IsNotSameInstanceAs(idsFromHeader)
	AssertThat(idsFromSource).IsEqualTo(idsFromHeader)
}

// Java: JSTypeColorIdHasherTest.returnsNativeId_forNativeTypes()
func (c *JSTypeColorIdHasherTest) TestReturnsNativeId_forNativeTypes() {
	registry := rhino.NewJSTypeRegistry(rhino.ErrorReporter_ALWAYS_THROWS_INSTANCE)
	c.hasher = NewJSTypeColorIdHasher(registry)

	for n, i := range JSTypeColorIdHasher_NATIVE_TYPE_TO_ID.All() {
		t := registry.GetNativeObjectType(n)
		AssertThat(c.hasher.HashObjectType(t)).IsEqualTo(i)
	}
}

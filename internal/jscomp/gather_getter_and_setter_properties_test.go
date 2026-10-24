// Port of com/google/javascript/jscomp/GatherGetterAndSetterPropertiesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/GatherGetterAndSetterPropertiesTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// GatherGetterAndSetterPropertiesTest holds tests for
// GatherGetterAndSetterProperties.
//
// Go-only: test-helper — Java @Test methods whose names do not start with
// "test" get a Test prefix, so that RunJUnit4 runs them.
type GatherGetterAndSetterPropertiesTest struct {
	CompilerTestCaseBase
}

// Java: GatherGetterAndSetterPropertiesTest()
func NewGatherGetterAndSetterPropertiesTest() *GatherGetterAndSetterPropertiesTest {
	c := &GatherGetterAndSetterPropertiesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestGatherGetterAndSetterPropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewGatherGetterAndSetterPropertiesTest() })
}

// Java: GatherGetterAndSetterPropertiesTest.getProcessor(Compiler)
func (c *GatherGetterAndSetterPropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewGatherGetterAndSetterProperties(compiler)
}

// Java: GatherGetterAndSetterPropertiesTest.defaultToNormal()
func (c *GatherGetterAndSetterPropertiesTest) TestDefaultToNormal() {
	c.TestSame("")
	AssertThat(c.GetLastAccessorSummary().GetKind("dne")).IsEqualTo(AccessorSummary_PropertyAccessKind_NORMAL)
}

// Java: GatherGetterAndSetterPropertiesTest.findGettersAndSettersInExterns()
func (c *GatherGetterAndSetterPropertiesTest) TestFindGettersAndSettersInExterns() {
	c.TestSame(c.Externs("({get x() {}, set y(v) {}})"), c.Srcs(""))

	AssertThat(c.GetLastCompiler().GetAccessorSummary().GetAccessors()).
		ContainsExactlyEntriesIn(
			gatherGetterAndSetterPropertiesTest_mapOf(
				"x", AccessorSummary_PropertyAccessKind_GETTER_ONLY, //
				"y", AccessorSummary_PropertyAccessKind_SETTER_ONLY))
}

// Java: GatherGetterAndSetterPropertiesTest.findGettersAndSettersInSources()
func (c *GatherGetterAndSetterPropertiesTest) TestFindGettersAndSettersInSources() {
	c.TestSame(c.Srcs("({get x() {}, set y(v) {}})"))

	AssertThat(c.GetLastCompiler().GetAccessorSummary().GetAccessors()).
		ContainsExactlyEntriesIn(
			gatherGetterAndSetterPropertiesTest_mapOf(
				"x", AccessorSummary_PropertyAccessKind_GETTER_ONLY, //
				"y", AccessorSummary_PropertyAccessKind_SETTER_ONLY))
}

// Java: GatherGetterAndSetterPropertiesTest.mergeInExterns()
func (c *GatherGetterAndSetterPropertiesTest) TestMergeInExterns() {
	c.TestSame(c.Externs("({get x() {}}); ({set x(v) {}});"), c.Srcs(""))
	AssertThat(c.GetLastAccessorSummary().GetKind("x")).
		IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
}

// Java: GatherGetterAndSetterPropertiesTest.mergeInSrcs()
func (c *GatherGetterAndSetterPropertiesTest) TestMergeInSrcs() {
	c.TestSame(c.Srcs("({get x() {}}); ({set x(v) {}});"))
	AssertThat(c.GetLastAccessorSummary().GetKind("x")).
		IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
}

// Java: GatherGetterAndSetterPropertiesTest.mergeExtersAndSrcs()
func (c *GatherGetterAndSetterPropertiesTest) TestMergeExtersAndSrcs() {
	c.TestSame(c.Externs("({get x() {}})"), c.Srcs("({set x(v) {}})"))
	AssertThat(c.GetLastAccessorSummary().GetKind("x")).
		IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
}

// Java: GatherGetterAndSetterPropertiesTest.detectObjectLiteral()
func (c *GatherGetterAndSetterPropertiesTest) TestDetectObjectLiteral() {
	c.TestSame(c.Srcs("var a = { set x(v) {} };"))
	AssertThat(c.GetLastAccessorSummary().GetKind("x")).IsEqualTo(AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	c.TestSame(c.Srcs("var a = { get x() {} };"))
	AssertThat(c.GetLastAccessorSummary().GetKind("x")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)
}

// Java: GatherGetterAndSetterPropertiesTest.detectClass()
func (c *GatherGetterAndSetterPropertiesTest) TestDetectClass() {
	c.TestSame(c.Srcs("class Class { set x(v) {} };"))
	AssertThat(c.GetLastAccessorSummary().GetKind("x")).IsEqualTo(AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	c.TestSame(c.Srcs("class Class { get x() {} };"))
	AssertThat(c.GetLastAccessorSummary().GetKind("x")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)
}

// Java: GatherGetterAndSetterPropertiesTest.detectClassStatic()
func (c *GatherGetterAndSetterPropertiesTest) TestDetectClassStatic() {
	c.TestSame(c.Srcs("class Class { static set x(v) {} };"))
	AssertThat(c.GetLastAccessorSummary().GetKind("x")).IsEqualTo(AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	c.TestSame(c.Srcs("class Class { static get x() {} };"))
	AssertThat(c.GetLastAccessorSummary().GetKind("x")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)
}

// Java: GatherGetterAndSetterPropertiesTest.detectDefinePropertyQuotedKeys()
func (c *GatherGetterAndSetterPropertiesTest) TestDetectDefinePropertyQuotedKeys() {
	c.TestSame(
		`Object.defineProperty(something, 'prop', {
    'get': function() {},
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)

	c.TestSame(
		`Object.defineProperty(something, 'prop', {
    'set': function(v) {}
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	c.TestSame(
		`Object.defineProperty(something, 'prop', {
    'get': function() {},
    'set': function(v) {}
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).
		IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
}

// Java: GatherGetterAndSetterPropertiesTest.detectDefinePropertyUnquotedKeys()
func (c *GatherGetterAndSetterPropertiesTest) TestDetectDefinePropertyUnquotedKeys() {
	c.TestSame(
		`Object.defineProperty(something, 'prop', {
    get: function() {},
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)

	c.TestSame(
		`Object.defineProperty(something, 'prop', {
    set: function(v) {}
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	c.TestSame(
		`Object.defineProperty(something, 'prop', {
    get: function() {},
    set: function(v) {}
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).
		IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
}

// Java: GatherGetterAndSetterPropertiesTest.detectDefinePropertyMemberFunctions()
func (c *GatherGetterAndSetterPropertiesTest) TestDetectDefinePropertyMemberFunctions() {
	c.TestSame(
		`Object.defineProperty(something, 'prop', {
    get() {},
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)

	c.TestSame(
		`Object.defineProperty(something, 'prop', {
    set(v) {}
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	c.TestSame(
		`Object.defineProperty(something, 'prop', {
    get() {},
    set(v) {}
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).
		IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
}

// Java: GatherGetterAndSetterPropertiesTest.detectDefinePropertiesQuotedKeys()
func (c *GatherGetterAndSetterPropertiesTest) TestDetectDefinePropertiesQuotedKeys() {
	c.TestSame(
		`Object.defineProperties(something, {
  'prop': {
    'get': function() {},
  }
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)

	c.TestSame(
		`Object.defineProperties(something, {
  'prop': {
    'set': function(v) {},
  }
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	c.TestSame(
		`Object.defineProperties(something, {
  'prop': {
    'get': function() {},
    'set': function(v) {},
  }
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).
		IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
}

// Java: GatherGetterAndSetterPropertiesTest.detectDefinePropertiesUnquotedKeys()
func (c *GatherGetterAndSetterPropertiesTest) TestDetectDefinePropertiesUnquotedKeys() {
	c.TestSame(
		`Object.defineProperties(something, {
  prop: {
    get: function() {},
  }
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)

	c.TestSame(
		`Object.defineProperties(something, {
  prop: {
    set: function(v) {},
  }
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	c.TestSame(
		`Object.defineProperties(something, {
  prop: {
    get: function() {},
    set: function(v) {},
  }
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).
		IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
}

// Java: GatherGetterAndSetterPropertiesTest.detectDefinePropertiesMemberFunctions()
func (c *GatherGetterAndSetterPropertiesTest) TestDetectDefinePropertiesMemberFunctions() {
	c.TestSame(
		`Object.defineProperties(something, {
  prop: {
    get() {},
  }
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)

	c.TestSame(
		`Object.defineProperties(something, {
  'prop': {
    set(v) {},
  }
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_SETTER_ONLY)

	c.TestSame(
		`Object.defineProperties(something, {
  prop: {
    get() {},
    set(v) {},
  }
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).
		IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
}

// Java: GatherGetterAndSetterPropertiesTest.detectJscompGlobalObject()
func (c *GatherGetterAndSetterPropertiesTest) TestDetectJscompGlobalObject() {
	c.TestSame(
		`$jscomp.global.Object.defineProperty(something, 'prop', {
    get: function() {},
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)

	c.TestSame(
		`$jscomp$global.Object.defineProperty(something, 'prop', {
    get: function() {}
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_ONLY)

	c.TestSame(
		`other.Object.defineProperty(something, 'prop', {
    get: function() {},
});
`)

	AssertThat(c.GetLastAccessorSummary().GetKind("prop")).IsEqualTo(AccessorSummary_PropertyAccessKind_NORMAL)
}

// Java: GatherGetterAndSetterPropertiesTest.noAssumePropertiesAreStaticallyAnalyzable()
func (c *GatherGetterAndSetterPropertiesTest) TestNoAssumePropertiesAreStaticallyAnalyzable() {
	options := NewCompilerOptions()
	options.SetAssumePropertiesAreStaticallyAnalyzable(false)
	compiler := NewCompiler()
	compiler.Init([]*SourceFile{}, []*SourceFile{}, options)

	NewGatherGetterAndSetterProperties(compiler).
		Process(compiler.GetExternsRoot(), compiler.GetJsRoot())

	AssertThat(compiler.GetAccessorSummary().GetKind("dne")).
		IsEqualTo(AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
}

// Java: GatherGetterAndSetterPropertiesTest.getLastAccessorSummary()
func (c *GatherGetterAndSetterPropertiesTest) GetLastAccessorSummary() *AccessorSummary {
	return c.GetLastCompiler().GetAccessorSummary()
}

// gatherGetterAndSetterPropertiesTest_mapOf is ImmutableMap.of(k1, v1, k2, v2).
//
// Go-only: test-helper.
func gatherGetterAndSetterPropertiesTest_mapOf(
	k1 string, v1 AccessorSummary_PropertyAccessKind,
	k2 string, v2 AccessorSummary_PropertyAccessKind) *jutil.LinkedHashMap[string, AccessorSummary_PropertyAccessKind] {
	return jutil.NewImmutableMapBuilder[string, AccessorSummary_PropertyAccessKind]().Put(k1, v1).Put(k2, v2).BuildOrThrow()
}

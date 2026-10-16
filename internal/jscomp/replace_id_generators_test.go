// Port of com/google/javascript/jscomp/ReplaceIdGeneratorsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/ReplaceIdGeneratorsTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/parser"
)

// ReplaceIdGeneratorsTest holds tests for ReplaceIdGenerators.
type ReplaceIdGeneratorsTest struct {
	CompilerTestCaseBase

	generatePseudoNames bool
	lastPass            *ReplaceIdGenerators
	// nil is Java null.
	previousMappings *string
}

// Java: ReplaceIdGeneratorsTest()
func NewReplaceIdGeneratorsTest() *ReplaceIdGeneratorsTest {
	c := &ReplaceIdGeneratorsTest{generatePseudoNames: false, lastPass: nil, previousMappings: nil}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestReplaceIdGeneratorsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewReplaceIdGeneratorsTest() })
}

// replaceIdGeneratorsTest_ptr boxes a non-null Java String.
//
// Go-only: test-helper.
func replaceIdGeneratorsTest_ptr(s string) *string { return &s }

// replaceIdGeneratorsTest_idTestMap is the anonymous RenamingMap of
// getProcessor.
type replaceIdGeneratorsTest_idTestMap struct {
	map_ *jutil.LinkedHashMap[string, string]
}

// Java: ReplaceIdGeneratorsTest.getProcessor(Compiler).new RenamingMap() {...}.get(String)
func (m *replaceIdGeneratorsTest_idTestMap) Get(value string) *string {
	replacement, ok := m.map_.Lookup(value)
	if ok {
		return &replacement
	}
	r := "unknown:" + value
	return &r
}

// Java: ReplaceIdGeneratorsTest.getProcessor(Compiler)
func (c *ReplaceIdGeneratorsTest) GetProcessor(compiler *Compiler) CompilerPass {
	idTestMap := &replaceIdGeneratorsTest_idTestMap{
		map_: jutil.NewImmutableMapBuilder[string, string]().
			Put("foo", ":foo:").
			Put("bar", ":bar:").
			BuildOrThrow(),
	}
	c.lastPass = NewReplaceIdGenerators(
		compiler,
		/* templateLiteralsAreTranspiled= */ compiler.
			GetOptions().
			NeedsTranspilationOf(parser.FeatureSet_Feature_TEMPLATE_LITERALS),
		jutil.NewImmutableMapBuilder[string, RenamingMap]().
			Put("goog.events.getUniqueId", RenamingToken_INCONSISTENT).
			Put("goog.place.getUniqueId", RenamingToken_INCONSISTENT).
			Put("id", idTestMap).
			Put("get.id", idTestMap).
			BuildOrThrow(),
		c.generatePseudoNames,
		c.previousMappings,
		/* xidHashFunction= */ nil)
	return c.lastPass
}

// Java: ReplaceIdGeneratorsTest.setUp()
func (c *ReplaceIdGeneratorsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.generatePseudoNames = false
	c.previousMappings = nil
}

// Java: ReplaceIdGeneratorsTest.getOptions()
func (c *ReplaceIdGeneratorsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_INJECT)
	return options
}

// Java: ReplaceIdGeneratorsTest.testBackwardCompat()
func (c *ReplaceIdGeneratorsTest) TestBackwardCompat() {
	c.TestWithPseudo(
		"foo.bar = goog.events.getUniqueId('foo_bar')", "foo.bar = 'a'", "foo.bar = 'foo_bar$0'")
}

// Java: ReplaceIdGeneratorsTest.testSerialization1()
func (c *ReplaceIdGeneratorsTest) TestSerialization1() {
	c.TestMap(
		`var x = goog.events.getUniqueId('xxx');
var y = goog.events.getUniqueId('yyy');
`,
		`var x = 'a';
var y = 'b';
`,
		`[goog.events.getUniqueId]

a:testcode:1:32
b:testcode:2:32
`)
}

// Java: ReplaceIdGeneratorsTest.testSerialization2()
func (c *ReplaceIdGeneratorsTest) TestSerialization2() {
	c.TestMap(
		`/** @idGenerator {consistent} */
xid = function() {};
f1 = xid('f1');
f1 = xid('f1')
`,
		`/** @idGenerator {consistent} */
xid = function() {};
f1 = 'a';
f1 = 'a'
`,
		`[xid]

a:f1
`)
}

// Java: ReplaceIdGeneratorsTest.testReusePreviousSerialization1()
func (c *ReplaceIdGeneratorsTest) TestReusePreviousSerialization1() {
	c.previousMappings = replaceIdGeneratorsTest_ptr(`[goog.events.getUniqueId]

previous1:testcode:1:32
previous2:testcode:2:32

[goog.place.getUniqueId]
`)
	c.TestMap(
		`var x = goog.events.getUniqueId('xxx');
var y = goog.events.getUniqueId('yyy');
`,
		`var x = 'previous1';
var y = 'previous2';
`,
		`[goog.events.getUniqueId]

previous1:testcode:1:32
previous2:testcode:2:32
`)
}

// Java: ReplaceIdGeneratorsTest.testReusePreviousSerialization2()
func (c *ReplaceIdGeneratorsTest) TestReusePreviousSerialization2() {
	c.previousMappings = replaceIdGeneratorsTest_ptr(`[goog.events.getUniqueId]

a:testcode:1:32
b:testcode:2:32

[goog.place.getUniqueId]
`)
	c.TestMap(
		`var x = goog.events.getUniqueId('xxx');
// new line to change location
var y = goog.events.getUniqueId('yyy');
`,
		`var x = 'a';
var y = 'c';
`,
		`[goog.events.getUniqueId]

a:testcode:1:32
c:testcode:3:32
`)
}

// Java: ReplaceIdGeneratorsTest.testReusePreviousSerializationConsistent1()
func (c *ReplaceIdGeneratorsTest) TestReusePreviousSerializationConsistent1() {
	c.previousMappings = replaceIdGeneratorsTest_ptr(`[cid]

a:f1
`)
	c.TestMap(
		`/** @idGenerator {consistent} */ cid = function() {};
f1 = cid('f1');
f1 = cid('f1')
`,
		`/** @idGenerator {consistent} */ cid = function() {};
f1 = 'a';
f1 = 'a'
`,
		`[cid]

a:f1
`)
}

// Java: ReplaceIdGeneratorsTest.testNullishCoalesce()
func (c *ReplaceIdGeneratorsTest) TestNullishCoalesce() {
	c.TestWithPseudo(
		`/** @idGenerator */ foo.getUniqueId = function() {x ?? y;};
foo.bar = foo.getUniqueId('foo_bar')
`,
		`/** @idGenerator */ foo.getUniqueId = function() {x ?? y;};
foo.bar = 'a'
`,
		`/** @idGenerator */ foo.getUniqueId = function() {x ?? y;};
foo.bar = 'foo_bar$0'
`)
}

// Java: ReplaceIdGeneratorsTest.testSimple()
func (c *ReplaceIdGeneratorsTest) TestSimple() {
	c.TestWithPseudo(
		`/** @idGenerator */ foo.getUniqueId = function() {};
foo.bar = foo.getUniqueId('foo_bar')
`,
		`/** @idGenerator */ foo.getUniqueId = function() {};
foo.bar = 'a'
`,
		`/** @idGenerator */ foo.getUniqueId = function() {};
foo.bar = 'foo_bar$0'
`)

	c.TestWithPseudo(
		`/** @idGenerator */ goog.events.getUniqueId = function() {};
foo1 = goog.events.getUniqueId('foo1');
foo1 = goog.events.getUniqueId('foo1');
`,
		`/** @idGenerator */ goog.events.getUniqueId = function() {};
foo1 = 'a';
foo1 = 'b';
`,
		`/** @idGenerator */ goog.events.getUniqueId = function() {};
foo1 = 'foo1$0';
foo1 = 'foo1$1';
`)
}

// Java: ReplaceIdGeneratorsTest.testIndirectCall()
func (c *ReplaceIdGeneratorsTest) TestIndirectCall() {
	c.TestWithPseudo(
		`/** @idGenerator */ foo.getUniqueId = function() {};
foo.bar = (0, foo.getUniqueId)('foo_bar')
`,
		`/** @idGenerator */ foo.getUniqueId = function() {};
foo.bar = 'a'
`,
		`/** @idGenerator */ foo.getUniqueId = function() {};
foo.bar = 'foo_bar$0'
`)
	// JSCompiler inserts JSCOMPILER_PRESERVE(...) in the "UselessCode" analysis.
	c.TestWithPseudo(
		`/** @idGenerator */ foo.getUniqueId = function() {};
foo.bar = (JSCOMPILER_PRESERVE(0), foo.getUniqueId)('foo_bar')
`,
		`/** @idGenerator */ foo.getUniqueId = function() {};
foo.bar = 'a'
`,
		`/** @idGenerator */ foo.getUniqueId = function() {};
foo.bar = 'foo_bar$0'
`)
}

// Java: ReplaceIdGeneratorsTest.testObjectLit()
func (c *ReplaceIdGeneratorsTest) TestObjectLit() {
	c.TestWithPseudo(
		`/** @idGenerator */ goog.id = function() {};
things = goog.id({foo1: 'test', 'foo bar': 'test'})
`,
		`/** @idGenerator */ goog.id = function() {};
things = {'a': 'test', 'b': 'test'}
`,
		`/** @idGenerator */ goog.id = function() {};
things = {'foo1$0': 'test', 'foo bar$1': 'test'}
`)
}

// Java: ReplaceIdGeneratorsTest.testObjectLit_mapped()
func (c *ReplaceIdGeneratorsTest) TestObjectLit_mapped() {
	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {mapped} */ id = function() {};
things = id({foo: 'test', 'bar': 'test'})
`,
		`/** @idGenerator {mapped} */ id = function() {};
things = {':foo:': 'test', ':bar:': 'test'}
`)
}

// Java: ReplaceIdGeneratorsTest.testObjectLit_xid()
func (c *ReplaceIdGeneratorsTest) TestObjectLit_xid() {
	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {xid} */ xid.object = function() {};
things = xid.object({foo: 'test', 'value': 'test'})
`,
		`/** @idGenerator {xid} */ xid.object = function() {};
things = {'QB6rXc': 'test', 'b6Lt6c': 'test'}
`)
}

// Java: ReplaceIdGeneratorsTest.testObjectLit_empty()
func (c *ReplaceIdGeneratorsTest) TestObjectLit_empty() {
	c.TestWithPseudo(
		"/** @idGenerator */ goog.id = function() {}; things = goog.id({})",
		"/** @idGenerator */ goog.id = function() {}; things = {}",
		"/** @idGenerator */ goog.id = function() {}; things = {}")
}

// Java: ReplaceIdGeneratorsTest.testObjectLit_function()
func (c *ReplaceIdGeneratorsTest) TestObjectLit_function() {
	c.TestWithPseudo(
		`/** @idGenerator */ goog.id = function() {};
things = goog.id({foo: function() {}})
`,
		`/** @idGenerator */ goog.id = function() {};
things = {'a': function() {}}
`,
		`/** @idGenerator */ goog.id = function() {};
things = {'foo$0': function() {}}
`)

	c.TestWithPseudo(
		`/** @idGenerator */ goog.id = function() {};
things = goog.id({foo: function*() {}})
`,
		`/** @idGenerator */ goog.id = function() {};
things = {'a': function*() {}}
`,
		`/** @idGenerator */ goog.id = function() {};
things = {'foo$0': function*() {}}
`)
}

// Java: ReplaceIdGeneratorsTest.testObjectLit_ES6()
func (c *ReplaceIdGeneratorsTest) TestObjectLit_ES6() {
	c.TestError(
		`/** @idGenerator */
goog.id = function() {};
things = goog.id({fooX() {}})
`,
		ReplaceIdGenerators_SHORTHAND_FUNCTION_NOT_SUPPORTED_IN_ID_GEN)

	c.TestError(
		`/** @idGenerator */
goog.id = function() {};
things = goog.id({['fooX']: 'test'})
`,
		ReplaceIdGenerators_COMPUTED_PROP_NOT_SUPPORTED_IN_ID_GEN)
}

// Java: ReplaceIdGeneratorsTest.testClass()
func (c *ReplaceIdGeneratorsTest) TestClass() {
	c.TestSame(
		c.Externs(""),
		c.Srcs(
			`/** @idGenerator */
goog.id = function() {};
things = goog.id(class fooBar{})
`),
		c.Warning(ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER))
}

// Java: ReplaceIdGeneratorsTest.testSimpleConsistent()
func (c *ReplaceIdGeneratorsTest) TestSimpleConsistent() {
	c.TestWithPseudo(
		"/** @idGenerator {consistent} */ cid = function() {}; foo.bar = cid('foo_bar')",
		"/** @idGenerator {consistent} */ cid = function() {}; foo.bar = 'a'",
		"/** @idGenerator {consistent} */ cid = function() {}; foo.bar = 'foo_bar$0'")

	c.TestWithPseudo(
		"/** @idGenerator {consistent} */ cid = function() {}; f1 = cid('f1'); f1 = cid('f1')",
		"/** @idGenerator {consistent} */ cid = function() {}; f1 = 'a'; f1 = 'a'",
		"/** @idGenerator {consistent} */ cid = function() {}; f1 = 'f1$0'; f1 = 'f1$0'")

	c.TestWithPseudo(
		`/** @idGenerator {consistent} */ cid = function() {};
f1 = cid('f1');
f1 = cid('f1');
f1 = cid('f1')
`,
		`/** @idGenerator {consistent} */ cid = function() {};
f1 = 'a';
f1 = 'a';
f1 = 'a'
`,
		`/** @idGenerator {consistent} */ cid = function() {};
f1 = 'f1$0';
f1 = 'f1$0';
f1 = 'f1$0'
`)
}

// Java: ReplaceIdGeneratorsTest.testSimpleStable()
func (c *ReplaceIdGeneratorsTest) TestSimpleStable() {
	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {stable} */ sid = function() {};
foo.bar = sid('foo_bar')
`,
		`/** @idGenerator {stable} */ sid = function() {};
foo.bar = '125lGg'
`)

	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {stable} */ sid = function() {};
f1 = sid('f1');
f1 = sid('f1')
`,
		`/** @idGenerator {stable} */ sid = function() {};
f1 = 'AAAMiw';
f1 = 'AAAMiw'
`)
}

// Java: ReplaceIdGeneratorsTest.testSimpleXid()
func (c *ReplaceIdGeneratorsTest) TestSimpleXid() {
	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {xid} */ xid = function() {};
foo.bar = xid('foo')
`,
		`/** @idGenerator {xid} */ xid = function() {};
foo.bar = 'QB6rXc'
`)

	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {xid} */ xid = function() {};
f1 = xid('foo');
f1 = xid('foo')
`,
		`/** @idGenerator {xid} */ xid = function() {};
f1 = 'QB6rXc';
f1 = 'QB6rXc'
`)

	c.TestNonPseudoSupportingGenerator(
		"/** @idGenerator {xid} */ xid = function() {};\nf1 = xid(`foo`);\nf1 = xid(`foo`)\n",
		`/** @idGenerator {xid} */ xid = function() {};
f1 = 'QB6rXc';
f1 = 'QB6rXc'
`)
}

// Java: ReplaceIdGeneratorsTest.testVar()
func (c *ReplaceIdGeneratorsTest) TestVar() {
	c.TestWithPseudo(
		`/** @idGenerator {consistent} */ var cid = function() {};
foo.bar = cid('foo_bar')
`,
		`/** @idGenerator {consistent} */ var cid = function() {};
foo.bar = 'a'
`,
		`/** @idGenerator {consistent} */ var cid = function() {};
foo.bar = 'foo_bar$0'
`)

	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {stable} */ var sid = function() {};
foo.bar = sid('foo_bar')
`,
		`/** @idGenerator {stable} */ var sid = function() {};
foo.bar = '125lGg'
`)
}

// Java: ReplaceIdGeneratorsTest.testLet()
func (c *ReplaceIdGeneratorsTest) TestLet() {
	c.TestWithPseudo(
		`/** @idGenerator {consistent} */ let cid = function() {};
foo.bar = cid('foo_bar')
`,
		`/** @idGenerator {consistent} */ let cid = function() {};
foo.bar = 'a'
`,
		`/** @idGenerator {consistent} */ let cid = function() {};
foo.bar = 'foo_bar$0'
`)

	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {stable} */ let sid = function() {};
foo.bar = sid('foo_bar')
`,
		`/** @idGenerator {stable} */ let sid = function() {};
foo.bar = '125lGg'
`)
}

// Java: ReplaceIdGeneratorsTest.testConst()
func (c *ReplaceIdGeneratorsTest) TestConst() {
	c.TestWithPseudo(
		`/** @idGenerator {consistent} */ const cid = function() {};
foo.bar = cid('foo_bar')
`,
		`/** @idGenerator {consistent} */ const cid = function() {};
foo.bar = 'a'
`,
		`/** @idGenerator {consistent} */ const cid = function() {};
foo.bar = 'foo_bar$0'
`)

	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {stable} */ const cid = function() {};
foo.bar = cid('foo_bar')
`,
		`/** @idGenerator {stable} */ const cid = function() {};
foo.bar = '125lGg'
`)
}

// Java: ReplaceIdGeneratorsTest.testInObjLit()
func (c *ReplaceIdGeneratorsTest) TestInObjLit() {
	c.TestWithPseudo(
		`/** @idGenerator {consistent} */ cid = function() {};
foo.bar = {a: cid('foo_bar')}
`,
		`/** @idGenerator {consistent} */ cid = function() {};
foo.bar = {a: 'a'}
`,
		`/** @idGenerator {consistent} */ cid = function() {};
foo.bar = {a: 'foo_bar$0'}
`)

	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {stable} */ sid = function() {};
foo.bar = {a: sid('foo_bar')}
`,
		`/** @idGenerator {stable} */ sid = function() {};
foo.bar = {a: '125lGg'}
`)

	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {xid} */ get.xid = function() {};
foo.bar = {a: get.xid('foo')}
`,
		`/** @idGenerator {xid} */ get.xid = function() {};
foo.bar = {a: 'QB6rXc'}
`)
}

// Java: ReplaceIdGeneratorsTest.testInObjLit_mapped()
func (c *ReplaceIdGeneratorsTest) TestInObjLit_mapped() {
	c.TestWithPseudo("foo.bar = {a: id('foo')}", "foo.bar = {a: ':foo:'}", "foo.bar = {a: ':foo:'}")
}

// Java: ReplaceIdGeneratorsTest.testMapped()
func (c *ReplaceIdGeneratorsTest) TestMapped() {
	c.TestWithPseudo("foo.bar = id('foo');", "foo.bar = ':foo:';", "foo.bar = ':foo:';")
}

// Java: ReplaceIdGeneratorsTest.testMappedMap()
func (c *ReplaceIdGeneratorsTest) TestMappedMap() {
	c.TestMap(
		`foo.bar = id('foo');
foo.bar = id('foo');
`,
		`foo.bar = ':foo:';
foo.bar = ':foo:';
`,
		`[id]

:foo::foo
`)
}

// Java: ReplaceIdGeneratorsTest.testMapped2()
func (c *ReplaceIdGeneratorsTest) TestMapped2() {
	c.TestWithPseudo(
		`/** @idGenerator {mapped}*/ id = function() {};
foo.bar = function() { return id('foo'); };
`,
		`/** @idGenerator {mapped}*/ id = function() {};
foo.bar = function() { return ':foo:'; };
`,
		`/** @idGenerator {mapped}*/ id = function() {};
foo.bar = function() { return ':foo:'; };
`)
}

// Java: ReplaceIdGeneratorsTest.testTwoGenerators()
func (c *ReplaceIdGeneratorsTest) TestTwoGenerators() {
	c.TestWithPseudo(
		`/** @idGenerator */ var id1 = function() {};
/** @idGenerator */ var id2 = function() {};
f1 = id1('1');
f2 = id1('1');
f3 = id2('1');
f4 = id2('1');
`,
		`/** @idGenerator */ var id1 = function() {};
/** @idGenerator */ var id2 = function() {};
f1 = 'a';
f2 = 'b';
f3 = 'a';
f4 = 'b';
`,
		`/** @idGenerator */ var id1 = function() {};
/** @idGenerator */ var id2 = function() {};
f1 = '1$0';
f2 = '1$1';
f3 = '1$0';
f4 = '1$1';
`)
}

// Java: ReplaceIdGeneratorsTest.testMixedGenerators()
func (c *ReplaceIdGeneratorsTest) TestMixedGenerators() {
	c.TestWithPseudo(
		`/** @idGenerator */ var id1 = function() {};
/** @idGenerator {consistent} */ var id2 = function() {};
/** @idGenerator {stable} */ var id3 = function() {};
f1 = id1('1');
f2 = id1('1');
f3 = id2('1');
f4 = id2('1');
f5 = id3('1');
f6 = id3('1');
`,
		`/** @idGenerator */ var id1 = function() {};
/** @idGenerator {consistent} */ var id2 = function() {};
/** @idGenerator {stable} */ var id3 = function() {};
f1 = 'a';
f2 = 'b';
f3 = 'a';
f4 = 'a';
f5 = 'AAAAMQ';
f6 = 'AAAAMQ';
`,
		`/** @idGenerator */ var id1 = function() {};
/** @idGenerator {consistent} */ var id2 = function() {};
/** @idGenerator {stable} */ var id3 = function() {};
f1 = '1$0';
f2 = '1$1';
f3 = '1$0';
f4 = '1$0';
f5 = 'AAAAMQ';
f6 = 'AAAAMQ';
`)
}

// Java: ReplaceIdGeneratorsTest.testNonLiteralParam1()
func (c *ReplaceIdGeneratorsTest) TestNonLiteralParam1() {
	c.TestSame(
		`/** @idGenerator */ var id = function() {};
var x = 'foo';
id(x);
`,
		ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER)
}

// Java: ReplaceIdGeneratorsTest.testNonLiteralParam2()
func (c *ReplaceIdGeneratorsTest) TestNonLiteralParam2() {
	c.TestSame(
		`/** @idGenerator */ var id = function() {};
id('foo' + 'bar');
`,
		ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER)
}

// Java: ReplaceIdGeneratorsTest.testLocalCall()
func (c *ReplaceIdGeneratorsTest) TestLocalCall() {
	c.TestError(
		`/** @idGenerator */ var iid = function() {};
function Foo() { iid('foo'); }
`,
		ReplaceIdGenerators_NON_GLOBAL_ID_GENERATOR_CALL)
}

// Java: ReplaceIdGeneratorsTest.testConditionalCall()
func (c *ReplaceIdGeneratorsTest) TestConditionalCall() {
	c.TestError(
		`/** @idGenerator */
var xid = function() {};
while(0){ xid('foo');}
`,
		ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL)

	c.TestError(
		`/** @idGenerator */
var xid = function() {};
for(;;){ xid('foo');}
`,
		ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL)

	c.TestError(
		`/** @idGenerator */ var xid = function() {};
if(x) xid('foo');
`,
		ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL)

	c.TestWithPseudo(
		`/** @idGenerator {consistent} */ var xid = function() {};
function fb() {foo.bar = xid('foo_bar')}
`,
		`/** @idGenerator {consistent} */ var xid = function() {};
function fb() {foo.bar = 'a'}
`,
		`/** @idGenerator {consistent} */ var xid = function() {};
function fb() {foo.bar = 'foo_bar$0'}
`)

	c.TestNonPseudoSupportingGenerator(
		`/** @idGenerator {stable} */ var xid = function() {};
function fb() {foo.bar = xid('foo_bar')}
`,
		`/** @idGenerator {stable} */ var xid = function() {};
function fb() {foo.bar = '125lGg'}
`)

	c.TestError(
		`/** @idGenerator */
var xid = function() {};
for(x of [1, 2, 3]){ xid('foo');}
`,
		ReplaceIdGenerators_CONDITIONAL_ID_GENERATOR_CALL)
}

// Java: ReplaceIdGeneratorsTest.testConflictingIdGenerator()
func (c *ReplaceIdGeneratorsTest) TestConflictingIdGenerator() {
	c.SetExpectParseWarningsInThisTest()
	c.TestSame("/** @idGenerator \n @idGenerator {consistent} \n*/ var xid = function() {}; ")

	c.TestSame("/** @idGenerator {stable} \n @idGenerator \n*/ var xid = function() {}; ")

	c.TestSame(
		`/** @idGenerator {stable}
 @idGenerator {consistent}
 */
var xid = function() {};
`)
}

// Java: ReplaceIdGeneratorsTest.testConsistentIdGenUnderPseudoRenaming()
func (c *ReplaceIdGeneratorsTest) TestConsistentIdGenUnderPseudoRenaming() {
	c.TestWithPseudo(
		`/** @idGenerator {consistent} */ var cid = function() {};
if (x) {foo.bar = cid('foo_bar')}
`,
		`/** @idGenerator {consistent} */ var cid = function() {};
if (x) {foo.bar = 'a'}
`,
		`/** @idGenerator {consistent} */ var cid = function() {};
if (x) {foo.bar = 'foo_bar$0'}
`)
}

// Java: ReplaceIdGeneratorsTest.testUnknownMapping()
func (c *ReplaceIdGeneratorsTest) TestUnknownMapping() {
	c.TestSame(
		`/** @idGenerator {mapped} */
var unknownId = function() {};
function Foo() { unknownId('foo'); }
`,
		ReplaceIdGenerators_MISSING_NAME_MAP_FOR_GENERATOR)
}

// Java: ReplaceIdGeneratorsTest.testBadGenerator1()
func (c *ReplaceIdGeneratorsTest) TestBadGenerator1() {
	c.TestSame(
		`/** @idGenerator */ id = function() {};
foo.bar = id()
`,
		ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER)
}

// Java: ReplaceIdGeneratorsTest.testBadGenerator2()
func (c *ReplaceIdGeneratorsTest) TestBadGenerator2() {
	c.TestSame(
		`/** @idGenerator {consistent} */ id = function() {};
foo.bar = id()
`,
		ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER)
}

// Java: ReplaceIdGeneratorsTest.testBadGenerator3()
func (c *ReplaceIdGeneratorsTest) TestBadGenerator3() {
	c.TestSame(
		"/** @idGenerator {consistent} */ id = function() {};\nfoo.bar = id(`hello${ ' '}world`)\n",
		ReplaceIdGenerators_INVALID_GENERATOR_PARAMETER)
}

// Java: ReplaceIdGeneratorsTest.testTTLCall()
func (c *ReplaceIdGeneratorsTest) TestTTLCall() {
	c.TestMap(
		"/** @idGenerator {stable} */ function tagId(strings) {return strings[0];};\nfunction Foo() { tagId`foo`; }\n",
		`/** @idGenerator {stable} */ function tagId(strings) {return strings[0];};
function Foo() { 'AAGMxg'; }
`,
		`[tagId]

AAGMxg:foo
`)
	c.TestMap(
		"/** @idGenerator {stable} */ function tagId(strings, ...args) {return strings[0];};\nfunction Foo() { tagId`foo${1}bar${2}`; }\n",
		"/** @idGenerator {stable} */ function tagId(strings, ...args) {return strings[0];};\nfunction Foo() { tagId`AAGMxg${1}AAF8Ew${2}AAAAAA`; }\n",
		`[tagId]

AAGMxg:foo
AAF8Ew:bar
AAAAAA:
`)
}

// The application may or may not be transpiling, ensure we support transpiled mode
// Java: ReplaceIdGeneratorsTest.testTTLCall_transpiled()
func (c *ReplaceIdGeneratorsTest) TestTTLCall_transpiled() {
	c.EnableTranspile()
	c.TestMap(
		"/** @idGenerator {stable} */ function tagId(strings) {return strings[0];};\nfunction Foo() { tagId`foo`; }\n",
		`/** @const */
var $jscomp = $jscomp || {};
/** @const */
$jscomp.scope = {};
/**
 * @nosideeffects
 * @noinline
 * @param {!ITemplateArray} arrayStrings
 * @return {!ITemplateArray}
 */
$jscomp.createTemplateTagFirstArg = function(arrayStrings) {
  return $jscomp.createTemplateTagFirstArgWithRaw(arrayStrings, arrayStrings)
};
/**
 * @nosideeffects
 * @noinline
 * @param {!ITemplateArray} arrayStrings
 * @param {!ITemplateArray} rawArrayStrings
 * @return {!ITemplateArray}
 */
$jscomp.createTemplateTagFirstArgWithRaw = function(arrayStrings, rawArrayStrings) {
  arrayStrings.raw = rawArrayStrings;
  Object.freeze && (Object.freeze(arrayStrings), Object.freeze(rawArrayStrings));
  return (/** @type {!ITemplateArray} */ (arrayStrings));
};
/** @idGenerator {stable} */ function tagId(strings) {return strings[0];};
function Foo() { 'AAGMxg'; }
`,
		`[tagId]

AAGMxg:foo
`)
	c.EnableTranspile()
	c.TestMap(
		"/** @idGenerator {stable} */ function tagId(strings, ...args) {return strings[0];};\nfunction Foo() { tagId`foo${1}bar${2}`; }\n",
		"/** @idGenerator {stable} */ function tagId(strings, ...args) {return strings[0];};\nfunction Foo() { tagId`AAGMxg${1}AAF8Ew${2}AAAAAA`; }\n",
		`[tagId]

AAGMxg:foo
AAF8Ew:bar
AAAAAA:
`)
}

// Java: ReplaceIdGeneratorsTest.testTTLCall_invalid()
func (c *ReplaceIdGeneratorsTest) TestTTLCall_invalid() {
	c.TestSame(
		"/** @idGenerator {stable} */ function tagId(strings) {};\n"+
			// invalid escape sequence.  We don't support this though we could.
			"function Foo() { tagId`$\\underline{u}$`;}\n",
		ReplaceIdGenerators_INVALID_TEMPLATE_LITERAL_PARAMETER)
}

// Java: ReplaceIdGeneratorsTest.testMap(String, String, String)
func (c *ReplaceIdGeneratorsTest) TestMap(code string, expected string, expectedMap string) {
	c.Test(code, expected)
	AssertThat(jutil.StringTrim(c.lastPass.GetSerializedIdMappings())).IsEqualTo(jutil.StringTrim(expectedMap))
}

// Java: ReplaceIdGeneratorsTest.testWithPseudo(String, String, String)
func (c *ReplaceIdGeneratorsTest) TestWithPseudo(code string, expected string, expectedPseudo string) {
	c.generatePseudoNames = false
	c.Test(code, expected)
	c.generatePseudoNames = true
	c.Test(code, expectedPseudo)
}

// Java: ReplaceIdGeneratorsTest.testNonPseudoSupportingGenerator(String, String)
func (c *ReplaceIdGeneratorsTest) TestNonPseudoSupportingGenerator(code string, expected string) {
	c.generatePseudoNames = false
	c.Test(code, expected)
	c.generatePseudoNames = true
	c.Test(code, expected)
}

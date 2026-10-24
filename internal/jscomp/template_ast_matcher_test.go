// Port of com/google/javascript/jscomp/TemplateAstMatcherTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/TemplateAstMatcherTest.java)

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TemplateAstMatcherTest holds tests for TemplateAstMatcher.
type TemplateAstMatcherTest struct {
	t            *testing.T
	lastCompiler *Compiler
}

func (c *TemplateAstMatcherTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *TemplateAstMatcherTest) JUnitBefore()           {}
func (c *TemplateAstMatcherTest) JUnitAfter()            {}

func TestTemplateAstMatcherTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &TemplateAstMatcherTest{} })
}

// Java: TemplateAstMatcherTest.testMatches_primitives()
func (c *TemplateAstMatcherTest) TestMatches_primitives() {

	template :=
		`function template() {
  3;
}
`

	pair := c.compile("", template, "3")
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	pair = c.compile("", template, "5")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile("", template, "var foo = 3;")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstFirstChild().GetFirstChild())
	pair = c.compile("", template, "obj.foo();")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())

	template =
		`function template() {
  'str';
}
`
	pair = c.compile("", template, "'str'")
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile("", template, "'not_str'")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile("", template, "var foo = 'str';")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstFirstChild().GetFirstChild())
	pair = c.compile("", template, "obj.foo();")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())

	template =
		`function template() {
  true;
}
`
	pair = c.compile("", template, "true")
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile("", template, "!true")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile("", template, "false")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile("", template, "var foo = true;")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstFirstChild().GetFirstChild())
	pair = c.compile("", template, "!undefined")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
}

// Java: TemplateAstMatcherTest.testMatches_varDeclarations()
func (c *TemplateAstMatcherTest) TestMatches_varDeclarations() {

	template :=
		`function template() {
  var a = 3;
}
`
	var pair *templateAstMatcherTest_TestNodePair

	pair = c.compile("", template, "var a = 3;")
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstChild())
	// Make sure to let variable names differ in var declarations.
	pair = c.compile("", template, "var b = 3;")
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstChild())
	pair = c.compile("", template, "var a = 5;")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	pair = c.compile("", template, "5;")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())

	template =
		`function template() {
  var a = {};
}
`
	pair = c.compile("", template, "var a = {};")
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstChild())
	pair = c.compile("", template, "var a = {'a': 'b'};")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())

	template =
		`function template() {
  var a = {
    'a': 'b'
  };
}
`
	pair = c.compile("", template, "var a = {'a': 'b'};")
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstChild())
	pair = c.compile("", template, "var a = {};")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
}

// Java: TemplateAstMatcherTest.testMatches_templateParameterType()
func (c *TemplateAstMatcherTest) TestMatches_templateParameterType() {

	externs := ""
	template :=
		`/**
 * @param {string} foo
 */
function template(foo) {
  foo;
}
`

	pair := c.compile(externs, template, "'str'")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "'different_str'")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "var foo = 'str';")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstFirstChild().GetFirstChild())

	template =
		`/**
 * @param {*} foo
 */
function template(foo) {
  foo;
}
`
	pair = c.compile(externs, template, "'str'")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "3")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "new Object()")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())

	template =
		`/**
 * @param {string} foo
 * @param {number} bar
 */
function template(foo, bar) {
  bar + foo;
}
`
	pair = c.compile(externs, template, "'str'")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "3")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "new Object()")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "3 + ''")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "7 + 'str'")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())

	template =
		`/**
 * @param {string} string_literal_foo
 */
function template(string_literal_foo) {
  string_literal_foo;
}
`
	pair = c.compile("", template, "\"foo\"")
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile("", template, "\"foo\" + \"bar\"")
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile("", template, "3")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile("", template, "var s = \"3\"; s;")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstFirstChild().GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetSecondChild())
}

// Java: TemplateAstMatcherTest.testMatches_againstUnresolvedType_fromTemplate_withImplicitNullability_noMatch()
func (c *TemplateAstMatcherTest) TestMatches_againstUnresolvedType_fromTemplate_withImplicitNullability_noMatch() {

	c.assertMatches_againstUnresolvedType_fromTemplate_withTypeModifier("")
}

// Java: TemplateAstMatcherTest.testMatches_againstUnresolvedType_fromTemplate_withBangNullability_noMatch()
func (c *TemplateAstMatcherTest) TestMatches_againstUnresolvedType_fromTemplate_withBangNullability_noMatch() {

	c.assertMatches_againstUnresolvedType_fromTemplate_withTypeModifier("!")
}

// Java: TemplateAstMatcherTest.testMatches_againstUnresolvedType_fromTemplate_withQmarkNullability_noMatch()
func (c *TemplateAstMatcherTest) TestMatches_againstUnresolvedType_fromTemplate_withQmarkNullability_noMatch() {

	c.assertMatches_againstUnresolvedType_fromTemplate_withTypeModifier("?")
}

// Java: TemplateAstMatcherTest.testMatches_againstUnresolvedType_fromTemplate_inUnion_noMatch()
func (c *TemplateAstMatcherTest) TestMatches_againstUnresolvedType_fromTemplate_inUnion_noMatch() {

	// TODO(b/146173738): consider making this a match. It's certain that `str` is a subtype of
	// the union `!Some.Missing.Type|string`, even though the compiler does not know what type
	// `Some.Missing.Type` is.
	c.assertMatches_againstUnresolvedType_fromTemplate_withTypeModifier("string|?")
}

// Java: TemplateAstMatcherTest.testMatches_functionCall()
func (c *TemplateAstMatcherTest) TestMatches_functionCall() {

	externs :=
		`function foo() {};
function bar(arg) {};
`
	template :=
		`function template() {
  foo();
}
`
	pair := c.compile(externs, template, "foo();")
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	pair = c.compile(externs, template, "bar();")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	pair = c.compile(externs, template, "bar(foo());")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstFirstChild().GetLastChild())
}

// Java: TemplateAstMatcherTest.testMatches_functionCallWithArguments()
func (c *TemplateAstMatcherTest) TestMatches_functionCallWithArguments() {

	externs :=
		`/** @return {string} */
function foo() {};
/** @param {string} arg */
function bar(arg) {};
/**
 * @param {string} arg
 * @param {number arg2
 */
function baz(arg, arg2) {};
`
	template :=
		`function template() {
  bar('str');
}
`
	pair := c.compile(externs, template, "foo();")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "bar();")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "bar('str');")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "bar(foo());")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild().GetLastChild())

	template =
		`/** @param {string} str */
function template(str) {
  bar(str);
}
`
	pair = c.compile(externs, template, "foo();")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "bar('str');")
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "bar('str' + 'other_str');")
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "bar(String(3));")
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())

	template =
		`/**
 * @param {string} str
 * @param {number} num
 */
function template(str, num) {
  baz(str, num);
}
`
	pair = c.compile(externs, template, "foo();")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "baz('str', 3);")
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "baz('str' + 'other_str', 3 + 4);")
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
}

// Java: TemplateAstMatcherTest.testMatches_methodCall()
func (c *TemplateAstMatcherTest) TestMatches_methodCall() {

	externs :=
		`/** @return {string} */
function foo() {};
`
	template :=
		`/**
 * @param {string} str
 */
function template(str) {
  str.toString();
}
`

	pair := c.compile(externs, template, "'str'")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "'str'.toString()")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "foo().toString()")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
}

// Java: TemplateAstMatcherTest.testMatches_methodCallWithArguments()
func (c *TemplateAstMatcherTest) TestMatches_methodCallWithArguments() {

	externs :=
		`/** @constructor */
function AppContext() {}
AppContext.prototype.init = function() {};
/**
 * @param {string} arg
 */
AppContext.prototype.get = function(arg) {};
/**
 * @param {string} arg
 */
AppContext.prototype.getOrNull = function(arg) {};
`
	template :=
		`/**
 * @param {AppContext} context
 */
function template(context) {
  context.init();
}
`

	pair := c.compile(externs, template, "'str'")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
	pair = c.compile(externs, template, "var context = new AppContext(); context.init();")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstChild())
	pair = c.compile(externs, template, "var context = new AppContext(); context.get('str');")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstChild())

	template =
		`/**
 * @param {AppContext} context
 * @param {string} service
 */
function template(context, service) {
  context.get(service);
}
`

	pair = c.compile(externs, template, "var context = new AppContext(); context.init();")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstChild())
	pair = c.compile(externs, template, "var context = new AppContext(); context.get('s');")
	c.assertMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstChild())
	pair = c.compile(externs, template, "var context = new AppContext(); context.get(3);")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstChild())
	pair = c.compile(externs, template, "var context = new AppContext(); context.getOrNull('s');")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstChild())
}

// Java: TemplateAstMatcherTest.testMatches_instantiation()
func (c *TemplateAstMatcherTest) TestMatches_instantiation() {

	externs :=
		`/** @constructor */
function AppContext() {}
`
	template :=
		`function template() {
  new AppContext();
}
`

	pair := c.compile(externs, template, "var foo = new AppContext()")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstFirstChild().GetFirstChild())
	pair = c.compile(externs, template, "var foo = new Object()")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild().GetFirstChild())
}

// Java: TemplateAstMatcherTest.testMatches_propertyAccess()
func (c *TemplateAstMatcherTest) TestMatches_propertyAccess() {

	externs :=
		`/** @constructor */
function AppContext() {}
/** @type {string} */
AppContext.prototype.location;
`
	template :=
		`/**
 * @param {AppContext} context
 */
function template(context) {
  context.location;
}
`

	pair :=
		c.compile(externs, template, "var context = new AppContext(); context.location = '3';")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstFirstChild())
}

// Java: TemplateAstMatcherTest.testMatches_multiLineTemplates()
func (c *TemplateAstMatcherTest) TestMatches_multiLineTemplates() {

	externs :=
		`/** @constructor */
function AppContext() {}
/** @type {string} */
AppContext.prototype.location;
`
	template :=
		`/**
 * @param {AppContext} context
 * @param {string} str
 */
function template(context, str) {
  context.location = str;
  delete context.location;
}
`

	pair :=
		c.compile(externs, template, "var context = new AppContext(); context.location = '3';")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstFirstChild())

	pair =
		c.compile(
			externs, template, "var ac = new AppContext(); ac.location = '3'; delete ac.location;")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetSecondChild())

	// Ensure that if a variable is declared within the template and reused
	// across multiple statements, ensure that the test code matches the same
	// pattern. For example:
	// a := b()
	// fn(a)
	externs =
		`/** @param {string} arg */
function bar(arg) {};
`
	template =
		`function template() {
  var a = 'string';
  bar(a);
}
`

	pair = c.compile(externs, template, "var loc = 'string'; bar(loc);")
	c.assertMatch(pair.templateNode, pair.testNode.GetFirstChild())
	pair = c.compile(externs, template, "var loc = 'string'; bar('foo');")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	pair = c.compile(externs, template, "var baz = 'qux'; var loc = 'string'; bar(baz);")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetSecondChild())
}

// Java: TemplateAstMatcherTest.testMatches_subclasses()
func (c *TemplateAstMatcherTest) TestMatches_subclasses() {

	externs :=
		`/** @constructor */
function AppContext() {}
/** @type {string} */
AppContext.prototype.location;
/**
 * @constructor
 * @extends {AppContext}
 */
function SubAppContext() {}
`
	template :=
		`/**
 * @param {AppContext} context
 * @param {string} str
 */
function template(context, str) {
  context.location = str;
}
`

	pair :=
		c.compile(externs, template, "var context = new SubAppContext(); context.location = '3';")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstChild())
}

// Java: TemplateAstMatcherTest.testMatches_nonDefaultStrategy()
func (c *TemplateAstMatcherTest) TestMatches_nonDefaultStrategy() {

	externs :=
		`/** @constructor */
function AppContext() {}
/** @type {string} */
AppContext.prototype.location;
/**
 * @constructor
 * @extends {AppContext}
 */
function SubAppContext() {}
var context = new AppContext();
var subContext = new SubAppContext();
`
	template :=
		`/**
 * @param {!AppContext} context
 * @param {string} str
 */
function template(context, str) {
  context.location = str;
}
`

	pair := c.compile(externs, template, "subContext.location = '3';")
	c.assertMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstChild())
	c.assertMatchWithStrategy(
		pair.templateNode,
		pair.testNode.GetLastChild().GetFirstChild(),
		false,
		TypeMatchingStrategy_EXACT)
}

// Java: TemplateAstMatcherTest.testMatches_namespace_method()
func (c *TemplateAstMatcherTest) TestMatches_namespace_method() {

	externs :=
		`/** @const */ const ns = {};
/** @const */ ns.sub = {};
/** @return {boolean} */ ns.sub.method = function() {};
`

	template :=
		`/**
 * @param {typeof ns.sub} target
 */
function template(target) {
  target.method();
}
`

	pair := c.compile(externs, template, "var alias = ns.sub; alias.method();")
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstChild())
	c.assertNotMatch(pair.templateNode, pair.testNode.GetFirstFirstChild())
	c.assertMatch(pair.templateNode, pair.testNode.GetLastChild().GetFirstChild())
}

// Java: TemplateAstMatcherTest.assertMatches_againstUnresolvedType_fromTemplate_withTypeModifier(String)
func (c *TemplateAstMatcherTest) assertMatches_againstUnresolvedType_fromTemplate_withTypeModifier(modifier string) {
	pair := c.compile(
		"",
		strings.ReplaceAll(`/**
 * @param {MODIFIERSome.Missing.Type} foo
 */
function template(foo) {
  foo;
}
`, "MODIFIER", modifier),
		"'str'")
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot())
	c.assertNotMatch(pair.templateNode, pair.GetTestExprResultRoot().GetFirstChild())
}

// Java: TemplateAstMatcherTest.assertMatch(Node, Node, boolean)
func (c *TemplateAstMatcherTest) assertMatchWithShouldMatch(templateRoot rhino.Node, testNode rhino.Node, shouldMatch bool) {
	c.assertMatchWithStrategy(templateRoot, testNode, shouldMatch, TypeMatchingStrategy_LOOSE)
}

// Java: TemplateAstMatcherTest.assertMatch(Node, Node, boolean, TypeMatchingStrategy)
func (c *TemplateAstMatcherTest) assertMatchWithStrategy(
	templateRoot rhino.Node,
	testNode rhino.Node,
	shouldMatch bool,
	typeMatchingStrategy TypeMatchingStrategy) {
	matcher := NewTemplateAstMatcher(c.lastCompiler, templateRoot.GetFirstChild(), typeMatchingStrategy)
	sb := jutil.NewStringBuilder()
	sb.Append("The nodes should")
	if shouldMatch {
		sb.Append("")
	} else {
		sb.Append(" not")
	}
	sb.Append(" have matched.\n")
	sb.Append("Template node:\n").Append(templateRoot.ToStringTree()).Append("\n")
	sb.Append("Test node:\n").Append(testNode.GetParent().ToStringTree()).Append("\n")
	AssertWithMessage(sb.String()).That(matcher.Matches(testNode)).IsEqualTo(shouldMatch)
}

// Java: TemplateAstMatcherTest.assertMatch(Node, Node)
func (c *TemplateAstMatcherTest) assertMatch(templateRoot rhino.Node, testNode rhino.Node) {
	c.assertMatchWithShouldMatch(templateRoot, testNode, true)
}

// Java: TemplateAstMatcherTest.assertNotMatch(Node, Node)
func (c *TemplateAstMatcherTest) assertNotMatch(templateRoot rhino.Node, testNode rhino.Node) {
	c.assertMatchWithShouldMatch(templateRoot, testNode, false)
}

// compile compiles the template and test code. The code must be compiled
// together using the same Compiler in order for the JsSourceMatcher to work
// properly.
//
// Java: TemplateAstMatcherTest.compile(String, String, String)
func (c *TemplateAstMatcherTest) compile(externs string, template string, code string) *templateAstMatcherTest_TestNodePair {
	compiler := NewCompiler()
	c.lastCompiler = compiler
	compiler.DisableThreads()
	options := NewCompilerOptions()
	options.SetCheckTypes(true)
	options.SetChecksOnly(true)
	options.SetPreserveDetailedSourceInfo(true)

	templateNode := compiler.ParseOfSourceFile(SourceFile_FromCode("template", template))

	compiler.Compile(
		[]*SourceFile{SourceFile_FromCode("externs", externs)},
		[]*SourceFile{
			// The extra block allows easier separation of template and test
			// code.
			SourceFile_FromCode("test", "{"+code+"}")},
		options)
	root := compiler.GetRoot()
	return newTemplateAstMatcherTest_TestNodePair(templateNode, root.GetLastChild())
}

// templateAstMatcherTest_TestNodePair is TemplateAstMatcherTest.TestNodePair.
type templateAstMatcherTest_TestNodePair struct {
	templateNode rhino.Node
	testNode     rhino.Node
}

// Java: TemplateAstMatcherTest.TestNodePair(Node, Node)
func newTemplateAstMatcherTest_TestNodePair(template rhino.Node, root rhino.Node) *templateAstMatcherTest_TestNodePair {
	return &templateAstMatcherTest_TestNodePair{
		templateNode: template,
		testNode:     root.GetLastChild().GetFirstChild(),
	}
}

// Java: TemplateAstMatcherTest.TestNodePair.getTestExprResultRoot()
func (p *templateAstMatcherTest_TestNodePair) GetTestExprResultRoot() rhino.Node {
	if p.testNode.GetFirstChild().IsExprResult() {
		return p.testNode.GetFirstChild()
	}
	return rhino.NilNode
}

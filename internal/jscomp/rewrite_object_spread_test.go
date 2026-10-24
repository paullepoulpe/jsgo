// Port of com/google/javascript/jscomp/RewriteObjectSpreadTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/RewriteObjectSpreadTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteObjectSpreadTest has the test cases for RewriteObjectSpread
type RewriteObjectSpreadTest struct {
	CompilerTestCaseBase
}

// Java: RewriteObjectSpreadTest()
func NewRewriteObjectSpreadTest() *RewriteObjectSpreadTest {
	c := &RewriteObjectSpreadTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, NewTestExternsBuilder().AddObject().Build())
	return c
}

func TestRewriteObjectSpreadTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteObjectSpreadTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp() and enableTypeCheckBeforePass().
func (c *RewriteObjectSpreadTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
	c.EnableTypeCheckBeforePass()
}

// Java: RewriteObjectSpreadTest.customSetUp()
func (c *RewriteObjectSpreadTest) CustomSetUp() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2017)
}

// Java: RewriteObjectSpreadTest.enableTypeCheckBeforePass()
func (c *RewriteObjectSpreadTest) EnableTypeCheckBeforePass() {
	c.EnableNormalize()
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.AllowExternsChanges()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
}

// Java: RewriteObjectSpreadTest.getProcessor(Compiler)
func (c *RewriteObjectSpreadTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRewriteObjectSpread(compiler)
}

// Java: RewriteObjectSpreadTest.testObjectLiteralWithSpread()
func (c *RewriteObjectSpreadTest) TestObjectLiteralWithSpread() {
	c.Test("({first, ...spread});", "Object.assign({}, {first}, spread)")
	c.Test("({first, second, ...spread});", "Object.assign({}, {first, second}, spread)")
	c.Test("({...spread, last});", "Object.assign({}, spread, {last})")
	c.Test("({...spread, penultimate, last});", "Object.assign({}, spread, {penultimate, last})")
	c.Test(
		"({before, ...spread1, mid1, mid2, ...spread2, after});",
		"Object.assign({}, {before}, spread1, {mid1, mid2}, spread2, {after})")
	c.Test("({first, ...{...nested}});", "Object.assign({}, {first}, Object.assign({}, nested))")
	c.Test(
		"({first, [foo()]: baz(), ...spread});",
		"Object.assign({}, {first, [foo()]: baz()}, spread)")
	c.Test(
		"({type : 'object', ...spread && spread});",
		"Object.assign({}, {type:'object'}, spread && spread);")
}

// Java: RewriteObjectSpreadTest.testTyping_ofSpreadResult_isObject()
func (c *RewriteObjectSpreadTest) TestTyping_ofSpreadResult_isObject() {
	c.Test(
		`const first = 0;
const spread = {bar: 'str', qux: false};
const obj = ({first, ...spread});
`,
		`const first = 0;
const spread = {bar: 'str', qux: false};
const obj = Object.assign({}, {first}, spread)
`)

	lastCompiler := c.GetLastCompiler()

	obj := c.GetNodeMatchingQName(lastCompiler.GetJsRoot(), "obj")
	AssertNode(obj).HasColorThat().IsEqualTo(colors.StandardColors_TOP_OBJECT)
	AssertNode(obj.GetFirstFirstChild()).
		HasColorThat().
		IsEqualTo(
			CodeSubTree_FindFirstNode(
				c.GetLastCompiler().GetExternsRoot(),
				func(node rhino.Node) bool { return node.MatchesQualifiedName("Object.assign") }).
				GetNext().
				GetColor())
}

// Returns the first node (preorder) in the given AST that matches the given qualified name
// Java: RewriteObjectSpreadTest.getNodeMatchingQName(Node, String)
func (c *RewriteObjectSpreadTest) GetNodeMatchingQName(root rhino.Node, qname string) rhino.Node {
	if root.MatchesQualifiedName(qname) {
		return root
	}
	for child := root.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		result := c.GetNodeMatchingQName(child, qname)
		if !result.IsNil() {
			return result
		}
	}
	return rhino.NilNode
}

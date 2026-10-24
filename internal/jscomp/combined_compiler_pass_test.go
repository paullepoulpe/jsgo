// Port of com/google/javascript/jscomp/CombinedCompilerPassTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CombinedCompilerPassTest.java)

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CombinedCompilerPassTest holds tests for CombinedCompilerPass.
type CombinedCompilerPassTest struct {
	t        *testing.T
	compiler *Compiler
}

func (c *CombinedCompilerPassTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *CombinedCompilerPassTest) JUnitBefore()           { c.SetUp() }
func (c *CombinedCompilerPassTest) JUnitAfter()            {}

func TestCombinedCompilerPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &CombinedCompilerPassTest{} })
}

// combinedCompilerPassTest_createPostOrderAlphabet returns a Node tree with
// the post-order traversal a b c d e f g h i j k l m and the in-order
// traversal m d a b c h e f g l i j k:
//
//	          m
//	,---------|---------.
//	d         h         l
//	,--|--. ,--|--. ,--|--.
//	a b c   e f g   i j k
//
// Java: CombinedCompilerPassTest.createPostOrderAlphabet()
func combinedCompilerPassTest_createPostOrderAlphabet() *rhino.Node {
	a := rhino.Node_NewString("a")
	b := rhino.Node_NewString("b")
	c := rhino.Node_NewString("c")
	d := rhino.Node_NewString("d")
	e := rhino.Node_NewString("e")
	f := rhino.Node_NewString("f")
	g := rhino.Node_NewString("g")
	h := rhino.Node_NewString("h")
	i := rhino.Node_NewString("i")
	j := rhino.Node_NewString("j")
	k := rhino.Node_NewString("k")
	l := rhino.Node_NewString("l")
	m := rhino.Node_NewString("m")

	d.AddChildToBack(a)
	d.AddChildToBack(b)
	d.AddChildToBack(c)

	h.AddChildToBack(e)
	h.AddChildToBack(f)
	h.AddChildToBack(g)

	l.AddChildToBack(i)
	l.AddChildToBack(j)
	l.AddChildToBack(k)

	m.AddChildToBack(d)
	m.AddChildToBack(h)
	m.AddChildToBack(l)

	return m
}

// Java: CombinedCompilerPassTest.setUp()
func (c *CombinedCompilerPassTest) SetUp() {
	c.compiler = NewCompiler()
	c.compiler.InitOptions(NewCompilerOptions())
}

// combinedCompilerPassTest_ConcatTraversal concatenates contents of string
// nodes encountered in pre-order and post-order traversals. Abbreviates
// traversals by ignoring subtrees rooted with specified strings.
type combinedCompilerPassTest_ConcatTraversal struct {
	visited         *jutil.StringBuilder
	shouldTraversed *jutil.StringBuilder
	ignoring        *jutil.LinkedHashSet[string]
}

func newCombinedCompilerPassTest_ConcatTraversal() *combinedCompilerPassTest_ConcatTraversal {
	return &combinedCompilerPassTest_ConcatTraversal{
		visited:         jutil.NewStringBuilder(),
		shouldTraversed: jutil.NewStringBuilder(),
		ignoring:        jutil.NewLinkedHashSet[string](),
	}
}

// Java: CombinedCompilerPassTest.ConcatTraversal.ignore(String)
func (c *combinedCompilerPassTest_ConcatTraversal) Ignore(s string) *combinedCompilerPassTest_ConcatTraversal {
	c.ignoring.Add(s)
	return c
}

// Java: CombinedCompilerPassTest.ConcatTraversal.visit(NodeTraversal, Node, Node)
func (c *combinedCompilerPassTest_ConcatTraversal) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	AssertThat(n.GetToken()).IsEqualTo(rhino.Token_STRINGLIT)
	c.visited.Append(n.GetString())
}

// Java: CombinedCompilerPassTest.ConcatTraversal.shouldTraverse(NodeTraversal, Node, Node)
func (c *combinedCompilerPassTest_ConcatTraversal) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	AssertThat(n.GetToken()).IsEqualTo(rhino.Token_STRINGLIT)
	c.shouldTraversed.Append(n.GetString())
	return !c.ignoring.Contains(n.GetString())
}

// GetVisited returns strings concatenated during post-order traversal.
//
// Java: CombinedCompilerPassTest.ConcatTraversal.getVisited()
func (c *combinedCompilerPassTest_ConcatTraversal) GetVisited() string {
	return c.visited.String()
}

// GetShouldTraversed returns strings concatenated during pre-order
// traversal.
//
// Java: CombinedCompilerPassTest.ConcatTraversal.getShouldTraversed()
func (c *combinedCompilerPassTest_ConcatTraversal) GetShouldTraversed() string {
	return c.shouldTraversed.String()
}

// Java: CombinedCompilerPassTest.ConcatTraversal.getIgnoring()
func (c *combinedCompilerPassTest_ConcatTraversal) GetIgnoring() *jutil.LinkedHashSet[string] {
	return c.ignoring
}

// combinedCompilerPassTest_TestHelper is a collection of data for a
// traversal test. Contains the traversal callback and the expected pre- and
// post-order traversal results.
type combinedCompilerPassTest_TestHelper struct {
	traversal              *combinedCompilerPassTest_ConcatTraversal
	expectedVisited        string
	shouldTraverseExpected string
}

// Java: CombinedCompilerPassTest.TestHelper.getTraversal()
func (h *combinedCompilerPassTest_TestHelper) GetTraversal() *combinedCompilerPassTest_ConcatTraversal {
	return h.traversal
}

// Java: CombinedCompilerPassTest.TestHelper.checkResults()
func (h *combinedCompilerPassTest_TestHelper) CheckResults() {
	AssertWithMessage(
		"ConcatTraversal ignoring %s has unexpected visiting order", h.traversal.GetIgnoring()).
		That(h.traversal.GetVisited()).
		IsEqualTo(h.expectedVisited)

	AssertWithMessage(
		"ConcatTraversal ignoring %s has unexpected traversal order", h.traversal.GetIgnoring()).
		That(h.traversal.GetShouldTraversed()).
		IsEqualTo(h.shouldTraverseExpected)
}

// Java: CombinedCompilerPassTest.createStringTests()
func combinedCompilerPassTest_createStringTests() []*combinedCompilerPassTest_TestHelper {
	var tests []*combinedCompilerPassTest_TestHelper

	tests = append(tests, &combinedCompilerPassTest_TestHelper{newCombinedCompilerPassTest_ConcatTraversal(), "abcdefghijklm", "mdabchefglijk"})

	tests = append(tests, &combinedCompilerPassTest_TestHelper{newCombinedCompilerPassTest_ConcatTraversal().Ignore("d"), "efghijklm", "mdhefglijk"})

	tests = append(tests, &combinedCompilerPassTest_TestHelper{newCombinedCompilerPassTest_ConcatTraversal().Ignore("f"), "abcdeghijklm", "mdabchefglijk"})

	tests = append(tests, &combinedCompilerPassTest_TestHelper{newCombinedCompilerPassTest_ConcatTraversal().Ignore("m"), "", "m"})

	return tests
}

// Java: CombinedCompilerPassTest.testIndividualPasses()
func (c *CombinedCompilerPassTest) TestIndividualPasses() {
	for _, test := range combinedCompilerPassTest_createStringTests() {
		pass := NewCombinedCompilerPass(c.compiler, test.GetTraversal())
		pass.Process(nil, combinedCompilerPassTest_createPostOrderAlphabet())
		test.CheckResults()
	}
}

// Java: CombinedCompilerPassTest.testCombinedPasses()
func (c *CombinedCompilerPassTest) TestCombinedPasses() {
	tests := combinedCompilerPassTest_createStringTests()
	callbacks := make([]NodeTraversal_Callback, len(tests))
	i := 0
	for _, test := range tests {
		callbacks[i] = test.GetTraversal()
		i++
	}
	pass := NewCombinedCompilerPass(c.compiler, callbacks...)
	pass.Process(nil, combinedCompilerPassTest_createPostOrderAlphabet())
	for _, test := range tests {
		test.CheckResults()
	}
}

// combinedCompilerPassTest_ScopeRecordingCallback records the scopes
// visited during an AST traversal. Abbreviates traversals by ignoring
// subtrees rooted with specified NAME nodes.
type combinedCompilerPassTest_ScopeRecordingCallback struct {
	visitedScopes *jutil.LinkedHashSet[*rhino.Node]
	ignoring      *jutil.LinkedHashSet[string]
}

func newCombinedCompilerPassTest_ScopeRecordingCallback() *combinedCompilerPassTest_ScopeRecordingCallback {
	return &combinedCompilerPassTest_ScopeRecordingCallback{
		visitedScopes: jutil.NewLinkedHashSet[*rhino.Node](),
		ignoring:      jutil.NewLinkedHashSet[string](),
	}
}

// Java: CombinedCompilerPassTest.ScopeRecordingCallback.ignore(String)
func (c *combinedCompilerPassTest_ScopeRecordingCallback) Ignore(name string) {
	c.ignoring.Add(name)
}

// Java: CombinedCompilerPassTest.ScopeRecordingCallback.enterScope(NodeTraversal)
func (c *combinedCompilerPassTest_ScopeRecordingCallback) EnterScope(t *NodeTraversal) {
	c.visitedScopes.Add(t.GetScopeRoot())
}

// Java: CombinedCompilerPassTest.ScopeRecordingCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *combinedCompilerPassTest_ScopeRecordingCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return !n.IsName() || !c.ignoring.Contains(n.GetString())
}

// Java: CombinedCompilerPassTest.ScopeRecordingCallback.getVisitedScopes()
func (c *combinedCompilerPassTest_ScopeRecordingCallback) GetVisitedScopes() *jutil.LinkedHashSet[*rhino.Node] {
	return c.visitedScopes
}

// Java: CombinedCompilerPassTest.ScopeRecordingCallback.exitScope(NodeTraversal)
func (c *combinedCompilerPassTest_ScopeRecordingCallback) ExitScope(t *NodeTraversal) {}

// Java: CombinedCompilerPassTest.ScopeRecordingCallback.visit(NodeTraversal, Node, Node)
func (c *combinedCompilerPassTest_ScopeRecordingCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
}

// Java: CombinedCompilerPassTest.testScopes()
func (c *CombinedCompilerPassTest) TestScopes() {
	root := c.compiler.ParseTestCode("var y = function() { var x = function() { };}")

	c1 := newCombinedCompilerPassTest_ScopeRecordingCallback()
	c1.Ignore("y")
	c2 := newCombinedCompilerPassTest_ScopeRecordingCallback()
	c2.Ignore("x")
	c3 := newCombinedCompilerPassTest_ScopeRecordingCallback()

	pass := NewCombinedCompilerPass(c.compiler, c1, c2, c3)
	pass.Process(nil, root)

	AssertThat(c1.GetVisitedScopes().Size()).IsEqualTo(1)
	AssertThat(c2.GetVisitedScopes().Size()).IsEqualTo(3)
	AssertThat(c3.GetVisitedScopes().Size()).IsEqualTo(5)
}

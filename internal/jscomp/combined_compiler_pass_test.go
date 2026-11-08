// Port of com/google/javascript/jscomp/CombinedCompilerPassTest.java @ closure-compiler v20261104
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
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/ast"
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
func combinedCompilerPassTest_createPostOrderAlphabet() rhino.Node {
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
	visited         strings.Builder
	shouldTraversed strings.Builder
	ignoring        *OrderedSet[string]
}

func newCombinedCompilerPassTest_ConcatTraversal() *combinedCompilerPassTest_ConcatTraversal {
	return &combinedCompilerPassTest_ConcatTraversal{
		ignoring: NewOrderedSet[string](),
	}
}

// Java: CombinedCompilerPassTest.ConcatTraversal.ignore(String)
func (c *combinedCompilerPassTest_ConcatTraversal) Ignore(s string) *combinedCompilerPassTest_ConcatTraversal {
	c.ignoring.Add(s)
	return c
}

// Java: CombinedCompilerPassTest.ConcatTraversal.visit(NodeTraversal, Node, Node)
func (c *combinedCompilerPassTest_ConcatTraversal) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	AssertThat(n.GetToken()).IsEqualTo(rhino.Token_STRINGLIT)
	c.visited.WriteString(n.GetString())
}

// Java: CombinedCompilerPassTest.ConcatTraversal.shouldTraverse(NodeTraversal, Node, Node)
func (c *combinedCompilerPassTest_ConcatTraversal) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	AssertThat(n.GetToken()).IsEqualTo(rhino.Token_STRINGLIT)
	c.shouldTraversed.WriteString(n.GetString())
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
func (c *combinedCompilerPassTest_ConcatTraversal) GetIgnoring() *OrderedSet[string] {
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
		pass.Process(rhino.NilNode, combinedCompilerPassTest_createPostOrderAlphabet())
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
	pass.Process(rhino.NilNode, combinedCompilerPassTest_createPostOrderAlphabet())
	for _, test := range tests {
		test.CheckResults()
	}
}

// combinedCompilerPassTest_ScopeRecordingCallback records the scopes
// visited during an AST traversal. Abbreviates traversals by ignoring
// subtrees rooted with specified NAME nodes.
type combinedCompilerPassTest_ScopeRecordingCallback struct {
	visitedScopes *OrderedSet[rhino.Node]
	ignoring      *OrderedSet[string]
}

func newCombinedCompilerPassTest_ScopeRecordingCallback() *combinedCompilerPassTest_ScopeRecordingCallback {
	return &combinedCompilerPassTest_ScopeRecordingCallback{
		visitedScopes: NewOrderedSet[rhino.Node](),
		ignoring:      NewOrderedSet[string](),
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
func (c *combinedCompilerPassTest_ScopeRecordingCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return !n.IsName() || !c.ignoring.Contains(n.GetString())
}

// Java: CombinedCompilerPassTest.ScopeRecordingCallback.getVisitedScopes()
func (c *combinedCompilerPassTest_ScopeRecordingCallback) GetVisitedScopes() *OrderedSet[rhino.Node] {
	return c.visitedScopes
}

// Java: CombinedCompilerPassTest.ScopeRecordingCallback.exitScope(NodeTraversal)
func (c *combinedCompilerPassTest_ScopeRecordingCallback) ExitScope(t *NodeTraversal) {}

// Java: CombinedCompilerPassTest.ScopeRecordingCallback.visit(NodeTraversal, Node, Node)
func (c *combinedCompilerPassTest_ScopeRecordingCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
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
	pass.Process(rhino.NilNode, root)

	AssertThat(c1.GetVisitedScopes().Size()).IsEqualTo(1)
	AssertThat(c2.GetVisitedScopes().Size()).IsEqualTo(3)
	AssertThat(c3.GetVisitedScopes().Size()).IsEqualTo(5)
}

type testFilteredCallback struct {
	id         string
	tokens     *ast.TokenSet
	jsdoc      bool
	noST       bool
	onlySF     bool
	noVisit    bool
	shouldFn   func(t *NodeTraversal, n, parent rhino.Node) bool
	visitFn    func(t *NodeTraversal, n, parent rhino.Node)
	enterFn    func(t *NodeTraversal)
	events     []string
	sawNilMask bool
}

func (cb *testFilteredCallback) ShouldTraverse(t *NodeTraversal, n, parent rhino.Node) bool {
	cb.events = append(cb.events, fmt.Sprintf("ST:%s:%s", cb.id, n.GetToken()))
	if cb.shouldFn != nil {
		return cb.shouldFn(t, n, parent)
	}
	return true
}

func (cb *testFilteredCallback) Visit(t *NodeTraversal, n, parent rhino.Node) {
	if t.visitTokens == nil {
		cb.sawNilMask = true
	}
	cb.events = append(cb.events, fmt.Sprintf("V:%s:%s:jsdoc=%v", cb.id, n.GetToken(), n.HasJSDocInfo()))
	if cb.visitFn != nil {
		cb.visitFn(t, n, parent)
	}
}

func (cb *testFilteredCallback) EnterScope(t *NodeTraversal) {
	cb.events = append(cb.events, fmt.Sprintf("EN:%s:%s", cb.id, t.GetScopeRoot().GetToken()))
	if cb.enterFn != nil {
		cb.enterFn(t)
	}
}

func (cb *testFilteredCallback) ExitScope(t *NodeTraversal) {
	cb.events = append(cb.events, fmt.Sprintf("EX:%s:%s", cb.id, t.GetScopeRoot().GetToken()))
}

func (cb *testFilteredCallback) visitTokens() *ast.TokenSet { return cb.tokens }
func (cb *testFilteredCallback) visitJSDoc() bool           { return cb.jsdoc }
func (cb *testFilteredCallback) shouldTraverseOnlyScriptAndFunction() bool {
	return cb.onlySF
}

type testNoSTFilteredCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	tokens  *ast.TokenSet
	jsdoc   bool
	visitFn func(t *NodeTraversal, n, parent rhino.Node)
	events  []string
}

func (cb *testNoSTFilteredCallback) Visit(t *NodeTraversal, n, parent rhino.Node) {
	cb.events = append(cb.events, fmt.Sprintf("V:%s:jsdoc=%v", n.GetToken(), n.HasJSDocInfo()))
	if cb.visitFn != nil {
		cb.visitFn(t, n, parent)
	}
}

func (cb *testNoSTFilteredCallback) visitTokens() *ast.TokenSet { return cb.tokens }
func (cb *testNoSTFilteredCallback) visitJSDoc() bool           { return cb.jsdoc }

type testPreOrderCallback struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	shouldFn func(t *NodeTraversal, n, parent rhino.Node) bool
	events   []string
}

func (cb *testPreOrderCallback) ShouldTraverse(t *NodeTraversal, n, parent rhino.Node) bool {
	cb.events = append(cb.events, fmt.Sprintf("ST:%s", n.GetToken()))
	if cb.shouldFn != nil {
		return cb.shouldFn(t, n, parent)
	}
	return true
}

func attachTestJSDoc(n rhino.Node) {
	b := rhino.JSDocInfo_Builder_()
	b.RecordConstancy()
	n.SetJSDocInfo(b.Build())
}

func TestCombinedCompilerPass_VisitJSDocAndTokenFilter(t *testing.T) {
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	root := compiler.ParseTestCode("function f(a) { var x = 1 + 2; debugger; return f(a); }")

	// Attach JSDoc to both a branch node (ADD) and normal leaf nodes (NUMBER, DEBUGGER).
	var addNode, numNode, dbgNode rhino.Node
	NodeUtil_VisitPreOrder(root, func(n rhino.Node) {
		switch n.GetToken() {
		case rhino.Token_ADD:
			addNode = n
		case rhino.Token_NUMBER:
			if numNode.IsNil() {
				numNode = n
			}
		case rhino.Token_DEBUGGER:
			dbgNode = n
		}
	})
	attachTestJSDoc(addNode)
	attachTestJSDoc(numNode)
	attachTestJSDoc(dbgNode)

	callSet := makeVisitTokenSet(rhino.Token_CALL)
	fnSet := makeVisitTokenSet(rhino.Token_FUNCTION)

	makeCallbacks := func() (*testNoSTFilteredCallback, *testNoSTFilteredCallback, *testNoSTFilteredCallback) {
		jsdocOnly := &testNoSTFilteredCallback{tokens: &emptyTokenSet, jsdoc: true}
		callOnly := &testNoSTFilteredCallback{tokens: &callSet, jsdoc: false}
		fnAndJSDoc := &testNoSTFilteredCallback{tokens: &fnSet, jsdoc: true}
		return jsdocOnly, callOnly, fnAndJSDoc
	}

	s1, s2, s3 := makeCallbacks()
	NodeTraversal_Traverse(compiler, root, s1)
	NodeTraversal_Traverse(compiler, root, s2)
	NodeTraversal_Traverse(compiler, root, s3)

	if len(s1.events) != 3 {
		t.Fatalf("expected jsdocOnly to visit exactly 3 JSDoc nodes, got %v", s1.events)
	}

	c1, c2, c3 := makeCallbacks()
	sawUnfiltered := false
	c1.visitFn = func(nt *NodeTraversal, _, _ rhino.Node) {
		if nt.visitTokens == nil || !nt.visitJSDoc || !nt.noShouldTraverse || !nt.noBranchShouldTraverse {
			sawUnfiltered = true
		}
	}
	NewCombinedCompilerPass(compiler, c1, c2, c3).Process(rhino.NilNode, root)

	if sawUnfiltered {
		t.Fatalf("expected CombinedCompilerPass to keep visitTokens, visitJSDoc, noShouldTraverse, and noBranchShouldTraverse enabled")
	}
	if !slices.Equal(c1.events, s1.events) {
		t.Fatalf("jsdocOnly mismatch:\ngot:  %v\nwant: %v", c1.events, s1.events)
	}
	if !slices.Equal(c2.events, s2.events) {
		t.Fatalf("callOnly mismatch:\ngot:  %v\nwant: %v", c2.events, s2.events)
	}
	if !slices.Equal(c3.events, s3.events) {
		t.Fatalf("fnAndJSDoc mismatch:\ngot:  %v\nwant: %v", c3.events, s3.events)
	}
}

func TestCombinedCompilerPass_WaitingTokenMaskParityAndActiveFilters(t *testing.T) {
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	s1 := compiler.ParseSyntheticCode("s1.js", "function outer() { function inner() { callInner(); } callOuter(); } function sib() { return 1; }")
	s2 := compiler.ParseSyntheticCode("s2.js", "var z = 1; z = 2;")
	root := rhino.IR_Root(s1, s2)

	// Attach JSDoc to a leaf inside inner() and a statement in s2.
	NodeUtil_VisitPreOrder(root, func(n rhino.Node) {
		if n.IsNumber() {
			attachTestJSDoc(n)
		}
	})

	modSet := makeVisitTokenSet(rhino.Token_MODULE_BODY)
	callSet := makeVisitTokenSet(rhino.Token_CALL)
	retSet := makeVisitTokenSet(rhino.Token_RETURN)
	assignSet := makeVisitTokenSet(rhino.Token_ASSIGN)

	buildSuite := func() []NodeTraversal_Callback {
		// Callback 0: deactivates on s1 (SCRIPT), active on s2.
		cb0 := &testFilteredCallback{
			id:     "cb0",
			tokens: &modSet,
			shouldFn: func(_ *NodeTraversal, n, _ rhino.Node) bool {
				return n != s1
			},
		}
		// Callback 1: deactivates on outer() FUNCTION.
		cb1 := &testFilteredCallback{
			id:     "cb1",
			tokens: &callSet,
			jsdoc:  true,
			shouldFn: func(_ *NodeTraversal, n, _ rhino.Node) bool {
				return !n.IsFunction() || n.GetFirstChild().GetString() != "outer"
			},
		}
		// Callback 2: deactivates on inner() FUNCTION (same Token_FUNCTION as outer!).
		cb2 := &testFilteredCallback{
			id:     "cb2",
			tokens: &callSet,
			shouldFn: func(_ *NodeTraversal, n, _ rhino.Node) bool {
				return !n.IsFunction() || n.GetFirstChild().GetString() != "inner"
			},
		}
		// Callback 3: post-order only + JSDoc.
		cb3 := &testNoSTFilteredCallback{
			tokens: &retSet,
			jsdoc:  true,
		}
		// Callback 4: onlyScriptAndFunction ShouldTraverse.
		cb4 := &testFilteredCallback{
			id:     "cb4",
			tokens: &assignSet,
			onlySF: true,
		}
		return []NodeTraversal_Callback{cb0, cb1, cb2, cb3, cb4}
	}

	standalone := buildSuite()
	for _, cb := range standalone {
		NodeTraversal_Traverse(compiler, root, cb)
	}

	combined := buildSuite()
	NewCombinedCompilerPassOfList(compiler, combined).Process(rhino.NilNode, root)

	for i := range standalone {
		switch s := standalone[i].(type) {
		case *testFilteredCallback:
			c := combined[i].(*testFilteredCallback)
			if c.sawNilMask {
				t.Fatalf("callback %d saw nil t.visitTokens while subset callbacks were waiting", i)
			}
			if !slices.Equal(c.events, s.events) {
				t.Fatalf("callback %d mismatch:\ngot:  %v\nwant: %v", i, c.events, s.events)
			}
		case *testNoSTFilteredCallback:
			c := combined[i].(*testNoSTFilteredCallback)
			if !slices.Equal(c.events, s.events) {
				t.Fatalf("callback %d mismatch:\ngot:  %v\nwant: %v", i, c.events, s.events)
			}
		}
	}
}

func TestCombinedCompilerPass_NoShouldTraverseWhenDeactivatedAndLeafReactivation(t *testing.T) {
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	s1 := compiler.ParseSyntheticCode("s1.js", "var a = 1; function fn() { var b = a; }")
	s2 := compiler.ParseSyntheticCode("s2.js", "var c = 2;")
	root := rhino.IR_Root(s1, s2)

	scriptSet := makeVisitTokenSet(rhino.Token_MODULE_BODY)
	nameSet := makeVisitTokenSet(rhino.Token_NAME)

	pruner := &testFilteredCallback{
		id:     "pruner",
		tokens: &scriptSet,
		shouldFn: func(_ *NodeTraversal, n, _ rhino.Node) bool {
			return n != s1
		},
	}
	sawDisabledNoSTInS1 := false
	visitor := &testNoSTFilteredCallback{
		tokens: &nameSet,
		visitFn: func(nt *NodeTraversal, n, _ rhino.Node) {
			if nt.GetCurrentScript() == s1 && (!nt.noShouldTraverse || !nt.noBranchShouldTraverse || nt.visitTokens == nil) {
				sawDisabledNoSTInS1 = true
			}
		},
	}

	NewCombinedCompilerPass(compiler, pruner, visitor).Process(rhino.NilNode, root)
	if sawDisabledNoSTInS1 {
		t.Fatalf("expected noShouldTraverse, noBranchShouldTraverse, and visitTokens to stay active inside s1 when pruner is deactivated at s1")
	}
	// Verify pruner reactivated at exit of s1 and ran ShouldTraverse on s2.
	sawS2 := false
	for _, ev := range pruner.events {
		if ev == "ST:pruner:SCRIPT" {
			// Should be called on BOTH s1 and s2.
			if sawS2 {
				break
			}
			sawS2 = true
		}
	}
	stScriptCount := 0
	for _, ev := range pruner.events {
		if ev == "ST:pruner:SCRIPT" {
			stScriptCount++
		}
	}
	if stScriptCount != 2 {
		t.Fatalf("expected pruner.ShouldTraverse to run on both SCRIPT nodes, got count=%d events=%v", stScriptCount, pruner.events)
	}

	// Also verify allNoVisit leaf deactivation/reactivation (e.g. EMPTY `;` and DEBUGGER `debugger;`).
	leafRoot := compiler.ParseTestCode("; debugger; var afterLeaf = 1;")
	p1 := &testPreOrderCallback{
		shouldFn: func(_ *NodeTraversal, n, _ rhino.Node) bool {
			return !n.IsEmpty() && !n.IsDebugger()
		},
	}
	p2 := &testPreOrderCallback{}
	NewCombinedCompilerPass(compiler, p1, p2).Process(rhino.NilNode, leafRoot)

	sP1 := &testPreOrderCallback{shouldFn: p1.shouldFn}
	NodeTraversal_Traverse(compiler, leafRoot, sP1)
	if !slices.Equal(p1.events, sP1.events) {
		t.Fatalf("leaf pre-order reactivation mismatch:\ngot:  %v\nwant: %v", p1.events, sP1.events)
	}
}

func TestCombinedCompilerPass_HaltingErrorStopsTraversalWhenNoShouldTraverse(t *testing.T) {
	haltDiag := DiagnosticType_Error("JSC_TEST_HALTING_ERROR", "halting test error")

	t.Run("haltInVisit", func(t *testing.T) {
		compiler := NewCompiler()
		compiler.InitOptions(NewCompilerOptions())
		s1 := compiler.ParseSyntheticCode("s1.js", "first(); second();")
		s2 := compiler.ParseSyntheticCode("s2.js", "third();")
		root := rhino.IR_Root(s1, s2)

		callSet := makeVisitTokenSet(rhino.Token_CALL)
		var visitedCalls []string
		errCb := &testNoSTFilteredCallback{
			tokens: &callSet,
			visitFn: func(nt *NodeTraversal, n, _ rhino.Node) {
				name := n.GetFirstChild().GetString()
				visitedCalls = append(visitedCalls, name)
				if name == "first" {
					nt.Report(n, haltDiag)
				}
			},
		}
		otherCb := &testNoSTFilteredCallback{tokens: &callSet}

		NewCombinedCompilerPass(compiler, errCb, otherCb).Process(rhino.NilNode, root)
		if !slices.Equal(visitedCalls, []string{"first"}) {
			t.Fatalf("expected halting error in Visit to stop before sibling/next script, visited=%v", visitedCalls)
		}
	})

	t.Run("haltInEnterScope", func(t *testing.T) {
		compiler := NewCompiler()
		compiler.InitOptions(NewCompilerOptions())
		root := compiler.ParseTestCode("function bad() { insideBad(); } afterBad();")

		callSet := makeVisitTokenSet(rhino.Token_CALL)
		var visitedCalls []string
		scopeErrCb := &testFilteredCallback{
			id:     "scopeErr",
			tokens: &callSet,
			enterFn: func(nt *NodeTraversal) {
				if nt.GetScopeRoot().IsFunction() {
					nt.Report(nt.GetScopeRoot(), haltDiag)
				}
			},
			visitFn: func(_ *NodeTraversal, n, _ rhino.Node) {
				visitedCalls = append(visitedCalls, n.GetFirstChild().GetString())
			},
		}
		// Wrap a second noShouldTraverse callback; make scopeErrCb also noShouldTraverse via wrapper below.
		noSTVisit := &testNoSTFilteredCallback{
			tokens: &callSet,
			visitFn: func(_ *NodeTraversal, n, _ rhino.Node) {
				visitedCalls = append(visitedCalls, "noST:"+n.GetFirstChild().GetString())
			},
		}
		NewCombinedCompilerPass(compiler, &testScopedNoSTCallback{inner: scopeErrCb}, noSTVisit).Process(rhino.NilNode, root)
		if len(visitedCalls) != 0 {
			t.Fatalf("expected halting error in EnterScope to prevent visiting insideBad/afterBad, visited=%v", visitedCalls)
		}
	})
}

type testScopedNoSTCallback struct {
	NodeTraversal_AbstractScopedCallbackBase
	inner *testFilteredCallback
}

func (s *testScopedNoSTCallback) visitTokens() *ast.TokenSet  { return s.inner.tokens }
func (s *testScopedNoSTCallback) EnterScope(t *NodeTraversal) { s.inner.EnterScope(t) }
func (s *testScopedNoSTCallback) Visit(t *NodeTraversal, n, parent rhino.Node) {
	s.inner.Visit(t, n, parent)
}

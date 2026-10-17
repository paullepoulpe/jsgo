// Port of test/com/google/javascript/jscomp/NodeTraversalTest.java @ closure-compiler v20260415

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// nodeTraversalTestCallback is an anonymous NodeTraversal.Callback.
type nodeTraversalTestCallback struct {
	shouldTraverse func(t *NodeTraversal, n, parent *rhino.Node) bool
	visit          func(t *NodeTraversal, n, parent *rhino.Node)
}

func (c *nodeTraversalTestCallback) ShouldTraverse(t *NodeTraversal, n, parent *rhino.Node) bool {
	if c.shouldTraverse == nil {
		return true
	}
	return c.shouldTraverse(t, n, parent)
}

func (c *nodeTraversalTestCallback) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	if c.visit != nil {
		c.visit(t, n, parent)
	}
}

// nodeTraversalTestScopedCallback is an anonymous NodeTraversal.ScopedCallback.
type nodeTraversalTestScopedCallback struct {
	nodeTraversalTestCallback
	enterScope func(t *NodeTraversal)
	exitScope  func(t *NodeTraversal)
}

func (c *nodeTraversalTestScopedCallback) EnterScope(t *NodeTraversal) {
	if c.enterScope != nil {
		c.enterScope(t)
	}
}

func (c *nodeTraversalTestScopedCallback) ExitScope(t *NodeTraversal) {
	if c.exitScope != nil {
		c.exitScope(t)
	}
}

// nodeTraversalTest_ErrorManager is the tests' anonymous BasicErrorManager
// that collects the reported errors.
type nodeTraversalTest_ErrorManager struct {
	BasicErrorManagerBase
	errors *[]*JSError
}

func newNodeTraversalTest_ErrorManager(errors *[]*JSError) *nodeTraversalTest_ErrorManager {
	m := &nodeTraversalTest_ErrorManager{errors: errors}
	m.InitBasicErrorManagerBase(m)
	return m
}

func (m *nodeTraversalTest_ErrorManager) Report(level *CheckLevel, error_ *JSError) {
	*m.errors = append(*m.errors, error_)
}

func (m *nodeTraversalTest_ErrorManager) Println(level *CheckLevel, error_ *JSError) {}

func (m *nodeTraversalTest_ErrorManager) PrintSummary() {}

func TestNodeTraversal_Report(t *testing.T) {
	dt := DiagnosticType_Warning("FOO", "{0}, {1} - {2}")

	callback := &nodeTraversalTestCallback{
		shouldTraverse: func(nt *NodeTraversal, n, parent *rhino.Node) bool {
			nt.Report(n, dt, "Foo", "Bar", "Hello")
			return false
		},
		visit: func(nt *NodeTraversal, n, parent *rhino.Node) {
			panic(jutil.NewAssertionErrorNoMessage())
		},
	}

	var errors []*JSError
	compiler := NewCompiler()
	compiler.SetErrorManager(newNodeTraversalTest_ErrorManager(&errors))
	compiler.InitCompilerOptionsIfTesting()

	NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(callback).
		Traverse(rhino.NewNode(rhino.Token_EMPTY))

	if len(errors) != 1 {
		t.Fatalf("got %d errors, want 1", len(errors))
	}
	assertEqual(t, "description", errors[0].Description(), "Foo, Bar - Hello")
}

func TestNodeTraversal_ReportWithRange(t *testing.T) {
	dt := DiagnosticType_Warning("FOO", "{0}, {1} - {2}")
	getRootNode := func(n *rhino.Node) *rhino.Node {
		for n.IsGetProp() {
			n = n.GetFirstChild()
		}
		return n
	}
	callback := &nodeTraversalTestCallback{
		shouldTraverse: func(nt *NodeTraversal, n, parent *rhino.Node) bool {
			if n.IsGetProp() && parent.IsExprResult() {
				nt.ReportWithNodeRange(getRootNode(n), n, dt, "Foo", "Bar", "Hello")
			}
			return true
		},
	}

	var errors []*JSError
	compiler := NewCompiler()
	compiler.SetErrorManager(newNodeTraversalTest_ErrorManager(&errors))
	compiler.InitCompilerOptionsIfTesting()
	code := "a.b.c;"
	tree := scopeTestParse(t, compiler, code)

	NodeTraversal_Builder_().SetCompiler(compiler).SetCallback(callback).Traverse(tree)

	if len(errors) != 1 {
		t.Fatalf("got %d errors, want 1", len(errors))
	}
	error_ := errors[0]
	assertEqual(t, "description", error_.Description(), "Foo, Bar - Hello")
	assertEqual(t, "node source offset", error_.GetNodeSourceOffset(), 0)
	assertEqual(t, "length", error_.Length(), 5)
}

const nodeTraversalTest_TEST_EXCEPTION = "test me"

func TestNodeTraversal_UnexpectedException(t *testing.T) {
	var cb NodeTraversal_AbstractPostOrderCallbackInterface = func(nt *NodeTraversal, n, parent *rhino.Node) {
		panic(jutil.NewRuntimeException(nodeTraversalTest_TEST_EXCEPTION))
	}

	compiler := NewCompiler()

	code := "function foo() {}"
	tree := scopeTestParse(t, compiler, code)
	thrown := scopeTestCatch(func() {
		NodeTraversal_Builder_().SetCompiler(compiler).SetCallbackFunc(cb).Traverse(tree)
	})
	e, ok := thrown.(jutil.IRuntimeException)
	if !ok {
		t.Fatalf("Expected RuntimeException, got %v", thrown)
	}
	if !strings.HasPrefix(e.GetMessage(), "INTERNAL COMPILER ERROR.\nPlease report this problem.\n\ntest me") {
		t.Errorf("message = %q", e.GetMessage())
	}
}

func TestNodeTraversal_GetScopeRoot(t *testing.T) {
	compiler := NewCompiler()
	code := `var a;
function foo() {
  var b
}
`
	tree := scopeTestParse(t, compiler, code)
	NodeTraversal_Traverse(compiler, tree, &nodeTraversalTestScopedCallback{
		enterScope: func(nt *NodeTraversal) {
			root1 := nt.GetScopeRoot()
			scope2 := nt.GetScope()
			root2 := scope2.GetRootNode()
			if root2 != root1 {
				t.Errorf("scope root %v != %v", root2, root1)
			}
		},
	})
}

func TestNodeTraversal_GetScopeRoot_inEsModule(t *testing.T) {
	compiler := NewCompiler()
	code := `const x = 0;
export {x};
`
	tree := scopeTestParse(t, compiler, code)
	check := func(nt *NodeTraversal, n, parent *rhino.Node) {
		if n.IsModuleBody() || (parent != nil && parent.IsModuleBody()) {
			if tok := nt.GetScopeRoot().GetToken(); tok != rhino.Token_MODULE_BODY {
				t.Errorf("scope root token %v, want MODULE_BODY", tok)
			}
		}
	}
	NodeTraversal_Traverse(compiler, tree, &nodeTraversalTestCallback{
		shouldTraverse: func(nt *NodeTraversal, n, parent *rhino.Node) bool {
			check(nt, n, parent)
			return true
		},
		visit: check,
	})
}

func TestNodeTraversal_GetScopeRoot_inGoogModule(t *testing.T) {
	compiler := NewCompiler()
	code := `goog.module('a.b');
function foo() {
  var b
}
`
	tree := scopeTestParse(t, compiler, code)
	check := func(nt *NodeTraversal, n, parent *rhino.Node) {
		if n.IsModuleBody() || (parent != nil && parent.IsModuleBody()) {
			if tok := nt.GetScopeRoot().GetToken(); tok != rhino.Token_MODULE_BODY {
				t.Errorf("scope root token %v, want MODULE_BODY", tok)
			}
		}
	}
	NodeTraversal_Traverse(compiler, tree, &nodeTraversalTestCallback{
		shouldTraverse: func(nt *NodeTraversal, n, parent *rhino.Node) bool {
			check(nt, n, parent)
			return true
		},
		visit: check,
	})
}

func TestNodeTraversal_GetHoistScopeRoot(t *testing.T) {
	compiler := NewCompiler()
	code := `function foo() {
  if (true) { var XXX; }
}
`
	tree := scopeTestParse(t, compiler, code)
	NodeTraversal_Traverse(compiler, tree, &nodeTraversalTestCallback{
		visit: func(nt *NodeTraversal, n, parent *rhino.Node) {
			if n.IsName() && n.GetString() == "XXX" {
				root := nt.GetClosestHoistScopeRoot()
				if !NodeUtil_IsFunctionBlock(root) {
					t.Errorf("closest hoist scope root %v is not a function block", root)
				}

				nt.GetScope() // force scope creation

				root = nt.GetClosestHoistScopeRoot()
				if !NodeUtil_IsFunctionBlock(root) {
					t.Errorf("closest hoist scope root %v is not a function block", root)
				}
			}
		},
	})
}

// nodeTraversalTest_NameChangingCallback is NodeTraversalTest.NameChangingCallback.
type nodeTraversalTest_NameChangingCallback struct{}

func (nodeTraversalTest_NameChangingCallback) ShouldTraverse(t *NodeTraversal, n, parent *rhino.Node) bool {
	return true
}

func (nodeTraversalTest_NameChangingCallback) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	if n.IsName() && n.GetString() == "change" {
		n.SetString("xx")
		t.ReportCodeChange()
	}
}

func TestNodeTraversal_ReportChange1(t *testing.T) {
	code := `var change;
function foo() {
  var b
}
`
	nodeTraversalTest_assertChangesRecorded(t, code, nodeTraversalTest_NameChangingCallback{})
}

func TestNodeTraversal_ReportChange2(t *testing.T) {
	code := `var a;
function foo() {
  var change
}
`
	nodeTraversalTest_assertChangesRecorded(t, code, nodeTraversalTest_NameChangingCallback{})
}

func TestNodeTraversal_ReportChange3(t *testing.T) {
	code := `var a;
function foo() {
  var b
}
var change
`
	nodeTraversalTest_assertChangesRecorded(t, code, nodeTraversalTest_NameChangingCallback{})
}

func TestNodeTraversal_ReportChange4(t *testing.T) {
	code := `function foo() {
  function bar() {
    var change
  }
}
`
	nodeTraversalTest_assertChangesRecorded(t, code, nodeTraversalTest_NameChangingCallback{})
}

// Java: NodeTraversalTest.assertChangesRecorded(String, NodeTraversal.Callback)
func nodeTraversalTest_assertChangesRecorded(t *testing.T, code string, callback NodeTraversal_Callback) {
	t.Helper()
	externs := ""
	compiler := NewCompiler()
	tree := scopeTestParseRoots(t, compiler, externs, code)

	changeVerifier := NewChangeVerifier(compiler).Snapshot(tree)
	NodeTraversal_TraverseRoots(compiler, callback, tree.GetFirstChild(), tree.GetSecondChild())
	// ChangeVerifier throws when a change was not recorded.
	if thrown := scopeTestCatch(func() { changeVerifier.CheckRecordedChanges(tree) }); thrown != nil {
		t.Fatal(thrown.Error())
	}
}

func TestNodeTraversal_GetCurrentNode(t *testing.T) {
	compiler := NewCompiler()
	creator := NewSyntacticScopeCreator(compiler)
	callback := newNodeTraversalTest_ExpectNodeOnEnterScope(t)
	b := NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(callback).
		SetScopeCreator(creator)

	code := `var a;
function foo() {
  var b;
}
`

	tree := scopeTestParse(t, compiler, code)
	topScope := creator.CreateScope(tree, nil).(*Scope)

	// Calling #traverseWithScope uses the given scope but starts traversal at
	// the given node.
	callback.expect(tree.GetFirstChild(), tree)
	b.TraverseWithScope(tree.GetFirstChild(), topScope)
	callback.assertEntered()

	// Calling #traverse creates a new scope with the given node as the root.
	callback.expect(tree.GetFirstChild(), tree.GetFirstChild())
	b.Traverse(tree.GetFirstChild())
	callback.assertEntered()

	// Calling #traverseAtScope starts traversal from the scope's root.
	fn := tree.GetSecondChild()
	fnScope := creator.CreateScope(fn, topScope).(*Scope)
	callback.expect(fn, fn)
	b.TraverseAtScope(fnScope)
	callback.assertEntered()
}

func TestNodeTraversal_TraverseAtScopeWithBlockScope(t *testing.T) {
	compiler := scopeTestNewCompiler(NewCompilerOptions())
	creator := NewSyntacticScopeCreator(compiler)
	callback := newNodeTraversalTest_ExpectNodeOnEnterScope(t)
	b := NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(callback).
		SetScopeCreator(creator)

	code := `function foo() {
  if (bar) {
    let x;
  }
}
`

	tree := scopeTestParse(t, compiler, code)
	topScope := creator.CreateScopeOfScope(tree, nil)

	innerBlock := tree. // script
				GetFirstChild(). // function
				GetLastChild().  // function body
				GetFirstChild(). // if
				GetLastChild()   // block

	blockScope := creator.CreateScopeOfScope(innerBlock, topScope)
	callback.expect(innerBlock, innerBlock)
	b.TraverseAtScope(blockScope)
	callback.assertEntered()
}

func TestNodeTraversal_TraverseAtScopeWithForScope(t *testing.T) {
	compiler := scopeTestNewCompiler(NewCompilerOptions())
	creator := NewSyntacticScopeCreator(compiler)
	callback := newNodeTraversalTest_ExpectNodeOnEnterScope(t)
	b := NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(callback).
		SetScopeCreator(creator)

	code := `function foo() {
  var b = [0];
  for (let a of b) {
    let x;
  }
}
`

	tree := scopeTestParse(t, compiler, code)
	topScope := creator.CreateScopeOfScope(tree, nil)

	forNode := tree. // script
				GetFirstChild(). // function
				GetLastChild().  // function body
				GetSecondChild() // for (first child is var b)

	innerBlock := forNode.GetLastChild()

	forScope := creator.CreateScopeOfScope(forNode, topScope)
	creator.CreateScopeOfScope(innerBlock, forScope)

	callback.expect(forNode, forNode)
	b.TraverseAtScope(forScope)
	callback.assertEntered()
}

func TestNodeTraversal_TraverseAtScopeWithSwitchScope(t *testing.T) {
	compiler := scopeTestNewCompiler(NewCompilerOptions())
	creator := NewSyntacticScopeCreator(compiler)
	callback := newNodeTraversalTest_ExpectNodeOnEnterScope(t)
	b := NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(callback).
		SetScopeCreator(creator)

	code := `function foo() {
  var b = [0];
  switch(b) {
    case 1:
       return b;
    case 2:
  }
}
`

	tree := scopeTestParse(t, compiler, code)
	topScope := creator.CreateScopeOfScope(tree, nil)

	innerBlock := tree. // script
				GetFirstChild().  // function
				GetLastChild().   // function body
				GetSecondChild(). // switch (first child is var b)
				GetSecondChild()  // switch block

	blockScope := creator.CreateScopeOfScope(innerBlock, topScope)
	callback.expect(innerBlock, innerBlock)
	b.TraverseAtScope(blockScope)
	callback.assertEntered()
}

func TestNodeTraversal_TraverseAtScopeWithModuleScope(t *testing.T) {
	compiler := scopeTestNewCompiler(NewCompilerOptions())
	creator := NewSyntacticScopeCreator(compiler)
	callback := newNodeTraversalTest_ExpectNodeOnEnterScope(t)
	b := NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(callback).
		SetScopeCreator(creator)

	code := `goog.module('example.module');

var x;
`

	tree := scopeTestParse(t, compiler, code)
	globalScope := creator.CreateScopeOfScope(tree, nil)
	moduleBody := tree.GetFirstChild()
	moduleScope := creator.CreateScopeOfScope(moduleBody, globalScope)

	callback.expect(moduleBody, moduleBody)

	b.TraverseAtScope(moduleScope)

	callback.assertEntered()
}

func TestNodeTraversal_TraverseAtScopeWithMemberFieldDefScope(t *testing.T) {
	compiler := NewCompiler()
	compiler.InitCompilerOptionsIfTesting()
	creator := NewSyntacticScopeCreator(compiler)
	callback := newNodeTraversalTest_ExpectNodeOnEnterScope(t)
	b := NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(callback).
		SetScopeCreator(creator)

	code := `class Foo {
  a = this.a;
}
class Bar extends Foo {
  b = super.a;
}
`

	tree := scopeTestParse(t, compiler, code)
	globalScope := creator.CreateScopeOfScope(tree, nil)
	memberFieldDefA := tree. // script
					GetFirstChild(). // class
					GetLastChild().  // class members
					GetFirstChild()  // member field def
	memberFieldDefAScope := creator.CreateScopeOfScope(memberFieldDefA, globalScope)

	callback.expect(memberFieldDefA, memberFieldDefA)
	b.TraverseAtScope(memberFieldDefAScope)
	callback.assertEntered()

	memberFieldDefB := tree. // script
					GetSecondChild(). // class
					GetLastChild().   // class members
					GetFirstChild()   // member field def
	memberFieldDefBScope := creator.CreateScopeOfScope(memberFieldDefB, globalScope)

	callback.expect(memberFieldDefB, memberFieldDefB)
	b.TraverseAtScope(memberFieldDefBScope)
	callback.assertEntered()
}

func TestNodeTraversal_TraverseAtScopeWithComputedFieldDefScope(t *testing.T) {
	compiler := NewCompiler()
	compiler.InitCompilerOptionsIfTesting()
	creator := NewSyntacticScopeCreator(compiler)
	callback := newNodeTraversalTest_ExpectNodeOnEnterScope(t)
	b := NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(callback).
		SetScopeCreator(creator)

	code := `class Foo {
  x = 'hi';
  [this.x] = this.x;
}
`

	tree := scopeTestParse(t, compiler, code)
	globalScope := creator.CreateScopeOfScope(tree, nil)
	computedFieldDef := tree. // script
					GetFirstChild(). // class
					GetLastChild().  // class members
					GetLastChild()   // computed field def
	computedFieldDefScope := creator.CreateScopeOfScope(computedFieldDef, globalScope)

	callback.expect(computedFieldDef, computedFieldDef)
	b.TraverseAtScope(computedFieldDefScope)
	callback.assertEntered()
}

func TestNodeTraversal_TraverseFieldDefScopeRootsInOrder(t *testing.T) {
	compiler := NewCompiler()
	code := `class Foo {
  x = 'hi';
  [this.x] = this.x;
}
`
	tree := scopeTestParse(t, compiler, code)
	callback := &nodeTraversalTest_TokenAccumulator{}
	NodeTraversal_Traverse(compiler, tree, callback)
	assertSliceEqual(t, "scope roots", callback.scopeRoots, []rhino.Token{
		rhino.Token_SCRIPT, rhino.Token_CLASS, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF})
}

func TestNodeTraversal_TraverseComputedFieldsInOrder(t *testing.T) {
	compiler := NewCompiler()
	compiler.InitCompilerOptionsIfTesting()
	callback := &nodeTraversalTest_TokenAccumulator{}
	creator := NewSyntacticScopeCreator(compiler)
	b := NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(callback).
		SetScopeCreator(creator)

	code := `class Foo {
  [this.x] = true;
}
`

	tree := scopeTestParse(t, compiler, code)
	globalScope := creator.CreateScopeOfScope(tree, nil)
	computedFieldDef := tree. // script
					GetFirstChild(). // class
					GetLastChild().  // class members
					GetLastChild()   // computed field def
	computedFieldDefScope := creator.CreateScopeOfScope(computedFieldDef, globalScope)

	b.TraverseAtScope(computedFieldDefScope)
	assertSliceEqual(t, "tokens", callback.tokens, []rhino.Token{rhino.Token_TRUE, rhino.Token_COMPUTED_FIELD_DEF})

	callback.tokens = nil
}

func TestNodeTraversal_GetVarAccessible(t *testing.T) {
	compiler := scopeTestNewCompiler(NewCompilerOptions())
	creator := NewSyntacticScopeCreator(compiler)
	callback := &nodeTraversalTest_AccessibleCallback{t: t}
	b := NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(callback).
		SetScopeCreator(creator)

	// variables are hoisted to their enclosing scope
	code := `var varDefinedInScript;
var foo = function(param) {
  var varDefinedInFoo;
  var baz = function() {
    var varDefinedInBaz;
  }
}
var bar = function() {
  var varDefinedInBar;
}
`

	// the function scope should have access to all variables defined before and in the function
	// scope
	tree := scopeTestParse(t, compiler, code)
	fooNode := tree. // script
				GetSecondChild().    // var foo declaration (first child is var varDefinedInScript)
				GetFirstFirstChild() // child of the var foo declaration is the foo function
	topScope := creator.CreateScopeOfScope(tree, nil)
	fooScope := creator.CreateScopeOfScope(fooNode, topScope)
	callback.expect(4)
	b.TraverseAtScope(fooScope)
	callback.assertAccessible(fooScope)

	// the function block scope should have access to all variables defined in the global, function,
	// and function block scopes
	fooBlockNode := fooNode.GetLastChild()
	fooBlockScope := creator.CreateScopeOfScope(fooBlockNode, fooScope)
	callback.expect(6)
	b.TraverseAtScope(fooBlockScope)
	callback.assertAccessible(fooBlockScope)

	// let and const variables are block scoped
	code = `var foo = function() {
  var varDefinedInFoo;
  var baz = function() {
    var varDefinedInBaz;
    let varDefinedInFoo; // shadows parent scope
  }
  let bar = 1;
}
`

	// the baz block scope has access to variables in its scope and parent scopes
	tree = scopeTestParse(t, compiler, code)
	fooNode = tree. // script
			GetFirstChild().     // var foo declaration (first child is var varDefinedInScript)
			GetFirstFirstChild() // child of the var foo declaration is the foo function
	fooBlockNode = fooNode.GetLastChild() // first child is param list of foo
	bazNode := fooBlockNode.GetSecondChild().GetFirstFirstChild()
	bazBlockNode := bazNode.GetLastChild()

	topScope = creator.CreateScopeOfScope(tree, nil)
	fooScope = creator.CreateScopeOfScope(fooNode, topScope)
	fooBlockScope = creator.CreateScopeOfScope(fooBlockNode, fooScope)
	bazScope := creator.CreateScopeOfScope(bazNode, fooBlockScope)
	bazBlockScope := creator.CreateScopeOfScope(bazBlockNode, bazScope)

	// bar, baz, foo, varDefinedInFoo(in baz function), varDefinedInBaz
	callback.expect(5)
	b.TraverseAtScope(bazBlockScope)
	callback.assertAccessible(bazBlockScope)
}

func TestNodeTraversal_TraverseEs6ScopeRoots_isLimitedToScope(t *testing.T) {
	compiler := NewCompiler()
	callback := &nodeTraversalTest_StringAccumulator{}

	code := `function foo() {
  'string in foo';
  function baz() {
    'string nested in baz';
  }
}
function bar() {
  'string in bar';
}
`

	tree := scopeTestParse(t, compiler, code)
	fooFunction := tree.GetFirstChild()

	// Traverse without entering nested scopes.
	NodeTraversal_TraverseScopeRoots(compiler, []*rhino.Node{fooFunction}, callback, false)
	assertSameElements(t, "strings", callback.strings, []string{"string in foo"})

	callback.strings = nil

	// Traverse *with* entering nested scopes, now also sees "string nested in baz".
	NodeTraversal_TraverseScopeRoots(compiler, []*rhino.Node{fooFunction}, callback, true)
	assertSameElements(t, "strings", callback.strings, []string{"string in foo", "string nested in baz"})
}

func TestNodeTraversal_TraverseEs6ScopeRoots_parentScopesWork(t *testing.T) {
	compiler := NewCompiler()
	callback := &nodeTraversalTest_LexicallyScopedVarsAccumulator{varNames: jutil.NewLinkedHashSet[string]()}

	code := `var varDefinedInScript;
var foo = function() {
  var varDefinedInFoo;
  var baz = function() {
    var varDefinedInBaz;
  }
}
var bar = function() {
  var varDefinedInBar;
}
`

	tree := scopeTestParse(t, compiler, code)
	fooFunction := tree.GetSecondChild().GetFirstFirstChild()

	// Traverse without entering nested scopes.
	NodeTraversal_TraverseScopeRoots(compiler, []*rhino.Node{fooFunction}, callback, false)
	assertSameElements(t, "var names", callback.varNames.ToSlice(),
		[]string{"varDefinedInScript", "foo", "bar", "varDefinedInFoo", "baz"})

	callback.varNames.Clear()

	// Traverse *with* entering nested scopes, now also sees "varDefinedInBaz".
	NodeTraversal_TraverseScopeRoots(compiler, []*rhino.Node{fooFunction}, callback, true)
	assertSameElements(t, "var names", callback.varNames.ToSlice(),
		[]string{"varDefinedInScript", "foo", "bar", "varDefinedInFoo", "baz", "varDefinedInBaz"})
}

func TestNodeTraversal_TraverseEs6ScopeRoots_callsEnterScope(t *testing.T) {
	compiler := NewCompiler()

	var scopesEntered []*rhino.Node

	testCallback := &nodeTraversalTestScopedCallback{
		enterScope: func(nt *NodeTraversal) {
			scopesEntered = append(scopesEntered, nt.GetScopeRoot())
		},
	}

	code := "function foo() { {} }"

	tree := scopeTestParse(t, compiler, code)
	fooFunction := tree.GetFirstChild()

	NodeTraversal_TraverseScopeRoots(compiler, []*rhino.Node{fooFunction}, testCallback, true)
	// Function, function's body, and the block inside it.
	assertEqual(t, "scopes entered", len(scopesEntered), 3)
}

func TestNodeTraversal_NodeTraversalInterruptable(t *testing.T) {
	t.Skip("needs Thread.interrupt(): goroutines cannot be interrupted (Platform.isThreadInterrupted is always false)")
}

// nodeTraversalTest_LexicallyScopedVarsAccumulator is the helper class used
// to collect all the vars from current scope and its parent scopes.
type nodeTraversalTest_LexicallyScopedVarsAccumulator struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	varNames *jutil.LinkedHashSet[string]
}

func (a *nodeTraversalTest_LexicallyScopedVarsAccumulator) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	firstScope := t.GetScope()
	if firstScope == nil {
		return
	}

	for scope := firstScope; scope != nil; scope = scope.GetParent() {
		for var_ := range scope.GetVarIterable() {
			a.varNames.Add(var_.GetName())
		}
	}
}

// nodeTraversalTest_StringAccumulator is NodeTraversalTest.StringAccumulator.
type nodeTraversalTest_StringAccumulator struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	strings []string
}

func (a *nodeTraversalTest_StringAccumulator) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	if n.IsStringLit() {
		a.strings = append(a.strings, n.GetString())
	}
}

// nodeTraversalTest_TokenAccumulator is NodeTraversalTest.TokenAccumulator.
type nodeTraversalTest_TokenAccumulator struct {
	NodeTraversal_AbstractScopedCallbackBase
	tokens     []rhino.Token
	scopeRoots []rhino.Token
}

func (a *nodeTraversalTest_TokenAccumulator) Visit(t *NodeTraversal, n, parent *rhino.Node) {
	a.tokens = append(a.tokens, n.GetToken())
}

func (a *nodeTraversalTest_TokenAccumulator) EnterScope(t *NodeTraversal) {
	a.scopeRoots = append(a.scopeRoots, t.GetScopeRoot().GetToken())
}

// nodeTraversalTest_ExpectNodeOnEnterScope is the helper class used to test
// getCurrentNode.
type nodeTraversalTest_ExpectNodeOnEnterScope struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	t         *testing.T
	node      *rhino.Node
	scopeRoot *rhino.Node
	entered   bool
}

func newNodeTraversalTest_ExpectNodeOnEnterScope(t *testing.T) *nodeTraversalTest_ExpectNodeOnEnterScope {
	return &nodeTraversalTest_ExpectNodeOnEnterScope{t: t}
}

func (c *nodeTraversalTest_ExpectNodeOnEnterScope) expect(node, scopeRoot *rhino.Node) {
	c.node = node
	c.scopeRoot = scopeRoot
	c.entered = false
}

func (c *nodeTraversalTest_ExpectNodeOnEnterScope) assertEntered() {
	c.t.Helper()
	if !c.entered {
		c.t.Error("expected the scope to be entered")
	}
}

func (c *nodeTraversalTest_ExpectNodeOnEnterScope) EnterScope(t *NodeTraversal) {
	if t.GetCurrentNode() != c.node {
		c.t.Errorf("current node %v, want %v", t.GetCurrentNode(), c.node)
	}
	if t.GetScopeRoot() != c.scopeRoot {
		c.t.Errorf("scope root %v, want %v", t.GetScopeRoot(), c.scopeRoot)
	}
	if c.node.IsForIn() || c.node.IsForOf() {
		c.node = c.node.GetLastChild()
		c.scopeRoot = c.scopeRoot.GetLastChild()
	}
	c.entered = true
}

func (c *nodeTraversalTest_ExpectNodeOnEnterScope) ExitScope(t *NodeTraversal) {}

func (c *nodeTraversalTest_ExpectNodeOnEnterScope) ShouldTraverse(t *NodeTraversal, n, parent *rhino.Node) bool {
	return !c.entered
}

func TestNodeTraversal_TraverseComputedFieldsInClass(t *testing.T) {
	compiler := NewCompiler()
	compiler.InitCompilerOptionsIfTesting()
	callback := &nodeTraversalTest_StringAccumulator{}

	code := `class Foo {
  ['in field lhs'] = 'in field rhs';
  ['in method lhs']() {
    'nested in method';
  }
}
`

	tree := scopeTestParse(t, compiler, code)

	NodeTraversal_Traverse(compiler, tree, callback)
	assertSliceEqual(t, "strings", callback.strings,
		[]string{"in field lhs", "in method lhs", "in field rhs", "nested in method"})

	callback.strings = nil
}

// nodeTraversalTest_AccessibleCallback is the helper class used to test
// accessible variables.
type nodeTraversalTest_AccessibleCallback struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	t             *testing.T
	numAccessible int
}

func (c *nodeTraversalTest_AccessibleCallback) expect(accessible int) {
	c.numAccessible = accessible
}

func (c *nodeTraversalTest_AccessibleCallback) assertAccessible(s *Scope) {
	c.t.Helper()
	if got := len(slices.Collect(s.GetAllAccessibleVariables())); got != c.numAccessible {
		c.t.Errorf("%d accessible variables, want %d", got, c.numAccessible)
	}
}

func (c *nodeTraversalTest_AccessibleCallback) EnterScope(t *NodeTraversal) {}

func (c *nodeTraversalTest_AccessibleCallback) ExitScope(t *NodeTraversal) {}

func (c *nodeTraversalTest_AccessibleCallback) ShouldTraverse(t *NodeTraversal, n, parent *rhino.Node) bool {
	return true
}

// Port of test/com/google/javascript/jscomp/SyntacticScopeCreatorTest.java @ closure-compiler v20260415

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// syntacticScopeCreatorTest holds the fields of SyntacticScopeCreatorTest;
// newSyntacticScopeCreatorTest is its setUp.
type syntacticScopeCreatorTest struct {
	t              *testing.T
	compiler       *Compiler
	scopeCreator   *SyntacticScopeCreator
	redeclarations map[string]int // the HashMultiset<String>
}

// syntacticScopeCreatorTest_RecordingRedeclarationHandler is
// SyntacticScopeCreatorTest.RecordingRedeclarationHandler.
type syntacticScopeCreatorTest_RecordingRedeclarationHandler struct {
	s *syntacticScopeCreatorTest
}

func (h *syntacticScopeCreatorTest_RecordingRedeclarationHandler) OnRedeclaration(
	s *Scope, name string, n *rhino.Node, input *CompilerInput) {
	h.s.redeclarations[name]++
}

func newSyntacticScopeCreatorTest_RecordingRedeclarationHandler(
	s *syntacticScopeCreatorTest) *syntacticScopeCreatorTest_RecordingRedeclarationHandler {
	return &syntacticScopeCreatorTest_RecordingRedeclarationHandler{s: s}
}

// Java: SyntacticScopeCreatorTest.getRoot(String)
func (s *syntacticScopeCreatorTest) getRoot(js string) *rhino.Node {
	s.t.Helper()
	root := s.compiler.ParseTestCode(js)
	if errors := s.compiler.GetErrors(); len(errors) != 0 {
		s.t.Fatalf("parse errors: %v", errors)
	}
	return root
}

// getScope is a helper to create a top-level scope from a JavaScript
// string.
func (s *syntacticScopeCreatorTest) getScope(js string) *Scope {
	return s.scopeCreator.CreateScopeOfScope(s.getRoot(js), nil)
}

// newSyntacticScopeCreatorTest is SyntacticScopeCreatorTest.setUp.
func newSyntacticScopeCreatorTest(t *testing.T) *syntacticScopeCreatorTest {
	s := &syntacticScopeCreatorTest{t: t}
	options := NewCompilerOptions()
	s.compiler = scopeTestNewCompiler(options)
	s.redeclarations = map[string]int{}
	handler := newSyntacticScopeCreatorTest_RecordingRedeclarationHandler(s)
	s.scopeCreator = NewSyntacticScopeCreatorWithRedeclarationHandler(s.compiler, handler)
	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	return s
}

// assertRedeclarations is assertThat(redeclarations).hasCount(name, count).
func (s *syntacticScopeCreatorTest) assertRedeclarations(name string, count int) {
	s.t.Helper()
	if got := s.redeclarations[name]; got != count {
		s.t.Errorf("redeclarations of %s: %d, want %d", name, got, count)
	}
}

// assertNoRedeclarations is assertThat(redeclarations).isEmpty().
func (s *syntacticScopeCreatorTest) assertNoRedeclarations() {
	s.t.Helper()
	s.assertRedeclarationsSize(0)
}

// assertRedeclarationsSize is assertThat(redeclarations).hasSize(size).
func (s *syntacticScopeCreatorTest) assertRedeclarationsSize(size int) {
	s.t.Helper()
	n := 0
	for _, c := range s.redeclarations {
		n += c
	}
	if n != size {
		s.t.Errorf("%d redeclarations (%v), want %d", n, s.redeclarations, size)
	}
}

// scopeTestAssertSame is assertThat(got).isSameInstanceAs(want) (and
// isEqualTo for identity-compared values).
func scopeTestAssertSame[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want the same instance as %v", got, want)
	}
}

// scopeTestAssertNotSame is assertThat(got).isNotSameInstanceAs(want).
func scopeTestAssertNotSame[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("got %v, want another instance", got)
	}
}

func TestSyntacticScopeCreator_VarRedeclaration1(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	s.getScope("var x; var x")
	s.assertRedeclarations("x", 1)
}

func TestSyntacticScopeCreator_VarRedeclaration2(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	s.getScope("var x; var x; var x;")
	s.assertRedeclarations("x", 2)
}

func TestSyntacticScopeCreator_VarRedeclaration3(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "var x; if (true) { var x; } var x;"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	block :=
		root.GetFirstChild(). // VAR
					GetNext().     // IF
					GetLastChild() // BLOCK
	jutil.CheckStateMsg1(block.IsBlock(), "%s", block)
	s.scopeCreator.CreateScopeOfScope(block, globalScope)

	s.assertRedeclarations("x", 2)
}

func TestSyntacticScopeCreator_VarRedeclaration4(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "var x; if (true) { var x; var x; }"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	block :=
		root.GetFirstChild(). // VAR
					GetNext().     // IF
					GetLastChild() // BLOCK
	jutil.CheckStateMsg1(block.IsBlock(), "%s", block)
	s.scopeCreator.CreateScopeOfScope(block, globalScope)

	s.assertRedeclarations("x", 2)
}

func TestSyntacticScopeCreator_VarRedeclaration_withDestructuring(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function foo() { var x; var [x] = 1; }"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	block := root.GetFirstChild().GetLastChild() // BLOCK
	jutil.CheckStateMsg1(block.IsBlock(), "%s", block)
	s.scopeCreator.CreateScopeOfScope(block, globalScope)

	s.assertRedeclarations("x", 1)
}

func TestSyntacticScopeCreator_VarRedeclaration5(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "if (true) { var x; var x; }"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	block :=
		root.GetFirstChild(). // IF
					GetLastChild() // BLOCK
	jutil.CheckStateMsg1(block.IsBlock(), "%s", block)
	s.scopeCreator.CreateScopeOfScope(block, globalScope)

	s.assertRedeclarations("x", 1)
}

func TestSyntacticScopeCreator_VarShadowsParam(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function f(p) { var p; }"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	function := root.GetFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(function, globalScope)

	body := function.GetLastChild()
	bodyScope := s.scopeCreator.CreateScopeOfScope(body, functionScope)

	assertSameElements(t, "globalScope vars", scopeTestVarNames(globalScope.GetVarIterable()), []string{"f"})
	assertSameElements(t, "functionScope vars", scopeTestVarNames(functionScope.GetVarIterable()), []string{"p"})

	// "var p" doesn't declare a new var, so there is no 'p' variable in the function body scope.
	assertEqual(t, "bodyScope var count", len(scopeTestVarNames(bodyScope.GetVarIterable())), 0)
}

func TestSyntacticScopeCreator_ParamShadowsFunctionName(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "var f = function g(g) { }"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	function := root.GetFirstChild().GetFirstFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(function, globalScope)

	body := function.GetLastChild()
	bodyScope := s.scopeCreator.CreateScopeOfScope(body, functionScope)

	assertSameElements(t, "globalScope vars", scopeTestVarNames(globalScope.GetVarIterable()), []string{"f"})
	assertSameElements(t, "functionScope vars", scopeTestVarNames(functionScope.GetVarIterable()), []string{"g"})
	assertEqual(t, "bodyScope var count", len(scopeTestVarNames(bodyScope.GetVarIterable())), 0)
}

func TestSyntacticScopeCreator_VarShadowsFunctionName(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "var f = function g() { var g; }"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	function := root.GetFirstChild().GetFirstFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(function, globalScope)

	body := function.GetLastChild()
	bodyScope := s.scopeCreator.CreateScopeOfScope(body, functionScope)

	assertSameElements(t, "globalScope vars", scopeTestVarNames(globalScope.GetVarIterable()), []string{"f"})
	assertSameElements(t, "functionScope vars", scopeTestVarNames(functionScope.GetVarIterable()), []string{"g"})

	// "var g" declares a new variable, which shadows the function name.
	assertSameElements(t, "bodyScope vars", scopeTestVarNames(bodyScope.GetVarIterable()), []string{"g"})
}

func TestSyntacticScopeCreator_ParamAndVarShadowFunctionName(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "var f = function g(g) { var g; }"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	function := root.GetFirstChild().GetFirstFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(function, globalScope)

	body := function.GetLastChild()
	bodyScope := s.scopeCreator.CreateScopeOfScope(body, functionScope)

	assertSameElements(t, "globalScope vars", scopeTestVarNames(globalScope.GetVarIterable()), []string{"f"})
	assertSameElements(t, "functionScope vars", scopeTestVarNames(functionScope.GetVarIterable()), []string{"g"})

	// "var g" doesn't declare a new var, so there is no 'g' variable in the function body scope.
	assertEqual(t, "bodyScope var count", len(scopeTestVarNames(bodyScope.GetVarIterable())), 0)
}

func TestSyntacticScopeCreator_VarRedeclaration1_inES6Module(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "export function f() { var x; var x; }"

	script := s.getRoot(js)
	global := s.scopeCreator.CreateScopeOfScope(script, nil)

	moduleBody := script.GetFirstChild()
	jutil.CheckState(moduleBody.IsModuleBody())
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, global)

	function := moduleBody.GetFirstFirstChild()
	jutil.CheckState(function.IsFunction())
	functionScope := s.scopeCreator.CreateScopeOfScope(function, moduleScope)

	functionBody := function.GetLastChild()
	s.scopeCreator.CreateScopeOfScope(functionBody, functionScope)

	s.assertRedeclarations("x", 1)
}

func TestSyntacticScopeCreator_VarRedeclaration2_inES6Module(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "export var x = 1; export var x = 2;"

	script := s.getRoot(js)
	global := s.scopeCreator.CreateScopeOfScope(script, nil)

	moduleBody := script.GetFirstChild()
	jutil.CheckState(moduleBody.IsModuleBody())
	s.scopeCreator.CreateScopeOfScope(moduleBody, global)

	s.assertRedeclarations("x", 1)
}

func TestSyntacticScopeCreator_Redeclaration3_inES6Module(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "export function f() { var x; if (true) { var x; var x; } var x; }"

	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	moduleBody := root.GetFirstChild()
	jutil.CheckState(moduleBody.IsModuleBody())
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, globalScope)

	function := moduleBody.GetFirstFirstChild()
	jutil.CheckState(function.IsFunction())
	functionScope := s.scopeCreator.CreateScopeOfScope(function, moduleScope)

	functionBody := function.GetLastChild()
	functionBlockScope := s.scopeCreator.CreateScopeOfScope(functionBody, functionScope)

	innerBlock :=
		functionBody.
			GetFirstChild(). // VAR
			GetNext().       // IF
			GetLastChild()   // BLOCK
	jutil.CheckStateMsg1(innerBlock.IsBlock(), "%s", innerBlock)
	s.scopeCreator.CreateScopeOfScope(innerBlock, functionBlockScope)

	s.assertRedeclarations("x", 3)
}

func TestSyntacticScopeCreator_LetRedeclaration1(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	s.getScope("let x; let x")
	s.assertRedeclarations("x", 1)
}

func TestSyntacticScopeCreator_LetRedeclaration2(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	s.getScope("let x; let x; let x;")
	s.assertRedeclarations("x", 2)
}

func TestSyntacticScopeCreator_LetRedeclaration3(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "let x; if (true) { let x; } let x;"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	block :=
		root.GetFirstChild(). // VAR
					GetNext().     // IF
					GetLastChild() // BLOCK
	jutil.CheckStateMsg1(block.IsBlock(), "%s", block)
	s.scopeCreator.CreateScopeOfScope(block, globalScope)

	s.assertRedeclarations("x", 1)
}

func TestSyntacticScopeCreator_LetRedeclaration3_withES6Module(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "export function f() { let x; if (true) { let x; } let x; }"

	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	moduleBody := root.GetFirstChild()
	jutil.CheckState(moduleBody.IsModuleBody())
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, globalScope)

	function := moduleBody.GetFirstFirstChild()
	jutil.CheckState(function.IsFunction())
	functionScope := s.scopeCreator.CreateScopeOfScope(function, moduleScope)

	functionBody := function.GetLastChild()
	functionBlockScope := s.scopeCreator.CreateScopeOfScope(functionBody, functionScope)

	innerBlock :=
		functionBody.
			GetFirstChild(). // VAR
			GetNext().       // IF
			GetLastChild()   // BLOCK
	s.scopeCreator.CreateScopeOfScope(innerBlock, functionBlockScope)

	s.assertRedeclarations("x", 1)
}

func TestSyntacticScopeCreator_LetRedeclaration4(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "let x; if (true) { let x; let x; }"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	block :=
		root.GetFirstChild(). // VAR
					GetNext().     // IF
					GetLastChild() // BLOCK
	jutil.CheckStateMsg1(block.IsBlock(), "%s", block)
	s.scopeCreator.CreateScopeOfScope(block, globalScope)

	s.assertRedeclarations("x", 1)
}

func TestSyntacticScopeCreator_LetRedeclaration5(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "if (true) { let x; let x; }"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	block :=
		root.GetFirstChild(). // IF
					GetLastChild() // BLOCK
	jutil.CheckStateMsg1(block.IsBlock(), "%s", block)
	s.scopeCreator.CreateScopeOfScope(block, globalScope)

	s.assertRedeclarations("x", 1)
}

func TestSyntacticScopeCreator_ArrayDestructuring(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	scope := s.getScope("var [x, y] = foo();")
	assertScope(t, scope).declares("x").directly()
	assertScope(t, scope).declares("y").directly()
}

func TestSyntacticScopeCreator_NestedArrayDestructuring(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	scope := s.getScope("var [x, [y,z]] = foo();")
	assertScope(t, scope).declares("x").directly()
	assertScope(t, scope).declares("y").directly()
	assertScope(t, scope).declares("z").directly()
}

func TestSyntacticScopeCreator_ArrayDestructuringWithName(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	scope := s.getScope("var a = 1, [x, y] = foo();")
	assertScope(t, scope).declares("a").directly()
	assertScope(t, scope).declares("x").directly()
	assertScope(t, scope).declares("y").directly()
}

func TestSyntacticScopeCreator_ArrayDestructuringLet(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`function foo() {
  var [a, b] = getVars()
  if (true) {
    let [x, y] = getLets()
  }
}
`
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	functionNode := root.GetFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(functionNode, globalScope)

	functionBlock := functionNode.GetLastChild()
	functionBlockScope := s.scopeCreator.CreateScopeOfScope(functionBlock, functionScope)

	assertScope(t, functionBlockScope).declares("a").directly()
	assertScope(t, functionBlockScope).declares("b").directly()
	assertScope(t, functionBlockScope).doesNotDeclare("x")
	assertScope(t, functionBlockScope).doesNotDeclare("y")

	var_ := functionBlock.GetFirstChild()
	ifStmt := var_.GetNext()
	ifBlock := ifStmt.GetLastChild()
	blockScope := s.scopeCreator.CreateScopeOfScope(ifBlock, functionBlockScope)

	// a and b are declared in the parent scope.
	assertScope(t, blockScope).declares("a").onSomeParent()
	assertScope(t, blockScope).declares("b").onSomeParent()

	// x and y are declared in this scope.
	assertScope(t, blockScope).declares("x").directly()
	assertScope(t, blockScope).declares("y").directly()
}

func TestSyntacticScopeCreator_ArrayDestructuringVarInBlock(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`function foo() {
  var [a, b] = getVars()
  if (true) {
    var [x, y] = getMoreVars()
  }
}
`
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	functionNode := root.GetFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(functionNode, globalScope)

	functionBlock := functionNode.GetLastChild()
	functionBlockScope := s.scopeCreator.CreateScopeOfScope(functionBlock, functionScope)

	assertScope(t, functionBlockScope).declares("a").directly()
	assertScope(t, functionBlockScope).declares("b").directly()
	assertScope(t, functionBlockScope).declares("x").directly()
	assertScope(t, functionBlockScope).declares("y").directly()
}

func TestSyntacticScopeCreator_ObjectDestructuring(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`function foo() {
  var {a, b} = bar()
}
`
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	functionNode := root.GetFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(functionNode, globalScope)

	functionBlock := functionNode.GetLastChild()
	functionBlockScope := s.scopeCreator.CreateScopeOfScope(functionBlock, functionScope)

	assertScope(t, functionBlockScope).declares("a").directly()
	assertScope(t, functionBlockScope).declares("b").directly()
}

func TestSyntacticScopeCreator_ObjectDestructuring2(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`function foo() {
  var {a: b = 1} = bar()
}
`
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	functionNode := root.GetFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(functionNode, globalScope)

	functionBlock := functionNode.GetLastChild()
	functionBlockScope := s.scopeCreator.CreateScopeOfScope(functionBlock, functionScope)

	assertScope(t, functionBlockScope).doesNotDeclare("a")
	assertScope(t, functionBlockScope).declares("b").directly()
}

func TestSyntacticScopeCreator_ObjectDestructuringComputedProp(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`function foo() {
  var {['s']: a} = bar()
}
`
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	functionNode := root.GetFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(functionNode, globalScope)

	functionBlock := functionNode.GetLastChild()
	functionBlockScope := s.scopeCreator.CreateScopeOfScope(functionBlock, functionScope)

	assertScope(t, functionBlockScope).declares("a").directly()
}

func TestSyntacticScopeCreator_ObjectDestructuringComputedPropParam(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function foo({['s']: a}) {}"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	functionNode := root.GetFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(functionNode, globalScope)
	assertScope(t, functionScope).declares("a").directly()
}

func TestSyntacticScopeCreator_ObjectDestructuringNested(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`function foo() {
  var {a:{b}} = bar()
}
`
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	functionNode := root.GetFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(functionNode, globalScope)

	functionBlock := functionNode.GetLastChild()
	functionBlockScope := s.scopeCreator.CreateScopeOfScope(functionBlock, functionScope)

	assertScope(t, functionBlockScope).doesNotDeclare("a")
	assertScope(t, functionBlockScope).declares("b").directly()
}

func TestSyntacticScopeCreator_ObjectDestructuringWithInitializer(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`function foo() {
  var {a=1} = bar()
}
`
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	functionNode := root.GetFirstChild()
	functionScope := s.scopeCreator.CreateScopeOfScope(functionNode, globalScope)

	functionBlock := functionNode.GetLastChild()
	functionBlockScope := s.scopeCreator.CreateScopeOfScope(functionBlock, functionScope)

	assertScope(t, functionBlockScope).declares("a").directly()
}

func TestSyntacticScopeCreator_ObjectDestructuringInForOfParam(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "{for (let {length: x} of gen()) {}}"
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	block := root.GetFirstChild()
	blockScope := s.scopeCreator.CreateScopeOfScope(block, globalScope)
	forOf := block.GetFirstChild()
	forOfScope := s.scopeCreator.CreateScopeOfScope(forOf, blockScope)

	assertScope(t, forOfScope).declares("x").directly()
}

func TestSyntacticScopeCreator_FunctionScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	scope :=
		s.getScope(
			`function foo() {}
var x = function bar(a1) {}
[function bar2() { var y; }]
if (true) { function z() {} }
`)
	assertScope(t, scope).declares("foo").directly()
	assertScope(t, scope).declares("x").directly()
	assertScope(t, scope).doesNotDeclare("z")

	// The following should not be declared in this scope
	assertScope(t, scope).doesNotDeclare("a1")
	assertScope(t, scope).doesNotDeclare("bar")
	assertScope(t, scope).doesNotDeclare("bar2")
	assertScope(t, scope).doesNotDeclare("y")
	assertScope(t, scope).doesNotDeclare("")
}

func TestSyntacticScopeCreator_ClassScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	scope :=
		s.getScope(
			`class Foo {}
var x = class Bar {}
[class Bar2 { constructor(a1) {} static y() {} }]
if (true) { class Z {} }
`)
	assertScope(t, scope).declares("Foo").directly()
	assertScope(t, scope).declares("x").directly()
	assertScope(t, scope).doesNotDeclare("Z")

	// The following should not be declared in this scope
	assertScope(t, scope).doesNotDeclare("a1")
	assertScope(t, scope).doesNotDeclare("Bar")
	assertScope(t, scope).doesNotDeclare("Bar2")
	assertScope(t, scope).doesNotDeclare("y")
	assertScope(t, scope).doesNotDeclare("")
}

func TestSyntacticScopeCreator_ScopeRootNode(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`function foo() {
 var x = 10
}
`
	root := s.getRoot(js)

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	scopeTestAssertSame(t, globalScope.GetRootNode(), root)
	assertEqual(t, "condition", globalScope.IsBlockScope(), false)
	scopeTestAssertSame(t, globalScope.GetClosestHoistScope(), globalScope)
	assertEqual(t, "condition", globalScope.IsHoistScope(), true)

	function := root.GetFirstChild()
	jutil.CheckStateMsg1(function.IsFunction(), "%s", function)
	functionScope := s.scopeCreator.CreateScopeOfScope(function, globalScope)

	fooBlockNode := NodeUtil_GetFunctionBody(function)
	fooScope := s.scopeCreator.CreateScopeOfScope(fooBlockNode, functionScope)
	scopeTestAssertSame(t, fooScope.GetRootNode(), fooBlockNode)
	assertEqual(t, "condition", fooScope.IsBlockScope(), true)
	scopeTestAssertSame(t, fooScope.GetClosestHoistScope(), fooScope)
	assertEqual(t, "condition", fooScope.IsHoistScope(), true)
	assertScope(t, fooScope).declares("x").directly()
}

func TestSyntacticScopeCreator_BlockScopeWithVar(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "if (true) { if (true) { var x; } }"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).declares("x").directly()

	firstLevelBlock := root.GetFirstChild().GetLastChild()
	firstLevelBlockScope := s.scopeCreator.CreateScopeOfScope(firstLevelBlock, globalScope)
	assertScope(t, firstLevelBlockScope).declares("x").onSomeParent()

	secondLevelBlock := firstLevelBlock.GetFirstChild().GetLastChild()
	secondLevelBLockScope := s.scopeCreator.CreateScopeOfScope(secondLevelBlock, firstLevelBlockScope)
	assertScope(t, secondLevelBLockScope).declares("x").onSomeParent()
}

func TestSyntacticScopeCreator_BlockScopeWithLet(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "if (true) { if (true) { let x; } }"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("x")

	firstLevelBlock := root.GetFirstChild().GetLastChild()
	firstLevelBlockScope := s.scopeCreator.CreateScopeOfScope(firstLevelBlock, globalScope)
	assertScope(t, firstLevelBlockScope).doesNotDeclare("x")

	secondLevelBlock := firstLevelBlock.GetFirstChild().GetLastChild()
	secondLevelBLockScope := s.scopeCreator.CreateScopeOfScope(secondLevelBlock, firstLevelBlockScope)
	assertScope(t, secondLevelBLockScope).declares("x").directly()
}

func TestSyntacticScopeCreator_BlockScopeWithClass(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "if (true) { if (true) { class X {} } }"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("X")

	firstLevelBlock := root.GetFirstChild().GetLastChild()
	firstLevelBlockScope := s.scopeCreator.CreateScopeOfScope(firstLevelBlock, globalScope)
	assertScope(t, firstLevelBlockScope).doesNotDeclare("X")

	secondLevelBlock := firstLevelBlock.GetFirstChild().GetLastChild()
	secondLevelBLockScope := s.scopeCreator.CreateScopeOfScope(secondLevelBlock, firstLevelBlockScope)
	assertScope(t, secondLevelBLockScope).declares("X").directly()
}

func TestSyntacticScopeCreator_ClassFieldsThisAndSuper(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`class Foo {
  a = this.a
  [this.a] = this.a
}
class Bar extends Foo {
  b = super.a
}
`

	root := s.getRoot(js)
	classFoo := root.GetFirstChild()
	classBar := root.GetLastChild()
	memberFieldDefA := classFoo.GetLastChild().GetFirstChild()
	memberFieldDefB := classBar.GetLastChild().GetFirstChild()
	computedFieldDef := classFoo.GetLastChild().GetLastChild()

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	fooScope := s.scopeCreator.CreateScopeOfScope(classFoo, globalScope)
	barScope := s.scopeCreator.CreateScopeOfScope(classBar, globalScope)
	memberFieldDefAScope := s.scopeCreator.CreateScopeOfScope(memberFieldDefA, fooScope)
	memberFieldDefBScope := s.scopeCreator.CreateScopeOfScope(memberFieldDefB, barScope)
	computedFieldDefRhsScope := s.scopeCreator.CreateScopeOfScope(computedFieldDef, fooScope)

	assertScope(t, globalScope).declares("Foo").directly()
	assertScope(t, globalScope).declares("Bar").directly()
	assertScope(t, globalScope).doesNotDeclare("this")
	assertScope(t, fooScope).doesNotDeclare("this")
	assertScope(t, barScope).doesNotDeclare("this")
	assertScope(t, memberFieldDefAScope).declares("this").directly()
	assertScope(t, computedFieldDefRhsScope).declares("this").directly()
	assertScope(t, globalScope).doesNotDeclare("super")
	assertScope(t, fooScope).doesNotDeclare("super")
	assertScope(t, barScope).doesNotDeclare("super")
	assertScope(t, memberFieldDefBScope).declares("super").directly()
}

func TestSyntacticScopeCreator_ClassStaticFieldsThisAndSuper(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`class Foo {
  static a = 2
  static b = this.b
  static [this.a] = this.a
}
class Bar extends Foo {
  static c = super.a + 1
}
`

	root := s.getRoot(js)
	classFoo := root.GetFirstChild()
	classBar := root.GetLastChild()
	memberFieldDefA := classFoo.GetLastChild().GetFirstChild()
	memberFieldDefB := classBar.GetLastChild().GetFirstChild()
	computedFieldDef := classFoo.GetLastChild().GetLastChild()

	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	fooScope := s.scopeCreator.CreateScopeOfScope(classFoo, globalScope)
	barScope := s.scopeCreator.CreateScopeOfScope(classBar, globalScope)
	memberFieldDefAScope := s.scopeCreator.CreateScopeOfScope(memberFieldDefA, fooScope)
	memberFieldDefBScope := s.scopeCreator.CreateScopeOfScope(memberFieldDefB, barScope)
	computedFieldDefRhsScope := s.scopeCreator.CreateScopeOfScope(computedFieldDef, fooScope)

	assertScope(t, globalScope).declares("Foo").directly()
	assertScope(t, globalScope).declares("Bar").directly()
	assertScope(t, globalScope).doesNotDeclare("this")
	assertScope(t, fooScope).doesNotDeclare("this")
	assertScope(t, barScope).doesNotDeclare("this")
	assertScope(t, memberFieldDefAScope).declares("this").directly()
	assertScope(t, computedFieldDefRhsScope).declares("this").directly()
	assertScope(t, globalScope).doesNotDeclare("super")
	assertScope(t, fooScope).doesNotDeclare("super")
	assertScope(t, barScope).doesNotDeclare("super")
	assertScope(t, memberFieldDefBScope).declares("super").directly()
}

func TestSyntacticScopeCreator_StaticBlockScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "class C {static {}}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)
	assertScope(t, globalScope).declares("C").directly()

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)
	assertScope(t, staticBlockScope).declares("C").onSomeParent()
	assertEqual(t, "condition", staticBlockScope.IsBlockScope(), true)
	scopeTestAssertSame(t, staticBlockScope.GetClosestHoistScope(), staticBlockScope)
	assertEqual(t, "condition", staticBlockScope.IsHoistScope(), true)
}

func TestSyntacticScopeCreator_StaticBlockScopeClassExpr(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "let c = class C {static {}}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild().GetFirstFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)
	assertScope(t, globalScope).declares("c").directly()
	assertScope(t, globalScope).doesNotDeclare("C")
	assertScope(t, classScope).declares("c").onSomeParent()
	assertScope(t, classScope).declares("C").directly()

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)
	assertScope(t, staticBlockScope).declares("C").onSomeParent()
	assertScope(t, staticBlockScope).declares("c").onSomeParent()
	assertEqual(t, "condition", staticBlockScope.IsBlockScope(), true)
	scopeTestAssertSame(t, staticBlockScope.GetClosestHoistScope(), staticBlockScope)
	assertEqual(t, "condition", staticBlockScope.IsHoistScope(), true)
}

func TestSyntacticScopeCreator_StaticBlockScopeWithLet(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "class C {static {let x;}}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)
	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, classScope).doesNotDeclare("x")

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)
	assertScope(t, staticBlockScope).declares("x").directly()
}

func TestSyntacticScopeCreator_StaticBlockScopeClassExprWithLet(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "let c = class C {static {let x;}}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild().GetFirstFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)
	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, classScope).doesNotDeclare("x")

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)
	assertScope(t, staticBlockScope).declares("x").directly()
}

func TestSyntacticScopeCreator_StaticBlockScopeWithVar(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "class C {static {var x;}} var y; "
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)
	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, classScope).doesNotDeclare("x")

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)
	assertScope(t, staticBlockScope).declares("x").directly()
}

func TestSyntacticScopeCreator_StaticBlockScopeClassExprWithVar(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "let c = class C {static {var x;}}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild().GetFirstFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)
	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, classScope).doesNotDeclare("x")

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)
	assertScope(t, staticBlockScope).declares("x").directly()
}

func TestSyntacticScopeCreator_StaticBlockScopeWithLetInnerBlock(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "class C {static {if (true) {let x;}}}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)
	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, classScope).doesNotDeclare("x")

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)
	assertScope(t, staticBlockScope).doesNotDeclare("x")

	ifBodyNode := classStaticBlockNode.GetFirstChild().GetLastChild()
	ifBodyScope := s.scopeCreator.CreateScopeOfScope(ifBodyNode, staticBlockScope)
	assertScope(t, ifBodyScope).declares("x").directly()
}

func TestSyntacticScopeCreator_StaticBlockScopeClassExprWithLetInnerBlock(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "let c = class C {static {if (true) {let x;}}}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild().GetFirstFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)
	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, classScope).doesNotDeclare("x")

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)
	assertScope(t, staticBlockScope).doesNotDeclare("x")

	ifBodyNode := classStaticBlockNode.GetFirstChild().GetLastChild()
	ifBodyScope := s.scopeCreator.CreateScopeOfScope(ifBodyNode, staticBlockScope)
	assertScope(t, ifBodyScope).declares("x").directly()
}

func TestSyntacticScopeCreator_StaticBlockScopeWithVarInnerBlock(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "class C {static {if (true) {var x;}}}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)
	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, classScope).doesNotDeclare("x")

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)
	assertScope(t, staticBlockScope).declares("x").directly()

	ifBodyNode := classStaticBlockNode.GetFirstChild().GetLastChild()
	ifBodyScope := s.scopeCreator.CreateScopeOfScope(ifBodyNode, staticBlockScope)
	assertScope(t, ifBodyScope).declares("x").onClosestContainerScope()
}

func TestSyntacticScopeCreator_StaticBlockScopeClassExprWithVarInnerBlock(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "let c = class C {static {if (true) {var x;}}}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild().GetFirstFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)
	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, classScope).doesNotDeclare("x")

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)
	assertScope(t, staticBlockScope).declares("x").directly()

	ifBodyNode := classStaticBlockNode.GetFirstChild().GetLastChild()
	ifBodyScope := s.scopeCreator.CreateScopeOfScope(ifBodyNode, staticBlockScope)
	assertScope(t, ifBodyScope).declares("x").onClosestContainerScope()
}

func TestSyntacticScopeCreator_ClassStaticBlockWithLoop(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "let f = class Foo { static { for (;;) { var x; }}}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	classDecl := root.GetFirstChild().GetFirstFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classDecl, globalScope)

	classStaticBlockNode := classDecl.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)

	forBodyNode := classStaticBlockNode.GetFirstChild().GetLastChild()
	forBodyScope := s.scopeCreator.CreateScopeOfScope(forBodyNode, staticBlockScope)

	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, classScope).doesNotDeclare("x")
	assertScope(t, forBodyScope).declares("x").onClosestContainerScope()
}

func TestSyntacticScopeCreator_StaticBlockThis(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`class Foo {
  static {
    this.x
  }
}
`
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	classFoo := root.GetFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classFoo, globalScope)
	classStaticBlockNode := classFoo.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classScope)

	assertEqual(t, "condition", staticBlockScope.IsBlockScope(), true)
	scopeTestAssertSame(t, staticBlockScope.GetClosestHoistScope(), staticBlockScope)
	assertEqual(t, "condition", staticBlockScope.IsHoistScope(), true)

	assertScope(t, globalScope).declares("Foo").directly()
	assertScope(t, staticBlockScope).declares("Foo").onSomeParent()
	assertScope(t, globalScope).doesNotDeclare("this")
	assertScope(t, classScope).doesNotDeclare("this")
	assertScope(t, staticBlockScope).declares("this").directly()
}

func TestSyntacticScopeCreator_ClassStaticBlockSuper(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`class Foo {
  static x
}
class Bar extends Foo {
  static {
    super.x = 'str'
  }
}
`
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	classFoo := root.GetFirstChild()
	classFooScope := s.scopeCreator.CreateScopeOfScope(classFoo, globalScope)
	classBar := root.GetLastChild()
	classBarScope := s.scopeCreator.CreateScopeOfScope(classBar, globalScope)
	classStaticBlockNode := classBar.GetLastChild().GetFirstChild()
	staticBlockScope := s.scopeCreator.CreateScopeOfScope(classStaticBlockNode, classBarScope)

	assertEqual(t, "condition", staticBlockScope.IsBlockScope(), true)
	scopeTestAssertSame(t, staticBlockScope.GetClosestHoistScope(), staticBlockScope)
	assertEqual(t, "condition", staticBlockScope.IsHoistScope(), true)

	assertScope(t, globalScope).declares("Foo").directly()
	assertScope(t, staticBlockScope).declares("Foo").onSomeParent()
	assertScope(t, globalScope).declares("Bar").directly()
	assertScope(t, staticBlockScope).declares("Bar").onSomeParent()
	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, globalScope).doesNotDeclare("super")
	assertScope(t, classFooScope).doesNotDeclare("x")
	assertScope(t, classFooScope).doesNotDeclare("super")
	assertScope(t, classBarScope).doesNotDeclare("x")
	assertScope(t, classBarScope).doesNotDeclare("super")
	assertScope(t, staticBlockScope).doesNotDeclare("x")
	assertScope(t, staticBlockScope).declares("super").directly()
}

func TestSyntacticScopeCreator_SwitchScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`switch (b) {
  case 1:
    b
  case 2:
    let c = 4
    c
}
`
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("c")

	switchNode := root.GetFirstChild()
	switchScope := s.scopeCreator.CreateScopeOfScope(switchNode.GetLastChild(), globalScope)
	assertScope(t, switchScope).declares("c").directly()
}

func TestSyntacticScopeCreator_ForLoopScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "for (let i = 0;;) { let x; }"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("i")
	assertScope(t, globalScope).doesNotDeclare("x")

	forNode := root.GetFirstChild()
	forScope := s.scopeCreator.CreateScopeOfScope(forNode, globalScope)
	assertScope(t, forScope).declares("i").directly()
	assertScope(t, forScope).doesNotDeclare("x")

	forBlock := forNode.GetLastChild()
	forBlockScope := s.scopeCreator.CreateScopeOfScope(forBlock, forScope)
	assertScope(t, forBlockScope).declares("i").onSomeParent()
	assertScope(t, forBlockScope).declares("x").directly()
}

func TestSyntacticScopeCreator_ForOfLoopScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "for (let i of arr) { let x; }"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("i")
	assertScope(t, globalScope).doesNotDeclare("x")

	forNode := root.GetFirstChild()
	forScope := s.scopeCreator.CreateScopeOfScope(forNode, globalScope)
	assertScope(t, forScope).declares("i").directly()
	assertScope(t, forScope).doesNotDeclare("x")

	forBlock := forNode.GetLastChild()
	forBlockScope := s.scopeCreator.CreateScopeOfScope(forBlock, forScope)
	assertScope(t, forBlockScope).declares("i").onSomeParent()
	assertScope(t, forBlockScope).declares("x").directly()
}

func TestSyntacticScopeCreator_FunctionArgument(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function f(x) { if (true) { let y = 3; } }"
	root := s.getRoot(js)
	global := s.scopeCreator.CreateScopeOfScope(root, nil)
	function := root.GetLastChild()
	jutil.CheckStateMsg1(function.IsFunction(), "%s", function)
	functionScope := s.scopeCreator.CreateScopeOfScope(function, global)

	functionBlock := NodeUtil_GetFunctionBody(function)
	fBlockScope := s.scopeCreator.CreateScopeOfScope(functionBlock, functionScope)

	assertScope(t, fBlockScope).declares("x").on(functionScope)
	assertScope(t, fBlockScope).doesNotDeclare("y")

	ifBlock := functionBlock.GetLastChild().GetLastChild()
	jutil.CheckStateMsg1(ifBlock.IsBlock(), "%s", ifBlock)
	blockScope := s.scopeCreator.CreateScopeOfScope(ifBlock, fBlockScope)
	assertScope(t, blockScope).declares("x").on(functionScope)
	assertScope(t, blockScope).declares("y").directly()
}

func TestSyntacticScopeCreator_TheArgumentsVariable(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function f() { if (true) { let arguments = 3; } }"
	root := s.getRoot(js)
	global := s.scopeCreator.CreateScopeOfScope(root, nil)

	function := root.GetFirstChild()
	jutil.CheckStateMsg1(function.IsFunction(), "%s", function)
	fScope := s.scopeCreator.CreateScopeOfScope(function, global)
	// Check "has" before "get" to ensure it works before lazy instantiation.
	assertEqual(t, "condition", fScope.HasOwnSlot("arguments"), true)
	assertEqual(t, "condition", fScope.HasSlot("arguments"), true)
	arguments := fScope.GetArgumentsVar()
	scopeTestAssertSame(t, fScope.GetVar("arguments"), arguments)

	fBlock := NodeUtil_GetFunctionBody(function)
	fBlockScope := s.scopeCreator.CreateScopeOfScope(fBlock, fScope)
	scopeTestAssertSame(t, fBlockScope.GetVar("arguments"), arguments)
	scopeTestAssertSame(t, fBlockScope.GetArgumentsVar(), arguments)

	ifBlock := fBlock.GetFirstChild().GetLastChild()
	blockScope := s.scopeCreator.CreateScopeOfScope(ifBlock, fBlockScope)
	assertScope(t, blockScope).declares("arguments").directly()
	scopeTestAssertSame(t, blockScope.GetArgumentsVar(), arguments)
	// Java: isNotEqualTo — Var equality is ScopedName.equals (name and scope root).
	if blockScope.GetVar("arguments").Equals(arguments) {
		t.Errorf("blockScope.getVar(arguments) equals the function's arguments var %v", arguments)
	}
}

func TestSyntacticScopeCreator_ArgumentsVariableInArrowFunction(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function outer() { var inner = () => { alert(0); } }"
	root := s.getRoot(js)
	global := s.scopeCreator.CreateScopeOfScope(root, nil)

	outer := root.GetFirstChild()
	jutil.CheckStateMsg1(outer.IsFunction(), "%s", outer)
	jutil.CheckStateMsg1(!outer.IsArrowFunction(), "%s", outer)
	outerFunctionScope := s.scopeCreator.CreateScopeOfScope(outer, global)
	arguments := outerFunctionScope.GetArgumentsVar()

	outerBody := NodeUtil_GetFunctionBody(outer)
	outerBodyScope := s.scopeCreator.CreateScopeOfScope(outerBody, outerFunctionScope)

	inner :=
		outerBody.
			GetFirstChild(). // VAR
			GetFirstChild(). // NAME
			GetFirstChild()  // FUNCTION
	jutil.CheckStateMsg1(inner.IsFunction(), "%s", inner)
	jutil.CheckStateMsg1(inner.IsArrowFunction(), "%s", inner)
	innerFunctionScope := s.scopeCreator.CreateScopeOfScope(inner, outerBodyScope)
	scopeTestAssertSame(t, innerFunctionScope.GetArgumentsVar(), arguments)
}

func TestSyntacticScopeCreator_TheThisVariable(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function f() { if (true) { function g() {} } }"
	root := s.getRoot(js)
	global := s.scopeCreator.CreateScopeOfScope(root, nil)

	function := root.GetFirstChild()
	jutil.CheckStateMsg1(function.IsFunction(), "%s", function)
	fScope := s.scopeCreator.CreateScopeOfScope(function, global)
	// Check "has" before "get" to ensure it works before lazy instantiation.
	assertEqual(t, "condition", fScope.HasOwnSlot("this"), true)
	assertEqual(t, "condition", fScope.HasSlot("this"), true)
	thisVar := fScope.GetVar("this")
	assertEqual(t, "condition", thisVar.IsThis(), true)

	fBlock := NodeUtil_GetFunctionBody(function)
	fBlockScope := s.scopeCreator.CreateScopeOfScope(fBlock, fScope)
	assertScope(t, fBlockScope).declares("this")
	scopeTestAssertSame(t, fBlockScope.GetVar("this"), thisVar)

	ifBlock := fBlock.GetFirstChild().GetLastChild()
	blockScope := s.scopeCreator.CreateScopeOfScope(ifBlock, fBlockScope)
	assertScope(t, blockScope).declares("this")
	scopeTestAssertSame(t, blockScope.GetVar("this"), thisVar)
	scopeTestAssertSame(t, blockScope.GetVar("this").GetScope(), rhino.StaticScope(fScope))

	gFunction := ifBlock.GetFirstChild()
	gScope := s.scopeCreator.CreateScopeOfScope(gFunction, blockScope)
	assertScope(t, gScope).declares("this").directly()
	scopeTestAssertSame(t, gScope.GetVar("this").GetScope(), rhino.StaticScope(gScope))
}

func TestSyntacticScopeCreator_TheSuperVariable(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function f() { if (true) { function g() {} } }"
	root := s.getRoot(js)
	global := s.scopeCreator.CreateScopeOfScope(root, nil)

	function := root.GetFirstChild()
	jutil.CheckStateMsg1(function.IsFunction(), "%s", function)
	fScope := s.scopeCreator.CreateScopeOfScope(function, global)
	// Check "has" before "get" to ensure it works before lazy instantiation.
	assertEqual(t, "condition", fScope.HasOwnSlot("super"), true)
	assertEqual(t, "condition", fScope.HasSlot("super"), true)
	thisVar := fScope.GetVar("super")

	fBlock := NodeUtil_GetFunctionBody(function)
	fBlockScope := s.scopeCreator.CreateScopeOfScope(fBlock, fScope)
	assertScope(t, fBlockScope).declares("super")
	scopeTestAssertSame(t, fBlockScope.GetVar("super"), thisVar)

	ifBlock := fBlock.GetFirstChild().GetLastChild()
	blockScope := s.scopeCreator.CreateScopeOfScope(ifBlock, fBlockScope)
	assertScope(t, blockScope).declares("super")
	scopeTestAssertSame(t, blockScope.GetVar("super"), thisVar)
	scopeTestAssertSame(t, blockScope.GetVar("super").GetScope(), rhino.StaticScope(fScope))

	gFunction := ifBlock.GetFirstChild()
	gScope := s.scopeCreator.CreateScopeOfScope(gFunction, blockScope)
	assertScope(t, gScope).declares("super").directly()
	scopeTestAssertSame(t, gScope.GetVar("super").GetScope(), rhino.StaticScope(gScope))
	scopeTestAssertNotSame(t, gScope.GetVar("super"), gScope.GetVar("this"))
}

func TestSyntacticScopeCreator_TheThisVariableInArrowFunction(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function outer() { var inner = () => this.x; }"
	root := s.getRoot(js)
	global := s.scopeCreator.CreateScopeOfScope(root, nil)

	outer := root.GetFirstChild()
	jutil.CheckStateMsg1(outer.IsFunction(), "%s", outer)
	jutil.CheckStateMsg1(!outer.IsArrowFunction(), "%s", outer)
	outerFunctionScope := s.scopeCreator.CreateScopeOfScope(outer, global)
	thisVar := outerFunctionScope.GetVar("this")

	outerBody := NodeUtil_GetFunctionBody(outer)
	outerBodyScope := s.scopeCreator.CreateScopeOfScope(outerBody, outerFunctionScope)

	inner :=
		outerBody.
			GetFirstChild(). // VAR
			GetFirstChild(). // NAME
			GetFirstChild()  // FUNCTION
	jutil.CheckStateMsg1(inner.IsFunction(), "%s", inner)
	jutil.CheckStateMsg1(inner.IsArrowFunction(), "%s", inner)
	innerFunctionScope := s.scopeCreator.CreateScopeOfScope(inner, outerBodyScope)
	scopeTestAssertSame(t, innerFunctionScope.GetVar("this"), thisVar)
}

func TestSyntacticScopeCreator_IsFunctionBlockScoped(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "if (true) { function f() {}; }"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("f")

	ifBlock := root.GetFirstChild().GetLastChild()
	blockScope := s.scopeCreator.CreateScopeOfScope(ifBlock, globalScope)
	assertScope(t, blockScope).declares("f").directly()
}

func TestSyntacticScopeCreator_IsClassBlockScoped(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "if (true) { class X {}; }"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("X")

	ifBlock := root.GetFirstChild().GetLastChild()
	blockScope := s.scopeCreator.CreateScopeOfScope(ifBlock, globalScope)
	assertScope(t, blockScope).declares("X").directly()
}

func TestSyntacticScopeCreator_IsCatchBlockScoped(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "try { var x = 2; } catch (e) { var y = 3; let z = 4; }"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).declares("x").directly()
	assertScope(t, globalScope).declares("y").directly()
	assertScope(t, globalScope).doesNotDeclare("z")
	assertScope(t, globalScope).doesNotDeclare("e")

	tryBlock := root.GetFirstFirstChild()
	tryBlockScope := s.scopeCreator.CreateScopeOfScope(tryBlock, globalScope)
	assertScope(t, tryBlockScope).declares("x").onSomeParent()
	assertScope(t, tryBlockScope).declares("y").onSomeParent()
	assertScope(t, tryBlockScope).doesNotDeclare("z")
	assertScope(t, tryBlockScope).doesNotDeclare("e")

	catchBlock := tryBlock.GetNext()
	catchBlockScope := s.scopeCreator.CreateScopeOfScope(catchBlock, tryBlockScope)
	assertScope(t, catchBlockScope).declares("x").onSomeParent()
	assertScope(t, catchBlockScope).declares("y").onSomeParent()
	assertScope(t, catchBlockScope).declares("z").directly()
	assertScope(t, catchBlockScope).declares("e").directly()
}

func TestSyntacticScopeCreator_Import(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`import * as ns from 'm1'
import d from 'm2'
import {foo} from 'm3'
import {x as y} from 'm4'
`

	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertEqual(t, "globalScope var count", len(scopeTestVarNames(globalScope.GetVarIterable())), 0)

	moduleBody := root.GetFirstChild()
	jutil.CheckStateMsg1(moduleBody.IsModuleBody(), "%s", moduleBody)
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, globalScope)
	assertScope(t, moduleScope).declares("ns").directly()
	assertScope(t, moduleScope).declares("d").directly()
	assertScope(t, moduleScope).declares("foo").directly()
	assertScope(t, moduleScope).declares("y").directly()
	assertScope(t, moduleScope).doesNotDeclare("x")
}

func TestSyntacticScopeCreator_ImportAsSelf(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "import {x as x} from 'm';"

	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertEqual(t, "globalScope var count", len(scopeTestVarNames(globalScope.GetVarIterable())), 0)

	moduleBody := root.GetFirstChild()
	jutil.CheckStateMsg1(moduleBody.IsModuleBody(), "%s", moduleBody)
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, globalScope)
	assertScope(t, moduleScope).declares("x").directly()
}

func TestSyntacticScopeCreator_ImportDefault(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "import x from 'm';"

	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertEqual(t, "globalScope var count", len(scopeTestVarNames(globalScope.GetVarIterable())), 0)

	moduleBody := root.GetFirstChild()
	jutil.CheckStateMsg1(moduleBody.IsModuleBody(), "%s", moduleBody)
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, globalScope)
	assertScope(t, moduleScope).declares("x").directly()
}

func TestSyntacticScopeCreator_ModuleScoped(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "export function f() { var x; if (1) { let y; } }; var z;"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("f")
	assertScope(t, globalScope).doesNotDeclare("x")
	assertScope(t, globalScope).doesNotDeclare("y")
	assertScope(t, globalScope).doesNotDeclare("z")

	moduleBlock := root.GetFirstChild()
	moduleBlockScope := s.scopeCreator.CreateScopeOfScope(moduleBlock, globalScope)
	assertScope(t, moduleBlockScope).declares("f").directly()
	assertScope(t, moduleBlockScope).doesNotDeclare("x")
	assertScope(t, moduleBlockScope).doesNotDeclare("y")
	assertScope(t, moduleBlockScope).declares("z").directly()
}

func TestSyntacticScopeCreator_ExportDefault(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "export default function f() {};"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("f")

	moduleBlock := root.GetFirstChild()
	moduleBlockScope := s.scopeCreator.CreateScopeOfScope(moduleBlock, globalScope)
	assertScope(t, moduleBlockScope).declares("f").directly()
}

func TestSyntacticScopeCreator_ExportFrom(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "export {PI} from './n.js';"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("PI")

	moduleBlock := root.GetFirstChild()
	moduleBlockScope := s.scopeCreator.CreateScopeOfScope(moduleBlock, globalScope)
	assertScope(t, moduleBlockScope).doesNotDeclare("PI")
}

func TestSyntacticScopeCreator_VarAfterLet(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js :=
		`function f() {
  if (a) {
    let x
  }
  var y
}
`

	root := s.getRoot(js)
	global := s.scopeCreator.CreateScopeOfScope(root, nil)
	function := root.GetFirstChild()
	fScope := s.scopeCreator.CreateScopeOfScope(function, global)

	fBlock := root.GetFirstChild().GetLastChild()
	fBlockScope := s.scopeCreator.CreateScopeOfScope(fBlock, fScope)
	jutil.CheckNotNull(fBlockScope)
	assertScope(t, fBlockScope).doesNotDeclare("x")
	assertScope(t, fBlockScope).declares("y").directly()

	ifBlock := fBlock.GetFirstChild().GetLastChild()
	ifBlockScope := s.scopeCreator.CreateScopeOfScope(ifBlock, fBlockScope)
	assertScope(t, ifBlockScope).declares("x").directly()
	assertScope(t, ifBlockScope).declares("y").onSomeParent()
}

func TestSyntacticScopeCreator_SimpleFunctionParam(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function f(x) {}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	fNode := root.GetFirstChild()
	jutil.CheckStateMsg1(fNode.IsFunction(), "%s", fNode)
	fScope := s.scopeCreator.CreateScopeOfScope(fNode, globalScope)
	assertScope(t, fScope).declares("x").directly()

	fBlock := NodeUtil_GetFunctionBody(fNode)
	fBlockScope := s.scopeCreator.CreateScopeOfScope(fBlock, fScope)
	assertScope(t, fBlockScope).declares("x").on(fScope)
}

func TestSyntacticScopeCreator_OnlyOneDeclaration(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function f(x) { if (!x) var x = 6; }"
	root := s.getRoot(js)
	fNode := root.GetFirstChild()
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	fScope := s.scopeCreator.CreateScopeOfScope(fNode, globalScope)
	assertScope(t, fScope).declares("x").directly()

	fBlock := fNode.GetLastChild()
	fBlockScope := s.scopeCreator.CreateScopeOfScope(fBlock, fScope)
	assertScope(t, fBlockScope).declares("x").on(fScope)

	ifBlock := fBlock.GetFirstChild().GetLastChild()
	ifBlockScope := s.scopeCreator.CreateScopeOfScope(ifBlock, fBlockScope)
	assertScope(t, ifBlockScope).declares("x").on(fScope)
}

func TestSyntacticScopeCreator_CatchInFunction(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "function f(e) { try {} catch (e) {} }"
	root := s.getRoot(js)
	fNode := root.GetFirstChild()
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	fScope := s.scopeCreator.CreateScopeOfScope(fNode, globalScope)
	assertScope(t, fScope).declares("e").directly()

	fBlock := fNode.GetLastChild()
	fBlockScope := s.scopeCreator.CreateScopeOfScope(fBlock, fScope)
	tryBlock := fBlock.GetFirstFirstChild()
	tryScope := s.scopeCreator.CreateScopeOfScope(tryBlock, fBlockScope)
	catchBlock := tryBlock.GetNext()
	catchScope := s.scopeCreator.CreateScopeOfScope(catchBlock, tryScope)
	assertScope(t, catchScope).declares("e").directly()
}

func TestSyntacticScopeCreator_FunctionName(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "var f = function foo() {}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).declares("f").directly()
	assertScope(t, globalScope).doesNotDeclare("foo")

	fNode := root.GetFirstChild().GetFirstFirstChild()
	fScope := s.scopeCreator.CreateScopeOfScope(fNode, globalScope)
	assertScope(t, fScope).declares("f").onSomeParent()
	assertScope(t, fScope).declares("foo").directly()
}

func TestSyntacticScopeCreator_FunctionNameMatchesParamName1(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "var f = function foo(foo) {}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).declares("f").directly()
	assertScope(t, globalScope).doesNotDeclare("foo")

	fNode := root.GetFirstChild().GetFirstFirstChild()
	fScope := s.scopeCreator.CreateScopeOfScope(fNode, globalScope)
	assertScope(t, fScope).declares("f").onSomeParent()
	assertScope(t, fScope).declares("foo").directly()

	// The parameter 'foo', not the function name, is the declaration of the variable 'foo' in this
	// scope.
	assertEqual(t, "token", fScope.GetVar("foo").GetNode().GetParent().GetToken(), rhino.Token_PARAM_LIST)
}

func TestSyntacticScopeCreator_FunctionNameMatchesParamName2(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "var f = function foo(x = foo, foo) {}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).declares("f").directly()
	assertScope(t, globalScope).doesNotDeclare("foo")

	fNode := root.GetFirstChild().GetFirstFirstChild()
	fScope := s.scopeCreator.CreateScopeOfScope(fNode, globalScope)
	assertScope(t, fScope).declares("f").onSomeParent()
	assertScope(t, fScope).declares("foo").directly()

	// The parameter 'foo', not the function name, is the declaration of the variable 'foo' in this
	// scope.
	assertEqual(t, "token", fScope.GetVar("foo").GetNode().GetParent().GetToken(), rhino.Token_PARAM_LIST)
}

func TestSyntacticScopeCreator_ClassName(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "var Clazz = class Foo {}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).declares("Clazz").directly()
	assertScope(t, globalScope).doesNotDeclare("Foo")

	classNode := root.GetFirstChild().GetFirstFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classNode, globalScope)
	assertScope(t, classScope).declares("Clazz").onSomeParent()
	assertScope(t, classScope).declares("Foo").directly()
}

func TestSyntacticScopeCreator_FunctionExpressionInForLoopInitializer(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("for (function foo() {};;) {}")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("foo")

	forNode := root.GetFirstChild()
	forScope := s.scopeCreator.CreateScopeOfScope(forNode, globalScope)
	assertScope(t, forScope).doesNotDeclare("foo")

	fNode := forNode.GetFirstChild()
	fScope := s.scopeCreator.CreateScopeOfScope(fNode, forScope)
	assertScope(t, fScope).declares("foo").directly()
}

func TestSyntacticScopeCreator_ClassExpressionInForLoopInitializer(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("for (class Clazz {};;) {}")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("Clazz")

	forNode := root.GetFirstChild()
	forScope := s.scopeCreator.CreateScopeOfScope(forNode, globalScope)
	assertScope(t, forScope).doesNotDeclare("Clazz")

	classNode := forNode.GetFirstChild()
	classScope := s.scopeCreator.CreateScopeOfScope(classNode, forScope)
	assertScope(t, classScope).declares("Clazz").directly()
}

func TestSyntacticScopeCreator_ClassDeclarationInExportDefault(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	js := "export default class Clazz {}"
	root := s.getRoot(js)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("Clazz")

	moduleBody := root.GetFirstChild()
	jutil.CheckStateMsg1(moduleBody.IsModuleBody(), "%s", moduleBody)
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, globalScope)
	assertScope(t, moduleScope).declares("Clazz").directly()
}

func TestSyntacticScopeCreator_VarsInModulesNotGlobal(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("goog.module('example'); var x;")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("x")

	moduleBody := root.GetFirstChild()
	jutil.CheckStateMsg1(moduleBody.IsModuleBody(), "%s", moduleBody)
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, globalScope)
	assertScope(t, moduleScope).declares("x").directly()
}

func TestSyntacticScopeCreator_GoogModuleDeclaresImplicitExports(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("goog.module('example');")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("exports")

	moduleBody := root.GetFirstChild()
	jutil.CheckStateMsg1(moduleBody.IsModuleBody(), "%s", moduleBody)
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, globalScope)
	assertScope(t, moduleScope).declares("exports").directly()
}

func TestSyntacticScopeCreator_GoogModuleCanOverrideImplicitExports(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("goog.module('example'); var exports = {};")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).doesNotDeclare("exports")

	moduleBody := root.GetFirstChild()
	jutil.CheckStateMsg1(moduleBody.IsModuleBody(), "%s", moduleBody)
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, globalScope)
	assertScope(t, moduleScope).declares("exports").directly()
}

func TestSyntacticScopeCreator_GoogProvideOfNameInScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("goog.provide('foo');")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).declares("foo").directly()
	assertEqual(t, "condition", globalScope.GetVar("foo").IsImplicitGoogNamespace(), true)
}

func TestSyntacticScopeCreator_GoogProvideOfNamespaceInScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("goog.provide('foo.bar');")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	assertScope(t, globalScope).declares("foo").directly()
	assertScope(t, globalScope).doesNotDeclare("foo.bar")
}

func TestSyntacticScopeCreator_TwoGoogProvidesOfNameInScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("goog.provide('foo.bar'); goog.provide('foo.baz');")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).declares("foo").directly()
	assertEqual(t, "condition", globalScope.GetVar("foo").IsImplicitGoogNamespace(), true)
	s.assertNoRedeclarations()
}

func TestSyntacticScopeCreator_TwoGoogProvidesOfNameInScope_withTreatProvidesAsRedeclarations(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	s.scopeCreator =
		NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations(
			s.compiler,
			newSyntacticScopeCreatorTest_RecordingRedeclarationHandler(s),
			/* treatProvidesAsRedeclarations= */ true)

	root := s.getRoot("goog.provide('foo.bar'); goog.provide('foo.baz');")
	s.scopeCreator.CreateScopeOfScope(root, nil)

	s.assertNoRedeclarations()
}

func TestSyntacticScopeCreator_VarFollowedByGoogProvidesOfNameInScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("var foo; goog.provide('foo.bar');")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertScope(t, globalScope).declares("foo").directly()
	assertEqual(t, "condition", globalScope.GetVar("foo").IsImplicitGoogNamespace(), false)
	s.assertNoRedeclarations()
}

func TestSyntacticScopeCreator_VarFollowedByGoogProvidesOfNameInScope_withTreatProvidesAsRedeclarations(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	s.scopeCreator =
		NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations(
			s.compiler,
			newSyntacticScopeCreatorTest_RecordingRedeclarationHandler(s),
			/* treatProvidesAsRedeclarations= */ true)

	root := s.getRoot("var foo; goog.provide('foo.bar'); goog.provide('foo');")
	s.scopeCreator.CreateScopeOfScope(root, nil)

	s.assertRedeclarationsSize(2)
}

func TestSyntacticScopeCreator_LegacyGoogModuleNamespaceInScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("goog.module('foo.bar'); goog.module.declareLegacyNamespace();")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	assertScope(t, globalScope).declares("foo").directly()
	assertScope(t, globalScope).doesNotDeclare("foo.bar")

	moduleBody := root.GetFirstChild()
	jutil.CheckStateMsg1(moduleBody.IsModuleBody(), "%s", moduleBody)
	moduleScope := s.scopeCreator.CreateScopeOfScope(moduleBody, globalScope)
	assertScope(t, moduleScope).declares("foo").on(globalScope)
}

func TestSyntacticScopeCreator_TwoLegacyNamespaceGoogModules(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	file1 := s.getRoot("goog.module('foo.bar'); goog.module.declareLegacyNamespace();")
	file2 := s.getRoot("goog.module('foo.baz'); goog.module.declareLegacyNamespace();")

	globalScope := s.scopeCreator.CreateScopeOfScope(rhino.IR_Root(file1, file2), nil)
	assertScope(t, globalScope).declares("foo").directly()
	assertEqual(t, "condition", globalScope.GetVar("foo").IsImplicitGoogNamespace(), true)
	s.assertNoRedeclarations()
}

func TestSyntacticScopeCreator_TwoLegacyNamespaceGoogModules_withTreatProvidesAsRedeclarations(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	s.scopeCreator =
		NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations(
			s.compiler,
			newSyntacticScopeCreatorTest_RecordingRedeclarationHandler(s),
			/* treatProvidesAsRedeclarations= */ true)

	file1 := s.getRoot("goog.module('foo.bar'); goog.module.declareLegacyNamespace();")
	file2 := s.getRoot("goog.module('foo.baz'); goog.module.declareLegacyNamespace();")

	globalScope := s.scopeCreator.CreateScopeOfScope(rhino.IR_Root(file1, file2), nil)
	assertScope(t, globalScope).declares("foo").directly()
	assertEqual(t, "condition", globalScope.GetVar("foo").IsImplicitGoogNamespace(), true)
	s.assertNoRedeclarations()
}

func TestSyntacticScopeCreator_VarFollowedByLegacyGoogModuleNamespace_withTreatProvidesAsRedeclarations(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	s.scopeCreator =
		NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations(
			s.compiler,
			newSyntacticScopeCreatorTest_RecordingRedeclarationHandler(s),
			/* treatProvidesAsRedeclarations= */ true)

	file1 := s.getRoot("var foo;")
	file2 := s.getRoot("goog.module('foo.bar'); goog.module.declareLegacyNamespace();")
	globalScope := s.scopeCreator.CreateScopeOfScope(rhino.IR_Root(file1, file2), nil)

	assertScope(t, globalScope).declares("foo").directly()
	assertEqual(t, "condition", globalScope.GetVar("foo").IsImplicitGoogNamespace(), false)
	s.assertRedeclarationsSize(1)
}

func TestSyntacticScopeCreator_BundledLegacyGoogModuleNamespaceInScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root :=
		s.getRoot(
			`goog.loadModule(function(exports) {
  goog.module('foo.bar')
  goog.module.declareLegacyNamespace()
  return exports
})
`)
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	assertScope(t, globalScope).declares("foo").directly()
	assertScope(t, globalScope).doesNotDeclare("foo.bar")
}

func TestSyntacticScopeCreator_NonLegacyGoogModuleNamespace_notInScope(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("goog.module('foo.bar');")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)

	assertScope(t, globalScope).doesNotDeclare("foo")
}

func TestSyntacticScopeCreator_GoogProvideDeclaresStrength(t *testing.T) {
	s := newSyntacticScopeCreatorTest(t)
	root := s.getRoot("goog.provide('foo');")
	globalScope := s.scopeCreator.CreateScopeOfScope(root, nil)
	assertEqual(t, "value", globalScope.GetVar("foo").GetImplicitGoogNamespaceStrength(), rhino.StaticSourceFile_SourceKind_STRONG)
}

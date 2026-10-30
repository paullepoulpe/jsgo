// Port of com/google/javascript/jscomp/SyntacticScopeCreator.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/MemoizedScopeCreator.java @ closure-compiler v20261021

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

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ScopeCreator defines how objects capable of creating scopes from the parse
// tree behave.
type ScopeCreator interface {
	// CreateScope creates a Scope object.
	//
	// n is the root node (either a FUNCTION node, a SCRIPT node, or a
	// synthetic block node whose children are all SCRIPT nodes); parent is
	// the parent Scope object (may be null).
	//
	// Go-only: language-adapter — implementations whose Java override narrows
	// the result (SyntacticScopeCreator returns Scope, TypedScopeCreator
	// TypedScope) keep this signature; callers assert (`.(*Scope)`).
	//
	// Java: ScopeCreator.createScope(Node, AbstractScope<?, ?>)
	CreateScope(n rhino.Node, parent AbstractScope) AbstractScope
}

// ScopeCreator_ASSERT_NO_SCOPES_CREATED is a scope creator that can be
// provided to NodeTraversal to ensure that no scopes are actually ever
// created. This is in contrast to the default SyntacticScopeCreator, which
// will create scopes on demand.
var ScopeCreator_ASSERT_NO_SCOPES_CREATED ScopeCreator = &scopeCreator_assertNoScopesCreated{}

// scopeCreator_assertNoScopesCreated is the anonymous class of
// ASSERT_NO_SCOPES_CREATED.
type scopeCreator_assertNoScopesCreated struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: ScopeCreator.ASSERT_NO_SCOPES_CREATED.createScope(Node, AbstractScope<?, ?>)
func (*scopeCreator_assertNoScopesCreated) CreateScope(n rhino.Node, parent AbstractScope) AbstractScope {
	panic(jutil.NewIllegalStateExceptionMsg(
		"ScopeCreator.ASSERT_NO_SCOPES_CREATED cannot create child scopes."))
}

// Go-only: this pass reads declarations through the typed view of package
// ast (fn.Params(), decl.FirstDeclarator(), catch.Param(), imp.Bindings()...)
// instead of child positions: each declaration shape has its own small
// method (declareNameDeclaration, declareParam, declarePattern...) where
// Java visits every left-hand side with one generic switch. The walk that
// finds the declarations (ScanVars) stays on rhino nodes: it visits the
// children of a dozen statement kinds alike, filtered by one token table
// (scanVarsTokenKind), and must read the next child before the current one
// is handled, because a redeclaration handler may replace it. The body of
// a for loop is a block (the grammar says so), so the Go-only code that
// scanned a loop body that is not a block is gone.

// SyntacticScopeCreator scans the parse tree to create a Scope object
// containing all the variable declarations in that scope. This class adds
// support for block-level scopes introduced in ECMAScript 6.
//
// This implementation is not thread-safe.
type SyntacticScopeCreator struct {
	compiler                      AbstractCompiler
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler
	treatProvidesAsRedeclarations bool
	lastInputId                   *rhino.InputId
	lastInput                     *CompilerInput

	// Go-only: performance — the vars this creator declares, which become
	// garbage with its scopes at the end of a pass, are allocated in chunks;
	// varScratch collects a scope's vars while it is scanned, so that its
	// varList is allocated once, at its final size.
	// scopes and varLists do the same for the scopes themselves and for the
	// var lists of small scopes (carved from shared pointer chunks); a Var
	// keeps its Scope alive anyway, so the scopes die with the vars.
	vars           jutil.Slab[Var]
	scopes         jutil.Slab[Scope]
	varLists       []*Var
	varScratch     []*Var
	scanning       bool
	skipGlobalScan bool
}

// syntacticScopeCreator_varListChunk is the size of the pointer chunks var
// lists are carved from; a list longer than a quarter of it is allocated on
// its own.
const syntacticScopeCreator_varListChunk = 508 // 4 KB with the 8-byte header of a pointerful allocation

// newVarList returns a zeroed []*Var of length and capacity n.
func (c *SyntacticScopeCreator) newVarList(n int) []*Var {
	if n > syntacticScopeCreator_varListChunk/4 {
		return make([]*Var, n)
	}
	if len(c.varLists) < n {
		c.varLists = jutil.SlabChunk[*Var](syntacticScopeCreator_varListChunk)
	}
	l := c.varLists[:n:n]
	c.varLists = c.varLists[n:]
	return l
}

// The arguments variable is special, in that it's declared for every function,
// but not explicitly declared.
const syntacticScopeCreator_ARGUMENTS = "arguments"

var SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER SyntacticScopeCreator_RedeclarationHandler = NewSyntacticScopeCreator_DefaultRedeclarationHandler()

// Java: SyntacticScopeCreator(AbstractCompiler)
func NewSyntacticScopeCreator(compiler AbstractCompiler) *SyntacticScopeCreator {
	return NewSyntacticScopeCreatorWithRedeclarationHandler(compiler, SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER)
}

// Java: SyntacticScopeCreator(AbstractCompiler, RedeclarationHandler)
func NewSyntacticScopeCreatorWithRedeclarationHandler(
	compiler AbstractCompiler, redeclarationHandler SyntacticScopeCreator_RedeclarationHandler) *SyntacticScopeCreator {
	return NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations(
		compiler, redeclarationHandler /* treatProvidesAsRedeclarations= */, false)
}

// NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations:
// treatProvidesAsRedeclarations is whether to report goog.provide/legacy
// goog.module initializations of a namespace to the redeclaration handler.
// For example, `var foo; goog.provide('foo.bar');` will report the provide
// as a redeclaration if this is enabled.
//
// Java: SyntacticScopeCreator(AbstractCompiler, RedeclarationHandler, boolean)
func NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations(
	compiler AbstractCompiler,
	redeclarationHandler SyntacticScopeCreator_RedeclarationHandler,
	treatProvidesAsRedeclarations bool) *SyntacticScopeCreator {
	return &SyntacticScopeCreator{
		compiler:                      compiler,
		redeclarationHandler:          redeclarationHandler,
		treatProvidesAsRedeclarations: treatProvidesAsRedeclarations,
	}
}

// CreateScope is the ScopeCreator method; Java narrows its result to Scope
// (see ScopeCreator).
//
// Java: SyntacticScopeCreator.createScope(Node, AbstractScope<?, ?>)
func (c *SyntacticScopeCreator) CreateScope(n rhino.Node, parent AbstractScope) AbstractScope {
	// Go-only: language-adapter — the cast (Scope) parent, which lets null
	// through.
	var p *Scope
	if parent != nil {
		p = parent.(*Scope)
	}
	return c.CreateScopeOfScope(n, p)
}

// Java: SyntacticScopeCreator.createScope(Node, Scope)
func (c *SyntacticScopeCreator) CreateScopeOfScope(n rhino.Node, parent *Scope) *Scope {
	var scope *Scope
	if parent == nil {
		scope = newScope(n)
		if c.skipGlobalScan {
			return scope
		}
	} else {
		scope = initScopeWithParent(c.scopes.New(), parent, n)
	}
	// Go-only: performance — the scanner, which nothing keeps, is a local
	// variable rather than a heap object (one per scope created).
	var scanner SyntacticScopeCreator_ScopeScanner
	initSyntacticScopeCreator_ScopeScanner(&scanner, c.compiler, c.redeclarationHandler, scope, nil, c.treatProvidesAsRedeclarations)
	scanner.creator = c
	if c.scanning {
		scanner.Populate()
	} else {
		c.scanning = true
		scope.varList = c.varScratch[:0]
		scanner.Populate()
		vs := scope.varList
		scope.varList = nil
		if len(vs) > 0 {
			scope.varList = c.newVarList(len(vs))
			copy(scope.varList, vs)
			clear(vs)
		}
		c.varScratch = vs[:0]
		c.scanning = false
	}
	scope.scanDone = true
	return scope
}

// syntacticScopeChildHasBlockDecls reports whether the statement c declares
// a name in the block scope around it: let, const, a function or class
// declaration, a catch (its parameter), an export, or one of these in a
// synthetic block, which has no scope of its own, or under a label.
//
// Go-only: performance — a block scope whose statements declare nothing is
// skipped, or created empty without a scan, and remembered as such
// (SetNoBlockScopeDecls).
func syntacticScopeChildHasBlockDecls(c rhino.Node) bool {
	switch c.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CATCH, rhino.Token_EXPORT:
		return true
	case rhino.Token_FUNCTION:
		if NodeUtil_IsFunctionExpression(c) {
			return false
		}
		_, named := ast.MustFunction(c).Name()
		return named
	case rhino.Token_CLASS:
		if NodeUtil_IsClassExpression(c) {
			return false
		}
		_, named := ast.MustClass(c).Name()
		return named
	case rhino.Token_BLOCK:
		return syntacticSyntheticBlockHasBlockDecls(ast.MustBlock(c))
	case rhino.Token_LABEL:
		return syntacticScopeChildHasBlockDecls(ast.MustLabel(c).Body().Node())
	case rhino.Token_IF:
		// The condition is an expression; only a branch that is a synthetic
		// block can declare a name here.
		ifs := ast.MustIf(c)
		if syntacticSyntheticBlockHasBlockDecls(ifs.Then()) {
			return true
		}
		els, ok := ifs.Else()
		return ok && syntacticSyntheticBlockHasBlockDecls(els)
	case rhino.Token_DO, rhino.Token_WHILE:
		return syntacticSyntheticBlockHasBlockDecls(ast.MustLoop(c).Body())
	case rhino.Token_WITH:
		return syntacticSyntheticBlockHasBlockDecls(ast.MustWith(c).Body())
	default:
		return false
	}
}

// syntacticSyntheticBlockHasBlockDecls reports whether b is a synthetic
// block (one the compiler added, with no scope of its own) whose statements
// declare a name in the block scope around it.
func syntacticSyntheticBlockHasBlockDecls(b ast.Block) bool {
	if !b.Node().IsSyntheticBlock() {
		return false
	}
	for stmt := range b.Statements() {
		if syntacticScopeChildHasBlockDecls(stmt.Node()) {
			return true
		}
	}
	return false
}

// syntacticForLetOrConst returns the let or const declared in the head of
// the for loop n: `for (let i = 0; ...)`, `for (const x of xs)`.
func syntacticForLetOrConst(n rhino.Node) (ast.NameDeclaration, bool) {
	var decl ast.NameDeclaration
	var ok bool
	if loop, isEnhanced := ast.AsEnhancedFor(n); isEnhanced {
		decl, ok = loop.Target().AsNameDeclaration()
	} else if init, hasInit := ast.MustFor(n).Init(); hasInit {
		decl, ok = init.AsNameDeclaration()
	}
	return decl, ok && decl.Token() != rhino.Token_VAR
}

// syntacticScopeHasBlockDecls reports whether n, the root of a block scope
// (a block or a for loop), may declare names; it reports true for other
// nodes. A false answer is remembered on n (SetNoBlockScopeDecls).
//
// The body of a for loop is a block, with a scope of its own: only a let or
// const in the loop's head is declared in the loop's scope.
func syntacticScopeHasBlockDecls(n rhino.Node) bool {
	if n.HasNoBlockScopeDecls() {
		return false
	}
	if NodeUtil_IsAnyFor(n) {
		if _, ok := syntacticForLetOrConst(n); ok {
			return true
		}
		n.SetNoBlockScopeDecls()
		return false
	}
	if n.IsBlock() {
		// (The children are statements, or the catch of a try.)
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if scanVarsTokenKind[uint8(c.GetToken())] == 1 && syntacticScopeChildHasBlockDecls(c) {
				return true
			}
		}
		n.SetNoBlockScopeDecls()
		return false
	}
	return true
}

// syntacticFunctionBlockHasNoDecls reports whether n, the body of a
// function or a class static block, has no statement that may declare a
// name. A true answer is remembered on n (SetNoBlockScopeDecls).
func syntacticFunctionBlockHasNoDecls(n rhino.Node) bool {
	if n.HasNoBlockScopeDecls() {
		return true
	}
	for stmt := range ast.MustBlock(n).Statements() {
		if scanVarsTokenKind[uint8(stmt.Token())] == 1 {
			return false
		}
	}
	n.SetNoBlockScopeDecls()
	return true
}

// SyntacticScopeCreator_ScopeScanner is a class to traverse the AST looking
// for name definitions and add them to the Scope.
type SyntacticScopeCreator_ScopeScanner struct {
	scope                         *Scope
	compiler                      AbstractCompiler
	creator                       *SyntacticScopeCreator
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler
	treatProvidesAsRedeclarations bool
	isFnBlock                     bool

	// Will be null, when a detached node is traversed.
	inputId         *rhino.InputId
	input           *CompilerInput
	inputIdResolved bool
	inputResolved   bool
	changeRootSet   map[rhino.Node]struct{}
}

// initSyntacticScopeCreator_ScopeScanner is the body of the ScopeScanner
// constructor.
//
// Java: SyntacticScopeCreator.ScopeScanner(AbstractCompiler, RedeclarationHandler, Scope, Set<Node>, boolean)
func initSyntacticScopeCreator_ScopeScanner(
	s *SyntacticScopeCreator_ScopeScanner,
	compiler AbstractCompiler,
	redeclarationHandler SyntacticScopeCreator_RedeclarationHandler,
	scope *Scope,
	changeRootSet map[rhino.Node]struct{},
	treatProvidesAsRedeclarations bool) {
	s.compiler = compiler
	s.redeclarationHandler = redeclarationHandler
	s.scope = scope
	s.changeRootSet = changeRootSet
	jutil.CheckState(changeRootSet == nil || scope.IsGlobal())
	s.treatProvidesAsRedeclarations = treatProvidesAsRedeclarations
	s.isFnBlock = scope.parent != nil && scope.parent.rootNode.IsFunction() && scope.rootNode.IsBlock()
}

func (s *SyntacticScopeCreator_ScopeScanner) setInputId(id *rhino.InputId) {
	s.inputId = id
	s.inputIdResolved = true
	s.input = nil
	s.inputResolved = false
}

func (s *SyntacticScopeCreator_ScopeScanner) getInput() *CompilerInput {
	if s.inputResolved {
		return s.input
	}
	return s.resolveInputSlow()
}

func (s *SyntacticScopeCreator_ScopeScanner) resolveInputSlow() *CompilerInput {
	if !s.inputIdResolved {
		if s.scope.parent != nil {
			if inp, ok := s.scope.inheritedInput(); ok {
				s.input = inp
				s.inputResolved = true
				return inp
			}
		}
		s.inputId = NodeUtil_GetInputId(s.scope.GetRootNode())
		s.inputIdResolved = true
	}
	if s.creator != nil && s.inputId != nil && s.inputId == s.creator.lastInputId {
		s.input = s.creator.lastInput
	} else {
		s.input = s.compiler.GetInput(s.inputId)
		if s.creator != nil && s.inputId != nil {
			s.creator.lastInputId = s.inputId
			s.creator.lastInput = s.input
		}
	}
	s.inputResolved = true
	return s.input
}

// Java: SyntacticScopeCreator.ScopeScanner.populate()
func (s *SyntacticScopeCreator_ScopeScanner) Populate() {
	n := s.scope.GetRootNode()
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		fn := ast.MustFunction(n)

		// Args: Declare function variables
		for param := range fn.Params().Params() {
			s.declareParam(s.scope, param)
		}

		// Bleed the function name into the scope, if it hasn't been declared in the outer scope
		// and the name isn't already in the scope via the param list.
		if name, ok := fn.Name(); ok && NodeUtil_IsFunctionExpression(n) {
			s.declareName(s.scope, name)
		}

		// Since we create a separate scope for body, stop scanning here
		return
	case rhino.Token_CLASS:
		// Bleed the class name into the scope, if it hasn't
		// been declared in the outer scope.
		if name, ok := ast.MustClass(n).Name(); ok && NodeUtil_IsClassExpression(n) {
			s.declareName(s.scope, name)
		}
		return
	case rhino.Token_ROOT, rhino.Token_SCRIPT:
		// n is the global scope
		if !s.scope.IsGlobal() {
			panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("%s", s.scope)))
		}
		s.ScanVars(n, s.scope, s.scope)
		return
	case rhino.Token_MODULE_BODY:
		s.ScanVars(n, s.scope, s.scope)
		return
	case rhino.Token_FOR, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_FOR_IN:
		// Go-only: performance — Java's scanVars(n, null, scope), which can
		// only declare the let or const of the loop's head: the body is a
		// block, with a scope of its own.
		if n.HasNoBlockScopeDecls() {
			return
		}
		if decl, ok := syntacticForLetOrConst(n); ok {
			s.declareNameDeclaration(s.scope, decl)
		}
		if len(s.scope.varList) == 0 {
			n.SetNoBlockScopeDecls()
		}
		return
	case rhino.Token_SWITCH_BODY:
		s.ScanVars(n, nil, s.scope)
		return
	case rhino.Token_BLOCK:
		if NodeUtil_IsFunctionBlock(n) || NodeUtil_IsClassStaticBlock(n) {
			if !syntacticFunctionBlockHasNoDecls(n) {
				s.ScanVars(n, s.scope, s.scope)
			}
			return
		}
		// Go-only: performance — Java's scanVars(n, null, scope), on the
		// children that may declare a block-scoped name only. A block that
		// has none is remembered, unless a child was scanned: a change below
		// that child would not clear the mark on this block.
		if n.HasNoBlockScopeDecls() {
			return
		}
		canCacheNoDecls := true
		// (The children are statements, or the catch of a try.)
		for c := n.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			if scanVarsTokenKind[uint8(c.GetToken())] == 1 && syntacticScopeChildHasBlockDecls(c) {
				canCacheNoDecls = false
				s.ScanVars(c, nil, s.scope)
			}
			c = next
		}
		if canCacheNoDecls && len(s.scope.varList) == 0 {
			n.SetNoBlockScopeDecls()
		}
		return
	case rhino.Token_COMPUTED_FIELD_DEF, rhino.Token_MEMBER_FIELD_DEF:
		// MEMBER_FIELD_DEF and COMPUTED_FIELD_DEF scopes only created to scope `this` and `super`
		// correctly
		return
	default:
		panic(jutil.NewRuntimeException("Illegal scope root: " + jutil.ToJavaString(n)))
	}
}

// DeclareLHS declares in sc the names that n binds: n is a var, let or
// const, a parameter list, a catch, an import, a destructuring declarator,
// a default value, a for-in or for-of loop, or an assignment.
//
// Go-only: the scanner calls the typed method for each shape directly.
//
// Java: SyntacticScopeCreator.ScopeScanner.declareLHS(Scope, Node)
func (s *SyntacticScopeCreator_ScopeScanner) DeclareLHS(sc *Scope, n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		s.declareNameDeclaration(sc, ast.MustNameDeclaration(n))
	case rhino.Token_PARAM_LIST:
		for param := range ast.MustParamList(n).Params() {
			s.declareParam(sc, param)
		}
	case rhino.Token_CATCH:
		s.declareCatchParam(sc, ast.MustCatch(n))
	case rhino.Token_IMPORT:
		s.declareImport(sc, ast.MustImport(n))
	case rhino.Token_DESTRUCTURING_LHS:
		s.declarePattern(sc, ast.MustDestructuringLhs(n).Pattern())
	case rhino.Token_DEFAULT_VALUE:
		s.declareTarget(sc, ast.MustDefaultValue(n).Target())
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		target := ast.MustEnhancedFor(n).Target()
		if decl, ok := target.AsNameDeclaration(); ok {
			s.declareNameDeclaration(sc, decl)
		} else if t, ok := target.AsTarget(); ok {
			s.declareTarget(sc, t)
		}
	default:
		assignment, ok := ast.AsAssignment(n)
		jutil.CheckArgumentMsg1(ok, "%s", n)
		s.declareTarget(sc, assignment.Target())
	}
}

// declareName declares the name name in sc.
func (s *SyntacticScopeCreator_ScopeScanner) declareName(sc *Scope, name ast.Name) {
	s.DeclareVar(sc, name.Node())
}

// declareNameDeclaration declares in sc the names of a var, let or const:
// `var a = 1, {b, c} = d;` declares a, b and c.
func (s *SyntacticScopeCreator_ScopeScanner) declareNameDeclaration(sc *Scope, decl ast.NameDeclaration) {
	// The next declarator is read first: the redeclaration handler may move
	// this one (Normalize turns a duplicate `var a = b` into `a = b`, where
	// the next sibling of a is b).
	for declarator := range ast.ReadAhead(decl.FirstDeclarator()) {
		if name, ok := declarator.AsName(); ok {
			s.declareName(sc, name)
		} else if lhs, ok := declarator.AsDestructuringLhs(); ok {
			s.declarePattern(sc, lhs.Pattern())
		}
	}
}

// declareParam declares in sc the names of a function parameter: a name or
// a pattern, either one with a default value (`a = 1`), or a rest (`...a`).
func (s *SyntacticScopeCreator_ScopeScanner) declareParam(sc *Scope, param ast.Param) {
	if name, ok := param.AsName(); ok {
		s.declareName(sc, name)
	} else if pattern, ok := param.AsPattern(); ok {
		s.declarePattern(sc, pattern)
	} else if withDefault, ok := param.AsDefaultValue(); ok {
		s.declareTarget(sc, withDefault.Target())
	} else if rest, ok := param.AsIterRest(); ok {
		s.declareTarget(sc, rest.Target())
	}
}

// declareCatchParam declares in sc the parameter of catch (a name or a
// pattern); `catch {}` has none.
func (s *SyntacticScopeCreator_ScopeScanner) declareCatchParam(sc *Scope, catch ast.Catch) {
	if param, ok := catch.Param(); ok {
		s.declareTarget(sc, param.AsTarget())
	}
}

// declareImport declares in sc the names an import binds, in order: its
// default, then the local names of its specifiers or its `* as name`.
func (s *SyntacticScopeCreator_ScopeScanner) declareImport(sc *Scope, imp ast.Import) {
	if name, ok := imp.Default(); ok {
		s.declareName(sc, name)
	}
	bindings, ok := imp.Bindings()
	if !ok {
		return
	}
	if star, ok := bindings.AsImportStar(); ok {
		s.DeclareVar(sc, star.Node())
	} else if specs, ok := bindings.AsImportSpecs(); ok {
		for spec := range specs.Specs() {
			s.declareName(sc, spec.Local())
		}
	}
}

// declareTarget declares in sc the names of a declaration's target: a name,
// or the names in a pattern.
func (s *SyntacticScopeCreator_ScopeScanner) declareTarget(sc *Scope, target ast.Target) {
	if name, ok := target.AsName(); ok {
		s.declareName(sc, name)
	} else if pattern, ok := target.AsPattern(); ok {
		s.declarePattern(sc, pattern)
	} else if cast, ok := target.AsCast(); ok {
		// Only in an assignment (see DeclareLHS); the cast holds a target.
		s.declareTarget(sc, ast.MustTarget(cast.Expr().Node()))
	} else {
		// A property (a.b, a[i]) declares nothing: DeclareVar rejects it.
		s.DeclareVar(sc, target.Node())
	}
}

// declarePattern declares in sc the names in a destructuring pattern, left
// to right.
func (s *SyntacticScopeCreator_ScopeScanner) declarePattern(sc *Scope, pattern ast.Pattern) {
	if array, ok := pattern.AsArrayPattern(); ok {
		// [a, b = 1, , ...c]
		for elem := range array.Elements() {
			if target, ok := elem.AsTarget(); ok {
				s.declareTarget(sc, target)
			} else if withDefault, ok := elem.AsDefaultValue(); ok {
				s.declareTarget(sc, withDefault.Target())
			} else if rest, ok := elem.AsIterRest(); ok {
				s.declareTarget(sc, rest.Target())
			} // else a hole, which declares nothing
		}
	} else if object, ok := pattern.AsObjectPattern(); ok {
		// {a, b: c, d = 1, [e]: f, ...g}; {a} is {a: a}, with a NAME value.
		for prop := range object.Props() {
			if key, ok := prop.AsStringKey(); ok {
				s.declarePropValue(sc, key.Value())
			} else if computed, ok := prop.AsComputedProp(); ok {
				if value, ok := computed.Value(); ok {
					s.declarePropValue(sc, value)
				}
			} else if rest, ok := prop.AsObjectRest(); ok {
				s.declareTarget(sc, rest.Target())
			}
		}
	}
}

// declarePropValue declares in sc the names of the value of a property in
// an object pattern: a target, with or without a default value.
func (s *SyntacticScopeCreator_ScopeScanner) declarePropValue(sc *Scope, value ast.PropValue) {
	if withDefault, ok := value.AsDefaultValue(); ok {
		s.declareTarget(sc, withDefault.Target())
	} else if target, ok := value.AsTarget(); ok {
		s.declareTarget(sc, target)
	}
}

// declareFunctionDeclaration declares in sc the name of fn if fn is a
// function declaration; a function expression declares nothing outside
// itself. The function's children belong to its own scopes.
func (s *SyntacticScopeCreator_ScopeScanner) declareFunctionDeclaration(sc *Scope, fn ast.Function) {
	if NodeUtil_IsFunctionExpression(fn.Node()) {
		return
	}
	// A declaration without a name is invalid, but allow it so the checks
	// can catch it.
	if name, ok := fn.Name(); ok {
		s.declareName(sc, name)
	}
}

// declareClassDeclaration declares in sc the name of cls if cls is a class
// declaration; a class expression declares nothing outside itself. The
// class's children belong to its own scopes.
func (s *SyntacticScopeCreator_ScopeScanner) declareClassDeclaration(sc *Scope, cls ast.Class) {
	if NodeUtil_IsClassExpression(cls.Node()) {
		return
	}
	// A declaration without a name is invalid, but allow it so the checks
	// can catch it.
	if name, ok := cls.Name(); ok {
		s.declareName(sc, name)
	}
}

// scanVarsTokenKind tells ScanVars which children to look at, by token: 1
// for the declarations and the statements that may contain one, 2 for an
// expression statement (at the top of a script, goog.provide('a.b') declares
// a), 0 for the nodes that declare nothing.
var scanVarsTokenKind = [256]uint8{
	rhino.Token_VAR:          1,
	rhino.Token_LET:          1,
	rhino.Token_CONST:        1,
	rhino.Token_IMPORT:       1,
	rhino.Token_EXPORT:       1,
	rhino.Token_FUNCTION:     1,
	rhino.Token_CLASS:        1,
	rhino.Token_CATCH:        1,
	rhino.Token_SCRIPT:       1,
	rhino.Token_MODULE_BODY:  1,
	rhino.Token_ROOT:         1,
	rhino.Token_BLOCK:        1,
	rhino.Token_SWITCH_BODY:  1,
	rhino.Token_FOR:          1,
	rhino.Token_FOR_IN:       1,
	rhino.Token_FOR_OF:       1,
	rhino.Token_FOR_AWAIT_OF: 1,
	rhino.Token_DO:           1,
	rhino.Token_WHILE:        1,
	rhino.Token_WITH:         1,
	rhino.Token_IF:           1,
	rhino.Token_LABEL:        1,
	rhino.Token_TRY:          1,
	rhino.Token_SWITCH:       1,
	rhino.Token_CASE:         1,
	rhino.Token_DEFAULT_CASE: 1,
	rhino.Token_EXPR_RESULT:  2,
}

// scanVarsBlockCreatesScope reports whether the BLOCK n has a scope of its
// own (NodeUtil.createsBlockScope): not a synthetic block in a block or a
// script, nor the body of a case or a catch (a catch shares the scope of
// the block that holds it in its try).
func scanVarsBlockCreatesScope(n rhino.Node) bool {
	parent := n.GetParent()
	if parent.IsNil() {
		return false
	}
	pt := parent.GetToken()
	if (pt == rhino.Token_BLOCK || pt == rhino.Token_SCRIPT) && n.IsSyntheticBlock() {
		return false
	}
	return pt != rhino.Token_CASE && pt != rhino.Token_DEFAULT_CASE && pt != rhino.Token_CATCH
}

// ScanVars scans and gather variables declarations under a Node.
//
// n is the node; hoistScope the scope that is the hoist target for vars, if
// we are scanning for vars; blockScope the scope that is the hoist target
// for block-level declarations, if we are scanning for block level
// declarations.
//
// Java: SyntacticScopeCreator.ScopeScanner.scanVars(Node, Scope, Scope)
func (s *SyntacticScopeCreator_ScopeScanner) ScanVars(n rhino.Node, hoistScope *Scope, blockScope *Scope) {
	switch n.GetToken() {
	case rhino.Token_VAR:
		if hoistScope != nil {
			s.declareNameDeclaration(hoistScope, ast.MustNameDeclaration(n))
		}
		return
	case rhino.Token_LET, rhino.Token_CONST:
		// Only declare when scope is the current lexical scope
		if blockScope != nil {
			s.declareNameDeclaration(blockScope, ast.MustNameDeclaration(n))
		}
		return
	case rhino.Token_IMPORT:
		s.declareImport(hoistScope, ast.MustImport(n))
		return
	case rhino.Token_EXPORT:
		// The body of an EXPORT can be a declaration, in the case of
		// export var/let/const/function/class name ...
		s.ScanVars(ast.MustExport(n).Body().Node(), hoistScope, blockScope)
		return
	case rhino.Token_FUNCTION:
		if blockScope != nil {
			s.declareFunctionDeclaration(blockScope, ast.MustFunction(n))
		}
		return // should not examine function's children
	case rhino.Token_CLASS:
		if blockScope != nil {
			s.declareClassDeclaration(blockScope, ast.MustClass(n))
		}
		return // should not examine class's children
	case rhino.Token_CATCH:
		catch := ast.MustCatch(n)
		if blockScope != nil {
			s.declareCatchParam(blockScope, catch)
		}
		// A new scope is not created for this BLOCK because there is a scope
		// created for the BLOCK above the CATCH
		s.ScanVars(catch.Body().Node(), hoistScope, blockScope)
		return // only one child to scan
	case rhino.Token_SCRIPT:
		if s.changeRootSet != nil {
			if _, ok := s.changeRootSet[n]; !ok {
				// If there is a changeRootSet configured, that means
				// a partial update is being done and we should skip
				// any SCRIPT that aren't being asked for.
				return
			}
		}
		s.setInputId(n.GetInputId())
		// Go-only: performance — the loop over children below, with the
		// statements of a script that may declare an implicit goog namespace
		// (Java: the EXPR_RESULT case).
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			ctok := child.GetToken()
			if scanVarsTokenKind[uint8(ctok)] != 0 {
				switch ctok {
				case rhino.Token_EXPR_RESULT:
					if stmt := ast.MustExprStmt(child); stmt.Expr().Token() == rhino.Token_CALL {
						s.declareImplicitGoogNamespace(hoistScope, stmt)
					}
				case rhino.Token_VAR:
					if hoistScope != nil {
						s.declareNameDeclaration(hoistScope, ast.MustNameDeclaration(child))
					}
				case rhino.Token_LET, rhino.Token_CONST:
					if blockScope != nil {
						s.declareNameDeclaration(blockScope, ast.MustNameDeclaration(child))
					}
				case rhino.Token_FUNCTION:
					if blockScope != nil {
						s.declareFunctionDeclaration(blockScope, ast.MustFunction(child))
					}
				case rhino.Token_CLASS:
					if blockScope != nil {
						s.declareClassDeclaration(blockScope, ast.MustClass(child))
					}
				default:
					s.ScanVars(child, hoistScope, blockScope)
				}
			}
			child = next
		}
		return
	case rhino.Token_MODULE_BODY:
		// Module bodies are not part of global scope, but may declare an implicit goog namespace.
		if hoistScope.IsGlobal() {
			first, ok := ast.MustModuleBody(n).FirstStatement()
			if ok && syntacticScopeCreator_ScopeScanner_isLegacyGoogModule(first) {
				s.DeclareImplicitGoogNamespaceFromCall(hoistScope, first.Node())
			}
			return
		}
	case rhino.Token_EXPR_RESULT:
		if n.GetParent().IsScript() {
			s.declareImplicitGoogNamespace(hoistScope, ast.MustExprStmt(n))
		}
		return
	case rhino.Token_BLOCK:
		if blockScope != nil {
			if n != blockScope.GetRootNode() && scanVarsBlockCreatesScope(n) {
				if hoistScope == nil {
					// We only enter new blocks when scanning for hoisted vars
					return
				}
				blockScope = nil
			}
		} else if hoistScope == nil && scanVarsBlockCreatesScope(n) {
			// We only enter new blocks when scanning for hoisted vars
			return
		}
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_SWITCH_BODY:
		if blockScope == nil || n != blockScope.GetRootNode() {
			if hoistScope == nil {
				// We only enter new blocks when scanning for hoisted vars
				return
			}
			blockScope = nil
		}
	case rhino.Token_ROOT, rhino.Token_DO, rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_IF,
		rhino.Token_LABEL, rhino.Token_TRY, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
	default:
		return
	}

	// Variables can only occur in statement-level nodes, so
	// we only need to traverse children in a couple special cases.
	//
	// Go-only: performance — the children that can declare nothing are
	// skipped by token (scanVarsTokenKind), and the cases below that
	// declare, or that would return right away, are handled without a call.
	for child := n.GetFirstChild(); !child.IsNil(); {
		next := child.GetNext()
		ctok := child.GetToken()
		if scanVarsTokenKind[uint8(ctok)] == 1 {
			switch ctok {
			case rhino.Token_VAR:
				if hoistScope != nil {
					s.declareNameDeclaration(hoistScope, ast.MustNameDeclaration(child))
				}
			case rhino.Token_LET, rhino.Token_CONST:
				if blockScope != nil {
					s.declareNameDeclaration(blockScope, ast.MustNameDeclaration(child))
				}
			case rhino.Token_FUNCTION:
				if blockScope != nil {
					s.declareFunctionDeclaration(blockScope, ast.MustFunction(child))
				}
			case rhino.Token_CLASS:
				if blockScope != nil {
					s.declareClassDeclaration(blockScope, ast.MustClass(child))
				}
			case rhino.Token_BLOCK:
				if hoistScope != nil || !scanVarsBlockCreatesScope(child) {
					s.ScanVars(child, hoistScope, blockScope)
				}
			case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_SWITCH_BODY:
				if hoistScope != nil {
					s.ScanVars(child, hoistScope, nil)
				}
			default:
				s.ScanVars(child, hoistScope, blockScope)
			}
		}
		child = next
	}
}

// declareImplicitGoogNamespace declares the implicit namespace of stmt, a
// statement of a script, if it is goog.provide('a.b.c') or a bundled legacy
// goog.module:
//
//	goog.loadModule(function(exports) {
//	  goog.module('a.b.c');
//	  goog.module.declareLegacyNamespace();
//	  ...
func (s *SyntacticScopeCreator_ScopeScanner) declareImplicitGoogNamespace(hoistScope *Scope, stmt ast.ExprStmt) {
	if NodeUtil_IsGoogProvideCall(stmt.Node()) {
		s.DeclareImplicitGoogNamespaceFromCall(hoistScope.GetGlobalScope(), stmt.Node())
		return
	}
	call, ok := stmt.Expr().AsCall()
	if !ok || !NodeUtil_IsBundledGoogModuleCall(call.Node()) {
		return
	}
	arg, _ := call.OnlyArg() // a bundled goog.module call has one argument
	fn, ok := arg.AsFunction()
	if !ok {
		return
	}
	body, ok := fn.Body().AsBlock()
	if !ok {
		return
	}
	moduleCall, ok := body.FirstStatement()
	if ok && syntacticScopeCreator_ScopeScanner_isLegacyGoogModule(moduleCall) {
		s.DeclareImplicitGoogNamespaceFromCall(hoistScope.GetGlobalScope(), moduleCall.Node())
	}
}

// DeclareVar declares a variable.
//
// sc is the scope to declare the variable in, n the node corresponding to
// the variable name.
//
// Java: SyntacticScopeCreator.ScopeScanner.declareVar(Scope, Node)
func (s *SyntacticScopeCreator_ScopeScanner) DeclareVar(sc *Scope, n rhino.Node) {
	if !n.IsName() && !n.IsImportStar() {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Invalid node for declareVar: %s", n)))
	}

	id := n.GetStringID()
	// Because of how we scan the variables, it is possible to encounter
	// the same var declared name node twice. Bail out in this case.
	// TODO(johnlenz): Hash lookups are not free and building scopes are already expensive.
	// Restructure the scope building to avoid this check.
	v := sc.getOwnSlotByID(id)
	if v != nil {
		if v.GetNode() == n {
			return
		} else if v.IsImplicitGoogNamespace() {
			// this is replacing an implicit provide. treat this as the actual declaration.
			sc.Undeclare(v)
			v = nil
		}
	}

	var parentVar *Var
	if v == nil && s.isFnBlock {
		parentVar = sc.parent.lookupOwnVarByID(id)
	}

	input := s.getInput()
	if v != nil ||
		(parentVar != nil && parentVar.IsParam()) ||
		(id == implicitArgumentsID &&
			(sc.rootNode.IsFunction() || s.isFnBlock)) {
		s.redeclarationHandler.OnRedeclaration(sc, n.GetString(), n, input)
	} else {
		if id <= rhino.RhinoStringPool_EmptyID {
			panic(jutil.NewIllegalArgumentException())
		}
		if parentVar != nil && !sc.isBleedingFunctionNameID(id) {
			panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Illegal shadow: %s", n)))
		}
		var_ := initVar(s.creator.vars.New(), id, n, sc, len(sc.varList), input, rhino.NilNode)
		sc.addVar(var_)
	}
}

// DeclareImplicitGoogNamespaceFromCall declares the implicit namespace 'a'
// from 'goog.provide('a.b.c');
//
// Explicit syntactical definitions like 'var a = {};' supersede this
// definition if present.
//
// exprCall is a goog.module or goog.provide call (the statement).
//
// Java: SyntacticScopeCreator.ScopeScanner.declareImplicitGoogNamespaceFromCall(Scope, Node)
func (s *SyntacticScopeCreator_ScopeScanner) DeclareImplicitGoogNamespaceFromCall(sc *Scope, exprCall rhino.Node) {
	call := ast.MustCall(ast.MustExprStmt(exprCall).Expr().Node())
	var namespaceNode ast.StringLit
	arg, ok := call.FirstArg()
	if ok {
		namespaceNode, ok = arg.AsStringLit()
	}
	if !ok {
		// invalid call missing an argument, we warn elsewhere
		return
	}
	namespace := namespaceNode.Value()
	root := NodeUtil_GetRootOfQualifiedNameString(namespace)

	if root == "" {
		return
	}
	// Conditionally report a redeclaration if the root name has already been declared as a normal
	// variable (normal meaning a var/function/class/etc. declaration).
	if s.treatProvidesAsRedeclarations {
		existing := sc.getOwnSlotV(root)
		if existing != nil && !existing.IsImplicitGoogNamespace() {
			s.redeclarationHandler.OnRedeclaration(
				sc, root, existing.GetNode(), s.getInput())
		}
	}
	sc.DeclareImplicitGoogNamespaceIfAbsent(root, namespaceNode.Node())
}

// isLegacyGoogModule reports whether stmt is goog.module('a.b.c') followed
// by goog.module.declareLegacyNamespace().
//
// Java: SyntacticScopeCreator.ScopeScanner.isLegacyGoogModule(Node)
func syntacticScopeCreator_ScopeScanner_isLegacyGoogModule(stmt ast.Statement) bool {
	if !NodeUtil_IsGoogModuleCall(stmt.Node()) {
		return false
	}
	next, ok := stmt.Next()
	return ok && NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(next.Node())
}

// SyntacticScopeCreator_RedeclarationHandler is the interface for
// injectable duplicate handling.
type SyntacticScopeCreator_RedeclarationHandler interface {
	// Java: SyntacticScopeCreator.RedeclarationHandler.onRedeclaration(Scope, String, Node, CompilerInput)
	OnRedeclaration(s *Scope, name string, n rhino.Node, input *CompilerInput)
}

// SyntacticScopeCreator_DefaultRedeclarationHandler is the default handler
// for duplicate declarations.
type SyntacticScopeCreator_DefaultRedeclarationHandler struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: SyntacticScopeCreator.DefaultRedeclarationHandler()
func NewSyntacticScopeCreator_DefaultRedeclarationHandler() *SyntacticScopeCreator_DefaultRedeclarationHandler {
	return &SyntacticScopeCreator_DefaultRedeclarationHandler{}
}

// Java: SyntacticScopeCreator.DefaultRedeclarationHandler.onRedeclaration(Scope, String, Node, CompilerInput)
func (h *SyntacticScopeCreator_DefaultRedeclarationHandler) OnRedeclaration(
	s *Scope, name string, n rhino.Node, input *CompilerInput) {
}

// MemoizedScopeCreator memoizes a scope creator.
//
// This allows you to make multiple passes, without worrying about the
// expense of generating Scope objects over and over again.
//
// On the other hand, you also have to be more aware of what your passes are
// doing. Scopes are memoized stupidly, so if the underlying tree changes, the
// scope may be out of sync.
type MemoizedScopeCreator struct {
	untypedScopes     map[rhino.Node]*Scope
	syntacticDelegate *SyntacticScopeCreator
}

// NewMemoizedScopeCreator: delegate is the real source of Scope objects.
//
// Java: MemoizedScopeCreator(ScopeCreator)
func NewMemoizedScopeCreator(delegate *SyntacticScopeCreator) *MemoizedScopeCreator {
	return &MemoizedScopeCreator{
		untypedScopes:     make(map[rhino.Node]*Scope),
		syntacticDelegate: jutil.CheckNotNull(delegate),
	}
}

func (m *MemoizedScopeCreator) createUntypedScope(n rhino.Node, parent *Scope) *Scope {
	if scope := m.untypedScopes[n]; scope != nil {
		jutil.CheckState(parent == scope.parent)
		return scope
	}
	scope := m.syntacticDelegate.CreateScopeOfScope(n, parent)
	m.untypedScopes[n] = scope
	return scope
}

// CreateScopeOfScope returns the memoized concrete *Scope for n without
// boxing through AbstractScope.
func (m *MemoizedScopeCreator) CreateScopeOfScope(n rhino.Node, parent *Scope) *Scope {
	return m.createUntypedScope(n, parent)
}

// Java: MemoizedScopeCreator.createScope(Node, AbstractScope<?, ?>)
func (m *MemoizedScopeCreator) CreateScope(n rhino.Node, parent AbstractScope) AbstractScope {
	var p *Scope
	if parent != nil {
		p = parent.(*Scope)
	}
	return m.createUntypedScope(n, p)
}

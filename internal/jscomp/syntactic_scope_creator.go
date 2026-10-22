// Port of com/google/javascript/jscomp/SyntacticScopeCreator.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	arena                         *scopeArena
	poolable                      bool
}

func (c *SyntacticScopeCreator) EnableArenaPool() {
	if c != nil {
		c.poolable = true
	}
}

func (c *SyntacticScopeCreator) ReleaseArena() {
	if c == nil {
		return
	}
	if c.arena != nil {
		c.arena.pooled = true
		c.arena.release()
		c.arena = nil
	}
	c.lastInputId = nil
	c.lastInput = nil
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
func (c *SyntacticScopeCreator) CreateScope(n *rhino.Node, parent AbstractScope) AbstractScope {
	// Go-only: language-adapter — the cast (Scope) parent, which lets null
	// through.
	var p *Scope
	if parent != nil {
		p = parent.(*Scope)
	}
	return c.CreateScopeOfScope(n, p)
}

// Java: SyntacticScopeCreator.createScope(Node, Scope)
func (c *SyntacticScopeCreator) CreateScopeOfScope(n *rhino.Node, parent *Scope) *Scope {
	if c.arena == nil {
		if c.poolable {
			c.arena = acquireScopeArena()
		} else {
			c.arena = &scopeArena{}
		}
	}
	var scope *Scope
	if parent == nil {
		scope = newScopeInArena(n, c.arena)
	} else {
		scope = newScopeWithParentInArena(parent, n, c.arena)
	}
	// Go-only: performance — the scanner, which nothing keeps, is a local
	// variable rather than a heap object (one per scope created).
	var scanner SyntacticScopeCreator_ScopeScanner
	initSyntacticScopeCreator_ScopeScanner(&scanner, c.compiler, c.redeclarationHandler, scope, nil, c.treatProvidesAsRedeclarations)
	scanner.creator = c
	scanner.Populate()
	scope.scanDone = true
	return scope
}

func syntacticScopeChildHasBlockDecls(c *rhino.Node) bool {
	switch c.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CATCH, rhino.Token_EXPORT:
		return true
	case rhino.Token_FUNCTION:
		return !NodeUtil_IsFunctionExpression(c) && c.GetFirstChild().GetString() != ""
	case rhino.Token_CLASS:
		return !NodeUtil_IsClassExpression(c) && c.GetFirstChild().GetString() != ""
	case rhino.Token_BLOCK:
		if !c.IsSyntheticBlock() {
			return false
		}
		for gc := c.GetFirstChild(); gc != nil; gc = gc.GetNext() {
			if syntacticScopeChildHasBlockDecls(gc) {
				return true
			}
		}
		return false
	case rhino.Token_LABEL, rhino.Token_IF, rhino.Token_DO, rhino.Token_WHILE, rhino.Token_WITH:
		for gc := c.GetFirstChild(); gc != nil; gc = gc.GetNext() {
			if scanVarsTokenKind[uint8(gc.GetToken())] == 1 && syntacticScopeChildHasBlockDecls(gc) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func syntacticScopeHasBlockDecls(n *rhino.Node) bool {
	if n.HasNoBlockScopeDecls() {
		return false
	}
	if NodeUtil_IsAnyFor(n) {
		first := n.GetFirstChild()
		if first != nil && (first.IsLet() || first.IsConst()) {
			return true
		}
		if last := n.GetLastChild(); last != nil && !last.IsBlock() && syntacticScopeChildHasBlockDecls(last) {
			return true
		}
		n.SetNoBlockScopeDecls()
		return false
	}
	if n.IsBlock() {
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			if scanVarsTokenKind[uint8(c.GetToken())] == 1 && syntacticScopeChildHasBlockDecls(c) {
				return true
			}
		}
		n.SetNoBlockScopeDecls()
		return false
	}
	return true
}

func syntacticFunctionBlockHasNoDecls(n *rhino.Node) bool {
	if n.HasNoBlockScopeDecls() {
		return true
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if scanVarsTokenKind[uint8(c.GetToken())] == 1 {
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
	parentBloom                   uint64

	// Will be null, when a detached node is traversed.
	inputId         *rhino.InputId
	input           *CompilerInput
	inputIdResolved bool
	inputResolved   bool
	changeRootSet   *jutil.LinkedHashSet[*rhino.Node]

	smallVars [8]*Var
	smallLen  int
}

// Java: SyntacticScopeCreator.ScopeScanner(AbstractCompiler, RedeclarationHandler, Scope, Set<Node>, boolean)
func newSyntacticScopeCreator_ScopeScanner(
	compiler AbstractCompiler,
	redeclarationHandler SyntacticScopeCreator_RedeclarationHandler,
	scope *Scope,
	changeRootSet *jutil.LinkedHashSet[*rhino.Node],
	treatProvidesAsRedeclarations bool) *SyntacticScopeCreator_ScopeScanner {
	s := &SyntacticScopeCreator_ScopeScanner{}
	initSyntacticScopeCreator_ScopeScanner(s, compiler, redeclarationHandler, scope, changeRootSet, treatProvidesAsRedeclarations)
	return s
}

// initSyntacticScopeCreator_ScopeScanner is the body of the ScopeScanner
// constructor.
func initSyntacticScopeCreator_ScopeScanner(
	s *SyntacticScopeCreator_ScopeScanner,
	compiler AbstractCompiler,
	redeclarationHandler SyntacticScopeCreator_RedeclarationHandler,
	scope *Scope,
	changeRootSet *jutil.LinkedHashSet[*rhino.Node],
	treatProvidesAsRedeclarations bool) {
	s.compiler = compiler
	s.redeclarationHandler = redeclarationHandler
	s.scope = scope
	s.changeRootSet = changeRootSet
	jutil.CheckState(changeRootSet == nil || scope.IsGlobal())
	s.treatProvidesAsRedeclarations = treatProvidesAsRedeclarations
	s.isFnBlock = scope.parent != nil && scope.parent.rootNode.IsFunction() && scope.rootNode.IsBlock()
	if s.isFnBlock {
		s.parentBloom = scope.parent.varBloom
	} else {
		s.parentBloom = 0
	}
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

func (s *SyntacticScopeCreator_ScopeScanner) flushSmallVarsExact() {
	sc := s.scope
	k := s.smallLen
	s.smallLen = 0
	if sc.arena != nil {
		sc.varList = sc.arena.allocVarList(k, k)
	} else {
		sc.varList = make([]*Var, k)
	}
	copy(sc.varList, s.smallVars[:k])
	sc.varCount = int32(k)
	sc.varsInitialized = true
}

func (s *SyntacticScopeCreator_ScopeScanner) flushSmallVarsForMore() {
	sc := s.scope
	k := s.smallLen
	s.smallLen = 0
	n := largeScopeVarSlotCount(sc)
	if n <= k {
		n = k + 4
	}
	sc.initLargeVars(n)
	sc.varsInitialized = true
	mask := uint32(len(sc.varHash) - 1)
	for i := 0; i < k; i++ {
		v := s.smallVars[i]
		sc.varList = append(sc.varList, v)
		slot := scopeNameHash(v.name) & mask
		for sc.varHash[slot] != 0 {
			slot = (slot + 1) & mask
		}
		sc.varHash[slot] = int32(i + 1)
	}
	sc.varCount = int32(k)
}

// Java: SyntacticScopeCreator.ScopeScanner.populate()
func (s *SyntacticScopeCreator_ScopeScanner) Populate() {
	s.populate()
	if s.smallLen > 0 {
		s.flushSmallVarsExact()
	}
}

func (s *SyntacticScopeCreator_ScopeScanner) populate() {
	n := s.scope.GetRootNode()
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		fnNameNode := n.GetFirstChild()
		args := fnNameNode.GetNext()

		// Args: Declare function variables
		jutil.CheckState(args.IsParamList())
		for a := args.GetFirstChild(); a != nil; a = a.GetNext() {
			if a.IsName() {
				s.DeclareVar(s.scope, a)
			} else {
				s.declareLhsNodesHelper(s.scope, a)
			}
		}

		// Bleed the function name into the scope, if it hasn't been declared in the outer scope
		// and the name isn't already in the scope via the param list.
		fnName := fnNameNode.GetString()
		if fnName != "" && NodeUtil_IsFunctionExpression(n) {
			s.DeclareVar(s.scope, fnNameNode)
		}

		// Since we create a separate scope for body, stop scanning here
		return
	case rhino.Token_CLASS:
		classNameNode := n.GetFirstChild()
		// Bleed the class name into the scope, if it hasn't
		// been declared in the outer scope.
		if !classNameNode.IsEmpty() && NodeUtil_IsClassExpression(n) {
			s.DeclareVar(s.scope, classNameNode)
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
		if n.HasNoBlockScopeDecls() {
			return
		}
		first := n.GetFirstChild()
		if first != nil && (first.IsLet() || first.IsConst()) {
			s.DeclareLHS(s.scope, first)
		}
		last := n.GetLastChild()
		if last != nil && !last.IsBlock() && syntacticScopeChildHasBlockDecls(last) {
			s.ScanVars(last, nil, s.scope)
		}
		if s.smallLen == 0 && s.scope.varCount == 0 && (last == nil || last.IsBlock()) {
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
		} else {
			if n.HasNoBlockScopeDecls() {
				return
			}
			canCacheNoDecls := true
			for c := n.GetFirstChild(); c != nil; {
				next := c.GetNext()
				if scanVarsTokenKind[uint8(c.GetToken())] == 1 && syntacticScopeChildHasBlockDecls(c) {
					canCacheNoDecls = false
					s.ScanVars(c, nil, s.scope)
				}
				c = next
			}
			if canCacheNoDecls && s.smallLen == 0 && s.scope.varCount == 0 {
				n.SetNoBlockScopeDecls()
			}
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

// Java: SyntacticScopeCreator.ScopeScanner.declareLHS(Scope, Node)
func (s *SyntacticScopeCreator_ScopeScanner) DeclareLHS(sc *Scope, n *rhino.Node) {
	if n.HasOneChild() && n.GetFirstChild().IsName() {
		// NAME is most common and trivial case.  This code is hot so it is worth special casing.
		// to avoid extra GC and cpu cycles.
		s.DeclareVar(sc, n.GetFirstChild())
	} else {
		jutil.CheckArgumentMsg1(
			NodeUtil_IsNameDeclaration(n) ||
				n.IsParamList() ||
				NodeUtil_IsAssignmentOp(n) ||
				n.IsCatch() ||
				n.IsDestructuringLhs() ||
				n.IsDefaultValue() ||
				n.IsImport() ||
				NodeUtil_IsEnhancedFor(n),
			"%s",
			n)
		s.declareLhsNodesHelper(sc, n)
	}
}

func (s *SyntacticScopeCreator_ScopeScanner) declareLhsNodesHelper(sc *Scope, n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_IMPORT:
		s.declareLhsNodesHelper(sc, n.GetFirstChild())
		s.declareLhsNodesHelper(sc, n.GetSecondChild())
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET, rhino.Token_OBJECT_PATTERN, rhino.Token_ARRAY_PATTERN,
		rhino.Token_PARAM_LIST, rhino.Token_IMPORT_SPECS:
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			if child.IsName() {
				s.DeclareVar(sc, child)
			} else {
				s.declareLhsNodesHelper(sc, child)
			}
		}
	case rhino.Token_DESTRUCTURING_LHS, rhino.Token_DEFAULT_VALUE, rhino.Token_CATCH, rhino.Token_ITER_REST,
		rhino.Token_OBJECT_REST, rhino.Token_CAST:
		s.declareLhsNodesHelper(sc, n.GetFirstChild())
	case rhino.Token_IMPORT_SPEC, rhino.Token_COMPUTED_PROP, rhino.Token_STRING_KEY:
		s.declareLhsNodesHelper(sc, n.GetLastChild())
	case rhino.Token_NAME, rhino.Token_IMPORT_STAR, rhino.Token_GETPROP, rhino.Token_GETELEM:
		s.DeclareVar(sc, n)
	case rhino.Token_EMPTY:
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		s.declareLhsNodesHelper(sc, n.GetFirstChild())
	default:
		if NodeUtil_IsAssignmentOp(n) {
			s.declareLhsNodesHelper(sc, n.GetFirstChild())
		} else {
			panic(jutil.NewIllegalStateExceptionMsg("Invalid node in lhs: " + n.String()))
		}
	}
}

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

func scanVarsBlockCreatesScope(n *rhino.Node) bool {
	parent := n.GetParent()
	if parent == nil {
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
func (s *SyntacticScopeCreator_ScopeScanner) ScanVars(n *rhino.Node, hoistScope *Scope, blockScope *Scope) {
	switch n.GetToken() {
	case rhino.Token_VAR:
		if hoistScope != nil {
			s.DeclareLHS(hoistScope, n)
		}
		return
	case rhino.Token_LET, rhino.Token_CONST:
		// Only declare when scope is the current lexical scope
		if blockScope != nil {
			s.DeclareLHS(blockScope, n)
		}
		return
	case rhino.Token_IMPORT:
		s.DeclareLHS(hoistScope, n)
		return
	case rhino.Token_EXPORT:
		// The first child of an EXPORT can be a declaration, in the case of
		// export var/let/const/function/class name ...
		s.ScanVars(n.GetFirstChild(), hoistScope, blockScope)
		return
	case rhino.Token_FUNCTION:
		if NodeUtil_IsFunctionExpression(n) || blockScope == nil {
			return
		}

		fnName := n.GetFirstChild().GetString()
		if fnName == "" {
			// This is invalid, but allow it so the checks can catch it.
			return
		}
		s.DeclareVar(blockScope, n.GetFirstChild())
		return // should not examine function's children
	case rhino.Token_CLASS:
		if NodeUtil_IsClassExpression(n) || blockScope == nil {
			return
		}
		className := n.GetFirstChild().GetString()
		if className == "" {
			// This is invalid, but allow it so the checks can catch it.
			return
		}
		s.DeclareVar(blockScope, n.GetFirstChild())
		return // should not examine class's children
	case rhino.Token_CATCH:
		jutil.CheckStateMsg1(n.HasTwoChildren(), "%s", n)
		// the first child is the catch var and the second child
		// is the code block
		if blockScope != nil {
			s.DeclareLHS(blockScope, n)
		}
		// A new scope is not created for this BLOCK because there is a scope
		// created for the BLOCK above the CATCH
		block := n.GetSecondChild()
		s.ScanVars(block, hoistScope, blockScope)
		return // only one child to scan
	case rhino.Token_SCRIPT:
		if s.changeRootSet != nil && !s.changeRootSet.Contains(n) {
			// If there is a changeRootSet configured, that means
			// a partial update is being done and we should skip
			// any SCRIPT that aren't being asked for.
			return
		}
		s.setInputId(n.GetInputId())
		for child := n.GetFirstChild(); child != nil; {
			next := child.GetNext()
			if scanVarsTokenKind[uint8(child.GetToken())] != 0 {
				s.ScanVars(child, hoistScope, blockScope)
			}
			child = next
		}
		return
	case rhino.Token_MODULE_BODY:
		// Module bodies are not part of global scope, but may declare an implicit goog namespace.
		if hoistScope.IsGlobal() {
			expr := n.GetFirstChild()
			if expr != nil && syntacticScopeCreator_ScopeScanner_isLegacyGoogModule(expr) {
				s.DeclareImplicitGoogNamespaceFromCall(hoistScope, expr)
			}
			return
		}
	case rhino.Token_EXPR_RESULT:
		if !n.GetParent().IsScript() {
			return
		}
		if NodeUtil_IsGoogProvideCall(n) {
			s.DeclareImplicitGoogNamespaceFromCall(hoistScope.GetGlobalScope(), n)
		} else if NodeUtil_IsBundledGoogModuleCall(n.GetFirstChild()) &&
			n.GetFirstChild().GetSecondChild().IsFunction() {
			// e.g.
			// goog.loadModule(function(exports) {
			//   goog.module('foo.bar');
			fn := n.GetFirstChild().GetSecondChild()
			moduleCall := NodeUtil_GetFunctionBody(fn).GetFirstChild()
			if syntacticScopeCreator_ScopeScanner_isLegacyGoogModule(moduleCall) {
				s.DeclareImplicitGoogNamespaceFromCall(hoistScope.GetGlobalScope(), moduleCall)
			}
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

	for child := n.GetFirstChild(); child != nil; {
		next := child.GetNext()
		if scanVarsTokenKind[uint8(child.GetToken())] == 1 {
			switch child.GetToken() {
			case rhino.Token_VAR:
				if hoistScope != nil {
					if first := child.GetFirstChild(); first != nil && first.GetNext() == nil && first.IsName() {
						s.DeclareVar(hoistScope, first)
					} else {
						s.DeclareLHS(hoistScope, child)
					}
				}
			case rhino.Token_LET, rhino.Token_CONST:
				if blockScope != nil {
					if first := child.GetFirstChild(); first != nil && first.GetNext() == nil && first.IsName() {
						s.DeclareVar(blockScope, first)
					} else {
						s.DeclareLHS(blockScope, child)
					}
				}
			case rhino.Token_FUNCTION:
				if blockScope != nil && !NodeUtil_IsFunctionExpression(child) {
					if fnNameNode := child.GetFirstChild(); fnNameNode.GetString() != "" {
						s.DeclareVar(blockScope, fnNameNode)
					}
				}
			case rhino.Token_CLASS:
				if blockScope != nil && !NodeUtil_IsClassExpression(child) {
					if classNameNode := child.GetFirstChild(); classNameNode.GetString() != "" {
						s.DeclareVar(blockScope, classNameNode)
					}
				}
			default:
				s.ScanVars(child, hoistScope, blockScope)
			}
		}
		child = next
	}
}

// DeclareVar declares a variable.
//
// sc is the scope to declare the variable in, n the node corresponding to
// the variable name.
//
// Java: SyntacticScopeCreator.ScopeScanner.declareVar(Scope, Node)
func (s *SyntacticScopeCreator_ScopeScanner) DeclareVar(sc *Scope, n *rhino.Node) {
	if !n.IsName() && !n.IsImportStar() {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Invalid node for declareVar: %s", n)))
	}

	name := n.GetString()
	// Because of how we scan the variables, it is possible to encounter
	// the same var declared name node twice. Bail out in this case.
	// TODO(johnlenz): Hash lookups are not free and building scopes are already expensive.
	// Restructure the scope building to avoid this check.
	bloom := scopeBloomMask(name)
	alreadyMaybePresent := (sc.varBloom & bloom) == bloom
	var v *Var
	var h uint32
	var probeSlot uint32
	if sc.varHash != nil {
		h = scopeNameHash(name)
		mask := uint32(len(sc.varHash) - 1)
		probeSlot = h & mask
		if alreadyMaybePresent {
			for {
				idx := sc.varHash[probeSlot]
				if idx == 0 {
					break
				}
				if cand := sc.varList[idx-1]; scopeNameEq(cand.name, name) {
					v = cand
					break
				}
				probeSlot = (probeSlot + 1) & mask
			}
		} else {
			for sc.varHash[probeSlot] != 0 {
				probeSlot = (probeSlot + 1) & mask
			}
		}
		if v == nil && sc.hasImplicitVars {
			if imp := AbstractScope_ImplicitVar_Of(name); imp != nil {
				v = sc.GetOwnImplicitSlot(imp)
			}
		}
	} else {
		if alreadyMaybePresent {
			for i := 0; i < s.smallLen; i++ {
				if scopeNameEq(s.smallVars[i].name, name) {
					v = s.smallVars[i]
					break
				}
			}
		}
		if v == nil && sc.hasImplicitVars {
			if imp := AbstractScope_ImplicitVar_Of(name); imp != nil {
				v = sc.GetOwnImplicitSlot(imp)
			}
		}
	}
	if v != nil {
		if v.GetNode() == n {
			return
		} else if v.IsImplicitGoogNamespace() {
			// this is replacing an implicit provide. treat this as the actual declaration.
			sc.Undeclare(v)
			v = nil
			if sc.varHash != nil {
				mask := uint32(len(sc.varHash) - 1)
				probeSlot = h & mask
				for sc.varHash[probeSlot] != 0 {
					probeSlot = (probeSlot + 1) & mask
				}
			}
		}
	}

	isFnBlock := s.isFnBlock
	var parentVar *Var
	if v == nil && (s.parentBloom&bloom) == bloom {
		parentVar = sc.parent.lookupOwnVarWithBloom(name, bloom)
	}

	input := s.getInput()
	if v != nil ||
		(parentVar != nil && parentVar.IsParam()) ||
		(name == syntacticScopeCreator_ARGUMENTS &&
			(sc.rootNode.IsFunction() || isFnBlock)) {
		if s.redeclarationHandler != SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER && sc.varHash == nil && s.smallLen > 0 {
			s.flushSmallVarsForMore()
		}
		s.redeclarationHandler.OnRedeclaration(sc, name, n, input)
	} else {
		if name == "" {
			panic(jutil.NewIllegalArgumentException())
		}
		if parentVar != nil && !sc.IsBleedingFunctionName(name) {
			panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Illegal shadow: %s", n)))
		}
		sc.varBloom |= bloom
		if sc.varHash == nil {
			if s.smallLen < len(s.smallVars) {
				var var_ *Var
				if sc.arena != nil {
					var_ = sc.arena.allocVar()
					initVar(var_, name, n, sc, s.smallLen, input, nil)
				} else {
					var_ = NewVar(name, n, sc, s.smallLen, input, nil)
				}
				s.smallVars[s.smallLen] = var_
				s.smallLen++
				return
			}
			s.flushSmallVarsForMore()
			h = scopeNameHash(name)
			mask := uint32(len(sc.varHash) - 1)
			probeSlot = h & mask
			for sc.varHash[probeSlot] != 0 {
				probeSlot = (probeSlot + 1) & mask
			}
		}
		var var_ *Var
		idx := int(sc.varCount)
		if sc.arena != nil {
			var_ = sc.arena.allocVar()
			initVar(var_, name, n, sc, idx, input, nil)
		} else {
			var_ = NewVar(name, n, sc, idx, input, nil)
		}
		sc.appendVarToInitialized(var_, probeSlot)
	}
}

// DeclareImplicitGoogNamespaceFromCall declares the implicit namespace 'a'
// from 'goog.provide('a.b.c');
//
// Explicit syntactical definitions like 'var a = {};' supersede this
// definition if present.
//
// exprCall is a goog.module or goog.provide call.
//
// Java: SyntacticScopeCreator.ScopeScanner.declareImplicitGoogNamespaceFromCall(Scope, Node)
func (s *SyntacticScopeCreator_ScopeScanner) DeclareImplicitGoogNamespaceFromCall(sc *Scope, exprCall *rhino.Node) {
	namespaceNode := exprCall.GetFirstChild().GetSecondChild()
	if namespaceNode == nil || !namespaceNode.IsString() {
		// invalid call missing an argument, we warn elsewhere
		return
	}
	namespace := namespaceNode.GetString()
	root := NodeUtil_GetRootOfQualifiedNameString(namespace)

	if root == "" {
		return
	}
	if sc.varHash == nil {
		s.flushSmallVarsForMore()
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
	sc.DeclareImplicitGoogNamespaceIfAbsent(root, namespaceNode)
}

// Function body declarations are not allowed to shadow
// function parameters.
//
// Java: SyntacticScopeCreator.ScopeScanner.isShadowingAllowed(String, Scope)
func syntacticScopeCreator_ScopeScanner_isShadowingAllowed(name string, s *Scope) bool {
	if s.IsFunctionBlockScope() {
		maybeParam := s.GetParent().getOwnSlotV(name)
		return maybeParam == nil || !maybeParam.IsParam()
	}
	return true
}

// Java: SyntacticScopeCreator.ScopeScanner.isLegacyGoogModule(Node)
func syntacticScopeCreator_ScopeScanner_isLegacyGoogModule(n *rhino.Node) bool {
	if !NodeUtil_IsGoogModuleCall(n) {
		return false
	}
	return n.GetNext() != nil && NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n.GetNext())
}

// SyntacticScopeCreator_RedeclarationHandler is the interface for
// injectable duplicate handling.
type SyntacticScopeCreator_RedeclarationHandler interface {
	// Java: SyntacticScopeCreator.RedeclarationHandler.onRedeclaration(Scope, String, Node, CompilerInput)
	OnRedeclaration(s *Scope, name string, n *rhino.Node, input *CompilerInput)
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
	s *Scope, name string, n *rhino.Node, input *CompilerInput) {
}

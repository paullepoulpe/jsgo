// Port of com/google/javascript/jscomp/NodeTraversal.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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

type scopeStackEntry struct {
	node         rhino.Node
	scope        AbstractScope
	untypedScope *Scope
	lookupScope  *Scope
}

var nodeTraversalBlockScopeToken = [256]bool{
	rhino.Token_FOR:          true,
	rhino.Token_FOR_IN:       true,
	rhino.Token_FOR_OF:       true,
	rhino.Token_FOR_AWAIT_OF: true,
	rhino.Token_SWITCH_BODY:  true,
	rhino.Token_CLASS:        true,
}

func (t *NodeTraversal) isBlockScopeBlock(n rhino.Node) bool {
	if t.mayContainSyntheticBlocks && n.IsSyntheticBlock() {
		// Don't create block scope for synthetic blocks.
		return false
	}
	parent := n.GetParent()
	if parent.IsNil() {
		return false
	}
	pt := parent.GetToken()
	// Don't create block scope for switch cases or catch blocks.
	return pt != rhino.Token_CASE && pt != rhino.Token_DEFAULT_CASE && pt != rhino.Token_CATCH
}

// CreatesBlockScope is like NodeUtil.createsBlockScope, but with a fast path
// for isSyntheticBlock.
//
// Java: NodeTraversal.createsBlockScope(Node)
func (t *NodeTraversal) CreatesBlockScope(n rhino.Node) bool {
	tok := n.GetToken()
	if tok == rhino.Token_BLOCK {
		return t.isBlockScopeBlock(n)
	}
	return nodeTraversalBlockScopeToken[uint8(tok)]
}

type cfgStackEntry struct {
	root rhino.Node
	cfg  *ControlFlowGraph[rhino.Node]
}

// cfgScopeStack tracks and lazily computes ControlFlowGraphs for enclosing CFG
// root scopes during a traversal without virtual self-dispatch.
type cfgScopeStack struct {
	cfgs    []cfgStackEntry
	cfgsBuf [8]cfgStackEntry

	// Go-only: performance — the graph of cfgs[i] is computed into
	// cfgStorage[i], which is released when cfgs[i] is popped and reused for
	// the next CFG root at the same depth.
	cfgStorage []*ReusableCfg
}

func (s *cfgScopeStack) resetCfgStack() {
	for i := range s.cfgs {
		if s.cfgs[i].cfg != nil {
			s.cfgStorage[i].release()
		}
	}
	clear(s.cfgs)
	s.cfgs = s.cfgsBuf[:0]
}

func (s *cfgScopeStack) pushCfgRoot(root rhino.Node) {
	if s.cfgs == nil {
		s.cfgs = s.cfgsBuf[:0]
	}
	s.cfgs = append(s.cfgs, cfgStackEntry{root: root})
}

func (s *cfgScopeStack) pushScopeRoot(scopeRoot rhino.Node) {
	if NodeUtil_IsValidCfgRoot(scopeRoot) {
		s.pushCfgRoot(scopeRoot)
	}
}

func (s *cfgScopeStack) popCfgRoot() {
	if len(s.cfgs) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	idx := len(s.cfgs) - 1
	if s.cfgs[idx].cfg != nil {
		s.cfgStorage[idx].release()
	}
	s.cfgs[idx] = cfgStackEntry{}
	s.cfgs = s.cfgs[:idx]
}

func (s *cfgScopeStack) popScopeRoot(scopeRoot rhino.Node) {
	if NodeUtil_IsValidCfgRoot(scopeRoot) {
		s.popCfgRoot()
	}
}

// PeekCurrentScopeCfg returns the already-computed ControlFlowGraph for the
// current scope if one has been materialized, or nil otherwise.
func (s *cfgScopeStack) PeekCurrentScopeCfg() *ControlFlowGraph[rhino.Node] {
	idx := len(s.cfgs) - 1
	if idx < 0 {
		return nil
	}
	return s.cfgs[idx].cfg
}

// GetControlFlowGraph gets the control flow graph for the current JS scope.
func (s *cfgScopeStack) GetControlFlowGraph(compiler AbstractCompiler) *ControlFlowGraph[rhino.Node] {
	idx := len(s.cfgs) - 1
	jutil.CheckState(idx >= 0)
	if cfg := s.cfgs[idx].cfg; cfg != nil {
		return cfg
	}
	for len(s.cfgStorage) <= idx {
		s.cfgStorage = append(s.cfgStorage, &ReusableCfg{})
	}
	result := s.cfgStorage[idx].Compute(
		compiler,
		s.cfgs[idx].root,
		/* shouldTraverseFunctions= */ false,
		/* edgeAnnotations= */ true)
	s.cfgs[idx].cfg = result
	return result
}

// GetEnclosingFunction examines the functions stack for the last instance
// of a function node. When possible, prefer this method over
// NodeUtil.getEnclosingFunction() because this in general looks at less
// nodes.
//
// Java: NodeTraversal.getEnclosingFunction()
func (t *NodeTraversal) GetEnclosingFunction() rhino.Node {
	return t.currentFunction
}

func (t *NodeTraversal) pushScopeQuietly(node rhino.Node) {
	if t.currentNode.IsNil() {
		panic(jutil.NewNullPointerException()) // checkNotNull(currentNode)
	}
	if node.IsNil() {
		panic(jutil.NewNullPointerException()) // checkNotNull(node)
	}
	releaseInline := len(t.scopes) == len(t.scopesBuf) && &t.scopes[0] == &t.scopesBuf[0]
	t.scopes = append(t.scopes, scopeStackEntry{node: node})
	if releaseInline {
		clear(t.scopesBuf[:])
	}
}

// PushScope creates a new scope (e.g. when entering a function).
//
// Java: NodeTraversal.pushScope(Node)
func (t *NodeTraversal) PushScope(node rhino.Node) {
	t.pushScopeQuietly(node)
	if t.scopeCallback != nil {
		t.scopeCallback.EnterScope(t)
	}
}

// PushScopeOfScope creates a new scope (e.g. when entering a function).
//
// Java: NodeTraversal.pushScope(AbstractScope<?, ?>)
func (t *NodeTraversal) PushScopeOfScope(s AbstractScope) {
	t.PushScopeOfScopeWithQuietly(s, false)
}

// PushScopeOfScopeWithQuietly creates a new scope (e.g. when entering a
// function); quietly means don't fire an enterScope callback.
//
// Java: NodeTraversal.pushScope(AbstractScope<?, ?>, boolean)
func (t *NodeTraversal) PushScopeOfScopeWithQuietly(s AbstractScope, quietly bool) {
	if t.currentNode.IsNil() {
		panic(jutil.NewNullPointerException()) // checkNotNull(currentNode)
	}
	var root rhino.Node
	var u *Scope
	if s != nil {
		root = s.GetRootNode()
		u, _ = s.(*Scope)
	}
	releaseInline := len(t.scopes) == len(t.scopesBuf) && &t.scopes[0] == &t.scopesBuf[0]
	t.scopes = append(t.scopes, scopeStackEntry{node: root, scope: s, untypedScope: u, lookupScope: u})
	if releaseInline {
		clear(t.scopesBuf[:])
	}
	if !quietly && t.scopeCallback != nil {
		t.scopeCallback.EnterScope(t)
	}
}

// Java: NodeTraversal.popScope()
func (t *NodeTraversal) PopScope() {
	t.PopScopeWithQuietly(false)
}

// PopScopeWithQuietly pops back to the previous scope (e.g. when leaving a
// function); quietly means don't fire the exitScope callback.
//
// Java: NodeTraversal.popScope(boolean)
func (t *NodeTraversal) PopScopeWithQuietly(quietly bool) {
	if !quietly && t.scopeCallback != nil {
		t.scopeCallback.ExitScope(t)
	}
	t.scopes[len(t.scopes)-1] = scopeStackEntry{}
	t.scopes = t.scopes[:len(t.scopes)-1]
}

// GetScopeRoot returns the current scope's root.
//
// Java: NodeTraversal.getScopeRoot()
func (t *NodeTraversal) GetScopeRoot() rhino.Node {
	roots := len(t.scopes)
	if roots > 0 {
		e := &t.scopes[roots-1]
		if !e.node.IsNil() {
			return e.node
		}
		return e.scope.GetRootNode()
	} else {
		return rhino.NilNode
	}
}

// Java: NodeTraversal.getScopeDepth()
func (t *NodeTraversal) GetScopeDepth() int {
	depth := len(t.scopes)
	jutil.CheckState(depth > 0)
	return depth - 1 // Use 0-based scope depth to be consistent within the compiler
}

// Java: NodeTraversal.getScope()
func (t *NodeTraversal) GetScope() *Scope {
	idx := len(t.scopes) - 1
	if s := t.scopes[idx].untypedScope; s != nil {
		return s
	}
	return t.getScopeSlow(idx)
}

func (t *NodeTraversal) LookupVar(name string) *Var {
	idx := len(t.scopes) - 1
	if s := t.scopes[idx].lookupScope; s != nil {
		return s.GetVar(name)
	}
	return t.resolveLookupScope(idx).GetVar(name)
}

func (t *NodeTraversal) LookupVarByID(id uint32) *Var {
	idx := len(t.scopes) - 1
	if s := t.scopes[idx].lookupScope; s != nil {
		return s.GetVarByID(id)
	}
	return t.resolveLookupScope(idx).GetVarByID(id)
}

func (t *NodeTraversal) LookupVarByNode(n rhino.Node) *Var {
	return t.LookupVarByID(n.GetStringID())
}

func (t *NodeTraversal) resolveLookupScope(rootDepth int) *Scope {
	e := &t.scopes[rootDepth]
	if e.lookupScope != nil {
		return e.lookupScope
	}
	if e.untypedScope != nil {
		if e.untypedScope.scanDone && len(e.untypedScope.varList) == 0 && !e.untypedScope.hasImplicitVars && e.untypedScope.lookupParent != nil {
			e.lookupScope = e.untypedScope.lookupParent
		} else {
			e.lookupScope = e.untypedScope
		}
		return e.lookupScope
	}
	if t.syntacticScopeCreator != nil && e.scope == nil && !e.node.IsNil() && rootDepth > 0 {
		if e.node.IsBlock() {
			if NodeUtil_IsFunctionBlock(e.node) {
				if !NodeUtil_IsBundledGoogModuleScopeRoot(e.node) && syntacticFunctionBlockHasNoDecls(e.node) {
					s := t.resolveLookupScope(rootDepth - 1)
					t.scopes[rootDepth].lookupScope = s
					return s
				}
			} else if !e.node.IsClassMembers() && !NodeUtil_IsClassStaticBlock(e.node) && !syntacticScopeHasBlockDecls(e.node) {
				s := t.resolveLookupScope(rootDepth - 1)
				t.scopes[rootDepth].lookupScope = s
				return s
			}
		} else if NodeUtil_IsAnyFor(e.node) && !syntacticScopeHasBlockDecls(e.node) {
			s := t.resolveLookupScope(rootDepth - 1)
			t.scopes[rootDepth].lookupScope = s
			return s
		}
	}
	s := t.getScopeSlow(rootDepth)
	if s.scanDone && len(s.varList) == 0 && !s.hasImplicitVars && s.lookupParent != nil {
		t.scopes[rootDepth].lookupScope = s.lookupParent
	} else {
		t.scopes[rootDepth].lookupScope = s
	}
	return t.scopes[rootDepth].lookupScope
}

func (t *NodeTraversal) getScopeSlow(rootDepth int) *Scope {
	e := &t.scopes[rootDepth]
	if e.untypedScope != nil {
		return e.untypedScope
	}
	if e.scope != nil {
		s := e.scope.Untyped()
		e.untypedScope = s
		if e.lookupScope == nil {
			e.lookupScope = s
		}
		return s
	}
	if t.syntacticScopeCreator != nil {
		if rootDepth == 0 && t.lazyGlobalScope != nil {
			s := t.lazyGlobalScope.get()
			t.scopes[0].scope = s
			t.scopes[0].untypedScope = s
			t.scopes[0].lookupScope = s
			return s
		}
		var parentScope *Scope
		if rootDepth > 0 {
			parentScope = t.getScopeSlow(rootDepth - 1)
		}
		s := t.syntacticScopeCreator.CreateScopeOfScope(e.node, parentScope)
		t.scopes[rootDepth].scope = s
		t.scopes[rootDepth].untypedScope = s
		t.scopes[rootDepth].lookupScope = s
		return s
	}
	s := t.GetAbstractScopeAt(rootDepth).Untyped()
	t.scopes[rootDepth].untypedScope = s
	if t.scopes[rootDepth].lookupScope == nil {
		t.scopes[rootDepth].lookupScope = s
	}
	return s
}

// Java: NodeTraversal.getTypedScope()
func (t *NodeTraversal) GetTypedScope() *TypedScope {
	return t.GetAbstractScope().Typed()
}

// GetAbstractScope gets the current scope.
//
// Java: NodeTraversal.getAbstractScope()
func (t *NodeTraversal) GetAbstractScope() AbstractScope {
	return t.GetAbstractScopeAt(len(t.scopes) - 1)
}

// GetAbstractScopeAt recusively visits uninstanciated scopes and
// instanciates and stores any necessary scopes to return the requested
// scope.
//
// Java: NodeTraversal.getAbstractScope(int)
func (t *NodeTraversal) GetAbstractScopeAt(rootDepth int) AbstractScope {
	e := &t.scopes[rootDepth]
	if e.scope != nil {
		return e.scope
	}
	if t.syntacticScopeCreator != nil {
		return t.getScopeSlow(rootDepth)
	}
	node := e.node
	// The root scope has a null parent.
	var parentScope AbstractScope
	if rootDepth > 0 {
		parentScope = t.GetAbstractScopeAt(rootDepth - 1)
	}
	scope := t.GetScopeCreator().CreateScope(node, parentScope)
	if u, ok := scope.(*Scope); ok {
		t.scopes[rootDepth].untypedScope = u
	}
	t.scopes[rootDepth].scope = scope
	return scope
}

// Java: NodeTraversal.isHoistScope()
func (t *NodeTraversal) IsHoistScope() bool {
	return nodeTraversal_isHoistScopeRootNode(t.GetScopeRoot())
}

// Java: NodeTraversal.getClosestHoistScopeRoot()
func (t *NodeTraversal) GetClosestHoistScopeRoot() rhino.Node {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		rootNode := t.scopes[i].node
		if nodeTraversal_isHoistScopeRootNode(rootNode) {
			return rootNode
		}
	}

	return rhino.NilNode
}

// Java: NodeTraversal.getClosestContainerScope()
func (t *NodeTraversal) GetClosestContainerScope() AbstractScope {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		rootNode := t.scopes[i].node
		if !NodeUtil_CreatesBlockScope(rootNode) {
			return t.GetAbstractScopeAt(i)
		}
	}
	return nil
}

// Java: NodeTraversal.getClosestHoistScope()
func (t *NodeTraversal) GetClosestHoistScope() AbstractScope {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		rootNode := t.scopes[i].node
		if nodeTraversal_isHoistScopeRootNode(rootNode) {
			return t.GetAbstractScopeAt(i)
		}
	}
	return nil
}

// Java: NodeTraversal.isHoistScopeRootNode(Node)
func nodeTraversal_isHoistScopeRootNode(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_MODULE_BODY, rhino.Token_ROOT, rhino.Token_SCRIPT:
		return true
	default:
		return NodeUtil_IsFunctionBlock(n)
	}
}

func (t *NodeTraversal) closestScopeHoistRoot() (int, rhino.Node) {
	for i := len(t.scopes) - 1; i > 0; i-- {
		r := t.scopes[i].node
		if r.IsNil() {
			r = t.scopes[i].scope.GetRootNode()
		}
		switch r.GetToken() {
		case rhino.Token_FUNCTION, rhino.Token_MODULE_BODY:
			return i, r
		case rhino.Token_BLOCK:
			if NodeUtil_IsFunctionBlock(r) || NodeUtil_IsClassStaticBlock(r) {
				return i, r
			}
		}
	}
	if len(t.scopes) > 0 {
		return 0, t.scopes[0].node
	}
	return 0, rhino.NilNode
}

// GetClosestScopeRootNodeBindingThisOrSuper returns the closest scope
// binding the `this` or `super` keyword
//
// Java: NodeTraversal.getClosestScopeRootNodeBindingThisOrSuper()
func (t *NodeTraversal) GetClosestScopeRootNodeBindingThisOrSuper() rhino.Node {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		rootNode := t.scopes[i].node
		switch rootNode.GetToken() {
		case rhino.Token_FUNCTION:
			if rootNode.IsArrowFunction() {
				continue
			}
			return rootNode
		case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF, rhino.Token_CLASS,
			rhino.Token_MODULE_BODY, rhino.Token_ROOT:
			return rootNode
		case rhino.Token_BLOCK:
			if NodeUtil_IsClassStaticBlock(rootNode) {
				return rootNode
			}
			continue
		default:
			continue
		}
	}
	return rhino.NilNode
}

// Java: NodeTraversal.getScopeCreator()
func (t *NodeTraversal) GetScopeCreator() ScopeCreator {
	if t.scopeCreator == nil {
		if t.syntacticScopeCreator != nil {
			t.scopeCreator = t.syntacticScopeCreator
		} else {
			t.scopeCreator = NewSyntacticScopeCreator(t.compiler)
		}
	}
	return t.scopeCreator
}

// InGlobalScope determines whether the traversal is currently in the global
// scope. Note that this returns false in a global block scope.
//
// Java: NodeTraversal.inGlobalScope()
func (t *NodeTraversal) InGlobalScope() bool {
	return t.GetScopeDepth() == 0
}

// InModuleScope determines whether the traversal is currently in a module
// scope.
//
// Java: NodeTraversal.inModuleScope()
func (t *NodeTraversal) InModuleScope() bool {
	return NodeUtil_IsModuleScopeRoot(t.GetScopeRoot())
}

// Java: NodeTraversal.inGlobalOrModuleScope()
func (t *NodeTraversal) InGlobalOrModuleScope() bool {
	return t.InGlobalScope() || t.InModuleScope()
}

// InFunctionBlockScope determines whether the traversal is currently in the
// scope of the block of a function.
//
// Java: NodeTraversal.inFunctionBlockScope()
func (t *NodeTraversal) InFunctionBlockScope() bool {
	return NodeUtil_IsFunctionBlock(t.GetScopeRoot())
}

// InGlobalHoistScope determines whether the hoist scope of the current
// traversal is global.
//
// Java: NodeTraversal.inGlobalHoistScope()
func (t *NodeTraversal) InGlobalHoistScope() bool {
	return t.currentHoistScopeRoot.IsNil()
}

// InModuleHoistScope determines whether the hoist scope of the current
// traversal is global.
//
// Java: NodeTraversal.inModuleHoistScope()
func (t *NodeTraversal) InModuleHoistScope() bool {
	moduleRoot := t.currentHoistScopeRoot
	if moduleRoot.IsNil() {
		// in global hoist scope
		return false
	}
	if moduleRoot.IsFunction() {
		// For wrapped modules, the function block is the module scope root.
		moduleRoot = moduleRoot.GetLastChild()
	}
	return NodeUtil_IsModuleScopeRoot(moduleRoot)
}

// InitScopeRoots prefills the scopes stack up to a given spot in the AST.
// Allows for starting traversal at any spot while still having correct scope
// state.
//
// Java: NodeTraversal.initScopeRoots(Node)
func (t *NodeTraversal) InitScopeRoots(n rhino.Node) {
	t.initScopeRoots(n, false)
}

func (t *NodeTraversal) initScopeRoots(n rhino.Node, quietly bool) {
	// Go-only: language-adapter — the Deque filled with addFirst: roots are
	// collected innermost first and pushed outermost first.
	var queuedScopeRootsBuf [16]rhino.Node
	queuedScopeRoots := queuedScopeRootsBuf[:0]
	for !n.IsNil() {
		if t.IsScopeRoot(n) {
			queuedScopeRoots = append(queuedScopeRoots, n)
		}
		n = n.GetParent()
	}
	if quietly {
		for i := len(queuedScopeRoots) - 1; i >= 0; i-- {
			t.pushScopeQuietly(queuedScopeRoots[i])
		}
	} else {
		for i := len(queuedScopeRoots) - 1; i >= 0; i-- {
			t.PushScope(queuedScopeRoots[i])
		}
	}
}

// Java: NodeTraversal.isScopeRoot(Node)
func (t *NodeTraversal) IsScopeRoot(n rhino.Node) bool {
	if n.IsRoot() && n.GetParent().IsNil() {
		return true
	} else if n.IsFunction() {
		return true
	} else if NodeUtil_CreatesBlockScope(n) {
		return true
	}
	return false
}

// nodeTraversal_getParent is `s.getParent()` on an AbstractScope<?, ?>.
//
// Go-only: language-adapter — getParent() is not on the wildcard interface
// (see abstract_scope.go); getParentScope() is final and returns getParent().
func nodeTraversal_getParent(s AbstractScope) AbstractScope {
	if p := s.GetParentScope(); p != nil {
		return p.(AbstractScope)
	}
	return nil
}

// TraverseAtScope traverses a parse tree recursively with a scope, starting
// at that scope's root. Omits children of the scope root that are traversed
// in the outer scope (specifically, non-bleeding function and class name
// nodes, class extends clauses, and computed property keys).
//
// Java: NodeTraversal.traverseAtScope(AbstractScope<?, ?>)
func (t *NodeTraversal) TraverseAtScope(s AbstractScope) {
	n := s.GetRootNode()
	t.InitTraversal(n)
	t.currentNode = n
	// Go-only: language-adapter — the Deque used as a stack: parents are
	// collected innermost first and pushed outermost first.
	var parentScopesBuf [16]AbstractScope
	parentScopes := parentScopesBuf[:0]
	temp := nodeTraversal_getParent(s)
	for temp != nil {
		parentScopes = append(parentScopes, temp)
		temp = nodeTraversal_getParent(temp)
	}
	for i := len(parentScopes) - 1; i >= 0; i-- {
		t.PushScopeOfScopeWithQuietly(parentScopes[i], true)
	}

	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		if t.shouldTraverseScriptOrFunction(n, rhino.NilNode) {
			t.PushScopeOfScope(s)

			fnName := n.GetFirstChild()
			args := fnName.GetNext()
			body := args.GetNext()
			if !NodeUtil_IsFunctionDeclaration(n) {
				// Only traverse the function name if it's a bleeding function expression name.
				t.TraverseBranch(fnName, n)
			}
			t.TraverseBranch(args, n)
			t.TraverseBranch(body, n)

			t.PopScope()
			t.visitNode(n, rhino.NilNode)
		}
	case rhino.Token_CLASS:
		if t.shouldTraverse(n, rhino.NilNode) {
			t.PushScopeOfScope(s)

			className := n.GetFirstChild()
			body := n.GetLastChild()

			if NodeUtil_IsClassExpression(n) {
				// Only traverse the class name if it's a bleeding class expression name.
				t.TraverseBranch(className, n)
			}
			// Omit the extends node, which is in the outer scope. Computed property keys are already
			// excluded by handleClassMembers.
			t.TraverseBranch(body, n)

			t.PopScope()
			t.visitNode(n, rhino.NilNode)
		}
	case rhino.Token_BLOCK, rhino.Token_SWITCH_BODY:
		if t.shouldTraverse(n, rhino.NilNode) {
			t.PushScopeOfScope(s)

			// traverseBranch is not called here to avoid re-creating the block scope.
			t.TraverseChildren(n)

			t.PopScope()
			t.visitNode(n, rhino.NilNode)
		}
	case rhino.Token_MEMBER_FIELD_DEF:
		t.PushScopeOfScope(s)
		if t.shouldTraverse(n, rhino.NilNode) {

			// traverseBranch is not called here to avoid re-creating the MemberFieldDef scope.
			t.TraverseChildren(n)

			t.visitNode(n, rhino.NilNode)
		}
		t.PopScope()
	case rhino.Token_COMPUTED_FIELD_DEF:
		t.PushScopeOfScope(s)
		if t.shouldTraverse(n, rhino.NilNode) {

			t.TraverseBranch(n.GetLastChild(), n)

			t.visitNode(n, rhino.NilNode)
		}
		t.PopScope()
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_FOR:
		if t.shouldTraverse(n, rhino.NilNode) {
			t.PushScopeOfScope(s)

			forAssignmentParam := n.GetFirstChild()
			forIterableParam := forAssignmentParam.GetNext()
			forBodyScope := forIterableParam.GetNext()
			t.TraverseBranch(forAssignmentParam, n)
			t.TraverseBranch(forIterableParam, n)
			t.TraverseBranch(forBodyScope, n)

			t.PopScope()
			t.visitNode(n, rhino.NilNode)
		}
	default:
		jutil.CheckStateMsg1(
			s.IsGlobal() || s.IsModuleScope(), "Expected global or module scope. Got: (%s)", s)
		t.TraverseWithScope(n, s)
	}
}

// NodeTraversal_TraverseAtScope traverses the given scope.
//
// Java: NodeTraversal.traverseAtScope(AbstractCompiler, Callback, AbstractScope<?, ?>)
func NodeTraversal_TraverseAtScope(compiler AbstractCompiler, cb NodeTraversal_Callback, scope AbstractScope) {
	b := NodeTraversal_Builder{compiler: compiler, callback: cb}
	b.TraverseAtScope(scope)
}

func traverseScriptUnderRoot(
	t *NodeTraversal,
	scopeRoot rhino.Node,
	scriptParent rhino.Node,
	script rhino.Node,
	lazyGlobal *lazyGlobalScope,
) {
	defer func() {
		if r := recover(); r != nil {
			t.ThrowUnexpectedException(jutil.AsThrowable(r))
		}
	}()
	t.lazyGlobalScope = lazyGlobal
	t.InitTraversal(scopeRoot)
	t.currentNode = scopeRoot
	t.PushScope(scopeRoot)
	if t.shouldTraverse(scopeRoot, rhino.NilNode) {
		if scriptParent != scopeRoot {
			if t.shouldTraverse(scriptParent, scopeRoot) {
				t.TraverseBranch(script, scriptParent)
				t.currentNode = scriptParent
				t.visitNode(scriptParent, scopeRoot)
			}
		} else {
			t.TraverseBranch(script, scriptParent)
		}
		t.currentNode = scopeRoot
		t.visitNode(scopeRoot, rhino.NilNode)
	}
	t.PopScope()
}

// Java: NodeTraversal.traverseScopeRoot(Node)
func (t *NodeTraversal) TraverseScopeRoot(scopeRoot rhino.Node) {
	t.traverseScopeRoot(scopeRoot, false)
}

func (t *NodeTraversal) traverseScopeRoot(scopeRoot rhino.Node, quietly bool) {
	defer func() {
		if r := recover(); r != nil {
			t.ThrowUnexpectedException(jutil.AsThrowable(r))
		}
	}()
	clear(t.scopes)
	t.scopes = t.scopes[:0]
	t.InitTraversal(scopeRoot)
	t.currentNode = scopeRoot
	t.initScopeRoots(scopeRoot.GetParent(), quietly)
	t.TraverseBranch(scopeRoot, scopeRoot.GetParent())
}

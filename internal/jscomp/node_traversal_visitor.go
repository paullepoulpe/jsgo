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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type nodeTraversal_noShouldTraverse interface {
	nodeTraversalNoShouldTraverse()
}

type nodeTraversal_noVisit interface {
	nodeTraversalNoVisit()
}

// nodeTraversal_noVisitLeaves is an internal marker interface for callbacks
// whose Visit method never performs any work on normal leaf nodes
// (n.GetFirstChild() == nil && nodeTraversalTokenKind[uint8(n.GetToken())] == 0).
type nodeTraversal_noVisitLeaves interface {
	noVisitLeaves()
}

// nodeTraversal_onlyNameLeaves is an internal marker interface for callbacks
// whose Visit method on normal leaf nodes (n.GetFirstChild() == nil &&
// nodeTraversalTokenKind[uint8(n.GetToken())] == 0) only acts on Token_NAME.
type nodeTraversal_onlyNameLeaves interface {
	onlyNameLeaves()
}

// nodeTraversal_visitTokens is an internal interface for callbacks whose
// Visit method only acts on nodes with the tokens of the set it returns
// (the typed visitor, node_traversal_visitor.go): the traversal does not
// call Visit for other nodes where it has their token at hand.
type nodeTraversal_visitTokens interface {
	visitTokens() *ast.TokenSet
}

// NodeTraversal_Callback is the callback for tree-based traversals.
type NodeTraversal_Callback interface {
	// ShouldTraverse visits a node in preorder (before its children) and
	// decides whether the node and its children should be traversed.
	//
	// If this method returns true, the node will be visited by Visit in
	// postorder and its children will be visited by both ShouldTraverse in
	// preorder and by Visit in postorder.
	//
	// If this method returns false, the node will not be visited by Visit and
	// its children will neither be visited by ShouldTraverse nor Visit.
	//
	// Siblings are always visited left-to-right.
	//
	// Implementations can have side-effects (e.g. modify the parse tree).
	// Removing the current node is legal, but removing or reordering nodes
	// above the current node may cause nodes to be visited twice or not at
	// all.
	//
	// t is the current traversal, n the current node, parent the parent of
	// the current node. Returns whether the children of this node should be
	// visited.
	//
	// Java: NodeTraversal.Callback.shouldTraverse(NodeTraversal, Node, Node)
	ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool

	// Visit visits a node in postorder (after its children). A node is
	// visited in postorder iff ShouldTraverse returned true for its parent
	// and itself. In particular, the root node is never visited in postorder.
	//
	// Siblings are always visited left-to-right.
	//
	// Implementations can have side-effects (e.g. modify the parse tree).
	// Removing the current node is legal, but removing or reordering nodes
	// above the current node may cause nodes to be visited twice or not at
	// all.
	//
	// t is the current traversal, n the current node, parent the parent of
	// the current node.
	//
	// Java: NodeTraversal.Callback.visit(NodeTraversal, Node, Node)
	Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node)
}

// NodeTraversal_ScopedCallback is a callback that also knows about scope
// changes.
type NodeTraversal_ScopedCallback interface {
	NodeTraversal_Callback

	// EnterScope is called immediately after entering a new scope. The new
	// scope can be accessed through t.GetScope()
	//
	// Java: NodeTraversal.ScopedCallback.enterScope(NodeTraversal)
	EnterScope(t *NodeTraversal)

	// ExitScope is called immediately before exiting a scope. The ending
	// scope can be accessed through t.GetScope()
	//
	// Java: NodeTraversal.ScopedCallback.exitScope(NodeTraversal)
	ExitScope(t *NodeTraversal)
}

// NodeTraversal_AbstractPostOrderCallback is the abstract callback to visit
// all nodes in postorder. Note: Do not create anonymous subclasses of this.
// Instead, write a lambda expression which will be interpreted as an
// AbstractPostOrderCallbackInterface.
type NodeTraversal_AbstractPostOrderCallback interface {
	NodeTraversal_Callback
}

// NodeTraversal_AbstractPostOrderCallbackBase holds the methods of
// AbstractPostOrderCallback; subclasses embed it and define Visit.
type NodeTraversal_AbstractPostOrderCallbackBase struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// ShouldTraverse is final in Java.
//
// Java: NodeTraversal.AbstractPostOrderCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractPostOrderCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

func (NodeTraversal_AbstractPostOrderCallbackBase) nodeTraversalNoShouldTraverse() {}

// NodeTraversal_ExternsSkippingCallback is the abstract callback to visit
// all non-extern nodes in postorder. Note: Even though type-summary nodes are
// included under the externs roots, they are traversed by this callback.
type NodeTraversal_ExternsSkippingCallback interface {
	NodeTraversal_Callback
}

// NodeTraversal_ExternsSkippingCallbackBase holds the methods of
// ExternsSkippingCallback; subclasses embed it and define Visit.
type NodeTraversal_ExternsSkippingCallbackBase struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// ShouldTraverse is final in Java.
//
// Java: NodeTraversal.ExternsSkippingCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_ExternsSkippingCallbackBase) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return !n.IsScript() || !n.IsFromExterns() || NodeUtil_IsFromTypeSummary(n)
}

// NodeTraversal_AbstractPostOrderCallbackInterface is the abstract callback
// to visit all nodes in postorder.
//
// Go-only: language-adapter — a @FunctionalInterface, so a Go func type.
//
// Java: NodeTraversal.AbstractPostOrderCallbackInterface.visit(NodeTraversal, Node, Node)
type NodeTraversal_AbstractPostOrderCallbackInterface func(t *NodeTraversal, n rhino.Node, parent rhino.Node)

// ShouldTraverse implements NodeTraversal_Callback directly on the func type
// so SetCallbackFunc does not need to allocate a wrapper struct on the heap.
func (f NodeTraversal_AbstractPostOrderCallbackInterface) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// Visit implements NodeTraversal_Callback directly on the func type.
func (f NodeTraversal_AbstractPostOrderCallbackInterface) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	f(t, n, parent)
}

// NodeTraversal_AbstractPreOrderCallback is the abstract callback to visit
// all nodes in preorder.
type NodeTraversal_AbstractPreOrderCallback interface {
	NodeTraversal_Callback
}

// NodeTraversal_AbstractPreOrderCallbackBase holds the methods of
// AbstractPreOrderCallback; subclasses embed it and define ShouldTraverse.
type NodeTraversal_AbstractPreOrderCallbackBase struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Visit is final in Java.
//
// Java: NodeTraversal.AbstractPreOrderCallback.visit(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractPreOrderCallbackBase) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
}

func (NodeTraversal_AbstractPreOrderCallbackBase) nodeTraversalNoVisit() {}

// NodeTraversal_AbstractScopedCallback is the abstract scoped callback to
// visit all nodes in postorder.
type NodeTraversal_AbstractScopedCallback interface {
	NodeTraversal_ScopedCallback
}

// NodeTraversal_AbstractScopedCallbackBase holds the methods of
// AbstractScopedCallback; subclasses embed it, define Visit and may
// redefine EnterScope and ExitScope.
type NodeTraversal_AbstractScopedCallbackBase struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// ShouldTraverse is final in Java.
//
// Java: NodeTraversal.AbstractScopedCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractScopedCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

func (NodeTraversal_AbstractScopedCallbackBase) nodeTraversalNoShouldTraverse() {}

// Java: NodeTraversal.AbstractScopedCallback.enterScope(NodeTraversal)
func (b *NodeTraversal_AbstractScopedCallbackBase) EnterScope(t *NodeTraversal) {}

// Java: NodeTraversal.AbstractScopedCallback.exitScope(NodeTraversal)
func (b *NodeTraversal_AbstractScopedCallbackBase) ExitScope(t *NodeTraversal) {}

type cfgStackEntry struct {
	root rhino.Node
	cfg  ControlFlowGraph[rhino.Node]
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
func (s *cfgScopeStack) PeekCurrentScopeCfg() ControlFlowGraph[rhino.Node] {
	idx := len(s.cfgs) - 1
	if idx < 0 {
		return nil
	}
	return s.cfgs[idx].cfg
}

// GetControlFlowGraph gets the control flow graph for the current JS scope.
//
// Java: NodeTraversal.AbstractCfgCallback.getControlFlowGraph(AbstractCompiler)
func (s *cfgScopeStack) GetControlFlowGraph(compiler AbstractCompiler) ControlFlowGraph[rhino.Node] {
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

// NodeTraversal_AbstractCfgCallback is a traversal base class that tracks
// and caches the ControlFlowGraph (CFG) during the traversal. The CFGs are
// constructed lazily.
type NodeTraversal_AbstractCfgCallback interface {
	NodeTraversal_ScopedCallback

	// Java: NodeTraversal.AbstractCfgCallback.getControlFlowGraph(AbstractCompiler)
	GetControlFlowGraph(compiler AbstractCompiler) ControlFlowGraph[rhino.Node]

	// Java: NodeTraversal.AbstractCfgCallback.enterScopeWithCfg(NodeTraversal)
	EnterScopeWithCfg(t *NodeTraversal)

	// Java: NodeTraversal.AbstractCfgCallback.exitScopeWithCfg(NodeTraversal)
	ExitScopeWithCfg(t *NodeTraversal)

	GetNodeTraversal_AbstractCfgCallbackBase() *NodeTraversal_AbstractCfgCallbackBase
}

// NodeTraversal_AbstractCfgCallbackBase holds the fields and methods of
// AbstractCfgCallback for external/test callbacks; production passes embed
// cfgScopeStack directly.
type NodeTraversal_AbstractCfgCallbackBase struct {
	cfgScopeStack
	self NodeTraversal_AbstractCfgCallback
}

// Java: NodeTraversal.AbstractCfgCallback()
func (b *NodeTraversal_AbstractCfgCallbackBase) InitNodeTraversal_AbstractCfgCallbackBase(self NodeTraversal_AbstractCfgCallback) {
	b.self = self
	b.cfgs = b.cfgsBuf[:0]
}

func (b *NodeTraversal_AbstractCfgCallbackBase) GetNodeTraversal_AbstractCfgCallbackBase() *NodeTraversal_AbstractCfgCallbackBase {
	return b
}

// Java: NodeTraversal.AbstractCfgCallback.enterScope(NodeTraversal)
func (b *NodeTraversal_AbstractCfgCallbackBase) EnterScope(t *NodeTraversal) {
	b.pushScopeRoot(t.GetScopeRoot())
	b.self.EnterScopeWithCfg(t)
}

// Java: NodeTraversal.AbstractCfgCallback.exitScope(NodeTraversal)
func (b *NodeTraversal_AbstractCfgCallbackBase) ExitScope(t *NodeTraversal) {
	b.self.ExitScopeWithCfg(t)
	b.popScopeRoot(t.GetScopeRoot())
}

// Java: NodeTraversal.AbstractCfgCallback.enterScopeWithCfg(NodeTraversal)
func (b *NodeTraversal_AbstractCfgCallbackBase) EnterScopeWithCfg(t *NodeTraversal) {}

// Java: NodeTraversal.AbstractCfgCallback.exitScopeWithCfg(NodeTraversal)
func (b *NodeTraversal_AbstractCfgCallbackBase) ExitScopeWithCfg(t *NodeTraversal) {}

// Java: NodeTraversal.AbstractCfgCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractCfgCallbackBase) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// Java: NodeTraversal.AbstractCfgCallback.visit(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractCfgCallbackBase) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
}

// NodeTraversal_AbstractShallowCallback is the abstract callback to visit
// all nodes but not traverse into function bodies.
type NodeTraversal_AbstractShallowCallback interface {
	NodeTraversal_Callback
}

// NodeTraversal_AbstractShallowCallbackBase holds the methods of
// AbstractShallowCallback; subclasses embed it and define Visit.
type NodeTraversal_AbstractShallowCallbackBase struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// ShouldTraverse is final in Java.
//
// Java: NodeTraversal.AbstractShallowCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractShallowCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// We do want to traverse the name of a named function, but we don't
	// want to traverse the arguments or body.
	return parent.IsNil() || !parent.IsFunction() || n == parent.GetFirstChild()
}

// NodeTraversal_AbstractShallowStatementCallback is the abstract callback to
// visit all structure and statement nodes but doesn't traverse into
// functions or expressions.
type NodeTraversal_AbstractShallowStatementCallback interface {
	NodeTraversal_Callback
}

// NodeTraversal_AbstractShallowStatementCallbackBase holds the methods of
// AbstractShallowStatementCallback; subclasses embed it and define Visit.
type NodeTraversal_AbstractShallowStatementCallbackBase struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// ShouldTraverse is final in Java.
//
// Java: NodeTraversal.AbstractShallowStatementCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractShallowStatementCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return NodeUtil_IsShallowStatementTree(parent)
}

// moduleTraversal tracks the current file or bundled goog.module during a
// traversal (Java: NodeTraversal.AbstractModuleCallback) without self-dispatch.
type moduleTraversal struct {
	compiler          AbstractCompiler
	moduleMetadataMap *ModuleMetadataMap

	currentModule *ModuleMetadataMap_ModuleMetadata
	scopeRoot     rhino.Node
	inLoadModule  bool
}

func (m *moduleTraversal) initModuleTraversal(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) {
	m.compiler = compiler
	m.moduleMetadataMap = moduleMetadataMap
}

var nodeTraversal_AbstractModuleCallback_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")

//go:noinline
func failNodeTraversalNull() {
	panic(jutil.NewNullPointerException())
}

// enterNode updates module state on preorder entry and reports whether a new
// module was entered (in which case m.currentModule and m.scopeRoot identify it).
func (m *moduleTraversal) enterNode(t *NodeTraversal, n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		m.currentModule = m.moduleMetadataMap.GetModulesByPath().Get(t.GetInput().GetPath().String())
		if m.currentModule == nil {
			failNodeTraversalNull()
		}
		if n.HasChildren() && n.GetFirstChild().IsModuleBody() {
			m.scopeRoot = n.GetFirstChild()
		} else {
			m.scopeRoot = n
		}
		return true
	case rhino.Token_BLOCK:
		if NodeUtil_IsBundledGoogModuleScopeRoot(n) {
			m.scopeRoot = n
			m.inLoadModule = true
		}
	case rhino.Token_CALL:
		if m.inLoadModule && nodeTraversal_AbstractModuleCallback_GOOG_MODULE.Matches(n.GetFirstChild()) {
			newModule := m.moduleMetadataMap.GetModulesByGoogNamespace().Get(n.GetLastChild().GetString())
			if newModule == nil {
				failNodeTraversalNull()
			}
			// In the event of multiple goog.module statements (an error), don't call enterModule
			// more than once.
			if newModule != m.currentModule {
				m.currentModule = newModule
				return true
			}
		}
	default:
	}
	return false
}

// exitNode updates module state on postorder exit and returns the exited module
// and scope root when leaving a file or bundled goog.module (or nil, NilNode).
func (m *moduleTraversal) exitNode(t *NodeTraversal, n rhino.Node) (exitedModule *ModuleMetadataMap_ModuleMetadata, exitedScopeRoot rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if m.currentModule == nil {
			failNodeTraversalNull()
		}
		exitedModule, exitedScopeRoot = m.currentModule, m.scopeRoot
		m.currentModule = nil
		m.scopeRoot = rhino.NilNode
	case rhino.Token_BLOCK:
		if NodeUtil_IsBundledGoogModuleScopeRoot(n) {
			if m.currentModule == nil {
				failNodeTraversalNull()
			}
			exitedModule, exitedScopeRoot = m.currentModule, m.scopeRoot
			m.scopeRoot = n.GetGrandparent().GetGrandparent()
			m.inLoadModule = false
			m.currentModule = m.moduleMetadataMap.GetModulesByPath().Get(t.GetInput().GetPath().String())
			if m.currentModule == nil {
				failNodeTraversalNull()
			}
		}
	default:
	}
	return exitedModule, exitedScopeRoot
}

// NodeTraversal_AbstractChangedScopeCallback is the callback that fires on
// changed scopes.
type NodeTraversal_AbstractChangedScopeCallback interface {
	NodeTraversal_AbstractPreOrderCallback

	// Java: NodeTraversal.AbstractChangedScopeCallback.enterChangedScopeRoot(AbstractCompiler, Node)
	EnterChangedScopeRoot(compiler AbstractCompiler, root rhino.Node)
}

// NodeTraversal_AbstractChangedScopeCallbackBase holds the methods of
// AbstractChangedScopeCallback (which extends AbstractPreOrderCallback);
// subclasses embed it, call InitNodeTraversal_AbstractChangedScopeCallbackBase
// with themselves and define EnterChangedScopeRoot.
type NodeTraversal_AbstractChangedScopeCallbackBase struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	self NodeTraversal_AbstractChangedScopeCallback
}

// InitNodeTraversal_AbstractChangedScopeCallbackBase is the implicit
// AbstractChangedScopeCallback() constructor; self is the subclass.
//
// Java: NodeTraversal.AbstractChangedScopeCallback()
func (b *NodeTraversal_AbstractChangedScopeCallbackBase) InitNodeTraversal_AbstractChangedScopeCallbackBase(
	self NodeTraversal_AbstractChangedScopeCallback) {
	b.self = self
}

// ShouldTraverse is final in Java.
//
// Java: NodeTraversal.AbstractChangedScopeCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractChangedScopeCallbackBase) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if ChangeTracker_IsChangeScopeRoot(n) && t.GetCompiler().HasScopeChanged(n) {
		b.self.EnterChangedScopeRoot(t.GetCompiler(), n)
	}
	return true
}

// Java: NodeTraversal.builder()
func NodeTraversal_Builder_() *NodeTraversal_Builder {
	return newNodeTraversal_Builder()
}

// NodeTraversal_Builder is the builder.
type NodeTraversal_Builder struct {
	callback                                       NodeTraversal_Callback
	compiler                                       AbstractCompiler
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

// Java: NodeTraversal.Builder()
func newNodeTraversal_Builder() *NodeTraversal_Builder {
	return &NodeTraversal_Builder{obeyDestructuringAndDefaultValueExecutionOrder: false}
}

// Java: NodeTraversal.Builder.setCallback(Callback)
func (b *NodeTraversal_Builder) SetCallback(x NodeTraversal_Callback) *NodeTraversal_Builder {
	b.callback = x
	return b
}

var nodeTraversal_nilPostOrderCallback NodeTraversal_Callback = &nodeTraversal_Builder_postOrderCallback{}

// SetCallbackFunc sets a post-order callback given as a function (a Java
// lambda for AbstractPostOrderCallbackInterface).
//
// Java: NodeTraversal.Builder.setCallback(AbstractPostOrderCallbackInterface)
func (b *NodeTraversal_Builder) SetCallbackFunc(x NodeTraversal_AbstractPostOrderCallbackInterface) *NodeTraversal_Builder {
	if x != nil {
		b.callback = x
	} else {
		b.callback = nodeTraversal_nilPostOrderCallback
	}
	return b
}

// nodeTraversal_Builder_postOrderCallback is the anonymous
// AbstractPostOrderCallback of setCallback(AbstractPostOrderCallbackInterface).
type nodeTraversal_Builder_postOrderCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	x NodeTraversal_AbstractPostOrderCallbackInterface
}

// Java: NodeTraversal.Builder.setCallback(AbstractPostOrderCallbackInterface) (anonymous AbstractPostOrderCallback).visit(NodeTraversal, Node, Node)
func (c *nodeTraversal_Builder_postOrderCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	c.x(t, n, parent)
}

// Java: NodeTraversal.Builder.setCompiler(AbstractCompiler)
func (b *NodeTraversal_Builder) SetCompiler(x AbstractCompiler) *NodeTraversal_Builder {
	b.compiler = x
	return b
}

// Java: NodeTraversal.Builder.setScopeCreator(ScopeCreator)
func (b *NodeTraversal_Builder) SetScopeCreator(x ScopeCreator) *NodeTraversal_Builder {
	b.scopeCreator = x
	return b
}

// Java: NodeTraversal.Builder.setObeyDestructuringAndDefaultValueExecutionOrder(boolean)
func (b *NodeTraversal_Builder) SetObeyDestructuringAndDefaultValueExecutionOrder(x bool) *NodeTraversal_Builder {
	b.obeyDestructuringAndDefaultValueExecutionOrder = x
	return b
}

// Java: NodeTraversal.Builder.build()
func (b *NodeTraversal_Builder) Build() *NodeTraversal {
	return newNodeTraversal(b)
}

// Java: NodeTraversal.Builder.traverse(Node)
func (b *NodeTraversal_Builder) Traverse(root rhino.Node) {
	t := acquireNodeTraversal(b)
	t.Traverse(root)
	releaseNodeTraversal(t)
}

// Java: NodeTraversal.Builder.traverseAtScope(AbstractScope<?, ?>)
func (b *NodeTraversal_Builder) TraverseAtScope(scope AbstractScope) {
	t := acquireNodeTraversal(b)
	t.TraverseAtScope(scope)
	releaseNodeTraversal(t)
}

// Java: NodeTraversal.Builder.traverseRoots(Node, Node)
func (b *NodeTraversal_Builder) TraverseRoots(externs rhino.Node, root rhino.Node) {
	t := acquireNodeTraversal(b)
	t.TraverseRoots(externs, root)
	releaseNodeTraversal(t)
}

// Java: NodeTraversal.Builder.traverseWithScope(Node, AbstractScope<?, ?>)
func (b *NodeTraversal_Builder) TraverseWithScope(root rhino.Node, s AbstractScope) {
	t := acquireNodeTraversal(b)
	t.TraverseWithScope(root, s)
	releaseNodeTraversal(t)
}

func (b *NodeTraversal_Builder) TraverseScopeRoot(scopeNode rhino.Node) {
	t := acquireNodeTraversal(b)
	t.TraverseScopeRoot(scopeNode)
	releaseNodeTraversal(t)
}

func (b *NodeTraversal_Builder) traverseScriptUnderRoot(
	scopeRoot rhino.Node,
	scriptParent rhino.Node,
	script rhino.Node,
	lazyGlobal *lazyGlobalScope,
) {
	t := acquireNodeTraversal(b)
	traverseScriptUnderRoot(t, scopeRoot, scriptParent, script, lazyGlobal)
	releaseNodeTraversal(t)
}

// Go-only: traversal with a typed visitor (internal/ast).

// NodeTraversal_TraverseVisitor traverses root like NodeTraversal_Traverse
// with a post-order callback: after the children of a node, it calls v's
// handler for the node's kind, if v has one. The traversal skips the nodes of
// the kinds v has no handler for, so a visitor costs no more than a callback
// that switches on the token.
func NodeTraversal_TraverseVisitor(compiler AbstractCompiler, root rhino.Node, v *ast.Visitor[*NodeTraversal]) {
	NodeTraversal_Traverse(compiler, root, &nodeTraversal_visitorCallback{v: v, kinds: v.Tokens()})
}

// NodeTraversal_TraverseRootsVisitor traverses externs and root like
// NodeTraversal_TraverseRoots with a post-order ast.Visitor.
func NodeTraversal_TraverseRootsVisitor(compiler AbstractCompiler, externs rhino.Node, root rhino.Node, v *ast.Visitor[*NodeTraversal]) {
	NodeTraversal_TraverseRoots(compiler, &nodeTraversal_visitorCallback{v: v, kinds: v.Tokens()}, externs, root)
}

// nodeTraversal_visitorCallback adapts an ast.Visitor to the callback
// interface.
type nodeTraversal_visitorCallback struct {
	v     *ast.Visitor[*NodeTraversal]
	kinds ast.TokenSet // the tokens v has a handler for
}

func (c *nodeTraversal_visitorCallback) visitTokens() *ast.TokenSet { return &c.kinds }

func (c *nodeTraversal_visitorCallback) ShouldTraverse(*NodeTraversal, rhino.Node, rhino.Node) bool {
	return true
}

func (c *nodeTraversal_visitorCallback) Visit(t *NodeTraversal, n rhino.Node, _ rhino.Node) {
	c.v.Visit(t, n)
}

func (c *nodeTraversal_visitorCallback) nodeTraversalNoShouldTraverse() {}

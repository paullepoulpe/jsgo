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

type cfgStackEntry struct {
	root rhino.Node
	cfg  ControlFlowGraph[rhino.Node]
}

// NodeTraversal_AbstractCfgCallbackBase holds the fields and methods of
// AbstractCfgCallback; subclasses embed it and call
// InitNodeTraversal_AbstractCfgCallbackBase with themselves.
type NodeTraversal_AbstractCfgCallbackBase struct {
	self NodeTraversal_AbstractCfgCallback

	cfgs    []cfgStackEntry
	cfgsBuf [8]cfgStackEntry

	// Go-only: performance — the graph of cfgs[i] is computed into
	// cfgStorage[i], which is released when cfgs[i] is popped and reused for
	// the next CFG root at the same depth.
	cfgStorage []*ReusableCfg
}

// InitNodeTraversal_AbstractCfgCallbackBase is the implicit
// AbstractCfgCallback() constructor; self is the subclass.
//
// Java: NodeTraversal.AbstractCfgCallback()
func (b *NodeTraversal_AbstractCfgCallbackBase) InitNodeTraversal_AbstractCfgCallbackBase(self NodeTraversal_AbstractCfgCallback) {
	b.self = self
	b.cfgs = b.cfgsBuf[:0]
}

// GetNodeTraversal_AbstractCfgCallbackBase gives access to the base through
// the interface.
func (b *NodeTraversal_AbstractCfgCallbackBase) GetNodeTraversal_AbstractCfgCallbackBase() *NodeTraversal_AbstractCfgCallbackBase {
	return b
}

// PeekCurrentScopeCfg returns the already-computed ControlFlowGraph for the
// current scope if one has been materialized, or nil otherwise.
func (b *NodeTraversal_AbstractCfgCallbackBase) PeekCurrentScopeCfg() ControlFlowGraph[rhino.Node] {
	idx := len(b.cfgs) - 1
	if idx < 0 {
		return nil
	}
	return b.cfgs[idx].cfg
}

// GetControlFlowGraph gets the control flow graph for the current JS scope.
//
// Java: NodeTraversal.AbstractCfgCallback.getControlFlowGraph(AbstractCompiler)
func (b *NodeTraversal_AbstractCfgCallbackBase) GetControlFlowGraph(compiler AbstractCompiler) ControlFlowGraph[rhino.Node] {
	idx := len(b.cfgs) - 1
	jutil.CheckState(idx >= 0)
	if cfg := b.cfgs[idx].cfg; cfg != nil {
		return cfg
	}
	for len(b.cfgStorage) <= idx {
		b.cfgStorage = append(b.cfgStorage, &ReusableCfg{})
	}
	result := b.cfgStorage[idx].Compute(
		compiler,
		b.cfgs[idx].root,
		/* shouldTraverseFunctions= */ false,
		/* edgeAnnotations= */ true)
	b.cfgs[idx].cfg = result
	return result
}

// EnterScope is final in Java.
//
// Java: NodeTraversal.AbstractCfgCallback.enterScope(NodeTraversal)
func (b *NodeTraversal_AbstractCfgCallbackBase) EnterScope(t *NodeTraversal) {
	currentScopeRoot := t.GetScopeRoot()
	if NodeUtil_IsValidCfgRoot(currentScopeRoot) {
		b.cfgs = append(b.cfgs, cfgStackEntry{root: currentScopeRoot})
	}
	b.self.EnterScopeWithCfg(t)
}

// ExitScope is final in Java.
//
// Java: NodeTraversal.AbstractCfgCallback.exitScope(NodeTraversal)
func (b *NodeTraversal_AbstractCfgCallbackBase) ExitScope(t *NodeTraversal) {
	b.self.ExitScopeWithCfg(t)
	currentScopeRoot := t.GetScopeRoot()
	if NodeUtil_IsValidCfgRoot(currentScopeRoot) {
		if len(b.cfgs) == 0 {
			panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
		}
		idx := len(b.cfgs) - 1
		if b.cfgs[idx].cfg != nil {
			b.cfgStorage[idx].release()
		}
		b.cfgs[idx] = cfgStackEntry{}
		b.cfgs = b.cfgs[:idx]
	}
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

// NodeTraversal_AbstractModuleCallback is the abstract callback that knows
// when a global script, goog.provide file, goog.module, goog.loadModule, ES
// module or CommonJS module is entered or exited. This includes both whole
// file modules and bundled modules, as well as files in the global scope.
type NodeTraversal_AbstractModuleCallback interface {
	NodeTraversal_Callback

	// EnterModule is called when the traversal enters a global file or
	// module. currentModule is the entered global file or module;
	// moduleScopeRoot the root scope for the entered module or SCRIPT for
	// global files.
	//
	// Java: NodeTraversal.AbstractModuleCallback.enterModule(ModuleMetadata, Node)
	EnterModule(currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot rhino.Node)

	// ExitModule is called when the traversal exits a global file or module.
	// oldModule is the exited global file or module; moduleScopeRoot the root
	// scope for the exited module or SCRIPT for global files.
	//
	// Java: NodeTraversal.AbstractModuleCallback.exitModule(ModuleMetadata, Node)
	ExitModule(oldModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot rhino.Node)

	// Java: NodeTraversal.AbstractModuleCallback.shouldTraverse(NodeTraversal, Node, ModuleMetadata, Node)
	ShouldTraverseWithModule(
		t *NodeTraversal, n rhino.Node, currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot rhino.Node) bool

	// Java: NodeTraversal.AbstractModuleCallback.visit(NodeTraversal, Node, ModuleMetadata, Node)
	VisitWithModule(
		t *NodeTraversal, n rhino.Node, currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot rhino.Node)

	GetNodeTraversal_AbstractModuleCallbackBase() *NodeTraversal_AbstractModuleCallbackBase
}

// NodeTraversal_AbstractModuleCallbackBase holds the fields and methods of
// AbstractModuleCallback; subclasses embed it and call
// InitNodeTraversal_AbstractModuleCallbackBase with themselves.
type NodeTraversal_AbstractModuleCallbackBase struct {
	self NodeTraversal_AbstractModuleCallback

	compiler          AbstractCompiler
	moduleMetadataMap *ModuleMetadataMap

	currentModule *ModuleMetadataMap_ModuleMetadata
	scopeRoot     rhino.Node
	inLoadModule  bool
}

// InitNodeTraversal_AbstractModuleCallbackBase is the
// AbstractModuleCallback constructor; self is the subclass.
//
// Java: NodeTraversal.AbstractModuleCallback(AbstractCompiler, ModuleMetadataMap)
func (b *NodeTraversal_AbstractModuleCallbackBase) InitNodeTraversal_AbstractModuleCallbackBase(
	self NodeTraversal_AbstractModuleCallback, compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) {
	b.self = self
	b.compiler = compiler
	b.moduleMetadataMap = moduleMetadataMap
}

// GetNodeTraversal_AbstractModuleCallbackBase gives access to the base
// through the interface.
func (b *NodeTraversal_AbstractModuleCallbackBase) GetNodeTraversal_AbstractModuleCallbackBase() *NodeTraversal_AbstractModuleCallbackBase {
	return b
}

// Java: NodeTraversal.AbstractModuleCallback.enterModule(ModuleMetadata, Node)
func (b *NodeTraversal_AbstractModuleCallbackBase) EnterModule(
	currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot rhino.Node) {
}

// Java: NodeTraversal.AbstractModuleCallback.exitModule(ModuleMetadata, Node)
func (b *NodeTraversal_AbstractModuleCallbackBase) ExitModule(
	oldModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot rhino.Node) {
}

var nodeTraversal_AbstractModuleCallback_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")

//go:noinline
func failNodeTraversalNull() {
	panic(jutil.NewNullPointerException())
}

// ShouldTraverse is final in Java.
//
// Java: NodeTraversal.AbstractModuleCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractModuleCallbackBase) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		b.currentModule = b.moduleMetadataMap.GetModulesByPath().Get(t.GetInput().GetPath().String())
		if b.currentModule == nil {
			failNodeTraversalNull()
		}
		if n.HasChildren() && n.GetFirstChild().IsModuleBody() {
			b.scopeRoot = n.GetFirstChild()
		} else {
			b.scopeRoot = n
		}
		b.self.EnterModule(b.currentModule, b.scopeRoot)
	case rhino.Token_BLOCK:
		if NodeUtil_IsBundledGoogModuleScopeRoot(n) {
			b.scopeRoot = n
			b.inLoadModule = true
		}
	case rhino.Token_CALL:
		if b.inLoadModule && nodeTraversal_AbstractModuleCallback_GOOG_MODULE.Matches(n.GetFirstChild()) {
			newModule := b.moduleMetadataMap.GetModulesByGoogNamespace().Get(n.GetLastChild().GetString())
			if newModule == nil {
				failNodeTraversalNull()
			}
			// In the event of multiple goog.module statements (an error), don't call enterModule
			// more than once.
			if newModule != b.currentModule {
				b.currentModule = newModule
				b.self.EnterModule(b.currentModule, b.scopeRoot)
			}
		}
	default:
	}
	return b.self.ShouldTraverseWithModule(t, n, b.currentModule, b.scopeRoot)
}

// ShouldTraverseWithModule: see Callback.ShouldTraverse.
//
// t is the current traversal, n the current node, currentModule the current
// module, or null if not inside a module (e.g. AST root), moduleScopeRoot the
// root scope for the current module, or null if not inside a module (e.g.
// AST root). Returns whether the children of this node should be visited.
//
// Java: NodeTraversal.AbstractModuleCallback.shouldTraverse(NodeTraversal, Node, ModuleMetadata, Node)
func (b *NodeTraversal_AbstractModuleCallbackBase) ShouldTraverseWithModule(
	t *NodeTraversal,
	n rhino.Node,
	currentModule *ModuleMetadataMap_ModuleMetadata,
	moduleScopeRoot rhino.Node) bool {
	return true
}

// Visit is final in Java.
//
// Java: NodeTraversal.AbstractModuleCallback.visit(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractModuleCallbackBase) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if b.currentModule == nil {
			failNodeTraversalNull()
		}
		b.self.ExitModule(b.currentModule, b.scopeRoot)
		b.currentModule = nil
		b.scopeRoot = rhino.NilNode
	case rhino.Token_BLOCK:
		if NodeUtil_IsBundledGoogModuleScopeRoot(n) {
			if b.currentModule == nil {
				failNodeTraversalNull()
			}
			b.self.ExitModule(b.currentModule, b.scopeRoot)
			b.scopeRoot = n.GetGrandparent().GetGrandparent()
			b.inLoadModule = false
			b.currentModule = b.moduleMetadataMap.GetModulesByPath().Get(t.GetInput().GetPath().String())
			if b.currentModule == nil {
				failNodeTraversalNull()
			}
		}
	default:
	}

	b.self.VisitWithModule(t, n, b.currentModule, b.scopeRoot)
}

// VisitWithModule: see Callback.Visit.
//
// t is the current traversal, n the current node, currentModule the current
// module, or null if not inside a module (e.g. AST root), moduleScopeRoot the
// root scope for the current module, or null if not inside a module (e.g.
// AST root).
//
// Java: NodeTraversal.AbstractModuleCallback.visit(NodeTraversal, Node, ModuleMetadata, Node)
func (b *NodeTraversal_AbstractModuleCallbackBase) VisitWithModule(
	t *NodeTraversal,
	n rhino.Node,
	currentModule *ModuleMetadataMap_ModuleMetadata,
	moduleScopeRoot rhino.Node) {
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

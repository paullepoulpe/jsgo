// Port of com/google/javascript/jscomp/NodeTraversal.java @ closure-compiler v20260415

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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Porting compiler passes against NodeTraversal (read this first).
//
// Callbacks. NodeTraversal.Callback is the Go interface
// NodeTraversal_Callback (ShouldTraverse + Visit); ScopedCallback is
// NodeTraversal_ScopedCallback (adds EnterScope + ExitScope). As in Java, a
// traversal fires enterScope/exitScope only if its callback is a
// NodeTraversal_ScopedCallback, which in Go means "has EnterScope and
// ExitScope methods".
//
// Abstract callback classes. Each abstract class Foo below is an interface
// NodeTraversal_Foo plus a struct NodeTraversal_FooBase holding Foo's fields
// and method bodies (PORTING.md §4). A pass that extends Foo embeds the base
// and defines the abstract methods:
//
//	// Java: final class MyPass extends NodeTraversal.AbstractPostOrderCallback
//	type MyPass struct {
//		NodeTraversal_AbstractPostOrderCallbackBase
//		compiler AbstractCompiler
//	}
//	func (p *MyPass) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) { ... }
//
//	NodeTraversal_Traverse(compiler, root, p) // p is a NodeTraversal_Callback
//
// The bases are:
//
//	AbstractPostOrderCallback        ShouldTraverse returns true; define Visit.
//	ExternsSkippingCallback          ShouldTraverse skips extern SCRIPTs; define Visit.
//	AbstractPreOrderCallback         Visit is empty; define ShouldTraverse.
//	AbstractScopedCallback           ShouldTraverse true, empty EnterScope/ExitScope
//	                                 (overridable); define Visit.
//	AbstractShallowCallback          ShouldTraverse stops at function params/bodies; define Visit.
//	AbstractShallowStatementCallback ShouldTraverse stops below statements; define Visit.
//	AbstractCfgCallback              *self*; EnterScope/ExitScope maintain the CFG stack and
//	                                 call the overridable EnterScopeWithCfg/ExitScopeWithCfg;
//	                                 ShouldTraverse/Visit are overridable defaults.
//	AbstractModuleCallback           *self*; ShouldTraverse/Visit track the module and call
//	                                 the overridable EnterModule/ExitModule and
//	                                 ShouldTraverseWithModule/VisitWithModule (the Java
//	                                 4-argument shouldTraverse/visit overloads).
//	AbstractChangedScopeCallback     *self*; ShouldTraverse calls the abstract
//	                                 EnterChangedScopeRoot; Visit is empty (pre-order).
//
// The bases marked *self* make virtual calls into the subclass, so the
// subclass constructor must hand itself to the base (the Go form of
// `super(...)`), before anything can call it:
//
//	// Java: class CoalesceVariableNames extends NodeTraversal.AbstractCfgCallback
//	type CoalesceVariableNames struct {
//		NodeTraversal_AbstractCfgCallbackBase
//		...
//	}
//	func NewCoalesceVariableNames(...) *CoalesceVariableNames {
//		c := &CoalesceVariableNames{}
//		c.InitNodeTraversal_AbstractCfgCallbackBase(c)
//		...
//		return c
//	}
//	func (c *CoalesceVariableNames) EnterScopeWithCfg(t *NodeTraversal) { ... }
//
// AbstractModuleCallback's constructor arguments go to
// InitNodeTraversal_AbstractModuleCallbackBase(self, compiler,
// moduleMetadataMap); subclasses read the protected `compiler` field
// directly (it is promoted). Java's final methods (e.g. ShouldTraverse of
// AbstractPostOrderCallback) must not be redefined by the subclass: Go cannot
// forbid it, and a redefinition would silently change the traversal.
//
// A lambda for NodeTraversal.AbstractPostOrderCallbackInterface is a Go func
// of type NodeTraversal_AbstractPostOrderCallbackInterface, passed to
// Builder.SetCallbackFunc (Java: setCallback(AbstractPostOrderCallbackInterface)).
//
// Names. NodeTraversal.builder() is NodeTraversal_Builder_() (the type is
// NodeTraversal_Builder); the static traverse(compiler, root, cb),
// traverseRoots(compiler, cb, externs, root) and traverseScopeRoots(...) are
// NodeTraversal_Traverse, NodeTraversal_TraverseRoots and
// NodeTraversal_TraverseScopeRoots. report(Node, Node, DiagnosticType,
// String...) is ReportWithNodeRange (cf. JSError_MakeWithNodeRange);
// reportCodeChange(Node) is ReportCodeChangeWithNode. getSourceName()
// returns a *string (a Java String that may be null). getScope() is a
// *Scope, getTypedScope() a *TypedScope, getAbstractScope() the wildcard
// AbstractScope.

// NodeTraversal allows an iteration through the nodes in the parse tree, and
// facilitates the optimizations on the parse tree.
type NodeTraversal struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCallback                                  NodeTraversal_ScopedCallback
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool

	// Contains the current node
	currentNode *rhino.Node

	currentHoistScopeRoot *rhino.Node

	// Contains the current FUNCTION node if there is one, otherwise null.
	currentFunction *rhino.Node

	// Contains the enclosing SCRIPT node if there is one, otherwise null.
	currentScript *rhino.Node

	// The change scope for the current node being visiteds
	currentChangeScope *rhino.Node

	// The chain scope for the currentNode being visited. Scopes are relatively
	// expensive to build so they are built lazily. The list contains
	// instanciated `AbstractScope` or the `Node` representing the root of the
	// scope.
	//
	// If `AbstractScope` is seen all previous entries will also have been
	// instanciated due to the heirarchical nature of scopes.
	//
	// Go-only: language-adapter — each element is a *rhino.Node or an
	// AbstractScope.
	scopes []any

	// The current source file name
	sourceName *string

	// The current input
	inputId *rhino.InputId

	compilerInput *CompilerInput
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
	ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool

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
	Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node)
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
func (b *NodeTraversal_AbstractPostOrderCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return true
}

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
func (b *NodeTraversal_ExternsSkippingCallbackBase) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return !n.IsScript() || !n.IsFromExterns() || NodeUtil_IsFromTypeSummary(n)
}

// NodeTraversal_AbstractPostOrderCallbackInterface is the abstract callback
// to visit all nodes in postorder.
//
// Go-only: language-adapter — a @FunctionalInterface, so a Go func type.
//
// Java: NodeTraversal.AbstractPostOrderCallbackInterface.visit(NodeTraversal, Node, Node)
type NodeTraversal_AbstractPostOrderCallbackInterface func(t *NodeTraversal, n *rhino.Node, parent *rhino.Node)

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
func (b *NodeTraversal_AbstractPreOrderCallbackBase) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
}

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
func (b *NodeTraversal_AbstractScopedCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return true
}

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
	GetControlFlowGraph(compiler AbstractCompiler) ControlFlowGraph[*rhino.Node]

	// Java: NodeTraversal.AbstractCfgCallback.enterScopeWithCfg(NodeTraversal)
	EnterScopeWithCfg(t *NodeTraversal)

	// Java: NodeTraversal.AbstractCfgCallback.exitScopeWithCfg(NodeTraversal)
	ExitScopeWithCfg(t *NodeTraversal)

	GetNodeTraversal_AbstractCfgCallbackBase() *NodeTraversal_AbstractCfgCallbackBase
}

// NodeTraversal_AbstractCfgCallbackBase holds the fields and methods of
// AbstractCfgCallback; subclasses embed it and call
// InitNodeTraversal_AbstractCfgCallbackBase with themselves.
type NodeTraversal_AbstractCfgCallbackBase struct {
	self NodeTraversal_AbstractCfgCallback

	// Go-only: language-adapter — the ArrayDeque<Object> used as a stack
	// (push/peek/pop at the end); each element is a *rhino.Node (a CFG root
	// whose graph is not built yet) or a ControlFlowGraph[*rhino.Node].
	cfgs []any
}

// InitNodeTraversal_AbstractCfgCallbackBase is the implicit
// AbstractCfgCallback() constructor; self is the subclass.
//
// Java: NodeTraversal.AbstractCfgCallback()
func (b *NodeTraversal_AbstractCfgCallbackBase) InitNodeTraversal_AbstractCfgCallbackBase(self NodeTraversal_AbstractCfgCallback) {
	b.self = self
	b.cfgs = nil
}

// GetNodeTraversal_AbstractCfgCallbackBase gives access to the base through
// the interface.
func (b *NodeTraversal_AbstractCfgCallbackBase) GetNodeTraversal_AbstractCfgCallbackBase() *NodeTraversal_AbstractCfgCallbackBase {
	return b
}

// GetControlFlowGraph gets the control flow graph for the current JS scope.
//
// Java: NodeTraversal.AbstractCfgCallback.getControlFlowGraph(AbstractCompiler)
func (b *NodeTraversal_AbstractCfgCallbackBase) GetControlFlowGraph(compiler AbstractCompiler) ControlFlowGraph[*rhino.Node] {
	var result ControlFlowGraph[*rhino.Node]
	var o any
	if len(b.cfgs) > 0 {
		o = b.cfgs[len(b.cfgs)-1]
	}
	jutil.CheckState(o != nil)
	if cfgRoot, ok := o.(*rhino.Node); ok {
		result = ControlFlowAnalysis_Builder_().
			SetCompiler(compiler).
			SetCfgRoot(cfgRoot).
			SetIncludeEdgeAnnotations(true).
			ComputeCfg()
		b.cfgs = b.cfgs[:len(b.cfgs)-1]
		b.cfgs = append(b.cfgs, result)
	} else {
		result = o.(ControlFlowGraph[*rhino.Node])
	}
	return result
}

// EnterScope is final in Java.
//
// Java: NodeTraversal.AbstractCfgCallback.enterScope(NodeTraversal)
func (b *NodeTraversal_AbstractCfgCallbackBase) EnterScope(t *NodeTraversal) {
	currentScopeRoot := t.GetScopeRoot()
	if NodeUtil_IsValidCfgRoot(currentScopeRoot) {
		b.cfgs = append(b.cfgs, currentScopeRoot)
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
		b.cfgs[len(b.cfgs)-1] = nil
		b.cfgs = b.cfgs[:len(b.cfgs)-1]
	}
}

// Java: NodeTraversal.AbstractCfgCallback.enterScopeWithCfg(NodeTraversal)
func (b *NodeTraversal_AbstractCfgCallbackBase) EnterScopeWithCfg(t *NodeTraversal) {}

// Java: NodeTraversal.AbstractCfgCallback.exitScopeWithCfg(NodeTraversal)
func (b *NodeTraversal_AbstractCfgCallbackBase) ExitScopeWithCfg(t *NodeTraversal) {}

// Java: NodeTraversal.AbstractCfgCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractCfgCallbackBase) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return true
}

// Java: NodeTraversal.AbstractCfgCallback.visit(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractCfgCallbackBase) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
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
func (b *NodeTraversal_AbstractShallowCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	// We do want to traverse the name of a named function, but we don't
	// want to traverse the arguments or body.
	return parent == nil || !parent.IsFunction() || n == parent.GetFirstChild()
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
func (b *NodeTraversal_AbstractShallowStatementCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
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
	EnterModule(currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot *rhino.Node)

	// ExitModule is called when the traversal exits a global file or module.
	// oldModule is the exited global file or module; moduleScopeRoot the root
	// scope for the exited module or SCRIPT for global files.
	//
	// Java: NodeTraversal.AbstractModuleCallback.exitModule(ModuleMetadata, Node)
	ExitModule(oldModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot *rhino.Node)

	// Java: NodeTraversal.AbstractModuleCallback.shouldTraverse(NodeTraversal, Node, ModuleMetadata, Node)
	ShouldTraverseWithModule(
		t *NodeTraversal, n *rhino.Node, currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot *rhino.Node) bool

	// Java: NodeTraversal.AbstractModuleCallback.visit(NodeTraversal, Node, ModuleMetadata, Node)
	VisitWithModule(
		t *NodeTraversal, n *rhino.Node, currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot *rhino.Node)

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
	scopeRoot     *rhino.Node
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
	currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot *rhino.Node) {
}

// Java: NodeTraversal.AbstractModuleCallback.exitModule(ModuleMetadata, Node)
func (b *NodeTraversal_AbstractModuleCallbackBase) ExitModule(
	oldModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot *rhino.Node) {
}

var nodeTraversal_AbstractModuleCallback_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")

// ShouldTraverse is final in Java.
//
// Java: NodeTraversal.AbstractModuleCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractModuleCallbackBase) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		b.currentModule = b.moduleMetadataMap.GetModulesByPath().Get(t.GetInput().GetPath().String())
		jutil.CheckNotNull(b.currentModule)
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
			jutil.CheckNotNull(newModule)
			// In the event of multiple goog.module statements (an error), don't call enterModule
			// more than once.
			if !base.JSCompObjects_Identical(newModule, b.currentModule) {
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
	n *rhino.Node,
	currentModule *ModuleMetadataMap_ModuleMetadata,
	moduleScopeRoot *rhino.Node) bool {
	return true
}

// Visit is final in Java.
//
// Java: NodeTraversal.AbstractModuleCallback.visit(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractModuleCallbackBase) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		jutil.CheckNotNull(b.currentModule)
		b.self.ExitModule(b.currentModule, b.scopeRoot)
		b.currentModule = nil
		b.scopeRoot = nil
	case rhino.Token_BLOCK:
		if NodeUtil_IsBundledGoogModuleScopeRoot(n) {
			jutil.CheckNotNull(b.currentModule)
			b.self.ExitModule(b.currentModule, b.scopeRoot)
			b.scopeRoot = n.GetGrandparent().GetGrandparent()
			b.inLoadModule = false
			b.currentModule = b.moduleMetadataMap.GetModulesByPath().Get(t.GetInput().GetPath().String())
			jutil.CheckNotNull(b.currentModule)
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
	n *rhino.Node,
	currentModule *ModuleMetadataMap_ModuleMetadata,
	moduleScopeRoot *rhino.Node) {
}

// NodeTraversal_AbstractChangedScopeCallback is the callback that fires on
// changed scopes.
type NodeTraversal_AbstractChangedScopeCallback interface {
	NodeTraversal_AbstractPreOrderCallback

	// Java: NodeTraversal.AbstractChangedScopeCallback.enterChangedScopeRoot(AbstractCompiler, Node)
	EnterChangedScopeRoot(compiler AbstractCompiler, root *rhino.Node)
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
func (b *NodeTraversal_AbstractChangedScopeCallbackBase) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
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

// SetCallbackFunc sets a post-order callback given as a function (a Java
// lambda for AbstractPostOrderCallbackInterface).
//
// Java: NodeTraversal.Builder.setCallback(AbstractPostOrderCallbackInterface)
func (b *NodeTraversal_Builder) SetCallbackFunc(x NodeTraversal_AbstractPostOrderCallbackInterface) *NodeTraversal_Builder {
	b.callback = &nodeTraversal_Builder_postOrderCallback{x: x}
	return b
}

// nodeTraversal_Builder_postOrderCallback is the anonymous
// AbstractPostOrderCallback of setCallback(AbstractPostOrderCallbackInterface).
type nodeTraversal_Builder_postOrderCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	x NodeTraversal_AbstractPostOrderCallbackInterface
}

// Java: NodeTraversal.Builder.setCallback(AbstractPostOrderCallbackInterface) (anonymous AbstractPostOrderCallback).visit(NodeTraversal, Node, Node)
func (c *nodeTraversal_Builder_postOrderCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
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
func (b *NodeTraversal_Builder) Traverse(root *rhino.Node) {
	b.Build().Traverse(root)
}

// Java: NodeTraversal.Builder.traverseAtScope(AbstractScope<?, ?>)
func (b *NodeTraversal_Builder) TraverseAtScope(scope AbstractScope) {
	b.Build().TraverseAtScope(scope)
}

// Java: NodeTraversal.Builder.traverseRoots(Node, Node)
func (b *NodeTraversal_Builder) TraverseRoots(externs *rhino.Node, root *rhino.Node) {
	b.Build().TraverseRoots(externs, root)
}

// Java: NodeTraversal.Builder.traverseWithScope(Node, AbstractScope<?, ?>)
func (b *NodeTraversal_Builder) TraverseWithScope(root *rhino.Node, s AbstractScope) {
	b.Build().TraverseWithScope(root, s)
}

// Java: NodeTraversal(Builder)
func newNodeTraversal(builder *NodeTraversal_Builder) *NodeTraversal {
	t := &NodeTraversal{}
	t.compiler = jutil.CheckNotNull(builder.compiler)
	t.callback = jutil.CheckNotNull(builder.callback)
	if scopedCallback, ok := t.callback.(NodeTraversal_ScopedCallback); ok {
		t.scopeCallback = scopedCallback
	} else {
		t.scopeCallback = nil
	}
	if builder.scopeCreator == nil {
		t.scopeCreator = NewSyntacticScopeCreator(t.compiler)
	} else {
		t.scopeCreator = builder.scopeCreator
	}
	t.obeyDestructuringAndDefaultValueExecutionOrder = builder.obeyDestructuringAndDefaultValueExecutionOrder
	return t
}

// Java: NodeTraversal.throwUnexpectedException(Throwable)
func (t *NodeTraversal) ThrowUnexpectedException(unexpectedException jutil.Throwable) {
	// If there's an unexpected exception, try to get the
	// line number of the code that caused it.
	message := jutil.GetMessageOrNull(unexpectedException)

	// TODO(user): It is possible to get more information if currentNode or
	// its parent is missing. We still have the scope stack in which it is still
	// very useful to find out at least which function caused the exception.
	if t.currentScript != nil {
		parentContext := ""
		if t.currentNode != nil {
			parentContext = t.FormatNodeContext("Parent", t.currentNode.GetParent())
		}
		m := jutil.ToJavaString(jutil.GetMessageOrNull(unexpectedException)) +
			"\n" +
			t.FormatNodeContext("Node", t.currentNode) +
			parentContext
		message = &m
	}
	// Go-only: language-adapter — AbstractCompiler.throwInternalError takes a
	// Go string, in which "" stands for a null message.
	msg := ""
	if message != nil {
		msg = *message
	}
	t.compiler.ThrowInternalError(msg, unexpectedException)
}

// Java: NodeTraversal.formatNodeContext(String, Node)
func (t *NodeTraversal) FormatNodeContext(label string, n *rhino.Node) string {
	if n == nil {
		return "  " + label + ": NULL"
	}
	return "  " + label + "(" + n.ToStringWithFlags(false, false, false) + "): " + t.FormatNodePosition(n)
}

// Traverse traverses a parse tree recursively.
//
// Java: NodeTraversal.traverse(Node)
func (t *NodeTraversal) Traverse(root *rhino.Node) {
	jutil.Try(func() {
		t.InitTraversal(root)
		t.currentNode = root
		t.PushScope(root)
		// null parent ensures that the shallow callbacks will traverse root
		t.TraverseBranch(root, nil)
		t.PopScope()
	}, func(unexpectedException jutil.Throwable) {
		t.ThrowUnexpectedException(unexpectedException)
	})
}

// NodeTraversal_Traverse traverses using the SyntacticScopeCreator
//
// Java: NodeTraversal.traverse(AbstractCompiler, Node, Callback)
func NodeTraversal_Traverse(compiler AbstractCompiler, root *rhino.Node, cb NodeTraversal_Callback) {
	NodeTraversal_Builder_().SetCompiler(compiler).SetCallback(cb).Traverse(root)
}

// Java: NodeTraversal.traverseRoots(Node, Node)
func (t *NodeTraversal) TraverseRoots(externs *rhino.Node, root *rhino.Node) {
	jutil.Try(func() {
		scopeRoot := externs.GetParent()
		jutil.CheckNotNull(scopeRoot)

		t.InitTraversal(scopeRoot)
		t.currentNode = scopeRoot
		t.PushScope(scopeRoot)

		t.TraverseBranch(externs, scopeRoot)
		jutil.CheckState(root.GetParent() == scopeRoot)
		t.TraverseBranch(root, scopeRoot)

		t.PopScope()
	}, func(unexpectedException jutil.Throwable) {
		t.ThrowUnexpectedException(unexpectedException)
	})
}

// Java: NodeTraversal.traverseRoots(AbstractCompiler, Callback, Node, Node)
func NodeTraversal_TraverseRoots(
	compiler AbstractCompiler, cb NodeTraversal_Callback, externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Builder_().SetCompiler(compiler).SetCallback(cb).TraverseRoots(externs, root)
}

const nodeTraversal_MISSING_SOURCE = "[source unknown]"

// Java: NodeTraversal.formatNodePosition(Node)
func (t *NodeTraversal) FormatNodePosition(n *rhino.Node) string {
	sourceFileName := t.GetBestSourceFileName(n)
	if sourceFileName == nil {
		return nodeTraversal_MISSING_SOURCE + "\n"
	}

	lineNumber := n.GetLineno()
	columnNumber := n.GetCharno()
	src := t.compiler.GetSourceLine(sourceFileName, lineNumber)
	if src == nil {
		missing := nodeTraversal_MISSING_SOURCE
		src = &missing
	}
	return *sourceFileName + ":" + strconv.Itoa(lineNumber) + ":" + strconv.Itoa(columnNumber) + "\n" + *src + "\n"
}

// TraverseWithScope traverses a parse tree recursively with a scope,
// starting with the given root. This should only be used in the global
// scope or module scopes. Otherwise, use TraverseAtScope.
//
// Java: NodeTraversal.traverseWithScope(Node, AbstractScope<?, ?>)
func (t *NodeTraversal) TraverseWithScope(root *rhino.Node, s AbstractScope) {
	jutil.CheckStateMsg1(s.IsGlobal() || s.IsModuleScope(), "%s", s)
	jutil.Try(func() {
		t.InitTraversal(root)
		t.currentNode = root
		t.PushScopeOfScope(s)
		t.TraverseBranch(root, nil)
		t.PopScope()
	}, func(unexpectedException jutil.Throwable) {
		t.ThrowUnexpectedException(unexpectedException)
	})
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
	var parentScopes []AbstractScope
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
		if t.callback.ShouldTraverse(t, n, nil) {
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
			t.callback.Visit(t, n, nil)
		}
	case rhino.Token_CLASS:
		if t.callback.ShouldTraverse(t, n, nil) {
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
			t.callback.Visit(t, n, nil)
		}
	case rhino.Token_BLOCK, rhino.Token_SWITCH_BODY:
		if t.callback.ShouldTraverse(t, n, nil) {
			t.PushScopeOfScope(s)

			// traverseBranch is not called here to avoid re-creating the block scope.
			t.TraverseChildren(n)

			t.PopScope()
			t.callback.Visit(t, n, nil)
		}
	case rhino.Token_MEMBER_FIELD_DEF:
		t.PushScopeOfScope(s)
		if t.callback.ShouldTraverse(t, n, nil) {

			// traverseBranch is not called here to avoid re-creating the MemberFieldDef scope.
			t.TraverseChildren(n)

			t.callback.Visit(t, n, nil)
		}
		t.PopScope()
	case rhino.Token_COMPUTED_FIELD_DEF:
		t.PushScopeOfScope(s)
		if t.callback.ShouldTraverse(t, n, nil) {

			t.TraverseBranch(n.GetLastChild(), n)

			t.callback.Visit(t, n, nil)
		}
		t.PopScope()
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_FOR:
		if t.callback.ShouldTraverse(t, n, nil) {
			t.PushScopeOfScope(s)

			forAssignmentParam := n.GetFirstChild()
			forIterableParam := forAssignmentParam.GetNext()
			forBodyScope := forIterableParam.GetNext()
			t.TraverseBranch(forAssignmentParam, n)
			t.TraverseBranch(forIterableParam, n)
			t.TraverseBranch(forBodyScope, n)

			t.PopScope()
			t.callback.Visit(t, n, nil)
		}
	default:
		jutil.CheckStateMsg1(
			s.IsGlobal() || s.IsModuleScope(), "Expected global or module scope. Got: (%s)", s)
		t.TraverseWithScope(n, s)
	}
}

// NodeTraversal_TraverseScopeRoots traverses *just* the contents of provided
// scope nodes (and optionally scopes nested within them)
//
// Java: NodeTraversal.traverseScopeRoots(AbstractCompiler, List<Node>, Callback, boolean)
func NodeTraversal_TraverseScopeRoots(
	compiler AbstractCompiler, scopeNodes []*rhino.Node, cb NodeTraversal_Callback, traverseNested bool) {

	scb := &nodeTraversal_TraverseScopeRootsCallback{cb: cb, traverseNested: traverseNested}
	scopeCreator := NewMemoizedScopeCreator(NewSyntacticScopeCreator(compiler))

	for _, scopeNode := range scopeNodes {
		scb.scopeNode = scopeNode
		NodeTraversal_Builder_().
			SetCompiler(compiler).
			SetCallback(scb).
			SetScopeCreator(scopeCreator).
			Build().
			TraverseScopeRoot(scopeNode)
	}
}

// nodeTraversal_TraverseScopeRootsCallback is the local class
// TraverseScopeRootsCallback of traverseScopeRoots (cb and traverseNested
// are its captured variables).
type nodeTraversal_TraverseScopeRootsCallback struct {
	insideScopeNode bool
	scopeNode       *rhino.Node

	cb             NodeTraversal_Callback
	traverseNested bool
}

// Java: NodeTraversal.traverseScopeRoots.TraverseScopeRootsCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *nodeTraversal_TraverseScopeRootsCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if c.scopeNode == n {
		c.insideScopeNode = true
	}
	return (c.traverseNested || c.scopeNode == n || !ChangeTracker_IsChangeScopeRoot(n)) &&
		c.cb.ShouldTraverse(t, n, parent)
}

// Java: NodeTraversal.traverseScopeRoots.TraverseScopeRootsCallback.visit(NodeTraversal, Node, Node)
func (c *nodeTraversal_TraverseScopeRootsCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if c.scopeNode == n {
		c.insideScopeNode = false
	}
	c.cb.Visit(t, n, parent)
}

// Java: NodeTraversal.traverseScopeRoots.TraverseScopeRootsCallback.enterScope(NodeTraversal)
func (c *nodeTraversal_TraverseScopeRootsCallback) EnterScope(t *NodeTraversal) {
	if scopedCallback, ok := c.cb.(NodeTraversal_ScopedCallback); c.insideScopeNode && ok {
		scopedCallback.EnterScope(t)
	}
}

// Java: NodeTraversal.traverseScopeRoots.TraverseScopeRootsCallback.exitScope(NodeTraversal)
func (c *nodeTraversal_TraverseScopeRootsCallback) ExitScope(t *NodeTraversal) {
	if scopedCallback, ok := c.cb.(NodeTraversal_ScopedCallback); c.insideScopeNode && ok {
		scopedCallback.ExitScope(t)
	}
}

// Java: NodeTraversal.traverseScopeRoot(Node)
func (t *NodeTraversal) TraverseScopeRoot(scopeRoot *rhino.Node) {
	jutil.Try(func() {
		t.InitTraversal(scopeRoot)
		t.currentNode = scopeRoot
		t.InitScopeRoots(scopeRoot.GetParent())
		t.TraverseBranch(scopeRoot, scopeRoot.GetParent())
	}, func(unexpectedException jutil.Throwable) {
		t.ThrowUnexpectedException(unexpectedException)
	})
}

// Java: NodeTraversal.getCompiler()
func (t *NodeTraversal) GetCompiler() AbstractCompiler {
	return t.compiler
}

// GetSourceName gets the current input source name: a string that may be
// empty, but not null.
//
// Go-only: language-adapter — a *string, since the Java String comes from
// Node.getSourceFileName(), which may be null.
//
// Java: NodeTraversal.getSourceName()
func (t *NodeTraversal) GetSourceName() *string {
	if t.sourceName == nil {
		if t.currentScript != nil {
			t.sourceName = t.currentScript.GetSourceFileName()
		} else {
			empty := ""
			t.sourceName = &empty
		}
	}
	return t.sourceName
}

// GetInput gets the current input source.
//
// Java: NodeTraversal.getInput()
func (t *NodeTraversal) GetInput() *CompilerInput {
	inputId := t.GetInputId()
	if t.compilerInput == nil && inputId != nil {
		t.compilerInput = t.compiler.GetInput(inputId)
	}
	return t.compilerInput
}

// GetChunk gets the current input chunk.
//
// Java: NodeTraversal.getChunk()
func (t *NodeTraversal) GetChunk() *JSChunk {
	input := t.GetInput()
	if input == nil {
		return nil
	}
	return input.GetChunk()
}

// GetCurrentNode returns the node currently being traversed.
//
// Java: NodeTraversal.getCurrentNode()
func (t *NodeTraversal) GetCurrentNode() *rhino.Node {
	return t.currentNode
}

// Java: NodeTraversal.handleScript(Node, Node)
func (t *NodeTraversal) HandleScript(n *rhino.Node, parent *rhino.Node) {
	if Platform_IsThreadInterrupted() {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
	}
	t.SetChangeScope(n)

	t.currentNode = n
	t.currentScript = n
	t.ClearScriptState()
	if t.callback.ShouldTraverse(t, n, parent) {
		t.TraverseChildren(n)
		t.currentNode = n
		t.callback.Visit(t, n, parent)
	}
	t.SetChangeScope(nil)
}

// Java: NodeTraversal.handleFunction(Node, Node)
func (t *NodeTraversal) HandleFunction(n *rhino.Node, parent *rhino.Node) {
	changeScope := t.currentChangeScope
	t.SetChangeScope(n)
	t.currentNode = n
	if t.callback.ShouldTraverse(t, n, parent) {
		t.TraverseFunction(n, parent)
		t.currentNode = n
		t.callback.Visit(t, n, parent)
	}
	t.SetChangeScope(changeScope)
}

// HandleModule traverses a module.
//
// Java: NodeTraversal.handleModule(Node, Node)
func (t *NodeTraversal) HandleModule(n *rhino.Node, parent *rhino.Node) {
	t.currentHoistScopeRoot = n
	t.PushScope(n)
	t.currentNode = n
	if t.callback.ShouldTraverse(t, n, parent) {
		t.TraverseChildren(n)
		t.currentNode = n
		t.callback.Visit(t, n, parent)
	}
	t.PopScope()
	// Module bodies don't nest
	t.currentHoistScopeRoot = nil
}

// Java: NodeTraversal.handleDestructuringOrDefaultValue(Node, Node)
func (t *NodeTraversal) HandleDestructuringOrDefaultValue(n *rhino.Node, parent *rhino.Node) {
	t.currentNode = n
	if t.callback.ShouldTraverse(t, n, parent) {
		first := n.GetFirstChild()
		second := first.GetNext()

		if second != nil {
			jutil.CheckStateMsg1(second.GetNext() == nil, "%s", second)
			t.TraverseBranch(second, n)
		}
		t.TraverseBranch(first, n)

		t.currentNode = n
		t.callback.Visit(t, n, parent)
	}
}

// TraverseBranch traverses a branch.
//
// Java: NodeTraversal.traverseBranch(Node, Node)
func (t *NodeTraversal) TraverseBranch(n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		t.HandleScript(n, parent)
		return
	case rhino.Token_FUNCTION:
		t.HandleFunction(n, parent)
		return
	case rhino.Token_MODULE_BODY:
		t.HandleModule(n, parent)
		return
	case rhino.Token_CLASS:
		t.HandleClass(n, parent)
		return
	case rhino.Token_CLASS_MEMBERS:
		t.HandleClassMembers(n, parent)
		return
	case rhino.Token_DEFAULT_VALUE, rhino.Token_DESTRUCTURING_LHS:
		// TODO(nickreid): Handle ASSIGN with destructuring target.
		if t.obeyDestructuringAndDefaultValueExecutionOrder {
			t.HandleDestructuringOrDefaultValue(n, parent)
			return
		}
	default:
	}

	t.currentNode = n
	if !t.callback.ShouldTraverse(t, n, parent) {
		return
	}

	createsBlockScope := NodeUtil_CreatesBlockScope(n)
	previousHoistScopeRoot := t.currentHoistScopeRoot
	if createsBlockScope {
		t.PushScope(n)
		if NodeUtil_IsClassStaticBlock(n) {
			t.currentHoistScopeRoot = n
		}
	}

	/*
	 * Intentionally inlined call to traverseChildren.
	 *
	 * <p>Calling traverseChildren here would double the maximum stack depth seen during
	 * compilation. That can cause stack overflows on some very deep ASTs (e.g. b/188616350).
	 * Inlining side-steps the issue.
	 */
	for child := n.GetFirstChild(); child != nil; {
		// child could be replaced, in which case our child node
		// would no longer point to the true next
		next := child.GetNext()
		t.TraverseBranch(child, n)
		child = next
	}

	if createsBlockScope {
		t.PopScope()
		t.currentHoistScopeRoot = previousHoistScopeRoot
	}

	t.currentNode = n
	t.callback.Visit(t, n, parent)
}

// TraverseFunction traverses a function.
//
// Java: NodeTraversal.traverseFunction(Node, Node)
func (t *NodeTraversal) TraverseFunction(n *rhino.Node, parent *rhino.Node) {
	fnName := n.GetFirstChild()
	// NOTE: If a function declaration is the root of a traversal, then we will treat it as a
	// function expression (since 'parent' is null, even though 'n' actually has a parent node) and
	// traverse the function name before entering the scope, rather than afterwards. Removing the
	// null check for 'parent' seems safe, but causes a rare crash when traverseScopeRoots is called
	// by PeepholeOptimizationsPass on a function that somehow doesn't actually have a parent at all
	// (presumably because it's already been removed from the AST?) so that doesn't actually work.
	// This does not actually change anything, though, since rooting a traversal at a function node
	// causes the function scope to be entered twice (unless using #traverseAtScope, which doesn't
	// call this method), so that the name node is just always traversed inside the scope anyway.
	isFunctionDeclaration := parent != nil && NodeUtil_IsFunctionDeclaration(n)

	if isFunctionDeclaration {
		// Function declarations are in the scope containing the declaration.
		t.TraverseBranch(fnName, n)
	}

	t.currentNode = n

	// may nest within class static blocks
	previousHoistScopeRoot := t.currentHoistScopeRoot
	t.currentHoistScopeRoot = n
	previousFunction := t.currentFunction
	t.currentFunction = n
	t.PushScope(n)

	if !isFunctionDeclaration {
		// Function expression names are only accessible within the function
		// scope.
		t.TraverseBranch(fnName, n)
	}

	args := fnName.GetNext()
	body := args.GetNext()

	// Args
	t.TraverseBranch(args, n)

	// Body
	// ES6 "arrow" function may not have a block as a body.
	t.TraverseBranch(body, n)

	t.PopScope()
	t.currentFunction = previousFunction
	t.currentHoistScopeRoot = previousHoistScopeRoot
}

// HandleClass traverses a class. Note that we traverse some of the child
// nodes slightly out of order to ensure children are visited in the correct
// scope. The following children are in the outer scope: (1) the 'extends'
// clause, (2) any computed method keys, (3) the class name for class
// declarations only (class expression names are traversed in the class
// scope). This requires that we visit the extends node (second child) and
// any computed member keys (grandchildren of the last, body, child) before
// visiting the name (first child) or body (last child).
//
// Java: NodeTraversal.handleClass(Node, Node)
func (t *NodeTraversal) HandleClass(n *rhino.Node, parent *rhino.Node) {
	t.currentNode = n
	if !t.callback.ShouldTraverse(t, n, parent) {
		return
	}

	className := n.GetFirstChild()
	extendsClause := className.GetNext()
	body := extendsClause.GetNext()

	isClassExpression := NodeUtil_IsClassExpression(n)

	t.TraverseBranch(extendsClause, n)

	if !isClassExpression {
		// Class declarations are in the scope containing the declaration.
		t.TraverseBranch(className, n)
	}

	t.currentNode = n
	t.PushScope(n)

	if isClassExpression {
		// Class expression names are only accessible within the function
		// scope.
		t.TraverseBranch(className, n)
	}

	// Body
	t.TraverseBranch(body, n)

	t.PopScope()

	t.currentNode = n
	t.callback.Visit(t, n, parent)
}

// HandleClassMembers traverses class members, excluding keys of computed
// props.
//
// Java: NodeTraversal.handleClassMembers(Node, Node)
func (t *NodeTraversal) HandleClassMembers(n *rhino.Node, parent *rhino.Node) {
	t.currentNode = n
	if !t.callback.ShouldTraverse(t, n, parent) {
		return
	}

	for child := n.GetFirstChild(); child != nil; {
		next := child.GetNext() // see traverseChildren
		if child.IsComputedProp() || child.IsComputedFieldDef() {
			t.TraverseBranch(child.GetFirstChild(), child)
		}
		child = next
	}

	for child := n.GetFirstChild(); child != nil; {
		next := child.GetNext() // see traverseChildren

		switch child.GetToken() {
		case rhino.Token_COMPUTED_PROP:
			t.currentNode = n

			if t.callback.ShouldTraverse(t, child, n) {
				t.TraverseBranch(child.GetLastChild(), child)
				t.currentNode = n
				t.callback.Visit(t, child, n)
			}
		case rhino.Token_COMPUTED_FIELD_DEF:
			t.currentNode = n

			previousHoistScopeRoot := t.currentHoistScopeRoot
			t.currentHoistScopeRoot = n
			t.PushScope(child)

			if t.callback.ShouldTraverse(t, child, n) {
				if child.HasTwoChildren() { // No RHS to traverse in `[x];` computed field case
					t.TraverseBranch(child.GetLastChild(), child)
				}
				t.currentNode = n
				t.callback.Visit(t, child, n)
			}

			t.PopScope()
			t.currentHoistScopeRoot = previousHoistScopeRoot
		case rhino.Token_MEMBER_FIELD_DEF:
			t.HandleMemberFieldDef(n, child)
		case rhino.Token_BLOCK, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_VARIABLE_DEF,
			rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
			t.TraverseBranch(child, n)
		default:
			panic(jutil.NewIllegalStateExceptionMsg("Invalid class member: " + child.GetToken().String()))
		}
		child = next
	}

	t.currentNode = n
	t.callback.Visit(t, n, parent)
}

// Java: NodeTraversal.handleMemberFieldDef(Node, Node)
func (t *NodeTraversal) HandleMemberFieldDef(n *rhino.Node, child *rhino.Node) {
	previousHoistScopeRoot := t.currentHoistScopeRoot
	t.currentHoistScopeRoot = n
	t.PushScope(child)
	t.TraverseBranch(child, n)
	t.PopScope()
	t.currentHoistScopeRoot = previousHoistScopeRoot
}

// Java: NodeTraversal.traverseChildren(Node)
func (t *NodeTraversal) TraverseChildren(n *rhino.Node) {
	for child := n.GetFirstChild(); child != nil; {
		// child could be replaced, in which case our child node
		// would no longer point to the true next
		next := child.GetNext()
		t.TraverseBranch(child, n)
		child = next
	}
}

// GetEnclosingFunction examines the functions stack for the last instance
// of a function node. When possible, prefer this method over
// NodeUtil.getEnclosingFunction() because this in general looks at less
// nodes.
//
// Java: NodeTraversal.getEnclosingFunction()
func (t *NodeTraversal) GetEnclosingFunction() *rhino.Node {
	return t.currentFunction
}

// PushScope creates a new scope (e.g. when entering a function).
//
// Java: NodeTraversal.pushScope(Node)
func (t *NodeTraversal) PushScope(node *rhino.Node) {
	if t.currentNode == nil {
		panic(jutil.NewNullPointerException()) // checkNotNull(currentNode)
	}
	if node == nil {
		panic(jutil.NewNullPointerException()) // checkNotNull(node)
	}
	t.scopes = append(t.scopes, node)
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
	if t.currentNode == nil {
		panic(jutil.NewNullPointerException()) // checkNotNull(currentNode)
	}
	t.scopes = append(t.scopes, s)
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
	t.scopes[len(t.scopes)-1] = nil
	t.scopes = t.scopes[:len(t.scopes)-1]
}

// GetNodeRootFromScopeObj: scopes are stored in a list of objects. These
// can either be instanciated `AbstractScope` or the `Node` representing the
// root of the scope.
//
// Java: NodeTraversal.getNodeRootFromScopeObj(Object)
func (t *NodeTraversal) GetNodeRootFromScopeObj(root any) *rhino.Node {
	if node, ok := root.(*rhino.Node); ok {
		return node
	}
	return root.(AbstractScope).GetRootNode()
}

// GetScopeRoot returns the current scope's root.
//
// Java: NodeTraversal.getScopeRoot()
func (t *NodeTraversal) GetScopeRoot() *rhino.Node {
	roots := len(t.scopes)
	if roots > 0 {
		return t.GetNodeRootFromScopeObj(t.scopes[roots-1])
	} else {
		return nil
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
	return t.GetAbstractScope().Untyped()
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

	o := t.scopes[rootDepth]
	if node, ok := o.(*rhino.Node); ok {
		// The root scope has a null parent.
		var parentScope AbstractScope
		if rootDepth > 0 {
			parentScope = t.GetAbstractScopeAt(rootDepth - 1)
		}
		scope := t.scopeCreator.CreateScope(node, parentScope)
		t.scopes[rootDepth] = scope
		return scope
	} else {
		return o.(AbstractScope)
	}
}

// Java: NodeTraversal.isHoistScope()
func (t *NodeTraversal) IsHoistScope() bool {
	return nodeTraversal_isHoistScopeRootNode(t.GetScopeRoot())
}

// Java: NodeTraversal.getClosestHoistScopeRoot()
func (t *NodeTraversal) GetClosestHoistScopeRoot() *rhino.Node {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		rootNode := t.GetNodeRootFromScopeObj(t.scopes[i])
		if nodeTraversal_isHoistScopeRootNode(rootNode) {
			return rootNode
		}
	}

	return nil
}

// Java: NodeTraversal.getClosestContainerScope()
func (t *NodeTraversal) GetClosestContainerScope() AbstractScope {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		rootNode := t.GetNodeRootFromScopeObj(t.scopes[i])
		if !NodeUtil_CreatesBlockScope(rootNode) {
			return t.GetAbstractScopeAt(i)
		}
	}
	return nil
}

// Java: NodeTraversal.getClosestHoistScope()
func (t *NodeTraversal) GetClosestHoistScope() AbstractScope {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		rootNode := t.GetNodeRootFromScopeObj(t.scopes[i])
		if nodeTraversal_isHoistScopeRootNode(rootNode) {
			return t.GetAbstractScopeAt(i)
		}
	}
	return nil
}

// Java: NodeTraversal.isHoistScopeRootNode(Node)
func nodeTraversal_isHoistScopeRootNode(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_MODULE_BODY, rhino.Token_ROOT, rhino.Token_SCRIPT:
		return true
	default:
		return NodeUtil_IsFunctionBlock(n)
	}
}

// GetClosestScopeRootNodeBindingThisOrSuper returns the closest scope
// binding the `this` or `super` keyword
//
// Java: NodeTraversal.getClosestScopeRootNodeBindingThisOrSuper()
func (t *NodeTraversal) GetClosestScopeRootNodeBindingThisOrSuper() *rhino.Node {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		rootNode := t.GetNodeRootFromScopeObj(t.scopes[i])
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
	return nil
}

// Java: NodeTraversal.getScopeCreator()
func (t *NodeTraversal) GetScopeCreator() ScopeCreator {
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
	return t.currentHoistScopeRoot == nil
}

// InModuleHoistScope determines whether the hoist scope of the current
// traversal is global.
//
// Java: NodeTraversal.inModuleHoistScope()
func (t *NodeTraversal) InModuleHoistScope() bool {
	moduleRoot := t.currentHoistScopeRoot
	if moduleRoot == nil {
		// in global hoist scope
		return false
	}
	if moduleRoot.IsFunction() {
		// For wrapped modules, the function block is the module scope root.
		moduleRoot = moduleRoot.GetLastChild()
	}
	return NodeUtil_IsModuleScopeRoot(moduleRoot)
}

// Report reports a diagnostic (error or warning)
//
// Java: NodeTraversal.report(Node, DiagnosticType, String...)
func (t *NodeTraversal) Report(n *rhino.Node, diagnosticType *DiagnosticType, arguments ...string) {
	error_ := JSError_MakeWithNode(n, diagnosticType, arguments...)
	t.compiler.Report(error_)
}

// ReportWithNodeRange reports a diagnostic (error or warning) This variant
// is particularly useful for reporting on GETPROP nodes whose length is
// otherwise just the last component (`abc` in `foo.abc`).
//
// Java: NodeTraversal.report(Node, Node, DiagnosticType, String...)
func (t *NodeTraversal) ReportWithNodeRange(start *rhino.Node, end *rhino.Node, diagnosticType *DiagnosticType, arguments ...string) {
	error_ := JSError_MakeWithNodeRange(start, end, diagnosticType, arguments...)
	t.compiler.Report(error_)
}

// Java: NodeTraversal.reportCodeChange()
func (t *NodeTraversal) ReportCodeChange() {
	changeScope := t.currentChangeScope
	if changeScope == nil {
		panic(jutil.NewNullPointerException()) // checkNotNull(changeScope)
	}
	jutil.CheckStateMsg1(ChangeTracker_IsChangeScopeRoot(changeScope), "%s", changeScope)
	t.compiler.ReportChangeToChangeScope(changeScope)
}

// Java: NodeTraversal.reportCodeChange(Node)
func (t *NodeTraversal) ReportCodeChangeWithNode(n *rhino.Node) {
	t.compiler.ReportChangeToEnclosingScope(n)
}

// GetCurrentScript returns the SCRIPT node enclosing the current scope.
//
// Throws UnsupportedOperationException if inside a NodeTraversal that does
// not support this operation, or if called while visiting a ROOT node during
// a global traversal. Only traversals that begin at a ROOT node (which
// includes most usages of NodeTraversal) will support this operation. For
// example, NodeTraversal.Builder#traverseAtScope(AbstractScope<?, ?>) is
// unsupported because it can begin at any arbitrary scope root like a
// FUNCTION.
//
// Java: NodeTraversal.getCurrentScript()
func (t *NodeTraversal) GetCurrentScript() *rhino.Node {
	if t.currentScript == nil {
		panic(jutil.NewUnsupportedOperationExceptionMsg("getCurrentScript not supported"))
	}
	return t.currentScript
}

// SetChangeScope: n is the current change scope, should be null when the
// traversal is complete.
//
// Java: NodeTraversal.setChangeScope(Node)
func (t *NodeTraversal) SetChangeScope(n *rhino.Node) {
	t.currentChangeScope = n
}

// Java: NodeTraversal.getEnclosingScript(Node)
func (t *NodeTraversal) GetEnclosingScript(n *rhino.Node) *rhino.Node {
	for n != nil && !n.IsScript() {
		n = n.GetParent()
	}
	return n
}

// Java: NodeTraversal.initTraversal(Node)
func (t *NodeTraversal) InitTraversal(traversalRoot *rhino.Node) {
	if Platform_IsThreadInterrupted() {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
	}
	hoistScopeRoot := NodeUtil_GetEnclosingHoistScopeRoot(traversalRoot)
	t.currentHoistScopeRoot = hoistScopeRoot
	// "null" is the global hoist scope root, but we want the current script if any
	// for the "change scope"
	changeScopeStart := traversalRoot
	if hoistScopeRoot != nil {
		changeScopeStart = hoistScopeRoot
	}
	changeScope := ChangeTracker_GetEnclosingChangeScopeRoot(changeScopeStart)
	t.SetChangeScope(changeScope)

	var enclosingFunction *rhino.Node
	if hoistScopeRoot != nil {
		enclosingFunction = NodeUtil_GetEnclosingFunction(hoistScopeRoot)
	}
	t.currentFunction = enclosingFunction

	script := t.GetEnclosingScript(changeScope)
	t.currentScript = script
	t.ClearScriptState()
}

// InitScopeRoots prefills the scopes stack up to a given spot in the AST.
// Allows for starting traversal at any spot while still having correct scope
// state.
//
// Java: NodeTraversal.initScopeRoots(Node)
func (t *NodeTraversal) InitScopeRoots(n *rhino.Node) {
	// Go-only: language-adapter — the Deque filled with addFirst: roots are
	// collected innermost first and pushed outermost first.
	var queuedScopeRoots []*rhino.Node
	for n != nil {
		if t.IsScopeRoot(n) {
			queuedScopeRoots = append(queuedScopeRoots, n)
		}
		n = n.GetParent()
	}
	for i := len(queuedScopeRoots) - 1; i >= 0; i-- {
		t.PushScope(queuedScopeRoots[i])
	}
}

// Java: NodeTraversal.isScopeRoot(Node)
func (t *NodeTraversal) IsScopeRoot(n *rhino.Node) bool {
	if n.IsRoot() && n.GetParent() == nil {
		return true
	} else if n.IsFunction() {
		return true
	} else if NodeUtil_CreatesBlockScope(n) {
		return true
	}
	return false
}

// ClearScriptState is used to clear any cached state with regard to the
// current script and should be called before traversing a SCRIPT rooted
// subtree.
//
// Java: NodeTraversal.clearScriptState()
func (t *NodeTraversal) ClearScriptState() {
	t.inputId = nil
	t.sourceName = nil
	t.compilerInput = nil
}

// Java: NodeTraversal.getInputId()
func (t *NodeTraversal) GetInputId() *rhino.InputId {
	if t.currentScript != nil && t.inputId == nil {
		t.inputId = t.currentScript.GetInputId()
	}
	return t.inputId
}

// Java: NodeTraversal.getBestSourceFileName(Node)
func (t *NodeTraversal) GetBestSourceFileName(n *rhino.Node) *string {
	if n == nil {
		return t.GetSourceName()
	}
	return n.GetSourceFileName()
}

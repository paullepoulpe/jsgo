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

// Porting compiler passes against NodeTraversal (read this first).
//
// Callbacks. NodeTraversal.Callback is the Go interface
// NodeTraversal_Callback (ShouldTraverse + Visit); ScopedCallback is
// NodeTraversal_ScopedCallback (adds EnterScope + ExitScope). As in Java, a
// traversal fires enterScope/exitScope only if its callback is a
// NodeTraversal_ScopedCallback, which in Go means "has EnterScope and
// ExitScope methods".
//
// Abstract callback bases. A pass that extends a stateless abstract callback
// embeds its base struct and defines the abstract methods:
//
//	// Java: final class MyPass extends NodeTraversal.AbstractPostOrderCallback
//	type MyPass struct {
//		NodeTraversal_AbstractPostOrderCallbackBase
//		compiler AbstractCompiler
//	}
//	func (p *MyPass) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) { ... }
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
//
// Passes needing lazy per-scope CFGs embed cfgScopeStack directly; passes
// tracking modules embed moduleTraversal directly. Java's final methods (e.g.
// ShouldTraverse of AbstractPostOrderCallback) must not be redefined by the
// subclass: Go cannot forbid it, and a redefinition would silently change the
// traversal.
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

type nodeTraversal_noShouldTraverse interface {
	nodeTraversalNoShouldTraverse()
}

type nodeTraversal_onlyScriptAndFunctionShouldTraverse interface {
	shouldTraverseOnlyScriptAndFunction() bool
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

// NodeTraversal_AbstractPostOrderCallbackBase holds the methods of
// AbstractPostOrderCallback; subclasses embed it and define Visit.
type NodeTraversal_AbstractPostOrderCallbackBase struct{}

// ShouldTraverse is final in Java.
//
// Java: NodeTraversal.AbstractPostOrderCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractPostOrderCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

func (NodeTraversal_AbstractPostOrderCallbackBase) nodeTraversalNoShouldTraverse() {}

// NodeTraversal_ExternsSkippingCallbackBase holds the methods of
// ExternsSkippingCallback; subclasses embed it and define Visit.
// Note: Even though type-summary nodes are included under the externs roots,
// they are traversed by this callback.
type NodeTraversal_ExternsSkippingCallbackBase struct{}

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

var nodeTraversal_nilPostOrderCallback NodeTraversal_Callback = &nodeTraversal_Builder_postOrderCallback{}

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

// NodeTraversal_AbstractPreOrderCallbackBase holds the methods of
// AbstractPreOrderCallback; subclasses embed it and define ShouldTraverse.
type NodeTraversal_AbstractPreOrderCallbackBase struct{}

// Visit is final in Java.
//
// Java: NodeTraversal.AbstractPreOrderCallback.visit(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractPreOrderCallbackBase) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
}

func (NodeTraversal_AbstractPreOrderCallbackBase) nodeTraversalNoVisit() {}

// NodeTraversal_AbstractScopedCallbackBase holds the methods of
// AbstractScopedCallback; subclasses embed it, define Visit and may
// redefine EnterScope and ExitScope.
type NodeTraversal_AbstractScopedCallbackBase struct{}

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

// NodeTraversal_AbstractShallowCallbackBase holds the methods of
// AbstractShallowCallback (visits all nodes without traversing into function
// bodies); subclasses embed it and define Visit.
type NodeTraversal_AbstractShallowCallbackBase struct{}

// ShouldTraverse is final in Java.
//
// Java: NodeTraversal.AbstractShallowCallback.shouldTraverse(NodeTraversal, Node, Node)
func (b *NodeTraversal_AbstractShallowCallbackBase) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// We do want to traverse the name of a named function, but we don't
	// want to traverse the arguments or body.
	return parent.IsNil() || !parent.IsFunction() || n == parent.GetFirstChild()
}

// NodeTraversal_AbstractShallowStatementCallbackBase holds the methods of
// AbstractShallowStatementCallback (visits all structure and statement nodes
// without traversing into functions or expressions); subclasses embed it and
// define Visit.
type NodeTraversal_AbstractShallowStatementCallbackBase struct{}

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
		m.currentModule = m.moduleMetadataMap.ModuleByPath(t.GetInput().GetPath().String())
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
			newModule := m.moduleMetadataMap.ModuleByGoogNamespace(n.GetLastChild().GetString())
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
			m.currentModule = m.moduleMetadataMap.ModuleByPath(t.GetInput().GetPath().String())
			if m.currentModule == nil {
				failNodeTraversalNull()
			}
		}
	default:
	}
	return exitedModule, exitedScopeRoot
}

// NodeTraversal_TraverseScopeRoots traverses *just* the contents of provided
// scope nodes (and optionally scopes nested within them)
//
// Java: NodeTraversal.traverseScopeRoots(AbstractCompiler, List<Node>, Callback, boolean)
func NodeTraversal_TraverseScopeRoots(
	compiler AbstractCompiler, scopeNodes []rhino.Node, cb NodeTraversal_Callback, traverseNested bool) {
	if len(scopeNodes) == 0 {
		return
	}

	b := NodeTraversal_Builder{
		compiler:     compiler,
		callback:     cb,
		scopeCreator: NewMemoizedScopeCreator(NewSyntacticScopeCreator(compiler)),
	}
	t := acquireNodeTraversal(&b)
	defer releaseNodeTraversal(t)
	t.stopAtNestedChangeScope = !traverseNested
	origNoShouldTraverse := t.noShouldTraverse
	origNoVisit := t.noVisit
	origVisitTokens := t.visitTokens
	for _, scopeNode := range scopeNodes {
		t.noShouldTraverse = origNoShouldTraverse
		t.noVisit = origNoVisit
		t.visitTokens = origVisitTokens
		t.scopeTraversalRoot = scopeNode
		t.traverseScopeRoot(scopeNode, true)
	}
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

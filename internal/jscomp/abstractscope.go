// Port of com/google/javascript/jscomp/AbstractScope.java
// @ closure-compiler v20260415
//
// AbstractScope<S,V> is the generic base for Scope and TypedScope in Java.
// Go cannot express the same recursive generic bound, so we represent it as
// an interface (AbstractScope) together with a concrete base struct
// (AbstractScopeBase) that captures the shared fields and provides the
// methods that are identical in Scope and TypedScope.
//
// Java inheritance chain:
//   AbstractScope<S,V> ← Scope (S=Scope,V=Var)
//   AbstractScope<S,V> ← TypedScope (S=TypedScope,V=TypedVar)
//
// Go embedding chain (mirrored):
//   AbstractScopeBase is embedded in UntypedScope (Scope.java port).
//   TypedScope (scope.go) carries equivalent fields directly; a future
//   refactor may embed AbstractScopeBase there too.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ImplicitVar enumerates the implicit variables that are implicitly defined
// in every function scope without being explicitly declared.
// Mirrors AbstractScope.ImplicitVar in Java.
type ImplicitVar int

const (
	// ImplicitVarArguments is the "arguments" object implicit in functions.
	ImplicitVarArguments ImplicitVar = iota
	// ImplicitVarExports is the module "exports" implicit in goog.module scopes.
	ImplicitVarExports
	// ImplicitVarSuper is the "super" keyword implicit in class-body scopes.
	ImplicitVarSuper
	// ImplicitVarThis is the "this" keyword implicit in function/class scopes.
	ImplicitVarThis
)

// implicitVarName returns the JavaScript name for an ImplicitVar value.
// Mirrors ImplicitVar.name in Java.
func implicitVarName(v ImplicitVar) string {
	switch v {
	case ImplicitVarArguments:
		return "arguments"
	case ImplicitVarExports:
		return "exports"
	case ImplicitVarSuper:
		return "super"
	case ImplicitVarThis:
		return "this"
	}
	return ""
}

// implicitVarOf maps a JavaScript name to its ImplicitVar constant, or
// returns -1 if the name is not an implicit variable.
// Mirrors ImplicitVar.of(String) in Java.
func implicitVarOf(name string) (ImplicitVar, bool) {
	switch name {
	case "arguments":
		return ImplicitVarArguments, true
	case "super":
		return ImplicitVarSuper, true
	case "this":
		return ImplicitVarThis, true
	case "exports":
		return ImplicitVarExports, true
	}
	return 0, false
}

// implicitVarIsMadeByScope reports whether this kind of implicit variable
// is created and owned by the given scope root node.
// Mirrors ImplicitVar.isMadeByScope(AbstractScope) in Java.
func implicitVarIsMadeByScope(v ImplicitVar, rootNode *rhino.Node) bool {
	if rootNode == nil {
		return false
	}
	switch v {
	case ImplicitVarExports:
		// exports is implicit only in a goog.module scope.
		return rootNode.Token() == rhino.MODULE_BODY
	case ImplicitVarSuper, ImplicitVarThis:
		// super and this are implicit in non-arrow functions and class static blocks.
		tok := rootNode.Token()
		return tok == rhino.FUNCTION || tok == rhino.CLASS
	case ImplicitVarArguments:
		// arguments is implicit in non-arrow functions only.
		return rootNode.Token() == rhino.FUNCTION
	}
	return false
}

// AbstractScope is the interface satisfied by both UntypedScope (port of
// Java's Scope) and the existing TypedScope. It provides the root-node and
// parent-scope accessors that are defined final on AbstractScope<S,V> in Java.
// Mirrors the StaticScope/AbstractScope contract in Java.
type AbstractScope interface {
	// GetRootNode returns the AST node that is the root of this scope.
	// Mirrors AbstractScope.getRootNode() in Java.
	GetRootNode() *rhino.Node

	// GetParentScope returns the parent AbstractScope, or nil for global.
	// Mirrors AbstractScope.getParentScope() in Java.
	GetParentScope() AbstractScope

	// IsGlobal reports whether this scope has no parent.
	// Mirrors AbstractScope.isGlobal() in Java.
	IsGlobal() bool

	// IsLocal reports whether this scope has a parent.
	// Mirrors AbstractScope.isLocal() in Java.
	IsLocal() bool

	// GetDepth returns the nesting depth (0 = global).
	// Mirrors AbstractScope.getDepth() in Java.
	GetDepth() int

	// IsBlockScope reports whether the root is a block-creating node.
	// Mirrors AbstractScope.isBlockScope() in Java.
	IsBlockScope() bool

	// IsFunctionScope reports whether the root is a FUNCTION node.
	// Mirrors AbstractScope.isFunctionScope() in Java.
	IsFunctionScope() bool

	// IsModuleScope reports whether the root is a MODULE_BODY node.
	// Mirrors AbstractScope.isModuleScope() in Java.
	IsModuleScope() bool

	// IsFunctionBlockScope reports whether the root is the body block of
	// a FUNCTION node (i.e. the directly-enclosing block inside a function).
	// Mirrors AbstractScope.isFunctionBlockScope() in Java.
	IsFunctionBlockScope() bool

	// IsHoistScope reports whether var declarations in this scope would be
	// hoisted to this scope (function, function-block, global, module, or
	// static-block scopes).
	// Mirrors AbstractScope.isHoistScope() in Java.
	IsHoistScope() bool
}

// AbstractScopeBase implements the shared fields and logic for both
// UntypedScope (Java's Scope) and could eventually be embedded in TypedScope.
// It stores the rootNode (from Java AbstractScope) and provides most of the
// "final" methods in Java's AbstractScope that don't need generics.
//
// Not embedded in the existing TypedScope to avoid a disruptive refactor, but
// provides a parallel implementation that satisfies AbstractScope.
type AbstractScopeBase struct {
	rootNode *rhino.Node
}

// GetRootNode returns the scope's root AST node.
// Mirrors AbstractScope.getRootNode() (final) in Java.
func (b *AbstractScopeBase) GetRootNode() *rhino.Node { return b.rootNode }

// IsBlockScope reports whether the root creates a block scope.
// Mirrors AbstractScope.isBlockScope() (final) in Java via NodeUtil.createsBlockScope.
func (b *AbstractScopeBase) IsBlockScope() bool {
	return nodeCreatesBlockScope(b.rootNode)
}

// IsFunctionBlockScope reports whether the root is a function body block.
// Mirrors AbstractScope.isFunctionBlockScope() (final) in Java via
// NodeUtil.isFunctionBlock.
func (b *AbstractScopeBase) IsFunctionBlockScope() bool {
	if b.rootNode == nil {
		return false
	}
	// A function body block: BLOCK whose parent is FUNCTION.
	if b.rootNode.Token() != rhino.BLOCK {
		return false
	}
	p := b.rootNode.Parent()
	return p != nil && p.Token() == rhino.FUNCTION
}

// IsFunctionScope reports whether the root is a FUNCTION node.
// Mirrors AbstractScope.isFunctionScope() (final) in Java.
func (b *AbstractScopeBase) IsFunctionScope() bool {
	return b.rootNode != nil && b.rootNode.Token() == rhino.FUNCTION
}

// IsModuleScope reports whether the root is a MODULE_BODY node.
// Mirrors AbstractScope.isModuleScope() (final) in Java.
func (b *AbstractScopeBase) IsModuleScope() bool {
	return b.rootNode != nil && b.rootNode.Token() == rhino.MODULE_BODY
}

// IsHoistScope reports whether var declarations within this scope are
// hoisted to this scope (not to an enclosing one).
// Mirrors AbstractScope.isHoistScope() in Java.
func (b *AbstractScopeBase) IsHoistScope() bool {
	return b.IsFunctionScope() ||
		b.IsFunctionBlockScope() ||
		b.rootNode == nil || // global: no parent
		b.IsModuleScope()
}

// nodeCreatesBlockScope mirrors NodeUtil.createsBlockScope for the minimal set
// of tokens we need.  The full Java version handles a wider token set; this is
// the subset relevant to ES2015+ block scoping.
func nodeCreatesBlockScope(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.BLOCK,
		rhino.FOR, rhino.FOR_IN, rhino.FOR_OF,
		rhino.SWITCH:
		return true
	}
	return false
}

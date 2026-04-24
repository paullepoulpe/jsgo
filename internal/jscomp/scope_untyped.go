// Port of com/google/javascript/jscomp/Scope.java
// @ closure-compiler v20260415
//
// Scope is the concrete untyped scope that extends AbstractScope<Scope,Var>.
// It stores variable declarations as UntypedVar instances keyed by name.
//
// Java class: public final class Scope extends AbstractScope<Scope, Var>
//
// Go port strategy:
//   UntypedScope embeds AbstractScopeBase (which holds rootNode and the
//   scope-predicate methods) and adds the parent pointer, depth, and a
//   vars map keyed by name. The name "UntypedScope" avoids collision with
//   the existing TypedScope in scope.go and the SymbolScope in symboltable.go.
//
// Factory functions CreateGlobalUntypedScope / CreateChildUntypedScope mirror
// the Java static factory methods createGlobalScope / createChildScope.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// UntypedScope is the Go port of com/google/javascript/jscomp/Scope.java.
// It holds variables declared in one syntactic scope of a JavaScript program,
// without type information (that is the job of TypedScope in scope.go).
type UntypedScope struct {
	AbstractScopeBase

	parent *UntypedScope
	depth  int

	// vars holds variables declared in this scope.
	// Starts as nil (empty ImmutableMap.of() sentinel in Java).
	vars map[string]*UntypedVar

	// implicitVars lazily caches implicit variables (arguments, this, super, exports).
	implicitVars map[ImplicitVar]*UntypedVar
}

// CreateGlobalUntypedScope creates a root scope at the given SCRIPT or ROOT node.
// Mirrors Scope.createGlobalScope(Node) in Java.
func CreateGlobalUntypedScope(rootNode *rhino.Node) *UntypedScope {
	return &UntypedScope{
		AbstractScopeBase: AbstractScopeBase{rootNode: rootNode},
		parent:            nil,
		depth:             0,
	}
}

// CreateChildUntypedScope creates a child scope under parent at rootNode.
// Mirrors Scope.createChildScope(Scope, Node) in Java.
func CreateChildUntypedScope(parent *UntypedScope, rootNode *rhino.Node) *UntypedScope {
	return &UntypedScope{
		AbstractScopeBase: AbstractScopeBase{rootNode: rootNode},
		parent:            parent,
		depth:             parent.depth + 1,
	}
}

// ---------------------------------------------------------------------------
// AbstractScope interface implementation
// ---------------------------------------------------------------------------

// GetParentScope satisfies AbstractScope.
// Mirrors AbstractScope.getParentScope() / getParent() in Java.
func (s *UntypedScope) GetParentScope() AbstractScope {
	if s.parent == nil {
		return nil
	}
	return s.parent
}

// GetParent returns the typed parent pointer (nil for global).
func (s *UntypedScope) GetParent() *UntypedScope { return s.parent }

// IsGlobal reports whether this scope has no parent.
// Mirrors AbstractScope.isGlobal() in Java.
func (s *UntypedScope) IsGlobal() bool { return s.parent == nil }

// IsLocal reports whether this scope has a parent.
// Mirrors AbstractScope.isLocal() in Java.
func (s *UntypedScope) IsLocal() bool { return s.parent != nil }

// GetDepth returns the nesting depth (0 = global).
// Mirrors Scope.getDepth() in Java.
func (s *UntypedScope) GetDepth() int { return s.depth }

// IsHoistScope is delegated to AbstractScopeBase, which re-computes from the
// root node. Overriding here for the interface check; the embedded base
// provides the implementation.
func (s *UntypedScope) IsHoistScope() bool { return s.AbstractScopeBase.IsHoistScope() }

// ---------------------------------------------------------------------------
// Variable declaration and lookup
// ---------------------------------------------------------------------------

// Declare adds a variable to this scope.
// Mirrors Scope.declare(String, Node, CompilerInput) in Java.
// The CompilerInput parameter is elided since it is not yet ported.
func (s *UntypedScope) Declare(name string, nameNode *rhino.Node) *UntypedVar {
	if s.vars == nil {
		s.vars = make(map[string]*UntypedVar, 4)
	}
	idx := len(s.vars)
	v := newUntypedVar(name, nameNode, s, idx)
	s.vars[name] = v
	return v
}

// DeclareImplicitGoogNamespaceIfAbsent declares an implicit goog namespace
// (from goog.provide / legacy goog.module) in this scope if not already
// declared. Mirrors Scope.declareImplicitGoogNamespaceIfAbsent in Java.
func (s *UntypedScope) DeclareImplicitGoogNamespaceIfAbsent(name string, definition *rhino.Node) *UntypedVar {
	if existing := s.GetOwnSlot(name); existing != nil {
		if existing.IsImplicitGoogNamespace() {
			existing.implicitGoogNamespace = definition // update definition
		}
		return existing
	}
	if s.vars == nil {
		s.vars = make(map[string]*UntypedVar, 4)
	}
	v := createUntypedImplicitGoogNamespace(name, s, definition)
	s.vars[name] = v
	return v
}

// GetOwnSlot returns the variable with the given name declared directly in
// this scope, or nil. Does not walk parent scopes.
// Mirrors AbstractScope.getOwnSlot(String) in Java.
func (s *UntypedScope) GetOwnSlot(name string) *UntypedVar {
	if s.vars != nil {
		if v, ok := s.vars[name]; ok {
			return v
		}
	}
	// Check implicit variables.
	iv, ok := implicitVarOf(name)
	if !ok {
		return nil
	}
	if !implicitVarIsMadeByScope(iv, s.rootNode) {
		return nil
	}
	return s.getOrMakeImplicitVar(iv)
}

// GetSlot returns the variable with the given name from this or any parent
// scope, or nil.
// Mirrors AbstractScope.getSlot(String) / getVar(String) in Java.
func (s *UntypedScope) GetSlot(name string) *UntypedVar {
	for sc := s; sc != nil; sc = sc.parent {
		if v := sc.GetOwnSlot(name); v != nil {
			return v
		}
	}
	return nil
}

// HasOwnSlot reports whether this scope directly declares the given name.
// Mirrors AbstractScope.hasOwnSlot(String) in Java.
func (s *UntypedScope) HasOwnSlot(name string) bool {
	return s.GetOwnSlot(name) != nil
}

// HasSlot reports whether the given name is declared in this or any parent scope.
// Mirrors AbstractScope.hasSlot(String) in Java.
func (s *UntypedScope) HasSlot(name string) bool {
	return s.GetSlot(name) != nil
}

// RemoveVar removes the variable with the given name from this scope.
// Mirrors Scope.undeclare(Var) in Java. Used by SyntacticScopeCreator when
// an implicit goog namespace is superseded by an explicit declaration.
func (s *UntypedScope) RemoveVar(name string) {
	if s.vars != nil {
		delete(s.vars, name)
	}
}

// GetVarIterable returns all non-implicit vars declared in this scope.
// Mirrors AbstractScope.getVarIterable() in Java.
func (s *UntypedScope) GetVarIterable() []*UntypedVar {
	out := make([]*UntypedVar, 0, len(s.vars))
	for _, v := range s.vars {
		out = append(out, v)
	}
	return out
}

// GetVarCount returns the number of explicitly-declared variables.
// Mirrors AbstractScope.getVarCount() in Java.
func (s *UntypedScope) GetVarCount() int {
	return len(s.vars)
}

// GetGlobalScope walks up to and returns the global scope.
// Mirrors AbstractScope.getGlobalScope() in Java.
func (s *UntypedScope) GetGlobalScope() *UntypedScope {
	result := s
	for result.parent != nil {
		result = result.parent
	}
	return result
}

// GetClosestHoistScope returns the nearest hoist scope at or above this scope.
// Mirrors AbstractScope.getClosestHoistScope() in Java.
func (s *UntypedScope) GetClosestHoistScope() *UntypedScope {
	for cur := s; cur != nil; cur = cur.parent {
		if cur.IsHoistScope() {
			return cur
		}
	}
	return nil
}

// getOrMakeImplicitVar lazily creates the implicit variable for the given
// ImplicitVar kind and caches it. Mirrors AbstractScope.getOwnImplicitSlot
// in Java.
func (s *UntypedScope) getOrMakeImplicitVar(iv ImplicitVar) *UntypedVar {
	if s.implicitVars == nil {
		s.implicitVars = make(map[ImplicitVar]*UntypedVar)
	}
	if v, ok := s.implicitVars[iv]; ok {
		return v
	}
	v := newUntypedVar(implicitVarName(iv), nil, s, -1)
	s.implicitVars[iv] = v
	return v
}

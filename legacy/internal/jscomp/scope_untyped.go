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
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
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

// Go-only: language-adapter — CreateGlobalUntypedScope is the Go name for
// Scope.createGlobalScope(Node); the "Untyped" prefix distinguishes it from
// the TypedScope factory CreateGlobalTypedScope.
func CreateGlobalUntypedScope(rootNode *rhino.Node) *UntypedScope {
	return &UntypedScope{
		AbstractScopeBase: AbstractScopeBase{rootNode: rootNode},
		parent:            nil,
		depth:             0,
	}
}

// Go-only: language-adapter — CreateChildUntypedScope is the Go name for
// Scope.createChildScope(Scope, Node); the "Untyped" prefix avoids collision
// with the TypedScope child-creation path.
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

// Go-only: language-adapter — GetParentScope returns AbstractScope so UntypedScope
// satisfies the AbstractScope interface; Java's getParentScope() is final on
// AbstractScope and delegates to the concrete getParent().
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

// Go-only: language-adapter — IsLocal is inherited from AbstractScope in Java (final method).
func (s *UntypedScope) IsLocal() bool { return s.parent != nil }

// GetDepth returns the nesting depth (0 = global).
// Mirrors Scope.getDepth() in Java.
func (s *UntypedScope) GetDepth() int { return s.depth }

// Go-only: language-adapter — IsHoistScope delegates to AbstractScopeBase; in Java this
// is a final method on AbstractScope, but Go interface dispatch requires explicit delegation.
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

// Go-only: language-adapter — GetSlot delegates to getSlot/getVar which is
// a final method on AbstractScope; Go requires explicit implementation since
// UntypedScope cannot inherit AbstractScope's generic getSlot(String).
func (s *UntypedScope) GetSlot(name string) *UntypedVar {
	for sc := s; sc != nil; sc = sc.parent {
		if v := sc.GetOwnSlot(name); v != nil {
			return v
		}
	}
	return nil
}

// Go-only: language-adapter — HasOwnSlot is a final method on AbstractScope in Java;
// Go requires explicit delegation since UntypedScope cannot inherit it.
func (s *UntypedScope) HasOwnSlot(name string) bool {
	return s.GetOwnSlot(name) != nil
}

// Go-only: language-adapter — HasSlot is a final method on AbstractScope in Java;
// Go requires explicit delegation since UntypedScope cannot inherit it.
func (s *UntypedScope) HasSlot(name string) bool {
	return s.GetSlot(name) != nil
}

// Go-only: language-adapter — RemoveVar is the Go name for Scope.undeclare(Var);
// it takes a string rather than a Var pointer since UntypedVar management
// differs from Java's final undeclare(V var) which checks getScope() == this.
func (s *UntypedScope) RemoveVar(name string) {
	if s.vars != nil {
		delete(s.vars, name)
	}
}

// Go-only: language-adapter — GetVarIterable returns a slice instead of Java's Iterable<V>;
// Java's AbstractScope.getVarIterable() returns Iterable backed by a LinkedHashMap.
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

// Go-only: language-adapter — GetGlobalScope returns *UntypedScope instead of the
// Java generic S type; Java's AbstractScope.getGlobalScope() returns S via thisScope().
func (s *UntypedScope) GetGlobalScope() *UntypedScope {
	result := s
	for result.parent != nil {
		result = result.parent
	}
	return result
}

// Go-only: language-adapter — GetClosestHoistScope returns *UntypedScope instead of
// the Java generic S; Java's AbstractScope.getClosestHoistScope() returns S via thisScope().
func (s *UntypedScope) GetClosestHoistScope() *UntypedScope {
	for cur := s; cur != nil; cur = cur.parent {
		if cur.IsHoistScope() {
			return cur
		}
	}
	return nil
}

// Go-only: language-adapter — getOrMakeImplicitVar combines Java's
// getOwnImplicitSlot(ImplicitVar) and makeImplicitVar(ImplicitVar) into
// a single lazy-init helper; Java uses computeIfAbsent on the implicitVars map.
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

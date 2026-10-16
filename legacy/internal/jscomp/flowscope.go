// Port of com/google/javascript/jscomp/type/FlowScope.java
// @ closure-compiler v20260415 (narrow; flat-map representation;
// concrete implementation in linkedflowscope.go).
//
// FlowScope.java is a Java interface (extends StaticTypedScope and
// LatticeElement). In Go, the interface is collapsed into the concrete
// FlowScope struct (the Go equivalent of LinkedFlowScope) since Go does
// not need the interface layer for the passes currently ported. The
// interface contract methods are implemented as methods on *FlowScope.
//
// A FlowScope maps names to their flow-refined JSType at a given
// program point. It's the lattice element passed around by
// DataFlowAnalysis. Missing names fall through to the declared type
// from the surrounding TypedScope.

package jscomp

import "github.com/paullepoulpe/jsgo/legacy/internal/rhino/jstype"

// FlowScope is immutable by convention — callers always construct a
// new FlowScope via Bind / BindAll / inferSlotType / JoinFlowScopes
// rather than mutating an existing one. It corresponds to
// LinkedFlowScope in the Java source; FlowScope.java is the interface
// it implements (see linkedflowscope.go for the concrete methods).
type FlowScope struct {
	types          map[string]jstype.JSType
	syntacticScope *TypedScope // mirrors LinkedFlowScope.syntacticScope
}

// Go-only: language-adapter — Get returns the flow-refined type for `name`,
// or nil if the name isn't refined (callers should fall back to the
// declaring scope).
func (f *FlowScope) Get(name string) jstype.JSType {
	if f == nil {
		return nil
	}
	return f.types[name]
}

// inferSlotType returns a new FlowScope with `name` refined to `t`.
// Mirrors FlowScope.inferSlotType (the interface method): functional
// update, receiver unchanged. When `t` is nil the name is removed.
func (f *FlowScope) inferSlotType(name string, t jstype.JSType) *FlowScope {
	return f.Bind(name, t)
}

// inferQualifiedSlot narrows the flow type for a qualified name.
// `name` is the dotted qualified-name key (e.g. "obj.prop"). When
// `inferredType` is nil the name is dropped from the scope. Delegates
// to inferSlotType. Mirrors FlowScope.inferQualifiedSlot (the interface
// method; narrow: global-scope guard and lazy declaration are omitted).
func (f *FlowScope) inferQualifiedSlot(name string, inferredType jstype.JSType) *FlowScope {
	return f.inferSlotType(name, inferredType)
}

// getDeclarationScope returns the underlying TypedScope for this flow
// scope. Mirrors FlowScope.getDeclarationScope (the interface method).
func (f *FlowScope) getDeclarationScope() *TypedScope {
	if f == nil {
		return nil
	}
	return f.syntacticScope
}

// withSyntacticScope returns a new FlowScope anchored to the given
// syntactic scope. Mirrors FlowScope.withSyntacticScope (the interface
// method).
func (f *FlowScope) withSyntacticScope(scope *TypedScope) *FlowScope {
	out := createChildFlowScope(f)
	out.syntacticScope = scope
	return out
}

// Port of com/google/javascript/jscomp/type/FlowScope.java +
// com/google/javascript/jscomp/type/LinkedFlowScope.java
// @ closure-compiler v20260415 (narrow; flat-map representation with
// lattice parent/child support, JoinOp, inferSlotType,
// inferQualifiedSlot, optimize; no symbol-source tracking).
//
// A FlowScope maps names to their flow-refined JSType at a given
// program point. It's the lattice element passed around by
// DataFlowAnalysis. Missing names fall through to the declared type
// from the surrounding TypedScope — so a FlowScope only has to carry
// entries for names whose type differs from their declaration.
//
// Representation: a flat map (not a linked parent chain) so that all
// existing callers that access the map directly continue to work.
// createChildFlowScope copies the parent map into the child, which
// is semantically equivalent to a parent chain with O(n) copy cost
// instead of O(depth) lookup — acceptable given scope sizes here.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino/jstype"

// FlowScope is immutable by convention — callers always construct a
// new FlowScope via Bind / BindAll / inferSlotType / JoinFlowScopes
// rather than mutating an existing one.
type FlowScope struct {
	types map[string]jstype.JSType
}

// EmptyFlowScope is the bottom of the lattice — no names refined.
func EmptyFlowScope() *FlowScope {
	return &FlowScope{types: map[string]jstype.JSType{}}
}

// createChildFlowScope returns a new FlowScope that inherits all
// refinements from `parent`. Mirrors the LinkedFlowScope lattice
// parent/child relationship; here we copy eagerly (flat map) so
// the child is independent and all callers continue to work without
// change.
func createChildFlowScope(parent *FlowScope) *FlowScope {
	if parent == nil {
		return EmptyFlowScope()
	}
	out := &FlowScope{types: make(map[string]jstype.JSType, len(parent.types))}
	for k, v := range parent.types {
		out.types[k] = v
	}
	return out
}

// Get returns the flow-refined type for `name`, or nil if the name
// isn't refined (callers should fall back to the declaring scope).
func (f *FlowScope) Get(name string) jstype.JSType {
	if f == nil {
		return nil
	}
	return f.types[name]
}

// inferSlotType returns a new FlowScope with `name` refined to `t`.
// Mirrors LinkedFlowScope.inferSlotType: functional update, receiver
// unchanged. When `t` is nil the name is removed from the scope.
func (f *FlowScope) inferSlotType(name string, t jstype.JSType) *FlowScope {
	return f.Bind(name, t)
}

// inferQualifiedSlot narrows the flow type for a qualified name.
// `name` is the dotted qualified-name key (e.g. "obj.prop"). When
// `inferredType` is nil the name is dropped from the scope. Delegates
// to inferSlotType. Mirrors LinkedFlowScope.inferQualifiedSlot (narrow:
// global-scope guard and lazy declaration are omitted — jsgo's
// TypedScopeCreator handles declarations up-front).
func (f *FlowScope) inferQualifiedSlot(name string, inferredType jstype.JSType) *FlowScope {
	return f.inferSlotType(name, inferredType)
}

// optimize returns a compacted FlowScope by detecting when `f` is
// equal to its semantic parent (i.e. adds no new refinements). In the
// flat-map representation we can only compact if we have a reference to
// a parent to compare against; callers that want compaction pass the
// parent explicitly. If `parent` is nil or unequal, `f` is returned
// unchanged. Mirrors the "trim consecutive identical scopes" logic in
// Java's LinkedFlowScope.
func (f *FlowScope) optimize(parent *FlowScope) *FlowScope {
	if f == nil || parent == nil {
		return f
	}
	if f.Equal(parent) {
		return parent
	}
	return f
}

// Bind produces a copy of `f` with `name` remapped to `t`. The
// receiver is never mutated — mirrors the persistent-map contract
// LinkedFlowScope documents in Java.
func (f *FlowScope) Bind(name string, t jstype.JSType) *FlowScope {
	if f == nil {
		f = EmptyFlowScope()
	}
	out := &FlowScope{types: make(map[string]jstype.JSType, len(f.types)+1)}
	for k, v := range f.types {
		out.types[k] = v
	}
	if t != nil {
		out.types[name] = t
	} else {
		delete(out.types, name)
	}
	return out
}

// BindAll layers every entry of `refs` onto `f`.
func (f *FlowScope) BindAll(refs map[string]jstype.JSType) *FlowScope {
	if len(refs) == 0 {
		return f
	}
	out := &FlowScope{types: make(map[string]jstype.JSType, len(f.types)+len(refs))}
	if f != nil {
		for k, v := range f.types {
			out.types[k] = v
		}
	}
	for k, v := range refs {
		if v == nil {
			delete(out.types, k)
			continue
		}
		out.types[k] = v
	}
	return out
}

// Equal compares two FlowScopes by content. Nil scopes are treated
// as empty.
func (f *FlowScope) Equal(other *FlowScope) bool {
	a, b := f, other
	if a == nil {
		a = EmptyFlowScope()
	}
	if b == nil {
		b = EmptyFlowScope()
	}
	if len(a.types) != len(b.types) {
		return false
	}
	for k, va := range a.types {
		vb, ok := b.types[k]
		if !ok {
			return false
		}
		if va == nil || vb == nil {
			if va != vb {
				return false
			}
			continue
		}
		if va.TypeName() != vb.TypeName() {
			return false
		}
	}
	return true
}

// JoinOp is the stateful lattice-join accumulator. Mirrors
// LinkedFlowScope.FlowScopeJoinOp: callers invoke JoinFlow for each
// incoming-edge scope, then call Finish to get the joined result.
// Implements the FlowJoiner pattern from Java's DataFlowAnalysis.
type JoinOp struct {
	result *FlowScope
}

// NewJoinOp returns an empty JoinOp ready to accumulate edge scopes.
func NewJoinOp() *JoinOp { return &JoinOp{} }

// JoinFlow merges `input` into the accumulator. Mirrors
// FlowScopeJoinOp.joinFlow: the first call seeds the result; subsequent
// calls widen it via JoinFlowScopes.
func (j *JoinOp) JoinFlow(input *FlowScope) {
	if j.result == nil {
		j.result = input
		return
	}
	j.result = JoinFlowScopes(j.result, input)
}

// Finish returns the accumulated join result, or EmptyFlowScope if no
// inputs were joined. Mirrors FlowScopeJoinOp.finish.
func (j *JoinOp) Finish() *FlowScope {
	if j.result == nil {
		return EmptyFlowScope()
	}
	return j.result
}

// JoinFlowScopes computes the lattice-join of two FlowScopes: for
// every name refined in either input, the output refines to the
// union of the two sides. Names refined on only one branch drop back
// to their declared type on the join — represented by omitting them
// from the merged map. Mirrors LinkedFlowScope.joinOp.apply.
func JoinFlowScopes(a, b *FlowScope) *FlowScope {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	out := &FlowScope{types: map[string]jstype.JSType{}}
	for k, va := range a.types {
		vb, ok := b.types[k]
		if !ok {
			// Missing on one side → falls back to declared elsewhere;
			// no refinement survives the join.
			continue
		}
		if va.TypeName() == vb.TypeName() {
			out.types[k] = va
			continue
		}
		out.types[k] = jstype.NewUnion(va, vb)
	}
	return out
}

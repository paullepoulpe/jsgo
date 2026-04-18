// Port of com/google/javascript/jscomp/type/FlowScope.java +
// com/google/javascript/jscomp/type/LinkedFlowScope.java
// @ closure-compiler v20260415 (narrow; single-map representation,
// no "outer" chaining, no symbol-source tracking).
//
// A FlowScope maps names to their flow-refined JSType at a given
// program point. It's the lattice element passed around by
// DataFlowAnalysis. Missing names fall through to the declared type
// from the surrounding TypedScope — so a FlowScope only has to carry
// entries for names whose type differs from their declaration.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino/jstype"

// FlowScope is immutable by convention — callers always construct a
// new FlowScope via Inherit / Bind / Join rather than mutating an
// existing one.
type FlowScope struct {
	types map[string]jstype.JSType
}

// EmptyFlowScope is the bottom of the lattice — no names refined.
func EmptyFlowScope() *FlowScope {
	return &FlowScope{types: map[string]jstype.JSType{}}
}

// Get returns the flow-refined type for `name`, or nil if the name
// isn't refined (callers should fall back to the declaring scope).
func (f *FlowScope) Get(name string) jstype.JSType {
	if f == nil {
		return nil
	}
	return f.types[name]
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

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

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// FlowScope is immutable by convention — callers always construct a
// new FlowScope via Inherit / Bind / Join rather than mutating an
// existing one.
//
// `scope`, when non-nil, is the TypedScope the FlowScope hangs off.
// JoinFlowScopes uses it to resolve the declared type of a name that
// was refined on only one incoming edge so loop back-edges don't
// silently drop refinements: the "missing" side contributes the
// declared type instead of being treated as ⊥. Mirrors Java's
// LinkedFlowScope.joinOp.apply, which similarly consults the static
// scope for the other-side type.
type FlowScope struct {
	types map[string]jstype.JSType
	scope *TypedScope
}

// EmptyFlowScope is the bottom of the lattice — no names refined.
func EmptyFlowScope() *FlowScope {
	return &FlowScope{types: map[string]jstype.JSType{}}
}

// NewFlowScope constructs an empty FlowScope anchored to `scope` so
// JoinFlowScopes can resolve declared types for the fallback path.
func NewFlowScope(scope *TypedScope) *FlowScope {
	return &FlowScope{types: map[string]jstype.JSType{}, scope: scope}
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
	out := &FlowScope{types: make(map[string]jstype.JSType, len(f.types)+1), scope: f.scope}
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
	var sc *TypedScope
	base := 0
	if f != nil {
		sc = f.scope
		base = len(f.types)
	}
	out := &FlowScope{types: make(map[string]jstype.JSType, base+len(refs)), scope: sc}
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
// union of the two sides. For names refined on only one branch we
// consult the anchored TypedScope for that name's declared type —
// when it exists and differs from the refined side, the result is a
// widened union. This is what lets loop back-edges propagate
// refinements past the loop header: the pre-header arrives with the
// declared type (no overlay entry), the back-edge arrives with the
// body-refined type, and the join is the union of both. Mirrors
// LinkedFlowScope.joinOp.apply.
func JoinFlowScopes(a, b *FlowScope) *FlowScope {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	var sc *TypedScope
	if a.scope != nil {
		sc = a.scope
	} else {
		sc = b.scope
	}
	out := &FlowScope{types: map[string]jstype.JSType{}, scope: sc}
	// Names present in both sides — union the refinements.
	for k, va := range a.types {
		if vb, ok := b.types[k]; ok {
			if va == nil || vb == nil {
				continue
			}
			if va.TypeName() == vb.TypeName() {
				out.types[k] = va
				continue
			}
			out.types[k] = jstype.NewUnion(va, vb)
			continue
		}
		joinWithDeclared(out, sc, k, va)
	}
	for k, vb := range b.types {
		if _, seen := a.types[k]; seen {
			continue
		}
		joinWithDeclared(out, sc, k, vb)
	}
	return out
}

// joinWithDeclared records the result of joining `refined` with the
// declared type of `name` into `out`. The "other side" of the join
// implicitly carries the declared type because it didn't refine the
// slot; widening the refinement back to (refined ∪ declared) is
// what propagates a body-side narrowing across a loop back-edge.
// Qualified (dotted) names have no declared backing so they're
// dropped — matches Java's property-slot join.
func joinWithDeclared(out *FlowScope, sc *TypedScope, name string, refined jstype.JSType) {
	if sc == nil || isDottedName(name) {
		return
	}
	declared := sc.GetRawDeclaredType(name)
	if declared == nil {
		return
	}
	if declared.TypeName() == refined.TypeName() {
		out.types[name] = refined
		return
	}
	out.types[name] = jstype.NewUnion(refined, declared)
}

// isDottedName reports whether `name` is a qualified-name key like
// "obj.foo" rather than a bare identifier. Qualified-name refinements
// don't have a declared-type fallback in the TypedScope, so JOIN drops
// them when they're present on only one side — mirrors Java's slot-
// based join where property slots lack a static backing.
func isDottedName(name string) bool {
	return strings.IndexByte(name, '.') >= 0
}

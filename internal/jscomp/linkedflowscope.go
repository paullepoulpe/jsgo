// Port of com/google/javascript/jscomp/LinkedFlowScope.java
// @ closure-compiler v20260415 (narrow; flat-map representation with
// lattice parent/child support, JoinOp, inferSlotType,
// inferQualifiedSlot, optimize; no symbol-source tracking).
//
// LinkedFlowScope is the concrete implementation of the FlowScope
// interface (see flowscope.go). In the Java source it uses a linked
// parent chain (PMap / OverlayScope) for efficient copy-on-write; the
// Go port uses a flat map copied eagerly in createChildFlowScope,
// which is semantically equivalent at the cost of O(n) copy instead of
// O(depth) lookup — acceptable given scope sizes in the ported passes.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino/jstype"

// Go-only: language-adapter — EmptyFlowScope is the bottom of the lattice — no names refined.
// Java uses createEntryLattice; Go uses a simpler nil/empty map.
func EmptyFlowScope() *FlowScope {
	return &FlowScope{types: map[string]jstype.JSType{}}
}

// Go-only: language-adapter — createChildFlowScope returns a new FlowScope that inherits all
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

// Go-only: language-adapter — Bind produces a copy of `f` with `name` remapped to `t`. The
// receiver is never mutated — mirrors the persistent-map contract
// LinkedFlowScope documents in Java (via OverlayScope PMap updates).
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

// Go-only: language-adapter — BindAll layers every entry of `refs` onto `f`.
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
// as empty. Mirrors LinkedFlowScope.equals (structural equality).
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

// Go-only: language-adapter — NewJoinOp returns an empty JoinOp ready to accumulate edge scopes.
// Java creates FlowScopeJoinOp via its constructor with an inputProvider.
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

// createEntryLattice creates the initial FlowScope for the entry of a
// function or script. Mirrors LinkedFlowScope.createEntryLattice.
func createEntryLattice(scope *TypedScope) *FlowScope {
	f := EmptyFlowScope()
	f.syntacticScope = scope
	return f
}

// getParentScope returns the parent of the syntactic scope.
// Mirrors LinkedFlowScope.getParentScope.
func (f *FlowScope) getParentScope() *TypedScope {
	if f == nil || f.syntacticScope == nil {
		return nil
	}
	return f.syntacticScope.parent
}

// getSlot resolves a slot by name, first checking flow refinements
// then falling back to the syntactic scope. Mirrors
// LinkedFlowScope.getSlot.
func (f *FlowScope) getSlot(name string) jstype.JSType {
	if t := f.Get(name); t != nil {
		return t
	}
	if f.syntacticScope != nil {
		return f.syntacticScope.GetDeclaredType(name)
	}
	return nil
}

// getOwnSlot checks for a slot exclusively in the flow overlay (not
// in the declared scope). Mirrors LinkedFlowScope.getOwnSlot.
func (f *FlowScope) getOwnSlot(name string) jstype.JSType {
	return f.Get(name)
}

// equals checks structural equality of two FlowScopes.
// Mirrors LinkedFlowScope.equals.
func (f *FlowScope) equals(other *FlowScope) bool {
	return f.Equal(other)
}

// hashCode returns a hash of the FlowScope contents.
// Mirrors LinkedFlowScope.hashCode (approximate).
func (f *FlowScope) hashCode() int {
	if f == nil || len(f.types) == 0 {
		return 0
	}
	h := 0
	for k, v := range f.types {
		h ^= len(k)
		if v != nil {
			h ^= len(v.TypeName())
		}
	}
	return h
}

// getCommonParentDeclarationScope returns the common ancestor scope of
// two FlowScopes. Mirrors
// LinkedFlowScope.getCommonParentDeclarationScope.
func getCommonParentDeclarationScope(left, right *FlowScope) *TypedScope {
	if left == nil {
		if right != nil {
			return right.syntacticScope
		}
		return nil
	}
	if right == nil {
		return left.syntacticScope
	}
	// Walk up left's chain until we find a scope that contains right's
	// scope. This approximates Java's TypedScope.isChildOf walk.
	l := left.syntacticScope
	r := right.syntacticScope
	if l == r {
		return l
	}
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	// Return the shallower of the two scopes as an approximation.
	ld, rd := 0, 0
	for s := l; s != nil; s = s.parent {
		ld++
	}
	for s := r; s != nil; s = s.parent {
		rd++
	}
	if ld <= rd {
		return l
	}
	return r
}

// Go-only: language-adapter — JoinFlowScopes computes the lattice-join of two FlowScopes.
// Java's join logic lives inside LinkedFlowScope.FlowScopeJoinOp.joinFlow
// and the private join() static method — extracted here for clarity.
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
		// Two FunctionTypes with the same arity merge piecewise so
		// the joined slot stays callable. Mirrors FunctionType.try-
		// MergeFunctionPiecewise (FunctionType.java:847) used by
		// JSType.getLeastSupertype.
		if fa, fok := va.(*jstype.FunctionType); fok {
			if fb, fok2 := vb.(*jstype.FunctionType); fok2 && len(fa.Params) == len(fb.Params) {
				out.types[k] = mergeFunctionTypePiecewise(fa, fb)
				continue
			}
		}
		out.types[k] = jstype.NewUnion(va, vb)
	}
	return out
}

// mergeFunctionTypePiecewise builds a FunctionType whose params and
// return are the union of `fa` and `fb`'s corresponding components.
// Port of FunctionType.tryMergeFunctionPiecewise @ closure-compiler v20260415.
// Go-only: language-adapter — exposed at package level so JoinFlowScopes
// can reuse the merge.
func mergeFunctionTypePiecewise(fa, fb *jstype.FunctionType) *jstype.FunctionType {
	mergedParams := make([]jstype.JSType, len(fa.Params))
	for i := range fa.Params {
		switch {
		case fa.Params[i] == nil && fb.Params[i] == nil:
			mergedParams[i] = nil
		case fa.Params[i] == nil:
			mergedParams[i] = fb.Params[i]
		case fb.Params[i] == nil:
			mergedParams[i] = fa.Params[i]
		case fa.Params[i].TypeName() == fb.Params[i].TypeName():
			mergedParams[i] = fa.Params[i]
		default:
			mergedParams[i] = jstype.NewUnion(fa.Params[i], fb.Params[i])
		}
	}
	var ret jstype.JSType
	switch {
	case fa.Return == nil && fb.Return == nil:
		ret = nil
	case fa.Return == nil:
		ret = fb.Return
	case fb.Return == nil:
		ret = fa.Return
	case fa.Return.TypeName() == fb.Return.TypeName():
		ret = fa.Return
	default:
		ret = jstype.NewUnion(fa.Return, fb.Return)
	}
	return jstype.NewFunctionType(mergedParams, ret)
}

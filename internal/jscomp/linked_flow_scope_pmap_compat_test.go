package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type flowScopePMapCalls struct {
	get, plus, minus, reconcile, equivalent int
}

// This persistent-map implementation delegates storage to HAMT but keeps its
// own implementation and identity through updates. Binary operations accept
// only the same implementation, as required by PMap.
type flowScopeWrappedPMap[K, V any] struct {
	*rhino.HamtPMap[K, V]
	calls *flowScopePMapCalls
}

func (m *flowScopeWrappedPMap[K, V]) Get(key K) V {
	m.calls.get++
	return m.HamtPMap.Get(key)
}

func (m *flowScopeWrappedPMap[K, V]) updated(next *rhino.HamtPMap[K, V]) rhino.PMap[K, V] {
	if next == m.HamtPMap {
		return m
	}
	return &flowScopeWrappedPMap[K, V]{HamtPMap: next, calls: m.calls}
}

func (m *flowScopeWrappedPMap[K, V]) Plus(key K, value V) rhino.PMap[K, V] {
	m.calls.plus++
	return m.updated(m.HamtPMap.Plus(key, value).(*rhino.HamtPMap[K, V]))
}

func (m *flowScopeWrappedPMap[K, V]) Minus(key K) rhino.PMap[K, V] {
	m.calls.minus++
	return m.updated(m.HamtPMap.Minus(key).(*rhino.HamtPMap[K, V]))
}

func (m *flowScopeWrappedPMap[K, V]) Reconcile(other rhino.PMap[K, V], joiner rhino.PMap_Reconciler[K, V]) rhino.PMap[K, V] {
	m.calls.reconcile++
	return m.updated(m.HamtPMap.Reconcile(other.(*flowScopeWrappedPMap[K, V]).HamtPMap, joiner).(*rhino.HamtPMap[K, V]))
}

func (m *flowScopeWrappedPMap[K, V]) Equivalent(other rhino.PMap[K, V], equivalent func(V, V) bool) bool {
	m.calls.equivalent++
	return m.HamtPMap.Equivalent(other.(*flowScopeWrappedPMap[K, V]).HamtPMap, equivalent)
}

func flowScopeCompatFixture() (*Compiler, *TypedScope, *TypedScope, rhino.JSType, rhino.JSType) {
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	fn := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(), rhino.IR_Block())
	global := TypedScope_CreateGlobalScope(rhino.IR_Root(rhino.IR_ScriptOfNodes(fn)))
	local := NewTypedScope(global, fn)
	local.Declare("x", nil, nil, nil, true)
	registry := compiler.GetTypeRegistry()
	return compiler, global, local, registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE), registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
}

func TestLinkedFlowScopeCustomPMapConstructorsAndUpdates(t *testing.T) {
	compiler, global, local, number, text := flowScopeCompatFixture()
	slotCalls, scopeCalls := &flowScopePMapCalls{}, &flowScopePMapCalls{}
	emptySlots := &flowScopeWrappedPMap[string, *LinkedFlowScope_OverlaySlot]{
		HamtPMap: rhino.HamtPMap_Empty[string, *LinkedFlowScope_OverlaySlot](), calls: slotCalls,
	}
	overlay := NewLinkedFlowScope_OverlayScopeWithSlots(local, emptySlots)
	numberOverlay := overlay.Infer("x", number)
	if !emptySlots.IsEmpty() || numberOverlay.GetSlot("x").GetType() != number || numberOverlay.Infer("x", number) != numberOverlay {
		t.Fatal("custom slot-map inference changed persistence, value, or no-op identity")
	}
	if _, ok := any(numberOverlay.slots).(*flowScopeWrappedPMap[string, *LinkedFlowScope_OverlaySlot]); !ok {
		t.Fatal("custom slot-map implementation was replaced")
	}
	emptyScopes := &flowScopeWrappedPMap[*TypedScope, *LinkedFlowScope_OverlayScope]{
		HamtPMap: rhino.HamtPMap_Empty[*TypedScope, *LinkedFlowScope_OverlayScope](), calls: scopeCalls,
	}
	maps := emptyScopes.Plus(local, numberOverlay)
	flow := NewLinkedFlowScope(compiler, maps, local, local)
	if flow.TrimScopes(local) != maps || flow.GetSlot("x").(rhino.StaticTypedSlot).GetType() != number {
		t.Fatal("custom scope-map construction or unchanged trimming lost map identity/value")
	}
	updated := flow.InferSlotType("x", text).(*LinkedFlowScope)
	if updated.GetSlot("x").(rhino.StaticTypedSlot).GetType() != text || flow.GetSlot("x").(rhino.StaticTypedSlot).GetType() != number {
		t.Fatal("scope-map update mutated its persistent predecessor")
	}
	if _, ok := any(updated.scopes).(*flowScopeWrappedPMap[*TypedScope, *LinkedFlowScope_OverlayScope]); !ok {
		t.Fatal("custom scope-map implementation was replaced")
	}
	trimmed := updated.WithSyntacticScope(global).(*LinkedFlowScope)
	if !trimmed.scopes.IsEmpty() || updated.scopes.IsEmpty() || scopeCalls.minus == 0 || scopeCalls.get == 0 || slotCalls.get == 0 || slotCalls.plus == 0 {
		t.Fatal("custom map operations or persistent child-scope trimming were bypassed")
	}
}

func TestLinkedFlowScopeNativeAndCustomPMapJoin(t *testing.T) {
	for _, custom := range []bool{false, true} {
		label := "native"
		if custom {
			label = "custom"
		}
		t.Run(label, func(t *testing.T) {
			compiler, _, local, number, text := flowScopeCompatFixture()
			slotCalls, scopeCalls := &flowScopePMapCalls{}, &flowScopePMapCalls{}
			var slots rhino.PMap[string, *LinkedFlowScope_OverlaySlot] = rhino.HamtPMap_Empty[string, *LinkedFlowScope_OverlaySlot]()
			var scopes rhino.PMap[*TypedScope, *LinkedFlowScope_OverlayScope] = rhino.HamtPMap_Empty[*TypedScope, *LinkedFlowScope_OverlayScope]()
			if custom {
				slots = &flowScopeWrappedPMap[string, *LinkedFlowScope_OverlaySlot]{HamtPMap: slots.(*rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]), calls: slotCalls}
				scopes = &flowScopeWrappedPMap[*TypedScope, *LinkedFlowScope_OverlayScope]{HamtPMap: scopes.(*rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]), calls: scopeCalls}
			}
			leftOverlay := NewLinkedFlowScope_OverlayScopeWithSlots(local, slots).Infer("x", number)
			rightOverlay := NewLinkedFlowScope_OverlayScopeWithSlots(local, slots).Infer("x", text)
			left := NewLinkedFlowScope(compiler, scopes.Plus(local, leftOverlay), local, local)
			right := NewLinkedFlowScope(compiler, scopes.Plus(local, rightOverlay), local, local)
			copyLeft := NewLinkedFlowScope(compiler, left.scopes, local, local)
			if !left.Equals(copyLeft) || left.Equals(right) {
				t.Fatal("map equivalence failed for equal and different slot types")
			}
			joined := LinkedFlowScope_Join(left, right, local)
			joinedFlow := NewLinkedFlowScope(compiler, joined, local, local)
			got := joinedFlow.GetSlot("x").(rhino.StaticTypedSlot).GetType()
			if got.DiffersFrom(number.GetLeastSupertype(text)) {
				t.Fatal("map reconciliation did not join conflicting inferred slot types")
			}
			joiner := NewLinkedFlowScope_FlowScopeJoinOp(compiler)
			joiner.JoinFlow(left)
			joiner.JoinFlow(right)
			if !joinedFlow.Equals(joiner.Finish()) {
				t.Fatal("public join and flow join operator disagree")
			}
			if custom {
				if _, ok := joined.(*flowScopeWrappedPMap[*TypedScope, *LinkedFlowScope_OverlayScope]); !ok {
					t.Fatal("join replaced the custom map implementation")
				}
				if scopeCalls.reconcile == 0 || scopeCalls.equivalent == 0 || slotCalls.reconcile == 0 || slotCalls.equivalent == 0 {
					t.Fatal("custom binary map operations were bypassed")
				}
			}
		})
	}
}

func TestLinkedFlowScopeCustomEmptySlots(t *testing.T) {
	_, _, local, number, _ := flowScopeCompatFixture()
	calls := &flowScopePMapCalls{}
	empty := &flowScopeWrappedPMap[string, *LinkedFlowScope_OverlaySlot]{
		HamtPMap: rhino.HamtPMap_Empty[string, *LinkedFlowScope_OverlaySlot](), calls: calls,
	}
	switch target := any(&LinkedFlowScope_EMPTY_SLOTS).(type) {
	case *rhino.PMap[string, *LinkedFlowScope_OverlaySlot]:
		original := *target
		t.Cleanup(func() { *target = original })
		*target = empty
	default:
		t.Fatal("exported empty-slot map no longer accepts the PMap interface")
	}
	overlay := NewLinkedFlowScope_OverlayScope(local)
	if any(overlay.slots) != any(empty) || overlay.Infer("x", number).GetSlot("x").GetType() != number || calls.plus != 1 {
		t.Fatal("overlay constructor ignored the configured persistent empty-slot map")
	}
}

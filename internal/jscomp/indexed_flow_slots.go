// Copyright 2026 The jsgo Authors.
// Licensed under the Apache License, Version 2.0.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/flowstore"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"os"
)

// This opt-in changes only inferred-slot storage. The surrounding type inference,
// declaration, trimming and type-equality rules remain in LinkedFlowScope.
var indexedFlowSlotsEnabled = os.Getenv("JSGO_EXPERIMENT_INDEXED_FLOW_SLOTS") == "1"

// Initialization stays on the existing typed-scope-building goroutine. The owner
// synchronizes later symbol interning and reads; immutable versions need no lock.
func emptyFlowSlots(scope *TypedScope) rhino.PMap[string, *LinkedFlowScope_OverlaySlot] {
	if !indexedFlowSlotsEnabled {
		return LinkedFlowScope_EMPTY_SLOTS
	}
	if scope.indexedFlowSlots == nil {
		scope.indexedFlowSlots = flowstore.NewIDs()
	}
	return flowstore.New[*LinkedFlowScope_OverlaySlot](scope.indexedFlowSlots)
}

func normalizeFlowSlots(scope *TypedScope, slots rhino.PMap[string, *LinkedFlowScope_OverlaySlot]) rhino.PMap[string, *LinkedFlowScope_OverlaySlot] {
	if !indexedFlowSlotsEnabled {
		return slots
	}
	if _, ok := slots.(*flowstore.Store[*LinkedFlowScope_OverlaySlot]); ok {
		return slots
	}
	result := emptyFlowSlots(scope)
	for key := range slots.Keys() {
		result = result.Plus(key, slots.Get(key))
	}
	return result
}

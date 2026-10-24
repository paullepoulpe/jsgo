package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
	"unsafe"
)

func TestCompactVarCrossClassEquality(t *testing.T) {
	// 48: the name node is a 16-byte handle (rhino/node_store.go).
	if unsafe.Sizeof(Var{}) != 48 {
		t.Fatalf("Var has %d bytes, want 48", unsafe.Sizeof(Var{}))
	}
	if unsafe.Sizeof(TypedVar{}) != 72 {
		t.Fatalf("TypedVar has %d bytes, want 72", unsafe.Sizeof(TypedVar{}))
	}
	root := rhino.IR_Root()
	scope := Scope_CreateGlobalScope(root)
	a, b := NewVar("x", rhino.NilNode, scope, 0, nil, rhino.NilNode), NewVar("x", rhino.NilNode, scope, 1, nil, rhino.NilNode)
	named := ScopedName_Of("x", root)
	for _, other := range []ScopedName{a, b, named} {
		if !a.Equals(other) || !other.Equals(a) || a.HashCode() != other.HashCode() {
			t.Fatal("equality/hash diverged from ScopedName")
		}
	}
	// Hash zero is a valid value; a lazy-cache sentinel must not alter it.
	zero := NewVar("", rhino.NilNode, scope, -1, nil, rhino.NilNode)
	for i := 0; i < 5; i++ {
		if zero.HashCode() != ScopedName_Of("", root).HashCode() {
			t.Fatal("zero name hash changed value")
		}
	}
}

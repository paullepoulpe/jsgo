package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
	"unsafe"
)

func TestCompactVarCrossClassEquality(t *testing.T) {
	if unsafe.Sizeof(Var{}) != 64 {
		t.Fatalf("Var has %d bytes, want 64", unsafe.Sizeof(Var{}))
	}
	root := rhino.IR_Root()
	scope := Scope_CreateGlobalScope(root)
	a, b := NewVar("x", nil, scope, 0, nil, nil), NewVar("x", nil, scope, 1, nil, nil)
	named := ScopedName_Of("x", root)
	for _, other := range []ScopedName{a, b, named} {
		if !a.Equals(other) || !other.Equals(a) || a.HashCode() != other.HashCode() {
			t.Fatal("equality/hash diverged from ScopedName")
		}
	}
	// Hash zero is a valid value; a lazy-cache sentinel must not alter it.
	zero := NewVar("", nil, scope, -1, nil, nil)
	for i := 0; i < 5; i++ {
		if zero.HashCode() != ScopedName_Of("", root).HashCode() {
			t.Fatal("zero name hash changed value")
		}
	}
}

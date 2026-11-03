package jscomp

import (
	"reflect"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestCompactVarCrossClassEquality(t *testing.T) {
	if sz := reflect.TypeFor[Var]().Size(); sz != 32 {
		t.Fatalf("Var has %d bytes, want 32", sz)
	}
	if sz := reflect.TypeFor[TypedVar]().Size(); sz != 48 {
		t.Fatalf("TypedVar has %d bytes, want 48", sz)
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

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestVarHashMatchesScopedName(t *testing.T) {
	for _, name := range []string{"", "x", "long_identifier_reused_by_passes", "λ", jutil.StringFromUTF16([]uint16{0xD800, 'x'})} {
		for _, index := range []int{-1, 0, 1, int(varIndexMask) - 1} {
			scope := Scope_CreateGlobalScope(rhino.IR_Root())
			v := NewVar(name, rhino.NilNode, scope, index, nil, rhino.NilNode)
			want := ScopedName_Of(name, scope.GetRootNode()).HashCode()
			for repeat := 0; repeat < 3; repeat++ {
				if v.HashCode() != want {
					t.Fatalf("hash for %q: %d, want %d", name, v.HashCode(), want)
				}
			}
			if v.GetIndex() != index {
				t.Fatalf("index %d became %d", index, v.GetIndex())
			}
		}
	}
}

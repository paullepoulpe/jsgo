package jscomp

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestScopeDeclarationKeysSurvivePublicVarReinitialization(t *testing.T) {
	for _, tc := range []struct {
		name  string
		count int
	}{{"small", 2}, {"hashed", 12}} {
		t.Run(tc.name, func(t *testing.T) {
			count := tc.count
			scope := Scope_CreateGlobalScope(rhino.IR_Root())
			for i := 0; i < count; i++ {
				name := string(rune('a' + i))
				scope.Declare(name, rhino.IR_Name(name), nil)
			}
			view := scope.GetAllSymbols()
			original := scope.GetVar("a")
			original.InitAbstractVarBase("renamed", original.GetNode(), scope, original.GetIndex(), original.GetInput(), nil)
			if scope.GetVar("a") != original || scope.GetVar("renamed") != nil || !scope.HasOwnSlot("a") || !scope.HasSlot("a") || scope.HasSlot("renamed") || scope.GetOwnSlot("a") != original {
				t.Fatal("public reinitialization changed the declaration key")
			}
			alias := NewVar("different", nil, scope, count, nil, nil)
			scope.DeclareInternal("alias", alias)
			if scope.GetVar("alias") != alias || scope.GetVar("different") != nil || scope.GetVarCount() != count+1 {
				t.Fatal("independent declaration key was lost")
			}
			got := slices.Collect(view)
			if len(got) != count+1 || got[0] != original || got[count] != alias {
				t.Fatal("delayed symbol view lost ordered identities after promotion")
			}
			scope.ClearVarsInternal()
			replacement := scope.Declare("replacement", rhino.IR_Name("replacement"), nil)
			if got := slices.Collect(view); len(got) != 1 || got[0] != replacement || scope.GetVarCount() != 1 {
				t.Fatal("promoted symbol view stopped following clear and redeclaration")
			}
		})
	}
	t.Run("independent_key", func(t *testing.T) {
		scope := Scope_CreateGlobalScope(rhino.IR_Root())
		first := scope.Declare("first", rhino.IR_Name("first"), nil)
		view := scope.GetVarIterable()
		value := NewVar("different", nil, scope, 1, nil, nil)
		scope.DeclareInternal("key", value)
		if scope.GetVar("first") != first || scope.GetVar("key") != value || scope.GetVar("different") != nil || !scope.HasSlot("key") || scope.GetVarCount() != 2 {
			t.Fatal("independently supplied key was replaced by the variable name")
		}
		if got := slices.Collect(view); len(got) != 2 || got[0] != first || got[1] != value {
			t.Fatal("independent key promotion lost the existing ordered view")
		}
	})

}

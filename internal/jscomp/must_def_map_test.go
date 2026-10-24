package jscomp

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Differential state operations check TOP/BOTTOM, value-equal Vars, copy
// isolation, and insertion order against the original Java-semantics map.
func TestMustDefDenseStateDifferential(t *testing.T) {
	for _, count := range []int{1, 8, 32, 64, 65} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			scope := Scope_CreateGlobalScope(rhino.IR_Root())
			vars := jutil.NewLinkedHashMap[string, *Var]()
			for i := 0; i < count; i++ {
				name := fmt.Sprintf("v%d", i)
				vars.Put(name, NewVar(name, rhino.Node{}, scope, i, nil, rhino.Node{}))
			}
			list := vars.KeySlice()
			got := newMustDefMap(newMustDefLayout(vars))
			want := jutil.NewObjectLinkedHashMap[*Var, *MustBeReachingVariableDef_Definition]()
			rng := rand.New(rand.NewSource(int64(count)))
			for step := 0; step < 500; step++ {
				v := vars.Get(list[rng.Intn(count)])
				var d *MustBeReachingVariableDef_Definition
				if rng.Intn(3) != 0 {
					d = NewMustBeReachingVariableDef_Definition(rhino.IR_Name(fmt.Sprint(step)))
				}
				got.Put(v, d)
				want.Put(v, d)
				var order []*Var
				for k, value := range got.All() {
					order = append(order, k)
					w, ok := want.Lookup(k)
					if !ok || value != w {
						t.Fatalf("step %d: value mismatch", step)
					}
				}
				if len(order) != want.Size() {
					t.Fatalf("step %d: size mismatch", step)
				}
				i := 0
				for k := range want.Keys() {
					if order[i] != k {
						t.Fatalf("step %d: insertion order mismatch", step)
					}
					i++
				}
				if got.HashCode() != want.HashCode() {
					t.Fatalf("step %d: hash mismatch", step)
				}
				copy := got.clone()
				if !copy.equals(&got) {
					t.Fatal("copy differs")
				}
				old := got.Get(v)
				copy.Put(v, NewMustBeReachingVariableDef_Definition(rhino.IR_Name("copy")))
				if got.Get(v) != old {
					t.Fatal("copy modified original")
				}
				alias := NewVar(v.GetName(), rhino.Node{}, scope, -1, nil, rhino.Node{})
				if got.Get(alias) != d || !got.ContainsKey(alias) {
					t.Fatal("ScopedName equality differs")
				}
			}
			// A value-equal non-canonical key can enter through public lattice
			// constructors; falling back must preserve its original key object.
			v := vars.Get(list[0])
			alias := NewVar(v.GetName(), rhino.Node{}, scope, -1, nil, rhino.Node{})
			got.Put(alias, nil)
			want.Put(alias, nil)
			if got.Get(v) != nil || got.HashCode() != want.HashCode() {
				t.Fatal("fallback differs")
			}
		})
	}
}

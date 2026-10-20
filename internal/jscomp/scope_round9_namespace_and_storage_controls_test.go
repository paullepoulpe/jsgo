package jscomp

import (
	"reflect"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type round9NamespaceRetainingHandler struct {
	scope    *Scope
	variable *Var
	calls    int
}

func (h *round9NamespaceRetainingHandler) OnRedeclaration(scope *Scope, name string, _ *rhino.Node, _ *CompilerInput) {
	h.calls++
	h.scope = scope
	h.variable = scope.GetVar(name)
}

// Existing public creator option invokes its callback from the goog.provide
// path, separately from ordinary duplicate-name declarations. No new pooling
// API or private storage field is used by this regression.
func TestScopeProvideCallbackPreservesRetainedScopeAfterRelease(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("var ns; goog.provide('ns.child');")
	if root.GetParent() != nil {
		root.Detach()
	}
	handler := &round9NamespaceRetainingHandler{}
	creator := NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations(fixture.compiler, handler, true)
	scope := creator.CreateScopeOfScope(root, nil)
	variable := scope.GetVar("ns")
	if handler.calls != 1 || handler.scope != scope || handler.variable != variable || variable == nil {
		t.Fatal("fixture did not invoke the existing goog.provide redeclaration callback")
	}
	creator.ReleaseArena()
	if handler.scope.GetRootNode() != root || handler.scope.GetVar("ns") != variable || handler.variable.GetName() != "ns" {
		t.Fatal("release cleared Scope or Var retained by goog.provide redeclaration callback")
	}
}

func TestScopeStorageKeepsLookupOrderAcrossCollisionsAndMutation(t *testing.T) {
	scope := Scope_CreateGlobalScope(rhino.IR_Root())
	names := []string{"collision13", "collision63", "collision323", "collision416", "collision432", "collision709", "collision870", "collision1013", "collision1129", "collision1199", "collision1200", "collision1405", "collision1451", "collision1580", "collision2032", "collision2236", "collision2284", "collision2418", "collision2422", "collision2446", "collision2521", "collision2599", "collision2787", "collision2849", "collision2920", "collision3025", "collision3055", "collision3104", "collision3468", "collision3541", "collision3727", "collision3839", "collision3889", "collision4093", "collision4275", "collision4421"}
	vars := make(map[string]*Var, len(names))
	for _, name := range names {
		vars[name] = scope.Declare(name, rhino.IR_Name(name), nil)
	}
	view := scope.GetVarIterable()
	deleted := map[string]bool{names[0]: true, names[17]: true, names[len(names)-1]: true}
	for _, name := range []string{names[0], names[17], names[len(names)-1]} {
		scope.Undeclare(vars[name])
	}
	var want []string
	for _, name := range names {
		if deleted[name] {
			if scope.GetVar(name) != nil || scope.HasOwnSlot(name) {
				t.Fatal("removed collision key remained visible")
			}
			continue
		}
		want = append(want, name)
		if scope.GetVar(name) != vars[name] || scope.GetOwnSlot(name) != vars[name] || !scope.HasSlot(name) {
			t.Fatal("collision removal lost a surviving binding")
		}
	}
	for _, name := range []string{"afterA", "afterB", "afterC"} {
		scope.Declare(name, rhino.IR_Name(name), nil)
		want = append(want, name)
	}
	var got []string
	for v := range view {
		got = append(got, v.GetName())
	}
	if !reflect.DeepEqual(got, want) || scope.GetVarCount() != len(want) {
		t.Fatalf("live insertion-order view after remove/grow: got%v want%v", got, want)
	}
	all := scope.GetAllSymbols()
	scope.ClearVarsInternal()
	replacement := scope.Declare("replacement", rhino.IR_Name("replacement"), nil)
	var final []*Var
	for v := range all {
		final = append(final, v)
	}
	if len(final) != 1 || final[0] != replacement || scope.GetVarCount() != 1 || scope.GetVar(names[1]) != nil {
		t.Fatal("clear/redeclare changed existing live-view behavior or retained a stale hash entry")
	}
}

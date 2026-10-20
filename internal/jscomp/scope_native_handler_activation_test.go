package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// This activation control is candidate-only: these storage fields were added
// upstream. Public callback and declaration-key controls live in separate files.
func TestSyntacticScopeCreator_CustomHandlerWithoutRedeclarationsKeepsNativeStorage(t *testing.T) {
	for _, tc := range []struct {
		name  string
		code  string
		count int
	}{{"small", "var a,b;", 2}, {"hashed", "var a,b,c,d,e,f,g,h,i,j,k,l;", 12}} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newSyntacticScopeCreatorTest(t)
			handler := &scopeCallbackVisibilityHandler{callback: func(*Scope, string, *rhino.Node, *CompilerInput) {
				t.Fatal("unexpected redeclaration")
			}}
			scope := NewSyntacticScopeCreatorWithRedeclarationHandler(fixture.compiler, handler).CreateScopeOfScope(fixture.getRoot(tc.code), nil)
			if scope.keyedVars || scope.vars.HasBacking() || scope.GetVarCount() != tc.count || len(scope.varList) != tc.count {
				t.Fatal("a custom handler allocated keyed state without a callback")
			}
			if (scope.varHash != nil) != (tc.count > 8) {
				t.Fatal("custom handler bypassed native small/hash storage")
			}
		})
	}
}

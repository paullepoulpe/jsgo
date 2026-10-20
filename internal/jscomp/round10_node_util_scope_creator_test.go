package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// A public ScopeCreator may provide bindings independently of declaration
// syntax. Even an empty block must reach its creator when collecting symbols.
type round10DeclaredScopeCreator struct {
	target   *rhino.Node
	calls    []*rhino.Node
	provided *Var
}

func (c *round10DeclaredScopeCreator) CreateScope(root *rhino.Node, parent AbstractScope) AbstractScope {
	c.calls = append(c.calls, root)
	scope := Scope_CreateChildScope(parent.(*Scope), root)
	if root == c.target {
		c.provided = scope.Declare("provided", rhino.IR_Name("provided"), nil)
	}
	return scope
}

func TestNodeUtilAllVarsHonorsCustomCreatorForDeclarationFreeScopes(t *testing.T) {
	for _, nested := range []bool{false, true} {
		name := "empty_function_body"
		code := "function f() {}"
		if nested {
			name = "empty_nested_block"
			code = "function f() { {} }"
		}
		t.Run(name, func(t *testing.T) {
			fixture := newSyntacticScopeCreatorTest(t)
			script := fixture.getRoot(code)
			fn := script.GetFirstChild()
			global := fixture.scopeCreator.CreateScopeOfScope(script, nil)
			function := fixture.scopeCreator.CreateScopeOfScope(fn, global)
			body := fn.GetLastChild()
			target := body
			wantCalls := []*rhino.Node{body}
			if nested {
				target = body.GetFirstChild()
				wantCalls = append(wantCalls, target)
			}
			creator := &round10DeclaredScopeCreator{target: target}
			vars := NodeUtil_GetAllVarsDeclaredInFunction(fixture.compiler, creator, function)
			if len(creator.calls) != len(wantCalls) {
				t.Fatalf("public collector skipped custom scope creation: got %d calls, want %d", len(creator.calls), len(wantCalls))
			}
			for i, root := range wantCalls {
				if creator.calls[i] != root {
					t.Fatalf("custom scope creation changed order at call %d", i)
				}
			}
			if creator.provided == nil || vars.GetAllVariables().Get("provided") != creator.provided {
				t.Fatal("public collector dropped a binding supplied for a declaration-free scope")
			}
			found := 0
			for _, v := range vars.GetAllVariablesInOrder() {
				if v == creator.provided {
					found++
				}
			}
			if found != 1 {
				t.Fatalf("provided variable occurs %d times in declaration order, want 1", found)
			}
		})
	}
}

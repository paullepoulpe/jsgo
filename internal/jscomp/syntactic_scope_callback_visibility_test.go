package jscomp

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type scopeCallbackVisibilityHandler struct {
	callback func(*Scope, string, *rhino.Node, *CompilerInput)
}

func (h *scopeCallbackVisibilityHandler) OnRedeclaration(s *Scope, name string, n *rhino.Node, input *CompilerInput) {
	h.callback(s, name, n, input)
}

// A slice-bearing value is a valid handler, but comparing two interfaces
// containing it panics. It must follow the same callback path as a pointer.
type scopeCallbackVisibilityValueHandler struct {
	marker   []byte
	callback func(*Scope, string, *rhino.Node, *CompilerInput)
}

func (h scopeCallbackVisibilityValueHandler) OnRedeclaration(s *Scope, name string, n *rhino.Node, input *CompilerInput) {
	h.callback(s, name, n, input)
}

// Embedding the stock implementation does not make an overriding wrapper a
// stock no-op callback.
type scopeCallbackVisibilityWrapper struct {
	SyntacticScopeCreator_DefaultRedeclarationHandler
	callback func(*Scope, string, *rhino.Node, *CompilerInput)
}

func (h *scopeCallbackVisibilityWrapper) OnRedeclaration(s *Scope, name string, n *rhino.Node, input *CompilerInput) {
	h.callback(s, name, n, input)
}

func scopeCallbackVisibilitySetDefault(t *testing.T, handler SyntacticScopeCreator_RedeclarationHandler) {
	t.Helper()
	previous := SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER
	t.Cleanup(func() { SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER = previous })
	SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER = handler
}

func scopeCallbackVisibilityCheckKeys(t *testing.T, scope *Scope, want []string) {
	t.Helper()
	var got []string
	for key := range scope.vars.Keys() {
		got = append(got, key)
	}
	if !slices.Equal(got, want) || scope.GetVarCount() != len(want) {
		t.Fatalf("scope keys/count = %v/%d, want %v/%d", got, scope.GetVarCount(), want, len(want))
	}
}

func TestSyntacticScopeCreator_DefaultHandlerCallbackSeesStagedBindings(t *testing.T) {
	types := []struct {
		name string
		make func(func(*Scope, string, *rhino.Node, *CompilerInput)) SyntacticScopeCreator_RedeclarationHandler
	}{
		{"pointer", func(f func(*Scope, string, *rhino.Node, *CompilerInput)) SyntacticScopeCreator_RedeclarationHandler {
			return &scopeCallbackVisibilityHandler{callback: f}
		}},
		{"noncomparable_value", func(f func(*Scope, string, *rhino.Node, *CompilerInput)) SyntacticScopeCreator_RedeclarationHandler {
			return scopeCallbackVisibilityValueHandler{marker: []byte{1}, callback: f}
		}},
		{"overriding_wrapper", func(f func(*Scope, string, *rhino.Node, *CompilerInput)) SyntacticScopeCreator_RedeclarationHandler {
			return &scopeCallbackVisibilityWrapper{callback: f}
		}},
	}
	for _, typ := range types {
		t.Run(typ.name, func(t *testing.T) {
			fixture := newSyntacticScopeCreatorTest(t)
			root := fixture.getRoot("var a,b,c,d,e,f,g,h; var a; var tail; var b;")
			initial := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
			var originalA, originalC, replacementB *Var
			calls := 0
			handler := typ.make(func(scope *Scope, name string, _ *rhino.Node, input *CompilerInput) {
				calls++
				switch calls {
				case 1:
					if name != "a" {
						t.Fatalf("first redeclaration = %q, want a", name)
					}
					scopeCallbackVisibilityCheckKeys(t, scope, initial)
					var iterated []*Var
					for variable := range scope.GetVarIterable() {
						iterated = append(iterated, variable)
					}
					for i, key := range initial {
						variable := scope.GetVar(key)
						if variable == nil || variable.GetIndex() != i || iterated[i] != variable {
							t.Fatalf("binding %s index/iteration identity differs at %d", key, i)
						}
					}
					originalA = scope.GetVar("a")
					originalC = scope.GetVar("c")
					originalA.InitAbstractVarBase("renamedA", originalA.GetNode(), scope, originalA.GetIndex(), originalA.GetInput(), nil)
					scope.DeclareInternal("alias", originalA)
					replacementB = NewVar("replacementB", nil, scope, 1, input, nil)
					scope.DeclareInternal("b", replacementB)
				case 2:
					if name != "b" {
						t.Fatalf("second redeclaration = %q, want b", name)
					}
					scopeCallbackVisibilityCheckKeys(t, scope, append(slices.Clone(initial), "alias", "tail"))
					if scope.GetVar("a") != originalA || scope.GetVar("alias") != originalA || scope.GetVar("renamedA") != nil {
						t.Fatal("callback mutation changed original declaration keys")
					}
					if scope.GetVar("b") != replacementB || scope.GetVar("replacementB") != nil || scope.GetVar("tail").GetIndex() != 9 {
						t.Fatal("replacement or subsequent declaration index was not preserved")
					}
					scope.UndeclareInteral(originalC)
					scope.DeclareInternal("after", NewVar("differentAfter", nil, scope, scope.GetVarCount(), input, nil))
				default:
					t.Fatalf("unexpected redeclaration %q", name)
				}
			})
			scopeCallbackVisibilitySetDefault(t, handler)
			scope := NewSyntacticScopeCreator(fixture.compiler).CreateScopeOfScope(root, nil)
			if calls != 2 {
				t.Fatalf("redeclaration callbacks = %d, want 2", calls)
			}
			scopeCallbackVisibilityCheckKeys(t, scope, []string{"a", "b", "d", "e", "f", "g", "h", "alias", "tail", "after"})
			if scope.GetVar("c") != nil || scope.GetVar("a") != originalA || scope.GetVar("b") != replacementB {
				t.Fatal("scanner republished bindings changed by the callback")
			}
		})
	}
}

func TestSyntacticScopeCreator_DefaultHandlerCallbackCanClearBindings(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("var first; var second; var first; var later;")
	calls := 0
	scopeCallbackVisibilitySetDefault(t, &scopeCallbackVisibilityHandler{callback: func(scope *Scope, name string, _ *rhino.Node, input *CompilerInput) {
		calls++
		scopeCallbackVisibilityCheckKeys(t, scope, []string{"first", "second"})
		scope.ClearVarsInternal()
		scope.DeclareInternal("alias", NewVar("different", nil, scope, 0, input, nil))
	}})
	scope := NewSyntacticScopeCreator(fixture.compiler).CreateScopeOfScope(root, nil)
	if calls != 1 {
		t.Fatalf("redeclaration callbacks = %d, want 1", calls)
	}
	scopeCallbackVisibilityCheckKeys(t, scope, []string{"alias", "later"})
	if scope.GetVar("first") != nil || scope.GetVar("second") != nil || scope.GetVar("later").GetIndex() != 1 {
		t.Fatal("cleared staged bindings reappeared or later index changed")
	}
}

func TestSyntacticScopeCreator_DefaultHandlerCallbackPanicIsPreserved(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("var a; var b; var a;")
	sentinel := &struct{ message string }{"callback panic"}
	scopeCallbackVisibilitySetDefault(t, &scopeCallbackVisibilityHandler{callback: func(scope *Scope, _ string, _ *rhino.Node, _ *CompilerInput) {
		scopeCallbackVisibilityCheckKeys(t, scope, []string{"a", "b"})
		panic(sentinel)
	}})
	defer func() {
		if got := recover(); got != sentinel {
			t.Fatalf("panic = %v, want original callback sentinel", got)
		}
	}()
	NewSyntacticScopeCreator(fixture.compiler).CreateScopeOfScope(root, nil)
	t.Fatal("callback panic was not propagated")
}

func TestSyntacticScopeCreator_StockDefaultHandlerVariants(t *testing.T) {
	for _, typedNil := range []bool{false, true} {
		name := "nonnull"
		if typedNil {
			name = "typed_nil"
		}
		t.Run(name, func(t *testing.T) {
			fixture := newSyntacticScopeCreatorTest(t)
			root := fixture.getRoot("var a; var b; var a; var c;")
			var handler *SyntacticScopeCreator_DefaultRedeclarationHandler
			if !typedNil {
				handler = NewSyntacticScopeCreator_DefaultRedeclarationHandler()
			}
			scopeCallbackVisibilitySetDefault(t, handler)
			scope := NewSyntacticScopeCreator(fixture.compiler).CreateScopeOfScope(root, nil)
			scopeCallbackVisibilityCheckKeys(t, scope, []string{"a", "b", "c"})
			if scope.GetVar("c").GetIndex() != 2 {
				t.Fatal("stock handler changed declaration index")
			}
		})
	}
}

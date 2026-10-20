package jscomp

import (
	"fmt"
	"iter"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Uses preexisting APIs so the same roots can run on the historical baseline.
type auditMemoizedScopeDelegate struct {
	calls     int
	makeScope func(*rhino.Node) AbstractScope
}

func (d *auditMemoizedScopeDelegate) CreateScope(n *rhino.Node, parent AbstractScope) AbstractScope {
	d.calls++
	if d.makeScope != nil {
		return d.makeScope(n)
	}
	return Scope_CreateGlobalScope(n)
}

func TestMemoizedScopeCreatorValueCopiesShareNewEntries(t *testing.T) {
	for _, seed := range []int{0, 1, 8, 9} {
		t.Run(fmt.Sprint(seed), func(t *testing.T) {
			delegate := &auditMemoizedScopeDelegate{}
			original := NewMemoizedScopeCreator(delegate)
			for range seed {
				original.CreateScope(rhino.IR_Root(), nil)
			}
			copied := *original
			n := rhino.IR_Root()
			want := original.CreateScope(n, nil)
			if got := copied.CreateScope(n, nil); got != want || delegate.calls != seed+1 {
				t.Fatalf("copied creator returned %p rather than shared %p; delegate calls=%d, want %d", got, want, delegate.calls, seed+1)
			}
		})
	}
}

type auditRetainingRedeclarationHandler struct {
	scope        *Scope
	v            *Var
	node         *rhino.Node
	view         iter.Seq[*Var]
	implicit     *Var
	afterCapture func()
}

func (h *auditRetainingRedeclarationHandler) OnRedeclaration(scope *Scope, name string, _ *rhino.Node, _ *CompilerInput) {
	if h.scope == nil {
		h.scope, h.v = scope, scope.GetVar(name)
		h.node = h.v.GetNode()
		h.view = scope.GetAllSymbols()
		h.implicit = scope.MakeImplicitVar(AbstractScope_ImplicitVar_ARGUMENTS)
		if h.afterCapture != nil {
			h.afterCapture()
		}
	}
}

func TestInferConstsPreservesScopeRetainedByDefaultHandler(t *testing.T) {
	for _, mode := range []string{"return", "panic", "reenter"} {
		t.Run(mode, func(t *testing.T) {
			compiler := NewCompiler()
			compiler.InitOptions(NewCompilerOptions())
			h := &auditRetainingRedeclarationHandler{}
			previous := SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER
			t.Cleanup(func() { SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER = previous })
			SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER = h
			js := rhino.IR_Root(rhino.IR_ScriptOfNodes(
				rhino.IR_Var(rhino.IR_Name("a")), rhino.IR_Var(rhino.IR_Name("a"))))
			sentinel := jutil.NewRuntimeException("scope-owner callback panic")
			if mode == "panic" {
				h.afterCapture = func() { panic(sentinel) }
			} else if mode == "reenter" {
				h.afterCapture = func() {
					NewInferConsts(compiler).Process(rhino.IR_Root(), rhino.IR_Root(rhino.IR_ScriptOfNodes(rhino.IR_Var(rhino.IR_Name("inner")))))
				}
			}
			func() {
				if mode == "panic" {
					defer func() {
						got := recover()
						wrapped, ok := got.(jutil.IRuntimeException)
						// InferConsts enters NodeTraversal, whose established public
						// contract wraps callback failures with node context.
						if !ok || wrapped.GetCause() != sentinel || !strings.HasPrefix(wrapped.GetMessage(), "INTERNAL COMPILER ERROR.\nPlease report this problem.\n\n"+sentinel.GetMessage()) {
							t.Fatalf("callback panic wrapping/cause changed: %v", got)
						}
					}()
				}
				NewInferConsts(compiler).Process(rhino.IR_Root(), js)
				if mode == "panic" {
					t.Fatal("callback panic was swallowed")
				}
			}()
			check := func() {
				t.Helper()
				if h.scope == nil || h.v == nil || h.scope.GetRootNode() != js || h.v.GetName() != "a" || h.v.GetNode() != h.node || h.scope.GetVar("a") != h.v {
					t.Fatal("pass release invalidated scope or variable retained through the public handler")
				}
				if h.implicit == nil || h.implicit.GetName() != "arguments" || h.implicit.GetScope() != h.scope {
					t.Fatal("pass release invalidated an implicit variable retained by the callback")
				}
				var got []*Var
				for v := range h.view {
					got = append(got, v)
				}
				if len(got) != 1 || got[0] != h.v {
					t.Fatal("pass release invalidated the public symbol view")
				}
			}
			check()
			// A later pass must not overwrite the published backing slots.
			NewInferConsts(compiler).Process(rhino.IR_Root(), rhino.IR_Root(rhino.IR_ScriptOfNodes(rhino.IR_Var(rhino.IR_Name("other")))))
			check()
		})
	}
}

type auditMemoizedReentrantScope struct {
	*Scope
	beforeParent func()
}

func (s *auditMemoizedReentrantScope) GetParentScope() rhino.StaticScope {
	if f := s.beforeParent; f != nil {
		s.beforeParent = nil
		f()
	}
	return s.Scope.GetParentScope()
}

func TestMemoizedScopeCreatorLastHitSurvivesCustomParentReentry(t *testing.T) {
	delegate := &auditMemoizedScopeDelegate{makeScope: func(n *rhino.Node) AbstractScope {
		return &auditMemoizedReentrantScope{Scope: Scope_CreateGlobalScope(n)}
	}}
	creator := NewMemoizedScopeCreator(delegate)
	root, other := rhino.IR_Root(), rhino.IR_Root()
	first := creator.CreateScope(root, nil).(*auditMemoizedReentrantScope)
	var nested AbstractScope
	first.beforeParent = func() { nested = creator.CreateScope(other, nil) }
	if got := creator.CreateScope(root, nil); got != first || nested == nil || nested == first || delegate.calls != 2 {
		t.Fatal("custom parent lookup reentry replaced the outer memoized scope")
	}
}

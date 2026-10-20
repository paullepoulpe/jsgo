package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestInlineVariablesPublicNoInlineNamedExpressions(t *testing.T) {
	for _, source := range []string{
		"/** @noinline */ var outer = function inner() { return inner; };",
		"/** @noinline */ var outer = class Inner {};",
	} {
		t.Run(source, func(t *testing.T) {
			fixture := newSyntacticScopeCreatorTest(t)
			script := fixture.getRoot(source)
			expression := script.GetFirstChild().GetFirstChild().GetFirstChild()
			name := expression.GetFirstChild()
			if expression.GetJSDocInfo() != nil || name.GetJSDocInfo() != nil {
				t.Fatal("fixture must inherit its annotation from the enclosing variable declaration")
			}
			info := NodeUtil_GetBestJSDocInfo(name)
			if info == nil || !info.IsNoInline() {
				t.Fatal("fixture does not exercise the existing recursive JSDoc path")
			}
			scope := Scope_CreateChildScope(Scope_CreateGlobalScope(script), expression)
			variable := scope.Declare(name.GetString(), name, nil)
			if !InlineVariables_HasNoInlineAnnotation(variable) {
				t.Fatal("named expression lost its enclosing @noinline annotation")
			}
		})
	}
}

func TestOptimizeCallsPublicReferenceMapOwnsGlobalScope(t *testing.T) {
	for _, mode := range []string{"deferred-getter", "saved-var", "callback-panic"} {
		t.Run(mode, func(t *testing.T) {
			fixture := newSyntacticScopeCreatorTest(t)
			script := fixture.getRoot("var keep = 1; keep;")
			externs := rhino.IR_Root()
			combined := rhino.IR_Root(externs, script)
			fixture.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
			fixture.compiler.SetAccessorSummary(newAccessorSummaryAssuming(false))
			var retained *OptimizeCalls_ReferenceMap
			var savedScope *Scope
			var savedVar *Var
			calls := 0
			pass := NewOptimizeCalls(fixture.compiler, []OptimizeCalls_CallGraphCompilerPass{
				OptimizeCalls_CallGraphCompilerPassFunc(func(_, _ *rhino.Node, references *OptimizeCalls_ReferenceMap) {
					calls++
					retained = references
					if mode == "saved-var" {
						savedScope = references.GetGlobalScope()
						savedVar = savedScope.GetVar("keep")
						if savedVar == nil {
							t.Fatal("fixture did not collect the declared global")
						}
					}
					if mode == "callback-panic" {
						panic("retained-map")
					}
				}),
			}, true)
			var caught any
			func() {
				defer func() { caught = recover() }()
				pass.Process(externs, script)
			}()
			if mode == "callback-panic" {
				if caught != "retained-map" {
					t.Fatalf("panic = %v, want retained-map", caught)
				}
			} else if caught != nil {
				t.Fatalf("unexpected panic: %v", caught)
			}
			if calls != 1 || retained == nil {
				t.Fatal("custom pass was not called exactly once")
			}
			// Deliberately obtain the scope only after Process has returned or
			// unwound. Publishing the map must protect this future observation.
			scope := retained.GetGlobalScope()
			if scope == nil || scope.GetRootNode() != combined {
				t.Fatal("retained map's global scope changed at pass exit")
			}
			variable := scope.GetVar("keep")
			if variable == nil || variable.GetName() != "keep" {
				t.Fatal("retained global Var was cleared at pass exit")
			}
			if mode == "saved-var" && (scope != savedScope || variable != savedVar || savedVar.GetNameNode() != script.GetFirstChild().GetFirstChild()) {
				t.Fatal("published scope or Var identity changed at pass exit")
			}
		})
	}
}

type round9ConventionGetterCompiler struct {
	*Compiler
	calls int
}

func (c *round9ConventionGetterCompiler) GetCodingConvention() CodingConvention {
	c.calls++
	return c.Compiler.GetCodingConvention()
}

func TestPassPublicCodingConventionGetterTiming(t *testing.T) {
	t.Run("inline-expert", func(t *testing.T) {
		previous, _, variable, references, _, _ := inlineExpertOwnershipFixture()
		compiler := &round9ConventionGetterCompiler{Compiler: previous.outer.compiler.(*Compiler)}
		behavior := NewInlineVariables_InliningBehavior(NewInlineVariables(compiler, InlineVariables_Mode_ALL))
		if compiler.calls != 0 {
			t.Errorf("constructor called convention getter %d times", compiler.calls)
		}
		compiler.calls = 0
		behavior.CreateVarExpert(variable, references, false)
		if compiler.calls != 2 {
			t.Errorf("expert convention getter calls = %d, want 2", compiler.calls)
		}
	})
	t.Run("reference-callback", func(t *testing.T) {
		fixture := newSyntacticScopeCreatorTest(t)
		compiler := &round9ConventionGetterCompiler{Compiler: fixture.compiler}
		callback := NewOptimizeCalls_ReferenceMapBuildingCallback(NewOptimizeCalls(compiler, nil, true), NewOptimizeCalls_ReferenceMap())
		if compiler.calls != 0 {
			t.Errorf("constructor called convention getter %d times", compiler.calls)
		}
		compiler.calls = 0
		callback.Visit(nil, rhino.IR_Call(rhino.IR_Name("ordinaryCall")), nil)
		callback.Visit(nil, rhino.IR_Call(rhino.IR_Name("anotherCall")), nil)
		if compiler.calls != 2 {
			t.Errorf("CALL convention getter calls = %d, want 2", compiler.calls)
		}
	})
}

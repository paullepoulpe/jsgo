package jscomp

import (
	"fmt"
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Compare complete fixed-point solutions against the original AST-interpreted
// transfer functions, not only the resulting optimized JavaScript. This also
// checks dependency metadata that lattice Equals deliberately ignores.
func TestNativeFunctionTransferFactsDifferential(t *testing.T) {
	cases := []string{
		`function f(a,b) { var x=a; var y=b; x=y+1; use(x); }`,
		`function f(a,b) { var x=a; while(b) { x=x+1; b--; } return x; }`,
		`function f(a,b) { var x=a; if (b) { x=b; } else { x=2; } return x; }`,
		`function f(a,b) { var x=a; var y=x+b; x=2; return y; }`,
		`function f(a,b) { var x=a; var y=x+b; arguments[0]=b; return y; }`,
		`function f(a,b) { var x=a; var y=x+b; use(arguments); return y; }`,
		`function f(a,b) { var x=a; var y=x+outside; b && (x=y); b ? (x=a) : (x=b); return x; }`,
		`function f(a,b) { var x=a; a?.[x=b]; b?.(x=a); return x; }`,
		`function f(a,b) { var x=a; try { x=b(); use(x); } catch(e) { x=e; } finally { use(a); } return x; }`,
		`function f(a,b) { var x=a; for(x of b) { use(x); } return x; }`,
		`function f(a,b) { let x=a; for(let y of b) { x=y; use(x); } return x; }`,
		`function f(a,b) { var x=a; var y=b; [x,y=a] = b; return x+y; }`,
		`function f(a,b) { let {x,y=a}=b; x += a; return x+y; }`,
		`function f(a,b) { let x=a; use(() => x); var y=b; y++; return y; }`,
		`function f(a,b) { var x=a; var y=b; x = (y=a, y+1); return x+y; }`,
		`function f(a,b) { var x=a; switch(b) { case 1: x=b; break; default: x=a; } return x; }`,
	}
	// More than 64 variables exercises the existing sparse lattice fallback.
	large := "function f(a,b) {"
	for i := 0; i < 70; i++ {
		large += fmt.Sprintf("var x%d=a;", i)
	}
	cases = append(cases, large+"while(b) { x0=x69+b; b--; } return x0; }")
	for i, src := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			compiler, scope := functionFactsTestScope(t, src)
			creator := NewSyntacticScopeCreator(compiler)
			all, escaped := flowSensitiveInlineVariables_collectFunctionFacts(compiler, creator, scope)
			vars := all.GetAllVariables()
			cfg := ControlFlowAnalysis_Builder_().SetCompiler(compiler).SetCfgRoot(scope.GetRootNode()).SetIncludeEdgeAnnotations(true).ComputeCfg()
			facts := newFunctionTransferFacts(vars, escaped)
			defer facts.invalidate()

			genericMust := NewMustBeReachingVariableDef(cfg, compiler, escaped, vars)
			genericMust.Analyze()
			mustStates := make(map[*rhino.Node]*DataFlowAnalysis_LinearFlowState[*MustBeReachingVariableDef_MustDef])
			for _, n := range cfg.GetDiGraphNodes() {
				mustStates[n.GetValue()] = dataFlowAnalysis_linearFlowStateOf[*MustBeReachingVariableDef_MustDef](n)
			}
			nativeMust := NewMustBeReachingVariableDef(cfg, compiler, escaped, vars)
			nativeMust.nativeFacts = facts
			nativeMust.Analyze()
			for _, n := range cfg.GetDiGraphNodes() {
				want := mustStates[n.GetValue()]
				got := dataFlowAnalysis_linearFlowStateOf[*MustBeReachingVariableDef_MustDef](n)
				if got.GetStepCount() != want.GetStepCount() {
					t.Fatal("must-def convergence schedule changed")
				}
				nativeFactsCompareMust(t, want.GetIn(), got.GetIn())
				nativeFactsCompareMust(t, want.GetOut(), got.GetOut())
			}

			genericUses := NewMaybeReachingVariableUse(cfg, escaped, vars)
			genericUses.Analyze()
			useStates := make(map[*rhino.Node]*DataFlowAnalysis_LinearFlowState[*MaybeReachingVariableUse_ReachingUses])
			for _, n := range cfg.GetDiGraphNodes() {
				useStates[n.GetValue()] = dataFlowAnalysis_linearFlowStateOf[*MaybeReachingVariableUse_ReachingUses](n)
			}
			nativeUses := NewMaybeReachingVariableUse(cfg, escaped, vars)
			nativeUses.nativeFacts = facts
			nativeUses.Analyze()
			for _, n := range cfg.GetDiGraphNodes() {
				want := useStates[n.GetValue()]
				got := dataFlowAnalysis_linearFlowStateOf[*MaybeReachingVariableUse_ReachingUses](n)
				if got.GetStepCount() != want.GetStepCount() {
					t.Fatal("may-use convergence schedule changed")
				}
				if !got.GetIn().Equals(want.GetIn()) || !got.GetOut().Equals(want.GetOut()) {
					t.Fatal("may-use fixed point changed")
				}
				for v := range vars.Values() {
					if !slices.Equal(slices.Collect(want.GetIn().Get(v)), slices.Collect(got.GetIn().Get(v))) || !slices.Equal(slices.Collect(want.GetOut().Get(v)), slices.Collect(got.GetOut().Get(v))) {
						t.Fatal("may-use order changed")
					}
				}
			}
			if len(facts.must) == 0 || len(facts.uses) == 0 {
				t.Fatal("native programs were not exercised")
			}
			for i, v := range facts.vars {
				if facts.symbols[v.GetName()] != uint32(i) || vars.Get(v.GetName()) != v {
					t.Fatal("canonical Var identity changed")
				}
			}
		})
	}
}

func nativeFactsCompareMust(t *testing.T, want, got *MustBeReachingVariableDef_MustDef) {
	t.Helper()
	var wantKeys, gotKeys []*Var
	for v := range want.reachingDef.All() {
		wantKeys = append(wantKeys, v)
	}
	for v := range got.reachingDef.All() {
		gotKeys = append(gotKeys, v)
	}
	if !slices.Equal(wantKeys, gotKeys) {
		t.Fatal("must-def canonical keys or insertion order changed")
	}
	for _, v := range wantKeys {
		a, b := want.reachingDef.Get(v), got.reachingDef.Get(v)
		if (a == nil) != (b == nil) {
			t.Fatal("must-def TOP/BOTTOM changed")
		}
		if a != nil && (a.node != b.node || a.unknownDependencies != b.unknownDependencies || !slices.Equal(slices.Collect(a.depends.All()), slices.Collect(b.depends.All()))) {
			t.Fatal("must-def dependency metadata changed")
		}
	}
}

func TestNativeFunctionTransferFactsInvalidation(t *testing.T) {
	compiler, scope := functionFactsTestScope(t, `function f(a) { var x=a; x=a; return x; }`)
	creator := NewSyntacticScopeCreator(compiler)
	all, escaped := flowSensitiveInlineVariables_collectFunctionFacts(compiler, creator, scope)
	vars := all.GetAllVariables()
	cfg := ControlFlowAnalysis_Builder_().SetCompiler(compiler).SetCfgRoot(scope.GetRootNode()).SetIncludeEdgeAnnotations(true).ComputeCfg()
	d := NewMustBeReachingVariableDef(cfg, compiler, escaped, vars)
	facts := newFunctionTransferFacts(vars, escaped)
	d.nativeFacts = facts
	var node *rhino.Node
	for _, n := range cfg.GetDiGraphNodes() {
		v := n.GetValue()
		if v != nil && v.IsExprResult() && v.GetFirstChild().IsAssign() {
			node = v
		}
	}
	if node == nil {
		t.Fatal("assignment CFG node missing")
	}
	input := d.CreateEntryLattice()
	before := d.FlowThrough(node, input)
	if before.reachingDef.Get(vars.Get("x")).unknownDependencies {
		t.Fatal("unexpected unknown dependency before mutation")
	}
	facts.invalidate()
	if facts.active || facts.must != nil || facts.uses != nil || facts.vars != nil {
		t.Fatal("generation retained programs after invalidation")
	}
	node.GetFirstChild().GetLastChild().SetString("outside")
	got := d.FlowThrough(node, input)
	generic := NewMustBeReachingVariableDef(cfg, compiler, escaped, vars)
	want := generic.FlowThrough(node, input)
	nativeFactsCompareMust(t, want, got)
	if !got.reachingDef.Get(vars.Get("x")).unknownDependencies {
		t.Fatal("invalidated context used stale dependencies")
	}
	fresh := newFunctionTransferFacts(vars, escaped)
	d.nativeFacts = fresh
	defer fresh.invalidate()
	nativeFactsCompareMust(t, want, d.FlowThrough(node, input))
}

// Escaped membership is captured once, by ScopedName equality, so equal aliases
// remain escaped without replacing canonical Vars in either real consumer.
func TestNativeFunctionTransferFactsEscapedAlias(t *testing.T) {
	compiler, scope := functionFactsTestScope(t, `function f(a) { var x=a; return x; }`)
	all, _ := flowSensitiveInlineVariables_collectFunctionFacts(compiler, NewSyntacticScopeCreator(compiler), scope)
	vars := all.GetAllVariables()
	x := vars.Get("x")
	alias := NewVar(x.GetName(), x.GetNode(), x.GetScope().(*Scope), 0, nil, nil)
	escaped := jutil.NewObjectLinkedHashSet[*Var]()
	escaped.Add(alias)
	facts := newFunctionTransferFacts(vars, escaped)
	defer facts.invalidate()
	if !facts.escaped[facts.symbols["x"]] || facts.vars[facts.symbols["x"]] != x {
		t.Fatal("escaped alias lost canonical identity")
	}
}

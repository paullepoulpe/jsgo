package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Candidate-only activation control; do not overlay this file onto pre-reuse
// stock sources. Public semantic controls live in cfg_public_ownership_test.go.
func TestCFGOwnershipStockCoalesceRetainsPrivateReuse(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("function first(a) { var x = a; return x; } function second(b) { var y = b; return y; }")
	fixture.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	pass := NewCoalesceVariableNames(fixture.compiler, false)
	pass.Process(rhino.IR_Root(), root)
	if pass.cfgPublished {
		t.Fatal("exact stock pass accidentally published its private CFG")
	}
	if pass.scopeCreator.cache == nil || pass.scopeCreator.cache.escaped {
		t.Fatal("stock finder accidentally published its memoized creator cache")
	}
	delegate, ok := pass.scopeCreator.delegate.(*SyntacticScopeCreator)
	if !ok || delegate.arenaEscaped || delegate.arena != nil {
		t.Fatal("stock finder disabled ordinary scope arena release")
	}
	if len(pass.scopeStates) != 1 {
		t.Fatalf("sequential sibling functions used %d analysis slots, want 1", len(pass.scopeStates))
	}
	last := &pass.scopeStates[0].liveness
	if last.GetAllVariables().Get("b") == nil || last.GetAllVariables().Get("a") != nil {
		t.Fatal("private analysis slot did not process/reset for the second sibling")
	}
	found := false
	for _, cache := range pass.reusableCFAs {
		if cache != nil && last.GetCfg() == &cache.cfg {
			found = true
		}
	}
	if !found {
		t.Fatal("real stock traversal did not use the private reusable CFG cache")
	}
}

// Exercise the transfer boundary itself, in addition to the full-pass
// activation above and the stock-compatible public lifetime controls.
func TestCFGOwnershipBorrowedTransferRetainsGraph(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("function first(a) { return a; } function second(b) { return b; }")
	pass := NewCheckUnreachableCode(fixture.compiler)
	nt := &NodeTraversal{callback: pass}
	pass.cfgs = append(pass.cfgs, root.GetFirstChild())
	borrowed := pass.getControlFlowGraph(nt, fixture.compiler)
	entry := borrowed.GetEntry()
	value := entry.GetValue()
	marker := &struct{ name string }{"borrowed-then-published"}
	entry.SetAnnotation(marker)
	if len(pass.reusableCFAs) != 1 || pass.reusableCFAs[0] == nil {
		t.Fatal("fixture did not activate the exact stock borrowed path")
	}
	published := pass.GetControlFlowGraph(fixture.compiler)
	if published != borrowed || published.GetEntry() != entry || pass.reusableCFAs[0] != nil {
		t.Fatal("publication did not transfer the exact existing graph/cache slot")
	}
	pass.cfgs[0] = root.GetFirstChild().GetNext()
	next := pass.getControlFlowGraph(nt, fixture.compiler)
	if next == published || published.GetEntry() != entry || entry.GetValue() != value || entry.GetAnnotation() != marker {
		t.Fatal("later sibling reset transferred graph storage")
	}
	pass.InitNodeTraversal_AbstractCfgCallbackBase(pass)
	pass.cfgs = append(pass.cfgs, root.GetFirstChild().GetNext())
	if pass.canReuseControlFlowGraph(nt, fixture.compiler) || pass.getControlFlowGraph(nt, fixture.compiler) == published {
		t.Fatal("reinitialization forgot the ownership publication boundary")
	}
	if entry.GetValue() != value || entry.GetAnnotation() != marker {
		t.Fatal("reinitialization invalidated retained graph values")
	}
}

// Candidate-only because ReleaseArena is new. Publishing before any arena is
// allocated must also protect scopes allocated later by the same creator.
func TestCFGOwnershipPublishedCreatorsSurviveReleaseAndFutureScopes(t *testing.T) {
	for _, memoized := range []bool{false, true} {
		name := "syntactic"
		if memoized {
			name = "memoized"
		}
		t.Run(name, func(t *testing.T) {
			fixture := newSyntacticScopeCreatorTest(t)
			ssc := NewSyntacticScopeCreator(fixture.compiler)
			var creator ScopeCreator = ssc
			release := ssc.ReleaseArena
			if memoized {
				memo := NewMemoizedScopeCreator(ssc)
				creator, release = memo, memo.ReleaseArena
			}
			if ssc.arena != nil {
				t.Fatal("fixture acquired an arena before publication")
			}
			nt := &NodeTraversal{scopeCreator: creator}
			if nt.GetScopeCreator() != creator {
				t.Fatal("provided creator identity changed")
			}
			makeBody := func(root *rhino.Node) *Scope {
				global := creator.CreateScope(root, nil)
				function := creator.CreateScope(root.GetFirstChild(), global)
				return creator.CreateScope(root.GetFirstChild().GetLastChild(), function).Untyped()
			}
			firstRoot := fixture.getRoot("function first(a) { var x = a; return x; }")
			first := makeBody(firstRoot)
			original := first.GetVar("x")
			if original == nil {
				t.Fatal("fixture has no retained local")
			}
			originalNode := original.GetNode()
			release()
			if memoized && creator.CreateScope(first.GetRootNode(), first.GetParent()) != first {
				t.Fatal("release cleared a published memoized scope identity")
			}
			second := makeBody(fixture.getRoot("function second(b) { var x = b; return x; }"))
			fresh := second.GetVar("x")
			if fresh == nil || fresh == original {
				t.Fatal("future allocation reused a published Var")
			}
			freshNode := fresh.GetNode()
			release()
			if first.GetVar("x") != original || original.GetNode() != originalNode || original.GetName() != "x" || second.GetVar("x") != fresh || fresh.GetNode() != freshNode {
				t.Fatal("release/future allocation changed published scope or Var values")
			}
		})
	}
}

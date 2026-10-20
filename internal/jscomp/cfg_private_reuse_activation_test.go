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

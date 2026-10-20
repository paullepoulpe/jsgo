package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// These positive controls use only APIs present before the reusable CFG change.
type cfgOwnershipObserver struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler *Compiler
	graphs   []ControlFlowGraph[*rhino.Node]
	entries  []DiGraph_DiGraphNode[*rhino.Node, ControlFlowGraph_Branch]
	values   []*rhino.Node
	markers  []*struct{ name string }
}

func (c *cfgOwnershipObserver) EnterScopeWithCfg(nt *NodeTraversal) {
	if !nt.GetScopeRoot().IsFunction() {
		return
	}
	g := c.GetControlFlowGraph(c.compiler)
	entry := g.GetEntry()
	marker := &struct{ name string }{nt.GetScopeRoot().GetFirstChild().GetString()}
	entry.SetAnnotation(marker)
	c.graphs = append(c.graphs, g)
	c.entries = append(c.entries, entry)
	c.values = append(c.values, entry.GetValue())
	c.markers = append(c.markers, marker)
}

func TestCFGPublicOwnershipAcrossSiblingScopes(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("function first(a) { return a; } function second(b) { if (b) return b; }")
	cb := &cfgOwnershipObserver{compiler: fixture.compiler}
	cb.InitNodeTraversal_AbstractCfgCallbackBase(cb)
	NodeTraversal_Traverse(fixture.compiler, root, cb)
	if len(cb.graphs) != 2 {
		t.Fatalf("captured %d function graphs, want 2", len(cb.graphs))
	}
	if cb.graphs[0] == cb.graphs[1] {
		t.Fatal("public sibling CFGs alias reusable storage")
	}
	for i, g := range cb.graphs {
		if g.GetEntry() != cb.entries[i] || g.GetEntry().GetValue() != cb.values[i] || cb.entries[i].GetAnnotation() != cb.markers[i] {
			t.Fatalf("retained graph %d changed entry/value/annotation after a sibling", i)
		}
	}
}

func TestCFGPublicTypeInferenceComputeCfgOwnsResult(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("function first(a) { return a; } function second(b) { if (b) return b; }")
	pass := &TypeInferencePass{compiler: fixture.compiler}
	first := pass.ComputeCfg(root.GetFirstChild())
	entry, value := first.GetEntry(), first.GetEntry().GetValue()
	marker := &struct{ name string }{"retained"}
	entry.SetAnnotation(marker)
	second := pass.ComputeCfg(root.GetFirstChild().GetNext())
	if first == second {
		t.Fatal("public ComputeCfg results alias reusable storage")
	}
	if first.GetEntry() != entry || entry.GetValue() != value || entry.GetAnnotation() != marker {
		t.Fatal("second ComputeCfg changed the first result")
	}
}

type cfgOwnershipLatticeSnapshot struct {
	state           *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice]
	in, out         *LiveVariablesAnalysis_LiveVariableLattice
	inText, outText string
}

type cfgOwnershipCoalesceObserver struct {
	*CoalesceVariableNames
	analyses      []*LiveVariablesAnalysis
	firstNames    []string
	firstVars     []*Var
	firstLattices []cfgOwnershipLatticeSnapshot
}

// Leave the base's self as the original concrete pass. The actual traversal
// callback is nevertheless an expert wrapper, so concrete self alone is not an
// ownership discriminator.
func (c *cfgOwnershipCoalesceObserver) EnterScope(nt *NodeTraversal) {
	c.CoalesceVariableNames.EnterScope(nt)
	if !nt.GetScopeRoot().IsFunction() {
		return
	}
	c.analyses = append(c.analyses, c.liveness)
	if len(c.analyses) != 1 {
		return
	}
	for _, v := range c.liveness.GetAllVariablesInOrder() {
		c.firstNames = append(c.firstNames, v.GetName())
		c.firstVars = append(c.firstVars, v)
	}
	// Publish the pass's analyzed graph through its existing public API.
	// Retaining its annotations must not borrow a later function's buffers.
	graph := c.GetControlFlowGraph(c.compiler)
	for _, node := range graph.GetNodes() {
		state := node.GetAnnotation().(*DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice])
		c.firstLattices = append(c.firstLattices, cfgOwnershipLatticeSnapshot{
			state: state, in: state.GetIn(), out: state.GetOut(),
			inText: state.GetIn().String(), outText: state.GetOut().String(),
		})
	}
}

func TestCFGPublicCoalesceExpertRetainsAnalysisAnnotations(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	// These simple functions already satisfy the unique-binding normalized
	// precondition; no unrelated optimization is needed to exercise ownership.
	root := fixture.getRoot("function first(a) { var x = a; return x; } function second(b) { var y = b; if (b) y = 2; return y; }")
	fixture.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	pass := NewCoalesceVariableNames(fixture.compiler, false)
	cb := &cfgOwnershipCoalesceObserver{CoalesceVariableNames: pass}
	builder := NodeTraversal_Builder{compiler: fixture.compiler, callback: cb, scopeCreator: pass.scopeCreator}
	builder.Traverse(root)
	if len(cb.analyses) != 2 {
		t.Fatalf("captured %d analyses, want 2", len(cb.analyses))
	}
	if cb.analyses[0] == cb.analyses[1] {
		t.Fatal("expert callback's retained liveness was reset for its sibling")
	}
	first := cb.analyses[0]
	for i, name := range cb.firstNames {
		if first.GetAllVariables().Get(name) != cb.firstVars[i] || first.GetAllVariablesInOrder()[i] != cb.firstVars[i] {
			t.Fatal("retained analysis Var map/order changed")
		}
	}
	if len(cb.firstLattices) == 0 {
		t.Fatal("fixture did not capture dataflow annotations")
	}
	for _, saved := range cb.firstLattices {
		if saved.state.GetIn() != saved.in || saved.state.GetOut() != saved.out || saved.in.String() != saved.inText || saved.out.String() != saved.outText {
			t.Fatal("sibling analysis changed retained flow-state/lattice values")
		}
	}
}

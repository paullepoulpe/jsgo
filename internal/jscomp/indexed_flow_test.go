package jscomp

import (
	"cmp"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type indexedFlowTestLattice struct{ bits uint64 }

func (s *indexedFlowTestLattice) Equals(other any) bool {
	x, ok := other.(*indexedFlowTestLattice)
	return ok && x != nil && s.bits == x.bits
}

type indexedFlowTestAnalysis struct {
	DataFlowAnalysisBase[int, *indexedFlowTestLattice]
	forward, branched bool
	events            []string
}

func (a *indexedFlowTestAnalysis) IsForward() bool  { return a.forward }
func (a *indexedFlowTestAnalysis) IsBranched() bool { return a.branched }
func (a *indexedFlowTestAnalysis) CreateInitialEstimateLattice() *indexedFlowTestLattice {
	a.events = append(a.events, "initial")
	return &indexedFlowTestLattice{}
}
func (a *indexedFlowTestAnalysis) CreateEntryLattice() *indexedFlowTestLattice {
	a.events = append(a.events, "entry")
	return &indexedFlowTestLattice{bits: 1 << 63}
}
func (a *indexedFlowTestAnalysis) FlowThrough(node int, input *indexedFlowTestLattice) *indexedFlowTestLattice {
	a.events = append(a.events, fmt.Sprintf("flow:%d:%x", node, input.bits))
	return &indexedFlowTestLattice{bits: input.bits | 1<<uint(node%32)}
}
func (a *indexedFlowTestAnalysis) CreateFlowBrancher(node int, output *indexedFlowTestLattice) DataFlowAnalysis_FlowBrancher[*indexedFlowTestLattice] {
	a.events = append(a.events, fmt.Sprintf("brancher:%d:%x", node, output.bits))
	return func(branch ControlFlowGraph_Branch) *indexedFlowTestLattice {
		a.events = append(a.events, fmt.Sprintf("branch:%d:%d", node, branch))
		return &indexedFlowTestLattice{bits: output.bits | 1<<uint(48+branch)}
	}
}
func (a *indexedFlowTestAnalysis) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*indexedFlowTestLattice] {
	a.events = append(a.events, "joiner")
	return &indexedFlowTestJoiner{analysis: a}
}

type indexedFlowTestJoiner struct {
	analysis *indexedFlowTestAnalysis
	bits     uint64
}

func (j *indexedFlowTestJoiner) JoinFlow(input *indexedFlowTestLattice) {
	j.analysis.events = append(j.analysis.events, fmt.Sprintf("join:%x", input.bits))
	j.bits |= input.bits
}
func (j *indexedFlowTestJoiner) Finish() *indexedFlowTestLattice {
	j.analysis.events = append(j.analysis.events, "finish")
	return &indexedFlowTestLattice{bits: j.bits}
}

type indexedFlowTestCFG struct {
	*ControlFlowGraphBase[int]
	priorities bool
}

func (g *indexedFlowTestCFG) GetOptionalNodeComparator(forward bool) func(a, b DiGraph_DiGraphNode[int, ControlFlowGraph_Branch]) int {
	if !g.priorities {
		return nil
	}
	return func(a, b DiGraph_DiGraphNode[int, ControlFlowGraph_Branch]) int {
		result := cmp.Compare(a.GetValue()%3, b.GetValue()%3)
		if !forward {
			result = -result
		}
		return result
	}
}

func TestIndexedFlowGraphEngineDifferential(t *testing.T) {
	for seed := int64(0); seed < 20; seed++ {
		for _, prioritized := range []bool{false, true} {
			for mode := 0; mode < 3; mode++ { // forward, backward, branched forward
				t.Run(fmt.Sprintf("seed=%d/priority=%t/mode=%d", seed, prioritized, mode), func(t *testing.T) {
					rng := rand.New(rand.NewSource(seed))
					cfg := &indexedFlowTestCFG{ControlFlowGraphBase: NewControlFlowGraph(1, true, true), priorities: prioritized}
					count := 3 + rng.Intn(18)
					for value := 2; value <= count; value++ {
						cfg.CreateDiGraphNode(value)
					}
					// Parallel edges, exception branches, cycles, unreachable nodes
					// and arbitrary edge insertion order all occur in this corpus.
					for edge := 0; edge < count*3; edge++ {
						from, to := 1+rng.Intn(count), rng.Intn(count+1)
						cfg.Connect(from, ControlFlowGraph_Branch(rng.Intn(5)), to)
					}
					generic := &indexedFlowTestAnalysis{forward: mode != 1, branched: mode == 2}
					generic.InitDataFlowAnalysisBase(generic, cfg)
					generic.analyzeGeneric()
					type stateSnapshot struct {
						in, out uint64
						steps   int32
					}
					var want []stateSnapshot
					for _, node := range cfg.GetDiGraphNodes() {
						state := dataFlowAnalysis_linearFlowStateOf[*indexedFlowTestLattice](node)
						want = append(want, stateSnapshot{state.in.bits, state.out.bits, state.stepCount})
					}
					var wantBranches []uint64
					if generic.branched {
						for _, edge := range cfg.GetDiGraphEdges() {
							wantBranches = append(wantBranches, dataFlowAnalysis_latticeOf[*indexedFlowTestLattice](edge).bits)
						}
					}
					cfg.ClearNodeAnnotations()
					cfg.ClearEdgeAnnotations()
					native := &indexedFlowTestAnalysis{forward: mode != 1, branched: mode == 2}
					native.InitDataFlowAnalysisBase(native, cfg)
					topology := newIndexedFlowTopology[int](cfg)
					if topology == nil {
						t.Fatal("supported test CFG did not produce an indexed topology")
					}
					result := newIndexedFlowResult(topology, native)
					result.compute()
					for _, node := range cfg.GetDiGraphNodes() {
						if node.GetAnnotation() != nil {
							t.Fatal("compute mutated graph annotations before publication")
						}
					}
					result.publish()
					for id, node := range cfg.GetDiGraphNodes() {
						state := dataFlowAnalysis_linearFlowStateOf[*indexedFlowTestLattice](node)
						got := stateSnapshot{state.in.bits, state.out.bits, state.stepCount}
						if got != want[id] {
							t.Fatalf("node %d got %+v, want %+v", id, got, want[id])
						}
					}
					if native.branched {
						for id, edge := range cfg.GetDiGraphEdges() {
							if got := dataFlowAnalysis_latticeOf[*indexedFlowTestLattice](edge).bits; got != wantBranches[id] {
								t.Fatalf("edge %d got %x, want %x", id, got, wantBranches[id])
							}
						}
					}
					if !slices.Equal(native.events, generic.events) {
						limit := min(len(native.events), len(generic.events))
						for i := 0; i < limit; i++ {
							if native.events[i] != generic.events[i] {
								t.Fatalf("event %d got %s, want %s", i, native.events[i], generic.events[i])
							}
						}
						t.Fatalf("event counts got %d, want %d", len(native.events), len(generic.events))
					}
				})
			}
		}
	}
}

type indexedMustDefEntry struct {
	variable *Var
	node     *rhino.Node
	bottom   bool
	unknown  bool
	depends  []*Var
}

func indexedMustDefEntries(state *MustBeReachingVariableDef_MustDef) []indexedMustDefEntry {
	var entries []indexedMustDefEntry
	for variable, definition := range state.reachingDef.All() {
		entry := indexedMustDefEntry{variable: variable, bottom: definition == nil}
		if definition != nil {
			entry.node, entry.unknown = definition.node, definition.unknownDependencies
			entry.depends = slices.Collect(definition.depends.All())
		}
		entries = append(entries, entry)
	}
	return entries
}

func TestIndexedFlowMustDefDifferential(t *testing.T) {
	cases := []string{
		"var x; if (param1) { x=1; } else { x=2; } return x;",
		"var x=param1,y=0; while (param2) { y=x; x=y+1; } return y;",
		"var x=1; try { x=param1(); } catch (e) { x=2; } finally { param2(x); } return x;",
		"var x=0,y=1; [x = y = 3] = param1; x &&= y; x ??= y; return x;",
		"var x=1; function nested(){ return x; } arguments[0]=3; return nested();",
		"var x=1; for (const [a,b] of param1) { if(a) continue; x=b; } return x;",
		strings.Repeat("var padding=0;", 80) + "var x=1; return x;",
	}
	// Force the >64-variable sparse/alias fallback as well as canonical IDs.
	var many strings.Builder
	for i := 0; i < 70; i++ {
		fmt.Fprintf(&many, "var value%d=%d;", i, i)
	}
	cases = append(cases, many.String()+"while(param1){ value1=value2; } return value1;")
	for id, source := range cases {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			tester := ReachingUseDefTester_Create()
			tester.ComputeReachingDef(source)
			analysis := tester.reachingDef
			analysis.analyzeGeneric()
			type snapshot struct {
				in, out []indexedMustDefEntry
				steps   int32
			}
			var want []snapshot
			for _, node := range analysis.cfg.GetDiGraphNodes() {
				state := dataFlowAnalysis_linearFlowStateOf[*MustBeReachingVariableDef_MustDef](node)
				want = append(want, snapshot{indexedMustDefEntries(state.in), indexedMustDefEntries(state.out), state.stepCount})
			}
			if !analysis.AnalyzeIndexed() {
				t.Fatal("MustDef CFG unexpectedly fell back")
			}
			for id, node := range analysis.cfg.GetDiGraphNodes() {
				state := dataFlowAnalysis_linearFlowStateOf[*MustBeReachingVariableDef_MustDef](node)
				got := snapshot{indexedMustDefEntries(state.in), indexedMustDefEntries(state.out), state.stepCount}
				if !reflect.DeepEqual(got, want[id]) {
					t.Fatalf("node %d (%v) state/order/dependencies differ", id, node.GetValue())
				}
			}
		})
	}
}

func TestIndexedFlowQueueDifferential(t *testing.T) {
	for _, prioritized := range []bool{false, true} {
		rng := rand.New(rand.NewSource(20261018))
		var compare func(int32, int32) int
		if prioritized {
			compare = func(a, b int32) int { return cmp.Compare(a%5, b%5) }
		}
		original := newDataFlowAnalysis_UniqueQueue(compare)
		native := indexedFlowQueue{compare: compare, queued: make([]bool, 31)}
		for step := 0; step < 10000; step++ {
			if !original.IsEmpty() && rng.Intn(3) == 0 {
				if got, want := native.remove(), original.RemoveFirst(); got != want {
					t.Fatalf("priority=%t/step=%d got %d, want %d", prioritized, step, got, want)
				}
			} else {
				id := int32(rng.Intn(len(native.queued)))
				native.add(id)
				original.Add(id)
			}
			if native.empty() != original.IsEmpty() {
				t.Fatalf("priority=%t/step=%d emptiness differs", prioritized, step)
			}
		}
	}
}

type indexedMayUseEntry struct {
	variable *Var
	uses     []*rhino.Node
}

func indexedMayUseEntries(state *MaybeReachingVariableUse_ReachingUses, analysis *MaybeReachingVariableUse) []indexedMayUseEntry {
	var entries []indexedMayUseEntry
	for variable := range analysis.allVarsInFn.Values() {
		entries = append(entries, indexedMayUseEntry{variable: variable, uses: slices.Collect(state.Get(variable))})
	}
	return entries
}

func TestIndexedFlowMaybeUseDifferential(t *testing.T) {
	cases := []string{
		"var x=1; if (param1) { param2(x); } else { x=2; } return x;",
		"var x=1,y=0; while (param1) { y=x; x=y+1; } return x+y;",
		"var x=1; try { param1(x); x=2; } catch(e) { param2(x); } finally { param2(x); } return x;",
		"var x=0,y=1; [x = y = 3] = param1; x &&= y; x ??= y; return x;",
	}
	for id, source := range cases {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			tester := ReachingUseDefTester_Create()
			tester.ComputeReachingUses(source, false)
			analysis := tester.reachingUse
			analysis.analyzeGeneric()
			type snapshot struct {
				in, out []indexedMayUseEntry
				steps   int32
			}
			var want []snapshot
			for _, node := range analysis.cfg.GetDiGraphNodes() {
				state := dataFlowAnalysis_linearFlowStateOf[*MaybeReachingVariableUse_ReachingUses](node)
				want = append(want, snapshot{indexedMayUseEntries(state.in, analysis), indexedMayUseEntries(state.out, analysis), state.stepCount})
			}
			if !analysis.AnalyzeIndexed() {
				t.Fatal("MaybeUse CFG unexpectedly fell back")
			}
			for id, node := range analysis.cfg.GetDiGraphNodes() {
				state := dataFlowAnalysis_linearFlowStateOf[*MaybeReachingVariableUse_ReachingUses](node)
				got := snapshot{indexedMayUseEntries(state.in, analysis), indexedMayUseEntries(state.out, analysis), state.stepCount}
				if !reflect.DeepEqual(got, want[id]) {
					t.Fatalf("node %d (%v) reaching-use state/order differs", id, node.GetValue())
				}
			}
		})
	}
}

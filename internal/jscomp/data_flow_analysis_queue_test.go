package jscomp

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// A membership flag changes only the duplicate-detection representation. In
// particular, equal-priority nodes must retain PriorityQueue's tie behavior,
// and a deque must retain insertion order even after removal and reinsertion.
func TestDataFlowAnalysisAnnotationQueueDifferential(t *testing.T) {
	type node struct {
		priority int
		queued   bool
	}
	for _, prioritized := range []bool{false, true} {
		for seed := int64(0); seed < 20; seed++ {
			t.Run(fmt.Sprintf("priority=%t/seed=%d", prioritized, seed), func(t *testing.T) {
				rng := rand.New(rand.NewSource(seed))
				nodes := make([]node, 31)
				for i := range nodes {
					nodes[i].priority = i % 5
				}
				var comparator func(*node, *node) int
				if prioritized {
					comparator = func(a, b *node) int { return a.priority - b.priority }
				}
				original := newDataFlowAnalysis_UniqueQueue(comparator)
				optimized := newDataFlowAnalysis_UniqueQueueWithMembership(comparator, func(n *node) *bool { return &n.queued })
				for step := 0; step < 1000; step++ {
					switch rng.Intn(10) {
					case 0:
						original.Clear()
						optimized.Clear()
					case 1, 2, 3:
						if got, want := optimized.RemoveFirst(), original.RemoveFirst(); got != want {
							t.Fatalf("step %d: removal got %p, want %p", step, got, want)
						}
					default:
						n := &nodes[rng.Intn(len(nodes))]
						original.Add(n)
						optimized.Add(n)
					}
					if optimized.IsEmpty() != original.IsEmpty() {
						t.Fatalf("step %d: emptiness differs", step)
					}
					for i := range nodes {
						if got, want := nodes[i].queued, original.seenSet.Contains(&nodes[i]); got != want {
							t.Fatalf("step %d: node %d membership got %t, want %t", step, i, got, want)
						}
					}
				}
				for !original.IsEmpty() {
					if got, want := optimized.RemoveFirst(), original.RemoveFirst(); got != want {
						t.Fatalf("drain: removal got %p, want %p", got, want)
					}
				}
				if !optimized.IsEmpty() {
					t.Fatal("optimized queue did not drain")
				}
			})
		}
	}
}

// Reinitialization must not read membership from annotations owned by an old
// analysis. They may have been cleared or replaced after an interruption.
func TestDataFlowAnalysisReinitializeAfterAnnotationReplacement(t *testing.T) {
	for _, prioritized := range []bool{false, true} {
		for _, replacement := range []string{"cleared", "other-lattice"} {
			t.Run(fmt.Sprintf("priority=%t/%s", prioritized, replacement), func(t *testing.T) {
				entry := &dataFlowAnalysisTest_DivergentAnalysis_Counter{}
				cfg := NewControlFlowGraph(entry, true, true)
				cfg.CreateDiGraphNode(&dataFlowAnalysisTest_DivergentAnalysis_Counter{})
				analysis := newDataFlowAnalysisTest_DivergentAnalysis(cfg)
				if prioritized {
					// A tied comparator also checks the priority queue reset path.
					analysis.workQueue.priorityQueue = jutil.NewPriorityQueue(func(a, b DiGraph_DiGraphNode[*dataFlowAnalysisTest_DivergentAnalysis_Counter, ControlFlowGraph_Branch]) int {
						return 0
					})
				}
				analysis.Initialize()
				if analysis.workQueue.IsEmpty() {
					t.Fatal("test requires interrupted work remaining")
				}
				cfg.ClearNodeAnnotations()
				if replacement == "other-lattice" {
					for _, node := range cfg.GetDiGraphNodes() {
						node.SetAnnotation(newDataFlowAnalysis_LinearFlowState(1, 2))
					}
				}
				analysis.Initialize()
				seen := make(map[DiGraph_DiGraphNode[*dataFlowAnalysisTest_DivergentAnalysis_Counter, ControlFlowGraph_Branch]]bool)
				for !analysis.workQueue.IsEmpty() {
					node := analysis.workQueue.RemoveFirst()
					if seen[node] || node == cfg.GetImplicitReturn() {
						t.Fatal("invalid or duplicate queued node after reinitialization")
					}
					seen[node] = true
					state := dataFlowAnalysis_linearFlowStateOf[*dataFlowAnalysisTest_DivergentAnalysis_Step](node)
					if state == nil || state.queued || state.GetStepCount() != 0 {
						t.Fatal("flow state was not replaced and dequeued correctly")
					}
				}
				if len(seen) != 2 {
					t.Fatalf("queued %d nodes, want 2", len(seen))
				}
			})
		}
	}
}

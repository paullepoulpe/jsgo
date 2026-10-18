package jscomp

import (
	"fmt"
	"math/rand"
	"testing"
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

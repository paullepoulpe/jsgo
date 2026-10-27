package jscomp

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: fast mode counts the AST sizes behind the optimization loop's
// size heuristic (PhaseOptimizer.Loop, useSizeHeuristicToStopOptimizationLoop)
// on worker goroutines. Counting is a pure read of the tree and the sums are
// exact, so the loop stops exactly where the sequential count would make it
// stop.

// fastLoopSizeWorkers is the number of workers for the loop's AST-size
// counts, or 1 to count sequentially.
func fastLoopSizeWorkers(compiler AbstractCompiler) int {
	opts := compiler.GetOptions()
	if !opts.isFastMode() {
		return 1
	}
	return parallelThreadCount(compiler, 1<<30)
}

// countAstSizes returns NodeUtil_CountAstSize of every node. With more than
// one worker the trees are cut into subtrees (a breadth-first frontier) that
// workers count; each node's size is 1 per expanded node plus the sizes of
// its frontier subtrees.
func countAstSizes(nodes []rhino.Node, workers int) []int {
	sizes := make([]int, len(nodes))
	if workers <= 1 {
		for i, n := range nodes {
			sizes[i] = NodeUtil_CountAstSize(n)
		}
		return sizes
	}
	type item struct {
		n     rhino.Node
		owner int32
	}
	target := 64 * workers
	frontier := make([]item, 0, len(nodes))
	for i, n := range nodes {
		frontier = append(frontier, item{n, int32(i)})
	}
	// Expand level by level until there is enough work to share.
	for depth := 0; len(frontier) < target && depth < 48; depth++ {
		next := make([]item, 0, 2*len(frontier))
		expanded := false
		for _, it := range frontier {
			if it.n.GetFirstChild().IsNil() {
				sizes[it.owner]++
				continue
			}
			expanded = true
			sizes[it.owner]++
			for c := it.n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				next = append(next, item{c, it.owner})
			}
		}
		frontier = next
		if !expanded {
			break
		}
	}
	if len(frontier) == 0 {
		return sizes
	}
	partial := make([][]int, workers)
	var next atomic.Int64
	var wg sync.WaitGroup
	const chunk = 16
	for w := 0; w < workers; w++ {
		partial[w] = make([]int, len(nodes))
		wg.Add(1)
		go func(acc []int) {
			defer wg.Done()
			for {
				start := int(next.Add(chunk)) - chunk
				if start >= len(frontier) {
					return
				}
				end := min(start+chunk, len(frontier))
				for _, it := range frontier[start:end] {
					acc[it.owner] += NodeUtil_CountAstSize(it.n)
				}
			}
		}(partial[w])
	}
	wg.Wait()
	for _, acc := range partial {
		for i, v := range acc {
			sizes[i] += v
		}
	}
	return sizes
}

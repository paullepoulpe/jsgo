package rhino

import (
	"runtime"
	"sync"
	"testing"
	"weak"
)

func TestNodeChunksDistinctConcurrent(t *testing.T) {
	const workers, count = 4, 1000
	groups := make([][]*Node, workers)
	var wg sync.WaitGroup
	for worker := range groups {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < count; i++ {
				groups[worker] = append(groups[worker], NewNode(Token_NAME))
			}
		}(worker)
	}
	wg.Wait()
	seen := make(map[*Node]bool, workers*count)
	for _, group := range groups {
		for _, n := range group {
			if seen[n] {
				t.Fatal("node allocation reused a live identity")
			}
			seen[n] = true
			if n.GetToken() != Token_NAME || n.GetLineno() != -1 || n.GetParent() != nil {
				t.Fatal("new node was not initialized")
			}
		}
	}
	runtime.KeepAlive(groups)
}

func discardedChunkNode() weak.Pointer[Node] {
	// Move beyond any chunk shared with package-initialization globals.
	for i := 0; i < 256; i++ {
		NewNode(Token_EMPTY)
	}
	n := NewNode(Token_NAME)
	n.AddChildToBack(NewNode(Token_NUMBER))
	return weak.Make(n)
}

func TestNodeChunksDoNotRootDiscardedAST(t *testing.T) {
	node := discardedChunkNode()
	for attempt := 0; attempt < 10; attempt++ {
		runtime.GC()
		if node.Value() == nil {
			return
		}
	}
	t.Fatal("the current allocation chunk kept a discarded AST alive")
}

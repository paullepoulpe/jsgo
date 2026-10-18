package rhino

import (
	"sync"
	"weak"
)

// Go-only: performance — batch node allocation without a permanent strong
// root. A weak reference lets the current chunk and its AST become garbage
// when callers release their nodes. The lock also supports independent
// parsers/compilers without sharing unsynchronized allocation state.
var nodeChunks struct {
	sync.Mutex
	current weak.Pointer[nodeChunk]
}

type nodeChunk struct {
	nodes [128]Node
	used  int
}

func newNodeInChunk(token Token) *Node {
	nodeChunks.Lock()
	chunk := nodeChunks.current.Value()
	if chunk == nil || chunk.used == len(chunk.nodes) {
		chunk = &nodeChunk{}
		nodeChunks.current = weak.Make(chunk)
	}
	n := &chunk.nodes[chunk.used]
	chunk.used++
	nodeChunks.Unlock()
	n.token = token
	n.linenoCharno = -1
	return n
}

// Node_AllocationBoundary prevents a partial chunk from linking the lifetimes
// of successive compilations or package-level constant ASTs. Existing nodes
// stay valid; only the allocation group for future nodes changes.
// Go-only: performance — call before starting an independent compilation.
func Node_AllocationBoundary() {
	nodeChunks.Lock()
	nodeChunks.current = weak.Pointer[nodeChunk]{}
	nodeChunks.Unlock()
}

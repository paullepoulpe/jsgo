package rhino

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Go-only: node storage.
//
// Nodes live in fixed-size chunks and link to each other with 32-bit ids
// instead of pointers (see Node), so that once the node carries no Go
// pointer at all the garbage collector never traces the AST. A *Node handed
// to passes is an interior pointer into a chunk; the accessors translate the
// ids back through a process-wide chunk table.
//
// Id scheme: id = chunkIndex<<12 | slot, with slot in 1..4095 (slot 0 of
// every chunk is the chunk header, so id 0 is nil and every valid id is
// non-zero). Chunk indices are handed out once per process and never
// reused: a stale id (one whose chunk was released) always hits an empty
// table slot and panics instead of reading another chunk's node.
//
// Lifetime: the table holds the chunk addresses as plain integers that the
// GC does not trace. A chunk is kept alive by its generation (nodeGen,
// strong list of chunks), and every chunk points back at its generation,
// so any *Node still held by anyone keeps its whole generation (and so
// every chunk its links can reach) alive. When nothing references a
// generation any more, its cleanup clears the table slots. A generation
// starts at NewNodeGeneration (one per Compiler, and one per compaction);
// nodes are always allocated in the current generation, whichever Compiler
// made it, which is only a lifetime grouping: ids are global and links
// across generations are valid as long as both generations are alive.
// Limit: a link from a live generation into a generation that was released
// (only possible for a Node that outlived the Compiler that made it, in a
// process that keeps compiling, i.e. tests) reads an empty slot and panics.
//
// Concurrency: allocation is lock-free (one atomic add on the current
// chunk's bump counter); chunk creation and table publication take a mutex
// and the table slot is written atomically, so goroutines parsing in
// parallel allocate concurrently and ids are stable for the whole process.

type nodeID uint32

const (
	nodeChunkShift = 12
	nodeChunkSlots = 1 << nodeChunkShift
	nodeChunkMask  = nodeChunkSlots - 1
	nodeChunkCount = 1 << (32 - nodeChunkShift)
	nodeSize       = unsafe.Sizeof(Node{})
)

// nodeChunkHeader occupies slot 0 of a chunk; it is exactly one node wide.
type nodeChunkHeader struct {
	owner *nodeGen
	index uint32
	n     atomic.Uint32 // bump allocator: the next free slot
	_     [nodeSize - 16]byte
}

type nodeChunk struct {
	hdr   nodeChunkHeader
	nodes [nodeChunkSlots - 1]Node
}

var (
	_ [unsafe.Sizeof(nodeChunkHeader{}) - nodeSize]byte
	_ [nodeSize - unsafe.Sizeof(nodeChunkHeader{})]byte
	_ [unsafe.Offsetof(nodeChunk{}.nodes) - nodeSize]byte
)

// nodeChunkTable maps a chunk index to the chunk's address. The slots are
// uintptr on purpose (not traced; see the lifetime note above) and are
// read and written through atomic pointer views.
var nodeChunkTable [nodeChunkCount]uintptr

func nodeChunkSlot(k uint32) *atomic.Pointer[nodeChunk] {
	return (*atomic.Pointer[nodeChunk])(unsafe.Pointer(&nodeChunkTable[k]))
}

// nodeAt translates an id into the node it names.
func nodeAt(id nodeID) *Node {
	if id == 0 {
		return nil
	}
	c := nodeChunkSlot(uint32(id) >> nodeChunkShift).Load()
	if c == nil {
		staleNodeID(id)
	}
	return (*Node)(unsafe.Add(unsafe.Pointer(c), uintptr(id&nodeChunkMask)*nodeSize))
}

//go:noinline
func staleNodeID(id nodeID) {
	panic(fmt.Sprintf("rhino.Node: stale node id %#x (its chunk was released)", uint32(id)))
}

// idOf returns the id of n (0 for nil). A Node that is not in a chunk (a
// by-value template Node) has no id and cannot be linked into a tree.
func idOf(n *Node) nodeID {
	if n == nil {
		return 0
	}
	if n.self == 0 {
		notInChunk(n)
	}
	return n.self
}

//go:noinline
func notInChunk(n *Node) {
	panic(fmt.Sprintf("rhino.Node: node %s is not in a chunk and cannot be linked", n.token))
}

// nodeGen is a generation of chunks that is released as a whole.
type nodeGen struct {
	mu     sync.Mutex
	cur    atomic.Pointer[nodeChunk]
	chunks []*nodeChunk
}

var (
	nodeCurrentGen atomic.Pointer[nodeGen]
	nodeInitialGen *nodeGen // the generation of the package-level sentinel nodes; never released
	nodeGenMu      sync.Mutex
	nodeChunkNext  atomic.Uint32 // the next chunk index; index 0 is never used
)

// NodeGeneration is the public handle of a generation (see Compiler).
type NodeGeneration struct{ g *nodeGen }

// NewNodeGeneration starts a new generation and makes it the one new nodes
// are allocated in.
func NewNodeGeneration() NodeGeneration {
	g := newNodeGen()
	nodeCurrentGen.Store(g)
	return NodeGeneration{g}
}

func newNodeGen() *nodeGen {
	g := &nodeGen{}
	g.addChunk()
	return g
}

// releaseNodeChunk clears the table slot of a chunk whose generation is no
// longer referenced (the generation's cleanup, one per chunk).
func releaseNodeChunk(k uint32) {
	nodeChunkSlot(k).Store(nil)
}

// addChunk allocates and publishes a new chunk for g (g.mu held or g not
// yet shared) and makes it the allocation chunk.
func (g *nodeGen) addChunk() *nodeChunk {
	k := nodeChunkNext.Add(1)
	if k >= nodeChunkCount-1 { // the last index is reserved for nodeForwardMarker
		panic("rhino.Node: node id space exhausted (2^32 nodes allocated in this process)")
	}
	c := new(nodeChunk)
	c.hdr.owner = g
	c.hdr.index = k
	c.hdr.n.Store(1) // slot 0 is the header
	g.chunks = append(g.chunks, c)
	nodeChunkSlot(k).Store(c)
	g.cur.Store(c)
	runtime.AddCleanup(g, releaseNodeChunk, k)
	return c
}

// currentGen returns the allocation generation, creating the initial one on
// first use (package-level sentinel nodes are allocated during init).
func currentGen() *nodeGen {
	if g := nodeCurrentGen.Load(); g != nil {
		return g
	}
	return initialGen()
}

//go:noinline
func initialGen() *nodeGen {
	nodeGenMu.Lock()
	defer nodeGenMu.Unlock()
	if g := nodeCurrentGen.Load(); g != nil {
		return g
	}
	g := newNodeGen()
	nodeInitialGen = g
	nodeCurrentGen.Store(g)
	return g
}

// allocNode returns a zeroed node in the current generation with its id set.
func allocNode() *Node {
	g := currentGen()
	for {
		c := g.cur.Load()
		i := c.hdr.n.Add(1) - 1
		if i < nodeChunkSlots {
			n := (*Node)(unsafe.Add(unsafe.Pointer(c), uintptr(i)*nodeSize))
			n.self = nodeID(c.hdr.index<<nodeChunkShift | i)
			return n
		}
		g.refill(c)
	}
}

//go:noinline
func (g *nodeGen) refill(full *nodeChunk) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.cur.Load() != full {
		return // another goroutine already replaced it
	}
	g.addChunk()
}

func (n *Node) nextNode() *Node   { return nodeAt(n.next) }
func (n *Node) firstNode() *Node  { return nodeAt(n.first) }
func (n *Node) parentNode() *Node { return nodeAt(n.parent) }
func (n *Node) prevNode() *Node   { return nodeAt(n.previous) }

// NodeChunksAllocated returns the number of chunks allocated so far in the
// process (each holds up to NodeChunkCapacity nodes); the difference between
// two readings bounds the nodes allocated in between.
func NodeChunksAllocated() int {
	return int(nodeChunkNext.Load())
}

// NodeChunkCapacity is the number of nodes a chunk holds.
const NodeChunkCapacity = nodeChunkSlots - 1

// NodeChunksLive returns the number of chunks whose table slot is still
// set, i.e. chunks of generations that have not been released (for tests:
// a generation is released once nothing references any of its nodes).
func NodeChunksLive() int {
	live := 0
	for k := uint32(1); k < nodeChunkNext.Load(); k++ {
		if nodeChunkSlot(k).Load() != nil {
			live++
		}
	}
	return live
}

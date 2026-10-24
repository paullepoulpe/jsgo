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
// Id scheme: id = chunkIndex<<9 | slot, with slot in 1..479 (slot 0 of
// every chunk is the chunk header, so id 0 is nil and every valid id is
// non-zero; slots 480..511 are never used). Chunk indices are handed out
// once per process and never reused: a stale id (one whose chunk was
// released) always hits an empty table slot and panics instead of reading
// another chunk's node.
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
	// A chunk holds nodeChunkNodes-1 nodes behind a one-node-wide header and
	// is sized to stay under Go's 32 KiB small-object limit: a small object
	// with pointers is accounted and scanned only up to its last pointer
	// (the header's owner, see below), while a large one counts whole, so
	// 480 x 56 B + the 8-byte malloc header fits the 27264-byte size class
	// with 1.4 % slack and the GC sees 8 bytes of it.
	nodeChunkShift = 9
	nodeChunkSlots = 1 << nodeChunkShift // id slot bits; slots nodeChunkNodes.. are never used
	nodeChunkNodes = 480
	nodeChunkMask  = nodeChunkSlots - 1
	nodeChunkCount = 1 << (32 - nodeChunkShift)
	nodeSize       = unsafe.Sizeof(Node{})
	nodeChunkClass = 27264 // the Go size class the chunk is allocated in
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
	nodes [nodeChunkNodes - 1]Node
}

var (
	_ [unsafe.Sizeof(nodeChunkHeader{}) - nodeSize]byte
	_ [nodeSize - unsafe.Sizeof(nodeChunkHeader{})]byte
	_ [unsafe.Offsetof(nodeChunk{}.nodes) - nodeSize]byte
	_ [nodeChunkClass - 8 - unsafe.Sizeof(nodeChunk{})]byte // fits its size class with the malloc header
	_ [unsafe.Sizeof(nodeChunk{}) + 8 - 24576 - 1]byte      // and is not a smaller class
)

// nodeChunkTable maps a chunk index to the chunk's address. The slots are
// uintptr on purpose (not traced; see the lifetime note above) and are
// read and written through atomic pointer views. 64 MiB of BSS, of which
// only the touched pages exist (8 B per chunk allocated).
var nodeChunkTable [nodeChunkCount]uintptr

// nodeChunkSlot is the read view of a slot. Slots are only ever WRITTEN
// as integers (setNodeChunkSlot): a pointer store would carry a GC write
// barrier that shades the slot's previous value, and a cleanup clears the
// slot of a chunk the GC has already freed.
func nodeChunkSlot(k uint32) *atomic.Pointer[nodeChunk] {
	return (*atomic.Pointer[nodeChunk])(unsafe.Pointer(&nodeChunkTable[k]))
}

func setNodeChunkSlot(k uint32, c *nodeChunk) {
	atomic.StoreUintptr(&nodeChunkTable[k], uintptr(unsafe.Pointer(c)))
}

// nodeAt translates an id into the node it names.
func nodeAt(id nodeID) *Node {
	if id == 0 {
		return nil
	}
	c := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&nodeChunkTable[uint32(id)>>nodeChunkShift])))
	if c == nil {
		panic("rhino.Node: stale node id (its chunk was released)")
	}
	return (*Node)(unsafe.Add(c, uintptr(id&nodeChunkMask)*nodeSize))
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
	// peers: generations this one has links into (see (*Node).ref). They
	// are kept alive with this one, since its nodes resolve ids into them.
	peers    map[*nodeGen]struct{}
	lastPeer atomic.Pointer[nodeGen]
}

// ref returns the id of y for a link from n to y, and ties n's generation
// to y's when they differ: a link is an id, which the GC does not trace,
// so the generation that holds the link must keep the target's alive. In
// one compilation every node is in one generation (or in a compaction's
// successor, which never links back); the cross-generation case is a
// process that compiles several times with interleaved compilers (a test
// harness parsing expected output with a second compiler while the first
// one is still being worked on).
func (n *Node) ref(y *Node) nodeID {
	if y == nil {
		return 0
	}
	if y.self == 0 {
		notInChunk(y)
	}
	if (n.self^y.self)>>nodeChunkShift != 0 {
		crossChunkRef(n, y)
	}
	return y.self
}

func crossChunkRef(n, y *Node) {
	a := nodeChunkSlot(uint32(n.self) >> nodeChunkShift).Load()
	b := nodeChunkSlot(uint32(y.self) >> nodeChunkShift).Load()
	if a == nil || b == nil {
		if a == nil {
			staleNodeID(n.self)
		}
		staleNodeID(y.self)
	}
	if ga, gb := a.hdr.owner, b.hdr.owner; ga != gb {
		ga.tie(gb)
	}
}

//go:noinline
func (g *nodeGen) tie(peer *nodeGen) {
	if g.lastPeer.Load() == peer {
		return
	}
	// Both directions: the parent link goes one way, the child links the
	// other. The initial generation (never released) never needs to hold a
	// peer, and must not (it would pin that peer forever).
	if g != nodeInitialGen {
		g.mu.Lock()
		if g.peers == nil {
			g.peers = make(map[*nodeGen]struct{})
		}
		g.peers[peer] = struct{}{}
		g.mu.Unlock()
	}
	if peer != nodeInitialGen {
		peer.mu.Lock()
		if peer.peers == nil {
			peer.peers = make(map[*nodeGen]struct{})
		}
		peer.peers[g] = struct{}{}
		peer.mu.Unlock()
	}
	g.lastPeer.Store(peer)
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
	setNodeChunkSlot(k, nil)
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
	setNodeChunkSlot(k, c)
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
		if i < nodeChunkNodes {
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
const NodeChunkCapacity = nodeChunkNodes - 1

// NodeChunksLive returns the number of chunks whose table slot is still
// set, i.e. chunks of generations that have not been released (for tests:
// a generation is released once nothing references any of its nodes).
func NodeChunksLive() int {
	live := 0
	for k := uint32(1); k < nodeChunkNext.Load(); k++ {
		// Read the slot as an integer: a released generation's chunks may
		// already be freed while their cleanup has not yet cleared the
		// slots, and a Go pointer to a freed object must never exist.
		if atomic.LoadUintptr(&nodeChunkTable[k]) != 0 {
			live++
		}
	}
	return live
}

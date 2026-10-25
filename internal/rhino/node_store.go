package rhino

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: node storage.
//
// A Node is a 32-bit id. The node's data (nodeData, 32 bytes, no Go
// pointers: links are ids too) lives in a NodeStore, which belongs to one
// compilation, so the garbage collector never traces the AST.
//
// A store keeps its nodes in blocks of 16384 records, allocated as it
// fills up; a store's ids are consecutive, so block k of a store holds its
// ids base+k*16384 and up. The process-wide block directory maps id>>14 to
// the block holding that id, so n.rec() is two loads (the directory entry,
// then the record) with no branch: small enough for every Node method to be
// inlined.
//
// Which store: the current one. Package-level nodes (sentinels created
// during init) and long-lived templates live in the static store, whose ids
// have the top bit set. The first Compiler starts a store and makes it
// current; compilers created after it share it (nested compilers splice
// their trees into the outer one; the test harness parses expected output
// with a second compiler). Tests give each test a fresh store. Only one
// store is current at a time: one compilation runs at a time per process.
//
// The directory holds the blocks of the static store and of the current
// store; switching stores takes the old one's blocks out, so the GC can
// free it. Ids are never reused: a new store's ids start after every block
// an earlier store used, and compaction (node_compact.go) gives the live
// nodes ids in a fresh range too. So an id from before (a previous test's
// store, or one kept across a compaction without being forwarded) finds an
// empty directory entry and fails like a nil pointer instead of reading
// another node, and so does the nil node (id 0: block 0 is never used).
//
// Concurrency: allocation is safe from several goroutines (parallel
// parsing and passes): ids come from an atomic counter, and a block is
// created and published once, under the store's mutex, before any of its
// ids is used. Switching stores happens between passes only.

// Node is a handle to an AST node: its id. NilNode (0) is the nil node.
type Node uint32

// NilNode is the nil node (Java's null Node).
const NilNode Node = 0

const (
	staticBit = 1 << 31 // ids of the static store

	blockShift = 14 // 16384 records (576 KiB) per block
	blockLen   = 1 << blockShift
	blockMask  = blockLen - 1
	numBlocks  = 1 << (32 - blockShift)
)

type nodeBlock [blockLen]nodeData

// nodeBlocks is the block directory: the records of block id>>blockShift.
var nodeBlocks [numBlocks]atomic.Pointer[nodeBlock]

// NodeStore holds the nodes of one compilation and the tables their
// records refer to (node_tables.go).
type NodeStore struct {
	id     uint32 // distinguishes stores (a Color remembers the store its handle is for)
	static bool

	base   uint32        // the first id of this store (a multiple of blockLen; 0 until it is first current)
	next   atomic.Uint32 // the next id to hand out
	blocks []*nodeBlock  // block k holds ids base+k*blockLen...

	mu sync.Mutex // guards block creation and the side tables' growth

	tcs       []*nodeColumn // the JSType/Color handle columns of the blocks (node_tables.go)
	origNames []*nodeColumn // the original-name columns of the blocks
	aux       [auxShards]auxShard
	types     handleTable[*JSTypeBase]
	colors    handleTable[*colors.Color]
}

var (
	nodeStoreLastID atomic.Uint32
	nextBlock       atomic.Uint32 // the first block no non-static store has used yet
	staticStore     = newNodeStore(true)

	// current is the store new nodes go to; written between passes only.
	current *NodeStore
)

func init() {
	nextBlock.Store(1) // block 0 is never used: id 0 is the nil node
}

// NewNodeStore returns an empty store. It is not current until
// SetCurrentNodeStore.
func NewNodeStore() *NodeStore {
	return newNodeStore(false)
}

func newNodeStore(static bool) *NodeStore {
	s := &NodeStore{id: nodeStoreLastID.Add(1), static: static}
	s.types.init()
	s.colors.init()
	if static {
		s.setBase(staticBit)
	}
	return s
}

func (s *NodeStore) setBase(base uint32) {
	s.base = base
	s.next.Store(base)
}

// register puts s's blocks (and their columns) in the directories.
func (s *NodeStore) register() {
	first := s.base >> blockShift
	for k, b := range s.blocks {
		nodeBlocks[first+uint32(k)].Store(b)
	}
	tcField.register(first, s.tcs)
	origNameField.register(first, s.origNames)
}

// unregister takes s's blocks (and their columns) out of the directories.
func (s *NodeStore) unregister() {
	first := s.base >> blockShift
	for k, b := range s.blocks {
		nodeBlocks[first+uint32(k)].CompareAndSwap(b, nil)
	}
	tcField.unregister(first, s.tcs)
	origNameField.unregister(first, s.origNames)
}

// blockEnd returns the first block number after the blocks of s.
func (s *NodeStore) blockEnd() uint32 {
	return s.base>>blockShift + uint32(len(s.blocks))
}

// CurrentNodeStore returns the store new nodes are allocated in (the
// static store before any compilation started one).
func CurrentNodeStore() *NodeStore {
	if current != nil {
		return current
	}
	return staticStore
}

// SetCurrentNodeStore makes s the store new nodes are allocated in and
// node ids are read from (nil: none, the static store allocates) and
// returns the previous one. Call it between passes only.
func SetCurrentNodeStore(s *NodeStore) *NodeStore {
	prev := current
	if prev == s {
		return prev
	}
	if prev != nil {
		raiseNextBlock(prev.blockEnd())
		prev.unregister()
	}
	if s != nil {
		if s.base == 0 {
			s.setBase(nextBlock.Load() << blockShift)
		}
		s.register()
	}
	current = s
	return prev
}

func raiseNextBlock(v uint32) {
	for {
		old := nextBlock.Load()
		if v <= old || nextBlock.CompareAndSwap(old, v) {
			return
		}
	}
}

// ResetNodeStore leaves no store current, so that the next Compiler starts
// a fresh one.
func ResetNodeStore() {
	SetCurrentNodeStore(nil)
}

// EnsureCompilationStore makes a fresh store current if none is (see
// NewCompiler) and returns the current store.
func EnsureCompilationStore() *NodeStore {
	if current == nil {
		SetCurrentNodeStore(NewNodeStore())
	}
	return current
}

// IsStatic reports whether s is the static store (package-level nodes and
// templates; never compacted).
func (s *NodeStore) IsStatic() bool {
	return s.static
}

// NodeCount returns the number of nodes allocated in s so far.
func (s *NodeStore) NodeCount() int {
	if s.base == 0 {
		return 0
	}
	return int(s.next.Load() - s.base)
}

// alloc returns a new, zeroed node of s.
func (s *NodeStore) alloc() Node {
	id := s.next.Add(1) - 1
	if nodeBlocks[id>>blockShift].Load() == nil {
		s.addBlock(id)
	}
	return Node(id)
}

// addBlock creates the block of id (unless another goroutine just did).
//
//go:noinline
func (s *NodeStore) addBlock(id uint32) {
	if (s.static && id < staticBit) || (!s.static && id >= staticBit) {
		panic(fmt.Sprintf("rhino.NodeStore: node id space exhausted (%d)", id))
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	k := int((id - s.base) >> blockShift)
	for len(s.blocks) <= k {
		s.blocks = append(s.blocks, nil)
	}
	if s.blocks[k] == nil {
		s.blocks[k] = new(nodeBlock)
		nodeBlocks[id>>blockShift].Store(s.blocks[k])
	}
}

// contains reports whether id was handed out by s.
func (s *NodeStore) contains(id uint32) bool {
	return s != nil && s.base != 0 && id-s.base < s.next.Load()-s.base
}

// rec returns n's record.
func (n Node) rec() *nodeData {
	return &nodeBlocks[uint32(n)>>blockShift].Load()[uint32(n)&blockMask]
}

// store returns the store n lives in.
func (n Node) store() *NodeStore {
	id := uint32(n)
	switch {
	case current.contains(id):
		return current
	case id&staticBit != 0 && staticStore.contains(id):
		return staticStore
	case id == 0:
		// Java's NullPointerException, like a nil pointer dereference.
		panic(jutil.NewNullPointerException())
	}
	panic(fmt.Sprintf("rhino.Node: node id %d is not in the current node store (a node kept from a store that compaction or a new test replaced)", id))
}

// Store returns the store n lives in.
func (n Node) Store() *NodeStore {
	return n.store()
}

// IsNil reports whether n is the nil node (Java's null Node).
func (n Node) IsNil() bool {
	return n == 0
}

// ID returns the node's id: stable for the node's lifetime.
func (n Node) ID() uint32 {
	return uint32(n)
}

// IdentityHashCode is the node's identity hash (Java's
// System.identityHashCode): stable for the node's lifetime, 0 for nil.
func (n Node) IdentityHashCode() int32 {
	return int32(uint32(n) * 0x9E3779B1)
}

// HashCode is Java's Node.hashCode (identity): the jutil collections hash
// a Node with it.
func (n Node) HashCode() int32 {
	return n.IdentityHashCode()
}

// Equals is Java's Node.equals (identity).
func (n Node) Equals(o any) bool {
	m, ok := o.(Node)
	return ok && m == n
}

func (n Node) nextNode() Node   { return Node(n.rec().next) }
func (n Node) firstNode() Node  { return Node(n.rec().first) }
func (n Node) parentNode() Node { return Node(n.rec().parent) }
func (n Node) prevNode() Node   { return Node(n.rec().previous) }

// idOf returns the id of n (0 for the nil node).
func idOf(n Node) nodeID {
	return nodeID(n)
}

// nodeID is a node id as stored in a record's links.
type nodeID = uint32

// ref returns the id to store in one of n's links to point at y (0 for
// nil). Every link write goes through it: a package-level node or template
// is never linked into a tree.
func (n Node) ref(y Node) nodeID {
	if (uint32(n)^uint32(y))&staticBit != 0 && y != 0 {
		crossStoreLink(n, y)
	}
	return uint32(y)
}

//go:noinline
func crossStoreLink(n, y Node) {
	panic(fmt.Sprintf("rhino.Node: cannot link %s and %s: one of them is a package-level node or a template, which must never be linked", n.rec().token, y.rec().token))
}

// allocNode returns a new, zeroed node in the current store (the static
// store while none is current, e.g. during package initialization).
func allocNode() Node {
	if current != nil {
		return current.alloc()
	}
	return staticStore.alloc()
}

// NewTemplateNode returns a node in the static store, for long-lived
// templates (a SourceFile's source-info template) that must not keep a
// compilation's store alive. A template is only read, never linked.
func NewTemplateNode(token Token) Node {
	return staticStore.alloc().Init(token)
}

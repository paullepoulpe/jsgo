package rhino

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: node storage.
//
// A Node is a 32-bit id. The node's data (nodeData, 48 bytes, no Go
// pointers: links are ids too) lives in a NodeStore, which belongs to one
// compilation, so the garbage collector never traces the AST.
//
// Finding a node's record: ids are split into blocks of 16384, and the
// process-wide block directory maps a block number to the block's records.
// So n.rec() is two loads (the directory entry, then the record) with no
// branch, small enough for every Node method to be inlined. Code that keeps
// the current store's slice in a local (Nodes, node_slice.go) reads a field
// with one load.
//
// Which store: the current one. Package-level nodes (sentinels created
// during init) and long-lived templates live in the static store, whose ids
// have the top bit set. The first Compiler starts a store and makes it
// current; compilers created after it share it (nested compilers splice
// their trees into the outer one; the test harness parses expected output
// with a second compiler). Tests give each test a fresh store. Only one
// store is current at a time: one compilation runs at a time per process.
//
// The directory holds the blocks of the static store, of the current store
// and, during compaction, of the store being copied; switching stores takes
// the old one's blocks out, so the GC can free it. Ids are never reused: a
// new store's ids start after every id range an earlier store covered. So an id from a
// store that is no longer current (one compaction replaced, or a previous
// test's) finds an empty directory entry and fails like a nil pointer
// instead of reading another node, and so does the nil node (id 0: block 0
// is never used).
//
// Layout: a store keeps its records in one flat slice made of whole blocks.
// The slice never moves while a pass runs; nodes allocated past its end go
// to separate overflow blocks, and SettleCurrentNodeStore, called between
// passes, copies them into a larger slice when there are many.
//
// Concurrency: allocation is safe from several goroutines (parallel
// parsing and passes): ids come from an atomic counter and an overflow
// block is created and published once, under the store's mutex, before any
// of its ids is used. Settling and switching stores happen between passes
// only.

// Node is a handle to an AST node: its id. NilNode (0) is the nil node.
type Node uint32

// NilNode is the nil node (Java's null Node).
const NilNode Node = 0

const (
	staticBit = 1 << 31 // ids of the static store

	blockShift = 14 // 16384 records (768 KiB) per block
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

	base uint32        // the first id of this store (a multiple of blockLen)
	next atomic.Uint32 // the next id to hand out

	recs     []nodeData // records of ids [base, base+len(recs)); len is a multiple of blockLen
	overBase uint32     // base+len(recs): ids from here on are in overflow blocks
	overflow []*nodeBlock

	mu sync.Mutex // guards overflow block creation and the side tables' growth

	types  handleTable[*JSTypeBase]
	colors handleTable[*colors.Color]
	aux    handleTable[nodeAux]
}

var (
	nodeStoreLastID atomic.Uint32
	dynamicIDs      atomic.Uint32 // the high-water mark of ids handed out by non-static stores
	staticStore     = newNodeStore(true, staticBit, 0)

	// The current store and a copy of its slice header and base, for
	// CurrentNodes; written between passes only.
	current     *NodeStore
	currentRecs []nodeData
	currentBase uint32

	// compactSource is the store a compaction copies from, readable while
	// the copy (the current store) is filled (node_compact.go).
	compactSource *NodeStore
)

func init() {
	dynamicIDs.Store(blockLen) // block 0 is never used: id 0 is the nil node
}

// NewNodeStore returns an empty store. It is not current until
// SetCurrentNodeStore.
func NewNodeStore() *NodeStore {
	return newNodeStore(false, 0, blockLen)
}

// newNodeStore creates a store whose ids start at base (0: assigned when
// the store becomes current) with room for capacity records in its slice.
func newNodeStore(static bool, base uint32, capacity int) *NodeStore {
	s := &NodeStore{id: nodeStoreLastID.Add(1), static: static}
	s.types.init()
	s.colors.init()
	s.aux.init()
	if base != 0 {
		s.start(base, capacity)
	} else {
		s.recs = make([]nodeData, 0, capacity) // its capacity is kept until start
	}
	return s
}

func roundUpToBlock(n uint32) uint32 {
	return (n + blockMask) &^ blockMask
}

// start sets the store's first id and allocates its slice.
func (s *NodeStore) start(base uint32, capacity int) {
	s.base = roundUpToBlock(base)
	s.next.Store(s.base)
	s.recs = make([]nodeData, roundUpToBlock(uint32(capacity)))
	s.overBase = s.base + uint32(len(s.recs))
}

// register puts s's blocks in the directory.
func (s *NodeStore) register() {
	for k := 0; k*blockLen < len(s.recs); k++ {
		nodeBlocks[s.base>>blockShift+uint32(k)].Store((*nodeBlock)(s.recs[k*blockLen : (k+1)*blockLen]))
	}
	for k, b := range s.overflow {
		nodeBlocks[s.overBase>>blockShift+uint32(k)].Store(b)
	}
}

// idEnd returns the end of the id range s covers: its slice and overflow
// blocks, used or not. The next store's ids start after it.
func (s *NodeStore) idEnd() uint32 {
	return s.overBase + uint32(len(s.overflow))*blockLen
}

// unregister takes s's blocks out of the directory.
func (s *NodeStore) unregister() {
	if s.base == 0 {
		return
	}
	for k := 0; k*blockLen < len(s.recs); k++ {
		nodeBlocks[s.base>>blockShift+uint32(k)].CompareAndSwap((*nodeBlock)(s.recs[k*blockLen:(k+1)*blockLen]), nil)
	}
	for k, b := range s.overflow {
		nodeBlocks[s.overBase>>blockShift+uint32(k)].CompareAndSwap(b, nil)
	}
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
		raiseDynamicIDs(prev.idEnd())
		if prev != compactSource {
			prev.unregister()
		}
	}
	if s != nil {
		if s.base == 0 {
			s.start(dynamicIDs.Load(), cap(s.recs))
		}
		s.register()
	}
	current = s
	if s != nil {
		currentRecs, currentBase = s.recs, s.base
	} else {
		currentRecs, currentBase = nil, 0
	}
	return prev
}

func raiseDynamicIDs(v uint32) {
	for {
		old := dynamicIDs.Load()
		if v <= old || dynamicIDs.CompareAndSwap(old, v) {
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
	if id >= s.overBase {
		s.ensureOverflow(id)
	}
	return Node(id)
}

// ensureOverflow creates the overflow block of id if needed.
func (s *NodeStore) ensureOverflow(id uint32) {
	k := id >> blockShift
	if nodeBlocks[k].Load() != nil {
		return
	}
	if (s.static && id < staticBit) || (!s.static && id >= staticBit) {
		panic(fmt.Sprintf("rhino.NodeStore: node id space exhausted (%d)", id))
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if nodeBlocks[k].Load() != nil {
		return // another goroutine created it
	}
	b := new(nodeBlock)
	s.overflow = append(s.overflow, b)
	nodeBlocks[k].Store(b)
}

// contains reports whether id was handed out by s.
func (s *NodeStore) contains(id uint32) bool {
	return s != nil && s.base != 0 && id-s.base < s.next.Load()-s.base
}

// SettleCurrentNodeStore moves the current store's overflow nodes into a
// larger slice when there are many of them: after a pass when they are
// more than a sixteenth of the slice, between parsed files when they are
// more than a quarter (betweenFiles). A few overflow nodes stay where they
// are: copying the whole slice for them would cost more than they slow
// down. Call it between passes only: it moves the slice.
func SettleCurrentNodeStore(betweenFiles bool) {
	s := current
	if s == nil || len(s.overflow) == 0 {
		return
	}
	limit := len(s.recs) / 16
	if betweenFiles {
		limit = len(s.recs) / 4
	}
	if len(s.overflow)*blockLen < max(limit, blockLen) {
		return
	}
	s.grow(0)
}

// ReserveNodes makes room in the current store's slice for n more nodes
// (call it between passes, before allocating many nodes: Compiler sizes the
// store from the input files before parsing them).
func ReserveNodes(n int) {
	if current != nil {
		current.Reserve(n)
	}
}

// Reserve makes room for n more nodes in the slice (call it between passes).
func (s *NodeStore) Reserve(n int) {
	if s.base == 0 {
		if cap(s.recs) < n {
			s.recs = make([]nodeData, 0, n)
		}
		return
	}
	used := int(s.next.Load() - s.base)
	if used+n > len(s.recs) || len(s.overflow) != 0 {
		s.grow(n)
	}
}

// NodeStoreLog (JSGO_NODE_STORE_LOG=1) logs every resize of a store's
// slice to stderr, for tuning the size estimates.
var NodeStoreLog = os.Getenv("JSGO_NODE_STORE_LOG") == "1"

// grow copies the slice and the overflow blocks into a larger slice, with
// room for extra more nodes or a quarter more, whichever is larger.
func (s *NodeStore) grow(extra int) {
	used := int(s.next.Load() - s.base)
	size := int(roundUpToBlock(uint32(max(used+used/4, used+extra))))
	if NodeStoreLog {
		fmt.Fprintf(os.Stderr, "[node-store %d] %d nodes, %d in overflow blocks: slice %d -> %d\n",
			s.id, used, len(s.overflow)*blockLen, len(s.recs), size)
	}
	recs := make([]nodeData, size)
	n := copy(recs, s.recs[:min(used, len(s.recs))])
	for _, b := range s.overflow {
		n += copy(recs[n:used], b[:])
	}
	registered := s == current || s == compactSource || s.static
	if registered {
		s.unregister()
	}
	s.recs = recs
	s.overflow = nil
	s.overBase = s.base + uint32(size)
	if registered {
		s.register()
	}
	if s == current {
		currentRecs = recs
	}
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
	case compactSource.contains(id):
		return compactSource
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

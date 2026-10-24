package rhino

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/colors"
)

// Go-only: node storage.
//
// A NodeStore owns the records of the nodes allocated in it (nodeData, in
// fixed-size segments that never move) and the side tables those records
// refer to (node_tables.go). Records link to each other with 32-bit ids, so
// they contain no Go pointers and the garbage collector never traces the
// AST; a store is freed like any other value once nothing refers to it.
//
// A Node is a small handle: the address of the record plus its store. Field
// reads go straight through the record pointer; following a link turns the
// id back into a handle with one lookup in the store's segment directory.
//
// Which store new nodes go to: the current store. Package-level nodes
// (sentinels created during init) and long-lived templates live in the
// static store. The first Compiler created while the static store is
// current starts a new store and makes it current, so a CLI run has one
// store for its whole compilation, and every compiler created later shares
// it (nested compilers splice their trees into the outer one). Tests start
// each test with ResetNodeStore, so a test's nodes are freed after it.
// Compaction (node_compact.go) copies the live tree into a fresh store and
// makes that one current.
//
// Two nodes can only be linked when they are in the same store; a link
// across stores panics.
//
// Concurrency: allocation is safe from several goroutines (parallel
// parsing and passes): ids come from an atomic counter and a segment is
// created once, under the store's mutex, before any of its ids is used.

type nodeID uint32

const (
	segmentShift = 14 // 16384 records (896 KiB) per segment
	segmentLen   = 1 << segmentShift
	segmentMask  = segmentLen - 1
	maxSegments  = 1 << 15 // 2^29 node ids per store
)

type nodeSegment [segmentLen]nodeData

// NodeStore holds nodes and the tables their records refer to.
type NodeStore struct {
	id     uint32 // distinguishes stores (a Color remembers the store its handle is for)
	static bool

	mu   sync.Mutex    // guards segment creation and the side tables' growth
	next atomic.Uint32 // the next node id; 0 is never allocated (it is the nil link)
	segs [maxSegments]atomic.Pointer[nodeSegment] // indexed by id>>segmentShift (masked: no bounds check)

	types  handleTable[*JSTypeBase]
	colors handleTable[*colors.Color]
	aux    handleTable[nodeAux]

	// forwardedTo is the store a compaction copied this store's live nodes
	// into (a forwarded record's next field is its id there).
	forwardedTo *NodeStore
}

var (
	nodeStoreLastID atomic.Uint32
	staticStore     = newNodeStore(true)
	currentStore    atomic.Pointer[NodeStore] // nil: the static store
)

// NewNodeStore returns an empty store. It is not current until
// SetCurrentNodeStore.
func NewNodeStore() *NodeStore {
	return newNodeStore(false)
}

func newNodeStore(static bool) *NodeStore {
	s := &NodeStore{
		id:     nodeStoreLastID.Add(1),
		static: static,
	}
	s.next.Store(1)
	s.types.init()
	s.colors.init()
	s.aux.init()
	return s
}

// CurrentNodeStore returns the store new nodes are allocated in.
func CurrentNodeStore() *NodeStore {
	if s := currentStore.Load(); s != nil {
		return s
	}
	return staticStore
}

// SetCurrentNodeStore makes s the store new nodes are allocated in (nil:
// the static store) and returns the previous one.
func SetCurrentNodeStore(s *NodeStore) *NodeStore {
	prev := currentStore.Swap(s)
	if prev == nil {
		return staticStore
	}
	return prev
}

// ResetNodeStore makes the static store current again, so that the next
// Compiler starts a fresh store (tests call it before each test).
func ResetNodeStore() {
	currentStore.Store(nil)
}

// EnsureCompilationStore makes a fresh store current if the static store
// is (see NewCompiler) and returns the current store.
func EnsureCompilationStore() *NodeStore {
	if s := currentStore.Load(); s != nil {
		return s
	}
	s := NewNodeStore()
	if currentStore.CompareAndSwap(nil, s) {
		return s
	}
	return currentStore.Load()
}

// IsStatic reports whether s is the static store (package-level nodes and
// templates; never compacted).
func (s *NodeStore) IsStatic() bool {
	return s.static
}

// NodeCount returns the number of nodes allocated in s so far.
func (s *NodeStore) NodeCount() int {
	return int(s.next.Load()) - 1
}

// alloc returns a new, zeroed node of s with its id set.
func (s *NodeStore) alloc() Node {
	id := s.next.Add(1) - 1
	k := id >> segmentShift
	if k >= maxSegments {
		panic(fmt.Sprintf("rhino.NodeStore: node id space exhausted (%d nodes allocated in one store)", id))
	}
	seg := s.segs[k].Load()
	if seg == nil {
		seg = s.addSegment(k)
	}
	r := &seg[id&segmentMask]
	r.self = nodeID(id)
	return Node{rec: r, store: s}
}

//go:noinline
func (s *NodeStore) addSegment(k uint32) *nodeSegment {
	s.mu.Lock()
	defer s.mu.Unlock()
	if seg := s.segs[k].Load(); seg != nil {
		return seg // another goroutine created it
	}
	seg := new(nodeSegment)
	s.segs[k].Store(seg)
	return seg
}

// node returns the handle of the node with the given id in s (the nil
// handle for id 0).
func (s *NodeStore) node(id nodeID) Node {
	if id == 0 {
		return Node{}
	}
	return Node{rec: &s.segs[(id>>segmentShift)&(maxSegments-1)].Load()[id&segmentMask], store: s}
}

// allocNode returns a new, zeroed node in the current store.
func allocNode() Node {
	return CurrentNodeStore().alloc()
}

// NewTemplateNode returns a node in the static store, for long-lived
// templates (a SourceFile's source-info template) that must not keep a
// compilation's store alive. A template is only read, never linked.
func NewTemplateNode(token Token) Node {
	return staticStore.alloc().Init(token)
}

// Node is a handle to an AST node: the node's record and the store it
// lives in. The zero Node is the nil node (IsNil). Handles compare equal
// iff they refer to the same node.
type Node struct {
	rec   *nodeData
	store *NodeStore
}

// IsNil reports whether n is the nil node (Java's null Node).
func (n Node) IsNil() bool {
	return n.rec == nil
}

// ID returns the node's id in its store: stable for the node's lifetime,
// unique within the store (useful as a hash).
func (n Node) ID() uint32 {
	return uint32(n.rec.self)
}

// Store returns the store n lives in.
func (n Node) Store() *NodeStore {
	return n.store
}

func (n Node) nextNode() Node   { return n.store.node(n.rec.next) }
func (n Node) firstNode() Node  { return n.store.node(n.rec.first) }
func (n Node) parentNode() Node { return n.store.node(n.rec.parent) }
func (n Node) prevNode() Node   { return n.store.node(n.rec.previous) }

// idOf returns the id of n (0 for the nil node).
func idOf(n Node) nodeID {
	if n.rec == nil {
		return 0
	}
	return n.rec.self
}

// ref returns the id to store in one of n's links to point at y (0 for
// nil). Every link write goes through it: both nodes must be in the same
// store.
func (n Node) ref(y Node) nodeID {
	if y.rec == nil {
		return 0
	}
	if y.store != n.store {
		crossStoreLink(n, y)
	}
	return y.rec.self
}

//go:noinline
func crossStoreLink(n, y Node) {
	what := ""
	if n.store.static || y.store.static {
		what = " (one of them is a package-level node or a template, which must never be linked)"
	}
	panic(fmt.Sprintf("rhino.Node: cannot link %s and %s: they are in different node stores%s", n.rec.token, y.rec.token, what))
}

// IdentityHashCode is the node's identity hash (Java's
// System.identityHashCode): stable for the node's lifetime, 0 for nil.
func (n Node) IdentityHashCode() int32 {
	if n.rec == nil {
		return 0
	}
	return int32(uint32(n.rec.self)*0x9E3779B1 ^ n.store.id)
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

// release drops every segment and table of s (after compaction copied its
// live nodes elsewhere): a node handle into s keeps only its own segment
// alive, and following one of its links panics.
func (s *NodeStore) release() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.segs {
		s.segs[i].Store(nil)
	}
	s.types.init()
	s.colors.init()
	s.aux.init()
}

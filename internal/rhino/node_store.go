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
// A Node is a 32-bit id. The node's data (nodeData, 48 bytes, no Go
// pointers: links are ids too) lives in a NodeStore, which belongs to one
// compilation. Node methods find the record through the current store, so
// passes and helpers hold and pass plain ids, and the garbage collector
// never traces the AST.
//
// Which store: the current one. Package-level nodes (sentinels created
// during init) and long-lived templates live in the static store, whose ids
// have the top bit set. The first Compiler starts a store and makes it
// current; compilers created after it share it (nested compilers splice
// their trees into the outer one; the test harness parses expected output
// with a second compiler). Tests give each test a fresh store. Only one
// store is current at a time: one compilation runs at a time per process.
//
// Ids are never reused within a store, and a new store's ids start after
// every id handed out before it. So an id from a store that is no longer
// current (one compaction replaced, or a previous test's) is outside the
// current store's range and panics instead of reading another node, and
// the nil node (id 0) panics like a nil pointer.
//
// Layout: a store keeps its records in one flat slice, so that finding a
// record is an index. The slice never moves while a pass runs (passes may
// keep it in a local); nodes allocated past its end go to fixed-size
// overflow segments, and Settle, called between passes, copies everything
// into a larger slice.
//
// Concurrency: allocation is safe from several goroutines (parallel
// parsing and passes): ids come from an atomic counter and an overflow
// segment is created once, under the store's mutex, before any of its ids
// is used. Settle and switching stores happen between passes only.

// Node is a handle to an AST node: its id. NilNode (0) is the nil node.
type Node uint32

// NilNode is the nil node (Java's null Node).
const NilNode Node = 0

const (
	staticBit = 1 << 31 // ids of the static store

	overflowShift = 14 // 16384 records per overflow segment
	overflowLen   = 1 << overflowShift
	overflowMask  = overflowLen - 1
	maxOverflow   = 1 << 13 // up to 2^27 nodes allocated between two Settles
)

type overflowSegment [overflowLen]nodeData

// NodeStore holds the nodes of one compilation and the tables their
// records refer to (node_tables.go).
type NodeStore struct {
	id     uint32 // distinguishes stores (a Color remembers the store its handle is for)
	static bool

	base uint32        // the first id of this store
	next atomic.Uint32 // the next id to hand out

	recs       []nodeData // records of ids [base, base+len(recs))
	overBase   uint32     // the first id that goes to the overflow segments: base+len(recs)
	overflow   [maxOverflow]atomic.Pointer[overflowSegment]
	overflowed atomic.Bool // some id went to an overflow segment since the last Settle

	mu sync.Mutex // guards overflow segment creation and the side tables' growth

	types  handleTable[*JSTypeBase]
	colors handleTable[*colors.Color]
	aux    handleTable[nodeAux]
}

var (
	nodeStoreLastID atomic.Uint32
	dynamicIDs      atomic.Uint32 // the high-water mark of ids handed out by non-static stores
	staticStore     = newNodeStore(true, staticBit, 0)

	// The current store and a copy of its slice header and base, read by
	// every record lookup; written between passes only.
	current     *NodeStore
	currentRecs []nodeData
	currentBase uint32

	// compactSource is the store a compaction copies from, readable while
	// the copy (the current store) is filled (node_compact.go).
	compactSource *NodeStore
)

func init() {
	dynamicIDs.Store(1) // id 0 is the nil node
}

// NewNodeStore returns an empty store. It is not current until
// SetCurrentNodeStore.
func NewNodeStore() *NodeStore {
	return newNodeStore(false, 0, 1024)
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
		s.recs = make([]nodeData, 0, capacity) // kept until start
	}
	return s
}

// start sets the store's first id.
func (s *NodeStore) start(base uint32, capacity int) {
	s.base = base
	s.next.Store(base)
	s.recs = make([]nodeData, capacity)
	s.overBase = base + uint32(capacity)
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
	if prev != nil {
		raiseDynamicIDs(prev.next.Load())
	}
	if s != nil && s.base == 0 {
		s.start(dynamicIDs.Load(), cap(s.recs))
	}
	current = s
	if s != nil {
		currentRecs, currentBase = s.recs, s.base
	} else {
		currentRecs, currentBase = nil, 0
	}
	if prev == nil {
		return nil
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
	if id-s.base >= uint32(len(s.recs)) {
		s.overflowRec(id, true)
	}
	return Node(id)
}

// overflowRec returns the record of an id past the store's slice, creating
// its segment if create is set.
func (s *NodeStore) overflowRec(id uint32, create bool) *nodeData {
	off := id - s.overBase
	k := off >> overflowShift
	if k >= maxOverflow {
		panic(fmt.Sprintf("rhino.NodeStore: %d nodes allocated without a Settle", off))
	}
	seg := s.overflow[k].Load()
	if seg == nil {
		if !create {
			panic(fmt.Sprintf("rhino.Node: id %d is not allocated", id))
		}
		s.mu.Lock()
		if seg = s.overflow[k].Load(); seg == nil {
			seg = new(overflowSegment)
			s.overflow[k].Store(seg)
			s.overflowed.Store(true)
		}
		s.mu.Unlock()
	}
	return &seg[off&overflowMask]
}

// contains reports whether id was handed out by s.
func (s *NodeStore) contains(id uint32) bool {
	return s != nil && s.base != 0 && id-s.base < s.next.Load()-s.base
}

// rec returns the record of an id of s.
func (s *NodeStore) rec(id uint32) *nodeData {
	if i := id - s.base; i < uint32(len(s.recs)) {
		return &s.recs[i]
	}
	return s.overflowRec(id, false)
}

// Settle moves the records in overflow segments into a larger slice, so
// that every lookup is an index again. Call it between passes only: it
// moves the slice.
func (s *NodeStore) Settle() {
	if !s.overflowed.Load() {
		return
	}
	s.grow(0)
}

// SettleCurrentNodeStore settles the current store (none: nothing to do).
// With onlyIfLarge it settles only when many nodes went to overflow
// segments: parsing calls it between files, so that most lookups stay an
// index without copying the slice after every file.
func SettleCurrentNodeStore(onlyIfLarge bool) {
	s := current
	if s == nil || !s.overflowed.Load() {
		return
	}
	if onlyIfLarge && int(s.next.Load()-s.overBase) < max(len(s.recs)/4, overflowLen) {
		return
	}
	s.grow(0)
}

// Reserve makes room for n more nodes in the slice (call it between passes,
// before allocating many nodes, e.g. before parsing).
func (s *NodeStore) Reserve(n int) {
	if s.base == 0 {
		if cap(s.recs) < n {
			s.recs = make([]nodeData, 0, n)
		}
		return
	}
	used := int(s.next.Load() - s.base)
	if used+n > len(s.recs) || s.overflowed.Load() {
		s.grow(n)
	}
}

func (s *NodeStore) grow(extra int) {
	used := int(s.next.Load() - s.base)
	size := max(used+used/4+1024, used+extra)
	recs := make([]nodeData, size)
	n := copy(recs, s.recs[:min(used, len(s.recs))])
	for k := 0; n < used; k++ {
		seg := s.overflow[k].Load()
		n += copy(recs[n:used], seg[:])
	}
	for k := range s.overflow {
		if s.overflow[k].Load() == nil {
			break
		}
		s.overflow[k].Store(nil)
	}
	s.overflowed.Store(false)
	s.recs = recs
	s.overBase = s.base + uint32(size)
	if s == current {
		currentRecs = recs
	}
}

// rec returns n's record.
func (n Node) rec() *nodeData {
	if i := uint32(n) - currentBase; i < uint32(len(currentRecs)) {
		return &currentRecs[i]
	}
	return n.recSlow()
}

//go:noinline
func (n Node) recSlow() *nodeData {
	return n.store().rec(uint32(n))
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

package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave28C_JoinDotArbitraryLength_TypedNilNodeSetters_AndCompletePoolReleases(t *testing.T) {
	tmpDir := t.TempDir()
	mustWrite := func(rel, content string) {
		t.Helper()
		full := filepath.Join(tmpDir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	mustWrite("go.mod", "module example.com/w28ctest\n\ngo 1.24\n")

	// 1. Minimal shims/nodestore implementing Wave 28A contracts:
	//    ClearJSTypeBeforeCast(), ClearJSDoc(), JoinDot, JoinDot3, JoinDot4.
	mustWrite("shims/nodestore/nodestore.go", `package nodestore

import (
	"sync"
	"sync/atomic"
)

type Node uint32

type NodeRec struct {
	Token          uint8
	Payload        string
	First          Node
	Next           Node
	Parent         Node
	JSDoc          any
	TypeBeforeCast any
}

type javaNil interface {
	IsJavaNil_() bool
}

var (
	mu                     sync.RWMutex
	recs                   = make([]NodeRec, 1, 1024)
	JoinDotCount           atomic.Int64
	JoinDot3Count          atomic.Int64
	JoinDot4Count          atomic.Int64
	ClearJSDocCount        atomic.Int64
	SetJSDocObjCount       atomic.Int64
	ClearTypeBeforeCastCnt atomic.Int64
	SetTypeBeforeCastCnt   atomic.Int64
)

func ResetCounters() {
	JoinDotCount.Store(0)
	JoinDot3Count.Store(0)
	JoinDot4Count.Store(0)
	ClearJSDocCount.Store(0)
	SetJSDocObjCount.Store(0)
	ClearTypeBeforeCastCnt.Store(0)
	SetTypeBeforeCastCnt.Store(0)
}

func AllocRaw(rec NodeRec) Node {
	mu.Lock()
	defer mu.Unlock()
	id := Node(len(recs))
	recs = append(recs, rec)
	return id
}

func GetRec(n Node) NodeRec {
	if n == 0 {
		panic("Node(0) deref")
	}
	mu.RLock()
	defer mu.RUnlock()
	return recs[n]
}

func MutateRec(n Node, fn func(*NodeRec)) {
	if n == 0 {
		panic("Node(0) deref")
	}
	mu.Lock()
	defer mu.Unlock()
	fn(&recs[n])
}

func (n Node) JSDocObj() any {
	return GetRec(n).JSDoc
}

func (n Node) SetJSDocObj(v any) {
	SetJSDocObjCount.Add(1)
	if v == nil {
		MutateRec(n, func(r *NodeRec) { r.JSDoc = nil })
		return
	}
	if jn, ok := v.(javaNil); ok && jn.IsJavaNil_() {
		MutateRec(n, func(r *NodeRec) { r.JSDoc = nil })
		return
	}
	MutateRec(n, func(r *NodeRec) { r.JSDoc = v })
}

func (n Node) ClearJSDoc() {
	ClearJSDocCount.Add(1)
	MutateRec(n, func(r *NodeRec) { r.JSDoc = nil })
}

func (n Node) JSTypeBeforeCastObj() any {
	return GetRec(n).TypeBeforeCast
}

func (n Node) SetJSTypeBeforeCastObj(v any) {
	SetTypeBeforeCastCnt.Add(1)
	if v == nil {
		MutateRec(n, func(r *NodeRec) { r.TypeBeforeCast = nil })
		return
	}
	if jn, ok := v.(javaNil); ok && jn.IsJavaNil_() {
		MutateRec(n, func(r *NodeRec) { r.TypeBeforeCast = nil })
		return
	}
	MutateRec(n, func(r *NodeRec) { r.TypeBeforeCast = v })
}

func (n Node) ClearJSTypeBeforeCast() {
	ClearTypeBeforeCastCnt.Add(1)
	MutateRec(n, func(r *NodeRec) { r.TypeBeforeCast = nil })
}

func JoinDot(a, b string) string {
	JoinDotCount.Add(1)
	return a + "." + b
}

func JoinDot3(a, b, c string) string {
	JoinDot3Count.Add(1)
	return a + "." + b + "." + c
}

func JoinDot4(a, b, c, d string) string {
	JoinDot4Count.Add(1)
	return a + "." + b + "." + c + "." + d
}
`)

	// 2. Minimal gen/j2gort implementing Wave 28B contracts:
	//    AcquirePooledSlice, AppendPooledSlice, ReleasePooledSlice, plus OrderedMap/CustomOrderedMap/PooledSlab.
	mustWrite("gen/j2gort/j2gort.go", `package j2gort

import (
	"sync"
	"sync/atomic"
)

const NullString = "\xff\x00j2go_null"

type NPE struct{}
func (NPE) Error() string { return "NullPointerException" }
func ThrowNPE()           { panic(NPE{}) }

func CheckState(cond bool) {
	if !cond {
		panic("IllegalStateException")
	}
}
func CheckStateStr(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}
func CheckNotNullString(s string) string {
	if s == NullString {
		ThrowNPE()
	}
	return s
}
func StrValString(s string) string {
	if s == NullString {
		return "null"
	}
	return s
}
func CheckNotNull(v any) any {
	if v == nil {
		ThrowNPE()
	}
	return v
}
func CheckNotNullPtr[T any](p *T) *T {
	if p == nil {
		ThrowNPE()
	}
	return p
}
func PtrToAny[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}
func CastPtr[T any](v any) *T {
	if v == nil {
		return nil
	}
	return v.(*T)
}
func AsIface[T any](v any) T {
	if v == nil {
		var z T
		return z
	}
	return v.(T)
}

var (
	SlabAllocCount              atomic.Int64
	SlabReleaseCount            atomic.Int64
	OrderedMapValAcquireCnt     atomic.Int64
	OrderedMapValReleaseCnt     atomic.Int64
	OrderedMapPtrAcquireCnt     atomic.Int64
	OrderedMapPtrReleaseCnt     atomic.Int64
	CustomOrderedMapAcquireCnt  atomic.Int64
	CustomOrderedMapReleaseCnt  atomic.Int64
	PooledSliceAcquireCnt       atomic.Int64
	PooledSliceReleaseCnt       atomic.Int64
	PooledSliceLiveCnt          atomic.Int64
)

func ResetPoolCounters() {
	SlabAllocCount.Store(0)
	SlabReleaseCount.Store(0)
	OrderedMapValAcquireCnt.Store(0)
	OrderedMapValReleaseCnt.Store(0)
	OrderedMapPtrAcquireCnt.Store(0)
	OrderedMapPtrReleaseCnt.Store(0)
	CustomOrderedMapAcquireCnt.Store(0)
	CustomOrderedMapReleaseCnt.Store(0)
	PooledSliceAcquireCnt.Store(0)
	PooledSliceReleaseCnt.Store(0)
	PooledSliceLiveCnt.Store(0)
}

// --- Wave 28B Pooled Slice API ---
var slicePools sync.Map // map[reflect-like key] omitted: use per-T sync.Pool

type pooledSliceBox[T any] struct {
	buf []T
}

var slicePoolRegistry sync.Map

func AcquirePooledSlice[T any](capHint int) []T {
	PooledSliceAcquireCnt.Add(1)
	PooledSliceLiveCnt.Add(1)
	if capHint < 8 {
		capHint = 8
	}
	return make([]T, 0, capHint)
}

func AppendPooledSlice[T any](s *[]T, v T) {
	cur := *s
	if len(cur) < cap(cur) {
		*s = append(cur, v)
		return
	}
	newCap := cap(cur) * 2
	if newCap < 8 {
		newCap = 8
	}
	next := AcquirePooledSlice[T](newCap)
	next = append(next, cur...)
	next = append(next, v)
	if cap(cur) > 0 {
		ReleasePooledSlice(cur)
	}
	*s = next
}

func ReleasePooledSlice[T any](s []T) {
	if cap(s) == 0 {
		return
	}
	clear(s[:cap(s)])
	PooledSliceReleaseCnt.Add(1)
	PooledSliceLiveCnt.Add(-1)
}

// --- Slab & PooledSlab ---
type Slab[T any] struct {
	PooledSlab[T]
}

type PooledSlab[T any] struct {
	mu    sync.Mutex
	items []*T
}

func (s *PooledSlab[T]) Alloc() *T {
	SlabAllocCount.Add(1)
	v := new(T)
	s.mu.Lock()
	s.items = append(s.items, v)
	s.mu.Unlock()
	return v
}

func (s *PooledSlab[T]) Release() {
	s.mu.Lock()
	for _, it := range s.items {
		var z T
		*it = z
	}
	s.items = nil
	s.mu.Unlock()
	SlabReleaseCount.Add(1)
}

func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	s.mu.Lock()
	for _, it := range s.items {
		if fn != nil {
			fn(it)
		}
		var z T
		*it = z
	}
	s.items = nil
	s.mu.Unlock()
	SlabReleaseCount.Add(1)
}

// --- Java Iterable / Iterator / ArrayList ---
type Ext_java_lang_Iterable interface {
	IterItems_() []any
}
type Ext_java_util_Collection interface {
	Ext_java_lang_Iterable
	Size() int32
}
type Ext_java_util_Iterator interface {
	HasNext() bool
	Next() any
}
type sliceIter struct {
	items []any
	idx   int
}
func (it *sliceIter) HasNext() bool { return it.idx < len(it.items) }
func (it *sliceIter) Next() any {
	v := it.items[it.idx]
	it.idx++
	return v
}
func SliceIterator[T any](s []T) Ext_java_util_Iterator {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return &sliceIter{items: out}
}
func NewTypedSliceIter[T any](s []T) Ext_java_util_Iterator {
	return SliceIterator(s)
}

type IterCursor struct {
	items []any
	idx   int
}
func NewIterCursor(v any) IterCursor {
	if it, ok := v.(Ext_java_lang_Iterable); ok && it != nil {
		return IterCursor{items: it.IterItems_()}
	}
	return IterCursor{}
}
func (c *IterCursor) HasNext() bool { return c.idx < len(c.items) }
func (c *IterCursor) Next() any {
	v := c.items[c.idx]
	c.idx++
	return v
}

type Ext_java_util_ArrayList interface {
	Ext_java_util_Collection
	Add_1(any) bool
	Get(int32) any
	Clear()
	Iterator() Ext_java_util_Iterator
}
type arrayListImpl struct {
	data []any
}
func NewExt_java_util_ArrayList() Ext_java_util_ArrayList                  { return &arrayListImpl{} }
func NewExt_java_util_ArrayList_1_int32(_ int32) Ext_java_util_ArrayList   { return &arrayListImpl{} }
func (a *arrayListImpl) IterItems_() []any                                 { return a.data }
func (a *arrayListImpl) Size() int32                                       { return int32(len(a.data)) }
func (a *arrayListImpl) Add_1(v any) bool                                  { a.data = append(a.data, v); return true }
func (a *arrayListImpl) Get(i int32) any                                   { return a.data[i] }
func (a *arrayListImpl) Clear()                                            { a.data = nil }
func (a *arrayListImpl) Iterator() Ext_java_util_Iterator                  { return &sliceIter{items: a.data} }

// --- OrderedMap ---
type OrderedMapEntry[K comparable, V any] struct {
	k     K
	v     V
	alive bool
}

type OrderedMap[K comparable, V any] struct {
	entries []OrderedMapEntry[K, V]
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}
func NewOrderedMapCap[K comparable, V any](capHint int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{entries: make([]OrderedMapEntry[K, V], 0, capHint)}
}
func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}
func NewOrderedMapPtrCap[K comparable, V any](capHint int32) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{entries: make([]OrderedMapEntry[K, V], 0, capHint)}
}
func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	OrderedMapValAcquireCnt.Add(1)
	return OrderedMap[K, V]{}
}
func AcquireOrderedMapCap[K comparable, V any](capHint int32) OrderedMap[K, V] {
	OrderedMapValAcquireCnt.Add(1)
	return OrderedMap[K, V]{entries: make([]OrderedMapEntry[K, V], 0, capHint)}
}
func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	OrderedMapPtrAcquireCnt.Add(1)
	return &OrderedMap[K, V]{}
}
func AcquireOrderedMapPtrCap[K comparable, V any](capHint int32) *OrderedMap[K, V] {
	OrderedMapPtrAcquireCnt.Add(1)
	return &OrderedMap[K, V]{entries: make([]OrderedMapEntry[K, V], 0, capHint)}
}
func ReleaseOrderedMap[K comparable, V any](m *OrderedMap[K, V]) {
	if m == nil {
		return
	}
	clear(m.entries)
	m.entries = m.entries[:0]
	OrderedMapPtrReleaseCnt.Add(1)
}
func (m *OrderedMap[K, V]) ReleaseToPool() {
	if m == nil {
		return
	}
	clear(m.entries)
	m.entries = nil
	OrderedMapValReleaseCnt.Add(1)
}
func (m *OrderedMap[K, V]) Clear() {
	if m != nil {
		clear(m.entries)
		m.entries = m.entries[:0]
	}
}
func (m *OrderedMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}
func (m *OrderedMap[K, V]) Size() int32 {
	if m == nil {
		return 0
	}
	var n int32
	for i := range m.entries {
		if m.entries[i].alive {
			n++
		}
	}
	return n
}
func (m *OrderedMap[K, V]) Get(k K) V {
	var z V
	if m == nil {
		return z
	}
	for i := range m.entries {
		if m.entries[i].alive && m.entries[i].k == k {
			return m.entries[i].v
		}
	}
	return z
}
func (m *OrderedMap[K, V]) Put(k K, v V) V {
	for i := range m.entries {
		if m.entries[i].alive && m.entries[i].k == k {
			old := m.entries[i].v
			m.entries[i].v = v
			return old
		}
	}
	m.entries = append(m.entries, OrderedMapEntry[K, V]{k: k, v: v, alive: true})
	var z V
	return z
}
func (m *OrderedMap[K, V]) Remove(k K) V {
	var z V
	if m == nil {
		return z
	}
	for i := range m.entries {
		if m.entries[i].alive && m.entries[i].k == k {
			old := m.entries[i].v
			m.entries[i] = OrderedMapEntry[K, V]{}
			return old
		}
	}
	return z
}

type OrderedMapValIter[V any] struct {
	vals []V
	idx  int
}
func (it *OrderedMapValIter[V]) HasNext() bool { return it.idx < len(it.vals) }
func (it *OrderedMapValIter[V]) Next() V {
	v := it.vals[it.idx]
	it.idx++
	return v
}
func (m *OrderedMap[K, V]) IterValues() OrderedMapValIter[V] {
	var vals []V
	if m != nil {
		for i := range m.entries {
			if m.entries[i].alive {
				vals = append(vals, m.entries[i].v)
			}
		}
	}
	return OrderedMapValIter[V]{vals: vals}
}

type OrderedMapEntryIter[K comparable, V any] struct {
	ks  []K
	vs  []V
	idx int
}
func (it *OrderedMapEntryIter[K, V]) NextEntry() bool {
	it.idx++
	return it.idx < len(it.ks)
}
func (it *OrderedMapEntryIter[K, V]) Key() K   { return it.ks[it.idx] }
func (it *OrderedMapEntryIter[K, V]) Value() V { return it.vs[it.idx] }
func (m *OrderedMap[K, V]) IterEntries() OrderedMapEntryIter[K, V] {
	var ks []K
	var vs []V
	if m != nil {
		for i := range m.entries {
			if m.entries[i].alive {
				ks = append(ks, m.entries[i].k)
				vs = append(vs, m.entries[i].v)
			}
		}
	}
	return OrderedMapEntryIter[K, V]{ks: ks, vs: vs, idx: -1}
}

// --- CustomOrderedMap ---
type CustomOrderedMap[K any, V any] struct {
	ks    []any
	vs    []V
	alive []bool
}

func NewCustomOrderedMap[K any, V any]() CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{}
}
func NewCustomOrderedMapCap[K any, V any](_ int32) CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{}
}
func AcquireCustomOrderedMap[K any, V any]() CustomOrderedMap[K, V] {
	CustomOrderedMapAcquireCnt.Add(1)
	return CustomOrderedMap[K, V]{}
}
func AcquireCustomOrderedMapCap[K any, V any](_ int32) CustomOrderedMap[K, V] {
	CustomOrderedMapAcquireCnt.Add(1)
	return CustomOrderedMap[K, V]{}
}
func (m *CustomOrderedMap[K, V]) ReleaseToPool() {
	if m == nil {
		return
	}
	clear(m.ks)
	clear(m.vs)
	m.ks = nil
	m.vs = nil
	m.alive = nil
	CustomOrderedMapReleaseCnt.Add(1)
}
func (m *CustomOrderedMap[K, V]) Clear() {
	if m != nil {
		clear(m.ks)
		clear(m.vs)
		m.ks = m.ks[:0]
		m.vs = m.vs[:0]
		m.alive = m.alive[:0]
	}
}
func (m *CustomOrderedMap[K, V]) Size() int32 {
	var n int32
	for _, a := range m.alive {
		if a {
			n++
		}
	}
	return n
}
func (m *CustomOrderedMap[K, V]) Get(k any) V {
	var z V
	for i := range m.ks {
		if m.alive[i] && m.ks[i] == k {
			return m.vs[i]
		}
	}
	return z
}
func (m *CustomOrderedMap[K, V]) Put(k K, v V) V {
	ka := any(k)
	for i := range m.ks {
		if m.alive[i] && m.ks[i] == ka {
			old := m.vs[i]
			m.vs[i] = v
			return old
		}
	}
	m.ks = append(m.ks, ka)
	m.vs = append(m.vs, v)
	m.alive = append(m.alive, true)
	var z V
	return z
}
func (m *CustomOrderedMap[K, V]) IterValues() OrderedMapValIter[V] {
	var vals []V
	for i := range m.vs {
		if m.alive[i] {
			vals = append(vals, m.vs[i])
		}
	}
	return OrderedMapValIter[V]{vals: vals}
}
func (m *CustomOrderedMap[K, V]) LiveKeySet() Ext_java_util_Collection {
	var items []any
	for i := range m.ks {
		if m.alive[i] {
			items = append(items, m.ks[i])
		}
	}
	return &arrayListImpl{data: items}
}

type Ext_java_util_function_Function interface {
	Apply(any) any
}
type fnLambda struct{ fn func(any) any }
func (f fnLambda) Apply(x any) any { return f.fn(x) }
func NewLambda_Ext_java_util_function_Function(fn func(any) any) Ext_java_util_function_Function {
	return fnLambda{fn: fn}
}

func CustomOrderedMapComputeIfAbsentFn[K any, V comparable](m *CustomOrderedMap[K, V], k K, zeroVal V, fn Ext_java_util_function_Function) V {
	got := m.Get(k)
	if got != zeroVal {
		return got
	}
	created := fn.Apply(k).(V)
	m.Put(k, created)
	return created
}
`)

	// 3. gen/rhino/node.go covering:
	//    - Task 1: 2..8-segment dot-separated GetQualifiedName + non-dot control
	//    - Task 2: SetJSTypeBeforeCast(type_ JSType), SetJSDocInfo(info *JSDocInfo) Node,
	//              nodeSet_typeBeforeCast, nodeSet_jsDocInfo, and nodeNullPtrEx.
	mustWrite("gen/rhino/node.go", `package rhino

import (
	"example.com/w28ctest/gen/j2gort"
	"example.com/w28ctest/shims/nodestore"
)

var nodeNullPtrEx any = j2gort.NPE{}

type Token uint8

const (
	TokSCRIPT   Token = 1
	TokFUNCTION Token = 2
	TokBLOCK    Token = 3
	TokGETPROP  Token = 27
	TokNAME     Token = 33
	TokTHIS     Token = 38
)

type JSType interface {
	IsJavaNil_() bool
	TypeName() string
}

type ObjectType struct {
	name string
}

func NewObjectType(name string) *ObjectType  { return &ObjectType{name: name} }
func (t *ObjectType) IsJavaNil_() bool       { return t == nil }
func (t *ObjectType) TypeName() string       { return t.name }

type JSDocInfo struct {
	Summary string
}

func (d *JSDocInfo) IsJavaNil_() bool { return d == nil }

type Node nodestore.Node

func (n Node) IsJavaNil_() bool       { return n == 0 }
func (n Node) GetToken() Token        { return Token(nodestore.GetRec(nodestore.Node(n)).Token + 1) }
func (n Node) TokenSwitchOrd_() int32 { return int32(nodestore.GetRec(nodestore.Node(n)).Token) }
func (n Node) IsScript() bool         { return n != 0 && n.GetToken() == TokSCRIPT }
func (n Node) IsFunction() bool       { return n != 0 && n.GetToken() == TokFUNCTION }
func (n Node) IsBlock() bool          { return n != 0 && n.GetToken() == TokBLOCK }
func (n Node) IsName() bool           { return n != 0 && n.GetToken() == TokNAME }
func (n Node) IsGetProp() bool        { return n != 0 && n.GetToken() == TokGETPROP }
func (n Node) GetFirstChild() Node    { return Node(nodestore.GetRec(nodestore.Node(n)).First) }
func (n Node) nodeGet_first() Node    { return n.GetFirstChild() }
func (n Node) GetNext() Node          { return Node(nodestore.GetRec(nodestore.Node(n)).Next) }
func (n Node) GetParent() Node        { return Node(nodestore.GetRec(nodestore.Node(n)).Parent) }
func (n Node) GetString() string      { return nodestore.GetRec(nodestore.Node(n)).Payload }

func (n Node) AddChildToBack(c Node) {
	nodestore.MutateRec(nodestore.Node(c), func(r *nodestore.NodeRec) { r.Parent = nodestore.Node(n) })
	first := n.GetFirstChild()
	if first == 0 {
		nodestore.MutateRec(nodestore.Node(n), func(r *nodestore.NodeRec) { r.First = nodestore.Node(c) })
		return
	}
	cur := first
	for cur.GetNext() != 0 {
		cur = cur.GetNext()
	}
	nodestore.MutateRec(nodestore.Node(cur), func(r *nodestore.NodeRec) { r.Next = nodestore.Node(c) })
}

func NewNode_1(tok Token) Node {
	return Node(nodestore.AllocRaw(nodestore.NodeRec{Token: uint8(tok) - 1}))
}
func NewNode_StringNode_2(tok Token, s string) Node {
	return Node(nodestore.AllocRaw(nodestore.NodeRec{Token: uint8(tok) - 1, Payload: s}))
}
func NewGetProp(recv Node, prop string) Node {
	gp := Node(nodestore.AllocRaw(nodestore.NodeRec{Token: uint8(TokGETPROP) - 1, Payload: prop}))
	gp.AddChildToBack(recv)
	return gp
}

// Raw j2go nodeGet/nodeSet accessors matching ClassEmitter.java output
func (this_ Node) nodeGet_jsDocInfo() *JSDocInfo {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if d := nodestore.Node(this_).JSDocObj(); d != nil {
		return d.(*JSDocInfo)
	}
	return nil
}
func (this_ Node) nodeSet_jsDocInfo(v *JSDocInfo) {
	if v == nil {
		nodestore.Node(this_).SetJSDocObj(nil)
	} else {
		nodestore.Node(this_).SetJSDocObj(v)
	}
}
func (this_ Node) nodeGet_typeBeforeCast() JSType {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if t := nodestore.Node(this_).JSTypeBeforeCastObj(); t != nil {
		return t.(JSType)
	}
	return nil
}
func (this_ Node) nodeSet_typeBeforeCast(v JSType) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if v == nil || v.IsJavaNil_() {
		nodestore.Node(this_).SetJSTypeBeforeCastObj(nil)
	} else {
		nodestore.Node(this_).SetJSTypeBeforeCastObj(v)
	}
}

func (this_ Node) GetJSTypeBeforeCast() JSType {
	return this_.nodeGet_typeBeforeCast()
}
func (this_ Node) SetJSTypeBeforeCast(type_ JSType) {
	_ = type_
	this_.nodeSet_typeBeforeCast(type_)
}

func (this_ Node) GetJSDocInfo() *JSDocInfo {
	return this_.nodeGet_jsDocInfo()
}
func (this_ Node) SetJSDocInfo(info *JSDocInfo) Node {
	_ = info
	this_.nodeSet_jsDocInfo(info)
	return this_
}

// 2..8-part dot-separated concatenation in GetQualifiedName()
func (this_ Node) GetQualifiedName() string {
	switch this_.TokenSwitchOrd_() {
	case int32(32):
		var name string = this_.GetString()
		if j2gort.CheckNotNullString(name) == "" {
			return j2gort.NullString
		}
		return name
	case int32(26):
		var p1 string = this_.GetString()
		var n1 Node = this_.nodeGet_first()
		if !(n1.IsGetProp()) {
			var left string = n1.GetQualifiedName()
			if left == j2gort.NullString {
				return j2gort.NullString
			}
			return (j2gort.StrValString(left) + "." + j2gort.StrValString(p1))
		}
		var p2 string = n1.GetString()
		var n2 Node = n1.nodeGet_first()
		if !(n2.IsGetProp()) {
			var left_2 string = n2.GetQualifiedName()
			if left_2 == j2gort.NullString {
				return j2gort.NullString
			}
			return (j2gort.StrValString(left_2) + "." + j2gort.StrValString(p2) + "." + j2gort.StrValString(p1))
		}
		var p3 string = n2.GetString()
		var n3 Node = n2.nodeGet_first()
		if !(n3.IsGetProp()) {
			var left_3 string = n3.GetQualifiedName()
			if left_3 == j2gort.NullString {
				return j2gort.NullString
			}
			return (j2gort.StrValString(left_3) + "." + j2gort.StrValString(p3) + "." + j2gort.StrValString(p2) + "." + j2gort.StrValString(p1))
		}
		var p4 string = n3.GetString()
		var n4 Node = n3.nodeGet_first()
		if !(n4.IsGetProp()) {
			var left_4 string = n4.GetQualifiedName()
			if left_4 == j2gort.NullString {
				return j2gort.NullString
			}
			return (j2gort.StrValString(left_4) + "." + j2gort.StrValString(p4) + "." + j2gort.StrValString(p3) + "." + j2gort.StrValString(p2) + "." + j2gort.StrValString(p1))
		}
		var p5 string = n4.GetString()
		var n5 Node = n4.nodeGet_first()
		if !(n5.IsGetProp()) {
			var left_5 string = n5.GetQualifiedName()
			if left_5 == j2gort.NullString {
				return j2gort.NullString
			}
			return (j2gort.StrValString(left_5) + "." + j2gort.StrValString(p5) + "." + j2gort.StrValString(p4) + "." + j2gort.StrValString(p3) + "." + j2gort.StrValString(p2) + "." + j2gort.StrValString(p1))
		}
		var p6 string = n5.GetString()
		var n6 Node = n5.nodeGet_first()
		if !(n6.IsGetProp()) {
			var left_6 string = n6.GetQualifiedName()
			if left_6 == j2gort.NullString {
				return j2gort.NullString
			}
			return (j2gort.StrValString(left_6) + "." + j2gort.StrValString(p6) + "." + j2gort.StrValString(p5) + "." + j2gort.StrValString(p4) + "." + j2gort.StrValString(p3) + "." + j2gort.StrValString(p2) + "." + j2gort.StrValString(p1))
		}
		var p7 string = n6.GetString()
		var n7 Node = n6.nodeGet_first()
		var left_7 string = n7.GetQualifiedName()
		if left_7 == j2gort.NullString {
			return j2gort.NullString
		}
		return (j2gort.StrValString(left_7) + "." + j2gort.StrValString(p7) + "." + j2gort.StrValString(p6) + "." + j2gort.StrValString(p5) + "." + j2gort.StrValString(p4) + "." + j2gort.StrValString(p3) + "." + j2gort.StrValString(p2) + "." + j2gort.StrValString(p1))
	case int32(37):
		return "this"
	default:
		return j2gort.NullString
	}
}

// Negative control: slash-separated chain must NOT be rewritten to JoinDot!
func SlashPath4(a, b, c, d string) string {
	return a + "/" + b + "/" + c + "/" + d
}
`)

	// 4. gen/jscomp covering Task 3a (GlobalNamespace), Task 3b (InlineFunctions),
	//    Task 3c (ReferenceCollection.rest + ReferenceCollector), and Task 3d (CrossChunkReferenceCollector + CrossChunkCodeMotion).
	mustWrite("gen/jscomp/scope_and_traversal.go", `package jscomp

import (
	"example.com/w28ctest/gen/j2gort"
	pkg_rhino "example.com/w28ctest/gen/rhino"
)

type AbstractCompiler interface {
	IsJavaNil_() bool
}
type Compiler struct{}
func (c *Compiler) IsJavaNil_() bool { return c == nil }

type CompilerPass interface {
	Implements_jscomp_CompilerPass()
	Process_2(externs, root pkg_rhino.Node)
}

type Var struct {
	name  string
	node  pkg_rhino.Node
	scope *Scope
}
func NewVar(name string, n pkg_rhino.Node, s *Scope) *Var {
	v := &Var{}
	InitVar(v, name, n, s)
	return v
}
func InitVar(this_ *Var, name string, n pkg_rhino.Node, s *Scope) {
	this_.name = name
	this_.node = n
	this_.scope = s
}
func (v *Var) GetName() string         { return v.name }
func (v *Var) GetScope() *Scope        { return v.scope }
func (v *Var) GetNode() pkg_rhino.Node { return v.node }

type Scope struct {
	parent *Scope
	root   pkg_rhino.Node
	vars   *j2gort.OrderedMap[string, *Var]
}
func NewScope_2(parent *Scope, root pkg_rhino.Node) *Scope {
	s := &Scope{}
	InitScope_2(s, parent, root)
	return s
}
func InitScope_2(this_ *Scope, parent *Scope, root pkg_rhino.Node) {
	this_.initFields_Scope()
	this_.parent = parent
	this_.root = root
}
func (this_ *Scope) initFields_Scope() {
	this_.vars = j2gort.NewOrderedMapPtrCap[string, *Var](int32(16))
}
func (s *Scope) GetRootNode_0() pkg_rhino.Node { return s.root }
func (s *Scope) Declare(name string, n pkg_rhino.Node, v *Var) {
	s.vars.Put(name, v)
}
func (s *Scope) GetVar(name string) *Var {
	if found := s.vars.Get(name); found != nil {
		return found
	}
	if s.parent != nil {
		return s.parent.GetVar(name)
	}
	return nil
}

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope(n pkg_rhino.Node, parent *Scope) *Scope
}

type SyntacticScopeCreator struct {
	compiler AbstractCompiler
}
func (s *SyntacticScopeCreator) IsJavaNil_() bool { return s == nil }
func NewSyntacticScopeCreator_1(c AbstractCompiler) *SyntacticScopeCreator {
	s := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_1(s, c)
	return s
}
func InitSyntacticScopeCreator_1(this_ *SyntacticScopeCreator, c AbstractCompiler) {
	this_.compiler = c
}
func (this_ *SyntacticScopeCreator) CreateScope(n pkg_rhino.Node, parent *Scope) *Scope {
	s := NewScope_2(parent, n)
	for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
		if c.IsName() {
			s.Declare(c.GetString(), c, NewVar(c.GetString(), c, s))
		}
	}
	return s
}

type NodeTraversal_Callback interface {
	ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool
	Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node)
}
type NodeTraversal_ScopedCallback interface {
	NodeTraversal_Callback
	EnterScope(t *NodeTraversal)
	ExitScope(t *NodeTraversal)
}

type NodeTraversal struct {
	compiler     AbstractCompiler
	callback     NodeTraversal_Callback
	scopeCreator ScopeCreator
	scopes       []*Scope
}
func (t *NodeTraversal) GetScope() *Scope { return t.scopes[len(t.scopes)-1] }
func (t *NodeTraversal) pushScope(n pkg_rhino.Node) {
	var parent *Scope
	if len(t.scopes) > 0 {
		parent = t.scopes[len(t.scopes)-1]
	}
	s := t.scopeCreator.CreateScope(n, parent)
	t.scopes = append(t.scopes, s)
	if sc, ok := t.callback.(NodeTraversal_ScopedCallback); ok {
		sc.EnterScope(t)
	}
}
func (t *NodeTraversal) popScope() {
	if sc, ok := t.callback.(NodeTraversal_ScopedCallback); ok {
		sc.ExitScope(t)
	}
	t.scopes = t.scopes[:len(t.scopes)-1]
}
func (t *NodeTraversal) walk(n, parent pkg_rhino.Node) {
	if n == 0 || !t.callback.ShouldTraverse_3(t, n, parent) {
		return
	}
	pushed := parent == 0 || n.IsFunction() || n.IsBlock()
	if pushed {
		t.pushScope(n)
	}
	for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
		t.walk(c, n)
	}
	t.callback.Visit_3(t, n, parent)
	if pushed {
		t.popScope()
	}
}

type NodeTraversal_Builder struct {
	compiler     AbstractCompiler
	callback     NodeTraversal_Callback
	scopeCreator ScopeCreator
}
func NodeTraversal_S_Builder() *NodeTraversal_Builder { return &NodeTraversal_Builder{} }
func (b *NodeTraversal_Builder) SetCompiler(c AbstractCompiler) *NodeTraversal_Builder {
	b.compiler = c
	return b
}
func (b *NodeTraversal_Builder) SetCallback_1_NodeTraversal_Callback(cb NodeTraversal_Callback) *NodeTraversal_Builder {
	b.callback = cb
	return b
}
func (b *NodeTraversal_Builder) SetScopeCreator(sc ScopeCreator) *NodeTraversal_Builder {
	b.scopeCreator = sc
	return b
}
func (b *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	sc := b.scopeCreator
	if sc == nil {
		sc = NewSyntacticScopeCreator_1(b.compiler)
	}
	nt := &NodeTraversal{compiler: b.compiler, callback: b.callback, scopeCreator: sc}
	nt.walk(root, 0)
}
func (b *NodeTraversal_Builder) TraverseRoots(externs, root pkg_rhino.Node) {
	b.Traverse(externs)
	b.Traverse(root)
}
func NodeTraversal_Traverse(c AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	NodeTraversal_S_Builder().SetCompiler(c).SetCallback_1_NodeTraversal_Callback(cb).Traverse(root)
}
`)

	// Task 3a: GlobalNamespace + GlobalNamespace_Name.refsForNodeMap + GlobalNamespace_ModuleNameMap
	mustWrite("gen/jscomp/global_namespace.go", `package jscomp

import (
	"example.com/w28ctest/gen/j2gort"
	pkg_rhino "example.com/w28ctest/gen/rhino"
)

type GlobalNamespace_Ref struct {
	node  pkg_rhino.Node
	scope *Scope
}

type GlobalNamespace_Name struct {
	baseName       string
	ref0           *GlobalNamespace_Ref
	ref1           *GlobalNamespace_Ref
	refsForNodeMap *j2gort.OrderedMap[pkg_rhino.Node, *GlobalNamespace_Ref]
}

func NewGlobalNamespace_Name(name string) *GlobalNamespace_Name {
	return &GlobalNamespace_Name{baseName: name}
}

func (this_ *GlobalNamespace_Name) AddRef(node pkg_rhino.Node, scope *Scope) {
	ref := &GlobalNamespace_Ref{node: node, scope: scope}
	if this_.refsForNodeMap != nil {
		this_.refsForNodeMap.Put(node, ref)
		return
	}
	if this_.ref0 == nil {
		this_.ref0 = ref
		return
	}
	if this_.ref1 == nil {
		this_.ref1 = ref
		return
	}
	r0 := this_.ref0
	r1 := this_.ref1
	this_.ref0 = nil
	this_.ref1 = nil
	this_.refsForNodeMap = j2gort.NewOrderedMapPtrCap[pkg_rhino.Node, *GlobalNamespace_Ref](int32(8))
	this_.refsForNodeMap.Put(r0.node, r0)
	this_.refsForNodeMap.Put(r1.node, r1)
	this_.refsForNodeMap.Put(node, ref)
}

type ModuleMetadata struct {
	id string
}

type GlobalNamespace_ModuleNameMap struct {
	names j2gort.OrderedMap[string, *GlobalNamespace_Name]
}

func NewGlobalNamespace_ModuleNameMap() *GlobalNamespace_ModuleNameMap {
	m := &GlobalNamespace_ModuleNameMap{}
	m.initFields_GlobalNamespace_ModuleNameMap()
	return m
}
func (this_ *GlobalNamespace_ModuleNameMap) initFields_GlobalNamespace_ModuleNameMap() {
	this_.names = j2gort.NewOrderedMap[string, *GlobalNamespace_Name]()
}

type GlobalNamespace struct {
	compiler        AbstractCompiler
	root            pkg_rhino.Node
	globalNames     []*GlobalNamespace_Name
	nameMap         j2gort.OrderedMap[string, *GlobalNamespace_Name]
	nameMapByModule j2gort.CustomOrderedMap[*ModuleMetadata, *GlobalNamespace_ModuleNameMap]
	generated       bool
}

func NewGlobalNamespace_2(c AbstractCompiler, root pkg_rhino.Node) *GlobalNamespace {
	ns := &GlobalNamespace{compiler: c, root: root}
	ns.initFields_GlobalNamespace()
	return ns
}
func (this_ *GlobalNamespace) initFields_GlobalNamespace() {
	this_.nameMap = j2gort.NewOrderedMap[string, *GlobalNamespace_Name]()
	this_.nameMapByModule = j2gort.NewCustomOrderedMap[*ModuleMetadata, *GlobalNamespace_ModuleNameMap]()
}

func (this_ *GlobalNamespace) EnsureGenerated() {
	if !this_.generated {
		this_.priv_GlobalNamespace_process_3001()
		this_.generated = true
	}
}

func (this_ *GlobalNamespace) priv_GlobalNamespace_process_3001() {
	var builder *NodeTraversal_Builder = NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(&GlobalNamespace_BuildGlobalNamespace{outer_: this_})
	builder.Traverse(this_.root)
}

func (this_ *GlobalNamespace) PutModuleRef(mod *ModuleMetadata, name string, node pkg_rhino.Node, scope *Scope) {
	mm := this_.nameMapByModule.Get(mod)
	if mm == nil {
		mm = NewGlobalNamespace_ModuleNameMap()
		this_.nameMapByModule.Put(mod, mm)
	}
	nm := mm.names.Get(name)
	if nm == nil {
		nm = NewGlobalNamespace_Name(name)
		mm.names.Put(name, nm)
	}
	nm.AddRef(node, scope)
}

type GlobalNamespace_BuildGlobalNamespace struct {
	outer_ *GlobalNamespace
}
func (*GlobalNamespace_BuildGlobalNamespace) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}
func (b *GlobalNamespace_BuildGlobalNamespace) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		nm := b.outer_.nameMap.Get(n.GetString())
		if nm == nil {
			nm = NewGlobalNamespace_Name(n.GetString())
			b.outer_.globalNames = append(b.outer_.globalNames, nm)
			b.outer_.nameMap.Put(n.GetString(), nm)
		}
		nm.AddRef(n, t.GetScope())
	}
}

type InlineAndCollapseProperties struct {
	compiler  AbstractCompiler
	namespace *GlobalNamespace
	Count     int
}
func (*InlineAndCollapseProperties) Implements_jscomp_CompilerPass() {}
func (this_ *InlineAndCollapseProperties) Process_2(_ pkg_rhino.Node, root pkg_rhino.Node) {
	this_.namespace = NewGlobalNamespace_2(this_.compiler, root)
	this_.namespace.EnsureGenerated()
	mod := &ModuleMetadata{id: "m0"}
	for c := root.GetFirstChild(); c != 0; c = c.GetNext() {
		this_.namespace.PutModuleRef(mod, "modSym", c, nil)
	}
	this_.Count = int(this_.namespace.nameMap.Size())
}
`)

	// Task 3b: InlineFunctions + InlineFunctions_FunctionState (including DisallowInlining_3 & Remove)
	mustWrite("gen/jscomp/inline_functions.go", `package jscomp

import (
	"example.com/w28ctest/gen/j2gort"
	pkg_rhino "example.com/w28ctest/gen/rhino"
)

type InlineFunctions_Reference struct {
	callNode pkg_rhino.Node
	scope    *Scope
}

type InlineFunctions_FunctionState struct {
	name       string
	inline     bool
	lowersCost bool
	ref0       *InlineFunctions_Reference
	ref1       *InlineFunctions_Reference
	references *j2gort.OrderedMap[pkg_rhino.Node, *InlineFunctions_Reference]
}

func (this_ *InlineFunctions_FunctionState) AddReference(ref *InlineFunctions_Reference) {
	if this_.references != nil {
		this_.references.Put(ref.callNode, ref)
		return
	}
	if this_.ref0 == nil {
		this_.ref0 = ref
		return
	}
	if this_.ref1 == nil {
		this_.ref1 = ref
		return
	}
	r0 := this_.ref0
	r1 := this_.ref1
	this_.ref0 = nil
	this_.ref1 = nil
	this_.references = j2gort.NewOrderedMapPtrCap[pkg_rhino.Node, *InlineFunctions_Reference](int32(8))
	this_.references.Put(r0.callNode, r0)
	this_.references.Put(r1.callNode, r1)
	this_.references.Put(ref.callNode, ref)
}

func (this_ *InlineFunctions_FunctionState) DisallowInlining_3() {
	this_.inline = false
	this_.ref0 = nil
	this_.ref1 = nil
	this_.references = nil
}

type InlineFunctions struct {
	compiler     AbstractCompiler
	fns          j2gort.OrderedMap[string, *InlineFunctions_FunctionState]
	anonFns      j2gort.OrderedMap[pkg_rhino.Node, string]
	InlinedCount int
}

func (*InlineFunctions) Implements_jscomp_CompilerPass() {}

func NewInlineFunctions(c AbstractCompiler) *InlineFunctions {
	p := &InlineFunctions{compiler: c}
	p.initFields_InlineFunctions()
	return p
}
func (this_ *InlineFunctions) initFields_InlineFunctions() {
	this_.fns = j2gort.NewOrderedMap[string, *InlineFunctions_FunctionState]()
	this_.anonFns = j2gort.NewOrderedMap[pkg_rhino.Node, string]()
}

func (this_ *InlineFunctions) trimCandidatesUsingOnCost() {
	var toRemove []string
	cur := this_.fns.IterEntries()
	for cur.NextEntry() {
		fs := cur.Value()
		if !fs.lowersCost || !fs.inline {
			toRemove = append(toRemove, cur.Key())
		}
	}
	for i := 0; i < len(toRemove); i++ {
		this_.fns.Remove(toRemove[i])
	}
}

func (this_ *InlineFunctions) Process_2(_ pkg_rhino.Node, root pkg_rhino.Node) {
	this_.fns.Put("keepFn", &InlineFunctions_FunctionState{name: "keepFn", inline: true, lowersCost: true})
	this_.fns.Put("trimCostFn", &InlineFunctions_FunctionState{name: "trimCostFn", inline: true, lowersCost: false})
	this_.fns.Put("disallowFn", &InlineFunctions_FunctionState{name: "disallowFn", inline: true, lowersCost: true})
	this_.anonFns.Put(root, "keepFn")

	NodeTraversal_Traverse(this_.compiler, root, NewInlineFunctions_FindCandidatesReferences(this_))

	if d := this_.fns.Get("disallowFn"); d != nil {
		d.DisallowInlining_3()
	}
	this_.trimCandidatesUsingOnCost()

	cur := this_.fns.IterValues()
	for cur.HasNext() {
		fs := cur.Next()
		if fs.references != nil {
			this_.InlinedCount += int(fs.references.Size())
		}
	}
}

type InlineFunctions_FindCandidatesReferences struct {
	owner *InlineFunctions
}
func NewInlineFunctions_FindCandidatesReferences(o *InlineFunctions) *InlineFunctions_FindCandidatesReferences {
	return &InlineFunctions_FindCandidatesReferences{owner: o}
}
func (*InlineFunctions_FindCandidatesReferences) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}
func (f *InlineFunctions_FindCandidatesReferences) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		cur := f.owner.fns.IterValues()
		for cur.HasNext() {
			cur.Next().AddReference(&InlineFunctions_Reference{callNode: n, scope: t.GetScope()})
		}
	}
}
`)

	// Task 3c + 3d: ReferenceCollection (raw ArrayList before native_slice_lowering!), ReferenceCollector,
	//InferConsts (non-escaping), SymbolTable (escaping control), CrossChunkReferenceCollector, and CrossChunkCodeMotion!
	mustWrite("gen/jscomp/references_and_cccm.go", `package jscomp

import (
	"example.com/w28ctest/gen/j2gort"
	pkg_rhino "example.com/w28ctest/gen/rhino"
)

type BasicBlock struct {
	root pkg_rhino.Node
}
func NewBasicBlock(r pkg_rhino.Node) *BasicBlock {
	b := &BasicBlock{}
	InitBasicBlock(b, r)
	return b
}
func InitBasicBlock(this_ *BasicBlock, r pkg_rhino.Node) {
	this_.root = r
}

type Reference struct {
	nameNode   pkg_rhino.Node
	basicBlock *BasicBlock
	scope      *Scope
}
func NewReference(n pkg_rhino.Node, b *BasicBlock, s *Scope) *Reference {
	r := &Reference{}
	InitReference(r, n, b, s)
	return r
}
func InitReference(this_ *Reference, n pkg_rhino.Node, b *BasicBlock, s *Scope) {
	this_.nameNode = n
	this_.basicBlock = b
	this_.scope = s
}
func (r *Reference) GetNode() pkg_rhino.Node   { return r.nameNode }
func (r *Reference) GetBasicBlock() *BasicBlock { return r.basicBlock }
func (r *Reference) GetScope() *Scope          { return r.scope }
func (r *Reference) GetSymbol() *Var           { return r.scope.GetVar(r.nameNode.GetString()) }

// Raw ReferenceCollection before native_slice_lowering lowers rest to []*Reference!
type ReferenceCollection struct {
	ref0 *Reference
	ref1 *Reference
	ref2 *Reference
	ref3 *Reference
	rest j2gort.Ext_java_util_ArrayList
}

func NewReferenceCollection() *ReferenceCollection {
	rc := &ReferenceCollection{}
	InitReferenceCollection(rc)
	return rc
}
func InitReferenceCollection(_ *ReferenceCollection) {}

func (this_ *ReferenceCollection) Size() int32 {
	if this_.rest != nil {
		return this_.rest.Size()
	}
	if this_.ref0 == nil {
		return 0
	}
	if this_.ref1 == nil {
		return 1
	}
	if this_.ref2 == nil {
		return 2
	}
	if this_.ref3 == nil {
		return 3
	}
	return 4
}

func (this_ *ReferenceCollection) Get(index int32) *Reference {
	if this_.rest != nil {
		return j2gort.CastPtr[Reference](this_.rest.Get(index))
	}
	switch index {
	case 0:
		return this_.ref0
	case 1:
		return this_.ref1
	case 2:
		return this_.ref2
	case 3:
		return this_.ref3
	default:
		panic("oob")
	}
}

func (this_ *ReferenceCollection) Add(reference *Reference) {
	if this_.rest != nil {
		this_.rest.Add_1(j2gort.PtrToAny(reference))
	} else if this_.ref0 == nil {
		this_.ref0 = reference
	} else if this_.ref1 == nil {
		this_.ref1 = reference
	} else if this_.ref2 == nil {
		this_.ref2 = reference
	} else if this_.ref3 == nil {
		this_.ref3 = reference
	} else {
		var list j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList_1_int32(int32(8))
		_ = list
		list.Add_1(j2gort.PtrToAny(this_.ref0))
		list.Add_1(j2gort.PtrToAny(this_.ref1))
		list.Add_1(j2gort.PtrToAny(this_.ref2))
		list.Add_1(j2gort.PtrToAny(this_.ref3))
		list.Add_1(j2gort.PtrToAny(reference))
		this_.rest = list
		this_.ref0 = nil
		this_.ref1 = nil
		this_.ref2 = nil
		this_.ref3 = nil
	}
}

type ReferenceCollector_Behavior interface {
	BeforeTraverseNode(t *NodeTraversal, n, parent pkg_rhino.Node)
	AfterVisitNode(t *NodeTraversal, n, parent pkg_rhino.Node, addedRef *Reference)
}

type doNothingBehavior struct{}
func (doNothingBehavior) BeforeTraverseNode(_ *NodeTraversal, _, _ pkg_rhino.Node)               {}
func (doNothingBehavior) AfterVisitNode(_ *NodeTraversal, _, _ pkg_rhino.Node, _ *Reference)     {}
func StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR() ReferenceCollector_Behavior {
	return doNothingBehavior{}
}

var lambdaSingleton_rc_new j2gort.Ext_java_util_function_Function
func init() {
	lambdaSingleton_rc_new = j2gort.NewLambda_Ext_java_util_function_Function(func(lp any) any {
		k := j2gort.CastPtr[Var](lp)
		_ = k
		return NewReferenceCollection()
	})
}

type ReferenceCollector struct {
	compiler     AbstractCompiler
	behavior     ReferenceCollector_Behavior
	scopeCreator ScopeCreator
	referenceMap j2gort.CustomOrderedMap[*Var, *ReferenceCollection]
	curBlock     *BasicBlock
}
func (*ReferenceCollector) Implements_jscomp_CompilerPass() {}

func NewReferenceCollector_3(c AbstractCompiler, b ReferenceCollector_Behavior, sc ScopeCreator) *ReferenceCollector {
	rc := &ReferenceCollector{}
	InitReferenceCollector_3(rc, c, b, sc)
	return rc
}
func InitReferenceCollector_3(this_ *ReferenceCollector, c AbstractCompiler, b ReferenceCollector_Behavior, sc ScopeCreator) {
	this_.compiler = c
	this_.behavior = b
	this_.scopeCreator = sc
	this_.referenceMap = j2gort.NewCustomOrderedMap[*Var, *ReferenceCollection]()
}

func (this_ *ReferenceCollector) Process_2(externs, root pkg_rhino.Node) {
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(this_).SetScopeCreator(this_.scopeCreator).TraverseRoots(externs, root)
}
func (this_ *ReferenceCollector) Process_1(root pkg_rhino.Node) {
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(this_).SetScopeCreator(this_.scopeCreator).Traverse(root)
}
func (this_ *ReferenceCollector) ShouldTraverse_3(t *NodeTraversal, n, parent pkg_rhino.Node) bool {
	if n.IsScript() || n.IsBlock() {
		this_.curBlock = NewBasicBlock(n)
	}
	this_.behavior.BeforeTraverseNode(t, n, parent)
	return true
}
func (this_ *ReferenceCollector) Visit_3(t *NodeTraversal, n, parent pkg_rhino.Node) {
	var added *Reference
	if n.IsName() {
		if v := t.GetScope().GetVar(n.GetString()); v != nil {
			var rVal Reference
			InitReference(&rVal, n, this_.curBlock, t.GetScope())
			added = &rVal
			this_.priv_addRef(v, added)
		}
	}
	this_.behavior.AfterVisitNode(t, n, parent, added)
}
func (this_ *ReferenceCollector) priv_addRef(v *Var, r *Reference) {
	var rc *ReferenceCollection = j2gort.CustomOrderedMapComputeIfAbsentFn(&this_.referenceMap, v, nil, lambdaSingleton_rc_new)
	rc.Add(r)
}
func (this_ *ReferenceCollector) GetAllSymbols() j2gort.Ext_java_util_Collection {
	return this_.referenceMap.LiveKeySet()
}
func (this_ *ReferenceCollector) GetReferences(v *Var) *ReferenceCollection {
	return this_.referenceMap.Get(v)
}

// Non-escaping ReferenceCollector caller -> releases slabs + pooled rc_.rest slice!
type InferConsts struct {
	compiler  AbstractCompiler
	TotalRefs int32
}
func (*InferConsts) Implements_jscomp_CompilerPass() {}
func (this_ *InferConsts) Process_2(externs, root pkg_rhino.Node) {
	var val_1 ReferenceCollector
	InitReferenceCollector_3(&val_1, this_.compiler, StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR(), NewSyntacticScopeCreator_1(this_.compiler))
	var collector *ReferenceCollector = &val_1
	_ = collector
	collector.Process_2(externs, root)
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(collector.GetAllSymbols()))
	for cur.HasNext() {
		v := j2gort.CastPtr[Var](cur.Next())
		this_.TotalRefs += collector.GetReferences(v).Size()
	}
}

// Escaping ReferenceCollector caller (SymbolTable) -> MUST NOT call ReleaseSlabs_()!
type SymbolTable struct {
	saved []*ReferenceCollection
}
func (st *SymbolTable) AddSymbolsFrom(rc *ReferenceCollector) {
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(rc.GetAllSymbols()))
	for cur.HasNext() {
		v := j2gort.CastPtr[Var](cur.Next())
		st.saved = append(st.saved, rc.GetReferences(v))
	}
}
func BuildSymbolTable(c AbstractCompiler, externs, root pkg_rhino.Node) *SymbolTable {
	var val_81 ReferenceCollector
	InitReferenceCollector_3(&val_81, c, StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR(), NewSyntacticScopeCreator_1(c))
	refCollector := &val_81
	refCollector.Process_2(externs, root)
	st := &SymbolTable{}
	st.AddSymbolsFrom(refCollector)
	return st
}

// Task 3d: CrossChunkReferenceCollector + CrossChunkCodeMotion
type CrossChunkReferenceCollector_TopLevelStatement struct {
	statementNode            pkg_rhino.Node
	declaredNameReference    *Reference
	nonDeclarationReferences []*Reference
}
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetDeclaredNameReference() *Reference {
	return j2gort.CheckNotNullPtr(s.declaredNameReference)
}
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetNonDeclarationReferenceCount() int32 {
	return int32(len(s.nonDeclarationReferences))
}
func (s *CrossChunkReferenceCollector_TopLevelStatement) GetNonDeclarationReference(i int32) *Reference {
	return s.nonDeclarationReferences[i]
}

type CrossChunkReferenceCollector struct {
	topLevelStatements     []*CrossChunkReferenceCollector_TopLevelStatement
	compiler               AbstractCompiler
	referenceCollector     *ReferenceCollector
	topLevelStatementDraft *CrossChunkReferenceCollector_TopLevelStatement
}
func (*CrossChunkReferenceCollector) Implements_jscomp_CompilerPass() {}

func NewCrossChunkReferenceCollector(compiler AbstractCompiler, creator ScopeCreator) *CrossChunkReferenceCollector {
	c := &CrossChunkReferenceCollector{}
	InitCrossChunkReferenceCollector(c, compiler, creator)
	return c
}
func InitCrossChunkReferenceCollector(this_ *CrossChunkReferenceCollector, compiler AbstractCompiler, creator ScopeCreator) {
	this_.compiler = compiler
	this_.referenceCollector = NewReferenceCollector_3(compiler, &CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior{outer_: this_}, creator)
}
func (this_ *CrossChunkReferenceCollector) Process_2(externs, root pkg_rhino.Node) {
	this_.referenceCollector.Process_2(externs, root)
}
func (this_ *CrossChunkReferenceCollector) Process_1(root pkg_rhino.Node) {
	this_.referenceCollector.Process_1(root)
}
func (this_ *CrossChunkReferenceCollector) GetTopLevelStatementCount() int32 {
	return int32(len(this_.topLevelStatements))
}
func (this_ *CrossChunkReferenceCollector) GetTopLevelStatement(i int32) *CrossChunkReferenceCollector_TopLevelStatement {
	return this_.topLevelStatements[i]
}
func (this_ *CrossChunkReferenceCollector) GetReferences(v *Var) *ReferenceCollection {
	return this_.referenceCollector.GetReferences(v)
}

type CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior struct {
	outer_ *CrossChunkReferenceCollector
}
func (b *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) BeforeTraverseNode(_ *NodeTraversal, n, parent pkg_rhino.Node) {
	if parent != 0 && parent.IsScript() {
		b.outer_.topLevelStatementDraft = &CrossChunkReferenceCollector_TopLevelStatement{statementNode: n}
	}
}
func (b *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) AfterVisitNode(_ *NodeTraversal, n, _ pkg_rhino.Node, added *Reference) {
	if b.outer_.topLevelStatementDraft == nil {
		return
	}
	if n == b.outer_.topLevelStatementDraft.statementNode {
		b.outer_.topLevelStatements = append(b.outer_.topLevelStatements, b.outer_.topLevelStatementDraft)
		b.outer_.topLevelStatementDraft = nil
	} else if added != nil {
		if b.outer_.topLevelStatementDraft.declaredNameReference == nil {
			b.outer_.topLevelStatementDraft.declaredNameReference = added
		} else {
			b.outer_.topLevelStatementDraft.nonDeclarationReferences = append(b.outer_.topLevelStatementDraft.nonDeclarationReferences, added)
		}
	}
}

type CrossChunkCodeMotion struct {
	compiler    AbstractCompiler
	chunkCount  int32
	MovedRefs   int32
	cccmLogOpen bool
}
func (*CrossChunkCodeMotion) Implements_jscomp_CompilerPass() {}

type CrossChunkCodeMotion_GlobalSymbolCollector struct {
	outer_ *CrossChunkCodeMotion
	total  int32
}
func (g *CrossChunkCodeMotion_GlobalSymbolCollector) CollectGlobalSymbols(referenceCollector *CrossChunkReferenceCollector) {
	for i := int32(0); i < referenceCollector.GetTopLevelStatementCount(); i++ {
		stmt := referenceCollector.GetTopLevelStatement(i)
		if stmt.declaredNameReference != nil {
			// Reads live *Reference + *Var + *ReferenceCollection AFTER referenceCollector.Process_1(root) returned!
			sym := stmt.GetDeclaredNameReference().GetSymbol()
			g.total += referenceCollector.GetReferences(sym).Size()
		}
		for j := int32(0); j < stmt.GetNonDeclarationReferenceCount(); j++ {
			if stmt.GetNonDeclarationReference(j).GetNode() != 0 {
				g.total++
			}
		}
	}
}

func (this_ *CrossChunkCodeMotion) Process_2(externs, root pkg_rhino.Node) {
	_ = externs
	// Mirrors Java try-with-resources IIFE in gen/jscomp/cross_chunk_code_motion.go!
	func() {
		this_.cccmLogOpen = true
		defer func() {
			this_.cccmLogOpen = false
		}()
		if this_.chunkCount > int32(1) {
			var val_1 CrossChunkReferenceCollector
			InitCrossChunkReferenceCollector(&val_1, this_.compiler, NewSyntacticScopeCreator_1(this_.compiler))
			var referenceCollector *CrossChunkReferenceCollector = &val_1
			_ = referenceCollector
			referenceCollector.Process_1(root)
			var val_2 CrossChunkCodeMotion_GlobalSymbolCollector
			val_2.outer_ = this_
			var collector *CrossChunkCodeMotion_GlobalSymbolCollector = &val_2
			collector.CollectGlobalSymbols(referenceCollector)
			this_.MovedRefs = collector.total
		}
	}()
}
`)

	allPasses := passes.AllPasses()
	for _, relPkg := range []string{"gen/rhino", "gen/jscomp"} {
		pkgDir := filepath.Join(tmpDir, relPkg)
		res1, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 1 on %s failed: %v", relPkg, err)
		}
		if res1.TotalSites == 0 || res1.FilesModified == 0 {
			t.Fatalf("expected transformations on %s Run 1, got %+v", relPkg, res1)
		}
		// Run 2 MUST produce 0 sites and 0 modified files (100% Run-2 idempotency)!
		res2, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 2 on %s failed: %v", relPkg, err)
		}
		if res2.TotalSites != 0 || res2.FilesModified != 0 {
			t.Fatalf("expected 0 sites and 0 modified files on %s Run 2, got %+v", relPkg, res2)
		}
	}

	readLowered := func(rel string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(tmpDir, rel))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	normWS := func(s string) string {
		return strings.Join(strings.Fields(s), " ")
	}

	// Verify Task 1 + Task 2 in gen/rhino/node.go
	rhinoCode := readLowered("gen/rhino/node.go")
	for _, want := range []string{
		"nodestore.JoinDot(",
		"nodestore.JoinDot3(",
		"nodestore.JoinDot4(",
		"nodestore.JoinDot(nodestore.JoinDot4(",
		"nodestore.JoinDot3(nodestore.JoinDot4(",
		"nodestore.JoinDot4(nodestore.JoinDot4(",
		"nodestore.JoinDot(nodestore.JoinDot4(nodestore.JoinDot4(",
		"if type_ == nil || type_.IsJavaNil_() { nodestore.Node(this_).ClearJSTypeBeforeCast() return }",
		"if info == nil { nodestore.Node(this_).ClearJSDoc() return this_ }",
		`a + "/" + b + "/" + c + "/" + d`,
	} {
		if !strings.Contains(normWS(rhinoCode), want) {
			t.Fatalf("expected lowered gen/rhino/node.go to contain %q:\n%s", want, rhinoCode)
		}
	}

	// Verify Task 3a in gen/jscomp/global_namespace.go
	gnCode := readLowered("gen/jscomp/global_namespace.go")
	for _, want := range []string{
		"j2gort.AcquireOrderedMapPtrCap[pkg_rhino.Node, *GlobalNamespace_Ref](int32(8))",
		"for cur_ := this_.nameMap.IterValues(); cur_.HasNext();",
		"j2gort.ReleaseOrderedMap(nm_.refsForNodeMap)",
		"nm_.refsForNodeMap = nil",
		"this_.nameMap.ReleaseToPool()",
		"this_.nameMapByModule.ReleaseToPool()",
	} {
		if !strings.Contains(normWS(gnCode), want) {
			t.Fatalf("expected lowered global_namespace.go to contain %q:\n%s", want, gnCode)
		}
	}

	// Verify Task 3b in gen/jscomp/inline_functions.go
	ifCode := readLowered("gen/jscomp/inline_functions.go")
	for _, want := range []string{
		"j2gort.AcquireOrderedMapPtrCap[pkg_rhino.Node, *InlineFunctions_Reference](int32(8))",
		"for cur_ := this_.fns.IterValues(); cur_.HasNext();",
		"j2gort.ReleaseOrderedMap(fs_.references)",
		"fs_.references = nil",
		"this_.fns.ReleaseToPool()",
		"this_.anonFns.ReleaseToPool()",
	} {
		if !strings.Contains(normWS(ifCode), want) {
			t.Fatalf("expected lowered inline_functions.go to contain %q:\n%s", want, ifCode)
		}
	}

	// Verify Task 3c + 3d in gen/jscomp/references_and_cccm.go
	cccmCode := readLowered("gen/jscomp/references_and_cccm.go")
	for _, want := range []string{
		"j2gort.AcquirePooledSlice[*Reference](8)",
		"j2gort.AppendPooledSlice(&this_.rest, reference)",
		"j2gort.ReleasePooledSlice(rc_.rest)",
		"func (this_ *CrossChunkReferenceCollector) ReleaseSlabs_()",
		"referenceCollector.ownsScopeCreator_ = true",
		"defer referenceCollector.ReleaseSlabs_()",
	} {
		if !strings.Contains(normWS(cccmCode), want) {
			t.Fatalf("expected lowered references_and_cccm.go to contain %q:\n%s", want, cccmCode)
		}
	}
	// CrossChunkReferenceCollector.Process_2 / Process_1 MUST NOT release itself before caller reads TopLevelStatements!
	if strings.Contains(normWS(cccmCode), "func (this_ *CrossChunkReferenceCollector) Process_2(externs, root pkg_rhino.Node) { defer") ||
		strings.Contains(normWS(cccmCode), "func (this_ *CrossChunkReferenceCollector) Process_1(root pkg_rhino.Node) { defer") {
		t.Fatalf("CrossChunkReferenceCollector.Process_1/2 MUST NOT defer ReleaseSlabs_() prematurely:\n%s", cccmCode)
	}

	// Write runtime semantic + pool accounting + -race test across lowered packages!
	mustWrite("gen/jscomp/wave28c_rt_test.go", `package jscomp

import (
	"sync"
	"testing"

	"example.com/w28ctest/gen/j2gort"
	pkg_rhino "example.com/w28ctest/gen/rhino"
	"example.com/w28ctest/shims/nodestore"
)

func buildManyRefAST() (pkg_rhino.Node, pkg_rhino.Node) {
	ext := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
	root := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
	// 12 references per BLOCK so ReferenceCollection.rest spills (cap 8 -> grows to 16 via AppendPooledSlice!)
	// AND GlobalNamespace_Name.refsForNodeMap + InlineFunctions_FunctionState.references spill (>2 refs)!
	for s := 0; s < 3; s++ {
		blk := pkg_rhino.NewNode_1(pkg_rhino.TokBLOCK)
		for i := 0; i < 12; i++ {
			blk.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "hotVar"))
		}
		root.AddChildToBack(blk)
	}
	return ext, root
}

func TestWave28C_RuntimeSemanticsPoolAccountingAndRace(t *testing.T) {
	// 1. Task 1: 2..8-segment GetQualifiedName via JoinDot / JoinDot3 / JoinDot4!
	nodestore.ResetCounters()
	cur := pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "proto")
	segs := []string{"google", "ads", "googleads", "v17", "services", "CampaignService", "Mutate"}
	want := "proto"
	for _, seg := range segs {
		cur = pkg_rhino.NewGetProp(cur, seg)
		want += "." + seg
		if got := cur.GetQualifiedName(); got != want {
			t.Fatalf("GetQualifiedName mismatch: got %q, want %q", got, want)
		}
	}
	if nodestore.JoinDot4Count.Load() == 0 {
		t.Fatalf("expected JoinDot4 to fire on >=4-segment qualified names")
	}

	// 2. Task 2: Typed-nil & untyped-nil fast paths on SetJSTypeBeforeCast and SetJSDocInfo (zero boxing!)
	nodestore.ResetCounters()
	n := pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "x")
	var typedNilType *pkg_rhino.ObjectType = nil
	var ifaceTypedNil pkg_rhino.JSType = typedNilType

	n.SetJSTypeBeforeCast(nil)
	n.SetJSTypeBeforeCast(ifaceTypedNil)
	if nodestore.ClearTypeBeforeCastCnt.Load() != 2 || nodestore.SetTypeBeforeCastCnt.Load() != 0 {
		t.Fatalf("expected 2 ClearJSTypeBeforeCast and 0 SetJSTypeBeforeCastObj calls on nil/typed-nil, got clear=%d set=%d",
			nodestore.ClearTypeBeforeCastCnt.Load(), nodestore.SetTypeBeforeCastCnt.Load())
	}
	liveType := pkg_rhino.NewObjectType("string")
	n.SetJSTypeBeforeCast(liveType)
	if nodestore.SetTypeBeforeCastCnt.Load() != 1 || n.GetJSTypeBeforeCast().TypeName() != "string" {
		t.Fatalf("live SetJSTypeBeforeCast failed")
	}
	n.SetJSTypeBeforeCast(ifaceTypedNil)
	if n.GetJSTypeBeforeCast() != nil {
		t.Fatalf("expected typed-nil SetJSTypeBeforeCast to clear stored type")
	}

	// Node(0).SetJSTypeBeforeCast(nil) must still panic with NPE!
	panicked := false
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		pkg_rhino.Node(0).SetJSTypeBeforeCast(nil)
	}()
	if !panicked {
		t.Fatalf("expected Node(0).SetJSTypeBeforeCast(nil) to panic")
	}

	if ret := n.SetJSDocInfo(nil); ret != n {
		t.Fatalf("SetJSDocInfo(nil) must return receiver Node")
	}
	if nodestore.ClearJSDocCount.Load() != 1 || nodestore.SetJSDocObjCount.Load() != 0 {
		t.Fatalf("expected 1 ClearJSDoc and 0 SetJSDocObj on nil *JSDocInfo, got clear=%d set=%d",
			nodestore.ClearJSDocCount.Load(), nodestore.SetJSDocObjCount.Load())
	}
	doc := &pkg_rhino.JSDocInfo{Summary: "@const"}
	if ret := n.SetJSDocInfo(doc); ret != n || n.GetJSDocInfo() != doc {
		t.Fatalf("live SetJSDocInfo failed")
	}
	n.SetJSDocInfo(nil)
	if n.GetJSDocInfo() != nil {
		t.Fatalf("SetJSDocInfo(nil) failed to clear JSDoc")
	}

	// 3a. GlobalNamespace.ReleaseSlabs_(): exact pool return for refsForNodeMap, nameMap, and nameMapByModule!
	ext, root := buildManyRefAST()
	j2gort.ResetPoolCounters()
	iacp := &InlineAndCollapseProperties{compiler: &Compiler{}}
	iacp.Process_2(ext, root)
	if iacp.Count != 1 {
		t.Fatalf("expected 1 global name in GlobalNamespace, got %d", iacp.Count)
	}
	if j2gort.OrderedMapPtrAcquireCnt.Load() < 2 ||
		j2gort.OrderedMapPtrReleaseCnt.Load() != j2gort.OrderedMapPtrAcquireCnt.Load() {
		t.Fatalf("GlobalNamespace refsForNodeMap + Scope.vars pool leak: acq=%d rel=%d",
			j2gort.OrderedMapPtrAcquireCnt.Load(), j2gort.OrderedMapPtrReleaseCnt.Load())
	}
	if j2gort.OrderedMapValReleaseCnt.Load() < 2 || j2gort.CustomOrderedMapReleaseCnt.Load() < 1 {
		t.Fatalf("expected nameMap, ModuleNameMap.names, and nameMapByModule ReleaseToPool, got ordValRel=%d custRel=%d",
			j2gort.OrderedMapValReleaseCnt.Load(), j2gort.CustomOrderedMapReleaseCnt.Load())
	}

	// 3b. InlineFunctions.Process_2: exact pool return across kept, cost-trimmed, and disallowed FunctionStates!
	j2gort.ResetPoolCounters()
	inf := NewInlineFunctions(&Compiler{})
	inf.Process_2(ext, root)
	if inf.InlinedCount != 36 {
		t.Fatalf("expected 36 references on keepFn, got %d", inf.InlinedCount)
	}
	if j2gort.OrderedMapPtrAcquireCnt.Load() < 4 ||
		j2gort.OrderedMapPtrReleaseCnt.Load() != j2gort.OrderedMapPtrAcquireCnt.Load() {
		t.Fatalf("InlineFunctions fs.references pool leak across kept/trimmed/disallowed: acq=%d rel=%d",
			j2gort.OrderedMapPtrAcquireCnt.Load(), j2gort.OrderedMapPtrReleaseCnt.Load())
	}
	if j2gort.OrderedMapValReleaseCnt.Load() != 2 {
		t.Fatalf("expected both fns and anonFns ReleaseToPool (2), got %d", j2gort.OrderedMapValReleaseCnt.Load())
	}

	// 3c. ReferenceCollection.rest pooled slice: exact 0 live pooled slices after InferConsts,
	//     and live retention on escaping BuildSymbolTable!
	j2gort.ResetPoolCounters()
	ic := &InferConsts{compiler: &Compiler{}}
	ic.Process_2(ext, root)
	if ic.TotalRefs != 36 {
		t.Fatalf("expected 36 collected references in InferConsts, got %d", ic.TotalRefs)
	}
	if j2gort.PooledSliceAcquireCnt.Load() < 6 || j2gort.PooledSliceLiveCnt.Load() != 0 {
		t.Fatalf("expected ReferenceCollection.rest pooled slice acquire >= 6 (initial + growth) and live == 0 after ReleaseSlabs_, got acq=%d rel=%d live=%d",
			j2gort.PooledSliceAcquireCnt.Load(), j2gort.PooledSliceReleaseCnt.Load(), j2gort.PooledSliceLiveCnt.Load())
	}

	st := BuildSymbolTable(&Compiler{}, ext, root)
	if len(st.saved) != 3 || st.saved[0].Size() != 12 || st.saved[0].Get(11).GetNode() == 0 {
		t.Fatalf("escaping SymbolTable ReferenceCollection corrupted: saved=%d", len(st.saved))
	}

	// 3d. CrossChunkReferenceCollector + CrossChunkCodeMotion:
	//     - Direct CrossChunkReferenceCollector caller retains live TopLevelStatements after Process_1!
	directCCRC := NewCrossChunkReferenceCollector(&Compiler{}, NewSyntacticScopeCreator_1(&Compiler{}))
	directCCRC.Process_1(root)
	if directCCRC.GetTopLevelStatementCount() != 3 || directCCRC.GetTopLevelStatement(0).GetDeclaredNameReference().GetSymbol().GetName() != "hotVar" {
		t.Fatalf("CrossChunkReferenceCollector prematurely released TopLevelStatements during Process_1!")
	}
	directCCRC.ReleaseSlabs_()

	//     - CrossChunkCodeMotion.Process_2 reads TopLevelStatements after Process_1 and releases ALL slabs on exit!
	j2gort.ResetPoolCounters()
	cccm := &CrossChunkCodeMotion{compiler: &Compiler{}, chunkCount: 2}
	cccm.Process_2(ext, root)
	if cccm.MovedRefs != 3*(12+11) {
		t.Fatalf("CrossChunkCodeMotion MovedRefs mismatch: got %d, want %d", cccm.MovedRefs, 3*(12+11))
	}
	if j2gort.PooledSliceLiveCnt.Load() != 0 || j2gort.SlabReleaseCount.Load() < 4 ||
		j2gort.OrderedMapPtrReleaseCnt.Load() != j2gort.OrderedMapPtrAcquireCnt.Load() {
		t.Fatalf("CrossChunkCodeMotion failed to release CrossChunkReferenceCollector + ScopeCreator slabs: sliceLive=%d slabRel=%d ptrAcq=%d ptrRel=%d",
			j2gort.PooledSliceLiveCnt.Load(), j2gort.SlabReleaseCount.Load(),
			j2gort.OrderedMapPtrAcquireCnt.Load(), j2gort.OrderedMapPtrReleaseCnt.Load())
	}

	// 4. Concurrent -race stress test across all Wave 28C paths!
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e, r := buildManyRefAST()
			localNode := pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "n")
			gp := pkg_rhino.NewGetProp(pkg_rhino.NewGetProp(pkg_rhino.NewGetProp(pkg_rhino.NewGetProp(localNode, "a"), "b"), "c"), "d")
			for iter := 0; iter < 15; iter++ {
				if gp.GetQualifiedName() != "n.a.b.c.d" {
					panic("bad qname")
				}
				localNode.SetJSTypeBeforeCast(pkg_rhino.NewObjectType("num"))
				localNode.SetJSTypeBeforeCast(ifaceTypedNil)
				localNode.SetJSDocInfo(&pkg_rhino.JSDocInfo{Summary: "doc"})
				localNode.SetJSDocInfo(nil)

				p1 := &InlineAndCollapseProperties{compiler: &Compiler{}}
				p1.Process_2(e, r)
				p2 := NewInlineFunctions(&Compiler{})
				p2.Process_2(e, r)
				p3 := &InferConsts{compiler: &Compiler{}}
				p3.Process_2(e, r)
				p4 := &CrossChunkCodeMotion{compiler: &Compiler{}, chunkCount: 2}
				p4.Process_2(e, r)
			}
		}()
	}
	wg.Wait()
}
`)

	cmd := exec.Command("go", "test", "-race", "-v", "./...")
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lowered Wave 28C runtime test with -race failed: %v\n%s", err, string(out))
	}
}

func TestWave28C_JoinDotFallbackWithoutJoinDot4AndPreLoweredUpgradeIdempotency(t *testing.T) {
	tmpDir := t.TempDir()
	mustWrite := func(rel, content string) {
		t.Helper()
		full := filepath.Join(tmpDir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	mustWrite("go.mod", "module example.com/w28cfallback\n\ngo 1.24\n")
	// nodestore stub WITHOUT JoinDot4 -> 4-part must fold to JoinDot(JoinDot3(a, b, c), d)
	// and 5-part must fold to JoinDot3(JoinDot3(a, b, c), d, e)!
	mustWrite("shims/nodestore/nodestore.go", `package nodestore
type Node uint32
func JoinDot(a, b string) string     { return a + "." + b }
func JoinDot3(a, b, c string) string { return a + "." + b + "." + c }
`)
	mustWrite("gen/rhino/node.go", `package rhino
import "example.com/w28cfallback/shims/nodestore"
type Node nodestore.Node
func Join4(a, b, c, d string) string    { return a + "." + b + "." + c + "." + d }
func Join5(a, b, c, d, e string) string { return a + "." + b + "." + c + "." + d + "." + e }
`)

	allPasses := passes.AllPasses()
	rhinoDir := filepath.Join(tmpDir, "gen/rhino")
	res1, err := passes.RunOnPackageDir(rhinoDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if res1.TotalSites != 2 {
		t.Fatalf("expected 2 JoinDot fallback sites on Run 1, got %+v", res1)
	}
	res2, err := passes.RunOnPackageDir(rhinoDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on Run 2, got %+v", res2)
	}

	b, err := os.ReadFile(filepath.Join(rhinoDir, "node.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Join(strings.Fields(string(b)), " ")
	for _, want := range []string{
		"nodestore.JoinDot(nodestore.JoinDot3(a, b, c), d)",
		"nodestore.JoinDot3(nodestore.JoinDot3(a, b, c), d, e)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected fallback %q in:\n%s", want, string(b))
		}
	}
}

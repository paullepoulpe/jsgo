package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func setupWave20CSynthWorkspace(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	goMod := "module example.com/w20ctest\n\ngo 1.24.0\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	nsDir := filepath.Join(tmpDir, "shims", "nodestore")
	if err := os.MkdirAll(nsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	nsSrc := `package nodestore

import (
	"math"
	"sync"
	"sync/atomic"
)

type Node uint32
type StringID uint32

const NilNode Node = 0

type NodeData struct {
	Parent       Node
	First        Node
	Next         Node
	Previous     Node
	LinenoCharno int32
	Payload      uint32
	PayloadHi    uint32
	PropFlags    uint32
	SourceFile   uint32
	Token        uint8
	KindFlags    uint8
}

type storeState struct {
	mu      sync.Mutex
	nodes   []NodeData
	strings []string
	strIdx  map[string]StringID
	files   []any
}

var currentStore = &storeState{
	nodes:   make([]NodeData, 1),
	strings: []string{""},
	strIdx:  map[string]StringID{"": 0},
	files:   []any{nil},
}

var LocalInternCalls atomic.Int64
var GlobalInternCalls atomic.Int64

func ResetCounters() {
	LocalInternCalls.Store(0)
	GlobalInternCalls.Store(0)
}

func Intern(s string) StringID {
	GlobalInternCalls.Add(1)
	currentStore.mu.Lock()
	defer currentStore.mu.Unlock()
	if id, ok := currentStore.strIdx[s]; ok {
		return id
	}
	id := StringID(len(currentStore.strings))
	currentStore.strings = append(currentStore.strings, s)
	currentStore.strIdx[s] = id
	return id
}

func String(id StringID) string {
	currentStore.mu.Lock()
	defer currentStore.mu.Unlock()
	return currentStore.strings[id]
}

func AllocNodeFastRec() (Node, *NodeData) {
	currentStore.mu.Lock()
	defer currentStore.mu.Unlock()
	id := Node(len(currentStore.nodes))
	currentStore.nodes = append(currentStore.nodes, NodeData{})
	return id, &currentStore.nodes[id]
}

func (n Node) Rec() *NodeData {
	return &currentStore.nodes[n]
}

func (n Node) Token() uint8      { return currentStore.nodes[n].Token }
func (n Node) KindFlags() uint8  { return currentStore.nodes[n].KindFlags }
func (n Node) Payload() uint32   { return currentStore.nodes[n].Payload }
func (n Node) PropFlags() uint32 { return currentStore.nodes[n].PropFlags }
func (n Node) SourceFile() uint32 { return currentStore.nodes[n].SourceFile }
func (n Node) FirstChild() Node  { return currentStore.nodes[n].First }
func (n Node) Next() Node        { return currentStore.nodes[n].Next }
func (n Node) Previous() Node    { return currentStore.nodes[n].Previous }

func (n Node) NumberVal() float64 {
	rec := &currentStore.nodes[n]
	bits := uint64(rec.Payload) | (uint64(rec.PayloadHi) << 32)
	return math.Float64frombits(bits)
}

func (n Node) SetNumberVal(d float64) {
	rec := &currentStore.nodes[n]
	bits := math.Float64bits(d)
	rec.Payload = uint32(bits)
	rec.PayloadHi = uint32(bits >> 32)
}

func (n Node) SetSourceFileObj(sf any) {
	if sf == nil {
		currentStore.nodes[n].SourceFile = 0
		return
	}
	currentStore.mu.Lock()
	defer currentStore.mu.Unlock()
	for i, existing := range currentStore.files {
		if existing == sf {
			currentStore.nodes[n].SourceFile = uint32(i)
			return
		}
	}
	id := uint32(len(currentStore.files))
	currentStore.files = append(currentStore.files, sf)
	currentStore.nodes[n].SourceFile = id
}

func (n Node) ClonePropsFrom(other Node) bool {
	if other == 0 {
		return true
	}
	rec := &currentStore.nodes[n]
	orec := &currentStore.nodes[other]
	rec.SourceFile = orec.SourceFile
	rec.PropFlags = orec.PropFlags
	return true
}

type NodeAllocator struct {
	tmplSourceFile uint32
	tmplPropFlags  uint32
	l1Cache        map[string]StringID
}

func (a *NodeAllocator) InitLocal() {
	if a.l1Cache == nil {
		a.l1Cache = make(map[string]StringID, 16)
	}
}

func (a *NodeAllocator) SetTemplate(tmpl Node) {
	if tmpl == 0 {
		a.tmplSourceFile = 0
		a.tmplPropFlags = 0
		return
	}
	rec := tmpl.Rec()
	a.tmplSourceFile = rec.SourceFile
	a.tmplPropFlags = rec.PropFlags
}

func (a *NodeAllocator) Flush() {}

func (a *NodeAllocator) Intern(s string) StringID {
	LocalInternCalls.Add(1)
	if a.l1Cache != nil {
		if id, ok := a.l1Cache[s]; ok {
			return id
		}
	}
	id := Intern(s)
	if a.l1Cache != nil {
		a.l1Cache[s] = id
	}
	return id
}

func (a *NodeAllocator) allocBase(tok uint8, tmpl Node) (Node, *NodeData) {
	raw, rec := AllocNodeFastRec()
	rec.Token = tok
	rec.LinenoCharno = -1
	if tmpl != 0 {
		trec := tmpl.Rec()
		rec.SourceFile = trec.SourceFile
		rec.PropFlags = trec.PropFlags
	} else {
		rec.SourceFile = a.tmplSourceFile
		rec.PropFlags = a.tmplPropFlags
	}
	return raw, rec
}

func (a *NodeAllocator) AllocNode1(tok uint8, tmpl Node) Node {
	raw, _ := a.allocBase(tok, tmpl)
	return raw
}

func (a *NodeAllocator) AllocNode2(tok uint8, c1 Node, tmpl Node) Node {
	raw, rec := a.allocBase(tok, tmpl)
	rec.First = c1
	r1 := c1.Rec()
	r1.Previous = c1
	r1.Parent = raw
	return raw
}

func (a *NodeAllocator) AllocNode3(tok uint8, c1, c2 Node, tmpl Node) Node {
	raw, rec := a.allocBase(tok, tmpl)
	rec.First = c1
	r1 := c1.Rec()
	r1.Next = c2
	r1.Previous = c2
	r1.Parent = raw
	r2 := c2.Rec()
	r2.Previous = c1
	r2.Parent = raw
	return raw
}

func (a *NodeAllocator) AllocNode4(tok uint8, c1, c2, c3 Node, tmpl Node) Node {
	raw, rec := a.allocBase(tok, tmpl)
	rec.First = c1
	r1 := c1.Rec()
	r1.Next = c2
	r1.Previous = c3
	r1.Parent = raw
	r2 := c2.Rec()
	r2.Next = c3
	r2.Previous = c1
	r2.Parent = raw
	r3 := c3.Rec()
	r3.Previous = c2
	r3.Parent = raw
	return raw
}

func (a *NodeAllocator) AllocStringNode(tok uint8, strID StringID, tmpl Node) Node {
	raw, rec := a.allocBase(tok, tmpl)
	rec.KindFlags = 3
	rec.Payload = uint32(strID)
	return raw
}

func (a *NodeAllocator) AllocNumberNode(tok uint8, d float64, tmpl Node) Node {
	raw, rec := a.allocBase(tok, tmpl)
	rec.KindFlags = 1
	raw.SetNumberVal(d)
	return raw
}
`
	if err := os.WriteFile(filepath.Join(nsDir, "node_store.go"), []byte(nsSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	j2Dir := filepath.Join(tmpDir, "j2gort")
	if err := os.MkdirAll(j2Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	j2Src := `package j2gort

import (
	"fmt"
	"strings"
	"sync/atomic"
)

const NullString = "\x00"

var (
	SlabChunkAllocs         atomic.Int64
	PooledSlabReleaseCalls  atomic.Int64
	OrderedMapAcquireCalls  atomic.Int64
	OrderedMapReleaseCalls  atomic.Int64
	TypedSliceIterCalls     atomic.Int64
)

func ResetTestMetrics() {
	SlabChunkAllocs.Store(0)
	PooledSlabReleaseCalls.Store(0)
	OrderedMapAcquireCalls.Store(0)
	OrderedMapReleaseCalls.Store(0)
	TypedSliceIterCalls.Store(0)
}

func ThrowNPE() { panic("NullPointerException") }

func CheckStateStr(b bool, msg string) {
	if !b {
		panic("IllegalStateException: " + msg)
	}
}

func CheckStateMsg(b bool, msg any) {
	if !b {
		panic(fmt.Sprintf("IllegalStateException: %v", msg))
	}
}

func CheckNotNull[T any](x *T) *T {
	if x == nil {
		ThrowNPE()
	}
	return x
}

func CheckNotNullIface[T any](x T) T {
	if IfaceEquals[T](x, nil) {
		ThrowNPE()
	}
	return x
}

type JavaNilable interface {
	IsJavaNil_() bool
}

func IfaceEquals[T any](a any, b any) bool {
	aNil := a == nil
	if !aNil {
		if jn, ok := a.(JavaNilable); ok && jn.IsJavaNil_() {
			aNil = true
		}
	}
	bNil := b == nil
	if !bNil {
		if jn, ok := b.(JavaNilable); ok && jn.IsJavaNil_() {
			bNil = true
		}
	}
	if aNil && bNil {
		return true
	}
	if aNil != bNil {
		return false
	}
	return a == b
}

func TryCast[T any](x any) (T, bool) {
	var zero T
	if x == nil {
		return zero, false
	}
	if jn, ok := x.(JavaNilable); ok && jn.IsJavaNil_() {
		return zero, false
	}
	v, ok := x.(T)
	return v, ok
}

func Cast[T any](x any) T {
	var zero T
	if x == nil {
		return zero
	}
	if jn, ok := x.(JavaNilable); ok && jn.IsJavaNil_() {
		return zero
	}
	return x.(T)
}

func ToAny[T any](v T) any { return v }

func JStrLen(s string) int32 { return int32(len(s)) }

type Slab[T any] struct {
	buf []T
}

func (s *Slab[T]) Alloc() *T {
	i := len(s.buf)
	if i == cap(s.buf) {
		SlabChunkAllocs.Add(1)
		c := make([]T, 1, 64)
		s.buf = c
		return &c[0]
	}
	b := s.buf[:i+1]
	s.buf = b
	return &b[i]
}

func (s *Slab[T]) Release() {
	if s != nil {
		s.buf = nil
	}
}

type PooledSlab[T any] struct {
	buf []T
}

func (s *PooledSlab[T]) Alloc() *T {
	i := len(s.buf)
	if i == cap(s.buf) {
		SlabChunkAllocs.Add(1)
		c := make([]T, 1, 64)
		s.buf = c
		return &c[0]
	}
	b := s.buf[:i+1]
	s.buf = b
	return &b[i]
}

func (s *PooledSlab[T]) Release() {
	if s != nil {
		clear(s.buf)
		s.buf = s.buf[:0]
		PooledSlabReleaseCalls.Add(1)
	}
}

// --- Iterators & ArrayList ---

type Ext_java_util_Iterator interface {
	IsJavaNil_() bool
	HasNext() bool
	Next() any
	Remove()
}

type Ext_java_lang_Iterable interface {
	IsJavaNil_() bool
	Iterator() Ext_java_util_Iterator
}

type Ext_java_util_Collection interface {
	Ext_java_lang_Iterable
	Size() int32
	IsEmpty() bool
	Clear()
	Add_1(any) bool
}

type Ext_java_util_List interface {
	Ext_java_util_Collection
	Get(int32) any
}

type Ext_java_util_ArrayList interface {
	Ext_java_util_List
}

type Ext_java_util_function_BiConsumer interface {
	IsJavaNil_() bool
	Accept(any, any)
}

type typedSliceIter[E any] struct {
	items []E
	idx   int
}

func NewTypedSliceIter[E any](s []E) Ext_java_util_Iterator {
	TypedSliceIterCalls.Add(1)
	return &typedSliceIter[E]{items: s}
}

func (it *typedSliceIter[E]) IsJavaNil_() bool { return it == nil }
func (it *typedSliceIter[E]) HasNext() bool    { return it.idx < len(it.items) }
func (it *typedSliceIter[E]) Next() any {
	v := it.items[it.idx]
	it.idx++
	return v
}
func (it *typedSliceIter[E]) Remove() {
	panic("UnsupportedOperationException")
}

func NewSliceIter[E any](s []E) Ext_java_util_Iterator {
	return &typedSliceIter[E]{items: s}
}

type IterCursor struct {
	it Ext_java_util_Iterator
}

func NewIterCursor(it Ext_java_util_Iterator) IterCursor { return IterCursor{it: it} }
func (c *IterCursor) HasNext() bool                      { return c.it.HasNext() }
func (c *IterCursor) Next() any                          { return c.it.Next() }

type arrayListImpl struct {
	items []any
}

func NewExt_java_util_ArrayList() Ext_java_util_ArrayList {
	return &arrayListImpl{items: make([]any, 0)}
}
func NewExt_java_util_ArrayList_1_int32(c int32) Ext_java_util_ArrayList {
	return &arrayListImpl{items: make([]any, 0, c)}
}
func MakeNonNilSlice[T any](c int32) []T {
	return make([]T, 0, c)
}
func PtrToAny[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}
func CastPtr[T any](x any) *T {
	if x == nil {
		return nil
	}
	return x.(*T)
}
func (l *arrayListImpl) IsJavaNil_() bool { return l == nil }
func (l *arrayListImpl) Size() int32      { return int32(len(l.items)) }
func (l *arrayListImpl) IsEmpty() bool    { return len(l.items) == 0 }
func (l *arrayListImpl) Get(i int32) any  { return l.items[i] }
func (l *arrayListImpl) Add_1(e any) bool { l.items = append(l.items, e); return true }
func (l *arrayListImpl) Clear()           { clear(l.items); l.items = l.items[:0] }
func (l *arrayListImpl) Iterator() Ext_java_util_Iterator {
	return &arrayListIter{list: l, lastRet: -1}
}

type arrayListIter struct {
	list    *arrayListImpl
	cursor  int
	lastRet int
}

func (it *arrayListIter) IsJavaNil_() bool { return it == nil }
func (it *arrayListIter) HasNext() bool    { return it.cursor < len(it.list.items) }
func (it *arrayListIter) Next() any {
	v := it.list.items[it.cursor]
	it.lastRet = it.cursor
	it.cursor++
	return v
}
func (it *arrayListIter) Remove() {
	if it.lastRet < 0 {
		panic("IllegalStateException")
	}
	idx := it.lastRet
	copy(it.list.items[idx:], it.list.items[idx+1:])
	it.list.items = it.list.items[:len(it.list.items)-1]
	it.cursor = idx
	it.lastRet = -1
}

func FormatSliceJavaList[E any](s []E) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, v := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte(']')
	return b.String()
}

func SliceToJavaList[E any](s []E) Ext_java_util_List {
	if s == nil {
		return nil
	}
	items := make([]any, len(s))
	for i, v := range s {
		items[i] = v
	}
	return &arrayListImpl{items: items}
}

func (l *arrayListImpl) String() string {
	var b strings.Builder
	b.WriteByte('[')
	for i, v := range l.items {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte(']')
	return b.String()
}

type Ext_com_google_common_base_MoreObjects_ToStringHelper struct {
	parts []string
}

func Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(_ any) *Ext_com_google_common_base_MoreObjects_ToStringHelper {
	return &Ext_com_google_common_base_MoreObjects_ToStringHelper{}
}

func (h *Ext_com_google_common_base_MoreObjects_ToStringHelper) Add_2_string_any(k string, v any) *Ext_com_google_common_base_MoreObjects_ToStringHelper {
	h.parts = append(h.parts, fmt.Sprintf("%s=%v", k, v))
	return h
}

func (h *Ext_com_google_common_base_MoreObjects_ToStringHelper) String() string {
	return fmt.Sprint(h.parts)
}

// --- Java Map + OrderedMap + IdentityOrderedMap ---

type Ext_java_util_Map interface {
	IsJavaNil_() bool
	Get(any) any
	Put(any, any) any
	Remove(any) any
	ContainsKey(any) bool
	Clear()
	Size() int32
	IsEmpty() bool
}

type simpleJavaMap struct {
	keys []any
	vals []any
}

func NewLinkedHashMap() Ext_java_util_Map              { return &simpleJavaMap{} }
func NewIdentityHashMap() *Ext_java_util_IdentityHashMap { return &Ext_java_util_IdentityHashMap{} }

func (m *simpleJavaMap) IsJavaNil_() bool { return m == nil }
func (m *simpleJavaMap) Size() int32      { return int32(len(m.keys)) }
func (m *simpleJavaMap) IsEmpty() bool    { return len(m.keys) == 0 }
func (m *simpleJavaMap) Clear()           { m.keys, m.vals = nil, nil }
func (m *simpleJavaMap) ContainsKey(k any) bool {
	for _, ek := range m.keys {
		if ek == k {
			return true
		}
	}
	return false
}
func (m *simpleJavaMap) Get(k any) any {
	for i, ek := range m.keys {
		if ek == k {
			return m.vals[i]
		}
	}
	return nil
}
func (m *simpleJavaMap) Put(k, v any) any {
	for i, ek := range m.keys {
		if ek == k {
			old := m.vals[i]
			m.vals[i] = v
			return old
		}
	}
	m.keys = append(m.keys, k)
	m.vals = append(m.vals, v)
	return nil
}
func (m *simpleJavaMap) Remove(k any) any {
	for i, ek := range m.keys {
		if ek == k {
			old := m.vals[i]
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			m.vals = append(m.vals[:i], m.vals[i+1:]...)
			return old
		}
	}
	return nil
}

type Ext_java_util_IdentityHashMap struct {
	simpleJavaMap
}

type mapEntry[K comparable, V any] struct {
	key     K
	val     V
	deleted bool
}

type OrderedMap[K comparable, V any] struct {
	entries []mapEntry[K, V]
	live    int32
	init    bool
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{init: true}
}

func NewOrderedMapCap[K comparable, V any](capHint int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{entries: make([]mapEntry[K, V], 0, capHint), init: true}
}

func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	OrderedMapAcquireCalls.Add(1)
	return OrderedMap[K, V]{entries: make([]mapEntry[K, V], 0, 8), init: true}
}

func (m *OrderedMap[K, V]) ReleaseToPool() {
	if m == nil || !m.init {
		return
	}
	clear(m.entries)
	m.entries = nil
	m.live = 0
	m.init = false
	OrderedMapReleaseCalls.Add(1)
}

func (m *OrderedMap[K, V]) IsNil() bool   { return m == nil || !m.init }
func (m *OrderedMap[K, V]) Size() int32   { return m.live }
func (m *OrderedMap[K, V]) IsEmpty() bool { return m.live == 0 }
func (m *OrderedMap[K, V]) Clear() {
	clear(m.entries)
	m.entries = m.entries[:0]
	m.live = 0
}
func (m *OrderedMap[K, V]) GetOk(k K) (V, bool) {
	for i := range m.entries {
		if !m.entries[i].deleted && m.entries[i].key == k {
			return m.entries[i].val, true
		}
	}
	var zero V
	return zero, false
}
func (m *OrderedMap[K, V]) Get(k K) V {
	v, _ := m.GetOk(k)
	return v
}
func (m *OrderedMap[K, V]) ContainsKey(k K) bool {
	_, ok := m.GetOk(k)
	return ok
}
func (m *OrderedMap[K, V]) PutOk(k K, v V) (V, bool) {
	for i := range m.entries {
		if !m.entries[i].deleted && m.entries[i].key == k {
			old := m.entries[i].val
			m.entries[i].val = v
			return old, true
		}
	}
	m.entries = append(m.entries, mapEntry[K, V]{key: k, val: v})
	m.live++
	var zero V
	return zero, false
}
func (m *OrderedMap[K, V]) Put(k K, v V) V {
	old, _ := m.PutOk(k, v)
	return old
}
func (m *OrderedMap[K, V]) RemoveOk(k K) (V, bool) {
	for i := range m.entries {
		if !m.entries[i].deleted && m.entries[i].key == k {
			old := m.entries[i].val
			var zk K
			var zv V
			m.entries[i] = mapEntry[K, V]{key: zk, val: zv, deleted: true}
			m.live--
			return old, true
		}
	}
	var zero V
	return zero, false
}
func (m *OrderedMap[K, V]) Remove(k K) V {
	old, _ := m.RemoveOk(k)
	return old
}

type IdentityOrderedMap[K comparable, V any] struct {
	OrderedMap[K, V]
}

func NewIdentityOrderedMap[K comparable, V any]() IdentityOrderedMap[K, V] {
	return IdentityOrderedMap[K, V]{OrderedMap: NewOrderedMap[K, V]()}
}

func NewIdentityOrderedMapCap[K comparable, V any](c int32) IdentityOrderedMap[K, V] {
	return IdentityOrderedMap[K, V]{OrderedMap: NewOrderedMapCap[K, V](c)}
}

type IdentityMapEntryCursor[K comparable, V any] struct {
	entries []mapEntry[K, V]
	idx     int
	k       K
	v       V
}

func (m *IdentityOrderedMap[K, V]) IterEntries() IdentityMapEntryCursor[K, V] {
	return IdentityMapEntryCursor[K, V]{entries: m.entries}
}

func (c *IdentityMapEntryCursor[K, V]) NextEntry() bool {
	for c.idx < len(c.entries) {
		e := c.entries[c.idx]
		c.idx++
		if !e.deleted {
			c.k = e.key
			c.v = e.val
			return true
		}
	}
	return false
}

func (c *IdentityMapEntryCursor[K, V]) Key() K   { return c.k }
func (c *IdentityMapEntryCursor[K, V]) Value() V { return c.v }
`
	if err := os.WriteFile(filepath.Join(j2Dir, "runtime.go"), []byte(j2Src), 0o644); err != nil {
		t.Fatal(err)
	}

	rhinoDir := filepath.Join(tmpDir, "rhino")
	if err := os.MkdirAll(rhinoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	rhinoSrc := `package rhino

import (
	"fmt"
	"math"
	"strings"

	"example.com/w20ctest/j2gort"
	"example.com/w20ctest/shims/nodestore"
)

type Node nodestore.Node
type Token uint8
type InputId struct{ Name string }

func (i *InputId) IsJavaNil_() bool { return i == nil }

func StaticGet_Token_ADD() Token       { return 14 }
func StaticGet_Token_POS() Token       { return 22 }
func StaticGet_Token_NAME() Token       { return 33 }
func StaticGet_Token_NUMBER() Token     { return 34 }
func StaticGet_Token_STRINGLIT() Token  { return 36 }
func StaticGet_Token_STRING_KEY() Token { return 37 }
func StaticGet_Token_HOOK() Token       { return 68 }
func StaticGet_Token_BLOCK() Token     { return 91 }
func StaticGet_Token_SCRIPT() Token    { return 95 }
func StaticGet_Token_BANG() Token      { return 99 }

type StaticSourceFile interface {
	IsJavaNil_() bool
	GetName() string
}

type SimpleSourceFile struct{ Name string }

func (s *SimpleSourceFile) IsJavaNil_() bool { return s == nil }
func (s *SimpleSourceFile) GetName() string  { return s.Name }

func NewNode_1(tok Token) Node {
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(tok) - 1
	rec.LinenoCharno = -1
	return Node(raw)
}

func NewNode_2(tok Token, c1 Node) Node {
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(tok) - 1
	rec.LinenoCharno = -1
	rec.First = nodestore.Node(c1)
	r1 := nodestore.Node(c1).Rec()
	r1.Previous = nodestore.Node(c1)
	r1.Parent = raw
	return Node(raw)
}

func NewNode_3(tok Token, c1, c2 Node) Node {
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(tok) - 1
	rec.LinenoCharno = -1
	rec.First = nodestore.Node(c1)
	r1 := nodestore.Node(c1).Rec()
	r1.Next = nodestore.Node(c2)
	r1.Previous = nodestore.Node(c2)
	r1.Parent = raw
	r2 := nodestore.Node(c2).Rec()
	r2.Previous = nodestore.Node(c1)
	r2.Parent = raw
	return Node(raw)
}

func NewNode_4(tok Token, c1, c2, c3 Node) Node {
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(tok) - 1
	rec.LinenoCharno = -1
	rec.First = nodestore.Node(c1)
	r1 := nodestore.Node(c1).Rec()
	r1.Next = nodestore.Node(c2)
	r1.Previous = nodestore.Node(c3)
	r1.Parent = raw
	r2 := nodestore.Node(c2).Rec()
	r2.Next = nodestore.Node(c3)
	r2.Previous = nodestore.Node(c1)
	r2.Parent = raw
	r3 := nodestore.Node(c3).Rec()
	r3.Previous = nodestore.Node(c2)
	r3.Parent = raw
	return Node(raw)
}

func Node_NewString_2(tok Token, str string) Node {
	if str == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(tok) - 1
	rec.LinenoCharno = -1
	rec.KindFlags = 3
	rec.Payload = uint32(nodestore.Intern(str))
	return Node(raw)
}

func Node_NewString_1(str string) Node {
	return Node_NewString_2(StaticGet_Token_STRINGLIT(), str)
}

func IR_String(s string) Node {
	return Node_NewString_2(StaticGet_Token_STRINGLIT(), s)
}

func IR_Number(d float64) Node {
	if math.IsNaN(d) || d < 0 || (d == 0 && 1/d < 0) {
		j2gort.CheckStateMsg(false, d)
	}
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(StaticGet_Token_NUMBER()) - 1
	rec.LinenoCharno = -1
	rec.KindFlags = 1
	nodestore.Node(raw).SetNumberVal(d)
	return Node(raw)
}

func (n Node) SetStaticSourceFile(sf StaticSourceFile) {
	nodestore.Node(n).SetSourceFileObj(sf)
}

func (n Node) ClonePropsFrom(other Node) Node {
	nodestore.Node(n).ClonePropsFrom(nodestore.Node(other))
	return n
}

func (n Node) SetLinenoCharno(lineno int32, charno int32) Node {
	rec := nodestore.Node(n).Rec()
	rec.LinenoCharno = (lineno << 12) | (charno & 0xfff)
	return n
}

func (n Node) SetLength(length int32) {
	_ = length
}

func (n Node) Summary() string {
	var b strings.Builder
	var walk func(nodestore.Node)
	walk = func(cur nodestore.Node) {
		rec := cur.Rec()
		fmt.Fprintf(&b, "(tok=%d kind=%d lc=%d sf=%d pf=%d", cur.Token(), cur.KindFlags(), rec.LinenoCharno, cur.SourceFile(), cur.PropFlags())
		switch cur.KindFlags() {
		case 1:
			fmt.Fprintf(&b, " num=%v", cur.NumberVal())
		case 3:
			fmt.Fprintf(&b, " str=%s", nodestore.String(nodestore.StringID(cur.Payload())))
		}
		for c := cur.FirstChild(); c != 0; c = c.Next() {
			b.WriteByte(' ')
			walk(c)
		}
		b.WriteByte(')')
	}
	walk(nodestore.Node(n))
	return b.String()
}
`
	if err := os.WriteFile(filepath.Join(rhinoDir, "node.go"), []byte(rhinoSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	return tmpDir
}

func runCmdInDir(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go %s failed: %v\n%s", strings.Join(args, " "), err, string(out))
	}
}

// --- Task 1 + Task 2.4: Preserve defaultScopeCreator_ slabs across releaseNodeTraversal_ & elide zeroed initFields_* ---

func TestWave20C_NodeTraversalPreservesScopeCreatorSlabsAndElidesZeroInitFields(t *testing.T) {
	tmpDir := setupWave20CSynthWorkspace(t)
	jscompDir := filepath.Join(tmpDir, "jscomp")
	if err := os.MkdirAll(jscompDir, 0o755); err != nil {
		t.Fatal(err)
	}

	jscompSrc := `package jscomp

import (
	"example.com/w20ctest/j2gort"
	pkg_rhino "example.com/w20ctest/rhino"
)

type CompilerOptions struct{ marker string }

func (o *CompilerOptions) GetSyntheticBlockStartMarker() string {
	if o == nil || o.marker == "" {
		return j2gort.NullString
	}
	return o.marker
}

type AbstractCompiler interface {
	IsJavaNil_() bool
	GetOptions() *CompilerOptions
}

type CompilerInput struct{ name string }

type AbstractScope interface {
	IsJavaNil_() bool
	GetRootNode_0() pkg_rhino.Node
}

type Var struct {
	name  string
	scope *Scope
	index int32
}

func (v *Var) initFields_Var() {
	this_ := v
	this_.name = ""
	this_.scope = nil
	this_.index = 0
}

func NewVar(name string, scope *Scope, index int32) *Var {
	this_ := &Var{}
	InitVar(this_, name, scope, index)
	return this_
}

func InitVar(this_ *Var, name string, scope *Scope, index int32) {
	this_.initFields_Var()
	this_.name = name
	this_.scope = scope
	this_.index = index
}

type AbstractScopeBase struct {
	var0         *Var
	var1         *Var
	vars         any
	implicitVars any
	zeroDepth    int32
	zeroFlag     bool
}

func (this_ *AbstractScopeBase) initFields_AbstractScope() {
	this_.var0 = nil
	this_.var1 = nil
	this_.vars = nil
	this_.implicitVars = nil
	this_.zeroDepth = 0
	this_.zeroFlag = false
}

func InitAbstractScope(this_ *AbstractScopeBase) {
	this_.initFields_AbstractScope()
}

type Scope struct {
	AbstractScopeBase
	parent   *Scope
	rootNode pkg_rhino.Node
}

func (s *Scope) IsJavaNil_() bool              { return s == nil }
func (s *Scope) GetRootNode_0() pkg_rhino.Node { return s.rootNode }

func (this_ *Scope) initFields_Scope() {
	this_.parent = nil
}

func NewScope_1(rootNode pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	InitScope_1(this_, rootNode)
	return this_
}

func InitScope_1(this_ *Scope, rootNode pkg_rhino.Node) {
	InitAbstractScope(&this_.AbstractScopeBase)
	this_.initFields_Scope()
	this_.parent = nil
	this_.rootNode = rootNode
}

func NewScope_2(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	InitScope_2(this_, parent, rootNode)
	return this_
}

func InitScope_2(this_ *Scope, parent *Scope, rootNode pkg_rhino.Node) {
	InitAbstractScope(&this_.AbstractScopeBase)
	this_.initFields_Scope()
	this_.parent = parent
	this_.rootNode = rootNode
}

func Scope_CreateGlobalScope(rootNode pkg_rhino.Node) *Scope {
	return NewScope_1(rootNode)
}

func Scope_CreateChildScope(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	return NewScope_2(parent, rootNode)
}

func (this_ *Scope) Declare(name string) *Var {
	v := NewVar(name, this_, 0)
	this_.var0 = v
	return v
}

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope_2_Node_AbstractScope(pkg_rhino.Node, AbstractScope) any
}

type SyntacticScopeCreator_RedeclarationHandler interface {
	IsJavaNil_() bool
}

type defaultRedeclHandler struct{}

func (h *defaultRedeclHandler) IsJavaNil_() bool { return h == nil }

var defaultHandler SyntacticScopeCreator_RedeclarationHandler = &defaultRedeclHandler{}

func StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER() SyntacticScopeCreator_RedeclarationHandler {
	return defaultHandler
}

type SyntacticScopeCreator_ScopeScanner struct {
	scope                         *Scope
	compiler                      AbstractCompiler
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler
	treatProvidesAsRedeclarations bool
	inputId                       *pkg_rhino.InputId
}

func (s *SyntacticScopeCreator_ScopeScanner) Reset() {
	s.scope = nil
	s.inputId = nil
}

type SyntacticScopeCreator struct {
	compiler                      AbstractCompiler
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler
	treatProvidesAsRedeclarations bool
	reusableScanner               *SyntacticScopeCreator_ScopeScanner
	scannerInUse                  bool
}

func (s *SyntacticScopeCreator) IsJavaNil_() bool { return s == nil }

func NewSyntacticScopeCreator_1(compiler AbstractCompiler) *SyntacticScopeCreator {
	s := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_1(s, compiler)
	return s
}

func InitSyntacticScopeCreator_1(this_ *SyntacticScopeCreator, compiler AbstractCompiler) {
	this_.compiler = compiler
	this_.redeclarationHandler = StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER()
	this_.treatProvidesAsRedeclarations = false
	this_.reusableScanner = &SyntacticScopeCreator_ScopeScanner{
		compiler:             compiler,
		redeclarationHandler: this_.redeclarationHandler,
	}
}

func (this_ *SyntacticScopeCreator) CreateScope_2_Node_Scope(n pkg_rhino.Node, parent *Scope) *Scope {
	var sc *Scope
	if parent == nil {
		sc = Scope_CreateGlobalScope(n)
	} else {
		sc = Scope_CreateChildScope(parent, n)
	}
	sc.Declare("local")
	if !this_.scannerInUse {
		this_.scannerInUse = true
		this_.reusableScanner.scope = sc
		this_.reusableScanner.Reset()
		this_.scannerInUse = false
	}
	return sc
}

func (this_ *SyntacticScopeCreator) CreateScope_2_Node_AbstractScope(n pkg_rhino.Node, parent AbstractScope) any {
	var p *Scope
	if parent != nil && !parent.IsJavaNil_() {
		p = parent.(*Scope)
	}
	return this_.CreateScope_2_Node_Scope(n, p)
}

type NodeTraversal_Callback interface {
	IsJavaNil_() bool
	ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool
	Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type NodeTraversal_ScopedCallback interface {
	NodeTraversal_Callback
	EnterScope(*NodeTraversal)
	ExitScope(*NodeTraversal)
}

type NodeTraversal struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCallback                                  NodeTraversal_ScopedCallback
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
	mayContainSyntheticBlocks                      bool
	currentNode                                    pkg_rhino.Node
	currentHoistScopeRoot                          pkg_rhino.Node
	currentFunction                                pkg_rhino.Node
	currentScript                                  pkg_rhino.Node
	currentChangeScope                             pkg_rhino.Node
	scopeNodes                                     []pkg_rhino.Node
	scopeObjs                                      []AbstractScope
	scopeDepth                                     int32
	sourceName                                     string
	inputId                                        *pkg_rhino.InputId
	compilerInput                                  *CompilerInput
}

var emptyScopeNodes = make([]pkg_rhino.Node, 0)
var emptyScopeObjs = make([]AbstractScope, 0)

func clinit_NodeTraversal() {}

func (this_ *NodeTraversal) initFields_NodeTraversal() {
	this_.scopeNodes = emptyScopeNodes
	this_.scopeObjs = emptyScopeObjs
	this_.scopeDepth = 0
	this_.sourceName = j2gort.NullString
}

func NewNodeTraversal(builder *NodeTraversal_Builder) *NodeTraversal {
	clinit_NodeTraversal()
	this_ := &NodeTraversal{}
	InitNodeTraversal(this_, builder)
	return this_
}

func InitNodeTraversal(this_ *NodeTraversal, builder *NodeTraversal_Builder) {
	this_.initFields_NodeTraversal()
	this_.compiler = j2gort.CheckNotNullIface(builder.compiler)
	this_.callback = j2gort.CheckNotNullIface(builder.callback)
	scopedCallback, ok_1 := j2gort.TryCast[NodeTraversal_ScopedCallback](this_.callback)
	if ok_1 {
		this_.scopeCallback = scopedCallback
	} else {
		this_.scopeCallback = nil
	}
	if j2gort.IfaceEquals[ScopeCreator](builder.scopeCreator, nil) {
		this_.scopeCreator = NewSyntacticScopeCreator_1(this_.compiler)
	} else {
		this_.scopeCreator = builder.scopeCreator
	}
	this_.obeyDestructuringAndDefaultValueExecutionOrder = builder.obeyDestructuringAndDefaultValueExecutionOrder
	options := this_.compiler.GetOptions()
	this_.mayContainSyntheticBlocks = (options == nil) || (options.GetSyntheticBlockStartMarker() != j2gort.NullString)
}

func NodeTraversal_S_Builder() *NodeTraversal_Builder {
	return &NodeTraversal_Builder{}
}

func (this_ *NodeTraversal) priv_NodeTraversal_ensureScopeCapacity_3840() {
	if this_.scopeDepth == int32(len(this_.scopeNodes)) {
		newLen := int32(8)
		if this_.scopeDepth > 0 {
			newLen = this_.scopeDepth << 1
		}
		ns := make([]pkg_rhino.Node, newLen)
		copy(ns, this_.scopeNodes)
		this_.scopeNodes = ns
		no := make([]AbstractScope, newLen)
		copy(no, this_.scopeObjs)
		this_.scopeObjs = no
	}
}

func (this_ *NodeTraversal) pushScope(n pkg_rhino.Node) {
	this_.priv_NodeTraversal_ensureScopeCapacity_3840()
	this_.scopeNodes[this_.scopeDepth] = n
	this_.scopeObjs[this_.scopeDepth] = nil
	this_.scopeDepth++
	if !(j2gort.IfaceEquals[NodeTraversal_ScopedCallback](this_.scopeCallback, nil)) {
		this_.scopeCallback.EnterScope(this_)
	}
}

func (this_ *NodeTraversal) popScope() {
	if !(j2gort.IfaceEquals[NodeTraversal_ScopedCallback](this_.scopeCallback, nil)) {
		this_.scopeCallback.ExitScope(this_)
	}
	this_.scopeDepth--
	this_.scopeNodes[this_.scopeDepth] = 0
	this_.scopeObjs[this_.scopeDepth] = nil
}

func (this_ *NodeTraversal) GetScopeRoot() pkg_rhino.Node {
	return this_.scopeNodes[this_.scopeDepth-1]
}

func (this_ *NodeTraversal) GetScope() *Scope {
	idx := this_.scopeDepth - 1
	if this_.scopeObjs[idx] == nil {
		var parent AbstractScope
		if idx > 0 {
			parent = this_.scopeObjs[idx-1]
		}
		this_.scopeObjs[idx] = this_.scopeCreator.CreateScope_2_Node_AbstractScope(this_.scopeNodes[idx], parent).(AbstractScope)
	}
	return this_.scopeObjs[idx].(*Scope)
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseBranch_3833(n pkg_rhino.Node, parent pkg_rhino.Node) {
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverse_3817(root pkg_rhino.Node) {
	this_.currentNode = root
	this_.pushScope(root)
	this_.priv_NodeTraversal_traverseBranch_3833(root, 0)
	this_.popScope()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseRoots_3818(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.currentNode = root
	this_.pushScope(externs)
	this_.priv_NodeTraversal_traverseBranch_3833(root, externs)
	this_.popScope()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseWithScope_3820(root pkg_rhino.Node, s AbstractScope) {
	this_.currentNode = root
	this_.priv_NodeTraversal_ensureScopeCapacity_3840()
	this_.scopeNodes[this_.scopeDepth] = s.GetRootNode_0()
	this_.scopeObjs[this_.scopeDepth] = s
	this_.scopeDepth++
	this_.priv_NodeTraversal_traverseBranch_3833(root, 0)
	this_.popScope()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseAtScope_3821(s AbstractScope) {
	this_.priv_NodeTraversal_traverseWithScope_3820(s.GetRootNode_0(), s)
}

func NodeTraversal_Traverse(compiler AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	clinit_NodeTraversal()
	NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).Traverse(root)
}

type NodeTraversal_Builder struct {
	callback                                       NodeTraversal_Callback
	compiler                                       AbstractCompiler
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

func (this_ *NodeTraversal_Builder) SetCompiler(x AbstractCompiler) *NodeTraversal_Builder {
	this_.compiler = x
	return this_
}

func (this_ *NodeTraversal_Builder) SetCallback_1_NodeTraversal_Callback(x NodeTraversal_Callback) *NodeTraversal_Builder {
	this_.callback = x
	return this_
}

func (this_ *NodeTraversal_Builder) Build() *NodeTraversal {
	return NewNodeTraversal(this_)
}

func (this_ *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverse_3817(root)
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "node_traversal.go"), []byte(jscompSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] < 2 || stats1.ByPass["node_traversal_lowering"] < 6 {
		t.Fatalf("unexpected Run 1 stats: %+v", stats1.ByPass)
	}

	// Run 2 MUST be 100% idempotent (0 sites, 0 modified files)!
	stats2, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites & 0 modified files on Run 2, got %+v", stats2)
	}

	loweredBytes, err := os.ReadFile(filepath.Join(jscompDir, "node_traversal.go"))
	if err != nil {
		t.Fatal(err)
	}
	lowered := string(loweredBytes)
	for _, want := range []string{
		"slab_Scope",
		"slab_Var",
		"savedScanner_ := t_.defaultScopeCreator_.reusableScanner",
		"savedScanner_.compiler = nil",
		"savedScanner_.redeclarationHandler = nil",
		"savedScanner_.Reset()",
		"*t_ = NodeTraversal{}",
		"t_.defaultScopeCreator_.reusableScanner = savedScanner_",
		"this_.sourceName = j2gort.NullString",
	} {
		if !strings.Contains(lowered, want) {
			t.Fatalf("missing %q in lowered node_traversal.go:\n%s", want, lowered)
		}
	}
	if strings.Contains(lowered, "savedCreator_") {
		t.Fatalf("releaseNodeTraversal_ must not preserve un-released defaultScopeCreator_ slabs across pool reuse:\n%s", lowered)
	}
	// Redundant zero-value writes in initFields_AbstractScope must be elided!
	if strings.Contains(lowered, "this_.implicitVars = nil") || strings.Contains(lowered, "this_.var0 = nil") {
		t.Fatalf("expected zero-value writes in initFields_AbstractScope to be elided:\n%s", lowered)
	}
	// Reused NodeTraversal struct MUST keep this_.scopeDepth = 0 in initFields_NodeTraversal!
	if !strings.Contains(lowered, "this_.scopeDepth = 0") {
		t.Fatalf("reused NodeTraversal MUST retain this_.scopeDepth = 0 in initFields_NodeTraversal:\n%s", lowered)
	}

	rtTest := `package jscomp

import (
	"testing"

	"example.com/w20ctest/j2gort"
	pkg_rhino "example.com/w20ctest/rhino"
)

type testCompiler struct{ opts *CompilerOptions }

func (c *testCompiler) IsJavaNil_() bool             { return c == nil }
func (c *testCompiler) GetOptions() *CompilerOptions { return c.opts }

type scopeAllocCallback struct{ scopes []*Scope }

func (cb *scopeAllocCallback) IsJavaNil_() bool { return cb == nil }
func (cb *scopeAllocCallback) ShouldTraverse_3(t *NodeTraversal, n, p pkg_rhino.Node) bool {
	return true
}
func (cb *scopeAllocCallback) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) {
	cb.scopes = append(cb.scopes, t.GetScope())
}

func TestDefaultScopeCreatorScannerPreservedAndSlabsDetached(t *testing.T) {
	j2gort.ResetTestMetrics()
	comp := &testCompiler{opts: &CompilerOptions{}}
	cb := &scopeAllocCallback{}

	// Deterministically verify releaseNodeTraversal_ preserves reusableScanner on the exact same
	// *NodeTraversal instance while detaching defaultScopeCreator_ slabs (zero pinned scopes in pool).
	tInst := &NodeTraversal{}
	b := NodeTraversal_Builder{compiler: comp, callback: cb}
	InitNodeTraversal(tInst, &b)
	tInst.priv_NodeTraversal_traverse_3817(1)
	firstScope := cb.scopes[0]
	scanner1 := tInst.defaultScopeCreator_.reusableScanner
	if scanner1 == nil {
		t.Fatalf("expected non-nil reusableScanner")
	}
	releaseNodeTraversal_(tInst)
	_ = nodeTraversalPool_.Get()
	if tInst.defaultScopeCreator_.compiler != nil || tInst.defaultScopeCreator_.redeclarationHandler != nil {
		t.Fatalf("releaseNodeTraversal_ retained compiler/handler references")
	}
	if tInst.defaultScopeCreator_.reusableScanner != scanner1 || scanner1.compiler != nil {
		t.Fatalf("releaseNodeTraversal_ failed to preserve/clear reusableScanner")
	}
	for iter := 0; iter < 50; iter++ {
		cb.scopes = cb.scopes[:0]
		InitNodeTraversal(tInst, &b)
		tInst.priv_NodeTraversal_traverse_3817(2)
		if tInst.defaultScopeCreator_.reusableScanner != scanner1 {
			t.Fatalf("reusableScanner reallocated on iter %d", iter)
		}
		releaseNodeTraversal_(tInst)
		_ = nodeTraversalPool_.Get()
	}
	if firstScope.var0 == nil || firstScope.var0.name != "local" || firstScope.var1 != nil {
		t.Fatalf("escaped firstScope corrupted across traversals: %+v", firstScope)
	}
	if len(cb.scopes) != 1 || cb.scopes[0].var0 == nil || cb.scopes[0].var0.name != "local" || cb.scopes[0].var1 != nil {
		t.Fatalf("unexpected scope state: %+v", cb.scopes)
	}

	// Also exercise full pooled NodeTraversal_Traverse under -race.
	for i := 0; i < 30; i++ {
		cb.scopes = cb.scopes[:0]
		NodeTraversal_Traverse(comp, 2, cb)
	}
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "nt_wave20c_test.go"), []byte(rtTest), 0o644); err != nil {
		t.Fatal(err)
	}
	runCmdInDir(t, tmpDir, "test", "-race", "-v", "./jscomp")
}

// --- Task 2.1: IRFactory.nodeAlloc_.Intern + JsDocInfoParser shared *nodestore.NodeAllocator ---

func TestWave20C_IRFactoryInternAndJsDocInfoParserNodeAllocator(t *testing.T) {
	tmpDir := setupWave20CSynthWorkspace(t)
	origDir := filepath.Join(tmpDir, "origparsing")
	optDir := filepath.Join(tmpDir, "optparsing")
	for _, d := range []string{origDir, optDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	makeParsingFiles := func(pkgName string) map[string]string {
		return map[string]string{
			"ir_factory.go": `package ` + pkgName + `

import (
	"example.com/w20ctest/j2gort"
	pkg_rhino "example.com/w20ctest/rhino"
)

type IRFactory struct {
	sourceFile   pkg_rhino.StaticSourceFile
	templateNode pkg_rhino.Node
}

func (this_ *IRFactory) initFields_IRFactory() { _ = this_ }

func NewIRFactory(sf pkg_rhino.StaticSourceFile) *IRFactory {
	this_ := &IRFactory{}
	InitIRFactory(this_, sf)
	return this_
}

func InitIRFactory(this_ *IRFactory, sf pkg_rhino.StaticSourceFile) {
	this_.initFields_IRFactory()
	this_.sourceFile = sf
	this_.templateNode = this_.createTemplateNode()
}

func (this_ *IRFactory) createTemplateNode() pkg_rhino.Node {
	n := pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT())
	n.SetStaticSourceFile(this_.sourceFile)
	return n
}

func IRFactory_TransformTree(sf pkg_rhino.StaticSourceFile, docText string) pkg_rhino.Node {
	var irFactory *IRFactory = NewIRFactory(sf)
	n1 := irFactory.NewStringNode_1("alpha")
	n2 := irFactory.NewStringNode_2(pkg_rhino.StaticGet_Token_NAME(), "alpha")
	docNode := irFactory.parseJsDoc(docText)
	return irFactory.NewNode_4(pkg_rhino.StaticGet_Token_HOOK(), n1, n2, docNode)
}

func (this_ *IRFactory) parseJsDoc(docText string) pkg_rhino.Node {
	var jsdocParser *JsDocInfoParser = NewJsDocInfoParser(this_.sourceFile, docText)
	return jsdocParser.ParseDocType()
}

func (this_ *IRFactory) NewNode_1(tok pkg_rhino.Token) pkg_rhino.Node {
	return pkg_rhino.NewNode_1(tok).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewNode_4(tok pkg_rhino.Token, c1, c2, c3 pkg_rhino.Node) pkg_rhino.Node {
	return pkg_rhino.NewNode_4(tok, c1, c2, c3).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewStringNode_1(val string) pkg_rhino.Node {
	return pkg_rhino.IR_String(val).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewStringNode_2(tok pkg_rhino.Token, val string) pkg_rhino.Node {
	return pkg_rhino.Node_NewString_2(tok, val).ClonePropsFrom(this_.templateNode)
}

var _ = j2gort.NullString
`,
			"js_doc_info_parser.go": `package ` + pkgName + `

import (
	"example.com/w20ctest/j2gort"
	pkg_rhino "example.com/w20ctest/rhino"
)

type JsDocTokenStream struct {
	lineno int32
	charno int32
}

func (s *JsDocTokenStream) GetLineno() int32 { return s.lineno }
func (s *JsDocTokenStream) GetCharno() int32 { return s.charno }

type JsDocInfoParser struct {
	stream       *JsDocTokenStream
	templateNode pkg_rhino.Node
	docText      string
}

func NewJsDocInfoParser(sf pkg_rhino.StaticSourceFile, docText string) *JsDocInfoParser {
	this_ := &JsDocInfoParser{}
	this_.stream = &JsDocTokenStream{lineno: 7, charno: 3}
	this_.docText = docText
	tmpl := pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT())
	tmpl.SetStaticSourceFile(sf)
	this_.templateNode = tmpl
	return this_
}

func JsDocInfoParser_ParseTypeString(typeStr string) pkg_rhino.Node {
	p := NewJsDocInfoParser(nil, typeStr)
	return p.ParseDocType()
}

func (this_ *JsDocInfoParser) ParseDocType() pkg_rhino.Node {
	a := this_.priv_JsDocInfoParser_newNode_1(pkg_rhino.StaticGet_Token_BANG())
	b := this_.priv_JsDocInfoParser_newStringNode_1(this_.docText)
	c := this_.priv_JsDocInfoParser_parseFieldName_1("propKey")
	w := this_.priv_JsDocInfoParser_wrapNode_4(pkg_rhino.StaticGet_Token_POS(), b, 8, 2)
	return pkg_rhino.NewNode_4(pkg_rhino.StaticGet_Token_BLOCK(), a, w, c).ClonePropsFrom(this_.templateNode)
}

func (this_ *JsDocInfoParser) priv_JsDocInfoParser_newNode_1(type_ pkg_rhino.Token) pkg_rhino.Node {
	return pkg_rhino.NewNode_1(type_).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno()).ClonePropsFrom(this_.templateNode)
}

func (this_ *JsDocInfoParser) priv_JsDocInfoParser_wrapNode_4(type_ pkg_rhino.Token, n pkg_rhino.Node, lineno int32, charno int32) pkg_rhino.Node {
	if n == 0 {
		return 0
	}
	return pkg_rhino.NewNode_2(type_, n).SetLinenoCharno(lineno, charno).ClonePropsFrom(this_.templateNode)
}

func (this_ *JsDocInfoParser) priv_JsDocInfoParser_newStringNode_1(s string) pkg_rhino.Node {
	var n pkg_rhino.Node = pkg_rhino.Node_NewString_1(s).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno()).ClonePropsFrom(this_.templateNode)
	n.SetLength(j2gort.JStrLen(s))
	return n
}

func (this_ *JsDocInfoParser) priv_JsDocInfoParser_parseFieldName_1(fieldText string) pkg_rhino.Node {
	var string_ string = fieldText
	if len(string_) == 0 {
		return 0
	}
	var n pkg_rhino.Node = pkg_rhino.Node_NewString_2(pkg_rhino.StaticGet_Token_STRING_KEY(), string_).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno()).ClonePropsFrom(this_.templateNode)
	return n
}

var _ = j2gort.NullString
`,
		}
	}

	for fn, src := range makeParsingFiles("origparsing") {
		if err := os.WriteFile(filepath.Join(origDir, fn), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for fn, src := range makeParsingFiles("optparsing") {
		if err := os.WriteFile(filepath.Join(optDir, fn), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on optparsing failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] < 9 {
		t.Fatalf("expected >= 9 slab_alloc_lowering sites on optparsing, got %+v", stats1.ByPass)
	}
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on optparsing failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites & 0 modified files on Run 2, got %+v", stats2)
	}

	irBytes, _ := os.ReadFile(filepath.Join(optDir, "ir_factory.go"))
	jsdocBytes, _ := os.ReadFile(filepath.Join(optDir, "js_doc_info_parser.go"))
	if !strings.Contains(string(irBytes), "this_.nodeAlloc_.Intern(") ||
		!strings.Contains(string(irBytes), "jsdocParser.nodeAlloc_ = &this_.nodeAlloc_") {
		t.Fatalf("missing nodeAlloc_.Intern or jsdocParser.nodeAlloc_ wiring in ir_factory.go:\n%s", string(irBytes))
	}
	if !strings.Contains(string(jsdocBytes), "nodeAlloc_") ||
		!strings.Contains(string(jsdocBytes), "this_.nodeAlloc_.AllocStringNode(") ||
		!strings.Contains(string(jsdocBytes), "this_.nodeAlloc_.Intern(") {
		t.Fatalf("missing JsDocInfoParser nodeAlloc_ lowering in js_doc_info_parser.go:\n%s", string(jsdocBytes))
	}

	parityTest := `package parsing_parity_test

import (
	"testing"

	"example.com/w20ctest/optparsing"
	"example.com/w20ctest/origparsing"
	pkg_rhino "example.com/w20ctest/rhino"
	"example.com/w20ctest/shims/nodestore"
)

func TestIRFactoryAndJsDocNodeAllocatorParity(t *testing.T) {
	sf1 := &pkg_rhino.SimpleSourceFile{Name: "a.js"}
	origTree := origparsing.IRFactory_TransformTree(sf1, "string")

	nodestore.ResetCounters()
	optTree := optparsing.IRFactory_TransformTree(sf1, "string")

	// Ignore sf integer ID difference by checking standalone + L1 interning hit!
	if nodestore.LocalInternCalls.Load() != 4 {
		t.Fatalf("expected 4 LocalInternCalls (2 in IRFactory + 2 in JsDocInfoParser), got %d", nodestore.LocalInternCalls.Load())
	}
	// "alpha" was interned twice on the same IRFactory -> L1 cache should hit on the 2nd "alpha",
	// so GlobalInternCalls is only 3 ("alpha" + "string" + "propKey") instead of 4!
	if nodestore.GlobalInternCalls.Load() != 3 {
		t.Fatalf("expected 3 GlobalInternCalls due to NodeAllocator.Intern L1 cache hit, got %d", nodestore.GlobalInternCalls.Load())
	}
	if origTree.Summary() != optTree.Summary() {
		t.Fatalf("IRFactory + JsDocInfoParser AST parity mismatch:\norig: %s\nopt:  %s", origTree.Summary(), optTree.Summary())
	}

	// Standalone JsDocInfoParser_ParseTypeString (where nodeAlloc_ == nil) must fall back cleanly!
	sOrig := origparsing.JsDocInfoParser_ParseTypeString("number").Summary()
	sOpt := optparsing.JsDocInfoParser_ParseTypeString("number").Summary()
	if sOrig != sOpt {
		t.Fatalf("standalone ParseTypeString mismatch:\norig: %s\nopt:  %s", sOrig, sOpt)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "parsing_parity_test.go"), []byte(parityTest), 0o644); err != nil {
		t.Fatal(err)
	}
	runCmdInDir(t, tmpDir, "test", "-v", "-run", "TestIRFactoryAndJsDocNodeAllocatorParity", ".")
}

// --- Task 2.2 + Task 2.3: LinkedDirectedGraphBase.nodes Pool + RemoveUnusedCode PooledSlabs ---

func TestWave20C_LinkedDirectedGraphNodesPoolAndRemoveUnusedCodeSlabs(t *testing.T) {
	tmpDir := setupWave20CSynthWorkspace(t)
	graphDir := filepath.Join(tmpDir, "graph")
	rucDir := filepath.Join(tmpDir, "jscompruc")
	for _, d := range []string{graphDir, rucDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	graphSrc := `package graph

import "example.com/w20ctest/j2gort"

type Graph[N any, E any] interface {
	Implements_graph_Graph()
	IsJavaNil_() bool
}

type LinkedDirectedGraph_LinkedDiGraphNode[N any, E any] struct{ val N }

func NewLinkedDirectedGraph_LinkedDiGraphNode[N any, E any](v N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	this_ := &LinkedDirectedGraph_LinkedDiGraphNode[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphNode(this_, v)
	return this_
}

func InitLinkedDirectedGraph_LinkedDiGraphNode[N any, E any](this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E], v N) {
	this_.val = v
}

type LinkedDirectedGraphBase[N comparable, E any] struct {
	nodes j2gort.OrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]
}

func (*LinkedDirectedGraphBase[N, E]) Implements_graph_Graph() {}
func (this_ *LinkedDirectedGraphBase[N, E]) IsJavaNil_() bool  { return this_ == nil }

func (this_ *LinkedDirectedGraphBase[N, E]) initFields_LinkedDirectedGraph() {
	this_.nodes = j2gort.NewOrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]()
}

func InitLinkedDirectedGraph[N comparable, E any](this_ *LinkedDirectedGraphBase[N, E]) {
	this_.initFields_LinkedDirectedGraph()
}

func (this_ *LinkedDirectedGraphBase[N, E]) CreateNode(v N) any {
	if existing := this_.nodes.Get(v); existing != nil {
		return existing
	}
	n := NewLinkedDirectedGraph_LinkedDiGraphNode[N, E](v)
	this_.nodes.Put(v, n)
	return n
}

func (this_ *LinkedDirectedGraphBase[N, E]) NodeCount() int32 {
	return this_.nodes.Size()
}
`
	rucSrc := `package jscompruc

import (
	"example.com/w20ctest/j2gort"
	pkg_rhino "example.com/w20ctest/rhino"
)

type RemoveUnusedCode_VarInfo interface {
	IsJavaNil_() bool
	MarkReferenced() bool
}

type RemoveUnusedCode_CanonicalUnremovableVarInfo struct{}

func (v *RemoveUnusedCode_CanonicalUnremovableVarInfo) IsJavaNil_() bool     { return v == nil }
func (v *RemoveUnusedCode_CanonicalUnremovableVarInfo) MarkReferenced() bool { return false }

func NewRemoveUnusedCode_CanonicalUnremovableVarInfo() *RemoveUnusedCode_CanonicalUnremovableVarInfo {
	return &RemoveUnusedCode_CanonicalUnremovableVarInfo{}
}

var canonicalSingleton = NewRemoveUnusedCode_CanonicalUnremovableVarInfo()

type RemoveUnusedCode_Assign struct {
	outer_     *RemoveUnusedCode
	assignNode pkg_rhino.Node
}

func NewRemoveUnusedCode_Assign(outer_ *RemoveUnusedCode, n pkg_rhino.Node) *RemoveUnusedCode_Assign {
	this_ := &RemoveUnusedCode_Assign{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_Assign(this_, n)
	return this_
}

func InitRemoveUnusedCode_Assign(this_ *RemoveUnusedCode_Assign, n pkg_rhino.Node) {
	this_.assignNode = n
}

type RemoveUnusedCode_RealVarInfo struct {
	outer_     *RemoveUnusedCode
	referenced bool
	lastAssign *RemoveUnusedCode_Assign
}

func (v *RemoveUnusedCode_RealVarInfo) IsJavaNil_() bool { return v == nil }
func (v *RemoveUnusedCode_RealVarInfo) MarkReferenced() bool {
	if v.referenced {
		return false
	}
	v.referenced = true
	return true
}

func NewRemoveUnusedCode_RealVarInfo(outer_ *RemoveUnusedCode) *RemoveUnusedCode_RealVarInfo {
	this_ := &RemoveUnusedCode_RealVarInfo{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_RealVarInfo(this_)
	return this_
}

func InitRemoveUnusedCode_RealVarInfo(this_ *RemoveUnusedCode_RealVarInfo) {
	this_.referenced = false
}

type RemoveUnusedCode_RemovableBuilder struct {
	outer_ *RemoveUnusedCode
	count  int32
}

func NewRemoveUnusedCode_RemovableBuilder(outer_ *RemoveUnusedCode, c int32) *RemoveUnusedCode_RemovableBuilder {
	this_ := &RemoveUnusedCode_RemovableBuilder{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_RemovableBuilder(this_, c)
	return this_
}

func InitRemoveUnusedCode_RemovableBuilder(this_ *RemoveUnusedCode_RemovableBuilder, c int32) {
	this_.count = c
}

type RemoveUnusedCode struct {
	panicOnEnd bool
	workCount  int32
}

func (this_ *RemoveUnusedCode) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	var canonical RemoveUnusedCode_VarInfo = canonicalSingleton
	_ = canonical
	for i := int32(0); i < 10; i++ {
		vi := NewRemoveUnusedCode_RealVarInfo(this_)
		vi.lastAssign = NewRemoveUnusedCode_Assign(this_, root)
		rb := NewRemoveUnusedCode_RemovableBuilder(this_, i)
		if vi.MarkReferenced() {
			this_.workCount += rb.count + 1
		}
	}
	if this_.panicOnEnd {
		panic("ruc_boom")
	}
}

var _ = j2gort.NullString
`
	if err := os.WriteFile(filepath.Join(graphDir, "linked_directed_graph.go"), []byte(graphSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(rucDir, "remove_unused_code.go"), []byte(rucSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.Run([]string{graphDir, rucDir}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] < 6 {
		t.Fatalf("expected >= 6 slab_alloc_lowering sites across graph+jscompruc, got %+v", stats1.ByPass)
	}
	stats2, err := passes.Run([]string{graphDir, rucDir}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites & 0 modified files on Run 2, got %+v", stats2)
	}

	gotGraphBytes, _ := os.ReadFile(filepath.Join(graphDir, "linked_directed_graph.go"))
	gotGraph := string(gotGraphBytes)
	for _, want := range []string{
		"j2gort.AcquireOrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]()",
		"this_.nodes.ReleaseToPool()",
		"func (this_ *LinkedDirectedGraphBase[N, E]) ReleaseSlabs_()",
	} {
		if !strings.Contains(gotGraph, want) {
			t.Fatalf("missing %q in lowered linked_directed_graph.go:\n%s", want, gotGraph)
		}
	}

	gotRucBytes, _ := os.ReadFile(filepath.Join(rucDir, "remove_unused_code.go"))
	gotRuc := string(gotRucBytes)
	for _, want := range []string{
		"slab_RemoveUnusedCode_RealVarInfo",
		"slab_RemoveUnusedCode_Assign",
		"slab_RemoveUnusedCode_RemovableBuilder",
		"func (this_ *RemoveUnusedCode) ReleaseSlabs_()",
		"defer this_.ReleaseSlabs_()",
	} {
		if !strings.Contains(gotRuc, want) {
			t.Fatalf("missing %q in lowered remove_unused_code.go:\n%s", want, gotRuc)
		}
	}
	if strings.Contains(gotRuc, "slab_RemoveUnusedCode_CanonicalUnremovableVarInfo") {
		t.Fatalf("singleton CanonicalUnremovableVarInfo must NOT be slab-allocated:\n%s", gotRuc)
	}

	rtTest := `package graph_ruc_test

import (
	"testing"

	"example.com/w20ctest/graph"
	"example.com/w20ctest/j2gort"
	"example.com/w20ctest/jscompruc"
)

func TestGraphNodesPoolAndRemoveUnusedCodeSlabs(t *testing.T) {
	j2gort.ResetTestMetrics()

	var g graph.LinkedDirectedGraphBase[string, int32]
	graph.InitLinkedDirectedGraph(&g)
	g.CreateNode("a")
	g.CreateNode("b")
	if g.NodeCount() != 2 {
		t.Fatalf("NodeCount = %d, want 2", g.NodeCount())
	}
	g.ReleaseSlabs_()
	if j2gort.OrderedMapAcquireCalls.Load() != 1 || j2gort.OrderedMapReleaseCalls.Load() != 1 {
		t.Fatalf("expected 1 AcquireOrderedMap and 1 ReleaseToPool, got acquire=%d release=%d",
			j2gort.OrderedMapAcquireCalls.Load(), j2gort.OrderedMapReleaseCalls.Load())
	}

	j2gort.ResetTestMetrics()
	ruc := &jscompruc.RemoveUnusedCode{}
	ruc.Process_2(0, 1)
	if j2gort.PooledSlabReleaseCalls.Load() != 3 {
		t.Fatalf("expected 3 PooledSlab releases from RemoveUnusedCode.Process_2, got %d", j2gort.PooledSlabReleaseCalls.Load())
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "graph_ruc_test.go"), []byte(rtTest), 0o644); err != nil {
		t.Fatal(err)
	}
	runCmdInDir(t, tmpDir, "test", "-v", "-run", "TestGraphNodesPoolAndRemoveUnusedCodeSlabs", ".")
}

// --- Task 3: Unlock .Iterator() on ArrayList<E> fields (including ReferenceCollection.rest) ---

func TestWave20C_NativeSliceIteratorAndReferenceCollection(t *testing.T) {
	tmpDir := setupWave20CSynthWorkspace(t)
	rcDir := filepath.Join(tmpDir, "refpkg")
	if err := os.MkdirAll(rcDir, 0o755); err != nil {
		t.Fatal(err)
	}

	rcSrc := `package refpkg

import "example.com/w20ctest/j2gort"

type Reference struct {
	Name string
}

func (r *Reference) String() string { return r.Name }

type ReferenceCollection struct {
	first    *Reference
	second   *Reference
	rest     j2gort.Ext_java_util_List
	tempScan j2gort.Ext_java_util_ArrayList
	mutItems j2gort.Ext_java_util_ArrayList
}

func (this_ *ReferenceCollection) initFields_ReferenceCollection() {
	this_.rest = nil
	this_.tempScan = j2gort.NewExt_java_util_ArrayList()
	this_.mutItems = j2gort.NewExt_java_util_ArrayList()
}

func NewReferenceCollection() *ReferenceCollection {
	this_ := &ReferenceCollection{}
	InitReferenceCollection(this_)
	return this_
}

func InitReferenceCollection(this_ *ReferenceCollection) {
	this_.initFields_ReferenceCollection()
}

func (this_ *ReferenceCollection) IsJavaNil_() bool { return this_ == nil }

func (this_ *ReferenceCollection) Add(r *Reference) {
	if this_.first == nil {
		this_.first = r
	} else if this_.second == nil {
		this_.second = r
	} else {
		if this_.rest == nil {
			this_.rest = j2gort.NewExt_java_util_ArrayList_1_int32(4)
		}
		this_.rest.Add_1(j2gort.PtrToAny(r))
	}
	this_.tempScan.Add_1(j2gort.PtrToAny(r))
	this_.mutItems.Add_1(j2gort.PtrToAny(r))
}

func (this_ *ReferenceCollection) Iterator() j2gort.Ext_java_util_Iterator {
	if this_.rest != nil {
		return this_.rest.Iterator()
	}
	return j2gort.NewTypedSliceIter([]*Reference(nil))
}

func (this_ *ReferenceCollection) ConcatTempNames() string {
	out := ""
	cur := j2gort.NewIterCursor(this_.tempScan.Iterator())
	for cur.HasNext() {
		raw := cur.Next()
		r := j2gort.CastPtr[Reference](raw)
		out += r.Name
	}
	return out
}

func (this_ *ReferenceCollection) RemoveFirstMutItem() int32 {
	it := this_.mutItems.Iterator()
	if it.HasNext() {
		_ = it.Next()
		it.Remove()
	}
	return this_.mutItems.Size()
}

func (this_ *ReferenceCollection) String() string {
	var refsView any = nil
	if this_.rest != nil {
		refsView = this_.rest
	}
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_any("references", refsView).String()
}
`
	if err := os.WriteFile(filepath.Join(rcDir, "reference_collection.go"), []byte(rcSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(rcDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on refpkg failed: %v", err)
	}
	if stats1.ByPass["native_slice_lowering"] < 2 {
		t.Fatalf("expected >= 2 native_slice_lowering sites (rest + tempScan), got %+v", stats1.ByPass)
	}
	stats2, err := passes.RunOnPackageDir(rcDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on refpkg failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites & 0 modified files on Run 2, got %+v", stats2)
	}

	gotBytes, _ := os.ReadFile(filepath.Join(rcDir, "reference_collection.go"))
	got := string(gotBytes)
	for _, want := range []string{
		"rest     []*Reference",
		"tempScan []*Reference",
		"mutItems j2gort.Ext_java_util_ArrayList",
		"j2gort.NewTypedSliceIter(this_.rest)",
		"j2gort.SliceToJavaList(this_.rest)",
		"range this_.tempScan",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in lowered reference_collection.go:\n%s", want, got)
		}
	}

	rtTest := `package refpkg_test

import (
	"testing"

	"example.com/w20ctest/j2gort"
	"example.com/w20ctest/refpkg"
)

func TestReferenceCollectionSliceIteratorRuntime(t *testing.T) {
	j2gort.ResetTestMetrics()
	rc := refpkg.NewReferenceCollection()
	rc.Add(&refpkg.Reference{Name: "r1"})
	rc.Add(&refpkg.Reference{Name: "r2"})
	rc.Add(&refpkg.Reference{Name: "r3"})
	rc.Add(&refpkg.Reference{Name: "r4"})

	var restNames []string
	for it := rc.Iterator(); it.HasNext(); {
		restNames = append(restNames, it.Next().(*refpkg.Reference).Name)
	}
	if len(restNames) != 2 || restNames[0] != "r3" || restNames[1] != "r4" {
		t.Fatalf("unexpected Iterator() elements: %v", restNames)
	}
	if j2gort.TypedSliceIterCalls.Load() != 1 {
		t.Fatalf("expected NewTypedSliceIter to be called once, got %d", j2gort.TypedSliceIterCalls.Load())
	}
	if got := rc.ConcatTempNames(); got != "r1r2r3r4" {
		t.Fatalf("ConcatTempNames = %q, want r1r2r3r4", got)
	}
	if got := rc.String(); got != "[references=[r3, r4]]" {
		t.Fatalf("String() = %q, want [references=[r3, r4]]", got)
	}
	// mutItems uses it.Remove(), so it stayed *Ext_java_util_ArrayList and supports Remove()!
	if rem := rc.RemoveFirstMutItem(); rem != 3 {
		t.Fatalf("RemoveFirstMutItem = %d, want 3", rem)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "refpkg_test.go"), []byte(rtTest), 0o644); err != nil {
		t.Fatal(err)
	}
	runCmdInDir(t, tmpDir, "test", "-v", "-run", "TestReferenceCollectionSliceIteratorRuntime", ".")
}

// --- Task 4: Unbox base.LinkedIdentityHashMap[K, V] onto j2gort.IdentityOrderedMap[K, V] ---

func TestWave20C_UnboxLinkedIdentityHashMap(t *testing.T) {
	tmpDir := setupWave20CSynthWorkspace(t)
	baseDir := filepath.Join(tmpDir, "base")
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		t.Fatal(err)
	}

	baseSrc := `package base

import (
	"sync/atomic"

	"example.com/w20ctest/j2gort"
)

var EntryAllocs atomic.Int64

type LinkedIdentityHashMap_Entry[K comparable, V any] struct {
	key   K
	value V
	prev  *LinkedIdentityHashMap_Entry[K, V]
	next  *LinkedIdentityHashMap_Entry[K, V]
}

func NewLinkedIdentityHashMap_Entry[K comparable, V any](key K, value V) *LinkedIdentityHashMap_Entry[K, V] {
	EntryAllocs.Add(1)
	return &LinkedIdentityHashMap_Entry[K, V]{key: key, value: value}
}

type LinkedIdentityHashMap[K comparable, V any] struct {
	innerMap j2gort.IdentityOrderedMap[K, *LinkedIdentityHashMap_Entry[K, V]]
	head     *LinkedIdentityHashMap_Entry[K, V]
	tail     *LinkedIdentityHashMap_Entry[K, V]
}

func (this_ *LinkedIdentityHashMap[K, V]) IsJavaNil_() bool { return this_ == nil }

func (this_ *LinkedIdentityHashMap[K, V]) initFields_LinkedIdentityHashMap() {
	this_.innerMap = j2gort.NewIdentityOrderedMap[K, *LinkedIdentityHashMap_Entry[K, V]]()
}

func NewLinkedIdentityHashMap[K comparable, V any]() *LinkedIdentityHashMap[K, V] {
	this_ := &LinkedIdentityHashMap[K, V]{}
	InitLinkedIdentityHashMap(this_)
	return this_
}

func InitLinkedIdentityHashMap[K comparable, V any](this_ *LinkedIdentityHashMap[K, V]) {
	this_.initFields_LinkedIdentityHashMap()
}

func (this_ *LinkedIdentityHashMap[K, V]) Get(key K) V {
	entry := this_.innerMap.Get(key)
	if entry == nil {
		if _, ok := any((*V)(nil)).(*string); ok {
			return any(j2gort.NullString).(V)
		}
		var zero V
		return zero
	}
	return entry.value
}

func (this_ *LinkedIdentityHashMap[K, V]) Put(key K, newValue V) V {
	if entry := this_.innerMap.Get(key); entry != nil {
		old := entry.value
		entry.value = newValue
		return old
	}
	entry := NewLinkedIdentityHashMap_Entry(key, newValue)
	this_.innerMap.Put(key, entry)
	if this_.head == nil {
		this_.head = entry
	} else {
		this_.tail.next = entry
		entry.prev = this_.tail
	}
	this_.tail = entry
	if _, ok := any((*V)(nil)).(*string); ok {
		return any(j2gort.NullString).(V)
	}
	var zero V
	return zero
}

func (this_ *LinkedIdentityHashMap[K, V]) Remove(key K) V {
	entry := this_.innerMap.Remove(key)
	if entry == nil {
		if _, ok := any((*V)(nil)).(*string); ok {
			return any(j2gort.NullString).(V)
		}
		var zero V
		return zero
	}
	if entry.prev != nil {
		entry.prev.next = entry.next
	} else {
		this_.head = entry.next
	}
	if entry.next != nil {
		entry.next.prev = entry.prev
	} else {
		this_.tail = entry.prev
	}
	return entry.value
}

func (this_ *LinkedIdentityHashMap[K, V]) ForEach(consumer j2gort.Ext_java_util_function_BiConsumer) {
	for entry := this_.head; entry != nil; entry = entry.next {
		consumer.Accept(j2gort.ToAny(entry.key), j2gort.ToAny(entry.value))
	}
}

func (this_ *LinkedIdentityHashMap[K, V]) ComputeIfAbsent(key K, fn func(K) V) V {
	v := this_.Get(key)
	if j2gort.IfaceEquals[V](v, nil) {
		nv := fn(key)
		if !j2gort.IfaceEquals[V](nv, nil) {
			this_.Put(key, nv)
			return nv
		}
	}
	return v
}

func (this_ *LinkedIdentityHashMap[K, V]) Clear() {
	this_.innerMap.Clear()
	this_.head = nil
	this_.tail = nil
}

func (this_ *LinkedIdentityHashMap[K, V]) ContainsKey(key K) bool { return this_.innerMap.ContainsKey(key) }
func (this_ *LinkedIdentityHashMap[K, V]) Size() int32            { return this_.innerMap.Size() }
func (this_ *LinkedIdentityHashMap[K, V]) IsEmpty() bool          { return this_.innerMap.IsEmpty() }
`
	if err := os.WriteFile(filepath.Join(baseDir, "linked_identity_hash_map.go"), []byte(baseSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(baseDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on base failed: %v", err)
	}
	if stats1.ByPass["j2gort_peephole"] < 1 {
		t.Fatalf("expected j2gort_peephole to unbox LinkedIdentityHashMap, got %+v", stats1.ByPass)
	}
	stats2, err := passes.RunOnPackageDir(baseDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on base failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites & 0 modified files on Run 2, got %+v", stats2)
	}

	gotBytes, _ := os.ReadFile(filepath.Join(baseDir, "linked_identity_hash_map.go"))
	got := string(gotBytes)
	if !strings.Contains(got, "innerMap j2gort.IdentityOrderedMap[K, V]") {
		t.Fatalf("expected unboxed innerMap j2gort.IdentityOrderedMap[K, V], got:\n%s", got)
	}
	if strings.Contains(got, "this_.head") || strings.Contains(got, "this_.tail") {
		t.Fatalf("expected head/tail linked list fields to be eliminated:\n%s", got)
	}

	rtTest := `package base_test

import (
	"fmt"
	"testing"

	"example.com/w20ctest/base"
	"example.com/w20ctest/j2gort"
)

type keyObj struct{ label string }

type recordingBiConsumer struct{ items []string }

func (c *recordingBiConsumer) IsJavaNil_() bool { return c == nil }
func (c *recordingBiConsumer) Accept(k, v any) {
	c.items = append(c.items, fmt.Sprintf("%s:%v", k.(*keyObj).label, v))
}

func TestUnboxedLinkedIdentityHashMapSemantics(t *testing.T) {
	base.EntryAllocs.Store(0)
	k1 := &keyObj{label: "a"}
	k1Clone := &keyObj{label: "a"} // distinct pointer with equal contents -> identity inequality!
	k2 := &keyObj{label: "b"}
	k3 := &keyObj{label: "c"}

	// 1. String values (must return NullString on miss!)
	ms := base.NewLinkedIdentityHashMap[*keyObj, string]()
	if got := ms.Get(k1); got != j2gort.NullString {
		t.Fatalf("missing string Get = %q, want NullString", got)
	}
	if old := ms.Put(k1, "v1"); old != j2gort.NullString {
		t.Fatalf("fresh string Put old = %q, want NullString", old)
	}
	if ms.ContainsKey(k1Clone) {
		t.Fatalf("Identity map must compare pointer identity, not struct field equality")
	}
	ms.Put(k2, "v2")
	ms.Put(k3, "v3")
	// In-place update preserves insertion order [a, b, c]!
	if old := ms.Put(k1, "v1_updated"); old != "v1" {
		t.Fatalf("update Put old = %q, want v1", old)
	}
	// Remove k2 and reinsert k2 -> order becomes [a, c, b]!
	if rem := ms.Remove(k2); rem != "v2" {
		t.Fatalf("Remove(k2) = %q, want v2", rem)
	}
	ms.Put(k2, "v2_tail")

	rec := &recordingBiConsumer{}
	ms.ForEach(rec)
	if fmt.Sprint(rec.items) != "[a:v1_updated c:v3 b:v2_tail]" {
		t.Fatalf("ForEach order = %v, want [a:v1_updated c:v3 b:v2_tail]", rec.items)
	}

	// 2. LinkedIdentityHashSet pattern: Put(k, nil) stores a nil value while keeping key present!
	sentinel := new(int)
	mset := base.NewLinkedIdentityHashMap[*keyObj, *int]()
	mset.Put(k1, sentinel)
	if prev := mset.Put(k1, nil); prev != sentinel {
		t.Fatalf("LinkedIdentityHashSet Remove pattern Put(k1, nil) prev = %v, want sentinel", prev)
	}
	if !mset.ContainsKey(k1) || mset.Get(k1) != nil {
		t.Fatalf("Put(k1, nil) must keep k1 in map with nil value")
	}
	if prev := mset.Put(k1, nil); prev != nil {
		t.Fatalf("second Put(k1, nil) must return nil")
	}

	// Zero heap Entry allocations!
	if got := base.EntryAllocs.Load(); got != 0 {
		t.Fatalf("EntryAllocs = %d, want 0", got)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "base_unbox_test.go"), []byte(rtTest), 0o644); err != nil {
		t.Fatal(err)
	}
	runCmdInDir(t, tmpDir, "test", "-v", "-run", "TestUnboxedLinkedIdentityHashMapSemantics", ".")
}

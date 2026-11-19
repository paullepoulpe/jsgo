package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave22C_BlockCachedTraversal_TieredRefMap_InlineVariablesAndGuardedContextPooling(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := "module example.com/w22ctest\n\ngo 1.24.0\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	nsDir := filepath.Join(tmpDir, "shims", "nodestore")
	j2Dir := filepath.Join(tmpDir, "gen", "j2gort")
	rhinoDir := filepath.Join(tmpDir, "gen", "rhino")
	jscompDir := filepath.Join(tmpDir, "gen", "jscomp")
	parsingDir := filepath.Join(tmpDir, "gen", "parsing")
	unsafeParsingDir := filepath.Join(tmpDir, "gen", "unsafeparsing")
	for _, d := range []string{nsDir, j2Dir, rhinoDir, jscompDir, parsingDir, unsafeParsingDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	// Synthetic shims/nodestore exposing BlockShift, NodeBlock, BlockOf, (*NodeBlock).FirstNextTok,
	// AllocNode1Cached, and AllocStringNodeCached.
	nsSrc := `package nodestore

import "sync/atomic"

const BlockShift = 14
const blockLen = 1 << BlockShift
const blockMask = blockLen - 1

type Node uint32
type StringID uint32

type NodeData struct {
	First  Node
	Next   Node
	Parent Node
	Str    StringID
	Tmpl   Node
	Token  uint8
}

type NodeBlock [blockLen]NodeData

var blocks [16]atomic.Pointer[NodeBlock]
var stringsPool = []string{""}

var (
	BlockOfCalls               atomic.Int64
	FirstNextTokDirectCalls    atomic.Int64
	BlockFirstNextTokCalls     atomic.Int64
	AllocNode1CachedCalls      atomic.Int64
	AllocNode1Calls            atomic.Int64
	AllocStringNodeCachedCalls atomic.Int64
	AllocStringNodeCalls       atomic.Int64
)

func ResetNodestoreMetrics() {
	BlockOfCalls.Store(0)
	FirstNextTokDirectCalls.Store(0)
	BlockFirstNextTokCalls.Store(0)
	AllocNode1CachedCalls.Store(0)
	AllocNode1Calls.Store(0)
	AllocStringNodeCachedCalls.Store(0)
	AllocStringNodeCalls.Store(0)
}

func EnsureBlock(blkIdx uint32) *NodeBlock {
	if b := blocks[blkIdx].Load(); b != nil {
		return b
	}
	nb := &NodeBlock{}
	if blocks[blkIdx].CompareAndSwap(nil, nb) {
		return nb
	}
	return blocks[blkIdx].Load()
}

func AllocAt(blkIdx, slot uint32, tok uint8) Node {
	b := EnsureBlock(blkIdx)
	id := Node((blkIdx << BlockShift) | (slot & blockMask))
	b[slot&blockMask] = NodeData{Token: tok}
	return id
}

func AppendChild(parent, child Node) {
	pb := blocks[uint32(parent)>>BlockShift].Load()
	cb := blocks[uint32(child)>>BlockShift].Load()
	cb[uint32(child)&blockMask].Parent = parent
	cb[uint32(child)&blockMask].Next = 0
	first := pb[uint32(parent)&blockMask].First
	if first == 0 {
		pb[uint32(parent)&blockMask].First = child
		return
	}
	cur := first
	for {
		curB := blocks[uint32(cur)>>BlockShift].Load()
		nx := curB[uint32(cur)&blockMask].Next
		if nx == 0 {
			curB[uint32(cur)&blockMask].Next = child
			return
		}
		cur = nx
	}
}

func BlockOf(n Node) *NodeBlock {
	BlockOfCalls.Add(1)
	b := blocks[uint32(n)>>BlockShift].Load()
	_ = &b[0]
	return b
}

func (b *NodeBlock) FirstNextTok(n Node) (first, next Node, tok uint8) {
	BlockFirstNextTokCalls.Add(1)
	rec := &b[uint32(n)&blockMask]
	return rec.First, rec.Next, rec.Token
}

func (n Node) Rec() *NodeData {
	b := blocks[uint32(n)>>BlockShift].Load()
	return &b[uint32(n)&blockMask]
}

func (n Node) FirstNextTok() (first, next Node, tok uint8) {
	FirstNextTokDirectCalls.Add(1)
	rec := n.Rec()
	return rec.First, rec.Next, rec.Token
}

type NodeAllocator struct {
	nextSlot uint32
	tmpl     Node
}

func (a *NodeAllocator) InitLocal() {
	EnsureBlock(1)
	if a.nextSlot == 0 {
		a.nextSlot = 1
	}
}

func (a *NodeAllocator) SetTemplate(tmpl Node) {
	a.tmpl = tmpl
}

func (a *NodeAllocator) Flush() {}

func (a *NodeAllocator) Intern(s string) StringID {
	stringsPool = append(stringsPool, s)
	return StringID(len(stringsPool) - 1)
}

func (a *NodeAllocator) allocRaw(tok uint8, tmpl Node) Node {
	b := EnsureBlock(1)
	slot := a.nextSlot
	a.nextSlot++
	b[slot] = NodeData{Token: tok, Tmpl: tmpl}
	return Node((1 << BlockShift) | slot)
}

func (a *NodeAllocator) AllocNode1Cached(tok uint8) Node {
	AllocNode1CachedCalls.Add(1)
	return a.allocRaw(tok, a.tmpl)
}

func (a *NodeAllocator) AllocNode1(tok uint8, tmpl Node) Node {
	AllocNode1Calls.Add(1)
	return a.allocRaw(tok, tmpl)
}

func (a *NodeAllocator) AllocNode2(tok uint8, c1, tmpl Node) Node {
	n := a.allocRaw(tok, tmpl)
	AppendChild(n, c1)
	return n
}

func (a *NodeAllocator) AllocNode3(tok uint8, c1, c2, tmpl Node) Node {
	n := a.allocRaw(tok, tmpl)
	AppendChild(n, c1)
	AppendChild(n, c2)
	return n
}

func (a *NodeAllocator) AllocNode4(tok uint8, c1, c2, c3, tmpl Node) Node {
	n := a.allocRaw(tok, tmpl)
	AppendChild(n, c1)
	AppendChild(n, c2)
	AppendChild(n, c3)
	return n
}

func (a *NodeAllocator) AllocStringNodeCached(tok uint8, strID StringID) Node {
	AllocStringNodeCachedCalls.Add(1)
	n := a.allocRaw(tok, a.tmpl)
	n.Rec().Str = strID
	return n
}

func (a *NodeAllocator) AllocStringNode(tok uint8, strID StringID, tmpl Node) Node {
	AllocStringNodeCalls.Add(1)
	n := a.allocRaw(tok, tmpl)
	n.Rec().Str = strID
	return n
}

func (a *NodeAllocator) AllocNumberNode(tok uint8, val float64, tmpl Node) Node {
	_ = val
	return a.allocRaw(tok, tmpl)
}

func LookupString(id StringID) string {
	return stringsPool[id]
}
`
	if err := os.WriteFile(filepath.Join(nsDir, "nodestore.go"), []byte(nsSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	j2Src := `package j2gort

import (
	"sync"
	"sync/atomic"
)

const NullString = "\x00"

var (
	PooledSlabChunkAllocs        atomic.Int64
	PooledSlabReleaseCalls       atomic.Int64
	PooledSlabReleaseWithCalls   atomic.Int64
	CustomOrderedMapTier0Acq     atomic.Int64
	CustomOrderedMapTier2CapAcq  atomic.Int64
	CustomOrderedMapLastCapHint  atomic.Int64
	CustomOrderedMapReleases     atomic.Int64
	OrderedMapReleases           atomic.Int64
)

func ResetW22CMetrics() {
	PooledSlabChunkAllocs.Store(0)
	PooledSlabReleaseCalls.Store(0)
	PooledSlabReleaseWithCalls.Store(0)
	CustomOrderedMapTier0Acq.Store(0)
	CustomOrderedMapTier2CapAcq.Store(0)
	CustomOrderedMapLastCapHint.Store(0)
	CustomOrderedMapReleases.Store(0)
	OrderedMapReleases.Store(0)
}

func ThrowNPE() { panic("NullPointerException") }
func JStrLen(s string) int32 { return int32(len(s)) }
func ValueOfDouble(v float64) any { return v }
func UnboxDouble(v any) float64 { return v.(float64) }

func CheckNotNullIface[I interface{ IsJavaNil_() bool }](v I) I {
	if any(v) == nil || v.IsJavaNil_() {
		panic("NPE")
	}
	return v
}

func TryCast[T any](v any) (T, bool) {
	var z T
	if v == nil {
		return z, false
	}
	r, ok := v.(T)
	return r, ok
}

func AsIface[T any](v any) T {
	var z T
	if v == nil {
		return z
	}
	return v.(T)
}

func CastPtr[T any](v any) *T {
	if v == nil {
		return nil
	}
	return v.(*T)
}

func PtrToAny[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}

type Slab[T any] struct{ buf []T }

func (s *Slab[T]) Alloc() *T {
	if len(s.buf) == cap(s.buf) {
		c := make([]T, 1, 32)
		s.buf = c
		return &c[0]
	}
	i := len(s.buf)
	s.buf = s.buf[:i+1]
	return &s.buf[i]
}

func (s *Slab[T]) Release() {
	if s != nil {
		s.buf = nil
	}
}

var pooledChunkPools sync.Map

func chunkPoolFor[T any]() *sync.Pool {
	var z *T
	if p, ok := pooledChunkPools.Load(z); ok {
		return p.(*sync.Pool)
	}
	np := &sync.Pool{}
	actual, _ := pooledChunkPools.LoadOrStore(z, np)
	return actual.(*sync.Pool)
}

type PooledSlab[T any] struct {
	buf    []T
	chunks [][]T
}

func (s *PooledSlab[T]) Alloc() *T {
	i := len(s.buf)
	if i == cap(s.buf) {
		if cap(s.buf) > 0 {
			s.chunks = append(s.chunks, s.buf)
		}
		var c []T
		if v := chunkPoolFor[T]().Get(); v != nil {
			p := v.(*[]T)
			c = (*p)[:1]
		} else {
			PooledSlabChunkAllocs.Add(1)
			c = make([]T, 1, 16)
		}
		s.buf = c
		return &c[0]
	}
	b := s.buf[:i+1]
	s.buf = b
	return &b[i]
}

func (s *PooledSlab[T]) Release() {
	if s == nil {
		return
	}
	PooledSlabReleaseCalls.Add(1)
	pool := chunkPoolFor[T]()
	for idx := range s.chunks {
		ch := s.chunks[idx]
		clear(ch)
		ch = ch[:0]
		pool.Put(&ch)
		s.chunks[idx] = nil
	}
	s.chunks = s.chunks[:0]
	if cap(s.buf) > 0 {
		ch := s.buf
		clear(ch)
		ch = ch[:0]
		pool.Put(&ch)
		s.buf = nil
	}
}

func (s *PooledSlab[T]) ReleaseWith(cleanup func(*T)) {
	if s == nil {
		return
	}
	if cleanup == nil {
		s.Release()
		return
	}
	PooledSlabReleaseWithCalls.Add(1)
	pool := chunkPoolFor[T]()
	for idx := range s.chunks {
		ch := s.chunks[idx]
		for j := range ch {
			cleanup(&ch[j])
		}
		clear(ch)
		ch = ch[:0]
		pool.Put(&ch)
		s.chunks[idx] = nil
	}
	s.chunks = s.chunks[:0]
	if cap(s.buf) > 0 {
		ch := s.buf
		for j := range ch {
			cleanup(&ch[j])
		}
		clear(ch)
		ch = ch[:0]
		pool.Put(&ch)
		s.buf = nil
	}
}

type OrderedMap[K comparable, V any] struct {
	keys []K
	vals []V
	idx  map[K]int
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] { return OrderedMap[K, V]{} }
func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] { return OrderedMap[K, V]{} }
func (m *OrderedMap[K, V]) Put(k K, v V) {
	if m.idx == nil {
		m.idx = make(map[K]int)
	}
	if i, ok := m.idx[k]; ok {
		m.vals[i] = v
		return
	}
	m.idx[k] = len(m.keys)
	m.keys = append(m.keys, k)
	m.vals = append(m.vals, v)
}
func (m *OrderedMap[K, V]) Get(k K) V {
	var z V
	if m == nil || m.idx == nil {
		return z
	}
	if i, ok := m.idx[k]; ok {
		return m.vals[i]
	}
	return z
}
func (m *OrderedMap[K, V]) Clear() {
	if m != nil {
		clear(m.keys)
		clear(m.vals)
		m.keys = m.keys[:0]
		m.vals = m.vals[:0]
		clear(m.idx)
	}
}
func (m *OrderedMap[K, V]) ReleaseToPool() {
	if m != nil {
		OrderedMapReleases.Add(1)
		*m = OrderedMap[K, V]{}
	}
}

type FastCustomKey interface {
	IsJavaNil_() bool
	HashCode() int32
	Equals(any) bool
}

type FastCustomKeyComparable interface {
	FastCustomKey
	comparable
}

type CustomOrderedMap[K FastCustomKeyComparable, V any] struct {
	keys  []K
	vals  []V
	tier2 bool
}

func NewCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{}
}

func NewCustomOrderedMapCap[K FastCustomKeyComparable, V any](capHint int32) CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{
		keys:  make([]K, 0, capHint),
		vals:  make([]V, 0, capHint),
		tier2: capHint > 512,
	}
}

func AcquireCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
	CustomOrderedMapTier0Acq.Add(1)
	return CustomOrderedMap[K, V]{}
}

func AcquireCustomOrderedMapCap[K FastCustomKeyComparable, V any](capHint int32) CustomOrderedMap[K, V] {
	CustomOrderedMapTier2CapAcq.Add(1)
	CustomOrderedMapLastCapHint.Store(int64(capHint))
	return NewCustomOrderedMapCap[K, V](capHint)
}

func (m *CustomOrderedMap[K, V]) Size() int32 {
	if m == nil {
		return 0
	}
	return int32(len(m.keys))
}

func (m *CustomOrderedMap[K, V]) IsTier2() bool {
	return m != nil && m.tier2
}

func (m *CustomOrderedMap[K, V]) Get(k K) V {
	var z V
	if m == nil {
		return z
	}
	for i, ek := range m.keys {
		if ek == k || (!ek.IsJavaNil_() && ek.Equals(k)) {
			return m.vals[i]
		}
	}
	return z
}

func (m *CustomOrderedMap[K, V]) Put(k K, v V) {
	for i, ek := range m.keys {
		if ek == k || (!ek.IsJavaNil_() && ek.Equals(k)) {
			m.vals[i] = v
			return
		}
	}
	m.keys = append(m.keys, k)
	m.vals = append(m.vals, v)
}

func (m *CustomOrderedMap[K, V]) Clear() {
	if m != nil {
		clear(m.keys)
		clear(m.vals)
		m.keys = m.keys[:0]
		m.vals = m.vals[:0]
	}
}

func (m *CustomOrderedMap[K, V]) KeysSlice() []K {
	if m == nil {
		return nil
	}
	return m.keys
}

func (m *CustomOrderedMap[K, V]) ReleaseToPool() {
	if m != nil {
		CustomOrderedMapReleases.Add(1)
		*m = CustomOrderedMap[K, V]{}
	}
}

type Deque[T any] struct {
	items []T
}

func (d *Deque[T]) Push(v T)    { d.items = append(d.items, v) }
func (d *Deque[T]) AddLast(v T) { d.items = append(d.items, v) }
func (d *Deque[T]) Peek() T {
	var z T
	if len(d.items) == 0 {
		return z
	}
	return d.items[len(d.items)-1]
}
func (d *Deque[T]) PeekLast() T { return d.Peek() }
func (d *Deque[T]) GetLast() T  { return d.Peek() }
func (d *Deque[T]) Pop() T {
	n := len(d.items) - 1
	v := d.items[n]
	var z T
	d.items[n] = z
	d.items = d.items[:n]
	return v
}
func (d *Deque[T]) RemoveLast() T { return d.Pop() }
func (d *Deque[T]) IsEmpty() bool { return len(d.items) == 0 }
`
	if err := os.WriteFile(filepath.Join(j2Dir, "runtime.go"), []byte(j2Src), 0o644); err != nil {
		t.Fatal(err)
	}

	rhinoSrc := `package rhino

import "example.com/w22ctest/shims/nodestore"

const (
	TokSCRIPT            int32 = 1
	TokFUNCTION          int32 = 2
	TokMODULE_BODY       int32 = 3
	TokCLASS             int32 = 4
	TokCLASS_MEMBERS     int32 = 5
	TokDEFAULT_VALUE     int32 = 6
	TokDESTRUCTURING_LHS int32 = 7
	TokBLOCK             int32 = 8
	TokFOR               int32 = 9
	TokIF                int32 = 10
	TokAND               int32 = 11
	TokOR                int32 = 12
	TokNAME              int32 = 13
	TokSTRINGLIT         int32 = 14
	TokSTRING_KEY        int32 = 15
	TokNUMBER            int32 = 16
)

func StaticGet_Token_STRINGLIT() int32  { return TokSTRINGLIT }
func StaticGet_Token_STRING_KEY() int32 { return TokSTRING_KEY }
func StaticGet_Token_NUMBER() int32     { return TokNUMBER }

type Node nodestore.Node
type InputId struct{ Name string }

func (n Node) IsJavaNil_() bool { return n == 0 }

func (n Node) FirstNextTok_() (Node, Node, uint8) {
	f, nx, tok := nodestore.Node(n).FirstNextTok()
	return Node(f), Node(nx), tok
}

func (n Node) TokenSwitchOrd_() int32 {
	if n == 0 {
		panic("NPE")
	}
	t := nodestore.Node(n).Rec().Token
	if t == 0xFF {
		panic("NPE")
	}
	return int32(t) + 1
}

func (n Node) GetFirstChild() Node {
	if n == 0 {
		panic("NPE")
	}
	return Node(nodestore.Node(n).Rec().First)
}

func (n Node) GetNext() Node {
	if n == 0 {
		panic("NPE")
	}
	return Node(nodestore.Node(n).Rec().Next)
}

func (n Node) GetParent() Node {
	if n == 0 {
		panic("NPE")
	}
	return Node(nodestore.Node(n).Rec().Parent)
}

func (n Node) SetLinenoCharno(line, col int32) Node {
	_, _ = line, col
	return n
}

func (n Node) SetLength(l int32) { _ = l }

func (n Node) ClonePropsFrom(tmpl Node) Node {
	nodestore.Node(n).Rec().Tmpl = nodestore.Node(tmpl)
	return n
}

func (n Node) TemplateOf() Node {
	return Node(nodestore.Node(n).Rec().Tmpl)
}

func (n Node) StringVal() string {
	return nodestore.LookupString(nodestore.Node(n).Rec().Str)
}

func NewNode_1(type_ int32) Node {
	return Node(nodestore.AllocAt(2, 1, uint8(type_)-1))
}
func NewNode_2(type_ int32, c1 Node) Node {
	p := NewNode_1(type_)
	nodestore.AppendChild(nodestore.Node(p), nodestore.Node(c1))
	return p
}
func NewNode_3(type_ int32, c1, c2 Node) Node {
	p := NewNode_2(type_, c1)
	nodestore.AppendChild(nodestore.Node(p), nodestore.Node(c2))
	return p
}
func NewNode_4(type_ int32, c1, c2, c3 Node) Node {
	p := NewNode_3(type_, c1, c2)
	nodestore.AppendChild(nodestore.Node(p), nodestore.Node(c3))
	return p
}
func Node_NewString_1(s string) Node {
	_ = s
	return NewNode_1(TokSTRINGLIT)
}
func Node_NewString_2(type_ int32, s string) Node {
	_ = s
	return NewNode_1(type_)
}
func IR_Number(d float64) Node {
	_ = d
	return NewNode_1(TokNUMBER)
}
`
	if err := os.WriteFile(filepath.Join(rhinoDir, "node.go"), []byte(rhinoSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1. Parsing package testing IRFactory AllocNode1Cached + AllocStringNodeCached AND JsDocInfoParser non-cached safety!
	irAndJsDocSrc := `package parsing

import (
	"example.com/w22ctest/gen/j2gort"
	pkg_rhino "example.com/w22ctest/gen/rhino"
)

type TokenStream struct{}
func (s *TokenStream) GetLineno() int32 { return 1 }
func (s *TokenStream) GetCharno() int32 { return 0 }

type JsDocInfoParser struct {
	templateNode pkg_rhino.Node
	stream       *TokenStream
}

func InitJsDocInfoParser(this_ *JsDocInfoParser, stream *TokenStream, tmpl pkg_rhino.Node) {
	this_.stream = stream
	this_.templateNode = tmpl
}

func (this_ *JsDocInfoParser) priv_JsDocInfoParser_newNode_1(type_ int32) pkg_rhino.Node {
	return pkg_rhino.NewNode_1(type_).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno()).ClonePropsFrom(this_.templateNode)
}

func (this_ *JsDocInfoParser) priv_JsDocInfoParser_wrapNode_2(type_ int32, n pkg_rhino.Node, lineno int32, charno int32) pkg_rhino.Node {
	if n == 0 {
		return pkg_rhino.Node(0)
	}
	return pkg_rhino.NewNode_2(type_, n).SetLinenoCharno(lineno, charno).ClonePropsFrom(this_.templateNode)
}

func (this_ *JsDocInfoParser) priv_JsDocInfoParser_newStringNode_3(s string) pkg_rhino.Node {
	var n pkg_rhino.Node = pkg_rhino.Node_NewString_1(s).SetLinenoCharno(this_.stream.GetLineno(), this_.stream.GetCharno()).ClonePropsFrom(this_.templateNode)
	n.SetLength(j2gort.JStrLen(s))
	return n
}

type ProgramTree struct{}

type IRFactory struct {
	templateNode               pkg_rhino.Node
	defaultTemplateNode        pkg_rhino.Node
	closureUnawareTemplateNode pkg_rhino.Node
}

func (this_ *IRFactory) initFields_IRFactory() {}

func NewIRFactory(defaultTmpl, cuTmpl pkg_rhino.Node) *IRFactory {
	this_ := &IRFactory{}
	InitIRFactory(this_, defaultTmpl, cuTmpl)
	return this_
}

func InitIRFactory(this_ *IRFactory, defaultTmpl, cuTmpl pkg_rhino.Node) {
	this_.initFields_IRFactory()
	this_.defaultTemplateNode = defaultTmpl
	this_.closureUnawareTemplateNode = cuTmpl
	this_.templateNode = this_.defaultTemplateNode
}

func (this_ *IRFactory) ParseJsDocString(jsdocTmpl pkg_rhino.Node, s string) pkg_rhino.Node {
	var jsdoc JsDocInfoParser
	InitJsDocInfoParser(&jsdoc, &TokenStream{}, jsdocTmpl)
	return jsdoc.priv_JsDocInfoParser_newStringNode_3(s)
}

func IRFactory_TransformTree(defaultTmpl, cuTmpl, jsdocTmpl pkg_rhino.Node) (pkg_rhino.Node, pkg_rhino.Node, pkg_rhino.Node, pkg_rhino.Node) {
	var irFactory *IRFactory = NewIRFactory(defaultTmpl, cuTmpl)
	_ = irFactory
	n1 := irFactory.NewNode_1(pkg_rhino.TokNAME)
	irFactory.SwitchClosureUnaware(true)
	s1 := irFactory.NewStringNode_1("hello")
	jn := irFactory.ParseJsDocString(jsdocTmpl, "docstr")
	irFactory.SwitchClosureUnaware(false)
	s2 := irFactory.NewStringNode_2(pkg_rhino.TokNAME, "world")
	return n1, s1, jn, s2
}

func (this_ *IRFactory) SwitchClosureUnaware(cu bool) {
	if cu {
		this_.templateNode = this_.closureUnawareTemplateNode
	} else {
		this_.templateNode = this_.defaultTemplateNode
	}
}

func (this_ *IRFactory) NewNode_1(type_ int32) pkg_rhino.Node {
	return pkg_rhino.NewNode_1(type_).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewNode_2(type_ int32, child1 pkg_rhino.Node) pkg_rhino.Node {
	return pkg_rhino.NewNode_2(type_, child1).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewNode_3(type_ int32, child1, child2 pkg_rhino.Node) pkg_rhino.Node {
	return pkg_rhino.NewNode_3(type_, child1, child2).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewNode_4(type_ int32, child1, child2, child3 pkg_rhino.Node) pkg_rhino.Node {
	return pkg_rhino.NewNode_4(type_, child1, child2, child3).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewStringNode_1(value string) pkg_rhino.Node {
	return pkg_rhino.Node_NewString_1(value).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewStringNode_2(type_ int32, value string) pkg_rhino.Node {
	return pkg_rhino.Node_NewString_2(type_, value).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewNumberNode(value any) pkg_rhino.Node {
	return pkg_rhino.IR_Number(j2gort.UnboxDouble(value)).ClonePropsFrom(this_.templateNode)
}
`
	if err := os.WriteFile(filepath.Join(parsingDir, "ir_factory.go"), []byte(irAndJsDocSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// Negative soundness case for IRFactory cached template: taking `&this_.templateNode` aliases templateNode!
	unsafeIRSrc := strings.Replace(irAndJsDocSrc, "package parsing", "package unsafeparsing", 1) + `
func (this_ *IRFactory) AliasTemplatePtr() *pkg_rhino.Node {
	return &this_.templateNode
}
`
	if err := os.WriteFile(filepath.Join(unsafeParsingDir, "ir_factory.go"), []byte(unsafeIRSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 2. jscomp package testing:
	//  - Block-cached child iteration in traverseBranch + traverseChildren + gc_ grandchild loop
	//  - NodeTraversal_TraverseScopeRoots pooling + dynamic CanReleaseSyntacticScopes_ release
	//  - ReferenceCollector.Process_2 Tier-2 AcquireCustomOrderedMapCap(1024) vs Tier-0 Process_1
	//  - InlineVariables inner ephemeral structs PooledSlab + static singleton preservation
	//  - GuardedCallback_Context PooledSlab + EMPTY singleton preservation + nested stack discipline
	jscompNodeTraversalSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w22ctest/gen/j2gort"
	pkg_rhino "example.com/w22ctest/gen/rhino"
)

var NewNodeTraversalHeapCalls atomic.Int64

type CompilerOptions struct{}
func (o *CompilerOptions) GetSyntheticBlockStartMarker() string { return j2gort.NullString }

type AbstractCompiler interface {
	IsJavaNil_() bool
	GetOptions() *CompilerOptions
}

type CompilerInput struct{}

type AbstractScope interface {
	IsJavaNil_() bool
	GetRootNode_0() pkg_rhino.Node
	GetParent() any
}

func ToIface_AbstractScope(s *Scope) AbstractScope {
	if s == nil {
		return nil
	}
	return s
}

type AbstractScopeBase struct {
	self_    AbstractScope
	rootNode pkg_rhino.Node
}

func (b *AbstractScopeBase) BindSelf_(s any) { b.self_ = s.(AbstractScope) }
func InitAbstractScope(this_ *AbstractScopeBase, root pkg_rhino.Node) { this_.rootNode = root }

type Var struct {
	name  string
	scope *Scope
	index int32
}

func (v *Var) IsJavaNil_() bool  { return v == nil }
func (v *Var) HashCode() int32   { return v.index }
func (v *Var) Equals(o any) bool { x, ok := o.(*Var); return ok && x == v }
func (this_ *Var) initFields_Var() {}

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

type Scope struct {
	AbstractScopeBase
	parent *Scope
	vars   j2gort.OrderedMap[string, *Var]
}

func (this_ *Scope) BindSelf_(s any)               { this_.AbstractScopeBase.BindSelf_(s) }
func (this_ *Scope) IsJavaNil_() bool              { return this_ == nil }
func (this_ *Scope) GetRootNode_0() pkg_rhino.Node { return this_.rootNode }
func (this_ *Scope) GetParent() any                { return j2gort.PtrToAny(this_.parent) }
func (this_ *Scope) IsHoistScope() bool            { return true }
func (this_ *Scope) initFields_Scope() {
	this_.vars = j2gort.NewOrderedMap[string, *Var]()
}

func NewScope_1(root pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_1(this_, root)
	return this_
}

func InitScope_1(this_ *Scope, root pkg_rhino.Node) {
	InitAbstractScope(&this_.AbstractScopeBase, root)
	this_.initFields_Scope()
	this_.parent = nil
}

func NewScope_2(parent *Scope, root pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_2(this_, parent, root)
	return this_
}

func InitScope_2(this_ *Scope, parent *Scope, root pkg_rhino.Node) {
	InitAbstractScope(&this_.AbstractScopeBase, root)
	this_.initFields_Scope()
	this_.parent = parent
}

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope_2_Node_AbstractScope(pkg_rhino.Node, AbstractScope) any
}

type SyntacticScopeCreator_RedeclarationHandler interface {
	IsJavaNil_() bool
}

type defaultHandlerImpl struct{}
func (d *defaultHandlerImpl) IsJavaNil_() bool { return d == nil }
var defHandler SyntacticScopeCreator_RedeclarationHandler = &defaultHandlerImpl{}
func StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER() SyntacticScopeCreator_RedeclarationHandler {
	return defHandler
}

type SyntacticScopeCreator_ScopeScanner struct {
	scope                         *Scope
	compiler                      AbstractCompiler
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler
	treatProvidesAsRedeclarations bool
}

func (s *SyntacticScopeCreator_ScopeScanner) Reset() { s.scope = nil }

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
	this_.reusableScanner = &SyntacticScopeCreator_ScopeScanner{
		compiler:             compiler,
		redeclarationHandler: this_.redeclarationHandler,
	}
}

func (this_ *SyntacticScopeCreator) CreateScope_2_Node_Scope(n pkg_rhino.Node, parent *Scope) *Scope {
	var sc *Scope
	if parent == nil {
		sc = NewScope_1(n)
	} else {
		sc = NewScope_2(parent, n)
	}
	v := NewVar("x", sc, int32(n))
	sc.vars.Put("x", v)
	return sc
}

func (this_ *SyntacticScopeCreator) CreateScope_2_Node_AbstractScope(n pkg_rhino.Node, parent AbstractScope) any {
	var p *Scope
	if parent != nil && !parent.IsJavaNil_() {
		p = parent.(*Scope)
	}
	return this_.CreateScope_2_Node_Scope(n, p)
}

type MemoizedScopeCreator struct {
	delegate          ScopeCreator
	scopesByScopeRoot j2gort.OrderedMap[pkg_rhino.Node, AbstractScope]
}

func (m *MemoizedScopeCreator) IsJavaNil_() bool { return m == nil }
func (this_ *MemoizedScopeCreator) initFields_MemoizedScopeCreator() {
	this_.scopesByScopeRoot = j2gort.NewOrderedMap[pkg_rhino.Node, AbstractScope]()
}

func InitMemoizedScopeCreator(this_ *MemoizedScopeCreator, delegate ScopeCreator) {
	this_.initFields_MemoizedScopeCreator()
	this_.delegate = delegate
}

func (this_ *MemoizedScopeCreator) CreateScope_2_Node_AbstractScope(n pkg_rhino.Node, parent AbstractScope) any {
	if s := this_.scopesByScopeRoot.Get(n); s != nil {
		return s
	}
	s := this_.delegate.CreateScope_2_Node_AbstractScope(n, parent).(AbstractScope)
	this_.scopesByScopeRoot.Put(n, s)
	return s
}

type NodeTraversal_Callback interface {
	IsJavaNil_() bool
	ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool
	Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

func ToIface_NodeTraversal_Callback(cb NodeTraversal_Callback) NodeTraversal_Callback { return cb }

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

func clinit_NodeTraversal() {}

func (this_ *NodeTraversal) initFields_NodeTraversal() {
	this_.scopeDepth = 0
	this_.sourceName = j2gort.NullString
}

func NewNodeTraversal(builder *NodeTraversal_Builder) *NodeTraversal {
	NewNodeTraversalHeapCalls.Add(1)
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
	if builder.scopeCreator == nil || builder.scopeCreator.IsJavaNil_() {
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

func (this_ *NodeTraversal) priv_NodeTraversal_ensureScopeCapacity_1() {
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

func (this_ *NodeTraversal) priv_NodeTraversal_pushScope_2(n pkg_rhino.Node) {
	this_.priv_NodeTraversal_ensureScopeCapacity_1()
	this_.scopeNodes[this_.scopeDepth] = n
	this_.scopeObjs[this_.scopeDepth] = nil
	this_.scopeDepth++
	if this_.scopeCallback != nil && !this_.scopeCallback.IsJavaNil_() {
		this_.scopeCallback.EnterScope(this_)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_popScope_3() {
	if this_.scopeCallback != nil && !this_.scopeCallback.IsJavaNil_() {
		this_.scopeCallback.ExitScope(this_)
	}
	this_.scopeDepth--
	this_.scopeNodes[this_.scopeDepth] = 0
	this_.scopeObjs[this_.scopeDepth] = nil
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

func (this_ *NodeTraversal) GetScopeRoot() pkg_rhino.Node {
	return this_.scopeNodes[this_.scopeDepth-1]
}

func (this_ *NodeTraversal) priv_NodeTraversal_createsBlockScope_4(n pkg_rhino.Node) bool {
	switch n.TokenSwitchOrd_() {
	case int32(8) /* BLOCK */, int32(9) /* FOR */, int32(4) /* CLASS */:
		return true
	default:
		return false
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleScript_10(n, parent pkg_rhino.Node) {
	this_.currentNode = n
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		this_.priv_NodeTraversal_traverseChildren_20(n)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
}
func (this_ *NodeTraversal) priv_NodeTraversal_handleFunction_11(n, parent pkg_rhino.Node) {
	this_.priv_NodeTraversal_pushScope_2(n)
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		this_.priv_NodeTraversal_traverseChildren_20(n)
		this_.callback.Visit_3(this_, n, parent)
	}
	this_.priv_NodeTraversal_popScope_3()
}
func (this_ *NodeTraversal) priv_NodeTraversal_handleModule_12(n, parent pkg_rhino.Node) {
	this_.priv_NodeTraversal_handleScript_10(n, parent)
}
func (this_ *NodeTraversal) priv_NodeTraversal_handleClass_13(n, parent pkg_rhino.Node) {
	this_.priv_NodeTraversal_handleScript_10(n, parent)
}
func (this_ *NodeTraversal) priv_NodeTraversal_handleClassMembers_14(n, parent pkg_rhino.Node) {
	this_.priv_NodeTraversal_handleScript_10(n, parent)
}
func (this_ *NodeTraversal) priv_NodeTraversal_handleDestructuringOrDefaultValue_15(n, parent pkg_rhino.Node) {
	this_.priv_NodeTraversal_handleScript_10(n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseBranch_19(n pkg_rhino.Node, parent pkg_rhino.Node) {
	switch n.TokenSwitchOrd_() {
	case int32(1):
		this_.priv_NodeTraversal_handleScript_10(n, parent)
		return
	case int32(2):
		this_.priv_NodeTraversal_handleFunction_11(n, parent)
		return
	case int32(3):
		this_.priv_NodeTraversal_handleModule_12(n, parent)
		return
	case int32(4):
		this_.priv_NodeTraversal_handleClass_13(n, parent)
		return
	case int32(5):
		this_.priv_NodeTraversal_handleClassMembers_14(n, parent)
		return
	case int32(6), int32(7):
		if this_.obeyDestructuringAndDefaultValueExecutionOrder {
			this_.priv_NodeTraversal_handleDestructuringOrDefaultValue_15(n, parent)
			return
		}
	default:
	}
	this_.currentNode = n
	if !(this_.callback.ShouldTraverse_3(this_, n, parent)) {
		return
	}
	var createsBlockScope bool = this_.priv_NodeTraversal_createsBlockScope_4(n)
	if createsBlockScope {
		this_.priv_NodeTraversal_pushScope_2(n)
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for child != 0 {
			var next pkg_rhino.Node = child.GetNext()
			_ = next
			this_.priv_NodeTraversal_traverseBranch_19(child, n)
			child = next
		}
	}
	if createsBlockScope {
		this_.priv_NodeTraversal_popScope_3()
	}
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseChildren_20(n pkg_rhino.Node) {
	var child pkg_rhino.Node = n.GetFirstChild()
	_ = child
	for child != 0 {
		var next pkg_rhino.Node = child.GetNext()
		_ = next
		this_.priv_NodeTraversal_traverseBranch_19(child, n)
		child = next
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverse_21(root pkg_rhino.Node) {
	this_.currentNode = root
	this_.priv_NodeTraversal_pushScope_2(root)
	this_.priv_NodeTraversal_traverseBranch_19(root, 0)
	this_.priv_NodeTraversal_popScope_3()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseRoots_22(externs, root pkg_rhino.Node) {
	this_.currentNode = externs
	this_.priv_NodeTraversal_pushScope_2(externs)
	this_.priv_NodeTraversal_traverseBranch_19(externs, 0)
	this_.priv_NodeTraversal_traverseBranch_19(root, 0)
	this_.priv_NodeTraversal_popScope_3()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseWithScope_23(root pkg_rhino.Node, s AbstractScope) {
	this_.currentNode = root
	this_.priv_NodeTraversal_ensureScopeCapacity_1()
	this_.scopeNodes[this_.scopeDepth] = s.GetRootNode_0()
	this_.scopeObjs[this_.scopeDepth] = s
	this_.scopeDepth++
	this_.priv_NodeTraversal_traverseBranch_19(root, 0)
	this_.priv_NodeTraversal_popScope_3()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseAtScope_24(s AbstractScope) {
	this_.priv_NodeTraversal_traverseWithScope_23(s.GetRootNode_0(), s)
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseScopeRoot_25(scopeRoot pkg_rhino.Node) {
	this_.currentNode = scopeRoot
	if scopeRoot.GetParent() != 0 {
		this_.priv_NodeTraversal_pushScope_2(scopeRoot.GetParent())
	}
	this_.priv_NodeTraversal_traverseBranch_19(scopeRoot, scopeRoot.GetParent())
}

func NodeTraversal_Traverse(compiler AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).Traverse(root)
}

func NodeTraversal_TraverseScopeRoots(compiler AbstractCompiler, scopeNodes []pkg_rhino.Node, cb NodeTraversal_Callback, traverseNested bool) {
	clinit_NodeTraversal()
	_ = traverseNested
	var val_8 MemoizedScopeCreator
	InitMemoizedScopeCreator(&val_8, NewSyntacticScopeCreator_1(compiler))
	var scopeCreator *MemoizedScopeCreator = &val_8
	_ = scopeCreator
	for _, scopeNode := range scopeNodes {
		NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).SetScopeCreator(scopeCreator).Build().priv_NodeTraversal_traverseScopeRoot_25(scopeNode)
	}
}

type NodeTraversal_Builder struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

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
func (b *NodeTraversal_Builder) SetObeyDestructuringAndDefaultValueExecutionOrder(v bool) *NodeTraversal_Builder {
	b.obeyDestructuringAndDefaultValueExecutionOrder = v
	return b
}
func (b *NodeTraversal_Builder) Build() *NodeTraversal { return NewNodeTraversal(b) }
func (b *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	b.Build().priv_NodeTraversal_traverse_21(root)
}
func (b *NodeTraversal_Builder) TraverseRoots(externs, root pkg_rhino.Node) {
	b.Build().priv_NodeTraversal_traverseRoots_22(externs, root)
}
`

	jscompPassesSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w22ctest/gen/j2gort"
	pkg_rhino "example.com/w22ctest/gen/rhino"
)

type BasicBlock struct {
	parent *BasicBlock
	root   pkg_rhino.Node
}

func (b *BasicBlock) initFields_BasicBlock() {}
func NewBasicBlock(parent *BasicBlock, root pkg_rhino.Node) *BasicBlock {
	this_ := &BasicBlock{}
	InitBasicBlock(this_, parent, root)
	return this_
}
func InitBasicBlock(this_ *BasicBlock, parent *BasicBlock, root pkg_rhino.Node) {
	this_.initFields_BasicBlock()
	this_.parent = parent
	this_.root = root
}

type Reference struct {
	node       pkg_rhino.Node
	basicBlock *BasicBlock
	scope      *Scope
}

func (r *Reference) initFields_Reference() {}
func NewReference_3(node pkg_rhino.Node, t *NodeTraversal, basicBlock *BasicBlock) *Reference {
	this_ := &Reference{}
	InitReference_3(this_, node, t, basicBlock)
	return this_
}
func InitReference_3(this_ *Reference, node pkg_rhino.Node, t *NodeTraversal, basicBlock *BasicBlock) {
	this_.initFields_Reference()
	this_.node = node
	this_.basicBlock = basicBlock
	this_.scope = t.GetScope()
}

type ReferenceCollection struct {
	rest []*Reference
}

func (rc *ReferenceCollection) initFields_ReferenceCollection() {}
func NewReferenceCollection() *ReferenceCollection {
	this_ := &ReferenceCollection{}
	InitReferenceCollection(this_)
	return this_
}
func InitReferenceCollection(this_ *ReferenceCollection) {
	this_.initFields_ReferenceCollection()
}
func (rc *ReferenceCollection) Add(r *Reference) { rc.rest = append(rc.rest, r) }
func (rc *ReferenceCollection) Size() int32      { return int32(len(rc.rest)) }
func (rc *ReferenceCollection) Get(i int32) *Reference { return rc.rest[i] }

type ReferenceCollector_Behavior interface {
	IsJavaNil_() bool
	AfterExitScope(*NodeTraversal, *ReferenceCollector)
}

type noopBehavior struct{}
func (n *noopBehavior) IsJavaNil_() bool                                 { return n == nil }
func (n *noopBehavior) AfterExitScope(*NodeTraversal, *ReferenceCollector) {}
var doNothingBeh ReferenceCollector_Behavior = &noopBehavior{}
func StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR() ReferenceCollector_Behavior { return doNothingBeh }

type ReferenceCollector struct {
	referenceMap j2gort.CustomOrderedMap[*Var, *ReferenceCollection]
	blockStack   *j2gort.Deque[*BasicBlock]
	behavior     ReferenceCollector_Behavior
	scopeCreator ScopeCreator
	compiler     AbstractCompiler
	callback     *ReferenceCollector_CollectorCallback
}

func clinit_ReferenceCollector() {}

func (this_ *ReferenceCollector) initFields_ReferenceCollector() {
	this_.referenceMap = j2gort.NewCustomOrderedMap[*Var, *ReferenceCollection]()
	this_.blockStack = &j2gort.Deque[*BasicBlock]{}
	this_.callback = NewReferenceCollector_CollectorCallback(this_)
}

func NewReferenceCollector_3(compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator) *ReferenceCollector {
	this_ := &ReferenceCollector{}
	InitReferenceCollector_3(this_, compiler, behavior, creator)
	return this_
}

func InitReferenceCollector_3(this_ *ReferenceCollector, compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator) {
	this_.initFields_ReferenceCollector()
	this_.compiler = compiler
	this_.behavior = behavior
	this_.scopeCreator = creator
}

func (this_ *ReferenceCollector) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(this_.callback).SetScopeCreator(this_.scopeCreator).TraverseRoots(externs, root)
}

func (this_ *ReferenceCollector) Process_1(root pkg_rhino.Node) {
	_ = root
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(this_.callback).SetScopeCreator(this_.scopeCreator).Traverse(root)
}

func (this_ *ReferenceCollector) GetAllSymbols() []*Var {
	return this_.referenceMap.KeysSlice()
}

func (this_ *ReferenceCollector) GetReferences(v *Var) *ReferenceCollection {
	return this_.referenceMap.Get(v)
}

func (this_ *ReferenceCollector) addRef(v *Var, n pkg_rhino.Node, t *NodeTraversal) {
	col := this_.referenceMap.Get(v)
	if col == nil {
		col = NewReferenceCollection()
		this_.referenceMap.Put(v, col)
	}
	var val_r Reference
	InitReference_3(&val_r, n, t, this_.blockStack.GetLast())
	var ref *Reference = &val_r
	col.Add(ref)
}

type ReferenceCollector_CollectorCallback struct {
	outer_ *ReferenceCollector
}

func NewReferenceCollector_CollectorCallback(outer_ *ReferenceCollector) *ReferenceCollector_CollectorCallback {
	return &ReferenceCollector_CollectorCallback{outer_: outer_}
}
func (cb *ReferenceCollector_CollectorCallback) IsJavaNil_() bool { return cb == nil }
func (cb *ReferenceCollector_CollectorCallback) ShouldTraverse_3(t *NodeTraversal, n, p pkg_rhino.Node) bool {
	return true
}
func (cb *ReferenceCollector_CollectorCallback) EnterScope(t *NodeTraversal) {
	cb.outer_.blockStack.AddLast(NewBasicBlock(cb.outer_.blockStack.PeekLast(), t.GetScopeRoot()))
}
func (cb *ReferenceCollector_CollectorCallback) ExitScope(t *NodeTraversal) {
	cb.outer_.behavior.AfterExitScope(t, cb.outer_)
	cb.outer_.blockStack.RemoveLast()
}
func (cb *ReferenceCollector_CollectorCallback) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) {
	sc := t.GetScope()
	if v := sc.vars.Get("x"); v != nil {
		cb.outer_.addRef(v, n, t)
	}
}

// --- InlineVariables ---

type InlineVariables struct {
	compiler      AbstractCompiler
	inlinedCount  int
	retriedCount  int
	idHash_       atomic.Int32
}

var (
	static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS *InlineVariables_InlineVarAnalysis
	static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT   *InlineVariables_VarExpert
	clinitDone_InlineVariables                                bool
)

func clinit_InlineVariables() {
	if !clinitDone_InlineVariables {
		clinitDone_InlineVariables = true
		clinitSlow_InlineVariables()
	}
}

func clinitSlow_InlineVariables() {
	static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS = NewInlineVariables_InlineVarAnalysis(int8(0), nil, nil, 0, nil, nil)
	static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT = NewInlineVariables_VarExpert_1(static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS)
}

func StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS() *InlineVariables_InlineVarAnalysis {
	clinit_InlineVariables()
	return static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS
}

func StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT() *InlineVariables_VarExpert {
	clinit_InlineVariables()
	return static_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT
}

func (this_ *InlineVariables) initFields_InlineVariables() {}
func (this_ *InlineVariables) IsJavaNil_() bool            { return this_ == nil }

func NewInlineVariables(compiler AbstractCompiler) *InlineVariables {
	clinit_InlineVariables()
	return &InlineVariables{compiler: compiler}
}

func (this_ *InlineVariables) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	clinit_ReferenceCollector()
	var val_1 ReferenceCollector
	InitReferenceCollector_3(&val_1, this_.compiler, NewInlineVariables_InliningBehavior(this_), NewSyntacticScopeCreator_1(this_.compiler))
	var callback *ReferenceCollector = &val_1
	_ = callback
	callback.Process_2(externs, root)
}

type InlineVariables_InlineVarAnalysis struct {
	kind                int8
	aliasedVar          *Var
	expert              *InlineVariables_VarExpert
	initValue           pkg_rhino.Node
	initialization      *Reference
	singleReadReference *Reference
	idHash_             atomic.Int32
}

func (a *InlineVariables_InlineVarAnalysis) initFields_InlineVariables_InlineVarAnalysis() {}

func NewInlineVariables_InlineVarAnalysis(kind int8, aliasedVar *Var, expert *InlineVariables_VarExpert, initValue pkg_rhino.Node, initialization *Reference, singleReadReference *Reference) *InlineVariables_InlineVarAnalysis {
	this_ := &InlineVariables_InlineVarAnalysis{}
	InitInlineVariables_InlineVarAnalysis(this_, kind, aliasedVar, expert, initValue, initialization, singleReadReference)
	return this_
}

func InitInlineVariables_InlineVarAnalysis(this_ *InlineVariables_InlineVarAnalysis, kind int8, aliasedVar *Var, expert *InlineVariables_VarExpert, initValue pkg_rhino.Node, initialization *Reference, singleReadReference *Reference) {
	this_.initFields_InlineVariables_InlineVarAnalysis()
	this_.kind = kind
	this_.aliasedVar = aliasedVar
	this_.expert = expert
	this_.initValue = initValue
	this_.initialization = initialization
	this_.singleReadReference = singleReadReference
}

func InlineVariables_InlineVarAnalysis_CreateWaitForAlias(aliasedVar *Var) *InlineVariables_InlineVarAnalysis {
	return NewInlineVariables_InlineVarAnalysis(int8(3), aliasedVar, nil, 0, nil, nil)
}

func InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(expert *InlineVariables_VarExpert, initValue pkg_rhino.Node) *InlineVariables_InlineVarAnalysis {
	return NewInlineVariables_InlineVarAnalysis(int8(5), nil, expert, initValue, nil, nil)
}

func InlineVariables_InlineVarAnalysis_CreateInlineSingleRead(expert *InlineVariables_VarExpert, initialization, singleReadReference *Reference) *InlineVariables_InlineVarAnalysis {
	return NewInlineVariables_InlineVarAnalysis(int8(6), nil, expert, 0, initialization, singleReadReference)
}

type InlineVariables_VarExpert struct {
	fixedAnalysis *InlineVariables_InlineVarAnalysis
	behavior      *InlineVariables_InliningBehavior
	v             *Var
	referenceInfo *ReferenceCollection
	declaration   *Reference
	flagConst     bool
	flagArgs      bool
	idHash_       atomic.Int32
}

func (e *InlineVariables_VarExpert) initFields_InlineVariables_VarExpert() {}

func NewInlineVariables_VarExpert_1(fixedAnalysis *InlineVariables_InlineVarAnalysis) *InlineVariables_VarExpert {
	this_ := &InlineVariables_VarExpert{}
	InitInlineVariables_VarExpert_1(this_, fixedAnalysis)
	return this_
}

func InitInlineVariables_VarExpert_1(this_ *InlineVariables_VarExpert, fixedAnalysis *InlineVariables_InlineVarAnalysis) {
	this_.initFields_InlineVariables_VarExpert()
	this_.fixedAnalysis = fixedAnalysis
}

func NewInlineVariables_VarExpert_5(behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, flagConst bool, flagArgs bool) *InlineVariables_VarExpert {
	this_ := &InlineVariables_VarExpert{}
	InitInlineVariables_VarExpert_5(this_, behavior, v, referenceInfo, flagConst, flagArgs)
	return this_
}

func InitInlineVariables_VarExpert_5(this_ *InlineVariables_VarExpert, behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, flagConst bool, flagArgs bool) {
	this_.initFields_InlineVariables_VarExpert()
	this_.behavior = behavior
	this_.v = v
	this_.referenceInfo = referenceInfo
	this_.declaration = referenceInfo.Get(0)
	this_.flagConst = flagConst
	this_.flagArgs = flagArgs
}

func (this_ *InlineVariables_VarExpert) Analyze() *InlineVariables_InlineVarAnalysis {
	if this_.fixedAnalysis != nil {
		return this_.fixedAnalysis
	}
	if this_.referenceInfo.Size() == 1 {
		return InlineVariables_InlineVarAnalysis_CreateWaitForAlias(this_.v)
	}
	if this_.referenceInfo.Size() == 2 {
		return InlineVariables_InlineVarAnalysis_CreateInlineSingleRead(this_, this_.declaration, this_.referenceInfo.Get(1))
	}
	return InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(this_, this_.declaration.node)
}

func (this_ *InlineVariables_VarExpert) Reanalyze(aliased *Var) *InlineVariables_InlineVarAnalysis {
	_ = aliased
	return InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(this_, this_.declaration.node)
}

type InlineVariables_InliningBehavior struct {
	outer_ *InlineVariables
}

func NewInlineVariables_InliningBehavior(outer_ *InlineVariables) *InlineVariables_InliningBehavior {
	return &InlineVariables_InliningBehavior{outer_: outer_}
}

func (this_ *InlineVariables_InliningBehavior) IsJavaNil_() bool { return this_ == nil }

func (this_ *InlineVariables_InliningBehavior) createExpert(v *Var, col *ReferenceCollection) *InlineVariables_VarExpert {
	if col == nil || col.Size() == 0 {
		return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT()
	}
	return NewInlineVariables_VarExpert_5(this_, v, col, true, false)
}

func (this_ *InlineVariables_InliningBehavior) AfterExitScope(t *NodeTraversal, rc *ReferenceCollector) {
	_ = t
	for _, v := range rc.GetAllSymbols() {
		exp := this_.createExpert(v, rc.GetReferences(v))
		ana := exp.Analyze()
		if ana.kind == 3 {
			rh := NewInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_, v, exp)
			rh.Complete(ana.aliasedVar)
		} else if ana.kind >= 4 {
			this_.outer_.inlinedCount++
		}
	}
}

type InlineVariables_InliningBehavior_AliasInlineRetryHandler struct {
	outer_ *InlineVariables_InliningBehavior
	v      *Var
	expert *InlineVariables_VarExpert
}

func (r *InlineVariables_InliningBehavior_AliasInlineRetryHandler) initFields_InlineVariables_InliningBehavior_AliasInlineRetryHandler() {}

func NewInlineVariables_InliningBehavior_AliasInlineRetryHandler(outer_ *InlineVariables_InliningBehavior, v *Var, expert *InlineVariables_VarExpert) *InlineVariables_InliningBehavior_AliasInlineRetryHandler {
	this_ := &InlineVariables_InliningBehavior_AliasInlineRetryHandler{}
	this_.outer_ = outer_
	InitInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_, v, expert)
	return this_
}

func InitInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler, v *Var, expert *InlineVariables_VarExpert) {
	this_.initFields_InlineVariables_InliningBehavior_AliasInlineRetryHandler()
	this_.v = v
	this_.expert = expert
}

func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) Complete(aliased *Var) {
	ana := this_.expert.Reanalyze(aliased)
	if ana.kind >= 4 {
		this_.outer_.outer_.retriedCount++
	}
}
`

	jscompGuardedSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w22ctest/gen/j2gort"
	pkg_rhino "example.com/w22ctest/gen/rhino"
)

type GuardedCallback interface {
	NodeTraversal_Callback
	VisitGuarded(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type GuardedCallbackBase struct {
	self_        GuardedCallback
	compiler     AbstractCompiler
	contextStack j2gort.Deque[*GuardedCallback_Context]
	idHash_      atomic.Int32
}

func (this_ *GuardedCallbackBase) BindSelf_(s any) { this_.self_ = s.(GuardedCallback) }
func (this_ *GuardedCallbackBase) IsJavaNil_() bool { return this_ == nil }
func InitGuardedCallback(this_ *GuardedCallbackBase, compiler AbstractCompiler) {
	this_.compiler = compiler
	this_.contextStack = j2gort.Deque[*GuardedCallback_Context]{}
}

func (this_ *GuardedCallbackBase) ShouldTraverse_3(traversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = traversal
	if parent == 0 {
		this_.contextStack.Push(StaticGet_GuardedCallback_Context_EMPTY())
	} else {
		this_.contextStack.Push(this_.contextStack.Peek().Descend(this_.compiler, parent, n))
	}
	return true
}

func (this_ *GuardedCallbackBase) Visit_3(traversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	this_.self_.VisitGuarded(traversal, n, parent)
	this_.contextStack.Pop()
}

func (this_ *GuardedCallbackBase) CurrentContextDepthAndSafe() (int, bool) {
	ctx := this_.contextStack.Peek()
	safe := ctx.safe
	depth := 0
	for c := ctx; c != nil && c.conditional != 0; c = c.linked {
		depth++
	}
	return depth, safe
}

type GuardedCallback_Context struct {
	conditional pkg_rhino.Node
	safe        bool
	linked      *GuardedCallback_Context
	idHash_     atomic.Int32
}

var static_GuardedCallback_Context_EMPTY *GuardedCallback_Context
var clinitDone_GuardedContext bool

func clinit_GuardedCallback_Context() {
	if !clinitDone_GuardedContext {
		clinitDone_GuardedContext = true
		clinitSlow_GuardedCallback_Context()
	}
}

func clinitSlow_GuardedCallback_Context() {
	static_GuardedCallback_Context_EMPTY = NewGuardedCallback_Context(0, false, nil)
}

func StaticGet_GuardedCallback_Context_EMPTY() *GuardedCallback_Context {
	clinit_GuardedCallback_Context()
	return static_GuardedCallback_Context_EMPTY
}

func (c *GuardedCallback_Context) initFields_GuardedCallback_Context() {}

func NewGuardedCallback_Context(conditional pkg_rhino.Node, safe bool, linked *GuardedCallback_Context) *GuardedCallback_Context {
	clinit_GuardedCallback_Context()
	this_ := &GuardedCallback_Context{}
	InitGuardedCallback_Context(this_, conditional, safe, linked)
	return this_
}

func InitGuardedCallback_Context(this_ *GuardedCallback_Context, conditional pkg_rhino.Node, safe bool, linked *GuardedCallback_Context) {
	this_.initFields_GuardedCallback_Context()
	this_.conditional = conditional
	this_.safe = safe
	this_.linked = linked
}

func (this_ *GuardedCallback_Context) Link(newConditional pkg_rhino.Node, newSafe bool) *GuardedCallback_Context {
	return NewGuardedCallback_Context(newConditional, newSafe, func() *GuardedCallback_Context {
		if this_.conditional != 0 {
			return this_
		}
		return nil
	}())
}

func (this_ *GuardedCallback_Context) Propagate(newSafe bool) *GuardedCallback_Context {
	if newSafe == this_.safe {
		return this_
	}
	return NewGuardedCallback_Context(this_.conditional, newSafe, this_.linked)
}

func (this_ *GuardedCallback_Context) Descend(compiler AbstractCompiler, parent pkg_rhino.Node, child pkg_rhino.Node) *GuardedCallback_Context {
	_ = compiler
	first := child == parent.GetFirstChild()
	switch parent.TokenSwitchOrd_() {
	case pkg_rhino.TokIF:
		if first {
			return this_.Link(parent, true)
		}
		return StaticGet_GuardedCallback_Context_EMPTY()
	case pkg_rhino.TokAND:
		if first {
			return this_.Link(parent, true)
		}
		return this_
	case pkg_rhino.TokOR:
		if first {
			return this_.Link(parent, false)
		}
		return this_
	default:
		return this_.Propagate(false)
	}
}

type PolyfillTraverser struct {
	GuardedCallbackBase
	guardedSafeHits int
	maxChainDepth   int
}

func NewPolyfillTraverser(compiler AbstractCompiler) *PolyfillTraverser {
	t := &PolyfillTraverser{}
	t.BindSelf_(t)
	InitGuardedCallback(&t.GuardedCallbackBase, compiler)
	return t
}

func (this_ *PolyfillTraverser) VisitGuarded(t *NodeTraversal, n, parent pkg_rhino.Node) {
	_, _ = t, parent
	d, safe := this_.CurrentContextDepthAndSafe()
	if safe {
		this_.guardedSafeHits++
	}
	if d > this_.maxChainDepth {
		this_.maxChainDepth = d
	}
}

type PeepScopeFreeCallback struct {
	visited int
}

func (p *PeepScopeFreeCallback) IsJavaNil_() bool { return p == nil }
func (p *PeepScopeFreeCallback) ShouldTraverse_3(t *NodeTraversal, n, par pkg_rhino.Node) bool {
	return true
}
func (p *PeepScopeFreeCallback) Visit_3(t *NodeTraversal, n, par pkg_rhino.Node) {
	_ = t.GetScopeRoot()
	p.visited++
}

type ScopeSavingCallback struct {
	saved []*Scope
}

func (s *ScopeSavingCallback) IsJavaNil_() bool { return s == nil }
func (s *ScopeSavingCallback) ShouldTraverse_3(t *NodeTraversal, n, par pkg_rhino.Node) bool {
	return true
}
func (s *ScopeSavingCallback) Visit_3(t *NodeTraversal, n, par pkg_rhino.Node) {
	s.saved = append(s.saved, t.GetScope())
}
`
	for name, content := range map[string]string{
		"node_traversal.go":   jscompNodeTraversalSrc,
		"inline_variables.go": jscompPassesSrc,
		"guarded_callback.go": jscompGuardedSrc,
	} {
		if err := os.WriteFile(filepath.Join(jscompDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	allPasses := passes.AllPasses()

	// Run 1 across parsing, unsafeparsing, and jscomp!
	for _, targetDir := range []string{parsingDir, unsafeParsingDir, jscompDir} {
		stats1, err := passes.RunOnPackageDir(targetDir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 1 on %s failed: %v", targetDir, err)
		}
		if stats1.TotalSites == 0 {
			t.Fatalf("expected > 0 sites on Run 1 for %s, got %+v", targetDir, stats1.ByPass)
		}
		// Run 2 MUST be 100% idempotent (0 sites, 0 modified files)!
		stats2, err := passes.RunOnPackageDir(targetDir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 2 on %s failed: %v", targetDir, err)
		}
		if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
			t.Fatalf("expected 0 sites & 0 modified files on Run 2 for %s, got %+v", targetDir, stats2)
		}
	}

	// Verify IRFactory vs JsDocInfoParser vs unsafeParsingDir lowering!
	irLoweredBytes, err := os.ReadFile(filepath.Join(parsingDir, "ir_factory.go"))
	if err != nil {
		t.Fatal(err)
	}
	irLowered := string(irLoweredBytes)
	for _, want := range []string{
		"this_.nodeAlloc_.AllocNode1Cached(",
		"this_.nodeAlloc_.AllocStringNodeCached(",
		"this_.nodeAlloc_.AllocNode1(",      // JsDocInfoParser MUST still pass explicit templateNode!
		"this_.nodeAlloc_.AllocStringNode(", // JsDocInfoParser MUST still pass explicit templateNode!
	} {
		if !strings.Contains(irLowered, want) {
			t.Fatalf("missing %q in lowered parsing/ir_factory.go:\n%s", want, irLowered)
		}
	}

	unsafeLoweredBytes, err := os.ReadFile(filepath.Join(unsafeParsingDir, "ir_factory.go"))
	if err != nil {
		t.Fatal(err)
	}
	unsafeLowered := string(unsafeLoweredBytes)
	if strings.Contains(unsafeLowered, "AllocNode1Cached(") || strings.Contains(unsafeLowered, "AllocStringNodeCached(") {
		t.Fatalf("UNSOUND: unsafeparsing aliasing &this_.templateNode must NOT use *Cached allocators:\n%s", unsafeLowered)
	}

	// Verify structural lowering in jscomp!
	ntLoweredBytes, err := os.ReadFile(filepath.Join(jscompDir, "node_traversal.go"))
	if err != nil {
		t.Fatal(err)
	}
	ntLowered := string(ntLoweredBytes)
	for _, want := range []string{
		"nodestore.BlockShift",
		"nodestore.BlockOf(",
		"curBlk_.FirstNextTok(",
		"traverseScopeRootPooled_(",
		"defer scopeCreator.ReleaseSlabs_()",
	} {
		if !strings.Contains(ntLowered, want) {
			t.Fatalf("missing %q in lowered jscomp/node_traversal.go:\n%s", want, ntLowered)
		}
	}
	if strings.Count(ntLowered, "nodestore.BlockOf(") < 4 {
		t.Fatalf("expected BlockOf in both traverseBranch and traverseChildren (outer + gc_ loops), got %d:\n%s",
			strings.Count(ntLowered, "nodestore.BlockOf("), ntLowered)
	}

	ivLoweredBytes, err := os.ReadFile(filepath.Join(jscompDir, "inline_variables.go"))
	if err != nil {
		t.Fatal(err)
	}
	ivLowered := string(ivLoweredBytes)
	for _, want := range []string{
		"j2gort.AcquireCustomOrderedMap[*Var, *ReferenceCollection]()",
		"j2gort.AcquireCustomOrderedMapCap[*Var, *ReferenceCollection](int32(1024))",
		"slab_InlineVariables_VarExpert",
		"slab_InlineVariables_InlineVarAnalysis",
		"slab_InlineVariables_InliningBehavior_AliasInlineRetryHandler",
		"func (this_ *InlineVariables) ReleaseSlabs_()",
		"defer this_.ReleaseSlabs_()",
	} {
		if !strings.Contains(ivLowered, want) {
			t.Fatalf("missing %q in lowered jscomp/inline_variables.go:\n%s", want, ivLowered)
		}
	}

	gcLoweredBytes, err := os.ReadFile(filepath.Join(jscompDir, "guarded_callback.go"))
	if err != nil {
		t.Fatal(err)
	}
	gcLowered := string(gcLoweredBytes)
	for _, want := range []string{
		"slab_GuardedCallback_Context",
		"j2gort.PooledSlab[GuardedCallback_Context]",
		"descendWithSlab_(&this_.slab_GuardedCallback_Context,",
		"if this_.contextStack.IsEmpty() {",
		"this_.slab_GuardedCallback_Context.Release()",
	} {
		if !strings.Contains(gcLowered, want) {
			t.Fatalf("missing %q in lowered jscomp/guarded_callback.go:\n%s", want, gcLowered)
		}
	}

	// Execute end-to-end behavioral & -race test on the lowered packages!
	rtTestSrc := `package w22c_rt_test

import (
	"sync"
	"testing"

	"example.com/w22ctest/gen/j2gort"
	"example.com/w22ctest/gen/jscomp"
	"example.com/w22ctest/gen/parsing"
	pkg_rhino "example.com/w22ctest/gen/rhino"
	"example.com/w22ctest/shims/nodestore"
)

type testComp struct{ opts *jscomp.CompilerOptions }
func (c *testComp) IsJavaNil_() bool                 { return c == nil }
func (c *testComp) GetOptions() *jscomp.CompilerOptions { return c.opts }

func TestWave22CRuntimeBehaviorAndRace(t *testing.T) {
	// 1. IRFactory AllocNode1Cached + AllocStringNodeCached vs JsDocInfoParser distinct templateNode!
	nodestore.ResetNodestoreMetrics()
	defTmpl := pkg_rhino.Node(nodestore.AllocAt(3, 1, 0))
	cuTmpl := pkg_rhino.Node(nodestore.AllocAt(3, 2, 0))
	jsdocTmpl := pkg_rhino.Node(nodestore.AllocAt(3, 3, 0))

	n1, s1, jn, s2 := parsing.IRFactory_TransformTree(defTmpl, cuTmpl, jsdocTmpl)
	if n1.TemplateOf() != defTmpl {
		t.Fatalf("n1 template = %v, want default %v", n1.TemplateOf(), defTmpl)
	}
	if s1.TemplateOf() != cuTmpl || s1.StringVal() != "hello" {
		t.Fatalf("s1 = (tmpl=%v, str=%q), want (%v, hello)", s1.TemplateOf(), s1.StringVal(), cuTmpl)
	}
	if jn.TemplateOf() != jsdocTmpl || jn.StringVal() != "docstr" {
		t.Fatalf("JsDocInfoParser node template = %v (str=%q), want jsdocTmpl %v", jn.TemplateOf(), jn.StringVal(), jsdocTmpl)
	}
	if s2.TemplateOf() != defTmpl || s2.StringVal() != "world" {
		t.Fatalf("s2 after restoring defaultTemplateNode = (tmpl=%v, str=%q), want (%v, world)", s2.TemplateOf(), s2.StringVal(), defTmpl)
	}
	if nodestore.AllocNode1CachedCalls.Load() != 1 || nodestore.AllocStringNodeCachedCalls.Load() != 2 {
		t.Fatalf("expected 1 AllocNode1Cached & 2 AllocStringNodeCached in IRFactory, got %d & %d",
			nodestore.AllocNode1CachedCalls.Load(), nodestore.AllocStringNodeCachedCalls.Load())
	}
	if nodestore.AllocStringNodeCalls.Load() != 1 {
		t.Fatalf("expected 1 non-cached AllocStringNode in JsDocInfoParser, got %d", nodestore.AllocStringNodeCalls.Load())
	}

	// 2. Block-cached child iteration across 200 siblings in Block 4 + cross-block jump to Block 5!
	script := nodestore.AllocAt(4, 1, uint8(pkg_rhino.TokSCRIPT)-1)
	blockNode := nodestore.AllocAt(4, 2, uint8(pkg_rhino.TokBLOCK)-1)
	nodestore.AppendChild(script, blockNode)
	for i := uint32(0); i < 100; i++ {
		c := nodestore.AllocAt(4, 10+i, uint8(pkg_rhino.TokNAME)-1)
		nodestore.AppendChild(blockNode, c)
	}
	for i := uint32(0); i < 50; i++ {
		c := nodestore.AllocAt(5, 10+i, uint8(pkg_rhino.TokNAME)-1)
		nodestore.AppendChild(blockNode, c)
	}

	nodestore.ResetNodestoreMetrics()
	comp := &testComp{opts: &jscomp.CompilerOptions{}}
	peepCb := &jscomp.PeepScopeFreeCallback{}
	jscomp.NodeTraversal_Traverse(comp, pkg_rhino.Node(script), peepCb)
	if peepCb.VisitedCount_() != 152 {
		t.Fatalf("expected 152 visited nodes (script + block + 150 leaves), got %d", peepCb.VisitedCount_())
	}
	// 100 siblings in block 4 + 50 siblings in block 5 = ONLY 2 BlockOf lookups for all 150 leaf siblings (+1 for blockNode under script = 3 total)!
	if gotBlk := nodestore.BlockOfCalls.Load(); gotBlk != 3 {
		t.Fatalf("expected exactly 3 BlockOf calls across 151 children (98%% block-cache hit rate), got %d", gotBlk)
	}
	if gotFast := nodestore.BlockFirstNextTokCalls.Load(); gotFast != 151 {
		t.Fatalf("expected 151 (*NodeBlock).FirstNextTok calls, got %d", gotFast)
	}

	// 3. NodeTraversal_TraverseScopeRoots must allocate 0 heap NodeTraversal objects and release MemoizedScopeCreator when scope-free!
	jscomp.NewNodeTraversalHeapCalls.Store(0)
	j2gort.ResetW22CMetrics()
	fn1 := nodestore.AllocAt(4, 500, uint8(pkg_rhino.TokFUNCTION)-1)
	fn2 := nodestore.AllocAt(4, 501, uint8(pkg_rhino.TokFUNCTION)-1)
	nodestore.AppendChild(script, fn1)
	nodestore.AppendChild(script, fn2)
	roots := []pkg_rhino.Node{pkg_rhino.Node(fn1), pkg_rhino.Node(fn2)}
	for iter := 0; iter < 20; iter++ {
		jscomp.NodeTraversal_TraverseScopeRoots(comp, roots, peepCb, false)
	}
	if heapNT := jscomp.NewNodeTraversalHeapCalls.Load(); heapNT != 0 {
		t.Fatalf("expected 0 NewNodeTraversal heap allocations in TraverseScopeRoots, got %d", heapNT)
	}
	if j2gort.OrderedMapReleases.Load() < 20 {
		t.Fatalf("expected MemoizedScopeCreator + Scope OrderedMaps released for scope-free callback, got %d", j2gort.OrderedMapReleases.Load())
	}

	// Escaping callback passed to NodeTraversal_TraverseScopeRoots MUST keep saved *Scope valid!
	savingCb := &jscomp.ScopeSavingCallback{}
	jscomp.NodeTraversal_TraverseScopeRoots(comp, roots, savingCb, true)
	if len(savingCb.SavedScopes_()) == 0 {
		t.Fatalf("expected ScopeSavingCallback to capture scopes")
	}
	savedRoot := savingCb.SavedScopes_()[0].GetRootNode_0()
	for iter := 0; iter < 15; iter++ {
		jscomp.NodeTraversal_TraverseScopeRoots(comp, roots, peepCb, false)
	}
	if savingCb.SavedScopes_()[0].GetRootNode_0() != savedRoot {
		t.Fatalf("escaped *Scope corrupted by subsequent TraverseScopeRoots!")
	}

	// 4. ReferenceCollector Tier-0 (Process_1) vs Tier-2 (Process_2) + InlineVariables PooledSlabs!
	j2gort.ResetW22CMetrics()
	externs := pkg_rhino.Node(nodestore.AllocAt(4, 600, uint8(pkg_rhino.TokSCRIPT)-1))
	rcLocal := jscomp.NewReferenceCollector_3(comp, jscomp.StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR(), jscomp.NewSyntacticScopeCreator_1(comp))
	rcLocal.Process_1(pkg_rhino.Node(fn1))
	if j2gort.CustomOrderedMapTier2CapAcq.Load() != 0 || j2gort.CustomOrderedMapTier0Acq.Load() != 1 {
		t.Fatalf("Process_1 must stay Tier-0 (tier0=%d, tier2=%d)",
			j2gort.CustomOrderedMapTier0Acq.Load(), j2gort.CustomOrderedMapTier2CapAcq.Load())
	}

	iv := jscomp.NewInlineVariables(comp)
	for iter := 0; iter < 30; iter++ {
		iv.Process_2(externs, pkg_rhino.Node(script))
	}
	if j2gort.CustomOrderedMapTier2CapAcq.Load() != 30 || j2gort.CustomOrderedMapLastCapHint.Load() != 1024 {
		t.Fatalf("expected 30 Tier-2 AcquireCustomOrderedMapCap(1024) calls in Process_2, got count=%d capHint=%d",
			j2gort.CustomOrderedMapTier2CapAcq.Load(), j2gort.CustomOrderedMapLastCapHint.Load())
	}
	// Static InlineVariables singletons MUST remain intact after 30 ReleaseSlabs_() cycles!
	if jscomp.StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT() == nil ||
		jscomp.StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS() == nil {
		t.Fatalf("static InlineVariables singletons nil!")
	}

	// 5. GuardedCallback_Context PooledSlab across deeply nested IF / AND / OR AST!
	ifNode := nodestore.AllocAt(4, 700, uint8(pkg_rhino.TokIF)-1)
	andNode := nodestore.AllocAt(4, 701, uint8(pkg_rhino.TokAND)-1)
	condLeaf := nodestore.AllocAt(4, 702, uint8(pkg_rhino.TokNAME)-1)
	rhsLeaf := nodestore.AllocAt(4, 703, uint8(pkg_rhino.TokNAME)-1)
	thenBlk := nodestore.AllocAt(4, 704, uint8(pkg_rhino.TokBLOCK)-1)
	nodestore.AppendChild(andNode, condLeaf)
	nodestore.AppendChild(andNode, rhsLeaf)
	nodestore.AppendChild(ifNode, andNode)
	nodestore.AppendChild(ifNode, thenBlk)

	j2gort.ResetW22CMetrics()
	for iter := 0; iter < 25; iter++ {
		pt := jscomp.NewPolyfillTraverser(comp)
		jscomp.NodeTraversal_Traverse(comp, pkg_rhino.Node(ifNode), pt)
		if pt.GuardedSafeHits_() == 0 || pt.MaxChainDepth_() < 2 {
			t.Fatalf("unexpected GuardedCallback context tracking: safeHits=%d maxDepth=%d",
				pt.GuardedSafeHits_(), pt.MaxChainDepth_())
		}
	}
	if rel := j2gort.PooledSlabReleaseCalls.Load(); rel < 25 {
		t.Fatalf("expected >= 25 GuardedCallback_Context slab releases, got %d", rel)
	}
	if jscomp.StaticGet_GuardedCallback_Context_EMPTY().IsSafe_() {
		t.Fatalf("EMPTY GuardedCallback_Context singleton corrupted!")
	}

	// 6. Concurrent -race stress across NodeTraversal + InlineVariables + GuardedCallback!
	var wg sync.WaitGroup
	for g := 0; g < 6; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := &testComp{opts: &jscomp.CompilerOptions{}}
			for i := 0; i < 20; i++ {
				localIV := jscomp.NewInlineVariables(c)
				localIV.Process_2(externs, pkg_rhino.Node(script))
				localPT := jscomp.NewPolyfillTraverser(c)
				jscomp.NodeTraversal_Traverse(c, pkg_rhino.Node(ifNode), localPT)
			}
		}()
	}
	wg.Wait()
}
`
	// Export small read-only test accessors in jscompDir so external test package can inspect counters.
	jscompTestExportSrc := `package jscomp

func (p *PeepScopeFreeCallback) VisitedCount_() int    { return p.visited }
func (s *ScopeSavingCallback) SavedScopes_() []*Scope  { return s.saved }
func (p *PolyfillTraverser) GuardedSafeHits_() int     { return p.guardedSafeHits }
func (p *PolyfillTraverser) MaxChainDepth_() int       { return p.maxChainDepth }
func (c *GuardedCallback_Context) IsSafe_() bool       { return c.safe }
`
	if err := os.WriteFile(filepath.Join(jscompDir, "test_exports.go"), []byte(jscompTestExportSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "wave22c_rt_test.go"), []byte(rtTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-race", "-v", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race in wave22c workspace failed: %v\n%s", err, string(out))
	}
}

func TestWave22C_RealGenDryRunSoundness(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/j2gort is currently mid-edit; skipping real gen/ DryRun")
	}

	realJscompDir := filepath.Join(repoRoot, "transpiled", "gen", "jscomp")
	ntSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "node_traversal.go"))
	mdnuSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "make_declared_names_unique.go"))
	normSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "normalize.go"))
	crmSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "closure_rewrite_module.go"))
	fsivSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "flow_sensitive_inline_variables.go"))
	nsNodeSrc, _ := os.ReadFile(filepath.Join(repoRoot, "transpiled", "shims", "nodestore", "node.go"))
	hasBlockOf := strings.Contains(string(nsNodeSrc), "func BlockOf(")
	ntPreLowered := (!hasBlockOf || strings.Contains(string(ntSrc), "BlockOf(")) &&
		strings.Contains(string(ntSrc), "traverseScopeRootPooled_") &&
		strings.Contains(string(ntSrc), "unwrapNodeTraversalCallback_") &&
		strings.Contains(string(ntSrc), "canReleaseSyntacticScopesForCallback_") &&
		strings.Contains(string(mdnuSrc), "CanReleaseSyntacticScopes_") &&
		strings.Contains(string(normSrc), "func (*Normalize_NormalizeStatements) CanReleaseSyntacticScopes_()") &&
		strings.Contains(string(crmSrc), "func (*ClosureRewriteModule_ScriptPreprocessor) CanReleaseSyntacticScopes_()") &&
		strings.Contains(string(fsivSrc), "func (*FlowSensitiveInlineVariables) CanReleaseSyntacticScopes_()")

	dryStats, err := passes.RunOnPackageDir(
		realJscompDir,
		[]passes.Pass{passes.NewSlabAllocLoweringPass(), passes.NewNodeTraversalLoweringPass()},
		passes.Options{DryRun: true},
	)
	if err != nil {
		t.Fatalf("DryRun slab_alloc_lowering + node_traversal_lowering on real gen/jscomp failed: %v", err)
	}
	if !ntPreLowered && dryStats.ByPass["node_traversal_lowering"] < 1 {
		t.Fatalf("expected >= 1 node_traversal_lowering sites on real gen/jscomp, got %+v", dryStats.ByPass)
	}
	if ntPreLowered && dryStats.ByPass["node_traversal_lowering"] != 0 {
		t.Fatalf("expected 0 node_traversal_lowering sites on already-lowered gen/jscomp, got %+v", dryStats.ByPass)
	}
}


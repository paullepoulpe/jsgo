package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave23C_EndToEndLeaksAndHotspots(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := "module example.com/w23ctest\n\ngo 1.24\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	j2Dir := filepath.Join(tmpDir, "gen", "j2gort")
	rhinoDir := filepath.Join(tmpDir, "gen", "rhino")
	graphDir := filepath.Join(tmpDir, "gen", "graph")
	jscompDir := filepath.Join(tmpDir, "gen", "jscomp")
	for _, d := range []string{j2Dir, rhinoDir, graphDir, jscompDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	j2Src := `package j2gort

import (
	"sync"
	"sync/atomic"
)

const NullString = "\x00"

type ClinitState struct{}

var (
	SlabReleaseCalls       atomic.Int64
	SlabResetCalls         atomic.Int64
	SlabReleaseWithCalls   atomic.Int64
	OrderedMapReleases     atomic.Int64
	CustomOrderedMapPutOps atomic.Int64
	CustomOrderedMapRelOps atomic.Int64
)

func ResetMetrics() {
	SlabReleaseCalls.Store(0)
	SlabResetCalls.Store(0)
	SlabReleaseWithCalls.Store(0)
	OrderedMapReleases.Store(0)
	CustomOrderedMapPutOps.Store(0)
	CustomOrderedMapRelOps.Store(0)
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
	p, _ := v.(*T)
	return p
}

func AsIface[T any](v any) T {
	var z T
	if v == nil {
		return z
	}
	r, _ := v.(T)
	return r
}

func TryCast[T any](v any) (T, bool) {
	var z T
	if v == nil {
		return z, false
	}
	r, ok := v.(T)
	return r, ok
}

func CheckNotNull(v any) any {
	if v == nil {
		panic("NPE")
	}
	return v
}

func CheckNotNullPtr[T any](p *T) *T {
	if p == nil {
		panic("NPE")
	}
	return p
}

type Slab[T any] struct{ buf []T }

func (s *Slab[T]) Alloc() *T {
	s.buf = append(s.buf, *new(T))
	return &s.buf[len(s.buf)-1]
}
func (s *Slab[T]) Release() {
	if s != nil {
		clear(s.buf)
		s.buf = nil
	}
}

var slabPools sync.Map

func poolFor[T any]() *sync.Pool {
	var zero *T
	key := any(zero)
	if p, ok := slabPools.Load(key); ok {
		return p.(*sync.Pool)
	}
	np := &sync.Pool{}
	actual, _ := slabPools.LoadOrStore(key, np)
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
		if v := poolFor[T]().Get(); v != nil {
			p := v.(*[]T)
			c = (*p)[:1]
		} else {
			c = make([]T, 1, 16)
		}
		s.buf = c
		return &c[0]
	}
	b := s.buf[:i+1]
	s.buf = b
	return &b[i]
}

func (s *PooledSlab[T]) Reset() {
	if s == nil {
		return
	}
	SlabResetCalls.Add(1)
	pool := poolFor[T]()
	for idx := range s.chunks {
		ch := s.chunks[idx]
		clear(ch)
		ch = ch[:0]
		pool.Put(&ch)
		s.chunks[idx] = nil
	}
	s.chunks = s.chunks[:0]
	if len(s.buf) > 0 {
		clear(s.buf)
		s.buf = s.buf[:0]
	}
}

func (s *PooledSlab[T]) Release() {
	if s == nil {
		return
	}
	SlabReleaseCalls.Add(1)
	pool := poolFor[T]()
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
	SlabReleaseWithCalls.Add(1)
	pool := poolFor[T]()
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

type OrderedMapValCursor[V any] struct {
	vals []V
	pos  int
}

func (c *OrderedMapValCursor[V]) HasNext() bool { return c.pos < len(c.vals) }
func (c *OrderedMapValCursor[V]) Next() V {
	v := c.vals[c.pos]
	c.pos++
	return v
}

type OrderedMap[K comparable, V any] struct {
	keys []K
	vals []V
	idx  map[K]int
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V]     { return OrderedMap[K, V]{} }
func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] { return OrderedMap[K, V]{} }
func NewOrderedMapPtrCap[K comparable, V any](_ int32) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}
func AcquireOrderedMapPtrCap[K comparable, V any](_ int32) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}
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
func (m *OrderedMap[K, V]) IterValues() OrderedMapValCursor[V] {
	if m == nil {
		return OrderedMapValCursor[V]{}
	}
	return OrderedMapValCursor[V]{vals: m.vals}
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
	keys []K
	vals []V
}

func NewCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{}
}
func AcquireCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{}
}
func AcquireCustomOrderedMapCap[K FastCustomKeyComparable, V any](capHint int32) CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{
		keys: make([]K, 0, capHint),
		vals: make([]V, 0, capHint),
	}
}
func (m *CustomOrderedMap[K, V]) Size() int32 {
	if m == nil {
		return 0
	}
	return int32(len(m.keys))
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
	CustomOrderedMapPutOps.Add(1)
	for i, ek := range m.keys {
		if ek == k || (!ek.IsJavaNil_() && ek.Equals(k)) {
			m.vals[i] = v
			return
		}
	}
	m.keys = append(m.keys, k)
	m.vals = append(m.vals, v)
}
func (m *CustomOrderedMap[K, V]) ReleaseToPool() {
	if m != nil {
		CustomOrderedMapRelOps.Add(1)
		*m = CustomOrderedMap[K, V]{}
	}
}

type CustomOrderedSet[T FastCustomKeyComparable] struct {
	items []T
}

func NewCustomOrderedSet[T FastCustomKeyComparable]() CustomOrderedSet[T] {
	return CustomOrderedSet[T]{}
}
func (s *CustomOrderedSet[T]) Add(v T) bool {
	if s.Contains(v) {
		return false
	}
	s.items = append(s.items, v)
	return true
}
func (s *CustomOrderedSet[T]) Contains(v T) bool {
	for _, x := range s.items {
		if x == v || (!x.IsJavaNil_() && x.Equals(v)) {
			return true
		}
	}
	return false
}

type OrderedSet[T comparable] struct {
	items []T
}

func NewOrderedSet[T comparable]() OrderedSet[T] { return OrderedSet[T]{} }
func (s *OrderedSet[T]) Add(v T) bool {
	for _, x := range s.items {
		if x == v {
			return false
		}
	}
	s.items = append(s.items, v)
	return true
}

type Deque[T any] struct {
	items []T
}

func (d *Deque[T]) Push(v T)     { d.items = append(d.items, v) }
func (d *Deque[T]) AddFirst(v T) { d.items = append(d.items, v) }
func (d *Deque[T]) Peek() T {
	var z T
	if len(d.items) == 0 {
		return z
	}
	return d.items[len(d.items)-1]
}
func (d *Deque[T]) PeekFirst() T { return d.Peek() }
func (d *Deque[T]) Pop() T {
	n := len(d.items) - 1
	v := d.items[n]
	var z T
	d.items[n] = z
	d.items = d.items[:n]
	return v
}
func (d *Deque[T]) RemoveFirst() T { return d.Pop() }
func (d *Deque[T]) IsEmpty() bool  { return len(d.items) == 0 }

type Ext_java_lang_Iterable interface {
	SliceItems_() []any
}

type Ext_java_util_List interface {
	Ext_java_lang_Iterable
}

type SliceIterable struct {
	Items []any
}

func (s SliceIterable) SliceItems_() []any { return s.Items }

type IterCursor struct {
	items []any
	pos   int
}

func NewIterCursor(it any) IterCursor {
	if s, ok := it.(Ext_java_lang_Iterable); ok && s != nil {
		return IterCursor{items: s.SliceItems_()}
	}
	return IterCursor{}
}
func (c *IterCursor) HasNext() bool { return c.pos < len(c.items) }
func (c *IterCursor) Next() any {
	v := c.items[c.pos]
	c.pos++
	return v
}
func (c *IterCursor) NextNode() uint32 {
	v := c.Next()
	return v.(uint32)
}

type Ext_com_google_common_base_Predicate interface {
	Apply(any) bool
}

type predFunc struct{ fn func(any) bool }

func (p *predFunc) Apply(v any) bool { return p.fn(v) }

func NewLambda_Ext_com_google_common_base_Predicate(fn func(any) bool) Ext_com_google_common_base_Predicate {
	return &predFunc{fn: fn}
}

func Ext_com_google_common_base_Predicates_AlwaysTrue() Ext_com_google_common_base_Predicate {
	return &predFunc{fn: func(any) bool { return true }}
}

type Ext_java_util_function_Consumer interface {
	Accept(any)
}
`
	if err := os.WriteFile(filepath.Join(j2Dir, "runtime.go"), []byte(j2Src), 0o644); err != nil {
		t.Fatal(err)
	}

	rhinoSrc := `package rhino

const (
	TokSCRIPT    int32 = 1
	TokFUNCTION  int32 = 2
	TokBLOCK     int32 = 3
	TokVAR       int32 = 4
	TokLET       int32 = 5
	TokCONST     int32 = 6
	TokFOR       int32 = 7
	TokOBJECTLIT int32 = 8
	TokNAME      int32 = 9
	TokNUMBER    int32 = 10
	TokIF        int32 = 11
	TokAND       int32 = 12
)

type Node uint32

type nodeRec struct {
	tok    int32
	str    string
	parent Node
	first  Node
	last   Node
	next   Node
}

var arena = []nodeRec{{}}

func ResetArena() {
	arena = arena[:1]
}

func AllocNode(tok int32, str string) Node {
	id := Node(len(arena))
	arena = append(arena, nodeRec{tok: tok, str: str})
	return id
}

func AppendChild(parent, child Node) {
	arena[child].parent = parent
	if arena[parent].first == 0 {
		arena[parent].first = child
		arena[parent].last = child
	} else {
		arena[arena[parent].last].next = child
		arena[parent].last = child
	}
}

func (n Node) IsJavaNil_() bool    { return n == 0 }
func (n Node) GetToken() int32     { return arena[n].tok }
func (n Node) GetString() string   { return arena[n].str }
func (n Node) GetParent() Node     { return arena[n].parent }
func (n Node) GetFirstChild() Node { return arena[n].first }
func (n Node) GetLastChild() Node  { return arena[n].last }
func (n Node) GetNext() Node       { return arena[n].next }
func (n Node) IsScript() bool      { return n != 0 && arena[n].tok == TokSCRIPT }
func (n Node) IsFunction() bool    { return n != 0 && arena[n].tok == TokFUNCTION }
func (n Node) IsBlock() bool       { return n != 0 && arena[n].tok == TokBLOCK }
func (n Node) IsObjectLit() bool   { return n != 0 && arena[n].tok == TokOBJECTLIT }
func (n Node) IsName() bool        { return n != 0 && arena[n].tok == TokNAME }
func (n Node) IsIf() bool          { return n != 0 && arena[n].tok == TokIF }
func (n Node) IsAnd() bool         { return n != 0 && arena[n].tok == TokAND }
`
	if err := os.WriteFile(filepath.Join(rhinoDir, "node.go"), []byte(rhinoSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	graphSrc := `package graph

import "example.com/w23ctest/gen/j2gort"

type Graph[N any, E any] interface {
	Implements_graph_Graph()
	IsJavaNil_() bool
}

type UndiGraph[N comparable, E any] interface {
	Implements_graph_UndiGraph()
	Graph[N, E]
	CreateUndirectedGraphNode(N)
	Connect(N, E, N)
}

type LinkedDirectedGraph_LinkedDiGraphNode[N any, E any] struct {
	val N
}

func NewLinkedDirectedGraph_LinkedDiGraphNode[N any, E any](v N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	this_ := &LinkedDirectedGraph_LinkedDiGraphNode[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphNode[N, E](this_, v)
	return this_
}
func InitLinkedDirectedGraph_LinkedDiGraphNode[N any, E any](this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E], v N) {
	this_.val = v
}

type LinkedDirectedGraph_LinkedDiGraphEdge[N any, E any] struct {
	tag E
}

func NewLinkedDirectedGraph_LinkedDiGraphEdge[N any, E any](e E) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	this_ := &LinkedDirectedGraph_LinkedDiGraphEdge[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphEdge[N, E](this_, e)
	return this_
}
func InitLinkedDirectedGraph_LinkedDiGraphEdge[N any, E any](this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E], e E) {
	this_.tag = e
}

type LinkedDirectedGraphBase[N comparable, E any] struct {
	nodes j2gort.OrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]
}

func (*LinkedDirectedGraphBase[N, E]) Implements_graph_Graph() {}
func (g *LinkedDirectedGraphBase[N, E]) IsJavaNil_() bool      { return g == nil }

func (this_ *LinkedDirectedGraphBase[N, E]) initFields_LinkedDirectedGraph() {
	this_.nodes = j2gort.NewOrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]()
}

func InitLinkedDirectedGraph[N comparable, E any](this_ *LinkedDirectedGraphBase[N, E]) {
	this_.initFields_LinkedDirectedGraph()
}

func (this_ *LinkedDirectedGraphBase[N, E]) CreateNode(v N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	n := NewLinkedDirectedGraph_LinkedDiGraphNode[N, E](v)
	this_.nodes.Put(v, n)
	return n
}

func (this_ *LinkedDirectedGraphBase[N, E]) CreateEdge(e E) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	return NewLinkedDirectedGraph_LinkedDiGraphEdge[N, E](e)
}

type LinkedUndirectedGraph_LinkedUndirectedGraphNode[N comparable, E any] struct {
	value N
}

func NewLinkedUndirectedGraph_LinkedUndirectedGraphNode[N comparable, E any](v N) *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] {
	this_ := &LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]{}
	InitLinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E](this_, v)
	return this_
}
func InitLinkedUndirectedGraph_LinkedUndirectedGraphNode[N comparable, E any](this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E], v N) {
	this_.value = v
}

type LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N comparable, E any] struct {
	value E
}

func NewLinkedUndirectedGraph_LinkedUndirectedGraphEdge[N comparable, E any](e E) *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E] {
	this_ := &LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]{}
	InitLinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E](this_, e)
	return this_
}
func InitLinkedUndirectedGraph_LinkedUndirectedGraphEdge[N comparable, E any](this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E], e E) {
	this_.value = e
}

type LinkedUndirectedGraph[N comparable, E any] struct {
	nodes j2gort.OrderedMap[N, *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]]
}

func (*LinkedUndirectedGraph[N, E]) Implements_graph_Graph()     {}
func (*LinkedUndirectedGraph[N, E]) Implements_graph_UndiGraph() {}
func (g *LinkedUndirectedGraph[N, E]) IsJavaNil_() bool          { return g == nil }

func (this_ *LinkedUndirectedGraph[N, E]) initFields_LinkedUndirectedGraph() {
	this_.nodes = j2gort.NewOrderedMap[N, *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]]()
}

func NewLinkedUndirectedGraph[N comparable, E any]() *LinkedUndirectedGraph[N, E] {
	this_ := &LinkedUndirectedGraph[N, E]{}
	this_.initFields_LinkedUndirectedGraph()
	return this_
}

func (this_ *LinkedUndirectedGraph[N, E]) CreateUndirectedGraphNode(v N) {
	if this_.nodes.Get(v) == nil {
		n := NewLinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E](v)
		this_.nodes.Put(v, n)
	}
}

func (this_ *LinkedUndirectedGraph[N, E]) Connect(a N, e E, b N) {
	this_.CreateUndirectedGraphNode(a)
	this_.CreateUndirectedGraphNode(b)
	_ = NewLinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E](e)
}

type GraphColoring[N comparable, E any] interface {
	IsJavaNil_() bool
	GetGraph() UndiGraph[N, E]
}

type GreedyGraphColoring[N comparable, E any] struct {
	g UndiGraph[N, E]
}

func NewGreedyGraphColoring[N comparable, E any](g UndiGraph[N, E]) GraphColoring[N, E] {
	return &GreedyGraphColoring[N, E]{g: g}
}
func (c *GreedyGraphColoring[N, E]) IsJavaNil_() bool         { return c == nil }
func (c *GreedyGraphColoring[N, E]) GetGraph() UndiGraph[N, E] { return c.g }
`
	if err := os.WriteFile(filepath.Join(graphDir, "graph.go"), []byte(graphSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	if _, err := passes.RunOnPackageDir(graphDir, allPasses, passes.Options{}); err != nil {
		t.Fatalf("graph Run 1 failed: %v", err)
	}
	gStats2, err := passes.RunOnPackageDir(graphDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("graph Run 2 failed: %v", err)
	}
	if gStats2.TotalSites != 0 || gStats2.FilesModified != 0 {
		t.Fatalf("expected 0 edits on graph Run 2, got %+v", gStats2)
	}

	graphLoweredBytes, err := os.ReadFile(filepath.Join(graphDir, "graph.go"))
	if err != nil {
		t.Fatal(err)
	}
	graphLowered := string(graphLoweredBytes)
	for _, want := range []string{
		"slab_LinkedUndirectedGraph_LinkedUndirectedGraphNode j2gort.PooledSlab[LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]]",
		"slab_LinkedUndirectedGraph_LinkedUndirectedGraphEdge j2gort.PooledSlab[LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]]",
		"func (this_ *LinkedUndirectedGraph[N, E]) ReleaseSlabs_()",
		"this_.nodes = j2gort.AcquireOrderedMap[N, *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]]()",
	} {
		if !strings.Contains(graphLowered, want) {
			t.Fatalf("missing %q in lowered graph.go:\n%s", want, graphLowered)
		}
	}

	// Now create jscomp files exercising all 4 Wave 23C fixes!
	jscompCoreSrc := `package jscomp

import (
	"example.com/w23ctest/gen/j2gort"
	pkg_graph "example.com/w23ctest/gen/graph"
	pkg_rhino "example.com/w23ctest/gen/rhino"
)

type CodingConvention interface {
	Implements_jscomp_CodingConvention()
	IsExported_2(string, bool) bool
	IsPropertyRenameFunction(pkg_rhino.Node) bool
}

type DefaultCodingConvention struct{}

func (*DefaultCodingConvention) Implements_jscomp_CodingConvention()          {}
func (*DefaultCodingConvention) IsExported_2(name string, _ bool) bool        { return len(name) > 0 && name[0] == '_' }
func (*DefaultCodingConvention) IsPropertyRenameFunction(pkg_rhino.Node) bool { return false }

type AbstractCompiler interface {
	Implements_jscomp_AbstractCompiler()
	IsJavaNil_() bool
	GetCodingConvention() CodingConvention
	HasScopeChanged(pkg_rhino.Node) bool
}

type Compiler struct {
	jsRoot  pkg_rhino.Node
	conv    CodingConvention
	changed bool
}

func NewCompiler(jsRoot pkg_rhino.Node) *Compiler {
	return &Compiler{jsRoot: jsRoot, conv: &DefaultCodingConvention{}, changed: true}
}
func (*Compiler) Implements_jscomp_AbstractCompiler()             {}
func (c *Compiler) IsJavaNil_() bool                              { return c == nil }
func (c *Compiler) GetCodingConvention() CodingConvention         { return c.conv }
func (c *Compiler) HasScopeChanged(_ pkg_rhino.Node) bool         { return c.changed }
func (c *Compiler) GetRoot() pkg_rhino.Node                       { return c.jsRoot }

type ControlFlowGraph[N comparable] interface {
	pkg_graph.Graph[N, int32]
}

type ControlFlowGraphImpl struct {
	pkg_graph.LinkedDirectedGraphBase[pkg_rhino.Node, int32]
}

func NewControlFlowGraphImpl(root pkg_rhino.Node) ControlFlowGraph[pkg_rhino.Node] {
	g := &ControlFlowGraphImpl{}
	pkg_graph.InitLinkedDirectedGraph[pkg_rhino.Node, int32](&g.LinkedDirectedGraphBase)
	g.CreateNode(root)
	for c := root.GetFirstChild(); c != 0; c = c.GetNext() {
		g.CreateNode(c)
		g.CreateEdge(1)
	}
	return g
}

type CfaBuilder struct {
	comp AbstractCompiler
	root pkg_rhino.Node
}

func ControlFlowAnalysis_S_Builder() *CfaBuilder                  { return &CfaBuilder{} }
func (b *CfaBuilder) SetCompiler(c AbstractCompiler) *CfaBuilder  { b.comp = c; return b }
func (b *CfaBuilder) SetCfgRoot(r pkg_rhino.Node) *CfaBuilder     { b.root = r; return b }
func (b *CfaBuilder) SetTraverseFunctions(bool) *CfaBuilder       { return b }
func (b *CfaBuilder) SetIncludeEdgeAnnotations(bool) *CfaBuilder  { return b }
func (b *CfaBuilder) ComputeCfg() ControlFlowGraph[pkg_rhino.Node] {
	return NewControlFlowGraphImpl(b.root)
}

func (this_ *Compiler) ComputeCFG() ControlFlowGraph[pkg_rhino.Node] {
	return ControlFlowAnalysis_S_Builder().SetCompiler(this_).SetCfgRoot(this_.jsRoot).ComputeCfg()
}

func DotFormatter_ToDot_2(_ pkg_rhino.Node, cfg ControlFlowGraph[pkg_rhino.Node]) string {
	if cfg == nil || cfg.IsJavaNil_() {
		return ""
	}
	return "digraph{}"
}

func DotFormatter_AppendDot(_ pkg_rhino.Node, cfg ControlFlowGraph[pkg_rhino.Node], out *[]string) {
	if cfg != nil && !cfg.IsJavaNil_() {
		*out = append(*out, "dot")
	}
}

func (this_ *Compiler) GetAstDotGraph() string {
	if this_.jsRoot != 0 {
		var cfg ControlFlowGraph[pkg_rhino.Node] = ControlFlowAnalysis_S_Builder().SetCompiler(this_).SetCfgRoot(this_.jsRoot).SetTraverseFunctions(true).ComputeCfg()
		_ = cfg
		return DotFormatter_ToDot_2(this_.jsRoot, cfg)
	}
	return ""
}

type AbstractCommandLineRunnerBase struct {
	compiler *Compiler
	printAst bool
	out      []string
}

func (this_ *AbstractCommandLineRunnerBase) ProcessResults() int32 {
	if this_.printAst {
		var cfg ControlFlowGraph[pkg_rhino.Node] = this_.compiler.ComputeCFG()
		_ = cfg
		DotFormatter_AppendDot(this_.compiler.GetRoot(), cfg, &this_.out)
		return 0
	}
	return 0
}

type AbstractVar interface {
	Implements_jscomp_AbstractVar()
	IsJavaNil_() bool
	GetName() string
	GetNameNode() pkg_rhino.Node
	IsGlobal() bool
	IsExtern() bool
	GetScope() any
}

type Var struct {
	name     string
	nameNode pkg_rhino.Node
	scope    *Scope
	extern   bool
}

func (*Var) Implements_jscomp_AbstractVar()     {}
func (v *Var) IsJavaNil_() bool                 { return v == nil }
func (v *Var) HashCode() int32                  { return int32(len(v.name)) }
func (v *Var) Equals(o any) bool                { ov, ok := o.(*Var); return ok && v == ov }
func (v *Var) GetName() string                  { return v.name }
func (v *Var) GetNameNode() pkg_rhino.Node      { return v.nameNode }
func (v *Var) IsGlobal() bool                   { return v.scope != nil && v.scope.IsGlobal() }
func (v *Var) IsExtern() bool                   { return v.extern }
func (v *Var) GetScope() any                    { return v.scope }

func NewVar(name string, nameNode pkg_rhino.Node, scope *Scope, extern bool) *Var {
	this_ := &Var{}
	InitVar(this_, name, nameNode, scope, extern)
	return this_
}

func InitVar(this_ *Var, name string, nameNode pkg_rhino.Node, scope *Scope, extern bool) {
	this_.name = name
	this_.nameNode = nameNode
	this_.scope = scope
	this_.extern = extern
}

type AbstractScope interface {
	Implements_jscomp_AbstractScope()
	IsJavaNil_() bool
	GetRootNode_0() pkg_rhino.Node
	GetParent() any
	IsGlobal() bool
	IsModuleScope() bool
	IsFunctionBlockScope() bool
	IsFunctionScope() bool
	GetVarCount() int32
	GetVar(string) any
	GetVarIterable() j2gort.Ext_java_lang_Iterable
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

func (*AbstractScopeBase) Implements_jscomp_AbstractScope() {}
func (s *AbstractScopeBase) IsJavaNil_() bool               { return s == nil }
func (s *AbstractScopeBase) BindSelf_(x any)                { s.self_ = x.(AbstractScope) }
func (s *AbstractScopeBase) GetRootNode_0() pkg_rhino.Node  { return s.rootNode }

type Scope struct {
	AbstractScopeBase
	parent *Scope
	var0   *Var
	var1   *Var
	var2   *Var
	var3   *Var
	var4   *Var
	var5   *Var
	vars   *j2gort.OrderedMap[string, *Var]
	count  int32
}

func (s *Scope) IsJavaNil_() bool              { return s == nil }
func (s *Scope) GetParent() any                { return s.parent }
func (s *Scope) IsGlobal() bool                { return s.parent == nil }
func (s *Scope) IsModuleScope() bool           { return false }
func (s *Scope) IsFunctionScope() bool         { return s.rootNode.IsFunction() }
func (s *Scope) IsFunctionBlockScope() bool    { return s.rootNode.IsBlock() && s.parent != nil && s.parent.IsFunctionScope() }
func (s *Scope) GetVarCount() int32            { return s.count }

func NewScope_1(root pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	this_.rootNode = root
	return this_
}

func NewScope_2(parent *Scope, root pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_2(this_, parent, root)
	return this_
}

func InitScope_2(this_ *Scope, parent *Scope, root pkg_rhino.Node) {
	this_.rootNode = root
	this_.parent = parent
}

func Scope_CreateGlobalScope(root pkg_rhino.Node) *Scope {
	return NewScope_1(root)
}

func (this_ *Scope) Declare(name string, nameNode pkg_rhino.Node, extern bool) *Var {
	v := NewVar(name, nameNode, this_, extern)
	this_.count++
	if this_.vars != nil {
		this_.vars.Put(name, v)
		return v
	}
	switch {
	case this_.var0 == nil:
		this_.var0 = v
	case this_.var1 == nil:
		this_.var1 = v
	case this_.var2 == nil:
		this_.var2 = v
	case this_.var3 == nil:
		this_.var3 = v
	case this_.var4 == nil:
		this_.var4 = v
	case this_.var5 == nil:
		this_.var5 = v
	default:
		this_.vars = j2gort.NewOrderedMapPtrCap[string, *Var](16)
		this_.vars.Put(this_.var0.name, this_.var0)
		this_.vars.Put(this_.var1.name, this_.var1)
		this_.vars.Put(this_.var2.name, this_.var2)
		this_.vars.Put(this_.var3.name, this_.var3)
		this_.vars.Put(this_.var4.name, this_.var4)
		this_.vars.Put(this_.var5.name, this_.var5)
		this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5 = nil, nil, nil, nil, nil, nil
		this_.vars.Put(name, v)
	}
	return v
}

func (this_ *Scope) GetOwnSlot(name string) *Var {
	if this_.vars != nil {
		return this_.vars.Get(name)
	}
	for _, v := range []*Var{this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5} {
		if v != nil && v.name == name {
			return v
		}
	}
	return nil
}

func (this_ *Scope) GetVar(name string) any {
	for cur := this_; cur != nil; cur = cur.parent {
		if v := cur.GetOwnSlot(name); v != nil {
			return v
		}
	}
	return nil
}

func (this_ *Scope) GetVarIterable() j2gort.Ext_java_lang_Iterable {
	var items []any
	if this_.vars != nil {
		cur := this_.vars.IterValues()
		for cur.HasNext() {
			items = append(items, cur.Next())
		}
		return j2gort.SliceIterable{Items: items}
	}
	for _, v := range []*Var{this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5} {
		if v != nil {
			items = append(items, v)
		}
	}
	return j2gort.SliceIterable{Items: items}
}

type ScopeCreator interface {
	Implements_jscomp_ScopeCreator()
	IsJavaNil_() bool
	CreateScope(pkg_rhino.Node, AbstractScope) AbstractScope
}

func ToIface_ScopeCreator(sc ScopeCreator) ScopeCreator { return sc }

type SyntacticScopeCreator_RedeclarationHandler interface{}
type defaultRedeclHandler struct{}

func StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER() SyntacticScopeCreator_RedeclarationHandler {
	return defaultRedeclHandler{}
}

type SyntacticScopeCreator_ScopeScanner struct {
	compiler                      AbstractCompiler
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler
	treatProvidesAsRedeclarations bool
	scope                         *Scope
}

func (s *SyntacticScopeCreator_ScopeScanner) Reset() { s.scope = nil }

func (s *SyntacticScopeCreator_ScopeScanner) Populate(n pkg_rhino.Node, extern bool) {
	for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
		tok := c.GetToken()
		if tok == pkg_rhino.TokVAR || tok == pkg_rhino.TokLET || tok == pkg_rhino.TokCONST {
			nm := c.GetFirstChild()
			if nm != 0 && nm.IsName() {
				s.scope.Declare(nm.GetString(), nm, extern)
			}
		} else if !c.IsFunction() && !c.IsBlock() {
			s.Populate(c, extern)
		}
	}
}

type SyntacticScopeCreator struct {
	compiler                      AbstractCompiler
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler
	treatProvidesAsRedeclarations bool
	reusableScanner               *SyntacticScopeCreator_ScopeScanner
	scannerInUse                  bool
}

func clinit_SyntacticScopeCreator()                              {}
func (*SyntacticScopeCreator) Implements_jscomp_ScopeCreator()   {}
func (s *SyntacticScopeCreator) IsJavaNil_() bool                { return s == nil }

func NewSyntacticScopeCreator_1(compiler AbstractCompiler) *SyntacticScopeCreator {
	this_ := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_1(this_, compiler)
	return this_
}

func InitSyntacticScopeCreator_1(this_ *SyntacticScopeCreator, compiler AbstractCompiler) {
	this_.compiler = compiler
	this_.redeclarationHandler = StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER()
	this_.reusableScanner = &SyntacticScopeCreator_ScopeScanner{compiler: compiler}
}

func (this_ *SyntacticScopeCreator) CreateScope(n pkg_rhino.Node, parent AbstractScope) AbstractScope {
	var s *Scope
	if parent == nil || parent.IsJavaNil_() {
		s = Scope_CreateGlobalScope(n)
	} else {
		s = NewScope_2(parent.(*Scope), n)
	}
	sc := &SyntacticScopeCreator_ScopeScanner{scope: s}
	sc.Populate(n, false)
	return s
}

type MemoizedScopeCreator struct {
	scopesByScopeRoot j2gort.OrderedMap[pkg_rhino.Node, AbstractScope]
	delegate          ScopeCreator
}

func (*MemoizedScopeCreator) Implements_jscomp_ScopeCreator() {}
func (m *MemoizedScopeCreator) IsJavaNil_() bool              { return m == nil }

func (this_ *MemoizedScopeCreator) initFields_MemoizedScopeCreator() {
	this_.scopesByScopeRoot = j2gort.NewOrderedMap[pkg_rhino.Node, AbstractScope]()
}

func InitMemoizedScopeCreator(this_ *MemoizedScopeCreator, delegate ScopeCreator) {
	this_.initFields_MemoizedScopeCreator()
	this_.delegate = delegate
}

func (this_ *MemoizedScopeCreator) CreateScope(n pkg_rhino.Node, parent AbstractScope) AbstractScope {
	if existing := this_.scopesByScopeRoot.Get(n); existing != nil {
		return existing
	}
	s := this_.delegate.CreateScope(n, parent)
	this_.scopesByScopeRoot.Put(n, s)
	return s
}

func NodeUtil_IsNameDeclaration(n pkg_rhino.Node) bool {
	if n == 0 {
		return false
	}
	t := n.GetToken()
	return t == pkg_rhino.TokVAR || t == pkg_rhino.TokLET || t == pkg_rhino.TokCONST
}

func NodeUtil_IsAnyFor(n pkg_rhino.Node) bool {
	return n != 0 && n.GetToken() == pkg_rhino.TokFOR
}

func NodeUtil_GetRValueOfLValue(n pkg_rhino.Node) pkg_rhino.Node {
	if n == 0 {
		return 0
	}
	p := n.GetParent()
	if NodeUtil_IsNameDeclaration(p) {
		return n.GetFirstChild()
	}
	return 0
}

func NodeUtil_IsValidCfgRoot(n pkg_rhino.Node) bool {
	return n.IsScript() || n.IsFunction()
}

type NodeUtil_AllVarsDeclaredInFunction struct {
	vars map[string]*Var
}

func (a *NodeUtil_AllVarsDeclaredInFunction) GetAllVariables() map[string]*Var { return a.vars }

func NodeUtil_GetAllVarsDeclaredInFunction(compiler AbstractCompiler, creator ScopeCreator, fnScope *Scope) *NodeUtil_AllVarsDeclaredInFunction {
	out := make(map[string]*Var)
	fnNode := fnScope.GetRootNode_0()
	body := fnNode.GetLastChild()
	if body != 0 {
		bs := creator.CreateScope(body, fnScope).(*Scope)
		it := j2gort.NewIterCursor(bs.GetVarIterable())
		for it.HasNext() {
			v := it.Next().(*Var)
			out[v.name] = v
		}
	}
	return &NodeUtil_AllVarsDeclaredInFunction{vars: out}
}

func ChangeTracker_IsChangeScopeRoot(n pkg_rhino.Node) bool {
	return n.IsScript() || n.IsFunction()
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "core.go"), []byte(jscompCoreSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	nodeTraversalSrc := `package jscomp

import (
	"example.com/w23ctest/gen/j2gort"
	pkg_rhino "example.com/w23ctest/gen/rhino"
)

func clinit_NodeTraversal() {}

type NodeTraversal_Callback interface {
	Implements_jscomp_NodeTraversal_Callback()
	IsJavaNil_() bool
	ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool
	Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type NodeTraversal_ScopedCallback interface {
	NodeTraversal_Callback
	EnterScope(*NodeTraversal)
	ExitScope(*NodeTraversal)
}

type NodeTraversal_AbstractPostOrderCallbackInterface interface {
	Implements_jscomp_NodeTraversal_AbstractPostOrderCallbackInterface()
	IsJavaNil_() bool
	Visit(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type lambda_NodeTraversal_AbstractPostOrderCallbackInterface struct {
	fn func(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

func (l *lambda_NodeTraversal_AbstractPostOrderCallbackInterface) IsJavaNil_() bool { return l == nil }
func (*lambda_NodeTraversal_AbstractPostOrderCallbackInterface) Implements_jscomp_NodeTraversal_AbstractPostOrderCallbackInterface() {
}
func (l *lambda_NodeTraversal_AbstractPostOrderCallbackInterface) Visit(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	l.fn(t, n, parent)
}
func NewLambda_NodeTraversal_AbstractPostOrderCallbackInterface(fn func(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)) NodeTraversal_AbstractPostOrderCallbackInterface {
	if fn == nil {
		return nil
	}
	return &lambda_NodeTraversal_AbstractPostOrderCallbackInterface{fn: fn}
}

type NodeTraversal_AbstractPostOrderCallbackBase struct{}

func (*NodeTraversal_AbstractPostOrderCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (b *NodeTraversal_AbstractPostOrderCallbackBase) IsJavaNil_() bool                        { return b == nil }
func (b *NodeTraversal_AbstractPostOrderCallbackBase) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}
func InitNodeTraversal_AbstractPostOrderCallback(_ *NodeTraversal_AbstractPostOrderCallbackBase) {}

type NodeTraversal_AbstractScopedCallbackBase struct{}

func (*NodeTraversal_AbstractScopedCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (b *NodeTraversal_AbstractScopedCallbackBase) IsJavaNil_() bool                        { return b == nil }
func (b *NodeTraversal_AbstractScopedCallbackBase) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}
func (b *NodeTraversal_AbstractScopedCallbackBase) Visit_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {
}

type NodeTraversal_AbstractCfgCallback interface {
	NodeTraversal_ScopedCallback
	EnterScopeWithCfg(*NodeTraversal)
	ExitScopeWithCfg(*NodeTraversal)
}

type NodeTraversal_AbstractCfgCallbackBase struct {
	self_ NodeTraversal_AbstractCfgCallback
	cfgs  j2gort.Deque[ControlFlowGraph[pkg_rhino.Node]]
}

func (*NodeTraversal_AbstractCfgCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (b *NodeTraversal_AbstractCfgCallbackBase) IsJavaNil_() bool                        { return b == nil }
func (b *NodeTraversal_AbstractCfgCallbackBase) BindSelf_(s any)                         { b.self_ = s.(NodeTraversal_AbstractCfgCallback) }
func (b *NodeTraversal_AbstractCfgCallbackBase) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}
func (b *NodeTraversal_AbstractCfgCallbackBase) Visit_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {
}
func (b *NodeTraversal_AbstractCfgCallbackBase) GetControlFlowGraph(_ AbstractCompiler) ControlFlowGraph[pkg_rhino.Node] {
	return b.cfgs.Peek()
}
func (b *NodeTraversal_AbstractCfgCallbackBase) EnterScope(t *NodeTraversal) {
	if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) {
		b.cfgs.Push(ControlFlowAnalysis_S_Builder().SetCompiler(t.GetCompiler()).SetCfgRoot(t.GetScopeRoot()).ComputeCfg())
	}
	b.self_.EnterScopeWithCfg(t)
}
func (b *NodeTraversal_AbstractCfgCallbackBase) ExitScope(t *NodeTraversal) {
	b.self_.ExitScopeWithCfg(t)
	if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) {
		b.cfgs.Pop()
	}
}

type NodeTraversal_Builder struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

func NodeTraversal_S_Builder() *NodeTraversal_Builder { return &NodeTraversal_Builder{} }
func (this_ *NodeTraversal_Builder) SetCompiler(c AbstractCompiler) *NodeTraversal_Builder {
	this_.compiler = c
	return this_
}
func (this_ *NodeTraversal_Builder) SetCallback_1_NodeTraversal_Callback(cb NodeTraversal_Callback) *NodeTraversal_Builder {
	this_.callback = cb
	return this_
}
func (this_ *NodeTraversal_Builder) SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(x NodeTraversal_AbstractPostOrderCallbackInterface) *NodeTraversal_Builder {
	return this_.SetCallback_1_NodeTraversal_Callback(NewNodeTraversal__Anon_66(this_, x))
}
func (this_ *NodeTraversal_Builder) SetScopeCreator(sc ScopeCreator) *NodeTraversal_Builder {
	this_.scopeCreator = sc
	return this_
}
func (this_ *NodeTraversal_Builder) Build() *NodeTraversal {
	t := &NodeTraversal{}
	InitNodeTraversal(t, this_)
	return t
}
func (this_ *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverse_1(root)
}
func (this_ *NodeTraversal_Builder) TraverseRoots(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverseRoots_2(externs, root)
}

type NodeTraversal__Anon_66 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_  *NodeTraversal_Builder
	cap_0_x NodeTraversal_AbstractPostOrderCallbackInterface
}

func NewNodeTraversal__Anon_66(outer_ *NodeTraversal_Builder, cap_p_x NodeTraversal_AbstractPostOrderCallbackInterface) *NodeTraversal__Anon_66 {
	this_ := &NodeTraversal__Anon_66{outer_: outer_, cap_0_x: cap_p_x}
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	return this_
}

func (this_ *NodeTraversal__Anon_66) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.cap_0_x.Visit(t, n, parent)
}

type NodeTraversal__Local_67_TraverseScopeRootsCallback struct {
	scopeNode            pkg_rhino.Node
	insideScopeNode      bool
	cap_0_traverseNested bool
	cap_1_cb             NodeTraversal_Callback
}

func (*NodeTraversal__Local_67_TraverseScopeRootsCallback) Implements_jscomp_NodeTraversal_Callback() {
}
func (c *NodeTraversal__Local_67_TraverseScopeRootsCallback) IsJavaNil_() bool { return c == nil }
func NewNodeTraversal__Local_67_TraverseScopeRootsCallback(traverseNested bool, cb NodeTraversal_Callback) *NodeTraversal__Local_67_TraverseScopeRootsCallback {
	return &NodeTraversal__Local_67_TraverseScopeRootsCallback{cap_0_traverseNested: traverseNested, cap_1_cb: cb}
}
func (this_ *NodeTraversal__Local_67_TraverseScopeRootsCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	if this_.scopeNode == n {
		this_.insideScopeNode = true
	}
	return (this_.cap_0_traverseNested || this_.scopeNode == n || !ChangeTracker_IsChangeScopeRoot(n)) && this_.cap_1_cb.ShouldTraverse_3(t, n, parent)
}
func (this_ *NodeTraversal__Local_67_TraverseScopeRootsCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	if this_.scopeNode == n {
		this_.insideScopeNode = false
	}
	this_.cap_1_cb.Visit_3(t, n, parent)
}
func (this_ *NodeTraversal__Local_67_TraverseScopeRootsCallback) EnterScope(t *NodeTraversal) {
	if this_.insideScopeNode {
		if sc, ok := j2gort.TryCast[NodeTraversal_ScopedCallback](this_.cap_1_cb); ok {
			sc.EnterScope(t)
		}
	}
}
func (this_ *NodeTraversal__Local_67_TraverseScopeRootsCallback) ExitScope(t *NodeTraversal) {
	if this_.insideScopeNode {
		if sc, ok := j2gort.TryCast[NodeTraversal_ScopedCallback](this_.cap_1_cb); ok {
			sc.ExitScope(t)
		}
	}
}

type NodeTraversal struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
	scopeNodes                                     []pkg_rhino.Node
	scopeObjs                                      []AbstractScope
	scopeDepth                                     int32
}

func InitNodeTraversal(this_ *NodeTraversal, b *NodeTraversal_Builder) {
	this_.compiler = b.compiler
	this_.callback = b.callback
	this_.obeyDestructuringAndDefaultValueExecutionOrder = b.obeyDestructuringAndDefaultValueExecutionOrder
	if b.scopeCreator == nil || b.scopeCreator.IsJavaNil_() {
		this_.scopeCreator = NewSyntacticScopeCreator_1(this_.compiler)
	} else {
		this_.scopeCreator = b.scopeCreator
	}
}

func (this_ *NodeTraversal) GetCompiler() AbstractCompiler  { return this_.compiler }
func (this_ *NodeTraversal) GetScopeCreator() ScopeCreator  { return this_.scopeCreator }
func (this_ *NodeTraversal) InGlobalScope() bool            { return this_.scopeDepth <= 1 }
func (this_ *NodeTraversal) InFunctionBlockScope() bool     { return this_.GetScope().IsFunctionBlockScope() }
func (this_ *NodeTraversal) GetScopeRoot() pkg_rhino.Node   { return this_.scopeNodes[this_.scopeDepth-1] }

func (this_ *NodeTraversal) GetScope() *Scope {
	idx := this_.scopeDepth - 1
	s := this_.scopeObjs[idx]
	if s == nil {
		var parent AbstractScope
		if idx > 0 {
			parent = this_.getScopeAt(idx - 1)
		}
		s = this_.scopeCreator.CreateScope(this_.scopeNodes[idx], parent)
		this_.scopeObjs[idx] = s
	}
	return s.(*Scope)
}

func (this_ *NodeTraversal) getScopeAt(idx int32) AbstractScope {
	s := this_.scopeObjs[idx]
	if s == nil {
		var parent AbstractScope
		if idx > 0 {
			parent = this_.getScopeAt(idx - 1)
		}
		s = this_.scopeCreator.CreateScope(this_.scopeNodes[idx], parent)
		this_.scopeObjs[idx] = s
	}
	return s
}

func (this_ *NodeTraversal) priv_NodeTraversal_ensureScopeCapacity_0() {
	if this_.scopeDepth == int32(len(this_.scopeNodes)) {
		n := int(this_.scopeDepth)*2 + 4
		ns := make([]pkg_rhino.Node, n)
		copy(ns, this_.scopeNodes)
		this_.scopeNodes = ns
		os := make([]AbstractScope, n)
		copy(os, this_.scopeObjs)
		this_.scopeObjs = os
	}
}

func (this_ *NodeTraversal) pushScopeNode(n pkg_rhino.Node) {
	this_.priv_NodeTraversal_ensureScopeCapacity_0()
	this_.scopeNodes[this_.scopeDepth] = n
	this_.scopeObjs[this_.scopeDepth] = nil
	this_.scopeDepth++
	if sc, ok := this_.callback.(NodeTraversal_ScopedCallback); ok {
		sc.EnterScope(this_)
	}
}

func (this_ *NodeTraversal) popScopeNode() {
	if sc, ok := this_.callback.(NodeTraversal_ScopedCallback); ok {
		sc.ExitScope(this_)
	}
	this_.scopeDepth--
	this_.scopeNodes[this_.scopeDepth] = 0
	this_.scopeObjs[this_.scopeDepth] = nil
}

func (this_ *NodeTraversal) walk(n pkg_rhino.Node, parent pkg_rhino.Node) {
	if n == 0 {
		return
	}
	createsScope := n.IsFunction() || n.IsBlock()
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	if createsScope {
		this_.pushScopeNode(n)
	}
	for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
		this_.walk(c, n)
	}
	if createsScope {
		this_.popScopeNode()
	}
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverse_1(root pkg_rhino.Node) {
	if !this_.callback.ShouldTraverse_3(this_, root, 0) {
		return
	}
	this_.pushScopeNode(root)
	for c := root.GetFirstChild(); c != 0; c = c.GetNext() {
		this_.walk(c, root)
	}
	this_.popScopeNode()
	this_.callback.Visit_3(this_, root, 0)
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseRoots_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.pushScopeNode(root)
	for c := externs.GetFirstChild(); c != 0; c = c.GetNext() {
		this_.walk(c, externs)
	}
	for c := root.GetFirstChild(); c != 0; c = c.GetNext() {
		this_.walk(c, root)
	}
	this_.popScopeNode()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseScopeRoot_3(scopeRoot pkg_rhino.Node) {
	var chain []pkg_rhino.Node
	for p := scopeRoot.GetParent(); p != 0; p = p.GetParent() {
		if p.IsScript() || p.IsFunction() || p.IsBlock() {
			chain = append(chain, p)
		}
	}
	for i := len(chain) - 1; i >= 0; i-- {
		this_.priv_NodeTraversal_ensureScopeCapacity_0()
		this_.scopeNodes[this_.scopeDepth] = chain[i]
		this_.scopeObjs[this_.scopeDepth] = nil
		this_.scopeDepth++
	}
	this_.walk(scopeRoot, scopeRoot.GetParent())
	for this_.scopeDepth > 0 {
		this_.scopeDepth--
		this_.scopeNodes[this_.scopeDepth] = 0
		this_.scopeObjs[this_.scopeDepth] = nil
	}
}

func NodeTraversal_Traverse(compiler AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).Traverse(root)
}

func NodeTraversal_TraverseRoots(compiler AbstractCompiler, cb NodeTraversal_Callback, externs pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).TraverseRoots(externs, root)
}

func NodeTraversal_TraverseScopeRoots(compiler AbstractCompiler, scopeNodes j2gort.Ext_java_util_List, cb NodeTraversal_Callback, traverseNested bool) {
	clinit_NodeTraversal()
	var scb *NodeTraversal__Local_67_TraverseScopeRootsCallback = NewNodeTraversal__Local_67_TraverseScopeRootsCallback(traverseNested, cb)
	_ = scb
	var val_8 MemoizedScopeCreator
	InitMemoizedScopeCreator(&val_8, NewSyntacticScopeCreator_1(compiler))
	var scopeCreator *MemoizedScopeCreator = &val_8
	_ = scopeCreator
	cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(scopeNodes))
	for cur_10.HasNext() {
		scopeNode := pkg_rhino.Node(cur_10.NextNode())
		scb.scopeNode = scopeNode
		NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(scb).SetScopeCreator(scopeCreator).Build().priv_NodeTraversal_traverseScopeRoot_3(scopeNode)
	}
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "node_traversal.go"), []byte(nodeTraversalSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	passesSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w23ctest/gen/j2gort"
	pkg_graph "example.com/w23ctest/gen/graph"
	pkg_rhino "example.com/w23ctest/gen/rhino"
)

// --- 1. MakeDeclaredNamesUnique & PostOrder CallbackUnwrap ---

type MakeDeclaredNamesUnique_Renamer interface {
	Implements_jscomp_MakeDeclaredNamesUnique_Renamer()
	AddDeclaredName(string, bool)
	CreateForChildScope() MakeDeclaredNamesUnique_Renamer
}

type MakeDeclaredNamesUnique_ContextualRenamer struct {
	names map[string]bool
}

func NewContextualRenamer() *MakeDeclaredNamesUnique_ContextualRenamer {
	return &MakeDeclaredNamesUnique_ContextualRenamer{names: make(map[string]bool)}
}
func (*MakeDeclaredNamesUnique_ContextualRenamer) Implements_jscomp_MakeDeclaredNamesUnique_Renamer() {
}
func (r *MakeDeclaredNamesUnique_ContextualRenamer) AddDeclaredName(name string, _ bool) {
	r.names[name] = true
}
func (r *MakeDeclaredNamesUnique_ContextualRenamer) CreateForChildScope() MakeDeclaredNamesUnique_Renamer {
	return NewContextualRenamer()
}

type MakeDeclaredNamesUnique struct {
	NodeTraversal_AbstractScopedCallbackBase
	renamerStack j2gort.Deque[MakeDeclaredNamesUnique_Renamer]
	rootRenamer  MakeDeclaredNamesUnique_Renamer
}

func NewMakeDeclaredNamesUnique() *MakeDeclaredNamesUnique {
	return &MakeDeclaredNamesUnique{rootRenamer: NewContextualRenamer()}
}

func (this_ *MakeDeclaredNamesUnique) EnterScope(t *NodeTraversal) {
	var renamer MakeDeclaredNamesUnique_Renamer
	if this_.renamerStack.IsEmpty() {
		renamer = this_.rootRenamer
	} else {
		renamer = this_.renamerStack.Peek().CreateForChildScope()
	}
	this_.priv_MakeDeclaredNamesUnique_findDeclaredNames_1(t, renamer)
	this_.renamerStack.Push(renamer)
}

func (this_ *MakeDeclaredNamesUnique) ExitScope(t *NodeTraversal) {
	_ = t
	if !this_.renamerStack.IsEmpty() {
		this_.renamerStack.Pop()
	}
}

func (this_ *MakeDeclaredNamesUnique) priv_MakeDeclaredNamesUnique_findDeclaredNames_1(t *NodeTraversal, renamer MakeDeclaredNamesUnique_Renamer) {
	cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(t.GetScope().GetVarIterable()))
	for cur_4.HasNext() {
		rawElem_3 := cur_4.Next()
		v := j2gort.CastPtr[Var](rawElem_3)
		_ = v
		renamer.AddDeclaredName(v.GetName(), false)
	}
}

// Negative soundness control: ContextualRenameInverter retains *Var in DeclaredVariableInfo!
type MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo interface {
	Implements_jscomp_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo()
}

type MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo struct {
	var_ *Var
}

func (*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) Implements_jscomp_MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo() {
}

type MakeDeclaredNamesUnique_ContextualRenameInverter struct {
	NodeTraversal_AbstractScopedCallbackBase
	saved []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) EnterScope(t *NodeTraversal) {
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(t.GetScope().GetVarIterable()))
	for cur.HasNext() {
		v := j2gort.CastPtr[Var](cur.Next())
		this_.saved = append(this_.saved, &MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo{var_: v})
	}
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) ExitScope(_ *NodeTraversal) {}

// PostOrder interface concrete implementations + lambda caller:
type EphemeralPostOrderCb struct {
	emptyFnCount int
}

func (*EphemeralPostOrderCb) Implements_jscomp_NodeTraversal_AbstractPostOrderCallbackInterface() {}
func (c *EphemeralPostOrderCb) IsJavaNil_() bool                                                  { return c == nil }
func (c *EphemeralPostOrderCb) Visit(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsBlock() && t.GetScope().GetVarCount() == 0 {
		c.emptyFnCount++
	}
}

type EscapingPostOrderCb struct {
	savedScope *Scope
}

func (*EscapingPostOrderCb) Implements_jscomp_NodeTraversal_AbstractPostOrderCallbackInterface() {}
func (c *EscapingPostOrderCb) IsJavaNil_() bool                                                  { return c == nil }
func (c *EscapingPostOrderCb) Visit(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsBlock() {
		c.savedScope = t.GetScope()
	}
}

type J2clConstantHoisterPass struct {
	compiler AbstractCompiler
	hoisted  int
}

func J2clConstantHoisterPass_IsHoistable(t *NodeTraversal, n pkg_rhino.Node) bool {
	return n.IsBlock() && t.GetScope().GetVarCount() >= 0
}

func (this_ *J2clConstantHoisterPass) Process_2(_ pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(NewLambda_NodeTraversal_AbstractPostOrderCallbackInterface(func(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
		_ = parent
		if J2clConstantHoisterPass_IsHoistable(t, n) {
			this_.hoisted++
		}
	})).Traverse(root)
}

// --- 2. FlowSensitiveInlineVariables & CoalesceVariableNames ---

type LiveVariablesAnalysis struct {
	vars map[string]*Var
}

func NewLiveVariablesAnalysis(cfg ControlFlowGraph[pkg_rhino.Node], fnScope *Scope, blkScope *Scope, comp AbstractCompiler, sc ScopeCreator, allVars *NodeUtil_AllVarsDeclaredInFunction) *LiveVariablesAnalysis {
	_, _, _, _ = cfg, fnScope, blkScope, comp
	_ = sc
	return &LiveVariablesAnalysis{vars: allVars.GetAllVariables()}
}
func (l *LiveVariablesAnalysis) Analyze()                        {}
func (l *LiveVariablesAnalysis) GetAllVariables() map[string]*Var { return l.vars }

type MustBeReachingVariableDef struct{}
type MaybeReachingVariableUse struct{}
type FlowSensitiveInlineVariables_Candidate struct{ name string }

type FlowSensitiveInlineVariables struct {
	compiler     AbstractCompiler
	cfg          ControlFlowGraph[pkg_rhino.Node]
	candidates   []*FlowSensitiveInlineVariables_Candidate
	reachingDef  *MustBeReachingVariableDef
	reachingUses *MaybeReachingVariableUse
	inlinedFns   int
}

func (*FlowSensitiveInlineVariables) Implements_jscomp_NodeTraversal_Callback() {}
func (f *FlowSensitiveInlineVariables) IsJavaNil_() bool                        { return f == nil }
func (f *FlowSensitiveInlineVariables) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}
func (this_ *FlowSensitiveInlineVariables) EnterScope(t *NodeTraversal) {
	if t.InGlobalScope() || !t.GetScope().IsFunctionBlockScope() {
		return
	}
	fnRoot := t.GetScopeRoot().GetParent()
	this_.cfg = ControlFlowAnalysis_S_Builder().SetCompiler(this_.compiler).SetCfgRoot(fnRoot).SetIncludeEdgeAnnotations(true).ComputeCfg()
	this_.reachingDef = &MustBeReachingVariableDef{}
	this_.candidates = []*FlowSensitiveInlineVariables_Candidate{{name: "x"}}
	this_.reachingUses = &MaybeReachingVariableUse{}
	if this_.cfg != nil && len(this_.candidates) > 0 {
		this_.inlinedFns++
	}
}
func (this_ *FlowSensitiveInlineVariables) ExitScope(t *NodeTraversal) {
	_ = t
}
func (this_ *FlowSensitiveInlineVariables) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}
func (this_ *FlowSensitiveInlineVariables) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(this_).TraverseRoots(externs, root)
}

type CoalesceVariableNames struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler   AbstractCompiler
	colorings  j2gort.Deque[pkg_graph.GraphColoring[*Var, int32]]
}

func NewCoalesceVariableNames(comp AbstractCompiler) *CoalesceVariableNames {
	c := &CoalesceVariableNames{compiler: comp}
	c.BindSelf_(c)
	return c
}
func (this_ *CoalesceVariableNames) EnterScopeWithCfg(t *NodeTraversal) {
	if !t.GetScope().IsFunctionScope() {
		return
	}
	ug := pkg_graph.NewLinkedUndirectedGraph[*Var, int32]()
	v1 := &Var{name: "a"}
	v2 := &Var{name: "b"}
	ug.Connect(v1, 1, v2)
	this_.colorings.Push(pkg_graph.NewGreedyGraphColoring[*Var, int32](ug))
}
func (this_ *CoalesceVariableNames) ExitScopeWithCfg(t *NodeTraversal) {
	if !t.GetScope().IsFunctionScope() {
		return
	}
	this_.colorings.Pop()
}

// --- 3. DeadAssignmentsElimination ---

type DeadAssignmentsElimination_BailoutInformation struct {
	containsFunction        bool
	containsRemovableAssign bool
}

type DeadAssignmentsElimination struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler      AbstractCompiler
	liveness      *LiveVariablesAnalysis
	functionStack j2gort.Deque[*DeadAssignmentsElimination_BailoutInformation]
	eliminated    int
}

func NewDeadAssignmentsElimination(comp AbstractCompiler) *DeadAssignmentsElimination {
	d := &DeadAssignmentsElimination{compiler: comp}
	d.BindSelf_(d)
	return d
}

func (this_ *DeadAssignmentsElimination) EnterScopeWithCfg(t *NodeTraversal) {
	if t.InFunctionBlockScope() {
		this_.functionStack.AddFirst(&DeadAssignmentsElimination_BailoutInformation{containsRemovableAssign: true})
	}
}

func (this_ *DeadAssignmentsElimination) ExitScopeWithCfg(t *NodeTraversal) {
	if t.InFunctionBlockScope() {
		this_.priv_DeadAssignmentsElimination_eliminateDeadAssignments_2319(t)
		this_.functionStack.RemoveFirst()
	}
}

func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_eliminateDeadAssignments_2319(t *NodeTraversal) {
	if !this_.compiler.HasScopeChanged(t.GetScopeRoot().GetParent()) {
		return
	}
	cur := this_.functionStack.PeekFirst()
	if cur.containsFunction || !cur.containsRemovableAssign {
		return
	}
	blockScope := t.GetScope()
	functionScope := j2gort.CastPtr[Scope](blockScope.GetParent())
	cfg := this_.GetControlFlowGraph(this_.compiler)
	clinit_SyntacticScopeCreator()
	var val_3 SyntacticScopeCreator
	InitSyntacticScopeCreator_1(&val_3, this_.compiler)
	var scopeCreator *SyntacticScopeCreator = &val_3
	_ = scopeCreator
	var allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction = NodeUtil_GetAllVarsDeclaredInFunction(this_.compiler, scopeCreator, functionScope)
	_ = allVarsDeclaredInFunction
	this_.liveness = NewLiveVariablesAnalysis(cfg, functionScope, blockScope, this_.compiler, scopeCreator, allVarsDeclaredInFunction)
	this_.liveness.Analyze()
	if len(this_.liveness.GetAllVariables()) > 0 {
		this_.eliminated++
	}
}

func (this_ *DeadAssignmentsElimination) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// --- 4. ReferenceCollector, InlineObjectLiterals, GuardedCallback, PolyfillUsageFinder ---

type BasicBlock struct{ root pkg_rhino.Node }

func NewBasicBlock(root pkg_rhino.Node) *BasicBlock {
	b := &BasicBlock{}
	b.root = root
	return b
}

type Reference struct {
	nameNode   pkg_rhino.Node
	basicBlock *BasicBlock
	scope      *Scope
}

func InitReference(r *Reference, n pkg_rhino.Node, bb *BasicBlock, s *Scope) {
	r.nameNode = n
	r.basicBlock = bb
	r.scope = s
}
func (r *Reference) GetNode() pkg_rhino.Node { return r.nameNode }

type ReferenceCollection struct {
	first *Reference
	rest  []*Reference
}

func InitReferenceCollection(_ *ReferenceCollection) {}
func (rc *ReferenceCollection) Add(r *Reference) {
	if rc.first == nil {
		rc.first = r
	} else {
		rc.rest = append(rc.rest, r)
	}
}
func (rc *ReferenceCollection) GetInitializingReference() *Reference { return rc.first }

type ReferenceMap interface {
	GetReferences(any) any
}

type ReferenceCollector_Behavior interface {
	AfterExitScope(*NodeTraversal, ReferenceMap)
}

func clinit_ReferenceCollector() {}

type ReferenceCollector struct {
	compiler     AbstractCompiler
	behavior     ReferenceCollector_Behavior
	scopeCreator ScopeCreator
	varFilter    j2gort.Ext_com_google_common_base_Predicate
	referenceMap j2gort.CustomOrderedMap[*Var, *ReferenceCollection]
}

func (this_ *ReferenceCollector) initFields_ReferenceCollector() {
	this_.referenceMap = j2gort.NewCustomOrderedMap[*Var, *ReferenceCollection]()
}

func InitReferenceCollector_3(this_ *ReferenceCollector, compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator) {
	InitReferenceCollector_4(this_, compiler, behavior, creator, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue())
}

func InitReferenceCollector_4(this_ *ReferenceCollector, compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator, varFilter j2gort.Ext_com_google_common_base_Predicate) {
	this_.initFields_ReferenceCollector()
	this_.compiler = compiler
	this_.behavior = behavior
	this_.scopeCreator = creator
	this_.varFilter = varFilter
}

func (this_ *ReferenceCollector) GetReferences(v any) any {
	vv := j2gort.CastPtr[Var](v)
	if vv == nil {
		return nil
	}
	return j2gort.PtrToAny(this_.referenceMap.Get(vv))
}

type rcCallback struct {
	NodeTraversal_AbstractScopedCallbackBase
	rc *ReferenceCollector
}

func (c *rcCallback) EnterScope(_ *NodeTraversal) {}
func (c *rcCallback) ExitScope(t *NodeTraversal) {
	c.rc.behavior.AfterExitScope(t, c.rc)
}
func (c *rcCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if !n.IsName() {
		return
	}
	v := j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString()))
	if v == nil || !c.rc.varFilter.Apply(v) {
		return
	}
	var ref Reference
	InitReference(&ref, n, NewBasicBlock(t.GetScopeRoot()), t.GetScope())
	col := c.rc.referenceMap.Get(v)
	if col == nil {
		col = &ReferenceCollection{}
		InitReferenceCollection(col)
		c.rc.referenceMap.Put(v, col)
	}
	col.Add(&ref)
}

func (this_ *ReferenceCollector) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetScopeCreator(this_.scopeCreator).SetCallback_1_NodeTraversal_Callback(&rcCallback{rc: this_}).TraverseRoots(externs, root)
}

type NodeUtil_Visitor interface {
	Visit(pkg_rhino.Node)
}

func NodeUtil_VisitPreOrder(n pkg_rhino.Node, v NodeUtil_Visitor) {
	if n == 0 {
		return
	}
	v.Visit(n)
	for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
		NodeUtil_VisitPreOrder(c, v)
	}
}

type InlineObjectLiterals struct {
	compiler     AbstractCompiler
	SplitNames   []string
}

func NewInlineObjectLiterals(comp AbstractCompiler) *InlineObjectLiterals {
	return &InlineObjectLiterals{compiler: comp}
}

func (this_ *InlineObjectLiterals) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	clinit_ReferenceCollector()
	var val_1 ReferenceCollector
	InitReferenceCollector_3(&val_1, this_.compiler, NewInlineObjectLiterals_InliningBehavior(this_), NewSyntacticScopeCreator_1(this_.compiler))
	var callback *ReferenceCollector = &val_1
	_ = callback
	callback.Process_2(externs, root)
}

type InlineObjectLiterals_InliningBehavior struct {
	outer_    *InlineObjectLiterals
	staleVars j2gort.CustomOrderedSet[*Var]
}

func (this_ *InlineObjectLiterals_InliningBehavior) initFields_InlineObjectLiterals_InliningBehavior() {
	this_.staleVars = j2gort.NewCustomOrderedSet[*Var]()
}

func NewInlineObjectLiterals_InliningBehavior(outer_ *InlineObjectLiterals) *InlineObjectLiterals_InliningBehavior {
	this_ := &InlineObjectLiterals_InliningBehavior{}
	this_.outer_ = outer_
	InitInlineObjectLiterals_InliningBehavior(this_)
	return this_
}

func InitInlineObjectLiterals_InliningBehavior(this_ *InlineObjectLiterals_InliningBehavior) {
	this_.initFields_InlineObjectLiterals_InliningBehavior()
}

func (this_ *InlineObjectLiterals_InliningBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	_ = t
	_ = referenceMap
	cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(t.GetScope().GetVarIterable()))
	for cur_3.HasNext() {
		rawElem_2 := cur_3.Next()
		v := j2gort.CastPtr[Var](rawElem_2)
		_ = v
		if this_.priv_InliningBehavior_isVarInlineForbidden_1(v) {
			continue
		}
		var referenceInfo *ReferenceCollection = j2gort.CastPtr[ReferenceCollection](referenceMap.GetReferences(j2gort.PtrToAny(v)))
		_ = referenceInfo
		if referenceInfo == nil {
			continue
		}
		if this_.priv_InliningBehavior_isInlinableObject_2(referenceInfo) {
			this_.staleVars.Add(v)
			var init *Reference = referenceInfo.GetInitializingReference()
			_ = init
			this_.priv_InliningBehavior_splitObject_3(v, init, referenceInfo)
		}
	}
}

func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_isVarInlineForbidden_1(v *Var) bool {
	return v.IsGlobal() || v.IsExtern() || this_.outer_.compiler.GetCodingConvention().IsExported_2(v.GetName(), true) || this_.staleVars.Contains(v)
}

func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_isInlinableObject_2(refs *ReferenceCollection) bool {
	init := refs.GetInitializingReference()
	if init == nil {
		return false
	}
	val := NodeUtil_GetRValueOfLValue(init.GetNode())
	return val != 0 && val.IsObjectLit()
}

func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_splitObject_3(v *Var, init *Reference, _ *ReferenceCollection) {
	val := NodeUtil_GetRValueOfLValue(init.GetNode())
	NodeUtil_VisitPreOrder(val, NewInlineObjectLiterals__Anon_54(this_, j2gort.CastPtr[Scope](v.GetScope())))
	this_.outer_.SplitNames = append(this_.outer_.SplitNames, v.GetName())
}

type InlineObjectLiterals__Anon_54 struct {
	outer_      *InlineObjectLiterals_InliningBehavior
	cap_0_scope *Scope
}

func NewInlineObjectLiterals__Anon_54(outer_ *InlineObjectLiterals_InliningBehavior, cap_p_scope *Scope) *InlineObjectLiterals__Anon_54 {
	this_ := &InlineObjectLiterals__Anon_54{}
	this_.outer_ = outer_
	this_.cap_0_scope = cap_p_scope
	return this_
}

func (this_ *InlineObjectLiterals__Anon_54) Visit(node pkg_rhino.Node) {
	if node.IsName() {
		if sv := j2gort.CastPtr[Var](this_.cap_0_scope.GetVar(node.GetString())); sv != nil {
			this_.outer_.staleVars.Add(sv)
		}
	}
}

// GuardedCallback + PolyfillUsageFinder:
type GuardedCallback interface {
	NodeTraversal_Callback
	VisitGuarded(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type GuardedCallbackBase struct {
	self_        GuardedCallback
	compiler     AbstractCompiler
	contextStack j2gort.Deque[*GuardedCallback_Context]
}

func (*GuardedCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (g *GuardedCallbackBase) IsJavaNil_() bool                        { return g == nil }
func (g *GuardedCallbackBase) BindSelf_(s any)                         { g.self_ = s.(GuardedCallback) }

func InitGuardedCallback(this_ *GuardedCallbackBase, compiler AbstractCompiler) {
	this_.compiler = compiler
	this_.contextStack = j2gort.Deque[*GuardedCallback_Context]{}
}

func (this_ *GuardedCallbackBase) ShouldTraverse_3(_ *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	if this_.contextStack.IsEmpty() {
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

type GuardedCallback_Context struct {
	conditional pkg_rhino.Node
	safe        bool
	linked      *GuardedCallback_Context
}

var (
	static_GuardedCallback_Context_EMPTY *GuardedCallback_Context
	clinitDone_GuardedCallback_Context   atomic.Bool
)

func clinit_GuardedCallback_Context() {
	if clinitDone_GuardedCallback_Context.CompareAndSwap(false, true) {
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

func NewGuardedCallback_Context(conditional pkg_rhino.Node, safe bool, linked *GuardedCallback_Context) *GuardedCallback_Context {
	clinit_GuardedCallback_Context()
	this_ := &GuardedCallback_Context{}
	InitGuardedCallback_Context(this_, conditional, safe, linked)
	return this_
}

func InitGuardedCallback_Context(this_ *GuardedCallback_Context, conditional pkg_rhino.Node, safe bool, linked *GuardedCallback_Context) {
	this_.conditional = conditional
	this_.safe = safe
	this_.linked = linked
}

func (this_ *GuardedCallback_Context) Link(cond pkg_rhino.Node, safe bool) *GuardedCallback_Context {
	return NewGuardedCallback_Context(cond, safe, this_)
}

func (this_ *GuardedCallback_Context) Propagate(safe bool) *GuardedCallback_Context {
	if safe == this_.safe {
		return this_
	}
	return NewGuardedCallback_Context(this_.conditional, safe, this_.linked)
}

func (this_ *GuardedCallback_Context) Descend(_ AbstractCompiler, parent pkg_rhino.Node, child pkg_rhino.Node) *GuardedCallback_Context {
	if parent.IsIf() {
		if child == parent.GetFirstChild() {
			return this_.Link(child, false)
		}
		return this_.Propagate(true)
	}
	if parent.IsAnd() {
		return this_.Link(child, child != parent.GetFirstChild())
	}
	return this_
}

type PolyfillUsageFinder struct {
	compiler AbstractCompiler
}

func (this_ *PolyfillUsageFinder) priv_PolyfillUsageFinder_isVisible_1(t *NodeTraversal, name string) bool {
	return j2gort.CastPtr[Var](t.GetScope().GetVar(name)) == nil
}

func (this_ *PolyfillUsageFinder) TraverseOnlyGuarded(root pkg_rhino.Node, consumer j2gort.Ext_java_util_function_Consumer) {
	NodeTraversal_Traverse(this_.compiler, root, NewPolyfillUsageFinder_Traverser(this_, this_.compiler, consumer))
}

type PolyfillUsageFinder_Traverser struct {
	GuardedCallbackBase
	outer_   *PolyfillUsageFinder
	consumer j2gort.Ext_java_util_function_Consumer
}

func NewPolyfillUsageFinder_Traverser(outer_ *PolyfillUsageFinder, comp AbstractCompiler, consumer j2gort.Ext_java_util_function_Consumer) *PolyfillUsageFinder_Traverser {
	t := &PolyfillUsageFinder_Traverser{outer_: outer_, consumer: consumer}
	t.BindSelf_(t)
	InitGuardedCallback(&t.GuardedCallbackBase, comp)
	return t
}

func (this_ *PolyfillUsageFinder_Traverser) VisitGuarded(traversal *NodeTraversal, node pkg_rhino.Node, _ pkg_rhino.Node) {
	if node.IsName() && this_.outer_.priv_PolyfillUsageFinder_isVisible_1(traversal, node.GetString()) {
		if this_.consumer != nil {
			this_.consumer.Accept(node.GetString())
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "passes.go"), []byte(passesSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// Run 1:
	stats1, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("jscomp Run 1 failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] < 12 || stats1.ByPass["node_traversal_lowering"] < 8 {
		t.Fatalf("unexpected Run 1 stats: %+v", stats1.ByPass)
	}

	// Run 2 MUST be 100% idempotent (0 sites, 0 modified files)!
	stats2, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("jscomp Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites & 0 modified files on Run 2, got %+v", stats2)
	}

	ntBytes, err := os.ReadFile(filepath.Join(jscompDir, "node_traversal.go"))
	if err != nil {
		t.Fatal(err)
	}
	ntLowered := string(ntBytes)
	for _, want := range []string{
		"func unwrapNodeTraversalCallback_(cb NodeTraversal_Callback) any",
		"func (this_ *NodeTraversal__Anon_66) UnwrapCallback_() any",
		"traverseScopeRootPooled_(compiler, scb, scopeCreator, scopeNode, canRel_)",
		"if canRel_ {\n\t\tt_.canReleaseScopes_ = true\n\t}",
		"newReleasableLambda_NodeTraversal_AbstractPostOrderCallbackInterface",
	} {
		if !strings.Contains(ntLowered, want) {
			t.Fatalf("missing %q in lowered node_traversal.go:\n%s", want, ntLowered)
		}
	}
	if strings.Contains(ntLowered, "func (*NodeTraversal__Anon_66) CanReleaseSyntacticScopes_()") {
		t.Fatalf("UNSOUND: NodeTraversal__Anon_66 must NOT implement CanReleaseSyntacticScopes_() directly!")
	}

	coreBytes, err := os.ReadFile(filepath.Join(jscompDir, "core.go"))
	if err != nil {
		t.Fatal(err)
	}
	coreLowered := string(coreBytes)
	if strings.Count(coreLowered, "defer r_.ReleaseSlabs_()") < 2 {
		t.Fatalf("expected defer r_.ReleaseSlabs_() in both Compiler.GetAstDotGraph and AbstractCommandLineRunnerBase.ProcessResults:\n%s", coreLowered)
	}

	passesBytes, err := os.ReadFile(filepath.Join(jscompDir, "passes.go"))
	if err != nil {
		t.Fatal(err)
	}
	pLowered := string(passesBytes)
	for _, want := range []string{
		"func (*MakeDeclaredNamesUnique) CanReleaseSyntacticScopes_()",
		"func (*EphemeralPostOrderCb) CanReleaseSyntacticScopes_()",
		"func (*FlowSensitiveInlineVariables) CanReleaseSyntacticScopes_()",
		"func (*DeadAssignmentsElimination) CanReleaseSyntacticScopes_()",
		"func (*PolyfillUsageFinder_Traverser) CanReleaseSyntacticScopes_()",
		"this_.candidates = nil",
		"this_.reachingDef = nil",
		"this_.reachingUses = nil",
		"poppedColoring_ := this_.colorings.Pop()",
		"poppedColoring_.GetGraph()",
		"this_.liveness = nil\n\t\tscopeCreator.ReleaseSlabs_()",
		"slab_InlineObjectLiterals_InliningBehavior",
		"slab_InlineObjectLiterals__Anon_54",
		"func (this_ *InlineObjectLiterals) ReleaseSlabs_()",
		"this_.slab_GuardedCallback_Context.Reset()",
		"func (this_ *GuardedCallbackBase) ReleaseSlabs_()",
		"defer tr_.ReleaseSlabs_()",
	} {
		if !strings.Contains(pLowered, want) {
			t.Fatalf("missing %q in lowered passes.go:\n%s", want, pLowered)
		}
	}
	for _, forbidden := range []string{
		"func (*MakeDeclaredNamesUnique_ContextualRenameInverter) CanReleaseSyntacticScopes_()",
		"func (*EscapingPostOrderCb) CanReleaseSyntacticScopes_()",
	} {
		if strings.Contains(pLowered, forbidden) {
			t.Fatalf("UNSOUND marker %q found in lowered passes.go:\n%s", forbidden, pLowered)
		}
	}

	// End-to-end runtime + -race verification!
	rtTestSrc := `package jscomp

import (
	"reflect"
	"sync"
	"testing"

	"example.com/w23ctest/gen/j2gort"
	pkg_rhino "example.com/w23ctest/gen/rhino"
)

func buildTestAST() (externs, script, fn1, fn2 pkg_rhino.Node) {
	pkg_rhino.ResetArena()
	externs = pkg_rhino.AllocNode(pkg_rhino.TokSCRIPT, "externs")
	extVar := pkg_rhino.AllocNode(pkg_rhino.TokVAR, "")
	extName := pkg_rhino.AllocNode(pkg_rhino.TokNAME, "extObj")
	pkg_rhino.AppendChild(extName, pkg_rhino.AllocNode(pkg_rhino.TokOBJECTLIT, ""))
	pkg_rhino.AppendChild(extVar, extName)
	pkg_rhino.AppendChild(externs, extVar)

	script = pkg_rhino.AllocNode(pkg_rhino.TokSCRIPT, "main")
	globVar := pkg_rhino.AllocNode(pkg_rhino.TokVAR, "")
	globName := pkg_rhino.AllocNode(pkg_rhino.TokNAME, "globObj")
	pkg_rhino.AppendChild(globName, pkg_rhino.AllocNode(pkg_rhino.TokOBJECTLIT, ""))
	pkg_rhino.AppendChild(globVar, globName)
	pkg_rhino.AppendChild(script, globVar)

	makeFn := func(fnName string, extraVars int) pkg_rhino.Node {
		fn := pkg_rhino.AllocNode(pkg_rhino.TokFUNCTION, fnName)
		blk := pkg_rhino.AllocNode(pkg_rhino.TokBLOCK, "")
		// Inlinable object literal: var obj = {};
		vDecl := pkg_rhino.AllocNode(pkg_rhino.TokVAR, "")
		vName := pkg_rhino.AllocNode(pkg_rhino.TokNAME, "obj")
		objLit := pkg_rhino.AllocNode(pkg_rhino.TokOBJECTLIT, "")
		pkg_rhino.AppendChild(objLit, pkg_rhino.AllocNode(pkg_rhino.TokNAME, "capturedStale"))
		pkg_rhino.AppendChild(vName, objLit)
		pkg_rhino.AppendChild(vDecl, vName)
		pkg_rhino.AppendChild(blk, vDecl)

		// Stale var referenced inside objLit: var capturedStale = {};
		sDecl := pkg_rhino.AllocNode(pkg_rhino.TokVAR, "")
		sName := pkg_rhino.AllocNode(pkg_rhino.TokNAME, "capturedStale")
		pkg_rhino.AppendChild(sName, pkg_rhino.AllocNode(pkg_rhino.TokOBJECTLIT, ""))
		pkg_rhino.AppendChild(sDecl, sName)
		pkg_rhino.AppendChild(blk, sDecl)

		// Non-object vars (must be filtered out by inlineObjectLiteralsVarFilter_!):
		for i := 0; i < extraVars; i++ {
			nd := pkg_rhino.AllocNode(pkg_rhino.TokVAR, "")
			nn := pkg_rhino.AllocNode(pkg_rhino.TokNAME, "num_"+string(rune('a'+i)))
			pkg_rhino.AppendChild(nn, pkg_rhino.AllocNode(pkg_rhino.TokNUMBER, "1"))
			pkg_rhino.AppendChild(nd, nn)
			pkg_rhino.AppendChild(blk, nd)
		}
		pkg_rhino.AppendChild(fn, blk)
		return fn
	}

	fn1 = makeFn("f1", 2) // <= 6 vars (var0..var5 path)
	fn2 = makeFn("f2", 7) // > 6 vars (vars OrderedMap path)
	pkg_rhino.AppendChild(script, fn1)
	pkg_rhino.AppendChild(script, fn2)
	return
}

func TestWave23CRuntimeBehaviorAndLeaks(t *testing.T) {
	externs, script, fn1, fn2 := buildTestAST()
	comp := NewCompiler(script)

	// 1. MakeDeclaredNamesUnique via TraverseScopeRoots releases MemoizedScopeCreator + Scope/Var slabs!
	j2gort.ResetMetrics()
	roots := j2gort.SliceIterable{Items: []any{uint32(fn1), uint32(fn2)}}
	for i := 0; i < 15; i++ {
		NodeTraversal_TraverseScopeRoots(comp, roots, NewMakeDeclaredNamesUnique(), false)
	}
	if j2gort.OrderedMapReleases.Load() < 15 || j2gort.SlabReleaseWithCalls.Load() < 15 {
		t.Fatalf("expected MakeDeclaredNamesUnique + TraverseScopeRoots to release MemoizedScopeCreator & Scope slabs, got mapRel=%d slabRelWith=%d",
			j2gort.OrderedMapReleases.Load(), j2gort.SlabReleaseWithCalls.Load())
	}

	// PostOrder callback via Builder.SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface:
	j2gort.ResetMetrics()
	ephemPO := &EphemeralPostOrderCb{}
	NodeTraversal_S_Builder().SetCompiler(comp).SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(ephemPO).Traverse(script)
	if j2gort.SlabReleaseWithCalls.Load() == 0 {
		t.Fatalf("expected EphemeralPostOrderCb wrapped in NodeTraversal__Anon_66 to release scopes!")
	}

	j2gort.ResetMetrics()
	escPO := &EscapingPostOrderCb{}
	NodeTraversal_S_Builder().SetCompiler(comp).SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(escPO).Traverse(script)
	if j2gort.SlabReleaseWithCalls.Load() != 0 || escPO.savedScope == nil {
		t.Fatalf("EscapingPostOrderCb must NOT release scopes (slabRelWith=%d savedScope=%v)",
			j2gort.SlabReleaseWithCalls.Load(), escPO.savedScope)
	}

	j2gort.ResetMetrics()
	hoister := &J2clConstantHoisterPass{compiler: comp}
	hoister.Process_2(externs, script)
	if hoister.hoisted == 0 || j2gort.SlabReleaseWithCalls.Load() == 0 {
		t.Fatalf("expected J2clConstantHoisterPass lambda to execute & release scopes (hoisted=%d, rel=%d)",
			hoister.hoisted, j2gort.SlabReleaseWithCalls.Load())
	}

	// 2. FlowSensitiveInlineVariables, Compiler.GetAstDotGraph, ProcessResults, CoalesceVariableNames:
	j2gort.ResetMetrics()
	fsiv := &FlowSensitiveInlineVariables{compiler: comp}
	fsiv.Process_2(externs, script)
	if fsiv.inlinedFns != 2 || fsiv.cfg != nil || fsiv.candidates != nil || j2gort.SlabReleaseCalls.Load() < 2 || j2gort.SlabReleaseWithCalls.Load() == 0 {
		t.Fatalf("FlowSensitiveInlineVariables leak: inlined=%d cfg=%v rel=%d relWith=%d",
			fsiv.inlinedFns, fsiv.cfg, j2gort.SlabReleaseCalls.Load(), j2gort.SlabReleaseWithCalls.Load())
	}

	j2gort.ResetMetrics()
	if dot := comp.GetAstDotGraph(); dot != "digraph{}" || j2gort.SlabReleaseCalls.Load() < 2 {
		t.Fatalf("GetAstDotGraph failed to release CFG: dot=%q rel=%d", dot, j2gort.SlabReleaseCalls.Load())
	}
	runner := &AbstractCommandLineRunnerBase{compiler: comp, printAst: true}
	if rc := runner.ProcessResults(); rc != 0 || j2gort.SlabReleaseCalls.Load() < 4 {
		t.Fatalf("ProcessResults failed to release CFG: rc=%d rel=%d", rc, j2gort.SlabReleaseCalls.Load())
	}

	j2gort.ResetMetrics()
	cvn := NewCoalesceVariableNames(comp)
	NodeTraversal_Traverse(comp, script, cvn)
	if j2gort.OrderedMapReleases.Load() < 4 {
		t.Fatalf("CoalesceVariableNames failed to release CFG + LinkedUndirectedGraph maps, got %d", j2gort.OrderedMapReleases.Load())
	}

	// 3. DeadAssignmentsElimination releases BOTH inner scopeCreator AND outer defaultScopeCreator_!
	j2gort.ResetMetrics()
	dae := NewDeadAssignmentsElimination(comp)
	dae.Process_2(externs, script)
	if dae.eliminated != 2 || dae.liveness != nil || j2gort.SlabReleaseWithCalls.Load() < 3 {
		t.Fatalf("DeadAssignmentsElimination leak: eliminated=%d liveness=%v relWith=%d",
			dae.eliminated, dae.liveness, j2gort.SlabReleaseWithCalls.Load())
	}

	// 4. InlineObjectLiterals mechanical PooledSlab pooling + staleVars parity:
	j2gort.ResetMetrics()
	iol := NewInlineObjectLiterals(comp)
	iol.Process_2(externs, script)
	// In each function, "obj" is split first and records "capturedStale" in staleVars so "capturedStale" is NOT split!
	wantSplit := []string{"obj", "obj"}
	if !reflect.DeepEqual(iol.SplitNames, wantSplit) {
		t.Fatalf("InlineObjectLiterals SplitNames = %v, want %v", iol.SplitNames, wantSplit)
	}
	if j2gort.SlabReleaseCalls.Load() < 2 {
		t.Fatalf("expected InlineObjectLiterals.Process_2 to defer ReleaseSlabs_() on InliningBehavior & Anon_54 slabs, got %d", j2gort.SlabReleaseCalls.Load())
	}

	// 5. GuardedCallbackBase Reset() vs PolyfillUsageFinder.TraverseOnlyGuarded defer ReleaseSlabs_():
	ifTree := pkg_rhino.AllocNode(pkg_rhino.TokSCRIPT, "guarded")
	for k := 0; k < 10; k++ {
		ifN := pkg_rhino.AllocNode(pkg_rhino.TokIF, "")
		andN := pkg_rhino.AllocNode(pkg_rhino.TokAND, "")
		pkg_rhino.AppendChild(andN, pkg_rhino.AllocNode(pkg_rhino.TokNAME, "Promise"))
		pkg_rhino.AppendChild(andN, pkg_rhino.AllocNode(pkg_rhino.TokNAME, "Symbol"))
		pkg_rhino.AppendChild(ifN, andN)
		pkg_rhino.AppendChild(ifN, pkg_rhino.AllocNode(pkg_rhino.TokBLOCK, ""))
		pkg_rhino.AppendChild(ifTree, ifN)
	}
	j2gort.ResetMetrics()
	puf := &PolyfillUsageFinder{compiler: comp}
	puf.TraverseOnlyGuarded(ifTree, nil)
	if resets := j2gort.SlabResetCalls.Load(); resets != 20 {
		t.Fatalf("expected GuardedCallbackBase.Visit_3 to call PooledSlab.Reset() after each of the 20 top-level slab subtrees (slabCtxCount_ == 0), got %d", resets)
	}
	// Release() on GuardedCallback_Context slab runs ONCE via defer tr_.ReleaseSlabs_()!
	if rels := j2gort.SlabReleaseCalls.Load(); rels < 1 {
		t.Fatalf("expected PolyfillUsageFinder.TraverseOnlyGuarded to defer tr_.ReleaseSlabs_(), got %d", rels)
	}
	if StaticGet_GuardedCallback_Context_EMPTY().safe {
		t.Fatalf("EMPTY singleton corrupted!")
	}

	// Concurrent -race stress test:
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := NewCompiler(script)
			for iter := 0; iter < 10; iter++ {
				NodeTraversal_TraverseScopeRoots(c, roots, NewMakeDeclaredNamesUnique(), false)
				f := &FlowSensitiveInlineVariables{compiler: c}
				f.Process_2(externs, script)
				d := NewDeadAssignmentsElimination(c)
				d.Process_2(externs, script)
				o := NewInlineObjectLiterals(c)
				o.Process_2(externs, script)
				pf := &PolyfillUsageFinder{compiler: c}
				pf.TraverseOnlyGuarded(ifTree, nil)
			}
		}()
	}
	wg.Wait()
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "wave23c_rt_test.go"), []byte(rtTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-timeout=20s", "-race", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race in wave23c workspace failed: %v\n%s", err, string(out))
	}
}

func TestWave23C_RealGenDryRunAndIdempotency(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	j2Cmd := exec.Command("go", "build", "./gen/j2gort")
	j2Cmd.Dir = filepath.Join(repoRoot, "transpiled")
	if j2Cmd.Run() != nil {
		t.Skip("gen/j2gort is currently mid-edit; skipping real gen/ DryRun")
	}

	colRtSrc, _ := os.ReadFile(filepath.Join(repoRoot, "transpiled", "gen", "j2gort", "collections_runtime.go"))
	extRtSrc, _ := os.ReadFile(filepath.Join(repoRoot, "transpiled", "gen", "j2gort", "ext_collections_rt.go"))
	ovRtSrc, _ := os.ReadFile(filepath.Join(repoRoot, "transpiled", "gen", "j2gort", "lowering_helpers_overlay.go"))
	j2AllSrc := string(colRtSrc) + "\n" + string(extRtSrc) + "\n" + string(ovRtSrc)
	hasSlabResetInGen := strings.Contains(j2AllSrc, "func (s *PooledSlab[T]) Reset()")
	hasReleaseOrderedMapInGen := strings.Contains(j2AllSrc, "func ReleaseOrderedMap[")
	hasPooledSliceInGen := strings.Contains(j2AllSrc, "func ReleasePooledSlice[")
	hasTrimStaleInGen := strings.Contains(j2AllSrc, "func TrimStaleSlabPools(")
	hasPurgeInGen := strings.Contains(j2AllSrc, "func PurgeAllRuntimePools(")
	nsDir := filepath.Join(repoRoot, "transpiled", "shims", "nodestore")
	nsEntries, _ := os.ReadDir(nsDir)
	var nsAll strings.Builder
	for _, ne := range nsEntries {
		if !ne.IsDir() && strings.HasSuffix(ne.Name(), ".go") && !strings.HasSuffix(ne.Name(), "_test.go") {
			b, _ := os.ReadFile(filepath.Join(nsDir, ne.Name()))
			nsAll.Write(b)
			nsAll.WriteByte('\n')
		}
	}
	hasRelJSTypesInNS := strings.Contains(nsAll.String(), "func ReleaseCurrentStoreJSTypes(")

	realGraphDir := filepath.Join(repoRoot, "transpiled", "gen", "graph")
	lugSrc, _ := os.ReadFile(filepath.Join(realGraphDir, "linked_undirected_graph.go"))
	ldgSrc, _ := os.ReadFile(filepath.Join(realGraphDir, "linked_directed_graph.go"))
	graphPreLowered := strings.Contains(string(lugSrc), "func (this_ *LinkedUndirectedGraph[N, E]) ReleaseSlabs_()") &&
		(!hasSlabResetInGen || strings.Contains(string(ldgSrc), ".slab_LinkedDirectedGraph_LinkedDiGraphNode.Reset()")) &&
		strings.Contains(string(ldgSrc), "this_.dfaEdgeIndex = int32(-1)") &&
		(!hasPooledSliceInGen || strings.Contains(string(ldgSrc), "ReleasePooledSlice(n_.inEdges)"))

	gDry, err := passes.RunOnPackageDir(
		realGraphDir,
		[]passes.Pass{passes.NewSlabAllocLoweringPass()},
		passes.Options{DryRun: true},
	)
	if err != nil {
		t.Fatalf("DryRun slab_alloc_lowering on real gen/graph failed: %v", err)
	}
	if !graphPreLowered && gDry.ByPass["slab_alloc_lowering"] < 1 {
		t.Fatalf("expected >= 1 slab_alloc_lowering sites on real gen/graph, got %+v", gDry.ByPass)
	}
	if graphPreLowered && gDry.ByPass["slab_alloc_lowering"] != 0 {
		t.Fatalf("expected 0 slab_alloc_lowering sites on pre-lowered gen/graph, got %+v", gDry.ByPass)
	}

	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/rhino is awaiting regeneration after shims/nodestore edit; skipping real gen/jscomp DryRun")
	}

	realJscompDir := filepath.Join(repoRoot, "transpiled", "gen", "jscomp")
	ntSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "node_traversal.go"))
	gnSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "global_namespace.go"))
	iolSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "inline_object_literals.go"))
	gcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "guarded_callback.go"))
	mdnuSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "make_declared_names_unique.go"))
	sndSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "script_node_deserializer.go"))
	normSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "normalize.go"))
	rvSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "rename_vars.go"))
	sscSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "syntactic_scope_creator.go"))
	compSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "compiler.go"))
	tipSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "type_inference_pass.go"))
	tasSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "typed_ast_serializer.go"))
	cttcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "convert_types_to_colors.go"))
	rcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "reference_collector.go"))
	rucSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "remove_unused_code.go"))
	ccrcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "cross_chunk_reference_collector.go"))
	pmcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "peephole_minimize_conditions.go"))
	ocSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "optimize_calls.go"))
	infSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "inline_functions.go"))
	pdSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "process_defines.go"))
	dpaeSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "dead_property_assignment_elimination.go"))
	tscSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "typed_scope_creator.go"))
	ftbSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "function_type_builder.go"))
	tiSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "type_inference.go"))
	poSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "phase_optimizer.go"))
	cfgSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "control_flow_graph.go"))
	lvaSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "live_variables_analysis.go"))
	pfiSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "pure_function_identifier.go"))
	ivSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "inline_variables.go"))
	mcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "minimized_condition.go"))
	spSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "string_pool.go"))
	jscompFullyLowered := strings.Contains(string(iolSrc), "slab_InlineObjectLiterals_InliningBehavior") &&
		strings.Contains(string(sscSrc), "allocFastScopeChild_") &&
		strings.Contains(string(ivSrc), "scratchVarExpert_") &&
		strings.Contains(string(mcSrc), "lazyTok_") &&
		strings.Contains(string(spSrc), "InternOwned(") &&
		strings.Contains(string(infSrc), "slab_InlineFunctions_Reference") &&
		strings.Contains(string(dpaeSrc), "slab_DeadPropertyAssignmentElimination_Property") &&
		strings.Contains(string(cfgSrc), "slab_LinearFlowState") &&
		strings.Contains(string(lvaSrc), "slab_LiveVariableLattice") &&
		strings.Contains(string(pfiSrc), "slab_SideEffectPropagation") &&
		!strings.Contains(string(tscSrc), "slab_TypedScope") &&
		strings.Contains(string(tscSrc), "acquireNormalScopeBuilder_") &&
		strings.Contains(string(ftbSrc), "var fb_ pkg_rhino.FunctionType_Builder") &&
		strings.Contains(string(tiSrc), ".InitBuilder_(&fb_") &&
		strings.Contains(string(mdnuSrc), "CanReleaseSyntacticScopes_") &&
		strings.Contains(string(mdnuSrc), "slab_ContextualRenamer") &&
		strings.Contains(string(pmcSrc), "slab_MinimizedCondition_MeasuredNode") &&
		strings.Contains(string(ocSrc), "slab_OptimizeCalls_NodeList") &&
		strings.Contains(string(ocSrc), "NewSyntacticScopeCreator_1") &&
		strings.Contains(string(gcSrc), "slabCtxCount_") &&
		strings.Contains(string(ntSrc), "canReleaseSyntacticScopesForCallback_") &&
		strings.Contains(string(gnSrc), "func (this_ *GlobalNamespace) ReleaseSlabs_()") &&
		strings.Contains(string(gnSrc), "slab_GlobalNamespace_Ref") &&
		strings.Contains(string(rucSrc), "slab_RemoveUnusedCode_VarInfo") &&
		strings.Contains(string(ccrcSrc), "func (this_ *CrossChunkReferenceCollector) ReleaseSlabs_()") &&
		strings.Contains(string(sndSrc), "nodeAlloc_") &&
		strings.Contains(string(normSrc), "func (*Normalize_NormalizeStatements) CanReleaseSyntacticScopes_()") &&
		strings.Contains(string(rvSrc), "ReleaseSlabs_()") &&
		(!hasReleaseOrderedMapInGen || (strings.Contains(string(sscSrc), "ReleaseOrderedMap(") &&
			strings.Contains(string(gnSrc), "ReleaseOrderedMap("))) &&
		(!hasPooledSliceInGen || (strings.Contains(string(rcSrc), "ReleasePooledSlice(") &&
			strings.Contains(string(rucSrc), "ReleasePooledSlice("))) &&
		(!hasTrimStaleInGen || (strings.Contains(string(compSrc), "TrimStaleSlabPools(0)") &&
			strings.Contains(string(poSrc), "defer j2gort.TrimStaleSlabPools(14)"))) &&
		(!hasPurgeInGen || (strings.Count(string(compSrc), "PurgeAllRuntimePools(") == 1 &&
			strings.Count(string(tipSrc), "PurgeAllRuntimePools(") == 1 &&
			strings.Contains(string(tasSrc), "PurgeAllRuntimePools(") &&
			!strings.Contains(string(pdSrc), "PurgeAllRuntimePools("))) &&
		(!hasRelJSTypesInNS || strings.Contains(string(cttcSrc), "ReleaseCurrentStoreJSTypes(")) &&
		(!hasSlabResetInGen || strings.Contains(string(gcSrc), ".slab_GuardedCallback_Context.Reset()"))

	jDry, err := passes.RunOnPackageDir(
		realJscompDir,
		[]passes.Pass{passes.NewSlabAllocLoweringPass(), passes.NewNodeTraversalLoweringPass()},
		passes.Options{DryRun: true},
	)
	if err != nil {
		t.Fatalf("DryRun on real gen/jscomp failed: %v", err)
	}
	if !jscompFullyLowered && jDry.TotalSites < 1 {
		t.Fatalf("expected >= 1 sites on un-regenerated gen/jscomp, got %+v", jDry.ByPass)
	}
	if jscompFullyLowered && jDry.TotalSites != 0 {
		t.Fatalf("expected 0 sites on already-regenerated gen/jscomp, got %+v", jDry)
	}
}

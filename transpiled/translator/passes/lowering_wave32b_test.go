package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave32B_CfgSlabsPhasePurgesInlineFunctionsBuilderStackAndNodeLambdas(t *testing.T) {
	tmpDir := t.TempDir()
	goMod := "module example.com/w32btest\n\ngo 1.24\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	j2Dir := filepath.Join(tmpDir, "gen", "j2gort")
	rhinoDir := filepath.Join(tmpDir, "gen", "rhino")
	graphDir := filepath.Join(tmpDir, "gen", "graph")
	jscompDir := filepath.Join(tmpDir, "gen", "jscomp")
	jscompChecksBoolDir := filepath.Join(tmpDir, "gen", "jscomp_checksonly")
	jscompEscDir := filepath.Join(tmpDir, "gen", "jscomp_esc")
	for _, d := range []string{j2Dir, rhinoDir, graphDir, jscompDir, jscompChecksBoolDir, jscompEscDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	j2Src := `package j2gort

import "sync/atomic"

var (
	PurgeCalls          atomic.Int64
	SlabAllocCalls      atomic.Int64
	SlabResetCalls      atomic.Int64
	SlabReleaseCalls    atomic.Int64
	NodeLambdaCalls     atomic.Int64
	EventLog            []string
)

func ResetWave32BMetrics() {
	PurgeCalls.Store(0)
	SlabAllocCalls.Store(0)
	SlabResetCalls.Store(0)
	SlabReleaseCalls.Store(0)
	NodeLambdaCalls.Store(0)
	EventLog = nil
}

func PurgeAllRuntimePools() {
	PurgeCalls.Add(1)
	EventLog = append(EventLog, "purge")
}
func TrimStaleSlabPools(int) {}

func CheckState(cond bool) {
	if !cond {
		panic("CheckState")
	}
}
func CheckNotNullString(s string) string { return s }
func CheckNotNullNode[N ~uint32](n N) N {
	if n == 0 {
		panic("NPE")
	}
	return n
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
func AsIface[I any](v any) I {
	if v == nil {
		var z I
		return z
	}
	return v.(I)
}
func MakeNonNilSlice[T any](capHint int32) []T {
	if capHint <= 0 {
		return []T{}
	}
	return make([]T, 0, capHint)
}
func SliceRemoveAt[T any](s []T, idx int32) []T {
	last := len(s) - 1
	copy(s[idx:], s[idx+1:])
	var z T
	s[last] = z
	return s[:last]
}

type Slab[T any] struct{ items []T }

func (s *Slab[T]) Alloc() *T {
	s.items = append(s.items, *new(T))
	return &s.items[len(s.items)-1]
}
func (s *Slab[T]) Release() { s.items = nil }

type PooledSlab[T any] struct{ items []T }

func (s *PooledSlab[T]) Alloc() *T {
	SlabAllocCalls.Add(1)
	s.items = append(s.items, *new(T))
	return &s.items[len(s.items)-1]
}
func (s *PooledSlab[T]) Reset() {
	SlabResetCalls.Add(1)
	clear(s.items)
	s.items = s.items[:0]
}
func (s *PooledSlab[T]) Release() {
	SlabReleaseCalls.Add(1)
	EventLog = append(EventLog, "slab_release")
	clear(s.items)
	s.items = nil
}
func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	for i := range s.items {
		if fn != nil {
			fn(&s.items[i])
		}
	}
	s.Release()
}

type OrderedMap[K comparable, V any] struct {
	keys []K
	m    map[K]V
}
type OrderedMapValCursor[V any] struct {
	vals []V
	pos  int
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V)}
}
func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V)}
}
func AcquireOrderedMapCap[K comparable, V any](int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V)}
}
func AcquireOrderedMapPtrCap[K comparable, V any](int32) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{m: make(map[K]V)}
}
func (om *OrderedMap[K, V]) Put(k K, v V) {
	if om.m == nil {
		om.m = make(map[K]V)
	}
	if _, exists := om.m[k]; !exists {
		om.keys = append(om.keys, k)
	}
	om.m[k] = v
}
func (om OrderedMap[K, V]) Get(k K) V { return om.m[k] }
func (om OrderedMap[K, V]) Size() int32 { return int32(len(om.m)) }
func (om OrderedMap[K, V]) IsEmpty() bool { return len(om.m) == 0 }
func (om *OrderedMap[K, V]) Clear() {
	clear(om.m)
	om.keys = om.keys[:0]
}
func (om *OrderedMap[K, V]) ReleaseToPool() {
	om.m = nil
	om.keys = nil
}
func (om *OrderedMap[K, V]) IterValues() OrderedMapValCursor[V] {
	vals := make([]V, 0, len(om.keys))
	for _, k := range om.keys {
		vals = append(vals, om.m[k])
	}
	return OrderedMapValCursor[V]{vals: vals}
}
func (c *OrderedMapValCursor[V]) HasNext() bool { return c.pos < len(c.vals) }
func (c *OrderedMapValCursor[V]) Next() V {
	v := c.vals[c.pos]
	c.pos++
	return v
}
func ReleaseOrderedMap[K comparable, V any](om *OrderedMap[K, V]) {
	if om != nil {
		om.m = nil
		om.keys = nil
	}
}

type Deque[T any] struct{ items []T }

func (d *Deque[T]) Push(v T)    { d.items = append([]T{v}, d.items...) }
func (d *Deque[T]) Peek() T     { return d.items[0] }
func (d *Deque[T]) Pop() T      { v := d.items[0]; d.items = d.items[1:]; return v }
func (d *Deque[T]) IsEmpty() bool { return len(d.items) == 0 }

type NodePredicate interface {
	TestNode_(uint32) bool
}

type Ext_com_google_common_base_Predicate interface {
	Implements_j2gort_Ext_com_google_common_base_Predicate()
	IsJavaNil_() bool
	Apply(any) bool
}

type Ext_java_util_function_Predicate interface {
	Implements_j2gort_Ext_java_util_function_Predicate()
	IsJavaNil_() bool
	Test(any) bool
}

type nodeGuavaPred struct{ fn func(uint32) bool }

func (*nodeGuavaPred) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (p *nodeGuavaPred) IsJavaNil_() bool                                      { return p == nil }
func (p *nodeGuavaPred) TestNode_(n uint32) bool                               { return p.fn(n) }
func (p *nodeGuavaPred) Apply(v any) bool                                      { return p.fn(v.(uint32)) }

func NewNodeLambda_Ext_com_google_common_base_Predicate(fn func(uint32) bool) Ext_com_google_common_base_Predicate {
	NodeLambdaCalls.Add(1)
	return &nodeGuavaPred{fn: fn}
}

type nodeJavaPred struct{ fn func(uint32) bool }

func (*nodeJavaPred) Implements_j2gort_Ext_java_util_function_Predicate() {}
func (p *nodeJavaPred) IsJavaNil_() bool                                  { return p == nil }
func (p *nodeJavaPred) TestNode_(n uint32) bool                           { return p.fn(n) }
func (p *nodeJavaPred) Test(v any) bool                                   { return p.fn(v.(uint32)) }

func NewNodeLambda_Ext_java_util_function_Predicate(fn func(uint32) bool) Ext_java_util_function_Predicate {
	NodeLambdaCalls.Add(1)
	return &nodeJavaPred{fn: fn}
}

var alwaysTruePred Ext_com_google_common_base_Predicate = &nodeGuavaPred{fn: func(uint32) bool { return true }}

func Ext_com_google_common_base_Predicates_AlwaysTrue() Ext_com_google_common_base_Predicate {
	return alwaysTruePred
}

func PredTestNode(p any, n uint32) bool {
	if p == alwaysTruePred {
		return true
	}
	if np, ok := p.(NodePredicate); ok {
		return np.TestNode_(n)
	}
	return p.(Ext_com_google_common_base_Predicate).Apply(n)
}
`
	if err := os.WriteFile(filepath.Join(j2Dir, "j2gort.go"), []byte(j2Src), 0o644); err != nil {
		t.Fatal(err)
	}

	rhinoSrc := `package rhino

type Token uint8

const (
	TokROOT     Token = 1
	TokSCRIPT   Token = 2
	TokFUNCTION Token = 3
	TokBLOCK    Token = 4
	TokNAME     Token = 5
	TokRETURN   Token = 6
	TokSUPER    Token = 7
)

type Node uint32

type nodeEntry struct {
	tok      Token
	str      string
	parent   Node
	children []Node
}

var nodes = []nodeEntry{{}}

func ResetNodes() { nodes = []nodeEntry{{}} }

func MakeNode(tok Token, str string, children ...Node) Node {
	id := Node(len(nodes))
	nodes = append(nodes, nodeEntry{tok: tok, str: str, children: children})
	for _, ch := range children {
		nodes[ch].parent = id
	}
	return id
}

func (n Node) GetToken() Token   { return nodes[n].tok }
func (n Node) RawToken_() uint8  { return uint8(nodes[n].tok) }
func (n Node) GetString() string { return nodes[n].str }
func (n Node) IsFunction() bool  { return n != 0 && nodes[n].tok == TokFUNCTION }
func (n Node) IsName() bool      { return n != 0 && nodes[n].tok == TokNAME }
func (n Node) IsReturn() bool    { return n != 0 && nodes[n].tok == TokRETURN }
func (n Node) IsSuper() bool     { return n != 0 && nodes[n].tok == TokSUPER }
func (n Node) GetFirstChild() Node {
	if n == 0 || len(nodes[n].children) == 0 {
		return 0
	}
	return nodes[n].children[0]
}
func (n Node) GetNext() Node {
	p := nodes[n].parent
	if p == 0 {
		return 0
	}
	sibs := nodes[p].children
	for i, c := range sibs {
		if c == n && i+1 < len(sibs) {
			return sibs[i+1]
		}
	}
	return 0
}
`
	if err := os.WriteFile(filepath.Join(rhinoDir, "rhino.go"), []byte(rhinoSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1A: LinkedDirectedGraphBase with ResetForReuse (freeNodes/freeEdges loop -> PooledSlab.Reset()).
	graphSrc := `package graph

import "example.com/w32btest/gen/j2gort"

type LinkedDirectedGraph[N comparable, E any] interface {
	IsJavaNil_() bool
	CreateNode(N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E]
	Connect(N, E, N)
	ResetForReuse(bool, bool) bool
}

type LinkedDirectedGraph_LinkedDiGraphNode[N comparable, E any] struct {
	value          N
	firstOutEdge   *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]
	outEdges       []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]
	useAnnotation  bool
}

func NewLinkedDirectedGraph_LinkedDiGraphNode[N comparable, E any](v N, useAnn bool) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	n := &LinkedDirectedGraph_LinkedDiGraphNode[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphNode(n, v, useAnn)
	return n
}
func InitLinkedDirectedGraph_LinkedDiGraphNode[N comparable, E any](n *LinkedDirectedGraph_LinkedDiGraphNode[N, E], v N, useAnn bool) {
	n.value = v
	n.useAnnotation = useAnn
}
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) Reinit(v N, useAnn bool) {
	n.value = v
	n.useAnnotation = useAnn
}
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) ClearForPool() {
	var z N
	n.value = z
	n.firstOutEdge = nil
	n.outEdges = nil
}

type LinkedDirectedGraph_LinkedDiGraphEdge[N comparable, E any] struct {
	src           *LinkedDirectedGraph_LinkedDiGraphNode[N, E]
	value         E
	dest          *LinkedDirectedGraph_LinkedDiGraphNode[N, E]
	useAnnotation bool
}

func NewLinkedDirectedGraph_LinkedDiGraphEdge[N comparable, E any](src *LinkedDirectedGraph_LinkedDiGraphNode[N, E], v E, dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E], useAnn bool) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	e := &LinkedDirectedGraph_LinkedDiGraphEdge[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphEdge(e, src, v, dest, useAnn)
	return e
}
func InitLinkedDirectedGraph_LinkedDiGraphEdge[N comparable, E any](e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E], src *LinkedDirectedGraph_LinkedDiGraphNode[N, E], v E, dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E], useAnn bool) {
	e.src = src
	e.value = v
	e.dest = dest
	e.useAnnotation = useAnn
}
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) Reinit(src *LinkedDirectedGraph_LinkedDiGraphNode[N, E], v E, dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E], useAnn bool) {
	e.src = src
	e.value = v
	e.dest = dest
	e.useAnnotation = useAnn
}
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) ClearForPool() {
	var z E
	e.src = nil
	e.dest = nil
	e.value = z
}

type LinkedDirectedGraphBase[N comparable, E any] struct {
	nodes              j2gort.OrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]
	freeNodes          []*LinkedDirectedGraph_LinkedDiGraphNode[N, E]
	freeEdges          []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]
	useNodeAnnotations bool
	useEdgeAnnotations bool
}

func (this_ *LinkedDirectedGraphBase[N, E]) IsJavaNil_() bool { return this_ == nil }

func (this_ *LinkedDirectedGraphBase[N, E]) initFields_LinkedDirectedGraph() {
	this_.nodes = j2gort.NewOrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]()
}

func InitLinkedDirectedGraph[N comparable, E any](this_ *LinkedDirectedGraphBase[N, E], useNodeAnn, useEdgeAnn bool) {
	this_.initFields_LinkedDirectedGraph()
	this_.useNodeAnnotations = useNodeAnn
	this_.useEdgeAnnotations = useEdgeAnn
}

func NewLinkedDirectedGraph[N comparable, E any](useNodeAnn, useEdgeAnn bool) *LinkedDirectedGraphBase[N, E] {
	g := &LinkedDirectedGraphBase[N, E]{}
	InitLinkedDirectedGraph(g, useNodeAnn, useEdgeAnn)
	return g
}

func (this_ *LinkedDirectedGraphBase[N, E]) CreateNode(v N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	n := this_.nodes.Get(v)
	if n == nil {
		if this_.freeNodes != nil && len(this_.freeNodes) != 0 {
			last := int32(len(this_.freeNodes) - 1)
			n = this_.freeNodes[last]
			this_.freeNodes = j2gort.SliceRemoveAt(this_.freeNodes, last)
			n.Reinit(v, this_.useNodeAnnotations)
		} else {
			n = NewLinkedDirectedGraph_LinkedDiGraphNode[N, E](v, this_.useNodeAnnotations)
		}
		this_.nodes.Put(v, n)
	}
	return n
}

func (this_ *LinkedDirectedGraphBase[N, E]) Connect(s N, v E, d N) {
	sn := this_.CreateNode(s)
	dn := this_.CreateNode(d)
	var edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]
	if this_.freeEdges != nil && len(this_.freeEdges) != 0 {
		last := int32(len(this_.freeEdges) - 1)
		edge = this_.freeEdges[last]
		this_.freeEdges = j2gort.SliceRemoveAt(this_.freeEdges, last)
		edge.Reinit(sn, v, dn, this_.useEdgeAnnotations)
	} else {
		edge = NewLinkedDirectedGraph_LinkedDiGraphEdge[N, E](sn, v, dn, this_.useEdgeAnnotations)
	}
	sn.outEdges = append(sn.outEdges, edge)
}

func (this_ *LinkedDirectedGraphBase[N, E]) ResetForReuse(useNodeAnnotations bool, useEdgeAnnotations bool) bool {
	var size int32 = this_.nodes.Size()
	if size > int32(4096) {
		return false
	}
	if size > int32(0) {
		if this_.freeNodes == nil {
			this_.freeNodes = j2gort.MakeNonNilSlice[*LinkedDirectedGraph_LinkedDiGraphNode[N, E]](size)
		}
		if this_.freeEdges == nil {
			this_.freeEdges = j2gort.MakeNonNilSlice[*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]](size * 2)
		}
		cur := this_.nodes.IterValues()
		for cur.HasNext() {
			node := cur.Next()
			for i := 0; i < len(node.outEdges); i++ {
				e := node.outEdges[i]
				e.ClearForPool()
				if len(this_.freeEdges) < 8192 {
					this_.freeEdges = append(this_.freeEdges, e)
				}
			}
			node.ClearForPool()
			if len(this_.freeNodes) < 4096 {
				this_.freeNodes = append(this_.freeNodes, node)
			}
		}
		this_.nodes.Clear()
	}
	this_.useNodeAnnotations = useNodeAnnotations
	this_.useEdgeAnnotations = useEdgeAnnotations
	return true
}
`
	if err := os.WriteFile(filepath.Join(graphDir, "linked_directed_graph.go"), []byte(graphSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// Main jscomp package covering 1A (cfgStack), 1B (Phase Purges + ProcessDefines Mode_CHECK vs OPTIMIZE),
	// 1C (InlineFunctions_Reference slab), 2 (NodeTraversal_Builder stack lowering), 3 (Stateless & Stateful NewNodeLambda_*).
	jscompNodeTravSrc := `package jscomp

import (
	"sync/atomic"

	pkg_graph "example.com/w32btest/gen/graph"
	"example.com/w32btest/gen/j2gort"
	pkg_rhino "example.com/w32btest/gen/rhino"
)

type AbstractCompiler interface {
	IsJavaNil_() bool
}

type ScopeCreator interface {
	IsJavaNil_() bool
}

type AbstractScope interface {
	IsJavaNil_() bool
	GetRootNode_0() pkg_rhino.Node
}

type Scope struct {
	root pkg_rhino.Node
}

func (s *Scope) IsJavaNil_() bool             { return s == nil }
func (s *Scope) GetRootNode_0() pkg_rhino.Node { return s.root }

type JSChunk struct{}

type ControlFlowGraph[N comparable] interface {
	pkg_graph.LinkedDirectedGraph[N, int]
}

type ControlFlowGraphBase[N comparable] struct {
	pkg_graph.LinkedDirectedGraphBase[N, int]
}

func NewControlFlowGraph[N comparable](entry N) ControlFlowGraph[N] {
	c := &ControlFlowGraphBase[N]{}
	pkg_graph.InitLinkedDirectedGraph(&c.LinkedDirectedGraphBase, true, true)
	c.CreateNode(entry)
	return c
}

func NodeUtil_IsValidCfgRoot(n pkg_rhino.Node) bool {
	return n != 0 && (n.GetToken() == pkg_rhino.TokROOT || n.IsFunction())
}

type NodeTraversal_Callback interface {
	ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool
	Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type NodeTraversal_ScopedCallback interface {
	NodeTraversal_Callback
	EnterScope(*NodeTraversal)
	ExitScope(*NodeTraversal)
}

type NodeTraversal_AbstractPostOrderCallbackInterface interface {
	IsJavaNil_() bool
	Visit(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type NodeTraversal struct {
	compiler     AbstractCompiler
	callback     NodeTraversal_Callback
	scopeCreator ScopeCreator
	scopeRoot    pkg_rhino.Node
	obey         bool
}

func (t *NodeTraversal) GetScopeRoot() pkg_rhino.Node { return t.scopeRoot }
func (t *NodeTraversal) GetScope() *Scope             { return &Scope{root: t.scopeRoot} }
func (t *NodeTraversal) GetChunk() *JSChunk           { return nil }

func clinit_NodeTraversal() {}

type NodeTraversal_Builder struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCreator                                   ScopeCreator
	idHash_                                        atomic.Int32
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

func (this_ *NodeTraversal_Builder) initFields_NodeTraversal_Builder() {}

var BuilderHeapAllocs atomic.Int64

func NewNodeTraversal_Builder(c AbstractCompiler, cb NodeTraversal_Callback) *NodeTraversal_Builder {
	BuilderHeapAllocs.Add(1)
	this_ := &NodeTraversal_Builder{}
	InitNodeTraversal_Builder(this_, c, cb)
	return this_
}

func InitNodeTraversal_Builder(this_ *NodeTraversal_Builder, c AbstractCompiler, cb NodeTraversal_Callback) {
	this_.initFields_NodeTraversal_Builder()
	this_.compiler = c
	this_.callback = cb
}

func NodeTraversal_S_Builder(c AbstractCompiler, cb NodeTraversal_Callback) *NodeTraversal_Builder {
	clinit_NodeTraversal()
	return NewNodeTraversal_Builder(c, cb)
}

func (this_ *NodeTraversal_Builder) SetCompiler(x AbstractCompiler) *NodeTraversal_Builder {
	this_.compiler = x
	return this_
}

func (this_ *NodeTraversal_Builder) SetCallback_1_NodeTraversal_Callback(x NodeTraversal_Callback) *NodeTraversal_Builder {
	this_.callback = x
	return this_
}

type NodeTraversal__Anon_66 struct {
	outer_ *NodeTraversal_Builder
	cb     NodeTraversal_AbstractPostOrderCallbackInterface
}

func (a *NodeTraversal__Anon_66) ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool {
	return true
}
func (a *NodeTraversal__Anon_66) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) { a.cb.Visit(t, n, p) }

// Escaping setter: stores this_ into heap-allocated NodeTraversal__Anon_66 — MUST NOT be stack-lowered!
func (this_ *NodeTraversal_Builder) SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(x NodeTraversal_AbstractPostOrderCallbackInterface) *NodeTraversal_Builder {
	this_.callback = &NodeTraversal__Anon_66{outer_: this_, cb: x}
	return this_
}

func (this_ *NodeTraversal_Builder) SetScopeCreator(x ScopeCreator) *NodeTraversal_Builder {
	this_.scopeCreator = x
	return this_
}

func (this_ *NodeTraversal_Builder) SetObeyDestructuringAndDefaultValueExecutionOrder(x bool) *NodeTraversal_Builder {
	this_.obeyDestructuringAndDefaultValueExecutionOrder = x
	return this_
}

func runSimpleTraverse(c AbstractCompiler, cb NodeTraversal_Callback, sc ScopeCreator, obey bool, root pkg_rhino.Node) {
	nt := &NodeTraversal{compiler: c, callback: cb, scopeCreator: sc, scopeRoot: root, obey: obey}
	scb, isScoped := cb.(NodeTraversal_ScopedCallback)
	var walk func(n, p pkg_rhino.Node)
	walk = func(n, p pkg_rhino.Node) {
		isScope := NodeUtil_IsValidCfgRoot(n)
		prev := nt.scopeRoot
		if isScope {
			nt.scopeRoot = n
			if isScoped {
				scb.EnterScope(nt)
			}
		}
		if cb.ShouldTraverse_3(nt, n, p) {
			for ch := n.GetFirstChild(); ch != 0; ch = ch.GetNext() {
				walk(ch, n)
			}
			cb.Visit_3(nt, n, p)
		}
		if isScope {
			if isScoped {
				scb.ExitScope(nt)
			}
			nt.scopeRoot = prev
		}
	}
	walk(root, 0)
}

func (this_ *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	runSimpleTraverse(this_.compiler, this_.callback, this_.scopeCreator, this_.obeyDestructuringAndDefaultValueExecutionOrder, root)
}

func (this_ *NodeTraversal_Builder) TraverseRoots(externs pkg_rhino.Node, root pkg_rhino.Node) {
	if externs != 0 {
		runSimpleTraverse(this_.compiler, this_.callback, this_.scopeCreator, this_.obeyDestructuringAndDefaultValueExecutionOrder, externs)
	}
	runSimpleTraverse(this_.compiler, this_.callback, this_.scopeCreator, this_.obeyDestructuringAndDefaultValueExecutionOrder, root)
}

func (this_ *NodeTraversal_Builder) TraverseAtScope(scope AbstractScope) {
	runSimpleTraverse(this_.compiler, this_.callback, this_.scopeCreator, this_.obeyDestructuringAndDefaultValueExecutionOrder, scope.GetRootNode_0())
}

func (this_ *NodeTraversal_Builder) TraverseWithScope(root pkg_rhino.Node, _ AbstractScope) {
	runSimpleTraverse(this_.compiler, this_.callback, this_.scopeCreator, this_.obeyDestructuringAndDefaultValueExecutionOrder, root)
}

// 1A: NodeTraversal_AbstractCfgCallbackBase with cfgStack!
type NodeTraversal_AbstractCfgCallback interface {
	NodeTraversal_ScopedCallback
	EnterScopeWithCfg(*NodeTraversal)
	ExitScopeWithCfg(*NodeTraversal)
}

type NodeTraversal_AbstractCfgCallbackBase struct {
	self_    NodeTraversal_AbstractCfgCallback
	cfgStack j2gort.Deque[ControlFlowGraph[pkg_rhino.Node]]
}

func (this_ *NodeTraversal_AbstractCfgCallbackBase) BindSelf_(s any) {
	this_.self_ = s.(NodeTraversal_AbstractCfgCallback)
}
func (this_ *NodeTraversal_AbstractCfgCallbackBase) GetControlFlowGraph() ControlFlowGraph[pkg_rhino.Node] {
	return this_.cfgStack.Peek()
}
func (this_ *NodeTraversal_AbstractCfgCallbackBase) EnterScope(t *NodeTraversal) {
	if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) {
		cfg := NewControlFlowGraph(t.GetScopeRoot())
		for ch := t.GetScopeRoot().GetFirstChild(); ch != 0; ch = ch.GetNext() {
			cfg.Connect(t.GetScopeRoot(), 1, ch)
		}
		this_.cfgStack.Push(cfg)
	}
	this_.self_.EnterScopeWithCfg(t)
}
func (this_ *NodeTraversal_AbstractCfgCallbackBase) ExitScope(t *NodeTraversal) {
	defer func() {
		if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) {
			this_.cfgStack.Pop()
		}
	}()
	this_.self_.ExitScopeWithCfg(t)
}
func (this_ *NodeTraversal_AbstractCfgCallbackBase) EnterScopeWithCfg(*NodeTraversal) {}
func (this_ *NodeTraversal_AbstractCfgCallbackBase) ExitScopeWithCfg(*NodeTraversal)  {}
func (this_ *NodeTraversal_AbstractCfgCallbackBase) ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool {
	return true
}
func (this_ *NodeTraversal_AbstractCfgCallbackBase) Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) {
}

type CheckUnreachableCode struct {
	NodeTraversal_AbstractCfgCallbackBase
	CheckedScopes int
}

func NewCheckUnreachableCode() *CheckUnreachableCode {
	c := &CheckUnreachableCode{}
	c.BindSelf_(c)
	return c
}
func (this_ *CheckUnreachableCode) EnterScopeWithCfg(t *NodeTraversal) {
	if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) {
		_ = this_.GetControlFlowGraph()
		this_.CheckedScopes++
	}
}

func ExerciseFluentBuilderChains(comp AbstractCompiler, cb NodeTraversal_Callback, sc ScopeCreator, root pkg_rhino.Node, scope AbstractScope) {
	NewNodeTraversal_Builder(comp, cb).SetScopeCreator(sc).SetObeyDestructuringAndDefaultValueExecutionOrder(true).Traverse(root)
	NodeTraversal_S_Builder(comp, cb).SetScopeCreator(sc).TraverseAtScope(scope)
	NewNodeTraversal_Builder(comp, cb).TraverseRoots(0, root)
	NewNodeTraversal_Builder(comp, cb).SetScopeCreator(sc).TraverseWithScope(root, scope)
}

type postOrderFunc struct {
	fn func(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

func (p *postOrderFunc) IsJavaNil_() bool { return p == nil }
func (p *postOrderFunc) Visit(t *NodeTraversal, n, par pkg_rhino.Node) {
	if p.fn != nil {
		p.fn(t, n, par)
	}
}

// Negative case: SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface captures *NodeTraversal_Builder!
func ExerciseEscapingPostOrderBuilder(comp AbstractCompiler, root pkg_rhino.Node) {
	NewNodeTraversal_Builder(comp, nil).SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(&postOrderFunc{}).Traverse(root)
}

type BuilderHelperHolder struct {
	compiler     AbstractCompiler
	callback     NodeTraversal_Callback
	scopeCreator ScopeCreator
}

func NewBuilderHelperHolder(c AbstractCompiler, cb NodeTraversal_Callback) *BuilderHelperHolder {
	return &BuilderHelperHolder{compiler: c, callback: cb}
}

func (this_ *BuilderHelperHolder) priv_BuilderHelperHolder_createTraversalBuilder_99() *NodeTraversal_Builder {
	return NodeTraversal_S_Builder(this_.compiler, this_.callback).SetScopeCreator(this_.scopeCreator).SetObeyDestructuringAndDefaultValueExecutionOrder(true)
}

func (this_ *BuilderHelperHolder) ProcessAtScope(scope AbstractScope) {
	this_.priv_BuilderHelperHolder_createTraversalBuilder_99().TraverseAtScope(scope)
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "node_traversal.go"), []byte(jscompNodeTravSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1B: Phase Purges including ProcessDefines (Mode_CHECK vs Mode_OPTIMIZE) + LIFO verification!
	jscompPhasesSrc := `package jscomp

import (
	"example.com/w32btest/gen/j2gort"
	pkg_rhino "example.com/w32btest/gen/rhino"
)

type Compiler struct{}

func (c *Compiler) IsJavaNil_() bool { return c == nil }

func (this_ *Compiler) ParseInputs() pkg_rhino.Node {
	j2gort.TrimStaleSlabPools(1)
	return 1
}

func (this_ *Compiler) PerformTranspilationAndOptimizations(seg int) {
	_ = seg
	j2gort.TrimStaleSlabPools(2)
}

type TypeInferencePass struct{}

func (this_ *TypeInferencePass) InferAllScopes(root pkg_rhino.Node) {
	_ = root
}

type TypedAstSerializer struct{}

func (this_ *TypedAstSerializer) SerializeRoots(externs, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer j2gort.PurgeAllRuntimePools()
}

type ProcessDefines_Mode uint8

const (
	enum_ProcessDefines_Mode_CHECK    ProcessDefines_Mode = 1
	enum_ProcessDefines_Mode_OPTIMIZE ProcessDefines_Mode = 2
)

func StaticGet_ProcessDefines_Mode_CHECK() ProcessDefines_Mode    { return enum_ProcessDefines_Mode_CHECK }
func StaticGet_ProcessDefines_Mode_OPTIMIZE() ProcessDefines_Mode { return enum_ProcessDefines_Mode_OPTIMIZE }

type FakeNamespace struct {
	slab j2gort.PooledSlab[int]
}

func (ns *FakeNamespace) ReleaseSlabs_() {
	ns.slab.Release()
}

type ProcessDefines struct {
	mode      ProcessDefines_Mode
	namespace *FakeNamespace
}

func NewProcessDefines(mode ProcessDefines_Mode) *ProcessDefines {
	ns := &FakeNamespace{}
	_ = ns.slab.Alloc()
	return &ProcessDefines{mode: mode, namespace: ns}
}

func (this_ *ProcessDefines) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer func() {
		if this_.namespace != nil {
			this_.namespace.ReleaseSlabs_()
			this_.namespace = nil
		}
	}()
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "phase_purges.go"), []byte(jscompPhasesSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1C: InlineFunctions + InlineFunctions_FindCandidatesReferences + InlineFunctions_Reference
	jscompInlineFnsSrc := `package jscomp

import (
	"example.com/w32btest/gen/j2gort"
	pkg_rhino "example.com/w32btest/gen/rhino"
)

type FunctionInjector_InliningMode uint8

const InliningModeDirect FunctionInjector_InliningMode = 1

type FunctionInjector_ReferenceBase struct {
	callNode pkg_rhino.Node
	scope    *Scope
	chunk    *JSChunk
	mode     FunctionInjector_InliningMode
}

func InitFunctionInjector_Reference(this_ *FunctionInjector_ReferenceBase, callNode pkg_rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) {
	this_.callNode = callNode
	this_.scope = scope
	this_.chunk = chunk
	this_.mode = mode
}

type InlineFunctions_Reference struct {
	FunctionInjector_ReferenceBase
	inlined bool
}

func (this_ *InlineFunctions_Reference) initFields_InlineFunctions_Reference() {}

func NewInlineFunctions_Reference(callNode pkg_rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *InlineFunctions_Reference {
	this_ := &InlineFunctions_Reference{}
	InitInlineFunctions_Reference(this_, callNode, scope, chunk, mode)
	return this_
}

func InitInlineFunctions_Reference(this_ *InlineFunctions_Reference, callNode pkg_rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) {
	InitFunctionInjector_Reference(&this_.FunctionInjector_ReferenceBase, callNode, scope, chunk, mode)
	this_.initFields_InlineFunctions_Reference()
}

type InlineFunctions_FunctionState struct {
	references *j2gort.OrderedMap[pkg_rhino.Node, *InlineFunctions_Reference]
}

type InlineFunctions struct {
	compiler   AbstractCompiler
	fns        j2gort.OrderedMap[string, *InlineFunctions_FunctionState]
	InlinedCnt int
}

func NewInlineFunctions(c AbstractCompiler) *InlineFunctions {
	return &InlineFunctions{
		compiler: c,
		fns:      j2gort.AcquireOrderedMapCap[string, *InlineFunctions_FunctionState](16),
	}
}

type InlineFunctions_FindCandidatesReferences struct {
	outer_ *InlineFunctions
}

func NewInlineFunctions_FindCandidatesReferences(outer_ *InlineFunctions) *InlineFunctions_FindCandidatesReferences {
	return &InlineFunctions_FindCandidatesReferences{outer_: outer_}
}

func (this_ *InlineFunctions_FindCandidatesReferences) ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool {
	return true
}

func (this_ *InlineFunctions_FindCandidatesReferences) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		fs := this_.outer_.fns.Get("fn")
		if fs != nil {
			cand := NewInlineFunctions_Reference(n, t.GetScope(), t.GetChunk(), InliningModeDirect)
			if fs.references == nil {
				fs.references = j2gort.AcquireOrderedMapPtrCap[pkg_rhino.Node, *InlineFunctions_Reference](8)
			}
			fs.references.Put(n, cand)
		}
	}
}

func (this_ *InlineFunctions) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	this_.fns.Put("fn", &InlineFunctions_FunctionState{})
	defer func() {
		for cur := this_.fns.IterValues(); cur.HasNext(); {
			fs := cur.Next()
			if fs != nil && fs.references != nil {
				j2gort.ReleaseOrderedMap(fs.references)
				fs.references = nil
			}
		}
		this_.fns.ReleaseToPool()
	}()
	NewNodeTraversal_Builder(this_.compiler, NewInlineFunctions_FindCandidatesReferences(this_)).Traverse(root)
	fs := this_.fns.Get("fn")
	if fs != nil && fs.references != nil {
		this_.InlinedCnt += int(fs.references.Size())
	}
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "inline_functions.go"), []byte(jscompInlineFnsSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// Task 3: Stateless NewNodeLambda_* + FunctionInjector DoesFunctionMeetMinimumRequirements & referencesEval!
	jscompNodeLambdaSrc := `package jscomp

import (
	"example.com/w32btest/gen/j2gort"
	pkg_rhino "example.com/w32btest/gen/rhino"
)

func NodeUtil_Has(node pkg_rhino.Node, pred j2gort.Ext_com_google_common_base_Predicate, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = node
	_ = pred
	_ = traverseChildrenPred
	if j2gort.PredTestNode(pred, uint32(node)) {
		return true
	}
	if !(j2gort.PredTestNode(traverseChildrenPred, uint32(node))) {
		return false
	}
	for c := node.GetFirstChild(); c != 0; c = c.GetNext() {
		if NodeUtil_Has(c, pred, traverseChildrenPred) {
			return true
		}
	}
	return false
}

var evalPredSingleton j2gort.Ext_com_google_common_base_Predicate

func init() {
	evalPredSingleton = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(raw uint32) bool {
		n := pkg_rhino.Node(raw)
		return n.IsName() && n.GetString() == "eval"
	})
}

// Stateless NewNodeLambda_* calls inside methods/functions -> hoisted to package singletons!
func StatelessReturnPred() j2gort.Ext_com_google_common_base_Predicate {
	return j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_1 uint32) bool {
		input := pkg_rhino.Node(rawNode_1)
		return input != 0 && input.IsReturn()
	})
}

func StatelessJavaNamePred() j2gort.Ext_java_util_function_Predicate {
	return j2gort.NewNodeLambda_Ext_java_util_function_Predicate(func(rawNode_2 uint32) bool {
		n := pkg_rhino.Node(rawNode_2)
		return n.IsName()
	})
}

// Negative case: captures outer parameter wantName and returned directly -> MUST NOT be hoisted!
func StatefulReturnedPred(wantName string) j2gort.Ext_com_google_common_base_Predicate {
	return j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_3 uint32) bool {
		n := pkg_rhino.Node(rawNode_3)
		return n.IsName() && n.GetString() == wantName
	})
}

type FunctionInjector struct{}

func (this_ *FunctionInjector) DoesFunctionMeetMinimumRequirements(fnName string, fnNode pkg_rhino.Node) bool {
	_ = this_
	var block pkg_rhino.Node = fnNode.GetFirstChild().GetNext()
	var fnRecursionName string = fnNode.GetFirstChild().GetString()
	var blocksInjection j2gort.Ext_com_google_common_base_Predicate = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_1 uint32) bool {
		n := pkg_rhino.Node(rawNode_1)
		if n.IsName() {
			return (j2gort.CheckNotNullString(n.GetString()) == "eval") ||
				(fnName != "" && j2gort.CheckNotNullString(n.GetString()) == fnName) ||
				(fnRecursionName != "" && j2gort.CheckNotNullString(n.GetString()) == fnRecursionName)
		} else if n.IsSuper() || n.RawToken_() == uint8(pkg_rhino.TokSUPER) {
			return true
		}
		return false
	})
	_ = blocksInjection
	return !(NodeUtil_Has(block, blocksInjection, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue()))
}

func (this_ *FunctionInjector) ReferencesEval(fn pkg_rhino.Node) bool {
	_ = this_
	var result bool = NodeUtil_Has(fn, evalPredSingleton, j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_5 uint32) bool {
		n := pkg_rhino.Node(rawNode_5)
		return (!n.IsFunction()) || (j2gort.CheckNotNullNode(n) == fn)
	}))
	return result
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "function_injector.go"), []byte(jscompNodeLambdaSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// Variant jscomp_checksonly: ProcessDefines with `checksOnly bool` field!
	checksBoolSrc := `package jscomp

import (
	"example.com/w32btest/gen/j2gort"
	pkg_rhino "example.com/w32btest/gen/rhino"
)

type ProcessDefines struct {
	checksOnly bool
}

func (this_ *ProcessDefines) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	j2gort.TrimStaleSlabPools(9)
}
`
	if err := os.WriteFile(filepath.Join(jscompChecksBoolDir, "process_defines.go"), []byte(checksBoolSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// Negative case jscomp_esc: NodeTraversal_AbstractCfgCallbackBase where cfgStack escapes!
	escCfgSrc := `package jscomp

import (
	"example.com/w32btest/gen/j2gort"
	pkg_rhino "example.com/w32btest/gen/rhino"
)

type NodeTraversal struct{}

type NodeTraversal_AbstractCfgCallbackBase struct {
	cfgStack j2gort.Deque[any]
}

func (this_ *NodeTraversal_AbstractCfgCallbackBase) LeakStack() *j2gort.Deque[any] {
	return &this_.cfgStack
}

func (this_ *NodeTraversal_AbstractCfgCallbackBase) ExitScope(_ *NodeTraversal) {
	_ = pkg_rhino.Node(0)
	this_.cfgStack.Pop()
}
`
	if err := os.WriteFile(filepath.Join(jscompEscDir, "escaping_cfg.go"), []byte(escCfgSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()

	// Run 1 + Run 2 idempotency on graph, jscomp, jscomp_checksonly, jscomp_esc!
	for _, target := range []struct {
		name      string
		dir       string
		wantMinS1 int
	}{
		{"graph", graphDir, 4},
		{"jscomp", jscompDir, 10},
		{"jscomp_checksonly", jscompChecksBoolDir, 1},
		{"jscomp_esc", jscompEscDir, 0},
	} {
		s1, err := passes.RunOnPackageDir(target.dir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 1 on %s failed: %v", target.name, err)
		}
		if s1.TotalSites < target.wantMinS1 {
			t.Fatalf("Run 1 on %s: TotalSites = %d, want >= %d (byPass=%+v)", target.name, s1.TotalSites, target.wantMinS1, s1.ByPass)
		}
		s2, err := passes.RunOnPackageDir(target.dir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 2 (idempotency) on %s failed: %v", target.name, err)
		}
		if s2.TotalSites != 0 || s2.FilesModified != 0 {
			t.Fatalf("Run 2 on %s NOT IDEMPOTENT: TotalSites=%d FilesModified=%d byPass=%+v",
				target.name, s2.TotalSites, s2.FilesModified, s2.ByPass)
		}
	}

	// Verify AST shapes in lowered files.
	gotGraphBytes, err := os.ReadFile(filepath.Join(graphDir, "linked_directed_graph.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotGraph := string(gotGraphBytes)
	for _, want := range []string{
		"this_.slab_LinkedDirectedGraph_LinkedDiGraphEdge.Reset()",
		"this_.slab_LinkedDirectedGraph_LinkedDiGraphNode.Reset()",
		"this_.nodes.ReleaseToPool()",
		"func (this_ *LinkedDirectedGraphBase[N, E]) ReleaseSlabs_()",
	} {
		if !strings.Contains(gotGraph, want) {
			t.Fatalf("missing %q in lowered graph/linked_directed_graph.go:\n%s", want, gotGraph)
		}
	}

	gotNTBytes, err := os.ReadFile(filepath.Join(jscompDir, "node_traversal.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotNT := string(gotNTBytes)
	for _, want := range []string{
		"poppedCfg_ := this_.cfgStack.Pop()",
		"r_.ReleaseSlabs_()",
		"var ntb_ NodeTraversal_Builder",
		"InitNodeTraversal_Builder(&ntb_, comp, cb)",
		"(&ntb_).SetScopeCreator(sc)",
		"(&ntb_).Traverse(root)",
		"(&ntb_1_).TraverseAtScope(scope)",
		"(&ntb_2_).TraverseRoots(0, root)",
		"(&ntb_3_).TraverseWithScope(root, scope)",
		"(&ntb_).TraverseAtScope(scope)",
		"InitNodeTraversal_Builder(&ntb_, this_.compiler, this_.callback)",
	} {
		if !strings.Contains(gotNT, want) {
			t.Fatalf("missing %q in lowered jscomp/node_traversal.go:\n%s", want, gotNT)
		}
	}
	// Escaping SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface MUST remain heap-constructed!
	if !strings.Contains(gotNT, "NewNodeTraversal_Builder(comp, nil).SetCallback_1_NodeTraversal_AbstractPostOrderCallbackInterface(") {
		t.Fatalf("expected escaping AbstractPostOrderCallbackInterface chain to remain untouched:\n%s", gotNT)
	}

	gotPhasesBytes, err := os.ReadFile(filepath.Join(jscompDir, "phase_purges.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotPhases := string(gotPhasesBytes)
	if strings.Contains(gotPhases, "TrimStaleSlabPools") {
		t.Fatalf("expected TrimStaleSlabPools stripped in phase_purges.go:\n%s", gotPhases)
	}
	if got := strings.Count(gotPhases, "PurgeAllRuntimePools()"); got != 3 {
		t.Fatalf("expected exactly 3 PurgeAllRuntimePools() calls in phase_purges.go (PerformTranspilationAndOptimizations defer + InferAllScopes entry + SerializeRoots entry), got %d:\n%s", got, gotPhases)
	}

	gotChecksBoolBytes, err := os.ReadFile(filepath.Join(jscompChecksBoolDir, "process_defines.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotChecksBool := string(gotChecksBoolBytes)
	if strings.Contains(gotChecksBool, "PurgeAllRuntimePools") || strings.Contains(gotChecksBool, "TrimStaleSlabPools") {
		t.Fatalf("expected TrimStaleSlabPools & PurgeAllRuntimePools stripped in jscomp_checksonly:\n%s", gotChecksBool)
	}

	gotEscBytes, err := os.ReadFile(filepath.Join(jscompEscDir, "escaping_cfg.go"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(gotEscBytes), "ReleaseSlabs_") {
		t.Fatalf("expected escaping cfgStack NOT to be rewritten:\n%s", string(gotEscBytes))
	}

	gotIFBytes, err := os.ReadFile(filepath.Join(jscompDir, "inline_functions.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotIF := string(gotIFBytes)
	for _, want := range []string{
		"slab_InlineFunctions_Reference j2gort.PooledSlab[InlineFunctions_Reference]",
		"allocPooled_InlineFunctions_Reference(&this_.outer_.slab_InlineFunctions_Reference,",
		"func (this_ *InlineFunctions) ReleaseSlabs_()",
		"defer this_.ReleaseSlabs_()",
	} {
		if !strings.Contains(gotIF, want) {
			t.Fatalf("missing %q in lowered inline_functions.go:\n%s", want, gotIF)
		}
	}

	gotFIBytes, err := os.ReadFile(filepath.Join(jscompDir, "function_injector.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotFI := string(gotFIBytes)
	for _, want := range []string{
		"var lambdaSingleton_peep_function_injector_1 j2gort.Ext_com_google_common_base_Predicate",
		"var lambdaSingleton_peep_function_injector_2 j2gort.Ext_java_util_function_Predicate",
		"nodeUtilHasPredAlwaysTrueFn_(",
		"nodeUtilHasTravFn_(",
		"func nodeUtilHasPredAlwaysTrueFn_(",
		"func nodeUtilHasTravFn_(",
	} {
		if !strings.Contains(gotFI, want) {
			t.Fatalf("missing %q in lowered function_injector.go:\n%s", want, gotFI)
		}
	}
	// StatefulReturnedPred must NOT be hoisted to a singleton!
	if !strings.Contains(gotFI, "func StatefulReturnedPred(wantName string)") ||
		strings.Contains(gotFI, "lambdaSingleton_peep_function_injector_3") {
		t.Fatalf("StatefulReturnedPred must NOT be hoisted to a package singleton:\n%s", gotFI)
	}

	// Verify with `go build -gcflags=-m` that `ntb_` in ExerciseFluentBuilderChains and `BuilderHelperHolder.ProcessAtScope`
	// NEVER escapes to the heap (`moved to heap: ntb_` must NOT appear)!
	escCmd := exec.Command("go", "build", "-gcflags=-m", "./gen/jscomp")
	escCmd.Dir = tmpDir
	escOutBytes, err := escCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build -gcflags=-m ./gen/jscomp failed: %v\n%s", err, string(escOutBytes))
	}
	escOut := string(escOutBytes)
	if strings.Contains(escOut, "moved to heap: ntb_") {
		t.Fatalf("ntb_ escaped to the heap under -gcflags=-m:\n%s", escOut)
	}
	if !strings.Contains(escOut, "func literal does not escape") {
		t.Fatalf("expected unboxed NodeUtil_Has func literals not to escape under -gcflags=-m:\n%s", escOut)
	}

	// Runtime behavioral & zero-alloc verification test inside tmpDir!
	rtTestSrc := `package w32b_rt_test

import (
	"reflect"
	"testing"

	"example.com/w32btest/gen/graph"
	"example.com/w32btest/gen/j2gort"
	"example.com/w32btest/gen/jscomp"
	"example.com/w32btest/gen/rhino"
)

func TestRuntimeSemanticsAndZeroAlloc(t *testing.T) {
	// 1A: LinkedDirectedGraph.ResetForReuse resets both PooledSlabs (including >8192 edges!)
	j2gort.ResetWave32BMetrics()
	g := graph.NewLinkedDirectedGraph[int, int](true, true)
	g.Connect(1, 10, 2)
	g.Connect(2, 20, 3)
	if !g.ResetForReuse(false, false) {
		t.Fatal("expected ResetForReuse == true")
	}
	if got := j2gort.SlabResetCalls.Load(); got != 2 {
		t.Fatalf("after ResetForReuse: SlabResetCalls = %d, want 2 (node + edge slabs)", got)
	}

	// 1A: CheckUnreachableCode via NodeTraversal_AbstractCfgCallbackBase releases per-scope CFG slabs on ExitScope!
	rhino.ResetNodes()
	fn1 := rhino.MakeNode(rhino.TokFUNCTION, "f1", rhino.MakeNode(rhino.TokNAME, "f1"), rhino.MakeNode(rhino.TokBLOCK, ""))
	fn2 := rhino.MakeNode(rhino.TokFUNCTION, "f2", rhino.MakeNode(rhino.TokNAME, "f2"), rhino.MakeNode(rhino.TokBLOCK, ""))
	root := rhino.MakeNode(rhino.TokROOT, "", fn1, fn2)
	j2gort.ResetWave32BMetrics()
	chk := jscomp.NewCheckUnreachableCode()
	comp := &jscomp.Compiler{}
	jscomp.NewNodeTraversal_Builder(comp, chk).Traverse(root)
	if chk.CheckedScopes != 3 {
		t.Fatalf("CheckedScopes = %d, want 3 (root + 2 functions)", chk.CheckedScopes)
	}
	// Each of the 3 scopes has 2 PooledSlabs (nodes + edges) released on ExitScope => 6 SlabReleaseCalls!
	if got := j2gort.SlabReleaseCalls.Load(); got != 6 {
		t.Fatalf("after CheckUnreachableCode: SlabReleaseCalls = %d, want 6", got)
	}

	// 1B: Phase Purges (3 exact sites: InferAllScopes entry, SerializeRoots entry, PerformTranspilationAndOptimizations defer; ProcessDefines & ParseInputs unpurged!)
	j2gort.ResetWave32BMetrics()
	_ = comp.ParseInputs()
	if got := j2gort.PurgeCalls.Load(); got != 0 {
		t.Fatalf("after ParseInputs: PurgeCalls = %d, want 0", got)
	}
	tip := &jscomp.TypeInferencePass{}
	tip.InferAllScopes(root)
	if got := j2gort.PurgeCalls.Load(); got != 1 {
		t.Fatalf("after InferAllScopes: PurgeCalls = %d, want 1 (entry only)", got)
	}

	j2gort.ResetWave32BMetrics()
	pdCheck := jscomp.NewProcessDefines(jscomp.StaticGet_ProcessDefines_Mode_CHECK())
	pdCheck.Process_2(0, root)
	wantLog := []string{"slab_release"}
	if !reflect.DeepEqual(j2gort.EventLog, wantLog) {
		t.Fatalf("ProcessDefines(CHECK) EventLog = %v, want %v (must NOT purge in ProcessDefines)", j2gort.EventLog, wantLog)
	}

	j2gort.ResetWave32BMetrics()
	pdOpt := jscomp.NewProcessDefines(jscomp.StaticGet_ProcessDefines_Mode_OPTIMIZE())
	pdOpt.Process_2(0, root)
	wantOptLog := []string{"slab_release"}
	if !reflect.DeepEqual(j2gort.EventLog, wantOptLog) {
		t.Fatalf("ProcessDefines(OPTIMIZE) EventLog = %v, want %v (must NOT purge in OPT!)", j2gort.EventLog, wantOptLog)
	}

	j2gort.ResetWave32BMetrics()
	tas := &jscomp.TypedAstSerializer{}
	tas.SerializeRoots(0, root)
	comp.PerformTranspilationAndOptimizations(0)
	if got := j2gort.PurgeCalls.Load(); got != 2 {
		t.Fatalf("after SerializeRoots + PerformTranspilationAndOptimizations: PurgeCalls = %d, want 2", got)
	}

	// 1C: InlineFunctions_Reference PooledSlab allocation and deferred ReleaseSlabs_()
	j2gort.ResetWave32BMetrics()
	inl := jscomp.NewInlineFunctions(comp)
	inl.Process_2(0, root)
	if inl.InlinedCnt != 2 {
		t.Fatalf("InlineFunctions.InlinedCnt = %d, want 2", inl.InlinedCnt)
	}
	if got := j2gort.SlabAllocCalls.Load(); got != 2 {
		t.Fatalf("InlineFunctions SlabAllocCalls = %d, want 2", got)
	}
	if got := j2gort.SlabReleaseCalls.Load(); got != 1 {
		t.Fatalf("InlineFunctions SlabReleaseCalls = %d, want 1", got)
	}

	// 2: Fluent NodeTraversal_Builder chains + helper method chain work identically with 0 builder heap allocs!
	leaf := rhino.MakeNode(rhino.TokNAME, "x")
	holderHolderRun(comp, leaf)

	// 3: Stateless NewNodeLambda_* singletons allocate 0 times per call!
	j2gort.ResetWave32BMetrics()
	p1 := jscomp.StatelessReturnPred()
	p2 := jscomp.StatelessReturnPred()
	jp1 := jscomp.StatelessJavaNamePred()
	jp2 := jscomp.StatelessJavaNamePred()
	if p1 != p2 || jp1 != jp2 {
		t.Fatal("expected identical singleton instances across calls")
	}
	if got := j2gort.NodeLambdaCalls.Load(); got != 0 {
		t.Fatalf("Stateless NewNodeLambda_* allocated %d times at call time, want 0", got)
	}
	retNode := rhino.MakeNode(rhino.TokRETURN, "")
	if !j2gort.PredTestNode(p1, uint32(retNode)) || j2gort.PredTestNode(p1, uint32(leaf)) {
		t.Fatal("StatelessReturnPred semantics mismatch")
	}
	spA := jscomp.StatefulReturnedPred("alpha")
	spB := jscomp.StatefulReturnedPred("beta")
	nameAlpha := rhino.MakeNode(rhino.TokNAME, "alpha")
	if !j2gort.PredTestNode(spA, uint32(nameAlpha)) || j2gort.PredTestNode(spB, uint32(nameAlpha)) {
		t.Fatal("StatefulReturnedPred capture semantics mismatch")
	}

	// 3: FunctionInjector DoesFunctionMeetMinimumRequirements + ReferencesEval allocate ZERO heap objects!
	fi := &jscomp.FunctionInjector{}
	cleanFn := rhino.MakeNode(rhino.TokFUNCTION, "",
		rhino.MakeNode(rhino.TokNAME, "myFn"),
		rhino.MakeNode(rhino.TokBLOCK, "", rhino.MakeNode(rhino.TokNAME, "okVar")),
	)
	recFn := rhino.MakeNode(rhino.TokFUNCTION, "",
		rhino.MakeNode(rhino.TokNAME, "recFn"),
		rhino.MakeNode(rhino.TokBLOCK, "", rhino.MakeNode(rhino.TokNAME, "recFn")),
	)
	evalFn := rhino.MakeNode(rhino.TokFUNCTION, "",
		rhino.MakeNode(rhino.TokNAME, "outerFn"),
		rhino.MakeNode(rhino.TokBLOCK, "", rhino.MakeNode(rhino.TokNAME, "eval")),
	)
	nestedEvalFn := rhino.MakeNode(rhino.TokFUNCTION, "",
		rhino.MakeNode(rhino.TokNAME, "outerFn2"),
		rhino.MakeNode(rhino.TokBLOCK, "", evalFn),
	)
	if !fi.DoesFunctionMeetMinimumRequirements("myFn", cleanFn) {
		t.Fatal("expected cleanFn to meet minimum requirements")
	}
	if fi.DoesFunctionMeetMinimumRequirements("recFn", recFn) {
		t.Fatal("expected recFn to fail minimum requirements due to recursion")
	}
	if !fi.ReferencesEval(evalFn) {
		t.Fatal("expected evalFn to reference eval")
	}
	if fi.ReferencesEval(nestedEvalFn) {
		t.Fatal("expected nestedEvalFn NOT to reference eval directly (inner function boundary)")
	}

	allocs := testing.AllocsPerRun(200, func() {
		_ = fi.DoesFunctionMeetMinimumRequirements("myFn", cleanFn)
		_ = fi.ReferencesEval(nestedEvalFn)
		_ = jscomp.StatelessReturnPred()
		_ = jscomp.StatelessJavaNamePred()
	})
	if allocs != 0 {
		t.Fatalf("expected 0 heap allocs per run for FunctionInjector + StatelessNodeLambda, got %v", allocs)
	}
}

type scopeWrapper struct{ n rhino.Node }

func (s *scopeWrapper) IsJavaNil_() bool          { return s == nil }
func (s *scopeWrapper) GetRootNode_0() rhino.Node { return s.n }

type noOpCb struct{}

func (noOpCb) ShouldTraverse_3(*jscomp.NodeTraversal, rhino.Node, rhino.Node) bool { return true }
func (noOpCb) Visit_3(*jscomp.NodeTraversal, rhino.Node, rhino.Node)               {}

func holderHolderRun(comp jscomp.AbstractCompiler, leaf rhino.Node) {
	sw := &scopeWrapper{n: leaf}
	cb := &noOpCb{}
	holder := jscomp.NewBuilderHelperHolder(comp, cb)
	jscomp.BuilderHeapAllocs.Store(0)
	jscomp.ExerciseFluentBuilderChains(comp, cb, nil, leaf, sw)
	holder.ProcessAtScope(sw)
	if got := jscomp.BuilderHeapAllocs.Load(); got != 0 {
		panic("expected 0 NewNodeTraversal_Builder calls in stack-lowered chains")
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w32b_rt_test.go"), []byte(rtTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	rtCmd := exec.Command("go", "test", "-v", "-run", "TestRuntimeSemanticsAndZeroAlloc", ".")
	rtCmd.Dir = tmpDir
	if out, err := rtCmd.CombinedOutput(); err != nil {
		t.Fatalf("runtime test failed: %v\n%s", err, string(out))
	}
}

func TestWave32B_RealGenDryRunCheck(t *testing.T) {
	repoRoot := filepath.Join("..", "..", "..")
	realJscompDir := filepath.Join(repoRoot, "transpiled", "gen", "jscomp")
	realGraphDir := filepath.Join(repoRoot, "transpiled", "gen", "graph")
	if _, err := os.Stat(filepath.Join(realJscompDir, "inline_functions.go")); err != nil {
		t.Skip("real gen/jscomp not present")
	}

	for _, dir := range []string{realGraphDir, realJscompDir} {
		dryStats, err := passes.RunOnPackageDir(dir, passes.AllPasses(), passes.Options{DryRun: true})
		if err != nil {
			t.Fatalf("DryRun AllPasses on %s failed: %v", dir, err)
		}
		t.Logf("DryRun AllPasses on %s: TotalSites=%d ByPass=%+v", filepath.Base(dir), dryStats.TotalSites, dryStats.ByPass)
	}
}

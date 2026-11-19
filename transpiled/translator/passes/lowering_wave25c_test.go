package passes_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave25C_EndToEndScopeAndSlabLifecycleAndRace(t *testing.T) {
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

	mustWrite("go.mod", "module example.com/w25ctest\n\ngo 1.24\n")

	// 1. Minimal shims/nodestore
	mustWrite("shims/nodestore/nodestore.go", `package nodestore

import "sync"

type Node uint32

type NodeRec struct {
	Token   uint8
	Payload string
	First   Node
	Next    Node
	Parent  Node
}

var (
	mu   sync.RWMutex
	recs = make([]NodeRec, 1, 1024)
)

func AllocRaw(rec NodeRec) Node {
	mu.Lock()
	defer mu.Unlock()
	id := Node(len(recs))
	recs = append(recs, rec)
	return id
}

func GetRec(n Node) NodeRec {
	mu.RLock()
	defer mu.RUnlock()
	return recs[n]
}

func MutateRec(n Node, fn func(*NodeRec)) {
	mu.Lock()
	defer mu.Unlock()
	fn(&recs[n])
}
`)

	// 2. Minimal gen/j2gort with PooledSlab, OrderedMap (Acquire + ReleaseToPool), OrderedSet, Deque, Supplier, etc.
	mustWrite("gen/j2gort/j2gort.go", `package j2gort

import "sync/atomic"

const NullString = "\xff\x00j2go_null"

type NPE struct{}
func (NPE) Error() string { return "NullPointerException" }

func ThrowNPE() { panic(NPE{}) }
func CheckNotNull[T any](v T) T {
	if any(v) == nil {
		ThrowNPE()
	}
	return v
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

type Ext_java_lang_Iterable interface {
	IterItems_() []any
}
type Ext_java_util_Set interface {
	Ext_java_lang_Iterable
	AddStr_(string)
	Clear()
	IsEmpty() bool
}

type StringSet struct {
	items []string
}

func NewStringSet() *StringSet { return &StringSet{} }
func (s *StringSet) IterItems_() []any {
	out := make([]any, len(s.items))
	for i, v := range s.items {
		out[i] = v
	}
	return out
}
func (s *StringSet) AddStr_(v string) { s.items = append(s.items, v) }
func (s *StringSet) Clear()           { s.items = s.items[:0] }
func (s *StringSet) IsEmpty() bool    { return len(s.items) == 0 }
func (s *StringSet) Count() int       { return len(s.items) }

func SetAddStr(dst Ext_java_util_Set, v string) {
	dst.AddStr_(v)
}

type IterSlice struct{ Items []any }
func (s *IterSlice) IterItems_() []any { return s.Items }

type JIterCursor struct {
	items []any
	idx   int
}

func NewIterCursor(it Ext_java_lang_Iterable) JIterCursor {
	if it == nil {
		ThrowNPE()
	}
	return JIterCursor{items: it.IterItems_()}
}
func (c *JIterCursor) HasNext() bool { return c.idx < len(c.items) }
func (c *JIterCursor) Next() any {
	v := c.items[c.idx]
	c.idx++
	return v
}

var (
	SlabAllocCount         atomic.Int64
	SlabReleaseCount       atomic.Int64
	OrderedMapReleaseCount atomic.Int64
)

type OrderedMap[K comparable, V any] struct {
	keys []K
	vals []V
	idx  map[K]int
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{idx: make(map[K]int)}
}
func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return NewOrderedMap[K, V]()
}
func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	m := NewOrderedMap[K, V]()
	return &m
}
func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	return NewOrderedMapPtr[K, V]()
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
	if m == nil || m.idx == nil {
		var z V
		return z
	}
	if i, ok := m.idx[k]; ok {
		return m.vals[i]
	}
	var z V
	return z
}
func (m *OrderedMap[K, V]) Size() int32   { if m == nil { return 0 }; return int32(len(m.keys)) }
func (m *OrderedMap[K, V]) IsEmpty() bool { return m == nil || len(m.keys) == 0 }
func (m *OrderedMap[K, V]) Clear() {
	if m == nil {
		return
	}
	clear(m.keys)
	clear(m.vals)
	m.keys = m.keys[:0]
	m.vals = m.vals[:0]
	clear(m.idx)
}
func (m *OrderedMap[K, V]) ReleaseToPool() {
	OrderedMapReleaseCount.Add(1)
	m.Clear()
}
func (m *OrderedMap[K, V]) LiveValues() Ext_java_lang_Iterable {
	out := make([]any, len(m.vals))
	for i, v := range m.vals {
		out[i] = v
	}
	return &IterSlice{Items: out}
}

type OrderedMapValCursor[K comparable, V any] struct {
	vals []V
	i    int
}
func (m *OrderedMap[K, V]) IterValues() OrderedMapValCursor[K, V] {
	if m == nil {
		return OrderedMapValCursor[K, V]{}
	}
	return OrderedMapValCursor[K, V]{vals: m.vals}
}
func (c *OrderedMapValCursor[K, V]) HasNext() bool { return c.i < len(c.vals) }
func (c *OrderedMapValCursor[K, V]) Next() V {
	v := c.vals[c.i]
	c.i++
	return v
}

type Deque[T any] struct{ s []T }
func (d *Deque[T]) IsEmpty() bool { return len(d.s) == 0 }
func (d *Deque[T]) Push(v T)      { d.s = append(d.s, v) }
func (d *Deque[T]) Pop() T        { v := d.s[len(d.s)-1]; d.s = d.s[:len(d.s)-1]; return v }
func (d *Deque[T]) Peek() T       { return d.s[len(d.s)-1] }
func (d *Deque[T]) Clear()        { clear(d.s); d.s = d.s[:0] }

type Slab[T any] struct{ items []*T }
func (s *Slab[T]) Alloc() *T {
	SlabAllocCount.Add(1)
	v := new(T)
	s.items = append(s.items, v)
	return v
}

type PooledSlab[T any] struct{ items []*T }
func (s *PooledSlab[T]) Alloc() *T {
	SlabAllocCount.Add(1)
	v := new(T)
	s.items = append(s.items, v)
	return v
}
func (s *PooledSlab[T]) Reset() {
	for _, it := range s.items {
		var z T
		*it = z
	}
	s.items = s.items[:0]
}
func (s *PooledSlab[T]) Release() {
	SlabReleaseCount.Add(1)
	s.Reset()
}
func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	for _, it := range s.items {
		fn(it)
	}
	s.Release()
}
`)

	// 3. Minimal gen/rhino
	mustWrite("gen/rhino/rhino.go", `package rhino

import "example.com/w25ctest/shims/nodestore"

type Token uint8

const (
	TokSCRIPT   Token = 1
	TokFUNCTION Token = 2
	TokBLOCK    Token = 3
	TokNAME     Token = 4
)

type Node nodestore.Node

func (n Node) IsJavaNil_() bool    { return n == 0 }
func (n Node) GetToken() Token     { return Token(nodestore.GetRec(nodestore.Node(n)).Token + 1) }
func (n Node) IsScript() bool      { return n != 0 && n.GetToken() == TokSCRIPT }
func (n Node) IsFunction() bool    { return n != 0 && n.GetToken() == TokFUNCTION }
func (n Node) IsBlock() bool       { return n != 0 && n.GetToken() == TokBLOCK }
func (n Node) IsName() bool        { return n != 0 && n.GetToken() == TokNAME }
func (n Node) GetFirstChild() Node { return Node(nodestore.GetRec(nodestore.Node(n)).First) }
func (n Node) GetNext() Node       { return Node(nodestore.GetRec(nodestore.Node(n)).Next) }
func (n Node) GetParent() Node     { return Node(nodestore.GetRec(nodestore.Node(n)).Parent) }
func (n Node) GetString() string   { return nodestore.GetRec(nodestore.Node(n)).Payload }
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
`)

	// 4. Minimal gen/jscomp/scope.go (with AddVarNamesTo, IsGoogModuleExports, GetImplicitGoogNamespaceStrength, MemoizedScopeCreator)
	mustWrite("gen/jscomp/scope.go", `package jscomp

import (
	"example.com/w25ctest/gen/j2gort"
	pkg_rhino "example.com/w25ctest/gen/rhino"
)

type AbstractCompiler interface {
	Implements_jscomp_AbstractCompiler()
	IsJavaNil_() bool
}

type Compiler struct{}
func (*Compiler) Implements_jscomp_AbstractCompiler() {}
func (c *Compiler) IsJavaNil_() bool                  { return c == nil }

type AbstractVar interface {
	Implements_jscomp_AbstractVar()
	IsJavaNil_() bool
	GetName() string
	GetScope() AbstractScope
}

type Var struct {
	name  string
	node  pkg_rhino.Node
	scope *Scope
}

func (*Var) Implements_jscomp_AbstractVar()             {}
func (v *Var) IsJavaNil_() bool                         { return v == nil }
func (v *Var) GetName() string                          { return v.name }
func (v *Var) GetNode() pkg_rhino.Node                  { return v.node }
func (v *Var) GetNameNode() pkg_rhino.Node              { return v.node }
func (v *Var) GetScopeRoot() pkg_rhino.Node             { return v.scope.rootNode }
func (v *Var) IsLocal() bool                            { return v.scope != nil && v.scope.parent != nil }
func (v *Var) IsGoogModuleExports() bool                { return v.name == "exports" }
func (v *Var) GetImplicitGoogNamespaceStrength() int32  { return 1 }
func (v *Var) GetScope() AbstractScope                  { return v.scope }

func NewVar(name string, node pkg_rhino.Node, scope *Scope) *Var {
	v := &Var{}
	InitVar(v, name, node, scope)
	return v
}
func InitVar(v *Var, name string, node pkg_rhino.Node, scope *Scope) {
	v.name = name
	v.node = node
	v.scope = scope
}

type AbstractScope interface {
	Implements_jscomp_AbstractScope()
	IsJavaNil_() bool
	GetRootNode() pkg_rhino.Node
	Untyped() *Scope
}

type AbstractScopeBase struct {
	rootNode pkg_rhino.Node
}
func (*AbstractScopeBase) Implements_jscomp_AbstractScope() {}
func (s *AbstractScopeBase) GetRootNode() pkg_rhino.Node    { return s.rootNode }
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
}

func (s *Scope) IsJavaNil_() bool                        { return s == nil }
func (s *Scope) Untyped() *Scope                         { return s }
func (s *Scope) GetParent() AbstractScope                { if s.parent == nil { return nil }; return s.parent }
func (s *Scope) GetClosestContainerScope() *Scope        { return s }
func (s *Scope) GetClosestHoistScope() *Scope            { return s }
func (s *Scope) HasSameContainerScope(o AbstractScope) bool {
	return o != nil && s == o.Untyped()
}
func (s *Scope) IsGlobal() bool { return s.parent == nil }

func (this_ *Scope) AddVarNamesTo(dst j2gort.Ext_java_util_Set) {
	for _, v := range []*Var{this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5} {
		if v != nil {
			j2gort.SetAddStr(dst, v.name)
		}
	}
	if this_.vars != nil {
		cur := this_.vars.IterValues()
		for cur.HasNext() {
			j2gort.SetAddStr(dst, cur.Next().name)
		}
	}
}

func NewScope(parent *Scope, root pkg_rhino.Node) *Scope {
	s := &Scope{}
	InitScope(s, parent, root)
	return s
}
func InitScope(s *Scope, parent *Scope, root pkg_rhino.Node) {
	s.parent = parent
	s.rootNode = root
}

func (this_ *Scope) Declare(name string, n pkg_rhino.Node, v *Var) *Var {
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
		this_.vars = j2gort.NewOrderedMapPtr[string, *Var]()
		for _, old := range []*Var{this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5} {
			this_.vars.Put(old.name, old)
		}
		this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5 = nil, nil, nil, nil, nil, nil
		this_.vars.Put(name, v)
	}
	return v
}

func (this_ *Scope) GetVar(name string) any {
	if this_.vars != nil {
		return j2gort.PtrToAny(this_.vars.Get(name))
	}
	for _, v := range []*Var{this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5} {
		if v != nil && v.name == name {
			return v
		}
	}
	if this_.parent != nil {
		return this_.parent.GetVar(name)
	}
	return nil
}
func (this_ *Scope) GetSlot(name string) any    { return this_.GetVar(name) }
func (this_ *Scope) GetOwnSlot(name string) any { return this_.GetVar(name) }

func (this_ *Scope) GetVarIterable() j2gort.Ext_java_lang_Iterable {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	if this_.vars != nil {
		return this_.vars.LiveValues()
	}
	var items []any
	for _, v := range []*Var{this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5} {
		if v != nil {
			items = append(items, v)
		}
	}
	return &j2gort.IterSlice{Items: items}
}

type ScopeCreator interface {
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
}
func (s *SyntacticScopeCreator_ScopeScanner) Reset() {}

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
	this_.reusableScanner = &SyntacticScopeCreator_ScopeScanner{compiler: compiler, redeclarationHandler: this_.redeclarationHandler}
}

func NewSyntacticScopeCreator_2(compiler AbstractCompiler, h SyntacticScopeCreator_RedeclarationHandler) *SyntacticScopeCreator {
	s := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_2(s, compiler, h)
	return s
}

func InitSyntacticScopeCreator_2(this_ *SyntacticScopeCreator, compiler AbstractCompiler, h SyntacticScopeCreator_RedeclarationHandler) {
	this_.compiler = compiler
	this_.redeclarationHandler = h
	this_.reusableScanner = &SyntacticScopeCreator_ScopeScanner{compiler: compiler, redeclarationHandler: h}
}

func (this_ *SyntacticScopeCreator) CreateScope(n pkg_rhino.Node, parent AbstractScope) AbstractScope {
	var p *Scope
	if parent != nil && !parent.IsJavaNil_() {
		p = parent.Untyped()
	}
	sc := NewScope(p, n)
	for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
		if c.IsName() {
			sc.Declare(c.GetString(), c, NewVar(c.GetString(), c, sc))
		}
	}
	return sc
}

type MemoizedScopeCreator struct {
	delegate          ScopeCreator
	scopesByScopeRoot j2gort.OrderedMap[pkg_rhino.Node, AbstractScope]
}

func (m *MemoizedScopeCreator) IsJavaNil_() bool { return m == nil }

func NewMemoizedScopeCreator(delegate ScopeCreator) *MemoizedScopeCreator {
	m := &MemoizedScopeCreator{}
	InitMemoizedScopeCreator(m, delegate)
	return m
}

func InitMemoizedScopeCreator(this_ *MemoizedScopeCreator, delegate ScopeCreator) {
	this_.delegate = delegate
	this_.scopesByScopeRoot = j2gort.NewOrderedMap[pkg_rhino.Node, AbstractScope]()
}

func (this_ *MemoizedScopeCreator) CreateScope(n pkg_rhino.Node, parent AbstractScope) AbstractScope {
	if existing := this_.scopesByScopeRoot.Get(n); existing != nil {
		return existing
	}
	sc := this_.delegate.CreateScope(n, parent)
	this_.scopesByScopeRoot.Put(n, sc)
	return sc
}
`)

	// 5. Minimal gen/jscomp/node_traversal.go (including NodeTraversal_TraverseScopeRoots + post-order adapter!)
	mustWrite("gen/jscomp/node_traversal.go", `package jscomp

import (
	"example.com/w25ctest/gen/j2gort"
	pkg_rhino "example.com/w25ctest/gen/rhino"
)

type NodeTraversal_Callback interface {
	Implements_jscomp_NodeTraversal_Callback()
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

type NodeTraversal_AbstractPostOrderCallbackInterface interface {
	Visit(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type NodeTraversal__Anon_1 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	cap_0_callback NodeTraversal_AbstractPostOrderCallbackInterface
}
func (this_ *NodeTraversal__Anon_1) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	this_.cap_0_callback.Visit(t, n, p)
}

func WrapPostOrderInterface_(target NodeTraversal_AbstractPostOrderCallbackInterface) NodeTraversal_Callback {
	return &NodeTraversal__Anon_1{cap_0_callback: target}
}

type NodeTraversal_AbstractPostOrderCallbackBase struct{}
func (*NodeTraversal_AbstractPostOrderCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (b *NodeTraversal_AbstractPostOrderCallbackBase) IsJavaNil_() bool                        { return b == nil }
func (*NodeTraversal_AbstractPostOrderCallbackBase) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}

type NodeTraversal_AbstractScopedCallbackBase struct{}
func (*NodeTraversal_AbstractScopedCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (b *NodeTraversal_AbstractScopedCallbackBase) IsJavaNil_() bool                        { return b == nil }
func (*NodeTraversal_AbstractScopedCallbackBase) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}
func (*NodeTraversal_AbstractScopedCallbackBase) Visit_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {}

type ControlFlowGraph[N any] interface {
	IsJavaNil_() bool
}
type cfgImpl struct{ slab j2gort.PooledSlab[int] }
func (c *cfgImpl) IsJavaNil_() bool { return c == nil }
func (c *cfgImpl) ReleaseSlabs_()   { c.slab.Release() }

type NodeTraversal_AbstractCfgCallback interface {
	NodeTraversal_ScopedCallback
	EnterScopeWithCfg(*NodeTraversal)
	ExitScopeWithCfg(*NodeTraversal)
}

type NodeTraversal_AbstractCfgCallbackBase struct {
	self_ NodeTraversal_AbstractCfgCallback
	cfgs  []any
}
func (*NodeTraversal_AbstractCfgCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (b *NodeTraversal_AbstractCfgCallbackBase) IsJavaNil_() bool                        { return b == nil }
func (b *NodeTraversal_AbstractCfgCallbackBase) BindSelf_(s any)                         { b.self_ = s.(NodeTraversal_AbstractCfgCallback) }
func (*NodeTraversal_AbstractCfgCallbackBase) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}
func (b *NodeTraversal_AbstractCfgCallbackBase) EnterScope(t *NodeTraversal) {
	c := &cfgImpl{}
	_ = c.slab.Alloc()
	b.cfgs = append(b.cfgs, c)
	b.self_.EnterScopeWithCfg(t)
}
func (b *NodeTraversal_AbstractCfgCallbackBase) ExitScope(t *NodeTraversal) {
	b.self_.ExitScopeWithCfg(t)
	last := b.cfgs[len(b.cfgs)-1]
	b.cfgs = b.cfgs[:len(b.cfgs)-1]
	if r, ok := last.(interface{ ReleaseSlabs_() }); ok {
		r.ReleaseSlabs_()
	}
}

func clinit_NodeTraversal() {}

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
func (this_ *NodeTraversal_Builder) SetScopeCreator(sc ScopeCreator) *NodeTraversal_Builder {
	this_.scopeCreator = sc
	return this_
}
func (this_ *NodeTraversal_Builder) Build() *NodeTraversal {
	nt := &NodeTraversal{}
	InitNodeTraversal(nt, this_)
	return nt
}
func (this_ *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverse_1(root)
}
func (this_ *NodeTraversal_Builder) TraverseRoots(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverseRoots_2(externs, root)
}

type NodeTraversal struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	curNode                                        pkg_rhino.Node
	scopeNodes                                     []pkg_rhino.Node
	scopeObjs                                      []AbstractScope
	scopeDepth                                     int32
	cfgs                                           any
	scopeCreator                                   ScopeCreator
	scopedCallback                                 NodeTraversal_ScopedCallback
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

func InitNodeTraversal(this_ *NodeTraversal, b *NodeTraversal_Builder) {
	this_.compiler = b.compiler
	this_.callback = b.callback
	this_.obeyDestructuringAndDefaultValueExecutionOrder = b.obeyDestructuringAndDefaultValueExecutionOrder
	if sc, ok := b.callback.(NodeTraversal_ScopedCallback); ok {
		this_.scopedCallback = sc
	}
	if b.scopeCreator == nil || b.scopeCreator.IsJavaNil_() {
		this_.scopeCreator = NewSyntacticScopeCreator_1(this_.compiler)
	} else {
		this_.scopeCreator = b.scopeCreator
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_ensureScopeCapacity_1() {
	if int(this_.scopeDepth) == len(this_.scopeNodes) {
		n := len(this_.scopeNodes)*2 + 4
		ns := make([]pkg_rhino.Node, n)
		copy(ns, this_.scopeNodes)
		this_.scopeNodes = ns
		os := make([]AbstractScope, n)
		copy(os, this_.scopeObjs)
		this_.scopeObjs = os
	}
}

func (this_ *NodeTraversal) pushScopeNode(n pkg_rhino.Node) {
	this_.priv_NodeTraversal_ensureScopeCapacity_1()
	this_.scopeNodes[this_.scopeDepth] = n
	this_.scopeObjs[this_.scopeDepth] = nil
	this_.scopeDepth++
	if this_.scopedCallback != nil {
		this_.scopedCallback.EnterScope(this_)
	}
}

func (this_ *NodeTraversal) popScope() {
	if this_.scopedCallback != nil {
		this_.scopedCallback.ExitScope(this_)
	}
	this_.scopeDepth--
	this_.scopeNodes[this_.scopeDepth] = 0
	this_.scopeObjs[this_.scopeDepth] = nil
}

func (this_ *NodeTraversal) GetScope() *Scope {
	var sc AbstractScope
	for i := int(this_.scopeDepth) - 1; i >= 0; i-- {
		if this_.scopeObjs[i] != nil {
			sc = this_.scopeObjs[i]
			for j := i + 1; j < int(this_.scopeDepth); j++ {
				sc = this_.scopeCreator.CreateScope(this_.scopeNodes[j], sc)
				this_.scopeObjs[j] = sc
			}
			return sc.Untyped()
		}
	}
	for j := 0; j < int(this_.scopeDepth); j++ {
		sc = this_.scopeCreator.CreateScope(this_.scopeNodes[j], sc)
		this_.scopeObjs[j] = sc
	}
	return sc.Untyped()
}

func (this_ *NodeTraversal) GetClosestContainerScope() *Scope { return this_.GetScope().GetClosestContainerScope() }
func (this_ *NodeTraversal) GetClosestHoistScope() *Scope     { return this_.GetScope().GetClosestHoistScope() }
func (this_ *NodeTraversal) GetScopeRoot() pkg_rhino.Node     { return this_.scopeNodes[this_.scopeDepth-1] }
func (this_ *NodeTraversal) InGlobalScope() bool              { return this_.scopeDepth <= 1 }

func (this_ *NodeTraversal) priv_NodeTraversal_traverse_1(root pkg_rhino.Node) {
	this_.pushScopeNode(root)
	this_.priv_NodeTraversal_traverseBranch_2(root, 0)
	this_.popScope()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseRoots_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.pushScopeNode(root)
	if externs != 0 {
		this_.priv_NodeTraversal_traverseBranch_2(externs, 0)
	}
	this_.priv_NodeTraversal_traverseBranch_2(root, 0)
	this_.popScope()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseScopeRoot_3(scopeRoot pkg_rhino.Node) {
	this_.pushScopeNode(scopeRoot)
	this_.priv_NodeTraversal_traverseBranch_2(scopeRoot, 0)
	this_.popScope()
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseBranch_2(n pkg_rhino.Node, parent pkg_rhino.Node) {
	this_.curNode = n
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	createsScope := n.IsBlock()
	if createsScope {
		this_.pushScopeNode(n)
	}
	for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
		this_.priv_NodeTraversal_traverseBranch_2(c, n)
	}
	this_.curNode = n
	this_.callback.Visit_3(this_, n, parent)
	if createsScope {
		this_.popScope()
	}
}

func NodeTraversal_Traverse(compiler AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).Traverse(root)
}

func NodeTraversal_TraverseRoots(compiler AbstractCompiler, cb NodeTraversal_Callback, externs pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).TraverseRoots(externs, root)
}

type NodeTraversal__Local_1_TraverseScopeRootsCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	scopeNode            pkg_rhino.Node
	insideScopeNode      bool
	cap_0_traverseNested bool
	cap_1_cb             NodeTraversal_Callback
}

func (l *NodeTraversal__Local_1_TraverseScopeRootsCallback) Visit_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {}

func NodeTraversal_TraverseScopeRoots(compiler AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback, traverseNested bool) {
	clinit_NodeTraversal()
	var scopeCreator *MemoizedScopeCreator = NewMemoizedScopeCreator(NewSyntacticScopeCreator_1(compiler))
	_ = scopeCreator
	localCb := &NodeTraversal__Local_1_TraverseScopeRootsCallback{
		scopeNode:            root,
		cap_0_traverseNested: traverseNested,
		cap_1_cb:             cb,
	}
	_ = localCb
	NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).SetScopeCreator(scopeCreator).Build().priv_NodeTraversal_traverseScopeRoot_3(root)
}
`)

	// 6. CombinedCompilerPass + Constituent callbacks (CheckUnreachableCode, ImplicitNullabilityCheck, ContextualRenameInverter, post-order interface callback)
	mustWrite("gen/jscomp/combined_compiler_pass.go", `package jscomp

import pkg_rhino "example.com/w25ctest/gen/rhino"

type CombinedCompilerPass_CallbackWrapper struct {
	callback       NodeTraversal_Callback
	scopedCallback NodeTraversal_ScopedCallback
	waiting        pkg_rhino.Node
}

func NewCombinedCompilerPass_CallbackWrapper(cb NodeTraversal_Callback) *CombinedCompilerPass_CallbackWrapper {
	w := &CombinedCompilerPass_CallbackWrapper{callback: cb}
	if sc, ok := cb.(NodeTraversal_ScopedCallback); ok {
		w.scopedCallback = sc
	}
	return w
}

type CombinedCompilerPass struct {
	callbacks []*CombinedCompilerPass_CallbackWrapper
	compiler  AbstractCompiler
}

func NewCombinedCompilerPass(compiler AbstractCompiler, cbs ...NodeTraversal_Callback) *CombinedCompilerPass {
	ws := make([]*CombinedCompilerPass_CallbackWrapper, len(cbs))
	for i, cb := range cbs {
		ws[i] = NewCombinedCompilerPass_CallbackWrapper(cb)
	}
	return &CombinedCompilerPass{compiler: compiler, callbacks: ws}
}

func (*CombinedCompilerPass) Implements_jscomp_NodeTraversal_Callback() {}
func (c *CombinedCompilerPass) IsJavaNil_() bool                        { return c == nil }

func (this_ *CombinedCompilerPass) EnterScope(t *NodeTraversal) {
	for _, w := range this_.callbacks {
		if w != nil && w.scopedCallback != nil {
			w.scopedCallback.EnterScope(t)
		}
	}
}

func (this_ *CombinedCompilerPass) ExitScope(t *NodeTraversal) {
	for _, w := range this_.callbacks {
		if w != nil && w.scopedCallback != nil {
			w.scopedCallback.ExitScope(t)
		}
	}
}

func (this_ *CombinedCompilerPass) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	anyActive := false
	for _, w := range this_.callbacks {
		if w != nil && w.callback.ShouldTraverse_3(t, n, parent) {
			anyActive = true
		}
	}
	return anyActive
}

func (this_ *CombinedCompilerPass) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	for _, w := range this_.callbacks {
		if w != nil {
			w.callback.Visit_3(t, n, parent)
		}
	}
}

func (this_ *CombinedCompilerPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_TraverseRoots(this_.compiler, this_, externs, root)
}
`)

	// 7. Wave 25C Callbacks + Pass-Boundary ScopeCreator & GlobalNamespace Owners
	mustWrite("gen/jscomp/wave25c_passes.go", `package jscomp

import (
	"example.com/w25ctest/gen/j2gort"
	pkg_rhino "example.com/w25ctest/gen/rhino"
)

// --- Constituent callbacks run inside CombinedCompilerPass ---
type CheckUnreachableCode struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler AbstractCompiler
	Checked  int
}

func NewCheckUnreachableCode(c AbstractCompiler) *CheckUnreachableCode {
	u := &CheckUnreachableCode{compiler: c}
	u.BindSelf_(u)
	return u
}
func (u *CheckUnreachableCode) EnterScopeWithCfg(t *NodeTraversal) {
	if t.GetScope() != nil {
		u.Checked++
	}
}
func (*CheckUnreachableCode) ExitScopeWithCfg(_ *NodeTraversal)                          {}
func (*CheckUnreachableCode) Visit_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {}

type ImplicitNullabilityCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
	Found    int
}

func (c *ImplicitNullabilityCheck) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		if v := j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString())); v != nil && v.IsLocal() {
			c.Found++
		}
	}
}

type Es6CheckModulePostOrder struct {
	compiler AbstractCompiler
	Count    int
}

func (e *Es6CheckModulePostOrder) Visit(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() && t.GetScope().GetSlot(n.GetString()) != nil {
		e.Count++
	}
}

// Negative control: retains *Var in slice!
type MakeDeclaredNamesUnique_ContextualRenameInverter struct {
	NodeTraversal_AbstractScopedCallbackBase
	saved []*Var
}

func (inv *MakeDeclaredNamesUnique_ContextualRenameInverter) EnterScope(t *NodeTraversal) {
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(t.GetScope().GetVarIterable()))
	for cur.HasNext() {
		inv.saved = append(inv.saved, j2gort.CastPtr[Var](cur.Next()))
	}
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter) ExitScope(_ *NodeTraversal) {}

// Negative control: LoopClosureTransformer retains *Var!
type Es6RewriteBlockScopedDeclaration_LoopClosureTransformer struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	vars []*Var
}
func (l *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		if v := j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString())); v != nil {
			l.vars = append(l.vars, v)
		}
	}
}

// --- Action 2: Widened ephemeral callbacks ---
type TypedScopeCreator_FirstOrderFunctionAnalyzer struct {
	NodeTraversal_AbstractScopedCallbackBase
	compiler AbstractCompiler
	Seen     int
}

func (a *TypedScopeCreator_FirstOrderFunctionAnalyzer) EnterScope(t *NodeTraversal) {
	names := j2gort.NewStringSet()
	t.GetScope().AddVarNamesTo(names)
	a.Seen += names.Count()
}
func (*TypedScopeCreator_FirstOrderFunctionAnalyzer) ExitScope(_ *NodeTraversal) {}
func (a *TypedScopeCreator_FirstOrderFunctionAnalyzer) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		container := t.GetClosestContainerScope()
		hoist := t.GetClosestHoistScope()
		if v := j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString())); v != nil {
			if container.HasSameContainerScope(v.GetScope()) && hoist != nil {
				a.Seen++
			}
		}
	}
}

type TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}
func (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) Visit_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {
}

type TypedScopeCreator struct {
	compiler AbstractCompiler
}

func (this_ *TypedScopeCreator) CreateInitialScope(root pkg_rhino.Node) int {
	var scopeCreator *MemoizedScopeCreator = NewMemoizedScopeCreator(NewSyntacticScopeCreator_1(this_.compiler))
	_ = scopeCreator
	analyzer := &TypedScopeCreator_FirstOrderFunctionAnalyzer{compiler: this_.compiler}
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(analyzer).SetScopeCreator(scopeCreator).Traverse(root)
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(&TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable{}).SetScopeCreator(scopeCreator).Traverse(root)
	return analyzer.Seen
}

type ClosureRewriteModule_ScriptPreprocessor struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	ExportsHits int32
}

func (p *ClosureRewriteModule_ScriptPreprocessor) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		if exportsVar := j2gort.CastPtr[Var](t.GetScope().GetOwnSlot(n.GetString())); exportsVar != nil && exportsVar.IsGoogModuleExports() {
			p.ExportsHits += exportsVar.GetImplicitGoogNamespaceStrength()
		}
	}
}

// --- Action 3A: Normalize (inline SetScopeCreator with NewSyntacticScopeCreator_2) ---
type Normalize_DuplicateDeclarationHandler struct {
	outer_ *Normalize
}

func NewNormalize_DuplicateDeclarationHandler(o *Normalize) *Normalize_DuplicateDeclarationHandler {
	return &Normalize_DuplicateDeclarationHandler{outer_: o}
}

type Normalize_PropagateConstantPropertyOverVars struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *Normalize
}

func NewNormalize_PropagateConstantPropertyOverVars(o *Normalize) *Normalize_PropagateConstantPropertyOverVars {
	return &Normalize_PropagateConstantPropertyOverVars{outer_: o}
}

func (p *Normalize_PropagateConstantPropertyOverVars) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() && t.GetScope().GetVar(n.GetString()) != nil {
		p.outer_.VisitedNames++
	}
}

type Normalize struct {
	compiler     AbstractCompiler
	VisitedNames int
}

func (this_ *Normalize) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_S_Builder().
		SetCompiler(this_.compiler).
		SetCallback_1_NodeTraversal_Callback(ToIface_NodeTraversal_Callback(NewNormalize_PropagateConstantPropertyOverVars(this_))).
		SetScopeCreator(ToIface_ScopeCreator(NewSyntacticScopeCreator_2(this_.compiler, NewNormalize_DuplicateDeclarationHandler(this_)))).
		TraverseRoots(externs, root)
}

// --- Action 3B: VarCheck (local scopeCreator via this_.priv_VarCheck_createScopeCreator_1()) ---
type VarCheck_RedeclHandler struct {
	outer_ *VarCheck
}

func NewVarCheck_RedeclHandler(o *VarCheck) *VarCheck_RedeclHandler {
	return &VarCheck_RedeclHandler{outer_: o}
}

type VarCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler    AbstractCompiler
	strictCheck bool
	CheckedVars int
}

func (this_ *VarCheck) priv_VarCheck_createScopeCreator_1() *SyntacticScopeCreator {
	if this_.strictCheck {
		return NewSyntacticScopeCreator_1(this_.compiler)
	}
	return NewSyntacticScopeCreator_2(this_.compiler, NewVarCheck_RedeclHandler(this_))
}

func (this_ *VarCheck) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() && t.GetScope().GetVar(n.GetString()) != nil {
		this_.CheckedVars++
	}
}

func (this_ *VarCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	var scopeCreator ScopeCreator = ToIface_ScopeCreator(this_.priv_VarCheck_createScopeCreator_1())
	_ = scopeCreator
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(this_).SetScopeCreator(scopeCreator).TraverseRoots(externs, root)
}

// --- Action 3C: ConstCheck (retains *Var in initializedConstants during Process_2) ---
type ConstCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler             AbstractCompiler
	initializedConstants j2gort.OrderedMap[*Var, bool]
	ConstHits            int
}

func NewConstCheck(c AbstractCompiler) *ConstCheck {
	return &ConstCheck{
		compiler:             c,
		initializedConstants: j2gort.NewOrderedMap[*Var, bool](),
	}
}

func (this_ *ConstCheck) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		if v := j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString())); v != nil {
			this_.initializedConstants.Put(v, true)
			this_.ConstHits++
		}
	}
}

func (this_ *ConstCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_TraverseRoots(this_.compiler, this_, externs, root)
}

// --- Action 3D: ClosureRewriteModule (retains globalScope + exportsToInline across ScriptRecorder -> ScriptUpdater -> declareSyntheticExterns!) ---
type ClosureRewriteModule_ScriptDescription struct {
	exportsToInline []*Var
}

type ClosureRewriteModule_GlobalRewriteState struct {
	scriptDescriptionsByGoogModuleNamespace j2gort.OrderedMap[string, *ClosureRewriteModule_ScriptDescription]
}

type ClosureRewriteModule_ScriptRecorder struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *ClosureRewriteModule
}

func NewClosureRewriteModule_ScriptRecorder(o *ClosureRewriteModule) *ClosureRewriteModule_ScriptRecorder {
	return &ClosureRewriteModule_ScriptRecorder{outer_: o}
}

func (r *ClosureRewriteModule_ScriptRecorder) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		if v := j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString())); v != nil {
			r.outer_.currentScript.exportsToInline = append(r.outer_.currentScript.exportsToInline, v)
			r.outer_.rewriteState.scriptDescriptionsByGoogModuleNamespace.Put("mod", r.outer_.currentScript)
		}
	}
}

type ClosureRewriteModule_ScriptUpdater struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *ClosureRewriteModule
}

func NewClosureRewriteModule_ScriptUpdater(o *ClosureRewriteModule) *ClosureRewriteModule_ScriptUpdater {
	return &ClosureRewriteModule_ScriptUpdater{outer_: o}
}

func (u *ClosureRewriteModule_ScriptUpdater) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if t.InGlobalScope() {
		u.outer_.globalScope = t.GetScope()
	}
	if n.IsName() && len(u.outer_.currentScript.exportsToInline) > 0 {
		u.outer_.UpdateModuleBody(n)
	}
}

type ClosureRewriteModule struct {
	compiler               AbstractCompiler
	globalScope            *Scope
	currentScript          *ClosureRewriteModule_ScriptDescription
	scriptStack            j2gort.Deque[*ClosureRewriteModule_ScriptDescription]
	rewriteState           *ClosureRewriteModule_GlobalRewriteState
	unrecognizedRequires   []string
	DeclaredSyntheticCount int
}

// Package-internal exported method (called strictly inside wave25c_passes.go, mirroring real closure_rewrite_module.go):
func (this_ *ClosureRewriteModule) UpdateModuleBody(_ pkg_rhino.Node) {
	if this_.currentScript.exportsToInline[0].GetName() == "" {
		panic("ScriptRecorder *Var corrupted before UpdateModuleBody ran!")
	}
}

func (this_ *ClosureRewriteModule) priv_ClosureRewriteModule_declareSyntheticExterns_1() {
	for _, req := range this_.unrecognizedRequires {
		// Reads this_.globalScope AFTER ScriptUpdater's TraverseRoots has returned!
		if this_.globalScope == nil || this_.globalScope.GetRootNode() == 0 {
			panic("globalScope was nil or zeroed before declareSyntheticExterns!")
		}
		if this_.globalScope.GetVar(req) == nil {
			this_.DeclaredSyntheticCount++
		}
	}
}

func (this_ *ClosureRewriteModule) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.currentScript = &ClosureRewriteModule_ScriptDescription{}
	this_.scriptStack.Push(this_.currentScript)
	for c := root.GetFirstChild(); c != 0; c = c.GetNext() {
		if c.IsScript() {
			NodeTraversal_Traverse(this_.compiler, c, NewClosureRewriteModule_ScriptRecorder(this_))
		}
	}
	NodeTraversal_TraverseRoots(this_.compiler, NewClosureRewriteModule_ScriptUpdater(this_), externs, root)
	this_.priv_ClosureRewriteModule_declareSyntheticExterns_1()
}

// --- Action 3E: GlobalNamespace + Owning Passes (InlineAndCollapseProperties, InlineAliases, ProcessDefines, PolymerPass, local GlobalNamespace) ---
type GlobalNamespace_Ref struct {
	node  pkg_rhino.Node
	scope *Scope
}

type GlobalNamespace_Name struct {
	name string
	refs []*GlobalNamespace_Ref
}

type GlobalNamespace_BuildGlobalNamespace struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *GlobalNamespace
}

func NewGlobalNamespace_BuildGlobalNamespace(o *GlobalNamespace) *GlobalNamespace_BuildGlobalNamespace {
	return &GlobalNamespace_BuildGlobalNamespace{outer_: o}
}

func (b *GlobalNamespace_BuildGlobalNamespace) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		nm := b.outer_.nameMap.Get(n.GetString())
		if nm == nil {
			nm = &GlobalNamespace_Name{name: n.GetString()}
			b.outer_.nameMap.Put(n.GetString(), nm)
			b.outer_.globalNames = append(b.outer_.globalNames, nm)
		}
		nm.refs = append(nm.refs, &GlobalNamespace_Ref{node: n, scope: t.GetScope()})
	}
}

type GlobalNamespace struct {
	compiler        AbstractCompiler
	root            pkg_rhino.Node
	externsRoot     pkg_rhino.Node
	generated       bool
	globalNames     []*GlobalNamespace_Name
	nameMap         j2gort.OrderedMap[string, *GlobalNamespace_Name]
	nameMapByModule j2gort.OrderedMap[string, *GlobalNamespace_Name]
}

func NewGlobalNamespace_2(compiler AbstractCompiler, externs pkg_rhino.Node, root pkg_rhino.Node) *GlobalNamespace {
	return &GlobalNamespace{
		compiler:        compiler,
		externsRoot:     externs,
		root:            root,
		nameMap:         j2gort.NewOrderedMap[string, *GlobalNamespace_Name](),
		nameMapByModule: j2gort.NewOrderedMap[string, *GlobalNamespace_Name](),
	}
}

func (this_ *GlobalNamespace) EnsureGenerated() {
	if !this_.generated {
		this_.priv_GlobalNamespace_process_1()
		this_.generated = true
	}
}

func (this_ *GlobalNamespace) priv_GlobalNamespace_process_1() {
	var builder *NodeTraversal_Builder = NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(NewGlobalNamespace_BuildGlobalNamespace(this_))
	builder.TraverseRoots(this_.externsRoot, this_.root)
}

func (this_ *GlobalNamespace) GetSlot(name string) *GlobalNamespace_Name {
	this_.EnsureGenerated()
	return this_.nameMap.Get(name)
}

type InlineAndCollapseProperties struct {
	compiler                   AbstractCompiler
	namespace                  *GlobalNamespace
	testAggressiveInliningOnly bool
	Collapsed                  int
}

func (this_ *InlineAndCollapseProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.namespace = NewGlobalNamespace_2(this_.compiler, externs, root)
	if nm := this_.namespace.GetSlot("g1"); nm != nil && len(nm.refs) > 0 {
		if nm.refs[0].scope == nil || nm.refs[0].scope.GetRootNode() == 0 {
			panic("GlobalNamespace_Ref.scope zeroed prematurely before InlineAndCollapseProperties finished!")
		}
		this_.Collapsed++
	}
}

type InlineAndCollapseProperties_InlineAliases struct {
	compiler  AbstractCompiler
	namespace *GlobalNamespace
	Aliased   int
}

func (this_ *InlineAndCollapseProperties_InlineAliases) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.namespace = NewGlobalNamespace_2(this_.compiler, externs, root)
	if nm := this_.namespace.GetSlot("g1"); nm != nil {
		this_.Aliased += len(nm.refs)
	}
}

type NamespaceSupplier interface {
	Get() *GlobalNamespace
}

type ProcessDefines struct {
	compiler           AbstractCompiler
	namespace          *GlobalNamespace
	namespaceSupplier  NamespaceSupplier
	defineByDefineName j2gort.OrderedMap[string, *GlobalNamespace_Name]
	DefinesFound       int
}

func NewProcessDefines(c AbstractCompiler, sup NamespaceSupplier) *ProcessDefines {
	return &ProcessDefines{
		compiler:           c,
		namespaceSupplier:  sup,
		defineByDefineName: j2gort.NewOrderedMap[string, *GlobalNamespace_Name](),
	}
}

func (this_ *ProcessDefines) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	if this_.namespaceSupplier != nil {
		this_.namespace = this_.namespaceSupplier.Get()
	} else {
		this_.namespace = NewGlobalNamespace_2(this_.compiler, externs, root)
	}
	if nm := this_.namespace.GetSlot("g1"); nm != nil && len(nm.refs) > 0 {
		if nm.refs[0].scope == nil || nm.refs[0].scope.GetRootNode() == 0 {
			panic("ProcessDefines namespace ref scope zeroed!")
		}
		this_.defineByDefineName.Put("g1", nm)
		this_.DefinesFound++
	}
}

type PolymerPass_BehaviorExtractor struct{ tag int }

type PolymerPass struct {
	compiler          AbstractCompiler
	globalNames       *GlobalNamespace
	behaviorExtractor *PolymerPass_BehaviorExtractor
	PolymerCount      int
}

func (this_ *PolymerPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.behaviorExtractor = &PolymerPass_BehaviorExtractor{tag: 1}
	this_.globalNames = NewGlobalNamespace_2(this_.compiler, externs, root)
	if nm := this_.globalNames.GetSlot("g1"); nm != nil {
		this_.PolymerCount += len(nm.refs)
	}
}

func InspectLocalGlobalNamespace(c AbstractCompiler, externs pkg_rhino.Node, root pkg_rhino.Node) {
	ns := NewGlobalNamespace_2(c, externs, root)
	ns.EnsureGenerated()
	if nm := ns.GetSlot("g1"); nm == nil || nm.refs[0].scope.GetRootNode() == 0 {
		panic("local GlobalNamespace failed")
	}
}
`)

	// Run SlabAllocLoweringPass + NodeTraversalLoweringPass (Run 1) on gen/jscomp!
	jscompDir := filepath.Join(tmpDir, "gen", "jscomp")
	passList := []passes.Pass{
		passes.NewJ2gortPeepholePass(),
		passes.NewSlabAllocLoweringPass(),
		passes.NewNodeTraversalLoweringPass(),
	}
	res1, err := passes.RunOnPackageDir(jscompDir, passList, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if res1.ByPass["slab_alloc_lowering"] < 8 || res1.ByPass["node_traversal_lowering"] < 6 {
		t.Fatalf("expected substantial Wave 25C lowering sites in both passes, got %+v", res1.ByPass)
	}

	// Verify Run-2 idempotency (0 sites, 0 modified files)!
	res2, err := passes.RunOnPackageDir(jscompDir, passList, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 modified files on Run 2 (idempotency), got %+v", res2)
	}

	readLowered := func(name string) string {
		t.Helper()
		b, rErr := os.ReadFile(filepath.Join(jscompDir, name))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}

	// 1. Inspect node_traversal.go & combined_compiler_pass.go
	ntCode := readLowered("node_traversal.go")
	for _, want := range []string{
		"func canReleaseSyntacticScopesForCallback_(cb NodeTraversal_Callback) bool",
		"CanReleaseSyntacticScopesDynamic_() bool",
		"this_.canReleaseScopes_ = canReleaseSyntacticScopesForCallback_(",
		"canRel_ := canReleaseSyntacticScopesForCallback_(",
	} {
		if !strings.Contains(ntCode, want) {
			t.Fatalf("expected lowered node_traversal.go to contain %q\n--- node_traversal.go ---\n%s", want, ntCode)
		}
	}

	ccpCode := readLowered("combined_compiler_pass.go")
	if !strings.Contains(ccpCode, "func (this_ *CombinedCompilerPass) CanReleaseSyntacticScopesDynamic_() bool") {
		t.Fatalf("expected CombinedCompilerPass.CanReleaseSyntacticScopesDynamic_() in combined_compiler_pass.go:\n%s", ccpCode)
	}
	if strings.Contains(ccpCode, "func (*CombinedCompilerPass) CanReleaseSyntacticScopes_()") {
		t.Fatalf("CombinedCompilerPass MUST NOT receive unconditional CanReleaseSyntacticScopes_()!")
	}

	// 2. Inspect scope.go (MemoizedScopeCreator.ReleaseSlabs_() releasing scopesByScopeRoot!)
	scopeCode := readLowered("scope.go")
	if !strings.Contains(scopeCode, "this_.scopesByScopeRoot.ReleaseToPool()") {
		t.Fatalf("expected MemoizedScopeCreator.ReleaseSlabs_() to release scopesByScopeRoot to pool:\n%s", scopeCode)
	}

	// 3. Inspect wave25c_passes.go positive & negative soundness markers + slab releases!
	passesCode := readLowered("wave25c_passes.go")
	for _, want := range []string{
		"func (*CheckUnreachableCode) CanReleaseSyntacticScopes_()",
		"func (*ImplicitNullabilityCheck) CanReleaseSyntacticScopes_()",
		"func (*TypedScopeCreator_FirstOrderFunctionAnalyzer) CanReleaseSyntacticScopes_()",
		"func (*ClosureRewriteModule_ScriptPreprocessor) CanReleaseSyntacticScopes_()",
		"func (*Normalize_PropagateConstantPropertyOverVars) CanReleaseSyntacticScopes_()",
		"func (*VarCheck) CanReleaseSyntacticScopes_()",
		"this_.initializedConstants.Clear()",
		"this_.globalScope = nil",
		"scopeCreator_",
		"InitSyntacticScopeCreator_1(&this_.scopeCreator_",
		"func (this_ *GlobalNamespace) ReleaseSlabs_()",
		"!this_.testAggressiveInliningOnly && this_.namespace != nil",
		"ownsNamespace_ := this_.namespaceSupplier == nil && this_.namespace == nil",
		"this_.globalNames.ReleaseSlabs_()",
		"defer ns.ReleaseSlabs_()",
	} {
		if !strings.Contains(passesCode, want) {
			t.Fatalf("expected lowered wave25c_passes.go to contain %q\n--- wave25c_passes.go ---\n%s", want, passesCode)
		}
	}
	for _, forbidden := range []string{
		"func (*MakeDeclaredNamesUnique_ContextualRenameInverter) CanReleaseSyntacticScopes_()",
		"func (*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) CanReleaseSyntacticScopes_()",
		"func (*GlobalNamespace_BuildGlobalNamespace) CanReleaseSyntacticScopes_()",
		"func (*ClosureRewriteModule_ScriptRecorder) CanReleaseSyntacticScopes_()",
		"func (*ClosureRewriteModule_ScriptUpdater) CanReleaseSyntacticScopes_()",
		"func (*ConstCheck) CanReleaseSyntacticScopes_()",
	} {
		if strings.Contains(passesCode, forbidden) {
			t.Fatalf("unsound CanReleaseSyntacticScopes_() emitted: %q!", forbidden)
		}
	}

	// 4. Write runtime behavioral, soundness, slab-reclamation, and -race test in gen/jscomp!
	mustWrite("gen/jscomp/wave25c_rt_test.go", `package jscomp

import (
	"sync"
	"testing"

	"example.com/w25ctest/gen/j2gort"
	pkg_rhino "example.com/w25ctest/gen/rhino"
)

func buildWave25CSampleAST() (pkg_rhino.Node, pkg_rhino.Node) {
	ext := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
	ext.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "extSym"))

	root := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
	subScript := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
	subScript.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "g1"))
	subScript.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "exports"))
	blk := pkg_rhino.NewNode_1(pkg_rhino.TokBLOCK)
	blk.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "locA"))
	blk.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "locB"))
	subScript.AddChildToBack(blk)
	root.AddChildToBack(subScript)
	root.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "g1"))
	root.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "exports"))
	root.AddChildToBack(blk)
	return ext, root
}

type fixedNamespaceSupplier struct{ ns *GlobalNamespace }
func (s *fixedNamespaceSupplier) Get() *GlobalNamespace { return s.ns }

func TestWave25C_CombinedCompilerPassDynamicRelease(t *testing.T) {
	comp := &Compiler{}
	ext, root := buildWave25CSampleAST()

	// Edge cases: nil receiver, empty callbacks, nil wrapper element MUST return false!
	var nilPass *CombinedCompilerPass
	if nilPass.CanReleaseSyntacticScopesDynamic_() {
		t.Fatal("nil CombinedCompilerPass must return false")
	}
	emptyPass := &CombinedCompilerPass{compiler: comp}
	if emptyPass.CanReleaseSyntacticScopesDynamic_() {
		t.Fatal("empty CombinedCompilerPass must return false")
	}
	nilWrapperPass := &CombinedCompilerPass{
		compiler:  comp,
		callbacks: []*CombinedCompilerPass_CallbackWrapper{nil},
	}
	if nilWrapperPass.CanReleaseSyntacticScopesDynamic_() {
		t.Fatal("CombinedCompilerPass with nil wrapper must return false")
	}

	// Positive case: CheckUnreachableCode + ImplicitNullabilityCheck + wrapped post-order interface + nested CombinedCompilerPass!
	cuc := NewCheckUnreachableCode(comp)
	inc := &ImplicitNullabilityCheck{compiler: comp}
	es6 := WrapPostOrderInterface_(&Es6CheckModulePostOrder{compiler: comp})
	nested := NewCombinedCompilerPass(comp, inc, es6)
	combinedSafe := NewCombinedCompilerPass(comp, cuc, nested)
	if !combinedSafe.CanReleaseSyntacticScopesDynamic_() {
		t.Fatal("expected CombinedCompilerPass of releasable callbacks to return true")
	}

	beforeRel := j2gort.SlabReleaseCount.Load()
	combinedSafe.Process_2(ext, root)
	afterRel := j2gort.SlabReleaseCount.Load()
	if afterRel <= beforeRel {
		t.Fatalf("expected CombinedCompilerPass.Process_2 to release SyntacticScopeCreator slabs (before=%d, after=%d)", beforeRel, afterRel)
	}
	if cuc.Checked == 0 || inc.Found == 0 {
		t.Fatalf("expected callbacks to run: Checked=%d Found=%d", cuc.Checked, inc.Found)
	}

	// NodeTraversal_TraverseScopeRoots with CombinedCompilerPass MUST ALSO release MemoizedScopeCreator + SyntacticScopeCreator slabs!
	beforeMapRel := j2gort.OrderedMapReleaseCount.Load()
	beforeRel = j2gort.SlabReleaseCount.Load()
	NodeTraversal_TraverseScopeRoots(comp, root, combinedSafe, true)
	if j2gort.SlabReleaseCount.Load() <= beforeRel || j2gort.OrderedMapReleaseCount.Load() <= beforeMapRel {
		t.Fatalf("expected NodeTraversal_TraverseScopeRoots with CombinedCompilerPass to release MemoizedScopeCreator + SyntacticScopeCreator slabs")
	}

	// Negative soundness case: CombinedCompilerPass containing ContextualRenameInverter MUST return false and retain live *Var!
	inv := &MakeDeclaredNamesUnique_ContextualRenameInverter{}
	combinedMixed := NewCombinedCompilerPass(comp, NewCheckUnreachableCode(comp), inv)
	if combinedMixed.CanReleaseSyntacticScopesDynamic_() {
		t.Fatal("CombinedCompilerPass containing ContextualRenameInverter MUST return false!")
	}
	combinedMixed.Process_2(ext, root)
	if len(inv.saved) == 0 || inv.saved[0].GetName() == "" {
		t.Fatalf("expected ContextualRenameInverter retained *Var to remain intact after CombinedCompilerPass, got %+v", inv.saved)
	}
}

func TestWave25C_WidenedCallbacksAndPassBoundarySlabReleases(t *testing.T) {
	comp := &Compiler{}
	ext, root := buildWave25CSampleAST()

	// 1. ClosureRewriteModule_ScriptPreprocessor (ephemeral IsGoogModuleExports + GetImplicitGoogNamespaceStrength)
	before := j2gort.SlabReleaseCount.Load()
	prep := &ClosureRewriteModule_ScriptPreprocessor{}
	NodeTraversal_Traverse(comp, root, prep)
	if prep.ExportsHits == 0 || j2gort.SlabReleaseCount.Load() <= before {
		t.Fatalf("ScriptPreprocessor failed or did not release slabs: hits=%d", prep.ExportsHits)
	}

	// 2. TypedScopeCreator.CreateInitialScope (MemoizedScopeCreator + FirstOrderFunctionAnalyzer)
	before = j2gort.SlabReleaseCount.Load()
	beforeMap := j2gort.OrderedMapReleaseCount.Load()
	tsc := &TypedScopeCreator{compiler: comp}
	if seen := tsc.CreateInitialScope(root); seen == 0 {
		t.Fatal("expected TypedScopeCreator.CreateInitialScope seen > 0")
	}
	if j2gort.SlabReleaseCount.Load() <= before || j2gort.OrderedMapReleaseCount.Load() <= beforeMap {
		t.Fatal("expected TypedScopeCreator.CreateInitialScope to release MemoizedScopeCreator map + inner SyntacticScopeCreator slabs")
	}

	// 3. Normalize.Process_2 (inline SetScopeCreator(NewSyntacticScopeCreator_2(...)))
	before = j2gort.SlabReleaseCount.Load()
	norm := &Normalize{compiler: comp}
	norm.Process_2(ext, root)
	if norm.VisitedNames == 0 || j2gort.SlabReleaseCount.Load() <= before {
		t.Fatalf("Normalize.Process_2 failed to visit or release slabs: visited=%d", norm.VisitedNames)
	}

	// 4. VarCheck.Process_2 (local scopeCreator from this_.priv_VarCheck_createScopeCreator_1())
	before = j2gort.SlabReleaseCount.Load()
	vc := &VarCheck{compiler: comp}
	vc.Process_2(ext, root)
	if vc.CheckedVars == 0 || j2gort.SlabReleaseCount.Load() <= before {
		t.Fatalf("VarCheck.Process_2 failed to check vars or release slabs: checked=%d", vc.CheckedVars)
	}

	// 5. ConstCheck.Process_2 (clears initializedConstants and releases pass-scoped SyntacticScopeCreator)
	before = j2gort.SlabReleaseCount.Load()
	cc := NewConstCheck(comp)
	cc.Process_2(ext, root)
	if cc.ConstHits == 0 || !cc.initializedConstants.IsEmpty() || j2gort.SlabReleaseCount.Load() <= before {
		t.Fatalf("ConstCheck.Process_2 failed cleanup or slab release: hits=%d empty=%v", cc.ConstHits, cc.initializedConstants.IsEmpty())
	}

	// 6. ClosureRewriteModule.Process_2 (retains globalScope & exportsToInline through declareSyntheticExterns, then cleans up & releases slabs!)
	before = j2gort.SlabReleaseCount.Load()
	crm := &ClosureRewriteModule{
		compiler:             comp,
		rewriteState:         &ClosureRewriteModule_GlobalRewriteState{scriptDescriptionsByGoogModuleNamespace: j2gort.NewOrderedMap[string, *ClosureRewriteModule_ScriptDescription]()},
		unrecognizedRequires: []string{"missingSyntheticSym"},
	}
	crm.Process_2(ext, root)
	if crm.DeclaredSyntheticCount != 1 || crm.globalScope != nil || crm.currentScript != nil ||
		!crm.scriptStack.IsEmpty() || !crm.rewriteState.scriptDescriptionsByGoogModuleNamespace.IsEmpty() {
		t.Fatalf("ClosureRewriteModule cleanup mismatch: declared=%d gs=%v cs=%v stackEmpty=%v rsEmpty=%v",
			crm.DeclaredSyntheticCount, crm.globalScope, crm.currentScript, crm.scriptStack.IsEmpty(), crm.rewriteState.scriptDescriptionsByGoogModuleNamespace.IsEmpty())
	}
	if j2gort.SlabReleaseCount.Load() <= before {
		t.Fatal("expected ClosureRewriteModule.Process_2 to release pass-scoped SyntacticScopeCreator slabs on exit")
	}

	// 7. GlobalNamespace + Owning Passes!
	// 7a. InlineAndCollapseProperties normal mode releases slabs and nils namespace:
	before = j2gort.SlabReleaseCount.Load()
	iacp := &InlineAndCollapseProperties{compiler: comp}
	iacp.Process_2(ext, root)
	if iacp.Collapsed != 1 || iacp.namespace != nil || j2gort.SlabReleaseCount.Load() <= before {
		t.Fatalf("InlineAndCollapseProperties normal mode failed: collapsed=%d ns=%v", iacp.Collapsed, iacp.namespace)
	}
	// 7b. InlineAndCollapseProperties testAggressiveInliningOnly == true MUST retain namespace!
	iacpTest := &InlineAndCollapseProperties{compiler: comp, testAggressiveInliningOnly: true}
	iacpTest.Process_2(ext, root)
	if iacpTest.namespace == nil || iacpTest.namespace.GetSlot("g1").refs[0].scope.GetRootNode() == 0 {
		t.Fatal("InlineAndCollapseProperties testAggressiveInliningOnly must retain live GlobalNamespace!")
	}
	iacpTest.namespace.ReleaseSlabs_()

	// 7c. InlineAndCollapseProperties_InlineAliases:
	before = j2gort.SlabReleaseCount.Load()
	ia := &InlineAndCollapseProperties_InlineAliases{compiler: comp}
	ia.Process_2(ext, root)
	if ia.Aliased == 0 || ia.namespace != nil || j2gort.SlabReleaseCount.Load() <= before {
		t.Fatalf("InlineAliases failed: aliased=%d ns=%v", ia.Aliased, ia.namespace)
	}

	// 7d. ProcessDefines with owned namespace releases slabs; with external namespaceSupplier retains caller's namespace!
	before = j2gort.SlabReleaseCount.Load()
	pdOwned := NewProcessDefines(comp, nil)
	pdOwned.Process_2(ext, root)
	if pdOwned.DefinesFound != 1 || pdOwned.namespace != nil || !pdOwned.defineByDefineName.IsEmpty() || j2gort.SlabReleaseCount.Load() <= before {
		t.Fatalf("ProcessDefines owned namespace failed cleanup/release")
	}

	sharedNS := NewGlobalNamespace_2(comp, ext, root)
	sharedNS.EnsureGenerated()
	before = j2gort.SlabReleaseCount.Load()
	pdSupplied := NewProcessDefines(comp, &fixedNamespaceSupplier{ns: sharedNS})
	pdSupplied.Process_2(ext, root)
	if j2gort.SlabReleaseCount.Load() != before || sharedNS.GetSlot("g1").refs[0].scope.GetRootNode() == 0 {
		t.Fatal("ProcessDefines with external namespaceSupplier MUST NOT release caller-owned GlobalNamespace!")
	}
	sharedNS.ReleaseSlabs_()

	// 7e. PolymerPass + non-escaping local GlobalNamespace:
	before = j2gort.SlabReleaseCount.Load()
	pp := &PolymerPass{compiler: comp}
	pp.Process_2(ext, root)
	if pp.PolymerCount == 0 || pp.globalNames != nil || pp.behaviorExtractor != nil || j2gort.SlabReleaseCount.Load() <= before {
		t.Fatalf("PolymerPass failed cleanup/release")
	}

	before = j2gort.SlabReleaseCount.Load()
	InspectLocalGlobalNamespace(comp, ext, root)
	if j2gort.SlabReleaseCount.Load() <= before {
		t.Fatal("expected local non-escaping GlobalNamespace to release slabs on exit")
	}

	// 8. Concurrent -race stress test across CombinedCompilerPass, Normalize, VarCheck, and GlobalNamespace!
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 15; i++ {
				cp := NewCombinedCompilerPass(comp, NewCheckUnreachableCode(comp), &ImplicitNullabilityCheck{compiler: comp})
				cp.Process_2(ext, root)
				n := &Normalize{compiler: comp}
				n.Process_2(ext, root)
				p := &InlineAndCollapseProperties{compiler: comp}
				p.Process_2(ext, root)
			}
		}()
	}
	wg.Wait()
}
`)

	cmd := exec.Command("go", "test", "-timeout=30s", "-race", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race in wave25c workspace failed: %v\n%s", err, string(out))
	}
}

func TestWave25C_RealGenJscompCoverageAndIdempotency(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/rhino or gen/j2gort is awaiting regeneration; skipping real gen/jscomp Wave25C check")
	}

	realJscompDir := filepath.Join(repoRoot, "transpiled", "gen", "jscomp")
	entries, err := os.ReadDir(realJscompDir)
	if err != nil {
		t.Fatal(err)
	}

	// Copy real gen/jscomp/*.go into an isolated TempDir module using `replace` to repoRoot
	// so we can verify Run 1 lowered output, `go build`, AND Run 2 0-site idempotency without touching gen/!
	tmpDir := t.TempDir()
	origModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/w25crealgen", 1) +
		fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, sErr := os.ReadFile(filepath.Join(repoRoot, "go.sum")); sErr == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}
	copyDir := filepath.Join(tmpDir, "jscomp")
	if err := os.MkdirAll(copyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") {
			continue
		}
		data, rErr := os.ReadFile(filepath.Join(realJscompDir, ent.Name()))
		if rErr != nil {
			t.Fatal(rErr)
		}
		if wErr := os.WriteFile(filepath.Join(copyDir, ent.Name()), data, 0o644); wErr != nil {
			t.Fatal(wErr)
		}
	}

	passList := []passes.Pass{
		passes.NewSlabAllocLoweringPass(),
		passes.NewNodeTraversalLoweringPass(),
	}

	res1, err := passes.RunOnPackageDir(copyDir, passList, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on copied real gen/jscomp failed: %v", err)
	}
	t.Logf("Wave 25C Real gen/jscomp Run 1 sites: total=%d modifiedFiles=%d byPass=%+v",
		res1.TotalSites, res1.FilesModified, res1.ByPass)

	readCopy := func(name string) string {
		b, rErr := os.ReadFile(filepath.Join(copyDir, name))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}

	// 1. Action 1: CombinedCompilerPass dynamic CanReleaseSyntacticScopesDynamic_ + node_traversal helper
	nt := readCopy("node_traversal.go")
	if !strings.Contains(nt, "func canReleaseSyntacticScopesForCallback_(cb NodeTraversal_Callback) bool") {
		t.Fatal("expected canReleaseSyntacticScopesForCallback_ in lowered real node_traversal.go")
	}
	ccp := readCopy("combined_compiler_pass.go")
	if !strings.Contains(ccp, "func (this_ *CombinedCompilerPass) CanReleaseSyntacticScopesDynamic_() bool") ||
		strings.Contains(ccp, "func (*CombinedCompilerPass) CanReleaseSyntacticScopes_()") {
		t.Fatal("CombinedCompilerPass dynamic marker mismatch in real combined_compiler_pass.go")
	}

	// 2. Action 2: Widened callbacks + negative soundness controls in real gen/jscomp
	tsc := readCopy("typed_scope_creator.go")
	if !strings.Contains(tsc, "func (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) CanReleaseSyntacticScopes_()") {
		t.Fatal("expected TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable to receive CanReleaseSyntacticScopes_()")
	}
	if !strings.Contains(tsc, "ReleaseSlabs_()") {
		t.Fatal("expected TypedScopeCreator.CreateInitialScope to release local MemoizedScopeCreator slabs")
	}
	crm := readCopy("closure_rewrite_module.go")
	if !strings.Contains(crm, "func (*ClosureRewriteModule_ScriptPreprocessor) CanReleaseSyntacticScopes_()") {
		t.Fatal("expected ClosureRewriteModule_ScriptPreprocessor to receive CanReleaseSyntacticScopes_()")
	}
	if strings.Contains(crm, "func (*ClosureRewriteModule_ScriptRecorder) CanReleaseSyntacticScopes_()") ||
		strings.Contains(crm, "func (*ClosureRewriteModule_ScriptUpdater) CanReleaseSyntacticScopes_()") {
		t.Fatal("ClosureRewriteModule ScriptRecorder/ScriptUpdater MUST NOT receive CanReleaseSyntacticScopes_()")
	}
	if !strings.Contains(crm, "sc_.ReleaseSlabs_()") || !strings.Contains(crm, "this_.globalScope = nil") {
		t.Fatal("expected ClosureRewriteModule.Process_2 pass-boundary ScopeCreator release")
	}

	// 3. Action 3: Normalize, VarCheck, ConstCheck, MemoizedScopeCreator, and GlobalNamespace owners
	if !strings.Contains(readCopy("normalize.go"), "ReleaseSlabs_()") {
		t.Fatal("expected Normalize.Process_2 inline SetScopeCreator to release SyntacticScopeCreator slabs")
	}
	if !strings.Contains(readCopy("var_check.go"), "ReleaseSlabs_()") {
		t.Fatal("expected VarCheck.Process_2 to release SyntacticScopeCreator slabs")
	}
	cc := readCopy("const_check.go")
	if strings.Contains(cc, "func (*ConstCheck) CanReleaseSyntacticScopes_()") ||
		!strings.Contains(cc, "this_.initializedConstants.Clear()") ||
		!strings.Contains(cc, "sc_.ReleaseSlabs_()") {
		t.Fatal("ConstCheck.Process_2 lowering mismatch")
	}
	if !strings.Contains(readCopy("memoized_scope_creator.go"), "ReleaseToPool()") {
		t.Fatal("expected MemoizedScopeCreator.ReleaseSlabs_() to release scopesByScopeRoot to pool")
	}
	gn := readCopy("global_namespace.go")
	if strings.Contains(gn, "func (*GlobalNamespace_BuildGlobalNamespace) CanReleaseSyntacticScopes_()") ||
		!strings.Contains(gn, "scopeCreator_") ||
		!strings.Contains(gn, "InitSyntacticScopeCreator_1(&this_.scopeCreator_") ||
		!strings.Contains(gn, "func (this_ *GlobalNamespace) ReleaseSlabs_()") {
		t.Fatal("GlobalNamespace slab release lowering mismatch")
	}
	iacp := readCopy("inline_and_collapse_properties.go")
	if !strings.Contains(iacp, "!this_.testAggressiveInliningOnly && this_.namespace != nil") ||
		!strings.Contains(iacp, "this_.namespace.ReleaseSlabs_()") {
		t.Fatal("InlineAndCollapseProperties GlobalNamespace slab release mismatch")
	}
	pd := readCopy("process_defines.go")
	if !strings.Contains(pd, "ownsNamespace_ := this_.namespaceSupplier == nil && this_.namespace == nil") ||
		!strings.Contains(pd, "this_.namespace.ReleaseSlabs_()") {
		t.Fatal("ProcessDefines GlobalNamespace slab release mismatch")
	}
	if !strings.Contains(readCopy("polymer_pass.go"), "this_.globalNames.ReleaseSlabs_()") {
		t.Fatal("PolymerPass GlobalNamespace slab release mismatch")
	}

	// 4. Compile ALL 450+ lowered real jscomp files with `go build ./jscomp`!
	buildCmd := exec.Command("go", "build", "./jscomp")
	buildCmd.Dir = tmpDir
	if out, bErr := buildCmd.CombinedOutput(); bErr != nil {
		t.Fatalf("go build ./jscomp on lowered real gen/jscomp failed: %v\n%s", bErr, string(out))
	}

	// 5. Run 2 MUST be 100% idempotent (0 sites, 0 modified files)!
	res2, err := passes.RunOnPackageDir(copyDir, passList, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on lowered real gen/jscomp failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 modified files on Run 2 of real gen/jscomp, got %+v", res2)
	}
}

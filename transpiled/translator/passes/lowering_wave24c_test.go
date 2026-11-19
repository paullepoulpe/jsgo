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

func TestWave24C_EndToEndBehaviorZeroAllocAndRace(t *testing.T) {
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

	mustWrite("go.mod", "module example.com/w24ctest\n\ngo 1.24\n")

	// 1. Minimal shims/nodestore with NodeAllocator, Intern, Epoch, SetOrigNameID, ClonePropsFrom.
	mustWrite("shims/nodestore/nodestore.go", `package nodestore

import (
	"math"
	"sync"
	"sync/atomic"
)

type Node uint32
type StringID uint32

const EmptyStringID StringID = 1

type NodeRec struct {
	Token      uint8
	KindFlags  uint8
	Payload    uint32
	SourceFile uint32
	PropFlags  uint32
	OrigName   StringID
	DoubleVal  float64
	First      Node
	Next       Node
	Parent     Node
}

var (
	mu           sync.RWMutex
	recs         = make([]NodeRec, 1, 4096)
	strs         = []string{"", ""}
	strMap       = map[string]StringID{"": EmptyStringID}
	epoch        atomic.Uint64
	GlobalIntern atomic.Int64
	Alloc1Calls  atomic.Int64
	AllocStrCall atomic.Int64
	AllocNumCall atomic.Int64
	FlushCalls   atomic.Int64
)

func init() { epoch.Store(1) }

func Epoch() uint64     { return epoch.Load() }
func BumpEpochForTest() { epoch.Add(1) }

func Intern(s string) StringID {
	if len(s) == 0 {
		return EmptyStringID
	}
	GlobalIntern.Add(1)
	mu.Lock()
	defer mu.Unlock()
	if id, ok := strMap[s]; ok {
		return id
	}
	id := StringID(len(strs))
	strs = append(strs, s)
	strMap[s] = id
	return id
}

func ResolveString(id StringID) string {
	if id == 0 {
		return "\xff\x00j2go_null"
	}
	if id == EmptyStringID {
		return ""
	}
	mu.RLock()
	defer mu.RUnlock()
	return strs[id]
}

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

func (n Node) SetOrigNameID(v StringID) {
	MutateRec(n, func(r *NodeRec) { r.OrigName = v })
}

func (n Node) OrigNameID() StringID {
	return GetRec(n).OrigName
}

type NodeAllocator struct {
	inited     bool
	tmpl       Node
	tmplSF     uint32
	tmplProps  uint32
	cache      [256]struct {
		s  string
		id StringID
	}
}

func (a *NodeAllocator) InitLocal() {
	a.inited = true
	a.tmpl = 0
	a.tmplSF = 0
	a.tmplProps = 0
}

func (a *NodeAllocator) Flush() {
	FlushCalls.Add(1)
	*a = NodeAllocator{}
}

func (a *NodeAllocator) SetTemplate(tmpl Node) {
	a.tmpl = tmpl
	if tmpl == 0 {
		a.tmplSF = 0
		a.tmplProps = 0
		return
	}
	r := GetRec(tmpl)
	a.tmplSF = r.SourceFile
	a.tmplProps = r.PropFlags
}

func (a *NodeAllocator) Intern(s string) StringID {
	if len(s) == 0 {
		return EmptyStringID
	}
	idx := (uint32(len(s))*2654435761 ^ uint32(s[0])) & 255
	e := &a.cache[idx]
	if e.id != 0 && e.s == s {
		return e.id
	}
	id := Intern(s)
	e.s = s
	e.id = id
	return id
}

func (a *NodeAllocator) stamp(r *NodeRec, tmpl Node) {
	if a.tmpl == tmpl {
		r.SourceFile = a.tmplSF
		r.PropFlags = a.tmplProps
	} else if tmpl != 0 {
		tr := GetRec(tmpl)
		r.SourceFile = tr.SourceFile
		r.PropFlags = tr.PropFlags
	}
}

func (a *NodeAllocator) AllocNode1(tok uint8, tmpl Node) Node {
	Alloc1Calls.Add(1)
	r := NodeRec{Token: tok}
	a.stamp(&r, tmpl)
	return AllocRaw(r)
}

func (a *NodeAllocator) AllocStringNode(tok uint8, strID StringID, tmpl Node) Node {
	AllocStrCall.Add(1)
	r := NodeRec{Token: tok, KindFlags: 3, Payload: uint32(strID)}
	a.stamp(&r, tmpl)
	return AllocRaw(r)
}

func (a *NodeAllocator) AllocNumberNode(tok uint8, val float64, tmpl Node) Node {
	AllocNumCall.Add(1)
	_ = math.Float64bits(val)
	r := NodeRec{Token: tok, KindFlags: 1, DoubleVal: val}
	a.stamp(&r, tmpl)
	return AllocRaw(r)
}
`)

	// 2. Minimal gen/j2gort with JIterCursor, OrderedMap, OrderedMapValCursor, PooledSlab, Deque, NPE/IllegalState.
	mustWrite("gen/j2gort/j2gort.go", `package j2gort

import "sync/atomic"

const NullString = "\xff\x00j2go_null"

type NPE struct{}
func (NPE) Error() string { return "NullPointerException" }

type IllegalState struct{ Msg string }
func (e IllegalState) Error() string { return "IllegalStateException: " + e.Msg }

func ThrowNPE()          { panic(NPE{}) }
func CheckState(b bool)  { if !b { panic(IllegalState{}) } }
func CheckNotNull[T any](v T) T {
	if any(v) == nil {
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
func AsIface[T any](v any) T {
	if v == nil {
		var z T
		return z
	}
	return v.(T)
}

type Ext_java_lang_Iterable interface {
	IterItems_() []any
}
type Ext_java_util_Set interface {
	Ext_java_lang_Iterable
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

type OrderedMap[K comparable, V any] struct {
	keys []K
	vals []V
	idx  map[K]int
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{idx: make(map[K]int)}
}
func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	m := NewOrderedMap[K, V]()
	return &m
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
func (d *Deque[T]) AddFirst(v T)  { d.s = append(d.s, v) }
func (d *Deque[T]) RemoveFirst()  { d.s = d.s[:len(d.s)-1] }
func (d *Deque[T]) PeekFirst() T  { return d.s[len(d.s)-1] }

var (
	SlabAllocCount   atomic.Int64
	SlabReleaseCount atomic.Int64
)

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

	// 3. Minimal gen/rhino using shims/nodestore.
	mustWrite("gen/rhino/rhino.go", `package rhino

import (
	"math"

	"example.com/w24ctest/gen/j2gort"
	"example.com/w24ctest/shims/nodestore"
)

type Token uint8

const (
	TokSCRIPT   Token = 1
	TokFUNCTION Token = 2
	TokBLOCK    Token = 3
	TokNAME     Token = 4
	TokNUMBER   Token = 5
	TokSTRING   Token = 6
	TokTRUE     Token = 7
)

func StaticGet_Token_SCRIPT() Token    { return TokSCRIPT }
func StaticGet_Token_FUNCTION() Token  { return TokFUNCTION }
func StaticGet_Token_BLOCK() Token     { return TokBLOCK }
func StaticGet_Token_NAME() Token      { return TokNAME }
func StaticGet_Token_NUMBER() Token    { return TokNUMBER }
func StaticGet_Token_STRINGLIT() Token { return TokSTRING }
func StaticGet_Token_TRUE() Token      { return TokTRUE }

type Node nodestore.Node

func (n Node) IsJavaNil_() bool     { return n == 0 }
func (n Node) GetToken() Token      { return Token(nodestore.GetRec(nodestore.Node(n)).Token + 1) }
func (n Node) IsScript() bool       { return n != 0 && n.GetToken() == TokSCRIPT }
func (n Node) IsFunction() bool     { return n != 0 && n.GetToken() == TokFUNCTION }
func (n Node) IsBlock() bool        { return n != 0 && n.GetToken() == TokBLOCK }
func (n Node) IsName() bool         { return n != 0 && n.GetToken() == TokNAME }
func (n Node) IsNumber() bool       { return n != 0 && n.GetToken() == TokNUMBER }
func (n Node) IsStringLit() bool    { return n != 0 && n.GetToken() == TokSTRING }
func (n Node) GetFirstChild() Node  { return Node(nodestore.GetRec(nodestore.Node(n)).First) }
func (n Node) GetNext() Node        { return Node(nodestore.GetRec(nodestore.Node(n)).Next) }
func (n Node) GetParent() Node      { return Node(nodestore.GetRec(nodestore.Node(n)).Parent) }
func (n Node) FirstNextTok_() (Node, Node, uint8) {
	r := nodestore.GetRec(nodestore.Node(n))
	return Node(r.First), Node(r.Next), r.Token
}
func (n Node) GetString() string {
	return nodestore.ResolveString(nodestore.StringID(nodestore.GetRec(nodestore.Node(n)).Payload))
}
func (n Node) GetOriginalName() string {
	return nodestore.ResolveString(nodestore.Node(n).OrigNameID())
}
func (n Node) GetDouble() float64   { return nodestore.GetRec(nodestore.Node(n)).DoubleVal }
func (n Node) GetSourceFileID() uint32 { return nodestore.GetRec(nodestore.Node(n)).SourceFile }
func (n Node) GetIsInClosureUnawareSubtree() bool {
	return nodestore.GetRec(nodestore.Node(n)).PropFlags&1 != 0
}
func (n Node) SetIsInClosureUnawareSubtree(v bool) {
	nodestore.MutateRec(nodestore.Node(n), func(r *nodestore.NodeRec) {
		if v {
			r.PropFlags |= 1
		} else {
			r.PropFlags &^= 1
		}
	})
}
func (n Node) SetSourceFileID(id uint32) {
	nodestore.MutateRec(nodestore.Node(n), func(r *nodestore.NodeRec) { r.SourceFile = id })
}
func (n Node) ClonePropsFrom(other Node) Node {
	or := nodestore.GetRec(nodestore.Node(other))
	nodestore.MutateRec(nodestore.Node(n), func(r *nodestore.NodeRec) {
		r.SourceFile = or.SourceFile
		r.PropFlags = or.PropFlags
	})
	return n
}
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
	if tok == 0 {
		j2gort.ThrowNPE()
	}
	return Node(nodestore.AllocRaw(nodestore.NodeRec{Token: uint8(tok) - 1}))
}

func NewNode_StringNode_2(tok Token, s string) Node {
	if tok == 0 || s == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	return Node(nodestore.AllocRaw(nodestore.NodeRec{Token: uint8(tok) - 1, KindFlags: 3, Payload: uint32(nodestore.Intern(s))}))
}

func IR_Number(d float64) Node {
	if math.IsNaN(d) || d < 0 || (d == 0 && math.Signbit(d)) {
		panic(j2gort.IllegalState{Msg: "negative or NaN number"})
	}
	return Node(nodestore.AllocRaw(nodestore.NodeRec{Token: uint8(TokNUMBER) - 1, KindFlags: 1, DoubleVal: d}))
}

type RhinoStringPool_LazyInternedStringList struct {
	pool []string
}

func NewRhinoStringPool_LazyInternedStringList(items []string) *RhinoStringPool_LazyInternedStringList {
	return &RhinoStringPool_LazyInternedStringList{pool: items}
}
func (l *RhinoStringPool_LazyInternedStringList) Size() int32 { return int32(len(l.pool)) }
func (l *RhinoStringPool_LazyInternedStringList) Get(off int32) string {
	return l.pool[off]
}

func Node_NewString_3(tok Token, sp *RhinoStringPool_LazyInternedStringList, off int32) Node {
	if tok == 0 {
		j2gort.ThrowNPE()
	}
	var id nodestore.StringID
	if s := sp.Get(off); s != j2gort.NullString {
		id = nodestore.Intern(s)
	}
	return Node(nodestore.AllocRaw(nodestore.NodeRec{Token: uint8(tok) - 1, KindFlags: 3, Payload: uint32(id)}))
}

func (n Node) SetOriginalNameFromStringPool(sp *RhinoStringPool_LazyInternedStringList, off int32) {
	if n == 0 {
		j2gort.ThrowNPE()
	}
	s := sp.Get(off)
	if s == j2gort.NullString {
		nodestore.Node(n).SetOrigNameID(0)
	} else {
		nodestore.Node(n).SetOrigNameID(nodestore.Intern(s))
	}
}
`)

	// 4. Minimal gen/jscomp: Scope, TypedScope, SyntacticScopeCreator, NodeTraversal, callbacks, pass-boundary passes, StringPool, ScriptNodeDeserializer.
	mustWrite("gen/jscomp/scope.go", `package jscomp

import (
	"example.com/w24ctest/gen/j2gort"
	pkg_rhino "example.com/w24ctest/gen/rhino"
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

func (*Var) Implements_jscomp_AbstractVar() {}
func (v *Var) IsJavaNil_() bool             { return v == nil }
func (v *Var) GetName() string              { return v.name }
func (v *Var) GetNode() pkg_rhino.Node      { return v.node }
func (v *Var) GetNameNode() pkg_rhino.Node  { return v.node }
func (v *Var) GetScopeRoot() pkg_rhino.Node { return v.scope.rootNode }
func (v *Var) IsLocal() bool                { return v.scope != nil && v.scope.parent != nil }
func (v *Var) GetScope() AbstractScope      { return v.scope }

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
func (s *Scope) GetClosestContainerScope() AbstractScope { return s }
func (s *Scope) HasSameContainerScope(o AbstractScope) bool {
	return o != nil && s == o.Untyped()
}
func (s *Scope) IsFunctionBlockScope() bool { return s.rootNode.IsBlock() }
func (s *Scope) IsGlobal() bool             { return s.parent == nil }
func (s *Scope) GetVarCount() int32 {
	if s.vars != nil {
		return int32(len(s.vars.LiveValues().IterItems_()))
	}
	var c int32
	for _, v := range []*Var{s.var0, s.var1, s.var2, s.var3, s.var4, s.var5} {
		if v != nil {
			c++
		}
	}
	return c
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

func (this_ *Scope) Undeclare(v *Var) {
	if this_.vars != nil {
		return
	}
	switch v {
	case this_.var0:
		this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5 = this_.var1, this_.var2, this_.var3, this_.var4, this_.var5, nil
	case this_.var1:
		this_.var1, this_.var2, this_.var3, this_.var4, this_.var5 = this_.var2, this_.var3, this_.var4, this_.var5, nil
	}
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

var ScopeGetVarIterableCalls int64

func (this_ *Scope) GetVarIterable() j2gort.Ext_java_lang_Iterable {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	ScopeGetVarIterableCalls++
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

func (this_ *Scope) GetAllSymbols() j2gort.Ext_java_lang_Iterable {
	return this_.GetVarIterable()
}

func ToIface_AbstractScope(s *Scope) AbstractScope {
	if s == nil {
		return nil
	}
	return s
}
`)

	mustWrite("gen/jscomp/typed_scope.go", `package jscomp

import (
	"example.com/w24ctest/gen/j2gort"
	pkg_rhino "example.com/w24ctest/gen/rhino"
)

type TypedVar struct {
	name string
}

type TypedScope struct {
	AbstractScopeBase
	var0 *TypedVar
	var1 *TypedVar
	var2 *TypedVar
	var3 *TypedVar
	var4 *TypedVar
	var5 *TypedVar
	vars *j2gort.OrderedMap[string, *TypedVar]
}

func (s *TypedScope) IsJavaNil_() bool       { return s == nil }
func (s *TypedScope) Untyped() *Scope        { return nil }
func NewTypedScope(root pkg_rhino.Node) *TypedScope {
	s := &TypedScope{}
	s.rootNode = root
	return s
}
func (this_ *TypedScope) GetVarIterable() j2gort.Ext_java_lang_Iterable {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	if this_.vars != nil {
		return this_.vars.LiveValues()
	}
	var items []any
	for _, v := range []*TypedVar{this_.var0, this_.var1, this_.var2, this_.var3, this_.var4, this_.var5} {
		if v != nil {
			items = append(items, v)
		}
	}
	return &j2gort.IterSlice{Items: items}
}
`)

	mustWrite("gen/jscomp/node_traversal.go", `package jscomp

import (
	"example.com/w24ctest/gen/j2gort"
	pkg_rhino "example.com/w24ctest/gen/rhino"
)

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope(pkg_rhino.Node, AbstractScope) AbstractScope
}

type SyntacticScopeCreator_RedeclarationHandler interface{}
type defaultRedeclHandler struct{}

func StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER() SyntacticScopeCreator_RedeclarationHandler {
	return defaultRedeclHandler{}
}
func clinit_SyntacticScopeCreator() {}

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

func (this_ *NodeTraversal) GetScopeRoot() pkg_rhino.Node { return this_.scopeNodes[this_.scopeDepth-1] }
func (this_ *NodeTraversal) InGlobalScope() bool          { return this_.scopeDepth <= 1 }
func (this_ *NodeTraversal) InFunctionBlockScope() bool   { return this_.GetScopeRoot().IsBlock() }

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
`)

	mustWrite("gen/jscomp/passes.go", `package jscomp

import (
	"sync/atomic"

	"example.com/w24ctest/gen/j2gort"
	pkg_rhino "example.com/w24ctest/gen/rhino"
)

// 1A. MakeDeclaredNamesUnique with JIterCursor over GetVarIterable() AND a non-scope JIterCursor in same method!
type MakeDeclaredNamesUnique_Renamer interface {
	Implements_jscomp_MakeDeclaredNamesUnique_Renamer()
	AddDeclaredName(string, bool)
	CreateForChildScope() MakeDeclaredNamesUnique_Renamer
}

type testRenamer struct{ names []string }
func (*testRenamer) Implements_jscomp_MakeDeclaredNamesUnique_Renamer() {}
func (r *testRenamer) AddDeclaredName(n string, _ bool)                 { r.names = append(r.names, n) }
func (r *testRenamer) CreateForChildScope() MakeDeclaredNamesUnique_Renamer {
	return &testRenamer{}
}

type MakeDeclaredNamesUnique struct {
	NodeTraversal_AbstractScopedCallbackBase
	renamerStack j2gort.Deque[MakeDeclaredNamesUnique_Renamer]
	rootRenamer  MakeDeclaredNamesUnique_Renamer
	extraTags    j2gort.Ext_java_lang_Iterable
	SeenVars     []string
}

func NewMakeDeclaredNamesUnique() *MakeDeclaredNamesUnique {
	return &MakeDeclaredNamesUnique{
		rootRenamer: &testRenamer{},
		extraTags:   &j2gort.IterSlice{Items: []any{"tagA"}},
	}
}

func (this_ *MakeDeclaredNamesUnique) EnterScope(t *NodeTraversal) {
	// Non-scope JIterCursor in the same block must NOT confuse ephemeral scope cursor tracking!
	curTags := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.extraTags))
	for curTags.HasNext() {
		_ = curTags.Next()
	}
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
	this_.renamerStack.Pop()
}

func (this_ *MakeDeclaredNamesUnique) priv_MakeDeclaredNamesUnique_findDeclaredNames_1(t *NodeTraversal, renamer MakeDeclaredNamesUnique_Renamer) {
	cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(t.GetScope().GetVarIterable()))
	for cur_4.HasNext() {
		rawElem_3 := cur_4.Next()
		v := j2gort.CastPtr[Var](rawElem_3)
		_ = v
		renamer.AddDeclaredName(v.GetName(), false)
		this_.SeenVars = append(this_.SeenVars, v.GetName())
	}
}

// Negative soundness control 1: ContextualRenameInverter stores *Var!
type MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo struct{ var_ *Var }
type MakeDeclaredNamesUnique_ContextualRenameInverter struct {
	NodeTraversal_AbstractScopedCallbackBase
	saved []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) EnterScope(t *NodeTraversal) {
	cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(t.GetScope().GetVarIterable()))
	for cur_6.HasNext() {
		rawElem_5 := cur_6.Next()
		var_ := j2gort.CastPtr[Var](rawElem_5)
		this_.saved = append(this_.saved, &MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo{var_: var_})
	}
}
func (*MakeDeclaredNamesUnique_ContextualRenameInverter) ExitScope(_ *NodeTraversal) {}

// Negative soundness control 2: GlobalNamespace_BuildGlobalNamespace stores *Scope!
type GlobalNamespace_Ref struct{ scope *Scope }
type GlobalNamespace_BuildGlobalNamespace struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	refs []*GlobalNamespace_Ref
}
func (b *GlobalNamespace_BuildGlobalNamespace) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		b.refs = append(b.refs, &GlobalNamespace_Ref{scope: t.GetScope()})
	}
}

// 1B. PolyfillUsageFinder_Traverser (depth-4 helper chain + ephemVar helper forwarding!)
type GuardedCallback interface {
	NodeTraversal_Callback
	VisitGuarded(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}
type GuardedCallbackBase struct {
	self_       GuardedCallback
	guardChecks j2gort.Ext_java_lang_Iterable
}
func (*GuardedCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (g *GuardedCallbackBase) IsJavaNil_() bool                        { return g == nil }
func (g *GuardedCallbackBase) BindSelf_(s any)                         { g.self_ = s.(GuardedCallback) }
func (this_ *GuardedCallbackBase) ShouldTraverse_3(t *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	_ = t
	if this_.guardChecks != nil {
		cur := j2gort.NewIterCursor(this_.guardChecks)
		for cur.HasNext() {
			_ = cur.Next()
		}
	}
	return true
}
func (this_ *GuardedCallbackBase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	this_.self_.VisitGuarded(t, n, p)
}

type PolyfillUsageFinder struct {
	compiler AbstractCompiler
	Hits     int
}
func PolyfillHelperIsLocal(v *Var) bool {
	return v != nil && v.IsLocal() && v.GetScopeRoot() != 0
}
func (p *PolyfillUsageFinder) priv_PolyfillUsageFinder_isPolyfillVisibleInScope_4(t *NodeTraversal, name string) bool {
	s := t.GetScope()
	v := j2gort.CastPtr[Var](s.GetSlot(name))
	if v != nil && s.HasSameContainerScope(v.GetScope()) && s == v.GetScope().Untyped() {
		return !PolyfillHelperIsLocal(v)
	}
	return v == nil
}
func (p *PolyfillUsageFinder) priv_PolyfillUsageFinder_maybeCreate_3(t *NodeTraversal, n pkg_rhino.Node) {
	if p.priv_PolyfillUsageFinder_isPolyfillVisibleInScope_4(t, n.GetString()) {
		p.Hits++
	}
}

type PolyfillUsageFinder_Traverser struct {
	GuardedCallbackBase
	outer_ *PolyfillUsageFinder
}
func NewPolyfillUsageFinder_Traverser(o *PolyfillUsageFinder) *PolyfillUsageFinder_Traverser {
	tr := &PolyfillUsageFinder_Traverser{outer_: o}
	tr.BindSelf_(tr)
	return tr
}
func (tr *PolyfillUsageFinder_Traverser) VisitGuarded(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	tr.priv_Traverser_visitGetPropChain_2(t, n)
}
func (tr *PolyfillUsageFinder_Traverser) priv_Traverser_visitGetPropChain_2(t *NodeTraversal, n pkg_rhino.Node) {
	if n.IsName() {
		tr.outer_.priv_PolyfillUsageFinder_maybeCreate_3(t, n)
	}
}

// 1C. FlowSensitiveInlineVariables with SideEffectPredicate + pre-lowered if-guarded cfg nil defer.
type MustBeReachingVariableDef struct{}
type MaybeReachingVariableUse struct{}
type FlowSensitiveInlineVariables_SideEffectPredicate struct {
	outer_       *FlowSensitiveInlineVariables
	namesToCheck j2gort.Ext_java_util_Set
	idHash_      atomic.Int32
}
type FlowSensitiveInlineVariables struct {
	compiler            AbstractCompiler
	sideEffectPredicate *FlowSensitiveInlineVariables_SideEffectPredicate
	cfg                 ControlFlowGraph[pkg_rhino.Node]
	candidates          j2gort.Ext_java_util_Set
	reachingDef         *MustBeReachingVariableDef
	reachingUses        *MaybeReachingVariableUse
	EnteredScopes       int
}
func (*FlowSensitiveInlineVariables) Implements_jscomp_NodeTraversal_Callback() {}
func (f *FlowSensitiveInlineVariables) IsJavaNil_() bool                        { return f == nil }
func (*FlowSensitiveInlineVariables) ShouldTraverse_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) bool {
	return true
}
func (*FlowSensitiveInlineVariables) Visit_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {}
func (this_ *FlowSensitiveInlineVariables) EnterScope(t *NodeTraversal) {
	if t.InGlobalScope() || !t.GetScope().IsFunctionBlockScope() {
		return
	}
	defer func() {
		if this_.cfg != nil && !this_.cfg.IsJavaNil_() {
			if r_, ok_ := any(this_.cfg).(interface{ ReleaseSlabs_() }); ok_ {
				r_.ReleaseSlabs_()
			}
			this_.cfg = nil
		}
		this_.candidates = nil
		this_.reachingDef = nil
		this_.reachingUses = nil
	}()
	c := &cfgImpl{}
	_ = c.slab.Alloc()
	this_.cfg = c
	this_.reachingDef = &MustBeReachingVariableDef{}
	this_.reachingUses = &MaybeReachingVariableUse{}
	this_.EnteredScopes++
}
func (this_ *FlowSensitiveInlineVariables) ExitScope(_ *NodeTraversal) {
	if this_.cfg != nil && !this_.cfg.IsJavaNil_() {
		this_.cfg = nil
	}
}

// 1D. DeadAssignmentsElimination embedding NodeTraversal_AbstractCfgCallbackBase.
type LiveVariablesAnalysis struct{ count int32 }
type DeadAssignmentsElimination_BailoutInformation struct{ ok bool }
type DeadAssignmentsElimination struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler      AbstractCompiler
	liveness      *LiveVariablesAnalysis
	functionStack j2gort.Deque[*DeadAssignmentsElimination_BailoutInformation]
	Eliminated    int
}
func NewDeadAssignmentsElimination(c AbstractCompiler) *DeadAssignmentsElimination {
	d := &DeadAssignmentsElimination{compiler: c}
	d.BindSelf_(d)
	return d
}
func (this_ *DeadAssignmentsElimination) Visit_3(_ *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {}
func (this_ *DeadAssignmentsElimination) EnterScopeWithCfg(t *NodeTraversal) {
	if t.InFunctionBlockScope() {
		this_.functionStack.AddFirst(&DeadAssignmentsElimination_BailoutInformation{ok: true})
	}
}
func (this_ *DeadAssignmentsElimination) ExitScopeWithCfg(t *NodeTraversal) {
	if t.InFunctionBlockScope() {
		this_.priv_DeadAssignmentsElimination_eliminateDeadAssignments_1(t)
		this_.functionStack.RemoveFirst()
	}
}
func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_eliminateDeadAssignments_1(t *NodeTraversal) {
	sc := NewSyntacticScopeCreator_1(this_.compiler)
	_ = sc
	this_.liveness = &LiveVariablesAnalysis{count: t.GetScope().GetVarCount()}
	if this_.liveness.count >= 0 {
		this_.Eliminated++
	}
}

// 1E. Pass-boundary ScopeCreator release passes (callbacks retain *Scope/*Var across Traverse inside Process_2/3!):
type Es6RewriteBlockScopedDeclaration_LoopClosureTransformer struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *Es6RewriteBlockScopedDeclaration
	vars   []*Var
}
func NewEs6RewriteBlockScopedDeclaration_LoopClosureTransformer(o *Es6RewriteBlockScopedDeclaration) *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer {
	return &Es6RewriteBlockScopedDeclaration_LoopClosureTransformer{outer_: o}
}
func (tr *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		if v := j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString())); v != nil {
			tr.vars = append(tr.vars, v)
		}
	}
}
func (tr *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) transform() {
	for _, v := range tr.vars {
		if v.GetName() == "" {
			panic("corrupted var before transform finished!")
		}
		tr.outer_.TransformedVars++
	}
}
type Es6RewriteBlockScopedDeclaration struct {
	compiler        AbstractCompiler
	TransformedVars int
}
func (this_ *Es6RewriteBlockScopedDeclaration) Process_2(_ pkg_rhino.Node, root pkg_rhino.Node) {
	transformer := NewEs6RewriteBlockScopedDeclaration_LoopClosureTransformer(this_)
	NodeTraversal_Traverse(this_.compiler, root, transformer)
	transformer.transform()
}

type InlineFunctions_FunctionState struct{ refScope *Scope }
type InlineFunctions_FindCandidatesReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *InlineFunctions
}
func NewInlineFunctions_FindCandidatesReferences(o *InlineFunctions) *InlineFunctions_FindCandidatesReferences {
	return &InlineFunctions_FindCandidatesReferences{outer_: o}
}
func (f *InlineFunctions_FindCandidatesReferences) Visit_3(t *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() {
		f.outer_.fns.Put(n.GetString(), &InlineFunctions_FunctionState{refScope: t.GetScope()})
	}
}
type InlineFunctions struct {
	fns      j2gort.OrderedMap[string, *InlineFunctions_FunctionState]
	anonFns  j2gort.OrderedMap[pkg_rhino.Node, string]
	compiler AbstractCompiler
	Inlined  int
}
func NewInlineFunctions(c AbstractCompiler) *InlineFunctions {
	return &InlineFunctions{
		compiler: c,
		fns:      j2gort.NewOrderedMap[string, *InlineFunctions_FunctionState](),
		anonFns:  j2gort.NewOrderedMap[pkg_rhino.Node, string](),
	}
}
func (this_ *InlineFunctions) Process_2(_ pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_Traverse(this_.compiler, root, NewInlineFunctions_FindCandidatesReferences(this_))
	cur := this_.fns.IterValues()
	for cur.HasNext() {
		st := cur.Next()
		if st.refScope.GetRootNode() == 0 {
			panic("scope released too early before InlineFunctions finished!")
		}
		this_.Inlined++
	}
}

type PureGraph interface {
	IsJavaNil_() bool
}
type pureGraphImpl struct{ slab j2gort.PooledSlab[int] }
func (g *pureGraphImpl) IsJavaNil_() bool { return g == nil }
func (g *pureGraphImpl) ReleaseSlabs_()   { g.slab.Release() }

type PureFunctionIdentifier_FunctionBodyAnalyzer struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer_ *PureFunctionIdentifier
}
func NewPureFunctionIdentifier_FunctionBodyAnalyzer(o *PureFunctionIdentifier) *PureFunctionIdentifier_FunctionBodyAnalyzer {
	return &PureFunctionIdentifier_FunctionBodyAnalyzer{outer_: o}
}
func (a *PureFunctionIdentifier_FunctionBodyAnalyzer) EnterScope(t *NodeTraversal) {
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(t.GetScope().GetVarIterable()))
	for cur.HasNext() {
		raw := cur.Next()
		v := j2gort.CastPtr[Var](raw)
		a.outer_.summariesByName.Put(v.GetName(), t.GetScope())
	}
}
func (*PureFunctionIdentifier_FunctionBodyAnalyzer) ExitScope(_ *NodeTraversal) {}

type PureFunctionIdentifier struct {
	compiler                             AbstractCompiler
	reverseCallGraph                     PureGraph
	summariesByName                      j2gort.OrderedMap[string, *Scope]
	summariesForAllNamesOfFunctionByNode j2gort.OrderedMap[pkg_rhino.Node, int]
	AnalyzedVars                         int
}
func NewPureFunctionIdentifier(c AbstractCompiler) *PureFunctionIdentifier {
	pg := &pureGraphImpl{}
	_ = pg.slab.Alloc()
	return &PureFunctionIdentifier{
		compiler:                             c,
		reverseCallGraph:                     pg,
		summariesByName:                      j2gort.NewOrderedMap[string, *Scope](),
		summariesForAllNamesOfFunctionByNode: j2gort.NewOrderedMap[pkg_rhino.Node, int](),
	}
}
func (this_ *PureFunctionIdentifier) Process_3(_ pkg_rhino.Node, root pkg_rhino.Node, _ *OptimizeCalls_ReferenceMap) {
	NodeTraversal_Traverse(this_.compiler, root, NewPureFunctionIdentifier_FunctionBodyAnalyzer(this_))
	cur := this_.summariesByName.IterValues()
	for cur.HasNext() {
		sc := cur.Next()
		if sc.GetRootNode() == 0 {
			panic("scope released before propagateSideEffects!")
		}
		this_.AnalyzedVars++
	}
}

type OptimizeCalls_ReferenceMap struct{ globalScope *Scope }
func (r *OptimizeCalls_ReferenceMap) GetGlobalScope() *Scope { return r.globalScope }
type OptimizeCalls_CallGraphCompilerPass interface {
	Process_3(pkg_rhino.Node, pkg_rhino.Node, *OptimizeCalls_ReferenceMap)
}
type OptimizeCalls_ReferenceMapBuildingCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *OptimizeCalls
	refs   *OptimizeCalls_ReferenceMap
}
func NewOptimizeCalls_ReferenceMapBuildingCallback(o *OptimizeCalls, r *OptimizeCalls_ReferenceMap) *OptimizeCalls_ReferenceMapBuildingCallback {
	return &OptimizeCalls_ReferenceMapBuildingCallback{outer_: o, refs: r}
}
func (cb *OptimizeCalls_ReferenceMapBuildingCallback) Visit_3(t *NodeTraversal, _ pkg_rhino.Node, _ pkg_rhino.Node) {
	if t.InGlobalScope() {
		cb.refs.globalScope = t.GetScope()
	}
}
type OptimizeParameters struct {
	globalScope *Scope
	CheckedVars int32
}
func (this_ *OptimizeParameters) Process_3(_ pkg_rhino.Node, _ pkg_rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	this_.globalScope = refMap.GetGlobalScope()
	this_.CheckedVars = this_.globalScope.GetVarCount()
}
type OptimizeCalls struct {
	compiler AbstractCompiler
	passes   []OptimizeCalls_CallGraphCompilerPass
}
func (this_ *OptimizeCalls) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	references := &OptimizeCalls_ReferenceMap{}
	NodeTraversal_TraverseRoots(this_.compiler, NewOptimizeCalls_ReferenceMapBuildingCallback(this_, references), externs, root)
	for _, p := range this_.passes {
		p.Process_3(externs, root, references)
	}
}
`)

	mustWrite("gen/jscomp/string_pool.go", `package jscomp

import (
	"example.com/w24ctest/gen/j2gort"
	pkg_rhino "example.com/w24ctest/gen/rhino"
)

type StringPool struct {
	pool      *pkg_rhino.RhinoStringPool_LazyInternedStringList
	maxLength int32
}

func NewStringPool(maxLength int32, items []string) *StringPool {
	sp := &StringPool{}
	InitStringPool(sp, maxLength, items)
	return sp
}

func InitStringPool(this_ *StringPool, maxLength int32, items []string) {
	_ = j2gort.NullString
	this_.maxLength = maxLength
	this_.pool = pkg_rhino.NewRhinoStringPool_LazyInternedStringList(items)
}

func (this_ *StringPool) Get(offset int32) string {
	return this_.pool.Get(offset)
}

func (this_ *StringPool) GetInternedStrings() *pkg_rhino.RhinoStringPool_LazyInternedStringList {
	return this_.pool
}
`)

	mustWrite("gen/jscomp/script_node_deserializer.go", `package jscomp

import (
	"example.com/w24ctest/gen/j2gort"
	pkg_rhino "example.com/w24ctest/gen/rhino"
)

type SourceFile struct{ ID uint32 }

type SerializedAstNode struct {
	Kind       int32
	StrPtr     int32
	OrigPtr    int32
	DoubleVal  float64
	SourceFile int32
	Shadow     bool
	Children   []*SerializedAstNode
}

func (a *SerializedAstNode) GetStringValuePointer() int32  { return a.StrPtr }
func (a *SerializedAstNode) GetOriginalNamePointer() int32 { return a.OrigPtr }
func (a *SerializedAstNode) GetDoubleValue() float64       { return a.DoubleVal }
func (a *SerializedAstNode) GetSourceFile() int32          { return a.SourceFile }

type ScriptNodeDeserializer struct {
	sourceFile *SourceFile
	rootProto  *SerializedAstNode
	stringPool *StringPool
	filePool   []*SourceFile
}

func NewScriptNodeDeserializer(sf *SourceFile, root *SerializedAstNode, sp *StringPool, fp []*SourceFile) *ScriptNodeDeserializer {
	return &ScriptNodeDeserializer{sourceFile: sf, rootProto: root, stringPool: sp, filePool: fp}
}

func (this_ *ScriptNodeDeserializer) DeserializeNew() pkg_rhino.Node {
	return NewScriptNodeDeserializer_Runner(this_).Run()
}

func (this_ *ScriptNodeDeserializer) createTemplate(sf *SourceFile) pkg_rhino.Node {
	t := pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT())
	if sf != nil {
		t.SetSourceFileID(sf.ID)
	}
	return t
}

func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_setOriginalNameIfPresent_1(astNode *SerializedAstNode, n pkg_rhino.Node) {
	if astNode.GetOriginalNamePointer() != 0 {
		n.SetOriginalNameFromStringPool(this_.stringPool.GetInternedStrings(), astNode.GetOriginalNamePointer())
	}
}

func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_stringNode_2(token pkg_rhino.Token, n *SerializedAstNode, templateNode pkg_rhino.Node) pkg_rhino.Node {
	str := pkg_rhino.Node_NewString_3(token, this_.stringPool.GetInternedStrings(), n.GetStringValuePointer())
	str.ClonePropsFrom(templateNode)
	return str
}

func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_newNodeFromTemplate_3(token pkg_rhino.Token, templateNode pkg_rhino.Node) pkg_rhino.Node {
	return pkg_rhino.NewNode_1(token).ClonePropsFrom(templateNode)
}

func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_deserializeSingleNode_4(n *SerializedAstNode, templateNode pkg_rhino.Node) pkg_rhino.Node {
	var numberNode pkg_rhino.Node
	_ = numberNode
	switch n.Kind {
	case 107:
		return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_3(pkg_rhino.StaticGet_Token_SCRIPT(), templateNode)
	case 1:
		numberNode = pkg_rhino.IR_Number(n.GetDoubleValue())
		numberNode.ClonePropsFrom(templateNode)
		return numberNode
	case 2:
		return this_.priv_ScriptNodeDeserializer_stringNode_2(pkg_rhino.StaticGet_Token_STRINGLIT(), n, templateNode)
	case 12:
		return this_.priv_ScriptNodeDeserializer_stringNode_2(pkg_rhino.StaticGet_Token_NAME(), n, templateNode)
	case 4:
		return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_3(pkg_rhino.StaticGet_Token_TRUE(), templateNode)
	case 97:
		return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_3(pkg_rhino.StaticGet_Token_BLOCK(), templateNode)
	}
	j2gort.CheckState(false)
	return 0
}

type ScriptNodeDeserializer_Runner struct {
	outer_ *ScriptNodeDeserializer
}

func NewScriptNodeDeserializer_Runner(outer_ *ScriptNodeDeserializer) *ScriptNodeDeserializer_Runner {
	r := &ScriptNodeDeserializer_Runner{outer_: outer_}
	InitScriptNodeDeserializer_Runner(r)
	return r
}

func InitScriptNodeDeserializer_Runner(this_ *ScriptNodeDeserializer_Runner) {
	_ = this_
}

func (this_ *ScriptNodeDeserializer_Runner) Run() pkg_rhino.Node {
	tmpl := this_.outer_.createTemplate(this_.outer_.sourceFile)
	return this_.priv_Runner_visit_5(this_.outer_.rootProto, nil, nil, tmpl)
}

func (this_ *ScriptNodeDeserializer_Runner) priv_Runner_visit_5(astNode *SerializedAstNode, ctxA any, ctxB any, sourceFileTemplate pkg_rhino.Node) pkg_rhino.Node {
	if sourceFileTemplate == 0 || astNode.GetSourceFile() != 0 {
		sourceFileTemplate = this_.outer_.createTemplate(this_.outer_.filePool[astNode.GetSourceFile()-1])
	}
	var n pkg_rhino.Node = this_.outer_.priv_ScriptNodeDeserializer_deserializeSingleNode_4(astNode, sourceFileTemplate)
	if astNode.Shadow && len(astNode.Children) > 0 {
		shadowTmpl := pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT()).ClonePropsFrom(sourceFileTemplate)
		shadowTmpl.SetIsInClosureUnawareSubtree(true)
		shadowChild := this_.priv_Runner_visit_5(astNode.Children[0], ctxA, ctxB, shadowTmpl)
		this_.outer_.priv_ScriptNodeDeserializer_setOriginalNameIfPresent_1(astNode.Children[0], shadowChild)
		n.AddChildToBack(shadowChild)
		for i := 1; i < len(astNode.Children); i++ {
			ch := this_.priv_Runner_visit_5(astNode.Children[i], ctxA, ctxB, sourceFileTemplate)
			n.AddChildToBack(ch)
			this_.outer_.priv_ScriptNodeDeserializer_setOriginalNameIfPresent_1(astNode.Children[i], ch)
		}
		return n
	}
	for _, child := range astNode.Children {
		dc := this_.priv_Runner_visit_5(child, ctxA, ctxB, sourceFileTemplate)
		n.AddChildToBack(dc)
		this_.outer_.priv_ScriptNodeDeserializer_setOriginalNameIfPresent_1(child, dc)
	}
	return n
}
`)

	mustWrite("gen/jscomp/scope_iter_helpers.go", `package jscomp

import "example.com/w24ctest/gen/j2gort"

func IterateScopeVarsHelper(s *Scope, undeclareFirst bool) []string {
	var out []string
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(s.GetVarIterable()))
	for cur.HasNext() {
		raw := cur.Next()
		v := j2gort.CastPtr[Var](raw)
		if v.GetName() == "skip" {
			continue
		}
		if v.GetName() == "stop" {
			break
		}
		if undeclareFirst && len(out) == 0 {
			s.Undeclare(v)
		}
		out = append(out, v.GetName())
	}
	return out
}

func IterateTypedScopeVarsHelper(ts *TypedScope) int {
	n := 0
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(ts.GetVarIterable()))
	for cur.HasNext() {
		raw := cur.Next()
		tv := j2gort.CastPtr[TypedVar](raw)
		if tv.name != "" {
			n++
		}
	}
	return n
}

func SumScopeVarNameLengths(s *Scope) int {
	sum := 0
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(s.GetVarIterable()))
		for cur.HasNext() {
			raw := cur.Next()
			v := j2gort.CastPtr[Var](raw)
			sum += len(v.GetName())
		}
	}
	return sum
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
	if res1.ByPass["slab_alloc_lowering"] < 6 || res1.ByPass["node_traversal_lowering"] < 6 {
		t.Fatalf("expected substantial lowering sites in both passes, got %+v", res1.ByPass)
	}

	// Verify Run-2 idempotency (0 sites, 0 modified files)!
	res2, err := passes.RunOnPackageDir(jscompDir, passList, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 modified files on Run 2 (idempotency), got %+v", res2)
	}

	// Inspect lowered AST text onpasses.go, scope.go, string_pool.go, script_node_deserializer.go!
	passesCodeBytes, err := os.ReadFile(filepath.Join(jscompDir, "passes.go"))
	if err != nil {
		t.Fatal(err)
	}
	passesCode := string(passesCodeBytes)
	for _, want := range []string{
		"func (*MakeDeclaredNamesUnique) CanReleaseSyntacticScopes_()",
		"func (*PolyfillUsageFinder_Traverser) CanReleaseSyntacticScopes_()",
		"func (*FlowSensitiveInlineVariables) CanReleaseSyntacticScopes_()",
		"func (*DeadAssignmentsElimination) CanReleaseSyntacticScopes_()",
		"t.GetScope()).iterVars_()",
	} {
		if !strings.Contains(passesCode, want) {
			t.Fatalf("expected lowered passes.go to contain %q\n--- passes.go ---\n%s", want, passesCode)
		}
	}
	for _, forbidden := range []string{
		"func (*MakeDeclaredNamesUnique_ContextualRenameInverter) CanReleaseSyntacticScopes_()",
		"func (*GlobalNamespace_BuildGlobalNamespace) CanReleaseSyntacticScopes_()",
		"func (*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) CanReleaseSyntacticScopes_()",
		"func (*InlineFunctions_FindCandidatesReferences) CanReleaseSyntacticScopes_()",
		"func (*PureFunctionIdentifier_FunctionBodyAnalyzer) CanReleaseSyntacticScopes_()",
		"func (*OptimizeCalls_ReferenceMapBuildingCallback) CanReleaseSyntacticScopes_()",
	} {
		if strings.Contains(passesCode, forbidden) {
			t.Fatalf("unsound marker emitted: %q found in passes.go!", forbidden)
		}
	}

	sndBytes, err := os.ReadFile(filepath.Join(jscompDir, "script_node_deserializer.go"))
	if err != nil {
		t.Fatal(err)
	}
	sndCode := string(sndBytes)
	for _, want := range []string{
		"nodeAlloc_",
		"InitLocal()",
		"Flush()",
		"SetTemplate(",
		"AllocNode1(",
		"AllocStringNode(",
		"AllocNumberNode(",
		"SetOrigNameID(",
		"internAt_(",
	} {
		if !strings.Contains(sndCode, want) {
			t.Fatalf("expected lowered script_node_deserializer.go to contain %q\n--- snd ---\n%s", want, sndCode)
		}
	}

	// Write comprehensive runtime behavioral, zero-alloc, NPE/Exception parity, and -race test inside gen/jscomp!
	mustWrite("gen/jscomp/wave24c_rt_test.go", `package jscomp

import (
	"math"
	"sync"
	"testing"

	"example.com/w24ctest/gen/j2gort"
	pkg_rhino "example.com/w24ctest/gen/rhino"
	"example.com/w24ctest/shims/nodestore"
)

func buildSampleScript() (pkg_rhino.Node, pkg_rhino.Node) {
	ext := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
	ext.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "extSym"))

	script := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
	script.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "g1"))
	blk := pkg_rhino.NewNode_1(pkg_rhino.TokBLOCK)
	blk.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "locA"))
	blk.AddChildToBack(pkg_rhino.NewNode_StringNode_2(pkg_rhino.TokNAME, "locB"))
	script.AddChildToBack(blk)
	return ext, script
}

func TestWave24C_ScopeReleasesAndPassBoundaries(t *testing.T) {
	comp := &Compiler{}
	ext, script := buildSampleScript()

	// 1. Ephemeral / transient-clearing callbacks MUST release SyntacticScopeCreator slabs on Traverse exit!
	for _, tc := range []struct {
		name string
		run  func()
	}{
		{"MakeDeclaredNamesUnique", func() {
			cb := NewMakeDeclaredNamesUnique()
			NodeTraversal_Traverse(comp, script, cb)
			if len(cb.SeenVars) != 3 {
				t.Fatalf("MakeDeclaredNamesUnique SeenVars = %v, want 3 vars", cb.SeenVars)
			}
		}},
		{"PolyfillUsageFinder_Traverser", func() {
			puf := &PolyfillUsageFinder{compiler: comp}
			NodeTraversal_Traverse(comp, script, NewPolyfillUsageFinder_Traverser(puf))
		}},
		{"FlowSensitiveInlineVariables", func() {
			fsiv := &FlowSensitiveInlineVariables{compiler: comp, sideEffectPredicate: &FlowSensitiveInlineVariables_SideEffectPredicate{}}
			NodeTraversal_Traverse(comp, script, fsiv)
			if fsiv.EnteredScopes != 1 {
				t.Fatalf("FlowSensitiveInlineVariables EnteredScopes = %d, want 1", fsiv.EnteredScopes)
			}
		}},
		{"DeadAssignmentsElimination", func() {
			dae := NewDeadAssignmentsElimination(comp)
			NodeTraversal_Traverse(comp, script, dae)
			if dae.Eliminated != 1 {
				t.Fatalf("DeadAssignmentsElimination Eliminated = %d, want 1", dae.Eliminated)
			}
		}},
	} {
		before := j2gort.SlabReleaseCount.Load()
		tc.run()
		after := j2gort.SlabReleaseCount.Load()
		if after <= before {
			t.Fatalf("%s: expected SyntacticScopeCreator slabs to be released (before=%d, after=%d)", tc.name, before, after)
		}
	}

	// 2. ContextualRenameInverter & GlobalNamespace_BuildGlobalNamespace MUST retain live *Var/*Scope past Traverse!
	inv := &MakeDeclaredNamesUnique_ContextualRenameInverter{}
	NodeTraversal_Traverse(comp, script, inv)
	if len(inv.saved) != 3 || inv.saved[0].var_.GetName() != "g1" {
		t.Fatalf("ContextualRenameInverter retained vars corrupted: %+v", inv.saved)
	}
	gn := &GlobalNamespace_BuildGlobalNamespace{}
	NodeTraversal_Traverse(comp, script, gn)
	if len(gn.refs) == 0 || gn.refs[0].scope.GetRootNode() == 0 {
		t.Fatalf("GlobalNamespace retained scopes corrupted!")
	}

	// 3. Pass-boundary ScopeCreator slab releases for Es6RewriteBlockScopedDeclaration, InlineFunctions, PureFunctionIdentifier, OptimizeCalls!
	for _, tc := range []struct {
		name string
		run  func()
	}{
		{"Es6RewriteBlockScopedDeclaration", func() {
			p := &Es6RewriteBlockScopedDeclaration{compiler: comp}
			p.Process_2(ext, script)
			if p.TransformedVars == 0 {
				t.Fatalf("expected TransformedVars > 0")
			}
		}},
		{"InlineFunctions", func() {
			p := NewInlineFunctions(comp)
			p.Process_2(ext, script)
			if p.Inlined == 0 || !p.fns.IsEmpty() {
				t.Fatalf("expected Inlined > 0 and fns cleared on exit, got Inlined=%d isEmpty=%v", p.Inlined, p.fns.IsEmpty())
			}
		}},
		{"PureFunctionIdentifier", func() {
			p := NewPureFunctionIdentifier(comp)
			p.Process_3(ext, script, nil)
			if p.AnalyzedVars == 0 || p.reverseCallGraph != nil || !p.summariesByName.IsEmpty() {
				t.Fatalf("expected PureFunctionIdentifier cleanup, got AnalyzedVars=%d rcg=%v", p.AnalyzedVars, p.reverseCallGraph)
			}
		}},
		{"OptimizeCalls+OptimizeParameters", func() {
			op := &OptimizeParameters{}
			oc := &OptimizeCalls{compiler: comp, passes: []OptimizeCalls_CallGraphCompilerPass{op}}
			oc.Process_2(ext, script)
			if op.CheckedVars != 1 || op.globalScope != nil {
				t.Fatalf("expected OptimizeParameters CheckedVars=1 and globalScope nilled, got %d / %v", op.CheckedVars, op.globalScope)
			}
		}},
	} {
		before := j2gort.SlabReleaseCount.Load()
		tc.run()
		after := j2gort.SlabReleaseCount.Load()
		if after <= before {
			t.Fatalf("%s: expected pass-boundary ReleaseSlabs_() (before=%d, after=%d)", tc.name, before, after)
		}
	}
}

func TestWave24C_ScopeVarCursorZeroAllocAndParity(t *testing.T) {
	root := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
	ScopeGetVarIterableCalls = 0

	// Test 0..8 vars (inline slots var0..var5 AND OrderedMap >6 vars), continue/break, Undeclare mid-loop!
	for count := 0; count <= 8; count++ {
		sc := NewScope(nil, root)
		for i := 0; i < count; i++ {
			nm := string(rune('a' + i))
			sc.Declare(nm, root, NewVar(nm, root, sc))
		}
		got := IterateScopeVarsHelper(sc, false)
		if len(got) != count {
			t.Fatalf("count=%d: got %v", count, got)
		}
	}
	if ScopeGetVarIterableCalls != 0 {
		t.Fatalf("expected 0 calls to GetVarIterable() after iterVars_() lowering, got %d", ScopeGetVarIterableCalls)
	}

	// Undeclare during iteration MUST visit all originally-present inline vars (snapshot parity with ImmutableList_Of_*):
	scMut := NewScope(nil, root)
	scMut.Declare("v0", root, NewVar("v0", root, scMut))
	scMut.Declare("v1", root, NewVar("v1", root, scMut))
	scMut.Declare("v2", root, NewVar("v2", root, scMut))
	gotMut := IterateScopeVarsHelper(scMut, true)
	if len(gotMut) != 3 || gotMut[0] != "v0" || gotMut[1] != "v1" || gotMut[2] != "v2" {
		t.Fatalf("Undeclare mid-iteration broke snapshot order: got %v", gotMut)
	}

	// TypedScope cursor lowering:
	ts := NewTypedScope(root)
	ts.var0 = &TypedVar{name: "t0"}
	ts.var1 = &TypedVar{name: "t1"}
	if n := IterateTypedScopeVarsHelper(ts); n != 2 {
		t.Fatalf("TypedScope iterVars_() = %d, want 2", n)
	}

	// Nil *Scope MUST panic with j2gort.NPE!
	didNPE := false
	func() {
		defer func() {
			if _, ok := recover().(j2gort.NPE); ok {
				didNPE = true
			}
		}()
		_ = IterateScopeVarsHelper(nil, false)
	}()
	if !didNPE {
		t.Fatalf("expected NPE when iterating nil *Scope")
	}

	// Zero-alloc check via testing.AllocsPerRun!
	scBench := NewScope(nil, root)
	for i := 0; i < 6; i++ {
		nm := string(rune('a' + i))
		scBench.Declare(nm, root, NewVar(nm, root, scBench))
	}
	allocs := testing.AllocsPerRun(100, func() {
		if sum := SumScopeVarNameLengths(scBench); sum != 6 {
			t.Fatal(sum)
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs/op for lowered Scope.GetVarIterable() loop, got %v", allocs)
	}
}

func TestWave24C_ScriptNodeDeserializerAndRace(t *testing.T) {
	sp := NewStringPool(16, []string{"", "alpha", "origAlpha", "beta", j2gort.NullString})
	sf1 := &SourceFile{ID: 101}
	sf2 := &SourceFile{ID: 202}
	fp := []*SourceFile{sf1, sf2}

	proto := &SerializedAstNode{
		Kind:   107,
		Shadow: true,
		Children: []*SerializedAstNode{
			{Kind: 12, StrPtr: 1, OrigPtr: 2},                        // shadow child (ClosureUnawareSubtree = true)
			{Kind: 1, DoubleVal: 0.0},                                // +0.0 fast path (restored sf1, ClosureUnawareSubtree = false!)
			{Kind: 1, DoubleVal: 42.5, SourceFile: 2},                // switched SourceFile sf2
			{Kind: 2, StrPtr: 3, OrigPtr: 4},                         // string "beta", NullString origName -> OrigNameID 0
			{Kind: 4},                                                // TRUE leaf
		},
	}

	nodestore.Alloc1Calls.Store(0)
	nodestore.AllocStrCall.Store(0)
	nodestore.AllocNumCall.Store(0)
	nodestore.FlushCalls.Store(0)
	nodestore.GlobalIntern.Store(0)

	snd := NewScriptNodeDeserializer(sf1, proto, sp, fp)
	root := snd.DeserializeNew()

	if nodestore.Alloc1Calls.Load() < 2 || nodestore.AllocStrCall.Load() < 2 || nodestore.AllocNumCall.Load() < 2 || nodestore.FlushCalls.Load() != 1 {
		t.Fatalf("expected NodeAllocator fast-path calls, got alloc1=%d str=%d num=%d flush=%d",
			nodestore.Alloc1Calls.Load(), nodestore.AllocStrCall.Load(), nodestore.AllocNumCall.Load(), nodestore.FlushCalls.Load())
	}

	c0 := root.GetFirstChild()
	if c0.GetString() != "alpha" || c0.GetOriginalName() != "origAlpha" || !c0.GetIsInClosureUnawareSubtree() || c0.GetSourceFileID() != 101 {
		t.Fatalf("shadow child mismatch: str=%q orig=%q shadow=%v sf=%d",
			c0.GetString(), c0.GetOriginalName(), c0.GetIsInClosureUnawareSubtree(), c0.GetSourceFileID())
	}
	c1 := c0.GetNext()
	if !c1.IsNumber() || c1.GetDouble() != 0.0 || c1.GetIsInClosureUnawareSubtree() || c1.GetSourceFileID() != 101 {
		t.Fatalf("sibling after shadow child failed template restore: shadow=%v sf=%d", c1.GetIsInClosureUnawareSubtree(), c1.GetSourceFileID())
	}
	c2 := c1.GetNext()
	if !c2.IsNumber() || c2.GetDouble() != 42.5 || c2.GetSourceFileID() != 202 {
		t.Fatalf("sf2 number child mismatch: val=%v sf=%d", c2.GetDouble(), c2.GetSourceFileID())
	}
	c3 := c2.GetNext()
	if !c3.IsStringLit() || c3.GetString() != "beta" || c3.GetOriginalName() != j2gort.NullString {
		t.Fatalf("c3 string/null-orig mismatch: str=%q orig=%q", c3.GetString(), c3.GetOriginalName())
	}

	// Second DeserializeNew() on same StringPool MUST hit StringPool.internAt_ atomic cache (0 GlobalIntern calls)!
	internBefore := nodestore.GlobalIntern.Load()
	_ = snd.DeserializeNew()
	if nodestore.GlobalIntern.Load() != internBefore {
		t.Fatalf("expected 0 global Intern calls on cached StringPool, got delta=%d", nodestore.GlobalIntern.Load()-internBefore)
	}

	// Epoch bump (simulating RewindBaseline) MUST safely fallback without stale ID crash!
	nodestore.BumpEpochForTest()
	rootAfterBump := snd.DeserializeNew()
	if rootAfterBump.GetFirstChild().GetString() != "alpha" {
		t.Fatalf("expected valid string after Epoch bump")
	}

	// Negative / -0.0 / NaN numbers MUST still panic with IllegalStateException just like IR_Number!
	for _, badNum := range []float64{-1.5, math.Copysign(0, -1), math.NaN()} {
		didPanic := false
		func() {
			defer func() {
				if _, ok := recover().(j2gort.IllegalState); ok {
					didPanic = true
				}
			}()
			badSnd := NewScriptNodeDeserializer(sf1, &SerializedAstNode{Kind: 1, DoubleVal: badNum}, sp, fp)
			_ = badSnd.DeserializeNew()
		}()
		if !didPanic {
			t.Fatalf("expected IllegalStateException for bad number %v", badNum)
		}
	}

	// Concurrent -race stress test on shared ScriptNodeDeserializer & StringPool!
	var wg sync.WaitGroup
	for g := 0; g < 6; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 25; i++ {
				r := snd.DeserializeNew()
				if r.GetFirstChild().GetString() != "alpha" {
					panic("concurrent DeserializeNew mismatch")
				}
			}
		}()
	}
	wg.Wait()
}
`)

	cmd := exec.Command("go", "test", "-timeout=30s", "-race", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race in wave24c workspace failed: %v\n%s", err, string(out))
	}
}

func TestWave24C_RealGenJscompCoverageAndIdempotency(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/rhino or gen/j2gort is awaiting regeneration; skipping real gen/jscomp Wave24C check")
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
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/w24crealgen", 1) +
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
	t.Logf("Wave 24C Real gen/jscomp Run 1 sites: total=%d modifiedFiles=%d byPass=%+v",
		res1.TotalSites, res1.FilesModified, res1.ByPass)

	readCopy := func(name string) string {
		b, rErr := os.ReadFile(filepath.Join(copyDir, name))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}

	// 1. Verify CanReleaseSyntacticScopes_ positive & negative soundness across real gen/jscomp!
	mdnu := readCopy("make_declared_names_unique.go")
	if !strings.Contains(mdnu, "func (*MakeDeclaredNamesUnique) CanReleaseSyntacticScopes_()") {
		t.Fatalf("expected MakeDeclaredNamesUnique to receive CanReleaseSyntacticScopes_() in real gen/jscomp")
	}
	if strings.Contains(mdnu, "func (*MakeDeclaredNamesUnique_ContextualRenameInverter) CanReleaseSyntacticScopes_()") {
		t.Fatalf("ContextualRenameInverter MUST NOT receive CanReleaseSyntacticScopes_()")
	}
	if !strings.Contains(mdnu, ".iterVars_()") {
		t.Fatalf("expected MakeDeclaredNamesUnique GetVarIterable() loops to lower to iterVars_()")
	}

	if !strings.Contains(readCopy("flow_sensitive_inline_variables.go"), "func (*FlowSensitiveInlineVariables) CanReleaseSyntacticScopes_()") {
		t.Fatalf("expected FlowSensitiveInlineVariables to receive CanReleaseSyntacticScopes_()")
	}
	if !strings.Contains(readCopy("dead_assignments_elimination.go"), "func (*DeadAssignmentsElimination) CanReleaseSyntacticScopes_()") {
		t.Fatalf("expected DeadAssignmentsElimination to receive CanReleaseSyntacticScopes_()")
	}
	if !strings.Contains(readCopy("polyfill_usage_finder.go"), "func (*PolyfillUsageFinder_Traverser) CanReleaseSyntacticScopes_()") {
		t.Fatalf("expected PolyfillUsageFinder_Traverser to receive CanReleaseSyntacticScopes_()")
	}
	if strings.Contains(readCopy("global_namespace.go"), "func (*GlobalNamespace_BuildGlobalNamespace) CanReleaseSyntacticScopes_()") {
		t.Fatalf("GlobalNamespace_BuildGlobalNamespace MUST NOT receive CanReleaseSyntacticScopes_()")
	}

	// 2. Verify pass-boundary SyntacticScopeCreator releases & soundness on the 4 scope-retaining passes!
	es6 := readCopy("es6_rewrite_block_scoped_declaration.go")
	if strings.Contains(es6, "func (*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) CanReleaseSyntacticScopes_()") ||
		!strings.Contains(es6, "defer sc_.ReleaseSlabs_()") {
		t.Fatalf("Es6RewriteBlockScopedDeclaration lowering mismatch")
	}
	inf := readCopy("inline_functions.go")
	if strings.Contains(inf, "func (*InlineFunctions_FindCandidatesReferences) CanReleaseSyntacticScopes_()") ||
		!strings.Contains(inf, "this_.fns.ReleaseToPool()") || !strings.Contains(inf, "sc_.ReleaseSlabs_()") {
		t.Fatalf("InlineFunctions lowering mismatch")
	}
	pfi := readCopy("pure_function_identifier.go")
	if strings.Contains(pfi, "func (*PureFunctionIdentifier_FunctionBodyAnalyzer) CanReleaseSyntacticScopes_()") ||
		!strings.Contains(pfi, "this_.reverseCallGraph = nil") || !strings.Contains(pfi, "sc_.ReleaseSlabs_()") ||
		!strings.Contains(pfi, ".iterVars_()") {
		t.Fatalf("PureFunctionIdentifier lowering mismatch")
	}
	oc := readCopy("optimize_calls.go")
	op := readCopy("optimize_parameters.go")
	if strings.Contains(oc, "func (*OptimizeCalls_ReferenceMapBuildingCallback) CanReleaseSyntacticScopes_()") ||
		!strings.Contains(oc, "references.globalScope = nil") || !strings.Contains(oc, "sc_.ReleaseSlabs_()") ||
		!strings.Contains(op, "this_.globalScope = nil") {
		t.Fatalf("OptimizeCalls / OptimizeParameters lowering mismatch")
	}

	// 3. Verify Scope & TypedScope iterVars_() synthesis and GetVarIterable() lowering across jscomp!
	if !strings.Contains(readCopy("scope.go"), "type scopeVarCursor_[V any] struct") ||
		!strings.Contains(readCopy("scope.go"), "func (this_ *Scope) iterVars_()") {
		t.Fatalf("expected scopeVarCursor_ and (*Scope).iterVars_() in scope.go")
	}
	hasTypedScopeIterVars := strings.Contains(readCopy("typed_scope.go"), "func (this_ *TypedScope) iterVars_()") &&
		strings.Contains(readCopy("typed_scope_creator.go"), ".iterVars_()")
	hasDirectResolveOwnVarTypes := strings.Contains(readCopy("typed_scope.go"), "func (this_ *TypedScope) ResolveOwnVarTypes(") &&
		strings.Contains(readCopy("typed_scope_creator.go"), ".ResolveOwnVarTypes(")
	if !hasTypedScopeIterVars && !hasDirectResolveOwnVarTypes {
		t.Fatalf("expected either (*TypedScope).iterVars_() or direct (*TypedScope).ResolveOwnVarTypes(...) in typed_scope.go and typed_scope_creator.go")
	}
	for _, fName := range []string{
		"inline_object_literals.go",
		"inline_variables.go",
		"scoped_aliases.go",
		"variable_reference_check.go",
		"rename_vars.go",
		"denormalize.go",
		"node_util.go",
		"infer_consts.go",
	} {
		if !strings.Contains(readCopy(fName), ".iterVars_()") {
			t.Fatalf("expected %s to lower GetVarIterable()/GetAllSymbols() to .iterVars_()", fName)
		}
	}

	// 4. Verify ScriptNodeDeserializer + StringPool batch NodeAllocator & StringID cache lowering!
	snd := readCopy("script_node_deserializer.go")
	for _, want := range []string{
		"nodeAlloc_",
		"InitLocal()",
		"Flush()",
		"SetTemplate(",
		"AllocNode1(",
		"AllocStringNode(",
		"AllocNumberNode(",
		"SetOrigNameID(",
		"internAt_(",
	} {
		if !strings.Contains(snd, want) {
			t.Fatalf("expected lowered real script_node_deserializer.go to contain %q", want)
		}
	}
	if !strings.Contains(readCopy("string_pool.go"), "internAt_(") {
		t.Fatalf("expected lowered real string_pool.go to contain internAt_")
	}

	// 5. Compile ALL 450+ lowered real jscomp files with `go build ./jscomp`!
	buildCmd := exec.Command("go", "build", "./jscomp")
	buildCmd.Dir = tmpDir
	if out, bErr := buildCmd.CombinedOutput(); bErr != nil {
		t.Fatalf("go build ./jscomp on lowered real gen/jscomp failed: %v\n%s", bErr, string(out))
	}

	// 6. Run 2 MUST be 100% idempotent (0 sites, 0 modified files)!
	res2, err := passes.RunOnPackageDir(copyDir, passList, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on lowered real gen/jscomp failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 modified files on Run 2 of real gen/jscomp, got %+v", res2)
	}
}

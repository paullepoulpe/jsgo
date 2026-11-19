package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func setupWave21CSynthWorkspace(t *testing.T) (tmpDir, jscompDir string) {
	t.Helper()
	tmpDir = t.TempDir()

	goMod := "module example.com/w21ctest\n\ngo 1.24.0\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	j2Dir := filepath.Join(tmpDir, "j2gort")
	rhinoDir := filepath.Join(tmpDir, "rhino")
	jscompDir = filepath.Join(tmpDir, "jscomp")
	for _, d := range []string{j2Dir, rhinoDir, jscompDir} {
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

var (
	PooledSlabChunkAllocs      atomic.Int64
	PooledSlabReleaseCalls     atomic.Int64
	PooledSlabReleaseWithCalls atomic.Int64
	OrderedMapAcquires         atomic.Int64
	OrderedMapReleases         atomic.Int64
	CustomOrderedMapAcquires   atomic.Int64
	CustomOrderedMapReleases   atomic.Int64
)

func ResetWave21CMetrics() {
	PooledSlabChunkAllocs.Store(0)
	PooledSlabReleaseCalls.Store(0)
	PooledSlabReleaseWithCalls.Store(0)
	OrderedMapAcquires.Store(0)
	OrderedMapReleases.Store(0)
	CustomOrderedMapAcquires.Store(0)
	CustomOrderedMapReleases.Store(0)
}

func CheckNotNullString(s string) string {
	if s == NullString {
		panic("NPE")
	}
	return s
}

func CheckArgument(b bool) {
	if !b {
		panic("IllegalArgumentException")
	}
}

func CheckState(b bool) {
	if !b {
		panic("IllegalStateException")
	}
}

func CheckNotNullIface[I interface{ IsJavaNil_() bool }](v I) I {
	if any(v) == nil || v.IsJavaNil_() {
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

func IfaceEquals[I interface{ IsJavaNil_() bool }](a, b I) bool {
	if any(a) == any(b) {
		return true
	}
	if any(a) == nil {
		return b.IsJavaNil_()
	}
	if any(b) == nil {
		return a.IsJavaNil_()
	}
	return a.IsJavaNil_() && b.IsJavaNil_()
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

type Slab[T any] struct {
	buf []T
}

func (s *Slab[T]) Alloc() *T {
	i := len(s.buf)
	if i == cap(s.buf) {
		c := make([]T, 1, 32)
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
			c = make([]T, 1, 32)
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

func (s *PooledSlab[T]) HasLiveChunks() bool {
	return len(s.buf) > 0 || cap(s.buf) > 0 || len(s.chunks) > 0
}

type OrderedMap[K comparable, V any] struct {
	keys []K
	vals []V
	idx  map[K]int
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func NewOrderedMapCap[K comparable, V any](capHint int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{
		keys: make([]K, 0, capHint),
		vals: make([]V, 0, capHint),
		idx:  make(map[K]int, capHint),
	}
}

func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}

func NewOrderedMapPtrCap[K comparable, V any](capHint int32) *OrderedMap[K, V] {
	m := NewOrderedMapCap[K, V](capHint)
	return &m
}

func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	OrderedMapAcquires.Add(1)
	return NewOrderedMapCap[K, V](8)
}

func AcquireOrderedMapCap[K comparable, V any](capHint int32) OrderedMap[K, V] {
	OrderedMapAcquires.Add(1)
	return NewOrderedMapCap[K, V](capHint)
}

func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	m := AcquireOrderedMap[K, V]()
	return &m
}

func AcquireOrderedMapPtrCap[K comparable, V any](capHint int32) *OrderedMap[K, V] {
	m := AcquireOrderedMapCap[K, V](capHint)
	return &m
}

func (m *OrderedMap[K, V]) Put(k K, v V) {
	if m.idx == nil {
		m.idx = make(map[K]int, 8)
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

func (m *OrderedMap[K, V]) IsEmpty() bool {
	return m == nil || len(m.keys) == 0
}

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
	if m == nil {
		return
	}
	OrderedMapReleases.Add(1)
	*m = OrderedMap[K, V]{}
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

type Ext_java_util_function_Function interface {
	IsJavaNil_() bool
	Apply(any) any
}

type fnLambda struct{ fn func(any) any }

func (l *fnLambda) IsJavaNil_() bool { return l == nil }
func (l *fnLambda) Apply(x any) any  { return l.fn(x) }

func NewLambda_Ext_java_util_function_Function(fn func(any) any) Ext_java_util_function_Function {
	return &fnLambda{fn: fn}
}

type CustomOrderedMap[K FastCustomKeyComparable, V any] struct {
	keys []K
	vals []V
}

func NewCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{}
}

func NewCustomOrderedMapCap[K FastCustomKeyComparable, V any](capHint int32) CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{
		keys: make([]K, 0, capHint),
		vals: make([]V, 0, capHint),
	}
}

func AcquireCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
	CustomOrderedMapAcquires.Add(1)
	return NewCustomOrderedMapCap[K, V](8)
}

func AcquireCustomOrderedMapCap[K FastCustomKeyComparable, V any](capHint int32) CustomOrderedMap[K, V] {
	CustomOrderedMapAcquires.Add(1)
	return NewCustomOrderedMapCap[K, V](capHint)
}

func AcquireCustomOrderedMapPtr[K FastCustomKeyComparable, V any]() *CustomOrderedMap[K, V] {
	m := AcquireCustomOrderedMap[K, V]()
	return &m
}

func AcquireCustomOrderedMapPtrCap[K FastCustomKeyComparable, V any](capHint int32) *CustomOrderedMap[K, V] {
	m := AcquireCustomOrderedMapCap[K, V](capHint)
	return &m
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

func (m *CustomOrderedMap[K, V]) KeysSlice() []K {
	if m == nil {
		return nil
	}
	return m.keys
}

func (m *CustomOrderedMap[K, V]) ReleaseToPool() {
	if m == nil {
		return
	}
	CustomOrderedMapReleases.Add(1)
	*m = CustomOrderedMap[K, V]{}
}

func CustomOrderedMapComputeIfAbsentFn[K FastCustomKeyComparable, V any](
	m *CustomOrderedMap[K, V],
	k K,
	zeroVal V,
	fn Ext_java_util_function_Function,
) V {
	for i, ek := range m.keys {
		if ek == k || (!ek.IsJavaNil_() && ek.Equals(k)) {
			return m.vals[i]
		}
	}
	raw := fn.Apply(k)
	if raw == nil {
		return zeroVal
	}
	v := raw.(V)
	m.Put(k, v)
	return v
}

type OrderedSet[T comparable] struct {
	items []T
}

func NewOrderedSet[T comparable]() OrderedSet[T] { return OrderedSet[T]{} }
func (s *OrderedSet[T]) Add(v T)                 { s.items = append(s.items, v) }
func (s *OrderedSet[T]) Clear()                  { clear(s.items); s.items = s.items[:0] }
func (s *OrderedSet[T]) Len() int                { return len(s.items) }

type Deque[T any] struct {
	items []T
}

func (d *Deque[T]) AddLast(v T) { d.items = append(d.items, v) }
func (d *Deque[T]) PeekLast() T {
	var z T
	if len(d.items) == 0 {
		return z
	}
	return d.items[len(d.items)-1]
}
func (d *Deque[T]) GetLast() T { return d.PeekLast() }
func (d *Deque[T]) RemoveLast() T {
	v := d.items[len(d.items)-1]
	var z T
	d.items[len(d.items)-1] = z
	d.items = d.items[:len(d.items)-1]
	return v
}
`
	if err := os.WriteFile(filepath.Join(j2Dir, "runtime.go"), []byte(j2Src), 0o644); err != nil {
		t.Fatal(err)
	}

	rhinoSrc := `package rhino

type Node uint32

func (n Node) IsJavaNil_() bool { return n == 0 }

type InputId struct{ Name string }
`
	if err := os.WriteFile(filepath.Join(rhinoDir, "node.go"), []byte(rhinoSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	return tmpDir, jscompDir
}

func TestWave21C_SyntacticScopeNodeTraversalAndReferenceCollectorPooling(t *testing.T) {
	tmpDir, jscompDir := setupWave21CSynthWorkspace(t)

	scopeAndVarSrc := `package jscomp

import (
	"fmt"

	"example.com/w21ctest/j2gort"
	pkg_rhino "example.com/w21ctest/rhino"
)

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
	GetDepth() int32
	GetParent() any
}

type AbstractVar interface {
	IsJavaNil_() bool
	GetName() string
}

func ToIface_AbstractScope(s *Scope) AbstractScope {
	if s == nil {
		return nil
	}
	return s
}

func ToIface_AbstractVar(v *Var) AbstractVar {
	if v == nil {
		return nil
	}
	return v
}

type AbstractScopeBase struct {
	self_    AbstractScope
	rootNode pkg_rhino.Node
}

func (b *AbstractScopeBase) BindSelf_(s any) { b.self_ = s.(AbstractScope) }

func InitAbstractScope(this_ *AbstractScopeBase, rootNode pkg_rhino.Node) {
	this_.rootNode = rootNode
}

type Var struct {
	name     string
	nameNode pkg_rhino.Node
	scope    *Scope
	index    int32
	input    *CompilerInput
}

func (v *Var) IsJavaNil_() bool  { return v == nil }
func (v *Var) GetName() string   { return v.name }
func (v *Var) GetScope() *Scope  { return v.scope }
func (v *Var) HashCode() int32   { return int32(len(v.name))*31 + v.index }
func (v *Var) Equals(o any) bool { x, ok := o.(*Var); return ok && x == v }

func (this_ *Var) initFields_Var() {}

func NewVar(name string, nameNode pkg_rhino.Node, scope *Scope, index int32, input *CompilerInput, implicit pkg_rhino.Node) *Var {
	this_ := &Var{}
	InitVar(this_, name, nameNode, scope, index, input, implicit)
	return this_
}

func InitVar(this_ *Var, name string, nameNode pkg_rhino.Node, scope *Scope, index int32, input *CompilerInput, implicit pkg_rhino.Node) {
	_ = implicit
	this_.initFields_Var()
	this_.name = name
	this_.nameNode = nameNode
	this_.scope = scope
	this_.index = index
	this_.input = input
}

type Scope struct {
	AbstractScopeBase
	parent *Scope
	depth  int32
	var0   *Var
	var1   *Var
	var2   *Var
	var3   *Var
	var4   *Var
	var5   *Var
	vars   *j2gort.OrderedMap[string, *Var]
}

func (this_ *Scope) BindSelf_(s any)               { this_.AbstractScopeBase.BindSelf_(s) }
func (this_ *Scope) IsJavaNil_() bool              { return this_ == nil }
func (this_ *Scope) GetRootNode_0() pkg_rhino.Node { return this_.rootNode }
func (this_ *Scope) GetDepth() int32               { return this_.depth }
func (this_ *Scope) GetParent() any                { return j2gort.PtrToAny(this_.parent) }
func (this_ *Scope) IsHoistScope() bool            { return true }
func (this_ *Scope) initFields_Scope()             {}

func NewScope_1(rootNode pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_1(this_, rootNode)
	return this_
}

func InitScope_1(this_ *Scope, rootNode pkg_rhino.Node) {
	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_Scope()
	this_.parent = nil
	this_.depth = 0
}

func NewScope_2(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_2(this_, parent, rootNode)
	return this_
}

func InitScope_2(this_ *Scope, parent *Scope, rootNode pkg_rhino.Node) {
	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_Scope()
	this_.parent = parent
	this_.depth = parent.depth + 1
}

func Scope_CreateGlobalScope(rootNode pkg_rhino.Node) *Scope {
	return NewScope_1(rootNode)
}

func Scope_CreateChildScope(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	return NewScope_2(parent, rootNode)
}

func (this_ *Scope) GetOwnDeclaredVar(name string) *Var {
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

func (this_ *Scope) PutOwnDeclaredVar(name string, var_ *Var) {
	if this_.vars != nil {
		this_.vars.Put(name, var_)
		return
	}
	if this_.var0 == nil || this_.var0.name == name {
		this_.var0 = var_
		return
	}
	if this_.var1 == nil || this_.var1.name == name {
		this_.var1 = var_
		return
	}
	if this_.var2 == nil || this_.var2.name == name {
		this_.var2 = var_
		return
	}
	if this_.var3 == nil || this_.var3.name == name {
		this_.var3 = var_
		return
	}
	if this_.var4 == nil || this_.var4.name == name {
		this_.var4 = var_
		return
	}
	if this_.var5 == nil || this_.var5.name == name {
		this_.var5 = var_
		return
	}
	this_.vars = j2gort.NewOrderedMapPtrCap[string, *Var](int32(16))
	this_.vars.Put(this_.var0.name, this_.var0)
	this_.vars.Put(this_.var1.name, this_.var1)
	this_.vars.Put(this_.var2.name, this_.var2)
	this_.vars.Put(this_.var3.name, this_.var3)
	this_.vars.Put(this_.var4.name, this_.var4)
	this_.vars.Put(this_.var5.name, this_.var5)
	this_.vars.Put(name, var_)
	this_.var0 = nil
	this_.var1 = nil
	this_.var2 = nil
	this_.var3 = nil
	this_.var4 = nil
	this_.var5 = nil
}

func (this_ *Scope) Declare(name string, nameNode pkg_rhino.Node, input *CompilerInput) *Var {
	v := NewVar(name, nameNode, this_, 0, input, 0)
	this_.PutOwnDeclaredVar(name, v)
	return v
}

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope_2_Node_AbstractScope(pkg_rhino.Node, AbstractScope) any
}

func ToIface_ScopeCreator(sc ScopeCreator) ScopeCreator { return sc }

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
	// Declare 8 vars so Scope.vars spillover (> 6 vars) is exercised on every created Scope!
	for i := 0; i < 8; i++ {
		sc.Declare(fmt.Sprintf("v_%d_%d", n, i), n, nil)
	}
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
	return j2gort.PtrToAny(this_.CreateScope_2_Node_Scope(n, p))
}

type MemoizedScopeCreator struct {
	scopesByScopeRoot j2gort.OrderedMap[pkg_rhino.Node, AbstractScope]
	delegate          ScopeCreator
}

func (m *MemoizedScopeCreator) IsJavaNil_() bool { return m == nil }

func (this_ *MemoizedScopeCreator) initFields_MemoizedScopeCreator() {
	this_.scopesByScopeRoot = j2gort.NewOrderedMap[pkg_rhino.Node, AbstractScope]()
}

func NewMemoizedScopeCreator(delegate ScopeCreator) *MemoizedScopeCreator {
	this_ := &MemoizedScopeCreator{}
	InitMemoizedScopeCreator(this_, delegate)
	return this_
}

func InitMemoizedScopeCreator(this_ *MemoizedScopeCreator, delegate ScopeCreator) {
	this_.initFields_MemoizedScopeCreator()
	this_.delegate = delegate
}

func (this_ *MemoizedScopeCreator) CreateScope_2_Node_AbstractScope(n pkg_rhino.Node, parent AbstractScope) any {
	s := this_.scopesByScopeRoot.Get(n)
	if s == nil || s.IsJavaNil_() {
		s = j2gort.AsIface[AbstractScope](this_.delegate.CreateScope_2_Node_AbstractScope(n, parent))
		this_.scopesByScopeRoot.Put(n, s)
	}
	return s
}
`

	nodeTraversalSrc := `package jscomp

import (
	"example.com/w21ctest/j2gort"
	pkg_rhino "example.com/w21ctest/rhino"
)

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

type NodeTraversal_AbstractPostOrderCallbackBase struct{}

func (b *NodeTraversal_AbstractPostOrderCallbackBase) IsJavaNil_() bool { return b == nil }
func (b *NodeTraversal_AbstractPostOrderCallbackBase) ShouldTraverse_3(t *NodeTraversal, n, p pkg_rhino.Node) bool {
	return true
}

type NodeTraversal struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCallback                                  NodeTraversal_ScopedCallback
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
	mayContainSyntheticBlocks                      bool
	currentNode                                    pkg_rhino.Node
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
	if n > 1 {
		this_.pushScope(n - 1)
		this_.priv_NodeTraversal_traverseBranch_3833(n-1, n)
		this_.popScope()
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

func NodeTraversal_TraverseRoots(compiler AbstractCompiler, cb NodeTraversal_Callback, externs pkg_rhino.Node, root pkg_rhino.Node) {
	clinit_NodeTraversal()
	NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).TraverseRoots(externs, root)
}

func NodeTraversal_TraverseAtScope(compiler AbstractCompiler, cb NodeTraversal_Callback, scope AbstractScope) {
	clinit_NodeTraversal()
	NodeTraversal_S_Builder().SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).TraverseAtScope(scope)
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

func (this_ *NodeTraversal_Builder) SetScopeCreator(x ScopeCreator) *NodeTraversal_Builder {
	this_.scopeCreator = x
	return this_
}

func (this_ *NodeTraversal_Builder) SetObeyDestructuringAndDefaultValueExecutionOrder(x bool) *NodeTraversal_Builder {
	this_.obeyDestructuringAndDefaultValueExecutionOrder = x
	return this_
}

func (this_ *NodeTraversal_Builder) Build() *NodeTraversal {
	return NewNodeTraversal(this_)
}

func (this_ *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverse_3817(root)
}

func (this_ *NodeTraversal_Builder) TraverseRoots(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverseRoots_3818(externs, root)
}

func (this_ *NodeTraversal_Builder) TraverseAtScope(scope AbstractScope) {
	this_.Build().priv_NodeTraversal_traverseAtScope_3821(scope)
}
`

	refCollectorAndPassesSrc := `package jscomp

import (
	"fmt"

	"example.com/w21ctest/j2gort"
	pkg_rhino "example.com/w21ctest/rhino"
)

type BasicBlock struct {
	parent *BasicBlock
	root   pkg_rhino.Node
}

func NewBasicBlock(parent *BasicBlock, root pkg_rhino.Node) *BasicBlock {
	this_ := &BasicBlock{}
	InitBasicBlock(this_, parent, root)
	return this_
}

func InitBasicBlock(this_ *BasicBlock, parent *BasicBlock, root pkg_rhino.Node) {
	this_.parent = parent
	this_.root = root
}

func (b *BasicBlock) GetRoot() pkg_rhino.Node { return b.root }

type Reference struct {
	nameNode   pkg_rhino.Node
	basicBlock *BasicBlock
	scope      *Scope
}

func InitReference_3(this_ *Reference, n pkg_rhino.Node, t *NodeTraversal, bb *BasicBlock) {
	this_.nameNode = n
	this_.basicBlock = bb
	this_.scope = t.GetScope()
}

type ReferenceCollection struct {
	ref0 *Reference
	ref1 *Reference
	rest []*Reference
}

func NewReferenceCollection() *ReferenceCollection {
	this_ := &ReferenceCollection{}
	InitReferenceCollection(this_)
	return this_
}

func InitReferenceCollection(this_ *ReferenceCollection) {}

func (rc *ReferenceCollection) Add(r *Reference) {
	if rc.rest != nil {
		rc.rest = append(rc.rest, r)
		return
	}
	if rc.ref0 == nil {
		rc.ref0 = r
		return
	}
	if rc.ref1 == nil {
		rc.ref1 = r
		return
	}
	rc.rest = []*Reference{rc.ref0, rc.ref1, r}
	rc.ref0 = nil
	rc.ref1 = nil
}

func (rc *ReferenceCollection) Size() int {
	if rc.rest != nil {
		return len(rc.rest)
	}
	if rc.ref0 == nil {
		return 0
	}
	if rc.ref1 == nil {
		return 1
	}
	return 2
}

type ReferenceCollector_Behavior interface {
	IsJavaNil_() bool
	AfterExitScope(*NodeTraversal, *ReferenceCollector)
}

type doNothingBehavior struct{}

func (b *doNothingBehavior) IsJavaNil_() bool                                  { return b == nil }
func (b *doNothingBehavior) AfterExitScope(*NodeTraversal, *ReferenceCollector) {}

var doNothingSingleton ReferenceCollector_Behavior = &doNothingBehavior{}

func StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR() ReferenceCollector_Behavior {
	return doNothingSingleton
}

var lambdaSingleton_reference_collector_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_reference_collector_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_9 any) any {
		k := j2gort.CastPtr[Var](lp_9)
		_ = k
		return NewReferenceCollection()
	})
}

type ReferenceCollector struct {
	referenceMap j2gort.CustomOrderedMap[*Var, *ReferenceCollection]
	blockStack   *j2gort.Deque[*BasicBlock]
	behavior     ReferenceCollector_Behavior
	scopeCreator ScopeCreator
	compiler     AbstractCompiler
	callback     *ReferenceCollector_CollectorCallback
}

func clinit_ReferenceCollector() {}

func (this_ *ReferenceCollector) IsJavaNil_() bool { return this_ == nil }

func (this_ *ReferenceCollector) initFields_ReferenceCollector() {
	this_.referenceMap = j2gort.NewCustomOrderedMap[*Var, *ReferenceCollection]()
	this_.blockStack = &j2gort.Deque[*BasicBlock]{}
	this_.callback = NewReferenceCollector_CollectorCallback(this_)
}

func NewReferenceCollector_3(compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator) *ReferenceCollector {
	clinit_ReferenceCollector()
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

func (this_ *ReferenceCollector) priv_ReferenceCollector_createTraversalBuilder_4529() *NodeTraversal_Builder {
	return NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(ToIface_NodeTraversal_Callback(this_.callback)).SetScopeCreator(this_.scopeCreator).SetObeyDestructuringAndDefaultValueExecutionOrder(true)
}

func (this_ *ReferenceCollector) Process_1(root pkg_rhino.Node) {
	this_.priv_ReferenceCollector_createTraversalBuilder_4529().Traverse(root)
}

func (this_ *ReferenceCollector) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.priv_ReferenceCollector_createTraversalBuilder_4529().TraverseRoots(externs, root)
}

func (this_ *ReferenceCollector) GetReferences(v *Var) *ReferenceCollection {
	return this_.referenceMap.Get(v)
}

func (this_ *ReferenceCollector) GetAllSymbols() []*Var {
	return this_.referenceMap.KeysSlice()
}

func (this_ *ReferenceCollector) priv_ReferenceCollector_pushNewBlock_4531(root pkg_rhino.Node) {
	this_.blockStack.AddLast(NewBasicBlock(this_.blockStack.PeekLast(), root))
}

func (this_ *ReferenceCollector) priv_ReferenceCollector_popLastBlock_4532() {
	_ = this_.blockStack.RemoveLast()
}

func (this_ *ReferenceCollector) priv_ReferenceCollector_addReference_4533(v *Var, n pkg_rhino.Node, t *NodeTraversal) *Reference {
	var collection *ReferenceCollection = j2gort.CustomOrderedMapComputeIfAbsentFn(&this_.referenceMap, v, nil, lambdaSingleton_reference_collector_1)
	_ = collection
	var val_10 Reference
	InitReference_3(&val_10, n, t, this_.blockStack.GetLast())
	var reference *Reference = &val_10
	collection.Add(reference)
	return reference
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
	cb.outer_.priv_ReferenceCollector_pushNewBlock_4531(t.GetScopeRoot())
}
func (cb *ReferenceCollector_CollectorCallback) ExitScope(t *NodeTraversal) {
	cb.outer_.behavior.AfterExitScope(t, cb.outer_)
	cb.outer_.priv_ReferenceCollector_popLastBlock_4532()
}
func (cb *ReferenceCollector_CollectorCallback) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) {
	sc := t.GetScope()
	v := sc.GetOwnDeclaredVar(fmt.Sprintf("v_%d_6", sc.GetRootNode_0()))
	if v != nil {
		// Add 3 references so ReferenceCollection.rest slice is also exercised!
		cb.outer_.priv_ReferenceCollector_addReference_4533(v, n, t)
		cb.outer_.priv_ReferenceCollector_addReference_4533(v, n, t)
		cb.outer_.priv_ReferenceCollector_addReference_4533(v, n, t)
	}
}

// --- NodeTraversal Callbacks: Positive vs Negative Soundness ---

// 1. Positive Fast Path 1: never touches GetScope().
type NoScopeCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	visits int
}

func (cb *NoScopeCallback) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) {
	cb.visits++
}

// 2. Positive Fast Path 2: calls t.GetScope(), has outer_ *VarCheckPass (all fields scope/var-free).
type VarCheckPass struct {
	compiler AbstractCompiler
	seenVars int
}

type VarCheckPass_Callback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *VarCheckPass
}

func (cb *VarCheckPass_Callback) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) {
	sc := t.GetScope()
	if sc.GetOwnDeclaredVar(fmt.Sprintf("v_%d_7", sc.GetRootNode_0())) != nil {
		cb.outer_.seenVars++
	}
}

// 3. Negative Case A: OptimizeCalls-style callback storing *Scope in a field read AFTER TraverseRoots returns!
type OptimizeCalls_ReferenceMapBuildingCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	globalScope *Scope
}

func (cb *OptimizeCalls_ReferenceMapBuildingCallback) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) {
	if cb.globalScope == nil {
		cb.globalScope = t.GetScope()
	}
}

// 4. Negative Case B: Lambda wrapper holding arbitrary callback interface.
type NodeTraversal__Anon_66 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	cap_0_cb NodeTraversal_Callback
}

func (cb *NodeTraversal__Anon_66) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) {
	cb.cap_0_cb.Visit_3(t, n, p)
}

// 5. Negative Case C: Mutates a package-level global during Visit_3!
var globalVisitCounter int

type GlobalEscapingCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

func (cb *GlobalEscapingCallback) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) {
	if t.GetScope() != nil {
		globalVisitCounter++
	}
}

// --- Passes Exercising RemoveUnusedCode, CoalesceVariableNames, InlineVariables, VariableReferenceCheck, Escaping SymbolTable ---

type RemoveUnusedCode_RemovableBuilder struct {
	outer_ *RemoveUnusedCode
	tag    int32
}

func NewRemoveUnusedCode_RemovableBuilder(outer_ *RemoveUnusedCode) *RemoveUnusedCode_RemovableBuilder {
	this_ := &RemoveUnusedCode_RemovableBuilder{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_RemovableBuilder(this_)
	return this_
}

func InitRemoveUnusedCode_RemovableBuilder(this_ *RemoveUnusedCode_RemovableBuilder) {
	this_.tag = 1
}

type RemoveUnusedCode struct {
	compiler               AbstractCompiler
	scopeCreator           *SyntacticScopeCreator
	allFunctionParamScopes []*Scope
	varInfoKeys            []*Var
}

func NewRemoveUnusedCode(compiler AbstractCompiler) *RemoveUnusedCode {
	return &RemoveUnusedCode{
		compiler:     compiler,
		scopeCreator: NewSyntacticScopeCreator_1(compiler),
	}
}

func (this_ *RemoveUnusedCode) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	sc := this_.scopeCreator.CreateScope_2_Node_Scope(root, nil)
	this_.allFunctionParamScopes = append(this_.allFunctionParamScopes, sc)
	b := NewRemoveUnusedCode_RemovableBuilder(this_)
	_ = b
}

type CoalesceVariableNames struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler     AbstractCompiler
	scopeCreator *MemoizedScopeCreator
}

func NewCoalesceVariableNames(compiler AbstractCompiler) *CoalesceVariableNames {
	return &CoalesceVariableNames{
		compiler:     compiler,
		scopeCreator: NewMemoizedScopeCreator(NewSyntacticScopeCreator_1(compiler)),
	}
}

func (this_ *CoalesceVariableNames) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) {
	_ = t.GetScope()
}

func (this_ *CoalesceVariableNames) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(this_).SetScopeCreator(ToIface_ScopeCreator(this_.scopeCreator)).Traverse(root)
}

type InlineVariables struct {
	compiler AbstractCompiler
	seenRefs int
}

type InlineVariables_InliningBehavior struct {
	outer_ *InlineVariables
}

func NewInlineVariables_InliningBehavior(outer_ *InlineVariables) *InlineVariables_InliningBehavior {
	return &InlineVariables_InliningBehavior{outer_: outer_}
}

func (b *InlineVariables_InliningBehavior) IsJavaNil_() bool { return b == nil }
func (b *InlineVariables_InliningBehavior) AfterExitScope(t *NodeTraversal, rc *ReferenceCollector) {
	for _, v := range rc.GetAllSymbols() {
		if col := rc.GetReferences(v); col != nil {
			b.outer_.seenRefs += col.Size()
		}
	}
}

func (this_ *InlineVariables) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	clinit_ReferenceCollector()
	var val_1 ReferenceCollector
	InitReferenceCollector_3(&val_1, this_.compiler, NewInlineVariables_InliningBehavior(this_), NewSyntacticScopeCreator_1(this_.compiler))
	var callback *ReferenceCollector = &val_1
	_ = callback
	callback.Process_2(externs, root)
}

type VariableReferenceCheck struct {
	compiler               AbstractCompiler
	blocksWithDeclarations j2gort.OrderedSet[*BasicBlock]
}

type VariableReferenceCheck_ReferenceCheckingBehavior struct {
	outer_ *VariableReferenceCheck
}

func NewVariableReferenceCheck_ReferenceCheckingBehavior(outer_ *VariableReferenceCheck) *VariableReferenceCheck_ReferenceCheckingBehavior {
	return &VariableReferenceCheck_ReferenceCheckingBehavior{outer_: outer_}
}

func (b *VariableReferenceCheck_ReferenceCheckingBehavior) IsJavaNil_() bool { return b == nil }
func (b *VariableReferenceCheck_ReferenceCheckingBehavior) AfterExitScope(t *NodeTraversal, rc *ReferenceCollector) {
	for _, v := range rc.GetAllSymbols() {
		if col := rc.GetReferences(v); col != nil && len(col.rest) > 0 {
			b.outer_.blocksWithDeclarations.Add(col.rest[0].basicBlock)
		}
	}
}

func (this_ *VariableReferenceCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	{
		clinit_ReferenceCollector()
		var val_1 ReferenceCollector
		InitReferenceCollector_3(&val_1, this_.compiler, NewVariableReferenceCheck_ReferenceCheckingBehavior(this_), NewSyntacticScopeCreator_1(this_.compiler))
		(&val_1).Process_2(externs, root)
	}
}

// Negative Escape Case for ReferenceCollector: Compiler.BuildKnownSymbolTable!
type SymbolTable struct {
	saved map[string]*ReferenceCollection
}

func (st *SymbolTable) AddSymbolsFrom(rc *ReferenceCollector) {
	if st.saved == nil {
		st.saved = make(map[string]*ReferenceCollection)
	}
	for _, v := range rc.GetAllSymbols() {
		st.saved[v.GetName()] = rc.GetReferences(v)
	}
}

func BuildKnownSymbolTable(compiler AbstractCompiler, root pkg_rhino.Node) *SymbolTable {
	st := &SymbolTable{}
	clinit_ReferenceCollector()
	var val_81 ReferenceCollector
	InitReferenceCollector_3(&val_81, compiler, StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR(), NewSyntacticScopeCreator_1(compiler))
	var refCollector *ReferenceCollector = &val_81
	_ = refCollector
	refCollector.Process_1(root)
	st.AddSymbolsFrom(refCollector)
	return st
}
`

	files := map[string]string{
		"scope_and_var.go":  scopeAndVarSrc,
		"node_traversal.go": nodeTraversalSrc,
		"ref_collector.go":  refCollectorAndPassesSrc,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(jscompDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] < 10 || stats1.ByPass["node_traversal_lowering"] < 8 {
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

	// Verify structural lowering in scope_and_var.go.
	svBytes, err := os.ReadFile(filepath.Join(jscompDir, "scope_and_var.go"))
	if err != nil {
		t.Fatal(err)
	}
	sv := string(svBytes)
	for _, want := range []string{
		"slab_Scope",
		"j2gort.PooledSlab[Scope]",
		"slab_Var",
		"j2gort.PooledSlab[Var]",
		"varSlab_ *j2gort.PooledSlab[Var]",
		"j2gort.AcquireOrderedMapPtrCap[string, *Var](int32(16))",
		"func (this_ *SyntacticScopeCreator) ReleaseSlabs_()",
		"this_.slab_Scope.ReleaseWith(func(s_ *Scope)",
		"s_.vars.ReleaseToPool()",
		"j2gort.AcquireOrderedMap[pkg_rhino.Node, AbstractScope]()",
		"func (this_ *MemoizedScopeCreator) ReleaseSlabs_()",
		"this_.scopesByScopeRoot.ReleaseToPool()",
	} {
		if !strings.Contains(sv, want) {
			t.Fatalf("missing %q in lowered scope_and_var.go:\n%s", want, sv)
		}
	}

	// Verify structural lowering in node_traversal.go.
	ntBytes, err := os.ReadFile(filepath.Join(jscompDir, "node_traversal.go"))
	if err != nil {
		t.Fatal(err)
	}
	nt := string(ntBytes)
	for _, want := range []string{
		"canReleaseScopes_",
		"CanReleaseSyntacticScopes_()",
		"if t_.canReleaseScopes_ {",
		"t_.defaultScopeCreator_.ReleaseSlabs_()",
		"savedScanner_ := t_.defaultScopeCreator_.reusableScanner",
		"*t_ = NodeTraversal{}",
	} {
		if !strings.Contains(nt, want) {
			t.Fatalf("missing %q in lowered node_traversal.go:\n%s", want, nt)
		}
	}

	// Verify structural lowering + soundness guards in ref_collector.go.
	rcBytes, err := os.ReadFile(filepath.Join(jscompDir, "ref_collector.go"))
	if err != nil {
		t.Fatal(err)
	}
	rc := string(rcBytes)
	for _, want := range []string{
		"slab_Reference",
		"j2gort.PooledSlab[Reference]",
		"slab_BasicBlock",
		"j2gort.PooledSlab[BasicBlock]",
		"slab_ReferenceCollection",
		"j2gort.PooledSlab[ReferenceCollection]",
		"ownsScopeCreator_",
		"j2gort.AcquireCustomOrderedMap[*Var, *ReferenceCollection]()",
		"this_.slab_ReferenceCollection.ReleaseWith(func(rc_ *ReferenceCollection)",
		"this_.referenceMap.ReleaseToPool()",
		"callback.ownsScopeCreator_ = true",
		"defer callback.ReleaseSlabs_()",
		"this_.blocksWithDeclarations.Clear()",
		"func (*NoScopeCallback) CanReleaseSyntacticScopes_()",
		"func (*VarCheckPass_Callback) CanReleaseSyntacticScopes_()",
		"this_.scopeCreator = nil",
	} {
		if !strings.Contains(rc, want) {
			t.Fatalf("missing %q in lowered ref_collector.go:\n%s", want, rc)
		}
	}

	// Soundness checks: NONE of the escaping callbacks or BuildKnownSymbolTable may release!
	for _, forbidden := range []string{
		"func (*ReferenceCollector_CollectorCallback) CanReleaseSyntacticScopes_()",
		"func (*OptimizeCalls_ReferenceMapBuildingCallback) CanReleaseSyntacticScopes_()",
		"func (*NodeTraversal__Anon_66) CanReleaseSyntacticScopes_()",
		"func (*GlobalEscapingCallback) CanReleaseSyntacticScopes_()",
		"defer refCollector.ReleaseSlabs_()",
		"refCollector.ownsScopeCreator_ = true",
	} {
		if strings.Contains(rc, forbidden) {
			t.Fatalf("UNSOUND lowering found %q in lowered ref_collector.go:\n%s", forbidden, rc)
		}
	}

	// Runtime behavioral + -race verification!
	rtTest := `package jscomp

import (
	"sync"
	"testing"

	"example.com/w21ctest/j2gort"
)

type mockComp struct{ opts *CompilerOptions }

func (c *mockComp) IsJavaNil_() bool             { return c == nil }
func (c *mockComp) GetOptions() *CompilerOptions { return c.opts }

func TestWave21CRuntimePoolingAndSoundness(t *testing.T) {
	j2gort.ResetWave21CMetrics()
	comp := &mockComp{opts: &CompilerOptions{}}

	// 1. Scope-free callback (VarCheckPass_Callback) MUST recycle Scope + Var + Scope.vars slabs!
	vc := &VarCheckPass{compiler: comp}
	vcCb := &VarCheckPass_Callback{outer_: vc}
	for i := 0; i < 40; i++ {
		NodeTraversal_Traverse(comp, 3, vcCb)
	}
	if vc.seenVars != 120 {
		t.Fatalf("expected 120 seenVars, got %d", vc.seenVars)
	}
	if j2gort.PooledSlabReleaseWithCalls.Load() < 40 || j2gort.OrderedMapReleases.Load() < 120 {
		t.Fatalf("expected >=40 ReleaseWith and >=120 OrderedMapReleases, got ReleaseWith=%d OrderedMapReleases=%d",
			j2gort.PooledSlabReleaseWithCalls.Load(), j2gort.OrderedMapReleases.Load())
	}
	// Without PooledSlab recycling, 40 traversals * (1 Scope chunk + 1 Var chunk) = 80 chunk allocs;
	// even with Go -race randomly dropping ~25% of sync.Pool.Put calls, chunk allocs stay well below 80!
	if chunks := j2gort.PooledSlabChunkAllocs.Load(); chunks >= 55 {
		t.Fatalf("expected < 55 PooledSlab chunk allocs across 40 traversals under -race (80 unpooled), got %d", chunks)
	}

	// 2. Escaping callback (OptimizeCalls_ReferenceMapBuildingCallback):
	//    - MUST keep cb.globalScope + its Vars valid after traversal returns!
	//    - MUST detach defaultScopeCreator_ from nodeTraversalPool_ (0 live chunks pinned in pooled NodeTraversal)!
	escCb := &OptimizeCalls_ReferenceMapBuildingCallback{}
	tInst := acquireNodeTraversal_(comp, escCb, nil, false)
	if tInst.canReleaseScopes_ {
		t.Fatalf("OptimizeCalls_ReferenceMapBuildingCallback must have canReleaseScopes_ == false")
	}
	tInst.priv_NodeTraversal_traverse_3817(2)
	savedGlobal := escCb.globalScope
	if savedGlobal == nil || savedGlobal.GetOwnDeclaredVar("v_1_7") == nil {
		t.Fatalf("missing expected declared var on escaped globalScope")
	}
	releaseNodeTraversal_(tInst)
	// Pop tInst from pool (if retained) to verify 0 live chunks pinned, and exercise 30 pooled traversals!
	if tPop, _ := nodeTraversalPool_.Get().(*NodeTraversal); tPop != nil {
		if tPop.defaultScopeCreator_.slab_Scope.HasLiveChunks() || tPop.defaultScopeCreator_.slab_Var.HasLiveChunks() {
			t.Fatalf("defaultScopeCreator_ leaked live chunks inside pooled NodeTraversal!")
		}
		nodeTraversalPool_.Put(tPop)
	}
	for i := 0; i < 30; i++ {
		NodeTraversal_Traverse(comp, 3, escCb)
	}
	if savedGlobal.GetOwnDeclaredVar("v_1_7") == nil || savedGlobal.GetOwnDeclaredVar("v_1_7").GetName() != "v_1_7" {
		t.Fatalf("escaped savedGlobal corrupted after subsequent traversals!")
	}

	// 3. RemoveUnusedCode & CoalesceVariableNames release scopeCreator and nil it!
	ruc := NewRemoveUnusedCode(comp)
	ruc.Process_2(0, 5)
	if ruc.scopeCreator != nil || ruc.allFunctionParamScopes != nil {
		t.Fatalf("RemoveUnusedCode failed to nil scopeCreator / allFunctionParamScopes")
	}
	cvn := NewCoalesceVariableNames(comp)
	cvn.Process_2(0, 4)
	if cvn.scopeCreator != nil {
		t.Fatalf("CoalesceVariableNames failed to nil scopeCreator")
	}

	// 4. Non-escaping ReferenceCollector passes (InlineVariables + VariableReferenceCheck) recycle slabs + CustomOrderedMap!
	customRelBefore := j2gort.CustomOrderedMapReleases.Load()
	iv := &InlineVariables{compiler: comp}
	vrc := &VariableReferenceCheck{compiler: comp}
	for i := 0; i < 25; i++ {
		iv.Process_2(4, 3)
		vrc.Process_2(4, 3)
		if vrc.blocksWithDeclarations.Len() != 0 {
			t.Fatalf("VariableReferenceCheck.blocksWithDeclarations not cleared!")
		}
	}
	if iv.seenRefs == 0 {
		t.Fatalf("InlineVariables saw 0 references")
	}
	if got := j2gort.CustomOrderedMapReleases.Load() - customRelBefore; got != 50 {
		t.Fatalf("expected 50 CustomOrderedMap releases from 25*(InlineVariables+VariableReferenceCheck), got %d", got)
	}

	// 5. Escaping BuildKnownSymbolTable MUST keep SymbolTable references valid across subsequent ReferenceCollector runs!
	symTab := BuildKnownSymbolTable(comp, 3)
	rcSaved := symTab.saved["v_3_6"]
	if rcSaved == nil || rcSaved.Size() != 3 || rcSaved.rest[0].scope.GetRootNode_0() != 3 {
		t.Fatalf("unexpected BuildKnownSymbolTable saved refs: %+v", rcSaved)
	}
	for i := 0; i < 20; i++ {
		iv.Process_2(4, 2)
	}
	if rcSaved.Size() != 3 || rcSaved.rest[0].scope.GetRootNode_0() != 3 {
		t.Fatalf("BuildKnownSymbolTable escaped ReferenceCollection corrupted by pooled ReferenceCollector!")
	}

	// 6. Concurrent -race stress test across pooled NodeTraversal + ReferenceCollector + SyntacticScopeCreator!
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localComp := &mockComp{opts: &CompilerOptions{}}
			localVC := &VarCheckPass{compiler: localComp}
			localCb := &VarCheckPass_Callback{outer_: localVC}
			localIV := &InlineVariables{compiler: localComp}
			for iter := 0; iter < 30; iter++ {
				NodeTraversal_Traverse(localComp, 3, localCb)
				localIV.Process_2(4, 3)
			}
		}()
	}
	wg.Wait()
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "wave21c_rt_test.go"), []byte(rtTest), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-race", "-v", "./jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race ./jscomp failed: %v\n%s", err, string(out))
	}
}

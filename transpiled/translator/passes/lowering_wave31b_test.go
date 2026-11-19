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

func TestWave31B_SlabAndLoggingAndSetLoweringAndRuntime(t *testing.T) {
	tmpDir := t.TempDir()
	goMod := "module example.com/w31btest\n\ngo 1.24\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	j2Dir := filepath.Join(tmpDir, "gen", "j2gort")
	rhinoDir := filepath.Join(tmpDir, "gen", "rhino")
	diagDir := filepath.Join(tmpDir, "gen", "diagnostic")
	jscompDir := filepath.Join(tmpDir, "gen", "jscomp")
	for _, d := range []string{j2Dir, rhinoDir, diagDir, jscompDir} {
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
	ReleaseOrdMapCalls  atomic.Int64
	SupplierLambdaCalls atomic.Int64
	LinkedHashSetCalls  atomic.Int64
)

func ResetWave31BMetrics() {
	PurgeCalls.Store(0)
	SlabAllocCalls.Store(0)
	SlabResetCalls.Store(0)
	SlabReleaseCalls.Store(0)
	ReleaseOrdMapCalls.Store(0)
	SupplierLambdaCalls.Store(0)
	LinkedHashSetCalls.Store(0)
}

func PurgeAllRuntimePools()  { PurgeCalls.Add(1) }
func TrimStaleSlabPools(int) {}
func CheckState(cond bool) {
	if !cond {
		panic("CheckState")
	}
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
func CheckNotNullNode[N ~uint32](n N) N {
	if n == 0 {
		panic("NPE")
	}
	return n
}
func AssignIdentHash(p *atomic.Int32) int32 {
	if v := p.Load(); v != 0 {
		return v
	}
	p.Store(1)
	return 1
}

type Slab[T any] struct{ buf []T }

func (s *Slab[T]) Alloc() *T {
	s.buf = append(s.buf, *new(T))
	return &s.buf[len(s.buf)-1]
}
func (s *Slab[T]) Release() { s.buf = nil }

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

type Ext_java_util_function_Supplier interface {
	Implements_j2gort_Ext_java_util_function_Supplier()
	IsJavaNil_() bool
	Get() any
}

type fnSupplier struct{ fn func() any }

func (*fnSupplier) Implements_j2gort_Ext_java_util_function_Supplier() {}
func (s *fnSupplier) IsJavaNil_() bool                                 { return s == nil }
func (s *fnSupplier) Get() any                                         { return s.fn() }

func NewLambda_Ext_java_util_function_Supplier(fn func() any) Ext_java_util_function_Supplier {
	SupplierLambdaCalls.Add(1)
	return &fnSupplier{fn: fn}
}

func BoxString(s string) any   { return s }
func CastString(v any) string  { return v.(string) }
func StrValString(s string) string { return s }
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

type Ext_java_lang_Iterable interface {
	Implements_j2gort_Ext_java_lang_Iterable()
	IsJavaNil_() bool
}
type Ext_java_util_Collection interface {
	Ext_java_lang_Iterable
	Implements_j2gort_Ext_java_util_Collection()
}
type Ext_java_util_Set interface {
	Ext_java_util_Collection
	Implements_j2gort_Ext_java_util_Set()
	Add_1(any) bool
	Contains(any) bool
	Size() int32
	IterItems_() []any
}
type Ext_com_google_common_collect_ImmutableSet interface {
	Ext_java_util_Set
	Implements_j2gort_Ext_com_google_common_collect_ImmutableSet()
}

type linkedHashSet struct {
	items []any
}

func (*linkedHashSet) Implements_j2gort_Ext_java_lang_Iterable()   {}
func (*linkedHashSet) Implements_j2gort_Ext_java_util_Collection() {}
func (*linkedHashSet) Implements_j2gort_Ext_java_util_Set()        {}
func (*linkedHashSet) Implements_j2gort_Ext_com_google_common_collect_ImmutableSet() {}
func (s *linkedHashSet) IsJavaNil_() bool { return s == nil }
func (s *linkedHashSet) Size() int32      { return int32(len(s.items)) }
func (s *linkedHashSet) IterItems_() []any { return s.items }
func (s *linkedHashSet) Contains(v any) bool {
	for _, x := range s.items {
		if x == v {
			return true
		}
	}
	return false
}
func (s *linkedHashSet) Add_1(v any) bool {
	if s.Contains(v) {
		return false
	}
	s.items = append(s.items, v)
	return true
}

func NewExt_java_util_LinkedHashSet() *linkedHashSet {
	LinkedHashSetCalls.Add(1)
	return &linkedHashSet{}
}
func NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(col Ext_java_util_Collection) *linkedHashSet {
	LinkedHashSetCalls.Add(1)
	out := &linkedHashSet{}
	if st, ok := col.(Ext_java_util_Set); ok && st != nil {
		for _, v := range st.IterItems_() {
			out.Add_1(v)
		}
	}
	return out
}
func Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(col Ext_java_util_Collection) Ext_com_google_common_collect_ImmutableSet {
	out := &linkedHashSet{}
	if st, ok := col.(Ext_java_util_Set); ok && st != nil {
		for _, v := range st.IterItems_() {
			out.Add_1(v)
		}
	}
	return out
}
func Ext_com_google_common_collect_ImmutableSet_Of_1(v any) Ext_com_google_common_collect_ImmutableSet {
	out := &linkedHashSet{}
	out.Add_1(v)
	return out
}
func SetAddStr(s Ext_java_util_Set, v string) bool { return s.Add_1(v) }

type IterCursor struct {
	items []any
	pos   int
}

func NewIterCursor(col any) IterCursor {
	if st, ok := col.(Ext_java_util_Set); ok && st != nil {
		return IterCursor{items: st.IterItems_()}
	}
	return IterCursor{}
}
func (c *IterCursor) HasNext() bool { return c.pos < len(c.items) }
func (c *IterCursor) Next() any {
	v := c.items[c.pos]
	c.pos++
	return v
}

type OrderedSet[T comparable] struct {
	items []T
}
type OrderedSetCursor[T comparable] struct {
	items []T
	pos   int
}

func NewOrderedSet[T comparable]() OrderedSet[T]      { return OrderedSet[T]{} }
func NewOrderedSetPtr[T comparable]() *OrderedSet[T]  { return &OrderedSet[T]{} }
func NewOrderedSetPtrCap[T comparable](int32) *OrderedSet[T] { return &OrderedSet[T]{} }
func OrderedSetFromSet[T comparable](src OrderedSet[T]) OrderedSet[T] {
	cp := make([]T, len(src.items))
	copy(cp, src.items)
	return OrderedSet[T]{items: cp}
}
func OrderedSetPtrFromSet[T comparable](src *OrderedSet[T]) *OrderedSet[T] {
	os := OrderedSetFromSet(*src)
	return &os
}
func OrderedSetFromStrCol(col any) OrderedSet[string] {
	var out OrderedSet[string]
	if st, ok := col.(Ext_java_util_Set); ok && st != nil {
		for _, v := range st.IterItems_() {
			out.Add(v.(string))
		}
	}
	return out
}
func OrderedSetPtrFromStrCol(col any) *OrderedSet[string] {
	s := OrderedSetFromStrCol(col)
	return &s
}
func OrderedSetFromCol[T comparable](col any) OrderedSet[T] {
	var out OrderedSet[T]
	if st, ok := col.(Ext_java_util_Set); ok && st != nil {
		for _, v := range st.IterItems_() {
			out.Add(v.(T))
		}
	}
	return out
}
func OrderedSetPtrFromCol[T comparable](col any) *OrderedSet[T] {
	s := OrderedSetFromCol[T](col)
	return &s
}
func (s *OrderedSet[T]) Contains(v T) bool {
	for _, x := range s.items {
		if x == v {
			return true
		}
	}
	return false
}
func (s *OrderedSet[T]) Add(v T) bool {
	if s.Contains(v) {
		return false
	}
	s.items = append(s.items, v)
	return true
}
func (s OrderedSet[T]) IsEmpty() bool { return len(s.items) == 0 }
func (s OrderedSet[T]) Size() int32   { return int32(len(s.items)) }
func (s *OrderedSet[T]) Iter() OrderedSetCursor[T] {
	if s == nil {
		panic("NPE")
	}
	return OrderedSetCursor[T]{items: s.items}
}
func (c *OrderedSetCursor[T]) HasNext() bool { return c.pos < len(c.items) }
func (c *OrderedSetCursor[T]) Next() T {
	v := c.items[c.pos]
	c.pos++
	return v
}
func (s OrderedSet[T]) ToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
	out := &linkedHashSet{}
	for _, x := range s.items {
		out.Add_1(x)
	}
	return out
}
func (s OrderedSet[T]) ToJavaSet() Ext_java_util_Set {
	out := &linkedHashSet{}
	for _, x := range s.items {
		out.Add_1(x)
	}
	return out
}

type OrderedMap[K comparable, V any] struct {
	m map[K]V
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V)}
}
func (om *OrderedMap[K, V]) Put(k K, v V) {
	if om.m == nil {
		om.m = make(map[K]V)
	}
	om.m[k] = v
}
func (om OrderedMap[K, V]) Get(k K) V         { return om.m[k] }
func (om OrderedMap[K, V]) ContainsKey(k K) bool { _, ok := om.m[k]; return ok }
func OrderedMapGetStr[K comparable](om OrderedMap[K, string], k K) string {
	if v, ok := om.m[k]; ok {
		return v
	}
	return ""
}
func ReleaseOrderedMap[K comparable, V any](om *OrderedMap[K, V]) {
	ReleaseOrdMapCalls.Add(1)
	if om != nil {
		om.m = nil
	}
}

type Deque[T any] struct {
	items []T
}

func (d *Deque[T]) Push(v T)    { d.items = append([]T{v}, d.items...) }
func (d *Deque[T]) Peek() T     { return d.items[0] }
func (d *Deque[T]) Pop() T      { v := d.items[0]; d.items = d.items[1:]; return v }
func (d *Deque[T]) AddLast(v T) { d.items = append(d.items, v) }
func (d *Deque[T]) RemoveFirst() T { v := d.items[0]; d.items = d.items[1:]; return v }
func (d *Deque[T]) IsEmpty() bool { return len(d.items) == 0 }
func (d *Deque[T]) Size() int32   { return int32(len(d.items)) }

func DequeFromCol[T any](col any) Deque[T] {
	var d Deque[T]
	if st, ok := col.(Ext_java_util_Set); ok && st != nil {
		for _, v := range st.IterItems_() {
			d.AddLast(v.(T))
		}
	}
	return d
}
func DequeFromOrderedSet[T comparable](s OrderedSet[T]) Deque[T] {
	cp := make([]T, len(s.items))
	copy(cp, s.items)
	return Deque[T]{items: cp}
}
func DequeFromOrderedSetPtr[T comparable](s *OrderedSet[T]) Deque[T] {
	return DequeFromOrderedSet(*s)
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
	TokIF       Token = 5
	TokAND      Token = 6
	TokOR       Token = 7
	TokNOT      Token = 8
	TokNAME     Token = 9
	TokVAR      Token = 10
)

func StaticGet_Token_NOT() Token { return TokNOT }
func StaticGet_Token_AND() Token { return TokAND }
func StaticGet_Token_OR() Token  { return TokOR }

type Node uint32

type nodeData struct {
	tok    Token
	str    string
	parent Node
	children []Node
}

var arena = []nodeData{{}}

func ResetArena() { arena = []nodeData{{}} }

func NewNode_1(tok Token) Node {
	id := Node(len(arena))
	arena = append(arena, nodeData{tok: tok})
	return id
}

func AllocNamed(tok Token, s string) Node {
	id := Node(len(arena))
	arena = append(arena, nodeData{tok: tok, str: s})
	return id
}

func (n Node) Srcref(Node) Node     { return n }
func (n Node) GetToken() Token      { return arena[n].tok }
func (n Node) TokenSwitchOrd_() int32 { return int32(arena[n].tok) }
func (n Node) GetString() string    { return arena[n].str }
func (n Node) SetString(s string)   { arena[n].str = s }
func (n Node) GetParent() Node      { return arena[n].parent }
func (n Node) HasParent() bool      { return arena[n].parent != 0 }
func (n Node) IsIf() bool           { return n != 0 && arena[n].tok == TokIF }
func (n Node) IsAnd() bool          { return n != 0 && arena[n].tok == TokAND }
func (n Node) IsNot() bool          { return n != 0 && arena[n].tok == TokNOT }
func (n Node) IsName() bool         { return n != 0 && arena[n].tok == TokNAME }
func (n Node) IsFunction() bool     { return n != 0 && arena[n].tok == TokFUNCTION }
func (n Node) IsScript() bool       { return n != 0 && arena[n].tok == TokSCRIPT }
func (n Node) GetFirstChild() Node {
	if n == 0 || len(arena[n].children) == 0 {
		return 0
	}
	return arena[n].children[0]
}
func (n Node) GetLastChild() Node {
	if n == 0 || len(arena[n].children) == 0 {
		return 0
	}
	return arena[n].children[len(arena[n].children)-1]
}
func (n Node) GetSecondChild() Node {
	if n == 0 || len(arena[n].children) < 2 {
		return 0
	}
	return arena[n].children[1]
}
func (n Node) GetNext() Node {
	p := arena[n].parent
	if p == 0 {
		return 0
	}
	ch := arena[p].children
	for i, c := range ch {
		if c == n && i+1 < len(ch) {
			return ch[i+1]
		}
	}
	return 0
}
func (n Node) AddChildToBack(c Node) {
	arena[c].parent = n
	arena[n].children = append(arena[n].children, c)
}
func (n Node) Detach() Node {
	p := arena[n].parent
	if p != 0 {
		ch := arena[p].children
		for i, c := range ch {
			if c == n {
				arena[p].children = append(ch[:i], ch[i+1:]...)
				break
			}
		}
		arena[n].parent = 0
	}
	return n
}
func (n Node) DetachChildren() {
	for _, c := range arena[n].children {
		arena[c].parent = 0
	}
	arena[n].children = arena[n].children[:0]
}
func (n Node) ReplaceWith(repl Node) {
	p := arena[n].parent
	if p == 0 {
		return
	}
	for i, c := range arena[p].children {
		if c == n {
			arena[p].children[i] = repl
			arena[repl].parent = p
			arena[n].parent = 0
			return
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(rhinoDir, "rhino.go"), []byte(rhinoSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	diagSrc := `package diagnostic

import "example.com/w31btest/gen/j2gort"

type LogFile interface {
	Implements_diagnostic_LogFile()
	IsJavaNil_() bool
	IsLogging() bool
	Log_1_Ext_java_util_function_Supplier(j2gort.Ext_java_util_function_Supplier)
}

type MemLogFile struct {
	Enabled bool
	Entries []string
}

func (*MemLogFile) Implements_diagnostic_LogFile() {}
func (l *MemLogFile) IsJavaNil_() bool             { return l == nil }
func (l *MemLogFile) IsLogging() bool              { return l != nil && l.Enabled }
func (l *MemLogFile) Log_1_Ext_java_util_function_Supplier(s j2gort.Ext_java_util_function_Supplier) {
	if l != nil && l.Enabled {
		l.Entries = append(l.Entries, s.Get().(string))
	}
}
`
	if err := os.WriteFile(filepath.Join(diagDir, "diag.go"), []byte(diagSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1A: Phase-boundary PurgeAllRuntimePools() + stripped mid-compilation PurgeAllRuntimePools/TrimStaleSlabPools.
	phaseSrc := `package jscomp

import (
	"example.com/w31btest/gen/j2gort"
	pkg_rhino "example.com/w31btest/gen/rhino"
)

type AbstractCompiler interface {
	IsJavaNil_() bool
}

type CompilerOptions_SegmentOfCompilationToRun uint8

type Compiler struct{}

func (c *Compiler) IsJavaNil_() bool { return c == nil }

func (this_ *Compiler) ParseInputs() pkg_rhino.Node {
	j2gort.TrimStaleSlabPools(1)
	return pkg_rhino.Node(1)
}

func (this_ *Compiler) PerformTranspilationAndOptimizations(seg CompilerOptions_SegmentOfCompilationToRun) {
	_ = seg
	j2gort.PurgeAllRuntimePools()
}

type TypeInferencePass struct{}

func (this_ *TypeInferencePass) InferAllScopes(inferenceRoot pkg_rhino.Node) int {
	_ = inferenceRoot
	j2gort.TrimStaleSlabPools(2)
	return 42
}

type TypedAstSerializer struct{}

func (this_ *TypedAstSerializer) SerializeRoots(externsRoot pkg_rhino.Node, jsRoot pkg_rhino.Node) int {
	_ = externsRoot
	_ = jsRoot
	defer j2gort.PurgeAllRuntimePools()
	return 7
}

type VariableReferenceCheck struct{}

func (this_ *VariableReferenceCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer j2gort.PurgeAllRuntimePools()
	j2gort.TrimStaleSlabPools(3)
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "phase_boundaries.go"), []byte(phaseSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1B: GuardedCallbackBase (testing upgrade from Wave 23C/30 form on Token.ROOT multi-script AST!).
	guardedSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w31btest/gen/j2gort"
	pkg_rhino "example.com/w31btest/gen/rhino"
)

type NodeTraversal struct{}

type NodeTraversal_Callback interface {
	ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool
	Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

func NodeTraversal_Traverse(_ AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	nt := &NodeTraversal{}
	var walk func(n, p pkg_rhino.Node)
	walk = func(n, p pkg_rhino.Node) {
		if !cb.ShouldTraverse_3(nt, n, p) {
			return
		}
		for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
			walk(c, n)
		}
		cb.Visit_3(nt, n, p)
	}
	walk(root, 0)
}

type GuardedCallback interface {
	NodeTraversal_Callback
	VisitGuarded(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type GuardedCallbackBase struct {
	self_                        GuardedCallback
	compiler                     AbstractCompiler
	contextStack                 j2gort.Deque[*GuardedCallback_Context]
	idHash_                      atomic.Int32
	slab_GuardedCallback_Context j2gort.PooledSlab[GuardedCallback_Context]
}

func (this_ *GuardedCallbackBase) BindSelf_(s any) { this_.self_ = s.(GuardedCallback) }

func InitGuardedCallback(this_ *GuardedCallbackBase, compiler AbstractCompiler) {
	this_.compiler = compiler
	this_.contextStack = j2gort.Deque[*GuardedCallback_Context]{}
}

func (this_ *GuardedCallbackBase) ShouldTraverse_3(traversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = traversal
	_ = n
	_ = parent
	if parent == 0 {
		this_.contextStack.Push(StaticGet_GuardedCallback_Context_EMPTY())
	} else {
		this_.contextStack.Push(this_.contextStack.Peek().descendWithSlab_(&this_.slab_GuardedCallback_Context, this_.compiler, parent, n))
	}
	return true
}

func (this_ *GuardedCallbackBase) Visit_3(traversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = n
	_ = parent
	this_.self_.VisitGuarded(traversal, n, parent)
	this_.contextStack.Pop()
	if this_.contextStack.IsEmpty() {
		this_.slab_GuardedCallback_Context.Reset()
	}
}

func (this_ *GuardedCallbackBase) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_GuardedCallback_Context.Release()
}

type GuardedCallback_Context struct {
	linked      *GuardedCallback_Context
	conditional pkg_rhino.Node
	idHash_     atomic.Int32
	safe        bool
}

var (
	static_GuardedCallback_Context_EMPTY *GuardedCallback_Context
	clinitDone_GuardedCtx                atomic.Bool
)

func clinit_GuardedCallback_Context() {
	if clinitDone_GuardedCtx.CompareAndSwap(false, true) {
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
func allocPooled_GuardedCallback_Context(slab_ *j2gort.PooledSlab[GuardedCallback_Context], conditional pkg_rhino.Node, safe bool, linked *GuardedCallback_Context) *GuardedCallback_Context {
	clinit_GuardedCallback_Context()
	this_ := slab_.Alloc()
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
func (this_ *GuardedCallback_Context) linkWithSlab_(slab_ *j2gort.PooledSlab[GuardedCallback_Context], cond pkg_rhino.Node, safe bool) *GuardedCallback_Context {
	return allocPooled_GuardedCallback_Context(slab_, cond, safe, this_)
}
func (this_ *GuardedCallback_Context) Propagate(safe bool) *GuardedCallback_Context {
	if safe == this_.safe {
		return this_
	}
	return NewGuardedCallback_Context(this_.conditional, safe, this_.linked)
}
func (this_ *GuardedCallback_Context) propagateWithSlab_(slab_ *j2gort.PooledSlab[GuardedCallback_Context], safe bool) *GuardedCallback_Context {
	if safe == this_.safe {
		return this_
	}
	return allocPooled_GuardedCallback_Context(slab_, this_.conditional, safe, this_.linked)
}
func (this_ *GuardedCallback_Context) Descend(_ AbstractCompiler, parent pkg_rhino.Node, child pkg_rhino.Node) *GuardedCallback_Context {
	if parent.IsIf() {
		if child == parent.GetFirstChild() {
			return this_.Link(parent, true)
		}
		return StaticGet_GuardedCallback_Context_EMPTY()
	}
	if parent.IsAnd() {
		if child == parent.GetFirstChild() {
			return this_.Link(parent, true)
		}
		return this_
	}
	if parent.IsScript() || parent.GetToken() == pkg_rhino.TokROOT || parent.GetToken() == pkg_rhino.TokBLOCK {
		return StaticGet_GuardedCallback_Context_EMPTY()
	}
	return this_.Propagate(false)
}
func (this_ *GuardedCallback_Context) descendWithSlab_(slab_ *j2gort.PooledSlab[GuardedCallback_Context], _ AbstractCompiler, parent pkg_rhino.Node, child pkg_rhino.Node) *GuardedCallback_Context {
	if parent.IsIf() {
		if child == parent.GetFirstChild() {
			return this_.linkWithSlab_(slab_, parent, true)
		}
		return StaticGet_GuardedCallback_Context_EMPTY()
	}
	if parent.IsAnd() {
		if child == parent.GetFirstChild() {
			return this_.linkWithSlab_(slab_, parent, true)
		}
		return this_
	}
	if parent.IsScript() || parent.GetToken() == pkg_rhino.TokROOT || parent.GetToken() == pkg_rhino.TokBLOCK {
		return StaticGet_GuardedCallback_Context_EMPTY()
	}
	return this_.propagateWithSlab_(slab_, false)
}

type CustomGuardedVisitor struct {
	GuardedCallbackBase
	SafeCount int
}

func NewCustomGuardedVisitor(c AbstractCompiler) *CustomGuardedVisitor {
	v := &CustomGuardedVisitor{}
	v.BindSelf_(v)
	InitGuardedCallback(&v.GuardedCallbackBase, c)
	return v
}
func (v *CustomGuardedVisitor) VisitGuarded(_ *NodeTraversal, n pkg_rhino.Node, _ pkg_rhino.Node) {
	if n.IsName() && v.contextStack.Peek().safe {
		v.SafeCount++
	}
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "guarded_callback.go"), []byte(guardedSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1C-1: MinimizedCondition + PeepholeMinimizeConditions.
	minCondSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w31btest/gen/j2gort"
	pkg_rhino "example.com/w31btest/gen/rhino"
)

func NodeUtil_Precedence(tok pkg_rhino.Token) int32 {
	switch tok {
	case pkg_rhino.TokOR:
		return 4
	case pkg_rhino.TokAND:
		return 5
	case pkg_rhino.TokNOT:
		return 13
	default:
		return 15
	}
}

type MinimizedCondition_MinimizationStyle uint8

const (
	stylePreferUnnegated MinimizedCondition_MinimizationStyle = 1
	styleAllowLeadingNot MinimizedCondition_MinimizationStyle = 2
)

func StaticGet_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED() MinimizedCondition_MinimizationStyle {
	return stylePreferUnnegated
}
func StaticGet_MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT() MinimizedCondition_MinimizationStyle {
	return styleAllowLeadingNot
}

type MinimizedCondition struct {
	positive *MinimizedCondition_MeasuredNode
	negative *MinimizedCondition_MeasuredNode
	idHash_  atomic.Int32
}

var static_MinimizedCondition_UNOPTIMIZED_NEG *MinimizedCondition_MeasuredNode
var clinitDone_MC atomic.Bool

func clinit_MinimizedCondition() {
	if clinitDone_MC.CompareAndSwap(false, true) {
		clinitSlow_MinimizedCondition()
	}
}
func clinitSlow_MinimizedCondition() {
	static_MinimizedCondition_UNOPTIMIZED_NEG = NewMinimizedCondition_MeasuredNode(0, -1, nil, nil, nil, 2147483647, true)
}
func StaticGet_MinimizedCondition_UNOPTIMIZED_NEG() *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition()
	return static_MinimizedCondition_UNOPTIMIZED_NEG
}

func NewMinimizedCondition(p *MinimizedCondition_MeasuredNode, n *MinimizedCondition_MeasuredNode) *MinimizedCondition {
	clinit_MinimizedCondition()
	this_ := &MinimizedCondition{}
	InitMinimizedCondition(this_, p, n)
	return this_
}
func InitMinimizedCondition(this_ *MinimizedCondition, p *MinimizedCondition_MeasuredNode, n *MinimizedCondition_MeasuredNode) {
	_ = p
	_ = n
	this_.positive = p
	this_.negative = n.priv_MeasuredNode_change_3801()
}

func MinimizedCondition_IsMinimizable(n pkg_rhino.Node) bool {
	switch n.GetToken() {
	case pkg_rhino.TokNOT, pkg_rhino.TokAND, pkg_rhino.TokOR:
		return true
	default:
		return true
	}
}

func MinimizedCondition_FromConditionNode(n pkg_rhino.Node) *MinimizedCondition {
	clinit_MinimizedCondition()
	_ = n
	if !n.HasParent() {
		j2gort.CheckState(false)
	}
	switch n.GetToken() {
	case pkg_rhino.TokNOT, pkg_rhino.TokAND, pkg_rhino.TokOR:
		return MinimizedCondition_ComputeMinimizedCondition(n)
	default:
		return MinimizedCondition_Unoptimized(n)
	}
}

func (this_ *MinimizedCondition) GetMinimized(style MinimizedCondition_MinimizationStyle) *MinimizedCondition_MeasuredNode {
	_ = style
	if style == StaticGet_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED() || this_.positive.node.IsNot() || this_.positive.length <= this_.negative.length {
		return this_.positive
	}
	return this_.negative.priv_MeasuredNode_addNot_3800()
}

func MinimizedCondition_Unoptimized(n pkg_rhino.Node) *MinimizedCondition {
	clinit_MinimizedCondition()
	_ = n
	j2gort.CheckNotNullNode(n.GetParent())
	var val_1 MinimizedCondition_MeasuredNode
	InitMinimizedCondition_MeasuredNode(&val_1, n, -1, nil, nil, nil, 0, false)
	var pos *MinimizedCondition_MeasuredNode = &val_1
	return NewMinimizedCondition(pos, StaticGet_MinimizedCondition_UNOPTIMIZED_NEG())
}

func MinimizedCondition_PickBest(a, b *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	if a.length <= b.length {
		return a
	}
	return b
}

func MinimizedCondition_ComputeMinimizedCondition(n pkg_rhino.Node) *MinimizedCondition {
	clinit_MinimizedCondition()
	switch n.GetToken() {
	case pkg_rhino.TokNOT:
		sub := MinimizedCondition_ComputeMinimizedCondition(n.GetFirstChild())
		pos := MinimizedCondition_PickBest(MinimizedCondition_MeasuredNode_AddNode_2(n, sub.positive), sub.negative)
		neg := MinimizedCondition_PickBest(sub.negative.priv_MeasuredNode_negate_3798(), sub.positive)
		return NewMinimizedCondition(pos, neg)
	case pkg_rhino.TokAND, pkg_rhino.TokOR:
		compTok := pkg_rhino.StaticGet_Token_AND()
		if n.IsAnd() {
			compTok = pkg_rhino.StaticGet_Token_OR()
		}
		compNode := pkg_rhino.NewNode_1(compTok).Srcref(n)
		left := MinimizedCondition_ComputeMinimizedCondition(n.GetFirstChild())
		right := MinimizedCondition_ComputeMinimizedCondition(n.GetLastChild())
		posTree := MinimizedCondition_MeasuredNode_AddNode_3(n, left.positive, right.positive)
		compTree := MinimizedCondition_MeasuredNode_AddNode_3(compNode, left.negative, right.negative)
		pos := MinimizedCondition_PickBest(posTree, compTree.priv_MeasuredNode_negate_3798())
		neg := MinimizedCondition_PickBest(posTree.priv_MeasuredNode_negate_3798(), compTree.priv_MeasuredNode_change_3801())
		return NewMinimizedCondition(pos, neg)
	default:
		pos := MinimizedCondition_MeasuredNode_ForNode(n)
		neg := pos.priv_MeasuredNode_negate_3798()
		return NewMinimizedCondition(pos, neg)
	}
}

type MinimizedCondition_MeasuredNode struct {
	c0         *MinimizedCondition_MeasuredNode
	c1         *MinimizedCondition_MeasuredNode
	c2         *MinimizedCondition_MeasuredNode
	node       pkg_rhino.Node
	length     int32
	childCount int32
	idHash_    atomic.Int32
	changed    bool
}

func NewMinimizedCondition_MeasuredNode(n pkg_rhino.Node, childCount int32, c0, c1, c2 *MinimizedCondition_MeasuredNode, len_ int32, ch bool) *MinimizedCondition_MeasuredNode {
	this_ := &MinimizedCondition_MeasuredNode{}
	InitMinimizedCondition_MeasuredNode(this_, n, childCount, c0, c1, c2, len_, ch)
	return this_
}
func InitMinimizedCondition_MeasuredNode(this_ *MinimizedCondition_MeasuredNode, n pkg_rhino.Node, childCount int32, c0, c1, c2 *MinimizedCondition_MeasuredNode, len_ int32, ch bool) {
	this_.node = n
	this_.childCount = childCount
	this_.c0 = c0
	this_.c1 = c1
	this_.c2 = c2
	this_.length = len_
	this_.changed = ch
}
func (this_ *MinimizedCondition_MeasuredNode) IsChanged() bool { return this_.changed }
func (this_ *MinimizedCondition_MeasuredNode) IsNot() bool     { return this_.node.IsNot() }
func (this_ *MinimizedCondition_MeasuredNode) WithoutNot() *MinimizedCondition_MeasuredNode {
	var first *MinimizedCondition_MeasuredNode
	if this_.childCount >= 0 {
		first = this_.c0
	} else {
		first = MinimizedCondition_MeasuredNode_ForNode(this_.node.GetFirstChild())
	}
	return first.priv_MeasuredNode_change_3801()
}
func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_negate_3798() *MinimizedCondition_MeasuredNode {
	if this_.node.IsNot() {
		return this_.WithoutNot()
	}
	return this_.priv_MeasuredNode_addNot_3800()
}
func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_updateToken_3799(tok pkg_rhino.Token) *MinimizedCondition_MeasuredNode {
	newNode := pkg_rhino.NewNode_1(tok).Srcref(this_.node)
	return NewMinimizedCondition_MeasuredNode(newNode, this_.childCount, this_.c0, this_.c1, this_.c2, this_.length, true)
}
func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_addNot_3800() *MinimizedCondition_MeasuredNode {
	cost := this_.length + 1
	if NodeUtil_Precedence(this_.node.GetToken()) < NodeUtil_Precedence(pkg_rhino.TokNOT) {
		cost += 2
	}
	return NewMinimizedCondition_MeasuredNode(pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_NOT()).Srcref(this_.node), 1, this_, nil, nil, cost, true)
}
func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_change_3801() *MinimizedCondition_MeasuredNode {
	if this_.IsChanged() {
		return this_
	}
	return NewMinimizedCondition_MeasuredNode(this_.node, this_.childCount, this_.c0, this_.c1, this_.c2, this_.length, true)
}
func MinimizedCondition_MeasuredNode_AddNode_2(parent pkg_rhino.Node, c0 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	cost := c0.length
	if parent.IsNot() {
		cost++
	}
	return NewMinimizedCondition_MeasuredNode(parent, 1, c0, nil, nil, cost, c0.changed)
}
func MinimizedCondition_MeasuredNode_AddNode_3(parent pkg_rhino.Node, c0, c1 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	cost := c0.length + c1.length
	return NewMinimizedCondition_MeasuredNode(parent, 2, c0, c1, nil, cost, c0.changed || c1.changed)
}
func MinimizedCondition_MeasuredNode_AddNode_4(parent pkg_rhino.Node, c0, c1, c2 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	cost := c0.length + c1.length + c2.length
	return NewMinimizedCondition_MeasuredNode(parent, 3, c0, c1, c2, cost, c0.changed || c1.changed || c2.changed)
}
func MinimizedCondition_MeasuredNode_ForNode(n pkg_rhino.Node) *MinimizedCondition_MeasuredNode {
	return NewMinimizedCondition_MeasuredNode(n, -1, nil, nil, nil, 0, false)
}
func (this_ *MinimizedCondition_MeasuredNode) WillChange(original pkg_rhino.Node) bool {
	return original != this_.node || this_.IsChanged()
}
func (this_ *MinimizedCondition_MeasuredNode) BuildReplacement() pkg_rhino.Node {
	if this_.childCount >= 0 {
		this_.node.DetachChildren()
		if this_.childCount >= 1 {
			this_.node.AddChildToBack(this_.c0.BuildReplacement().Detach())
		}
		if this_.childCount >= 2 {
			this_.node.AddChildToBack(this_.c1.BuildReplacement().Detach())
		}
	}
	return this_.node
}
func (this_ *MinimizedCondition_MeasuredNode) ApplyTo(original pkg_rhino.Node) pkg_rhino.Node {
	repl := this_.BuildReplacement()
	if original != repl {
		repl.Detach()
		original.ReplaceWith(repl)
	}
	return repl
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "minimized_condition.go"), []byte(minCondSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	pmcSrc := `package jscomp

import pkg_rhino "example.com/w31btest/gen/rhino"

type AbstractPeepholeOptimization interface {
	OptimizeSubtree(pkg_rhino.Node) pkg_rhino.Node
	EndTraversal()
}

type AbstractPeepholeOptimizationBase struct {
	self_    AbstractPeepholeOptimization
	compiler AbstractCompiler
}

func (b *AbstractPeepholeOptimizationBase) BindSelf_(s any) {
	b.self_ = s.(AbstractPeepholeOptimization)
}
func (b *AbstractPeepholeOptimizationBase) EndTraversal() {
	b.compiler = nil
}

type PeepholeMinimizeConditions struct {
	AbstractPeepholeOptimizationBase
	late bool
}

func NewPeepholeMinimizeConditions(late bool) *PeepholeMinimizeConditions {
	p := &PeepholeMinimizeConditions{late: late}
	p.BindSelf_(p)
	return p
}

func (this_ *PeepholeMinimizeConditions) OptimizeSubtree(node pkg_rhino.Node) pkg_rhino.Node {
	if node.IsIf() {
		return this_.priv_PeepholeMinimizeConditions_tryMinimizeIf_4126(node)
	}
	return node
}

func (this_ *PeepholeMinimizeConditions) priv_PeepholeMinimizeConditions_tryMinimizeIf_4126(n pkg_rhino.Node) pkg_rhino.Node {
	origCond := n.GetFirstChild()
	if MinimizedCondition_IsMinimizable(origCond) {
		minCond := MinimizedCondition_FromConditionNode(origCond)
		shortCond := minCond.GetMinimized(StaticGet_MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT())
		if shortCond.IsNot() {
			unneg := shortCond.WithoutNot()
			if unneg.WillChange(origCond) {
				unneg.ApplyTo(origCond)
			}
		} else if shortCond.WillChange(origCond) {
			shortCond.ApplyTo(origCond)
		}
	}
	return n
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "peephole_minimize_conditions.go"), []byte(pmcSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1C-2: OptimizeCalls + PureFunctionIdentifier + OptimizeParameters.
	optCallsSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w31btest/gen/j2gort"
	pkg_rhino "example.com/w31btest/gen/rhino"
)

type OptimizeCalls_CallGraphCompilerPass interface {
	Process_3(pkg_rhino.Node, pkg_rhino.Node, *OptimizeCalls_ReferenceMap)
}

type OptimizeCalls struct {
	compiler AbstractCompiler
	passes   []OptimizeCalls_CallGraphCompilerPass
}

func (this_ *OptimizeCalls) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	if len(this_.passes) == 0 {
		return
	}
	var val_1 OptimizeCalls_ReferenceMap
	InitOptimizeCalls_ReferenceMap(&val_1)
	var references *OptimizeCalls_ReferenceMap = &val_1
	_ = references
	defer func() {
		references.globalScope = nil
	}()
	references.AddNameReference("foo", root)
	references.AddPropReference("foo", externs)
	for _, pass := range this_.passes {
		pass.Process_3(externs, root, references)
	}
}

type OptimizeCalls_NodeList struct {
	nodes   []pkg_rhino.Node
	idHash_ atomic.Int32
}

func NewOptimizeCalls_NodeList() *OptimizeCalls_NodeList {
	this_ := &OptimizeCalls_NodeList{}
	InitOptimizeCalls_NodeList(this_)
	return this_
}
func InitOptimizeCalls_NodeList(this_ *OptimizeCalls_NodeList) {
	this_.nodes = []pkg_rhino.Node{}
}
func (this_ *OptimizeCalls_NodeList) Add(n pkg_rhino.Node) {
	this_.nodes = append(this_.nodes, n)
}
func (this_ *OptimizeCalls_NodeList) AddAll(other *OptimizeCalls_NodeList) {
	this_.nodes = append(this_.nodes, other.nodes...)
}
func (this_ *OptimizeCalls_NodeList) Size() int32   { return int32(len(this_.nodes)) }
func (this_ *OptimizeCalls_NodeList) IsEmpty() bool { return len(this_.nodes) == 0 }

type OptimizeCalls_ReferenceMap struct {
	globalScope any
	names       j2gort.OrderedMap[string, *OptimizeCalls_NodeList]
	props       j2gort.OrderedMap[string, *OptimizeCalls_NodeList]
	idHash_     atomic.Int32
}

func InitOptimizeCalls_ReferenceMap(this_ *OptimizeCalls_ReferenceMap) {
	this_.names = j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList]()
	this_.props = j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList]()
}
func (this_ *OptimizeCalls_ReferenceMap) AddNameReference(name string, n pkg_rhino.Node) {
	var refs *OptimizeCalls_NodeList = this_.names.Get(name)
	if refs == nil {
		refs = NewOptimizeCalls_NodeList()
		this_.names.Put(name, refs)
	}
	refs.Add(n)
}
func (this_ *OptimizeCalls_ReferenceMap) AddPropReference(name string, n pkg_rhino.Node) {
	var refs *OptimizeCalls_NodeList = this_.props.Get(name)
	if refs == nil {
		refs = NewOptimizeCalls_NodeList()
		this_.props.Put(name, refs)
	}
	refs.Add(n)
}

type PureFunctionIdentifier struct {
	MergedCount int32
}

func (this_ *PureFunctionIdentifier) Process_3(_, _ pkg_rhino.Node, referenceMap *OptimizeCalls_ReferenceMap) {
	this_.priv_PureFunctionIdentifier_populateDatastructuresForAnalysisTraversal_4485(referenceMap)
}
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_populateDatastructuresForAnalysisTraversal_4485(referenceMap *OptimizeCalls_ReferenceMap) {
	existing := referenceMap.names.Get("foo")
	propRefs := referenceMap.props.Get("foo")
	if existing != nil && propRefs != nil {
		var merged *OptimizeCalls_NodeList = NewOptimizeCalls_NodeList()
		merged.AddAll(existing)
		merged.AddAll(propRefs)
		this_.MergedCount = merged.Size()
	}
}

func OptimizeParameters_MayReferenceParamBeforeBody(fn pkg_rhino.Node) bool {
	var namesByNames j2gort.OrderedMap[string, *OptimizeCalls_NodeList] = j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList]()
	_ = namesByNames
	for c := fn.GetFirstChild(); c != 0; c = c.GetNext() {
		if c.IsName() {
			list := namesByNames.Get(c.GetString())
			if list == nil {
				list = NewOptimizeCalls_NodeList()
				namesByNames.Put(c.GetString(), list)
			}
			list.Add(c)
			if list.Size() > 1 {
				return true
			}
		}
	}
	return false
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "optimize_calls.go"), []byte(optCallsSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1C-3: MakeDeclaredNamesUnique + ContextualRenamerBase + BoilerplateRenamer + Normalize.
	mdnuSrc := `package jscomp

import (
	"fmt"
	"sync/atomic"

	"example.com/w31btest/gen/j2gort"
	pkg_rhino "example.com/w31btest/gen/rhino"
)

type ScopeTestNode struct {
	Root   pkg_rhino.Node
	Global bool
	Vars   []string
}

func (t *NodeTraversal) GetScopeRoot() pkg_rhino.Node { return curScopeNode.Root }
func (t *NodeTraversal) InGlobalScope() bool          { return curScopeNode.Global }

var curScopeNode ScopeTestNode

type MakeDeclaredNamesUnique_Renamer interface {
	AddDeclaredName(string, bool)
	GetReplacementName(string) string
	CreateForChildScope(pkg_rhino.Node, bool) MakeDeclaredNamesUnique_Renamer
	GetHoistRenamer() MakeDeclaredNamesUnique_Renamer
}

type MakeDeclaredNamesUnique_ContextualRenamer interface {
	MakeDeclaredNamesUnique_Renamer
	Implements_ContextualRenamer()
}

type MakeDeclaredNamesUnique struct {
	renamerStack   j2gort.Deque[MakeDeclaredNamesUnique_Renamer]
	rootRenamer    MakeDeclaredNamesUnique_Renamer
	assertOnChange bool
}

func NewMakeDeclaredNamesUnique(renamer MakeDeclaredNamesUnique_Renamer, assertOnChange bool) *MakeDeclaredNamesUnique {
	return &MakeDeclaredNamesUnique{rootRenamer: renamer, assertOnChange: assertOnChange}
}

func (this_ *MakeDeclaredNamesUnique) EnterScope(t *NodeTraversal) {
	_ = t
	var declarationRoot pkg_rhino.Node = t.GetScopeRoot()
	var renamer MakeDeclaredNamesUnique_Renamer = nil
	if this_.renamerStack.IsEmpty() {
		renamer = this_.rootRenamer
	} else {
		var hoist bool = !declarationRoot.IsFunction()
		renamer = this_.renamerStack.Peek().CreateForChildScope(t.GetScopeRoot(), hoist)
	}
	this_.renamerStack.Push(renamer)
	for _, v := range curScopeNode.Vars {
		this_.renamerStack.Peek().AddDeclaredName(v, false)
	}
}

func (this_ *MakeDeclaredNamesUnique) ExitScope(t *NodeTraversal) {
	_ = t
	if !(t.InGlobalScope()) {
		this_.renamerStack.Pop()
	}
}

func (this_ *MakeDeclaredNamesUnique) Resolve(name string) string {
	return this_.renamerStack.Peek().GetReplacementName(name)
}

type MakeDeclaredNamesUnique_Builder struct {
	renamer        MakeDeclaredNamesUnique_Renamer
	assertOnChange bool
}

func MakeDeclaredNamesUnique_S_Builder() *MakeDeclaredNamesUnique_Builder {
	return &MakeDeclaredNamesUnique_Builder{}
}
func (b *MakeDeclaredNamesUnique_Builder) WithRenamer(r MakeDeclaredNamesUnique_Renamer) *MakeDeclaredNamesUnique_Builder {
	b.renamer = r
	return b
}
func (b *MakeDeclaredNamesUnique_Builder) WithAssertOnChange(a bool) *MakeDeclaredNamesUnique_Builder {
	b.assertOnChange = a
	return b
}
func (b *MakeDeclaredNamesUnique_Builder) Build() *MakeDeclaredNamesUnique {
	if b.renamer == nil {
		b.renamer = NewMakeDeclaredNamesUnique_ContextualRenamer()
	}
	return NewMakeDeclaredNamesUnique(b.renamer, b.assertOnChange)
}

type MakeDeclaredNamesUnique_ContextualRenamer_NameUsage struct {
	counts map[string]int32
}

type MakeDeclaredNamesUnique_ContextualRenamerBase struct {
	self_        MakeDeclaredNamesUnique_ContextualRenamer
	nameUsage    *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage
	declarations j2gort.OrderedMap[string, string]
	hoistRenamer MakeDeclaredNamesUnique_Renamer
	scopeRoot    pkg_rhino.Node
	idHash_      atomic.Int32
	global       bool
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) BindSelf_(s any) {
	this_.self_ = s.(MakeDeclaredNamesUnique_ContextualRenamer)
}
func (*MakeDeclaredNamesUnique_ContextualRenamerBase) Implements_ContextualRenamer() {}

func NewMakeDeclaredNamesUnique_ContextualRenamer() MakeDeclaredNamesUnique_ContextualRenamer {
	this_ := &MakeDeclaredNamesUnique_ContextualRenamerBase{
		global:       true,
		nameUsage:    &MakeDeclaredNamesUnique_ContextualRenamer_NameUsage{counts: make(map[string]int32)},
		declarations: j2gort.NewOrderedMap[string, string](),
	}
	this_.BindSelf_(this_)
	this_.hoistRenamer = this_
	return this_
}

func NewMakeDeclaredNamesUnique_ContextualRenamer_4(scopeRoot pkg_rhino.Node, nameUsage *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage, hoistingTargetScope bool, parent MakeDeclaredNamesUnique_Renamer) MakeDeclaredNamesUnique_ContextualRenamer {
	this_ := &MakeDeclaredNamesUnique_ContextualRenamerBase{}
	this_.BindSelf_(this_)
	InitMakeDeclaredNamesUnique_ContextualRenamer_4(this_, scopeRoot, nameUsage, hoistingTargetScope, parent)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenamer_4(this_ *MakeDeclaredNamesUnique_ContextualRenamerBase, scopeRoot pkg_rhino.Node, nameUsage *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage, hoistingTargetScope bool, parent MakeDeclaredNamesUnique_Renamer) {
	this_.declarations = j2gort.NewOrderedMap[string, string]()
	this_.scopeRoot = scopeRoot
	this_.global = false
	this_.nameUsage = nameUsage
	if hoistingTargetScope {
		this_.hoistRenamer = this_.self_
	} else {
		this_.hoistRenamer = parent.GetHoistRenamer()
	}
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) CreateForChildScope(scopeRoot pkg_rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	return j2gort.AsIface[MakeDeclaredNamesUnique_Renamer](NewMakeDeclaredNamesUnique_ContextualRenamer_4(scopeRoot, this_.nameUsage, hoistingTargetScope, j2gort.AsIface[MakeDeclaredNamesUnique_Renamer](this_.self_)))
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) AddDeclaredName(name string, _ bool) {
	if this_.global {
		if _, ok := this_.nameUsage.counts[name]; !ok {
			this_.nameUsage.counts[name] = 1
		}
		return
	}
	if !this_.declarations.ContainsKey(name) {
		prev := this_.nameUsage.counts[name]
		this_.nameUsage.counts[name] = prev + 1
		if prev > 0 {
			this_.declarations.Put(name, fmt.Sprintf("%s$jscomp$%d", name, prev))
		} else {
			this_.declarations.Put(name, "")
		}
	}
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) GetReplacementName(oldName string) string {
	return j2gort.OrderedMapGetStr(this_.declarations, oldName)
}
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return this_.hoistRenamer
}

type MakeDeclaredNamesUnique_InlineRenamer struct {
	prefix string
}

func (r *MakeDeclaredNamesUnique_InlineRenamer) AddDeclaredName(string, bool) {}
func (r *MakeDeclaredNamesUnique_InlineRenamer) GetReplacementName(s string) string {
	return s + "$inline$" + r.prefix
}
func (r *MakeDeclaredNamesUnique_InlineRenamer) CreateForChildScope(_ pkg_rhino.Node, _ bool) MakeDeclaredNamesUnique_Renamer {
	return &MakeDeclaredNamesUnique_InlineRenamer{prefix: r.prefix}
}
func (r *MakeDeclaredNamesUnique_InlineRenamer) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return r
}

type MakeDeclaredNamesUnique_BoilerplateRenamer struct {
	MakeDeclaredNamesUnique_ContextualRenamerBase
	idPrefix string
}

func NewMakeDeclaredNamesUnique_BoilerplateRenamer(prefix string) *MakeDeclaredNamesUnique_BoilerplateRenamer {
	b := &MakeDeclaredNamesUnique_BoilerplateRenamer{idPrefix: prefix}
	b.nameUsage = &MakeDeclaredNamesUnique_ContextualRenamer_NameUsage{counts: make(map[string]int32)}
	b.declarations = j2gort.NewOrderedMap[string, string]()
	b.global = true
	b.BindSelf_(b)
	b.hoistRenamer = b
	return b
}
func (this_ *MakeDeclaredNamesUnique_BoilerplateRenamer) CreateForChildScope(_ pkg_rhino.Node, _ bool) MakeDeclaredNamesUnique_Renamer {
	return &MakeDeclaredNamesUnique_InlineRenamer{prefix: this_.idPrefix}
}

type Normalize struct {
	assertOnChange bool
	LastRenamed    string
}

func (this_ *Normalize) Process_2(_ pkg_rhino.Node, root pkg_rhino.Node) {
	var renamer *MakeDeclaredNamesUnique = MakeDeclaredNamesUnique_S_Builder().WithAssertOnChange(this_.assertOnChange).Build()
	_ = renamer
	nt := &NodeTraversal{}
	curScopeNode = ScopeTestNode{Root: root, Global: true, Vars: []string{"x"}}
	renamer.EnterScope(nt)
	fnNode := pkg_rhino.AllocNamed(pkg_rhino.TokFUNCTION, "f")
	curScopeNode = ScopeTestNode{Root: fnNode, Global: false, Vars: []string{"x"}}
	renamer.EnterScope(nt)
	this_.LastRenamed = renamer.Resolve("x")
	renamer.ExitScope(nt)
	curScopeNode = ScopeTestNode{Root: root, Global: true}
	renamer.ExitScope(nt)
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "make_declared_names_unique.go"), []byte(mdnuSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 2: CompilerDecisionsLog + 1-hop priv_* Supplier logging wrappers.
	logTestSrc := `package jscomp

import (
	pkg_diagnostic "example.com/w31btest/gen/diagnostic"
	"example.com/w31btest/gen/j2gort"
)

type CompilerDecisionsLog struct {
	enabled bool
	entries []string
}

func (l *CompilerDecisionsLog) IsJavaNil_() bool { return l == nil }
func (l *CompilerDecisionsLog) IsLogging() bool  { return l != nil && l.enabled }
func (l *CompilerDecisionsLog) Log_1_Ext_java_util_function_Supplier(s j2gort.Ext_java_util_function_Supplier) {
	if l.IsLogging() {
		l.entries = append(l.entries, j2gort.CastString(s.Get()))
	}
}

type DirectDecisionPass struct {
	decisionsLog *CompilerDecisionsLog
}

func (this_ *DirectDecisionPass) Run(name string) {
	this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxString("direct:" + name)
	}))
}

type CollapseOuter struct {
	decisionsLog pkg_diagnostic.LogFile
}

type CollapsePropsHelper struct {
	outer_ *CollapseOuter
}

func (this_ *CollapsePropsHelper) priv_CollapsePropsHelper_logDecisionForName_1(name string, msgSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = name
	_ = msgSupplier
	if this_.outer_.decisionsLog.IsLogging() {
		this_.outer_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return j2gort.BoxString(name + ": " + j2gort.CastString(msgSupplier.Get()))
		}))
	}
}

func (this_ *CollapsePropsHelper) CheckName(name string, reason string) {
	this_.priv_CollapsePropsHelper_logDecisionForName_1(name, j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxString("skip " + reason)
	}))
}

type GlobalNsOuter struct {
	decisionsLog pkg_diagnostic.LogFile
}

type GlobalNsName struct {
	outer_ *GlobalNsOuter
	full   string
}

type Inlinability uint8

func StaticGet_Inlinability_DO_NOT_INLINE() Inlinability { return 1 }

func (this_ *GlobalNsName) priv_Name_logChildNamesDecision_2(inl Inlinability, reasonSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = inl
	_ = reasonSupplier
	if (this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging() {
		this_.outer_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return j2gort.BoxString(this_.full + ":" + j2gort.CastString(reasonSupplier.Get()))
		}))
	}
}

func (this_ *GlobalNsName) CheckChildren(sets string) {
	this_.priv_Name_logChildNamesDecision_2(StaticGet_Inlinability_DO_NOT_INLINE(), j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxString("sets=" + sets)
	}))
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "decisions_log.go"), []byte(logTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// 3: 1-arg NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection in FunctionArgumentInjector + DeadPropertyAssignmentElimination.
	setTestSrc := `package jscomp

import (
	"example.com/w31btest/gen/j2gort"
	pkg_rhino "example.com/w31btest/gen/rhino"
)

type FunctionArgumentInjector struct{}

func (this_ *FunctionArgumentInjector) priv_FunctionArgumentInjector_findParametersReferencedAfterSideEffect_2877(parameters j2gort.Ext_com_google_common_collect_ImmutableSet, root pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = parameters
	_ = root
	var locals j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](parameters))
	_ = locals
	FunctionArgumentInjector_GatherLocalNames(root, locals)
	return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(locals)
}

func FunctionArgumentInjector_GatherLocalNames(n pkg_rhino.Node, names j2gort.Ext_java_util_Set) {
	_ = n
	_ = names
	if n.IsFunction() {
		return
	}
	if n.IsName() && n.GetParent().GetToken() == pkg_rhino.TokVAR {
		j2gort.SetAddStr(names, n.GetString())
	}
	for c := n.GetFirstChild(); c != 0; c = c.GetNext() {
		FunctionArgumentInjector_GatherLocalNames(c, names)
	}
}

type DeadPropertyAssignmentElimination_Property struct {
	name     string
	read     bool
	children j2gort.Ext_java_util_Set
}

func NewDeadPropertyAssignmentElimination_Property(name string) *DeadPropertyAssignmentElimination_Property {
	this_ := &DeadPropertyAssignmentElimination_Property{}
	InitDeadPropertyAssignmentElimination_Property(this_, name)
	return this_
}

func InitDeadPropertyAssignmentElimination_Property(this_ *DeadPropertyAssignmentElimination_Property, name string) {
	this_.name = name
	this_.children = j2gort.NewExt_java_util_LinkedHashSet()
}

func (this_ *DeadPropertyAssignmentElimination_Property) AddChild(c *DeadPropertyAssignmentElimination_Property) {
	this_.children.Add_1(j2gort.PtrToAny(c))
}

func (this_ *DeadPropertyAssignmentElimination_Property) MarkChildrenRead() {
	var propertiesSet j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.children))
	_ = propertiesSet
	var propertyQueue j2gort.Deque[*DeadPropertyAssignmentElimination_Property] = j2gort.DequeFromCol[*DeadPropertyAssignmentElimination_Property](propertiesSet)
	_ = propertyQueue
	propertiesSet.Add_1(this_)
	for !(propertyQueue.IsEmpty()) {
		var childProperty *DeadPropertyAssignmentElimination_Property = propertyQueue.RemoveFirst()
		childProperty.read = true
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](childProperty.children)))
		for cur_5.HasNext() {
			rawElem_4 := cur_5.Next()
			grandchildProperty := j2gort.CastPtr[DeadPropertyAssignmentElimination_Property](rawElem_4)
			if propertiesSet.Add_1(j2gort.PtrToAny(grandchildProperty)) {
				propertyQueue.AddLast(grandchildProperty)
			}
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "set_lowering_targets.go"), []byte(setTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	targetPasses := []passes.Pass{
		passes.NewNativeSetLoweringPass(),
		passes.NewJ2gortPeepholePass(),
		passes.NewSlabAllocLoweringPass(),
	}

	stats1, err := passes.RunOnPackageDir(jscompDir, targetPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Wave 31B Run 1 failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] < 15 || stats1.ByPass["j2gort_peephole"] < 3 || stats1.ByPass["native_set_lowering"] < 10 {
		t.Fatalf("unexpected Wave 31B Run 1 stats: %+v", stats1.ByPass)
	}

	// Run 2 MUST be 100% idempotent (0 sites, 0 modified files)!
	stats2, err := passes.RunOnPackageDir(jscompDir, targetPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Wave 31B Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("Wave 31B Run 2 NOT IDEMPOTENT: got %+v", stats2)
	}

	readLowered := func(name string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(jscompDir, name))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	// Verify 1A phase boundaries:
	pbCode := readLowered("phase_boundaries.go")
	if strings.Contains(pbCode, "TrimStaleSlabPools(") {
		t.Fatalf("TrimStaleSlabPools not stripped:\n%s", pbCode)
	}
	if got := strings.Count(pbCode, "PurgeAllRuntimePools()"); got != 3 {
		t.Fatalf("expected exactly 3 PurgeAllRuntimePools() calls across the 3 phase boundaries (InferAllScopes entry, SerializeRoots entry, PerformTranspilationAndOptimizations defer), got %d:\n%s", got, pbCode)
	}

	// Verify 3 set lowering eliminated 100% of LinkedHashSet_1_Ext_java_util_Collection in findParametersReferencedAfterSideEffect & MarkChildrenRead
	// while keeping the original exported FunctionArgumentInjector_GatherLocalNames intact for external callers:
	setCode := readLowered("set_lowering_targets.go")
	for _, want := range []string{
		"j2gort.OrderedSetFromStrCol(",
		"ordSet_FunctionArgumentInjector_GatherLocalNames(root, &locals)",
		"func FunctionArgumentInjector_GatherLocalNames(n pkg_rhino.Node, names j2gort.Ext_java_util_Set)",
		"j2gort.OrderedSetFromSet(this_.children)",
		"j2gort.DequeFromOrderedSet(propertiesSet)",
	} {
		if !strings.Contains(setCode, want) {
			t.Fatalf("missing %q in lowered set_lowering_targets.go:\n%s", want, setCode)
		}
	}

	// End-to-end runtime + -race test!
	rtSrc := `package jscomp_test

import (
	"testing"

	"example.com/w31btest/gen/j2gort"
	"example.com/w31btest/gen/jscomp"
	pkg_rhino "example.com/w31btest/gen/rhino"
)

func TestWave31BRuntimeSemantics(t *testing.T) {
	// 1A: Phase boundaries (3 exact sites: InferAllScopes entry, SerializeRoots entry, PerformTranspilationAndOptimizations defer)
	j2gort.ResetWave31BMetrics()
	c := &jscomp.Compiler{}
	_ = c.ParseInputs()
	if j2gort.PurgeCalls.Load() != 0 {
		t.Fatalf("ParseInputs must NOT purge runtime pools, got %d", j2gort.PurgeCalls.Load())
	}
	tip := &jscomp.TypeInferencePass{}
	_ = tip.InferAllScopes(1)
	if j2gort.PurgeCalls.Load() != 1 {
		t.Fatalf("after InferAllScopes (entry only), purge count = %d, want 1", j2gort.PurgeCalls.Load())
	}
	tas := &jscomp.TypedAstSerializer{}
	_ = tas.SerializeRoots(1, 1)
	if j2gort.PurgeCalls.Load() != 2 {
		t.Fatalf("after SerializeRoots (entry), purge count = %d, want 2", j2gort.PurgeCalls.Load())
	}
	c.PerformTranspilationAndOptimizations(0)
	if j2gort.PurgeCalls.Load() != 3 {
		t.Fatalf("after PerformTranspilationAndOptimizations (defer), purge count = %d, want 3", j2gort.PurgeCalls.Load())
	}
	vrc := &jscomp.VariableReferenceCheck{}
	vrc.Process_2(1, 1)
	if j2gort.PurgeCalls.Load() != 3 {
		t.Fatalf("VariableReferenceCheck must NOT purge runtime pools, got %d", j2gort.PurgeCalls.Load())
	}

	// 1B: GuardedCallbackBase whole-AST Token.ROOT traversal across multiple scripts
	pkg_rhino.ResetArena()
	root := pkg_rhino.NewNode_1(pkg_rhino.TokROOT)
	for s := 0; s < 4; s++ {
		scr := pkg_rhino.NewNode_1(pkg_rhino.TokSCRIPT)
		for k := 0; k < 3; k++ {
			ifN := pkg_rhino.NewNode_1(pkg_rhino.TokIF)
			andN := pkg_rhino.NewNode_1(pkg_rhino.TokAND)
			andN.AddChildToBack(pkg_rhino.AllocNamed(pkg_rhino.TokNAME, "guard"))
			andN.AddChildToBack(pkg_rhino.AllocNamed(pkg_rhino.TokNAME, "guardedUse"))
			ifN.AddChildToBack(andN)
			blk := pkg_rhino.NewNode_1(pkg_rhino.TokBLOCK)
			blk.AddChildToBack(pkg_rhino.AllocNamed(pkg_rhino.TokNAME, "uncondBody"))
			ifN.AddChildToBack(blk)
			scr.AddChildToBack(ifN)
		}
		root.AddChildToBack(scr)
	}
	j2gort.ResetWave31BMetrics()
	gv := jscomp.NewCustomGuardedVisitor(c)
	jscomp.NodeTraversal_Traverse(c, root, gv)
	if gv.SafeCount != 24 {
		t.Fatalf("GuardedCallback SafeCount = %d, want 24", gv.SafeCount)
	}
	// 4 scripts * 3 IFs = 12 per-IF slab resets while inside Token.ROOT, plus 1 Release() when Token.ROOT exits!
	if resets := j2gort.SlabResetCalls.Load(); resets != 12 {
		t.Fatalf("expected 12 per-guard slab Reset() calls during Token.ROOT traversal, got %d", resets)
	}
	if rels := j2gort.SlabReleaseCalls.Load(); rels != 1 {
		t.Fatalf("expected 1 slab Release() call on Token.ROOT exit even without caller ReleaseSlabs_(), got %d", rels)
	}

	// 1C-1: PeepholeMinimizeConditions + MinimizedCondition_MeasuredNode
	pkg_rhino.ResetArena()
	ifStmt := pkg_rhino.NewNode_1(pkg_rhino.TokIF)
	not1 := pkg_rhino.NewNode_1(pkg_rhino.TokNOT)
	not2 := pkg_rhino.NewNode_1(pkg_rhino.TokNOT)
	not2.AddChildToBack(pkg_rhino.AllocNamed(pkg_rhino.TokNAME, "condA"))
	not1.AddChildToBack(not2)
	ifStmt.AddChildToBack(not1)
	ifStmt.AddChildToBack(pkg_rhino.NewNode_1(pkg_rhino.TokBLOCK))
	j2gort.ResetWave31BMetrics()
	pmc := jscomp.NewPeepholeMinimizeConditions(false)
	pmc.OptimizeSubtree(ifStmt)
	if j2gort.SlabAllocCalls.Load() == 0 {
		t.Fatalf("expected PeepholeMinimizeConditions to allocate MeasuredNodes from PooledSlab")
	}
	pmc.EndTraversal()
	if j2gort.SlabReleaseCalls.Load() != 2 {
		t.Fatalf("expected EndTraversal() to release both MinimizedCondition slabs, got %d", j2gort.SlabReleaseCalls.Load())
	}

	// 1C-2: OptimizeCalls_NodeList
	j2gort.ResetWave31BMetrics()
	pfi := &jscomp.PureFunctionIdentifier{}
	oc := &jscomp.OptimizeCalls{}
	_ = oc
	fnNode := pkg_rhino.NewNode_1(pkg_rhino.TokFUNCTION)
	fnNode.AddChildToBack(pkg_rhino.AllocNamed(pkg_rhino.TokNAME, "dup"))
	fnNode.AddChildToBack(pkg_rhino.AllocNamed(pkg_rhino.TokNAME, "dup"))
	if !jscomp.OptimizeParameters_MayReferenceParamBeforeBody(fnNode) {
		t.Fatalf("expected duplicate parameter detection to return true")
	}
	if j2gort.SlabAllocCalls.Load() == 0 || j2gort.SlabReleaseCalls.Load() == 0 {
		t.Fatalf("expected OptimizeParameters_MayReferenceParamBeforeBody to use and release PooledSlab")
	}
	_ = pfi

	// 1C-3: MakeDeclaredNamesUnique ContextualRenamerBase vs BoilerplateRenamer
	j2gort.ResetWave31BMetrics()
	norm := &jscomp.Normalize{}
	norm.Process_2(0, root)
	if norm.LastRenamed != "x$jscomp$1" {
		t.Fatalf("expected Normalize ContextualRenamer rename 'x$jscomp$1', got %q", norm.LastRenamed)
	}
	if j2gort.SlabAllocCalls.Load() != 1 || j2gort.SlabResetCalls.Load() != 1 || j2gort.SlabReleaseCalls.Load() != 1 {
		t.Fatalf("expected 1 ContextualRenamer slab alloc, 1 Reset, 1 Release; got alloc=%d reset=%d rel=%d",
			j2gort.SlabAllocCalls.Load(), j2gort.SlabResetCalls.Load(), j2gort.SlabReleaseCalls.Load())
	}
	// Verify BoilerplateRenamer subclass still dispatches to InlineRenamer!
	bpMdnu := jscomp.MakeDeclaredNamesUnique_S_Builder().WithRenamer(jscomp.NewMakeDeclaredNamesUnique_BoilerplateRenamer("bp_")).Build()
	nt := &jscomp.NodeTraversal{}
	bpMdnu.EnterScope(nt)
	bpMdnu.EnterScope(nt)
	if got := bpMdnu.Resolve("v"); got != "v$inline$bp_" {
		t.Fatalf("BoilerplateRenamer child scope broken: got %q, want %q", got, "v$inline$bp_")
	}

	// 3: Native set lowering on 1-arg LinkedHashSet_1_Ext_java_util_Collection + external caller compatibility!
	j2gort.ResetWave31BMetrics()
	extSet := j2gort.NewExt_java_util_LinkedHashSet()
	vStmt := pkg_rhino.NewNode_1(pkg_rhino.TokVAR)
	vStmt.AddChildToBack(pkg_rhino.AllocNamed(pkg_rhino.TokNAME, "loc1"))
	j2gort.ResetWave31BMetrics()
	// External caller using original exported signature still works:
	jscomp.FunctionArgumentInjector_GatherLocalNames(vStmt, extSet)
	if !extSet.Contains("loc1") {
		t.Fatalf("external FunctionArgumentInjector_GatherLocalNames failed to add 'loc1'")
	}

	// MarkChildrenRead: ZERO LinkedHashSet allocations!
	pRoot := jscomp.NewDeadPropertyAssignmentElimination_Property("root")
	pChild := jscomp.NewDeadPropertyAssignmentElimination_Property("child")
	pGrand := jscomp.NewDeadPropertyAssignmentElimination_Property("grand")
	pRoot.AddChild(pChild)
	pChild.AddChild(pGrand)
	j2gort.ResetWave31BMetrics()
	pRoot.MarkChildrenRead()
	if j2gort.LinkedHashSetCalls.Load() != 0 {
		t.Fatalf("expected 0 LinkedHashSet allocations during MarkChildrenRead(), got %d", j2gort.LinkedHashSetCalls.Load())
	}
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "wave31b_rt_test.go"), []byte(rtSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	// Also add an internal package test for unexported fields in decisions_log.go & optimize_calls.go!
	intTestSrc := `package jscomp

import (
	"testing"

	pkg_diagnostic "example.com/w31btest/gen/diagnostic"
	"example.com/w31btest/gen/j2gort"
	pkg_rhino "example.com/w31btest/gen/rhino"
)

func TestWave31BInternalLoggingAndOptimizeCalls(t *testing.T) {
	j2gort.ResetWave31BMetrics()
	dd := &DirectDecisionPass{decisionsLog: &CompilerDecisionsLog{enabled: false}}
	dd.Run("fnA")
	mem := &pkg_diagnostic.MemLogFile{Enabled: false}
	cp := &CollapsePropsHelper{outer_: &CollapseOuter{decisionsLog: mem}}
	cp.CheckName("a.b", "escaped")
	gnNil := &GlobalNsName{outer_: &GlobalNsOuter{decisionsLog: nil}, full: "ns.x"}
	gnNil.CheckChildren("2")
	gnOff := &GlobalNsName{outer_: &GlobalNsOuter{decisionsLog: mem}, full: "ns.x"}
	gnOff.CheckChildren("2")
	if got := j2gort.SupplierLambdaCalls.Load(); got != 0 {
		t.Fatalf("expected 0 Supplier lambda allocations when IsLogging() is false/nil, got %d", got)
	}

	// Enable logging and verify messages are recorded faithfully!
	dd.decisionsLog.enabled = true
	mem.Enabled = true
	dd.Run("fnA")
	cp.CheckName("a.b", "escaped")
	gnOff.CheckChildren("2")
	if len(dd.decisionsLog.entries) != 1 || dd.decisionsLog.entries[0] != "direct:fnA" {
		t.Fatalf("unexpected CompilerDecisionsLog entries: %+v", dd.decisionsLog.entries)
	}
	if len(mem.Entries) != 2 || mem.Entries[0] != "a.b: skip escaped" || mem.Entries[1] != "ns.x:sets=2" {
		t.Fatalf("unexpected MemLogFile entries: %+v", mem.Entries)
	}

	// OptimizeCalls + PureFunctionIdentifier slab test:
	j2gort.ResetWave31BMetrics()
	pfi := &PureFunctionIdentifier{}
	oc := &OptimizeCalls{compiler: &Compiler{}, passes: []OptimizeCalls_CallGraphCompilerPass{pfi}}
	oc.Process_2(pkg_rhino.Node(1), pkg_rhino.Node(1))
	if pfi.MergedCount != 2 {
		t.Fatalf("expected PureFunctionIdentifier MergedCount == 2, got %d", pfi.MergedCount)
	}
	if j2gort.SlabAllocCalls.Load() != 3 || j2gort.SlabReleaseCalls.Load() != 1 {
		t.Fatalf("expected 3 OptimizeCalls_NodeList slab allocs and 1 Release, got alloc=%d rel=%d",
			j2gort.SlabAllocCalls.Load(), j2gort.SlabReleaseCalls.Load())
	}

	// FunctionArgumentInjector in-package findParametersReferencedAfterSideEffect: ZERO LinkedHashSet allocations!
	pkg_rhino.ResetArena()
	fnBody := pkg_rhino.NewNode_1(pkg_rhino.TokBLOCK)
	varNode := pkg_rhino.NewNode_1(pkg_rhino.TokVAR)
	varNode.AddChildToBack(pkg_rhino.AllocNamed(pkg_rhino.TokNAME, "innerVar"))
	fnBody.AddChildToBack(varNode)
	params := j2gort.Ext_com_google_common_collect_ImmutableSet_Of_1("p0")
	fai := &FunctionArgumentInjector{}
	j2gort.ResetWave31BMetrics()
	res := fai.priv_FunctionArgumentInjector_findParametersReferencedAfterSideEffect_2877(params, fnBody)
	if j2gort.LinkedHashSetCalls.Load() != 0 {
		t.Fatalf("expected 0 LinkedHashSet_1_Collection allocations, got %d", j2gort.LinkedHashSetCalls.Load())
	}
	if !res.Contains("p0") || !res.Contains("innerVar") {
		t.Fatalf("expected result ImmutableSet to contain both 'p0' and 'innerVar'")
	}
}
`
	if err := os.WriteFile(filepath.Join(jscompDir, "wave31b_internal_test.go"), []byte(intTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-timeout=30s", "-race", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race on lowered Wave 31B workspace failed: %v\n%s", err, string(out))
	}
}

func TestWave31B_RealGenJscompAllThreePassesIdempotencyAndBuild(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/rhino or gen/j2gort is mid-edit; skipping real gen/jscomp copy check")
	}
	checkCmd := exec.Command("go", "build", "./transpiled/gen/rhino", "./transpiled/gen/diagnostic")
	checkCmd.Dir = repoRoot
	if err := checkCmd.Run(); err != nil {
		t.Skipf("gen/rhino or gen/diagnostic is mid-edit (%v); skipping real gen/jscomp copy check", err)
	}

	realJscompDir := filepath.Join(repoRoot, "transpiled", "gen", "jscomp")
	entries, err := os.ReadDir(realJscompDir)
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()
	origModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/realw31b", 1) +
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
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") || strings.HasSuffix(ent.Name(), "_test.go") {
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

	w31bPasses := []passes.Pass{
		passes.NewNativeSetLoweringPass(),
		passes.NewSlabAllocLoweringPass(),
		passes.NewJ2gortPeepholePass(),
	}
	res1, err := passes.RunOnPackageDir(copyDir, w31bPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on copied real gen/jscomp failed: %v", err)
	}
	t.Logf("Wave 31B Real gen/jscomp Run 1 sites: total=%d files=%d byPass=%+v", res1.TotalSites, res1.FilesModified, res1.ByPass)

	checkContains := func(file, sub string) {
		t.Helper()
		b, rErr := os.ReadFile(filepath.Join(copyDir, file))
		if rErr != nil {
			t.Fatal(rErr)
		}
		if !strings.Contains(string(b), sub) {
			t.Fatalf("expected %s to contain %q", file, sub)
		}
	}
	checkNotContains := func(file, sub string) {
		t.Helper()
		b, rErr := os.ReadFile(filepath.Join(copyDir, file))
		if rErr != nil {
			t.Fatal(rErr)
		}
		if strings.Contains(string(b), sub) {
			t.Fatalf("expected %s NOT to contain %q", file, sub)
		}
	}

	// Task 1A: 3 exact phase boundaries (PerformTranspilationAndOptimizations defer, InferAllScopes entry only, SerializeRoots entry only)
	checkContains("compiler.go", "defer j2gort.PurgeAllRuntimePools()")
	checkContains("type_inference_pass.go", "j2gort.PurgeAllRuntimePools()")
	checkNotContains("type_inference_pass.go", "defer j2gort.PurgeAllRuntimePools()")
	checkContains("typed_ast_serializer.go", "j2gort.PurgeAllRuntimePools()")
	checkNotContains("process_defines.go", "PurgeAllRuntimePools")
	checkNotContains("variable_reference_check.go", "PurgeAllRuntimePools")

	// Task 1B: GuardedCallbackBase slabCtxCount_ reset + root Release
	checkContains("guarded_callback.go", "slabCtxCount_")
	checkContains("guarded_callback.go", ".slab_GuardedCallback_Context.Release()")
	checkContains("guarded_callback.go", ".slab_GuardedCallback_Context.Reset()")

	// Task 1C: 3 OPT pass-local slabs
	checkContains("peephole_minimize_conditions.go", "slab_MinimizedCondition_MeasuredNode")
	checkContains("peephole_minimize_conditions.go", "minimizedConditionFromConditionNodeWithSlab_")
	checkContains("peephole_minimize_conditions.go", "func (this_ *PeepholeMinimizeConditions) EndTraversal()")
	checkContains("optimize_calls.go", "slab_OptimizeCalls_NodeList")
	checkContains("optimize_calls.go", "allocPooled_OptimizeCalls_NodeList")
	if pfiBytes, err := os.ReadFile(filepath.Join(copyDir, "pure_function_identifier.go")); err == nil &&
		!strings.Contains(string(pfiBytes), "scratchRvalueLeaves") {
		checkContains("pure_function_identifier.go", "allocPooled_OptimizeCalls_NodeList(&referenceMap.slab_OptimizeCalls_NodeList)")
	}
	if optParamBytes, err := os.ReadFile(filepath.Join(copyDir, "optimize_parameters.go")); err == nil &&
		!strings.Contains(string(optParamBytes), "OptimizeParameters_CollectParamListNames(") {
		checkContains("optimize_parameters.go", "allocPooled_OptimizeCalls_NodeList(&slabNodeList_)")
	}
	checkContains("make_declared_names_unique.go", "slab_ContextualRenamer")
	checkContains("make_declared_names_unique.go", "createForChildScopeWithSlab_")
	checkContains("normalize.go", "defer renamer.ReleaseSlabs_()")

	// Task 2: Supplier-based decisionsLog wrapper calls guarded by IsLogging()
	checkContains("inline_and_collapse_properties.go", "this_.outer_.decisionsLog.IsLogging() {")
	checkContains("inline_and_collapse_properties.go", "this_.priv_CollapseProperties_logDecisionForName_")
	checkContains("global_namespace.go", "this_.outer_.decisionsLog.IsLogging() {\n\t\t\tthis_.priv_Name_logChildNamesDecision_")
	checkContains("strip_code.go", "if this_.decisionsLog.IsLogging() {\n\t\t\t\tthis_.priv_Strip_logStripName_")

	// Task 3: 1-arg NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection lowered in both OPT hotspots
	checkNotContains("function_argument_injector.go", "NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection")
	if faiBytes, err := os.ReadFile(filepath.Join(copyDir, "function_argument_injector.go")); err == nil &&
		!strings.Contains(string(faiBytes), "FunctionArgumentInjector_GatherLocalNames(root, nil)") {
		checkContains("function_argument_injector.go", "ordSet_FunctionArgumentInjector_GatherLocalNames(root, &locals)")
	}
	checkContains("function_argument_injector.go", "func FunctionArgumentInjector_GatherLocalNames(")
	checkNotContains("dead_property_assignment_elimination.go", "NewExt_java_util_LinkedHashSet_1_Ext_java_util_Collection")
	checkContains("dead_property_assignment_elimination.go", "j2gort.OrderedSetFromSet(this_.children)")
	checkContains("dead_property_assignment_elimination.go", "j2gort.DequeFromOrderedSet(propertiesSet)")

	// Build lowered real gen/jscomp!
	buildCmd := exec.Command("go", "build", "./jscomp")
	buildCmd.Dir = tmpDir
	if out, bErr := buildCmd.CombinedOutput(); bErr != nil {
		t.Fatalf("go build ./jscomp on lowered real gen/jscomp failed: %v\n%s", bErr, string(out))
	}

	// Run 2 on lowered real gen/jscomp MUST produce 0 sites and 0 modified files!
	res2, err := passes.RunOnPackageDir(copyDir, w31bPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on lowered real gen/jscomp failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("Run 2 on real gen/jscomp NOT IDEMPOTENT: got %+v", res2)
	}
}

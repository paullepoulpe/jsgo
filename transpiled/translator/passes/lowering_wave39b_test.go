package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

const wave39bMockJ2gortSrc = `package j2gort

import (
	"sync/atomic"
)

const NullString = "\xff\x00j2go_null"

var (
	PooledSlabAllocCount   atomic.Int64
	PooledSlabReleaseCount atomic.Int64
	OrderedMapReleaseCount atomic.Int64
)

type ClinitState struct {
	done uint32
}

func EnterClinit(s *ClinitState, done *uint32) bool {
	if atomic.LoadUint32(done) != 0 {
		return false
	}
	return true
}

func ExitClinit(s *ClinitState, done *uint32) {
	atomic.StoreUint32(done, 1)
}

func AssignIdentHash(p *atomic.Int32) int32 {
	p.Store(1)
	return 1
}

func NewStringArray(n int32) []string {
	return make([]string, n)
}

func CheckNotNullString(s string) string {
	if s == NullString {
		panic("NPE")
	}
	return s
}

func CheckNotNullPtr[T any](p *T) *T {
	if p == nil {
		panic("NPE")
	}
	return p
}

func CheckNotNullIface[I interface{ IsJavaNil_() bool }](v I) I {
	if any(v) == nil || v.IsJavaNil_() {
		panic("NPE")
	}
	return v
}

func CheckNotNullNode[N ~uint32](n N) N {
	if n == 0 {
		panic("NPE")
	}
	return n
}

func CheckState(cond bool) {
	if !cond {
		panic("IllegalStateException")
	}
}

func CheckArgument(cond bool) {
	if !cond {
		panic("IllegalArgumentException")
	}
}

type Slab[T any] struct {
	items []*T
}

func (s *Slab[T]) Alloc() *T {
	v := new(T)
	s.items = append(s.items, v)
	return v
}

func (s *Slab[T]) Release() {
	for i := range s.items {
		var zero T
		*s.items[i] = zero
	}
	s.items = s.items[:0]
}

type PooledSlab[T any] struct {
	items []*T
}

func (s *PooledSlab[T]) Alloc() *T {
	PooledSlabAllocCount.Add(1)
	v := new(T)
	s.items = append(s.items, v)
	return v
}

func (s *PooledSlab[T]) Reset() {
	for i := range s.items {
		var zero T
		*s.items[i] = zero
	}
	s.items = s.items[:0]
}

func (s *PooledSlab[T]) Release() {
	if s == nil {
		return
	}
	PooledSlabReleaseCount.Add(1)
	for i := range s.items {
		var zero T
		*s.items[i] = zero
	}
	s.items = nil
}

func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	if s == nil {
		return
	}
	for i := range s.items {
		if fn != nil {
			fn(s.items[i])
		}
	}
	s.Release()
}

type OrderedMap[K comparable, V any] struct {
	keys []K
	vals []V
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}

func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}

func ReleaseOrderedMap[K comparable, V any](m *OrderedMap[K, V]) {
	if m == nil {
		return
	}
	OrderedMapReleaseCount.Add(1)
	m.Clear()
}

func (m *OrderedMap[K, V]) Clear() {
	clear(m.keys)
	clear(m.vals)
	m.keys = m.keys[:0]
	m.vals = m.vals[:0]
}

func (m *OrderedMap[K, V]) ReleaseToPool() {
	OrderedMapReleaseCount.Add(1)
	m.Clear()
}

func (m *OrderedMap[K, V]) Get(k K) V {
	for i, ek := range m.keys {
		if ek == k {
			return m.vals[i]
		}
	}
	var zero V
	return zero
}

func (m *OrderedMap[K, V]) Put(k K, v V) V {
	for i, ek := range m.keys {
		if ek == k {
			old := m.vals[i]
			m.vals[i] = v
			return old
		}
	}
	m.keys = append(m.keys, k)
	m.vals = append(m.vals, v)
	var zero V
	return zero
}

type CustomOrderedMap[K comparable, V any] struct {
	OrderedMap[K, V]
}

func NewCustomOrderedMap[K comparable, V any]() CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{}
}

type OrderedSet[E comparable] struct {
	items []E
}

func NewOrderedSet[E comparable]() OrderedSet[E] {
	return OrderedSet[E]{}
}

func (s *OrderedSet[E]) Add(v E) bool {
	for _, x := range s.items {
		if x == v {
			return false
		}
	}
	s.items = append(s.items, v)
	return true
}

type SetCursor[E any] struct {
	items []E
	idx   int
}

func (s *OrderedSet[E]) Iter() SetCursor[E] {
	return SetCursor[E]{items: s.items}
}

func (c *SetCursor[E]) HasNext() bool { return c.idx < len(c.items) }
func (c *SetCursor[E]) Next() E {
	v := c.items[c.idx]
	c.idx++
	return v
}

type Deque[E any] struct {
	items []E
}

func NewDequeCap[E any](cap int32) Deque[E] {
	return Deque[E]{items: make([]E, 0, cap)}
}

func (d *Deque[E]) Size() int32   { return int32(len(d.items)) }
func (d *Deque[E]) Push(v E)      { d.items = append([]E{v}, d.items...) }
func (d *Deque[E]) AddLast(v E)   { d.items = append(d.items, v) }
func (d *Deque[E]) Peek() E {
	if len(d.items) == 0 {
		var zero E
		return zero
	}
	return d.items[0]
}
func (d *Deque[E]) Pop() E {
	v := d.items[0]
	d.items = d.items[1:]
	return v
}
func (d *Deque[E]) Iter() SetCursor[E] {
	return SetCursor[E]{items: d.items}
}

type Ext_java_util_BitSet struct {
	bits uint64
}

func NewExt_java_util_BitSet_1(nbits int32) *Ext_java_util_BitSet {
	return &Ext_java_util_BitSet{}
}

func InitExt_java_util_BitSet_1(b *Ext_java_util_BitSet, nbits int32) {
	b.bits = 0
}

func (b *Ext_java_util_BitSet) Set_1(idx int32) {
	b.bits |= 1 << uint(idx)
}

func (b *Ext_java_util_BitSet) Or(other *Ext_java_util_BitSet) {
	b.bits |= other.bits
}

func (b *Ext_java_util_BitSet) IsEmpty() bool {
	return b.bits == 0
}
`

func TestWave39B_ScopeClinitSlabResilienceAndParity(t *testing.T) {
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

	mustWrite("go.mod", "module example.com/w39b_scope\n\ngo 1.23\n")
	mustWrite("gen/j2gort/j2gort.go", wave39bMockJ2gortSrc)
	mustWrite("gen/rhino/node.go", `package rhino

type Node uint32
`)

	jscompScopeSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w39b_scope/gen/j2gort"
	pkg_rhino "example.com/w39b_scope/gen/rhino"
)

type AbstractScope interface {
	IsJavaNil_() bool
	Implements_jscomp_AbstractScope()
}

func ToIface_AbstractScope(s *Scope) AbstractScope { return s }

type ScopeCreator interface {
	IsJavaNil_() bool
	Implements_jscomp_ScopeCreator()
}

type Scope struct {
	rootNode pkg_rhino.Node
	parent   *Scope
	vars     *j2gort.OrderedMap[string, *Var]
	depth    int32
}

func (this_ *Scope) BindSelf_(s any)               { _ = s }
func (this_ *Scope) IsJavaNil_() bool              { return this_ == nil }
func (*Scope) Implements_jscomp_AbstractScope()    {}

var static_Scope_EMPTY_VAR_NAMES []string
var static_Scope_SINGLE_ASCII_VAR_NAMES [][]string

func StaticGet_Scope_EMPTY_VAR_NAMES() []string {
	clinit_Scope()
	return static_Scope_EMPTY_VAR_NAMES
}

var clinitDone_Scope uint32
var clinitState_Scope j2gort.ClinitState
var ClinitScopeRunCount atomic.Int32

func clinit_Scope() {
	if atomic.LoadUint32(&clinitDone_Scope) != 0 {
		return
	}
	clinitSlow_Scope()
}

func clinitSlow_Scope() {
	if !j2gort.EnterClinit(&clinitState_Scope, &clinitDone_Scope) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Scope, &clinitDone_Scope)
	ClinitScopeRunCount.Add(1)
	static_Scope_EMPTY_VAR_NAMES = j2gort.NewStringArray(int32(0))
	static_Scope_SINGLE_ASCII_VAR_NAMES = make([][]string, int32(128))
	static_Scope_SINGLE_ASCII_VAR_NAMES[int32('a')] = []string{"a"}
}

func (this_ *Scope) initFields_Scope() { _ = this_ }

func NewScope_1(rootNode pkg_rhino.Node) *Scope {
	clinit_Scope()
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_1(this_, rootNode)
	return this_
}

func InitScope_1(this_ *Scope, rootNode pkg_rhino.Node) {
	_ = rootNode
	this_.initFields_Scope()
	this_.rootNode = rootNode
	this_.parent = nil
	this_.depth = int32(0)
}

func NewScope_2(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	clinit_Scope()
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_2(this_, parent, rootNode)
	return this_
}

func InitScope_2(this_ *Scope, parent *Scope, rootNode pkg_rhino.Node) {
	_ = parent
	_ = rootNode
	this_.initFields_Scope()
	this_.rootNode = rootNode
	this_.parent = parent
	this_.depth = parent.depth + int32(1)
}

func Scope_CreateGlobalScope(rootNode pkg_rhino.Node) *Scope {
	clinit_Scope()
	_ = rootNode
	return NewScope_1(rootNode)
}

func Scope_CreateChildScope(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	clinit_Scope()
	_ = parent
	_ = rootNode
	return NewScope_2(parent, rootNode)
}

func (this_ *Scope) Declare(name string, nameNode pkg_rhino.Node) *Var {
	if this_.vars == nil {
		this_.vars = j2gort.NewOrderedMapPtr[string, *Var]()
	}
	v := NewVar(name, nameNode, this_, int32(0))
	this_.vars.Put(name, v)
	return v
}

type Var struct {
	name     string
	nameNode pkg_rhino.Node
	scope    *Scope
	index    int32
}

func (this_ *Var) initFields_Var() { _ = this_ }

func NewVar(name string, nameNode pkg_rhino.Node, scope *Scope, index int32) *Var {
	this_ := &Var{}
	InitVar(this_, name, nameNode, scope, index)
	return this_
}

func InitVar(this_ *Var, name string, nameNode pkg_rhino.Node, scope *Scope, index int32) {
	this_.initFields_Var()
	this_.name = name
	this_.nameNode = nameNode
	this_.scope = scope
	this_.index = index
}

// Negative control: genuine singleton instance stored in static_* variable of type *SingletonGuardNode
// MUST still be disqualified from slabbing!
type SingletonGuardNode struct {
	tag int32
}

var static_SingletonGuardNode_INSTANCE *SingletonGuardNode

func init() {
	static_SingletonGuardNode_INSTANCE = NewSingletonGuardNode(int32(99))
}

func NewSingletonGuardNode(tag int32) *SingletonGuardNode {
	this_ := &SingletonGuardNode{}
	InitSingletonGuardNode(this_, tag)
	return this_
}

func InitSingletonGuardNode(this_ *SingletonGuardNode, tag int32) {
	this_.tag = tag
}
`

	jscompSscSrc := `package jscomp

import (
	pkg_rhino "example.com/w39b_scope/gen/rhino"
)

type SyntacticScopeCreator struct {
	tag int32
}

func (this_ *SyntacticScopeCreator) IsJavaNil_() bool          { return this_ == nil }
func (*SyntacticScopeCreator) Implements_jscomp_ScopeCreator() {}

func NewSyntacticScopeCreator_1() *SyntacticScopeCreator {
	return &SyntacticScopeCreator{tag: 1}
}

func (this_ *SyntacticScopeCreator) CreateScope_2_Node_Scope(n pkg_rhino.Node, parent *Scope) *Scope {
	_ = n
	_ = parent
	var scope *Scope
	if parent == nil {
		scope = Scope_CreateGlobalScope(n)
	} else {
		scope = Scope_CreateChildScope(parent, n)
	}
	scope.Declare("x", n)
	_ = NewSingletonGuardNode(this_.tag)
	return scope
}
`

	mustWrite("orig/jscomp/scope.go", strings.ReplaceAll(jscompScopeSrc, "example.com/w39b_scope/gen/", "example.com/w39b_scope/gen/"))
	mustWrite("orig/jscomp/syntactic_scope_creator.go", jscompSscSrc)
	mustWrite("opt/jscomp/scope.go", jscompScopeSrc)
	mustWrite("opt/jscomp/syntactic_scope_creator.go", jscompSscSrc)

	optDir := filepath.Join(tmpDir, "opt/jscomp")
	slabPass := []passes.Pass{passes.NewSlabAllocLoweringPass()}

	s1, err := passes.RunOnPackageDir(optDir, slabPass, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero slab sites on Run 1, got %+v", s1)
	}

	// Run 2 MUST produce strictly 0 sites!
	s2, err := passes.RunOnPackageDir(optDir, slabPass, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d ByPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	gotSscBytes, err := os.ReadFile(filepath.Join(optDir, "syntactic_scope_creator.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotSsc := string(gotSscBytes)
	gotScopeBytes, err := os.ReadFile(filepath.Join(optDir, "scope.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotScope := string(gotScopeBytes)

	for _, want := range []string{
		"slab_Scope j2gort.PooledSlab[Scope]",
		"slab_Var   j2gort.PooledSlab[Var]",
		"allocPooledVia_Scope_CreateGlobalScope(&this_.slab_Scope, n)",
		"allocPooledVia_Scope_CreateChildScope(&this_.slab_Scope, parent, n)",
		"this_.slab_Scope.ReleaseWith(",
	} {
		if !strings.Contains(gotSsc, want) {
			t.Errorf("syntactic_scope_creator.go missing %q in:\n%s", want, gotSsc)
		}
	}
	if strings.Contains(gotSsc, "slab_SingletonGuardNode") {
		t.Errorf("genuine singleton SingletonGuardNode must NOT be slabbed:\n%s", gotSsc)
	}
	for _, want := range []string{
		"func allocPooled_Scope_1(slab_ *j2gort.PooledSlab[Scope], rootNode pkg_rhino.Node) *Scope {\n\tclinit_Scope()",
		"func allocPooled_Scope_2(slab_ *j2gort.PooledSlab[Scope], parent *Scope, rootNode pkg_rhino.Node) *Scope {\n\tclinit_Scope()",
		"j2gort.AcquireOrderedMapPtr[string, *Var]()",
	} {
		if !strings.Contains(gotScope, want) {
			t.Errorf("scope.go missing %q in:\n%s", want, gotScope)
		}
	}

	diffTest := `package scope_diff_test

import (
	"testing"

	"example.com/w39b_scope/gen/j2gort"
	pkg_rhino "example.com/w39b_scope/gen/rhino"
	opt "example.com/w39b_scope/opt/jscomp"
	orig "example.com/w39b_scope/orig/jscomp"
)

func TestScopeClinitAndPooledSlabParity(t *testing.T) {
	j2gort.PooledSlabAllocCount.Store(0)
	j2gort.PooledSlabReleaseCount.Store(0)
	j2gort.OrderedMapReleaseCount.Store(0)

	osc := orig.NewSyntacticScopeCreator_1()
	og := osc.CreateScope_2_Node_Scope(pkg_rhino.Node(1), nil)
	oc := osc.CreateScope_2_Node_Scope(pkg_rhino.Node(2), og)

	xsc := opt.NewSyntacticScopeCreator_1()
	xg := xsc.CreateScope_2_Node_Scope(pkg_rhino.Node(1), nil)
	xc := xsc.CreateScope_2_Node_Scope(pkg_rhino.Node(2), xg)

	if opt.ClinitScopeRunCount.Load() != 1 {
		t.Fatalf("expected clinit_Scope to execute once, got %d", opt.ClinitScopeRunCount.Load())
	}
	if len(opt.StaticGet_Scope_EMPTY_VAR_NAMES()) != len(orig.StaticGet_Scope_EMPTY_VAR_NAMES()) {
		t.Fatalf("EMPTY_VAR_NAMES mismatch")
	}
	_ = oc
	_ = xc
	if j2gort.PooledSlabAllocCount.Load() != 4 {
		t.Fatalf("expected 4 PooledSlab allocations (2 Scopes + 2 Vars), got %d", j2gort.PooledSlabAllocCount.Load())
	}
	xsc.ReleaseSlabs_()
	if j2gort.PooledSlabReleaseCount.Load() != 2 {
		t.Fatalf("expected 2 PooledSlab releases (slab_Scope + slab_Var), got %d", j2gort.PooledSlabReleaseCount.Load())
	}
	if j2gort.OrderedMapReleaseCount.Load() != 2 {
		t.Fatalf("expected 2 Scope.vars OrderedMap releases, got %d", j2gort.OrderedMapReleaseCount.Load())
	}
}
`
	mustWrite("scope_diff_test.go", diffTest)
	cmd := exec.Command("go", "test", "-race", "-count=1", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("scope differential runtime test failed: %v\n%s", err, string(out))
	}
}

func TestWave39B_HasNonNegativeLengthPeepholeAndNPEParity(t *testing.T) {
	t.Run("WithHasNonNegativeLength", func(t *testing.T) {
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

		mustWrite("go.mod", "module example.com/w39b_peep\n\ngo 1.23\n")
		mustWrite("gen/j2gort/j2gort.go", wave39bMockJ2gortSrc)
		mustWrite("shims/nodestore/node.go", `package nodestore

type Node uint32

var lengths = map[Node]int32{
	1: 0,
	2: 42,
	3: 70000,
	4: -1,
}

func (n Node) Length() int32 {
	if n == 0 {
		panic("nodestore Node(0)")
	}
	return lengths[n]
}

func (n Node) HasNonNegativeLength() bool {
	if n == 0 {
		panic("NPE")
	}
	return lengths[n] >= 0
}
`)
		mustWrite("gen/rhino/node.go", `package rhino

import "example.com/w39b_peep/shims/nodestore"

type Node uint32

func (this_ Node) nodeGet_length() int32 {
	if this_ == 0 {
		panic("NPE")
	}
	return nodestore.Node(this_).Length()
}

func (this_ Node) GetLength() int32 {
	return this_.nodeGet_length()
}

func (this_ Node) GetLineno() int32 {
	if this_ == 0 {
		panic("NPE")
	}
	return 1
}
`)

		appSrc := `package jscomp

import (
	"example.com/w39b_peep/gen/j2gort"
	pkg_rhino "example.com/w39b_peep/gen/rhino"
)

type TypedAstSerializer struct {
	curNode pkg_rhino.Node
}

func (this_ *TypedAstSerializer) CheckPos(n pkg_rhino.Node) int32 {
	_ = n
	if !(n.GetLength() >= int32(0)) {
		j2gort.CheckState(false)
	}
	return n.GetLineno()
}

func (this_ *TypedAstSerializer) CheckSourceInfo(n pkg_rhino.Node) bool {
	if n.GetLength() < int32(0) {
		return false
	}
	return 0 <= this_.curNode.GetLength() && !(0 > n.GetLength())
}

type CodeBuilder struct{ n int32 }

func (cb *CodeBuilder) GetLength() int32 { return cb.n }

var sideEffectCalls int

func nextNodeWithSideEffect(n pkg_rhino.Node) pkg_rhino.Node {
	sideEffectCalls++
	return n
}

func NegativeControls(cb *CodeBuilder, n pkg_rhino.Node) bool {
	a := cb.GetLength() >= int32(0)
	b := nextNodeWithSideEffect(n).GetLength() >= int32(0)
	c := n.GetLength() > int32(0)
	return a && b && c
}
`
		mustWrite("orig/jscomp/serializer.go", appSrc)
		mustWrite("opt/jscomp/serializer.go", appSrc)

		optDir := filepath.Join(tmpDir, "opt/jscomp")
		allPasses := passes.AllPasses()

		s1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 1 failed: %v", err)
		}
		if s1.TotalSites == 0 {
			t.Fatalf("expected non-zero peephole sites on Run 1, got %+v", s1)
		}

		s2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 2 failed: %v", err)
		}
		if s2.TotalSites != 0 || s2.FilesModified != 0 {
			t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d ByPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
		}

		gotBytes, err := os.ReadFile(filepath.Join(optDir, "serializer.go"))
		if err != nil {
			t.Fatal(err)
		}
		got := string(gotBytes)

		for _, want := range []string{
			`"example.com/w39b_peep/shims/nodestore"`,
			"if !nodestore.Node(n).HasNonNegativeLength()",
			"nodestore.Node(this_.curNode).HasNonNegativeLength()",
			"cb.GetLength() >= int32(0)",
			"nextNodeWithSideEffect(n).GetLength() >= int32(0)",
			"n.GetLength() > int32(0)",
		} {
			if !strings.Contains(got, want) {
				t.Errorf("serializer.go missing %q in:\n%s", want, got)
			}
		}

		diffTest := `package peep_diff_test

import (
	"testing"

	pkg_rhino "example.com/w39b_peep/gen/rhino"
	opt "example.com/w39b_peep/opt/jscomp"
	orig "example.com/w39b_peep/orig/jscomp"
)

func catchPanicKind(fn func()) (panicked bool, isState bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			if r == "IllegalStateException" {
				isState = true
			}
		}
	}()
	fn()
	return false, false
}

func TestHasNonNegativeLengthDifferentialParity(t *testing.T) {
	os := &orig.TypedAstSerializer{}
	xs := &opt.TypedAstSerializer{}

	for _, nd := range []pkg_rhino.Node{0, 1, 2, 3, 4} {
		op, oState := catchPanicKind(func() { os.CheckPos(nd) })
		xp, xState := catchPanicKind(func() { xs.CheckPos(nd) })
		if op != xp || oState != xState {
			t.Fatalf("CheckPos(%d) panic mismatch: orig=(%v,%v) opt=(%v,%v)", nd, op, oState, xp, xState)
		}
	}

	os2 := &orig.TypedAstSerializer{}
	xs2 := &opt.TypedAstSerializer{}
	// Zero curNode on CheckSourceInfo(1) must panic with NPE on both!
	op0, _ := catchPanicKind(func() { os2.CheckSourceInfo(1) })
	xp0, _ := catchPanicKind(func() { xs2.CheckSourceInfo(1) })
	if !op0 || !xp0 {
		t.Fatalf("expected NPE on zero curNode: orig=%v opt=%v", op0, xp0)
	}
	// CheckSourceInfo(0) must panic with NPE on both!
	opN0, _ := catchPanicKind(func() { os2.CheckSourceInfo(0) })
	xpN0, _ := catchPanicKind(func() { xs2.CheckSourceInfo(0) })
	if !opN0 || !xpN0 {
		t.Fatalf("expected NPE on n==0 in CheckSourceInfo: orig=%v opt=%v", opN0, xpN0)
	}
	// CheckSourceInfo(4) returns false before touching curNode!
	if os2.CheckSourceInfo(4) != xs2.CheckSourceInfo(4) {
		t.Fatalf("CheckSourceInfo(4) mismatch")
	}
}
`
		mustWrite("peep_diff_test.go", diffTest)
		cmd := exec.Command("go", "test", "-race", "-count=1", "./...")
		cmd.Dir = tmpDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("peephole differential runtime test failed: %v\n%s", err, string(out))
		}
	})

	t.Run("FallbackWithoutHasNonNegativeLength", func(t *testing.T) {
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

		mustWrite("go.mod", "module example.com/w39b_fallback\n\ngo 1.23\n")
		mustWrite("gen/j2gort/j2gort.go", wave39bMockJ2gortSrc)
		mustWrite("shims/nodestore/node.go", `package nodestore

type Node uint32

func (n Node) Length() int32 { return int32(n) }
`)
		mustWrite("gen/rhino/node.go", `package rhino

import "example.com/w39b_fallback/shims/nodestore"

type Node uint32

func (this_ Node) GetLength() int32 { return nodestore.Node(this_).Length() }
`)
		mustWrite("gen/jscomp/s.go", `package jscomp

import pkg_rhino "example.com/w39b_fallback/gen/rhino"

func Check(n pkg_rhino.Node) bool {
	return n.GetLength() >= int32(0)
}
`)

		res, err := passes.RunOnPackageDir(filepath.Join(tmpDir, "gen/jscomp"), []passes.Pass{passes.NewJ2gortPeepholePass()}, passes.Options{})
		if err != nil {
			t.Fatalf("unexpected error when HasNonNegativeLength absent: %v", err)
		}
		if res.TotalSites != 0 {
			t.Fatalf("expected 0 sites when HasNonNegativeLength absent, got %+v", res)
		}
	})
}

func TestWave39B_OptPassLocalSlabsAndScratchReuseParity(t *testing.T) {
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

	mustWrite("go.mod", "module example.com/w39b_opt\n\ngo 1.23\n")
	mustWrite("gen/j2gort/j2gort.go", wave39bMockJ2gortSrc)
	mustWrite("gen/rhino/node.go", `package rhino

type Node uint32

type nodeInfo struct {
	name     string
	isFn     bool
	isAssign bool
	isExpr   bool
	isVar    bool
	parent   Node
	first    Node
	next     Node
}

var tree = map[Node]nodeInfo{}

func ResetTree() { clear(tree) }

func DefNode(id Node, name string, isFn, isVar, isAssign, isExpr bool, parent, first, next Node) {
	tree[id] = nodeInfo{
		name: name, isFn: isFn, isVar: isVar, isAssign: isAssign, isExpr: isExpr,
		parent: parent, first: first, next: next,
	}
}

func (n Node) IsFunction() bool { return tree[n].isFn }
func (n Node) IsName() bool     { return tree[n].name != "" && !tree[n].isVar && !tree[n].isAssign && !tree[n].isExpr && !tree[n].isFn }
func (n Node) IsAssign() bool   { return tree[n].isAssign }
func (n Node) IsExprResult() bool { return tree[n].isExpr }
func (n Node) IsVar() bool      { return tree[n].isVar }
func (n Node) GetString() string { return tree[n].name }
func (n Node) GetParent() Node  { return tree[n].parent }
func (n Node) GetFirstChild() Node { return tree[n].first }
func (n Node) GetNext() Node    { return tree[n].next }
func (n Node) HasChildren() bool { return tree[n].first != 0 }
`)

	globalNsSrc := `package jscomp

import (
	"example.com/w39b_opt/gen/j2gort"
	pkg_rhino "example.com/w39b_opt/gen/rhino"
)

type GlobalNamespace_NameProp struct {
	kind int32
}

var (
	namePropGetSet = &GlobalNamespace_NameProp{kind: 1}
	namePropOther  = &GlobalNamespace_NameProp{kind: 2}
	namePropFn     = &GlobalNamespace_NameProp{kind: 3}
)

func StaticGet_GlobalNamespace_NameProp_GET_SET() *GlobalNamespace_NameProp      { return namePropGetSet }
func StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT() *GlobalNamespace_NameProp { return namePropOther }
func StaticGet_GlobalNamespace_NameProp_FUNCTION() *GlobalNamespace_NameProp     { return namePropFn }

type GlobalNamespace struct {
	SeenNames []string
	SeenKinds []int32
}

type GlobalNamespace_BuildGlobalNamespace struct {
	outer_ *GlobalNamespace
}

func NewGlobalNamespace_BuildGlobalNamespace(outer_ *GlobalNamespace) *GlobalNamespace_BuildGlobalNamespace {
	this_ := &GlobalNamespace_BuildGlobalNamespace{}
	this_.outer_ = outer_
	return this_
}

func (this_ *GlobalNamespace_BuildGlobalNamespace) CollectKey(n pkg_rhino.Node, mode int32) {
	var name string = j2gort.NullString
	var type_ *GlobalNamespace_NameProp
	var analysis *GlobalNamespace_ObjLitStringKeyAnalysis = this_.priv_BuildGlobalNamespace_createObjLitStringKeyAnalysis_3064(n, mode)
	_ = analysis
	name = analysis.NameString()
	type_ = analysis.NameType()
	if name != j2gort.NullString && type_ != nil {
		this_.outer_.SeenNames = append(this_.outer_.SeenNames, name)
		this_.outer_.SeenKinds = append(this_.outer_.SeenKinds, type_.kind)
	}
}

func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_createObjLitStringKeyAnalysis_3064(stringKeyNode pkg_rhino.Node, mode int32) *GlobalNamespace_ObjLitStringKeyAnalysis {
	switch mode {
	case 1:
		return GlobalNamespace_ObjLitStringKeyAnalysis_ForObjLitAssignment(stringKeyNode.GetString(), StaticGet_GlobalNamespace_NameProp_FUNCTION())
	case 2:
		return GlobalNamespace_ObjLitStringKeyAnalysis_ForObjectDefineProperty(stringKeyNode.GetString())
	default:
		return GlobalNamespace_ObjLitStringKeyAnalysis_ForNonReference()
	}
}

type GlobalNamespace_ObjLitStringKeyAnalysis struct {
	nameString string
	nameType   *GlobalNamespace_NameProp
}

func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) initFields_GlobalNamespace_ObjLitStringKeyAnalysis() {
	_ = this_
}

func NewGlobalNamespace_ObjLitStringKeyAnalysis(nameString string, nameType *GlobalNamespace_NameProp) *GlobalNamespace_ObjLitStringKeyAnalysis {
	this_ := &GlobalNamespace_ObjLitStringKeyAnalysis{}
	InitGlobalNamespace_ObjLitStringKeyAnalysis(this_, nameString, nameType)
	return this_
}

func InitGlobalNamespace_ObjLitStringKeyAnalysis(this_ *GlobalNamespace_ObjLitStringKeyAnalysis, nameString string, nameType *GlobalNamespace_NameProp) {
	this_.initFields_GlobalNamespace_ObjLitStringKeyAnalysis()
	this_.nameString = nameString
	this_.nameType = nameType
}

func GlobalNamespace_ObjLitStringKeyAnalysis_ForObjectDefineProperty(nameString string) *GlobalNamespace_ObjLitStringKeyAnalysis {
	return NewGlobalNamespace_ObjLitStringKeyAnalysis(j2gort.CheckNotNullString(nameString), StaticGet_GlobalNamespace_NameProp_GET_SET())
}

func GlobalNamespace_ObjLitStringKeyAnalysis_ForObjLitAssignment(nameString string, nameType *GlobalNamespace_NameProp) *GlobalNamespace_ObjLitStringKeyAnalysis {
	return NewGlobalNamespace_ObjLitStringKeyAnalysis(j2gort.CheckNotNullString(nameString), nameType)
}

func GlobalNamespace_ObjLitStringKeyAnalysis_ForNonReference() *GlobalNamespace_ObjLitStringKeyAnalysis {
	return NewGlobalNamespace_ObjLitStringKeyAnalysis(j2gort.NullString, StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT())
}

func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) NameString() string { return this_.nameString }
func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) NameType() *GlobalNamespace_NameProp {
	return this_.nameType
}
`

	nodeIteratorsSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w39b_opt/gen/j2gort"
	pkg_rhino "example.com/w39b_opt/gen/rhino"
)

type AbstractCompiler interface {
	IsJavaNil_() bool
}

type MockCompiler struct{}

func (c *MockCompiler) IsJavaNil_() bool { return c == nil }

type NodeIterators_FunctionlessLocalScope struct {
	ancestors     []pkg_rhino.Node
	ancestorCount int32
	idHash_       atomic.Int32
}

func (this_ *NodeIterators_FunctionlessLocalScope) initFields_NodeIterators_FunctionlessLocalScope() {
	_ = this_
}

func NewNodeIterators_FunctionlessLocalScope_3(a0 pkg_rhino.Node, a1 pkg_rhino.Node, a2 pkg_rhino.Node) *NodeIterators_FunctionlessLocalScope {
	this_ := &NodeIterators_FunctionlessLocalScope{}
	InitNodeIterators_FunctionlessLocalScope_3(this_, a0, a1, a2)
	return this_
}

func InitNodeIterators_FunctionlessLocalScope_3(this_ *NodeIterators_FunctionlessLocalScope, a0 pkg_rhino.Node, a1 pkg_rhino.Node, a2 pkg_rhino.Node) {
	this_.initFields_NodeIterators_FunctionlessLocalScope()
	this_.ancestors = make([]pkg_rhino.Node, int32(16))
	var valid int32 = int32(0)
	if !(a0.IsFunction()) {
		valid = int32(1)
		if !(a1.IsFunction()) {
			valid = int32(2)
			if !(a2.IsFunction()) {
				valid = int32(3)
			}
		}
	}
	this_.ancestorCount = valid
	if valid == int32(3) {
		this_.ancestors[0] = a2
		this_.ancestors[1] = a1
		this_.ancestors[2] = a0
	} else if valid == int32(2) {
		this_.ancestors[0] = a1
		this_.ancestors[1] = a0
	} else if valid == int32(1) {
		this_.ancestors[0] = a0
	}
}

func (this_ *NodeIterators_FunctionlessLocalScope) pushAncestor(n pkg_rhino.Node) {
	c := this_.ancestorCount
	this_.ancestors[c] = n
	this_.ancestorCount = c + 1
}

func (this_ *NodeIterators_FunctionlessLocalScope) HasNext() bool {
	c := this_.ancestorCount
	return !((c == int32(1)) && (this_.ancestors[int32(0)].GetNext() == 0))
}

func (this_ *NodeIterators_FunctionlessLocalScope) NextNode() pkg_rhino.Node {
	for {
		this_.ancestorCount--
		topIdx := this_.ancestorCount
		current := this_.ancestors[topIdx]
		this_.ancestors[topIdx] = pkg_rhino.Node(0)
		nextSibling := current.GetNext()
		if nextSibling == 0 {
			current = this_.ancestors[topIdx-1]
			if current.IsFunction() {
				continue
			}
			return current
		}
		current = nextSibling
		this_.pushAncestor(current)
		if current.IsFunction() {
			continue
		}
		for current.HasChildren() {
			current = current.GetFirstChild()
			this_.pushAncestor(current)
			if current.IsFunction() {
				break
			}
		}
		if current.IsFunction() {
			continue
		}
		return current
	}
}

func (this_ *NodeIterators_FunctionlessLocalScope) Current() pkg_rhino.Node {
	return this_.ancestors[this_.ancestorCount-1]
}

type NodeIterators_LocalVarMotion struct {
	compiler            AbstractCompiler
	iterator            *NodeIterators_FunctionlessLocalScope
	varName             string
	lookAhead           pkg_rhino.Node
	idHash_             atomic.Int32
	valueHasSideEffects bool
}

func (this_ *NodeIterators_LocalVarMotion) initFields_NodeIterators_LocalVarMotion() {
	this_.varName = j2gort.NullString
}

func NewNodeIterators_LocalVarMotion(compiler AbstractCompiler, nameNode pkg_rhino.Node, iterator *NodeIterators_FunctionlessLocalScope) *NodeIterators_LocalVarMotion {
	this_ := &NodeIterators_LocalVarMotion{}
	InitNodeIterators_LocalVarMotion(this_, compiler, nameNode, iterator)
	return this_
}

func InitNodeIterators_LocalVarMotion(this_ *NodeIterators_LocalVarMotion, compiler AbstractCompiler, nameNode pkg_rhino.Node, iterator *NodeIterators_FunctionlessLocalScope) {
	this_.initFields_NodeIterators_LocalVarMotion()
	j2gort.CheckArgument(nameNode.IsName())
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.varName = nameNode.GetString()
	this_.valueHasSideEffects = false
	this_.iterator = iterator
	this_.advance(true)
}

func NodeIterators_LocalVarMotion_ForVar(compiler AbstractCompiler, name pkg_rhino.Node, var_ pkg_rhino.Node, block pkg_rhino.Node) *NodeIterators_LocalVarMotion {
	j2gort.CheckArgument(var_.IsVar())
	return NewNodeIterators_LocalVarMotion(compiler, name, NewNodeIterators_FunctionlessLocalScope_3(name, var_, block))
}

func NodeIterators_LocalVarMotion_ForAssign(compiler AbstractCompiler, name pkg_rhino.Node, assign pkg_rhino.Node, expr pkg_rhino.Node, block pkg_rhino.Node) *NodeIterators_LocalVarMotion {
	j2gort.CheckArgument(assign.IsAssign())
	j2gort.CheckArgument(expr.IsExprResult())
	return NewNodeIterators_LocalVarMotion(compiler, name, NewNodeIterators_FunctionlessLocalScope_3(assign, expr, block))
}

func (this_ *NodeIterators_LocalVarMotion) HasNext() bool {
	return this_.lookAhead != 0
}

func (this_ *NodeIterators_LocalVarMotion) NextNode() pkg_rhino.Node {
	next := this_.lookAhead
	this_.advance(false)
	return next
}

func (this_ *NodeIterators_LocalVarMotion) advance(atStart bool) {
	if !atStart {
		if this_.lookAhead == 0 {
			return
		}
		cur := this_.iterator.Current()
		if cur.IsName() && cur.GetString() == this_.varName {
			this_.lookAhead = 0
			return
		}
	}
	if !this_.iterator.HasNext() {
		this_.lookAhead = 0
		return
	}
	this_.lookAhead = this_.iterator.NextNode()
}

type InlineVariables struct {
	compiler AbstractCompiler
}

func NewInlineVariables(compiler AbstractCompiler) *InlineVariables {
	return &InlineVariables{compiler: compiler}
}

type InlineVariables_InliningBehavior struct {
	outer_ *InlineVariables
}

func NewInlineVariables_InliningBehavior(outer_ *InlineVariables) *InlineVariables_InliningBehavior {
	return &InlineVariables_InliningBehavior{outer_: outer_}
}

func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_canMoveModerately_3399(initNode, refNode pkg_rhino.Node) bool {
	var it *NodeIterators_LocalVarMotion = nil
	_ = it
	parent := initNode.GetParent()
	if parent.IsVar() {
		it = NodeIterators_LocalVarMotion_ForVar(this_.outer_.compiler, initNode, parent, parent.GetParent())
	} else if parent.IsAssign() {
		gp := parent.GetParent()
		it = NodeIterators_LocalVarMotion_ForAssign(this_.outer_.compiler, initNode, parent, gp, gp.GetParent())
	}
	for it.HasNext() {
		curNode := it.NextNode()
		if curNode == refNode {
			return true
		}
	}
	return false
}
`

	cccmSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w39b_opt/gen/j2gort"
	pkg_rhino "example.com/w39b_opt/gen/rhino"
)

type CompilerPass interface {
	Implements_jscomp_CompilerPass()
}

type JSChunk struct {
	Idx_ int32
}

type VarSym struct {
	Name_ string
}

type RefSym struct {
	V_      *VarSym
	Chunk_  *JSChunk
	IsDecl_ bool
}

type CrossChunkCodeMotion struct {
	compiler               AbstractCompiler
	chunkCount             int32
	chunkInsertionPointMap j2gort.OrderedMap[*JSChunk, pkg_rhino.Node]
	InputRefs              []RefSym
	ProcessedCycles        int32
	GuardedRefs            int32
	idHash_                atomic.Int32
}

func (*CrossChunkCodeMotion) Implements_jscomp_CompilerPass() {}

func NewCrossChunkCodeMotion(compiler AbstractCompiler, chunkCount int32) *CrossChunkCodeMotion {
	return &CrossChunkCodeMotion{
		compiler:               compiler,
		chunkCount:             chunkCount,
		chunkInsertionPointMap: j2gort.NewOrderedMap[*JSChunk, pkg_rhino.Node](),
	}
}

func (this_ *CrossChunkCodeMotion) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	func() {
		if this_.chunkCount > int32(1) {
			var val_2 CrossChunkCodeMotion_GlobalSymbolCollector
			val_2.outer_ = this_
			InitCrossChunkCodeMotion_GlobalSymbolCollector(&val_2)
			var collector *CrossChunkCodeMotion_GlobalSymbolCollector = &val_2
			collector.Collect(this_.InputRefs)
			this_.moveGlobalSymbols(collector)
			this_.countGuards(collector)
		}
	}()
}

func (this_ *CrossChunkCodeMotion) moveGlobalSymbols(collector *CrossChunkCodeMotion_GlobalSymbolCollector) {
	var val_7 CrossChunkCodeMotion_OrderAndCombineGlobalSymbols
	val_7.outer_ = this_
	InitCrossChunkCodeMotion_OrderAndCombineGlobalSymbols(&val_7, collector)
	orderer := &val_7
	orderer.OrderAndCombine()
	cur := orderer.stronglyConnectedSymbols.Iter()
	for cur.HasNext() {
		cyc := cur.Next()
		cyc.MoveDeclarationStatements()
		this_.ProcessedCycles++
	}
}

func (this_ *CrossChunkCodeMotion) countGuards(collector *CrossChunkCodeMotion_GlobalSymbolCollector) {
	cur := collector.symbolStack.Iter()
	for cur.HasNext() {
		gs := cur.Next()
		ic := gs.instanceofReferencesToGuard.Iter()
		for ic.HasNext() {
			ir := ic.Next()
			if ir.GetChunk() != nil {
				this_.GuardedRefs++
			}
		}
	}
}

type CrossChunkCodeMotion_GlobalSymbolCollector struct {
	outer_             *CrossChunkCodeMotion
	globalSymbolforVar j2gort.CustomOrderedMap[*VarSym, *CrossChunkCodeMotion_GlobalSymbol]
	symbolStack        j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]
}

func InitCrossChunkCodeMotion_GlobalSymbolCollector(this_ *CrossChunkCodeMotion_GlobalSymbolCollector) {
	this_.globalSymbolforVar = j2gort.NewCustomOrderedMap[*VarSym, *CrossChunkCodeMotion_GlobalSymbol]()
	this_.symbolStack = j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]{}
}

func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) getGlobalSymbol(v *VarSym) *CrossChunkCodeMotion_GlobalSymbol {
	gs := this_.globalSymbolforVar.Get(v)
	if gs == nil {
		gs = NewCrossChunkCodeMotion_GlobalSymbol(this_.outer_, v)
		this_.globalSymbolforVar.Put(v, gs)
		this_.symbolStack.Push(gs)
	}
	return gs
}

func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) Collect(refs []RefSym) {
	var prevDecl *CrossChunkCodeMotion_GlobalSymbol
	for _, r := range refs {
		gs := this_.getGlobalSymbol(r.V_)
		if r.IsDecl_ {
			dsg := gs.AddDeclarationStatement(r.Chunk_)
			if prevDecl != nil && prevDecl != gs {
				dsg.referencedGlobalSymbols.Add(prevDecl)
				prevDecl.referencingGlobalSymbols.Add(gs)
			}
			prevDecl = gs
		} else {
			gs.chunksWithImmovableReferences.Set_1(r.Chunk_.Idx_)
			gs.instanceofReferencesToGuard.Push(CrossChunkCodeMotion_InstanceofReference_Immovable(r.Chunk_))
		}
	}
}

type CrossChunkCodeMotion_GlobalSymbol struct {
	outer_                        *CrossChunkCodeMotion
	var_                          *VarSym
	dsgStack                      j2gort.Deque[*CrossChunkCodeMotion_DeclarationStatementGroup]
	chunksWithImmovableReferences *j2gort.Ext_java_util_BitSet
	referencingGlobalSymbols      j2gort.OrderedSet[*CrossChunkCodeMotion_GlobalSymbol]
	instanceofReferencesToGuard   j2gort.Deque[*CrossChunkCodeMotion_InstanceofReference]
	preorderNumber                int32
	assignedToScc                 bool
}

func NewCrossChunkCodeMotion_GlobalSymbol(outer_ *CrossChunkCodeMotion, var_ *VarSym) *CrossChunkCodeMotion_GlobalSymbol {
	this_ := &CrossChunkCodeMotion_GlobalSymbol{}
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_GlobalSymbol(this_, var_)
	return this_
}

func InitCrossChunkCodeMotion_GlobalSymbol(this_ *CrossChunkCodeMotion_GlobalSymbol, var_ *VarSym) {
	this_.dsgStack = j2gort.Deque[*CrossChunkCodeMotion_DeclarationStatementGroup]{}
	this_.chunksWithImmovableReferences = j2gort.NewExt_java_util_BitSet_1(this_.outer_.chunkCount)
	this_.referencingGlobalSymbols = j2gort.NewOrderedSet[*CrossChunkCodeMotion_GlobalSymbol]()
	this_.instanceofReferencesToGuard = j2gort.Deque[*CrossChunkCodeMotion_InstanceofReference]{}
	this_.preorderNumber = -1
	this_.var_ = var_
}

func (this_ *CrossChunkCodeMotion_GlobalSymbol) AddDeclarationStatement(chunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroup {
	dsg := NewCrossChunkCodeMotion_DeclarationStatementGroup(this_, chunk)
	this_.dsgStack.Push(dsg)
	return dsg
}

type CrossChunkCodeMotion_GlobalSymbolCycle struct {
	outer_  *CrossChunkCodeMotion
	symbols j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]
}

func NewCrossChunkCodeMotion_GlobalSymbolCycle(outer_ *CrossChunkCodeMotion) *CrossChunkCodeMotion_GlobalSymbolCycle {
	this_ := &CrossChunkCodeMotion_GlobalSymbolCycle{}
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_GlobalSymbolCycle(this_)
	return this_
}

func InitCrossChunkCodeMotion_GlobalSymbolCycle(this_ *CrossChunkCodeMotion_GlobalSymbolCycle) {
	this_.symbols = j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]{}
}

func (this_ *CrossChunkCodeMotion_GlobalSymbolCycle) MoveDeclarationStatements() {
	var cycles []*CrossChunkCodeMotion_DeclarationStatementGroupCycle
	cur := this_.symbols.Iter()
	for cur.HasNext() {
		sym := cur.Next()
		dc := sym.dsgStack.Iter()
		for dc.HasNext() {
			dsg := dc.Next()
			c := NewCrossChunkCodeMotion_DeclarationStatementGroupCycle(this_.outer_, dsg.currentChunk)
			c.dsgs.AddLast(dsg)
			cycles = append(cycles, c)
		}
	}
	_ = cycles
}

type CrossChunkCodeMotion_DeclarationStatementGroup struct {
	declaredGlobalSymbol    *CrossChunkCodeMotion_GlobalSymbol
	referencedGlobalSymbols j2gort.OrderedSet[*CrossChunkCodeMotion_GlobalSymbol]
	currentChunk            *JSChunk
}

func NewCrossChunkCodeMotion_DeclarationStatementGroup(declaredGlobalSymbol *CrossChunkCodeMotion_GlobalSymbol, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroup {
	this_ := &CrossChunkCodeMotion_DeclarationStatementGroup{}
	InitCrossChunkCodeMotion_DeclarationStatementGroup(this_, declaredGlobalSymbol, currentChunk)
	return this_
}

func InitCrossChunkCodeMotion_DeclarationStatementGroup(this_ *CrossChunkCodeMotion_DeclarationStatementGroup, declaredGlobalSymbol *CrossChunkCodeMotion_GlobalSymbol, currentChunk *JSChunk) {
	this_.referencedGlobalSymbols = j2gort.NewOrderedSet[*CrossChunkCodeMotion_GlobalSymbol]()
	this_.declaredGlobalSymbol = declaredGlobalSymbol
	this_.currentChunk = currentChunk
}

type CrossChunkCodeMotion_DeclarationStatementGroupCycle struct {
	outer_       *CrossChunkCodeMotion
	currentChunk *JSChunk
	dsgs         j2gort.Deque[*CrossChunkCodeMotion_DeclarationStatementGroup]
}

func NewCrossChunkCodeMotion_DeclarationStatementGroupCycle(outer_ *CrossChunkCodeMotion, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroupCycle {
	this_ := &CrossChunkCodeMotion_DeclarationStatementGroupCycle{}
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_DeclarationStatementGroupCycle(this_, currentChunk)
	return this_
}

func InitCrossChunkCodeMotion_DeclarationStatementGroupCycle(this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle, currentChunk *JSChunk) {
	this_.currentChunk = currentChunk
	this_.dsgs = j2gort.Deque[*CrossChunkCodeMotion_DeclarationStatementGroup]{}
}

type CrossChunkCodeMotion_OrderAndCombineGlobalSymbols struct {
	outer_                   *CrossChunkCodeMotion
	collector                *CrossChunkCodeMotion_GlobalSymbolCollector
	componentContents        j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]
	componentRoots           j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]
	stronglyConnectedSymbols j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbolCycle]
	preorderCounter          int32
}

func InitCrossChunkCodeMotion_OrderAndCombineGlobalSymbols(this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols, collector *CrossChunkCodeMotion_GlobalSymbolCollector) {
	this_.collector = collector
	this_.stronglyConnectedSymbols = j2gort.NewDequeCap[*CrossChunkCodeMotion_GlobalSymbolCycle](collector.symbolStack.Size())
}

func (this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) OrderAndCombine() {
	cur := this_.collector.symbolStack.Iter()
	for cur.HasNext() {
		gs := cur.Next()
		if gs.preorderNumber < 0 {
			this_.ProcessGlobalSymbol(gs)
		}
	}
}

func (this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) ProcessGlobalSymbol(symbol *CrossChunkCodeMotion_GlobalSymbol) {
	symbol.preorderNumber = this_.preorderCounter
	this_.preorderCounter++
	this_.componentRoots.Push(symbol)
	this_.componentContents.Push(symbol)
	cur := symbol.referencingGlobalSymbols.Iter()
	for cur.HasNext() {
		ref := cur.Next()
		if ref.preorderNumber < 0 {
			this_.ProcessGlobalSymbol(ref)
		} else if !ref.assignedToScc {
			for this_.componentRoots.Peek().preorderNumber > ref.preorderNumber {
				this_.componentRoots.Pop()
			}
		}
	}
	if this_.componentRoots.Peek() == symbol {
		this_.componentRoots.Pop()
		var val_46 CrossChunkCodeMotion_GlobalSymbolCycle
		val_46.outer_ = this_.outer_
		InitCrossChunkCodeMotion_GlobalSymbolCycle(&val_46)
		var cycle *CrossChunkCodeMotion_GlobalSymbolCycle = &val_46
		for {
			connected := this_.componentContents.Pop()
			cycle.symbols.AddLast(connected)
			connected.assignedToScc = true
			if connected == symbol {
				break
			}
		}
		this_.stronglyConnectedSymbols.AddLast(cycle)
	}
}

type CrossChunkCodeMotion_InstanceofReference struct {
	chunk *JSChunk
}

func NewCrossChunkCodeMotion_InstanceofReference(chunk *JSChunk) *CrossChunkCodeMotion_InstanceofReference {
	this_ := &CrossChunkCodeMotion_InstanceofReference{}
	InitCrossChunkCodeMotion_InstanceofReference(this_, chunk)
	return this_
}

func InitCrossChunkCodeMotion_InstanceofReference(this_ *CrossChunkCodeMotion_InstanceofReference, chunk *JSChunk) {
	this_.chunk = chunk
}

func CrossChunkCodeMotion_InstanceofReference_Immovable(chunk *JSChunk) *CrossChunkCodeMotion_InstanceofReference {
	return NewCrossChunkCodeMotion_InstanceofReference(chunk)
}

func (this_ *CrossChunkCodeMotion_InstanceofReference) GetChunk() *JSChunk {
	return this_.chunk
}
`

	for _, pkgPrefix := range []string{"orig", "opt"} {
		mustWrite(pkgPrefix+"/jscomp/global_namespace.go", globalNsSrc)
		mustWrite(pkgPrefix+"/jscomp/node_iterators.go", nodeIteratorsSrc)
		mustWrite(pkgPrefix+"/jscomp/cross_chunk_code_motion.go", cccmSrc)
		mustWrite(pkgPrefix+"/jscomp/export_shim.go", `package jscomp

import pkg_rhino "example.com/w39b_opt/gen/rhino"

func (b *InlineVariables_InliningBehavior) PrivMoveCheck(initNode, refNode pkg_rhino.Node) bool {
	return b.priv_InliningBehavior_canMoveModerately_3399(initNode, refNode)
}
`)
	}

	optDir := filepath.Join(tmpDir, "opt/jscomp")
	allPasses := passes.AllPasses()

	s1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero Wave 39B OPT sites on Run 1, got %+v", s1)
	}

	// Run 2 MUST produce strictly 0 sites and 0 modified files!
	s2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d ByPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	gotGn, err := os.ReadFile(filepath.Join(optDir, "global_namespace.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"scratchObjLitKeyAnalysis_ GlobalNamespace_ObjLitStringKeyAnalysis",
		"allocScratch_GlobalNamespace_ObjLitStringKeyAnalysis",
		"allocScratchVia_GlobalNamespace_ObjLitStringKeyAnalysis_ForObjLitAssignment(&this_.scratchObjLitKeyAnalysis_",
		"allocScratchVia_GlobalNamespace_ObjLitStringKeyAnalysis_ForObjectDefineProperty(&this_.scratchObjLitKeyAnalysis_",
		"allocScratchVia_GlobalNamespace_ObjLitStringKeyAnalysis_ForNonReference(&this_.scratchObjLitKeyAnalysis_",
	} {
		if !strings.Contains(string(gotGn), want) {
			t.Errorf("global_namespace.go missing %q in:\n%s", want, string(gotGn))
		}
	}

	gotNi, err := os.ReadFile(filepath.Join(optDir, "node_iterators.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"scratchLocalScope_     NodeIterators_FunctionlessLocalScope",
		"scratchLocalVarMotion_ NodeIterators_LocalVarMotion",
		"initReuse_NodeIterators_FunctionlessLocalScope_3",
		"allocReuseVia_NodeIterators_LocalVarMotion_ForVar(&this_.scratchLocalVarMotion_, &this_.scratchLocalScope_",
		"allocReuseVia_NodeIterators_LocalVarMotion_ForAssign(&this_.scratchLocalVarMotion_, &this_.scratchLocalScope_",
	} {
		if !strings.Contains(string(gotNi), want) {
			t.Errorf("node_iterators.go missing %q in:\n%s", want, string(gotNi))
		}
	}

	gotCccm, err := os.ReadFile(filepath.Join(optDir, "cross_chunk_code_motion.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"slab_GlobalSymbol",
		"slab_GlobalSymbolCycle",
		"slab_DeclarationStatementGroup",
		"slab_DeclarationStatementGroupCycle",
		"slab_InstanceofReference",
		"allocPooled_CrossChunkCodeMotion_GlobalSymbol(&this_.outer_.slab_GlobalSymbol",
		"allocPooled_CrossChunkCodeMotion_GlobalSymbolCycle(&this_.outer_.slab_GlobalSymbolCycle",
		"defer this_.ReleaseSlabs_()",
		".ReleaseWith(",
	} {
		if !strings.Contains(string(gotCccm), want) {
			t.Errorf("cross_chunk_code_motion.go missing %q in:\n%s", want, string(gotCccm))
		}
	}

	diffTest := `package opt_diff_test

import (
	"testing"

	"example.com/w39b_opt/gen/j2gort"
	pkg_rhino "example.com/w39b_opt/gen/rhino"
	opt "example.com/w39b_opt/opt/jscomp"
	orig "example.com/w39b_opt/orig/jscomp"
)

func TestOptLoweringsDifferentialParity(t *testing.T) {
	// 1. GlobalNamespace_ObjLitStringKeyAnalysis scratch reuse parity
	pkg_rhino.ResetTree()
	pkg_rhino.DefNode(10, "foo", false, false, false, false, 0, 0, 0)
	pkg_rhino.DefNode(11, "bar", false, false, false, false, 0, 0, 0)

	ogn := &orig.GlobalNamespace{}
	ob := orig.NewGlobalNamespace_BuildGlobalNamespace(ogn)
	ob.CollectKey(10, 1)
	ob.CollectKey(11, 2)
	ob.CollectKey(11, 3)

	xgn := &opt.GlobalNamespace{}
	xb := opt.NewGlobalNamespace_BuildGlobalNamespace(xgn)
	xb.CollectKey(10, 1)
	xb.CollectKey(11, 2)
	xb.CollectKey(11, 3)

	if len(ogn.SeenNames) != len(xgn.SeenNames) || len(ogn.SeenKinds) != len(xgn.SeenKinds) {
		t.Fatalf("GlobalNamespace SeenNames/SeenKinds len mismatch: %v vs %v", ogn, xgn)
	}
	for i := range ogn.SeenNames {
		if ogn.SeenNames[i] != xgn.SeenNames[i] || ogn.SeenKinds[i] != xgn.SeenKinds[i] {
			t.Fatalf("GlobalNamespace Seen mismatch at %d: orig=(%q,%d) opt=(%q,%d)", i, ogn.SeenNames[i], ogn.SeenKinds[i], xgn.SeenNames[i], xgn.SeenKinds[i])
		}
	}

	// 2. NodeIterators_FunctionlessLocalScope + LocalVarMotion reuse parity
	// Build: block(1) -> var(2, next=expr(4)) -> initName(3); expr(4) -> targetName(5)
	pkg_rhino.ResetTree()
	pkg_rhino.DefNode(1, "", false, false, false, false, 0, 2, 0)
	pkg_rhino.DefNode(2, "", false, true, false, false, 1, 3, 4)
	pkg_rhino.DefNode(3, "v", false, false, false, false, 2, 0, 0)
	pkg_rhino.DefNode(4, "", false, false, false, true, 1, 5, 0)
	pkg_rhino.DefNode(5, "v", false, false, false, false, 4, 0, 0)
	pkg_rhino.DefNode(6, "other", false, false, false, false, 4, 0, 0)

	oBeh := orig.NewInlineVariables_InliningBehavior(orig.NewInlineVariables(&orig.MockCompiler{}))
	xBeh := opt.NewInlineVariables_InliningBehavior(opt.NewInlineVariables(&opt.MockCompiler{}))
	for iter := 0; iter < 5; iter++ {
		if oBeh.PrivMoveCheck(3, 5) != xBeh.PrivMoveCheck(3, 5) {
			t.Fatalf("iter %d PrivMoveCheck(3,5) mismatch", iter)
		}
		if oBeh.PrivMoveCheck(3, 6) != xBeh.PrivMoveCheck(3, 6) {
			t.Fatalf("iter %d PrivMoveCheck(3,6) mismatch", iter)
		}
	}

	// 3. CrossChunkCodeMotion PooledSlabs + ReleaseWith parity
	j2gort.PooledSlabAllocCount.Store(0)
	j2gort.PooledSlabReleaseCount.Store(0)

	oc0, oc1 := &orig.JSChunk{Idx_: 0}, &orig.JSChunk{Idx_: 1}
	ov1, ov2 := &orig.VarSym{Name_: "a"}, &orig.VarSym{Name_: "b"}
	occcm := orig.NewCrossChunkCodeMotion(&orig.MockCompiler{}, 2)
	occcm.InputRefs = []orig.RefSym{
		{V_: ov1, Chunk_: oc0, IsDecl_: true},
		{V_: ov2, Chunk_: oc1, IsDecl_: true},
		{V_: ov1, Chunk_: oc1, IsDecl_: false},
	}
	occcm.Process_2(0, 0)

	xc0, xc1 := &opt.JSChunk{Idx_: 0}, &opt.JSChunk{Idx_: 1}
	xv1, xv2 := &opt.VarSym{Name_: "a"}, &opt.VarSym{Name_: "b"}
	xcccm := opt.NewCrossChunkCodeMotion(&opt.MockCompiler{}, 2)
	xcccm.InputRefs = []opt.RefSym{
		{V_: xv1, Chunk_: xc0, IsDecl_: true},
		{V_: xv2, Chunk_: xc1, IsDecl_: true},
		{V_: xv1, Chunk_: xc1, IsDecl_: false},
	}
	xcccm.Process_2(0, 0)

	if occcm.ProcessedCycles != xcccm.ProcessedCycles || occcm.GuardedRefs != xcccm.GuardedRefs {
		t.Fatalf("CCCM mismatch: orig=(%d,%d) opt=(%d,%d)", occcm.ProcessedCycles, occcm.GuardedRefs, xcccm.ProcessedCycles, xcccm.GuardedRefs)
	}
	if j2gort.PooledSlabAllocCount.Load() == 0 {
		t.Fatalf("expected non-zero PooledSlab allocations in CCCM")
	}
	if j2gort.PooledSlabReleaseCount.Load() != 5 {
		t.Fatalf("expected 5 PooledSlab releases in CCCM ReleaseSlabs_, got %d", j2gort.PooledSlabReleaseCount.Load())
	}
}
`
	mustWrite("opt_diff_test.go", diffTest)
	cmd := exec.Command("go", "test", "-race", "-count=1", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("OPT differential runtime test failed: %v\n%s", err, string(out))
	}
}

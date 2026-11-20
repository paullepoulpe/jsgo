package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func wave45bWriteFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func wave45bRunAllAndCheckIdempotent(t *testing.T, pkgDir string) *passes.Stats {
	t.Helper()
	allPasses := passes.AllPasses()
	s1, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("AllPasses Run 1 on %s failed: %v", pkgDir, err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero sites on Run 1 for %s", pkgDir)
	}
	s2, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("AllPasses Run 2 on %s failed: %v", pkgDir, err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("AllPasses Run 2 on %s NOT idempotent: TotalSites=%d FilesModified=%d byPass=%+v", pkgDir, s2.TotalSites, s2.FilesModified, s2.ByPass)
	}
	return s1
}

const wave45bJ2gortStub = `package j2gort

import "sync/atomic"

const NullString = "\x00NULL\x00"

type ClinitState struct{}

func EnterClinit(*ClinitState, *uint32) bool { return true }
func ExitClinit(*ClinitState, *uint32)       {}
func RegisterTypePredicate(string, string, func(any) bool) {}
func InstanceOf[T any](v any) bool { _, ok := v.(T); return ok }
func AssignIdentHash(h *atomic.Int32) int32 { h.Store(1); return 1 }
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
func CheckNotNullPtr[T any](p *T) *T {
	if p == nil {
		panic("NPE")
	}
	return p
}
func CheckNotNullNode[T ~uint32](n T) T {
	if n == 0 {
		panic("NPE")
	}
	return n
}
func CheckArgumentMsg(b bool, _ any) {
	if !b {
		panic("IAE")
	}
}
func CheckStateMsg(b bool, _ any) {
	if !b {
		panic("ISE")
	}
}
func BoxNode[T ~uint32](n T) any { return n }

var (
	PooledSlabAllocs   atomic.Int64
	PooledSlabReleases atomic.Int64
	CustomMapOps       atomic.Int64
)

func ResetCounters() {
	PooledSlabAllocs.Store(0)
	PooledSlabReleases.Store(0)
	CustomMapOps.Store(0)
}

type Slab[T any] struct{ items []T }

func (s *Slab[T]) Alloc() *T { return new(T) }
func (s *Slab[T]) Reset()    { s.items = nil }
func (s *Slab[T]) Release()  { s.items = nil }

type PooledSlab[T any] struct {
	Allocs   int32
	Releases int32
}

func (s *PooledSlab[T]) Alloc() *T {
	s.Allocs++
	PooledSlabAllocs.Add(1)
	return new(T)
}
func (s *PooledSlab[T]) Reset() {}
func (s *PooledSlab[T]) Release() {
	s.Releases++
	PooledSlabReleases.Add(1)
}
func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	_ = fn
	s.Releases++
	PooledSlabReleases.Add(1)
}

type OrderedMap[K comparable, V any] struct {
	m map[K]V
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V)}
}
func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V)}
}
func (o *OrderedMap[K, V]) Get(k K) V {
	if o == nil || o.m == nil {
		var z V
		return z
	}
	return o.m[k]
}
func (o *OrderedMap[K, V]) Put(k K, v V) V {
	if o.m == nil {
		o.m = make(map[K]V)
	}
	old := o.m[k]
	o.m[k] = v
	return old
}
func (o *OrderedMap[K, V]) Clear()         { clear(o.m) }
func (o *OrderedMap[K, V]) ReleaseToPool() { clear(o.m) }

type OrderedMapIter[V any] struct {
	vals []V
	idx  int
}

func (it *OrderedMapIter[V]) HasNext() bool { return it.idx < len(it.vals) }
func (it *OrderedMapIter[V]) Next() V {
	v := it.vals[it.idx]
	it.idx++
	return v
}
func (o *OrderedMap[K, V]) IterValues() OrderedMapIter[V] {
	vals := make([]V, 0, len(o.m))
	for _, v := range o.m {
		vals = append(vals, v)
	}
	return OrderedMapIter[V]{vals: vals}
}
func ReleaseOrderedMap[K comparable, V any](m *OrderedMap[K, V]) {
	if m != nil {
		m.ReleaseToPool()
	}
}

type CustomOrderedMap[K comparable, V any] struct {
	m map[K]V
}

func NewCustomOrderedMap[K comparable, V any]() CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{m: make(map[K]V)}
}
func (c *CustomOrderedMap[K, V]) Get(k K) V {
	CustomMapOps.Add(1)
	if c == nil || c.m == nil {
		var z V
		return z
	}
	return c.m[k]
}
func (c *CustomOrderedMap[K, V]) Put(k K, v V) V {
	CustomMapOps.Add(1)
	if c.m == nil {
		c.m = make(map[K]V)
	}
	old := c.m[k]
	c.m[k] = v
	return old
}
`

const wave45bRhinoStub = `package rhino

type Node uint32

var nodeOnlyChild = map[Node]Node{}
var nodeIsNameDecl = map[Node]bool{}

func RegisterNode(n Node, onlyChild Node, isNameDecl bool) {
	nodeOnlyChild[n] = onlyChild
	nodeIsNameDecl[n] = isNameDecl
}

func (n Node) GetOnlyChild() Node {
	c, ok := nodeOnlyChild[n]
	if !ok {
		panic("GetOnlyChild: node has no unique child")
	}
	return c
}
`

func TestWave45BAllOptimizations(t *testing.T) {
	tmpDir := t.TempDir()

	jscompFiles := map[string]string{
		"go.mod":                 "module example.com/w45b\n\ngo 1.23\n",
		"gen/j2gort/j2gort.go":   wave45bJ2gortStub,
		"gen/rhino/node.go":      wave45bRhinoStub,
		"opt/jscomp/scope.go": `package jscomp

import (
	"sync/atomic"

	"example.com/w45b/gen/j2gort"
	pkg_rhino "example.com/w45b/gen/rhino"
)

type AbstractScope interface {
	IsJavaNil_() bool
}

type AbstractScopeBase struct {
	self_        AbstractScope
	implicitVars any
	rootNode     pkg_rhino.Node
	idHash_      atomic.Int32
}

func (this_ *AbstractScopeBase) BindSelf_(s any) {
	this_.self_ = s.(AbstractScope)
}
func (this_ *AbstractScopeBase) IsJavaNil_() bool { return this_ == nil }
func (this_ *AbstractScopeBase) CheckRootScope()  {}

func InitAbstractScope(this_ *AbstractScopeBase, rootNode pkg_rhino.Node) {
	this_.implicitVars = nil
	this_.rootNode = j2gort.CheckNotNullNode(rootNode)
}

type Var struct {
	name string
}

// Start in Wave 44B layout (varSlab_ at offset 104, Cache Line 1) so Wave 45B moves it to offset 48 (Cache Line 0).
type Scope struct {
	AbstractScopeBase
	parent   *Scope
	var0     *Var
	var1     *Var
	var2     *Var
	var3     *Var
	var4     *Var
	var5     *Var
	vars     *j2gort.OrderedMap[string, *Var]
	varSlab_ *j2gort.PooledSlab[Var]
	depth    int32
}

func (this_ *Scope) IsJavaNil_() bool { return this_ == nil }
`,
		"opt/jscomp/type_inference.go": `package jscomp

import (
	"example.com/w45b/gen/j2gort"
	pkg_rhino "example.com/w45b/gen/rhino"
)

type TypedScope struct {
	AbstractScopeBase
	parent   *TypedScope
	depth    int32
	isBottom bool
}

func (this_ *TypedScope) BindSelf_(s any) {
	this_.AbstractScopeBase.BindSelf_(s)
}
func (this_ *TypedScope) IsJavaNil_() bool { return this_ == nil }
func (this_ *TypedScope) IsBottom() bool   { return this_.isBottom }
func (this_ *TypedScope) initFields_TypedScope() {
	_ = this_
}

func NewTypedScope_2_Node_bool(rootNode pkg_rhino.Node, isBottom bool) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_2_Node_bool(this_, rootNode, isBottom)
	return this_
}

func allocPooled_TypedScope_2_Node_bool(slab_ *j2gort.PooledSlab[TypedScope], rootNode pkg_rhino.Node, isBottom bool) *TypedScope {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitTypedScope_2_Node_bool(this_, rootNode, isBottom)
	return this_
}

func InitTypedScope_2_Node_bool(this_ *TypedScope, rootNode pkg_rhino.Node, isBottom bool) {
	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_TypedScope()
	this_.CheckRootScope()
	this_.parent = nil
	this_.depth = int32(0)
	this_.isBottom = isBottom
}

func TypedScope_CreateLatticeBottom(rootNode pkg_rhino.Node) *TypedScope {
	return NewTypedScope_2_Node_bool(rootNode, true)
}

func allocPooledVia_TypedScope_CreateLatticeBottom(slab_ *j2gort.PooledSlab[TypedScope], rootNode pkg_rhino.Node) *TypedScope {
	return allocPooled_TypedScope_2_Node_bool(slab_, rootNode, true)
}

type typedScopeArena_ struct {
	slab_TypedScope j2gort.PooledSlab[TypedScope]
}

type TypedScopeCreator struct {
	tsArena_ *typedScopeArena_
}

type LinkedFlowScope_Pool struct {
	resets int32
}

func (p *LinkedFlowScope_Pool) Reset() { p.resets++ }

type LinkedFlowScope struct {
	functionScope *TypedScope
}

func LinkedFlowScope_CreateEntryLattice_3(scope *TypedScope, pool *LinkedFlowScope_Pool) *LinkedFlowScope {
	_ = pool
	return &LinkedFlowScope{functionScope: scope}
}

type TypeInference struct {
	scopeCreator  *TypedScopeCreator
	bottomScope   *LinkedFlowScope
	flowScopePool *LinkedFlowScope_Pool
}

func NewTypeInference(sc *TypedScopeCreator, root pkg_rhino.Node) *TypeInference {
	this_ := &TypeInference{}
	InitTypeInference(this_, sc, root)
	return this_
}

func InitTypeInference(this_ *TypeInference, sc *TypedScopeCreator, root pkg_rhino.Node) {
	this_.scopeCreator = sc
	this_.flowScopePool = &LinkedFlowScope_Pool{}
	// Shape 1: already-lowered Wave 40B call into tsArena_.slab_TypedScope
	this_.bottomScope = LinkedFlowScope_CreateEntryLattice_3(
		allocPooledVia_TypedScope_CreateLatticeBottom(&this_.scopeCreator.tsArena_.slab_TypedScope, root),
		this_.flowScopePool,
	)
}

func (this_ *TypeInference) ResetForScope(root pkg_rhino.Node) {
	this_.flowScopePool.Reset()
	// Shape 2: freshly translated TypedScope_CreateLatticeBottom call
	this_.bottomScope = LinkedFlowScope_CreateEntryLattice_3(
		TypedScope_CreateLatticeBottom(root),
		this_.flowScopePool,
	)
}

func (this_ *TypeInference) FinishScope() {
	this_.bottomScope = nil
	this_.flowScopePool.Reset()
}
`,
		"opt/jscomp/global_namespace.go": `package jscomp

import (
	"example.com/w45b/gen/j2gort"
	pkg_rhino "example.com/w45b/gen/rhino"
)

type GlobalNamespace_SourceKind uint8

const (
	gnSourceKind_CODE   GlobalNamespace_SourceKind = 1
	gnSourceKind_EXTERN GlobalNamespace_SourceKind = 2
)

type GlobalNamespace_Ref struct {
	node pkg_rhino.Node
}

type GlobalNamespace struct {
	nameMap                  j2gort.OrderedMap[string, *GlobalNamespace_Name]
	globalNames              []*GlobalNamespace_Name
	sourceKind               GlobalNamespace_SourceKind
	slab_GlobalNamespace_Ref j2gort.PooledSlab[GlobalNamespace_Ref]
}

func NewGlobalNamespace() *GlobalNamespace {
	return &GlobalNamespace{
		nameMap:    j2gort.AcquireOrderedMap[string, *GlobalNamespace_Name](),
		sourceKind: gnSourceKind_CODE,
	}
}

func (this_ *GlobalNamespace) CreateNameForTesting(name string) *GlobalNamespace_Name {
	return NewGlobalNamespace_Name(this_, name, nil, gnSourceKind_CODE)
}

func (this_ *GlobalNamespace) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.globalNames = nil
	this_.nameMap.ReleaseToPool()
	this_.slab_GlobalNamespace_Ref.Release()
}

type GlobalNamespace_BuildGlobalNamespace struct {
	outer_ *GlobalNamespace
}

func (this_ *GlobalNamespace_BuildGlobalNamespace) GetOrCreateName(name string) *GlobalNamespace_Name {
	if existing := this_.outer_.nameMap.Get(name); existing != nil {
		return existing
	}
	var node *GlobalNamespace_Name
	if dot := lastDotIdx(name); dot >= 0 {
		parent := this_.GetOrCreateName(name[:dot])
		node = parent.AddProperty(name[dot+1:], this_.outer_.sourceKind)
		this_.outer_.nameMap.Put(name, node)
	} else {
		node = NewGlobalNamespace_Name(this_.outer_, name, nil, this_.outer_.sourceKind)
		this_.outer_.globalNames = append(this_.outer_.globalNames, node)
		this_.outer_.nameMap.Put(name, node)
	}
	return node
}

func lastDotIdx(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return i
		}
	}
	return -1
}

type GlobalNamespace_Name struct {
	outer_     *GlobalNamespace
	baseName   string
	parent     *GlobalNamespace_Name
	fullName   string
	prop0      *GlobalNamespace_Name
	propCount  int32
	sourceKind GlobalNamespace_SourceKind
}

func (this_ *GlobalNamespace_Name) initFields_GlobalNamespace_Name() {
	this_.baseName = j2gort.NullString
	this_.fullName = j2gort.NullString
}

func NewGlobalNamespace_Name(outer_ *GlobalNamespace, name string, parent *GlobalNamespace_Name, sourceKind GlobalNamespace_SourceKind) *GlobalNamespace_Name {
	this_ := &GlobalNamespace_Name{}
	this_.outer_ = outer_
	InitGlobalNamespace_Name(this_, name, parent, sourceKind)
	return this_
}

func InitGlobalNamespace_Name(this_ *GlobalNamespace_Name, name string, parent *GlobalNamespace_Name, sourceKind GlobalNamespace_SourceKind) {
	this_.initFields_GlobalNamespace_Name()
	this_.baseName = name
	this_.parent = parent
	if parent == nil {
		this_.fullName = name
	} else {
		this_.fullName = j2gort.NullString
	}
	this_.sourceKind = sourceKind
}

func (this_ *GlobalNamespace_Name) AddProperty(name string, sourceKind GlobalNamespace_SourceKind) *GlobalNamespace_Name {
	_ = name
	_ = sourceKind
	var val_10 GlobalNamespace_Name
	val_10.outer_ = this_.outer_
	InitGlobalNamespace_Name(&val_10, name, this_, sourceKind)
	var node *GlobalNamespace_Name = &val_10
	_ = node
	this_.prop0 = node
	this_.propCount++
	return node
}
`,
		"opt/jscomp/cross_chunk_code_motion.go": `package jscomp

import "example.com/w45b/gen/j2gort"

type CrossChunkCodeMotion struct{}

type CrossChunkCodeMotion_GlobalSymbol struct {
	var_ *Var
}

type CrossChunkCodeMotion_GlobalSymbolCollector struct {
	outer_             *CrossChunkCodeMotion
	globalSymbolforVar j2gort.CustomOrderedMap[*Var, *CrossChunkCodeMotion_GlobalSymbol]
}

func NewCrossChunkCodeMotion_GlobalSymbolCollector(outer_ *CrossChunkCodeMotion) *CrossChunkCodeMotion_GlobalSymbolCollector {
	this_ := &CrossChunkCodeMotion_GlobalSymbolCollector{outer_: outer_}
	this_.globalSymbolforVar = j2gort.NewCustomOrderedMap[*Var, *CrossChunkCodeMotion_GlobalSymbol]()
	return this_
}

func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) GetGlobalSymbol(v *Var) *CrossChunkCodeMotion_GlobalSymbol {
	sym := this_.globalSymbolforVar.Get(v)
	if sym == nil {
		sym = &CrossChunkCodeMotion_GlobalSymbol{var_: v}
		this_.globalSymbolforVar.Put(v, sym)
	}
	return sym
}
`,
		"opt/jscomp/remove_unused_code.go": `package jscomp

import (
	"example.com/w45b/gen/j2gort"
	pkg_rhino "example.com/w45b/gen/rhino"
)

func NodeUtil_IsNameDeclaration(n pkg_rhino.Node) bool {
	return n != 0 && (n == 100 || n == 102)
}

type RemoveUnusedCode struct {
	slab_RemoveUnusedCode_NameDeclarationStatement j2gort.PooledSlab[RemoveUnusedCode_NameDeclarationStatement]
}

type RemoveUnusedCode_Removable interface {
	IsJavaNil_() bool
}

type RemoveUnusedCode_Continuation struct {
	node pkg_rhino.Node
}

type RemoveUnusedCode_RemovableBase struct {
	self_                           RemoveUnusedCode_Removable
	outer_                          *RemoveUnusedCode
	propertyName                    string
	targetNode                      pkg_rhino.Node
	continuation0                   RemoveUnusedCode_Continuation
	continuation1                   RemoveUnusedCode_Continuation
	restContinuations               []RemoveUnusedCode_Continuation
	isPrototypeDotPropertyReference bool
	isThisDotPropertyReference      bool
}

func (this_ *RemoveUnusedCode_RemovableBase) IsJavaNil_() bool { return this_ == nil }

func (this_ *RemoveUnusedCode_RemovableBase) initFields_RemoveUnusedCode_Removable() {
	this_.propertyName = j2gort.NullString
}

func InitRemoveUnusedCode_Removable(this_ *RemoveUnusedCode_RemovableBase, targetNode pkg_rhino.Node, builder *RemoveUnusedCode_RemovableBuilder) {
	this_.initFields_RemoveUnusedCode_Removable()
	this_.continuation0 = builder.continuation0
	this_.continuation1 = builder.continuation1
	this_.restContinuations = builder.restContinuations
	this_.propertyName = builder.propertyName
	this_.isPrototypeDotPropertyReference = builder.isPrototypeDotPropertyReference
	this_.isThisDotPropertyReference = builder.isThisDotPropertyReference
	this_.targetNode = targetNode
}

type RemoveUnusedCode_RemovableBuilder struct {
	outer_                          *RemoveUnusedCode
	propertyName                    string
	continuation0                   RemoveUnusedCode_Continuation
	continuation1                   RemoveUnusedCode_Continuation
	restContinuations               []RemoveUnusedCode_Continuation
	isPrototypeDotPropertyReference bool
	isThisDotPropertyReference      bool
}

func (this_ *RemoveUnusedCode_RemovableBuilder) priv_RemovableBuilder_finish_4772(x any) {
	_ = x
	this_.propertyName = j2gort.NullString
	this_.continuation0 = RemoveUnusedCode_Continuation{}
	this_.continuation1 = RemoveUnusedCode_Continuation{}
	this_.restContinuations = nil
	this_.isPrototypeDotPropertyReference = false
	this_.isThisDotPropertyReference = false
}

func (this_ *RemoveUnusedCode_RemovableBuilder) BuildNameDeclarationStatement(declarationStatement pkg_rhino.Node) *RemoveUnusedCode_NameDeclarationStatement {
	_ = declarationStatement
	rucOuter_ := this_.outer_
	res_ := rucOuter_.slab_RemoveUnusedCode_NameDeclarationStatement.Alloc()
	res_.self_ = res_
	res_.outer_ = rucOuter_
	InitRemoveUnusedCode_NameDeclarationStatement(res_, this_, declarationStatement)
	this_.priv_RemovableBuilder_finish_4772(nil)
	return res_
}

type RemoveUnusedCode_NameDeclarationStatement struct {
	RemoveUnusedCode_RemovableBase
	outer_               *RemoveUnusedCode
	declarationStatement pkg_rhino.Node
}

func (this_ *RemoveUnusedCode_NameDeclarationStatement) initFields_RemoveUnusedCode_NameDeclarationStatement() {
	_ = this_
}

func allocPooled_RemoveUnusedCode_NameDeclarationStatement(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, declarationStatement pkg_rhino.Node) *RemoveUnusedCode_NameDeclarationStatement {
	this_ := outer_.slab_RemoveUnusedCode_NameDeclarationStatement.Alloc()
	this_.self_ = this_
	this_.outer_ = outer_
	InitRemoveUnusedCode_NameDeclarationStatement(this_, builder, declarationStatement)
	return this_
}

func InitRemoveUnusedCode_NameDeclarationStatement(this_ *RemoveUnusedCode_NameDeclarationStatement, builder *RemoveUnusedCode_RemovableBuilder, declarationStatement pkg_rhino.Node) {
	_ = builder
	_ = declarationStatement
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, declarationStatement.GetOnlyChild(), builder)
	this_.initFields_RemoveUnusedCode_NameDeclarationStatement()
	if !(NodeUtil_IsNameDeclaration(declarationStatement)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(declarationStatement))
	}
	this_.declarationStatement = declarationStatement
}
`,
		"opt/jscomp/w45b_runtime_test.go": `package jscomp

import (
	"reflect"
	"testing"

	"example.com/w45b/gen/j2gort"
	pkg_rhino "example.com/w45b/gen/rhino"
)

func assertPanics(t *testing.T, label string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s: expected panic, got none", label)
		}
	}()
	fn()
}

func TestWave45BRuntimeParity(t *testing.T) {
	// 1. Task 2: Scope Cache-Line 0 layout: AbstractScopeBase (0..40), parent (40..48), varSlab_ (48..56), var0 (56..64), ..., depth (112..116).
	scopeT := reflect.TypeFor[Scope]()
	if reflect.TypeFor[*Scope]().Size() == 8 {
		order := []string{"AbstractScopeBase", "parent", "varSlab_", "var0", "var1", "var2", "var3", "var4", "var5", "vars", "depth"}
		var prevEnd uintptr
		for i, fname := range order {
			sf, ok := scopeT.FieldByName(fname)
			if !ok {
				t.Fatalf("missing Scope.%s", fname)
			}
			if i > 0 && sf.Offset != prevEnd {
				t.Fatalf("padding detected before Scope.%s: offset=%d prevEnd=%d", fname, sf.Offset, prevEnd)
			}
			prevEnd = sf.Offset + sf.Type.Size()
		}
		for _, cl0 := range []string{"AbstractScopeBase", "parent", "varSlab_", "var0"} {
			sf, _ := scopeT.FieldByName(cl0)
			if sf.Offset+sf.Type.Size() > 64 {
				t.Fatalf("Scope.%s ends at %d, want <= 64 (Cache Line 0)", cl0, sf.Offset+sf.Type.Size())
			}
		}
	}

	// 2. Task 1: TypeInference embedded bottomTypedScope_ reuse (ZERO tsArena_.slab_TypedScope allocations!).
	j2gort.ResetCounters()
	tsc := &TypedScopeCreator{tsArena_: &typedScopeArena_{}}
	ti := NewTypeInference(tsc, pkg_rhino.Node(10))
	if ti.bottomScope == nil || ti.bottomScope.functionScope != &ti.bottomTypedScope_ ||
		!ti.bottomScope.functionScope.IsBottom() || ti.bottomScope.functionScope.rootNode != 10 ||
		ti.bottomScope.functionScope.self_ != &ti.bottomTypedScope_ {
		t.Fatalf("unexpected initial bottomScope: %+v", ti.bottomScope)
	}
	ti.FinishScope()
	if ti.bottomScope != nil || ti.bottomTypedScope_.rootNode != 0 || ti.bottomTypedScope_.isBottom || ti.bottomTypedScope_.self_ != nil {
		t.Fatalf("FinishScope did not clear bottomTypedScope_: %+v", ti.bottomTypedScope_)
	}
	for i := pkg_rhino.Node(20); i < 30; i++ {
		ti.ResetForScope(i)
		if ti.bottomScope.functionScope != &ti.bottomTypedScope_ || !ti.bottomScope.functionScope.IsBottom() || ti.bottomScope.functionScope.rootNode != i {
			t.Fatalf("ResetForScope(%d) failed: %+v", i, ti.bottomScope.functionScope)
		}
		ti.FinishScope()
	}
	if got := tsc.tsArena_.slab_TypedScope.Allocs; got != 0 {
		t.Fatalf("expected 0 tsArena_.slab_TypedScope allocs for bottomScope, got %d", got)
	}

	// 3. Task 3: GlobalNamespace_Name PooledSlab + nil-outer safety + ReleaseSlabs_().
	j2gort.ResetCounters()
	gn := NewGlobalNamespace()
	bgn := &GlobalNamespace_BuildGlobalNamespace{outer_: gn}
	leaf := bgn.GetOrCreateName("foo.bar.baz")
	if leaf == nil || leaf.baseName != "baz" || leaf.parent == nil || leaf.parent.baseName != "bar" ||
		leaf.parent.parent == nil || leaf.parent.parent.baseName != "foo" || leaf.outer_ != gn {
		t.Fatalf("unexpected GlobalNamespace_Name tree: %+v", leaf)
	}
	tn := gn.CreateNameForTesting("qux")
	if tn == nil || tn.baseName != "qux" || tn.outer_ != gn {
		t.Fatalf("unexpected CreateNameForTesting: %+v", tn)
	}
	if gn.slab_GlobalNamespace_Name.Allocs != 4 {
		t.Fatalf("expected 4 slab_GlobalNamespace_Name allocs (foo, bar, baz, qux), got %d", gn.slab_GlobalNamespace_Name.Allocs)
	}
	gn.ReleaseSlabs_()
	if gn.slab_GlobalNamespace_Ref.Releases != 1 || gn.slab_GlobalNamespace_Name.Releases != 1 {
		t.Fatalf("expected both Ref and Name slabs released, got ref=%d name=%d",
			gn.slab_GlobalNamespace_Ref.Releases, gn.slab_GlobalNamespace_Name.Releases)
	}
	// Nil-outer standalone Name.AddProperty must not panic.
	standalone := NewGlobalNamespace_Name(nil, "root", nil, gnSourceKind_CODE)
	child := standalone.AddProperty("child", gnSourceKind_EXTERN)
	if child == nil || child.baseName != "child" || child.parent != standalone || child.sourceKind != gnSourceKind_EXTERN {
		t.Fatalf("nil-outer AddProperty failed: %+v", child)
	}

	// 4. Task 4: CrossChunkCodeMotion_GlobalSymbolCollector.globalSymbolforVar -> OrderedMap[*Var, *CrossChunkCodeMotion_GlobalSymbol].
	j2gort.ResetCounters()
	gsc := NewCrossChunkCodeMotion_GlobalSymbolCollector(&CrossChunkCodeMotion{})
	v1, v2 := &Var{name: "a"}, &Var{name: "b"}
	s1 := gsc.GetGlobalSymbol(v1)
	s1Again := gsc.GetGlobalSymbol(v1)
	s2 := gsc.GetGlobalSymbol(v2)
	if s1 != s1Again || s1 == s2 || s1.var_ != v1 || s2.var_ != v2 {
		t.Fatalf("unexpected GlobalSymbol lookup results: s1=%p s1Again=%p s2=%p", s1, s1Again, s2)
	}
	if j2gort.CustomMapOps.Load() != 0 {
		t.Fatalf("expected 0 CustomOrderedMap ops after identity OrderedMap upgrade, got %d", j2gort.CustomMapOps.Load())
	}

	// 5. Task 5: RemoveUnusedCode_NameDeclarationStatement direct zeroed-slab init.
	pkg_rhino.RegisterNode(100, 101, true)
	pkg_rhino.RegisterNode(102, 103, true)
	pkg_rhino.RegisterNode(200, 201, false)
	ruc := &RemoveUnusedCode{}
	bld := &RemoveUnusedCode_RemovableBuilder{
		outer_:                          ruc,
		propertyName:                    "propA",
		continuation0:                   RemoveUnusedCode_Continuation{node: 11},
		continuation1:                   RemoveUnusedCode_Continuation{node: 12},
		restContinuations:               []RemoveUnusedCode_Continuation{{node: 13}},
		isPrototypeDotPropertyReference: true,
		isThisDotPropertyReference:      true,
	}
	nds1 := bld.BuildNameDeclarationStatement(100)
	if nds1.declarationStatement != 100 || nds1.targetNode != 101 || nds1.outer_ != ruc ||
		nds1.RemoveUnusedCode_RemovableBase.outer_ != ruc || nds1.self_ != nds1 ||
		nds1.propertyName != "propA" || nds1.continuation0.node != 11 || nds1.continuation1.node != 12 ||
		len(nds1.restContinuations) != 1 || nds1.restContinuations[0].node != 13 ||
		!nds1.isPrototypeDotPropertyReference || !nds1.isThisDotPropertyReference {
		t.Fatalf("unexpected nds1: %+v", nds1)
	}
	// Builder was reset by finish; test zero-continuation fast path via allocPooled helper.
	nds2 := allocPooled_RemoveUnusedCode_NameDeclarationStatement(ruc, bld, 102)
	if nds2.declarationStatement != 102 || nds2.targetNode != 103 ||
		nds2.propertyName != j2gort.NullString || nds2.continuation0 != (RemoveUnusedCode_Continuation{}) ||
		nds2.isPrototypeDotPropertyReference || nds2.isThisDotPropertyReference {
		t.Fatalf("unexpected nds2: %+v", nds2)
	}
	assertPanics(t, "non-name declaration CheckArgumentMsg", func() {
		_ = bld.BuildNameDeclarationStatement(200)
	})
	assertPanics(t, "GetOnlyChild panic on unregistered node", func() {
		_ = bld.BuildNameDeclarationStatement(999)
	})
}
`,
	}

	wave45bWriteFiles(t, tmpDir, jscompFiles)
	optDir := filepath.Join(tmpDir, "opt", "jscomp")
	wave45bRunAllAndCheckIdempotent(t, optDir)

	readFile := func(name string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(optDir, name))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	// Verify Task 1 AST (TypeInference embedded bottomTypedScope_ reuse)
	tiSrc := readFile("type_inference.go")
	for _, want := range []string{
		"bottomTypedScope_ TypedScope",
		"func reinitBottomTypedScope_(",
		"reinitBottomTypedScope_(&this_.bottomTypedScope_, root)",
		"this_.bottomTypedScope_ = TypedScope{}",
	} {
		if !strings.Contains(tiSrc, want) {
			t.Fatalf("type_inference.go missing %q:\n%s", want, tiSrc)
		}
	}
	if strings.Count(tiSrc, "reinitBottomTypedScope_(&this_.bottomTypedScope_, root)") != 2 {
		t.Fatalf("expected 2 reinitBottomTypedScope_ call sites (InitTypeInference + ResetForScope):\n%s", tiSrc)
	}

	// Verify Task 3 AST (GlobalNamespace_Name PooledSlab)
	gnSrc := readFile("global_namespace.go")
	for _, want := range []string{
		"slab_GlobalNamespace_Name j2gort.PooledSlab[GlobalNamespace_Name]",
		"func allocPooled_GlobalNamespace_Name(",
		"allocPooled_GlobalNamespace_Name(this_, name, nil, gnSourceKind_CODE)",
		"allocPooled_GlobalNamespace_Name(this_.outer_, name, nil, this_.outer_.sourceKind)",
		"val_10 := allocPooled_GlobalNamespace_Name(this_.outer_, name, this_, sourceKind)",
		"this_.slab_GlobalNamespace_Name.Release()",
	} {
		if !strings.Contains(gnSrc, want) {
			t.Fatalf("global_namespace.go missing %q:\n%s", want, gnSrc)
		}
	}

	// Verify Task 4 AST (CrossChunkCodeMotion_GlobalSymbolCollector.globalSymbolforVar -> OrderedMap)
	cccmSrc := readFile("cross_chunk_code_motion.go")
	for _, want := range []string{
		"j2gort.OrderedMap[*Var, *CrossChunkCodeMotion_GlobalSymbol]",
		"j2gort.NewOrderedMap[*Var, *CrossChunkCodeMotion_GlobalSymbol]()",
	} {
		if !strings.Contains(cccmSrc, want) {
			t.Fatalf("cross_chunk_code_motion.go missing %q:\n%s", want, cccmSrc)
		}
	}
	if strings.Contains(cccmSrc, "CustomOrderedMap") {
		t.Fatalf("cross_chunk_code_motion.go still contains CustomOrderedMap:\n%s", cccmSrc)
	}

	// Verify Task 5 AST (RemoveUnusedCode_NameDeclarationStatement direct zeroed-slab init)
	rucSrc := readFile("remove_unused_code.go")
	if strings.Count(rucSrc, "InitRemoveUnusedCode_NameDeclarationStatement(") != 1 {
		t.Fatalf("expected InitRemoveUnusedCode_NameDeclarationStatement to appear only in its own decl (0 calls left in BuildNameDeclarationStatement / allocPooled), got:\n%s", rucSrc)
	}
	if strings.Count(rucSrc, "targetNode_ := declarationStatement.GetOnlyChild()") != 2 {
		t.Fatalf("expected direct zeroed-slab init in both BuildNameDeclarationStatement and allocPooled_RemoveUnusedCode_NameDeclarationStatement:\n%s", rucSrc)
	}

	// Compile & run runtime parity test in temp module.
	cmd := exec.Command("go", "test", "-v", "-p", "1", "-parallel", "1", "./opt/jscomp")
	cmd.Dir = tmpDir
	cmd.Env = append(os.Environ(), "GOFLAGS=-p=2")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Wave 45B runtime test failed: %v\n%s", err, string(out))
	}
}

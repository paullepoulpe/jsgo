package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

const wave40bMockJ2gortSrc = `package j2gort

import (
	"fmt"
	"sync/atomic"
)

const NullString = "\xff\x00j2go_null"

var (
	PooledSlabAllocCount   atomic.Int64
	PooledSlabReleaseCount atomic.Int64
	OrderedMapReleaseCount atomic.Int64
)

func ResetCounters() {
	PooledSlabAllocCount.Store(0)
	PooledSlabReleaseCount.Store(0)
	OrderedMapReleaseCount.Store(0)
}

type ClinitState struct {
	done uint32
}

func EnterClinit(s *ClinitState, done *uint32) bool {
	return atomic.LoadUint32(done) == 0
}

func ExitClinit(s *ClinitState, done *uint32) {
	atomic.StoreUint32(done, 1)
}

func AssignIdentHash(p *atomic.Int32) int32 {
	p.Store(1)
	return 1
}

func CheckNotNull[T any](v T) T { return v }

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

func BoxInt(v int32) any      { return v }
func BoxString(s string) any  { return s }
func PtrToAny[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}
func ToAny[T any](v T) any { return v }
func StrVal(v any) string  { return fmt.Sprint(v) }
func StrValString(s string) string {
	if s == NullString {
		return "null"
	}
	return s
}
func Ext_java_lang_String_S_Format_2(tmpl string, args []any) string {
	return fmt.Sprintf(tmpl, args...)
}

type Ext_java_util_function_Supplier interface {
	IsJavaNil_() bool
	Get() any
}

type lambdaSupplier struct {
	fn func() any
}

func (l *lambdaSupplier) IsJavaNil_() bool { return l == nil }
func (l *lambdaSupplier) Get() any         { return l.fn() }

func NewLambda_Ext_java_util_function_Supplier(fn func() any) Ext_java_util_function_Supplier {
	return &lambdaSupplier{fn: fn}
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

func AcquireOrderedMapCap[K comparable, V any](cap int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{
		keys: make([]K, 0, cap),
		vals: make([]V, 0, cap),
	}
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

func (m *OrderedMap[K, V]) IsEmpty() bool {
	return len(m.keys) == 0
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

type MapValueCursor[V any] struct {
	vals []V
	idx  int
}

func (m *OrderedMap[K, V]) IterValues() MapValueCursor[V] {
	return MapValueCursor[V]{vals: m.vals}
}

func (c *MapValueCursor[V]) HasNext() bool { return c.idx < len(c.vals) }
func (c *MapValueCursor[V]) Next() V {
	v := c.vals[c.idx]
	c.idx++
	return v
}
`

func wave40bWriteFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func wave40bRunGoTest(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go test failed in %s: %v\n%s", dir, err, string(out))
	}
}

func TestWave40B_DiagnosticLogFileDeadAllocationGuard(t *testing.T) {
	tmpDir := t.TempDir()

	diagSrc := `package diagnostic

import "fmt"

type LogFile interface {
	IsJavaNil_() bool
	Implements_diagnostic_LogFile()
	IsLogging() bool
	Log_1_string(value string) LogFile
	Log_2_string_SliceOfObject(template string, values []any) LogFile
	Log_1_Ext_java_util_function_Supplier(supplier interface{ Get() any }) LogFile
	LogJson_1_string(json string) LogFile
}

type TestLogFile struct {
	Enabled bool
	Entries []string
}

func (l *TestLogFile) IsJavaNil_() bool               { return l == nil }
func (*TestLogFile) Implements_diagnostic_LogFile()   {}
func (l *TestLogFile) IsLogging() bool                { return l != nil && l.Enabled }
func (l *TestLogFile) Log_1_string(v string) LogFile {
	if l.IsLogging() {
		l.Entries = append(l.Entries, v)
	}
	return l
}
func (l *TestLogFile) Log_2_string_SliceOfObject(tmpl string, vals []any) LogFile {
	if l.IsLogging() {
		l.Entries = append(l.Entries, fmt.Sprintf(tmpl, vals...))
	}
	return l
}
func (l *TestLogFile) Log_1_Ext_java_util_function_Supplier(sup interface{ Get() any }) LogFile {
	if l.IsLogging() {
		l.Entries = append(l.Entries, fmt.Sprint(sup.Get()))
	}
	return l
}
func (l *TestLogFile) LogJson_1_string(js string) LogFile {
	if l.IsLogging() {
		l.Entries = append(l.Entries, "json:"+js)
	}
	return l
}
`

	graphSrc := `package graph

type GraphvizGraph interface {
	IsJavaNil_() bool
	Implements_graph_GraphvizGraph()
	GraphLabel() string
}

type DiGraph struct {
	Label string
}

func (g *DiGraph) IsJavaNil_() bool                { return g == nil }
func (*DiGraph) Implements_graph_GraphvizGraph()   {}
func (g *DiGraph) GraphLabel() string              { return g.Label }

func Upcast_GraphvizGraph(g *DiGraph) GraphvizGraph { return g }
`

	jscompSrc := `package jscomp

import (
	"sync/atomic"

	pkg_diagnostic "example.com/w40b_log/gen/diagnostic"
	pkg_graph "example.com/w40b_log/gen/graph"
	"example.com/w40b_log/gen/j2gort"
)

var (
	ToGraphvizCalls     atomic.Int64
	ToDotCalls          atomic.Int64
	SupplierBodyCalls   atomic.Int64
	ImpureRecvCalls     atomic.Int64
	ImpureArgSideEffect atomic.Int64
)

func ResetLogTestCounters() {
	ToGraphvizCalls.Store(0)
	ToDotCalls.Store(0)
	SupplierBodyCalls.Store(0)
	ImpureRecvCalls.Store(0)
	ImpureArgSideEffect.Store(0)
}

type ChunkGraph struct {
	Name string
}

func (cg *ChunkGraph) ToGraphvizGraph() *pkg_graph.DiGraph {
	ToGraphvizCalls.Add(1)
	return &pkg_graph.DiGraph{Label: cg.Name}
}

func DotFormatter_ToDot_1(g pkg_graph.GraphvizGraph) string {
	ToDotCalls.Add(1)
	return "dot:" + g.GraphLabel()
}

type StripCode_StripCodeConfigRecord struct {
	Tag string
}

func NewStripCode_StripCodeConfigRecord(tag string) *StripCode_StripCodeConfigRecord {
	return &StripCode_StripCodeConfigRecord{Tag: tag}
}

type PassBase struct {
	decisionsLog pkg_diagnostic.LogFile
}

func (b *PassBase) GetPassBase() *PassBase { return b }

type LogHarness struct {
	PassBase
	log        pkg_diagnostic.LogFile
	chunkGraph *ChunkGraph
}

func (this_ *LogHarness) createOrReopenLog() pkg_diagnostic.LogFile {
	ImpureRecvCalls.Add(1)
	return this_.log
}

func impureMutatingMessage() string {
	ImpureArgSideEffect.Add(1)
	return "impure-arg"
}

func (this_ *LogHarness) RunHotLogs(chunkGraphLog pkg_diagnostic.LogFile, colorGraph *pkg_graph.DiGraph, plainIdent string) {
	// 1. Compiler.initChunks shape (with Upcast_GraphvizGraph + ToGraphvizGraph + DotFormatter_ToDot_1).
	chunkGraphLog.Log_1_string(DotFormatter_ToDot_1(pkg_graph.Upcast_GraphvizGraph(this_.chunkGraph.ToGraphvizGraph())))

	// 2. AmbiguateProperties shape.
	chunkGraphLog.Log_1_string(DotFormatter_ToDot_1(pkg_graph.Upcast_GraphvizGraph(colorGraph)))

	// 3. ChunkGraphAwareLicenseTracker varargs []any slice shape.
	this_.log.Log_2_string_SliceOfObject("Chunk %s depends on %d", []any{"c0", j2gort.BoxInt(int32(7))})

	// 4. Supplier lambda shape (closure body has arbitrary logic, only NewLambda_* construction is at call site).
	this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		SupplierBodyCalls.Add(1)
		rec := NewStripCode_StripCodeConfigRecord("strip")
		return "decision:" + rec.Tag
	}))

	// 5. Base-getter receiver + LogJson_* + string concat shape.
	this_.GetPassBase().decisionsLog.LogJson_1_string("{" + plainIdent + "}")

	// --- Negative controls ---
	// 6a. 1-arg plain literal -> 0 allocs already, must NOT be wrapped.
	chunkGraphLog.Log_1_string("plain-literal")
	// 6b. 1-arg plain ident -> 0 allocs already, must NOT be wrapped.
	chunkGraphLog.Log_1_string(plainIdent)
	// 6c. Impure receiver call -> must NOT be wrapped (or side effect would be skipped / duplicated).
	this_.createOrReopenLog().Log_2_string_SliceOfObject("impure-recv %s", []any{plainIdent})
	// 6d. Impure argument call -> must NOT be wrapped.
	this_.log.Log_1_string(impureMutatingMessage())
	// 6e. Already guarded by IsLogging() -> must NOT be double-wrapped.
	if this_.log.IsLogging() {
		this_.log.Log_2_string_SliceOfObject("already %s", []any{plainIdent})
	}
}

func (this_ *LogHarness) RunPureAllocSitesOnly(chunkGraphLog pkg_diagnostic.LogFile, colorGraph *pkg_graph.DiGraph, plainIdent string) {
	chunkGraphLog.Log_1_string(DotFormatter_ToDot_1(pkg_graph.Upcast_GraphvizGraph(this_.chunkGraph.ToGraphvizGraph())))
	chunkGraphLog.Log_1_string(DotFormatter_ToDot_1(pkg_graph.Upcast_GraphvizGraph(colorGraph)))
	this_.log.Log_2_string_SliceOfObject("Chunk %s depends on %d", []any{"c0", j2gort.BoxInt(int32(7))})
	this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		SupplierBodyCalls.Add(1)
		return "decision:" + plainIdent
	}))
	this_.GetPassBase().decisionsLog.LogJson_1_string("{" + plainIdent + "}")
}
`

	parityTestSrc := `package parity_test

import (
	"reflect"
	"testing"

	pkg_diagnostic "example.com/w40b_log/gen/diagnostic"
	pkg_graph "example.com/w40b_log/gen/graph"
	optjs "example.com/w40b_log/opt/jscomp"
	origjs "example.com/w40b_log/orig/jscomp"
)

func TestLogGuardBehaviorAndZeroAllocs(t *testing.T) {
	// 1. Disabled logging: zero eager ToGraphviz / ToDot / Supplier calls, and 0 heap allocs on pure sites!
	optjs.ResetLogTestCounters()
	noOp1 := &pkg_diagnostic.TestLogFile{Enabled: false}
	noOp2 := &pkg_diagnostic.TestLogFile{Enabled: false}
	noOp3 := &pkg_diagnostic.TestLogFile{Enabled: false}
	cg := &pkg_graph.DiGraph{Label: "colors"}
	hOpt := &optjs.LogHarness{}
	hOpt.InitForTest("chunks", noOp1, noOp2)
	hOpt.RunHotLogs(noOp3, cg, "k")

	if got := optjs.ToGraphvizCalls.Load(); got != 0 {
		t.Fatalf("expected 0 ToGraphvizGraph calls when !IsLogging(), got %d", got)
	}
	if got := optjs.ToDotCalls.Load(); got != 0 {
		t.Fatalf("expected 0 DotFormatter_ToDot_1 calls when !IsLogging(), got %d", got)
	}
	if got := optjs.SupplierBodyCalls.Load(); got != 0 {
		t.Fatalf("expected 0 supplier body calls when !IsLogging(), got %d", got)
	}
	// Impure negative controls MUST still execute exactly once!
	if got := optjs.ImpureRecvCalls.Load(); got != 1 {
		t.Fatalf("expected ImpureRecvCalls == 1, got %d", got)
	}
	if got := optjs.ImpureArgSideEffect.Load(); got != 1 {
		t.Fatalf("expected ImpureArgSideEffect == 1, got %d", got)
	}

	allocs := testing.AllocsPerRun(100, func() {
		hOpt.RunPureAllocSitesOnly(noOp3, cg, "k")
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs/op on guarded pure LogFile sites when !IsLogging(), got %v", allocs)
	}

	// 2. Enabled logging: exact parity between orig and opt!
	origL1, origL2, origL3 := &pkg_diagnostic.TestLogFile{Enabled: true}, &pkg_diagnostic.TestLogFile{Enabled: true}, &pkg_diagnostic.TestLogFile{Enabled: true}
	hOrig := &origjs.LogHarness{}
	hOrig.InitForTest("chunks", origL1, origL2)
	hOrig.RunHotLogs(origL3, cg, "k")

	optL1, optL2, optL3 := &pkg_diagnostic.TestLogFile{Enabled: true}, &pkg_diagnostic.TestLogFile{Enabled: true}, &pkg_diagnostic.TestLogFile{Enabled: true}
	hOpt2 := &optjs.LogHarness{}
	hOpt2.InitForTest("chunks", optL1, optL2)
	hOpt2.RunHotLogs(optL3, cg, "k")

	if !reflect.DeepEqual(origL1.Entries, optL1.Entries) ||
		!reflect.DeepEqual(origL2.Entries, optL2.Entries) ||
		!reflect.DeepEqual(origL3.Entries, optL3.Entries) {
		t.Fatalf("enabled log output mismatch:\norig=(%v, %v, %v)\nopt=(%v, %v, %v)",
			origL1.Entries, origL2.Entries, origL3.Entries,
			optL1.Entries, optL2.Entries, optL3.Entries)
	}
}
`

	helperMethods := `
func (this_ *LogHarness) SetLogsForTest(log, dec pkg_diagnostic.LogFile) {
	this_.log = log
	this_.decisionsLog = dec
}

func (this_ *LogHarness) InitForTest(cgName string, log, dec pkg_diagnostic.LogFile) {
	this_.chunkGraph = &ChunkGraph{Name: cgName}
	this_.log = log
	this_.decisionsLog = dec
}
`

	wave40bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                 "module example.com/w40b_log\n\ngo 1.23\n",
		"gen/j2gort/j2gort.go":   wave40bMockJ2gortSrc,
		"gen/diagnostic/log.go":  diagSrc,
		"gen/graph/graph.go":     graphSrc,
		"orig/jscomp/harness.go": jscompSrc + helperMethods,
		"opt/jscomp/harness.go":  jscompSrc + helperMethods,
		"parity/parity_test.go":  parityTestSrc,
	})

	optDir := filepath.Join(tmpDir, "opt/jscomp")
	s1, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.ByPass["j2gort_peephole"] < 10 {
		t.Fatalf("expected >= 10 j2gort_peephole sites (10 pure LogFile calls across 2 methods), got %+v", s1.ByPass)
	}

	// Run 2 MUST be 100% idempotent!
	s2, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d ByPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	raw, err := os.ReadFile(filepath.Join(optDir, "harness.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	for _, want := range []string{
		"if chunkGraphLog.IsLogging() {",
		"if this_.log.IsLogging() {",
		"if this_.decisionsLog.IsLogging() {",
		"if this_.GetPassBase().decisionsLog.IsLogging() {",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected transformed harness.go to contain %q, got:\n%s", want, got)
		}
	}
	if strings.Contains(got, "this_.createOrReopenLog().IsLogging()") {
		t.Fatalf("impure receiver createOrReopenLog() MUST NOT be duplicated in IsLogging() guard:\n%s", got)
	}

	wave40bRunGoTest(t, tmpDir)
}

func TestWave40B_RhinoJSTypeRegistrySlabsAndJscompClearRelease(t *testing.T) {
	tmpDir := t.TempDir()

	rhinoSrc := `package rhino

import (
	"sync/atomic"

	"example.com/w40b_rhino/gen/j2gort"
)

var (
	UnionClinitRuns    atomic.Int32
	FunctionClinitRuns atomic.Int32
	UnionSelfBound     atomic.Int32
	FunctionSelfBound  atomic.Int32
)

type JSType interface {
	IsJavaNil_() bool
	Implements_rhino_JSType()
	GetJSTypeBase() *JSTypeBase
}

type ObjectType interface {
	JSType
	Implements_rhino_ObjectType()
}

type JSTypeBase struct {
	self_    any
	registry *JSTypeRegistry
}

func (b *JSTypeBase) BindSelf_(s any)          { b.self_ = s }
func (b *JSTypeBase) GetJSTypeBase() *JSTypeBase { return b }

type ProxyObjectTypeBase struct {
	JSTypeBase
}

type JSTypeRegistry struct {
	tag int32
}

func NewJSTypeRegistry(tag int32) *JSTypeRegistry {
	return &JSTypeRegistry{tag: tag}
}

var clinitDone_UnionType uint32
var clinitState_UnionType j2gort.ClinitState

func clinit_UnionType() {
	if atomic.LoadUint32(&clinitDone_UnionType) != 0 {
		return
	}
	if !j2gort.EnterClinit(&clinitState_UnionType, &clinitDone_UnionType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_UnionType, &clinitDone_UnionType)
	UnionClinitRuns.Add(1)
}

type UnionType struct {
	JSTypeBase
	altCount int32
}

func (this_ *UnionType) BindSelf_(s any) {
	this_.JSTypeBase.BindSelf_(s)
	UnionSelfBound.Add(1)
}
func (this_ *UnionType) IsJavaNil_() bool         { return this_ == nil }
func (*UnionType) Implements_rhino_JSType()       {}
func (this_ *UnionType) AltCount() int32          { return this_.altCount }
func (this_ *UnionType) SelfMatches() bool        { return this_.self_ == this_ }

func NewUnionType(builder *UnionType_Builder) *UnionType {
	clinit_UnionType()
	this_ := &UnionType{}
	this_.BindSelf_(this_)
	InitUnionType(this_, builder)
	return this_
}

func InitUnionType(this_ *UnionType, builder *UnionType_Builder) {
	this_.registry = builder.registry
	this_.altCount = builder.size
}

type UnionType_Builder struct {
	registry *JSTypeRegistry
	size     int32
	fast     bool
}

func NewUnionType_Builder(reg *JSTypeRegistry, size int32, fast bool) *UnionType_Builder {
	return &UnionType_Builder{registry: reg, size: size, fast: fast}
}

func (this_ *UnionType_Builder) Build() JSType {
	if this_.fast {
		return NewUnionType(this_)
	}
	return NewUnionType(this_)
}

var clinitDone_FunctionType uint32
var clinitState_FunctionType j2gort.ClinitState

func clinit_FunctionType() {
	if atomic.LoadUint32(&clinitDone_FunctionType) != 0 {
		return
	}
	if !j2gort.EnterClinit(&clinitState_FunctionType, &clinitDone_FunctionType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FunctionType, &clinitDone_FunctionType)
	FunctionClinitRuns.Add(1)
}

type PrototypeObjectType_BuilderBase struct {
	registry *JSTypeRegistry
}

type FunctionType_Builder struct {
	PrototypeObjectType_BuilderBase
	arity int32
}

func NewFunctionType_Builder(reg *JSTypeRegistry, arity int32) *FunctionType_Builder {
	b := &FunctionType_Builder{arity: arity}
	b.registry = reg
	return b
}

type FunctionType struct {
	JSTypeBase
	arity int32
	name  string
}

func (this_ *FunctionType) BindSelf_(s any) {
	this_.JSTypeBase.BindSelf_(s)
	FunctionSelfBound.Add(1)
}
func (this_ *FunctionType) IsJavaNil_() bool       { return this_ == nil }
func (*FunctionType) Implements_rhino_JSType()     {}
func (*FunctionType) Implements_rhino_ObjectType() {}
func (this_ *FunctionType) SelfMatches() bool      { return this_.self_ == this_ }
func (this_ *FunctionType) Name() string           { return this_.name }

func NewFunctionType_1(builder *FunctionType_Builder) *FunctionType {
	clinit_FunctionType()
	this_ := &FunctionType{}
	this_.BindSelf_(this_)
	InitFunctionType_1(this_, builder)
	return this_
}

func InitFunctionType_1(this_ *FunctionType, builder *FunctionType_Builder) {
	this_.registry = builder.registry
	this_.arity = builder.arity
	this_.name = "builder-fn"
}

func NewFunctionType_2(registry *JSTypeRegistry, name string, arity int32) *FunctionType {
	clinit_FunctionType()
	this_ := &FunctionType{}
	this_.BindSelf_(this_)
	InitFunctionType_2(this_, registry, name, arity)
	return this_
}

func InitFunctionType_2(this_ *FunctionType, registry *JSTypeRegistry, name string, arity int32) {
	this_.registry = registry
	this_.name = name
	this_.arity = arity
}

func (this_ *FunctionType_Builder) Build() *FunctionType {
	return NewFunctionType_1(this_)
}

func (this_ *JSTypeRegistry) CreateAnonymousFunctionType(arity int32) *FunctionType {
	return NewFunctionType_2(this_, "anon", arity)
}

type TemplatizedType struct {
	ProxyObjectTypeBase
	baseType ObjectType
	paramCnt int32
}

func (this_ *TemplatizedType) BindSelf_(s any)          { this_.ProxyObjectTypeBase.BindSelf_(s) }
func (this_ *TemplatizedType) IsJavaNil_() bool         { return this_ == nil }
func (*TemplatizedType) Implements_rhino_JSType()       {}
func (*TemplatizedType) Implements_rhino_ObjectType()   {}
func (this_ *TemplatizedType) ParamCount() int32        { return this_.paramCnt }

func NewTemplatizedType(registry *JSTypeRegistry, baseType ObjectType, paramCnt int32) *TemplatizedType {
	this_ := &TemplatizedType{}
	this_.BindSelf_(this_)
	InitTemplatizedType(this_, registry, baseType, paramCnt)
	return this_
}

func InitTemplatizedType(this_ *TemplatizedType, registry *JSTypeRegistry, baseType ObjectType, paramCnt int32) {
	this_.registry = registry
	this_.baseType = baseType
	this_.paramCnt = paramCnt
}

func (this_ *JSTypeRegistry) CreateTemplatizedType(baseType ObjectType, paramCnt int32) *TemplatizedType {
	return NewTemplatizedType(this_, baseType, paramCnt)
}

func (this_ *TemplatizedType) ResolveInternal() *TemplatizedType {
	return NewTemplatizedType(this_.GetJSTypeBase().registry, this_.baseType, this_.paramCnt+1)
}
`

	jscompCompSrc := `package jscomp

import (
	"example.com/w40b_rhino/gen/j2gort"
	pkg_rhino "example.com/w40b_rhino/opt/rhino"
)

type Compiler struct {
	typeRegistry *pkg_rhino.JSTypeRegistry
	name         string
}

func NewCompiler(reg *pkg_rhino.JSTypeRegistry) *Compiler {
	return &Compiler{typeRegistry: reg, name: j2gort.CheckNotNullString("comp")}
}

func (this_ *Compiler) ClearJSTypeRegistry() {
	this_.typeRegistry = nil
}
`

	runtimeTestSrc := `package rhino_slab_test

import (
	"testing"

	"example.com/w40b_rhino/gen/j2gort"
	"example.com/w40b_rhino/opt/jscomp"
	"example.com/w40b_rhino/opt/rhino"
)

func TestJSTypeRegistrySlabsAndClearRelease(t *testing.T) {
	j2gort.ResetCounters()
	reg := rhino.NewJSTypeRegistry(42)

	u1 := rhino.NewUnionType_Builder(reg, 2, true).Build().(*rhino.UnionType)
	u2 := rhino.NewUnionType_Builder(reg, 3, false).Build().(*rhino.UnionType)
	f1 := rhino.NewFunctionType_Builder(reg, 1).Build()
	f2 := reg.CreateAnonymousFunctionType(4)
	t1 := reg.CreateTemplatizedType(f1, 1)
	t2 := t1.ResolveInternal()

	if !u1.SelfMatches() || !u2.SelfMatches() || !f1.SelfMatches() || !f2.SelfMatches() {
		t.Fatal("BindSelf_ was not preserved on pooled UnionType / FunctionType")
	}
	if rhino.UnionClinitRuns.Load() != 1 || rhino.FunctionClinitRuns.Load() != 1 {
		t.Fatalf("clinit not preserved: union=%d fn=%d", rhino.UnionClinitRuns.Load(), rhino.FunctionClinitRuns.Load())
	}
	if t2.ParamCount() != 2 {
		t.Fatalf("expected resolved TemplatizedType paramCnt=2, got %d", t2.ParamCount())
	}
	if got := j2gort.PooledSlabAllocCount.Load(); got != 6 {
		t.Fatalf("expected all 6 allocations (2 UnionType + 2 FunctionType + 2 TemplatizedType) to hit PooledSlab, got %d", got)
	}
	if got := j2gort.PooledSlabReleaseCount.Load(); got != 0 {
		t.Fatalf("expected 0 slab releases before ClearJSTypeRegistry, got %d", got)
	}

	comp := jscomp.NewCompiler(reg)
	comp.ClearJSTypeRegistry()
	if got := j2gort.PooledSlabReleaseCount.Load(); got != 3 {
		t.Fatalf("expected ClearJSTypeRegistry to release all 3 JSTypeRegistry slabs, got %d", got)
	}
	// Calling ClearJSTypeRegistry again when typeRegistry == nil must be a safe no-op.
	comp.ClearJSTypeRegistry()
	if got := j2gort.PooledSlabReleaseCount.Load(); got != 3 {
		t.Fatalf("expected second ClearJSTypeRegistry to be a no-op, got %d", got)
	}
}
`

	wave40bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                   "module example.com/w40b_rhino\n\ngo 1.23\n",
		"gen/j2gort/j2gort.go":     wave40bMockJ2gortSrc,
		"opt/rhino/jstype.go":      rhinoSrc,
		"opt/jscomp/compiler.go":   jscompCompSrc,
		"testpkg/registry_test.go": runtimeTestSrc,
	})

	rhinoDir := filepath.Join(tmpDir, "opt/rhino")
	jscompDir := filepath.Join(tmpDir, "opt/jscomp")

	// Run 1 on rhino, then jscomp (so jscomp sees JSTypeRegistry.ReleaseSlabs_()).
	r1, err := passes.RunOnPackageDir(rhinoDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("rhino Run 1 failed: %v", err)
	}
	if r1.TotalSites == 0 {
		t.Fatalf("expected non-zero sites on rhino Run 1, got %+v", r1)
	}
	j1, err := passes.RunOnPackageDir(jscompDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("jscomp Run 1 failed: %v", err)
	}
	if j1.TotalSites == 0 {
		t.Fatalf("expected non-zero sites on jscomp Run 1, got %+v", j1)
	}

	// Run 2 on BOTH packages MUST produce 0 sites!
	r2, err := passes.RunOnPackageDir(rhinoDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("rhino Run 2 failed: %v", err)
	}
	if r2.TotalSites != 0 || r2.FilesModified != 0 {
		t.Fatalf("rhino Run 2 NOT IDEMPOTENT: %+v", r2)
	}
	j2, err := passes.RunOnPackageDir(jscompDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("jscomp Run 2 failed: %v", err)
	}
	if j2.TotalSites != 0 || j2.FilesModified != 0 {
		t.Fatalf("jscomp Run 2 NOT IDEMPOTENT: %+v", j2)
	}

	rhinoOutBytes, err := os.ReadFile(filepath.Join(rhinoDir, "jstype.go"))
	if err != nil {
		t.Fatal(err)
	}
	rhinoOut := string(rhinoOutBytes)
	for _, want := range []string{
		"slab_UnionType",
		"slab_FunctionType",
		"slab_TemplatizedType",
		"allocPooled_UnionType(&this_.registry.slab_UnionType, this_)",
		"allocPooled_FunctionType_1(&this_.registry.slab_FunctionType, this_)",
		"allocPooled_FunctionType_2(&this_.slab_FunctionType, this_",
		"allocPooled_TemplatizedType(&this_.slab_TemplatizedType, this_",
		"allocPooled_TemplatizedType(&this_.registry.slab_TemplatizedType,",
		"func (this_ *JSTypeRegistry) ReleaseSlabs_()",
	} {
		if !strings.Contains(rhinoOut, want) {
			t.Fatalf("expected rhino/jstype.go to contain %q, got:\n%s", want, rhinoOut)
		}
	}

	jscompOutBytes, err := os.ReadFile(filepath.Join(jscompDir, "compiler.go"))
	if err != nil {
		t.Fatal(err)
	}
	jscompOut := string(jscompOutBytes)
	if !strings.Contains(jscompOut, "this_.typeRegistry.ReleaseSlabs_()") {
		t.Fatalf("expected jscomp/compiler.go to call this_.typeRegistry.ReleaseSlabs_(), got:\n%s", jscompOut)
	}

	wave40bRunGoTest(t, tmpDir)
}

func TestWave40B_TypedScopeAndTypedVarArenaLifetimeAndOverflowRelease(t *testing.T) {
	tmpDir := t.TempDir()

	rhinoNodeSrc := `package rhino

type Node uint32
`

	typedScopeSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w40b_tscope/gen/j2gort"
	pkg_rhino "example.com/w40b_tscope/gen/rhino"
)

var OverflowResetCount atomic.Int64

type TypedScope_Overflow struct {
	vars *j2gort.OrderedMap[string, *TypedVar]
}

func NewTypedScope_Overflow() *TypedScope_Overflow {
	return &TypedScope_Overflow{
		vars: j2gort.AcquireOrderedMapPtr[string, *TypedVar](),
	}
}

func (o *TypedScope_Overflow) ResetForPool() {
	OverflowResetCount.Add(1)
}

type TypedVar struct {
	name  string
	node  pkg_rhino.Node
	scope *TypedScope
}

func (this_ *TypedVar) BindSelf_(s any) { _ = s }
func (this_ *TypedVar) Name() string    { return this_.name }

func NewTypedVar(name string, node pkg_rhino.Node, scope *TypedScope) *TypedVar {
	this_ := &TypedVar{}
	this_.BindSelf_(this_)
	InitTypedVar(this_, name, node, scope)
	return this_
}

func InitTypedVar(this_ *TypedVar, name string, node pkg_rhino.Node, scope *TypedScope) {
	this_.name = name
	this_.node = node
	this_.scope = scope
}

type TypedScope struct {
	parent   *TypedScope
	rootNode pkg_rhino.Node
	var0     *TypedVar
	overflow *TypedScope_Overflow
	depth    int32
	isBottom bool
}

func (this_ *TypedScope) BindSelf_(s any)         { _ = s }
func (this_ *TypedScope) GetParent() *TypedScope  { return this_.parent }
func (this_ *TypedScope) IsBottom() bool          { return this_.isBottom }

func NewTypedScope_5(parent *TypedScope, rootNode pkg_rhino.Node, useOverflow bool) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_5(this_, parent, rootNode, useOverflow)
	return this_
}

func InitTypedScope_5(this_ *TypedScope, parent *TypedScope, rootNode pkg_rhino.Node, useOverflow bool) {
	this_.parent = parent
	this_.rootNode = rootNode
	this_.depth = parent.depth + 1
	if useOverflow {
		this_.overflow = NewTypedScope_Overflow()
	}
}

func NewTypedScope_2_Node_bool(rootNode pkg_rhino.Node, isBottom bool) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_2_Node_bool(this_, rootNode, isBottom)
	return this_
}

func InitTypedScope_2_Node_bool(this_ *TypedScope, rootNode pkg_rhino.Node, isBottom bool) {
	this_.rootNode = rootNode
	this_.isBottom = isBottom
}

func TypedScope_CreateGlobalScope(rootNode pkg_rhino.Node) *TypedScope {
	return NewTypedScope_2_Node_bool(rootNode, false)
}

func TypedScope_CreateLatticeBottom(rootNode pkg_rhino.Node) *TypedScope {
	return NewTypedScope_2_Node_bool(rootNode, true)
}

func (this_ *TypedScope) Declare(name string, n pkg_rhino.Node) *TypedVar {
	var v *TypedVar = NewTypedVar(name, n, this_)
	if this_.var0 == nil {
		this_.var0 = v
	} else {
		if this_.overflow == nil {
			this_.overflow = NewTypedScope_Overflow()
		}
		this_.overflow.vars.Put(name, v)
	}
	return v
}

func (this_ *TypedScope) MakeImplicitVar(name string) *TypedVar {
	return NewTypedVar(name, pkg_rhino.Node(1), this_)
}

func (this_ *TypedScope) LookupOwn(name string) *TypedVar {
	if this_.var0 != nil && this_.var0.name == name {
		return this_.var0
	}
	if this_.overflow != nil && this_.overflow.vars != nil {
		return this_.overflow.vars.Get(name)
	}
	return nil
}
`

	typedScopeCreatorSrc := `package jscomp

import pkg_rhino "example.com/w40b_tscope/gen/rhino"

type TypedScopeCreator struct {
	compiler *Compiler
}

func NewTypedScopeCreator_1(compiler *Compiler) *TypedScopeCreator {
	this_ := &TypedScopeCreator{}
	InitTypedScopeCreator_1(this_, compiler)
	return this_
}

func InitTypedScopeCreator_1(this_ *TypedScopeCreator, compiler *Compiler) {
	_ = compiler
	this_.compiler = compiler
}

func (this_ *TypedScopeCreator) CreateInitialScope(root pkg_rhino.Node) *TypedScope {
	var s *TypedScope = TypedScope_CreateGlobalScope(root)
	_ = s
	s.Declare("Object", root)
	s.Declare("Array", root)
	return s
}

func (this_ *TypedScopeCreator) priv_TypedScopeCreator_createScopeInternal_1(n pkg_rhino.Node, parent *TypedScope) *TypedScope {
	return NewTypedScope_5(parent, n, true)
}

type TypeInference struct {
	scopeCreator *TypedScopeCreator
	bottomScope  *TypedScope
}

func NewTypeInference(sc *TypedScopeCreator, fnRoot pkg_rhino.Node) *TypeInference {
	this_ := &TypeInference{}
	InitTypeInference(this_, sc, fnRoot)
	return this_
}

func InitTypeInference(this_ *TypeInference, sc *TypedScopeCreator, fnRoot pkg_rhino.Node) {
	this_.scopeCreator = sc
	this_.bottomScope = TypedScope_CreateLatticeBottom(fnRoot)
}

func (this_ *TypeInference) ResetForScope(fnRoot pkg_rhino.Node) {
	this_.bottomScope = TypedScope_CreateLatticeBottom(fnRoot)
}
`

	compilerSrc := `package jscomp

type Compiler struct {
	typedScopeCreator *TypedScopeCreator
	topScope          *TypedScope
}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (this_ *Compiler) GetTypedScopeCreator() *TypedScopeCreator {
	if this_.typedScopeCreator == nil {
		this_.typedScopeCreator = NewTypedScopeCreator_1(this_)
	}
	return this_.typedScopeCreator
}

func (this_ *Compiler) ClearTypedScopeCreator() {
	this_.typedScopeCreator = nil
}

func (this_ *Compiler) SetTopScope(x *TypedScope) {
	_ = x
	this_.topScope = x
}

func (this_ *Compiler) GetTopScope() *TypedScope {
	return this_.topScope
}
`

	lifetimeTestSrc := `package tscope_lifetime_test

import (
	"testing"

	"example.com/w40b_tscope/gen/j2gort"
	pkg_rhino "example.com/w40b_tscope/gen/rhino"
	"example.com/w40b_tscope/opt/jscomp"
)

func TestTypedScopeArenaDeferredReleaseLifetimes(t *testing.T) {
	// Scenario A: DefaultPassConfig order:
	//   GetTypedScopeCreator() -> SetTopScope(top) -> ClearTypedScopeCreator() -> post-check Declare on topScope -> SetTopScope(nil)
	j2gort.ResetCounters()
	jscomp.OverflowResetCount.Store(0)

	comp := jscomp.NewCompiler()
	tsc := comp.GetTypedScopeCreator()
	top := tsc.CreateInitialScope(pkg_rhino.Node(1))
	child := tsc.CreateChildForTest(pkg_rhino.Node(2), top)
	child.Declare("localA", pkg_rhino.Node(3))
	child.Declare("localB", pkg_rhino.Node(4))
	_ = child.MakeImplicitVar("arguments")
	grandChild := jscomp.NewTypedScope_5(child, pkg_rhino.Node(5), false)
	grandChild.Declare("nested", pkg_rhino.Node(6))

	ti := jscomp.NewTypeInference(tsc, pkg_rhino.Node(2))
	ti.ResetForScope(pkg_rhino.Node(5))

	comp.SetTopScope(top)
	allocsBeforeClear := j2gort.PooledSlabAllocCount.Load()
	if allocsBeforeClear < 10 {
		t.Fatalf("expected >= 10 PooledSlab allocations for TypedScope + TypedVar, got %d", allocsBeforeClear)
	}

	// Step 1: ClearTypedScopeCreator runs BEFORE ClearTopScope! Arena MUST NOT release yet!
	comp.ClearTypedScopeCreator()
	if got := j2gort.PooledSlabReleaseCount.Load(); got != 0 {
		t.Fatalf("PREMATURE RELEASE on ClearTypedScopeCreator while topScope != nil: releases=%d", got)
	}
	// Es6RewriteModules / ClosureRewriteModule still declares global module vars on topScope!
	modVar := comp.GetTopScope().Declare("module$exports$foo", pkg_rhino.Node(7))
	if modVar.Name() != "module$exports$foo" || comp.GetTopScope().LookupOwn("Object") == nil || child.LookupOwn("localB") == nil {
		t.Fatal("topScope / child TypedVars corrupted after ClearTypedScopeCreator()")
	}
	if j2gort.PooledSlabAllocCount.Load() != allocsBeforeClear+1 {
		t.Fatal("expected post-ClearTypedScopeCreator topScope.Declare to still allocate from tsArena_.slab_TypedVar")
	}

	// Step 2: SetTopScope(nil) now clears topScope -> BOTH are nil -> arena MUST release now!
	comp.SetTopScope(nil)
	if got := j2gort.PooledSlabReleaseCount.Load(); got != 2 {
		t.Fatalf("expected both slab_TypedScope and slab_TypedVar (2 slabs) released after SetTopScope(nil), got %d", got)
	}
	if got := j2gort.OrderedMapReleaseCount.Load(); got < 2 {
		t.Fatalf("expected TypedScope_Overflow.vars OrderedMaps released, got %d", got)
	}
	if got := jscomp.OverflowResetCount.Load(); got < 2 {
		t.Fatalf("expected TypedScope_Overflow.ResetForPool called during slab cleanup, got %d", got)
	}

	// Scenario B: Reverse order (SetTopScope(nil) BEFORE ClearTypedScopeCreator()).
	j2gort.ResetCounters()
	comp2 := jscomp.NewCompiler()
	tsc2 := comp2.GetTypedScopeCreator()
	top2 := tsc2.CreateInitialScope(pkg_rhino.Node(10))
	comp2.SetTopScope(top2)
	comp2.SetTopScope(nil)
	if got := j2gort.PooledSlabReleaseCount.Load(); got != 0 {
		t.Fatalf("PREMATURE RELEASE on SetTopScope(nil) while typedScopeCreator != nil: releases=%d", got)
	}
	comp2.ClearTypedScopeCreator()
	if got := j2gort.PooledSlabReleaseCount.Load(); got != 2 {
		t.Fatalf("expected 2 slabs released once ClearTypedScopeCreator() runs, got %d", got)
	}

	// Scenario C: Standalone NewTypedScopeCreator_1(compiler) (TypeCheck.processForTesting).
	j2gort.ResetCounters()
	comp3 := jscomp.NewCompiler()
	standalone := jscomp.NewTypedScopeCreator_1(comp3)
	top3 := standalone.CreateInitialScope(pkg_rhino.Node(20))
	comp3.SetTopScope(top3)
	comp3.SetTopScope(nil)
	if got := j2gort.PooledSlabReleaseCount.Load(); got != 0 {
		t.Fatalf("standalone TypedScopeCreator arena MUST NOT be released by compiler.SetTopScope(nil): releases=%d", got)
	}
	if top3.LookupOwn("Object") == nil {
		t.Fatal("standalone TypedScopeCreator top3 corrupted")
	}
}
`

	tscTestExport := `
func (this_ *TypedScopeCreator) CreateChildForTest(n pkg_rhino.Node, parent *TypedScope) *TypedScope {
	return this_.priv_TypedScopeCreator_createScopeInternal_1(n, parent)
}
`

	wave40bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                             "module example.com/w40b_tscope\n\ngo 1.23\n",
		"gen/j2gort/j2gort.go":               wave40bMockJ2gortSrc,
		"gen/rhino/node.go":                  rhinoNodeSrc,
		"opt/jscomp/typed_scope.go":          typedScopeSrc,
		"opt/jscomp/typed_scope_creator.go":  typedScopeCreatorSrc + tscTestExport,
		"opt/jscomp/compiler.go":             compilerSrc,
		"testpkg/tscope_lifetime_test.go":    lifetimeTestSrc,
	})

	optDir := filepath.Join(tmpDir, "opt/jscomp")
	s1, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero sites on Run 1, got %+v", s1)
	}

	// Run 2 MUST be 100% idempotent!
	s2, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d ByPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	compBytes, err := os.ReadFile(filepath.Join(optDir, "compiler.go"))
	if err != nil {
		t.Fatal(err)
	}
	compOut := string(compBytes)
	for _, want := range []string{
		"tsArena_",
		"*typedScopeArena_",
		"this_.tsArena_ = this_.typedScopeCreator.tsArena_",
		"func (this_ *Compiler) maybeReleaseTypedScopeArena_()",
		"this_.typedScopeCreator == nil && this_.topScope == nil",
	} {
		if !strings.Contains(compOut, want) {
			t.Fatalf("expected compiler.go to contain %q, got:\n%s", want, compOut)
		}
	}

	wave40bRunGoTest(t, tmpDir)
}

func TestWave40B_GlobalNamespaceSyntacticScopeCreatorRecycling(t *testing.T) {
	tmpDir := t.TempDir()

	rhinoSrc := `package rhino

type Node uint32

type StaticScope interface {
	IsJavaNil_() bool
	Implements_rhino_StaticScope()
}
`

	jscompScopeAndNtSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w40b_gn/gen/j2gort"
	pkg_rhino "example.com/w40b_gn/gen/rhino"
)

var ScannerAllocCount atomic.Int64

type AbstractCompiler interface {
	IsJavaNil_() bool
	Implements_jscomp_AbstractCompiler()
}

type CompilerPass interface {
	IsJavaNil_() bool
	Implements_jscomp_CompilerPass()
	Process_2(externs pkg_rhino.Node, root pkg_rhino.Node)
}

type AbstractScope interface {
	IsJavaNil_() bool
	Implements_jscomp_AbstractScope()
}

type Scope struct {
	rootNode pkg_rhino.Node
	parent   *Scope
	vars     *j2gort.OrderedMap[string, *Var]
}

func (s *Scope) BindSelf_(v any)             { _ = v }
func (s *Scope) IsJavaNil_() bool            { return s == nil }
func (*Scope) Implements_jscomp_AbstractScope() {}
func (s *Scope) Root() pkg_rhino.Node        { return s.rootNode }

type Var struct {
	name  string
	scope *Scope
}

func (v *Var) BindSelf_(x any)  { _ = x }
func (v *Var) Name() string     { return v.name }

type ScopeCreator interface {
	IsJavaNil_() bool
	Implements_jscomp_ScopeCreator()
	CreateScope(n pkg_rhino.Node, parent *Scope) *Scope
}

type SyntacticScopeCreator_ScopeScanner struct {
	compiler AbstractCompiler
}

func NewSyntacticScopeCreator_ScopeScanner(c AbstractCompiler) *SyntacticScopeCreator_ScopeScanner {
	ScannerAllocCount.Add(1)
	return &SyntacticScopeCreator_ScopeScanner{compiler: c}
}

type SyntacticScopeCreator struct {
	compiler        AbstractCompiler
	reusableScanner *SyntacticScopeCreator_ScopeScanner
	slab_Scope      j2gort.PooledSlab[Scope]
	slab_Var        j2gort.PooledSlab[Var]
}

func (s *SyntacticScopeCreator) IsJavaNil_() bool             { return s == nil }
func (*SyntacticScopeCreator) Implements_jscomp_ScopeCreator() {}

func NewSyntacticScopeCreator_1(compiler AbstractCompiler) *SyntacticScopeCreator {
	s := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_1(s, compiler)
	return s
}

func InitSyntacticScopeCreator_1(this_ *SyntacticScopeCreator, compiler AbstractCompiler) {
	this_.compiler = compiler
	this_.reusableScanner = NewSyntacticScopeCreator_ScopeScanner(compiler)
}

func (this_ *SyntacticScopeCreator) CreateScope(n pkg_rhino.Node, parent *Scope) *Scope {
	s := this_.slab_Scope.Alloc()
	s.rootNode = n
	s.parent = parent
	s.vars = j2gort.AcquireOrderedMapPtr[string, *Var]()
	v := this_.slab_Var.Alloc()
	v.name = "globalSym"
	v.scope = s
	s.vars.Put(v.name, v)
	return s
}

func (this_ *SyntacticScopeCreator) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_Scope.ReleaseWith(func(s_ *Scope) {
		if s_.vars != nil {
			j2gort.ReleaseOrderedMap(s_.vars)
			s_.vars = nil
		}
		*s_ = Scope{}
	})
	this_.slab_Var.ReleaseWith(func(v_ *Var) { *v_ = Var{} })
}

type NodeTraversal_Callback interface {
	IsJavaNil_() bool
	Implements_jscomp_NodeTraversal_Callback()
	ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool
	Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node)
}

type NodeTraversal struct {
	compiler     AbstractCompiler
	callback     NodeTraversal_Callback
	scopeCreator ScopeCreator
	curScope     *Scope
}

func (t *NodeTraversal) GetScope() *Scope {
	if t.curScope == nil {
		sc := t.scopeCreator
		if sc == nil {
			sc = NewSyntacticScopeCreator_1(t.compiler)
		}
		t.curScope = sc.CreateScope(pkg_rhino.Node(1), nil)
	}
	return t.curScope
}

type NodeTraversal_Builder struct {
	compiler     AbstractCompiler
	callback     NodeTraversal_Callback
	scopeCreator ScopeCreator
}

func NodeTraversal_S_Builder() *NodeTraversal_Builder {
	return &NodeTraversal_Builder{}
}

func InitNodeTraversal_Builder(this_ *NodeTraversal_Builder) {
	*this_ = NodeTraversal_Builder{}
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

func (b *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	t := &NodeTraversal{compiler: b.compiler, callback: b.callback, scopeCreator: b.scopeCreator}
	if b.callback.ShouldTraverse_3(t, root, 0) {
		b.callback.Visit_3(t, root, 0)
	}
}

func (b *NodeTraversal_Builder) TraverseRoots(externs pkg_rhino.Node, root pkg_rhino.Node) {
	b.Traverse(externs)
	b.Traverse(root)
}
`

	globalNamespaceSrc := `package jscomp

import (
	"example.com/w40b_gn/gen/j2gort"
	pkg_rhino "example.com/w40b_gn/gen/rhino"
)

type GlobalNamespace_Ref struct {
	scope *Scope
	node  pkg_rhino.Node
}

func NewGlobalNamespace_Ref(scope *Scope, node pkg_rhino.Node) *GlobalNamespace_Ref {
	r := &GlobalNamespace_Ref{}
	InitGlobalNamespace_Ref(r, scope, node)
	return r
}

func InitGlobalNamespace_Ref(this_ *GlobalNamespace_Ref, scope *Scope, node pkg_rhino.Node) {
	this_.scope = scope
	this_.node = node
}

type GlobalNamespace_Name struct {
	name           string
	declaration    *GlobalNamespace_Ref
	refsForNodeMap *j2gort.OrderedMap[pkg_rhino.Node, *GlobalNamespace_Ref]
}

type GlobalNamespace_ModuleNameMap struct {
	names j2gort.OrderedMap[string, *GlobalNamespace_Name]
}

type GlobalNamespace struct {
	compiler    AbstractCompiler
	externsRoot pkg_rhino.Node
	root        pkg_rhino.Node
	globalNames []*GlobalNamespace_Name
	nameMap     j2gort.OrderedMap[string, *GlobalNamespace_Name]
	generated   bool
}

func (g *GlobalNamespace) IsJavaNil_() bool              { return g == nil }
func (*GlobalNamespace) Implements_rhino_StaticScope()   {}

func NewGlobalNamespace_3_AbstractCompiler_Node_Node(compiler AbstractCompiler, externsRoot pkg_rhino.Node, root pkg_rhino.Node) *GlobalNamespace {
	g := &GlobalNamespace{}
	InitGlobalNamespace_3_AbstractCompiler_Node_Node(g, compiler, externsRoot, root)
	return g
}

func InitGlobalNamespace_3_AbstractCompiler_Node_Node(this_ *GlobalNamespace, compiler AbstractCompiler, externsRoot pkg_rhino.Node, root pkg_rhino.Node) {
	this_.compiler = compiler
	this_.externsRoot = externsRoot
	this_.root = root
	this_.nameMap = j2gort.AcquireOrderedMapCap[string, *GlobalNamespace_Name](8)
}

func (this_ *GlobalNamespace) HasExternsRoot() bool {
	return this_.externsRoot != 0
}

func (this_ *GlobalNamespace) EnsureGenerated() {
	if !this_.generated {
		this_.priv_GlobalNamespace_process_3053()
	}
}

func (this_ *GlobalNamespace) ResetGeneratedForRetraverseTest() {
	this_.generated = false
}

func (this_ *GlobalNamespace) FirstDeclScopeRoot() pkg_rhino.Node {
	this_.EnsureGenerated()
	nm := this_.nameMap.Get("globalSym")
	if nm == nil || nm.declaration == nil || nm.declaration.scope == nil {
		return 0
	}
	return nm.declaration.scope.Root()
}

func (this_ *GlobalNamespace) priv_GlobalNamespace_process_3053() {
	var traversal *NodeTraversal_Builder = NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(NewGlobalNamespace_BuildGlobalNamespace(this_))
	_ = traversal
	if this_.HasExternsRoot() {
		traversal.TraverseRoots(this_.externsRoot, this_.root)
	} else {
		traversal.Traverse(this_.root)
	}
	this_.generated = true
}

type GlobalNamespace_BuildGlobalNamespace struct {
	outer_ *GlobalNamespace
}

func NewGlobalNamespace_BuildGlobalNamespace(outer_ *GlobalNamespace) *GlobalNamespace_BuildGlobalNamespace {
	return &GlobalNamespace_BuildGlobalNamespace{outer_: outer_}
}

func (b *GlobalNamespace_BuildGlobalNamespace) IsJavaNil_() bool                          { return b == nil }
func (*GlobalNamespace_BuildGlobalNamespace) Implements_jscomp_NodeTraversal_Callback()   {}
func (b *GlobalNamespace_BuildGlobalNamespace) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	ref := NewGlobalNamespace_Ref(t.GetScope(), n)
	nm := &GlobalNamespace_Name{
		name:           "globalSym",
		declaration:    ref,
		refsForNodeMap: j2gort.AcquireOrderedMapPtr[pkg_rhino.Node, *GlobalNamespace_Ref](),
	}
	nm.refsForNodeMap.Put(n, ref)
	b.outer_.nameMap.Put("globalSym", nm)
	b.outer_.globalNames = append(b.outer_.globalNames, nm)
	return true
}
func (b *GlobalNamespace_BuildGlobalNamespace) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
}

type InlineAndCollapseProperties struct {
	compiler                   AbstractCompiler
	namespace                  *GlobalNamespace
	testAggressiveInliningOnly bool
}

func (p *InlineAndCollapseProperties) IsJavaNil_() bool                 { return p == nil }
func (*InlineAndCollapseProperties) Implements_jscomp_CompilerPass()    {}

func NewInlineAndCollapseProperties(c AbstractCompiler) *InlineAndCollapseProperties {
	return &InlineAndCollapseProperties{compiler: c}
}

func (this_ *InlineAndCollapseProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.namespace = NewGlobalNamespace_3_AbstractCompiler_Node_Node(this_.compiler, externs, root)
	if this_.namespace.FirstDeclScopeRoot() == 0 {
		panic("missing decl scope during InlineAndCollapseProperties.Process_2")
	}
}

type Compiler struct {
	transpilationNamespace pkg_rhino.StaticScope
}

func (c *Compiler) IsJavaNil_() bool                    { return c == nil }
func (*Compiler) Implements_jscomp_AbstractCompiler()   {}

func NewCompiler() *Compiler { return &Compiler{} }

func (this_ *Compiler) GetTranspilationNamespace() pkg_rhino.StaticScope {
	if this_.transpilationNamespace == nil || this_.transpilationNamespace.IsJavaNil_() {
		gn := NewGlobalNamespace_3_AbstractCompiler_Node_Node(this_, pkg_rhino.Node(1), pkg_rhino.Node(2))
		gn.EnsureGenerated()
		this_.transpilationNamespace = gn
	}
	return this_.transpilationNamespace
}
`

	gnRuntimeTestSrc := `package gn_slab_test

import (
	"testing"

	"example.com/w40b_gn/gen/j2gort"
	pkg_rhino "example.com/w40b_gn/gen/rhino"
	"example.com/w40b_gn/opt/jscomp"
)

func TestGlobalNamespaceScopeCreatorLazyInitAndPassRelease(t *testing.T) {
	j2gort.ResetCounters()
	jscomp.ScannerAllocCount.Store(0)

	comp := jscomp.NewCompiler()

	// 1. Non-escaping GlobalNamespace in InlineAndCollapseProperties.Process_2:
	//    Ref.scope stays live during the pass, then ReleaseSlabs_() recycles scopeCreator_ slabs on exit!
	pass := jscomp.NewInlineAndCollapseProperties(comp)
	pass.Process_2(pkg_rhino.Node(10), pkg_rhino.Node(20))
	if got := j2gort.PooledSlabAllocCount.Load(); got < 2 {
		t.Fatalf("expected Scope + Var allocations on GlobalNamespace.scopeCreator_ slabs, got %d", got)
	}
	if got := j2gort.PooledSlabReleaseCount.Load(); got < 2 {
		t.Fatalf("expected InlineAndCollapseProperties.Process_2 exit to release scopeCreator_ slabs, got %d", got)
	}

	// 2. Lazy-init check: retraversing on the same GlobalNamespace reuses reusableScanner (1 allocation, not 2!).
	jscomp.ScannerAllocCount.Store(0)
	gn := jscomp.NewGlobalNamespace_3_AbstractCompiler_Node_Node(comp, 0, pkg_rhino.Node(30))
	gn.EnsureGenerated()
	gn.ResetGeneratedForRetraverseTest()
	gn.EnsureGenerated()
	if got := jscomp.ScannerAllocCount.Load(); got != 1 {
		t.Fatalf("expected lazy-init reusableScanner to allocate once across 2 traversals, got %d", got)
	}
	gn.ReleaseSlabs_()

	// 3. Compiler.transpilationNamespace MUST NOT be released!
	j2gort.ResetCounters()
	tns := comp.GetTranspilationNamespace().(*jscomp.GlobalNamespace)
	if got := j2gort.PooledSlabReleaseCount.Load(); got != 0 {
		t.Fatalf("transpilationNamespace MUST NOT call ReleaseSlabs_(), got releases=%d", got)
	}
	if tns.FirstDeclScopeRoot() == 0 {
		t.Fatal("transpilationNamespace scope root became nil")
	}
}
`

	wave40bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                            "module example.com/w40b_gn\n\ngo 1.23\n",
		"gen/j2gort/j2gort.go":              wave40bMockJ2gortSrc,
		"gen/rhino/node.go":                 rhinoSrc,
		"opt/jscomp/scope_and_nt.go":        jscompScopeAndNtSrc,
		"opt/jscomp/global_namespace.go":    globalNamespaceSrc,
		"testpkg/gn_scope_creator_test.go":  gnRuntimeTestSrc,
	})

	optDir := filepath.Join(tmpDir, "opt/jscomp")
	s1, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero sites on Run 1, got %+v", s1)
	}

	// Run 2 MUST be 100% idempotent!
	s2, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d ByPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	gnBytes, err := os.ReadFile(filepath.Join(optDir, "global_namespace.go"))
	if err != nil {
		t.Fatal(err)
	}
	gnOut := string(gnBytes)
	if strings.Contains(gnOut, "func (*GlobalNamespace_BuildGlobalNamespace) CanReleaseSyntacticScopes_()") {
		t.Fatal("GlobalNamespace_BuildGlobalNamespace MUST NOT receive CanReleaseSyntacticScopes_()")
	}
	for _, want := range []string{
		"scopeCreator_ SyntacticScopeCreator",
		"if this_.scopeCreator_.reusableScanner == nil {",
		"InitSyntacticScopeCreator_1(&this_.scopeCreator_, this_.compiler)",
		".SetScopeCreator(&this_.scopeCreator_)",
		"this_.scopeCreator_.ReleaseSlabs_()",
		"this_.namespace.ReleaseSlabs_()",
	} {
		if !strings.Contains(gnOut, want) {
			t.Fatalf("expected global_namespace.go to contain %q, got:\n%s", want, gnOut)
		}
	}

	wave40bRunGoTest(t, tmpDir)
}

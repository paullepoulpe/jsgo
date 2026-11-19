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

func TestWave36B_CheckNotNullPeepholeSpecializationAndRace(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/j2gort not ready right now; skipping")
	}

	tmpDir := t.TempDir()
	origModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/w36bcnn", 1) +
		fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, sErr := os.ReadFile(filepath.Join(repoRoot, "go.sum")); sErr == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	mustWrite := func(rel, content string) {
		t.Helper()
		full := filepath.Join(tmpDir, filepath.FromSlash(rel))
		if mkErr := os.MkdirAll(filepath.Dir(full), 0o755); mkErr != nil {
			t.Fatal(mkErr)
		}
		if wErr := os.WriteFile(full, []byte(content), 0o644); wErr != nil {
			t.Fatal(wErr)
		}
	}

	mustWrite("rhino/node.go", `package rhino

type Node uint32

func (n Node) IsJavaNil_() bool { return n == 0 }

type TokenID uint32
`)

	mustWrite("sample/sample.go", `package sample

import (
	"example.com/w36bcnn/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

type Widget struct {
	Val int32
}

type CustomPtr *Widget
type CustomSlice []int32

func CheckString(s string) string {
	j2gort.CheckNotNull(s)
	a := j2gort.CheckNotNull(s)
	b := j2gort.CheckNotNull[string](a)
	return b
}

func CheckPtr(w *Widget) *Widget {
	j2gort.CheckNotNull(w)
	a := j2gort.CheckNotNull(w)
	b := j2gort.CheckNotNull[*Widget](a)
	return b
}

func CheckSlice(xs []int32) []int32 {
	j2gort.CheckNotNull(xs)
	a := j2gort.CheckNotNull(xs)
	b := j2gort.CheckNotNull[[]int32](a)
	return b
}

func CheckNode(n rhino.Node) rhino.Node {
	j2gort.CheckNotNull(n)
	a := j2gort.CheckNotNull(n)
	b := j2gort.CheckNotNull[rhino.Node](a)
	return b
}

func CheckUntypedNils(flag bool) (*Widget, []int32) {
	if flag {
		return j2gort.CheckNotNull[*Widget](nil), nil
	}
	return nil, j2gort.CheckNotNull[[]int32](nil)
}

// Negative guards: exact return types and non-Node uint32 semantics MUST be preserved!
func NegativeReturnAndTypeGuards(s string, cp CustomPtr, cs CustomSlice, tok rhino.TokenID) (any, CustomPtr, CustomSlice, rhino.TokenID) {
	anyVal := j2gort.CheckNotNull[any](s)
	anyVal = int32(7)
	outP := j2gort.CheckNotNull(cp)
	outS := j2gort.CheckNotNull(cs)
	outTok := j2gort.CheckNotNull(tok)
	return anyVal, outP, outS, outTok
}
`)

	peepholeOnly := []passes.Pass{passes.NewJ2gortPeepholePass()}
	pkgDir := filepath.Join(tmpDir, "sample")

	s1, err := passes.RunOnPackageDir(pkgDir, peepholeOnly, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero CheckNotNull peephole specializations on Run 1")
	}

	s2, err := passes.RunOnPackageDir(pkgDir, peepholeOnly, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d byPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	gotBytes, err := os.ReadFile(filepath.Join(pkgDir, "sample.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"j2gort.CheckNotNullString(s)",
		"j2gort.CheckNotNullPtr(w)",
		"j2gort.CheckNotNullSlice(xs)",
		"j2gort.CheckNotNullNode(n)",
		"j2gort.CheckNotNullPtr[Widget](nil)",
		"j2gort.CheckNotNullSlice[int32](nil)",
		"j2gort.CheckNotNull[any](s)",
		"j2gort.CheckNotNull(cp)",
		"j2gort.CheckNotNull(cs)",
		"j2gort.CheckNotNull(tok)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in specialized sample.go:\n%s", want, got)
		}
	}

	mustWrite("sample/sample_test.go", `package sample_test

import (
	"testing"

	"example.com/w36bcnn/rhino"
	"example.com/w36bcnn/sample"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

func mustPanic(t *testing.T, label string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s: expected NPE panic, got none", label)
		}
	}()
	fn()
}

func TestSpecializedCheckNotNullSemantics(t *testing.T) {
	if got := sample.CheckString("ok"); got != "ok" {
		t.Fatalf("CheckString = %q, want ok", got)
	}
	if got := sample.CheckString(""); got != "" {
		t.Fatalf("CheckString empty = %q, want empty", got)
	}
	mustPanic(t, "NullString", func() { sample.CheckString(j2gort.NullString) })

	w := &sample.Widget{Val: 42}
	if got := sample.CheckPtr(w); got != w {
		t.Fatalf("CheckPtr mismatch")
	}
	mustPanic(t, "nil ptr", func() { sample.CheckPtr(nil) })

	xs := []int32{1, 2, 3}
	if got := sample.CheckSlice(xs); len(got) != 3 || got[2] != 3 {
		t.Fatalf("CheckSlice mismatch: %v", got)
	}
	if got := sample.CheckSlice([]int32{}); got == nil || len(got) != 0 {
		t.Fatalf("CheckSlice non-nil empty mismatch")
	}
	mustPanic(t, "nil slice", func() { sample.CheckSlice(nil) })

	if got := sample.CheckNode(rhino.Node(9)); got != 9 {
		t.Fatalf("CheckNode = %d, want 9", got)
	}
	mustPanic(t, "zero Node", func() { sample.CheckNode(0) })

	mustPanic(t, "untyped nil ptr", func() { sample.CheckUntypedNils(true) })
	mustPanic(t, "untyped nil slice", func() { sample.CheckUntypedNils(false) })

	// TokenID(0) is NOT rhino.Node and must NOT panic under CheckNotNull!
	anyV, cp, cs, tok := sample.NegativeReturnAndTypeGuards("hi", sample.CustomPtr(w), sample.CustomSlice(xs), rhino.TokenID(0))
	if anyV != int32(7) || cp != sample.CustomPtr(w) || len(cs) != 3 || tok != 0 {
		t.Fatalf("NegativeReturnAndTypeGuards mismatch: %v %v %v %v", anyV, cp, cs, tok)
	}
}
`)

	cmd := exec.Command("go", "test", "-race", "./sample")
	cmd.Dir = tmpDir
	if out, runErr := cmd.CombinedOutput(); runErr != nil {
		t.Fatalf("go test -race ./sample failed: %v\n%s", runErr, string(out))
	}
}

func TestWave36B_ImmutableListIterStrSliceLoweringAndRace(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	mustWrite := func(rel, content string) {
		t.Helper()
		full := filepath.Join(tmpDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	mustWrite("go.mod", "module example.com/w36bimmlist\n\ngo 1.24\n")
	mustWrite("j2gort/j2gort.go", `package j2gort

const NullString = "\xff\x00j2go_null"

type Ext_java_lang_Iterable interface {
	IsJavaNil_() bool
	Iterator() Ext_java_util_Iterator
}

type Ext_java_util_Iterator interface {
	IsJavaNil_() bool
	HasNext() bool
	Next() any
}

type Ext_com_google_common_collect_ImmutableList interface {
	Ext_java_lang_Iterable
	IterSlice_() []any
	IterStrSlice_() ([]string, bool)
	Size() int32
}

type strBackedImmList struct {
	strs []string
}

func NewStrImmList(s ...string) Ext_com_google_common_collect_ImmutableList {
	if s == nil {
		s = []string{}
	}
	return &strBackedImmList{strs: s}
}

func (l *strBackedImmList) IsJavaNil_() bool                 { return l == nil }
func (l *strBackedImmList) Size() int32                      { return int32(len(l.strs)) }
func (l *strBackedImmList) IterStrSlice_() ([]string, bool)  { return l.strs, true }
func (l *strBackedImmList) IterSlice_() []any {
	out := make([]any, len(l.strs))
	for i, s := range l.strs {
		out[i] = s
	}
	return out
}
func (l *strBackedImmList) Iterator() Ext_java_util_Iterator {
	return &sliceIter{items: l.IterSlice_()}
}

type anyBackedImmList struct {
	items []any
}

func NewAnyImmList(items ...any) Ext_com_google_common_collect_ImmutableList {
	return &anyBackedImmList{items: items}
}

func NewTypedNilImmList() Ext_com_google_common_collect_ImmutableList {
	var p *anyBackedImmList
	return p
}

func (l *anyBackedImmList) IsJavaNil_() bool                { return l == nil }
func (l *anyBackedImmList) Size() int32                     { return int32(len(l.items)) }
func (l *anyBackedImmList) IterStrSlice_() ([]string, bool) { return nil, false }
func (l *anyBackedImmList) IterSlice_() []any               { return l.items }
func (l *anyBackedImmList) Iterator() Ext_java_util_Iterator {
	return &sliceIter{items: l.items}
}

type sliceIter struct {
	items []any
	idx   int
}

func (it *sliceIter) IsJavaNil_() bool { return it == nil }
func (it *sliceIter) HasNext() bool    { return it.idx < len(it.items) }
func (it *sliceIter) Next() any {
	v := it.items[it.idx]
	it.idx++
	return v
}

type IterCursor struct {
	strs  []string
	items []any
	idx   int
	isStr bool
}

func NewIterCursor(src any) IterCursor {
	if src == nil {
		panic("NPE")
	}
	if l, ok := src.(Ext_com_google_common_collect_ImmutableList); ok {
		if l.IsJavaNil_() {
			panic("NPE")
		}
		if ss, okS := l.IterStrSlice_(); okS {
			return IterCursor{strs: ss, isStr: true}
		}
		return IterCursor{items: l.IterSlice_()}
	}
	panic("unsupported iterable")
}

func (c *IterCursor) HasNext() bool {
	if c.isStr {
		return c.idx < len(c.strs)
	}
	return c.idx < len(c.items)
}

func (c *IterCursor) NextStr() string {
	if c.isStr {
		s := c.strs[c.idx]
		c.idx++
		return s
	}
	v := c.items[c.idx]
	c.idx++
	return v.(string)
}

func (c *IterCursor) Next() any {
	if c.isStr {
		s := c.strs[c.idx]
		c.idx++
		return s
	}
	v := c.items[c.idx]
	c.idx++
	return v
}

func CheckNotNull[T any](v T) T {
	if any(v) == nil {
		panic("NPE")
	}
	if c, ok := any(v).(interface{ IsJavaNil_() bool }); ok && c.IsJavaNil_() {
		panic("NPE")
	}
	return v
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

func CastString(v any) string {
	if v == nil {
		return NullString
	}
	return v.(string)
}

func AsIface[T any](v any) T {
	if v == nil {
		var zero T
		return zero
	}
	return v.(T)
}

func ImmListFromTypedSlice[T comparable](s []T) Ext_com_google_common_collect_ImmutableList {
	return nil
}

func ImmListToSliceNullable[T any](l Ext_com_google_common_collect_ImmutableList) []T {
	if l == nil || l.IsJavaNil_() {
		return nil
	}
	raw := l.IterSlice_()
	out := make([]T, len(raw))
	for i, v := range raw {
		out[i] = v.(T)
	}
	return out
}
`)

	mustWrite("consumer/consumer.go", `package consumer

import "example.com/w36bimmlist/j2gort"

type Provider interface {
	Fetch() j2gort.Ext_com_google_common_collect_ImmutableList
}

type ProviderImpl struct {
	calls int32
	List  j2gort.Ext_com_google_common_collect_ImmutableList
}

func (p *ProviderImpl) Fetch() j2gort.Ext_com_google_common_collect_ImmutableList {
	p.calls++
	return p.List
}

// 1. NewIterCursor + NextStr() with nested ImmutableList<String> loop, continue, break, and mutated loop var!
func JoinNestedCursor(p Provider, pfxProv Provider) []string {
	var out []string
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](p.Fetch())))
		for cur_1.HasNext() {
			item := cur_1.NextStr()
			_ = item
			if item == "skip" {
				continue
			}
			{
				cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(pfxProv.Fetch()))
				for cur_2.HasNext() {
					pfx := cur_2.NextStr()
					_ = pfx
					if pfx == "stop" {
						item = item + "!"
						break
					}
					item = pfx + ":" + item
				}
			}
			out = append(out, item)
		}
	}
	return out
}

// 2. Range over IterSlice_() + CastString on interface return
func CollectRangeIterSlice(p Provider) string {
	res := ""
	for _, raw := range j2gort.CheckNotNull(p.Fetch()).IterSlice_() {
		_ = raw
		s := j2gort.CastString(raw)
		res += "[" + s + "]"
	}
	return res
}

// 3. Iterator() + CastString(it.Next()) on interface return
func CollectIterator(p Provider) string {
	res := ""
	it := j2gort.CheckNotNull(p.Fetch()).Iterator()
	_ = it
	for it.HasNext() {
		s := j2gort.CastString(it.Next())
		res += "(" + s + ")"
	}
	return res
}

// 4. Negative guard: labeled loop MUST NOT duplicate body!
func LabeledGuard(p Provider) int {
	n := 0
outer:
	for _, raw := range j2gort.CheckNotNull(p.Fetch()).IterSlice_() {
		s := j2gort.CastString(raw)
		if s == "x" {
			break outer
		}
		n++
	}
	return n
}
`)

	immPass := []passes.Pass{passes.NewImmutableListLoweringPass()}
	pkgDir := filepath.Join(tmpDir, "consumer")

	s1, err := passes.RunOnPackageDir(pkgDir, immPass, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}

	s2, err := passes.RunOnPackageDir(pkgDir, immPass, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d", s2.TotalSites, s2.FilesModified)
	}

	gotBytes, err := os.ReadFile(filepath.Join(pkgDir, "consumer.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	if s1.TotalSites != 4 {
		t.Fatalf("expected 4 IterStrSlice_ loop rewrites (2 nested + 1 range + 1 iterator), got %d:\n%s", s1.TotalSites, got)
	}
	if strings.Count(got, "IterStrSlice_()") != 5 {
		t.Fatalf("expected 5 IterStrSlice_() call sites (1 outer + 2 inner copies + 1 range + 1 iterator) in consumer.go, got:\n%s", got)
	}

	mustWrite("consumer/consumer_test.go", `package consumer

import (
	"reflect"
	"testing"

	"example.com/w36bimmlist/j2gort"
)

func mustPanicNPE(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("expected NPE panic")
		}
	}()
	fn()
}

func TestIterStrSliceFastAndFallbackPaths(t *testing.T) {
	strMain := j2gort.NewStrImmList("a", "skip", "b")
	strPfx := j2gort.NewStrImmList("p1", "stop", "p2")
	anyMain := j2gort.NewAnyImmList("a", "skip", "b")
	anyPfx := j2gort.NewAnyImmList("p1", "stop", "p2")

	p1 := &ProviderImpl{List: strMain}
	gotFast := JoinNestedCursor(p1, &ProviderImpl{List: strPfx})
	want := []string{"p1:a!", "p1:b!"}
	if !reflect.DeepEqual(gotFast, want) || p1.calls != 1 {
		t.Fatalf("fast path = %v (calls=%d), want %v (calls=1)", gotFast, p1.calls, want)
	}

	// Verify 0 allocations on the IterStrSlice_() fast path when not appending!
	emptyProv := &ProviderImpl{List: j2gort.NewStrImmList("skip", "skip")}
	pfxProvFast := &ProviderImpl{List: strPfx}
	allocs := testing.AllocsPerRun(50, func() {
		_ = JoinNestedCursor(emptyProv, pfxProvFast)
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs on IterStrSlice_() fast path, got %v", allocs)
	}

	p2 := &ProviderImpl{List: anyMain}
	gotSlow := JoinNestedCursor(p2, &ProviderImpl{List: anyPfx})
	if !reflect.DeepEqual(gotSlow, want) || p2.calls != 1 {
		t.Fatalf("fallback path = %v (calls=%d), want %v (calls=1)", gotSlow, p2.calls, want)
	}

	if r1, r2 := CollectRangeIterSlice(&ProviderImpl{List: strMain}), CollectRangeIterSlice(&ProviderImpl{List: anyMain}); r1 != "[a][skip][b]" || r2 != r1 {
		t.Fatalf("CollectRangeIterSlice mismatch: %q vs %q", r1, r2)
	}
	if i1, i2 := CollectIterator(&ProviderImpl{List: strMain}), CollectIterator(&ProviderImpl{List: anyMain}); i1 != "(a)(skip)(b)" || i2 != i1 {
		t.Fatalf("CollectIterator mismatch: %q vs %q", i1, i2)
	}
	if n := LabeledGuard(&ProviderImpl{List: strMain}); n != 3 {
		t.Fatalf("LabeledGuard = %d, want 3", n)
	}

	// Nil and typed-nil ImmutableList MUST still throw NPE!
	mustPanicNPE(t, func() { CollectRangeIterSlice(&ProviderImpl{List: nil}) })
	mustPanicNPE(t, func() { CollectIterator(&ProviderImpl{List: j2gort.NewTypedNilImmList()}) })
}
`)

	cmd := exec.Command("go", "test", "-race", "./consumer")
	cmd.Dir = tmpDir
	if out, runErr := cmd.CombinedOutput(); runErr != nil {
		t.Fatalf("go test -race ./consumer failed: %v\n%s", runErr, string(out))
	}
}

func TestWave36B_SyntheticAllSlabTargetsAndRace(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/j2gort not ready right now; skipping")
	}

	tmpDir := t.TempDir()
	origModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/w36bslabs", 1) +
		fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, sErr := os.ReadFile(filepath.Join(repoRoot, "go.sum")); sErr == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	mustWrite := func(rel, content string) {
		t.Helper()
		full := filepath.Join(tmpDir, filepath.FromSlash(rel))
		if mkErr := os.MkdirAll(filepath.Dir(full), 0o755); mkErr != nil {
			t.Fatal(mkErr)
		}
		if wErr := os.WriteFile(full, []byte(content), 0o644); wErr != nil {
			t.Fatal(wErr)
		}
	}

	mustWrite("rhino/node.go", `package rhino

type Node uint32

func (n Node) IsJavaNil_() bool { return n == 0 }
`)

	// Task 1: ControlFlowGraphBase.ResetCfgForPool upgrading/emitting slab_LinearFlowState.Release()
	mustWrite("jscomp/cfg_dfa.go", `package jscomp

import (
	"example.com/w36bslabs/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

type cfgNode struct {
	ann any
}

func (n *cfgNode) SetAnnotation(a any) { n.ann = a }
func (n *cfgNode) GetAnnotation() any  { return n.ann }

type ControlFlowGraph[N comparable] interface {
	IsJavaNil_() bool
	GetControlFlowGraphBase() *ControlFlowGraphBase[N]
}

type ControlFlowGraphBase[N comparable] struct {
	entry                N
	nodes                []*cfgNode
	slab_LinearFlowState j2gort.PooledSlab[DataFlowAnalysis_LinearFlowState]
}

func (this_ *ControlFlowGraphBase[N]) IsJavaNil_() bool { return this_ == nil }
func (this_ *ControlFlowGraphBase[N]) GetControlFlowGraphBase() *ControlFlowGraphBase[N] {
	return this_
}
func (this_ *ControlFlowGraphBase[N]) ResetCfgForPool() bool {
	this_.slab_LinearFlowState.Reset()
	this_.nodes = this_.nodes[:0]
	return true
}

type DataFlowAnalysis_LinearFlowState struct {
	stepCount int32
	inL       any
	outL      any
}

func NewDataFlowAnalysis_LinearFlowState(inL any, outL any) *DataFlowAnalysis_LinearFlowState {
	this_ := &DataFlowAnalysis_LinearFlowState{}
	InitDataFlowAnalysis_LinearFlowState(this_, inL, outL)
	return this_
}

func InitDataFlowAnalysis_LinearFlowState(this_ *DataFlowAnalysis_LinearFlowState, inL any, outL any) {
	this_.stepCount = 0
	this_.inL = inL
	this_.outL = outL
}

type DataFlowAnalysisBase[N comparable, L any] struct {
	cfg ControlFlowGraph[N]
}

func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_initialize_2310(n *cfgNode, initIn L, initOut L) {
	n.SetAnnotation(NewDataFlowAnalysis_LinearFlowState(initIn, initOut))
}

var _ = rhino.Node(0)
`)

	// Task 4A: MakeDeclaredNamesUnique ContextualRenamer lazy declarations + ContextualRenameInverter slabs
	mustWrite("jscomp/mdnu.go", `package jscomp

import (
	"fmt"

	"example.com/w36bslabs/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

type AbstractCompiler interface {
	IsJavaNil_() bool
}

type ScopeCreator interface {
	IsJavaNil_() bool
}

type Var struct {
	name  string
	depth int32
}

func (v *Var) GetName() string { return v.name }

type Scope struct {
	depth int32
	vars  []*Var
}

func (s *Scope) GetDepth() int32 { return s.depth }

type NodeTraversal struct {
	scopeRoot rhino.Node
	inGlobal  bool
	scope     *Scope
}

func (t *NodeTraversal) GetScopeRoot() rhino.Node { return t.scopeRoot }
func (t *NodeTraversal) InGlobalScope() bool      { return t.inGlobal }
func (t *NodeTraversal) GetScope() *Scope         { return t.scope }
func (t *NodeTraversal) InFunctionBlockScope() bool {
	return !t.inGlobal
}

type MakeDeclaredNamesUnique_Renamer interface {
	IsJavaNil_() bool
	AddDeclaredName(name string, hoisted bool)
	GetReplacementName(oldName string) string
	CreateForChildScope(scopeRoot rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer
}

type MakeDeclaredNamesUnique_ContextualRenamer interface {
	MakeDeclaredNamesUnique_Renamer
}

type MakeDeclaredNamesUnique_ContextualRenamer_NameUsage struct {
	counts map[string]int32
}

func NewMakeDeclaredNamesUnique_ContextualRenamer_NameUsage() *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage {
	return &MakeDeclaredNamesUnique_ContextualRenamer_NameUsage{counts: make(map[string]int32)}
}

type MakeDeclaredNamesUnique_ContextualRenamerBase struct {
	self_        any
	nameUsage    *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage
	declarations *j2gort.OrderedMap[string, string]
	scopeRoot    rhino.Node
	declName0    string
	declRep0     string
	declName1    string
	declRep1     string
	declName2    string
	declRep2     string
	declName3    string
	declRep3     string
	declCount    int32
	global       bool
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) BindSelf_(s any)  { this_.self_ = s }
func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) IsJavaNil_() bool { return this_ == nil }

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) initFields_MakeDeclaredNamesUnique_ContextualRenamer() {
	_ = this_
	this_.declName0 = j2gort.NullString
	this_.declRep0 = j2gort.NullString
	this_.declName1 = j2gort.NullString
	this_.declRep1 = j2gort.NullString
	this_.declName2 = j2gort.NullString
	this_.declRep2 = j2gort.NullString
	this_.declName3 = j2gort.NullString
	this_.declRep3 = j2gort.NullString
	this_.declarations = &j2gort.OrderedMap[string, string]{}
}

func NewMakeDeclaredNamesUnique_ContextualRenamer() MakeDeclaredNamesUnique_ContextualRenamer {
	this_ := &MakeDeclaredNamesUnique_ContextualRenamerBase{}
	this_.BindSelf_(this_)
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenamer()
	this_.global = true
	this_.nameUsage = NewMakeDeclaredNamesUnique_ContextualRenamer_NameUsage()
	return this_
}

func NewMakeDeclaredNamesUnique_ContextualRenamer_4(scopeRoot rhino.Node, nameUsage *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage, hoistingTargetScope bool, parent MakeDeclaredNamesUnique_Renamer) MakeDeclaredNamesUnique_ContextualRenamer {
	this_ := &MakeDeclaredNamesUnique_ContextualRenamerBase{}
	this_.BindSelf_(this_)
	InitMakeDeclaredNamesUnique_ContextualRenamer_4(this_, scopeRoot, nameUsage, hoistingTargetScope, parent)
	return this_
}

func InitMakeDeclaredNamesUnique_ContextualRenamer_4(this_ *MakeDeclaredNamesUnique_ContextualRenamerBase, scopeRoot rhino.Node, nameUsage *MakeDeclaredNamesUnique_ContextualRenamer_NameUsage, hoistingTargetScope bool, parent MakeDeclaredNamesUnique_Renamer) {
	_ = hoistingTargetScope
	_ = parent
	this_.initFields_MakeDeclaredNamesUnique_ContextualRenamer()
	this_.scopeRoot = scopeRoot
	this_.global = false
	this_.nameUsage = nameUsage
	this_.declCount = 0
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) CreateForChildScope(scopeRoot rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	return j2gort.AsIface[MakeDeclaredNamesUnique_Renamer](NewMakeDeclaredNamesUnique_ContextualRenamer_4(scopeRoot, this_.nameUsage, hoistingTargetScope, j2gort.AsIface[MakeDeclaredNamesUnique_Renamer](this_.self_)))
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) priv_ContextualRenamer_putLocalDeclaration_3831(name string, newName string) {
	c := this_.declCount
	switch c {
	case int32(0):
		this_.declName0 = name
		this_.declRep0 = newName
		this_.declCount = 1
	case int32(1):
		this_.declName1 = name
		this_.declRep1 = newName
		this_.declCount = 2
	case int32(2):
		this_.declName2 = name
		this_.declRep2 = newName
		this_.declCount = 3
	case int32(3):
		this_.declName3 = name
		this_.declRep3 = newName
		this_.declCount = 4
	case int32(4):
		this_.declarations.Put(this_.declName0, this_.declRep0)
		this_.declarations.Put(this_.declName1, this_.declRep1)
		this_.declarations.Put(this_.declName2, this_.declRep2)
		this_.declarations.Put(this_.declName3, this_.declRep3)
		this_.declarations.Put(name, newName)
		this_.declCount = 5
	default:
		this_.declarations.Put(name, newName)
	}
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) AddDeclaredName(name string, hoisted bool) {
	_ = hoisted
	if this_.global {
		this_.nameUsage.counts[name] = 1
		return
	}
	n := this_.nameUsage.counts[name]
	this_.nameUsage.counts[name] = n + 1
	rep := j2gort.NullString
	if n > 0 {
		rep = fmt.Sprintf("%s$jscomp$%d", name, n)
	}
	this_.priv_ContextualRenamer_putLocalDeclaration_3831(name, rep)
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) GetReplacementName(oldName string) string {
	c := this_.declCount
	if c == 0 {
		return j2gort.NullString
	}
	if c <= 4 {
		if oldName == this_.declName0 {
			return this_.declRep0
		}
		if c > 1 && oldName == this_.declName1 {
			return this_.declRep1
		}
		if c > 2 && oldName == this_.declName2 {
			return this_.declRep2
		}
		if c > 3 && oldName == this_.declName3 {
			return this_.declRep3
		}
		return j2gort.NullString
	}
	return j2gort.OrderedMapPtrGetStr(this_.declarations, oldName)
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenamerBase) String() string {
	var declView *j2gort.OrderedMap[string, string] = this_.declarations
	if this_.declCount > 0 && this_.declCount <= 4 {
		declView = j2gort.NewOrderedMapPtrCap[string, string](this_.declCount)
		declView.Put(this_.declName0, this_.declRep0)
	}
	return fmt.Sprintf("decls=%d", declView.Size())
}

var _ = fmt.Sprint

type MakeDeclaredNamesUnique struct {
	renamerStack j2gort.Deque[MakeDeclaredNamesUnique_Renamer]
	rootRenamer  MakeDeclaredNamesUnique_Renamer
}

func NewMakeDeclaredNamesUnique(root MakeDeclaredNamesUnique_Renamer) *MakeDeclaredNamesUnique {
	return &MakeDeclaredNamesUnique{rootRenamer: root}
}

func (this_ *MakeDeclaredNamesUnique) EnterScope(t *NodeTraversal) {
	var renamer MakeDeclaredNamesUnique_Renamer
	if this_.renamerStack.IsEmpty() {
		renamer = this_.rootRenamer
	} else {
		renamer = this_.renamerStack.Peek().CreateForChildScope(t.GetScopeRoot(), false)
	}
	this_.renamerStack.Push(renamer)
	for _, v := range t.GetScope().vars {
		this_.renamerStack.Peek().AddDeclaredName(v.GetName(), false)
	}
}

func (this_ *MakeDeclaredNamesUnique) ExitScope(t *NodeTraversal) {
	if !(t.InGlobalScope()) {
		this_.renamerStack.Pop()
	}
}

// ContextualRenameInverter and its 4 pass-bounded inner structs:
type MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo interface {
	IsJavaNil_() bool
	GetName() string
}

type MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo struct {
	outer_        *MakeDeclaredNamesUnique_ContextualRenameInverter
	currentName   string
	preferredName string
}

func (r *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) IsJavaNil_() bool {
	return r == nil
}
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(outer_ *MakeDeclaredNamesUnique_ContextualRenameInverter, cur string, pref string) *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo {
	this_ := &MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo{}
	this_.outer_ = outer_
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(this_, cur, pref)
	return this_
}
func InitMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo, cur string, pref string) {
	this_.currentName = cur
	this_.preferredName = pref
}

type MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo struct {
	var_         *Var
	name         string
	renamingInfo *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo
}

func (d *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) IsJavaNil_() bool {
	return d == nil
}
func (d *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) GetName() string {
	return d.name
}
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(v *Var, ri *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo {
	this_ := &MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo{}
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(this_, v, ri)
	return this_
}
func InitMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo, v *Var, ri *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) {
	this_.var_ = v
	this_.name = v.GetName()
	this_.renamingInfo = ri
}

type MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo struct {
	name string
}

func (u *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) IsJavaNil_() bool {
	return u == nil
}
func (u *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) GetName() string {
	return u.name
}
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(name string) *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo {
	this_ := &MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo{}
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(this_, name)
	return this_
}
func InitMakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo, name string) {
	this_.name = name
}

type MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext struct {
	parent *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext
	scope  *Scope
	decls  []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo
}

func NewMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(s *Scope, decls []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo, parent *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext {
	this_ := &MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext{}
	InitMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(this_, s, decls, parent)
	return this_
}
func InitMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(this_ *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext, s *Scope, decls []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo, parent *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) {
	this_.scope = s
	this_.decls = decls
	this_.parent = parent
}

type MakeDeclaredNamesUnique_ContextualRenameInverter struct {
	compiler          AbstractCompiler
	scopeContextStack j2gort.Deque[*MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext]
	invertedCount     int32
}

func NewMakeDeclaredNamesUnique_ContextualRenameInverter(c AbstractCompiler) *MakeDeclaredNamesUnique_ContextualRenameInverter {
	return &MakeDeclaredNamesUnique_ContextualRenameInverter{compiler: c}
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) EnterScope(s *Scope) MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
	var decls []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo
	for _, v := range s.vars {
		ri := NewMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(this_, v.name, v.name)
		dv := NewMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(v, ri)
		decls = append(decls, dv)
	}
	sc := NewMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(s, decls, this_.scopeContextStack.Peek())
	this_.scopeContextStack.Push(sc)
	return NewMakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo("ext")
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) ExitScope() {
	sc := this_.scopeContextStack.Pop()
	this_.invertedCount += int32(len(sc.decls))
}

func (this_ *MakeDeclaredNamesUnique_ContextualRenameInverter) Process_2(externs rhino.Node, js rhino.Node) {
	_ = externs
	_ = js
	s := &Scope{depth: 1, vars: []*Var{{name: "a$jscomp$1"}, {name: "b$jscomp$1"}}}
	_ = this_.EnterScope(s)
	this_.ExitScope()
}
`)

	// Task 4B: CrossChunkReferenceCollector_TopLevelStatement + CrossChunkCodeMotion
	mustWrite("jscomp/ccrc.go", `package jscomp

import "example.com/w36bslabs/rhino"

type JSChunk struct {
	id int32
}

type CrossChunkReferenceCollector_TopLevelStatement struct {
	outer_        *CrossChunkReferenceCollector
	originalOrder int32
	chunk         *JSChunk
	statementNode rhino.Node
}

func NewCrossChunkReferenceCollector_TopLevelStatement(outer_ *CrossChunkReferenceCollector, order int32, chunk *JSChunk, stmt rhino.Node) *CrossChunkReferenceCollector_TopLevelStatement {
	this_ := &CrossChunkReferenceCollector_TopLevelStatement{}
	this_.outer_ = outer_
	InitCrossChunkReferenceCollector_TopLevelStatement(this_, order, chunk, stmt)
	return this_
}

func InitCrossChunkReferenceCollector_TopLevelStatement(this_ *CrossChunkReferenceCollector_TopLevelStatement, order int32, chunk *JSChunk, stmt rhino.Node) {
	this_.originalOrder = order
	this_.chunk = chunk
	this_.statementNode = stmt
}

type CrossChunkReferenceCollector struct {
	statementCounter       int32
	topLevelStatements     []*CrossChunkReferenceCollector_TopLevelStatement
	topLevelStatementDraft *CrossChunkReferenceCollector_TopLevelStatement
}

func NewCrossChunkReferenceCollector(compiler AbstractCompiler, creator ScopeCreator) *CrossChunkReferenceCollector {
	_ = compiler
	_ = creator
	return &CrossChunkReferenceCollector{}
}

func (this_ *CrossChunkReferenceCollector) priv_CrossChunkReferenceCollector_initializeDraftStatement_2281(chunk *JSChunk, statementNode rhino.Node) *CrossChunkReferenceCollector_TopLevelStatement {
	var val_4 CrossChunkReferenceCollector_TopLevelStatement
	val_4.outer_ = this_
	InitCrossChunkReferenceCollector_TopLevelStatement(&val_4, func() int32 {
		old_5 := this_.statementCounter
		this_.statementCounter = old_5 + 1
		return old_5
	}(), chunk, statementNode)
	var draft *CrossChunkReferenceCollector_TopLevelStatement = &val_4
	return draft
}

func (this_ *CrossChunkReferenceCollector) Process_1(root rhino.Node) {
	ch := &JSChunk{id: 1}
	for i := rhino.Node(1); i <= root; i++ {
		d := this_.priv_CrossChunkReferenceCollector_initializeDraftStatement_2281(ch, i)
		this_.topLevelStatementDraft = d
		this_.topLevelStatements = append(this_.topLevelStatements, d)
		this_.topLevelStatementDraft = nil
	}
}

func (this_ *CrossChunkReferenceCollector) GetTopLevelStatementCount() int32 {
	return int32(len(this_.topLevelStatements))
}

func (this_ *CrossChunkReferenceCollector) GetTopLevelStatement(i int32) *CrossChunkReferenceCollector_TopLevelStatement {
	return this_.topLevelStatements[i]
}

type CrossChunkCodeMotion struct {
	compiler AbstractCompiler
	movedSum int32
}

func (this_ *CrossChunkCodeMotion) IsJavaNil_() bool { return this_ == nil }

func (this_ *CrossChunkCodeMotion) Process_2(externs rhino.Node, root rhino.Node) {
	_ = externs
	var referenceCollector *CrossChunkReferenceCollector = NewCrossChunkReferenceCollector(this_.compiler, nil)
	_ = referenceCollector
	referenceCollector.Process_1(root)
	sum := int32(0)
	for i := int32(0); i < referenceCollector.GetTopLevelStatementCount(); i++ {
		st := referenceCollector.GetTopLevelStatement(i)
		sum += int32(st.statementNode) + st.originalOrder
	}
	this_.movedSum = sum
}
`)

	// Task 4C: OptimizeParameters + DeadAssignmentsElimination + InlineVariables
	mustWrite("jscomp/opt_dae_iv.go", `package jscomp

import (
	"example.com/w36bslabs/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

type OptimizeCalls_NodeList struct {
	nodes []rhino.Node
}

type OptimizeCalls_ReferenceMap struct {
	names j2gort.OrderedMap[string, *OptimizeCalls_NodeList]
	props j2gort.OrderedMap[string, *OptimizeCalls_NodeList]
}

type OptimizeParameters_CandidateAnalysis struct {
	outer_ *OptimizeParameters
	safe   bool
}

func NewOptimizeParameters_CandidateAnalysis(outer_ *OptimizeParameters) *OptimizeParameters_CandidateAnalysis {
	this_ := &OptimizeParameters_CandidateAnalysis{outer_: outer_}
	InitOptimizeParameters_CandidateAnalysis(this_)
	return this_
}
func InitOptimizeParameters_CandidateAnalysis(this_ *OptimizeParameters_CandidateAnalysis) {
	this_.safe = false
}

type OptimizeParameters_Parameter struct {
	arg          rhino.Node
	shouldRemove bool
}

func NewOptimizeParameters_Parameter(arg rhino.Node, shouldRemove bool) *OptimizeParameters_Parameter {
	this_ := &OptimizeParameters_Parameter{}
	InitOptimizeParameters_Parameter(this_, arg, shouldRemove)
	return this_
}
func InitOptimizeParameters_Parameter(this_ *OptimizeParameters_Parameter, arg rhino.Node, shouldRemove bool) {
	this_.arg = arg
	this_.shouldRemove = shouldRemove
}

type OptimizeParameters_ParameterList struct {
	params []*OptimizeParameters_Parameter
}

func NewOptimizeParameters_ParameterList() *OptimizeParameters_ParameterList {
	this_ := &OptimizeParameters_ParameterList{}
	InitOptimizeParameters_ParameterList(this_)
	return this_
}
func InitOptimizeParameters_ParameterList(this_ *OptimizeParameters_ParameterList) {
	this_.params = nil
}

type OptimizeParameters struct {
	compiler  AbstractCompiler
	optimized int32
}

func (this_ *OptimizeParameters) priv_OptimizeParameters_analyzeCandidate_4062(kind string, key string, refs *OptimizeCalls_NodeList) *OptimizeParameters_CandidateAnalysis {
	_ = kind
	_ = key
	var val_4 OptimizeParameters_CandidateAnalysis
	val_4.outer_ = this_
	InitOptimizeParameters_CandidateAnalysis(&val_4)
	var analysis *OptimizeParameters_CandidateAnalysis = &val_4
	analysis.safe = len(refs.nodes) > 0
	return analysis
}

func (this_ *OptimizeParameters) priv_OptimizeParameters_analyzeCandidateName_4060(key string, refs *OptimizeCalls_NodeList) *OptimizeParameters_CandidateAnalysis {
	return this_.priv_OptimizeParameters_analyzeCandidate_4062("name", key, refs)
}

func (this_ *OptimizeParameters) priv_OptimizeParameters_analyzeCandidateProperty_4061(key string, refs *OptimizeCalls_NodeList) *OptimizeParameters_CandidateAnalysis {
	return this_.priv_OptimizeParameters_analyzeCandidate_4062("prop", key, refs)
}

func (this_ *OptimizeParameters) priv_OptimizeParameters_findFixedArguments_4070(refs *OptimizeCalls_NodeList) *OptimizeParameters_ParameterList {
	var val_12 OptimizeParameters_ParameterList
	InitOptimizeParameters_ParameterList(&val_12)
	var parameters *OptimizeParameters_ParameterList = &val_12
	for _, n := range refs.nodes {
		p := NewOptimizeParameters_Parameter(n, n%2 == 0)
		parameters.params = append(parameters.params, p)
	}
	return parameters
}

func (this_ *OptimizeParameters) priv_OptimizeParameters_tryEliminateConstantArgs_4069(refs *OptimizeCalls_NodeList) {
	_ = refs
	pl := this_.priv_OptimizeParameters_findFixedArguments_4070(refs)
	for _, p := range pl.params {
		if p.shouldRemove {
			this_.optimized++
		}
	}
}

func (this_ *OptimizeParameters) Process_3(externs rhino.Node, root rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	_ = externs
	_ = root
	var toOptimize []*OptimizeCalls_NodeList
	{
		cur_1 := refMap.names.IterEntries()
		for cur_1.NextEntry() {
			ca := this_.priv_OptimizeParameters_analyzeCandidateName_4060(cur_1.Key(), cur_1.Value())
			if ca.safe {
				toOptimize = append(toOptimize, cur_1.Value())
			}
		}
	}
	{
		cur_2 := refMap.props.IterEntries()
		for cur_2.NextEntry() {
			ca2 := this_.priv_OptimizeParameters_analyzeCandidateProperty_4061(cur_2.Key(), cur_2.Value())
			if ca2.safe {
				toOptimize = append(toOptimize, cur_2.Value())
			}
		}
	}
	for _, refs := range toOptimize {
		this_.priv_OptimizeParameters_tryEliminateConstantArgs_4069(refs)
	}
}

type DeadAssignmentsElimination_BailoutInformation struct {
	containsFunction        bool
	containsRemovableAssign bool
}

func NewDeadAssignmentsElimination_BailoutInformation() *DeadAssignmentsElimination_BailoutInformation {
	this_ := &DeadAssignmentsElimination_BailoutInformation{}
	InitDeadAssignmentsElimination_BailoutInformation(this_)
	return this_
}
func InitDeadAssignmentsElimination_BailoutInformation(this_ *DeadAssignmentsElimination_BailoutInformation) {
	this_.containsFunction = false
	this_.containsRemovableAssign = false
}

type DeadAssignmentsElimination struct {
	functionStack j2gort.Deque[*DeadAssignmentsElimination_BailoutInformation]
	eliminated    int32
}

func (this_ *DeadAssignmentsElimination) EnterScopeWithCfg(t *NodeTraversal) {
	if t.InFunctionBlockScope() {
		this_.functionStack.AddFirst(NewDeadAssignmentsElimination_BailoutInformation())
	}
}

func (this_ *DeadAssignmentsElimination) ExitScopeWithCfg(t *NodeTraversal) {
	if t.InFunctionBlockScope() {
		top := this_.functionStack.PeekFirst()
		if !top.containsFunction && top.containsRemovableAssign {
			this_.eliminated++
		}
		this_.functionStack.RemoveFirst()
	}
}

func (this_ *DeadAssignmentsElimination) Process_2(externs rhino.Node, root rhino.Node) {
	_ = externs
	t := &NodeTraversal{scopeRoot: root, inGlobal: false}
	this_.EnterScopeWithCfg(t)
	this_.functionStack.PeekFirst().containsRemovableAssign = true
	this_.ExitScopeWithCfg(t)
}

type InlineVariables_VarExpert struct {
	v       *Var
	aliased *Var
}

type InlineVariables_InlineVarAnalysis struct {
	inlined bool
}

type InlineVariables_InliningBehavior_AliasInlineRetryHandler struct {
	outer_ *InlineVariables_InliningBehavior
	v      *Var
	expert *InlineVariables_VarExpert
}

type InlineVariables struct {
	slab_InlineVariables_VarExpert                                j2gort.PooledSlab[InlineVariables_VarExpert]
	slab_InlineVariables_InlineVarAnalysis                        j2gort.PooledSlab[InlineVariables_InlineVarAnalysis]
	slab_InlineVariables_InliningBehavior_AliasInlineRetryHandler j2gort.PooledSlab[InlineVariables_InliningBehavior_AliasInlineRetryHandler]
	inlinedCount                                                  int32
}

type InlineVariables_InliningBehavior struct {
	outer_                            *InlineVariables
	varToAliasRetryHandlersMap        j2gort.OrderedMap[*Var, *InlineVariables_InliningBehavior_AliasInlineRetryHandler]
	currentScopeHandledVarAnalysesMap j2gort.OrderedMap[*Var, *InlineVariables_InlineVarAnalysis]
}

func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_doInlinesForScope_3385(vars []*Var) {
	this_.currentScopeHandledVarAnalysesMap.Clear()
	for _, v := range vars {
		exp := this_.outer_.slab_InlineVariables_VarExpert.Alloc()
		exp.v = v
		if v.depth > 0 {
			// Simulate inner-scope alias waiting on outer var!
			rh := this_.outer_.slab_InlineVariables_InliningBehavior_AliasInlineRetryHandler.Alloc()
			rh.outer_ = this_
			rh.v = v
			rh.expert = exp
			this_.varToAliasRetryHandlersMap.Put(vars[0], rh)
			continue
		}
		an := this_.outer_.slab_InlineVariables_InlineVarAnalysis.Alloc()
		an.inlined = true
		this_.outer_.inlinedCount++
		this_.currentScopeHandledVarAnalysesMap.Put(v, an)
		if rh := this_.varToAliasRetryHandlersMap.Remove(v); rh != nil && rh.expert.v != nil {
			this_.outer_.inlinedCount++
		}
	}
}
`)

	allPasses := passes.AllPasses()
	jscompDir := filepath.Join(tmpDir, "jscomp")

	s1, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero Wave 36B sites on Run 1")
	}

	s2, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d byPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	readGot := func(name string) string {
		b, rErr := os.ReadFile(filepath.Join(jscompDir, name))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}

	cfgGot := readGot("cfg_dfa.go")
	if !strings.Contains(cfgGot, "this_.slab_LinearFlowState.Release()") || strings.Contains(cfgGot, "this_.slab_LinearFlowState.Reset()") {
		t.Fatalf("expected slab_LinearFlowState.Release() in ResetCfgForPool:\n%s", cfgGot)
	}

	mdnuGot := readGot("mdnu.go")
	for _, want := range []string{
		"this_.declarations = nil",
		"AcquireOrderedMapPtrCap[string, string](8)",
		"ReleaseOrderedMapPtr",
		"slab_ScopeContext",
		"slab_DeclaredVariableInfo",
		"slab_RenamingInfo",
		"slab_UndeclaredVariableInfo",
	} {
		if !strings.Contains(mdnuGot, want) {
			t.Fatalf("missing %q in lowered mdnu.go:\n%s", want, mdnuGot)
		}
	}

	ccrcGot := readGot("ccrc.go")
	for _, want := range []string{
		"slab_TopLevelStatement",
		"this_.slab_TopLevelStatement.Alloc()",
		"this_.slab_TopLevelStatement.Release()",
		"defer referenceCollector.ReleaseSlabs_()",
	} {
		if !strings.Contains(ccrcGot, want) {
			t.Fatalf("missing %q in lowered ccrc.go:\n%s", want, ccrcGot)
		}
	}

	optGot := readGot("opt_dae_iv.go")
	for _, want := range []string{
		"slab_CandidateAnalysis",
		"this_.slab_CandidateAnalysis.Reset()",
		"slab_ParameterList",
		"slab_Parameter",
		"slab_BailoutInformation",
		"this_.slab_BailoutInformation.Reset()",
		"this_.outer_.slab_InlineVariables_VarExpert.Reset()",
	} {
		if !strings.Contains(optGot, want) {
			t.Fatalf("missing %q in lowered opt_dae_iv.go:\n%s", want, optGot)
		}
	}

	mustWrite("jscomp/w36b_runtime_test.go", `package jscomp

import (
	"sync"
	"testing"

	"example.com/w36bslabs/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

func TestWave36BConcurrentRuntimeAndSlabs(t *testing.T) {
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for iter := 0; iter < 25; iter++ {
				// 1. CFG LinearFlowState Release on ResetCfgForPool
				cfg := &ControlFlowGraphBase[rhino.Node]{entry: 1}
				dfa := &DataFlowAnalysisBase[rhino.Node, int32]{cfg: cfg}
				cn := &cfgNode{}
				dfa.priv_DataFlowAnalysis_initialize_2310(cn, 10, 20)
				if st, ok := cn.GetAnnotation().(*DataFlowAnalysis_LinearFlowState); !ok || st.inL != int32(10) {
					t.Errorf("unexpected DFA state")
				}
				cn.SetAnnotation(nil)
				cfg.ResetCfgForPool()

				// 2. MakeDeclaredNamesUnique <=4 inline & >4 overflow scopes + String()
				rootRen := NewMakeDeclaredNamesUnique_ContextualRenamer()
				if cr0, ok := rootRen.(*MakeDeclaredNamesUnique_ContextualRenamerBase); !ok || cr0.String() != "decls=0" {
					t.Errorf("empty String() = %q", cr0.String())
				}
				mdnu := NewMakeDeclaredNamesUnique(rootRen)
				mdnu.EnterScope(&NodeTraversal{scopeRoot: 1, inGlobal: true, scope: &Scope{vars: []*Var{{name: "a"}, {name: "b"}, {name: "c"}, {name: "d"}, {name: "e"}, {name: "f"}}}})
				// Small child scope (<=4 vars -> 0 OrderedMap allocs!)
				mdnu.EnterScope(&NodeTraversal{scopeRoot: 2, inGlobal: false, scope: &Scope{vars: []*Var{{name: "a"}, {name: "b"}}}})
				if got := mdnu.renamerStack.Peek().GetReplacementName("a"); got != "a$jscomp$1" {
					t.Errorf("small scope replacement = %q", got)
				}
				mdnu.ExitScope(&NodeTraversal{scopeRoot: 2, inGlobal: false})
				// Overflow child scope (>4 vars -> pooled OrderedMapPtr!)
				mdnu.EnterScope(&NodeTraversal{scopeRoot: 3, inGlobal: false, scope: &Scope{vars: []*Var{{name: "a"}, {name: "b"}, {name: "c"}, {name: "d"}, {name: "e"}, {name: "f"}}}})
				if got := mdnu.renamerStack.Peek().GetReplacementName("f"); got != "f$jscomp$1" {
					t.Errorf("overflow scope replacement = %q", got)
				}
				mdnu.ExitScope(&NodeTraversal{scopeRoot: 3, inGlobal: false})
				mdnu.ReleaseSlabs_()

				// ContextualRenameInverter
				cri := NewMakeDeclaredNamesUnique_ContextualRenameInverter(nil)
				cri.Process_2(0, 1)
				if cri.invertedCount != 2 {
					t.Errorf("invertedCount = %d", cri.invertedCount)
				}

				// 3. CrossChunkCodeMotion + TopLevelStatement slab
				cccm := &CrossChunkCodeMotion{}
				cccm.Process_2(0, 5)
				if cccm.movedSum != (1+0)+(2+1)+(3+2)+(4+3)+(5+4) {
					t.Errorf("movedSum = %d", cccm.movedSum)
				}

				// 4. OptimizeParameters + DeadAssignmentsElimination + InlineVariables
				rm := &OptimizeCalls_ReferenceMap{
					names: j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList](),
					props: j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList](),
				}
				rm.names.Put("fn1", &OptimizeCalls_NodeList{nodes: []rhino.Node{2, 3, 4}})
				rm.props.Put("p1", &OptimizeCalls_NodeList{nodes: []rhino.Node{6, 7}})
				op := &OptimizeParameters{}
				op.Process_3(0, 1, rm)
				if op.optimized != 3 {
					t.Errorf("optimized = %d, want 3", op.optimized)
				}

				dae := &DeadAssignmentsElimination{}
				dae.Process_2(0, 1)
				if dae.eliminated != 1 {
					t.Errorf("eliminated = %d, want 1", dae.eliminated)
				}

				iv := &InlineVariables{}
				ib := &InlineVariables_InliningBehavior{
					outer_:                            iv,
					varToAliasRetryHandlersMap:        j2gort.NewOrderedMap[*Var, *InlineVariables_InliningBehavior_AliasInlineRetryHandler](),
					currentScopeHandledVarAnalysesMap: j2gort.NewOrderedMap[*Var, *InlineVariables_InlineVarAnalysis](),
				}
				outerV := &Var{name: "outer", depth: 0}
				innerAlias := &Var{name: "alias", depth: 1}
				// Inner scope exits first: registers retry handler on outerV (must NOT reset slab yet!)
				ib.priv_InliningBehavior_doInlinesForScope_3385([]*Var{outerV, innerAlias}[1:])
				ib.varToAliasRetryHandlersMap.Clear()
				ib.varToAliasRetryHandlersMap.Put(outerV, &InlineVariables_InliningBehavior_AliasInlineRetryHandler{
					outer_: ib,
					v:      innerAlias,
					expert: &InlineVariables_VarExpert{v: innerAlias},
				})
				// Outer scope exits: resolves retry handler -> map becomes empty -> resets slabs!
				ib.priv_InliningBehavior_doInlinesForScope_3385([]*Var{outerV})
				if iv.inlinedCount != 2 {
					t.Errorf("inlinedCount = %d, want 2", iv.inlinedCount)
				}
				iv.slab_InlineVariables_VarExpert.Release()
				iv.slab_InlineVariables_InlineVarAnalysis.Release()
				iv.slab_InlineVariables_InliningBehavior_AliasInlineRetryHandler.Release()
			}
		}()
	}
	wg.Wait()
}
`)

	cmd := exec.Command("go", "test", "-race", "./jscomp")
	cmd.Dir = tmpDir
	if out, runErr := cmd.CombinedOutput(); runErr != nil {
		t.Fatalf("go test -race ./jscomp failed: %v\n%s", runErr, string(out))
	}
}

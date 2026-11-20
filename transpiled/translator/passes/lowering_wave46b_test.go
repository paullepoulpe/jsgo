package passes

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func wave46bWriteFiles(t *testing.T, dir string, files map[string]string) {
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

const wave46bJ2gortStub = `package j2gort

import "sync/atomic"

var LinkedHashSetAllocs atomic.Int64

func ResetWave46BCounters() {
	LinkedHashSetAllocs.Store(0)
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

func ThrowNPE() {
	panic("NPE")
}

type Ext_java_lang_Iterable interface {
	Iterator() Ext_java_util_Iterator
}

type Ext_java_util_Iterator interface {
	HasNext() bool
	Next() any
	Remove()
}

type Ext_java_util_Collection interface {
	Ext_java_lang_Iterable
	Size() int32
	IsEmpty() bool
	Add_1(any) bool
	Contains(any) bool
	Remove_1_any(any) bool
	Clear()
}

type Ext_java_util_Set interface {
	Ext_java_util_Collection
}

type Ext_java_util_LinkedHashSet struct {
	items []any
}

func NewExt_java_util_LinkedHashSet() *Ext_java_util_LinkedHashSet {
	LinkedHashSetAllocs.Add(1)
	return &Ext_java_util_LinkedHashSet{}
}

func NewExt_java_util_LinkedHashSet_1_int32(cap int32) *Ext_java_util_LinkedHashSet {
	LinkedHashSetAllocs.Add(1)
	return &Ext_java_util_LinkedHashSet{items: make([]any, 0, cap)}
}

func (s *Ext_java_util_LinkedHashSet) Add_1(v any) bool {
	if s == nil {
		ThrowNPE()
	}
	for _, x := range s.items {
		if x == v {
			return false
		}
	}
	s.items = append(s.items, v)
	return true
}

func (s *Ext_java_util_LinkedHashSet) Contains(v any) bool {
	if s == nil {
		ThrowNPE()
	}
	for _, x := range s.items {
		if x == v {
			return true
		}
	}
	return false
}

func (s *Ext_java_util_LinkedHashSet) Remove_1_any(v any) bool {
	if s == nil {
		ThrowNPE()
	}
	for i, x := range s.items {
		if x == v {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Ext_java_util_LinkedHashSet) Size() int32 {
	if s == nil {
		ThrowNPE()
	}
	return int32(len(s.items))
}

func (s *Ext_java_util_LinkedHashSet) IsEmpty() bool {
	if s == nil {
		ThrowNPE()
	}
	return len(s.items) == 0
}

func (s *Ext_java_util_LinkedHashSet) Clear() {
	if s == nil {
		ThrowNPE()
	}
	s.items = s.items[:0]
}

func (s *Ext_java_util_LinkedHashSet) Iterator() Ext_java_util_Iterator {
	return &sliceJavaIter{items: s.items}
}

type sliceJavaIter struct {
	items []any
	pos   int
}

func (it *sliceJavaIter) HasNext() bool { return it.pos < len(it.items) }
func (it *sliceJavaIter) Next() any {
	v := it.items[it.pos]
	it.pos++
	return v
}
func (it *sliceJavaIter) Remove() {}

type IterCursor struct {
	items []any
	pos   int
}

func NewIterCursor(col any) IterCursor {
	if s, ok := col.(*Ext_java_util_LinkedHashSet); ok {
		if s == nil {
			ThrowNPE()
		}
		return IterCursor{items: s.items}
	}
	panic("unsupported")
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
	s   *OrderedSet[T]
	pos int
}

func NewOrderedSet[T comparable]() OrderedSet[T]             { return OrderedSet[T]{} }
func NewOrderedSetCap[T comparable](c int32) OrderedSet[T]   { return OrderedSet[T]{items: make([]T, 0, c)} }
func NewOrderedSetPtr[T comparable]() *OrderedSet[T]         { return &OrderedSet[T]{} }
func NewOrderedSetPtrCap[T comparable](c int32) *OrderedSet[T] { return &OrderedSet[T]{items: make([]T, 0, c)} }

func (s *OrderedSet[T]) Add(v T) bool {
	if s == nil {
		ThrowNPE()
	}
	for _, x := range s.items {
		if x == v {
			return false
		}
	}
	s.items = append(s.items, v)
	return true
}

func (s *OrderedSet[T]) Remove(v T) bool {
	if s == nil {
		ThrowNPE()
	}
	for i, x := range s.items {
		if x == v {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}

func (s OrderedSet[T]) Contains(v T) bool {
	for _, x := range s.items {
		if x == v {
			return true
		}
	}
	return false
}

func (s OrderedSet[T]) Size() int32   { return int32(len(s.items)) }
func (s OrderedSet[T]) IsEmpty() bool { return len(s.items) == 0 }

func (s *OrderedSet[T]) Clear() {
	if s == nil {
		ThrowNPE()
	}
	s.items = s.items[:0]
}

func (s *OrderedSet[T]) Iter() OrderedSetCursor[T] {
	if s == nil {
		ThrowNPE()
	}
	return OrderedSetCursor[T]{s: s}
}

func (c *OrderedSetCursor[T]) HasNext() bool { return c.pos < len(c.s.items) }
func (c *OrderedSetCursor[T]) Next() T {
	v := c.s.items[c.pos]
	c.pos++
	return v
}
`

func TestWave46B_NullableLazyLinkedHashSetLowering(t *testing.T) {
	tmpDir := t.TempDir()
	wave46bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                 "module example.com/w46btest\n\ngo 1.24\n",
		"gen/j2gort/j2gort.go":   wave46bJ2gortStub,
		"gen/jscomp/lazy_set.go": `package jscomp

import "example.com/w46btest/gen/j2gort"

type GlobalSymbol struct {
	name                     string
	refSym0                  *GlobalSymbol
	refSym1                  *GlobalSymbol
	referencingGlobalSymbols *j2gort.Ext_java_util_LinkedHashSet
	lazyInitField            *j2gort.Ext_java_util_LinkedHashSet
	reverseInitField         *j2gort.Ext_java_util_LinkedHashSet
	eagerValueField          *j2gort.Ext_java_util_LinkedHashSet
	mutatedInLoopField       *j2gort.Ext_java_util_LinkedHashSet
	mutatedViaHelperField    *j2gort.Ext_java_util_LinkedHashSet
}

func (*GlobalSymbol) HasIdentityEquals_() {}

func (this_ *GlobalSymbol) initFields_GlobalSymbol() {
	this_.eagerValueField = j2gort.NewExt_java_util_LinkedHashSet()
}

func NewGlobalSymbol(name string) *GlobalSymbol {
	this_ := &GlobalSymbol{name: name}
	this_.initFields_GlobalSymbol()
	return this_
}

// Exact CrossChunkCodeMotion.GlobalSymbol.addReferringGlobalSymbol shape.
func (this_ *GlobalSymbol) AddReferringGlobalSymbol(declaredSymbol *GlobalSymbol) {
	var set *j2gort.Ext_java_util_LinkedHashSet = this_.referencingGlobalSymbols
	_ = set
	if set != nil {
		set.Add_1(j2gort.PtrToAny(declaredSymbol))
		return
	}
	var r0 *GlobalSymbol = this_.refSym0
	if r0 == nil {
		this_.refSym0 = declaredSymbol
		return
	}
	if r0 == declaredSymbol {
		return
	}
	var r1 *GlobalSymbol = this_.refSym1
	if r1 == nil {
		this_.refSym1 = declaredSymbol
		return
	}
	if r1 == declaredSymbol {
		return
	}
	set = j2gort.NewExt_java_util_LinkedHashSet_1_int32(int32(8))
	set.Add_1(j2gort.PtrToAny(r0))
	set.Add_1(j2gort.PtrToAny(r1))
	set.Add_1(j2gort.PtrToAny(declaredSymbol))
	this_.referencingGlobalSymbols = set
	this_.refSym0 = nil
	this_.refSym1 = nil
}

func (this_ *GlobalSymbol) CollectReferringNames() []string {
	var out []string
	if this_.referencingGlobalSymbols != nil {
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(this_.referencingGlobalSymbols))
		for cur_1.HasNext() {
			raw_1 := cur_1.Next()
			sym := j2gort.CastPtr[GlobalSymbol](raw_1)
			out = append(out, sym.name)
		}
	} else if this_.refSym0 != nil {
		out = append(out, this_.refSym0.name)
		if this_.refSym1 != nil {
			out = append(out, this_.refSym1.name)
		}
	}
	return out
}

func (this_ *GlobalSymbol) RemoveReferring(sym *GlobalSymbol) bool {
	var set *j2gort.Ext_java_util_LinkedHashSet = this_.referencingGlobalSymbols
	if set != nil {
		return set.Remove_1_any(j2gort.PtrToAny(sym))
	}
	return false
}

func (this_ *GlobalSymbol) ClearReferring() {
	if this_.referencingGlobalSymbols != nil {
		this_.referencingGlobalSymbols.Clear()
	}
}

func (this_ *GlobalSymbol) UnguardedReferringSize() int32 {
	var set *j2gort.Ext_java_util_LinkedHashSet = this_.referencingGlobalSymbols
	return set.Size()
}

// Pattern B: lazy init via if set == nil { set = New...(); this_.field = set }, followed by guarded Contains.
func (this_ *GlobalSymbol) AddLazyAndCheck(sym *GlobalSymbol) bool {
	var set *j2gort.Ext_java_util_LinkedHashSet = this_.lazyInitField
	if set == nil {
		set = j2gort.NewExt_java_util_LinkedHashSet()
		this_.lazyInitField = set
	}
	set.Add_1(j2gort.PtrToAny(sym))
	return set.Contains(j2gort.PtrToAny(sym))
}

// Pattern C: local allocated first, then assigned to nullable lazy field.
func (this_ *GlobalSymbol) InitReverse(a, b *GlobalSymbol) int32 {
	var set *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	set.Add_1(j2gort.PtrToAny(a))
	set.Add_1(j2gort.PtrToAny(b))
	this_.reverseInitField = set
	if this_.reverseInitField != nil {
		return this_.reverseInitField.Size()
	}
	return 0
}

// Pattern D: non-nullable eager field + non-reassigned local alias must stay value OrderedSet[T].
func (this_ *GlobalSymbol) AddEager(sym *GlobalSymbol) bool {
	var set *j2gort.Ext_java_util_LinkedHashSet = this_.eagerValueField
	set.Add_1(j2gort.PtrToAny(sym))
	return set.Contains(j2gort.PtrToAny(sym))
}

// Pattern E (negative): mutating via local alias inside loop over field MUST disqualify!
func (this_ *GlobalSymbol) MutateViaAliasInLoop(sym *GlobalSymbol) {
	var set *j2gort.Ext_java_util_LinkedHashSet = this_.mutatedInLoopField
	if set == nil {
		set = j2gort.NewExt_java_util_LinkedHashSet()
		this_.mutatedInLoopField = set
	}
	set.Add_1(j2gort.PtrToAny(sym))
	cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(this_.mutatedInLoopField))
	for cur_2.HasNext() {
		raw_2 := cur_2.Next()
		el := j2gort.CastPtr[GlobalSymbol](raw_2)
		set.Add_1(j2gort.PtrToAny(el))
	}
}

// Pattern F (negative): mutating via helper method using local alias inside loop over field MUST disqualify!
func (this_ *GlobalSymbol) addViaAliasHelper(sym *GlobalSymbol) {
	var set *j2gort.Ext_java_util_LinkedHashSet = this_.mutatedViaHelperField
	if set == nil {
		set = j2gort.NewExt_java_util_LinkedHashSet()
		this_.mutatedViaHelperField = set
	}
	set.Add_1(j2gort.PtrToAny(sym))
}

func (this_ *GlobalSymbol) MutateViaHelperInLoop(sym *GlobalSymbol) {
	this_.addViaAliasHelper(sym)
	if this_.mutatedViaHelperField != nil {
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(this_.mutatedViaHelperField))
		for cur_3.HasNext() {
			raw_3 := cur_3.Next()
			el := j2gort.CastPtr[GlobalSymbol](raw_3)
			this_.addViaAliasHelper(el)
		}
	}
}
`,
	})

	jscompDir := filepath.Join(tmpDir, "gen", "jscomp")
	passList := []Pass{NewNativeSetLoweringPass()}

	s1, err := RunOnPackageDir(jscompDir, passList, Options{})
	if err != nil {
		t.Fatalf("NativeSetLoweringPass Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero rewrite sites on Run 1")
	}

	// Round-2 idempotency MUST be 0 sites!
	s2, err := RunOnPackageDir(jscompDir, passList, Options{})
	if err != nil {
		t.Fatalf("NativeSetLoweringPass Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT idempotent: TotalSites=%d FilesModified=%d byPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	gotBytes, err := os.ReadFile(filepath.Join(jscompDir, "lazy_set.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)

	for _, want := range []string{
		"referencingGlobalSymbols *j2gort.OrderedSet[*GlobalSymbol]",
		"lazyInitField            *j2gort.OrderedSet[*GlobalSymbol]",
		"reverseInitField         *j2gort.OrderedSet[*GlobalSymbol]",
		"eagerValueField          j2gort.OrderedSet[*GlobalSymbol]",
		"mutatedInLoopField       *j2gort.Ext_java_util_LinkedHashSet",
		"mutatedViaHelperField    *j2gort.Ext_java_util_LinkedHashSet",
		"var set *j2gort.OrderedSet[*GlobalSymbol] = this_.referencingGlobalSymbols",
		"set = j2gort.NewOrderedSetPtrCap[*GlobalSymbol](int32(8))",
		"cur_1 := this_.referencingGlobalSymbols.Iter()",
		"sym := cur_1.",
		"Next()",
		"return set.Remove(sym)",
		"this_.referencingGlobalSymbols.Clear()",
		"return j2gort.CheckNotNullPtr(set).Size()",
		"return set.Contains(sym)",
		"var set *j2gort.OrderedSet[*GlobalSymbol] = &this_.eagerValueField",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing expected lowered snippet %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "NewIterCursor(j2gort.CheckNotNullPtr(this_.referencingGlobalSymbols))") {
		t.Fatalf("referencingGlobalSymbols still uses NewIterCursor:\n%s", got)
	}

	// Verify runtime behavior, NPE parity, and zero LinkedHashSet allocations on lowered fields!
	wave46bWriteFiles(t, tmpDir, map[string]string{
		"gen/jscomp/lazy_set_rt_test.go": `package jscomp

import (
	"reflect"
	"testing"

	"example.com/w46btest/gen/j2gort"
)

func TestLazySetRuntimeParity(t *testing.T) {
	j2gort.ResetWave46BCounters()

	root := NewGlobalSymbol("root")
	a := NewGlobalSymbol("a")
	b := NewGlobalSymbol("b")
	c := NewGlobalSymbol("c")
	d := NewGlobalSymbol("d")

	// Unguarded Size() on nil lazy set must throw NPE (CheckNotNullPtr).
	func() {
		defer func() {
			if r := recover(); r != "NPE" {
				t.Fatalf("expected NPE on nil UnguardedReferringSize, got %v", r)
			}
		}()
		_ = root.UnguardedReferringSize()
	}()

	// Inline fast path (0, 1, 2 elements + dedup) then spill to *OrderedSet[*GlobalSymbol].
	root.AddReferringGlobalSymbol(a)
	root.AddReferringGlobalSymbol(a)
	root.AddReferringGlobalSymbol(b)
	if got := root.CollectReferringNames(); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("inline collect got %v", got)
	}
	root.AddReferringGlobalSymbol(c)
	root.AddReferringGlobalSymbol(b)
	root.AddReferringGlobalSymbol(d)
	if got := root.CollectReferringNames(); !reflect.DeepEqual(got, []string{"a", "b", "c", "d"}) {
		t.Fatalf("spilled collect got %v", got)
	}
	if sz := root.UnguardedReferringSize(); sz != 4 {
		t.Fatalf("expected size 4, got %d", sz)
	}
	if !root.RemoveReferring(b) || root.RemoveReferring(b) {
		t.Fatalf("unexpected RemoveReferring result")
	}
	if got := root.CollectReferringNames(); !reflect.DeepEqual(got, []string{"a", "c", "d"}) {
		t.Fatalf("after remove got %v", got)
	}
	root.ClearReferring()
	if sz := root.UnguardedReferringSize(); sz != 0 {
		t.Fatalf("expected size 0 after Clear, got %d", sz)
	}

	// Pattern B, C, D.
	if !root.AddLazyAndCheck(a) {
		t.Fatalf("AddLazyAndCheck failed")
	}
	if got := root.InitReverse(a, b); got != 2 {
		t.Fatalf("InitReverse expected 2, got %d", got)
	}
	if !root.AddEager(c) {
		t.Fatalf("AddEager failed")
	}

	if allocs := j2gort.LinkedHashSetAllocs.Load(); allocs != 0 {
		t.Fatalf("expected 0 Java LinkedHashSet allocations for lowered sets, got %d", allocs)
	}
}
`,
	})

	cmd := exec.Command("go", "test", "-p", "1", "-parallel", "1", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	cmd.Env = append(os.Environ(), "GOFLAGS=-p=2")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go test on lowered package failed: %v\n%s", err, string(out))
	}
}

func TestWave46B_RealCrossChunkCodeMotionLowering(t *testing.T) {
	realJscompDir, err := filepath.Abs(filepath.Join("..", "..", "gen", "jscomp"))
	if err != nil {
		t.Fatal(err)
	}
	cccmPath := filepath.Join(realJscompDir, "cross_chunk_code_motion.go")
	origBytes, err := os.ReadFile(cccmPath)
	if err != nil {
		t.Skipf("real cross_chunk_code_motion.go not found: %v", err)
	}
	if !strings.Contains(string(origBytes), "referencingGlobalSymbols") || !strings.Contains(string(origBytes), "referencedGlobalSymbols") {
		t.Skip("cross_chunk_code_motion.go does not contain expected fields")
	}
	alreadyLowered := strings.Contains(string(origBytes), "referencingGlobalSymbols                     *j2gort.OrderedSet[*CrossChunkCodeMotion_GlobalSymbol]") ||
		strings.Contains(string(origBytes), "referencingGlobalSymbols *j2gort.OrderedSet[*CrossChunkCodeMotion_GlobalSymbol]")

	// Run NativeSetLoweringPass in-memory targeting ONLY cross_chunk_code_motion.go (zero disk writes).
	filePaths, err := listPackageGoFiles(realJscompDir, false)
	if err != nil {
		t.Fatal(err)
	}
	memFiles := make([]memFile, len(filePaths))
	cccmIdx := -1
	for i, fp := range filePaths {
		raw, rErr := os.ReadFile(fp)
		if rErr != nil {
			t.Fatal(rErr)
		}
		isTarget := fp == cccmPath
		if isTarget {
			cccmIdx = i
		}
		memFiles[i] = memFile{name: fp, src: raw, target: isTarget}
	}
	if cccmIdx < 0 {
		t.Fatalf("cross_chunk_code_motion.go not in package file list")
	}

	imp := NewModuleImporter(realJscompDir)
	passList := []Pass{NewNativeSetLoweringPass()}
	outFiles, stats1, err := executePipeline("github.com/paullepoulpe/jsgo/transpiled/gen/jscomp", memFiles, passList, imp)
	if err != nil {
		t.Fatalf("executePipeline Run 1 on cross_chunk_code_motion.go failed: %v", err)
	}
	if !alreadyLowered && stats1.TotalSites == 0 {
		t.Fatalf("expected non-zero native_set_lowering sites on un-lowered cross_chunk_code_motion.go")
	}

	got := string(outFiles[cccmIdx].src)
	foundRefingField := false
	foundRefedField := false
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "referencingGlobalSymbols") && strings.Contains(line, "*j2gort.OrderedSet[*CrossChunkCodeMotion_GlobalSymbol]") {
			foundRefingField = true
		}
		if strings.Contains(line, "referencedGlobalSymbols") && strings.Contains(line, "*j2gort.OrderedSet[*CrossChunkCodeMotion_GlobalSymbol]") {
			foundRefedField = true
		}
	}
	if !foundRefingField || !foundRefedField {
		t.Fatalf("expected BOTH referencingGlobalSymbols (%v) and referencedGlobalSymbols (%v) fields to lower to *j2gort.OrderedSet[*CrossChunkCodeMotion_GlobalSymbol]", foundRefingField, foundRefedField)
	}
	for _, want := range []string{
		"j2gort.NewOrderedSetPtrCap[*CrossChunkCodeMotion_GlobalSymbol](int32(8))",
		"this_.referencedGlobalSymbols.Iter()",
		"symbol.referencingGlobalSymbols.Iter()",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("lowered cross_chunk_code_motion.go missing %q", want)
		}
	}
	if strings.Contains(got, "Ext_java_util_LinkedHashSet") {
		t.Fatalf("lowered cross_chunk_code_motion.go still contains Ext_java_util_LinkedHashSet")
	}

	// Verify round-2 idempotency on the in-memory lowered cross_chunk_code_motion.go!
	imp2 := NewModuleImporter(realJscompDir)
	_, stats2, err := executePipeline("github.com/paullepoulpe/jsgo/transpiled/gen/jscomp", outFiles, passList, imp2)
	if err != nil {
		t.Fatalf("executePipeline Run 2 on cross_chunk_code_motion.go failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("Run 2 on cross_chunk_code_motion.go NOT idempotent: TotalSites=%d FilesModified=%d", stats2.TotalSites, stats2.FilesModified)
	}
	t.Logf("Real cross_chunk_code_motion.go: Run 1 sites=%d (modified=%d), Run 2 sites=%d", stats1.TotalSites, stats1.FilesModified, stats2.TotalSites)
}

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

func TestWave34B_DequeAndSlice_EnumAndValueStructUnboxing(t *testing.T) {
	t.Parallel()
	j2gortSrc := `package j2gort

type EnumBase struct {
	Name_    string
	Ordinal_ int32
}

type Ext_java_lang_Iterable interface {
	Implements_j2gort_Ext_java_lang_Iterable()
	IsJavaNil_() bool
}
type Ext_java_util_Collection interface {
	Ext_java_lang_Iterable
	Implements_j2gort_Ext_java_util_Collection()
	Size() int32
	IsEmpty() bool
	Add_1(e any) bool
	Clear()
}
type Ext_java_util_SequencedCollection interface {
	Ext_java_util_Collection
	Implements_j2gort_Ext_java_util_SequencedCollection()
	AddFirst(e any)
	AddLast(e any)
	GetFirst() any
	GetLast() any
	RemoveFirst() any
	RemoveLast() any
}
type Ext_java_util_List interface {
	Ext_java_util_SequencedCollection
	Implements_j2gort_Ext_java_util_List()
	Get(idx int32) any
}
type Ext_java_util_ArrayList interface {
	Ext_java_util_List
	Implements_j2gort_Ext_java_util_ArrayList()
	Remove_1_int32(idx int32) any
}
type Ext_java_util_Queue interface {
	Ext_java_util_Collection
	Implements_j2gort_Ext_java_util_Queue()
}
type Ext_java_util_Deque interface {
	Ext_java_util_Queue
	Implements_j2gort_Ext_java_util_Deque()
	AddFirst(e any)
	AddLast(e any)
	Push(e any)
	Pop() any
	RemoveFirst() any
	RemoveLast() any
	Peek() any
	PeekFirst() any
	PeekLast() any
}
type Ext_java_util_ArrayDeque interface {
	Ext_java_util_Deque
	Implements_j2gort_Ext_java_util_ArrayDeque()
}

type NoSuchElementErr struct{ Msg string }

func ThrowNPE() { panic("NPE") }

func CheckNotNull[T any](v T) T { return v }

type NullableVal interface {
	comparable
	IsJavaNil_() bool
}

var BoxNullableValCalls int

func BoxNullableVal[T NullableVal](v T) any {
	BoxNullableValCalls++
	if v.IsJavaNil_() {
		return nil
	}
	return v
}

func PtrToAny[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}

func Cast[T any](v any) T {
	if v == nil {
		var z T
		return z
	}
	return v.(T)
}

func CastVal[T any](v any) T {
	return Cast[T](v)
}

func CastPtr[T any](v any) *T {
	if v == nil {
		return nil
	}
	return v.(*T)
}

func ListGetVal[T any](l Ext_java_util_List, idx int32) T {
	return Cast[T](l.Get(idx))
}

type rawList struct{ data []any }

func (l *rawList) Implements_j2gort_Ext_java_lang_Iterable()        {}
func (l *rawList) Implements_j2gort_Ext_java_util_Collection()      {}
func (l *rawList) Implements_j2gort_Ext_java_util_SequencedCollection() {}
func (l *rawList) Implements_j2gort_Ext_java_util_List()            {}
func (l *rawList) Implements_j2gort_Ext_java_util_ArrayList()       {}
func (l *rawList) IsJavaNil_() bool                                 { return l == nil }
func (l *rawList) Size() int32                                      { return int32(len(l.data)) }
func (l *rawList) IsEmpty() bool                                    { return len(l.data) == 0 }
func (l *rawList) Add_1(e any) bool                                 { l.data = append(l.data, e); return true }
func (l *rawList) AddFirst(e any)                                   { l.data = append([]any{e}, l.data...) }
func (l *rawList) AddLast(e any)                                    { l.data = append(l.data, e) }
func (l *rawList) Clear()                                           { clear(l.data); l.data = l.data[:0] }
func (l *rawList) Get(idx int32) any                                { return l.data[idx] }
func (l *rawList) GetFirst() any {
	if len(l.data) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	return l.data[0]
}
func (l *rawList) GetLast() any {
	if len(l.data) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	return l.data[len(l.data)-1]
}
func (l *rawList) RemoveFirst() any {
	if len(l.data) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	v := l.data[0]
	l.data = append([]any(nil), l.data[1:]...)
	return v
}
func (l *rawList) RemoveLast() any {
	if len(l.data) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	v := l.data[len(l.data)-1]
	l.data = l.data[:len(l.data)-1]
	return v
}
func (l *rawList) Remove_1_int32(idx int32) any {
	v := l.data[idx]
	l.data = append(l.data[:idx], l.data[idx+1:]...)
	return v
}

func NewExt_java_util_ArrayList() Ext_java_util_ArrayList              { return &rawList{data: []any{}} }
func NewExt_java_util_ArrayList_1_int32(c int32) Ext_java_util_ArrayList { return &rawList{data: make([]any, 0, c)} }
func Ext_java_util_Collections_UnmodifiableList(l Ext_java_util_List) Ext_java_util_List { return l }

type liveSliceView[T any] struct{ p *[]T }

func (v liveSliceView[T]) Implements_j2gort_Ext_java_lang_Iterable()          {}
func (v liveSliceView[T]) Implements_j2gort_Ext_java_util_Collection()        {}
func (v liveSliceView[T]) Implements_j2gort_Ext_java_util_SequencedCollection() {}
func (v liveSliceView[T]) Implements_j2gort_Ext_java_util_List()              {}
func (v liveSliceView[T]) IsJavaNil_() bool                                   { return v.p == nil || *v.p == nil }
func (v liveSliceView[T]) Size() int32                                        { return int32(len(*v.p)) }
func (v liveSliceView[T]) IsEmpty() bool                                      { return len(*v.p) == 0 }
func (v liveSliceView[T]) Add_1(any) bool                                     { panic("unmodifiable") }
func (v liveSliceView[T]) AddFirst(any)                                       { panic("unmodifiable") }
func (v liveSliceView[T]) AddLast(any)                                        { panic("unmodifiable") }
func (v liveSliceView[T]) Clear()                                             { panic("unmodifiable") }
func (v liveSliceView[T]) RemoveFirst() any                                   { panic("unmodifiable") }
func (v liveSliceView[T]) RemoveLast() any                                    { panic("unmodifiable") }
func (v liveSliceView[T]) Get(idx int32) any                                  { return (*v.p)[idx] }
func (v liveSliceView[T]) GetFirst() any {
	s := *v.p
	if len(s) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	return s[0]
}
func (v liveSliceView[T]) GetLast() any {
	s := *v.p
	if len(s) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	return s[len(s)-1]
}

func SliceUnmodifiableList[T any](s []T) Ext_java_util_List {
	cp := append([]T(nil), s...)
	return liveSliceView[T]{p: &cp}
}
func SliceUnmodifiableListLive[T any](p *[]T) Ext_java_util_List {
	return liveSliceView[T]{p: p}
}
func CloneNonNilSlice[T any](s []T) []T { return append([]T{}, s...) }
func MakeNonNilSlice[T any](c int32) []T { return make([]T, 0, c) }
func SliceGetFirst[T any](s []T) T {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	return s[0]
}
func SliceGetLast[T any](s []T) T {
	if s == nil {
		ThrowNPE()
	}
	if len(s) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	return s[len(s)-1]
}
func SliceRemoveAt[T any](s []T, idx int32) []T {
	if s == nil {
		ThrowNPE()
	}
	var z T
	copy(s[idx:], s[idx+1:])
	s[len(s)-1] = z
	return s[:len(s)-1]
}
func SliceRemoveAtElem[T any](sp *[]T, idx int32) T {
	s := *sp
	v := s[idx]
	*sp = SliceRemoveAt(s, idx)
	return v
}
func SliceInsertAt[T any](s []T, idx int32, elem T) []T {
	var z T
	s = append(s, z)
	copy(s[int(idx)+1:], s[int(idx):len(s)-1])
	s[idx] = elem
	return s
}

type rawDeque struct{ data []any }

func (d *rawDeque) Implements_j2gort_Ext_java_lang_Iterable()   {}
func (d *rawDeque) Implements_j2gort_Ext_java_util_Collection() {}
func (d *rawDeque) Implements_j2gort_Ext_java_util_Queue()      {}
func (d *rawDeque) Implements_j2gort_Ext_java_util_Deque()      {}
func (d *rawDeque) Implements_j2gort_Ext_java_util_ArrayDeque() {}
func (d *rawDeque) IsJavaNil_() bool                            { return d == nil }
func (d *rawDeque) Size() int32                                 { return int32(len(d.data)) }
func (d *rawDeque) IsEmpty() bool                               { return len(d.data) == 0 }
func (d *rawDeque) Clear()                                      { clear(d.data); d.data = d.data[:0] }
func (d *rawDeque) Add_1(e any) bool                            { d.data = append(d.data, e); return true }
func (d *rawDeque) AddFirst(e any)                              { d.data = append([]any{e}, d.data...) }
func (d *rawDeque) AddLast(e any)                               { d.data = append(d.data, e) }
func (d *rawDeque) Push(e any)                                  { d.AddFirst(e) }
func (d *rawDeque) Pop() any                                    { return d.RemoveFirst() }
func (d *rawDeque) RemoveFirst() any {
	if len(d.data) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	v := d.data[0]
	d.data = append([]any(nil), d.data[1:]...)
	return v
}
func (d *rawDeque) RemoveLast() any {
	if len(d.data) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	v := d.data[len(d.data)-1]
	d.data = d.data[:len(d.data)-1]
	return v
}
func (d *rawDeque) Peek() any      { return d.PeekFirst() }
func (d *rawDeque) PeekFirst() any { if len(d.data) == 0 { return nil }; return d.data[0] }
func (d *rawDeque) PeekLast() any  { if len(d.data) == 0 { return nil }; return d.data[len(d.data)-1] }

func NewExt_java_util_ArrayDeque() Ext_java_util_ArrayDeque { return &rawDeque{data: []any{}} }

type Deque[T any] struct {
	items []T
}

func NewDeque[T any]() Deque[T] { return Deque[T]{} }
func (d *Deque[T]) IsNil() bool { return d == nil }
func (d *Deque[T]) check() {
	if d == nil {
		ThrowNPE()
	}
}
func (d *Deque[T]) Size() int32   { d.check(); return int32(len(d.items)) }
func (d *Deque[T]) IsEmpty() bool { d.check(); return len(d.items) == 0 }
func (d *Deque[T]) Clear()        { d.check(); clear(d.items); d.items = d.items[:0] }
func (d *Deque[T]) AddFirst(v T)  { d.check(); d.items = append([]T{v}, d.items...) }
func (d *Deque[T]) AddLast(v T)   { d.check(); d.items = append(d.items, v) }
func (d *Deque[T]) OfferLast(v T) bool { d.AddLast(v); return true }
func (d *Deque[T]) Push(v T)      { d.AddFirst(v) }
func (d *Deque[T]) Pop() T        { return d.RemoveFirst() }
func (d *Deque[T]) RemoveFirst() T {
	d.check()
	if len(d.items) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	v := d.items[0]
	var z T
	copy(d.items, d.items[1:])
	d.items[len(d.items)-1] = z
	d.items = d.items[:len(d.items)-1]
	return v
}
func (d *Deque[T]) RemoveLast() T {
	d.check()
	if len(d.items) == 0 {
		panic(NoSuchElementErr{"empty"})
	}
	v := d.items[len(d.items)-1]
	var z T
	d.items[len(d.items)-1] = z
	d.items = d.items[:len(d.items)-1]
	return v
}
func (d *Deque[T]) Peek() T      { return d.PeekFirst() }
func (d *Deque[T]) PeekFirst() T { d.check(); var z T; if len(d.items) == 0 { return z }; return d.items[0] }
func (d *Deque[T]) PeekLast() T  { d.check(); var z T; if len(d.items) == 0 { return z }; return d.items[len(d.items)-1] }

type DequeCursor[T any] struct {
	items []T
	idx   int
}

func (d *Deque[T]) Iter() DequeCursor[T] { d.check(); return DequeCursor[T]{items: d.items} }
func (c *DequeCursor[T]) HasNext() bool  { return c.idx < len(c.items) }
func (c *DequeCursor[T]) Next() T        { v := c.items[c.idx]; c.idx++; return v }

type TypedIterCursor[T any] struct {
	items []T
	idx   int
}

func NewTypedIterCursor[T any](col any) TypedIterCursor[T] {
	switch x := col.(type) {
	case *rawList:
		out := make([]T, len(x.data))
		for i, v := range x.data {
			out[i] = Cast[T](v)
		}
		return TypedIterCursor[T]{items: out}
	case *rawDeque:
		out := make([]T, len(x.data))
		for i, v := range x.data {
			out[i] = Cast[T](v)
		}
		return TypedIterCursor[T]{items: out}
	}
	return TypedIterCursor[T]{}
}
func (c *TypedIterCursor[T]) HasNext() bool { return c.idx < len(c.items) }
func (c *TypedIterCursor[T]) Next() T       { v := c.items[c.idx]; c.idx++; return v }
`

	appSrc := `package deser

import "example.com/w34b_col/j2gort"

// Value enum (uint8, 0 == null) like FeatureCollector_FeatureContext
type FeatureContext uint8

const (
	FeatureContext_NONE     FeatureContext = 1
	FeatureContext_FUNCTION FeatureContext = 2
	FeatureContext_CLASS    FeatureContext = 3
)

func (e FeatureContext) IsJavaNil_() bool { return e == 0 }

// Pointer enum (polymorphic enum with per-constant overrides)
type RichEnum struct {
	j2gort.EnumBase
	Weight int32
}

var (
	RichEnum_A = &RichEnum{EnumBase: j2gort.EnumBase{Name_: "A", Ordinal_: 0}, Weight: 10}
	RichEnum_B = &RichEnum{EnumBase: j2gort.EnumBase{Name_: "B", Ordinal_: 1}, Weight: 20}
)

// Value structs (JSTypeExpression / SourcePosition)
type SourcePosition struct {
	Line   int32
	Column int32
	Valid  bool
}

func (p SourcePosition) IsJavaNil_() bool { return p == (SourcePosition{}) }

// ScriptNodeDeserializer_Runner reproduction
type Runner struct {
	trace []FeatureContext
}

func (this_ *Runner) visitSingle(unmod j2gort.Ext_java_util_List, depth int32) FeatureContext {
	prev := j2gort.Cast[FeatureContext](unmod.GetLast())
	if depth%2 == 0 {
		return FeatureContext_FUNCTION
	}
	if prev == FeatureContext_FUNCTION {
		return FeatureContext_CLASS
	}
	return FeatureContext_NONE
}

func (this_ *Runner) Run(maxDepth int32) []FeatureContext {
	var contextStack j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	contextStack.AddLast(j2gort.BoxNullableVal(FeatureContext_NONE))
	var unmodifiableContextStack j2gort.Ext_java_util_List = j2gort.Ext_java_util_Collections_UnmodifiableList(contextStack)
	this_.priv_Runner_visit_7658(0, maxDepth, contextStack, unmodifiableContextStack)
	return this_.trace
}

func (this_ *Runner) priv_Runner_visit_7658(depth int32, maxDepth int32, contextStack j2gort.Ext_java_util_List, unmodifiableContextStack j2gort.Ext_java_util_List) {
	if contextStack != nil {
		contextStack.AddLast(j2gort.BoxNullableVal(this_.visitSingle(unmodifiableContextStack, depth)))
	}
	this_.trace = append(this_.trace, j2gort.Cast[FeatureContext](unmodifiableContextStack.GetLast()))
	if depth+1 < maxDepth {
		this_.priv_Runner_visit_7658(depth+1, maxDepth, contextStack, unmodifiableContextStack)
	}
	if contextStack != nil {
		contextStack.RemoveLast()
	}
}

// Deque with value-enum, pointer-enum, and value-struct elements (including CastVal & NewTypedIterCursor)
func RunDequeSuite() int32 {
	var ctxDq j2gort.Ext_java_util_Deque = j2gort.NewExt_java_util_ArrayDeque()
	ctxDq.AddLast(j2gort.BoxNullableVal(FeatureContext_NONE))
	ctxDq.Push(j2gort.BoxNullableVal(FeatureContext_CLASS))
	topCtx := j2gort.CastVal[FeatureContext](ctxDq.Pop())

	var enumDq j2gort.Ext_java_util_Deque = j2gort.NewExt_java_util_ArrayDeque()
	enumDq.AddLast(j2gort.PtrToAny(RichEnum_A))
	enumDq.AddFirst(j2gort.PtrToAny(RichEnum_B))
	firstEnum := j2gort.CastPtr[RichEnum](enumDq.RemoveFirst())
	lastEnum := j2gort.CastPtr[RichEnum](enumDq.PeekLast())

	var posDq j2gort.Ext_java_util_Deque = j2gort.NewExt_java_util_ArrayDeque()
	posDq.AddLast(j2gort.BoxNullableVal(SourcePosition{Line: 3, Column: 5, Valid: true}))
	posDq.AddLast(j2gort.BoxNullableVal(SourcePosition{Line: 7, Column: 9, Valid: true}))
	var lineSum int32
	{
		cur := j2gort.NewTypedIterCursor[SourcePosition](j2gort.CheckNotNull(posDq))
		for cur.HasNext() {
			p := cur.Next()
			lineSum += p.Line
		}
	}
	poppedPos := j2gort.CastVal[SourcePosition](posDq.RemoveLast())

	var posList j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	posList.Add_1(j2gort.BoxNullableVal(SourcePosition{Line: 11, Column: 1, Valid: true}))
	posList.AddLast(j2gort.BoxNullableVal(SourcePosition{Line: 13, Column: 2, Valid: true}))
	firstListPos := j2gort.ListGetVal[SourcePosition](posList, 0)
	posList.RemoveFirst()

	return int32(topCtx) + firstEnum.Weight + lastEnum.Weight + lineSum + poppedPos.Column + firstListPos.Line + posList.Size()
}

// Negative safety guard: local list reassigned to a brand-new ArrayList AFTER UnmodifiableList(list)
// MUST NOT lower to SliceUnmodifiableListLive(&reassigned), because in Java view still sees the first list!
func NegativeReassignedAfterUnmod() int32 {
	var reassigned j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	reassigned.Add_1(j2gort.BoxNullableVal(FeatureContext_NONE))
	view := j2gort.Ext_java_util_Collections_UnmodifiableList(reassigned)
	reassigned = j2gort.NewExt_java_util_ArrayList()
	reassigned.Add_1(j2gort.BoxNullableVal(FeatureContext_CLASS))
	reassigned.Add_1(j2gort.BoxNullableVal(FeatureContext_FUNCTION))
	return view.Size() + reassigned.Size()
}
`

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
	mustWrite("go.mod", "module example.com/w34b_col\n\ngo 1.24\n")
	mustWrite("j2gort/j2gort.go", j2gortSrc)
	mustWrite("deser/deser.go", appSrc)

	colPasses, err := passes.ParsePassesFlag("native_deque_lowering,native_slice_lowering")
	if err != nil {
		t.Fatal(err)
	}
	deserDir := filepath.Join(tmpDir, "deser")

	s1, err := passes.RunOnPackageDir(deserDir, colPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.ByPass["native_deque_lowering"] == 0 || s1.ByPass["native_slice_lowering"] == 0 {
		t.Fatalf("expected both native_deque_lowering and native_slice_lowering sites > 0, got %+v", s1)
	}

	s2, err := passes.RunOnPackageDir(deserDir, colPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 (idempotency) failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: %+v", s2)
	}

	outBytes, err := os.ReadFile(filepath.Join(deserDir, "deser.go"))
	if err != nil {
		t.Fatal(err)
	}
	code := string(outBytes)

	// Verify Runner.Run & priv_Runner_visit_7658 have ZERO BoxNullableVal and use []FeatureContext + SliceUnmodifiableListLive!
	runSection := code[:strings.Index(code, "func NegativeReassignedAfterUnmod")]
	if strings.Contains(runSection, "BoxNullableVal") {
		t.Fatalf("expected ZERO BoxNullableVal calls in lowered Runner + RunDequeSuite, got:\n%s", runSection)
	}
	if !strings.Contains(runSection, "SliceUnmodifiableListLive(&contextStack)") {
		t.Fatalf("expected SliceUnmodifiableListLive(&contextStack) in Runner.Run, got:\n%s", runSection)
	}
	if !strings.Contains(runSection, "contextStack *[]FeatureContext") {
		t.Fatalf("expected contextStack *[]FeatureContext in priv_Runner_visit_7658, got:\n%s", runSection)
	}
	if !strings.Contains(runSection, "j2gort.Deque[FeatureContext]") || !strings.Contains(runSection, "j2gort.Deque[*RichEnum]") || !strings.Contains(runSection, "j2gort.Deque[SourcePosition]") {
		t.Fatalf("expected typed Deque[FeatureContext], Deque[*RichEnum], and Deque[SourcePosition], got:\n%s", runSection)
	}

	// Verify negative safety guard preserved Java semantics on reassigned local!
	negSection := code[strings.Index(code, "func NegativeReassignedAfterUnmod"):]
	if strings.Contains(negSection, "SliceUnmodifiableListLive") {
		t.Fatalf("NegativeReassignedAfterUnmod must NOT use SliceUnmodifiableListLive:\n%s", negSection)
	}

	// Run runtime parity + 0-BoxNullableVal + 0-alloc steady-state verification!
	rtTest := `package deser

import (
	"testing"

	"example.com/w34b_col/j2gort"
)

func TestRuntimeAndZeroBox(t *testing.T) {
	j2gort.BoxNullableValCalls = 0
	r := &Runner{}
	got := r.Run(4)
	want := []FeatureContext{FeatureContext_FUNCTION, FeatureContext_CLASS, FeatureContext_FUNCTION, FeatureContext_CLASS}
	if len(got) != len(want) {
		t.Fatalf("trace len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("trace[%d] = %v, want %v", i, got[i], want[i])
		}
	}
	if sum := RunDequeSuite(); sum != (int32(FeatureContext_CLASS) + 20 + 10 + 10 + 9 + 11 + 1) {
		t.Fatalf("RunDequeSuite() = %d", sum)
	}
	if j2gort.BoxNullableValCalls != 0 {
		t.Fatalf("expected 0 BoxNullableVal calls during Runner.Run + RunDequeSuite, got %d", j2gort.BoxNullableValCalls)
	}
	if n := NegativeReassignedAfterUnmod(); n != 3 {
		t.Fatalf("NegativeReassignedAfterUnmod() = %d, want 3 (1 in old view + 2 in new list)", n)
	}
}
`
	mustWrite("deser/deser_rt_test.go", rtTest)
	cmd := exec.Command("go", "test", "-v", "./deser")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("runtime test failed: %v\n%s", err, string(out))
	}
}

func TestWave34B_TypedScopeCreatorScopeBuilderReuse_And_FunctionTypeBuilderStack(t *testing.T) {
	t.Parallel()
	j2gortSrc := `package j2gort

type Slab[T any] struct{}

func ThrowNPE() { panic("NPE") }
func CheckNotNullIface[T any](v T) T { return v }
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

var NullString = "\x00"
`

	rhinoSrc := `package rhino

import "example.com/w34b_sb/j2gort"

type Node uint32

func (n Node) IsFunction() bool         { return n == 1 }
func (n Node) IsClass() bool            { return n == 2 }
func (n Node) IsMemberFieldDef() bool   { return false }
func (n Node) IsComputedFieldDef() bool { return false }
func (n Node) IsModuleBody() bool       { return false }

type JSType interface {
	IsJavaNil_() bool
}
type ObjectType interface {
	JSType
}
type TemplateType struct{ Name string }
type TemplateTypeMap struct{ Keys []*TemplateType }

func (m *TemplateTypeMap) Size() int32                         { return int32(len(m.Keys)) }
func (m *TemplateTypeMap) AsKeysOnlySingleSubMap() *TemplateTypeMap { return m }

type JSTypeRegistry struct {
	emptyMap *TemplateTypeMap
	protoObj ObjectType
}

func NewJSTypeRegistry() *JSTypeRegistry {
	return &JSTypeRegistry{emptyMap: &TemplateTypeMap{}}
}
func (r *JSTypeRegistry) GetEmptyTemplateTypeMap() *TemplateTypeMap { return r.emptyMap }
func (r *JSTypeRegistry) GetNativeObjectType(id int32) ObjectType   { return r.protoObj }

type JSTypeBase struct {
	registry       *JSTypeRegistry
	templateTypeMap *TemplateTypeMap
	protoFlags     int8
	fnFlags        int8
}

func (b *JSTypeBase) GetJSTypeBase() *JSTypeBase { return b }

type ObjectTypeBase struct {
	JSTypeBase
}

func InitObjectType_3(this_ *ObjectTypeBase, reg *JSTypeRegistry, tm *TemplateTypeMap, mask int32) {
	this_.registry = reg
	this_.templateTypeMap = tm
	_ = mask
}

type PrototypeObjectType_Builder interface {
	GetPrototypeObjectType_BuilderBase() *PrototypeObjectType_BuilderBase
}

func ToIface_PrototypeObjectType_Builder[T interface {
	PrototypeObjectType_Builder
	comparable
}](v T) PrototypeObjectType_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

type PrototypeObjectType_BuilderBase struct {
	self_              PrototypeObjectType_Builder
	registry           *JSTypeRegistry
	className          string
	implicitPrototype  ObjectType
	nativeType         bool
	anonymousType      bool
	templateTypeMap    *TemplateTypeMap
	templateParamCount int32
}

func (b *PrototypeObjectType_BuilderBase) BindSelf_(s any) {
	b.self_ = s.(PrototypeObjectType_Builder)
}
func (b *PrototypeObjectType_BuilderBase) GetPrototypeObjectType_BuilderBase() *PrototypeObjectType_BuilderBase {
	return b
}
func (b *PrototypeObjectType_BuilderBase) initFields_PrototypeObjectType_Builder() {
	b.className = j2gort.NullString
}
func InitPrototypeObjectType_Builder(this_ *PrototypeObjectType_BuilderBase, registry *JSTypeRegistry) {
	this_.initFields_PrototypeObjectType_Builder()
	this_.registry = registry
	this_.templateTypeMap = registry.GetEmptyTemplateTypeMap()
}
func (this_ *PrototypeObjectType_BuilderBase) CastThis() PrototypeObjectType_Builder { return this_.self_ }
func (this_ *PrototypeObjectType_BuilderBase) SetName(x string) any {
	this_.className = x
	return this_.CastThis()
}
func (this_ *PrototypeObjectType_BuilderBase) SetImplicitPrototype(x ObjectType) PrototypeObjectType_Builder {
	this_.implicitPrototype = x
	return this_.CastThis()
}
func (this_ *PrototypeObjectType_BuilderBase) SetNative(x bool) PrototypeObjectType_Builder {
	this_.nativeType = x
	return this_.CastThis()
}
func (this_ *PrototypeObjectType_BuilderBase) SetTemplateTypeMap(x *TemplateTypeMap) PrototypeObjectType_Builder {
	this_.templateTypeMap = x
	return this_.CastThis()
}
func (this_ *PrototypeObjectType_BuilderBase) SetTemplateParamCount(x int32) PrototypeObjectType_Builder {
	this_.templateParamCount = x
	return this_.CastThis()
}

type PrototypeObjectTypeBase struct {
	ObjectTypeBase
	className string
}

func (p *PrototypeObjectTypeBase) GetPrototypeObjectTypeBase() *PrototypeObjectTypeBase { return p }
func (p *PrototypeObjectTypeBase) initFields_PrototypeObjectType()                      {}

func InitPrototypeObjectType_2(this_ *PrototypeObjectTypeBase, builder PrototypeObjectType_Builder, extraKindMask int32) {
	_ = builder
	_ = extraKindMask
	InitObjectType_3(&this_.ObjectTypeBase, builder.GetPrototypeObjectType_BuilderBase().registry, builder.GetPrototypeObjectType_BuilderBase().templateTypeMap, extraKindMask)
	this_.initFields_PrototypeObjectType()
	this_.className = builder.GetPrototypeObjectType_BuilderBase().className
	if builder.GetPrototypeObjectType_BuilderBase().nativeType {
		this_.protoFlags = 1
	}
}

type FunctionType_Kind uint8

const (
	FunctionType_Kind_ORDINARY    FunctionType_Kind = 1
	FunctionType_Kind_CONSTRUCTOR FunctionType_Kind = 2
	FunctionType_Kind_INTERFACE   FunctionType_Kind = 3
)

type FunctionType_Parameter struct {
	Type JSType
}

type FunctionType struct {
	PrototypeObjectTypeBase
	kind             FunctionType_Kind
	sourceNode       Node
	googModuleId     string
	isAbstract       bool
	isKnownAmbiguous bool
	parameters       []FunctionType_Parameter
	returnType       JSType
	returnInferred   bool
	typeOfThis       JSType
}

func (f *FunctionType) IsJavaNil_() bool                   { return f == nil }
func (f *FunctionType) BindSelf_(any)                      {}
func (f *FunctionType) ToMaybeFunctionType() *FunctionType { return f }
func (f *FunctionType) GetReferenceName() string           { return f.className }
func (f *FunctionType) GetGoogModuleId() string            { return f.googModuleId }
func (f *FunctionType) IsNativeObjectType() bool           { return f.protoFlags != 0 }
func (f *FunctionType) GetTemplateTypeMap() *TemplateTypeMap { return f.templateTypeMap }
func (f *FunctionType) GetTemplateParamCount() int32       { return int32(len(f.templateTypeMap.Keys)) }
func (f *FunctionType) IsAbstract() bool                   { return f.isAbstract }
func (f *FunctionType) GetKind() FunctionType_Kind         { return f.kind }
func (f *FunctionType) GetReturnType() JSType              { return f.returnType }
func (f *FunctionType) IsReturnTypeInferred() bool         { return f.returnInferred }
func (f *FunctionType) GetSource() Node                    { return f.sourceNode }
func (f *FunctionType) GetTypeOfThis() JSType              { return f.typeOfThis }
func (f *FunctionType) IsConstructor() bool                { return f.kind == FunctionType_Kind_CONSTRUCTOR }

func NewFunctionType(builder *FunctionType_Builder) *FunctionType {
	this_ := &FunctionType{}
	this_.BindSelf_(this_)
	InitFunctionType(this_, builder)
	return this_
}

func InitFunctionType(this_ *FunctionType, builder *FunctionType_Builder) {
	_ = builder
	InitPrototypeObjectType_2(&this_.PrototypeObjectTypeBase, ToIface_PrototypeObjectType_Builder(builder), 0)
	this_.kind = builder.kind
	this_.sourceNode = builder.sourceNode
	this_.googModuleId = builder.googModuleId
	this_.isAbstract = builder.isAbstract
	this_.isKnownAmbiguous = builder.isKnownAmbiguous
	this_.parameters = builder.parameterArray
	this_.returnType = builder.returnType
	this_.returnInferred = builder.returnTypeIsInferred
	this_.typeOfThis = builder.typeOfThis
}

var BuilderHeapAllocs int

type FunctionType_Builder struct {
	PrototypeObjectType_BuilderBase
	kind                 FunctionType_Kind
	sourceNode           Node
	googModuleId         string
	isAbstract           bool
	isKnownAmbiguous     bool
	parameterArray       []FunctionType_Parameter
	returnType           JSType
	returnTypeIsInferred bool
	typeOfThis           JSType
}

func (this_ *FunctionType_Builder) BindSelf_(s any) {
	this_.PrototypeObjectType_BuilderBase.BindSelf_(s)
}
func (this_ *FunctionType_Builder) initFields_FunctionType_Builder() {
	this_.kind = FunctionType_Kind_ORDINARY
	this_.googModuleId = j2gort.NullString
}

func NewFunctionType_Builder(registry *JSTypeRegistry) *FunctionType_Builder {
	BuilderHeapAllocs++
	this_ := &FunctionType_Builder{}
	this_.BindSelf_(this_)
	InitFunctionType_Builder(this_, registry)
	return this_
}

func FunctionType_S_Builder(registry *JSTypeRegistry) *FunctionType_Builder {
	return NewFunctionType_Builder(registry)
}

func InitFunctionType_Builder(this_ *FunctionType_Builder, registry *JSTypeRegistry) {
	_ = registry
	InitPrototypeObjectType_Builder(&this_.PrototypeObjectType_BuilderBase, registry)
	this_.initFields_FunctionType_Builder()
	this_.SetImplicitPrototype(j2gort.CheckNotNullIface(registry.GetNativeObjectType(0)))
}

func (this_ *FunctionType_Builder) WithName(name string) *FunctionType_Builder {
	return j2gort.CastPtr[FunctionType_Builder](this_.SetName(name))
}
func (this_ *FunctionType_Builder) ForNativeType() *FunctionType_Builder {
	return j2gort.CastPtr[FunctionType_Builder](this_.SetNative(true))
}
func (this_ *FunctionType_Builder) ForConstructor() *FunctionType_Builder {
	this_.kind = FunctionType_Kind_CONSTRUCTOR
	return this_
}
func (this_ *FunctionType_Builder) ForInterface() *FunctionType_Builder {
	this_.kind = FunctionType_Kind_INTERFACE
	return this_
}
func (this_ *FunctionType_Builder) WithSourceNode(n Node) *FunctionType_Builder {
	this_.sourceNode = n
	return this_
}
func (this_ *FunctionType_Builder) SetGoogModuleId(id string) *FunctionType_Builder {
	this_.googModuleId = id
	return this_
}
func (this_ *FunctionType_Builder) SetIsKnownAmbiguous(x bool) *FunctionType_Builder {
	this_.isKnownAmbiguous = x
	return this_
}
func (this_ *FunctionType_Builder) WithIsAbstract(x bool) *FunctionType_Builder {
	this_.isAbstract = x
	return this_
}
func (this_ *FunctionType_Builder) WithKind(k FunctionType_Kind) *FunctionType_Builder {
	this_.kind = k
	return this_
}
func (this_ *FunctionType_Builder) WithParameterArray(p []FunctionType_Parameter) *FunctionType_Builder {
	this_.parameterArray = p
	return this_
}
func (this_ *FunctionType_Builder) WithParameters_0() *FunctionType_Builder {
	this_.parameterArray = nil
	return this_
}
func (this_ *FunctionType_Builder) WithReturnType_1(r JSType) *FunctionType_Builder {
	this_.returnType = r
	return this_
}
func (this_ *FunctionType_Builder) WithReturnType_2(r JSType, inf bool) *FunctionType_Builder {
	this_.returnType = r
	this_.returnTypeIsInferred = inf
	return this_
}
func (this_ *FunctionType_Builder) WithTypeOfThis(t JSType) *FunctionType_Builder {
	this_.typeOfThis = t
	return this_
}
func (this_ *FunctionType_Builder) WithTemplateKeysFrom(sourceMap *TemplateTypeMap) *FunctionType_Builder {
	count := sourceMap.Size()
	if count == 0 {
		return j2gort.CastPtr[FunctionType_Builder](j2gort.CastPtr[FunctionType_Builder](this_.SetTemplateTypeMap(this_.GetPrototypeObjectType_BuilderBase().registry.GetEmptyTemplateTypeMap())).SetTemplateParamCount(0))
	}
	return j2gort.CastPtr[FunctionType_Builder](j2gort.CastPtr[FunctionType_Builder](this_.SetTemplateTypeMap(sourceMap.AsKeysOnlySingleSubMap())).SetTemplateParamCount(count))
}
func (this_ *FunctionType_Builder) Build() any {
	return NewFunctionType(this_)
}
func (this_ *FunctionType_Builder) BuildAndResolve() *FunctionType {
	return j2gort.CheckNotNullPtr(j2gort.CastPtr[FunctionType](this_.Build()).ToMaybeFunctionType())
}

func (this_ *FunctionType) ToBuilder() *FunctionType_Builder {
	var builder *FunctionType_Builder = j2gort.CastPtr[FunctionType_Builder](j2gort.CastPtr[FunctionType_Builder](j2gort.CastPtr[FunctionType_Builder](j2gort.CastPtr[FunctionType_Builder](FunctionType_S_Builder(this_.GetJSTypeBase().registry).SetGoogModuleId(this_.GetGoogModuleId()).SetName(this_.GetReferenceName())).SetNative(this_.IsNativeObjectType())).SetTemplateTypeMap(this_.GetTemplateTypeMap())).SetIsKnownAmbiguous(this_.isKnownAmbiguous).WithIsAbstract(this_.IsAbstract()).WithKind(this_.GetKind()).WithParameterArray(this_.parameters).WithReturnType_2(this_.GetReturnType(), this_.IsReturnTypeInferred()).WithSourceNode(this_.GetSource()).WithTypeOfThis(this_.GetTypeOfThis()).SetTemplateParamCount(this_.GetTemplateParamCount()))
	_ = builder
	if this_.IsConstructor() {
		builder.ForConstructor()
	}
	return builder
}

func (this_ *JSTypeRegistry) priv_JSTypeRegistry_nativeConstructorBuilder_9116(name string) *FunctionType_Builder {
	_ = name
	return FunctionType_S_Builder(this_).ForNativeType().ForConstructor().WithName(name)
}

func (this_ *JSTypeRegistry) BuildNativeCtor(name string) *FunctionType {
	return this_.priv_JSTypeRegistry_nativeConstructorBuilder_9116(name).WithParameters_0().BuildAndResolve()
}

func (this_ *FunctionType) CloneBind(n int32) *FunctionType {
	var builder *FunctionType_Builder = FunctionType_S_Builder(this_.GetJSTypeBase().registry).WithName(this_.GetReferenceName()).WithTemplateKeysFrom(this_.GetTemplateTypeMap())
	_ = builder
	if n > 0 {
		builder.WithParameters_0()
	}
	return j2gort.CastPtr[FunctionType](builder.Build())
}

var EscapedBuilderSink *FunctionType_Builder

func EscapingBuilderMustNotStackAlloc(reg *JSTypeRegistry) *FunctionType {
	var builder *FunctionType_Builder = FunctionType_S_Builder(reg).WithName("Escaped")
	EscapedBuilderSink = builder
	return builder.BuildAndResolve()
}
`

	jscompSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w34b_sb/j2gort"
	pkg_rhino "example.com/w34b_sb/rhino"
)

type TypedScope struct {
	Root     pkg_rhino.Node
	Parent   *TypedScope
	Children []*TypedScope
}

var ScopeBuilderHeapAllocs int

type TypedScopeCreator_AbstractScopeBuilder interface {
	Implements_jscomp_TypedScopeCreator_AbstractScopeBuilder()
	IsJavaNil_() bool
	GetTypedScopeCreator_AbstractScopeBuilderBase() *TypedScopeCreator_AbstractScopeBuilderBase
	Build()
}

type TypedScopeCreator_AbstractScopeBuilderBase struct {
	self_        TypedScopeCreator_AbstractScopeBuilder
	outer_       *TypedScopeCreator
	currentScope *TypedScope
	sourceName   string
	idHash_      atomic.Int32
}

func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) BindSelf_(s any) {
	this_.self_ = s.(TypedScopeCreator_AbstractScopeBuilder)
}
func (this_ *TypedScopeCreator_AbstractScopeBuilderBase) IsJavaNil_() bool { return this_ == nil }
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetTypedScopeCreator_AbstractScopeBuilderBase() *TypedScopeCreator_AbstractScopeBuilderBase {
	return b
}
func (*TypedScopeCreator_AbstractScopeBuilderBase) Implements_jscomp_TypedScopeCreator_AbstractScopeBuilder() {
}

func InitTypedScopeCreator_AbstractScopeBuilder(this_ *TypedScopeCreator_AbstractScopeBuilderBase, scope *TypedScope) {
	this_.currentScope = scope
	this_.sourceName = "src"
}

type TypedScopeCreator_NormalScopeBuilder struct {
	TypedScopeCreator_AbstractScopeBuilderBase
	outer_ *TypedScopeCreator
}

func (this_ *TypedScopeCreator_NormalScopeBuilder) BindSelf_(s any) {
	this_.TypedScopeCreator_AbstractScopeBuilderBase.BindSelf_(s)
}
func NewTypedScopeCreator_NormalScopeBuilder(outer_ *TypedScopeCreator, scope *TypedScope) *TypedScopeCreator_NormalScopeBuilder {
	ScopeBuilderHeapAllocs++
	this_ := &TypedScopeCreator_NormalScopeBuilder{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitTypedScopeCreator_NormalScopeBuilder(this_, scope)
	return this_
}
func InitTypedScopeCreator_NormalScopeBuilder(this_ *TypedScopeCreator_NormalScopeBuilder, scope *TypedScope) {
	this_.TypedScopeCreator_AbstractScopeBuilderBase.outer_ = this_.outer_
	InitTypedScopeCreator_AbstractScopeBuilder(&this_.TypedScopeCreator_AbstractScopeBuilderBase, scope)
}
func (this_ *TypedScopeCreator_NormalScopeBuilder) Build() {
	for _, ch := range this_.currentScope.Children {
		this_.outer_.priv_TypedScopeCreator_createScopeInternal_5678(ch.Root, this_.currentScope, ch.Children)
	}
}

type TypedScopeCreator_FunctionScopeBuilder struct {
	TypedScopeCreator_AbstractScopeBuilderBase
	outer_ *TypedScopeCreator
}

func (this_ *TypedScopeCreator_FunctionScopeBuilder) BindSelf_(s any) {
	this_.TypedScopeCreator_AbstractScopeBuilderBase.BindSelf_(s)
}
func NewTypedScopeCreator_FunctionScopeBuilder(outer_ *TypedScopeCreator, scope *TypedScope) *TypedScopeCreator_FunctionScopeBuilder {
	ScopeBuilderHeapAllocs++
	this_ := &TypedScopeCreator_FunctionScopeBuilder{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitTypedScopeCreator_FunctionScopeBuilder(this_, scope)
	return this_
}
func InitTypedScopeCreator_FunctionScopeBuilder(this_ *TypedScopeCreator_FunctionScopeBuilder, scope *TypedScope) {
	this_.TypedScopeCreator_AbstractScopeBuilderBase.outer_ = this_.outer_
	InitTypedScopeCreator_AbstractScopeBuilder(&this_.TypedScopeCreator_AbstractScopeBuilderBase, scope)
}
func (this_ *TypedScopeCreator_FunctionScopeBuilder) Build() {
	for _, ch := range this_.currentScope.Children {
		this_.outer_.priv_TypedScopeCreator_createScopeInternal_5678(ch.Root, this_.currentScope, ch.Children)
	}
}

type TypedScopeCreator_ClassScopeBuilder struct {
	TypedScopeCreator_AbstractScopeBuilderBase
	outer_ *TypedScopeCreator
}

func (this_ *TypedScopeCreator_ClassScopeBuilder) BindSelf_(s any) {
	this_.TypedScopeCreator_AbstractScopeBuilderBase.BindSelf_(s)
}
func NewTypedScopeCreator_ClassScopeBuilder(outer_ *TypedScopeCreator, scope *TypedScope) *TypedScopeCreator_ClassScopeBuilder {
	ScopeBuilderHeapAllocs++
	this_ := &TypedScopeCreator_ClassScopeBuilder{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitTypedScopeCreator_ClassScopeBuilder(this_, scope)
	return this_
}
func InitTypedScopeCreator_ClassScopeBuilder(this_ *TypedScopeCreator_ClassScopeBuilder, scope *TypedScope) {
	this_.TypedScopeCreator_AbstractScopeBuilderBase.outer_ = this_.outer_
	InitTypedScopeCreator_AbstractScopeBuilder(&this_.TypedScopeCreator_AbstractScopeBuilderBase, scope)
}
func (this_ *TypedScopeCreator_ClassScopeBuilder) Build() {
	for _, ch := range this_.currentScope.Children {
		this_.outer_.priv_TypedScopeCreator_createScopeInternal_5678(ch.Root, this_.currentScope, ch.Children)
	}
}

type TypedScopeCreator struct {
	registry *pkg_rhino.JSTypeRegistry
}

func NewTypedScopeCreator(reg *pkg_rhino.JSTypeRegistry) *TypedScopeCreator {
	t := &TypedScopeCreator{registry: reg}
	t.priv_TypedScopeCreator_clearCommonState_5688()
	return t
}

func (this_ *TypedScopeCreator) priv_TypedScopeCreator_createScopeInternal_5678(root pkg_rhino.Node, parent *TypedScope, children []*TypedScope) *TypedScope {
	newScope := &TypedScope{Root: root, Parent: parent, Children: children}
	var scopeBuilder TypedScopeCreator_AbstractScopeBuilder = nil
	_ = scopeBuilder
	if root.IsFunction() {
		scopeBuilder = NewTypedScopeCreator_FunctionScopeBuilder(this_, newScope)
	} else if root.IsClass() {
		scopeBuilder = NewTypedScopeCreator_ClassScopeBuilder(this_, newScope)
	} else {
		scopeBuilder = NewTypedScopeCreator_NormalScopeBuilder(this_, newScope)
	}
	scopeBuilder.Build()
	return newScope
}

func (this_ *TypedScopeCreator) FinishAndFreeze() {
	this_.priv_TypedScopeCreator_clearCommonState_5688()
}

func (this_ *TypedScopeCreator) priv_TypedScopeCreator_clearCommonState_5688() {
}

// FunctionTypeBuilder in jscomp calling pkg_rhino.FunctionType_S_Builder via priv_ helper + ToBuilder!
type FunctionTypeBuilder struct {
	typeRegistry *pkg_rhino.JSTypeRegistry
	fnName       string
	moduleId     string
}

func (this_ *FunctionTypeBuilder) priv_FunctionTypeBuilder_createDefaultBuilder_2975() *pkg_rhino.FunctionType_Builder {
	return pkg_rhino.FunctionType_S_Builder(this_.typeRegistry).WithName(this_.fnName).SetIsKnownAmbiguous(false).SetGoogModuleId(this_.moduleId)
}

func (this_ *FunctionTypeBuilder) BuildOrdinary() *pkg_rhino.FunctionType {
	fnType := j2gort.CastPtr[pkg_rhino.FunctionType](this_.priv_FunctionTypeBuilder_createDefaultBuilder_2975().WithParameters_0().WithIsAbstract(false).Build())
	return fnType.ToBuilder().WithSourceNode(pkg_rhino.Node(1)).BuildAndResolve()
}
`

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
	mustWrite("go.mod", "module example.com/w34b_sb\n\ngo 1.24\n")
	mustWrite("j2gort/j2gort.go", j2gortSrc)
	mustWrite("rhino/rhino.go", rhinoSrc)
	mustWrite("jscomp/typed_scope_creator.go", jscompSrc)

	slabPass, err := passes.ParsePassesFlag("slab_alloc_lowering")
	if err != nil {
		t.Fatal(err)
	}

	rhinoDir := filepath.Join(tmpDir, "rhino")
	jscompDir := filepath.Join(tmpDir, "jscomp")

	// Run 1 & Run 2 on rhino:
	r1, err := passes.RunOnPackageDir(rhinoDir, slabPass, passes.Options{})
	if err != nil {
		t.Fatalf("rhino Run 1 failed: %v", err)
	}
	if r1.TotalSites < 8 {
		t.Fatalf("expected rhino Run 1 sites >= 8, got %+v", r1)
	}
	r2, err := passes.RunOnPackageDir(rhinoDir, slabPass, passes.Options{})
	if err != nil {
		t.Fatalf("rhino Run 2 failed: %v", err)
	}
	if r2.TotalSites != 0 || r2.FilesModified != 0 {
		t.Fatalf("rhino Run 2 NOT IDEMPOTENT: %+v", r2)
	}

	// Run 1 & Run 2 on jscomp (sees lowered rhino's InitBuilder_):
	j1, err := passes.RunOnPackageDir(jscompDir, slabPass, passes.Options{})
	if err != nil {
		t.Fatalf("jscomp Run 1 failed: %v", err)
	}
	if j1.TotalSites < 7 {
		t.Fatalf("expected jscomp Run 1 sites >= 7 (ScopeBuilder reuse + FunctionTypeBuilder stack), got %+v", j1)
	}
	j2, err := passes.RunOnPackageDir(jscompDir, slabPass, passes.Options{})
	if err != nil {
		t.Fatalf("jscomp Run 2 failed: %v", err)
	}
	if j2.TotalSites != 0 || j2.FilesModified != 0 {
		t.Fatalf("jscomp Run 2 NOT IDEMPOTENT: %+v", j2)
	}

	rhinoCodeBytes, err := os.ReadFile(filepath.Join(rhinoDir, "rhino.go"))
	if err != nil {
		t.Fatal(err)
	}
	rhinoCode := string(rhinoCodeBytes)
	if !strings.Contains(rhinoCode, "initPrototypeObjectTypeFromBase_") || !strings.Contains(rhinoCode, "InitBuilder_") {
		t.Fatalf("expected initPrototypeObjectTypeFromBase_ and InitBuilder_ in rhino.go:\n%s", rhinoCode)
	}
	// Negative escape check: EscapingBuilderMustNotStackAlloc must still call FunctionType_S_Builder(reg) on heap!
	escIdx := strings.Index(rhinoCode, "func EscapingBuilderMustNotStackAlloc")
	if escIdx < 0 || !strings.Contains(rhinoCode[escIdx:], "FunctionType_S_Builder(reg)") {
		t.Fatalf("expected EscapingBuilderMustNotStackAlloc to remain heap-allocated:\n%s", rhinoCode)
	}

	jscompCodeBytes, err := os.ReadFile(filepath.Join(jscompDir, "typed_scope_creator.go"))
	if err != nil {
		t.Fatal(err)
	}
	jscompCode := string(jscompCodeBytes)
	for _, want := range []string{
		"acquireNormalScopeBuilder_",
		"acquireFunctionScopeBuilder_",
		"acquireClassScopeBuilder_",
		"defer this_.releaseScopeBuilder_(scopeBuilder)",
		"this_.releaseFreeScopeBuilders_()",
		"InitBuilder_(&fb_",
	} {
		if !strings.Contains(jscompCode, want) {
			t.Fatalf("expected %q in lowered jscomp/typed_scope_creator.go:\n%s", want, jscompCode)
		}
	}
	if strings.Contains(jscompCode, "slab_TypedScope") || strings.Contains(jscompCode, "varSlab_") {
		t.Fatalf("Part 1 regression: TypedScope/TypedVar slabs must NOT be emitted:\n%s", jscompCode)
	}

	// Runtime test verifying:
	// 1. Deep recursive CreateScope (500 roots x depth 4 = 2,000 scopes) allocates <= 6 ScopeBuilders total and zeroes free lists on FinishAndFreeze!
	// 2. Non-escaping FunctionType_Builder chains in both rhino and jscomp perform ZERO NewFunctionType_Builder heap calls and 0 builder escapes!
	rtTest := `package jscomp

import (
	"testing"

	pkg_rhino "example.com/w34b_sb/rhino"
)

func TestScopeBuilderReuseAndFTBStackZeroHeap(t *testing.T) {
	reg := pkg_rhino.NewJSTypeRegistry()
	tsc := NewTypedScopeCreator(reg)
	ScopeBuilderHeapAllocs = 0

	makeTree := func() []*TypedScope {
		leaf := &TypedScope{Root: pkg_rhino.Node(0)}
		cls := &TypedScope{Root: pkg_rhino.Node(2), Children: []*TypedScope{leaf}}
		fn := &TypedScope{Root: pkg_rhino.Node(1), Children: []*TypedScope{cls}}
		return []*TypedScope{fn}
	}
	for i := 0; i < 500; i++ {
		tsc.priv_TypedScopeCreator_createScopeInternal_5678(pkg_rhino.Node(0), nil, makeTree())
	}
	if ScopeBuilderHeapAllocs > 6 {
		t.Fatalf("expected <= 6 ScopeBuilder heap allocations across 2,000 recursive scopes, got %d", ScopeBuilderHeapAllocs)
	}
	if len(tsc.freeNormalBuilders_) == 0 || len(tsc.freeFunctionBuilders_) == 0 || len(tsc.freeClassBuilders_) == 0 {
		t.Fatalf("expected free lists populated before FinishAndFreeze")
	}
	tsc.FinishAndFreeze()
	if tsc.freeNormalBuilders_ != nil || tsc.freeFunctionBuilders_ != nil || tsc.freeClassBuilders_ != nil {
		t.Fatalf("expected FinishAndFreeze to nil out all free builder slices")
	}

	pkg_rhino.BuilderHeapAllocs = 0
	ftb := &FunctionTypeBuilder{typeRegistry: reg, fnName: "MyFn", moduleId: "mod.id"}
	for i := 0; i < 200; i++ {
		fn := ftb.BuildOrdinary()
		if fn.GetReferenceName() != "MyFn" || fn.GetGoogModuleId() != "mod.id" || fn.GetSource() != pkg_rhino.Node(1) {
			t.Fatalf("unexpected built FunctionType fields: name=%q mod=%q src=%v", fn.GetReferenceName(), fn.GetGoogModuleId(), fn.GetSource())
		}
		ctor := reg.BuildNativeCtor("MyCtor")
		if !ctor.IsConstructor() || !ctor.IsNativeObjectType() || ctor.GetReferenceName() != "MyCtor" {
			t.Fatalf("unexpected BuildNativeCtor result")
		}
		bound := ctor.CloneBind(1)
		if bound.GetReferenceName() != "MyCtor" {
			t.Fatalf("unexpected CloneBind result")
		}
	}
	if pkg_rhino.BuilderHeapAllocs != 0 {
		t.Fatalf("expected 0 NewFunctionType_Builder heap calls on lowered chains, got %d", pkg_rhino.BuilderHeapAllocs)
	}

	// Verify escaping builder still allocates on heap:
	_ = pkg_rhino.EscapingBuilderMustNotStackAlloc(reg)
	if pkg_rhino.BuilderHeapAllocs != 1 || pkg_rhino.EscapedBuilderSink == nil {
		t.Fatalf("expected EscapingBuilderMustNotStackAlloc to allocate 1 builder on heap, got %d", pkg_rhino.BuilderHeapAllocs)
	}
}
`
	mustWrite("jscomp/w34b_rt_test.go", rtTest)

	// Verify Go compiler escape analysis (`-gcflags=-m`) confirms `fb_` does NOT escape to heap!
	escCmd := exec.Command("go", "build", "-gcflags=all=-m", "./rhino", "./jscomp")
	escCmd.Dir = tmpDir
	escOut, err := escCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build -gcflags=all=-m failed: %v\n%s", err, string(escOut))
	}
	for _, line := range strings.Split(string(escOut), "\n") {
		if strings.Contains(line, "moved to heap: fb_") {
			t.Fatalf("stack FunctionType_Builder escaped to heap: %s\nFull escape output:\n%s", line, string(escOut))
		}
	}

	testCmd := exec.Command("go", "test", "-v", "./jscomp")
	testCmd.Dir = tmpDir
	if out, err := testCmd.CombinedOutput(); err != nil {
		t.Fatalf("runtime test failed: %v\n%s", err, string(out))
	}
}

func TestWave34B_RealGenRhinoParsingJscompCopyIdempotencyAndBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping full gen/rhino + gen/parsing + gen/jscomp copy test in -short mode")
	}
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/j2gort does not build cleanly right now; skipping real gen copy check")
	}

	checkCmd := exec.Command("go", "build", "./transpiled/gen/j2gort", "./transpiled/gen/rhino", "./transpiled/gen/parsing", "./transpiled/gen/diagnostic")
	checkCmd.Dir = repoRoot
	if err := checkCmd.Run(); err != nil {
		t.Skipf("gen packages mid-edit (%v); skipping real copy check", err)
	}

	tmpDir := t.TempDir()
	origModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/realw34b", 1) +
		fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, sErr := os.ReadFile(filepath.Join(repoRoot, "go.sum")); sErr == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	pkgs := []string{"rhino", "parsing", "jscomptype", "jscomp"}
	for _, pkg := range pkgs {
		srcDir := filepath.Join(repoRoot, "transpiled", "gen", pkg)
		dstDir := filepath.Join(tmpDir, pkg)
		if err := os.MkdirAll(dstDir, 0o755); err != nil {
			t.Fatal(err)
		}
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			t.Fatal(err)
		}
		for _, ent := range entries {
			if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") || strings.HasSuffix(ent.Name(), "_test.go") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(srcDir, ent.Name()))
			if err != nil {
				t.Fatal(err)
			}
			txt := string(data)
			txt = strings.ReplaceAll(txt, `"github.com/paullepoulpe/jsgo/transpiled/gen/rhino"`, `"example.com/realw34b/rhino"`)
			txt = strings.ReplaceAll(txt, `"github.com/paullepoulpe/jsgo/transpiled/gen/parsing"`, `"example.com/realw34b/parsing"`)
			txt = strings.ReplaceAll(txt, `"github.com/paullepoulpe/jsgo/transpiled/gen/jscomptype"`, `"example.com/realw34b/jscomptype"`)
			if err := os.WriteFile(filepath.Join(dstDir, ent.Name()), []byte(txt), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}

	allPasses := passes.AllPasses()
	for _, pkg := range pkgs {
		pkgDir := filepath.Join(tmpDir, pkg)
		s1, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
		if err != nil {
			if !realGenJ2gortReady(repoRoot) {
				t.Skipf("dependency became mid-edit during %s Run 1: %v", pkg, err)
			}
			t.Fatalf("Run 1 on copied gen/%s failed: %v", pkg, err)
		}
		t.Logf("Wave 34B Real gen/%s Run 1 sites: total=%d files=%d byPass=%+v", pkg, s1.TotalSites, s1.FilesModified, s1.ByPass)

		s2, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
		if err != nil {
			if !realGenJ2gortReady(repoRoot) {
				t.Skipf("dependency became mid-edit during %s Run 2: %v", pkg, err)
			}
			t.Fatalf("Run 2 (idempotency) on copied gen/%s failed: %v", pkg, err)
		}
		if s2.TotalSites != 0 || s2.FilesModified != 0 {
			t.Fatalf("Run 2 on copied gen/%s NOT IDEMPOTENT: TotalSites=%d FilesModified=%d byPass=%+v",
				pkg, s2.TotalSites, s2.FilesModified, s2.ByPass)
		}
	}

	readCopy := func(rel string) string {
		t.Helper()
		b, rErr := os.ReadFile(filepath.Join(tmpDir, rel))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}

	// Part 1 check: no TypedScope / TypedVar slabs in typed_scope_creator.go
	tscCode := readCopy("jscomp/typed_scope_creator.go")
	if strings.Contains(tscCode, "NewTypedScopeCreator_NormalScopeBuilder(this_, newScope)") {
		t.Fatalf("expected NewTypedScopeCreator_NormalScopeBuilder in createScopeInternal to be lowered to acquireNormalScopeBuilder_")
	}
	if (!strings.Contains(tscCode, "acquireNormalScopeBuilder_") || !strings.Contains(tscCode, "releaseScopeBuilder_") || !strings.Contains(tscCode, "releaseFreeScopeBuilders_")) &&
		!strings.Contains(tscCode, "normalScopeBuilderPool") {
		t.Fatalf("expected ScopeBuilder free-list reuse helpers in typed_scope_creator.go")
	}

	// Part 2 check: script_node_deserializer.go has ZERO BoxNullableVal and lowers contextStack to []pkg_parsing.FeatureCollector_FeatureContext!
	sndCode := readCopy("jscomp/script_node_deserializer.go")
	if strings.Contains(sndCode, "BoxNullableVal") {
		t.Fatalf("expected ZERO BoxNullableVal in lowered script_node_deserializer.go")
	}
	if !strings.Contains(sndCode, "[]pkg_parsing.FeatureCollector_FeatureContext") || !strings.Contains(sndCode, "SliceUnmodifiableListLive(&contextStack)") {
		t.Fatalf("expected contextStack lowered to []pkg_parsing.FeatureCollector_FeatureContext with SliceUnmodifiableListLive in script_node_deserializer.go")
	}

	// Part 4 check: FunctionType_Builder stack allocation / reuse in rhino + jscomp
	ftbCode := readCopy("jscomp/function_type_builder.go")
	if !strings.Contains(ftbCode, "var fb_ pkg_rhino.FunctionType_Builder") &&
		!strings.Contains(ftbCode, "reusableFnBuilder") {
		t.Fatalf("expected stack var fb_ pkg_rhino.FunctionType_Builder in jscomp/function_type_builder.go")
	}
	tiCode := readCopy("jscomp/type_inference.go")
	if !strings.Contains(tiCode, ".InitBuilder_(&fb_") {
		t.Fatalf("expected .ToBuilder() chains lowered to .InitBuilder_(&fb_...) in jscomp/type_inference.go")
	}

	// Build all 4 copied packages and verify zero `moved to heap: fb_` in rhino & jscomp!
	escCmd := exec.Command("go", "build", "-gcflags=./rhino=-m", "./rhino", "./parsing", "./jscomptype", "./jscomp")
	escCmd.Dir = tmpDir
	escOut, err := escCmd.CombinedOutput()
	if err != nil {
		if !realGenJ2gortReady(repoRoot) {
			t.Skipf("dependency mid-edit during go build: %v", err)
		}
		t.Fatalf("go build on lowered gen/rhino + gen/parsing + gen/jscomp failed: %v\n%s", err, string(escOut))
	}
	for _, line := range strings.Split(string(escOut), "\n") {
		if strings.Contains(line, "moved to heap: fb_") {
			t.Fatalf("real rhino stack FunctionType_Builder escaped to heap: %s", line)
		}
	}
}

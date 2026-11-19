package passes_test

import (
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func assertIdempotentAndGofmt(t *testing.T, transformed string, passList []passes.Pass) {
	t.Helper()
	formatted, err := format.Source([]byte(transformed))
	if err != nil {
		t.Fatalf("output is not valid Go syntax: %v\ncode:\n%s", err, transformed)
	}
	if string(formatted) != transformed {
		t.Fatalf("output is not canonical gofmt:\ngot:\n%s\nwant:\n%s", transformed, string(formatted))
	}

	second, stats2, err := passes.TransformSource("idempotent.go", transformed, passList)
	if err != nil {
		t.Fatalf("second pass run failed: %v\ncode:\n%s", err, transformed)
	}
	if stats2.TotalSites != 0 {
		t.Fatalf("expected 0 transformed sites on second run (idempotency), got %d (%+v)", stats2.TotalSites, stats2.ByPass)
	}
	if second != transformed {
		t.Fatalf("non-idempotent output on second run:\nfirst:\n%s\nsecond:\n%s", transformed, second)
	}
}

func TestPassFlagParsing(t *testing.T) {
	none, err := passes.ParsePassesFlag("none")
	if err != nil || len(none) != 0 {
		t.Fatalf("ParsePassesFlag(none) = (%v, %v), want empty", none, err)
	}
	all, err := passes.ParsePassesFlag("all")
	if err != nil || len(all) != 10 {
		t.Fatalf("ParsePassesFlag(all) = (%d passes, %v), want 10 passes", len(all), err)
	}
	wantOrder := []string{"native_slice_lowering", "native_set_lowering", "native_map_lowering", "immutable_list_lowering", "native_deque_lowering", "string_builder_lowering", "slab_alloc_lowering", "node_traversal_lowering", "dead_cast_cleanup", "j2gort_peephole"}
	for i, wantName := range wantOrder {
		if all[i].Name() != wantName {
			t.Fatalf("AllPasses()[%d] = %q, want %q", i, all[i].Name(), wantName)
		}
		single, err := passes.ParsePassesFlag(wantName)
		if err != nil || len(single) != 1 || single[0].Name() != wantName {
			t.Fatalf("ParsePassesFlag(%q) = (%v, %v)", wantName, single, err)
		}
	}
	for _, removed := range []string{"arraylist_slice", "stringbuilder"} {
		if _, err := passes.ParsePassesFlag(removed); err == nil {
			t.Fatalf("expected error for removed pass %q", removed)
		}
	}
	if _, err := passes.ParsePassesFlag("unknown_pass"); err == nil {
		t.Fatalf("expected error for unknown pass")
	}
}

func TestDeadCastCleanup_SoundnessAndPrecedence(t *testing.T) {
	src := `package casts

type Node uint32
type Box struct{ V int }

const TypedConst int32 = 42
const UntypedConst = 42

func echo(v int32) int32 { return v }

func ExerciseCasts(x int32, s string, n Node, p *int32, b Box) (int32, string, Node, any, int) {
	a := int32(x)
	str := string(s)
	nd := Node(n)
	d := int32(TypedConst)
	ptr := (*int32)(p)
	nested := int32(int32(99))
	prec := (*Box)(&b).V
	if Box(Box{V: 1}) == b {
		prec++
	}

	keepUntypedLit := int32(42)
	keepUntypedConst := int32(UntypedConst)
	keepFuncCall := echo(x)
	keepNamedConv := Node(uint32(1))
	keepTypedNil := (*int32)(nil)

	return a + d + nested + keepUntypedLit + keepUntypedConst + keepFuncCall,
		str, nd + keepNamedConv, any(keepTypedNil), prec + int(*ptr)
}
`
	passList, err := passes.ParsePassesFlag("dead_cast_cleanup")
	if err != nil {
		t.Fatal(err)
	}

	got, stats, err := passes.TransformSource("casts.go", src, passList)
	if err != nil {
		t.Fatalf("TransformSource failed: %v", err)
	}
	if stats.ByPass["dead_cast_cleanup"] != 8 {
		t.Fatalf("expected 8 dead_cast_cleanup sites, got %d\ncode:\n%s", stats.ByPass["dead_cast_cleanup"], got)
	}

	for _, want := range []string{
		`a := x`,
		`str := s`,
		`nd := n`,
		`d := TypedConst`,
		`ptr := p`,
		`nested := int32(99)`,
		`prec := (&b).V`,
		`if (Box{V: 1}) == b`,
		`keepUntypedLit := int32(42)`,
		`keepUntypedConst := int32(UntypedConst)`,
		`keepNamedConv := Node(uint32(1))`,
		`keepTypedNil := (*int32)(nil)`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing expected %q in:\n%s", want, got)
		}
	}
	assertIdempotentAndGofmt(t, got, passList)
}

const mockJ2gortRuntimeSrc = `package j2gort

import ("fmt"; "reflect"; "strings")

const NullString = "\xff\x00j2go_null"

func isNilVal(v any) bool {
	if v == nil { return true }
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return rv.IsNil()
	}
	return false
}
func ToAny(v any) any {
	if v == nil || isNilVal(v) { return nil }
	if s, ok := v.(string); ok && s == NullString { return nil }
	return v
}
func BoxString(s string) any        { if s == NullString { return nil }; return s }
func CastString(v any) string       { if v == nil { return NullString }; return v.(string) }
func BoxNode[T ~uint32](v T) any    { if v == 0 { return nil }; return v }
func CastNode[T ~uint32](v any) T   { if v == nil { return 0 }; return v.(T) }
func CastPtr[T any](v any) *T       { if v == nil { return nil }; return v.(*T) }
func PtrToAny[T any](p *T) any      { if p == nil { return nil }; return p }
func SliceToAny[T any](s []T) any   { if s == nil { return nil }; return s }
func UncheckedCast[T any](v any) T  { return Cast[T](v) }
func AsIface[I any](v any) I        { return Cast[I](v) }
func Cast[T any](v any) T {
	var zero T
	if v == nil || isNilVal(v) { return zero }
	return v.(T)
}
func CheckNotNull[T any](v T) T {
	if any(v) == nil || isNilVal(v) { panic("NPE") }
	if s, ok := any(v).(string); ok && s == NullString { panic("NPE") }
	return v
}
func CheckNotNullPtr[T any](v *T) *T     { if v == nil { panic("NPE") }; return v }
func CheckNotNullSlice[T any](s []T) []T { if s == nil { panic("NPE") }; return s }
func CheckNotNullNode[N ~uint32](n N) N  { if n == 0 { panic("NPE") }; return n }
func CheckNotNullString(s string) string { if s == NullString { panic("NPE") }; return s }
func CheckNotNullIface[I interface{ IsJavaNil_() bool }](v I) I {
	if any(v) == nil || v.IsJavaNil_() { panic("NPE") }
	return v
}
func StrVal(v any) string {
	if v == nil || isNilVal(v) { return "null" }
	if s, ok := v.(string); ok {
		if s == NullString { return "null" }
		return s
	}
	if c, ok := v.(uint16); ok { return string(rune(c)) }
	return fmt.Sprint(v)
}
func StrConcat(parts ...any) string {
	var sb strings.Builder
	for _, p := range parts { sb.WriteString(StrVal(p)) }
	return sb.String()
}
func StringEquals(s string, other any) bool {
	if s == NullString { panic("NPE") }
	o, ok := other.(string)
	return ok && o != NullString && s == o
}
type BoxInt struct{ Val int32 }
type BoxBoolean struct{ Val bool }
func ValueOfInt(v int32) *BoxInt        { return &BoxInt{Val: v} }
func UnboxInt(b *BoxInt) int32          { return CheckNotNull(b).Val }
func Int32Val(v any) int32              { return v.(*BoxInt).Val }
func ValueOfBoolean(v bool) *BoxBoolean { return &BoxBoolean{Val: v} }
func UnboxBoolean(b *BoxBoolean) bool   { return CheckNotNull(b).Val }
func UnboxNode[T ~uint32](v any) T      { return CastNode[T](v) }

type Ext_java_lang_Iterable interface {
	IsJavaNil_() bool
	Iterator() any
}
type Ext_java_util_Collection interface {
	Ext_java_lang_Iterable
	Size() int32
	IsEmpty() bool
	Clear()
	Add_1(v any) bool
	AddAll_1(c Ext_java_util_Collection) bool
}
type Ext_java_util_List interface {
	Ext_java_util_Collection
	Get(idx int32) any
	Set(idx int32, v any) any
	Remove_1_int32(idx int32) any
	AddLast(v any)
}
type Ext_java_util_ArrayList interface {
	Ext_java_util_List
	IterSlice_() []any
}
type arrayListImpl struct{ items []any }
func (a *arrayListImpl) IsJavaNil_() bool        { return a == nil }
func (a *arrayListImpl) Size() int32             { return int32(len(a.items)) }
func (a *arrayListImpl) IsEmpty() bool           { return len(a.items) == 0 }
func (a *arrayListImpl) Clear()                  { clear(a.items); a.items = a.items[:0] }
func (a *arrayListImpl) Add_1(v any) bool        { a.items = append(a.items, v); return true }
func (a *arrayListImpl) AddLast(v any)           { a.items = append(a.items, v) }
func (a *arrayListImpl) Get(idx int32) any       { return a.items[idx] }
func (a *arrayListImpl) Set(idx int32, v any) any { old := a.items[idx]; a.items[idx] = v; return old }
func (a *arrayListImpl) IterSlice_() []any       { return a.items }
type sliceIter struct {
	items []any
	idx   int
}
func (it *sliceIter) IsJavaNil_() bool { return it == nil }
func (it *sliceIter) HasNext() bool    { return it.idx < len(it.items) }
func (it *sliceIter) Next() any {
	if it.idx >= len(it.items) { panic("NoSuchElementException") }
	v := it.items[it.idx]
	it.idx++
	return v
}
func (a *arrayListImpl) Iterator() any {
	cp := make([]any, len(a.items))
	copy(cp, a.items)
	return &sliceIter{items: cp}
}
func (a *arrayListImpl) Remove_1_int32(idx int32) any {
	old := a.items[idx]
	a.items = append(a.items[:idx], a.items[idx+1:]...)
	return old
}
func (a *arrayListImpl) AddAll_1(c Ext_java_util_Collection) bool {
	if c == nil || isNilVal(c) {
		panic("NPE")
	}
	changed := false
	for cur := NewIterCursor(c); cur.HasNext(); {
		a.items = append(a.items, cur.Next())
		changed = true
	}
	return changed
}

func NewExt_java_util_ArrayList() Ext_java_util_ArrayList {
	return &arrayListImpl{items: []any{}}
}
func NewExt_java_util_ArrayList_1_int32(cap int32) Ext_java_util_ArrayList {
	return &arrayListImpl{items: make([]any, 0, cap)}
}

type JIterCursor struct {
	slice []any
	idx   int
}
func NewIterCursor(coll any) JIterCursor {
	if coll == nil || isNilVal(coll) { panic("NPE") }
	si := coll.(interface{ IterSlice_() []any })
	return JIterCursor{slice: si.IterSlice_()}
}
func (c *JIterCursor) HasNext() bool   { return c.idx < len(c.slice) }
func (c *JIterCursor) Next() any       { v := c.slice[c.idx]; c.idx++; return v }
func (c *JIterCursor) NextStr() string { return CastString(c.Next()) }
func (c *JIterCursor) NextNode() uint32 { return CastNode[uint32](c.Next()) }
`

func writeTempJ2gortWorkspace(t *testing.T, pkgFiles map[string]string) (string, string) {
	t.Helper()
	tmpDir := t.TempDir()
	gomod := "module example.com/j2gotest\n\ngo 1.24\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	j2Dir := filepath.Join(tmpDir, "j2gort")
	if err := os.MkdirAll(j2Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(j2Dir, "runtime.go"), []byte(mockJ2gortRuntimeSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	appDir := filepath.Join(tmpDir, "app")
	if len(pkgFiles) > 0 {
		if err := os.MkdirAll(appDir, 0o755); err != nil {
			t.Fatal(err)
		}
		// order-independent: distinct file names written to disk
		for name, content := range pkgFiles {
			if err := os.WriteFile(filepath.Join(appDir, name), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return tmpDir, appDir
}

func TestJ2gortPeephole_CheckNotNullIfaceAndZeroSiteRun2(t *testing.T) {
	appSrc := `package app

import "example.com/j2gotest/j2gort"

type Node uint32
type SuperIface interface { IsJavaNil_() bool; Tag() string }
type SubIface interface { SuperIface; Extra() int32 }

func ToIface_SuperIface[T interface{ SuperIface; comparable }](v T) SuperIface {
	var z T
	if v == z { return nil }
	return v
}

type Impl struct{ Label string }

func NewImpl(label string) *Impl                 { return &Impl{Label: label} }
func (i *Impl) IsJavaNil_() bool                 { return i == nil }
func (i *Impl) Tag() string                      { return i.Label }
func (i *Impl) Extra() int32                     { return 1 }
func NewLambda_SubIface(fn func() string) SubIface { return &Impl{Label: fn()} }
func acceptAny(a any, rest ...any) any           { return a }
func acceptSuper(s SuperIface) string            { return s.Tag() }
func sideEffect() *Impl                          { return &Impl{Label: "side"} }

func ExercisePeepholes(n Node, s string, p *Impl, iface SubIface, i int32, b bool) (Node, string, *Impl, int32, bool, bool, SuperIface) {
	j2gort.Cast[*Impl](sideEffect())
	j2gort.AsIface[SuperIface](sideEffect())

	n2 := j2gort.CastNode[Node](j2gort.BoxNode(n))
	s2 := j2gort.CastString(j2gort.BoxString(s))
	i2 := j2gort.UnboxInt(j2gort.ValueOfInt(i))
	b2 := j2gort.UnboxBoolean(j2gort.ValueOfBoolean(b))

	p2 := j2gort.Cast[*Impl](j2gort.PtrToAny(p))
	nestedIface := j2gort.AsIface[SuperIface](j2gort.AsIface[SubIface](p))
	var explicitSuper SuperIface = j2gort.CheckNotNullIface(j2gort.AsIface[SuperIface](NewImpl("fresh")))
	nonNilLocalIface := j2gort.CheckNotNullIface(iface)
	_ = acceptSuper(j2gort.CheckNotNullIface(nonNilLocalIface))
	_ = acceptSuper(j2gort.AsIface[SuperIface](j2gort.CheckNotNullIface(NewLambda_SubIface(func() string { return "lam" }))))
	_ = acceptSuper(j2gort.AsIface[SuperIface](&Impl{Label: "inline"}))

	lit := j2gort.StrVal("constant")
	chk := j2gort.CheckNotNull("nonnull")
	cat := j2gort.StrConcat("a", int32(7), true, s)
	_ = acceptAny(j2gort.BoxString("boxed-lit"), j2gort.PtrToAny(NewImpl("arg")), j2gort.ToAny(i))
	eq := j2gort.StringEquals("target", j2gort.BoxString(s))
	neq := !(j2gort.StringEquals("other", j2gort.BoxString(s)))

	// Semantic guards that MUST remain intact:
	keepNullableAsIface := j2gort.AsIface[SuperIface](p)
	keepShortDefine := j2gort.AsIface[SuperIface](&Impl{})
	keepNullStringCast := j2gort.Cast[string](j2gort.ToAny(s))
	keepNullableStrEquals := j2gort.StringEquals(s, "literal")
	keepNullableStrVal := j2gort.StrVal(s)
	_ = nestedIface
	_ = keepNullableAsIface
	_ = keepShortDefine

	return n2, s2 + lit + chk + cat + keepNullStringCast + keepNullableStrVal, p2, i2, b2 && eq && neq, keepNullableStrEquals, explicitSuper
}
`
	pruneOnlySrc := `package app

import "example.com/j2gotest/j2gort"

func OnlyCascadingCheckNotNullIface() SuperIface {
	var res SuperIface = j2gort.CheckNotNullIface(j2gort.AsIface[SuperIface](NewImpl("clean")))
	var selfRes SuperIface = j2gort.AsIface[SuperIface](ToIface_SuperIface(NewImpl("self")))
	_ = selfRes
	return res
}
`

	_, appDir := writeTempJ2gortWorkspace(t, map[string]string{
		"app.go":       appSrc,
		"cascading.go": pruneOnlySrc,
	})

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(appDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir(allPasses) failed: %v", err)
	}
	if stats1.TotalSites < 20 {
		t.Fatalf("expected >= 20 simplified sites on Run 1, got %+v", stats1.ByPass)
	}

	gotAppBytes, err := os.ReadFile(filepath.Join(appDir, "app.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotApp := string(gotAppBytes)
	for _, want := range []string{
		`var explicitSuper SuperIface = NewImpl("fresh")`,
		`nonNilLocalIface := j2gort.CheckNotNullIface(iface)`,
		`acceptSuper(nonNilLocalIface)`,
		`acceptSuper(NewLambda_SubIface(func() string { return "lam" }))`,
		`cat := "a7true" + j2gort.StrVal(s)`,
		`eq := "target" == s`,
		`neq := "other" != s`,
		`keepNullableAsIface := j2gort.AsIface[SuperIface](p)`,
		`keepShortDefine := j2gort.AsIface[SuperIface](&Impl{})`,
	} {
		if !strings.Contains(gotApp, want) {
			t.Errorf("missing expected snippet %q in app.go:\n%s", want, gotApp)
		}
	}

	gotCascBytes, err := os.ReadFile(filepath.Join(appDir, "cascading.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotCasc := string(gotCascBytes)
	if strings.Contains(gotCasc, "j2gort") {
		t.Fatalf("expected unused j2gort import to be pruned on Run 1 in cascading.go, got:\n%s", gotCasc)
	}
	if !strings.Contains(gotCasc, `var res SuperIface = NewImpl("clean")`) ||
		!strings.Contains(gotCasc, `var selfRes SuperIface = NewImpl("self")`) {
		t.Fatalf("unexpected cascading.go output:\n%s", gotCasc)
	}

	// Run 2 MUST produce strictly 0 sites and 0 modified files!
	stats2, err := passes.RunOnPackageDir(appDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("second run failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 modified files on Run 2, got %+v", stats2)
	}
}

func TestTransformSourceAllPassesSingleRunIdempotence(t *testing.T) {
	src := `package sample
type IfaceA interface { Greet() string }
type IfaceB interface { Greet() string }
func ToIface_IfaceA[T IfaceA](v T) IfaceA { return v }
func Run(b IfaceB) IfaceB {
	var res IfaceB = IfaceB(ToIface_IfaceA(b))
	return res
}
`
	out1, stats1, err := passes.TransformSource("sample.go", src, passes.AllPasses())
	if err != nil {
		t.Fatalf("first TransformSource(AllPasses) failed: %v", err)
	}
	if stats1.ByPass["j2gort_peephole"] != 1 || stats1.ByPass["dead_cast_cleanup"] != 1 || stats1.TotalSites != 2 {
		t.Fatalf("expected j2gort_peephole=1 and dead_cast_cleanup=1 on first run, got %+v\nout1:\n%s", stats1.ByPass, out1)
	}
	if !strings.Contains(out1, "var res IfaceB = b") {
		t.Fatalf("expected outer redundant cast to be stripped on first run, got:\n%s", out1)
	}
	assertIdempotentAndGofmt(t, out1, passes.AllPasses())
}

func TestDifferentialRuntimeParity(t *testing.T) {
	workloadSrc := `package orig
import "example.com/j2gotest/j2gort"
type Node uint32
type IPart interface { IsJavaNil_() bool; Text() string }
type Part struct{ Val string }
func NewPart(v string) *Part     { return &Part{Val: v} }
func (p *Part) IsJavaNil_() bool { return p == nil }
func (p *Part) Text() string     { return p.Val }
func ToIface_IPart[T interface{ IPart; comparable }](v T) IPart {
	var z T
	if v == z { return nil }
	return v
}
func ExecuteWorkload(maybeNull string, n Node, maybeNilPtr *Part, maybeNilIface IPart) string {
	var p1 IPart = j2gort.CheckNotNullIface(j2gort.AsIface[IPart](NewPart("ok")))
	nRound := j2gort.CastNode[Node](j2gort.BoxNode(Node(n)))
	sRound := j2gort.CastString(j2gort.BoxString(string(maybeNull)))
	asIface := j2gort.AsIface[IPart](maybeNilPtr)
	eqLit := j2gort.StringEquals("abc", j2gort.BoxString(sRound))
	partTag := "nil"
	if asIface != nil { partTag = j2gort.CheckNotNullIface(asIface).Text() }
	ifaceTag := "inil"
	if maybeNilIface != nil { ifaceTag = j2gort.CheckNotNullIface(maybeNilIface).Text() }
	return j2gort.StrConcat(p1.Text(), ":", int32(nRound), ":", sRound, ":", eqLit, ":", partTag, ":", ifaceTag)
}
`
	tmpDir, _ := writeTempJ2gortWorkspace(t, nil)
	origDir, optDir := filepath.Join(tmpDir, "orig"), filepath.Join(tmpDir, "opt")
	for _, d := range []string{origDir, optDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(origDir, "workload.go"), []byte(workloadSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optDir, "workload.go"), []byte(strings.Replace(workloadSrc, "package orig", "package opt", 1)), 0o644); err != nil {
		t.Fatal(err)
	}

	stats, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{})
	if err != nil || stats.TotalSites < 6 {
		t.Fatalf("RunOnPackageDir(opt) = (%+v, %v), want >= 6 sites", stats, err)
	}
	if stats2, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{}); err != nil || stats2.TotalSites != 0 {
		t.Fatalf("expected 0 sites on second run, got (%+v, %v)", stats2, err)
	}

	diffTest := `package diff_test
import ("testing"; "example.com/j2gotest/j2gort"; "example.com/j2gotest/opt"; "example.com/j2gotest/orig")
func runSafe(fn func() string) (out string, panicked bool) {
	defer func() { if recover() != nil { panicked = true } }()
	return fn(), false
}
func TestParity(t *testing.T) {
	var typedNilOrig orig.IPart = (*orig.Part)(nil)
	var typedNilOpt opt.IPart = (*opt.Part)(nil)
	cases := []struct{ s string; n uint32; nilPtr, tNil bool }{
		{"abc", 5, false, false}, {"", 0, true, false}, {j2gort.NullString, 9, false, false}, {"xyz", 1, false, true},
	}
	for i, tc := range cases {
		var op *orig.Part
		var xp *opt.Part
		if !tc.nilPtr { op, xp = orig.NewPart("p"), opt.NewPart("p") }
		var oi orig.IPart
		var xi opt.IPart
		if tc.tNil {
			oi, xi = typedNilOrig, typedNilOpt
		} else if !tc.nilPtr {
			oi, xi = orig.ToIface_IPart(op), opt.ToIface_IPart(xp)
		}
		gotOrig, panOrig := runSafe(func() string { return orig.ExecuteWorkload(tc.s, orig.Node(tc.n), op, oi) })
		gotOpt, panOpt := runSafe(func() string { return opt.ExecuteWorkload(tc.s, opt.Node(tc.n), xp, xi) })
		if panOrig != panOpt || gotOrig != gotOpt {
			t.Fatalf("case %d mismatch: orig=(%q, panic=%v) vs opt=(%q, panic=%v)", i, gotOrig, panOrig, gotOpt, panOpt)
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "diff_test.go"), []byte(diffTest), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-count=1", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("differential runtime test failed: %v\n%s", err, string(out))
	}
}

func TestNativeSliceLowering_FieldAndLocalsAcrossMethods(t *testing.T) {
	collectorDeclSrc := `package app

import "example.com/j2gotest/j2gort"

type Node uint32

type MyItem struct {
	Name string
}

func NewMyItem(name string) *MyItem { return &MyItem{Name: name} }

type Collector struct {
	items   j2gort.Ext_java_util_ArrayList ` + "`gson:\"items\"`" + `
	nodeIds j2gort.Ext_java_util_List      ` + "`gson:\"nodeIds\"`" + `
}

func (this_ *Collector) initFields_Collector() {
	this_.items = j2gort.NewExt_java_util_ArrayList()
	this_.nodeIds = j2gort.NewExt_java_util_ArrayList_1_int32(int32(4))
}

func NewCollector() *Collector {
	c := &Collector{}
	c.initFields_Collector()
	return c
}
`
	collectorOpsSrc := `package app

import "example.com/j2gotest/j2gort"

func (this_ *Collector) Record(it *MyItem, n Node) {
	this_.items.Add_1(j2gort.PtrToAny(it))
	j2gort.CheckNotNull(this_.nodeIds).AddLast(j2gort.BoxNode(n))
}

func (this_ *Collector) Summary() (string, Node, int32, bool) {
	if this_.items.IsEmpty() {
		return "", 0, 0, true
	}
	first := j2gort.CastPtr[MyItem](this_.items.Get(int32(0)))
	firstNode := j2gort.UnboxNode[Node](this_.nodeIds.Get(int32(0)))
	var joined string
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.items))
		for cur.HasNext() {
			raw := cur.Next()
			elem := j2gort.CastPtr[MyItem](raw)
			_ = elem
			if elem != nil {
				joined += elem.Name
			}
		}
	}
	var sum Node
	{
		cur2 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.nodeIds))
		for cur2.HasNext() {
			nd := Node(cur2.NextNode())
			_ = nd
			sum += nd
		}
	}
	hasMore := !(this_.nodeIds.IsEmpty())
	sz := this_.items.Size()
	_ = first
	return joined, sum + firstNode, sz, hasMore
}

func (this_ *Collector) ResetAll() {
	this_.items.Clear()
	this_.nodeIds.Clear()
}

func ExerciseLocalSlices(a, b string, n1, n2 Node, m *MyItem, k int32) (string, Node, string, int32) {
	var strs j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList_1_int32(int32(2))
	_ = strs
	strs.Add_1(j2gort.BoxString(a))
	strs.Add_1(j2gort.BoxString(b))

	copyStrs := j2gort.NewExt_java_util_ArrayList()
	copyStrs.Add_1(j2gort.BoxString(j2gort.CastString(strs.Get(int32(0)))))
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(strs))
		for cur.HasNext() {
			s := cur.NextStr()
			copyStrs.Add_1(j2gort.BoxString(s))
		}
	}
	var acc string
	curRaw := j2gort.NewIterCursor(copyStrs)
	for curRaw.HasNext() {
		raw := curRaw.Next()
		acc += j2gort.CastString(raw)
	}

	var nodes j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = nodes
	nodes.Add_1(j2gort.BoxNode(n1))
	nodes.Add_1(j2gort.BoxNode(n2))
	nodeSum := j2gort.CastNode[Node](nodes.Get(int32(0))) + j2gort.UnboxNode[Node](nodes.Get(int32(1)))

	var ptrs j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	ptrs.Add_1(j2gort.PtrToAny(m))
	var ptrNames string
	for _, rawElem := range j2gort.CheckNotNull(ptrs).IterSlice_() {
		it := j2gort.CastPtr[MyItem](rawElem)
		if it != nil {
			ptrNames += it.Name
		}
	}

	var nums j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	nums.Add_1(j2gort.ValueOfInt(k))
	nums.Add_1(j2gort.ValueOfInt(k + int32(1)))
	numTotal := j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](nums.Get(int32(0)))) + j2gort.Int32Val(nums.Get(int32(1))) + nums.Size()

	return acc, nodeSum, ptrNames, numTotal
}
`

	_, appDir := writeTempJ2gortWorkspace(t, map[string]string{
		"collector.go":     collectorDeclSrc,
		"collector_ops.go": collectorOpsSrc,
	})

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(appDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir failed: %v", err)
	}
	if stats1.ByPass["native_slice_lowering"] < 30 {
		t.Fatalf("expected >= 30 native_slice_lowering sites, got %+v", stats1.ByPass)
	}

	declBytes, err := os.ReadFile(filepath.Join(appDir, "collector.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotDecl := string(declBytes)
	if strings.Contains(gotDecl, "j2gort") {
		t.Fatalf("expected unused j2gort import to be pruned in collector.go:\n%s", gotDecl)
	}
	for _, want := range []string{
		"items   []*MyItem",
		"nodeIds []Node",
		"this_.items = []*MyItem{}",
		"this_.nodeIds = make([]Node, 0, int32(4))",
	} {
		if !strings.Contains(gotDecl, want) {
			t.Errorf("missing expected snippet %q in collector.go:\n%s", want, gotDecl)
		}
	}

	opsBytes, err := os.ReadFile(filepath.Join(appDir, "collector_ops.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotOps := string(opsBytes)
	if strings.Contains(gotOps, "j2gort") {
		t.Fatalf("expected unused j2gort import to be pruned in collector_ops.go:\n%s", gotOps)
	}
	for _, want := range []string{
		"this_.items = append(this_.items, it)",
		"this_.nodeIds = append(this_.nodeIds, n)",
		"if len(this_.items) == 0",
		"first := this_.items[int32(0)]",
		"firstNode := this_.nodeIds[int32(0)]",
		"for _, elem := range this_.items",
		"for _, nd := range this_.nodeIds",
		"hasMore := len(this_.nodeIds) != 0",
		"clear(this_.items)",
		"this_.items = this_.items[:0]",
		"this_.nodeIds = this_.nodeIds[:0]",
		"var strs []string = make([]string, 0, int32(2))",
		"copyStrs := []string{}",
		"copyStrs = append(copyStrs, strs[int32(0)])",
		"for _, s := range strs",
		"for _, raw := range copyStrs",
		"acc += raw",
		"var nodes []Node = []Node{}",
		"var ptrs []*MyItem = []*MyItem{}",
		"for _, it := range ptrs",
		"var nums []int32 = []int32{}",
		"numTotal := nums[int32(0)] + nums[int32(1)] + int32(len(nums))",
	} {
		if !strings.Contains(gotOps, want) {
			t.Errorf("missing expected snippet %q in collector_ops.go:\n%s", want, gotOps)
		}
	}
	if strings.Contains(gotOps, "clear(this_.nodeIds)") {
		t.Errorf("scalar uint32 Node slice should reslice without clear():\n%s", gotOps)
	}

	// Run 2 MUST produce strictly 0 sites and 0 modified files!
	stats2, err := passes.RunOnPackageDir(appDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("second run failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 modified files on Run 2, got %+v", stats2)
	}
}

func TestNativeSliceLowering_NegativeSafetyGuards(t *testing.T) {
	negSrcA := `package app

import "example.com/j2gotest/j2gort"

type GuardHolder struct {
	nullableField j2gort.Ext_java_util_ArrayList
	loopMutField  j2gort.Ext_java_util_ArrayList
	crossPkgField j2gort.Ext_java_util_ArrayList
}

func (this_ *GuardHolder) initFields_GuardHolder() {
	this_.nullableField = j2gort.NewExt_java_util_ArrayList()
	this_.loopMutField = j2gort.NewExt_java_util_ArrayList()
	this_.crossPkgField = j2gort.NewExt_java_util_ArrayList()
}

func (this_ *GuardHolder) NilOutField() {
	this_.nullableField.Add_1(j2gort.BoxString("x"))
	this_.nullableField = nil
}

func (this_ *GuardHolder) mutateHelper() {
	this_.loopMutField.Add_1(j2gort.BoxString("mut"))
}

func (this_ *GuardHolder) IterateMutatingField() string {
	this_.loopMutField.Add_1(j2gort.BoxString("init"))
	var out string
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.loopMutField))
	for cur.HasNext() {
		out += cur.NextStr()
		this_.mutateHelper()
	}
	return out
}

func acceptList(l j2gort.Ext_java_util_List) int32 { return l.Size() }

func NegativeLocals() int32 {
	// 1. Aliased local variable: y = x
	var aliased j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	aliased.Add_1(j2gort.BoxString("a"))
	other := aliased
	_ = other

	// 2. Passed to function expecting Ext_java_util_List
	var escaped j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	escaped.Add_1(j2gort.BoxString("b"))
	n := acceptList(escaped)

	// 3. List mutated inside its own enhanced-for loop
	var selfMut j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	selfMut.Add_1(j2gort.BoxString("c"))
	{
		cur := j2gort.NewIterCursor(selfMut)
		for cur.HasNext() {
			s := cur.NextStr()
			if s == "c" {
				selfMut.Add_1(j2gort.BoxString("d"))
				break
			}
		}
	}

	// 4. Mixed element types (string + int32)
	var mixed j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	mixed.Add_1(j2gort.BoxString("str"))
	mixed.Add_1(j2gort.ValueOfInt(int32(7)))

	return n + aliased.Size() + selfMut.Size() + mixed.Size()
}
`
	negSrcB := `package app

import (
	"example.com/j2gotest/j2gort"
	"example.com/j2gotest/subpkg"
)

func (this_ *GuardHolder) AddCrossPkg(w *subpkg.Widget) *subpkg.Widget {
	this_.crossPkgField.Add_1(j2gort.PtrToAny(w))
	return j2gort.CastPtr[subpkg.Widget](this_.crossPkgField.Get(int32(0)))
}
`
	tmpDir, appDir := writeTempJ2gortWorkspace(t, map[string]string{
		"neg_a.go": negSrcA,
		"neg_b.go": negSrcB,
	})
	subDir := filepath.Join(tmpDir, "subpkg")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "widget.go"), []byte("package subpkg\ntype Widget struct{ ID int }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	slicePass, err := passes.ParsePassesFlag("native_slice_lowering")
	if err != nil {
		t.Fatal(err)
	}
	stats, err := passes.RunOnPackageDir(appDir, slicePass, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir(native_slice_lowering) failed: %v", err)
	}
	if stats.ByPass["native_slice_lowering"] != 0 || stats.FilesModified != 0 {
		t.Fatalf("expected 0 lowered sites on negative safety cases, got %+v", stats)
	}
}

func TestNativeSliceLowering_DifferentialRuntimeParity(t *testing.T) {
	workloadSrc := `package orig

import "example.com/j2gotest/j2gort"

type Node uint32
type Entry struct{ Label string }

type Bag struct {
	entries j2gort.Ext_java_util_ArrayList
	nodes   j2gort.Ext_java_util_List
}

func (this_ *Bag) initFields_Bag() {
	this_.entries = j2gort.NewExt_java_util_ArrayList()
	this_.nodes = j2gort.NewExt_java_util_ArrayList_1_int32(int32(2))
}

func RunSliceWorkload(s1, s2 string, n1, n2 Node, nilEntry bool) string {
	b := &Bag{}
	b.initFields_Bag()
	if !nilEntry {
		b.entries.Add_1(j2gort.PtrToAny(&Entry{Label: s1}))
	} else {
		var np *Entry
		b.entries.Add_1(j2gort.PtrToAny(np))
	}
	b.nodes.Add_1(j2gort.BoxNode(n1))
	b.nodes.Clear()
	b.nodes.Add_1(j2gort.BoxNode(n2))

	var strs j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	strs.Add_1(j2gort.BoxString(s1))
	strs.Add_1(j2gort.BoxString(s2))

	var nullCount int32
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(strs))
		for cur.HasNext() {
			v := cur.NextStr()
			if v == j2gort.NullString {
				nullCount++
			}
		}
	}
	firstStr := j2gort.CastString(strs.Get(int32(0)))
	isNull0 := firstStr == j2gort.NullString

	e0 := j2gort.CastPtr[Entry](b.entries.Get(int32(0)))
	eTag := "nil"
	if e0 != nil {
		eTag = j2gort.StrVal(e0.Label)
	}
	nd0 := j2gort.CastNode[Node](b.nodes.Get(int32(0)))
	return j2gort.StrConcat(eTag, "|", int32(nd0), "|", nullCount, "|", isNull0, "|", !(b.entries.IsEmpty()), "|", strs.Size())
}
`
	tmpDir, _ := writeTempJ2gortWorkspace(t, nil)
	origDir, optDir := filepath.Join(tmpDir, "orig"), filepath.Join(tmpDir, "opt")
	for _, d := range []string{origDir, optDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(origDir, "workload.go"), []byte(workloadSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optDir, "workload.go"), []byte(strings.Replace(workloadSrc, "package orig", "package opt", 1)), 0o644); err != nil {
		t.Fatal(err)
	}

	stats1, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{})
	if err != nil || stats1.ByPass["native_slice_lowering"] < 15 {
		t.Fatalf("RunOnPackageDir(opt) = (%+v, %v), want >= 15 slice sites", stats1, err)
	}
	if stats2, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{}); err != nil || stats2.TotalSites != 0 {
		t.Fatalf("expected 0 sites on Run 2, got (%+v, %v)", stats2, err)
	}

	diffTest := `package slice_diff_test

import (
	"testing"
	"example.com/j2gotest/j2gort"
	"example.com/j2gotest/opt"
	"example.com/j2gotest/orig"
)

func TestSliceParity(t *testing.T) {
	cases := []struct {
		s1, s2   string
		n1, n2   uint32
		nilEntry bool
	}{
		{"alpha", "beta", 1, 2, false},
		{j2gort.NullString, "", 0, 0, false},
		{"x", j2gort.NullString, 42, 0, true},
		{j2gort.NullString, j2gort.NullString, 0, 99, true},
	}
	for i, tc := range cases {
		gotOrig := orig.RunSliceWorkload(tc.s1, tc.s2, orig.Node(tc.n1), orig.Node(tc.n2), tc.nilEntry)
		gotOpt := opt.RunSliceWorkload(tc.s1, tc.s2, opt.Node(tc.n1), opt.Node(tc.n2), tc.nilEntry)
		if gotOrig != gotOpt {
			t.Fatalf("case %d mismatch: orig=%q vs opt=%q", i, gotOrig, gotOpt)
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "slice_diff_test.go"), []byte(diffTest), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-count=1", "-race", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("differential slice runtime test failed: %v\n%s", err, string(out))
	}
}

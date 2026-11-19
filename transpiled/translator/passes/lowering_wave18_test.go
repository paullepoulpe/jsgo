package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

const wave18MockJ2gortSrc = `package j2gort

const NullString = "\xff\x00j2go_null"

func CheckNotNullString(s string) string {
	if s == NullString {
		panic("NPE")
	}
	return s
}

func StrValString(s string) string {
	if s == NullString {
		return "null"
	}
	return s
}

func StrValInt32(v int32) string { return "0" }

func CheckNotNullIface[I interface{ IsJavaNil_() bool }](v I) I {
	if any(v) == nil || v.IsJavaNil_() {
		panic("NPE")
	}
	return v
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
	var zero T
	if v == nil {
		return zero, false
	}
	r, ok := v.(T)
	return r, ok
}

func PtrToAny[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}

func CheckStateFmt(cond bool, format string, args ...any) {
	if !cond {
		panic(format)
	}
}

type Ext_java_util_BitSet struct {
	buf   [4]uint64
	extra []uint64
	n     int32
}

func (b *Ext_java_util_BitSet) IsJavaNil_() bool { return b == nil }

func InitExt_java_util_BitSet_0(b *Ext_java_util_BitSet) {
	b.n = 0
	b.extra = nil
}

func InitExt_java_util_BitSet_1(b *Ext_java_util_BitSet, nbits int32) {
	if nbits <= 256 {
		if nbits < 0 {
			bitSetFailNeg()
		}
		b.n = 0
		b.extra = nil
		return
	}
	bitSetGrowInit(b, nbits)
}

func bitSetFailNeg() { panic("NegativeArraySizeException") }

func bitSetGrowInit(b *Ext_java_util_BitSet, nbits int32) {
	need := int((int64(nbits) + 63) >> 6)
	b.extra = make([]uint64, 0, need)
	b.n = 0
}

func NewExt_java_util_BitSet() *Ext_java_util_BitSet {
	obj := &Ext_java_util_BitSet{}
	InitExt_java_util_BitSet_0(obj)
	return obj
}

func NewExt_java_util_BitSet_0() *Ext_java_util_BitSet {
	return NewExt_java_util_BitSet()
}

func NewExt_java_util_BitSet_1(nbits int32) *Ext_java_util_BitSet {
	obj := &Ext_java_util_BitSet{}
	InitExt_java_util_BitSet_1(obj, nbits)
	return obj
}

func (b *Ext_java_util_BitSet) wordAt(w int) uint64 {
	if w >= int(b.n) {
		return 0
	}
	if b.extra != nil {
		return b.extra[w]
	}
	return b.buf[w]
}

func (b *Ext_java_util_BitSet) setWord(w int, v uint64) {
	if b.extra != nil {
		b.extra[w] = v
	} else {
		b.buf[w] = v
	}
}

func (b *Ext_java_util_BitSet) ensureWord(w int) {
	if w < int(b.n) {
		return
	}
	if b.extra == nil && w < len(b.buf) {
		b.n = int32(w + 1)
		return
	}
	if w < cap(b.extra) {
		b.extra = b.extra[:w+1]
		b.n = int32(w + 1)
		return
	}
	nw := make([]uint64, w+1, max(w+1, max(8, cap(b.extra)*2)))
	if b.extra != nil {
		copy(nw, b.extra[:b.n])
	} else {
		copy(nw, b.buf[:b.n])
	}
	b.extra = nw
	b.n = int32(w + 1)
}

func (b *Ext_java_util_BitSet) trim() {
	n := int(b.n)
	for n > 0 && b.wordAt(n-1) == 0 {
		n--
	}
	b.n = int32(n)
	if b.extra != nil {
		b.extra = b.extra[:n]
	}
}

func (b *Ext_java_util_BitSet) Set_1(bit int32) {
	w := int(bit >> 6)
	b.ensureWord(w)
	b.setWord(w, b.wordAt(w)|(1<<(uint(bit)&63)))
}

func (b *Ext_java_util_BitSet) Set_3(from, to int32, val bool) {
	for i := from; i < to; i++ {
		if val {
			b.Set_1(i)
		}
	}
}

func (b *Ext_java_util_BitSet) Get_1(bit int32) bool {
	w := int(bit >> 6)
	return (b.wordAt(w) & (1 << (uint(bit) & 63))) != 0
}

func (b *Ext_java_util_BitSet) IsEmpty() bool {
	return b.n == 0
}

func (b *Ext_java_util_BitSet) Size() int32 {
	if b.extra != nil {
		return int32(cap(b.extra) * 64)
	}
	return 256
}

func (b *Ext_java_util_BitSet) Or(other *Ext_java_util_BitSet) {
	if other.n > b.n {
		b.ensureWord(int(other.n) - 1)
	}
	for i := 0; i < int(other.n); i++ {
		b.setWord(i, b.wordAt(i)|other.wordAt(i))
	}
	b.trim()
}

func (b *Ext_java_util_BitSet) And(other *Ext_java_util_BitSet) {
	common := min(int(b.n), int(other.n))
	for i := 0; i < common; i++ {
		b.setWord(i, b.wordAt(i)&other.wordAt(i))
	}
	for i := common; i < int(b.n); i++ {
		b.setWord(i, 0)
	}
	b.n = int32(common)
	b.trim()
}

func (b *Ext_java_util_BitSet) AndNot(other *Ext_java_util_BitSet) {
	common := min(int(b.n), int(other.n))
	for i := 0; i < common; i++ {
		b.setWord(i, b.wordAt(i)&^other.wordAt(i))
	}
	b.trim()
}

func (b *Ext_java_util_BitSet) Xor(other *Ext_java_util_BitSet) {
	if other.n > b.n {
		b.ensureWord(int(other.n) - 1)
	}
	for i := 0; i < int(other.n); i++ {
		b.setWord(i, b.wordAt(i)^other.wordAt(i))
	}
	b.trim()
}

func (b *Ext_java_util_BitSet) Intersects(other *Ext_java_util_BitSet) bool {
	common := min(int(b.n), int(other.n))
	for i := 0; i < common; i++ {
		if (b.wordAt(i) & other.wordAt(i)) != 0 {
			return true
		}
	}
	return false
}

func (b *Ext_java_util_BitSet) Clone() any {
	cp := NewExt_java_util_BitSet_1(b.n * 64)
	cp.Or(b)
	return cp
}

type Ext_com_google_protobuf_CodedInputStream interface{}
type Ext_com_google_protobuf_ExtensionRegistryLite interface{}

var AstNodeReleaseCount int

type Ext_com_google_javascript_jscomp_serialization_AstNode struct {
	Kind     int32
	Children []*Ext_com_google_javascript_jscomp_serialization_AstNode
	Released bool
}

func (n *Ext_com_google_javascript_jscomp_serialization_AstNode) IsJavaNil_() bool { return n == nil }

func (n *Ext_com_google_javascript_jscomp_serialization_AstNode) Release_() {
	if n != nil && !n.Released {
		n.Released = true
		AstNodeReleaseCount++
	}
}

func Ext_com_google_javascript_jscomp_serialization_AstNode_ParseFrom_2_Ext_com_google_protobuf_CodedInputStream_Ext_com_google_protobuf_ExtensionRegistryLite(
	stream Ext_com_google_protobuf_CodedInputStream,
	reg Ext_com_google_protobuf_ExtensionRegistryLite,
) *Ext_com_google_javascript_jscomp_serialization_AstNode {
	_ = stream
	_ = reg
	return &Ext_com_google_javascript_jscomp_serialization_AstNode{
		Kind: 107,
		Children: []*Ext_com_google_javascript_jscomp_serialization_AstNode{
			{Kind: 1},
		},
	}
}
`

func writeWave18Workspace(t *testing.T, pkgName string, pkgFiles map[string]string) (tmpRoot, pkgDir string) {
	t.Helper()
	tmpRoot = t.TempDir()
	goMod := "module example.com/wave18test\n\ngo 1.23\n"
	if err := os.WriteFile(filepath.Join(tmpRoot, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	j2Dir := filepath.Join(tmpRoot, "j2gort")
	if err := os.MkdirAll(j2Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(j2Dir, "runtime.go"), []byte(wave18MockJ2gortSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	pkgDir = filepath.Join(tmpRoot, pkgName)
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range pkgFiles {
		if err := os.WriteFile(filepath.Join(pkgDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return tmpRoot, pkgDir
}

func runGoTestAndRaceInDir(t *testing.T, dir string) {
	t.Helper()
	// 1. Standard non-race build to measure exact 0-alloc escape analysis & sync.Pool behavior.
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test failed: %v\n%s", err, string(out))
	}
	// 2. Race-enabled build to verify 100% -race concurrency safety.
	raceCmd := exec.Command("go", "test", "-race", "-run", "Semantics|Race", "./...")
	raceCmd.Dir = dir
	if out, err := raceCmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race failed: %v\n%s", err, string(out))
	}
}

func TestWave18Task1AB_IfaceEqualsAndStringLiterals(t *testing.T) {
	src := `package sample

import "example.com/wave18test/j2gort"

type MyIface interface {
	IsJavaNil_() bool
	ID() int
}

type Impl struct {
	id int
}

func (x *Impl) IsJavaNil_() bool { return x == nil }
func (x *Impl) ID() int          { return x.id }

type Holder struct {
	cb  MyIface
	alt MyIface
}

var sideEffects int

func impureGet(x MyIface) MyIface {
	sideEffects++
	return x
}

func CheckNil(h *Holder, x MyIface, flag bool) (bool, bool, bool, bool, bool) {
	a := j2gort.IfaceEquals[MyIface](x, nil)
	b := j2gort.IfaceEquals[MyIface](nil, h.cb)
	c := !(j2gort.IfaceEquals[MyIface](h.cb, nil))
	d := flag && !j2gort.IfaceEquals[MyIface](x, nil)
	e := j2gort.IfaceEquals[MyIface](h.cb, h.alt)
	return a, b, c, d, e
}

func KeepImpure(arr []MyIface, x MyIface) (bool, bool) {
	return !j2gort.IfaceEquals[MyIface](impureGet(x), nil), j2gort.IfaceEquals[MyIface](arr[0], nil)
}

func CheckStrings(dyn string) (string, string, string, string, string, string) {
	s1 := j2gort.CheckNotNullString("alpha")
	s2 := j2gort.StrValString("beta")
	s3 := j2gort.StrValString(j2gort.StrValInt32(7))
	keepNull1 := j2gort.StrValString("\xff\x00j2go_null")
	keepNull2 := j2gort.StrValString(j2gort.NullString)
	keepDyn := j2gort.CheckNotNullString(dyn)
	return s1, s2, s3, keepNull1, keepNull2, keepDyn
}
`
	testSrc := `package sample

import "testing"

func TestSemantics(t *testing.T) {
	var typedNil *Impl
	var ifaceTypedNil MyIface = typedNil
	obj1 := &Impl{id: 1}
	obj2 := &Impl{id: 2}

	cases := []MyIface{nil, ifaceTypedNil, obj1, obj2}
	for _, lhs := range cases {
		for _, rhs := range cases {
			h := &Holder{cb: lhs, alt: rhs}
			a, b, c, d, e := CheckNil(h, lhs, true)
			wantNilX := lhs == nil || lhs.IsJavaNil_()
			wantNilCb := lhs == nil || lhs.IsJavaNil_()
			if a != wantNilX || b != wantNilCb || c != !wantNilCb || d != !wantNilX {
				t.Fatalf("nil check mismatch for (%#v, %#v): got (%v,%v,%v,%v)", lhs, rhs, a, b, c, d)
			}
			wantEq := (lhs == nil || lhs.IsJavaNil_()) && (rhs == nil || rhs.IsJavaNil_())
			if !wantEq && lhs != nil && !lhs.IsJavaNil_() && rhs != nil && !rhs.IsJavaNil_() {
				wantEq = lhs == rhs
			}
			if e != wantEq {
				t.Fatalf("2-arg IfaceEquals mismatch for (%#v, %#v): got %v, want %v", lhs, rhs, e, wantEq)
			}
		}
	}

	sideEffects = 0
	gotImpure, gotIdx := KeepImpure([]MyIface{ifaceTypedNil}, obj1)
	if sideEffects != 1 || !gotImpure || !gotIdx {
		t.Fatalf("impure check mismatch: sideEffects=%d gotImpure=%v gotIdx=%v", sideEffects, gotImpure, gotIdx)
	}

	s1, s2, s3, k1, k2, kd := CheckStrings("ok")
	if s1 != "alpha" || s2 != "beta" || s3 != "0" || k1 != "null" || k2 != "null" || kd != "ok" {
		t.Fatalf("string check mismatch: %q %q %q %q %q %q", s1, s2, s3, k1, k2, kd)
	}
}
`
	tmpRoot, pkgDir := writeWave18Workspace(t, "sample", map[string]string{
		"sample.go":      src,
		"sample_test.go": testSrc,
	})

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.ByPass["j2gort_peephole"] < 8 {
		t.Fatalf("expected >= 8 j2gort_peephole sites, got %+v", stats1.ByPass)
	}

	stats2, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 {
		t.Fatalf("expected 0 sites on Run 2 (idempotency), got %d (%+v)", stats2.TotalSites, stats2.ByPass)
	}

	outBytes, err := os.ReadFile(filepath.Join(pkgDir, "sample.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(outBytes)
	for _, want := range []string{
		`x == nil || x.IsJavaNil_()`,
		`h.cb == nil || h.cb.IsJavaNil_()`,
		`h.cb != nil && !h.cb.IsJavaNil_()`,
		`flag && (x != nil && !x.IsJavaNil_())`,
		`h.cb == h.alt || ((h.cb == nil || h.cb.IsJavaNil_()) && (h.alt == nil || h.alt.IsJavaNil_()))`,
		`s1 := "alpha"`,
		`s2 := "beta"`,
		`j2gort.IfaceEquals[MyIface](impureGet(x), nil)`,
		`j2gort.IfaceEquals[MyIface](arr[0], nil)`,
		`j2gort.StrValString("\xff\x00j2go_null")`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in lowered output:\n%s", want, got)
		}
	}

	runGoTestAndRaceInDir(t, tmpRoot)
}

func TestWave18Task1CD_StackBitSetAndAstNodeRelease(t *testing.T) {
	lvaSrc := `package analysis

import "example.com/wave18test/j2gort"

type Lattice struct {
	liveSet *j2gort.Ext_java_util_BitSet
}

type Analyzer struct {
	saved *j2gort.Ext_java_util_BitSet
}

func (this_ *Analyzer) FlowThrough(size int32, input *Lattice) *Lattice {
	var gen *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet_1(input.liveSet.Size())
	_ = gen
	var kill *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet_1(input.liveSet.Size())
	_ = kill
	this_.priv_computeGenKill(size, gen, kill)
	out := &Lattice{liveSet: j2gort.NewExt_java_util_BitSet_1(size)}
	out.liveSet.Or(input.liveSet)
	out.liveSet.AndNot(kill)
	out.liveSet.Or(gen)
	return out
}

func (this_ *Analyzer) priv_computeGenKill(n int32, gen *j2gort.Ext_java_util_BitSet, kill *j2gort.Ext_java_util_BitSet) {
	if n > 0 {
		this_.priv_computeGenKill(n-1, gen, kill)
		gen.Set_1(n - 1)
		kill.Set_1(0)
	}
}

func (this_ *Analyzer) GetSmallestCoveringSubtree(count int32, dep *j2gort.Ext_java_util_BitSet) bool {
	var candidates *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet_1(count)
	_ = candidates
	candidates.Set_3(0, count, true)
	candidates.And(dep)
	if candidates.IsEmpty() {
		j2gort.CheckStateFmt(false, "No common dependency for %s", j2gort.PtrToAny(dep))
	}
	return candidates.Get_1(0)
}

func (this_ *Analyzer) MoveDeclarationStatements(count int32, sym *j2gort.Ext_java_util_BitSet) bool {
	var chunksWithImmovableReferences *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet_1(count)
	_ = chunksWithImmovableReferences
	chunksWithImmovableReferences.Or(sym)
	return this_.GetSmallestCoveringSubtree(count, chunksWithImmovableReferences)
}

func ZeroArgLocalBitSet() bool {
	var used *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet()
	_ = used
	used.Set_1(5)
	return used.Get_1(5)
}

func MustNotLowerBitSets(a *Analyzer) *j2gort.Ext_java_util_BitSet {
	var cloned *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet()
	cloned.Set_1(1)
	_ = cloned.Clone()

	var fieldEsc *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet_1(16)
	a.saved = fieldEsc

	var closureEsc *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet_1(16)
	fn := func() { closureEsc.Set_1(2) }
	fn()

	var retEsc *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet_1(16)
	retEsc.Set_1(3)
	return retEsc
}
`
	deserSrc := `package analysis

import "example.com/wave18test/j2gort"

type DeserializerRunner struct{}

func (this_ *DeserializerRunner) RunInline(stream j2gort.Ext_com_google_protobuf_CodedInputStream, reg j2gort.Ext_com_google_protobuf_ExtensionRegistryLite) int32 {
	var scriptNode int32 = this_.priv_Runner_visit(j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_ParseFrom_2_Ext_com_google_protobuf_CodedInputStream_Ext_com_google_protobuf_ExtensionRegistryLite(stream, reg), 10)
	return scriptNode
}

func (this_ *DeserializerRunner) RunLocalVar(stream j2gort.Ext_com_google_protobuf_CodedInputStream, reg j2gort.Ext_com_google_protobuf_ExtensionRegistryLite) int32 {
	var proto *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode = j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_ParseFrom_2_Ext_com_google_protobuf_CodedInputStream_Ext_com_google_protobuf_ExtensionRegistryLite(stream, reg)
	_ = proto
	return this_.priv_Runner_visit(proto, 20)
}

func (this_ *DeserializerRunner) priv_Runner_visit(astNode *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode, base int32) int32 {
	if astNode.Released {
		panic("AstNode released prematurely before visit finished!")
	}
	return astNode.Kind + base
}
`
	testSrc := `package analysis

import (
	"testing"

	"example.com/wave18test/j2gort"
)

func TestBitSetZeroAllocsAndRelease(t *testing.T) {
	a := &Analyzer{}
	dep := j2gort.NewExt_java_util_BitSet_1(64)
	dep.Set_1(0)

	allocs := testing.AllocsPerRun(100, func() {
		if !a.GetSmallestCoveringSubtree(16, dep) {
			t.Fatal("expected true")
		}
		if !ZeroArgLocalBitSet() {
			t.Fatal("expected true")
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 heap allocations for stack-lowered local BitSets, got %v", allocs)
	}
	if !a.MoveDeclarationStatements(16, dep) {
		t.Fatal("expected true")
	}

	inLat := &Lattice{liveSet: j2gort.NewExt_java_util_BitSet_1(32)}
	inLat.liveSet.Set_1(0)
	outLat := a.FlowThrough(4, inLat)
	if !outLat.liveSet.Get_1(3) {
		t.Fatal("expected bit 3 set in FlowThrough output")
	}

	j2gort.AstNodeReleaseCount = 0
	r := &DeserializerRunner{}
	if got := r.RunInline(nil, nil); got != 117 {
		t.Fatalf("RunInline = %d, want 117", got)
	}
	if got := r.RunLocalVar(nil, nil); got != 127 {
		t.Fatalf("RunLocalVar = %d, want 127", got)
	}
	if j2gort.AstNodeReleaseCount != 2 {
		t.Fatalf("AstNodeReleaseCount = %d, want 2", j2gort.AstNodeReleaseCount)
	}
}
`
	tmpRoot, pkgDir := writeWave18Workspace(t, "analysis", map[string]string{
		"lva.go":      lvaSrc,
		"deser.go":    deserSrc,
		"lva_test.go": testSrc,
	})

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.ByPass["j2gort_peephole"] < 7 {
		t.Fatalf("expected >= 7 j2gort_peephole sites (5 BitSets + 2 AstNode.Release_), got %+v", stats1.ByPass)
	}

	stats2, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 {
		t.Fatalf("expected 0 sites on Run 2 (idempotency), got %d (%+v)", stats2.TotalSites, stats2.ByPass)
	}

	lvaBytes, err := os.ReadFile(filepath.Join(pkgDir, "lva.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotLva := string(lvaBytes)
	for _, want := range []string{
		`var gen_bsVal j2gort.Ext_java_util_BitSet`,
		`j2gort.InitExt_java_util_BitSet_1(&gen_bsVal, input.liveSet.Size())`,
		`var gen *j2gort.Ext_java_util_BitSet = &gen_bsVal`,
		`var kill_bsVal j2gort.Ext_java_util_BitSet`,
		`var candidates_bsVal j2gort.Ext_java_util_BitSet`,
		`var chunksWithImmovableReferences_bsVal j2gort.Ext_java_util_BitSet`,
		`var used_bsVal j2gort.Ext_java_util_BitSet`,
		`j2gort.InitExt_java_util_BitSet_0(&used_bsVal)`,
		`var cloned *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet()`,
		`var fieldEsc *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet_1(16)`,
		`var closureEsc *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet_1(16)`,
		`var retEsc *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet_1(16)`,
	} {
		if !strings.Contains(gotLva, want) {
			t.Errorf("missing expected snippet %q in lva.go:\n%s", want, gotLva)
		}
	}

	deserBytes, err := os.ReadFile(filepath.Join(pkgDir, "deser.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotDeser := string(deserBytes)
	for _, want := range []string{
		`defer astRoot_.Release_()`,
		`defer proto.Release_()`,
	} {
		if !strings.Contains(gotDeser, want) {
			t.Errorf("missing expected snippet %q in deser.go:\n%s", want, gotDeser)
		}
	}

	runGoTestAndRaceInDir(t, tmpRoot)
}

func TestWave18Task2_NodeTraversalLowering_PoolingAndZeroAllocs(t *testing.T) {
	rhinoSrc := `package rhino

type Node uint32
type InputId struct{}
`
	jscompTraversalSrc := `package jscomp

import (
	"example.com/wave18test/j2gort"
	pkg_rhino "example.com/wave18test/rhino"
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
}

type Scope struct {
	root   pkg_rhino.Node
	parent *Scope
}

func (s *Scope) IsJavaNil_() bool              { return s == nil }
func (s *Scope) GetRootNode_0() pkg_rhino.Node { return s.root }

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope_2_Node_AbstractScope(pkg_rhino.Node, AbstractScope) any
}

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

func (s *SyntacticScopeCreator) CreateScope_2_Node_AbstractScope(n pkg_rhino.Node, parent AbstractScope) any {
	var p *Scope
	if parent != nil && !parent.IsJavaNil_() {
		p = parent.(*Scope)
	}
	sc := &Scope{root: n, parent: p}
	if !s.scannerInUse {
		s.scannerInUse = true
		s.reusableScanner.scope = sc
		s.reusableScanner.Reset()
		s.scannerInUse = false
	}
	return sc
}

type NodeTraversal_Callback interface {
	IsJavaNil_() bool
	ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool
	Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type NodeTraversal_ScopedCallback interface {
	NodeTraversal_Callback
	EnterScope(*NodeTraversal)
	ExitScope(*NodeTraversal)
}

type NodeTraversal struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCallback                                  NodeTraversal_ScopedCallback
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
	mayContainSyntheticBlocks                      bool
	currentNode                                    pkg_rhino.Node
	currentHoistScopeRoot                          pkg_rhino.Node
	currentFunction                                pkg_rhino.Node
	currentScript                                  pkg_rhino.Node
	currentChangeScope                             pkg_rhino.Node
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
		var newLen int32
		if this_.scopeDepth == 0 {
			newLen = 8
		} else {
			newLen = int32(len(this_.scopeNodes)) << 1
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

// Simulate traversing a chain of nested scopes of length root (root..1).
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

func (this_ *NodeTraversal_Builder) Build() *NodeTraversal {
	return NewNodeTraversal(this_)
}

func (this_ *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverse_3817(root)
}

func (this_ *NodeTraversal_Builder) TraverseRoots(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.Build().priv_NodeTraversal_traverseRoots_3818(externs, root)
}

func (this_ *NodeTraversal_Builder) TraverseWithScope(root pkg_rhino.Node, s AbstractScope) {
	this_.Build().priv_NodeTraversal_traverseWithScope_3820(root, s)
}

func (this_ *NodeTraversal_Builder) TraverseAtScope(scope AbstractScope) {
	this_.Build().priv_NodeTraversal_traverseAtScope_3821(scope)
}
`
	jscompTestSrc := `package jscomp

import (
	"sync"
	"testing"

	pkg_rhino "example.com/wave18test/rhino"
)

type mockCompiler struct {
	opts *CompilerOptions
}

func (c *mockCompiler) IsJavaNil_() bool             { return c == nil }
func (c *mockCompiler) GetOptions() *CompilerOptions { return c.opts }

type countingCallback struct {
	visits    int
	entered   int
	exited    int
	needScope bool
}

func (cb *countingCallback) IsJavaNil_() bool { return cb == nil }
func (cb *countingCallback) ShouldTraverse_3(t *NodeTraversal, n, p pkg_rhino.Node) bool {
	return true
}
func (cb *countingCallback) Visit_3(t *NodeTraversal, n, p pkg_rhino.Node) {
	cb.visits++
	if cb.needScope {
		if sc := t.GetScope(); sc == nil || sc.GetRootNode_0() != t.GetScopeRoot() {
			panic("scope root mismatch")
		}
	}
}
func (cb *countingCallback) EnterScope(t *NodeTraversal) { cb.entered++ }
func (cb *countingCallback) ExitScope(t *NodeTraversal)  { cb.exited++ }

func TestNodeTraversalPoolingZeroAllocs(t *testing.T) {
	comp := &mockCompiler{opts: &CompilerOptions{}}
	cb := &countingCallback{}
	b := NodeTraversal_S_Builder().SetCompiler(comp).SetCallback_1_NodeTraversal_Callback(cb)

	// Warm up sync.Pool + defaultScopeCreator_.reusableScanner once.
	NodeTraversal_Traverse(comp, 10, cb)

	allocs := testing.AllocsPerRun(200, func() {
		cb.visits, cb.entered, cb.exited = 0, 0, 0
		NodeTraversal_Traverse(comp, 12, cb)
		if cb.visits != 12 || cb.entered != 12 || cb.exited != 12 {
			t.Fatalf("unexpected counts: %+v", cb)
		}
		b.Traverse(8)
	})
	if allocs != 0 {
		t.Fatalf("expected 0 heap allocs/op for pooled NodeTraversal up to depth 12, got %v", allocs)
	}
}

func TestNodeTraversalSemanticsAndRace(t *testing.T) {
	comp := &mockCompiler{opts: &CompilerOptions{}}
	b := NodeTraversal_S_Builder().SetCompiler(comp)

	// Test deep nesting > 12 (depth = 20) + lazy SyntacticScopeCreator reuse.
	cbScope := &countingCallback{needScope: true}
	b.SetCallback_1_NodeTraversal_Callback(cbScope)
	NodeTraversal_Traverse(comp, 20, cbScope)
	if cbScope.visits != 20 || cbScope.entered != 20 || cbScope.exited != 20 {
		t.Fatalf("deep traversal mismatch: %+v", cbScope)
	}
	NodeTraversal_TraverseRoots(comp, cbScope, 30, 5)
	s0 := &Scope{root: 99}
	NodeTraversal_TraverseAtScope(comp, cbScope, s0)
	b.TraverseWithScope(4, s0)

	// Verify -race concurrency safety across goroutines sharing nodeTraversalPool_.
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localComp := &mockCompiler{opts: &CompilerOptions{}}
			localCb := &countingCallback{needScope: true}
			for iter := 0; iter < 50; iter++ {
				NodeTraversal_Traverse(localComp, 15, localCb)
			}
		}()
	}
	wg.Wait()
}
`
	tmpRoot, jscompDir := writeWave18Workspace(t, "jscomp", map[string]string{
		"node_traversal.go":      jscompTraversalSrc,
		"node_traversal_test.go": jscompTestSrc,
	})
	rhinoDir := filepath.Join(tmpRoot, "rhino")
	if err := os.MkdirAll(rhinoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(rhinoDir, "node.go"), []byte(rhinoSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.ByPass["node_traversal_lowering"] != 11 {
		t.Fatalf("expected 11 node_traversal_lowering sites, got %d (%+v)", stats1.ByPass["node_traversal_lowering"], stats1.ByPass)
	}

	// Run 2 MUST report 0 sites modified across ALL passes (100% idempotency)!
	stats2, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 {
		t.Fatalf("expected 0 sites on Run 2 (idempotency), got %d (%+v)", stats2.TotalSites, stats2.ByPass)
	}

	ntBytes, err := os.ReadFile(filepath.Join(jscompDir, "node_traversal.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotNT := string(ntBytes)
	for _, want := range []string{
		`scopeNodesBuf_`,
		`scopeObjsBuf_`,
		`defaultScopeCreator_`,
		`this_.scopeNodesBuf_[:12]`,
		`this_.scopeObjsBuf_[:12]`,
		`InitSyntacticScopeCreator_1(&this_.defaultScopeCreator_, this_.compiler)`,
		`var nodeTraversalPool_ = sync.Pool{`,
		`func acquireNodeTraversal_(`,
		`func releaseNodeTraversal_(`,
		`defer releaseNodeTraversal_(t_)`,
	} {
		if !strings.Contains(gotNT, want) {
			t.Errorf("missing expected snippet %q in lowered node_traversal.go:\n%s", want, gotNT)
		}
	}

	runGoTestAndRaceInDir(t, tmpRoot)
}

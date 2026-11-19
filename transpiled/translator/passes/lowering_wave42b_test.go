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

func wave42bWriteFiles(t *testing.T, root string, files map[string]string) {
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

func wave42bRealGenReady(repoRoot string) bool {
	for _, pkg := range []string{"j2gort", "rhino"} {
		cmd := exec.Command("go", "build", "./transpiled/gen/"+pkg)
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			return false
		}
	}
	return true
}

// TestWave42B_AbstractBaseDevirtExplicitGetBaseAndElseIfGuard verifies Task 1:
//   - Explicit `<x>.Get<I>Base()` expressions (including `.typeKindMask` field reads and chained
//     method calls) dominated by `base_<x>_` are rewritten to `base_<x>_`.
//   - `if ... { return ... } else if (left == nil || left.IsJavaNil_()) || (right == nil || right.IsJavaNil_()) { return ... }`
//     hoists BOTH `base_left_` and `base_right_` immediately after the nil-guard IfStmt.
//   - Non-returning `if` branches do NOT prematurely hoist nullable variables (verified at runtime with nil).
//   - Upgrading already-partially-lowered v40 code also converges cleanly, and Run 2 produces 0 sites.
func TestWave42B_AbstractBaseDevirtExplicitGetBaseAndElseIfGuard(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	j2gortSrc := `package j2gort

const NullString = "\xff\x00j2go_null"

func CheckNotNullIface[T any](v T) T { return v }
`

	rhinoSrc := `package rhino

import "sync/atomic"

var GetBaseCalls atomic.Int64

type JSType interface {
	Implements_rhino_JSType()
	IsJavaNil_() bool
	GetJSTypeBase() *JSTypeBase
	IsUnknownType() bool
	IsAllType() bool
	IsNoType() bool
	IsNoResolvedType() bool
	IsUnionType() bool
	GetReferenceTag() string
}

type JSTypeBase struct {
	self_          JSType
	typeKindMask   int32
	unknown        bool
	noResolved     bool
}

func (b *JSTypeBase) BindSelf_(s any)           { b.self_ = s.(JSType) }
func (b *JSTypeBase) IsJavaNil_() bool          { return b == nil }
func (*JSTypeBase) Implements_rhino_JSType()    {}
func (b *JSTypeBase) GetJSTypeBase() *JSTypeBase {
	GetBaseCalls.Add(1)
	return b
}
func (b *JSTypeBase) IsUnknownType() bool    { return b.unknown }
func (b *JSTypeBase) IsAllType() bool        { return (b.typeKindMask & 16) != 0 }
func (b *JSTypeBase) IsNoType() bool         { return (b.typeKindMask & 2) != 0 }
func (b *JSTypeBase) IsNoResolvedType() bool { return b.noResolved }
func (b *JSTypeBase) IsUnionType() bool      { return (b.typeKindMask & 64) != 0 }
func (b *JSTypeBase) GetReferenceTag() string { return "base" }

type NamedObjType struct {
	JSTypeBase
	tag string
}

func NewNamedObjType(mask int32, unk, noRes bool, tag string) *NamedObjType {
	t := &NamedObjType{tag: tag}
	t.typeKindMask = mask
	t.unknown = unk
	t.noResolved = noRes
	t.BindSelf_(t)
	return t
}
func (t *NamedObjType) IsJavaNil_() bool        { return t == nil }
func (*NamedObjType) Implements_rhino_JSType()  {}
func (t *NamedObjType) GetReferenceTag() string { return t.tag }

// Fresh SubtypeChecker.isSubtypeHelper pattern with explicit .GetJSTypeBase().typeKindMask
func IsSubtypeHelperFresh(subtype JSType, supertype JSType) bool {
	if subtype == supertype || ((subtype == nil || subtype.IsJavaNil_()) && (supertype == nil || supertype.IsJavaNil_())) {
		return true
	}
	var subMask int32 = subtype.GetJSTypeBase().typeKindMask
	var supMask int32 = supertype.GetJSTypeBase().typeKindMask
	if ((supertype.IsUnknownType() || (((supMask & 17) != 0) && supertype.IsAllType())) || subtype.IsUnknownType()) || (((subMask & 3) != 0) && subtype.IsNoType()) {
		return true
	}
	return subtype.IsUnionType() && supertype.IsUnionType()
}

// Fresh EqualityChecker.areEqualInternal pattern with if / else if nil guard
func AreEqualInternalFresh(left JSType, right JSType) bool {
	if left == right || ((left == nil || left.IsJavaNil_()) && (right == nil || right.IsJavaNil_())) {
		return true
	} else if (left == nil || left.IsJavaNil_()) || (right == nil || right.IsJavaNil_()) {
		return false
	}
	if left.IsNoResolvedType() && right.IsNoResolvedType() {
		return left.GetReferenceTag() == right.GetReferenceTag()
	}
	leftUnk := left.IsUnknownType()
	rightUnk := right.IsUnknownType()
	if leftUnk || rightUnk {
		return leftUnk && rightUnk
	}
	return left.IsUnionType() && right.IsUnionType()
}

// Already-partially-lowered v40 EqualityChecker pattern (upgrades cleanly!)
func AreEqualInternalV40Upgrade(left JSType, right JSType) bool {
	if left == right || ((left == nil || left.IsJavaNil_()) && (right == nil || right.IsJavaNil_())) {
		return true
	} else if (left == nil || left.IsJavaNil_()) || (right == nil || right.IsJavaNil_()) {
		return false
	}
	base_left_ := left.GetJSTypeBase()
	if base_left_.IsNoResolvedType() && right.GetJSTypeBase().IsNoResolvedType() {
		return left.GetReferenceTag() == right.GetReferenceTag()
	}
	leftUnk := base_left_.IsUnknownType()
	base_right_ := right.GetJSTypeBase()
	rightUnk := base_right_.IsUnknownType()
	return leftUnk == rightUnk && (base_left_.typeKindMask == right.GetJSTypeBase().typeKindMask)
}

// Safety: non-returning first branch MUST NOT allow else-if nil guard to hoist nullable x!
func NonReturningFirstBranchMustNotHoist(skipGuard bool, x JSType) int32 {
	side := int32(0)
	if skipGuard {
		side = 7
	} else if x == nil || x.IsJavaNil_() {
		return -1
	}
	if skipGuard {
		return side
	}
	if x.IsUnknownType() {
		return 1
	}
	return x.GetJSTypeBase().typeKindMask
}

// Single explicit GetJSTypeBase() must NOT be double-rewritten or trigger non-zero sites.
func SingleExplicitFieldRead(x JSType) int32 {
	return x.GetJSTypeBase().typeKindMask
}
`

	rtTestSrc := `package rhino

import "testing"

func TestDevirtRuntimeAndCallCounts(t *testing.T) {
	a := NewNamedObjType(16, false, true, "RefA")
	b := NewNamedObjType(16, false, true, "RefA")

	GetBaseCalls.Store(0)
	if !IsSubtypeHelperFresh(a, b) {
		t.Fatalf("expected IsSubtypeHelperFresh(a, b) == true")
	}
	if got := GetBaseCalls.Load(); got != 2 {
		t.Fatalf("expected exactly 2 GetJSTypeBase calls (1 subtype + 1 supertype), got %d", got)
	}

	GetBaseCalls.Store(0)
	if !AreEqualInternalFresh(a, b) {
		t.Fatalf("expected AreEqualInternalFresh(a, b) == true")
	}
	if got := GetBaseCalls.Load(); got != 2 {
		t.Fatalf("expected exactly 2 GetJSTypeBase calls in AreEqualInternalFresh, got %d", got)
	}

	GetBaseCalls.Store(0)
	if !AreEqualInternalV40Upgrade(a, b) {
		t.Fatalf("expected AreEqualInternalV40Upgrade(a, b) == true")
	}
	if got := GetBaseCalls.Load(); got != 2 {
		t.Fatalf("expected exactly 2 GetJSTypeBase calls in upgraded v40 AreEqualInternal, got %d", got)
	}

	// Nil safety checks!
	if AreEqualInternalFresh(a, nil) || AreEqualInternalFresh(nil, b) || !AreEqualInternalFresh(nil, nil) {
		t.Fatalf("unexpected nil equality behavior")
	}
	var typedNil *NamedObjType
	if AreEqualInternalFresh(a, typedNil) || !AreEqualInternalFresh(typedNil, nil) {
		t.Fatalf("unexpected typed-nil equality behavior")
	}
	if got := NonReturningFirstBranchMustNotHoist(true, nil); got != 7 {
		t.Fatalf("expected NonReturningFirstBranchMustNotHoist(true, nil) == 7 without panic, got %d", got)
	}
}
`

	wave42bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                 "module example.com/w42b_devirt\n\ngo 1.24\n",
		"gen/j2gort/j2gort.go":   j2gortSrc,
		"gen/rhino/rhino.go":     rhinoSrc,
	})

	rhinoDir := filepath.Join(tmpDir, "gen/rhino")
	allPasses := passes.AllPasses()

	s1, err := passes.RunOnPackageDir(rhinoDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.ByPass["abstract_base_devirt"] < 12 {
		t.Fatalf("expected >= 12 abstract_base_devirt sites on Run 1, got %+v", s1.ByPass)
	}

	s2, err := passes.RunOnPackageDir(rhinoDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: %+v", s2)
	}

	codeBytes, err := os.ReadFile(filepath.Join(rhinoDir, "rhino.go"))
	if err != nil {
		t.Fatal(err)
	}
	code := string(codeBytes)
	for _, want := range []string{
		"var subMask int32 = base_subtype_.typeKindMask",
		"var supMask int32 = base_supertype_.typeKindMask",
		"if base_left_.IsNoResolvedType() && base_right_.IsNoResolvedType()",
		"base_left_.typeKindMask == base_right_.typeKindMask",
	} {
		if !strings.Contains(code, want) {
			t.Fatalf("expected lowered rhino.go to contain %q:\n%s", want, code)
		}
	}

	wave42bWriteFiles(t, tmpDir, map[string]string{
		"gen/rhino/devirt_rt_test.go": rtTestSrc,
	})
	cmd := exec.Command("go", "test", "-race", "-v", "./gen/rhino")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race ./gen/rhino failed: %v\n%s", err, string(out))
	}
}

// TestWave42B_CallArgIIFEUnwrappingAndOrderSafety verifies Task 2:
//   - Call-argument zero-arg single-result IIFEs unwrap cleanly in ExprStmt, DeclStmt, AssignStmt,
//     ReturnStmt, multi-IIFE argument lists, nested pure wrapper calls, and goto-containing funcs.
//   - Impure preceding expressions (e.g. stateful builder chains) and receiver-mutating IIFEs are
//     NOT unwrapped prematurely, strictly preserving Go left-to-right evaluation order.
//   - Run 2 produces strictly 0 sites and 0 modified files.
func TestWave42B_CallArgIIFEUnwrappingAndOrderSafety(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	j2gortSrc := `package j2gort

const NullString = "\xff\x00j2go_null"

func StrValString(s string) string {
	if s == NullString {
		return "null"
	}
	return s
}
`

	serSrc := `package jscomp

import "example.com/w42b_iife/gen/j2gort"

type StringPoolBuilder struct {
	items []string
}

func (p *StringPoolBuilder) Put(s string) int32 {
	p.items = append(p.items, s)
	return int32(len(p.items))
}

type AstBuilder struct {
	origPtr int32
	lineCnt int32
	rawPtr  int32
	cookPtr int32
}

func (b *AstBuilder) SetOriginalNamePointer(v int32) *AstBuilder {
	b.origPtr = v
	return b
}
func (b *AstBuilder) SetLineCount(v int32) *AstBuilder {
	b.lineCnt = v
	return b
}
func (b *AstBuilder) SetRawStringPointer(v int32) *AstBuilder {
	b.rawPtr = v
	return b
}
func (b *AstBuilder) SetCookedStringPointer(v int32) *AstBuilder {
	b.cookPtr = v
	return b
}

type Serializer struct {
	pool *StringPoolBuilder
}

func (this_ *Serializer) VisitNode(builder *AstBuilder, origNameId int32) string {
	// 1. ExprStmt call-arg IIFE (TypedAstSerializer.visit line 264)
	builder.SetOriginalNamePointer(func() int32 {
		if origNameId <= int32(1) {
			return int32(0)
		}
		return this_.pool.Put("orig")
	}())

	// 2. DeclStmt + nested pure wrapper call-arg IIFE
	var norm string = j2gort.StrValString(func() string {
		if origNameId < 0 {
			return j2gort.NullString
		}
		return "ok"
	}())

	// 3. AssignStmt (:=) with TWO sibling IIFE arguments + &local prefix arg!
	var pair Pair
	tag := initPair(&pair, origNameId, func() int32 {
		if origNameId == 0 {
			return 10
		}
		return 20
	}(), func() int32 {
		if origNameId == 0 {
			return 30
		}
		return 40
	}())

	// 4. ReturnStmt call-arg IIFE
	return formatOut(norm, tag, func() int32 {
		if builder.origPtr > 0 {
			return builder.origPtr + 100
		}
		return 5
	}())
}

type Pair struct{ a, b int32 }

func initPair(dst *Pair, base, x, y int32) int32 {
	dst.a = x
	dst.b = y
	return base + x + y
}

func formatOut(s string, a, b int32) string {
	if a > b {
		return s + ":gt"
	}
	return s + ":le"
}

// 5. Function with goto: ExprStmt call-arg IIFE must wrap in block scope so goto never jumps over var!
func VisitWithGoto(builder *AstBuilder, skip bool, v int32) int32 {
	if skip {
		goto done
	}
	builder.SetLineCount(func() int32 {
		if v < 0 {
			return 0
		}
		return v + 1
	}())
done:
	return builder.lineCnt
}

// 6. SAFETY: Preceding stateful call SetRawStringPointer(pool.Put("raw")) MUST evaluate BEFORE cooked IIFE!
func ChainedStatefulMustPreserveOrder(b *AstBuilder, pool *StringPoolBuilder, hasCooked bool) {
	b.SetRawStringPointer(pool.Put("raw")).SetCookedStringPointer(func() int32 {
		if !hasCooked {
			return -1
		}
		return pool.Put("cooked")
	}())
}

// 7. SAFETY: Callee expression with function call getBuilder() MUST evaluate BEFORE argument IIFE,
// and variable-mutating IIFE must ALSO remain untouched!
func CalleeCallAndMutatedVarMustPreserveOrder() (string, int32, int32) {
	var trace string
	b1 := &AstBuilder{}
	b2 := &AstBuilder{}
	getBuilder := func() *AstBuilder {
		trace += "callee->"
		return b1
	}
	getBuilder().SetLineCount(func() int32 {
		trace += "arg"
		return 77
	}())
	target := b1
	target.SetOriginalNamePointer(func() int32 {
		target = b2
		return 99
	}())
	return trace, b1.lineCnt, b2.origPtr
}
`

	rtTestSrc := `package jscomp

import "testing"

func TestCallArgIIFERuntimeEquivalence(t *testing.T) {
	pool := &StringPoolBuilder{}
	s := &Serializer{pool: pool}
	b := &AstBuilder{}

	got1 := s.VisitNode(b, 1)
	if b.origPtr != 0 || got1 != "ok:gt" {
		t.Fatalf("unexpected VisitNode(1): origPtr=%d got=%q", b.origPtr, got1)
	}
	got2 := s.VisitNode(b, 5)
	if b.origPtr != 1 || got2 != "ok:le" {
		t.Fatalf("unexpected VisitNode(5): origPtr=%d got=%q", b.origPtr, got2)
	}

	bGoto := &AstBuilder{}
	if res := VisitWithGoto(bGoto, true, 10); res != 0 {
		t.Fatalf("expected skipped goto result 0, got %d", res)
	}
	if res := VisitWithGoto(bGoto, false, 10); res != 11 {
		t.Fatalf("expected non-skipped goto result 11, got %d", res)
	}

	// Verify stateful chain order: "raw" MUST get ID 1 and "cooked" MUST get ID 2!
	orderPool := &StringPoolBuilder{}
	bChain := &AstBuilder{}
	ChainedStatefulMustPreserveOrder(bChain, orderPool, true)
	if bChain.rawPtr != 1 || bChain.cookPtr != 2 {
		t.Fatalf("evaluation order violated! rawPtr=%d cookPtr=%d items=%v", bChain.rawPtr, bChain.cookPtr, orderPool.items)
	}

	// Verify callee call runs BEFORE arg IIFE ("callee->arg") and mutated receiver IIFE stays intact!
	trace, c1, c2 := CalleeCallAndMutatedVarMustPreserveOrder()
	if trace != "callee->arg" || c1 != 77 || c2 != 99 {
		t.Fatalf("evaluation order violated! trace=%q c1=%d c2=%d", trace, c1, c2)
	}
}
`

	wave42bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                   "module example.com/w42b_iife\n\ngo 1.24\n",
		"gen/j2gort/j2gort.go":     j2gortSrc,
		"gen/jscomp/serializer.go": serSrc,
	})

	jscompDir := filepath.Join(tmpDir, "gen/jscomp")
	allPasses := passes.AllPasses()

	s1, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.ByPass["j2gort_peephole"] < 6 {
		t.Fatalf("expected >= 6 j2gort_peephole sites on Run 1, got %+v", s1.ByPass)
	}

	s2, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: %+v", s2)
	}

	serBytes, err := os.ReadFile(filepath.Join(jscompDir, "serializer.go"))
	if err != nil {
		t.Fatal(err)
	}
	serCode := string(serBytes)
	if strings.Contains(serCode, "builder.SetOriginalNamePointer(func()") {
		t.Fatalf("expected builder.SetOriginalNamePointer IIFE to be unwrapped:\n%s", serCode)
	}
	if !strings.Contains(serCode, "builder.SetOriginalNamePointer(iifeArg_") {
		t.Fatalf("expected builder.SetOriginalNamePointer(iifeArg_*) in:\n%s", serCode)
	}
	if !strings.Contains(serCode, "target.SetOriginalNamePointer(func()") {
		t.Fatalf("mutated receiver target.SetOriginalNamePointer(func()...) MUST remain untouched:\n%s", serCode)
	}
	if strings.Contains(serCode, "initPair(&pair, origNameId, func()") {
		t.Fatalf("expected both initPair IIFE args to be unwrapped:\n%s", serCode)
	}
	// Unsafe sites MUST remain untouched as IIFEs!
	if !strings.Contains(serCode, "SetCookedStringPointer(func()") {
		t.Fatalf("stateful chain SetCookedStringPointer(func()...) MUST NOT be hoisted:\n%s", serCode)
	}

	wave42bWriteFiles(t, tmpDir, map[string]string{
		"gen/jscomp/ser_rt_test.go": rtTestSrc,
	})
	cmd := exec.Command("go", "test", "-race", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race ./gen/jscomp failed: %v\n%s", err, string(out))
	}
}

// TestWave42B_OptHotspotsInlineVariablesInlineFunctionsRemoveUnusedCode verifies Task 3:
//   1) InlineVariables CreateLazyNegative avoids promoteIvScratchVarExpert_ and folds known constants.
//   2) InlineFunctions pools FunctionState onto PooledSlab, pools FunctionInjector.callerFnStateCache,
//      and promotes >256-entry FunctionState.references maps into the largeFree OrderedMap pool tier.
//   3) RemoveUnusedCode replaces BindSelf_(any) with static this_.self_ = this_ (inlining allocPooled_*),
//      strips ToIface+CastPtr round-trips in Build*, guards zero-continuation writes, and tolerates
//      RemoveUnusedCode_Continuation BOTH with and without outer_!
func TestWave42B_OptHotspotsInlineVariablesInlineFunctionsRemoveUnusedCode(t *testing.T) {
	t.Parallel()

	for _, staticCont := range []bool{false, true} {
		name := "NonStaticContinuation"
		if staticCont {
			name = "StaticContinuationNoOuter"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()

			j2gortSrc := `package j2gort

import "sync/atomic"

const NullString = "\xff\x00j2go_null"

func CheckState(b bool) {
	if !b {
		panic("CheckState")
	}
}
func CheckArgument(b bool) {
	if !b {
		panic("CheckArgument")
	}
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

type Slab[T any] struct{ items []T }
func (s *Slab[T]) Alloc() *T {
	s.items = append(s.items, *new(T))
	return &s.items[len(s.items)-1]
}
func (s *Slab[T]) Release() { clear(s.items); s.items = nil }

type PooledSlab[T any] struct {
	items      []T
	AllocCalls atomic.Int64
	Releases   atomic.Int64
}
func (s *PooledSlab[T]) Alloc() *T {
	s.AllocCalls.Add(1)
	s.items = append(s.items, *new(T))
	return &s.items[len(s.items)-1]
}
func (s *PooledSlab[T]) Reset() { clear(s.items); s.items = s.items[:0] }
func (s *PooledSlab[T]) Release() {
	s.Releases.Add(1)
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

var (
	AcquireMapCapCalls    atomic.Int64
	ReleaseMapToPoolCalls atomic.Int64
	LargePtrMapPromotes   atomic.Int64
)

type OrderedMapEntry[K comparable, V any] struct {
	Key K
	Val V
}

type OrderedMap[K comparable, V any] struct {
	entries []OrderedMapEntry[K, V]
	idx     map[K]int
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V]        { return OrderedMap[K, V]{} }
func NewOrderedMapCap[K comparable, V any](_ int32) OrderedMap[K, V] { return OrderedMap[K, V]{} }
func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V]    { return &OrderedMap[K, V]{} }
func NewOrderedMapPtrCap[K comparable, V any](_ int32) *OrderedMap[K, V] { return &OrderedMap[K, V]{} }

func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] { return OrderedMap[K, V]{} }
func AcquireOrderedMapCap[K comparable, V any](c int32) OrderedMap[K, V] {
	AcquireMapCapCalls.Add(1)
	return OrderedMap[K, V]{idx: make(map[K]int, c)}
}
func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] { return &OrderedMap[K, V]{} }
func AcquireOrderedMapPtrCap[K comparable, V any](c int32) *OrderedMap[K, V] {
	if c >= 512 {
		LargePtrMapPromotes.Add(1)
	}
	return &OrderedMap[K, V]{idx: make(map[K]int, c)}
}
func ReleaseOrderedMap[K comparable, V any](m *OrderedMap[K, V]) {
	if m != nil {
		m.Clear()
	}
}
func (m *OrderedMap[K, V]) ReleaseToPool() {
	ReleaseMapToPoolCalls.Add(1)
	m.Clear()
}
func (m *OrderedMap[K, V]) Clear() {
	m.entries = nil
	m.idx = nil
}
func (m *OrderedMap[K, V]) Size() int32    { return int32(len(m.entries)) }
func (m *OrderedMap[K, V]) IsEmpty() bool  { return len(m.entries) == 0 }
func (m *OrderedMap[K, V]) ContainsKey(k K) bool {
	if m.idx == nil {
		return false
	}
	_, ok := m.idx[k]
	return ok
}
func (m *OrderedMap[K, V]) Get(k K) V {
	var zero V
	if m.idx == nil {
		return zero
	}
	if i, ok := m.idx[k]; ok {
		return m.entries[i].Val
	}
	return zero
}
func (m *OrderedMap[K, V]) Put(k K, v V) V {
	var zero V
	if m.idx == nil {
		m.idx = make(map[K]int)
	}
	if i, ok := m.idx[k]; ok {
		old := m.entries[i].Val
		m.entries[i].Val = v
		return old
	}
	m.idx[k] = len(m.entries)
	m.entries = append(m.entries, OrderedMapEntry[K, V]{Key: k, Val: v})
	return zero
}
func (m *OrderedMap[K, V]) PutAllMapPtr(other *OrderedMap[K, V]) {
	if other == nil {
		return
	}
	for _, e := range other.entries {
		m.Put(e.Key, e.Val)
	}
}
func (m *OrderedMap[K, V]) Remove(k K) V {
	var zero V
	if m.idx == nil {
		return zero
	}
	i, ok := m.idx[k]
	if !ok {
		return zero
	}
	old := m.entries[i].Val
	m.entries = append(m.entries[:i], m.entries[i+1:]...)
	delete(m.idx, k)
	for j := i; j < len(m.entries); j++ {
		m.idx[m.entries[j].Key] = j
	}
	return old
}

type MapValCursor[V any] struct {
	s []V
	i int
}
func (c *MapValCursor[V]) HasNext() bool { return c.i < len(c.s) }
func (c *MapValCursor[V]) Next() V       { v := c.s[c.i]; c.i++; return v }

func (m *OrderedMap[K, V]) IterValues() MapValCursor[V] {
	vals := make([]V, len(m.entries))
	for i, e := range m.entries {
		vals[i] = e.Val
	}
	return MapValCursor[V]{s: vals}
}
`

			rhinoSrc := `package rhino

type Node uint32

func (n Node) IsName() bool     { return n != 0 }
func (n Node) IsFunction() bool { return n != 0 }
`

			inlineVarsSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w42b_opt/gen/j2gort"
	pkg_rhino "example.com/w42b_opt/gen/rhino"
)

type AbstractCompiler interface{ IsCompiler() }
type Scope struct{}

type Var struct {
	name  string
	kind  int // 0=inlineWellDef, 1=lazyNegKnownFalse, 2=lazyNegUnknown, 3=waitAlias
	alias *Var
}

type Reference struct {
	node pkg_rhino.Node
}

type ReferenceCollection struct {
	refs         []*Reference
	wellDef      bool
	assignedOnce bool
	WellDefCalls int
}

func (rc *ReferenceCollection) Get(i int32) *Reference { return rc.refs[i] }
func (rc *ReferenceCollection) Size() int32            { return int32(len(rc.refs)) }
func (rc *ReferenceCollection) IsWellDefined() bool {
	rc.WellDefCalls++
	return rc.wellDef
}
func (rc *ReferenceCollection) IsAssignedOnceInLifetime() bool {
	return rc.assignedOnce
}

type InlineVariables struct {
	compiler     AbstractCompiler
	vars         []*Var
	cols         map[*Var]*ReferenceCollection
	inlinedNames []string
}

func (*InlineVariables) Implements_jscomp_CompilerPass() {}

var (
	static_NO_INLINE_UNSAFE = NewInlineVariables_InlineVarAnalysis(0, nil, nil, 0, nil, nil)
	static_NO_INLINE_SAFE   = NewInlineVariables_InlineVarAnalysis(1, nil, nil, 0, nil, nil)
	static_NO_EXPERT        = NewInlineVariables_VarExpert_1(static_NO_INLINE_UNSAFE)
)

func StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS() *InlineVariables_InlineVarAnalysis {
	return static_NO_INLINE_UNSAFE
}
func StaticGet_InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS() *InlineVariables_InlineVarAnalysis {
	return static_NO_INLINE_SAFE
}
func StaticGet_NO_EXPERT() *InlineVariables_VarExpert { return static_NO_EXPERT }

func (this_ *InlineVariables) Process_2(_ pkg_rhino.Node, _ pkg_rhino.Node) {
	beh := NewInlineVariables_InliningBehavior(this_)
	beh.priv_InliningBehavior_doInlinesForScope_1()
}

type InlineVariables_InlineVarAnalysis struct {
	aliasedVar          *Var
	expert              *InlineVariables_VarExpert
	initialization      *Reference
	singleReadReference *Reference
	initValue           pkg_rhino.Node
	idHash_             atomic.Int32
	kind                int8
}

func (this_ *InlineVariables_InlineVarAnalysis) initFields_InlineVariables_InlineVarAnalysis() {}

func NewInlineVariables_InlineVarAnalysis(kind int8, aliasedVar *Var, expert *InlineVariables_VarExpert, initValue pkg_rhino.Node, initialization *Reference, singleReadReference *Reference) *InlineVariables_InlineVarAnalysis {
	this_ := &InlineVariables_InlineVarAnalysis{}
	InitInlineVariables_InlineVarAnalysis(this_, kind, aliasedVar, expert, initValue, initialization, singleReadReference)
	return this_
}

func InitInlineVariables_InlineVarAnalysis(this_ *InlineVariables_InlineVarAnalysis, kind int8, aliasedVar *Var, expert *InlineVariables_VarExpert, initValue pkg_rhino.Node, initialization *Reference, singleReadReference *Reference) {
	this_.initFields_InlineVariables_InlineVarAnalysis()
	this_.kind = kind
	this_.aliasedVar = aliasedVar
	this_.expert = expert
	this_.initValue = initValue
	this_.initialization = initialization
	this_.singleReadReference = singleReadReference
}

func InlineVariables_InlineVarAnalysis_CreateLazyNegative(expert *InlineVariables_VarExpert) *InlineVariables_InlineVarAnalysis {
	return NewInlineVariables_InlineVarAnalysis(int8(2), nil, expert, 0, nil, nil)
}
func InlineVariables_InlineVarAnalysis_CreateWaitForAlias(aliasedVar *Var) *InlineVariables_InlineVarAnalysis {
	return NewInlineVariables_InlineVarAnalysis(int8(3), aliasedVar, nil, 0, nil, nil)
}
func InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(expert *InlineVariables_VarExpert, initValue pkg_rhino.Node) *InlineVariables_InlineVarAnalysis {
	return NewInlineVariables_InlineVarAnalysis(int8(5), nil, expert, initValue, nil, nil)
}

func (this_ *InlineVariables_InlineVarAnalysis) ShouldInline() bool             { return this_.kind >= 4 }
func (this_ *InlineVariables_InlineVarAnalysis) ShouldWaitForAliasedVar() bool  { return this_.kind == 3 }
func (this_ *InlineVariables_InlineVarAnalysis) GetAliasedVar() *Var            { return this_.aliasedVar }
func (this_ *InlineVariables_InlineVarAnalysis) IsSafeToInlineAliases() bool {
	switch this_.kind {
	case int8(0):
		return false
	case int8(1), int8(4), int8(5), int8(6):
		return true
	case int8(2):
		return this_.expert.IsWellDefinedAssignedOnce()
	default:
		panic("bad kind")
	}
}
func (this_ *InlineVariables_InlineVarAnalysis) PerformInline() {
	if this_.kind == 5 {
		this_.expert.InlineWellDefined(this_.initValue)
	}
}

type InlineVariables_VarExpert struct {
	fixedAnalysis                       *InlineVariables_InlineVarAnalysis
	behavior                            *InlineVariables_InliningBehavior
	v                                   *Var
	referenceInfo                       *ReferenceCollection
	declaration                         *Reference
	initializationReference             *Reference
	idHash_                             atomic.Int32
	isDeclaredOrInferredConstant        bool
	mayBeAParameterModifiedViaArguments bool
	isNeverAssigned                     int8
	isWellDefined                       int8
	isAssignedOnceInLifetime            int8
	isReferencedWeakly                  int8
	isWellDefinedAssignedOnce           int8
	initializationReferenceKnown        bool
	allReferencesAreValid               int8
}

func (this_ *InlineVariables_VarExpert) initFields_InlineVariables_VarExpert() {
	this_.isNeverAssigned = -1
	this_.isWellDefined = -1
	this_.isAssignedOnceInLifetime = -1
	this_.isReferencedWeakly = -1
	this_.isWellDefinedAssignedOnce = -1
	this_.allReferencesAreValid = -1
}

func NewInlineVariables_VarExpert_1(fixed *InlineVariables_InlineVarAnalysis) *InlineVariables_VarExpert {
	this_ := &InlineVariables_VarExpert{}
	InitInlineVariables_VarExpert_1(this_, fixed)
	return this_
}
func InitInlineVariables_VarExpert_1(this_ *InlineVariables_VarExpert, fixed *InlineVariables_InlineVarAnalysis) {
	this_.initFields_InlineVariables_VarExpert()
	this_.fixedAnalysis = fixed
}

func NewInlineVariables_VarExpert_5(behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, isConst bool, mayMod bool) *InlineVariables_VarExpert {
	this_ := &InlineVariables_VarExpert{}
	InitInlineVariables_VarExpert_5(this_, behavior, v, referenceInfo, isConst, mayMod)
	return this_
}
func InitInlineVariables_VarExpert_5(this_ *InlineVariables_VarExpert, behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, isConst bool, mayMod bool) {
	this_.initFields_InlineVariables_VarExpert()
	this_.fixedAnalysis = nil
	this_.behavior = behavior
	this_.v = v
	this_.referenceInfo = referenceInfo
	this_.declaration = referenceInfo.Get(0)
	this_.isDeclaredOrInferredConstant = isConst
	this_.mayBeAParameterModifiedViaArguments = mayMod
}

func (this_ *InlineVariables_VarExpert) priv_VarExpert_isWellDefined_1() bool {
	if this_.isWellDefined < 0 {
		if this_.referenceInfo.IsWellDefined() {
			this_.isWellDefined = 1
		} else {
			this_.isWellDefined = 0
		}
	}
	return this_.isWellDefined == 1
}

func (this_ *InlineVariables_VarExpert) priv_VarExpert_isAssignedOnceInLifetime_2() bool {
	if this_.isAssignedOnceInLifetime < 0 {
		if this_.referenceInfo.IsAssignedOnceInLifetime() {
			this_.isAssignedOnceInLifetime = 1
		} else {
			this_.isAssignedOnceInLifetime = 0
		}
	}
	return this_.isAssignedOnceInLifetime == 1
}

func (this_ *InlineVariables_VarExpert) IsWellDefinedAssignedOnce() bool {
	if this_.isWellDefinedAssignedOnce < 0 {
		if this_.isWellDefined != 0 && this_.isAssignedOnceInLifetime != 0 &&
			this_.priv_VarExpert_isWellDefined_1() && this_.priv_VarExpert_isAssignedOnceInLifetime_2() {
			this_.isWellDefinedAssignedOnce = 1
		} else {
			this_.isWellDefinedAssignedOnce = 0
		}
	}
	return this_.isWellDefinedAssignedOnce == 1
}

func (this_ *InlineVariables_VarExpert) InlineWellDefined(_ pkg_rhino.Node) {
	this_.behavior.outer_.inlinedNames = append(this_.behavior.outer_.inlinedNames, this_.v.name)
}

func (this_ *InlineVariables_VarExpert) Analyze() *InlineVariables_InlineVarAnalysis {
	if this_.fixedAnalysis != nil {
		return this_.fixedAnalysis
	}
	switch this_.v.kind {
	case 0:
		return InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(this_, this_.declaration.node)
	case 1:
		// Eagerly evaluate isWellDefined (false) before CreateLazyNegative!
		_ = this_.priv_VarExpert_isWellDefined_1()
		return InlineVariables_InlineVarAnalysis_CreateLazyNegative(this_)
	case 2:
		return InlineVariables_InlineVarAnalysis_CreateLazyNegative(this_)
	case 3:
		return InlineVariables_InlineVarAnalysis_CreateWaitForAlias(this_.v.alias)
	default:
		return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS()
	}
}

func (this_ *InlineVariables_VarExpert) Reanalyze(av *Var, ana *InlineVariables_InlineVarAnalysis) *InlineVariables_InlineVarAnalysis {
	_ = av
	if ana.IsSafeToInlineAliases() {
		return InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(this_, this_.declaration.node)
	}
	return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS()
}

type InlineVariables_InliningBehavior struct {
	outer_ *InlineVariables
}

func NewInlineVariables_InliningBehavior(outer_ *InlineVariables) *InlineVariables_InliningBehavior {
	return &InlineVariables_InliningBehavior{outer_: outer_}
}

func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_createVarExpert_3(v *Var, rc *ReferenceCollection) *InlineVariables_VarExpert {
	if rc == nil {
		return StaticGet_NO_EXPERT()
	}
	return NewInlineVariables_VarExpert_5(this_, v, rc, true, false)
}

func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_doInlinesForScope_1() {
	handled := make(map[*Var]*InlineVariables_InlineVarAnalysis)
	waiting := make(map[*Var][]*InlineVariables_InliningBehavior_AliasInlineRetryHandler)
	for _, v := range this_.outer_.vars {
		exp := this_.priv_InliningBehavior_createVarExpert_3(v, this_.outer_.cols[v])
		ana := exp.Analyze()
		if ana.ShouldWaitForAliasedVar() {
			rh := NewInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_, v, exp)
			av := ana.GetAliasedVar()
			if prev := handled[av]; prev != nil {
				rh.Handle(av, prev)
			} else {
				waiting[av] = append(waiting[av], rh)
			}
		} else {
			if ana.ShouldInline() {
				ana.PerformInline()
			}
			handled[v] = ana
			for _, rh := range waiting[v] {
				rh.Handle(v, ana)
			}
		}
	}
	_ = j2gort.NullString
}

type InlineVariables_InliningBehavior_AliasInlineRetryHandler struct {
	outer_  *InlineVariables_InliningBehavior
	v       *Var
	expert  *InlineVariables_VarExpert
	idHash_ atomic.Int32
}

func NewInlineVariables_InliningBehavior_AliasInlineRetryHandler(outer_ *InlineVariables_InliningBehavior, v *Var, expert *InlineVariables_VarExpert) *InlineVariables_InliningBehavior_AliasInlineRetryHandler {
	this_ := &InlineVariables_InliningBehavior_AliasInlineRetryHandler{}
	this_.outer_ = outer_
	InitInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_, v, expert)
	return this_
}
func InitInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler, v *Var, expert *InlineVariables_VarExpert) {
	this_.v = v
	this_.expert = expert
}
func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) Handle(av *Var, ana *InlineVariables_InlineVarAnalysis) {
	na := this_.expert.Reanalyze(av, ana)
	if na.ShouldInline() {
		na.PerformInline()
	}
}
`

			inlineFuncsSrc := `package jscomp

import (
	"example.com/w42b_opt/gen/j2gort"
	pkg_rhino "example.com/w42b_opt/gen/rhino"
)

type FunctionInjector struct {
	callerFnStateCache j2gort.OrderedMap[pkg_rhino.Node, pkg_rhino.Node]
}

func NewFunctionInjector() *FunctionInjector {
	fi := &FunctionInjector{}
	fi.callerFnStateCache = j2gort.NewOrderedMap[pkg_rhino.Node, pkg_rhino.Node]()
	return fi
}

func (this_ *FunctionInjector) priv_FunctionInjector_getCallerFnState_2931(callerFn pkg_rhino.Node) pkg_rhino.Node {
	var cached pkg_rhino.Node = this_.callerFnStateCache.Get(callerFn)
	if cached != 0 {
		return cached
	}
	cached = callerFn + 1000
	this_.callerFnStateCache.Put(callerFn, cached)
	return cached
}

type InlineFunctions_Reference struct {
	callNode pkg_rhino.Node
	scope    *Scope
}

func NewInlineFunctions_Reference(callNode pkg_rhino.Node, scope *Scope) *InlineFunctions_Reference {
	r := &InlineFunctions_Reference{}
	InitInlineFunctions_Reference(r, callNode, scope)
	return r
}
func InitInlineFunctions_Reference(this_ *InlineFunctions_Reference, callNode pkg_rhino.Node, scope *Scope) {
	this_.callNode = callNode
	this_.scope = scope
}

type InlineFunctions_FunctionState struct {
	inline     bool
	remove     bool
	ref0       *InlineFunctions_Reference
	ref1       *InlineFunctions_Reference
	references *j2gort.OrderedMap[pkg_rhino.Node, *InlineFunctions_Reference]
}

func NewInlineFunctions_FunctionState() *InlineFunctions_FunctionState {
	this_ := &InlineFunctions_FunctionState{}
	InitInlineFunctions_FunctionState(this_)
	return this_
}
func InitInlineFunctions_FunctionState(this_ *InlineFunctions_FunctionState) {
	this_.inline = true
	this_.remove = true
	this_.references = nil
}

func (this_ *InlineFunctions_FunctionState) AddReference(ref *InlineFunctions_Reference) {
	if this_.references != nil {
		this_.references.Put(ref.callNode, ref)
		return
	}
	if this_.ref0 == nil || this_.ref0.callNode == ref.callNode {
		this_.ref0 = ref
		return
	}
	if this_.ref1 == nil || this_.ref1.callNode == ref.callNode {
		this_.ref1 = ref
		return
	}
	r0 := this_.ref0
	r1 := this_.ref1
	this_.ref0 = nil
	this_.ref1 = nil
	this_.references = j2gort.NewOrderedMapPtrCap[pkg_rhino.Node, *InlineFunctions_Reference](int32(8))
	this_.references.Put(r0.callNode, r0)
	this_.references.Put(r1.callNode, r1)
	this_.references.Put(ref.callNode, ref)
}

type InlineFunctions_FindCandidatesReferences struct {
	outer_ *InlineFunctions
}

func (this_ *InlineFunctions_FindCandidatesReferences) VisitCalls(fnName string, calls []pkg_rhino.Node) {
	fs := this_.outer_.priv_InlineFunctions_getOrCreateFunctionState_1(fnName)
	for _, c := range calls {
		_ = this_.outer_.injector.priv_FunctionInjector_getCallerFnState_2931(c)
		ref := NewInlineFunctions_Reference(c, nil)
		fs.AddReference(ref)
	}
}

type InlineFunctions struct {
	compiler AbstractCompiler
	injector *FunctionInjector
	fns      j2gort.OrderedMap[string, *InlineFunctions_FunctionState]
	anonFns  j2gort.OrderedMap[pkg_rhino.Node, string]
}

func (*InlineFunctions) Implements_jscomp_CompilerPass() {}

func NewInlineFunctions() *InlineFunctions {
	p := &InlineFunctions{injector: NewFunctionInjector()}
	p.initFields_InlineFunctions()
	return p
}
func (this_ *InlineFunctions) initFields_InlineFunctions() {
	this_.fns = j2gort.NewOrderedMapCap[string, *InlineFunctions_FunctionState](int32(1024))
	this_.anonFns = j2gort.NewOrderedMap[pkg_rhino.Node, string]()
}

func (this_ *InlineFunctions) priv_InlineFunctions_getOrCreateFunctionState_1(name string) *InlineFunctions_FunctionState {
	fs := this_.fns.Get(name)
	if fs == nil {
		fs = NewInlineFunctions_FunctionState()
		this_.fns.Put(name, fs)
	}
	return fs
}

func (this_ *InlineFunctions) Process_2(_ pkg_rhino.Node, _ pkg_rhino.Node) {
	defer func() {
		for cur_ := this_.fns.IterValues(); cur_.HasNext(); {
			fs_ := cur_.Next()
			if fs_ != nil && fs_.references != nil {
				j2gort.ReleaseOrderedMap(fs_.references)
				fs_.references = nil
			}
		}
		this_.fns.ReleaseToPool()
		this_.anonFns.ReleaseToPool()
	}()
	calls := make([]pkg_rhino.Node, 270)
	for i := range calls {
		calls[i] = pkg_rhino.Node(i + 1)
	}
	finder := &InlineFunctions_FindCandidatesReferences{outer_: this_}
	finder.VisitCalls("hotHelper", calls)
}
`

			contStructSrc := `type RemoveUnusedCode_Continuation struct {
	outer_ *RemoveUnusedCode
	node   pkg_rhino.Node
}
func MakeContinuation(outer_ *RemoveUnusedCode, n pkg_rhino.Node) RemoveUnusedCode_Continuation {
	return RemoveUnusedCode_Continuation{outer_: outer_, node: n}
}`
			if staticCont {
				contStructSrc = `type RemoveUnusedCode_Continuation struct {
	node pkg_rhino.Node
}
func MakeContinuation(_ *RemoveUnusedCode, n pkg_rhino.Node) RemoveUnusedCode_Continuation {
	return RemoveUnusedCode_Continuation{node: n}
}`
			}

			rucSrc := fmt.Sprintf(`package jscomp

import (
	"sync/atomic"

	"example.com/w42b_opt/gen/j2gort"
	pkg_rhino "example.com/w42b_opt/gen/rhino"
)

var BindSelfCalls atomic.Int64

%s

type RemoveUnusedCode_Removable interface {
	Implements_jscomp_RemoveUnusedCode_Removable()
	IsJavaNil_() bool
	GetRemoveUnusedCode_RemovableBase() *RemoveUnusedCode_RemovableBase
	Remove()
	RemoveInternal()
}

func ToIface_RemoveUnusedCode_Removable[T interface {
	RemoveUnusedCode_Removable
	comparable
}](v T) RemoveUnusedCode_Removable {
	var z T
	if v == z {
		return nil
	}
	return v
}

type RemoveUnusedCode_RemovableBase struct {
	self_                           RemoveUnusedCode_Removable
	outer_                          *RemoveUnusedCode
	continuation0                   RemoveUnusedCode_Continuation
	continuation1                   RemoveUnusedCode_Continuation
	restContinuations               []RemoveUnusedCode_Continuation
	propertyName                    string
	targetNode                      pkg_rhino.Node
	idHash_                         atomic.Int32
	isPrototypeDotPropertyReference bool
	isThisDotPropertyReference      bool
	continuationsAreApplied         bool
	isRemoved                       bool
}

func (this_ *RemoveUnusedCode_RemovableBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(RemoveUnusedCode_Removable)
}
func (this_ *RemoveUnusedCode_RemovableBase) IsJavaNil_() bool { return this_ == nil }
func (b *RemoveUnusedCode_RemovableBase) GetRemoveUnusedCode_RemovableBase() *RemoveUnusedCode_RemovableBase {
	return b
}
func (*RemoveUnusedCode_RemovableBase) Implements_jscomp_RemoveUnusedCode_Removable() {}
func (this_ *RemoveUnusedCode_RemovableBase) RemoveInternal()                          { panic("abstract") }
func (this_ *RemoveUnusedCode_RemovableBase) Remove() {
	if !this_.isRemoved {
		this_.isRemoved = true
		this_.self_.RemoveInternal()
	}
}
func (this_ *RemoveUnusedCode_RemovableBase) initFields_RemoveUnusedCode_Removable() {
	_ = this_
	this_.propertyName = j2gort.NullString
}
func InitRemoveUnusedCode_Removable(this_ *RemoveUnusedCode_RemovableBase, targetNode pkg_rhino.Node, builder *RemoveUnusedCode_RemovableBuilder) {
	_ = targetNode
	_ = builder
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
	continuation0                   RemoveUnusedCode_Continuation
	continuation1                   RemoveUnusedCode_Continuation
	restContinuations               []RemoveUnusedCode_Continuation
	propertyName                    string
	idHash_                         atomic.Int32
	isPrototypeDotPropertyReference bool
	isThisDotPropertyReference      bool
}

func NewRemoveUnusedCode_RemovableBuilder(outer_ *RemoveUnusedCode) *RemoveUnusedCode_RemovableBuilder {
	this_ := &RemoveUnusedCode_RemovableBuilder{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_RemovableBuilder(this_)
	return this_
}
func InitRemoveUnusedCode_RemovableBuilder(this_ *RemoveUnusedCode_RemovableBuilder) {
	this_.propertyName = j2gort.NullString
}

func (this_ *RemoveUnusedCode_RemovableBuilder) AddContinuation(c RemoveUnusedCode_Continuation) *RemoveUnusedCode_RemovableBuilder {
	if this_.continuation0 == (RemoveUnusedCode_Continuation{}) {
		this_.continuation0 = c
	} else if this_.continuation1 == (RemoveUnusedCode_Continuation{}) {
		this_.continuation1 = c
	} else {
		this_.restContinuations = append(this_.restContinuations, c)
	}
	return this_
}

func (this_ *RemoveUnusedCode_RemovableBuilder) priv_RemovableBuilder_finish_4743(removable RemoveUnusedCode_Removable) RemoveUnusedCode_Removable {
	_ = removable
	this_.continuation0 = (RemoveUnusedCode_Continuation{})
	this_.continuation1 = (RemoveUnusedCode_Continuation{})
	this_.restContinuations = nil
	this_.propertyName = j2gort.NullString
	this_.isPrototypeDotPropertyReference = false
	this_.isThisDotPropertyReference = false
	this_.outer_.freeBuilder = this_
	return removable
}

func (this_ *RemoveUnusedCode_RemovableBuilder) BuildIndirectAssign(root pkg_rhino.Node, targetNode pkg_rhino.Node) *RemoveUnusedCode_IndirectAssign {
	return j2gort.CastPtr[RemoveUnusedCode_IndirectAssign](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(NewRemoveUnusedCode_IndirectAssign(this_.outer_, this_, root, targetNode))))
}
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildNameDeclarationStatement(decl pkg_rhino.Node) *RemoveUnusedCode_NameDeclarationStatement {
	return j2gort.CastPtr[RemoveUnusedCode_NameDeclarationStatement](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(NewRemoveUnusedCode_NameDeclarationStatement(this_.outer_, this_, decl))))
}
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildAssign(assignNode pkg_rhino.Node, propName string) *RemoveUnusedCode_Assign {
	this_.propertyName = propName
	return j2gort.CastPtr[RemoveUnusedCode_Assign](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(NewRemoveUnusedCode_Assign(this_.outer_, this_, assignNode))))
}
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildInstanceofName(n pkg_rhino.Node) RemoveUnusedCode_Removable {
	return ToIface_RemoveUnusedCode_Removable(j2gort.CastPtr[RemoveUnusedCode_InstanceofName](this_.priv_RemovableBuilder_finish_4743(ToIface_RemoveUnusedCode_Removable(NewRemoveUnusedCode_InstanceofName(this_.outer_, this_, n)))))
}

type RemoveUnusedCode_IndirectAssign struct {
	RemoveUnusedCode_RemovableBase
	outer_ *RemoveUnusedCode
	root   pkg_rhino.Node
}
func (this_ *RemoveUnusedCode_IndirectAssign) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_IndirectAssign) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_IndirectAssign) Implements_jscomp_RemoveUnusedCode_Removable() {}
func (this_ *RemoveUnusedCode_IndirectAssign) RemoveInternal() { this_.outer_.RemovedCount++ }
func NewRemoveUnusedCode_IndirectAssign(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, root pkg_rhino.Node, targetNode pkg_rhino.Node) *RemoveUnusedCode_IndirectAssign {
	this_ := &RemoveUnusedCode_IndirectAssign{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_IndirectAssign(this_, builder, root, targetNode)
	return this_
}
func InitRemoveUnusedCode_IndirectAssign(this_ *RemoveUnusedCode_IndirectAssign, builder *RemoveUnusedCode_RemovableBuilder, root pkg_rhino.Node, targetNode pkg_rhino.Node) {
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, targetNode, builder)
	this_.root = root
}

type RemoveUnusedCode_NameDeclarationStatement struct {
	RemoveUnusedCode_RemovableBase
	outer_ *RemoveUnusedCode
	decl   pkg_rhino.Node
}
func (this_ *RemoveUnusedCode_NameDeclarationStatement) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_NameDeclarationStatement) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_NameDeclarationStatement) Implements_jscomp_RemoveUnusedCode_Removable() {}
func (this_ *RemoveUnusedCode_NameDeclarationStatement) RemoveInternal() { this_.outer_.RemovedCount += 10 }
func NewRemoveUnusedCode_NameDeclarationStatement(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, decl pkg_rhino.Node) *RemoveUnusedCode_NameDeclarationStatement {
	this_ := &RemoveUnusedCode_NameDeclarationStatement{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_NameDeclarationStatement(this_, builder, decl)
	return this_
}
func InitRemoveUnusedCode_NameDeclarationStatement(this_ *RemoveUnusedCode_NameDeclarationStatement, builder *RemoveUnusedCode_RemovableBuilder, decl pkg_rhino.Node) {
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, decl, builder)
	this_.decl = decl
}

type RemoveUnusedCode_Assign struct {
	RemoveUnusedCode_RemovableBase
	outer_     *RemoveUnusedCode
	assignNode pkg_rhino.Node
}
func (this_ *RemoveUnusedCode_Assign) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_Assign) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_Assign) Implements_jscomp_RemoveUnusedCode_Removable() {}
func (this_ *RemoveUnusedCode_Assign) RemoveInternal() { this_.outer_.RemovedCount += 100 }
func NewRemoveUnusedCode_Assign(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, assignNode pkg_rhino.Node) *RemoveUnusedCode_Assign {
	this_ := &RemoveUnusedCode_Assign{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_Assign(this_, builder, assignNode)
	return this_
}
func InitRemoveUnusedCode_Assign(this_ *RemoveUnusedCode_Assign, builder *RemoveUnusedCode_RemovableBuilder, assignNode pkg_rhino.Node) {
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, assignNode, builder)
	this_.assignNode = assignNode
}

type RemoveUnusedCode_InstanceofName struct {
	RemoveUnusedCode_RemovableBase
	outer_ *RemoveUnusedCode
	node   pkg_rhino.Node
}
func (this_ *RemoveUnusedCode_InstanceofName) BindSelf_(s any) {
	_ = s
	this_.RemoveUnusedCode_RemovableBase.BindSelf_(s)
}
func (this_ *RemoveUnusedCode_InstanceofName) IsJavaNil_() bool { return this_ == nil }
func (*RemoveUnusedCode_InstanceofName) Implements_jscomp_RemoveUnusedCode_Removable() {}
func (this_ *RemoveUnusedCode_InstanceofName) RemoveInternal() { this_.outer_.RemovedCount += 1000 }
func NewRemoveUnusedCode_InstanceofName(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, n pkg_rhino.Node) *RemoveUnusedCode_InstanceofName {
	this_ := &RemoveUnusedCode_InstanceofName{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitRemoveUnusedCode_InstanceofName(this_, builder, n)
	return this_
}
func InitRemoveUnusedCode_InstanceofName(this_ *RemoveUnusedCode_InstanceofName, builder *RemoveUnusedCode_RemovableBuilder, n pkg_rhino.Node) {
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, n, builder)
	this_.node = n
}

type RemoveUnusedCode struct {
	freeBuilder  *RemoveUnusedCode_RemovableBuilder
	RemovedCount int
	ContCount    int
}

func (*RemoveUnusedCode) Implements_jscomp_CompilerPass() {}

func (this_ *RemoveUnusedCode) acquireBuilder() *RemoveUnusedCode_RemovableBuilder {
	b := this_.freeBuilder
	if b != nil {
		this_.freeBuilder = nil
		return b
	}
	return NewRemoveUnusedCode_RemovableBuilder(this_)
}

func (this_ *RemoveUnusedCode) Process_2(_ pkg_rhino.Node, _ pkg_rhino.Node) {
	r1 := this_.acquireBuilder().BuildIndirectAssign(1, 2)
	r2 := this_.acquireBuilder().AddContinuation(MakeContinuation(this_, 10)).AddContinuation(MakeContinuation(this_, 20)).AddContinuation(MakeContinuation(this_, 30)).BuildNameDeclarationStatement(3)
	r3 := this_.acquireBuilder().BuildAssign(4, "propX")
	r4 := this_.acquireBuilder().BuildInstanceofName(5)
	if r2.continuation0 != (RemoveUnusedCode_Continuation{}) && r2.continuation1 != (RemoveUnusedCode_Continuation{}) && len(r2.restContinuations) == 1 {
		this_.ContCount = 3
	}
	// Verify builder state was cleanly reset after 3-continuation + propertyName builds!
	r5 := this_.acquireBuilder().BuildIndirectAssign(6, 7)
	if r5.continuation0 != (RemoveUnusedCode_Continuation{}) || r5.propertyName != j2gort.NullString || r3.propertyName != "propX" {
		panic("builder state was not cleanly reset!")
	}
	r1.Remove()
	r2.Remove()
	r3.Remove()
	r4.Remove()
	r5.Remove()
}
`, contStructSrc)

			rtTestSrc := `package jscomp

import (
	"testing"

	"example.com/w42b_opt/gen/j2gort"
)

func TestWave42BOptHotspotsRuntime(t *testing.T) {
	// 1. InlineVariables: 50 LazyNegative (known false) + 50 LazyNegative (unknown, unaliased)
	//    + 1 LazyNegative (unknown, aliased & true) + 1 WaitForAlias!
	vTargetTrue := &Var{name: "targetTrue", kind: 2}
	vTargetFalse := &Var{name: "targetFalse", kind: 2}
	vWaitTrue := &Var{name: "waitTrue", kind: 3, alias: vTargetTrue}
	vWaitFalse := &Var{name: "waitFalse", kind: 3, alias: vTargetFalse}

	vars := []*Var{vWaitTrue, vWaitFalse}
	cols := map[*Var]*ReferenceCollection{
		vWaitTrue:    {refs: []*Reference{{node: 1}}, wellDef: true, assignedOnce: true},
		vWaitFalse:   {refs: []*Reference{{node: 2}}, wellDef: true, assignedOnce: true},
		vTargetTrue:  {refs: []*Reference{{node: 3}}, wellDef: true, assignedOnce: true},
		vTargetFalse: {refs: []*Reference{{node: 4}}, wellDef: true, assignedOnce: false},
	}
	var unaliasedUnknown []*ReferenceCollection
	for i := 0; i < 40; i++ {
		vk := &Var{name: "knownFalse", kind: 1}
		vu := &Var{name: "unaliasedLazy", kind: 2}
		rcU := &ReferenceCollection{refs: []*Reference{{node: 5}}, wellDef: true, assignedOnce: true}
		unaliasedUnknown = append(unaliasedUnknown, rcU)
		vars = append(vars, vk, vu)
		cols[vk] = &ReferenceCollection{refs: []*Reference{{node: 6}}, wellDef: false, assignedOnce: false}
		cols[vu] = rcU
	}
	vars = append(vars, vTargetTrue, vTargetFalse)

	iv := &InlineVariables{vars: vars, cols: cols}
	iv.Process_2(1, 2)

	if len(iv.inlinedNames) != 1 || iv.inlinedNames[0] != "waitTrue" {
		t.Fatalf("expected ONLY waitTrue to inline, got %v", iv.inlinedNames)
	}
	// ONLY the 2 AliasInlineRetryHandlers may promote scratchVarExpert_ — ZERO promotions from 82 LazyNegative vars!
	if got := iv.slab_InlineVariables_VarExpert.AllocCalls.Load(); got != 2 {
		t.Fatalf("expected ONLY 2 slab_InlineVariables_VarExpert allocs (for 2 RetryHandlers, 0 for 82 LazyNegatives!), got %d", got)
	}
	// The 40 knownFalse LazyNegatives folded directly to NO_INLINE_SELF_OR_ALIASES_ANALYSIS (0 slab allocs!),
	// leaving ONLY 2 WaitForAlias + 1 inlined waitTrue + 42 unknown LazyNegatives = 45 InlineVarAnalysis allocs!
	if got := iv.slab_InlineVariables_InlineVarAnalysis.AllocCalls.Load(); got != 45 {
		t.Fatalf("expected 45 InlineVarAnalysis slab allocs (40 known-false folded to singleton!), got %d", got)
	}
	// Unaliased LazyNegative vars MUST NEVER evaluate IsWellDefined()!
	for _, rcU := range unaliasedUnknown {
		if rcU.WellDefCalls != 0 {
			t.Fatalf("expected 0 IsWellDefined() calls on unaliased LazyNegative var, got %d", rcU.WellDefCalls)
		}
	}

	// 2. InlineFunctions + FunctionInjector:
	j2gort.AcquireMapCapCalls.Store(0)
	j2gort.ReleaseMapToPoolCalls.Store(0)
	j2gort.LargePtrMapPromotes.Store(0)

	inf := NewInlineFunctions()
	inf.Process_2(1, 2)

	if got := inf.slab_InlineFunctions_FunctionState.AllocCalls.Load(); got != 1 {
		t.Fatalf("expected FunctionState to allocate from slab_InlineFunctions_FunctionState (1), got %d", got)
	}
	if got := inf.slab_InlineFunctions_FunctionState.Releases.Load(); got != 1 {
		t.Fatalf("expected slab_InlineFunctions_FunctionState.Release() to run on Process_2 exit, got %d", got)
	}
	if got := j2gort.LargePtrMapPromotes.Load(); got != 1 {
		t.Fatalf("expected 270-ref FunctionState.references to promote once at Size()==256 to cap 512, got %d", got)
	}
	if got := j2gort.ReleaseMapToPoolCalls.Load(); got < 3 {
		t.Fatalf("expected fns + anonFns + injector.callerFnStateCache ReleaseToPool (>=3), got %d", got)
	}

	// 3. RemoveUnusedCode: virtual dispatch via static self_ = this_, builder reset, and continuations!
	BindSelfCalls.Store(0)
	ruc := &RemoveUnusedCode{}
	ruc.Process_2(1, 2)
	if ruc.RemovedCount != 1112 || ruc.ContCount != 3 {
		t.Fatalf("unexpected RemoveUnusedCode counts: RemovedCount=%d ContCount=%d", ruc.RemovedCount, ruc.ContCount)
	}
	if got := BindSelfCalls.Load(); got != 0 {
		t.Fatalf("expected 0 BindSelf_ interface boxing calls on slab-allocated Removables, got %d", got)
	}
}
`

			wave42bWriteFiles(t, tmpDir, map[string]string{
				"go.mod":                            "module example.com/w42b_opt\n\ngo 1.24\n",
				"gen/j2gort/j2gort.go":              j2gortSrc,
				"gen/rhino/rhino.go":                rhinoSrc,
				"gen/jscomp/inline_variables.go":    inlineVarsSrc,
				"gen/jscomp/inline_functions.go":    inlineFuncsSrc,
				"gen/jscomp/remove_unused_code.go":  rucSrc,
			})

			jscompDir := filepath.Join(tmpDir, "gen/jscomp")
			allPasses := passes.AllPasses()

			s1, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
			if err != nil {
				t.Fatalf("Run 1 failed: %v", err)
			}
			if s1.ByPass["slab_alloc_lowering"] == 0 {
				t.Fatalf("expected slab_alloc_lowering sites on Run 1, got %+v", s1)
			}

			s2, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
			if err != nil {
				t.Fatalf("Run 2 failed: %v", err)
			}
			if s2.TotalSites != 0 || s2.FilesModified != 0 {
				t.Fatalf("Run 2 NOT IDEMPOTENT: %+v", s2)
			}

			rucBytes, err := os.ReadFile(filepath.Join(jscompDir, "remove_unused_code.go"))
			if err != nil {
				t.Fatal(err)
			}
			rucLowered := string(rucBytes)
			for _, want := range []string{
				"res_ := rucOuter_.slab_RemoveUnusedCode_IndirectAssign.Alloc()",
				"res_ := rucOuter_.slab_RemoveUnusedCode_Assign.Alloc()",
				"res_ := rucOuter_.slab_RemoveUnusedCode_NameDeclarationStatement.Alloc()",
				"res_.self_ = res_",
				"this_.priv_RemovableBuilder_finish_4743(nil)",
				"if builder.continuation0 !=",
			} {
				if !strings.Contains(rucLowered, want) {
					t.Fatalf("expected lowered remove_unused_code.go to contain %q:\n%s", want, rucLowered)
				}
			}
			// Verify Build* methods no longer call allocPooled_RemoveUnusedCode_* wrappers!
			for _, badCall := range []string{
				"allocPooled_RemoveUnusedCode_IndirectAssign(",
				"allocPooled_RemoveUnusedCode_Assign(",
				"allocPooled_RemoveUnusedCode_NameDeclarationStatement(",
			} {
				if strings.Contains(rucLowered, badCall+"this_") {
					t.Fatalf("Build* MUST inline slab Alloc directly instead of calling %s:\n%s", badCall, rucLowered)
				}
			}

			wave42bWriteFiles(t, tmpDir, map[string]string{
				"gen/jscomp/opt_rt_test.go": rtTestSrc,
			})
			cmd := exec.Command("go", "test", "-race", "-v", "./gen/jscomp")
			cmd.Dir = tmpDir
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("go test -race ./gen/jscomp failed: %v\n%s", err, string(out))
			}
		})
	}
}

// TestWave42B_RealGenRhinoAndJscompIdempotencyAndTargets copies real transpiled/gen/rhino and
// transpiled/gen/jscomp into isolated modules, runs AllPasses() Run 1 + Run 2, verifies all Wave 42B
// rewrites landed cleanly, asserts Run 2 produces TotalSites == 0, and compiles both packages.
func TestWave42B_RealGenRhinoAndJscompIdempotencyAndTargets(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping full real gen/rhino + gen/jscomp build & idempotency check in -short mode")
	}
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !wave42bRealGenReady(repoRoot) {
		t.Skip("gen/rhino or gen/j2gort is mid-edit; skipping real gen copy check")
	}

	origModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	sumBytes, _ := os.ReadFile(filepath.Join(repoRoot, "go.sum"))

	copyAndRunRealPkg := func(pkgName string) string {
		t.Helper()
		srcDir := filepath.Join(repoRoot, "transpiled", "gen", pkgName)
		entries, rErr := os.ReadDir(srcDir)
		if rErr != nil {
			t.Fatal(rErr)
		}
		tmpDir := t.TempDir()
		goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/realw42b_"+pkgName, 1) +
			fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
		if wErr := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); wErr != nil {
			t.Fatal(wErr)
		}
		if len(sumBytes) > 0 {
			_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
		}
		dstDir := filepath.Join(tmpDir, pkgName)
		if mErr := os.MkdirAll(dstDir, 0o755); mErr != nil {
			t.Fatal(mErr)
		}
		for _, ent := range entries {
			if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") || strings.HasSuffix(ent.Name(), "_test.go") {
				continue
			}
			data, errR := os.ReadFile(filepath.Join(srcDir, ent.Name()))
			if errR != nil {
				t.Fatal(errR)
			}
			if errW := os.WriteFile(filepath.Join(dstDir, ent.Name()), data, 0o644); errW != nil {
				t.Fatal(errW)
			}
		}

		allPasses := passes.AllPasses()
		s1, err1 := passes.RunOnPackageDir(dstDir, allPasses, passes.Options{})
		if err1 != nil {
			t.Fatalf("AllPasses() Run 1 on real gen/%s failed: %v", pkgName, err1)
		}
		t.Logf("Wave 42B real gen/%s Run 1: total=%d files=%d byPass=%+v", pkgName, s1.TotalSites, s1.FilesModified, s1.ByPass)

		s2, err2 := passes.RunOnPackageDir(dstDir, allPasses, passes.Options{})
		if err2 != nil {
			t.Fatalf("AllPasses() Run 2 on real gen/%s failed: %v", pkgName, err2)
		}
		if s2.TotalSites != 0 || s2.FilesModified != 0 {
			t.Fatalf("Wave 42B real gen/%s Run 2 NOT IDEMPOTENT: total=%d files=%d byPass=%+v", pkgName, s2.TotalSites, s2.FilesModified, s2.ByPass)
		}

		buildCmd := exec.Command("go", "build", "./"+pkgName)
		buildCmd.Dir = tmpDir
		if out, bErr := buildCmd.CombinedOutput(); bErr != nil {
			t.Fatalf("go build ./%s on lowered real gen/%s failed: %v\n%s", pkgName, pkgName, bErr, string(out))
		}
		return dstDir
	}

	rhinoCopy := copyAndRunRealPkg("rhino")
	readFileStr := func(dir, name string) string {
		b, rErr := os.ReadFile(filepath.Join(dir, name))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}

	// Task 1 checks on real rhino
	subChecker := readFileStr(rhinoCopy, "subtype_checker.go")
	for _, want := range []string{
		"var subMask int32 = base_subtype_.typeKindMask",
		"var supMask int32 = base_supertype_.typeKindMask",
	} {
		if !strings.Contains(subChecker, want) {
			t.Fatalf("expected lowered real rhino/subtype_checker.go to contain %q", want)
		}
	}
	eqChecker := readFileStr(rhinoCopy, "equality_checker.go")
	if !strings.Contains(eqChecker, "if base_left_.IsNoResolvedType() && base_right_.IsNoResolvedType()") {
		t.Fatalf("expected lowered real rhino/equality_checker.go to hoist base_right_ before IsNoResolvedType()")
	}

	jscompCopy := copyAndRunRealPkg("jscomp")
	// Task 2 checks on real jscomp
	tas := readFileStr(jscompCopy, "typed_ast_serializer.go")
	if strings.Contains(tas, "SetOriginalNamePointer(func()") || !strings.Contains(tas, "SetOriginalNamePointer(iifeArg_") {
		t.Fatalf("expected SetOriginalNamePointer(func()...) in typed_ast_serializer.go to unwrap to iifeArg_*")
	}
	// Task 3 checks on real jscomp
	iv := readFileStr(jscompCopy, "inline_variables.go")
	if !strings.Contains(iv, "lazyRefInfo_") {
		t.Fatalf("expected lazyRefInfo_ in lowered real inline_variables.go")
	}
	inf := readFileStr(jscompCopy, "inline_functions.go")
	for _, want := range []string{
		"slab_InlineFunctions_FunctionState",
		"callerFnStateCache.ReleaseToPool()",
		"PutAllMapPtr(",
	} {
		if !strings.Contains(inf, want) {
			t.Fatalf("expected lowered real inline_functions.go to contain %q", want)
		}
	}
	ruc := readFileStr(jscompCopy, "remove_unused_code.go")
	for _, want := range []string{
		"this_.self_ = this_",
		"res_ := rucOuter_.slab_RemoveUnusedCode_IndirectAssign.Alloc()",
		"res_ := rucOuter_.slab_RemoveUnusedCode_Assign.Alloc()",
		"res_ := rucOuter_.slab_RemoveUnusedCode_NameDeclarationStatement.Alloc()",
		"res_.self_ = res_",
		".priv_RemovableBuilder_finish_",
	} {
		if !strings.Contains(ruc, want) {
			t.Fatalf("expected lowered real remove_unused_code.go to contain %q", want)
		}
	}
}

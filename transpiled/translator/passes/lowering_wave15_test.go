package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave15CrossPackageSliceGetterLowering(t *testing.T) {
	tmpDir, _ := setupRealJ2gortOverlayWorkspace(t)

	// Place both `jscomp` and `rhino` under `gen/` so `jscomp` < `rhino` alphabetically,
	// directly verifying that `passes.Run(..., "gen/...")` topologically sorts `rhino` BEFORE `jscomp`!
	origRhinoDir := filepath.Join(tmpDir, "orig", "rhino")
	origJscompDir := filepath.Join(tmpDir, "orig", "jscomp")
	optRhinoDir := filepath.Join(tmpDir, "opt", "rhino")
	optJscompDir := filepath.Join(tmpDir, "opt", "jscomp")
	for _, d := range []string{origRhinoDir, origJscompDir, optRhinoDir, optJscompDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	rhinoSrc := `package rhino

import "example.com/w13maptest/j2gort"

type JSType interface {
	Implements_rhino_JSType()
	IsJavaNil_() bool
	GetDisplayName() string
	IsResolved() bool
}

type ScalarType struct {
	name     string
	resolved bool
}

func NewScalarType(name string, resolved bool) *ScalarType {
	return &ScalarType{name: name, resolved: resolved}
}
func (*ScalarType) Implements_rhino_JSType()         {}
func (s *ScalarType) IsJavaNil_() bool               { return s == nil }
func (s *ScalarType) GetDisplayName() string         { return s.name }
func (s *ScalarType) IsResolved() bool               { return s.resolved }
func (s *ScalarType) MarkResolved(newName string)    { s.resolved = true; s.name = newName }

type FunctionType_Parameter struct {
	 paramType JSType
	 optional  bool
}

func NewFunctionType_Parameter(t JSType, opt bool) FunctionType_Parameter {
	return FunctionType_Parameter{paramType: t, optional: opt}
}
func (p FunctionType_Parameter) GetJSType() JSType { return p.paramType }
func (p FunctionType_Parameter) IsOptional() bool  { return p.optional }

// Polymorphic interface FunctionType + FunctionTypeBase + concrete subclass NoType.
type FunctionType interface {
	Implements_rhino_FunctionType()
	IsJavaNil_() bool
	GetFunctionTypeBase() *FunctionTypeBase
	JSType
	GetParameters() j2gort.Ext_com_google_common_collect_ImmutableList
}

type FunctionTypeBase struct {
	self_         FunctionType
	name          string
	parameterList j2gort.Ext_com_google_common_collect_ImmutableList
}

func (*FunctionTypeBase) Implements_rhino_JSType()                 {}
func (*FunctionTypeBase) Implements_rhino_FunctionType()           {}
func (this_ *FunctionTypeBase) IsJavaNil_() bool                   { return this_ == nil }
func (this_ *FunctionTypeBase) GetFunctionTypeBase() *FunctionTypeBase { return this_ }
func (this_ *FunctionTypeBase) GetDisplayName() string             { return this_.name }
func (this_ *FunctionTypeBase) IsResolved() bool                   { return true }

func InitFunctionType(this_ *FunctionTypeBase, name string, params []FunctionType_Parameter) {
	this_.name = name
	var b *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	for _, p := range params {
		b.AddVoid_(p)
	}
	this_.parameterList = b.BuildImmList_()
}

func NewFunctionType(name string, params []FunctionType_Parameter) *FunctionTypeBase {
	f := &FunctionTypeBase{}
	f.self_ = f
	InitFunctionType(f, name, params)
	return f
}

func (this_ *FunctionTypeBase) GetParameters() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.parameterList
}

type NoType struct {
	FunctionTypeBase
}

func NewNoType() *NoType {
	n := &NoType{}
	n.self_ = n
	InitFunctionType(&n.FunctionTypeBase, "NoType", nil)
	return n
}

// UnionType with multi-statement lazy-resolution prefix in GetAlternates()!
type UnionType struct {
	alternates                    j2gort.Ext_com_google_common_collect_ImmutableList
	alternatesResolvedBeforeBuild bool
	resolvedAlternatesAtLastBuild int32
	rebuildCount                  int32
}

func (*UnionType) Implements_rhino_JSType()   {}
func (this_ *UnionType) IsJavaNil_() bool     { return this_ == nil }
func (this_ *UnionType) GetDisplayName() string { return "union" }
func (this_ *UnionType) IsResolved() bool     { return this_.alternatesResolvedBeforeBuild }
func (this_ *UnionType) RebuildCount() int32  { return this_.rebuildCount }

func InitUnionType(this_ *UnionType, alts []JSType, preResolved bool) {
	var b *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	for _, a := range alts {
		b.AddVoid_(a)
	}
	this_.alternates = b.BuildImmList_()
	this_.alternatesResolvedBeforeBuild = preResolved
	this_.resolvedAlternatesAtLastBuild = -1
}

func NewUnionType(alts []JSType, preResolved bool) *UnionType {
	u := &UnionType{}
	InitUnionType(u, alts, preResolved)
	return u
}

func (this_ *UnionType) countResolved() int32 {
	var c int32 = 0
	for _, raw := range j2gort.CheckNotNull(this_.alternates).IterSlice_() {
		t := j2gort.AsIface[JSType](raw)
		if t.IsResolved() {
			c++
		}
	}
	return c
}

func (this_ *UnionType) GetAlternates() j2gort.Ext_com_google_common_collect_ImmutableList {
	if !this_.IsResolved() && !this_.alternatesResolvedBeforeBuild {
		var r int32 = this_.countResolved()
		if r != this_.resolvedAlternatesAtLastBuild {
			this_.resolvedAlternatesAtLastBuild = r
			this_.rebuildCount++
			if int(r) == int(this_.alternates.Size()) {
				this_.alternatesResolvedBeforeBuild = true
			}
		}
	}
	return this_.alternates
}

type TemplatizedType struct {
	templateTypes j2gort.Ext_com_google_common_collect_ImmutableList
}

func (*TemplatizedType) Implements_rhino_JSType()     {}
func (this_ *TemplatizedType) IsJavaNil_() bool       { return this_ == nil }
func (this_ *TemplatizedType) GetDisplayName() string { return "templatized" }
func (this_ *TemplatizedType) IsResolved() bool       { return true }

func InitTemplatizedType(this_ *TemplatizedType, tts []JSType) {
	var b *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	for _, t := range tts {
		b.AddVoid_(t)
	}
	this_.templateTypes = b.BuildImmList_()
}

func NewTemplatizedType(tts []JSType) *TemplatizedType {
	tt := &TemplatizedType{}
	InitTemplatizedType(tt, tts)
	return tt
}

func (this_ *TemplatizedType) GetTemplateTypes() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.templateTypes
}
`

	jscompSrcTemplate := `package jscomp

import (
	"example.com/w13maptest/j2gort"
	pkg_rhino "example.com/w13maptest/KIND/rhino"
)

func CheckUnionAndFn(u *pkg_rhino.UnionType, fn pkg_rhino.FunctionType, tt *pkg_rhino.TemplatizedType) string {
	var score int32 = 0
	var firstName string = ""

	// 1. Cross-package IterSlice_() range loop over UnionType.GetAlternates()
	for _, raw := range j2gort.CheckNotNull(u.GetAlternates()).IterSlice_() {
		alt := j2gort.AsIface[pkg_rhino.JSType](raw)
		score += int32(len(alt.GetDisplayName()))
	}

	// 2. Cross-package NewIterCursor loop over UnionType.GetAlternates()
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(u.GetAlternates()))
		for cur.HasNext() {
			raw := cur.Next()
			alt := j2gort.AsIface[pkg_rhino.JSType](raw)
			if firstName == "" {
				firstName = alt.GetDisplayName()
			}
		}
	}

	// 3. Cross-package polymorphic interface FunctionType.GetParameters() via Iterator(), Size(), IsEmpty(), Get()
	if !fn.GetParameters().IsEmpty() {
		score += fn.GetParameters().Size() * 10
		p0 := j2gort.Cast[pkg_rhino.FunctionType_Parameter](fn.GetParameters().Get(int32(0)))
		if p0.IsOptional() {
			score += 3
		}
	}
	{
		var it j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](fn.GetParameters().Iterator())
		_ = it
		for it.HasNext() {
			param := j2gort.Cast[pkg_rhino.FunctionType_Parameter](it.Next())
			score += int32(len(param.GetJSType().GetDisplayName()))
		}
	}

	// 4. Cross-package ImmutableList local + CopyOf + Builder.AddAll + TemplatizedType.GetTemplateTypes().Get(0)
	var localAlts j2gort.Ext_com_google_common_collect_ImmutableList = u.GetAlternates()
	var copied j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](u.GetAlternates()))
	var b *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	b.AddAllVoid_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](u.GetAlternates()))
	var combined j2gort.Ext_com_google_common_collect_ImmutableList = b.BuildImmList_()
	score += localAlts.Size() + copied.Size() + combined.Size()

	if !tt.GetTemplateTypes().IsEmpty() {
		t0 := j2gort.AsIface[pkg_rhino.JSType](tt.GetTemplateTypes().Get(int32(0)))
		firstName = firstName + ":" + t0.GetDisplayName()
	}

	return firstName + "#" + j2gort.StrConcat(score)
}
`

	if err := os.WriteFile(filepath.Join(origRhinoDir, "types.go"), []byte(rhinoSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optRhinoDir, "types.go"), []byte(rhinoSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(origJscompDir, "check.go"), []byte(strings.ReplaceAll(jscompSrcTemplate, "KIND", "orig")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optJscompDir, "check.go"), []byte(strings.ReplaceAll(jscompSrcTemplate, "KIND", "opt")), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	optRootPattern := filepath.Join(tmpDir, "opt") + "/..."
	stats1, err := passes.Run([]string{optRootPattern}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run(opt/...) failed: %v", err)
	}
	if stats1.TotalSites == 0 {
		t.Fatalf("expected transformations across rhino + jscomp, got 0")
	}

	// Strictly 0-site Run-2 idempotence across both packages!
	stats2, err := passes.Run([]string{optRootPattern}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on Run 2, got %+v", stats2)
	}

	gotRhinoBytes, err := os.ReadFile(filepath.Join(optRhinoDir, "types.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotRhino := string(gotRhinoBytes)
	for _, want := range []string{
		"Slice_GetParameters() []FunctionType_Parameter",
		"func (this_ *FunctionTypeBase) Slice_GetParameters() []FunctionType_Parameter",
		"func (this_ *UnionType) Slice_GetAlternates() []JSType",
		"func (this_ *TemplatizedType) Slice_GetTemplateTypes() []JSType",
	} {
		if !strings.Contains(gotRhino, want) {
			t.Errorf("rhino missing expected pattern %q in:\n%s", want, gotRhino)
		}
	}

	gotJscompBytes, err := os.ReadFile(filepath.Join(optJscompDir, "check.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotJscomp := string(gotJscompBytes)
	for _, want := range []string{
		".Slice_GetAlternates()",
		".Slice_GetParameters()",
		".Slice_GetTemplateTypes()",
	} {
		if !strings.Contains(gotJscomp, want) {
			t.Errorf("jscomp missing expected cross-package slice call %q in:\n%s", want, gotJscomp)
		}
	}
	for _, unwanted := range []string{
		".IterSlice_()",
		"NewIterCursor(",
		"ImmListFromTypedSlice(",
		"ImmListToSliceNullable",
		"u.GetAlternates()",
		"fn.GetParameters()",
		"tt.GetTemplateTypes()",
	} {
		if strings.Contains(gotJscomp, unwanted) {
			t.Errorf("jscomp still contains unlowered pattern %q in:\n%s", unwanted, gotJscomp)
		}
	}

	// Execute differential & zero-allocation benchmark test in temp workspace!
	diffTestSrc := `package cross_pkg_test

import (
	"testing"

	orig_jscomp "example.com/w13maptest/orig/jscomp"
	orig_rhino "example.com/w13maptest/orig/rhino"
	opt_jscomp "example.com/w13maptest/opt/jscomp"
	opt_rhino "example.com/w13maptest/opt/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

func TestCrossPackageEquivalenceAndZeroAlloc(t *testing.T) {
	oa1 := orig_rhino.NewScalarType("num", false)
	oa2 := orig_rhino.NewScalarType("str", true)
	ou := orig_rhino.NewUnionType([]orig_rhino.JSType{oa1, oa2}, false)
	ofn := orig_rhino.NewFunctionType("fn", []orig_rhino.FunctionType_Parameter{
		orig_rhino.NewFunctionType_Parameter(oa1, true),
		orig_rhino.NewFunctionType_Parameter(oa2, false),
	})
	ott := orig_rhino.NewTemplatizedType([]orig_rhino.JSType{oa2})

	pa1 := opt_rhino.NewScalarType("num", false)
	pa2 := opt_rhino.NewScalarType("str", true)
	pu := opt_rhino.NewUnionType([]opt_rhino.JSType{pa1, pa2}, false)
	pfn := opt_rhino.NewFunctionType("fn", []opt_rhino.FunctionType_Parameter{
		opt_rhino.NewFunctionType_Parameter(pa1, true),
		opt_rhino.NewFunctionType_Parameter(pa2, false),
	})
	ptt := opt_rhino.NewTemplatizedType([]opt_rhino.JSType{pa2})

	want1 := orig_jscomp.CheckUnionAndFn(ou, ofn, ott)
	got1 := opt_jscomp.CheckUnionAndFn(pu, pfn, ptt)
	if got1 != want1 || pu.RebuildCount() != ou.RebuildCount() {
		t.Fatalf("pass 1 mismatch: got (%q, rebuild=%d), want (%q, rebuild=%d)", got1, pu.RebuildCount(), want1, ou.RebuildCount())
	}

	// Resolve oa1 / pa1 and verify lazy resolution inside Slice_GetAlternates() triggers identically!
	oa1.MarkResolved("number")
	pa1.MarkResolved("number")
	want2 := orig_jscomp.CheckUnionAndFn(ou, ofn, ott)
	got2 := opt_jscomp.CheckUnionAndFn(pu, pfn, ptt)
	if got2 != want2 || pu.RebuildCount() != ou.RebuildCount() || !pu.IsResolved() {
		t.Fatalf("pass 2 lazy-resolve mismatch: got (%q, rebuild=%d, resolved=%v), want (%q, rebuild=%d)",
			got2, pu.RebuildCount(), pu.IsResolved(), want2, ou.RebuildCount())
	}

	// Verify NoType subclass implementing FunctionType interface also works seamlessly.
	if orig_jscomp.CheckUnionAndFn(ou, orig_rhino.NewNoType(), ott) != opt_jscomp.CheckUnionAndFn(pu, opt_rhino.NewNoType(), ptt) {
		t.Fatalf("NoType FunctionType interface dispatch mismatch")
	}

	// Verify nil UnionType / nil FunctionType still throw Java NullPointerException!
	assertNPE := func(label string, fn func()) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("%s: expected NullPointerException, got nil", label)
			}
			if _, ok := r.(*javalang.NullPointerException); !ok {
				t.Fatalf("%s: expected *javalang.NullPointerException, got %T (%v)", label, r, r)
			}
		}()
		fn()
	}
	assertNPE("nil UnionType", func() { opt_jscomp.CheckUnionAndFn(nil, pfn, ptt) })
	assertNPE("typed-nil FunctionType", func() {
		var nilFn *opt_rhino.FunctionTypeBase
		opt_jscomp.CheckUnionAndFn(pu, nilFn, ptt)
	})

	// Verify 0 ImmListFromTypedSlice / 0 IterSlice_ allocations on hot iteration + indexing loop!
	allocs := testing.AllocsPerRun(200, func() {
		s := 0
		for _, alt := range pu.Slice_GetAlternates() {
			s += len(alt.GetDisplayName())
		}
		var iface opt_rhino.FunctionType = pfn
		for _, p := range iface.Slice_GetParameters() {
			if p.IsOptional() {
				s++
			}
		}
		if s == 0 {
			t.Fatal("unreachable")
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocations for Slice_GetAlternates + Slice_GetParameters iteration, got %v", allocs)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "cross_pkg_test.go"), []byte(diffTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test in overlay workspace failed: %v\n%s", err, string(out))
	}
}

func TestWave15NativeDequeLoweringAndRuntimeSemantics(t *testing.T) {
	tmpDir, _ := setupRealJ2gortOverlayWorkspace(t)

	origDir := filepath.Join(tmpDir, "origdeq")
	optDir := filepath.Join(tmpDir, "optdeq")
	for _, d := range []string{origDir, optDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	deqSrc := `package origdeq

import (
	"example.com/w13maptest/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

type Item struct {
	Name string
	Val  int32
}

func NewItem(name string, val int32) *Item { return &Item{Name: name, Val: val} }
func (i *Item) IsJavaNil_() bool           { return i == nil }
func (i *Item) Equals(o any) bool {
	x, ok := o.(*Item)
	return ok && x != nil && x.Name == i.Name && x.Val == i.Val
}
func (i *Item) HashCode() int32 { return int32(len(i.Name))*31 + i.Val }

// 1. Eagerly-initialized struct field (lowered to inline 32-byte value j2gort.Deque[T]!) +
//    nullable struct field (lowered to *j2gort.Deque[string]!) + Node stack field!
type FlowAnalyzer struct {
	workQueue        j2gort.Ext_java_util_Deque
	lazyTags         *j2gort.Ext_java_util_ArrayDeque
	exceptionHandler j2gort.Ext_java_util_Deque
}

func (this_ *FlowAnalyzer) initFields_FlowAnalyzer() {
	this_.workQueue = j2gort.NewExt_java_util_ArrayDeque()
	this_.lazyTags = nil
	this_.exceptionHandler = j2gort.NewExt_java_util_ArrayDeque()
}

func NewFlowAnalyzer() *FlowAnalyzer {
	f := &FlowAnalyzer{}
	f.initFields_FlowAnalyzer()
	return f
}

func (this_ *FlowAnalyzer) AddTag(tag string) {
	if this_.lazyTags == nil {
		this_.lazyTags = j2gort.NewExt_java_util_ArrayDeque()
	}
	this_.lazyTags.AddLast(j2gort.BoxString(tag))
}

func (this_ *FlowAnalyzer) PollTagOrNull() string {
	if this_.lazyTags == nil {
		return j2gort.NullString
	}
	return j2gort.CastString(this_.lazyTags.PollFirst())
}

func (this_ *FlowAnalyzer) ResetTags() {
	this_.lazyTags = nil
}

func (this_ *FlowAnalyzer) HasTags() bool {
	return this_.lazyTags != nil && !this_.lazyTags.IsEmpty()
}

func (this_ *FlowAnalyzer) PushExcNode(n nodestore.Node) {
	j2gort.DequePushNode(this_.exceptionHandler, uint32(n))
}

func (this_ *FlowAnalyzer) TopExcNode() nodestore.Node {
	return nodestore.Node(j2gort.DequePeekFirstNode(this_.exceptionHandler))
}

func (this_ *FlowAnalyzer) PopExcNode(expected nodestore.Node) bool {
	if nodestore.Node(j2gort.DequePeekFirstNode(this_.exceptionHandler)) != expected {
		return false
	}
	j2gort.DequePopNode(this_.exceptionHandler)
	return true
}

func (this_ *FlowAnalyzer) RunWorklist(seed []*Item) int32 {
	this_.workQueue.Clear()
	for _, it := range seed {
		this_.workQueue.AddLast(j2gort.PtrToAny(it))
	}
	var localStack *j2gort.Ext_java_util_ArrayDeque = j2gort.NewExt_java_util_ArrayDeque_1_int32(int32(4))
	_ = localStack

	var sum int32 = 0
	for !this_.workQueue.IsEmpty() {
		cur := j2gort.CastPtr[Item](this_.workQueue.PollFirst())
		sum += cur.Val
		localStack.Push(j2gort.PtrToAny(cur))
		if cur.Val > 10 {
			this_.workQueue.AddFirst(j2gort.PtrToAny(NewItem(cur.Name+"'", cur.Val-10)))
		}
	}

	// Iterate localStack via NewIterCursor, then drain with Pop / RemoveLast / PeekLast!
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(localStack))
		for cur.HasNext() {
			raw := cur.Next()
			x := j2gort.CastPtr[Item](raw)
			sum += int32(len(x.Name))
		}
	}
	if localStack.Contains(j2gort.PtrToAny(NewItem("b", 2))) {
		localStack.Remove_1_any(j2gort.PtrToAny(NewItem("b", 2)))
		sum += 100
	}
	for localStack.Size() > 1 {
		top := j2gort.CastPtr[Item](localStack.Pop())
		bot := j2gort.CastPtr[Item](localStack.RemoveLast())
		sum += top.Val - bot.Val
	}
	if last := j2gort.CastPtr[Item](localStack.Peek()); last != nil {
		sum += last.Val * 3
	}
	return sum
}

// 2. Save/restore pointer-alias pattern (like ReferenceCollector.blockStack) + OrderedSet -> Queue!
type ScopedCollector struct {
	blockStack *j2gort.Ext_java_util_ArrayDeque
}

func InitScopedCollector(this_ *ScopedCollector, root string) {
	this_.blockStack = j2gort.NewExt_java_util_ArrayDeque()
	this_.blockStack.AddLast(j2gort.BoxString(root))
}

func NewScopedCollector(root string) *ScopedCollector {
	c := &ScopedCollector{}
	InitScopedCollector(c, root)
	return c
}

func (this_ *ScopedCollector) WithNestedScope(inner string) string {
	var oldStack *j2gort.Ext_java_util_ArrayDeque = this_.blockStack
	this_.blockStack = j2gort.NewExt_java_util_ArrayDeque()
	this_.blockStack.AddLast(j2gort.CastString(oldStack.GetFirst()))
	this_.blockStack.AddLast(j2gort.BoxString(inner))
	res := j2gort.CastString(this_.blockStack.GetLast()) + ":" + j2gort.StrConcat(this_.blockStack.Size())
	this_.blockStack = oldStack
	return res + ":" + j2gort.CastString(this_.blockStack.GetLast())
}

func DrainSetViaQueue(names []string) string {
	var s *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	for _, n := range names {
		j2gort.SetAddStr(s, n)
	}
	var q j2gort.Ext_java_util_Queue = j2gort.NewExt_java_util_ArrayDeque_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](s))
	var out string = ""
	for !q.IsEmpty() {
		out += j2gort.CastString(q.Remove_0()) + ","
	}
	return out
}
`

	if err := os.WriteFile(filepath.Join(origDir, "deque_sample.go"), []byte(deqSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	optDeqSrc := strings.Replace(deqSrc, "package origdeq", "package optdeq", 1)
	if err := os.WriteFile(filepath.Join(optDir, "deque_sample.go"), []byte(optDeqSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir(optdeq) failed: %v", err)
	}
	if stats1.ByPass["native_deque_lowering"] < 15 {
		t.Fatalf("expected >= 15 native_deque_lowering sites, got %+v", stats1.ByPass)
	}

	// Strictly 0-site Run-2 idempotence!
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on Run 2, got %+v", stats2)
	}

	gotBytes, err := os.ReadFile(filepath.Join(optDir, "deque_sample.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"workQueue        j2gort.Deque[*Item]",
		"lazyTags         *j2gort.Deque[string]",
		"exceptionHandler j2gort.Deque[nodestore.Node]",
		"blockStack *j2gort.Deque[string]",
		"j2gort.NewDequeCap[*Item](int32(4))",
		"j2gort.DequeFromOrderedSet(",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("optdeq missing expected lowered pattern %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "NewExt_java_util_ArrayDeque") {
		t.Errorf("optdeq still contains unlowered NewExt_java_util_ArrayDeque:\n%s", got)
	}

	rtTestSrc := `package w15_deque_test

import (
	"reflect"
	"testing"

	"example.com/w13maptest/j2gort"
	"example.com/w13maptest/optdeq"
	"example.com/w13maptest/origdeq"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

func TestDequeStructLayoutZeroAllocAndRingWraparound(t *testing.T) {
	if sz := reflect.TypeOf(j2gort.Deque[*optdeq.Item]{}).Size(); sz != 64 {
		t.Fatalf("Deque[*Item] size = %d bytes, want 64 bytes (32B header + [4]*Item inline buffer)", sz)
	}
	if nsz := reflect.TypeOf(j2gort.Deque[nodestore.Node]{}).Size(); nsz != 48 {
		t.Fatalf("Deque[Node] size = %d bytes, want 48 bytes (32B header + [4]Node inline buffer)", nsz)
	}

	// 0 heap allocations when empty!
	emptyAllocs := testing.AllocsPerRun(200, func() {
		var d j2gort.Deque[*optdeq.Item]
		if !d.IsEmpty() || d.Size() != 0 || d.Peek() != nil || d.Poll() != nil {
			t.Fatal("unexpected non-empty zero Deque")
		}
	})
	if emptyAllocs != 0 {
		t.Fatalf("expected 0 allocs for empty Deque, got %v", emptyAllocs)
	}

	// 0 per-element boxing allocations on Node / string / pointer ring-buffer wraparound!
	var nd j2gort.Deque[nodestore.Node]
	nd.InitCap(8)
	wrapAllocs := testing.AllocsPerRun(500, func() {
		nd.AddLast(nodestore.Node(11))
		nd.Push(nodestore.Node(22))
		_ = nd.Pop()
		_ = nd.PollFirst()
	})
	if wrapAllocs != 0 {
		t.Fatalf("expected 0 per-element allocs on pre-sized Node Deque, got %v", wrapAllocs)
	}

	// Stress ring-buffer wraparound (head > tail) + power-of-2 growth (8 -> 16 -> 32) + deleteAtOffset!
	var ring j2gort.Deque[int32]
	for round := int32(0); round < 20; round++ {
		ring.AddLast(round * 10)
		ring.AddLast(round*10 + 1)
		_ = ring.RemoveFirst() // advances head steadily around mask!
	}
	// Now ring has 20 elements wrapped across grown buffer.
	if ring.Size() != 20 {
		t.Fatalf("ring.Size() = %d, want 20", ring.Size())
	}
	for i := int32(0); i < ring.Size(); i++ {
		want := int32(20) + (i/2)*10 + (i % 2)
		// Verify monotonic sequence via At(i)
		_ = want
		if i > 0 && ring.At(i) <= ring.At(i-1) {
			t.Fatalf("ring order corrupted at %d: %d <= %d", i, ring.At(i), ring.At(i-1))
		}
	}
	// Remove middle elements from both front-half and back-half + Iterator.Remove()
	vFront := ring.At(2)
	vBack := ring.At(16)
	if !ring.RemoveFirstOccurrence(vFront) || ring.Contains(vFront) {
		t.Fatalf("failed RemoveFirstOccurrence(%d)", vFront)
	}
	if !ring.RemoveLastOccurrence(vBack) || ring.Contains(vBack) {
		t.Fatalf("failed RemoveLastOccurrence(%d)", vBack)
	}
	it := j2gort.DequeIterator(&ring)
	first := j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](it.Next()))
	it.Remove()
	if ring.Contains(first) || ring.Size() != 17 {
		t.Fatalf("Iterator.Remove failed: size=%d", ring.Size())
	}

	// Exact Java NPE vs NoSuchElementException checks!
	assertPanicsWith := func(label string, wantNPE bool, fn func()) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("%s: expected panic, got nil", label)
			}
			if wantNPE {
				if _, ok := r.(*javalang.NullPointerException); !ok {
					t.Fatalf("%s: expected NullPointerException, got %T (%v)", label, r, r)
				}
			} else {
				if _, ok := r.(*javalang.NoSuchElementException); !ok {
					t.Fatalf("%s: expected NoSuchElementException, got %T (%v)", label, r, r)
				}
			}
		}()
		fn()
	}

	var sd j2gort.Deque[string]
	if sd.Peek() != j2gort.NullString || sd.PollLast() != j2gort.NullString {
		t.Fatalf("empty string Deque Peek/Poll must return NullString")
	}
	assertPanicsWith("Pop empty", false, func() { sd.Pop() })
	assertPanicsWith("GetFirst empty", false, func() { sd.GetFirst() })
	assertPanicsWith("GetLast empty", false, func() { sd.GetLast() })
	assertPanicsWith("RemoveLast empty", false, func() { sd.RemoveLast() })
	assertPanicsWith("Push NullString", true, func() { sd.Push(j2gort.NullString) })
	assertPanicsWith("Push Node(0)", true, func() { nd.Push(nodestore.Node(0)) })
	var pd j2gort.Deque[*optdeq.Item]
	assertPanicsWith("AddLast nil ptr", true, func() { pd.AddLast(nil) })
	var nilDequePtr *j2gort.Deque[string]
	assertPanicsWith("nil *Deque.Size()", true, func() { nilDequePtr.Size() })
}

func TestDequeDifferentialEquivalence(t *testing.T) {
	faOrig := origdeq.NewFlowAnalyzer()
	faOpt := optdeq.NewFlowAnalyzer()

	if faOrig.PollTagOrNull() != faOpt.PollTagOrNull() || faOrig.HasTags() != faOpt.HasTags() {
		t.Fatal("initial nullable lazyTags mismatch")
	}
	faOrig.AddTag("alpha")
	faOrig.AddTag("beta")
	faOpt.AddTag("alpha")
	faOpt.AddTag("beta")
	if faOrig.HasTags() != faOpt.HasTags() || faOrig.PollTagOrNull() != faOpt.PollTagOrNull() {
		t.Fatal("populated lazyTags mismatch")
	}
	faOrig.ResetTags()
	faOpt.ResetTags()
	if faOrig.HasTags() != faOpt.HasTags() || faOrig.PollTagOrNull() != faOpt.PollTagOrNull() {
		t.Fatal("reset lazyTags mismatch")
	}

	faOrig.PushExcNode(nodestore.Node(5))
	faOrig.PushExcNode(nodestore.Node(9))
	faOpt.PushExcNode(nodestore.Node(5))
	faOpt.PushExcNode(nodestore.Node(9))
	if faOrig.TopExcNode() != faOpt.TopExcNode() || faOrig.PopExcNode(nodestore.Node(9)) != faOpt.PopExcNode(nodestore.Node(9)) || faOrig.TopExcNode() != faOpt.TopExcNode() {
		t.Fatal("Node stack operations mismatch")
	}

	origSeed := []*origdeq.Item{origdeq.NewItem("a", 15), origdeq.NewItem("b", 2), origdeq.NewItem("c", 23), origdeq.NewItem("d", 7)}
	optSeed := []*optdeq.Item{optdeq.NewItem("a", 15), optdeq.NewItem("b", 2), optdeq.NewItem("c", 23), optdeq.NewItem("d", 7)}
	if got, want := faOpt.RunWorklist(optSeed), faOrig.RunWorklist(origSeed); got != want {
		t.Fatalf("RunWorklist mismatch: got %d, want %d", got, want)
	}

	scOrig := origdeq.NewScopedCollector("global")
	scOpt := optdeq.NewScopedCollector("global")
	if got, want := scOpt.WithNestedScope("fn"), scOrig.WithNestedScope("fn"); got != want {
		t.Fatalf("WithNestedScope mismatch: got %q, want %q", got, want)
	}

	inputs := []string{"x", "y", "x", "z", "a"}
	if got, want := optdeq.DrainSetViaQueue(inputs), origdeq.DrainSetViaQueue(inputs); got != want {
		t.Fatalf("DrainSetViaQueue mismatch: got %q, want %q", got, want)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w15_deque_test.go"), []byte(rtTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestDeque", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test TestDeque failed: %v\n%s", err, string(out))
	}
}

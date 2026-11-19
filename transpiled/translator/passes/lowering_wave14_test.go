package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave14StructSlimmingAndNullablePointers(t *testing.T) {
	tmpDir, _ := setupRealJ2gortOverlayWorkspace(t)

	origDir := filepath.Join(tmpDir, "orig")
	optDir := filepath.Join(tmpDir, "opt")
	if err := os.MkdirAll(origDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(optDir, 0o755); err != nil {
		t.Fatal(err)
	}

	src := `package orig

import "example.com/w13maptest/j2gort"

type Var struct {
	Name string
	Slot int32
}

func NewVar(name string, slot int32) *Var { return &Var{Name: name, Slot: slot} }
func (v *Var) IsJavaNil_() bool           { return v == nil }
func (v *Var) HashCode() int32            { return int32(len(v.Name))*31 + v.Slot }
func (v *Var) Equals(o any) bool {
	x, ok := o.(*Var)
	return ok && x != nil && x.Name == v.Name && x.Slot == v.Slot
}

// AbstractScopeBase: 90% of Scopes only populate var0/var1 and leave vars == nil!
type AbstractScopeBase struct {
	var0 *Var
	var1 *Var
	vars *j2gort.Ext_java_util_LinkedHashMap
}

func (this_ *AbstractScopeBase) initFields_AbstractScopeBase() {
	this_.vars = nil
}

func NewAbstractScope() *AbstractScopeBase {
	s := &AbstractScopeBase{}
	s.initFields_AbstractScopeBase()
	return s
}

func (this_ *AbstractScopeBase) Declare(v *Var) {
	if this_.var0 == nil {
		this_.var0 = v
		return
	}
	if this_.var1 == nil {
		this_.var1 = v
		return
	}
	if this_.vars == nil {
		this_.vars = j2gort.NewExt_java_util_LinkedHashMap()
		j2gort.MapPutStr(this_.vars, this_.var0.Name, this_.var0)
		j2gort.MapPutStr(this_.vars, this_.var1.Name, this_.var1)
	}
	j2gort.MapPutStr(this_.vars, v.Name, v)
}

func (this_ *AbstractScopeBase) GetOwnVar(name string) *Var {
	if this_.vars != nil {
		return j2gort.CastPtr[Var](j2gort.MapGetStr(this_.vars, name))
	}
	if this_.var0 != nil && this_.var0.Name == name {
		return this_.var0
	}
	if this_.var1 != nil && this_.var1.Name == name {
		return this_.var1
	}
	return nil
}

func (this_ *AbstractScopeBase) ResetVars() {
	this_.var0 = nil
	this_.var1 = nil
	this_.vars = nil
}

func (this_ *AbstractScopeBase) HasOverflowMap() bool {
	return this_.vars != nil
}

// TypedScope: nullable set field initialized via ternary IIFE + CheckStateFmt(PtrToAny(...))
type TypedScope struct {
	AbstractScopeBase
	depth         int32
	reservedNames j2gort.Ext_java_util_Set
}

func InitTypedScope(this_ *TypedScope, parent *TypedScope, depth int32) {
	this_.initFields_AbstractScopeBase()
	this_.depth = depth
	this_.reservedNames = func() j2gort.Ext_java_util_Set {
		if parent == nil {
			return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.NewExt_java_util_LinkedHashSet())
		}
		return nil
	}()
}

func NewTypedScope(parent *TypedScope, depth int32) *TypedScope {
	t := &TypedScope{}
	InitTypedScope(t, parent, depth)
	return t
}

func (this_ *TypedScope) Reserve(name string) {
	j2gort.CheckStateFmt(this_.reservedNames != nil, "Not root scope: %s", j2gort.ToAny(this_.reservedNames))
	j2gort.SetAddStr(this_.reservedNames, name)
}

func (this_ *TypedScope) IsReserved(name string) bool {
	if this_.reservedNames == nil {
		return false
	}
	return j2gort.SetContainsStr(this_.reservedNames, name)
}

// MustDef: eagerly-initialized non-nullable struct field stays value CustomOrderedMap!
type MustDef struct {
	reachingDef j2gort.Ext_java_util_Map
}

func NewMustDef() *MustDef {
	m := &MustDef{}
	InitMustDef(m)
	return m
}

func InitMustDef(this_ *MustDef) {
	this_.reachingDef = j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.NewExt_java_util_LinkedHashMap())
}

func (this_ *MustDef) Record(v *Var, def string) {
	this_.reachingDef.Put_2(v, j2gort.BoxString(def))
}

func (this_ *MustDef) Lookup(v *Var) string {
	return j2gort.CastString(this_.reachingDef.Get(v))
}

// LocalSmallSet: standalone non-nullable local set uses stack LocalOrderedSet!
func DedupUpToFour(a, b, c string) int32 {
	var seen *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	_ = seen
	j2gort.SetAddStr(seen, a)
	j2gort.SetAddStr(seen, b)
	j2gort.SetAddStr(seen, c)
	var count int32 = 0
	{
		cur := j2gort.NewIterCursor(seen)
		for cur.HasNext() {
			s := cur.NextStr()
			count += int32(len(s))
		}
	}
	return seen.Size()*100 + count
}
`
	if err := os.WriteFile(filepath.Join(origDir, "scope.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	optSrc := strings.Replace(src, "package orig", "package opt", 1)
	if err := os.WriteFile(filepath.Join(optDir, "scope.go"), []byte(optSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir(opt) failed: %v", err)
	}
	if stats1.TotalSites == 0 {
		t.Fatalf("expected transformations, got 0")
	}

	// Run 2 idempotence check
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on Run 2, got %+v", stats2)
	}

	gotBytes, err := os.ReadFile(filepath.Join(optDir, "scope.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"vars *j2gort.OrderedMap[string, *Var]",
		"this_.vars = nil",
		"if this_.vars == nil {",
		"if this_.vars != nil {",
		"this_.vars = &j2gort.OrderedMap[string, *Var]{}",
		"reservedNames *j2gort.OrderedSet[string]",
		"return &j2gort.OrderedSet[string]{}",
		"reachingDef j2gort.CustomOrderedMap[*Var, string]",
		"var seen j2gort.LocalOrderedSet[string]",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing expected lowered pattern %q in:\n%s", want, got)
		}
	}

	sizeAndParityTest := `package w14_slim_test

import (
	"fmt"
	"reflect"
	"testing"

	"example.com/w13maptest/j2gort"
	"example.com/w13maptest/opt"
	"example.com/w13maptest/orig"
)

func TestStructByteSizesAndZeroAllocs(t *testing.T) {
	if sz := reflect.TypeFor[j2gort.OrderedMap[string, int32]]().Size(); sz != 40 {
		t.Fatalf("OrderedMap byte size = %d, want 40", sz)
	}
	if sz := reflect.TypeFor[j2gort.CustomOrderedMap[*opt.Var, string]]().Size(); sz != 64 {
		t.Fatalf("CustomOrderedMap byte size = %d, want 64", sz)
	}
	if sz := reflect.TypeFor[j2gort.OrderedSet[string]]().Size(); sz != 40 {
		t.Fatalf("OrderedSet byte size = %d, want 40", sz)
	}
	// AbstractScopeBase: var0 (8) + var1 (8) + vars (*OrderedMap = 8) = 24 bytes total!
	if sz := reflect.TypeFor[opt.AbstractScopeBase]().Size(); sz != 24 {
		t.Fatalf("AbstractScopeBase byte size = %d, want 24", sz)
	}
	// TypedScope: AbstractScopeBase (24) + depth (8 padded) + reservedNames (*OrderedSet = 8) = 40 bytes!
	if sz := reflect.TypeFor[opt.TypedScope]().Size(); sz != 40 {
		t.Fatalf("TypedScope byte size = %d, want 40", sz)
	}

	// LocalOrderedSet with <= 4 strings must do 0 heap allocations!
	if allocs := testing.AllocsPerRun(100, func() {
		_ = opt.DedupUpToFour("alpha", "beta", "alpha")
	}); allocs != 0 {
		t.Fatalf("DedupUpToFour allocs = %v, want 0", allocs)
	}
}

func TestScopeAndTypedScopeParity(t *testing.T) {
	os, xs := orig.NewAbstractScope(), opt.NewAbstractScope()
	if os.HasOverflowMap() || xs.HasOverflowMap() {
		t.Fatalf("fresh scope should not have overflow map")
	}
	for i := 0; i < 6; i++ {
		nm := fmt.Sprintf("v%d", i)
		os.Declare(orig.NewVar(nm, int32(i)))
		xs.Declare(opt.NewVar(nm, int32(i)))
		if os.HasOverflowMap() != xs.HasOverflowMap() {
			t.Fatalf("overflow state diverged at %d", i)
		}
	}
	for i := 0; i < 8; i++ {
		nm := fmt.Sprintf("v%d", i)
		ov, xv := os.GetOwnVar(nm), xs.GetOwnVar(nm)
		if (ov == nil) != (xv == nil) || (ov != nil && ov.Slot != xv.Slot) {
			t.Fatalf("lookup mismatch for %s: %v vs %v", nm, ov, xv)
		}
	}
	os.ResetVars()
	xs.ResetVars()
	if os.HasOverflowMap() || xs.HasOverflowMap() {
		t.Fatalf("reset scope should clear vars back to nil")
	}

	oRoot, xRoot := orig.NewTypedScope(nil, 0), opt.NewTypedScope(nil, 0)
	oChild, xChild := orig.NewTypedScope(oRoot, 1), opt.NewTypedScope(xRoot, 1)
	oRoot.Reserve("arguments")
	xRoot.Reserve("arguments")
	if !oRoot.IsReserved("arguments") || !xRoot.IsReserved("arguments") || oChild.IsReserved("arguments") || xChild.IsReserved("arguments") {
		t.Fatalf("reservedNames parity mismatch")
	}

	om, xm := orig.NewMustDef(), opt.NewMustDef()
	om.Record(orig.NewVar("a", 1), "rhs_a")
	xm.Record(opt.NewVar("a", 1), "rhs_a")
	if om.Lookup(orig.NewVar("a", 1)) != xm.Lookup(opt.NewVar("a", 1)) {
		t.Fatalf("MustDef lookup mismatch")
	}
	if orig.DedupUpToFour("x", "y", "x") != opt.DedupUpToFour("x", "y", "x") {
		t.Fatalf("DedupUpToFour mismatch")
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w14_slim_test.go"), []byte(sizeAndParityTest), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-count=1", "-race", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test failed: %v\n%s", err, string(out))
	}
}

func TestWave14GenericStructAndReceiverLowering(t *testing.T) {
	tmpDir, _ := setupRealJ2gortOverlayWorkspace(t)

	optDir := filepath.Join(tmpDir, "genpkg")
	if err := os.MkdirAll(optDir, 0o755); err != nil {
		t.Fatal(err)
	}

	src := `package genpkg

import "example.com/w13maptest/j2gort"

type GenericTag[K any] struct {
	Key K
}

type GenericContainer[K comparable, V any] struct {
	items   j2gort.Ext_java_util_ArrayList
	seen    *j2gort.Ext_java_util_LinkedHashSet
	idMap   *j2gort.Ext_java_util_IdentityHashMap
	immVals j2gort.Ext_com_google_common_collect_ImmutableList
}

func (this_ *GenericContainer[K, V]) initFields_GenericContainer() {
	this_.items = j2gort.NewExt_java_util_ArrayList()
	this_.seen = nil
	this_.idMap = nil
	this_.immVals = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
}

func NewGenericContainer[K comparable, V any]() *GenericContainer[K, V] {
	c := &GenericContainer[K, V]{}
	c.initFields_GenericContainer()
	return c
}

func (this_ *GenericContainer[K, V]) Push(tag *GenericTag[K], k K, v V) {
	this_.items.Add_1(j2gort.ToAny(v))
	if this_.seen == nil {
		this_.seen = j2gort.NewExt_java_util_LinkedHashSet()
	}
	this_.seen.Add_1(j2gort.PtrToAny(tag))
	if this_.idMap == nil {
		this_.idMap = j2gort.NewExt_java_util_IdentityHashMap()
	}
	this_.idMap.Put_2(j2gort.ToAny(k), j2gort.ToAny(v))
	this_.immVals = j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.ToAny(v))
}

func (this_ *GenericContainer[K, V]) FirstAndLookup(tag *GenericTag[K], k K) (V, bool, V, int32) {
	first := j2gort.Cast[V](this_.items.Get(0))
	has := this_.seen != nil && this_.seen.Contains(j2gort.PtrToAny(tag))
	mapped := j2gort.Cast[V](this_.idMap.Get(j2gort.ToAny(k)))
	immLen := this_.immVals.Size() + int32(len(func() []V { return []V{j2gort.Cast[V](this_.immVals.Get(0))} }()))
	return first, has, mapped, immLen
}

func UseInstantiated() (string, bool, string, int32) {
	c := NewGenericContainer[int32, string]()
	tag := &GenericTag[int32]{Key: 7}
	c.Push(tag, 7, "hello")
	c.items.Add_1(j2gort.BoxString("world"))
	return c.FirstAndLookup(tag, 7)
}
`
	if err := os.WriteFile(filepath.Join(optDir, "generic.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir(genpkg) failed: %v", err)
	}
	for _, passName := range []string{"native_slice_lowering", "native_set_lowering", "native_map_lowering", "immutable_list_lowering"} {
		if stats1.ByPass[passName] == 0 {
			t.Fatalf("expected %s to transform generic struct fields, got %+v", passName, stats1.ByPass)
		}
	}

	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 {
		t.Fatalf("expected 0 sites on Run 2, got %+v", stats2)
	}

	gotBytes, err := os.ReadFile(filepath.Join(optDir, "generic.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"items   []V",
		"seen    *j2gort.OrderedSet[*GenericTag[K]]",
		"idMap   *j2gort.IdentityOrderedMap[K, V]",
		"immVals []V",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing lowered generic field %q in:\n%s", want, got)
		}
	}

	testSrc := `package w14_gen_test

import (
	"testing"
	"example.com/w13maptest/genpkg"
)

func TestGenericContainerExecution(t *testing.T) {
	first, has, mapped, immLen := genpkg.UseInstantiated()
	if first != "hello" || !has || mapped != "hello" || immLen != 2 {
		t.Fatalf("unexpected generic execution result: (%q, %v, %q, %d)", first, has, mapped, immLen)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w14_gen_test.go"), []byte(testSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-count=1", "-race", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test failed: %v\n%s", err, string(out))
	}
}

func TestWave14ImmutableListLowering_FunctionTypeValueStructAndDifferentialParity(t *testing.T) {
	tmpDir, _ := setupRealJ2gortOverlayWorkspace(t)

	origDir := filepath.Join(tmpDir, "orig")
	optDir := filepath.Join(tmpDir, "opt")
	if err := os.MkdirAll(origDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(optDir, 0o755); err != nil {
		t.Fatal(err)
	}

	src := `package orig

import (
	"example.com/w13maptest/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

type Node uint32

func (n Node) RawNodeId_() uint32 { return uint32(n) }

type JSType interface {
	IsJavaNil_() bool
	TypeName() string
}

type SimpleType struct {
	name string
}

func NewSimpleType(name string) JSType { return &SimpleType{name: name} }
func (s *SimpleType) IsJavaNil_() bool { return s == nil }
func (s *SimpleType) TypeName() string { return s.name }

// 24-byte Go value struct representing lowered rhino.FunctionType_Parameter!
type FunctionType_Parameter struct {
	jsType     JSType
	isOptional bool
	isVariadic bool
}

func NewParameter(t JSType, opt, variadic bool) FunctionType_Parameter {
	return FunctionType_Parameter{jsType: t, isOptional: opt, isVariadic: variadic}
}

func (p FunctionType_Parameter) GetJSType() JSType { return p.jsType }
func (p FunctionType_Parameter) IsOptional() bool  { return p.isOptional }
func (p FunctionType_Parameter) IsVariadic() bool  { return p.isVariadic }
func (p FunctionType_Parameter) Equals(other any) bool {
	o, ok := other.(FunctionType_Parameter)
	if !ok {
		return false
	}
	return javalang.ObjectsEquals(p.jsType, o.jsType) && p.isOptional == o.isOptional && p.isVariadic == o.isVariadic
}
func (p FunctionType_Parameter) HashCode() int32 {
	var h int32 = 1
	if p.jsType != nil && !p.jsType.IsJavaNil_() {
		h = 31*h + int32(len(p.jsType.TypeName()))
	}
	if p.isOptional {
		h += 3
	}
	if p.isVariadic {
		h += 7
	}
	return h
}

type FunctionType interface {
	JSType
	GetFunctionTypeBase() *FunctionTypeBase
	GetParameters() j2gort.Ext_com_google_common_collect_ImmutableList
	GetMinArity() int32
}

type FunctionTypeBase struct {
	self_         FunctionType
	name          string
	parameterList j2gort.Ext_com_google_common_collect_ImmutableList
	extraTags     j2gort.Ext_com_google_common_collect_ImmutableList
}

type FunctionTypeImpl struct {
	FunctionTypeBase
}

func (f *FunctionTypeBase) IsJavaNil_() bool                      { return f == nil }
func (f *FunctionTypeBase) TypeName() string                      { return f.name }
func (f *FunctionTypeBase) GetFunctionTypeBase() *FunctionTypeBase { return f }

func NewFunctionTypeFromBuilder(name string, params []FunctionType_Parameter, useExternalFallback bool) FunctionType {
	impl := &FunctionTypeImpl{}
	impl.self_ = impl
	impl.name = name
	var tmpList j2gort.Ext_java_util_ArrayList = nil
	if !useExternalFallback {
		tmpList = j2gort.NewExt_java_util_ArrayList()
		for _, p := range params {
			tmpList.Add_1(j2gort.ToAny(p))
		}
	}
	impl.parameterList = func() j2gort.Ext_com_google_common_collect_ImmutableList {
		if tmpList == nil {
			b := &j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}
			for _, p := range params {
				b.AddVoid_(j2gort.ToAny(p))
			}
			return ExternalIdentityPass(b.BuildImmList_())
		}
		return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](tmpList))
	}()
	impl.extraTags = j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.BoxString("fn"), j2gort.BoxString(name))
	return impl
}

func ExternalIdentityPass(in j2gort.Ext_com_google_common_collect_ImmutableList) j2gort.Ext_com_google_common_collect_ImmutableList {
	return in
}

func (this_ *FunctionTypeBase) GetParameters() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.parameterList
}

func (this_ *FunctionTypeBase) GetMinArity() int32 {
	var min_ int32 = 0
	var size int32 = this_.parameterList.Size()
	for i := int32(0); i < size; i++ {
		var p FunctionType_Parameter = j2gort.Cast[FunctionType_Parameter](this_.parameterList.Get(i))
		if !p.IsOptional() && !p.IsVariadic() {
			min_ = i + 1
		}
	}
	return min_
}

// Iterator loop on this_.self_.GetParameters().Iterator() -> SliceIterator -> native for _, param := range!
func (this_ *FunctionTypeBase) FormatSignature() string {
	var out string = this_.name + "("
	var it j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_com_google_common_collect_UnmodifiableIterator](this_.self_.GetParameters().Iterator())
	_ = it
	first := true
	for it.HasNext() {
		var param FunctionType_Parameter = j2gort.Cast[FunctionType_Parameter](it.Next())
		_ = param
		if !first {
			out += ","
		}
		first = false
		out += param.GetJSType().TypeName()
	}
	var tagsSummary string = ""
	for _, rawTag := range j2gort.CheckNotNull(this_.extraTags).IterSlice_() {
		tag := j2gort.CastString(rawTag)
		tagsSummary += ":" + tag
	}
	return out + ")" + tagsSummary
}

// Same-package SubtypeChecker hot-path on FunctionType interface: 0 wrapper allocs, 0 element boxing!
func SubtypeCheckParams(sub FunctionType, sup FunctionType) bool {
	var subParams j2gort.Ext_com_google_common_collect_ImmutableList = sub.GetParameters()
	var supParams j2gort.Ext_com_google_common_collect_ImmutableList = sup.GetParameters()
	if subParams.Size() != supParams.Size() {
		return false
	}
	for i := int32(0); i < subParams.Size(); i++ {
		sp := j2gort.Cast[FunctionType_Parameter](subParams.Get(i))
		tp := j2gort.Cast[FunctionType_Parameter](supParams.Get(i))
		if sp.IsOptional() != tp.IsOptional() || sp.IsVariadic() != tp.IsVariadic() || sp.GetJSType().TypeName() != tp.GetJSType().TypeName() {
			return false
		}
	}
	return true
}

// Same-package caller exercising Equals, Contains, IndexOf, HashCode on value-struct elements!
func CheckParamArityMatch(sub FunctionType, sup FunctionType, probe FunctionType_Parameter) (bool, bool, int32, int32) {
	var subParams j2gort.Ext_com_google_common_collect_ImmutableList = sub.GetParameters()
	var supParams j2gort.Ext_com_google_common_collect_ImmutableList = sup.GetParameters()
	if subParams.IsEmpty() && supParams.IsEmpty() {
		return true, false, subParams.HashCode(), 0
	}
	eq := subParams.Equals(supParams)
	hasProbe := subParams.Contains(j2gort.ToAny(probe))
	idx := subParams.IndexOf_1(j2gort.ToAny(probe))
	return eq, hasProbe, subParams.HashCode(), idx
}

func ExerciseBuildersAndStaticOf(a, b string, n1, n2 Node) (string, int32, bool) {
	var empty j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	var nilList j2gort.Ext_com_google_common_collect_ImmutableList = nil
	if empty == nil || nilList != nil {
		return "BAD_NIL", -1, false
	}

	var b1 *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.Ext_com_google_common_collect_ImmutableList_BuilderWithExpectedSize(4)
	_ = b1
	b1.AddStrVoid_(a)
	b1.AddVoid_(j2gort.BoxString(b))
	var built j2gort.Ext_com_google_common_collect_ImmutableList = b1.BuildImmList_()

	var b2 *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	b2.AddAllVoid_(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](built))
	b2.AddStrVoid_("tail")
	var combined j2gort.Ext_com_google_common_collect_ImmutableList = b2.BuildImmList_()

	var chained j2gort.Ext_com_google_common_collect_ImmutableList = func() j2gort.Ext_com_google_common_collect_ImmutableList {
		cb := &j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}
		return cb.AddSelf_(j2gort.BoxNode(n1)).AddSelf_(j2gort.BoxNode(n2)).BuildImmList_()
	}()

	var nodeSum int32 = 0
	for it := j2gort.AsIface[j2gort.Ext_java_util_Iterator](chained.Iterator()); it.HasNext(); {
		nd := Node(j2gort.CastNode[uint32](it.Next()))
		nodeSum += int32(nd)
	}

	rev := combined.Reverse()
	sub := j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](combined.SubList(0, 2))
	eq := sub.Equals(built) && javalang.ObjectsEquals(built, sub)
	summary := j2gort.ListGetStr(combined, 0) + "|" + j2gort.CastString(rev.Get(0)) + "|" + j2gort.StrVal(empty)
	return summary, nodeSum + empty.Size(), eq
}

// Negative soundness guard: reference-identity comparison l1 == l2 must NOT lower!
func IdentityComparedList() bool {
	var l1 j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	var l2 j2gort.Ext_com_google_common_collect_ImmutableList = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	return l1 == l2
}
`
	if err := os.WriteFile(filepath.Join(origDir, "fntype.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	optSrc := strings.Replace(src, "package orig", "package opt", 1)
	if err := os.WriteFile(filepath.Join(optDir, "fntype.go"), []byte(optSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir(opt) failed: %v", err)
	}
	if stats1.ByPass["immutable_list_lowering"] < 20 {
		t.Fatalf("expected >= 20 immutable_list_lowering sites, got %+v", stats1.ByPass)
	}

	// Run 2 idempotence check!
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on Run 2, got %+v", stats2)
	}

	gotBytes, err := os.ReadFile(filepath.Join(optDir, "fntype.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"parameterList []FunctionType_Parameter",
		"extraTags     []string",
		"return j2gort.ImmListFromTypedSlice(this_.parameterList)",
		"func (this_ *FunctionTypeBase) Slice_GetParameters() []FunctionType_Parameter",
		"Slice_GetParameters() []FunctionType_Parameter",
		"j2gort.CheckNotNullIface(sub).Slice_GetParameters()",
		"for _, param := range",
		"for _, nd := range chained",
		"for _, tag := range this_.extraTags",
		"var b1 []string = make([]string, 0, 4)",
		"return l1 == l2",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing expected snippet %q in lowered fntype.go:\n%s", want, got)
		}
	}

	diffTest := `package w14_immlist_test

import (
	"reflect"
	"testing"

	"example.com/w13maptest/j2gort"
	"example.com/w13maptest/opt"
	"example.com/w13maptest/orig"
)

func catchPanic(fn func()) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	fn()
	return false
}

func TestImmutableListFunctionTypeParityAndZeroAllocCheck(t *testing.T) {
	if sz := reflect.TypeFor[opt.FunctionType_Parameter]().Size(); sz != 24 {
		t.Fatalf("FunctionType_Parameter byte size = %d, want 24", sz)
	}

	otNum, otStr := orig.NewSimpleType("number"), orig.NewSimpleType("string")
	xtNum, xtStr := opt.NewSimpleType("number"), opt.NewSimpleType("string")

	for _, fallback := range []bool{false, true} {
		ofn1 := orig.NewFunctionTypeFromBuilder("fnA", []orig.FunctionType_Parameter{
			orig.NewParameter(otNum, false, false),
			orig.NewParameter(otStr, true, false),
		}, fallback)
		xfn1 := opt.NewFunctionTypeFromBuilder("fnA", []opt.FunctionType_Parameter{
			opt.NewParameter(xtNum, false, false),
			opt.NewParameter(xtStr, true, false),
		}, fallback)

		ofn2 := orig.NewFunctionTypeFromBuilder("fnB", []orig.FunctionType_Parameter{
			orig.NewParameter(otNum, false, false),
			orig.NewParameter(otStr, true, false),
		}, fallback)
		xfn2 := opt.NewFunctionTypeFromBuilder("fnB", []opt.FunctionType_Parameter{
			opt.NewParameter(xtNum, false, false),
			opt.NewParameter(xtStr, true, false),
		}, fallback)

		if ofn1.GetMinArity() != xfn1.GetMinArity() || ofn1.GetMinArity() != 1 {
			t.Fatalf("GetMinArity mismatch: %d vs %d", ofn1.GetMinArity(), xfn1.GetMinArity())
		}
		if ofn1.(*orig.FunctionTypeImpl).FormatSignature() != xfn1.(*opt.FunctionTypeImpl).FormatSignature() {
			t.Fatalf("FormatSignature mismatch: %q vs %q", ofn1.(*orig.FunctionTypeImpl).FormatSignature(), xfn1.(*opt.FunctionTypeImpl).FormatSignature())
		}
		// Distinct SimpleType pointer with identical name/fields tests value-struct Equals!
		oProbe := orig.NewParameter(otStr, true, false)
		xProbe := opt.NewParameter(xtStr, true, false)
		oEq, oHas, oHc, oIdx := orig.CheckParamArityMatch(ofn1, ofn2, oProbe)
		xEq, xHas, xHc, xIdx := opt.CheckParamArityMatch(xfn1, xfn2, xProbe)
		if oEq != xEq || !xEq || oHas != xHas || !xHas || oHc != xHc || oIdx != xIdx || xIdx != 1 {
			t.Fatalf("CheckParamArityMatch mismatch: orig=(%v,%v,%d,%d) opt=(%v,%v,%d,%d)", oEq, oHas, oHc, oIdx, xEq, xHas, xHc, xIdx)
		}

		// Escape via GetParameters() interface must behave identically!
		if ofn1.GetParameters().Size() != xfn1.GetParameters().Size() {
			t.Fatalf("escaped GetParameters().Size() mismatch")
		}

		// SubtypeChecker-style SubtypeCheckParams must allocate strictly 0 heap objects!
		if !orig.SubtypeCheckParams(ofn1, ofn2) || !opt.SubtypeCheckParams(xfn1, xfn2) {
			t.Fatalf("SubtypeCheckParams mismatch")
		}
		if allocs := testing.AllocsPerRun(100, func() {
			_ = opt.SubtypeCheckParams(xfn1, xfn2)
		}); allocs != 0 {
			t.Fatalf("SubtypeCheckParams allocs = %v, want 0", allocs)
		}
	}

	// Empty FunctionType GetParameters() escape must also allocate 0 objects (singleton emptyImmutableList_)!
	xEmptyFn := opt.NewFunctionTypeFromBuilder("emptyFn", nil, false)
	if allocs := testing.AllocsPerRun(100, func() {
		_ = xEmptyFn.GetParameters()
	}); allocs != 0 {
		t.Fatalf("empty GetParameters() escape allocs = %v, want 0", allocs)
	}

	os1, on1, oeq1 := orig.ExerciseBuildersAndStaticOf("hello", "world", 10, 22)
	xs1, xn1, xeq1 := opt.ExerciseBuildersAndStaticOf("hello", "world", 10, 22)
	if os1 != xs1 || on1 != xn1 || oeq1 != xeq1 || !xeq1 {
		t.Fatalf("ExerciseBuildersAndStaticOf mismatch: orig=(%q,%d,%v) opt=(%q,%d,%v)", os1, on1, oeq1, xs1, xn1, xeq1)
	}

	// Null-element NPE parity on NullString and Node(0)!
	if catchPanic(func() { _, _, _ = orig.ExerciseBuildersAndStaticOf(j2gort.NullString, "ok", 1, 2) }) !=
		catchPanic(func() { _, _, _ = opt.ExerciseBuildersAndStaticOf(j2gort.NullString, "ok", 1, 2) }) {
		t.Fatalf("NullString NPE parity mismatch")
	}
	if catchPanic(func() { _, _, _ = orig.ExerciseBuildersAndStaticOf("ok", "ok", 0, 2) }) !=
		catchPanic(func() { _, _, _ = opt.ExerciseBuildersAndStaticOf("ok", "ok", 0, 2) }) {
		t.Fatalf("Node(0) NPE parity mismatch")
	}

	if !orig.IdentityComparedList() || !opt.IdentityComparedList() {
		t.Fatalf("IdentityComparedList should return true")
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w14_immlist_test.go"), []byte(diffTest), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-count=1", "-race", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test failed: %v\n%s", err, string(out))
	}
}

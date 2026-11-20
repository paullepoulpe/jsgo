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

func wave43bWriteFiles(t *testing.T, root string, files map[string]string) {
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

func wave43bMustBuild(t *testing.T, root string) {
	t.Helper()
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("pre-pass go build ./... failed: %v\n%s", err, string(out))
	}
}

func wave43bJ2gortStubSrc() string {
	return `package j2gort

import "sync/atomic"

const NullString = "\xff\x00j2go_null"

var (
	OrderedMapAcquireCalls  atomic.Int64
	OrderedMapReleaseCalls  atomic.Int64
	PooledSliceAcquireCalls atomic.Int64
	PooledSliceReleaseCalls atomic.Int64
)

type Ext_java_util_Map interface {
	IsJavaNil_() bool
	Size() int32
	IsEmpty() bool
	ContainsKey(any) bool
	Get(any) any
	Put_2(any, any) any
	Remove_1(any) any
	Clear()
}

type Ext_java_util_LinkedHashMap struct {
	m map[any]any
}

func (m *Ext_java_util_LinkedHashMap) IsJavaNil_() bool { return m == nil }
func (m *Ext_java_util_LinkedHashMap) Size() int32      { return int32(len(m.m)) }
func (m *Ext_java_util_LinkedHashMap) IsEmpty() bool    { return m == nil || len(m.m) == 0 }
func (m *Ext_java_util_LinkedHashMap) ContainsKey(k any) bool {
	_, ok := m.m[k]
	return ok
}
func (m *Ext_java_util_LinkedHashMap) Get(k any) any { return m.m[k] }
func (m *Ext_java_util_LinkedHashMap) Put_2(k, v any) any {
	if m.m == nil {
		m.m = make(map[any]any)
	}
	old := m.m[k]
	m.m[k] = v
	return old
}
func (m *Ext_java_util_LinkedHashMap) Remove_1(k any) any {
	old := m.m[k]
	delete(m.m, k)
	return old
}
func (m *Ext_java_util_LinkedHashMap) Clear() { clear(m.m) }

type Ext_java_util_HashMap struct {
	Ext_java_util_LinkedHashMap
}

func NewExt_java_util_LinkedHashMap() *Ext_java_util_LinkedHashMap {
	return &Ext_java_util_LinkedHashMap{m: make(map[any]any)}
}
func NewExt_java_util_LinkedHashMap_1_int32(capHint int32) *Ext_java_util_LinkedHashMap {
	return &Ext_java_util_LinkedHashMap{m: make(map[any]any, capHint)}
}
func Ext_com_google_common_collect_Maps_NewLinkedHashMapWithExpectedSize(expected int32) *Ext_java_util_LinkedHashMap {
	return &Ext_java_util_LinkedHashMap{m: make(map[any]any, expected)}
}
func NewExt_java_util_HashMap() *Ext_java_util_HashMap {
	return &Ext_java_util_HashMap{Ext_java_util_LinkedHashMap{m: make(map[any]any)}}
}
func NewExt_java_util_HashMap_1_int32(capHint int32) *Ext_java_util_HashMap {
	return &Ext_java_util_HashMap{Ext_java_util_LinkedHashMap{m: make(map[any]any, capHint)}}
}

func AsIface[I any](v any) I {
	if v == nil {
		var z I
		return z
	}
	return v.(I)
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

func CheckNotNullPtr[T any](p *T) *T {
	if p == nil {
		panic("NPE")
	}
	return p
}

func CheckNotNullString(s string) string {
	if s == NullString {
		panic("NullString")
	}
	return s
}

func CheckNotNullIface[T any](v T) T { return v }

func CheckArgumentMsg(cond bool, _ any) {
	if !cond {
		panic("IllegalArgumentException")
	}
}

func CheckArgumentFmt(cond bool, _ string, _ ...any) {
	if !cond {
		panic("IllegalArgumentException")
	}
}

func ObjEqCustomPtr[T any](a, b *T) bool { return a == b }

func SliceRemoveAt[T any](s []T, idx int32) []T {
	copy(s[idx:], s[idx+1:])
	var z T
	s[len(s)-1] = z
	return s[:len(s)-1]
}

func AcquirePooledSlice[T any](capHint int) []T {
	PooledSliceAcquireCalls.Add(1)
	if capHint <= 0 {
		return nil
	}
	return make([]T, 0, capHint)
}

func AppendPooledSlice[T any](s *[]T, v T) {
	*s = append(*s, v)
}

func ReleasePooledSlice[T any](s []T) {
	if cap(s) > 0 {
		PooledSliceReleaseCalls.Add(1)
		clear(s[:cap(s)])
	}
}

type Slab[T any] struct{ items []T }

func (s *Slab[T]) Alloc() *T {
	s.items = append(s.items, *new(T))
	return &s.items[len(s.items)-1]
}
func (s *Slab[T]) Release() { clear(s.items); s.items = nil }

type PooledSlab[T any] struct {
	items        []*T
	AllocCalls   atomic.Int64
	ResetCalls   atomic.Int64
	ReleaseCalls atomic.Int64
}

func (s *PooledSlab[T]) Alloc() *T {
	s.AllocCalls.Add(1)
	v := new(T)
	s.items = append(s.items, v)
	return v
}

func (s *PooledSlab[T]) Reset() {
	s.ResetCalls.Add(1)
	for _, it := range s.items {
		var z T
		*it = z
	}
	s.items = s.items[:0]
}

func (s *PooledSlab[T]) Release() {
	s.ReleaseCalls.Add(1)
	for _, it := range s.items {
		var z T
		*it = z
	}
	s.items = nil
}

func (s *PooledSlab[T]) ReleaseWith(cleanup func(*T)) {
	s.ReleaseCalls.Add(1)
	for _, it := range s.items {
		if cleanup != nil {
			cleanup(it)
		}
		var z T
		*it = z
	}
	s.items = nil
}

type OrderedMap[K comparable, V any] struct {
	m       map[K]V
	capHint int32
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] { return OrderedMap[K, V]{} }
func NewOrderedMapCap[K comparable, V any](capHint int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V, capHint), capHint: capHint}
}
func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] { return &OrderedMap[K, V]{} }
func NewOrderedMapPtrCap[K comparable, V any](capHint int32) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{m: make(map[K]V, capHint), capHint: capHint}
}
func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	OrderedMapAcquireCalls.Add(1)
	return OrderedMap[K, V]{}
}
func AcquireOrderedMapCap[K comparable, V any](capHint int32) OrderedMap[K, V] {
	OrderedMapAcquireCalls.Add(1)
	return OrderedMap[K, V]{m: make(map[K]V, capHint), capHint: capHint}
}
func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	OrderedMapAcquireCalls.Add(1)
	return &OrderedMap[K, V]{}
}
func AcquireOrderedMapPtrCap[K comparable, V any](capHint int32) *OrderedMap[K, V] {
	OrderedMapAcquireCalls.Add(1)
	return &OrderedMap[K, V]{m: make(map[K]V, capHint), capHint: capHint}
}
func ReleaseOrderedMap[K comparable, V any](m *OrderedMap[K, V]) {
	if m != nil {
		OrderedMapReleaseCalls.Add(1)
		clear(m.m)
	}
}
func (m *OrderedMap[K, V]) CapHint() int32 { return m.capHint }
func (m *OrderedMap[K, V]) IsNil() bool    { return m == nil }
func (m *OrderedMap[K, V]) Size() int32 {
	if m == nil {
		return 0
	}
	return int32(len(m.m))
}
func (m *OrderedMap[K, V]) IsEmpty() bool { return m == nil || len(m.m) == 0 }
func (m *OrderedMap[K, V]) Get(k K) V {
	var z V
	if m == nil || m.m == nil {
		return z
	}
	return m.m[k]
}
func (m *OrderedMap[K, V]) Put(k K, v V) V {
	if m.m == nil {
		m.m = make(map[K]V)
	}
	old := m.m[k]
	m.m[k] = v
	return old
}
func (m *OrderedMap[K, V]) PutVoid(k K, v V) {
	if m.m == nil {
		m.m = make(map[K]V)
	}
	m.m[k] = v
}
func (m *OrderedMap[K, V]) Remove(k K) V {
	var z V
	if m == nil || m.m == nil {
		return z
	}
	old := m.m[k]
	delete(m.m, k)
	return old
}
func (m *OrderedMap[K, V]) Clear() {
	if m != nil && m.m != nil {
		clear(m.m)
	}
}
func (m *OrderedMap[K, V]) ReleaseToPool() {
	OrderedMapReleaseCalls.Add(1)
	if m != nil && m.m != nil {
		clear(m.m)
	}
}

type CustomOrderedMap[K comparable, V any] struct {
	m map[K]V
}

func NewCustomOrderedMap[K comparable, V any]() CustomOrderedMap[K, V] {
	return CustomOrderedMap[K, V]{}
}
func (m *CustomOrderedMap[K, V]) Get(k K) V {
	var z V
	if m == nil || m.m == nil {
		return z
	}
	return m.m[k]
}
func (m *CustomOrderedMap[K, V]) Put(k K, v V) V {
	if m.m == nil {
		m.m = make(map[K]V)
	}
	old := m.m[k]
	m.m[k] = v
	return old
}
func (m *CustomOrderedMap[K, V]) Clear() {
	if m != nil && m.m != nil {
		clear(m.m)
	}
}

type OrderedSet[E comparable] struct {
	m map[E]struct{}
}

func NewOrderedSet[E comparable]() OrderedSet[E] { return OrderedSet[E]{} }
func (s *OrderedSet[E]) Add(v E) bool {
	if s.m == nil {
		s.m = make(map[E]struct{})
	}
	_, existed := s.m[v]
	s.m[v] = struct{}{}
	return !existed
}
func (s *OrderedSet[E]) Contains(v E) bool {
	if s == nil || s.m == nil {
		return false
	}
	_, ok := s.m[v]
	return ok
}

type Deque[E any] struct {
	s []E
}

func (d *Deque[E]) IsEmpty() bool { return len(d.s) == 0 }
func (d *Deque[E]) Push(v E)      { d.s = append(d.s, v) }
func (d *Deque[E]) Pop() E {
	last := len(d.s) - 1
	v := d.s[last]
	var z E
	d.s[last] = z
	d.s = d.s[:last]
	return v
}
func (d *Deque[E]) Clear() {
	clear(d.s)
	d.s = d.s[:0]
}

type PlainMap[K comparable, V any] struct {
	m       map[K]V
	capHint int32
}

func NewPlainMap[K comparable, V any](capHint int) *PlainMap[K, V] {
	return &PlainMap[K, V]{m: make(map[K]V, capHint), capHint: int32(capHint)}
}
func NewPlainMapVal[K comparable, V any]() PlainMap[K, V] { return PlainMap[K, V]{} }
func NewPlainMapCap[K comparable, V any](capHint int32) PlainMap[K, V] {
	return PlainMap[K, V]{m: make(map[K]V, capHint), capHint: capHint}
}
func NewPlainMapPtr[K comparable, V any]() *PlainMap[K, V] { return &PlainMap[K, V]{} }
func NewPlainMapPtrCap[K comparable, V any](capHint int32) *PlainMap[K, V] {
	return &PlainMap[K, V]{m: make(map[K]V, capHint), capHint: capHint}
}
func (m *PlainMap[K, V]) CapHint() int32 { return m.capHint }
func (m *PlainMap[K, V]) IsNil() bool    { return m == nil }
func (m *PlainMap[K, V]) Size() int32 {
	if m == nil {
		return 0
	}
	return int32(len(m.m))
}
func (m *PlainMap[K, V]) IsEmpty() bool { return m == nil || len(m.m) == 0 }
func (m *PlainMap[K, V]) Get(k K) V {
	var z V
	if m == nil || m.m == nil {
		return z
	}
	return m.m[k]
}
func (m *PlainMap[K, V]) Put(k K, v V) V {
	if m.m == nil {
		m.m = make(map[K]V)
	}
	old := m.m[k]
	m.m[k] = v
	return old
}
func (m *PlainMap[K, V]) PutVoid(k K, v V) {
	if m.m == nil {
		m.m = make(map[K]V)
	}
	m.m[k] = v
}
func (m *PlainMap[K, V]) Clear() {
	if m != nil && m.m != nil {
		clear(m.m)
	}
}
`
}

// TestWave43B_InlineVariablesZeroAllocLazyNegativeAndScopeReset verifies Task 1:
//   - Unconditional per-scope slab_InlineVariables_InlineVarAnalysis.Reset() in doInlinesForScope
//     even when varToAliasRetryHandlersMap stays non-empty due to an outer-scope alias!
//   - Zero-allocation CreateLazyNegative using ReferenceCollection.ivSafeAliases_ (fitting in
//     trailing struct padding with zero size increase!) + 3 static sentinels + activeRefMap_.
//   - Exact predicate call counts & caching across multiple aliases of the same outer variable.
//   - Run 2 idempotency (0 sites, 0 modified files) and `go test -race` pass.
func TestWave43B_InlineVariablesZeroAllocLazyNegativeAndScopeReset(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	rhinoSrc := `package rhino

type Node uint32
`

	jscompSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w43b_iv/gen/j2gort"
	pkg_rhino "example.com/w43b_iv/gen/rhino"
)

type AbstractCompiler interface{ IsCompiler() }
type NodeTraversal struct {
	scope *Scope
}

func (t *NodeTraversal) GetScope() *Scope { return t.scope }

type Scope struct {
	vars []*Var
}

type Var struct {
	name            string
	kind            int // 0=inlineWellDef, 1=lazyNegKnownFalse, 2=lazyNegWellDefKnown, 3=lazyNegUnknown, 4=waitAlias
	alias           *Var
	preEvalWellDef  bool
	preEvalAssigned bool
}

type Reference struct {
	node pkg_rhino.Node
}

type ReferenceCollection struct {
	ref0          *Reference
	ref1          *Reference
	ref2          *Reference
	ref3          *Reference
	rest          []*Reference
	idHash_       atomic.Int32
}

var (
	WellDefEvalCalls      atomic.Int64
	AssignedOnceEvalCalls atomic.Int64
	WellDefByRef0         = map[*Reference]bool{}
	AssignedOnceByRef0    = map[*Reference]bool{}
)

func (this_ *ReferenceCollection) initFields_ReferenceCollection() {
	_ = this_
}

func NewReferenceCollection() *ReferenceCollection {
	this_ := &ReferenceCollection{}
	InitReferenceCollection(this_)
	return this_
}

func InitReferenceCollection(this_ *ReferenceCollection) {
	this_.initFields_ReferenceCollection()
}

func (rc *ReferenceCollection) Size() int32 {
	if rc.ref0 == nil {
		return 0
	}
	return 1
}

func (rc *ReferenceCollection) Get(i int32) *Reference { return rc.ref0 }

func (rc *ReferenceCollection) IsWellDefined() bool {
	WellDefEvalCalls.Add(1)
	return WellDefByRef0[rc.ref0]
}

func (rc *ReferenceCollection) IsAssignedOnceInLifetime() bool {
	AssignedOnceEvalCalls.Add(1)
	return AssignedOnceByRef0[rc.ref0]
}

type ReferenceMap interface {
	GetReferences(any) any
}

type ReferenceCollector struct {
	referenceMap j2gort.CustomOrderedMap[*Var, *ReferenceCollection]
}

func (rc *ReferenceCollector) GetReferences(v any) any {
	return rc.referenceMap.Get(j2gort.CastPtr[Var](v))
}

type SimpleMultimap[K comparable, V any] struct {
	m map[K][]V
}

func (sm *SimpleMultimap[K, V]) IsEmpty() bool { return len(sm.m) == 0 }
func (sm *SimpleMultimap[K, V]) Put(k K, v V) {
	if sm.m == nil {
		sm.m = make(map[K][]V)
	}
	sm.m[k] = append(sm.m[k], v)
}
func (sm *SimpleMultimap[K, V]) RemoveAll(k K) []V {
	vs := sm.m[k]
	delete(sm.m, k)
	return vs
}
func (sm *SimpleMultimap[K, V]) Clear() { clear(sm.m) }

type InlineVariables struct {
	compiler     AbstractCompiler
	scopes       []*Scope
	collector    *ReferenceCollector
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
	for _, sc := range this_.scopes {
		beh.priv_InliningBehavior_doInlinesForScope_3393(&NodeTraversal{scope: sc}, this_.collector)
	}
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

func (this_ *InlineVariables_InlineVarAnalysis) ShouldInline() bool             { return int32(this_.kind) >= int32(4) }
func (this_ *InlineVariables_InlineVarAnalysis) ShouldWaitForAliasedVar() bool  { return int32(this_.kind) == int32(3) }
func (this_ *InlineVariables_InlineVarAnalysis) GetAliasedVar() *Var            { return this_.aliasedVar }
func (this_ *InlineVariables_InlineVarAnalysis) IsSafeToInlineAliases() bool {
	switch this_.kind {
	case int8(0), int8(3):
		return false
	case int8(1):
		return true
	case int8(2):
		return this_.expert.IsWellDefinedAssignedOnce()
	default:
		return true
	}
}

func (this_ *InlineVariables_InlineVarAnalysis) PerformInline() {
	if int32(this_.kind) == int32(5) {
		this_.expert.InlineWellDefined(this_.initValue)
	}
}

type InlineVariables_VarExpert struct {
	v                                   *Var
	referenceInfo                       *ReferenceCollection
	behavior                            *InlineVariables_InliningBehavior
	fixedAnalysis                       *InlineVariables_InlineVarAnalysis
	declaration                         *Reference
	initialization                      *Reference
	idHash_                             atomic.Int32
	mayBeAParameterModifiedViaArguments bool
	isDeclaredOrInferredConstant        bool
	initializationKnown                 bool
	isWellDefined                       int8
	isAssignedOnceInLifetime            int8
	isWellDefinedAssignedOnce           int8
	isNeverAssigned                     int8
	hasValidDeclaration                 int8
	hasValidInitialization              int8
	allReferencesAreValid               int8
}

func (this_ *InlineVariables_VarExpert) initFields_InlineVariables_VarExpert() {
	this_.isWellDefined = int8(-1)
	this_.isAssignedOnceInLifetime = int8(-1)
	this_.isWellDefinedAssignedOnce = int8(-1)
	this_.isNeverAssigned = int8(-1)
	this_.hasValidDeclaration = int8(-1)
	this_.hasValidInitialization = int8(-1)
	this_.allReferencesAreValid = int8(-1)
}

func NewInlineVariables_VarExpert_1(fixedAnalysis *InlineVariables_InlineVarAnalysis) *InlineVariables_VarExpert {
	this_ := &InlineVariables_VarExpert{}
	this_.initFields_InlineVariables_VarExpert()
	this_.fixedAnalysis = fixedAnalysis
	return this_
}

func NewInlineVariables_VarExpert_5(behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, mayBeAParameterModifiedViaArguments bool, isDeclaredOrInferredConstant bool) *InlineVariables_VarExpert {
	this_ := &InlineVariables_VarExpert{}
	InitInlineVariables_VarExpert_5(this_, behavior, v, referenceInfo, mayBeAParameterModifiedViaArguments, isDeclaredOrInferredConstant)
	return this_
}

func InitInlineVariables_VarExpert_5(this_ *InlineVariables_VarExpert, behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, mayBeAParameterModifiedViaArguments bool, isDeclaredOrInferredConstant bool) {
	this_.initFields_InlineVariables_VarExpert()
	this_.v = v
	this_.referenceInfo = referenceInfo
	this_.behavior = behavior
	this_.mayBeAParameterModifiedViaArguments = mayBeAParameterModifiedViaArguments
	this_.isDeclaredOrInferredConstant = isDeclaredOrInferredConstant
}

func (this_ *InlineVariables_VarExpert) priv_VarExpert_isWellDefined_1() bool {
	if this_.isWellDefined != int8(-1) {
		return this_.isWellDefined == int8(1)
	}
	if this_.referenceInfo.IsWellDefined() {
		this_.isWellDefined = int8(1)
		return true
	}
	this_.isWellDefined = int8(0)
	return false
}

func (this_ *InlineVariables_VarExpert) priv_VarExpert_isAssignedOnceInLifetime_1() bool {
	if this_.isAssignedOnceInLifetime != int8(-1) {
		return this_.isAssignedOnceInLifetime == int8(1)
	}
	if this_.referenceInfo.IsAssignedOnceInLifetime() {
		this_.isAssignedOnceInLifetime = int8(1)
		return true
	}
	this_.isAssignedOnceInLifetime = int8(0)
	return false
}

func (this_ *InlineVariables_VarExpert) IsWellDefinedAssignedOnce() bool {
	if this_.isWellDefinedAssignedOnce != int8(-1) {
		return this_.isWellDefinedAssignedOnce == int8(1)
	}
	if this_.isWellDefined == int8(0) || this_.isAssignedOnceInLifetime == int8(0) {
		this_.isWellDefinedAssignedOnce = int8(0)
		return false
	}
	if this_.priv_VarExpert_isWellDefined_1() && this_.priv_VarExpert_isAssignedOnceInLifetime_1() {
		this_.isWellDefinedAssignedOnce = int8(1)
		return true
	}
	this_.isWellDefinedAssignedOnce = int8(0)
	return false
}

func (this_ *InlineVariables_VarExpert) InlineWellDefined(_ pkg_rhino.Node) {
	this_.behavior.outer_.inlinedNames = append(this_.behavior.outer_.inlinedNames, this_.v.name)
}

func (this_ *InlineVariables_VarExpert) Analyze() *InlineVariables_InlineVarAnalysis {
	if this_.fixedAnalysis != nil {
		return this_.fixedAnalysis
	}
	if this_.v.preEvalWellDef {
		_ = this_.priv_VarExpert_isWellDefined_1()
	}
	if this_.v.preEvalAssigned {
		_ = this_.priv_VarExpert_isAssignedOnceInLifetime_1()
	}
	switch this_.v.kind {
	case 0:
		return InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(this_, 1)
	case 1, 2, 3:
		return InlineVariables_InlineVarAnalysis_CreateLazyNegative(this_)
	case 4:
		return InlineVariables_InlineVarAnalysis_CreateWaitForAlias(this_.v.alias)
	}
	return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS()
}

func (this_ *InlineVariables_VarExpert) ReanalyzeAfterAliasedVar(aliasedVar *Var, aliasedVarAnalysis *InlineVariables_InlineVarAnalysis) *InlineVariables_InlineVarAnalysis {
	if j2gort.ObjEqCustomPtr(aliasedVar, this_.v.alias) && aliasedVarAnalysis.IsSafeToInlineAliases() && this_.IsWellDefinedAssignedOnce() {
		return InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(this_, 2)
	}
	return StaticGet_InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS()
}

type InlineVariables_InliningBehavior struct {
	outer_                            *InlineVariables
	currentScopeHandledVarAnalysesMap j2gort.CustomOrderedMap[*Var, *InlineVariables_InlineVarAnalysis]
	varToAliasRetryHandlersMap        SimpleMultimap[*Var, *InlineVariables_InliningBehavior_AliasInlineRetryHandler]
}

func NewInlineVariables_InliningBehavior(outer_ *InlineVariables) *InlineVariables_InliningBehavior {
	return &InlineVariables_InliningBehavior{outer_: outer_}
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

func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) HandleAliasedVarCompletion(aliasedVar *Var, aliasedVarAnalysis *InlineVariables_InlineVarAnalysis) {
	a := this_.expert.ReanalyzeAfterAliasedVar(aliasedVar, aliasedVarAnalysis)
	if a.ShouldInline() {
		a.PerformInline()
	}
	this_.outer_.currentScopeHandledVarAnalysesMap.Put(this_.v, a)
}

func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_createVarExpert_1(v *Var, referenceInfo *ReferenceCollection) *InlineVariables_VarExpert {
	if referenceInfo == nil {
		return StaticGet_NO_EXPERT()
	}
	return NewInlineVariables_VarExpert_5(this_, v, referenceInfo, false, false)
}

func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_doInlinesForScope_3393(t *NodeTraversal, referenceMap ReferenceMap) {
	_ = t
	_ = referenceMap
	this_.currentScopeHandledVarAnalysesMap.Clear()
	for _, v := range t.GetScope().vars {
		var referenceInfo *ReferenceCollection = j2gort.CastPtr[ReferenceCollection](referenceMap.GetReferences(j2gort.PtrToAny(v)))
		expert := this_.priv_InliningBehavior_createVarExpert_1(v, referenceInfo)
		analysis := expert.Analyze()
		if analysis.ShouldWaitForAliasedVar() {
			rh := NewInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_, v, expert)
			av := analysis.GetAliasedVar()
			if avAnalysis := this_.currentScopeHandledVarAnalysesMap.Get(av); avAnalysis != nil {
				rh.HandleAliasedVarCompletion(av, avAnalysis)
			} else {
				this_.varToAliasRetryHandlersMap.Put(av, rh)
			}
		} else {
			if analysis.ShouldInline() {
				analysis.PerformInline()
			}
			this_.currentScopeHandledVarAnalysesMap.Put(v, analysis)
			for _, rh := range this_.varToAliasRetryHandlersMap.RemoveAll(v) {
				rh.HandleAliasedVarCompletion(v, analysis)
			}
		}
	}
}
`

	wave43bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                   "module example.com/w43b_iv\n\ngo 1.24\n",
		"gen/j2gort/j2gort.go":     wave43bJ2gortStubSrc(),
		"gen/rhino/rhino.go":       rhinoSrc,
		"gen/jscomp/inline_var.go": jscompSrc,
	})
	wave43bMustBuild(t, tmpDir)

	allPasses := passes.AllPasses()
	jscompDir := filepath.Join(tmpDir, "gen/jscomp")
	s1, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero sites on Run 1")
	}

	s2, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT idempotent: TotalSites=%d FilesModified=%d byPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	gotBytes, err := os.ReadFile(filepath.Join(jscompDir, "inline_var.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"ivSafeAliases_",
		"activeRefMap_",
		"staticIvLazyNegUnknown_",
		"staticIvLazyNegWellDef_",
		"staticIvLazyNegAssignedOnce_",
		"isSafeToInlineAliasesIv_",
		"rcCol.referenceMap.Get(aliasedVar)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in lowered inline_var.go:\n%s", want, got)
		}
	}

	runtimeTestSrc := `package jscomp

import (
	"reflect"
	"slices"
	"testing"
)

func TestWave43BInlineVariablesRuntimeAndZeroAlloc(t *testing.T) {
	// 1. Verify ReferenceCollection struct size did NOT grow beyond 64 bytes on 64-bit!
	if reflect.TypeFor[*Reference]().Size() == 8 {
		if sz := reflect.TypeFor[ReferenceCollection]().Size(); sz != 64 {
			t.Fatalf("sizeof(ReferenceCollection) = %d, want 64", sz)
		}
	}

	newRC := func(wellDef, assignedOnce bool) *ReferenceCollection {
		r := &Reference{node: 1}
		WellDefByRef0[r] = wellDef
		AssignedOnceByRef0[r] = assignedOnce
		rc := NewReferenceCollection()
		rc.ref0 = r
		return rc
	}

	col := &ReferenceCollector{}

	// Outer variable (lazy-negative unknown, safe to inline aliases: wellDef=true, assignedOnce=true).
	outerSafe := &Var{name: "outerSafe", kind: 3}
	col.referenceMap.Put(outerSafe, newRC(true, true))

	// Outer variable with pre-evaluated wellDef=true, assignedOnce unknown (true).
	outerPartial := &Var{name: "outerPartial", kind: 2, preEvalWellDef: true}
	col.referenceMap.Put(outerPartial, newRC(true, true))

	// Outer variable unsafe for aliases (wellDef=true, assignedOnce=false).
	outerUnsafe := &Var{name: "outerUnsafe", kind: 3}
	col.referenceMap.Put(outerUnsafe, newRC(true, false))

	// Inner scope 1: aliases outerSafe (keeps varToAliasRetryHandlersMap NON-EMPTY!) + 50 negative unaliased vars.
	alias1 := &Var{name: "alias1", kind: 4, alias: outerSafe}
	col.referenceMap.Put(alias1, newRC(true, true))
	alias2 := &Var{name: "alias2", kind: 4, alias: outerSafe}
	col.referenceMap.Put(alias2, newRC(true, true))
	aliasPartial := &Var{name: "aliasPartial", kind: 4, alias: outerPartial}
	col.referenceMap.Put(aliasPartial, newRC(true, true))
	aliasUnsafe := &Var{name: "aliasUnsafe", kind: 4, alias: outerUnsafe}
	col.referenceMap.Put(aliasUnsafe, newRC(true, true))

	inner1Vars := []*Var{alias1, alias2, aliasPartial, aliasUnsafe}
	for i := 0; i < 50; i++ {
		neg := &Var{name: "neg1", kind: 3}
		col.referenceMap.Put(neg, newRC(false, false))
		inner1Vars = append(inner1Vars, neg)
	}

	// Inner scope 2: 50 more negative unaliased vars while varToAliasRetryHandlersMap is STILL non-empty!
	var inner2Vars []*Var
	for i := 0; i < 50; i++ {
		neg := &Var{name: "neg2", kind: 3}
		col.referenceMap.Put(neg, newRC(true, true))
		inner2Vars = append(inner2Vars, neg)
	}

	WellDefEvalCalls.Store(0)
	AssignedOnceEvalCalls.Store(0)

	iv := &InlineVariables{
		scopes: []*Scope{
			{vars: inner1Vars},
			{vars: inner2Vars},
			{vars: []*Var{outerSafe, outerPartial, outerUnsafe}},
		},
		collector: col,
	}
	iv.Process_2(0, 1)

	wantInlined := []string{"alias1", "alias2", "aliasPartial"}
	if !slices.Equal(iv.inlinedNames, wantInlined) {
		t.Fatalf("inlinedNames = %v, want %v", iv.inlinedNames, wantInlined)
	}

	// All 103 CreateLazyNegative calls used static sentinels -> ONLY the 4 WaitForAlias + 3 InlineWellDefined = 7 slab allocs (NOT 110)!
	if got := iv.slab_InlineVariables_InlineVarAnalysis.AllocCalls.Load(); got != 7 {
		t.Fatalf("slab_InlineVariables_InlineVarAnalysis.AllocCalls = %d, want 7 (zero-alloc CreateLazyNegative!)", got)
	}

	// Unconditional per-scope Reset() ran on ALL 3 scopes (even scopes 1 & 2 when varToAliasRetryHandlersMap was non-empty)!
	if got := iv.slab_InlineVariables_InlineVarAnalysis.ResetCalls.Load(); got < 3 {
		t.Fatalf("slab_InlineVariables_InlineVarAnalysis.ResetCalls = %d, want >= 3", got)
	}

	// Verify ivSafeAliases_ caching: outerSafe had TWO aliases (alias1 + alias2), yet IsWellDefined() on outerSafe ran ONLY ONCE!
	// Total IsWellDefined calls:
	// - outerPartial preEvalWellDef: 1 (and ReanalyzeAfterAliasedVar skipped IsWellDefined on outerPartial via staticIvLazyNegWellDef_!)
	// - outerSafe (cached across alias1 & alias2): 1
	// - alias1, alias2, aliasPartial IsWellDefinedAssignedOnce: 3
	// - outerUnsafe: 1
	// Total = 6 (NOT 8+, and 0 for the 100 unaliased negative vars!).
	if got := WellDefEvalCalls.Load(); got != 6 {
		t.Fatalf("WellDefEvalCalls = %d, want 6", got)
	}
}
`
	wave43bWriteFiles(t, tmpDir, map[string]string{"gen/jscomp/w43b_iv_test.go": runtimeTestSrc})
	cmd := exec.Command("go", "test", "-race", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race failed: %v\n%s", err, string(out))
	}
}

// TestWave43B_ToleranceWhenCreateLazyNegativeChangedOrAbsent verifies 100% compatibility
// if Java/Patch changes CreateLazyNegative's parameter type (e.g., *ReferenceCollection) or removes it.
func TestWave43B_ToleranceWhenCreateLazyNegativeChangedOrAbsent(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	jscompSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w43b_tol/gen/j2gort"
	pkg_rhino "example.com/w43b_tol/gen/rhino"
)

var _ atomic.Int32

type Var struct{ name string }
type ReferenceCollection struct{ ok bool }
func (rc *ReferenceCollection) IsWellDefined() bool { return rc.ok }
func (rc *ReferenceCollection) IsAssignedOnceInLifetime() bool { return rc.ok }

type InlineVariables struct{}
func (*InlineVariables) Implements_jscomp_CompilerPass() {}
func (this_ *InlineVariables) Process_2(_ pkg_rhino.Node, _ pkg_rhino.Node) {
	b := &InlineVariables_InliningBehavior{outer_: this_}
	_ = InlineVariables_InlineVarAnalysis_CreateLazyNegative(&ReferenceCollection{ok: true})
	_ = b
}

type InlineVariables_InlineVarAnalysis struct {
	ref  *ReferenceCollection
	kind int8
}

func NewInlineVariables_InlineVarAnalysis(kind int8, ref *ReferenceCollection) *InlineVariables_InlineVarAnalysis {
	return &InlineVariables_InlineVarAnalysis{kind: kind, ref: ref}
}

// Changed signature: takes *ReferenceCollection instead of *InlineVariables_VarExpert!
func InlineVariables_InlineVarAnalysis_CreateLazyNegative(ref *ReferenceCollection) *InlineVariables_InlineVarAnalysis {
	return NewInlineVariables_InlineVarAnalysis(int8(2), ref)
}

type InlineVariables_VarExpert struct {
	v             *Var
	referenceInfo *ReferenceCollection
	behavior      *InlineVariables_InliningBehavior
	fixedAnalysis *InlineVariables_InlineVarAnalysis
}

func NewInlineVariables_VarExpert_5(behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, a, b bool) *InlineVariables_VarExpert {
	e := &InlineVariables_VarExpert{}
	InitInlineVariables_VarExpert_5(e, behavior, v, referenceInfo, a, b)
	return e
}

func InitInlineVariables_VarExpert_5(this_ *InlineVariables_VarExpert, behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, _, _ bool) {
	this_.behavior = behavior
	this_.v = v
	this_.referenceInfo = referenceInfo
}

type InlineVariables_InliningBehavior struct {
	outer_                            *InlineVariables
	currentScopeHandledVarAnalysesMap j2gort.CustomOrderedMap[*Var, *InlineVariables_InlineVarAnalysis]
	varToAliasRetryHandlersMap        j2gort.OrderedMap[*Var, *InlineVariables_InliningBehavior_AliasInlineRetryHandler]
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
`
	wave43bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":               "module example.com/w43b_tol\n\ngo 1.24\n",
		"gen/j2gort/j2gort.go": wave43bJ2gortStubSrc(),
		"gen/rhino/rhino.go":   "package rhino\ntype Node uint32\n",
		"gen/jscomp/iv_tol.go": jscompSrc,
	})

	allPasses := passes.AllPasses()
	jscompDir := filepath.Join(tmpDir, "gen/jscomp")
	if _, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{}); err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	s2, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT idempotent: %+v", s2)
	}

	cmd := exec.Command("go", "build", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, string(out))
	}
}

// TestWave43B_ScopeVarClosureModuleFinGraphsAndPatch86Maps verifies Task 2 + Patch 86 integration:
//   - Scope Cache-Line 0 field packing + allocFastScopeGlobal_ / allocFastScopeChild_ + allocOpt_Var fast path.
//   - Patch 86 MemoizedScopeCreator.scopesByScopeRoot (HashMap -> PlainMapCap preserving dynamic capacity!)
//     still emits ReleaseSlabs_() and releases delegate SyntacticScopeCreator slabs in TypedScopeCreator.CreateInitialScope.
//   - AsIface[Ext_java_util_Map](Maps_NewLinkedHashMapWithExpectedSize(dynamicCap)) lowers to NewOrderedMapCap(dynamicCap).
//   - ClosureRewriteModule_ScriptDescription offset-0 rootNode + PooledSlab + Process_2 defer cleanup.
//   - CoalesceVariableNames.colorings.Pop() graph ReleaseSlabs_() + Patch 86 reusable LiveVariablesAnalysis /
//     MaybeReachingVariableUse ResetForScope & ClearScopeState slab releases.
//   - AnalyzePrototypeProperties + CrossChunkMethodMotion graph/map/scopeCreator ReleaseSlabs_().
//   - LinkedDirectedGraph direct zeroed-slab node/edge init + pooled inEdges/outEdges slices.
func TestWave43B_ScopeVarClosureModuleFinGraphsAndPatch86Maps(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	rhinoSrc := `package rhino

type Node uint32
type StaticSourceFile_SourceKind uint8
type StaticScope interface{}
type StaticSlot interface{}
type StaticRef interface{}

func (n Node) IsRoot() bool   { return n == 1 }
func (n Node) IsScript() bool { return n == 2 }
func (n Node) TokenSwitchOrd_() int32 {
	if n >= 100 && n < 200 {
		return int32(32) // NAME
	}
	return int32(999)
}
`

	graphSrc := `package graph

import (
	"sync/atomic"

	"example.com/w43b_all/gen/j2gort"
)

type Graph[N comparable, E any] interface {
	Implements_graph_Graph()
}

type DiGraph[N comparable, E any] interface {
	Graph[N, E]
	IsJavaNil_() bool
}

type DiGraph_DiGraphEdge[N comparable, E any] interface {
	IsJavaNil_() bool
}

func ToIface_DiGraph_DiGraphEdge[N comparable, E any](e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) DiGraph_DiGraphEdge[N, E] {
	if e == nil {
		return nil
	}
	return e
}

type LinkedDirectedGraph[N comparable, E any] interface {
	DiGraph[N, E]
	Connect(N, E, N)
	GetNodeCount() int32
}

type LinkedDirectedGraphBase[N comparable, E any] struct {
	nodes              j2gort.OrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]
	freeNodes          []*LinkedDirectedGraph_LinkedDiGraphNode[N, E]
	freeEdges          []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]
	useNodeAnnotations bool
	useEdgeAnnotations bool
}

func (g *LinkedDirectedGraphBase[N, E]) IsJavaNil_() bool       { return g == nil }
func (*LinkedDirectedGraphBase[N, E]) Implements_graph_Graph() {}

func (this_ *LinkedDirectedGraphBase[N, E]) initFields_LinkedDirectedGraph() {
	this_.nodes = j2gort.NewOrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]()
}

type LinkedDirectedGraphImpl[N comparable, E any] struct {
	LinkedDirectedGraphBase[N, E]
}

func LinkedDirectedGraph_CreateWithoutAnnotations[N comparable, E any]() LinkedDirectedGraph[N, E] {
	g := &LinkedDirectedGraphImpl[N, E]{}
	g.initFields_LinkedDirectedGraph()
	return g
}

func (this_ *LinkedDirectedGraphBase[N, E]) GetNodeCount() int32 {
	return this_.nodes.Size()
}

func (this_ *LinkedDirectedGraphBase[N, E]) GetOrCreateNode(v N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	n := this_.nodes.Get(v)
	if n == nil {
		n = NewLinkedDirectedGraph_LinkedDiGraphNode[N, E](v, this_.useNodeAnnotations)
		this_.nodes.Put(v, n)
	}
	return n
}

func (this_ *LinkedDirectedGraphBase[N, E]) Connect(src N, edgeVal E, dst N) {
	sn := this_.GetOrCreateNode(src)
	dn := this_.GetOrCreateNode(dst)
	e := NewLinkedDirectedGraph_LinkedDiGraphEdge[N, E](sn, edgeVal, dn, this_.useEdgeAnnotations)
	sn.AddOutEdge(e)
	dn.AddInEdge(e)
}

type LinkedDirectedGraph_LinkedDiGraphNode[N comparable, E any] struct {
	value         N
	annotation    any
	firstInEdge   *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]
	inEdges       []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]
	firstOutEdge  *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]
	outEdges      []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]
	cachedId      string
	idHash_       atomic.Int32
	priority      int32
	astPosition   int32
	dfaIndex      int32
	useAnnotation bool
	inDfaQueue    bool
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) IsJavaNil_() bool { return this_ == nil }

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) initFields_LinkedDirectedGraph_LinkedDiGraphNode() {
	this_.cachedId = j2gort.NullString
	this_.priority = int32(-1)
	this_.astPosition = int32(-1)
	this_.dfaIndex = int32(-1)
}

func NewLinkedDirectedGraph_LinkedDiGraphNode[N comparable, E any](nodeValue N, useAnnotation bool) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	this_ := &LinkedDirectedGraph_LinkedDiGraphNode[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphNode[N, E](this_, nodeValue, useAnnotation)
	return this_
}

func InitLinkedDirectedGraph_LinkedDiGraphNode[N comparable, E any](this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E], nodeValue N, useAnnotation bool) {
	this_.initFields_LinkedDirectedGraph_LinkedDiGraphNode()
	this_.value = nodeValue
	this_.useAnnotation = useAnnotation
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) ClearForPool() {
	this_.firstInEdge = nil
	this_.inEdges = nil
	this_.firstOutEdge = nil
	this_.outEdges = nil
	var z N
	this_.value = z
	this_.annotation = nil
	this_.priority = int32(-1)
	this_.astPosition = int32(-1)
	this_.dfaIndex = int32(-1)
	this_.inDfaQueue = false
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) InEdgeCount() int32 {
	if this_.inEdges != nil {
		return int32(len(this_.inEdges))
	}
	if this_.firstInEdge != nil {
		return 1
	}
	return 0
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) OutEdgeCount() int32 {
	if this_.outEdges != nil {
		return int32(len(this_.outEdges))
	}
	if this_.firstOutEdge != nil {
		return 1
	}
	return 0
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) AddInEdge(edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) {
	if this_.inEdges != nil {
		this_.inEdges = append(this_.inEdges, edge)
	} else if this_.firstInEdge == nil {
		this_.firstInEdge = edge
	} else {
		this_.inEdges = make([]*LinkedDirectedGraph_LinkedDiGraphEdge[N, E], 0, int32(4))
		this_.inEdges = append(this_.inEdges, this_.firstInEdge)
		this_.inEdges = append(this_.inEdges, edge)
		this_.firstInEdge = nil
	}
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) AddOutEdge(edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) {
	if this_.outEdges != nil {
		this_.outEdges = append(this_.outEdges, edge)
	} else if this_.firstOutEdge == nil {
		this_.firstOutEdge = edge
	} else {
		this_.outEdges = make([]*LinkedDirectedGraph_LinkedDiGraphEdge[N, E], 0, int32(4))
		this_.outEdges = append(this_.outEdges, this_.firstOutEdge)
		this_.outEdges = append(this_.outEdges, edge)
		this_.firstOutEdge = nil
	}
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) RemoveInEdge(edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) {
	if this_.inEdges != nil {
		for i, size := int32(0), int32(len(this_.inEdges)); i < size; i++ {
			if j2gort.CheckNotNullPtr(this_.inEdges[i]) == edge {
				this_.inEdges = j2gort.SliceRemoveAt(this_.inEdges, i)
				remaining := int32(len(this_.inEdges))
				if remaining == 1 {
					this_.firstInEdge = this_.inEdges[0]
					this_.inEdges = nil
				} else if remaining == 0 {
					this_.inEdges = nil
				}
				return
			}
		}
	} else if this_.firstInEdge != nil && j2gort.CheckNotNullPtr(this_.firstInEdge) == edge {
		this_.firstInEdge = nil
	}
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) RemoveOutEdge(edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) {
	if this_.outEdges != nil {
		for i, size := int32(0), int32(len(this_.outEdges)); i < size; i++ {
			if j2gort.CheckNotNullPtr(this_.outEdges[i]) == edge {
				this_.outEdges = j2gort.SliceRemoveAt(this_.outEdges, i)
				remaining := int32(len(this_.outEdges))
				if remaining == 1 {
					this_.firstOutEdge = this_.outEdges[0]
					this_.outEdges = nil
				} else if remaining == 0 {
					this_.outEdges = nil
				}
				return
			}
		}
	} else if this_.firstOutEdge != nil && j2gort.CheckNotNullPtr(this_.firstOutEdge) == edge {
		this_.firstOutEdge = nil
	}
}

type LinkedDirectedGraph_LinkedDiGraphEdge[N comparable, E any] struct {
	sourceNode    *LinkedDirectedGraph_LinkedDiGraphNode[N, E]
	value         E
	destNode      *LinkedDirectedGraph_LinkedDiGraphNode[N, E]
	annotation    any
	idHash_       atomic.Int32
	dfaEdgeIndex  int32
	useAnnotation bool
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) IsJavaNil_() bool { return this_ == nil }

func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) initFields_LinkedDirectedGraph_LinkedDiGraphEdge() {
	this_.dfaEdgeIndex = int32(-1)
}

func NewLinkedDirectedGraph_LinkedDiGraphEdge[N comparable, E any](sourceNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, destNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], useAnnotation bool) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	this_ := &LinkedDirectedGraph_LinkedDiGraphEdge[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphEdge[N, E](this_, sourceNode, edgeValue, destNode, useAnnotation)
	return this_
}

func InitLinkedDirectedGraph_LinkedDiGraphEdge[N comparable, E any](this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E], sourceNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, destNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], useAnnotation bool) {
	this_.initFields_LinkedDirectedGraph_LinkedDiGraphEdge()
	this_.value = edgeValue
	this_.sourceNode = sourceNode
	this_.destNode = destNode
	this_.useAnnotation = useAnnotation
}

var UndirectedGraphReleaseCalls atomic.Int64

type DummyUndirectedGraph struct{}

func (g *DummyUndirectedGraph) ReleaseSlabs_() {
	UndirectedGraphReleaseCalls.Add(1)
}

type GraphColoring[N comparable, E any] interface {
	GetGraph() any
}

type GreedyGraphColoring[N comparable, E any] struct {
	g any
}

func NewGreedyColoring[N comparable, E any](g any) GraphColoring[N, E] {
	return &GreedyGraphColoring[N, E]{g: g}
}

func (c *GreedyGraphColoring[N, E]) GetGraph() any { return c.g }
`

	jscompScopeVarSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w43b_all/gen/j2gort"
	pkg_rhino "example.com/w43b_all/gen/rhino"
)

type AbstractCompiler interface {
	GetNumberOfInputs() int32
}

type CompilerInput struct{ name string }

func NodeUtil_CreatesScope(n pkg_rhino.Node) bool {
	return n >= 10 && n < 100
}

type ScopedName interface {
	Implements_jscomp_ScopedName()
	IsJavaNil_() bool
	GetScopedNameBase() *ScopedNameBase
}

type ScopedNameBase struct {
	name           string
	cachedHashCode int32
	scopeRoot      pkg_rhino.Node
}

func (b *ScopedNameBase) GetScopedNameBase() *ScopedNameBase { return b }
func (*ScopedNameBase) Implements_jscomp_ScopedName()        {}
func (b *ScopedNameBase) IsJavaNil_() bool                   { return b == nil }

func InitScopedName(this_ *ScopedNameBase, name string, scopeRoot pkg_rhino.Node) {
	this_.name = j2gort.CheckNotNullString(name)
	this_.scopeRoot = scopeRoot
}

type AbstractScope interface {
	Implements_jscomp_AbstractScope()
	IsJavaNil_() bool
	GetAbstractScopeBase() *AbstractScopeBase
	GetDepth() int32
	GetParent() any
	CheckChildScope(AbstractScope)
	CheckRootScope()
}

func ToIface_AbstractScope[T interface {
	AbstractScope
	comparable
}](v T) AbstractScope {
	var z T
	if v == z {
		return nil
	}
	return v
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
func (this_ *AbstractScopeBase) IsJavaNil_() bool                   { return this_ == nil }
func (b *AbstractScopeBase) GetAbstractScopeBase() *AbstractScopeBase { return b }
func (*AbstractScopeBase) Implements_jscomp_AbstractScope()         {}
func (this_ *AbstractScopeBase) GetRootNode_0() pkg_rhino.Node      { return this_.rootNode }
func (this_ *AbstractScopeBase) GetDepth() int32                    { panic("abstract") }
func (this_ *AbstractScopeBase) GetParent() any                     { panic("abstract") }

func (this_ *AbstractScopeBase) initFields_AbstractScope() {
	this_.implicitVars = nil
}

func InitAbstractScope(this_ *AbstractScopeBase, rootNode pkg_rhino.Node) {
	this_.initFields_AbstractScope()
	if rootNode == 0 {
		panic("NullNode")
	}
	this_.rootNode = rootNode
}

func (this_ *AbstractScopeBase) CheckChildScope(parent AbstractScope) {
	if parent == nil || parent.IsJavaNil_() {
		panic("NPE parent")
	}
	if !NodeUtil_CreatesScope(this_.rootNode) {
		j2gort.CheckArgumentMsg(false, this_.rootNode)
	}
	if this_.rootNode == parent.GetAbstractScopeBase().GetRootNode_0() {
		j2gort.CheckArgumentFmt(false, "same root")
	}
}

func (this_ *AbstractScopeBase) CheckRootScope() {
	if !(NodeUtil_CreatesScope(this_.rootNode) || this_.rootNode.IsScript() || this_.rootNode.IsRoot()) {
		j2gort.CheckArgumentMsg(false, this_.rootNode)
	}
}

type Scope struct {
	AbstractScopeBase
	parent *Scope
	var0   *Var
	var1   *Var
	var2   *Var
	var3   *Var
	var4   *Var
	var5   *Var
	vars   *j2gort.OrderedMap[string, *Var]
	depth  int32
}

func (this_ *Scope) BindSelf_(s any) {
	this_.AbstractScopeBase.BindSelf_(s)
}
func (this_ *Scope) IsJavaNil_() bool         { return this_ == nil }
func (*Scope) Implements_jscomp_AbstractScope() {}
func (this_ *Scope) GetDepth() int32          { return this_.depth }
func (this_ *Scope) GetParent() any           { return j2gort.PtrToAny(this_.parent) }
func (this_ *Scope) initFields_Scope()        {}

func NewScope_2(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_2(this_, parent, rootNode)
	return this_
}

func InitScope_2(this_ *Scope, parent *Scope, rootNode pkg_rhino.Node) {
	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_Scope()
	this_.CheckChildScope(ToIface_AbstractScope(parent))
	this_.parent = parent
	this_.depth = parent.GetDepth() + int32(1)
}

func NewScope_1(rootNode pkg_rhino.Node) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_1(this_, rootNode)
	return this_
}

func InitScope_1(this_ *Scope, rootNode pkg_rhino.Node) {
	InitAbstractScope(&this_.AbstractScopeBase, rootNode)
	this_.initFields_Scope()
	this_.CheckRootScope()
	this_.parent = nil
	this_.depth = int32(0)
}

func Scope_CreateGlobalScope(rootNode pkg_rhino.Node) *Scope {
	return NewScope_1(rootNode)
}

func Scope_CreateChildScope(parent *Scope, rootNode pkg_rhino.Node) *Scope {
	return NewScope_2(parent, rootNode)
}

func (this_ *Scope) Declare(name string, nameNode pkg_rhino.Node, input *CompilerInput) *Var {
	if this_.vars == nil {
		this_.vars = j2gort.NewOrderedMapPtr[string, *Var]()
	}
	v := NewVar(name, nameNode, this_, this_.vars.Size(), input, 0)
	this_.vars.Put(name, v)
	return v
}

type AbstractVar interface {
	Implements_jscomp_AbstractVar()
	IsJavaNil_() bool
	GetAbstractVarBase() *AbstractVarBase
	ScopedName
}

type AbstractVarBase struct {
	ScopedNameBase
	input                         *CompilerInput
	scope                         AbstractScope
	nameNode                      pkg_rhino.Node
	index                         int32
	implicitGoogNamespaceStrength pkg_rhino.StaticSourceFile_SourceKind
}

func (b *AbstractVarBase) GetAbstractVarBase() *AbstractVarBase { return b }
func (*AbstractVarBase) Implements_jscomp_AbstractVar()         {}
func (this_ *AbstractVarBase) IsJavaNil_() bool                 { return this_ == nil }
func (this_ *AbstractVarBase) initFields_AbstractVar()          {}

func InitAbstractVar(this_ *AbstractVarBase, name string, nameNode pkg_rhino.Node, scope AbstractScope, index int32, input *CompilerInput, implicitGoogNamespaceDefinition pkg_rhino.Node) {
	InitScopedName(&this_.ScopedNameBase, j2gort.CheckNotNullString(name), func() pkg_rhino.Node {
		if scope != nil && !scope.IsJavaNil_() {
			return scope.GetAbstractScopeBase().GetRootNode_0()
		}
		return pkg_rhino.Node(0)
	}())
	this_.initFields_AbstractVar()
	if !(index >= int32(-1)) {
		j2gort.CheckArgumentMsg(false, index)
	}
	if implicitGoogNamespaceDefinition != 0 {
		this_.nameNode = 0
		this_.implicitGoogNamespaceStrength = 1
	} else {
		this_.nameNode = nameNode
		this_.implicitGoogNamespaceStrength = 0
	}
	this_.scope = scope
	this_.index = index
	this_.input = input
}

type Var struct {
	AbstractVarBase
}

func (this_ *Var) IsJavaNil_() bool         { return this_ == nil }
func (*Var) Implements_jscomp_AbstractVar() {}
func (this_ *Var) initFields_Var()          {}

func NewVar(name string, nameNode pkg_rhino.Node, scope *Scope, index int32, input *CompilerInput, implicitGoogNamespaceDefinition pkg_rhino.Node) *Var {
	this_ := &Var{}
	InitVar(this_, name, nameNode, scope, index, input, implicitGoogNamespaceDefinition)
	return this_
}

func InitVar(this_ *Var, name string, nameNode pkg_rhino.Node, scope *Scope, index int32, input *CompilerInput, implicitGoogNamespaceDefinition pkg_rhino.Node) {
	InitAbstractVar(&this_.AbstractVarBase, name, nameNode, ToIface_AbstractScope(scope), index, input, implicitGoogNamespaceDefinition)
	this_.initFields_Var()
	if nameNode != 0 {
		switch nameNode.TokenSwitchOrd_() {
		case int32(123), int32(32), int32(119):
		default:
			panic("Invalid name node")
		}
	}
}

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope_2_Node_AbstractScope(pkg_rhino.Node, AbstractScope) any
}

type SyntacticScopeCreator struct {
	compiler AbstractCompiler
}

func (s *SyntacticScopeCreator) IsJavaNil_() bool { return s == nil }

func NewSyntacticScopeCreator_1(compiler AbstractCompiler) *SyntacticScopeCreator {
	return &SyntacticScopeCreator{compiler: compiler}
}

func (this_ *SyntacticScopeCreator) CreateScope_2_Node_AbstractScope(n pkg_rhino.Node, parent AbstractScope) any {
	return j2gort.PtrToAny(this_.CreateScope_2_Node_Scope(n, j2gort.CastPtr[Scope](parent)))
}

func (this_ *SyntacticScopeCreator) CreateScope_2_Node_Scope(n pkg_rhino.Node, parent *Scope) *Scope {
	var scope *Scope
	if parent == nil {
		scope = Scope_CreateGlobalScope(n)
	} else {
		scope = Scope_CreateChildScope(parent, n)
	}
	if n == 10 {
		scope.Declare("x", 101, nil)
	}
	return scope
}

// Patch 86 MemoizedScopeCreator: uses HashMap<Node, AbstractScope> with dynamic expectedScopes!
type MemoizedScopeCreator struct {
	scopesByScopeRoot *j2gort.Ext_java_util_HashMap
	delegate          ScopeCreator
}

func (m *MemoizedScopeCreator) IsJavaNil_() bool { return m == nil }
func (this_ *MemoizedScopeCreator) initFields_MemoizedScopeCreator() {}

func NewMemoizedScopeCreator_2(delegate ScopeCreator, expectedScopes int32) *MemoizedScopeCreator {
	this_ := &MemoizedScopeCreator{}
	InitMemoizedScopeCreator_2(this_, delegate, expectedScopes)
	return this_
}

func InitMemoizedScopeCreator_2(this_ *MemoizedScopeCreator, delegate ScopeCreator, expectedScopes int32) {
	this_.initFields_MemoizedScopeCreator()
	this_.scopesByScopeRoot = j2gort.NewExt_java_util_HashMap_1_int32(expectedScopes)
	this_.delegate = delegate
}

func (this_ *MemoizedScopeCreator) CreateScope_2_Node_AbstractScope(n pkg_rhino.Node, parent AbstractScope) any {
	s := j2gort.AsIface[AbstractScope](this_.scopesByScopeRoot.Get(n))
	if s == nil || s.IsJavaNil_() {
		s = j2gort.AsIface[AbstractScope](this_.delegate.CreateScope_2_Node_AbstractScope(n, parent))
		this_.scopesByScopeRoot.Put_2(n, s)
	}
	return s
}

// DefaultNameGenerator-style AsIface[Ext_java_util_Map](Maps_NewLinkedHashMapWithExpectedSize(dynamicCap))!
type DefaultNameGenerator struct {
	priorityLookupMap j2gort.Ext_java_util_Map
}

func (this_ *DefaultNameGenerator) initFields_DefaultNameGenerator() {}

func NewDefaultNameGenerator(n int32) *DefaultNameGenerator {
	d := &DefaultNameGenerator{}
	InitDefaultNameGenerator(d, n)
	return d
}

func InitDefaultNameGenerator(this_ *DefaultNameGenerator, n int32) {
	this_.initFields_DefaultNameGenerator()
	this_.priorityLookupMap = j2gort.AsIface[j2gort.Ext_java_util_Map](
		j2gort.Ext_com_google_common_collect_Maps_NewLinkedHashMapWithExpectedSize(max(int32(64), n*2)),
	)
	this_.priorityLookupMap.Put_2("a", int32(1))
}

// TypedScopeCreator.CreateInitialScope creating local MemoizedScopeCreator(SyntacticScopeCreator, max(4096, ...))!
type NodeTraversal_Callback interface {
	IsJavaNil_() bool
}

func clinit_NodeTraversal() {}

type NodeTraversal_Builder struct {
	compiler AbstractCompiler
	cb       NodeTraversal_Callback
	sc       ScopeCreator
}

func InitNodeTraversal_Builder(_ *NodeTraversal_Builder) {}
func NodeTraversal_S_Builder() *NodeTraversal_Builder     { return &NodeTraversal_Builder{} }
func (b *NodeTraversal_Builder) SetCompiler(c AbstractCompiler) *NodeTraversal_Builder {
	b.compiler = c
	return b
}
func (b *NodeTraversal_Builder) SetCallback_1_NodeTraversal_Callback(cb NodeTraversal_Callback) *NodeTraversal_Builder {
	b.cb = cb
	return b
}
func (b *NodeTraversal_Builder) SetScopeCreator(sc ScopeCreator) *NodeTraversal_Builder {
	b.sc = sc
	return b
}
func (b *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	if b.sc != nil {
		gs := j2gort.AsIface[AbstractScope](b.sc.CreateScope_2_Node_AbstractScope(1, nil))
		_ = b.sc.CreateScope_2_Node_AbstractScope(10, gs)
	}
}

func NodeTraversal_Traverse(compiler AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	sc := NewSyntacticScopeCreator_1(compiler)
	var b NodeTraversal_Builder
	b.SetCompiler(compiler).SetCallback_1_NodeTraversal_Callback(cb).SetScopeCreator(sc).Traverse(root)
}

type TypedScopeCreator_FirstOrderFunctionAnalyzer struct{}
func (f *TypedScopeCreator_FirstOrderFunctionAnalyzer) IsJavaNil_() bool { return f == nil }
func NewTypedScopeCreator_FirstOrderFunctionAnalyzer(_ *TypedScopeCreator) *TypedScopeCreator_FirstOrderFunctionAnalyzer {
	return &TypedScopeCreator_FirstOrderFunctionAnalyzer{}
}

type TypedScopeCreator struct {
	compiler AbstractCompiler
}

func (this_ *TypedScopeCreator) CreateInitialScope(root pkg_rhino.Node) int32 {
	var val_1 MemoizedScopeCreator
	InitMemoizedScopeCreator_2(&val_1, NewSyntacticScopeCreator_1(this_.compiler), max(int32(4096), this_.compiler.GetNumberOfInputs()*2))
	var scopeCreator *MemoizedScopeCreator = &val_1
	_ = scopeCreator
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(NewTypedScopeCreator_FirstOrderFunctionAnalyzer(this_))
		(&ntb_).SetScopeCreator(scopeCreator)
		(&ntb_).Traverse(root)
	}
	return int32(2)
}
`

	jscompPassesSrc := `package jscomp

import (
	"sync/atomic"

	pkg_graph "example.com/w43b_all/gen/graph"
	"example.com/w43b_all/gen/j2gort"
	pkg_rhino "example.com/w43b_all/gen/rhino"
)

// --- ClosureRewriteModule + ScriptDescription ---
type ClosureRewriteModule_AliasName struct{ newName string }
type ClosureRewriteModule_ExportDefinition struct{ name string }

type ClosureRewriteModule_GlobalRewriteState struct {
	scriptDescriptionsByGoogModuleNamespace j2gort.OrderedMap[string, *ClosureRewriteModule_ScriptDescription]
}

type ClosureRewriteModule_ScriptDescription struct {
	namespaceId             string
	contentsPrefix          string
	topLevelNames           j2gort.OrderedSet[string]
	childScripts            j2gort.Deque[*ClosureRewriteModule_ScriptDescription]
	namesToInlineByAlias    j2gort.OrderedMap[string, *ClosureRewriteModule_AliasName]
	defaultExport           *ClosureRewriteModule_ExportDefinition
	defaultExportLocalName  string
	namedExports            j2gort.OrderedSet[string]
	exportsToInline         j2gort.CustomOrderedMap[*Var, *ClosureRewriteModule_ExportDefinition]
	rootNode                pkg_rhino.Node
	idHash_                 atomic.Int32
	isModule                bool
	declareLegacyNamespace  bool
	willCreateExportsObject bool
	hasCreatedExportObject  bool
}

func (this_ *ClosureRewriteModule_ScriptDescription) initFields_ClosureRewriteModule_ScriptDescription() {
	_ = this_
	this_.namespaceId = j2gort.NullString
	this_.contentsPrefix = j2gort.NullString
	this_.topLevelNames = j2gort.NewOrderedSet[string]()
	this_.childScripts = j2gort.Deque[*ClosureRewriteModule_ScriptDescription]{}
	this_.namesToInlineByAlias = j2gort.NewOrderedMap[string, *ClosureRewriteModule_AliasName]()
	this_.defaultExportLocalName = j2gort.NullString
	this_.namedExports = j2gort.NewOrderedSet[string]()
	this_.exportsToInline = j2gort.NewCustomOrderedMap[*Var, *ClosureRewriteModule_ExportDefinition]()
}

func NewClosureRewriteModule_ScriptDescription() *ClosureRewriteModule_ScriptDescription {
	this_ := &ClosureRewriteModule_ScriptDescription{}
	InitClosureRewriteModule_ScriptDescription(this_)
	return this_
}

func InitClosureRewriteModule_ScriptDescription(this_ *ClosureRewriteModule_ScriptDescription) {
	this_.initFields_ClosureRewriteModule_ScriptDescription()
}

type ClosureRewriteModule struct {
	compiler      AbstractCompiler
	globalScope   *Scope
	scriptStack   j2gort.Deque[*ClosureRewriteModule_ScriptDescription]
	currentScript *ClosureRewriteModule_ScriptDescription
	rewriteState  *ClosureRewriteModule_GlobalRewriteState
	MatchedRoots  int32
}

func (*ClosureRewriteModule) Implements_jscomp_CompilerPass() {}

func (this_ *ClosureRewriteModule) Process_2(_ pkg_rhino.Node, root pkg_rhino.Node) {
	if this_.rewriteState == nil {
		this_.rewriteState = &ClosureRewriteModule_GlobalRewriteState{}
	}
	var descriptions []*ClosureRewriteModule_ScriptDescription
	for i := int32(0); i < 16; i++ {
		sd := NewClosureRewriteModule_ScriptDescription()
		sd.rootNode = root + pkg_rhino.Node(i)
		sd.isModule = (i%2 == 0)
		if sd.namespaceId != j2gort.NullString {
			panic("expected NullString namespaceId")
		}
		this_.currentScript = sd
		this_.scriptStack.Push(sd)
		this_.rewriteState.scriptDescriptionsByGoogModuleNamespace.Put("m", sd)
		descriptions = append(descriptions, sd)
		_ = this_.scriptStack.Pop()
	}
	for i, sd := range descriptions {
		if sd.rootNode == root+pkg_rhino.Node(i) {
			this_.MatchedRoots++
		}
	}
}

// --- AnalyzePrototypeProperties + CrossChunkMethodMotion ---
type JSChunkGraph struct{ chunks int32 }
func (cg *JSChunkGraph) GetChunkCount() int32 { return cg.chunks }

type JSChunk struct{ id int32 }

type AnalyzePrototypeProperties_NameInfo struct {
	name string
	v    *Var
}

type AnalyzePrototypeProperties_ProcessProperties struct {
	outer_ *AnalyzePrototypeProperties
}

func (p *AnalyzePrototypeProperties_ProcessProperties) IsJavaNil_() bool { return p == nil }

type AnalyzePrototypeProperties struct {
	compiler         AbstractCompiler
	chunkGraph       *JSChunkGraph
	symbolGraph      pkg_graph.LinkedDirectedGraph[*AnalyzePrototypeProperties_NameInfo, *JSChunk]
	propertyNameInfo j2gort.OrderedMap[string, *AnalyzePrototypeProperties_NameInfo]
	varNameInfo      j2gort.OrderedMap[string, *AnalyzePrototypeProperties_NameInfo]
	canModifyExterns bool
}

func (this_ *AnalyzePrototypeProperties) initFields_AnalyzePrototypeProperties() {
	this_.symbolGraph = pkg_graph.LinkedDirectedGraph_CreateWithoutAnnotations[*AnalyzePrototypeProperties_NameInfo, *JSChunk]()
	this_.propertyNameInfo = j2gort.NewOrderedMap[string, *AnalyzePrototypeProperties_NameInfo]()
	this_.varNameInfo = j2gort.NewOrderedMap[string, *AnalyzePrototypeProperties_NameInfo]()
}

func NewAnalyzePrototypeProperties(compiler AbstractCompiler, chunkGraph *JSChunkGraph, canModifyExterns bool) *AnalyzePrototypeProperties {
	a := &AnalyzePrototypeProperties{}
	a.initFields_AnalyzePrototypeProperties()
	a.compiler = compiler
	a.chunkGraph = chunkGraph
	a.canModifyExterns = canModifyExterns
	return a
}

func (this_ *AnalyzePrototypeProperties) Process_2(externRoot pkg_rhino.Node, root pkg_rhino.Node) {
	if !this_.canModifyExterns {
		NodeTraversal_Traverse(this_.compiler, externRoot, &AnalyzePrototypeProperties_ProcessProperties{outer_: this_})
	}
	NodeTraversal_Traverse(this_.compiler, root, &AnalyzePrototypeProperties_ProcessProperties{outer_: this_})
	n0 := &AnalyzePrototypeProperties_NameInfo{name: "root"}
	c0 := &JSChunk{id: 0}
	for _, prop := range []string{"p1", "p2", "p3", "p4"} {
		ni := &AnalyzePrototypeProperties_NameInfo{name: prop}
		this_.propertyNameInfo.Put(prop, ni)
		this_.symbolGraph.Connect(n0, c0, ni)
	}
}

func (this_ *AnalyzePrototypeProperties) GetCount() int32 {
	return this_.propertyNameInfo.Size()
}

type CrossChunkMethodMotion struct {
	compiler   AbstractCompiler
	chunkGraph *JSChunkGraph
	analyzer   *AnalyzePrototypeProperties
	Moved      int32
}

func (*CrossChunkMethodMotion) Implements_jscomp_CompilerPass() {}

func NewCrossChunkMethodMotion(compiler AbstractCompiler, chunkGraph *JSChunkGraph) *CrossChunkMethodMotion {
	c := &CrossChunkMethodMotion{}
	InitCrossChunkMethodMotion(c, compiler, chunkGraph)
	return c
}

func InitCrossChunkMethodMotion(this_ *CrossChunkMethodMotion, compiler AbstractCompiler, chunkGraph *JSChunkGraph) {
	this_.compiler = compiler
	this_.chunkGraph = chunkGraph
	this_.analyzer = NewAnalyzePrototypeProperties(compiler, chunkGraph, false)
}

func (this_ *CrossChunkMethodMotion) Process_2(externRoot pkg_rhino.Node, root pkg_rhino.Node) {
	if this_.chunkGraph.GetChunkCount() > 1 {
		this_.analyzer.Process_2(externRoot, root)
		this_.Moved = this_.analyzer.GetCount()
	}
}

// --- Patch 86 Reusable LiveVariablesAnalysis + DeadAssignmentsElimination + CoalesceVariableNames ---
type LiveVariablesAnalysis_LiveVariableLattice struct {
	bits int32
}

func NewLiveVariablesAnalysis_LiveVariableLattice(n int32) *LiveVariablesAnalysis_LiveVariableLattice {
	this_ := &LiveVariablesAnalysis_LiveVariableLattice{}
	InitLiveVariablesAnalysis_LiveVariableLattice(this_, n)
	return this_
}

func InitLiveVariablesAnalysis_LiveVariableLattice(this_ *LiveVariablesAnalysis_LiveVariableLattice, n int32) {
	this_.bits = n
}

type LiveVariablesAnalysis_LiveVariableJoinOp struct {
	outer_ *LiveVariablesAnalysis
	acc    int32
}

func InitLiveVariablesAnalysis_LiveVariableJoinOp(_ *LiveVariablesAnalysis_LiveVariableJoinOp) {}

func NewLiveVariablesAnalysis_LiveVariableJoinOp(outer_ *LiveVariablesAnalysis) *LiveVariablesAnalysis_LiveVariableJoinOp {
	j := &LiveVariablesAnalysis_LiveVariableJoinOp{outer_: outer_}
	InitLiveVariablesAnalysis_LiveVariableJoinOp(j)
	return j
}

type LiveVariablesAnalysis struct {
	emptyLattice *LiveVariablesAnalysis_LiveVariableLattice
	scratchGen   *LiveVariablesAnalysis_LiveVariableLattice
	scratchKill  *LiveVariablesAnalysis_LiveVariableLattice
}

func NewLiveVariablesAnalysis(n int32) *LiveVariablesAnalysis {
	l := &LiveVariablesAnalysis{}
	InitLiveVariablesAnalysis(l, n)
	return l
}

func InitLiveVariablesAnalysis(this_ *LiveVariablesAnalysis, n int32) {
	this_.priv_LiveVariablesAnalysis_initForScope_1(n)
}

func (this_ *LiveVariablesAnalysis) ResetForScope(n int32) {
	this_.priv_LiveVariablesAnalysis_initForScope_1(n)
}

func (this_ *LiveVariablesAnalysis) priv_LiveVariablesAnalysis_initForScope_1(n int32) {
	this_.emptyLattice = NewLiveVariablesAnalysis_LiveVariableLattice(n)
	this_.scratchGen = NewLiveVariablesAnalysis_LiveVariableLattice(n)
	this_.scratchKill = NewLiveVariablesAnalysis_LiveVariableLattice(n)
	_ = NewLiveVariablesAnalysis_LiveVariableJoinOp(this_)
}

func (this_ *LiveVariablesAnalysis) ClearScopeState() {
	// Patch 86 clearScopeState()
}

type DeadAssignmentsElimination struct {
	liveness *LiveVariablesAnalysis
}

func (*DeadAssignmentsElimination) Implements_jscomp_CompilerPass() {}

func (this_ *DeadAssignmentsElimination) Process_2(_ pkg_rhino.Node, _ pkg_rhino.Node) {
	this_.priv_DeadAssignmentsElimination_eliminateDeadAssignments_1(4)
	this_.priv_DeadAssignmentsElimination_eliminateDeadAssignments_1(8)
}

func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_eliminateDeadAssignments_1(n int32) {
	if this_.liveness == nil {
		this_.liveness = NewLiveVariablesAnalysis(n)
	} else {
		this_.liveness.ResetForScope(n)
	}
	defer func() {
		this_.liveness.ClearScopeState()
	}()
}

type CoalesceVariableNames struct {
	colorings        j2gort.Deque[pkg_graph.GraphColoring[*Var, any]]
	liveAnalyses     j2gort.Deque[*LiveVariablesAnalysis]
	freeLiveAnalyses j2gort.Deque[*LiveVariablesAnalysis]
	liveness         *LiveVariablesAnalysis
}

func (*CoalesceVariableNames) Implements_jscomp_CompilerPass() {}

func (this_ *CoalesceVariableNames) EnterScopeWithCfg() {
	if !this_.freeLiveAnalyses.IsEmpty() {
		pooled := this_.freeLiveAnalyses.Pop()
		pooled.ResetForScope(6)
		this_.liveness = pooled
	} else {
		this_.liveness = NewLiveVariablesAnalysis(6)
	}
	this_.liveAnalyses.Push(this_.liveness)
	this_.colorings.Push(pkg_graph.NewGreedyColoring[*Var, any](&pkg_graph.DummyUndirectedGraph{}))
}

func (this_ *CoalesceVariableNames) ExitScopeWithCfg() {
	this_.colorings.Pop()
	finished := this_.liveAnalyses.Pop()
	finished.ClearScopeState()
	this_.freeLiveAnalyses.Push(finished)
	if !this_.liveAnalyses.IsEmpty() {
		l := this_.liveAnalyses.Pop()
		this_.liveAnalyses.Push(l)
		this_.liveness = l
	} else {
		this_.liveness = nil
	}
}

func (this_ *CoalesceVariableNames) Process_2(_ pkg_rhino.Node, _ pkg_rhino.Node) {
	defer func() {
		this_.liveness = nil
	}()
	for i := 0; i < 5; i++ {
		this_.EnterScopeWithCfg()
		this_.ExitScopeWithCfg()
	}
}
`

	wave43bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                   "module example.com/w43b_all\n\ngo 1.24\n",
		"gen/j2gort/j2gort.go":     wave43bJ2gortStubSrc(),
		"gen/rhino/rhino.go":       rhinoSrc,
		"gen/graph/graph.go":       graphSrc,
		"gen/jscomp/scope_var.go":  jscompScopeVarSrc,
		"gen/jscomp/fin_passes.go": jscompPassesSrc,
	})
	wave43bMustBuild(t, tmpDir)

	allPasses := passes.AllPasses()
	for _, pkgRel := range []string{"gen/graph", "gen/jscomp"} {
		dir := filepath.Join(tmpDir, pkgRel)
		s1, err := passes.RunOnPackageDir(dir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 1 on %s failed: %v", pkgRel, err)
		}
		if s1.TotalSites == 0 {
			t.Fatalf("expected non-zero sites on Run 1 for %s", pkgRel)
		}
		s2, err := passes.RunOnPackageDir(dir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 2 on %s failed: %v", pkgRel, err)
		}
		if s2.TotalSites != 0 || s2.FilesModified != 0 {
			t.Fatalf("Run 2 on %s NOT idempotent: TotalSites=%d FilesModified=%d byPass=%+v", pkgRel, s2.TotalSites, s2.FilesModified, s2.ByPass)
		}
	}

	readPkgFile := func(rel string) string {
		b, err := os.ReadFile(filepath.Join(tmpDir, rel))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	graphGot := readPkgFile("gen/graph/graph.go")
	for _, want := range []string{
		"AcquirePooledSlice[*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]](4)",
		"AppendPooledSlice(&this_.inEdges",
		"AppendPooledSlice(&this_.outEdges",
		"ReleasePooledSlice(this_.inEdges)",
		"ReleasePooledSlice(n_.inEdges)",
		"this_.dfaEdgeIndex = int32(-1)",
	} {
		if !strings.Contains(graphGot, want) {
			t.Fatalf("missing %q in lowered graph.go:\n%s", want, graphGot)
		}
	}

	scopeGot := readPkgFile("gen/jscomp/scope_var.go")
	for _, want := range []string{
		"allocFastScopeGlobal_",
		"allocFastScopeChild_",
		"okNode_ :=",
		"NewPlainMapCap[pkg_rhino.Node, AbstractScope](expectedScopes)",
		"func (this_ *MemoizedScopeCreator) ReleaseSlabs_()",
		"defer scopeCreator.ReleaseSlabs_()",
		"NewOrderedMapCap[string, int32](max(int32(64), n*2))",
	} {
		if !strings.Contains(scopeGot, want) {
			t.Fatalf("missing %q in lowered scope_var.go:\n%s", want, scopeGot)
		}
	}

	finGot := readPkgFile("gen/jscomp/fin_passes.go")
	for _, want := range []string{
		"slab_ClosureRewriteModule_ScriptDescription",
		"allocPooled_ClosureRewriteModule_ScriptDescription",
		"poppedColoring_.GetGraph()",
		"func (this_ *AnalyzePrototypeProperties) ReleaseSlabs_()",
		"this_.analyzer.ReleaseSlabs_()",
		"this_.ReleaseSlabs_()",
	} {
		if !strings.Contains(finGot, want) {
			t.Fatalf("missing %q in lowered fin_passes.go:\n%s", want, finGot)
		}
	}

	runtimeTestSrc := `package jscomp

import (
	"reflect"
	"sync"
	"testing"

	pkg_graph "example.com/w43b_all/gen/graph"
	"example.com/w43b_all/gen/j2gort"
)

type mockCompiler struct{ inputs int32 }
func (m mockCompiler) GetNumberOfInputs() int32 { return m.inputs }

func assertPanics(t *testing.T, label string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic in %s", label)
		}
	}()
	fn()
}

func TestWave43BAllOptimizationsRuntime(t *testing.T) {
	// 1. Verify Scope contiguous pointer packing (Wave 44B): AbstractScopeBase, parent, var0..var5, vars, varSlab_
	// are contiguous with 0 internal padding before depth int32 (and parent, var0, var1 fit in Cache Line 0 [0..64)).
	scopeT := reflect.TypeFor[Scope]()
	if reflect.TypeFor[*Scope]().Size() == 8 {
		depthF, okD := scopeT.FieldByName("depth")
		if !okD {
			t.Fatalf("missing field depth on Scope")
		}
		for _, fname := range []string{"AbstractScopeBase", "parent", "var0", "var1"} {
			sf, ok := scopeT.FieldByName(fname)
			if !ok {
				t.Fatalf("missing field %s on Scope", fname)
			}
			if sf.Offset+sf.Type.Size() > 64 {
				t.Fatalf("Scope.%s ends at byte %d, want <= 64 (Cache Line 0!)", fname, sf.Offset+sf.Type.Size())
			}
		}
		for _, fname := range []string{"parent", "var0", "var1", "var2", "var3", "var4", "var5", "vars", "varSlab_"} {
			sf, ok := scopeT.FieldByName(fname)
			if !ok {
				t.Fatalf("missing field %s on Scope", fname)
			}
			if sf.Offset >= depthF.Offset {
				t.Fatalf("Scope.%s (offset %d) must precede Scope.depth (offset %d) to avoid internal padding", fname, sf.Offset, depthF.Offset)
			}
		}
	}

	// 2. Verify ClosureRewriteModule_ScriptDescription.rootNode is at byte offset 0!
	sdT := reflect.TypeFor[ClosureRewriteModule_ScriptDescription]()
	rf, ok := sdT.FieldByName("rootNode")
	if !ok || rf.Offset != 0 {
		t.Fatalf("ClosureRewriteModule_ScriptDescription.rootNode offset = %d (ok=%v), want 0", rf.Offset, ok)
	}

	// 3. Verify SyntacticScopeCreator fast scope + fast Var slab init + fallback error panics.
	ssc := NewSyntacticScopeCreator_1(mockCompiler{inputs: 3000})
	gs := ssc.CreateScope_2_Node_Scope(1, nil)
	cs := ssc.CreateScope_2_Node_Scope(10, gs)
	if cs.GetDepth() != 1 || cs.vars == nil || cs.vars.Get("x") == nil || cs.vars.Get("x").scopeRoot != 10 {
		t.Fatalf("unexpected child scope state: depth=%d vars=%v", cs.GetDepth(), cs.vars)
	}
	assertPanics(t, "invalid child scope node", func() {
		ssc.CreateScope_2_Node_Scope(999, gs)
	})
	assertPanics(t, "invalid var nameNode", func() {
		_ = allocOpt_Var(&ssc.slab_Var, "bad", 999, cs, 0, nil, 0)
	})
	ssc.ReleaseSlabs_()

	// 4. Verify TypedScopeCreator.CreateInitialScope + Patch 86 PlainMap capacity & ReleaseSlabs_() + DefaultNameGenerator cap!
	tsc := &TypedScopeCreator{compiler: mockCompiler{inputs: 3000}}
	if n := tsc.CreateInitialScope(1); n != 2 {
		t.Fatalf("CreateInitialScope scopes = %d, want 2", n)
	}
	msc := NewMemoizedScopeCreator_2(NewSyntacticScopeCreator_1(mockCompiler{inputs: 3000}), 6000)
	if msc.scopesByScopeRoot.CapHint() != 6000 {
		t.Fatalf("MemoizedScopeCreator PlainMap CapHint = %d, want 6000", msc.scopesByScopeRoot.CapHint())
	}
	msc.ReleaseSlabs_()

	dng := NewDefaultNameGenerator(50)
	if dng.priorityLookupMap.CapHint() != 100 {
		t.Fatalf("DefaultNameGenerator OrderedMap CapHint = %d, want 100", dng.priorityLookupMap.CapHint())
	}

	// 5. Verify ClosureRewriteModule ScriptDescription slab + Process_2 defer release.
	crm := &ClosureRewriteModule{compiler: mockCompiler{inputs: 100}}
	crm.Process_2(0, 500)
	if crm.MatchedRoots != 16 {
		t.Fatalf("MatchedRoots = %d, want 16", crm.MatchedRoots)
	}
	if crm.slab_ClosureRewriteModule_ScriptDescription.AllocCalls.Load() != 16 ||
		crm.slab_ClosureRewriteModule_ScriptDescription.ReleaseCalls.Load() != 1 {
		t.Fatalf("ScriptDescription slab alloc=%d release=%d, want 16 / 1",
			crm.slab_ClosureRewriteModule_ScriptDescription.AllocCalls.Load(),
			crm.slab_ClosureRewriteModule_ScriptDescription.ReleaseCalls.Load())
	}

	// 6. Verify CrossChunkMethodMotion + AnalyzePrototypeProperties + LinkedDirectedGraph pooled edge slices!
	j2gort.PooledSliceAcquireCalls.Store(0)
	j2gort.PooledSliceReleaseCalls.Store(0)
	ccmm := NewCrossChunkMethodMotion(mockCompiler{inputs: 100}, &JSChunkGraph{chunks: 4})
	ccmm.Process_2(1, 1)
	if ccmm.Moved != 4 || ccmm.analyzer != nil {
		t.Fatalf("CrossChunkMethodMotion Moved=%d analyzer=%v, want 4 / nil", ccmm.Moved, ccmm.analyzer)
	}
	if acq, rel := j2gort.PooledSliceAcquireCalls.Load(), j2gort.PooledSliceReleaseCalls.Load(); acq == 0 || rel != acq {
		t.Fatalf("LinkedDirectedGraph pooled edge slices acquire=%d release=%d, want >0 and equal!", acq, rel)
	}

	// 7. Verify CoalesceVariableNames.colorings.Pop() graph release + Patch 86 reusable DeadAssignmentsElimination / LiveVariablesAnalysis!
	pkg_graph.UndirectedGraphReleaseCalls.Store(0)
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cvn := &CoalesceVariableNames{}
			cvn.Process_2(0, 1)
			dae := &DeadAssignmentsElimination{}
			dae.Process_2(0, 1)
			if dae.liveness != nil {
				t.Errorf("expected DeadAssignmentsElimination.Process_2 to clear reusable liveness on exit")
			}
		}()
	}
	wg.Wait()
	if got := pkg_graph.UndirectedGraphReleaseCalls.Load(); got != 20 {
		t.Fatalf("UndirectedGraphReleaseCalls = %d, want 20 (5 scopes * 4 goroutines)", got)
	}
}
`
	wave43bWriteFiles(t, tmpDir, map[string]string{"gen/jscomp/w43b_all_test.go": runtimeTestSrc})
	cmd := exec.Command("go", "test", "-race", "-v", "./gen/...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race ./gen/... failed: %v\n%s", err, string(out))
	}
}


// TestWave43B_RealGenGraphTargetsAndIdempotency copies real transpiled/gen/graph into an isolated
// temp module, runs AllPasses() Run 1 + Run 2, verifies all Wave 43B graph optimizations fire
// cleanly, asserts Run 2 produces TotalSites == 0, and builds the package.
func TestWave43B_RealGenGraphTargetsAndIdempotency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real gen/graph check in -short mode")
	}
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
		goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/realw43b_"+pkgName, 1) +
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
		t.Logf("Wave 43B real gen/%s Run 1: total=%d files=%d byPass=%+v", pkgName, s1.TotalSites, s1.FilesModified, s1.ByPass)

		s2, err2 := passes.RunOnPackageDir(dstDir, allPasses, passes.Options{})
		if err2 != nil {
			t.Fatalf("AllPasses() Run 2 on real gen/%s failed: %v", pkgName, err2)
		}
		if s2.TotalSites != 0 || s2.FilesModified != 0 {
			t.Fatalf("Wave 43B real gen/%s Run 2 NOT IDEMPOTENT: total=%d files=%d byPass=%+v", pkgName, s2.TotalSites, s2.FilesModified, s2.ByPass)
		}

		buildCmd := exec.Command("go", "build", "-p", "2", "./"+pkgName)
		buildCmd.Dir = tmpDir
		if out, bErr := buildCmd.CombinedOutput(); bErr != nil {
			t.Fatalf("go build ./%s on lowered real gen/%s failed: %v\n%s", pkgName, pkgName, bErr, string(out))
		}
		return dstDir
	}

	graphCopy := copyAndRunRealPkg("graph")
	ldgBytes, err := os.ReadFile(filepath.Join(graphCopy, "linked_directed_graph.go"))
	if err != nil {
		t.Fatal(err)
	}
	ldg := string(ldgBytes)
	for _, want := range []string{
		"AcquirePooledSlice[*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]](4)",
		"AppendPooledSlice(&this_.inEdges",
		"AppendPooledSlice(&this_.outEdges",
		"ReleasePooledSlice(n_.inEdges)",
		"this_.dfaEdgeIndex = int32(-1)",
	} {
		if !strings.Contains(ldg, want) {
			t.Fatalf("expected lowered real graph/linked_directed_graph.go to contain %q", want)
		}
	}
}

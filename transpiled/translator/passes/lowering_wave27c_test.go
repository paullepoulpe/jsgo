package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave27C_PurgeAllRuntimePools_JSTypeRelease_And_CollectionLowerings(t *testing.T) {
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

	mustWrite("go.mod", "module example.com/w27ctest\n\ngo 1.24\n")

	// 1. Minimal shims/nodestore exporting ReleaseCurrentStoreJSTypes (Wave 27A signature)
	mustWrite("shims/nodestore/nodestore.go", `package nodestore

import "sync"

var (
	evMu   sync.Mutex
	events []string
)

func Record(ev string) {
	evMu.Lock()
	events = append(events, ev)
	evMu.Unlock()
}

func ResetEvents() {
	evMu.Lock()
	events = nil
	evMu.Unlock()
}

func SnapshotEvents() []string {
	evMu.Lock()
	defer evMu.Unlock()
	return append([]string(nil), events...)
}

func ReleaseCurrentStoreJSTypes() {
	Record("nodestore:ReleaseCurrentStoreJSTypes")
}
`)

	// 2. Minimal gen/rhino importing shims/nodestore (mirrors real gen/rhino -> shims/nodestore graph)
	mustWrite("gen/rhino/rhino.go", `package rhino

import "example.com/w27ctest/shims/nodestore"

type Node uint32

func (n Node) IsJavaNil_() bool { return n == 0 }
func (n Node) IsRoot() bool     { return n != 0 }
func TouchStore()               { nodestore.Record("rhino:touch") }
`)

	// 3. Minimal gen/j2gort exporting PurgeAllRuntimePools, TrimStaleSlabPools, Slab, PooledSlab,
	// OrderedMap, ImmutableMap, ImmutableSet, Sets_Difference, WrapTypedArrayList, ArrayList, LinkedHashMap
	mustWrite("gen/j2gort/j2gort.go", `package j2gort

import (
	"fmt"
	"sync"
)

const NullString = "\xff\x00j2go_null"

type NPE struct{}
func (NPE) Error() string { return "NullPointerException" }
func ThrowNPE()           { panic(NPE{}) }

type JavaNilChecker interface {
	IsJavaNil_() bool
}

func isJavaNil(v any) bool {
	if v == nil {
		return true
	}
	if nc, ok := v.(JavaNilChecker); ok {
		return nc.IsJavaNil_()
	}
	return false
}

func CheckNotNull[T any](v T) T {
	if isJavaNil(v) {
		ThrowNPE()
	}
	return v
}
func CheckNotNullPtr[T any](p *T) *T {
	if p == nil {
		ThrowNPE()
	}
	return p
}
func CheckNotNullString(s string) string {
	if s == NullString {
		ThrowNPE()
	}
	return s
}

func BoxString(s string) any {
	if s == NullString {
		return nil
	}
	return s
}
func CastString(v any) string {
	if v == nil {
		return NullString
	}
	return v.(string)
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
	if isJavaNil(v) {
		var z T
		return z
	}
	return v.(T)
}
func BoxNode[N ~uint32](n N) any {
	if n == 0 {
		return nil
	}
	return n
}
func CastNode[N ~uint32](v any) N {
	if v == nil {
		return 0
	}
	return v.(N)
}

var (
	EventMu sync.Mutex
	Events  []string
)

func RecordEvent(ev string) {
	EventMu.Lock()
	Events = append(Events, ev)
	EventMu.Unlock()
}

func ResetEvents() {
	EventMu.Lock()
	Events = nil
	EventMu.Unlock()
}

func SnapshotEvents() []string {
	EventMu.Lock()
	defer EventMu.Unlock()
	return append([]string(nil), Events...)
}

func TrimStaleSlabPools(keepRecent int) {
	RecordEvent(fmt.Sprintf("trim:%d", keepRecent))
}

func PurgeAllRuntimePools() {
	RecordEvent("purge_all_pools")
}

// Slab + PooledSlab
type Slab[T any] struct {
	items [][]T
}

func (s *Slab[T]) Alloc() *T {
	b := make([]T, 1)
	s.items = append(s.items, b)
	return &b[0]
}

type PooledSlab[T any] struct {
	items [][]T
}

func (s *PooledSlab[T]) Alloc() *T {
	b := make([]T, 1)
	s.items = append(s.items, b)
	return &b[0]
}

func (s *PooledSlab[T]) Release() {
	s.items = nil
	RecordEvent("pooled_slab_release")
}

func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	for _, b := range s.items {
		for i := range b {
			fn(&b[i])
		}
	}
	s.items = nil
	RecordEvent("pooled_slab_release_with")
}

// Collection interfaces & helpers
type Ext_java_lang_Iterable interface {
	IterSlice_() []any
}
type Ext_java_util_Collection interface {
	Ext_java_lang_Iterable
	Size() int32
	IsEmpty() bool
	Contains(any) bool
	ContainsAll(Ext_java_util_Collection) bool
}
type Ext_java_util_Set interface {
	Ext_java_util_Collection
}
type Ext_java_util_List interface {
	Ext_java_util_Collection
	Get(int32) any
}
type Ext_java_util_ArrayList interface {
	Ext_java_util_List
	Add_1(any) bool
	Clear()
}
type Ext_java_util_Map interface {
	Size() int32
	IsEmpty() bool
	Clear()
	Get(any) any
	Put_2(any, any) any
	Remove_1(any) any
	ContainsKey(any) bool
	KeySet_0() Ext_java_util_Set
	Values() Ext_java_util_Collection
}

type Ext_com_google_common_collect_ImmutableSet interface {
	Ext_java_util_Set
}

type Ext_com_google_common_collect_ImmutableMap interface {
	Ext_java_util_Map
}

type IterCursor struct {
	items []any
	idx   int
}

func NewIterCursor(it Ext_java_lang_Iterable) IterCursor {
	if isJavaNil(it) {
		ThrowNPE()
	}
	return IterCursor{items: it.IterSlice_()}
}
func (c *IterCursor) HasNext() bool { return c.idx < len(c.items) }
func (c *IterCursor) Next() any {
	v := c.items[c.idx]
	c.idx++
	return v
}
func (c *IterCursor) NextStr() string {
	return CastString(c.Next())
}

type arrayListImpl struct {
	items []any
}

func NewExt_java_util_ArrayList() Ext_java_util_ArrayList {
	return &arrayListImpl{}
}
func NewExt_java_util_ArrayList_1_int32(capHint int32) Ext_java_util_ArrayList {
	return &arrayListImpl{items: make([]any, 0, capHint)}
}
func (l *arrayListImpl) IsJavaNil_() bool     { return l == nil }
func (l *arrayListImpl) IterSlice_() []any    { return l.items }
func (l *arrayListImpl) Size() int32          { return int32(len(l.items)) }
func (l *arrayListImpl) IsEmpty() bool        { return len(l.items) == 0 }
func (l *arrayListImpl) Get(i int32) any      { return l.items[i] }
func (l *arrayListImpl) Clear()               { l.items = l.items[:0] }
func (l *arrayListImpl) Add_1(v any) bool     { l.items = append(l.items, v); return true }
func (l *arrayListImpl) Contains(v any) bool {
	for _, x := range l.items {
		if x == v {
			return true
		}
	}
	return false
}
func (l *arrayListImpl) ContainsAll(c Ext_java_util_Collection) bool {
	for _, x := range c.IterSlice_() {
		if !l.Contains(x) {
			return false
		}
	}
	return true
}

func MakeNonNilSlice[T any](c int32) []T          { return make([]T, 0, c) }
func CloneNonNilSlice[T any](s []T) []T           { out := make([]T, len(s)); copy(out, s); return out }
func CopyCollectionToSlice[T any](c Ext_java_util_Collection) []T {
	raw := c.IterSlice_()
	out := make([]T, len(raw))
	for i, v := range raw {
		if v != nil {
			out[i] = v.(T)
		}
	}
	return out
}

type typedArrayList[T any] struct {
	s []T
}

func WrapTypedArrayList[T any](s []T) Ext_java_util_ArrayList {
	if s == nil {
		return nil
	}
	return &typedArrayList[T]{s: s}
}
func (l *typedArrayList[T]) IsJavaNil_() bool { return l == nil }
func (l *typedArrayList[T]) Size() int32      { return int32(len(l.s)) }
func (l *typedArrayList[T]) IsEmpty() bool    { return len(l.s) == 0 }
func (l *typedArrayList[T]) Get(i int32) any {
	v := l.s[i]
	if isJavaNil(v) {
		return nil
	}
	return v
}
func (l *typedArrayList[T]) IterSlice_() []any {
	out := make([]any, len(l.s))
	for i := range l.s {
		out[i] = l.Get(int32(i))
	}
	return out
}
func (l *typedArrayList[T]) Add_1(any) bool                         { panic("unsupported") }
func (l *typedArrayList[T]) Clear()                                 { l.s = l.s[:0] }
func (l *typedArrayList[T]) Contains(v any) bool                    { return false }
func (l *typedArrayList[T]) ContainsAll(Ext_java_util_Collection) bool { return false }

type setSliceImpl struct {
	items []any
}

func NewStringImmSet(ss ...string) Ext_com_google_common_collect_ImmutableSet {
	items := make([]any, len(ss))
	for i, s := range ss {
		items[i] = s
	}
	return &setSliceImpl{items: items}
}
func (s *setSliceImpl) IsJavaNil_() bool  { return s == nil }
func (s *setSliceImpl) IterSlice_() []any { return s.items }
func (s *setSliceImpl) Size() int32       { return int32(len(s.items)) }
func (s *setSliceImpl) IsEmpty() bool     { return len(s.items) == 0 }
func (s *setSliceImpl) Contains(v any) bool {
	for _, x := range s.items {
		if x == v {
			return true
		}
	}
	return false
}
func (s *setSliceImpl) ContainsAll(c Ext_java_util_Collection) bool {
	for _, x := range c.IterSlice_() {
		if !s.Contains(x) {
			return false
		}
	}
	return true
}

type Ext_com_google_common_collect_Sets_SetView struct {
	items []any
}

func (v *Ext_com_google_common_collect_Sets_SetView) ImmutableCopy() Ext_com_google_common_collect_ImmutableSet {
	return &setSliceImpl{items: append([]any(nil), v.items...)}
}

func Ext_com_google_common_collect_Sets_Difference(a, b Ext_java_util_Set) *Ext_com_google_common_collect_Sets_SetView {
	var diff []any
	for _, x := range a.IterSlice_() {
		if !b.Contains(x) {
			diff = append(diff, x)
		}
	}
	return &Ext_com_google_common_collect_Sets_SetView{items: diff}
}

// LinkedHashMap
type mapEntry struct {
	k      any
	v      any
	active bool
}

type Ext_java_util_LinkedHashMap struct {
	entries []mapEntry
}

func NewExt_java_util_LinkedHashMap() *Ext_java_util_LinkedHashMap {
	return &Ext_java_util_LinkedHashMap{}
}
func (m *Ext_java_util_LinkedHashMap) IsJavaNil_() bool { return m == nil }
func (m *Ext_java_util_LinkedHashMap) Size() int32 {
	var n int32
	for _, e := range m.entries {
		if e.active {
			n++
		}
	}
	return n
}
func (m *Ext_java_util_LinkedHashMap) IsEmpty() bool { return m.Size() == 0 }
func (m *Ext_java_util_LinkedHashMap) Clear()        { m.entries = nil }
func (m *Ext_java_util_LinkedHashMap) Get(k any) any {
	for _, e := range m.entries {
		if e.active && e.k == k {
			return e.v
		}
	}
	return nil
}
func (m *Ext_java_util_LinkedHashMap) ContainsKey(k any) bool {
	for _, e := range m.entries {
		if e.active && e.k == k {
			return true
		}
	}
	return false
}
func (m *Ext_java_util_LinkedHashMap) Put_2(k, v any) any {
	for i := range m.entries {
		if m.entries[i].active && m.entries[i].k == k {
			old := m.entries[i].v
			m.entries[i].v = v
			return old
		}
	}
	m.entries = append(m.entries, mapEntry{k: k, v: v, active: true})
	return nil
}
func (m *Ext_java_util_LinkedHashMap) Remove_1(k any) any {
	for i := range m.entries {
		if m.entries[i].active && m.entries[i].k == k {
			old := m.entries[i].v
			m.entries[i].active = false
			return old
		}
	}
	return nil
}
func (m *Ext_java_util_LinkedHashMap) KeySet_0() Ext_java_util_Set {
	var ks []any
	for _, e := range m.entries {
		if e.active {
			ks = append(ks, e.k)
		}
	}
	return &setSliceImpl{items: ks}
}
func (m *Ext_java_util_LinkedHashMap) Values() Ext_java_util_Collection {
	var vs []any
	for _, e := range m.entries {
		if e.active {
			vs = append(vs, e.v)
		}
	}
	return &setSliceImpl{items: vs}
}

func MapPutStr(m Ext_java_util_Map, k string, v any) any {
	return m.Put_2(BoxString(k), v)
}
func MapGetStr(m Ext_java_util_Map, k string) any {
	return m.Get(BoxString(k))
}
func Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(m Ext_java_util_Map) Ext_com_google_common_collect_ImmutableMap {
	cp := NewExt_java_util_LinkedHashMap()
	ks := m.KeySet_0().IterSlice_()
	for _, k := range ks {
		cp.Put_2(k, m.Get(k))
	}
	return cp
}

// OrderedMap[K, V]
type ordEntry[K comparable, V any] struct {
	Key    K
	Val    V
	active bool
}

type OrderedMap[K comparable, V any] struct {
	entries []ordEntry[K, V]
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}
func NewOrderedMapCap[K comparable, V any](c int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{entries: make([]ordEntry[K, V], 0, c)}
}
func (m *OrderedMap[K, V]) Size() int32 {
	var n int32
	for _, e := range m.entries {
		if e.active {
			n++
		}
	}
	return n
}
func (m *OrderedMap[K, V]) IsEmpty() bool { return m.Size() == 0 }
func (m *OrderedMap[K, V]) Clear()        { m.entries = nil }
func (m *OrderedMap[K, V]) Get(k K) V {
	for _, e := range m.entries {
		if e.active && e.Key == k {
			return e.Val
		}
	}
	var z V
	return z
}
func (m *OrderedMap[K, V]) ContainsKey(k K) bool {
	for _, e := range m.entries {
		if e.active && e.Key == k {
			return true
		}
	}
	return false
}
func (m *OrderedMap[K, V]) Put(k K, v V) V {
	for i := range m.entries {
		if m.entries[i].active && m.entries[i].Key == k {
			old := m.entries[i].Val
			m.entries[i].Val = v
			return old
		}
	}
	m.entries = append(m.entries, ordEntry[K, V]{Key: k, Val: v, active: true})
	var z V
	return z
}
func (m *OrderedMap[K, V]) Remove(k K) V {
	for i := range m.entries {
		if m.entries[i].active && m.entries[i].Key == k {
			old := m.entries[i].Val
			m.entries[i].active = false
			return old
		}
	}
	var z V
	return z
}

type liveKeySetView[K comparable, V any] struct {
	m *OrderedMap[K, V]
}

func (s *liveKeySetView[K, V]) IsJavaNil_() bool { return s == nil || s.m == nil }
func (s *liveKeySetView[K, V]) Size() int32      { return s.m.Size() }
func (s *liveKeySetView[K, V]) IsEmpty() bool    { return s.m.IsEmpty() }
func (s *liveKeySetView[K, V]) IterSlice_() []any {
	var out []any
	for _, e := range s.m.entries {
		if e.active {
			out = append(out, e.Key)
		}
	}
	return out
}
func (s *liveKeySetView[K, V]) Contains(v any) bool {
	if v == nil {
		return false
	}
	k, ok := v.(K)
	return ok && s.m.ContainsKey(k)
}
func (s *liveKeySetView[K, V]) ContainsAll(c Ext_java_util_Collection) bool {
	for _, x := range c.IterSlice_() {
		if !s.Contains(x) {
			return false
		}
	}
	return true
}

func (m *OrderedMap[K, V]) LiveKeySet() Ext_java_util_Set {
	if m == nil {
		ThrowNPE()
	}
	return &liveKeySetView[K, V]{m: m}
}

func (m *OrderedMap[K, V]) LiveValues() Ext_java_util_Collection {
	var vs []any
	for _, e := range m.entries {
		if e.active {
			vs = append(vs, e.Val)
		}
	}
	return &setSliceImpl{items: vs}
}

func (m *OrderedMap[K, V]) ToImmutableMap() Ext_com_google_common_collect_ImmutableMap {
	cp := NewExt_java_util_LinkedHashMap()
	for _, e := range m.entries {
		if e.active {
			cp.Put_2(e.Key, e.Val)
		}
	}
	return cp
}
`)

	// 4. Target package jscomp exercising Task 1, Task 2, and Task 3 (including upgrade of legacy TrimStaleSlabPools!)
	mustWrite("gen/jscomp/phase_pools.go", `package jscomp

import (
	"example.com/w27ctest/gen/j2gort"
	pkg_rhino "example.com/w27ctest/gen/rhino"
)

type Compiler struct {
	hasTypeCheckRun bool
}

func (this_ *Compiler) AfterPass(name string) {
	j2gort.RecordEvent("afterPass:" + name)
}

func (this_ *Compiler) ClearJSTypeRegistry() {
	j2gort.RecordEvent("compiler:ClearJSTypeRegistry")
}

func (this_ *Compiler) HasTypeCheckingRun() bool {
	return this_.hasTypeCheckRun
}

// Task 1.1: Compiler.ParseInputs (including legacy TrimStaleSlabPools(0) upgrade!)
func (this_ *Compiler) ParseInputs() pkg_rhino.Node {
	defer j2gort.TrimStaleSlabPools(0)

	defer func() {
		this_.AfterPass("parseInputs")
	}()
	j2gort.RecordEvent("body:ParseInputs")
	return pkg_rhino.Node(1)
}

// Task 1.2: VariableReferenceCheck.Process_2 (LIFO defer order: ReleaseSlabs_() MUST execute BEFORE PurgeAllRuntimePools()!)
type ReferenceCollector struct{}

func (rc *ReferenceCollector) ReleaseSlabs_() {
	j2gort.RecordEvent("rc:ReleaseSlabs_")
}

func (rc *ReferenceCollector) Process_2(externs, root pkg_rhino.Node) {
	_ = externs
	_ = root
	j2gort.RecordEvent("body:ReferenceCollector.Process_2")
}

type VariableReferenceCheck struct {
	compiler *Compiler
}

func (this_ *VariableReferenceCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer j2gort.TrimStaleSlabPools(2)

	{
		var val_1 ReferenceCollector
		defer func() {
			(&val_1).ReleaseSlabs_()
		}()
		(&val_1).Process_2(externs, root)
	}
}

// Task 1.3: TypeInferencePass.InferAllScopes (sync PurgeAllRuntimePools at entry + defer PurgeAllRuntimePools on exit)
type TypedScope struct{}
type TypeInferencePass struct {
	topScope *TypedScope
}

func (this_ *TypeInferencePass) InferAllScopes(inferenceRoot pkg_rhino.Node) *TypedScope {
	_ = inferenceRoot
	j2gort.PurgeAllRuntimePools()
	defer j2gort.PurgeAllRuntimePools()
	j2gort.RecordEvent("body:InferAllScopes")
	if this_.topScope == nil {
		this_.topScope = &TypedScope{}
	}
	return this_.topScope
}

// Task 1.4: TypedAstSerializer.SerializeRoots (sync PurgeAllRuntimePools at entry)
type TypedAst struct{}
type TypedAstSerializer struct{}

func (this_ *TypedAstSerializer) SerializeRoots(externsRoot pkg_rhino.Node, jsRoot pkg_rhino.Node) *TypedAst {
	_ = externsRoot
	_ = jsRoot
	if !(externsRoot.IsRoot()) || !(jsRoot.IsRoot()) {
		j2gort.ThrowNPE()
	}
	j2gort.RecordEvent("body:SerializeRoots")
	return &TypedAst{}
}
`)

	// Separate file WITHOUT pre-existing nodestore import to verify Task 2 clean nodestore import injection!
	mustWrite("gen/jscomp/convert_types_to_colors.go", `package jscomp

import (
	"example.com/w27ctest/gen/j2gort"
	pkg_rhino "example.com/w27ctest/gen/rhino"
)

type ConvertTypesToColors struct {
	compiler *Compiler
}

func (this_ *ConvertTypesToColors) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.compiler.HasTypeCheckingRun()) {
		j2gort.RecordEvent("body:ConvertTypesToColors:noTypeCheck")
		this_.compiler.ClearJSTypeRegistry()
		return
	}
	j2gort.RecordEvent("body:ConvertTypesToColors:withTypeCheck")
	this_.compiler.ClearJSTypeRegistry()
}
`)

	// Task 3A: RemoveUnusedCode_VarInfo.restRemovables field-to-local move handoff + negative aliasing guard
	mustWrite("gen/jscomp/remove_unused_code.go", `package jscomp

import (
	"example.com/w27ctest/gen/j2gort"
)

type RemoveUnusedCode_Removable interface {
	IsJavaNil_() bool
	PropName() string
}

type simpleRemovable struct {
	name string
}

func (r *simpleRemovable) IsJavaNil_() bool { return r == nil }
func (r *simpleRemovable) PropName() string { return r.name }

type RemoveUnusedCode struct {
	independentRemoved []string
	allRemoved         []string
}

func (this_ *RemoveUnusedCode) considerForIndependentRemoval(r RemoveUnusedCode_Removable) {
	this_.independentRemoved = append(this_.independentRemoved, r.PropName())
}

type RemoveUnusedCode_VarInfo struct {
	outer_              *RemoveUnusedCode
	rem0                RemoveUnusedCode_Removable
	rem1                RemoveUnusedCode_Removable
	restRemovables      j2gort.Ext_java_util_ArrayList
	isEntirelyRemovable bool
}

func NewRemoveUnusedCode_VarInfo(outer_ *RemoveUnusedCode) *RemoveUnusedCode_VarInfo {
	v := &RemoveUnusedCode_VarInfo{outer_: outer_}
	InitRemoveUnusedCode_VarInfo(v)
	return v
}

func InitRemoveUnusedCode_VarInfo(this_ *RemoveUnusedCode_VarInfo) {
	this_.isEntirelyRemovable = true
}

func (this_ *RemoveUnusedCode_VarInfo) AddRemovable(removable RemoveUnusedCode_Removable) {
	if this_.isEntirelyRemovable {
		if this_.rem0 == nil || this_.rem0.IsJavaNil_() {
			this_.rem0 = removable
		} else if this_.rem1 == nil || this_.rem1.IsJavaNil_() {
			this_.rem1 = removable
		} else {
			if this_.restRemovables == nil {
				this_.restRemovables = j2gort.NewExt_java_util_ArrayList_1_int32(int32(4))
			}
			this_.restRemovables.Add_1(removable)
		}
	} else {
		this_.outer_.considerForIndependentRemoval(removable)
	}
}

func (this_ *RemoveUnusedCode_VarInfo) SetIsExplicitlyNotRemovable() {
	if this_.isEntirelyRemovable {
		this_.isEntirelyRemovable = false
		if this_.rem0 != nil && !this_.rem0.IsJavaNil_() {
			var r0 RemoveUnusedCode_Removable = this_.rem0
			_ = r0
			var r1 RemoveUnusedCode_Removable = this_.rem1
			_ = r1
			var rest j2gort.Ext_java_util_ArrayList = this_.restRemovables
			_ = rest
			this_.rem0 = nil
			this_.rem1 = nil
			this_.restRemovables = nil
			this_.outer_.considerForIndependentRemoval(r0)
			if r1 != nil && !r1.IsJavaNil_() {
				this_.outer_.considerForIndependentRemoval(r1)
				if rest != nil {
					var i int32 = int32(0)
					var n int32 = rest.Size()
					for ; i < n; i = i + 1 {
						this_.outer_.considerForIndependentRemoval(j2gort.AsIface[RemoveUnusedCode_Removable](rest.Get(i)))
					}
				}
			}
		}
	}
}

func (this_ *RemoveUnusedCode_VarInfo) RemoveAllRemovables() {
	if this_.rem0 != nil && !this_.rem0.IsJavaNil_() {
		this_.outer_.allRemoved = append(this_.outer_.allRemoved, this_.rem0.PropName())
		this_.rem0 = nil
		if this_.rem1 != nil && !this_.rem1.IsJavaNil_() {
			this_.outer_.allRemoved = append(this_.outer_.allRemoved, this_.rem1.PropName())
			this_.rem1 = nil
			if this_.restRemovables != nil {
				var i int32 = int32(0)
				var n int32 = this_.restRemovables.Size()
				for ; i < n; i = i + 1 {
					r := j2gort.AsIface[RemoveUnusedCode_Removable](this_.restRemovables.Get(i))
					this_.outer_.allRemoved = append(this_.outer_.allRemoved, r.PropName())
				}
				this_.restRemovables = nil
			}
		}
	}
}

// Negative soundness guard: field-to-local alias WITHOUT immediate field = nil reset MUST NOT lower!
type NegativeAliasedFieldHolder struct {
	items j2gort.Ext_java_util_ArrayList
}

func (h *NegativeAliasedFieldHolder) Exercise() int32 {
	if h.items == nil {
		h.items = j2gort.NewExt_java_util_ArrayList_1_int32(int32(4))
	}
	h.items.Add_1(j2gort.BoxString("a"))
	var alias j2gort.Ext_java_util_ArrayList = h.items
	h.items.Add_1(j2gort.BoxString("b"))
	return alias.Size()
}
`)

	// Task 3B: ModuleMapCreator.resolve local OrderedMap LiveKeySet + negative reassigned local guard
	mustWrite("gen/jscomp/module_map_creator.go", `package jscomp

import (
	"example.com/w27ctest/gen/j2gort"
)

type Module struct {
	name       string
	namespaces []string
}

func (m *Module) IsJavaNil_() bool { return m == nil }

type UnresolvedModule struct {
	mod *Module
}

func (u *UnresolvedModule) IsJavaNil_() bool { return u == nil }
func (u *UnresolvedModule) Resolve() *Module { return u.mod }

type ModuleMap struct {
	byPath j2gort.Ext_com_google_common_collect_ImmutableMap
	byNs   j2gort.Ext_com_google_common_collect_ImmutableMap
}

type ModuleMapCreator struct {
	unresolvedModules                   j2gort.OrderedMap[string, *UnresolvedModule]
	unresolvedModulesByClosureNamespace j2gort.OrderedMap[string, *UnresolvedModule]
}

func NewModuleMapCreator() *ModuleMapCreator {
	return &ModuleMapCreator{
		unresolvedModules:                   j2gort.NewOrderedMap[string, *UnresolvedModule](),
		unresolvedModulesByClosureNamespace: j2gort.NewOrderedMap[string, *UnresolvedModule](),
	}
}

func (this_ *ModuleMapCreator) Resolve() *ModuleMap {
	var resolvedModules j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = resolvedModules
	var resolvedClosureModules j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = resolvedClosureModules
	for doFirst := true; doFirst || (!(j2gort.AsIface[j2gort.Ext_java_util_Set](resolvedModules.KeySet_0()).ContainsAll(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.unresolvedModules.LiveKeySet())))); doFirst = false {
		var toResolve j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_Sets_Difference(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.unresolvedModules.LiveKeySet()), j2gort.AsIface[j2gort.Ext_java_util_Set](resolvedModules.KeySet_0())).ImmutableCopy()
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(toResolve))
		for cur.HasNext() {
			key := cur.NextStr()
			resolved := this_.unresolvedModules.Get(key).Resolve()
			j2gort.MapPutStr(resolvedModules, key, j2gort.PtrToAny(resolved))
			for _, ns := range resolved.namespaces {
				j2gort.MapPutStr(resolvedClosureModules, ns, j2gort.PtrToAny(resolved))
			}
		}
	}
	for doFirst2 := true; doFirst2 || (!(j2gort.AsIface[j2gort.Ext_java_util_Set](resolvedClosureModules.KeySet_0()).ContainsAll(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.unresolvedModulesByClosureNamespace.LiveKeySet())))); doFirst2 = false {
		var toResolve2 j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_Sets_Difference(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.unresolvedModulesByClosureNamespace.LiveKeySet()), j2gort.AsIface[j2gort.Ext_java_util_Set](resolvedClosureModules.KeySet_0())).ImmutableCopy()
		cur2 := j2gort.NewIterCursor(j2gort.CheckNotNull(toResolve2))
		for cur2.HasNext() {
			ns := cur2.NextStr()
			j2gort.MapPutStr(resolvedClosureModules, ns, j2gort.PtrToAny(this_.unresolvedModulesByClosureNamespace.Get(ns).Resolve()))
		}
	}
	return &ModuleMap{
		byPath: j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(resolvedModules),
		byNs:   j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(resolvedClosureModules),
	}
}

// Negative soundness guard: reassigned local map using KeySet_0() MUST NOT lower to value OrderedMap!
func NegativeReassignedKeySetMap(flag bool) bool {
	var m j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	j2gort.MapPutStr(m, "a", j2gort.BoxString("1"))
	if flag {
		m = j2gort.NewExt_java_util_LinkedHashMap()
	}
	return m.KeySet_0().IsEmpty()
}
`)

	// Task 3C: Timeline + ChangeTracker
	mustWrite("gen/jscomp/timeline.go", `package jscomp

import (
	"sync/atomic"

	"example.com/w27ctest/gen/j2gort"
)

type Timeline struct {
	eventsByTime  j2gort.Ext_java_util_Map
	eventsByValue j2gort.Ext_java_util_Map
	headEvent     *Timeline_Event
	idHash_       atomic.Int32
}

func (this_ *Timeline) IsJavaNil_() bool { return this_ == nil }

func (this_ *Timeline) initFields_Timeline() {
	_ = this_
	this_.eventsByTime = j2gort.NewExt_java_util_LinkedHashMap()
	this_.eventsByValue = j2gort.NewExt_java_util_LinkedHashMap()
	this_.headEvent = NewTimeline_Event(NewTimeline_Time("-beginning-"))
}

func NewTimeline() *Timeline {
	this_ := &Timeline{}
	InitTimeline(this_)
	return this_
}

func InitTimeline(this_ *Timeline) {
	this_.initFields_Timeline()
}

func (this_ *Timeline) Add(value any) {
	_ = value
	this_.priv_Timeline_addEvent_5237(value, this_.eventsByValue)
}

func (this_ *Timeline) Remove(value any) {
	_ = value
	var event *Timeline_Event = j2gort.CastPtr[Timeline_Event](this_.eventsByValue.Remove_1(value))
	_ = event
	if event != nil {
		if event.nextEvent != nil {
			event.nextEvent.previousEvent = event.previousEvent
		} else {
			this_.headEvent = event.previousEvent
		}
		event.previousEvent.nextEvent = event.nextEvent
		event.nextEvent = nil
		event.previousEvent = nil
	}
}

func (this_ *Timeline) Mark(timeName string) {
	_ = timeName
	this_.priv_Timeline_addEvent_5237(NewTimeline_Time(timeName), this_.eventsByTime)
}

func (this_ *Timeline) GetSince(timeName string) j2gort.Ext_java_util_List {
	_ = timeName
	var firstEvent *Timeline_Event = j2gort.CastPtr[Timeline_Event](this_.eventsByTime.Get(NewTimeline_Time(timeName)))
	_ = firstEvent
	if firstEvent == nil {
		return nil
	}
	var values j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = values
	for event := firstEvent; event != nil; event = event.nextEvent {
		if _, isTime := event.value.(*Timeline_Time); !isTime {
			values.Add_1(event.value)
		}
	}
	return values
}

func (this_ *Timeline) priv_Timeline_addEvent_5237(value any, eventsByKey j2gort.Ext_java_util_Map) {
	_ = value
	_ = eventsByKey
	var event *Timeline_Event = j2gort.CastPtr[Timeline_Event](eventsByKey.Get(value))
	_ = event
	if this_.headEvent == event {
		return
	}
	if event != nil {
		event.previousEvent.nextEvent = event.nextEvent
		event.nextEvent.previousEvent = event.previousEvent
		event.nextEvent = nil
	} else {
		event = NewTimeline_Event(value)
		eventsByKey.Put_2(value, j2gort.PtrToAny(event))
	}
	event.previousEvent = this_.headEvent
	this_.headEvent.nextEvent = event
	this_.headEvent = event
}

type Timeline_Event struct {
	nextEvent     *Timeline_Event
	previousEvent *Timeline_Event
	value         any
	idHash_       atomic.Int32
}

func (this_ *Timeline_Event) IsJavaNil_() bool { return this_ == nil }

func NewTimeline_Event(value any) *Timeline_Event {
	j2gort.CheckNotNull(value)
	return &Timeline_Event{value: value}
}

type Timeline_Time struct {
	name string
}

func (this_ *Timeline_Time) IsJavaNil_() bool { return this_ == nil }

func NewTimeline_Time(name string) *Timeline_Time {
	j2gort.CheckNotNullString(name)
	return &Timeline_Time{name: name}
}
`)

	mustWrite("gen/jscomp/change_tracker.go", `package jscomp

import (
	"example.com/w27ctest/gen/j2gort"
	pkg_rhino "example.com/w27ctest/gen/rhino"
)

type ChangeTracker struct {
	changeTimeline *Timeline
}

func NewChangeTracker() *ChangeTracker {
	return &ChangeTracker{changeTimeline: NewTimeline()}
}

func (this_ *ChangeTracker) ReportFunctionDeleted(n pkg_rhino.Node) {
	this_.changeTimeline.Remove(j2gort.BoxNode(n))
}

func (this_ *ChangeTracker) GetChangedScopeNodesForPass(passName string) j2gort.Ext_java_util_List {
	var changedScopeNodes j2gort.Ext_java_util_List = this_.changeTimeline.GetSince(passName)
	this_.changeTimeline.Mark(passName)
	return changedScopeNodes
}

func (this_ *ChangeTracker) RecordChange(n pkg_rhino.Node) {
	this_.changeTimeline.Add(j2gort.BoxNode(n))
}
`)

	jscompDir := filepath.Join(tmpDir, "gen", "jscomp")
	allPasses := passes.AllPasses()

	// Run 1: Apply all passes
	stats1, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] < 5 {
		t.Fatalf("expected >= 5 slab_alloc_lowering sites (4 stripped pool purges/trims + 1 ReleaseCurrentStoreJSTypes), got %+v", stats1.ByPass)
	}
	if stats1.ByPass["native_slice_lowering"] < 8 {
		t.Fatalf("expected >= 8 native_slice_lowering sites on RemoveUnusedCode_VarInfo.restRemovables, got %+v", stats1.ByPass)
	}
	if stats1.ByPass["native_map_lowering"] < 14 {
		t.Fatalf("expected >= 14 native_map_lowering sites on ModuleMapCreator + Timeline/ChangeTracker, got %+v", stats1.ByPass)
	}

	// Run 2: Verify 100% idempotency (0 sites, 0 files modified!)
	stats2, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected Run 2 idempotency (0 sites, 0 files modified), got %+v", stats2)
	}

	readLowered := func(name string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(jscompDir, name))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	normWS := func(s string) string {
		return strings.Join(strings.Fields(s), " ")
	}

	// Verify Task 1 strips TrimStaleSlabPools + ParseInputs/VariableReferenceCheck PurgeAllRuntimePools while keeping InferAllScopes & SerializeRoots entry purges
	phaseCode := readLowered("phase_pools.go")
	if strings.Contains(phaseCode, "j2gort.TrimStaleSlabPools(") {
		t.Fatalf("expected all mid-compilation TrimStaleSlabPools calls to be stripped:\n%s", phaseCode)
	}
	if got := strings.Count(phaseCode, "j2gort.PurgeAllRuntimePools()"); got != 2 {
		t.Fatalf("expected 2 phase-boundary PurgeAllRuntimePools() calls (InferAllScopes entry + SerializeRoots entry; ParseInputs & VariableReferenceCheck unpurged), got %d:\n%s", got, phaseCode)
	}

	// Verify Task 2 AST shape & clean nodestore import
	cttcCode := readLowered("convert_types_to_colors.go")
	for _, want := range []string{
		`"example.com/w27ctest/shims/nodestore"`,
		"func (this_ *ConvertTypesToColors) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) { _ = externs _ = root defer nodestore.ReleaseCurrentStoreJSTypes()",
	} {
		if !strings.Contains(normWS(cttcCode), want) {
			t.Fatalf("expected convert_types_to_colors.go to contain %q:\n%s", want, cttcCode)
		}
	}

	// Verify Task 3A AST shape + negative guard preservation
	rucCode := readLowered("remove_unused_code.go")
	for _, want := range []string{
		"restRemovables []RemoveUnusedCode_Removable",
		"var rest []RemoveUnusedCode_Removable = this_.restRemovables",
		"items j2gort.Ext_java_util_ArrayList",
		"var alias j2gort.Ext_java_util_ArrayList = h.items",
	} {
		if !strings.Contains(normWS(rucCode), want) {
			t.Fatalf("expected remove_unused_code.go to contain %q:\n%s", want, rucCode)
		}
	}

	// Verify Task 3B AST shape + negative guard preservation
	mmcCode := readLowered("module_map_creator.go")
	for _, want := range []string{
		"var resolvedModules j2gort.OrderedMap[string, *Module] = j2gort.NewOrderedMap[string, *Module]()",
		"var resolvedClosureModules j2gort.OrderedMap[string, *Module] = j2gort.NewOrderedMap[string, *Module]()",
		"resolvedModules.LiveKeySet()",
		"resolvedClosureModules.LiveKeySet()",
		"resolvedModules.ToImmutableMap()",
		"resolvedClosureModules.ToImmutableMap()",
		"var m j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()",
	} {
		if !strings.Contains(normWS(mmcCode), want) {
			t.Fatalf("expected module_map_creator.go to contain %q:\n%s", want, mmcCode)
		}
	}

	// Verify Task 3C AST shape in timeline.go + change_tracker.go
	tlCode := readLowered("timeline.go")
	ctCode := readLowered("change_tracker.go")
	for _, want := range []string{
		"eventsByTime j2gort.OrderedMap[string, *Timeline_Event]",
		"eventsByNode_ j2gort.OrderedMap[pkg_rhino.Node, *Timeline_Event]",
		"return j2gort.WrapTypedArrayList(values)",
		"func (this_ *Timeline) AddNode_(n pkg_rhino.Node)",
		"func (this_ *Timeline) RemoveNode_(n pkg_rhino.Node)",
	} {
		if !strings.Contains(normWS(tlCode), want) {
			t.Fatalf("expected timeline.go to contain %q:\n%s", want, tlCode)
		}
	}
	for _, want := range []string{
		"this_.changeTimeline.RemoveNode_(n)",
		"this_.changeTimeline.AddNode_(n)",
	} {
		if !strings.Contains(normWS(ctCode), want) {
			t.Fatalf("expected change_tracker.go to contain %q:\n%s", want, ctCode)
		}
	}

	// Write runtime semantic + LIFO defer order + -race test!
	mustWrite("gen/jscomp/wave27c_rt_test.go", `package jscomp

import (
	"reflect"
	"testing"

	"example.com/w27ctest/gen/j2gort"
	pkg_rhino "example.com/w27ctest/gen/rhino"
	"example.com/w27ctest/shims/nodestore"
)

func TestWave27C_LifoDefersAndRuntimeSemantics(t *testing.T) {
	// 1. Compiler.ParseInputs: NO pool purge after ParseInputs (avoids synchronous madvise thrashing)!
	j2gort.ResetEvents()
	comp := &Compiler{hasTypeCheckRun: true}
	_ = comp.ParseInputs()
	wantParse := []string{"body:ParseInputs", "afterPass:parseInputs"}
	if got := j2gort.SnapshotEvents(); !reflect.DeepEqual(got, wantParse) {
		t.Fatalf("ParseInputs event order = %v, want %v", got, wantParse)
	}

	// 2. VariableReferenceCheck.Process_2: inner block defer ReleaseSlabs_() runs WITHOUT mid-compilation pool destruction (reused by InferConsts)!
	j2gort.ResetEvents()
	vrc := &VariableReferenceCheck{compiler: comp}
	vrc.Process_2(pkg_rhino.Node(1), pkg_rhino.Node(2))
	wantVRC := []string{"body:ReferenceCollector.Process_2", "rc:ReleaseSlabs_"}
	if got := j2gort.SnapshotEvents(); !reflect.DeepEqual(got, wantVRC) {
		t.Fatalf("VariableReferenceCheck.Process_2 event order = %v, want %v", got, wantVRC)
	}

	// 3. TypeInferencePass.InferAllScopes: PurgeAllRuntimePools at entry ONLY!
	j2gort.ResetEvents()
	tip := &TypeInferencePass{}
	_ = tip.InferAllScopes(pkg_rhino.Node(1))
	wantTIP := []string{"purge_all_pools", "body:InferAllScopes"}
	if got := j2gort.SnapshotEvents(); !reflect.DeepEqual(got, wantTIP) {
		t.Fatalf("TypeInferencePass.InferAllScopes event order = %v, want %v", got, wantTIP)
	}

	// 4. TypedAstSerializer.SerializeRoots: PurgeAllRuntimePools at entry!
	j2gort.ResetEvents()
	tas := &TypedAstSerializer{}
	_ = tas.SerializeRoots(pkg_rhino.Node(1), pkg_rhino.Node(2))
	wantTAS := []string{"purge_all_pools", "body:SerializeRoots"}
	if got := j2gort.SnapshotEvents(); !reflect.DeepEqual(got, wantTAS) {
		t.Fatalf("TypedAstSerializer.SerializeRoots event order = %v, want %v", got, wantTAS)
	}

	// 5. ConvertTypesToColors.Process_2: nodestore.ReleaseCurrentStoreJSTypes runs on exit on BOTH paths!
	for _, tcRun := range []bool{false, true} {
		j2gort.ResetEvents()
		nodestore.ResetEvents()
		c := &Compiler{hasTypeCheckRun: tcRun}
		cttc := &ConvertTypesToColors{compiler: c}
		cttc.Process_2(pkg_rhino.Node(1), pkg_rhino.Node(2))
		nsEv := nodestore.SnapshotEvents()
		if !reflect.DeepEqual(nsEv, []string{"nodestore:ReleaseCurrentStoreJSTypes"}) {
			t.Fatalf("ConvertTypesToColors(hasTypeCheck=%v) nodestore events = %v", tcRun, nsEv)
		}
		j2Ev := j2gort.SnapshotEvents()
		if len(j2Ev) != 2 || j2Ev[1] != "compiler:ClearJSTypeRegistry" {
			t.Fatalf("ConvertTypesToColors(hasTypeCheck=%v) j2gort events = %v", tcRun, j2Ev)
		}
	}

	// 6. RemoveUnusedCode_VarInfo.restRemovables field-to-local move handoff semantics
	ruc := &RemoveUnusedCode{}
	vi := NewRemoveUnusedCode_VarInfo(ruc)
	for _, nm := range []string{"p0", "p1", "p2", "p3", "p4"} {
		vi.AddRemovable(&simpleRemovable{name: nm})
	}
	vi.SetIsExplicitlyNotRemovable()
	// Adding after SetIsExplicitlyNotRemovable must go directly to independentRemoved and see restRemovables == nil
	vi.AddRemovable(&simpleRemovable{name: "p5"})
	wantIndep := []string{"p0", "p1", "p2", "p3", "p4", "p5"}
	if !reflect.DeepEqual(ruc.independentRemoved, wantIndep) {
		t.Fatalf("independentRemoved = %v, want %v", ruc.independentRemoved, wantIndep)
	}

	vi2 := NewRemoveUnusedCode_VarInfo(ruc)
	for _, nm := range []string{"q0", "q1", "q2", "q3"} {
		vi2.AddRemovable(&simpleRemovable{name: nm})
	}
	vi2.RemoveAllRemovables()
	if !reflect.DeepEqual(ruc.allRemoved, []string{"q0", "q1", "q2", "q3"}) {
		t.Fatalf("allRemoved = %v", ruc.allRemoved)
	}

	// 7. ModuleMapCreator.Resolve local OrderedMap LiveKeySet semantics
	mmc := NewModuleMapCreator()
	mA := &Module{name: "a.js", namespaces: []string{"ns.a1", "ns.a2"}}
	mB := &Module{name: "b.js", namespaces: []string{"ns.b"}}
	mmc.unresolvedModules.Put("a.js", &UnresolvedModule{mod: mA})
	mmc.unresolvedModules.Put("b.js", &UnresolvedModule{mod: mB})
	mmc.unresolvedModulesByClosureNamespace.Put("ns.a1", &UnresolvedModule{mod: mA})
	mmc.unresolvedModulesByClosureNamespace.Put("ns.a2", &UnresolvedModule{mod: mA})
	mmc.unresolvedModulesByClosureNamespace.Put("ns.b", &UnresolvedModule{mod: mB})
	mm := mmc.Resolve()
	if mm.byPath.Size() != 2 || mm.byNs.Size() != 3 {
		t.Fatalf("ModuleMap sizes = (%d, %d), want (2, 3)", mm.byPath.Size(), mm.byNs.Size())
	}

	// 8. Timeline + ChangeTracker node fast path + freelist recycling + generic fallback
	ct := NewChangeTracker()
	if got := ct.GetChangedScopeNodesForPass("peephole"); got != nil {
		t.Fatalf("first GetChangedScopeNodesForPass must return nil, got %v", got)
	}
	ct.RecordChange(pkg_rhino.Node(10))
	ct.RecordChange(pkg_rhino.Node(20))
	ct.RecordChange(pkg_rhino.Node(30))
	ct.RecordChange(pkg_rhino.Node(10)) // move 10 to head!
	ct.ReportFunctionDeleted(pkg_rhino.Node(20)) // remove & recycle 20!
	ct.RecordChange(pkg_rhino.Node(40)) // reuses recycled Timeline_Event!

	changed := ct.GetChangedScopeNodesForPass("peephole")
	if changed == nil || changed.Size() != 3 {
		t.Fatalf("expected 3 changed nodes [30, 10, 40], got %v", changed)
	}
	gotNodes := []pkg_rhino.Node{
		j2gort.CastNode[pkg_rhino.Node](changed.Get(0)),
		j2gort.CastNode[pkg_rhino.Node](changed.Get(1)),
		j2gort.CastNode[pkg_rhino.Node](changed.Get(2)),
	}
	wantNodes := []pkg_rhino.Node{30, 10, 40}
	if !reflect.DeepEqual(gotNodes, wantNodes) {
		t.Fatalf("Timeline node order = %v, want %v", gotNodes, wantNodes)
	}
	if again := ct.GetChangedScopeNodesForPass("peephole"); again == nil || !again.IsEmpty() {
		t.Fatalf("expected non-nil empty list when no changes since mark, got %v", again)
	}

	// Generic non-Node Timeline compatibility
	genTl := NewTimeline()
	genTl.Mark("m1")
	genTl.Add("alpha")
	genTl.Add("beta")
	genTl.Remove("alpha")
	genSince := genTl.GetSince("m1")
	if genSince == nil || genSince.Size() != 1 || genSince.Get(0) != "beta" {
		t.Fatalf("generic Timeline fallback mismatch: %v", genSince)
	}
}
`)

	testCmd := exec.Command("go", "test", "-race", "-v", "./gen/jscomp")
	testCmd.Dir = tmpDir
	if out, err := testCmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race on lowered Wave 27C package failed: %v\n%s", err, string(out))
	}
}

func TestWave27C_RealGenJscompDryRun(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/rhino or gen/j2gort is mid-edit; skipping real gen/jscomp DryRun")
	}

	realJscompDir := filepath.Join(repoRoot, "transpiled", "gen", "jscomp")
	rucSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "remove_unused_code.go"))
	mmcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "module_map_creator.go"))
	tlSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "timeline.go"))

	slicePreLowered := strings.Contains(string(rucSrc), "restRemovables []RemoveUnusedCode_Removable") ||
		strings.Contains(string(rucSrc), "restRemovables                      []RemoveUnusedCode_Removable")
	mapPreLowered := strings.Contains(string(mmcSrc), "resolvedModules.LiveKeySet()") &&
		strings.Contains(string(tlSrc), "eventsByNode_")

	dryStats, err := passes.RunOnPackageDir(
		realJscompDir,
		[]passes.Pass{
			passes.NewNativeSliceLoweringPass(),
			passes.NewNativeMapLoweringPass(),
			passes.NewSlabAllocLoweringPass(),
		},
		passes.Options{DryRun: true},
	)
	if err != nil {
		t.Fatalf("DryRun Wave 27C passes on real gen/jscomp failed: %v", err)
	}
	if !slicePreLowered && dryStats.ByPass["native_slice_lowering"] < 8 {
		t.Fatalf("expected >= 8 native_slice_lowering sites on un-regenerated gen/jscomp (RemoveUnusedCode_VarInfo.restRemovables), got %+v", dryStats.ByPass)
	}
	if !mapPreLowered && dryStats.ByPass["native_map_lowering"] < 14 {
		t.Fatalf("expected >= 14 native_map_lowering sites on un-regenerated gen/jscomp (ModuleMapCreator + Timeline/ChangeTracker), got %+v", dryStats.ByPass)
	}
}

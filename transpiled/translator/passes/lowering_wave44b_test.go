package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func wave44bWriteFiles(t *testing.T, dir string, files map[string]string) {
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

func wave44bRunAllAndCheckIdempotent(t *testing.T, pkgDir string) *passes.Stats {
	t.Helper()
	allPasses := passes.AllPasses()
	s1, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("AllPasses Run 1 on %s failed: %v", pkgDir, err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero sites on Run 1 for %s", pkgDir)
	}
	s2, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("AllPasses Run 2 on %s failed: %v", pkgDir, err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("AllPasses Run 2 on %s NOT idempotent: TotalSites=%d FilesModified=%d byPass=%+v", pkgDir, s2.TotalSites, s2.FilesModified, s2.ByPass)
	}
	return s1
}

const wave44bJ2gortStub = `package j2gort

import "sync/atomic"

const NullString = "\x00NULL\x00"

type ClinitState struct{}

func EnterClinit(*ClinitState, *uint32) bool { return true }
func ExitClinit(*ClinitState, *uint32)       {}
func RegisterTypePredicate(string, string, func(any) bool) {}
func InstanceOf[T any](v any) bool { _, ok := v.(T); return ok }
func AssignIdentHash(h *atomic.Int32) int32 { h.Store(1); return 1 }
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
func TryCast[T any](v any) (T, bool) {
	if v == nil {
		var z T
		return z, false
	}
	x, ok := v.(T)
	return x, ok
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
func CheckNotNullString(s string) string {
	if s == NullString {
		panic("NPE")
	}
	return s
}
func CheckNotNullSlice[T any](s []T) []T {
	if s == nil {
		panic("NPE")
	}
	return s
}
func CheckArgument(b bool) {
	if !b {
		panic("IAE")
	}
}
func CheckArgumentMsg(b bool, _ any) {
	if !b {
		panic("IAE")
	}
}
func CheckState(b bool) {
	if !b {
		panic("ISE")
	}
}
func CheckStateMsg(b bool, _ any) {
	if !b {
		panic("ISE")
	}
}
func BoxNode(n any) any         { return n }
func NewStringArray(n int32) []string { return make([]string, n) }
func ConstructShimThrowable(msg string) error { return errStr(msg) }

type errStr string
func (e errStr) Error() string { return string(e) }

type Ext_java_lang_Iterable interface{ Iterator() Ext_java_util_Iterator }
type Ext_java_util_Iterator interface {
	HasNext() bool
	Next() any
}
type Ext_java_util_Collection interface {
	Ext_java_lang_Iterable
	Size() int32
	IsEmpty() bool
	Clear()
}
type Ext_java_util_Set interface {
	Ext_java_util_Collection
	Contains(any) bool
	Add_1(any) bool
}
type Ext_java_util_HashSet interface {
	Ext_java_util_Set
}
type hashSetImpl struct{ m map[any]struct{} }
func NewExt_java_util_HashSet() Ext_java_util_HashSet { return &hashSetImpl{m: map[any]struct{}{}} }
func (h *hashSetImpl) Size() int32 { return int32(len(h.m)) }
func (h *hashSetImpl) IsEmpty() bool { return len(h.m) == 0 }
func (h *hashSetImpl) Clear() { clear(h.m) }
func (h *hashSetImpl) Contains(k any) bool { _, ok := h.m[k]; return ok }
func (h *hashSetImpl) Add_1(k any) bool { _, ok := h.m[k]; h.m[k] = struct{}{}; return !ok }
func (h *hashSetImpl) Iterator() Ext_java_util_Iterator { return nil }

type Ext_java_util_Map interface {
	Get(any) any
	Put_2(any, any) any
	Remove_1(any) any
	Clear()
	Size() int32
	IsEmpty() bool
	KeySet_0() Ext_java_util_Set
	ComputeIfAbsent(any, Ext_java_util_function_Function) any
}
type Ext_java_util_HashMap interface{ Ext_java_util_Map }
type Ext_java_util_LinkedHashMap struct{ m map[any]any }

func NewExt_java_util_HashMap() Ext_java_util_HashMap { return &Ext_java_util_LinkedHashMap{m: map[any]any{}} }
func NewExt_java_util_HashMap_1_int32(int32) Ext_java_util_HashMap { return &Ext_java_util_LinkedHashMap{m: map[any]any{}} }
func NewExt_java_util_LinkedHashMap() *Ext_java_util_LinkedHashMap { return &Ext_java_util_LinkedHashMap{m: map[any]any{}} }
func (m *Ext_java_util_LinkedHashMap) Get(k any) any { return m.m[k] }
func (m *Ext_java_util_LinkedHashMap) Put_2(k, v any) any { old := m.m[k]; m.m[k] = v; return old }
func (m *Ext_java_util_LinkedHashMap) Remove_1(k any) any { old := m.m[k]; delete(m.m, k); return old }
func (m *Ext_java_util_LinkedHashMap) Clear() { clear(m.m) }
func (m *Ext_java_util_LinkedHashMap) Size() int32 { return int32(len(m.m)) }
func (m *Ext_java_util_LinkedHashMap) IsEmpty() bool { return len(m.m) == 0 }
func (m *Ext_java_util_LinkedHashMap) KeySet_0() Ext_java_util_Set { return NewExt_java_util_HashSet() }
func (m *Ext_java_util_LinkedHashMap) ComputeIfAbsent(k any, fn Ext_java_util_function_Function) any {
	if v := m.m[k]; v != nil {
		return v
	}
	v := fn.Apply(k)
	m.m[k] = v
	return v
}
func MapGetStr(m Ext_java_util_Map, k string) any { return m.Get(k) }
func MapPutStr(m Ext_java_util_Map, k string, v any) any { return m.Put_2(k, v) }
func MapRemoveStr(m Ext_java_util_Map, k string) any { return m.Remove_1(k) }

type Ext_java_util_function_Function interface{ Apply(any) any }
type fnLambda func(any) any
func (f fnLambda) Apply(x any) any { return f(x) }
func NewLambda_Ext_java_util_function_Function(f func(any) any) Ext_java_util_function_Function { return fnLambda(f) }

type Ext_com_google_common_base_Predicate interface {
	Apply(any) bool
	Test(any) bool
}
type alwaysTruePred struct{}
var alwaysTrueInst Ext_com_google_common_base_Predicate = &alwaysTruePred{}
func (*alwaysTruePred) Apply(any) bool { return true }
func (*alwaysTruePred) Test(any) bool  { return true }
func Ext_com_google_common_base_Predicates_AlwaysTrue() Ext_com_google_common_base_Predicate { return alwaysTrueInst }

type nodePredLambda struct{ fn func(uint32) bool }
func (p *nodePredLambda) Apply(v any) bool { return p.fn(v.(uint32)) }
func (p *nodePredLambda) Test(v any) bool  { return p.fn(v.(uint32)) }
func NewNodeLambda_Ext_com_google_common_base_Predicate(fn func(uint32) bool) Ext_com_google_common_base_Predicate {
	return &nodePredLambda{fn: fn}
}

var PredTestNodeCalls atomic.Int64

func PredTestNode(p Ext_com_google_common_base_Predicate, n uint32) bool {
	PredTestNodeCalls.Add(1)
	if np, ok := p.(*nodePredLambda); ok {
		return np.fn(n)
	}
	return p.Apply(n)
}
func UnwrapNodePred(p Ext_com_google_common_base_Predicate) func(uint32) bool {
	if np, ok := p.(*nodePredLambda); ok && np != nil {
		return np.fn
	}
	return nil
}
func IsAlwaysTruePred(p any) bool {
	_, ok := p.(*alwaysTruePred)
	return ok
}
func RefEquals(a, b any) bool { return a == b }

type Slab[T any] struct{ items []T }
func (s *Slab[T]) Alloc() *T { return new(T) }
func (s *Slab[T]) Reset()    { s.items = nil }
func (s *Slab[T]) Release()  { s.items = nil }

type PooledSlab[T any] struct{ allocs, releases int32 }
func (s *PooledSlab[T]) Alloc() *T { s.allocs++; return new(T) }
func (s *PooledSlab[T]) Reset()    {}
func (s *PooledSlab[T]) Release()  { s.releases++ }
func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) { _ = fn; s.releases++ }

type Deque[T any] struct{ items []T }
func (d *Deque[T]) AddLast(v T)   { d.items = append(d.items, v) }
func (d *Deque[T]) PeekLast() T   { var z T; if len(d.items) == 0 { return z }; return d.items[len(d.items)-1] }
func (d *Deque[T]) RemoveLast() T { v := d.items[len(d.items)-1]; d.items = d.items[:len(d.items)-1]; return v }
func (d *Deque[T]) GetLast() T    { return d.items[len(d.items)-1] }
func (d *Deque[T]) IsEmpty() bool { return len(d.items) == 0 }

type PlainMap[K comparable, V any] struct{ m map[K]V }
func NewPlainMap[K comparable, V any](capHint int) *PlainMap[K, V] {
	return &PlainMap[K, V]{m: make(map[K]V, capHint)}
}
func NewPlainMapVal[K comparable, V any]() PlainMap[K, V] { return PlainMap[K, V]{m: make(map[K]V)} }
func NewPlainMapCap[K comparable, V any](capHint int32) PlainMap[K, V] {
	return PlainMap[K, V]{m: make(map[K]V, capHint)}
}
func NewPlainMapPtr[K comparable, V any]() *PlainMap[K, V] { return &PlainMap[K, V]{m: make(map[K]V)} }
func NewPlainMapPtrCap[K comparable, V any](capHint int32) *PlainMap[K, V] {
	return &PlainMap[K, V]{m: make(map[K]V, capHint)}
}
func (p *PlainMap[K, V]) IsNil() bool { return p == nil }
func (p *PlainMap[K, V]) Size() int32 { if p == nil { return 0 }; return int32(len(p.m)) }
func (p *PlainMap[K, V]) IsEmpty() bool { return p.Size() == 0 }
func (p *PlainMap[K, V]) Clear() { if p != nil { clear(p.m) } }
func (p *PlainMap[K, V]) Get(k K) V { if p == nil { var z V; return z }; return p.m[k] }
func (p *PlainMap[K, V]) Put(k K, v V) V {
	if p.m == nil {
		p.m = make(map[K]V)
	}
	old := p.m[k]
	p.m[k] = v
	return old
}
func (p *PlainMap[K, V]) PutVoid(k K, v V) {
	if p.m == nil {
		p.m = make(map[K]V)
	}
	p.m[k] = v
}
func (p *PlainMap[K, V]) Remove(k K) V { old := p.m[k]; delete(p.m, k); return old }

type OrderedMap[K comparable, V any] struct {
	m        map[K]V
	order    []K
	released bool
}
func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] { return OrderedMap[K, V]{} }
func NewOrderedMapCap[K comparable, V any](c int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V, c)}
}
func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] { return &OrderedMap[K, V]{} }
func NewOrderedMapPtrCap[K comparable, V any](c int32) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{m: make(map[K]V, c)}
}
func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] { return OrderedMap[K, V]{} }
func AcquireOrderedMapCap[K comparable, V any](c int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V, c)}
}
func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] { return &OrderedMap[K, V]{} }
func AcquireOrderedMapPtrCap[K comparable, V any](c int32) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{m: make(map[K]V, c)}
}
func (o *OrderedMap[K, V]) ReleaseToPool() { o.released = true; clear(o.m); o.order = nil }
func (o *OrderedMap[K, V]) Size() int32    { return int32(len(o.m)) }
func (o *OrderedMap[K, V]) IsEmpty() bool  { return len(o.m) == 0 }
func (o *OrderedMap[K, V]) Clear()         { clear(o.m); o.order = o.order[:0] }
func (o *OrderedMap[K, V]) Get(k K) V      { return o.m[k] }
func (o *OrderedMap[K, V]) Put(k K, v V) V {
	if o.m == nil {
		o.m = make(map[K]V)
	}
	old, exists := o.m[k]
	o.m[k] = v
	if !exists {
		o.order = append(o.order, k)
	}
	return old
}
func (o *OrderedMap[K, V]) LiveKeySet() Ext_java_util_Set { return NewExt_java_util_HashSet() }
func OrderedMapComputeIfAbsentFn[K comparable, V comparable](m *OrderedMap[K, V], k K, zero V, fn Ext_java_util_function_Function) V {
	if v := m.Get(k); v != zero {
		return v
	}
	v := fn.Apply(k).(V)
	m.Put(k, v)
	return v
}

type FastCustomKey interface {
	HashCode() int32
	Equals(any) bool
}

type CustomOrderedMap[K FastCustomKey, V any] struct {
	keys []K
	vals []V
}
func NewCustomOrderedMap[K FastCustomKey, V any]() CustomOrderedMap[K, V] { return CustomOrderedMap[K, V]{} }
func NewCustomOrderedMapCap[K FastCustomKey, V any](int32) CustomOrderedMap[K, V] { return CustomOrderedMap[K, V]{} }
func AcquireCustomOrderedMap[K FastCustomKey, V any]() CustomOrderedMap[K, V] { return CustomOrderedMap[K, V]{} }
func AcquireCustomOrderedMapCap[K FastCustomKey, V any](int32) CustomOrderedMap[K, V] { return CustomOrderedMap[K, V]{} }
func (c *CustomOrderedMap[K, V]) ReleaseToPool() { c.keys = nil; c.vals = nil }
func (c *CustomOrderedMap[K, V]) Size() int32    { return int32(len(c.keys)) }
func (c *CustomOrderedMap[K, V]) Clear()         { c.keys = c.keys[:0]; c.vals = c.vals[:0] }
func (c *CustomOrderedMap[K, V]) Get(k K) V {
	for i, ek := range c.keys {
		if ek.Equals(k) {
			return c.vals[i]
		}
	}
	var z V
	return z
}
func (c *CustomOrderedMap[K, V]) Put(k K, v V) V {
	for i, ek := range c.keys {
		if ek.Equals(k) {
			old := c.vals[i]
			c.vals[i] = v
			return old
		}
	}
	c.keys = append(c.keys, k)
	c.vals = append(c.vals, v)
	var z V
	return z
}
func (c *CustomOrderedMap[K, V]) LiveKeySet() Ext_java_util_Set { return NewExt_java_util_HashSet() }
func CustomOrderedMapComputeIfAbsentFn[K FastCustomKey, V comparable](m *CustomOrderedMap[K, V], k K, zero V, fn Ext_java_util_function_Function) V {
	if v := m.Get(k); v != zero {
		return v
	}
	v := fn.Apply(k).(V)
	m.Put(k, v)
	return v
}

type OrderedSet[T comparable] struct{ m map[T]struct{} }
func NewOrderedSet[T comparable]() OrderedSet[T] { return OrderedSet[T]{} }
func (s *OrderedSet[T]) Add(v T) bool {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	_, ok := s.m[v]
	s.m[v] = struct{}{}
	return !ok
}
func (s *OrderedSet[T]) Contains(v T) bool { _, ok := s.m[v]; return ok }
func (s *OrderedSet[T]) Clear()            { clear(s.m) }
func (s *OrderedSet[T]) IsEmpty() bool     { return len(s.m) == 0 }
`

// TestWave44B_PropertyMapPropStorePlainMap verifies Task 2:
// PropertyMap_PropStore.map_ (lazy nil field assigned from local NewExt_java_util_HashMap_1_int32(16))
// lowers to *j2gort.PlainMap[string, *Property] with full runtime parity and round-2 idempotency.
func TestWave44B_PropertyMapPropStorePlainMap(t *testing.T) {
	tmpDir := t.TempDir()
	wave44bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":             "module example.com/w44b_propstore\n\ngo 1.24.0\n",
		"gen/j2gort/rt.go":   wave44bJ2gortStub,
		"gen/rhino/props.go": `package rhino

import (
	"sync/atomic"

	"example.com/w44b_propstore/gen/j2gort"
)

type Property struct {
	name string
}
func (p *Property) GetName() string { return p.name }

type PropertyMap struct {
	store *PropertyMap_PropStore
}

type PropertyMap_PropStore struct {
	names   []string
	items   []*Property
	map_    j2gort.Ext_java_util_HashMap
	size    int32
	idHash_ atomic.Int32
}

func NewPropertyMap_PropStore(p0, p1, p2 *Property) *PropertyMap_PropStore {
	this_ := &PropertyMap_PropStore{}
	InitPropertyMap_PropStore(this_, p0, p1, p2)
	return this_
}

func InitPropertyMap_PropStore(this_ *PropertyMap_PropStore, p0, p1, p2 *Property) {
	nArr := j2gort.NewStringArray(4)
	nArr[0], nArr[1], nArr[2] = p0.GetName(), p1.GetName(), p2.GetName()
	this_.names = nArr
	arr := make([]*Property, 4)
	arr[0], arr[1], arr[2] = p0, p1, p2
	this_.items = arr
	this_.size = 3
}

func (this_ *PropertyMap) GetOwnProperty_1_string(propertyName string) *Property {
	var s *PropertyMap_PropStore = this_.store
	if s == nil {
		return nil
	}
	if s.map_ != nil {
		return j2gort.CastPtr[Property](j2gort.MapGetStr(s.map_, propertyName))
	}
	for i := int32(0); i < s.size; i++ {
		if s.names[i] == propertyName {
			return s.items[i]
		}
	}
	return nil
}

func (this_ *PropertyMap) PutProperty_2_string_Property(name string, newProp *Property) {
	var s *PropertyMap_PropStore = this_.store
	if s.map_ != nil {
		var oldProp *Property = j2gort.CastPtr[Property](j2gort.MapPutStr(s.map_, name, j2gort.PtrToAny(newProp)))
		if oldProp == nil {
			n := s.size
			if n == int32(len(s.items)) {
				s.items = append(s.items, make([]*Property, n)...)
				s.names = append(s.names, make([]string, n)...)
			}
			s.names[n] = name
			s.items[n] = newProp
			s.size = n + 1
		} else {
			for i := int32(0); i < s.size; i++ {
				if s.items[i] == oldProp {
					s.items[i] = newProp
					break
				}
			}
		}
		return
	}
	for i := int32(0); i < s.size; i++ {
		if s.names[i] == name {
			s.items[i] = newProp
			return
		}
	}
	n := s.size
	if n == int32(len(s.items)) {
		s.items = append(s.items, make([]*Property, n)...)
		s.names = append(s.names, make([]string, n)...)
	}
	s.names[n] = name
	s.items[n] = newProp
	newSize := n + 1
	s.size = newSize
	if newSize > int32(8) {
		var m j2gort.Ext_java_util_HashMap = j2gort.NewExt_java_util_HashMap_1_int32(int32(16))
		_ = m
		for i := int32(0); i < newSize; i++ {
			j2gort.MapPutStr(m, s.names[i], j2gort.PtrToAny(s.items[i]))
		}
		r_64 := s
		r_64.map_ = m
	}
}

func (this_ *PropertyMap) RemoveProperty(name string) bool {
	var s *PropertyMap_PropStore = this_.store
	if s.map_ != nil {
		var removed *Property = j2gort.CastPtr[Property](j2gort.MapRemoveStr(s.map_, name))
		if removed == nil {
			return false
		}
		last := s.size - 1
		for i := int32(0); i < s.size; i++ {
			if s.items[i] == removed {
				copy(s.items[i:last], s.items[i+1:s.size])
				copy(s.names[i:last], s.names[i+1:s.size])
				s.items[last] = nil
				s.names[last] = j2gort.NullString
				s.size = last
				if last <= int32(8) {
					r_74 := s
					r_74.map_ = nil
				}
				break
			}
		}
		return true
	}
	return false
}
`,
	})

	rhinoDir := filepath.Join(tmpDir, "gen/rhino")
	wave44bRunAllAndCheckIdempotent(t, rhinoDir)

	gotB, err := os.ReadFile(filepath.Join(rhinoDir, "props.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotB)
	for _, want := range []string{
		"map_    *j2gort.PlainMap[string, *Property]",
		"var m *j2gort.PlainMap[string, *Property] = j2gort.NewPlainMapPtrCap[string, *Property](int32(16))",
		"m.PutVoid(",
		"s.map_.Get(propertyName)",
		"s.map_.Put(name, newProp)",
		"s.map_.Remove(name)",
		"r_64.map_ = m",
		"r_74.map_ = nil",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected lowered props.go to contain %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "Ext_java_util_HashMap") || strings.Contains(got, "MapPutStr") {
		t.Fatalf("expected Ext_java_util_HashMap / MapPutStr eliminated in props.go:\n%s", got)
	}

	wave44bWriteFiles(t, tmpDir, map[string]string{
		"gen/rhino/props_rt_test.go": `package rhino

import (
	"fmt"
	"testing"
)

func TestPropStoreRuntime(t *testing.T) {
	pm := &PropertyMap{store: NewPropertyMap_PropStore(&Property{name: "p0"}, &Property{name: "p1"}, &Property{name: "p2"})}
	for i := 3; i < 12; i++ {
		nm := fmt.Sprintf("p%d", i)
		pm.PutProperty_2_string_Property(nm, &Property{name: nm})
	}
	if pm.store.map_ == nil || pm.store.map_.Size() != 12 {
		t.Fatalf("expected PlainMap of size 12, got %v", pm.store.map_)
	}
	repl := &Property{name: "p5"}
	pm.PutProperty_2_string_Property("p5", repl)
	if pm.GetOwnProperty_1_string("p5") != repl {
		t.Fatalf("expected updated p5")
	}
	for i := 11; i >= 8; i-- {
		if !pm.RemoveProperty(fmt.Sprintf("p%d", i)) {
			t.Fatalf("expected remove p%d to succeed", i)
		}
	}
	if pm.store.map_ != nil {
		t.Fatalf("expected map_ reset to nil at size <= 8")
	}
	if pm.GetOwnProperty_1_string("p5") != repl {
		t.Fatalf("expected p5 still found via linear scan")
	}
}
`,
	})
	cmd := exec.Command("go", "test", "-v", "./gen/rhino")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("runtime test failed: %v\n%s", err, out)
	}
}

// TestWave44B_ScopeLayoutDirectSlabInitCanonicalVarAndDaeFsivNodeUtilHas verifies Tasks 1, 3, 4,
// plus Patch 86 DAE/FSIV Scope release fixes and NodeUtil_Has UnwrapNodePred fast dispatch.
func TestWave44B_ScopeLayoutDirectSlabInitCanonicalVarAndDaeFsivNodeUtilHas(t *testing.T) {
	tmpDir := t.TempDir()
	wave44bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":           "module example.com/w44b_jscomp\n\ngo 1.24.0\n",
		"gen/j2gort/rt.go": wave44bJ2gortStub,
		"gen/rhino/node.go": `package rhino

type Token uint8

const (
	Token_NAME Token = iota + 1
	Token_GETPROP
	Token_FUNCTION
	Token_PARAM_LIST
	Token_OBJECT_PATTERN
	Token_BLOCK
	Token_DO
	Token_WHILE
	Token_FOR
	Token_FOR_OF
	Token_FOR_AWAIT_OF
	Token_FOR_IN
)

func StaticGet_Token_DO() Token           { return Token_DO }
func StaticGet_Token_WHILE() Token        { return Token_WHILE }
func StaticGet_Token_FOR() Token          { return Token_FOR }
func StaticGet_Token_FOR_OF() Token       { return Token_FOR_OF }
func StaticGet_Token_FOR_AWAIT_OF() Token { return Token_FOR_AWAIT_OF }
func StaticGet_Token_FOR_IN() Token       { return Token_FOR_IN }

type Node uint32

type nodeData struct {
	tok        Token
	parent     Node
	firstChild Node
	next       Node
}

var arena = map[Node]nodeData{}

func RegisterNode(id Node, tok Token, parent, firstChild, next Node) {
	arena[id] = nodeData{tok: tok, parent: parent, firstChild: firstChild, next: next}
}

func (n Node) GetToken() Token            { return arena[n].tok }
func (n Node) GetParent() Node            { return arena[n].parent }
func (n Node) HasParent() bool            { return arena[n].parent != 0 }
func (n Node) GetFirstChild() Node        { return arena[n].firstChild }
func (n Node) GetNext() Node              { return arena[n].next }
func (n Node) IsFunction() bool           { return arena[n].tok == Token_FUNCTION }
func (n Node) IsName() bool               { return arena[n].tok == Token_NAME }
func (n Node) IsGetProp() bool            { return arena[n].tok == Token_GETPROP }
func (n Node) IsParamList() bool          { return arena[n].tok == Token_PARAM_LIST }
func (n Node) IsDestructuringPattern() bool { return arena[n].tok == Token_OBJECT_PATTERN }
func (n Node) IsRoot() bool               { return n == 1 }
func (n Node) IsScript() bool             { return n == 2 }
func (n Node) TokenSwitchOrd_() int32     { return int32(arena[n].tok) }
`,
		"gen/jscomp/scope_var_rc_iv.go": `package jscomp

import (
	"sync/atomic"

	"example.com/w44b_jscomp/gen/j2gort"
	pkg_rhino "example.com/w44b_jscomp/gen/rhino"
)

type AbstractCompiler interface{ Implements_jscomp_AbstractCompiler() }
type CompilerInput struct{}

type AbstractScope interface {
	Implements_jscomp_AbstractScope()
	IsJavaNil_() bool
	GetAbstractScopeBase() *AbstractScopeBase
}
type AbstractScopeBase struct {
	self_        AbstractScope
	implicitVars any
	rootNode     pkg_rhino.Node
	idHash_      atomic.Int32
}
func (b *AbstractScopeBase) BindSelf_(s any)                       { b.self_ = s.(AbstractScope) }
func (b *AbstractScopeBase) IsJavaNil_() bool                      { return b == nil }
func (b *AbstractScopeBase) GetAbstractScopeBase() *AbstractScopeBase { return b }
func (b *AbstractScopeBase) GetOwnSlot(name string) any {
	if v, ok := b.implicitVars.(*Var); ok && v != nil && v.name == name {
		return v
	}
	return nil
}

// Intentionally start in the Wave 43B buggy layout (depth int32 between parent and varSlab_).
type Scope struct {
	AbstractScopeBase
	parent   *Scope
	depth    int32
	varSlab_ *j2gort.PooledSlab[Var]
	var0     *Var
	var1     *Var
	var2     *Var
	var3     *Var
	var4     *Var
	var5     *Var
	vars     *j2gort.OrderedMap[string, *Var]
}
func (*Scope) Implements_jscomp_AbstractScope() {}
func (this_ *Scope) IsJavaNil_() bool           { return this_ == nil }
func (this_ *Scope) initFields_Scope()          {}

type ScopedNameBase struct {
	name      string
	scopeRoot pkg_rhino.Node
}
var VarHashCodeCalls atomic.Int64
var VarEqualsCalls atomic.Int64
func (s *ScopedNameBase) HashCode() int32 { VarHashCodeCalls.Add(1); return 1 }
func (s *ScopedNameBase) Equals(o any) bool {
	VarEqualsCalls.Add(1)
	ov, ok := o.(*Var)
	return ok && ov != nil && ov.name == s.name && ov.scopeRoot == s.scopeRoot
}

type AbstractVarBase struct {
	ScopedNameBase
	input    *CompilerInput
	scope    AbstractScope
	nameNode pkg_rhino.Node
	index    int32
}
type Var struct {
	AbstractVarBase
}
func (v *Var) GetName() string  { return v.name }
func (v *Var) IsJavaNil_() bool { return v == nil }

type ScopeCreator interface {
	Implements_jscomp_ScopeCreator()
}
type SyntacticScopeCreator struct {
	compiler   AbstractCompiler
	slab_Scope j2gort.PooledSlab[Scope]
	slab_Var   j2gort.PooledSlab[Var]
}
func (*SyntacticScopeCreator) Implements_jscomp_ScopeCreator() {}
func clinit_SyntacticScopeCreator() {}
func NewSyntacticScopeCreator_1(c AbstractCompiler) *SyntacticScopeCreator {
	s := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_1(s, c)
	return s
}
func InitSyntacticScopeCreator_1(this_ *SyntacticScopeCreator, c AbstractCompiler) {
	this_.compiler = c
}
func (this_ *SyntacticScopeCreator) ReleaseSlabs_() {
	this_.slab_Scope.Release()
	this_.slab_Var.Release()
}

// --- Task 3c: BasicBlock ---
type BasicBlock struct {
	parent     *BasicBlock
	root       pkg_rhino.Node
	idHash_    atomic.Int32
	isFunction bool
	isLoop     bool
}
func (this_ *BasicBlock) BindSelf_(s any)      { _ = s }
func (this_ *BasicBlock) initFields_BasicBlock() { _ = this_ }
func NewBasicBlock(parent *BasicBlock, root pkg_rhino.Node) *BasicBlock {
	this_ := &BasicBlock{}
	InitBasicBlock(this_, parent, root)
	return this_
}
func allocPooled_BasicBlock(slab_ *j2gort.PooledSlab[BasicBlock], parent *BasicBlock, root pkg_rhino.Node) *BasicBlock {
	this_ := slab_.Alloc()
	InitBasicBlock(this_, parent, root)
	return this_
}
func InitBasicBlock(this_ *BasicBlock, parent *BasicBlock, root pkg_rhino.Node) {
	_ = parent
	_ = root
	this_.initFields_BasicBlock()
	this_.parent = parent
	this_.root = root
	this_.isFunction = root.IsFunction()
	if root.HasParent() {
		var pType pkg_rhino.Token = root.GetParent().GetToken()
		_ = pType
		this_.isLoop = ((((((pType == pkg_rhino.StaticGet_Token_DO()) || (pType == pkg_rhino.StaticGet_Token_WHILE())) || (pType == pkg_rhino.StaticGet_Token_FOR())) || (pType == pkg_rhino.StaticGet_Token_FOR_OF())) || (pType == pkg_rhino.StaticGet_Token_FOR_AWAIT_OF())) || (pType == pkg_rhino.StaticGet_Token_FOR_IN()))
	} else {
		this_.isLoop = false
	}
}

// --- Task 4: ReferenceCollector.referenceMap (already-lowered CustomOrderedMap -> OrderedMap + implicit var fallback)
//     and InlineVariables_InliningBehavior.currentScopeHandledVarAnalysesMap (fresh Ext_java_util_Map -> OrderedMap) ---
type ReferenceCollection struct {
	count int32
}
func NewReferenceCollection() *ReferenceCollection { return &ReferenceCollection{} }

type ReferenceCollector struct {
	referenceMap             j2gort.CustomOrderedMap[*Var, *ReferenceCollection]
	slab_BasicBlock          j2gort.PooledSlab[BasicBlock]
	slab_ReferenceCollection j2gort.PooledSlab[ReferenceCollection]
}
func (this_ *ReferenceCollector) initFields_ReferenceCollector() {
	this_.referenceMap = j2gort.AcquireCustomOrderedMap[*Var, *ReferenceCollection]()
}
func (this_ *ReferenceCollector) Process_2(externs, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if this_.referenceMap.Size() == 0 {
		this_.referenceMap.ReleaseToPool()
		this_.referenceMap = j2gort.AcquireCustomOrderedMapCap[*Var, *ReferenceCollection](int32(1024))
	}
}
func (this_ *ReferenceCollector) GetAllSymbols() any {
	return j2gort.AsIface[j2gort.Ext_java_util_Set](this_.referenceMap.LiveKeySet())
}
func (this_ *ReferenceCollector) GetReferences(arg_2 any) any {
	v := j2gort.CastPtr[Var](arg_2)
	_ = v
	return j2gort.PtrToAny(this_.referenceMap.Get(v))
}
func (this_ *ReferenceCollector) AddRef(v *Var) {
	rc := this_.referenceMap.Get(v)
	if rc == nil {
		rc = this_.slab_ReferenceCollection.Alloc()
		this_.referenceMap.Put(v, rc)
	}
	rc.count++
}
func (this_ *ReferenceCollector) ReleaseSlabs_() {
	this_.slab_BasicBlock.Release()
	this_.slab_ReferenceCollection.Release()
	this_.referenceMap.ReleaseToPool()
}

type InlineVariables_InlineVarAnalysis struct{ ok bool }
type InlineVariables_InliningBehavior struct {
	currentScopeHandledVarAnalysesMap j2gort.Ext_java_util_Map
}
func (this_ *InlineVariables_InliningBehavior) initFields_InlineVariables_InliningBehavior() {
	this_.currentScopeHandledVarAnalysesMap = j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.NewExt_java_util_HashMap())
}
func (this_ *InlineVariables_InliningBehavior) PutAndGet(v *Var, a *InlineVariables_InlineVarAnalysis) *InlineVariables_InlineVarAnalysis {
	this_.currentScopeHandledVarAnalysesMap.Clear()
	this_.currentScopeHandledVarAnalysesMap.Put_2(j2gort.PtrToAny(v), j2gort.PtrToAny(a))
	return j2gort.CastPtr[InlineVariables_InlineVarAnalysis](this_.currentScopeHandledVarAnalysesMap.Get(j2gort.PtrToAny(v)))
}
`,
		"gen/jscomp/ruc.go": `package jscomp

import (
	"sync/atomic"

	"example.com/w44b_jscomp/gen/j2gort"
	pkg_rhino "example.com/w44b_jscomp/gen/rhino"
)

type RemoveUnusedCode struct {
	slab_RemoveUnusedCode_IndirectAssign j2gort.PooledSlab[RemoveUnusedCode_IndirectAssign]
	slab_RemoveUnusedCode_VarInfo        j2gort.PooledSlab[RemoveUnusedCode_VarInfo]
	freeBuilder                          *RemoveUnusedCode_RemovableBuilder
}

type RemoveUnusedCode_Continuation struct {
	scope *Scope
	node  pkg_rhino.Node
}

type RemoveUnusedCode_Removable interface {
	Implements_jscomp_RemoveUnusedCode_Removable()
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
func (*RemoveUnusedCode_RemovableBase) Implements_jscomp_RemoveUnusedCode_Removable() {}
func (this_ *RemoveUnusedCode_RemovableBase) initFields_RemoveUnusedCode_Removable() {
	_ = this_
	this_.propertyName = j2gort.NullString
	this_.continuationsAreApplied = false
	this_.isRemoved = false
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
	isPrototypeDotPropertyReference bool
	isThisDotPropertyReference      bool
}
func (this_ *RemoveUnusedCode_RemovableBuilder) priv_RemovableBuilder_finish_4767(r RemoveUnusedCode_Removable) RemoveUnusedCode_Removable {
	if this_.continuation0 != (RemoveUnusedCode_Continuation{}) {
		this_.continuation0 = RemoveUnusedCode_Continuation{}
		if this_.continuation1 != (RemoveUnusedCode_Continuation{}) {
			this_.continuation1 = RemoveUnusedCode_Continuation{}
			this_.restContinuations = nil
		}
	}
	if this_.propertyName != j2gort.NullString {
		this_.propertyName = j2gort.NullString
	}
	this_.isPrototypeDotPropertyReference = false
	this_.isThisDotPropertyReference = false
	this_.outer_.freeBuilder = this_
	return r
}
func (this_ *RemoveUnusedCode_RemovableBuilder) BuildIndirectAssign(root pkg_rhino.Node, targetNode pkg_rhino.Node) *RemoveUnusedCode_IndirectAssign {
	_ = root
	_ = targetNode
	rucOuter_ := this_.outer_
	res_ := rucOuter_.slab_RemoveUnusedCode_IndirectAssign.Alloc()
	res_.self_ = res_
	res_.outer_ = rucOuter_
	InitRemoveUnusedCode_IndirectAssign(res_, this_, root, targetNode)
	this_.priv_RemovableBuilder_finish_4767(nil)
	return res_
}

type RemoveUnusedCode_IndirectAssign struct {
	RemoveUnusedCode_RemovableBase
	outer_ *RemoveUnusedCode
	root   pkg_rhino.Node
}
func (*RemoveUnusedCode_IndirectAssign) Implements_jscomp_RemoveUnusedCode_Removable() {}
func (this_ *RemoveUnusedCode_IndirectAssign) initFields_RemoveUnusedCode_IndirectAssign() {
	_ = this_
}
func allocPooled_RemoveUnusedCode_IndirectAssign(outer_ *RemoveUnusedCode, builder *RemoveUnusedCode_RemovableBuilder, root pkg_rhino.Node, targetNode pkg_rhino.Node) *RemoveUnusedCode_IndirectAssign {
	this_ := outer_.slab_RemoveUnusedCode_IndirectAssign.Alloc()
	this_.self_ = this_
	this_.outer_ = outer_
	InitRemoveUnusedCode_IndirectAssign(this_, builder, root, targetNode)
	return this_
}
func InitRemoveUnusedCode_IndirectAssign(this_ *RemoveUnusedCode_IndirectAssign, builder *RemoveUnusedCode_RemovableBuilder, root pkg_rhino.Node, targetNode pkg_rhino.Node) {
	_ = builder
	_ = root
	_ = targetNode
	this_.RemoveUnusedCode_RemovableBase.outer_ = this_.outer_
	InitRemoveUnusedCode_Removable(&this_.RemoveUnusedCode_RemovableBase, targetNode, builder)
	this_.initFields_RemoveUnusedCode_IndirectAssign()
	var rootParent pkg_rhino.Node = root.GetParent()
	_ = rootParent
	if !(rootParent.IsDestructuringPattern() || rootParent.IsParamList()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(rootParent))
	}
	if !(targetNode.IsName() || targetNode.IsGetProp()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(targetNode))
	}
	this_.root = root
}

type RemoveUnusedCode_VarInfo struct {
	outer_                              *RemoveUnusedCode
	varName                             string
	canonicalUnremovable                bool
	isEntirelyRemovable                 bool
	hasNonLocalOrNonLiteralValue        bool
}
func (this_ *RemoveUnusedCode_VarInfo) initFields_RemoveUnusedCode_VarInfo() {
	_ = this_
	this_.varName = j2gort.NullString
}
func InitRemoveUnusedCode_VarInfo_1(this_ *RemoveUnusedCode_VarInfo, varName string) {
	_ = varName
	this_.initFields_RemoveUnusedCode_VarInfo()
	this_.varName = varName
	this_.canonicalUnremovable = false
	this_.isEntirelyRemovable = true
}
func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_getVarInfo_4732(var_ *Var) *RemoveUnusedCode_VarInfo {
	val_14 := this_.slab_RemoveUnusedCode_VarInfo.Alloc()
	val_14.outer_ = this_
	InitRemoveUnusedCode_VarInfo_1(val_14, var_.GetName())
	var varInfo *RemoveUnusedCode_VarInfo = val_14
	return varInfo
}
`,
		"gen/jscomp/dae_fsiv_nodeutil.go": `package jscomp

import (
	"example.com/w44b_jscomp/gen/j2gort"
	pkg_rhino "example.com/w44b_jscomp/gen/rhino"
)

type NodeTraversal struct {
	scope *Scope
}
func (t *NodeTraversal) GetScope() *Scope { return t.scope }
type NodeTraversal_Builder struct{}
func clinit_NodeTraversal() {}
func InitNodeTraversal(*NodeTraversal) {}
func StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER() any { return nil }

type NodeTraversal_Callback interface {
	Implements_jscomp_NodeTraversal_Callback()
	ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool
	Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}
type NodeTraversal_ScopedCallback interface {
	NodeTraversal_Callback
	Implements_jscomp_NodeTraversal_ScopedCallback()
	EnterScope(*NodeTraversal)
	ExitScope(*NodeTraversal)
}

type NodeUtil_AllVarsDeclaredInFunction struct{ s *Scope }
func NodeUtil_GetAllVarsDeclaredInFunction(c AbstractCompiler, sc ScopeCreator, fnScope *Scope) *NodeUtil_AllVarsDeclaredInFunction {
	_ = c
	_ = sc
	return &NodeUtil_AllVarsDeclaredInFunction{s: fnScope}
}

type LiveVariablesAnalysis struct {
	scope   *Scope
	cleared int32
}
func NewLiveVariablesAnalysis(cfg any, fnScope, blkScope *Scope, c AbstractCompiler, sc ScopeCreator, allVars *NodeUtil_AllVarsDeclaredInFunction) *LiveVariablesAnalysis {
	_ = cfg; _ = blkScope; _ = c; _ = sc; _ = allVars
	return &LiveVariablesAnalysis{scope: fnScope}
}
func (l *LiveVariablesAnalysis) ResetForScope(cfg any, fnScope, blkScope *Scope, c AbstractCompiler, sc ScopeCreator, allVars *NodeUtil_AllVarsDeclaredInFunction) {
	_ = cfg; _ = blkScope; _ = c; _ = sc; _ = allVars
	l.scope = fnScope
}
func (l *LiveVariablesAnalysis) ClearScopeState() { l.scope = nil; l.cleared++ }
func (l *LiveVariablesAnalysis) Analyze()         {}

type DeadAssignmentsElimination struct {
	compiler AbstractCompiler
	liveness *LiveVariablesAnalysis
}
func (*DeadAssignmentsElimination) Implements_jscomp_NodeTraversal_Callback()       {}
func (*DeadAssignmentsElimination) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (this_ *DeadAssignmentsElimination) ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool { return true }
func (this_ *DeadAssignmentsElimination) Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) {}
func (this_ *DeadAssignmentsElimination) EnterScope(*NodeTraversal) {}
func (this_ *DeadAssignmentsElimination) ExitScope(t *NodeTraversal) {
	this_.priv_DeadAssignmentsElimination_eliminateDeadAssignments_2341(t)
}
func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_eliminateDeadAssignments_2341(t *NodeTraversal) {
	blockScope := t.GetScope()
	functionScope := blockScope.parent
	clinit_SyntacticScopeCreator()
	var val_3 SyntacticScopeCreator
	InitSyntacticScopeCreator_1(&val_3, this_.compiler)
	var scopeCreator *SyntacticScopeCreator = &val_3
	_ = scopeCreator
	var allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction = NodeUtil_GetAllVarsDeclaredInFunction(this_.compiler, scopeCreator, functionScope)
	_ = allVarsDeclaredInFunction
	if this_.liveness == nil {
		this_.liveness = NewLiveVariablesAnalysis(nil, functionScope, blockScope, this_.compiler, scopeCreator, allVarsDeclaredInFunction)
	} else {
		this_.liveness.ResetForScope(nil, functionScope, blockScope, this_.compiler, scopeCreator, allVarsDeclaredInFunction)
	}
	defer func() {
		this_.liveness.ClearScopeState()
	}()
	this_.liveness.Analyze()
}

type MustBeReachingVariableDef struct{ s *Scope }
func NewMustBeReachingVariableDef(s *Scope) *MustBeReachingVariableDef { return &MustBeReachingVariableDef{s: s} }
func (m *MustBeReachingVariableDef) ResetForScope(s *Scope)            { m.s = s }
func (m *MustBeReachingVariableDef) ClearScopeState()                  { m.s = nil }

type MaybeReachingVariableUse struct{ s *Scope }
func NewMaybeReachingVariableUse(s *Scope) *MaybeReachingVariableUse { return &MaybeReachingVariableUse{s: s} }
func (m *MaybeReachingVariableUse) ResetForScope(s *Scope)           { m.s = s }
func (m *MaybeReachingVariableUse) ClearScopeState()                 { m.s = nil }

type FlowSensitiveInlineVariables struct {
	compiler             AbstractCompiler
	cfg                  *Scope
	reachingDef          *MustBeReachingVariableDef
	reachingUses         *MaybeReachingVariableUse
	reusableReachingDef  *MustBeReachingVariableDef
	reusableReachingUses *MaybeReachingVariableUse
	reusableEscaped      j2gort.Ext_java_util_HashSet
}
func (*FlowSensitiveInlineVariables) Implements_jscomp_NodeTraversal_Callback()       {}
func (*FlowSensitiveInlineVariables) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (this_ *FlowSensitiveInlineVariables) ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool { return true }
func (this_ *FlowSensitiveInlineVariables) Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) {}
func (this_ *FlowSensitiveInlineVariables) ExitScope(*NodeTraversal) {}
func (this_ *FlowSensitiveInlineVariables) initFields_FlowSensitiveInlineVariables() {
	this_.reusableEscaped = j2gort.NewExt_java_util_HashSet()
}
func (this_ *FlowSensitiveInlineVariables) EnterScope(t *NodeTraversal) {
	defer func() {
		if this_.reachingDef != nil {
			this_.reachingDef.ClearScopeState()
		}
		if this_.reachingUses != nil {
			this_.reachingUses.ClearScopeState()
		}
		if !(this_.reusableEscaped.IsEmpty()) {
			this_.reusableEscaped.Clear()
		}
		this_.cfg = nil
		this_.reachingDef = nil
		this_.reachingUses = nil
	}()
	s := t.GetScope()
	this_.cfg = s
	if this_.reusableReachingDef == nil {
		this_.reusableReachingDef = NewMustBeReachingVariableDef(s)
	} else {
		this_.reusableReachingDef.ResetForScope(s)
	}
	this_.reachingDef = this_.reusableReachingDef
	if this_.reusableReachingUses == nil {
		this_.reusableReachingUses = NewMaybeReachingVariableUse(s)
	} else {
		this_.reusableReachingUses.ResetForScope(s)
	}
	this_.reachingUses = this_.reusableReachingUses
}

// --- NodeUtil_Has ---
func NodeUtil_Has(node pkg_rhino.Node, pred j2gort.Ext_com_google_common_base_Predicate, traverseChildrenPred j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = node
	_ = pred
	_ = traverseChildrenPred
	if j2gort.PredTestNode(pred, uint32(node)) {
		return true
	}
	if !(j2gort.PredTestNode(traverseChildrenPred, uint32(node))) {
		return false
	}
	{
		var c pkg_rhino.Node = node.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if NodeUtil_Has(c, pred, traverseChildrenPred) {
				return true
			}
		}
	}
	return false
}
`,
	})

	jscompDir := filepath.Join(tmpDir, "gen/jscomp")
	wave44bRunAllAndCheckIdempotent(t, jscompDir)

	readFile := func(name string) string {
		b, err := os.ReadFile(filepath.Join(jscompDir, name))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	svGot := readFile("scope_var_rc_iv.go")
	for _, want := range []string{
		"referenceMap             j2gort.OrderedMap[*Var, *ReferenceCollection]",
		"AcquireOrderedMap[*Var, *ReferenceCollection]()",
		"AcquireOrderedMapCap[*Var, *ReferenceCollection](int32(1024))",
		"sc_.implicitVars != nil",
		"currentScopeHandledVarAnalysesMap j2gort.OrderedMap[*Var, *InlineVariables_InlineVarAnalysis]",
		"p_ := root.GetParent()",
		"switch p_.GetToken()",
	} {
		if !strings.Contains(svGot, want) {
			t.Fatalf("expected scope_var_rc_iv.go to contain %q:\n%s", want, svGot)
		}
	}
	// Verify allocPooled_BasicBlock no longer calls InitBasicBlock.
	if idx := strings.Index(svGot, "func allocPooled_BasicBlock"); idx >= 0 {
		end := strings.Index(svGot[idx:], "func InitBasicBlock")
		if end > 0 && strings.Contains(svGot[idx:idx+end], "InitBasicBlock(") {
			t.Fatalf("expected allocPooled_BasicBlock to inline zeroed-slab init without calling InitBasicBlock:\n%s", svGot[idx:idx+end])
		}
	}

	rucGot := readFile("ruc.go")
	for _, want := range []string{
		"res_.RemoveUnusedCode_RemovableBase.outer_ = rucOuter_",
		"val_14.varName = var_.GetName()",
		"val_14.isEntirelyRemovable = true",
		"if builder.continuation0 !=",
	} {
		if !strings.Contains(rucGot, want) {
			t.Fatalf("expected ruc.go to contain %q:\n%s", want, rucGot)
		}
	}
	if strings.Contains(rucGot, "InitRemoveUnusedCode_VarInfo_1(val_14") {
		t.Fatalf("expected InitRemoveUnusedCode_VarInfo_1(val_14, ...) replaced with direct stores:\n%s", rucGot)
	}
	if strings.Contains(rucGot, "InitRemoveUnusedCode_IndirectAssign(res_") ||
		strings.Contains(rucGot, "InitRemoveUnusedCode_IndirectAssign(this_,") {
		t.Fatalf("expected BuildIndirectAssign and allocPooled_RemoveUnusedCode_IndirectAssign to inline direct zeroed-slab init:\n%s", rucGot)
	}

	daeGot := readFile("dae_fsiv_nodeutil.go")
	for _, want := range []string{
		"this_.liveness.ClearScopeState()",
		"scopeCreator.ReleaseSlabs_()",
		"func (*DeadAssignmentsElimination) CanReleaseSyntacticScopes_()",
		"func (*FlowSensitiveInlineVariables) CanReleaseSyntacticScopes_()",
		"j2gort.UnwrapNodePred(pred)",
		"j2gort.IsAlwaysTruePred(traverseChildrenPred)",
		"nodeUtilHasPredAlwaysTrueFn_(",
		"nodeUtilHasBothFn_(",
		"nodeUtilHasSlowAlwaysTrue_(",
		"nodeUtilHasSlow_(",
	} {
		if !strings.Contains(daeGot, want) {
			t.Fatalf("expected dae_fsiv_nodeutil.go to contain %q:\n%s", want, daeGot)
		}
	}
	if strings.Contains(daeGot, "this_.liveness = nil") {
		t.Fatalf("DeadAssignmentsElimination must call ClearScopeState() without nilling reusable this_.liveness:\n%s", daeGot)
	}

	wave44bWriteFiles(t, tmpDir, map[string]string{
		"gen/jscomp/wave44b_rt_test.go": `package jscomp

import (
	"reflect"
	"testing"

	"example.com/w44b_jscomp/gen/j2gort"
	pkg_rhino "example.com/w44b_jscomp/gen/rhino"
)

type customSlowPred struct{ want uint32 }
func (p *customSlowPred) Apply(v any) bool { return v.(uint32) == p.want }
func (p *customSlowPred) Test(v any) bool  { return v.(uint32) == p.want }

func TestWave44BRuntimeParity(t *testing.T) {
	// 1. Scope zero internal padding + contiguous pointers before depth int32
	scopeT := reflect.TypeFor[Scope]()
	if reflect.TypeFor[*Scope]().Size() == 8 {
		order := []string{"AbstractScopeBase", "parent", "var0", "var1", "var2", "var3", "var4", "var5", "vars", "varSlab_", "depth"}
		var prevEnd uintptr
		for i, fname := range order {
			sf, ok := scopeT.FieldByName(fname)
			if !ok {
				t.Fatalf("missing Scope.%s", fname)
			}
			if i > 0 && sf.Offset != prevEnd {
				t.Fatalf("padding detected before Scope.%s: offset=%d prevEnd=%d", fname, sf.Offset, prevEnd)
			}
			prevEnd = sf.Offset + sf.Type.Size()
		}
	}

	// 2. BasicBlock fast zeroed-slab init
	pkg_rhino.RegisterNode(10, pkg_rhino.Token_WHILE, 0, 11, 0)
	pkg_rhino.RegisterNode(11, pkg_rhino.Token_BLOCK, 10, 12, 0)
	pkg_rhino.RegisterNode(12, pkg_rhino.Token_FUNCTION, 11, 0, 0)
	var bbSlab j2gort.PooledSlab[BasicBlock]
	b0 := allocPooled_BasicBlock(&bbSlab, nil, 10)
	if b0.parent != nil || b0.isFunction || b0.isLoop || b0.root != 10 {
		t.Fatalf("unexpected b0: %+v", b0)
	}
	b1 := allocPooled_BasicBlock(&bbSlab, b0, 11)
	if b1.parent != b0 || b1.isFunction || !b1.isLoop || b1.root != 11 {
		t.Fatalf("unexpected b1: %+v", b1)
	}
	b2 := allocPooled_BasicBlock(&bbSlab, b1, 12)
	if b2.parent != b1 || !b2.isFunction || b2.isLoop || b2.root != 12 {
		t.Fatalf("unexpected b2: %+v", b2)
	}

	// 3. IndirectAssign + VarInfo direct zeroed-slab init
	pkg_rhino.RegisterNode(20, pkg_rhino.Token_PARAM_LIST, 0, 21, 0)
	pkg_rhino.RegisterNode(21, pkg_rhino.Token_NAME, 20, 0, 0)
	ruc := &RemoveUnusedCode{}
	bld := &RemoveUnusedCode_RemovableBuilder{
		outer_:                     ruc,
		propertyName:               "foo",
		isThisDotPropertyReference: true,
		continuation0:              RemoveUnusedCode_Continuation{node: 21},
	}
	ia := bld.BuildIndirectAssign(21, 21)
	if ia.outer_ != ruc || ia.RemoveUnusedCode_RemovableBase.outer_ != ruc || ia.root != 21 || ia.targetNode != 21 ||
		ia.propertyName != "foo" || !ia.isThisDotPropertyReference || ia.isPrototypeDotPropertyReference || ia.continuation0.node != 21 {
		t.Fatalf("unexpected IndirectAssign: %+v", ia)
	}
	v := &Var{AbstractVarBase: AbstractVarBase{ScopedNameBase: ScopedNameBase{name: "x"}, index: 0}}
	vi := ruc.priv_RemoveUnusedCode_getVarInfo_4732(v)
	if vi.outer_ != ruc || vi.varName != "x" || !vi.isEntirelyRemovable || vi.canonicalUnremovable {
		t.Fatalf("unexpected VarInfo: %+v", vi)
	}

	// 4. Canonical-Var OrderedMap (0 HashCode/Equals calls) + implicit-var MakeImplicitVar(EXPORTS) fallback
	sc := &Scope{}
	sc.BindSelf_(sc)
	exportsCanon := &Var{AbstractVarBase: AbstractVarBase{ScopedNameBase: ScopedNameBase{name: "exports"}, scope: sc, index: -1}}
	sc.implicitVars = exportsCanon
	exportsUncached := &Var{AbstractVarBase: AbstractVarBase{ScopedNameBase: ScopedNameBase{name: "exports"}, scope: sc, index: -1}}

	rc := &ReferenceCollector{}
	rc.initFields_ReferenceCollector()
	rc.Process_2(0, 0)
	VarHashCodeCalls.Store(0)
	VarEqualsCalls.Store(0)
	rc.AddRef(v)
	rc.AddRef(exportsCanon)
	if gotRC := j2gort.CastPtr[ReferenceCollection](rc.GetReferences(v)); gotRC == nil || gotRC.count != 1 {
		t.Fatalf("expected canonical var lookup hit, got %+v", gotRC)
	}
	if gotExp := j2gort.CastPtr[ReferenceCollection](rc.GetReferences(exportsUncached)); gotExp == nil || gotExp.count != 1 {
		t.Fatalf("expected uncached MakeImplicitVar(EXPORTS) lookup hit via Scope fallback, got %+v", gotExp)
	}
	if VarHashCodeCalls.Load() != 0 || VarEqualsCalls.Load() != 0 {
		t.Fatalf("expected 0 HashCode/Equals calls on OrderedMap[*Var, ...], got hash=%d eq=%d", VarHashCodeCalls.Load(), VarEqualsCalls.Load())
	}

	ib := &InlineVariables_InliningBehavior{}
	ib.initFields_InlineVariables_InliningBehavior()
	an := &InlineVariables_InlineVarAnalysis{ok: true}
	if gotAn := ib.PutAndGet(v, an); gotAn != an {
		t.Fatalf("expected InlineVarAnalysis hit")
	}
	if VarHashCodeCalls.Load() != 0 || VarEqualsCalls.Load() != 0 {
		t.Fatalf("expected 0 HashCode/Equals calls on currentScopeHandledVarAnalysesMap")
	}

	// 5. DAE reusable liveness preserved + ClearScopeState called
	fnScope := &Scope{}
	blkScope := &Scope{parent: fnScope}
	dae := &DeadAssignmentsElimination{}
	dae.ExitScope(&NodeTraversal{scope: blkScope})
	firstLiveness := dae.liveness
	if firstLiveness == nil || firstLiveness.scope != nil || firstLiveness.cleared == 0 {
		t.Fatalf("expected reusable liveness kept non-nil with scope cleared, got %+v", firstLiveness)
	}
	dae.ExitScope(&NodeTraversal{scope: blkScope})
	if dae.liveness != firstLiveness {
		t.Fatalf("expected same reusable liveness instance across scopes")
	}

	// 6. NodeUtil_Has UnwrapNodePred + IsAlwaysTruePred fast dispatch (0 PredTestNode calls for NodeLambda+AlwaysTrue!)
	j2gort.PredTestNodeCalls.Store(0)
	fastPred := j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(n uint32) bool { return n == 12 })
	if !NodeUtil_Has(10, fastPred, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue()) {
		t.Fatalf("expected NodeUtil_Has fastPred hit on node 12")
	}
	if j2gort.PredTestNodeCalls.Load() != 0 {
		t.Fatalf("expected 0 PredTestNode calls for NodeLambda + AlwaysTrue, got %d", j2gort.PredTestNodeCalls.Load())
	}
	stopAtBlock := j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(n uint32) bool { return n != 11 })
	if NodeUtil_Has(10, fastPred, stopAtBlock) {
		t.Fatalf("expected NodeUtil_Has to stop at node 11")
	}
	if j2gort.PredTestNodeCalls.Load() != 0 {
		t.Fatalf("expected 0 PredTestNode calls for NodeLambda + NodeLambda, got %d", j2gort.PredTestNodeCalls.Load())
	}
	if !NodeUtil_Has(10, &customSlowPred{want: 12}, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue()) {
		t.Fatalf("expected customSlowPred fallback hit on node 12")
	}
}
`,
	})
	cmd := exec.Command("go", "test", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("runtime test failed: %v\n%s", err, out)
	}
}

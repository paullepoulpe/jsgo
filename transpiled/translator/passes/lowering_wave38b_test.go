package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

const wave38bMockJ2gortSrc = `package j2gort

import (
	"fmt"
	"reflect"
	"strings"
)

const NullString = "\xff\x00j2go_null"

func ThrowNPE() { panic("NPE") }

func isNilVal(v any) bool {
	if v == nil {
		return true
	}
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return rv.IsNil()
	}
	return false
}

func ToAny(v any) any {
	if v == nil || isNilVal(v) {
		return nil
	}
	if s, ok := v.(string); ok && s == NullString {
		return nil
	}
	return v
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
	if v == nil || isNilVal(v) {
		return nil
	}
	return v.(*T)
}

func Cast[T any](v any) T {
	var zero T
	if v == nil || isNilVal(v) {
		return zero
	}
	return v.(T)
}

func UncheckedCast[T any](v any) T { return Cast[T](v) }
func AsIface[I any](v any) I       { return Cast[I](v) }

func CheckNotNull[T any](v T) T {
	if any(v) == nil || isNilVal(v) {
		panic("NPE")
	}
	if s, ok := any(v).(string); ok && s == NullString {
		panic("NPE")
	}
	return v
}

func CheckNotNullPtr[T any](v *T) *T {
	if v == nil {
		panic("NPE")
	}
	return v
}

func CheckNotNullString(s string) string {
	if s == NullString {
		panic("NPE")
	}
	return s
}

func CheckNotNullIface[I interface{ IsJavaNil_() bool }](v I) I {
	if any(v) == nil || v.IsJavaNil_() {
		panic("NPE")
	}
	return v
}

func RefEquals(a, b any) bool {
	if a == nil || isNilVal(a) {
		return b == nil || isNilVal(b)
	}
	if b == nil || isNilVal(b) {
		return false
	}
	return a == b
}

func StrVal(v any) string {
	if v == nil || isNilVal(v) {
		return "null"
	}
	if s, ok := v.(string); ok {
		if s == NullString {
			return "null"
		}
		return s
	}
	return fmt.Sprint(v)
}

func StrConcat(parts ...any) string {
	var sb strings.Builder
	for _, p := range parts {
		sb.WriteString(StrVal(p))
	}
	return sb.String()
}

type BoxInt struct{ Val int32 }

func ValueOfInt(v int32) *BoxInt { return &BoxInt{Val: v} }
func UnboxInt(b *BoxInt) int32   { return CheckNotNullPtr(b).Val }
func Int32Val(v any) int32       { return v.(*BoxInt).Val }

// --- Java Collection / Map / Set mocks ---

type Ext_java_util_function_Function interface {
	IsJavaNil_() bool
	Apply(v any) any
}

type fnImpl struct{ f func(any) any }

func (fn *fnImpl) IsJavaNil_() bool { return fn == nil }
func (fn *fnImpl) Apply(v any) any  { return fn.f(v) }
func NewLambda_Ext_java_util_function_Function(f func(any) any) Ext_java_util_function_Function {
	return &fnImpl{f: f}
}

type Ext_java_lang_Iterable interface {
	IsJavaNil_() bool
	IterSlice_() []any
}

type Ext_java_util_Collection interface {
	Ext_java_lang_Iterable
	Size() int32
	IsEmpty() bool
	Clear()
	Add_1(v any) bool
	Contains(v any) bool
	Remove_1_any(v any) bool
}

type Ext_java_util_Set interface {
	Ext_java_util_Collection
}

type Ext_java_util_LinkedHashSet struct {
	items []any
}

func (s *Ext_java_util_LinkedHashSet) IsJavaNil_() bool { return s == nil }
func (s *Ext_java_util_LinkedHashSet) Size() int32      { return int32(len(s.items)) }
func (s *Ext_java_util_LinkedHashSet) IsEmpty() bool    { return len(s.items) == 0 }
func (s *Ext_java_util_LinkedHashSet) Clear()           { s.items = s.items[:0] }
func (s *Ext_java_util_LinkedHashSet) IterSlice_() []any {
	cp := make([]any, len(s.items))
	copy(cp, s.items)
	return cp
}
func (s *Ext_java_util_LinkedHashSet) Contains(v any) bool {
	for _, x := range s.items {
		if x == v {
			return true
		}
	}
	return false
}
func (s *Ext_java_util_LinkedHashSet) Add_1(v any) bool {
	if s.Contains(v) {
		return false
	}
	s.items = append(s.items, v)
	return true
}
func (s *Ext_java_util_LinkedHashSet) Remove_1_any(v any) bool {
	for i, x := range s.items {
		if x == v {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}

type Ext_java_util_HashSet struct {
	Ext_java_util_LinkedHashSet
}

func NewExt_java_util_HashSet() *Ext_java_util_HashSet { return &Ext_java_util_HashSet{} }
func NewExt_java_util_HashSet_1_int32(cap int32) *Ext_java_util_HashSet {
	if cap < 0 {
		panic("IllegalArgumentException")
	}
	return &Ext_java_util_HashSet{Ext_java_util_LinkedHashSet{items: make([]any, 0, cap)}}
}
func Ext_com_google_common_collect_Sets_NewHashSet_0() *Ext_java_util_HashSet {
	return &Ext_java_util_HashSet{}
}
func Ext_com_google_common_collect_Sets_NewHashSetWithExpectedSize(cap int32) *Ext_java_util_HashSet {
	return NewExt_java_util_HashSet_1_int32(cap)
}
func NewExt_java_util_LinkedHashSet() *Ext_java_util_LinkedHashSet {
	return &Ext_java_util_LinkedHashSet{}
}

type JIterCursor struct {
	slice []any
	idx   int
}

func NewIterCursor(coll any) JIterCursor {
	if coll == nil || isNilVal(coll) {
		panic("NPE")
	}
	return JIterCursor{slice: coll.(Ext_java_lang_Iterable).IterSlice_()}
}
func (c *JIterCursor) HasNext() bool    { return c.idx < len(c.slice) }
func (c *JIterCursor) Next() any        { v := c.slice[c.idx]; c.idx++; return v }
func (c *JIterCursor) NextStr() string  { return CastString(c.Next()) }

type mapEntry struct {
	k any
	v any
}

type Ext_java_util_Map interface {
	IsJavaNil_() bool
	Size() int32
	IsEmpty() bool
	Clear()
	ContainsKey(k any) bool
	Get(k any) any
	GetOrDefault(k any, def any) any
	Put_2(k any, v any) any
	PutIfAbsent(k any, v any) any
	Remove_1(k any) any
	ComputeIfAbsent(k any, fn Ext_java_util_function_Function) any
	entries_() []mapEntry
}

type Ext_java_util_LinkedHashMap struct {
	ents []mapEntry
}

func (m *Ext_java_util_LinkedHashMap) IsJavaNil_() bool  { return m == nil }
func (m *Ext_java_util_LinkedHashMap) Size() int32       { return int32(len(m.ents)) }
func (m *Ext_java_util_LinkedHashMap) IsEmpty() bool     { return len(m.ents) == 0 }
func (m *Ext_java_util_LinkedHashMap) Clear()            { m.ents = m.ents[:0] }
func (m *Ext_java_util_LinkedHashMap) entries_() []mapEntry { return m.ents }
func (m *Ext_java_util_LinkedHashMap) ContainsKey(k any) bool {
	for _, e := range m.ents {
		if e.k == k {
			return true
		}
	}
	return false
}
func (m *Ext_java_util_LinkedHashMap) Get(k any) any {
	for _, e := range m.ents {
		if e.k == k {
			return e.v
		}
	}
	return nil
}
func (m *Ext_java_util_LinkedHashMap) GetOrDefault(k any, def any) any {
	for _, e := range m.ents {
		if e.k == k {
			return e.v
		}
	}
	return def
}
func (m *Ext_java_util_LinkedHashMap) Put_2(k any, v any) any {
	for i := range m.ents {
		if m.ents[i].k == k {
			old := m.ents[i].v
			m.ents[i].v = v
			return old
		}
	}
	m.ents = append(m.ents, mapEntry{k: k, v: v})
	return nil
}
func (m *Ext_java_util_LinkedHashMap) PutIfAbsent(k any, v any) any {
	for i := range m.ents {
		if m.ents[i].k == k {
			if m.ents[i].v == nil {
				m.ents[i].v = v
				return nil
			}
			return m.ents[i].v
		}
	}
	m.ents = append(m.ents, mapEntry{k: k, v: v})
	return nil
}
func (m *Ext_java_util_LinkedHashMap) Remove_1(k any) any {
	for i, e := range m.ents {
		if e.k == k {
			m.ents = append(m.ents[:i], m.ents[i+1:]...)
			return e.v
		}
	}
	return nil
}
func (m *Ext_java_util_LinkedHashMap) ComputeIfAbsent(k any, fn Ext_java_util_function_Function) any {
	if existing := m.Get(k); existing != nil {
		return existing
	}
	computed := fn.Apply(k)
	if computed != nil {
		m.Put_2(k, computed)
	}
	return computed
}

type Ext_java_util_HashMap struct {
	Ext_java_util_LinkedHashMap
}

func NewExt_java_util_HashMap() *Ext_java_util_HashMap { return &Ext_java_util_HashMap{} }
func NewExt_java_util_HashMap_1_int32(cap int32) *Ext_java_util_HashMap {
	if cap < 0 {
		panic("IllegalArgumentException")
	}
	return &Ext_java_util_HashMap{Ext_java_util_LinkedHashMap{ents: make([]mapEntry, 0, cap)}}
}
func Ext_com_google_common_collect_Maps_NewHashMap_0() *Ext_java_util_HashMap {
	return &Ext_java_util_HashMap{}
}
func Ext_com_google_common_collect_Maps_NewHashMapWithExpectedSize(cap int32) *Ext_java_util_HashMap {
	return NewExt_java_util_HashMap_1_int32(cap)
}
func NewExt_java_util_LinkedHashMap() *Ext_java_util_LinkedHashMap {
	return &Ext_java_util_LinkedHashMap{}
}

type MapKeyStrCursor struct {
	keys []string
	idx  int
}

func NewMapKeyStrCursor(m Ext_java_util_Map) MapKeyStrCursor {
	ents := m.entries_()
	ks := make([]string, len(ents))
	for i, e := range ents {
		ks[i] = CastString(e.k)
	}
	return MapKeyStrCursor{keys: ks}
}
func (c *MapKeyStrCursor) HasNext() bool  { return c.idx < len(c.keys) }
func (c *MapKeyStrCursor) Next() string   { k := c.keys[c.idx]; c.idx++; return k }
func (c *MapKeyStrCursor) NextStr() string { return c.Next() }

// --- OrderedMap / OrderedSet / LocalOrderedSet ---

type OrderedMapKeyCursor[K comparable] struct {
	keys []K
	idx  int
}

func (c *OrderedMapKeyCursor[K]) HasNext() bool { return c.idx < len(c.keys) }
func (c *OrderedMapKeyCursor[K]) Next() K       { k := c.keys[c.idx]; c.idx++; return k }

type OrderedMap[K comparable, V any] struct {
	keys   []K
	vals   map[K]V
	isInit bool
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{vals: make(map[K]V), isInit: true}
}
func NewOrderedMapCap[K comparable, V any](cap int32) OrderedMap[K, V] {
	if cap < 0 {
		panic("IllegalArgumentException")
	}
	return OrderedMap[K, V]{vals: make(map[K]V, cap), isInit: true}
}
func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	om := NewOrderedMap[K, V]()
	return &om
}
func NewOrderedMapPtrCap[K comparable, V any](cap int32) *OrderedMap[K, V] {
	om := NewOrderedMapCap[K, V](cap)
	return &om
}
func (m OrderedMap[K, V]) IsNil() bool           { return !m.isInit }
func (m OrderedMap[K, V]) Size() int32           { return int32(len(m.vals)) }
func (m OrderedMap[K, V]) IsEmpty() bool         { return len(m.vals) == 0 }
func (m OrderedMap[K, V]) ContainsKey(k K) bool  { _, ok := m.vals[k]; return ok }
func (m OrderedMap[K, V]) Get(k K) V             { return m.vals[k] }
func (m OrderedMap[K, V]) GetOk(k K) (V, bool)   { v, ok := m.vals[k]; return v, ok }
func (m OrderedMap[K, V]) GetOrDefault(k K, d V) V {
	if v, ok := m.vals[k]; ok {
		return v
	}
	return d
}
func (m *OrderedMap[K, V]) Put(k K, v V) V {
	if m.vals == nil {
		m.vals = make(map[K]V)
		m.isInit = true
	}
	old, existed := m.vals[k]
	if !existed {
		m.keys = append(m.keys, k)
	}
	m.vals[k] = v
	return old
}
func (m *OrderedMap[K, V]) Clear() {
	clear(m.vals)
	m.keys = m.keys[:0]
}
func (m *OrderedMap[K, V]) IterKeys() OrderedMapKeyCursor[K] {
	cp := make([]K, len(m.keys))
	copy(cp, m.keys)
	return OrderedMapKeyCursor[K]{keys: cp}
}
func OrderedMapGetStr[K comparable](m OrderedMap[K, string], k K) string {
	if v, ok := m.vals[k]; ok {
		return v
	}
	return NullString
}
func OrderedMapPtrGetStr[K comparable](m *OrderedMap[K, string], k K) string {
	return OrderedMapGetStr(*CheckNotNullPtr(m), k)
}

type OrderedSetCursor[T comparable] struct {
	items []T
	idx   int
}

func (c *OrderedSetCursor[T]) HasNext() bool { return c.idx < len(c.items) }
func (c *OrderedSetCursor[T]) Next() T       { v := c.items[c.idx]; c.idx++; return v }

type OrderedSet[T comparable] struct {
	items []T
	seen  map[T]bool
}

func NewOrderedSet[T comparable]() OrderedSet[T] {
	return OrderedSet[T]{seen: make(map[T]bool)}
}
func NewOrderedSetCap[T comparable](cap int32) OrderedSet[T] {
	return OrderedSet[T]{items: make([]T, 0, cap), seen: make(map[T]bool, cap)}
}
func (s OrderedSet[T]) Size() int32        { return int32(len(s.items)) }
func (s OrderedSet[T]) IsEmpty() bool      { return len(s.items) == 0 }
func (s OrderedSet[T]) Contains(v T) bool  { return s.seen[v] }
func (s *OrderedSet[T]) Add(v T) bool {
	if s.seen == nil {
		s.seen = make(map[T]bool)
	}
	if s.seen[v] {
		return false
	}
	s.seen[v] = true
	s.items = append(s.items, v)
	return true
}
func (s *OrderedSet[T]) Remove(v T) bool {
	if !s.seen[v] {
		return false
	}
	delete(s.seen, v)
	for i, x := range s.items {
		if x == v {
			s.items = append(s.items[:i], s.items[i+1:]...)
			break
		}
	}
	return true
}
func (s *OrderedSet[T]) Clear() {
	clear(s.seen)
	s.items = s.items[:0]
}
func (s *OrderedSet[T]) Iter() OrderedSetCursor[T] {
	cp := make([]T, len(s.items))
	copy(cp, s.items)
	return OrderedSetCursor[T]{items: cp}
}

type LocalOrderedSet[T comparable] struct {
	OrderedSet[T]
}

func NewLocalOrderedSet[T comparable]() LocalOrderedSet[T] {
	return LocalOrderedSet[T]{OrderedSet: NewOrderedSet[T]()}
}

// --- Wave 38 PlainMap / PlainSet ---

type PlainMap[K comparable, V any] struct {
	m      map[K]V
	isInit bool
}

func NewPlainMap[K comparable, V any](capHint int) PlainMap[K, V] {
	if capHint < 0 {
		panic("IllegalArgumentException")
	}
	return PlainMap[K, V]{m: make(map[K]V, capHint), isInit: true}
}
func NewPlainMapVal[K comparable, V any]() PlainMap[K, V]        { return NewPlainMap[K, V](0) }
func NewPlainMapCap[K comparable, V any](cap int32) PlainMap[K, V] { return NewPlainMap[K, V](int(cap)) }
func NewPlainMapPtr[K comparable, V any]() *PlainMap[K, V] {
	pm := NewPlainMap[K, V](0)
	return &pm
}
func NewPlainMapPtrCap[K comparable, V any](cap int32) *PlainMap[K, V] {
	pm := NewPlainMap[K, V](int(cap))
	return &pm
}
func (p *PlainMap[K, V]) InitCap(capHint int) {
	if capHint < 0 {
		panic("IllegalArgumentException")
	}
	p.m = make(map[K]V, capHint)
	p.isInit = true
}
func (p PlainMap[K, V]) IsNil() bool          { return !p.isInit }
func (p PlainMap[K, V]) IsJavaNil_() bool     { return !p.isInit }
func (p PlainMap[K, V]) Size() int32          { return int32(len(p.m)) }
func (p PlainMap[K, V]) IsEmpty() bool        { return len(p.m) == 0 }
func (p PlainMap[K, V]) ContainsKey(k K) bool { _, ok := p.m[k]; return ok }
func (p PlainMap[K, V]) Get(k K) V            { return p.m[k] }
func (p PlainMap[K, V]) GetOk(k K) (V, bool)  { v, ok := p.m[k]; return v, ok }
func (p PlainMap[K, V]) GetOrDefault(k K, d V) V {
	if v, ok := p.m[k]; ok {
		return v
	}
	return d
}
func (p *PlainMap[K, V]) ensure() {
	if p.m == nil {
		p.m = make(map[K]V)
		p.isInit = true
	}
}
func (p *PlainMap[K, V]) Put(k K, v V) V {
	p.ensure()
	old := p.m[k]
	p.m[k] = v
	return old
}
func (p *PlainMap[K, V]) PutOk(k K, v V) (V, bool) {
	p.ensure()
	old, ok := p.m[k]
	p.m[k] = v
	return old, ok
}
func (p *PlainMap[K, V]) PutVoid(k K, v V) {
	p.ensure()
	p.m[k] = v
}
func (p *PlainMap[K, V]) PutIfAbsent(k K, v V) V {
	p.ensure()
	if cur, ok := p.m[k]; ok {
		return cur
	}
	p.m[k] = v
	var zero V
	return zero
}
func (p *PlainMap[K, V]) PutIfAbsentOk(k K, v V) (V, bool) {
	p.ensure()
	if cur, ok := p.m[k]; ok {
		return cur, true
	}
	p.m[k] = v
	var zero V
	return zero, false
}
func (p *PlainMap[K, V]) Remove(k K) V {
	old := p.m[k]
	delete(p.m, k)
	return old
}
func (p *PlainMap[K, V]) RemoveOk(k K, _ ...any) (V, bool) {
	old, ok := p.m[k]
	delete(p.m, k)
	return old, ok
}
func (p *PlainMap[K, V]) Clear() {
	clear(p.m)
}

func PlainMapGetStr[K comparable](p PlainMap[K, string], k K) string {
	if v, ok := p.m[k]; ok {
		return v
	}
	return NullString
}
func PlainMapPtrGetStr[K comparable](p *PlainMap[K, string], k K) string {
	return PlainMapGetStr(*CheckNotNullPtr(p), k)
}
func PlainMapPutStr[K comparable](p *PlainMap[K, string], k K, v string) string {
	p.ensure()
	old, ok := p.m[k]
	p.m[k] = v
	if !ok {
		return NullString
	}
	return old
}
func PlainMapPutIfAbsentStr[K comparable](p *PlainMap[K, string], k K, v string) string {
	p.ensure()
	if cur, ok := p.m[k]; ok && cur != NullString {
		return cur
	}
	p.m[k] = v
	return NullString
}
func PlainMapRemoveStr[K comparable](p *PlainMap[K, string], k K) string {
	old, ok := p.m[k]
	delete(p.m, k)
	if !ok {
		return NullString
	}
	return old
}
func PlainMapComputeIfAbsent[K comparable, V comparable](p *PlainMap[K, V], k K, zeroVal V, fn func(K) V) V {
	if cur, ok := p.m[k]; ok && cur != zeroVal {
		return cur
	}
	val := fn(k)
	if val != zeroVal {
		p.PutVoid(k, val)
	}
	return val
}
func PlainMapComputeIfAbsentFn[K comparable, V comparable](p *PlainMap[K, V], k K, fn func(K) V) V {
	var zeroVal V
	return PlainMapComputeIfAbsent(p, k, zeroVal, fn)
}

type PlainSet[T comparable] struct {
	m map[T]struct{}
}

func NewPlainSet[T comparable](capHint int) PlainSet[T] {
	if capHint < 0 {
		panic("IllegalArgumentException")
	}
	return PlainSet[T]{m: make(map[T]struct{}, capHint)}
}
func NewPlainSetVal[T comparable]() PlainSet[T]          { return NewPlainSet[T](0) }
func NewPlainSetCap[T comparable](cap int32) PlainSet[T] { return NewPlainSet[T](int(cap)) }
func NewPlainSetPtr[T comparable]() *PlainSet[T] {
	ps := NewPlainSet[T](0)
	return &ps
}
func NewPlainSetPtrCap[T comparable](cap int32) *PlainSet[T] {
	ps := NewPlainSet[T](int(cap))
	return &ps
}
func (s *PlainSet[T]) InitCap(capHint int) {
	if capHint < 0 {
		panic("IllegalArgumentException")
	}
	s.m = make(map[T]struct{}, capHint)
}
func (s PlainSet[T]) Size() int32       { return int32(len(s.m)) }
func (s PlainSet[T]) IsEmpty() bool     { return len(s.m) == 0 }
func (s PlainSet[T]) Contains(v T) bool { _, ok := s.m[v]; return ok }
func (s *PlainSet[T]) ensure() {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
}
func (s *PlainSet[T]) Add(v T) bool {
	s.ensure()
	if _, ok := s.m[v]; ok {
		return false
	}
	s.m[v] = struct{}{}
	return true
}
func (s *PlainSet[T]) AddVoid(v T) {
	s.ensure()
	s.m[v] = struct{}{}
}
func (s *PlainSet[T]) Remove(v T) bool {
	if _, ok := s.m[v]; !ok {
		return false
	}
	delete(s.m, v)
	return true
}
func (s *PlainSet[T]) Clear() {
	clear(s.m)
}
`

func TestWave38B_AbstractBaseDevirtAndCrossPkgSoundness(t *testing.T) {
	rhinoSrc := `package rhino

import "MODROOT/j2gort"

type JSType interface {
	Implements_rhino_JSType()
	IsJavaNil_() bool
	GetJSTypeBase() *JSTypeBase
	IsUnknownType() bool
	IsAllType() bool
	IsNoType() bool
	IsUnionType() bool
	HasKindMask(int32) bool
	IsResolved() bool
	GetDisplayName() string
	CanBeCalled() bool
	CollapseUnion() JSType
}

func ToIface_JSType[T interface {
	JSType
	comparable
}](v T) JSType {
	var z T
	if v == z {
		return nil
	}
	return v
}

type JSTypeBase struct {
	self_        JSType
	typeKindMask int32
	resolved     bool
}

func (b *JSTypeBase) BindSelf_(s any)           { b.self_ = s.(JSType) }
func (b *JSTypeBase) IsJavaNil_() bool          { return b == nil }
func (b *JSTypeBase) GetJSTypeBase() *JSTypeBase { return b }
func (*JSTypeBase) Implements_rhino_JSType()    {}

func (b *JSTypeBase) IsUnknownType() bool    { return (b.typeKindMask & 1) != 0 }
func (b *JSTypeBase) IsAllType() bool        { return (b.typeKindMask & 2) != 0 }
func (b *JSTypeBase) IsNoType() bool         { return (b.typeKindMask & 4) != 0 }
func (b *JSTypeBase) IsUnionType() bool      { return (b.typeKindMask & 8) != 0 }
func (b *JSTypeBase) HasKindMask(m int32) bool { return (b.typeKindMask & m) != 0 }
func (b *JSTypeBase) IsResolved() bool       { return b.resolved }
func (b *JSTypeBase) GetDisplayName() string { return "base" }
func (b *JSTypeBase) CanBeCalled() bool      { return false }
func (b *JSTypeBase) CollapseUnion() JSType  { return b.self_ }

type ObjectType interface {
	Implements_rhino_ObjectType()
	IsJavaNil_() bool
	GetObjectTypeBase() *ObjectTypeBase
	JSType
	HasCachedValues() bool
	GetReferenceName() string
}

type ObjectTypeBase struct {
	JSTypeBase
	cached  bool
	refName string
}

func (b *ObjectTypeBase) BindSelf_(s any)                   { b.JSTypeBase.BindSelf_(s) }
func (b *ObjectTypeBase) IsJavaNil_() bool                  { return b == nil }
func (b *ObjectTypeBase) GetObjectTypeBase() *ObjectTypeBase { return b }
func (*ObjectTypeBase) Implements_rhino_ObjectType()        {}
func (*ObjectTypeBase) Implements_rhino_JSType()            {}
func (b *ObjectTypeBase) HasCachedValues() bool             { return b.cached }
func (b *ObjectTypeBase) GetReferenceName() string          { return b.refName }

type FunctionType struct {
	ObjectTypeBase
	customLabel string
}

func NewFunctionType(mask int32, label string, cached bool) *FunctionType {
	f := &FunctionType{customLabel: label}
	f.typeKindMask = mask
	f.resolved = true
	f.cached = cached
	f.refName = "fn:" + label
	f.BindSelf_(f)
	return f
}

func (f *FunctionType) IsJavaNil_() bool                 { return f == nil }
func (*FunctionType) Implements_rhino_ObjectType()       {}
func (*FunctionType) Implements_rhino_JSType()           {}
func (f *FunctionType) GetDisplayName() string           { return f.customLabel }
func (f *FunctionType) CanBeCalled() bool                { return true }
func (f *FunctionType) GetReferenceName() string         { return "override:" + f.customLabel }

type UnionType struct {
	JSTypeBase
	first JSType
}

func NewUnionType(first JSType) *UnionType {
	u := &UnionType{first: first}
	u.typeKindMask = 8
	u.resolved = true
	u.BindSelf_(u)
	return u
}

func (u *UnionType) IsJavaNil_() bool       { return u == nil }
func (*UnionType) Implements_rhino_JSType() {}
func (u *UnionType) CollapseUnion() JSType  { return u.first }

// SimpleSlot: GetDeclaration() is NOT overridden in rhino, ONLY overridden in jscomp.SymbolTable_Symbol!
type SimpleSlot interface {
	Implements_rhino_SimpleSlot()
	IsJavaNil_() bool
	GetSimpleSlotBase() *SimpleSlotBase
	GetName() string
	GetDeclaration() string
}

type SimpleSlotBase struct {
	self_ SimpleSlot
	name  string
}

func NewSimpleSlot(name string) SimpleSlot {
	s := &SimpleSlotBase{name: name}
	s.self_ = s
	return s
}

func (b *SimpleSlotBase) BindSelf_(s any)                   { b.self_ = s.(SimpleSlot) }
func (b *SimpleSlotBase) IsJavaNil_() bool                  { return b == nil }
func (b *SimpleSlotBase) GetSimpleSlotBase() *SimpleSlotBase { return b }
func (*SimpleSlotBase) Implements_rhino_SimpleSlot()        {}
func (b *SimpleSlotBase) GetName() string                   { return b.name }
func (b *SimpleSlotBase) GetDeclaration() string            { return "baseDecl:" + b.name }

func ReadSlotInRhino(slot SimpleSlot) string {
	return slot.GetName() + "|" + slot.GetDeclaration()
}

func CheckPairWithNilGuard(left JSType, right JSType) string {
	if left == right || ((left == nil || left.IsJavaNil_()) && (right == nil || right.IsJavaNil_())) {
		return "same"
	}
	if (left == nil || left.IsJavaNil_()) || (right == nil || right.IsJavaNil_()) {
		return "one-nil"
	}
	if left.IsUnknownType() || right.IsUnknownType() {
		return "unk"
	}
	if left.HasKindMask(2) || right.IsResolved() {
		return left.GetDisplayName() + ":" + j2gort.StrVal(left.CanBeCalled())
	}
	return "other"
}

func SubtypeHelperStyle(subtype JSType, supertype JSType, obj ObjectType, guardFlag bool, maybeNil JSType) string {
	if ((subtype == supertype) || supertype.IsUnknownType()) || subtype.IsUnknownType() {
		return "fast"
	}
	if subtype.IsAllType() || supertype.IsNoType() {
		return "bounds"
	}
	collapsed := subtype.CollapseUnion()
	if collapsed.IsUnionType() || subtype.HasKindMask(2) || supertype.HasKindMask(8) {
		return "union"
	}
	// Guarded nullable parameter must NOT be hoisted above guardFlag!
	if guardFlag {
		if maybeNil.IsUnknownType() && maybeNil.IsResolved() {
			return "guarded"
		}
	}
	if obj != nil && !obj.IsJavaNil_() {
		if obj.HasCachedValues() && obj.IsResolved() {
			return obj.GetReferenceName() + "|" + subtype.GetDisplayName()
		}
	}
	return "none"
}
`

	jscompSrc := `package jscomp

import pkg_rhino "MODROOT/rhino"

type SymbolTable_Symbol struct {
	pkg_rhino.SimpleSlotBase
	declTag string
}

func NewSymbolTable_Symbol(name, declTag string) *SymbolTable_Symbol {
	s := &SymbolTable_Symbol{declTag: declTag}
	InitSimpleSlotCompat(&s.SimpleSlotBase, name)
	s.BindSelf_(s)
	return s
}

func InitSimpleSlotCompat(b *pkg_rhino.SimpleSlotBase, name string) {
	*b = *(pkg_rhino.NewSimpleSlot(name).GetSimpleSlotBase())
}

func (s *SymbolTable_Symbol) IsJavaNil_() bool               { return s == nil }
func (*SymbolTable_Symbol) Implements_rhino_SimpleSlot()     {}
func (s *SymbolTable_Symbol) GetDeclaration() string         { return "symDecl:" + s.declTag }

type AbstractScope interface {
	Implements_jscomp_AbstractScope()
	IsJavaNil_() bool
	GetAbstractScopeBase() *AbstractScopeBase
	GetParentScope() AbstractScope
	IsLocal() bool
	IsGlobal() bool
	ScopeKindTag() string
}

type AbstractScopeBase struct {
	self_  AbstractScope
	parent AbstractScope
}

func (b *AbstractScopeBase) BindSelf_(s any)                        { b.self_ = s.(AbstractScope) }
func (b *AbstractScopeBase) IsJavaNil_() bool                       { return b == nil }
func (b *AbstractScopeBase) GetAbstractScopeBase() *AbstractScopeBase { return b }
func (*AbstractScopeBase) Implements_jscomp_AbstractScope()         {}
func (b *AbstractScopeBase) GetParentScope() AbstractScope          { return b.parent }
func (b *AbstractScopeBase) IsLocal() bool                          { return b.parent != nil }
func (b *AbstractScopeBase) IsGlobal() bool                         { return b.parent == nil }
func (b *AbstractScopeBase) ScopeKindTag() string                   { return "abstract" }

type TypedScope struct {
	AbstractScopeBase
}

func NewTypedScope(parent AbstractScope) *TypedScope {
	s := &TypedScope{}
	s.parent = parent
	s.BindSelf_(s)
	return s
}

func (s *TypedScope) IsJavaNil_() bool                   { return s == nil }
func (*TypedScope) Implements_jscomp_AbstractScope()     {}
func (s *TypedScope) ScopeKindTag() string               { return "typed" }

func WalkScopesAndTypes(sc AbstractScope, t pkg_rhino.JSType, slot pkg_rhino.SimpleSlot) string {
	// Reassigned local variable must NOT hoist a stale base pointer!
	cur := sc
	depth := 0
	for cur != nil && !cur.IsJavaNil_() && cur.IsLocal() {
		depth++
		cur = cur.GetParentScope()
	}
	isGlob := cur != nil && !cur.IsJavaNil_() && cur.IsGlobal()
	tag := sc.ScopeKindTag()
	tCheck := t.IsUnknownType() || t.HasKindMask(4) || !t.IsResolved()
	if tCheck {
		tag += ":t"
	}
	return tag + "|" + pkg_rhino.ReadSlotInRhino(slot) + "|" + slot.GetDeclaration() + "|" + func() string {
		if isGlob && depth > 0 {
			return "ok"
		}
		return "root"
	}()
}
`

	tmpDir := t.TempDir()
	writePkgSet := func(prefix string) {
		modImport := "example.com/w38b_devirt/" + prefix
		mustWrite := func(rel, content string) {
			t.Helper()
			full := filepath.Join(tmpDir, prefix, rel)
			if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
				t.Fatal(err)
			}
			replaced := strings.ReplaceAll(content, "MODROOT", modImport)
			if err := os.WriteFile(full, []byte(replaced), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		mustWrite("j2gort/runtime.go", wave38bMockJ2gortSrc)
		mustWrite("rhino/rhino.go", rhinoSrc)
		mustWrite("jscomp/jscomp.go", jscompSrc)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/w38b_devirt\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writePkgSet("orig")
	writePkgSet("opt")

	optPattern := filepath.Join(tmpDir, "opt") + "/..."
	allPasses := passes.AllPasses()

	s1, err := passes.Run([]string{optPattern}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.ByPass["abstract_base_devirt"] < 15 {
		t.Fatalf("expected >= 15 abstract_base_devirt sites on Run 1, got %+v", s1.ByPass)
	}

	// Run 2 MUST produce strictly 0 sites and 0 modified files!
	s2, err := passes.Run([]string{optPattern}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d ByPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	gotRhinoBytes, err := os.ReadFile(filepath.Join(tmpDir, "opt", "rhino", "rhino.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotRhino := string(gotRhinoBytes)

	for _, want := range []string{
		"base_left_ := left.GetJSTypeBase()",
		"base_right_ := right.GetJSTypeBase()",
		"base_subtype_ := subtype.GetJSTypeBase()",
		"base_supertype_ := supertype.GetJSTypeBase()",
		"slot.GetSimpleSlotBase().GetName()",
		"slot.GetDeclaration()",
		"left.GetDisplayName()",
		"left.CanBeCalled()",
		"subtype.CollapseUnion()",
		"maybeNil.GetJSTypeBase().IsUnknownType()",
		"obj.GetObjectTypeBase().HasCachedValues()",
		"obj.GetObjectTypeBase().IsResolved()",
		"obj.GetReferenceName()",
	} {
		if !strings.Contains(gotRhino, want) {
			t.Errorf("opt/rhino/rhino.go missing expected %q in:\n%s", want, gotRhino)
		}
	}
	// Cross-package override safety: SimpleSlot.GetDeclaration() is overridden in jscomp.SymbolTable_Symbol,
	// so rhino MUST NOT devirtualize slot.GetDeclaration() to slot.GetSimpleSlotBase().GetDeclaration()!
	for _, unwanted := range []string{
		"GetSimpleSlotBase().GetDeclaration()",
		"GetJSTypeBase().GetDisplayName()",
		"GetJSTypeBase().CanBeCalled()",
		"GetJSTypeBase().CollapseUnion()",
		"GetObjectTypeBase().GetReferenceName()",
		"base_maybeNil_",
	} {
		if strings.Contains(gotRhino, unwanted) {
			t.Errorf("opt/rhino/rhino.go contains unsound devirt %q in:\n%s", unwanted, gotRhino)
		}
	}

	gotJscompBytes, err := os.ReadFile(filepath.Join(tmpDir, "opt", "jscomp", "jscomp.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotJscomp := string(gotJscompBytes)
	for _, want := range []string{
		"cur.GetAbstractScopeBase().IsLocal()",
		"cur.GetAbstractScopeBase().GetParentScope()",
		"cur.GetAbstractScopeBase().IsGlobal()",
		"sc.ScopeKindTag()",
		"base_t_ := t.GetJSTypeBase()",
		"base_t_.IsUnknownType()",
	} {
		if !strings.Contains(gotJscomp, want) {
			t.Errorf("opt/jscomp/jscomp.go missing expected %q in:\n%s", want, gotJscomp)
		}
	}
	if strings.Contains(gotJscomp, "base_cur_") || strings.Contains(gotJscomp, "GetAbstractScopeBase().ScopeKindTag()") {
		t.Fatalf("opt/jscomp/jscomp.go hoisted mutated cur or devirtualized overridden ScopeKindTag:\n%s", gotJscomp)
	}

	diffTest := `package devirt_diff_test

import (
	"testing"

	opt_jscomp "example.com/w38b_devirt/opt/jscomp"
	opt_rhino "example.com/w38b_devirt/opt/rhino"
	orig_jscomp "example.com/w38b_devirt/orig/jscomp"
	orig_rhino "example.com/w38b_devirt/orig/rhino"
)

func runStr(fn func() string) (res string, panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	return fn(), false
}

func TestDevirtRuntimeParity(t *testing.T) {
	ofn1 := orig_rhino.NewFunctionType(2, "Alpha", true)
	ofn2 := orig_rhino.NewFunctionType(0, "Beta", true)
	ou := orig_rhino.NewUnionType(ofn1)
	var oTypedNil orig_rhino.JSType = (*orig_rhino.FunctionType)(nil)

	xfn1 := opt_rhino.NewFunctionType(2, "Alpha", true)
	xfn2 := opt_rhino.NewFunctionType(0, "Beta", true)
	xu := opt_rhino.NewUnionType(xfn1)
	var xTypedNil opt_rhino.JSType = (*opt_rhino.FunctionType)(nil)

	pairCases := [][2]int{
		{0, 0}, {0, 1}, {1, 2}, {2, 3}, {3, 3},
	}
	pickOrig := func(k int) orig_rhino.JSType {
		switch k {
		case 0:
			return nil
		case 1:
			return oTypedNil
		case 2:
			return ofn1
		default:
			return ofn2
		}
	}
	pickOpt := func(k int) opt_rhino.JSType {
		switch k {
		case 0:
			return nil
		case 1:
			return xTypedNil
		case 2:
			return xfn1
		default:
			return xfn2
		}
	}
	for i, pc := range pairCases {
		go1, p1 := runStr(func() string { return orig_rhino.CheckPairWithNilGuard(pickOrig(pc[0]), pickOrig(pc[1])) })
		gx1, px1 := runStr(func() string { return opt_rhino.CheckPairWithNilGuard(pickOpt(pc[0]), pickOpt(pc[1])) })
		if go1 != gx1 || p1 != px1 {
			t.Fatalf("CheckPair case %d mismatch: orig=(%q,%v) opt=(%q,%v)", i, go1, p1, gx1, px1)
		}
	}

	// Verify guardFlag=false with nil maybeNil does NOT panic, while guardFlag=true DOES panic identically!
	for _, gf := range []bool{false, true} {
		go2, p2 := runStr(func() string { return orig_rhino.SubtypeHelperStyle(ou, ofn2, ofn1, gf, nil) })
		gx2, px2 := runStr(func() string { return opt_rhino.SubtypeHelperStyle(xu, xfn2, xfn1, gf, nil) })
		if go2 != gx2 || p2 != px2 {
			t.Fatalf("SubtypeHelperStyle(gf=%v) mismatch: orig=(%q,%v) opt=(%q,%v)", gf, go2, p2, gx2, px2)
		}
	}

	// Verify cross-package SymbolTable_Symbol override of SimpleSlot.GetDeclaration()!
	oRoot := orig_jscomp.NewTypedScope(nil)
	oChild := orig_jscomp.NewTypedScope(oRoot)
	oSym := orig_jscomp.NewSymbolTable_Symbol("mySym", "line42")

	xRoot := opt_jscomp.NewTypedScope(nil)
	xChild := opt_jscomp.NewTypedScope(xRoot)
	xSym := opt_jscomp.NewSymbolTable_Symbol("mySym", "line42")

	gotOrig := orig_jscomp.WalkScopesAndTypes(oChild, ofn2, oSym)
	gotOpt := opt_jscomp.WalkScopesAndTypes(xChild, xfn2, xSym)
	if gotOrig != gotOpt || !containsStr(gotOpt, "symDecl:line42") {
		t.Fatalf("WalkScopesAndTypes mismatch: orig=%q opt=%q", gotOrig, gotOpt)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexStr(s, sub) >= 0)
}

func indexStr(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "devirt_diff_test.go"), []byte(diffTest), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-count=1", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("differential runtime test failed: %v\n%s", err, string(out))
	}
}

func TestWave38B_PeepholeEqualsJSTypeAndUpcastFallback(t *testing.T) {
	rhinoSrc := `package rhino

import "MODROOT/j2gort"

type JSType interface {
	Implements_rhino_JSType()
	IsJavaNil_() bool
	GetJSTypeBase() *JSTypeBase
	Equals(any) bool
	EqualsJSType(JSType) bool
	IsUnknownType() bool
}

func ToIface_JSType[T interface {
	JSType
	comparable
}](v T) JSType {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JSType[T JSType](v T) JSType {
	if any(v) == nil || v.IsJavaNil_() {
		return nil
	}
	return v
}

type JSTypeBase struct {
	self_ JSType
	id    int32
	unk   bool
}

func (b *JSTypeBase) BindSelf_(s any)           { b.self_ = s.(JSType) }
func (b *JSTypeBase) IsJavaNil_() bool          { return b == nil }
func (b *JSTypeBase) GetJSTypeBase() *JSTypeBase { return b }
func (*JSTypeBase) Implements_rhino_JSType()    {}
func (b *JSTypeBase) IsUnknownType() bool       { return b.unk }

func (b *JSTypeBase) EqualsJSType(that JSType) bool {
	if b == nil {
		panic("NPE")
	}
	if that == nil || that.IsJavaNil_() {
		return false
	}
	return b.id == that.GetJSTypeBase().id
}

func (b *JSTypeBase) Equals(other any) bool {
	if j2gort.RefEquals(b.self_, other) {
		return true
	}
	jt, ok := other.(JSType)
	return ok && b.EqualsJSType(jt)
}

type ObjectType interface {
	Implements_rhino_ObjectType()
	IsJavaNil_() bool
	GetObjectTypeBase() *ObjectTypeBase
	JSType
}

func Upcast_ObjectType[T ObjectType](v T) ObjectType {
	if any(v) == nil || v.IsJavaNil_() {
		return nil
	}
	return v
}

type ObjectTypeBase struct {
	JSTypeBase
}

func (b *ObjectTypeBase) BindSelf_(s any)                   { b.JSTypeBase.BindSelf_(s) }
func (b *ObjectTypeBase) IsJavaNil_() bool                  { return b == nil }
func (b *ObjectTypeBase) GetObjectTypeBase() *ObjectTypeBase { return b }
func (*ObjectTypeBase) Implements_rhino_ObjectType()        {}
func (*ObjectTypeBase) Implements_rhino_JSType()            {}

type FunctionType struct {
	ObjectTypeBase
}

func NewFunctionType(id int32) *FunctionType {
	f := &FunctionType{}
	f.id = id
	f.BindSelf_(f)
	return f
}

func (f *FunctionType) IsJavaNil_() bool           { return f == nil }
func (*FunctionType) Implements_rhino_ObjectType() {}
func (*FunctionType) Implements_rhino_JSType()     {}

func (this_ *FunctionType) MatchesOther(that *FunctionType, obj ObjectType, jt JSType, rawAny any) (bool, JSType, ObjectType, JSType) {
	eqPtr := this_.Equals(j2gort.PtrToAny(that))
	eqObj := obj.Equals(jt)
	eqAny := jt.Equals(rawAny)
	upNullable := j2gort.AsIface[JSType](that)
	upShort := j2gort.AsIface[ObjectType](&FunctionType{})
	upNested := j2gort.AsIface[JSType](j2gort.AsIface[ObjectType](that))
	keepAny := j2gort.AsIface[JSType](rawAny)
	_ = keepAny
	return (eqPtr || eqObj) && !eqAny, upNullable, upShort, upNested
}
`

	jscompSrc := `package jscomp

import (
	"MODROOT/j2gort"
	pkg_rhino "MODROOT/rhino"
)

type LegacyIface interface {
	IsJavaNil_() bool
	Ping() int32
}

type LegacyImpl struct{ v int32 }

func (l *LegacyImpl) IsJavaNil_() bool { return l == nil }
func (l *LegacyImpl) Ping() int32      { return l.v }

func CheckTypesInJscomp(a pkg_rhino.JSType, b pkg_rhino.ObjectType, fn1 *pkg_rhino.FunctionType, fn2 *pkg_rhino.FunctionType, leg *LegacyImpl) (bool, pkg_rhino.JSType, LegacyIface) {
	if a.IsUnknownType() || b.IsUnknownType() {
		return false, nil, nil
	}
	same1 := a.Equals(b)
	same2 := fn1.Equals(j2gort.PtrToAny(fn2))
	upCast := j2gort.AsIface[pkg_rhino.JSType](fn2)
	keepNoUpcast := j2gort.AsIface[LegacyIface](leg)
	return same1 && same2, upCast, keepNoUpcast
}
`

	tmpDir := t.TempDir()
	writeSet := func(prefix string) {
		modImport := "example.com/w38b_peep/" + prefix
		mustWrite := func(rel, content string) {
			t.Helper()
			full := filepath.Join(tmpDir, prefix, rel)
			if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(full, []byte(strings.ReplaceAll(content, "MODROOT", modImport)), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		mustWrite("j2gort/runtime.go", wave38bMockJ2gortSrc)
		mustWrite("rhino/rhino.go", rhinoSrc)
		mustWrite("jscomp/jscomp.go", jscompSrc)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/w38b_peep\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeSet("orig")
	writeSet("opt")

	optPattern := filepath.Join(tmpDir, "opt") + "/..."
	allPasses := passes.AllPasses()

	s1, err := passes.Run([]string{optPattern}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.ByPass["j2gort_peephole"] < 7 {
		t.Fatalf("expected >= 7 j2gort_peephole sites on Run 1, got %+v", s1.ByPass)
	}

	// Run 2 MUST produce strictly 0 sites and 0 modified files!
	s2, err := passes.Run([]string{optPattern}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d ByPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	gotRhinoBytes, err := os.ReadFile(filepath.Join(tmpDir, "opt", "rhino", "rhino.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotRhino := string(gotRhinoBytes)
	for _, want := range []string{
		"this_.GetJSTypeBase().EqualsJSType(ToIface_JSType(that))",
		"obj.GetObjectTypeBase().GetJSTypeBase().EqualsJSType(jt)",
		"jt.GetJSTypeBase().Equals(rawAny)",
		"upNullable := Upcast_JSType(that)",
		"upShort := Upcast_ObjectType(&FunctionType{})",
		"upNested := Upcast_JSType(that)",
		"keepAny := j2gort.AsIface[JSType](rawAny)",
	} {
		if !strings.Contains(gotRhino, want) {
			t.Errorf("opt/rhino/rhino.go missing expected %q in:\n%s", want, gotRhino)
		}
	}

	gotJscompBytes, err := os.ReadFile(filepath.Join(tmpDir, "opt", "jscomp", "jscomp.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotJscomp := string(gotJscompBytes)
	for _, want := range []string{
		".EqualsJSType(b)",
		"fn1.GetJSTypeBase().EqualsJSType(pkg_rhino.ToIface_JSType(fn2))",
		"upCast := pkg_rhino.Upcast_JSType(fn2)",
		"keepNoUpcast := j2gort.AsIface[LegacyIface](leg)",
	} {
		if !strings.Contains(gotJscomp, want) {
			t.Errorf("opt/jscomp/jscomp.go missing expected %q in:\n%s", want, gotJscomp)
		}
	}

	diffTest := `package peep_diff_test

import (
	"testing"

	opt_jscomp "example.com/w38b_peep/opt/jscomp"
	opt_rhino "example.com/w38b_peep/opt/rhino"
	orig_jscomp "example.com/w38b_peep/orig/jscomp"
	orig_rhino "example.com/w38b_peep/orig/rhino"
)

func TestPeepholeParity(t *testing.T) {
	of1 := orig_rhino.NewFunctionType(10)
	of2 := orig_rhino.NewFunctionType(10)
	var oNilFn *orig_rhino.FunctionType = nil

	xf1 := opt_rhino.NewFunctionType(10)
	xf2 := opt_rhino.NewFunctionType(10)
	var xNilFn *opt_rhino.FunctionType = nil

	ob1, ou1, _, on1 := of1.MatchesOther(of2, of2, of1, of2)
	xb1, xu1, _, xn1 := xf1.MatchesOther(xf2, xf2, xf1, xf2)
	if ob1 != xb1 || (ou1 == nil) != (xu1 == nil) || (on1 == nil) != (xn1 == nil) {
		t.Fatalf("MatchesOther non-nil mismatch: orig=(%v,%v,%v) opt=(%v,%v,%v)", ob1, ou1, on1, xb1, xu1, xn1)
	}

	// Critical typed-nil check: passing nil *FunctionType MUST produce untyped nil interface (== nil)!
	ob2, ou2, _, on2 := of1.MatchesOther(oNilFn, of2, of1, nil)
	xb2, xu2, _, xn2 := xf1.MatchesOther(xNilFn, xf2, xf1, nil)
	if ob2 != xb2 || ou2 != nil || xu2 != nil || on2 != nil || xn2 != nil {
		t.Fatalf("MatchesOther nil *FunctionType typed-nil regression: orig=(%v,%v,%v) opt=(%v,%v,%v)", ob2, ou2, on2, xb2, xu2, xn2)
	}

	oj1, oup1, oleg1 := orig_jscomp.CheckTypesInJscomp(of1, of2, of1, oNilFn, nil)
	xj1, xup1, xleg1 := opt_jscomp.CheckTypesInJscomp(xf1, xf2, xf1, xNilFn, nil)
	if oj1 != xj1 || oup1 != nil || xup1 != nil || oleg1 != nil || xleg1 != nil {
		t.Fatalf("CheckTypesInJscomp nil mismatch: orig=(%v,%v,%v) opt=(%v,%v,%v)", oj1, oup1, oleg1, xj1, xup1, xleg1)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "peep_diff_test.go"), []byte(diffTest), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-count=1", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("differential runtime test failed: %v\n%s", err, string(out))
	}
}

func TestWave38B_PlainMapAndPlainSetLoweringAndParity(t *testing.T) {
	appSrc := `package app

import "MODROOT/j2gort"

type Entry struct {
	Label string
}

type LookupTracker struct {
	byName      j2gort.Ext_java_util_Map
	seenTags    *j2gort.Ext_java_util_HashSet
	orderedById j2gort.Ext_java_util_Map
	linkedCache *j2gort.Ext_java_util_LinkedHashMap
	linkedNames *j2gort.Ext_java_util_LinkedHashSet
}

func (this_ *LookupTracker) initFields_LookupTracker() {
	this_.byName = j2gort.NewExt_java_util_HashMap_1_int32(int32(16))
	this_.seenTags = j2gort.Ext_com_google_common_collect_Sets_NewHashSetWithExpectedSize(int32(8))
	this_.orderedById = j2gort.NewExt_java_util_HashMap()
	this_.linkedCache = j2gort.NewExt_java_util_LinkedHashMap()
	this_.linkedNames = j2gort.NewExt_java_util_LinkedHashSet()
}

func NewLookupTracker() *LookupTracker {
	t := &LookupTracker{}
	t.initFields_LookupTracker()
	return t
}

func (this_ *LookupTracker) RecordAndQuery(keys []string, vals []int32, probe string) (int32, bool, string) {
	this_.byName.Clear()
	this_.seenTags.Clear()

	localCounts := j2gort.NewExt_java_util_HashMap()
	localVisited := j2gort.NewExt_java_util_HashSet()

	for i, k := range keys {
		e := &Entry{Label: k}
		// Standalone ExprStmt Put_2 and Add_1 -> must lower to PutVoid and AddVoid!
		this_.byName.Put_2(j2gort.BoxString(k), j2gort.PtrToAny(e))
		this_.seenTags.Add_1(j2gort.BoxString(k))
		localCounts.Put_2(j2gort.BoxString(k), j2gort.ValueOfInt(vals[i]))
		localVisited.Add_1(j2gort.BoxString(k))

		this_.orderedById.Put_2(j2gort.BoxString(k), j2gort.PtrToAny(e))
		this_.linkedCache.Put_2(j2gort.BoxString(k), j2gort.PtrToAny(e))
		this_.linkedNames.Add_1(j2gort.BoxString(k))
	}

	// Non-void Put / Add where return value IS inspected must keep .Put(...) / .Add(...)!
	prev := j2gort.CastPtr[Entry](this_.byName.Put_2(j2gort.BoxString(probe), j2gort.PtrToAny(&Entry{Label: "probe"})))
	addedFresh := localVisited.Add_1(j2gort.BoxString(probe))

	// ComputeIfAbsent on PlainMap:
	comp := j2gort.CastPtr[Entry](this_.byName.ComputeIfAbsent(
		j2gort.BoxString("computed"),
		j2gort.NewLambda_Ext_java_util_function_Function(func(rawK any) any {
			ks := j2gort.CastString(rawK)
			return j2gort.PtrToAny(&Entry{Label: "gen:" + ks})
		}),
	))

	var score int32 = this_.byName.Size() + this_.seenTags.Size() + localVisited.Size()
	if localCounts.ContainsKey(j2gort.BoxString(probe)) {
		score += j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](localCounts.Get(j2gort.BoxString(probe))))
	}
	if prev != nil {
		score += int32(len(prev.Label))
	}
	if comp != nil && this_.seenTags.Contains(j2gort.BoxString(probe)) {
		_ = this_.seenTags.Remove_1_any(j2gort.BoxString(probe))
		score += int32(len(comp.Label))
	}

	// Iterate orderedById -> forces orderedById to stay OrderedMap!
	var orderTrace string
	{
		cur := j2gort.NewMapKeyStrCursor(this_.orderedById)
		for cur.HasNext() {
			k := cur.NextStr()
			_ = k
			orderTrace += k + ","
		}
	}
	if !this_.linkedCache.IsEmpty() && !this_.linkedNames.IsEmpty() {
		orderTrace += "L"
	}

	return score, addedFresh, orderTrace
}
`

	tmpDir := t.TempDir()
	writeSet := func(prefix string) {
		modImport := "example.com/w38b_plain/" + prefix
		mustWrite := func(rel, content string) {
			t.Helper()
			full := filepath.Join(tmpDir, prefix, rel)
			if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(full, []byte(strings.ReplaceAll(content, "MODROOT", modImport)), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		mustWrite("j2gort/runtime.go", wave38bMockJ2gortSrc)
		mustWrite("app/app.go", appSrc)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com/w38b_plain\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeSet("orig")
	writeSet("opt")

	optAppDir := filepath.Join(tmpDir, "opt", "app")
	allPasses := passes.AllPasses()

	s1, err := passes.RunOnPackageDir(optAppDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.ByPass["native_map_lowering"] == 0 || s1.ByPass["native_set_lowering"] == 0 {
		t.Fatalf("expected non-zero map & set lowerings on Run 1, got %+v", s1.ByPass)
	}

	// Run 2 MUST produce strictly 0 sites and 0 modified files!
	s2, err := passes.RunOnPackageDir(optAppDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d ByPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	gotBytes, err := os.ReadFile(filepath.Join(optAppDir, "app.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)

	for _, want := range []string{
		"byName      j2gort.PlainMap[string, *Entry]",
		"seenTags    j2gort.PlainSet[string]",
		"orderedById j2gort.OrderedMap[string, *Entry]",
		"linkedCache j2gort.OrderedMap[string, *Entry]",
		"linkedNames j2gort.OrderedSet[string]",
		"j2gort.NewPlainMap[string, *Entry](16)",
		"j2gort.NewPlainSet[string](8)",
		"j2gort.NewPlainMap[string, int32](0)",
		"j2gort.NewPlainSet[string](0)",
		"this_.byName.PutVoid(k, e)",
		"this_.seenTags.AddVoid(k)",
		"localCounts.PutVoid(k, vals[i])",
		"localVisited.AddVoid(k)",
		"this_.byName.Put(probe, &Entry{Label: \"probe\"})",
		"localVisited.Add(probe)",
		"j2gort.PlainMapComputeIfAbsent(&this_.byName,",
		"localCounts.GetOk(probe)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("opt/app/app.go missing expected %q in:\n%s", want, got)
		}
	}

	diffTest := `package plain_diff_test

import (
	"testing"

	"example.com/w38b_plain/opt/app"
	orig "example.com/w38b_plain/orig/app"
)

func TestPlainMapAndPlainSetRuntimeParity(t *testing.T) {
	ot := orig.NewLookupTracker()
	xt := app.NewLookupTracker()

	cases := []struct {
		keys  []string
		vals  []int32
		probe string
	}{
		{[]string{"a", "b", "c"}, []int32{10, 20, 30}, "b"},
		{[]string{"x", "y"}, []int32{-5, 99}, "missing"},
		{[]string{"dup", "dup"}, []int32{1, 7}, "dup"},
	}
	for i, tc := range cases {
		os, oa, otr := ot.RecordAndQuery(tc.keys, tc.vals, tc.probe)
		xs, xa, xtr := xt.RecordAndQuery(tc.keys, tc.vals, tc.probe)
		if os != xs || oa != xa || otr != xtr {
			t.Fatalf("case %d mismatch: orig=(%d,%v,%q) vs opt=(%d,%v,%q)", i, os, oa, otr, xs, xa, xtr)
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "plain_diff_test.go"), []byte(diffTest), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-race", "-count=1", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("differential runtime test failed: %v\n%s", err, string(out))
	}
}
